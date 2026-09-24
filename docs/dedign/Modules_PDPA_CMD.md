# Prompt for Generating PDPA Module in Clean Architecture + DDD + Event-Driven Architecture

## ROLE
You are a Senior Go Developer specializing in Clean Architecture + Domain-Driven Design (DDD) and Modular Monolith Pattern with Event-Driven Architecture (EDA) expertise.

## CONTEXT

**Project:** icmongolang  
**Go Version:** 1.22+  
**Architecture:** Clean Architecture + DDD + Event-Driven Architecture  
**Structure:** Modular Monolith (33 modules)  
**Pattern:** Domain Events + Kafka + Async Workers

### Standard Module Structure
```
internal/modules/{module_name}/
├── domain/
│   ├── entity/{entity}.go
│   ├── value_object/{vo}.go
│   ├── repository/{entity}_repository.go
│   ├── event/{event}.go              # Domain Events
│   ├── service/{service}.go
│   └── errors/errors.go
├── application/
│   ├── command/{verb}_{entity}.go    # Commands (Write)
│   ├── query/{verb}_{entity}.go      # Queries (Read)
│   ├── event_handler/{event}_handler.go
│   └── mocks/{entity}_repository_mock.go
├── infrastructure/
│   ├── persistence/postgres/
│   │   ├── models.go
│   │   └── {entity}_repo_impl.go
│   ├── messaging/kafka/
│   │   ├── producer.go
│   │   └── consumer/{consumer}.go
│   └── cache/redis/
│       └── {entity}_cache.go
├── interfaces/
│   ├── http/
│   │   ├── {entity}_handler.go
│   │   ├── dto.go
│   │   └── routes.go
│   └── websocket/
│       └── {entity}_broadcaster.go
└── module.go
```

### Shared Packages
- `icmongolang/pkg/logger` (structured logging)
- `icmongolang/pkg/validator` (input validation)
- `icmongolang/pkg/httputil` (JSON responses)
- `icmongolang/pkg/eventbus` (domain event bus)
- `github.com/google/uuid`
- `github.com/shopspring/decimal`
- `gorm.io/gorm`
- `github.com/go-chi/chi/v5`
- `github.com/IBM/sarama` (Kafka)
- `github.com/go-redis/redis/v8`

---

# TASK

Create a new module named **"pdpa"** (Personal Data Protection Act Compliance System) with full Event-Driven Architecture implementation.

## Business Requirements

### เป้าหมายหลัก (Core Objectives)
1. **บันทึกข้อมูลประวัติความยินยอม** - Record consent history of users
2. **รายงานประวัติการใช้งานระบบตามหลัก PDPA** - Report system usage history per PDPA
3. **แจ้งเตือนทาง Email** - Email notifications for consent changes and DSAR status
4. **การประมวลผลข้อมูลด้วย LLM AI** - LLM AI analysis for data processing reports
5. **การยินยอมรายข้อ** - Granular consent for:
   - การเก็บข้อมูล (Data collection)
   - การลบข้อมูล (Data deletion)
   - ระบบบัญชีผู้ใช้งาน (User account system)
   - ประวัติการใช้งานระบบ (System usage logs)
   - ประวัติการทำธุรกรรม (Transaction history)

### PDPA Compliance Requirements
- **Privacy Policy Management** - Version-controlled policy display on all pages
- **Cookie Consent Banner** - Accept/Reject buttons with equal prominence, no pre-ticked boxes
- **Granular Consent** - Separate purposes: Necessary, Analytics, Marketing
- **DSAR (Data Subject Access Request)** - Access, Erasure, Withdraw Consent
- **Data Encryption** - Password hashing (Bcrypt), PII encryption
- **Privacy by Design** - Role-based access control, audit trails

### Data Retention Policy
- Suspended accounts: Retain for 1 year, then auto-delete
- Terminated accounts: Immediate deletion upon confirmation
- Consent logs: Retain per legal requirements, anonymize when required

---

## Entities

