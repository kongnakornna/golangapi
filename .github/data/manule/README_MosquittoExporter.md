# คู่มือการใช้งาน Mosquitto Exporter

`mosquitto-exporter` เป็น Prometheus exporter ที่เชื่อมกับ MQTT broker (Mosquitto) ผ่าน **MQTT protocol** แล้วอ่าน metrics จาก topic `$SYS/#` — เปิดเป็น `/metrics` สำหรับ Prometheus (ภาพรวม exporters ดู `README_Exporters.md`)

**Role ในระบบ:** monitor สถานะของ MQTT broker (clients, messages, bytes) ที่ใช้กับระบบ IoT

```
┌───────────────┐   connect tcp://mqtt:1883   ┌───────────────────┐   ┌────────────┐
│  Mosquitto    │ ◀─────────────────────────  │ mosquitto-exporter│ ─▶│ Prometheus │
│  broker :1883 │   subscribe $SYS/#          │ :9234             │   │ :9090      │
└───────────────┘                             └───────────────────┘   └─────┬──────┘
                                                                             ▼
                                                                        Grafana
```

---

## 1. การตั้งค่า (docker-compose.yml:205-212)

```yaml
mosquitto-exporter:
  image: sapcc/mosquitto-exporter:latest
  command: ["--endpoint=tcp://mqtt:1883", "--bind-address=0.0.0.0:9234"]
  expose:
    - 9234
  depends_on:
    - mqtt
  restart: unless-stopped
```

**Flags:**
- `--endpoint` — MQTT endpoint (`tcp://mqtt:1883` — port 1883 ของ broker ภายใน network)
- `--bind-address` — ที่อยู่ที่ exporter ฟัง (`0.0.0.0:9234`)

---

## 2. รันและตรวจสอบ

```bash
docker compose up -d mosquitto-exporter

# ดู metrics
docker compose exec prometheus wget -qO- http://mosquitto-exporter:9234/metrics | head -40
```

Prometheus scrape job (`monitoring/prometheus.yml:18-20`):
```yaml
- job_name: "mosquitto-exporter"
  static_configs:
    - targets: ["mosquitto-exporter:9234"]
```

---

## 3. Metrics ที่สำคัญ

| Metric | ความหมาย |
|--------|----------|
| `broker_clients_connected` | clients ที่เชื่อมอยู่ |
| `broker_clients_disconnected` | clients ที่ disconnect |
| `broker_clients_total` | clients ทั้งหมด (connected + persistent) |
| `broker_clients_maximum` | clients สูงสุดพร้อมกัน |
| `broker_clients_expired` | client session ที่หมดอายุ |
| `broker_messages_received` | จำนวนข้อความที่ broker รับ (cumulative) |
| `broker_messages_sent` | จำนวนข้อความที่ broker ส่ง |
| `broker_messages_publish_received` | PUBLISH ที่รับ |
| `broker_messages_publish_sent` | PUBLISH ที่ส่ง |
| `broker_messages_retained` | retained messages ค้างอยู่ |
| `broker_messages_dropped` | ข้อความที่ drop |
| `broker_bytes_received` / `broker_bytes_sent` | ปริมาณข้อมูลเข้า/ออก (cumulative) |
| `broker_subscriptions` | จำนวน subscriptions ทั้งหมด |
| `broker_topics` | จำนวน topics ที่มี active subscriber |
| `broker_messages_stored` | ข้อความใน queue (persistent) |

> metrics มาจาก **`$SYS/broker/...`** — ต้องเปิด `sys_interval` > 0 ใน mosquitto config (ดู `mqtt/mosquitto.conf`) ไม่งั้นจะไม่มี data

---

## 4. Query ตัวอย่าง

```promql
# clients ที่เชื่อมอยู่ตอนนี้
broker_clients_connected

# throughput messages/s
rate(broker_messages_received[5m])

# bytes/s ใน/ออก
rate(broker_bytes_received[5m])
rate(broker_bytes_sent[5m])

# subscriptions ที่ active
sum(broker_subscriptions)
```

---

## 5. Edge Cases & ข้อควรระวัง

1. **ขึ้นกับ `sys_interval` ของ broker** — ต้องเปิด MQTT `$SYS` ใน `mosquitto.conf` ไม่งั้น exporter เห็นค่าศูนย์/ว่าง
2. **Exporter เป็น MQTT client จริง** — นับเป็น 1 ใน `broker_clients_connected` ด้วย
3. **port 1883 (internal)** — ไม่ใช่ 1885 ที่ host map ไว้ (ดู README_MQTT.md)
4. **ไม่ส่ง username/password** — ถ้า broker เปิด auth ต้องเพิ่ม flag `--username/--password`
5. **ค่า cumulative เป็น counter** — ใช้ `rate()`/`increase()` ไม่ใช่ค่าดิบ
6. **expose 9234** — เข้าถึงได้จาก docker network เท่านั้น

---

## 6. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| metric เป็น 0/ไม่มี data | ตรวจ `sys_interval` ใน mosquitto.conf |
| connect fail ตอน start | mqtt container ยังไม่ up หรือ broker port ผิด |
| target down | exporter ไม่รัน → `docker compose logs mosquitto-exporter` |
| ค่าเพิ่มไม่หยุด | เป็น counter — ใช้ `rate()` |
