# คู่มือการใช้งานโมดูล InfluxDB (Time Series)

โมดูล `internal/modules/influxdb` ทำหน้าที่อ่าน/เขียนข้อมูล **Time Series** (ข้อมูลแบบอนุกรมเวลา เช่น ค่าอุณหภูมิ, กระแส, แรงดันจากอุปกรณ์ IoT) ลงใน **InfluxDB 2.x** โดยใช้ภาษา query แบบ **Flux** ผ่าน library `github.com/influxdata/influxdb-client-go/v2`

ฟีเจอร์หลักของโมดูล:
- **Write** — เขียน data point ใหม่ (measurement + fields + tags)
- **Query** — ดึงข้อมูลดิบตาม measurement/field/ช่วงเวลา พร้อม limit/offset
- **Device Chart** — ดึงข้อมูลเพื่อวาดกราฟ (เรียงข้อมูล + ค่าล่าสุด)
- **Statistics** — คำนวณสถิติ (mean, median, mode, last, first, stddev, variance, percentile)

---

## 1. สถาปัตยกรรม (Architecture)

โมดูลแบ่งเป็น 3 ชั้น เช่นเดียวกับโมดูลอื่นในโปรเจกต์:

```
Client / Postman / Frontend
        │ HTTP (chi router)
        ▼
┌─────────────────────────────────────────────────┐
│ delivery/http                                     │
│   handler.go  → รับ HTTP request / ส่ง response    │
│   routes.go   → ลงทะเบียน route ใต้ /api/influx     │
└─────────────────────────────────────────────────┘
        │
        ▼
┌─────────────────────────────────────────────────┐
│ usecase                                          │
│   usecase.go  → ประมวลผล แปลง timezone เป็น UTC+7 │
└─────────────────────────────────────────────────┘
        │
        ▼
┌─────────────────────────────────────────────────┐
│ pkg/influxdb (client.go)                          │
│  influxdb-client-go/v2  → WriteAPI + QueryAPI      │
│  สร้าง Flux query + executeQuery                   │
└─────────────────────────────────────────────────┘
        │
        ▼
        InfluxDB 2.x (port 8086 / host mapping 9087)
```

### 1.1 ตารางไฟล์

| ไฟล์ | บทบาท |
|------|-------|
| `internal/modules/influxdb/presenter/presenter.go` | DTO — Request / Response ของทุก endpoint |
| `internal/modules/influxdb/usecase/usecase.go` | Business logic — เรียก client, แปลงเวลาเป็นโซนไทย |
| `internal/modules/influxdb/delivery/http/handler.go` | HTTP handler + ตรวจสอบค่า + error mapping |
| `internal/modules/influxdb/delivery/http/routes.go` | ลงทะเบียน route ใต้ `/api/influx` |
| `pkg/influxdb/client.go` | Wrapper ของ influxdb-client-go/v2 — สร้าง Flux query, execute query |
| `pkg/helpers/format.go` | Helper เรื่อง timezone (Asia/Bangkok) และรูปแบบเวลา |

### 1.2 จุดที่โมดูลถูก connect (wiring)

- **`internal/server/server.go`** — สร้าง `influxClient` ผ่าน `influxdb.NewInfluxClient(cfg, log)`; ถ้า `INFLUXDB_URL` หรือ `INFLUXDB_TOKEN` ไม่ครบ (หรือเชื่อมไม่ได้) จะได้ `nil` และระบบหลักยังทำงานต่อได้
- **`internal/server/handlers.go`** — ถ้า `influxClient != nil` จะสร้าง UseCase + Handler แล้วลง route ผ่าน `MapInfluxRoutes(apiRouter, influxHandler, mw)` ไม่งั้นข้ามพร้อม log คำเตือน
- Base path ของ API ทั้งหมดคือ `/api` → endpoints นี้อยู่ใต้ `/api/influx/...`

---

## 2. แนวคิดของ InfluxDB ที่ควรรู้ก่อนใช้

InfluxDB เก็บข้อมูลเป็น **Point** ซึ่งประกอบด้วย:

| ส่วน | ตัวอย่าง | คำอธิบาย |
|------|----------|----------|
| Measurement | `temperature` | ชื่อกลุ่มข้อมูล (คล้าย table) |
| Tag | `device="PLC-01"`, `location="A1"` | ใช้ระบุ attribute (indexed — ใช้ filter ได้เร็ว) |
| Field | `value=27.5`, `status="on"` | ค่าที่วัดจริง (ตัวเลข/สตริง) |
| Timestamp | `2026-08-14T10:00:00Z` | เวลาเกิดเหตุการณ์ |

