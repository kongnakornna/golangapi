# 📘 Application Design Document
## ระบบ ERP+SaaS สำหรับฟาร์มเห็ดอัจฉริยะ (Smart Farm ERP)

---

# สารบัญ

1. [บทนำ (Introduction)](#1-บทนำ-introduction)
2. [บทนิยาม (Definitions)](#2-บทนิยาม-definitions)
3. [บทหัวข้อ (Topics)](#3-บทหัวข้อ-topics)
4. [การออกแบบ Workflow](#4-การออกแบบ-workflow)
5. [Case Study](#5-case-study)
6. [โครงสร้างโฟลเดอร์](#6-โครงสร้างโฟลเดอร์)
7. [หลักการทำงาน (Concept)](#7-หลักการทำงาน-concept)
8. [Workflow และ Dataflow](#8-workflow-และ-dataflow)
9. [Code Template](#9-code-template)
10. [Checklist Module](#10-checklist-module)
11. [Security Code](#11-security-code)
12. [Load Test](#12-load-test)
13. [สรุป](#13-สรุป)
14. [คู่มือการทดสอบ](#14-คู่มือการทดสอบ)
15. [คู่มือการใช้งาน](#15-คู่มือการใช้งาน)
16. [คู่มือการบำรุงรักษา](#16-คู่มือการบำรุงรักษา)
17. [คู่มือการขยาย/แก้ไข](#17-คู่มือการขยายแก้ไข)
18. [Git Flow](#18-git-flow)
19. [Code Review & PR](#19-code-review--pr)
20. [CI/CD](#20-cicd)
21. [Root Cause Analysis (RCA)](#21-root-cause-analysis-rca)

---

# 1. บทนำ (Introduction)

## 1.1 Business Model Canvas (BMC)

### 🎯 Business Model Canvas — Smart Mushroom Farm ERP+SaaS

| องค์ประกอบ | รายละเอียด |
|---|---|
| **Key Partners** | ผู้ผลิตเซ็นเซอร์ IoT, ผู้ให้บริการ Cloud (AWS/GCP), ธนาคาร (Payment Gateway), หน่วยงานราชการ (อย., GAP), ซัพพลายเออร์ขี้เลื่อย/หัวเชื้อ, ขนส่ง (Kerry, Flash) |
| **Key Activities** | พัฒนาแพลตฟอร์ม ERP+SaaS, ดูแลระบบ Multi-tenant, ให้บริการ IoT Integration, วิเคราะห์ข้อมูล AI, อบรม/สนับสนุนลูกค้า |
| **Key Resources** | ทีมพัฒนา Go, Infrastructure Cloud, องค์ความรู้ฟาร์มเห็ด, โมเดล AI พยากรณ์ผลผลิต, ฐานข้อมูลลูกค้า |
| **Value Propositions** | 1) ERP ครบวงจรสำหรับฟาร์มเห็ด 2) ใช้ฟรีแบบ SaaS Multi-tenant 3) IoT ควบคุมโรงเพาะ 4) AI พยากรณ์ผลผลิต 5) รองรับ GAP/GMP/อย. 6) ต้นทุนต่ำกว่า ERP ทั่วไป 80% |
| **Customer Relationships** | Self-service + Customer Success, Training ออนไลน์, Community ฟาร์มเห็ด, Support 24/7 |
| **Channels** | Web App, Mobile App (iOS/Android), LINE OA, Facebook, Partner (สหกรณ์การเกษตร), Direct Sales |
| **Customer Segments** | 1) ฟาร์มเห็ดขนาดเล็ก (1-10 โรง) 2) ฟาร์มเห็ดขนาดกลาง (10-50 โรง) 3) สหกรณ์/กลุ่มเกษตรกร 4) โรงงานแปรรูปเห็ด 5) ผู้รับงานภาครัฐ (ส่งเห็ดให้โครงการ) |
| **Cost Structure** | Cloud Infrastructure (30%), R&D (40%), Sales & Marketing (15%), Support (10%), Admin (5%) |
| **Revenue Streams** | 1) SaaS Subscription (Free/Pro/Enterprise) 2) IoT Hardware (ขาย/เช่า) 3) AI Forecast Add-on 4) Setup & Training Fee 5) Transaction Fee (Payment Gateway) 6) Data Insight Report |

### 📊 Revenue Model (SaaS Tiers)

| Tier | ราคา/เดือน | ฟีเจอร์ |
|---|---|---|
| **Free** | 0 บาท | 1 โรงเพาะ, 1 user, Basic Inventory |
| **Pro** | 990 บาท | 10 โรงเพาะ, 10 users, Production + CRM + POS |
| **Business** | 2,990 บาท | 50 โรงเพาะ, 50 users, + IoT + AI Forecast |
| **Enterprise** | 9,990 บาท | ไม่จำกัด, + Custom Module + SLA 99.9% |

---

# 2. บทนิยาม (Definitions)

| คำศัพท์ | ความหมาย |
|---|---|
| **ERP** | Enterprise Resource Planning — ระบบรวมกระบวนการทางธุรกิจ |
| **SaaS** | Software as a Service — ซอฟต์แวร์ให้บริการผ่าน Cloud |
| **Multi-tenant** | สถาปัตยกรรมที่ 1 ระบบให้บริการหลายองค์กร (tenant) แยกข้อมูลกัน |
| **Clean Architecture** | สถาปัตยกรรมแยก Layer: domain → application → infrastructure → interfaces |
| **DDD** | Domain-Driven Design — ออกแบบโดยยึด Domain เป็นศูนย์กลาง |
| **Aggregate Root** | Entity หลักที่ควบคุม Entity อื่นใน Aggregate เดียวกัน |
| **Value Object** | Object ที่ไม่มี identity เปลี่ยนไม่ได้ (immutable) |
| **Repository** | Interface สำหรับเข้าถึง Data Layer |
| **Tenant ID** | รหัสองค์กร ใช้แยกข้อมูลในระบบ Multi-tenant |
| **Batch/Lot** | รุ่นการผลิต ใช้ติดตามย้อนกลับ |
| **BOM** | Bill of Materials — สูตรการผลิต |
| **OEE** | Overall Equipment Effectiveness — ประสิทธิภาพโดยรวมของเครื่องจักร |
| **GAP** | Good Agricultural Practice — มาตรฐานการเกษตรดี |
| **GMP** | Good Manufacturing Practice — มาตรฐานการผลิต |
| **SLA** | Service Level Agreement — ข้อตกลงระดับบริการ |
| **CAPA** | Corrective and Preventive Action — มาตรการแก้ไขและป้องกัน |
| **KPI** | Key Performance Indicator — ตัวชี้วัดผลงาน |
| **CQRS** | Command Query Responsibility Segregation — แยกอ่าน/เขียน |
| **Event-Driven** | สถาปัตยกรรมที่สื่อสารผ่าน Event |
| **RCA** | Root Cause Analysis — การวิเคราะห์สาเหตุราก |

---

# 3. บทหัวข้อ (Topics)

## 3.1 โครงสร้างการทำงาน (Architecture Overview)

```
┌─────────────────────────────────────────────────────────────┐
│                    PRESENTATION LAYER                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │  Web App │  │ Mobile   │  │  POS     │  │ IoT Dash │    │
│  │  (React) │  │ (Flutter)│  │ (Tablet) │  │ (Web)    │    │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘    │
└───────┼─────────────┼─────────────┼─────────────┼──────────┘
        │             │             │             │
        └─────────────┴──────┬──────┴─────────────┘
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                    API GATEWAY (Kong/Nginx)                  │
│  ┌──────────┬──────────┬──────────┬──────────┐              │
│  │  Auth    │  Rate    │  Tenant  │  Logging │              │
│  │  JWT     │  Limit   │  Resolve │  Tracing │              │
│  └──────────┴──────────┴──────────┴──────────┘              │
└─────────────────────────┬───────────────────────────────────┘
                          ▼
┌─────────────────────────────────────────────────────────────┐
│              APPLICATION LAYER (Go Modular Monolith)         │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  erp/inventory  │  erp/production  │  erp/finance   │   │
│  ├──────────────────────────────────────────────────────┤   │
│  │  crm/crm        │  pos/pos         │  ecommerce     │   │
│  ├──────────────────────────────────────────────────────┤   │
│  │  logistics      │  iot/iot         │  forecast      │   │
│  ├──────────────────────────────────────────────────────┤   │
│  │  reporting      │  KPI/Analytics                      │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────┬───────────────────────────────────┘
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                    INFRASTRUCTURE LAYER                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │PostgreSQL│  │  Redis   │  │  Kafka   │  │ClickHouse│    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │  MQTT    │  │  MinIO   │  │  ML      │  │  Elastic │    │
│  │  Broker  │  │  (S3)    │  │  Python  │  │  Search  │    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
└─────────────────────────────────────────────────────────────┘
```

## 3.2 วัตถุประสงค์ (Objectives)

1. **รวมระบบ**: รวมทุกกระบวนการฟาร์มเห็ดไว้ในที่เดียว
2. **ลดต้นทุน**: ลดต้นทุนซอฟต์แวร์ 80% ด้วย SaaS
3. **เพิ่มผลผลิต**: ใช้ IoT + AI เพิ่มผลผลิต 30%
4. **รองรับมาตรฐาน**: GAP, GMP, อย., HACCP
5. **ขยายได้**: Modular Monolith → Microservices ได้ในอนาคต
6. **Multi-tenant**: 1 ระบบให้บริการหลายฟาร์ม

## 3.3 กลุ่มเป้าหมาย (Target Users)

| กลุ่ม | ลักษณะ | Module ที่ใช้ |
|---|---|---|
| ฟาร์มขนาดเล็ก | 1-10 โรง, 1-5 users | Inventory, POS, Reporting |
| ฟาร์มขนาดกลาง | 10-50 โรง, 10-50 users | + Production, CRM, Logistics |
| สหกรณ์ | 50+ โรง, 100+ users | + IoT, AI Forecast, Finance |
| โรงงานแปรรูป | ผลิตภัณฑ์แปรรูป | + Quality, GMP, Batch Tracking |
| ผู้รับงานภาครัฐ | ส่งโครงการ | + Government Contract |

## 3.4 ความรู้พื้นฐาน (Prerequisites)

- **ภาษา Go** 1.21+
- **PostgreSQL** 15+
- **Docker & Docker Compose**
- **Clean Architecture + DDD**
- **REST API / gRPC**
- **Kafka / MQTT**
- **Git / CI/CD**

## 3.5 เนื้อหาโดยย่อ

| Module | วัตถุประสงค์ | ประโยชน์ |
|---|---|---|
| **erp/inventory** | จัดการสต็อกก้อนเชื้อ/เห็ด | ลดของเสีย, รู้สต็อก real-time |
| **erp/production** | วางแผนผลิต | เพิ่มผลผลิต, ลดต้นทุน |
| **erp/finance** | บัญชี/ต้นทุน | รู้กำไรจริงต่อก้อน |
| **crm/crm** | ดูแลลูกค้า | เพิ่มยอดขายซ้ำ |
| **pos/pos** | ขายหน้าร้าน | ตัดสต็อกอัตโนมัติ |
| **ecommerce** | ขายออนไลน์ | เข้าถึงลูกค้ากว้าง |
| **logistics** | ขนส่ง | ส่งถึงมือทันเวลา |
| **iot/iot** | เซ็นเซอร์ | ควบคุมโรงเพาะ |
| **forecast** | AI พยากรณ์ | วางแผนล่วงหน้า |
| **reporting** | รายงาน/KPI | ตัดสินใจด้วยข้อมูล |

---

# 4. การออกแบบ Workflow

## 4.1 Workflow หลัก: วงจรชีวิตก้อนเชื้อ → ขายเห็ด

```
┌──────────────────────────────────────────────────────────────────┐
│              MUSHROOM FARM LIFECYCLE WORKFLOW                     │
└──────────────────────────────────────────────────────────────────┘

[1] จัดซื้อวัสดุ              [2] ผลิตก้อนเชื้อ           [3] บ่มเชื้อ
    │                            │                          │
    ▼                            ▼                          ▼
┌─────────┐              ┌─────────────┐            ┌─────────────┐
│Procure- │              │ Production  │            │ Incubation  │
│ment     │─────────────▶│ Order       │───────────▶│ (21-30 วัน) │
└─────────┘              └─────────────┘            └──────┬──────┘
    │                          │                           │
    │ PR → PO → GR             │ BOM → WO → QC             │
    ▼                          ▼                           ▼
┌─────────┐              ┌─────────────┐            ┌─────────────┐
│Supplier │              │ Inventory   │            │ Stock       │
│Payment  │              │ (Raw Mat)   │            │ (Incubated) │
└─────────┘              └─────────────┘            └──────┬──────┘
                                                            │
[4] เปิดดอก                  [5] เก็บเกี่ยว              [6] แพ็ค/ขาย
    │                            │                          │
    ▼                            ▼                          ▼
┌─────────────┐            ┌─────────────┐            ┌─────────────┐
│ IoT Control │            │ Harvest     │            │ QC + Grade  │
│ (Temp/Hum)  │───────────▶│ Record      │───────────▶│ A/B/C       │
└─────────────┘            └─────────────┘            └──────┬──────┘
                                                              │
                              ┌───────────────────────────────┼──────────┐
                              ▼                               ▼          ▼
                        ┌──────────┐                   ┌──────────┐ ┌────────┐
                        │ POS      │                   │E-commerce│ │Direct  │
                        │ (หน้าร้าน)│                   │(ออนไลน์)  │ │(โรงงาน) │
                        └────┬─────┘                   └────┬─────┘ └───┬────┘
                             │                              │           │
                             └──────────────┬───────────────┴───────────┘
                                            ▼
                                    ┌─────────────┐
                                    │ Logistics   │
                                    │ (ส่งของ)     │
                                    └──────┬──────┘
                                           ▼
                                    ┌─────────────┐
                                    │ Finance     │
                                    │ (Invoice)   │
                                    └──────┬──────┘
                                           ▼
                                    ┌─────────────┐
                                    │ Reporting   │
                                    │ (KPI)       │
                                    └─────────────┘
```

## 4.2 Workflow การขายแบบ Multi-channel

```
┌─────────────────────────────────────────────────────────────┐
│                    ORDER-TO-CASH WORKFLOW                    │
└─────────────────────────────────────────────────────────────┘

Customer Order
      │
      ├── POS (Walk-in) ──────┐
      ├── E-commerce ─────────┤
      ├── Phone/Line ─────────┤
      └── B2B Contract ───────┤
                              ▼
                    ┌─────────────────┐
                    │ Order Capture   │
                    │ (Validate)      │
                    └────────┬────────┘
                             ▼
                    ┌─────────────────┐
                    │ Stock Check     │
                    │ (Inventory)     │
                    └────────┬────────┘
                             │
              ┌──────────────┼──────────────┐
              ▼              ▼              ▼
        ┌──────────┐  ┌──────────┐  ┌──────────┐
        │มีสต็อก    │  │สต็อกน้อย  │  │ไม่มีสต็อก │
        │→ Reserve │  │→ Backorder│  │→ สั่งผลิต │
        └────┬─────┘  └────┬─────┘  └────┬─────┘
             └─────────────┼─────────────┘
                           ▼
                  ┌─────────────────┐
                  │ Pick & Pack     │
                  │ (QC + Grade)    │
                  └────────┬────────┘
                           ▼
                  ┌─────────────────┐
                  │ Ship            │
                  │ (Logistics)     │
                  └────────┬────────┘
                           ▼
                  ┌─────────────────┐
                  │ Deliver         │
                  │ (POD)           │
                  └────────┬────────┘
                           ▼
                  ┌─────────────────┐
                  │ Invoice         │
                  │ (Finance)       │
                  └────────┬────────┘
                           ▼
                  ┌─────────────────┐
                  │ Payment         │
                  │ (Reconcile)     │
                  └────────┬────────┘
                           ▼
                  ┌─────────────────┐
                  │ Update KPI      │
                  │ (Reporting)     │
                  └─────────────────┘
```

## 4.3 Workflow การพยากรณ์และวางแผน

```
┌─────────────────────────────────────────────────────────────┐
│              DEMAND FORECAST & PLANNING WORKFLOW             │
└─────────────────────────────────────────────────────────────┘

Historical Data (ClickHouse)
      │
      ├── ยอดขายย้อนหลัง 12 เดือน
      ├── ข้อมูลอากาศ (จาก IoT)
      ├── ฤดูกาล/เทศกาล
      └── แคมเปญการตลาด
              │
              ▼
      ┌─────────────────┐
      │ Feature Eng.    │
      │ (Python/ML)     │
      └────────┬────────┘
               ▼
      ┌─────────────────┐
      │ Train Model     │
      │ (Prophet/LSTM)  │
      └────────┬────────┘
               ▼
      ┌─────────────────┐
      │ Predict Demand  │
      │ (30/60/90 วัน)  │
      └────────┬────────┘
               ▼
      ┌─────────────────┐
      │ MRP Calculation │
      │ (ต้องการวัสดุ)   │
      └────────┬────────┘
               ▼
      ┌─────────────────┐
      │ Production Plan │
      │ (แผนผลิต)        │
      └────────┬────────┘
               ▼
      ┌─────────────────┐
      │ Purchase Plan   │
      │ (แผนจัดซื้อ)     │
      └────────┬────────┘
               ▼
      ┌─────────────────┐
      │ Notify Manager  │
      │ (Approval)      │
      └─────────────────┘
```

---

# 5. Case Study

## 5.1 Case Study: ฟาร์มเห็ดนางฟ้า "ภูเขียว" จ.ชัยภูมิ

### ปัญหาก่อนใช้ระบบ
- ฟาร์มขนาด 30 โรงเพาะ
- จดบันทึกในสมุด → ข้อมูลหาย, คำนวณต้นทุนผิด
- ขาดสต็อกบ่อย → ผลผลิตลด 20%
- ส่งของช้า → ลูกค้ายกเลิก 15%
- ไม่รู้ต้นทุนจริงต่อก้อน → ตั้งราคาผิด

### แนวทางแก้ไขด้วยระบบ

| ปัญหา | Module ที่ใช้ | ผลลัพธ์ |
|---|---|---|
| บันทึกในสมุด | erp/inventory + Mobile App | ข้อมูล real-time |
| ขาดสต็อก | Reorder Point + AI Forecast | ลดขาดสต็อก 90% |
| ส่งช้า | logistics + IoT | ส่งตรงเวลา 98% |
| ไม่รู้ต้นทุน | erp/finance | รู้ต้นทุนต่อก้อน |
| ลูกค้าหนี | crm + POS | ยอดขายซ้ำ +40% |

### ผลลัพธ์หลังใช้ 6 เดือน
- ผลผลิตเพิ่ม **35%** (500 → 675 กก./วัน)
- ต้นทุนลด **22%** (จาก 45 → 35 บาท/กก.)
- กำไรเพิ่ม **60%** (จาก 150,000 → 240,000 บาท/เดือน)
- ลูกค้าประจำเพิ่ม **3 เท่า**

## 5.2 Case Study: โรงงานแปรรูปเห็ด "เมืองเลย"

### ปัญหา
- ต้องได้ GMP/อย. → เอกสารกระจัดกระจาย
- ติดตาม Batch ไม่ได้ → recall ไม่ได้
- ข้อร้องเรียนลูกค้า → แก้ช้า

### แนวทางแก้ไข

| ปัญหา | Module | ผลลัพธ์ |
|---|---|---|
| GMP | Quality Control + Audit Trail | ผ่าน GMP ใน 3 เดือน |
| Batch Tracking | Batch Tracking + QR | Recall ได้ใน 1 ชม. |
| ข้อร้องเรียน | CRM + CAPA | ปิดภายใน 48 ชม. |

## 5.3 Case Study: รับงานภาครัฐ — สหกรณ์การเกษตร

### ปัญหา
- ประมูลงานราชการบ่อย → พลาด deadline
- เปลี่ยนเจ้าหน้าที่ → ติดตามขาด

### แนวทางแก้ไข

| ปัญหา | Module | ผลลัพธ์ |
|---|---|---|
| พลาด deadline | Deadline Automation | พลาด 0 ครั้ง |
| เปลี่ยนคน | Personnel Transfer | ฟื้นตัวใน 7 วัน |
| ชนะประมูลน้อย | Tender Pipeline | ชนะเพิ่ม 2 เท่า |

---

# 6. โครงสร้างโฟลเดอร์

## 6.1 Folder Structure หลัก

```
smart-mushroom-erp/
├── cmd/
│   └── api/
│       └── main.go                    # Entry point
├── internal/
│   ├── modules/
│   │   ├── erp/
│   │   │   ├── inventory/             # ✅ พร้อมใช้
│   │   │   ├── production/            # ⏳ กำลังพัฒนา
│   │   │   └── finance/               # ⏳ กำลังพัฒนา
│   │   ├── crm/crm/
│   │   ├── pos/pos/
│   │   ├── ecommerce/ecommerce/
│   │   ├── logistics/logistics/
│   │   ├── iot/iot/
│   │   ├── forecast/forecast/
│   │   └── reporting/reporting/
│   ├── shared/                        # Shared code
│   │   ├── config/
│   │   ├── database/
│   │   ├── logger/
│   │   ├── middleware/
│   │   ├── errors/
│   │   ├── events/
│   │   └── utils/
│   └── platform/                      # Platform services
│       ├── auth/
│       ├── tenant/
│       ├── storage/
│       └── messaging/
├── pkg/                               # Public packages
│   ├── uuid/
│   ├── decimal/
│   └── validator/
├── migrations/                        # DB migrations
├── deployments/
│   ├── docker/
│   ├── k8s/
│   └── terraform/
├── docs/
├── tests/
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── scripts/
├── .env.example
├── docker-compose.yml
├── Makefile
├── go.mod
└── README.md
```

## 6.2 โครงสร้าง Module มาตรฐาน (ใช้ทุก Module)

```
internal/modules/{prefix}/{module}/
├── domain/                            # 🎯 Business Logic
│   ├── entity/                        # Aggregate Roots & Entities
│   │   ├── {aggregate}.go
│   │   └── {entity}.go
│   ├── value_object/                  # Immutable Value Objects
│   │   ├── {vo1}.go
│   │   └── {vo2}.go
│   ├── repository/                    # Repository Interfaces
│   │   └── {aggregate}_repository.go
│   ├── service/                       # Domain Services
│   │   └── {module}_domain_service.go
│   ├── event/                         # Domain Events
│   │   └── {event}.go
│   └── errors/
│       └── errors.go
├── application/                       # 🔄 Use Cases
│   ├── command/                       # Write operations
│   │   ├── create_{aggregate}.go
│   │   └── update_{aggregate}.go
│   ├── query/                         # Read operations
│   │   └── get_{aggregate}.go
│   ├── dto/                           # Data Transfer Objects
│   │   ├── request.go
│   │   └── response.go
│   └── port/                          # Ports (interfaces)
│       └── {port}.go
├── infrastructure/                    # 🔧 Implementation
│   ├── persistence/
│   │   ├── postgres/
│   │   │   ├── {aggregate}_repo_impl.go
│   │   │   └── models.go
│   │   ├── redis/
│   │   │   └── {module}_cache.go
│   │   └── clickhouse/
│   │       └── {module}_analytics.go
│   ├── messaging/
│   │   ├── kafka_producer.go
│   │   └── kafka_consumer.go
│   └── external/
│       └── {external_service}.go
├── interfaces/                        # 🌐 API Layer
│   ├── http/
│   │   ├── handler/
│   │   │   └── {aggregate}_handler.go
│   │   ├── request/
│   │   │   └── {aggregate}_request.go
│   │   ├── response/
│   │   │   └── {aggregate}_response.go
│   │   └── routes.go
│   ├── grpc/
│   │   └── {module}.proto
│   └── middleware/
│       └── tenant.go
├── module.go                          # Module wire-up
└── README.md
```

---

# 7. หลักการทำงาน (Concept)

## 7.1 Clean Architecture Layers

```
┌─────────────────────────────────────────────────────────┐
│  INTERFACES (HTTP/gRPC/CLI)                              │
│  ┌───────────────────────────────────────────────────┐  │
│  │  APPLICATION (Use Cases, DTO)                      │  │
│  │  ┌─────────────────────────────────────────────┐  │  │
│  │  │  DOMAIN (Entity, VO, Repository Interface)  │  │  │
│  │  │  ┌───────────────────────────────────────┐  │  │  │
│  │  │  │  Business Rules (Pure Go)             │  │  │  │
│  │  │  └───────────────────────────────────────┘  │  │  │
│  │  └─────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────┘  │
│  INFRASTRUCTURE (DB, Cache, MQ, External)                │
└─────────────────────────────────────────────────────────┘
        ▲                                        ▲
        │ Dependency Rule:                      │
        │ Outer → Inner (ห้ามย้อนกลับ)            │
        │ Domain ห้าม import อะไรนอกจาก stdlib    │
```

## 7.2 Dependency Rule

| Layer | Import ได้ | Import ไม่ได้ |
|---|---|---|
| **Domain** | stdlib, pkg | application, infrastructure, interfaces |
| **Application** | domain, pkg | infrastructure, interfaces |
| **Infrastructure** | domain, application, pkg | interfaces |
| **Interfaces** | ทุก layer | - |

## 7.3 Multi-tenant Pattern

```go
// ทุก Repository ต้องรับ tenant_id
type ProductRepository interface {
    FindByID(ctx context.Context, tenantID, id string) (*Product, error)
    FindAll(ctx context.Context, tenantID string, filter Filter) ([]*Product, error)
    Save(ctx context.Context, tenantID string, product *Product) error
}

// Middleware inject tenant จาก JWT
func TenantMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tenantID := extractTenantFromJWT(c)
        ctx := context.WithValue(c.Request.Context(), TenantKey, tenantID)
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}
```

## 7.4 Event-Driven Pattern

```go
// Domain Event
type StockAdjustedEvent struct {
    TenantID  string
    ProductID string
    Quantity  decimal.Decimal
    Reason    string
    Timestamp time.Time
}

// Publish หลัง commit
func (s *AdjustStockUseCase) Execute(ctx context.Context, cmd AdjustStockCommand) error {
    // 1. Load aggregate
    product, _ := s.repo.FindByID(ctx, cmd.TenantID, cmd.ProductID)
    
    // 2. Execute business logic
    if err := product.AdjustStock(cmd.Quantity, cmd.Reason); err != nil {
        return err
    }
    
    // 3. Save
    if err := s.repo.Save(ctx, cmd.TenantID, product); err != nil {
        return err
    }
    
    // 4. Publish event
    event := StockAdjustedEvent{...}
    return s.eventBus.Publish(ctx, "stock.adjusted", event)
}
```

---

# 8. Workflow และ Dataflow

## 8.1 Dataflow: สร้างสินค้า (Create Product)

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│  Client  │───▶│  Router  │───▶│Middleware│───▶│ Handler  │
│ (Postman)│    │  (Gin)   │    │(Tenant)  │    │          │
└──────────┘    └──────────┘    └──────────┘    └────┬─────┘
                                                     │
                                                     ▼
                                              ┌──────────────┐
                                              │  Use Case    │
                                              │ CreateProduct│
                                              └──────┬───────┘
                                                     │
                                    ┌────────────────┼────────────────┐
                                    ▼                ▼                ▼
                            ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
                            │  Validate    │ │  Domain      │ │  Repository  │
                            │  (DTO)       │ │  Logic       │ │  (Save)      │
                            └──────────────┘ └──────────────┘ └──────┬───────┘
                                                                     │
                                                                     ▼
                                                            ┌──────────────┐
                                                            │  PostgreSQL  │
                                                            └──────┬───────┘
                                                                   │
                                                                   ▼
                                                            ┌──────────────┐
                                                            │Kafka Producer│
                                                            │product.created│
                                                            └──────┬───────┘
                                                                   │
                                                                   ▼
                                                            ┌──────────────┐
                                                            │  Consumers   │
                                                            │(Reporting,   │
                                                            │ Forecast)    │
                                                            └──────────────┘
```

## 8.2 Dataflow: ปรับสต็อก (Adjust Stock)

```
┌─────────────────────────────────────────────────────────────────┐
│                  ADJUST STOCK DATAFLOW                            │
└─────────────────────────────────────────────────────────────────┘

[1] HTTP Request
    POST /api/v1/inventory/stock/adjust
    Body: { product_id, quantity, reason }
    Header: Authorization: Bearer <JWT>
         │
         ▼
[2] Middleware Chain
    ├── AuthMiddleware → validate JWT
    ├── TenantMiddleware → extract tenant_id
    ├── RateLimitMiddleware → check quota
    └── LoggingMiddleware → log request
         │
         ▼
[3] Handler
    ├── Parse request → AdjustStockRequest
    ├── Validate input
    └── Call UseCase
         │
         ▼
[4] Use Case: AdjustStock
    ├── Load Product (Repository)
    │   └── SELECT * FROM products WHERE tenant_id=$1 AND id=$2
    ├── Load Stock Balance
    │   └── SELECT * FROM stock_balances WHERE ...
    ├── Execute Domain Logic
    │   ├── product.AdjustStock(qty, reason)
    │   ├── Validate: qty != 0
    │   ├── Check: stock ไม่ติดลบ
    │   └── Create StockMovement
    ├── Save (Transaction)
    │   ├── BEGIN
    │   ├── UPDATE stock_balances SET quantity = ...
    │   ├── INSERT INTO stock_movements (...)
    │   ├── INSERT INTO audit_trail (...)
    │   └── COMMIT
    └── Publish Event
        └── Kafka: stock.adjusted
         │
         ▼
[5] Response
    { "success": true, "data": { movement_id, new_balance } }
         │
         ▼
[6] Async Consumers
    ├── Reporting → update KPI
    ├── Forecast → update feature store
    └── Notification → alert ถ้า stock ต่ำ
```

## 8.3 Dataflow: IoT → Alert → Action

```
┌─────────────────────────────────────────────────────────────────┐
│                    IoT DATAFLOW                                  │
└─────────────────────────────────────────────────────────────────┘

[1] Sensor (ESP32)
    ├── อุณหภูมิ 32°C
    ├── ความชื้น 85%
    └── CO2 800 ppm
         │
         ▼ MQTT
[2] MQTT Broker (Mosquitto)
    └── Topic: farm/{tenant_id}/{device_id}/reading
         │
         ▼
[3] MQTT Consumer (Go)
    ├── Parse payload
    ├── Validate schema
    └── Publish to Kafka
         │
         ▼
[4] Kafka: iot.readings
         │
         ├──────────────────────┐
         ▼                      ▼
[5a] Time-series Writer    [5b] Alert Engine
    └── ClickHouse             ├── Check rules
        INSERT readings        │   └── temp > 30°C?
                               ├── Create Alert
                               └── Publish: iot.alert
                                        │
                                        ▼
                               [6] Alert Consumer
                                   ├── Send LINE OA
                                   ├── Send Push
                                   └── Auto Control
                                       └── เปิดพัดลม
```

---

# 9. Code Template

## 9.1 Domain Entity: Product (Aggregate Root)

```go
// internal/modules/erp/inventory/domain/entity/product.go
package entity

import (
    "errors"
    "time"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)

// Product is the Aggregate Root for inventory management.
// Product เป็น Aggregate Root สำหรับการจัดการสินค้าคงคลัง
type Product struct {
    ID          string          // รหัสสินค้า
    TenantID    string          // รหัสองค์กร (Multi-tenant)
    SKU         string          // รหัส SKU
    Name        string          // ชื่อสินค้า
    Description string          // คำอธิบาย
    Category    string          // หมวดหมู่
    Unit        string          // หน่วยนับ
    CostPrice   decimal.Decimal // ราคาทุน
    SalePrice   decimal.Decimal // ราคาขาย
    MinStock    decimal.Decimal // สต็อกขั้นต่ำ
    MaxStock    decimal.Decimal // สต็อกสูงสุด
    Status      ProductStatus   // สถานะ
    CreatedAt   time.Time       // วันที่สร้าง
    UpdatedAt   time.Time       // วันที่แก้ไข
    Version     int             // Optimistic Lock
}

// ProductStatus represents the status of a product.
// ProductStatus แทนสถานะของสินค้า
type ProductStatus string

const (
    ProductStatusActive   ProductStatus = "active"   // ใช้งาน
    ProductStatusInactive ProductStatus = "inactive" // ไม่ใช้งาน
    ProductStatusDeleted  ProductStatus = "deleted"  // ลบ
)

// NewProduct creates a new Product with validation.
// NewProduct สร้าง Product ใหม่พร้อม validation
func NewProduct(tenantID, sku, name, unit string, costPrice, salePrice decimal.Decimal) (*Product, error) {
    // Validate tenant ID
    // ตรวจสอบรหัสองค์กร
    if tenantID == "" {
        return nil, errors.New("tenant_id is required")
    }

    // Validate SKU
    // ตรวจสอบ SKU
    if sku == "" {
        return nil, errors.New("sku is required")
    }

    // Validate name
    // ตรวจสอบชื่อ
    if name == "" {
        return nil, errors.New("name is required")
    }

    // Validate prices
    // ตรวจสอบราคา
    if costPrice.IsNegative() || salePrice.IsNegative() {
        return nil, errors.New("prices cannot be negative")
    }

    now := time.Now()
    return &Product{
        ID:        uuid.NewString(),
        TenantID:  tenantID,
        SKU:       sku,
        Name:      name,
        Unit:      unit,
        CostPrice: costPrice,
        SalePrice: salePrice,
        Status:    ProductStatusActive,
        CreatedAt: now,
        UpdatedAt: now,
        Version:   1,
    }, nil
}

// UpdatePrice updates the product price with validation.
// UpdatePrice อัปเดตราคาสินค้าพร้อม validation
func (p *Product) UpdatePrice(costPrice, salePrice decimal.Decimal) error {
    // Check if prices are valid
    // ตรวจสอบว่าราคาถูกต้อง
    if costPrice.IsNegative() || salePrice.IsNegative() {
        return errors.New("prices cannot be negative")
    }

    // Sale price should be greater than cost price
    // ราคาขายควรมากกว่าราคาทุน
    if salePrice.LessThan(costPrice) {
        return errors.New("sale price cannot be less than cost price")
    }

    p.CostPrice = costPrice
    p.SalePrice = salePrice
    p.UpdatedAt = time.Now()
    p.Version++
    return nil
}

// CalculateMargin calculates the profit margin.
// CalculateMargin คำนวณกำไรขั้นต้น
func (p *Product) CalculateMargin() decimal.Decimal {
    if p.CostPrice.IsZero() {
        return decimal.Zero
    }
    // Margin = (Sale - Cost) / Sale * 100
    // กำไร = (ราคาขาย - ราคาทุน) / ราคาขาย * 100
    return p.SalePrice.Sub(p.CostPrice).Div(p.SalePrice).Mul(decimal.NewFromInt(100))
}

// IsLowStock checks if the stock is below minimum.
// IsLowStock ตรวจสอบว่าสต็อกต่ำกว่าขั้นต่ำหรือไม่
func (p *Product) IsLowStock(currentStock decimal.Decimal) bool {
    return currentStock.LessThanOrEqual(p.MinStock)
}
```

## 9.2 Domain Repository Interface

```go
// internal/modules/erp/inventory/domain/repository/product_repository.go
package repository

import (
    "context"

    "github.com/yourorg/smart-mushroom-erp/internal/modules/erp/inventory/domain/entity"
)

// ProductRepository defines the contract for product persistence.
// ProductRepository กำหนดสัญญาสำหรับการจัดเก็บสินค้า
type ProductRepository interface {
    // FindByID retrieves a product by ID within a tenant.
    // FindByID ค้นหาสินค้าด้วย ID ภายในองค์กร
    FindByID(ctx context.Context, tenantID, id string) (*entity.Product, error)

    // FindBySKU retrieves a product by SKU within a tenant.
    // FindBySKU ค้นหาสินค้าด้วย SKU ภายในองค์กร
    FindBySKU(ctx context.Context, tenantID, sku string) (*entity.Product, error)

    // FindAll retrieves all products with pagination.
    // FindAll ค้นหาสินค้าทั้งหมดพร้อม pagination
    FindAll(ctx context.Context, tenantID string, filter ProductFilter) ([]*entity.Product, int64, error)

    // Save creates or updates a product.
    // Save สร้างหรืออัปเดตสินค้า
    Save(ctx context.Context, tenantID string, product *entity.Product) error

    // Delete soft-deletes a product.
    // Delete ลบสินค้าแบบ soft delete
    Delete(ctx context.Context, tenantID, id string) error
}

// ProductFilter defines filter criteria for product queries.
// ProductFilter กำหนดเงื่อนไขการกรองสำหรับการค้นหาสินค้า
type ProductFilter struct {
    Search   string
    Category string
    Status   string
    Page     int
    PageSize int
}
```

## 9.3 Application Use Case

```go
// internal/modules/erp/inventory/application/command/create_product.go
package command

import (
    "context"
    "fmt"

    "github.com/shopspring/decimal"
    "github.com/yourorg/smart-mushroom-erp/internal/modules/erp/inventory/domain/entity"
    "github.com/yourorg/smart-mushroom-erp/internal/modules/erp/inventory/domain/repository"
    "github.com/yourorg/smart-mushroom-erp/internal/shared/events"
)

// CreateProductCommand contains the data needed to create a product.
// CreateProductCommand เก็บข้อมูลที่จำเป็นสำหรับการสร้างสินค้า
type CreateProductCommand struct {
    TenantID    string          // รหัสองค์กร
    SKU         string          // รหัส SKU
    Name        string          // ชื่อสินค้า
    Description string          // คำอธิบาย
    Category    string          // หมวดหมู่
    Unit        string          // หน่วยนับ
    CostPrice   decimal.Decimal // ราคาทุน
    SalePrice   decimal.Decimal // ราคาขาย
    MinStock    decimal.Decimal // สต็อกขั้นต่ำ
    MaxStock    decimal.Decimal // สต็อกสูงสุด
}

// CreateProductResult contains the result of creating a product.
// CreateProductResult เก็บผลลัพธ์ของการสร้างสินค้า
type CreateProductResult struct {
    ProductID string // รหัสสินค้าที่สร้าง
    SKU       string // รหัส SKU
}

// CreateProductUseCase handles the product creation logic.
// CreateProductUseCase จัดการตรรกะการสร้างสินค้า
type CreateProductUseCase struct {
    productRepo repository.ProductRepository // Repository สำหรับสินค้า
    eventBus    events.EventBus              // Event Bus สำหรับ publish events
}

// NewCreateProductUseCase creates a new use case instance.
// NewCreateProductUseCase สร้าง instance ของ use case ใหม่
func NewCreateProductUseCase(
    productRepo repository.ProductRepository,
    eventBus events.EventBus,
) *CreateProductUseCase {
    return &CreateProductUseCase{
        productRepo: productRepo,
        eventBus:    eventBus,
    }
}

// Execute runs the create product use case.
// Execute รัน use case การสร้างสินค้า
func (uc *CreateProductUseCase) Execute(ctx context.Context, cmd CreateProductCommand) (*CreateProductResult, error) {
    // Step 1: Check if SKU already exists
    // ขั้นตอนที่ 1: ตรวจสอบว่า SKU มีอยู่แล้วหรือไม่
    existing, err := uc.productRepo.FindBySKU(ctx, cmd.TenantID, cmd.SKU)
    if err != nil && !errors.Is(err, repository.ErrNotFound) {
        return nil, fmt.Errorf("check sku: %w", err)
    }
    if existing != nil {
        return nil, ErrSKUAlreadyExists
    }

    // Step 2: Create domain entity
    // ขั้นตอนที่ 2: สร้าง domain entity
    product, err := entity.NewProduct(
        cmd.TenantID,
        cmd.SKU,
        cmd.Name,
        cmd.Unit,
        cmd.CostPrice,
        cmd.SalePrice,
    )
    if err != nil {
        return nil, fmt.Errorf("create product entity: %w", err)
    }

    // Step 3: Set additional fields
    // ขั้นตอนที่ 3: ตั้งค่าฟิลด์เพิ่มเติม
    product.Description = cmd.Description
    product.Category = cmd.Category
    product.MinStock = cmd.MinStock
    product.MaxStock = cmd.MaxStock

    // Step 4: Save to repository
    // ขั้นตอนที่ 4: บันทึกไปยัง repository
    if err := uc.productRepo.Save(ctx, cmd.TenantID, product); err != nil {
        return nil, fmt.Errorf("save product: %w", err)
    }

    // Step 5: Publish domain event
    // ขั้นตอนที่ 5: Publish domain event
    event := events.ProductCreatedEvent{
        TenantID:  cmd.TenantID,
        ProductID: product.ID,
        SKU:       product.SKU,
        Name:      product.Name,
        Timestamp: product.CreatedAt,
    }
    if err := uc.eventBus.Publish(ctx, "product.created", event); err != nil {
        // Log error but don't fail the operation
        // Log error แต่ไม่ทำให้ operation ล้มเหลว
        log.Printf("failed to publish event: %v", err)
    }

    // Step 6: Return result
    // ขั้นตอนที่ 6: คืนค่าผลลัพธ์
    return &CreateProductResult{
        ProductID: product.ID,
        SKU:       product.SKU,
    }, nil
}
```

## 9.4 HTTP Handler

```go
// internal/modules/erp/inventory/interfaces/http/handler/product_handler.go
package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/yourorg/smart-mushroom-erp/internal/modules/erp/inventory/application/command"
    "github.com/yourorg/smart-mushroom-erp/internal/shared/middleware"
)

// ProductHandler handles HTTP requests for products.
// ProductHandler จัดการ HTTP requests สำหรับสินค้า
type ProductHandler struct {
    createUC *command.CreateProductUseCase // Use case สำหรับสร้างสินค้า
}

// NewProductHandler creates a new product handler.
// NewProductHandler สร้าง handler สำหรับสินค้าใหม่
func NewProductHandler(createUC *command.CreateProductUseCase) *ProductHandler {
    return &ProductHandler{createUC: createUC}
}

// CreateProductRequest represents the HTTP request body.
// CreateProductRequest แทน request body ของ HTTP
type CreateProductRequest struct {
    SKU         string  `json:"sku" binding:"required"`         // รหัส SKU
    Name        string  `json:"name" binding:"required"`        // ชื่อสินค้า
    Description string  `json:"description"`                    // คำอธิบาย
    Category    string  `json:"category"`                       // หมวดหมู่
    Unit        string  `json:"unit" binding:"required"`        // หน่วยนับ
    CostPrice   float64 `json:"cost_price" binding:"required"`  // ราคาทุน
    SalePrice   float64 `json:"sale_price" binding:"required"`  // ราคาขาย
    MinStock    float64 `json:"min_stock"`                      // สต็อกขั้นต่ำ
    MaxStock    float64 `json:"max_stock"`                      // สต็อกสูงสุด
}

// CreateProduct handles POST /api/v1/inventory/products
// CreateProduct จัดการ POST /api/v1/inventory/products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    // Step 1: Parse request body
    // ขั้นตอนที่ 1: Parse request body
    var req CreateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "invalid request: " + err.Error(),
        })
        return
    }

    // Step 2: Get tenant ID from context (injected by middleware)
    // ขั้นตอนที่ 2: ดึง tenant ID จาก context (injected โดย middleware)
    tenantID := middleware.GetTenantID(c)

    // Step 3: Build command
    // ขั้นตอนที่ 3: สร้าง command
    cmd := command.CreateProductCommand{
        TenantID:    tenantID,
        SKU:         req.SKU,
        Name:        req.Name,
        Description: req.Description,
        Category:    req.Category,
        Unit:        req.Unit,
        CostPrice:   decimal.NewFromFloat(req.CostPrice),
        SalePrice:   decimal.NewFromFloat(req.SalePrice),
        MinStock:    decimal.NewFromFloat(req.MinStock),
        MaxStock:    decimal.NewFromFloat(req.MaxStock),
    }

    // Step 4: Execute use case
    // ขั้นตอนที่ 4: รัน use case
    result, err := h.createUC.Execute(c.Request.Context(), cmd)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }

    // Step 5: Return response
    // ขั้นตอนที่ 5: คืนค่า response
    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "data": gin.H{
            "product_id": result.ProductID,
            "sku":        result.SKU,
        },
    })
}
```

## 9.5 Repository Implementation (PostgreSQL)

```go
// internal/modules/erp/inventory/infrastructure/persistence/postgres/product_repo_impl.go
package postgres

