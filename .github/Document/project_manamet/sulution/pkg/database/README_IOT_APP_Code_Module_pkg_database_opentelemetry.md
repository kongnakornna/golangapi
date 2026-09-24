# Module 20: pkg/opentelemetry

## สำหรับโฟลเดอร์ `pkg/opentelemetry/`

ไฟล์ที่เกี่ยวข้อง:
- `client.go` - การสร้างและจัดการ TracerProvider, MeterProvider, LoggerProvider
- `trace.go` - การสร้าง spans, propagation, และ manual instrumentation
- `metric.go` - การสร้าง instruments (Counter, Histogram, UpDownCounter)
- `log.go` - การตั้งค่า Log Bridge (otelslog)
- `config.go` - การตั้งค่า exporters (OTLP gRPC/HTTP, stdout)
- `middleware.go` - HTTP/gRPC middleware สำหรับ auto-instrumentation
- `gorm.go` - การรวม OpenTelemetry กับ GORM plugin
- `resource.go` - การตั้งค่า resource attributes (service.name, version, environment)
- `example_main.go` - ตัวอย่างการใช้งานครบวงจร

---

## หลักการ (Concept)

### OpenTelemetry คืออะไร?
OpenTelemetry (หรือ OTel) คือเฟรมเวิร์กสำหรับ observability แบบ open-source ที่รวบรวมมาตรฐานการเก็บ telemetry data (traces, metrics, logs) เข้าด้วยกัน[reference:0]พัฒนาโดย Cloud Native Computing Foundation (CNCF) เพื่อแก้ปัญหาความหลากหลายของ vendor และเครื่องมือ observability[reference:1]OpenTelemetry ไม่ใช่เครื่องมือ backend สำหรับ visualization แต่เป็นกลไกในการ instrument โค้ดแอปพลิเคชันเพื่อทำให้ระบบสามารถสังเกตการณ์ได้ (observable)[reference:2]

OpenTelemetry เกิดขึ้นจากการรวมกันของสองโปรเจกต์หลักคือ OpenTracing และ OpenCensus โดยมีเป้าหมายเพื่อ:
1. **Vendor-agnostic** - เปลี่ยน backend observability ได้โดยไม่ต้องแก้โค้ด
2. **Standardized telemetry** - กำหนด semantic conventions สำหรับชื่อ attributes ที่สอดคล้องกัน
3. **Multiple signals** - รองรับทั้ง traces, metrics, และ logs ใน API เดียว

**ข้อห้ามสำคัญ:** ห้ามใช้ OpenTelemetry แทน logging library โดยตรง เพราะ OpenTelemetry Logs Bridge API ออกแบบมาเพื่อ bridge ระหว่าง logging libraries (เช่น slog, zap, logrus) กับ OpenTelemetry ecosystem ไม่ใช่เป็น logging API เอง หากต้องการ logging ให้ใช้ logging library ปรกติแล้ว bridge เข้า OTel[reference:3]

### มีกี่แบบ? (Deployment Models)

OpenTelemetry มีแนวทางการใช้งานหลักๆ 3 รูปแบบ:

| รูปแบบ | คำอธิบาย | เหมาะกับ |
|--------|----------|----------|
| **Direct to Backend** | SDK ส่ง telemetry ตรงไปยัง backend (Jaeger, Prometheus, Datadog, etc.) | Development, ทดสอบ, production ขนาดเล็ก |
| **Via OpenTelemetry Collector** | SDK ส่งไปยัง Collector ก่อน แล้ว Collector forward ไปยัง backend(s) ต่างๆ | Production (best practice)[reference:4] |
| **Auto-Instrumentation (eBPF)** | ใช้ eBPF เพื่อ instrument โดยไม่ต้องแก้ไขโค้ดหรือ recompile[reference:5] | Legacy applications, ไม่อยากแก้ไขโค้ด |

**รูปแบบของ SDK Exporter:**

| Exporter | Protocol | Port (default) | ใช้กับ |
|----------|----------|----------------|--------|
| **OTLP/gRPC** | gRPC + protobuf | 4317 | Collector, high-performance[reference:6] |
| **OTLP/HTTP** | HTTP + protobuf/JSON | 4318 | Firewall-friendly, HTTP-only environments[reference:7] |
| **stdout** | Console output | - | Development, debugging[reference:8] |
| **Jaeger** | Jaeger thrift | 6831, 14268 | Legacy Jaeger instances |

### ใช้อย่างไร / นำไปใช้กรณีไหน

**กรณีใช้งาน:**
- **Distributed tracing** - 追踪 request ข้าม microservices เพื่อหา bottleneck และ troubleshoot
- **Performance monitoring** - เก็บ latency, throughput, error rate ของ API และ database
- **Correlate logs with traces** - เชื่อมโยง logs เข้ากับ trace IDs เพื่อ debugging ที่มีประสิทธิภาพ[reference:9]
- **Service dependency mapping** - visualize service topology อัตโนมัติ
- **SLO/SLI monitoring** - เก็บ metrics สำหรับ service level objectives

**รูปแบบการใช้งาน OpenTelemetry ใน Go:**
1. **Manual instrumentation** - ใช้ OpenTelemetry API โดยตรง ควบคุมได้ละเอียด[reference:10]
2. **Instrumentation libraries** - ใช้ middleware/library ที่มีอยู่แล้วสำหรับ net/http, gRPC, GORM, Gin[reference:11]
3. **Auto-instrumentation (eBPF)** - ใช้ eBPF โดยไม่ต้องแก้ไขโค้ด[reference:12]

