package http

import (
	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/document"

	"github.com/go-chi/chi/v5"
)

// MapDocumentRoute registers document routes.
// กำหนดเส้นทาง API สำหรับเอกสาร
func MapDocumentRoute(router *chi.Mux, h document.Handlers, mw *middleware.MiddlewareManager) {
	router.Route("/api/documents", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())
			r.Post("/", h.Upload())
			r.Get("/", h.List())
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.Download())
				r.Delete("/", h.Delete())
			})
		})
	})
}