import (
    "context"
    "database/sql"
    "errors"
    "fmt"

    "github.com/yourorg/smart-mushroom-erp/internal/modules/erp/inventory/domain/entity"
    "github.com/yourorg/smart-mushroom-erp/internal/modules/erp/inventory/domain/repository"
)

// ProductRepositoryImpl implements ProductRepository using PostgreSQL.
// ProductRepositoryImpl implement ProductRepository ด้วย PostgreSQL
type ProductRepositoryImpl struct {
    db *sql.DB // Database connection
}

// NewProductRepository creates a new repository instance.
// NewProductRepository สร้าง instance ของ repository ใหม่
func NewProductRepository(db *sql.DB) repository.ProductRepository {
    return &ProductRepositoryImpl{db: db}
}

// FindByID retrieves a product by ID within a tenant.
// FindByID ค้นหาสินค้าด้วย ID ภายในองค์กร
func (r *ProductRepositoryImpl) FindByID(ctx context.Context, tenantID, id string) (*entity.Product, error) {
    // SQL query with tenant isolation
    // SQL query พร้อมการแยก tenant
    query := `
        SELECT id, tenant_id, sku, name, description, category, unit,
               cost_price, sale_price, min_stock, max_stock,
               status, created_at, updated_at, version
        FROM products
        WHERE tenant_id = $1 AND id = $2 AND status != 'deleted'
    `

    var p entity.Product
    var costPrice, salePrice, minStock, maxStock float64

    // Execute query
    // รัน query
    err := r.db.QueryRowContext(ctx, query, tenantID, id).Scan(
        &p.ID, &p.TenantID, &p.SKU, &p.Name, &p.Description, &p.Category,
        &p.Unit, &costPrice, &salePrice, &minStock, &maxStock,
        &p.Status, &p.CreatedAt, &p.UpdatedAt, &p.Version,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, repository.ErrNotFound
        }
        return nil, fmt.Errorf("query product: %w", err)
    }

    // Convert float64 to decimal
    // แปลง float64 เป็น decimal
    p.CostPrice = decimal.NewFromFloat(costPrice)
    p.SalePrice = decimal.NewFromFloat(salePrice)
    p.MinStock = decimal.NewFromFloat(minStock)
    p.MaxStock = decimal.NewFromFloat(maxStock)

    return &p, nil
}

