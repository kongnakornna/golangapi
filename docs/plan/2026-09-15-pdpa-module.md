# แผนพัฒนา PDPA Module

- **โจทย์**: พัฒนา PDPA Module ใหม่ภายใน `internal/modules/pdpa/` ด้วย Clean Architecture + DDD + Event-Driven Architecture (EDA) บน Go backend
- **Ticket**: PDPA
- **วันที่รวบรวมแผน**: 2026-09-15
- **หมายเหตุ**: แผนนี้จัดเรียงเป็น Flow ของการทำงาน (Foundation → Domain → Infrastructure → Application → Interfaces) ส่วนเลข `PDPA-N` เป็นหมายเลขรายการย่อยเพื่อใช้อ้างอิงในงานจริง

---

## ภาพรวมการแก้ไข

| เลข | รายการงาน | เลเยอร์ | ไฟล์หลัก |
|-----|-----------|---------|----------|
| PDPA-1 | สร้าง DDL + Seed ฐานข้อมูล | Foundation / Infra | `migrations/001_initial_pdpa_schema.sql`, `migrations/002_seed_initial_policy.sql` |
| PDPA-2 | สร้าง Domain Entities / Aggregates | Domain | `internal/modules/pdpa/domain/entity/**` |
| PDPA-3 | สร้าง Value Objects | Domain | `internal/modules/pdpa/domain/value_object/**` |
| PDPA-4 | สร้าง Domain Events | Domain | `internal/modules/pdpa/domain/event/**` |
| PDPA-5 | สร้าง Repository Interfaces | Domain | `internal/modules/pdpa/domain/repository/**` |
| PDPA-6 | สร้าง Domain Services | Domain | `internal/modules/pdpa/domain/service/**` |
| PDPA-7 | สร้าง Domain Errors (sentinel) | Domain | `internal/modules/pdpa/domain/errors/**` |
| PDPA-8 | Implement Repository: Consent | Infrastructure | `internal/modules/pdpa/infrastructure/persistence/postgres/**` |
| PDPA-9 | Implement Repository: DSAR Request | Infrastructure | `internal/modules/pdpa/infrastructure/persistence/postgres/**` |
| PDPA-10 | Implement Repository: Policy / AuditTrail / UserAccountStatus | Infrastructure | `internal/modules/pdpa/infrastructure/persistence/postgres/**` |
| PDPA-11 | Outbox Repository + Dispatcher | Infrastructure | `internal/modules/pdpa/infrastructure/persistence/postgres/outbox_repo_impl.go` (NEW) |
| PDPA-12 | Kafka Producer | Infrastructure | `internal/modules/pdpa/infrastructure/messaging/kafka/producer.go` |
| PDPA-13 | Redis Cache | Infrastructure | `internal/modules/pdpa/infrastructure/cache/redis/**` |
| PDPA-14 | External Clients | Infrastructure | `internal/modules/pdpa/infrastructure/external/**` |
| PDPA-15 | Application Ports | Application | `internal/modules/pdpa/application/port/**` (NEW) |
| PDPA-16 | Consent Commands | Application | `internal/modules/pdpa/application/command/**` |
| PDPA-17 | DSAR Commands | Application | `internal/modules/pdpa/application/command/**` |
| PDPA-18 | Account Commands | Application | `internal/modules/pdpa/application/command/**` |
| PDPA-19 | Policy + Admin Commands | Application | `internal/modules/pdpa/application/command/**` |
| PDPA-20 | Queries | Application | `internal/modules/pdpa/application/query/**` |
| PDPA-21 | DTOs | Application | `internal/modules/pdpa/application/dto/**` |
| PDPA-22 | Event Handlers | Application | `internal/modules/pdpa/application/event_handler/**` (NEW) |
| PDPA-23 | HTTP Handlers + Routes | Interfaces | `internal/modules/pdpa/interfaces/**` + `interfaces/http/routes.go` |
| PDPA-24 | Module Wiring + Background | Interfaces / Entry | `internal/modules/pdpa/interfaces/http/module.go` |