ข้อมูลถูกเก็บใน **Bucket** (เช่น `icmongolang` หรือ `AIRCOM1`) ภายใต้ **Org** (เช่น `icmongolang`)

โครงสร้างข้อมูลมาตรฐานที่โมดูลนี้ใช้:
- field `_value` — ค่าของข้อมูล
- field `_time` — เวลา
- field `_measurement`, `_field` — ใช้ filter

---

## 3. การตั้งค่า Configuration

### 3.1 ตัวแปรสภาพแวดล้อม (`.env`)

```env
# ---- InfluxDB 2.x ----
INFLUXDB_URL=http://localhost:9087      # ระวัง: ใน Docker จะถูก override เป็น http://influxdb:8086
INFLUXDB_USERNAME=admin                 # ใช้ Login UI ของ InfluxDB (ไม่ได้ใช้ยิง API)
INFLUXDB_PASSWORD=admin1234
INFLUXDB_TOKEN=icmongolang-super-secret-token   # ใช้ยิง API จริง (ต้องไม่ว่าง)
INFLUXDB_ORG=icmongolang                # Organization
INFLUXDB_BUCKET=icmongolang             # Bucket หลัก (default เมื่อไม่ระบุใน request)
INFLUXDB_TIMEOUT=30                     # HTTP timeout หน่วยวินาที (default 30)
```

> **สำคัญ:** `INFLUXDB_URL` + `INFLUXDB_TOKEN` เป็นคู่ที่ต้องครบถ้วน ถ้าขาดตัวใดตัวหนึ่ง `NewInfluxClient` จะ return error ทันที (`influxdb config missing`) และโมดูลจะถูกข้าม

### 3.2 ค่าเริ่มต้นใน `config/config.default.yml`

```yaml
influxdb:
  url: "http://localhost:9087"
  token: ""
  org: "cmon_org"
  bucket: "AIRCOM1"
  timeout: 30
  username: ""
  password: ""
```

### 3.3 หมายเหตุ port

- ระบบ InfluxDB จริงเปิด port `8086`
- ในเครื่อง dev จะ map มาที่ `localhost:9087` (native run) — ดู `docker-compose.yml` ที่ map `9087:8086`

---

## 4. การยืนยันตัวตน (Authentication)

**ทุก endpoint ของโมดูลนี้ต้องใช้ JWT** — route ทั้งหมดผ่าน middleware:
`Verifier(true)` → `Authenticator()` → `CurrentUser()` → `ActiveUser()`
(ดู `internal/middleware/jwtauth.go`)

ส่ง token ใน header:

```
Authorization: Bearer <ACCESS_TOKEN>
```

**วิธีขอ token** (POST `/api/auth/login` — multipart/form-data):

```bash
curl -X POST http://localhost:5000/api/auth/login \
  -F "username=root@gmail.com" \
  -F "password=root_password"
```

นำค่า `access_token` ที่ได้ไปใช้ (ตัวอย่างแทนด้วย `$TOKEN`)

---

## 5. API Reference

### 5.1 `POST /api/influx/write` — เขียน data point

- **ต้อง auth**
- เขียน point 1 จุดด้วย timestamp = เวลาปัจจุบัน (`time.Now()`)

**Request body:**

```json
{
  "measurement": "temperature",
  "fields": { "value": 27.5, "humidity": 62 },
  "tags": { "device": "PLC-01", "location": "A1" }
}
```

| Field | ประเภท | บังคับ | คำอธิบาย |
|-------|--------|-------|----------|
| `measurement` | string | ✅ | ชื่อ measurement (เช่น `temperature`, `power`) |
| `fields` | map[string]interface{} | ✅ | ค่าที่จะเก็บ (ตัวเลข/สตริง/boolean) |
| `tags` | map[string]string | ❌ | Attribute ใช้กรองข้อมูล (ต้องเป็น string เท่านั้น) |

```bash
curl -X POST http://localhost:5000/api/influx/write \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "measurement": "temperature",
    "fields": { "value": 27.5 },
    "tags": { "device": "PLC-01" }
  }'
```

**Response 200:**

```json
{ "status": "ok" }
```

> หมายเหตุ: ใช้ `WriteAPI` + `Flush()` ทันที → ข้อมูลถูกส่งไป InfluxDB ทันที (write synchronously)

