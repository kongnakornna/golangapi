# Module 18: pkg/timescaledb

## สำหรับโฟลเดอร์ `pkg/timescaledb/`

ไฟล์ที่เกี่ยวข้อง:
- `client.go` - การเชื่อมต่อกับ TimescaleDB (PostgreSQL)
- `hypertable.go` - สร้างและจัดการ hypertable
- `writer.go` - การเขียนข้อมูล time-series
- `query.go` - การ query พร้อม time-bucket functions
- `compression.go` - ตั้งค่า compression และ retention

---

## หลักการ (Concept)

### TimescaleDB คืออะไร?
TimescaleDB คือส่วนขยาย (extension) ของ PostgreSQL ที่ออกแบบมาเพื่อจัดการข้อมูล Time Series โดยเฉพาะ โดยทำงานบน PostgreSQL เต็มรูปแบบ ทำให้สามารถใช้ SQL ปกติ + ฟังก์ชัน time-series เพิ่มเติม เช่น `time_bucket()`, `first()`, `last()` และการจัดการ lifecycle อัตโนมัติ

### มีกี่แบบ? (Formats / Deployment)

**รูปแบบการติดตั้ง:**
1. **TimescaleDB Community (Open Source)** - ฟรี, รองรับ single node
2. **TimescaleDB Enterprise** - รองรับ clustering, replication, backup ขั้นสูง
3. **Timescale Cloud (Managed Service)** - บริการ fully-managed บน cloud

**โครงสร้างข้อมูลหลัก:**
- **Hypertable** - ตารางหลักที่ถูกแบ่งเป็น chunks อัตโนมัติตามเวลา (partition by time) และตัวเลือก partition โดย space (เช่น device_id)
- **Chunk** - ส่วนย่อยของ hypertable แต่ละ chunk เป็นตาราง PostgreSQL ปกติ ช่วยในการลบข้อมูลเก่า และ compression
- **Continuous Aggregate** - materialized view ที่ถูก refresh อัตโนมัติ ใช้สำหรับ pre-aggregate ข้อมูล (คล้าย downsampling)
- **Compression Policy** - การบีบอัดข้อมูลแบบ segment-by-column อัตโนมัติ ตามนโยบายเวลา
- **Retention Policy** - การลบ chunks ที่เก่ากว่าเวลาที่กำหนด

### ใช้อย่างไร / นำไปใช้กรณีไหน

**กรณีใช้งาน:**
- IoT / sensor data ที่ต้องการ SQL เต็มรูปแบบ (JOIN, subquery, window functions)
- ระบบการเงินที่ต้องการความถูกต้องแม่นยำและ ACID
- แทนที่ PostgreSQL เดิมที่มีข้อมูล time-series จำนวนมาก
- ต้องการใช้เครื่องมือ BI/analytics ที่เชื่อมต่อ PostgreSQL (Tableau, Metabase, Superset)
- ข้อมูลที่ต้องมีการ update หรือ delete บางจุด (TimescaleDB รองรับการ update ได้ดีกว่า InfluxDB)

**รูปแบบการเขียนข้อมูล:**
```sql
-- Hypertable เหมือนตารางปกติ
CREATE TABLE metrics (
    time        TIMESTAMPTZ NOT NULL,
    device_id   TEXT NOT NULL,
    cpu_usage   DOUBLE PRECISION,
    temperature FLOAT
);
SELECT create_hypertable('metrics', 'time', chunk_time_interval => INTERVAL '1 day');
```

### ประโยชน์ที่ได้รับ
- ใช้ SQL เต็มรูปแบบ (JOIN, CTE, window functions) กับข้อมูล time-series
- รองรับการ update และ delete แบบระบุแถวได้ (ไม่เหมือน InfluxDB)
- ACID compliance - เหมาะกับงานการเงินและระบบที่ต้องการความถูกต้อง
- ชุดฟังก์ชัน time-series ที่มีประสิทธิภาพ (`time_bucket()`, `locf()`, `interpolate()`)
- Compression อัตราส่วนสูง (90-98% ประหยัดพื้นที่)
- ทำงานร่วมกับ PostgreSQL ecosystem ทั้งหมด (extensions, tools)

### ข้อควรระวัง
- **การเลือก chunk_time_interval** - ควรทำให้แต่ละ chunk มีขนาดประมาณ 25% ของ RAM (เช่น ข้อมูล 1 วัน ถ้าขนาด 1GB/chunk ควรมี RAM 4GB+)
- **Space partitioning** - ถ้าใช้ partition key เพิ่ม (เช่น device_id) cardinality ของค่านั้นควรมีจำนวนจำกัด (<10,000) มิฉะนั้นจะสร้าง chunks มากเกิน
- **Continuous aggregate** - การ refresh อาจใช้ทรัพยากรสูง ถ้าช่วง refresh ถี่เกินไป
- **Index** - ควรมี index ที่ (time, device_id) เสมอ การ query ที่ไม่ระบุ time range จะช้ามาก

