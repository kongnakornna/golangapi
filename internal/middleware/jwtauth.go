package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"icmongolang/internal/models"
	"icmongolang/internal/modules/auth/domain"
	"icmongolang/pkg/apikey"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/jwe"
	"icmongolang/pkg/jwt"
	"icmongolang/pkg/responses"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// ─── context keys ─────────────────────────────────────────────
var (
	TokenCtxKey         = &contextKey{"Token"}
	IdCtxKey            = &contextKey{"Id"}
	EmailCtxKey         = &contextKey{"Email"}
	ErrorCtxKey         = &contextKey{"Error"}
	UserCtxKey          = &contextKey{"User"}
	DownloadTokenCtxKey = &contextKey{"DownloadToken"}
	ClaimsCtxKey        = &contextKey{"Claims"}
)

type contextKey struct{ name string }

func (k *contextKey) String() string { return "jwtauth context value " + k.name }

// ─── Claims ───────────────────────────────────────────────────
// Claims ที่ handler ฝั่ง HTTP อ่าน
type Claims struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Email     string      `json:"email"`
	Role      domain.Role `json:"role"`
	SessionID string      `json:"sid"`
}

// ClaimsFromContext อ่าน claims จาก context
func ClaimsFromContext(ctx context.Context) *Claims {
	if c, ok := ctx.Value(ClaimsCtxKey).(*Claims); ok {
		return c
	}
	return nil
}

