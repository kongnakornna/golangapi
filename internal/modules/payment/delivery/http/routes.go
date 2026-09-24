package http

import (
	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/payment"

	"github.com/go-chi/chi/v5"
)

func MapPaymentRoute(router *chi.Mux, h payment.Handlers, mw *middleware.MiddlewareManager) {
	router.Route("/payments", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())

			r.Post("/", h.Create())
			r.Get("/", h.List())
			r.Post("/search", h.Search())
			r.Get("/outstanding/{customerId}", h.GetOutstanding())
			r.Get("/history/{customerId}", h.GetHistory())

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.Get())
				r.Post("/refund", h.Refund())
				r.Put("/cancel", h.Cancel())
			})

			r.Get("/invoice/{invoiceId}", h.GetByInvoice())
		})
	})
}

func MapReceiptRoute(router *chi.Mux, h payment.Handlers, mw *middleware.MiddlewareManager) {
	router.Route("/receipts", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.GetReceipt())
				r.Get("/pdf", h.GetReceiptPDF())
				r.Put("/cancel", h.CancelReceipt())
			})

			r.Get("/payment/{paymentId}", h.GetReceiptByPayment())
		})
	})
}