### ข้อดี
- SQL ที่ทุกคนรู้จัก ลดการเรียนรู้
- JOIN ระหว่าง hypertable และตารางปกติทำได้ง่าย
- รองรับ window functions และ advanced analytics
- ข้อมูลถูกเก็บในรูปแบบ row-based ทำให้ update/delete ทำได้
- มี continuous aggregate สำหรับ downsampling อัตโนมัติ

### ข้อเสีย
- การเขียนข้อมูล (insert) ช้ากว่า InfluxDB โดยเฉพาะที่ปริมาณสูง (>100k event/วินาที ต่อ node)
- ต้องการ tuning PostgreSQL (shared_buffers, effective_cache_size, work_mem)
- การ query แบบ time-series เฉพาะทาง (เช่น moving average, fill missing) ต้องเขียน SQL ที่ซับซ้อนกว่า Flux
- ต้องการพื้นที่จัดเก็บมากกว่า InfluxDB ก่อน compression (แต่หลัง compression ใกล้เคียงกัน)

### ข้อห้าม
**ห้ามใช้ TimescaleDB แทน relational database ปกติ** สำหรับข้อมูลที่ไม่ใช่ time-series (เช่น ตาราง user, product) เพราะ hypertable ออกแบบมาให้ insert แบบ append-only เป็นหลัก การ update/delete บ่อยครั้งจะทำลายประสิทธิภาพของ chunk partitioning ให้แยกตาราง time-series ออกจากตาราง master data แล้ว JOIN ทีหลัง[reference: TimescaleDB best practices]

---

## การออกแบบ Workflow และ Dataflow

```mermaid
flowchart LR
    A[IoT Devices] --> B[API Gateway]
    B --> C[Go App]
    C -->|INSERT| D[(TimescaleDB)]
    D --> E[Hypertable + Chunks]
    E --> F[Compression Policy]
    E --> G[Retention Policy]
    E --> H[Continuous Aggregate]
    H --> I[Query via SQL]
    I --> J[Dashboard / Grafana]
```

**Dataflow ใน Go application:**
1. Config → สร้าง connection pool ไปยัง PostgreSQL (pgx)
2. Insert → ใช้ prepared statement หรือ COPY protocol สำหรับ bulk insert
3. Query → สร้าง SQL พร้อม `time_bucket()` และ window functions
4. Background jobs → ตั้งค่า policies (compression, retention, refresh continuous aggregate)

---

## ตัวอย่างโค้ดที่รันได้จริง

### โครงสร้างโปรเจกต์
```
pkg/timescaledb/
├── client.go
├── hypertable.go
├── writer.go
├── query.go
└── example_main.go
```

### 1. การติดตั้ง TimescaleDB ด้วย Docker
```bash
docker run -d --name timescaledb -p 5432:5432 \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=metricsdb \
  timescale/timescaledb:latest-pg14
```

### 2. ติดตั้ง Go driver
```bash
go get github.com/jackc/pgx/v5
```

### 3. ตัวอย่างโค้ด: การสร้าง hypertable

```go
// client.go
package timescaledb

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v5/pgxpool"
    "time"
)

type TimescaleDB struct {
    pool *pgxpool.Pool
}

func NewTimescaleDB(connString string) (*TimescaleDB, error) {
    config, err := pgxpool.ParseConfig(connString)
    if err != nil {
        return nil, err
    }
    // ตั้งค่า connection pool
    config.MaxConns = 20
    config.MinConns = 5
    pool, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        return nil, err
    }
    return &TimescaleDB{pool: pool}, nil
}

func (db *TimescaleDB) Close() {
    db.pool.Close()
}

// CreateHypertable สร้าง hypertable สำหรับ metrics
func (db *TimescaleDB) CreateHypertable(ctx context.Context) error {
    // สร้างตารางปกติก่อน
    createTableSQL := `
        CREATE TABLE IF NOT EXISTS metrics (
            time        TIMESTAMPTZ NOT NULL,
            device_id   TEXT NOT NULL,
            sensor_type TEXT NOT NULL,
            value       DOUBLE PRECISION NOT NULL,
            metadata    JSONB
        );
    `
    if _, err := db.pool.Exec(ctx, createTableSQL); err != nil {
        return err
    }

    // แปลงเป็น hypertable (ใช้ time เป็น partition key, chunk size 1 วัน)
    createHypertableSQL := `
        SELECT create_hypertable('metrics', 'time', 
            chunk_time_interval => INTERVAL '1 day',
            if_not_exists => TRUE
        );
    `
    if _, err := db.pool.Exec(ctx, createHypertableSQL); err != nil {
        return err
    }

    // สร้าง index เพื่อเพิ่มประสิทธิภาพ
    createIndexSQL := `
        CREATE INDEX IF NOT EXISTS idx_metrics_device_time 
        ON metrics (device_id, time DESC);
    `
    _, err := db.pool.Exec(ctx, createIndexSQL)
    return err
}
```

