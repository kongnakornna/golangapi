# Infrastructure & Utility Modules API

> Modules: Dashboard / Document / Email / Batch / I18n / WOS / Report / VectorData / Elasticsearch
> Auth: JWT Bearer (ยกเว้นที่ระบุว่า public)

---

## Response Envelope (dashboard, document, email, batch, i18n, wos)

```json
{ "data": <payload>, "is_success": true }
```
Error:
```json
{ "error": { "status": 404, "statusText": "not_found", "msg": "..." }, "is_success": false }
```
> ⚠️ vectordata / elasticsearch ใช้ error format ต่าง: `{ "error": "<message>" }`

---

## 1. Dashboard — `/api/dashboard`

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| GET | `/api/dashboard/stats` | `GetDashboardStats` | Auth |
| GET | `/api/dashboard/revenue` | `GetRevenueChart` | Auth |
| GET | `/api/dashboard/top-parts` | `GetTopParts` | Auth |
| GET | `/api/dashboard/job-status` | `GetJobStatusSummary` | Auth |

### GET `/api/dashboard/stats`
Response:
```json
{ "data": { "totalDevices": 120, "onlineDevices": 95, "activeAlerts": 7, "todayCommands": 345 }, "is_success": true }
```

### GET `/api/dashboard/revenue`
Query: `period` (optional, enum `daily|weekly|monthly`, default `monthly`)

### GET `/api/dashboard/top-parts`
Query: `limit` (handler hardcode `limit := 10` จริง — ไม่ได้อ่าน query)

---

## 2. Document — `/api/documents`

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| POST | `/api/documents/` | `Upload` | Auth |
| GET | `/api/documents/` | `List` | Auth |
| GET | `/api/documents/{id}` | `Download` | Auth |
| DELETE | `/api/documents/{id}` | `Delete` | Auth |

### POST `/api/documents/`
Body (`DocumentUploadRequest`, ทั้ง 3 required):
```json
{ "filename": "report.pdf", "mimeType": "application/pdf", "size": 2048 }
```
Response (echo ของ request — handler ยังไม่ได้เรียก use case จริง)

### GET `/api/documents/{id}`
Response `DocumentResponse`: `id`, `filename`, `originalName`, `mimeType`, `size`, `createdAt`

---

## 3. Email — `/api/email`

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| POST | `/api/email/send` | `SendEmail` | Auth |
| GET | `/api/email/logs` | `ListEmailLogs` | Auth |
| GET | `/api/email/logs/{id}` | `GetEmailLog` | Auth |
| GET | `/api/email/config` | `GetEmailConfig` | Auth |
| PUT | `/api/email/config` | `UpdateEmailConfig` | Auth |

### POST `/api/email/send`
Body (`EmailSendRequest`): `to` (required, email), `cc`, `bcc`, `subject` (required), `body` (required)
```json
{ "to": "user@example.com", "cc": "", "bcc": "", "subject": "แจ้งเตือน", "body": "เนื้อหา" }
```

### GET `/api/email/config`
Response `EmailConfigResponse`:
```json
{
  "data": {
    "smtpHost": "smtp.example.com", "smtpPort": 587, "smtpUser": "noreply",
    "fromEmail": "noreply@example.com", "fromName": "ICMono", "isActive": true
  },
  "is_success": true
}
```

---

## 4. Batch — `/api/batch`

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| POST | `/api/batch/jobs` | `CreateJob` | Auth |
| GET | `/api/batch/jobs` | `ListJobs` | Auth |
| GET | `/api/batch/jobs/{id}` | `GetJob` | Auth |
| PUT | `/api/batch/jobs/{id}` | `UpdateJob` | Auth |
| DELETE | `/api/batch/jobs/{id}` | `DeleteJob` | Auth |
| POST | `/api/batch/jobs/{id}/run` | `RunJobNow` | Auth |
| GET | `/api/batch/jobs/{id}/logs` | `GetJobLogs` | Auth |

### POST `/api/batch/jobs`
Body (`BatchJobRequest`): `name` (required), `type` (required), `config`, `schedule`
```json
{ "name": "Nightly Sync", "type": "sync", "config": "{\"source\":\"api\"}", "schedule": "0 2 * * *" }
```

### GET `/api/batch/jobs/{id}`
Response `BatchJobResponse`:
```json
{
  "data": {
    "id": "...", "name": "Nightly Sync", "type": "sync", "status": "completed",
    "config": "{\"source\":\"api\"}", "schedule": "0 2 * * *",
    "totalCount": 1000, "successCount": 998, "failCount": 2,
    "startedAt": "...", "finishedAt": "...", "createdAt": "..."
  },
  "is_success": true
}
```

---

