# คู่มือการใช้งานโมดูล Elasticsearch (Vector Search)

โมดูล `internal/modules/elasticsearch` ทำหน้าที่เป็น **Semantic / Vector Search** ให้กับระบบ IC-MON
โดยทำงานร่วมกับ 2 ส่วนหลัก:

1. **LLM Embedding Client** (`pkg/llm`) — แปลงข้อความเป็นเวกเตอร์ (embedding) ผ่าน API แบบ OpenAI-compatible (รองรับ Ollama, OpenAI, Azure OpenAI, ฯลฯ)
2. **Elasticsearch Client** (`pkg/elasticsearch`) — เก็บเวกเตอร์ใน field `dense_vector` และค้นหาแบบ **cosine similarity (KNN)**

---

## 1. สถาปัตยกรรม (Architecture)

โมดูลแบ่งเป็น 3 ชั้น ตามรูปแบบของทั้งโปรเจกต์:

```
Client / Postman / Frontend
        │ HTTP (chi router)
        ▼
┌─────────────────────────────────────────────────┐
│ delivery/http                                     │
│   handler.go  → รับ HTTP request / ส่ง response    │
│   routes.go   → ลงทะเบียน route ใต้ /api/elasticsearch │
└─────────────────────────────────────────────────┘
        │
        ▼
┌─────────────────────────────────────────────────┐
│ usecase                                          │
│   usecase.go  → ประสานงาน LLM + Elasticsearch     │
└─────────────────────────────────────────────────┘
        │
        ▼
┌───────────────────────────────┐   ┌───────────────────────┐
│ pkg/llm (embedding)           │   │ pkg/elasticsearch      │
│  OpenAI-compatible /v1        │   │  REST client           │
└───────────────────────────────┘   └───────────────────────┘
        │                               │
        ▼                               ▼
  Ollama / OpenAI                Elasticsearch (9200)
```

### 1.1 ตารางไฟล์

| ไฟล์ | บทบาท |
|------|-------|
| `internal/modules/elasticsearch/presenter/presenter.go` | DTO — โครงสร้าง Request / Response ที่ใช้ใน API |
| `internal/modules/elasticsearch/usecase/usecase.go` | Business logic — orchestrate LLM embedding + ES indexing/search |
| `internal/modules/elasticsearch/delivery/http/handler.go` | HTTP handler + ตรวจสอบข้อมูลเบื้องต้น + error mapping |
| `internal/modules/elasticsearch/delivery/http/routes.go` | ลงทะเบียน route พร้อม middleware |
| `pkg/elasticsearch/client.go` | HTTP client สำหรับยิง REST API ไป Elasticsearch |
| `pkg/llm/client.go` | HTTP client สำหรับสร้าง embedding (OpenAI-compatible) |

### 1.2 จุดที่โมดูลถูก connect (wiring)

- **`internal/server/server.go`** — สร้าง `esClient` (นิ่งได้ถ้าเชื่อมไม่สำเร็จ) และ `llmClient`
- **`internal/server/handlers.go`** — ถ้าทั้ง `esClient` และ `llmClient` ไม่เป็น nil จะสร้าง UseCase + Handler แล้วลง route ผ่าน `MapElasticsearchRoutes(apiRouter, esHandler, mw)` ไม่งั้นข้าม (skip) พร้อม log คำเตือน
- Base path ของ API ทั้งหมดคือ `/api` → endpoints นี้อยู่ใต้ `/api/elasticsearch/...`

---

## 2. ข้อกำหนดของระบบ (Prerequisites)

| Component | รายละเอียด |
|-----------|-----------|
| Elasticsearch | เวอร์ชันที่รองรับ field `dense_vector` (>= 7.x สำหรับ KNN ผ่าน `knn` query; ใน docker-compose ใช้ **8.13.4**) |
| LLM Embedding Provider | API แบบ OpenAI-compatible ที่มี endpoint `POST {base_url}/embeddings` เช่น **Ollama** (ค่า default), OpenAI, Azure OpenAI |

