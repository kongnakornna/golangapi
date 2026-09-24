# 🤖 AI Prompt สำหรับสร้าง PDPA Module + Writing Plans

## 📋 Prompt Template (คัดลอกไปใช้ได้เลย)

```markdown
# คำสั่ง: สร้างระบบ PDPA Module ตาม Clean Architecture + DDD

## 🎯 วัตถุประสงค์
สร้างโมดูลระบบจัดการสิทธิ์ส่วนบุคคล (PDPA Module) สำหรับภาษา Go 
โดยใช้สถาปัตยกรรม Clean Architecture + Domain-Driven Design (DDD)

---

## 📐 สถาปัตยกรรมที่ต้องการ

### Layered Architecture (4 Layers)
1. **Domain Layer** - Entities, Value Objects, Repository Interfaces, Domain Services, Domain Errors
2. **Application Layer** - Use Cases (Business Logic)
3. **Infrastructure Layer** - PostgreSQL (GORM), Redis, Kafka, Elasticsearch, Scheduler
4. **Interface Layer** - HTTP Handlers (Gin), Routes, WebSocket, Localization (i18n)

### ข้อกำหนดสำคัญ
- ✅ ใช้ **Prefix `pdpa_`** ในทุกตารางฐานข้อมูล
- ✅ แยก **Domain Entities** ออกจาก **Persistence Models** อย่างชัดเจน
- ✅ ใช้ **Value Objects** สำหรับ Type Safety
- ✅ ใช้ **Repository Pattern** สำหรับ Data Access
- ✅ ใช้ **Use Case Pattern** สำหรับ Business Logic
- ✅ ใช้ **Dependency Injection** ผ่าน Constructor

---

## 🗄️ Database Schema ที่ต้องการ

### ตารางทั้งหมด (Prefix: pdpa_)
| ตาราง | คำอธิบาย |
|-------|----------|
| `pdpa_policies` | เวอร์ชันนโยบายความเป็นส่วนตัว |
| `pdpa_purposes` | วัตถุประสงค์ในการขอข้อมูล (NECESSARY, ANALYTICS, MARKETING) |
| `pdpa_consents` | บันทึกการให้ความยินยอม |
| `pdpa_user_requests` | คำร้องขอใช้สิทธิ์ (DSAR) |
| `pdpa_user_account_statuses` | สถานะบัญชีผู้ใช้ (ACTIVE, SUSPENDED, TERMINATED, DELETED) |
| `pdpa_audit_trails` | บันทึกการดำเนินการทั้งหมด |
| `pdpa_request_responses` | ไฟล์ตอบกลับคำร้อง |

### ความสัมพันธ์
- `pdpa_consents.purpose_code` → `pdpa_purposes.code` (FK)
- `pdpa_request_responses.request_id` → `pdpa_user_requests.id` (FK, CASCADE)

---

## 📁 โครงสร้างโปรเจกต์ที่ต้องการ

```
icmongolang/
├── cmd/
│   ├── api/main.go              # REST API Server
│   ├── scheduler/main.go        # Cron Job สำหรับ Auto-deletion
│   ├── migrate/main.go          # Database Migration
│   └── workers/
│       ├── dsar/main.go         # DSAR Processing Worker
│       ├── email/main.go        # Email Sending Worker
│       ├── llm/main.go          # LLM Analysis Worker
│       ├── revoke/main.go       # Revoke Consent Worker
│       └── account/main.go      # Account Event Worker
├── internal/modules/pdpa/
│   ├── domain/
│   │   ├── entity/              # ConsentLog, DSARRequest, AuditTrail, UserAccountStatus
│   │   ├── value_object/        # ConsentPurpose, ConsentStatus, DSARType, DSARStatus, AccountStatus
│   │   ├── repository/          # Repository Interfaces
│   │   ├── service/             # DeletionPolicyService, ConsentValidator, BlockchainService
│   │   ├── event/               # Domain Events
│   │   └── errors/              # Domain Errors
│   ├── application/             # Use Cases ทั้งหมด
│   ├── infrastructure/
│   │   ├── persistence/
│   │   │   ├── postgres/        # GORM Models + Repository Implementations
│   │   │   └── redis/           # Consent Cache
│   │   ├── messaging/
│   │   │   ├── kafka_producer.go
│   │   │   └── consumers/       # Kafka Consumer Workers
│   │   ├── search/elasticsearch/# Consent Indexer
│   │   ├── services/
│   │   │   ├── email/           # SMTP Sender
│   │   │   ├── llm/             # OpenAI Client
│   │   │   └── blockchain/      # Ethereum Client
│   │   └── scheduler/           # Cleanup Job
│   └── interfaces/
│       ├── http/                # Handlers + Routes
│       ├── websocket/           # Real-time Hub
│       └── localization/        # i18n (EN/TH)
├── migrations/                  # SQL Migrations
├── go.mod
├── .env
└── docker-compose.yml
```

---

## 🔧 Use Cases ที่ต้องสร้าง

### Consent Management
1. **RecordConsentUseCase** - บันทึกการให้ความยินยอม
2. **RevokeConsentUseCase** - เพิกถอนความยินยอม
3. **GetConsentHistoryUseCase** - ดึงประวัติความยินยอม

### DSAR (Data Subject Access Request)
4. **SubmitDSARUseCase** - ยื่นคำร้องขอใช้สิทธิ์ (สร้าง OTP)
5. **ProcessDSARUseCase** - ประมวลผลคำร้อง (Approve/Reject)
6. **GetDSARStatusUseCase** - ตรวจสอบสถานะคำร้อง

### Deletion Management
7. **ImmediateDeletionUseCase** - ลบข้อมูลทันที (Terminated + Confirmed)
8. **AutoDeleteExpiredConsentsUseCase** - ลบข้อมูลอัตโนมัติ (Suspended + หมดอายุ)

### Account Events
9. **HandleAccountEventUseCase** - จัดการ Event: Suspended, Terminated, ConfirmDeletion

### Reporting
10. **GetAdminReportUseCase** - รายงานสำหรับ Admin/DPO

---

## 🔌 Integration ที่ต้องรองรับ

### Message Queue (Kafka Topics)
| Topic | Producer | Consumer |
|-------|----------|----------|
| `pdpa.consent.log` | API | Audit Logger |
| `pdpa.consent.revoked` | API | Revoke Worker |
| `pdpa.dsar.request` | API | DSAR Worker |
| `pdpa.data.deleted` | API/Scheduler | Notification |
| `pdpa.account.event` | External | Account Worker |
| `pdpa.email.send` | Various | Email Worker |
| `pdpa.llm.analyze` | DSAR Worker | LLM Worker |

### External Services
- **PostgreSQL** - Primary Database
- **Redis** - Consent Cache (TTL: 30 days)
- **Kafka** - Async Messaging
- **Elasticsearch** - Consent Search Index
- **Ethereum** - Immutable Audit Trail (Blockchain)
- **SMTP** - Email Notifications
- **OpenAI** - LLM Analysis

---

## 🎯 Business Rules ที่ต้อง implement

### Consent Rules
```
- Consent หมดอายุใน 1 ปี (จากวันที่ให้)
- ไม่สามารถ Revoke consent ที่ไม่ได้อยู่ในสถานะ GRANTED
- บัญชีที่ SUSPENDED/TERMINATED ไม่สามารถให้ consent ใหม่ได้
- Purpose ต้องเป็น: NECESSARY, ANALYTICS, หรือ MARKETING เท่านั้น
```

### Account Status Rules
```
- SUSPENDED: ข้อมูลเก็บไว้ 1 ปี (Retention Period)
- TERMINATED: รอการยืนยันการลบ
- DELETED: ข้อมูลถูกลบ/Anonymize แล้ว
- Immediate Deletion: ต้องเป็น TERMINATED + มี DeletionConfirmedAt
- Auto Deletion: ต้องเป็น SUSPENDED + หมด RetentionDeadline
```

### DSAR Rules
```
- DSAR Types: ACCESS, ERASURE, WITHDRAW_CONSENT
- Status Flow: PENDING → PROCESSING → COMPLETED/REJECTED
- ต้องมี OTP ยืนยัน (หมดอายุใน 15 นาที)
- Rate Limiting: จำกัดจำนวนคำร้องต่อวัน
```

### Data Deletion Rules
```
- Consent/DSAR: Hard Delete (ลบออกจาก DB)
- User Data: Anonymization (email → deleted_{uuid}_{timestamp})
- Cache: ลบ Redis keys ทั้งหมดของ user
- Search Index: ลบ Elasticsearch documents
- Blockchain: บันทึก Proof of Deletion (ไม่ลบ)
```

---

## 📝 รูปแบบโค้ดที่ต้องการ

### Domain Entity Example
```go
type ConsentLog struct {
    ID            uuid.UUID
    UserID        uuid.UUID
    SessionID     string
    Purpose       valueobject.ConsentPurpose
    Status        valueobject.ConsentStatus
    IPAddress     string
    UserAgent     string
    GrantedAt     time.Time
    ExpiresAt     time.Time
    RevokedAt     *time.Time
    AutoDeletedAt *time.Time
}

