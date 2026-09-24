# Module 29: `pkg/exporter_prometheus` (ระบบส่ง Metrics ให้ Prometheus)

## สำหรับโฟลเดอร์ `pkg/exporter_prometheus/`

ไฟล์ที่เกี่ยวข้อง:

- `config.go`       – โครงสร้างการตั้งค่า Prometheus (enable, namespace, subsystem, metrics path)
- `metrics.go`      – ประกาศ metrics ที่ใช้ในระบบ (Counter, Gauge, Histogram, Summary)
- `middleware.go`   – Middleware สำหรับเก็บ HTTP request metrics (duration, status, method, path)
- `prometheus.go`   – ฟังก์ชันเริ่มต้นและลงทะเบียน metrics กับ Prometheus registry
- `registry.go`     – Singleton registry และ helper functions

---

## หลักการ (Concept)

### คืออะไร?

โมดูลสำหรับรวบรวม metrics ของแอปพลิเคชัน (จำนวน request, response time, error rate, ข้อมูล business เช่น จำนวนอุปกรณ์ IoT ออนไลน์) และ expose ผ่าน HTTP endpoint `/metrics` ให้ Prometheus ดึงไปเก็บ เพื่อใช้ในการเฝ้าระวัง (monitoring) และแจ้งเตือน (alerting)

### มีกี่แบบ?

- **Counter** – ค่าเพิ่มขึ้นเรื่อย ๆ (จำนวน request, จำนวน error)
- **Gauge** – ค่าที่ขึ้นลงได้ (จำนวน connection, อุณหภูมิ, memory usage)
- **Histogram** – กระจายค่าตาม bucket (response time, request size)
- **Summary** – ค่า percentile (quantiles) เหมาะกับ latency

**ข้อห้ามสำคัญ:**  
- ห้ามใช้ label ที่มีค่าไม่จำกัด (เช่น `user_id`) เพราะจะสร้าง series มากเกินไป  
- ห้ามอัปเดต metrics ในทุก request หากไม่จำเป็น (อาจทำให้ช้า)  
- ห้าม expose metrics ที่มีข้อมูล sensitive (token, password)

### ใช้อย่างไร / นำไปใช้กรณีไหน

- ตรวจสอบสุขภาพของ API (RPS, error rate, latency)  
- แจ้งเตือนเมื่อ error rate สูง หรือ response time นาน  
- ดูแนวโน้มการใช้งาน (จำนวน user active, device online)  
- ร่วมกับ Grafana เพื่อทำ dashboard

### ประโยชน์ที่ได้รับ

- เข้าใจพฤติกรรมของระบบแบบ real-time  
- ตรวจจับปัญหาก่อน user สังเกต  
- ช่วยในการวางแผน capacity  
- มาตรฐานอุตสาหกรรม (Prometheus + Grafana)

### ข้อควรระวัง

- Prometheus จะ scrape ทุกช่วงเวลา (default 15-30 วินาที) ต้องแน่ใจว่า endpoint `/metrics` เร็ว  
- Label ที่ใช้บ่อยควรจำกัด cardinality  
- การใช้ Histogram กับ bucket ที่เหมาะสม (เช่น response time bucket 0.01, 0.05, 0.1, 0.5, 1, 2, 5)

### ข้อดี

- เปิดรับ standard มาก รองรับ tooling หลากหลาย  
- ใช้งานง่าย ใช้ Go client library โดยตรง  
- ไม่ต้องใช้ external service (แค่ HTTP endpoint)

### ข้อเสีย

- ข้อมูลไม่ถูกเก็บถาวรใน Prometheus เอง (ต้องมี remote storage)  
- การ query ข้อมูลซับซ้อน (ใช้ PromQL)  
- อาจต้องปรับแต่ง performance ถ้า metrics เยอะมาก

### ข้อห้าม

- ห้ามใช้ Prometheus แทน logging system  
- ห้าม expose metrics endpoint โดยไม่มีการป้องกัน (อาจถูก ddos)  
- ห้ามเก็บ metrics ความถี่สูงเกิน 1 ครั้งต่อวินาที (อาจ overload)

---

## การออกแบบ Workflow และ Dataflow

```mermaid
flowchart LR
    A[HTTP Request] --> B[Prometheus Middleware]
    B --> C[บันทึก metrics: request count, duration, status]
    C --> D[Business Logic]
    D --> E[Business Metrics: device count, alarm count]
    E --> F[Prometheus Registry]
    F --> G[/metrics endpoint]
    G --> H[Prometheus Server]
    H --> I[Grafana Dashboard]
    H --> J[Alertmanager]
```

