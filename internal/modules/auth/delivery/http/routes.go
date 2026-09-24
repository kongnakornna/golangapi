package http

import (
	"icmongolang/internal/modules/auth"
	"icmongolang/internal/middleware"
	"time"

	"github.com/go-chi/chi/v5"
)

func (h *authHandler) Routes(r chi.Router) {
	// สร้าง rate limit middleware
	rateLimitConfig := middleware.RateLimitConfig{
		RequestsPerSecond: 50,               // 50 requests per second | 50 คำขอต่อวินาที
		Burst:             100,              // Burst up to 100 | อนุญาตให้ส่งสูงสุด 100 คำขอในครั้งเดียว
		CleanupInterval:   15 * time.Minute, // Clean up every 15 minutes | ทำความสะอาดทุก 15 นาที
	}
	rateLimiter := middleware.NewRateLimitMiddleware(rateLimitConfig)

	// ใช้ rate limit เฉพาะ route login
	r.Group(func(r chi.Router) {
		r.Use(rateLimiter.Handler)  // ใช้ rate limit middleware
		r.Post("/login", h.Login()) // route login
		r.Post("/signin", h.SignIn())
	})

	// // routes อื่นๆ ไม่มี rate limit
	// r.Post("/signup", h.SignUp())
	// r.Post("/logout", h.Logout())
	// r.Post("/refresh", h.Refresh())
}

func MapAuthRoute(router *chi.Mux, h auth.Handlers, mw *middleware.MiddlewareManager) {
	rateLimitConfig := middleware.RateLimitConfig{
		RequestsPerSecond: 10,
		Burst:             20,
		CleanupInterval:   15 * time.Minute,
	}
	rateLimiter := middleware.NewRateLimitMiddleware(rateLimitConfig)

	// Auth routes
	router.Route("/auth", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Use(rateLimiter.Handler)
			r.Post("/login", h.Login())
			r.Post("/signin", h.SignIn())
			r.Get("/publickey", h.GetPublicKey())
			r.Get("/verifyemail", h.VerifyEmail())
			r.Post("/forgotpassword", h.ForgotPassword())
			r.Patch("/resetpassword", h.ResetPassword())
		})
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(false))
			r.Use(mw.Authenticator())
			r.Get("/refresh", h.RefreshToken())
			r.Get("/logout", h.Logout())
			r.Get("/logoutall", h.LogoutAllToken())
			r.Post("/logoutall", h.LogoutAllToken())
		})
	})
}
