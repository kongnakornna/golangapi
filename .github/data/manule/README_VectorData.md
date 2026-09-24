# คู่มือการใช้งานโมดูล Vector Data (Vector Database)

โมดูล `internal/modules/vectordata` เป็น **ระบบ Vector Database** ของโปรเจกต์ IC-MON
ที่ออกแบบให้ทำงานได้กับ **มากกว่า 1 backend** ผ่าน abstraction layer `pkg/vectordb`
ทำให้สลับ backend ได้โดยแก้ config เท่านั้น โดยไม่ต้องแก้โค้ด:

| Provider | ชนิด | รายละเอียด |
|----------|------|-------------|
| `pgvector` | PostgreSQL + pgvector extension | เก็บเวกเตอร์ในคอลัมน์ `vector` + HNSW index (ค่า default) |
| `elasticsearch` | Elasticsearch `dense_vector` | ใช้ index ที่รองรับ KNN cosine similarity |

โมดูลทำงานร่วมกับ **LLM Embedding Client** (`pkg/llm`) เพื่อแปลงข้อความเป็นเวกเตอร์ก่อนจัดเก็บ/ค้นหา

---

## 1. สถาปัตยกรรม (Architecture)

```
Client / Postman / Frontend
        │ HTTP (chi router)
        ▼
┌─────────────────────────────────────────────────┐
│ delivery/http (internal/modules/vectordata)      │
│   handler.go  → รับ HTTP request / ส่ง response   │
│   routes.go   → ลงทะเบียน route ใต้ /api/vectordata │
└─────────────────────────────────────────────────┘
        │
        ▼
┌─────────────────────────────────────────────────┐
│ usecase                                          │
│   usecase.go  → ประสานงาน LLM + VectorDB          │
└─────────────────────────────────────────────────┘
        │
        ▼
┌─────────────────────────────────────────────┐
│ pkg/vectordb (abstraction layer)             │
│   VectorDB interface                        │
│   ┌───────────────┐   ┌──────────────────┐  │
│   │ es_backend.go │   │ pg_backend.go    │  │
│   │ (Elasticsearch)│  │ (pgvector)       │  │
│   └───────────────┘   └──────────────────┘  │
└─────────────────────────────────────────────┘
        │                    │
        ▼                    ▼
  Elasticsearch (9200)   PostgreSQL (5435)
                        + pgvector extension
```

### 1.1 ตารางไฟล์

| ไฟล์ | บทบาท |
|------|-------|
| `pkg/vectordb/vectordb.go` | Interface `VectorDB` + factory `New()` เลือก backend ตาม config |
| `pkg/vectordb/es_backend.go` | Backend Elasticsearch (wrap `pkg/elasticsearch.Client`) |
| `pkg/vectordb/pg_backend.go` | Backend PostgreSQL + pgvector (raw SQL ผ่าน GORM) |
| `internal/modules/vectordata/presenter/presenter.go` | DTO — โครงสร้าง Request / Response |
| `internal/modules/vectordata/usecase/usecase.go` | Business logic — embed ผ่าน LLM + สั่งงาน VectorDB + sample docs สำหรับ seed |
| `internal/modules/vectordata/delivery/http/handler.go` | HTTP handler + ตรวจสอบข้อมูล + error mapping |
| `internal/modules/vectordata/delivery/http/routes.go` | ลงทะเบียน route พร้อม middleware |
| `cmd/vectordata.go` | CLI command `vectordata seed` สำหรับ seed ข้อมูลตัวอย่าง |

### 1.2 จุดที่โมดูลถูก connect (wiring)

- **`internal/server/handlers.go`** — สร้าง `vdb` ผ่าน `vectordb.New(&cfg.VectorDB, esClient, db, logger)` ถ้าสำเร็จจะลง route ผ่าน `MapVectorDataRoutes(apiRouter, vdHandler, mw)` พร้อม log provider ที่ใช้ ถ้าล้มเหลวข้ามไปพร้อม log error
- **`cmd/migrate.go`** — ถ้า `vectorDb.provider == "pgvector"` จะเรียก `EnsureIndex` เพื่อสร้าง extension `vector` + ตาราง + HNSW index อัตโนมัติตอน `go run main.go migrate`
- Base path ของ API ทั้งหมดคือ `/api` → endpoints นี้อยู่ใต้ `/api/vectordata/...`
- ตรวจสอบสถานะได้ที่ `/api/health` (field `vectordata`)