// Save creates or updates a product.
// Save สร้างหรืออัปเดตสินค้า
func (r *ProductRepositoryImpl) Save(ctx context.Context, tenantID string, p *entity.Product) error {
    // Check if product exists
    // ตรวจสอบว่าสินค้ามีอยู่หรือไม่
    existing, err := r.FindByID(ctx, tenantID, p.ID)
    if err != nil && !errors.Is(err, repository.ErrNotFound) {
        return fmt.Errorf("check existing: %w", err)
    }

    if existing == nil {
        // INSERT new product
        // INSERT สินค้าใหม่
        query := `
            INSERT INTO products (
                id, tenant_id, sku, name, description, category, unit,
                cost_price, sale_price, min_stock, max_stock,
                status, created_at, updated_at, version
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
        `
        _, err = r.db.ExecContext(ctx, query,
            p.ID, p.TenantID, p.SKU, p.Name, p.Description, p.Category,
            p.Unit, p.CostPrice, p.SalePrice, p.MinStock, p.MaxStock,
            p.Status, p.CreatedAt, p.UpdatedAt, p.Version,
        )
        if err != nil {
            return fmt.Errorf("insert product: %w", err)
        }
        return nil
    }

    // UPDATE existing product with optimistic locking
    // UPDATE สินค้าที่มีอยู่พร้อม optimistic locking
    query := `
        UPDATE products
        SET sku = $1, name = $2, description = $3, category = $4,
            unit = $5, cost_price = $6, sale_price = $7,
            min_stock = $8, max_stock = $9, status = $10,
            updated_at = $11, version = version + 1
        WHERE tenant_id = $12 AND id = $13 AND version = $14
    `
    result, err := r.db.ExecContext(ctx, query,
        p.SKU, p.Name, p.Description, p.Category,
        p.Unit, p.CostPrice, p.SalePrice,
        p.MinStock, p.MaxStock, p.Status,
        time.Now(), tenantID, p.ID, p.Version,
    )
    if err != nil {
        return fmt.Errorf("update product: %w", err)
    }

    // Check if update affected any rows (optimistic lock)
    // ตรวจสอบว่า update มีผลกับแถวใดหรือไม่ (optimistic lock)
    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("rows affected: %w", err)
    }
    if rows == 0 {
        return repository.ErrConcurrentUpdate
    }

    return nil
}
```

## 9.6 Module Wire-up

```go
// internal/modules/erp/inventory/module.go
package inventory