// ─── Verifier (RS256 JWT) ─────────────────────────────────────
func (mw *MiddlewareManager) Verifier(requireAccessToken bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token := TokenFromRequest(r)

			if token == "" {
				err := httpErrors.ErrTokenNotFound(errors.New("token not found"))
				ctx = context.WithValue(ctx, ErrorCtxKey, err)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			var publicKey string
			if requireAccessToken {
				publicKey = mw.cfg.Jwt.AccessTokenPublicKey
			} else {
				publicKey = mw.cfg.Jwt.RefreshTokenPublicKey
			}

			id, email, err := jwt.ParseTokenRS256(token, publicKey)
			ctx = context.WithValue(ctx, TokenCtxKey, token)
			ctx = context.WithValue(ctx, IdCtxKey, id)
			ctx = context.WithValue(ctx, EmailCtxKey, email)
			ctx = context.WithValue(ctx, ErrorCtxKey, err)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ─── VerifierJWE (Ed25519 + ECDH-ES) ──────────────────────────
// ⭐ สำคัญ: set ClaimsCtxKey ให้ downstream (RequireRole) ใช้ได้
func (mw *MiddlewareManager) VerifierJWE() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token := TokenFromRequest(r)

			if token == "" {
				ctx = context.WithValue(ctx, ErrorCtxKey,
					httpErrors.ErrTokenNotFound(errors.New("token not found")))
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if mw.jweKP == nil {
				ctx = context.WithValue(ctx, ErrorCtxKey,
					httpErrors.ErrInvalidJWTToken(errors.New("JWE keys not configured")))
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			claims, err := jwe.OpenFromCookie(mw.jweKP, token)
			if err != nil {
				ctx = context.WithValue(ctx, ErrorCtxKey,
					httpErrors.ErrInvalidJWTToken(err))
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// user-level revoke
			if mw.redisRepo != nil {
				userRevokeKey := fmt.Sprintf("revoke:user:%s", claims.UserID)
				if v, err := mw.redisRepo.GetString(ctx, userRevokeKey); err == nil && v != "" {
					ctx = context.WithValue(ctx, ErrorCtxKey,
						httpErrors.ErrNotFoundRefreshTokenRedis(errors.New("user sessions revoked")))
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			// token-level revoke (HMAC fingerprint)
			if len(mw.apiKeyPepper) > 0 {
				fp := apikey.Hash(mw.apiKeyPepper, token)
				revokeKey := fmt.Sprintf("revoke:refresh:%s", fp)
				if v, err := mw.redisRepo.GetString(ctx, revokeKey); err == nil && v != "" {
					ctx = context.WithValue(ctx, ErrorCtxKey,
						httpErrors.ErrNotFoundRefreshTokenRedis(errors.New("token revoked")))
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			// ── map role ──
			role := domain.Role(claims.Role)
			if !role.IsValid() {
				role = domain.RoleUser
			}

			// ── set context ──
			ctx = context.WithValue(ctx, TokenCtxKey, token)
			ctx = context.WithValue(ctx, IdCtxKey, claims.UserID)
			ctx = context.WithValue(ctx, EmailCtxKey, claims.Email)
			ctx = context.WithValue(ctx, ErrorCtxKey, nil)
			ctx = context.WithValue(ctx, ClaimsCtxKey, &Claims{
				ID:        claims.UserID,
				UserID:    claims.UserID,
				Email:     claims.Email,
				Role:      role,
				SessionID: claims.SessionID,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ─── VerifierForReport ────────────────────────────────────────
func (mw *MiddlewareManager) VerifierForReport() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if token := tokenFromHeaderOrCookie(r); token != "" {
				id, email, err := jwt.ParseTokenRS256(token, mw.cfg.Jwt.AccessTokenPublicKey)
				ctx = context.WithValue(ctx, TokenCtxKey, token)
				ctx = context.WithValue(ctx, IdCtxKey, id)
				ctx = context.WithValue(ctx, EmailCtxKey, email)
				ctx = context.WithValue(ctx, ErrorCtxKey, err)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if dt := r.URL.Query().Get("token"); dt != "" {
				if mw.jweKP != nil {
					claims, err := jwt.ParseDownloadTokenJWE(mw.jweKP, dt)
					if err == nil {
						ctx = context.WithValue(ctx, DownloadTokenCtxKey, claims)
						ctx = context.WithValue(ctx, TokenCtxKey, dt)
						ctx = context.WithValue(ctx, IdCtxKey, claims.Id)
						ctx = context.WithValue(ctx, EmailCtxKey, "")
						ctx = context.WithValue(ctx, ErrorCtxKey, nil)
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
				claims, err := jwt.ParseDownloadTokenHS256(dt, mw.cfg.Jwt.DownloadTokenSecret)
				if err != nil {
					ctx = context.WithValue(ctx, ErrorCtxKey, err)
					ctx = context.WithValue(ctx, IdCtxKey, "")
				} else {
					ctx = context.WithValue(ctx, DownloadTokenCtxKey, claims)
					ctx = context.WithValue(ctx, TokenCtxKey, dt)
					ctx = context.WithValue(ctx, IdCtxKey, claims.Id)
					ctx = context.WithValue(ctx, EmailCtxKey, "")
					ctx = context.WithValue(ctx, ErrorCtxKey, nil)
				}
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			ctx = context.WithValue(ctx, ErrorCtxKey,
				httpErrors.ErrTokenNotFound(errors.New("token not found")))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ─── Authenticator ────────────────────────────────────────────
// ตรวจ error ที่ Verifier เก็บไว้ → 401 ถ้ามี
func (mw *MiddlewareManager) Authenticator() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err, _ := r.Context().Value(ErrorCtxKey).(error); err != nil {
				_ = render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ─── CurrentUser ──────────────────────────────────────────────
func (mw *MiddlewareManager) CurrentUser() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			id, _ := ctx.Value(IdCtxKey).(string)
			err, _ := ctx.Value(ErrorCtxKey).(error)

			if err != nil || id == "" {
				_ = render.Render(w, r, responses.CreateErrorResponse(httpErrors.ParseErrors(err)))
				return
			}

			idParsed, err := uuid.Parse(id)
			if err != nil {
				_ = render.Render(w, r, responses.CreateErrorResponse(
					httpErrors.ErrInvalidJWTClaims(errors.New("cannot convert id to uuid"))))
				return
			}

			user, err := mw.usersUC.Get(ctx, idParsed)
			if err != nil {
				_ = render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}

			ctx = context.WithValue(ctx, UserCtxKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ─── SuperUser / ActiveUser ───────────────────────────────────
func (mw *MiddlewareManager) SuperUser() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user, err := GetUserFromCtx(ctx)
			if err != nil {
				_ = render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
			if !mw.usersUC.IsSuper(ctx, *user) {
				_ = render.Render(w, r, responses.CreateErrorResponse(
					httpErrors.ErrNotEnoughPrivileges(errors.New("not super user"))))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (mw *MiddlewareManager) ActiveUser() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user, err := GetUserFromCtx(ctx)
			if err != nil {
				_ = render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
			if !mw.usersUC.IsActive(ctx, *user) {
				_ = render.Render(w, r, responses.CreateErrorResponse(
					httpErrors.ErrInactiveUser(errors.New("user inactive"))))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ─── Token extractors ─────────────────────────────────────────
const AccessTokenCookieName = "access_token"
const RefreshTokenCookieName = "refresh_token"

func TokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

func TokenFromCookie(r *http.Request, name string) string {
	if c, err := r.Cookie(name); err == nil && c.Value != "" {
		return c.Value
	}
	return ""
}

func tokenFromHeaderOrCookie(r *http.Request) string {
	if t := TokenFromHeader(r); t != "" {
		return t
	}
	return TokenFromCookie(r, AccessTokenCookieName)
}

func RefreshTokenFromRequest(r *http.Request) string {
	if t := TokenFromHeader(r); t != "" {
		return t
	}
	return TokenFromCookie(r, RefreshTokenCookieName)
}

func TokenFromRequest(r *http.Request) string {
	if t := tokenFromHeaderOrCookie(r); t != "" {
		return t
	}
	if t := r.URL.Query().Get("access_token"); t != "" {
		return t
	}
	return r.URL.Query().Get("token")
}

func GetUserFromCtx(ctx context.Context) (*models.SdUser, error) {
	user, ok := ctx.Value(UserCtxKey).(*models.SdUser)
	if !ok {
		return nil, httpErrors.ErrUnauthorized(errors.New("cannot convert user from context"))
	}
	return user, nil
}