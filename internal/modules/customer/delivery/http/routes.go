package http

import (
	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/customer"

	"github.com/go-chi/chi/v5"
)

func MapCustomerRoute(router *chi.Mux, h customer.Handlers, mw *middleware.MiddlewareManager) {
	router.Route("/customer", func(r chi.Router) {
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
			})
		})
	})

	router.Route("/car", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())
			r.Get("/", h.ListCars())
			r.Post("/", h.CreateCar())
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.GetCar())
				r.Put("/", h.UpdateCar())
				r.Delete("/", h.DeleteCar())
			})
		})
	})
}
