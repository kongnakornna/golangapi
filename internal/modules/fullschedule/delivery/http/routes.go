package http

import (
	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/fullschedule"

	"github.com/go-chi/chi/v5"
)

// MapMasterRoutes mounts the Group/Zone/Area master data routes.
// จดทะเบียนเส้นทาง Master Data (Group/Zone/Area)
func MapMasterRoutes(router chi.Router, mh fullschedule.MasterHandlers, mw *middleware.MiddlewareManager) {
	router.Group(func(pr chi.Router) {
		pr.Use(mw.Verifier(true), mw.Authenticator(), mw.CurrentUser(), mw.ActiveUser())

		pr.Route("/groups", func(r chi.Router) {
			r.Post("/", mh.CreateGroup())
			r.Get("/", mh.ListGroups())
			r.Route("/{id}", func(sr chi.Router) {
				sr.Get("/", mh.GetGroup())
				sr.Put("/", mh.UpdateGroup())
				sr.Delete("/", mh.DeleteGroup())
			})
		})

		pr.Route("/zones", func(r chi.Router) {
			r.Post("/", mh.CreateZone())
			r.Get("/", mh.ListZones())
			r.Route("/{id}", func(sr chi.Router) {
				sr.Get("/", mh.GetZone())
				sr.Put("/", mh.UpdateZone())
				sr.Delete("/", mh.DeleteZone())
			})
		})

		pr.Route("/areas", func(r chi.Router) {
			r.Post("/", mh.CreateArea())
			r.Get("/", mh.ListAreas())
			r.Route("/{id}", func(sr chi.Router) {
				sr.Get("/", mh.GetArea())
				sr.Put("/", mh.UpdateArea())
				sr.Delete("/", mh.DeleteArea())
				sr.Post("/devices", mh.MapAreaDevices())
				sr.Get("/devices", mh.AreaDevices())
			})
		})

		pr.Get("/fullschedule/scope/preview", mh.ScopePreview())
	})
}

// MapFullscheduleRoutes mounts the Full Schedule routes under /api/fullschedule.
// จดทะเบียนเส้นทางโมดูล Full Schedule
func MapFullscheduleRoutes(router chi.Router, h fullschedule.Handlers, mw *middleware.MiddlewareManager) {
	router.Route("/fullschedule", func(r chi.Router) {
		r.Group(func(pr chi.Router) {
			pr.Use(mw.Verifier(true), mw.Authenticator(), mw.CurrentUser(), mw.ActiveUser())

			pr.Post("/", h.Create())
			pr.Get("/", h.List())
			pr.Get("/devices", h.ListDevices())
			pr.Get("/history", h.GetHistory())
			pr.Get("/report", h.GetReport())
			pr.Get("/listschedulepage", h.ListSchedulePage())
			pr.Get("/scheduleall", h.ListScheduleAll())
			pr.Get("/scheduledevicepage", h.ListScheduleDevicePage())

			pr.Route("/{id}", func(sr chi.Router) {
				sr.Get("/", h.GetByID())
				sr.Put("/", h.Update())
				sr.Delete("/", h.Delete())
				sr.Put("/status", h.SetStatus())
				sr.Put("/day-status", h.UpdateDayStatus())
				sr.Post("/trigger", h.Trigger())
				sr.Post("/settings", h.SetSettings())
				sr.Get("/settings", h.GetSettings())
				sr.Get("/history", h.GetHistoryBySchedule())
			})
		})
	})
}