### ประโยชน์ที่ได้รับ
- **Vendor lock-in prevention** - เปลี่ยน backend observability ได้โดยไม่แก้โค้ด (เปลี่ยนแค่ exporter config)
- **Standardized telemetry** - semantic conventions ทำให้เครื่องมือต่างๆ อ่านข้อมูล telemetry ได้ถูกต้อง[reference:13]
- **Single API for three signals** - เรียนรู้ API ชุดเดียวสำหรับ traces, metrics, logs
- **Context propagation อัตโนมัติ** - trace context ถูก propagate ผ่าน HTTP headers, gRPC metadata
- **Rich ecosystem** - instrumentation libraries สำหรับทุก major frameworks (Gin, Echo, gRPC, GORM)
- **Active community** - CNCF graduated project มีการพัฒนาต่อเนื่อง

### ข้อควรระวัง
- **Performance overhead** - การสร้าง spans และ metrics มี overhead (โดยเฉพาะ在高 throughput scenarios)
- **Sampling** - ควรใช้ sampling เพื่อลด overhead ใน production (head sampling หรือ tail sampling)
- **Cardinality explosion** - ระวังการใช้ attribute values ที่มี cardinality สูง (เช่น user_id, session_id)
- **Logs bridge complexity** - Logs API ยังอยู่ในสถานะ experimental (beta) มีโอกาสเปลี่ยนแปลง[reference:14]
- **eBPF instrumentation** - ต้องการ kernel version สูง (>=4.18) และต้อง enable feature gate[reference:15]
- **Context propagation** - ต้องมั่นใจว่า context ถูก propagate อย่างถูกต้องข้าม goroutines และ services

### ข้อดี
- **Standard de facto** - เป็นมาตรฐาน open-source ที่ได้รับการยอมรับมากที่สุดในวงการ observability
- **Native integration** - ใช้ร่วมกับ Grafana, Jaeger, Prometheus, Datadog, New Relic ได้
- **Collector เป็นตัวกรองและ enrich data** - ก่อนส่งไปยัง multiple backends
- **Semantic conventions** - ช่วยให้การ query และ correlation ทำได้ง่าย
- **Support for OpenTelemetry Protocol (OTLP)** - protocol ที่ efficient และ extensible

### ข้อเสีย
- **Learning curve** - ต้องเข้าใจ concepts หลายอย่าง (spans, traces, metrics instruments, exporters)
- **Logging complexity** - Logs API ยังไม่ stable เท่า traces และ metrics
- **eBPF ยังเป็น beta** - Go auto-instrumentation ด้วย eBPF ยังอยู่ในช่วง beta และต้องการ feature gate[reference:16]
- **Configuration boilerplate** - การตั้งค่า SDK ตั้งแต่เริ่มต้นมีโค้ดจำนวนมาก (แต่ใช้ helper libraries ได้)
- **Documentation กระจาย** - เอกสารมีหลายส่วน (spec, SDK, contrib, collector) อาจทำให้สับสน

### ข้อห้าม
**ห้ามใช้ OpenTelemetry เป็น primary logging API แทน slog, zap, logrus** เพราะ OpenTelemetry Logs Bridge API ออกแบบมาให้เป็น bridge ระหว่าง existing logging libraries กับ OpenTelemetry ecosystem ไม่ใช่ logging API สำหรับ application โดยตรง[reference:17]การใช้ OpenTelemetry logs โดยตรงจะทำให้:
1. ขาด features ของ logging libraries (levels, structured logging, output formats)
2. ต้อง migrate โค้ด logging ทั้งหมดเมื่อเปลี่ยน vendor
3. ไม่สามารถใช้ประโยชน์จาก logging ecosystem ที่มีอยู่ได้

**แนวทางที่ถูกต้อง:** ใช้ slog, zap, หรือ logrus สำหรับ logging แล้วใช้ otelslog, otelzap เป็น bridge เพื่อ inject trace context[reference:18]

---

## การออกแบบ Workflow และ Dataflow

```mermaid
flowchart TB
    subgraph Application["Go Application"]
        A[Business Logic]
        B[HTTP Handler]
        C[Database Call]
        D[gRPC Client]
    end
    
    subgraph Instrumentation["OpenTelemetry SDK"]
        E[Trace API]
        F[Metric API]
        G[Log Bridge<br/>(otelslog)]
        H[Context Propagation]
    end
    
    subgraph Exporters["Exporters"]
        I[OTLP/gRPC<br/>:4317]
        J[OTLP/HTTP<br/>:4318]
        K[stdout]
    end
    
    subgraph Collector["OpenTelemetry Collector"]
        L[Receiver]
        M[Processor<br/>(Batch, Filter, Sampling)]
        N[Exporter]
    end
    
    subgraph Backend["Backend & Visualization"]
        O[Jaeger]
        P[Prometheus]
        Q[Grafana]
    end
    
    A & B & C & D --> E & F & G
    E & F & G --> H
    H --> I & J & K
    I & J --> L
    L --> M --> N
    N --> O & P
    O & P --> Q
```

**Dataflow ใน Go application (Push model with Collector):**
1. **Initialize SDK** - สร้าง TracerProvider, MeterProvider, LoggerProvider พร้อม exporters
2. **Set global providers** - ใช้ `otel.SetTracerProvider()`, `otel.SetMeterProvider()`
3. **Instrument code** - เรียก `tracer.Start()` เพื่อสร้าง spans, record metrics ด้วย instruments
4. **Propagate context** - ใช้ context ในการส่งผ่าน trace information ข้าม functions/services
5. **Export telemetry** - SDK ส่ง telemetry ไปยัง Collector (หรือ backend โดยตรง)
6. **Visualize** - Grafana query จาก Jaeger/Prometheus เพื่อแสดงข้อมูล


## ตัวอย่างโค้ดที่รันได้จริง

