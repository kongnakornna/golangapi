package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"icmongolang/internal/modules/auth/domain"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/responses"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// ─── RequireRole ──────────────────────────────────────────────
// Chain AFTER Verifier + Authenticator + WithClaims
func (mw *MiddlewareManager) RequireRole(min domain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := ClaimsFromContext(r.Context())
			if claims == nil {
				_ = render.Render(w, r,
					responses.CreateErrorResponse(httpErrors.ErrUnauthorized(errors.New("missing claims"))))
				return
			}
			if !claims.Role.AtLeast(min) {
				_ = render.Render(w, r,
					responses.CreateErrorResponse(httpErrors.ErrNotEnoughPrivileges(errors.New("insufficient role"))))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ─── RequireAnyRole ───────────────────────────────────────────
func (mw *MiddlewareManager) RequireAnyRole(roles ...domain.Role) func(http.Handler) http.Handler {
	allowed := make(map[domain.Role]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := ClaimsFromContext(r.Context())
			if claims == nil {
				_ = render.Render(w, r,
					responses.CreateErrorResponse(httpErrors.ErrUnauthorized(errors.New("missing claims"))))
				return
			}
			if _, ok := allowed[claims.Role]; !ok {
				_ = render.Render(w, r,
					responses.CreateErrorResponse(httpErrors.ErrNotEnoughPrivileges(errors.New("role not allowed"))))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ─── RequireSuperUser ─────────────────────────────────────────
func (mw *MiddlewareManager) RequireSuperUser() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := GetUserFromCtx(r.Context())
			if err != nil {
				_ = render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrUnauthorized(err)))
				return
			}
			if !user.IsSuperUser {
				_ = render.Render(w, r,
					responses.CreateErrorResponse(httpErrors.ErrNotEnoughPrivileges(errors.New("superuser required"))))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ─── WithClaims ───────────────────────────────────────────────
// ใช้เฉพาะกรณีใช้ Verifier(RS256) เท่านั้น
// ถ้าใช้ VerifierJWE() อยู่แล้ว ไม่ต้องใช้ WithClaims
func (mw *MiddlewareManager) WithClaims() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			id, _ := ctx.Value(IdCtxKey).(string)
			email, _ := ctx.Value(EmailCtxKey).(string)
			authErr, _ := ctx.Value(ErrorCtxKey).(error)

			if authErr != nil || id == "" {
				next.ServeHTTP(w, r)
				return
			}

			idParsed, err := uuid.Parse(id)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			user, err := mw.usersUC.Get(ctx, idParsed)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			claims := &Claims{
				ID:     id,
				UserID: id,
				Email:  email,
				Role:   domain.RoleFromID(user.RoleID),
			}

			ctx = context.WithValue(ctx, ClaimsCtxKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ─── WS origin guard ──────────────────────────────────────────
func (mw *MiddlewareManager) WSOriginGuard(allowedOrigins []string) func(http.HandlerFunc) http.HandlerFunc {
	return WSOriginGuardFunc(allowedOrigins)
}

func WSOriginGuardFunc(allowedOrigins []string) func(http.HandlerFunc) http.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowed[strings.TrimRight(o, "/")] = struct{}{}
	}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			origin := strings.TrimRight(r.Header.Get("Origin"), "/")
			if origin == "" {
				http.Error(w, "missing origin", http.StatusForbidden)
				return
			}
			if _, ok := allowed[origin]; !ok {
				http.Error(w, "origin not allowed", http.StatusForbidden)
				return
			}
			next(w, r)
		}
	}
}