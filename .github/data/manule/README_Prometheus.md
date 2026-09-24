# คู่มือการใช้งาน Prometheus (Monitoring)

Prometheus เป็นระบบเก็บ metrics/alerting ของโปรเจกต์ `icmongolang` — scrape ข้อมูลจาก exporter ต่าง ๆ (PostgreSQL, Redis, Kafka, MQTT, Elasticsearch, InfluxDB, Ollama, Go backend) และนำเสนอผ่าน **Grafana**

**สถาปัตยกรรม:**

```
┌───────────────┐   scrape /metrics   ┌───────────────┐
│  Exporters    │ ──────────────────▶ │  Prometheus   │
│  kafka :9308  │                     │  :9090        │
│  redis :9121  │                     └───────┬───────┘
│  postgres:9187│                             │
│  mosquitto:9234│                            │  query
│  influxdb:8086 │                            ▼
│  elasticsearch:9114│                  ┌───────────────┐
│  ollama-exporter:9309│                │  Grafana :3000│
│  backend:5000  │                      └───────────────┘
│  blackbox:9115 │
└───────────────┘
```

---

## 1. การตั้งค่า

### 1.1 Docker Compose (`docker-compose.yml:94-110`)
```yaml
prometheus:
  image: prom/prometheus:latest
  ports:
    - "9090:9090"
  volumes:
    - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml:ro
    - app-prometheus-data:/prometheus
  depends_on: [kafka-exporter, redis-exporter, postgres-exporter, blackbox-exporter, mosquitto-exporter, influxdb, elasticsearch-exporter, ollama-exporter]
```

### 1.2 รัน
```bash
docker compose up -d prometheus
```

### 1.3 UI
- **Prometheus UI:** http://localhost:9090
- **Targets:** http://localhost:9090/targets
- **Query:** http://localhost:9090/graph

---

## 2. `monitoring/prometheus.yml`

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s
```

### 2.1 Scrape jobs

| job_name | target | หมายเหตุ |
|----------|--------|----------|
| `kafka-exporter` | `kafka-exporter:9308` | |
| `redis-exporter` | `redis-exporter:9121` | |
| `postgres-exporter` | `postgres-exporter:9187` | |
| `mosquitto-exporter` | `mosquitto-exporter:9234` | |
| `influxdb` | `influxdb:8086` (`/metrics`) | |
| `elasticsearch-exporter` | `elasticsearch-exporter:9114` | |
| `ollama-exporter` | `ollama-exporter:9309` | |
| `backend` | `backend:5000` (`/metrics`) | Go app prometheus metrics |
| `backend-blackbox` | `http://backend:5000/health`, `/apimetric` | probe ผ่าน blackbox `http_2xx` |

### 2.2 Blackbox relabel
- `__address__` → `__param_target` → `instance`
- `__address__` ถูกแทนเป็น `blackbox-exporter:9115` (Prometheus เรียก probe ที่ blackbox แทน)

> **⚠️ ข้อควรระวัง:** job `backend` และ `backend-blackbox` ชี้ไปที่ hostname **`backend`** แต่ docker-compose ตั้งชื่อ service ว่า **`api`** — งานทั้งสองจะ **scrape ไม่ได้ (DNS fail)** จนกว่าจะ rename/alias service

---

## 3. Metrics ที่เก็บ

### 3.1 จาก Go backend (`/metrics`)
- `http_requests_total{method,path,status}` — จำนวน request (middleware monitoring)
- `http_request_duration_seconds{method,path}` — เวลาตอบสนอง
- `go_goroutines`, `process_start_time_seconds`, `process_resident_memory_bytes` — runtime
- `ws_connected_clients`, `ws_active_rooms`, `ws_active_topics`, `ws_connections_total`, `ws_messages_broadcast_total{type}`, `ws_messages_received_total{type}`, `ws_messages_sent_total`, `ws_messages_dropped_total` — WebSocket (pkg/websocket/metrics.go)

### 3.2 จาก `/apimetric` (JSON snapshot)
`total_requests`, `active_requests`, `total_errors`, `error_rate`, `uptime_seconds`, `qps`

### 3.3 จาก exporters
แต่ละ exporter มี metric ของตัวเอง (ดูคู่มือ exporter แยก: `README_Exporters.md`)

---

## 4. Query ตัวอย่าง

```promql
# HTTP request rate (5m)
sum(rate(http_requests_total{job="backend"}[5m]))

# Error rate 4xx/5xx
sum(rate(http_requests_total{job="backend",status=~"4..|5.."}[5m])) /
sum(rate(http_requests_total{job="backend"}[5m]))

# p95 latency
histogram_quantile(0.95, sum by (le) (rate(http_request_duration_seconds_bucket{job="backend"}[5m])))

# Uptime (วัน)
time() - process_start_time_seconds{job="backend"}

# MQTT clients connected
broker_clients_connected

# Redis memory %
redis_memory_used_bytes / redis_memory_max_bytes * 100

# PostgreSQL connections
sum(pg_stat_database_numbackends)
```

---

## 5. Edge Cases & ข้อควรระวัง

1. **ไม่มี rule_files / alerting config** — ระบบยังไม่มี alert rules; ใช้ Grafana alert เองถ้าต้องการ
2. **backend hostname ผิด** — ต้อง rename service เป็น `backend` หรือเพิ่ม network alias ไม่งั้น 2 jobs fail
3. **expose (ไม่ publish port)** ของ exporter ส่วนใหญ่ — เข้าถึงได้จาก docker network เท่านั้น (ยกเว้น prometheus 9090)
4. **scrape_interval 15s** — data granularity จำกัด 15s
5. **blackbox ใช้ module `http_2xx`** จาก default config ของ image (ไม่มี config file แยก)
6. **influxdb ใช้ `/metrics`** (มี `metrics_path` ระบุ) — influxdb 2.7 expose prometheus endpoint ที่ `/metrics`
7. **volume `app-prometheus-data`** — ข้อมูล time-series เก็บไว้; ล้างด้วย `docker compose down -v`

---

## 6. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| target down ใน UI | exporter container ไม่รัน → `docker compose ps`; ตรวจ depends_on |
| `backend`/`backend-blackbox` down | hostname ผิด (ใช้ `backend` แต่ service ชื่อ `api`) |
| query ไม่มี data | เช็คที่ `/targets` ว่า target UP; ตรวจ labels job ชื่อถูกต้อง |
| ข้อมูลหายหลัง restart | ใช้ volume `app-prometheus-data` (ห้าม `-v` ถ้าไม่ต้องการล้าง) |
