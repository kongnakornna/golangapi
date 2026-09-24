package middleware

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

// Declare metrics without auto‑registration
var (
	goroutines = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_goroutines",
		Help: "Number of goroutines",
	})
	memStats = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "go_memstats",
		Help: "Go memory statistics",
	}, []string{"type"})
	gcPauseTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "go_gc_pause_total_seconds",
		Help: "Total GC pause in seconds",
	})
	cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "system_cpu_usage_percent",
		Help: "System CPU usage percentage",
	})
	memUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "system_memory_usage_bytes",
		Help: "System memory usage in bytes",
	})
	netBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "system_network_bytes_total",
		Help: "Network bytes received/transmitted",
	}, []string{"direction"})
	redisMetrics = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "redis_stats",
		Help: "Redis statistics",
	}, []string{"stat"})
)

func init() {
	// Register metrics and ignore "already registered" errors
	mustRegister := func(c prometheus.Collector) {
		if err := prometheus.Register(c); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				log.Printf("⚠️ prometheus register failed: %v", err)
			}
		}
	}
	mustRegister(memStats)
	mustRegister(gcPauseTotal)
	mustRegister(cpuUsage)
	mustRegister(memUsage)
	mustRegister(netBytes)
	mustRegister(redisMetrics)
}

// StartMetricsCollector launches background collectors (unchanged)
func StartMetricsCollector(redisClient *redis.Client) {
	// … the rest of the function is the same as before …
	// (copy your existing StartMetricsCollector body here)
}