### 4. ตัวอย่างโค้ด: การเขียนข้อมูล

```go
// writer.go
package timescaledb

import (
    "context"
    "time"
)

type Metric struct {
    Time       time.Time
    DeviceID   string
    SensorType string
    Value      float64
    Metadata   map[string]interface{}
}

// WriteMetric เขียน 1 record
func (db *TimescaleDB) WriteMetric(ctx context.Context, m Metric) error {
    sql := `
        INSERT INTO metrics (time, device_id, sensor_type, value, metadata)
        VALUES ($1, $2, $3, $4, $5)
    `
    _, err := db.pool.Exec(ctx, sql, m.Time, m.DeviceID, m.SensorType, m.Value, m.Metadata)
    return err
}

// WriteBatch เขียนหลาย record แบบ batch
func (db *TimescaleDB) WriteBatch(ctx context.Context, metrics []Metric) error {
    batch := &pgx.Batch{}
    sql := `INSERT INTO metrics (time, device_id, sensor_type, value, metadata) VALUES ($1, $2, $3, $4, $5)`
    for _, m := range metrics {
        batch.Queue(sql, m.Time, m.DeviceID, m.SensorType, m.Value, m.Metadata)
    }
    results := db.pool.SendBatch(ctx, batch)
    defer results.Close()
    for range metrics {
        if _, err := results.Exec(); err != nil {
            return err
        }
    }
    return nil
}

// BulkCopy ใช้ COPY protocol เพื่อความเร็วสูงสุด
func (db *TimescaleDB) BulkCopy(ctx context.Context, metrics []Metric) error {
    copyFrom := pgx.CopyFromSlice(len(metrics), func(i int) ([]any, error) {
        m := metrics[i]
        return []any{m.Time, m.DeviceID, m.SensorType, m.Value, m.Metadata}, nil
    })
    _, err := db.pool.CopyFrom(ctx, pgx.Identifier{"metrics"}, []string{"time", "device_id", "sensor_type", "value", "metadata"}, copyFrom)
    return err
}
```

### 5. ตัวอย่างโค้ด: การ query และ time_bucket

```go
// query.go
package timescaledb

import (
    "context"
    "time"
)

type AggregatedMetric struct {
    Bucket    time.Time
    DeviceID  string
    AvgValue  float64
    MaxValue  float64
    MinValue  float64
    Count     int
}

// GetHourlyAverage ดึงค่าเฉลี่ยรายชั่วโมงของ device_id ที่ระบุ
func (db *TimescaleDB) GetHourlyAverage(ctx context.Context, deviceID string, startTime, endTime time.Time) ([]AggregatedMetric, error) {
    sql := `
        SELECT 
            time_bucket('1 hour', time) AS bucket,
            device_id,
            AVG(value) AS avg_value,
            MAX(value) AS max_value,
            MIN(value) AS min_value,
            COUNT(*) AS count
        FROM metrics
        WHERE device_id = $1
            AND time >= $2 AND time < $3
        GROUP BY bucket, device_id
        ORDER BY bucket ASC
    `
    rows, err := db.pool.Query(ctx, sql, deviceID, startTime, endTime)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []AggregatedMetric
    for rows.Next() {
        var m AggregatedMetric
        err := rows.Scan(&m.Bucket, &m.DeviceID, &m.AvgValue, &m.MaxValue, &m.MinValue, &m.Count)
        if err != nil {
            return nil, err
        }
        results = append(results, m)
    }
    return results, nil
}

// GetLastValues ดึงค่าล่าสุดของแต่ละ device และ sensor type
func (db *TimescaleDB) GetLastValues(ctx context.Context) (map[string]float64, error) {
    sql := `
        SELECT DISTINCT ON (device_id, sensor_type) 
            device_id, sensor_type, value, time
        FROM metrics
        ORDER BY device_id, sensor_type, time DESC
    `
    rows, err := db.pool.Query(ctx, sql)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    result := make(map[string]float64)
    for rows.Next() {
        var deviceID, sensorType string
        var value float64
        var t time.Time
        if err := rows.Scan(&deviceID, &sensorType, &value, &t); err != nil {
            return nil, err
        }
        key := deviceID + ":" + sensorType
        result[key] = value
    }
    return result, nil
}
```