---

### 5.2 `GET /api/influx/query` — ค้นหาข้อมูลดิบ (Query Parameters)

- **ต้อง auth**
- อ่านค่า filter จาก URL query parameters (วิธีเดียวกับ POST /filters)

| Parameter | ประเภท | บังคับ | ค่าเริ่มต้น | คำอธิบาย |
|-----------|--------|-------|------------|----------|
| `measurement` | string | ✅ | — | ชื่อ measurement |
| `field` | string | ✅ | — | ชื่อ field (เช่น `value`) |
| `start` | string | ❌ | `-1h` | เวลาเริ่ม (relative เช่น `-1h`, `-30d` หรือ absolute `2026-01-01T00:00:00Z`) |
| `stop` | string | ❌ | `now()` | เวลาสิ้นสุด |
| `limit` | int | ❌ | `1000` | จำนวน record สูงสุด |
| `offset` | int | ❌ | `0` | ข้าม record ต้น (pagination) |

```bash
curl -G http://localhost:5000/api/influx/query \
  -H "Authorization: Bearer $TOKEN" \
  --data-urlencode "measurement=temperature" \
  --data-urlencode "field=value" \
  --data-urlencode "start=-2h" \
  --data-urlencode "stop=now()" \
  --data-urlencode "limit=100"
```

**Response 200** (array ของ data point — เวลาถูกแปลงเป็นโซน **Asia/Bangkok UTC+7** แล้ว รูปแบบ `2006-01-02 15:04:05`):

```json
[
  { "time": "2026-08-14 13:30:00", "value": 27.5 },
  { "time": "2026-08-14 13:30:10", "value": 27.6 },
  { "time": "2026-08-14 13:30:20", "value": 27.4 }
]
```

**Flux query ที่ถูกสร้าง** (ดู `pkg/influxdb/client.go:184`):

```flux
from(bucket: "icmongolang")
  |> range(start: -2h, stop: now())
  |> filter(fn: (r) => r["_measurement"] == "temperature")
  |> filter(fn: (r) => r["_field"] == "value")
  |> limit(n: 100, offset: 0)
  |> yield(name: "filtered_data")
```

---

### 5.3 `POST /api/influx/filters` — ค้นหาข้อมูลดิบ (JSON body)

- **ต้อง auth**
- เหมือน `GET /query` แต่รับค่าเป็น JSON body แทน

**Request body:**

```json
{
  "measurement": "temperature",
  "field": "value",
  "bucket": "icmongolang",
  "start": "-1h",
  "stop": "now()",
  "limit": 100,
  "offset": 0
}
```

| Field | ประเภท | บังคับ | ค่าเริ่มต้น | คำอธิบาย |
|-------|--------|-------|------------|----------|
| `measurement` | string | ✅ | — | ชื่อ measurement |
| `field` | string | ✅ | — | ชื่อ field |
| `bucket` | string | ❌ | bucket จาก config | ชื่อ bucket (เช่น `AIRCOM1`, `icmongolang`) |
| `start` | string | ❌ | `-1h` | เวลาเริ่ม |
| `stop` | string | ❌ | `now()` | เวลาสิ้นสุด |
| `limit` | int | ❌ | `1000` | จำนวน record สูงสุด |
| `offset` | int | ❌ | `0` | ข้าม record ต้น |

```bash
curl -X POST http://localhost:5000/api/influx/filters \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "measurement": "temperature", "field": "value", "start": "-2h", "limit": 100 }'
```

**Response 200:** เหมือน `GET /query` — array ของ `{ "time", "value" }`

---

### 5.4 `POST /api/influx/devicechart` — ข้อมูลสำหรับวาดกราฟอุปกรณ์

- **ต้อง auth**
- ใช้ Query เดียวกับ `/filters` แต่แปลงผลลัพธ์เป็นโครงสร้าง chart พร้อมข้อมูลเมตา (ค่าล่าสุด, เวลาล่าสุด)
- ค่าตัวเลขจะถูกแปลงเป็น **int64** (ทศนิยมถูกตัด) → เหมาะกับข้อมูลกราฟที่ต้องการความเรียบง่าย

**Request body:** เหมือน `QueryFilterRequest` ทุกประการ

```json
{
  "measurement": "temperature",
  "field": "value",
  "bucket": "icmongolang",
  "start": "-1h",
  "stop": "now()",
  "limit": 100
}
```