### โครงสร้างโปรเจกต์
```
pkg/opentelemetry/
├── client.go          # SDK initialization
├── trace.go           # Span and tracer utilities
├── metric.go          # Metric instruments
├── log.go             # Log bridge setup
├── config.go          # Configuration management
├── middleware.go      # HTTP/gRPC middleware
├── gorm.go            # GORM integration
├── resource.go        # Resource attributes
├── config.go          # Config structs
└── example_main.go    # Complete example
```

### 1. การติดตั้ง Go packages

```bash
# Core SDK
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/sdk
go get go.opentelemetry.io/otel/trace
go get go.opentelemetry.io/otel/metric
go get go.opentelemetry.io/otel/sdk/metric
go get go.opentelemetry.io/otel/sdk/trace

# Exporters
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp
go get go.opentelemetry.io/otel/exporters/stdout/stdouttrace
go get go.opentelemetry.io/otel/exporters/stdout/stdoutmetric
go get go.opentelemetry.io/otel/exporters/prometheus

# Contrib instrumentation
go get go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp
go get go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc
go get go.opentelemetry.io/contrib/bridges/otelslog

# GORM plugin
go get gorm.io/plugin/opentelemetry/tracing

# Semantic conventions
go get go.opentelemetry.io/otel/semconv/v1.26.0
```

### 2. การติดตั้ง OpenTelemetry Collector (Option - Recommended)

```yaml
# docker-compose.yml
version: '3.8'
services:
  # OpenTelemetry Collector
  otel-collector:
    image: otel/opentelemetry-collector-contrib:latest
    container_name: otel-collector
    command: ["--config=/etc/otel-collector-config.yaml"]
    volumes:
      - ./otel-collector-config.yaml:/etc/otel-collector-config.yaml
    ports:
      - "4317:4317"   # OTLP gRPC receiver
      - "4318:4318"   # OTLP HTTP receiver
      - "8888:8888"   # Prometheus metrics (self)
    restart: unless-stopped

  # Jaeger (trace backend)
  jaeger:
    image: jaegertracing/all-in-one:latest
    container_name: jaeger
    ports:
      - "16686:16686"   # UI
      - "14250:14250"   # gRPC receiver for Collector
    environment:
      - COLLECTOR_OTLP_ENABLED=true
    restart: unless-stopped

  # Prometheus (metrics backend)
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    restart: unless-stopped

  # Grafana
  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana_data:/var/lib/grafana
    restart: unless-stopped

volumes:
  grafana_data:
```

Configuration Collector (`otel-collector-config.yaml`):
```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch:
    timeout: 1s
    send_batch_size: 1024
  memory_limiter:
    check_interval: 1s
    limit_mib: 512

exporters:
  jaeger:
    endpoint: jaeger:14250
    tls:
      insecure: true
  prometheus:
    endpoint: "0.0.0.0:8889"
    namespace: myapp

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [jaeger]
    metrics:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [prometheus]
```

รันด้วย:
```bash
docker-compose up -d
```

### 3. ตัวอย่างโค้ด: Configuration

```go
// config.go
package opentelemetry

import (
    "os"
    "time"
)

type Config struct {
    // Service identification
    ServiceName    string
    ServiceVersion string
    Environment    string
    
    // Exporter settings
    ExporterType   string   // "otlp-grpc", "otlp-http", "stdout", "jaeger"
    OTLPEndpoint   string
    OTLPInsecure   bool
    
    // Trace settings
    TraceSamplingRatio float64
    TraceBatchTimeout  time.Duration
    
    // Metric settings
    MetricInterval time.Duration
    
    // Log settings
    LogLevel string
}

func DefaultConfig() Config {
    return Config{
        ServiceName:        "myapp",
        ServiceVersion:     "1.0.0",
        Environment:        "development",
        ExporterType:       "otlp-grpc",
        OTLPEndpoint:       "localhost:4317",
        OTLPInsecure:       true,
        TraceSamplingRatio: 1.0,
        TraceBatchTimeout:  5 * time.Second,
        MetricInterval:     15 * time.Second,
        LogLevel:           "info",
    }
}

func LoadConfigFromEnv() Config {
    cfg := DefaultConfig()
    
    if name := os.Getenv("OTEL_SERVICE_NAME"); name != "" {
        cfg.ServiceName = name
    }
    if version := os.Getenv("OTEL_SERVICE_VERSION"); version != "" {
        cfg.ServiceVersion = version
    }
    if env := os.Getenv("ENVIRONMENT"); env != "" {
        cfg.Environment = env
    }
    if endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); endpoint != "" {
        cfg.OTLPEndpoint = endpoint
    }
    
    return cfg
}
```

### 4. ตัวอย่างโค้ด: Resource Attributes

```go
// resource.go
package opentelemetry

import (
    "context"
    
    "go.opentelemetry.io/otel/sdk/resource"
    semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// NewResource creates resource with service metadata
func NewResource(cfg Config) (*resource.Resource, error) {
    return resource.New(context.Background(),
        // Pull from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME env vars
        resource.WithFromEnv(),
        // Add SDK information
        resource.WithTelemetrySDK(),
        // Add custom resource attributes
        resource.WithAttributes(
            semconv.ServiceName(cfg.ServiceName),
            semconv.ServiceVersion(cfg.ServiceVersion),
            semconv.DeploymentEnvironment(cfg.Environment),
        ),
    )
}
```

### 5. ตัวอย่างโค้ด: Client (SDK Initialization)

