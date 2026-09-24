# PDPA Module Implementation Plan

Status: Draft — based on gap analysis of `docs/Modules_PDPA_0.md` (canonical) + `docs/Modules_PDPA_CMD.md` vs. current `internal/modules/pdpa` code. Out of scope: `docs/Modules_PDPA.md` (3164-line duplicate of `_0`, skipped).

## 1. Architecture Diagram

Clean Architecture + DDD + Modular Monolith. Strict dependency rule (CMD 1-6): `interfaces` → `application` → `domain`; `infrastructure` imports only `domain` interfaces. No `gorm`, `gin`, `chi`, `sarama`, `redis`, `elasticsearch` in `domain` (CMD 7).

```
                          ┌───────────────────────────────────────────────┐
                          │                  HTTP (interfaces)             │
                          │   handlers/* (10)  http/routes.go             │
                          └─────────────────────┬─────────────────────────┘
                                                │ Request → Command/Query
                          ┌─────────────────────▼─────────────────────────┐
                          │                application                     │
                          │  command (10)   query (5)   dto (13)           │
                          │  event_handler (NEW)   port/  (NEW)           │
                          └──────────────┬──────────────────────┬─────────┘
                                         │                      │
   ┌─────────────────────────────────────▼──┐   ┌───────────────▼──────────────┐
   │                 domain                 │   │           infrastructure      │
   │  entity / value_object / repository /  │   │  persistence/postgres (impls) │
   │  event / service / errors              │   │  outbox_repo_impl.go (NEW)    │
   │                                        │   │  messaging/kafka producer.go  │
   │  Command → DomainService               │   │  cache/redis consent_cache.go│
   └────────────────────────────────────────┘   │  scheduler/ + es + search    │
                                                │  localization                │
                                                └──────────────────────────────┘

   ═══ Transactional outbox (NEW, from scratch — repoint none; no precedent) ═══

    Command (application/command)
        │ 1. DomainService mutates aggregate (injectable clock)
        ▼
    same DB tx ──► [consent/subject/dsar/policy tables]  (writer)
           └──► [outbox table] append {topic, key, payload, status=pending}
        │
    OutboxDispatcher (scheduler/outbox_dispatcher_job.go)  -- polls pending rows
        │ claim batch (status=dispatching, atomic UPDATE ... WHERE status=pending)
        ▼
    Kafka producer (messaging/kafka/producer.go) PublishMessage("pdpa.{entity}.{action}")
        │ on ack → mark dispatched; on fail → retry w/ backoff, max-deliveries → DLQ topic
        ▼
    Event handlers (application/event_handler/*)  -- idempotent, correlationID, PII-redacted
        ├─ consent.recorded / dsar.status.changed / account.status.updated
        ├─ data.deleted              (NEW)
        ├─ llm.analysis.request      (NEW)
        ├─ email.notification        (NEW)
        └─ blockchain.record         (NEW)

   Atomicity contract (CMD 4, 12): commands no longer call producer.PublishMessage
   directly; they persist domain state + matched outbox row in one tx. The poller
   publishes and acks, so every published event corresponds to a committed write.
```

## 2. File Tree