**คำอธิบาย:**  
1. Middleware จับทุก request บันทึก method, path, status, duration  
2. Business logic อัปเดต metrics ที่เกี่ยวข้อง (เช่น จำนวน device online)  
3. Metrics ทั้งหมดถูกรวบรวมใน registry  
4. Prometheus server ดึงข้อมูลผ่าน `/metrics` ทุก interval  
5. Alertmanager ส่งการแจ้งเตือนตามกฎที่ตั้งไว้

---

## ตัวอย่างโค้ดที่รันได้จริง

### 1. การติดตั้ง Dependencies

```bash
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp
```

### 2. เพิ่ม configuration ใน `config/config.go`

```go
type PrometheusConfig struct {
    Enabled     bool   `mapstructure:"enabled"`
    MetricsPath string `mapstructure:"metricsPath"`
    Namespace   string `mapstructure:"namespace"`
    Subsystem   string `mapstructure:"subsystem"`
}

type Config struct {
    // ... existing fields
    Prometheus PrometheusConfig `mapstructure:"prometheus"`
}
```

และเพิ่มใน `config-local.yml`:

```yaml
prometheus:
  enabled: true
  metricsPath: "/metrics"
  namespace: "icmongolang"
  subsystem: "api"
```

### 3. โค้ด `pkg/exporter_prometheus/config.go`

```go
package exporter_prometheus

import "icmongolang/config"

type PromConfig struct {
    Enabled     bool
    MetricsPath string
    Namespace   string
    Subsystem   string
}

func FromAppConfig(cfg *config.Config) PromConfig {
    return PromConfig{
        Enabled:     cfg.Prometheus.Enabled,
        MetricsPath: cfg.Prometheus.MetricsPath,
        Namespace:   cfg.Prometheus.Namespace,
        Subsystem:   cfg.Prometheus.Subsystem,
    }
}
```

### 4. โค้ด `pkg/exporter_prometheus/metrics.go`

```go
package exporter_prometheus

import (
    "github.com/prometheus/client_golang/prometheus"
)

// Metrics struct holds all custom metrics
type Metrics struct {
    // HTTP metrics
    HttpRequestsTotal   *prometheus.CounterVec
    HttpRequestDuration *prometheus.HistogramVec
    RequestsInFlight    prometheus.Gauge

    // Business metrics (example)
    DevicesOnline   prometheus.Gauge
    AlertsTotal     *prometheus.CounterVec

    // System metrics (optional)
    GoGoroutines prometheus.Gauge
}

// NewMetrics creates and registers all metrics
func NewMetrics(namespace, subsystem string) *Metrics {
    m := &Metrics{
        HttpRequestsTotal: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Namespace: namespace,
                Subsystem: subsystem,
                Name:      "requests_total",
                Help:      "Total number of HTTP requests",
            },
            []string{"method", "path", "status"},
        ),
        HttpRequestDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Namespace: namespace,
                Subsystem: subsystem,
                Name:      "request_duration_seconds",
                Help:      "HTTP request duration in seconds",
                Buckets:   prometheus.DefBuckets,
            },
            []string{"method", "path"},
        ),
        RequestsInFlight: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Namespace: namespace,
                Subsystem: subsystem,
                Name:      "requests_in_flight",
                Help:      "Current number of HTTP requests being processed",
            },
        ),
        DevicesOnline: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Namespace: namespace,
                Subsystem: "iot",
                Name:      "devices_online",
                Help:      "Number of IoT devices currently online",
            },
        ),
        AlertsTotal: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Namespace: namespace,
                Subsystem: "iot",
                Name:      "alerts_total",
                Help:      "Total number of alerts by severity",
            },
            []string{"severity"},
        ),
        GoGoroutines: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Namespace: namespace,
                Subsystem: "go",
                Name:      "goroutines",
                Help:      "Number of goroutines",
            },
        ),
    }

    // Register all metrics
    prometheus.MustRegister(m.HttpRequestsTotal)
    prometheus.MustRegister(m.HttpRequestDuration)
    prometheus.MustRegister(m.RequestsInFlight)
    prometheus.MustRegister(m.DevicesOnline)
    prometheus.MustRegister(m.AlertsTotal)
    prometheus.MustRegister(m.GoGoroutines)

    return m
}
```