---

## 3. การตั้งค่า Configuration

โมดูลอ่านค่าจาก config 2 บล็อก ได้แก่ `elasticsearch` และ `llm` (จาก `config/config.go`)

### 3.1 ตัวแปรสภาพแวดล้อม (`.env` / `.env.example`)

```env
# ---- Elasticsearch ----
ELASTICSEARCH_ADDRESSES=http://localhost:9200   # รองรับหลาย address คั่นด้วยเครื่องหมายจุลภาค
ELASTICSEARCH_USERNAME=                          # ใส่ username ถ้า ES เปิดใช้ auth (เช่น user ใน xpack)
ELASTICSEARCH_PASSWORD=
ELASTICSEARCH_INDEX=vector_embeddings            # ชื่อ index ที่ใช้เก็บเวกเตอร์
ELASTICSEARCH_TIMEOUT=10                         # HTTP timeout หน่วยวินาที (default 10)

# ---- LLM Embedding ----
LLM_BASE_URL=http://localhost:11434/v1           # Ollama เปิด OpenAI-compatible API ไว้ที่ /v1
LLM_API_KEY=
LLM_MODEL=nomic-embed-text                       # ชื่อโมเดล embedding
LLM_TIMEOUT=30                                   # HTTP timeout หน่วยวินาที (default 30)
LLM_DIMS=768                                     # ขนาดเวกเตอร์ของโมเดล (nomic-embed-text = 768)
```

> **สำคัญ:** ค่า `LLM_DIMS` ต้องตรงกับขนาดเวกเตอร์ที่โมเดลสร้างจริง มิฉะนั้น Elasticsearch จะ reject ตอน index (ข้อผิดพลาด `dense_vector_dims_exception`)

### 3.2 ค่าเริ่มต้นใน code (`config/config.default.yml`)

```yaml
elasticsearch:
  addresses: ["http://localhost:9200"]
  username: ""
  password: ""
  index: "vector_embeddings"
  timeout: 10

llm:
  base_url: "http://localhost:11434/v1"
  api_key: ""
  model: "nomic-embed-text"
  timeout: 30
  dims: 768
```

---

## 4. การยืนยันตัวตน (Authentication)

| Endpoint | ต้อง Login? |
|----------|-------------|
| `GET /health` | ❌ ไม่ต้อง (Public) |
| 3 Endpoint ที่เหลือ | ✅ ต้องใช้ JWT access token |

Route ที่ต้อง auth จะผ่าน middleware เรียงตามลำดับ:
`Verifier(true)` → `Authenticator()` → `CurrentUser()` → `ActiveUser()`
(ดู `internal/middleware/jwtauth.go`)

ส่ง token ใน header:

```
Authorization: Bearer <ACCESS_TOKEN>
```

**วิธีขอ token** (POST `/api/auth/login` — ใช้ multipart/form-data):

```bash
curl -X POST http://localhost:5000/api/auth/login \
  -H "Content-Type: multipart/form-data" \
  -F "username=root@gmail.com" \
  -F "password=root_password"
```

นำค่า `access_token` ที่ได้ไปใช้ในคำสั่งด้านล่าง (ตัวอย่างแทนด้วย `$TOKEN`)

> ดู Swagger UI ได้ที่ `http://localhost:5000/swagger/index.html` (หมวด `elasticsearch`)

---

## 5. API Reference

### 5.1 `GET /api/elasticsearch/health` — ตรวจสุขภาพคลัสเตอร์ ES

- **ไม่ต้อง auth**
- เรียก `GET /_cluster/health` ไปยัง Elasticsearch

**Request:** ไม่มี body

```bash
curl http://localhost:5000/api/elasticsearch/health
```

**Response 200** (ตัวอย่างจริง ขึ้นอยู่กับสถานะคลัสเตอร์):

```json
{
  "cluster_name": "docker-cluster",
  "status": "green",
  "timed_out": false,
  "number_of_nodes": 1,
  "number_of_data_nodes": 1,
  "active_primary_shards": 1,
  "active_shards": 1,
  "unassigned_shards": 0,
  ...
}
```

