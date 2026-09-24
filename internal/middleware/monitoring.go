package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})
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
		UptimeSeconds:  uptime.Seconds(),
		QPS:            float64(total) / uptime.Seconds(),
	}
}

// MetricsSnapshot สแนปชอตตัวชี้วัด
type MetricsSnapshot struct {
	TotalRequests  uint64  `json:"total_requests"`
	ActiveRequests int64   `json:"active_requests"`
	TotalErrors    uint64  `json:"total_errors"`
	ErrorRate      float64 `json:"error_rate"`
	UptimeSeconds  float64 `json:"uptime_seconds"`
	QPS            float64 `json:"qps"`
}

// MetricsHandler ตัวจัดการเอนด์พอยต์ตัวชี้วัด (custom JSON)
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := GetMetricsSnapshot()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// MonitoringMiddleware – Prometheus + custom counters
func MonitoringMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		duration := time.Since(start).Seconds()
		path := r.URL.Path
		method := r.Method
		status := ww.Status()

		httpRequestsTotal.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)

		// Keep original counters for custom /metrics endpoint
		GlobalMetrics.TotalRequests.Add(1)
		GlobalMetrics.ActiveRequests.Add(1)
		defer GlobalMetrics.ActiveRequests.Add(-1)
		if status >= 400 {
			GlobalMetrics.TotalErrors.Add(1)
		}
		w.Header().Set("X-Response-Time", strconv.FormatInt(int64(duration*1000), 10)+"ms")
	})
}