## 5. I18n — `/api/i18n`

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| GET | `/api/i18n/translations` | `GetTranslations` | Auth |
| POST | `/api/i18n/translations` | `CreateTranslation` | Auth |
| GET | `/api/i18n/translations/{key}` | `GetTranslationByKey` | Auth |
| PUT | `/api/i18n/translations/{key}` | `UpdateTranslation` | Auth |
| DELETE | `/api/i18n/translations/{key}` | `DeleteTranslation` | Auth |

### GET `/api/i18n/translations`
Query: `locale` (optional) — ถ้ามี filter ตาม locale, ไม่มี → ทั้งหมด (limit 1000)
Response `TranslationResponse[]`: `id`, `locale`, `key`, `value`, `createdAt`, `updatedAt`

### GET `/api/i18n/translations/{key}`
Query: `locale` (**required** — ถ้าไม่มี → 400)

### POST `/api/i18n/translations`
Body (`TranslationRequest`): `locale` (required, 2–10 chars), `key` (required, 1–255), `value` (required)
```json
{ "locale": "th", "key": "greeting", "value": "สวัสดี" }
```

---

## 6. WOS (Web Order System) — `/api/wos`

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| POST | `/api/wos/orders` | `CreateOrder` | Auth |
| GET | `/api/wos/orders` | `ListOrders` | Auth |
| GET | `/api/wos/orders/{id}` | `GetOrder` | Auth |
| PUT | `/api/wos/orders/{id}/status` | `UpdateOrderStatus` | Auth |

### POST `/api/wos/orders`
Body (`WosOrderRequest`): `customerName` (required), `customerEmail` (required, email), `customerPhone`, `items` (required, JSON), `totalAmount` (required, float), `notes`
```json
{
  "customerName": "สมชาย", "customerEmail": "somchai@example.com", "customerPhone": "081-234-5678",
  "items": [ { "sku": "OIL-01", "qty": 2 } ], "totalAmount": 1000.0, "notes": ""
}
```

### GET `/api/wos/orders/{id}`
Response `WosOrderResponse`: `id`, `orderNumber`, `customerName`, `customerEmail`, `items`, `totalAmount`, `status`, `createdAt`, `updatedAt`

### PUT `/api/wos/orders/{id}/status`
Query: `status` (**required** — enum `pending|confirmed|shipped|delivered|cancelled`); ไม่มี → 400

---

## 7. Report — `/reports` (ไม่มี prefix `/api`)

> Access: PDF routes ใช้ `VerifierForReport()` — ยอมรับ **full access token** ใน header **หรือ** short-lived download token แบบ `?token=<download_token>`. `POST /reports/token` ใช้ **full access token เท่านั้น**

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| GET | `/reports/daily-sales/pdf` | `DailySalesPDF` | full หรือ download token |
| GET | `/reports/inventory-summary/pdf` | `InventorySummaryPDF` | full หรือ download token |
| GET | `/reports/customer-list/pdf` | `CustomerListPDF` | full หรือ download token |
| GET | `/reports/invoice/pdf` | `InvoicePDF` | full หรือ download token |
| GET | `/reports/credit-note/pdf` | `CreditNotePDF` | full หรือ download token |
| GET | `/reports/debit-note/pdf` | `DebitNotePDF` | full หรือ download token |
| POST | `/reports/token` | `CreateDownloadToken` | full access token เท่านั้น |

### POST `/reports/token` — สร้าง download token (อายุ 5 นาที)
Body (`downloadTokenRequest`): `type` (required — enum `daily_sales|inventory_summary|customer_list|invoice|credit_note|debit_note`), `source` (optional; **required ถ้า type=invoice** — quotation id)
```json
{ "type": "invoice", "source": "9f8e2c1a-0000-4000-8000-000000000099" }
```
Response:
```json
{ "token": "<jwt-download-token>", "token_type": "download", "expires_in": 300, "type": "invoice" }
```

### GET `/reports/daily-sales/pdf`
Query: `date` (optional, format `YYYY-MM-DD`; default = วันนี้)
Response: PDF inline `daily_sales_YYYYMMDD.pdf`

### GET `/reports/invoice/pdf`
Query: `source` (**required** — quotation id, uuid); ไมมี → 400
Response: PDF inline `invoice_<InvoiceNo>.pdf`

> token ผูกกับ `type` เดียว (+ `source` สำหรับ invoice) — ข้ามใช้ไม่ได้ (`downloadScopeOK`)

---

