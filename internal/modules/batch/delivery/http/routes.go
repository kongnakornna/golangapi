package http

import (
	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/batch"

	"github.com/go-chi/chi/v5"
)

// MapBatchRoute registers batch job routes.
// กำหนดเส้นทาง API สำหรับงานแบตช์
func MapBatchRoute(router *chi.Mux, h batch.Handlers, mw *middleware.MiddlewareManager) {
	router.Route("/api/batch", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())
			r.Post("/jobs", h.CreateJob())
			r.Get("/jobs", h.ListJobs())
			r.Route("/jobs/{id}", func(r chi.Router) {
				r.Get("/", h.GetJob())
				r.Put("/", h.UpdateJob())
				r.Delete("/", h.DeleteJob())
				r.Post("/run", h.RunJobNow())
				r.Get("/logs", h.GetJobLogs())
			})
		})
	})
}
