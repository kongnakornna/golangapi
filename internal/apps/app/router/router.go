package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	v1 "github.com/vadxq/go-rest-starter/internal/apps/app/router/v1"
	authhandler "github.com/vadxq/go-rest-starter/internal/core/auth/handler"
	healthmodule "github.com/vadxq/go-rest-starter/internal/core/health"
	healthhandler "github.com/vadxq/go-rest-starter/internal/core/health/handler"
	userhandler "github.com/vadxq/go-rest-starter/internal/core/user/handler"
	custommiddleware "github.com/vadxq/go-rest-starter/internal/transport/middleware"
	"github.com/vadxq/go-rest-starter/pkg/logger"
)

// RouterGroup คำจำกัดความประเภทกลุ่มเส้นทาง
type RouterGroup struct {
	Router     chi.Router
	Middleware []func(http.Handler) http.Handler
}

// RouterConfig การกำหนดค่าเส้นทาง
type RouterConfig struct {
	UserHandler   *userhandler.UserHandler
	AuthHandler   *authhandler.AuthHandler
	HealthHandler *healthhandler.HealthHandler
	JWTSecret     string
	Logger        logger.Logger
}

// Setup ตั้งค่าเส้นทาง API ทั้งหมด
func Setup(r chi.Router, config RouterConfig) {
	// ใช้มิดเดิลแวร์ทั่วโลก
	applyGlobalMiddleware(r, config.Logger)

	// เส้นทางเอกสาร API
	v1.SetupSwaggerRoutes(r)

	// การตรวจสอบสุขภาพและการตรวจสอบสถานะ
	healthmodule.RegisterRoutes(r, config.HealthHandler)

	// API v1
	setupV1Routes(r, config)
}

// applyGlobalMiddleware ใช้มิดเดิลแวร์ทั่วโลก
func applyGlobalMiddleware(r chi.Router, log logger.Logger) {
	// มิดเดิลแวร์พื้นฐาน
	r.Use(middleware.RequestID)                    // ID คำขอ
	r.Use(middleware.RealIP)                       // IP จริง
	r.Use(custommiddleware.RequestContext)         // บริบทคำขอ
	r.Use(custommiddleware.LoggingMiddleware(log)) // บันทึก
	r.Use(custommiddleware.RecoveryMiddleware)     // การกู้คืน
	r.Use(middleware.Timeout(60 * time.Second))    // หมดเวลา
	r.Use(middleware.CleanPath)                    // ล้างเส้นทาง
	r.Use(middleware.StripSlashes)                 // ตัดเครื่องหมายทับท้าย

	// มิดเดิลแวร์ความปลอดภัย
	r.Use(custommiddleware.CORSMiddleware) // ข้ามแหล่งที่มา
	r.Use(securityHeaders)                 // ส่วนหัวความปลอดภัย

	// มิดเดิลแวร์จำกัดอัตรา
	rateLimiter := custommiddleware.NewRateLimitMiddleware(custommiddleware.DefaultRateLimitConfig)
	r.Use(rateLimiter.Handler) // จำกัดอัตรา
}

// setupV1Routes ตั้งค่าเส้นทาง API v1
func setupV1Routes(r chi.Router, config RouterConfig) {
	// กำหนดเส้นทางที่ไม่ต้องตรวจสอบสิทธิ์
	excludePaths := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/refresh",
		"/swagger",
		"/health",
		"/version",
		"/status",
	}

	// สร้างการกำหนดค่า JWT
	jwtConfig := &custommiddleware.JWTConfig{
		Secret:       config.JWTSecret,
		ExcludePaths: excludePaths,
	}

	// เส้นทางพื้นฐาน API v1
	r.Route("/api/v1", func(r chi.Router) {
		v1Config := v1.RouterConfig{
			UserHandler: config.UserHandler,
			AuthHandler: config.AuthHandler,
			JWTSecret:   config.JWTSecret,
		}
		// กลุ่มเส้นทางสาธารณะ - ไม่ต้องการตรวจสอบสิทธิ์
		v1.SetupPublicRoutes(r, v1Config)
		// กลุ่มเส้นทางที่ได้รับการป้องกัน - ต้องการตรวจสอบสิทธิ์
		v1.SetupProtectedRoutes(r, v1Config, jwtConfig, config.Logger)
	})
}

// securityHeaders เพิ่มส่วนหัว HTTP ที่เกี่ยวข้องกับความปลอดภัย
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ป้องกันการตรวจจับชนิด MIME
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// เปิดใช้งานการป้องกัน XSS
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		// ป้องกันการคลิกจี้
		w.Header().Set("X-Frame-Options", "DENY")
		// การรักษาความปลอดภัยการส่งข้อมูลแบบเข้มงวด
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		// นโยบายอ้างอิง
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		// นโยบายความปลอดภัยของเนื้อหา
		cspValue := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data:; " +
			"style-src 'self' 'unsafe-inline'; " +
			"font-src 'self' data:;"
		w.Header().Set("Content-Security-Policy", cspValue)

		next.ServeHTTP(w, r)
	})
}