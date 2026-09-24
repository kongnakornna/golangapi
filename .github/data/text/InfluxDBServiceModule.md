# เอกสารฉบับสมบูรณ์: InfluxDB Service Module ด้วย DDD + Clean Architecture ในภาษา Go

---

## สารบัญ

1. [บทนำ](#1-บทนำ)
2. [สถาปัตยกรรมโดยรวม](#2-สถาปัตยกรรมโดยรวม)
3. [โครงสร้างโปรเจกต์](#3-โครงสร้างโปรเจกต์)
4. [เลเยอร์ Domain](#4-เลเยอร์-domain)
5. [เลเยอร์ Application (Use Cases)](#5-เลเยอร์-application-use-cases)
6. [เลเยอร์ Infrastructure (Repository)](#6-เลเยอร์-infrastructure-repository)
7. [เลเยอร์ Interface/Transport](#7-เลเยอร์-interfacetransport)
8. [Dependency Injection (Wire)](#8-dependency-injection-wire)
9. [การเขียนข้อมูล (Write)](#9-การเขียนข้อมูล-write)
10. [การสอบถามข้อมูล (Query)](#10-การสอบถามข้อมูล-query)
11. [การจัดการ Error](#11-การจัดการ-error)
12. [การทดสอบ (Testing)](#12-การทดสอบ-testing)
13. [แนวทางปฏิบัติที่ดีที่สุด](#13-แนวทางปฏิบัติที่ดีที่สุด)
14. [InfluxDB API Integration](#14-influxdb-api-integration)
15. [การจัดการ Connection และ Retry](#15-การจัดการ-connection-และ-retry)
16. [Worker Pool สำหรับ Batch Writing](#16-worker-pool-สำหรับการเขียนข้อมูลแบบ-batch)
17. [Graceful Shutdown และ Lifecycle Management](#17-graceful-shutdown-และ-lifecycle-management)
18. [Middleware และ Observability](#18-middleware-และ-observability)
19. [Deployment](#19-deployment)
20. [API Documentation (Swagger/OpenAPI)](#20-api-documentation-swaggeropenapi)
21. [สรุปและแนวทางเพิ่มเติม](#21-สรุปและแนวทางเพิ่มเติม)

---

## 1. บทนำ

เอกสารนี้เป็นแนวทางในการสร้าง **InfluxDB Service Module** สำหรับภาษา Go โดยใช้หลักการ **Domain-Driven Design (DDD)** ร่วมกับ **Clean Architecture** เพื่อให้ได้ระบบที่:

- **แยกความรับผิดชอบ** (Separation of Concerns) อย่างชัดเจน
- **สามารถทดสอบได้** (Testable)
- **บำรุงรักษาง่าย** (Maintainable)
- **เป็นอิสระจากเฟรมเวิร์กและฐานข้อมูล** (Framework & Database Independent)

### 1.1 หลักการสำคัญ

Clean Architecture หรือ Onion Architecture (Hexagonal/Ports & Adapters) มีหลักการสำคัญคือ **Dependency Rule**: การอ้างอิง dependencies ต้องชี้เข้าหาศูนย์กลาง (core business logic) เท่านั้น

```
外层 (Interface/Transport) → 中层 (Application) → 内层 (Domain)
Infrastructure → ขึ้นอยู่กับ interfaces ที่ Domain/Application กำหนด
```

Control flow จะไหลออกจากศูนย์กลาง: inbound adapters → application layer → domain layer → outbound adapters

### 1.2 InfluxDB API Overview

InfluxDB มี API ที่ครอบคลุมสำหรับการจัดการข้อมูลและระบบ โดยหลักๆ แล้วจะอยู่ภายใต้เส้นทาง `/api/v2/` พอร์ตเริ่มต้นคือ `8086` นอกจากนี้ยังมี API แบบย้อนหลังสำหรับ InfluxDB v1 อีกด้วย

| หมวดหมู่ | API Endpoint | คำอธิบาย |
|---------|--------------|-----------|
| **การเขียนข้อมูล** | `POST /api/v2/write` | เขียนข้อมูลลง InfluxDB |
| **การสอบถามข้อมูล** | `POST /api/v2/query` | สอบถามข้อมูลด้วยภาษา Flux |
| **การลบข้อมูล** | `POST /api/v2/delete` | ลบข้อมูลตามเงื่อนไข |
| **บัคเก็ต** | `GET/POST/PATCH/DELETE /api/v2/buckets` | จัดการบัคเก็ต |
| **การตรวจสอบสถานะ** | `GET /health`, `GET /ping` | ตรวจสอบสถานะระบบ |
| **ความเข้ากันได้ v1** | `GET /query`, `POST /write` | รองรับ InfluxQL และ v1 |

---

## 2. สถาปัตยกรรมโดยรวม

```
┌─────────────────────────────────────────────────────────────────┐
│                        cmd/                                     │
│                   (Entrypoint & Bootstrap)                      │
└───────────────────────────┬─────────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────────┐
│                    internal/di/                                  │
│               (Dependency Injection Wiring)                      │
└───────────────────────────┬─────────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────────┐
│                 internal/adapters/in/                           │
│          (Inbound Adapters: HTTP, gRPC, Jobs)                   │
└───────────────────────────┬─────────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────────┐
│              internal/core/application/                         │
│       (Use Cases / Command-Query Separation)                    │
└───────────────────────────┬─────────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────────┐
│               internal/core/domain/                             │
│           (Entities, Aggregates, Domain Services)               │
│                                                                 │
│    ┌─────────────────────────────────────────────────────┐      │
│    │              internal/core/ports/                   │      │
│    │      (Repository Interfaces - ports)                │      │
│    └─────────────────────────────────────────────────────┘      │
└───────────────────────────┬─────────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────────┐
│              internal/adapters/out/                             │
│     (Outbound Adapters: InfluxDB Repository Implementations)   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 3. โครงสร้างโปรเจกต์

```
project-root/
├── cmd/
│   └── api/
│       ├── main.go                 # Entrypoint
│       └── bootstrap/
│           └── bootstrap.go        # Application bootstrap
│
├── internal/
│   ├── core/                       # Core business logic
│   │   ├── domain/                 # Domain layer
│   │   │   ├── metric.go           # Entity
│   │   │   ├── metric_repository.go # Repository interface (port)
│   │   │   └── errors.go           # Domain errors
│   │   │
│   │   ├── application/            # Application layer (Use Cases)
│   │   │   ├── command/
│   │   │   │   ├── write_metric.go
│   │   │   │   └── delete_metric.go
│   │   │   ├── query/
│   │   │   │   ├── query_metric.go
│   │   │   │   └── query_aggregate.go
│   │   │   └── interfaces.go       # Application ports
│   │   │
│   │   └── ports/                  # Port interfaces
│   │       └── influxdb_port.go    # (Alternative: all ports here)
│   │
│   ├── adapters/
│   │   ├── in/                     # Inbound adapters
│   │   │   ├── http/
│   │   │   │   ├── handler.go
│   │   │   │   └── dto.go
│   │   │   └── jobs/
│   │   │       └── metric_writer.go
│   │   │
│   │   └── out/                    # Outbound adapters
│   │       └── influxdb/
│   │           ├── client.go       # InfluxDB client wrapper
│   │           ├── metric_repository.go  # Repository implementation
│   │           ├── converter.go    # Domain ↔ InfluxDB conversion
│   │           └── config.go       # Connection config
│   │
│   ├── di/                         # Dependency Injection
│   │   └── wire.go
│   │
│   └── pkg/                        # Shared utilities
│       ├── logger/
│       ├── errors/
│       └── context/
│
├── configs/
│   └── config.yaml                 # Configuration
│
├── go.mod
├── go.sum
├── Makefile
└── docker-compose.yml
```

---

## 4. เลเยอร์ Domain

Domain layer คือ **หัวใจของระบบ** ประกอบด้วย business entities, value objects, aggregates, และ repository interfaces

### 4.1 Entity: Metric

```go
// internal/core/domain/metric.go
package domain

import (
    "time"
)

// Metric represents a time-series data point in the domain
type Metric struct {
    Measurement string                 `json:"measurement"`
    Tags        map[string]string      `json:"tags"`
    Fields      map[string]interface{} `json:"fields"`
    Timestamp   time.Time              `json:"timestamp"`
}

// NewMetric creates a new Metric entity with validation
func NewMetric(measurement string, tags map[string]string, 
    fields map[string]interface{}, timestamp time.Time) (*Metric, error) {
    
    if measurement == "" {
        return nil, ErrInvalidMeasurement
    }
    if len(fields) == 0 {
        return nil, ErrNoFields
    }
    if timestamp.IsZero() {
        timestamp = time.Now()
    }
    
    return &Metric{
        Measurement: measurement,
        Tags:        tags,
        Fields:      fields,
        Timestamp:   timestamp,
    }, nil
}

// Validate checks if the metric is valid
func (m *Metric) Validate() error {
    if m.Measurement == "" {
        return ErrInvalidMeasurement
    }
    if len(m.Fields) == 0 {
        return ErrNoFields
    }
    return nil
}
```

### 4.2 Repository Interface (Port)

```go
// internal/core/domain/metric_repository.go
package domain

import "context"

// MetricRepository defines the port for metric data access
// This is the interface that the infrastructure layer must implement
type MetricRepository interface {
    // Write writes a single metric to InfluxDB
    Write(ctx context.Context, metric *Metric) error
    
    // WriteBatch writes multiple metrics to InfluxDB
    WriteBatch(ctx context.Context, metrics []*Metric) error
    
    // Query retrieves metrics based on filter criteria
    Query(ctx context.Context, filter *MetricFilter) ([]*Metric, error)
    
    // QueryAggregate retrieves aggregated metrics
    QueryAggregate(ctx context.Context, query *AggregateQuery) (*AggregateResult, error)
    
    // Delete deletes metrics matching the filter
    Delete(ctx context.Context, filter *DeleteFilter) error
    
    // Health checks the InfluxDB connection
    Health(ctx context.Context) error
}

// MetricFilter defines query filter criteria
type MetricFilter struct {
    Measurement string
    Start       time.Time
    End         time.Time
    Tags        map[string]string
    Limit       int
    Offset      int
}

// AggregateQuery defines aggregation query
type AggregateQuery struct {
    Measurement string
    Start       time.Time
    End         time.Time
    Tags        map[string]string
    Function    string   // e.g., "mean", "sum", "max", "min", "count"
    Field       string
    GroupBy     []string // e.g., ["tag", "1h"]
}

// AggregateResult contains aggregated data
type AggregateResult struct {
    Series []Series `json:"series"`
}

type Series struct {
    Tags   map[string]string      `json:"tags"`
    Values []SeriesValue          `json:"values"`
}

type SeriesValue struct {
    Time  time.Time   `json:"time"`
    Value interface{} `json:"value"`
}

// DeleteFilter defines deletion criteria
type DeleteFilter struct {
    Measurement string
    Start       time.Time
    End         time.Time
    Tags        map[string]string
}
```

### 4.3 Domain Errors

```go
// internal/core/domain/errors.go
package domain

import "errors"

var (
    ErrInvalidMeasurement = errors.New("invalid measurement name")
    ErrNoFields           = errors.New("no fields provided")
    ErrMetricNotFound     = errors.New("metric not found")
    ErrInvalidFilter      = errors.New("invalid filter criteria")
    ErrWriteFailed        = errors.New("failed to write metric")
    ErrQueryFailed        = errors.New("failed to query metrics")
    ErrDeleteFailed       = errors.New("failed to delete metrics")
    ErrConnectionFailed   = errors.New("failed to connect to InfluxDB")
)
```

---

## 5. เลเยอร์ Application (Use Cases)

Application layer จัดการ use cases และ orchestrate business logic โดยอาศัย interfaces ที่ domain layer กำหนด

### 5.1 Command: Write Metric

```go
// internal/core/application/command/write_metric.go
package command

import (
    "context"
    "time"
    
    "your-project/internal/core/domain"
)

// WriteMetricCommand represents the command to write a metric
type WriteMetricCommand struct {
    Measurement string
    Tags        map[string]string
    Fields      map[string]interface{}
    Timestamp   *time.Time // optional
}

// WriteMetricHandler handles the write metric command
type WriteMetricHandler struct {
    repo domain.MetricRepository
}

// NewWriteMetricHandler creates a new handler
func NewWriteMetricHandler(repo domain.MetricRepository) *WriteMetricHandler {
    return &WriteMetricHandler{repo: repo}
}

// Handle executes the write metric command
func (h *WriteMetricHandler) Handle(ctx context.Context, cmd WriteMetricCommand) error {
    // 1. Create domain entity
    timestamp := time.Now()
    if cmd.Timestamp != nil {
        timestamp = *cmd.Timestamp
    }
    
    metric, err := domain.NewMetric(cmd.Measurement, cmd.Tags, cmd.Fields, timestamp)
    if err != nil {
        return err
    }
    
    // 2. Validate
    if err := metric.Validate(); err != nil {
        return err
    }
    
    // 3. Write to repository
    return h.repo.Write(ctx, metric)
}
```

### 5.2 Command: Write Batch

```go
// internal/core/application/command/write_batch_metric.go
package command

type WriteBatchMetricCommand struct {
    Metrics []WriteMetricCommand
}

type WriteBatchMetricHandler struct {
    repo domain.MetricRepository
}

func NewWriteBatchMetricHandler(repo domain.MetricRepository) *WriteBatchMetricHandler {
    return &WriteBatchMetricHandler{repo: repo}
}

func (h *WriteBatchMetricHandler) Handle(ctx context.Context, cmd WriteBatchMetricCommand) error {
    metrics := make([]*domain.Metric, 0, len(cmd.Metrics))
    
    for _, m := range cmd.Metrics {
        timestamp := time.Now()
        if m.Timestamp != nil {
            timestamp = *m.Timestamp
        }
        
        metric, err := domain.NewMetric(m.Measurement, m.Tags, m.Fields, timestamp)
        if err != nil {
            return err
        }
        metrics = append(metrics, metric)
    }
    
    return h.repo.WriteBatch(ctx, metrics)
}
```

### 5.3 Query: Query Metric

```go
// internal/core/application/query/query_metric.go
package query

import (
    "context"
    "time"
    
    "your-project/internal/core/domain"
)

// QueryMetricQuery represents the query to retrieve metrics
type QueryMetricQuery struct {
    Measurement string
    Start       time.Time
    End         time.Time
    Tags        map[string]string
    Limit       int
    Offset      int
}

// QueryMetricResult contains the query result
type QueryMetricResult struct {
    Metrics []*domain.Metric `json:"metrics"`
    Total   int              `json:"total"`
}

// QueryMetricHandler handles the query metric query
type QueryMetricHandler struct {
    repo domain.MetricRepository
}

func NewQueryMetricHandler(repo domain.MetricRepository) *QueryMetricHandler {
    return &QueryMetricHandler{repo: repo}
}

func (h *QueryMetricHandler) Handle(ctx context.Context, q QueryMetricQuery) (*QueryMetricResult, error) {
    filter := &domain.MetricFilter{
        Measurement: q.Measurement,
        Start:       q.Start,
        End:         q.End,
        Tags:        q.Tags,
        Limit:       q.Limit,
        Offset:      q.Offset,
    }
    
    metrics, err := h.repo.Query(ctx, filter)
    if err != nil {
        return nil, err
    }
    
    return &QueryMetricResult{
        Metrics: metrics,
        Total:   len(metrics),
    }, nil
}
```

### 5.4 Query: Aggregate

```go
// internal/core/application/query/query_aggregate.go
package query

type QueryAggregateQuery struct {
    Measurement string
    Start       time.Time
    End         time.Time
    Tags        map[string]string
    Function    string
    Field       string
    GroupBy     []string
}

type QueryAggregateHandler struct {
    repo domain.MetricRepository
}

func NewQueryAggregateHandler(repo domain.MetricRepository) *QueryAggregateHandler {
    return &QueryAggregateHandler{repo: repo}
}

func (h *QueryAggregateHandler) Handle(ctx context.Context, q QueryAggregateQuery) (*domain.AggregateResult, error) {
    query := &domain.AggregateQuery{
        Measurement: q.Measurement,
        Start:       q.Start,
        End:         q.End,
        Tags:        q.Tags,
        Function:    q.Function,
        Field:       q.Field,
        GroupBy:     q.GroupBy,
    }
    
    return h.repo.QueryAggregate(ctx, query)
}
```

### 5.5 Application Services (Optional)

```go
// internal/core/application/interfaces.go
package application

import (
    "context"
    "time"
    
    "your-project/internal/core/domain"
)

// MetricService defines the application service interface
// This can be used as a facade for all metric operations
type MetricService interface {
    WriteMetric(ctx context.Context, measurement string, tags map[string]string, 
        fields map[string]interface{}, timestamp *time.Time) error
    WriteMetrics(ctx context.Context, metrics []MetricInput) error
    QueryMetrics(ctx context.Context, filter MetricFilterInput) ([]*domain.Metric, error)
    QueryAggregate(ctx context.Context, query AggregateQueryInput) (*domain.AggregateResult, error)
    DeleteMetrics(ctx context.Context, filter DeleteFilterInput) error
}

// DTOs for the service interface
type MetricInput struct {
    Measurement string
    Tags        map[string]string
    Fields      map[string]interface{}
    Timestamp   *time.Time
}

type MetricFilterInput struct {
    Measurement string
    Start       time.Time
    End         time.Time
    Tags        map[string]string
    Limit       int
    Offset      int
}

type AggregateQueryInput struct {
    Measurement string
    Start       time.Time
    End         time.Time
    Tags        map[string]string
    Function    string
    Field       string
    GroupBy     []string
}

type DeleteFilterInput struct {
    Measurement string
    Start       time.Time
    End         time.Time
    Tags        map[string]string
}
```

---

## 6. เลเยอร์ Infrastructure (Repository)

Infrastructure layer **implement interfaces** ที่ domain layer กำหนด

### 6.1 InfluxDB Client Configuration

```go
// internal/adapters/out/influxdb/config.go
package influxdb

import (
    "fmt"
    "time"
)

// Config holds InfluxDB connection configuration
type Config struct {
    URL          string        `yaml:"url"`
    Token        string        `yaml:"token"`
    Organization string        `yaml:"organization"`
    Bucket       string        `yaml:"bucket"`
    Timeout      time.Duration `yaml:"timeout"`
    MaxRetries   int           `yaml:"max_retries"`
    BatchSize    int           `yaml:"batch_size"`
}

// Validate validates the configuration
func (c *Config) Validate() error {
    if c.URL == "" {
        return fmt.Errorf("influxdb url is required")
    }
    if c.Token == "" {
        return fmt.Errorf("influxdb token is required")
    }
    if c.Organization == "" {
        return fmt.Errorf("influxdb organization is required")
    }
    if c.Bucket == "" {
        return fmt.Errorf("influxdb bucket is required")
    }
    if c.Timeout == 0 {
        c.Timeout = 30 * time.Second
    }
    if c.MaxRetries == 0 {
        c.MaxRetries = 3
    }
    if c.BatchSize == 0 {
        c.BatchSize = 1000
    }
    return nil
}
```

### 6.2 InfluxDB Client Wrapper

```go
// internal/adapters/out/influxdb/client.go
package influxdb

import (
    "context"
    "fmt"
    "crypto/tls"
    "net/http"
    "time"
    
    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"
)

// Client wraps the InfluxDB client
type Client struct {
    client     influxdb2.Client
    writeAPI   api.WriteAPIBlocking
    queryAPI   api.QueryAPI
    config     *Config
    httpClient *http.Client
}

// NewClient creates a new InfluxDB client
func NewClient(config *Config) (*Client, error) {
    if err := config.Validate(); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }
    
    // Create HTTP client with connection pool
    httpClient := &http.Client{
        Timeout: config.Timeout,
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
            IdleConnTimeout:     90 * time.Second,
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: false,
            },
            WriteBufferSize: 4096,
            ReadBufferSize:  4096,
        },
    }
    
    // Create InfluxDB client with custom HTTP client
    client := influxdb2.NewClientWithOptions(
        config.URL,
        config.Token,
        influxdb2.DefaultOptions().
            SetHTTPClient(httpClient).
            SetBatchSize(config.BatchSize).
            SetFlushInterval(1000), // 1 second
    )
    
    writeAPI := client.WriteAPIBlocking(config.Organization, config.Bucket)
    queryAPI := client.QueryAPI(config.Organization)
    
    return &Client{
        client:     client,
        writeAPI:   writeAPI,
        queryAPI:   queryAPI,
        config:     config,
        httpClient: httpClient,
    }, nil
}

// Close closes the client connection
func (c *Client) Close() {
    if c.client != nil {
        c.client.Close()
    }
}

// Health checks the connection
func (c *Client) Health(ctx context.Context) error {
    // Ping the server
    _, err := c.client.Health(ctx)
    return err
}

// Ready checks if the server is ready
func (c *Client) Ready(ctx context.Context) error {
    _, err := c.client.Ready(ctx)
    return err
}
```

### 6.3 Metric Repository Implementation

```go
// internal/adapters/out/influxdb/metric_repository.go
package influxdb

import (
    "context"
    "fmt"
    "time"
    
    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"
    "github.com/influxdata/influxdb-client-go/v2/api/query"
    
    "your-project/internal/core/domain"
    "your-project/internal/pkg/retry"
)

// MetricRepository implements domain.MetricRepository
type MetricRepository struct {
    client *Client
}

// NewMetricRepository creates a new metric repository
func NewMetricRepository(client *Client) *MetricRepository {
    return &MetricRepository{client: client}
}

// Write writes a single metric
func (r *MetricRepository) Write(ctx context.Context, metric *domain.Metric) error {
    point := influxdb2.NewPoint(
        metric.Measurement,
        metric.Tags,
        metric.Fields,
        metric.Timestamp,
    )
    
    // Use retry
    return retry.DoWithRetry(ctx, retry.DefaultConfig(), func() error {
        return r.client.writeAPI.WritePoint(ctx, point)
    })
}

// WriteBatch writes multiple metrics
func (r *MetricRepository) WriteBatch(ctx context.Context, metrics []*domain.Metric) error {
    points := make([]*influxdb2.Point, 0, len(metrics))
    
    for _, metric := range metrics {
        point := influxdb2.NewPoint(
            metric.Measurement,
            metric.Tags,
            metric.Fields,
            metric.Timestamp,
        )
        points = append(points, point)
    }
    
    return retry.DoWithRetry(ctx, retry.DefaultConfig(), func() error {
        return r.client.writeAPI.WritePoint(ctx, points...)
    })
}

// Query retrieves metrics based on filter
func (r *MetricRepository) Query(ctx context.Context, filter *domain.MetricFilter) ([]*domain.Metric, error) {
    // Build Flux query
    flux := r.buildQuery(filter)
    
    // Execute query
    result, err := r.client.queryAPI.Query(ctx, flux)
    if err != nil {
        return nil, fmt.Errorf("query failed: %w", err)
    }
    defer result.Close()
    
    // Parse results
    metrics := make([]*domain.Metric, 0)
    
    for result.Next() {
        record := result.Record()
        
        metric := &domain.Metric{
            Measurement: record.Measurement(),
            Tags:        record.Tags(),
            Timestamp:   record.Time(),
            Fields: map[string]interface{}{
                record.Field(): record.Value(),
            },
        }
        
        metrics = append(metrics, metric)
    }
    
    if result.Err() != nil {
        return nil, fmt.Errorf("query error: %w", result.Err())
    }
    
    return metrics, nil
}

// QueryAggregate retrieves aggregated metrics
func (r *MetricRepository) QueryAggregate(ctx context.Context, query *domain.AggregateQuery) (*domain.AggregateResult, error) {
    // Build Flux aggregation query
    flux := r.buildAggregateQuery(query)
    
    // Execute query
    result, err := r.client.queryAPI.Query(ctx, flux)
    if err != nil {
        return nil, fmt.Errorf("aggregate query failed: %w", err)
    }
    defer result.Close()
    
    // Parse results
    series := make([]domain.Series, 0)
    currentSeries := domain.Series{
        Tags:   make(map[string]string),
        Values: make([]domain.SeriesValue, 0),
    }
    
    for result.Next() {
        record := result.Record()
        
        currentSeries.Values = append(currentSeries.Values, domain.SeriesValue{
            Time:  record.Time(),
            Value: record.Value(),
        })
        
        for k, v := range record.Tags() {
            currentSeries.Tags[k] = v
        }
    }
    
    if result.Err() != nil {
        return nil, fmt.Errorf("aggregate query error: %w", result.Err())
    }
    
    if len(currentSeries.Values) > 0 {
        series = append(series, currentSeries)
    }
    
    return &domain.AggregateResult{
        Series: series,
    }, nil
}

// Delete deletes metrics matching the filter
func (r *MetricRepository) Delete(ctx context.Context, filter *domain.DeleteFilter) error {
    predicate := r.buildDeletePredicate(filter)
    
    err := r.client.client.DeleteAPI().Delete(
        ctx,
        filter.Start,
        filter.End,
        predicate,
        r.client.config.Organization,
        r.client.config.Bucket,
    )
    if err != nil {
        return fmt.Errorf("delete failed: %w", err)
    }
    
    return nil
}

// Health checks the connection
func (r *MetricRepository) Health(ctx context.Context) error {
    return r.client.Health(ctx)
}

// buildQuery builds a Flux query from a filter
func (r *MetricRepository) buildQuery(filter *domain.MetricFilter) string {
    flux := fmt.Sprintf(`from(bucket: "%s")
        |> range(start: %s, stop: %s)`,
        r.client.config.Bucket,
        filter.Start.Format(time.RFC3339),
        filter.End.Format(time.RFC3339),
    )
    
    if filter.Measurement != "" {
        flux += fmt.Sprintf(` |> filter(fn: (r) => r._measurement == "%s")`, 
            filter.Measurement)
    }
    
    for k, v := range filter.Tags {
        flux += fmt.Sprintf(` |> filter(fn: (r) => r.%s == "%s")`, k, v)
    }
    
    if filter.Limit > 0 {
        flux += fmt.Sprintf(` |> limit(n: %d)`, filter.Limit)
    }
    
    if filter.Offset > 0 {
        flux += fmt.Sprintf(` |> offset(n: %d)`, filter.Offset)
    }
    
    return flux
}

// buildAggregateQuery builds a Flux aggregation query
func (r *MetricRepository) buildAggregateQuery(query *domain.AggregateQuery) string {
    flux := fmt.Sprintf(`from(bucket: "%s")
        |> range(start: %s, stop: %s)
        |> filter(fn: (r) => r._measurement == "%s")
        |> filter(fn: (r) => r._field == "%s")`,
        r.client.config.Bucket,
        query.Start.Format(time.RFC3339),
        query.End.Format(time.RFC3339),
        query.Measurement,
        query.Field,
    )
    
    for k, v := range query.Tags {
        flux += fmt.Sprintf(` |> filter(fn: (r) => r.%s == "%s")`, k, v)
    }
    
    flux += fmt.Sprintf(` |> %s(column: "_value")`, query.Function)
    
    if len(query.GroupBy) > 0 {
        groupStr := ""
        for i, g := range query.GroupBy {
            if i > 0 {
                groupStr += ", "
            }
            groupStr += fmt.Sprintf(`"%s"`, g)
        }
        flux += fmt.Sprintf(` |> group(columns: [%s])`, groupStr)
    }
    
    return flux
}

// buildDeletePredicate builds a delete predicate
func (r *MetricRepository) buildDeletePredicate(filter *domain.DeleteFilter) string {
    predicates := make([]string, 0)
    
    if filter.Measurement != "" {
        predicates = append(predicates, fmt.Sprintf(`_measurement="%s"`, filter.Measurement))
    }
    
    for k, v := range filter.Tags {
        predicates = append(predicates, fmt.Sprintf(`%s="%s"`, k, v))
    }
    
    if len(predicates) == 0 {
        return ""
    }
    
    result := ""
    for i, p := range predicates {
        if i > 0 {
            result += " and "
        }
        result += p
    }
    
    return result
}
```

### 6.4 Domain-Infrastructure Converter

```go
// internal/adapters/out/influxdb/converter.go
package influxdb

import (
    "time"
    
    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
    
    "your-project/internal/core/domain"
)

// ToPoint converts a domain Metric to an InfluxDB Point
func ToPoint(metric *domain.Metric) *influxdb2.Point {
    return influxdb2.NewPoint(
        metric.Measurement,
        metric.Tags,
        metric.Fields,
        metric.Timestamp,
    )
}

// ToDomainMetric converts InfluxDB record to domain Metric
func ToDomainMetric(record *influxdb2.QueryRecord) *domain.Metric {
    return &domain.Metric{
        Measurement: record.Measurement(),
        Tags:        record.Tags(),
        Fields: map[string]interface{}{
            record.Field(): record.Value(),
        },
        Timestamp: record.Time(),
    }
}
```

---

## 7. เลเยอร์ Interface/Transport

HTTP layer รับ request, เรียก use cases, และส่ง response กลับ

### 7.1 HTTP Handler

```go
// internal/adapters/in/http/handler.go
package http

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"
    
    "your-project/internal/core/application/command"
    "your-project/internal/core/application/query"
)

// MetricHandler handles HTTP requests for metrics
type MetricHandler struct {
    writeHandler      *command.WriteMetricHandler
    writeBatchHandler *command.WriteBatchMetricHandler
    queryHandler      *query.QueryMetricHandler
    aggregateHandler  *query.QueryAggregateHandler
}

// NewMetricHandler creates a new metric handler
func NewMetricHandler(
    writeHandler *command.WriteMetricHandler,
    writeBatchHandler *command.WriteBatchMetricHandler,
    queryHandler *query.QueryMetricHandler,
    aggregateHandler *query.QueryAggregateHandler,
) *MetricHandler {
    return &MetricHandler{
        writeHandler:      writeHandler,
        writeBatchHandler: writeBatchHandler,
        queryHandler:      queryHandler,
        aggregateHandler:  aggregateHandler,
    }
}

// WriteMetric handles POST /api/v1/metrics
func (h *MetricHandler) WriteMetric(w http.ResponseWriter, r *http.Request) {
    var req WriteMetricRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    
    cmd := command.WriteMetricCommand{
        Measurement: req.Measurement,
        Tags:        req.Tags,
        Fields:      req.Fields,
        Timestamp:   req.Timestamp,
    }
    
    if err := h.writeHandler.Handle(r.Context(), cmd); err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    respondJSON(w, http.StatusCreated, map[string]string{"status": "success"})
}

// QueryMetrics handles GET /api/v1/metrics
func (h *MetricHandler) QueryMetrics(w http.ResponseWriter, r *http.Request) {
    measurement := r.URL.Query().Get("measurement")
    start, _ := time.Parse(time.RFC3339, r.URL.Query().Get("start"))
    end, _ := time.Parse(time.RFC3339, r.URL.Query().Get("end"))
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
    
    tags := parseTags(r.URL.Query().Get("tags"))
    
    q := query.QueryMetricQuery{
        Measurement: measurement,
        Start:       start,
        End:         end,
        Tags:        tags,
        Limit:       limit,
        Offset:      offset,
    }
    
    result, err := h.queryHandler.Handle(r.Context(), q)
    if err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    respondJSON(w, http.StatusOK, result)
}

// HealthCheck handles GET /health
func (h *MetricHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
    respondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

// ReadinessCheck handles GET /ready
func (h *MetricHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
    respondJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

// parseTags parses tags from query string
func parseTags(tagsStr string) map[string]string {
    tags := make(map[string]string)
    // Implementation depends on format
    return tags
}
```

### 7.2 DTOs

```go
// internal/adapters/in/http/dto.go
package http

import "time"

// WriteMetricRequest represents the request body for writing a metric
type WriteMetricRequest struct {
    Measurement string                 `json:"measurement"`
    Tags        map[string]string      `json:"tags"`
    Fields      map[string]interface{} `json:"fields"`
    Timestamp   *time.Time             `json:"timestamp,omitempty"`
}

// WriteMetricResponse represents the response for writing a metric
type WriteMetricResponse struct {
    Status  string `json:"status"`
    Message string `json:"message,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
    Error   string `json:"error"`
    Code    int    `json:"code"`
    Message string `json:"message"`
}

// respondJSON writes a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

// respondError writes an error response
func respondError(w http.ResponseWriter, status int, message string) {
    respondJSON(w, status, ErrorResponse{
        Error:   "error",
        Code:    status,
        Message: message,
    })
}
```

### 7.3 HTTP Server

```go
// internal/adapters/in/http/server.go
package http

import (
    "fmt"
    "net/http"
    "time"
)

type HTTPServerConfig struct {
    Port         int           `yaml:"port"`
    ReadTimeout  time.Duration `yaml:"read_timeout"`
    WriteTimeout time.Duration `yaml:"write_timeout"`
    IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type Server struct {
    *http.Server
    handler *MetricHandler
}

func NewServer(handler *MetricHandler, config HTTPServerConfig) *Server {
    mux := http.NewServeMux()
    
    // Register routes
    mux.HandleFunc("POST /api/v1/metrics", handler.WriteMetric)
    mux.HandleFunc("GET /api/v1/metrics", handler.QueryMetrics)
    mux.HandleFunc("GET /health", handler.HealthCheck)
    mux.HandleFunc("GET /ready", handler.ReadinessCheck)
    
    return &Server{
        Server: &http.Server{
            Addr:         fmt.Sprintf(":%d", config.Port),
            Handler:      mux,
            ReadTimeout:  config.ReadTimeout,
            WriteTimeout: config.WriteTimeout,
            IdleTimeout:  config.IdleTimeout,
        },
        handler: handler,
    }
}
```

---

## 8. Dependency Injection (Wire)

ใช้ Google Wire สำหรับ dependency injection

### 8.1 Wire Definition

```go
// internal/di/wire.go
//go:build wireinject
// +build wireinject

package di

import (
    "context"
    
    "github.com/google/wire"
    
    "your-project/internal/adapters/in/http"
    "your-project/internal/adapters/out/influxdb"
    "your-project/internal/core/application/command"
    "your-project/internal/core/application/query"
    "your-project/internal/core/domain"
)

// App contains all application dependencies
type App struct {
    MetricHandler *http.MetricHandler
    InfluxClient  *influxdb.Client
}

// InitializeApp wires up all dependencies
func InitializeApp(ctx context.Context, config *Config) (*App, error) {
    wire.Build(
        // Infrastructure
        influxdb.NewClient,
        influxdb.NewMetricRepository,
        
        // Domain - Repository interface binding
        wire.Bind(new(domain.MetricRepository), new(*influxdb.MetricRepository)),
        
        // Application - Commands
        command.NewWriteMetricHandler,
        command.NewWriteBatchMetricHandler,
        
        // Application - Queries
        query.NewQueryMetricHandler,
        query.NewQueryAggregateHandler,
        
        // HTTP Handlers
        http.NewMetricHandler,
        
        // App
        wire.Struct(new(App), "*"),
    )
    return &App{}, nil
}
```

### 8.2 Wire Generated Code

```go
// internal/di/wire_gen.go (auto-generated by wire)
// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//go:build !wireinject

package di

import (
    "context"
    "your-project/internal/adapters/in/http"
    "your-project/internal/adapters/out/influxdb"
    "your-project/internal/core/application/command"
    "your-project/internal/core/application/query"
)

func InitializeApp(ctx context.Context, config *Config) (*App, error) {
    influxdbClient, err := influxdb.NewClient(&config.InfluxDB)
    if err != nil {
        return nil, err
    }
    metricRepository := influxdb.NewMetricRepository(influxdbClient)
    writeMetricHandler := command.NewWriteMetricHandler(metricRepository)
    writeBatchMetricHandler := command.NewWriteBatchMetricHandler(metricRepository)
    queryMetricHandler := query.NewQueryMetricHandler(metricRepository)
    queryAggregateHandler := query.NewQueryAggregateHandler(metricRepository)
    metricHandler := http.NewMetricHandler(
        writeMetricHandler,
        writeBatchMetricHandler,
        queryMetricHandler,
        queryAggregateHandler,
    )
    return &App{
        MetricHandler: metricHandler,
        InfluxClient:  influxdbClient,
    }, nil
}
```

---

## 9. การเขียนข้อมูล (Write)

### 9.1 การใช้ WriteAPIBlocking (แบบ synchronous)

```go
// การใช้ WriteAPIBlocking
writeAPI := client.WriteAPIBlocking(org, bucket)

// สร้าง point
p := influxdb2.NewPoint(
    "stat",
    map[string]string{"unit": "temperature"},
    map[string]interface{}{"avg": 24.5, "max": 45},
    time.Now(),
)

// เขียนข้อมูล
err := writeAPI.WritePoint(context.Background(), p)
```

### 9.2 การเขียนแบบ Batch

```go
func (r *MetricRepository) WriteBatch(ctx context.Context, metrics []*domain.Metric) error {
    points := make([]*influxdb2.Point, 0, len(metrics))
    
    for _, metric := range metrics {
        point := influxdb2.NewPoint(
            metric.Measurement,
            metric.Tags,
            metric.Fields,
            metric.Timestamp,
        )
        points = append(points, point)
    }
    
    return r.client.writeAPI.WritePoint(ctx, points...)
}
```

### 9.3 API Reference

| API Endpoint | Method | คำอธิบาย |
|-------------|--------|-----------|
| `/api/v2/write` | POST | เขียนข้อมูลลง InfluxDB (รองรับ InfluxDB 1.8.0+) |

---

## 10. การสอบถามข้อมูล (Query)

### 10.1 การใช้ QueryAPI

```go
// สร้าง query client
queryAPI := client.QueryAPI(org)

// สร้าง Flux query
query := `from(bucket:"example-bucket")
    |> range(start: -1h)
    |> filter(fn: (r) => r._measurement == "stat")
    |> filter(fn: (r) => r._field == "avg")`

// Execute query
result, err := queryAPI.Query(ctx, query)
if err != nil {
    return err
}
defer result.Close()

// Parse results
for result.Next() {
    record := result.Record()
    fmt.Printf("Time: %v, Value: %v\n", record.Time(), record.Value())
}
```

### 10.2 Flux Query Examples

```go
// Query with time range
func buildTimeRangeQuery(bucket, start, stop string) string {
    return fmt.Sprintf(`from(bucket: "%s")
        |> range(start: %s, stop: %s)`, bucket, start, stop)
}

// Query with filters
func buildFilteredQuery(bucket, start, measurement string, tags map[string]string) string {
    query := fmt.Sprintf(`from(bucket: "%s")
        |> range(start: %s)
        |> filter(fn: (r) => r._measurement == "%s")`, bucket, start, measurement)
    
    for k, v := range tags {
        query += fmt.Sprintf(` |> filter(fn: (r) => r.%s == "%s")`, k, v)
    }
    
    return query
}

// Aggregation query
func buildAggregationQuery(bucket, start, measurement, field, function string) string {
    return fmt.Sprintf(`from(bucket: "%s")
        |> range(start: %s)
        |> filter(fn: (r) => r._measurement == "%s")
        |> filter(fn: (r) => r._field == "%s")
        |> %s(column: "_value")`, bucket, start, measurement, field, function)
}
```

### 10.3 API Reference

| API Endpoint | Method | คำอธิบาย |
|-------------|--------|-----------|
| `/api/v2/query` | POST | สอบถามข้อมูลด้วยภาษา Flux |
| `/api/v2/query/suggestions` | GET | ดูคำแนะนำสำหรับการสอบถาม |
| `/query` | GET | สอบถามข้อมูลแบบ InfluxQL (v1 compatibility) |

---

## 11. การจัดการ Error

### 11.1 Domain Errors

```go
// internal/core/domain/errors.go
package domain

import "errors"

var (
    ErrInvalidMeasurement = errors.New("invalid measurement name")
    ErrNoFields           = errors.New("no fields provided")
    ErrMetricNotFound     = errors.New("metric not found")
    ErrInvalidFilter      = errors.New("invalid filter criteria")
    ErrWriteFailed        = errors.New("failed to write metric")
    ErrQueryFailed        = errors.New("failed to query metrics")
    ErrDeleteFailed       = errors.New("failed to delete metrics")
    ErrConnectionFailed   = errors.New("failed to connect to InfluxDB")
)
```

### 11.2 Error Wrapping

```go
func (r *MetricRepository) Write(ctx context.Context, metric *domain.Metric) error {
    point := influxdb2.NewPoint(
        metric.Measurement,
        metric.Tags,
        metric.Fields,
        metric.Timestamp,
    )
    
    if err := r.client.writeAPI.WritePoint(ctx, point); err != nil {
        return fmt.Errorf("%w: %v", domain.ErrWriteFailed, err)
    }
    
    return nil
}
```

---

## 12. การทดสอบ (Testing)

### 12.1 Mock Repository

```go
// internal/core/domain/mocks/metric_repository_mock.go
//go:generate mockgen -source=metric_repository.go -destination=mocks/metric_repository_mock.go -package=mocks

package mocks

import (
    "context"
    
    "github.com/golang/mock/gomock"
    
    "your-project/internal/core/domain"
)

// MockMetricRepository is a mock of MetricRepository interface
type MockMetricRepository struct {
    ctrl     *gomock.Controller
    recorder *MockMetricRepositoryMockRecorder
}

// ... (generated by mockgen)
```

### 12.2 Unit Test for Use Case

```go
// internal/core/application/command/write_metric_test.go
package command_test

import (
    "context"
    "testing"
    
    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"
    
    "your-project/internal/core/application/command"
    "your-project/internal/core/domain"
    "your-project/internal/core/domain/mocks"
)

func TestWriteMetricHandler_Handle_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockMetricRepository(ctrl)
    
    mockRepo.EXPECT().
        Write(gomock.Any(), gomock.Any()).
        Return(nil)
    
    handler := command.NewWriteMetricHandler(mockRepo)
    
    cmd := command.WriteMetricCommand{
        Measurement: "test",
        Tags:        map[string]string{"env": "prod"},
        Fields:      map[string]interface{}{"value": 100},
    }
    
    err := handler.Handle(context.Background(), cmd)
    
    assert.NoError(t, err)
}

func TestWriteMetricHandler_Handle_InvalidMeasurement(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockMetricRepository(ctrl)
    
    handler := command.NewWriteMetricHandler(mockRepo)
    
    cmd := command.WriteMetricCommand{
        Measurement: "", // invalid
        Tags:        map[string]string{"env": "prod"},
        Fields:      map[string]interface{}{"value": 100},
    }
    
    err := handler.Handle(context.Background(), cmd)
    
    assert.Error(t, err)
    assert.Equal(t, domain.ErrInvalidMeasurement, err)
}
```

### 12.3 Integration Test with Testcontainers

```go
// internal/adapters/out/influxdb/metric_repository_test.go
package influxdb_test

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
    
    "your-project/internal/adapters/out/influxdb"
    "your-project/internal/core/domain"
)

func TestMetricRepository_WriteAndQuery(t *testing.T) {
    ctx := context.Background()
    
    req := testcontainers.ContainerRequest{
        Image:        "influxdb:2.7",
        ExposedPorts: []string{"8086/tcp"},
        WaitingFor:   wait.ForLog("Listening on HTTP"),
        Env: map[string]string{
            "INFLUXDB_DB":            "test",
            "INFLUXDB_ADMIN_USER":    "admin",
            "INFLUXDB_ADMIN_PASSWORD": "password",
        },
    }
    
    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    require.NoError(t, err)
    defer container.Terminate(ctx)
    
    host, _ := container.Host(ctx)
    port, _ := container.MappedPort(ctx, "8086")
    
    config := &influxdb.Config{
        URL:          "http://" + host + ":" + port.Port(),
        Token:        "my-token",
        Organization: "my-org",
        Bucket:       "test",
    }
    
    client, err := influxdb.NewClient(config)
    require.NoError(t, err)
    defer client.Close()
    
    repo := influxdb.NewMetricRepository(client)
    
    metric, err := domain.NewMetric(
        "test_measurement",
        map[string]string{"env": "test"},
        map[string]interface{}{"value": 42},
        time.Now(),
    )
    require.NoError(t, err)
    
    err = repo.Write(ctx, metric)
    assert.NoError(t, err)
    
    filter := &domain.MetricFilter{
        Measurement: "test_measurement",
        Start:       time.Now().Add(-1 * time.Hour),
        End:         time.Now().Add(1 * time.Hour),
        Tags:        map[string]string{"env": "test"},
    }
    
    metrics, err := repo.Query(ctx, filter)
    assert.NoError(t, err)
    assert.NotEmpty(t, metrics)
}
```

---

## 13. แนวทางปฏิบัติที่ดีที่สุด

### 13.1 Dependency Rule

**หลักการ**: Dependencies ต้องชี้เข้าหาศูนย์กลางเท่านั้น

```
✅ ถูกต้อง: transport → application → domain
✅ ถูกต้อง: infrastructure → domain (implement interface)
❌ ผิด: domain → infrastructure
❌ ผิด: application → transport
```

### 13.2 Interface Segregation

กำหนด interface ให้เล็กและเฉพาะเจาะจง:

```go
// ✅ ดี: แยก interface ตามความรับผิดชอบ
type MetricWriter interface {
    Write(ctx context.Context, metric *Metric) error
    WriteBatch(ctx context.Context, metrics []*Metric) error
}

type MetricReader interface {
    Query(ctx context.Context, filter *MetricFilter) ([]*Metric, error)
    QueryAggregate(ctx context.Context, query *AggregateQuery) (*AggregateResult, error)
}

type MetricDeleter interface {
    Delete(ctx context.Context, filter *DeleteFilter) error
}

// ❌ ไม่ดี: interface ใหญ่เกินไป
type MetricRepository interface {
    Write(ctx context.Context, metric *Metric) error
    WriteBatch(ctx context.Context, metrics []*Metric) error
    Query(ctx context.Context, filter *MetricFilter) ([]*Metric, error)
    QueryAggregate(ctx context.Context, query *AggregateQuery) (*AggregateResult, error)
    Delete(ctx context.Context, filter *DeleteFilter) error
    Health(ctx context.Context) error
}
```

### 13.3 Context Management

```go
// ✅ ดี: ส่ง context ไปทุก layer
func (r *MetricRepository) Write(ctx context.Context, metric *domain.Metric) error {
    return r.client.writeAPI.WritePoint(ctx, point)
}

// ❌ ไม่ดี: ไม่ใช้ context หรือใช้ context.Background()
func (r *MetricRepository) Write(metric *domain.Metric) error {
    return r.client.writeAPI.WritePoint(context.Background(), point)
}
```

### 13.4 Error Handling

```go
// ✅ ดี: wrap error ด้วย context
if err != nil {
    return fmt.Errorf("failed to write metric: %w", err)
}

// ✅ ดี: ใช้ domain errors
if metric.Measurement == "" {
    return domain.ErrInvalidMeasurement
}

// ❌ ไม่ดี: return raw error
if err != nil {
    return err
}
```

---

## 14. InfluxDB API Integration

### 14.1 API Endpoints ที่รองรับ

| หมวดหมู่ | API Endpoint | Method | คำอธิบาย | Implementation |
|---------|--------------|--------|-----------|----------------|
| **การเขียนข้อมูล** | `/api/v2/write` | POST | เขียนข้อมูลลง InfluxDB | `MetricRepository.Write()` |
| **การสอบถามข้อมูล** | `/api/v2/query` | POST | สอบถามด้วย Flux | `MetricRepository.Query()` |
| **การลบข้อมูล** | `/api/v2/delete` | POST | ลบข้อมูลตามเงื่อนไข | `MetricRepository.Delete()` |
| **บัคเก็ต** | `/api/v2/buckets` | GET/POST | จัดการบัคเก็ต | `BucketRepository` |
| **สถานะระบบ** | `/health`, `/ping` | GET | ตรวจสอบสถานะ | `Client.Health()` |

### 14.2 การใช้งาน API ผ่าน Repository

```go
// การเขียนข้อมูล
func (r *MetricRepository) Write(ctx context.Context, metric *domain.Metric) error {
    point := influxdb2.NewPoint(metric.Measurement, metric.Tags, metric.Fields, metric.Timestamp)
    return r.client.writeAPI.WritePoint(ctx, point)
}

// การสอบถามข้อมูล
func (r *MetricRepository) Query(ctx context.Context, filter *domain.MetricFilter) ([]*domain.Metric, error) {
    flux := r.buildQuery(filter)
    result, err := r.client.queryAPI.Query(ctx, flux)
    // ... parse results
}

// การลบข้อมูล
func (r *MetricRepository) Delete(ctx context.Context, filter *domain.DeleteFilter) error {
    predicate := r.buildDeletePredicate(filter)
    return r.client.client.DeleteAPI().Delete(ctx, filter.Start, filter.End, 
        predicate, r.client.config.Organization, r.client.config.Bucket)
}
```

---

## 15. การจัดการ Connection และ Retry

### 15.1 Retry with Exponential Backoff

```go
// internal/pkg/retry/retry.go
package retry

import (
    "context"
    "fmt"
    "time"
)

type Config struct {
    MaxRetries      int
    InitialInterval time.Duration
    MaxInterval     time.Duration
    Multiplier      float64
}

func DefaultConfig() Config {
    return Config{
        MaxRetries:      3,
        InitialInterval: 100 * time.Millisecond,
        MaxInterval:     5 * time.Second,
        Multiplier:      2.0,
    }
}

func DoWithRetry(ctx context.Context, cfg Config, fn func() error) error {
    var err error
    interval := cfg.InitialInterval
    
    for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
        if attempt > 0 {
            select {
            case <-ctx.Done():
                return ctx.Err()
            case <-time.After(interval):
            }
            interval = time.Duration(float64(interval) * cfg.Multiplier)
            if interval > cfg.MaxInterval {
                interval = cfg.MaxInterval
            }
        }
        
        err = fn()
        if err == nil {
            return nil
        }
        
        if !isRetryable(err) {
            return err
        }
    }
    
    return fmt.Errorf("max retries exceeded: %w", err)
}

func isRetryable(err error) bool {
    // Connection errors, timeouts, 5xx
    return true // placeholder
}
```

### 15.2 นำไปใช้ใน Repository

```go
func (r *MetricRepository) Write(ctx context.Context, metric *domain.Metric) error {
    point := influxdb2.NewPoint(metric.Measurement, metric.Tags, metric.Fields, metric.Timestamp)
    
    return retry.DoWithRetry(ctx, retry.DefaultConfig(), func() error {
        return r.client.writeAPI.WritePoint(ctx, point)
    })
}
```

---

## 16. Worker Pool สำหรับการเขียนข้อมูลแบบ Batch

```go
// internal/adapters/out/influxdb/worker_pool.go
package influxdb

import (
    "context"
    "sync"
    "time"
    
    "your-project/internal/core/domain"
)

type BatchWriter struct {
    repo          domain.MetricRepository
    workers       []*worker
    wg            sync.WaitGroup
    stopChan      chan struct{}
    batchSize     int
    flushInterval time.Duration
}

type worker struct {
    id         int
    workChan   chan *domain.Metric
    batch      []*domain.Metric
    repo       domain.MetricRepository
    batchSize  int
    mu         sync.Mutex
}

func NewBatchWriter(repo domain.MetricRepository, numWorkers int, batchSize int, flushInterval time.Duration) *BatchWriter {
    bw := &BatchWriter{
        repo:          repo,
        workers:       make([]*worker, numWorkers),
        stopChan:      make(chan struct{}),
        batchSize:     batchSize,
        flushInterval: flushInterval,
    }
    
    for i := 0; i < numWorkers; i++ {
        w := &worker{
            id:        i,
            workChan:  make(chan *domain.Metric, 100),
            batch:     make([]*domain.Metric, 0, batchSize),
            repo:      repo,
            batchSize: batchSize,
        }
        bw.workers[i] = w
        bw.wg.Add(1)
        go w.run(bw.stopChan)
    }
    
    return bw
}

func (w *worker) run(stopChan <-chan struct{}) {
    defer w.wg.Done()
    ticker := time.NewTicker(flushInterval)
    defer ticker.Stop()
    
    for {
        select {
        case <-stopChan:
            w.flush()
            return
        case metric := <-w.workChan:
            w.addMetric(metric)
        case <-ticker.C:
            w.flush()
        }
    }
}

func (w *worker) addMetric(metric *domain.Metric) {
    w.mu.Lock()
    defer w.mu.Unlock()
    w.batch = append(w.batch, metric)
    if len(w.batch) >= w.batchSize {
        w.flushLocked()
    }
}

func (w *worker) flush() {
    w.mu.Lock()
    defer w.mu.Unlock()
    w.flushLocked()
}

func (w *worker) flushLocked() {
    if len(w.batch) == 0 {
        return
    }
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    _ = w.repo.WriteBatch(ctx, w.batch)
    w.batch = w.batch[:0]
}

func (bw *BatchWriter) Write(metric *domain.Metric) {
    // Round-robin load balancing
    workerIdx := int(time.Now().UnixNano()) % len(bw.workers)
    bw.workers[workerIdx].workChan <- metric
}

func (bw *BatchWriter) Close() error {
    close(bw.stopChan)
    bw.wg.Wait()
    return nil
}
```

---

## 17. Graceful Shutdown และ Lifecycle Management

### 17.1 Application Bootstrap

```go
// cmd/api/bootstrap/bootstrap.go
package bootstrap

import (
    "context"
    "fmt"
    "log/slog"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "your-project/internal/adapters/in/http"
    "your-project/internal/adapters/out/influxdb"
    "your-project/internal/di"
)

func Run(configPath string) error {
    config, err := LoadConfig(configPath)
    if err != nil {
        return fmt.Errorf("failed to load config: %w", err)
    }
    
    ctx := context.Background()
    app, err := di.InitializeApp(ctx, config)
    if err != nil {
        return fmt.Errorf("failed to initialize app: %w", err)
    }
    
    server := http.NewServer(app.MetricHandler, config.HTTP)
    
    return runWithGracefulShutdown(server, app.InfluxClient)
}

func runWithGracefulShutdown(server *http.Server, influxClient *influxdb.Client) error {
    stopChan := make(chan os.Signal, 1)
    signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
    
    errChan := make(chan error, 1)
    go func() {
        slog.Info("starting HTTP server", "addr", server.Addr)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            errChan <- err
        }
    }()
    
    select {
    case err := <-errChan:
        return fmt.Errorf("server error: %w", err)
    case <-stopChan:
        slog.Info("received shutdown signal, gracefully stopping...")
        
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        if err := server.Shutdown(ctx); err != nil {
            slog.Error("server shutdown error", "error", err)
        }
        
        influxClient.Close()
        slog.Info("influxdb client closed")
        return nil
    }
}
```

---

## 18. Middleware และ Observability

### 18.1 Logging Middleware

```go
// internal/adapters/in/http/middleware/logging.go
package middleware

import (
    "log/slog"
    "net/http"
    "time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        next.ServeHTTP(wrapped, r)
        
        slog.Info("HTTP request",
            "method", r.Method,
            "path", r.URL.Path,
            "status", wrapped.statusCode,
            "duration", time.Since(start),
            "remote_addr", r.RemoteAddr,
        )
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

### 18.2 Recovery Middleware

```go
// internal/adapters/in/http/middleware/recovery.go
package middleware

import (
    "log/slog"
    "net/http"
    "runtime/debug"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                slog.Error("panic recovered", "error", err, "stack", string(debug.Stack()))
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

### 18.3 Metrics Middleware (Prometheus)

```go
// internal/adapters/in/http/middleware/metrics.go
package middleware

import (
    "net/http"
    "strconv"
    "time"
    
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )
)

func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        next.ServeHTTP(wrapped, r)
        
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(wrapped.statusCode)
        
        httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
        httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
    })
}
```

---

## 19. Deployment

### 19.1 Dockerfile

```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/api .
COPY --from=builder /app/configs ./configs

EXPOSE 8080

CMD ["./api", "-config", "configs/config.yaml"]
```

### 19.2 docker-compose.yml

```yaml
version: '3.8'

services:
  influxdb:
    image: influxdb:2.7
    container_name: influxdb
    ports:
      - "8086:8086"
    environment:
      - DOCKER_INFLUXDB_INIT_MODE=setup
      - DOCKER_INFLUXDB_INIT_USERNAME=admin
      - DOCKER_INFLUXDB_INIT_PASSWORD=password
      - DOCKER_INFLUXDB_INIT_ORG=myorg
      - DOCKER_INFLUXDB_INIT_BUCKET=mybucket
      - DOCKER_INFLUXDB_INIT_RETENTION=1w
      - DOCKER_INFLUXDB_INIT_ADMIN_TOKEN=my-super-secret-token
    volumes:
      - influxdb-data:/var/lib/influxdb2
    healthcheck:
      test: ["CMD", "influx", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  app:
    build: .
    container_name: influxdb-service
    ports:
      - "8080:8080"
    depends_on:
      influxdb:
        condition: service_healthy
    environment:
      - INFLUXDB_URL=http://influxdb:8086
      - INFLUXDB_TOKEN=my-super-secret-token
      - INFLUXDB_ORG=myorg
      - INFLUXDB_BUCKET=mybucket

  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    depends_on:
      - prometheus
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin

volumes:
  influxdb-data:
```

---

## 20. API Documentation (Swagger/OpenAPI)

### 20.1 ติดตั้งและตั้งค่า

```bash
go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/http-swagger
```

### 20.2 เพิ่ม Annotation ใน Handler

```go
// internal/adapters/in/http/handler.go

// WriteMetric godoc
// @Summary Write a new metric
// @Description Write a single metric point to InfluxDB
// @Tags metrics
// @Accept json
// @Produce json
// @Param request body WriteMetricRequest true "Metric data"
// @Success 201 {object} WriteMetricResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/metrics [post]
func (h *MetricHandler) WriteMetric(w http.ResponseWriter, r *http.Request) {
    // ... implementation
}

// QueryMetrics godoc
// @Summary Query metrics
// @Description Retrieve metrics based on filter criteria
// @Tags metrics
// @Accept json
// @Produce json
// @Param measurement query string false "Measurement name"
// @Param start query string true "Start time (RFC3339)"
// @Param end query string true "End time (RFC3339)"
// @Param tags query string false "Tags (format: key1=value1,key2=value2)"
// @Param limit query int false "Limit" default(100)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} QueryMetricResult
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/metrics [get]
func (h *MetricHandler) QueryMetrics(w http.ResponseWriter, r *http.Request) {
    // ... implementation
}
```

### 20.3 เปิด Swagger UI

```go
// internal/adapters/in/http/server.go (เพิ่มเติม)
import (
    _ "your-project/docs"
    httpSwagger "github.com/swaggo/http-swagger"
)

func NewServer(handler *MetricHandler, config HTTPServerConfig) *Server {
    mux := http.NewServeMux()
    
    // ... other routes
    
    mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
    
    return &Server{...}
}
```

### 20.4 สร้างเอกสาร

```bash
swag init -g cmd/api/main.go -o docs
```

---

## 21. สรุปและแนวทางเพิ่มเติม

เอกสารนี้ได้นำเสนอการสร้าง **InfluxDB Service Module** อย่างสมบูรณ์แบบด้วย **Go + DDD + Clean Architecture** ครอบคลุมตั้งแต่:

- Domain Entities และ Repository Interfaces
- Application Use Cases (Command/Query)
- Infrastructure Implementation (InfluxDB)
- HTTP Transport และ Middleware
- Dependency Injection (Wire)
- Error Handling, Retry, Batch Writing
- Graceful Shutdown และ Deployment

### 21.1 แนวทางเพิ่มเติม

1. **Event Sourcing** - ใช้ domain events เพื่อติดตามการเปลี่ยนแปลง
2. **CQRS** - แยก command และ query stores
3. **API Versioning** - รองรับหลายเวอร์ชันของ API
4. **Rate Limiting** - ป้องกันการใช้งานเกิน
5. **Circuit Breaker** - เพิ่มความทนทานต่อความล้มเหลว
6. **Feature Flags** - ควบคุมการเปิดใช้งานฟีเจอร์แบบไดนามิก
7. **Distributed Tracing** - ติดตาม request ข้ามบริการ
8. **Profiling** - ใช้ pprof เพื่อวิเคราะห์ประสิทธิภาพ

### 21.2 เครื่องมือแนะนำ

| เครื่องมือ | การใช้งาน |
|-----------|-----------|
| **Wire** | Dependency Injection |
| **Mockgen** | สร้าง mocks สำหรับ testing |
| **Testcontainers** | Integration testing |
| **Prometheus + Grafana** | Monitoring |
| **Jaeger** | Distributed Tracing |
| **Swaggo** | API Documentation |
| **golangci-lint** | Code quality |

### 21.3 การตรวจสอบความสมบูรณ์ก่อน Deploy

- ✅ ทุกเลเยอร์เป็นไปตาม Dependency Rule
- ✅ Repository interface ถูก implement ครบถ้วน
- ✅ มีการจัดการ context และ timeout
- ✅ มี graceful shutdown
- ✅ มี health check และ readiness probe
- ✅ มี logging และ metrics
- ✅ มีการทดสอบ unit และ integration
- ✅ มีการจัดการ error ที่เหมาะสม
- ✅ มี API documentation (Swagger)

---

**จบเอกสารฉบับสมบูรณ์**