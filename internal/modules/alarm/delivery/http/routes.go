package http

import (
	"icmongolang/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func MapAlarmRoutes(router chi.Router, h *AlarmHandler, mw *middleware.MiddlewareManager) {
	router.Route("/alarm", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true), mw.Authenticator(), mw.CurrentUser(), mw.ActiveUser())
			r.Post("/validate", h.ValidateAlarm)
			r.Post("/validate/en", h.ValidateAlarmEn)
			r.Post("/validate/th", h.ValidateAlarmTh)
		})
	})
}