```go
// client.go
package opentelemetry

import (
    "context"
    "fmt"
    "time"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
    "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
    "go.opentelemetry.io/otel/sdk/trace"
    sdkmetric "go.opentelemetry.io/otel/sdk/metric"
    "go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type OpenTelemetryClient struct {
    TracerProvider *trace.TracerProvider
    MeterProvider  *sdkmetric.MeterProvider
    ShutdownFunc   func(context.Context) error
}

// InitializeSDK initializes OpenTelemetry SDK with configured exporters
func InitializeSDK(ctx context.Context, cfg Config) (*OpenTelemetryClient, error) {
    // Create resource
    res, err := NewResource(cfg)
    if err != nil {
        return nil, fmt.Errorf("failed to create resource: %w", err)
    }
    
    // Create trace exporter
    traceExporter, err := createTraceExporter(cfg)
    if err != nil {
        return nil, fmt.Errorf("failed to create trace exporter: %w", err)
    }
    
    // Create trace provider
    tp := trace.NewTracerProvider(
        trace.WithBatcher(traceExporter,
            trace.WithBatchTimeout(cfg.TraceBatchTimeout),
        ),
        trace.WithResource(res),
        trace.WithSampler(trace.TraceIDRatioBased(cfg.TraceSamplingRatio)),
    )
    
    // Create metric exporter
    metricExporter, err := createMetricExporter(cfg)
    if err != nil {
        return nil, fmt.Errorf("failed to create metric exporter: %w", err)
    }
    
    // Create meter provider
    mp := sdkmetric.NewMeterProvider(
        sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter,
            sdkmetric.WithInterval(cfg.MetricInterval),
        )),
        sdkmetric.WithResource(res),
    )
    
    // Set global providers
    otel.SetTracerProvider(tp)
    otel.SetMeterProvider(mp)
    
    return &OpenTelemetryClient{
        TracerProvider: tp,
        MeterProvider:  mp,
        ShutdownFunc: func(ctx context.Context) error {
            if err := tp.Shutdown(ctx); err != nil {
                return err
            }
            if err := mp.Shutdown(ctx); err != nil {
                return err
            }
            return nil
        },
    }, nil
}

func createTraceExporter(cfg Config) (trace.SpanExporter, error) {
    switch cfg.ExporterType {
    case "otlp-grpc":
        return otlptracegrpc.New(context.Background(),
            otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint),
            otlptracegrpc.WithInsecure(),
        )
    case "otlp-http":
        return otlptracehttp.New(context.Background(),
            otlptracehttp.WithEndpoint(cfg.OTLPEndpoint),
            otlptracehttp.WithInsecure(),
        )
    case "stdout":
        return stdouttrace.New(stdouttrace.WithPrettyPrint())
    default:
        return otlptracegrpc.New(context.Background(),
            otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint),
            otlptracegrpc.WithInsecure(),
        )
    }
}

func createMetricExporter(cfg Config) (sdkmetric.Exporter, error) {
    // For simplicity, return stdout exporter
    // In production, use OTLP exporter similar to trace
    return &stdoutMetricExporter{}, nil
}

type stdoutMetricExporter struct{}

func (e *stdoutMetricExporter) Export(ctx context.Context, rm *sdkmetric.ResourceMetrics) error {
    fmt.Printf("Metrics exported: %+v\n", rm)
    return nil
}

func (e *stdoutMetricExporter) ForceFlush(ctx context.Context) error {
    return nil
}

func (e *stdoutMetricExporter) Shutdown(ctx context.Context) error {
    return nil
}

func (e *stdoutMetricExporter) Temporality(ik sdkmetric.InstrumentKind) metricdata.Temporality {
    return metricdata.DeltaTemporality
}

func (e *stdoutMetricExporter) Aggregation(ik sdkmetric.InstrumentKind) sdkmetric.Aggregation {
    return sdkmetric.DefaultAggregationSelector(ik)
}
```

### 6. ตัวอย่างโค้ด: Trace Utilities

```go
// trace.go
package opentelemetry

import (
    "context"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/trace"
)

// GetTracer returns a tracer with the given name
func GetTracer(name string) trace.Tracer {
    return otel.Tracer(name)
}

// StartSpan creates a new span with attributes
func StartSpan(ctx context.Context, tracerName, spanName string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
    tracer := GetTracer(tracerName)
    return tracer.Start(ctx, spanName, trace.WithAttributes(attrs...))
}

// AddSpanError records an error on the current span
func AddSpanError(span trace.Span, err error) {
    if err == nil {
        return
    }
    span.RecordError(err)
    span.SetStatus(codes.Error, err.Error())
}

// AddSpanAttributes adds attributes to current span
func AddSpanAttributes(span trace.Span, attrs ...attribute.KeyValue) {
    span.SetAttributes(attrs...)
}

// Example function with manual instrumentation
func ExampleInstrumentedFunction(ctx context.Context) {
    ctx, span := StartSpan(ctx, "myapp", "ExampleInstrumentedFunction",
        attribute.String("custom.attribute", "value"),
    )
    defer span.End()
    
    // Business logic here...
    
    // If error occurs
    // AddSpanError(span, err)
}
```

### 7. ตัวอย่างโค้ด: Metrics Instruments

