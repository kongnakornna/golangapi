package http

import (
	"github.com/go-chi/chi/v5"

	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/auth/domain"
	"icmongolang/internal/modules/keys"
)

func MapRoutes(r chi.Router, h keys.Handlers, mw *middleware.MiddlewareManager) {
	r.Route("/api-keys", func(r chi.Router) {
		// ⭐ ใช้ VerifierJWE (set ClaimsCtxKey) ไม่ใช่ Verifier(true)
		// ถ้า login ออก RS256 ให้ใช้ Verifier(true) + WithClaims() แทน
		r.Use(mw.VerifierJWE())
		r.Use(mw.Authenticator())
		r.Use(mw.RequireRole(domain.RoleManager))

		r.Post("/", h.CreateAPIKey())
		r.Get("/", h.ListAPIKeys())
		r.Get("/{id}", h.GetAPIKey())
		r.Post("/{id}/rotate", h.RotateAPIKey())
		r.Delete("/{id}", h.RevokeAPIKey())
	})
}