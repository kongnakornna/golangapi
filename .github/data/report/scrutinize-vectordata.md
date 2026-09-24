# Scrutiny Report — โมดูล Vector Data (`internal/modules/vectordata` + `pkg/vectordb`)

> สถานะ: **fix-then-ship**
> วันที่: 2026-08-14
> ขอบเขต: abstraction layer `pkg/vectordb`, backend `pgvector` + `elasticsearch`, โมดูล `vectordata` (CRUD/search/seed), CLI `vectordata seed`, wiring ใน `internal/server/handlers.go` + `cmd/migrate.go`

---

## 1. Intent — โมดูลนี้พยายามทำอะไร

สร้างระบบ Vector Database ที่สลับ backend ได้ด้วย config (`pgvector` default / `elasticsearch`) พร้อม API CRUD + semantic search + seed ข้อมูลตัวอย่าง และ CLI seed — ให้โปรเจกต์มี vector search โดยไม่ผูกกับ backend ตัวใดตัวหนึ่ง

**คำถาม simpler-alternative:** ระบบนี้จำเป็นไหม? — มีเหตุผลรองรับบางส่วน แต่จุดอ่อนสำคัญคือมี **โมดูล ES อยู่ก่อนแล้ว** (`internal/modules/elasticsearch/usecase/usecase.go:14-20`) ซึ่งทำ LLM-embed + ES vector search อยู่แล้ว (EmbedAndIndex/IndexVector/SearchVector) ผ่าน `pkg/elasticsearch` ตัวเดียวกัน ตอนนี้มี **สองกองซ้อนที่ทำสิ่งเดียวกันบน ES** แต่มี behavior ต่างกัน (ดู Finding #3, #4) — เป็น maintenance/contract burden ที่ไม่จำเป็นต้องซ้อนกัน ทางเลือกที่สวยกว่า: ให้โมดูล ES เดิม delegating ผ่าน `vectordb` เดียวกัน หรือ deprecate โมดูล ES เดิม แล้วให้ `vectordata` เป็น single path ตัวเดียว

---

## 2. Trace — เดินตามเส้นทางจริง

### 2.1 Wiring
`internal/server/handlers.go:214` → `vectordb.New(&cfg.VectorDB, esClient, db, logger)` → `handlers.go:218` `NewVectorDataUseCase(vdb, llmClient, cfg, logger)` → `handlers.go:220` `MapVectorDataRoutes` ใต้ `/api/vectordata`; `handlers.go:353` health flag ใช้ `vdb != nil`

### 2.2 pgvector path (default)
`POST /documents` → `handler.Create` (`delivery/http/handler.go:80`) → `usecase.Create` (`usecase.go:67`) → ถ้าไม่มี embedding เรียก `llm.Embed` → `vdb.EnsureIndex(len(embedding))` → `pgBackend.Index` (`pg_backend.go:100`) → SQL `INSERT ... ?::vector RETURNING id`

Search: `handler.Search` → `usecase.Search` (`usecase.go:186`) → `llm.Embed(query)` → `pgBackend.SearchKNN` (`pg_backend.go:191`) → SQL `1 - (embedding <=> ?::vector) ORDER BY embedding <=> ?::vector`

### 2.3 ES path
`vectordb.New` provider=elasticsearch → `newESBackend` (`es_backend.go:19`) → ทุก method delegate ไป `b.client.*` ซึ่งใช้ `client.index` ที่ตั้งจาก **`cfg.Elasticsearch.Index`** (`pkg/elasticsearch/client.go:58`) — **ไม่ใช่** `cfg.VectorDB.Index`

### 2.4 CLI seed
`cmd/vectordata.go:73` → `vdUC.Seed(ctx, nil)` → `usecase.Seed` (`usecase.go:220`) — default ใหม่ = `len(sampleDocs)` (5) ตามที่แก้แล้ว

---

## 3. Verify — ยืนยันกับสิ่งที่อ้างไว้

อ้างอิงการทดสอบจริงทั้งหมด (live บน `localhost:5000`, DB `icmongolang`, Ollama `nomic-embed-text`):

---

## 4. Findings (เรียงตาม severity)

### 🔴 BLOCKER-1: ES provider **ไม่ฟัง** `VECTOR_DB_INDEX` — ข้อมูลไปลง index ผิดตัว

- **Finding:** `esBackend.index` เป็น dead code — ตั้งค่าจาก `cfg.VectorDB.Index` (`es_backend.go:20-24`) แต่ **ไม่ถูกใช้ใน method ใดเลย**; ทุก operation เรียก `b.client.*` ซึ่งใช้ `client.index` = `cfg.Elasticsearch.Index` (`client.go:58`) เสมอ
- **Why:** ทำลายสัญญาหลักของโมดูล "สลับ backend ด้วย config" และทำให้ `README_VectorData.md:126` (อ้างว่า `VECTOR_DB_INDEX` ใช้กับ ES ได้) ผิด
- **Evidence:** รัน CLI ด้วย `VECTOR_DB_INDEX=vector_documents_test` → log บอก `Seeded into "vector_documents_test"` แต่ ES log บอก `index "vector_embeddings" ready`; ลบ `vector_documents_test` → 404 `index_not_found`, ลบ `vector_embeddings` → acknowledged (ข้อมูลไปอยู่ที่ `vector_embeddings` จริง)
- **Fix:** ให้ `pkg/elasticsearch` รับ index เป็น parameter ต่อ method หรือสร้าง client ตัวแยกสำหรับ vectordb ด้วย `cfg.VectorDB.Index` แล้วตรวจสอบ/แก้ README §3.2

### 🟠 MAJOR-2: Partial update พังทั้งสองแบบ (wipes ข้อมูล + 500)

- **Finding:** `usecase.Update` (`usecase.go:143-175`) ทำ **replace แบบ all-fields** แต่ handler รับ payload บางส่วนได้โดยไม่มี validation (`UpdateDocumentRequest` ไม่มี `validate:"required"`)
- **Why:** (ก) PUT ที่ส่งแค่ `content` → ล้าง `document_id`/`source_type`/`metadata` เป็น `""` (data loss); (ข) PUT ที่ส่งแค่ `metadata`/`source_type` → embedding เป็น nil → `vecToString(nil)` = `"[]"` → SQL `'[]'::vector` error → 500
- **Evidence:**
  - `PUT /documents/2` `{"content":"..."}` → 200 → `GET /documents/2` คืน `document_id:""` (ข้อมูลหาย)
  - `PUT /documents/2` `{"metadata":"only metadata change"}` → 500 `ERROR: vector must have at least 1 dimension` (ยืนยันผ่าน psql: `SELECT '[]'::vector` → error เดียวกัน)
  - `PUT /documents/999999` → 200 `{"indexed":true}` ทั้งที่ไม่มี row (ดู Finding #3)
- **Fix:** Fetch เอกสารเดิมก่อน แล้ว merge เฉพาะ field ที่ส่งมา (หรือบังคับ payload เต็ม + 400 ถ้าขาด); ถ้า content และ embedding ว่าง ให้คง embedding เดิมไว้ ไม่ส่ง `[]`

### 🟠 MAJOR-3: Error semantics ต่างกันระหว่าง provider + ถูก handler ปิดบัง

- **Finding:** (ก) pgvector `UPDATE`/`DELETE` ไม่ตรวจ `RowsAffected` → ทำบน id ที่ไม่มี คืน success; ES คืน 404→500; (ข) handler.Get (`handler.go:110-114`) map **ทุก error → 404** รวมถึง DB error จริง; (ค) `GET /documents/abc` (id ไม่ใช่ตัวเลข) ผ่าน error ของ PG cast ออกมาเป็น body 404
- **Why:** เรียก API เดียวกันบน backend คนละตัวได้ผลต่างกัน ทำลาย "สลับ backend" และทำให้ทีม debug ผิดทาง (DB พังถูกมองเป็น not-found)
- **Evidence:**
  - pg: `DELETE /documents/999999` → 200 `{"status":"ok"}`; `PUT /documents/999999` (payload เต็ม) → 200
  - ES: update/delete id ที่ไม่มี → `404 index_not_found` → handler ตอบ 500
  - `GET /documents/abc` → body `ERROR: invalid input syntax for type bigint` (SQLSTATE 22P02)
- **Fix:** ใช้ sentinel error `ErrNotFound` ตรวจ `RowsAffected` ใน pg_backend; handler map เฉพาะ sentinel → 404, อย่างอื่น → 500; validate id เป็นตัวเลขก่อนส่ง SQL

### 🟠 MAJOR-4: ES backend **ทิ้ง `metadata` ตลอด**

- **Finding:** `esclient.VectorDoc` (`client.go:18-25`) **ไม่มี field `Metadata`** → `esBackend.Index/Get/List/Update` ไม่มีทางเก็บ metadata ได้ (`es_backend.go:40-106`)
- **Why:** Cross-provider contract แตก — pgvector เก็บ metadata ได้ครบ, ES เก็บไม่ได้เลย (เงียบ ๆ ไม่ error); ผู้ใช้ย้าย backend แล้วข้อมูลหาย
- **Evidence:** อ่าน struct `VectorDoc` ใน `pkg/elasticsearch/client.go:18-25` ไม่พบ `metadata`; `es_backend.go` ทุก method สร้าง `esclient.VectorDoc` โดยไม่ map `doc.Metadata`
- **Fix:** เพิ่ม `Metadata string json:"metadata,omitempty"` ใน `esclient.VectorDoc` + mapping ใน `esBackend.Index`/`Update`/`Get`/`List` (ทั้ง index mapping และ ES Update doc ต้องรวม field ด้วย)

### 🟠 MAJOR-5: Seed ไม่ idempotent — รันซ้ำสะสม duplicates

- **Finding:** `usecase.Seed` (`usecase.go:220-252`) insert ไม่มีเงื่อนไข ไม่ dedup ตาม `document_id`
- **Why:** `README_VectorData.md:307` แนะนำ "ทำ POST /seed หรือ re-index ใหม่" เหมือนรันซ้ำได้ปลอดภัย แต่ความจริงรันซ้ำทุกครั้งได้ docs ซ้ำ (observed: `doc-temp-sensor` มี 3 rows)
- **Evidence:** ระหว่างเทสต์ รัน HTTP seed + CLI seed → ตารางมี `doc-temp-sensor` 3 row (psql `SELECT ... WHERE document_id='doc-temp-sensor'`)
- **Fix:** ใน Seed ให้ `DELETE`/ข้ามถ้ามี `document_id` เดิมอยู่แล้ว (หรือ document ว่าต้องล้างก่อน seed)

### 🟡 MINOR-6: `cmd/migrate.go` สั่ง `CREATE EXTENSION IF NOT EXISTS "vector"` ไว้ unconditionally

- **Finding:** `createExtensions` (`migrate.go:69`) รัน extension `vector` เสมอ แม้ provider เป็น elasticsearch
- **Why:** deployment ที่ใช้ ES อย่างเดียวกับ postgres image ธรรมดา (ไม่มี pgvector) จะเจอ warning ตอน migrate โดยไม่จำเป็น (ไม่ fatal เพราะเป็น Warnf แต่สร้าง noise/confusion)
- **Fix:** ย้าย `CREATE EXTENSION vector` ไปไว้ในเงื่อนไข `provider == "pgvector"` หรือเงื่อนไขเดียวกับ EnsureIndex

### 🟡 MINOR-7: ไม่มี upper bound สำหรับ `k` / `limit` จากผู้ใช้

- **Finding:** `SearchRequest.K` และ `limit` จาก query parameter ผ่านเข้า SQL `LIMIT ?` / ES `num_candidates=k*10` ตรง ๆ ไม่มี cap (`pg_backend.go:191-193`, `client.go:176-183`)
- **Why:** `k` ใหญ่ (เช่น 1e6) → ES `num_candidates` เกิน max (10000) → 500; pg `LIMIT` ใหญ่ → ค่าใช้จ่าย query สูง; เป็นช่อง DoS เบา ๆ
- **Fix:** cap `k`/`limit` (เช่น max 100) ที่ handler หรือ usecase

### ⚪ NIT-8: `Create` บังคับ `content` แม้ส่ง `embedding` มาด้วย

- **Finding:** `handler.Create` (`handler.go:86-89`) ปฏิเสธทุก request ที่ `content` ว่าง — ทำให้เก็บ embedding-only / metadata-only doc ไม่ได้ทั้งที่ interface รองรับ
- **Why:** ผู้ใช้ที่มี vector อยู่แล้ว (เช่น ย้ายข้อมูล) ยังต้องสมมติ content ปลอมขึ้นมา
- **Fix:** อนุญาต content ว่างเมื่อมี `embedding` ไม่ว่าง

### ⚪ NIT-9: ไม่มี automated tests เลย

- **Finding:** `glob internal/modules/vectordata/**/*_test.go` และ `pkg/vectordb/*_test.go` → ไม่มีไฟล์
- **Why:** จุดที่ละเอียดอ่อนสุด (GORM `Raw` + `?::vector` binding, `RowsAffected`, empty-vector edge, ES/ID mapping) ถูกยืนยันด้วยการเทสต์ manual เท่านั้น — TC1-21 ใน README เป็น manual checklist
- **Fix:** เพิ่ม unit test อย่างน้อยสำหรับ `vecToString` (รวม empty → ปัจจุบันคืน `"[]"` ซึ่งเป็นต้นตอ Finding #2), `pgBackend.Update` (RowsAffected), และ usecase.Update (merge พฤติกรรม)

---

## 5. สิ่งที่ยืนยันแล้วว่าทำงานถูกต้อง (ไม่ใช่ปัญหา)

- CRUD + search + seed ทำงาน end-to-end บน pgvector (health, create-index, documents, list, get, update เต็ม payload, delete, search, seed — เทสต์จริงผ่าน)
- CLI `vectordata seed` seed ครบ 5 docs แล้ว (แก้ default เรียบร้อย)
- Build / vet / gofmt ผ่าน
- Ollama คืน vector เหมือนกันสำหรับ Thai ล้วน = ข้อจำกัดของ model (documented ใน README) — **ไม่ใช่บั๊กของโค้ด**

---

## 6. Verdict

**fix-then-ship** — เหตุผลหลัก: ES provider ทำงานผิดสัญญา (ignore `VECTOR_DB_INDEX` + ทิ้ง `metadata`) ซึ่งคือหัวใจของ abstraction ที่โมดูลนี้สร้างมาเพื่อแก้ปัญหา และ partial update มี data loss + 500 จริงที่ยืนยันได้บนระบบที่รันอยู่
