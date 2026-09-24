# การแตก Sub-tasks สำหรับ PDPA Module

## ✅ ขั้นตอนที่ 0: ตรวจสอบ Input

**รหัสตั๋วหลัก:** PDPA  
**รายละเอียดงาน:** สร้าง Module PDPA ตาม Clean Architecture + DDD + EDA (จากเอกสารที่ให้มา)

---

## 🔍 ขั้นตอนที่ 1: วิเคราะห์ว่าต้องสร้าง Database task หรือไม่

**สร้าง Database task** เนื่องจาก:
- ✅ สร้างตารางใหม่ (DDL) — `pdpa_consents`, `pdpa_dsar_requests`, `pdpa_user_account_statuses`, `pdpa_audit_trails`, `pdpa_privacy_policies`, `pdpa_outbox`
- ✅ เพิ่ม Migration script
- ✅ INSERT ข้อมูล seed (purposes, initial policy)

---

## 📋 Sub-tasks ที่สร้าง

---

### PDPA-1: [database] สร้าง Schema และ Migration สำหรับ PDPA Module

```
📁 file: migrations/001_initial_pdpa_schema.sql, migrations/002_seed_initial_policy.sql
📥 Input: -
📤 Output: -
🔁 Flow: -
🧪 Unit: -
📘🔗 Swagger: -
🎯 Objective:
- สร้างตาราง pdpa_consents สำหรับบันทึกประวัติความยินยอม
  - คอลัมน์: id, user_id, session_id, purpose_code, status, ip_address, user_agent, granted_at, expires_at, revoked_at, auto_deleted_at, created_at, updated_at
  - Index: idx_pdpa_consents_user_id, idx_pdpa_consents_user_purpose, idx_pdpa_consents_status, idx_pdpa_consents_active
  - Constraint: chk_pdpa_consents_purpose, chk_pdpa_consents_status
- สร้างตาราง pdpa_dsar_requests สำหรับคำร้องขอใช้สิทธิ์
  - คอลัมน์: id, user_id, request_type, status, requested_at, completed_at, data_payload, data_hash, blockchain_tx, rejection_reason, otp_hash, otp_expires_at, otp_verified, otp_attempts, ip_address, user_agent
  - Index: idx_pdpa_dsar_user_id, idx_pdpa_dsar_status, idx_pdpa_dsar_rate_limit
- สร้างตาราง pdpa_user_account_statuses สำหรับสถานะบัญชี
  - คอลัมน์: id, user_id, status, suspended_at, terminated_at, deletion_confirmed_at, retention_deadline, auto_deleted_at
  - Unique Index: uq_pdpa_account_user_id
- สร้างตาราง pdpa_audit_trails สำหรับ audit logs
  - คอลัมน์: id, user_id, action, details (JSONB), ip_address, user_agent, created_at
  - Index: idx_pdpa_audit_user_id, idx_pdpa_audit_action, idx_pdpa_audit_created_at, idx_pdpa_audit_details (GIN)
- สร้างตาราง pdpa_privacy_policies สำหรับนโยบายความเป็นส่วนตัว
  - คอลัมน์: id, version, title, content, effective_date, is_active
  - Unique Index: uq_pdpa_policy_version, uq_pdpa_policy_active
- สร้างตาราง pdpa_outbox สำหรับ Transactional Outbox Pattern
  - คอลัมน์: id, aggregate_type, aggregate_id, event_type, topic, payload (JSONB), status, retry_count, last_error, created_at, published_at
  - Index: idx_pdpa_outbox_status_created (WHERE status='PENDING')
- สร้างตาราง pdpa_processed_events สำหรับ Idempotency
  - คอลัมน์: event_id, handler_name, processed_at, expires_at
  - Primary Key: (event_id, handler_name)
- สร้างตาราง pdpa_purposes สำหรับ master data
  - Seed: NECESSARY, ANALYTICS, MARKETING, ACCOUNT_SYSTEM, USAGE_LOGS, TRANSACTION_HISTORY
- สร้าง Trigger สำหรับ auto-update updated_at
- Seed ข้อมูล Privacy Policy เวอร์ชันแรก (v1.0.0)
```

---

### PDPA-2: [Repo+Service+Handler] Value Objects และ Domain Errors

```
📁 file: internal/modules/pdpa/domain/valueobject/consent_purpose.go, 
         internal/modules/pdpa/domain/valueobject/consent_status.go,
         internal/modules/pdpa/domain/valueobject/dsar_type.go,
         internal/modules/pdpa/domain/valueobject/dsar_status.go,
         internal/modules/pdpa/domain/valueobject/account_status.go,
         internal/modules/pdpa/domain/errors/errors.go
📥 Input: string (raw purpose/status/type)
📤 Output: Value Object, error
🔁 Flow: -
🧪 Unit: unit test & mockery
📘🔗 Swagger: -
🎯 Objective:
- repo: -
- service: -
- handler: -
- valueobject:
  - ConsentPurpose (string-based): NECESSARY, ANALYTICS, MARKETING, ACCOUNT_SYSTEM, USAGE_LOGS, TRANSACTION_HISTORY
    - ParseConsentPurpose(s string) (ConsentPurpose, error)
    - IsValid() bool, IsRequired() bool, String() string, MarshalJSON(), UnmarshalJSON()
  - ConsentStatus: GRANTED, REVOKED, EXPIRED, DELETED
    - IsValid(), IsActive(), String(), MarshalJSON()
  - DSARType: ACCESS, ERASURE, WITHDRAW_CONSENT
    - ParseDSARType(), IsValid(), String(), MarshalJSON()
  - DSARStatus: PENDING, PROCESSING, COMPLETED, REJECTED
    - IsValid(), IsTerminal(), String(), MarshalJSON()
  - AccountStatus: ACTIVE, SUSPENDED, TERMINATED, DELETED
    - IsValid(), IsActive(), IsSuspended(), IsTerminated(), IsDeleted(), String(), MarshalJSON()
- errors:
  - Sentinel errors: ErrConsentNotFound, ErrConsentAlreadyRevoked, ErrConsentExpired, ErrInvalidPurpose, ErrRevokeNotAllowed
  - ErrDSARNotFound, ErrDSARAlreadyProcessed, ErrDSARNotPending, ErrOTPInvalid, ErrOTPExpired, ErrRateLimitExceeded
  - ErrAccountNotFound, ErrAccountNotActive, ErrAccountNotSuspended, ErrAccountNotTerminated, ErrDeletionNotConfirmed
  - ErrImmediateDeletionNotAllowed, ErrAlreadyDeleted, ErrPolicyNotFound, ErrPolicyVersionExists
  - ErrInvalidInput, ErrUnauthorized, ErrInternal
```