### 6. ตัวอย่างการตั้งค่า Compression และ Retention Policies

```go
// policy.go
package timescaledb

import "context"

// EnableCompression เปิดใช้งาน compression สำหรับ hypertable (ข้อมูลเก่ากว่า 7 วัน)
func (db *TimescaleDB) EnableCompression(ctx context.Context) error {
    // ตั้งค่า segment by (column ที่ใช้จัดกลุ่ม compression)
    compressSQL := `SELECT add_compression_policy('metrics', INTERVAL '7 days', if_not_exists => true);`
    _, err := db.pool.Exec(ctx, compressSQL)
    return err
}

// SetRetentionPolicy ลบข้อมูลที่เก่ากว่า 30 วัน
func (db *TimescaleDB) SetRetentionPolicy(ctx context.Context) error {
    retentionSQL := `SELECT add_retention_policy('metrics', INTERVAL '30 days', if_not_exists => true);`
    _, err := db.pool.Exec(ctx, retentionSQL)
    return err
}

// CreateContinuousAggregate สร้าง materialized view สำหรับ downsampling รายชั่วโมง
func (db *TimescaleDB) CreateContinuousAggregate(ctx context.Context) error {
    createCaggSQL := `
        CREATE MATERIALIZED VIEW IF NOT EXISTS metrics_hourly
        WITH (timescaledb.continuous) AS
        SELECT 
            time_bucket('1 hour', time) AS bucket,
            device_id,
            sensor_type,
            AVG(value) AS avg_value,
            MAX(value) AS max_value,
            MIN(value) AS min_value
        FROM metrics
        GROUP BY bucket, device_id, sensor_type;
    `
    if _, err := db.pool.Exec(ctx, createCaggSQL); err != nil {
        return err
    }
    // ตั้งค่า refresh policy (ทุก 2 ชั่วโมง)
    refreshSQL := `SELECT add_continuous_aggregate_policy('metrics_hourly',
        start_offset => INTERVAL '3 hours',
        end_offset => INTERVAL '1 hour',
        schedule_interval => INTERVAL '2 hours');`
    _, err := db.pool.Exec(ctx, refreshSQL)
    return err
}
```

### 7. ตัวอย่างการใช้งานรวมใน HTTP server