---

## 2. ข้อกำหนดของระบบ (Prerequisites)

| Component | รายละเอียด |
|-----------|-----------|
| PostgreSQL (pgvector) | ถ้าใช้ `pgvector` ต้องใช้ PostgreSQL ที่ติดตั้ง extension `vector` ไว้ — ใน docker-compose ใช้ image **`pgvector/pgvector:pg15`** |
| Elasticsearch | ถ้าใช้ `elasticsearch` ต้องใช้ ES ที่รองรับ `dense_vector` (>= 7.x; docker-compose ใช้ 8.13.4) |
| LLM Embedding Provider | API แบบ OpenAI-compatible ที่มี `POST {base_url}/embeddings` เช่น **Ollama** (ค่า default) |

> ⚠️ ถ้าใช้ `pgvector` บนเครื่องที่ติดตั้ง PostgreSQL เอง (ไม่ได้ผ่าน docker-compose)
> ต้องติดตั้ง pgvector extension ก่อน: `CREATE EXTENSION IF NOT EXISTS vector;`
> (คำสั่งนี้จะถูก run อัตโนมัติตอน migrate / EnsureIndex แต่ถ้า image ไม่มีตัวติดตั้งจะ error)

---

## 3. การตั้งค่า Configuration

โมดูลอ่านค่าจากบล็อก `vectorDb` และ `llm` (จาก `config/config.go`):

```yaml
llm:
  base_url: "http://localhost:11434/v1"   # Ollama /v1 (OpenAI-compatible)
  api_key: ""
  model: "nomic-embed-text"
  timeout: 30
  dims: 768

vectorDb:
  provider: "pgvector"        # "elasticsearch" | "pgvector"
  index: "vector_documents"   # ชื่อ ES index / PG table
  dims: 768                   # มิติของเวกเตอร์
```

### 3.1 ตัวแปรสภาพแวดล้อม (`.env` / `.env.example`)

ระบบ BindEnvs แปลง `mapstructure` tag เป็น environment variable อัตโนมัติ:

```env
# ---- Vector Database ----
VECTOR_DB_PROVIDER=pgvector            # elasticsearch | pgvector
VECTOR_DB_INDEX=vector_documents
VECTOR_DB_DIMS=768
```

| Config | Env | ค่า default | ความหมาย |
|--------|-----|-------------|----------|
| `vectorDb.provider` | `VECTOR_DB_PROVIDER` | `pgvector` | เลือก backend |
| `vectorDb.index` | `VECTOR_DB_INDEX` | `vector_documents` | ชื่อ index/table |
| `vectorDb.dims` | `VECTOR_DB_DIMS` | `768` | มิติของเวกเตอร์ (ต้องตรงกับ dims ของโมเดล embedding) |

> ดูค่า `llm.*` เพิ่มเติมได้ใน `README_LLM.md`

### 3.2 ตัวอย่าง: สลับไปใช้ Elasticsearch

```env
VECTOR_DB_PROVIDER=elasticsearch
VECTOR_DB_INDEX=vector_documents   # index เฉพาะของโมดูล vectordata
```

ต้องมี Elasticsearch เปิดอยู่ (docker-compose มี service `elasticsearch`) และมี `ELASTICSEARCH_ADDRESSES` ใน config

> หมายเหตุ: โมดูล vectordata (`/api/vectordata/*`) และโมดูล elasticsearch (`/api/elasticsearch/*`) ใช้ **index แยกกัน** — vectordata ใช้ `cfg.VectorDB.Index` ส่วนโมดูล elasticsearch ใช้ `cfg.Elasticsearch.Index` — การ seed/search ของทั้งสองไม่ชนกัน

---

## 4. API Reference

### 4.1 Endpoints (ทั้งหมดอยู่ใต้ `/api/vectordata`)