---

## ส่วนที่ 1 — Foundation

### PDPA-1: Migrations (DDL + Seed)

- **Location**: `migrations/`
- **จาก → เป็น**: ยังไม่มี → สร้าง `migrations/001_initial_pdpa_schema.sql` (ทำ DDL) และ `migrations/002_seed_initial_policy.sql` (ทำ seed)

**รายละเอียด**

- สร้างตารางหลักทั้งหมดของ PDPA Module:
  - `pdpa_consents`
  - `pdpa_dsar_requests`
  - `pdpa_user_account_statuses`
  - `pdpa_audit_trails`
  - `pdpa_privacy_policies`
  - `pdpa_outbox`
  - `pdpa_processed_events`
  - `pdpa_purposes`
- ออกแบบ Index + FK ให้สอดคล้องกับ Query ที่กำหนดไว้
- `pdpa_outbox` ต้องมีครบ field ที่ใช้ใน Outbox pattern เช่น `id`, `event_type`, `payload`, `status`, `retry_count`, `created_at`
- `pdpa_processed_events` ใช้อนุกรมการประมวลผล Event แบบ idempotent
- Seed ข้อมูล Policy เริ่มต้น (initial privacy policy) ไว้ใน `002`
- **Acceptance**: `make migrate` หรือคำสั่ง migration ตามสคริปต์ของโปรเจกต์รันผ่าน ไม่มี conflict ของตารางซ้ำ

---

## ส่วนที่ 2 — Domain

### PDPA-2: Domain Entities / Aggregates

- **Location**: `internal/modules/pdpa/domain/entity/**`
- **จาก → เป็น**: ไม่มี → สร้าง Entity และ Aggregate ของ PDPA

**รายละเอียด**

- Entity หลัก: `Consent` (aggregate root), `DSARRequest`, `PrivacyPolicy`, `UserAccountStatus`, `AuditTrail` (ถ้าจำเป็นใน domain)
- Aggregate `Consent` เป็นตัวเก็บธุรกรรมหลัก: entity ตัวเดียว + 1 row ใน `pdpa_outbox` ภายใน transaction เดียว
- domain ห้าม import ตาม CMD 7: **ห้ามใช้ `gorm`, `gin`, `chi`, `sarama`, `redis`, `elasticsearch` ใน layer domain**

### PDPA-3: Value Objects

- **Location**: `internal/modules/pdpa/domain/value_object/**`
- **จาก → เป็น**: ไม่มี → สร้าง Value Objects

**รายละเอียด**

- Value Objects เท่าที่จำเป็น เช่น `Purpose`, `ConsentStatus`, `DSARStatus`, `ContactChannel`, `PolicyVersion` เป็นต้น
- เน้น immutable + method ประกอบ business logic ของ value เอง

### PDPA-4: Domain Events

- **Location**: `internal/modules/pdpa/domain/event/**`
- **จาก → เป็น**: ไม่มี → สร้าง Domain Events

**รายละเอียด**

- Event ที่เกิดจากการทำงานหลัก เช่น Consent ที่บันทึก, DSAR ที่ Submit/Process, สถานะบัญชีที่เปลี่ยน (suspend/terminate/confirm deletion), Policy ที่ publish
- Event ต้องบรรจุ payload พอที่ consumer จะทำงานต่อได้ โดยไม่ต้องไปอ่าน DB ใหม่ (ตามหลัก EDA)

### PDPA-5: Repository Interfaces

- **Location**: `internal/modules/pdpa/domain/repository/**`
- **จาก → เป็น**: ไม่มี → สร้าง interface ของ repository ทุกตัวที่ domain ต้องการ

**รายละเอียด**

- ประกาศ interface อย่างน้อย: `ConsentRepository`, `DSARRequestRepository`, `PolicyRepository`, `AuditTrailRepository`, `UserAccountStatusRepository`, `OutboxRepository`
- interface ให้อยู่ฝั่ง domain เพราะ domain เป็นเจ้านายของ contract (ในแง่นี้สอดคล้องกับทิศทาง dependency ของโปรเจกต์)

