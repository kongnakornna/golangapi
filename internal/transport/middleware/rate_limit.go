package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimitConfig โครงสร้างการกำหนดค่าการจำกัดอัตรา
type RateLimitConfig struct {
	RequestsPerSecond int           // จำนวนคำขอที่อนุญาตต่อวินาที
	Burst             int           // จำนวนคำขอที่อนุญาตให้ส่งแบบกระชุ
	CleanupInterval   time.Duration // ช่วงเวลาทำความสะอาดบันทึกที่หมดอายุ
}

// DefaultRateLimitConfig ค่ากำหนดเริ่มต้นสำหรับการจำกัดอัตรา
var DefaultRateLimitConfig = RateLimitConfig{
	RequestsPerSecond: 10,
	Burst:             20,
	CleanupInterval:   10 * time.Minute,
}

// rateLimiter ตัวจำกัดอัตราสำหรับแต่ละ IP
type rateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimitMiddleware มิดเดิลแวร์จำกัดอัตราตาม IP
type RateLimitMiddleware struct {
	config   RateLimitConfig
	limiters map[string]*rateLimiter
	mu       sync.RWMutex
}

// NewRateLimitMiddleware สร้างมิดเดิลแวร์จำกัดอัตราใหม่
func NewRateLimitMiddleware(config RateLimitConfig) *RateLimitMiddleware {
	rlm := &RateLimitMiddleware{
		config:   config,
		limiters: make(map[string]*rateLimiter),
	}

	// เริ่ม goroutine สำหรับทำความสะอาด
	go rlm.cleanup()

	return rlm
}

// Handler ฟังก์ชันจัดการของมิดเดิลแวร์จำกัดอัตรา
func (rlm *RateLimitMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ดึง IP ของไคลเอนต์
		ip := getClientIP(r)

		// ดึงหรือสร้างตัวจำกัดสำหรับ IP นี้
		limiter := rlm.getLimiter(ip)

		// ตรวจสอบว่าอนุญาตคำขอนี้หรือไม่
		if !limiter.Allow() {
			writeRateLimitResponse(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getLimiter ดึงหรือสร้างตัวจำกัดสำหรับ IP ที่กำหนด
func (rlm *RateLimitMiddleware) getLimiter(ip string) *rate.Limiter {
	rlm.mu.Lock()
	defer rlm.mu.Unlock()

	limiterInfo, exists := rlm.limiters[ip]
	if !exists {
		limiterInfo = &rateLimiter{
			limiter: rate.NewLimiter(
				rate.Limit(rlm.config.RequestsPerSecond),
				rlm.config.Burst,
			),
			lastSeen: time.Now(),
		}
		rlm.limiters[ip] = limiterInfo
	} else {
		limiterInfo.lastSeen = time.Now()
	}

	return limiterInfo.limiter
}

// cleanup ทำความสะอาดตัวจำกัดที่หมดอายุเป็นระยะ
func (rlm *RateLimitMiddleware) cleanup() {
	ticker := time.NewTicker(rlm.config.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rlm.mu.Lock()
		cutoff := time.Now().Add(-rlm.config.CleanupInterval * 2)

		for ip, limiterInfo := range rlm.limiters {
			if limiterInfo.lastSeen.Before(cutoff) {
				delete(rlm.limiters, ip)
			}
		}
		rlm.mu.Unlock()
	}
}

// writeRateLimitResponse เขียนการตอบสนองเมื่อถูกจำกัดอัตรา
func writeRateLimitResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-RateLimit-Limit", "10")
	w.Header().Set("X-RateLimit-Remaining", "0")
	w.Header().Set("Retry-After", "60")
	w.WriteHeader(http.StatusTooManyRequests)

	response := `{
		"error": {
			"type": "RATE_LIMIT_EXCEEDED",
			"message": "ส่งคำขอถี่เกินไป กรุณาลองใหม่ภายหลัง",
			"details": "Rate limit exceeded. Please try again later."
		}
	}`

	w.Write([]byte(response))
}

// getClientIP ดึง IP จริงของไคลเอนต์
func getClientIP(r *http.Request) string {
	// ตรวจสอบหัว X-Forwarded-For
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// ใช้เฉพาะ IP แรก
		if idx := len(xff); idx > 0 {
			return xff[:idx]
		}
	}

	// ตรวจสอบหัว X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// ตรวจสอบ X-Forwarded-For อีกครั้ง (กรณีไม่มีการตัด)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// ใช้ RemoteAddr เป็นค่าเริ่มต้น
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