```go
// metric.go
package opentelemetry

import (
    "context"
    "sync"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/metric"
)

// Metrics struct holds all application metrics
type Metrics struct {
    // Counter - only increases
    RequestsTotal   metric.Int64Counter
    
    // UpDownCounter - can increase and decrease
    ActiveRequests  metric.Int64UpDownCounter
    
    // Histogram - for latency distributions
    RequestDuration metric.Float64Histogram
    
    // Gauge (via Async instruments)
    QueueSize       metric.Int64Gauge
    
    mu sync.RWMutex
}

// NewMetrics creates and registers all metrics
func NewMetrics(serviceName string) (*Metrics, error) {
    meter := otel.Meter(serviceName)
    
    m := &Metrics{}
    
    var err error
    
    // Create counter
    m.RequestsTotal, err = meter.Int64Counter(
        "requests_total",
        metric.WithDescription("Total number of HTTP requests"),
        metric.WithUnit("{request}"),
    )
    if err != nil {
        return nil, err
    }
    
    // Create up down counter
    m.ActiveRequests, err = meter.Int64UpDownCounter(
        "active_requests",
        metric.WithDescription("Number of active requests"),
        metric.WithUnit("{request}"),
    )
    if err != nil {
        return nil, err
    }
    
    // Create histogram
    m.RequestDuration, err = meter.Float64Histogram(
        "request_duration_seconds",
        metric.WithDescription("Request duration in seconds"),
        metric.WithUnit("s"),
    )
    if err != nil {
        return nil, err
    }
    
    // Create gauge (async callback)
    m.QueueSize, err = meter.Int64Gauge(
        "queue_size",
        metric.WithDescription("Current queue size"),
        metric.WithUnit("{item}"),
    )
    if err != nil {
        return nil, err
    }
    
    return m, nil
}

// RecordRequest records metrics for an HTTP request
func (m *Metrics) RecordRequest(ctx context.Context, method, endpoint, status string, duration float64) {
    attrs := []attribute.KeyValue{
        attribute.String("http.method", method),
        attribute.String("http.route", endpoint),
        attribute.String("http.status_code", status),
    }
    
    m.RequestsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
    m.RequestDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
}

// AddActiveRequest increments active request counter
func (m *Metrics) AddActiveRequest(ctx context.Context) {
    m.ActiveRequests.Add(ctx, 1)
}

// RemoveActiveRequest decrements active request counter
func (m *Metrics) RemoveActiveRequest(ctx context.Context) {
    m.ActiveRequests.Add(ctx, -1)
}

// RegisterQueueGauge registers a callback for queue size gauge
func (m *Metrics) RegisterQueueGauge(getQueueSize func() int64) error {
    _, err := m.QueueSize.RegisterCallback(func(ctx context.Context, obs metric.Int64Observer) error {
        obs.Observe(getQueueSize())
        return nil
    })
    return err
}
```

### 8. ตัวอย่างโค้ด: Log Bridge (otelslog)

```go
// log.go
package opentelemetry

import (
    "context"
    "log/slog"
    "os"
    
    "go.opentelemetry.io/contrib/bridges/otelslog"
    "go.opentelemetry.io/otel"
)

// SetupLogBridge configures slog with OpenTelemetry bridge
// This injects trace context into log records automatically
func SetupLogBridge(cfg Config) error {
    // Create OTel handler that wraps slog
    otelHandler := otelslog.NewHandler(otel.GetTracerProvider())
    
    // Create JSON handler for output
    jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: getLogLevel(cfg.LogLevel),
    })
    
    // Chain handlers: OTel -> JSON
    handler := otelslog.NewHandlerWithConfig(
        otel.GetTracerProvider(),
        otelslog.HandlerConfig{
            // Inject trace/span IDs into log records automatically
            WithSpanContext: true,
        },
    )
    
    // Or use multi-handler to write to both
    multiHandler := slog.NewMultiHandler(otelHandler, jsonHandler)
    slog.SetDefault(slog.New(multiHandler))
    
    return nil
}

func getLogLevel(level string) slog.Level {
    switch level {
    case "debug":
        return slog.LevelDebug
    case "info":
        return slog.LevelInfo
    case "warn":
        return slog.LevelWarn
    case "error":
        return slog.LevelError
    default:
        return slog.LevelInfo
    }
}

// LogWithTraceContext logs a message with automatic trace context injection
func LogWithTraceContext(ctx context.Context, level slog.Level, msg string, args ...any) {
    // Use slog.LogAttrs which automatically extracts trace context from ctx
    // when using otelslog bridge
    slog.Log(ctx, level, msg, args...)
}
```

### 9. ตัวอย่างโค้ด: HTTP Middleware (otelhttp)

```go
// middleware.go
package opentelemetry

import (
    "net/http"
    
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/propagation"
)

// SetupGlobalPropagator sets up context propagation for HTTP headers
func SetupGlobalPropagator() {
    otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{},
        propagation.Baggage{},
    ))
}

// NewHTTPHandler wraps an http.Handler with OpenTelemetry instrumentation
// Automatically creates spans and records metrics for each request
func NewHTTPHandler(handler http.Handler, operation string, opts ...otelhttp.Option) http.Handler {
    defaultOpts := []otelhttp.Option{
        otelhttp.WithTracerProvider(otel.GetTracerProvider()),
        otelhttp.WithMeterProvider(otel.GetMeterProvider()),
        otelhttp.WithPropagators(otel.GetTextMapPropagator()),
        otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
            return operation + " " + r.Method
        }),
    }
    defaultOpts = append(defaultOpts, opts...)
    return otelhttp.NewHandler(handler, operation, defaultOpts...)
}

// NewHTTPClient returns an http.Client with OpenTelemetry instrumentation
// Automatically injects trace context into outgoing requests
func NewHTTPClient() *http.Client {
    return &http.Client{
        Transport: otelhttp.NewTransport(http.DefaultTransport),
    }
}

// HTTPMiddleware returns middleware for standard net/http
func HTTPMiddleware(metrics *Metrics) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return NewHTTPHandler(next, "http.server")
    }
}
```

### 10. ตัวอย่างโค้ด: GORM Integration