```go
// main.go
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"
    "yourproject/pkg/timescaledb"
)

var db *timescaledb.TimescaleDB

func main() {
    var err error
    db, err = timescaledb.NewTimescaleDB("postgres://postgres:password@localhost:5432/metricsdb?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    ctx := context.Background()
    if err := db.CreateHypertable(ctx); err != nil {
        log.Printf("Warning: %v", err)
    }

    http.HandleFunc("/metrics", postMetrics)
    http.HandleFunc("/query/hourly", queryHourly)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func postMetrics(w http.ResponseWriter, r *http.Request) {
    var req struct {
        DeviceID   string                 `json:"device_id"`
        SensorType string                 `json:"sensor_type"`
        Value      float64                `json:"value"`
        Metadata   map[string]interface{} `json:"metadata"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), 400)
        return
    }
    m := timescaledb.Metric{
        Time:       time.Now(),
        DeviceID:   req.DeviceID,
        SensorType: req.SensorType,
        Value:      req.Value,
        Metadata:   req.Metadata,
    }
    if err := db.WriteMetric(r.Context(), m); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func queryHourly(w http.ResponseWriter, r *http.Request) {
    deviceID := r.URL.Query().Get("device_id")
    if deviceID == "" {
        http.Error(w, "missing device_id", 400)
        return
    }
    end := time.Now()
    start := end.Add(-24 * time.Hour)
    results, err := db.GetHourlyAverage(r.Context(), deviceID, start, end)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(results)
}
```

---

## วิธีใช้งาน module นี้

1. **ติดตั้ง TimescaleDB** (ใช้ Docker หรือติดตั้งบน Ubuntu ด้วย `apt install timescaledb-2-postgresql-14`)
2. **สร้าง database** และเปิด extension:
   ```sql
   CREATE DATABASE metricsdb;
   \c metricsdb
   CREATE EXTENSION IF NOT EXISTS timescaledb;
   ```
3. **ติดตั้ง Go driver**:
   ```bash
   go get github.com/jackc/pgx/v5
   ```
4. **คัดลอกโค้ด** ไฟล์ `client.go`, `hypertable.go`, `writer.go`, `query.go`, `policy.go` ไปไว้ใน `pkg/timescaledb/`
5. **ปรับ connection string** ให้ถูกต้อง
6. **เรียกใช้งาน** ตามตัวอย่างใน `main.go`

---

## ตารางสรุป TimescaleDB Components

| Component | คำอธิบาย | ตัวอย่าง |
|-----------|----------|----------|
| **Hypertable** | ตารางหลักที่ถูกแบ่งเป็น chunks อัตโนมัติ | `CREATE TABLE conditions ...; SELECT create_hypertable('conditions', 'time');` |
| **Chunk** | ส่วนย่อยของ hypertable (partition by time) | `_timescaledb_internal._hyper_1_2_chunk` |
| **Chunk Time Interval** | ระยะเวลาของแต่ละ chunk | `chunk_time_interval => INTERVAL '7 days'` |
| **Space Partition** | การ partition เพิ่มเติมด้วยคอลัมน์ (เช่น device_id) | `partitioning_column => 'device_id', number_partitions => 4` |
| **Continuous Aggregate** | Materialized view สำหรับ pre-aggregate | `CREATE MATERIALIZED VIEW hourly_avg WITH (timescaledb.continuous) AS ...` |
| **Compression** | บีบอัดข้อมูล segment-by-column | `ALTER TABLE metrics SET (timescaledb.compress = true); SELECT add_compression_policy('metrics', INTERVAL '7d');` |
| **Retention Policy** | ลบ chunks ที่เก่ากว่าเวลาที่กำหนด | `SELECT add_retention_policy('metrics', INTERVAL '30 days');` |
| **Time Bucket** | ฟังก์ชันจับกลุ่มตามเวลา | `time_bucket('1 hour', time)` |
| **Locf / Interpolate** | เติมค่าที่ขาดหาย | `locf(avg_value)` หรือ `interpolate(avg_value)` |

---

## แบบฝึกหัดท้าย module (3 ข้อ)

### ข้อ 1: การออก Schema และ Bulk Insert
บริษัทมีเซ็นเซอร์ 500 ตัว ส่งค่าอุณหภูมิทุก 5 วินาที จงออกแบบ hypertable ที่เหมาะสม (กำหนด chunk_time_interval, index) และเขียนฟังก์ชัน `BulkInsertFromCSV(filePath string)` ที่อ่านไฟล์ CSV (time, device_id, temp, humidity) แล้วใช้ `COPY` protocol เขียนข้อมูลลง TimescaleDB

### ข้อ 2: การ Query แบบ Window Function และ Time Bucket
ให้เขียนฟังก์ชัน `GetMovingAverage(deviceID string, duration time.Duration, windowMinutes int)` ที่คืนค่า moving average ของ temperature ทุกๆ duration (เช่น ทุก 1 ชั่วโมง) โดยใช้ window function `AVG(...) OVER (ORDER BY time ROWS BETWEEN x PRECEDING AND CURRENT ROW)` หรือใช้ time_bucket ร่วมกับ self-join

### ข้อ 3: Continuous Aggregate และ Compression
จาก hypertable `sensor_data` ที่มี 10 million rows ต่อวัน:
- สร้าง continuous aggregate สำหรับข้อมูลรายชั่วโมง (avg, max, min)
- ตั้งค่า compression policy ให้บีบอัดข้อมูลที่เก่ากว่า 3 วัน
- ตั้งค่า retention policy ให้ลบข้อมูลดิบ (raw) ที่เก่ากว่า 30 วัน แต่เก็บ aggregate ไว้ 1 ปี
- เขียน SQL และ Go code สำหรับสร้าง policies เหล่านี้

---

## แหล่งอ้างอิง

- [TimescaleDB Official Documentation](https://docs.timescale.com/)
- [TimescaleDB + Go (pgx) Guide](https://docs.timescale.com/use-timescale/latest/integrations/go/)
- [pgx GitHub Repository](https://github.com/jackc/pgx)
- [TimescaleDB Best Practices for Schema Design](https://docs.timescale.com/use-timescale/latest/schema-management/)
- [Continuous Aggregates Guide](https://docs.timescale.com/use-timescale/latest/continuous-aggregates/)

---

**หมายเหตุ:** module นี้ครบถ้วนสำหรับ `pkg/timescaledb` สำหรับระบบ gobackend หากต้องการ module เพิ่มเติม (เช่น `pkg/clickhouse`, `pkg/questdb`) โปรดแจ้ง