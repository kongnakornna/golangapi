# คู่มือการตั้งค่า Node-RED Monitoring บน Grafana

คู่มือนี้อธิบายวิธีติดตั้งและใช้งาน **Node-RED Dashboard** ใน Grafana เพื่อ monitor Node-RED ที่รันอยู่ที่ `http://localhost:1881/`

**สถาปัตยกรรม:**

```
┌──────────────────────────┐
│  Node-RED :1881          │  ← ติดตั้ง node-red-contrib-prometheus
│  /metrics endpoint       │
└───────────┬──────────────┘
            │ scrape (PromQL)
┌───────────▼──────────────┐   PromQL query   ┌──────────────┐   browser
│  Prometheus :9090        │ ◀─────────────── │ Grafana :3000│──────────▶ http://localhost:3000
└──────────────────────────┘                  └──────────────┘
```

---

## 1. Prerequisites

### 1.1 ติดตั้ง `node-red-contrib-prometheus` ใน Node-RED

Node-RED ไม่มี built-in Prometheus metrics — ต้องติดตั้ง plugin เพิ่ม

**วิธีที่ 1: ติดตั้งผ่าน UI**
1. เปิด Node-RED ที่ `http://localhost:1881/`
2. ไปที่ **Menu > Manage Palette > Install**
3. ค้นหา `node-red-contrib-prometheus`
4. คลิก **Install** แล้ว Confirm
5. **Deploy** flow ใดก็ได้ (ต้อง deploy อย่างน้อย 1 ครั้งหลังติดตั้ง)

**วิธีที่ 2: ติดตั้งผ่าน CLI**
```bash
cd ~/.node-red
npm install node-red-contrib-prometheus
# restart Node-RED
node-red
```

### 1.2 ตรวจสอบว่า metrics endpoint ทำงาน

หลังติดตั้งและ deploy แล้ว ทดสอบ:
```bash
curl http://localhost:1881/metrics
```

จะได้ output แบบ Prometheus text format:
```
# HELP node_red_info Node-RED version info
# TYPE node_red_info gauge
node_red_info{version="3.1.0"} 1
# HELP node_red_nodes_total Total number of nodes
# TYPE node_red_nodes_total gauge
node_red_nodes_total 15
...
```

> **ถ้าไม่เห็น output** → ต้อง deploy flow อย่างน้อย 1 ครั้ง หลังติดตั้ง plugin

---

## 2. Prometheus Configuration

### 2.1 เพิ่ม scrape target ใน `monitoring/prometheus.yml`

ไฟล์ถูกเพิ่มไว้แล้ว:

```yaml
# Node-RED – requires node-red-contrib-prometheus plugin installed in Node-RED
- job_name: "nodered"
  metrics_path: /metrics
  static_configs:
    - targets: ["host.docker.internal:1881"]
```

> `host.docker.internal` ใช้ได้บน Docker Desktop (Windows/Mac) — Prometheus ใน Docker จะเข้าถึง Node-RED บน host ได้

### 2.2 ถ้ารัน Node-RED ใน Docker

ถ้า Node-RED อยู่ใน Docker network เดียวกับ Prometheus ให้เปลี่ยน target:

```yaml
- job_name: "nodered"
  metrics_path: /metrics
  static_configs:
    - targets: ["nodered:1881"]   # ชื่อ service ใน docker-compose
```

### 2.3 Restart Prometheus

```bash
docker compose restart prometheus
```

### 2.4 ตรวจสอบ target ใน Prometheus UI

เปิด `http://localhost:9090/targets` → ดูว่า job `nodered` ขึ้น **UP** สีเขียว

---

## 3. Grafana Dashboard

### 3.1 Provisioning

ไฟล์ provisioning อยู่ที่ `monitoring/grafana/provisioning/dashboards/nodered.yml`:

```yaml
apiVersion: 1
providers:
  - name: nodered
    orgId: 1
    folder: Node-RED
    type: file
    disableDeletion: false
    updateIntervalSeconds: 10
    allowUiUpdates: true
    options:
      path: /var/lib/grafana/dashboards/nodered
      foldersFromFilesStructure: false
```

Dashboard JSON อยู่ที่ `monitoring/grafana/dashboards/nodered/nodered-dashboard.json`

Grafana จะโหลดอัตโนมัติเมื่อ restart:
```bash
docker compose restart grafana
```

### 3.2 เปิดดู Dashboard

1. เปิด `http://localhost:3000`
2. Login: `admin` / `admin`
3. Sidebar → **Dashboards** → folder **Node-RED**
4. เลือก **Node-RED Monitoring**

---

## 4. Dashboard Panels

### 4.1 Node-RED Overview

| Panel | Metric | ความหมาย |
|-------|--------|----------|
| Node-RED Up | `up{job="nodered"}` | สถานะ UP/DOWN (สีเขียว/แดง) |
| Node-RED Uptime | `process_uptime_seconds{job="nodered"}` | เวลาที่ Node-RED ทำงานมา |
| Total Nodes | `node_red_nodes_total{job="nodered"}` | จำนวน nodes ทั้งหมด |
| Total Flows | `node_red_flows_total{job="nodered"}` | จำนวน flows ทั้งหมด |

### 4.2 Message Traffic

| Panel | Metric | ความหมาย |
|-------|--------|----------|
| Messages per Second | `rate(node_red_msg_received_total[5m])` / `rate(node_red_msg_sent_total[5m])` | อัตราข้อความเข้า/ออก (ต่อวินาที) |
| Total Messages | `node_red_msg_received_total` / `node_red_msg_sent_total` | จำนวนข้อความสะสม |

### 4.3 Nodes by State

| Panel | Metric | ความหมาย |
|-------|--------|----------|
| Nodes by State | `node_red_nodes_connected_total` / `node_red_nodes_disabled_total` | nodes ที่เชื่อมต่อ vs ปิดใช้งาน |
| Nodes by Type | `sum by (type) (node_red_nodes_connected_total)` | จำนวน nodes แยกตามชนิด |

### 4.4 Message Processing

| Panel | Metric | ความหมาย |
|-------|--------|----------|
| Message Processing Latency | `histogram_quantile(0.5/0.95/0.99, rate(node_red_msg_processing_time_seconds_bucket[5m]))` | เวลาประมวลผลข้อความ (p50/p95/p99) |
| Message Processing Rate | `rate(node_red_msg_processing_time_seconds_count[5m])` | อัตราการประมวลผล (msg/sec) |

### 4.5 Runtime Info

| Panel | Metric | ความหมาย |
|-------|--------|----------|
| Node-RED Version | `node_red_info{job="nodered"}` | เวอร์ชัน Node-RED |
| Active Flows | `node_red_flows_total` | จำนวน flows ที่ทำงานอยู่ |
| Active Nodes | `node_red_nodes_total` | จำนวน nodes ที่ทำงานอยู่ |
| Memory Usage | `process_resident_memory_bytes{job="nodered"}` | หน่วยความจำที่ใช้ |

---

## 5. PromQL Reference

```bash
# สถานะ Node-RED
up{job="nodered"}

# จำนวน nodes/flows
node_red_nodes_total{job="nodered"}
node_red_flows_total{job="nodered"}

# อัตราข้อความ (rate ต่อ 5 นาที)
rate(node_red_msg_received_total{job="nodered"}[5m])
rate(node_red_msg_sent_total{job="nodered"}[5m])

# จำนวนข้อความสะสม
node_red_msg_received_total{job="nodered"}
node_red_msg_sent_total{job="nodered"}

# เวลาประมวลผล (histogram quantile)
histogram_quantile(0.95, rate(node_red_msg_processing_time_seconds_bucket{job="nodered"}[5m]))

# Uptime
process_uptime_seconds{job="nodered"}

# Memory
process_resident_memory_bytes{job="nodered"}

# Nodes แยกตาม type
sum by (type) (node_red_nodes_connected_total{job="nodered"})

# Nodes แยกตาม state
node_red_nodes_connected_total{job="nodered"}
node_red_nodes_disabled_total{job="nodered"}
```

