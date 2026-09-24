# คู่มือการใช้งาน Exporters (Prometheus Exporters)

Exporters เป็น bridge ที่แปลง metrics ของระบบต่าง ๆ ให้เป็น Prometheus format — โปรเจกต์ `icmongolang` รัน exporters หลายตัวใน docker-compose และให้ Prometheus scrape

**ภาพรวม:**

| Exporter | Container | Port | Target metric | Image |
|----------|-----------|------|---------------|-------|
| Kafka | `kafka-exporter` | 9308 | kafka:9092 | `danielqsj/kafka-exporter` |
| Redis | `redis-exporter` | 9121 | redis:6379 | `oliver006/redis_exporter` |
| PostgreSQL | `postgres-exporter` | 9187 | db:5435 | `prometheuscommunity/postgres-exporter` |
| Blackbox | `blackbox-exporter` | 9115 | probe HTTP/TCP | `prom/blackbox-exporter` |
| Elasticsearch | `elasticsearch-exporter` | 9114 | elasticsearch:9200 | `prometheuscommunity/elasticsearch-exporter` |
| Mosquitto | `mosquitto-exporter` | 9234 | tcp://mqtt:1883 | `sapcc/mosquitto-exporter` |

> ทั้งหมดถูก scrape โดย `monitoring/prometheus.yml` (ดู `README_Prometheus.md`)

---

## 1. Kafka Exporter (danielqsj/kafka-exporter)

### docker-compose.yml:85-92
```yaml
kafka-exporter:
  image: danielqsj/kafka-exporter:latest
  command: ["--kafka.server=kafka:9092", "--web.listen-address=:9308"]
  expose: [9308]
  depends_on: [kafka]
```

### Metric ตัวอย่าง
- `kafka_brokers`
- `kafka_topic_partitions{topic}`
- `kafka_topic_partition_current_offset{topic,partition}`
- `kafka_topic_partition_oldest_offset`
- `kafka_consumergroup_current_offset{consumergroup,topic,partition}`
- `kafka_consumergroup_lag{consumergroup,topic,partition}` — **ใช้ monitor consumer lag**

---

## 2. Redis Exporter (oliver006/redis_exporter)

### docker-compose.yml:112-119
```yaml
redis-exporter:
  image: oliver006/redis_exporter:latest
  command: ["--redis.addr=redis:6379"]
  expose: [9121]
  depends_on: [redis]
```

### Metric ตัวอย่าง
- `redis_up`
- `redis_memory_used_bytes`
- `redis_memory_max_bytes`
- `redis_connected_clients`
- `redis_connected_slaves`
- `redis_keyspace_hits_total` / `redis_keyspace_misses_total`
- `redis_db_keys{db="0"}` / `redis_db_keys{db="1"}`

---

## 3. PostgreSQL Exporter (prometheuscommunity/postgres-exporter)

### docker-compose.yml:121-129
```yaml
postgres-exporter:
  image: prometheuscommunity/postgres-exporter:latest
  environment:
    - DATA_SOURCE_NAME=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5435/${POSTGRES_DBNAME}?sslmode=disable
  expose: [9187]
  depends_on: [db]
```

> ⚠️ ชี้ไปที่ port **5435** (port ของ db container) และใช้ env `POSTGRES_USER/PASSWORD/DBNAME`

### Metric ตัวอย่าง
- `pg_up`
- `pg_stat_database_numbackends{datname}`
- `pg_stat_database_connections{datname}`
- `pg_stat_database_tup_fetched` / `pg_stat_database_tup_inserted`
- `pg_settings_max_connections`
- `pg_locks_count{datname}`

---

## 4. Blackbox Exporter (prom/blackbox-exporter)

### docker-compose.yml:131-135
```yaml
blackbox-exporter:
  image: prom/blackbox-exporter:latest
  expose: [9115]
```

ไม่มีการ mount config — ใช้ default module (`icmp`, `tcp_connect`, `http_2xx`, ...)

