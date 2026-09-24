package http

import (
	"icmongolang/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// MapAPIManagerRoutes registers API-manager routes under /api/apimanager.
func MapAPIManagerRoutes(router chi.Router, h *APIManagerHandler, mw *middleware.MiddlewareManager) {
	router.Route("/apimanager", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())

			r.Post("/key", h.CreateKey)
			r.Delete("/key/{key_id}", h.RevokeKey)
			r.Get("/ratelimit/{key_id}", h.CheckRateLimit)
			r.Get("/usage/{key_id}", h.Usage)
		})
	})
}
