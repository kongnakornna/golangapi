package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// SetupSwaggerRoutes ตั้งค่าเส้นทางเอกสาร Swagger
func SetupSwaggerRoutes(r chi.Router) {
	// ตั้งค่าเส้นทาง UI ของ Swagger
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // URL ของ JSON API ของ Swagger
	))

	// เปลี่ยนเส้นทางเส้นทางหลักไปยัง UI ของ Swagger
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})
}
