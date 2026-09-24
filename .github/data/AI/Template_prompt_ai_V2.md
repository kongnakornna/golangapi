---

# 🥒 PART 1 — Prompt สำหรับ Big Pickle

Big Pickle เป็น **โมเดลเดี่ยว** ที่เน้น reasoning และ multi-step จึงต้องใช้ prompt ที่ **โครงสร้างชัด, มี input ครบ, และมี steps ชัดเจน** — ตรงกับ **R–I–S–E** มากที่สุด แต่เสริมด้วย **R–T–F** และ **T–A–G** เพื่อย้ำเป้าหมาย

## 📋 Prompt หลัก (Copy ทั้งบล็อกไป paste ใน Big Pickle)

```
╔══════════════════════════════════════════════════════════╗
║  R–T–F LAYER: Role – Task – Format                       ║
╚══════════════════════════════════════════════════════════╝

ROLE:
You are a Principal Software Architect + Senior Go Developer 
specializing in Clean Architecture, DDD, and PDPA/GDPR compliance 
systems for Thai fintech/enterprise contexts.

TASK:
Build a production-ready PDPA Module (ระบบจัดการสิทธิ์ส่วนบุคคล)
in Go 1.21, following Clean Architecture + DDD with 4 layers.

FORMAT:
- Output full runnable code per file (no pseudo-code)
- File path header before each code block
- Order: Domain → Application → Infrastructure → Interface
- Include: go.mod, .env, docker-compose.yml, migrations/*.sql

╔══════════════════════════════════════════════════════════╗
║  T–A–G LAYER: Task – Action – Goal                       ║
╚══════════════════════════════════════════════════════════╝

TASK:
Create the complete PDPA Module.

ACTION:
- Follow the 4-layer Clean Architecture strictly
- Use Value Objects, Domain Errors, Repository Pattern
- Apply Constructor-based Dependency Injection
- Use Context propagation in all layers
- Verify each layer builds: `go build ./...`

GOAL:
A module that:
- Compiles 100% on Go 1.21
- Runs via `docker-compose up`
- Passes `go vet` and `golangci-lint`
- Is integration-ready with existing user management system

╔══════════════════════════════════════════════════════════╗
║  B–A–B LAYER: Before – After – Bridge                    ║
╚══════════════════════════════════════════════════════════╝

BEFORE (สถานะปัจจุบัน):
- No consent management system
- No DSAR workflow (manual Excel)
- No audit trail
- Risk of PDPA non-compliance

AFTER (สถานะเป้าหมาย):
- Automated consent (grant/revoke/expire)
- DSAR workflow: submit → verify OTP → process → notify
- Blockchain-backed audit trail
- Auto-deletion policy (SUSPENDED + 1yr / TERMINATED + confirmed)
- Real-time WebSocket notifications

BRIDGE (ขั้นตอน):
1. Build PDPA as isolated bounded context
2. Connect to existing system via Kafka events
3. Migrate legacy consent data
4. Roll out with strangler fig pattern
5. Monitor: `consent_granted_total`, `dsar_processed_total`

╔══════════════════════════════════════════════════════════╗
║  C–A–R–E LAYER: Context – Action – Result – Example      ║
╚══════════════════════════════════════════════════════════╝

CONTEXT:
Thai company needs PDPA พ.ศ. 2562 + GDPR compliance.
Tech stack: Go 1.21, Gin, GORM, PostgreSQL 15, Redis 7, 
Kafka, Elasticsearch 8, Ethereum, SMTP, OpenAI.
Multi-tenant, high-volume, needs immutable proof.

ACTION:
Implement a module with:
- 7 tables (prefix `pdpa_`)
- 10 Use Cases
- 5 Kafka Workers + 1 Scheduler + 1 API Server
- WebSocket Hub
- i18n (EN/TH)
- Blockchain audit (Keccak256 → Smart Contract)

RESULT:
- Consent lifecycle: GRANTED → REVOKED/EXPIRED → DELETED
- DSAR responds within 3 days (vs 30-day legal limit)
- Every state change recorded in Blockchain
- Zero business logic in infrastructure layer

EXAMPLE (แนวทางอ้างอิง):
Model the design after GDPR-compliant systems like OneTrust 
and Transcend, but simplified for Thai context:
- Separate purpose per consent (ไม่ bundle)
- Retention policy per status
- Anonymization instead of hard delete for user PII

╔══════════════════════════════════════════════════════════╗
║  R–I–S–E LAYER: Role – Input – Steps – Expectation       ║
╚══════════════════════════════════════════════════════════╝

ROLE:
You are Principal Architect + Senior Go Developer.

INPUT (ข้อมูลที่ให้):

[Tech Stack]
- Go 1.21, Gin 1.9, GORM 2, PostgreSQL 15
- Redis 7, Kafka 3.x, Elasticsearch 8, Ethereum (Goerli/Mainnet)

[Database Schema - 7 Tables with prefix `pdpa_`]
- pdpa_policies, pdpa_purposes, pdpa_consents, 
  pdpa_user_requests, pdpa_user_account_statuses,
  pdpa_audit_trails, pdpa_request_responses

[10 Use Cases]
1. RecordConsent      6. GetDSARStatus
2. RevokeConsent      7. ImmediateDeletion
3. GetConsentHistory  8. AutoDeleteExpired
4. SubmitDSAR         9. HandleAccountEvent
5. ProcessDSAR       10. GetAdminReport

[Kafka Topics - 7]
pdpa.consent.log, pdpa.consent.revoked, pdpa.dsar.request,
pdpa.data.deleted, pdpa.account.event, pdpa.email.send, pdpa.llm.analyze

[Business Rules]
- Consent expires in 1 year
- Cannot revoke non-GRANTED consent
- SUSPENDED → retention 1 year
- TERMINATED + Confirmed → immediate delete
- DSAR OTP expires in 15 min

STEPS (ทำตามลำดับ):

Phase 1: Foundation
  1.1 Create go.mod, .env, docker-compose.yml
  1.2 Write migrations/001_initial_pdpa_schema.sql

Phase 2: Domain Layer
  2.1 Value Objects (5 files)
  2.2 Entities (4 files)
  2.3 Repository Interfaces (4 files)
  2.4 Domain Services (3 files)
  2.5 Domain Errors (1 file)

Phase 3: Application Layer
  3.1 DTOs
  3.2 Use Cases 1-10

Phase 4: Infrastructure Layer
  4.1 GORM Models (7)
  4.2 Repository Implementations (4)
  4.3 Redis Cache
  4.4 Kafka Producer + 5 Consumers
  4.5 Elasticsearch Indexer
  4.6 Services: Email, LLM, Blockchain

Phase 5: Interface Layer
  5.1 HTTP Handlers (4)
  5.2 Routes
  5.3 WebSocket Hub + Handler
  5.4 i18n (EN/TH)

Phase 6: Entry Points
  6.1 cmd/api/main.go
  6.2 cmd/scheduler/main.go
  6.3 cmd/workers/{dsar,email,llm,revoke,account}/main.go
  6.4 cmd/migrate/main.go

Phase 7: Testing
  7.1 Unit tests for Domain Services
  7.2 Unit tests for Use Cases (with mocks)
  7.3 Integration test for full flow

EXPECTATION (ความคาดหวัง):
- ✅ Code compiles with `go build ./...`
- ✅ `docker-compose up` brings up all services
- ✅ Test coverage > 70% for Use Cases
- ✅ No pseudo-code, no `// TODO` in critical paths
- ✅ Every layer respects dependency direction
- ✅ Domain layer has ZERO infrastructure imports
- ✅ Response includes self-check checklist at the end

