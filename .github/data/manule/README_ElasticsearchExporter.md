# คู่มือการใช้งาน Elasticsearch Exporter

`elasticsearch-exporter` เป็น Prometheus exporter ที่ดึง metrics จาก Elasticsearch และเปิดเป็น `/metrics` สำหรับ Prometheus (ภาพรวม exporters ทั้งหมดดู `README_Exporters.md`)

**Role ในระบบ:** monitor สภาพ Elasticsearch (cluster health, node, disk, indices, JVM) ที่ใช้เป็น data store ของ ELK stack

```
┌────────────────────┐  scrape   ┌──────────────────────┐   ┌───────────────┐
│ Elasticsearch :9200│ ◀─────── │ elasticsearch-exporter│ ─▶│ Prometheus    │
│                    │  metrics  │ :9114                │   │ :9090         │
└────────────────────┘           └──────────────────────┘   └───────┬───────┘
                                                                     ▼
                                                                Grafana :3000
```

---

## 1. การตั้งค่า (docker-compose.yml:156-167)

```yaml
elasticsearch-exporter:
  image: prometheuscommunity/elasticsearch-exporter:latest
  command:
    - "--es.uri=http://elasticsearch:9200"
    - "--es.all"
    - "--es.indices"
  expose:
    - 9114
  depends_on:
    elasticsearch:
      condition: service_healthy
  restart: unless-stopped
```

**Flags:**
- `--es.uri` — ที่อยู่ ES (ชี้ไป `elasticsearch:9200`)
- `--es.all` — เก็บ cluster + node metrics ทั้งหมด
- `--es.indices` — เก็บ metrics แยกราย index

**Health check ของ elasticsearch (yml:145-154):** `curl http://localhost:9200` → `service_healthy` → exporter รอจนกว่า ES พร้อม

---

## 2. รันและตรวจสอบ

```bash
docker compose up -d elasticsearch-exporter

# ดู metrics โดยตรง
docker compose exec prometheus wget -qO- http://elasticsearch-exporter:9114/metrics | head -40
# หรือจาก host
curl -s http://localhost:9114/metrics | head
```

Prometheus scrape job (`monitoring/prometheus.yml:27-29`):
```yaml
- job_name: "elasticsearch-exporter"
  static_configs:
    - targets: ["elasticsearch-exporter:9114"]
```

---

## 3. Metrics ที่สำคัญ

| Metric | ความหมาย |
|--------|----------|
| `elasticsearch_up` | exporter ต่อ ES ได้ไหม (1/0) |
| `elasticsearch_cluster_health_status` | 0=red, 1=yellow, 2=green |
| `elasticsearch_cluster_health_number_of_nodes` | จำนวน node |
| `elasticsearch_cluster_health_active_shards` | active shards |
| `elasticsearch_filesystem_data_available_bytes` | พื้นที่ data ที่เหลือ |
| `elasticsearch_filesystem_data_free_bytes` | พื้นที่ว่าง |
| `elasticsearch_filesystem_data_size_bytes` | ขนาดพื้นที่ทั้งหมด |
| `elasticsearch_indices_docs_count{index}` | จำนวน doc ต่อ index |
| `elasticsearch_indices_docs_deleted{index}` | doc ที่ลบ (soft delete) |
| `elasticsearch_indices_store_size_bytes{index}` | ขนาด data ต่อ index |
| `elasticsearch_jvm_memory_used_bytes{node,area}` | JVM heap/ non-heap |
| `elasticsearch_jvm_gc_collection_seconds_count{node}` | GC count |

---

## 4. Query ตัวอย่าง (Prometheus/Grafana)

```promql
# cluster สมบูรณ์ไหม (1 = green)
elasticsearch_cluster_health_status

# พื้นที่ดิสก์เหลือ %
1 - (elasticsearch_filesystem_data_available_bytes / elasticsearch_filesystem_data_size_bytes)

# docs ทั้งหมดต่อ index
sum by (index) (elasticsearch_indices_docs_count)

# JVM heap ใช้ไปเท่าไร
sum by (node) (elasticsearch_jvm_memory_used_bytes{area="heap"})
```

---

## 5. Edge Cases & ข้อควรระวัง

1. **rันเมื่อ ES healthy เท่านั้น** (`service_healthy`) — ถ้า ES down, exporter container จะไม่ start
2. **`--es.all` + `--es.indices`** อาจทำให้ metric เยอะมากถ้ามีหลาย indices — เลือกเฉพาะที่จำเป็นได้
3. **ES ในโปรเจกต์ปิด xpack security** (`xpack.security.enabled=false`) — ไม่ต้อง user/password; ถ้าเปิด security ต้องเพิ่ม `--es.username/--es.password`
4. **expose 9114** — เข้าถึงได้จาก docker network เท่านั้น
5. **scrape ผ่าน Prometheus เท่านั้น** — อย่าไปเปิด publish port ถ้าไม่จำเป็น (ผ่าน prometheus:9090 เป็นทางเข้า)

---

## 6. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| container ยังไม่ start | ES ไม่ผ่าน health → ตรวจ `docker compose ps` ของ elasticsearch |
| target down ใน Prometheus | exporter ไม่รัน หรือ network ชื่อผิด |
| metric ว่าง | ES ยังไม่พร้อมตอน scrape แรก; รอ 15s ดูใหม่ |
| `elasticsearch_up == 0` | ES ถูกหยุด หรือ `--es.uri` ผิด |