### 5. โค้ด `pkg/exporter_prometheus/middleware.go`

```go
package exporter_prometheus

import (
    "net/http"
    "strconv"
    "time"

    "github.com/go-chi/chi/v5/middleware"
)

// HTTPMiddleware returns a chi middleware that records Prometheus metrics
func HTTPMiddleware(metrics *Metrics) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Increase in-flight counter
            metrics.RequestsInFlight.Inc()
            defer metrics.RequestsInFlight.Dec()

            // Wrap response writer to capture status code
            ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

            start := time.Now()
            defer func() {
                duration := time.Since(start).Seconds()
                // Record duration histogram
                metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
                // Record request total counter with status
                status := strconv.Itoa(ww.Status())
                metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
            }()

            next.ServeHTTP(ww, r)
        })
    }
}
```

### 6. โค้ด `pkg/exporter_prometheus/prometheus.go` (Registry และ helper)

```go
package exporter_prometheus

import (
    "net/http"

    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// RegisterMetricsHandler returns an HTTP handler for the /metrics endpoint
func RegisterMetricsHandler() http.Handler {
    return promhttp.Handler()
}

// UpdateDevicesOnline updates the gauge for online devices (example business metric)
func (m *Metrics) SetDevicesOnline(count float64) {
    m.DevicesOnline.Set(count)
}

// IncrementAlerts increments alert counter by severity
func (m *Metrics) IncrementAlerts(severity string) {
    m.AlertsTotal.WithLabelValues(severity).Inc()
}
```

### 7. การรวมเข้ากับ `internal/server/handlers.go` (หรือ `server.go`)

```go
// ในฟังก์ชัน New(...) ของ internal/server
import (
    // ... existing imports
    "icmongolang/pkg/exporter_prometheus"
)

func New(db *gorm.DB, redisClient *redis.Client, taskRedisClient *asynq.Client, cfg *config.Config, logger logger.Logger) (*chi.Mux, error) {
    r := chi.NewRouter()
    // ... existing middleware (logger, recoverer, cors, etc.)

    // ==== Prometheus Setup ====
    promCfg := exporter_prometheus.FromAppConfig(cfg)
    var promMetrics *exporter_prometheus.Metrics
    if promCfg.Enabled {
        promMetrics = exporter_prometheus.NewMetrics(promCfg.Namespace, promCfg.Subsystem)
        // Add Prometheus middleware to record HTTP metrics
        r.Use(exporter_prometheus.HTTPMiddleware(promMetrics))
        // Mount /metrics endpoint
        r.Get(promCfg.MetricsPath, func(w http.ResponseWriter, r *http.Request) {
            exporter_prometheus.RegisterMetricsHandler().ServeHTTP(w, r)
        })
        logger.Infof("Prometheus metrics enabled at %s", promCfg.MetricsPath)
    }

    // ... rest of routes (API, swagger, etc.)

    // Optional: expose promMetrics to other packages (via context or global)
    // e.g., store in a global variable or pass to usecases
    return r, nil
}
```

### 8. ตัวอย่างการอัปเดต business metrics จาก usecase หรือ worker

```go
// ใน internal/modules/iot/usecase หรือ worker
func (uc *IotUsecase) UpdateDeviceStatus(ctx context.Context, deviceID string, online bool) error {
    // ... logic
    if online {
        // สมมติมี global หรือ access ถึง promMetrics
        promMetrics.SetDevicesOnline(currentOnlineCount)
    }
    return nil
}

// เมื่อมี alarm เกิดขึ้น
func (uc *IotUsecase) ProcessAlarm(deviceID string, severity string) {
    // ... logic
    promMetrics.IncrementAlerts(severity)
}
```

---

## วิธีใช้งาน module นี้

1. เปิดใช้งาน Prometheus ใน `config-local.yml` (`enabled: true`)  
2. รันแอปพลิเคชัน `go run main.go serve`  
3. เข้าถึง `http://localhost:8080/metrics` จะเห็นข้อมูล metrics  
4. ตั้งค่า Prometheus server ให้ scrape ที่ endpoint นี้ (เช่น `http://app:8080/metrics`)  
5. สร้าง dashboard ใน Grafana โดยใช้ PromQL query ตัวอย่าง:
   - `rate(icmongolang_api_requests_total[1m])` – request rate  
   - `histogram_quantile(0.95, sum(rate(icmongolang_api_request_duration_seconds_bucket[5m])) by (le, method, path))` – P95 latency