import (
    "database/sql"

    "github.com/yourorg/smart-mushroom-erp/internal/modules/erp/inventory/application/command"
    "github.com/yourorg/smart-mushroom-erp/internal/modules/erp/inventory/infrastructure/persistence/postgres"
    "github.com/yourorg/smart-mushroom-erp/internal/modules/erp/inventory/interfaces/http/handler"
    "github.com/yourorg/smart-mushroom-erp/internal/shared/events"
)

// Module represents the inventory module.
// Module แทนโมดูลสินค้าคงคลัง
type Module struct {
    // Use Cases
    // Use Cases
    CreateProductUC *command.CreateProductUseCase

    // Handlers
    // Handlers
    ProductHandler *handler.ProductHandler
}

// NewModule initializes the inventory module with all dependencies.
// NewModule เริ่มต้นโมดูลสินค้าคงคลังพร้อม dependencies ทั้งหมด
func NewModule(db *sql.DB, eventBus events.EventBus) *Module {
    // Infrastructure layer
    // Infrastructure layer
    productRepo := postgres.NewProductRepository(db)

    // Application layer
    // Application layer
    createProductUC := command.NewCreateProductUseCase(productRepo, eventBus)

    // Interfaces layer
    // Interfaces layer
    productHandler := handler.NewProductHandler(createProductUC)

    return &Module{
        CreateProductUC: createProductUC,
        ProductHandler:  productHandler,
    }
}