---

### PDPA-3: [Repo+Service+Handler] Domain Events และ Event Metadata

```
📁 file: internal/modules/pdpa/domain/event/event_metadata.go,
         internal/modules/pdpa/domain/event/consent_events.go,
         internal/modules/pdpa/domain/event/dsar_events.go,
         internal/modules/pdpa/domain/event/account_events.go,
         internal/modules/pdpa/domain/event/data_events.go
📥 Input: -
📤 Output: Event structs, Envelope
🔁 Flow: -
🧪 Unit: unit test & mockery
📘🔗 Swagger: -
🎯 Objective:
- repo: -
- service: -
- handler: -
- event:
  - EventMetadata: EventID, CorrelationID, CausationID, TraceID, OccurredAt, Version, Source
    - NewMetadata(correlationID, causationID, traceID string) EventMetadata
  - Envelope: EventType, Metadata, Payload
    - NewEnvelope(eventType string, meta EventMetadata, payload interface{}) Envelope
  - Kafka Topics constants:
    - TopicConsentGranted, TopicConsentRevoked
    - TopicDSARSubmitted, TopicDSARProcessing, TopicDSAROTPVerified, TopicDSARCompleted, TopicDSARRejected
    - TopicAccountSuspended, TopicAccountTerminated
    - TopicDeletionRequested, TopicDataDeleted
    - TopicLLMAnalysis, TopicEmailNotification, TopicBlockchainRecord, TopicPolicyPublished
    - TopicDLQPrefix
  - Typed Payloads:
    - ConsentGrantedPayload, ConsentRevokedPayload
    - DSARSubmittedPayload, DSARProcessingPayload, DSAROTPVerifiedPayload, DSARCompletedPayload, DSARRejectedPayload
    - AccountSuspendedPayload, AccountTerminatedPayload
    - DataDeletionRequestedPayload, DataDeletedPayload
    - EmailNotificationPayload, BlockchainRecordPayload, PolicyPublishedPayload
```

---

### PDPA-4: [Repo+Service+Handler] Domain Entities (Aggregate Roots)

```
📁 file: internal/modules/pdpa/domain/entity/consent_log.go,
         internal/modules/pdpa/domain/entity/dsar_request.go,
         internal/modules/pdpa/domain/entity/user_account_status.go,
         internal/modules/pdpa/domain/entity/audit_trail.go,
         internal/modules/pdpa/domain/entity/privacy_policy.go
📥 Input: -
📤 Output: Entity structs, error
🔁 Flow: -
🧪 Unit: unit test & mockery
📘🔗 Swagger: -
🎯 Objective:
- repo: -
- service: -
- handler: -
- entity:
  - ConsentLog (Aggregate Root):
    - Fields: ID, UserID, SessionID, Purpose, Status, ConsentedAt, ExpiresAt, RevokedAt, AutoDeletedAt, IPAddress, UserAgent
    - NewConsentLog(userID, purpose, now, ip, userAgent) (*ConsentLog, error) — inject clock
    - Revoke(now time.Time) error — ตรวจสอบ status + IsRequired
    - MarkDeleted(now time.Time)
    - IsActive(now time.Time) bool
  - DSARRequest (Aggregate Root):
    - Fields: ID, UserID, RequestType, Status, RequestedAt, CompletedAt, DataPayload, DataHash, BlockchainTx, RejectionReason, OTPHash, OTPExpiresAt, OTPVerified, OTPAttempts, IPAddress, UserAgent
    - NewDSARRequest(userID, requestType, now, ip, userAgent) (*DSARRequest, string, error) — คืน plaintext OTP
    - VerifyOTP(otp string, now time.Time) error — hash comparison, attempt counting
    - MarkProcessing() error — requires OTPVerified
    - MarkCompleted(payload []byte, dataHash, txHash string, now time.Time) error
    - MarkRejected(reason string) error
    - generateOTP() — crypto/rand 6 digits
    - hashOTP(otp string) — SHA-256
  - UserAccountStatus (Aggregate Root):
    - Fields: ID, UserID, Status, SuspendedAt, TerminatedAt, DeletionConfirmedAt, RetentionDeadline, AutoDeletedAt, UpdatedAt
    - NewUserAccountStatus(userID, now) (*UserAccountStatus, error)
    - Suspend(now time.Time, retentionYears int) error
    - Terminate(now time.Time) error
    - ConfirmDeletion(now time.Time) error — requires TERMINATED
    - MarkDeleted(now time.Time)
    - CanImmediateDelete() bool
    - IsReadyForAutoDeletion(now time.Time) bool
  - AuditTrail:
    - Fields: ID, UserID, Action, Details, IPAddress, UserAgent, CreatedAt
    - NewAuditTrail(userID, action, details) *AuditTrail
    - NewAuditTrailWithMeta(userID, action, details, ip, userAgent) *AuditTrail
    - Action constants: CONSENT_GRANTED, CONSENT_REVOKED, DSAR_SUBMITTED, ...
  - PrivacyPolicy:
    - Fields: ID, Version, Title, Content, EffectiveDate, IsActive, CreatedAt, UpdatedAt
    - NewPrivacyPolicy(version, title, content, effectiveDate, now) (*PrivacyPolicy, error)
    - Activate(now time.Time)
    - Deactivate(now time.Time)
```

---