### PDPA-6: Domain Services

- **Location**: `internal/modules/pdpa/domain/service/**`
- **จาก → เป็น**: ไม่มี → สร้าง Domain Services

**รายละเอียด**

- DomainService กลายเป็นตัวกลาง: command → เรียก DomainService → mutate aggregate → เก็บ entity + เขียน `pdpa_outbox` ใน tx เดียว
- ใส่ business rules: ความถูกต้องของ consent, การเทียบ purpose, การ hash ข้อมูลส่วนบุคคลด้วย `anon_salt` (ค่าจาก secret)
- Command **ห้าม** เรียก `producer.PublishMessage` ตรงๆ (CMD 4, CMD 12) — ให้ DomainService/Outbox เป็นคนจัดการแทน

### PDPA-7: Domain Errors

- **Location**: `internal/modules/pdpa/domain/errors/**`
- **จาก → เป็น**: ไม่มี → สร้าง sentinel errors

**รายละเอียด**

- สร้างข้อผิดพลาดแบบ sentinel ตาม CMD 10 เช่น `ErrConsentNotFound`, `ErrDSARNotFound`, `ErrPolicyNotFound`, `ErrInvalidConsentStatus` เป็นต้น
- มาแบบ `errors.New` + `errors.Is/As` ใช้เทียบที่ usecase/handler

---

## ส่วนที่ 3 — Infrastructure

### PDPA-8: Consent Repository Impl

- **Location**: `internal/modules/pdpa/infrastructure/persistence/postgres/**`
- **จาก → เป็น**: ไม่มี → implement `ConsentRepository`

**รายละเอียด**

- Implement ให้ครบตาม interface จาก PDPA-5
- ใช้วิธีการ persist แบบที่โปรเจกต์ใช้อยู่เดิม (พยายามไม่ใช้ gorm ถ้า domain rule ห้าม; ถ้าโปรเจกต์ใช้ gorm ที่อื่น ให้คงรูปแบบเดิม)
- ตรวจสอบให้ผ่านการ rule CMD 1–6 ที่โปรเจกต์ใช้

### PDPA-9: DSAR Request Repository Impl

- **Location**: `internal/modules/pdpa/infrastructure/persistence/postgres/**`
- **จาก → เป็น**: ไม่มี → implement `DSARRequestRepository`

**รายละเอียด**

- CRUD + flow สถานะของ DSAR (submit → process → close/complete)
- รองรับ Query สำหรับ `GET /api/v1/pdpa/dsar` (list)

### PDPA-10: Policy / AuditTrail / UserAccountStatus Repository Impls

- **Location**: `internal/modules/pdpa/infrastructure/persistence/postgres/**`
- **จาก → เป็น**: ไม่มี → implement repository ที่เหลือ (`PolicyRepository`, `AuditTrailRepository`, `UserAccountStatusRepository`)

**รายละเอียด**

- Policy: รองรับ getActive / listVersions / publish
- AuditTrail: เขียน log การกระทำ + query สำหรับ `GetReport` / `ListAudit`
- UserAccountStatus: ใช้ใน suspend / terminate / confirm deletion

### PDPA-11: Outbox Repository + Dispatcher

- **Location**: `internal/modules/pdpa/infrastructure/persistence/postgres/outbox_repo_impl.go` (NEW)
- **จาก → เป็น**: ไม่มี → สร้าง outbox repo impl + กลไก dispatch

**รายละเอียด**

- Flow มาตรฐาน: command → DomainService → เขียน entity + 1 row `pdpa_outbox` ใน **transaction เดียว** → `OutboxDispatcher` (polling) ติดตาม row ที่ยังไม่ถูกส่ง → ส่งเข้า Kafka producer → อัปเดตสถานะเป็น success / ack; ถ้าล้มเหลว → retry (นับ `retry_count`) → สุดท้ายเข้า DLQ
- กัน duplicate event ด้วย `pdpa_processed_events`
- **Acceptance**: เมื่อ command สำเร็จ event จะถูกส่งออกไปยัง topic ตรง ไม่มีการเรียก producer จาก command ตรงๆ

