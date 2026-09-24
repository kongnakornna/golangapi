# คู่มือการใช้งาน Grafana (Dashboard)

Grafana เป็น dashboard UI สำหรับแสดงผล metrics จาก **Prometheus** ของโปรเจกต์ `icmongolang` — มีการ provisioning อัตโนมัติ (datasource + dashboards) ผ่าน `monitoring/grafana/`

**สถาปัตยกรรม:**

```
┌──────────────────┐   PromQL query   ┌──────────────┐   browser
│  Prometheus :9090│ ◀─────────────── │ Grafana :3000│──────────▶ http://localhost:3000
└──────────────────┘                  └──────────────┘
```

---

## 1. การตั้งค่า (docker-compose.yml:278-292)

```yaml
grafana:
  image: grafana/grafana:latest
  ports:
    - 3000:3000
  environment:
    - GF_SECURITY_ADMIN_USER=admin
    - GF_SECURITY_ADMIN_PASSWORD=admin
    - GF_USERS_ALLOW_SIGN_UP=false
  volumes:
    - ./monitoring/grafana/provisioning:/etc/grafana/provisioning:ro
    - ./monitoring/grafana/dashboards:/var/lib/grafana/dashboards:ro
    - app-grafana-data:/var/lib/grafana
  depends_on:
    - prometheus
```

### รัน
```bash
docker compose up -d grafana
```

### เข้าใช้งาน
- URL: **http://localhost:3000**
- User: `admin` / Password: `admin` (จาก env)
- ลงทะเบียนถูกปิด (`GF_USERS_ALLOW_SIGN_UP=false`)

---

## 2. Provisioning (สร้างอัตโนมัติ)

### 2.1 Data Source — `monitoring/grafana/provisioning/datasources/prometheus.yml`
```yaml
apiVersion: 1
datasources:
  - name: Prometheus
    uid: prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
    editable: true
```

### 2.2 Dashboards — `monitoring/grafana/provisioning/dashboards/*.yml`
แต่ละไฟล์ชี้ไปที่ JSON ใน `monitoring/grafana/dashboards/`:

| Provision ไฟล์ | Dashboard JSON | ครอบคลุม |
|----------------|----------------|----------|
| `golang.yml` | `golang/golang-dashboard.json` | Go backend (`http_requests_total`, duration, goroutines, `/apimetric`) |
| `websocket.yml` | `websocket/websocket-dashboard.json` | WebSocket (clients, rooms, topics, messages) |
| `redis.yml` | `redis/redis-dashboard.json` | Redis |
| `postgres.yml` | `postgres/postgres-dashboard.json` | PostgreSQL |
| `kafka.yml` | `kafka/kafka-dashboard.json` | Kafka |
| `mqtt.yml` | `mqtt/mqtt-dashboard.json` | MQTT broker |
| `influxdb.yml` | `influxdb/influxdb-dashboard.json` | InfluxDB |
| `elasticsearch.yml` | `elasticsearch/elasticsearch-dashboard.json` | Elasticsearch |
| `ollama.yml` | `ollama/ollama-dashboard.json` | Ollama |
| `nodered.yml` | `nodered/nodered-dashboard.json` | Node-RED |

ตัวอย่าง (`provisioning/dashboards/golang.yml`):
```yaml
apiVersion: 1
providers:
  - name: golang
    orgId: 1
    folder: go
    type: file
    disableDeletion: false
    editable: true
    options:
      path: /var/lib/grafana/dashboards/golang
```

---

## 3. การใช้งาน

### 3.1 ดู dashboard
- Sidebar → **Dashboards** → folder ตามชื่อ (go, websocket, redis, ...)
- เปลี่ยนช่วงเวลา (บนขวา) เพื่อดูย้อนหลัง
- dashboard ใช้ datasource `Prometheus` (uid `prometheus`)

### 3.2 Query ด้วยตัวเอง
- **Explore** (sidebar) → เลือก datasource Prometheus → พิมพ์ PromQL
- ตัวอย่าง: `sum(rate(http_requests_total[5m]))`

### 3.3 สร้าง/แก้ dashboard
- dashboards ที่ provisioning เป็น **editable: true** — แก้ได้ แต่การแก้จะ**ถูก overwrite เมื่อ restart container** (อ่านจากไฟล์ JSON)
- ถ้าต้องการคงการแก้ ต้องแก้ไฟล์ JSON ใน `monitoring/grafana/dashboards/`

---

## 4. Edge Cases & ข้อควรระวัง

1. **admin/admin ตั้งผ่าน env** — เปลี่ยนใน docker-compose หรือใช้ volume data
2. **dashboard ถูก overwrite เมื่อ restart** — แก้ที่ไฟล์ JSON ต้นทางเท่านั้น
3. **ขึ้นกับ Prometheus** (`depends_on: prometheus`) — ถ้า Prometheus down dashboard จะไม่แสดง data
4. **provisioning mount เป็น `:ro`** — แก้ใน container ไม่ได้
5. **ยังไม่มี alert rules** — ระบบ alert ยังไม่ตั้ง (ต้อง config เพิ่มเอง)
6. **GF_USERS_ALLOW_SIGN_UP=false** — ไม่มีหน้าลงทะเบียน
7. **`app-grafana-data` volume** — เก็บ user/settings ที่ไม่ใช่ provisioning

---

## 5. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| login ไม่ผ่าน | user/pass ดูจาก `GF_SECURITY_ADMIN_*`; ถ้าเปลี่ยนไปแล้วลืม → ลบ volume `app-grafana-data` (หรือ reset admin) |
| dashboard ไม่โผล่ | ตรวจ provisioning path / JSON ซ้ำกัน; restart grafana |
| "Datasource not found" | ตรวจ `datasources/prometheus.yml` ชื่อ/uid ตรงกับ dashboard |
| ข้อมูลว่างใน panel | Prometheus target down (ดู `README_Prometheus.md`) หรือ labels ผิด |
| dashboard ไม่ขึ้น auto | อาจเป็นกรณี `uid` ซ้ำใน JSON → แก้ uid หรือลบ grafana volume |
