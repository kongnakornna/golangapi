package http

import (
	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/job"

	"github.com/go-chi/chi/v5"
)

// MapJobRoute registers Job module routes.
// กำหนดเส้นทาง API สำหรับโมดูลใบรับงานซ่อม
func MapJobRoute(router *chi.Mux, h job.Handlers, mw *middleware.MiddlewareManager) {
	router.Route("/job", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())
			r.Get("/", h.List())
			r.Post("/", h.Create())
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.GetByID())
				r.Put("/", h.Update())
				r.Delete("/", h.Delete())
				r.Put("/status", h.ChangeStatus())
				r.Get("/history", h.GetStatusHistory())
				r.Get("/report", h.GetReport())
				r.Get("/pdf", h.GetPDF())
				r.Get("/picking/pdf", h.GetPickingPDF())
				r.Get("/delivery/pdf", h.GetDeliveryPDF())
				r.Post("/services", h.AddService())
				r.Get("/services", h.GetServices())
				r.Post("/parts", h.AddPart())
				r.Get("/parts", h.GetParts())
			})
		})
	})
}