### PDPA-12: Kafka Producer

- **Location**: `internal/modules/pdpa/infrastructure/messaging/kafka/producer.go`
- **จาก → เป็น**: ไม่มี → สร้าง Kafka producer

**รายละเอียด**

- ใช้ **13 topic** สำหรับ PDPA Module (ปรับชื่อก่อน implement — ดูหัวข้อ Open Decisions ข้อ 3)
- `producer.PublishMessage` ถูกเรียกจาก OutboxDispatcher / event_handler เท่านั้น ไม่ใช่จาก command
- implement ให้รองรับ side effect: ack / fail / DLQ
- domain layer ห้าม import `sarama` (CMD 7)

### PDPA-13: Redis Cache

- **Location**: `internal/modules/pdpa/infrastructure/cache/redis/**`
- **จาก → เป็น**: ไม่มี → สร้าง cache ฝั่ง read

**รายละเอียด**

- ใช้ทำ cache ของข้อมูลที่อ่านบ่อย เช่น policy ที่ active / consent สรุป
- domain layer ห้าม import `redis` (CMD 7)

### PDPA-14: External Clients

- **Location**: `internal/modules/pdpa/infrastructure/external/**`
- **จาก → เป็น**: ไม่มี → สร้าง client ติดต่อระบบภายนอก

**รายละเอียด**

- client ที่ PDPA ต้องเรียก เช่น ระบบ user/account (ถ้ามี) หรือระบบที่เกี่ยวข้องกับ DSAR
- ครอบด้วย interface เพื่อให้ test ได้ง่าย

---

## ส่วนที่ 4 — Application

### PDPA-15: Application Ports

- **Location**: `internal/modules/pdpa/application/port/**` (NEW)
- **จาก → เป็น**: ไม่มี → สร้าง port interfaces

**รายละเอียด**

- สร้าง interface ที่ application ต้องการจากภายนอก (repo, producer, cache, outbox, external client)
- ทิศทาง dependency ยังชี้เข้า domain ตามธรรมเนียมของโปรเจกต์

### PDPA-16: Consent Commands

- **Location**: `internal/modules/pdpa/application/command/**`
- **จาก → เป็น**: ไม่มี → สร้าง command กลุ่ม consent (เช่น `record_consent` / `create_consent` — ดู Open Decisions ข้อ 2)

**รายละเอียด**

- command เป็น entry point ของ application ฝั่ง write (รวมทั้งหมดของทั้ง module มี 10 command)
- เรียก DomainService ให้เลือก aggregate เอง และ **ห้าม** เรียก producer ตรงๆ
- ตรวจสอบให้ command pass กฎ CMD 1–6, CMD 12

### PDPA-17: DSAR Commands

- **Location**: `internal/modules/pdpa/application/command/**`
- **จาก → เป็น**: ไม่มี → สร้าง command กลุ่ม DSAR

**รายละเอียด**

- เช่น `submit_dsar`, `process_dsar`, `close_dsar` (สอดคล้องกับ endpoint `POST /api/v1/pdpa/dsar` และ `POST /api/v1/pdpa/dsar/{id}/process`)
- เหมือนเดิม: เขียน entity + outbox ใน tx เดียว

### PDPA-18: Account Commands

- **Location**: `internal/modules/pdpa/application/command/**`
- **จาก → เป็น**: ไม่มี → สร้าง command กลุ่ม account status

**รายละเอียด**

- สอดคล้อง endpoint ของ AccountHandler: suspend / terminate / confirm deletion
- update `pdpa_user_account_statuses` + ปล่อย event ที่เกี่ยวข้อง

### PDPA-19: Policy + Admin Commands

