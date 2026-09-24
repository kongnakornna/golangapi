# คู่มือการใช้งาน Ollama (+ Ollama Exporter)

**Ollama** เป็น LLM server รัน local ในโปรเจกต์ `icmongolang` — ให้บริการโมเดล (default `qwen3:4b`) ให้กับ backend และ `ai` service ผ่าน API `:11434` (OpenAI-compatible `/v1`)

**Ollama Exporter** เป็น Prometheus exporter ตัวเล็กที่เขียนเองด้วย Go (`cmd/ollama-exporter`) — ดึง metrics จาก Ollama API แล้วเปิด `/metrics:9309`

**สถาปัตยกรรม:**

```
┌─────────────┐  LLM_BASE_URL=http://ollama:11434/v1   ┌──────────────────┐
│ Go backend  │ ─────────────────────────────────────▶ │                  │
│             │  OLLAMA_BASE_URL=http://ollama:11434   │      Ollama      │
│  ai service │ ─────────────────────────────────────▶ │  :11434          │
│  :8001       │                                       │  (local LLM)     │
└─────────────┘                                       └────────┬─────────┘
                                                              │ /api/version, /api/ps, /api/tags
                                              scrape 15s      ▼
                                        ┌─────────────┐  ┌───────────────┐
                                        │ollama-      │─▶│ Prometheus    │
                                        │exporter:9309│  │ :9090         │
                                        └─────────────┘  └───────┬───────┘
                                                                  ▼
                                                             Grafana
```

---

## 1. Ollama Server

### docker-compose.yml:234-248
```yaml
ollama:
  image: ollama/ollama:latest
  ports:
    - 11434:11434
  environment:
    - OLLAMA_KEEP_ALIVE=30m        # โมเดลค้างในหน่วยความจำ 30 นาทีหลังใช้
  volumes:
    - app-ollama-data:/root/.ollama  # เก็บโมเดลที่ download
  healthcheck:
    test: ["CMD", "ollama", "list"]
    interval: 30s
    timeout: 10s
    retries: 5
    start_period: 15s
```

### รันและจัดการโมเดล
```bash
docker compose up -d ollama

# download โมเดลเริ่มต้น (ถ้ายังไม่มี)
docker compose exec ollama ollama pull qwen3:4b

# ดูโมเดลที่มี
docker compose exec ollama ollama list

# ทดสอบ API
curl -s http://localhost:11434/api/version
# {"version":"0.x.x"}

# ทดสอบ generate
curl -s http://localhost:11434/api/generate -d '{"model":"qwen3:4b","prompt":"สวัสดี","stream":false}'
```

### OpenAI-compatible endpoint
`http://ollama:11434/v1/chat/completions` — backend ใช้ผ่าน env `LLM_BASE_URL=http://ollama:11434/v1`

---

## 2. Ollama Exporter (`cmd/ollama-exporter/main.go`)

### 2.1 docker-compose.yml:250-261
```yaml
ollama-exporter:
  build:
    context: .
    dockerfile: Dockerfile.ollama-exporter
  environment:
    - OLLAMA_URL=http://ollama:11434
    - LISTEN_ADDR=:9309
  expose: [9309]
  depends_on: [ollama]
```

### 2.2 ตั้งค่า (env/flag)
| env | flag | default |
|-----|------|---------|
| `OLLAMA_URL` | `--ollama.url` | `http://localhost:11434` |
| `LISTEN_ADDR` | `--listen` | `:9309` |
| (ไม่มี env) | `--interval` | `15s` (scrape interval) |

### 2.3 วิธีทำงาน (`main.go`)
- loop ทุก `interval` (15s) เรียก 3 endpoints:
  - `/api/version` → `ollama_up`, `ollama_version_info`, `ollama_scrape_duration_seconds`
  - `/api/ps` → โมเดลที่โหลดอยู่: `ollama_loaded_models`, `ollama_model_size_bytes`, `ollama_model_vram_bytes`
  - `/api/tags` → โมเดลใน catalog: `ollama_catalog_models`, `ollama_model_modified_seconds`
- ตัวใด fail → `ollama_up = 0` (และไม่ตั้ง metric นั้น)
- metric ถูก `Reset()` ก่อน scrape ทุกครั้ง → แสดงเฉพาะค่าล่าสุด

