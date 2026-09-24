package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"icmongolang/internal/models"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/jwt"
	"icmongolang/pkg/responses"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

var (
	TokenCtxKey         = &contextKey{"Token"}
	IdCtxKey            = &contextKey{"Id"}
	EmailCtxKey         = &contextKey{"Email"}
	ErrorCtxKey         = &contextKey{"Error"}
	UserCtxKey          = &contextKey{"User"}
	DownloadTokenCtxKey = &contextKey{"DownloadToken"}
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "jwtauth context value " + k.name
}

func (mw *MiddlewareManager) Verifier(requireAccessToken bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token := TokenFromRequest(r)
			if token == "" {
				err := httpErrors.ErrTokenNotFound(errors.New("not found token in header"))
				ctx = context.WithValue(ctx, ErrorCtxKey, err)
			} else {
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
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// VerifierForReport authenticates /reports/*/pdf requests via either the
// Authorization header (RS256 access token) or a short-lived ?token= download
// token. The download token lets a browser navigate straight to a PDF link
// without exposing the long-lived access token in the URL.
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

			ctx = context.WithValue(ctx, ErrorCtxKey, httpErrors.ErrTokenNotFound(errors.New("not found token in header")))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (mw *MiddlewareManager) Authenticator() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err, _ := r.Context().Value(ErrorCtxKey).(error)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (mw *MiddlewareManager) CurrentUser() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			id, _ := ctx.Value(IdCtxKey).(string)
			err, _ := ctx.Value(ErrorCtxKey).(error)

			if err != nil || id == "" {
				render.Render(w, r, responses.CreateErrorResponse(httpErrors.ParseErrors(err)))
				return
			}

			idParsed, err := uuid.Parse(id)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrInvalidJWTClaims(errors.New("can not convert id to uuid from id in token"))))
				return
			}

			user, err := mw.usersUC.Get(ctx, idParsed)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}

			ctx = context.WithValue(ctx, UserCtxKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (mw *MiddlewareManager) SuperUser() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user, err := GetUserFromCtx(ctx)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
			if !mw.usersUC.IsSuper(ctx, *user) {
				render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotEnoughPrivileges(errors.New("user is not super user"))))
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
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
			if !mw.usersUC.IsActive(ctx, *user) {
				render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrInactiveUser(errors.New("user inactive"))))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func TokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

// AccessTokenCookieName is the cookie name checked when the Authorization
// header is absent, so cookie-based clients can authenticate without changes.
const AccessTokenCookieName = "access_token"

// RefreshTokenCookieName is the cookie used for the refresh token on login.
const RefreshTokenCookieName = "refresh_token"

// TokenFromCookie returns the value of the named cookie, or "" when absent.
func TokenFromCookie(r *http.Request, name string) string {
	if c, err := r.Cookie(name); err == nil && c.Value != "" {
		return c.Value
	}
	return ""
}

// tokenFromHeaderOrCookie returns the token from the Authorization header or,
// as a fallback, from the access_token cookie.
func tokenFromHeaderOrCookie(r *http.Request) string {
	if token := TokenFromHeader(r); token != "" {
		return token
	}
	return TokenFromCookie(r, AccessTokenCookieName)
}

// RefreshTokenFromRequest returns the refresh token from the Authorization
// header or, as a fallback, from the refresh_token cookie.
func RefreshTokenFromRequest(r *http.Request) string {
	if token := TokenFromHeader(r); token != "" {
		return token
	}
	return TokenFromCookie(r, RefreshTokenCookieName)
}

// TokenFromRequest returns the access token from the Authorization header, the
// access_token cookie, or the access_token/token query parameter, in that order
// of preference.
func TokenFromRequest(r *http.Request) string {
	if token := tokenFromHeaderOrCookie(r); token != "" {
		return token
	}
	if token := r.URL.Query().Get("access_token"); token != "" {
		return token
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