- **Location**: `internal/modules/pdpa/application/command/**`
- **จาก → เป็น**: ไม่มี → สร้าง command กลุ่ม policy และ admin

**รายละเอียด**

- Policy: publish version ใหม่ (สอดคล้อง `Publish` ของ PolicyHandler)
- Admin: สร้างข้อมูล report / เรียกอะไรที่ต้องใช้สถานะ snapshot

### PDPA-20: Queries

- **Location**: `internal/modules/pdpa/application/query/**`
- **จาก → เป็น**: ไม่มี → สร้าง query (รวมทั้งหมด 5 query)

**รายละเอียด**

- query ฝั่ง read เช่น รายการ DSAR, รายการ audit, report สรุป, policy ฉบับ active
- **สำคัญ**: ก่อนลบ / เปลี่ยนชื่อ query file ใด ให้ยืนยันชุด 4 query ที่ระบุไว้ใน G10 ก่อน (ดู Open Decisions ข้อ 1)

### PDPA-21: DTOs

- **Location**: `internal/modules/pdpa/application/dto/**`
- **จาก → เป็น**: ไม่มี → สร้าง DTO (รวมทั้งหมด 13 ตัว)

**รายละเอียด**

- ครอบ request/response ที่ handlers ใช้ เพื่อไม่ให้ layer interface ผูกกับ entity ตรงๆ
- แยก DTO อินพุตกับเอาต์พุต

### PDPA-22: Event Handlers

- **Location**: `internal/modules/pdpa/application/event_handler/**` (NEW)
- **จาก → เป็น**: ไม่มี → สร้าง event handler

**รายละเอียด**

- รับ event ที่มาจาก domain / outbox และทำงานต่อ เช่น อัปเดต read model, ส่ง notification, เรียก external client
- ตรวจ idempotency ผ่าน `pdpa_processed_events`

---

## ส่วนที่ 5 — Interfaces (HTTP)

### PDPA-23: HTTP Handlers + Routes

- **Location**: `internal/modules/pdpa/interfaces/**` + `interfaces/http/routes.go`
- **จาก → เป็น**: ไม่มี → สร้าง handlers และลง routes

**รายละเอียด**

- สร้าง handler รวมทั้งหมด **10 ตัว** ตาม interface layer (รวม minimum ที่ระบุไว้ด้านล่าง)
- รูปแบบอ้างอิงจาก modules ที่มีอยู่แล้ว (`job`, `wos`, `items`, `purchaseorder`, `queue`, `websocket`) ซึ่งใช้ `usecase/repository/presenter/handler/delivery/http`

**Endpoints ที่ยืนยันแล้ว**

- **DSARHandler**
  - `POST /api/v1/pdpa/dsar` — submit DSAR
  - `GET /api/v1/pdpa/dsar` — list DSAR
  - `POST /api/v1/pdpa/dsar/{id}/process` — process DSAR (admin)
- **AccountHandler** — suspend / terminate / confirm deletion
- **PolicyHandler** — `GetActive` / `ListVersions` / `Publish`
- **AdminHandler** — `GetReport` / `ListAudit`
- helper: `extractUserID(r *http.Request) (uuid.UUID, error)` ใช้ดึง user id จาก request

**Routes**

- ลงทะเบียนทุก route ใน `interfaces/http/routes.go` ตาม prefix เดิมของโปรเจกต์

### PDPA-24: Module Wiring + Background

- **Location**: `internal/modules/pdpa/interfaces/http/module.go`
- **จาก → เป็น**: ไม่มี → สร้าง module wiring

**รายละเอียด**

- `Handlers` struct: รวมทุก handler
- `RegisterRoutes`: เรียกเพื่อลง route
- `Config` struct: config ของ PDPA module (รวม secret เช่น `anon_salt`)
- `NewModule`: constructor รวมทุก dependency (repo, service, producer, cache, handler)
- `StartBackground`: เริ่มงาน background เช่น `OutboxDispatcher`
- `Shutdown`: ปิดงาน background เรียบร้อย (graceful shutdown)

