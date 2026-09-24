# คู่มือการใช้งาน LLM (Large Language Model)

โปรเจกต์ `icmongolang` มีระบบ AI/LLM **2 ชุด** แยกกัน:

| ชุด | ที่อยู่ | บทบาท | ผู้เรียก |
|-----|--------|-------|---------|
| **Embedding client** | `pkg/llm/` (Go) | แปลง text → vector embedding (OpenAI-compatible `/v1/embeddings`) | elasticsearch vector search module |
| **Chat service** | `ai/` (Go + langchaingo) | ให้บริการแชต/ถาม-ตอบผ่าน Ollama บน port 8001 | ผู้ใช้/backend เรียกตรง |

ทั้งสองชุดคุยกับ **Ollama** ที่รันเป็น service ใน docker-compose (ดู `README_Ollama.md`)

```
┌──────────────┐  POST /v1/embeddings  ┌──────────────┐
│  pkg/llm     │ ────────────────────▶ │    Ollama    │  nomic-embed-text
│  (embedding) │                       │   :11434     │  (dims 768)
└──────┬───────┘                       └──────┬───────┘
       │                                     │
┌──────▼─────────┐  POST /v1/chat           │
│ ai service :8001│ ────────────────────────┤  qwen3:4b / deepseek-r1:7b
└────────────────┘                          │
```

---

## 1. การตั้งค่า Configuration

### 1.1 Embedding client — `config/config.default.yml:103-110`
```yaml
llm:
  # Embedding provider (OpenAI-compatible /v1 endpoint, served by Ollama)
  base_url: "http://localhost:11434/v1"
  api_key: ""
  model: "nomic-embed-text"
  timeout: 30
  dims: 768
```

Env vars (ผ่าน `config.go` BindEnvs): `LLM_BASE_URL`, `LLM_API_KEY`, `LLM_MODEL`, `LLM_TIMEOUT`, `LLM_DIMS`

> ใน docker-compose backend ตั้ง `LLM_BASE_URL=http://ollama:11434/v1`

### 1.2 Chat service (`ai/`) — env
| env | default | หมายเหตุ |
|-----|---------|----------|
| `PORT` | `8001` | port ที่ ai service ฟัง |
| `OLLAMA_BASE_URL` | `http://localhost:11434` | ที่อยู่ Ollama (ไม่มี `/v1`) |
| `OLLAMA_MODEL` | `deepseek-r1:7b` (code) / `qwen3:4b` (compose) | โมเดลแชต |
| `AI_TIMEOUT` | `180s` | HTTP timeout ไป Ollama |

### 1.3 โมเดลที่ต้อง pull ใน Ollama
```bash
docker compose exec ollama ollama pull nomic-embed-text   # embedding (dims 768)
docker compose exec ollama ollama pull qwen3:4b           # chat (ถ้าใช้ตัวนี้)
```

---

## 2. Embedding Client — `pkg/llm/client.go`

เป็น **OpenAI-compatible embedding client** — ใช้ได้กับ OpenAI, Azure OpenAI, Ollama, หรือเซิร์ฟเวอร์ local ที่ให้ API แบบเดียวกัน

### 2.1 การสร้าง client
```go
// internal/server/server.go:122
llmClient := llm.NewClient(&cfg.LLM, log)
```
- `timeout` ≤ 0 → default 30s
- `dims` ≤ 0 → default 1536
- `baseURL` ตัด `/` ท้ายออก แล้วต่อ `/embeddings`

### 2.2 Method
| Method | หน้าที่ |
|--------|---------|
| `Dims() int` | ขนาด embedding (สำหรับสร้าง index) |
| `Model() string` | ชื่อโมเดล |
| `Embed(ctx, text) ([]float32, error)` | แปลง text → vector |

### 2.3 `Embed()` ทำงานอย่างไร (`client.go:61-117`)
1. POST `{baseURL}/embeddings` body `{"model": ..., "input": text}`
2. ตั้ง header `Content-Type: application/json`, `Accept: application/json`, และ `Authorization: Bearer <api_key>` ถ้ามี key
3. status ไม่ใช่ 2xx → error `"llm embeddings returned %d: %s"`
4. decode `{"data":[{"embedding":[...]}]}` → คืน `data[0].Embedding`
5. error cases: `llm request failed`, `llm response decode`, `llm error: <message>`, `llm returned empty embedding`

---

## 3. Vector Search (LLM + Elasticsearch)

LLM embedding ถูกใช้ใน **elasticsearch module** (`/api/elasticsearch`) เพื่อทำ semantic search — endpoints ใช้ JWT (ยกเว้น `/health`):