**Response 200 (`DeviceChartResponse`):**

```json
{
  "bucket": "icmongolang",
  "field": "value",
  "info": {
    "bucket": "icmongolang",
    "measurement": "temperature",
    "result": "last",
    "table": 0,
    "field": "value",
    "start": "-1h",
    "stop": "now()",
    "time": "2026-08-14 13:30:20",
    "value": 27
  },
  "data": [27, 27, 27],
  "date": ["2026-08-14 13:30:00", "2026-08-14 13:30:10", "2026-08-14 13:30:20"],
  "name": "value",
  "cache": "cache"
}
```

| Field | คำอธิบาย |
|-------|----------|
| `info.value` | ค่าล่าสุด (int64) |
| `info.time` | เวลาล่าสุด (Bangkok timezone) |
| `data` | array ค่าทั้งหมด (int64) — เรียงตามเวลาที่ query กลับมา |
| `date` | array เวลาของแต่ละจุด (รูปแบบ `2006-01-02 15:04:05`) |
| `name` | ชื่อ field |

```bash
curl -X POST http://localhost:5000/api/influx/devicechart \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "measurement": "temperature", "field": "value", "start": "-1h" }'
```

> หมายเหตุ: field `result` และ `cache` เป็นค่า hardcode (`"last"`, `"cache"`) ใน code — ยังไม่ dynamic

---

### 5.5 `POST /api/influx/statistics` — คำนวณสถิติ

- **ต้อง auth**
- คำนวณค่าสถิติจากข้อมูลตามช่วงเวลา โดยเลือกวิธี (`aggregate`) ที่ต้องการ

**Request body:**

```json
{
  "measurement": "temperature",
  "field": "value",
  "bucket": "icmongolang",
  "start": "-15s",
  "stop": "now()",
  "aggregate": "mean",
  "windowPeriod": "15s",
  "percentile": 95
}
```

| Field | ประเภท | บังคับ | ค่าเริ่มต้น | คำอธิบาย |
|-------|--------|-------|------------|----------|
| `measurement` | string | ✅ | — | ชื่อ measurement |
| `bucket` | string | ✅ | — | ชื่อ bucket |
| `field` | string | ✅ | — | ชื่อ field |
| `start` | string | ❌ | `-15s` | เวลาเริ่ม |
| `stop` | string | ❌ | `now()` | เวลาสิ้นสุด |
| `aggregate` | string | ✅ | (ถ้าไม่ใส่เป็น `last`) | วิธีคำนวณ — ดูตารางด้านล่าง |
| `windowPeriod` | string | ❌ | `15s` | ขนาด window (ใช้กับ mean/median) |
| `percentile` | float64 | ❌ | `0.95` | ค่า percentile (ใช้กับ aggregate=`percentile`) |

**ค่าที่รองรับของ `aggregate`:**

| ค่า | Flux | ความหมาย |
|-----|------|----------|
| `mean` / `average` | `aggregateWindow(fn: mean)` | ค่าเฉลี่ย ต่อ window |
| `median` | `aggregateWindow(fn: median)` | ค่ามัธยฐาน ต่อ window |
| `mode` | count + top(1) | ค่าที่เกิดซ้ำมากที่สุด |
| `last` | `last()` | ค่าล่าสุด (default) |
| `first` | `first()` | ค่าแรกสุด |
| `stddev` / `standarddeviation` | `stddev()` | ส่วนเบี่ยงเบนมาตรฐาน |
| `variance` | `variance()` | ความแปรปรวน |
| `percentile` | `percentile(percentile: X)` | เปอร์เซ็นไทล์ (ค่าเริ่มต้น 0.95) |

**Response 200 (`StatisticsResponse`) — 2 รูปแบบ:**

ถ้าได้ผลลัพธ์ **1 จุด** → อยู่ใน `value`:

```json
{
  "aggregate": "mean",
  "value": 27.45
}
```

ถ้าได้ผลลัพธ์ **หลายจุด** (เช่น aggregateWindow หลาย window) → อยู่ใน `data`:

```json
{
  "aggregate": "mean",
  "data": [
    { "time": "2026-08-14 13:30:00", "value": 27.4 },
    { "time": "2026-08-14 13:30:15", "value": 27.6 }
  ]
}
```

```bash
curl -X POST http://localhost:5000/api/influx/statistics \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "measurement": "temperature", "field": "value", "bucket": "icmongolang", "aggregate": "mean", "start": "-1h" }'
```

