package http

import (
	"icmongolang/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// MapFlowEngineRoutes registers flow engine routes under /api/flowengine.
func MapFlowEngineRoutes(router chi.Router, h *FlowEngineHandler, mw *middleware.MiddlewareManager) {
	router.Route("/flowengine", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())

			r.Post("/flow", h.Create)
			r.Get("/flow/{id}", h.Get)
			r.Post("/validate", h.Validate)
			r.Post("/trigger", h.Trigger)
		})
	})
}