### 1. ConsentLog (Aggregate Root)
```go
type ConsentLog struct {
    ID            uuid.UUID
    UserID        uuid.UUID
    SessionID     string
    Purpose       ConsentPurpose    // VO: NECESSARY, ANALYTICS, MARKETING
    Status        ConsentStatus     // VO: GRANTED, REVOKED, EXPIRED, DELETED
    IPAddress     string
    UserAgent     string
    GrantedAt     time.Time
    ExpiresAt     time.Time
    RevokedAt     *time.Time
    AutoDeletedAt *time.Time
}
```

### 2. DSARRequest (Aggregate Root)
```go
type DSARRequest struct {
    ID              uuid.UUID
    UserID          uuid.UUID
    RequestType     DSARType        // VO: ACCESS, ERASURE, WITHDRAW_CONSENT
    Status          DSARStatus      // VO: PENDING, PROCESSING, COMPLETED, REJECTED
    RequestedAt     time.Time
    CompletedAt     *time.Time
    DataPayload     []byte
    RejectionReason string
    OTPCode         string
    OTPExpiredAt    time.Time
    IPAddress       string
    UserAgent       string
}
```

### 3. UserAccountStatus (Aggregate Root)
```go
type UserAccountStatus struct {
    ID                   uuid.UUID
    UserID               uuid.UUID
    Status               AccountStatus  // VO: ACTIVE, SUSPENDED, TERMINATED, DELETED
    SuspendedAt          *time.Time
    TerminatedAt         *time.Time
    DeletionConfirmedAt  *time.Time
    RetentionDeadline    *time.Time
    AutoDeletedAt        *time.Time
    UpdatedAt            time.Time
}
```

### 4. AuditTrail (Entity)
```go
type AuditTrail struct {
    ID        uuid.UUID
    UserID    *uuid.UUID
    Action    string
    Details   interface{}
    IPAddress string
    UserAgent string
    CreatedAt time.Time
}
```

### 5. PrivacyPolicy (Entity)
```go
type PrivacyPolicy struct {
    ID            uuid.UUID
    Version       string
    Title         string
    Content       string
    EffectiveDate time.Time
    IsActive      bool
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

---

## Value Objects

```go
// ConsentPurpose
type ConsentPurpose string
const (
    PurposeNecessary ConsentPurpose = "NECESSARY"
    PurposeAnalytics ConsentPurpose = "ANALYTICS"
    PurposeMarketing ConsentPurpose = "MARKETING"
)

// ConsentStatus
type ConsentStatus string
const (
    ConsentGranted ConsentStatus = "GRANTED"
    ConsentRevoked ConsentStatus = "REVOKED"
    ConsentExpired ConsentStatus = "EXPIRED"
    ConsentDeleted ConsentStatus = "DELETED"
)

// DSARType
type DSARType string
const (
    DSARTypeAccess          DSARType = "ACCESS"
    DSARTypeErasure         DSARType = "ERASURE"
    DSARTypeWithdrawConsent DSARType = "WITHDRAW_CONSENT"
)

// DSARStatus
type DSARStatus string
const (
    DSARStatusPending    DSARStatus = "PENDING"
    DSARStatusProcessing DSARStatus = "PROCESSING"
    DSARStatusCompleted  DSARStatus = "COMPLETED"
    DSARStatusRejected   DSARStatus = "REJECTED"
)

// AccountStatus
type AccountStatus string
const (
    AccountActive     AccountStatus = "ACTIVE"
    AccountSuspended  AccountStatus = "SUSPENDED"
    AccountTerminated AccountStatus = "TERMINATED"
    AccountDeleted    AccountStatus = "DELETED"
)
```

---

## Domain Events (Event-Driven Architecture)

### Event Definitions

```go
// domain/event/consent_events.go
type ConsentGrantedEvent struct {
    EventID   uuid.UUID
    UserID    uuid.UUID
    Purpose   string
    GrantedAt time.Time
    Metadata  EventMetadata
}

type ConsentRevokedEvent struct {
    EventID   uuid.UUID
    UserID    uuid.UUID
    Purpose   string
    RevokedAt time.Time
    Metadata  EventMetadata
}