---

## การติดตั้ง

```bash
go get github.com/prometheus/client_golang/prometheus
```

---

## การตั้งค่า configuration

ตัวอย่าง `config-local.yml`:

```yaml
prometheus:
  enabled: true
  metricsPath: "/metrics"
  namespace: "icmongolang"
  subsystem: "api"
```

Environment variables:

```bash
PROMETHEUS_ENABLED=true
PROMETHEUS_METRICSPATH=/metrics
PROMETHEUS_NAMESPACE=icmongolang
PROMETHEUS_SUBSYSTEM=api
```

---

## การรวมกับ GORM (เสริม)

ใช้ GORM เพื่อเก็บข้อมูล metrics ที่ต้องการวิเคราะห์ระยะยาว (เช่น ส่งเข้า InfluxDB หรือ PostgreSQL) แต่ Prometheus จะใช้ scrape แทนการ query ฐานข้อมูลโดยตรง

```go
// ตัวอย่าง: บันทึกการส่ง email แต่ละครั้ง
type EmailMetric struct {
    ID        uint
    Success   bool
    Duration  float64
    CreatedAt time.Time
}
```

---

## Design file / table SQL ที่เกี่ยวข้อง

ไม่มีตาราง SQL สำหรับ Prometheus metrics โดยตรง (เก็บใน TSDB ของ Prometheus)

---

## Return เป็น REST API (endpoint `/metrics`)

Response format: text/plain; version=0.0.4

```
# HELP icmongolang_api_requests_total Total number of HTTP requests
# TYPE icmongolang_api_requests_total counter
icmongolang_api_requests_total{method="GET",path="/api/users",status="200"} 1234
icmongolang_api_requests_total{method="POST",path="/api/auth/login",status="401"} 56
...
```

---

## การใช้งานจริง

**Scenario:** ระบบ API มีผู้ใช้เพิ่มขึ้นเรื่อย ๆ  
- ใช้ Prometheus metrics เพื่อ monitor request rate และ error rate  
- ตั้ง alert เมื่อ error rate > 1% หรือ P95 latency > 500ms  
- เมื่อ alert ทำงาน, ทีม DevOps จะได้รับแจ้งทาง Slack หรือ PagerDuty  
- ใช้ Grafana dashboard เพื่อดูแนวโน้มและวิเคราะห์ performance

---

## ตารางสรุป Components

| Component        | เทคโนโลยี                        | หน้าที่                                     |
|------------------|----------------------------------|---------------------------------------------|
| Prometheus Registry | client_golang/prometheus       | เก็บ metrics ทั้งหมด                       |
| HTTP Middleware  | custom + promhttp               | บันทึก request metrics (duration, count)   |
| Business Metrics | custom Gauge/CounterVec         | เก็บ metrics เฉพาะระบบ (device online)     |
| Metrics Endpoint | promhttp.Handler()              | expose /metrics ให้ Prometheus scrape      |
| Config           | Viper                           | เปิด/ปิด และตั้งค่า namespace               |

---

## แบบฝึกหัดท้าย module (5 ข้อ)

1. **เพิ่ม bucket ที่กำหนดเอง** สำหรับ Histogram ของ request duration (เช่น 0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10)  
2. **เพิ่ม metrics สำหรับ database query** (จำนวน query, duration) โดยใช้ GORM callbacks  
3. **สร้าง Grafana dashboard** ตัวอย่างที่แสดง RPS, error rate, และ P95 latency  
4. **เขียน Prometheus alert rule** สำหรับเงื่อนไข `icmongolang_api_requests_total{status=~"5.."}` rate > 0.1 ต่อนาที  
5. **เพิ่ม middleware** ที่บันทึก request size และ response size เป็น histogram  

---

## แหล่งอ้างอิง

- [Prometheus Go client](https://github.com/prometheus/client_golang)  
- [Prometheus instrumentation best practices](https://prometheus.io/docs/practices/instrumentation/)  
- [PromQL tutorial](https://prometheus.io/docs/prometheus/latest/querying/basics/)  

---

**หมายเหตุ:** module นี้ครบถ้วนสำหรับ `pkg/exporter_prometheus` สามารถเปิด/ปิดได้ผ่าน config และทำงานร่วมกับ Chi router ได้ทันที สามารถเพิ่ม business metrics เพิ่มเติมตามความต้องการของระบบ