// RegisterRoutes registers HTTP routes for the module.
// RegisterRoutes ลงทะเบียน HTTP routes สำหรับโมดูล
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
    products := router.Group("/products")
    {
        products.POST("", m.ProductHandler.CreateProduct)
        // products.GET("", m.ProductHandler.ListProducts)
        // products.GET("/:id", m.ProductHandler.GetProduct)
        // products.PUT("/:id", m.ProductHandler.UpdateProduct)
        // products.DELETE("/:id", m.ProductHandler.DeleteProduct)
    }
}
```

## 9.7 Main Entry Point

```go
// cmd/api/main.go
package main

import (
    "context"
    "database/sql"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
    "github.com/yourorg/smart-mushroom-erp/internal/modules/erp/inventory"
    "github.com/yourorg/smart-mushroom-erp/internal/shared/config"
    "github.com/yourorg/smart-mushroom-erp/internal/shared/events"
    "github.com/yourorg/smart-mushroom-erp/internal/shared/middleware"
)

func main() {
    // Load configuration
    // โหลด configuration
    cfg := config.Load()

    // Connect to database
    // เชื่อมต่อฐานข้อมูล
    db, err := sql.Open("postgres", cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("failed to connect database: %v", err)
    }
    defer db.Close()

    // Initialize event bus (Kafka)
    // เริ่มต้น event bus (Kafka)
    eventBus := events.NewKafkaEventBus(cfg.KafkaBrokers)

    // Initialize modules
    // เริ่มต้นโมดูล
    inventoryModule := inventory.NewModule(db, eventBus)

    // Setup router
    // ตั้งค่า router
    router := gin.Default()
    router.Use(middleware.TenantMiddleware())
    router.Use(middleware.AuthMiddleware(cfg.JWTSecret))

    // Register routes
    // ลงทะเบียน routes
    api := router.Group("/api/v1")
    inventoryModule.RegisterRoutes(api.Group("/inventory"))

    // Health check
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    // Start server
    // เริ่ม server
    srv := &http.Server{
        Addr:    ":" + cfg.Port,
        Handler: router,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("server error: %v", err)
        }
    }()

    log.Printf("server started on port %s", cfg.Port)

    // Graceful shutdown
    // ปิด server อย่าง graceful
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("server shutdown: %v", err)
    }

    log.Println("server stopped")
}
```

---

# 10. Checklist Module

## 10.1 Module Completion Checklist

| # | รายการ | สถานะ |
|---|---|---|
| 1 | Domain Entity ครบ (Aggregate Root + VO) | ☐ |
| 2 | Repository Interface ใน domain | ☐ |
| 3 | Repository Implementation ใน infrastructure | ☐ |
| 4 | Use Case (Command + Query) | ☐ |
| 5 | DTO (Request + Response) | ☐ |
| 6 | HTTP Handler | ☐ |
| 7 | Routes Registration | ☐ |
| 8 | Tenant Middleware | ☐ |
| 9 | Auth Middleware | ☐ |
| 10 | Validation (Input) | ☐ |
| 11 | Error Handling | ☐ |
| 12 | Domain Events | ☐ |
| 13 | Kafka Producer/Consumer | ☐ |
| 14 | Unit Tests (≥80% coverage) | ☐ |
| 15 | Integration Tests | ☐ |
| 16 | API Documentation (Swagger) | ☐ |
| 17 | Migration Scripts | ☐ |
| 18 | Audit Trail | ☐ |
| 19 | Logging | ☐ |
| 20 | Metrics (Prometheus) | ☐ |

## 10.2 Module ตามลำดับความสำคัญ

| ลำดับ | Module | สถานะ | ความสำคัญ |
|---|---|---|---|
| 1 | erp/inventory | ✅ | 🔴 สูง |
| 2 | erp/production | ⏳ | 🔴 สูง |
| 3 | erp/finance | ⏳ | 🔴 สูง |
| 4 | crm/crm | ⏳ | 🟡 กลาง |
| 5 | pos/pos | ⏳ | 🟡 กลาง |
| 6 | iot/iot | ⏳ | 🟡 กลาง |
| 7 | ecommerce | ⏳ | 🟢 ต่ำ |
| 8 | logistics | ⏳ | 🟢 ต่ำ |
| 9 | forecast | ⏳ | 🟢 ต่ำ |
| 10 | reporting | ⏳ | 🟡 กลาง |

---

# 11. Security Code

## 11.1 Security Checklist

| # | รายการ | รายละเอียด |
|---|---|---|
| 1 | **Authentication** | JWT + Refresh Token |
| 2 | **Authorization** | RBAC (Role-Based Access Control) |
| 3 | **Tenant Isolation** | ทุก query ต้องมี tenant_id |
| 4 | **Input Validation** | Validate ทุก input |
| 5 | **SQL Injection** | ใช้ Parameterized Query |
| 6 | **XSS** | Sanitize output |
| 7 | **CSRF** | ใช้ CSRF Token |
| 8 | **Rate Limiting** | จำกัด request/นาที |
| 9 | **HTTPS** | บังคับ HTTPS |
| 10 | **Password Hashing** | bcrypt (cost ≥ 12) |
| 11 | **Secret Management** | ใช้ Vault/AWS Secrets |
| 12 | **Audit Log** | บันทึกทุก action |
| 13 | **Data Encryption** | AES-256 at rest |
| 14 | **PII Protection** | Mask ข้อมูลส่วนตัว |
| 15 | **Dependency Scan** | Trivy/Snyk |

## 11.2 ตัวอย่าง Code: Tenant Isolation

```go
// internal/shared/middleware/tenant.go
package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

// ContextKey is a custom type for context keys.
// ContextKey เป็น type กำหนดเองสำหรับ context keys
type ContextKey string

const (
    TenantIDKey ContextKey = "tenant_id" // Key สำหรับ tenant ID
    UserIDKey   ContextKey = "user_id"   // Key สำหรับ user ID
    RoleKey     ContextKey = "role"      // Key สำหรับ role
)