---

## ส่วนที่ 6 — Swagger

- **คำสั่งจากโจทย์** (ท้าย `Modules_plan_PDPA_Task.md`): **สร้าง swagger ทั้งหมดใน PDPA Module**
- สร้าง/อัปเดตไฟล์ swagger ให้ครบทุก endpoint ของ PDPA Module ตามรูปแบบที่โปรเจกต์ใช้อยู่เดิม (ทั้งใน handler comment หรือไฟล์ doc ตามที่โปรเจกต์ทำ)

---

## ส่วนที่ 7 — Unit Tests

- **ขอบเขต (confirmed)**: เขียนเทสเฉพาะ layer สำคัญ — **application + domain** เป็นหลัก
- **ไม่เขียน**: เทสฝั่ง infrastructure (Postgres repo, Kafka producer/outbox dispatcher, Redis cache, external clients) — ข้ามตามการตกลง แม้ว่า spec ราย subtask บางตัวจะเขียน `🧪 Unit: unit test & mockery` ก็ให้ถือว่า **ขอบเขตนี้แทนที่**
- ใช้ mockery / test double สำหรับ repo, producer, cache ที่ application เรียก
- **Acceptance**: coverage ≥ 80% (CMD 13) เฉพาะกับ layer ที่เทส

---

## ส่วนที่ 8 — กฎและเกณฑ์ผ่านงาน (Checks)

- **CMD 1–6**: ทิศทาง dependency ตามที่โปรเจกต์ใช้ (domain เป็นแกน, application → domain, infrastructure/interface เข้าหา application/domain)
- **CMD 7**: domain ห้าม import `gorm`, `gin`, `chi`, `sarama`, `redis`, `elasticsearch`
- **CMD 10**: sentinel errors + ใช้ `errors.Is/As`
- **CMD 12 / CMD 4**: command ห้ามเรียก `producer.PublishMessage` ตรงๆ — ผ่าน OutboxDispatcher/event_handler
- **CMD 13**: coverage ≥ 80%
- **CMD 15**: comment สองภาษา (ไทย + อังกฤษ) ในทุกโค้ดใหม่
- **Acceptance รวม**:
  - `terraform plan` (หรือคำสั่ง infra ที่โปรเจกต์ใช้) รันผ่าน
  - ตั้ง secret ครบรวม `anon_salt` (ความยาว 48) และ MSK credentials
  - `go vet` + `golangci-lint` สะอาด

---

## Open Decisions (รายการที่ต้องยืนยันก่อน implement)

1. **G10 — ชุด 4 query**: ก่อนลบ/เปลี่ยนชื่อ query file ใด ให้ยืนยันชุด query 4 ตัวที่ระบุไว้ใน G10 ก่อน
2. **G5 — ชื่อ command**: ยืนยันเจตนาระหว่าง `record_consent` กับ `create_consent` (ใช้ชื่อเดียว)
3. **G13 — ชื่อ Kafka topic**: ยืนยันชื่อทั้ง 13 topic ก่อนใส่ prefix
4. **ตำแหน่ง mocks / test doubles**: ยังไม่ชัดเจนว่า `application/mocks` จะอยู่ตรงไหน (ไม่พบ `application/mocks` ปัจจุบัน และ `repository/` ยังว่าง) — ระบุไว้ในแผนเป็น open decision ถึงแม้จะไม่บล็อกการเขียนแผน แต่เป็น **blocker ก่อนเริ่มเขียนโค้ด Phase 0**
5. **เอกสารอ้างอิงต้นทาง**: `docs/Modules_PDPA_0.md` และ `docs/Modules_PDPA_CMD.md` หาไม่พบที่ path คาดไว้ — แผนนี้สร้างจาก `Modules_plan_PDPA.md` + `Modules_plan_PDPA_Task.md` ที่อ่านได้ ณ วันที่รวบรวม ถ้ามีเอกสารอัปเดตให้ทบทวนแผนก่อน implement