```
internal/modules/pdpa
├── module.go
├── domain
│   ├── entity/                          # Consent, ConsentLog, DataSubject, DataSubjectLog, Policy, DSAR
│   ├── value_object/
│   │   ├── consent_status.go
│   │   ├── dsar_status.go               # P1: string-based (spec)
│   │   ├── account_status.go            # P1: string-based (spec)
│   │   ├── dsar_type.go                 # P1: string-based (spec)
│   │   └── purpose.go
│   ├── repository/                      # interfaces (domain owns them)
│   │   ├── consent_repo.go
│   │   ├── consent_log_repo.go
│   │   ├── subject_repo.go
│   │   ├── dsar_repo.go
│   │   └── policy_repo.go               # P1: NEW (spec gap)
│   ├── event/
│   │   ├── consent_events.go
│   │   ├── dsar_events.go
│   │   ├── event_metadata.go            # P1: NEW — correlationID, occurredAt, entity, action
│   │   └── outbox_event.go              # P0: NEW — payload envelope for outbox rows
│   ├── service/
│   │   ├── consent_service.go
│   │   └── dsar_service.go
│   └── errors/                          # sentinel errors (CMD 10)
├── application
│   ├── command/                         # 10 existing (all publish via producer today — P0 fix)
│   │   ├── create_consent_command.go
│   │   ├── create_dsar_command.go
│   │   ├── verify_otp_command.go
│   │   ├── complete_dsar_command.go
│   │   ├── cancel_dsar_command.go
│   │   ├── create_data_subject_command.go
│   │   ├── update_data_subject_command.go
│   │   ├── change_consent_status_command.go
│   │   ├── delete_consent_command.go
│   │   └── record_consent_command.go    # P1: NEW (spec gap)
│   ├── query/                           # 5 read-only (spec mandates 4 — P1: trim/align)
│   ├── dto/                             # 13 DTOs (keep)
│   ├── event_handler/                   # P1: NEW dir
│   │   ├── consent_recorded_handler.go
│   │   ├── dsar_status_changed_handler.go
│   │   └── ... (idempotent + DLQ hooks)
│   └── port/                            # P0: NEW — boundaries owned by application
│       ├── outbox_repo.go               # Append(ctx,e) / Claim(ctx,limit) / Complete(ctx,id) / Fail(ctx,id)
│       └── clock.go                     # NewClock() func() time.Time (P0)
├── infrastructure
│   ├── persistence/postgres/
│   │   ├── consent_repo_impl.go
│   │   ├── consent_log_repo_impl.go
│   │   ├── subject_repo_impl.go
│   │   ├── dsar_repo_impl.go
│   │   ├── policy_repo_impl.go          # P1: NEW
│   │   └── outbox_repo_impl.go          # P0: NEW (snake_case per CMD FORMAT)
│   ├── persistence/models/              # 7 GORM models, pdpa_* prefix (CMD 11)
│   ├── messaging/kafka/
│   │   ├── producer.go                  # P0: wrap PublishMessage (exists in messaging/?)
│   │   └── consumer.go                  # subscription for event handlers
│   ├── cache/redis/
│   │   └── consent_cache.go             # P0: NEW — keys pdpa:{entity}:{id} (CMD 9)
│   ├── scheduler/
│   │   ├── outbox_dispatcher_job.go     # P0: NEW poller
│   │   └── consent_cleanup_job.go       # P0: NEW async worker (spec)
│   ├── es/                              # (dir only — no .go yet)
│   ├── search/                          # (dir only — no .go yet)
│   └── localization/locales/
├── repository/                          # currently EMPTY — resolve mocks placement (open item)
└── interfaces
    ├── handlers/                        # 10 handlers + handler_helpers.go
    └── http/                            # routes.go
```

## 3. Gaps Summary (verified)

| ID | Sev | Gap |
|----|-----|-----|
| G1 | P0 | No transactional outbox anywhere in repo (`outbox`/`Outbox` grep = 0 hits). Commands publish via `producer.PublishMessage` directly, not atomic with DB persist (CMD 4, 12) |
| G2 | P0 | OTP stored raw (no sha256), TTL 5m not spec 15m (`OTPTTL`), attempt limit only in-memory `sync.Map` (not `OTPMaxAttempt=5`, not shared per-subject) |
| G3 | P0 | `NewConsentLog` 3-arg; spec requires injectable `now` (5-arg: userID, purpose, now, ip, userAgent) |
| G4 | P0 | DSAR lifecycle incomplete — no `MarkProcessing` / `MarkCompleted` / `OTPVerified` state persistence in commands |
| G5 | P1 | Missing `RecordConsentCommand` |
| G6 | P1 | `policy_repo.go` interface missing in domain (policy exists in models/migrations) |
| G7 | P1 | String-based VOs for `DSARStatus`, `AccountStatus`, `DSARType` not implemented |
| G8 | P1 | `event/event_metadata.go` missing (correlationID/tracing, CMD 14) |
| G9 | P1 | `application/event_handler/` directory missing |
| G10 | P1 | Query set = 5; spec = 4 (trim/align + keep read-only) |
| G11 | P2 | Handler hardening: idempotency, DLQ, PII redaction, `pdpa.` topic prefix (spec 13-topic list) |
| G12 | P2 | Missing handlers: `data.deleted`, `llm.analysis.request`, `email.notification`, `blockchain.record` |
| G13 | P2 | ConsentLog: behaviors `Revoke(now)`, `MarkDeleted(now)`, `IsActive(now)` not wired |
| G14 | P2 | Terraform: pdpa-msk (kafka 3.6.0, replication 3, min.insync 2, partitions=6, retention 168h/100GiB) + pdpa-secrets (random_password anon_salt len 48) |
| G15 | P2 | Infrastructure dirs `es`, `search`, `cache/redis`, `scheduler`, `messaging/kafka` empty — build consent_cache, cleanup job, es/search adapters |