### PDPA-5: [Repo+Service+Handler] Domain Services

```
📁 file: internal/modules/pdpa/domain/service/consent_validator.go,
         internal/modules/pdpa/domain/service/deletion_policy_service.go,
         internal/modules/pdpa/domain/service/anonymization_service.go,
         internal/modules/pdpa/domain/service/blockchain_service.go
📥 Input: -
📤 Output: -
🔁 Flow: -
🧪 Unit: unit test & mockery
📘🔗 Swagger: -
🎯 Objective:
- repo: -
- service:
  - ConsentValidator:
    - ValidatePurposes(purposes map[string]bool) ([]ConsentPurpose, error) — ต้องมี NECESSARY
    - ValidateRevoke(purpose ConsentPurpose) error — ห้าม revoke NECESSARY
  - DeletionPolicyService:
    - RetentionPolicy: SuspendedRetentionYears, RevokedConsentDays
    - DefaultRetentionPolicy() RetentionPolicy
    - CanImmediateDeletion(status *UserAccountStatus) bool
    - IsReadyForAutoDeletion(status *UserAccountStatus, now time.Time) bool
    - CalculateRetentionDeadline(suspendedAt time.Time) time.Time
    - ConsentAutoDeleteCutoff(now time.Time) time.Time
  - AnonymizationService:
    - AnonymizeEmail(email string) string — keep domain, hash local part
    - AnonymizePhone(phone string) string — keep last 4 digits
    - AnonymizeIP(ip string) string — mask last octet
    - hash(input string) string — SHA-256 + salt
  - BlockchainService (interface):
    - RecordHash(ctx context.Context, hash string) (txHash string, err error)
- handler: -
```

---

### PDPA-6: [Repo+Service+Handler] Repository Interfaces และ Cache Interface

```
📁 file: internal/modules/pdpa/domain/repository/consent_repo.go,
         internal/modules/pdpa/domain/repository/dsar_repo.go,
         internal/modules/pdpa/domain/repository/account_status_repo.go,
         internal/modules/pdpa/domain/repository/audit_repo.go,
         internal/modules/pdpa/domain/repository/outbox_repo.go,
         internal/modules/pdpa/domain/repository/cache.go
📥 Input: -
📤 Output: Repository interfaces
🔁 Flow: -
🧪 Unit: unit test & mockery
📘🔗 Swagger: -
🎯 Objective:
- repo:
  - ConsentRepository interface:
    - Save(ctx, *ConsentLog) error
    - Update(ctx, *ConsentLog) error
    - FindByID(ctx, id) (*ConsentLog, error)
    - FindByUserID(ctx, userID) ([]ConsentLog, error)
    - FindLatestByUserAndPurpose(ctx, userID, purpose) (*ConsentLog, error)
    - FindActiveByUserAndPurpose(ctx, userID, purpose) (*ConsentLog, error)
    - FindReadyForAutoDeletion(ctx, cutoff, limit) ([]ConsentLog, error)
    - DeleteByUserID(ctx, userID) error
    - DeleteByID(ctx, id) error
    - IsConsentActive(ctx, userID, purpose) (bool, error)
  - DSARRepository interface:
    - Save, Update, FindByID, FindByUserID, FindByStatus
    - CountByUserAndTypeInPeriod(ctx, userID, type, from) (int64, error)
    - DeleteByUserID(ctx, userID) error
  - UserAccountStatusRepository interface:
    - Save, FindByUserID, FindReadyForAutoDeletion, ExistsByUserID
  - AuditRepository interface:
    - Save, FindByFilter(ctx, AuditFilter) ([]AuditTrail, error)
    - CountByAction(ctx, from, to) (map[string]int64, error)
  - PolicyRepository interface:
    - Save, Update, FindByID, FindActive, FindAll, FindByVersion, DeactivateAll
  - OutboxRepository interface:
    - Save, FetchPending, MarkPublished, MarkFailed, CleanupPublished
  - UnitOfWork interface:
    - Do(ctx, fn func(ctx) error) error
- service: -
- handler: -
- cache:
  - Cache interface:
    - Set(ctx, userID, purpose, status) error
    - Get(ctx, userID, purpose) (ConsentStatus, error)
    - DeleteAllByUser(ctx, userID) error
```

---

### PDPA-7: [Repo+Service+Handler] Infrastructure - PostgreSQL Models และ UnitOfWork

```
📁 file: internal/modules/pdpa/infrastructure/persistence/postgres/models.go,
         internal/modules/pdpa/infrastructure/persistence/postgres/uow.go
📥 Input: -
📤 Output: GORM models, UnitOfWork implementation
🔁 Flow: -
🧪 Unit: unit test & mockery
📘🔗 Swagger: -
🎯 Objective:
- repo: -
- service: -
- handler: -
- models:
  - PdpAConsentModel (TableName: pdpa_consents)
  - PdpADSARModel (TableName: pdpa_dsar_requests)
  - PdpAAccountStatusModel (TableName: pdpa_user_account_statuses)
  - PdpAAuditTrailModel (TableName: pdpa_audit_trails)
  - PdpAPolicyModel (TableName: pdpa_privacy_policies)
  - PdpAOutboxModel (TableName: pdpa_outbox)
- uow:
  - UnitOfWorkImpl struct { db *gorm.DB }
  - NewUnitOfWork(db *gorm.DB) *UnitOfWorkImpl
  - Do(ctx, fn) error — GORM transaction wrapper
  - GetDB(ctx, fallback) *gorm.DB — ดึง tx จาก context ถ้ามี
```

---

### PDPA-8: [Repo+Service+Handler] Infrastructure - Consent Repository Implementation