type DSARSubmittedEvent struct {
    EventID     uuid.UUID
    DSARID      uuid.UUID
    UserID      uuid.UUID
    RequestType string
    RequestedAt time.Time
    Metadata    EventMetadata
}

type DSARCompletedEvent struct {
    EventID     uuid.UUID
    DSARID      uuid.UUID
    UserID      uuid.UUID
    RequestType string
    CompletedAt time.Time
    Metadata    EventMetadata
}

type AccountSuspendedEvent struct {
    EventID          uuid.UUID
    UserID           uuid.UUID
    SuspendedAt      time.Time
    RetentionYears   int
    Metadata         EventMetadata
}

type AccountTerminatedEvent struct {
    EventID      uuid.UUID
    UserID       uuid.UUID
    TerminatedAt time.Time
    Metadata     EventMetadata
}

type DataDeletionRequestedEvent struct {
    EventID   uuid.UUID
    UserID    uuid.UUID
    Reason    string
    RequestedAt time.Time
    Metadata  EventMetadata
}

type DataDeletedEvent struct {
    EventID   uuid.UUID
    UserID    uuid.UUID
    DeletedAt time.Time
    DeleteType string // IMMEDIATE, AUTO, DSAR
    Metadata  EventMetadata
}

type EventMetadata struct {
    CorrelationID string
    CausationID   string
    TraceID       string
    Timestamp     time.Time
    Version       int
}
```

### Kafka Topics

```
pdpa.consent.granted
pdpa.consent.revoked
pdpa.dsar.submitted
pdpa.dsar.completed
pdpa.account.suspended
pdpa.account.terminated
pdpa.data.deletion_requested
pdpa.data.deleted
pdpa.audit.trail
pdpa.llm.analysis.requested
pdpa.email.notification
pdpa.blockchain.record
```

---

## Use Cases

### Commands (Write Operations)

1. **RecordConsentCommand** - บันทึกความยินยอม
   - Input: UserID, SessionID, Purposes[], IP, UserAgent
   - Events: `ConsentGrantedEvent`

2. **RevokeConsentCommand** - ถอนความยินยอม
   - Input: UserID, Purpose, IP, UserAgent
   - Events: `ConsentRevokedEvent`

3. **SubmitDSARCommand** - ส่งคำร้องขอใช้สิทธิ์
   - Input: UserID, RequestType, IP, UserAgent
   - Events: `DSARSubmittedEvent`

4. **ProcessDSARCommand** - ประมวลผลคำร้อง
   - Input: DSARID, Action, Payload
   - Events: `DSARCompletedEvent`

5. **ConfirmDeletionCommand** - ยืนยันการลบข้อมูล
   - Input: UserID, AdminID
   - Events: `DataDeletionRequestedEvent`

6. **ImmediateDeletionCommand** - ลบข้อมูลทันที
   - Input: UserID, Reason
   - Events: `DataDeletedEvent`

7. **AutoDeleteExpiredCommand** - ลบข้อมูลอัตโนมัติ (Scheduled)
   - Input: RetentionYears
   - Events: `DataDeletedEvent`

8. **HandleAccountSuspendedCommand** - จัดการระงับบัญชี
   - Input: UserID, SuspendedAt, RetentionYears
   - Events: `AccountSuspendedEvent`

9. **HandleAccountTerminatedCommand** - จัดการยกเลิกบัญชี
   - Input: UserID, TerminatedAt
   - Events: `AccountTerminatedEvent`

10. **PublishPrivacyPolicyCommand** - เผยแพร่นโยบาย
    - Input: Version, Title, Content, EffectiveDate
    - Events: `PrivacyPolicyPublishedEvent`

### Queries (Read Operations)

1. **GetConsentHistoryQuery** - ดูประวัติความยินยอม
2. **GetDSARStatusQuery** - ดูสถานะคำร้อง
3. **GetAdminReportQuery** - รายงานสำหรับ Admin
4. **GetActivePrivacyPolicyQuery** - ดูนโยบายที่ใช้งาน
5. **GetUserConsentStatusQuery** - ดูสถานะความยินยอมปัจจุบัน

### Event Handlers (Async)

1. **ConsentGrantedHandler** - ส่ง email ยืนยัน, index Elasticsearch, cache Redis
2. **ConsentRevokedHandler** - อัปเดต cache, แจ้งเตือน
3. **DSARSubmittedHandler** - ส่ง OTP email, สร้าง LLM analysis job
4. **DataDeletionRequestedHandler** - ลบจากทุก service, บันทึก blockchain
5. **AccountSuspendedHandler** - ตั้งเวลา auto-delete, แจ้งเตือน
6. **LLMAnalysisHandler** - วิเคราะห์ข้อมูลด้วย LLM
7. **EmailNotificationHandler** - ส่ง email
8. **BlockchainRecordHandler** - บันทึก audit log ลง blockchain

---

## API Endpoints

### Consent Management
```
POST   /api/v1/pdpa/consent              # บันทึกความยินยอม
DELETE /api/v1/pdpa/consent              # ถอนความยินยอม
GET    /api/v1/pdpa/consent/history      # ดูประวัติ
GET    /api/v1/pdpa/consent/status       # ดูสถานะปัจจุบัน
```

### DSAR
```
POST   /api/v1/pdpa/dsar                 # ส่งคำร้อง
GET    /api/v1/pdpa/dsar/:id             # ดูสถานะ
GET    /api/v1/pdpa/dsar                 # ดูรายการทั้งหมด
POST   /api/v1/pdpa/dsar/:id/verify-otp  # ยืนยัน OTP
POST   /api/v1/pdpa/dsar/:id/process     # ประมวลผล (Admin)
```

### Account Management
```
POST   /api/v1/pdpa/account/suspend      # ระงับบัญชี
POST   /api/v1/pdpa/account/terminate    # ยกเลิกบัญชี
POST   /api/v1/pdpa/account/confirm-deletion  # ยืนยันการลบ
```

### Privacy Policy
```
GET    /api/v1/pdpa/policy               # ดูนโยบายปัจจุบัน
GET    /api/v1/pdpa/policy/versions      # ดูเวอร์ชันทั้งหมด
POST   /api/v1/pdpa/policy               # สร้างเวอร์ชันใหม่ (Admin)
PUT    /api/v1/pdpa/policy/:id/activate  # เปิดใช้งาน (Admin)
```

### Admin Reports
```
GET    /api/v1/pdpa/admin/reports        # รายงาน
GET    /api/v1/pdpa/admin/audit-trails   # Audit logs
GET    /api/v1/pdpa/admin/statistics     # สถิติ
```

### WebSocket
```
WS     /api/v1/pdpa/ws/dsar-status       # Real-time DSAR status
```

---

## Constraints

1. **Domain Layer ห้าม import:** gorm, gin, chi, sarama, redis, elasticsearch
2. **Entity ต้องมี constructor** `New{Entity}` + validation
3. **เปลี่ยน state ผ่าน behavior methods เท่านั้น** (ไม่มี setter)
4. **Repository เป็น interface เท่านั้น**
5. **Errors เป็น sentinel errors**
6. **Comment 2 ภาษา** (ไทย/English)
7. **Test coverage ≥ 80%**
8. **Table prefix:** `pdpa_`
9. **Redis key pattern:** `pdpa:{entity}:{id}`
10. **Kafka topic pattern:** `pdpa.{entity}.{action}`
11. **ใช้ Transactional Outbox Pattern** สำหรับ event publishing
12. **ใช้ Correlation ID** ในการ trace events
13. **Implement Idempotency** สำหรับ event handlers
14. **ใช้ Circuit Breaker** สำหรับ external services (LLM, Email)
15. **ทุก Command ต้อง publish domain event**

---

# FORMAT

ตอบในรูปแบบ:

## 1. Architecture Diagram
```
[แสดงโครงสร้าง Event-Driven Architecture]
[แสดง data flow ระหว่าง components]
[แสดง Kafka topics และ consumers]
```

## 2. File Tree
```
internal/modules/pdpa/
├── domain/
│   ├── entity/
│   │   ├── consent_log.go
│   │   ├── dsar_request.go
│   │   ├── user_account_status.go
│   │   ├── audit_trail.go
│   │   └── privacy_policy.go
│   ├── value_object/
│   │   ├── consent_purpose.go
│   │   ├── consent_status.go
│   │   ├── dsar_type.go
│   │   ├── dsar_status.go
│   │   └── account_status.go
│   ├── event/
│   │   ├── consent_events.go
│   │   ├── dsar_events.go
│   │   ├── account_events.go
│   │   └── event_metadata.go
│   ├── repository/
│   │   ├── consent_repository.go
│   │   ├── dsar_repository.go
│   │   ├── account_status_repository.go
│   │   ├── audit_repository.go
│   │   └── policy_repository.go
│   ├── service/
│   │   ├── consent_validator.go
│   │   ├── deletion_policy_service.go
│   │   └── anonymization_service.go
│   └── errors/
│       └── errors.go
├── application/
│   ├── command/
│   │   ├── record_consent.go
│   │   ├── revoke_consent.go
│   │   ├── submit_dsar.go
│   │   ├── process_dsar.go
│   │   ├── confirm_deletion.go
│   │   ├── immediate_deletion.go
│   │   ├── auto_delete_expired.go
│   │   ├── handle_account_suspended.go
│   │   ├── handle_account_terminated.go
│   │   └── publish_privacy_policy.go
│   ├── query/
│   │   ├── get_consent_history.go
│   │   ├── get_dsar_status.go
│   │   ├── get_admin_report.go
│   │   └── get_active_privacy_policy.go
│   ├── event_handler/
│   │   ├── consent_granted_handler.go
│   │   ├── consent_revoked_handler.go
│   │   ├── dsar_submitted_handler.go
│   │   ├── data_deletion_handler.go
│   │   ├── account_suspended_handler.go
│   │   ├── llm_analysis_handler.go
│   │   ├── email_notification_handler.go
│   │   └── blockchain_record_handler.go
│   ├── dto/
│   │   └── dto.go
│   └── mocks/
│       └── repository_mocks.go
├── infrastructure/
│   ├── persistence/
│   │   ├── postgres/
│   │   │   ├── models.go
│   │   │   ├── consent_repo_impl.go
│   │   │   ├── dsar_repo_impl.go
│   │   │   ├── account_status_repo_impl.go
│   │   │   ├── audit_repo_impl.go
│   │   │   ├── policy_repo_impl.go
│   │   │   └── outbox_repo_impl.go
│   │   └── redis/
│   │       └── consent_cache.go
│   ├── messaging/
│   │   ├── kafka/
│   │   │   ├── producer.go
│   │   │   ├── event_serializer.go
│   │   │   └── consumer/
│   │   │       ├── dsar_consumer.go
│   │   │       ├── deletion_consumer.go
│   │   │       ├── llm_consumer.go
│   │   │       ├── email_consumer.go
│   │   │       └── blockchain_consumer.go
│   │   └── outbox/
│   │       └── outbox_publisher.go
│   ├── search/
│   │   └── elasticsearch/
│   │       └── consent_indexer.go
│   ├── scheduler/
│   │   └── consent_cleanup_job.go
│   └── external/
│       ├── llm_client.go
│       ├── email_client.go
│       └── blockchain_client.go
├── interfaces/
│   ├── http/
│   │   ├── consent_handler.go
│   │   ├── dsar_handler.go
│   │   ├── account_handler.go
│   │   ├── policy_handler.go
│   │   ├── admin_handler.go
│   │   ├── dto.go
│   │   └── routes.go
│   ├── websocket/
│   │   └── dsar_broadcaster.go
│   └── localization/
│       ├── locales/
│       │   ├── en.toml
│       │   └── th.toml
│       └── i18n.go
└── module.go
```

## 3. Code Implementation
[สร้างโค้ดจริงสำหรับทุกไฟล์]

### 3.1 Domain Layer
- Entities พร้อม behavior methods
- Value Objects พร้อม validation
- Domain Events พร้อม metadata
- Repository Interfaces
- Domain Services
- Sentinel Errors

### 3.2 Application Layer
- Commands (Write Use Cases) พร้อม event publishing
- Queries (Read Use Cases)
- Event Handlers พร้อม idempotency
- DTOs

### 3.3 Infrastructure Layer
- PostgreSQL Repository Implementations
- Outbox Pattern Implementation
- Kafka Producer/Consumers
- Redis Cache
- Elasticsearch Indexer
- Scheduler
- External Service Clients

### 3.4 Interface Layer
- HTTP Handlers
- Routes
- WebSocket Broadcaster
- Localization (i18n)

### 3.5 Module Entry Point
- module.go พร้อม dependency injection

## 4. Database Migrations
```sql
-- migrations/001_initial_pdpa_schema.sql
[สร้าง schema ทั้งหมดพร้อม prefix pdpa_]
[รวม outbox table สำหรับ Transactional Outbox Pattern]
```

## 5. Docker Compose
```yaml
[สร้าง docker-compose.yml สำหรับ development]
```

## 6. Testing Strategy
```
[อธิบายกลยุทธ์การทดสอบ]
- Unit Tests สำหรับ Domain Layer
- Integration Tests สำหรับ Repository
- Event Handler Tests พร้อม idempotency
- End-to-End Tests
```

## 7. Deployment Considerations
```
[อธิบายการ deploy]
- Horizontal Scaling
- Kafka Consumer Groups
- Database Connection Pooling
- Graceful Shutdown
- Health Checks
```

---

## Additional Requirements

### Transactional Outbox Pattern
```go
// ทุก Command ต้องบันทึก event ลง outbox table ใน transaction เดียวกัน
// แล้วมี background process คอย publish ไป Kafka
type OutboxEvent struct {
    ID            uuid.UUID
    AggregateType string
    AggregateID   uuid.UUID
    EventType     string
    Payload       []byte
    CreatedAt     time.Time
    PublishedAt   *time.Time
    Status        string // PENDING, PUBLISHED, FAILED
}
```

### Idempotency
```go
// Event handlers ต้องตรวจสอบว่าประมวลผล event นี้ไปแล้วหรือยัง
type ProcessedEvent struct {
    EventID     uuid.UUID
    HandlerName string
    ProcessedAt time.Time
}
```

### Correlation & Tracing
```go
// ทุก event ต้องมี metadata สำหรับ tracing
type EventMetadata struct {
    CorrelationID string  // trace ทั้ง flow
    CausationID   string  // event ที่ทำให้เกิด event นี้
    TraceID       string  // distributed tracing
    Timestamp     time.Time
    Version       int
}
```

### Circuit Breaker
```go
// External services ต้องมี circuit breaker
type CircuitBreakerConfig struct {
    MaxRequests   uint32
    Interval      time.Duration
    Timeout       time.Duration
    ReadyToTrip   func(counts Counts) bool
    OnStateChange func(name string, from State, to State)
}
```

### Retry Policy
```go
// Event handlers ต้องมี retry policy
type RetryPolicy struct {
    MaxRetries      int
    InitialInterval time.Duration
    MaxInterval     time.Duration
    Multiplier      float64
}
```

### Dead Letter Queue
```go
// Events ที่ process ไม่สำเร็จต้องเข้า DLQ
// Kafka topic: pdpa.dlq.{original_topic}
```

---

## Output Requirements

1. **ครบถ้วน**: สร้างทุกไฟล์ที่ระบุใน file tree
2. **Production-Ready**: โค้ดต้องพร้อมใช้งานจริง
3. **Well-Documented**: Comment 2 ภาษา (ไทย/English)
4. **Testable**: มี interfaces สำหรับ mocking
5. **Event-Driven**: ทุก state change ต้อง publish event
6. **Idempotent**: Event handlers ต้อง idempotent
7. **Observable**: มี logging, metrics, tracing
8. **Resilient**: มี circuit breaker, retry, DLQ
9. **Scalable**: รองรับ horizontal scaling
10. **Secure**: มี input validation, encryption, RBAC