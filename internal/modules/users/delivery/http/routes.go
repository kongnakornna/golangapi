package http

import (
	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/users"

	"github.com/go-chi/chi/v5"
)

// MapUserRoute — register ทั้งหมดของโมดูล users
func MapUserRoute(router *chi.Mux, h users.Handlers, mw *middleware.MiddlewareManager) {
	// ==========================================================
	// PUBLIC
	// ==========================================================
	router.Post("/register", h.Register())
	router.Post("/signin", h.SignInEmail())
	router.Post("/login", h.SignInUsername())

	// ==========================================================
	// PROTECTED — /user
	// ==========================================================
	router.Route("/user", func(r chi.Router) {
		r.Use(mw.Verifier(true))
		r.Use(mw.Authenticator())
		r.Use(mw.CurrentUser())
		r.Use(mw.ActiveUser())

		// ---- Self ----
		r.Get("/me", h.Me())
		r.Put("/me", h.UpdateMe())
		r.Patch("/me/updatepass", h.UpdatePasswordMe())
		r.Post("/me/avatar", h.UploadAvatar())
		r.Get("/profile/{id}", h.Profile())

		// ==========================================================
		// Admin only
		// ==========================================================
		r.Group(func(admin chi.Router) {
			admin.Use(mw.SuperUser())

			// Users
			admin.Get("/", h.GetMulti())
			admin.Post("/", h.Create())
			admin.Get("/list", h.ListUsers())
			admin.Get("/statistics", h.Statistics())
			admin.Get("/notify/{channel}", h.NotifyList())
			admin.Patch("/{id}/role", h.UpdateRole())
			admin.Patch("/{id}/activestatus", h.UpdateActiveStatus())

			// Access Menus
			admin.Route("/access-menus", func(r chi.Router) {
				r.Get("/", h.ListAccessMenus())
				r.Post("/", h.CreateAccessMenu())
				r.Get("/{id}", h.GetAccessMenu())
				r.Put("/{id}", h.UpdateAccessMenu())
				r.Delete("/{id}", h.DeleteAccessMenu())
			})

			// Files
			admin.Route("/files", func(r chi.Router) {
				r.Get("/", h.ListFiles())
				r.Post("/", h.CreateFile())
				r.Get("/{id}", h.GetFile())
				r.Put("/{id}", h.UpdateFile())
				r.Delete("/{id}", h.DeleteFile())
			})

			// Logs
			admin.Route("/logs", func(r chi.Router) {
				r.Get("/", h.ListLogs())
				r.Post("/", h.CreateLog())
				r.Get("/{log_type_id}/{id}", h.GetLog())
				r.Delete("/{log_type_id}/{id}", h.DeleteLog())
			})

			// Log Types
			admin.Route("/log-types", func(r chi.Router) {
				r.Get("/", h.ListLogTypes())
				r.Post("/", h.CreateLogType())
				r.Get("/{id}", h.GetLogType())
				r.Put("/{id}", h.UpdateLogType())
				r.Delete("/{id}", h.DeleteLogType())
			})

			// Modules — ⚠️ /sort must come before /{id}
			admin.Route("/modules", func(r chi.Router) {
				r.Get("/", h.ListModules())
				r.Post("/", h.CreateModule())
				r.Patch("/sort", h.BulkUpdateModuleSort())
				r.Get("/{id}", h.GetModule())
				r.Put("/{id}", h.UpdateModule())
				r.Delete("/{id}", h.DeleteModule())
			})

			// Menus — ⚠️ /sort must come before /{id}
			admin.Route("/menus", func(r chi.Router) {
				r.Get("/", h.ListMenus())
				r.Post("/", h.CreateMenu())
				r.Patch("/sort", h.BulkUpdateMenuSort())
				r.Get("/{id}", h.GetMenu())
				r.Put("/{id}", h.UpdateMenu())
				r.Delete("/{id}", h.DeleteMenu())
			})
		})

		// ---- Per-id admin write ----
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.Get())
			r.Group(func(write chi.Router) {
				write.Use(mw.SuperUser())
				write.Delete("/", h.Delete())
				write.Put("/", h.Update())
				write.Patch("/updatepass", h.UpdatePassword())
				write.Get("/logoutall", h.LogoutAllAdmin())
			})
		})
	})
}