### 2.4 Metrics

| Metric | Labels | ความหมาย |
|--------|--------|----------|
| `ollama_up` | — | 1=Ollama reachable, 0=ไม่ |
| `ollama_scrape_duration_seconds` | — | เวลา scrape ล่าสุด |
| `ollama_version_info` | version | เวอร์ชัน (1) |
| `ollama_loaded_models` | model, family, parameter_size, quantization | โมเดลที่โหลดใน VRAM (1) |
| `ollama_model_size_bytes` | model | ขนาดโมเดลบนดิสก์ |
| `ollama_model_vram_bytes` | model | VRAM ที่ใช้ (size_vram) |
| `ollama_catalog_models` | model | โมเดลทั้งหมดที่ pull แล้ว (1) |
| `ollama_model_modified_seconds` | model | unix timestamp ที่แก้ล่าสุด |

### 2.5 Scrape job (`monitoring/prometheus.yml:31-33`)
```yaml
- job_name: "ollama-exporter"
  static_configs:
    - targets: ["ollama-exporter:9309"]
```

---

## 3. `ai` Service (FastAPI ฝั่ง LLM wrapper)

### docker-compose.yml:263-276
```yaml
ai:
  build: { context: ./ai }
  environment:
    - OLLAMA_BASE_URL=http://ollama:11434
    - OLLAMA_MODEL=${OLLAMA_MODEL:-qwen3:4b}
    - PORT=8001
    - AI_TIMEOUT=300s
  ports:
    - 8001:8001
  depends_on:
    ollama:
      condition: service_healthy
```
- Python service (ใน `ai/`) ทำหน้าที่เป็น wrapper/แยก business logic จาก Ollama API ตรง
- ใช้ `OLLAMA_MODEL` default `qwen3:4b`

---

## 4. Query ตัวอย่าง (Grafana/Prometheus)

```promql
# Ollama ขึ้นไหม
ollama_up

# โมเดลที่โหลดอยู่
sum(ollama_loaded_models)

# VRAM ทั้งหมดที่ใช้
sum(ollama_model_vram_bytes)

# โมเดลขนาดบนดิสก์
sum by (model) (ollama_model_size_bytes)
```

Grafana dashboard มี provisioning ไว้แล้ว (`monitoring/grafana/provisioning/dashboards/ollama.yml` → `ollama-dashboard.json`)

---

## 5. Edge Cases & ข้อควรระวัง

1. **pull โมเดลครั้งแรกต้องทำเอง** — compose ไม่ได้ pull อัตโนมัติ (ต้อง `ollama pull qwen3:4b`)
2. **OLLAMA_KEEP_ALIVE=30m** — โมเดลถูกค้างใน RAM/VRAM 30 นาทีหลังใช้งาน (ลด latency แต่กินหน่วยความจำ)
3. **ตัวแรกที่เรียก LLM จะช้า** — ต้องโหลดโมเดลเข้าหน่วยความจำ (cold start)
4. **exporter `Reset()` metric ทุก scrape** — ค่าหายถ้า scrape fail; ระวัง `avg_over_time` บน gauge แบบนี้
5. **`OLLAMA_URL` default `localhost`** — ถ้ารัน exporter นอก compose ต้องตั้ง env ให้ชี้ container
6. **โมเดล default `qwen3:4b`** — ตั้งผ่าน `.env` `OLLAMA_MODEL`
7. **ai service รอ health ของ Ollama ก่อน** (`condition: service_healthy`)

---

## 6. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| `ollama_up == 0` | Ollama down หรือ `OLLAMA_URL` ผิด |
| response ช้า / timeout | โมเดลต้องโหลดครั้งแรก; หรือ `AI_TIMEOUT` สั้นไป |
| `ollama pull` fail | เช็ค internet / disk space (โมเดลหลาย GB) |
| backend เรียก LLM ไม่ได้ | ตรวจ `LLM_BASE_URL` และ network |
| ai service ไม่ start | Ollama ยังไม่ healthy → รอ 15-30s |
| dashboard ว่าง | exporter not running / ยังไม่มีโมเดล (catalog ว่าง) |
