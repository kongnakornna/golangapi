package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// Metrics ตัวชี้วัดประสิทธิภาพพื้นฐาน
type Metrics struct {
	TotalRequests  atomic.Uint64
	ActiveRequests atomic.Int64
	TotalErrors    atomic.Uint64
	StartTime      time.Time
}

// GlobalMetrics อินสแตนซ์ตัวชี้วัดส่วนกลาง
var GlobalMetrics = &Metrics{
	StartTime: time.Now(),
}

// MonitoringMiddleware มิดเดิลแวร์ตรวจสอบ (เวอร์ชันย่อ)
func MonitoringMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// เพิ่มจำนวนนับ
		GlobalMetrics.TotalRequests.Add(1)
		GlobalMetrics.ActiveRequests.Add(1)
		defer GlobalMetrics.ActiveRequests.Add(-1)

		// ห่อ ResponseWriter
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// ดำเนินการคำขอ
		next.ServeHTTP(ww, r)

		// บันทึกข้อผิดพลาด
		if ww.Status() >= 400 {
			GlobalMetrics.TotalErrors.Add(1)
		}

		// เพิ่ม Header เวลาตอบสนอง
		duration := time.Since(start)
		w.Header().Set("X-Response-Time", strconv.FormatInt(duration.Milliseconds(), 10)+"ms")
	})
}

// GetMetricsSnapshot รับสแนปชอตตัวชี้วัด
func GetMetricsSnapshot() MetricsSnapshot {
	uptime := time.Since(GlobalMetrics.StartTime)
	total := GlobalMetrics.TotalRequests.Load()
	errors := GlobalMetrics.TotalErrors.Load()

	var errorRate float64
	if total > 0 {
		errorRate = float64(errors) / float64(total) * 100
	}

	return MetricsSnapshot{
		TotalRequests:  total,
		ActiveRequests: GlobalMetrics.ActiveRequests.Load(),
		TotalErrors:    errors,
		ErrorRate:      errorRate,
		Uptime:         uptime,
		QPS:            float64(total) / uptime.Seconds(),
	}
}

// MetricsSnapshot สแนปชอตตัวชี้วัด
type MetricsSnapshot struct {
	TotalRequests  uint64        `json:"total_requests"`
	ActiveRequests int64         `json:"active_requests"`
	TotalErrors    uint64        `json:"total_errors"`
	ErrorRate      float64       `json:"error_rate"`
	Uptime         time.Duration `json:"uptime_seconds"`
	QPS            float64       `json:"qps"`
}

// MetricsHandler ตัวจัดการเอนด์พอยต์ตัวชี้วัด
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := GetMetricsSnapshot()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
