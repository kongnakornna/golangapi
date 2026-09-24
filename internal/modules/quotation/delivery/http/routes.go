package http

import (
	"icmongolang/internal/modules/quotation"
	"icmongolang/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func MapQuotationRoute(router *chi.Mux, h quotation.Handlers, mw *middleware.MiddlewareManager) {
	router.Route("/quotation", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())
			r.Get("/", h.GetMulti())
			r.Post("/", h.Create())
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.Get())
				r.Delete("/", h.Delete())
				r.Put("/", h.Update())
				r.Get("/pdf", h.GetPDF())
			})
		})
	})
}
