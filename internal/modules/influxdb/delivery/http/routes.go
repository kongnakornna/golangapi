package http

import (
	"icmongolang/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// MapInfluxRoutes ลงทะเบียน routes ของ InfluxDB ภายใต้ /api/influx
func MapInfluxRoutes(router chi.Router, h *InfluxHandler, mw *middleware.MiddlewareManager) {
	router.Route("/influx", func(r chi.Router) {
		// Routes ที่ต้องการ authentication (ทั้งหมด)
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())

			r.Post("/write", h.WriteData)
			r.Get("/query", h.QueryGetFilter)
			r.Post("/devicechart", h.Querydevicechart)
			r.Post("/filters", h.QueryFilters)
			r.Post("/statistics", h.QueryStatistics)
		})
	})
}