## 4. Implementation Phases

### Phase 0 — Outbox foundation (G1)
1. `domain/repository/outbox_repo.go` (interface: Append, Claim, Complete, Fail) + `domain/event/outbox_event.go` (envelope: ID, Topic, Key, Payload, Status, Attempts, CreatedAt, DispatchedAt).
2. `application/port/clock.go` + `outbox_repo.go` (application-side boundary).
3. `infrastructure/persistence/postgres/outbox_repo_impl.go` + migration `NNN_create_pdpa_outbox.sql` (`pdpa_outbox_events`, idx on `(status, created_at)`).
4. `infrastructure/scheduler/outbox_dispatcher_job.go`: poll → claim (atomic `UPDATE ... WHERE status='pending' RETURNING`) → publish → `Complete`; fail → backoff, mark `Failed` after max-deliveries, DLQ topic `pdpa.dlq`.
5. Refactor all 10 commands: remove direct `producer.PublishMessage`; persist + `outboxRepo.Append` in same tx.

### Phase 1 — Domain/spec conformance (G2-G10)
6. OTP sha256 + `OTPTTL=15m`, `OTPMaxAttempt=5`; replace in-memory `sync.Map` with Redis/DB-backed attempts, TTL expiry + mark `OTPVerified`.
7. `NewConsentLog(userID, purpose, now, ip, userAgent)` — inject clock everywhere; wire `verify_otp_command` / `complete_dsar_command` / `cancel_dsar_command` to persist `Processing/Completed` states.
8. Add `RecordConsentCommand`, wire to `ConsentLog.Create`.
9. Add domain `policy_repo.go` + `infrastructure/persistence/postgres/policy_repo_impl.go`.
10. Convert `DSARStatus`, `AccountStatus`, `DSARType` VOs to string-based (align names to spec) with validation.
11. Add `event/event_metadata.go` (correlationID, occurredAt, source).
12. Create `application/event_handler/`; register handlers for consent recorded / dsar status changed; align query set to spec (4).

### Phase 2 — Robustness / hardening (G11-G15)
13. Idempotency keys + DLQ wiring in handlers; PII redaction on payloads; `pdpa.` topic prefix; add missing handlers (`data.deleted`, `llm.analysis.request`, `email.notification`, `blockchain.record`).
14. Wrap `ConsentLog.Revoke(now)/MarkDeleted(now)/IsActive(now)`.
15. Terraform `pdpa-msk` + `pdpa-secrets`.
16. Build `consent_cache.go`, `consent_cleanup_job.go`, es/search adapters.

## 5. Outbox Design Notes (from-scratch, no precedent)

- Same tx Append: handler implements `domain.RepositoryTx`; both entity persist and outbox Append share `*gorm.DB`/session (check module's existing tx mechanism in infrastructure).
- Exactly-once-ish: outbox rows give at-least-once; idempotency keys + dedupe index `(event_id, handler)` make it effectively-once.
- Dispatcher: single consumer group `pdpa-outbox`; claim batch (e.g. 100); commit after publish-ack.
- DLQ: `pdpa.dlq.{entity}`; poison rows marked `Failed` with last error.

## 6. Missing Decision (blocker before Phase 0 coding)

- **Mocks placement unresolved**: `application/mocks` not found; `repository/` dir is empty; no `_test.go` globs confirmed. Must locate existing test-double convention in a sibling module before wiring `New{Entity}` constructor tests (CMD 13: coverage ≥ 80%).

## 7. Acceptance Criteria

- Every command path: entity committed AND outbox row appended in one tx; no direct publisher call in `application/command` (CMD 4, 12).
- `pdpa.{entity}.{action}` topics confirmed on publish; dispatcher + DLQ operational.
- OTP: sha256 only; 15m TTL; 5 attempts enforced per subject across restart.
- `Revoke/MarkDeleted/IsActive` + injectable clock; CMD 15 (bilingual comments) on all new code.
- `go test ./internal/modules/pdpa/...` ≥ 80% coverage (CMD 13); `go vet` + `golangci-lint` clean.
- Terraform plan applies; secrets include `anon_salt` (48) + MSK creds.

## 8. Open Items

- Confirm exact spec 4-query set before deleting/renaming query files (G10).
- Confirm `record_consent` vs `create_consent` intent (G5) with the CMD spec list of commands.
- Confirm Kafka topic names from spec (13-topic list) before prefixing.me