## 8. VectorData — `/api/vectordata`

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| GET | `/api/vectordata/health` | `Health` | **Public** |
| POST | `/api/vectordata/create-index` | `CreateIndex` | Auth |
| POST | `/api/vectordata/documents` | `Create` | Auth |
| GET | `/api/vectordata/documents` | `List` | Auth |
| GET | `/api/vectordata/documents/{id}` | `Get` | Auth |
| PUT | `/api/vectordata/documents/{id}` | `Update` | Auth |
| DELETE | `/api/vectordata/documents/{id}` | `Delete` | Auth |
| POST | `/api/vectordata/search` | `Search` | Auth |
| POST | `/api/vectordata/seed` | `Seed` | Auth |

> Error format: `{ "error": "<message>" }` พร้อม status (400/404/500)

### GET `/api/vectordata/health` (public)
Response: `{ "provider": "...", "status": "ok" }`

### POST `/api/vectordata/create-index`
Body (`CreateIndexRequest`): `dims` (int)
```json
{ "dims": 768 }
```

### POST `/api/vectordata/documents`
Body: ต้องมี `content` หรือ `embedding` อย่างน้อยหนึ่ง → 400 ถ้าว่างทั้งคู่
```json
{ "document_id": "doc-001", "content": "คู่มือเซนเซอร์ IoT", "source_type": "manual", "metadata": "{\"cat\":\"iot\"}" }
```
Response `DocumentResponse`: `id`, `document_id`, `content`, `source_type`, `metadata`, `dim`, `indexed`

### GET `/api/vectordata/documents`
Query: `limit` (default 10, max 100), `offset` (default 0)

### POST `/api/vectordata/search`
Body (`SearchRequest`): `query` (**required**), `k` (int)
```json
{ "query": "เซนเซอร์ทำงานผิดปกติ", "k": 5 }
```
Response `SearchResponse`: `query`, `k`, `count`, `hits[]` (`document_id`, `content`, `source_type`, `metadata`, `score`)

### POST `/api/vectordata/seed`
Body (`SeedRequest`): `count` (int) → Response `{ "provider": "...", "index": "...", "seeded": 5 }`

---

## 9. Elasticsearch — `/api/elasticsearch`

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| GET | `/api/elasticsearch/health` | `Health` | **Public** |
| POST | `/api/elasticsearch/create-index` | `CreateIndex` | Auth |
| POST | `/api/elasticsearch/embeddings` | `EmbedAndIndex` | Auth |
| POST | `/api/elasticsearch/vectors` | `IndexVector` | Auth |
| POST | `/api/elasticsearch/search` | `SearchVector` | Auth |

> Error format: `{ "error": "..." }` (400/500)

### GET `/api/elasticsearch/health` (public)
Response: map ของ cluster health

### POST `/api/elasticsearch/create-index`
Body: `dims` (int) → `{ "status": "ok", "message": "index ready" }`

### POST `/api/elasticsearch/embeddings`
Body (`EmbeddingRequest`): `document_id`, `content` (**required** → 400 ถ้าว่าง), `source_type`, `metadata`
```json
{ "document_id": "doc-001", "content": "คู่มือระบบ", "source_type": "manual", "metadata": "{}" }
```
Response `EmbeddingResponse`: `document_id`, `content`, `source_type`, `index`, `model`, `dim`, `indexed`

### POST `/api/elasticsearch/vectors`
Body (`IndexVectorRequest`): `document_id`, `content` (**required**), `embedding` (**required**, []float32, ไม่ว่าง), `source_type`, `metadata`
```json
{ "document_id": "doc-002", "content": "เนื้อหา", "embedding": [0.1, 0.2, 0.3], "source_type": "manual", "metadata": "{}" }
```

### POST `/api/elasticsearch/search`
Body: `query` (**required** → 400 ถ้าว่าง), `k` (int)
```json
{ "query": "แจ้งเตือนเกินขีดจำกัด", "k": 5 }
```

---

## สรุป access

| Module | Public | Auth (full token) | Auth (full หรือ download token) |
|--------|--------|-------------------|----------------------------------|
| dashboard | — | ทั้งหมด | — |
| document | — | ทั้งหมด | — |
| email | — | ทั้งหมด | — |
| batch | — | ทั้งหมด | — |
| i18n | — | ทั้งหมด | — |
| wos | — | ทั้งหมด | — |
| report | — | `/reports/token` | `/reports/*/pdf` |
| vectordata | `/health` | ที่เหลือ | — |
| elasticsearch | `/health` | ที่เหลือ | — |

## หมายเหตุ stub/ไม่สมบูรณ์ (จาก code จริง)
- `document.Upload`, `batch.CreateJob`, `wos.CreateOrder`, `i18n.CreateTranslation`, `email.UpdateEmailConfig`, `batch.UpdateJob` ฯลฯ มี handler แค่ echo/stub
- `document.List`, `batch.ListJobs`, `email.ListEmailLogs`, `wos.ListOrders` ใช้ hardcode limit/offset (50, 0)