| Method | Path | Auth | คำอธิบาย |
|--------|------|------|----------|
| GET | `/vectordata/health` | ❌ Public | ตรวจสอบสถานะ backend (provider + up/down) |
| POST | `/vectordata/create-index` | ✅ JWT | สร้าง/ยืนยัน index หรือ table + HNSW index |
| POST | `/vectordata/documents` | ✅ JWT | สร้างเอกสาร — embed ผ่าน LLM (หรือรับ vector สำเร็จรูป) แล้วจัดเก็บ |
| GET | `/vectordata/documents` | ✅ JWT | แสดงรายการแบบมี pagination (`?limit=&offset=`) |
| GET | `/vectordata/documents/{id}` | ✅ JWT | ดูเอกสารตาม id |
| PUT | `/vectordata/documents/{id}` | ✅ JWT | อัปเดตเอกสาร (re-embed ถ้าเปลี่ยน content) |
| DELETE | `/vectordata/documents/{id}` | ✅ JWT | ลบเอกสาร |
| POST | `/vectordata/search` | ✅ JWT | Semantic search (embed query + cosine KNN) |
| POST | `/vectordata/seed` | ✅ JWT | Seed เอกสารตัวอย่าง (IoT sensor + alarm policy) |

> Auth ใช้ chain `Verifier(true) → Authenticator → CurrentUser → ActiveUser` เหมือนโมดูลอื่น
> ต้องส่ง `Authorization: Bearer <token>` (ล็อกอินผ่าน `/api/auth/login`)

### 4.2 ตัวอย่างการใช้งาน (PowerShell + curl.exe)

ขั้นแรกต้อง login เพื่อเอา token:

```powershell
$body = '{"email":"admin@admin.com","password":"admin"}'
Set-Content -Path "$env:TEMP\login.json" -Value $body
$resp = Invoke-RestMethod -Method Post -Uri "http://localhost:5000/api/auth/login" -ContentType "application/json" -Body (Get-Content "$env:TEMP\login.json" -Raw)
$token = $resp.access_token
$headers = @{ Authorization = "Bearer $token" }
```

**1) ตรวจสอบสถานะ**

```powershell
curl.exe -s http://localhost:5000/api/vectordata/health
```

**2) สร้าง index**

```powershell
curl.exe -s -X POST http://localhost:5000/api/vectordata/create-index -H "Authorization: Bearer $token" -H "Content-Type: application/json" -d "{\"dims\":768}"
```

**3) สร้างเอกสาร (embed อัตโนมัติผ่าน LLM)**

```powershell
$doc = '{"document_id":"doc-001","content":"เซนเซอร์วัดอุณหภูมิแจ้งเตือนเมื่อเกิน 60 องศา","source_type":"iot-manual"}'
Set-Content -Path "$env:TEMP\doc.json" -Value $doc
curl.exe -s -X POST http://localhost:5000/api/vectordata/documents -H "Authorization: Bearer $token" -H "Content-Type: application/json" -d "@$env:TEMP\doc.json"
```

**4) สร้างเอกสารด้วย vector สำเร็จรูป**

```json
{
  "document_id": "doc-002",
  "content": "คู่มือเซนเซอร์ความชื้น",
  "embedding": [0.123, -0.456, 0.789]
}
```

**5) ค้นหาแบบ semantic**

```powershell
$q = '{"query":"เซนเซอร์ตัวไหนแจ้งเตือนเรื่องอุณหภูมิ","k":3}'
Set-Content -Path "$env:TEMP\search.json" -Value $q
curl.exe -s -X POST http://localhost:5000/api/vectordata/search -H "Authorization: Bearer $token" -H "Content-Type: application/json" -d "@$env:TEMP\search.json"
```

**6) Seed ข้อมูลตัวอย่าง**

```powershell
curl.exe -s -X POST http://localhost:5000/api/vectordata/seed -H "Authorization: Bearer $token" -H "Content-Type: application/json" -d "{\"count\":5}"
```

---

## 5. CLI Command (Seed ผ่าน Command Line)

Seed ข้อมูลตัวอย่างโดยไม่ต้องล็อกอิน ผ่าน command:

```bash
go run main.go vectordata seed
```

- อ่าน provider จาก `vectorDb.provider` ใน config โดยอัตโนมัติ
- ถ้าเป็น `pgvector` จะเชื่อม PostgreSQL; ถ้าเป็น `elasticsearch` จะเชื่อม ES
- ต้องการให้ **LLM embedding service ทำงาน** (`go run main.go serve` ขึ้นก่อน หรือ Ollama เปิดอยู่)
- เรียก `EnsureIndex` อัตโนมัติก่อน seed

ข้อมูลตัวอย่าง (5 ชิ้น): คู่มือเซนเซอร์อุณหภูมิ / ความชื้น / แรงสั่นสะเทือน / power meter + นโยบายการแจ้งเตือน

