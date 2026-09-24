package http

import (
	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/dashboard"

	"github.com/go-chi/chi/v5"
)

// MapDashboardRoute registers dashboard routes.
// กำหนดเส้นทาง API สำหรับแดชบอร์ด
func MapDashboardRoute(router *chi.Mux, h dashboard.Handlers, mw *middleware.MiddlewareManager) {
	router.Route("/api/dashboard", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())
			r.Get("/stats", h.GetDashboardStats())
			r.Get("/revenue", h.GetRevenueChart())
			r.Get("/top-parts", h.GetTopParts())
			r.Get("/job-status", h.GetJobStatusSummary())
		})
	})
}