```go
// gorm.go
package opentelemetry

import (
    "gorm.io/gorm"
    "gorm.io/plugin/opentelemetry/tracing"
)

// SetupGORMWithTracing initializes GORM with OpenTelemetry tracing
// Automatically traces all database operations including queries,
// creates, updates, deletes, and transactions[reference:19]
func SetupGORMWithTracing(db *gorm.DB) error {
    // Create OTel plugin with options
    plugin := tracing.NewPlugin(
        tracing.WithTracerProvider(otel.GetTracerProvider()),
        tracing.WithDBName("myapp_db"),
    )
    
    // Use the plugin
    if err := db.Use(plugin); err != nil {
        return err
    }
    
    return nil
}

// GORMPluginConfig allows custom configuration
type GORMPluginConfig struct {
    DBName         string
    ServiceName    string
    IncludeQuery   bool  // Include SQL query in span attributes
}

// SetupGORMWithCustomTracing configures GORM with custom options
func SetupGORMWithCustomTracing(db *gorm.DB, cfg GORMPluginConfig) error {
    opts := []tracing.Option{
        tracing.WithTracerProvider(otel.GetTracerProvider()),
        tracing.WithDBName(cfg.DBName),
    }
    
    if !cfg.IncludeQuery {
        opts = append(opts, tracing.WithoutQueryVariables())
    }
    
    plugin := tracing.NewPlugin(opts...)
    return db.Use(plugin)
}
```

### 11. ตัวอย่างการใช้งานรวมใน HTTP server

```go
// main.go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"
    
    "yourproject/pkg/opentelemetry"
    
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"
)

func main() {
    // Load config
    cfg := opentelemetry.LoadConfigFromEnv()
    
    // Initialize OpenTelemetry SDK
    ctx := context.Background()
    otelClient, err := opentelemetry.InitializeSDK(ctx, cfg)
    if err != nil {
        log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
    }
    defer func() {
        shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        otelClient.ShutdownFunc(shutdownCtx)
    }()
    
    // Setup global propagator for context propagation
    opentelemetry.SetupGlobalPropagator()
    
    // Setup log bridge
    if err := opentelemetry.SetupLogBridge(cfg); err != nil {
        log.Printf("Warning: Failed to setup log bridge: %v", err)
    }
    
    // Create metrics
    metrics, err := opentelemetry.NewMetrics(cfg.ServiceName)
    if err != nil {
        log.Fatalf("Failed to create metrics: %v", err)
    }
    
    // Create HTTP server with middleware
    mux := http.NewServeMux()
    
    // Instrumented endpoints
    mux.Handle("/api/users", opentelemetry.NewHTTPHandler(
        http.HandlerFunc(handleUsers(metrics)),
        "GET /api/users",
    ))
    
    mux.Handle("/api/health", opentelemetry.NewHTTPHandler(
        http.HandlerFunc(handleHealth),
        "GET /api/health",
    ))
    
    // Start server
    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }
    
    go func() {
        log.Println("Server starting on :8080")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()
    
    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)
    <-quit
    
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    server.Shutdown(shutdownCtx)
}

func handleUsers(metrics *opentelemetry.Metrics) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract current span from context
        span := trace.SpanFromContext(r.Context())
        defer span.End()
        
        // Add custom attributes to span
        span.SetAttributes(attribute.String("user.id", "12345"))
        
        // Increment active requests counter
        metrics.AddActiveRequest(r.Context())
        defer metrics.RemoveActiveRequest(r.Context())
        
        // Record request duration
        start := time.Now()
        defer func() {
            duration := time.Since(start).Seconds()
            metrics.RecordRequest(r.Context(), r.Method, "/api/users", "200", duration)
        }()
        
        // Simulate database call with span
        _, dbSpan := opentelemetry.StartSpan(r.Context(), "db", "SELECT * FROM users")
        time.Sleep(50 * time.Millisecond) // Simulate DB query
        dbSpan.End()
        
        // Log with trace context
        slog.InfoContext(r.Context(), "Handled users request")
        
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"users":[]}`))
    }
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ok"}`))
}
```


## วิธีใช้งาน module นี้

1. **ติดตั้ง OpenTelemetry Collector และ backends** (Jaeger, Prometheus, Grafana) ตาม docker-compose ข้างต้น
2. **ติดตั้ง Go packages** ตามที่ระบุในหัวข้อ "การติดตั้ง Go packages"
3. **คัดลอกโค้ด** ไฟล์ `client.go`, `trace.go`, `metric.go`, `log.go`, `middleware.go`, `gorm.go`, `resource.go`, `config.go` ไปไว้ใน `pkg/opentelemetry/`
4. **ปรับ configuration** ตาม environment ของคุณ (ใช้ environment variables)
5. **Initialize SDK** ที่จุดเริ่มต้นของ application
6. **Wrap HTTP handlers** ด้วย `opentelemetry.NewHTTPHandler()`
7. **ใช้ context propagation** โดยส่ง context ผ่าน functions และ services
8. **Record metrics** ด้วย metrics struct
9. **Log ด้วย slog** เพื่อให้ trace context ถูก inject อัตโนมัติ


## การติดตั้ง

```bash
# Create module
go mod init yourproject

# Install core OpenTelemetry packages
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/sdk
go get go.opentelemetry.io/otel/trace
go get go.opentelemetry.io/otel/metric
go get go.opentelemetry.io/otel/sdk/metric
go get go.opentelemetry.io/otel/sdk/trace

# Install exporters
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp
go get go.opentelemetry.io/otel/exporters/stdout/stdouttrace

# Install instrumentation libraries
go get go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp
go get go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc
go get go.opentelemetry.io/contrib/bridges/otelslog

# Install GORM plugin
go get gorm.io/plugin/opentelemetry/tracing

# Install semantic conventions
go get go.opentelemetry.io/otel/semconv/v1.26.0

# Install OpenTelemetry Collector builder (optional)
go install github.com/open-telemetry/opentelemetry-collector-contrib/cmd/telemetrygen@latest