```
📁 file: internal/modules/pdpa/infrastructure/persistence/postgres/consent_repo_impl.go
📥 Input: Entity models, context
📤 Output: Database operations
🔁 Flow: -
🧪 Unit: unit test & mockery
📘🔗 Swagger: -
🎯 Objective:
- repo:
  - ConsentRepoImpl struct { db *gorm.DB }
  - NewConsentRepository(db *gorm.DB) *ConsentRepoImpl
  - Save(ctx, *ConsentLog) — INSERT
  - Update(ctx, *ConsentLog) — UPDATE status, revoked_at, auto_deleted_at (คืน ErrConsentNotFound ถ้า RowsAffected=0)
  - FindByID — SELECT WHERE id
  - FindByUserID — SELECT WHERE user_id ORDER BY granted_at DESC
  - FindLatestByUserAndPurpose — SELECT ORDER BY granted_at DESC LIMIT 1
  - FindActiveByUserAndPurpose — SELECT WHERE status=GRANTED AND expires_at > now()
  - FindReadyForAutoDeletion — SELECT WHERE status IN (REVOKED,EXPIRED) AND revoked_at <= cutoff LIMIT
  - DeleteByUserID — DELETE WHERE user_id
  - DeleteByID — DELETE WHERE id
  - IsConsentActive — COUNT WHERE status=GRANTED AND expires_at > now()
  - toConsentModel / toConsentEntity / toConsentEntities mappers
- service: -
- handler: -
```

---

### PDPA-9: [Repo+Service+Handler] Infrastructure - DSAR Repository Implementation

```
📁 file: internal/modules/pdpa/infrastructure/persistence/postgres/dsar_repo_impl.go
📥 Input: Entity models, context
📤 Output: Database operations
🔁 Flow: -
🧪 Unit: unit test & mockery
📘🔗 Swagger: -
🎯 Objective:
- repo:
  - DSARRepoImpl struct { db *gorm.DB }
  - NewDSARRepository(db *gorm.DB) *DSARRepoImpl
  - Save(ctx, *DSARRequest) — INSERT
  - Update(ctx, *DSARRequest) — UPDATE status, completed_at, data_payload, data_hash, blockchain_tx, rejection_reason, otp_verified, otp_attempts
  - FindByID — SELECT WHERE id
  - FindByUserID — SELECT WHERE user_id ORDER BY requested_at DESC
  - FindByStatus — SELECT WHERE status ORDER BY requested_at ASC LIMIT OFFSET
  - CountByUserAndTypeInPeriod — COUNT WHERE user_id AND request_type AND requested_at >= from
  - DeleteByUserID — DELETE WHERE user_id
  - toDSARModel / toDSAREntity / toDSAREntities mappers
- service: -
- handler: -
```

---

### PDPA-10: [Repo+Service+Handler] Infrastructure - Account Status, Audit, Policy, Outbox Repositories

```
📁 file: internal/modules/pdpa/infrastructure/persistence/postgres/account_status_repo_impl.go,
         internal/modules/pdpa/infrastructure/persistence/postgres/audit_repo_impl.go,
         internal/modules/pdpa/infrastructure/persistence/postgres/policy_repo_impl.go,
         internal/modules/pdpa/infrastructure/persistence/postgres/outbox_repo_impl.go
📥 Input: Entity models, context
📤 Output: Database operations
🔁 Flow: -
🧪 Unit: unit test & mockery
📘🔗 Swagger: -
🎯 Objective:
- repo:
  - AccountStatusRepoImpl:
    - Save — UPSERT (ON CONFLICT user_id)
    - FindByUserID, FindReadyForAutoDeletion, ExistsByUserID
  - AuditRepoImpl:
    - Save — INSERT with JSON marshal details
    - FindByFilter — SELECT with dynamic WHERE (user_id, action, from, to)
    - CountByAction — GROUP BY action
  - PolicyRepoImpl:
    - Save, Update, FindByID, FindActive (WHERE is_active=true), FindAll, FindByVersion, DeactivateAll
  - OutboxRepoImpl:
    - Save — INSERT
    - FetchPending — SELECT with FOR UPDATE SKIP LOCKED (multi-replica safe)
    - MarkPublished — UPDATE status=PUBLISHED, published_at
    - MarkFailed — UPDATE status=FAILED, last_error, retry_count+1
    - CleanupPublished — DELETE WHERE status=PUBLISHED AND published_at < cutoff
- service: -
- handler: -
```

---

### PDPA-11: [Repo+Service+Handler] Infrastructure - Redis Cache และ Idempotency Store

```
📁 file: internal/modules/pdpa/infrastructure/cache/redis/consent_cache.go,
         internal/modules/pdpa/infrastructure/cache/redis/idempotency_store.go
📥 Input: -
📤 Output: Cache operations
🔁 Flow: -
🧪 Unit: unit test & mockery
📘🔗 Swagger: -
🎯 Objective:
- repo: -
- service: -
- handler: -
- cache:
  - ConsentCacheImpl:
    - NewConsentCache(client *redis.Client, ttl time.Duration) *ConsentCacheImpl
    - key(userID, purpose) — "pdpa:consent:{userID}:{purpose}"
    - Set(ctx, userID, purpose, status) — SETEX
    - Get(ctx, userID, purpose) — GET (คืน "", nil ถ้า redis.Nil)
    - DeleteAllByUser(ctx, userID) — SCAN + DEL pattern
  - IdempotencyStore:
    - NewIdempotencyStore(client *redis.Client) *IdempotencyStore
    - key(eventID, handler) — "pdpa:idem:{handler}:{eventID}"
    - IsProcessed(ctx, eventID, handler) — EXISTS
    - MarkProcessed(ctx, eventID, handler, ttl) — SET NX (atomic)
```

---

### PDPA-12: [Repo+Service+Handler] Infrastructure - Kafka Producer และ Outbox Publisher

