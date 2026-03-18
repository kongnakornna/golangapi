package handler

import (
	"context"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	httpx "github.com/kongnakornna/golangapi/internal/transport/httpx"
	"github.com/kongnakornna/golangapi/pkg/logger"
)

// HealthHandler ตัวจัดการตรวจสอบสุขภาพระบบ
type HealthHandler struct {
	db     *gorm.DB
	redis  *redis.Client
	logger logger.Logger
}

// NewHealthHandler สร้างตัวจัดการตรวจสอบสุขภาพระบบ
func NewHealthHandler(db *gorm.DB, redis *redis.Client, log logger.Logger) *HealthHandler {
	return &HealthHandler{
		db:     db,
		redis:  redis,
		logger: log,
	}
}

// HealthStatus โครงสร้างสถานะสุขภาพ
type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version"`
	Uptime    string            `json:"uptime,omitempty"`
}

var startTime = time.Now()

// Health ตรวจสอบสุขภาพพื้นฐาน
// @Summary ตรวจสอบสุขภาพ
// @Description ตรวจสอบสถานะพื้นฐานของแอปพลิเคชัน
// @Tags health
// @Produce json
// @Success 200 {object} HealthStatus
// @Router /health [get]
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	status := &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    time.Since(startTime).String(),
		Services:  make(map[string]string),
	}

	httpx.JSON(w, r, http.StatusOK, status)
}

// DetailedHealth ตรวจสอบสุขภาพโดยละเอียด
// @Summary ตรวจสอบสุขภาพโดยละเอียด
// @Description ตรวจสอบสถานะของแอปพลิเคชันและบริการที่ต้องพึ่งพา
// @Tags health
// @Produce json
// @Success 200 {object} HealthStatus
// @Success 503 {object} HealthStatus "บริการไม่พร้อมใช้งาน"
// @Router /health/detailed [get]
func (h *HealthHandler) DetailedHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	status := &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    time.Since(startTime).String(),
		Services:  make(map[string]string),
	}

	// ตรวจสอบการเชื่อมต่อฐานข้อมูล
	dbStatus := h.checkDatabase(ctx)
	status.Services["database"] = dbStatus

	// ตรวจสอบการเชื่อมต่อ Redis
	redisStatus := h.checkRedis(ctx)
	status.Services["redis"] = redisStatus

	// กำหนดสถานะโดยรวม
	if dbStatus != "healthy" || redisStatus != "healthy" {
		status.Status = "unhealthy"
		httpx.JSON(w, r, http.StatusServiceUnavailable, status)
		return
	}

	httpx.JSON(w, r, http.StatusOK, status)
}

// Ready ตรวจสอบความพร้อม
// @Summary ตรวจสอบความพร้อม
// @Description ตรวจสอบว่าแอปพลิเคชันพร้อมรับคำขอหรือไม่
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Success 503 {object} map[string]interface{} "บริการยังไม่พร้อม"
// @Router /ready [get]
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	ready := true
	checks := make(map[string]interface{})

	// ตรวจสอบฐานข้อมูล
	if h.checkDatabase(ctx) != "healthy" {
		ready = false
		checks["database"] = "not ready"
	} else {
		checks["database"] = "ready"
	}

	// ตรวจสอบ Redis
	if h.checkRedis(ctx) != "healthy" {
		ready = false
		checks["redis"] = "not ready"
	} else {
		checks["redis"] = "ready"
	}

	response := map[string]interface{}{
		"ready":     ready,
		"timestamp": time.Now(),
		"checks":    checks,
	}

	if ready {
		httpx.JSON(w, r, http.StatusOK, response)
	} else {
		httpx.JSON(w, r, http.StatusServiceUnavailable, response)
	}
}

// Live ตรวจสอบการทำงาน
// @Summary ตรวจสอบการทำงาน
// @Description ตรวจสอบว่าแอปพลิเคชันยังทำงานอยู่ (ความสามารถในการตอบสนองพื้นฐาน)
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /live [get]
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"alive":     true,
		"timestamp": time.Now(),
	}
	httpx.JSON(w, r, http.StatusOK, response)
}