// Domain Methods
func (c *ConsentLog) Revoke() error
func (c *ConsentLog) MarkDeleted()
func (c *ConsentLog) IsActive() bool
```

### Use Case Example
```go
type RecordConsentUseCase struct {
    consentRepo       repository.ConsentRepository
    accountStatusRepo repository.UserAccountStatusRepository
    auditRepo         repository.AuditRepository
    cache             redis.ConsentCache
    producer          messaging.KafkaProducer
    blockchainSvc     service.BlockchainService
}

func (uc *RecordConsentUseCase) Execute(ctx context.Context, input RecordConsentInput) error {
    // 1. Validate account status
    // 2. Create consent log
    // 3. Save to DB
    // 4. Update cache
    // 5. Publish event
    // 6. Record to blockchain
    // 7. Create audit trail
}
```

### Repository Implementation Example
```go
type consentRepoImpl struct {
    db *gorm.DB
}

func (r *consentRepoImpl) Save(ctx context.Context, log *entity.ConsentLog) error {
    model := &PdpAConsentModel{
        ID:          log.ID,
        UserID:      log.UserID,
        PurposeCode: log.Purpose.String(),
        Status:      string(log.Status),
        // ...
    }
    return r.db.WithContext(ctx).Create(model).Error
}
```

---

## ✅ Deliverables ที่ต้องส่งมอบ

### 1. Domain Layer
- [ ] Entities: ConsentLog, DSARRequest, AuditTrail, UserAccountStatus
- [ ] Value Objects: ConsentPurpose, ConsentStatus, DSARType, DSARStatus, AccountStatus
- [ ] Repository Interfaces: ConsentRepository, DSARRepository, AuditRepository, UserAccountStatusRepository
- [ ] Domain Services: DeletionPolicyService, ConsentValidator, BlockchainService
- [ ] Domain Errors: ครบทุก error case

### 2. Application Layer
- [ ] Use Cases ทั้ง 10 รายการ
- [ ] DTOs สำหรับ Input/Output
- [ ] Event Publishing Logic

### 3. Infrastructure Layer
- [ ] GORM Models (Prefix: pdpa_)
- [ ] Repository Implementations
- [ ] Redis Cache Implementation
- [ ] Kafka Producer + Consumers
- [ ] Elasticsearch Indexer
- [ ] Scheduler Job
- [ ] Email/LLM/Blockchain Services

### 4. Interface Layer
- [ ] HTTP Handlers: Consent, DSAR, Deletion, Admin
- [ ] Routes Registration
- [ ] WebSocket Hub + Handler
- [ ] i18n (EN/TH)

### 5. Entry Points
- [ ] cmd/api/main.go
- [ ] cmd/scheduler/main.go
- [ ] cmd/workers/*/main.go (5 workers)
- [ ] cmd/migrate/main.go

### 6. Configuration
- [ ] go.mod (พร้อม dependencies)
- [ ] .env (ตัวอย่าง)
- [ ] docker-compose.yml (PostgreSQL, Redis, Kafka, Elasticsearch)
- [ ] migrations/001_initial_pdpa_schema.sql

---

## 📊 Writing Plans ที่ต้องสร้าง

1. สร้างโครงสร้างโปรเจกต์ + go.mod
2. สร้าง Database Migration (SQL)
3. สร้าง Domain Entities + Value Objects
4. สร้าง Domain Errors
5. สร้าง Repository Interfaces
6. สร้าง Domain Services (DeletionPolicyService, etc.)
7. เขียน Unit Tests สำหรับ Domain
8. สร้าง Use Cases ทั้งหมด
9. สร้าง DTOs
10. เขียน Unit Tests สำหรับ Use Cases
11. สร้าง GORM Models
12. สร้าง Repository Implementations
13. สร้าง Redis Cache
14. สร้าง Kafka Producer/Consumers
15. สร้าง Elasticsearch Indexer
16. สร้าง Email/LLM/Blockchain Services
17. สร้าง HTTP Handlers
18. สร้าง Routes
19. สร้าง WebSocket Hub
20. สร้าง i18n
21. สร้าง cmd/api/main.go
22. สร้าง cmd/scheduler/main.go
23. สร้าง cmd/workers/*/main.go
24. Integration Tests
25. Docker Compose Setup
26. Documentation

---

## 🚨 ข้อกำหนดเพิ่มเติม
1. **Error Handling**: ใช้ Domain Errors ที่ชัดเจน, ไม่ leak infrastructure errors
2. **Context Propagation**: ส่ง `context.Context` ไปทุก layer
3. **Logging**: ใช้ structured logging
4. **Graceful Shutdown**: รองรับ SIGINT/SIGTERM
5. **Configuration**: ใช้ environment variables ทั้งหมด
6. **Security**: 
   - Hash OTP codes
   - Validate JWT
   - Rate limiting สำหรับ DSAR
7. **Performance**:
   - ใช้ Index ในทุก query ที่สำคัญ
   - Cache consent status ใน Redis
   - Async processing ผ่าน Kafka

---

## 📤 Output Format ที่ต้องการ

กรุณาสร้างโค้ดตามลำดับ:
1. **go.mod** - dependencies ทั้งหมด
2. **migrations/*.sql** - database schema
3. **Domain Layer** - entities, value objects, repositories, services, errors
4. **Application Layer** - use cases ทั้งหมด
5. **Infrastructure Layer** - models, repos, cache, messaging, services
6. **Interface Layer** - handlers, routes, websocket, i18n
7. **Entry Points** - main.go files
8. **Configuration** - .env, docker-compose.yml

พร้อมเขียน comment อธิบายในส่วนที่ซับซ้อน

---


เริ่มสร้างได้เลย!
```

