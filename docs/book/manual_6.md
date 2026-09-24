# 🤖 คู่มือ AI Prompt Templates สำหรับงาน Go Module ทุกประเภท
## The Complete AI Prompt Library for Go Module Work

> **เวอร์ชัน 1.0 (เมษายน 2026)**
> คลัง Prompt สำหรับใช้กับ AI (Claude, ChatGPT, Gemini, Cursor, Copilot)
> ครอบคลุมทุกงานใน 5 เล่ม: สร้าง → แก้ไข → ขาย → Deploy → Maintain

---

# สารบัญ

**ภาคที่ 1: หลักการเขียน Prompt ที่ได้ผล**
1. [หลักการ 7 ข้อ](#บทที่-1-หลักการ)
2. [โครงสร้าง Prompt มาตรฐาน](#บทที่-2-โครงสร้าง)

**ภาคที่ 2: Prompt สำหรับสร้าง Module ใหม่ (เล่ม 1)**
3. [Prompt สร้าง Module ครบชุด](#บทที่-3-สร้าง-module)
4. [Prompt สร้าง Domain Layer](#บทที่-4-domain)
5. [Prompt สร้าง Application Layer](#บทที่-5-application)
6. [Prompt สร้าง Infrastructure](#บทที่-6-infrastructure)
7. [Prompt สร้าง Interface](#บทที่-7-interface)

**ภาคที่ 3: Prompt สำหรับแก้ไข Module (เล่ม 2)**
8. [Prompt วิเคราะห์ Impact](#บทที่-8-impact)
9. [Prompt เพิ่ม Field](#บทที่-9-add-field)
10. [Prompt Migration](#บทที่-10-migration)
11. [Prompt Refactor](#บทที่-11-refactor)

**ภาคที่ 4: Prompt สำหรับขาย Module (เล่ม 3)**
12. [Prompt เขียน Documentation](#บทที่-12-docs)
13. [Prompt สร้าง Sales Kit](#บทที่-13-sales)
14. [Prompt Pricing & ROI](#บทที่-14-pricing)

**ภาคที่ 5: Prompt สำหรับทดสอบ & Deploy (เล่ม 4)**
15. [Prompt เขียน Test](#บทที่-15-test)
16. [Prompt CI/CD](#บทที่-16-cicd)
17. [Prompt Deployment](#บทที่-17-deploy)

**ภาคที่ 6: Prompt สำหรับ Maintenance (เล่ม 5)**
18. [Prompt Profiling & Performance](#บทที่-18-profiling)
19. [Prompt Scaling](#บทที่-19-scaling)
20. [Prompt Incident Response](#บทที่-20-incident)

**ภาคผนวก**
- [A. Master Prompt Template](#ภาคผนวก-a)
- [B. Context Injection Pattern](#ภาคผนวก-b)
- [C. Quick Copy-Paste Library](#ภาคผนวก-c)

---

# บทที่ 1: หลักการเขียน Prompt ที่ได้ผล

## 1.1 หลักการ 7 ข้อ

### หลักการ 1: ระบุ Role และ Context

```
❌ แย่: "เขียน Go code ให้หน่อย"

✅ ดี: "คุณเป็น Senior Go Developer ที่เชี่ยวชาญ
Clean Architecture + DDD ในโปรเจกต์ Modular Monolith
ต่อไปนี้คือ context ของโปรเจกต์..."
```

### หลักการ 2: ระบุ Format ของ Output

```
❌ แย่: "อธิบายวิธีทำ"

✅ ดี: "ตอบในรูปแบบ:
1. สรุป 2-3 ประโยค
2. Code ตัวอย่าง (รันได้)
3. คำอธิบายแต่ละส่วน
4. ข้อควรระวัง
5. Unit test"
```

### หลักการ 3: ให้ Examples

```
❌ แย่: "สร้าง handler"

✅ ดี: "สร้าง handler ตาม pattern นี้:
[แนบตัวอย่าง handler ที่มีอยู่]
ต้องการ handler ใหม่สำหรับ X"
```

### หลักการ 4: ระบุ Constraints

```
❌ แย่: "เขียน code"

✅ ดี: "ข้อจำกัด:
- Go 1.22+
- ห้ามใช้ external package นอกจาก:
  github.com/google/uuid
  github.com/shopspring/decimal
- Test coverage ≥ 80%
- Comment 2 ภาษา (ไทย/อังกฤษ)"
```

### หลักการ 5: ให้ Context ของโปรเจกต์

```
❌ แย่: "สร้าง module"

✅ ดี: "โปรเจกต์: icmongolang
- 33 modules ที่ใช้ pattern เดียวกัน
- Structure: internal/modules/{name}/
- Layer: domain, application, infrastructure, interfaces
- Shared pkg: pkg/logger, pkg/validator, pkg/httputil
- Pattern อ้างอิง: internal/modules/payment/
[แนบไฟล์ตัวอย่าง]"
```

### หลักการ 6: ใช้ Chain-of-Thought

```
❌ แย่: "แก้บั๊กนี้"

✅ ดี: "แก้บั๊กนี้ โดยทำตามขั้นตอน:
1. วิเคราะห์ root cause (5 Whys)
2. เสนอ 2-3 solutions พร้อม trade-offs
3. แนะนำ solution ที่ดีที่สุด
4. เขียน test ที่ reproduce bug
5. เขียน fix
6. อธิบายว่าทำไม fix นี้ถูกต้อง"
```

### หลักการ 7: Iterative Refinement

```
Round 1: "สร้าง module payment"
Round 2: "ปรับ Domain Layer ให้รองรับ partial refund"
Round 3: "เพิ่ม test สำหรับ edge case ที่กล่าวถึง"
Round 4: "ทำให้ error messages consistent กับ module อื่น"
```

## 1.2 โครงสร้าง Prompt มาตรฐาน

```
┌─────────────────────────────────────────┐
│  1. ROLE                                │
│     "คุณเป็น..."                        │
├─────────────────────────────────────────┤
│  2. CONTEXT                             │
│     "โปรเจกต์..."                       │
│     "อ้างอิงจาก..."                     │
├─────────────────────────────────────────┤
│  3. TASK                                │
│     "สร้าง/แก้ไข/วิเคราะห์..."          │
├─────────────────────────────────────────┤
│  4. CONSTRAINTS                         │
│     "ข้อจำกัด:..."                      │
├─────────────────────────────────────────┤
│  5. FORMAT                              │
│     "ตอบในรูปแบบ:..."                   │
├─────────────────────────────────────────┤
│  6. EXAMPLES                            │
│     "ตัวอย่าง:..."                      │
├─────────────────────────────────────────┤
│  7. VERIFICATION                        │
│     "ตรวจสอบว่า..."                     │
└─────────────────────────────────────────┘
```

---

# บทที่ 2: โครงสร้าง Prompt มาตรฐาน

## 2.1 Master Template

```markdown
# ROLE
คุณเป็น {expert_role} ที่เชี่ยวชาญ {domain}

# CONTEXT
โปรเจกต์: {project_name}
Stack: {tech_stack}
Pattern: {architecture_pattern}
อ้างอิง: {reference_files_or_modules}

# TASK
{clear_task_description}

# CONSTRAINTS
- {constraint_1}
- {constraint_2}
- ...

# FORMAT
ตอบในรูปแบบ:
1. {section_1}
2. {section_2}
...

# EXAMPLES
ตัวอย่าง input: {example_input}
ตัวอย่าง output: {example_output}

# VERIFICATION
ก่อนส่งคำตอบ ให้ตรวจสอบ:
- [ ] {check_1}
- [ ] {check_2}
```

## 2.2 Context Injection Pattern

**สำหรับใช้กับ AI ที่ไม่จำ context:**

```
[CONTEXT SNAPSHOT]
┌─────────────────────────────────────────┐
│ Project: icmongolang                    │
│ Go: 1.22                                │
│ Pattern: Clean Architecture + DDD       │
│ Structure:                              │
│   internal/modules/{name}/              │
│     ├── domain/                         │
│     ├── application/                    │
│     ├── infrastructure/                 │
│     └── interfaces/                     │
│ Shared packages:                        │
│   pkg/logger, pkg/validator, ...        │
│ Reference module: internal/modules/payment/ │
└─────────────────────────────────────────┘
[END CONTEXT]
```

---

# บทที่ 3: Prompt สร้าง Module ครบชุด

## 3.1 Master Prompt: สร้าง Module ใหม่ทั้งหมด

```markdown
# ROLE
คุณเป็น Senior Go Developer ที่เชี่ยวชาญ Clean Architecture + DDD
และ Modular Monolith Pattern

# CONTEXT
โปรเจกต์: icmongolang
Go version: 1.22+
Architecture: Clean Architecture + DDD
Structure: Modular Monolith (33 modules)

โครงสร้าง Module มาตรฐาน:
internal/modules/{module_name}/
├── domain/
│   ├── entity/{entity}.go
│   ├── value_object/{vo}.go
│   ├── repository/{entity}_repository.go
│   └── errors/errors.go
├── application/
│   ├── {verb}_{entity}.go
│   └── mocks/{entity}_repository_mock.go
├── infrastructure/
│   └── persistence/postgres/
│       ├── models.go
│       └── {entity}_repo_impl.go
├── interfaces/
│   └── http/
│       ├── {entity}_handler.go
│       ├── dto.go
│       └── routes.go
└── module.go

Shared packages ที่ใช้ได้:
- icmongolang/pkg/logger (structured logging)
- icmongolang/pkg/validator (input validation)
- icmongolang/pkg/httputil (JSON responses)
- github.com/google/uuid
- github.com/shopspring/decimal
- gorm.io/gorm
- github.com/go-chi/chi/v5

# TASK
สร้าง Module ใหม่ชื่อ "{module_name}" ที่มี:

## Business Requirements
{business_requirements}

## Entities
{entity_list}

## Use Cases
{use_case_list}

## API Endpoints
{api_endpoints}

# CONSTRAINTS
1. Domain Layer ห้าม import: gorm, gin, chi, sarama, redis
2. Entity ต้องมี constructor `New{Entity}` + validation
3. เปลี่ยน state ผ่าน behavior methods เท่านั้น (ไม่มี setter)
4. Repository เป็น interface เท่านั้น
5. Errors เป็น sentinel errors
6. Comment 2 ภาษา (ไทย/อังกฤษ)
7. Test coverage ≥ 80%
8. Table prefix: {module_name}_
9. Redis key: {module_name}:{entity}:{id}

# FORMAT
ตอบในรูปแบบ:

## 1. Architecture Diagram
[แสดงโครงสร้าง module]

## 2. File Tree
```
internal/modules/{module_name}/
├── ... (ระบุทุกไฟล์)
```

## 3. Code ทั้งหมด
### 3.1 Domain Layer
- entity/{entity}.go (พร้อม test)
- value_object/{vo}.go (พร้อม test)
- repository/{entity}_repository.go
- errors/errors.go

### 3.2 Application Layer
- application/{verb}_{entity}.go (ทุก use case)
- application/mocks/{entity}_repository_mock.go

### 3.3 Infrastructure Layer
- infrastructure/persistence/postgres/models.go
- infrastructure/persistence/postgres/{entity}_repo_impl.go

### 3.4 Interface Layer
- interfaces/http/{entity}_handler.go
- interfaces/http/dto.go
- interfaces/http/routes.go

### 3.5 Module Bootstrap
- module.go
- Wire-up ใน main.go

## 4. Migration SQL
migrations/YYYYMMDD_{module_name}_init.sql

## 5. Tests
- Unit tests ทุก layer
- Integration test
- Handler test

## 6. Documentation
- README.md ของ module
- API examples

# VERIFICATION
ก่อนส่งคำตอบ ให้ตรวจสอบ:
- [ ] Domain ไม่มี GORM tag
- [ ] Entity มี constructor + behavior methods
- [ ] Repository interface อยู่ใน domain
- [ ] Use case ละ 1 ไฟล์
- [ ] Handler ดึง user_id จาก context
- [ ] ทุก error เป็น sentinel
- [ ] Migration มี prefix
- [ ] Test ครอบ happy + error paths
- [ ] Comment ครบ 2 ภาษา
```

## 3.2 ตัวอย่างการใช้งานจริง

```markdown
# ตัวอย่าง: สร้าง Module "Inventory"

# ROLE
คุณเป็น Senior Go Developer ที่เชี่ยวชาญ Clean Architecture + DDD

# CONTEXT
[ตาม Master Prompt ด้านบน]

# TASK
สร้าง Module "inventory" สำหรับจัดการสต็อกสินค้า

## Business Requirements
1. เพิ่ม/ลบ/แก้ไขสินค้าในสต็อก
2. ปรับสต็อกเมื่อมีการขาย
3. แจ้งเตือนเมื่อสต็อกต่ำ
4. รายงานสต็อกตามหมวดหมู่

## Entities
- Product (aggregate root)
  - id, sku, name, price, quantity, min_quantity
  - status: ACTIVE, INACTIVE, OUT_OF_STOCK
- StockMovement (entity ย่อย)
  - id, product_id, quantity, type, reason, created_at

## Use Cases
- CreateProduct
- UpdateProduct
- AdjustStock (increase/decrease)
- GetProduct
- ListProducts
- GetLowStockProducts

## API Endpoints
- POST   /api/v1/inventory/products
- GET    /api/v1/inventory/products/{id}
- GET    /api/v1/inventory/products
- PUT    /api/v1/inventory/products/{id}
- POST   /api/v1/inventory/products/{id}/adjust
- GET    /api/v1/inventory/products/low-stock

# [ให้ AI สร้างตาม Format ที่กำหนด]
```

---

# บทที่ 4: Prompt สร้าง Domain Layer

## 4.1 Prompt สร้าง Entity

```markdown
# ROLE
คุณเป็น DDD Expert ที่เชี่ยวชาญการออกแบบ Aggregate

# CONTEXT
โปรเจกต์: icmongolang
Module: {module_name}
Pattern: Domain-Driven Design
Reference: internal/modules/payment/domain/entity/payment.go

# TASK
สร้าง Entity ชื่อ "{EntityName}" ที่มี:

## Fields
{field_list}

## Invariants (กฎที่ต้องเป็นจริงเสมอ)
{invariants}

## State Transitions
{state_transitions}

## Behavior Methods
{behavior_methods}

# CONSTRAINTS
1. ห้าม import GORM, HTTP, framework
2. ใช้ `github.com/google/uuid` สำหรับ ID
3. ใช้ `github.com/shopspring/decimal` สำหรับเงิน
4. Constructor `New{Entity}` validate ทุก invariant
5. เปลี่ยน state ผ่าน behavior method เท่านั้น
6. Sentinel errors สำหรับทุก error case
7. Comment 2 ภาษา

# FORMAT
## 1. Entity Struct
```go
type {EntityName} struct {
    // ...
}
```

## 2. Constructor
```go
func New{EntityName}(...) (*{EntityName}, error) {
    // validation + invariants
}
```

## 3. Behavior Methods
[ทุก method พร้อม comment 2 ภาษา]

## 4. Query Methods
[IsValid, CanX, etc.]

## 5. Unit Tests
```go
func Test{EntityName}_Success(t *testing.T) { ... }
func Test{EntityName}_ValidationErrors(t *testing.T) { ... }
func Test{EntityName}_StateTransitions(t *testing.T) { ... }
```

# EXAMPLE INPUT
สร้าง Entity "Subscription":
- Fields: id, user_id, plan (FREE/PRO/ENTERPRISE), status (ACTIVE/CANCELLED/EXPIRED), start_date, end_date
- Invariants: end_date > start_date, plan ต้อง valid
- Transitions: ACTIVE → CANCELLED, ACTIVE → EXPIRED
- Methods: Activate, Cancel, Expire, IsActive, DaysRemaining
```

## 4.2 Prompt สร้าง Value Object

```markdown
# ROLE
คุณเป็น DDD Expert ที่เชี่ยวชาญ Value Objects

# CONTEXT
Module: {module_name}
Reference: internal/modules/payment/domain/value_object/status.go

# TASK
สร้าง Value Object "{VO_Name}" ที่มี:

## Possible Values
{values}

## Validation Rules
{validation_rules}

## Behavior Methods
{methods}

# CONSTRAINTS
1. เป็น immutable
2. Type: string หรือ custom type
3. มี IsValid() method
4. มี String() method
5. Comment 2 ภาษา

# FORMAT
## 1. Type + Constants
## 2. Methods
## 3. Unit Tests (100% coverage)

# EXAMPLE
สร้าง Value Object "SubscriptionPlan":
- Values: FREE, PRO, ENTERPRISE
- Methods: IsValid, String, Price, Features, CanUpgradeTo
```

---

# บทที่ 5: Prompt สร้าง Application Layer

## 5.1 Prompt สร้าง Use Case

```markdown
# ROLE
คุณเป็น Senior Go Developer ที่เชี่ยวชาญ Use Case Pattern

# CONTEXT
Module: {module_name}
Layer: Application
Reference: internal/modules/payment/application/create_payment.go

# TASK
สร้าง Use Case ชื่อ "{Verb}{Entity}UseCase"

## Input
{input_fields}

## Output
{output_fields}

## Business Logic
{business_logic}

## Side Effects
- [ ] Publish event
- [ ] Send notification
- [ ] Update cache
- [ ] Audit log

## Dependencies
- Repository: {repo_interface}
- Services: {service_list}
- Logger: pkg/logger

# CONSTRAINTS
1. 1 use case = 1 ไฟล์
2. Execute() เป็น method เดียว
3. Return error ไม่ panic
4. ไม่มี SQL, HTTP, framework
5. ใช้ logger สำหรับ structured logging
6. Comment 2 ภาษา

# FORMAT
## 1. Struct + Constructor
## 2. Input/Output DTO
## 3. Execute() Method
## 4. Unit Tests (with mock)
## 5. Integration Test (ถ้าจำเป็น)

# EXAMPLE
สร้าง Use Case "CancelSubscriptionUseCase":
- Input: subscription_id, user_id, reason
- Business Logic: 
  1. Load subscription
  2. Check ownership
  3. Check can cancel
  4. Cancel
  5. Refund (if within 7 days)
  6. Publish event
- Output: refund_amount, cancel_date
```

## 5.2 Prompt สร้าง Mock

```markdown
# ROLE
คุณเป็น Go Developer ที่เชี่ยวชาญ Testing

# TASK
สร้าง Mock สำหรับ Repository interface "{RepositoryName}"

# CONSTRAINTS
1. ใช้ github.com/stretchr/testify/mock
2. ทุก method ต้อง mock ได้
3. Handle nil return values
4. Comment อธิบาย

# FORMAT
```go
// File: internal/modules/{module}/application/mocks/{repo}_mock.go
package mocks

type {RepositoryName}Mock struct {
    mock.Mock
}

// ... all methods
```

# EXAMPLE
สร้าง Mock สำหรับ "SubscriptionRepository":
- Methods: Save, FindByID, FindByUserID, Update, Delete
```

---

# บทที่ 6: Prompt สร้าง Infrastructure

## 6.1 Prompt สร้าง GORM Model

```markdown
# ROLE
คุณเป็น Go Developer ที่เชี่ยวชาญ GORM

# CONTEXT
Module: {module_name}
Table prefix: {module_name}_

# TASK
สร้าง GORM Model สำหรับ Entity "{EntityName}"

## Fields
{fields_with_types}

## Relations
{relations}

## Indexes
{indexes}

# CONSTRAINTS
1. Type: uuid สำหรับ PK
2. Type: numeric(x,y) สำหรับเงิน
3. Type: jsonb สำหรับ metadata
4. มี TableName() method คืน prefix
5. Index สำหรับ FK, status, created_at
6. Comment 2 ภาษา

# FORMAT
## 1. Model Struct
## 2. TableName Method
## 3. Conversion (toEntity, toModel)
## 4. Integration Test

# EXAMPLE
สร้าง Model สำหรับ Subscription:
- id: UUID PK
- user_id: UUID (index)
- plan: varchar(20)
- status: varchar(20) (index)
- start_date, end_date: timestamp
- metadata: jsonb
- Table: subscriptions_subscriptions
```

## 6.2 Prompt สร้าง Repository Implementation

```markdown
# ROLE
คุณเป็น Go Developer ที่เชี่ยวชาญ GORM + Repository Pattern

# CONTEXT
Module: {module_name}
Interface: internal/modules/{module}/domain/repository/{entity}_repository.go
Reference: internal/modules/payment/infrastructure/persistence/postgres/payment_repo_impl.go

# TASK
สร้าง Repository Implementation ที่ implement interface

# METHODS TO IMPLEMENT
{method_list}

# CONSTRAINTS
1. ใช้ GORM กับ context ทุกครั้ง: db.WithContext(ctx)
2. Map errors เป็น domain errors
3. ไม่มี business logic
4. Transaction เมื่อจำเป็น
5. Comment 2 ภาษา

# FORMAT
## 1. Struct + Constructor
## 2. Methods (ตาม interface)
## 3. Conversion functions (toModel, toEntity)
## 4. Integration Tests (with testcontainers)

# EXAMPLE
Implement SubscriptionRepository:
```go
type subscriptionRepoImpl struct {
    db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *subscriptionRepoImpl { ... }

func (r *subscriptionRepoImpl) Save(ctx context.Context, s *entity.Subscription) error { ... }
// etc.
```
```

---

# บทที่ 7: Prompt สร้าง Interface

## 7.1 Prompt สร้าง HTTP Handler

```markdown
# ROLE
คุณเป็น Go Developer ที่เชี่ยวชาญ REST API + Chi Router

# CONTEXT
Module: {module_name}
Layer: Interfaces/HTTP
Reference: internal/modules/payment/interfaces/http/payment_handler.go

# TASK
สร้าง HTTP Handler สำหรับ Entity "{EntityName}"

## Endpoints
{endpoints}

## Use Cases
{use_cases}

# CONSTRAINTS
1. ใช้ chi router
2. ดึง user_id จาก context เท่านั้น
3. Validate request ด้วย pkg/validator
4. Response ใช้ pkg/httputil.JSON
5. Error mapping ตาม convention
6. ไม่มี business logic
7. Comment 2 ภาษา

# FORMAT
## 1. Handler Struct + Constructor
## 2. Request DTO (with validation tags)
## 3. Response DTO
## 4. Handler Methods (ทุก endpoint)
## 5. Error Mapping
## 6. Handler Tests

# EXAMPLE
สร้าง Subscription Handler:
- POST /api/v1/subscriptions
- GET /api/v1/subscriptions/{id}
- PUT /api/v1/subscriptions/{id}
- POST /api/v1/subscriptions/{id}/cancel
```

## 7.2 Prompt สร้าง Routes

```markdown
# ROLE
คุณเป็น Go Developer

# TASK
สร้าง Routes registration สำหรับ Module "{module_name}"

# ENDPOINTS
{endpoints}

# CONSTRAINTS
1. ทุก route ต้องผ่าน authMiddleware
2. Path prefix: /api/v1/{module_name}/
3. ใช้ chi.Router

# FORMAT
```go
// File: interfaces/http/routes.go
package http

type Handlers struct {
    {Entity}: *{Entity}Handler
}

func RegisterRoutes(r chi.Router, h *Handlers, authMW func(http.Handler) http.Handler) {
    // ...
}
```
```

---

# บทที่ 8: Prompt วิเคราะห์ Impact

## 8.1 Master Prompt: Impact Analysis

```markdown
# ROLE
คุณเป็น Senior Architect ที่เชี่ยวชาญ Risk Assessment

# CONTEXT
โปรเจกต์: icmongolang
Module ที่จะแก้: {module_name}
Change ที่ต้องการ: {change_description}

# TASK
วิเคราะห์ผลกระทบของการเปลี่ยนแปลงนี้ อย่างละเอียด

# INPUT
- Change: {change}
- Reason: {reason}
- Urgency: {urgency}

# ANALYSIS REQUIRED

## 1. Scope Analysis
- ไฟล์ที่จะต้องแก้ (ระบุทั้งหมด)
- Layer ที่กระทบ
- Module ที่ depend
- External clients

## 2. Risk Assessment
- Risk level (Low/Medium/High/Critical)
- Failure modes
- Data loss potential
- Downtime potential
- Rollback difficulty

## 3. Impact Matrix
| Area | Impact | Mitigation |
|------|--------|------------|
| Code | | |
| DB | | |
| API | | |
| Kafka | | |
| Cache | | |
| Monitoring | | |
| Cost | | |

## 4. Dependencies
- Upstream: อะไร depend on this?
- Downstream: this depend on อะไร?
- Circular dependencies?

## 5. Timeline Estimate
- Development: X days
- Testing: Y days
- Deployment: Z days
- Total: X+Y+Z

## 6. Rollback Plan
- Method: (git revert / flag / restore)
- Time to rollback: X minutes
- Data safety: (safe / risky / loss)

## 7. Recommendation
- GO / NO-GO
- Conditions
- Alternative approaches

# FORMAT
ตอบเป็น report ที่มีหัวข้อตาม ANALYSIS REQUIRED

# EXAMPLE INPUT
Change: เปลี่ยน Payment.Amount จาก float64 → decimal.Decimal
Reason: Precision issues in accounting
Urgency: High (compliance deadline)

# ให้ AI วิเคราะห์
```

---

# บทที่ 9: Prompt เพิ่ม Field

## 9.1 Prompt เพิ่ม Field แบบปลอดภัย

```markdown
# ROLE
คุณเป็น Senior Go Developer ที่เชี่ยวชาญ Database Migration

# CONTEXT
Module: {module_name}
Entity: {entity_name}
Current fields: {current_fields}

# TASK
เพิ่ม Field ใหม่ "{field_name}" ใน Entity

## Field Spec
- Name: {field_name}
- Type: {type}
- Required: {yes/no}
- Default: {default_value}
- Nullable: {yes/no}

# CONSTRAINTS
1. Backward compatible (existing rows ต้องใช้ได้)
2. Migration ต้องไม่ lock table
3. Rollback ได้ง่าย
4. Zero downtime

# FORMAT
## 1. Domain Entity Changes
```go
// เพิ่ม field + comment
```

## 2. GORM Model Changes
```go
// เพิ่ม column + tag
```

## 3. Migration SQL
```sql
-- Safe migration (nullable + default)
```

## 4. Repository Changes
```go
// toModel, toEntity
```

## 5. DTO Changes
```go
// Request, Response
```

## 6. Handler Changes
```go
// Parse + validate
```

## 7. Tests
- Unit test สำหรับ field ใหม่
- Integration test

## 8. Rollback Plan
- ไม่ต้อง rollback (additive)
- หรือ: DROP COLUMN if needed

# EXAMPLE
เพิ่ม "notes" (optional string, max 500) ใน Payment entity
```

## 9.2 Prompt เพิ่ม Field ที่เป็น Enum

```markdown
# TASK
เพิ่ม Enum field "{field_name}" ใน Entity "{EntityName}"

## Possible Values
{values}

## Default
{default}

## Migration Strategy
- Add column (nullable)
- Backfill existing rows with default
- Add NOT NULL constraint (if needed)

# FORMAT
## 1. Value Object
## 2. Entity field
## 3. GORM column
## 4. Migration (3 steps)
## 5. Backfill script
## 6. Validation
## 7. Tests
```

---

# บทที่ 10: Prompt Migration

## 10.1 Prompt สร้าง Migration

```markdown
# ROLE
คุณเป็น Database Expert ที่เชี่ยวชาญ PostgreSQL + Zero-Downtime

# CONTEXT
Module: {module_name}
Database: PostgreSQL 15+
Current schema: {current_schema}

# TASK
สร้าง Migration สำหรับการเปลี่ยนแปลงนี้:
{change_description}

# CONSTRAINTS
1. Zero-downtime (ไม่ lock table)
2. Idempotent (รันซ้ำได้)
3. Rollback script ครบ
4. Index ที่จำเป็น
5. Constraints ที่จำเป็น
6. Comment อธิบาย

# FORMAT
## 1. Migration File (up)
```sql
-- migrations/YYYYMMDD_{module}_{description}.sql
```

## 2. Rollback File (down)
```sql
-- migrations/YYYYMMDD_{module}_{description}.down.sql
```

## 3. Verification Queries
```sql
-- ตรวจสอบว่าสำเร็จ
```

## 4. Rollback Test Plan
- Test on staging
- Expected outcome

## 5. Estimated Time
- Small table (< 100K rows): X seconds
- Medium (100K-10M): Y minutes
- Large (> 10M): Z hours

# EXAMPLE
Migration: เปลี่ยน payment.user_id จาก UUID → BIGINT
- Table size: 10M rows
- Requires: Expand-Contract pattern
```

## 10.2 Prompt Migration แบบ Data Backfill

```markdown
# ROLE
คุณเป็น Database Expert

# TASK
เขียน Script backfill ข้อมูลสำหรับ column ใหม่

## Context
- Table: {table_name}
- Rows: {row_count}
- New column: {new_column}
- Old column: {old_column}
- Transformation: {logic}

# CONSTRAINTS
1. Batch processing (10K rows/batch)
2. Sleep ระหว่าง batch (0.1s)
3. ใช้ FOR UPDATE SKIP LOCKED
4. Monitorable (log progress)
5. Resumable (save last_id)
6. Rollback-able

# FORMAT
## 1. SQL Script (DO block)
## 2. Go Script (สำหรับ complex logic)
## 3. Monitoring Queries
## 4. Verification
## 5. Rollback
```

---

# บทที่ 11: Prompt Refactor

## 11.1 Prompt Refactor Code

```markdown
# ROLE
คุณเป็น Refactoring Expert ที่เชี่ยวชาญ Go

# CONTEXT
ไฟล์ที่จะ refactor: {file_path}
ปัญหาปัจจุบัน: {problem}
เป้าหมาย: {goal}

# CURRENT CODE
```go
{paste_current_code}
```

# TASK
Refactor โค้ดนี้ตามเป้าหมาย

## Requirements
1. Behavior ไม่เปลี่ยน (test ต้องผ่าน)
2. ลด complexity
3. เพิ่ม readability
4. ลด duplication
5. ตาม Go idioms

# CONSTRAINTS
1. ห้าม breaking API (public methods)
2. Test ทุกอย่างต้องผ่าน
3. Comment 2 ภาษา
4. ไม่เพิ่ม dependency

# FORMAT
## 1. Analysis (ปัญหาปัจจุบัน)
## 2. Refactoring Plan
## 3. Refactored Code
## 4. Explanation (ทำไมดีกว่า)
## 5. Tests (ยืนยันว่า behavior เดิม)
## 6. Benchmark (ถ้า performance concern)

# EXAMPLE INPUT
Refactor ไฟล์: internal/modules/payment/application/create_payment.go
ปัญหา: Execute() ยาว 200 บรรทัด, ทำหลายอย่าง
เป้าหมาย: แยกเป็น helper methods, แต่ละ method < 30 บรรทัด
```

## 11.2 Prompt Strangler Fig Refactor

```markdown
# ROLE
คุณเป็น Architect ที่เชี่ยวชาญ Legacy Migration

# CONTEXT
Legacy code: {legacy_module}
Target: {target_architecture}
Reason: {reason}

# TASK
ออกแบบ Strangler Fig pattern เพื่อ migrate แบบปลอดภัย

# REQUIRED
## 1. Mapping
- Legacy functions → New use cases
- Data flow changes

## 2. Phase Plan
- Phase 1: Add new (parallel)
- Phase 2: Dual write
- Phase 3: Switch read
- Phase 4: Switch write
- Phase 5: Remove old

## 3. Feature Flags
- Flag names
- Rollout strategy (5% → 25% → 100%)
- Rollback trigger

## 4. Testing Strategy
- Characterization tests ก่อน
- Parallel run verification
- Success metrics

## 5. Timeline
- Per phase duration
- Total timeline

## 6. Risk & Mitigation

# FORMAT
[Report ตาม REQUIRED]
```

---

# บทที่ 12: Prompt เขียน Documentation

## 12.1 Prompt README ที่ขายได้

```markdown
# ROLE
คุณเป็น Technical Writer + Developer ที่เชี่ยวชาญ Developer Experience

# CONTEXT
Module: {module_name}
Target audience: {audience}
Goal: {goal}

# INPUT
- Features: {features}
- Performance: {benchmarks}
- Pricing: {pricing}
- Unique value: {usp}

# TASK
เขียน README.md ที่:
1. ทำให้ developer อยากใช้ภายใน 10 วินาที
2. บอก value ชัดเจน
3. มี quick start ที่รันได้
4. แสดง social proof
5. CTA ที่ชัดเจน

# FORMAT
## 1. Header (Title + Tagline)
## 2. Badges
## 3. Features (top 5)
## 4. Quick Start (5 min)
## 5. Performance Table
## 6. Use Cases
## 7. Pricing (ถ้ามี)
## 8. Documentation Links
## 9. Support Channels
## 10. License
## 11. Star History

# EXAMPLE
Module: payment-module
Target: Go developers building e-commerce
Goal: Get 100 GitHub stars in first month
Features: Multi-provider, Refund, Event-driven
Performance: 5000 TPS, 45ms P95
```

## 12.2 Prompt Getting Started Guide

```markdown
# ROLE
คุณเป็น Technical Writer

# TASK
เขียน Getting Started Guide ที่ใช้เวลา 5 นาที

# INPUT
- Module: {module_name}
- Prerequisites: {prerequisites}
- First use case: {first_use_case}

# FORMAT
## Prerequisites (bullet list)
## Step 1: Install (1 command)
## Step 2: Get credentials (1 step)
## Step 3: First request (10 lines of code)
## Step 4: Verify (1 command + expected output)
## Step 5: Next steps (3 links)
## Troubleshooting (top 3 errors)

# CONSTRAINTS
1. Code ทุกอย่างต้องรันได้จริง
2. Output ทุกอย่างต้องชัดเจน
3. Error messages ต้องตรงกับที่ user จะเจอ
4. Screenshots (ถ้าจำเป็น)
5. ใช้ภาษาเข้าใจง่าย
```

## 12.3 Prompt Architecture Decision Record

```markdown
# ROLE
คุณเป็น Senior Architect

# TASK
เขียน ADR (Architecture Decision Record) สำหรับ decision นี้:
{decision_description}

# INPUT
- Context: {context}
- Decision: {decision}
- Alternatives: {alternatives}
- Consequences: {consequences}

# FORMAT
# ADR-{NNN}: {Title}

## Status
[Proposed | Accepted | Deprecated | Superseded]

## Date
YYYY-MM-DD

## Context
[ทำไมต้องตัดสินใจ - 3-5 ประโยค]

## Decision
[ตัดสินใจอะไร - ชัดเจน]

## Consequences
### Positive
- [ข้อดี]

### Negative
- [ข้อเสีย]

### Neutral
- [ผลกระทบ]

## Alternatives Considered
### Alternative 1: {name}
- Pros: [list]
- Cons: [list]
- Rejected because: [reason]

## References
- [Link 1]
- [Link 2]
```

---

# บทที่ 13: Prompt สร้าง Sales Kit

## 13.1 Prompt Product One-Pager

```markdown
# ROLE
คุณเป็น Product Marketing Manager ที่เชี่ยวชาญ Developer Tools

# CONTEXT
Product: {product_name}
Target: {target_market}
Price: {pricing}

# TASK
เขียน Product One-Pager (1 หน้า A4)

# INPUT
- Problem solved: {problem}
- Solution: {solution}
- Key features: {features}
- Pricing: {pricing}
- Social proof: {proof}

# FORMAT
┌──────────────────────────────────────┐
│ [LOGO] {Product Name}                │
│ {Tagline - 1 บรรทัด}                 │
├──────────────────────────────────────┤
│                                      │
│ ✅ {Benefit 1}                       │
│ ✅ {Benefit 2}                       │
│ ✅ {Benefit 3}                       │
│ ✅ {Benefit 4}                       │
│                                      │
├──────────────────────────────────────┤
│ 📊 PERFORMANCE                       │
│ • {Metric 1}                         │
│ • {Metric 2}                         │
├──────────────────────────────────────┤
│ 💰 PRICING                           │
│ • {Tier 1}: {price}                  │
│ • {Tier 2}: {price}                  │
├──────────────────────────────────────┤
│ 📞 CONTACT                           │
│ {email} | {website}                  │
└──────────────────────────────────────┘

# CONSTRAINTS
1. Tagline: max 60 characters
2. Benefits: action-oriented
3. Numbers > adjectives
4. CTA clear
```

## 13.2 Prompt Case Study

```markdown
# ROLE
คุณเป็น Content Marketer ที่เชี่ยวชาญ B2B Tech

# TASK
เขียน Case Study สำหรับลูกค้า

# INPUT
- Customer: {customer_name}
- Industry: {industry}
- Challenge: {challenge}
- Solution: {solution}
- Results: {results}

# FORMAT
# How {Customer} {Achieved Result} with {Your Product}

## Executive Summary
[2-3 ประโยค]

## The Challenge
[ปัญหา + context - 200 words]

## The Solution
[ใช้ product ยังไง - 200 words]

## The Results
| Metric | Before | After | Change |
|--------|--------|-------|--------|
| | | | |

## Quote
> "{testimonial}"
> — {Name}, {Title} at {Company}

## Technical Details
[Implementation specifics - 150 words]

## Try It Yourself
[CTA]

# CONSTRAINTS
1. Specific numbers
2. Real quotes (หรือ realistic)
3. Focus on customer (not product)
4. Length: 800-1200 words
```

## 13.3 Prompt Comparison Chart

```markdown
# ROLE
คุณเป็น Competitive Analyst

# TASK
สร้าง Comparison Chart กับ competitors

# INPUT
- Your product: {yours}
- Competitors: {competitors}
- Features to compare: {features}

# FORMAT
| Feature | You | {Comp 1} | {Comp 2} | {Comp 3} |
|---------|-----|----------|----------|----------|
| {Feature 1} | ✅ | ✅ | ❌ | ✅ |
| {Feature 2} | ✅ | ❌ | ❌ | ✅ |
| Performance | {value} | {value} | {value} | {value} |
| Pricing | {price} | {price} | {price} | {price} |

## Analysis
- **Best for {use case 1}**: {recommendation}
- **Best for {use case 2}**: {recommendation}

# CONSTRAINTS
1. Fair comparison (ไม่ bias)
2. Numbers ตรวจสอบได้
3. Cite sources
4. Focus on value, not FUD
```

---

# บทที่ 14: Prompt Pricing & ROI

## 14.1 Prompt Pricing Strategy

```markdown
# ROLE
คุณเป็น Pricing Strategist ที่เชี่ยวชาญ SaaS/Developer Tools

# CONTEXT
Product: {product_name}
Target: {target_market}
Competitors: {competitors}
Costs: {costs}

# TASK
ออกแบบ Pricing Strategy

## Required Analysis
1. Value-based pricing (ไม่ใช่ cost-plus)
2. Market positioning
3. Tier design
4. Discount policy
5. Grandfathering policy

# INPUT
- Value delivered: {value}
- Customer savings: {savings}
- Competitor prices: {competitor_prices}
- Our costs: {costs}

# FORMAT

## 1. Market Analysis
- Competitor pricing
- Price sensitivity
- Willingness to pay

## 2. Value-Based Pricing
- Customer ROI
- Our price as % of value
- Justification

## 3. Tier Design
| Tier | Price | Features | Target |
|------|-------|----------|--------|
| Free | $0 | | |
| Pro | $X | | |
| Business | $Y | | |
| Enterprise | Custom | | |

## 4. Pricing Psychology
- Anchor pricing
- Decoy effect
- Charm pricing
- Annual discount

## 5. Discount Policy
- Annual: X%
- Startup: X%
- Non-profit: X%
- Volume: X%

## 6. Grandfathering
- Policy
- Communication
- Timeline

## 7. Projection
- Year 1: $X
- Year 2: $Y
- Year 3: $Z
```

## 14.2 Prompt ROI Calculator

```markdown
# ROLE
คุณเป็น Business Analyst

# TASK
สร้าง ROI Calculator สำหรับลูกค้า

# INPUT
- Customer pain: {pain}
- Current cost: {current_cost}
- Our product: {product}
- Our price: {price}

# FORMAT
# ROI Analysis: {Product} vs Current State

## Current State (Without Us)
| Item | Monthly Cost |
|------|--------------|
| {Cost 1} | $X |
| {Cost 2} | $Y |
| **Total** | **$X+Y** |

## With {Product}
| Item | Monthly Cost |
|------|--------------|
| License | ${price} |
| Integration (one-time) | $Z |
| **Total** | **${price}** |

## Savings
- **Monthly**: $X+Y - price
- **Annual**: (X+Y-price) × 12
- **Payback period**: {days}

## 3-Year TCO
| Year | Cost | Savings | Net |
|------|------|---------|-----|
| 1 | | | |
| 2 | | | |
| 3 | | | |
| **Total** | | | |

## ROI
- Year 1: X%
- 3-Year: Y%

## Assumptions
- [Assumption 1]
- [Assumption 2]

# CONSTRAINTS
1. Conservative estimates
2. Cite assumptions
3. Comparable numbers
4. Verifiable
```

---

# บทที่ 15: Prompt เขียน Test

## 15.1 Prompt Unit Test (Domain Layer)

```markdown
# ROLE
คุณเป็น Go Test Expert ที่เชี่ยวชาญ table-driven tests

# CONTEXT
Module: {module_name}
File: domain/entity/{entity}.go
Function to test: {function_name}

# CURRENT CODE
```go
{paste_code}
```

# TASK
เขียน unit tests ที่:

## Coverage Required
- [ ] Happy path
- [ ] All validation errors
- [ ] All state transitions (valid + invalid)
- [ ] Edge cases (zero, nil, max, negative)
- [ ] Boundary conditions

# CONSTRAINTS
1. Coverage 100% สำหรับ domain
2. ใช้ testify (assert, require)
3. Table-driven tests
4. Deterministic (ไม่ flaky)
5. Comment 2 ภาษา

# FORMAT
```go
package {pkg}_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNew{Entity}_Success(t *testing.T) { ... }

func TestNew{Entity}_ValidationErrors(t *testing.T) {
    tests := []struct{
        name    string
        // input
        wantErr error
    }{
        // cases
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test
        })
    }
}

func Test{Entity}_{Behavior}(t *testing.T) { ... }
```

# EXAMPLE
Test file: internal/modules/payment/domain/entity/payment_test.go
Functions: NewPayment, StartProcessing, MarkSuccess, Refund
```

## 15.2 Prompt Use Case Test (with Mock)

```markdown
# ROLE
คุณเป็น Go Test Expert

# CONTEXT
Use Case: {usecase}
Dependencies: {repos}, {services}

# TASK
เขียน unit test สำหรับ use case ด้วย mock

# SCENARIOS TO COVER
1. Success (happy path)
2. Validation errors
3. Repository errors (not found, duplicate, DB error)
4. Authorization errors
5. Business rule violations
6. Edge cases

# CONSTRAINTS
1. Mock repository (testify/mock)
2. Test idempotency
3. Test context cancellation
4. Assert expectations
5. Coverage 90%+

# FORMAT
```go
package application_test

func Test{UseCase}_Success(t *testing.T) { ... }
func Test{UseCase}_ValidationError(t *testing.T) { ... }
func Test{UseCase}_RepositoryError(t *testing.T) { ... }
func Test{UseCase}_Unauthorized(t *testing.T) { ... }
```

# EXAMPLE
Use Case: CreatePaymentUseCase
Dependencies: PaymentRepository, EventBus, Logger
```

## 15.3 Prompt Integration Test

```markdown
# ROLE
คุณเป็น Go Test Expert ที่เชี่ยวชาญ Testcontainers

# CONTEXT
Component: {component}
Dependencies: {db}, {redis}, {kafka}

# TASK
เขียน integration test ด้วย testcontainers

# FORMAT
```go
//go:build integration

package {pkg}_test

import (
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestMain(m *testing.M) {
    // setup containers
    code := m.Run()
    // cleanup
    os.Exit(code)
}

func Test{Component}_Integration(t *testing.T) { ... }
```

# CONSTRAINTS
1. Use build tag `integration`
2. Auto-cleanup
3. Skip in short mode
4. Real dependencies (not mocks)
5. Realistic data size

# EXAMPLE
Component: PaymentRepository
Dependencies: PostgreSQL 15
```

## 15.4 Prompt Load Test (k6)

```markdown
# ROLE
คุณเป็น Performance Engineer ที่เชี่ยวชาญ k6

# CONTEXT
API: {endpoint}
Expected load: {load}
SLO: {slo}

# TASK
เขียน k6 load test

## Scenarios
1. Smoke (1 user, 1 min)
2. Load (expected, 10 min)
3. Stress (2x expected)
4. Spike (10x, 1 min)
5. Soak (expected, 2 hours)

# FORMAT
```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  scenarios: {
    smoke: {
      executor: 'constant-vus',
      vus: 1,
      duration: '1m',
    },
    load: {
      executor: 'ramping-vus',
      stages: [
        { duration: '2m', target: 100 },
        { duration: '5m', target: 100 },
        { duration: '2m', target: 0 },
      ],
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<200'],
    http_req_failed: ['rate<0.01'],
  },
};

export default function () {
  // test logic
}
```

# CONSTRAINTS
1. Realistic scenarios
2. SLO-based thresholds
3. Custom metrics
4. CI integration
5. Clear pass/fail criteria
```

---

# บทที่ 16: Prompt CI/CD

## 16.1 Prompt GitHub Actions Workflow

```markdown
# ROLE
คุณเป็น DevOps Engineer ที่เชี่ยวชาญ GitHub Actions

# CONTEXT
Project: {project_name}
Language: Go
Deployment: {target}

# TASK
สร้าง GitHub Actions workflow ที่:

## Required
- [ ] Lint (golangci-lint)
- [ ] Test (unit + integration)
- [ ] Coverage gate (80%+)
- [ ] Security scan (gosec, trivy)
- [ ] Build binary
- [ ] Build container
- [ ] Push to registry
- [ ] Deploy (dev/staging/prod)

# CONSTRAINTS
1. Fast (< 15 min total)
2. Cache dependencies
3. Parallel jobs เมื่อทำได้
4. Fail fast
5. Clear output

# FORMAT
```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:

env:
  GO_VERSION: '1.22'

jobs:
  lint: { ... }
  test: { ... }
  security: { ... }
  build: { ... }
  docker: { ... }
  deploy: { ... }
```
```

## 16.2 Prompt GitLab CI

```markdown
# ROLE
คุณเป็น DevOps Engineer

# TASK
สร้าง .gitlab-ci.yml ที่มี stages:
- lint, test, security, build, deploy

# FORMAT
```yaml
stages:
  - lint
  - test
  - security
  - build
  - deploy

variables:
  GO_VERSION: "1.22"

lint:
  stage: lint
  # ...

test:
  stage: test
  # ...
```
```

---

# บทที่ 17: Prompt Deployment

## 17.1 Prompt Kubernetes Manifests

```markdown
# ROLE
คุณเป็น K8s Expert

# CONTEXT
App: {app_name}
Image: {image}
Port: {port}

# TASK
สร้าง K8s manifests ที่:
- Production-ready
- Zero-downtime deploy
- Auto-scaling
- Health checks
- Security best practices

# REQUIRED MANIFESTS
1. Deployment
2. Service
3. Ingress
4. ConfigMap
5. Secret (template)
6. HPA
7. PDB
8. NetworkPolicy (optional)

# CONSTRAINTS
1. runAsNonRoot
2. readOnlyRootFilesystem
3. Resource limits
4. Liveness + Readiness + Startup probes
5. Anti-affinity
6. PodDisruptionBudget

# FORMAT
[YAML files ทั้งหมดพร้อม comment 2 ภาษา]
```

## 17.2 Prompt Canary Deployment

```markdown
# ROLE
คุณเป็น Deployment Expert

# TASK
ออกแบบ Canary Deployment Strategy

# INPUT
- App: {app}
- Risk level: {risk}
- Blast radius: {blast_radius}

# FORMAT
## 1. Strategy
- Stage 1: 5% (5 min)
- Stage 2: 25% (10 min)
- Stage 3: 50% (15 min)
- Stage 4: 100%

## 2. Metrics to Watch
- Error rate
- Latency (P95)
- Success rate
- Custom business metric

## 3. Auto-Rollback Triggers
- If error rate > 1% for 5 min
- If P95 > 500ms for 5 min
- If success rate < 99%

## 4. Implementation
- Argo Rollouts / Flagger config
- Analysis templates
- Alert config

## 5. Manual Override
- Pause/resume
- Promote
- Rollback

## 6. Communication Plan
- Before: announce
- During: status updates
- After: summary
```

## 17.3 Prompt Zero-Downtime Migration

```markdown
# ROLE
คุณเป็น Database + Deployment Expert

# TASK
ออกแบบ Zero-Downtime Migration

# INPUT
- Change: {change}
- Table size: {size}
- Impact: {impact}

# REQUIRED
## 1. Expand-Contract Plan
[Phase 1: Expand] - Add new
[Phase 2: Migrate] - Backfill + dual write
[Phase 3: Switch] - Read from new
[Phase 4: Contract] - Drop old

## 2. Each Phase
- Actions
- Verification
- Rollback plan
- Duration

## 3. Migration SQL
- Add column (nullable)
- Create index (CONCURRENTLY)
- Backfill script (batch)
- Add constraint (NOT VALID → VALIDATE)

## 4. Code Changes
- Dual write
- Dual read + verify
- Switch read
- Switch write

## 5. Monitoring
- What to watch
- Alert thresholds
- Rollback triggers

## 6. Timeline
| Phase | Duration | Owner |
|-------|----------|-------|
```

---

# บทที่ 18: Prompt Profiling & Performance

## 18.1 Prompt Performance Analysis

```markdown
# ROLE
คุณเป็น Performance Engineer ที่เชี่ยวชาญ Go

# CONTEXT
App: {app}
Issue: {issue}
Metric target: {target}

# INPUT
- Current P95: {current_p95}
- Target P95: {target_p95}
- Load: {load}

# TASK
วิเคราะห์ performance และเสนอ optimization

## Analysis Required
1. Bottleneck identification
2. Root cause (5 Whys)
3. Optimization options (with trade-offs)
4. Recommended solution
5. Implementation plan
6. Verification plan

# FORMAT
## 1. Profiling Plan
- CPU profile command
- Memory profile command
- What to look for

## 2. Likely Bottlenecks
- Based on symptoms

## 3. Optimization Options
| Option | Effort | Impact | Risk |
|--------|--------|--------|------|
| | | | |

## 4. Recommended Solution
[Detailed]

## 5. Implementation
```go
// Before
// After
```

## 6. Benchmark
```go
func BenchmarkX(b *testing.B) { ... }
```

## 7. Verification
- Metrics to watch
- Expected improvement
- Rollback if not better
```

## 18.2 Prompt Cache Strategy

```markdown
# ROLE
คุณเป็น Caching Expert

# CONTEXT
Module: {module}
Data: {data_type}
Access pattern: {pattern}

# TASK
ออกแบบ Caching Strategy

## Analysis
1. Cache type (Redis, in-memory, CDN)
2. Cache pattern (cache-aside, write-through, write-behind)
3. TTL strategy
4. Invalidation strategy
5. Key naming

## Format
### 1. Cache Decision
- Should we cache? (yes/no + reason)
- What to cache?
- What NOT to cache?

### 2. Cache Key Design
```
{module}:{entity}:{id}
{module}:{entity}:list:{filter_hash}
```

### 3. TTL Strategy
| Data | TTL | Reason |
|------|-----|--------|
| | | |

### 4. Invalidation
- On write: delete key
- On update: delete list caches
- Pattern: `{module}:{entity}:*`

### 5. Stampede Prevention
- Single-flight
- Pre-warming
- Distributed lock

### 6. Implementation
```go
// Get with cache-aside pattern
func (s *Service) Get(ctx, id) (*Entity, error) { ... }
```

### 7. Metrics
- Hit rate target: > 80%
- Latency improvement
- Memory cost

### 8. Monitoring
```promql
cache_hits / (cache_hits + cache_misses)
```
```

---

# บทที่ 19: Prompt Scaling

## 19.1 Prompt Scaling Strategy

```markdown
# ROLE
คุณเป็น Scalability Architect

# CONTEXT
App: {app}
Current load: {current_load}
Target load: {target_load}
Bottleneck: {bottleneck}

# TASK
ออกแบบ Scaling Strategy

## Analysis Required
1. Current bottleneck (verified ด้วย data)
2. Scaling options (vertical vs horizontal)
3. Recommended approach
4. Implementation phases
5. Cost impact
6. Risk mitigation

# FORMAT
## 1. Current State
- Metrics (CPU, mem, DB, etc.)
- Bottleneck identified
- Evidence

## 2. Target State
- Required capacity
- SLO to maintain

## 3. Scaling Options
| Option | Pros | Cons | Cost |
|--------|------|------|------|
| Vertical | | | |
| Horizontal | | | |
| Caching | | | |
| Async | | | |
| DB replica | | | |
| Sharding | | | |

## 4. Recommended Plan
[Phased approach]

### Phase 1 (Week 1-2): {action}
### Phase 2 (Week 3-4): {action}
### Phase 3 (Week 5-8): {action}

## 5. Cost Analysis
- Before: $X/mo
- After: $Y/mo
- Cost per user: $Z

## 6. Implementation
[Code + config]

## 7. Monitoring
- Metrics to watch
- Alert thresholds
- Success criteria
```

## 19.2 Prompt Database Scaling

```markdown
# ROLE
คุณเป็น Database Expert

# CONTEXT
DB: PostgreSQL
Size: {size}
Load: {load}
Issue: {issue}

# TASK
ออกแบบ Database Scaling Strategy

## Options to Evaluate
1. Read replicas
2. Connection pooling (PgBouncer)
3. Partitioning
4. Sharding
5. CQRS
6. Cache layer

## Format
### 1. Current Bottleneck
- Query patterns
- Slow queries
- Resource usage

### 2. Option Analysis
| Option | Effort | Impact | Complexity |
|--------|--------|--------|------------|
| | | | |

### 3. Recommended Plan
[Phased]

### 4. Implementation
- Migration SQL
- Code changes
- Config changes

### 5. Monitoring
- Replication lag
- Query latency
- Connection count
```

---

# บทที่ 20: Prompt Incident Response

## 20.1 Prompt Incident Analysis

```markdown
# ROLE
คุณเป็น SRE Expert ที่เชี่ยวชาญ Incident Response

# CONTEXT
Incident: {description}
Severity: {severity}
Duration: {duration}
Impact: {impact}

# INPUT
- Symptom: {symptom}
- Timeline: {timeline}
- Metrics: {metrics}
- Logs: {logs}

# TASK
วิเคราะห์ incident และเสนอ action plan

## Required Output
1. Root cause (5 Whys)
2. Immediate mitigation
3. Permanent fix
4. Prevention
5. Postmortem

# FORMAT
## 1. Incident Summary
- What happened
- Impact
- Duration

## 2. Timeline
| Time | Event |
|------|-------|

## 3. Root Cause Analysis (5 Whys)
1. Why? → ...
2. Why? → ...
3. Why? → ...
4. Why? → ...
5. Why? → **Root Cause**

## 4. Immediate Mitigation
- Action 1
- Action 2

## 5. Permanent Fix
- Change 1
- Change 2

## 6. Prevention
- Monitoring 1
- Process 1
- Test 1

## 7. Action Items
| # | Action | Owner | Due | Priority |
|---|--------|-------|-----|----------|

## 8. Postmortem Doc
[Full blameless postmortem]
```

## 20.2 Prompt Runbook Creation

```markdown
# ROLE
คุณเป็น SRE ที่เชี่ยวชาญ Runbook Writing

# TASK
เขียน Runbook สำหรับ issue: {issue}

# CONTEXT
- Alert name: {alert}
- Component: {component}
- Severity: {severity}

# FORMAT
# Runbook: {Issue Name}

## Alert
`{AlertName}`

## Symptoms
- What users see
- What metrics show
- What logs show

## Quick Diagnosis
```bash
# Check 1
command

# Check 2
command
```

## Likely Causes
1. Cause A (X%)
2. Cause B (Y%)
3. Cause C (Z%)

## Resolution Steps

### If Cause A
```bash
# Fix commands
```

### If Cause B
```bash
# Fix commands
```

### If Cause C
```bash
# Fix commands
```

## Verification
- [ ] Metric back to normal
- [ ] No new errors
- [ ] Smoke test passes

## Rollback
```bash
# If fix doesn't work
kubectl rollout undo deployment/xxx
```

## Escalation
- 15 min: @team-lead
- 30 min: @manager
- 1 hour: @vp-engineering

## Related
- [Dashboard]
- [Logs]
- [Previous incidents]
```

---

# ภาคผนวก A: Master Prompt Template

## A.1 Universal Master Prompt

```markdown
# ═══════════════════════════════════════════
# ROLE
# ═══════════════════════════════════════════
คุณเป็น {expert_role}

# ═══════════════════════════════════════════
# CONTEXT
# ═══════════════════════════════════════════
โปรเจกต์: {project}
Stack: {stack}
Pattern: {pattern}
อ้างอิง: {reference}

# ═══════════════════════════════════════════
# TASK
# ═══════════════════════════════════════════
{task}

# ═══════════════════════════════════════════
# CONSTRAINTS
# ═══════════════════════════════════════════
{constraints}

# ═══════════════════════════════════════════
# FORMAT
# ═══════════════════════════════════════════
{format}

# ═══════════════════════════════════════════
# EXAMPLES
# ═══════════════════════════════════════════
{examples}

# ═══════════════════════════════════════════
# VERIFICATION
# ═══════════════════════════════════════════
ก่อนส่งคำตอบ ตรวจสอบ:
{checklist}
```

## A.2 Prompt for Cursor/Copilot (Short)

```markdown
# Context
Project: icmongolang (Go, Clean Architecture + DDD)
Reference: internal/modules/payment/

# Task
สร้าง {file_type} สำหรับ {module_name} ตาม pattern เดียวกับ payment module

# Constraints
- Comment 2 ภาษา
- Test coverage 80%+
- Domain ห้าม import GORM

# Expected Output
[ระบุไฟล์ที่ต้องการ]
```

---

# ภาคผนวก B: Context Injection Pattern

## B.1 Full Context Block (สำหรับ Long Session)

```markdown
[PROJECT CONTEXT - KEEP IN MEMORY]

## Project: icmongolang
- Go 1.22+
- Clean Architecture + DDD
- Modular Monolith (33 modules)

## Structure
```
internal/modules/{name}/
├── domain/
│   ├── entity/
│   ├── value_object/
│   ├── repository/
│   └── errors/
├── application/
│   └── mocks/
├── infrastructure/
│   └── persistence/postgres/
├── interfaces/
│   └── http/
└── module.go
```

## Dependency Rules
| Layer | Can Import | Cannot Import |
|-------|-----------|---------------|
| Domain | stdlib, uuid, decimal | gorm, gin, chi, redis, kafka |
| Application | Domain | Infrastructure, Interfaces |
| Infrastructure | Domain, Application | Interfaces |
| Interfaces | All | – |

## Shared Packages
- icmongolang/pkg/logger
- icmongolang/pkg/validator
- icmongolang/pkg/httputil
- icmongolang/pkg/jwt

## Conventions
- Table: {module}_{plural}
- Redis key: {module}:{entity}:{id}
- Kafka topic: {module}.{entity}.{action}
- Migration: YYYYMMDD_{module}_{desc}.sql

## Error Mapping
| Domain Error | HTTP |
|--------------|------|
| ErrNotFound | 404 |
| ErrAlreadyExists | 409 |
| ErrInvalid* | 400 |
| ErrUnauthorized | 401 |

[END CONTEXT]
```

## B.2 Per-Task Context

```markdown
[CONTEXT FOR THIS TASK]
Module: {module_name}
File: {file_path}
Related files:
- {file_1}
- {file_2}

Previous decisions:
- {decision_1}
- {decision_2}

[END CONTEXT]

# TASK
{task}
```

---

# ภาคผนวก C: Quick Copy-Paste Library

## C.1 Prompt สร้าง Entity (1-line)

```
สร้าง Go entity ชื่อ "{Name}" ตาม pattern ของ icmongolang:
- Fields: {fields}
- Invariants: {invariants}
- ใช้ uuid, decimal
- มี constructor + behavior methods
- Comment 2 ภาษา
- Test 100%
```

## C.2 Prompt สร้าง Use Case (1-line)

```
สร้าง Go use case "{Verb}{Entity}" ตาม pattern ของ payment module:
- Input: {input}
- Output: {output}
- Business logic: {logic}
- Side effects: {side_effects}
- Test with mock
```

## C.3 Prompt สร้าง Handler (1-line)

```
สร้าง chi HTTP handler สำหรับ {Entity} ตาม pattern ของ payment_handler.go:
- Endpoints: {endpoints}
- Auth middleware required
- Validate request
- Error mapping per convention
```

## C.4 Prompt สร้าง Migration (1-line)

```
สร้าง PostgreSQL migration สำหรับ:
- Change: {change}
- Zero-downtime
- Include rollback
- Idempotent
- Batch backfill if needed
```

## C.5 Prompt Debug (1-line)

```
วิเคราะห์ bug:
- Symptom: {symptom}
- Context: {context}
- ใช้ 5 Whys หา root cause
- เสนอ fix + test
```

## C.6 Prompt Review (1-line)

```
Review โค้ดนี้ตาม criteria:
- Clean Architecture compliance
- Test coverage
- Security issues
- Performance concerns
- Naming & readability
[แนบโค้ด]
```

## C.7 Prompt Refactor (1-line)

```
Refactor โค้ดนี้:
- Behavior เดิม
- ลด complexity
- แยก method < 30 บรรทัด
- Test ผ่าน
[แนบโค้ด]
```

## C.8 Prompt Test (1-line)

```
เขียน test สำหรับ:
- Function: {function}
- Coverage: happy + errors + edge cases
- Table-driven
- testify
```

## C.9 Prompt Docs (1-line)

```
เขียน {doc_type} สำหรับ:
- Module: {module}
- Audience: {audience}
- Format: {format}
- Length: {length}
```

## C.10 Prompt Incident (1-line)

```
วิเคราะห์ incident:
- Symptom: {symptom}
- Impact: {impact}
- Timeline: {timeline}
- Root cause (5 Whys)
- Action items
```

---

## 📋 ตารางสรุป Prompt ทั้งหมด

| หมวด | Prompt | ใช้เมื่อ |
|------|--------|----------|
| **สร้าง** | Master Module | สร้าง module ใหม่ทั้งหมด |
| | Entity | สร้าง domain entity |
| | Value Object | สร้าง VO |
| | Use Case | สร้าง business logic |
| | Repository | สร้าง data access |
| | Handler | สร้าง HTTP endpoint |
| | Migration | สร้าง DB schema |
| **แก้ไข** | Impact Analysis | ประเมินก่อนแก้ |
| | Add Field | เพิ่ม field ปลอดภัย |
| | Migration | เปลี่ยน type/schema |
| | Refactor | ปรับปรุงโค้ด |
| **ขาย** | README | เอกสารขาย |
| | Sales Kit | One-pager, case study |
| | Pricing | ตั้งราคา + ROI |
| **Test/Deploy** | Unit Test | Test domain/use case |
| | Integration | Test infra |
| | Load Test | Test performance |
| | CI/CD | Pipeline |
| | K8s | Manifests |
| **Maintain** | Profile | วิเคราะห์ performance |
| | Scaling | วางแผน scale |
| | Incident | จัดการ incident |
| | Runbook | เขียนคู่มือ |

---

## 🎯 Tips การใช้ Prompt ให้ได้ผล

### 1. Iterative Refinement

```
Round 1: "สร้าง module X"
Round 2: "ปรับ domain ให้ validate Y"
Round 3: "เพิ่ม test สำหรับ Z"
Round 4: "ทำให้ error messages consistent"
```

### 2. Give Feedback

```
❌ "ไม่ถูก แก้ใหม่"
✅ "ข้อ 3 ไม่ถูก เพราะ violate dependency rule.
   Domain ไม่ควร import gorm.
   กรุณาแก้ให้ domain ไม่มี GORM tag"
```

### 3. Ask for Alternatives

```
"เสนอ 3 solutions พร้อม trade-offs:
- Performance
- Complexity  
- Maintainability
- Migration effort"
```

### 4. Verify Before Accept

```
"ก่อนส่งคำตอบ ให้ตรวจสอบ:
- [ ] รันได้จริง
- [ ] Test ผ่าน
- [ ] ไม่มี import ที่ห้าม
- [ ] Comment ครบ 2 ภาษา"
```

### 5. Chain Prompts

```
1. "วิเคราะห์ requirement"
2. "ออกแบบ domain model"
3. "Implement"
4. "เขียน test"
5. "Review"
```

---

## 📚 Prompt Library ตามสถานการณ์

### สถานการณ์: เพิ่งเริ่มโปรเจกต์ใหม่

```
1. Master Module Prompt (สร้าง structure)
2. Domain Layer Prompt (entities)
3. Application Layer Prompt (use cases)
4. Infrastructure Prompt (repos)
5. Interface Prompt (handlers)
6. Test Prompt (unit tests)
7. Migration Prompt (schema)
```

### สถานการณ์: แก้ Bug Production

```
1. Incident Analysis Prompt
2. Debug Prompt (5 Whys)
3. Test Prompt (reproduce bug)
4. Fix Prompt
5. Postmortem Prompt
```

### สถานการณ์: เพิ่ม Feature

```
1. Impact Analysis Prompt
2. Domain Design Prompt
3. Use Case Prompt
4. Handler Prompt
5. Test Prompt
6. Migration Prompt
7. Documentation Prompt
```

### สถานการณ์: ขาย Module

```
1. README Prompt
2. Sales Kit Prompt
3. Pricing Prompt
4. Case Study Prompt
5. Product One-Pager Prompt
```

### สถานการณ์: Scale ระบบ

```
1. Performance Analysis Prompt
2. Profiling Prompt
3. Database Scaling Prompt
4. Cache Strategy Prompt
5. K8s Scaling Prompt
```

### สถานการณ์: Incident

```
1. Incident Analysis Prompt
2. Root Cause Prompt
3. Runbook Prompt
4. Postmortem Prompt
5. Prevention Prompt
```

---

## 💡 Advanced Prompt Techniques

### Technique 1: Persona Stacking

```
คุณเป็นทั้ง:
1. Senior Go Developer (code quality)
2. SRE (production concerns)
3. Security Engineer (threats)
4. Cost Optimizer (efficiency)

ให้ความเห็นจากทุกมุมมอง
```

### Technique 2: Devil's Advocate

```
เสนอ solution ที่ดีที่สุด
จากนั้นเสนอ:
- จุดอ่อนของ solution นี้
- วิธีที่มันจะ fail
- Alternative ที่ดีกว่า
```

### Technique 3: Constraint-Based Design

```
ออกแบบโดยมี constraints:
- Go 1.22 เท่านั้น
- ห้าม external dependencies
- Coverage 90%+
- Latency < 50ms
- Memory < 100MB

จะออกแบบยังไง?
```

### Technique 4: Iterative Refinement Chain

```
Prompt 1: "สร้าง draft"
Prompt 2: "วิเคราะห์ draft นี้ หาจุดอ่อน"
Prompt 3: "ปรับปรุงตาม feedback"
Prompt 4: "Verify ว่าดีขึ้นจริง"
```

### Technique 5: Multi-Perspective Review

```
Review โค้ดนี้จาก 4 มุมมอง:

1. Developer: readability, maintainability
2. SRE: production concerns, monitoring
3. Security: vulnerabilities, threats
4. Business: cost, time to market

แต่ละมุมมอง ให้:
- จุดแข็ง
- จุดอ่อน
- คำแนะนำ
```

---

## 📝 ข้อมูลเอกสาร

**ชื่อเอกสาร:** AI Prompt Templates Library for Go Module Work
**เวอร์ชัน:** 1.0
**วันที่:** เมษายน 2026
**จำนวนหน้า:** ~90 หน้า (ประมาณ)
**ระดับ:** Intermediate - Advanced

**ผู้อ่านเป้าหมาย:**
- Go Developer ที่ใช้ AI เป็น copilot
- Tech Lead ที่ต้อง instruct AI
- Team ที่ต้องการ standard prompts

**AI Platforms ที่ใช้ได้:**
- Claude (Anthropic)
- ChatGPT (OpenAI)
- Gemini (Google)
- Cursor IDE
- GitHub Copilot
- Codeium

**เอกสารที่เกี่ยวข้อง:**
- เล่ม 1: สร้าง Module
- เล่ม 2: แก้ไข Module
- เล่ม 3: ขาย Module
- เล่ม 4: Test & Deploy
- เล่ม 5: Maintain & Scale
- **เอกสารนี้**: AI Prompt Templates (ใช้ร่วมกับทุกเล่ม)

---

**END OF PROMPT LIBRARY**

> 🎯 **วิธีใช้:**
> 1. Copy prompt ที่ตรงกับงาน
> 2. แทนที่ `{placeholder}` ด้วยข้อมูลจริง
> 3. Paste ให้ AI (Claude, ChatGPT, Cursor)
> 4. Iterate จนได้ผลลัพธ์ที่ต้องการ
> 5. Verify ก่อนใช้จริง

> 💡 **Tips:**
> - ใช้ Context Block สำหรับ long session
> - ใช้ 1-line prompt สำหรับงานเร็ว
> - Chain prompts สำหรับงานซับซ้อน
> - Verify ทุก output ของ AI

---

**พิมพ์เมื่อ:** เมษายน 2026
**ผู้จัดทำ:** ทีมสถาปัตยกรรมซอฟต์แวร์ icmongolang
**ติดต่อ:** kongnakornjantakun@gmail.com
**Line:** kongnakornna

---

**📞 พร้อมเพย์:** 0955088091
*(สนับสนุนการสร้างเนื้อหาแนวนี้ต่อ)*