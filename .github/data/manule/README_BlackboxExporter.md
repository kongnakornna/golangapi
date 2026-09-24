# คู่มือการใช้งาน Blackbox Exporter

`blackbox-exporter` เป็น Prometheus exporter ที่ **probe (ตรวจสอบ) ความพร้อมใช้งาน** ของ endpoint ภายนอกผ่าน TCP/HTTP/ICMP/DNS — ในโปรเจกต์ `icmongolang` ใช้ตรวจ `/health` และ `/apimetric` ของ Go backend (ภาพรวม exporters ดู `README_Exporters.md`)

```
                        probe                        ┌──────────────────────┐
┌──────────────┐  /probe?module=http_2xx ──────────▶ │ blackbox-exporter    │
│  Prometheus  │  &target=...             ──────────▶ │ :9115                │
│  :9090       │ ◀─────────── metrics ─────────────── │  │                  │
└──────────────┘                                     └──┼──────────────────┘
                                                        │ HTTP request
                                                        ▼
                                               http://backend:5000/health
                                               http://backend:5000/apimetric
```

**ข้อดี:** วัดความพร้อมใช้งานจริงของ service (availability/latency) ไม่ใช่แค่ "exporter ต่อได้"

---

## 1. การตั้งค่า (docker-compose.yml:131-135)

```yaml
blackbox-exporter:
  image: prom/blackbox-exporter:latest
  expose:
    - 9115
  restart: unless-stopped
```

> **⚠️ ไม่ได้ mount config** (`-c /etc/blackbox_exporter/config.yml`) — ใช้ **default module** ของ image: `http_2xx`, `icmp`, `tcp_connect`, `dns_tcp`, `dns_udp`, `ssh_banner`, `smtp_banner`, ... ถ้าต้องการ module กำหนดเองต้อง mount config file

---

## 2. การใช้งานร่วมกับ Prometheus

`monitoring/prometheus.yml:42-56`:
```yaml
- job_name: "backend-blackbox"
  metrics_path: /probe
  params:
    module: [http_2xx]
  static_configs:
    - targets:
        - "http://backend:5000/health"
        - "http://backend:5000/apimetric"
  relabel_configs:
    - source_labels: [__address__]
      target_label: __param_target
    - source_labels: [__param_target]
      target_label: instance
    - target_label: __address__
      replacement: blackbox-exporter:9115
```

**Relabel อธิบาย:**
1. `__address__` (เดิมคือ target URL) → เอาไปเป็น `__param_target` (URL ที่จะ probe)
2. `__param_target` → copy เป็น label `instance` (รู้ว่ากำลัง probe ตัวไหน)
3. `__address__` ถูกเขียนทับเป็น `blackbox-exporter:9115` — Prometheus เรียก blackbox แทน target จริง

---

## 3. Probe ตรง ๆ

```bash
docker compose exec prometheus wget -qO- \
  'http://blackbox-exporter:9115/probe?module=http_2xx&target=http://backend:5000/health' | head -30
```

**Test metric:** `probe_success 1` หมายถึง HTTP 200

---

## 4. Metrics ที่สำคัญ (module http_2xx)

| Metric | ความหมาย |
|--------|----------|
| `probe_success` | 1 = probe สำเร็จ (HTTP 2xx) |
| `probe_duration_seconds` | เวลา probe ทั้งหมด (latency) |
| `probe_http_status_code` | HTTP status ที่ได้ |
| `probe_http_duration_seconds{phase="connect"}` | เวลา TCP connect |
| `probe_http_duration_seconds{phase="tls"}` | เวลา TLS handshake |
| `probe_http_duration_seconds{phase="processing"}` | เวลา server processing |
| `probe_http_duration_seconds{phase="transfer"}` | เวลาโอน body |
| `probe_http_version` | HTTP version |
| `probe_ssl_earliest_cert_expiry` | วันหมดอายุ cert (สำคัญถ้ามี TLS) |
| `probe_http_redirects` | จำนวน redirects |

---

## 5. Query ตัวอย่าง

```promql
# availability (%) ของ /health ใน 24 ชม.
100 * avg_over_time(probe_success{instance="http://backend:5000/health"}[24h])

# latency p95
histogram_quantile(0.95, sum by (le, instance) (rate(probe_http_duration_seconds_bucket[5m])))

# error count 5 นาที
sum(increase(probe_success{job="backend-blackbox"} == 0[5m]))

# SSL cert เหลือกี่วัน
(probe_ssl_earliest_cert_expiry - time()) / 86400
```

---

## 6. Edge Cases & ข้อควรระวัง

1. **module default จาก image** — ไม่มี custom config (เช่น timeout/custom headers)
2. **`probe_success` ดูเป็น step** — ใช้ `avg_over_time` เพื่อดู availability
3. **instance label = URL ที่ probe** — ใช้ filter นี้ใน Grafana
4. **ระวัง DNS fail** — target `http://backend:5000/...` ใช้ hostname `backend` ซึ่ง compose ตั้งชื่อ service เป็น `api` (ดู README_Prometheus.md §2.2) → probe อาจ fail ด้วย DNS
5. **probe ทุก 15s** (scrape_interval) — availability ขั้นต่ำละเอียดแค่ 15s
6. **ถ้า backend ต้องการ auth** — blackbox http_2xx default ไม่ส่ง Bearer; ต้อง custom config + `bearer_token`

---

## 7. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| `probe_success == 0` | endpoint ตอบ error หรือ DNS ผิด |
| target down ใน Prometheus | blackbox container ไม่รัน |
| `probe_dns_lookup` ช้า | ดู DNS resolver ของ container |
| query ว่าง | ตรวจชื่อ job (`backend-blackbox`) และ instance label |