### การใช้งาน
Prometheus เรียก `/probe?module=http_2xx&target=<url>`:
```promql
probe_success{job="backend-blackbox"}
probe_duration_seconds{job="backend-blackbox"}
probe_http_status_code{job="backend-blackbox"}
```
โปรเจกต์นี้ probe 2 endpoint: `http://backend:5000/health` และ `http://backend:5000/apimetric`

---

## 5. Elasticsearch Exporter (prometheuscommunity/elasticsearch-exporter)

### docker-compose.yml:156-167
```yaml
elasticsearch-exporter:
  image: prometheuscommunity/elasticsearch-exporter:latest
  command:
    - "--es.uri=http://elasticsearch:9200"
    - "--es.all"
    - "--es.indices"
  expose: [9114]
  depends_on:
    elasticsearch:
      condition: service_healthy
```

### flag ที่ใช้
- `--es.uri` — ที่อยู่ Elasticsearch
- `--es.all` — เก็บ cluster + nodes
- `--es.indices` — เก็บ per-index metrics

### Metric ตัวอย่าง
- `elasticsearch_cluster_health_status`
- `elasticsearch_cluster_health_number_of_nodes`
- `elasticsearch_filesystem_data_available_bytes` / `_free_bytes`
- `elasticsearch_indices_docs_count{index}`
- `elasticsearch_jvm_memory_used_bytes{node}`

---

## 6. Mosquitto Exporter (sapcc/mosquitto-exporter)

### docker-compose.yml:205-212
```yaml
mosquitto-exporter:
  image: sapcc/mosquitto-exporter:latest
  command: ["--endpoint=tcp://mqtt:1883", "--bind-address=0.0.0.0:9234"]
  expose: [9234]
  depends_on: [mqtt]
```

### Metric ตัวอย่าง
- `broker_clients_connected`
- `broker_clients_disconnected`
- `broker_clients_total`
- `broker_messages_received` / `broker_messages_sent`
- `broker_messages_publish_received`
- `broker_bytes_received` / `broker_bytes_sent`

---

## 7. วิธีตรวจสอบ

```bash
# exporter ทุกตัว up?
docker compose ps

# ทดสอบ scrape โดยตรง (จาก container)
docker compose exec prometheus wget -qO- http://redis-exporter:9121/metrics | head
docker compose exec prometheus wget -qO- http://postgres-exporter:9187/metrics | head
docker compose exec prometheus wget -qO- http://mosquitto-exporter:9234/metrics | head

# ผ่าน UI
open http://localhost:9090/targets   # ดูสถานะทุก target
```

---

## 8. Edge Cases & ข้อควรระวัง

1. **expose (ไม่ publish port)** — exporters เข้าถึงได้จาก docker network เท่านั้น; ต้องรันผ่าน compose ไม่ใช่แยก
2. **ขึ้นกับ service เป้าหมาย** — ใช้ `depends_on` แต่ส่วนใหญ่เป็นแค่ start order ไม่ใช่ health condition (ยกเว้น elasticsearch)
3. **postgres-exporter ใช้ `DATA_SOURCE_NAME`** — ถ้า env เปลี่ยน DB ต้อง restart
4. **kafka-exporter ติดตาม consumergroup ได้** — ใช้คู่กับ `kafka_consumergroup_lag` เพื่อ alert lag
5. **mosquitto-exporter เชื่อม `tcp://mqtt:1883`** — ถ้า broker ต้องการ auth ต้องเพิ่ม flag
6. **blackbox ไม่มี config แยก** — ถ้าต้องการ module แบบกำหนดเองต้อง mount config

---

## 9. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| target down | exporter container ขึ้นไม่ทัน → ตรวจ `docker compose logs <exporter>` |
| `postgres-exporter` ขึ้น fail | `DATA_SOURCE_NAME` env ไม่ถูกต้อง / db ยังไม่พร้อม |
| `elasticsearch-exporter` down | รอ health ของ elasticsearch (condition: service_healthy) |
| query ว่างใน Grafana | ตรวจชื่อ metric ผ่าน `/metrics` ตรง ๆ |
| kafka metric ไม่มา | เช็ค `--kafka.server` ถูกต้อง + kafka up |
