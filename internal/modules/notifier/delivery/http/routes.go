package http

import (
	"icmongolang/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// MapNotifierRoutes registers notifier routes under /api/notifier.
func MapNotifierRoutes(router chi.Router, h *NotifierHandler, mw *middleware.MiddlewareManager) {
	router.Route("/notifier", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())

			r.Get("/health", h.Health)
			r.Post("/dispatch", h.Dispatch)
		})
	})
}