---

## 6. ไฟล์ที่เกี่ยวข้อง

| ไฟล์ | บทบาท |
|------|-------|
| `monitoring/grafana/dashboards/nodered/nodered-dashboard.json` | Dashboard JSON (15 panels) |
| `monitoring/grafana/provisioning/dashboards/nodered.yml` | Provisioning config สำหรับ Grafana |
| `monitoring/prometheus.yml` | เพิ่ม scrape target `nodered` (line 42-45) |

---

## 7. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| Dashboard แสดง "No data" | 1) ยังไม่ได้ติดตั้ง `node-red-contrib-prometheus` 2) ยังไม่ได้ Deploy flow 3) Prometheus target DOWN |
| Prometheus target แสดง DOWN | เช็ค `http://localhost:1881/metrics` จาก host; ถ้าใช้ Docker อาจต้องเปลี่ยน target เป็น IP จริง |
| `host.docker.internal` ใช้ไม่ได้ | Linux: เพิ่ม `extra_hosts: ["host.docker.internal:host-gateway"]` ใน docker-compose.yml |
| Dashboard ไม่โผล่ใน Grafana | เช็ค provisioning path; restart grafana; เช็ค logs `docker compose logs grafana` |
| Metrics บางตัวหาย | Plugin version ต่างกัน metrics ชื่อต่างกัน; เช็ค `curl localhost:1881/metrics | grep node_red` |
| Memory แสดง 0 | บาง plugin ไม่ export `process_resident_memory_bytes`; อาจต้องใช้ `node_exporter` แทน |

---

## 8. ตัวอย่าง curl

```bash
# ตรวจสอบ metrics endpoint
curl -s http://localhost:1881/metrics | head -20

# ตรวจสอบ Prometheus target
curl -s http://localhost:9090/api/v1/targets | jq '.data.activeTargets[] | select(.labels.job=="nodered")'

# ตรวจสอบ Grafana datasource
curl -s -u admin:admin http://localhost:3000/api/datasources | jq '.[] | select(.uid=="prometheus")'
```

---

## 9. เพิ่มเติม: ถ้า Node-RED อยู่ใน Docker

ถ้าต้องการรัน Node-RED ใน Docker ให้เพิ่มใน `docker-compose.yml`:

```yaml
nodered:
  image: nodered/node-red:latest
  ports:
    - "1881:1880"
  volumes:
    - app-nodered-data:/data
  restart: unless-stopped

nodered-exporter:
  # ไม่ต้องใช้ separate exporter ถ้าใช้ node-red-contrib-prometheus
  # แต่ถ้าต้องการ node_exporter สำหรับ OS metrics:
  image: prom/node-exporter:latest
  expose:
    - 9100
  restart: unless-stopped
```

แล้วเพิ่มใน `prometheus.yml`:
```yaml
- job_name: "nodered-node-exporter"
  static_configs:
    - targets: ["nodered-exporter:9100"]
```

---

## 10. สรุปคำสั่งทั้งหมด

```bash
# 1. ติดตั้ง plugin ใน Node-RED (ผ่าน UI หรือ CLI)
cd ~/.node-red && npm install node-red-contrib-prometheus

# 2. Deploy flow ใน Node-RED (อย่างน้อย 1 flow)

# 3. Restart Prometheus + Grafana
docker compose restart prometheus grafana

# 4. ตรวจสอบ
curl http://localhost:1881/metrics        # metrics endpoint
curl http://localhost:9090/targets        # Prometheus target
# เปิด http://localhost:3000             # Grafana dashboard
```