---

## 6. โครงสร้างข้อมูลใน Database

### 6.1 pgvector — ตาราง `vector_documents`

```sql
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS vector_documents (
    id          BIGSERIAL PRIMARY KEY,
    document_id VARCHAR(255) NOT NULL,
    content     TEXT NOT NULL,
    embedding   vector(768),
    source_type VARCHAR(100) NOT NULL DEFAULT '',
    metadata    TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS vector_documents_embedding_hnsw
    ON vector_documents USING hnsw (embedding vector_cosine_ops);
```

- เก็บเวกเตอร์ในคอลัมน์ `embedding` ชนิด `vector(dims)` — ค่า `dims` ต้องตรงกับมิติของ embedding จาก LLM
- ใช้ **HNSW index** + `vector_cosine_ops` เพื่อค้นหา cosine similarity ได้เร็ว
- Query ภายใน: คำนวณ score จาก `1 - (embedding <=> $1::vector)` (distance `<=>` = cosine distance)
- ตาราง + index สร้างอัตโนมัติโดย `EnsureIndex` (เรียกจาก migrate / create-index / seed)

### 6.2 Elasticsearch — index `vector_documents` (dense_vector)

Mapping ที่สร้างอัตโนมัติ:

```json
{
  "mappings": {
    "properties": {
      "document_id": { "type": "keyword" },
      "content":     { "type": "text" },
      "embedding":   { "type": "dense_vector", "dims": 768, "index": true, "similarity": "cosine" },
      "source_type": { "type": "keyword" },
      "timestamp":   { "type": "date" }
    }
  }
}
```

---

## 7. Edge Cases

| กรณี | พฤติกรรม |
|------|----------|
| `provider` ไม่มีค่า (ว่าง) | ใช้ `pgvector` เป็น default |
| ไม่มี Elasticsearch แต่ตั้ง `provider: elasticsearch` | `vectordb.New` คืน error → routes ถูกข้าม พร้อม log `❌ Failed to init vector DB` |
| ตั้ง `provider: pgvector` แต่ db ไม่พร้อม / image ไม่มี extension | error ตอน `EnsureIndex` → response 500 พร้อมข้อความชี้ชัด เช่น `pgvector: enable extension` |
| `dims` ใน config ไม่ตรงกับมิติจริงของ embedding | ข้อมูลอาจ insert ไม่ได้ (pgvector จะ error เรื่องมิติ vector) → ตั้ง `VECTOR_DB_DIMS` ให้ตรงกับโมเดล |
| ส่ง `embedding` มาด้วยใน request | ระบบใช้ vector ที่ส่งมาเลย ไม่เรียก LLM (ประหยัดรอบ request) |
| ส่ง request ที่ไม่มี `content` | ตอบ 400 `content is required` |
| `GET /documents` ไม่มี query | default `limit=10, offset=0` |
| `POST /search` ไม่มี `query` | ตอบ 400 `query is required` |
| LLM ยังไม่เปิด / model ไม่มี | `POST /documents`, `/search`, `/seed` ตอบ 500 พร้อม error จาก LLM (เช่น `llm embeddings returned ...`) |
| มีตาราง/index อยู่แล้ว | `EnsureIndex` เป็น idempotent — ไม่ error |
| Query เป็นข้อความไทยล้วน (ไม่มีตัวอักษรละติน/ตัวเลข) | `nomic-embed-text` (Ollama) อาจคืนเวกเตอร์เหมือนกันสำหรับข้อความไทยต่างกัน (ข้อจำกัดของ tokenizer/model) → ผล search อาจเหมือนเดิม; ตรวจได้โดย embed ข้อความต่างกัน 2 ข้อความผ่าน `/v1/embeddings` แล้วเทียบเวกเตอร์; แนวทางแก้: ใช้โมเดล embed ที่รองรับภาษาไทยดีกว่า หรือเติมคำ/ตัวอักษรละตินประกอบบริบท |

---

## 8. Troubleshooting

### 8.1 `pgvector: enable extension` error ตอน migrate

**สาเหตุ:** PostgreSQL ที่ใช้ไม่มี pgvector extension (เช่น ใช้ image `postgres:*` ธรรมดา)

