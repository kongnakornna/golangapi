package http

import (
	"icmongolang/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// MapRealtimeRoutes registers realtime routes under /api/realtime.
func MapRealtimeRoutes(router chi.Router, h *RealtimeHandler, mw *middleware.MiddlewareManager) {
	router.Route("/realtime", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())

			r.Get("/health", h.Health)
			r.Post("/publish", h.Publish)
		})
	})
}