```go
// internal/modules/elasticsearch/delivery/http/routes.go
r.Get("/health", h.Health)                     // ไม่ต้อง auth
r.Post("/create-index", h.CreateIndex)          // JWT
r.Post("/embeddings", h.EmbedAndIndex)          // JWT
r.Post("/vectors", h.IndexVector)               // JWT
r.Post("/search", h.SearchVector)               // JWT
```

### 3.1 `POST /api/elasticsearch/create-index`
```json
{ "dims": 768 }
```
สร้าง/ยืนยัน index `vector_embeddings` (จาก `cfg.Elasticsearch.Index`) ด้วย mapping `dense_vector` dims=768, similarity=cosine
- ไม่ส่ง dims → ใช้ `llm.Dims()` (768)
- index มีอยู่แล้ว → คืน success (จับ `resource_already_exists_exception`)

### 3.2 `POST /api/elasticsearch/embeddings` — embed + index เอกสาร
```json
{ "document_id": "doc-001", "content": "เซ็นเซอร์อุณหภูมิทำงานผิดปกติ", "source_type": "iot_alarm" }
```
ขั้นตอน (`usecase.go:51-88`):
1. `EnsureIndex` (dims จาก client)
2. `llm.Embed(content)` → vector
3. `es.Index(VectorDoc)` → เก็บใน ES
4. response:
```json
{
  "document_id": "doc-001",
  "content": "เซ็นเซอร์อุณหภูมิทำงานผิดปกติ",
  "source_type": "iot_alarm",
  "index": "vector_embeddings",
  "model": "nomic-embed-text",
  "dim": 768,
  "indexed": true
}
```

### 3.3 `POST /api/elasticsearch/vectors` — index vector ที่คำนวณไว้แล้ว
```json
{ "document_id": "doc-002", "content": "ข้อความ", "embedding": [0.1, 0.2, ...] }
```
ไม่เรียก LLM — เก็บ vector ตรง (response ไม่มี `model`)

### 3.4 `POST /api/elasticsearch/search` — semantic search
```json
{ "query": "ปัญหาเซ็นเซอร์", "k": 10 }
```
1. `llm.Embed(query)` → vector
2. `es.SearchKNN(embedding, k)` → knn search (`field=embedding`, `num_candidates = max(k*10, 100)`)
3. response:
```json
{
  "query": "ปัญหาเซ็นเซอร์",
  "k": 10,
  "hits": [
    { "document_id": "doc-001", "content": "...", "source_type": "iot_alarm", "score": 0.87 }
  ]
}
```
- `k` ≤ 0 → default 10

### 3.5 โครงสร้าง doc ใน ES (`pkg/elasticsearch/client.go`)
| field | type |
|-------|------|
| `document_id` | keyword |
| `content` | text |
| `embedding` | dense_vector (dims, index, similarity=cosine) |
| `source_type` | keyword |
| `timestamp` | date |

---

## 4. Chat Service — `ai/main.go`

แยก module ตัวเอง (Go + langchaingo `v0.1.14`) รันเป็น container `ai` (port 8001) — wrapper บาง ๆ เหนือ Ollama chat

### 4.1 endpoints
| Method | Path | Body | Response |
|--------|------|------|----------|
| GET | `/health` | — | `{"status":"ok","model":"..."}` |
| POST | `/v1/chat` | `{"messages":[{"role":"user","content":"..."}]}` | `{"reply":"...","model":"..."}` |
| POST | `/v1/generate` | `{"question":"...","context":"..."}` | `{"reply":"...","model":"..."}` |

### 4.2 `/v1/chat` — multi-turn
- รับ `messages` (ระบบ+user+assistant)
- ถ้า message แรกไม่ใช่ `role=system` → แทรก system prompt default
- เรียก Ollama ด้วย `MaxTokens=1024`, `Temperature=0.2`, `TopP=0.8`, `RunnerNumCtx=2048`
- response: `{"reply": resp.Choices[0].Content, "model": ...}`

### 4.3 `/v1/generate` — Q&A พร้อม context (RAG-ish)
- ใช้ `LLMChain` + prompt template:
```
[systemPrompt]

Context:
{{.context}}

Question:
{{.question}}
```
- เรียก `chains.Call(...)` ด้วย `MaxTokens=1024`
- เหมาะกับ: เอาข้อมูลจาก IoT/sensor มาใส่ context แล้วให้โมเดลตอบ

### 4.4 System prompt เริ่มต้น
```
You are an IoT operations assistant for the icmongolang platform.
Provide concise, factual answers about sensor data, alarms, MQTT topics, and platform events.
Use the provided context when relevant.
Be brief: keep answers under 150 words unless the user asks for detail.
```

### 4.5 รัน (docker-compose.yml:263-276)
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