**แก้ไข:**
- ใช้ image `pgvector/pgvector:pg15` (ดู docker-compose `db` service)
- หรือติดตั้ง pgvector ด้วยตัวเอง: https://github.com/pgvector/pgvector
- ถ้า container เก่า (สร้างจาก image เก่า) ต้อง `docker compose down && docker compose up -d db` หรือลบ volume `app-postgres-data` (⚠️ สำรองข้อมูลก่อน)

### 8.2 สลับ provider แล้วข้อมูลเดิมไม่อยู่

**สาเหตุ:** pgvector เก็บใน PostgreSQL ส่วน Elasticsearch เก็บใน index ES — เป็นคนละที่กัน

**แก้ไข:** ทำ `POST /api/vectordata/seed` หรือ re-index ใหม่ใน backend ที่เลือก

### 8.3 Embedding ใช้งานไม่ได้ (`llm embeddings returned ...`)

- ตรวจว่า Ollama เปิด: `curl.exe -s http://localhost:11434/api/version`
- ตรวจว่า pull โมเดลแล้ว: `ollama list` (ต้องมี `nomic-embed-text`)
- ตรวจ config `llm.base_url` ต้องลงท้ายด้วย `/v1` → `POST {base_url}/embeddings` ต้องตอบ 200
- ดูรายละเอียดเพิ่มเติมใน `README_LLM.md`

### 8.4 อยากตรวจสอบข้อมูลตรง ๆ (pgvector)

```powershell
docker exec -it icmongolang-db-1 psql -U icm -d icm -p 5435 -c "SELECT id, document_id, source_type, left(content,50) FROM vector_documents;"
```

> (เปลี่ยนชื่อ container / user / db ตาม `.env` ของคุณ)

---

## 9. Test Cases

| # | กลุ่ม | ขั้นตอน | ผลลัพธ์ที่คาดหวัง |
|---|------|---------|------------------|
| TC1 | Health | `GET /api/vectordata/health` | 200 `{"provider":"pgvector","status":"up",...}` |
| TC2 | Health | เปิดโดยไม่ตั้ง provider (default) | ใช้ `pgvector` |
| TC3 | Index | `POST /api/vectordata/create-index` ด้วย `{"dims":768}` | 200 `{"status":"ok"}` |
| TC4 | Index | เรียก create-index ซ้ำ 2 ครั้ง | ไม่ error (idempotent) |
| TC5 | CRUD | `POST /documents` ด้วย content | 200 มี `id`, `dim=768`, `indexed=true` |
| TC6 | CRUD | `POST /documents` ไร้ content | 400 `content is required` |
| TC7 | CRUD | `POST /documents` พร้อม `embedding` สำเร็จรูป | 200 ไม่ต้องเรียก LLM |
| TC8 | CRUD | `GET /documents/{id}` ตาม id จาก TC5 | 200 ได้เอกสารคืน |
| TC9 | CRUD | `GET /documents/{id}` ด้วย id ที่ไม่มี | 404 `document ... not found` |
| TC10 | CRUD | `GET /documents?limit=2&offset=0` | 200 `total`, `limit`, `offset`, `documents` |
| TC11 | CRUD | `PUT /documents/{id}` เปลี่ยน content | 200 `indexed=true` + re-embed |
| TC12 | CRUD | `DELETE /documents/{id}` | 200 `{"status":"ok"}` แล้ว GET คืน 404 |
| TC13 | Search | `POST /search` ด้วยคำถามที่ใกล้เคียงเอกสาร | 200 มี `hits` เรียง score มากไปน้อย |
| TC14 | Search | `POST /search` ไร้ query | 400 `query is required` |
| TC15 | Auth | เรียก CRUD โดยไม่ส่ง token | 401 |
| TC16 | Seed | `POST /seed` ด้วย `{"count":5}` | 200 `{"seeded":5}` |
| TC17 | Seed | CLI: `go run main.go vectordata seed` | log `🌱 Seeded 5 sample documents` |
| TC18 | ES | เปลี่ยน `VECTOR_DB_PROVIDER=elasticsearch` แล้ว health | 200 `{"provider":"elasticsearch","status":"up"}` |
| TC19 | Health | `/api/health` | มี field `vectordata: true` |
| TC20 | Search | ค้นหาคำที่ไม่เกี่ยวข้อง | 200 `hits` ว่างหรือ score ต่ำมาก |
| TC21 | Migration | `go run main.go migrate` กับ `pgvector` | ตาราง `vector_documents` + HNSW index สร้างสำเร็จ |
