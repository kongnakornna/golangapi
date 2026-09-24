package http

import (
	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/i18n"

	"github.com/go-chi/chi/v5"
)

// MapI18nRoute registers i18n routes.
// กำหนดเส้นทาง API สำหรับการแปลภาษา
func MapI18nRoute(router *chi.Mux, h i18n.Handlers, mw *middleware.MiddlewareManager) {
	router.Route("/api/i18n", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())
			r.Get("/translations", h.GetTranslations())
			r.Post("/translations", h.CreateTranslation())
			r.Route("/translations/{key}", func(r chi.Router) {
				r.Get("/", h.GetTranslationByKey())
				r.Put("/", h.UpdateTranslation())
				r.Delete("/", h.DeleteTranslation())
			})
		})
	})
}
