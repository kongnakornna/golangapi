package http

import (
	"icmongolang/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// MapElasticsearchRoutes registers routes under /api/elasticsearch.
func MapElasticsearchRoutes(router chi.Router, h *ElasticsearchHandler, mw *middleware.MiddlewareManager) {
	router.Route("/elasticsearch", func(r chi.Router) {
		r.Get("/health", h.Health)

		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())

			r.Post("/create-index", h.CreateIndex)
			r.Post("/embeddings", h.EmbedAndIndex)
			r.Post("/vectors", h.IndexVector)
			r.Post("/search", h.SearchVector)
		})
	})
}