// checkDatabase ตรวจสอบสถานะการเชื่อมต่อฐานข้อมูล
func (h *HealthHandler) checkDatabase(ctx context.Context) string {
	if h.db == nil {
		return "unavailable"
	}

	sqlDB, err := h.db.DB()
	if err != nil {
		h.logger.Error("获取数据库连接失败", "error", err)
		return "error"
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		h.logger.Error("数据库ping失败", "error", err)
		return "unhealthy"
	}

	return "healthy"
}

// checkRedis ตรวจสอบสถานะการเชื่อมต่อ Redis
func (h *HealthHandler) checkRedis(ctx context.Context) string {
	if h.redis == nil {
		return "unavailable"
	}

	if err := h.redis.Ping(ctx).Err(); err != nil {
		h.logger.Error("Redis ping失败", "error", err)
		return "unhealthy"
	}

	return "healthy"
}

// Readiness พร็อบความพร้อมสำหรับ K8s
func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	h.Ready(w, r)
}

// Liveness พร็อบการทำงานสำหรับ K8s
func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	h.Live(w, r)
}

// SystemInfo ข้อมูลระบบ
func (h *HealthHandler) SystemInfo(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	systemInfo := map[string]interface{}{
		"runtime": map[string]interface{}{
			"version":    runtime.Version(),
			"goroutines": runtime.NumGoroutine(),
			"cpu_count":  runtime.NumCPU(),
			"goos":       runtime.GOOS,
			"goarch":     runtime.GOARCH,
		},
		"memory": map[string]interface{}{
			"alloc_mb":       float64(m.Alloc) / 1024 / 1024,
			"total_alloc_mb": float64(m.TotalAlloc) / 1024 / 1024,
			"sys_mb":         float64(m.Sys) / 1024 / 1024,
			"num_gc":         m.NumGC,
			"heap_alloc_mb":  float64(m.HeapAlloc) / 1024 / 1024,
			"heap_sys_mb":    float64(m.HeapSys) / 1024 / 1024,
		},
		"application": map[string]interface{}{
			"version": "1.0.0",
			"uptime":  time.Since(startTime).String(),
			"started": startTime.Format(time.RFC3339),
		},
		"timestamp": time.Now().Unix(),
	}

	httpx.JSON(w, r, http.StatusOK, systemInfo)
}

// CheckDependencies ตรวจสอบบริการที่ต้องพึ่งพาทั้งหมด
func (h *HealthHandler) CheckDependencies(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	type DependencyStatus struct {
		Name         string        `json:"name"`
		Status       string        `json:"status"`
		ResponseTime time.Duration `json:"response_time_ms"`
		Error        string        `json:"error,omitempty"`
	}

	var dependencies []DependencyStatus
	var wg sync.WaitGroup
	var mu sync.Mutex

	// ตรวจสอบฐานข้อมูล
	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		status := "healthy"
		var errMsg string

		if h.db != nil {
			sqlDB, err := h.db.DB()
			if err != nil {
				status = "error"
				errMsg = err.Error()
			} else if err := sqlDB.PingContext(ctx); err != nil {
				status = "unhealthy"
				errMsg = err.Error()
			}
		} else {
			status = "unavailable"
		}

		mu.Lock()
		dependencies = append(dependencies, DependencyStatus{
			Name:         "postgresql",
			Status:       status,
			ResponseTime: time.Since(start) / time.Millisecond,
			Error:        errMsg,
		})
		mu.Unlock()
	}()

	// ตรวจสอบ Redis
	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		status := "healthy"
		var errMsg string

		if h.redis != nil {
			if err := h.redis.Ping(ctx).Err(); err != nil {
				status = "unhealthy"
				errMsg = err.Error()
			}
		} else {
			status = "unavailable"
		}

		mu.Lock()
		dependencies = append(dependencies, DependencyStatus{
			Name:         "redis",
			Status:       status,
			ResponseTime: time.Since(start) / time.Millisecond,
			Error:        errMsg,
		})
		mu.Unlock()
	}()

	wg.Wait()

	// กำหนดสถานะโดยรวม
	overallStatus := "healthy"
	for _, dep := range dependencies {
		if dep.Status != "healthy" {
			overallStatus = "degraded"
			if dep.Status == "unhealthy" || dep.Status == "error" {
				overallStatus = "unhealthy"
				break
			}
		}
	}

	response := map[string]interface{}{
		"status":       overallStatus,
		"dependencies": dependencies,
		"timestamp":    time.Now().Unix(),
	}

	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	httpx.JSON(w, r, statusCode, response)
}
