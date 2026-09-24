package http

import (
	"icmongolang/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// MapAuditLogRoutes registers audit log routes under /api/auditlog.
func MapAuditLogRoutes(router chi.Router, h *AuditLogHandler, mw *middleware.MiddlewareManager) {
	router.Route("/auditlog", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())

			r.Post("/", h.Create)
			r.Get("/verify/{id}", h.Verify)
			r.Get("/chain", h.Chain)
		})
	})
}