ตัวอย่าง:
```bash
curl -s http://localhost:8001/health
curl -s http://localhost:8001/v1/chat -d '{"messages":[{"role":"user","content":"สรุปสถานะเซ็นเซอร์"}]}'
curl -s http://localhost:8001/v1/generate -d '{"question":"อุณหภูมิเกินไหม?","context":"Sensor A: 42C, Sensor B: 25C"}'
```

---

## 5. Edge Cases & ข้อควรระวัง

1. **embedding กับ chat ใช้คนละ endpoint** — `pkg/llm` ต่อ `/v1/embeddings`, `ai` ต่อ `/v1/chat`/`/api/generate`; ทั้งคู่ชี้ Ollama
2. **`dims` ต้องตรงกับโมเดล** — `nomic-embed-text` = 768; ถ้าเปลี่ยนโมเดล (เช่น 1024-dim) ต้องสร้าง index ใหม่ (เปลี่ยน dims ไม่ได้บน index เดิม)
3. **LLM client ถูกสร้างเสมอ** (`server.go:122`) — แต่ถ้า config `llm` ว่าง/ES ไม่เชื่อมต่อ ฟีเจอร์นี้เงียบไป (route ไม่ถูก map)
4. **elasticsearch routes เปิดเมื่อ ES + LLM ทั้งคู่ไม่ nil** (`handlers.go:201`) — ถ้าอันใดอันหนึ่ง fail จะไม่มี routes
5. **Ollama ต้อง pull โมเดลเอง** — ถ้ายังไม่มี `nomic-embed-text`/`qwen3:4b` จะ error ตอนเรียก
6. **`deepseek-r1:7b` (code default) ≠ `qwen3:4b` (compose default)** — ดูว่าใช้ค่าไหนจริงใน env
7. **ai service เป็น module แยก** — มี `go.mod` ของตัวเอง (`module icmongolang/ai`) ต้อง build ใน `ai/` ไม่ใช่ root
8. **ไม่มีการ rate limit/streaming** — `/v1/chat` และ `/v1/generate` ตอบแบบเต็มข้อความ (ไม่ SSE)
9. **`ai_timeout` ที่ compose (300s) มากกว่า code default (180s)** — env ชนะ
10. **ผลลัพธ์ AI เปลี่ยนได้** — temperature 0.2 ช่วยให้คงที่ แต่ไม่ deterministic 100%

---

## 6. Test Cases (กรณีทดสอบ)

> หมายเหตุ: `POST` ต้องมี body เสมอ — ถ้าเรียก GET จะได้ **405** (พฤติกรรมปกติของ Ollama)

### 6.1 Ollama API (ผ่าน `localhost:11434`)

| # | Method | Path | Body | ผลที่คาดหวัง |
|---|--------|------|------|--------------|
| TC1 | GET | `/api/version` | — | `200` → `{"version":"0.x.x"}` |
| TC2 | GET | `/v1/models` | — | `200` → `{"object":"list","data":[...]}` (มีโมเดลที่ pull แล้ว) |
| TC3 | POST | `/v1/embeddings` | `{"model":"nomic-embed-text","input":"hello"}` | `200` → `data[0].embedding` ยาว **768** ตัว |
| TC4 | GET | `/v1/embeddings` | — | `405 method not allowed` (รับแค่ POST) |
| TC5 | POST | `/v1/embeddings` | `{"model":"nomic-embed-text","input":""}` | `200` → embedding ของ empty string (Ollama ไม่ error) |
| TC6 | POST | `/v1/embeddings` | `{"model":"โมเดลที่ไม่มี"}` | `400` → error `model "<name>" not found` |
| TC7 | GET | `/v1/chat/completions` | — | `405` (รับแค่ POST) |
| TC8 | POST | `/v1/chat/completions` | `{"model":"qwen3:4b","messages":[{"role":"user","content":"hi"}]}` | `200` → `choices[0].message.content` |