═══════════════════════════════════════════════════════════
🧠 REASONING (Big Pickle Specific)
═══════════════════════════════════════════════════════════
Before writing code, use extended thinking:
1. List all assumptions
2. Explain schema FK relationships
3. Explain layer dependency direction
4. Identify trade-offs (Hard Delete vs Anonymize)
5. Plan implementation order

For each Use Case, explain:
- Why this operation sequence?
- What are failure modes?
- Partial failure handling (DB saved, Kafka failed)?

═══════════════════════════════════════════════════════════
🚀 BEGIN
═══════════════════════════════════════════════════════════
Start with Architecture Decision Summary → then Phase 1 → 7.
```

## 💡 เคล็ดลับเสริมสำหรับ Big Pickle

| สถานการณ์ | คำสั่งเพิ่ม |
|-----------|-------------|
| ต้องการ reasoning ลึก | ใส่ `Ultrathink` ท้าย prompt |
| Output ยาวเกิน 32K | แบ่ง 3 session: Domain / App+Infra / Interface |
| Context หลุด | Paste ไฟล์ที่มีอยู่แล้วลงใน prompt ตรงๆ |
| ลืม requirement | ใช้ checklist ท้าย prompt |

---

# 🔧 PART 2 — Prompt สำหรับ OpenCode

OpenCode เป็น **Agent Harness** ที่อ่าน repo จริง จึงใช้ **Custom Commands + AGENTS.md** แทนการ paste prompt ยาว

## 📁 โครงสร้างไฟล์ที่ต้องสร้าง

```
your-project/
├── AGENTS.md                          # R–I–S–E ฝังไว้ที่นี่
├── .opencode/
│   ├── config.json
│   └── commands/
│       ├── pdpa-domain.md             # R–T–F + T–A–G
│       ├── pdpa-app.md                # R–T–F + C–A–R–E
│       ├── pdpa-infra.md              # R–T–F + C–A–R–E
│       ├── pdpa-interface.md          # R–T–F + C–A–R–E
│       └── pdpa-deploy.md             # B–A–B
```

---

## 1️⃣ `AGENTS.md` (ใช้ R–I–S–E)

```markdown
# AGENTS.md — PDPA Module