// TenantMiddleware extracts tenant_id from JWT and injects into context.
// TenantMiddleware ดึง tenant_id จาก JWT และใส่ใน context
func TenantMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get Authorization header
        // ดึง Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "missing authorization header",
            })
            return
        }

        // Extract token
        // แยก token
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "invalid authorization format",
            })
            return
        }

        // Parse JWT
        // Parse JWT
        token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return []byte(getJWTSecret()), nil
        })
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "invalid token",
            })
            return
        }

        // Extract claims
        // ดึง claims
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "invalid claims",
            })
            return
        }

        tenantID, _ := claims["tenant_id"].(string)
        userID, _ := claims["user_id"].(string)
        role, _ := claims["role"].(string)

        if tenantID == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "missing tenant_id",
            })
            return
        }

        // Inject into context
        // ใส่ใน context
        ctx := context.WithValue(c.Request.Context(), TenantIDKey, tenantID)
        ctx = context.WithValue(ctx, UserIDKey, userID)
        ctx = context.WithValue(ctx, RoleKey, role)

        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}

// GetTenantID extracts tenant ID from context.
// GetTenantID ดึง tenant ID จาก context
func GetTenantID(c *gin.Context) string {
    if v, ok := c.Request.Context().Value(TenantIDKey).(string); ok {
        return v
    }
    return ""
}

// GetUserID extracts user ID from context.
// GetUserID ดึง user ID จาก context
func GetUserID(c *gin.Context) string {
    if v, ok := c.Request.Context().Value(UserIDKey).(string); ok {
        return v
    }
    return ""
}

// GetRole extracts role from context.
// GetRole ดึง role จาก context
func GetRole(c *gin.Context) string {
    if v, ok := c.Request.Context().Value(RoleKey).(string); ok {
        return v
    }
    return ""
}
```

## 11.3 RBAC Middleware

```go
// internal/shared/middleware/rbac.go
package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// RequireRole checks if the user has the required role.
// RequireRole ตรวจสอบว่าผู้ใช้มี role ที่ต้องการหรือไม่
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := GetRole(c)
        for _, role := range roles {
            if userRole == role {
                c.Next()
                return
            }
        }
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
            "error": "insufficient permissions",
        })
    }
}

// Usage in routes
// การใช้งานใน routes
func setupRoutes(r *gin.Engine) {
    api := r.Group("/api/v1")
    api.Use(TenantMiddleware())

    // Admin only
    // เฉพาะ admin
    admin := api.Group("/admin")
    admin.Use(RequireRole("admin"))
    {
        admin.GET("/users", listUsers)
    }

    // Manager + Admin
    // manager + admin
    products := api.Group("/products")
    products.Use(RequireRole("admin", "manager"))
    {
        products.POST("", createProduct)
    }
}
```

---

# 12. Load Test

## 12.1 Load Test Plan

| รายการ | รายละเอียด |
|---|---|
| **เครื่องมือ** | k6, Vegeta, JMeter |
| **เป้าหมาย** | 1,000 concurrent users |
| **RPS** | 5,000 requests/second |
| **Response Time** | P95 < 200ms |
| **Error Rate** | < 0.1% |
| **Duration** | 30 นาที |

## 12.2 k6 Script

```javascript
// tests/load/inventory_test.js
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
// Metrics กำหนดเอง
const errorRate = new Rate('errors');

