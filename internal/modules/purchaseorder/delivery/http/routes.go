package http

import (
	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/purchaseorder"

	"github.com/go-chi/chi/v5"
)

func MapPurchaseOrderRoute(router chi.Router, h purchaseorder.Handlers, mw *middleware.MiddlewareManager) {
	router.Route("/purchase-orders", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())

			r.Get("/", h.List())
			r.Post("/", h.Create())
			r.Post("/from-quotation/{quotationId}", h.CreateFromQuotation())
			r.Get("/suggestions/{jobId}", h.GetSuggestions())

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.GetByID())
				r.Put("/", h.Update())
				r.Delete("/", h.Delete())
				r.Post("/send", h.Send())
				r.Put("/confirm", h.Confirm())
				r.Post("/receive", h.Receive())
				r.Put("/cancel", h.Cancel())
				r.Get("/pdf", h.GetPDF())
				r.Get("/history", h.GetStatusHistory())
			})
		})
	})
}