## ROLE
You are working on a Go 1.21 PDPA Module for Thai PDPA + GDPR 
compliance. Follow Clean Architecture + DDD strictly.

## INPUT (Project Context)
- Module path: `internal/modules/pdpa/`
- 4 Layers: domain → application → infrastructure → interfaces
- Dependency rule: interfaces → application → domain ← infrastructure
- Database: PostgreSQL 15 (prefix `pdpa_` on ALL tables)
- Cache: Redis 7 (TTL 30 days)
- Queue: Kafka (7 topics — see below)
- Search: Elasticsearch 8
- Blockchain: Ethereum (Keccak256 proof)

## STEPS (Standard Workflow)
1. Read this file + relevant files in the repo
2. Plan changes in Plan mode (Shift+Tab)
3. Execute in Build mode
4. Run `go build ./...` after each change
5. Run `go vet ./...` and `golangci-lint run`

## EXPECTATION (Definition of Done)
- [ ] Code compiles: `go build ./...` exits 0
- [ ] Lint clean: `golangci-lint run` exits 0
- [ ] Tests pass: `go test ./... -v`
- [ ] No pseudo-code, no `// TODO`
- [ ] Domain layer has ZERO infrastructure imports
- [ ] Every entity has `NewXxx()` constructor

## Coding Standards
- Go 1.21 idiomatic
- `context.Context` in all public methods
- `fmt.Errorf("...: %w", err)` for wrapping
- Constructor injection: `NewXxxUseCase(...)`
- Domain errors: `domainerrors.ErrXxx`
- No global state, no singleton

## Database Rules
- ALL tables MUST be `pdpa_` prefixed
- Separate `entity` (domain) from `model` (infra)
- Explicit `CONSTRAINT` names for FKs
- Index on every FK column

## PDPA Business Rules
- Consent expires in 1 year from GRANTED
- Cannot revoke non-GRANTED consent
- SUSPENDED account → 1-year retention
- TERMINATED + DeletionConfirmed → immediate delete
- DSAR OTP expires in 15 min
- Every state change → audit trail

## Kafka Topics
- `pdpa.consent.log`      (producer: API)
- `pdpa.consent.revoked`  (producer: API → consumer: revoke worker)
- `pdpa.dsar.request`     (producer: API → consumer: dsar worker)
- `pdpa.data.deleted`     (producer: API/Scheduler)
- `pdpa.account.event`    (producer: external → consumer: account worker)
- `pdpa.email.send`       (producer: various → consumer: email worker)
- `pdpa.llm.analyze`      (producer: dsar worker → consumer: llm worker)

## Forbidden
- ❌ Business logic in infrastructure layer
- ❌ Importing infrastructure from domain
- ❌ Pseudo-code
- ❌ Skipping audit trail for state changes
```

---

## 2️⃣ `.opencode/config.json`

```json
{
  "instructions": [
    "AGENTS.md",
    ".opencode/commands/*.md"
  ],
  "model": "opencode/big-pickle",
  "permissions": { "edit": true, "bash": true }
}
```

---

## 3️⃣ Custom Commands (ใช้ R–T–F + T–A–G / C–A–R–E)

### `.opencode/commands/pdpa-domain.md`

```markdown
---
description: สร้าง PDPA Domain Layer
---

