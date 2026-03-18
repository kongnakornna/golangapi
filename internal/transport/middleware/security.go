package middleware

import (
	"net/http"
	"strings"
)

// SecurityConfig การกำหนดค่าความปลอดภัย
type SecurityConfig struct {
	// การกำหนดค่า CSP
	ContentSecurityPolicy string
	// การกำหนดค่า HSTS
	StrictTransportSecurity string
	// ต้นทางที่อนุญาต
	AllowedOrigins []string
	// เปิดใช้งานคุณสมบัติความปลอดภัยต่างๆ หรือไม่
	EnableCSP       bool
	EnableHSTS      bool
	EnableXSS       bool
	EnableNoSniff   bool
	EnableFrameDeny bool
	EnableReferrer  bool
}

// DefaultSecurityConfig การกำหนดค่าความปลอดภัยเริ่มต้น
var DefaultSecurityConfig = SecurityConfig{
	ContentSecurityPolicy: "default-src 'self'; " +
		"script-src 'self' 'unsafe-inline'; " +
		"style-src 'self' 'unsafe-inline'; " +
		"img-src 'self' data: https:; " +
		"font-src 'self' data:; " +
		"connect-src 'self'; " +
		"media-src 'self'; " +
		"object-src 'none'; " +
		"frame-ancestors 'none'; " +
		"base-uri 'self'; " +
		"form-action 'self'; " +
		"upgrade-insecure-requests;",
	StrictTransportSecurity: "max-age=31536000; includeSubDomains; preload",
	AllowedOrigins:          []string{"https://icmon.com"},
	EnableCSP:               true,
	EnableHSTS:              true,
	EnableXSS:               true,
	EnableNoSniff:           true,
	EnableFrameDeny:         true,
	EnableReferrer:          true,
}

// SecurityMiddleware มิดเดิลแวร์ความปลอดภัย
func SecurityMiddleware(config *SecurityConfig) func(http.Handler) http.Handler {
	if config == nil {
		config = &DefaultSecurityConfig
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ป้องกันการตรวจจับประเภท MIME
			if config.EnableNoSniff {
				w.Header().Set("X-Content-Type-Options", "nosniff")
			}

			// การป้องกัน XSS
			if config.EnableXSS {
				w.Header().Set("X-XSS-Protection", "1; mode=block")
			}

			// ป้องกัน Clickjacking
			if config.EnableFrameDeny {
				w.Header().Set("X-Frame-Options", "DENY")
			}

			// HTTP Strict Transport Security (ใช้ได้เฉพาะบน HTTPS)
			if config.EnableHSTS && (r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https") {
				w.Header().Set("Strict-Transport-Security", config.StrictTransportSecurity)
			}

			// นโยบาย Referrer
			if config.EnableReferrer {
				w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			}

			// นโยบายความปลอดภัยของเนื้อหา
			if config.EnableCSP {
				w.Header().Set("Content-Security-Policy", config.ContentSecurityPolicy)
			}

			// นโยบายสิทธิ์ (Feature Policy / Permissions Policy)
			w.Header().Set("Permissions-Policy",
				"accelerometer=(), "+
					"camera=(), "+
					"geolocation=(), "+
					"gyroscope=(), "+
					"magnetometer=(), "+
					"microphone=(), "+
					"payment=(), "+
					"usb=()")

			// ป้องกันเบราว์เซอร์แคชข้อมูลที่ละเอียดอ่อน
			if strings.Contains(r.URL.Path, "/api/") {
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
			}

			next.ServeHTTP(w, r)
		})
	}
}

// BasicSecurityHeaders มิดเดิลแวร์ส่วนหัวความปลอดภัยพื้นฐาน (แบบง่าย)
func BasicSecurityHeaders(next http.Handler) http.Handler {
	return SecurityMiddleware(&DefaultSecurityConfig)(next)
}

// NoCacheMiddleware มิดเดิลแวร์ปิดใช้งานแคช
func NoCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}

// SecureRedirectMiddleware มิดเดิลแวร์เปลี่ยนเส้นทาง HTTPS
func SecureRedirectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ตรวจสอบว่าเป็น HTTPS หรือยัง
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			next.ServeHTTP(w, r)
			return
		}

		// สร้าง URL HTTPS
		target := "https://" + r.Host + r.URL.Path
		if r.URL.RawQuery != "" {
			target += "?" + r.URL.RawQuery
		}

		// ดำเนินการเปลี่ยนเส้นทางถาวร 301
		http.Redirect(w, r, target, http.StatusMovedPermanently)
	})
}