```
📁 file: internal/modules/pdpa/infrastructure/messaging/kafka/producer.go,
         internal/modules/pdpa/infrastructure/messaging/kafka/event_serializer.go,
         internal/modules/pdpa/infrastructure/messaging/outbox/outbox_publisher.go
📥 Input: Event payloads, context
📤 Output: Kafka messages
🔁 Flow: -
🧪 Unit: unit test & mockery
📘🔗 Swagger: -
🎯 Objective:
- repo: -
- service: -
- handler: -
- producer:
  - Producer interface: Publish(ctx, topic, key string, payload []byte) error, Close() error
  - SyncProducerImpl (sarama):
    - NewSyncProducer(brokers []string, log logger.Logger) (*SyncProducerImpl, error)
    - Publish — SendMessage with idempotent=true, WaitForAll
- serializer:
  - SerializeEnvelope(env event.Envelope) ([]byte, error)
  - DeserializeEnvelope(data []byte) (*event.Envelope, error)
- outbox publisher:
  - OutboxPublisher:
    - NewOutboxPublisher(outboxRepo, producer, interval, batchSize, logger)
    - Run(ctx) — ticker loop
    - drain(ctx) — FetchPending → Publish → MarkPublished / MarkFailed
    - Cleanup(ctx) — ลบ PUBLISHED events เก่ากว่า 7 วัน
```

---

### PDPA-13: [Repo+Service+Handler] Infrastructure - External Clients (Email, LLM, Blockchain)

```
📁 file: internal/modules/pdpa/infrastructure/external/email_client.go,
         internal/modules/pdpa/infrastructure/external/llm_client.go,
         internal/modules/pdpa/infrastructure/external/blockchain_client.go
📥 Input: -
📤 Output: -
🔁 Flow: -
🧪 Unit: unit test & mockery
📘🔗 Swagger: -
🎯 Objective:
- repo: -
- service: -
- handler: -
- email client:
  - SMTPEmailClient:
    - NewEmailClient(from string, log logger.Logger) *SMTPEmailClient
    - Send(ctx, to, subject, template, locale string, vars map[string]interface{}) error
- llm client:
  - HTTPLLMClient:
    - NewLLMClient(baseURL, apiKey, model string, log logger.Logger) *HTTPLLMClient
    - AnalyzeDataWithContext(ctx, analysisType string, payload map[string]interface{}) (string, error)
- blockchain client:
  - MockBlockchainClient:
    - NewBlockchainClient(log logger.Logger) *MockBlockchainClient
    - Record(ctx, userID, action, dataHash string) (string, error)
```

---

### PDPA-14: [Repo+Service+Handler] Application Layer - DTOs

```
📁 file: internal/modules/pdpa/application/dto/dto.go
📥 Input: -
📤 Output: DTO structs
🔁 Flow: -
🧪 Unit: unit test & mockery
📘🔗 Swagger: -
🎯 Objective:
- repo: -
- service: -
- handler: -
- dto:
  - Consent:
    - GrantConsentRequest: UserID, Purpose, IPAddress, UserAgent, CorrelationID, TraceID
    - RevokeConsentRequest: UserID, Purpose, IPAddress, UserAgent, CorrelationID, TraceID
    - ConsentStatusResponse: UserID, Purpose, Status, IsActive, GrantedAt, ExpiresAt, RevokedAt
    - ConsentStatusListResponse: Consents[], Total
  - DSAR:
    - SubmitDSARRequest: UserID, RequestType, IPAddress, UserAgent, CorrelationID, TraceID
    - VerifyOTPRequest: UserID, DSARRequestID, OTPCode, CorrelationID, TraceID
    - ProcessDSARRequest: UserID, DSARRequestID, CorrelationID, TraceID
    - CompleteDSARRequest: UserID, DSARRequestID, ResponseData, CorrelationID, TraceID
    - RejectDSARRequest: UserID, DSARRequestID, Reason, CorrelationID, TraceID
    - DSARStatusResponse: ID, UserID, RequestType, Status, RequestedAt, CompletedAt, RejectionReason
    - DSARStatusListResponse: Requests[], Total
  - Account:
    - SuspendAccountRequest: UserID, RetentionYears, CorrelationID, TraceID
    - TerminateAccountRequest: UserID, CorrelationID, TraceID
    - ConfirmDeletionRequest: UserID, Reason, CorrelationID, TraceID
    - UserAccountStatusResponse: UserID, Status, RetentionDeadline, DeletionConfirmedAt, AutoDeletedAt
  - Policy:
    - PublishPolicyRequest: Version, Title, Content, EffectiveDate, CorrelationID, TraceID
    - PrivacyPolicyResponse: ID, Version, Title, Content, EffectiveDate, IsActive, UpdatedAt
  - Admin:
    - AdminReportResponse: From, To, ConsentByPurpose, ActionsByType, GeneratedAt
    - AuditTrailItem: ID, UserID, Action, Details, IPAddress, UserAgent, CreatedAt
```

---

### PDPA-15: [Repo+Service+Handler] Application Layer - Consent Commands

```
📁 file: internal/modules/pdpa/application/command/record_consent.go,
         internal/modules/pdpa/application/command/revoke_consent.go
📥 Input: GrantConsentRequest, RevokeConsentRequest
📤 Output: error
🔁 Flow: Handler → Validate → Entity → UnitOfWork (Save + Outbox) → Kafka
🧪 Unit: unit test & mockery
📘🔗 Swagger: เพิ่ม API documentation & endpoint test
🎯 Objective:
- repo:
  - ใช้ ConsentRepository.Save/Update, OutboxRepository.Save, UnitOfWork.Do
- service:
  - RecordConsentHandler:
    - Parse purpose → NewConsentLog → NewEnvelope(TopicConsentGranted)
    - UnitOfWork.Do: ConsentRepo.Save + OutboxRepo.Save (atomic)
  - RevokeConsentHandler:
    - Parse purpose → ตรวจ IsRequired → FindActiveByUserAndPurpose → consent.Revoke(now)
    - NewEnvelope(TopicConsentRevoked)
    - UnitOfWork.Do: ConsentRepo.Update + OutboxRepo.Save
- handler:
  - RecordConsentHandler.Handle(ctx, GrantConsentRequest) error
  - RevokeConsentHandler.Handle(ctx, RevokeConsentRequest) error
```

---

### PDPA-16: [Repo+Service+Handler] Application Layer - DSAR Commands

