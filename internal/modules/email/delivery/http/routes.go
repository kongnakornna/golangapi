package http

import (
	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/email"

	"github.com/go-chi/chi/v5"
)

// MapEmailRoute registers email routes.
// กำหนดเส้นทาง API สำหรับอีเมล
func MapEmailRoute(router *chi.Mux, h email.Handlers, mw *middleware.MiddlewareManager) {
	router.Route("/api/email", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())
			r.Post("/send", h.SendEmail())
			r.Get("/logs", h.ListEmailLogs())
			r.Get("/logs/{id}", h.GetEmailLog())
			r.Get("/config", h.GetEmailConfig())
			r.Put("/config", h.UpdateEmailConfig())
		})
	})
}
