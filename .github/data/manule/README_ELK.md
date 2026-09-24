# คู่มือการใช้งาน ELK Stack (Elasticsearch + Logstash + Kibana)

โปรเจกต์ `icmongolang` ใช้ **ELK Stack** เก็บ/ค้นหา/แสดงผล logs:
- **Elasticsearch** — จัดเก็บ + index logs (single-node)
- **Logstash** — รับ logs (TCP/Beats) แล้วเขียนเข้า Elasticsearch
- **Kibana** — UI ค้นหาและ visualize

**สถาปัตยกรรม:**

```
┌─────────────┐   tcp:5000 (host:5001)   ┌───────────┐
│ Go app      │ ───────────────────────▶ │ Logstash  │
│ / Beats     │ ── beats:5044 ─────────▶ │ :5000/5044│
└─────────────┘                          └─────┬─────┘
                                               │ output → index "logstash-YYYY.MM.dd"
                                               ▼
┌──────────────┐  health:9200   ┌──────────────┐   ┌──────────────┐
│ Elasticsearch│ ◀───────────── │   Kibana     │──▶│  browser     │
│ :9200        │                │ :5601        │   │ :5601        │
└──────────────┘                └──────────────┘   └──────────────┘
```

---

## 1. การตั้งค่า (docker-compose.yml)

### 1.1 Elasticsearch (:137-154)
```yaml
elasticsearch:
  image: docker.elastic.co/elasticsearch/elasticsearch:8.13.4
  environment:
    - discovery.type=single-node
    - xpack.security.enabled=false     # ปิด security → เข้าถึงตรงได้เลย
    - network.host=0.0.0.0
    - ES_JAVA_OPTS=-Xms512m -Xmx512m   # 512MB heap
    - ingest.geoip.downloader.enabled=false
  ports:
    - 9200:9200
  volumes:
    - app-elasticsearch-data:/usr/share/elasticsearch/data
  healthcheck:
    test: ["CMD-SHELL", "curl -sf http://localhost:9200/_cluster/health?wait_for_status=yellow&timeout=5s >/dev/null || exit 1"]
```

### 1.2 Logstash (:169-182)
```yaml
logstash:
  image: docker.elastic.co/logstash/logstash:8.13.4
  environment:
    - LS_JAVA_OPTS=-Xms256m -Xmx256m
  ports:
    - 5001:5000    # TCP input (host:5001)
    - 5044:5044    # Beats input (Filebeat)
    - 9600:9600    # Logstash monitoring API
  volumes:
    - ./elk/logstash/pipeline:/usr/share/logstash/pipeline
  depends_on:
    elasticsearch:
      condition: service_healthy
```

### 1.3 Kibana (:184-194)
```yaml
kibana:
  image: docker.elastic.co/kibana/kibana:8.13.4
  environment:
    - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
    - SERVER_NAME=kibana
  ports:
    - 5601:5601
  depends_on:
    elasticsearch:
      condition: service_healthy
```

### 1.4 เริ่มใช้งาน (Script `elk/elk.ps1`)
```powershell
.\elk\elk.ps1 up      # ขึ้น logstash + kibana (พร้อมตั้ง vm.max_map_count)
.\elk\elk.ps1 down    # หยุด
.\elk\elk.ps1 logs    # ดู logs
```
หรือตรง ๆ: `docker compose up -d logstash kibana`

> **สคริปต์จะตั้ง `vm.max_map_count=262144` ใน docker-desktop WSL** — จำเป็นสำหรับ Elasticsearch บน Docker Desktop

---

## 2. Logstash Pipeline (`elk/logstash/pipeline/logstash.conf`)

```ruby
input {
  tcp {
    port => 5000
  }
  beats {
    port => 5044
  }
}

filter {
}

output {
  elasticsearch {
    hosts => ["http://elasticsearch:9200"]
    index => "logstash-%{+YYYY.MM.dd}"
  }
}
```

- รับ input 2 ช่อง: **TCP :5000** และ **Beats :5044** (Filebeat)
- filter ว่าง → ส่งข้อมูลดิบ
- output → index รายวัน `logstash-2026.08.14`

### ทดสอบส่ง log
```powershell
# TCP (host:5001 → container 5000)
"hello icmon" | nc localhost 5001
# หรือผ่าน PowerShell
$tcp = New-Object System.Net.Sockets.TcpClient("localhost", 5001)
$stream = $tcp.GetStream()
$bytes = [Text.Encoding]::UTF8.GetBytes("hello icmon`n")
$stream.Write($bytes, 0, $bytes.Length)
$tcp.Close()
```

### ตรวจสอบว่า log เข้า ES แล้ว
```bash
curl -s "http://localhost:9200/logstash-*/_count"
curl -s "http://localhost:9200/logstash-*/_search?size=1&pretty"
```

---

## 3. Kibana

### 3.1 หน้าแรก
- URL: **http://localhost:5601**
- ไม่ต้อง login (ปิด xpack.security)

### 3.2 สร้าง Data View (ขั้นแรกสุด)
1. Menu → **Stack Management → Data Views**
2. **Create data view** → `logstash-*` → Save
3. ไป **Discover** เลือก data view `logstash-*` → ดู logs ได้

### 3.3 ใช้งาน
| หน้า | ใช้ทำ |
|------|-------|
| Discover | ค้นหา logs แบบ full-text + filter ตาม field |
| Analytics → Dashboard | สร้าง visualization/dashboard |
| Alerts | สร้าง alert (เช่น error เกิน threshold) |

---

## 4. Edge Cases & ข้อควรระวัง

1. **Elasticsearch ต้องการ vm.max_map_count สูง** — ใช้ `elk.ps1` หรือตั้ง WSL เอง ไม่งั้น container crash
2. **logstash TCP input ใช้ port 5000 ใน container → host 5001** — อย่าพิมพ์ nc ไป 5000
3. **index เป็นรายวัน** — ไฟล์ index pattern `logstash-*` ครอบทุกวัน
4. **ไม่มี TLS/auth** — เปิด `xpack.security.enabled=false`; เหมาะกับ dev เท่านั้น
5. **memory จำกัด** — ES 512MB + Logstash 256MB (ดูแล total memory ของ Docker)
6. **filter ว่าง** — logs เข้าแบบ raw; ถ้าต้องการ parse JSON/field ต้องเพิ่ม filter
7. **Kibana data view ต้องสร้างเองครั้งแรก** — ไม่ได้ provision อัตโนมัติ (ต่างจาก Grafana)
8. **Elasticsearch ยังเป็น data store ของ `pkg/elasticsearch` ของ app ด้วย** — อย่าลบ container โดยไม่ได้ดูว่า app ใช้อยู่

---

## 5. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| ES crash ตอน start | `vm.max_map_count` ต่ำ → `elk.ps1 up` หรือ WSL sysctl |
| Kibana ขึ้นช้ามาก | รอ ES healthy ก่อน; ดู `docker compose logs kibana` |
| `_count` = 0 ใน ES | ยังไม่มี log ส่งเข้า Logstash → ทดสอบด้วย nc :5001 |
| log เข้า TCP แต่ไม่ออก index | ตรวจ `docker compose logs logstash`; pipeline error |
| OOM | ลด `ES_JAVA_OPTS`/`LS_JAVA_OPTS` หรือเพิ่ม memory Docker |
| Discover ว่าง | ยังไม่สร้าง data view `logstash-*` |
