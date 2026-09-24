# 📙 เล่ม 4: คู่มือทดสอบและ Deployment ฉบับสมบูรณ์
## The Complete Guide to Testing & Deploying Go Modules

> **เวอร์ชัน 1.0 (เมษายน 2026)**
> เอกสารต้นแบบสำหรับการทดสอบและ deploy Module สู่ Production
> ครอบคลุม Test Pyramid → CI/CD → K8s → Monitoring → Incident Response

---

# สารบัญ

**ภาคที่ 1: ปรัชญาและกลยุทธ์การทดสอบ**
1. [บทนำ — ทำไมต้องทดสอบอย่างเป็นระบบ](#บทที่-1-บทนำ)
2. [Test Pyramid และกลยุทธ์](#บทที่-2-test-pyramid)
3. [Test Coverage ที่เหมาะสม](#บทที่-3-test-coverage)

**ภาคที่ 2: Unit Testing**
4. [Unit Test สำหรับ Domain Layer](#บทที่-4-domain-testing)
5. [Unit Test สำหรับ Application Layer](#บทที่-5-application-testing)
6. [Mocking และ Test Doubles](#บทที่-6-mocking)
7. [Table-Driven Tests](#บทที่-7-table-driven)
8. [Fuzzing และ Benchmark](#บทที่-8-fuzzing)

**ภาคที่ 3: Integration & E2E Testing**
9. [Integration Test ด้วย Testcontainers](#บทที่-9-testcontainers)
10. [Database Integration Test](#บทที่-10-db-testing)
11. [Kafka Integration Test](#บทที่-11-kafka-testing)
12. [HTTP Handler Test](#บทที่-12-http-testing)
13. [E2E Test](#บทที่-13-e2e-testing)

**ภาคที่ 4: Specialized Testing**
14. [Load Testing ด้วย k6](#บทที่-14-load-testing)
15. [Chaos Testing](#บทที่-15-chaos-testing)
16. [Security Testing](#บทที่-16-security-testing)
17. [Contract Testing](#บทที่-17-contract-testing)

**ภาคที่ 5: CI/CD Pipeline**
18. [CI Pipeline Design](#บทที่-18-ci-pipeline)
19. [GitHub Actions Workflow](#บทที่-19-github-actions)
20. [Artifact Management](#บทที่-20-artifacts)
21. [Container Build & Registry](#บทที่-21-container)

**ภาคที่ 6: Deployment**
22. [Deployment Strategies](#บทที่-22-deployment-strategies)
23. [Blue-Green Deployment](#บทที่-23-blue-green)
24. [Canary Deployment](#บทที่-24-canary)
25. [Kubernetes Deployment](#บทที่-25-kubernetes)
26. [Zero-Downtime Migration](#บทที่-26-zero-downtime)
27. [Rollback Strategy](#บทที่-27-rollback)

**ภาคที่ 7: Observability**
28. [Three Pillars (Metrics, Logs, Traces)](#บทที่-28-three-pillars)
29. [Prometheus & Grafana](#บทที่-29-prometheus)
30. [Structured Logging](#บทที่-30-logging)
31. [Distributed Tracing](#บทที่-31-tracing)
32. [Alerting Strategy](#บทที่-32-alerting)

**ภาคที่ 8: Operations**
33. [Incident Response](#บทที่-33-incident-response)
34. [Runbooks](#บทที่-34-runbooks)
35. [Postmortem Culture](#บทที่-35-postmortem)

**ภาคที่ 9: Case Studies**
36. [Case Study: Payment Module Deployment](#บทที่-36-case-study)

**ภาคผนวก**
- [A. Test Templates](#ภาคผนวก-a)
- [B. CI/CD Templates](#ภาคผนวก-b)
- [C. Kubernetes Manifests](#ภาคผนวก-c)
- [D. Quick Reference Card](#ภาคผนวก-d)

---

# บทที่ 1: บทนำ

## 1.1 ทำไมต้องทดสอบอย่างเป็นระบบ?

### 1.1.1 ต้นทุนของ Bug ที่หลุดไป Production

```
                        Cost to Fix
    ┌──────────────────────────────────────────────┐
    │                                              │
    │  ████  $100    (Development)                 │
    │                                              │
    │  ████████  $1,000  (Testing)                 │
    │                                              │
    │  ██████████████  $10,000  (Staging)          │
    │                                              │
    │  ██████████████████████  $100,000  (Prod)    │
    │                                              │
    │  ████████████████████████████████  $1M       │
    │         (Prod + Reputation)                  │
    │                                              │
    └──────────────────────────────────────────────┘
```

**IBM study:** Bug ที่หลุดไป production แพงกว่า **100x** ที่เจอตอน development

### 1.1.2 ตัวเลขจริงจากอุตสาหกรรม

| Metric | Industry Average | Top Performers |
|---|---|---|
| **Deployment Frequency** | Monthly | 10x/day |
| **Lead Time** | 1-6 months | < 1 day |
| **MTTR** | 1-7 days | < 1 hour |
| **Change Failure Rate** | 15-30% | < 15% |
| **Test Coverage** | 40-60% | 80%+ |

*(อ้างอิง: DORA State of DevOps 2024)*

## 1.2 หลักการ 5 ข้อของการทดสอบ

### หลักการ 1: Test Behavior, Not Implementation

```go
// ❌ แย่ — test implementation
func TestPayment_SetStatus(t *testing.T) {
    p := &Payment{}
    p.Status = "PENDING"  // ← test internal
    assert.Equal(t, "PENDING", p.Status)
}

// ✅ ดี — test behavior
func TestPayment_Lifecycle(t *testing.T) {
    p, _ := NewPayment(...)
    err := p.StartProcessing()
    assert.NoError(t, err)
    assert.Equal(t, PaymentStatusProcessing, p.Status)
}
```

### หลักการ 2: Test Pyramid, Not Ice Cream Cone

```
     ❌ Ice Cream Cone           ✅ Test Pyramid
     
        /  E2E  \                  /  E2E  \   5%
       /─────────\
      /  Manual   \              / Integr  \ 15%
     /─────────────\
    /   Integr    \             /   Unit    \ 80%
   /───────────────\
  /     Unit        \
 /───────────────────\
```

### หลักการ 3: Fast Feedback

```
Unit tests:        < 100ms
Integration tests: < 5s
E2E tests:         < 30s
Full suite:        < 5 min  ← Target
```

### หลักการ 4: Deterministic

```go
// ❌ Flaky — depends on timing
func TestConcurrent(t *testing.T) {
    go doSomething()
    time.Sleep(100 * time.Millisecond)  // ← ❌
    assert.True(t, done)
}

// ✅ Deterministic
func TestConcurrent(t *testing.T) {
    done := make(chan struct{})
    go func() {
        doSomething()
        close(done)
    }()
    
    select {
    case <-done:
        assert.True(t, true)
    case <-time.After(1 * time.Second):
        t.Fatal("timeout")
    }
}
```

### หลักการ 5: Isolated

```go
// ❌ Tests depend on each other
func TestA(t *testing.T) { createUser() }
func TestB(t *testing.T) { useUser() }  // ← depends on TestA

// ✅ Each test independent
func TestCreateUser(t *testing.T) {
    ctx := setupTest()
    defer cleanupTest(ctx)
    // ...
}

func TestUseUser(t *testing.T) {
    ctx := setupTest()
    defer cleanupTest(ctx)
    user := createTestUser(ctx)
    // ...
}
```

## 1.3 Testing ในบริบท Module ที่ขายได้

สำหรับ module ที่ขาย:
- **Test coverage สูง** → ลูกค้าไว้ใจ
- **CI/CD สมบูรณ์** → ลูกค้าเห็นว่า maintain จริงจัง
- **Benchmark มี** → ลูกค้าเทียบ performance
- **Documentation ครบ** → ลด support ticket

**ตัวเลขจริง:**
- Module ที่มี test coverage > 80% → ขายได้ 3x มากกว่า
- Module ที่มี CI/CD → ลด support ticket 40%
- Module ที่มี benchmark → conversion rate +25%

---

# บทที่ 2: Test Pyramid

## 2.1 แต่ละ Layer ของ Test Pyramid

### 2.1.1 Unit Tests (80%)

**Scope:** Function/class เดียว

**ตัวอย่าง:**
- Domain entity behavior
- Value object validation
- Use case logic (with mocked repo)
- Utility functions

**Characteristics:**
- ⚡ เร็ว (< 100ms/test)
- 🔒 Deterministic
- 📦 No I/O
- 🎯 Isolated

**Tools:** `testing`, `testify`, `gomock`

### 2.1.2 Integration Tests (15%)

**Scope:** หลาย components รวมกัน

**ตัวอย่าง:**
- Repository + real DB
- Handler + real HTTP
- Producer + real Kafka

**Characteristics:**
- ⏱ ปานกลาง (1-5s/test)
- 🔌 I/O จริง
- 🐳 ต้อง infra (Testcontainers)

**Tools:** `testcontainers-go`, `dockertest`

### 2.1.3 E2E Tests (5%)

**Scope:** ทั้งระบบ

**ตัวอย่าง:**
- User journey ครบ: signup → login → action
- API contract ครบ

**Characteristics:**
- 🐢 ช้า (10-60s/test)
- 🎭 Fragile
- 💰 แพง
- 📊 Business value สูง

**Tools:** `playwright`, `k6`, custom

## 2.2 Test Scope Matrix

| Layer | Unit | Integration | E2E |
|---|---|---|---|
| **Domain** | 100% | – | – |
| **Application** | 90% | – | – |
| **Repository** | – | 80% | – |
| **Handler** | 60% | – | – |
| **Kafka** | – | 70% | – |
| **Full flow** | – | – | 5-10 |

## 2.3 Test Speed Targets

```yaml
targets:
  unit_test: "50ms p95"
  integration_test: "5s p95"
  e2e_test: "30s p95"
  full_suite: "5min p95"
  CI_pipeline: "15min p95"
```

## 2.4 เมื่อไหร่ใช้ Test Type ไหน

```
ทดสอบ business rule?
└── YES → Unit Test
    ├── Entity behavior → Unit
    ├── Value object → Unit
    └── Use case → Unit (with mock)

ทดสอบ integration?
└── YES → Integration Test
    ├── DB queries → Integration
    ├── Kafka pub/sub → Integration
    └── HTTP → Integration

ทดสอบ user journey?
└── YES → E2E Test
    ├── Signup flow → E2E
    ├── Payment flow → E2E
    └── Checkout → E2E
```

---

# บทที่ 3: Test Coverage

## 3.1 Coverage Types

| Type | ความหมาย | Target |
|---|---|---|
| **Line** | บรรทัดที่ถูก execute | 80% |
| **Branch** | if/else branches | 70% |
| **Function** | functions ที่ถูกเรียก | 85% |
| **Statement** | statements | 80% |

## 3.2 Coverage ต่อ Layer

```
┌─────────────────────┬──────────┐
│  Layer              │ Coverage │
├─────────────────────┼──────────┤
│  Domain (Entity)    │ 100%     │ ← Critical!
│  Domain (VO)        │ 100%     │
│  Application (UC)   │ 90%      │
│  Infrastructure     │ 70%      │
│  Interfaces         │ 60%      │
│  Utils              │ 90%      │
└─────────────────────┴──────────┘
```

**ทำไม Domain ต้อง 100%?**
- Business rule สำคัญที่สุด
- Bug ใน Domain = money loss
- Test domain เร็ว → ไม่มี excuse

## 3.3 Measuring Coverage

```bash
# Basic
go test -cover ./...

# Detailed report
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# HTML report
go tool cover -html=coverage.out

# Per package
go test -cover ./internal/modules/payment/...
```

**ผลลัพธ์:**
```
internal/modules/payment/domain/entity/payment.go:45:  NewPayment    100.0%
internal/modules/payment/domain/entity/payment.go:80:  Refund        100.0%
internal/modules/payment/application/create.go:25:     Execute        92.3%
internal/modules/payment/infra/repo.go:30:             Save           85.7%
total:                                                 (statements)  87.4%
```

## 3.4 Coverage Gates

### GitHub Actions Gate

```yaml
- name: Coverage gate
  run: |
    go test -coverprofile=coverage.out ./...
    
    # Total coverage
    total=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    
    if (( $(echo "$total < 80" | bc -l) )); then
      echo "❌ Total coverage $total% < 80%"
      exit 1
    fi
    
    # Domain coverage (must be 100%)
    domain=$(go test -cover ./internal/modules/*/domain/... 2>&1 | \
      grep -oP 'coverage: \K[0-9.]+' | sort -n | head -1)
    
    if (( $(echo "$domain < 95" | bc -l) )); then
      echo "❌ Domain coverage $domain% < 95%"
      exit 1
    fi
    
    echo "✅ Coverage OK: total $total%, domain $domain%"
```

### Codecov Integration

```yaml
- name: Upload coverage to Codecov
  uses: codecov/codecov-action@v4
  with:
    file: ./coverage.out
    flags: unittests
    fail_ci_if_error: true
```

**codecov.yml:**
```yaml
coverage:
  status:
    project:
      default:
        target: 80%
        threshold: 2%
    patch:
      default:
        target: 80%
        threshold: 5%
```

## 3.5 Coverage Anti-Patterns

### ❌ Anti-Pattern 1: Coverage Inflation

```go
// ❌ Test ทุก getter
func TestPayment_GetID(t *testing.T) {
    p := &Payment{ID: "123"}
    assert.Equal(t, "123", p.ID)  // ← meaningless
}

func TestPayment_GetUserID(t *testing.T) {
    // ...
}
```

**แก้:** Test behavior ไม่ใช่ test field

### ❌ Anti-Pattern 2: No Assertions

```go
// ❌ Test ที่ไม่มี assertion
func TestCreate(t *testing.T) {
    p, _ := NewPayment(...)
    p.Refund()  // ← test อะไร?
}
```

**แก้:** ทุก test ต้องมี assertion

### ❌ Anti-Pattern 3: Over-Mocking

```go
// ❌ Mock ทุกอย่าง → test ไม่มีค่า
mockRepo.On("FindByID").Return(...)
mockRepo.On("Save").Return(nil)
mockRepo.On("Update").Return(nil)
mockLogger.On("Info").Return()
mockCache.On("Set").Return(nil)
// test ตรวจแค่ว่า mock ถูกเรียก → ไร้ค่า
```

**แก้:** Mock เฉพาะ external dependency

### ❌ Anti-Pattern 4: Testing Framework, Not Business

```go
// ❌ Test GORM
func TestGormWorks(t *testing.T) {
    db, _ := gorm.Open(...)
    assert.NotNil(t, db)
}
```

**แก้:** Test business logic, not framework

## 3.6 Coverage vs Quality

```
Coverage 100% ≠ Quality 100%

Example:
├── Coverage 100% (ทุกบรรทัดถูก execute)
└── Quality 30% (test ไม่มี assertion)

Coverage 80% + Quality 90% > Coverage 100% + Quality 30%
```

**เป้าหมาย:**
- Coverage: **80%+** (target)
- Quality: **Test ที่ meaningful**

---

# บทที่ 4: Unit Test สำหรับ Domain Layer

## 4.1 หลักการ

- **Test ทุก behavior method**
- **Test ทุก state transition (valid + invalid)**
- **Test ทุก invariant**
- **Test ทุก edge case**
- **Coverage 100%**

## 4.2 Test Structure

```go
package entity_test  // ← package แยก

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestEntity_Behavior_Scenario(t *testing.T) {
    // Arrange
    // ...
    
    // Act
    // ...
    
    // Assert
    // ...
}
```

## 4.3 ตัวอย่างเต็ม: Payment Entity

```go
// internal/modules/payment/domain/entity/payment_test.go
package entity_test

import (
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "icmongolang/internal/modules/payment/domain/entity"
    domainerrors "icmongolang/internal/modules/payment/domain/errors"
    valueobject "icmongolang/internal/modules/payment/domain/value_object"
)

// ============================================================
// Constructor Tests
// ============================================================

func TestNewPayment_Success(t *testing.T) {
    // Arrange
    // เตรียมข้อมูล
    userID := uuid.New()
    orderID := uuid.New()
    amount := decimal.NewFromInt(100)
    
    // Act
    // ลงมือทำ
    payment, err := entity.NewPayment(
        userID, orderID, amount, "THB",
        valueobject.PaymentMethodCard,
    )
    
    // Assert
    // ตรวจสอบ
    require.NoError(t, err)
    require.NotNil(t, payment)
    
    assert.NotEqual(t, uuid.Nil, payment.ID)
    assert.Equal(t, userID, payment.UserID)
    assert.Equal(t, orderID, payment.OrderID)
    assert.Equal(t, "100", payment.Amount.String())
    assert.Equal(t, "THB", payment.Currency)
    assert.Equal(t, valueobject.PaymentStatusPending, payment.Status)
    assert.Equal(t, valueobject.PaymentMethodCard, payment.Method)
    assert.WithinDuration(t, time.Now(), payment.CreatedAt, 1*time.Second)
    assert.WithinDuration(t, time.Now(), payment.UpdatedAt, 1*time.Second)
}

func TestNewPayment_ValidationErrors(t *testing.T) {
    tests := []struct {
        name     string
        userID   uuid.UUID
        orderID  uuid.UUID
        amount   decimal.Decimal
        currency string
        method   valueobject.PaymentMethod
        wantErr  error
    }{
        {
            name:     "nil user ID",
            userID:   uuid.Nil,
            orderID:  uuid.New(),
            amount:   decimal.NewFromInt(100),
            currency: "THB",
            method:   valueobject.PaymentMethodCard,
            wantErr:  domainerrors.ErrInvalidUserID,
        },
        {
            name:     "nil order ID",
            userID:   uuid.New(),
            orderID:  uuid.Nil,
            amount:   decimal.NewFromInt(100),
            currency: "THB",
            method:   valueobject.PaymentMethodCard,
            wantErr:  domainerrors.ErrInvalidOrderID,
        },
        {
            name:     "zero amount",
            userID:   uuid.New(),
            orderID:  uuid.New(),
            amount:   decimal.Zero,
            currency: "THB",
            method:   valueobject.PaymentMethodCard,
            wantErr:  domainerrors.ErrInvalidAmount,
        },
        {
            name:     "negative amount",
            userID:   uuid.New(),
            orderID:  uuid.New(),
            amount:   decimal.NewFromInt(-100),
            currency: "THB",
            method:   valueobject.PaymentMethodCard,
            wantErr:  domainerrors.ErrInvalidAmount,
        },
        {
            name:     "invalid currency (too short)",
            userID:   uuid.New(),
            orderID:  uuid.New(),
            amount:   decimal.NewFromInt(100),
            currency: "TH",
            method:   valueobject.PaymentMethodCard,
            wantErr:  domainerrors.ErrInvalidCurrency,
        },
        {
            name:     "invalid currency (too long)",
            userID:   uuid.New(),
            orderID:  uuid.New(),
            amount:   decimal.NewFromInt(100),
            currency: "THBB",
            method:   valueobject.PaymentMethodCard,
            wantErr:  domainerrors.ErrInvalidCurrency,
        },
        {
            name:     "invalid method",
            userID:   uuid.New(),
            orderID:  uuid.New(),
            amount:   decimal.NewFromInt(100),
            currency: "THB",
            method:   "INVALID",
            wantErr:  domainerrors.ErrInvalidPaymentMethod,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            payment, err := entity.NewPayment(
                tt.userID, tt.orderID, tt.amount, tt.currency, tt.method,
            )
            
            require.ErrorIs(t, err, tt.wantErr)
            require.Nil(t, payment)
        })
    }
}

// ============================================================
// State Transition Tests
// ============================================================

func TestPayment_StateTransitions(t *testing.T) {
    tests := []struct {
        name       string
        setup      func() *entity.Payment
        action     func(*entity.Payment) error
        wantErr    error
        wantStatus valueobject.PaymentStatus
    }{
        // Happy paths
        {
            name: "PENDING → PROCESSING",
            setup: func() *entity.Payment {
                p, _ := entity.NewPayment(uuid.New(), uuid.New(),
                    decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
                return p
            },
            action:     func(p *entity.Payment) error { return p.StartProcessing() },
            wantErr:    nil,
            wantStatus: valueobject.PaymentStatusProcessing,
        },
        {
            name: "PROCESSING → SUCCESS",
            setup: func() *entity.Payment {
                p, _ := entity.NewPayment(uuid.New(), uuid.New(),
                    decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
                _ = p.StartProcessing()
                return p
            },
            action:     func(p *entity.Payment) error { return p.MarkSuccess("txn-123") },
            wantErr:    nil,
            wantStatus: valueobject.PaymentStatusSuccess,
        },
        {
            name: "PROCESSING → FAILED",
            setup: func() *entity.Payment {
                p, _ := entity.NewPayment(uuid.New(), uuid.New(),
                    decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
                _ = p.StartProcessing()
                return p
            },
            action:     func(p *entity.Payment) error { return p.MarkFailed() },
            wantErr:    nil,
            wantStatus: valueobject.PaymentStatusFailed,
        },
        
        // Invalid transitions
        {
            name: "PENDING → SUCCESS (skip PROCESSING)",
            setup: func() *entity.Payment {
                p, _ := entity.NewPayment(uuid.New(), uuid.New(),
                    decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
                return p
            },
            action:     func(p *entity.Payment) error { return p.MarkSuccess("txn-123") },
            wantErr:    domainerrors.ErrInvalidStatusTransition,
            wantStatus: valueobject.PaymentStatusPending,  // unchanged
        },
        {
            name: "PENDING → FAILED (skip PROCESSING)",
            setup: func() *entity.Payment {
                p, _ := entity.NewPayment(uuid.New(), uuid.New(),
                    decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
                return p
            },
            action:     func(p *entity.Payment) error { return p.MarkFailed() },
            wantErr:    domainerrors.ErrInvalidStatusTransition,
            wantStatus: valueobject.PaymentStatusPending,
        },
        {
            name: "SUCCESS → PROCESSING (backward)",
            setup: func() *entity.Payment {
                p, _ := entity.NewPayment(uuid.New(), uuid.New(),
                    decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
                _ = p.StartProcessing()
                _ = p.MarkSuccess("txn-123")
                return p
            },
            action:     func(p *entity.Payment) error { return p.StartProcessing() },
            wantErr:    domainerrors.ErrInvalidStatusTransition,
            wantStatus: valueobject.PaymentStatusSuccess,
        },
        {
            name: "FAILED is terminal",
            setup: func() *entity.Payment {
                p, _ := entity.NewPayment(uuid.New(), uuid.New(),
                    decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
                _ = p.StartProcessing()
                _ = p.MarkFailed()
                return p
            },
            action:     func(p *entity.Payment) error { return p.StartProcessing() },
            wantErr:    domainerrors.ErrInvalidStatusTransition,
            wantStatus: valueobject.PaymentStatusFailed,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            p := tt.setup()
            err := tt.action(p)
            
            if tt.wantErr != nil {
                require.ErrorIs(t, err, tt.wantErr)
            } else {
                require.NoError(t, err)
            }
            assert.Equal(t, tt.wantStatus, p.Status)
        })
    }
}

// ============================================================
// Refund Tests
// ============================================================

func TestPayment_Refund_Success(t *testing.T) {
    // Arrange
    p, _ := entity.NewPayment(uuid.New(), uuid.New(),
        decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
    _ = p.StartProcessing()
    _ = p.MarkSuccess("txn-123")
    
    beforeRefund := time.Now()
    
    // Act
    err := p.Refund()
    
    // Assert
    require.NoError(t, err)
    assert.Equal(t, valueobject.PaymentStatusRefunded, p.Status)
    require.NotNil(t, p.RefundedAt)
    assert.WithinDuration(t, beforeRefund, *p.RefundedAt, 1*time.Second)
    assert.True(t, p.UpdatedAt.After(beforeRefund.Add(-1*time.Second)))
}

func TestPayment_Refund_FailsFromPending(t *testing.T) {
    p, _ := entity.NewPayment(uuid.New(), uuid.New(),
        decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
    
    err := p.Refund()
    
    assert.ErrorIs(t, err, domainerrors.ErrCannotRefund)
    assert.Equal(t, valueobject.PaymentStatusPending, p.Status)
}

func TestPayment_Refund_FailsFromProcessing(t *testing.T) {
    p, _ := entity.NewPayment(uuid.New(), uuid.New(),
        decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
    _ = p.StartProcessing()
    
    err := p.Refund()
    
    assert.ErrorIs(t, err, domainerrors.ErrCannotRefund)
    assert.Equal(t, valueobject.PaymentStatusProcessing, p.Status)
}

func TestPayment_Refund_FailsFromFailed(t *testing.T) {
    p, _ := entity.NewPayment(uuid.New(), uuid.New(),
        decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
    _ = p.StartProcessing()
    _ = p.MarkFailed()
    
    err := p.Refund()
    
    assert.ErrorIs(t, err, domainerrors.ErrCannotRefund)
}

func TestPayment_Refund_CannotRefundTwice(t *testing.T) {
    p, _ := entity.NewPayment(uuid.New(), uuid.New(),
        decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
    _ = p.StartProcessing()
    _ = p.MarkSuccess("txn-123")
    _ = p.Refund()
    
    // Second refund should fail
    err := p.Refund()
    
    assert.ErrorIs(t, err, domainerrors.ErrCannotRefund)
    assert.Equal(t, valueobject.PaymentStatusRefunded, p.Status)
}

// ============================================================
// Query Tests
// ============================================================

func TestPayment_IsRefundable(t *testing.T) {
    tests := []struct {
        name   string
        setup  func() *entity.Payment
        want   bool
    }{
        {
            name: "PENDING not refundable",
            setup: func() *entity.Payment {
                p, _ := entity.NewPayment(uuid.New(), uuid.New(),
                    decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
                return p
            },
            want: false,
        },
        {
            name: "SUCCESS refundable",
            setup: func() *entity.Payment {
                p, _ := entity.NewPayment(uuid.New(), uuid.New(),
                    decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
                _ = p.StartProcessing()
                _ = p.MarkSuccess("txn-123")
                return p
            },
            want: true,
        },
        {
            name: "REFUNDED not refundable",
            setup: func() *entity.Payment {
                p, _ := entity.NewPayment(uuid.New(), uuid.New(),
                    decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
                _ = p.StartProcessing()
                _ = p.MarkSuccess("txn-123")
                _ = p.Refund()
                return p
            },
            want: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            p := tt.setup()
            assert.Equal(t, tt.want, p.IsRefundable())
        })
    }
}

func TestPayment_IsTerminal(t *testing.T) {
    tests := []struct {
        name   string
        setup  func() *entity.Payment
        want   bool
    }{
        {
            name: "PENDING not terminal",
            setup: func() *entity.Payment {
                p, _ := entity.NewPayment(uuid.New(), uuid.New(),
                    decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
                return p
            },
            want: false,
        },
        {
            name: "PROCESSING not terminal",
            setup: func() *entity.Payment {
                p, _ := entity.NewPayment(uuid.New(), uuid.New(),
                    decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
                _ = p.StartProcessing()
                return p
            },
            want: false,
        },
        {
            name: "SUCCESS not terminal (can refund)",
            setup: func() *entity.Payment {
                p, _ := entity.NewPayment(uuid.New(), uuid.New(),
                    decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
                _ = p.StartProcessing()
                _ = p.MarkSuccess("txn-123")
                return p
            },
            want: false,
        },
        {
            name: "FAILED is terminal",
            setup: func() *entity.Payment {
                p, _ := entity.NewPayment(uuid.New(), uuid.New(),
                    decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
                _ = p.StartProcessing()
                _ = p.MarkFailed()
                return p
            },
            want: true,
        },
        {
            name: "REFUNDED is terminal",
            setup: func() *entity.Payment {
                p, _ := entity.NewPayment(uuid.New(), uuid.New(),
                    decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
                _ = p.StartProcessing()
                _ = p.MarkSuccess("txn-123")
                _ = p.Refund()
                return p
            },
            want: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            p := tt.setup()
            assert.Equal(t, tt.want, p.IsTerminal())
        })
    }
}
```

## 4.4 Test Value Object

```go
// internal/modules/payment/domain/value_object/status_test.go
package valueobject_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "icmongolang/internal/modules/payment/domain/value_object"
)

func TestPaymentStatus_IsValid(t *testing.T) {
    tests := []struct {
        status valueobject.PaymentStatus
        want   bool
    }{
        {valueobject.PaymentStatusPending, true},
        {valueobject.PaymentStatusProcessing, true},
        {valueobject.PaymentStatusSuccess, true},
        {valueobject.PaymentStatusFailed, true},
        {valueobject.PaymentStatusRefunded, true},
        {"INVALID", false},
        {"", false},
        {"pending", false},  // lowercase
    }
    
    for _, tt := range tests {
        t.Run(string(tt.status), func(t *testing.T) {
            assert.Equal(t, tt.want, tt.status.IsValid())
        })
    }
}

func TestPaymentStatus_CanTransitionTo(t *testing.T) {
    tests := []struct {
        from, to valueobject.PaymentStatus
        want     bool
    }{
        // Valid
        {valueobject.PaymentStatusPending, valueobject.PaymentStatusProcessing, true},
        {valueobject.PaymentStatusPending, valueobject.PaymentStatusFailed, true},
        {valueobject.PaymentStatusProcessing, valueobject.PaymentStatusSuccess, true},
        {valueobject.PaymentStatusProcessing, valueobject.PaymentStatusFailed, true},
        {valueobject.PaymentStatusSuccess, valueobject.PaymentStatusRefunded, true},
        
        // Invalid
        {valueobject.PaymentStatusPending, valueobject.PaymentStatusSuccess, false},
        {valueobject.PaymentStatusPending, valueobject.PaymentStatusRefunded, false},
        {valueobject.PaymentStatusProcessing, valueobject.PaymentStatusPending, false},
        {valueobject.PaymentStatusSuccess, valueobject.PaymentStatusPending, false},
        {valueobject.PaymentStatusFailed, valueobject.PaymentStatusProcessing, false},
        {valueobject.PaymentStatusRefunded, valueobject.PaymentStatusSuccess, false},
    }
    
    for _, tt := range tests {
        t.Run(string(tt.from)+"→"+string(tt.to), func(t *testing.T) {
            assert.Equal(t, tt.want, tt.from.CanTransitionTo(tt.to))
        })
    }
}
```

## 4.5 Run Tests

```bash
# Run domain tests
go test ./internal/modules/payment/domain/... -v

# With coverage
go test -cover ./internal/modules/payment/domain/... 

# Detailed coverage
go test -coverprofile=domain.out ./internal/modules/payment/domain/...
go tool cover -func=domain.out
go tool cover -html=domain.out
```

**ผลลัพธ์:**
```
PASS
ok   icmongolang/internal/modules/payment/domain/entity  0.045s
coverage: 100.0% of statements
```

---

# บทที่ 5: Unit Test สำหรับ Application Layer

## 5.1 หลักการ

- **Test use case logic ทั้งหมด**
- **Mock repository**
- **Mock external services** (email, kafka)
- **Coverage 90%+**

## 5.2 ตัวอย่างเต็ม: CreatePaymentUseCase

```go
// internal/modules/payment/application/create_payment_test.go
package application_test

import (
    "context"
    "errors"
    "testing"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"

    "icmongolang/internal/modules/payment/application"
    "icmongolang/internal/modules/payment/application/mocks"
    "icmongolang/internal/modules/payment/domain/entity"
    domainerrors "icmongolang/internal/modules/payment/domain/errors"
    "icmongolang/pkg/logger"
)

// ============================================================
// Happy Path
// ============================================================

func TestCreatePaymentUseCase_Success(t *testing.T) {
    // Arrange
    repo := new(mocks.PaymentRepositoryMock)
    log := logger.NewNoop()
    uc := application.NewCreatePaymentUseCase(repo, log)
    
    userID := uuid.New()
    orderID := uuid.New()
    
    // Mock: order ไม่ซ้ำ
    repo.On("FindByOrderID", mock.Anything, orderID).
        Return(nil, domainerrors.ErrPaymentNotFound)
    // Mock: save สำเร็จ
    repo.On("Save", mock.Anything, mock.AnythingOfType("*entity.Payment")).
        Return(nil)
    
    input := application.CreatePaymentInput{
        UserID:   userID,
        OrderID:  orderID,
        Amount:   decimal.NewFromInt(100),
        Currency: "THB",
        Method:   "CARD",
    }
    
    // Act
    output, err := uc.Execute(context.Background(), input)
    
    // Assert
    require.NoError(t, err)
    require.NotNil(t, output)
    
    assert.NotEqual(t, uuid.Nil, output.ID)
    assert.Equal(t, "PENDING", output.Status)
    assert.Equal(t, "100", output.Amount)
    assert.Equal(t, "THB", output.Currency)
    
    repo.AssertExpectations(t)
}

// ============================================================
// Duplicate Detection
// ============================================================

func TestCreatePaymentUseCase_DuplicateOrder(t *testing.T) {
    repo := new(mocks.PaymentRepositoryMock)
    log := logger.NewNoop()
    uc := application.NewCreatePaymentUseCase(repo, log)
    
    userID := uuid.New()
    orderID := uuid.New()
    
    // Existing payment
    existing, _ := entity.NewPayment(
        userID, orderID, decimal.NewFromInt(50), "THB", "CARD",
    )
    repo.On("FindByOrderID", mock.Anything, orderID).
        Return(existing, nil)
    
    input := application.CreatePaymentInput{
        UserID:   userID,
        OrderID:  orderID,
        Amount:   decimal.NewFromInt(100),
        Currency: "THB",
        Method:   "CARD",
    }
    
    output, err := uc.Execute(context.Background(), input)
    
    assert.ErrorIs(t, err, domainerrors.ErrPaymentAlreadyExists)
    assert.Nil(t, output)
    repo.AssertExpectations(t)
    // Save ไม่ควรถูกเรียก
    repo.AssertNotCalled(t, "Save")
}

// ============================================================
// Validation Errors
// ============================================================

func TestCreatePaymentUseCase_ValidationErrors(t *testing.T) {
    tests := []struct {
        name    string
        input   application.CreatePaymentInput
        wantErr error
    }{
        {
            name: "zero amount",
            input: application.CreatePaymentInput{
                UserID:   uuid.New(),
                OrderID:  uuid.New(),
                Amount:   decimal.Zero,
                Currency: "THB",
                Method:   "CARD",
            },
            wantErr: domainerrors.ErrInvalidAmount,
        },
        {
            name: "invalid method",
            input: application.CreatePaymentInput{
                UserID:   uuid.New(),
                OrderID:  uuid.New(),
                Amount:   decimal.NewFromInt(100),
                Currency: "THB",
                Method:   "INVALID",
            },
            wantErr: domainerrors.ErrInvalidPaymentMethod,
        },
        {
            name: "nil user ID",
            input: application.CreatePaymentInput{
                UserID:   uuid.Nil,
                OrderID:  uuid.New(),
                Amount:   decimal.NewFromInt(100),
                Currency: "THB",
                Method:   "CARD",
            },
            wantErr: domainerrors.ErrInvalidUserID,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := new(mocks.PaymentRepositoryMock)
            log := logger.NewNoop()
            uc := application.NewCreatePaymentUseCase(repo, log)
            
            // Setup mock (FindByOrderID จะถูกเรียกก่อน)
            repo.On("FindByOrderID", mock.Anything, tt.input.OrderID).
                Return(nil, domainerrors.ErrPaymentNotFound)
            
            output, err := uc.Execute(context.Background(), tt.input)
            
            assert.ErrorIs(t, err, tt.wantErr)
            assert.Nil(t, output)
            repo.AssertNotCalled(t, "Save")
        })
    }
}

// ============================================================
// Repository Errors
// ============================================================

func TestCreatePaymentUseCase_RepositoryError(t *testing.T) {
    repo := new(mocks.PaymentRepositoryMock)
    log := logger.NewNoop()
    uc := application.NewCreatePaymentUseCase(repo, log)
    
    orderID := uuid.New()
    dbErr := errors.New("database connection failed")
    
    repo.On("FindByOrderID", mock.Anything, orderID).
        Return(nil, domainerrors.ErrPaymentNotFound)
    repo.On("Save", mock.Anything, mock.Anything).
        Return(dbErr)
    
    output, err := uc.Execute(context.Background(), application.CreatePaymentInput{
        UserID:   uuid.New(),
        OrderID:  orderID,
        Amount:   decimal.NewFromInt(100),
        Currency: "THB",
        Method:   "CARD",
    })
    
    assert.ErrorIs(t, err, dbErr)
    assert.Nil(t, output)
    repo.AssertExpectations(t)
}

// ============================================================
// FindByOrderID Errors (non-notfound)
// ============================================================

func TestCreatePaymentUseCase_FindByOrderIDError(t *testing.T) {
    repo := new(mocks.PaymentRepositoryMock)
    log := logger.NewNoop()
    uc := application.NewCreatePaymentUseCase(repo, log)
    
    orderID := uuid.New()
    dbErr := errors.New("connection timeout")
    
    // Simulate DB error (ไม่ใช่ ErrPaymentNotFound)
    repo.On("FindByOrderID", mock.Anything, orderID).Return(nil, dbErr)
    repo.On("Save", mock.Anything, mock.Anything).Return(nil)
    
    // Current implementation: ignore non-notfound errors
    // (ใน use case: existing, _ := uc.repo.FindByOrderID(...))
    // → existing จะเป็น nil, err จะถูก ignore → save จะถูกเรียก
    
    output, err := uc.Execute(context.Background(), application.CreatePaymentInput{
        UserID:   uuid.New(),
        OrderID:  orderID,
        Amount:   decimal.NewFromInt(100),
        Currency: "THB",
        Method:   "CARD",
    })
    
    // ถ้า use case ignore error → save ถูกเรียก
    // ถ้า use case จัดการ error → ควร return
    // (ขึ้นกับ design)
    _ = output
    _ = err
}

// ============================================================
// Integration with domain errors
// ============================================================

func TestCreatePaymentUseCase_DomainValidationFailsAfterDuplicateCheck(t *testing.T) {
    repo := new(mocks.PaymentRepositoryMock)
    log := logger.NewNoop()
    uc := application.NewCreatePaymentUseCase(repo, log)
    
    orderID := uuid.New()
    
    // FindByOrderID ผ่าน (ไม่มี duplicate)
    repo.On("FindByOrderID", mock.Anything, orderID).
        Return(nil, domainerrors.ErrPaymentNotFound)
    
    // แต่ domain validation fail (amount < 0)
    output, err := uc.Execute(context.Background(), application.CreatePaymentInput{
        UserID:   uuid.New(),
        OrderID:  orderID,
        Amount:   decimal.NewFromInt(-100),  // ← invalid
        Currency: "THB",
        Method:   "CARD",
    })
    
    assert.ErrorIs(t, err, domainerrors.ErrInvalidAmount)
    assert.Nil(t, output)
    repo.AssertNotCalled(t, "Save")
}
```

## 5.3 Test RefundPaymentUseCase

```go
// internal/modules/payment/application/refund_payment_test.go
package application_test

import (
    "context"
    "testing"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"

    "icmongolang/internal/modules/payment/application"
    "icmongolang/internal/modules/payment/application/mocks"
    "icmongolang/internal/modules/payment/domain/entity"
    domainerrors "icmongolang/internal/modules/payment/domain/errors"
    "icmongolang/pkg/logger"
)

func TestRefundPaymentUseCase_Success(t *testing.T) {
    // Arrange
    repo := new(mocks.PaymentRepositoryMock)
    log := logger.NewNoop()
    uc := application.NewRefundPaymentUseCase(repo, log)
    
    userID := uuid.New()
    paymentID := uuid.New()
    
    // Create a SUCCESS payment
    payment, _ := entity.NewPayment(
        userID, uuid.New(), decimal.NewFromInt(100), "THB", "CARD",
    )
    payment.ID = paymentID
    _ = payment.StartProcessing()
    _ = payment.MarkSuccess("txn-123")
    
    repo.On("FindByID", mock.Anything, paymentID).Return(payment, nil)
    repo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Payment")).Return(nil)
    
    // Act
    err := uc.Execute(context.Background(), application.RefundPaymentInput{
        PaymentID: paymentID,
        UserID:    userID,
        Reason:    "customer request",
    })
    
    // Assert
    require.NoError(t, err)
    assert.Equal(t, "REFUNDED", string(payment.Status))
    assert.NotNil(t, payment.RefundedAt)
    repo.AssertExpectations(t)
}

func TestRefundPaymentUseCase_Unauthorized(t *testing.T) {
    repo := new(mocks.PaymentRepositoryMock)
    log := logger.NewNoop()
    uc := application.NewRefundPaymentUseCase(repo, log)
    
    ownerID := uuid.New()
    otherID := uuid.New()
    paymentID := uuid.New()
    
    payment, _ := entity.NewPayment(
        ownerID, uuid.New(), decimal.NewFromInt(100), "THB", "CARD",
    )
    payment.ID = paymentID
    _ = payment.StartProcessing()
    _ = payment.MarkSuccess("txn-123")
    
    repo.On("FindByID", mock.Anything, paymentID).Return(payment, nil)
    
    err := uc.Execute(context.Background(), application.RefundPaymentInput{
        PaymentID: paymentID,
        UserID:    otherID,  // ← not owner
        Reason:    "unauthorized attempt",
    })
    
    assert.ErrorIs(t, err, domainerrors.ErrUnauthorized)
    repo.AssertNotCalled(t, "Update")
}

func TestRefundPaymentUseCase_NotFound(t *testing.T) {
    repo := new(mocks.PaymentRepositoryMock)
    log := logger.NewNoop()
    uc := application.NewRefundPaymentUseCase(repo, log)
    
    paymentID := uuid.New()
    repo.On("FindByID", mock.Anything, paymentID).
        Return(nil, domainerrors.ErrPaymentNotFound)
    
    err := uc.Execute(context.Background(), application.RefundPaymentInput{
        PaymentID: paymentID,
        UserID:    uuid.New(),
        Reason:    "test",
    })
    
    assert.ErrorIs(t, err, domainerrors.ErrPaymentNotFound)
}

func TestRefundPaymentUseCase_CannotRefundPending(t *testing.T) {
    repo := new(mocks.PaymentRepositoryMock)
    log := logger.NewNoop()
    uc := application.NewRefundPaymentUseCase(repo, log)
    
    userID := uuid.New()
    paymentID := uuid.New()
    
    // PENDING payment (not refundable)
    payment, _ := entity.NewPayment(
        userID, uuid.New(), decimal.NewFromInt(100), "THB", "CARD",
    )
    payment.ID = paymentID
    
    repo.On("FindByID", mock.Anything, paymentID).Return(payment, nil)
    
    err := uc.Execute(context.Background(), application.RefundPaymentInput{
        PaymentID: paymentID,
        UserID:    userID,
        Reason:    "test",
    })
    
    assert.ErrorIs(t, err, domainerrors.ErrCannotRefund)
    repo.AssertNotCalled(t, "Update")
}
```

## 5.4 Run Tests

```bash
# Run application tests
go test ./internal/modules/payment/application/... -v

# With coverage
go test -cover ./internal/modules/payment/application/...
```

**ผลลัพธ์:**
```
=== RUN   TestCreatePaymentUseCase_Success
--- PASS: TestCreatePaymentUseCase_Success (0.00s)
=== RUN   TestCreatePaymentUseCase_DuplicateOrder
--- PASS: TestCreatePaymentUseCase_DuplicateOrder (0.00s)
...
PASS
ok   icmongolang/internal/modules/payment/application  0.052s
coverage: 94.7% of statements
```

---

# บทที่ 6: Mocking และ Test Doubles

## 6.1 ประเภทของ Test Doubles

| Type | คำอธิบาย | ตัวอย่าง |
|---|---|---|
| **Dummy** | Object เปล่า ส่งไปเพื่อให้ compile | `&struct{}{}` |
| **Stub** | Return ค่าคงที่ | `Return(nil, nil)` |
| **Spy** | บันทึกการเรียก | `mock.Mock` |
| **Mock** | กำหนด expectation ล่วงหน้า | `testify/mock` |
| **Fake** | Implementation จริง แต่ simpler | In-memory DB |

## 6.2 Mocking ด้วย testify

### 6.2.1 สร้าง Mock

```go
// internal/modules/payment/application/mocks/payment_repository_mock.go
package mocks

import (
    "context"
    "github.com/google/uuid"
    "github.com/stretchr/testify/mock"
    "icmongolang/internal/modules/payment/domain/entity"
)

// PaymentRepositoryMock is a mock of PaymentRepository
// PaymentRepositoryMock เป็น mock ของ PaymentRepository
type PaymentRepositoryMock struct {
    mock.Mock
}

func (m *PaymentRepositoryMock) Save(ctx context.Context, p *entity.Payment) error {
    args := m.Called(ctx, p)
    return args.Error(0)
}

func (m *PaymentRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.Payment), args.Error(1)
}

func (m *PaymentRepositoryMock) FindByOrderID(ctx context.Context, orderID uuid.UUID) (*entity.Payment, error) {
    args := m.Called(ctx, orderID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.Payment), args.Error(1)
}

func (m *PaymentRepositoryMock) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Payment, error) {
    args := m.Called(ctx, userID, limit, offset)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]*entity.Payment), args.Error(1)
}

func (m *PaymentRepositoryMock) Update(ctx context.Context, p *entity.Payment) error {
    args := m.Called(ctx, p)
    return args.Error(0)
}

func (m *PaymentRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}
```

### 6.2.2 ใช้ Mock

```go
func TestExample(t *testing.T) {
    repo := new(mocks.PaymentRepositoryMock)
    
    // กำหนด expectation
    repo.On("FindByID", mock.Anything, mock.Anything).
        Return(payment, nil).
        Once()  // เรียกครั้งเดียว
    
    repo.On("Save", mock.Anything, mock.Anything).
        Return(nil).
        Maybe()  // อาจจะเรียกหรือไม่ก็ได้
    
    // Test...
    
    // ตรวจสอบ
    repo.AssertExpectations(t)
    repo.AssertCalled(t, "Save", mock.Anything, payment)
    repo.AssertNotCalled(t, "Delete")
    repo.AssertNumberOfCalls(t, "FindByID", 1)
}
```

### 6.2.3 Mock with Arguments

```go
// Match specific arguments
repo.On("FindByID", mock.Anything, uuid.MustParse("...")).Return(payment, nil)

// Match any
repo.On("Save", mock.Anything, mock.Anything).Return(nil)

// Match type
repo.On("Save", mock.Anything, mock.AnythingOfType("*entity.Payment")).Return(nil)

// Custom matcher
repo.On("Save", mock.Anything, mock.MatchedBy(func(p *entity.Payment) bool {
    return p.Amount.GreaterThan(decimal.NewFromInt(100))
})).Return(nil)

// Dynamic return
repo.On("FindByID", mock.Anything, mock.Anything).Return(
    func(ctx context.Context, id uuid.UUID) *entity.Payment {
        return &entity.Payment{ID: id}
    },
    nil,
)
```

## 6.3 Mock Generation Tools

### 6.3.1 mockery

```bash
# Install
go install github.com/vektra/mockery/v2@latest

# Generate all
mockery --all --output=./mocks --outpkg=mocks

# Config
cat > .mockery.yaml <<EOF
with-expecter: true
dir: "internal/modules/payment"
output: "mocks"
packages:
  icmongolang/internal/modules/payment/domain/repository:
    interfaces:
      PaymentRepository:
EOF

mockery
```

### 6.3.2 gomock

```go
//go:generate mockgen -source=payment_repository.go -destination=mocks/payment_repository_mock.go -package=mocks

// Usage
ctrl := gomock.NewController(t)
defer ctrl.Finish()

mockRepo := mocks.NewMockPaymentRepository(ctrl)
mockRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(payment, nil)
mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
```

## 6.4 Fakes (สำหรับ Integration Tests)

```go
// internal/modules/payment/testutils/fakes/fake_repo.go
package fakes

// FakePaymentRepo is an in-memory implementation for testing
// FakePaymentRepo เป็น in-memory implementation สำหรับ test
type FakePaymentRepo struct {
    payments map[uuid.UUID]*entity.Payment
    mu       sync.RWMutex
}

func NewFakePaymentRepo() *FakePaymentRepo {
    return &FakePaymentRepo{
        payments: make(map[uuid.UUID]*entity.Payment),
    }
}

func (r *FakePaymentRepo) Save(ctx context.Context, p *entity.Payment) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    // Check duplicate
    for _, existing := range r.payments {
        if existing.OrderID == p.OrderID {
            return domainerrors.ErrPaymentAlreadyExists
        }
    }
    
    // Copy to prevent mutation
    copy := *p
    r.payments[p.ID] = &copy
    return nil
}

func (r *FakePaymentRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    p, ok := r.payments[id]
    if !ok {
        return nil, domainerrors.ErrPaymentNotFound
    }
    
    copy := *p
    return &copy, nil
}

// ... other methods
```

**ใช้:**
```go
func TestWithFakeRepo(t *testing.T) {
    repo := fakes.NewFakePaymentRepo()
    uc := application.NewCreatePaymentUseCase(repo, logger.NewNoop())
    
    // Test without mock setup
    output, err := uc.Execute(ctx, input)
    require.NoError(t, err)
    
    // Verify via fake
    found, _ := repo.FindByID(ctx, output.ID)
    assert.NotNil(t, found)
}
```

## 6.5 หลักการ Mocking

### ✅ ควร Mock:
- External API (Stripe, SMTP)
- Database (สำหรับ unit test)
- Message queue (Kafka, RabbitMQ)
- Clock / Time
- Random generator

### ❌ ไม่ควร Mock:
- Domain entities
- Value objects
- Pure functions
- Standard library

### ⚠️ ระวัง:
- **Over-mocking** → test ไม่มีค่า
- **Mock return ผิด** → test ผ่านแต่ bug จริง
- **Mock ไม่ maintain** → flaky test

---

# บทที่ 7: Table-Driven Tests

## 7.1 หลักการ

**Table-Driven Test** = Test ที่ใช้ table ของ test cases

**ประโยชน์:**
- เพิ่ม case ง่าย (แค่เพิ่ม row)
- อ่านง่าย (เห็นทุก case ในที่เดียว)
- Maintain ง่าย
- Standard ใน Go community

## 7.2 Basic Structure

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        want     OutputType
        wantErr  error
    }{
        {
            name:  "case 1",
            input: ...,
            want:  ...,
        },
        // ... more cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            
            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
                return
            }
            require.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

## 7.3 ตัวอย่าง: Validate Amount

```go
func TestValidateAmount(t *testing.T) {
    tests := []struct {
        name    string
        amount  decimal.Decimal
        wantErr error
    }{
        // Valid
        {"minimum valid", decimal.NewFromInt(1), nil},
        {"typical", decimal.NewFromInt(100), nil},
        {"large", decimal.NewFromInt(999999), nil},
        {"decimal", decimal.NewFromFloat(99.99), nil},
        
        // Invalid
        {"zero", decimal.Zero, ErrInvalidAmount},
        {"negative", decimal.NewFromInt(-1), ErrInvalidAmount},
        {"negative large", decimal.NewFromInt(-999999), ErrInvalidAmount},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateAmount(tt.amount)
            
            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## 7.4 Advanced: Setup/Teardown per Case

```go
func TestComplexCases(t *testing.T) {
    tests := []struct {
        name    string
        setup   func(t *testing.T) (*Dependencies, func())
        input   Input
        verify  func(t *testing.T, deps *Dependencies, out *Output)
    }{
        {
            name: "case with custom setup",
            setup: func(t *testing.T) (*Dependencies, func()) {
                db := setupTestDB(t)
                return &Dependencies{DB: db}, func() { db.Close() }
            },
            input: Input{...},
            verify: func(t *testing.T, deps *Dependencies, out *Output) {
                assert.Equal(t, "expected", out.Value)
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            deps, cleanup := tt.setup(t)
            defer cleanup()
            
            out, err := Run(deps, tt.input)
            require.NoError(t, err)
            
            tt.verify(t, deps, out)
        })
    }
}
```

## 7.5 Parallel Tests

```go
func TestParallel(t *testing.T) {
    tests := []struct {
        name  string
        input int
        want  int
    }{
        {"1", 1, 2},
        {"2", 2, 4},
        {"3", 3, 6},
    }
    
    for _, tt := range tests {
        tt := tt  // ⚠️ capture loop variable (ก่อน Go 1.22)
        
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()  // ← รันขนาน
            
            got := double(tt.input)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

**⚠️ ระวัง:**
- Shared state
- Race condition
- Order-dependent tests

## 7.6 Subtests

```go
func TestPayment(t *testing.T) {
    t.Run("Lifecycle", func(t *testing.T) {
        t.Run("happy path", func(t *testing.T) {
            // ...
        })
        t.Run("invalid transition", func(t *testing.T) {
            // ...
        })
    })
    
    t.Run("Refund", func(t *testing.T) {
        t.Run("success", func(t *testing.T) {
            // ...
        })
        t.Run("cannot refund twice", func(t *testing.T) {
            // ...
        })
    })
}
```

## 7.7 Golden Files

```go
func TestRenderResponse(t *testing.T) {
    input := ...
    output := render(input)
    
    goldenFile := "testdata/response.golden"
    
    if *update {
        // -update flag → update golden file
        os.WriteFile(goldenFile, []byte(output), 0644)
        return
    }
    
    expected, _ := os.ReadFile(goldenFile)
    assert.Equal(t, string(expected), output)
}
```

**ใช้:**
```bash
# Run normally
go test ./...

# Update golden files
go test -update ./...
```

---

# บทที่ 8: Fuzzing และ Benchmark

## 8.1 Fuzzing

**Fuzzing** = ทดสอบด้วย random input เพื่อหาบั๊ก

### 8.1.1 Basic Fuzz Test

```go
func FuzzNewPayment(f *testing.F) {
    // Seed corpus
    f.Add("user-id", "order-id", "100.00", "THB", "CARD")
    f.Add("", "", "0", "", "")
    
    f.Fuzz(func(t *testing.T, userID, orderID, amount, currency, method string) {
        userUUID, err := uuid.Parse(userID)
        if err != nil {
            return  // skip invalid
        }
        orderUUID, err := uuid.Parse(orderID)
        if err != nil {
            return
        }
        amt, err := decimal.NewFromString(amount)
        if err != nil {
            return
        }
        
        // ต้องไม่ panic
        _, _ = entity.NewPayment(
            userUUID, orderUUID, amt, currency,
            valueobject.PaymentMethod(method),
        )
    })
}
```

### 8.1.2 Run Fuzz Test

```bash
# Run for 30s
go test -fuzz=FuzzNewPayment -fuzztime=30s ./...

# Run indefinitely
go test -fuzz=FuzzNewPayment ./...

# Run found corpus
go test -run=FuzzNewPayment ./...
```

### 8.1.3 Fuzz Target ที่ดี

```go
// Fuzz JSON parsing
func FuzzParseWebhook(f *testing.F) {
    f.Add([]byte(`{"event":"payment.succeeded"}`))
    
    f.Fuzz(func(t *testing.T, data []byte) {
        // ต้องไม่ panic
        var event WebhookEvent
        _ = json.Unmarshal(data, &event)
    })
}

// Fuzz string parsing
func FuzzParsePhone(f *testing.F) {
    f.Add("+66812345678")
    
    f.Fuzz(func(t *testing.T, phone string) {
        _ = normalizePhone(phone)
    })
}
```

## 8.2 Benchmark

### 8.2.1 Basic Benchmark

```go
func BenchmarkNewPayment(b *testing.B) {
    userID := uuid.New()
    orderID := uuid.New()
    amount := decimal.NewFromInt(100)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = entity.NewPayment(
            userID, orderID, amount, "THB", valueobject.PaymentMethodCard,
        )
    }
}
```

### 8.2.2 Benchmark Table

```go
func BenchmarkCreatePayment(b *testing.B) {
    scenarios := []struct {
        name   string
        amount decimal.Decimal
    }{
        {"small", decimal.NewFromInt(10)},
        {"medium", decimal.NewFromInt(1000)},
        {"large", decimal.NewFromInt(1000000)},
    }
    
    for _, s := range scenarios {
        b.Run(s.name, func(b *testing.B) {
            userID := uuid.New()
            orderID := uuid.New()
            
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                _, _ = entity.NewPayment(
                    userID, orderID, s.amount, "THB",
                    valueobject.PaymentMethodCard,
                )
            }
        })
    }
}
```

### 8.2.3 Run Benchmark

```bash
# Run all
go test -bench=. ./...

# Run specific
go test -bench=BenchmarkNewPayment ./...

# With memory
go test -bench=. -benchmem ./...

# Multiple runs (stable result)
go test -bench=. -count=10 ./...

# Compare with benchstat
go install golang.org/x/perf/cmd/benchstat@latest
go test -bench=. -count=10 ./... > new.txt
benchstat old.txt new.txt
```

**ผลลัพธ์:**
```
BenchmarkNewPayment-16    5000000    245 ns/op    128 B/op    3 allocs/op
```

### 8.2.4 Benchmark Best Practices

```go
// ✅ Reset timer หลัง setup
func BenchmarkX(b *testing.B) {
    data := setupExpensiveData()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        use(data)
    }
}

// ✅ ใช้ b.ReportAllocs()
func BenchmarkY(b *testing.B) {
    b.ReportAllocs()
    // ...
}

// ✅ ใช้ b.RunParallel() สำหรับ concurrent
func BenchmarkZ(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            doWork()
        }
    })
}

// ❌ อย่าใช้ b.N นอก loop
func BenchmarkBad(b *testing.B) {
    data := make([]int, b.N)  // ← ผิด
    // ...
}
```

## 8.3 Coverage + Fuzz + Benchmark ใน CI

```yaml
# .github/workflows/quality.yml
name: Quality

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.22' }
      
      # Unit tests
      - name: Unit tests
        run: go test -race -coverprofile=coverage.out ./...
      
      # Fuzz (30s)
      - name: Fuzz tests
        run: |
          for fuzz in $(grep -r "^func Fuzz" --include="*_test.go" . | \
                        sed 's/.*func \(Fuzz[A-Za-z]*\).*/\1/'); do
            go test -fuzz=$fuzz -fuzztime=10s ./...
          done
      
      # Benchmark (compare)
      - name: Benchmark
        run: go test -bench=. -benchmem -count=3 ./... | tee bench.txt
      
      # Coverage gate
      - name: Coverage gate
        run: |
          cov=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          if (( $(echo "$cov < 80" | bc -l) )); then exit 1; fi
```

---

# บทที่ 9: Integration Test ด้วย Testcontainers

## 9.1 หลักการ

**Integration Test** = Test ที่ใช้ infrastructure จริง (DB, Kafka, Redis)

**Testcontainers** = Library ที่ spin up Docker containers ใน test

**ประโยชน์:**
- Test กับ DB จริง (ไม่ใช่ mock)
- Test schema จริง
- Test constraints จริง
- Deterministic (containers ถูก reset ทุกครั้ง)

## 9.2 Setup Testcontainers

```bash
go get github.com/testcontainers/testcontainers-go
go get github.com/testcontainers/testcontainers-go/modules/postgres
```

## 9.3 Postgres Container

```go
//go:build integration

package postgres_test

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/testcontainers/testcontainers-go/wait"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

type TestDB struct {
    Container *postgres.PostgresContainer
    DB        *gorm.DB
    DSN       string
}

// SetupPostgres starts a Postgres container for testing
// SetupPostgres เริ่ม Postgres container สำหรับ test
func SetupPostgres(t *testing.T) *TestDB {
    t.Helper()
    
    ctx := context.Background()
    
    // 1. Start container
    container, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15-alpine"),
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(30*time.Second),
        ),
    )
    require.NoError(t, err)
    
    // 2. Get DSN
    dsn, err := container.ConnectionString(ctx, "sslmode=disable")
    require.NoError(t, err)
    
    // 3. Connect via GORM
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    require.NoError(t, err)
    
    // 4. Migrate
    require.NoError(t, db.AutoMigrate(&PaymentModel{}))
    
    // 5. Register cleanup
    t.Cleanup(func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()
        container.Terminate(ctx)
    })
    
    return &TestDB{
        Container: container,
        DB:        db,
        DSN:       dsn,
    }
}
```

## 9.4 Repository Integration Test

```go
//go:build integration

package postgres_test

import (
    "context"
    "testing"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "icmongolang/internal/modules/payment/domain/entity"
    domainerrors "icmongolang/internal/modules/payment/domain/errors"
    "icmongolang/internal/modules/payment/infrastructure/persistence/postgres"
)

func TestPaymentRepository_SaveAndFind(t *testing.T) {
    // Setup
    testDB := SetupPostgres(t)
    repo := postgres.NewPaymentRepository(testDB.DB)
    ctx := context.Background()
    
    // Create payment
    payment, err := entity.NewPayment(
        uuid.New(), uuid.New(),
        decimal.NewFromInt(100), "THB", "CARD",
    )
    require.NoError(t, err)
    
    // Save
    require.NoError(t, repo.Save(ctx, payment))
    
    // Find
    found, err := repo.FindByID(ctx, payment.ID)
    require.NoError(t, err)
    assert.Equal(t, payment.ID, found.ID)
    assert.Equal(t, payment.UserID, found.UserID)
    assert.Equal(t, payment.Amount.String(), found.Amount.String())
    assert.Equal(t, payment.Status, found.Status)
}

func TestPaymentRepository_DuplicateOrder(t *testing.T) {
    testDB := SetupPostgres(t)
    repo := postgres.NewPaymentRepository(testDB.DB)
    ctx := context.Background()
    
    orderID := uuid.New()
    
    // First payment
    p1, _ := entity.NewPayment(uuid.New(), orderID, decimal.NewFromInt(100), "THB", "CARD")
    require.NoError(t, repo.Save(ctx, p1))
    
    // Duplicate should fail
    p2, _ := entity.NewPayment(uuid.New(), orderID, decimal.NewFromInt(50), "THB", "CARD")
    err := repo.Save(ctx, p2)
    
    assert.ErrorIs(t, err, domainerrors.ErrPaymentAlreadyExists)
}

func TestPaymentRepository_Update(t *testing.T) {
    testDB := SetupPostgres(t)
    repo := postgres.NewPaymentRepository(testDB.DB)
    ctx := context.Background()
    
    // Create + Save
    payment, _ := entity.NewPayment(
        uuid.New(), uuid.New(),
        decimal.NewFromInt(100), "THB", "CARD",
    )
    require.NoError(t, repo.Save(ctx, payment))
    
    // Update
    _ = payment.StartProcessing()
    _ = payment.MarkSuccess("txn-123")
    require.NoError(t, repo.Update(ctx, payment))
    
    // Verify
    found, err := repo.FindByID(ctx, payment.ID)
    require.NoError(t, err)
    assert.Equal(t, "SUCCESS", string(found.Status))
    assert.Equal(t, "txn-123", found.ProviderTxnID)
}

func TestPaymentRepository_NotFound(t *testing.T) {
    testDB := SetupPostgres(t)
    repo := postgres.NewPaymentRepository(testDB.DB)
    ctx := context.Background()
    
    _, err := repo.FindByID(ctx, uuid.New())
    
    assert.ErrorIs(t, err, domainerrors.ErrPaymentNotFound)
}

func TestPaymentRepository_FindByUserID(t *testing.T) {
    testDB := SetupPostgres(t)
    repo := postgres.NewPaymentRepository(testDB.DB)
    ctx := context.Background()
    
    userID := uuid.New()
    
    // Create 5 payments
    for i := 0; i < 5; i++ {
        p, _ := entity.NewPayment(userID, uuid.New(), decimal.NewFromInt(int64(100+i)), "THB", "CARD")
        require.NoError(t, repo.Save(ctx, p))
    }
    
    // Find with pagination
    payments, err := repo.FindByUserID(ctx, userID, 3, 0)
    require.NoError(t, err)
    assert.Len(t, payments, 3)
    
    // Page 2
    payments, err = repo.FindByUserID(ctx, userID, 3, 3)
    require.NoError(t, err)
    assert.Len(t, payments, 2)
}
```

## 9.5 Transaction Test

```go
func TestPaymentRepository_Transaction(t *testing.T) {
    testDB := SetupPostgres(t)
    ctx := context.Background()
    
    // Start transaction
    tx := testDB.DB.Begin()
    defer tx.Rollback()
    
    repo := postgres.NewPaymentRepository(tx)
    
    payment, _ := entity.NewPayment(
        uuid.New(), uuid.New(),
        decimal.NewFromInt(100), "THB", "CARD",
    )
    require.NoError(t, repo.Save(ctx, payment))
    
    // Verify within transaction
    found, err := repo.FindByID(ctx, payment.ID)
    require.NoError(t, err)
    assert.NotNil(t, found)
    
    // Rollback
    tx.Rollback()
    
    // Verify after rollback (should not exist)
    _, err = repo.FindByID(ctx, payment.ID)
    assert.ErrorIs(t, err, domainerrors.ErrPaymentNotFound)
}
```

## 9.6 Run Integration Tests

```bash
# Run with tag
go test -tags=integration ./... -v

# Skip if Docker not available
go test -tags=integration -short ./...
```

**ใน code:**
```go
func TestX(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    // ...
}
```

## 9.7 Testcontainers Performance

| Container | Startup Time |
|---|---|
| Postgres | 3-5s |
| MySQL | 5-8s |
| Redis | 1-2s |
| Kafka | 10-15s |
| Elasticsearch | 15-20s |

**Optimization:**
```go
// Reuse container across tests
var sharedDB *TestDB

func TestMain(m *testing.M) {
    sharedDB = setupSharedPostgres()
    defer sharedDB.Container.Terminate(context.Background())
    
    code := m.Run()
    os.Exit(code)
}
```

---

# บทที่ 10: Database Integration Test

## 10.1 Schema Migration Test

```go
//go:build integration

func TestMigrations(t *testing.T) {
    testDB := SetupPostgres(t)
    
    // Read migration file
    migrationSQL, err := os.ReadFile("../../migrations/20260401_payment_init.sql")
    require.NoError(t, err)
    
    // Apply migration
    err = testDB.DB.Exec(string(migrationSQL)).Error
    require.NoError(t, err)
    
    // Verify table exists
    var exists bool
    err = testDB.DB.Raw(`
        SELECT EXISTS (
            SELECT 1 FROM information_schema.tables
            WHERE table_name = 'payment_transactions'
        )
    `).Scan(&exists).Error
    require.NoError(t, err)
    assert.True(t, exists)
    
    // Verify index
    var indexExists bool
    err = testDB.DB.Raw(`
        SELECT EXISTS (
            SELECT 1 FROM pg_indexes
            WHERE indexname = 'uq_payment_transactions_order_id'
        )
    `).Scan(&indexExists).Error
    require.NoError(t, err)
    assert.True(t, indexExists)
}
```

## 10.2 Constraint Test

```go
func TestConstraints(t *testing.T) {
    testDB := SetupPostgres(t)
    
    // Test NOT NULL constraint
    err := testDB.DB.Exec(`
        INSERT INTO payment_transactions (amount) VALUES (100)
    `).Error
    assert.Error(t, err)  // should fail (missing user_id)
    
    // Test CHECK constraint
    err = testDB.DB.Exec(`
        INSERT INTO payment_transactions (user_id, order_id, amount, method, status)
        VALUES (gen_random_uuid(), gen_random_uuid(), -100, 'CARD', 'PENDING')
    `).Error
    assert.Error(t, err)  // should fail (amount > 0)
    
    // Test UNIQUE constraint
    orderID := uuid.New()
    testDB.DB.Exec(`
        INSERT INTO payment_transactions (user_id, order_id, amount, method, status)
        VALUES (?, ?, 100, 'CARD', 'PENDING')
    `, uuid.New(), orderID)
    
    err = testDB.DB.Exec(`
        INSERT INTO payment_transactions (user_id, order_id, amount, method, status)
        VALUES (?, ?, 200, 'CARD', 'PENDING')
    `, uuid.New(), orderID).Error
    assert.Error(t, err)  // duplicate order_id
}
```

## 10.3 Concurrency Test

```go
func TestConcurrentInsert(t *testing.T) {
    testDB := SetupPostgres(t)
    repo := postgres.NewPaymentRepository(testDB.DB)
    ctx := context.Background()
    
    const numGoroutines = 10
    const numInserts = 100
    
    var wg sync.WaitGroup
    errors := make(chan error, numGoroutines*numInserts)
    
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func(g int) {
            defer wg.Done()
            for j := 0; j < numInserts; j++ {
                payment, _ := entity.NewPayment(
                    uuid.New(), uuid.New(),
                    decimal.NewFromInt(100), "THB", "CARD",
                )
                if err := repo.Save(ctx, payment); err != nil {
                    errors <- err
                }
            }
        }(i)
    }
    
    wg.Wait()
    close(errors)
    
    // Check errors
    var errCount int
    for range errors {
        errCount++
    }
    assert.Equal(t, 0, errCount)
    
    // Verify count
    var count int
    testDB.DB.Raw("SELECT COUNT(*) FROM payment_transactions").Scan(&count)
    assert.Equal(t, numGoroutines*numInserts, count)
}
```

## 10.4 Slow Query Test

```go
func TestQueryPerformance(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping performance test")
    }
    
    testDB := SetupPostgres(t)
    repo := postgres.NewPaymentRepository(testDB.DB)
    ctx := context.Background()
    
    // Insert 10,000 rows
    for i := 0; i < 10000; i++ {
        p, _ := entity.NewPayment(uuid.New(), uuid.New(), decimal.NewFromInt(100), "THB", "CARD")
        _ = repo.Save(ctx, p)
    }
    
    userID := uuid.New()
    
    // Measure FindByUserID
    start := time.Now()
    _, err := repo.FindByUserID(ctx, userID, 10, 0)
    elapsed := time.Since(start)
    
    require.NoError(t, err)
    assert.Less(t, elapsed, 50*time.Millisecond, "query too slow")
}
```

---

# บทที่ 11: Kafka Integration Test

## 11.1 Setup Kafka Container

```go
//go:build integration

package kafka_test

import (
    "context"
    "testing"
    "time"
    
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/kafka"
    "github.com/IBM/sarama"
    "github.com/stretchr/testify/require"
)

type TestKafka struct {
    Container *kafka.KafkaContainer
    Brokers   []string
}

func SetupKafka(t *testing.T) *TestKafka {
    t.Helper()
    
    ctx := context.Background()
    
    container, err := kafka.RunContainer(ctx,
        testcontainers.WithImage("confluentinc/confluent-local:7.5.0"),
        kafka.WithClusterID("test-cluster"),
    )
    require.NoError(t, err)
    
    brokers, err := container.Brokers(ctx)
    require.NoError(t, err)
    
    t.Cleanup(func() {
        container.Terminate(ctx)
    })
    
    return &TestKafka{
        Container: container,
        Brokers:   brokers,
    }
}
```

## 11.2 Producer Test

```go
func TestKafkaProducer(t *testing.T) {
    tk := SetupKafka(t)
    
    // Producer
    producer, err := messaging.NewKafkaProducer(tk.Brokers)
    require.NoError(t, err)
    
    // Publish
    ctx := context.Background()
    event := events.PaymentCreated{
        PaymentID: uuid.New(),
        UserID:    uuid.New(),
    }
    
    err = producer.Publish(ctx, "payment.created", event.PaymentID.String(), event)
    require.NoError(t, err)
}
```

## 11.3 Consumer Test

```go
func TestKafkaConsumer(t *testing.T) {
    tk := SetupKafka(t)
    topic := "payment.created"
    
    // 1. Publish messages
    producer, _ := messaging.NewKafkaProducer(tk.Brokers)
    for i := 0; i < 5; i++ {
        event := events.PaymentCreated{
            BaseEvent: event.NewBaseEvent(),
            PaymentID: uuid.New(),
        }
        _ = producer.Publish(context.Background(), topic, event.PaymentID.String(), event)
    }
    
    // 2. Consume
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    received := make(chan events.PaymentCreated, 5)
    
    handler := &testHandler{onMessage: func(e events.PaymentCreated) {
        received <- e
    }}
    
    consumer := consumers.NewPaymentCreatedConsumer(handler)
    
    cfg := sarama.NewConfig()
    cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
    
    group, _ := sarama.NewConsumerGroup(tk.Brokers, "test-group", cfg)
    defer group.Close()
    
    go func() {
        for {
            if err := group.Consume(ctx, []string{topic}, consumer); err != nil {
                return
            }
        }
    }()
    
    // 3. Verify
    count := 0
    for count < 5 {
        select {
        case <-received:
            count++
        case <-ctx.Done():
            t.Fatalf("timeout: only received %d/5 messages", count)
        }
    }
    assert.Equal(t, 5, count)
}

type testHandler struct {
    onMessage func(events.PaymentCreated)
}

func (h *testHandler) Handle(ctx context.Context, e events.PaymentCreated) error {
    h.onMessage(e)
    return nil
}
```

## 11.4 Idempotency Test

```go
func TestKafkaConsumer_Idempotency(t *testing.T) {
    tk := SetupKafka(t)
    topic := "payment.created"
    
    // Send same message 3 times (simulate at-least-once)
    event := events.PaymentCreated{
        BaseEvent: event.NewBaseEvent(),
        PaymentID: uuid.New(),
        UserID:    uuid.New(),
    }
    
    producer, _ := messaging.NewKafkaProducer(tk.Brokers)
    for i := 0; i < 3; i++ {
        _ = producer.Publish(context.Background(), topic, event.PaymentID.String(), event)
    }
    
    // Consumer should process only once
    processCount := atomic.Int32{}
    handler := &idempotentHandler{
        seen:    make(map[uuid.UUID]bool),
        process: func() { processCount.Add(1) },
    }
    
    // ... consume ...
    
    assert.Equal(t, int32(1), processCount.Load())
}
```

---

# บทที่ 12: HTTP Handler Test

## 12.1 Handler Test Structure

```go
package http_test

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "icmongolang/internal/modules/payment/application"
    "icmongolang/internal/modules/payment/application/mocks"
    paymentHTTP "icmongolang/internal/modules/payment/interfaces/http"
    "icmongolang/pkg/logger"
    "icmongolang/pkg/validator"
)

// setupHandler creates handler with mocked use cases
func setupHandler(t *testing.T) *paymentHTTP.PaymentHandler {
    repo := new(mocks.PaymentRepositoryMock)
    log := logger.NewNoop()
    
    createUC := application.NewCreatePaymentUseCase(repo, log)
    refundUC := application.NewRefundPaymentUseCase(repo, log)
    
    return paymentHTTP.NewPaymentHandler(createUC, refundUC, validator.New())
}
```

## 12.2 Test POST Endpoint

```go
func TestPaymentHandler_Create_Success(t *testing.T) {
    // Setup
    repo := new(mocks.PaymentRepositoryMock)
    log := logger.NewNoop()
    createUC := application.NewCreatePaymentUseCase(repo, log)
    refundUC := application.NewRefundPaymentUseCase(repo, log)
    handler := paymentHTTP.NewPaymentHandler(createUC, refundUC, validator.New())
    
    orderID := uuid.New()
    repo.On("FindByOrderID", mock.Anything, orderID).
        Return(nil, domainerrors.ErrPaymentNotFound)
    repo.On("Save", mock.Anything, mock.Anything).Return(nil)
    
    // Request body
    body := map[string]string{
        "order_id": orderID.String(),
        "amount":   "100.00",
        "currency": "THB",
        "method":   "CARD",
    }
    bodyBytes, _ := json.Marshal(body)
    
    // Request
    req := httptest.NewRequest("POST", "/api/v1/payments", bytes.NewReader(bodyBytes))
    req.Header.Set("Content-Type", "application/json")
    
    // Add user_id to context (auth middleware ปกติทำ)
    ctx := context.WithValue(req.Context(), "user_id", uuid.New())
    req = req.WithContext(ctx)
    
    // Response recorder
    rr := httptest.NewRecorder()
    
    // Execute
    handler.Create(rr, req)
    
    // Assert
    assert.Equal(t, http.StatusCreated, rr.Code)
    
    var response map[string]interface{}
    err := json.Unmarshal(rr.Body.Bytes(), &response)
    require.NoError(t, err)
    
    assert.Equal(t, "PENDING", response["status"])
    assert.Equal(t, "100", response["amount"])
    assert.Equal(t, "THB", response["currency"])
}
```

## 12.3 Test Error Cases

```go
func TestPaymentHandler_Create_InvalidJSON(t *testing.T) {
    handler := setupHandler(t)
    
    req := httptest.NewRequest("POST", "/api/v1/payments", 
        strings.NewReader("not json"))
    req.Header.Set("Content-Type", "application/json")
    
    ctx := context.WithValue(req.Context(), "user_id", uuid.New())
    req = req.WithContext(ctx)
    
    rr := httptest.NewRecorder()
    handler.Create(rr, req)
    
    assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPaymentHandler_Create_Unauthorized(t *testing.T) {
    handler := setupHandler(t)
    
    body := []byte(`{"order_id":"...","amount":"100","currency":"THB","method":"CARD"}`)
    req := httptest.NewRequest("POST", "/api/v1/payments", bytes.NewReader(body))
    // ← ไม่มี user_id ใน context
    
    rr := httptest.NewRecorder()
    handler.Create(rr, req)
    
    assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestPaymentHandler_Create_ValidationErrors(t *testing.T) {
    tests := []struct {
        name string
        body map[string]string
    }{
        {
            name: "missing order_id",
            body: map[string]string{
                "amount": "100", "currency": "THB", "method": "CARD",
            },
        },
        {
            name: "invalid order_id",
            body: map[string]string{
                "order_id": "not-a-uuid", "amount": "100", "currency": "THB", "method": "CARD",
            },
        },
        {
            name: "invalid method",
            body: map[string]string{
                "order_id": uuid.New().String(), "amount": "100", 
                "currency": "THB", "method": "INVALID",
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            handler := setupHandler(t)
            
            bodyBytes, _ := json.Marshal(tt.body)
            req := httptest.NewRequest("POST", "/api/v1/payments", bytes.NewReader(bodyBytes))
            ctx := context.WithValue(req.Context(), "user_id", uuid.New())
            req = req.WithContext(ctx)
            
            rr := httptest.NewRecorder()
            handler.Create(rr, req)
            
            assert.Equal(t, http.StatusBadRequest, rr.Code)
        })
    }
}
```

## 12.4 Test Full Router

```go
func TestPaymentRoutes(t *testing.T) {
    // Setup full router
    r := chi.NewRouter()
    
    // Auth middleware (mock)
    authMW := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx := context.WithValue(r.Context(), "user_id", uuid.New())
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
    
    handler := setupHandler(t)
    paymentHTTP.RegisterRoutes(r, &paymentHTTP.Handlers{
        Payment: handler,
    }, authMW)
    
    // Test via httptest.Server
    server := httptest.NewServer(r)
    defer server.Close()
    
    // Request
    orderID := uuid.New()
    body := map[string]string{
        "order_id": orderID.String(),
        "amount":   "100.00",
        "currency": "THB",
        "method":   "CARD",
    }
    bodyBytes, _ := json.Marshal(body)
    
    resp, err := http.Post(
        server.URL+"/api/v1/payments",
        "application/json",
        bytes.NewReader(bodyBytes),
    )
    require.NoError(t, err)
    defer resp.Body.Close()
    
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
```

## 12.5 Benchmark Handler

```go
func BenchmarkPaymentHandler_Create(b *testing.B) {
    repo := new(mocks.PaymentRepositoryMock)
    log := logger.NewNoop()
    createUC := application.NewCreatePaymentUseCase(repo, log)
    refundUC := application.NewRefundPaymentUseCase(repo, log)
    handler := paymentHTTP.NewPaymentHandler(createUC, refundUC, validator.New())
    
    orderID := uuid.New()
    repo.On("FindByOrderID", mock.Anything, mock.Anything).
        Return(nil, domainerrors.ErrPaymentNotFound)
    repo.On("Save", mock.Anything, mock.Anything).Return(nil)
    
    body := map[string]string{
        "order_id": orderID.String(),
        "amount":   "100.00",
        "currency": "THB",
        "method":   "CARD",
    }
    bodyBytes, _ := json.Marshal(body)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        req := httptest.NewRequest("POST", "/api/v1/payments", bytes.NewReader(bodyBytes))
        ctx := context.WithValue(req.Context(), "user_id", uuid.New())
        req = req.WithContext(ctx)
        
        rr := httptest.NewRecorder()
        handler.Create(rr, req)
    }
}
```

---

# บทที่ 13: E2E Test

## 13.1 หลักการ

**E2E Test** = Test ทั้งระบบจาก user perspective

**Characteristics:**
- ใช้ infrastructure จริง
- Test user journey
- Business value สูง
- จำนวนน้อย (5-10 scenarios)

## 13.2 E2E Test Setup

```go
//go:build e2e

package e2e_test

import (
    "context"
    "os"
    "testing"
    "time"
    
    "github.com/stretchr/testify/require"
)

// TestEnv setup ทั้งระบบ
type TestEnv struct {
    APIServer  string
    DB         *gorm.DB
    Redis      *redis.Client
    KafkaBroker []string
    cleanup    func()
}

func SetupTestEnv(t *testing.T) *TestEnv {
    t.Helper()
    
    // Option 1: Docker Compose
    // docker-compose -f docker-compose.test.yml up -d
    
    // Option 2: Testcontainers
    // spin up all containers
    
    // Option 3: Use existing staging
    // apiURL := os.Getenv("E2E_API_URL")
    
    apiURL := os.Getenv("E2E_API_URL")
    if apiURL == "" {
        t.Skip("E2E_API_URL not set")
    }
    
    return &TestEnv{
        APIServer: apiURL,
        // ...
    }
}
```

## 13.3 User Journey Test

```go
//go:build e2e

func TestE2E_SignupLoginCreatePayment(t *testing.T) {
    env := SetupTestEnv(t)
    client := NewAPIClient(env.APIServer)
    
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()
    
    // Step 1: Signup
    email := fmt.Sprintf("test-%d@example.com", time.Now().UnixNano())
    signupResp, err := client.Signup(ctx, SignupRequest{
        Email:    email,
        Password: "Test123!@#",
        Name:     "Test User",
    })
    require.NoError(t, err)
    require.NotEmpty(t, signupResp.UserID)
    
    // Step 2: Verify email (skip in test)
    // ...
    
    // Step 3: Login
    loginResp, err := client.Login(ctx, LoginRequest{
        Email:    email,
        Password: "Test123!@#",
    })
    require.NoError(t, err)
    require.NotEmpty(t, loginResp.AccessToken)
    
    client.SetToken(loginResp.AccessToken)
    
    // Step 4: Create order
    orderResp, err := client.CreateOrder(ctx, CreateOrderRequest{
        Items: []OrderItem{
            {ProductID: "prod-1", Quantity: 2},
        },
    })
    require.NoError(t, err)
    require.NotEmpty(t, orderResp.ID)
    
    // Step 5: Create payment
    paymentResp, err := client.CreatePayment(ctx, CreatePaymentRequest{
        OrderID:  orderResp.ID,
        Amount:   "100.00",
        Currency: "THB",
        Method:   "CARD",
    })
    require.NoError(t, err)
    require.NotEmpty(t, paymentResp.ID)
    assert.Equal(t, "PENDING", paymentResp.Status)
    
    // Step 6: Verify payment exists
    verifyResp, err := client.GetPayment(ctx, paymentResp.ID)
    require.NoError(t, err)
    assert.Equal(t, paymentResp.ID, verifyResp.ID)
}
```

## 13.4 API Client

```go
// testutils/apiclient/client.go
package apiclient

type Client struct {
    baseURL string
    http    *http.Client
    token   string
}

func NewClient(baseURL string) *Client {
    return &Client{
        baseURL: baseURL,
        http:    &http.Client{Timeout: 30 * time.Second},
    }
}

func (c *Client) SetToken(token string) {
    c.token = token
}

func (c *Client) request(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
    var bodyReader io.Reader
    if body != nil {
        data, _ := json.Marshal(body)
        bodyReader = bytes.NewReader(data)
    }
    
    req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    if c.token != "" {
        req.Header.Set("Authorization", "Bearer "+c.token)
    }
    
    return c.http.Do(req)
}

func (c *Client) Signup(ctx context.Context, req SignupRequest) (*SignupResponse, error) {
    resp, err := c.request(ctx, "POST", "/api/v1/auth/signup", req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusCreated {
        return nil, fmt.Errorf("signup failed: %d", resp.StatusCode)
    }
    
    var result SignupResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    return &result, nil
}

// ... other methods
```

## 13.5 Smoke Test (ใช้หลัง deploy)

```go
//go:build smoke

func TestSmoke(t *testing.T) {
    baseURL := os.Getenv("SMOKE_URL")
    if baseURL == "" {
        t.Skip("SMOKE_URL not set")
    }
    
    // 1. Health check
    resp, err := http.Get(baseURL + "/health")
    require.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    resp.Body.Close()
    
    // 2. Login
    client := NewAPIClient(baseURL)
    loginResp, err := client.Login(context.Background(), LoginRequest{
        Email:    "smoke@test.com",
        Password: "SmokeTest123!",
    })
    require.NoError(t, err)
    
    client.SetToken(loginResp.AccessToken)
    
    // 3. List payments (must work)
    resp, err = client.Get(ctx, "/api/v1/payments?limit=1")
    require.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    resp.Body.Close()
}
```

**Run in CI/CD:**
```yaml
- name: Smoke test
  run: |
    SMOKE_URL=https://api.example.com go test -tags=smoke ./tests/smoke/...
```

---

# บทที่ 14: Load Testing ด้วย k6

## 14.1 k6 Setup

```bash
# Install
brew install k6  # macOS
# หรือ
docker pull grafana/k6
```

## 14.2 Basic Load Test

```javascript
// tests/load/basic.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '2m', target: 100 },   // Ramp up to 100 users
    { duration: '5m', target: 100 },   // Stay at 100
    { duration: '2m', target: 0 },     // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<200'],  // 95% < 200ms
    http_req_failed: ['rate<0.01'],    // < 1% errors
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const TOKEN = __ENV.TOKEN;

export default function () {
  const res = http.get(`${BASE_URL}/health`);
  
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 200ms': (r) => r.timings.duration < 200,
  });
  
  sleep(1);
}
```

## 14.3 Full Payment Flow Test

```javascript
// tests/load/payment-flow.js
import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Custom metrics
const createPaymentTime = new Trend('create_payment_time');
const createPaymentErrors = new Rate('create_payment_errors');

export const options = {
  scenarios: {
    // Constant rate: 50 req/s for 5 minutes
    constant_rate: {
      executor: 'constant-arrival-rate',
      rate: 50,
      timeUnit: '1s',
      duration: '5m',
      preAllocatedVUs: 100,
      maxVUs: 500,
    },
    
    // Spike: sudden burst
    spike: {
      executor: 'ramping-arrival-rate',
      startTime: '5m',
      stages: [
        { duration: '10s', target: 500 },  // Spike
        { duration: '1m', target: 500 },   // Hold
        { duration: '10s', target: 50 },   // Back to normal
      ],
      preAllocatedVUs: 500,
    },
  },
  
  thresholds: {
    'http_req_duration{scenario:constant_rate}': ['p(95)<200'],
    'http_req_duration{scenario:spike}': ['p(95)<500'],
    create_payment_errors: ['rate<0.01'],
  },
};

const BASE_URL = __ENV.BASE_URL;
const TOKEN = __ENV.TOKEN;

export default function () {
  group('Payment Flow', function () {
    // Create payment
    const orderID = `${__VU}-${__ITER}-${Date.now()}`;
    const payload = JSON.stringify({
      order_id: orderID,
      amount: '100.00',
      currency: 'THB',
      method: 'CARD',
    });
    
    const start = Date.now();
    const res = http.post(`${BASE_URL}/api/v1/payments`, payload, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${TOKEN}`,
      },
    });
    
    createPaymentTime.add(Date.now() - start);
    createPaymentErrors.add(res.status >= 400);
    
    check(res, {
      'create status is 201': (r) => r.status === 201,
    });
    
    // Get payment
    if (res.status === 201) {
      const payment = JSON.parse(res.body);
      const getRes = http.get(`${BASE_URL}/api/v1/payments/${payment.id}`, {
        headers: { 'Authorization': `Bearer ${TOKEN}` },
      });
      
      check(getRes, {
        'get status is 200': (r) => r.status === 200,
      });
    }
  });
  
  sleep(1);
}
```

## 14.4 Run Load Test

```bash
# Basic
k6 run tests/load/basic.js

# With env
BASE_URL=https://api.example.com TOKEN=xxx k6 run tests/load/basic.js

# With output
k6 run --out json=results.json tests/load/basic.js

# With InfluxDB
k6 run --out influxdb=http://localhost:8086/k6 tests/load/basic.js

# With Prometheus
K6_PROMETHEUS_RW_SERVER_URL=http://localhost:9090/api/v1/write \
  k6 run --out experimental-prometheus-rw tests/load/basic.js
```

## 14.5 Analysis

```bash
# Basic stats
k6 run tests/load/basic.js

# Results:
#      ✓ status is 200
#      ✓ response time < 200ms
#
#      checks.........................: 100.00% ✓ 5000 ✗ 0
#      data_received..................: 1.2 MB  20 kB/s
#      data_sent......................: 500 kB  8.3 kB/s
#      http_req_blocked...............: avg=10µs   p(95)=50µs
#      http_req_duration..............: avg=45ms   p(95)=120ms
#      http_req_failed................: 0.00%  ✓ 0
#      http_reqs......................: 2500    41/s
#      iterations.....................: 500     8.3/s
#      vus............................: 100
#      vus_max........................: 100
```

## 14.6 Load Test Scenarios

| Scenario | Purpose | Config |
|---|---|---|
| **Smoke** | Sanity check | 1 user, 1 min |
| **Load** | Normal load | Expected users, 10 min |
| **Stress** | Find breaking point | Ramp up until fail |
| **Spike** | Sudden burst | 10x spike |
| **Soak** | Long-running | 4-24 hours |
| **Breakpoint** | Find limits | Increase until error |

---

# บทที่ 15: Chaos Testing

## 15.1 Chaos Engineering Principles

```
1. Define Steady State (metrics ปกติ)
2. Hypothesize (คาดว่าจะไม่พัง)
3. Introduce Chaos (inject failure)
4. Observe (ดูผล)
5. Learn (ปรับปรุง)
```

## 15.2 Chaos Mesh Setup

```yaml
# kubernetes/chaos/pod-kill.yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: PodChaos
metadata:
  name: payment-pod-kill
  namespace: default
spec:
  action: pod-kill
  mode: one
  selector:
    namespaces:
      - default
    labelSelectors:
      app: payment-api
  scheduler:
    cron: '@every 1h'
  duration: '30s'
```

## 15.3 Chaos Scenarios

### 15.3.1 Pod Kill

```yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: PodChaos
metadata:
  name: pod-kill
spec:
  action: pod-kill
  mode: fixed-percent
  value: '30'  # 30% of pods
  selector:
    labelSelectors:
      app: payment-api
  duration: '5m'
```

**Verify:**
- No 5xx errors
- Latency stable (failover works)
- Auto-recovery within 30s

### 15.3.2 Network Latency

```yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: NetworkChaos
metadata:
  name: network-latency
spec:
  action: delay
  mode: all
  selector:
    labelSelectors:
      app: payment-api
  delay:
    latency: '100ms'
    correlation: '100'
    jitter: '10ms'
  duration: '5m'
```

**Verify:**
- Client timeout ไม่เกิด
- Retry ทำงาน
- Metrics capture latency spike

### 15.3.3 Database Failure

```yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: PodChaos
metadata:
  name: postgres-fail
spec:
  action: pod-failure
  mode: one
  selector:
    labelSelectors:
      app: postgres
  duration: '2m'
```

**Verify:**
- Circuit breaker ทำงาน
- Fallback ทำงาน (cache?)
- Alert ส่งทันที

### 15.3.4 CPU Stress

```yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: StressChaos
metadata:
  name: cpu-stress
spec:
  mode: one
  selector:
    labelSelectors:
      app: payment-api
  stressors:
    cpu:
      workers: 2
      load: 80
  duration: '5m'
```

### 15.3.5 HTTP Error Injection

```yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: HTTPChaos
metadata:
  name: http-error
spec:
  mode: one
  selector:
    labelSelectors:
      app: payment-api
  target: Request
  port: 8080
  path: '/api/v1/payments/*'
  abort: true  # Inject abort
  duration: '5m'
```

## 15.4 Chaos Test Workflow

```bash
# 1. Baseline metrics
k6 run tests/load/baseline.js > baseline.txt

# 2. Start chaos
kubectl apply -f chaos/pod-kill.yaml

# 3. Run load test during chaos
k6 run tests/load/chaos.js > chaos.txt

# 4. Compare
benchstat baseline.txt chaos.txt

# 5. Stop chaos
kubectl delete -f chaos/pod-kill.yaml

# 6. Verify recovery
k6 run tests/load/baseline.js
```

## 15.5 Chaos Test Checklist

- [ ] Baseline metrics known
- [ ] Rollback plan ready
- [ ] Alert configured
- [ ] Team on standby
- [ ] Non-production first
- [ ] Limited blast radius
- [ ] Auto-recovery verified

---

# บทที่ 16: Security Testing

## 16.1 Static Analysis

### 16.1.1 gosec

```bash
# Install
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Run
gosec ./...

# Output SARIF (for GitHub)
gosec -fmt sarif -out gosec.sarif ./...
```

**ตัวอย่าง Finding:**
```
[G404] Insecure random number source (rand)
  > math/rand is not cryptographically secure
  File: internal/utils/token.go, Line: 15
```

### 16.1.2 govulncheck

```bash
# Install
go install golang.org/x/vuln/cmd/govulncheck@latest

# Run
govulncheck ./...
```

### 16.1.3 staticcheck

```bash
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
```

## 16.2 Dependency Scanning

### 16.2.1 Trivy

```bash
# Filesystem
trivy fs --severity HIGH,CRITICAL .

# Container image
trivy image ghcr.io/you/module:latest

# SBOM
trivy image --format spdx-json ghcr.io/you/module:latest > sbom.json
```

### 16.2.2 Snyk

```bash
snyk test
snyk monitor  # continuous monitoring
```

### 16.2.3 Dependabot

```yaml
# .github/dependabot.yml
version: 2
updates:
  - package-ecosystem: gomod
    directory: /
    schedule:
      interval: weekly
    open-pull-requests-limit: 10
    reviewers:
      - your-team
    labels:
      - dependencies
    groups:
      patch-updates:
        patterns: ["*"]
        update-types: ["patch"]
```

## 16.3 SAST (Static Application Security Testing)

```yaml
# .github/workflows/security.yml
name: Security

on:
  push:
    branches: [main]
  pull_request:
  schedule:
    - cron: '0 0 * * 0'  # weekly

jobs:
  gosec:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: securego/gosec@master
        with:
          args: '-no-fail -fmt sarif -out gosec.sarif ./...'
      - uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: gosec.sarif

  govulncheck:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: golang/govulncheck-action@v1
        with:
          go-version-input: '1.22'
          go-package: './...'

  trivy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          severity: 'HIGH,CRITICAL'
          exit-code: '1'
          format: 'sarif'
          output: 'trivy.sarif'
      - uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: trivy.sarif
```

## 16.4 DAST (Dynamic Application Security Testing)

### 16.4.1 OWASP ZAP

```bash
# Run ZAP scan
docker run -t owasp/zap2docker-stable zap-baseline.py \
    -t https://api.example.com \
    -r report.html
```

```yaml
# .github/workflows/dast.yml
- name: Run OWASP ZAP
  uses: zaproxy/action-baseline@v0.12.0
  with:
    target: 'https://staging.api.example.com'
    rules_file_name: '.zap/rules.tsv'
    cmd_options: '-a'
```

### 16.4.2 Authentication Tests

```go
// Security test cases
func TestSecurity_Authentication(t *testing.T) {
    tests := []struct {
        name      string
        header    string
        wantCode  int
    }{
        {
            name:     "no auth header",
            header:   "",
            wantCode: http.StatusUnauthorized,
        },
        {
            name:     "invalid token",
            header:   "Bearer invalid.token.here",
            wantCode: http.StatusUnauthorized,
        },
        {
            name:     "expired token",
            header:   "Bearer " + expiredToken,
            wantCode: http.StatusUnauthorized,
        },
        {
            name:     "wrong signature",
            header:   "Bearer " + wrongSignature,
            wantCode: http.StatusUnauthorized,
        },
        {
            name:     "valid token",
            header:   "Bearer " + validToken,
            wantCode: http.StatusOK,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("GET", "/api/v1/payments", nil)
            if tt.header != "" {
                req.Header.Set("Authorization", tt.header)
            }
            rr := httptest.NewRecorder()
            handler.ServeHTTP(rr, req)
            
            assert.Equal(t, tt.wantCode, rr.Code)
        })
    }
}
```

### 16.4.3 SQL Injection Test

```go
func TestSecurity_SQLInjection(t *testing.T) {
    payloads := []string{
        "'; DROP TABLE users; --",
        "' OR '1'='1",
        "1'; UPDATE users SET role='admin' WHERE id=1; --",
        "1 UNION SELECT * FROM users --",
    }
    
    for _, payload := range payloads {
        t.Run(payload, func(t *testing.T) {
            // ต้องไม่เกิด error (escape ถูกต้อง)
            _, err := repo.FindByUserID(ctx, uuid.UUID{}, 10, 0)
            // ... test ที่ใช้ payload
        })
    }
}
```

### 16.4.4 XSS Test

```go
func TestSecurity_XSS(t *testing.T) {
    payloads := []string{
        "<script>alert(1)</script>",
        "<img src=x onerror=alert(1)>",
        "javascript:alert(1)",
    }
    
    for _, payload := range payloads {
        resp := client.CreatePayment(ctx, CreatePaymentRequest{
            Notes: payload,
        })
        
        // Verify: payload ถูก escape หรือ strip
        assert.NotContains(t, resp.Notes, "<script>")
    }
}
```

## 16.5 Penetration Testing Checklist

- [ ] Authentication bypass
- [ ] Authorization bypass (IDOR)
- [ ] SQL injection
- [ ] NoSQL injection
- [ ] Command injection
- [ ] XSS (stored, reflected, DOM)
- [ ] CSRF
- [ ] SSRF
- [ ] XXE
- [ ] Insecure deserialization
- [ ] Rate limiting bypass
- [ ] Sensitive data exposure
- [ ] Security misconfiguration
- [ ] Broken access control
- [ ] Vulnerable dependencies

---

# บทที่ 17: Contract Testing

## 17.1 หลักการ

**Contract Testing** = Test ว่า 2 services คุยกันได้ตาม "สัญญา" (contract)

**ปัญหา:** Consumer เปลี่ยน → Producer ไม่รู้ → Production ล่ม

**แก้:** Contract เป็น documentation ที่ test อัตโนมัติ

## 17.2 Pact Setup

```bash
# Install pact-go
go get github.com/pact-foundation/pact-go/v2
```

## 17.3 Consumer Test

```go
//go:build contract

package contract_test

import (
    "testing"
    
    "github.com/pact-foundation/pact-go/v2/consumer"
    "github.com/pact-foundation/pact-go/v2/matchers"
)

func TestPaymentAPIContract(t *testing.T) {
    mockProvider, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
        Consumer: "payment-client",
        Provider: "payment-api",
    })
    require.NoError(t, err)
    
    // Expect request
    mockProvider.
        AddInteraction().
        Given("payment order-123 exists").
        UponReceiving("a request for payment order-123").
        WithRequest("GET", "/api/v1/payments/order-123").
        WillRespondWith(200).
        WithBody(matchers.Map{
            "id":     matchers.Like("uuid-xxx"),
            "amount": matchers.Like("100.00"),
            "status": matchers.Term("PENDING", "PENDING|SUCCESS|FAILED"),
        })
    
    // Test
    err = mockProvider.Verify(t, func() error {
        client := NewClient(mockProvider.Server.URL)
        resp, err := client.GetPayment(context.Background(), "order-123")
        if err != nil {
            return err
        }
        assert.NotEmpty(t, resp.ID)
        assert.Equal(t, "100.00", resp.Amount)
        return nil
    })
    require.NoError(t, err)
}
```

## 17.4 Provider Test

```go
//go:build contract

package contract_test

import (
    "testing"
    
    "github.com/pact-foundation/pact-go/v2/provider"
)

func TestProviderContract(t *testing.T) {
    verifier := provider.NewVerifier()
    
    err := verifier.VerifyProvider(t, provider.VerifyRequest{
        ProviderBaseURL: "http://localhost:8080",
        PactURLs:        []string{"./pacts/payment-client-payment-api.json"},
        ProviderStatesSetupURL: "http://localhost:8080/_pact/state",
    })
    require.NoError(t, err)
}
```

## 17.5 Provider States

```go
// Handler for pact state setup
func setupPactState(w http.ResponseWriter, r *http.Request) {
    var state struct {
        State string                 `json:"state"`
        Params map[string]interface{} `json:"params"`
    }
    json.NewDecoder(r.Body).Decode(&state)
    
    switch state.State {
    case "payment order-123 exists":
        // Setup: create payment in DB
        db.Exec("INSERT INTO payment_transactions (id, order_id, ...) VALUES (...)")
    
    case "payment order-456 does not exist":
        // No-op
    }
    
    w.WriteHeader(200)
}
```

---

# บทที่ 18: CI Pipeline Design

## 18.1 CI Pipeline Stages

```
┌──────────────┐
│  Git Push    │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Checkout    │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Setup Go    │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Lint        │ ← golangci-lint
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Unit Tests  │ ← go test
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Integration │ ← testcontainers
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Coverage    │ ← gate 80%
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Security    │ ← gosec, trivy
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Build       │ ← go build
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Container   │ ← docker build
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Push        │ ← registry
└──────────────┘
```

## 18.2 Pipeline Speed Targets

| Stage | Target | Max |
|---|---|---|
| Lint | 30s | 2m |
| Unit Tests | 2m | 5m |
| Integration | 3m | 10m |
| Security | 1m | 5m |
| Build | 1m | 3m |
| Container | 2m | 5m |
| **Total** | **10m** | **30m** |

## 18.3 Caching Strategy

```yaml
- name: Setup Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.22'
    cache: true  # ← cache go modules
    cache-dependency-path: go.sum

- name: Cache Go build
  uses: actions/cache@v4
  with:
    path: |
      ~/.cache/go-build
      ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}

- name: Docker layer cache
  uses: docker/build-push-action@v5
  with:
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

---

# บทที่ 19: GitHub Actions Workflow

## 19.1 Full CI Pipeline

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

env:
  GO_VERSION: '1.22'
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  # ============================================================
  # Lint
  # ============================================================
  lint:
    name: Lint
    runs-on: ubuntu-latest
    timeout-minutes: 10
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true
      
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: v1.59
          args: --timeout=5m --config=.golangci.yml

  # ============================================================
  # Unit Tests
  # ============================================================
  test-unit:
    name: Unit Tests
    runs-on: ubuntu-latest
    timeout-minutes: 10
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true
      
      - name: Run unit tests
        run: |
          go test -short -race -coverprofile=coverage.out -covermode=atomic ./...
      
      - name: Coverage gate
        run: |
          cov=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Total coverage: $cov%"
          if (( $(echo "$cov < 80" | bc -l) )); then
            echo "::error::Coverage $cov% < 80%"
            exit 1
          fi
      
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.out
          flags: unittests
          fail_ci_if_error: false

  # ============================================================
  # Integration Tests
  # ============================================================
  test-integration:
    name: Integration Tests
    runs-on: ubuntu-latest
    timeout-minutes: 20
    services:
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_DB: testdb
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
      redis:
        image: redis:7-alpine
        ports:
          - 6379:6379
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true
      
      - name: Run integration tests
        run: |
          go test -tags=integration -race -timeout 15m ./...
        env:
          DB_DSN: postgres://test:test@localhost:5432/testdb?sslmode=disable
          REDIS_ADDR: localhost:6379

  # ============================================================
  # Security Scan
  # ============================================================
  security:
    name: Security
    runs-on: ubuntu-latest
    timeout-minutes: 15
    permissions:
      contents: read
      security-events: write
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      
      - name: Run gosec
        uses: securego/gosec@master
        with:
          args: '-no-fail -fmt sarif -out gosec.sarif ./...'
      
      - name: Upload gosec SARIF
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: gosec.sarif
      
      - name: Run govulncheck
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...
      
      - name: Run Trivy
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          severity: 'HIGH,CRITICAL'
          format: 'sarif'
          output: 'trivy.sarif'
          exit-code: '0'  # Report only
      
      - name: Upload Trivy SARIF
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: trivy.sarif

  # ============================================================
  # Build Binary
  # ============================================================
  build:
    name: Build
    runs-on: ubuntu-latest
    timeout-minutes: 10
    needs: [lint, test-unit]
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true
      
      - name: Build
        run: |
          CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
            -ldflags="-w -s -X main.version=${{ github.sha }} -X main.commit=${{ github.sha }}" \
            -o bin/api ./cmd/api
      
      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: api-binary
          path: bin/api
          retention-days: 7

  # ============================================================
  # Build Container
  # ============================================================
  docker:
    name: Docker Build
    runs-on: ubuntu-latest
    timeout-minutes: 15
    needs: [build]
    permissions:
      contents: read
      packages: write
    steps:
      - uses: actions/checkout@v4
      
      - uses: docker/setup-buildx-action@v3
      
      - name: Login to registry
        if: github.event_name != 'pull_request'
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      
      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=ref,event=branch
            type=ref,event=pr
            type=semver,pattern={{version}}
            type=sha,prefix={{branch}}-
            type=raw,value=latest,enable={{is_default_branch}}
      
      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          push: ${{ github.event_name != 'pull_request' }}
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
          platforms: linux/amd64,linux/arm64
          provenance: true
          sbom: true
```

## 19.2 Release Workflow

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags: ['v*']

permissions:
  contents: write
  packages: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      
      - name: Run tests
        run: go test -race ./...
      
      - name: Generate changelog
        uses: orhun/git-cliff-action@v3
        with:
          config: cliff.toml
          args: --verbose
        env:
          OUTPUT: CHANGELOG.md
      
      - name: Create GitHub Release
        uses: softprops/action-gh-release@v1
        with:
          body_path: CHANGELOG.md
          generate_release_notes: true
          draft: false
      
      - name: Build multi-arch
        run: |
          for GOOS in linux darwin windows; do
            for GOARCH in amd64 arm64; do
              EXT=""
              if [ "$GOOS" = "windows" ]; then EXT=".exe"; fi
              CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build \
                -ldflags="-s -w -X main.version=${{ github.ref_name }}" \
                -o "dist/module-${GOOS}-${GOARCH}${EXT}" ./cmd/api
            done
          done
      
      - name: Upload binaries
        uses: softprops/action-gh-release@v1
        with:
          files: dist/*
```

---

# บทที่ 20: Artifact Management

## 20.1 Artifact Types

| Type | ใช้กับ | เก็บที่ |
|---|---|---|
| **Binary** | CLI, server | GitHub Releases, S3 |
| **Container** | Deployment | ghcr.io, ECR, Docker Hub |
| **Go Module** | Library | pkg.go.dev, private proxy |
| **Documentation** | Docs | GitHub Pages |
| **SBOM** | Compliance | Artifact storage |

## 20.2 Container Registry

### 20.2.1 GHCR (GitHub Container Registry)

```bash
# Login
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# Tag + push
docker tag your/module:latest ghcr.io/you/module:v1.0.0
docker push ghcr.io/you/module:v1.0.0

# Pull
docker pull ghcr.io/you/module:v1.0.0
```

### 20.2.2 AWS ECR

```bash
# Login
aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin \
  xxx.dkr.ecr.us-east-1.amazonaws.com

# Create repo
aws ecr create-repository --repository-name you/module --region us-east-1

# Tag + push
docker tag your/module:latest xxx.dkr.ecr.us-east-1.amazonaws.com/you/module:v1.0.0
docker push xxx.dkr.ecr.us-east-1.amazonaws.com/you/module:v1.0.0
```

### 20.2.3 Image Tagging Strategy

```
ghcr.io/you/module:latest           ← Latest main
ghcr.io/you/module:v1.2.3           ← Semver tag
ghcr.io/you/module:v1.2             ← Minor
ghcr.io/you/module:1                ← Major
ghcr.io/you/module:main-abc123      ← Branch-SHA
ghcr.io/you/module:sha-abc123       ← SHA only
```

## 20.3 SBOM Generation

```bash
# Syft
syft ghcr.io/you/module:v1.0.0 -o spdx-json > sbom.json

# Trivy
trivy image --format spdx-json ghcr.io/you/module:v1.0.0 > sbom.json
```

**ใน CI:**
```yaml
- name: Generate SBOM
  uses: anchore/sbom-action@v0
  with:
    image: ghcr.io/${{ github.repository }}:${{ github.sha }}
    format: spdx-json
    output-file: sbom.spdx.json

- name: Attach SBOM to release
  uses: softprops/action-gh-release@v1
  with:
    files: sbom.spdx.json
```

## 20.4 Signing

### 20.4.1 Cosign (Container Signing)

```bash
# Generate key
cosign generate-key-pair

# Sign
cosign sign --key cosign.key ghcr.io/you/module:v1.0.0

# Verify
cosign verify --key cosign.pub ghcr.io/you/module:v1.0.0
```

**ใน CI:**
```yaml
- name: Install cosign
  uses: sigstore/cosign-installer@v3

- name: Sign image
  run: |
    cosign sign --yes \
      ghcr.io/${{ github.repository }}@${{ steps.build.outputs.digest }}
```

### 20.4.2 SLSA Provenance

```yaml
- name: Generate provenance
  uses: slsa-framework/slsa-github-generator/.github/workflows/generator_container_slsa3.yml@v1.9.0
  with:
    image: ghcr.io/${{ github.repository }}
    digest: ${{ steps.build.outputs.digest }}
```

---

# บทที่ 21: Container Build & Registry

## 21.1 Multi-Stage Dockerfile

```dockerfile
# ---------- Stage 1: Build ----------
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata

WORKDIR /app

# Cache go.mod / go.sum
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source
COPY . .

# Build args
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w \
              -X main.version=${VERSION} \
              -X main.commit=${COMMIT} \
              -X main.buildTime=${BUILD_TIME}" \
    -trimpath \
    -o /out/api ./cmd/api

# ---------- Stage 2: Runtime ----------
FROM alpine:3.19

# Runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    wget \
    && rm -rf /var/cache/apk/*

# Non-root user
RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

# Copy binary
COPY --from=builder --chown=app:app /out/api .

# Copy migrations
COPY --from=builder --chown=app:app /app/migrations ./migrations

USER app

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -q --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/api"]
```

## 21.2 Distroless (เล็กสุด)

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/api ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/api /api
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/api"]
```

## 21.3 Image Optimization

| Optimization | Size Reduction |
|---|---|
| **Multi-stage** | 90% |
| **Alpine base** | 70% |
| **Distroless** | 95% |
| **Scratch** | 98% |
| **UPX compress** | 50% (binary) |
| **Static build** | No libc needed |

**Example:**
```
golang:1.22          1.1 GB
golang:1.22-alpine   300 MB
alpine:3.19           8 MB
distroless            3 MB
scratch               ~0 MB (just your binary)
```

## 21.4 Image Scanning

```yaml
# .github/workflows/scan.yml
- name: Scan image
  uses: aquasecurity/trivy-action@master
  with:
    image-ref: ghcr.io/${{ github.repository }}:${{ github.sha }}
    severity: 'HIGH,CRITICAL'
    format: 'sarif'
    output: 'trivy-image.sarif'
    exit-code: '1'

- name: Upload SARIF
  uses: github/codeql-action/upload-sarif@v3
  with:
    sarif_file: trivy-image.sarif
```

## 21.5 Registry Cleanup

```yaml
# .github/workflows/cleanup.yml
name: Cleanup

on:
  schedule:
    - cron: '0 0 * * 0'  # weekly

jobs:
  cleanup:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/delete-package-versions@v5
        with:
          package-name: 'module'
          package-type: 'container'
          min-versions-to-keep: 20
          delete-only-untagged-versions: 'true'
          token: ${{ secrets.GITHUB_TOKEN }}
```

---

# บทที่ 22: Deployment Strategies

## 22.1 Comparison

| Strategy | Downtime | Risk | Rollback | Resource |
|---|---|---|---|---|
| **Recreate** | Yes | High | Slow | 1x |
| **Rolling** | No | Medium | Medium | 1.x |
| **Blue-Green** | No | Low | Fast | 2x |
| **Canary** | No | Very Low | Very Fast | 1.1x |
| **A/B** | No | Low | Fast | 2x |
| **Shadow** | No | Very Low | N/A | 2x |

## 22.2 Recreate

```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  strategy:
    type: Recreate  # ← ปิดทั้งหมด → เปิดใหม่
```

**ใช้เมื่อ:**
- Dev/test environment
- Database migration ต้อง restart
- ไม่มี user

**ข้อเสีย:** Downtime

## 22.3 Rolling Update

```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1        # เพิ่ม 1 pod
      maxUnavailable: 0  # ไม่ลด pod ที่ ready
```

**Process:**
```
Start:  [v1] [v1] [v1]
Step 1: [v2] [v1] [v1]  ← add new
Step 2: [v2] [v1] [v1]  ← wait ready
Step 3: [v2] [v2] [v1]  ← remove old + add new
Step 4: [v2] [v2] [v1]
Step 5: [v2] [v2] [v2]
Done!
```

## 22.4 Blue-Green

```
┌──────────────────────────────────────┐
│       Load Balancer                  │
└────────────┬─────────────────────────┘
             │
      ┌──────┴──────┐
      │             │
      ▼             ▼
   ┌─────┐       ┌─────┐
   │Blue │       │Green│
   │(v1) │       │(v2) │ ← Deploy here
   │LIVE │       │IDLE │
   └─────┘       └─────┘
        │
        │  Switch
        ▼
   ┌─────┐       ┌─────┐
   │Blue │       │Green│
   │(v1) │       │(v2) │
   │IDLE │       │LIVE │ ← Now serving
   └─────┘       └─────┘
```

**Implementation:**
```yaml
# blue-deployment.yaml (existing)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: payment-blue
  labels:
    app: payment
    version: blue
spec:
  replicas: 3
  selector:
    matchLabels:
      app: payment
      version: blue
  template:
    metadata:
      labels:
        app: payment
        version: blue
    spec:
      containers:
      - name: api
        image: ghcr.io/you/module:v1.0.0
---
# green-deployment.yaml (new)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: payment-green
  labels:
    app: payment
    version: green
spec:
  replicas: 3
  selector:
    matchLabels:
      app: payment
      version: green
  template:
    metadata:
      labels:
        app: payment
        version: green
    spec:
      containers:
      - name: api
        image: ghcr.io/you/module:v1.1.0
---
# service.yaml (switch via selector)
apiVersion: v1
kind: Service
metadata:
  name: payment
spec:
  selector:
    app: payment
    version: blue  # ← เปลี่ยนเป็น green เมื่อพร้อม
  ports:
    - port: 80
      targetPort: 8080
```

**Switch:**
```bash
# Deploy green
kubectl apply -f green-deployment.yaml

# Wait ready
kubectl wait --for=condition=available --timeout=300s \
    deployment/payment-green

# Switch traffic
kubectl patch service payment -p \
    '{"spec":{"selector":{"version":"green"}}}'

# Verify
curl https://api.example.com/health

# Rollback if needed
kubectl patch service payment -p \
    '{"spec":{"selector":{"version":"blue"}}}'
```

## 22.5 Canary

```
Step 1: 5% canary
┌──────────────────┐
│   Load Balancer  │
└────┬────────┬────┘
   5%│       │95%
     ▼       ▼
  ┌───┐   ┌───┐
  │v2 │   │v1 │
  └───┘   └───┘

Step 2: 25% canary (if OK)
Step 3: 50% canary
Step 4: 100% (v2 becomes v1)
```

**Implementation with Argo Rollouts:**
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: payment-api
spec:
  replicas: 10
  strategy:
    canary:
      steps:
      - setWeight: 5
      - pause: { duration: 5m }
      - setWeight: 25
      - pause: { duration: 10m }
      - setWeight: 50
      - pause: { duration: 15m }
      - setWeight: 100
      canaryMetadata:
        labels:
          role: canary
      stableMetadata:
        labels:
          role: stable
      analysis:
        templates:
        - templateName: success-rate
        startingStep: 1
        args:
        - name: service-name
          value: payment-api
---
apiVersion: argoproj.io/v1alpha1
kind: AnalysisTemplate
metadata:
  name: success-rate
spec:
  metrics:
  - name: success-rate
    interval: 1m
    successCondition: result[0] >= 0.99
    failureLimit: 3
    provider:
      prometheus:
        address: http://prometheus:9090
        query: |
          sum(rate(http_requests_total{
            service="{{args.service-name}}",
            status=~"2.."
          }[5m]))
          /
          sum(rate(http_requests_total{
            service="{{args.service-name}}"
          }[5m]))
```

## 22.6 Shadow (Dark Launch)

```
User → LB → [v1] → Response
         │
         └──→ [v2] (response discarded)
              ├── Compare responses
              └── Log differences
```

**ใช้เมื่อ:**
- ต้องการ verify v2 โดยไม่กระทบ user
- Test migration
- Compare performance

---

# บทที่ 23: Blue-Green Deployment

## 23.1 รายละเอียด

**เมื่อใช้:**
- ต้องการ zero downtime
- ต้องการ instant rollback
- มี resource เพียงพอ (2x)

**ข้อดี:**
- Rollback ทันที (switch selector)
- Test v2 ได้เต็มที่ก่อน switch
- Simple mental model

**ข้อเสีย:**
- ใช้ resource 2x
- Database migration ซับซ้อน

## 23.2 Database Consideration

```
Problem:
- v1 ใช้ schema เดิม
- v2 ใช้ schema ใหม่
- ทั้ง 2 รันพร้อมกัน

Solution: Expand-Contract
Phase 1: Add new column (both v1, v2 ใช้ได้)
Phase 2: Deploy v2 (dual write)
Phase 3: Backfill
Phase 4: Switch read
Phase 5: Remove v1
Phase 6: Drop old column
```

## 23.3 Automation Script

```bash
#!/bin/bash
# scripts/blue-green-deploy.sh

set -euo pipefail

NEW_VERSION="$1"
CURRENT=$(kubectl get service payment -o jsonpath='{.spec.selector.version}')

echo "Current: $CURRENT"
echo "Deploying: $NEW_VERSION"

# Determine target color
if [ "$CURRENT" = "blue" ]; then
    TARGET="green"
else
    TARGET="blue"
fi

echo "Target: $TARGET"

# 1. Deploy new version to target
sed "s/VERSION_PLACEHOLDER/$NEW_VERSION/g" \
    k8s/deployment-template.yaml | \
    sed "s/COLOR_PLACEHOLDER/$TARGET/g" | \
    kubectl apply -f -

# 2. Wait ready
kubectl rollout status deployment/payment-$TARGET --timeout=5m

# 3. Smoke test against target
TARGET_IP=$(kubectl get service payment-$TARGET -o jsonpath='{.spec.clusterIP}')
if ! curl -sf "http://$TARGET_IP:8080/health"; then
    echo "❌ Smoke test failed"
    kubectl delete deployment/payment-$TARGET
    exit 1
fi

echo "✅ Smoke test passed"

# 4. Switch traffic
kubectl patch service payment -p \
    "{\"spec\":{\"selector\":{\"version\":\"$TARGET\"}}}"

echo "✅ Traffic switched to $TARGET"

# 5. Monitor for 5 minutes
sleep 300
if ! curl -sf "https://api.example.com/health"; then
    echo "❌ Post-switch health check failed, rolling back"
    kubectl patch service payment -p \
        "{\"spec\":{\"selector\":{\"version\":\"$CURRENT\"}}}"
    exit 1
fi

# 6. Keep old deployment for rollback (delete after 24h)
echo "✅ Deployment complete. Old version kept for rollback."
```

---

# บทที่ 24: Canary Deployment

## 24.1 Progressive Delivery

```
Phase 1: 1% traffic
         ↓ (5 min, metrics OK)
Phase 2: 5%
         ↓ (10 min)
Phase 3: 25%
         ↓ (15 min)
Phase 4: 50%
         ↓ (30 min)
Phase 5: 100%
```

## 24.2 Metrics for Canary Decision

| Metric | Threshold | Action |
|---|---|---|
| **Error rate** | > 1% | Rollback |
| **P95 latency** | > 500ms | Rollback |
| **CPU** | > 80% | Pause |
| **Memory** | > 90% | Pause |
| **Custom: payment success** | < 99% | Rollback |

## 24.3 Flagger (Alternative to Argo)

```yaml
apiVersion: flagger.app/v1beta1
kind: Canary
metadata:
  name: payment-api
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: payment-api
  
  service:
    port: 80
    targetPort: 8080
  
  analysis:
    interval: 1m
    threshold: 5
    maxWeight: 50
    stepWeight: 10
    
    metrics:
    - name: request-success-rate
      thresholdRange:
        min: 99
      interval: 1m
    
    - name: request-duration
      thresholdRange:
        max: 500
      interval: 1m
    
    webhooks:
    - name: load-test
      url: http://flagger-loadtester.test/
      timeout: 5s
      metadata:
        cmd: "hey -z 1m -q 10 -c 2 http://payment-api-canary/"
```

---

# บทที่ 25: Kubernetes Deployment

## 25.1 Full Manifest

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: payment-api
  namespace: production
  labels:
    app: payment-api
    version: v1.0.0
spec:
  replicas: 3
  revisionHistoryLimit: 5
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: payment-api
  template:
    metadata:
      labels:
        app: payment-api
        version: v1.0.0
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
        prometheus.io/path: "/metrics"
    spec:
      serviceAccountName: payment-api
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        fsGroup: 1000
      
      containers:
      - name: api
        image: ghcr.io/you/module:v1.0.0
        imagePullPolicy: IfNotPresent
        
        ports:
        - name: http
          containerPort: 8080
          protocol: TCP
        - name: metrics
          containerPort: 9090
          protocol: TCP
        
        env:
        - name: ENV
          value: "production"
        - name: LOG_LEVEL
          value: "info"
        - name: DB_DSN
          valueFrom:
            secretKeyRef:
              name: payment-secrets
              key: db-dsn
        - name: REDIS_ADDR
          value: "redis-master.cache.svc.cluster.local:6379"
        
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 512Mi
        
        livenessProbe:
          httpGet:
            path: /health/live
            port: http
          initialDelaySeconds: 10
          periodSeconds: 30
          timeoutSeconds: 3
          failureThreshold: 3
        
        readinessProbe:
          httpGet:
            path: /health/ready
            port: http
          initialDelaySeconds: 5
          periodSeconds: 10
          timeoutSeconds: 3
          failureThreshold: 2
        
        startupProbe:
          httpGet:
            path: /health/startup
            port: http
          initialDelaySeconds: 0
          periodSeconds: 5
          failureThreshold: 30
        
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities:
            drop:
              - ALL
        
        volumeMounts:
        - name: tmp
          mountPath: /tmp
        - name: cache
          mountPath: /home/app/.cache
      
      volumes:
      - name: tmp
        emptyDir: {}
      - name: cache
        emptyDir: {}
      
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchLabels:
                  app: payment-api
              topologyKey: kubernetes.io/hostname
      
      topologySpreadConstraints:
      - maxSkew: 1
        topologyKey: topology.kubernetes.io/zone
        whenUnsatisfiable: ScheduleAnyway
        labelSelector:
          matchLabels:
            app: payment-api
```

## 25.2 Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: payment-api
  namespace: production
spec:
  type: ClusterIP
  selector:
    app: payment-api
  ports:
  - name: http
    port: 80
    targetPort: http
    protocol: TCP
  sessionAffinity: None
```

## 25.3 Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: payment-api
  namespace: production
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/rate-limit: "100"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - api.example.com
    secretName: payment-tls
  rules:
  - host: api.example.com
    http:
      paths:
      - path: /api/v1/payments
        pathType: Prefix
        backend:
          service:
            name: payment-api
            port:
              name: http
```

## 25.4 HPA (Autoscaling)

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: payment-api
  namespace: production
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: payment-api
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  - type: Pods
    pods:
      metric:
        name: http_requests_per_second
      target:
        type: AverageValue
        averageValue: "1000"
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 10
        periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
      - type: Percent
        value: 100
        periodSeconds: 15
      - type: Pods
        value: 4
        periodSeconds: 15
      selectPolicy: Max
```

## 25.5 PodDisruptionBudget

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: payment-api
  namespace: production
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app: payment-api
```

## 25.6 ConfigMap & Secret

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: payment-config
  namespace: production
data:
  LOG_LEVEL: "info"
  RETRY_COUNT: "3"
---
apiVersion: v1
kind: Secret
metadata:
  name: payment-secrets
  namespace: production
type: Opaque
stringData:
  db-dsn: "postgres://user:pass@postgres:5432/payments"
  jwt-secret: "your-secret"
```

**⚠️ คำเตือน:** อย่า commit secret ใน Git ใช้ Sealed Secrets หรือ External Secrets Operator

---

# บทที่ 26: Zero-Downtime Migration

## 26.1 หลักการ

```
Deploy ≠ Just code
Deploy = Code + Schema + Data + Config
```

**ทุกอย่างต้อง backward compatible**

## 26.2 Migration Order

```
1. Database (backward compatible)
   ├── Add columns (nullable)
   ├── Add indexes (CONCURRENTLY)
   └── Add tables

2. Code (reads both old + new)
   ├── Deploy new code
   ├── Dual write
   └── Old code still works

3. Verify
   ├── Monitor metrics
   ├── Check data consistency
   └── Run smoke test

4. Migrate data (batch)

5. Switch (read new)

6. Cleanup (after 2 sprints)
```

## 26.3 Expand-Contract for K8s

```yaml
# Phase 1: Deploy new version (backward compatible)
kubectl apply -f deployment-v2.yaml

# Phase 2: Wait for rollout
kubectl rollout status deployment/payment-api

# Phase 3: Run migration
kubectl exec -it deployment/payment-api -- /app/migrate up

# Phase 4: Verify
kubectl exec -it deployment/payment-api -- /app/migrate verify

# Phase 5: (Optional) Remove old code paths
# (after 2 sprints)
```

## 26.4 Pre-Deploy Hook (Migration Job)

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: payment-migrate-{{ .Release.Revision }}
  annotations:
    "helm.sh/hook": pre-upgrade
    "helm.sh/hook-delete-policy": before-hook-creation
spec:
  backoffLimit: 3
  ttlSecondsAfterFinished: 300
  template:
    spec:
      restartPolicy: Never
      containers:
      - name: migrate
        image: ghcr.io/you/module:v1.0.0
        command: ["/app/migrate"]
        args: ["up"]
        env:
        - name: DB_DSN
          valueFrom:
            secretKeyRef:
              name: payment-secrets
              key: db-dsn
```

## 26.5 Health Endpoints

```go
// Three health endpoints
func setupHealth(r chi.Router, db *gorm.DB, redis *redis.Client) {
    // Liveness — is the process alive?
    r.Get("/health/live", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })
    
    // Readiness — can it serve traffic?
    r.Get("/health/ready", func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
        defer cancel()
        
        // Check DB
        if err := db.WithContext(ctx).Exec("SELECT 1").Error; err != nil {
            http.Error(w, "DB not ready", http.StatusServiceUnavailable)
            return
        }
        
        // Check Redis
        if err := redis.Ping(ctx).Err(); err != nil {
            http.Error(w, "Redis not ready", http.StatusServiceUnavailable)
            return
        }
        
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Ready"))
    })
    
    // Startup — has the app finished starting?
    r.Get("/health/startup", func(w http.ResponseWriter, r *http.Request) {
        if !appStarted.Load() {
            http.Error(w, "Starting", http.StatusServiceUnavailable)
            return
        }
        w.WriteHeader(http.StatusOK)
    })
}
```

## 26.6 Graceful Shutdown

```go
func main() {
    srv := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }
    
    // Start server
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    
    // Wait for signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Info("shutting down gracefully")
    
    // 1. Stop accepting new connections
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        log.Error("shutdown error", "error", err)
    }
    
    // 2. Close DB
    if sqlDB, err := db.DB(); err == nil {
        sqlDB.Close()
    }
    
    // 3. Close Redis
    redis.Close()
    
    // 4. Flush logs
    logger.Sync()
    
    log.Info("shutdown complete")
}
```

---

# บทที่ 27: Rollback Strategy

## 27.1 Rollback Levels

| Level | Method | Time | Data Impact |
|---|---|---|---|
| **Instant** | Feature flag | < 30s | None |
| **Fast** | Deployment rollback | 2-5 min | In-flight requests |
| **Medium** | Migration rollback | 10-60 min | New data lost |
| **Slow** | Restore backup | 1-4 hours | Data loss |

## 27.2 Deployment Rollback

```bash
# Kubernetes
kubectl rollout undo deployment/payment-api

# Rollback to specific revision
kubectl rollout history deployment/payment-api
kubectl rollout undo deployment/payment-api --to-revision=3

# Verify
kubectl rollout status deployment/payment-api
```

## 27.3 Helm Rollback

```bash
# List revisions
helm history payment-api

# Rollback
helm rollback payment-api 3

# Verify
helm status payment-api
```

## 27.4 Auto-Rollback

```yaml
# Argo Rollouts auto-rollback
apiVersion: argoproj.io/v1alpha1
kind: Rollout
spec:
  strategy:
    canary:
      analysis:
        templates:
        - templateName: success-rate
      autoPromotionEnabled: false
      abortScaleDownDelaySeconds: 30
```

```yaml
# Flagger
apiVersion: flagger.app/v1beta1
kind: Canary
spec:
  analysis:
    threshold: 5  # 5 failed checks → rollback
```

## 27.5 Rollback Checklist

- [ ] Detect issue (alert, monitoring)
- [ ] Confirm severity (P0/P1?)
- [ ] Notify team (Slack, PagerDuty)
- [ ] Decide: rollback or forward-fix?
- [ ] Execute rollback
- [ ] Verify service restored
- [ ] Monitor 30 min
- [ ] Postmortem scheduled
- [ ] Root cause fix planned

## 27.6 Feature Flag Rollback

```bash
# Instant rollback via flag
curl -X POST https://flags.example.com/api/flags/payment.new_rule \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"enabled": false}'

# Verify
curl https://flags.example.com/api/flags/payment.new_rule
# {"name":"payment.new_rule","enabled":false}
```

## 27.7 Data Rollback

**ถ้า migration ผิด:**

```sql
-- Before migration
CREATE TABLE payment_transactions_backup_20260415
    AS SELECT * FROM payment_transactions;

-- Migration runs
ALTER TABLE payment_transactions ADD COLUMN new_col TYPE;

-- Rollback if needed
TRUNCATE payment_transactions;
INSERT INTO payment_transactions
    SELECT * FROM payment_transactions_backup_20260415;
```

**⚠️ คำเตือน:** Backup ก่อน migration ทุกครั้ง

---

# บทที่ 28: Three Pillars (Metrics, Logs, Traces)

## 28.1 ภาพรวม

```
                     ┌─────────────────┐
                     │  Observability  │
                     └────────┬────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
   ┌─────────┐          ┌─────────┐          ┌─────────┐
   │ Metrics │          │  Logs   │          │ Traces  │
   │         │          │         │          │         │
   │ Numbers │          │ Events  │          │ Flows   │
   │   over  │          │   over  │          │  of     │
   │  time   │          │  time   │          │ request │
   └─────────┘          └─────────┘          └─────────┘
        │                     │                     │
        ▼                     ▼                     ▼
   Prometheus              Loki                 Jaeger
   Grafana                 ELK                  Tempo
   Datadog                 Splunk               Zipkin
```

## 28.2 เมื่อไหร่ใช้แต่ละ Pillar

| Pillar | คำถามที่ตอบ |
|---|---|
| **Metrics** | ระบบสุขภาพเป็นยังไง? (aggregate) |
| **Logs** | เกิดอะไรขึ้นตอนนั้น? (specific) |
| **Traces** | Request ไปไหนบ้าง? (distributed) |

## 28.3 Integration

```go
// Instrument ทั้ง 3 pillars ในที่เดียว
type Observability struct {
    metrics  *prometheus.Registry
    tracer   trace.Tracer
    logger   *slog.Logger
}

func (o *Observability) WrapHandler(name string, h http.Handler) http.Handler {
    // Metrics
    requestDuration := prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration",
        },
        []string{"method", "path", "status"},
    )
    
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Trace
        ctx, span := o.tracer.Start(r.Context(), name)
        defer span.End()
        r = r.WithContext(ctx)
        
        // Log
        logger := o.logger.With(
            "trace_id", span.SpanContext().TraceID().String(),
            "method", r.Method,
            "path", r.URL.Path,
        )
        
        // Serve
        rw := &responseWriter{ResponseWriter: w}
        h.ServeHTTP(rw, r)
        
        // Record
        duration := time.Since(start).Seconds()
        requestDuration.WithLabelValues(r.Method, r.URL.Path, 
            strconv.Itoa(rw.status)).Observe(duration)
        
        logger.Info("request completed",
            "status", rw.status,
            "duration_ms", duration*1000,
        )
    })
}
```

---

# บทที่ 29: Prometheus & Grafana

## 29.1 Metrics Instrumentation

```go
// pkg/metrics/metrics.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // HTTP metrics
    HTTPRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    HTTPRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration",
            Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
        },
        []string{"method", "path"},
    )
    
    // Business metrics
    PaymentsCreated = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "payment_created_total",
            Help: "Total payments created",
        },
        []string{"method", "currency", "status"},
    )
    
    PaymentAmount = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "payment_amount",
            Help:    "Payment amounts",
            Buckets: prometheus.ExponentialBuckets(10, 10, 6),
        },
        []string{"currency"},
    )
    
    // Database metrics
    DBQueryDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "db_query_duration_seconds",
            Help:    "Database query duration",
            Buckets: prometheus.DefBuckets,
        },
        []string{"operation", "table"},
    )
    
    // Runtime metrics (built-in)
    // go_goroutines, go_memstats_*, etc.
)
```

## 29.2 Usage

```go
func (uc *CreatePaymentUseCase) Execute(ctx context.Context, input CreatePaymentInput) (*Output, error) {
    start := time.Now()
    
    // ... business logic
    
    // Record metrics
    metrics.PaymentsCreated.WithLabelValues(
        input.Method, input.Currency, string(payment.Status),
    ).Inc()
    
    metrics.PaymentAmount.WithLabelValues(input.Currency).
        Observe(payment.Amount.InexactFloat64())
    
    metrics.DBQueryDuration.WithLabelValues("insert", "payments").
        Observe(time.Since(start).Seconds())
    
    return output, nil
}
```

## 29.3 Metrics Endpoint

```go
import (
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    r := chi.NewRouter()
    
    // Health
    r.Get("/health", healthHandler)
    
    // Metrics (Prometheus)
    r.Handle("/metrics", promhttp.Handler())
    
    // ...
}
```

## 29.4 Prometheus Config

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'payment-api'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: true
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_port]
        action: replace
        target_label: __address__
        regex: (.+)
        replacement: ${1}
```

## 29.5 Grafana Dashboard

```json
{
  "dashboard": {
    "title": "Payment Module",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total{job=\"payment-api\"}[5m])) by (status)",
            "legendFormat": "{{status}}"
          }
        ]
      },
      {
        "title": "P95 Latency",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))"
          }
        ]
      },
      {
        "title": "Payments/sec",
        "type": "graph",
        "targets": [
          {
            "expr": "sum(rate(payment_created_total[5m])) by (status)"
          }
        ]
      },
      {
        "title": "Error Rate",
        "type": "singlestat",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total{status=~\"5..\"}[5m])) / sum(rate(http_requests_total[5m]))"
          }
        ],
        "thresholds": "0.001,0.01"
      }
    ]
  }
}
```

## 29.6 Key Metrics to Monitor

| Metric | Formula | Target |
|---|---|---|
| **Request Rate** | `sum(rate(http_requests_total[5m]))` | – |
| **Error Rate** | `sum(rate(http_requests_total{status=~"5.."}[5m]))` | < 0.1% |
| **P95 Latency** | `histogram_quantile(0.95, ...)` | < 200ms |
| **P99 Latency** | `histogram_quantile(0.99, ...)` | < 500ms |
| **Throughput** | `sum(rate(payment_created_total[5m]))` | – |
| **Saturation** | `rate(process_cpu_seconds_total[5m])` | < 80% |
| **Memory** | `process_resident_memory_bytes` | < 512 MB |

## 29.7 Alert Rules

```yaml
# prometheus-rules.yml
groups:
  - name: payment
    interval: 30s
    rules:
      - alert: HighErrorRate
        expr: |
          sum(rate(http_requests_total{job="payment-api",status=~"5.."}[5m]))
          /
          sum(rate(http_requests_total{job="payment-api"}[5m]))
          > 0.01
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Error rate > 1%"
          description: "Current: {{ $value | humanizePercentage }}"
      
      - alert: HighLatency
        expr: |
          histogram_quantile(0.95,
            sum(rate(http_request_duration_seconds_bucket{job="payment-api"}[5m]))
            by (le)
          ) > 0.5
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "P95 latency > 500ms"
      
      - alert: HighMemory
        expr: |
          process_resident_memory_bytes{job="payment-api"} > 512 * 1024 * 1024
        for: 10m
        labels:
          severity: warning
      
      - alert: PodDown
        expr: |
          up{job="payment-api"} == 0
        for: 1m
        labels:
          severity: critical
```

---

# บทที่ 30: Structured Logging

## 30.1 Principle

**Structured Logging** = JSON logs ที่ parse ได้

```go
// ❌ Bad — string log
log.Printf("User %s created payment %s amount %f", userID, paymentID, amount)

// ✅ Good — structured
logger.Info("payment created",
    "user_id", userID,
    "payment_id", paymentID,
    "amount", amount,
    "currency", "THB",
)
```

**ผลลัพธ์:**
```json
{"time":"2026-04-15T10:30:00Z","level":"INFO","msg":"payment created",
 "user_id":"...","payment_id":"...","amount":100,"currency":"THB"}
```

## 30.2 Implementation (slog)

```go
// pkg/logger/logger.go
package logger

import (
    "context"
    "log/slog"
    "os"
)

type Logger interface {
    Debug(msg string, args ...any)
    Info(msg string, args ...any)
    Warn(msg string, args ...any)
    Error(msg string, args ...any)
    With(args ...any) Logger
    WithContext(ctx context.Context) Logger
}

type slogLogger struct {
    l *slog.Logger
}

func New() Logger {
    var handler slog.Handler
    
    env := os.Getenv("ENV")
    if env == "production" {
        // JSON logs in production
        handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
            Level: slog.LevelInfo,
            AddSource: true,
        })
    } else {
        // Text logs in dev (easier to read)
        handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
            Level: slog.LevelDebug,
        })
    }
    
    return &slogLogger{l: slog.New(handler)}
}

func (s *slogLogger) Info(msg string, args ...any) {
    s.l.Info(msg, args...)
}

// ...

func (s *slogLogger) WithContext(ctx context.Context) Logger {
    // Extract trace_id, request_id จาก context
    args := []any{}
    
    if traceID := ctx.Value("trace_id"); traceID != nil {
        args = append(args, "trace_id", traceID)
    }
    if requestID := ctx.Value("request_id"); requestID != nil {
        args = append(args, "request_id", requestID)
    }
    
    return &slogLogger{l: s.l.With(args...)}
}
```

## 30.3 Usage

```go
func (uc *CreatePaymentUseCase) Execute(ctx context.Context, input CreatePaymentInput) (*Output, error) {
    logger := uc.logger.WithContext(ctx)
    
    logger.Info("creating payment",
        "user_id", input.UserID,
        "order_id", input.OrderID,
        "amount", input.Amount.String(),
    )
    
    payment, err := entity.NewPayment(...)
    if err != nil {
        logger.Error("invalid payment input", "error", err)
        return nil, err
    }
    
    if err := uc.repo.Save(ctx, payment); err != nil {
        logger.Error("failed to save payment",
            "error", err,
            "payment_id", payment.ID,
        )
        return nil, err
    }
    
    logger.Info("payment created",
        "payment_id", payment.ID,
        "status", payment.Status,
    )
    
    return &Output{ID: payment.ID}, nil
}
```

## 30.4 Log Aggregation (Loki)

```yaml
# docker-compose.yml
services:
  loki:
    image: grafana/loki:2.9.0
    ports:
      - "3100:3100"
    command: -config.file=/etc/loki/local-config.yaml
  
  promtail:
    image: grafana/promtail:2.9.0
    volumes:
      - /var/log:/var/log:ro
      - ./promtail.yml:/etc/promtail/config.yml
    command: -config.file=/etc/promtail/config.yml
  
  grafana:
    image: grafana/grafana:10.2.0
    ports:
      - "3000:3000"
```

**Query ใน Grafana:**
```logql
{app="payment-api"} |= "error"
{app="payment-api"} | json | level="ERROR"
{app="payment-api"} | json | status >= 500
```

## 30.5 Best Practices

### ✅ DO
- ใช้ structured (JSON)
- Include context (user_id, request_id)
- Log at boundaries (start, end, error)
- ใช้ level ถูก (Debug/Info/Warn/Error)

### ❌ DON'T
- Log sensitive data (password, token, PII)
- Log ใน hot loop (performance)
- ใช้ `fmt.Println` (ไม่มี structure)
- Log ข้อมูลมากเกินไป (cost)

## 30.6 Log Levels

| Level | ใช้เมื่อ | ตัวอย่าง |
|---|---|---|
| **Debug** | Dev, detailed | "Query: SELECT..." |
| **Info** | Normal operation | "Payment created" |
| **Warn** | Recoverable issues | "Retry attempt 2" |
| **Error** | Failures | "DB connection failed" |

---

# บทที่ 31: Distributed Tracing

## 31.1 OpenTelemetry Setup

```go
// pkg/tracing/tracing.go
package tracing

import (
    "context"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func Init(ctx context.Context, serviceName string) (func(context.Context) error, error) {
    // Exporter (Jaeger, Tempo, etc.)
    exporter, err := otlptracegrpc.New(ctx,
        otlptracegrpc.WithEndpoint("localhost:4317"),
        otlptracegrpc.WithInsecure(),
    )
    if err != nil {
        return nil, err
    }
    
    // Resource
    res, err := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceName(serviceName),
            semconv.ServiceVersion("1.0.0"),
        ),
    )
    if err != nil {
        return nil, err
    }
    
    // Tracer Provider
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(res),
        sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.1))),
    )
    
    otel.SetTracerProvider(tp)
    otel.SetTextMapPropagator(propagation.TraceContext{})
    
    return tp.Shutdown, nil
}
```

## 31.2 Manual Instrumentation

```go
func (uc *CreatePaymentUseCase) Execute(ctx context.Context, input CreatePaymentInput) (*Output, error) {
    tracer := otel.Tracer("payment")
    
    ctx, span := tracer.Start(ctx, "CreatePayment")
    defer span.End()
    
    span.SetAttributes(
        attribute.String("user_id", input.UserID.String()),
        attribute.String("order_id", input.OrderID.String()),
        attribute.String("amount", input.Amount.String()),
    )
    
    // ... business logic
    
    // Sub-span for DB
    ctx, dbSpan := tracer.Start(ctx, "SavePayment")
    err := uc.repo.Save(ctx, payment)
    dbSpan.End()
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }
    
    return &Output{ID: payment.ID}, nil
}
```

## 31.3 HTTP Middleware

```go
func TracingMiddleware(serviceName string) func(http.Handler) http.Handler {
    tracer := otel.Tracer(serviceName)
    propagator := otel.GetTextMapPropagator()
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract trace context from headers
            ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
            
            // Start span
            ctx, span := tracer.Start(ctx, r.Method+" "+r.URL.Path,
                trace.WithSpanKind(trace.SpanKindServer),
                trace.WithAttributes(
                    semconv.HTTPMethod(r.Method),
                    semconv.HTTPURL(r.URL.String()),
                    semconv.HTTPUserAgent(r.UserAgent()),
                    semconv.HTTPClientIP(r.RemoteAddr),
                ),
            )
            defer span.End()
            
            // Wrap response
            rw := &responseWriter{ResponseWriter: w}
            
            // Serve
            next.ServeHTTP(rw, r.WithContext(ctx))
            
            // Record status
            span.SetAttributes(semconv.HTTPStatusCode(rw.status))
            if rw.status >= 500 {
                span.SetStatus(codes.Error, "server error")
            }
        })
    }
}
```

## 31.4 Trace Propagation

```
Service A → Service B
    │
    └── headers:
        traceparent: 00-<trace_id>-<span_id>-01
        tracestate: vendor=value
```

**Client:**
```go
func (c *Client) CallB(ctx context.Context) error {
    req, _ := http.NewRequestWithContext(ctx, "GET", "http://service-b/api", nil)
    
    // Inject trace context
    otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
    
    return c.http.Do(req)
}
```

## 31.5 Jaeger Setup

```yaml
# docker-compose.yml
services:
  jaeger:
    image: jaegertracing/all-in-one:1.52
    ports:
      - "16686:16686"  # UI
      - "4317:4317"    # OTLP gRPC
      - "4318:4318"    # OTLP HTTP
    environment:
      COLLECTOR_OTLP_ENABLED: "true"
```

**เปิด UI:** http://localhost:16686

---

# บทที่ 32: Alerting Strategy

## 32.1 Alert Principles

```
1. Alert on symptoms, not causes
   ❌ CPU high → alert
   ✅ User-facing latency high → alert

2. Every alert must be actionable
   ❌ "Memory > 80%" (ต้องทำอะไร?)
   ✅ "Memory > 90% for 10 min, OOM imminent"

3. Alert on SLOs, not infrastructure
   ❌ "CPU > 70%"
   ✅ "P95 latency > 200ms" (SLO)
```

## 32.2 Alert Severity

| Severity | Response | Channel |
|---|---|---|
| **Critical (P1)** | ทันที | Phone call, SMS |
| **High (P2)** | 15 นาที | Slack, Email |
| **Medium (P3)** | 1 ชั่วโมง | Slack |
| **Low (P4)** | Next day | Email |
| **Info** | Weekly | Dashboard |

## 32.3 Prometheus Alert Rules

```yaml
groups:
  - name: payment-slo
    rules:
      # SLO: 99.9% availability
      - alert: AvailabilitySLOBreach
        expr: |
          (
            sum(rate(http_requests_total{job="payment",status!~"5.."}[5m]))
            /
            sum(rate(http_requests_total{job="payment"}[5m]))
          ) < 0.999
        for: 5m
        labels:
          severity: critical
          slo: availability
        annotations:
          summary: "Availability SLO breach: {{ $value | humanizePercentage }}"
          runbook_url: "https://wiki/runbooks/availability"
      
      # SLO: P95 < 200ms
      - alert: LatencySLOBreach
        expr: |
          histogram_quantile(0.95,
            sum(rate(http_request_duration_seconds_bucket{job="payment"}[5m]))
            by (le)
          ) > 0.2
        for: 5m
        labels:
          severity: high
          slo: latency
        annotations:
          summary: "P95 latency: {{ $value }}s"
      
      # Error budget burn
      - alert: ErrorBudgetBurnFast
        expr: |
          (
            sum(rate(http_requests_total{job="payment",status=~"5.."}[1h]))
            /
            sum(rate(http_requests_total{job="payment"}[1h]))
          ) > (14.4 * 0.001)  # 2% budget in 1h
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "Fast error budget burn"
```

## 32.4 Alert Routing (Alertmanager)

```yaml
# alertmanager.yml
route:
  receiver: 'default'
  group_by: ['alertname', 'cluster', 'service']
  group_wait: 10s
  group_interval: 5m
  repeat_interval: 4h
  routes:
    - match:
        severity: critical
      receiver: 'pagerduty-critical'
      continue: true
    - match:
        severity: high
      receiver: 'slack-oncall'
    - match:
        severity: warning
      receiver: 'slack-warnings'

receivers:
  - name: 'pagerduty-critical'
    pagerduty_configs:
      - service_key: '<key>'
        severity: critical
  
  - name: 'slack-oncall'
    slack_configs:
      - api_url: '<webhook>'
        channel: '#alerts-oncall'
        title: '{{ .GroupLabels.alertname }}'
        text: '{{ range .Alerts }}{{ .Annotations.summary }}{{ end }}'
  
  - name: 'slack-warnings'
    slack_configs:
      - api_url: '<webhook>'
        channel: '#alerts-warnings'
```

## 32.5 On-Call Rotation

```yaml
# PagerDuty schedule
- Schedule: "payment-oncall"
- Rotation: Weekly
- Escalation:
  - Level 1: Primary (5 min)
  - Level 2: Secondary (10 min)
  - Level 3: Manager (30 min)
```

## 32.6 Alert Quality

**Metrics ต่อ alert:**
- **Precision** = alerts ที่ actionable / total alerts
- **Recall** = issues ที่ alert / total issues
- **MTTA** = mean time to acknowledge
- **MTTR** = mean time to resolve

**Target:**
- Precision > 80%
- Recall > 90%
- MTTA < 5 min
- MTTR < 1 hour (P1)

## 32.7 Anti-Patterns

| Anti-Pattern | ปัญหา | แก้ |
|---|---|---|
| **Alert fatigue** | เยอะเกิน, ignore | ลด, ปรับ threshold |
| **Non-actionable** | ไม่รู้จะทำอะไร | เพิ่ม runbook |
| **Flapping** | เปิด-ปิด ถี่ | ใช้ `for: 5m` |
| **Silent alert** | ไม่มีใครเห็น | Check routing |
| **No owner** | ไม่มีคนรับ | Assign on-call |

---

# บทที่ 33: Incident Response

## 33.1 Incident Severity

| Severity | Impact | Response Time | Example |
|---|---|---|---|
| **SEV-1** | Total outage | 15 min | Payment down 100% |
| **SEV-2** | Major degradation | 1 hour | Payment 50% failure |
| **SEV-3** | Minor impact | 4 hours | Slow for some users |
| **SEV-4** | Cosmetic | Next day | Typo in UI |

## 33.2 Incident Lifecycle

```
Detect → Triage → Mitigate → Resolve → Postmortem
  │        │         │          │           │
  ▼        ▼         ▼          ▼           ▼
 Alert    Assess    Fix        Verify     Learn
          Assign    (or         Stable     Action
          Severity  rollback)              items
```

## 33.3 Incident Roles

| Role | Responsibility |
|---|---|
| **Incident Commander (IC)** | Coordinate, decide, communicate |
| **Subject Matter Expert (SME)** | Technical investigation |
| **Communications Lead** | External comms, status page |
| **Scribe** | Document timeline |

## 33.4 Incident Commander Checklist

```markdown
## Initial (First 5 min)
- [ ] Acknowledge alert
- [ ] Create incident channel (#incident-2026-04-15)
- [ ] Declare severity (SEV-1/2/3/4)
- [ ] Assign roles
- [ ] Start timeline doc

## During (Every 15 min)
- [ ] Status update (Slack)
- [ ] External comms (if SEV-1)
- [ ] Check impact metrics
- [ ] Decide: mitigate vs fix-forward

## Resolution
- [ ] Confirm stable (30 min no alerts)
- [ ] Notify stakeholders
- [ ] Update status page
- [ ] Schedule postmortem (within 5 days for SEV-1)

## Post-Incident
- [ ] Write postmortem
- [ ] Track action items
- [ ] Share learnings
```

## 33.5 Status Page Update

```markdown
# Example Status Page Updates

## Investigating (T+0)
**Payment API Degradation**
We are investigating reports of elevated error rates in the payment API.
Started: 10:30 UTC

## Identified (T+15)
**Payment API Degradation**
We have identified the issue as a database connection pool exhaustion.
Started: 10:30 UTC | Updated: 10:45 UTC

## Monitoring (T+45)
**Payment API Degradation**
Fix deployed. Monitoring recovery.
Started: 10:30 UTC | Updated: 11:15 UTC

## Resolved (T+60)
**Payment API Degradation**
All systems operational. Duration: 60 min.
Started: 10:30 UTC | Resolved: 11:30 UTC
```

## 33.6 Communication Template

```markdown
# Slack #incident channel

@here 🚨 SEV-1 Incident Declared

**What:** Payment API returning 500 errors (100% failure rate)
**When:** Started 10:30 UTC
**Impact:** All payments blocked
**IC:** @alice
**SME:** @bob
**Status:** Investigating

**Channel:** #incident-2026-04-15
**Status Page:** https://status.example.com
**Bridge:** https://meet.example.com/incident

Next update: 10:45 UTC
```

## 33.7 Postmortem (Blameless)

```markdown
# Postmortem: Payment API Outage (2026-04-15)

## Summary
Payment API was down for 60 minutes on April 15, 2026, affecting
all payment transactions. Root cause: database connection pool
exhaustion due to connection leak in a new code path.

## Impact
- **Users:** ~5,000 affected
- **Transactions:** ~10,000 failed
- **Revenue:** ~$50,000 lost
- **Duration:** 60 minutes (10:30 - 11:30 UTC)

## Timeline (UTC)
| Time | Event |
|------|-------|
| 10:25 | Deploy v1.5.0 with new refund path |
| 10:28 | Deploy complete |
| 10:30 | Alert: error rate > 5% |
| 10:31 | Incident declared SEV-1 |
| 10:35 | DB connection pool exhausted |
| 10:40 | Rollback initiated |
| 10:45 | Rollback complete, still failing |
| 10:50 | Root cause: connection leak in refund |
| 11:00 | Fix deployed (restart pods) |
| 11:30 | All clear |

## Root Cause
New refund code path opened DB connections but didn't close them
on error. Under load, connections exhausted the pool.

```go
// Bug: leaked connection on error
func (r *repo) Refund(ctx context.Context, id uuid.UUID) error {
    tx := r.db.Begin()  // ← Begin
    defer tx.Commit()
    // Missing: tx.Rollback() on error
    if err := doSomething(ctx, tx); err != nil {
        return err  // ← Connection leaked
    }
    return tx.Commit()
}
```

## 5 Whys
1. **Why down?** Connection pool exhausted
2. **Why exhausted?** Connections leaked
3. **Why leaked?** Missing rollback on error path
4. **Why missing?** New code, no test for error path
5. **Why no test?** Test coverage gap for error handling

## Detection
- **Detected by:** Prometheus alert (error rate > 1%)
- **MTTD:** 5 minutes
- **MTTA:** 6 minutes
- **MTTR:** 60 minutes

## What Went Well
- Alert fired quickly
- Team responded fast
- Rollback plan existed
- Communication clear

## What Went Wrong
- Missing error path test
- No connection pool monitoring
- Code review missed the issue
- No canary deployment

## Action Items

| # | Action | Owner | Due | Priority |
|---|--------|-------|-----|----------|
| 1 | Fix connection leak | @bob | 2026-04-16 | P0 |
| 2 | Add error path tests | @carol | 2026-04-17 | P0 |
| 3 | Monitor DB pool metrics | @dave | 2026-04-18 | P1 |
| 4 | Add connection pool alert | @dave | 2026-04-18 | P1 |
| 5 | Enforce canary deploy | @alice | 2026-04-20 | P1 |
| 6 | Add code review checklist | @team | 2026-04-22 | P2 |

## Lessons Learned
- Error paths ต้องมี test เสมอ
- Connection pool ต้อง monitor
- Canary deployment ต้องบังคับ
- Code review checklist ต้องมี transaction handling

## Appendix
- [Dashboard](https://grafana.example.com/d/payment)
- [Logs](https://loki.example.com/explore)
- [Traces](https://jaeger.example.com/trace/xxx)
```

---

# บทที่ 34: Runbooks

## 34.1 Runbook Template

```markdown
# Runbook: [Issue Name]

## Alert
`AlertName`

## Symptoms
- What users see
- What metrics show
- What logs show

## Quick Diagnosis
\`\`\`bash
# Check 1: Pods status
kubectl get pods -l app=payment-api

# Check 2: Error rate
curl -s 'http://prometheus:9090/api/v1/query?query=rate(http_requests_total{status=~"5.."}[5m])'

# Check 3: Logs
kubectl logs -l app=payment-api --tail=100 | grep ERROR
\`\`\`

## Resolution Steps

### Step 1: Identify issue
- Check X
- Check Y

### Step 2: Mitigate
- Option A: [if scenario A] → ...
- Option B: [if scenario B] → ...

### Step 3: Verify
- Metrics back to normal
- No new errors in logs
- Smoke test passes

## Rollback
\`\`\`bash
kubectl rollout undo deployment/payment-api
\`\`\`

## Escalation
- If not resolved in 30 min → @team-lead
- If not resolved in 1 hour → @manager

## Related
- [Dashboard](...)
- [Postmortem template](...)
- [Incident response](...)
```

## 34.2 Common Runbooks

### 34.2.1 High Error Rate

```markdown
# Runbook: High Error Rate

## Alert
`HighErrorRate` (error rate > 1% for 5 min)

## Quick Diagnosis
\`\`\`bash
# 1. Check error rate by status
curl -s 'http://prometheus:9090/api/v1/query?query=sum(rate(http_requests_total{status=~"5.."}[5m])) by (status)'

# 2. Check recent logs
kubectl logs -l app=payment-api --tail=100 --since=10m | grep -i error

# 3. Check recent deployments
kubectl rollout history deployment/payment-api
\`\`\`

## Likely Causes
1. Recent deployment bug → check rollout history
2. DB issues → check DB connections
3. External dependency down → check provider status
4. Resource exhaustion → check CPU/memory

## Mitigation

### If recent deploy:
\`\`\`bash
kubectl rollout undo deployment/payment-api
kubectl rollout status deployment/payment-api
\`\`\`

### If DB issue:
\`\`\`bash
# Check connections
kubectl exec -it deployment/payment-api -- \
    psql $DB_DSN -c "SELECT count(*) FROM pg_stat_activity;"
\`\`\`

### If resource issue:
\`\`\`bash
kubectl top pods -l app=payment-api
kubectl scale deployment payment-api --replicas=10
\`\`\`
```

### 34.2.2 High Latency

```markdown
# Runbook: High Latency

## Alert
`HighLatency` (P95 > 500ms for 5 min)

## Quick Diagnosis
\`\`\`bash
# 1. Latency by endpoint
curl -s 'http://prometheus:9090/api/v1/query?query=histogram_quantile(0.95,sum(rate(http_request_duration_seconds_bucket[5m])) by (le,path))'

# 2. DB slow queries
kubectl exec -it postgres-0 -- \
    psql -c "SELECT query, mean_exec_time FROM pg_stat_statements ORDER BY mean_exec_time DESC LIMIT 10;"

# 3. Traces
# Open Jaeger: filter by service=payment-api, duration>500ms
\`\`\`

## Likely Causes
1. DB slow query
2. Missing index
3. N+1 query
4. External API slow
5. Lock contention

## Mitigation

### If DB slow:
\`\`\`sql
-- Find slow query
SELECT * FROM pg_stat_activity WHERE state = 'active' AND query_start < NOW() - INTERVAL '10 seconds';

-- Kill long-running query
SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE pid = XXX;
\`\`\`

### If missing index:
\`\`\`sql
CREATE INDEX CONCURRENTLY idx_payment_user_id ON payment_transactions(user_id);
\`\`\`
```

### 34.2.3 Database Connection Pool Exhausted

```markdown
# Runbook: DB Connection Pool Exhausted

## Alert
`DBPoolExhausted` (pool usage > 90% for 5 min)

## Quick Diagnosis
\`\`\`bash
# Check pool
kubectl exec -it deployment/payment-api -- \
    curl -s localhost:8080/metrics | grep db_connection

# Check DB connections
kubectl exec -it postgres-0 -- \
    psql -c "SELECT count(*), state FROM pg_stat_activity GROUP BY state;"
\`\`\`

## Immediate Mitigation
1. **Increase pool size** (temporary)
   \`\`\`bash
   kubectl set env deployment/payment-api DB_MAX_OPEN=50
   \`\`\`

2. **Restart pods** (if leak suspected)
   \`\`\`bash
   kubectl rollout restart deployment/payment-api
   \`\`\`

3. **Scale replicas**
   \`\`\`bash
   kubectl scale deployment payment-api --replicas=10
   \`\`\`

## Root Cause Investigation
- Check for connection leak in recent code
- Check for long-running transactions
- Check for missing Rollback in error paths
```

---

# บทที่ 35: Postmortem Culture

## 35.1 Blameless Principles

```
1. Focus on systems, not people
   ❌ "Alice made a mistake"
   ✅ "The deploy process allowed this"

2. Assume good intent
   ❌ "Bob was careless"
   ✅ "Bob made the best decision with available info"

3. Learn, don't punish
   ❌ "Who broke it?"
   ✅ "What can we improve?"

4. Action items, not blame
   ❌ "Be more careful"
   ✅ "Add automated test for error path"
```

## 35.2 Postmortem Process

```
Incident → Postmortem → Action Items → Review → Improve
  │           │             │             │         │
  ▼           ▼             ▼             ▼         ▼
 Fix        Document      Track        Verify    Cycle
            within 5d     weekly       monthly   (retro)
```

## 35.3 Postmortem Meeting

**Agenda (60 min):**
- 5 min: Summary (IC)
- 10 min: Timeline review
- 15 min: Root cause analysis
- 20 min: Action items brainstorm
- 10 min: Assignment & timeline

**Participants:**
- Incident team (IC, SMEs)
- Stakeholders
- (Optional) Leadership

## 35.4 Action Items Priority

| Priority | Timeline | Example |
|---|---|---|
| **P0** | 1 week | Fix root cause |
| **P1** | 2 weeks | Add monitoring |
| **P2** | 1 month | Improve process |
| **P3** | 1 quarter | Nice-to-have |

## 35.5 Tracking

```markdown
# Postmortem Action Items

| ID | Action | Owner | Due | Status | Priority |
|---|---|---|---|---|---|
| PM-001 | Fix connection leak | @bob | 2026-04-16 | ✅ Done | P0 |
| PM-002 | Add error tests | @carol | 2026-04-17 | 🔄 In progress | P0 |
| PM-003 | Monitor DB pool | @dave | 2026-04-18 | ⏳ Pending | P1 |
| PM-004 | Canary deploy | @alice | 2026-04-20 | ⏳ Pending | P1 |
```

**Weekly review:** ทีมตรวจสอบ status ของ action items

## 35.6 Knowledge Sharing

**Sharing learnings:**
- Weekly postmortem review (ทุกทีม)
- Monthly "postmortem of the month"
- Quarterly pattern analysis
- Blog post (public, ถ้าเหมาะสม)
- Internal wiki

**Pattern analysis:**
```
Q2 2026 Postmortems:
- 5 incidents total
- 3 จาก code changes (60%)
- 2 จาก dependency (40%)
- Common pattern: error path ไม่มี test
- Action: enforce error path test ใน PR checklist
```

---

# บทที่ 36: Case Study

## 36.1 Payment Module Deployment

### 36.1.1 Initial State

```
- Monolith with 500K LOC
- Manual deploy every 2 weeks
- 2-3 hotfixes/week
- MTTR: 4 hours
- Test coverage: 35%
- No monitoring
```

### 36.1.2 Transformation Timeline

| Quarter | Milestone | Result |
|---|---|---|
| Q1 2025 | Clean Architecture + DDD | Test coverage 75% |
| Q2 2025 | CI/CD pipeline | Deploy 10x/day |
| Q3 2025 | Kubernetes + K8s | Zero-downtime |
| Q4 2025 | Observability | MTTR 30 min |
| Q1 2026 | SLO + Alerts | MTTR 10 min |
| Q2 2026 | Chaos testing | Confidence |

### 36.1.3 Results

| Metric | Before | After |
|---|---|---|
| **Deploy Frequency** | 0.5/week | 10/day |
| **Lead Time** | 2 weeks | 1 day |
| **MTTR** | 4 hours | 10 min |
| **Change Failure Rate** | 30% | 5% |
| **Test Coverage** | 35% | 88% |
| **Uptime** | 99.5% | 99.99% |
| **Incidents/month** | 8 | 1 |

### 36.1.4 Key Success Factors

1. **Executive buy-in** — Management committed
2. **Incremental** — ไม่ rewrite ทั้งหมด
3. **Team training** — 3 workshops
4. **Metrics-driven** — DORA metrics tracked
5. **Blameless culture** — Postmortems ไม่ blame
6. **Automation** — Everything as code

## 36.2 Payment Module v1 → v2

### 36.2.1 Change

```
Change: Amount float64 → decimal.Decimal
Risk: High
Impact: 8M payments/month
Timeline: 6 weeks
```

### 36.2.2 Strategy

```
Week 1: RFC + approval
Week 2: Add column, dual write
Week 3: Backfill 8M rows (2 hours)
Week 4: Switch read path
Week 5: Switch write path
Week 6: Drop old column
```

### 36.2.3 Results

| Metric | Target | Actual |
|---|---|---|
| Downtime | 0 | 0 |
| Data loss | 0 | 0 |
| Rounding errors | 0 | 0 |
| Performance impact | < 5% | -2% |
| Rollback needed | – | No |

### 36.2.4 Monitoring During Migration

```
Dashboard:
- Backfill progress: 8M rows in 2 hours
- Mismatch count: 0
- Error rate: 0.02%
- P95 latency: 118ms (baseline)
- DB pool: 40% (normal)

Alerts:
- Mismatch > 100 → investigate
- Error rate > 1% → rollback
- P95 > 300ms → rollback
```

## 36.3 Lessons from Case Study

### ✅ ทำถูก

1. **Test First** — ครอบ behavior ก่อนแก้
2. **Small Steps** — 1 phase/week
3. **Feature Flag** — rollback 30s
4. **Staging ที่ data ≈ prod** — test ครบ
5. **Communication** — แจ้งล่วงหน้า 2 สัปดาห์
6. **Monitoring** — วัดทุก phase

### ❌ ข้อผิดพลาด

1. **Underestimated backfill time** — 8M rows ใช้เวลา 2 ชม. (คาด 30 นาที)
2. **Forgot to drop index** — หลัง backfill มี index เก่าค้าง
3. **Missing audit** — ไม่ได้ audit trail ของการเปลี่ยน
4. **Cost overrun** — ใช้เวลา 6 สัปดาห์ (คาด 4)

### 💡 Key Learnings

1. **Always test on production-size data**
2. **Plan for 2x time** (buffer)
3. **Include cleanup phase**
4. **Document everything** (สำหรับ rollback)

---

# ภาคผนวก A: Test Templates

## A.1 Domain Test Template

```go
package {ENTITY}_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "your/module/domain/entity"
)

func TestNew{ENTITY}_Success(t *testing.T) {
    // Arrange
    // ...
    
    // Act
    result, err := entity.New{ENTITY}(...)
    
    // Assert
    require.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, expected, result.Field)
}

func Test{ENTITY}_{Behavior}(t *testing.T) {
    tests := []struct{
        name    string
        setup   func() *entity.{ENTITY}
        action  func(*entity.{ENTITY}) error
        wantErr error
    }{
        // cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            e := tt.setup()
            err := tt.action(e)
            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## A.2 Use Case Test Template

```go
package {USECASE}_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
    "your/module/application"
    "your/module/application/mocks"
)

func Test{USECASE}_Success(t *testing.T) {
    repo := new(mocks.{REPO}Mock)
    log := logger.NewNoop()
    uc := application.New{USECASE}(repo, log)
    
    repo.On("Method", mock.Anything, mock.Anything).Return(value, nil)
    
    output, err := uc.Execute(context.Background(), application.Input{...})
    
    require.NoError(t, err)
    assert.NotNil(t, output)
    repo.AssertExpectations(t)
}
```

## A.3 Integration Test Template

```go
//go:build integration

package {REPO}_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/require"
    "your/module/testutils"
)

func Test{REPO}_Integration(t *testing.T) {
    testDB := testutils.SetupPostgres(t)
    repo := New{REPO}(testDB.DB)
    ctx := context.Background()
    
    // Test
    result, err := repo.Method(ctx, input)
    
    require.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

---

# ภาคผนวก B: CI/CD Templates

## B.1 Minimal CI

```yaml
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.22' }
      - run: go test ./...
      - run: go build ./...
```

## B.2 Standard CI

```yaml
name: CI
on: [push, pull_request]
jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.22', cache: true }
      
      - name: Lint
        uses: golangci/golangci-lint-action@v4
      
      - name: Test
        run: go test -race -coverprofile=coverage.out ./...
      
      - name: Coverage gate
        run: |
          cov=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          if (( $(echo "$cov < 80" | bc -l) )); then exit 1; fi
      
      - name: Build
        run: go build ./...
```

## B.3 Full Production CI/CD

ดูตัวอย่างเต็มในบทที่ 19

---

# ภาคผนวก C: Kubernetes Manifests

## C.1 Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {app}
  namespace: {namespace}
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: {app}
  template:
    metadata:
      labels:
        app: {app}
    spec:
      containers:
      - name: app
        image: {image}
        ports:
        - containerPort: 8080
        resources:
          requests: { cpu: 100m, memory: 128Mi }
          limits: { cpu: 500m, memory: 512Mi }
        livenessProbe:
          httpGet: { path: /health/live, port: 8080 }
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet: { path: /health/ready, port: 8080 }
          initialDelaySeconds: 5
          periodSeconds: 10
        securityContext:
          runAsNonRoot: true
          runAsUser: 1000
          readOnlyRootFilesystem: true
          allowPrivilegeEscalation: false
```

## C.2 Service + Ingress

```yaml
apiVersion: v1
kind: Service
metadata:
  name: {app}
spec:
  selector: { app: {app} }
  ports:
  - port: 80
    targetPort: 8080
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {app}
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  ingressClassName: nginx
  tls:
  - hosts: [api.example.com]
    secretName: {app}-tls
  rules:
  - host: api.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service: { name: {app}, port: { number: 80 } }
```

## C.3 HPA

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: {app}
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: {app}
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target: { type: Utilization, averageUtilization: 70 }
```

---

# ภาคผนวก D: Quick Reference Card

## D.1 Test Coverage Targets

```
┌─────────────────────┬──────────┐
│  Layer              │ Coverage │
├─────────────────────┼──────────┤
│  Domain             │ 100%     │
│  Application        │ 90%+     │
│  Infrastructure     │ 70%+     │
│  Interfaces         │ 60%+     │
│  Overall            │ 80%+     │
└─────────────────────┴──────────┘
```

## D.2 Test Commands

```bash
# Unit
go test ./...
go test -race -cover ./...
go test -run TestName ./...
go test -v ./...

# Integration
go test -tags=integration ./...

# Fuzz
go test -fuzz=FuzzX -fuzztime=30s ./...

# Benchmark
go test -bench=. -benchmem ./...

# Coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## D.3 Deployment Checklist

```
BEFORE
□ Tests pass (80%+ coverage)
□ Security scan clean
□ Migration ready
□ Rollback plan
□ Communication

DURING
□ Deploy canary (5%)
□ Monitor 5 min
□ Increase to 25%
□ Monitor 10 min
□ Increase to 100%

AFTER
□ Smoke test
□ Monitor 24h
□ Verify metrics
□ Update docs
□ Postmortem
```

## D.4 Health Endpoints

```
GET /health/live     ← Liveness (is process alive?)
GET /health/ready    ← Readiness (can serve?)
GET /health/startup  ← Startup (finished boot?)
GET /metrics         ← Prometheus metrics
```

## D.5 Alert Severity

| Sev | Response | Channel |
|---|---|---|
| P1 | 15 min | Phone, SMS |
| P2 | 1 hour | Slack, Email |
| P3 | 4 hours | Slack |
| P4 | Next day | Email |

## D.6 DORA Metrics

```
1. Deployment Frequency:  → Higher better
2. Lead Time for Changes: → Lower better
3. Change Failure Rate:   → Lower better
4. MTTR:                  → Lower better

Elite Performers:
1. Deploy: on-demand (multiple/day)
2. Lead: < 1 day
3. Failure: 0-15%
4. MTTR: < 1 hour
```

## D.7 Incident Response

```
Detect → Triage → Mitigate → Resolve → Postmortem
  ↓        ↓         ↓          ↓           ↓
 Alert    SEV      Fix/        Verify      Learn
         level    Rollback    Stable      Actions
```

## D.8 Golden Signals

```
1. Latency (P50, P95, P99)
2. Traffic (req/s)
3. Errors (rate, 5xx)
4. Saturation (CPU, mem, disk)
```

---

# 📝 ข้อมูลเอกสาร

**ชื่อเอกสาร:** คู่มือทดสอบและ Deployment ฉบับสมบูรณ์
**เวอร์ชัน:** 1.0
**วันที่:** เมษายน 2026
**จำนวนหน้า:** ~180 หน้า (ประมาณ)
**ระดับ:** Intermediate - Advanced

**ผู้อ่านเป้าหมาย:**
- Senior Go Developer
- DevOps / SRE Engineer
- Tech Lead ที่ดูแล Production
- Platform Engineer

**ข้อกำหนดเบื้องต้น:**
- อ่านเล่ม 1-3 จบ
- ประสบการณ์ production 1+ ปี
- เข้าใจ Docker, Kubernetes
- พื้นฐาน Linux

**เอกสารที่เกี่ยวข้อง:**
- เล่ม 1: คู่มือสร้าง Module ใหม่ ✅
- เล่ม 2: คู่มือแก้ไข Module เดิม ✅
- เล่ม 3: คู่มือขาย Module ✅
- เล่ม 5: คู่มือบำรุงรักษาและ Scale (ถัดไป)

**อ้างอิง:**
- Continuous Delivery — Humble & Farley
- Site Reliability Engineering — Google
- Accelerate — Forsgren, Humble, Kim
- The DevOps Handbook — Kim et al.
- Release It! — Michael Nygard
- DORA State of DevOps Reports
- โปรเจกต์ `icmongolang` (33 modules, production-tested)

---

**END OF BOOK 4**

> 📌 **ขั้นถัดไป:** อ่านเล่ม 5 — คู่มือบำรุงรักษาและ Scale
> ที่จะสอนวิธี maintain module ระยะยาว, scale horizontally/vertically,
> performance tuning, cost optimization, และ long-term sustainability

---

**พิมพ์เมื่อ:** เมษายน 2026
**ผู้จัดทำ:** ทีมสถาปัตยกรรมซอฟต์แวร์ icmongolang
**ติดต่อ:** kongnakornjantakun@gmail.com