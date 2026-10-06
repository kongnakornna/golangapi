package http

import (
	"time"

	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/auth"

	"github.com/go-chi/chi/v5"
)

// MapAuthRoute registers all /auth routes on the given router.
//
// Routes are grouped into three tiers:
//
//	Tier 1 — Public (rate limited, no auth)
//	Tier 2 — Refresh-token only (uses refresh cookie, no access-token verify)
//	Tier 3 — Protected (access token required)
func MapAuthRoute(router *chi.Mux, h auth.Handlers, mw *middleware.MiddlewareManager) {
	rateLimitConfig := middleware.RateLimitConfig{
		RequestsPerSecond: 10,
		Burst:             20,
		CleanupInterval:   15 * time.Minute,
	}
	rateLimiter := middleware.NewRateLimitMiddleware(rateLimitConfig)

	router.Route("/auth", func(r chi.Router) {
		// Tier 1: Public
		r.Group(func(r chi.Router) {
			r.Use(rateLimiter.Handler)

			r.Post("/login", h.Login())
			r.Post("/signin", h.SignIn())
			r.Get("/publickey", h.GetPublicKey())
			r.Get("/verifyemail", h.VerifyEmail())
			r.Post("/forgotpassword", h.ForgotPassword())
			r.Patch("/resetpassword", h.ResetPassword())
		})

		// Tier 2: Refresh-token routes — no Verifier / Authenticator
		r.Group(func(r chi.Router) {
			r.Use(rateLimiter.Handler)

			r.Get("/refresh", h.RefreshToken())
			r.Post("/logout", h.Logout())
			r.Get("/logout", h.LogoutAllToken())
			r.Post("/logoutall", h.LogoutAllToken())
			r.Get("/logoutall", h.LogoutAllToken())
		})

		// Tier 3: Protected routes — add authenticated endpoints here
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(false))
			r.Use(mw.Authenticator()) 
		})
	})
}