---

## 🎯 วิธีใช้ Prompt นี้

### สำหรับ AI แต่ละตัว:

| AI | ข้อแนะนำ |
|-----|----------|
| **Claude** | ใช้ prompt นี้ได้เลย เพราะ context window ใหญ่ |
| **ChatGPT** | แบ่งเป็น 3-4 ส่วน (Domain, Application, Infrastructure, Interface) |
| **Gemini** | ใช้ prompt นี้ + ระบุ "สร้างทีละ layer" |

### การแบ่ง Phase สำหรับ ChatGPT:

```
Prompt 1: "สร้าง Domain Layer ตาม spec นี้: [paste domain section]"
Prompt 2: "สร้าง Application Layer ตาม spec นี้: [paste application section]"
Prompt 3: "สร้าง Infrastructure Layer ตาม spec นี้: [paste infrastructure section]"
Prompt 4: "สร้าง Interface Layer + Entry Points ตาม spec นี้: [paste interface section]"
```

### Tips เพิ่มเติม:

1. **ระบุ AI Model**: "ใช้ Go 1.21, GORM v2, Gin framework"
2. **ระบุ Style**: "ใช้ Clean Code, มี comment ภาษาไทย/อังกฤษ"
3. **ระบุ Testing**: "เขียน unit test สำหรับ use case ด้วย testify"
4. **ขอ Diagrams**: "สร้าง Mermaid diagram สำหรับ flow การทำงาน"

---

## 📁 ไฟล์ที่ได้จาก Prompt นี้

```
✅ ~50+ ไฟล์ Go
✅ 1 ไฟล์ SQL Migration
✅ go.mod + .env + docker-compose.yml
✅ Writing Plans 7 phases
✅ Unit Tests (optional)
```