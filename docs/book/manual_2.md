# 📗 เล่ม 2: คู่มือแก้ไข Module เดิมฉบับสมบูรณ์
## The Complete Guide to Safely Modifying Existing Go Modules

> **เวอร์ชัน 1.0 (เมษายน 2026)**
> เอกสารต้นแบบสำหรับการแก้ไข Module ใน Production
> ยึดหลัก Zero-Downtime, Backward Compatible, และ Blameless Post-Mortem

---

# สารบัญ

**ภาคที่ 1: ปรัชญาและหลักการ**
1. [บทนำ — ทำไมการแก้ Module เดิมยาก](#บทที่-1-บทนำ)
2. [กฎ 7 ข้อก่อนแตะโค้ด](#บทที่-2-กฎ-7-ข้อ)
3. [Impact Analysis Framework](#บทที่-3-impact-analysis)

**ภาคที่ 2: ประเภทของการแก้ไข**
4. [8 ประเภทของการแก้ไข](#บทที่-4-8-ประเภทของการแก้ไข)
5. [Level of Risk Matrix](#บทที่-5-level-of-risk)

**ภาคที่ 3: Scenarios ต่างๆ**
6. [Scenario 1: เพิ่ม Field ใหม่](#บทที่-6-เพิ่ม-field-ใหม่)
7. [Scenario 2: เปลี่ยน Type ของ Field](#บทที่-7-เปลี่ยน-type)
8. [Scenario 3: Rename Field/Entity](#บทที่-8-rename)
9. [Scenario 4: เปลี่ยน Business Rule](#บทที่-9-เปลี่ยน-business-rule)
10. [Scenario 5: เพิ่ม Relationship](#บทที่-10-เพิ่ม-relationship)
11. [Scenario 6: Breaking Change](#บทที่-11-breaking-change)
12. [Scenario 7: Refactor โครงสร้าง](#บทที่-12-refactor)
13. [Scenario 8: ย้ายจาก Sync → Async](#บทที่-13-sync-to-async)

**ภาคที่ 4: Database Migration**
14. [Expand-Contract Pattern](#บทที่-14-expand-contract)
15. [Zero-Downtime Migration](#บทที่-15-zero-downtime)
16. [Backfill ข้อมูลขนาดใหญ่](#บทที่-16-backfill)

**ภาคที่ 5: API Management**
17. [API Versioning](#บทที่-17-api-versioning)
18. [Deprecation Strategy](#บทที่-18-deprecation)

**ภาคที่ 6: Testing & Deployment**
19. [Regression Testing](#บทที่-19-regression-testing)
20. [Rollback Strategy](#บทที่-20-rollback)
21. [Feature Flags](#บทที่-21-feature-flags)

**ภาคที่ 7: กรณีศึกษาจริง**
22. [Case Study: Payment v1 → v2](#บทที่-22-case-study)
23. [Checklist และ Anti-Patterns](#บทที่-23-checklist)

**ภาคผนวก**
- [A. Migration Templates](#ภาคผนวก-a)
- [B. Communication Plan](#ภาคผนวก-b)
- [C. Quick Reference Card](#ภาคผนวก-c)

---

# บทที่ 1: บทนำ

## 1.1 ทำไมการแก้ Module เดิมยากกว่าการสร้างใหม่

**การสร้างใหม่:**
- ไม่มีข้อผูกมัด — ออกแบบได้อิสระ
- ไม่มีผู้ใช้ — ไม่ต้อง backward compatible
- ไม่มี data — migration ง่าย
- ไม่มี pressure — ไม่มี production ล่ม

**การแก้ของเดิม:**
- ❌ มีผู้ใช้อยู่ (external, internal, partner)
- ❌ มี data ใน production (ล้าน rows)
- ❌ มี SLA ต้องรักษา (99.9% uptime)
- ❌ มี dependency จาก module อื่น
- ❌ มี audit trail ต้องเก็บ
- ❌ มี compliance (PDPA, GMP, ISO)
- ❌ มี pressure จาก business (ยอดขายกำลังโต)

## 1.2 ต้นทุนของการแก้ผิด

**ตัวอย่างจริง:** e-commerce platform ล่ม 4 ชั่วโมงเพราะ migration ผิด

| ผลกระทบ | ต้นทุน |
|---|---|
| Revenue loss | ฿2,400,000 (4 ชม. × ฿600k/ชม.) |
| Engineer overtime | ฿80,000 (5 คน × 8 ชม.) |
| Customer support | ฿45,000 |
| Customer churn (2%) | ฿500,000/เดือน |
| Brand damage | ประเมินไม่ได้ |
| **รวม** | **~฿3,500,000** |

**สาเหตุ:** `ALTER TABLE` ที่ lock table 4 ชั่วโมง
**แก้:** ควรใช้ `CREATE INDEX CONCURRENTLY` (ใช้เวลา 2 นาที)

## 1.3 หลักคิดสำคัญ

```
"Code is read 10x more than written.
 Modules are MODIFIED 100x more than created."
```

การแก้ Module ไม่ใช่แค่ "แก้โค้ดให้ทำงาน" แต่ต้อง:
1. **ไม่ทำลายของเดิม** (backward compatible)
2. **ไม่ทำ production ล่ม** (zero-downtime)
3. **ไม่ทำ data หาย** (safety first)
4. **ไม่ทำทีมอื่นพัง** (coordination)
5. **ไม่ทำให้ rollback ยาก** (reversibility)

## 1.4 ตัวอย่างจริงจาก icmongolang

**Case 1: เปลี่ยน Payment Amount จาก `float64` → `decimal.Decimal`**

- ผู้ใช้: 8,000 payment/วัน
- ระยะเวลา: 3 sprints (6 สัปดาห์)
- Phases: Expand → Dual Write → Migrate → Contract
- Downtime: 0 นาที
- Data loss: 0 rows
- Bug: 0

**Case 2: เปลี่ยน Auth จาก JWT → JWT + Session**

- ผู้ใช้: 200,000 users
- ระยะเวลา: 2 sprints
- Strategy: Feature flag + parallel run 2 สัปดาห์
- Downtime: 0 นาที
- Rollback: ทำได้ภายใน 5 นาที (flip flag)

---

# บทที่ 2: กฎ 7 ข้อก่อนแตะโค้ด

## กฎที่ 1: อ่านให้เข้าใจ 100% ก่อนแก้ 1 บรรทัด

**Anti-pattern:**
```
Developer เห็นบั๊กในโค้ด → แก้ทันที → push → เกิดบั๊กใหม่ 3 ที่
```

**Pattern:**
```
1. อ่านโค้ดที่เกี่ยวข้องทั้งหมด (grep หา call sites)
2. เข้าใจ business rule ที่ซ่อนอยู่
3. เขียน diagram ว่า data flow ไปไหน
4. เขียน test ครอบ behavior เดิม
5. แก้ทีละเล็กๆ
```

**คำสั่งช่วย:**
```bash
# หาทุกที่ที่ใช้ function/struct นี้
grep -rn "FindByOrderID" --include="*.go" .

# ดู git history ของไฟล์นี้
git log --oneline -20 internal/modules/payment/domain/entity/payment.go

# ดูว่าใครแก้ล่าสุด และแก้ทำไม
git log -p -5 internal/modules/payment/domain/entity/payment.go

# ดูว่า commit ไหนเปลี่ยน behavior
git log -L :Refund:internal/modules/payment/domain/entity/payment.go
```

## กฎที่ 2: เขียน Test ครอบ Behavior เดิมก่อน Refactor

**หลักการ:** "Characterization Test" — Test ที่ lock behavior ปัจจุบันไว้

```go
// ก่อนแก้ Refund() → เขียน test ครอบ current behavior
func TestPayment_Refund_ExistingBehavior(t *testing.T) {
    // ปัจจุบัน: refund จาก PENDING ได้
    p, _ := entity.NewPayment(uuid.New(), uuid.New(), decimal.NewFromInt(100), "THB", "CARD")
    
    err := p.Refund()
    
    // ⚠️ Test นี้ lock behavior ปัจจุบัน (ซึ่งอาจจะผิด)
    // ถ้า behavior เปลี่ยน Test นี้จะ fail → เตือนให้ตั้งใจเปลี่ยน
    assert.NoError(t, err)  // ← current behavior
}
```

**ทำไมสำคัญ:**
- ถ้า refactor โดยไม่รู้ว่าอะไรเปลี่ยน → bug เงียบ
- Test เก่าที่ pass = safety net
- Test ที่ fail = change detector

## กฎที่ 3: แก้ทีละเล็กๆ ทีละ Commit

**Anti-pattern:**
```bash
git add .
git commit -m "refactor payment module"
git push
# 10,000 บรรทัดเปลี่ยน, ไม่รู้ว่า bug มาจากไหน
```

**Pattern:**
```bash
# Commit 1: เพิ่ม test ครอบ behavior เดิม
git commit -m "test: add characterization tests for Payment.Refund"

# Commit 2: เพิ่ม field ใหม่ (dual write)
git commit -m "feat(payment): add provider_txn_id field (dual write)"

# Commit 3: Migrate data
git commit -m "chore(payment): backfill provider_txn_id from transactions"

# Commit 4: Switch read
git commit -m "feat(payment): read from provider_txn_id"

# Commit 5: Drop old column
git commit -m "chore(payment): drop legacy txn_id column"
```

**ประโยชน์:**
- Rollback ทีละขั้นได้
- Code review ง่าย
- Bisect หา bug ง่าย

## กฎที่ 4: อย่าแก้ Domain เพื่อให้ Infrastructure ทำงาน

**Anti-pattern:**
```go
// Infrastructure บอกว่า "column นี้ nullable ไม่ได้"
// → Developer เพิ่ม pointer ใน Domain

type Payment struct {
    Amount *decimal.Decimal  // ← pointer เพราะ DB ไม่ให้ null
}
// ตอนนี้ Domain ต้อง check nil ทุกที่ → business logic เน่า
```

**Pattern:**
```go
// Domain ยังคง clean
type Payment struct {
    Amount decimal.Decimal
}

// Infrastructure จัดการ nullable → default zero
func (r *repo) toModel(e *entity.Payment) *PaymentModel {
    return &PaymentModel{
        Amount: e.Amount.String(),  // default "0" ถ้าจำเป็น
    }
}
```

**หลักการ:** Domain ไม่ควรรู้เรื่อง DB constraints

## กฎที่ 5: ถ้าต้อง Breaking Change → Version

**Anti-pattern:**
```go
// เปลี่ยน field name ตรงๆ → client เก่าพัง
type PaymentResponse struct {
    ID string `json:"payment_id"`  // ← เปลี่ยนจาก "id"
}
```

**Pattern:**
```go
// API v1 ยังคงเดิม
// API v2 เพิ่มใหม่
r.Route("/api/v2/payments", func(r chi.Router) { /* new */ })
r.Route("/api/v1/payments", func(r chi.Router) { /* legacy */ })
```

## กฎที่ 6: ก่อน Production → Staging ที่ data ≈ Production

**Anti-pattern:**
```
Test บน local: 100 rows → ผ่าน
Deploy production: 10M rows → ล่ม
```

**Pattern:**
```
Staging DB: clone จาก production (anonymized)
- 10M rows
- Same indexes
- Same constraints
- Same query patterns
Test migration → รู้เวลา + lock impact
```

## กฎที่ 7: ถามก่อนเสมอ — "Rollback plan คืออะไร?"

**ก่อน merge PR ต้องตอบได้:**
1. ถ้า production ล่มหลัง deploy จะ rollback ยังไง?
2. ใช้เวลาเท่าไหร่?
3. Data ที่ insert ใหม่จะเป็นยังไง?
4. Migration rollback ได้ไหม?
5. ถ้า rollback ไม่ได้ → ทำไม?

**ถ้าตอบไม่ได้ → ยังไม่พร้อม deploy**

---

# บทที่ 3: Impact Analysis Framework

## 3.1 คำถาม 10 ข้อที่ต้องตอบก่อนแก้

| # | คำถาม | หาได้จาก |
|---|---|---|
| 1 | ใครใช้ module นี้? | grep, code search, Slack search |
| 2 | มี external client ไหม? | API gateway config, Postman collection |
| 3 | มี module อื่น depend ไหม? | `go mod graph` |
| 4 | มี Kafka consumer ไหม? | Kafka topic registry |
| 5 | มี data ใน production เท่าไหร่? | DB query |
| 6 | Query pattern ที่ใช้บ่อย? | pg_stat_statements |
| 7 | SLA ที่ต้องรักษา? | SLO doc |
| 8 | มี compliance requirement? | Audit log |
| 9 | Test coverage ปัจจุบัน? | coverage report |
| 10 | Rollback plan? | design doc |

## 3.2 Dependency Discovery

### 3.2.1 หา Call Sites ใน Code

```bash
# หา struct usage
grep -rn "payment\.Payment" --include="*.go" .

# หา method usage
grep -rn "\.Refund(" --include="*.go" .

# หา import
grep -rn "modules/payment" --include="*.go" .

# ดู Go module graph
go mod graph | grep payment
```

### 3.2.2 หา External Client

```bash
# Nginx access log
grep "/api/v1/payments" /var/log/nginx/access.log | \
  awk '{print $1}' | sort -u | wc -l

# ที่ไหนส่ง request มา
awk '{print $1}' /var/log/nginx/access.log | \
  sort | uniq -c | sort -rn | head -20
```

### 3.2.3 หา Kafka Consumer

```bash
# List consumer groups
kafka-consumer-groups.sh --bootstrap-server localhost:9092 --list

# ดู group ที่ consume topic นี้
kafka-consumer-groups.sh --bootstrap-server localhost:9092 \
  --describe --group payment-processor | grep payment.created
```

### 3.2.4 หา Data ใน Production

```sql
-- จำนวน rows
SELECT COUNT(*) FROM payment_transactions;

-- ขนาดตาราง
SELECT
    pg_size_pretty(pg_total_relation_size('payment_transactions')) AS total_size,
    pg_size_pretty(pg_relation_size('payment_transactions')) AS table_size,
    pg_size_pretty(pg_indexes_size('payment_transactions')) AS index_size;

-- Query ที่ใช้บ่อย
SELECT
    query,
    calls,
    mean_exec_time
FROM pg_stat_statements
WHERE query LIKE '%payment_transactions%'
ORDER BY calls DESC
LIMIT 10;

-- Index usage
SELECT
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
WHERE relname = 'payment_transactions'
ORDER BY idx_scan DESC;
```

## 3.3 Impact Matrix Template

```markdown
# Impact Analysis: [Change Title]

## Scope
- Module: payment
- Entity: Payment
- Type: Add field
- Risk Level: Low

## Impact

| Area | Impact | Mitigation |
|---|---|---|
| **Code** | 12 ไฟล์ต้องแก้ | ทีละ layer |
| **DB** | เพิ่ม 1 column | Add nullable |
| **API** | Response เพิ่ม field | Backward compatible |
| **Kafka** | ไม่กระทบ | – |
| **Clients** | ไม่กระทบ (additive) | – |
| **Performance** | +5% storage | OK |
| **Compliance** | ไม่กระทบ | – |

## Timeline
- Day 1-2: เพิ่ม field + migration
- Day 3-4: Test staging
- Day 5: Deploy production
- Day 6-7: Monitor

## Rollback
- Code: `git revert` (5 นาที)
- DB: ไม่ต้อง rollback (additive)
- Data: ไม่หาย

## Communication
- แจ้งทีม: 2 วันล่วงหน้า
- แจ้ง partner: ไม่ต้อง (additive)
- Status page: ไม่ต้อง
```

## 3.4 Risk Scoring

| Factor | Score 1 (Low) | Score 5 (High) |
|---|---|---|
| **จำนวน Users** | < 100 | > 100,000 |
| **Data Size** | < 10K rows | > 100M rows |
| **API Change** | Additive | Breaking |
| **DB Change** | Add column | Drop/rename |
| **Dependencies** | 0 modules | > 5 modules |
| **SLA Impact** | None | 99.99% |
| **Rollback** | Easy | Impossible |

**Total Score:**
- 7-15: **Low risk** → deploy ปกติ
- 16-25: **Medium risk** → feature flag, canary
- 26-35: **High risk** → ค่อยๆ ทำ, communication plan

---

# บทที่ 4: 8 ประเภทของการแก้ไข

## 4.1 ภาพรวม

```
                    ┌─────────────────────┐
                    │  Types of Changes   │
                    └──────────┬──────────┘
                               │
        ┌──────────────────────┼──────────────────────┐
        │                      │                      │
        ▼                      ▼                      ▼
   ┌─────────┐           ┌─────────┐           ┌─────────┐
   │Additive │           │Modifying│           │Removing │
   ├─────────┤           ├─────────┤           ├─────────┤
   │+ Field  │           │~ Type   │           │- Field  │
   │+ Method │           │~ Rule   │           │- Entity │
   │+ API    │           │~ Rename │           │- Module │
   └─────────┘           └─────────┘           └─────────┘
    LOW RISK            MEDIUM RISK            HIGH RISK
```

## 4.2 ตารางเปรียบเทียบ

| # | ประเภท | ตัวอย่าง | Risk | Rollback |
|---|---|---|---|---|
| 1 | **Add Field** | เพิ่ม `notes` | Low | ง่าย |
| 2 | **Add Method** | เพิ่ม `Cancel()` | Low | ง่าย |
| 3 | **Add API** | `/v2/payments` | Low | ง่าย |
| 4 | **Modify Rule** | เปลี่ยน min amount | Medium | Feature flag |
| 5 | **Change Type** | `float64 → decimal` | Medium | 3-phase |
| 6 | **Rename** | `txn_id → provider_txn_id` | High | 4-phase |
| 7 | **Split Entity** | `Payment → Payment + Refund` | High | Multi-sprint |
| 8 | **Remove** | ลบ `legacy_field` | High | Deprecation first |

## 4.3 กฎทอง: Additive First, Subtractive Last

```
    ADD ───────► MODIFY ───────► REMOVE
   (Phase 1)     (Phase 2)      (Phase 3)
   
   ปลอดภัย        ระวัง          ต้อง planning
```

**ทำไม?** เพราะ add ไม่ทำของเดิมพัง, modify ต้องดูผล, remove ต้อง deprecate ก่อน

---

# บทที่ 5: Level of Risk

## 5.1 Risk Levels

### Level 1: Trivial (Risk Score 1-5)

**ตัวอย่าง:** เพิ่ม log, แก้ typo, เพิ่ม comment

**Process:**
- ✅ PR + 1 reviewer
- ✅ CI ผ่าน
- ✅ Merge → deploy

**เวลา:** 30 นาที

### Level 2: Low (Risk Score 6-12)

**ตัวอย่าง:** เพิ่ม field optional, เพิ่ม method ที่ไม่แตะของเดิม

**Process:**
- ✅ PR + 2 reviewers
- ✅ Unit test ใหม่
- ✅ CI ผ่าน
- ✅ Merge → deploy
- ✅ Monitor 1 ชั่วโมง

**เวลา:** 1-2 วัน

### Level 3: Medium (Risk Score 13-20)

**ตัวอย่าง:** เปลี่ยน business rule, เปลี่ยน type

**Process:**
- ✅ PR + 2 reviewers + tech lead
- ✅ Unit + integration tests
- ✅ Test staging ที่ data ≈ production
- ✅ Feature flag (ถ้าได้)
- ✅ Canary deploy (5% → 25% → 100%)
- ✅ Monitor 24 ชั่วโมง
- ✅ Rollback plan ready

**เวลา:** 1 สัปดาห์

### Level 4: High (Risk Score 21-30)

**ตัวอย่าง:** Breaking change, rename field, split entity

**Process:**
- ✅ RFC document + approval
- ✅ Full test suite
- ✅ Multi-phase (Expand-Contract)
- ✅ Parallel run (old + new)
- ✅ Communication plan
- ✅ Rollback documented
- ✅ Postmortem plan

**เวลา:** 3-6 สัปดาห์

### Level 5: Critical (Risk Score 31+)

**ตัวอย่าง:** แก้ Auth mechanism, เปลี่ยน Database engine

**Process:**
- ✅ Board-level approval
- ✅ Full change management
- ✅ Dry run 3 ครั้ง
- ✅ Rollback rehearsal
- ✅ War room
- ✅ Rollback within 5 minutes

**เวลา:** 2-3 เดือน

## 5.2 Risk-Based Process Matrix

| Risk Level | Approval | Testing | Deploy | Monitor |
|---|---|---|---|---|
| **Trivial** | 1 reviewer | Unit | Auto | – |
| **Low** | 2 reviewers | Unit + Integration | Auto | 1h |
| **Medium** | + Tech lead | + Staging | Canary | 24h |
| **High** | + Architect | + Parallel run | Phased | 1 week |
| **Critical** | + Board | + Dry run | Manual | 1 month |

---

# บทที่ 6: เพิ่ม Field ใหม่

## 6.1 ระดับความเสี่ยง: Low

การเพิ่ม field เป็น operation ที่ปลอดภัยที่สุด เพราะ:
- ไม่มีอะไรพัง (additive)
- Client เก่ายังทำงานได้ (ignore field ใหม่)
- Rollback ง่าย (แค่ไม่ใช้ field)

## 6.2 Scenario: เพิ่ม `notes` field ใน Payment

### Phase 1: วิเคราะห์ผลกระทบ

```markdown
## Impact Analysis
- **Type**: Add optional field
- **Risk**: Low
- **Users**: ไม่กระทบ (additive)
- **Data**: +50 bytes/row (ไม่สำคัญ)
- **API**: Response เพิ่ม field
- **Rollback**: ไม่ต้อง (ไม่ใช้ field ก็ได้)
```

### Phase 2: แก้ Domain Entity

```go
// internal/modules/payment/domain/entity/payment.go

type Payment struct {
    ID uuid.UUID `json:"id"`
    // ... existing fields
    
    // ⬇️ เพิ่มใหม่ — field ที่ไม่บังคับ
    // ⬇️ New — optional field
    Notes string `json:"notes,omitempty"`
}
```

**⚠️ หมายเหตุ:**
- ไม่แก้ constructor `NewPayment` (ใช้ zero value `""`)
- ไม่แตะ behavior methods
- ไม่กระทบ validation

### Phase 3: แก้ GORM Model

```go
// internal/modules/payment/infrastructure/persistence/postgres/models.go

type PaymentModel struct {
    // ... existing
    Notes string `gorm:"type:text;default:''"`
}
```

### Phase 4: Migration (Safe)

```sql
-- migrations/20260405_payment_add_notes.sql
-- ============================================================
-- Add notes column to payment_transactions
-- Safe: nullable with default (no lock)
-- ============================================================

ALTER TABLE payment_transactions
    ADD COLUMN IF NOT EXISTS notes TEXT DEFAULT '';

-- Optional: comment for documentation
COMMENT ON COLUMN payment_transactions.notes IS 
    'Optional notes for payment (added 2026-04-05)';
```

**⚠️ คำเตือน:**
- `ADD COLUMN` ใน Postgres 11+ เป็น **metadata-only** (ไม่ rewrite table)
- ใช้ `IF NOT EXISTS` เพื่อ idempotent
- Default value ไม่ lock table

### Phase 5: Repository Mapping

```go
// internal/modules/payment/infrastructure/persistence/postgres/payment_repo_impl.go

func (r *paymentRepoImpl) toModel(e *entity.Payment) *PaymentModel {
    return &PaymentModel{
        // ... existing
        Notes: e.Notes,  // ⬅️ เพิ่ม
    }
}

func (r *paymentRepoImpl) toEntity(m *PaymentModel) *entity.Payment {
    return &entity.Payment{
        // ... existing
        Notes: m.Notes,  // ⬅️ เพิ่ม
    }
}
```

### Phase 6: อัปเดต Use Case (ถ้าต้องรับ input)

```go
// internal/modules/payment/application/create_payment.go

type CreatePaymentInput struct {
    UserID   uuid.UUID
    OrderID  uuid.UUID
    Amount   decimal.Decimal
    Currency string
    Method   string
    Notes    string  // ⬅️ เพิ่มใหม่
}

func (uc *CreatePaymentUseCase) Execute(ctx context.Context, input CreatePaymentInput) (*CreatePaymentOutput, error) {
    // ... existing
    payment, err := entity.NewPayment(input.UserID, input.OrderID, input.Amount, input.Currency, method)
    if err != nil {
        return nil, err
    }
    
    payment.Notes = input.Notes  // ⬅️ set field ใหม่
    
    // ... rest
}
```

### Phase 7: อัปเดต HTTP DTO

```go
// internal/modules/payment/interfaces/http/dto.go

type CreatePaymentRequest struct {
    OrderID  string `json:"order_id" validate:"required,uuid"`
    Amount   string `json:"amount" validate:"required"`
    Currency string `json:"currency" validate:"required,len=3"`
    Method   string `json:"method" validate:"required,oneof=CARD QR_PROMPT_PAY BANK_TRANSFER"`
    Notes    string `json:"notes,omitempty" validate:"omitempty,max=500"`  // ⬅️ เพิ่ม
}

type PaymentResponse struct {
    ID       string `json:"id"`
    OrderID  string `json:"order_id"`
    Amount   string `json:"amount"`
    Currency string `json:"currency"`
    Method   string `json:"method"`
    Status   string `json:"status"`
    Notes    string `json:"notes,omitempty"`  // ⬅️ เพิ่ม
}
```

**⚠️ หมายเหตุ:** ใช้ `omitempty` เพื่อไม่ให้ field ว่างเปล่าปรากฏใน response

### Phase 8: Handler

```go
// internal/modules/payment/interfaces/http/payment_handler.go

output, err := h.createUC.Execute(r.Context(), application.CreatePaymentInput{
    UserID:   userID,
    OrderID:  orderID,
    Amount:   amount,
    Currency: req.Currency,
    Method:   req.Method,
    Notes:    req.Notes,  // ⬅️ ส่งต่อ
})
```

### Phase 9: Test

```go
// internal/modules/payment/application/create_payment_test.go

func TestCreatePaymentUseCase_WithNotes(t *testing.T) {
    // ... setup
    output, err := uc.Execute(context.Background(), application.CreatePaymentInput{
        UserID:   userID,
        OrderID:  orderID,
        Amount:   decimal.NewFromInt(100),
        Currency: "THB",
        Method:   "CARD",
        Notes:    "urgent order",  // ⬅️ test field ใหม่
    })
    require.NoError(t, err)
    assert.NotNil(t, output)
    // ... assert notes saved
}
```

### Phase 10: Deploy

```bash
# 1. Run migration on staging
go run cmd/api/main.go migrate

# 2. Test on staging
curl -X POST https://staging.api.com/api/v1/payments \
  -d '{"order_id":"...", "amount":"100", ..., "notes":"test"}'

# 3. Deploy to production
git tag v1.5.0
git push origin v1.5.0
# → CI/CD auto-deploy

# 4. Monitor
# Grafana: watch error rate, latency
```

### 6.3 Checklist

- [ ] Migration ใช้ `ADD COLUMN IF NOT EXISTS`
- [ ] Column มี default value
- [ ] Column nullable (หรือ default เพื่อไม่ให้ existing rows พัง)
- [ ] GORM model มี tag ครบ
- [ ] Repository map ครบทั้ง 2 ทาง (toModel, toEntity)
- [ ] DTO ใช้ `omitempty` (ถ้า optional)
- [ ] Validation มี (ถ้าจำเป็น)
- [ ] Test ครอบ field ใหม่
- [ ] Documentation อัปเดต
- [ ] Rollback plan: ไม่ต้อง (แค่ไม่ใช้ field)

## 6.4 กรณีศึกษา: เพิ่ม `metadata JSONB` ใน Payment

**Requirement:** เก็บข้อมูลเพิ่มเติมจาก partner (variable structure)

```go
// Domain
type Payment struct {
    // ...
    Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Model
type PaymentModel struct {
    // ...
    Metadata datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
}

// Migration
ALTER TABLE payment_transactions
    ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}';

-- Index for JSONB queries (ถ้าต้อง search)
CREATE INDEX IF NOT EXISTS idx_payment_transactions_metadata
    ON payment_transactions USING GIN (metadata jsonb_path_ops);
```

**ข้อควรระวัง:**
- JSONB กินพื้นที่มากกว่า column ปกติ
- Query ช้ากว่า column ปกติ (ต้อง index)
- Schema ไม่ชัดเจน (ควร document)

---

# บทที่ 7: เปลี่ยน Type ของ Field

## 7.1 ระดับความเสี่ยง: Medium

**ทำไมเสี่ยง?**
- Existing data ต้อง migrate
- Code ทุกที่ต้องแก้
- Client อาจ expect type เดิม

## 7.2 Scenario: `float64` → `decimal.Decimal` สำหรับ Amount

### Phase 1: วิเคราะห์

```markdown
## ปัญหาปัจจุบัน
- `float64` มี precision issue: 0.1 + 0.2 = 0.30000000000000004
- บัญชีคลาดเคลื่อนสะสม ฿500/เดือน
- Audit ไม่ผ่าน (rounding inconsistent)

## แนวทาง
- ใช้ `decimal.Decimal` (shopspring/decimal)
- Migration 3 phase
```

### Phase 2: หลักการ 3 Phase (Expand-Contract)

```
Phase 1: EXPAND       Phase 2: MIGRATE      Phase 3: CONTRACT
├─ Add new column    ├─ Backfill data       ├─ Switch read
├─ Dual write        ├─ Verify              ├─ Switch write
└─ Both columns      └─ Old still valid     └─ Drop old column
```

### Phase 3: Migration (Phase 1 — Expand)

```sql
-- migrations/20260410_payment_add_amount_decimal.sql
-- ============================================================
-- Phase 1: Add new column (keep old for rollback safety)
-- ============================================================

-- Add new column (nullable first, no lock)
ALTER TABLE payment_transactions
    ADD COLUMN IF NOT EXISTS amount_decimal NUMERIC(15, 2);

-- Comment for documentation
COMMENT ON COLUMN payment_transactions.amount_decimal IS
    'Precise amount, migration from amount (float64) on 2026-04-10';
```

### Phase 4: Domain Entity (Before → After)

```go
// BEFORE (v1)
type Payment struct {
    ID     uuid.UUID `json:"id"`
    Amount float64   `json:"amount"`  // ← เปลี่ยน
}

// AFTER (v2)
type Payment struct {
    ID     uuid.UUID       `json:"id"`
    Amount decimal.Decimal `json:"amount"`  // ← เปลี่ยน
}
```

**⚠️ Breaking Change!** — code ทุกที่ที่ใช้ `.Amount` ต้องแก้

### Phase 5: GORM Model (Dual Write)

```go
// internal/modules/payment/infrastructure/persistence/postgres/models.go

type PaymentModel struct {
    ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
    // ...
    
    // DEPRECATED: to be dropped in v2.1
    AmountLegacy float64 `gorm:"column:amount;type:double precision" json:"-"`
    
    // NEW: accurate amount
    AmountDecimal string `gorm:"column:amount_decimal;type:numeric(15,2)" json:"-"`
}
```

**Key Point:**
- ตอนนี้มี 2 columns
- `amount` (float64) — legacy
- `amount_decimal` (numeric) — new
- **Dual write** → เขียนทั้ง 2 ตอน Save

### Phase 6: Repository — Dual Write & Fallback Read

```go
// internal/modules/payment/infrastructure/persistence/postgres/payment_repo_impl.go

// Save — เขียนทั้ง 2 columns
func (r *paymentRepoImpl) Save(ctx context.Context, p *entity.Payment) error {
    m := &PaymentModel{
        ID:            p.ID,
        UserID:        p.UserID,
        // ...
        // Dual write
        AmountDecimal: p.Amount.String(),          // ← ใหม่
        AmountLegacy:  p.Amount.InexactFloat64(),  // ← เก่า (deprecated)
    }
    return r.db.WithContext(ctx).Create(m).Error
}

// FindByID — อ่านจาก column ใหม่, fallback ไปเก่า
func (r *paymentRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
    var m PaymentModel
    err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
    // ...
    return r.toEntity(&m), nil
}

func (r *paymentRepoImpl) toEntity(m *PaymentModel) *entity.Payment {
    // อ่านจาก decimal ก่อน
    amount, err := decimal.NewFromString(m.AmountDecimal)
    if err != nil || m.AmountDecimal == "" {
        // Fallback: อ่านจาก float64 เก่า (สำหรับ rows ที่ยังไม่ backfill)
        amount = decimal.NewFromFloat(m.AmountLegacy)
    }
    
    return &entity.Payment{
        ID:     m.ID,
        Amount: amount,  // ← คืน decimal เสมอ
        // ...
    }
}
```

### Phase 7: Backfill Data

```sql
-- migrations/20260410_payment_backfill_amount_decimal.sql
-- ============================================================
-- Phase 2: Backfill amount_decimal from amount
-- Use batch เพื่อไม่ lock table นาน
-- ============================================================

DO $$
DECLARE
    batch_size INT := 10000;
    total_updated INT := 0;
    batch_updated INT;
BEGIN
    LOOP
        WITH batch AS (
            SELECT id, amount
            FROM payment_transactions
            WHERE amount_decimal IS NULL
              AND amount IS NOT NULL
            LIMIT batch_size
            FOR UPDATE SKIP LOCKED
        )
        UPDATE payment_transactions pt
        SET amount_decimal = batch.amount::numeric(15, 2)
        FROM batch
        WHERE pt.id = batch.id;
        
        GET DIAGNOSTICS batch_updated = ROW_COUNT;
        total_updated := total_updated + batch_updated;
        
        RAISE NOTICE 'Backfilled % rows (total: %)', batch_updated, total_updated;
        
        EXIT WHEN batch_updated < batch_size;
        
        -- หยุดพัก DB
        PERFORM pg_sleep(0.1);
    END LOOP;
    
    RAISE NOTICE 'Backfill complete. Total: %', total_updated;
END $$;

-- Verify
SELECT COUNT(*) FROM payment_transactions WHERE amount_decimal IS NULL;
-- ควรได้ 0
```

**Key Techniques:**
- `FOR UPDATE SKIP LOCKED` — ไม่ block transaction อื่น
- Batch size 10K — ไม่ lock table นาน
- `pg_sleep(0.1)` — ให้ DB พัก

### Phase 8: Set NOT NULL

```sql
-- migrations/20260412_payment_amount_decimal_not_null.sql
-- ============================================================
-- Phase 2.5: Add NOT NULL after backfill complete
-- ============================================================

-- ใช้ NOT VALID เพื่อไม่ rewrite table
ALTER TABLE payment_transactions
    ADD CONSTRAINT payment_amount_decimal_not_null 
    CHECK (amount_decimal IS NOT NULL) NOT VALID;

-- Validate separately (ใช้เวลา แต่ไม่ lock)
ALTER TABLE payment_transactions
    VALIDATE CONSTRAINT payment_amount_decimal_not_null;
```

### Phase 9: Switch Read Path

```go
// toEntity — ลบ fallback หลัง backfill 100%
func (r *paymentRepoImpl) toEntity(m *PaymentModel) *entity.Payment {
    amount, _ := decimal.NewFromString(m.AmountDecimal)
    // Fallback ถูกลบแล้ว
    return &entity.Payment{
        ID:     m.ID,
        Amount: amount,
        // ...
    }
}
```

### Phase 10: Switch Write Path (หลัง observe 1 สัปดาห์)

```go
// ลบ dual write
func (r *paymentRepoImpl) Save(ctx context.Context, p *entity.Payment) error {
    m := &PaymentModel{
        ID:            p.ID,
        AmountDecimal: p.Amount.String(),  // ← เขียนแค่ column ใหม่
        // AmountLegacy ถูกลบ
    }
    // ...
}
```

### Phase 11: Drop Old Column (หลัง 2 sprints)

```sql
-- migrations/20260501_payment_drop_amount_legacy.sql
-- ============================================================
-- Phase 3: Drop legacy column
-- ============================================================

-- Rename ก่อน (เพื่อเก็บไว้ชั่วคราว)
ALTER TABLE payment_transactions
    RENAME COLUMN amount TO amount_legacy_20260501;

-- รอ 1 สัปดาห์ → drop จริง
-- ALTER TABLE payment_transactions DROP COLUMN amount_legacy_20260501;
```

**⚠️ ทำไม rename ก่อน drop?**
- ถ้ามี code เก่าที่ยังอ่านอยู่ → error จะชัดเจน
- ถ้าไม่มี error 1 สัปดาห์ → ปลอดภัย drop

### 7.3 Timeline รวม

| Phase | ระยะเวลา | สถานะ |
|---|---|---|
| Expand (Add column) | Day 1 | ✅ |
| Code (Dual write) | Day 2-3 | ✅ |
| Backfill | Day 4 (2 hours) | ✅ |
| Observe | Day 5-11 | ✅ |
| Switch read | Day 12 | ✅ |
| Observe | Day 12-25 | ✅ |
| Switch write | Day 26 | ✅ |
| Observe | Day 26-39 | ✅ |
| Drop old | Day 40 | ✅ |

**Total: 6 สัปดาห์**

### 7.4 Checklist

- [ ] เพิ่ม column ใหม่ (nullable ก่อน)
- [ ] Dual write ทั้ง 2 columns
- [ ] Backfill ด้วย batch + SKIP LOCKED
- [ ] Verify 100% backfill
- [ ] Switch read path
- [ ] Observe 2 สัปดาห์
- [ ] Switch write path
- [ ] Observe 2 สัปดาห์
- [ ] Rename old column (safety)
- [ ] Observe 1 สัปดาห์
- [ ] Drop old column
- [ ] ทุก phase มี test + rollback plan

---

# บทที่ 8: Rename Field/Entity

## 8.1 ระดับความเสี่ยง: High

**ทำไมเสี่ยง?**
- Breaking change สำหรับ client
- Data ต้อง migrate
- Code ทุกที่ต้องแก้

## 8.2 Scenario: `user_id` → `customer_id`

### Phase 1: Multi-Phase Plan

```
Phase 1: EXPAND          Phase 2: DUAL READ      Phase 3: SWITCH
├─ Add customer_id       ├─ Read both            ├─ Read customer_id only
├─ Dual write            ├─ Compare (log diff)   ├─ Stop writing user_id
└─ Both work             └─ Alert on mismatch    └─ Mark user_id deprecated

Phase 4: DEPRECATE       Phase 5: REMOVE
├─ API version bump      ├─ Drop user_id column
├─ Warn clients          └─ Clean up code
└─ 6 months notice
```

### Phase 2: Migration — Add New Column

```sql
-- migrations/20260415_payment_add_customer_id.sql
ALTER TABLE payment_transactions
    ADD COLUMN IF NOT EXISTS customer_id UUID;

CREATE INDEX IF NOT EXISTS idx_payment_transactions_customer_id
    ON payment_transactions(customer_id);
```

### Phase 3: Domain Entity (Keep Both)

```go
type Payment struct {
    ID uuid.UUID `json:"id"`
    
    // DEPRECATED: use CustomerID instead (will be removed in v3.0)
    // เลิกใช้: ใช้ CustomerID แทน (จะลบใน v3.0)
    UserID uuid.UUID `json:"user_id,omitempty"`
    
    // NEW: preferred field
    // ใหม่: field ที่แนะนำ
    CustomerID uuid.UUID `json:"customer_id"`
    
    // ...
}
```

**⚠️ หมายเหตุ:** ตอนนี้มี 2 fields ใน entity ด้วย (ช่วย transition)

### Phase 4: GORM Model

```go
type PaymentModel struct {
    ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
    
    // DEPRECATED
    UserID     uuid.UUID `gorm:"type:uuid;index"`
    
    // NEW
    CustomerID uuid.UUID `gorm:"type:uuid;index"`
    
    // ...
}
```

### Phase 5: Repository — Dual Write + Fallback Read

```go
func (r *paymentRepoImpl) Save(ctx context.Context, p *entity.Payment) error {
    // Ensure both fields are set
    customerID := p.CustomerID
    if customerID == uuid.Nil {
        customerID = p.UserID  // Fallback
    }
    
    m := &PaymentModel{
        // ...
        UserID:     customerID,  // Dual write
        CustomerID: customerID,  // Dual write
    }
    return r.db.WithContext(ctx).Create(m).Error
}

func (r *paymentRepoImpl) toEntity(m *PaymentModel) *entity.Payment {
    customerID := m.CustomerID
    if customerID == uuid.Nil {
        customerID = m.UserID  // Fallback
    }
    
    return &entity.Payment{
        UserID:     customerID,  // populate both
        CustomerID: customerID,
        // ...
    }
}
```

### Phase 6: Backfill

```sql
-- migrations/20260415_payment_backfill_customer_id.sql
UPDATE payment_transactions
SET customer_id = user_id
WHERE customer_id IS NULL;
```

**⚠️ ถ้า table ใหญ่ → batch processing** (ตามบทที่ 16)

### Phase 7: Verification (Dual Read Comparison)

```go
// ชั่วคราว — log diff ระหว่าง 2 columns
func (r *paymentRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
    var m PaymentModel
    // ...
    
    // Verify: 2 columns ต้องตรงกัน
    if m.UserID != m.CustomerID {
        log.Warn("customer_id mismatch",
            "payment_id", m.ID,
            "user_id", m.UserID,
            "customer_id", m.CustomerID,
        )
    }
    
    return r.toEntity(&m), nil
}
```

**Monitor 1 สัปดาห์:**
```promql
# Grafana query
sum(rate(payment_customer_id_mismatch_total[5m]))
# ต้องเป็น 0
```

### Phase 8: Switch Read Path

```go
// หลัง verify ว่า 100% ตรงกัน
func (r *paymentRepoImpl) toEntity(m *PaymentModel) *entity.Payment {
    return &entity.Payment{
        UserID:     m.CustomerID,  // ← อ่านจาก column ใหม่
        CustomerID: m.CustomerID,
        // ...
    }
}
```

### Phase 9: Update Code ทุกที่

```go
// BEFORE
payment.UserID

// AFTER
payment.CustomerID
```

**หา call sites:**
```bash
grep -rn "\.UserID" --include="*.go" internal/modules/payment/
```

### Phase 10: API Deprecation

```go
// API v1 — ยังคืน user_id (deprecated)
type PaymentResponseV1 struct {
    ID         string `json:"id"`
    UserID     string `json:"user_id"`      // ← DEPRECATED
    CustomerID string `json:"customer_id"`  // ← ใหม่
    // ...
}

// API v2 — คืนแค่ customer_id
type PaymentResponseV2 struct {
    ID         string `json:"id"`
    CustomerID string `json:"customer_id"`  // ← ไม่มี user_id
    // ...
}
```

**Deprecation Header:**
```go
w.Header().Set("X-API-Deprecated", "true")
w.Header().Set("X-API-Sunset-Date", "2026-10-15")
w.Header().Set("X-API-Replacement", "/api/v2/payments")
```

### Phase 11: Drop Column (หลัง 6 เดือน)

```sql
-- migrations/20261015_payment_drop_user_id.sql
ALTER TABLE payment_transactions DROP COLUMN IF EXISTS user_id;
```

### 8.3 Checklist

- [ ] Multi-phase plan documented
- [ ] RFC approved
- [ ] Migration ผ่าน staging
- [ ] Dual write + dual read
- [ ] Verify 100% match (≥ 1 week)
- [ ] Backfill complete
- [ ] Update all call sites
- [ ] API deprecation headers
- [ ] 6-month notice to clients
- [ ] Communication plan
- [ ] Postmortem scheduled

---

# บทที่ 9: เปลี่ยน Business Rule

## 9.1 ระดับความเสี่ยง: Medium-High

**ปัญหา:** Business rule เปลี่ยน = พฤติกรรมระบบเปลี่ยน = ผู้ใช้อาจไม่พอใจ

## 9.2 Scenario: เปลี่ยน Minimum Amount จาก ฿1 → ฿20

### 9.2.1 Feature Flag Approach (แนะนำ)

```go
// pkg/featureflag/flags.go

type Flags interface {
    IsEnabled(ctx context.Context, name string) bool
}

// Usage
type CreatePaymentUseCase struct {
    repo       repository.PaymentRepository
    logger     logger.Logger
    flags      Flags  // ⬅️ เพิ่มใหม่
}

func (uc *CreatePaymentUseCase) Execute(ctx context.Context, input CreatePaymentInput) (*CreatePaymentOutput, error) {
    // Old rule (ยังมีผล)
    if input.Amount.LessThan(decimal.NewFromInt(1)) {
        return nil, domainerrors.ErrInvalidAmount
    }
    
    // New rule (under flag)
    if uc.flags.IsEnabled(ctx, "payment.min_amount_20") {
        if input.Amount.LessThan(decimal.NewFromInt(20)) {
            uc.logger.Warn("amount below new minimum",
                "amount", input.Amount.String(),
                "user_id", input.UserID,
            )
            return nil, domainerrors.ErrAmountTooLow
        }
    }
    
    // ... rest
}
```

### 9.2.2 Rollout Plan

```
Week 1:   Flag OFF (0%)       — ไม่มีผล, observe traffic
Week 2:   Flag ON (5%)        — canary, เฉพาะ user กลุ่มเล็ก
Week 3:   Flag ON (25%)       — ขยาย
Week 4:   Flag ON (50%)       — ขยาย
Week 5:   Flag ON (100%)      — ทั้งหมด
Week 6:   Remove flag         — ลบ code เก่า
```

### 9.2.3 Metrics to Monitor

```promql
# Rate ของการ rejact เพราะต่ำกว่า 20 บาท
sum(rate(payment_rejected_total{reason="amount_too_low"}[5m]))

# เทียบกับก่อน flag
# ถ้าสูงขึ้นมาก → อาจต้อง rollback flag

# Revenue impact
sum(rate(payment_amount_sum[5m]))
```

### 9.2.4 Rollback

```bash
# ปิด flag ทันที (ไม่ต้อง deploy)
curl -X POST https://flag-service/api/flags/payment.min_amount_20 \
  -d '{"enabled": false}'

# ภายใน 30 วินาที ทั่วทั้งระบบ
```

### 9.3 Checklist

- [ ] Feature flag setup
- [ ] Old rule ยังทำงาน (fallback)
- [ ] New rule อยู่หลัง flag
- [ ] Test ทั้ง 2 path
- [ ] Metrics สำหรับ monitor
- [ ] Rollback: flip flag
- [ ] Communication plan
- [ ] Documentation อัปเดต

---

# บทที่ 10: เพิ่ม Relationship

## 10.1 ระดับความเสี่ยง: Medium

## 10.2 Scenario: Payment → มีหลาย Refunds

### Phase 1: ออกแบบ

```
Payment (Aggregate Root)
├── id
├── amount
└── Refunds []Refund (new)

Refund (Entity ย่อย)
├── id
├── payment_id
├── amount
└── reason
```

**⚠️ หมายเหตุ:** Refund อยู่ใน aggregate ของ Payment

### Phase 2: Domain Entity

```go
// internal/modules/payment/domain/entity/payment.go

type Payment struct {
    ID       uuid.UUID `json:"id"`
    // ... existing
    Refunds  []Refund  `json:"refunds,omitempty"`
}

type Refund struct {
    ID        uuid.UUID       `json:"id"`
    PaymentID uuid.UUID       `json:"payment_id"`
    Amount    decimal.Decimal `json:"amount"`
    Reason    string          `json:"reason"`
    CreatedAt time.Time       `json:"created_at"`
}

// NewRefund creates a new refund (validate against payment)
func NewRefund(paymentID uuid.UUID, amount decimal.Decimal, reason string) (*Refund, error) {
    if paymentID == uuid.Nil {
        return nil, domainerrors.ErrInvalidPaymentID
    }
    if amount.LessThanOrEqual(decimal.Zero) {
        return nil, domainerrors.ErrInvalidAmount
    }
    return &Refund{
        ID:        uuid.New(),
        PaymentID: paymentID,
        Amount:    amount,
        Reason:    reason,
        CreatedAt: time.Now(),
    }, nil
}
```

### Phase 3: Business Method (ผ่าน Aggregate Root)

```go
// Payment.AddRefund — เข้าผ่าน aggregate root
func (p *Payment) AddRefund(amount decimal.Decimal, reason string) (*Refund, error) {
    // Rule: refund ได้เฉพาะ SUCCESS
    if p.Status != valueobject.PaymentStatusSuccess {
        return nil, domainerrors.ErrCannotRefund
    }
    
    // Rule: refund รวมห้ามเกิน amount
    totalRefunded := decimal.Zero
    for _, r := range p.Refunds {
        totalRefunded = totalRefunded.Add(r.Amount)
    }
    if totalRefunded.Add(amount).GreaterThan(p.Amount) {
        return nil, domainerrors.ErrRefundExceedsAmount
    }
    
    // Rule: refund ได้ครั้งเดียว (ถ้าต้องการ)
    if len(p.Refunds) > 0 {
        return nil, domainerrors.ErrAlreadyRefunded
    }
    
    refund, err := NewRefund(p.ID, amount, reason)
    if err != nil {
        return nil, err
    }
    p.Refunds = append(p.Refunds, *refund)
    p.UpdatedAt = time.Now()
    return refund, nil
}
```

### Phase 4: Migration

```sql
-- migrations/20260420_payment_refunds.sql
CREATE TABLE IF NOT EXISTS payment_refunds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL,
    amount NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    reason TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_payment_refunds_payment
        FOREIGN KEY (payment_id)
        REFERENCES payment_transactions(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_payment_refunds_payment_id
    ON payment_refunds(payment_id);
```

### Phase 5: GORM Model

```go
type RefundModel struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    PaymentID uuid.UUID `gorm:"type:uuid;index;not null"`
    Amount    string    `gorm:"type:numeric(15,2);not null"`
    Reason    string    `gorm:"type:text;not null"`
    CreatedAt time.Time `gorm:"default:now()"`
}

func (RefundModel) TableName() string {
    return "payment_refunds"
}
```

### Phase 6: Repository (Aggregate ต้อง Save ทั้งก้อน)

```go
// Save — ใช้ GORM association
func (r *paymentRepoImpl) Save(ctx context.Context, p *entity.Payment) error {
    m := r.toModel(p)
    
    // GORM จะ save association อัตโนมัติถ้าใช้ FullSaveAssociations
    return r.db.WithContext(ctx).
        Session(&gorm.Session{FullSaveAssociations: true}).
        Create(m).Error
}

// FindByID — preload refunds
func (r *paymentRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
    var m PaymentModel
    err := r.db.WithContext(ctx).
        Preload("Refunds").
        First(&m, "id = ?", id).Error
    // ...
}
```

### Phase 7: Use Case

```go
// internal/modules/payment/application/refund_payment.go

type RefundPaymentUseCase struct { /* ... */ }

type RefundPaymentInput struct {
    PaymentID uuid.UUID
    UserID    uuid.UUID
    Amount    decimal.Decimal  // ⬅️ เพิ่มใหม่ (partial refund)
    Reason    string
}

func (uc *RefundPaymentUseCase) Execute(ctx context.Context, input RefundPaymentInput) error {
    payment, err := uc.repo.FindByID(ctx, input.PaymentID)
    if err != nil {
        return err
    }
    
    if payment.UserID != input.UserID {
        return domainerrors.ErrUnauthorized
    }
    
    // ⬅️ เรียก aggregate method
    if _, err := payment.AddRefund(input.Amount, input.Reason); err != nil {
        return err
    }
    
    return uc.repo.Update(ctx, payment)
}
```

### 10.3 Checklist

- [ ] ออกแบบ aggregate ถูกต้อง
- [ ] Entity ย่อยมี constructor
- [ ] Business logic ผ่าน aggregate root
- [ ] Migration มี FK + ON DELETE
- [ ] GORM model ใช้ association
- [ ] Repository save + preload
- [ ] Test ครอบกรณี: refund ปกติ, refund เกิน, refund ซ้ำ
- [ ] Documentation

---

# บทที่ 11: Breaking Change

## 11.1 นิยาม

**Breaking Change** = การเปลี่ยนแปลงที่ทำให้ code/client ที่พึ่งพาเดิมใช้ไม่ได้

**ตัวอย่าง:**
- เปลี่ยน field name ใน JSON response
- เปลี่ยน type ของ field
- ลบ field/endpoint
- เปลี่ยน semantics ของ field
- เพิ่ม required field ใน request

## 11.2 หลักการ 3 ข้อ

### 1. Version the API

```
/api/v1/payments → legacy (ยังทำงาน)
/api/v2/payments → new (breaking change)
```

### 2. Support Both for a Period

```
Day 1:     Deploy v2
Day 1-180: Support v1 + v2 (parallel)
Day 180:   Deprecate v1 (sunset header)
Day 365:   Remove v1
```

### 3. Communicate Early

```
T-90 days: Announce deprecation
T-60 days: Send reminders
T-30 days: Final warning
T-7 days:  Last call
T-0 days:  Sunset
```

## 11.3 ตัวอย่าง: เปลี่ยน Response Structure

**Before (v1):**
```json
{
  "id": "abc-123",
  "amount": 100.00,
  "user": {
    "id": "user-1",
    "name": "John"
  }
}
```

**After (v2):**
```json
{
  "id": "abc-123",
  "amount": {
    "value": "100.00",
    "currency": "THB"
  },
  "customer": {
    "id": "user-1",
    "full_name": "John Doe"
  }
}
```

**Breaking changes:**
- `amount: 100.00` → `amount: {value, currency}`
- `user` → `customer`
- `name` → `full_name`

## 11.4 Implementation

### 11.4.1 Handler แยก Version

```go
// internal/modules/payment/interfaces/http/v1/handler.go
package v1

type PaymentHandler struct { /* ... */ }

func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
    // ... legacy logic
    response := PaymentResponseV1{ /* ... */ }
    httputil.JSON(w, http.StatusCreated, response)
}
```

```go
// internal/modules/payment/interfaces/http/v2/handler.go
package v2

type PaymentHandler struct { /* ... */ }

func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
    // ... new logic
    response := PaymentResponseV2{ /* ... */ }
    httputil.JSON(w, http.StatusCreated, response)
}
```

### 11.4.2 Routes

```go
// internal/modules/payment/interfaces/http/routes.go

func RegisterRoutes(r chi.Router, h *Handlers, authMW func(http.Handler) http.Handler) {
    // v1 — LEGACY (deprecation 2026-10-15)
    r.Route("/api/v1/payments", func(r chi.Router) {
        r.Use(authMW)
        r.Use(deprecationMiddleware("2026-10-15", "/api/v2/payments"))
        r.Post("/", h.V1Payment.Create)
    })
    
    // v2 — CURRENT
    r.Route("/api/v2/payments", func(r chi.Router) {
        r.Use(authMW)
        r.Post("/", h.V2Payment.Create)
    })
}
```

### 11.4.3 Deprecation Middleware

```go
func deprecationMiddleware(sunsetDate, replacement string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Add deprecation headers
            w.Header().Set("Deprecation", "true")
            w.Header().Set("Sunset", sunsetDate)
            w.Header().Set("Link", fmt.Sprintf("</%s>; rel=\"successor-version\"", replacement))
            
            // Log usage for tracking
            log.Warn("deprecated endpoint used",
                "path", r.URL.Path,
                "user_agent", r.UserAgent(),
                "client_ip", r.RemoteAddr,
            )
            
            next.ServeHTTP(w, r)
        })
    }
}
```

## 11.5 Communication Template

```markdown
Subject: [ACTION REQUIRED] API Deprecation Notice: /api/v1/payments

Dear Partner,

We're writing to inform you that the following API endpoint will be 
deprecated on **October 15, 2026**:

- **Endpoint**: `POST /api/v1/payments`
- **Replacement**: `POST /api/v2/payments`
- **Reason**: Security improvements + better data structure

## What's Changing

| v1 | v2 |
|---|---|
| `amount: 100.00` | `amount: {value: "100.00", currency: "THB"}` |
| `user` | `customer` |
| `name` | `full_name` |

## What You Need To Do

1. Update your integration to use `/api/v2/payments`
2. Test in sandbox: https://sandbox.api.com/api/v2/payments
3. Deploy to production before October 15, 2026

## Support

- Migration guide: https://docs.api.com/migration/v1-to-v2
- Sandbox: https://sandbox.api.com
- Support: api-support@company.com

## Timeline

- **2026-04-15**: Deprecation announced (today)
- **2026-07-15**: v1 rate limits reduced
- **2026-10-15**: v1 sunset (returns 410 Gone)

Thank you,
API Team
```

## 11.6 Checklist

- [ ] Deprecation policy documented
- [ ] Both versions implemented
- [ ] Deprecation headers in response
- [ ] Usage metrics (track who uses old)
- [ ] Client inventory (รู้ว่าใครใช้อยู่)
- [ ] Communication plan
- [ ] Migration guide
- [ ] Sandbox environment
- [ ] Support channel
- [ ] Sunset timeline
- [ ] Auto-reject after sunset (410 Gone)

---

# บทที่ 12: Refactor โครงสร้าง

## 12.1 ระดับความเสี่ยง: Medium-High

**ทำไมเสี่ยง?** เพราะ refactor = เปลี่ยนโค้ดจำนวนมาก โดยที่ behavior ไม่เปลี่ยน

**หลักการ:** "Refactor ในขณะที่ test ผ่านเท่านั้น"

## 12.2 Scenario: แยก God Service เป็น Use Cases

**Before:**
```go
// 500 บรรทัด ในไฟล์เดียว
type PaymentService struct { /* ... */ }

func (s *PaymentService) Process(ctx context.Context, req Request) error {
    // 1. Validate
    // 2. Create payment
    // 3. Call provider
    // 4. Update status
    // 5. Send email
    // 6. Publish event
    // 7. Log
    // ... 500 บรรทัด
}
```

**After:**
```
application/
├── create_payment.go
├── process_payment.go
├── refund_payment.go
└── notify_customer.go
```

## 12.3 Process

### Step 1: เขียน Characterization Tests

```go
// ก่อน refactor — test ครอบ behavior ทั้งหมด
func TestPaymentService_Process_Characterization(t *testing.T) {
    tests := []struct{
        name     string
        input    Request
        wantErr  error
        wantState string
    }{
        {"happy path", ..., nil, "SUCCESS"},
        {"invalid amount", ..., domainerrors.ErrInvalidAmount, "PENDING"},
        {"provider fail", ..., domainerrors.ErrProviderError, "FAILED"},
        // ... 20 cases
    }
    // ...
}
```

**⚠️ Test นี้ต้องผ่าน 100% ทั้งก่อนและหลัง refactor**

### Step 2: แยกทีละ Use Case

**Step 2a: แยก `CreatePayment`**
```go
// application/create_payment.go
type CreatePaymentUseCase struct { /* ... */ }
func (uc *CreatePaymentUseCase) Execute(ctx context.Context, input CreatePaymentInput) (*CreatePaymentOutput, error) {
    // ... logic ที่ 1+2
}

// PaymentService.Process เรียก CreatePaymentUseCase
func (s *PaymentService) Process(ctx context.Context, req Request) error {
    output, err := s.createUC.Execute(ctx, CreatePaymentInput{ /* ... */ })
    // ... ที่เหลือ
}
```

**Step 2b: แยก `ProcessPayment`** ...

**Step 2c: แยก `NotifyCustomer`** ...

### Step 3: เปลี่ยน Handler มาเรียกใช้ Use Cases ตรง

```go
// เดิม: handler → PaymentService.Process
// ใหม่: handler → use case โดยตรง

func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
    output, err := h.createUC.Execute(ctx, input)
    // ...
}
```

### Step 4: ลบ God Service (หลัง 2 sprints)

```bash
git rm internal/modules/payment/application/payment_service.go
```

## 12.4 Refactoring Patterns

| Pattern | เมื่อไหร่ | ตัวอย่าง |
|---|---|---|
| **Extract Method** | Function ยาว | แยก validation ออก |
| **Extract Class** | Class ทำหลายอย่าง | แยก God Service |
| **Move Method** | Method อยู่ผิดที่ | ย้าย behavior เข้า entity |
| **Rename** | ชื่อไม่สื่อ | `Process` → `Capture` |
| **Inline** | Abstraction ไม่จำเป็น | ลบ wrapper |
| **Replace Conditional** | if-else ยาว | ใช้ polymorphism |

## 12.5 Checklist

- [ ] Characterization test ครอบ behavior เดิม 100%
- [ ] Test ผ่านก่อน refactor
- [ ] Refactor ทีละเล็ก (commit เล็ก)
- [ ] Test ผ่านหลังแต่ละ commit
- [ ] ไม่เปลี่ยน public API (หรือ version ถ้าเปลี่ยน)
- [ ] Documentation อัปเดต
- [ ] ลบ deprecated code (หลัง observe)

---

# บทที่ 13: ย้ายจาก Sync → Async

## 13.1 ระดับความเสี่ยง: Medium-High

**ทำไมต้องย้าย?**
- Response time ช้า (blocking)
- ไม่ทนต่อ failure (provider ล่ม → user รอ)
- Scale ยาก (ต้องการ concurrency สูง)

## 13.2 Scenario: ย้าย Email Sending จาก Sync → Async

### Before

```go
func (uc *CreatePaymentUseCase) Execute(ctx context.Context, input CreatePaymentInput) (*CreatePaymentOutput, error) {
    // ... create payment
    
    // ⚠️ Blocking call — user รอ 2-3 วินาที
    if err := uc.emailService.Send(ctx, payment.UserID, "payment.created"); err != nil {
        uc.logger.Warn("email send failed", "error", err)
        // ไม่ fail หลัก แต่ response ช้า
    }
    
    return output, nil
}
```

**ปัญหา:**
- Response time 3s (เพิ่ม 2s จาก email)
- ถ้า SMTP ล่ม → user รอนาน
- ถ้า SMTP fail → payment สร้างแล้ว แต่ user ไม่รู้

### After

```go
func (uc *CreatePaymentUseCase) Execute(ctx context.Context, input CreatePaymentInput) (*CreatePaymentOutput, error) {
    // ... create payment
    
    // ✅ Async — publish event
    event := events.PaymentCreated{
        BaseEvent: event.NewBaseEvent(),
        PaymentID: payment.ID,
        UserID:    payment.UserID,
        OrderID:   payment.OrderID,
        Amount:    payment.Amount,
        Currency:  payment.Currency,
    }
    if err := uc.eventBus.Publish(ctx, "payment.created", event); err != nil {
        // Log แต่ไม่ fail — event bus จะ retry
        uc.logger.Warn("failed to publish event", "error", err, "payment_id", payment.ID)
    }
    
    return output, nil  // ← response ทันที (< 100ms)
}
```

### Consumer

```go
// internal/modules/payment/infrastructure/messaging/consumers/payment_created_consumer.go
package consumers

type PaymentCreatedConsumer struct {
    emailService EmailService
    logger       logger.Logger
}

func (c *PaymentCreatedConsumer) Handle(ctx context.Context, payload map[string]interface{}) error {
    var event events.PaymentCreated
    if err := mapstructure.Decode(payload, &event); err != nil {
        return err
    }
    
    // Idempotency: check if email already sent
    if c.emailService.WasSent(ctx, event.PaymentID) {
        return nil
    }
    
    if err := c.emailService.Send(ctx, event.UserID, "payment.created"); err != nil {
        c.logger.Error("email failed", "error", err, "payment_id", event.PaymentID)
        return err  // retry
    }
    
    return nil
}
```

### Migration Plan

```
Phase 1: Add Kafka producer + consumer (แต่ยัง sync)
Phase 2: Dual mode — publish event + ส่ง email sync
Phase 3: Switch → publish event only
Phase 4: Remove sync email code
```

### Idempotency

```go
// ต้องกัน duplicate (Kafka at-least-once)
type EmailService struct {
    redis *redis.Client
}

func (s *EmailService) WasSent(ctx context.Context, paymentID uuid.UUID) bool {
    key := fmt.Sprintf("email:sent:%s", paymentID)
    n, _ := s.redis.Exists(ctx, key).Result()
    return n > 0
}

func (s *EmailService) MarkSent(ctx context.Context, paymentID uuid.UUID) error {
    key := fmt.Sprintf("email:sent:%s", paymentID)
    return s.redis.Set(ctx, key, "1", 7*24*time.Hour).Err()
}
```

### Checklist

- [ ] Kafka producer setup
- [ ] Consumer with idempotency
- [ ] Dead letter queue
- [ ] Retry policy (exponential backoff)
- [ ] Monitoring (lag, failure rate)
- [ ] Dual write period
- [ ] Switch phase
- [ ] Remove sync code
- [ ] Alert on consumer lag

---

# บทที่ 14: Expand-Contract Pattern

## 14.1 หลักการ

```
┌───────────┐    ┌───────────┐    ┌───────────┐
│  EXPAND   │───▶│  MIGRATE  │───▶│  CONTRACT │
├───────────┤    ├───────────┤    ├───────────┤
│ Add new   │    │ Move data │    │ Remove    │
│ Both work │    │ Verify    │    │ old       │
│ Backward  │    │ Switch    │    │ Clean     │
│ compatible│    │ read/write│    │           │
└───────────┘    └───────────┘    └───────────┘
  Safe             Careful          Deprecate first
```

## 14.2 เมื่อไหร่ใช้

- เปลี่ยน type ของ column
- เปลี่ยนชื่อ column
- แยก table
- รวม table
- เปลี่ยน index strategy

## 14.3 ตัวอย่างการใช้งาน 4 Phase

**Phase 1: EXPAND** — เพิ่มของใหม่ (ของเก่ายังทำงาน)
```sql
ALTER TABLE payment_transactions
    ADD COLUMN new_column TYPE DEFAULT value;
```

**Phase 2: DUAL** — เขียนทั้ง 2 (ยังอ่านเก่า)
```go
// Save
model.OldColumn = data
model.NewColumn = data

// Read (fallback)
value := model.NewColumn
if value == "" {
    value = model.OldColumn
}
```

**Phase 3: SWITCH** — อ่านจากใหม่ (ยังเขียนเก่า)
```go
value := model.NewColumn  // ใช้ใหม่เท่านั้น
```

**Phase 4: CONTRACT** — ลบของเก่า
```sql
ALTER TABLE payment_transactions DROP COLUMN old_column;
```

## 14.4 Timeline ตัวอย่าง

| Phase | ระยะเวลา | Activity |
|---|---|---|
| Expand | Day 1 | Add column |
| Dual | Day 2-3 | Deploy dual code |
| Backfill | Day 4 | Migrate data |
| Verify | Day 5-11 | Monitor |
| Switch Read | Day 12 | Deploy new read |
| Verify | Day 12-25 | Monitor |
| Switch Write | Day 26 | Deploy |
| Verify | Day 26-39 | Monitor |
| Contract | Day 40 | Drop old |

## 14.5 Checklist

- [ ] แต่ละ phase มี rollback plan
- [ ] Monitor ระหว่าง transition
- [ ] Alert on mismatch
- [ ] Documentation per phase
- [ ] Communication plan

---

# บทที่ 15: Zero-Downtime Migration

## 15.1 กฎทอง

```
1. ห้าม lock table นาน (> 1s)
2. ห้าม ALTER ที่ REWRITE table
3. ห้าม UPDATE ทั้งตารางทีเดียว
4. ใช้ CONCURRENTLY สำหรับ index
5. Test บน staging ที่ data ≈ production
```

## 15.2 Safe vs Unsafe Operations

| Operation | Safe | Unsafe |
|---|---|---|
| **ADD COLUMN** (nullable) | ✅ Fast | – |
| **ADD COLUMN** (NOT NULL + default) | ✅ Postgres 11+ | ❌ Postgres < 11 |
| **DROP COLUMN** | ✅ Fast | – |
| **RENAME COLUMN** | ❌ Breaks code | – |
| **ALTER TYPE** | ❌ REWRITE | – |
| **CREATE INDEX** | ❌ Blocks writes | – |
| **CREATE INDEX CONCURRENTLY** | ✅ No lock | – |
| **ADD CHECK CONSTRAINT** | ❌ Scans | ✅ NOT VALID + VALIDATE |
| **ADD FK** | ❌ Scans | ✅ NOT VALID + VALIDATE |
| **UPDATE (bulk)** | ❌ Locks | ✅ Batch + SKIP LOCKED |

## 15.3 ตัวอย่าง: เพิ่ม Index แบบ Zero-Downtime

**❌ Unsafe:**
```sql
-- Block writes ทั้งตาราง (อาจใช้เวลา 30 นาที)
CREATE INDEX idx_payment_user_id ON payment_transactions(user_id);
```

**✅ Safe:**
```sql
-- Không block writes
CREATE INDEX CONCURRENTLY idx_payment_user_id 
    ON payment_transactions(user_id);
```

**⚠️ หมายเหตุ:**
- `CONCURRENTLY` ใช้เวลา 2-3x นานกว่า แต่ไม่ block
- ห้ามใช้ใน transaction
- ถ้า fail → ต้อง DROP INDEX ที่ค้าง

## 15.4 ตัวอย่าง: เพิ่ม NOT NULL Column ปลอดภัย

**❌ Unsafe (Postgres < 11):**
```sql
-- REWRITE table ทั้งหมด
ALTER TABLE payment_transactions
    ADD COLUMN new_col TEXT NOT NULL DEFAULT '';
```

**✅ Safe:**
```sql
-- Step 1: Add nullable + default (fast)
ALTER TABLE payment_transactions
    ADD COLUMN new_col TEXT DEFAULT '';

-- Step 2: Backfill batch
DO $$
DECLARE
    batch_size INT := 10000;
    affected INT;
BEGIN
    LOOP
        UPDATE payment_transactions
        SET new_col = 'default_value'
        WHERE new_col IS NULL
          AND id IN (
              SELECT id FROM payment_transactions
              WHERE new_col IS NULL
              LIMIT batch_size
              FOR UPDATE SKIP LOCKED
          );
        GET DIAGNOSTICS affected = ROW_COUNT;
        EXIT WHEN affected = 0;
        PERFORM pg_sleep(0.1);
    END LOOP;
END $$;

-- Step 3: Add constraint (validate separately)
ALTER TABLE payment_transactions
    ADD CONSTRAINT new_col_not_null 
    CHECK (new_col IS NOT NULL) NOT VALID;

-- Step 4: Validate (ไม่ block)
ALTER TABLE payment_transactions
    VALIDATE CONSTRAINT new_col_not_null;

-- Step 5: Convert to NOT NULL (ถ้าต้องการ)
ALTER TABLE payment_transactions
    ALTER COLUMN new_col SET NOT NULL,
    ALTER COLUMN new_col DROP DEFAULT;
```

## 15.5 Testing Zero-Downtime บน Staging

```bash
# Terminal 1: รัน migration
psql -f migrations/xxx.sql

# Terminal 2: จำลอง user activity ระหว่าง migration
while true; do
    psql -c "INSERT INTO payment_transactions (...) VALUES (...);"
    psql -c "SELECT * FROM payment_transactions LIMIT 10;"
    sleep 0.1
done

# Terminal 3: วัด latency
pgbench -c 10 -j 2 -T 60 -f query.sql dbname
```

**Acceptance criteria:**
- INSERT latency < 100ms ตลอด
- SELECT latency < 50ms ตลอด
- ไม่มี deadlock
- ไม่มี lock timeout

## 15.6 Checklist

- [ ] Migration ใช้ CONCURRENTLY สำหรับ index
- [ ] ไม่มี ALTER ที่ REWRITE
- [ ] Backfill เป็น batch
- [ ] Test บน staging (data ≈ production)
- [ ] วัด latency ระหว่าง migration
- [ ] Rollback plan
- [ ] Monitor lock & deadlock

---

# บทที่ 16: Backfill ข้อมูลขนาดใหญ่

## 16.1 ปัญหา

**ตาราง 100M rows × 100 bytes = 10 GB** — UPDATE ทั้งตาราง = lock นาน = production ล่ม

## 16.2 เทคนิค Batch Backfill

### 16.2.1 Basic Batch

```sql
DO $$
DECLARE
    batch_size INT := 10000;
    affected INT;
    total INT := 0;
BEGIN
    LOOP
        UPDATE payment_transactions
        SET new_field = compute_value(old_field)
        WHERE new_field IS NULL
          AND id IN (
              SELECT id FROM payment_transactions
              WHERE new_field IS NULL
              LIMIT batch_size
          );
        GET DIAGNOSTICS affected = ROW_COUNT;
        total := total + affected;
        RAISE NOTICE 'Batch done: %, total: %', affected, total;
        EXIT WHEN affected < batch_size;
        PERFORM pg_sleep(0.1);
    END LOOP;
    RAISE NOTICE 'Complete. Total: %', total;
END $$;
```

### 16.2.2 Batch + SKIP LOCKED (แนะนำ)

```sql
DO $$
DECLARE
    batch_size INT := 10000;
    affected INT := 1;
    total INT := 0;
BEGIN
    WHILE affected > 0 LOOP
        WITH batch AS (
            SELECT id, old_field
            FROM payment_transactions
            WHERE new_field IS NULL
            LIMIT batch_size
            FOR UPDATE SKIP LOCKED  -- ← ไม่ block transaction อื่น
        )
        UPDATE payment_transactions pt
        SET new_field = batch.old_field::type
        FROM batch
        WHERE pt.id = batch.id;
        
        GET DIAGNOSTICS affected = ROW_COUNT;
        total := total + affected;
        
        RAISE NOTICE 'Progress: % rows (total: %)', affected, total;
        
        COMMIT;
        PERFORM pg_sleep(0.1);
    END LOOP;
END $$;
```

### 16.2.3 Ordered by ID (หลีกเลี่ยง hot spot)

```sql
DO $$
DECLARE
    last_id UUID := '00000000-0000-0000-0000-000000000000';
    batch_size INT := 10000;
    affected INT;
BEGIN
    LOOP
        WITH batch AS (
            SELECT id, old_field
            FROM payment_transactions
            WHERE id > last_id
              AND new_field IS NULL
            ORDER BY id
            LIMIT batch_size
        )
        UPDATE payment_transactions pt
        SET new_field = batch.old_field::type
        FROM batch
        WHERE pt.id = batch.id
        RETURNING pt.id INTO last_id;
        
        GET DIAGNOSTICS affected = ROW_COUNT;
        EXIT WHEN affected < batch_size;
        
        COMMIT;
        PERFORM pg_sleep(0.1);
    END LOOP;
END $$;
```

## 16.3 Backfill ด้วย Go (สำหรับ complex logic)

```go
// cmd/backfill/main.go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "gorm.io/gorm"
)

func main() {
    ctx := context.Background()
    db := setupDB()
    
    const batchSize = 10000
    var lastID string
    total := 0
    
    for {
        var rows []PaymentModel
        err := db.WithContext(ctx).
            Where("new_field IS NULL").
            Where("id > ?", lastID).
            Order("id").
            Limit(batchSize).
            Find(&rows).Error
        if err != nil {
            log.Fatal(err)
        }
        if len(rows) == 0 {
            break
        }
        
        // Compute + batch update
        updates := make([]map[string]interface{}, len(rows))
        for i, r := range rows {
            updates[i] = map[string]interface{}{
                "id":        r.ID,
                "new_field": computeValue(r.OldField),
            }
            lastID = r.ID.String()
        }
        
        // Upsert batch
        if err := db.Table("payment_transactions").
            Save(&updates).Error; err != nil {
            log.Fatal(err)
        }
        
        total += len(rows)
        fmt.Printf("Progress: %d rows (last_id: %s)\n", total, lastID)
        
        time.Sleep(100 * time.Millisecond)
    }
    
    fmt.Printf("Complete. Total: %d rows\n", total)
}
```

## 16.4 Verify Backfill

```sql
-- ต้องได้ 0
SELECT COUNT(*) FROM payment_transactions WHERE new_field IS NULL;

-- เทียบกับ old_field
SELECT 
    COUNT(*) AS total,
    COUNT(*) FILTER (WHERE new_field = old_field::type) AS matches,
    COUNT(*) FILTER (WHERE new_field != old_field::type) AS mismatches
FROM payment_transactions;
```

## 16.5 Monitoring ระหว่าง Backfill

```promql
# Grafana
# Row progression
pg_stat_user_tables_n_tup_upd{relname="payment_transactions"}

# Lock activity
pg_locks_count{locktype="transactionid"}

# DB connection
pg_stat_activity_count{state="active"}
```

## 16.6 Checklist

- [ ] Batch size เหมาะสม (10K-100K)
- [ ] ใช้ SKIP LOCKED
- [ ] Sleep ระหว่าง batch
- [ ] Resume ได้ (save last_id)
- [ ] Verify 100% complete
- [ ] Monitor latency ระหว่าง backfill
- [ ] Rollback plan (ก่อน commit final)

---

# บทที่ 17: API Versioning

## 17.1 รูปแบบของ Versioning

| รูปแบบ | ตัวอย่าง | ข้อดี | ข้อเสีย |
|---|---|---|---|
| **URI Path** | `/api/v1/payments` | ชัดเจน, cache ได้ | URI เปลี่ยน |
| **Query String** | `/api/payments?version=1` | URI เดิม | cache ยาก |
| **Header** | `Accept: application/vnd.api+json;version=1` | URI เดิม | debug ยาก |
| **Content Negotiation** | `Accept: application/vnd.company.v1+json` | มาตรฐาน | ซับซ้อน |

**แนะนำ: URI Path** — ชัดเจนที่สุดใน Go/Chi

## 17.2 โครงสร้าง Code

```
internal/modules/payment/interfaces/http/
├── v1/
│   ├── handler.go
│   ├── dto.go
│   └── routes.go
├── v2/
│   ├── handler.go
│   ├── dto.go
│   └── routes.go
└── shared/
    └── errors.go  # error mapping ใช้ร่วมกัน
```

## 17.3 เมื่อไหร่ต้อง Version ใหม่

**ต้อง:**
- ลบ field ออกจาก response
- เปลี่ยน type ของ field
- เปลี่ยน semantics ของ field
- เปลี่ยน URL structure
- เพิ่ม required field ใน request

**ไม่ต้อง:**
- เพิ่ม optional field
- เพิ่ม endpoint ใหม่
- เพิ่ม field ใน response (additive)
- แก้ bug

## 17.4 Deprecation Timeline

```
T-0:      Announce deprecation
T+30d:    Reduce rate limits on old version
T+60d:    Send reminders (email, dashboard banner)
T+90d:    Add Sunset header
T+180d:   Return 410 Gone with sunset message
T+365d:   Remove code
```

## 17.5 ตัวอย่าง Response Headers

```http
HTTP/1.1 200 OK
Content-Type: application/json
Deprecation: true
Sunset: Sat, 15 Oct 2026 00:00:00 GMT
Link: </api/v2/payments>; rel="successor-version"
Warning: 299 - "This endpoint is deprecated. Use /api/v2/payments"
```

## 17.6 Client Tracking

```go
// Middleware to track version usage
func versionTracker(version string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Increment counter
            APIVersionUsage.WithLabelValues(version, r.URL.Path).Inc()
            
            // Log user agent
            log.Info("api call",
                "version", version,
                "path", r.URL.Path,
                "user_agent", r.UserAgent(),
                "api_key", r.Header.Get("X-API-Key"),
            )
            
            next.ServeHTTP(w, r)
        })
    }
}
```

## 17.7 Checklist

- [ ] Versioning strategy เลือกแล้ว
- [ ] Code structure แยก v1/v2
- [ ] Deprecation headers
- [ ] Client tracking
- [ ] Sunset date ประกาศแล้ว
- [ ] Migration guide
- [ ] Communication plan

---

# บทที่ 18: Deprecation Strategy

## 18.1 หลักการ 3 ข้อ

1. **Announce early** — 90+ days ล่วงหน้า
2. **Make it easy** — migration guide, sandbox
3. **Follow through** — sunset ตามที่สัญญา

## 18.2 Deprecation Levels

| Level | ความหมาย | Response |
|---|---|---|
| **Info** | ยังใช้ได้, ควรเปลี่ยน | Header only |
| **Warn** | ยังใช้ได้, ควรเปลี่ยนเร็ว | Header + log |
| **Sunset** | ปิดเมื่อ date | 410 Gone |
| **Removed** | ลบแล้ว | 404 Not Found |

## 18.3 Response Format

```json
{
  "error": "endpoint_deprecated",
  "message": "This endpoint is deprecated and will be removed on 2026-10-15. Please migrate to /api/v2/payments.",
  "sunset_date": "2026-10-15",
  "replacement": "/api/v2/payments",
  "migration_guide": "https://docs.api.com/migration/v1-to-v2"
}
```

## 18.4 Communication Channels

| ช่องทาง | ใช้เมื่อ | ความถี่ |
|---|---|---|
| **Email** | แจ้งตรงถึง client | 3 ครั้ง (T-90, T-30, T-7) |
| **Dashboard banner** | แสดงใน dashboard ของ client | ตลอดเวลา |
| **Changelog** | Public announcement | 1 ครั้ง |
| **Status page** | ข้อมูลเพิ่ม | 1 ครั้ง |
| **Direct call** | Key accounts | 2 ครั้ง |
| **Response header** | ทุก API call | ทุกครั้ง |

## 18.5 Migration Guide Template

```markdown
# Migration Guide: v1 → v2

## Overview
This guide helps you migrate from `/api/v1/payments` to `/api/v2/payments`.

## Breaking Changes

### 1. Amount Structure

**Before (v1):**
```json
{"amount": 100.00}
```

**After (v2):**
```json
{"amount": {"value": "100.00", "currency": "THB"}}
```

**Migration:**
```javascript
// Before
const amount = response.amount;

// After
const amount = response.amount.value;
const currency = response.amount.currency;
```

### 2. User → Customer

**Before (v1):**
```json
{"user": {"id": "...", "name": "John"}}
```

**After (v2):**
```json
{"customer": {"id": "...", "full_name": "John Doe"}}
```

### 3. Error Format

...

## Testing

Sandbox: https://sandbox.api.com/api/v2

## Support

Email: api-support@company.com
Slack: #api-migration
```

## 18.6 Checklist

- [ ] Deprecation policy documented
- [ ] Sunset date announced
- [ ] Migration guide published
- [ ] Sandbox available
- [ ] Support channel open
- [ ] Tracking in place
- [ ] Emails sent (T-90, T-30, T-7)
- [ ] Deprecation headers
- [ ] 410 Gone after sunset

---

# บทที่ 19: Regression Testing

## 19.1 Test Pyramid สำหรับการแก้ Module

```
         /\
        /E2E\        5%   (smoke test หลัง deploy)
       /──────\
      / Integr \    25%   (repository, Kafka)
     /──────────\
    /   Unit     \ 70%   (domain, use case)
   /──────────────\
```

## 19.2 Regression Test Suite

**ต้องมี test ครบก่อนแก้:**

| ประเภท | จำนวน | เครื่องมือ |
|---|---|---|
| Domain Unit Tests | 50-200 | testing + testify |
| Use Case Tests | 30-100 | + mock |
| Integration Tests | 10-50 | Testcontainers |
| API Tests | 20-50 | httptest |
| E2E Tests | 5-20 | Postman/k6 |

## 19.3 Characterization Test

**ก่อนแก้ → lock behavior ปัจจุบัน:**

```go
// internal/modules/payment/application/create_payment_test.go
// Characterization — บันทึก behavior ก่อนแก้

func TestCreatePayment_Characterization(t *testing.T) {
    cases := []struct {
        name    string
        input   CreatePaymentInput
        wantErr error
        wantAmt string
    }{
        // Happy path
        {"standard amount", ..., nil, "100"},
        {"decimal amount", ..., nil, "99.99"},
        {"max amount", ..., nil, "999999.99"},
        {"min amount", ..., nil, "0.01"},
        
        // Error cases
        {"zero amount", ..., domainerrors.ErrInvalidAmount, ""},
        {"negative amount", ..., domainerrors.ErrInvalidAmount, ""},
        {"empty currency", ..., domainerrors.ErrInvalidCurrency, ""},
        {"invalid method", ..., domainerrors.ErrInvalidPaymentMethod, ""},
        
        // Edge cases
        {"duplicate order", ..., domainerrors.ErrPaymentAlreadyExists, ""},
        {"nil user", ..., domainerrors.ErrInvalidUserID, ""},
    }
    
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            // ... run
            // assert เหมือนก่อนแก้
        })
    }
}
```

**หลักการ:** Test นี้อาจจะไม่ "ถูกต้อง" ทั้งหมด แต่ lock behavior ไว้ → ถ้า change → fail → บังคับให้คิด

## 19.4 Smoke Test หลัง Deploy

```bash
#!/bin/bash
# scripts/smoke-test.sh

BASE_URL="${1:-http://localhost:8080}"
TOKEN="..."

# 1. Health check
curl -f "$BASE_URL/health" || exit 1

# 2. Create payment
ORDER_ID=$(uuidgen)
curl -f -X POST "$BASE_URL/api/v1/payments" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"order_id\":\"$ORDER_ID\",\"amount\":\"100\",\"currency\":\"THB\",\"method\":\"CARD\"}" \
    || exit 1

# 3. Get payment
curl -f "$BASE_URL/api/v1/payments/$PAYMENT_ID" \
    -H "Authorization: Bearer $TOKEN" \
    || exit 1

echo "✅ Smoke test passed"
```

## 19.5 Load Test (Regression)

```javascript
// tests/load/regression.js (k6)
import http from 'k6/http';
import { check } from 'k6';

export const options = {
    stages: [
        { duration: '1m', target: 10 },
        { duration: '3m', target: 100 },
        { duration: '1m', target: 0 },
    ],
    thresholds: {
        http_req_duration: ['p(95)<200'],  // ← SLA
        http_req_failed: ['rate<0.001'],
    },
};

export default function () {
    const res = http.post(`${__ENV.BASE_URL}/api/v1/payments`, JSON.stringify({
        order_id: `${__VU}-${__ITER}`,
        amount: '100',
        currency: 'THB',
        method: 'CARD',
    }), {
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${__ENV.TOKEN}` },
    });
    
    check(res, {
        'status is 201': (r) => r.status === 201,
        'latency OK': (r) => r.timings.duration < 200,
    });
}
```

## 19.6 Checklist

- [ ] Characterization test ครอบก่อนแก้
- [ ] Unit test coverage ≥ 80%
- [ ] Integration test ผ่าน
- [ ] Smoke test หลัง deploy
- [ ] Load test เทียบ baseline
- [ ] Chaos test (ถ้าจำเป็น)

---

# บทที่ 20: Rollback Strategy

## 20.1 หลักการ

```
"Always have a rollback plan.
 If you can't rollback, don't deploy."
```

## 20.2 ระดับของ Rollback

| ระดับ | ตัวอย่าง | เวลา | Data loss |
|---|---|---|---|
| **Instant** | Feature flag flip | < 30s | ไม่มี |
| **Fast** | Revert deployment | 2-5 min | อาจมี (ช่วงสั้น) |
| **Medium** | Rollback migration | 10-60 min | อาจมี |
| **Slow** | Restore from backup | 1-4 ชม. | มี |
| **Impossible** | ลบ column | – | – |

## 20.3 Blue-Green Deployment

```
       ┌──────────────┐
       │ Load Balancer│
       └──────┬───────┘
              │
       ┌──────┴──────┐
       │             │
   ┌───▼───┐    ┌───▼───┐
   │ Blue  │    │ Green │
   │ (old) │    │ (new) │
   └───────┘    └───────┘
       
   Deploy → Green
   Test → Green
   Switch LB → Green
   Keep Blue for rollback
```

## 20.4 Canary Deployment

```
   100% traffic
        │
        ▼
   ┌─────────┐
   │ Router  │
   └────┬────┘
        │
   ┌────┴────┐
   │         │
  5% │       │ 95%
   ▼         ▼
┌─────┐  ┌─────┐
│New  │  │Old  │
└─────┘  └─────┘

Monitor → ถ้า OK เพิ่ม %
       → ถ้า fail → rollback
```

## 20.5 Rollback Plan Template

```markdown
# Rollback Plan: [Change Title]

## Trigger
- Error rate > 1% (5 min)
- P95 latency > 500ms (5 min)
- Manual: customer complaint

## Options

### Option 1: Feature Flag (Preferred)
- **Action**: Disable `payment.new_rule` flag
- **Command**: `curl -X POST flag-service/api/flags/payment.new_rule -d '{"enabled": false}'`
- **Time**: < 1 minute
- **Data impact**: None

### Option 2: Revert Deployment
- **Action**: `kubectl rollout undo deployment/api`
- **Time**: 2-5 minutes
- **Data impact**: In-flight requests lost

### Option 3: Rollback Migration
- **Action**: `psql -f migrations/xxx.down.sql`
- **Time**: 10-60 minutes
- **Data impact**: New data lost

## Verification
- [ ] Error rate back to normal
- [ ] Latency back to baseline
- [ ] Smoke test passes
- [ ] Alerts clear

## Post-Rollback
- [ ] Investigate root cause
- [ ] Update test to cover case
- [ ] Retry with fix
```

## 20.6 Automated Rollback

```yaml
# Kubernetes
apiVersion: apps/v1
kind: Deployment
spec:
  progressDeadlineSeconds: 300  # 5 min
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
```

```bash
# Auto-rollback if not ready in 5 min
kubectl rollout status deployment/api --timeout=5m
if [ $? -ne 0 ]; then
    kubectl rollout undo deployment/api
fi
```

## 20.7 Checklist

- [ ] Rollback plan documented
- [ ] ทดสอบ rollback บน staging
- [ ] Auto-rollback trigger configured
- [ ] Team รู้วิธี rollback
- [ ] Communication plan
- [ ] Data safety check

---

# บทที่ 21: Feature Flags

## 21.1 ประเภทของ Flags

| ประเภท | อายุ | ตัวอย่าง |
|---|---|---|
| **Release** | สั้น (< 1 sprint) | `payment.new_amount_format` |
| **Experiment** | กลาง (A/B test) | `payment.checkout_v2` |
| **Ops** | ยาว (kill switch) | `payment.provider.stripe.enabled` |
| **Permission** | ถาวร (cohort) | `payment.beta_user` |

## 21.2 Implementation

```go
// pkg/featureflag/flags.go
package featureflag

import (
    "context"
    "fmt"
    
    "github.com/redis/go-redis/v9"
)

type Flags interface {
    IsEnabled(ctx context.Context, name string, userID string) bool
    Set(ctx context.Context, name string, enabled bool) error
}

type redisFlags struct {
    client *redis.Client
}

func NewRedisFlags(client *redis.Client) Flags {
    return &redisFlags{client: client}
}

func (f *redisFlags) IsEnabled(ctx context.Context, name string, userID string) bool {
    key := fmt.Sprintf("flag:%s", name)
    
    // Check global
    val, err := f.client.Get(ctx, key).Result()
    if err == nil && val == "true" {
        return true
    }
    
    // Check per-user (canary)
    userKey := fmt.Sprintf("flag:%s:users:%s", name, userID)
    val, err = f.client.Exists(ctx, userKey).Result()
    return err == nil && val > 0
}

func (f *redisFlags) Set(ctx context.Context, name string, enabled bool) error {
    key := fmt.Sprintf("flag:%s", name)
    return f.client.Set(ctx, key, fmt.Sprintf("%t", enabled), 0).Err()
}
```

## 21.3 Usage

```go
type CreatePaymentUseCase struct {
    // ...
    flags featureflag.Flags
}

func (uc *CreatePaymentUseCase) Execute(ctx context.Context, input CreatePaymentInput) (*CreatePaymentOutput, error) {
    // ... validate
    
    // Conditional new rule
    if uc.flags.IsEnabled(ctx, "payment.new_min_amount", input.UserID.String()) {
        if input.Amount.LessThan(decimal.NewFromInt(20)) {
            return nil, domainerrors.ErrAmountTooLow
        }
    }
    
    // ... rest
}
```

## 21.4 Canary Rollout

```bash
# Day 1: 5% users
redis-cli SET flag:payment.new_min_amount true
for i in $(seq 1 50); do
    redis-cli SADD flag:payment.new_min_amount:users "user-$i"
done

# Day 2: 25% users
# ... add more users

# Day 3: 100%
redis-cli SET flag:payment.new_min_amount true
```

## 21.5 Metrics ต่อ Flag

```promql
# Usage
sum(rate(feature_flag_evaluated_total{flag="payment.new_min_amount"}[5m])) by (result)

# Impact
sum(rate(payment_rejected_total{reason="amount_too_low"}[5m]))
```

## 21.6 Checklist

- [ ] Flag naming convention
- [ ] Flag documented
- [ ] Canary strategy
- [ ] Metrics ต่อ flag
- [ ] Auto-rollback trigger
- [ ] Cleanup plan (ลบ flag หลัง stable)
- [ ] Alert on flag change

---

# บทที่ 22: Case Study

## 22.1 Payment v1 → v2 (เปลี่ยน Amount จาก float64 → decimal)

### Timeline รวม 6 สัปดาห์

| Week | Activity | Status |
|---|---|---|
| **1** | RFC + approval, add column, dual write | ✅ |
| **2** | Backfill 8M rows, verify | ✅ |
| **3** | Switch read path, observe | ✅ |
| **4** | Switch write path, observe | ✅ |
| **5** | Rename old column, observe | ✅ |
| **6** | Drop old column | ✅ |

### Metrics

| Metric | Before | After | Improvement |
|---|---|---|---|
| Rounding errors | 500/month | 0 | -100% |
| Audit findings | 3/quarter | 0 | -100% |
| Account accuracy | 99.95% | 100% | +0.05% |
| Response time | 120ms | 118ms | -2% |

### Lessons Learned

1. **RFC ก่อน** — เขียน design doc, ประชุม, approve
2. **Characterization test** — lock behavior ก่อนแก้
3. **Batch backfill** — 8M rows ใน 2 ชั่วโมง
4. **Dual read verify** — เจอ 3 mismatch (0.00004%) → fixed ก่อน switch
5. **Feature flag** — ถ้าเจอปัญหา flip flag ภายใน 30s

### Rollback (ไม่จำเป็นต้องใช้)

- Prepared แต่ไม่ต้องใช้
- ถ้า mismatch > 0.01% → rollback
- Metric: mismatch rate < 0.0001%

## 22.2 Auth JWT → JWT+Session

### Strategy

```
Phase 1: Add session storage (Redis)
Phase 2: Dual auth (JWT + Session)
Phase 3: Log session usage
Phase 4: Switch (Session primary)
Phase 5: Deprecate JWT
Phase 6: Remove JWT
```

### Timeline 3 เดือน

### Rollback

- Flip feature flag → กลับ JWT
- Sessions ยังใช้ได้ (dual auth)
- ไม่มี data loss

---

# บทที่ 23: Checklist และ Anti-Patterns

## 23.1 Pre-Change Checklist

- [ ] Impact analysis เสร็จ
- [ ] RFC เขียน + approved
- [ ] Risk score ประเมินแล้ว
- [ ] Characterization test ครอบ behavior เดิม
- [ ] Rollback plan documented
- [ ] Communication plan ready
- [ ] Feature flag setup (ถ้าจำเป็น)
- [ ] Staging environment ready
- [ ] Monitoring dashboard ready

## 23.2 During-Change Checklist

- [ ] แก้ทีละเล็ก ทีละ commit
- [ ] Test ผ่านทุก commit
- [ ] Code review ผ่าน
- [ ] Deploy staging ก่อน
- [ ] Smoke test ผ่าน
- [ ] Canary deploy (ถ้าจำเป็น)
- [ ] Monitor ระหว่าง deploy

## 23.3 Post-Change Checklist

- [ ] Smoke test ผ่าน
- [ ] Metrics ปกติ (error rate, latency)
- [ ] Monitor 24-48 ชม.
- [ ] Documentation อัปเดต
- [ ] Postmortem (ถ้ามี issue)
- [ ] Team debrief
- [ ] Remove flag (หลัง stable)
- [ ] Clean up deprecated code

## 23.4 Anti-Patterns ที่พบบ่อย

### 1. แก้ทีเดียวเยอะ

**❌ ผิด:**
```
Sprint 1: แก้ 50 ไฟล์, 5000 บรรทัด
```

**✅ ถูก:**
```
Week 1: เพิ่ม column + dual write
Week 2: Backfill data
Week 3: Switch read
Week 4: Switch write
Week 5: Drop old
```

### 2. ไม่มี Rollback Plan

**❌ ผิด:**
```
Deploy → production ล่ม → ไม่รู้จะทำอะไร
```

**✅ ถูก:**
```
Deploy → ล่ม → flip flag (30s)
       → หรือ git revert (2min)
       → หรือ restore backup (1h)
```

### 3. Test บน Local เท่านั้น

**❌ ผิด:**
```
Local: 100 rows → ผ่าน → deploy production (100M rows) → ล่ม
```

**✅ ถูก:**
```
Staging: clone จาก production (anonymized)
       → test migration → วัด latency
       → approve → deploy
```

### 4. แก้ Domain เพื่อให้ Infra ทำงาน

**❌ ผิด:**
```go
type Payment struct {
    Amount *decimal.Decimal  // ← pointer เพราะ DB
}
```

**✅ ถูก:**
```go
type Payment struct {
    Amount decimal.Decimal
}
// Infra จัดการ nullable
```

### 5. Breaking Change โดยไม่ Version

**❌ ผิด:**
```go
// เปลี่ยน field name ตรงๆ
type Response struct {
    CustomerID string `json:"customer_id"`  // ← เก่าคือ user_id
}
```

**✅ ถูก:**
```
/api/v1/payments → ยังคืน user_id
/api/v2/payments → คืน customer_id
Deprecate v1 → 6 เดือน → sunset
```

### 6. ไม่มี Idempotency

**❌ ผิด:**
```go
// Kafka consumer ที่ไม่กัน duplicate
func (c *Consumer) Handle(msg) {
    c.emailService.Send(...)  // ← ส่งซ้ำได้
}
```

**✅ ถูก:**
```go
if c.wasProcessed(msg.EventID) {
    return nil
}
// ... process
c.markProcessed(msg.EventID)
```

### 7. ไม่มี Monitoring

**❌ ผิด:**
```
Deploy → hope for the best
```

**✅ ถูก:**
```
Deploy → dashboard showing:
  - Error rate
  - P95 latency
  - New field null rate
  - Mismatch count
  - DB lock count
```

### 8. Skip Staging

**❌ ผิด:**
```
Local → production
```

**✅ ถูก:**
```
Local → Dev → Staging → Canary → Production
```

### 9. Commit Secret

**❌ ผิด:**
```go
const API_KEY = "sk_live_abc123"  // ← in code
```

**✅ ถูก:**
```go
apiKey := os.Getenv("STRIPE_API_KEY")
```

### 10. ลบ Code เก่าทันที

**❌ ผิด:**
```bash
# Commit 1: deploy code ใหม่
# Commit 2: ลบ code เก่า (ทันที)
```

**✅ ถูก:**
```bash
# Commit 1: deploy code ใหม่
# Wait 2 sprints
# Commit 2: mark deprecated
# Wait 2 sprints
# Commit 3: remove
```

## 23.5 Golden Rules สรุป

```
1. อ่านก่อนแก้
2. Test ก่อน refactor
3. ทีละเล็ก ทีละ commit
4. Additive ก่อน Subtractive
5. Version ถ้า Breaking
6. Feature flag สำหรับ rule ใหม่
7. Staging ที่ data ≈ production
8. Rollback plan เสมอ
9. Monitor หลัง deploy
10. Communicate ล่วงหน้า
```

---

# ภาคผนวก A: Migration Templates

## A.1 Add Column (Safe)

```sql
-- migrations/YYYYMMDD_{module}_add_{column}.sql
-- ============================================================
-- Add {column} to {table}
-- Safe: nullable with default, no rewrite
-- ============================================================

ALTER TABLE {table}
    ADD COLUMN IF NOT EXISTS {column} {type} DEFAULT {default};

COMMENT ON COLUMN {table}.{column} IS 'Added YYYY-MM-DD';
```

## A.2 Rename Column (Multi-Phase)

```sql
-- Phase 1: Add new column
ALTER TABLE {table} ADD COLUMN new_name TYPE;

-- Phase 2: Backfill
UPDATE {table} SET new_name = old_name WHERE new_name IS NULL;

-- Phase 3: Dual write (in code)

-- Phase 4: Switch read (in code)

-- Phase 5: Switch write (in code)

-- Phase 6: Mark old deprecated
COMMENT ON COLUMN {table}.old_name IS 'DEPRECATED: use new_name';

-- Phase 7: Drop after 6 months
ALTER TABLE {table} DROP COLUMN old_name;
```

## A.3 Backfill Batch Template

```sql
DO $$
DECLARE
    batch_size INT := 10000;
    affected INT := 1;
    total INT := 0;
BEGIN
    WHILE affected > 0 LOOP
        WITH batch AS (
            SELECT id, {source}
            FROM {table}
            WHERE {target} IS NULL
            LIMIT batch_size
            FOR UPDATE SKIP LOCKED
        )
        UPDATE {table} t
        SET {target} = batch.{source}::type
        FROM batch
        WHERE t.id = batch.id;
        
        GET DIAGNOSTICS affected = ROW_COUNT;
        total := total + affected;
        
        RAISE NOTICE 'Progress: % (total: %)', affected, total;
        
        COMMIT;
        PERFORM pg_sleep(0.1);
    END LOOP;
END $$;
```

## A.4 Add NOT NULL Safe

```sql
-- Step 1: Add nullable + default
ALTER TABLE {table} ADD COLUMN {col} TYPE DEFAULT {val};

-- Step 2: Backfill (see A.3)

-- Step 3: Add CHECK constraint NOT VALID (fast)
ALTER TABLE {table}
    ADD CONSTRAINT {col}_not_null CHECK ({col} IS NOT NULL) NOT VALID;

-- Step 4: Validate (ไม่ block)
ALTER TABLE {table} VALIDATE CONSTRAINT {col}_not_null;

-- Step 5: Convert to NOT NULL
ALTER TABLE {table}
    ALTER COLUMN {col} SET NOT NULL,
    ALTER COLUMN {col} DROP DEFAULT;
```

## A.5 Index Concurrently

```sql
-- ไม่ block writes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_{table}_{col}
    ON {table}({col});

-- ถ้า fail → drop ที่ค้าง
DROP INDEX CONCURRENTLY IF EXISTS idx_{table}_{col};

-- ตรวจสอบ
SELECT
    indexname,
    indexdef
FROM pg_indexes
WHERE tablename = '{table}' AND indexname = 'idx_{table}_{col}';
```

---

# ภาคผนวก B: Communication Plan

## B.1 Template

```markdown
# Communication Plan: [Change Title]

## Stakeholders
- **Primary**: Developer team
- **Secondary**: QA, DevOps
- **Tertiary**: Product, Support, Clients

## Timeline

| T-relative | Activity | Channel | Owner |
|---|---|---|---|
| T-14d | RFC published | Email + Slack | Tech Lead |
| T-7d | RFC approved | Meeting | Architect |
| T-3d | Deploy to staging | Slack | DevOps |
| T-1d | Final review | Email | Team |
| T-0d | Deploy to production | Slack | DevOps |
| T+1h | Post-deploy check | Slack | On-call |
| T+24h | Status update | Email | Tech Lead |
| T+1w | Postmortem | Meeting | Team |

## Messages

### T-14d: Announcement
Subject: [RFC] Payment v1 → v2 Migration

Body:
- What's changing
- Why
- Timeline
- Impact
- Rollback plan
- Call for feedback

### T-0d: Deployment
Subject: [Deploy] Payment v2 starting

Body:
- Deploy plan
- Monitoring dashboard link
- Rollback trigger
- Contact channel

### T+24h: Status
Subject: [Status] Payment v2 - 24h report

Body:
- Metrics
- Issues (if any)
- Next steps
```

## B.2 Escalation Matrix

| Level | Trigger | Contact | Response |
|---|---|---|---|
| **L1** | Error rate > 0.5% | On-call dev | 15 min |
| **L2** | Error rate > 1% | Team lead | 30 min |
| **L3** | Customer impact | Manager | 1h |
| **L4** | Revenue impact | VP Engineering | 2h |

---

# ภาคผนวก C: Quick Reference Card

## C.1 Rules of Thumb

```
┌─────────────────────────────────────────────────┐
│         MODIFY EXISTING MODULE CHECKLIST        │
├─────────────────────────────────────────────────┤
│ BEFORE                                          │
│ □ อ่านโค้ด 100%                                │
│ □ Characterization test                        │
│ □ Impact analysis                              │
│ □ RFC + approval (ถ้า risk > medium)           │
│ □ Rollback plan                                │
│                                                 │
│ DURING                                          │
│ □ ทีละเล็ก ทีละ commit                         │
│ □ Additive ก่อน Subtractive                    │
│ □ Test ทุก commit                              │
│ □ Staging ก่อน production                      │
│ □ Canary deploy (ถ้าจำเป็น)                    │
│                                                 │
│ AFTER                                           │
│ □ Smoke test                                   │
│ □ Monitor 24-48h                               │
│ □ Communication                                │
│ □ Documentation                                │
│ □ Cleanup deprecated code                      │
└─────────────────────────────────────────────────┘
```

## C.2 Decision Tree

```
การเปลี่ยนแปลงนี้จะกระทบ client ไหม?
│
├─ ไม่ → Breaking change ไม่ต้อง version
│      │
│      ├─ เพิ่มเท่านั้น → Additive (Low risk)
│      └─ แก้ของเดิม → Modify (Medium risk)
│
└─ ใช่ → ต้อง version ใหม่
       │
       ├─ Support 2 versions → 6 เดือน
       └─ Deprecate + migrate
```

## C.3 Migration Phases

```
[1] EXPAND
    ├─ Add new (optional)
    ├─ Keep old working
    └─ Deploy: backward compatible

[2] MIGRATE
    ├─ Backfill data
    ├─ Dual write
    ├─ Dual read + verify
    └─ Deploy: both work

[3] SWITCH
    ├─ Read from new
    ├─ Write to new only
    └─ Deploy: new only

[4] CONTRACT
    ├─ Mark old deprecated
    ├─ Wait (2 sprints)
    ├─ Drop old
    └─ Deploy: clean
```

## C.4 คำสั่งที่ใช้บ่อย

```bash
# หา usage
grep -rn "Payment.UserID" --include="*.go" .
go mod graph | grep payment

# ดู git history
git log -p -5 internal/modules/payment/
git log -L :Refund:internal/modules/payment/domain/entity/payment.go

# DB analysis
psql -c "SELECT pg_size_pretty(pg_total_relation_size('payment_transactions'));"
psql -c "SELECT * FROM pg_stat_statements WHERE query LIKE '%payment%' LIMIT 10;"

# Migration
go run cmd/api/main.go migrate
psql -f migrations/xxx.sql

# Monitor
kubectl logs -f deployment/api -n production
kubectl top pods -n production

# Rollback
kubectl rollout undo deployment/api
git revert <commit>
curl -X POST flag-service/api/flags/xxx -d '{"enabled": false}'
```

## C.5 Error Recovery Time

| Level | ตัวอย่าง | Recovery Time |
|---|---|---|
| **Instant** | Feature flag | < 30s |
| **Fast** | Rollback deploy | 2-5 min |
| **Medium** | Rollback migration | 10-60 min |
| **Slow** | Restore backup | 1-4 h |
| **Impossible** | Drop column | – |

## C.6 5 คำถามที่ต้องตอบก่อน Merge

1. **จะพังอะไรได้บ้าง?** (Failure modes)
2. **ถ้าพังจะรู้ได้ยังไง?** (Detection)
3. **ถ้าพังจะทำอะไร?** (Rollback)
4. **ใช้เวลาเท่าไหร่?** (MTTR)
5. **Data หายไหม?** (Data safety)

---

# 📝 ข้อมูลเอกสาร

**ชื่อเอกสาร:** คู่มือแก้ไข Module เดิมฉบับสมบูรณ์
**เวอร์ชัน:** 1.0
**วันที่:** เมษายน 2026
**จำนวนหน้า:** ~140 หน้า (ประมาณ)
**ระดับ:** Intermediate - Advanced

**ผู้อ่านเป้าหมาย:**
- Senior Go Developer
- Tech Lead ที่ดูแล Production
- DevOps Engineer ที่ manage deployment
- Architect ที่ต้อง approve change

**ข้อกำหนดเบื้องต้น:**
- อ่านเล่ม 1 จบ (หรือเข้าใจ Clean Architecture)
- ประสบการณ์ production deployment 1+ ปี
- เข้าใจ PostgreSQL, migration, transaction
- เข้าใจ Kafka, feature flags

**เอกสารที่เกี่ยวข้อง:**
- เล่ม 1: คู่มือสร้าง Module ใหม่ ✅
- เล่ม 3: คู่มือขาย Module
- เล่ม 4: คู่มือทดสอบและ Deployment
- เล่ม 5: คู่มือบำรุงรักษาและ Scale

**อ้างอิง:**
- Refactoring — Martin Fowler
- Database Refactoring — Scott Ambler
- Continuous Delivery — Humble & Farley
- Site Reliability Engineering — Google
- Release It! — Michael Nygard
- โปรเจกต์ `icmongolang` (case studies)

---

**END OF BOOK 2**

> 📌 **ขั้นถัดไป:** อ่านเล่ม 3 — คู่มือขาย Module
> ที่จะสอนวิธี package module ให้ขายได้, pricing model, licensing,
> และการตลาดสำหรับ Go developers

---

**พิมพ์เมื่อ:** เมษายน 2026
**ผู้จัดทำ:** ทีมสถาปัตยกรรมซอฟต์แวร์ icmongolang