```
📁 file: internal/modules/pdpa/application/command/submit_dsar.go,
         internal/modules/pdpa/application/command/process_dsar.go,
         internal/modules/pdpa/application/command/verify_dsar_otp.go,
         internal/modules/pdpa/application/command/complete_dsar.go,
         internal/modules/pdpa/application/command/reject_dsar.go
📥 Input: SubmitDSARRequest, ProcessDSARRequest, VerifyOTPRequest, CompleteDSARRequest, RejectDSARRequest
📤 Output: *DSARRequest, error
🔁 Flow: Handler → Validate → Rate Limit → Entity → UnitOfWork (Save + Outbox) → Kafka
🧪 Unit: unit test & mockery
📘🔗 Swagger: เพิ่ม API documentation & endpoint test
🎯 Objective:
- repo:
  - ใช้ DSARRepository.Save/Update, OutboxRepository.Save, UnitOfWork.Do
- service:
  - SubmitDSARHandler:
    - Rate limit check (3/day/user/type)
    - NewDSARRequest → คืน plaintext OTP
    - Outbox: TopicDSARSubmitted + TopicEmailNotification (OTP)
  - ProcessDSARHandler:
    - FindByID → dsar.MarkProcessing() → DSARRepo.Update + Outbox(TopicDSARProcessing)
  - VerifyDSAROTPHandler:
    - Redis rate limit → FindByID → dsar.VerifyOTP() → DSARRepo.Update + Outbox(TopicDSAROTPVerified)
  - CompleteDSARHandler:
    - FindByID → hash payload → blockchain.RecordHash → dsar.MarkCompleted()
    - DSARRepo.Update + Outbox(TopicDSARCompleted)
  - RejectDSARHandler:
    - FindByID → dsar.MarkRejected(reason) → DSARRepo.Update + Outbox(TopicDSARRejected)
- handler:
  - แต่ละ handler มี Handle(ctx, request) error
```

---

### PDPA-17: [Repo+Service+Handler] Application Layer - Account Commands

```
📁 file: internal/modules/pdpa/application/command/handle_account_suspended.go,
         internal/modules/pdpa/application/command/handle_account_terminated.go,
         internal/modules/pdpa/application/command/confirm_deletion.go,
         internal/modules/pdpa/application/command/immediate_deletion.go,
         internal/modules/pdpa/application/command/auto_delete_expired.go
📥 Input: SuspendAccountRequest, TerminateAccountRequest, ConfirmDeletionRequest
📤 Output: error, *AutoDeleteExpiredResult
🔁 Flow: Handler → Find/Create Entity → Entity Method → UnitOfWork (Save + Outbox) → Kafka
🧪 Unit: unit test & mockery
📘🔗 Swagger: เพิ่ม API documentation & endpoint test
🎯 Objective:
- repo:
  - ใช้ UserAccountStatusRepository, ConsentRepository, DSARRepository, OutboxRepository, UnitOfWork
- service:
  - HandleAccountSuspendedHandler:
    - Find/Create UserAccountStatus → Suspend(now, retentionYears)
    - Save + Outbox(TopicAccountSuspended)
  - HandleAccountTerminatedHandler:
    - Find/Create → Terminate(now) → Save + Outbox(TopicAccountTerminated)
  - ConfirmDeletionHandler:
    - FindByUserID → ConfirmDeletion(now) → Save + Outbox(TopicDeletionRequested)
  - ImmediateDeletionHandler:
    - FindByUserID → MarkDeleted(now) → UnitOfWork.Do: Delete consents + Delete dsars + Save + Outbox(TopicDataDeleted)
  - AutoDeleteExpiredHandler:
    - Phase 1: ConsentRepo.FindReadyForAutoDeletion → MarkDeleted → Update
    - Phase 2: AccountRepo.FindReadyForAutoDeletion → deleteUserData (Delete consents + dsars + Save + Outbox)
    - คืน AutoDeleteExpiredResult { AccountsDeleted, ConsentsDeleted, Errors }
- handler:
  - แต่ละ handler มี Handle(ctx, request) error
```

---

### PDPA-18: [Repo+Service+Handler] Application Layer - Policy Command

```
📁 file: internal/modules/pdpa/application/command/publish_privacy_policy.go
📥 Input: PublishPolicyRequest
📤 Output: *PrivacyPolicy, error
🔁 Flow: Handler → Check version exists → NewPrivacyPolicy → Activate → UnitOfWork (DeactivateAll + Save + Outbox)
🧪 Unit: unit test & mockery
📘🔗 Swagger: เพิ่ม API documentation & endpoint test
🎯 Objective:
- repo:
  - ใช้ PolicyRepository (FindByVersion, DeactivateAll, Save), OutboxRepository, UnitOfWork
- service:
  - PublishPrivacyPolicyHandler:
    - ตรวจ version ซ้ำ → NewPrivacyPolicy → Activate(now)
    - UnitOfWork.Do: DeactivateAll + Save + Outbox(TopicPolicyPublished)
- handler:
  - Handle(ctx, PublishPolicyRequest) (*PrivacyPolicy, error)
```

---

### PDPA-19: [Repo+Service+Handler] Application Layer - Queries

```
📁 file: internal/modules/pdpa/application/query/get_consent_history.go,
         internal/modules/pdpa/application/query/get_dsar_status.go,
         internal/modules/pdpa/application/query/get_admin_report.go,
         internal/modules/pdpa/application/query/get_active_privacy_policy.go
📥 Input: userID, dsarID, from, to
📤 Output: DTOs
🔁 Flow: Handler → Repository → DTO mapping
🧪 Unit: unit test & mockery
📘🔗 Swagger: เพิ่ม API documentation & endpoint test
🎯 Objective:
- repo:
  - ใช้ ConsentRepository.FindByUserID
  - DSARRepository.FindByID, FindByUserID
  - AuditRepository.CountByAction, FindByFilter
  - PolicyRepository.FindActive, FindAll
- service:
  - GetConsentHistoryHandler.Handle(ctx, userID) (*ConsentStatusListResponse, error)
  - GetDSARStatusHandler.Handle(ctx, id) (*DSARStatusResponse, error)
  - ListDSARHandler.Handle(ctx, userID) (*DSARStatusListResponse, error)
  - GetAdminReportHandler.Handle(ctx, from, to) (*AdminReportResponse, error)
  - ListAuditTrailsHandler.Handle(ctx, filter) ([]AuditTrailItem, error)
  - GetActivePrivacyPolicyHandler.Handle(ctx) (*PrivacyPolicyResponse, error)
  - ListPoliciesHandler.Handle(ctx) ([]PrivacyPolicyResponse, error)
- handler: -
```