**Error:** 500 ถ้า ES ถูกปิดหรือไม่สามารถเชื่อมต่อ

---

### 5.2 `POST /api/elasticsearch/create-index` — สร้าง (หรือยืนยันว่า) index มีอยู่

- **ต้อง auth**
- สร้าง index พร้อม mapping ของ `dense_vector` (idempotent — ถ้า index มีอยู่แล้วจะถือว่าสำเร็จ)
- เรียก `EnsureIndex` ซึ่งใช้ค่า dims จาก:

```go
dims := u.llm.Dims()          // จาก LLM_DIMS
if req.Dims > 0 { dims = req.Dims }   // หรือระบุใน body
```

**Request body:**

```json
{ "dims": 768 }
```

`dims` ไม่บังคับ — ถ้าใส่ 0 หรือไม่ส่ง จะใช้ค่า `LLM_DIMS` (ถ้าเป็น 0 อีกจะ fallback เป็น **1536**)

```bash
curl -X POST http://localhost:5000/api/elasticsearch/create-index \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "dims": 768 }'
```

**Response 200:**

```json
{ "status": "ok", "message": "index ready" }
```

**Mapping ที่สร้าง (ดูจาก `pkg/elasticsearch/client.go:94`)**

```json
{
  "mappings": {
    "properties": {
      "document_id": { "type": "keyword" },
      "content":     { "type": "text" },
      "embedding":   {
        "type": "dense_vector",
        "dims": 768,
        "index": true,
        "similarity": "cosine"
      },
      "source_type": { "type": "keyword" },
      "timestamp":   { "type": "date" }
    }
  }
}
```

> หมายเหตุ: `resource_already_exists_exception` (HTTP 400) จะถูกถือว่าเป็นกรณีปกติและไม่ error

---

### 5.3 `POST /api/elasticsearch/embeddings` — แปลงข้อความเป็นเวกเตอร์ (ด้วย LLM) แล้วเก็บลง ES

- **ต้อง auth**
- ขั้นตอน (ดู `usecase.go:51`):
  1. `EnsureIndex` (ใช้ dims ของ LLM)
  2. `llm.Embed(ctx, content)` — สร้าง embedding
  3. `es.Index(ctx, doc)` — POST ไปที่ `/{index}/_doc` (ES สร้าง `_id` ให้อัตโนมัติ)
  4. `timestamp` ถูกเติมอัตโนมัติเป็นเวลาปัจจุบัน

**Request body:**

```json
{
  "document_id": "doc-001",
  "content": "ขั้นตอนการตั้งค่า Modbus TCP ใน PLC Siemens S7-1500",
  "source_type": "manual"
}
```

| Field | ประเภท | บังคับ | คำอธิบาย |
|-------|--------|-------|----------|
| `document_id` | string | ❌ | ID อ้างอิงเอกสารต้นทาง (สำหรับลิงก์กลับ) |
| `content` | string | ✅ | ข้อความที่จะ embed (ต้องไม่เว้นว่าง) |
| `source_type` | string | ❌ | ประเภทต้นทาง เช่น `manual`, `knowledge`, `report` |

```bash
curl -X POST http://localhost:5000/api/elasticsearch/embeddings \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "document_id": "doc-001",
    "content": "ขั้นตอนการตั้งค่า Modbus TCP ใน PLC Siemens S7-1500",
    "source_type": "manual"
  }'
```

**Response 200 (`EmbeddingResponse`):**

```json
{
  "document_id": "doc-001",
  "content": "ขั้นตอนการตั้งค่า Modbus TCP ใน PLC Siemens S7-1500",
  "source_type": "manual",
  "index": "vector_embeddings",
  "model": "nomic-embed-text",
  "dim": 768,
  "indexed": true
}
```

**หมายเหตุ:**
- field `_id` ของ ES ถูกสร้างอัตโนมัติ (ต่างจาก `document_id` ซึ่งเป็น field ใน document)
- `model` เป็น `omitempty` — จะหายไปถ้า LLM ไม่ได้ตั้งชื่อโมเดลไว้

