package http

import (
	"icmongolang/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// MapVectorDataRoutes registers routes under /api/vectordata.
func MapVectorDataRoutes(router chi.Router, h *VectorDataHandler, mw *middleware.MiddlewareManager) {
	router.Route("/vectordata", func(r chi.Router) {
		r.Get("/health", h.Health)

		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())

			r.Post("/create-index", h.CreateIndex)
			r.Post("/documents", h.Create)
			r.Get("/documents", h.List)
			r.Get("/documents/{id}", h.Get)
			r.Put("/documents/{id}", h.Update)
			r.Delete("/documents/{id}", h.Delete)
			r.Post("/search", h.Search)
			r.Post("/seed", h.Seed)
		})
	})
}