---

### PDPA-20: [Repo+Service+Handler] Application Layer - Event Handlers Base

```
📁 file: internal/modules/pdpa/application/event_handler/base_handler.go,
         internal/modules/pdpa/application/event_handler/pii_redactor.go
📥 Input: Envelope
📤 Output: error
🔁 Flow: -
🧪 Unit: unit test & mockery
📘🔗 Swagger: -
🎯 Objective:
- repo:
  - ใช้ AuditRepository.Save
- service:
  - Envelope struct: EventType, Metadata, Payload (json.RawMessage)
  - EventMetadata struct: EventID, CorrelationID, CausationID, TraceID, OccurredAt, Version, Source
  - DecodeEnvelope(data []byte) (*Envelope, error)
  - DecodePayload[T any](e *Envelope) (*T, error)
  - IdempotencyStore interface: IsProcessed, MarkProcessed
  - Metrics interface: IncHandlerProcessed, ObserveHandlerDuration
  - NopMetrics() Metrics
  - Handler interface: Handle(ctx, *Envelope) error, Name() string
  - BaseHandler struct: auditRepo, idempotency, metrics, redactor, logger, name
  - NewBaseHandler(name, auditRepo, idem, metrics, redactor, logger) BaseHandler
  - IsProcessed(ctx, eventID) (bool, error)
  - MarkProcessed(ctx, eventID, ttl) error
  - SaveAudit(ctx, userID, action, payload) error — PII redacted
  - Observe(ctx, meta, fn) error — metrics + structured logging
  - PIIRedactor:
    - NewPIIRedactor() *PIIRedactor
    - maskKeys: email, phone, otp, password, token, national_id, credit_card, ...
    - patterns: email regex, Thai phone, Thai national ID, credit card
    - RedactPayload(payload map[string]interface{}) map[string]interface{}
    - RedactJSONString(v interface{}) string
- handler: -
```

---

### PDPA-21: [Repo+Service+Handler] Application Layer - Event Handlers (Consent & DSAR)

```
📁 file: internal/modules/pdpa/application/event_handler/consent_granted_handler.go,
         internal/modules/pdpa/application/event_handler/consent_revoked_handler.go,
         internal/modules/pdpa/application/event_handler/dsar_submitted_handler.go,
         internal/modules/pdpa/application/event_handler/dsar_completed_handler.go
📥 Input: Envelope (from Kafka)
📤 Output: error
🔁 Flow: Decode → Idempotency Check → Business Logic → Audit → Mark Processed
🧪 Unit: unit test & mockery
📘🔗 Swagger: -
🎯 Objective:
- repo:
  - ใช้ AuditRepository.Save, Cache.Set, Kafka Producer
- service:
  - ConsentGrantedHandler:
    - Handle(ctx, *Envelope) error
    - Steps: Idempotency → cache.Set(GRANTED) → Publish(TopicEmailNotification) → SaveAudit → MarkProcessed
  - ConsentRevokedHandler:
    - Handle(ctx, *Envelope) error
    - Steps: Idempotency → cache.Set(REVOKED) → Publish(TopicEmailNotification) → SaveAudit → MarkProcessed
  - DSARSubmittedHandler:
    - Handle(ctx, *Envelope) error
    - Steps: Idempotency → Publish(TopicLLMAnalysis) for ACCESS/ERASURE → SaveAudit → MarkProcessed
  - DSARCompletedHandler:
    - Handle(ctx, *Envelope) error
    - Steps: Idempotency → Publish(TopicEmailNotification) → Publish(TopicBlockchainRecord) → SaveAudit → MarkProcessed
- handler: -
```

---

### PDPA-22: [Repo+Service+Handler] Application Layer - Event Handlers (Account, Data, LLM, Email, Blockchain)

```
📁 file: internal/modules/pdpa/application/event_handler/account_suspended_handler.go,
         internal/modules/pdpa/application/event_handler/account_terminated_handler.go,
         internal/modules/pdpa/application/event_handler/data_deletion_handler.go,
         internal/modules/pdpa/application/event_handler/llm_analysis_handler.go,
         internal/modules/pdpa/application/event_handler/email_notification_handler.go,
         internal/modules/pdpa/application/event_handler/blockchain_record_handler.go
📥 Input: Envelope (from Kafka)
📤 Output: error
🔁 Flow: Decode → Idempotency Check → Business Logic → Audit → Mark Processed
🧪 Unit: unit test & mockery
📘🔗 Swagger: -
🎯 Objective:
- repo:
  - ใช้ AuditRepository.Save, Kafka Producer, External Clients (LLM, Email, Blockchain)
- service:
  - AccountSuspendedHandler:
    - Handle: Idempotency → Publish(TopicEmailNotification) → SaveAudit → MarkProcessed
  - AccountTerminatedHandler:
    - Handle: Idempotency → Publish(TopicEmailNotification) → SaveAudit → MarkProcessed
  - DataDeletionHandler:
    - Handle: Idempotency → Publish(TopicDataDeleted) if immediate → Publish(TopicBlockchainRecord) → SaveAudit → MarkProcessed
  - LLMAnalysisHandler:
    - Handle: Idempotency → llm.AnalyzeDataWithContext() → SaveAudit → MarkProcessed
  - EmailNotificationHandler:
    - Handle: Idempotency → email.Send() → MarkProcessed
  - BlockchainRecordHandler:
    - Handle: Idempotency → client.Record() → SaveAudit → MarkProcessed
- handler: -
```

