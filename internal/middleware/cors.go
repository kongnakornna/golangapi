package middleware

import (
	"github.com/go-chi/cors"
)

func (mw *MiddlewareManager) Cors() cors.Options {
	origins := []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:8080"}
	if mw.cfg.Server.Mode == "production" {
		origins = []string{"https://" + mw.cfg.Server.BaseUrl}
	}
	return cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}
}