// Test configuration
// การตั้งค่า test
export const options = {
  stages: [
    { duration: '2m', target: 100 },   // Ramp up
    { duration: '5m', target: 500 },   // Stay
    { duration: '2m', target: 1000 },  // Peak
    { duration: '5m', target: 1000 },  // Stay
    { duration: '2m', target: 0 },     // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<200'],  // 95% < 200ms
    errors: ['rate<0.001'],            // Error < 0.1%
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const TOKEN = __ENV.JWT_TOKEN;

// Main test function
// ฟังก์ชัน test หลัก
export default function () {
  // Test 1: Create product
  // Test 1: สร้างสินค้า
  const createPayload = JSON.stringify({
    sku: `SKU-${__VU}-${__ITER}`,
    name: `Test Product ${__VU}`,
    unit: 'kg',
    cost_price: 100,
    sale_price: 150,
  });

  const createRes = http.post(`${BASE_URL}/api/v1/inventory/products`, createPayload, {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${TOKEN}`,
    },
  });

  check(createRes, {
    'create status 201': (r) => r.status === 201,
    'create response time < 200ms': (r) => r.timings.duration < 200,
  }) || errorRate.add(1);

  // Test 2: List products
  // Test 2: แสดงรายการสินค้า
  const listRes = http.get(`${BASE_URL}/api/v1/inventory/products?page=1&page_size=20`, {
    headers: { 'Authorization': `Bearer ${TOKEN}` },
  });

  check(listRes, {
    'list status 200': (r) => r.status === 200,
    'list response time < 100ms': (r) => r.timings.duration < 100,
  }) || errorRate.add(1);

  sleep(1);
}
```

## 12.3 รัน Load Test

```bash
# ติดตั้ง k6
brew install k6  # macOS
# หรือ
docker pull grafana/k6

# รัน test
k6 run --vus 100 --duration 30s tests/load/inventory_test.js

# รันพร้อม env
BASE_URL=https://api.mushroom-erp.com \
JWT_TOKEN=your_token \
k6 run tests/load/inventory_test.js

# รันด้วย Docker
docker run --rm -i grafana/k6 run - < tests/load/inventory_test.js
```

---

# 13. สรุป

## 13.1 ประโยชน์ที่ได้รับ

| # | ประโยชน์ | รายละเอียด |
|---|---|---|
| 1 | **รวมศูนย์** | ทุกกระบวนการในที่เดียว |
| 2 | **ลดต้นทุน** | SaaS ลดต้นทุน 80% |
| 3 | **เพิ่มผลผลิต** | IoT + AI เพิ่ม 30% |
| 4 | **มาตรฐาน** | GAP/GMP/อย. |
| 5 | **ขยายได้** | Modular → Microservices |
| 6 | **Multi-tenant** | 1 ระบบหลายฟาร์ม |
| 7 | **Real-time** | ข้อมูลทันที |
| 8 | **ตัดสินใจดี** | KPI + Dashboard |

## 13.2 ข้อควรระวัง

| # | ข้อควรระวัง | แนวทาง |
|---|---|---|
| 1 | ข้อมูลรั่วระหว่าง tenant | บังคับ tenant_id ทุก query |
| 2 | Performance ตก | Index + Cache + Load Test |
| 3 | Dependency ซับซ้อน | ใช้ Wire/DI |
| 4 | Migration ยาก | Version + Rollback plan |
| 5 | Security | Audit + Pen Test |
| 6 | Downtime | Zero-downtime deployment |

## 13.3 ข้อดี

- ✅ Clean Architecture → ทดสอบง่าย
- ✅ DDD → Business logic ชัด
- ✅ Modular → แยกทีมพัฒนา
- ✅ Event-Driven → Loose coupling
- ✅ Multi-tenant → ต้นทุนต่ำ

## 13.4 ข้อเสีย

- ❌ ซับซ้อนกว่า Monolith ปกติ
- ❌ ต้องมีทีมที่มีประสบการณ์
- ❌ ต้นทุนเริ่มต้นสูง
- ❌ ใช้เวลา setup นาน

## 13.5 ข้อห้าม

| # | ข้อห้าม | เหตุผล |
|---|---|---|
| 1 | ❌ ห้าม import infrastructure ใน domain | ผิด Dependency Rule |
| 2 | ❌ ห้าม query โดยไม่มี tenant_id | ข้อมูลรั่ว |
| 3 | ❌ ห้าม commit secret | Security |
| 4 | ❌ ห้ามใช้ SELECT * | Performance |
| 5 | ❌ ห้ามข้าม Layer | Architecture |
| 6 | ❌ ห้าม deploy โดยไม่ test | Quality |
| 7 | ❌ ห้าม hardcode config | ยืดหยุ่น |
| 8 | ❌ ห้ามใช้ float กับเงิน | ใช้ decimal |

## 13.6 ตัวอย่างโค้ดที่รันได้จริง

### Docker Compose

```yaml
# docker-compose.yml
version: '3.9'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: mushroom_erp
      POSTGRES_USER: erp_user
      POSTGRES_PASSWORD: erp_pass
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
    depends_on:
      - zookeeper

  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    ports:
      - "2181:2181"
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181

  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://erp_user:erp_pass@postgres:5432/mushroom_erp?sslmode=disable
      REDIS_URL: redis://redis:6379
      KAFKA_BROKERS: kafka:9092
      JWT_SECRET: your-secret-key
    depends_on:
      - postgres
      - redis
      - kafka

volumes:
  postgres_data:
```

### รันระบบ

```bash
# Clone
git clone https://github.com/yourorg/smart-mushroom-erp.git
cd smart-mushroom-erp

# Start infrastructure
docker-compose up -d postgres redis kafka

# Run migrations
make migrate-up

# Start API
go run cmd/api/main.go

# Test
curl -X POST http://localhost:8080/api/v1/inventory/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "sku": "MUSH-001",
    "name": "เห็ดนางฟ้าเกรด A",
    "unit": "kg",
    "cost_price": 35,
    "sale_price": 80
  }'
```

---

# 14. คู่มือการทดสอบ

## 14.1 Checklist Test

| # | ประเภท Test | รายการ | เครื่องมือ |
|---|---|---|---|
| 1 | **Unit Test** | Domain logic | Go testing |
| 2 | **Unit Test** | Use case | Go testing + mock |
| 3 | **Integration Test** | Repository + DB | Testcontainers |
| 4 | **Integration Test** | Kafka | Testcontainers |
| 5 | **E2E Test** | API flow | Go + httptest |
| 6 | **Load Test** | 1,000 users | k6 |
| 7 | **Security Test** | OWASP Top 10 | OWASP ZAP |
| 8 | **Tenant Test** | Data isolation | Custom |
| 9 | **Contract Test** | API contract | Pact |
| 10 | **Chaos Test** | Failure injection | Chaos Mesh |

## 14.2 Unit Test ตัวอย่าง

```go
// internal/modules/erp/inventory/domain/entity/product_test.go
package entity_test

import (
    "testing"

    "github.com/shopspring/decimal"
    "github.com/stretchr/testify/assert"
    "github.com/yourorg/smart-mushroom-erp/internal/modules/erp/inventory/domain/entity"
)

func TestNewProduct(t *testing.T) {
    tests := []struct {
        name      string
        tenantID  string
        sku       string
        productNm string
        unit      string
        cost      decimal.Decimal
        sale      decimal.Decimal
        wantErr   bool
    }{
        {
            name:      "valid product",
            tenantID:  "tenant-1",
            sku:       "SKU-001",
            productNm: "เห็ดนางฟ้า",
            unit:      "kg",
            cost:      decimal.NewFromInt(35),
            sale:      decimal.NewFromInt(80),
            wantErr:   false,
        },
        {
            name:      "missing tenant",
            tenantID:  "",
            sku:       "SKU-001",
            productNm: "เห็ดนางฟ้า",
            unit:      "kg",
            cost:      decimal.NewFromInt(35),
            sale:      decimal.NewFromInt(80),
            wantErr:   true,
        },
        {
            name:      "negative price",
            tenantID:  "tenant-1",
            sku:       "SKU-001",
            productNm: "เห็ดนางฟ้า",
            unit:      "kg",
            cost:      decimal.NewFromInt(-35),
            sale:      decimal.NewFromInt(80),
            wantErr:   true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            p, err := entity.NewProduct(tt.tenantID, tt.sku, tt.productNm, tt.unit, tt.cost, tt.sale)
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, p)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, p)
                assert.Equal(t, tt.sku, p.SKU)
                assert.Equal(t, tt.tenantID, p.TenantID)
            }
        })
    }
}

func TestProduct_CalculateMargin(t *testing.T) {
    p, _ := entity.NewProduct("t1", "SKU-001", "Test", "kg",
        decimal.NewFromInt(35), decimal.NewFromInt(80))

    margin := p.CalculateMargin()
    // (80 - 35) / 80 * 100 = 56.25
    assert.Equal(t, "56.25", margin.String())
}

func TestProduct_IsLowStock(t *testing.T) {
    p, _ := entity.NewProduct("t1", "SKU-001", "Test", "kg",
        decimal.NewFromInt(35), decimal.NewFromInt(80))
    p.MinStock = decimal.NewFromInt(100)

    assert.True(t, p.IsLowStock(decimal.NewFromInt(50)))
    assert.True(t, p.IsLowStock(decimal.NewFromInt(100)))
    assert.False(t, p.IsLowStock(decimal.NewFromInt(150)))
}
```

## 14.3 Integration Test ตัวอย่าง

```go
// tests/integration/inventory_test.go
package integration

import (
    "context"
    "database/sql"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/yourorg/smart-mushroom-erp/internal/modules/erp/inventory/domain/entity"
    "github.com/yourorg/smart-mushroom-erp/internal/modules/erp/inventory/infrastructure/persistence/postgres"
)

func setupPostgres(t *testing.T) (*sql.DB, func()) {
    ctx := context.Background()

    // Start PostgreSQL container
    // เริ่ม container PostgreSQL
    pgContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15-alpine"),
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
    )
    require.NoError(t, err)

    // Get connection string
    // ดึง connection string
    connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
    require.NoError(t, err)

    // Connect
    // เชื่อมต่อ
    db, err := sql.Open("postgres", connStr)
    require.NoError(t, err)

    // Run migration
    // รัน migration
    runMigration(t, db)

    cleanup := func() {
        db.Close()
        pgContainer.Terminate(ctx)
    }

    return db, cleanup
}

func TestProductRepository_SaveAndFind(t *testing.T) {
    db, cleanup := setupPostgres(t)
    defer cleanup()

    repo := postgres.NewProductRepository(db)
    ctx := context.Background()

    // Create product
    // สร้างสินค้า
    p, err := entity.NewProduct("tenant-1", "SKU-001", "เห็ดนางฟ้า", "kg",
        decimal.NewFromInt(35), decimal.NewFromInt(80))
    require.NoError(t, err)

    // Save
    // บันทึก
    err = repo.Save(ctx, "tenant-1", p)
    require.NoError(t, err)

    // Find
    // ค้นหา
    found, err := repo.FindByID(ctx, "tenant-1", p.ID)
    require.NoError(t, err)
    assert.Equal(t, p.ID, found.ID)
    assert.Equal(t, p.SKU, found.SKU)

    // Tenant isolation test
    // ทดสอบการแยก tenant
    _, err = repo.FindByID(ctx, "tenant-2", p.ID)
    assert.ErrorIs(t, err, postgres.ErrNotFound)
}
```

---

# 15. คู่มือการใช้งาน

## 15.1 การเริ่มต้นใช้งาน

### 15.1.1 สำหรับผู้ใช้ทั่วไป

1. **สมัครสมาชิก** ที่ https://app.mushroom-erp.com/register
2. **ยืนยันอีเมล** ผ่านลิงก์ที่ส่งไป
3. **สร้างองค์กร** (Tenant) ครั้งแรก
4. **เพิ่มโรงเพาะ** ในเมนู Settings → Farms
5. **เพิ่มสินค้า** ในเมนู Inventory → Products
6. **เริ่มบันทึก** การผลิต/ขาย

### 15.1.2 สำหรับ Admin

```bash
# 1. Login
curl -X POST https://api.mushroom-erp.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@farm.com","password":"password"}'

# Response: { "token": "eyJ...", "refresh_token": "..." }

# 2. ใช้ token
export TOKEN="eyJ..."

# 3. สร้างสินค้า
curl -X POST https://api.mushroom-erp.com/api/v1/inventory/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "sku": "MUSH-001",
    "name": "เห็ดนางฟ้าเกรด A",
    "unit": "kg",
    "cost_price": 35,
    "sale_price": 80
  }'
```

## 15.2 เมนูหลัก

| เมนู | หน้าที่ |
|---|---|
| **Dashboard** | ภาพรวม KPI |
| **Inventory** | สต็อกสินค้า |
| **Production** | วางแผนผลิต |
| **Sales** | ขาย/ออเดอร์ |
| **CRM** | ลูกค้า |
| **Finance** | บัญชี |
| **Reports** | รายงาน |
| **Settings** | ตั้งค่า |

## 15.3 การจัดการผู้ใช้

| Role | สิทธิ์ |
|---|---|
| **Owner** | ทุกอย่าง |
| **Admin** | จัดการระบบ |
| **Manager** | ดู/อนุมัติ |
| **Staff** | บันทึกข้อมูล |
| **Viewer** | ดูอย่างเดียว |

---

# 16. คู่มือการบำรุงรักษา

## 16.1 Daily Maintenance

| เวลา | งาน | เครื่องมือ |
|---|---|---|
| 06:00 | ตรวจสอบ health check | `/health` |
| 07:00 | ดู error log | Grafana/Loki |
| 08:00 | ตรวจสอบ KPI | Grafana |
| 12:00 | ตรวจสอบ disk usage | `df -h` |
| 18:00 | Backup database | pg_dump |
| 23:00 | ตรวจสอบ slow query | pg_stat_statements |

## 16.2 Weekly Maintenance

- ✅ Update dependencies (Dependabot)
- ✅ Review security alerts
- ✅ Check certificate expiry
- ✅ Review capacity planning
- ✅ Test backup restore

## 16.3 Monthly Maintenance

- ✅ Patch OS
- ✅ Rotate secrets
- ✅ Review access logs
- ✅ Update documentation
- ✅ Load test

## 16.4 Backup Strategy

```bash
# Daily backup
pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME -F c -f backup_$(date +%Y%m%d).dump

# Upload to S3
aws s3 cp backup_$(date +%Y%m%d).dump s3://mushroom-erp-backup/

# Retention: 7 daily, 4 weekly, 12 monthly
```

## 16.5 Monitoring

| Metric | Alert Threshold |
|---|---|
| CPU | > 80% (5 min) |
| Memory | > 85% (5 min) |
| Disk | > 90% |
| Error Rate | > 1% |
| P95 Latency | > 500ms |
| DB Connections | > 80% pool |

---

# 17. คู่มือการขยาย/แก้ไข

## 17.1 การเพิ่ม Module ใหม่

```bash
# 1. สร้างโครงสร้าง
mkdir -p internal/modules/{prefix}/{module}/{domain,application,infrastructure,interfaces}

# 2. สร้างไฟล์พื้นฐาน
touch internal/modules/{prefix}/{module}/module.go
touch internal/modules/{prefix}/{module}/domain/entity/{aggregate}.go
touch internal/modules/{prefix}/{module}/domain/repository/{aggregate}_repository.go
touch internal/modules/{prefix}/{module}/application/command/create_{aggregate}.go
touch internal/modules/{prefix}/{module}/infrastructure/persistence/postgres/{aggregate}_repo_impl.go
touch internal/modules/{prefix}/{module}/interfaces/http/handler/{aggregate}_handler.go

# 3. Register ใน main.go
```

## 17.2 การเพิ่ม Field ใหม่

```go
// 1. เพิ่มใน Entity
type Product struct {
    // ... existing fields
    NewField string `json:"new_field"`
}

// 2. เพิ่มใน Migration
// migrations/20240115120000_add_new_field_to_products.up.sql
ALTER TABLE products ADD COLUMN new_field VARCHAR(255);

// 3. เพิ่มใน DTO
type CreateProductRequest struct {
    // ...
    NewField string `json:"new_field"`
}

// 4. อัปเดต Repository
// ...
```

## 17.3 การเปลี่ยน Database

```go
// ใช้ Repository Pattern ทำให้เปลี่ยน DB ได้ง่าย
// เพียง implement interface ใหม่

// MySQL implementation
type ProductRepositoryMySQL struct {
    db *sql.DB
}

func (r *ProductRepositoryMySQL) FindByID(ctx context.Context, tenantID, id string) (*entity.Product, error) {
    // MySQL-specific implementation
}
```

## 17.4 การ Scale Out

```
Phase 1: Vertical Scaling
  └── เพิ่ม CPU/RAM

Phase 2: Read Replica
  ├── Master (Write)
  └── Replica (Read) × N

Phase 3: Sharding by tenant_id
  ├── Shard 1: tenant A-M
  └── Shard 2: tenant N-Z

Phase 4: Microservices
  ├── Inventory Service
  ├── Production Service
  └── ...
```

---

# 18. Git Flow

## 18.1 Branching Model

```
┌─────────────────────────────────────────────────────────────┐
│                      GIT FLOW                                │
└─────────────────────────────────────────────────────────────┘

main (production)
  │
  ├─── hotfix/urgent-fix ──────────────┐
  │                                     │
  └─── develop (staging)               │
        │                               │
        ├─── feature/inventory ────────┤
        │                               │
        ├─── feature/crm ──────────────┤
        │                               │
        ├─── release/v1.2.0 ───────────┘
        │
        └─── bugfix/fix-stock ─────────┘
```

## 18.2 Branch Types

| Branch | จาก | Merge ไป | ชื่อ |
|---|---|---|---|
| **main** | - | - | `main` |
| **develop** | main | main | `develop` |
| **feature** | develop | develop | `feature/{module}-{desc}` |
| **release** | develop | main + develop | `release/v{X.Y.Z}` |
| **hotfix** | main | main + develop | `hotfix/{desc}` |
| **bugfix** | develop | develop | `bugfix/{desc}` |

## 18.3 Commit Convention

```
<type>(<scope>): <subject>

<body>

<footer>
```

| Type | ความหมาย |
|---|---|
| `feat` | ฟีเจอร์ใหม่ |
| `fix` | แก้บั๊ก |
| `docs` | เอกสาร |
| `style` | จัดรูปแบบ |
| `refactor` | Refactor |
| `test` | Test |
| `chore` | งานทั่วไป |
| `perf` | Performance |

### ตัวอย่าง

```
feat(inventory): add product creation API

- Add CreateProductUseCase
- Add ProductRepository implementation
- Add HTTP handler with validation

Closes #123
```

## 18.4 Workflow

```bash
# 1. Start feature
git checkout develop
git pull origin develop
git checkout -b feature/inventory-create-product

# 2. Develop
git add .
git commit -m "feat(inventory): add product entity"

# 3. Push
git push origin feature/inventory-create-product

# 4. Create PR to develop
gh pr create --base develop --title "feat: inventory product" --body "..."

# 5. After merge, delete branch
git branch -d feature/inventory-create-product
```

---

# 19. Code Review & PR

## 19.1 PR Template

```markdown
<!-- .github/pull_request_template.md -->
## 📋 Description
<!-- อธิบายการเปลี่ยนแปลง -->

## 🎯 Type of Change
- [ ] 🐛 Bug fix
- [ ] ✨ New feature
- [ ] 💥 Breaking change
- [ ] 📝 Documentation
- [ ] ♻️ Refactoring
- [ ] ⚡ Performance

## 🔗 Related Issues
Closes #

## 🧪 How Has This Been Tested?
- [ ] Unit tests
- [ ] Integration tests
- [ ] Manual testing

## 📸 Screenshots (if applicable)

## ✅ Checklist
- [ ] Code follows style guide
- [ ] Self-review done
- [ ] Comments added
- [ ] Documentation updated
- [ ] No new warnings
- [ ] Tests added
- [ ] All tests pass
- [ ] Tenant isolation verified
- [ ] Security reviewed
- [ ] Performance considered

## 🔍 Reviewer Notes
```

## 19.2 Code Review Checklist

| # | หัวข้อ | รายการ |
|---|---|---|
| 1 | **Architecture** | ถูก Layer? Dependency Rule? |
| 2 | **Domain** | Business logic ถูก? |
| 3 | **Tenant** | มี tenant_id ทุก query? |
| 4 | **Security** | Validate input? SQL injection? |
| 5 | **Error Handling** | Handle ครบ? |
| 6 | **Testing** | Coverage ≥ 80%? |
| 7 | **Performance** | N+1? Index? |
| 8 | **Naming** | ชื่อสื่อความหมาย? |
| 9 | **Comments** | อธิบาย 2 ภาษา? |
| 10 | **Documentation** | อัปเดต? |

## 19.3 Review Process

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ Author   │───▶│ CI/CD    │───▶│ Reviewer │───▶│ Merge    │
│ Creates  │    │ Tests    │    │ Approves │    │          │
│ PR       │    │ Lint     │    │          │    │          │
└──────────┘    └──────────┘    └──────────┘    └──────────┘
                     │                │
                     ▼                ▼
                ┌──────────┐    ┌──────────┐
                │ Fail →   │    │ Request  │
                │ Fix      │    │ Changes  │
                └──────────┘    └──────────┘
```

## 19.4 Review Guidelines

**ผู้เขียน:**
- PR เล็ก (< 400 บรรทัด)
- อธิบายชัดเจน
- Self-review ก่อน
- ตอบ comment ทุกอัน

**ผู้รีวิว:**
- รีวิวภายใน 24 ชม.
- Comment สร้างสรรค์
- Approve เมื่อพร้อม
- Test ใน local (ถ้าจำเป็น)

---

# 20. CI/CD

## 20.1 Pipeline Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    CI/CD PIPELINE                            │
└─────────────────────────────────────────────────────────────┘

Push/PR
   │
   ▼
┌──────────┐
│  Lint    │───▶ golangci-lint
└────┬─────┘
     ▼
┌──────────┐
│  Test    │───▶ go test -race -cover
└────┬─────┘
     ▼
┌──────────┐
│  Build   │───▶ go build
└────┬─────┘
     ▼
┌──────────┐
│ Security │───▶ Trivy, gosec
└────┬─────┘
     ▼
┌──────────┐
│  Docker  │───▶ Build & Push
└────┬─────┘
     ▼
┌──────────┐
│ Deploy   │───▶ Staging → Production
│  Dev     │
└──────────┘
```

## 20.2 GitHub Actions

```yaml
# .github/workflows/ci.yml
name: CI/CD

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

env:
  GO_VERSION: '1.21'
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest

  test:
    name: Test
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_PASSWORD: test
          POSTGRES_DB: testdb
        ports: ['5432:5432']
      redis:
        image: redis:7-alpine
        ports: ['6379:6379']
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.out

  security:
    name: Security Scan
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run Trivy
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          scan-ref: '.'
          format: 'sarif'
          output: 'trivy-results.sarif'
      - name: Run gosec
        uses: securego/gosec@master
        with:
          args: './...'

  build:
    name: Build & Push
    runs-on: ubuntu-latest
    needs: [lint, test, security]
    if: github.event_name == 'push'
    steps:
      - uses: actions/checkout@v4
      - uses: docker/setup-buildx-action@v3
      - uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: |
            ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }}
            ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max

  deploy-staging:
    name: Deploy Staging
    runs-on: ubuntu-latest
    needs: build
    if: github.ref == 'refs/heads/develop'
    environment: staging
    steps:
      - uses: actions/checkout@v4
      - name: Deploy to K8s
        run: |
          kubectl set image deployment/api \
            api=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }} \
            -n staging

  deploy-production:
    name: Deploy Production
    runs-on: ubuntu-latest
    needs: build
    if: github.ref == 'refs/heads/main'
    environment: production
    steps:
      - uses: actions/checkout@v4
      - name: Deploy to K8s
        run: |
          kubectl set image deployment/api \
            api=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }} \
            -n production
```

## 20.3 Dockerfile

```dockerfile
# Multi-stage build
# Build แบบ multi-stage

# Stage 1: Builder
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
# ติดตั้ง dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod
# คัดลอก go mod
COPY go.mod go.sum ./
RUN go mod download

# Copy source
# คัดลอก source
COPY . .

# Build
# Build
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.version=$(git describe --tags --always)" \
    -o /app/api ./cmd/api

# Stage 2: Runtime
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
# สร้าง user ที่ไม่ใช่ root
RUN adduser -D -u 1000 appuser

WORKDIR /app

# Copy binary
# คัดลอก binary
COPY --from=builder /app/api .
COPY --from=builder /app/migrations ./migrations

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["./api"]
```

---

# 21. Root Cause Analysis (RCA)

## 21.1 ขั้นตอนการทำ RCA

### 21.1.1 ระบุปัญหา (Define the Problem)

**คำถาม:**
- เกิดอะไรขึ้น?
- เมื่อไหร่?
- ที่ไหน?
- ใครได้รับผลกระทบ?
- รุนแรงแค่ไหน?

**ตัวอย่าง:**
```
เมื่อเวลา 14:30 น. วันที่ 15 ม.ค. 2024
API /api/v1/inventory/products ตอบสนองช้า
P95 latency = 5 วินาที (ปกติ 200ms)
ผู้ใช้ 500 รายได้รับผลกระทบ
```

### 21.1.2 รวบรวมข้อมูล (Collect Data)

| แหล่งข้อมูล | รายละเอียด |
|---|---|
| **Logs** | Application logs, Access logs |
| **Metrics** | CPU, Memory, DB connections |
| **Traces** | Distributed tracing (Jaeger) |
| **Events** | Deployment, Config changes |
| **Users** | Reports จากผู้ใช้ |

### 21.1.3 ระบุสาเหตุที่เป็นไปได้ (Identify Possible Causes)

**ใช้ Fishbone Diagram:**

```
                    ┌── คน ──┐
                    │ ไม่มี training
                    │ ลา
                    └────────┘
┌── วิธี ──┐        │        ┌── เครื่อง ──┐
│ ไม่มี index│       │        │ CPU 100%   │
│ Query ซับซ้อน│─────┼────────│ RAM ไม่พอ  │
└──────────┘        │        └────────────┘
                    │
              ┌─────┴─────┐
              │  ปัญหา:   │
              │ API ช้า    │
              └─────┬─────┘
                    │
┌── วัสดุ ──┐        │        ┌── สภาพแวดล้อม ──┐
│ DB disk I/O│      │        │ Network latency  │
│ Network    │──────┼────────│ Cloud provider   │
└──────────┘        │        └──────────────────┘
                    │
              ┌── การวัด ──┐
              │ ไม่มี monitoring
              │ ไม่มี alert
              └────────────┘
```

### 21.1.4 หา "สาเหตุหลัก" (Find the Root Cause)

**ใช้ 5 Whys:**

```
1. ทำไม API ช้า?
   → เพราะ query ใช้เวลา 4 วินาที

2. ทำไม query ช้า?
   → เพราะ scan ทั้งตาราง (full table scan)

3. ทำไมต้อง scan ทั้งตาราง?
   → เพราะไม่มี index บน tenant_id

4. ทำไมไม่มี index?
   → เพราะ migration ไม่ได้สร้าง index

5. ทำไม migration ไม่สร้าง index?
   → เพราะ developer ลืม (Root Cause)
```

### 21.1.5 วางแผนและแก้ไข (Implement Solution)

| ระดับ | การแก้ไข |
|---|---|
| **Immediate** | เพิ่ม index ทันที |
| **Short-term** | แก้ migration, เพิ่ม test |
| **Long-term** | Code review checklist, monitoring |
| **Preventive** | Automated index detection |

### 21.1.6 ติดตามผล (Monitor)

- ✅ ตรวจสอบ P95 latency กลับมา < 200ms
- ✅ เพิ่ม alert สำหรับ query > 1s
- ✅ Review migration ใน PR
- ✅ เพิ่ม index ใน CI/CD

## 21.2 RCA Template

```markdown
# Root Cause Analysis Report

## 1. Incident Summary
- **Date**: 2024-01-15 14:30
- **Duration**: 45 นาที
- **Severity**: P1 (Critical)
- **Impact**: 500 users

## 2. Timeline
| เวลา | เหตุการณ์ |
|---|---|
| 14:30 | Alert: P95 > 5s |
| 14:35 | ทีมเริ่ม investigate |
| 14:45 | พบ full table scan |
| 15:00 | เพิ่ม index |
| 15:15 | ระบบกลับปกติ |

## 3. Root Cause
Missing index on `products.tenant_id` ทำให้ full table scan

## 4. 5 Whys
1. ทำไมช้า? → Query 4s
2. ทำไม query ช้า? → Full scan
3. ทำไม full scan? → ไม่มี index
4. ทำไมไม่มี index? → Migration ลืม
5. ทำไมลืม? → ไม่มี checklist

## 5. Action Items
| # | การแก้ไข | ผู้รับผิดชอบ | กำหนดเสร็จ |
|---|---|---|---|
| 1 | เพิ่ม index | DBA | 15 ม.ค. |
| 2 | แก้ migration | Dev | 16 ม.ค. |
| 3 | เพิ่ม checklist | Lead | 17 ม.ค. |
| 4 | เพิ่ม monitoring | DevOps | 18 ม.ค. |

## 6. Lessons Learned
- ต้องมี index review ใน PR
- ต้องมี query performance test
- ต้องมี alert สำหรับ slow query

## 7. Prevention
- Automated index detection
- Query plan analysis ใน CI
- Performance budget ใน test
```

---

# 📎 ภาคผนวก

## A. เอกสารอ้างอิง

- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design - Eric Evans](https://domainlanguage.com/ddd/)
- [Go Standard Project Layout](https://github.com/golang-standards/project-layout)
- [12 Factor App](https://12factor.net/)

## B. เครื่องมือที่ใช้

| ประเภท | เครื่องมือ |
|---|---|
| **ภาษา** | Go 1.21+ |
| **Web Framework** | Gin |
| **Database** | PostgreSQL 15 |
| **Cache** | Redis 7 |
| **Message Queue** | Kafka |
| **Time-series** | ClickHouse |
| **Search** | Elasticsearch |
| **Container** | Docker |
| **Orchestration** | Kubernetes |
| **CI/CD** | GitHub Actions |
| **Monitoring** | Prometheus + Grafana |
| **Tracing** | Jaeger |
| **Logging** | Loki |
| **Testing** | Testify, Testcontainers |
| **Load Test** | k6 |

## C. ติดต่อ

- **Repository**: https://github.com/yourorg/smart-mushroom-erp
- **Documentation**: https://docs.mushroom-erp.com
- **Support**: support@mushroom-erp.com

---

**เวอร์ชัน**: 1.0  
**วันที่**: 2024-01-15  
**ผู้จัดทำ**: Architecture Team  
**ประเภทเอกสาร**: Application Design Document

---

*เอกสารนี้เป็นต้นแบบสำหรับให้ AI สร้าง/แก้ไขโปรแกรมภาษา Go ทุกประเภท ใช้โครงสร้าง Clean Architecture + DDD พร้อม Layer ที่ชัดเจน, Cross-cutting Concerns ครบ, และ Pattern แยก Modules*