## ROLE (R)
Senior Go Developer specializing in DDD.

## TASK (T)
Create the complete Domain Layer for PDPA Module.

## FORMAT (F)
Full Go files, no pseudo-code. Run `go build` at the end.

## ACTION (A)
Create these files in `internal/modules/pdpa/domain/`:

**value_object/**
- consent_purpose.go    → ConsentPurpose (NECESSARY, ANALYTICS, MARKETING)
- consent_status.go     → ConsentStatus (GRANTED, REVOKED, EXPIRED, DELETED)
- dsar_type.go          → DSARType (ACCESS, ERASURE, WITHDRAW_CONSENT)
- dsar_status.go        → DSARStatus (PENDING, PROCESSING, COMPLETED, REJECTED)
- account_status.go     → AccountStatus (ACTIVE, SUSPENDED, TERMINATED)

**entity/**
- consent_log.go          → +Revoke() +MarkDeleted() +IsActive()
- dsar_request.go         → +VerifyOTP() +MarkProcessing() +MarkCompleted() +MarkRejected()
- audit_trail.go          → +NewAuditTrail() +NewAuditTrailWithMeta()
- user_account_status.go  → +SetSuspended() +SetTerminated() +ConfirmDeletion() +MarkDeleted()

**repository/** (interfaces only)
- consent_repository.go
- dsar_repository.go
- audit_repository.go
- user_account_status_repository.go

**service/**
- deletion_policy_service.go  → CanImmediateDeletion, IsReadyForAutoDeletion, AnonymizeUserData
- consent_validator.go        → ValidatePurpose
- blockchain_service.go       → interface only

**errors/**
- errors.go  → 13+ domain errors

**event/**
- event.go   → WebSocketEvent

## GOAL (G)
`go build ./internal/modules/pdpa/domain/...` exits 0.
NO imports from application/infrastructure.

## EXPECTATION
Read AGENTS.md. Follow all coding standards. Run build.
Report: files created + build result + self-check.
```

### `.opencode/commands/pdpa-app.md`

```markdown
---
description: สร้าง PDPA Application Layer (Use Cases)
---

## ROLE (R)
Senior Go Developer, DDD practitioner.

## TASK (T)
Create 10 Use Cases + DTOs.

## FORMAT (F)
Full Go files. Build must pass.

## CONTEXT (C)
Domain Layer already exists. Infra does NOT exist yet.
Use only interfaces from `domain/repository/`, `domain/service/`, 
and interfaces for cache/producer/indexer (define them here if needed).

## ACTION (A)
Create in `internal/modules/pdpa/application/`:

1. record_consent.go               — RecordConsentUseCase
2. revoke_consent.go               — RevokeConsentUseCase
3. get_consent_history.go          — GetConsentHistoryUseCase
4. submit_dsar.go                  — SubmitDSARUseCase (genOTP helper)
5. process_dsar.go                 — ProcessDSARUseCase (+ WebSocket broadcast)
6. get_dsar_status.go              — GetDSARStatusUseCase
7. immediate_deletion.go           — ImmediateDeletionUseCase (+ Blockchain proof)
8. auto_delete_expired_consents.go — AutoDeleteExpiredConsentsUseCase
9. handle_account_event.go         — HandleAccountEventUseCase
10. get_admin_report.go            — GetAdminReportUseCase
11. dto.go                         — all Input/Output DTOs

## RESULT (R)
All 10 Use Cases with:
- Struct + `NewXxx()` constructor
- `Execute(ctx, input)` method
- Domain errors for business violations
- Audit trail on every state change
- Kafka event publish after DB save

## EXAMPLE (E)
```
type RecordConsentUseCase struct {
    consentRepo       repository.ConsentRepository
    accountStatusRepo repository.UserAccountStatusRepository
    auditRepo         repository.AuditRepository
    cache             redis.ConsentCache
    producer          messaging.KafkaProducer
    blockchainSvc     service.BlockchainService
}

func (uc *RecordConsentUseCase) Execute(ctx context.Context, input RecordConsentInput) error {
    // 1. validate account status
    // 2. save consent
    // 3. cache
    // 4. kafka
    // 5. blockchain (async)
    // 6. audit
}
```

## EXPECTATION
Read AGENTS.md. `go build ./internal/modules/pdpa/application/...` exits 0.
```

### `.opencode/commands/pdpa-infra.md`

```markdown
---
description: สร้าง PDPA Infrastructure Layer
---

## ROLE (R)
Senior Go Developer, expert in GORM/Kafka/Redis/ES/Ethereum.

## TASK (T)
Implement all infrastructure adapters for PDPA Domain interfaces.

## FORMAT (F)
Full Go files. Build passes.

## CONTEXT (C)
Domain + Application already exist.
You must implement repository interfaces defined in domain.

## ACTION (A)
Create in `internal/modules/pdpa/infrastructure/`:

**persistence/postgres/**
- models.go                     → 7 GORM models, all with TableName() returning pdpa_*
- consent_repo_impl.go
- dsar_repo_impl.go
- audit_repo_impl.go
- user_account_status_repo_impl.go

**persistence/redis/**
- consent_cache.go              → ConsentCache interface + impl

**messaging/**
- kafka_producer.go             → KafkaProducer interface + impl
- consumers/dsar_worker.go
- consumers/email_worker.go
- consumers/llm_worker.go
- consumers/revoke_worker.go
- consumers/account_event_worker.go

**search/elasticsearch/**
- consent_indexer.go

**services/email/smtp.go**
**services/llm/openai.go**
**services/blockchain/ethereum.go**

**scheduler/consent_cleanup_job.go**

## RESULT (R)
All domain interfaces implemented. Repos use GORM.
Kafka uses sarama. Blockchain uses go-ethereum.

## EXAMPLE (E)
Consent repo mapping: `PdpAConsentModel` (infra) ↔ 
`entity.ConsentLog` (domain). No domain leaks in infra layer.

## EXPECTATION
`go build ./internal/modules/pdpa/infrastructure/...` exits 0.
Every domain interface has exactly one implementation.
```

### `.opencode/commands/pdpa-interface.md`

```markdown
---
description: สร้าง PDPA Interface Layer + Entry Points
---

## ROLE (R)
Senior Go Developer, expert in Gin + WebSocket.

## TASK (T)
Create HTTP handlers, routes, WebSocket hub, i18n, and all main.go.

## FORMAT (F)
Full Go files. Build passes.

## CONTEXT (C)
All lower layers exist. Ready to wire up.

## ACTION (A)
Create:

**interfaces/http/**
- consent_handler.go
- dsar_handler.go
- deletion_handler.go
- admin_handler.go
- websocket_handler.go
- routes.go
- dto.go

**interfaces/websocket/**
- hub.go                     → Hub + Client + WritePump/ReadPump

**interfaces/localization/**
- i18n.go                    → I18nManager with embed
- locales/en.toml
- locales/th.toml

**cmd/**
- api/main.go                → wire everything, graceful shutdown
- scheduler/main.go          → cron 0 2 * * *
- migrate/main.go            → run migrations
- workers/dsar/main.go
- workers/email/main.go
- workers/llm/main.go
- workers/revoke/main.go
- workers/account/main.go

## RESULT (R)
All binaries build. API exposes:
- POST   /api/v1/pdpa/consent
- DELETE /api/v1/pdpa/consent
- POST   /api/v1/pdpa/dsar
- POST   /api/v1/pdpa/deletion/confirm
- POST   /api/v1/pdpa/events/account
- GET    /api/v1/pdpa/admin/reports
- GET    /api/v1/pdpa/admin/audit
- GET    /ws (WebSocket with JWT auth)

## EXPECTATION
`go build ./...` exits 0 for entire project.
Graceful shutdown handles SIGINT/SIGTERM.
```

### `.opencode/commands/pdpa-deploy.md` (ใช้ B–A–B)

```markdown
---
description: สร้าง deployment artifacts + migration plan
---

## BEFORE (B)
Project compiles locally. No docker-compose. No .env.

## AFTER (A)
Full local-dev environment:
- docker-compose.yml with postgres, redis, kafka, zookeeper, elasticsearch
- .env with all vars (DB_DSN, REDIS_ADDR, KAFKA_BROKERS, etc.)
- migrations/001_initial_pdpa_schema.sql (idempotent)
- README with setup steps
- One command: `docker-compose up -d && go run cmd/migrate/main.go && go run cmd/api/main.go`

## BRIDGE (B)
Steps to produce:
1. Create `docker-compose.yml` (services: postgres:15, redis:7, 
   confluentinc/cp-kafka, elasticsearch:8.11.0)
2. Create `.env` with placeholders
3. Write SQL migration with `IF NOT EXISTS` + seed data
4. Create `cmd/migrate/main.go` that reads SQL and applies
5. Write `README.md` with:
   - Prerequisites
   - Setup commands
   - Endpoint list
   - Architecture diagram (Mermaid)

## EXPECTATION
User can run `docker-compose up -d` → `go run cmd/migrate/main.go` 
→ `go run cmd/api/main.go` and hit endpoints.
```

---

## 4️⃣ วิธีใช้ OpenCode กับ Prompt นี้

```bash
# 1. Setup ครั้งแรก
cd your-project
mkdir -p .opencode/commands
# สร้าง AGENTS.md + config.json + command files ตามด้านบน

# 2. เปิด OpenCode
opencode

# 3. ใน OpenCode TUI:
# กด Shift+Tab → Plan mode
/plan อ่าน AGENTS.md แล้วสรุปแผนการสร้าง PDPA Module ตาม commands ทั้ง 5

# ตรวจสอบแผน → กด Shift+Tab กลับ Build mode

# 4. Execute ทีละ layer
/pdpa-domain      # สร้าง Domain Layer
/pdpa-app         # สร้าง Application Layer
/pdpa-infra       # สร้าง Infrastructure Layer
/pdpa-interface   # สร้าง Interface + Entry points
/pdpa-deploy      # สร้าง docker-compose + .env + migration

# 5. Verify
"go build ./... && go test ./... && golangci-lint run"
```

---

# 📊 Mapping Table — Framework × Layer × Tool

| Framework | ใช้ที่ไหน | Big Pickle | OpenCode |
|-----------|----------|-----------|----------|
| **R–T–F** | กำหนด role/task/format | ✅ Prompt ต้น | ✅ ทุก command file |
| **T–A–G** | กำหนด task/action/goal | ✅ ส่วนที่ 2 | ✅ AGENTS.md + Goal section |
| **B–A–B** | เปลี่ยนสถานะโปรเจกต์ | ✅ ส่วนที่ 3 | ✅ `pdpa-deploy.md` |
| **C–A–R–E** | ให้บริบท + ตัวอย่าง | ✅ ส่วนที่ 4 | ✅ `pdpa-app/infra.md` |
| **R–I–S–E** | ให้ข้อมูล + steps + expectation | ✅ ส่วนที่ 5 | ✅ `AGENTS.md` ทั้งไฟล์ |

---

# 🎯 จุดต่างสำคัญระหว่าง 2 Tools

| ประเด็น | Big Pickle | OpenCode |
|---------|-----------|----------|
| **Framework หลัก** | R–I–S–E (เพราะต้องให้ input+steps ครบ) | R–T–F (เพราะเป็น command สั้นๆ) |
| **Framework รอง** | R–T–F + T–A–G | C–A–R–E + B–A–B |
| **Prompt ยาวแค่ไหน** | ยาว (paste ทั้งบล็อก) | สั้น (แยกเป็น command files) |
| **การจัดการ context** | Manual — paste เอง | Automatic — อ่าน repo |
| **Output** | Code blocks ใน chat | แก้ไฟล์จริง |
| **เหมาะกับ** | งานเดียวจบ / prototype | งานหลาย layer / production |

---

# ⚠️ ข้อควรระวัง

1. **Big Pickle context limit 200K** — ถ้า paste โค้ดเดิมเยอะ อาจชนขอบเขต → แบ่ง session
2. **Big Pickle output limit 32K tokens** — แต่ละ phase ควรจบใน 1 session
3. **OpenCode ต้องมี repo ก่อน** — ไม่ใช่ chat wrapper
4. **Plan mode สำคัญมาก** — อย่าข้ามไป Build ทันที
5. **ตรวจ `go build` ทุกครั้งหลังรัน command** — อย่าปล่อยให้ error สะสม

---
