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

// MonitoringMiddleware มิดเดิลแวร์ตรวจสอบ
func MonitoringMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		GlobalMetrics.TotalRequests.Add(1)
		GlobalMetrics.ActiveRequests.Add(1)
		defer GlobalMetrics.ActiveRequests.Add(-1)

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		if ww.Status() >= 400 {
			GlobalMetrics.TotalErrors.Add(1)
		}

		duration := time.Since(start)
		w.Header().Set("X-Response-Time", strconv.FormatInt(duration.Milliseconds(), 10)+"ms")
	})
}

// GetMetricsSnapshot returns current metrics
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
		UptimeSeconds:  uptime.Seconds(),
		QPS:            float64(total) / uptime.Seconds(),
	}
}

// MetricsSnapshot represents a point-in-time metrics snapshot
type MetricsSnapshot struct {
	TotalRequests  uint64  `json:"total_requests"`
	ActiveRequests int64   `json:"active_requests"`
	TotalErrors    uint64  `json:"total_errors"`
	ErrorRate      float64 `json:"error_rate"`
	UptimeSeconds  float64 `json:"uptime_seconds"`
	QPS            float64 `json:"qps"`
}

// MetricsHandler exposes the metrics endpoint
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(GlobalMetrics.StartTime).Seconds()
	total := GlobalMetrics.TotalRequests.Load()
	errors := GlobalMetrics.TotalErrors.Load()
	errorRate := 0.0
	if total > 0 {
		errorRate = float64(errors) / float64(total) * 100
	}
	resp := map[string]interface{}{
		"total_requests":  total,
		"active_requests": GlobalMetrics.ActiveRequests.Load(),
		"total_errors":    errors,
		"error_rate":      errorRate,
		"uptime_seconds":  uptime,
		"qps":             float64(total) / uptime,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