# For Docker setup
docker pull otel/opentelemetry-collector-contrib:latest
docker pull jaegertracing/all-in-one:latest
docker pull prom/prometheus:latest
docker pull grafana/grafana:latest
```


## การตั้งค่า configuration

### Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `OTEL_SERVICE_NAME` | Service name for telemetry | `myapp` |
| `OTEL_SERVICE_VERSION` | Service version | `1.0.0` |
| `ENVIRONMENT` | Deployment environment | `production` |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OTLP endpoint URL | `http://localhost:4317` |
| `OTEL_RESOURCE_ATTRIBUTES` | Additional resource attributes | `deployment.environment=prod` |
| `OTEL_TRACES_SAMPLER` | Sampling strategy | `parentbased_traceidratio` |
| `OTEL_TRACES_SAMPLER_ARG` | Sampling ratio (0-1) | `0.1` |
| `OTEL_EXPORTER_OTLP_HEADERS` | Custom headers for exporter | `X-API-Key=123` |

### Config struct

```go
// config.go
type Config struct {
    ServiceName        string
    ServiceVersion     string
    Environment        string
    ExporterType       string   // "otlp-grpc", "otlp-http", "stdout"
    OTLPEndpoint       string
    OTLPInsecure       bool
    TraceSamplingRatio float64
    TraceBatchTimeout  time.Duration
    MetricInterval     time.Duration
    LogLevel           string
}
```

### Loading config

```go
// Load from environment variables (recommended)
cfg := opentelemetry.LoadConfigFromEnv()

// Or load from config file
cfg := opentelemetry.DefaultConfig()
cfg.ServiceName = "myapp"
cfg.Environment = "production"
```


## การรวมกับ GORM

```go
import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/plugin/opentelemetry/tracing"
)

func initDB() (*gorm.DB, error) {
    // Connect to database
    db, err := gorm.Open(postgres.Open("postgres://..."), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    // Add OpenTelemetry plugin
    plugin := tracing.NewPlugin(
        tracing.WithTracerProvider(otel.GetTracerProvider()),
        tracing.WithDBName("myapp_db"),
    )
    if err := db.Use(plugin); err != nil {
        return nil, err
    }
    
    return db, nil
}

// Example usage with context
func getUserByID(ctx context.Context, db *gorm.DB, id string) (*User, error) {
    var user User
    // GORM automatically creates spans for this query
    err := db.WithContext(ctx).Where("id = ?", id).First(&user).Error
    return &user, err
}
```

**หมายเหตุ:** GORM OpenTelemetry plugin จะ trace ทุก database operations รวมถึง queries, creates, updates, deletes และ transactions โดยอัตโนมัติ[reference:20]


## การใช้งานจริง

### Example 1: Basic HTTP Server with Tracing and Metrics

```go
package main

import (
    "net/http"
    "yourproject/pkg/opentelemetry"
)

func main() {
    // Setup OTel
    cfg := opentelemetry.DefaultConfig()
    otelClient, _ := opentelemetry.InitializeSDK(context.Background(), cfg)
    defer otelClient.ShutdownFunc(context.Background())
    
    opentelemetry.SetupGlobalPropagator()
    metrics, _ := opentelemetry.NewMetrics(cfg.ServiceName)
    
    // Create server with instrumentation
    mux := http.NewServeMux()
    mux.Handle("/api/", opentelemetry.HTTPMiddleware(metrics)(http.DefaultServeMux))
    
    http.ListenAndServe(":8080", mux)
}
```

### Example 2: Distributed Tracing Across Services

```go
// Service A (client)
func callServiceB(ctx context.Context, url string) {
    client := opentelemetry.NewHTTPClient()
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    client.Do(req)  // Trace context automatically propagated via headers
}

// Service B (server) - automatically extracts trace context
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    // span automatically created and linked to parent trace
    // business logic...
}
```

### Example 3: Background Job with Tracing

```go
func processBatch(ctx context.Context, batch []Data) {
    ctx, span := opentelemetry.StartSpan(ctx, "processor", "ProcessBatch")
    defer span.End()
    
    for i, item := range batch {
        // Create child span for each item
        _, itemSpan := opentelemetry.StartSpan(ctx, "processor", "ProcessItem",
            attribute.Int("item.index", i),
        )
        processItem(item)
        itemSpan.End()
    }
}
```


## ตารางสรุป OpenTelemetry Components

| Component | คำอธิบาย | ตัวอย่าง |
|-----------|----------|----------|
| **TracerProvider** | 负责创建 Tracer (span generator) | `trace.NewTracerProvider()` |
| **Tracer** | ใช้สร้าง spans สำหรับ tracing | `otel.Tracer("myapp")` |
| **Span** | ตัวแทน unit of work มี start time, end time, attributes, events | `tracer.Start(ctx, "operation")` |
| **SpanContext** | ข้อมูลสำหรับ span identification (TraceID, SpanID, TraceFlags) | `span.SpanContext()` |
| **Propagator** | 负责 inject/extract trace context ข้าม services | `propagation.TraceContext{}` |
| **MeterProvider** | 负责สร้าง instruments สำหรับ metrics | `metric.NewMeterProvider()` |
| **Meter** | ใช้สร้าง metric instruments | `otel.Meter("myapp")` |
| **Counter** | Metric ที่เพิ่มขึ้นอย่างเดียว (cumulative) | `meter.Int64Counter("requests")`[reference:21] |
| **UpDownCounter** | Metric ที่เพิ่มและลดได้ | `meter.Int64UpDownCounter("active")`[reference:22] |
| **Histogram** | ใช้วัด distribution ของค่า (latency, size) | `meter.Float64Histogram("duration")`[reference:23] |
| **LoggerProvider** | 负责สร้าง Logger bridges | `log.NewLoggerProvider()` |
| **Log Bridge** | เชื่อมต่อ logging library เข้ากับ OTel (otelslog, otelzap) | `otelslog.NewHandler()`[reference:24] |
| **Exporter** | ส่ง telemetry ไปยัง backend หรือ Collector | `otlptracegrpc.New()` |
| **OTLP Collector** | Central component สำหรับรับ, process, ส่ง telemetry | OpenTelemetry Collector |
| **Resource** | Metadata ของ service ที่สร้าง telemetry | `resource.NewWithAttributes()` |
| **Semantic Conventions** | Standardized attribute names | `semconv.HTTPMethodKey`[reference:25] |
| **Batch Processor** | รวม spans/metrics ก่อนส่งเพื่อลด overhead | `trace.WithBatcher()` |