**คำสั่งรัน (PowerShell — ใช้ไฟล์ body เพื่อเลี่ยง quote):**
```powershell
Set-Content -LiteralPath "$env:TEMP\embody.json" -Value '{"model":"nomic-embed-text","input":"hello"}' -Encoding ascii
curl.exe -s http://localhost:11434/api/version
curl.exe -s http://localhost:11434/v1/models
curl.exe -s -o NUL -w "%{http_code}`n" http://localhost:11434/v1/embeddings          # คาด 405
curl.exe -s http://localhost:11434/v1/embeddings -H "Content-Type: application/json" -d "@$env:TEMP\embody.json"
```

### 6.2 Ai service (ผ่าน `localhost:8001`)

| # | Method | Path | Body | ผลที่คาดหวัง |
|---|--------|------|------|--------------|
| TC9 | GET | `/health` | — | `200` → `{"status":"ok","model":"..."}` |
| TC10 | POST | `/v1/chat` | `{"messages":[{"role":"user","content":"hi"}]}` | `200` → `{"reply":"...","model":"qwen3:4b"}` |
| TC11 | POST | `/v1/chat` | `{"messages":[]}` | `400` → `{"error":"messages is required"}` |
| TC12 | POST | `/v1/generate` | `{"question":"อุณหภูมิเกินไหม?","context":"Sensor A: 42C"}` | `200` → `{"reply":"...","model":"..."}` |
| TC13 | POST | `/v1/generate` | `{}` | `400` → `{"error":"question is required"}` |

```powershell
curl.exe -s http://localhost:8001/health
curl.exe -s http://localhost:8001/v1/generate -H "Content-Type: application/json" -d '{"question":"hi"}'
```

### 6.3 Elasticsearch vector search (ผ่าน app `:5000`) — ต้อง login ก่อน

เตรียม token:
```powershell
$login = curl.exe -s http://localhost:5000/api/auth/login -H "Content-Type: application/json" -d '{"username":"root","password":"root"}'
# เอา access_token ไปใช้ในทุก request ด้านล่าง
```

| # | Method | Path | Body | ผลที่คาดหวัง |
|---|--------|------|------|--------------|
| TC14 | GET | `/api/elasticsearch/health` | — | `200` → cluster health ของ ES (ไม่ต้อง auth) |
| TC15 | POST | `/api/elasticsearch/create-index` | `{"dims":768}` | `200` → `{"status":"ok","message":"index ready"}` |
| TC16 | POST | `/api/elasticsearch/embeddings` | `{"document_id":"doc-1","content":"เซ็นเซอร์อุณหภูมิร้อน","source_type":"iot"}` | `200` → `{"indexed":true,"dim":768,"model":"nomic-embed-text",...}` |
| TC17 | POST | `/api/elasticsearch/embeddings` | `{"content":""}` | `400` → `content is required` |
| TC18 | POST | `/api/elasticsearch/vectors` | `{"document_id":"doc-2","content":"x","embedding":[0.1,0.2]}` | `200` → `{"indexed":true,"dim":2,...}` (ไม่เรียก LLM) |
| TC19 | POST | `/api/elasticsearch/search` | `{"query":"อุณหภูมิ","k":10}` | `200` → `hits` มี doc ที่ embed ไว้ (score cosine) |
| TC20 | POST | `/api/elasticsearch/search` | `{"query":""}` | `400` → `query is required` |
| TC21 | POST | `/api/elasticsearch/search` | ไม่ส่ง token | `401` (route ต้อง JWT) |

```powershell
curl.exe -s http://localhost:5000/api/elasticsearch/health
curl.exe -s http://localhost:5000/api/elasticsearch/embeddings -H "Content-Type: application/json" -H "Authorization: Bearer $token" -d '{"document_id":"doc-1","content":"เซ็นเซอร์อุณหภูมิร้อน","source_type":"iot"}'
curl.exe -s http://localhost:5000/api/elasticsearch/search -H "Content-Type: application/json" -H "Authorization: Bearer $token" -d '{"query":"อุณหภูมิ","k":10}'
```

> หมายเหตุ: ยังไม่มี API delete แบบผ่าน route — ถ้าจะลบ doc ใช้ ES REST ตรง ๆ (`DELETE /vector_embeddings/_doc/{id}`)

---

## 7. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| `llm returned empty embedding` | โมเดล embedding ไม่มี/ตอบ 0 มิติ → pull โมเดลที่ถูกต้อง |
| `llm request failed: ... connection refused` | Ollama ไม่รัน หรือ `LLM_BASE_URL` ผิด |
| `llm embeddings returned 404` | path ผิด — ต้องลงท้าย `/v1` (base_url มี `/v1`) |
| `/api/elasticsearch/*` ไม่มี route (404) | ES client หรือ LLM client เป็น nil → ตรวจ config + log ตอน boot |
| `index ... resource_already_exists_exception` | ปกติ (ระบบถือเป็น success) |
| `ai` container ไม่ start | Ollama ยังไม่ healthy → รอ 15-30s |
| `/v1/chat` ช้า/timeout | cold start ของโมเดล; `AI_TIMEOUT` เพิ่มได้ |
| embed dim ไม่ตรงกับ index | เปลี่ยนโมเดล embedding แล้วลบ index เดิม แล้ว `create-index` ใหม่ |
| `docker compose up` ไม่ pull ai | ใช้ `docker compose build ai` ก่อนถ้าอิมเมจไม่ทัน |