---

### PDPA-23: [Repo+Service+Handler] Interfaces Layer - HTTP DTOs และ Handlers

```
📁 file: internal/modules/pdpa/interfaces/http/dto.go,
         internal/modules/pdpa/interfaces/http/consent_handler.go,
         internal/modules/pdpa/interfaces/http/dsar_handler.go,
         internal/modules/pdpa/interfaces/http/account_handler.go,
         internal/modules/pdpa/interfaces/http/policy_handler.go,
         internal/modules/pdpa/interfaces/http/admin_handler.go
📥 Input: HTTP Request (JSON body, query params, path params)
📤 Output: HTTP Response (JSON)
🔁 Flow: HTTP Request → Parse → Call Command/Query Handler → Response
🧪 Unit: unit test & mockery
📘🔗 Swagger: เพิ่ม API documentation & endpoint test
🎯 Objective:
- repo: -
- service: -
- handler:
  - ConsentHandler:
    - Record: POST /api/v1/pdpa/consent
    - Revoke: DELETE /api/v1/pdpa/consent
    - History: GET /api/v1/pdpa/consent/history
  - DSARHandler:
    - Submit: POST /api/v1/pdpa/dsar
    - Get: GET /api/v1/pdpa/dsar/{id}
    - List: GET /api/v1/pdpa/dsar
    - Process: POST /api/v1/pdpa/dsar/{id}/process (admin)
  - AccountHandler:
    - Suspend: POST /api/v1/pdpa/account/suspend (admin)
    - Terminate: POST /api/v1/pdpa/account/terminate (admin)
    - ConfirmDeletion: POST /api/v1/pdpa/account/confirm-deletion (admin)
  - PolicyHandler:
    - GetActive: GET /api/v1/pdpa/policy
    - ListVersions: GET /api/v1/pdpa/policy/versions
    - Publish: POST /api/v1/pdpa/policy (admin)
  - AdminHandler:
    - GetReport: GET /api/v1/pdpa/admin/reports (admin)
    - ListAudit: GET /api/v1/pdpa/admin/audit-trails (admin)
  - extractUserID(r *http.Request) (uuid.UUID, error) helper
```

---

### PDPA-24: [Repo+Service+Handler] Interfaces Layer - Routes และ Module Entry Point

```
📁 file: internal/modules/pdpa/interfaces/http/routes.go,
         internal/modules/pdpa/module.go
📥 Input: -
📤 Output: -
🔁 Flow: -
🧪 Unit: unit test & mockery
📘🔗 Swagger: เพิ่ม API documentation & endpoint test
🎯 Objective:
- repo: -
- service: -
- handler:
  - routes.go:
    - Handlers struct: Consent, DSAR, Account, Policy, Admin
    - Middleware type: func(http.Handler) http.Handler
    - RegisterRoutes(r chi.Router, h *Handlers, auth, adminOnly Middleware)
    - Routes:
      - POST /consent, DELETE /consent, GET /consent/history
      - POST /dsar, GET /dsar, GET /dsar/{id}
      - GET /policy, GET /policy/versions
      - Admin group: POST /dsar/{id}/process, POST /account/suspend, POST /account/terminate, POST /account/confirm-deletion, POST /policy, GET /admin/reports, GET /admin/audit-trails
- module.go:
  - Config struct: KafkaBrokers, AnonymizationSalt, RetentionYears, OutboxInterval
  - Module struct: repos, services, commands, queries, event handlers, infra, handlers
  - NewModule(db *gorm.DB, redisClient *redis.Client, cfg Config, log logger.Logger) (*Module, error)
    - Initialize repos, UoW, services, Kafka producer, outbox publisher, cache, idempotency, external clients
    - Wire commands, queries, event handlers
    - Build HTTP handlers
  - RegisterHTTP(r chi.Router, auth, admin Middleware)
  - StartBackground(ctx context.Context) — go OutboxPublisher.Run(ctx)
  - Shutdown() error — Producer.Close()
```

---

## 📊 สรุป Sub-tasks

| # | Ticket | Layer | คำอธิบาย |
|---|--------|-------|----------|
| 1 | PDPA-1 | Database | Schema + Migration + Seed |
| 2 | PDPA-2 | Domain | Value Objects + Errors |
| 3 | PDPA-3 | Domain | Events + Metadata |
| 4 | PDPA-4 | Domain | Entities (Aggregate Roots) |
| 5 | PDPA-5 | Domain | Domain Services |
| 6 | PDPA-6 | Domain | Repository Interfaces + Cache |
| 7 | PDPA-7 | Infrastructure | PostgreSQL Models + UnitOfWork |
| 8 | PDPA-8 | Infrastructure | Consent Repository Impl |
| 9 | PDPA-9 | Infrastructure | DSAR Repository Impl |
| 10 | PDPA-10 | Infrastructure | Account/Audit/Policy/Outbox Repo Impl |
| 11 | PDPA-11 | Infrastructure | Redis Cache + Idempotency |
| 12 | PDPA-12 | Infrastructure | Kafka Producer + Outbox Publisher |
| 13 | PDPA-13 | Infrastructure | External Clients (Email, LLM, Blockchain) |
| 14 | PDPA-14 | Application | DTOs |
| 15 | PDPA-15 | Application | Consent Commands |
| 16 | PDPA-16 | Application | DSAR Commands |
| 17 | PDPA-17 | Application | Account Commands |
| 18 | PDPA-18 | Application | Policy Command |
| 19 | PDPA-19 | Application | Queries |
| 20 | PDPA-20 | Application | Event Handler Base + PII Redactor |
| 21 | PDPA-21 | Application | Event Handlers (Consent & DSAR) |
| 22 | PDPA-22 | Application | Event Handlers (Account, Data, LLM, Email, Blockchain) |
| 23 | PDPA-23 | Interface | HTTP Handlers |
| 24 | PDPA-24 | Interface | Routes + Module Entry Point |

---
สร้าง  swagger  ทั้งหมดใน PDPA Module