---

### 5.4 `POST /api/elasticsearch/vectors` — เก็บเวกเตอร์ที่คำนวณสำเร็จแล้วโดยตรง (ไม่ต้องใช้ LLM)

- **ต้อง auth**
- เหมาะกับกรณีที่ embedding ถูกสร้างไว้ที่อื่น (เช่น offline / แบตช์) แล้วอยากยัดเข้า ES โดยตรง
- ใช้ `dims = len(req.Embedding)` ในการ EnsureIndex (ดู `usecase.go:91`)

**Request body:**

```json
{
  "document_id": "doc-002",
  "content": "มาตรฐานการเดินสายไฟฟ้า IEC 60364",
  "embedding": [0.0123, -0.0456, 0.0789, 0.0012, 0.0911],
  "source_type": "standard"
}
```

| Field | ประเภท | บังคับ | คำอธิบาย |
|-------|--------|-------|----------|
| `document_id` | string | ❌ | ID อ้างอิงเอกสารต้นทาง |
| `content` | string | ✅ | ข้อความต้นฉบับ (ต้องไม่ว่าง) |
| `embedding` | []float32 | ✅ | เวกเตอร์ที่คำนวณมาแล้ว (ต้องมี >= 1 ค่า) |
| `source_type` | string | ❌ | ประเภทต้นทาง |

```bash
curl -X POST http://localhost:5000/api/elasticsearch/vectors \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "document_id": "doc-002",
    "content": "มาตรฐานการเดินสายไฟฟ้า IEC 60364",
    "embedding": [0.0123, -0.0456, 0.0789],
    "source_type": "standard"
  }'
```

**Response 200:** เหมือน `EmbeddingResponse` แต่ `model` จะหายไป (omitempty) และ `dim` = ขนาดเวกเตอร์ที่ส่งเข้า

```json
{
  "document_id": "doc-002",
  "content": "มาตรฐานการเดินสายไฟฟ้า IEC 60364",
  "source_type": "standard",
  "index": "vector_embeddings",
  "dim": 3,
  "indexed": true
}
```

---

### 5.5 `POST /api/elasticsearch/search` — ค้นหาแบบ Semantic (Cosine Similarity)

- **ต้อง auth**
- ขั้นตอน (ดู `usecase.go:121`):
  1. Embed ข้อความคำถามด้วย LLM
  2. `es.SearchKNN(ctx, embedding, k)` — ส่ง query `knn` ไปที่ `/{index}/_search`
  3. ค่าเริ่มต้น `k = 10` ถ้าไม่ระบุหรือ <= 0
  4. `num_candidates = max(k*10, 100)`

**Request body:**

```json
{
  "query": "วิธีตั้งค่า Modbus กับ PLC",
  "k": 5
}
```

| Field | ประเภท | บังคับ | คำอธิบาย |
|-------|--------|-------|----------|
| `query` | string | ✅ | ข้อความคำถาม (ต้องไม่ว่าง) |
| `k` | int | ❌ | จำนวนผลลัพธ์ที่ต้องการ (default 10) |

```bash
curl -X POST http://localhost:5000/api/elasticsearch/search \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "query": "วิธีตั้งค่า Modbus กับ PLC", "k": 5 }'
```

**Response 200 (`SearchResponse`):**

```json
{
  "query": "วิธีตั้งค่า Modbus กับ PLC",
  "k": 5,
  "hits": [
    {
      "document_id": "doc-001",
      "content": "ขั้นตอนการตั้งค่า Modbus TCP ใน PLC Siemens S7-1500",
      "source_type": "manual",
      "score": 0.8123456
    },
    {
      "document_id": "doc-002",
      "content": "มาตรฐานการเดินสายไฟฟ้า IEC 60364",
      "source_type": "standard",
      "score": 0.654321
    }
  ]
}
```

