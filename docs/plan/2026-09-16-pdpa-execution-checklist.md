# PDPA Module — Execution Checklist (PDPA-1..24)

> Mapping: plan path/name -> actual repo file + สถานะ + ส่วนเบี่ยงเบน (deviation)
> อ้างอิง: `Modules_plan_PDPA_Task.md` (L886-913 sub-task rows), `2026-09-15-pdpa-module.md` (master L13-37), repo inventory 2026-09-16.
> สถานะ: DONE = มีไฟล์จริงแล้ว / ADD = ต้องสร้าง / UPDATE = ไฟล์มีแล้วต้องแก้ / VERIFY = ต้องตรวจก่อน

## โครงสร้างจริงของ `internal/modules/pdpa` (verified 2026-09-16)

```
internal/modules/pdpa/
├── modulesdoc.md, modulesdoc_pdpa.md
├── application/
│   ├── command/       10 files (PascalCase)
│   ├── dto/           13 files (lowercase)
│   ├── port/          clock.go, outbox_repo.go
│   └── query/         4 files (PascalCase)
├── domain/
│   ├── entity/        7 files
│   ├── errors/        errors.go
│   ├── event/         event.go, outbox_event.go
│   ├── repository/    6 files  (policy_repo.go MISSING)
│   ├── service/       blockchain_service.go, deletion_policy_service.go
│   └── valueobject/   5 files
├── infrastructure/
│   ├── persistence/pg_repository.go        (untracked — VERIFY)
│   ├── persistence/postgres/outbox_repo_impl.go
│   └── scheduler/outbox_dispatcher_job.go (+ _test.go)
└── interfaces/
    ├── handlers/      10 files + handler_helpers.go  (event handlers อยู่ที่นี่)
    └── http/routes.go
```

## Checklist รายข้อ

| # | งาน (plan) | path ตาม plan | สถานะจริง | หมายเหตุ / เบี่ยงเบน |
|---|---|---|---|---|
| 1 | Database Schema+Migration+Seed | `migrations/001_initial_pdpa_schema.sql`, `002_seed_initial_policy.sql` | DONE | ใช้ชื่อ `20260902..20260910_pdpa_*.sql` (ไม่ใช่ `00x_`) — คงชื่อเดิม 7 ไฟล์ M + 1 ใหม่ (20260910_processed_events) |
| 2 | VO + Errors | `domain/valueobject/*`, `domain/errors/errors.go` | DONE | ครบ 5 VO + errors.go |
| 3 | Events + Metadata | `domain/event/event.go`, `outbox_event.go` | DONE | 16 Kafka topic constants (NAME ยืนยัน; VALUES ยังไม่ verify) |
| 4 | Entities | `domain/entity/*` | DONE | ครบ 7 |
| 5 | Domain Services | `domain/service/*` | DONE | blockchain_service.go + deletion_policy_service.go |
| 6 | Repository Interfaces + Cache | `domain/repository/*` | **ADD** | มี 6; **ไม่มี policy_repo.go** → ADD `FindByVersion/DeactivateAll/Save` |
| 7 | Infrastructure (postgres models + UnitOfWorkImpl) | `infrastructure/persistence/postgres/models.go`, `uow.go` | **ADD (VERIFY)** | มีเฉพาะ outbox_repo_impl.go; ต้องอ่าน `pg_repository.go` ก่อน |
| 8 | DSAR Repo Impl | `infrastructure/persistence/postgres/dsar_repo_impl.go` | **ADD (VERIFY)** | ยังไม่มี |
| 9 | Account/Audit/Policy/Outbox Repo Impl | `.../postgres/*_repo_impl.go` | **ADD (VERIFY)** | ยังมีแค่ outbox_repo_impl.go |
| 10 | Redis Cache + Idempotency | `infrastructure/redis/*` | **ADD** | ไม่มี dir `redis` |
| 11 | Kafka Producer + Outbox Publisher | `infrastructure/kafka/producer.go`, `outbox_publisher.go` | **ADD** | มีแล้วแค่ scheduler/outbox_dispatcher_job.go (verify group/batch/DLQ/cleanup/dedupe) |
| 12 | External Clients (Email/LLM/Blockchain) | `infrastructure/external/*` | **ADD (VERIFY)** | ไม่มี dir `external` |
| 13 | Application DTOs | `application/dto/dto.go` + named | DONE+ADD | 13 ไฟล์ lowercase, ไม่มี dto.go; เพิ่ม 5 response DTO (PDPA-19) |
| 14 | Consent Commands | `application/command/Grant*`, `Revoke*` | DONE | GrantConsentCommand.go, RevokeConsentCommand.go |
| 15 | DSAR Commands | `application/command/*DSAR*` | DONE | Submit/Process/VerifyOTP/Complete/Reject |
| 16 | Account Commands | `application/command/Suspend*`, `Terminate*`, `ConfirmDeletion*` + immediate deletion + auto delete expired | PARTIAL (**ADD**) | มี 3; **ADD immediate deletion + autoDeleteExpired + `AutoDeleteExpiredResult`** |
| 17 | Policy Command | `application/command/PublishPrivacyPolicyCommand.go` | **ADD** | block: ต้องมี PolicyRepository (6) + UnitOfWork (7) ก่อน |
| 18 | Queries | `application/query/*` | DONE+ADD | มี 4; **ADD GetConsentHistory / GetActivePrivacyPolicy / ListPolicies** ตามต้องการ PDPA-19 |
| 19 | Event Handler Base + PII Redactor | `application/event_handler/base_handler.go`, `pii_redactor.go` | **ADD** | dir ยังไม่มี; handler จริงอยู่ `interfaces/handlers/` — ต้องตัดสินใจ reuse vs สร้างใหม่ |
| 20 | Event Handlers (Consent & DSAR) | `application/event_handler/*` | VERIFY | มีอยู่แล้ว: consent_granted/revoked, dsar_submitted/completed — อยู่ interfaces/handlers/ |
| 21 | Event Handlers (Account/Data/LLM/Email/Blockchain) | `application/event_handler/*` | VERIFY | มีอยู่แล้ว: account_suspended/terminated, data_deletion_requested, audit_trail, dlq |
| 22 | HTTP Handlers | `interfaces/http/handlers/*` (consent/dsar/account/policy/admin) | **ADD** | ยังไม่มี; ต้อง ADD 5 handler + `extractUserID` |
| 23 | Routes | `interfaces/http/routes.go` | **UPDATE** | routes.go มีแล้ว; ต้อง ADD `Handlers` struct + `RegisterRoutes` + `Middleware` |
| 24 | Module Entry Point | `module.go` | **ADD** | `NewModule/RegisterHTTP/StartBackground/Shutdown` + SWAGGO annotations ทุก endpoint |