**ตัวอย่างเพิ่มเติม:**

```bash
# ค่าล่าสุด
-d '{ "measurement": "power", "field": "value", "bucket": "icmongolang", "aggregate": "last", "start": "-5m" }'

# ค่าเฉลี่ยทุก 1 ชั่วโมง ย้อนหลัง 1 วัน
-d '{ "measurement": "power", "field": "value", "bucket": "icmongolang", "aggregate": "mean", "windowPeriod": "1h", "start": "-24h" }'

# เปอร์เซ็นไทล์ที่ 95
-d '{ "measurement": "latency", "field": "value", "bucket": "icmongolang", "aggregate": "percentile", "percentile": 95, "start": "-1h" }'
```

---

## 6. ตัวอย่างการใช้งานครบวงจร

```bash
# 1) Login เอา token
TOKEN=$(curl -s -X POST http://localhost:5000/api/auth/login \
  -F "username=root@gmail.com" -F "password=root_password" | jq -r .access_token)

# 2) เขียนข้อมูลอุณหภูมิจาก PLC-01
curl -X POST http://localhost:5000/api/influx/write \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{ "measurement": "temperature", "fields": { "value": 27.5 }, "tags": { "device": "PLC-01", "location": "A1" } }'

# 3) ดึงข้อมูลย้อนหลัง 1 ชั่วโมง
curl -G http://localhost:5000/api/influx/query \
  -H "Authorization: Bearer $TOKEN" \
  --data-urlencode "measurement=temperature" \
  --data-urlencode "field=value" \
  --data-urlencode "start=-1h"

# 4) วาดกราฟ
curl -X POST http://localhost:5000/api/influx/devicechart \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{ "measurement": "temperature", "field": "value", "start": "-1h" }'

# 5) ค่าเฉลี่ยย้อนหลัง 1 ชั่วโมง
curl -X POST http://localhost:5000/api/influx/statistics \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{ "measurement": "temperature", "field": "value", "bucket": "icmongolang", "aggregate": "mean", "start": "-1h" }'
```

---

## 7. การเรียกใช้จากโค้ด Go (โดยตรง ไม่ผ่าน HTTP)

```go
import (
    "icmongolang/config"
    "icmongolang/pkg/influxdb"
)

cfg := config.Load(...)           // โหลด config ที่มี InfluxDB block

client, err := influxdb.NewInfluxClient(cfg, logger)
if err != nil {
    // URL หรือ Token ไม่ครบ
}

// เขียนข้อมูล
client.WriteData("temperature", map[string]interface{}{"value": 27.5}, map[string]string{"device": "PLC-01"})

// query
params := influxdb.QueryParams{
    Measurement: "temperature",
    Field:       "value",
    Start:       "-1h",
    Stop:        "now()",
    Limit:       100,
}
rows, err := client.QueryFilterData(params)

// สถิติ
result := client.CalculateStatistics(influxdb.QueryParams{
    Measurement: "temperature",
    Field:       "value",
    Bucket:      "icmongolang",
    Start:       "-1h",
    Mean:        "mean",
    WindowPeriod: "15s",
})

client.Close()
```

---

## 8. พฤติกรรมและข้อควรระวัง (Edge Cases / Notes)