- `score` คือ cosine similarity (0–1) — ยิ่งใกล้ 1 ยิ่งคล้าย
- ผลลัพธ์ถูกเรียงจากคล้ายมากไปน้อยแล้ว (ES KNN)

---

## 6. ตัวอย่างการใช้งานครบวงจร

### 6.1 ขั้นพื้นฐาน

```bash
# 1) ตรวจสอบว่า ES และ LLM พร้อมใช้งาน
curl http://localhost:5000/api/elasticsearch/health

# 2) Login เพื่อเอา token
TOKEN=$(curl -s -X POST http://localhost:5000/api/auth/login \
  -F "username=root@gmail.com" -F "password=root_password" | jq -r .access_token)

# 3) สร้าง index (optional — ถ้าไม่มีจะสร้างให้อัตโนมัติตอน embed/index)
curl -X POST http://localhost:5000/api/elasticsearch/create-index \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{ "dims": 768 }'

# 4) เก็บความรู้ (embed ผ่าน LLM)
curl -X POST http://localhost:5000/api/elasticsearch/embeddings \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{ "document_id": "kb-1", "content": "คำสั่งเปิดปิดวาล์วด้วย Modbus", "source_type": "knowledge" }'

# 5) ค้นหา
curl -X POST http://localhost:5000/api/elasticsearch/search \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{ "query": "เปิดวาล์วอัตโนมัติยังไง", "k": 5 }'
```

### 6.2 ทดสอบด้วย Docker

โปรเจกต์มี `docker-compose.yml` ซึ่งรวม:
- **Elasticsearch 8.13.4** — port `9200`
- **Ollama** — ใช้ `nomic-embed-text` (dims 768)
- **elasticsearch-exporter** — port `9114` (สำหรับ Prometheus/Grafana)

```bash
docker compose up -d elasticsearch ollama
# ดึงโมเดล embedding (ครั้งแรก)
docker exec -it ollama ollama pull nomic-embed-text
```

> ในโหมด Docker ค่า `ELASTICSEARCH_ADDRESSES` จะถูก override เป็น `http://elasticsearch:9200` และ `LLM_BASE_URL` เป็น `http://ollama:11434/v1` (ดู `docker-compose.yml`)

---

## 7. พฤติกรรมและข้อควรระวัง (Edge Cases / Notes)

1. **โมดูลเป็น optional** — ถ้า ES หรือ LLM เชื่อมไม่สำเร็จ route ทั้งหมดจะถูกข้ามไป (log `⚠️ Elasticsearch client is nil – ES routes skipped`) ไม่ทำให้ระบบหลักล่ม
2. **`llm` เป็น nil** — `EmbedAndIndex` / `SearchVector` จะ return error `"llm client is nil"` และ ES client ไม่มี fallback นี้
3. **การ fallback ของ dims** — ถ้า LLM ไม่ได้ตั้ง `LLM_DIMS` จะเป็น 1536; ถ้าสร้าง index ด้วย `dims` หนึ่งแล้วสลับโมเดลที่มี dims ต่างกัน จะ index ไม่ได้ (ต้องลบ index แล้วสร้างใหม่)
4. **Failover หลาย address** — `Client.do()` ลอง address ทีละตัว; ถ้าตัวแรกตอบ 5xx และมี address มากกว่า 1 จะลองตัวถัดไป (ดู `client.go:270`)
5. **Basic Auth** — ถ้าตั้ง `ELASTICSEARCH_USERNAME/PASSWORD` client จะส่ง `Authorization: Basic` ทุก request
6. **`document_id` กับ `_id`** — ต่างกัน: `_id` ของ ES ถูก generate อัตโนมัติ ส่วน `document_id` คือข้อมูลที่เรากำหนดเอง
7. **Validation** — handler ตรวจสอบ `content`/`query`/`embedding` ไม่ว่างเอง (HTTP 400) แต่การ validate ด้วย tag `required` ใน presenter จะทำงานเมื่อ route ใช้ validator middleware ด้วย
8. **Error format** — error response มีรูปแบบ `{ "error": "<ข้อความ>" }` พร้อม HTTP status 400 (ข้อมูลไม่ครบ/ไม่ถูกต้อง) หรือ 500 (ข้อผิดพลาดภายใน/เชื่อมต่อ ES หรือ LLM ไม่ได้)

