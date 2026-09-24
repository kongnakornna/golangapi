package http

import (
	"icmongolang/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// MapControlRoutes registers control routes under /api/control.
func MapControlRoutes(router chi.Router, h *ControlHandler, mw *middleware.MiddlewareManager) {
	router.Route("/control", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())

			r.Get("/health", h.Health)
			r.Post("/execute", h.ExecuteControl)
			r.Post("/schedule/run", h.RunSchedule)
		})
	})
}