## แบบฝึกหัดท้าย module (5 ข้อ)

### ข้อ 1: การสร้างและใช้ Custom Span

จงเขียนฟังก์ชัน `ProcessOrder(ctx context.Context, orderID string, amount float64) error` ที่:
- สร้าง span ชื่อ `ProcessOrder` ด้วย tracer name "order-service"
- เพิ่ม attributes: `order.id`, `order.amount`
- Simulate work: สร้าง child span สำหรับ `ValidatePayment` และ `UpdateInventory`
- ถ้าเกิด error ให้บันทึก error บน span และ set status เป็น Error
- Record metrics: `orders_processed_total` (Counter) พร้อม labels `status` (success/failed)

### ข้อ 2: การตั้งค่า OpenTelemetry Collector และ Backend

จาก docker-compose ที่มี OpenTelemetry Collector, Jaeger, Prometheus, Grafana:
- จงเขียน Collector configuration (YAML) ที่รับ OTLP gRPC และ OTLP HTTP แล้วส่ง traces ไป Jaeger และ metrics ไป Prometheus
- เพิ่ม batch processor ที่ส่ง batch ทุก 2 วินาที หรือมีขนาด 2048 items
- ตั้งค่า memory limiter ที่ limit 256 MiB
- ทดสอบโดยใช้ `telemetrygen` สร้าง synthetic telemetry และ verify ใน Jaeger UI

### ข้อ 3: การสร้าง HTTP Middleware แบบกำหนดเอง

จงสร้าง middleware `otelMiddleware` ที่:
- สร้าง span สำหรับทุก incoming request (ใช้ `otelhttp.NewHandler` หรือ implement เอง)
- บันทึก HTTP method, URL path, status code, response size เป็น span attributes
- Record metrics: request duration (histogram) และ request count (counter)
- สำหรับ requests ที่ใช้เวลา > 1 วินาที ให้บันทึก event "slow_request" บน span
- ทดสอบโดยสร้าง HTTP server ด้วย middleware นี้และส่ง requests

### ข้อ 4: การรวม OpenTelemetry กับ GORM และ Correlation

จาก GORM model `User` และ `Order`:
- จง integrate GORM กับ OpenTelemetry โดยใช้ `gorm.io/plugin/opentelemetry/tracing`
- เขียนฟังก์ชัน `CreateOrderWithUser(ctx context.Context, db *gorm.DB, user User, order Order)` ที่:
  - สร้าง parent span "CreateOrderTransaction"
  - ใช้ `db.WithContext(ctx)` เพื่อ propagate trace context
  - สร้าง user และ order ใน database transaction
  - ถ้ามี error ให้ rollback transaction และบันทึก error บน span
- แสดงว่า trace สามารถเชื่อมโยงการสร้าง user และ order ใน transaction เดียวกันได้อย่างไร

### ข้อ 5: การตั้งค่า Sampling และ Performance Optimization

ระบบ production มีปริมาณ requests 100,000 requests/วินาที:
- จงอธิบายและเปรียบเทียบ head sampling vs tail sampling พร้อมข้อดี/ข้อเสีย
- ตั้งค่า OpenTelemetry SDK ให้ใช้ `parentbased_traceidratio` sampler ที่ sampling ratio = 0.01 (1%)
- เขียน Collector configuration ที่ใช้ tail sampling เพื่อ capture slow requests (duration > 1s) และ error requests (status >= 400) แม้ถูก head sampling ตัดทิ้ง
- ประเมิน overhead ที่เกิดขึ้นจากการทำ tail sampling เมื่อเทียบกับการไม่ทำ sampling


## แหล่งอ้างอิง

- [OpenTelemetry Official Documentation](https://opentelemetry.io/docs/)
- [OpenTelemetry Go SDK Documentation](https://opentelemetry.io/docs/instrumentation/go/)
- [OpenTelemetry Go GitHub Repository](https://github.com/open-telemetry/opentelemetry-go)
- [OpenTelemetry Go Contrib](https://github.com/open-telemetry/opentelemetry-go-contrib)
- [OpenTelemetry Collector Documentation](https://opentelemetry.io/docs/collector/)
- [OpenTelemetry Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/)
- [OpenTelemetry slog bridge (otelslog)](https://pkg.go.dev/go.opentelemetry.io/contrib/bridges/otelslog)[reference:26]
- [OpenTelemetry net/http instrumentation (otelhttp)](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp)[reference:27]
- [OpenTelemetry GORM plugin](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/gorm.io/gorm/otelgorm)[reference:28]
- [OpenTelemetry gRPC instrumentation (otelgrpc)](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc)[reference:29]
- [OpenTelemetry Metrics for Go](https://opentelemetry.io/docs/instrumentation/go/manual/#metrics)[reference:30]
- [OpenTelemetry Best Practices](https://opentelemetry.io/docs/specs/otel/overview/#best-practices)

---

**หมายเหตุ:** module นี้ครบถ้วนสำหรับ `pkg/opentelemetry` สำหรับระบบ gobackend หากต้องการ module เพิ่มเติม (เช่น `pkg/thanos`, `pkg/cortex`) โปรดแจ้ง