## Decisions/Deviations ที่ต้องบันทึกส่งท้าย

1. **Event handlers อยู่ `interfaces/handlers/` ไม่ใช่ `application/event_handler/`** (plan path ไม่ตรง) → ตัดสินใจ: reuse 10 ไฟล์เดิม + ADD base_handler/pii_redactor ไว้ `interfaces/handlers/` (กัน split-brain) — ยืนยันก่อนลงมือ
2. **Infra : เห็นแค่ outbox_repo_impl.go + pg_repository.go (untracked)** — ข้ออ้างเดิม "PDPA-7..13 เสร็จแล้ว" ยังไม่ยืนยัน → อ่าน `pg_repository.go` ก่อน ปรับ todo ถ้าต้อง implement ใหม่จากศูนย์
3. **policy_repo.go หาย** → PDPA-17/18 block; ต้องหา interface ที่เกี่ยวข้อง (อาจซ่อนใน pg_repository.go)
4. **Migration ชื่อ**: รักษา `2026090x_` เดิม แทน `001_/002_` ของ plan
5. **PascalCase command/query + lowercase DTO** (ตรงกันทั้ง repo; ไม่ใช่ตามตัวอักษร plan)
6. **`GetConsentStatusQuery.go` ถูกลบ (git D)** — ต้องเช็คว่า PDPA-19 handler ใช้ query ตัวไหนแทน แล้วบันทึก
7. **Mocks**: master เปิดคำถามข้อ 4 → commit ใช้ `domain/port/mocks/`
8. **Kafka topic VALUES** ยังไม่ verify — เช็ค ก่อน module.go Config รอ

## ลำดับลงมือ (align todo list)

1. เขียน checklist นี้ (จบ) → 2. อ่าน pg_repository.go → 3. fix outbox TableName() → 4. PDPA-6/7/8/9 (PolicyRepository + infra) → 5. PDPA-10/11/12 (redis/kafka/external) → 6. PDPA-16/17/18 commands+query → 7. PDPA-19/20/21 (base/pii+handlers) → 8. PDPA-22/23/24 HTTP+routes+module+swagger → 9. mockery mocks → 10. build/vet/lint/test → 11. รายงานไทย caveman