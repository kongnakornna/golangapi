package http

import (
	"icmongolang/internal/middleware"
	"icmongolang/internal/template"

	"github.com/go-chi/chi/v5"
)

func MapTemplateRoutes(router chi.Router, h template.Handlers, mw *middleware.MiddlewareManager) {
	router.Route("/templates", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true), mw.Authenticator(), mw.CurrentUser(), mw.ActiveUser())

			r.Get("/", h.GetMulti())
			r.Get("/{id}", h.Get())
			r.Post("/", h.Create())
			r.Put("/{id}", h.Update())
			r.Delete("/{id}", h.Delete())
		})
	})
}