---

## 8. การเรียกใช้จากโค้ด Go (โดยตรง ไม่ผ่าน HTTP)

โมดูลนี้ถูกออกแบบให้เรียกผ่าน HTTP เป็นหลัก แต่ component ย่อยสามารถ reuse ได้:

```go
import (
    "context"
    "icmongolang/config"
    "icmongolang/pkg/elasticsearch"
    "icmongolang/pkg/llm"
)

ctx := context.Background()

es, err := elasticsearch.NewClient(&cfg.Elasticsearch, logger)
llmClient := llm.NewClient(&cfg.LLM, logger)

// สร้าง embedding แล้วเก็บ
emb, err := llmClient.Embed(ctx, "ข้อความตัวอย่าง")
es.EnsureIndex(ctx, len(emb))
doc := &elasticsearch.VectorDoc{
    DocumentID: "doc-1",
    Content:    "ข้อความตัวอย่าง",
    Embedding:  emb,
    SourceType: "demo",
}
id, err := es.Index(ctx, doc)

// ค้นหา
hits, err := es.SearchKNN(ctx, emb, 10)
```

---

## 9. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| health ไม่ทำงาน / connection failed | ตรวจ ES เปิดอยู่ที่ `9200` หรือไม่ (`curl localhost:9200`); ดู `ELASTICSEARCH_ADDRESSES` |
| `dense_vector_dims_exception` ตอน index | `LLM_DIMS` กับขนาดเวกเตอร์จริงไม่ตรงกัน → ตั้งค่าให้ถูก แล้วลบ index เดิม (`DELETE /vector_embeddings`) เพื่อสร้างใหม่ |
| embedding ไม่มา / timeout | ตรวจ Ollama: `curl http://localhost:11434/v1/embeddings`; ดู `LLM_BASE_URL`, `LLM_MODEL` ว่ามีโมเดลจริงหรือยัง (`ollama list`) |
| `resource_already_exists_exception` | ปกติ — create-index เป็น idempotent ไม่ใช่ error |
| 401 Unauthorized | token หมดอายุ / ไม่ได้ส่ง `Authorization: Bearer` / ผู้ใช้ถูกปิด (ActiveUser) |
| search คืนผลว่าง | ยังไม่มีข้อมูลใน index, หรือ embedding ระหว่าง query กับ documents ต่างชุดกันอย่างสิ้นเชิง (score ต่ำ) |
| route หายไป (404) | ทั้ง ES และ LLM ต้องเชื่อมสำเร็จพร้อมกัน ไม่งั้น route ไม่ถูก register ดู log ตอน start |

---

## 10. โครงสร้างข้อมูลใน Elasticsearch

Document แต่ละตัวใน index (เช่น `vector_embeddings`) มีหน้าตาแบบนี้:

```json
{
  "_index": "vector_embeddings",
  "_id": "auto-generated-id",
  "_source": {
    "document_id": "doc-001",
    "content": "ขั้นตอนการตั้งค่า Modbus TCP ใน PLC Siemens S7-1500",
    "embedding": [0.0123, -0.0456, ...],
    "source_type": "manual",
    "timestamp": "2026-08-14T10:30:00.000Z"
  }
}
```

| Field | Type | หมายเหตุ |
|-------|------|----------|
| `document_id` | keyword | ใช้อ้างอิงเอกสารต้นทาง |
| `content` | text | ข้อความต้นฉบับ (แสดงผลในการค้นหา) |
| `embedding` | dense_vector (cosine) | เวกเตอร์ 768 dims (ตามโมเดล) |
| `source_type` | keyword | ใช้ filter กลุ่มข้อมูลได้ภายหลัง |
| `timestamp` | date | เติมอัตโนมัติเมื่อ index |