1. **โมดูลเป็น optional** — ถ้า config `url`/`token` ไม่ครบ หรือเชื่อมไม่ได้ → route ถูกข้าม ระบบหลักไม่ล่ม (log `⚠️ InfluxDB client is nil`)
2. **Timezone** — ทุก response ถูกแปลงเป็น **Asia/Bangkok (UTC+7)** ผ่าน `helpers.GetTimeLocation()`; ถ้าโหลด timezone ไม่ได้ จะ fallback เป็น fixed `UTC+7` (ดู `pkg/helpers/format.go`)
3. **ค่าเริ่มต้นของ query** — `start=-1h`, `stop=now()`, `limit=1000` (filters/query) และ `limit=100` (เฉพาะ `QueryFilterDataRs`) แต่ endpoint นี้ใช้ 1000
4. **`WriteData` ใช้เวลา `time.Now()` เสมอ** — ไม่มีพารามิเตอร์กำหนดเวลาใน HTTP API (แต่ `WritePoint` ใน Go รองรับการกำหนดเวลาเอง)
5. **`/devicechart` แปลงค่าเป็น int64** — ค่าทศนิยมจะถูกตัดทิ้ง (เช่น 27.9 → 27) อย่าใช้ endpoint นี้กับข้อมูลที่ต้องการความละเอียด
6. **`percentile`** — ค่าใน body ใช้เลขเปอร์เซ็นต์ตรงๆ เช่น 95, 99 (client จะใช้เป็น 0.95/0.99 ใน Flux ถ้าไม่ใส่จะ default 0.95)
7. **สถิติ `mean`/`median`** — ใช้ `aggregateWindow` ต้องระบุ `windowPeriod` เพื่อกำหนดขนาด window (default `15s`)
8. **`/statistics` response** — 1 จุด → field `value`; หลายจุด → field `data` (ดู `usecase.go:156`); field ทั้งสองเป็น `omitempty` จึงมีเพียงอันเดียว
9. **Error format** — `{ "error": "<ข้อความ>" }` พร้อม status 400 (JSON/body ผิด) หรือ 500 (query ผิด / เชื่อมต่อไม่ได้) ข้อความ error จาก Flux จะถูกส่งกลับตรงๆ
10. **ไม่มี health endpoint** — ต่างจาก Elasticsearch ที่มี `GET /health`; โมดูล InfluxDB มีเพียง 5 endpoints ข้างต้น

---

## 9. Docker

`docker-compose.yml` มี service **influxdb:2.7** map port `9087:8086` พร้อม volume `app-influxdb-data`:

```bash
docker compose up -d influxdb
```

หลังสร้างครั้งแรก เข้า UI ที่ `http://localhost:9087` (user: `admin` / `admin1234`) เพื่อ:
- สร้าง Organization (default `icmongolang`)
- สร้าง Bucket (default `icmongolang`)
- สร้าง All-Access API Token แล้วนำค่ามาใส่ `INFLUXDB_TOKEN`

> ในโหมด Docker ค่า `INFLUXDB_URL` จะถูก override เป็น `http://influxdb:8086` (เห็นใน `docker-compose.yml`)

---

## 10. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| โมดูลไม่ทำงาน (route 404) | ตรวจ log ตอน start ว่า `❌ Failed to connect InfluxDB`; เช็คว่า `INFLUXDB_URL` + `INFLUXDB_TOKEN` ครบหรือยัง |
| 401 Unauthorized | token JWT หมดอายุ / ไม่ได้ส่ง header `Authorization: Bearer` / ผู้ใช้ถูกปิด |
| query คืนผลว่าง | ยังไม่มีข้อมูลในช่วงเวลานั้น / measurement-field พิมพ์ผิด / bucket ผิด / ใช้ start-stop นอกช่วงข้อมูล |
| error จาก Flux (500) | ดูข้อความ error ที่ส่งกลับ — มักเกิดจาก syntax ช่วงเวลา (เช่น `start` เป็นค่า invalid), measurement/field ไม่มีช่องว่าง, หรือ `bucket` ไม่มีอยู่ |
| เวลาไม่ตรงกับไทย (UTC) | ตรวจว่าค่า `timezone` ใน server config เป็น `Asia/Bangkok` (ถ้าไม่ตั้ง ค่า default ใน helper คือ `Asia/Bangkok` อยู่แล้ว) |
| เขียนแล้ว query ไม่เจอ | การเขียนใช้ bucket/org จาก config — ตรวจว่า bucket ที่ query ตรงกับ bucket ที่เขียน (`INFLUXDB_BUCKET` กับค่า `bucket` ใน request) |
| port 9087 ติดไม่ได้ | ตรวจว่า container `influxdb` รันอยู่; `docker compose ps` และลอง `curl http://localhost:9087/health` |

---

## 11. โครงสร้างข้อมูลใน InfluxDB

Point ที่ถูกเขียนผ่าน `/influx/write` มีหน้าตาแบบนี้ (อ้างอิง InfluxDB 2.x line protocol):

```
temperature,device=PLC-01,location=A1 value=27.5 1755000000000000000
│          │ tag tags             │ field        │ timestamp (ns)
│          └─────────────────────────────────────┘
└ measurement
```

- **Bucket:** ค่าจาก `INFLUXDB_BUCKET` (default) หรือที่ระบุใน request
- **Org:** ค่าจาก `INFLUXDB_ORG`
- ข้อมูลที่ query กลับมา: `_time`, `_measurement`, `_field`, `_value` (+ tags ถ้ามี)
