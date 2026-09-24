package http

import (
	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/wos"

	"github.com/go-chi/chi/v5"
)

// MapWosRoute registers Web Order System routes.
// กำหนดเส้นทาง API สำหรับระบบสั่งซื้อออนไลน์
func MapWosRoute(router *chi.Mux, h wos.Handlers, mw *middleware.MiddlewareManager) {
	router.Route("/api/wos", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())
			r.Post("/orders", h.CreateOrder())
			r.Get("/orders", h.ListOrders())
			r.Route("/orders/{id}", func(r chi.Router) {
				r.Get("/", h.GetOrder())
				r.Put("/status", h.UpdateOrderStatus())
			})
		})
	})
}
