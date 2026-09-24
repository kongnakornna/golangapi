### 📗 คู่มือการดูแล แก้ไข และขยายระบบ (Operations & Maintenance Manual)

> **คู่มือนี้ใช้คู่กับ `Template_Module.md`**  
> `Template_Module.md` = "พิมพ์เขียว" (สร้างใหม่)  
> `คู่มือนี้` = "คู่มือช่าง" (ดูแล, ซ่อม, ต่อเติม)  
>
> **กลุ่มเป้าหมาย:** Developer, Tech Lead, DevOps, DBA ที่ดูแลระบบ `icmongolang`

---

## สารบัญคู่มือ

**ภาค 1 — การดูแลรักษา (Maintenance)**
1. [ปรัชญาการดูแลระบบ](#1-ปรัชญาการดูแลระบบ)
2. [Health Check & Monitoring](#2-health-check--monitoring)
3. [Version & Release Management](#3-version--release-management)
4. [Backup / Restore / Disaster Recovery](#4-backup--restore--disaster-recovery)
5. [Dependency Upgrade Playbook](#5-dependency-upgrade-playbook)
6. [Secret & Config Management](#6-secret--config-management)

**ภาค 2 — การแก้ไข (Edit Guide)**
7. [Workflow การแก้ไขที่ปลอดภัย](#7-workflow-การแก้ไขที่ปลอดภัย)
8. [แก้ Domain Layer](#8-แก้-domain-layer)
9. [แก้ Application Layer](#9-แก้-application-layer)
10. [แก้ Infrastructure Layer](#10-แก้-infrastructure-layer)
11. [แก้ Interface Layer](#11-แก้-interface-layer)
12. [แก้ Database Migration](#12-แก้-database-migration)
13. [Breaking Change Playbook](#13-breaking-change-playbook)

**ภาค 3 — การขยายฟีเจอร์ (Extension Guide)**
14. [เพิ่ม Use Case ใหม่](#14-เพิ่ม-use-case-ใหม่)
15. [เพิ่ม Entity / Aggregate ใหม่](#15-เพิ่ม-entity--aggregate-ใหม่)
16. [เพิ่ม Value Object ใหม่](#16-เพิ่ม-value-object-ใหม่)
17. [เพิ่ม Kafka Topic / Consumer](#17-เพิ่ม-kafka-topic--consumer)
18. [เพิ่ม Worker ใหม่](#18-เพิ่ม-worker-ใหม่)
19. [เพิ่ม Scheduler Job](#19-เพิ่ม-scheduler-job)
20. [เพิ่ม HTTP Endpoint](#20-เพิ่ม-http-endpoint)
21. [เพิ่ม WebSocket Event](#21-เพิ่ม-websocket-event)
22. [เพิ่ม Outbound Port ใหม่](#22-เพิ่ม-outbound-port-ใหม่)

**ภาค 4 — การเสริมหลัง (Advanced Extensions)**
23. [Observability Stack](#23-observability-stack)
24. [Resilience Patterns](#24-resilience-patterns)
25. [Feature Flags & A/B Testing](#25-feature-flags--ab-testing)
26. [Multi-tenancy](#26-multi-tenancy)
27. [Advanced Caching Strategy](#27-advanced-caching-strategy)
28. [Testing Strategy](#28-testing-strategy)
29. [CI/CD Pipeline](#29-cicd-pipeline)
30. [Performance Tuning](#30-performance-tuning)
31. [Security Hardening](#31-security-hardening)

**ภาคผนวก**
- [A. Decision Tree](#a-decision-tree--ฉันควรทำอะไรต่อ)
- [B. Runbook ฉุกเฉิน](#b-runbook-ฉุกเฉิน)
- [C. Glossary](#c-glossary)

---

# ภาค 1 — การดูแลรักษา (Maintenance)

## 1. ปรัชญาการดูแลระบบ

### 1.1 หลัก 5 ข้อ

| หลัก | ความหมาย | ตัวอย่างการปฏิบัติ |
| :--- | :--- | :--- |
| **Single Source of Truth** | ข้อมูลจริงมีที่เดียว | Domain entity เป็นเจ้าของ state, DB เป็น persistence |
| **Fail Fast, Fail Loud** | ผิดต้องรู้ทันที | validate ที่ domain, log ที่ infrastructure |
| **Idempotent by Design** | ทำซ้ำได้ผลเท่าเดิม | Kafka consumer ใช้ `MarkMessage`, use case ใช้ idempotency key |
| **Observable by Default** | ตรวจสอบได้ทุกจุด | ทุก use case มี log + metric + trace |
| **Reversible** | ย้อนกลับได้เสมอ | migration มี `down`, feature flag ทุกฟีเจอร์ใหม่ |

### 1.2 ระดับความรุนแรงของปัญหา

```
P0 ─ ระบบล่ม / ข้อมูลรั่ว / สูญหาย     → แก้ทันที ≤ 15 นาที
P1 ─ ฟีเจอร์หลักใช้ไม่ได้              → แก้ ≤ 4 ชั่วโมง
P2 ─ ฟีเจอร์รองเสีย / performance ตก   → แก้ ≤ 2 วัน
P3 ─ Bug เล็ก / cosmetic               → แก้ใน sprint ถัดไป
P4 ─ ปรับปรุง / refactor               → backlog
```

---

## 2. Health Check & Monitoring

### 2.1 Health Check Endpoints (บังคับมีทุก service)

**`interfaces/http/health_handler.go`**
```go
package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	db        interface{ Ping() error }
	redis     interface{ Ping(ctx context.Context) error }
	producer  interface{ Healthy() bool }
	startTime time.Time
}

func (h *HealthHandler) Liveness(c *gin.Context) {
	// ระบบยังหายใจ → return 200 เสมอ (ถ้า process ยังรัน)
	c.JSON(http.StatusOK, gin.H{"status": "alive", "uptime": time.Since(h.startTime).String()})
}

func (h *HealthHandler) Readiness(c *gin.Context) {
	// ระบบพร้อมรับ traffic → เช็ค dependency ทั้งหมด
	checks := map[string]string{}
	overall := http.StatusOK

	if err := h.db.Ping(); err != nil {
		checks["database"] = "fail: " + err.Error()
		overall = http.StatusServiceUnavailable
	} else {
		checks["database"] = "ok"
	}

	if !h.producer.Healthy() {
		checks["kafka"] = "fail"
		overall = http.StatusServiceUnavailable
	} else {
		checks["kafka"] = "ok"
	}

	c.JSON(overall, gin.H{"status": overall == 200, "checks": checks})
}
```

**Routes:**
```go
r.GET("/healthz", h.Liveness)   // liveness (K8s)
r.GET("/readyz",  h.Readiness)  // readiness (K8s)
r.GET("/metrics", prometheusHandler())  // Prometheus
```

### 2.2 Metrics ที่ต้องมีทุก Use Case

```go
var (
	ucDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "usecase_duration_seconds",
			Help:    "Use case execution time",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"module", "usecase", "status"},
	)
	ucErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "usecase_errors_total",
			Help: "Total use case errors",
		},
		[]string{"module", "usecase", "error_type"},
	)
)

// Wrapper ที่ใช้กับทุก use case
func Instrument(module, usecase string, fn func() error) error {
	timer := prometheus.NewTimer(ucDuration.WithLabelValues(module, usecase, "start"))
	err := fn()
	if err != nil {
		ucErrors.WithLabelValues(module, usecase, classify(err)).Inc()
		ucDuration.WithLabelValues(module, usecase, "error").Observe(timer.ObserveDuration().Seconds())
	} else {
		ucDuration.WithLabelValues(module, usecase, "success").Observe(timer.ObserveDuration().Seconds())
	}
	return err
}
```

### 2.3 Logging Standard (Structured JSON)

```go
// ใช้ pkg/logger เท่านั้น
logger.Info("consent recorded",
	logger.String("module", "pdpa"),
	logger.String("usecase", "RecordConsent"),
	logger.String("user_id", uid.String()),
	logger.String("purpose", "MARKETING"),
	logger.Duration("duration", elapsed),
)

// ❌ ห้ามใช้ log.Println, fmt.Println ใน production code
```

**Required fields ในทุก log:**
- `timestamp`, `level`, `service`, `module`, `usecase`
- `trace_id`, `span_id` (ถ้ามี tracing)
- `user_id` (ถ้ามี)
- `error` (ถ้าเป็น error log)

### 2.4 Alert Rules (Prometheus)

| Alert | Condition | Severity |
| :--- | :--- | :--- |
| `HighErrorRate` | `rate(usecase_errors_total[5m]) > 10` | P1 |
| `SlowUseCase` | `histogram_quantile(0.95, uc_duration) > 2s` | P2 |
| `KafkaConsumerLag` | `kafka_consumergroup_lag > 1000` | P2 |
| `DBConnectionPoolExhausted` | `db_pool_wait_count > 0` | P1 |
| `ServiceDown` | `up == 0` for 1m | P0 |

---

## 3. Version & Release Management

### 3.1 Semantic Versioning

```
v<MAJOR>.<MINOR>.<PATCH>[-<prerelease>]
```

| เปลี่ยน | ตัวอย่าง | ต้องทำ |
| :--- | :--- | :--- |
| MAJOR | `v1.0.0 → v2.0.0` | Breaking change, migration guide |
| MINOR | `v1.0.0 → v1.1.0` | เพิ่มฟีเจอร์, backward compatible |
| PATCH | `v1.0.0 → v1.0.1` | Bug fix เท่านั้น |

### 3.2 Git Branching Model

```
main ────●────●────●────●────●───→  production
          \    \    \    \    \
           \    \    \    \    └─ hotfix/P0-consent-race
            \    \    \    └────── release/v1.2.0
             \    \    └─────────── feature/pdpa-dsar-v2
              \    └─────────────── feature/consent-analytics
               └─────────────────── develop
```

**Commit Convention:**
```
<type>(<scope>): <subject>

feat(pdpa): add auto-delete for expired consents
fix(users): correct JWT expiry check
docs(template): update extension guide
refactor(domain): extract ValueObject validation
chore(deps): bump gin to 1.10.0
```

### 3.3 Release Checklist

```markdown
## Pre-release
- [ ] ทุก test ผ่าน (`go test ./...`)
- [ ] migration ใหม่มี down script
- [ ] CHANGELOG.md อัปเดต
- [ ] version bump ใน main.go
- [ ] config ใหม่เพิ่มใน .env.example
- [ ] API docs (swagger) regenerate
- [ ] Load test ผ่าน (ถ้ามี major change)

## Release
- [ ] Tag: `git tag -a v1.2.0 -m "Release v1.2.0"`
- [ ] Push: `git push origin v1.2.0`
- [ ] CI สร้าง docker image
- [ ] Deploy staging → smoke test → production
- [ ] Monitor 30 นาทีหลัง deploy

## Post-release
- [ ] Merge main → develop
- [ ] ปิด issues
- [ ] อัปเดต runbook
```

---

## 4. Backup / Restore / Disaster Recovery

### 4.1 Backup Matrix

| ข้อมูล | ความถี่ | เก็บที่ | Retention | RPO | RTO |
| :--- | :--- | :--- | :--- | :--- | :--- |
| PostgreSQL | ทุก 6 ชม. (full) + WAL (continuous) | S3 | 30 วัน | 5 นาที | 1 ชม. |
| Redis | ทุก 1 ชม. (RDB) | S3 | 7 วัน | 1 ชม. | 15 นาที |
| Elasticsearch | Snapshot ทุกวัน | S3 | 14 วัน | 1 วัน | 2 ชม. |
| Kafka | Retention 7 วัน | Kafka log | 7 วัน | 0 | – |
| Config/Secret | ทุกครั้งที่เปลี่ยน | Vault + Git | ∞ | 0 | 5 นาที |

### 4.2 PostgreSQL Backup Script

**`scripts/backup-pg.sh`**
```bash
#!/usr/bin/env bash
set -euo pipefail

DB_NAME="${DB_NAME:-pdpa_db}"
BACKUP_DIR="${BACKUP_DIR:-/backup/postgres}"
S3_BUCKET="${S3_BUCKET:-s3://my-backups/postgres}"
TS=$(date +%Y%m%d_%H%M%S)
FILE="${BACKUP_DIR}/${DB_NAME}_${TS}.dump"

mkdir -p "${BACKUP_DIR}"

# Full backup (custom format)
pg_dump -Fc -d "${DB_NAME}" -f "${FILE}"

# Verify
pg_restore --list "${FILE}" > /dev/null

# Compress + upload
gzip "${FILE}"
aws s3 cp "${FILE}.gz" "${S3_BUCKET}/"

# Cleanup local > 7 days
find "${BACKUP_DIR}" -name "*.dump.gz" -mtime +7 -delete

echo "Backup complete: ${FILE}.gz"
```

### 4.3 Restore Playbook

```bash
# 1. หยุด traffic (maintenance mode)
kubectl scale deploy/api --replicas=0

# 2. หยุด workers
kubectl scale deploy/worker-* --replicas=0

# 3. Drop & recreate
dropdb pdpa_db && createdb pdpa_db

# 4. Restore
aws s3 cp s3://my-backups/postgres/pdpa_db_20240101_030000.dump.gz .
gunzip pdpa_db_20240101_030000.dump.gz
pg_restore -Fc -d pdpa_db pdpa_db_20240101_030000.dump

# 5. Verify
psql pdpa_db -c "SELECT count(*) FROM pdpa_consents;"

# 6. Start services (ทีละตัว)
kubectl scale deploy/api --replicas=2
kubectl scale deploy/worker-account --replicas=1
# monitor 5 นาที ก่อน scale ตัวอื่น
```

### 4.4 Disaster Recovery Drill (รายไตรมาส)

```markdown
## DR Drill Checklist
- [ ] จำลอง DB ล่ม → restore จาก backup
- [ ] วัด RTO จริง (ควร < 2 ชม.)
- [ ] ทดสอบ failover Kafka → replica
- [ ] ทดสอบ Redis cluster promote
- [ ] ตรวจสอบ log/trace ครบระหว่าง incident
- [ ] บันทึกบทเรียน (post-mortem)
```

---

## 5. Dependency Upgrade Playbook

### 5.1 ตารางตรวจสอบ Dependency

```bash
# ดู outdated ทั้งหมด
go list -u -m all

# ดู vulnerability
govulncheck ./...

# ดู license
go-licenses check ./...
```

### 5.2 ขั้นตอน Upgrade (ปลอดภัย)

```bash
# 1. ตรวจสอบ breaking changes
#    อ่าน CHANGELOG ของ dependency

# 2. Upgrade minor/patch ก่อน (ปลอดภัย)
go get -u=patch ./...
go mod tidy
go test ./...

# 3. Upgrade major ทีละตัว
go get github.com/gin-gonic/gin@v1.10.0
go mod tidy
go test ./...

# 4. ถ้าผ่าน → commit
git add go.mod go.sum
git commit -m "chore(deps): bump gin to v1.10.0"

# 5. Rollback ถ้าพัง
git checkout go.mod go.sum
```

### 5.3 Critical Dependencies Matrix

| Package | Version ปัจจุบัน | Stability | หมายเหตุ |
| :--- | :--- | :--- | :--- |
| `gin-gonic/gin` | v1.9.1 | Stable | ระวัง v2 breaking |
| `gorm.io/gorm` | v1.25.5 | Stable | – |
| `IBM/sarama` | v1.41.3 | Stable | IBM fork ของ Shopify |
| `go-redis/redis/v8` | v8.11.5 | Stable | v9 มี breaking |
| `ethereum/go-ethereum` | v1.13.5 | Beta | API เปลี่ยนบ่อย |

---

## 6. Secret & Config Management

### 6.1 กฎการจัดการ Secret

```
✅ ทำ:
  - ใช้ HashiCorp Vault / AWS Secrets Manager
  - Inject ผ่าน env ตอน runtime
  - Rotate ทุก 90 วัน
  - Audit การเข้าถึง

❌ ห้าม:
  - commit .env ลง git
  - hardcode secret ใน source
  - log secret ทุกกรณี
  - แชร์ secret ผ่าน chat
```

### 6.2 Config Hierarchy

```
1. Code default (ค่าที่ปลอดภัย)
2. config.yaml (ใน repo, ไม่มี secret)
3. .env (local only, gitignore)
4. Environment variables (production)
5. Vault (secret จริง)
```

**`internal/config/config.go`**
```go
package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DB           DBConfig
	Redis        RedisConfig
	Kafka        KafkaConfig
	JWT          JWTConfig
	Observability ObservabilityConfig
}

type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

func Load() *Config {
	return &Config{
		DB: DBConfig{
			DSN:     mustEnv("DB_DSN"),
			MaxOpen: envInt("DB_MAX_OPEN", 25),
			MaxIdle: envInt("DB_MAX_IDLE", 10),
		},
		JWT: JWTConfig{
			Secret: mustEnv("JWT_SECRET"),  // panic ถ้าไม่มี
			TTL:    envDuration("JWT_TTL", 24*time.Hour),
		},
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("required env missing: " + key)
	}
	return v
}
```

### 6.3 Secret Rotation Procedure

```markdown
## JWT Secret Rotation
1. สร้าง secret ใหม่ (เก็บใน Vault)
2. Deploy code ที่รับ 2 secrets (เก่า + ใหม่)
   - verify: ลองทั้ง 2
   - sign: ใช้ใหม่เท่านั้น
3. รอ TTL ของ token เก่าหมดอายุ (24 ชม.)
4. ลบ secret เก่า
5. Deploy code ที่รับ secret เดียว
```

---

# ภาค 2 — การแก้ไข (Edit Guide)

## 7. Workflow การแก้ไขที่ปลอดภัย

### 7.1 ขั้นตอน 8 ขั้น

```
1. ระบุ Layer ที่ต้องแก้ (ใช้ decision tree ด้านล่าง)
2. เขียน test ที่ fail (TDD)
3. แก้ code
4. รัน test → ผ่าน
5. รัน build + vet + lint
6. ทดสอบ local (docker-compose up)
7. สร้าง PR + อธิบาย change
8. Code review + merge
```

### 7.2 Decision Tree — แก้ตรงไหน?

```
ปัญหาคืออะไร?
├── กฎทางธุรกิจเปลี่ยน → Domain Layer
│   ├── เงื่อนไขใน entity → entity/{{entity}}.go
│   ├── validation ค่า → value_object/
│   └── logic ข้าม entity → domain/service/
│
├── flow การทำงานเปลี่ยน → Application Layer
│   └── application/{{verb}}_{{entity}}.go
│
├── เปลี่ยน technology → Infrastructure Layer
│   ├── เปลี่ยน DB → persistence/postgres/
│   ├── เปลี่ยน message broker → messaging/
│   └── เปลี่ยน external service → services/
│
├── เปลี่ยน API/UI → Interface Layer
│   ├── เปลี่ยน route → interfaces/http/routes.go
│   ├── เปลี่ยน request/response → handler + dto
│   └── เปลี่ยน middleware → interfaces/middleware/
│
└── เปลี่ยน schema → migrations/
```

### 7.3 กฎการแก้ไข (ห้ามละเมิด)

| ❌ ห้าม | ✅ ทำแทน |
| :--- | :--- |
| แก้ domain ให้ import gorm | สร้าง repository impl |
| ใส่ business logic ใน handler | ย้ายไป use case |
| เรียก DB ตรงจาก use case | ผ่าน repository interface |
| แก้ migration ที่ deploy แล้ว | สร้าง migration ใหม่ |
| ลบ field ใน entity ทันที | deprecate + migrate |
| เปลี่ยน signature ของ public interface | สร้าง v2 หรือ overload |

---

## 8. แก้ Domain Layer

### 8.1 เพิ่ม Behavior ใน Entity

**กรณี:** ต้องการให้ entity ทำอะไรใหม่ได้

```go
// เดิม
func (c *ConsentLog) Revoke() error {
	if c.Status == valueobject.ConsentRevoked {
		return domainerrors.ErrConsentAlreadyRevoked
	}
	c.Status = valueobject.ConsentRevoked
	now := time.Now()
	c.RevokedAt = &now
	return nil
}

// เพิ่ม behavior ใหม่
func (c *ConsentLog) Extend(additional time.Duration) error {
	if c.Status != valueobject.ConsentGranted {
		return domainerrors.ErrConsentNotActive
	}
	if additional <= 0 {
		return domainerrors.ErrInvalidExtension
	}
	c.ExpiresAt = c.ExpiresAt.Add(additional)
	return nil
}
```

**ขั้นตอน:**
1. เขียน test ที่ `domain/entity/{{entity}}_test.go`
2. เพิ่ม method พร้อม validation
3. ถ้าต้อง persist → เพิ่มใน repository interface
4. อัปเดต repository impl

### 8.2 เพิ่ม Invariant ใน Constructor

```go
func NewConsentLog(userID uuid.UUID, ...) (*ConsentLog, error) {
	if userID == uuid.Nil {
		return nil, domainerrors.ErrInvalidUserID
	}
	if !purpose.IsValid() {
		return nil, domainerrors.ErrInvalidPurpose
	}
	// ...
	return &ConsentLog{...}, nil
}
```

> ⚠️ **Breaking change:** ถ้าเปลี่ยน signature ของ constructor ต้องแก้ทุก caller

### 8.3 เปลี่ยน Value Object

```go
// เพิ่มค่าที่อนุญาต
const (
	PurposeNecessary ConsentPurpose = "NECESSARY"
	PurposeAnalytics ConsentPurpose = "ANALYTICS"
	PurposeMarketing ConsentPurpose = "MARKETING"
	PurposeProfiling ConsentPurpose = "PROFILING"  // ใหม่
)

func (c ConsentPurpose) IsValid() bool {
	switch c {
	case PurposeNecessary, PurposeAnalytics, PurposeMarketing, PurposeProfiling:
		return true
	}
	return false
}
```

**ต้องอัปเดต:**
- [ ] Migration: `INSERT INTO pdpa_purposes` ค่าใหม่
- [ ] i18n: เพิ่ม translation
- [ ] Documentation

---

## 9. แก้ Application Layer

### 9.1 เพิ่ม Field ใน Input DTO

```go
type RecordConsentInput struct {
	UserID    uuid.UUID
	SessionID string
	Purposes  map[string]bool
	IPAddress string
	UserAgent string
	Source    string  // ใหม่: web/mobile/api
}
```

**ต้องอัปเดต:**
- [ ] Handler ที่ map จาก HTTP request
- [ ] Test
- [ ] Audit trail payload (ถ้าใช้)

### 9.2 เพิ่ม Use Case เข้า Flow เดิม

```go
// เดิม Execute() ทำ 3 อย่าง
func (uc *RecordConsentUseCase) Execute(ctx context.Context, input RecordConsentInput) error {
	// 1. validate
	// 2. save
	// 3. audit
}

// เพิ่ม step 4: publish event
func (uc *RecordConsentUseCase) Execute(ctx context.Context, input RecordConsentInput) error {
	// ... 1-3
	// 4. publish
	if err := uc.producer.PublishConsentEvent(ctx, log); err != nil {
		// log แต่ไม่ fail
		logger.Warn("publish failed", logger.Error(err))
	}
	return nil
}
```

**กฎ:** side effect ที่ไม่ critical → log warn แต่ไม่ fail  
side effect ที่ critical → return error (เช่น การชำระเงิน)

### 9.3 แยก Use Case ใหญ่ออกเป็นหลายตัว

**ก่อน:**
```go
func (uc *ConsentUseCase) Execute(ctx, input) error {
	// 200 บรรทัด
}
```

**หลัง:**
```go
type ConsentUseCase struct {
	validateUC *ValidateConsentUseCase
	saveUC     *SaveConsentUseCase
	notifyUC   *NotifyConsentUseCase
}

func (uc *ConsentUseCase) Execute(ctx, input) error {
	if err := uc.validateUC.Execute(ctx, input); err != nil { return err }
	if err := uc.saveUC.Execute(ctx, input); err != nil { return err }
	return uc.notifyUC.Execute(ctx, input)
}
```

---

## 10. แก้ Infrastructure Layer

### 10.1 เปลี่ยน Database (เช่น PostgreSQL → MongoDB)

**ขั้นตอน:**
1. สร้าง repository impl ใหม่ใน `persistence/mongodb/`
2. **ไม่แก้** domain repository interface
3. Wire-up ใน `module.go` เลือก impl ตาม env

```go
// module.go
var repo repository.ConsentRepository
switch os.Getenv("DB_DRIVER") {
case "mongo":
	repo = mongodb.NewConsentRepository(mongoClient)
default:
	repo = postgres.NewConsentRepository(db)
}
```

### 10.2 เพิ่ม Index ใน Elasticsearch

```go
// infrastructure/search/elasticsearch/consent_indexer.go
func (i *consentIndexer) ensureIndex(ctx context.Context) error {
	mapping := `{
		"mappings": {
			"properties": {
				"user_id":    {"type": "keyword"},
				"purpose":    {"type": "keyword"},
				"granted_at": {"type": "date"},
				"status":     {"type": "keyword"}
			}
		}
	}`
	res, err := i.client.Indices.Create(i.index,
		i.client.Indices.Create.WithBody(strings.NewReader(mapping)))
	// ...
}
```

> ⚠️ เปลี่ยน mapping ของ index ที่มีข้อมูลแล้ว → ต้อง reindex

### 10.3 ปรับ Kafka Producer

```go
// เพิ่ม headers สำหรับ tracing
func (k *kafkaProducer) Publish(ctx context.Context, topic, key string, payload interface{}) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
	}
	// เพิ่ม trace context
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		msg.Headers = append(msg.Headers, sarama.RecordHeader{
			Key:   []byte("trace_id"),
			Value: []byte(span.SpanContext().TraceID().String()),
		})
	}
	_, _, err := k.p.SendMessage(msg)
	return err
}
```

### 10.4 เปลี่ยน Retry Policy ของ Consumer

```go
// consumers/base.go
type RetryPolicy struct {
	MaxAttempts int
	Backoff     time.Duration
}

func (c *{{Topic}}Consumer) ConsumeClaim(s sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var lastErr error
		for attempt := 1; attempt <= c.retry.MaxAttempts; attempt++ {
			if err := c.handle(msg); err == nil {
				lastErr = nil
				break
			}
			time.Sleep(c.retry.Backoff * time.Duration(attempt))
		}
		if lastErr != nil {
			// ส่งเข้า DLQ
			c.dlq.Publish(msg)
		}
		s.MarkMessage(msg, "")
	}
	return nil
}
```

---

## 11. แก้ Interface Layer

### 11.1 เพิ่ม Endpoint ใหม่ใน Handler เดิม

```go
// interfaces/http/consent_handler.go
func (h *ConsentHandler) GetHistory(c *gin.Context) {
	uid := mustUserID(c)
	history, err := h.getHistoryUC.Execute(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": history})
}

// routes.go
pdpa.GET("/consent/history", h.Consent.GetHistory)
```

### 11.2 เปลี่ยน Request/Response Format (Breaking)

**ใช้ API Versioning:**
```go
v1 := r.Group("/api/v1")
v2 := r.Group("/api/v2")

v1.POST("/pdpa/consent", h.V1.RecordConsent)  // format เดิม
v2.POST("/pdpa/consent", h.V2.RecordConsent)  // format ใหม่
```

**หรือ Deprecation Header:**
```go
func (h *ConsentHandler) RecordConsent(c *gin.Context) {
	c.Header("Deprecation", "true")
	c.Header("Sunset", "2025-12-31")
	// ...
}
```

### 11.3 เพิ่ม Middleware

```go
// interfaces/middleware/request_id.go
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = uuid.New().String()
		}
		c.Set("request_id", rid)
		c.Header("X-Request-ID", rid)
		c.Next()
	}
}

// routes.go
r.Use(middleware.RequestID(), middleware.Logging(), middleware.Recovery())
```

---

## 12. แก้ Database Migration

### 12.1 กฎเหล็ก

```
❌ ห้ามแก้ migration ที่ deploy แล้ว
✅ สร้าง migration ใหม่เสมอ

❌ ห้าม DROP COLUMN ทันที
✅ deprecate → stop using → DROP (2-3 release)

❌ ห้ามเปลี่ยน type ที่ข้อมูลเยอะโดยไม่มี downtime plan
✅ สร้าง column ใหม่ → copy → swap
```

### 12.2 ตัวอย่าง Migration ที่ปลอดภัย

**เพิ่ม column:**
```sql
-- migrations/20240115_pdpa_add_consent_source.sql
ALTER TABLE pdpa_consents
	ADD COLUMN IF NOT EXISTS source VARCHAR(50) DEFAULT 'unknown';
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_pdpa_consents_source
	ON pdpa_consents (source);
```

**เปลี่ยนชื่อ column (แบบปลอดภัย):**
```sql
-- Step 1: เพิ่ม column ใหม่
ALTER TABLE pdpa_consents ADD COLUMN ip_addr VARCHAR(45);

-- Step 2: copy data (background)
UPDATE pdpa_consents SET ip_addr = ip_address WHERE ip_addr IS NULL;

-- Step 3: (หลัง deploy code ที่อ่าน ip_addr) DROP column เก่า
ALTER TABLE pdpa_consents DROP COLUMN ip_address;
```

**เพิ่ม NOT NULL:**
```sql
-- Step 1: เพิ่มแบบ nullable + default
ALTER TABLE pdpa_consents ADD COLUMN version INT DEFAULT 1;

-- Step 2: backfill (ถ้าจำเป็น)
UPDATE pdpa_consents SET version = 1 WHERE version IS NULL;

-- Step 3: บังคับ NOT NULL
ALTER TABLE pdpa_consents ALTER COLUMN version SET NOT NULL;
```

### 12.3 Migration Tool

```bash
# ใช้ golang-migrate
migrate -path migrations -database "postgres://..." up
migrate -path migrations -database "postgres://..." down 1
migrate -path migrations -database "postgres://..." version
```

**`cmd/migrate/main.go`**
```go
package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	m, err := migrate.New("file://migrations", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	switch cmd {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "version":
		v, dirty, _ := m.Version()
		log.Printf("version=%d dirty=%v", v, dirty)
		return
	}
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}
```

---

## 13. Breaking Change Playbook

### 13.1 ลำดับการเปลี่ยน (Expand-Contract)

```
Phase 1: EXPAND
├── เพิ่ม field/table/endpoint ใหม่
├── deploy code ที่เขียนทั้งเก่า+ใหม่
└── data ยังใช้เก่า

Phase 2: MIGRATE
├── backfill data
├── switch reads ไปใหม่
└── monitor

Phase 3: CONTRACT
├── หยุดเขียนเก่า
├── รอ 1-2 release
└── ลบของเก่า
```

### 13.2 ตัวอย่าง: เปลี่ยน Field Name

```
Day 0:  ADD COLUMN new_field
Day 1:  Code writes to BOTH old+new
Day 7:  Backfill old→new
Day 14: Code reads from new
Day 30: Code stops writing old
Day 45: DROP COLUMN old_field
```

### 13.3 Rollback Plan

ทุก breaking change ต้องมี:
```markdown
## Rollback Procedure
1. **Code:** `git revert <sha>` + deploy
2. **DB:** `migrate down 1` (ถ้ามี down script)
3. **Data:** restore จาก snapshot (ถ้าจำเป็น)
4. **Cache:** flush Redis keys: `DEL pdpa:consent:*`
5. **Kafka:** replay จาก offset ที่กำหนด
```

---

# ภาค 3 — การขยายฟีเจอร์ (Extension Guide)

## 14. เพิ่ม Use Case ใหม่

### 14.1 Checklist

```markdown
- [ ] ตั้งชื่อ: `{{verb}}_{{entity}}.go` (เช่น `archive_consent.go`)
- [ ] สร้าง Input/Output DTO ในไฟล์เดียวกัน
- [ ] Inject dependencies ผ่าน constructor
- [ ] Execute() เป็น entry เดียว, return error
- [ ] เขียน test (`application/{{verb}}_test.go`)
- [ ] Wire-up ใน `module.go`
- [ ] สร้าง handler ที่เรียก use case
- [ ] เพิ่ม route
- [ ] เพิ่ม audit trail (ถ้าเป็นการเปลี่ยน state)
- [ ] เพิ่ม metric
```

### 14.2 Template

```go
// application/archive_consent.go
package application

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/repository"
)

type ArchiveConsentUseCase struct {
	repo      repository.ConsentRepository
	auditRepo repository.AuditRepository
}

func NewArchiveConsentUseCase(
	repo repository.ConsentRepository,
	auditRepo repository.AuditRepository,
) *ArchiveConsentUseCase {
	return &ArchiveConsentUseCase{repo: repo, auditRepo: auditRepo}
}

type ArchiveConsentInput struct {
	UserID uuid.UUID
	Reason string
}

func (uc *ArchiveConsentUseCase) Execute(ctx context.Context, input ArchiveConsentInput) error {
	consents, err := uc.repo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return err
	}
	for i := range consents {
		c := &consents[i]
		if err := c.Archive(); err != nil {
			return err
		}
		if err := uc.repo.Save(ctx, c); err != nil {
			return err
		}
	}
	return nil
}
```

### 14.3 Wire-up

```go
// module.go
archiveUC := application.NewArchiveConsentUseCase(consentRepo, auditRepo)
handler := httpiface.NewConsentHandler(recordUC, revokeUC, archiveUC)
```

---

## 15. เพิ่ม Entity / Aggregate ใหม่

### 15.1 ขั้นตอน (10 ขั้น)

```
1. domain/entity/{{new_entity}}.go     ← Entity + behavior
2. domain/repository/{{new_entity}}_repository.go  ← Interface
3. domain/errors/errors.go             ← เพิ่ม error
4. migrations/YYYYMMDD_{{module}}_{{new_entity}}.sql  ← Table
5. infrastructure/persistence/postgres/models.go  ← GORM model
6. infrastructure/persistence/postgres/{{new_entity}}_repo_impl.go
7. application/{{verb}}_{{new_entity}}.go  ← UseCase
8. interfaces/http/{{new_entity}}_handler.go
9. interfaces/http/routes.go           ← Register
10. module.go                          ← Wire
```

### 15.2 ตัวอย่าง: เพิ่ม `ConsentTemplate` Entity

**1. Entity**
```go
// domain/entity/consent_template.go
package entity

type ConsentTemplate struct {
	ID        uuid.UUID
	Language  string
	Purpose   valueobject.ConsentPurpose
	Title     string
	Body      string
	Version   int
	IsActive  bool
	CreatedAt time.Time
}

func NewConsentTemplate(lang string, purpose valueobject.ConsentPurpose, title, body string) (*ConsentTemplate, error) {
	if !purpose.IsValid() {
		return nil, domainerrors.ErrInvalidPurpose
	}
	if title == "" || body == "" {
		return nil, domainerrors.ErrEmptyTemplate
	}
	return &ConsentTemplate{
		ID:        uuid.New(),
		Language:  lang,
		Purpose:   purpose,
		Title:     title,
		Body:      body,
		Version:   1,
		IsActive:  true,
		CreatedAt: time.Now(),
	}, nil
}
```

**2. Repository**
```go
// domain/repository/consent_template_repository.go
type ConsentTemplateRepository interface {
	Save(ctx context.Context, t *entity.ConsentTemplate) error
	FindActive(ctx context.Context, lang string, purpose valueobject.ConsentPurpose) (*entity.ConsentTemplate, error)
	ListByPurpose(ctx context.Context, purpose valueobject.ConsentPurpose) ([]entity.ConsentTemplate, error)
}
```

**3-10** ทำตาม pattern เดียวกับ Consent

---

## 16. เพิ่ม Value Object ใหม่

```go
// domain/value_object/language.go
package valueobject

type Language string

const (
	LanguageEN Language = "en"
	LanguageTH Language = "th"
	LanguageZH Language = "zh"
)

func (l Language) IsValid() bool {
	switch l {
	case LanguageEN, LanguageTH, LanguageZH:
		return true
	}
	return false
}

func (l Language) String() string { return string(l) }

func ParseLanguage(s string) (Language, error) {
	l := Language(s)
	if !l.IsValid() {
		return "", errors.New("invalid language: " + s)
	}
	return l, nil
}
```

---

## 17. เพิ่ม Kafka Topic / Consumer

### 17.1 ตั้งชื่อ Topic

```
<module>.<entity>.<event>
เช่น: pdpa.consent.expired, pdpa.dsar.approved
```

### 17.2 Checklist

```markdown
- [ ] Producer: เพิ่ม method ใน KafkaProducer interface
- [ ] Publish ที่ use case ที่เกี่ยวข้อง
- [ ] Consumer: สร้าง `consumers/{{topic}}_consumer.go`
- [ ] สร้าง `cmd/workers/{{topic}}/main.go`
- [ ] Docker compose: เพิ่ม service
- [ ] Monitor: consumer lag alert
- [ ] DLQ: ตั้งค่า dead letter queue
```

### 17.3 Producer + Consumer

```go
// Producer method
func (p *kafkaProducer) PublishConsentExpired(ctx context.Context, consent *entity.ConsentLog) error {
	return p.publish("pdpa.consent.expired", consent.UserID.String(), map[string]interface{}{
		"event_type": "consent.expired",
		"consent_id": consent.ID,
		"user_id":    consent.UserID,
		"expired_at": consent.ExpiresAt,
	})
}

// Consumer
// infrastructure/messaging/consumers/consent_expired_consumer.go
type ConsentExpiredHandler struct {
	notifyUC *application.NotifyExpiryUseCase
}

func (h *ConsentExpiredHandler) Handle(ctx context.Context, payload map[string]interface{}) error {
	userID, _ := uuid.Parse(payload["user_id"].(string))
	return h.notifyUC.Execute(ctx, application.NotifyExpiryInput{UserID: userID})
}
```

---

## 18. เพิ่ม Worker ใหม่

### 18.1 Template `cmd/workers/{{name}}/main.go`

```go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"

	"icmongolang/internal/modules/{{module_name}}/application"
	"icmongolang/internal/modules/{{module_name}}/infrastructure/messaging/consumers"
	"icmongolang/internal/modules/{{module_name}}/infrastructure/persistence/postgres"
)

func main() {
	_ = godotenv.Load()

	db := mustDB()
	repo := postgres.New{{Entity}}Repository(db)
	uc := application.New{{Verb}}{{Entity}}UseCase(repo, nil)

	handler := consumers.New{{Topic}}Consumer(uc)

	cfg := sarama.NewConfig()
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Return.Errors = true

	group, err := sarama.NewConsumerGroup(
		[]string{os.Getenv("KAFKA_BROKERS")},
		"{{module_name}}-{{topic}}-group",
		cfg,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer group.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle consumer errors
	go func() {
		for err := range group.Errors() {
			log.Printf("[{{topic}}] consumer error: %v", err)
		}
	}()

	// Consume loop
	go func() {
		for {
			if err := group.Consume(ctx, []string{"{{topic}}"}, handler); err != nil {
				log.Printf("[{{topic}}] consume error: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	log.Println("[{{topic}}] worker started")

	// Graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("[{{topic}}] shutting down...")
	cancel()
	time.Sleep(5 * time.Second) // ให้ in-flight messages จบ
}
```

---

## 19. เพิ่ม Scheduler Job

### 19.1 Job Template

```go
// infrastructure/scheduler/{{job}}_job.go
package scheduler

import (
	"context"
	"time"

	"icmongolang/pkg/logger"
)

type {{Job}}Job struct {
	// deps
	interval time.Duration
}

func New{{Job}}Job(/* deps */) *{{Job}}Job {
	return &{{Job}}Job{interval: 24 * time.Hour}
}

func (j *{{Job}}Job) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	start := time.Now()
	logger.Info("[{{job}}] start")

	if err := j.execute(ctx); err != nil {
		logger.Error("[{{job}}] failed", logger.Error(err))
		return
	}
	logger.Info("[{{job}}] done", logger.Duration("took", time.Since(start)))
}

func (j *{{Job}}Job) execute(ctx context.Context) error {
	// implementation
	return nil
}
```

### 19.2 Register ใน Scheduler

```go
// cmd/scheduler/main.go
c := cron.New(cron.WithChain(
	cron.Recover(cron.DefaultLogger),
	cron.SkipIfStillRunning(cron.DefaultLogger),
))

_, _ = c.AddFunc("0 2 * * *", job1.Run)  // ทุกวัน 02:00
_, _ = c.AddFunc("@every 6h", job2.Run)  // ทุก 6 ชม.
c.Start()
```

---

## 20. เพิ่ม HTTP Endpoint

### 20.1 Endpoint Naming

```
GET    /api/v1/pdpa/consents           # list
GET    /api/v1/pdpa/consents/:id       # get one
POST   /api/v1/pdpa/consents           # create
PUT    /api/v1/pdpa/consents/:id       # full update
PATCH  /api/v1/pdpa/consents/:id       # partial update
DELETE /api/v1/pdpa/consents/:id       # delete
```

### 20.2 Handler Template

```go
func (h *ConsentHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	consent, err := h.getUC.Execute(c.Request.Context(), application.GetConsentInput{ID: id})
	switch {
	case errors.Is(err, domainerrors.ErrConsentNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, consent)
}
```

### 20.3 Standard Response Format

```json
// Success
{
  "data": {...},
  "meta": {"request_id": "uuid", "timestamp": "..."}
}

// Error
{
  "error": {
    "code": "CONSENT_NOT_FOUND",
    "message": "consent record not found",
    "details": {"id": "uuid"}
  },
  "meta": {"request_id": "uuid", "timestamp": "..."}
}
```

---

## 21. เพิ่ม WebSocket Event

### 21.1 Event Types

```go
// domain/event/websocket_events.go
package event

type EventType string

const (
	EventDSARUpdated      EventType = "DSAR_UPDATED"
	EventConsentChanged   EventType = "CONSENT_CHANGED"
	EventDeletionDone     EventType = "DELETION_COMPLETED"
	EventAccountSuspended EventType = "ACCOUNT_SUSPENDED"
)

type WSMessage struct {
	Type      EventType   `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}
```

### 21.2 Broadcast จาก Use Case

```go
// ใน use case
uc.wsHub.SendToUser(userID.String(), event.WSMessage{
	Type:      event.EventConsentChanged,
	Timestamp: time.Now(),
	Payload: map[string]interface{}{
		"purpose": purpose.String(),
		"status":  "revoked",
	},
})
```

### 21.3 Client-side (JavaScript)

```javascript
const ws = new WebSocket(`wss://api.example.com/ws?token=${jwt}`);
ws.onmessage = (e) => {
  const msg = JSON.parse(e.data);
  switch (msg.type) {
    case 'DSAR_UPDATED':
      updateDSARUI(msg.payload);
      break;
    case 'CONSENT_CHANGED':
      showNotification(msg.payload);
      break;
  }
};
```

---

## 22. เพิ่ม Outbound Port ใหม่

### 22.1 ขั้นตอน (5 ขั้น)

```
1. ประกาศ interface ใน domain/service/{{port}}_port.go
2. Use case ใช้ port (dependency injection)
3. สร้าง impl ใน infrastructure/services/{{name}}/
4. Wire ใน module.go
5. เขียน mock สำหรับ test
```

### 22.2 ตัวอย่าง: เพิ่ม SMS Sender

**1. Port**
```go
// domain/service/sms_sender.go
package service

import "context"

type SMSSender interface {
	Send(ctx context.Context, phone, message string) error
}
```

**2. Use Case ใช้**
```go
type SubmitDSARUseCase struct {
	sms service.SMSSender
	// ...
}

func (uc *SubmitDSARUseCase) Execute(ctx context.Context, input SubmitDSARInput) error {
	// ...
	if err := uc.sms.Send(ctx, input.Phone, "OTP: "+otp); err != nil {
		logger.Warn("sms failed", logger.Error(err))
	}
	return nil
}
```

**3. Implementation**
```go
// infrastructure/services/sms/twilio.go
package sms

import (
	"context"
	"github.com/twilio/twilio-go"
)

type TwilioSender struct {
	client *twilio.RestClient
	from   string
}

func NewTwilioSender(accountSID, authToken, from string) *TwilioSender {
	return &TwilioSender{
		client: twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: accountSID,
			Password: authToken,
		}),
		from: from,
	}
}

func (s *TwilioSender) Send(ctx context.Context, phone, message string) error {
	_, err := s.client.Api.CreateMessage(&openapi.CreateMessageParams{
		To:   &phone,
		From: &s.from,
		Body: &message,
	})
	return err
}
```

**4. Mock สำหรับ Test**
```go
// infrastructure/services/sms/mock.go
type MockSender struct {
	Messages []string
	Err      error
}

func (m *MockSender) Send(ctx context.Context, phone, message string) error {
	if m.Err != nil {
		return m.Err
	}
	m.Messages = append(m.Messages, message)
	return nil
}
```

---

# ภาค 4 — การเสริมหลัง (Advanced Extensions)

## 23. Observability Stack

### 23.1 OpenTelemetry Integration

```go
// pkg/otel/setup.go
package otel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func Setup(ctx context.Context, serviceName, endpoint string) (func(), error) {
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.1))),
	)
	otel.SetTracerProvider(tp)

	return func() { _ = tp.Shutdown(ctx) }, nil
}
```

### 23.2 Instrument Use Case

```go
func (uc *RecordConsentUseCase) Execute(ctx context.Context, input RecordConsentInput) error {
	ctx, span := otel.Tracer("pdpa").Start(ctx, "RecordConsentUseCase.Execute")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", input.UserID.String()),
	)
	// ...
}
```

### 23.3 Stack แนะนำ

| Component | เครื่องมือ | บทบาท |
| :--- | :--- | :--- |
| Metrics | Prometheus + Grafana | ตัวเลข, alert |
| Logs | Loki / ELK | structured logs |
| Traces | Tempo / Jaeger | distributed tracing |
| Profiling | Pyroscope / pprof | CPU/memory |
| APM | OpenTelemetry Collector | รวมทุก signal |

---

## 24. Resilience Patterns

### 24.1 Circuit Breaker

```go
// pkg/resilience/circuit.go
package resilience

import "github.com/sony/gobreaker"

func NewCircuitBreaker(name string) *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        name,
		MaxRequests: 5,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
	})
}

// ใช้
func (c *BlockchainClient) RecordHash(ctx context.Context, data string) (string, error) {
	result, err := c.cb.Execute(func() (interface{}, error) {
		return c.doRecordHash(ctx, data)
	})
	if err != nil {
		return "", err
	}
	return result.(string), nil
}
```

### 24.2 Retry with Backoff

```go
// pkg/resilience/retry.go
func Retry(ctx context.Context, attempts int, baseDelay time.Duration, fn func() error) error {
	var lastErr error
	for i := 0; i < attempts; i++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
		}
		select {
		case <-time.After(baseDelay * time.Duration(1<<i)):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return lastErr
}
```

### 24.3 Bulkhead (แยก resource pool)

```go
// แยก worker pool สำหรับงานแต่ละประเภท
type Bulkhead struct {
	sem chan struct{}
}

func NewBulkhead(size int) *Bulkhead {
	return &Bulkhead{sem: make(chan struct{}, size)}
}

func (b *Bulkhead) Execute(ctx context.Context, fn func() error) error {
	select {
	case b.sem <- struct{}{}:
		defer func() { <-b.sem }()
		return fn()
	case <-ctx.Done():
		return ctx.Err()
	}
}
```

### 24.4 Timeout

```go
// ใช้ context.WithTimeout เสมอสำหรับ external call
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
resp, err := httpClient.Do(req.WithContext(ctx))
```

---

## 25. Feature Flags & A/B Testing

### 25.1 Feature Flag Interface

```go
// domain/service/feature_flag.go
package service

import "context"

type FeatureFlag interface {
	IsEnabled(ctx context.Context, key string, userID string) bool
	Variant(ctx context.Context, key string, userID string) string
}
```

### 25.2 Implementation (Unleash / Local)

```go
// infrastructure/services/featureflag/local.go
package featureflag

type LocalFlags struct {
	flags map[string]bool
}

func (f *LocalFlags) IsEnabled(ctx context.Context, key, userID string) bool {
	return f.flags[key]
}
```

### 25.3 ใช้ใน Use Case

```go
if uc.flags.IsEnabled(ctx, "pdpa.new_consent_flow", input.UserID.String()) {
	// flow ใหม่
} else {
	// flow เก่า
}
```

---

## 26. Multi-tenancy

### 26.1 กลยุทธ์ 3 แบบ

| Strategy | ข้อดี | ข้อเสีย | เหมาะกับ |
| :--- | :--- | :--- | :--- |
| **Shared DB, shared schema** | ง่าย, ประหยัด | security ต้องใช้ RLS | SaaS เล็ก |
| **Shared DB, separate schema** | แยกชัด, backup แยก | migrate ซับซ้อน | SaaS กลาง |
| **Separate DB** | แยกสมบูรณ์ | ต้นทุนสูง | Enterprise |

### 26.2 Row-Level Security (PostgreSQL)

```sql
-- เปิด RLS
ALTER TABLE pdpa_consents ENABLE ROW LEVEL SECURITY;

-- Policy: user เห็นแค่ tenant ของตัวเอง
CREATE POLICY tenant_isolation ON pdpa_consents
	USING (tenant_id = current_setting('app.current_tenant')::uuid);

-- Set tenant ตอน connect
SET app.current_tenant = 'tenant-uuid';
```

### 26.3 Middleware

```go
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": "tenant required"})
			return
		}
		c.Set("tenant_id", tenantID)
		c.Next()
	}
}
```

---

## 27. Advanced Caching Strategy

### 27.1 Cache Layers

```
┌─────────────────┐
│   HTTP Cache    │  Cache-Control, ETag
├─────────────────┤
│   CDN           │  CloudFront
├─────────────────┤
│   In-Process    │  sync.Map, LRU (ristretto)
├─────────────────┤
│   Redis         │  Shared cache
├─────────────────┤
│   Database      │  Materialized view
└─────────────────┘
```

### 27.2 Cache Invalidation Strategy

```go
// Write-through
func (uc *RecordConsentUseCase) Execute(ctx, input) error {
	// ... save to DB
	if err := uc.cache.Set(ctx, userID, purpose, status); err != nil {
		logger.Warn("cache set failed")
	}
	return nil
}

// Cache-aside (สำหรับ read)
func (uc *GetConsentUseCase) Execute(ctx, input) (*entity.Consent, error) {
	if cached, err := uc.cache.Get(ctx, input.ID); err == nil && cached != "" {
		return decode(cached), nil
	}
	consent, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	_ = uc.cache.Set(ctx, input.ID, encode(consent))
	return consent, nil
}
```

### 27.3 Anti-Pattern ที่ต้องหลีกเลี่ยง

```
❌ Cache stampede: TTL เดียวกันหมด
✅ Jitter TTL: ttl + rand(0, 10% ttl)

❌ Cache penetration: query ค่าที่ไม่มี
✅ Cache null value สั้นๆ (30s)

❌ Cache inconsistency: write DB สำเร็จ cache fail
✅ Retry + log + metric

❌ Cache key collision: user:123:purpose กับ user:123purpose
✅ ใช้ separator ชัดเจน: user:123:purpose:MARKETING
```

---

## 28. Testing Strategy

### 28.1 Test Pyramid

```
        ┌──────────┐
        │   E2E    │  ← 5%  (critical flows)
        ├──────────┤
        │Integration│  ← 15% (repo + use case)
        ├──────────┤
        │   Unit   │  ← 80%  (domain + app)
        └──────────┘
```

### 28.2 Unit Test (Domain)

```go
// domain/entity/consent_log_test.go
func TestConsentLog_Revoke(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *entity.ConsentLog
		wantErr error
	}{
		{
			name: "granted → revoked",
			setup: func() *entity.ConsentLog {
				return entity.NewConsentLog(uuid.New(), "", "MARKETING", "", "")
			},
			wantErr: nil,
		},
		{
			name: "already revoked → error",
			setup: func() *entity.ConsentLog {
				c := entity.NewConsentLog(uuid.New(), "", "MARKETING", "", "")
				_ = c.Revoke()
				return c
			},
			wantErr: domainerrors.ErrConsentAlreadyRevoked,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setup()
			err := c.Revoke()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}
```

### 28.3 Integration Test (Repository)

```go
// infrastructure/persistence/postgres/consent_repo_impl_test.go
func TestConsentRepo_SaveAndFind(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	db := setupTestDB(t)  // testcontainers
	defer db.Close()

	repo := postgres.NewConsentRepository(db)
	consent := entity.NewConsentLog(uuid.New(), "sess", "MARKETING", "127.0.0.1", "test")
	require.NoError(t, repo.Save(context.Background(), consent))

	found, err := repo.FindByID(context.Background(), consent.ID)
	require.NoError(t, err)
	require.Equal(t, consent.ID, found.ID)
}
```

### 28.4 E2E Test

```go
// tests/e2e/consent_flow_test.go
func TestConsentFlow_E2E(t *testing.T) {
	// 1. Start docker-compose
	// 2. POST /api/v1/pdpa/consent
	// 3. GET /api/v1/pdpa/consent/history
	// 4. Verify via DB query
}
```

### 28.5 Coverage Target

| Layer | Target |
| :--- | :--- |
| Domain | ≥ 90% |
| Application | ≥ 80% |
| Infrastructure | ≥ 60% |
| Interface | ≥ 50% |

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 29. CI/CD Pipeline

### 29.1 GitHub Actions

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test
          POSTGRES_DB: test_db
        ports: [5432:5432]
      redis:
        image: redis:7
        ports: [6379:6379]

    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
          cache: true

      - name: Lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest

      - name: Vet
        run: go vet ./...

      - name: Test
        run: go test -race -coverprofile=coverage.out ./...
        env:
          DB_DSN: postgres://postgres:test@localhost:5432/test_db?sslmode=disable
          REDIS_ADDR: localhost:6379

      - name: Security scan
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v4
```

### 29.2 Build & Deploy

```yaml
# .github/workflows/cd.yml
name: CD

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: docker/setup-buildx-action@v3
      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - uses: docker/build-push-action@v5
        with:
          push: true
          tags: ghcr.io/org/icmongolang:${{ github.ref_name }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to staging
        run: kubectl set image deploy/api api=ghcr.io/org/icmongolang:${{ github.ref_name }}
      - name: Smoke test
        run: ./scripts/smoke-test.sh
      - name: Deploy to production
        if: success()
        run: kubectl rollout status deploy/api
```

### 29.3 Dockerfile (Multi-stage)

```dockerfile
# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/api ./cmd/api

# Stage 2: Runtime
FROM gcr.io/distroless/static-debian12
COPY --from=builder /out/api /api
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/api"]
```

---

## 30. Performance Tuning

### 30.1 Database

```sql
-- หา slow query
SELECT query, calls, mean_exec_time, total_exec_time
FROM pg_stat_statements
ORDER BY mean_exec_time DESC LIMIT 20;

-- เพิ่ม index ที่ขาด
CREATE INDEX CONCURRENTLY idx_pdpa_consents_user_status
	ON pdpa_consents (user_id, status)
	WHERE status = 'GRANTED';

-- Analyze
ANALYZE pdpa_consents;
```

**Connection pool:**
```go
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(25)      // = CPU cores × 2
sqlDB.SetMaxIdleConns(10)
sqlDB.SetConnMaxLifetime(5 * time.Minute)
sqlDB.SetConnMaxIdleTime(1 * time.Minute)
```

### 30.2 Redis

```go
// ใช้ pipeline สำหรับ batch
pipe := client.Pipeline()
for _, id := range ids {
	pipe.Get(ctx, key(id))
}
cmds, _ := pipe.Exec(ctx)
```

### 30.3 Go Runtime

```go
// จำกัด GOMAXPROCS ใน container
import _ "go.uber.org/automaxprocs"

// pprof (เปิดเฉพาะ debug endpoint)
import _ "net/http/pprof"
go func() { http.ListenAndServe("localhost:6060", nil) }()
```

### 30.4 N+1 Query Fix

```go
// ❌ N+1
for _, user := range users {
	consents, _ := repo.FindByUserID(ctx, user.ID)
}

// ✅ Batch
consentsByUser, _ := repo.FindByUserIDs(ctx, userIDs)
for _, user := range users {
	consents := consentsByUser[user.ID]
}
```

---

## 31. Security Hardening

### 31.1 Input Validation

```go
// validate ที่ handler
if len(req.Purpose) > 50 || !validPurpose(req.Purpose) {
	c.JSON(400, gin.H{"error": "invalid"})
	return
}

// validate ที่ domain (defense in depth)
func NewConsentLog(...) (*ConsentLog, error) {
	if !purpose.IsValid() {
		return nil, domainerrors.ErrInvalidPurpose
	}
}
```

### 31.2 SQL Injection Prevention

```go
// ✅ ใช้ parameterized query เสมอ (GORM ทำให้อยู่แล้ว)
db.Where("user_id = ?", userID).Find(&consents)

// ❌ ห้าม string concatenation
db.Where("user_id = '" + userID + "'")  // NEVER
```

### 31.3 Authentication Best Practices

```go
// JWT
- HS256 → RS256 (asymmetric)
- ใส่ exp, iat, nbf, jti
- Rotation: refresh token + short-lived access
- Blacklist: Redis (jti, TTL = exp)

// Password
- bcrypt cost ≥ 12
- ตรวจ breach (HaveIBeenPwned API)
- Rate limit login: 5 attempts / 15 นาที
```

### 31.4 Rate Limiting

```go
// interfaces/middleware/rate_limit.go
func RateLimit(redis *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "rl:" + c.ClientIP() + ":" + c.FullPath()
		count, _ := redis.Incr(c, key).Result()
		if count == 1 {
			redis.Expire(c, key, window)
		}
		if count > int64(limit) {
			c.Header("Retry-After", window.String())
			c.AbortWithStatusJSON(429, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}

// ใช้
pdpa.POST("/consent", middleware.RateLimit(rdb, 10, time.Minute), h.RecordConsent)
```

### 31.5 Security Headers

```go
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}
```

---

# ภาคผนวก

## A. Decision Tree — ฉันควรทำอะไรต่อ?

```
เริ่ม
│
├── เจอ bug?
│   ├── เกิดจาก business rule → แก้ Domain + เขียน test
│   ├── เกิดจาก logic flow → แก้ Application
│   ├── เกิดจาก infra → แก้ Infrastructure + alert
│   └── เกิดจาก user input → แก้ Interface + validation
│
├── อยากเพิ่มฟีเจอร์?
│   ├── ฟีเจอร์ใหม่หมด → สร้าง module ใหม่ (ใช้ Template_Module.md)
│   ├── เพิ่มใน module เดิม → ดูภาค 3 (Extension Guide)
│   └── ปรับฟีเจอร์เดิม → ดูภาค 2 (Edit Guide)
│
├── performance ตก?
│   ├── DB ช้า → section 30.1
│   ├── API ช้า → instrument use case + trace
│   ├── Kafka lag → scale consumer
│   └── memory leak → pprof
│
└── incident?
    └── ดู Runbook ภาคผนวก B
```

---

## B. Runbook ฉุกเฉิน

### B.1 Service Down (P0)

```markdown
1. **ตรวจสอบ:** `curl /healthz`, `kubectl get pods`
2. **ถ้า pod crash loop:**
   - `kubectl logs <pod> --previous`
   - ดู error: config? secret? dependency?
3. **ถ้า dependency ล่ม:**
   - DB → ตรวจ connection pool, slow query
   - Redis → ตรวจ memory, network
   - Kafka → ตรวจ broker, consumer lag
4. **Rollback ถ้าจำเป็น:**
   - `kubectl rollout undo deploy/api`
5. **ประกาศ:** update status page
```

### B.2 Data Corruption

```markdown
1. **หยุดเขียนทันที:** scale down workers, set maintenance mode
2. **ประเมินขอบเขต:** query หา records ที่ผิด
3. **Snapshot:** backup ปัจจุบันก่อนแก้
4. **กู้คืน:** restore จาก backup + replay Kafka (ถ้าจำเป็น)
5. **Root cause:** หาว่ามาจาก bug ไหน, fix, เขียน test
```

### B.3 Security Incident

```markdown
1. **Isolate:** rotate secrets ทันที (JWT, DB, API keys)
2. **Revoke:** blacklist tokens ทั้งหมด
3. **Audit:** ตรวจ log ย้อนหลัง 30 วัน
4. **Notify:** DPO, legal, ผู้ใช้ที่กระทบ (ถ้าจำเป็นตาม PDPA)
5. **Post-mortem:** วิเคราะห์สาเหตุ + ปรับปรุง
```

### B.4 Kafka Lag พุ่งสูง

```markdown
1. **ตรวจสอบ:** consumer group lag
2. **Scale:** เพิ่ม consumer instances (≤ partitions)
3. **ตรวจ poison message:** ถ้ามี → ส่ง DLQ
4. **ถ้า throughput ไม่พอ:** เพิ่ม partitions, batch size
5. **ถ้า downstream ช้า:** ตรวจ DB, external API
```

---

## C. Glossary

| คำ | ความหมาย |
| :--- | :--- |
| **Aggregate Root** | Entity หลักที่เป็นประตูเข้าถึง entity ย่อยในกลุ่มเดียวกัน |
| **Application Layer** | ชั้นใช้ case — orchestrate business flow |
| **CQRS** | Command Query Responsibility Segregation — แยก read/write |
| **DLQ** | Dead Letter Queue — queue สำหรับ message ที่ประมวลผลไม่ได้ |
| **Domain Layer** | ชั้นที่บรรจุ business rules ล้วนๆ |
| **DTO** | Data Transfer Object — object สำหรับขนส่งข้อมูล |
| **Idempotent** | ทำซ้ำได้ผลเท่าเดิม |
| **Infrastructure Layer** | ชั้น adapter — DB, message broker, external service |
| **Interface Layer** | ชั้นรับ input — HTTP, WebSocket, CLI |
| **Port** | Interface ที่ domain ต้องการจากภายนอก |
| **Repository** | Abstraction สำหรับ persistent storage |
| **RPO** | Recovery Point Objective — ข้อมูลหายได้กี่นาที |
| **RTO** | Recovery Time Objective — ระบบล่มได้นานแค่ไหน |
| **Saga** | Pattern สำหรับ distributed transaction |
| **Value Object** | Object ที่ไม่มี identity, เทียบด้วยค่า |

---

## สรุป — วิธีใช้คู่มือนี้

| สถานการณ์ | เปิดที่ |
| :--- | :--- |
| สร้าง module ใหม่ | `Template_Module.md` |
| แก้ bug | ภาค 2 |
| เพิ่มฟีเจอร์ | ภาค 3 |
| ปรับ performance | ภาค 4 + section 30 |
| Incident | ภาคผนวก B |
| Upgrade dependency | section 5 |
| Onboard ทีมใหม่ | อ่านทั้งเล่ม |

> **ปรัชญาสุดท้าย:**  
> "โค้ดที่ดีคือโค้ดที่คนอื่นอ่านแล้วเข้าใจใน 5 นาที"  
> "ระบบที่ดีคือระบบที่พังตอนตี 3 แล้วเราซ่อมได้โดยไม่ต้องคิดเยอะ"