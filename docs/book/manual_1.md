# 📕 เล่ม 1: คู่มือสร้าง Module ใหม่ฉบับสมบูรณ์
## The Complete Guide to Creating Go Modules with Clean Architecture + DDD

> **เวอร์ชัน 1.0 (เมษายน 2026)**
> เอกสารต้นแบบสำหรับนักพัฒนา Go ทุกระดับ
> อ้างอิงจากโปรเจกต์ `icmongolang` 
> Modular Monolith + Clean Architecture + DDD

---

# สารบัญ

**ภาคที่ 1: พื้นฐานและปรัชญา**
1. [บทนำ — ทำไมต้อง Clean Architecture + DDD](#บทที่-1-บทนำ)
2. [ปรัชญาและหลักการ](#บทที่-2-ปรัชญาและหลักการ)

**ภาคที่ 2: การเตรียมการ**
3. [การเตรียมความพร้อม](#บทที่-3-การเตรียมความพร้อม)
4. [Naming Conventions และ Folder Structure](#บทที่-4-naming-conventions)

**ภาคที่ 3: ลงมือสร้าง Module 12 ขั้น**
5. [ออกแบบ Domain Model](#บทที่-5-ออกแบบ-domain-model)
6. [สร้าง Domain Layer](#บทที่-6-สร้าง-domain-layer)
7. [สร้าง Application Layer](#บทที่-7-สร้าง-application-layer)
8. [สร้าง Infrastructure Layer](#บทที่-8-สร้าง-infrastructure-layer)
9. [สร้าง Interface Layer](#บทที่-9-สร้าง-interface-layer)
10. [Migration และ Schema](#บทที่-10-migration-และ-schema)
11. [Wire-up และ Bootstrap](#บทที่-11-wire-up-และ-bootstrap)
12. [Testing](#บทที่-12-testing)

**ภาคที่ 4: กรณีศึกษาจริง**
13. [กรณีศึกษา: Payment Module](#บทที่-13-กรณีศึกษา-payment-module)
14. [Checklist และ Best Practices](#บทที่-14-checklist-และ-best-practices)

**ภาคผนวก**
- [A. Code Templates ทั้งหมด](#ภาคผนวก-a-code-templates)
- [B. Common Pitfalls](#ภาคผนวก-b-common-pitfalls)
- [C. Quick Reference Card](#ภาคผนวก-c-quick-reference-card)

---
# โครงสร้างโปรเจกต์ icmongolang (Structure Map)

รวบรวมจากการตรวจสอบ live filesystem ณ root ` icmongolang`
คำศัพท์ 3 กลุ่ม: **ระบบ** = Global infrastructure / root / อื่น ๆ ที่ไม่ใช่ตัว module, **รวม** = shared packages / ร่วมกันทุก module, **แยกรราย** = แต่ละ module แยกตามโฟลเดอร์ของตัวเอง

ไฟล์บางรายการถูกตัดออกจากแผนที่: `vendor/`, `.git/`, `app-postgres-data/`, `tmp/`

---

## ระบบ (โครงสร้างระดับ Global)

```
icmongolang
├─ main.go    # entrypoint หลัก
├─ go.mod / go.sum
├─ Makefile / build.ps1 / run.ps1
├─ air.cmd / air.ps1 / .air.toml / .air.toml.linux / .air_1.toml / .air_2.toml
├─ Dockerfile (dev / prod / ollama-exporter)
├─ docker-compose.yml (base / .dev / .prod)
├─ .env / .env.example / .gitattributes / .gitignore
├─ icmongolang.exe / ollama-exporter.exe / package-lock.json / queryex
├─ Note.md / NoteAI.md / README.md / README_IO.md
├─ .github/            # CI/CD workflows
├─ .gitlab-ci/         # GitLab CI
├─ .vscode/
├─ cmd/
│  ├─ api/
│  ├─ kafka/
│  ├─ ollama-exporter/
│  ├─ scheduler/
│  ├─ websocket/
│  ├─ workers/
│  ├─ initdata.go
│  ├─ migrate.go
│  ├─ root.go
│  ├─ serve.go
│  ├─ vectordata.go
│  └─ worker.go
├─ config/
│  ├─ config.default.yml
│  ├─ config.dev.yml
│  ├─ config.go
│  └─ verify_env_test.go
├─ db/
│  ├─ icmon.sql
│  ├─ public.sql
│  ├─ public_DB.sql
│  ├─ public_icmongolang.sql
│  ├─ public_icmongolang_v1.sql
│  └─ sd_iot_device.sql
├─ migrations/
│  ├─ 20250619_kafka_tables.sql
│  ├─ 20250619_websocket_tables.sql
│  ├─ 20260712_new_modules_schema.sql
│  ├─ 20260712_seed_data.sql
│  ├─ 20260827_slow_sql_device_indexes.sql
│  ├─ 20260829_fullschedule.sql
│  ├─ 20260830_fullschedule_scope.sql
│  ├─ 20260831_fullschedule_legacy_backfill.sql
│  ├─ 20260832_fs_schedule_sdiot_columns.sql
│  ├─ 20260833_fullschedule_baseline_audit.sql
│  ├─ 20260834_fullschedule_sort_order.sql
│  ├─ 20260835_fullschedule_baseline_gap.sql
│  ├─ 20260901_fs_schedule_demo_data.sql
│  ├─ 20260901_fs_schedule_history_settings_seed.sql
│  ├─ db.sql
│  ├─ icmon.sql
│  ├─ icmon_database.sql
│  ├─ public.sql
│  ├─ public_db.sql
│  └─ sd_user_access_menu.sql
├─ monitoring/
│  ├─ prometheus.yml
│  ├─ kafka-dashboard.json
│  └─ grafana/
│     ├─ dashboards/
│     │  ├─ elasticsearch/ *-dashboard.json
│     │  ├─ golang/ *-dashboard.json
│     │  ├─ health/ *-dashboard.json
│     │  ├─ influxdb/ *-dashboard.json
│     │  ├─ kafka/ *-dashboard.json
│     │  ├─ mqtt/ *-dashboard.json
│     │  ├─ nodered/ *-dashboard.json
│     │  ├─ ollama/ *-dashboard.json
│     │  ├─ postgres/ *-dashboard.json
│     │  ├─ redis/ *-dashboard.json
│     │  └─ websocket/ *-dashboard.json
│     └─ provisioning/
│        ├─ dashboards/
│        │  ├─ elasticsearch.yml / golang.yml / health.yml / influxdb.yml
│        │  ├─ kafka.yml / mqtt.yml / nodered.yml / ollama.yml
│        │  └─ postgres.yml / redis.yml / websocket.yml
│        └─ datasources/
│           └─ prometheus.yml
├─ elk/                # Elasticsearch / Logstash / Kibana
├─ influxdb/
├─ mqtt/
├─ linux/
├─ scripts/
├─ test/
├─ postman/
├─ docs/
│  ├─ structure-map.md            # ไฟล์นี้
│  ├─ docs.md
│  ├─ docs.go (swagger embed)
│  ├─ swagger.json / swagger.yaml
│  └─ modules_dev/
│     ├─ Template_Modules.md
│     ├─ Template_Modules_PDPA.md
│     ├─ Template_prompt_ai.md
│     ├─ PDPA_Modules_V4.md
│     └─ PDPA_Modules_V5.md
├─ ai/
├─ internal/            # Go source (ดู รวม และ แยกรราย ด้านล่าง)
└─ pkg/                 # shared packages (ดู รวม)
```

---

## รวม (โครงสร้างระดับ shared / ร่วมกัน)

```
internal/                       # application-level shared + modules
├─ server/
│  ├─ api_cache_adapter.go
│  ├─ handlers.go
│  └─ server.go
├─ delivery/
│  └─ rest/
│     ├─ router.go
│     └─ middleware/monitoring.go
├─ middleware/
│  ├─ cors.go / jwtauth.go / jwtauth_test.go / logging.go
│  ├─ middleware.go / monitoring.go / prometheus.go
│  ├─ rate_limit.go / security.go
├─ models/
│  ├─ air_control.go / base.go / device.go / models.go / payment.go
│  ├─ refresh_token.go / sd_user.go / sd_user_role.go
│  ├─ sd_user_roles_access.go / sd_user_roles_permission.go
│  ├─ user.go / ws_models.go
├─ repository/
│  ├─ pg.go / redis.go
│  └─ ReadMe_Modules.md
├─ distributor/
│  └─ distributor.go
├─ processor/
│  └─ processor.go
├─ usecase/
│  └─ usecase.go
├─ worker/
│  └─ worker.go
└─ template/                    # ต้นแบบสำหรับสร้าง module ใหม่
   ├─ handler.go / migration.sql / pg_repository.go / README.md
   ├─ redis_repository.go / usecase.go / worker.go
   ├─ delivery/http/handler.go / routes.go
   ├─ distributor/distributor.go
   ├─ models/model.go
   ├─ presenter/presenter.go
   ├─ processor/processor.go
   ├─ repository/pg_repository.go / redis_repository.go
   └─ usecase/usecase.go

pkg/                            # 21 shared packages
├─ cryptpass
├─ db
├─ elasticsearch
├─ emailTemplates
├─ helpers
├─ http-swagger
├─ httpErrors
├─ influxdb
├─ jwt
├─ kafka
├─ llm
├─ logger
├─ mqtt
├─ report
├─ responses
├─ secureRandom
├─ sendEmail
├─ transaction
├─ utils
├─ vectordb
└─ websocket
```

> หมายเหตุ: `internal/modules/*` เป็นกลุ่ม **แยกรราย** (ดูด้านล่าง)

---

## แยกรราย (โครงสร้างระดับ module — 33 modules)

### แบบที่ 1 — ลำดับชั้นเต็ม 6 ชั้น (alarm)
```
internal/modules/alarm/
├─ delivery/ + presenter/ + processor/ + models/
├─ repository/ + usecase/
└─ (handler / pg_repository / usecase อยู่ใต้ sub-folder ตามโครงสร้าง)
```

### แบบที่ 2 — 4 ชั้นแบบบาง (apimanager / control / flowengine)
```
internal/modules/<name>/
├─ delivery/ + presenter/
└─ repository/ + usecase/
```

### แบบที่ 3 — 5 ชั้น + provider (auditlog / notifier)
```
internal/modules/<name>/
├─ delivery/ + presenter/ + provider/
└─ repository/ + usecase/
```

### แบบที่ 4 — main file ที่ root (auth/batch/customer/dashboard/document/email/i18n/items/job/payment/settings/users/wos + ฯลฯ)
```
internal/modules/<name>/
├─ delivery/ + presenter/
├─ repository/ + usecase/
├─ handler.go
├─ pg_repository.go
└─ usecase.go
```

### แบบที่ 5 — จัดกลุ่มตาม service (elasticsearch / influxdb / mqtt / vectordata / realtime / websocket)
```
internal/modules/<name>/
├─ delivery/ + presenter/
└─ usecase/
```

### module พิเศษแบบละเอียด
```
internal/modules/fullschedule/
├─ delivery/ + executor/ + presenter/
├─ repository/ + scheduler/ + usecase/
└─ filter.go / handler.go / master.go / model.go / model_test.go
   pg_repository.go / usecase.go

internal/modules/iot/
├─ delivery/ + iothelper/ + models/ + presenter/ + repository/ + usecase/

internal/modules/kafka/
├─ delivery/ + models/ + repository/ + usecase/ + ws/

internal/modules/pdpa/
├─ application/ + domain/ + infrastructure/ + interfaces/ + repository/

internal/modules/queue/
├─ manager.go / manager_test.go / noop_queue.go / queue_test.go

internal/modules/report/
├─ delivery/ + usecase/ + company.go + handler.go

internal/modules/purchaseorder/  และ quotation/ และ wos/
├─ delivery/ + presenter/ + repository/ + usecase/
└─ handler.go / pg_repository.go / usecase.go

internal/modules/settings/
├─ delivery/ + presenter/ + repository/ + usecase/
└─ handler.go / repository.go / specs_alarm_logs.go / specs_device_schedule.go
   specs_integration.go / specs_master.go / spec_with.go / types.go / types_test.go / usecase.go

internal/modules/users/        # module อเนกประสงค์สุด
├─ delivery/ + distributor/ + doc/ + presenter/ + processor/
├─ repository/ + usecase/
└─ handler.go / pg_repository.go / pg_repository_test.go
   redis_repository.go / usecase.go / worker.go

internal/modules/websocket/
├─ delivery/ + models/ + presenter/ + repository/ + usecase/
└─ interface.go
```

### รายการ 33 modules (แยกตามโฟลเดอร์ `internal/modules/`)
```
alarm        apimanager   auditlog     auth         batch
control      customer     dashboard    document     elasticsearch
email        flowengine   fullschedule i18n         influxdb
iot          items        job          kafka        mqtt
notifier     payment      pdpa         purchaseorder queue
quotation    realtime     report       settings     users
vectordata   websocket    wos
```

---

## Template Prompt สำหรับสร้าง / แก้ไข Module (ใช้งานซ้ำได้)

เมื่อต้อง สร้าง module ใหม่ หรือ แก้ไข module ที่มีอยู่ ให้ใช้ prompt ต่อไปนี้

```
สร้าง/แก้ไขโมดูลใหม่ในโปรเจกต์ icmongolang

1. ตั้งชื่อโมดูล: <module_name>
2. ดูโครงสร้างที่ได้จาก docs\Template_Module.md 
3. สร้างโฟลเดอร์ที่จำเป็นภายใต้ internal\modules\<module_name>:
   - delivery\http\handler.go, routes.go
   - models\model.go
   - repository\pg_repository.go (และ redis_repository.go ถ้าต้องการใช้ Redis)
   - businesslogic\  (ถ้ามี logic แยก)
   - optional: processor\engine\session\template\ + es.go / influx.go / vectordb.go
4. ใช้ internal\modules\users และ internal\template เป็นตัวอย่างอ้างอิง
   รวมถึง internal\repository\ReadMe_Modules.md
5. ตรวจสอบว่า import package ใช้ชื่อจริงจาก pkg\ เท่านั้น:
   pkg\cryptpass, pkg\db, pkg\elasticsearch, pkg\emailTemplates, pkg\helpers,
   pkg\http-swagger, pkg\httpErrors, pkg\influxdb, pkg\jwt, pkg\kafka, pkg\llm,
   pkg\logger, pkg\mqtt, pkg\report, pkg\responses, pkg\secureRandom,
   pkg\sendEmail, pkg\transaction, pkg\utils, pkg\vectordb, pkg\websocket
   (ห้ามใช้ package ที่ไม่อยู่ในรายการนี้ - ตรวจสอบ live listing ด้วย)
6. ถ้าเป็น module ใหม่:
   - เพิ่ม schema/migration ไฟล์ใน migrations\ (ชื่อ YYYYMMDD_<module>_*.sql)
   - เส้นทาง route ลงทะเบียนใน internal\delivery\rest\router.go
7. ถ้าเป็น module เดิม: แก้ไขไฟล์เฉพาะในโครงสร้างเดิมของโมดูลนั้นเท่านั้น
   (ดู 'แยกรราย' ใน structure-map.md เพื่อหาแบบที่ตรงกับโมดูล)
8. สร้าง  Unit Test ตัวอย่าง
    internal\modules\<module_name>\myapp
            ├── go.mod
            ├── calculator.go
            └── calculator_test.go
9. รัน go build ./... และ/หรือ go test ./... ให้ผ่านก่อนจบ
```
---


# บทที่ 1: บทนำ

## 1.1 ปัญหาที่พบบ่อยในโปรเจกต์ Go ขนาดกลาง-ใหญ่

เมื่อทีมเริ่มพัฒนาโปรเจกต์ Go ที่มีมากกว่า 5 โมดูล มักเจอปัญหาซ้ำๆ เหล่านี้:

| ปัญหา | อาการ | ผลกระทบ |
|---|---|---|
| **Business logic กระจาย** | Handler มี if-else ยาว 200 บรรทัด | แก้ไขยาก ทดสอบไม่ได้ |
| **Coupling สูง** | เปลี่ยน DB ต้องแก้ handler ด้วย | Refactor ยาก |
| **ไม่มี Test** | `if err != nil { panic(err) }` | บั๊กใน production |
| **Duplicate code** | Copy-paste validation ทุก handler | Maintain หลายที่ |
| **Naming ยุ่งเหยิง** | `GetUserData`, `FetchUser`, `LoadUser` | สับสน |
| **Test ต้อง spin infra** | Test ต้องเปิด Postgres, Redis | CI ช้า, flaky |

## 1.2 แนวทางแก้: Clean Architecture + DDD

**Clean Architecture (Robert C. Martin)** = แยกโค้ดเป็นชั้น (Layers) โดยให้ Business Logic อยู่ตรงกลาง ไม่ผูกกับ Framework, Database, UI

**Domain-Driven Design (Eric Evans)** = ออกแบบซอฟต์แวร์ให้สะท้อน Domain จริงของธุรกิจ ผ่าน Ubiquitous Language, Aggregate, Value Object

**รวมกัน →** ซอฟต์แวร์ที่:
- ✅ Test ได้โดยไม่ต้อง spin infra
- ✅ เปลี่ยน DB/Framework ได้โดยไม่กระทบ business
- ✅ Business rule อยู่ที่เดียว ชัดเจน
- ✅ ทีมใหม่เข้าใจเร็ว

## 1.3 เมื่อไหร่ควรใช้ / ไม่ควรใช้

### ✅ ควรใช้เมื่อ:

- Business logic ซับซ้อน มี state machine
- ต้องรองรับ multi-tenant, multi-currency
- ทีม > 2 คน, โมดูล > 5
- ต้องสลับ technology ในอนาคต
- ต้อง compliance (PDPA, GMP, ISO)
- ต้องการ test coverage > 80%

### ❌ ไม่ควรใช้เมื่อ:

- CRUD ง่ายๆ ไม่มี business rules
- Prototype / PoC < 1 สัปดาห์
- Script ขนาดเล็ก (batch, migration tool)
- ทีมเดียว < 3 คน, โมดูล < 3

## 1.4 Trade-offs

| ข้อดี | ข้อเสีย |
|---|---|
| Test ได้ง่าย | Boilerplate มากขึ้น (~2x ไฟล์) |
| Maintain ง่ายระยะยาว | Learning curve สูง |
| เปลี่ยน tech ได้ | Setup ครั้งแรกใช้เวลา 1-2 วัน |
| ทีมขยายได้ | Over-engineering ถ้าโปรเจกต์เล็ก |

## 1.5 ตัวอย่างจริงจาก icmongolang

```
icmongolang/
├── internal/modules/       ← 33 modules ทุกตัวใช้ pattern เดียวกัน
│   ├── auth/
│   ├── payment/
│   ├── purchaseorder/
│   ├── iot/
│   └── ...
├── pkg/                    ← 21 shared packages
└── cmd/                    ← 5 entry points
```

ทุก module มี **โครงสร้างเดียวกัน** ทำให้:
- ย้าย developer ข้าม module ได้ทันที
- สร้าง module ใหม่ใช้เวลา 1-2 วัน
- Test/CI/CD ใช้ pipeline เดียวกัน

---

# บทที่ 2: ปรัชญาและหลักการ

## 2.1 Dependency Rule — กฎเหล็กเดียวที่ต้องจำ

```
┌─────────────────────────────────────────┐
│  Interfaces (HTTP/WS/gRPC/CLI)          │  ← รู้จักทุก Layer
├─────────────────────────────────────────┤
│  Infrastructure (Postgres/Redis/Kafka)  │  ← รู้จัก Domain + Application
├─────────────────────────────────────────┤
│  Application (Use Cases)                │  ← รู้จัก Domain เท่านั้น
├─────────────────────────────────────────┤
│  Domain (Entity, VO, Repository IF)     │  ← ❤️ ไม่รู้จักใครเลย
└─────────────────────────────────────────┘
        ↑
   Dependency ชี้เข้าข้างในเท่านั้น (Inward)
```

**ตาราง Dependency:**

| Layer | import ได้ | ห้าม import |
|---|---|---|
| **Domain** | `stdlib`, `uuid`, `decimal` | gorm, gin, chi, sarama, redis, kafka |
| **Application** | Domain | Infrastructure, Interfaces |
| **Infrastructure** | Domain, Application | Interfaces |
| **Interfaces** | ทุก Layer | – |

## 2.2 ทำไมกฎนี้สำคัญ?

**ตัวอย่าง: ถ้า Domain import GORM**

```go
// ❌ ผิด
package entity

import "gorm.io/gorm"  // ← GORM leak เข้า Domain

type Payment struct {
    ID     uuid.UUID
    Amount decimal.Decimal
    db     *gorm.DB  // ← Domain รู้จัก ORM แล้ว
}
```

**ผลเสีย:**
- Test Domain ต้อง mock GORM
- เปลี่ยน ORM ต้องแก้ Entity
- Entity ไม่ pure แล้ว

**✅ ถูกต้อง:**
```go
package entity

import (
    "time"
    "github.com/google/uuid"
    "github.com/shopspring/decimal"
    // ไม่มี GORM!
)

type Payment struct {
    ID     uuid.UUID
    Amount decimal.Decimal
    Status valueobject.PaymentStatus
}
```

## 2.3 Layer Responsibilities

### 2.3.1 Domain Layer — หัวใจ

**หน้าที่:** เก็บ business rules ที่ไม่เปลี่ยนตาม technology

**ประกอบด้วย:**
| ส่วน | บทบาท | ตัวอย่าง |
|---|---|---|
| **Entity** | Object ที่มี identity, มี lifecycle | `Payment`, `Order` |
| **Value Object** | Object ไม่มี identity, immutable | `Money`, `Address` |
| **Aggregate Root** | ประตูเข้า aggregate | `Order` (เข้า `OrderItem` ผ่าน `Order`) |
| **Repository Interface** | สัญญาเข้าถึงข้อมูล | `PaymentRepository` |
| **Domain Service** | Logic ที่ไม่ผูก entity ตัวเดียว | `TransferService` |
| **Domain Event** | เหตุการณ์ที่เกิดขึ้น | `PaymentSucceeded` |
| **Domain Error** | Error ที่ธุรกิจกำหนด | `ErrInsufficientBalance` |

**ตัวอย่าง Invariant:**
```go
func (p *Payment) Refund() error {
    // Invariant: refund ได้เฉพาะเมื่อ SUCCESS
    if p.Status != valueobject.PaymentStatusSuccess {
        return domainerrors.ErrCannotRefund
    }
    // Invariant: refund ได้ครั้งเดียว
    if p.RefundedAt != nil {
        return domainerrors.ErrAlreadyRefunded
    }
    now := time.Now()
    p.Status = valueobject.PaymentStatusRefunded
    p.RefundedAt = &now
    p.UpdatedAt = now
    return nil
}
```

### 2.3.2 Application Layer — ผู้ประสานงาน

**หน้าที่:** Orchestrate — เรียก Domain หลายตัว, จัดการ transaction, side effects

**ประกอบด้วย:**
- Use Case (1 ไฟล์ ต่อ 1 use case)
- Input/Output DTO
- Mapper (Entity ↔ DTO)

**กฏ:**
- ✅ เรียก Domain behavior
- ✅ เรียก Repository interface
- ✅ Publish Domain Event
- ❌ ห้ามมี SQL
- ❌ ห้ามมี HTTP
- ❌ ห้ามใช้ `*gorm.DB` ตรง

**ตัวอย่าง:**
```go
func (uc *CreatePaymentUseCase) Execute(ctx context.Context, input CreatePaymentInput) (*CreatePaymentOutput, error) {
    // 1. ตรวจ duplicate
    existing, _ := uc.repo.FindByOrderID(ctx, input.OrderID)
    if existing != nil {
        return nil, domainerrors.ErrPaymentAlreadyExists
    }

    // 2. สร้าง entity (validate อยู่ใน domain)
    payment, err := entity.NewPayment(input.UserID, input.OrderID, input.Amount, input.Currency, input.Method)
    if err != nil {
        return nil, err
    }

    // 3. Persist
    if err := uc.repo.Save(ctx, payment); err != nil {
        return nil, err
    }

    // 4. Publish event (side effect)
    _ = uc.eventBus.Publish(ctx, "payment.created", events.PaymentCreatedEvent{
        PaymentID: payment.ID,
        UserID:    payment.UserID,
    })

    // 5. Return output
    return &CreatePaymentOutput{ID: payment.ID, Status: string(payment.Status)}, nil
}
```

### 2.3.3 Infrastructure Layer — Adapter

**หน้าที่:** implement interface ที่ Domain/Application ต้องการ

**ประกอบด้วย:**
- GORM Model + Repository Implementation
- Redis Cache
- Kafka Producer/Consumer
- External API Client (Stripe, SMTP, LLM)

**กฏ:**
- ✅ implement interface
- ✅ map Model ↔ Entity
- ❌ ห้ามมี business logic

### 2.3.4 Interface Layer — ประตู

**หน้าที่:** รับ request จากโลกภายนอก แปลงเป็น use case call

**ประกอบด้วย:**
- HTTP Handler
- Route Registration
- Request/Response DTO
- Middleware (Auth, CORS, RateLimit)

**กฏ:**
- ✅ Validate input
- ✅ ดึง user จาก context
- ❌ ห้าม business logic
- ❌ ห้ามเรียก Repository ตรง

## 2.4 Domain-Driven Design Building Blocks

### 2.4.1 Entity vs Value Object

| Aspect | Entity | Value Object |
|---|---|---|
| **Identity** | มี ID (UUID) | ไม่มี ID |
| **Mutability** | เปลี่ยน state ได้ | Immutable |
| **Equality** | เทียบด้วย ID | เทียบด้วยค่า |
| **Lifecycle** | มี lifecycle | ไม่มี |
| **ตัวอย่าง** | `Payment`, `User` | `Money`, `Address` |

**ตัวอย่าง Value Object:**
```go
type Money struct {
    Amount   decimal.Decimal
    Currency string
}

// Immutable — คืนค่าใหม่เสมอ
func (m Money) Add(other Money) (Money, error) {
    if m.Currency != other.Currency {
        return Money{}, errors.New("currency mismatch")
    }
    return Money{Amount: m.Amount.Add(other.Amount), Currency: m.Currency}, nil
}
```

### 2.4.2 Aggregate Root

**นิยาม:** Entity หลักที่ควบคุม entity ย่อย ต้องเข้า aggregate ผ่าน root เท่านั้น

**ตัวอย่าง:** Order + OrderItems
```go
// ✅ เข้าผ่าน Order (root)
order.AddItem(product, qty, price)

// ❌ ห้ามสร้าง OrderItem ตรงๆ
item := NewOrderItem(orderID, ...)  // ← ผิด
```

**กฎ:**
- Aggregate root เป็นเจ้าของ business invariant
- Transaction ครอบ 1 aggregate
- Reference ระหว่าง aggregate → ใช้ ID ไม่ใช่ pointer

### 2.4.3 Ubiquitous Language

ใช้คำเดียวกันทั้งทีม ทั้งในโค้ด, เอกสาร, การคุย

| ❌ ผิด | ✅ ถูก |
|---|---|
| `getUserData()` | `findUserByID()` |
| `doPayment()` | `capturePayment()` |
| `processOrder()` | `placeOrder()` |
| `checkAndSave()` | `validateAndPersist()` |

## 2.5 SOLID ในบริบท Go

| หลัก | ความหมาย | ตัวอย่างใน Go |
|---|---|---|
| **S**RP | 1 อย่าง 1 หน้าที่ | 1 use case ต่อ 1 ไฟล์ |
| **O**CP | เปิดขยาย ปิดแก้ | เพิ่ม discount type ไม่แก้ existing |
| **L**SP | Subtype แทน Supertype ได้ | Interface implementation |
| **I**SP | Interface เล็ก เฉพาะทาง | `PaymentRepo` แยกจาก `RefundRepo` |
| **D**IP | พึ่ง abstraction | ใช้ interface ไม่ใช่ concrete |

---

# บทที่ 3: การเตรียมความพร้อม

## 3.1 ตรวจสอบ Environment

```bash
# Go version (ต้อง 1.22+)
go version
# ควรได้: go version go1.22.x

# เครื่องมือที่ต้องติดตั้ง
go install github.com/air-verse/air@latest              # hot-reload
go install github.com/swaggo/swag/cmd/swag@latest       # Swagger
go install github.com/securego/gosec/v2/cmd/gosec@latest # security scan
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest # lint

# ตรวจสอบ PATH
which air swag gosec golangci-lint

# Windows 11 ที่ AppLocker block exe → ใช้ go run แทน
# go run github.com/air-verse/air@latest
```

## 3.2 ตรวจสอบ Shared Packages ที่มีอยู่

```bash
# ดูว่า pkg/ มี package อะไรแล้ว (ต้องใช้จากรายการนี้เท่านั้น)
ls -1 pkg/

# ผลลัพธ์จาก icmongolang:
# cryptpass, db, elasticsearch, emailTemplates, helpers,
# http-swagger, httpErrors, influxdb, jwt, kafka, llm,
# logger, mqtt, report, responses, secureRandom, sendEmail,
# transaction, utils, vectordb, websocket
```

**กฎ:** ห้าม import package นอก list นี้ — ถ้าต้องการเพิ่ม ต้องสร้างใน `pkg/` ก่อน

## 3.3 ตรวจสอบ Modules ที่มีอยู่

```bash
# ดู modules ปัจจุบัน (เพื่อหลีกเลี่ยงชื่อซ้ำ)
ls -1 internal/modules/

# ตัวอย่าง icmongolang (33 modules):
# alarm, apimanager, auditlog, auth, batch, control, customer,
# dashboard, document, elasticsearch, email, flowengine,
# fullschedule, i18n, influxdb, iot, items, job, kafka, mqtt,
# notifier, payment, pdpa, purchaseorder, queue, quotation,
# realtime, report, settings, users, vectordata, websocket, wos
```

## 3.4 Environment Variables ที่ใช้บ่อย

```env
# Database
DB_DSN=postgres://user:pass@localhost:5432/db?sslmode=disable
DB_MAX_OPEN=25
DB_MAX_IDLE=10

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_GROUP_ID={module}-group

# JWT
JWT_SECRET=xxx
JWT_TTL=24h

# Module-specific
{MODULE}_RETENTION_DAYS=365
{MODULE}_DEFAULT_CURRENCY=THB
```

---

# บทที่ 4: Naming Conventions

## 4.1 ตาราง Naming ทั้งหมด

| ประเภท | รูปแบบ | ตัวอย่าง |
|---|---|---|
| **Module folder** | lowercase, single word | `payment`, `iot` |
| **Package** | lowercase, single word | `entity`, `valueobject` |
| **Entity (struct)** | PascalCase | `Payment`, `PurchaseOrder` |
| **File** | snake_case | `payment_repository.go` |
| **Table** | `{module}_{plural}` | `payment_transactions` |
| **Column** | snake_case | `user_id`, `created_at` |
| **Primary Key** | `id` (UUID) | – |
| **Foreign Key** | `{entity}_id` | `order_id` |
| **Index** | `idx_{table}_{col}` | `idx_payment_transactions_user_id` |
| **Unique** | `uq_{table}_{col}` | `uq_payment_transactions_order_id` |
| **Migration** | `YYYYMMDD_{module}_{desc}.sql` | `20260401_payment_init.sql` |
| **Kafka Topic** | `{module}.{entity}.{action}` | `payment.txn.created` |
| **Redis Key** | `{module}:{entity}:{id}` | `payment:txn:abc-123` |
| **Use Case** | `{Verb}{Entity}UseCase` | `CreatePaymentUseCase` |
| **Use Case File** | `{verb}_{entity}.go` | `create_payment.go` |
| **Handler Method** | `{HTTPVerb}` | `Create`, `Update`, `List` |
| **Domain Error** | `Err{Reason}` | `ErrInsufficientBalance` |

## 4.2 Folder Structure มาตรฐาน

### 4.2.1 โครงสร้างเต็ม (สำหรับ module ที่ซับซ้อน)

```
internal/modules/{module}/
│
├── domain/                                    # 🏛️ DOMAIN LAYER
│   ├── entity/
│   │   └── {entity}.go                       # Aggregate Root + Entity ย่อย
│   ├── value_object/
│   │   ├── {vo}.go                           # Immutable value objects
│   │   └── status.go                         # Enum-like status
│   ├── repository/
│   │   └── {entity}_repository.go            # Interface เท่านั้น
│   ├── service/
│   │   └── {domain}_service.go               # Stateless domain logic
│   ├── event/
│   │   └── event.go                          # Domain events
│   └── errors/
│       └── errors.go                         # Sentinel errors
│
├── application/                                # 🎯 APPLICATION LAYER
│   ├── {verb}_{entity}.go                    # 1 use case ต่อ 1 ไฟล์
│   ├── dto.go                                # Input/Output DTO
│   └── mappers.go                            # Entity ↔ DTO
│
├── infrastructure/                             # 🔧 INFRASTRUCTURE
│   ├── persistence/
│   │   ├── postgres/
│   │   │   ├── models.go                     # GORM models
│   │   │   └── {entity}_repo_impl.go         # Repository impl
│   │   └── redis/
│   │       └── {entity}_cache.go
│   ├── messaging/
│   │   ├── kafka_producer.go
│   │   └── consumers/
│   │       └── {topic}_consumer.go
│   └── external/
│       └── {service}.go                      # Stripe, SMTP, LLM ฯลฯ
│
├── interfaces/                                 # 🌐 INTERFACE LAYER
│   ├── http/
│   │   ├── {entity}_handler.go
│   │   ├── dto.go                            # HTTP-specific DTO
│   │   └── routes.go
│   └── middleware/
│       └── {specific}.go
│
└── module.go                                   # Composition Root
```

### 4.2.2 โครงสร้างแบบย่อ (สำหรับ module CRUD)

```
internal/modules/{module}/
├── domain/
│   ├── entity/{entity}.go
│   ├── repository/{entity}_repository.go
│   └── errors/errors.go
├── application/
│   └── {verb}_{entity}.go
├── infrastructure/postgres/
│   ├── models.go
│   └── {entity}_repo_impl.go
├── interfaces/http/
│   ├── {entity}_handler.go
│   └── routes.go
└── module.go
```

### 4.2.3 คำแนะนำ: เริ่มจากย่อ → ขยายเมื่อจำเป็น

- **เริ่ม:** ใช้โครงสร้างย่อ
- **เมื่อมี VO:** เพิ่ม `value_object/`
- **เมื่อมี Kafka:** เพิ่ม `messaging/`
- **เมื่อมี Redis:** เพิ่ม `persistence/redis/`

## 4.3 ตัวอย่าง Folder Structure: Payment Module

```
internal/modules/payment/
├── domain/
│   ├── entity/payment.go
│   ├── value_object/
│   │   ├── status.go
│   │   └── method.go
│   ├── repository/payment_repository.go
│   └── errors/errors.go
├── application/
│   ├── create_payment.go
│   ├── refund_payment.go
│   ├── get_payment.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/
│   │   ├── postgres/
│   │   │   ├── models.go
│   │   │   └── payment_repo_impl.go
│   │   └── redis/
│   │       └── payment_cache.go
│   └── messaging/
│       └── kafka_producer.go
├── interfaces/http/
│   ├── payment_handler.go
│   └── routes.go
└── module.go
```

---

# บทที่ 5: ออกแบบ Domain Model

> **ขั้นตอนนี้สำคัญที่สุด** — ใช้เวลา 30-50% ของโปรเจกต์ทั้งหมด

## 5.1 คำถาม 7 ข้อที่ต้องตอบ

1. **Aggregate Root** คืออะไร? (เช่น `Payment`)
2. **Entities ย่อย** ใน aggregate มีอะไร? (เช่น `Transaction`)
3. **Value Objects** อะไรบ้าง? (เช่น `Money`, `PaymentStatus`)
4. **State Transitions** มีอะไร? (เช่น PENDING → SUCCESS)
5. **Invariants** (กฎที่ต้องเป็นจริงเสมอ) คืออะไร?
6. **Domain Events** ที่เกิดขึ้นคืออะไร?
7. **Use Cases** ที่ต้องรองรับคืออะไร?

## 5.2 ตัวอย่างการวิเคราะห์: Payment Module

### 5.2.1 Business Requirements

```
บริษัทต้องการ:
1. ลูกค้าสร้าง payment สำหรับ order
2. ระบบประมวลผลผ่าน provider (Stripe, Omise)
3. ถ้าสำเร็จ → mark SUCCESS, publish event
4. ถ้าล้มเหลว → mark FAILED, notify
5. ถ้าลูกค้าขอ → refund ได้ (ครั้งเดียว)
6. ทุกอย่างต้องมี audit trail
```

### 5.2.2 Aggregate Design

```
Aggregate Root: Payment
├── id: UUID
├── user_id: UUID
├── order_id: UUID
├── amount: Money (VO)
├── method: PaymentMethod (VO)
├── status: PaymentStatus (VO)
├── provider_txn_id: string (nullable)
├── refunded_at: *time.Time
├── created_at, updated_at: time.Time
```

### 5.2.3 Value Objects

**PaymentStatus:**
```
States: PENDING | PROCESSING | SUCCESS | FAILED | REFUNDED
Transitions:
  PENDING     → PROCESSING, FAILED
  PROCESSING  → SUCCESS, FAILED
  SUCCESS     → REFUNDED
  FAILED      → (terminal)
  REFUNDED    → (terminal)
```

**PaymentMethod:**
```
Values: CARD | QR_PROMPT_PAY | BANK_TRANSFER
```

### 5.2.4 Invariants

1. `amount > 0` เสมอ
2. `order_id` ไม่ซ้ำ (1 order = 1 payment)
3. refund ได้เฉพาะ `SUCCESS`
4. refund ได้ครั้งเดียว
5. state transition ต้องถูกต้อง

### 5.2.5 Domain Events

- `PaymentCreated` — เมื่อสร้าง payment
- `PaymentSucceeded` — เมื่อประมวลผลสำเร็จ
- `PaymentFailed` — เมื่อล้มเหลว
- `PaymentRefunded` — เมื่อคืนเงิน

### 5.2.6 Use Cases

| Use Case | Input | Output | Side Effects |
|---|---|---|---|
| `CreatePayment` | user, order, amount, method | payment ID, status | Save, Publish `payment.created` |
| `ProcessPayment` | payment ID, provider | – | Update, Publish event |
| `RefundPayment` | payment ID, reason | – | Update, Publish `payment.refunded` |
| `GetPayment` | payment ID | payment details | – |
| `ListPayments` | user ID, pagination | payments | – |

## 5.3 แผนภาพ State Machine

```
                    StartProcessing()
       ┌───────────────────────────────┐
       │                               ▼
   ┌─────────┐   StartProcessing   ┌────────────┐
   │ PENDING │────────────────────▶│ PROCESSING │
   └────┬────┘                     └──────┬─────┘
        │                                 │
        │ Fail()                          │ MarkSuccess()
        │                                 │
        ▼                                 ▼
   ┌────────┐                       ┌─────────┐
   │ FAILED │                       │ SUCCESS │
   └────────┘                       └────┬────┘
    (terminal)                          │
                                        │ Refund()
                                        ▼
                                   ┌──────────┐
                                   │ REFUNDED │
                                   └──────────┘
                                    (terminal)
```

## 5.4 Ubiquitous Language Dictionary

ก่อนเขียนโค้ด ให้เขียน dictionary นี้ก่อน:

| คำ | ความหมาย | ใช้ในโค้ด |
|---|---|---|
| Payment | รายการชำระเงิน 1 รายการ | `Payment` (entity) |
| Transaction | การพยายามชำระ 1 ครั้ง | `Transaction` |
| Refund | การคืนเงิน | `Refund()` |
| Capture | การหักเงินจริง | `Capture()` |
| Void | การยกเลิกก่อน capture | `Void()` |
| Provider | ผู้ให้บริการ payment | `PaymentProvider` (interface) |

---

# บทที่ 6: สร้าง Domain Layer

## 6.1 สร้าง Value Objects

### 6.1.1 Status Value Object

```go
// internal/modules/payment/domain/value_object/status.go
package valueobject

// PaymentStatus represents the state of a payment
// PaymentStatus แทนสถานะของ payment
type PaymentStatus string

const (
    PaymentStatusPending    PaymentStatus = "PENDING"
    PaymentStatusProcessing PaymentStatus = "PROCESSING"
    PaymentStatusSuccess    PaymentStatus = "SUCCESS"
    PaymentStatusFailed     PaymentStatus = "FAILED"
    PaymentStatusRefunded   PaymentStatus = "REFUNDED"
)

// IsValid checks if the status is a known value
// IsValid ตรวจสอบว่าสถานะเป็นค่าที่รู้จัก
func (s PaymentStatus) IsValid() bool {
    switch s {
    case PaymentStatusPending, PaymentStatusProcessing,
         PaymentStatusSuccess, PaymentStatusFailed, PaymentStatusRefunded:
        return true
    }
    return false
}

// String returns the string representation
// String คืนค่า string representation
func (s PaymentStatus) String() string {
    return string(s)
}

// IsTerminal checks if this status is final (no further transition)
// IsTerminal ตรวจสอบว่าสถานะนี้สิ้นสุดหรือไม่
func (s PaymentStatus) IsTerminal() bool {
    return s == PaymentStatusFailed || s == PaymentStatusRefunded
}

// CanTransitionTo checks if the state can transition to next
// CanTransitionTo ตรวจสอบว่าสามารถเปลี่ยนไปสถานะถัดไปได้หรือไม่
func (s PaymentStatus) CanTransitionTo(next PaymentStatus) bool {
    transitions := map[PaymentStatus][]PaymentStatus{
        PaymentStatusPending:    {PaymentStatusProcessing, PaymentStatusFailed},
        PaymentStatusProcessing: {PaymentStatusSuccess, PaymentStatusFailed},
        PaymentStatusSuccess:    {PaymentStatusRefunded},
        PaymentStatusFailed:     {},
        PaymentStatusRefunded:   {},
    }

    allowed, ok := transitions[s]
    if !ok {
        return false
    }
    for _, a := range allowed {
        if a == next {
            return true
        }
    }
    return false
}
```

### 6.1.2 Method Value Object

```go
// internal/modules/payment/domain/value_object/method.go
package valueobject

// PaymentMethod represents the payment method
// PaymentMethod แทนวิธีชำระเงิน
type PaymentMethod string

const (
    PaymentMethodCard         PaymentMethod = "CARD"
    PaymentMethodQRPromptPay  PaymentMethod = "QR_PROMPT_PAY"
    PaymentMethodBankTransfer PaymentMethod = "BANK_TRANSFER"
)

// IsValid checks if the method is known
// IsValid ตรวจสอบว่าวิธีเป็นค่าที่รู้จัก
func (m PaymentMethod) IsValid() bool {
    switch m {
    case PaymentMethodCard, PaymentMethodQRPromptPay, PaymentMethodBankTransfer:
        return true
    }
    return false
}

// String returns the string representation
func (m PaymentMethod) String() string {
    return string(m)
}
```

### 6.1.3 Money Value Object (ถ้าต้องการ)

```go
// internal/modules/payment/domain/value_object/money.go
package valueobject

import (
    "errors"
    "github.com/shopspring/decimal"
)

// Money represents an amount with currency
// Money แทนจำนวนเงินพร้อมสกุล
type Money struct {
    Amount   decimal.Decimal
    Currency string
}

// NewMoney creates a new Money with validation
// NewMoney สร้าง Money ใหม่พร้อม validation
func NewMoney(amount decimal.Decimal, currency string) (Money, error) {
    if amount.LessThan(decimal.Zero) {
        return Money{}, errors.New("amount cannot be negative")
    }
    if currency == "" {
        return Money{}, errors.New("currency required")
    }
    return Money{Amount: amount, Currency: currency}, nil
}

// Add adds two Money (must be same currency)
// Add บวก Money สองตัว (สกุลต้องเดียวกัน)
func (m Money) Add(other Money) (Money, error) {
    if m.Currency != other.Currency {
        return Money{}, errors.New("currency mismatch")
    }
    return Money{Amount: m.Amount.Add(other.Amount), Currency: m.Currency}, nil
}

// IsZero checks if amount is zero
func (m Money) IsZero() bool {
    return m.Amount.IsZero()
}

// String returns formatted string
func (m Money) String() string {
    return m.Amount.StringFixed(2) + " " + m.Currency
}
```

## 6.2 สร้าง Domain Errors

```go
// internal/modules/payment/domain/errors/errors.go
package domainerrors

import "errors"

// Sentinel errors for the payment domain
// Sentinel errors สำหรับ payment domain
var (
    // Not found
    ErrPaymentNotFound = errors.New("payment not found")

    // Validation
    ErrInvalidUserID        = errors.New("invalid user id")
    ErrInvalidOrderID       = errors.New("invalid order id")
    ErrInvalidAmount        = errors.New("amount must be positive")
    ErrInvalidCurrency      = errors.New("invalid currency")
    ErrInvalidPaymentMethod = errors.New("invalid payment method")

    // Business rules
    ErrPaymentAlreadyExists    = errors.New("payment already exists for this order")
    ErrInvalidStatusTransition = errors.New("invalid status transition")
    ErrCannotRefund            = errors.New("payment cannot be refunded")
    ErrAlreadyRefunded         = errors.New("payment already refunded")

    // Authorization
    ErrUnauthorized = errors.New("unauthorized")
)
```

**หลักการ:**
- ใช้ `errors.New` (sentinel error)
- ไม่ include context (เช่น ID) — ให้ caller ห่อ
- ห้ามใช้ `fmt.Errorf` ใน Domain errors declaration
- ชื่อขึ้นต้น `Err`

## 6.3 สร้าง Entity (Aggregate Root)

```go
// internal/modules/payment/domain/entity/payment.go
package entity

import (
    "time"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"

    domainerrors "icmongolang/internal/modules/payment/domain/errors"
    valueobject "icmongolang/internal/modules/payment/domain/value_object"
)

// Payment is the aggregate root for payment transactions
// Payment เป็น aggregate root สำหรับรายการชำระเงิน
type Payment struct {
    // Identity
    ID uuid.UUID `json:"id"`

    // References (ใช้ ID ไม่ใช้ pointer ข้าม aggregate)
    UserID  uuid.UUID `json:"user_id"`
    OrderID uuid.UUID `json:"order_id"`

    // Value Objects
    Amount   decimal.Decimal               `json:"amount"`
    Currency string                        `json:"currency"`
    Method   valueobject.PaymentMethod     `json:"method"`
    Status   valueobject.PaymentStatus     `json:"status"`

    // Provider info
    ProviderTxnID string `json:"provider_txn_id,omitempty"`

    // Timestamps
    RefundedAt *time.Time `json:"refunded_at,omitempty"`
    CreatedAt  time.Time  `json:"created_at"`
    UpdatedAt  time.Time  `json:"updated_at"`
}

// NewPayment creates a new Payment aggregate
// NewPayment สร้าง Payment aggregate ใหม่
//
// Invariants enforced:
//   - userID, orderID must not be zero
//   - amount must be positive
//   - method must be valid
func NewPayment(
    userID, orderID uuid.UUID,
    amount decimal.Decimal,
    currency string,
    method valueobject.PaymentMethod,
) (*Payment, error) {
    // ตรวจสอบ user ID
    // Validate user ID
    if userID == uuid.Nil {
        return nil, domainerrors.ErrInvalidUserID
    }
    // ตรวจสอบ order ID
    // Validate order ID
    if orderID == uuid.Nil {
        return nil, domainerrors.ErrInvalidOrderID
    }
    // ตรวจสอบจำนวนเงิน
    // Validate amount
    if amount.LessThanOrEqual(decimal.Zero) {
        return nil, domainerrors.ErrInvalidAmount
    }
    // ตรวจสอบสกุลเงิน
    // Validate currency
    if len(currency) != 3 {
        return nil, domainerrors.ErrInvalidCurrency
    }
    // ตรวจสอบวิธีชำระ
    // Validate payment method
    if !method.IsValid() {
        return nil, domainerrors.ErrInvalidPaymentMethod
    }

    now := time.Now()
    return &Payment{
        ID:        uuid.New(),
        UserID:    userID,
        OrderID:   orderID,
        Amount:    amount,
        Currency:  currency,
        Method:    method,
        Status:    valueobject.PaymentStatusPending,
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}

// StartProcessing transitions to PROCESSING state
// StartProcessing เปลี่ยนไปสถานะ PROCESSING
func (p *Payment) StartProcessing() error {
    if !p.Status.CanTransitionTo(valueobject.PaymentStatusProcessing) {
        return domainerrors.ErrInvalidStatusTransition
    }
    p.Status = valueobject.PaymentStatusProcessing
    p.UpdatedAt = time.Now()
    return nil
}

// MarkSuccess marks the payment as SUCCESS with provider txn ID
// MarkSuccess ทำเครื่องหมายว่าชำระสำเร็จพร้อม provider txn ID
func (p *Payment) MarkSuccess(providerTxnID string) error {
    if !p.Status.CanTransitionTo(valueobject.PaymentStatusSuccess) {
        return domainerrors.ErrInvalidStatusTransition
    }
    p.Status = valueobject.PaymentStatusSuccess
    p.ProviderTxnID = providerTxnID
    p.UpdatedAt = time.Now()
    return nil
}

// MarkFailed marks the payment as FAILED
// MarkFailed ทำเครื่องหมายว่าล้มเหลว
func (p *Payment) MarkFailed() error {
    if !p.Status.CanTransitionTo(valueobject.PaymentStatusFailed) {
        return domainerrors.ErrInvalidStatusTransition
    }
    p.Status = valueobject.PaymentStatusFailed
    p.UpdatedAt = time.Now()
    return nil
}

// Refund refunds a SUCCESS payment
// Refund คืนเงิน payment ที่สำเร็จ
//
// Invariants enforced:
//   - must be in SUCCESS state
//   - must not have been refunded before
func (p *Payment) Refund() error {
    if !p.Status.CanTransitionTo(valueobject.PaymentStatusRefunded) {
        return domainerrors.ErrCannotRefund
    }
    if p.RefundedAt != nil {
        return domainerrors.ErrAlreadyRefunded
    }
    now := time.Now()
    p.Status = valueobject.PaymentStatusRefunded
    p.RefundedAt = &now
    p.UpdatedAt = now
    return nil
}

// IsRefundable checks if the payment can be refunded
// IsRefundable ตรวจสอบว่าสามารถคืนเงินได้หรือไม่
func (p *Payment) IsRefundable() bool {
    return p.Status == valueobject.PaymentStatusSuccess && p.RefundedAt == nil
}

// IsTerminal checks if the payment is in a terminal state
// IsTerminal ตรวจสอบว่า payment อยู่ในสถานะสิ้นสุดหรือไม่
func (p *Payment) IsTerminal() bool {
    return p.Status.IsTerminal()
}
```

## 6.4 สร้าง Repository Interface

```go
// internal/modules/payment/domain/repository/payment_repository.go
package repository

import (
    "context"

    "github.com/google/uuid"

    "icmongolang/internal/modules/payment/domain/entity"
)

// PaymentRepository defines the contract for payment persistence
// PaymentRepository กำหนดสัญญาสำหรับการจัดเก็บ payment
//
// NOTE: Interface only — implementation อยู่ใน Infrastructure layer
type PaymentRepository interface {
    // Save creates a new payment
    // Save สร้าง payment ใหม่
    Save(ctx context.Context, p *entity.Payment) error

    // FindByID finds a payment by ID
    // FindByID ค้นหา payment ด้วย ID
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error)

    // FindByOrderID finds a payment by order ID
    // FindByOrderID ค้นหา payment ด้วย order ID
    FindByOrderID(ctx context.Context, orderID uuid.UUID) (*entity.Payment, error)

    // FindByUserID finds all payments of a user with pagination
    // FindByUserID ค้นหา payment ทั้งหมดของผู้ใช้พร้อม pagination
    FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Payment, error)

    // Update updates an existing payment
    // Update อัปเดต payment ที่มีอยู่
    Update(ctx context.Context, p *entity.Payment) error

    // Delete soft-deletes a payment
    // Delete ลบ payment (soft delete)
    Delete(ctx context.Context, id uuid.UUID) error
}
```

## 6.5 สร้าง Domain Event (Optional)

```go
// internal/modules/payment/domain/event/events.go
package event

import (
    "time"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)

// BaseEvent is the base for all domain events
// BaseEvent เป็น base ของ domain events ทั้งหมด
type BaseEvent struct {
    EventID    uuid.UUID `json:"event_id"`
    OccurredAt time.Time `json:"occurred_at"`
}

// NewBaseEvent creates a new BaseEvent
func NewBaseEvent() BaseEvent {
    return BaseEvent{
        EventID:    uuid.New(),
        OccurredAt: time.Now(),
    }
}

// PaymentCreated is emitted when a payment is created
// PaymentCreated ถูก emit เมื่อสร้าง payment
type PaymentCreated struct {
    BaseEvent
    PaymentID uuid.UUID       `json:"payment_id"`
    UserID    uuid.UUID       `json:"user_id"`
    OrderID   uuid.UUID       `json:"order_id"`
    Amount    decimal.Decimal `json:"amount"`
    Currency  string          `json:"currency"`
}

// PaymentSucceeded is emitted when a payment succeeds
// PaymentSucceeded ถูก emit เมื่อ payment สำเร็จ
type PaymentSucceeded struct {
    BaseEvent
    PaymentID     uuid.UUID `json:"payment_id"`
    ProviderTxnID string    `json:"provider_txn_id"`
}

// PaymentFailed is emitted when a payment fails
// PaymentFailed ถูก emit เมื่อ payment ล้มเหลว
type PaymentFailed struct {
    BaseEvent
    PaymentID uuid.UUID `json:"payment_id"`
    Reason    string    `json:"reason"`
}

// PaymentRefunded is emitted when a payment is refunded
// PaymentRefunded ถูก emit เมื่อ payment ถูกคืนเงิน
type PaymentRefunded struct {
    BaseEvent
    PaymentID uuid.UUID `json:"payment_id"`
    RefundedBy uuid.UUID `json:"refunded_by,omitempty"`
}
```

## 6.6 เขียน Unit Test สำหรับ Domain

**ก่อนไปขั้นถัดไป ต้อง test domain ให้ผ่าน 100%**

```go
// internal/modules/payment/domain/entity/payment_test.go
package entity_test

import (
    "testing"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "icmongolang/internal/modules/payment/domain/entity"
    domainerrors "icmongolang/internal/modules/payment/domain/errors"
    valueobject "icmongolang/internal/modules/payment/domain/value_object"
)

func TestNewPayment_Success(t *testing.T) {
    userID := uuid.New()
    orderID := uuid.New()

    p, err := entity.NewPayment(userID, orderID, decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)

    require.NoError(t, err)
    assert.NotEqual(t, uuid.Nil, p.ID)
    assert.Equal(t, userID, p.UserID)
    assert.Equal(t, orderID, p.OrderID)
    assert.Equal(t, valueobject.PaymentStatusPending, p.Status)
    assert.Equal(t, "100", p.Amount.String())
}

func TestNewPayment_InvalidInputs(t *testing.T) {
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
            name:     "zero user ID",
            userID:   uuid.Nil,
            orderID:  uuid.New(),
            amount:   decimal.NewFromInt(100),
            currency: "THB",
            method:   valueobject.PaymentMethodCard,
            wantErr:  domainerrors.ErrInvalidUserID,
        },
        {
            name:     "zero order ID",
            userID:   uuid.New(),
            orderID:  uuid.Nil,
            amount:   decimal.NewFromInt(100),
            currency: "THB",
            method:   valueobject.PaymentMethodCard,
            wantErr:  domainerrors.ErrInvalidOrderID,
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
            name:     "zero amount",
            userID:   uuid.New(),
            orderID:  uuid.New(),
            amount:   decimal.Zero,
            currency: "THB",
            method:   valueobject.PaymentMethodCard,
            wantErr:  domainerrors.ErrInvalidAmount,
        },
        {
            name:     "invalid currency",
            userID:   uuid.New(),
            orderID:  uuid.New(),
            amount:   decimal.NewFromInt(100),
            currency: "TH",
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
            _, err := entity.NewPayment(tt.userID, tt.orderID, tt.amount, tt.currency, tt.method)
            assert.ErrorIs(t, err, tt.wantErr)
        })
    }
}

func TestPayment_HappyPath(t *testing.T) {
    // PENDING → PROCESSING → SUCCESS
    p, _ := entity.NewPayment(uuid.New(), uuid.New(), decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)

    require.NoError(t, p.StartProcessing())
    assert.Equal(t, valueobject.PaymentStatusProcessing, p.Status)

    require.NoError(t, p.MarkSuccess("txn_123"))
    assert.Equal(t, valueobject.PaymentStatusSuccess, p.Status)
    assert.Equal(t, "txn_123", p.ProviderTxnID)
    assert.True(t, p.IsTerminal())
}

func TestPayment_InvalidTransition(t *testing.T) {
    // ไม่สามารถ PENDING → SUCCESS โดยตรง
    p, _ := entity.NewPayment(uuid.New(), uuid.New(), decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)

    err := p.MarkSuccess("txn_123")
    assert.ErrorIs(t, err, domainerrors.ErrInvalidStatusTransition)
    assert.Equal(t, valueobject.PaymentStatusPending, p.Status)
}

func TestPayment_RefundHappyPath(t *testing.T) {
    p, _ := entity.NewPayment(uuid.New(), uuid.New(), decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
    _ = p.StartProcessing()
    _ = p.MarkSuccess("txn_123")

    require.NoError(t, p.Refund())
    assert.Equal(t, valueobject.PaymentStatusRefunded, p.Status)
    assert.NotNil(t, p.RefundedAt)
}

func TestPayment_CannotRefundTwice(t *testing.T) {
    p, _ := entity.NewPayment(uuid.New(), uuid.New(), decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)
    _ = p.StartProcessing()
    _ = p.MarkSuccess("txn_123")
    _ = p.Refund()

    err := p.Refund()
    assert.ErrorIs(t, err, domainerrors.ErrCannotRefund)
}

func TestPayment_CannotRefundPending(t *testing.T) {
    p, _ := entity.NewPayment(uuid.New(), uuid.New(), decimal.NewFromInt(100), "THB", valueobject.PaymentMethodCard)

    err := p.Refund()
    assert.ErrorIs(t, err, domainerrors.ErrCannotRefund)
}
```

**รัน test:**
```bash
go test ./internal/modules/payment/domain/... -v
# ต้องผ่าน 100%
```

---

# บทที่ 7: สร้าง Application Layer

## 7.1 หลักการของ Use Case

- **1 use case = 1 ไฟล์**
- **1 use case = 1 public method `Execute`**
- **Input/Output เป็น DTO เฉพาะของ use case**
- **ไม่มี SQL, HTTP, framework ในไฟล์**

## 7.2 Use Case: CreatePayment

```go
// internal/modules/payment/application/create_payment.go
package application

import (
    "context"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"

    "icmongolang/internal/modules/payment/domain/entity"
    domainerrors "icmongolang/internal/modules/payment/domain/errors"
    "icmongolang/internal/modules/payment/domain/repository"
    valueobject "icmongolang/internal/modules/payment/domain/value_object"
    "icmongolang/pkg/logger"
)

// CreatePaymentUseCase handles the creation of a new payment
// CreatePaymentUseCase จัดการการสร้าง payment ใหม่
type CreatePaymentUseCase struct {
    repo   repository.PaymentRepository
    logger logger.Logger
}

// NewCreatePaymentUseCase creates a new instance
// NewCreatePaymentUseCase สร้าง instance ใหม่
func NewCreatePaymentUseCase(
    repo repository.PaymentRepository,
    log logger.Logger,
) *CreatePaymentUseCase {
    return &CreatePaymentUseCase{
        repo:   repo,
        logger: log,
    }
}

// CreatePaymentInput is the input DTO
// CreatePaymentInput เป็น input DTO
type CreatePaymentInput struct {
    UserID   uuid.UUID
    OrderID  uuid.UUID
    Amount   decimal.Decimal
    Currency string
    Method   string
}

// CreatePaymentOutput is the output DTO
// CreatePaymentOutput เป็น output DTO
type CreatePaymentOutput struct {
    ID       uuid.UUID `json:"id"`
    Status   string    `json:"status"`
    Amount   string    `json:"amount"`
    Currency string    `json:"currency"`
}

// Execute runs the create payment use case
// Execute รัน use case สร้าง payment
func (uc *CreatePaymentUseCase) Execute(
    ctx context.Context,
    input CreatePaymentInput,
) (*CreatePaymentOutput, error) {
    // Step 1: Check for duplicate (1 order = 1 payment)
    // ขั้นที่ 1: ตรวจสอบซ้ำ (1 order = 1 payment)
    existing, _ := uc.repo.FindByOrderID(ctx, input.OrderID)
    if existing != nil {
        uc.logger.Warn("duplicate payment attempt",
            "order_id", input.OrderID,
            "existing_id", existing.ID,
        )
        return nil, domainerrors.ErrPaymentAlreadyExists
    }

    // Step 2: Create domain entity (validation อยู่ใน entity)
    // ขั้นที่ 2: สร้าง domain entity (validation อยู่ใน entity)
    method := valueobject.PaymentMethod(input.Method)
    payment, err := entity.NewPayment(
        input.UserID,
        input.OrderID,
        input.Amount,
        input.Currency,
        method,
    )
    if err != nil {
        uc.logger.Warn("invalid payment input", "error", err)
        return nil, err
    }

    // Step 3: Persist
    // ขั้นที่ 3: บันทึก
    if err := uc.repo.Save(ctx, payment); err != nil {
        uc.logger.Error("failed to save payment",
            "error", err,
            "order_id", input.OrderID,
        )
        return nil, err
    }

    // Step 4: Log success
    // ขั้นที่ 4: log ความสำเร็จ
    uc.logger.Info("payment created",
        "payment_id", payment.ID,
        "user_id", payment.UserID,
        "amount", payment.Amount.String(),
        "currency", payment.Currency,
    )

    // Step 5: Return output
    // ขั้นที่ 5: คืนค่า output
    return &CreatePaymentOutput{
        ID:       payment.ID,
        Status:   string(payment.Status),
        Amount:   payment.Amount.String(),
        Currency: payment.Currency,
    }, nil
}
```

## 7.3 Use Case: RefundPayment

```go
// internal/modules/payment/application/refund_payment.go
package application

import (
    "context"

    "github.com/google/uuid"

    domainerrors "icmongolang/internal/modules/payment/domain/errors"
    "icmongolang/internal/modules/payment/domain/repository"
    "icmongolang/pkg/logger"
)

// RefundPaymentUseCase handles refunding a payment
// RefundPaymentUseCase จัดการการคืนเงิน
type RefundPaymentUseCase struct {
    repo   repository.PaymentRepository
    logger logger.Logger
}

// NewRefundPaymentUseCase creates a new instance
func NewRefundPaymentUseCase(
    repo repository.PaymentRepository,
    log logger.Logger,
) *RefundPaymentUseCase {
    return &RefundPaymentUseCase{
        repo:   repo,
        logger: log,
    }
}

// RefundPaymentInput is the input DTO
type RefundPaymentInput struct {
    PaymentID uuid.UUID
    UserID    uuid.UUID // สำหรับ authorization check
    Reason    string
}

// Execute runs the refund use case
// Execute รัน use case คืนเงิน
func (uc *RefundPaymentUseCase) Execute(
    ctx context.Context,
    input RefundPaymentInput,
) error {
    // Step 1: Load aggregate
    // ขั้นที่ 1: โหลด aggregate
    payment, err := uc.repo.FindByID(ctx, input.PaymentID)
    if err != nil {
        return err
    }

    // Step 2: Authorization check
    // ขั้นที่ 2: ตรวจสอบสิทธิ์
    if payment.UserID != input.UserID {
        uc.logger.Warn("unauthorized refund attempt",
            "payment_id", input.PaymentID,
            "requester_id", input.UserID,
        )
        return domainerrors.ErrUnauthorized
    }

    // Step 3: Call domain behavior (validation อยู่ใน domain)
    // ขั้นที่ 3: เรียก domain behavior (validation อยู่ใน domain)
    if err := payment.Refund(); err != nil {
        return err
    }

    // Step 4: Persist
    // ขั้นที่ 4: บันทึก
    if err := uc.repo.Update(ctx, payment); err != nil {
        uc.logger.Error("failed to update payment", "error", err)
        return err
    }

    // Step 5: Log
    uc.logger.Info("payment refunded",
        "payment_id", payment.ID,
        "reason", input.Reason,
    )

    return nil
}
```

## 7.4 หลักการตั้งชื่อ Use Case

| Verb | ใช้เมื่อ | ตัวอย่าง |
|---|---|---|
| `Create` | สร้างใหม่ | `CreatePaymentUseCase` |
| `Update` | แก้ไขที่มีอยู่ | `UpdatePaymentUseCase` |
| `Delete` | ลบ | `DeletePaymentUseCase` |
| `Get` | ดึง 1 รายการ | `GetPaymentUseCase` |
| `List` | ดึงหลายรายการ | `ListPaymentsUseCase` |
| `Process` | ประมวลผล (async) | `ProcessPaymentUseCase` |
| `Refund` | verb เฉพาะ domain | `RefundPaymentUseCase` |
| `Cancel` | ยกเลิก | `CancelOrderUseCase` |
| `Approve` | อนุมัติ | `ApproveQuotationUseCase` |

## 7.5 Use Case Test (with Mock)

```go
// internal/modules/payment/application/create_payment_test.go
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

func TestCreatePaymentUseCase_Success(t *testing.T) {
    // Arrange
    repo := new(mocks.PaymentRepositoryMock)
    log := logger.NewNoop()
    uc := application.NewCreatePaymentUseCase(repo, log)

    userID := uuid.New()
    orderID := uuid.New()

    repo.On("FindByOrderID", mock.Anything, orderID).
        Return(nil, domainerrors.ErrPaymentNotFound)
    repo.On("Save", mock.Anything, mock.AnythingOfType("*entity.Payment")).
        Return(nil)

    // Act
    output, err := uc.Execute(context.Background(), application.CreatePaymentInput{
        UserID:   userID,
        OrderID:  orderID,
        Amount:   decimal.NewFromInt(100),
        Currency: "THB",
        Method:   "CARD",
    })

    // Assert
    require.NoError(t, err)
    assert.NotEqual(t, uuid.Nil, output.ID)
    assert.Equal(t, "PENDING", output.Status)
    assert.Equal(t, "100", output.Amount)
    assert.Equal(t, "THB", output.Currency)
    repo.AssertExpectations(t)
}

func TestCreatePaymentUseCase_DuplicateOrder(t *testing.T) {
    repo := new(mocks.PaymentRepositoryMock)
    log := logger.NewNoop()
    uc := application.NewCreatePaymentUseCase(repo, log)

    orderID := uuid.New()
    existing, _ := entity.NewPayment(uuid.New(), orderID, decimal.NewFromInt(50), "THB", "CARD")

    repo.On("FindByOrderID", mock.Anything, orderID).Return(existing, nil)

    output, err := uc.Execute(context.Background(), application.CreatePaymentInput{
        UserID:   uuid.New(),
        OrderID:  orderID,
        Amount:   decimal.NewFromInt(100),
        Currency: "THB",
        Method:   "CARD",
    })

    assert.ErrorIs(t, err, domainerrors.ErrPaymentAlreadyExists)
    assert.Nil(t, output)
    repo.AssertExpectations(t)
}
```

### 7.5.1 สร้าง Mock ด้วย testify

```go
// internal/modules/payment/application/mocks/payment_repository_mock.go
package mocks

import (
    "context"

    "github.com/google/uuid"
    "github.com/stretchr/testify/mock"

    "icmongolang/internal/modules/payment/domain/entity"
)

// PaymentRepositoryMock is a mock for PaymentRepository
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

---

# บทที่ 8: สร้าง Infrastructure Layer

## 8.1 GORM Model

```go
// internal/modules/payment/infrastructure/persistence/postgres/models.go
package postgres

import (
    "time"

    "github.com/google/uuid"
)

// PaymentModel is the GORM model for payment_transactions table
// PaymentModel เป็น GORM model สำหรับตาราง payment_transactions
//
// NOTE: Table name MUST have module prefix
type PaymentModel struct {
    ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    UserID        uuid.UUID  `gorm:"type:uuid;index;not null"`
    OrderID       uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null"`
    Amount        string     `gorm:"type:numeric(15,2);not null"`
    Currency      string     `gorm:"type:varchar(3);not null;default:'THB'"`
    Method        string     `gorm:"type:varchar(30);not null"`
    Status        string     `gorm:"type:varchar(20);not null;index"`
    ProviderTxnID string     `gorm:"type:varchar(100)"`
    RefundedAt    *time.Time
    CreatedAt     time.Time  `gorm:"default:now()"`
    UpdatedAt     time.Time  `gorm:"default:now()"`
}

// TableName returns the table name with module prefix
// TableName คืนชื่อตารางพร้อม prefix ของ module
func (PaymentModel) TableName() string {
    return "payment_transactions"
}
```

## 8.2 Repository Implementation

```go
// internal/modules/payment/infrastructure/persistence/postgres/payment_repo_impl.go
package postgres

import (
    "context"
    "errors"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
    "gorm.io/gorm"

    "icmongolang/internal/modules/payment/domain/entity"
    domainerrors "icmongolang/internal/modules/payment/domain/errors"
    valueobject "icmongolang/internal/modules/payment/domain/value_object"
)

// paymentRepoImpl implements PaymentRepository using GORM
// paymentRepoImpl implement PaymentRepository ด้วย GORM
type paymentRepoImpl struct {
    db *gorm.DB
}

// NewPaymentRepository creates a new repository instance
// NewPaymentRepository สร้าง instance ของ repository ใหม่
func NewPaymentRepository(db *gorm.DB) *paymentRepoImpl {
    return &paymentRepoImpl{db: db}
}

// Save creates a new payment
// Save สร้าง payment ใหม่
func (r *paymentRepoImpl) Save(ctx context.Context, p *entity.Payment) error {
    m := r.toModel(p)
    if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
        // ตรวจ unique violation
        if errors.Is(err, gorm.ErrDuplicatedKey) {
            return domainerrors.ErrPaymentAlreadyExists
        }
        return err
    }
    return nil
}

// FindByID finds a payment by ID
// FindByID ค้นหา payment ด้วย ID
func (r *paymentRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
    var m PaymentModel
    err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrPaymentNotFound
    }
    if err != nil {
        return nil, err
    }
    return r.toEntity(&m), nil
}

// FindByOrderID finds a payment by order ID
// FindByOrderID ค้นหา payment ด้วย order ID
func (r *paymentRepoImpl) FindByOrderID(ctx context.Context, orderID uuid.UUID) (*entity.Payment, error) {
    var m PaymentModel
    err := r.db.WithContext(ctx).First(&m, "order_id = ?", orderID).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrPaymentNotFound
    }
    if err != nil {
        return nil, err
    }
    return r.toEntity(&m), nil
}

// FindByUserID finds all payments of a user
// FindByUserID ค้นหา payment ทั้งหมดของผู้ใช้
func (r *paymentRepoImpl) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Payment, error) {
    var models []PaymentModel
    err := r.db.WithContext(ctx).
        Where("user_id = ?", userID).
        Order("created_at DESC").
        Limit(limit).
        Offset(offset).
        Find(&models).Error
    if err != nil {
        return nil, err
    }

    payments := make([]*entity.Payment, len(models))
    for i := range models {
        payments[i] = r.toEntity(&models[i])
    }
    return payments, nil
}

// Update updates an existing payment
// Update อัปเดต payment ที่มีอยู่
func (r *paymentRepoImpl) Update(ctx context.Context, p *entity.Payment) error {
    m := r.toModel(p)
    result := r.db.WithContext(ctx).Save(m)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return domainerrors.ErrPaymentNotFound
    }
    return nil
}

// Delete soft-deletes a payment
// Delete ลบ payment (soft delete)
func (r *paymentRepoImpl) Delete(ctx context.Context, id uuid.UUID) error {
    result := r.db.WithContext(ctx).Delete(&PaymentModel{}, "id = ?", id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return domainerrors.ErrPaymentNotFound
    }
    return nil
}

// toModel converts domain entity to GORM model
// toModel แปลง domain entity เป็น GORM model
func (r *paymentRepoImpl) toModel(e *entity.Payment) *PaymentModel {
    return &PaymentModel{
        ID:            e.ID,
        UserID:        e.UserID,
        OrderID:       e.OrderID,
        Amount:        e.Amount.String(),
        Currency:      e.Currency,
        Method:        string(e.Method),
        Status:        string(e.Status),
        ProviderTxnID: e.ProviderTxnID,
        RefundedAt:    e.RefundedAt,
        CreatedAt:     e.CreatedAt,
        UpdatedAt:     e.UpdatedAt,
    }
}

// toEntity converts GORM model to domain entity
// toEntity แปลง GORM model เป็น domain entity
func (r *paymentRepoImpl) toEntity(m *PaymentModel) *entity.Payment {
    amount, _ := decimal.NewFromString(m.Amount)
    return &entity.Payment{
        ID:            m.ID,
        UserID:        m.UserID,
        OrderID:       m.OrderID,
        Amount:        amount,
        Currency:      m.Currency,
        Method:        valueobject.PaymentMethod(m.Method),
        Status:        valueobject.PaymentStatus(m.Status),
        ProviderTxnID: m.ProviderTxnID,
        RefundedAt:    m.RefundedAt,
        CreatedAt:     m.CreatedAt,
        UpdatedAt:     m.UpdatedAt,
    }
}
```

## 8.3 Redis Cache (Optional)

```go
// internal/modules/payment/infrastructure/persistence/redis/payment_cache.go
package redis

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/redis/go-redis/v9"

    "icmongolang/internal/modules/payment/domain/entity"
)

// PaymentCache caches payment data in Redis
// PaymentCache cache ข้อมูล payment ใน Redis
type PaymentCache struct {
    client *redis.Client
    ttl    time.Duration
}

// NewPaymentCache creates a new cache instance
func NewPaymentCache(client *redis.Client, ttl time.Duration) *PaymentCache {
    if ttl == 0 {
        ttl = 10 * time.Minute
    }
    return &PaymentCache{client: client, ttl: ttl}
}

// key generates a Redis key with namespace
// key สร้าง Redis key พร้อม namespace
func (c *PaymentCache) key(id uuid.UUID) string {
    return fmt.Sprintf("payment:txn:%s", id.String())
}

// Set caches a payment
func (c *PaymentCache) Set(ctx context.Context, p *entity.Payment) error {
    data, err := json.Marshal(p)
    if err != nil {
        return err
    }
    return c.client.Set(ctx, c.key(p.ID), data, c.ttl).Err()
}

// Get retrieves a payment from cache
func (c *PaymentCache) Get(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
    data, err := c.client.Get(ctx, c.key(id)).Bytes()
    if err != nil {
        return nil, err
    }
    var p entity.Payment
    if err := json.Unmarshal(data, &p); err != nil {
        return nil, err
    }
    return &p, nil
}

// Delete removes a payment from cache
func (c *PaymentCache) Delete(ctx context.Context, id uuid.UUID) error {
    return c.client.Del(ctx, c.key(id)).Err()
}
```

## 8.4 Kafka Producer (Optional)

```go
// internal/modules/payment/infrastructure/messaging/kafka_producer.go
package messaging

import (
    "context"
    "encoding/json"
    "time"

    "github.com/IBM/sarama"
)

// Producer defines the interface for publishing events
// Producer กำหนด interface สำหรับ publish events
type Producer interface {
    Publish(ctx context.Context, topic string, key string, payload interface{}) error
}

type kafkaProducer struct {
    p sarama.SyncProducer
}

// NewKafkaProducer creates a new Kafka producer
// NewKafkaProducer สร้าง Kafka producer ใหม่
func NewKafkaProducer(brokers []string) (Producer, error) {
    cfg := sarama.NewConfig()
    cfg.Producer.Return.Successes = true
    cfg.Producer.RequiredAcks = sarama.WaitForAll
    cfg.Producer.Retry.Max = 5
    cfg.Producer.Compression = sarama.CompressionSnappy

    p, err := sarama.NewSyncProducer(brokers, cfg)
    if err != nil {
        return nil, err
    }
    return &kafkaProducer{p: p}, nil
}

// Publish sends a message to Kafka
// Publish ส่ง message ไปยัง Kafka
func (k *kafkaProducer) Publish(ctx context.Context, topic, key string, payload interface{}) error {
    data, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    _, _, err = k.p.SendMessage(&sarama.ProducerMessage{
        Topic:     topic,
        Key:       sarama.StringEncoder(key),
        Value:     sarama.ByteEncoder(data),
        Timestamp: time.Now(),
    })
    return err
}
```

---

# บทที่ 9: สร้าง Interface Layer

## 9.1 HTTP DTO

```go
// internal/modules/payment/interfaces/http/dto.go
package http

// CreatePaymentRequest is the HTTP request body for creating a payment
// CreatePaymentRequest เป็น HTTP request body สำหรับสร้าง payment
type CreatePaymentRequest struct {
    OrderID  string `json:"order_id" validate:"required,uuid"`
    Amount   string `json:"amount" validate:"required"`
    Currency string `json:"currency" validate:"required,len=3"`
    Method   string `json:"method" validate:"required,oneof=CARD QR_PROMPT_PAY BANK_TRANSFER"`
}

// RefundPaymentRequest is the HTTP request body for refunding
// RefundPaymentRequest เป็น HTTP request body สำหรับคืนเงิน
type RefundPaymentRequest struct {
    Reason string `json:"reason" validate:"required,min=3,max=500"`
}

// PaymentResponse is the HTTP response for a payment
// PaymentResponse เป็น HTTP response สำหรับ payment
type PaymentResponse struct {
    ID       string `json:"id"`
    OrderID  string `json:"order_id"`
    Amount   string `json:"amount"`
    Currency string `json:"currency"`
    Method   string `json:"method"`
    Status   string `json:"status"`
}
```

## 9.2 HTTP Handler

```go
// internal/modules/payment/interfaces/http/payment_handler.go
package http

import (
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    "github.com/shopspring/decimal"

    "icmongolang/internal/modules/payment/application"
    "icmongolang/pkg/httputil"
    "icmongolang/pkg/validator"
)

// PaymentHandler handles HTTP requests for payment module
// PaymentHandler จัดการ HTTP request สำหรับโมดูล payment
type PaymentHandler struct {
    createUC  *application.CreatePaymentUseCase
    refundUC  *application.RefundPaymentUseCase
    validator *validator.Validator
}

// NewPaymentHandler creates a new handler
func NewPaymentHandler(
    createUC *application.CreatePaymentUseCase,
    refundUC *application.RefundPaymentUseCase,
    v *validator.Validator,
) *PaymentHandler {
    return &PaymentHandler{
        createUC:  createUC,
        refundUC:  refundUC,
        validator: v,
    }
}

// Create handles POST /api/v1/payments
// Create จัดการ POST /api/v1/payments
func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
    // Step 1: Extract user from context (set by auth middleware)
    // ขั้นที่ 1: ดึง user จาก context (ตั้งโดย auth middleware)
    userID, ok := r.Context().Value("user_id").(uuid.UUID)
    if !ok {
        httputil.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
        return
    }

    // Step 2: Parse request body
    // ขั้นที่ 2: parse request body
    var req CreatePaymentRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
        return
    }

    // Step 3: Validate
    // ขั้นที่ 3: validate
    if err := h.validator.Validate(req); err != nil {
        httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        return
    }

    // Step 4: Parse IDs and amount
    // ขั้นที่ 4: parse ID และจำนวนเงิน
    orderID, err := uuid.Parse(req.OrderID)
    if err != nil {
        httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid order_id"})
        return
    }
    amount, err := decimal.NewFromString(req.Amount)
    if err != nil {
        httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid amount"})
        return
    }

    // Step 5: Execute use case
    // ขั้นที่ 5: รัน use case
    output, err := h.createUC.Execute(r.Context(), application.CreatePaymentInput{
        UserID:   userID,
        OrderID:  orderID,
        Amount:   amount,
        Currency: req.Currency,
        Method:   req.Method,
    })
    if err != nil {
        status := mapErrorToStatus(err)
        httputil.JSON(w, status, map[string]string{"error": err.Error()})
        return
    }

    // Step 6: Return response
    // ขั้นที่ 6: คืนค่า response
    httputil.JSON(w, http.StatusCreated, output)
}

// Refund handles POST /api/v1/payments/{id}/refund
// Refund จัดการ POST /api/v1/payments/{id}/refund
func (h *PaymentHandler) Refund(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value("user_id").(uuid.UUID)
    if !ok {
        httputil.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
        return
    }

    // Parse payment ID from URL
    // Parse payment ID จาก URL
    paymentID, err := uuid.Parse(chi.URLParam(r, "id"))
    if err != nil {
        httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payment id"})
        return
    }

    var req RefundPaymentRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
        return
    }
    if err := h.validator.Validate(req); err != nil {
        httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        return
    }

    err = h.refundUC.Execute(r.Context(), application.RefundPaymentInput{
        PaymentID: paymentID,
        UserID:    userID,
        Reason:    req.Reason,
    })
    if err != nil {
        status := mapErrorToStatus(err)
        httputil.JSON(w, status, map[string]string{"error": err.Error()})
        return
    }

    httputil.JSON(w, http.StatusOK, map[string]string{"message": "refunded"})
}

// mapErrorToStatus maps domain errors to HTTP status codes
// mapErrorToStatus map domain error เป็น HTTP status code
func mapErrorToStatus(err error) int {
    switch {
    case errors.Is(err, domainerrors.ErrPaymentNotFound):
        return http.StatusNotFound
    case errors.Is(err, domainerrors.ErrPaymentAlreadyExists):
        return http.StatusConflict
    case errors.Is(err, domainerrors.ErrUnauthorized):
        return http.StatusForbidden
    case errors.Is(err, domainerrors.ErrInvalidStatusTransition),
         errors.Is(err, domainerrors.ErrCannotRefund),
         errors.Is(err, domainerrors.ErrAlreadyRefunded),
         errors.Is(err, domainerrors.ErrInvalidAmount):
        return http.StatusBadRequest
    default:
        return http.StatusInternalServerError
    }
}
```

## 9.3 Routes

```go
// internal/modules/payment/interfaces/http/routes.go
package http

import (
    "net/http"

    "github.com/go-chi/chi/v5"
)

// Handlers aggregates all handlers of the payment module
// Handlers รวม handler ทั้งหมดของโมดูล payment
type Handlers struct {
    Payment *PaymentHandler
}

// RegisterRoutes registers all routes for the payment module
// RegisterRoutes ลงทะเบียน route ทั้งหมดสำหรับโมดูล payment
func RegisterRoutes(r chi.Router, h *Handlers, authMW func(http.Handler) http.Handler) {
    r.Route("/api/v1/payments", func(r chi.Router) {
        r.Use(authMW)
        r.Post("/", h.Payment.Create)
        r.Post("/{id}/refund", h.Payment.Refund)
    })
}
```

## 9.4 Error Mapping Table

| Domain Error | HTTP Status |
|---|---|
| `ErrPaymentNotFound` | 404 |
| `ErrPaymentAlreadyExists` | 409 |
| `ErrUnauthorized` | 403 |
| `ErrInvalid*` | 400 |
| `ErrCannotRefund` | 400 |
| `ErrAlreadyRefunded` | 400 |
| (default) | 500 |

---

# บทที่ 10: Migration และ Schema

## 10.1 Migration Template

```sql
-- migrations/20260401_payment_init.sql
-- ============================================================
-- Payment Module — Initial Schema
-- Prefix: payment_ (เพื่อแยกจากโมดูลอื่น)
-- ============================================================

CREATE TABLE IF NOT EXISTS payment_transactions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    order_id        UUID NOT NULL,
    amount          NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    currency        VARCHAR(3) NOT NULL DEFAULT 'THB',
    method          VARCHAR(30) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    provider_txn_id VARCHAR(100),
    refunded_at     TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Unique constraint: 1 order = 1 payment
CREATE UNIQUE INDEX IF NOT EXISTS uq_payment_transactions_order_id
    ON payment_transactions(order_id);

-- Index for common queries
CREATE INDEX IF NOT EXISTS idx_payment_transactions_user_id
    ON payment_transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_payment_transactions_status
    ON payment_transactions(status);
CREATE INDEX IF NOT EXISTS idx_payment_transactions_created_at
    ON payment_transactions(created_at DESC);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION payment_update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS payment_transactions_updated_at ON payment_transactions;
CREATE TRIGGER payment_transactions_updated_at
    BEFORE UPDATE ON payment_transactions
    FOR EACH ROW EXECUTE FUNCTION payment_update_updated_at();
```

## 10.2 Rollback Script

```sql
-- migrations/20260401_payment_init.down.sql
DROP TRIGGER IF EXISTS payment_transactions_updated_at ON payment_transactions;
DROP FUNCTION IF EXISTS payment_update_updated_at();
DROP TABLE IF EXISTS payment_transactions;
```

## 10.3 Migration Checklist

- [ ] ตั้งชื่อไฟล์ `YYYYMMDD_{module}_{desc}.sql`
- [ ] ทุกตารางมี prefix `{module}_`
- [ ] ใช้ `CREATE TABLE IF NOT EXISTS`
- [ ] มี index สำหรับ column ที่ query บ่อย (FK, status, created_at)
- [ ] มี CHECK constraint (เช่น `amount > 0`)
- [ ] มี UNIQUE constraint ถ้ามี business rule
- [ ] มี trigger `updated_at` (ถ้าจำเป็น)
- [ ] มี rollback script (`.down.sql`)
- [ ] ไม่ lock table นาน (หลีกเลี่ยง ALTER ที่ REWRITE)
- [ ] Test บน staging ที่มี data ≈ production

## 10.4 Index Strategy

| Query Pattern | Index Type |
|---|---|
| `WHERE user_id = ?` | B-tree บน `user_id` |
| `ORDER BY created_at DESC` | B-tree บน `created_at DESC` |
| `WHERE status = ?` (low cardinality) | B-tree (อาจ partial) |
| `WHERE status = 'PENDING'` (specific) | Partial index |
| Full-text search | GIN + `tsvector` |
| JSONB query | GIN + `jsonb_path_ops` |

## 10.5 Partial Index ตัวอย่าง

```sql
-- Partial index เฉพาะ pending payments
CREATE INDEX IF NOT EXISTS idx_payment_transactions_pending
    ON payment_transactions(created_at DESC)
    WHERE status = 'PENDING';
```

---

# บทที่ 11: Wire-up และ Bootstrap

## 11.1 Module Composition Root

```go
// internal/modules/payment/module.go
package payment

import (
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"

    "icmongolang/internal/modules/payment/application"
    "icmongolang/internal/modules/payment/infrastructure/persistence/postgres"
    redisRepo "icmongolang/internal/modules/payment/infrastructure/persistence/redis"
    paymentHTTP "icmongolang/internal/modules/payment/interfaces/http"
    "icmongolang/pkg/logger"
    "icmongolang/pkg/validator"
)

// Dependencies contains all external dependencies for the module
// Dependencies มี dependencies ภายนอกทั้งหมดของโมดูล
type Dependencies struct {
    DB    *gorm.DB
    Redis *redis.Client
    Log   logger.Logger
}

// Module represents the payment module
type Module struct {
    deps Dependencies
}

// NewModule creates a new payment module instance
func NewModule(deps Dependencies) *Module {
    return &Module{deps: deps}
}

// Init wires up all dependencies and registers routes
// Init wire-up dependencies ทั้งหมดและลงทะเบียน route
func (m *Module) Init(r chi.Router, authMW func(http.Handler) http.Handler) error {
    // ========================================
    // 1. Migrations (dev only — production ใช้ SQL migration)
    // ========================================
    if err := m.deps.DB.AutoMigrate(&postgres.PaymentModel{}); err != nil {
        return err
    }

    // ========================================
    // 2. Infrastructure Layer
    // ========================================
    paymentRepo := postgres.NewPaymentRepository(m.deps.DB)

    // Redis cache (optional)
    var paymentCache *redisRepo.PaymentCache
    if m.deps.Redis != nil {
        paymentCache = redisRepo.NewPaymentCache(m.deps.Redis, 10*time.Minute)
    }
    _ = paymentCache // ใช้ใน use case ถ้าต้องการ cache

    // ========================================
    // 3. Application Layer
    // ========================================
    createUC := application.NewCreatePaymentUseCase(paymentRepo, m.deps.Log)
    refundUC := application.NewRefundPaymentUseCase(paymentRepo, m.deps.Log)

    // ========================================
    // 4. Interface Layer
    // ========================================
    v := validator.New()
    handler := paymentHTTP.NewPaymentHandler(createUC, refundUC, v)

    // ========================================
    // 5. Register Routes
    // ========================================
    paymentHTTP.RegisterRoutes(r, &paymentHTTP.Handlers{
        Payment: handler,
    }, authMW)

    m.deps.Log.Info("payment module initialized")
    return nil
}
```

## 11.2 การ Wire-up ใน cmd/api/main.go

```go
// cmd/api/main.go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/joho/godotenv"
    "github.com/redis/go-redis/v9"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "icmongolang/internal/modules/payment"
    appMiddleware "icmongolang/internal/shared/middleware"
    "icmongolang/pkg/jwt"
    "icmongolang/pkg/logger"
)

func main() {
    _ = godotenv.Load()

    log := logger.New()

    // DB
    db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
    if err != nil {
        log.Fatal("db connection failed", "error", err)
    }

    // Redis
    rdb := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_ADDR"),
        Password: os.Getenv("REDIS_PASSWORD"),
    })

    // JWT
    jwtMaker, _ := jwt.NewRSAMaker(
        []byte(os.Getenv("JWT_PRIVATE_KEY")),
        []byte(os.Getenv("JWT_PUBLIC_KEY")),
    )
    authMW := appMiddleware.Auth(jwtMaker)

    // Router
    r := chi.NewRouter()
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    // ========================================
    // Register Modules
    // ========================================
    if err := payment.NewModule(payment.Dependencies{
        DB:    db,
        Redis: rdb,
        Log:   log,
    }).Init(r, authMW); err != nil {
        log.Fatal("payment module failed", "error", err)
    }

    // Health check
    r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    // Server
    srv := &http.Server{Addr: ":8080", Handler: r}
    go func() {
        log.Info("server starting", "addr", srv.Addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal("server error", "error", err)
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("forced shutdown", "error", err)
    }
}
```

---

# บทที่ 12: Testing

## 12.1 Test Pyramid

```
         /\
        /E2E\        5%   - Full stack
       /──────\
      / Integr \    15%   - Repository, DB
     /──────────\
    /   Unit     \ 80%   - Domain, Use Case
   /──────────────\
```

## 12.2 Unit Test Coverage ที่ต้องการ

| Layer | Coverage | เครื่องมือ |
|---|---|---|
| Domain (Entity, VO) | 100% | testing + testify |
| Application (Use Case) | 90%+ | testing + testify + mock |
| Infrastructure (Repo) | 70%+ | Testcontainers |
| Interface (Handler) | 60%+ | httptest |

## 12.3 Integration Test กับ Testcontainers

```go
//go:build integration

package postgres_test

import (
    "context"
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/testcontainers/testcontainers-go/wait"
    gormpg "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "icmongolang/internal/modules/payment/domain/entity"
    "icmongolang/internal/modules/payment/infrastructure/persistence/postgres"
)

func setupPostgres(t *testing.T) (*gorm.DB, func()) {
    ctx := context.Background()

    pgContainer, err := postgres.RunContainer(ctx,
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

    dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
    require.NoError(t, err)

    db, err := gorm.Open(gormpg.Open(dsn), &gorm.Config{})
    require.NoError(t, err)

    require.NoError(t, db.AutoMigrate(&postgres.PaymentModel{}))

    cleanup := func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()
        pgContainer.Terminate(ctx)
    }
    return db, cleanup
}

func TestPaymentRepository_SaveAndFind(t *testing.T) {
    db, cleanup := setupPostgres(t)
    defer cleanup()

    repo := postgres.NewPaymentRepository(db)
    ctx := context.Background()

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
    assert.Equal(t, payment.Amount.String(), found.Amount.String())
    assert.Equal(t, payment.Status, found.Status)
}
```

**รัน integration test:**
```bash
go test -tags=integration ./... -v
```

## 12.4 HTTP Handler Test

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
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"

    "icmongolang/internal/modules/payment/application"
    "icmongolang/internal/modules/payment/application/mocks"
    domainerrors "icmongolang/internal/modules/payment/domain/errors"
    paymentHTTP "icmongolang/internal/modules/payment/interfaces/http"
    "icmongolang/pkg/logger"
    "icmongolang/pkg/validator"
)

func TestCreatePaymentHandler_Success(t *testing.T) {
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

    // Request
    body := map[string]string{
        "order_id": orderID.String(),
        "amount":   "100.00",
        "currency": "THB",
        "method":   "CARD",
    }
    bodyBytes, _ := json.Marshal(body)

    req := httptest.NewRequest("POST", "/api/v1/payments", bytes.NewReader(bodyBytes))
    ctx := context.WithValue(req.Context(), "user_id", uuid.New())
    req = req.WithContext(ctx)
    w := httptest.NewRecorder()

    // Act
    handler.Create(w, req)

    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)

    var response map[string]interface{}
    json.NewDecoder(w.Body).Decode(&response)
    assert.Equal(t, "PENDING", response["status"])
}
```

## 12.5 Test Commands

```bash
# Unit test ทั้งหมด
go test ./...

# Unit test พร้อม verbose
go test -v ./...

# Unit test พร้อม race detector
go test -race ./...

# Unit test พร้อม coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Integration test (build tag)
go test -tags=integration ./...

# Test เฉพาะ module
go test ./internal/modules/payment/...

# Test เฉพาะ function
go test -run TestCreatePaymentUseCase_Success ./...
```

---

# บทที่ 13: กรณีศึกษา Payment Module

## 13.1 โครงสร้างไฟล์ทั้งหมด

```
internal/modules/payment/
├── domain/
│   ├── entity/
│   │   ├── payment.go                    # ✅ สร้างแล้ว
│   │   └── payment_test.go               # ✅ Test ครบ
│   ├── value_object/
│   │   ├── status.go                     # ✅ สร้างแล้ว
│   │   ├── method.go                     # ✅ สร้างแล้ว
│   │   └── money.go                      # (optional)
│   ├── repository/
│   │   └── payment_repository.go         # ✅ Interface
│   ├── event/
│   │   └── events.go                     # ✅ Domain events
│   └── errors/
│       └── errors.go                     # ✅ Sentinel errors
│
├── application/
│   ├── create_payment.go                 # ✅ Use case 1
│   ├── create_payment_test.go            # ✅ Test
│   ├── refund_payment.go                 # ✅ Use case 2
│   └── mocks/
│       └── payment_repository_mock.go    # ✅ Mock
│
├── infrastructure/
│   └── persistence/
│       ├── postgres/
│       │   ├── models.go                 # ✅ GORM model
│       │   └── payment_repo_impl.go      # ✅ Repo impl
│       └── redis/
│           └── payment_cache.go          # ✅ Cache
│
├── interfaces/
│   └── http/
│       ├── payment_handler.go            # ✅ Handler
│       ├── dto.go                        # ✅ Request/Response
│       └── routes.go                     # ✅ Route
│
└── module.go                              # ✅ Composition root
```

## 13.2 คำสั่งรันทั้งหมด (จาก 0 ถึงพร้อมใช้)

```bash
# 1. สร้าง folder structure
export MODULE=payment
mkdir -p internal/modules/$MODULE/{domain/{entity,value_object,repository,event,errors},application/mocks,infrastructure/persistence/{postgres,redis},interfaces/http}

# 2. สร้างไฟล์ทั้งหมด (ตามตัวอย่างในบทที่ 6-9)

# 3. สร้าง migration
touch migrations/$(date +%Y%m%d)_${MODULE}_init.sql

# 4. Wire-up ใน main.go (ตามบทที่ 11)

# 5. Download dependencies
go mod tidy

# 6. Migrate
go run cmd/api/main.go migrate

# 7. Run tests
go test ./internal/modules/payment/... -v

# 8. Build
go build ./...

# 9. Vet
go vet ./...

# 10. Security scan
gosec ./...

# 11. Run server
go run cmd/api/main.go serve

# 12. Test API
curl -X POST http://localhost:8080/api/v1/payments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "550e8400-e29b-41d4-a716-446655440000",
    "amount": "100.00",
    "currency": "THB",
    "method": "CARD"
  }'
```

## 13.3 Dependency Graph

```
                    ┌──────────────────┐
                    │  main.go         │
                    └────────┬─────────┘
                             │
                             ▼
                    ┌──────────────────┐
                    │  module.go       │
                    │  (Composition)   │
                    └────────┬─────────┘
                             │
        ┌────────────────────┼────────────────────┐
        │                    │                    │
        ▼                    ▼                    ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  Interfaces  │    │ Application  │    │ Infrastruct. │
│  (HTTP)      │───▶│ (Use Cases)  │◀───│ (GORM/Redis) │
└──────────────┘    └──────┬───────┘    └──────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │   Domain     │
                    │  (Entity/VO) │
                    └──────────────┘
```

---

# บทที่ 14: Checklist และ Best Practices

## 14.1 Domain Layer Checklist

- [ ] Entity มี constructor `New{Entity}` ที่ validate
- [ ] ไม่มี setter ตรง — เปลี่ยน state ผ่าน behavior method
- [ ] Value Object มี `IsValid()` และ logic methods
- [ ] Repository เป็น interface เท่านั้น
- [ ] Domain error เป็น sentinel error (ไม่ใช้ `fmt.Errorf`)
- [ ] Domain **ไม่ import** gorm, gin, chi, sarama, redis
- [ ] Entity ไม่มี GORM tag (tag อยู่ที่ Model)
- [ ] Unit test coverage = 100%
- [ ] ทุก state transition มี test ทั้ง valid และ invalid
- [ ] Ubiquitous language สอดคล้องทั้งทีม

## 14.2 Application Layer Checklist

- [ ] Use case ละ 1 ไฟล์ `{verb}_{entity}.go`
- [ ] มี Input/Output DTO แยกชัดเจน
- [ ] `Execute()` return `error` ไม่ panic
- [ ] ไม่มี SQL/HTTP/gorm ใน use case
- [ ] Logger ใช้ structured logging (ไม่ใช้ `fmt.Println`)
- [ ] Test ด้วย mock repository
- [ ] มี test ทั้ง happy path และ error cases

## 14.3 Infrastructure Layer Checklist

- [ ] GORM model มี `TableName()` + prefix `{module}_`
- [ ] Repository impl map model ↔ entity ถูกต้อง
- [ ] ไม่มี business logic ใน repository
- [ ] ทุก error map กลับเป็น domain error
- [ ] Connection pool config เหมาะสม
- [ ] Integration test ผ่าน
- [ ] Kafka message ใช้ JSON
- [ ] Redis key มี namespace `{module}:{entity}:{id}`

## 14.4 Interface Layer Checklist

- [ ] Handler ดึง `user_id` จาก context เท่านั้น
- [ ] Route ลงทะเบียนใน `routes.go`
- [ ] Error response เป็น JSON consistent
- [ ] Auth middleware ถูก apply
- [ ] Validation ใช้ library (validator)
- [ ] Response DTO ไม่ leak sensitive data (password, token)
- [ ] Handler test coverage ≥ 60%

## 14.5 Build & Test Checklist

```bash
go build ./...             # ✅ ผ่าน
go vet ./...               # ✅ ผ่าน
go test -race ./...        # ✅ ผ่าน
go test -cover ./...       # ✅ coverage ≥ 80%
gosec ./...                # ✅ ไม่มี HIGH/CRITICAL
golangci-lint run          # ✅ ผ่าน
```

## 14.6 Migration Checklist

- [ ] ตั้งชื่อไฟล์ `YYYYMMDD_{module}_{desc}.sql`
- [ ] ทุกตารางมี prefix `{module}_`
- [ ] มี index สำหรับ query ที่ใช้บ่อย
- [ ] มี FK constraint (ถ้าจำเป็น)
- [ ] มี CHECK constraint
- [ ] มี rollback script
- [ ] Test บน staging แล้ว

## 14.7 Documentation Checklist

- [ ] อัปเดต README ของ module
- [ ] เพิ่มตัวอย่าง .env ถ้ามี env ใหม่
- [ ] อัปเดต docker-compose ถ้ามี service ใหม่
- [ ] Swagger annotation ครบ
- [ ] Comment 2 ภาษา (ไทย/อังกฤษ)

## 14.8 Best Practices สรุป

### ✅ DO

1. **เริ่มจาก Domain** — ออกแบบ business ก่อน technology
2. **Test Domain 100%** — business rule ต้องถูกต้อง
3. **ใช้ Constructor** — ไม่ให้สร้าง Entity ผิด state
4. **Behavior methods** — เปลี่ยน state ผ่าน method
5. **Return error ไม่ panic** — ให้ caller จัดการ
6. **Sentinel errors** — เทียบได้ด้วย `errors.Is`
7. **Structured logging** — ไม่ใช้ `fmt.Println`
8. **Interface ที่เล็ก** — แยกตามการใช้งาน
9. **Immutable ที่ทำได้** — Value Object ต้อง immutable
10. **Comment 2 ภาษา** — ทีมไทย-เทศอ่านได้

### ❌ DON'T

1. **อย่า import framework ใน Domain**
2. **อย่าใส่ business logic ใน Handler หรือ Repository**
3. **อย่าใช้ `*gorm.DB` ใน Use Case**
4. **อย่า hardcode configuration**
5. **อย่าใช้ `float64` กับเงิน** — ใช้ `decimal.Decimal`
6. **อย่าใช้ global state**
7. **อย่าละเลย context** — ส่งทุกที่ที่ควรส่ง
8. **อย่า skip test**
9. **อย่า commit secret**
10. **อย่าใช้ `panic()` ใน business logic**

---

# ภาคผนวก A: Code Templates

## A.1 Domain Entity Template

```go
package entity

import (
    "time"
    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/{MODULE}/domain/errors"
)

type {ENTITY} struct {
    ID        uuid.UUID `json:"id"`
    // ... fields
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func New{ENTITY}(/* required args */) (*{ENTITY}, error) {
    // Validate invariants
    now := time.Now()
    return &{ENTITY}{
        ID:        uuid.New(),
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}

func (e *{ENTITY}) /*Behavior*/() error {
    e.UpdatedAt = time.Now()
    return nil
}
```

## A.2 Value Object Template

```go
package valueobject

type {VO} string

const (
    {VO}A {VO} = "A"
    {VO}B {VO} = "B"
)

func (v {VO}) IsValid() bool {
    switch v {
    case {VO}A, {VO}B:
        return true
    }
    return false
}

func (v {VO}) String() string {
    return string(v)
}
```

## A.3 Repository Interface Template

```go
package repository

import (
    "context"
    "github.com/google/uuid"
    "icmongolang/internal/modules/{MODULE}/domain/entity"
)

type {ENTITY}Repository interface {
    Save(ctx context.Context, e *entity.{ENTITY}) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.{ENTITY}, error)
    Update(ctx context.Context, e *entity.{ENTITY}) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

## A.4 Use Case Template

```go
package application

import (
    "context"
    "icmongolang/internal/modules/{MODULE}/domain/repository"
    "icmongolang/pkg/logger"
)

type {VERB}{ENTITY}UseCase struct {
    repo   repository.{ENTITY}Repository
    logger logger.Logger
}

func New{VERB}{ENTITY}UseCase(repo repository.{ENTITY}Repository, log logger.Logger) *{VERB}{ENTITY}UseCase {
    return &{VERB}{ENTITY}UseCase{repo: repo, logger: log}
}

type {VERB}{ENTITY}Input struct { /* fields */ }
type {VERB}{ENTITY}Output struct { /* fields */ }

func (uc *{VERB}{ENTITY}UseCase) Execute(ctx context.Context, input {VERB}{ENTITY}Input) (*{VERB}{ENTITY}Output, error) {
    // 1. Validate
    // 2. Load aggregate
    // 3. Call domain behavior
    // 4. Persist
    // 5. Side effects
    return &{VERB}{ENTITY}Output{}, nil
}
```

## A.5 GORM Model Template

```go
package postgres

import (
    "time"
    "github.com/google/uuid"
)

type {ENTITY}Model struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    // ... columns
    CreatedAt time.Time `gorm:"default:now()"`
    UpdatedAt time.Time `gorm:"default:now()"`
}

func ({ENTITY}Model) TableName() string {
    return "{MODULE}_{plural}"
}
```

## A.6 HTTP Handler Template

```go
package http

import (
    "encoding/json"
    "net/http"
    "github.com/google/uuid"
    "icmongolang/internal/modules/{MODULE}/application"
    "icmongolang/pkg/httputil"
    "icmongolang/pkg/validator"
)

type {ENTITY}Handler struct {
    {verb}UC   *application.{VERB}{ENTITY}UseCase
    validator *validator.Validator
}

func New{ENTITY}Handler({verb}UC *application.{VERB}{ENTITY}UseCase, v *validator.Validator) *{ENTITY}Handler {
    return &{ENTITY}Handler{{verb}UC: {verb}UC, validator: v}
}

func (h *{ENTITY}Handler) {HTTPVerb}(w http.ResponseWriter, r *http.Request) {
    // 1. Get user from context
    // 2. Parse body
    // 3. Validate
    // 4. Execute use case
    // 5. Return response
}
```

---

# ภาคผนวก B: Common Pitfalls

## B.1 10 ข้อผิดพลาดที่พบบ่อย

### 1. ใส่ GORM tag ใน Entity แทน Model

**❌ ผิด:**
```go
// internal/modules/payment/domain/entity/payment.go
type Payment struct {
    ID uuid.UUID `gorm:"primaryKey"`  // ← GORM leak เข้า Domain
}
```

**✅ ถูก:**
```go
// entity — ไม่มี tag GORM
type Payment struct {
    ID uuid.UUID `json:"id"`
}

// infrastructure/persistence/postgres/models.go — มี tag GORM
type PaymentModel struct {
    ID uuid.UUID `gorm:"type:uuid;primaryKey"`
}
```

### 2. ใช้ `float64` กับเงิน

**❌ ผิด:**
```go
Amount float64  // 0.1 + 0.2 = 0.30000000000000004
```

**✅ ถูก:**
```go
Amount decimal.Decimal  // แม่นยำ
```

### 3. ไม่ห่อ error จาก repository

**❌ ผิด:**
```go
if err := r.db.Create(&m).Error; err != nil {
    return err  // caller ได้ gorm error ไม่ใช่ domain error
}
```

**✅ ถูก:**
```go
if err := r.db.Create(&m).Error; err != nil {
    if errors.Is(err, gorm.ErrDuplicatedKey) {
        return domainerrors.ErrPaymentAlreadyExists
    }
    return err
}
```

### 4. Business logic ใน Handler

**❌ ผิด:**
```go
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    if amount.LessThan(decimal.NewFromInt(20)) {
        // ← business rule ใน Handler
        return error
    }
}
```

**✅ ถูก:**
```go
// ย้ายไป Domain หรือ Use Case
func (p *Payment) validateAmount() error {
    if p.Amount.LessThan(decimal.NewFromInt(20)) {
        return domainerrors.ErrAmountTooLow
    }
    return nil
}
```

### 5. Import ข้าม Layer

**❌ ผิด:**
```go
// Domain import Infrastructure
import "icmongolang/internal/modules/payment/infrastructure/postgres"
```

**✅ ถูก:**
```go
// Domain import แค่ stdlib + pkg
import (
    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)
```

### 6. Test ไม่ครอบ state transitions

**❌ ผิด:**
```go
func TestPayment(t *testing.T) {
    p, _ := NewPayment(...)
    // ไม่มี test สำหรับ invalid transitions
}
```

**✅ ถูก:**
```go
func TestPayment_InvalidTransitions(t *testing.T) {
    // ทดสอบ PENDING → SUCCESS ตรงๆ ต้อง fail
    // ทดสอบ refund จาก PENDING ต้อง fail
    // ทดสอบ refund 2 ครั้ง ต้อง fail
}
```

### 7. ใช้ `panic` ใน business logic

**❌ ผิด:**
```go
func (p *Payment) Refund() {
    if p.Status != StatusSuccess {
        panic("cannot refund")  // ← อย่า
    }
}
```

**✅ ถูก:**
```go
func (p *Payment) Refund() error {
    if p.Status != StatusSuccess {
        return domainerrors.ErrCannotRefund
    }
    return nil
}
```

### 8. ไม่มี context

**❌ ผิด:**
```go
func (r *repo) FindByID(id uuid.UUID) (*Payment, error) {
    return r.db.First(&m, id).Error  // ไม่มี context
}
```

**✅ ถูก:**
```go
func (r *repo) FindByID(ctx context.Context, id uuid.UUID) (*Payment, error) {
    return r.db.WithContext(ctx).First(&m, id).Error
}
```

### 9. Hardcode configuration

**❌ ผิด:**
```go
time.Sleep(30 * time.Second)  // hardcode
db, _ := gorm.Open("localhost:5432")  // hardcode
```

**✅ ถูก:**
```go
cfg := config.Load()  // อ่านจาก env
time.Sleep(cfg.RetryDelay)
```

### 10. ไม่มี log

**❌ ผิด:**
```go
if err != nil {
    return nil, err  // เงียบ
}
```

**✅ ถูก:**
```go
if err != nil {
    uc.logger.Error("save payment failed",
        "error", err,
        "order_id", input.OrderID,
    )
    return nil, err
}
```

## B.2 Anti-Patterns ที่ควรหลีกเลี่ยง

| Anti-Pattern | ปัญหา | แก้ |
|---|---|---|
| **Anemic Model** | Entity มีแต่ data ไม่มี behavior | ย้าย behavior เข้า Entity |
| **God Service** | 1 service ทำทุกอย่าง | แยกเป็น use case |
| **Shotgun Surgery** | แก้ 1 feature ต้องแก้ 10 ไฟล์ | รวม responsibility |
| **Feature Envy** | Method ใช้ data ของ class อื่น | ย้าย method |
| **Primitive Obsession** | ใช้ `string` แทน VO | สร้าง Value Object |
| **Long Parameter List** | Function รับ 8+ parameters | รวมเป็น struct |
| **God Object** | 1 struct มี 30+ fields | แยก Aggregate |

---

# ภาคผนวก C: Quick Reference Card

## C.1 โครงสร้าง Module

```
internal/modules/{module}/
├── domain/                ← ❤️ Business logic
│   ├── entity/            ← Aggregate Root
│   ├── value_object/      ← Immutable VO
│   ├── repository/        ← Interface
│   ├── service/           ← Domain service
│   ├── event/             ← Domain events
│   └── errors/            ← Sentinel errors
├── application/           ← 🎯 Use cases
│   ├── {verb}_{entity}.go ← 1 use case = 1 file
│   └── mocks/             ← Test mocks
├── infrastructure/        ← 🔧 Adapters
│   └── persistence/
│       ├── postgres/      ← GORM + Repo impl
│       └── redis/         ← Cache
├── interfaces/            ← 🌐 Inbound
│   └── http/
│       ├── handler.go
│       ├── dto.go
│       └── routes.go
└── module.go              ← Composition root
```

## C.2 คำสั่งที่ใช้บ่อย

```bash
# สร้าง module ใหม่
mkdir -p internal/modules/{module}/{domain/{entity,value_object,repository,errors},application/mocks,infrastructure/persistence/postgres,interfaces/http}

# Build
go build ./...

# Test
go test -race -cover ./...
go test -tags=integration ./...

# Lint
go vet ./...
golangci-lint run

# Security
gosec ./...

# Migration
go run cmd/api/main.go migrate

# Run
go run cmd/api/main.go serve
# หรือ
air
```

## C.3 Error Mapping

| Domain Error | HTTP Status |
|---|---|
| `ErrNotFound` | 404 |
| `ErrAlreadyExists` | 409 |
| `ErrInvalid*` | 400 |
| `ErrUnauthorized` | 401 |
| `ErrForbidden` | 403 |
| `ErrRateLimit` | 429 |
| (default) | 500 |

## C.4 Naming Quick Reference

| สิ่ง | Pattern | ตัวอย่าง |
|---|---|---|
| Module | lowercase | `payment` |
| Package | lowercase | `entity` |
| Entity | PascalCase | `Payment` |
| Table | `{module}_{plural}` | `payment_transactions` |
| Column | snake_case | `user_id` |
| Index | `idx_{table}_{col}` | `idx_payment_transactions_user_id` |
| Migration | `YYYYMMDD_{module}_{desc}.sql` | `20260401_payment_init.sql` |
| Kafka Topic | `{module}.{entity}.{action}` | `payment.txn.created` |
| Redis Key | `{module}:{entity}:{id}` | `payment:txn:abc` |
| Use Case | `{Verb}{Entity}UseCase` | `CreatePaymentUseCase` |
| Domain Error | `Err{Reason}` | `ErrInsufficientBalance` |

## C.5 Dependency Rules (MEMORIZE)

```
Domain         → stdlib, uuid, decimal
Application    → Domain
Infrastructure → Domain, Application
Interfaces     → ทุก Layer
```

**ห้าม:**
- Domain → gorm, gin, chi, sarama, redis, kafka
- Application → Infrastructure, Interfaces
- Infrastructure → Interfaces

## C.6 Checklist 30 วินาที

```
✅ go build ./...          ผ่าน
✅ go vet ./...            ผ่าน
✅ go test -race ./...     ผ่าน
✅ gosec ./...             ไม่มี HIGH/CRITICAL
✅ Domain ไม่ import framework
✅ 1 use case = 1 file
✅ Entity มี constructor
✅ ทุก error เป็น sentinel
✅ Migration มี prefix
✅ Log structured
```

---

# 📝 ข้อมูลเอกสาร

**ชื่อเอกสาร:** คู่มือสร้าง Module ใหม่ฉบับสมบูรณ์ (Complete Guide to Creating Go Modules)
**เวอร์ชัน:** 1.0
**วันที่:** เมษายน 2026
**จำนวนหน้า:** ~120 หน้า (ประมาณ)
**ระดับ:** Intermediate - Advanced

**ผู้อ่านเป้าหมาย:**
- Go Developer ที่ต้องการเรียนรู้ Clean Architecture
- Tech Lead ที่ต้องวางมาตรฐานทีม
- Architect ที่ออกแบบระบบ Modular Monolith

**ข้อกำหนดเบื้องต้น:**
- Go 1.22+
- พื้นฐาน PostgreSQL
- พื้นฐาน GORM
- เข้าใจ Interface และ Error handling ใน Go

**เอกสารที่เกี่ยวข้อง:**
- เล่ม 2: คู่มือแก้ไข Module เดิม
- เล่ม 3: คู่มือขาย Module
- เล่ม 4: คู่มือทดสอบและ Deployment
- เล่ม 5: คู่มือบำรุงรักษาและ Scale

**อ้างอิง:**
- Clean Architecture — Robert C. Martin
- Domain-Driven Design — Eric Evans
- Implementing DDD — Vaughn Vernon
- The Go Programming Language — Donovan & Kernighan
- โปรเจกต์ `icmongolang` (33 modules, 21 shared packages)

---

**END OF BOOK 1**

> 📌 **ขั้นถัดไป:** อ่านเล่ม 2 — คู่มือแก้ไข Module เดิม
> ที่จะสอนวิธีแก้ Module โดยไม่ทำลายของเดิม พร้อม Migration แบบ Zero-Downtime

---

**พิมพ์เมื่อ:** เมษายน 2026
**ผู้จัดทำ:** ทีมสถาปัตยกรรมซอฟต์แวร์ icmongolang