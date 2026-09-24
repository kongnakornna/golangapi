# 📘 Application Design Document (Extended Edition)
## ระบบ ERP+SaaS สำหรับฟาร์มเห็ดอัจฉริยะ (Smart Farm ERP) และโรงงาน SME

---

# สารบัญ (Table of Contents)

| # | หัวข้อ | หน้า |
|---|---|---|
| 1 | [บทนำ (Introduction) — Business Model Canvas](#1-บทนำ-introduction) | 1 |
| 2 | [บทนิยาม (Definitions)](#2-บทนิยาม-definitions) | 4 |
| 3 | [บทหัวข้อ (Topics)](#3-บทหัวข้อ-topics) | 6 |
| 4 | [การออกแบบ Workflow](#4-การออกแบบ-workflow) | 10 |
| 5 | [Case Study](#5-case-study) | 16 |
| 6 | [โครงสร้างโฟลเดอร์](#6-โครงสร้างโฟลเดอร์) | 20 |
| 7 | [หลักการทำงาน (Concept)](#7-หลักการทำงาน-concept) | 25 |
| 8 | [Workflow และ Dataflow](#8-workflow-และ-dataflow) | 30 |
| 9 | [Code Template](#9-code-template) | 36 |
| 10 | [Checklist Module](#10-checklist-module) | 55 |
| 11 | [Security Code](#11-security-code) | 58 |
| 12 | [Load Test](#12-load-test) | 63 |
| 13 | [สรุป](#13-สรุป) | 67 |
| 14 | [คู่มือการทดสอบ](#14-คู่มือการทดสอบ) | 72 |
| 15 | [คู่มือการใช้งาน](#15-คู่มือการใช้งาน) | 76 |
| 16 | [คู่มือการบำรุงรักษา](#16-คู่มือการบำรุงรักษา) | 80 |
| 17 | [คู่มือการขยาย/แก้ไข](#17-คู่มือการขยายแก้ไข) | 84 |
| 18 | [Git Flow](#18-git-flow) | 88 |
| 19 | [Code Review & PR](#19-code-review--pr) | 92 |
| 20 | [CI/CD](#20-cicd) | 96 |
| 21 | [Root Cause Analysis (RCA)](#21-root-cause-analysis-rca) | 101 |
| 22 | [SME Factory Extension](#22-sme-factory-extension) | 105 |
| 23 | [Dual-Domain Architecture](#23-dual-domain-architecture) | 112 |
| 24 | [Cross-Domain Integration](#24-cross-domain-integration) | 118 |

---

# 1. บทนำ (Introduction)

## 1.1 Business Model Canvas — Unified Platform

### 🎯 BMC: Smart Farm ERP + SME Factory ERP

| องค์ประกอบ | ฟาร์มเห็ดอัจฉริยะ | โรงงาน SME |
|---|---|---|
| **Key Partners** | ผู้ผลิตเซ็นเซอร์ IoT, ซัพพลายเออร์ขี้เลื่อย/หัวเชื้อ, ขนส่ง (Kerry/Flash), อย./GAP | ซัพพลายเออร์วัตถุดิบ, ผู้รับจ้างผลิต, ขนส่ง, ธนาคาร, กรมโรงงาน |
| **Key Activities** | พัฒนาแพลตฟอร์ม ERP+SaaS, IoT Integration, AI Forecast, Support 24/7 | พัฒนา ERP สำหรับการผลิต, MES Integration, Quality Control, Compliance |
| **Key Resources** | ทีม Go, Cloud Infrastructure, โมเดล AI, ฐานข้อมูลเกษตรกร | ทีม Go, Cloud, ML Models, ฐานข้อมูลโรงงาน |
| **Value Propositions** | 1) ERP ครบวงจรฟาร์มเห็ด 2) SaaS ฟรี/ราคาถูก 3) IoT ควบคุมโรงเพาะ 4) AI พยากรณ์ผลผลิต 5) GAP/GMP/อย. | 1) ERP ครบวงจรโรงงาน 2) MES + SCADA 3) Costing ต่อหน่วย 4) OEE Tracking 5) ISO 9001/14001 |
| **Customer Relationships** | Self-service, Community, Training ออนไลน์ | Dedicated CSM, On-site Training, 24/7 Support |
| **Channels** | Web, Mobile, LINE OA, Partner (สหกรณ์) | Web, Mobile, Direct Sales, System Integrator |
| **Customer Segments** | ฟาร์มเล็ก/กลาง/สหกรณ์/โรงงานแปรรูปเห็ด | โรงงานผลิต/แปรรูป/ประกอบ/OEM |
| **Cost Structure** | Cloud 30%, R&D 40%, S&M 15%, Support 10%, Admin 5% | Cloud 25%, R&D 35%, S&M 20%, Support 15%, Admin 5% |
| **Revenue Streams** | SaaS Subscription, IoT Hardware, AI Add-on, Setup Fee | SaaS Subscription, Implementation, Custom Module, Support Contract |

## 1.2 Dual-Domain Revenue Model

| Tier | ราคา/เดือน | ฟาร์มเห็ด | โรงงาน SME |
|---|---|---|---|
| **Free** | 0 บาท | 1 โรงเพาะ | - |
| **Starter** | 990 บาท | 10 โรงเพาะ + Inventory + POS | 1 สายการผลิต + BOM |
| **Pro** | 2,990 บาท | + Production + CRM + IoT | + MES + QC + Costing |
| **Business** | 9,990 บาท | + AI Forecast + Multi-farm | + OEE + SPC + Multi-line |
| **Enterprise** | 29,990 บาท | ไม่จำกัด + Custom + SLA | ไม่จำกัด + SCADA + MES Full |

## 1.3 Market Size (TAM/SAM/SOM)

```
┌─────────────────────────────────────────────────────────┐
│  TAM: 500,000 ฟาร์ม + 300,000 โรงงาน SME = 800,000     │
│  SAM: 50,000 ฟาร์ม + 30,000 โรงงาน (มี Smartphone)     │
│  SOM: 5,000 ฟาร์ม + 3,000 โรงงาน (3 ปีแรก)             │
└─────────────────────────────────────────────────────────┘
```

---

# 2. บทนิยาม (Definitions)

## 2.1 คำศัพท์ทั่วไป

| คำศัพท์ | ความหมาย |
|---|---|
| **ERP** | Enterprise Resource Planning — ระบบรวมกระบวนการธุรกิจ |
| **SaaS** | Software as a Service — ซอฟต์แวร์ให้บริการผ่าน Cloud |
| **Multi-tenant** | 1 ระบบให้บริการหลายองค์กร แยกข้อมูล |
| **Clean Architecture** | แยก Layer: domain → application → infrastructure → interfaces |
| **DDD** | Domain-Driven Design |
| **Aggregate Root** | Entity หลักควบคุม Entity อื่น |
| **Value Object** | Object ไม่มี identity, immutable |
| **Repository** | Interface เข้าถึง Data Layer |
| **CQRS** | Command Query Responsibility Segregation |
| **Event-Driven** | สื่อสารผ่าน Event |

## 2.2 คำศัพท์เฉพาะฟาร์มเห็ด

| คำศัพท์ | ความหมาย |
|---|---|
| **ก้อนเชื้อ** | Substrate Block — วัสดุเพาะเห็ด |
| **บ่มเชื้อ** | Incubation — ระยะเส้นใยเจริญ |
| **เปิดดอก** | Fruiting — ระยะออกดอก |
| **โรงเพาะ** | Growing House |
| **OEE** | Overall Equipment Effectiveness |
| **GAP** | Good Agricultural Practice |
| **GMP** | Good Manufacturing Practice |

## 2.3 คำศัพท์เฉพาะโรงงาน SME

| คำศัพท์ | ความหมาย |
|---|---|
| **BOM** | Bill of Materials — สูตรการผลิต |
| **MES** | Manufacturing Execution System |
| **SCADA** | Supervisory Control and Data Acquisition |
| **SPC** | Statistical Process Control |
| **OEE** | Availability × Performance × Quality |
| **WIP** | Work In Progress |
| **MRP** | Material Requirements Planning |
| **ISO 9001** | มาตรฐานระบบคุณภาพ |
| **CAPA** | Corrective and Preventive Action |
| **Andon** | ระบบแจ้งเตือนบนสายการผลิต |

## 2.4 คำศัพท์ทางเทคนิค

| คำศัพท์ | ความหมาย |
|---|---|
| **Kafka** | Distributed Event Streaming |
| **MQTT** | Message Queuing Telemetry Transport (IoT) |
| **ClickHouse** | Column-oriented DB สำหรับ Analytics |
| **Redis** | In-memory Cache |
| **Prometheus** | Monitoring System |
| **Grafana** | Visualization |
| **Jaeger** | Distributed Tracing |

---

# 3. บทหัวข้อ (Topics)

## 3.1 Architecture Overview — Unified Platform

```
┌──────────────────────────────────────────────────────────────────┐
│                     PRESENTATION LAYER                            │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────┐ │
│  │ Web App  │ │ Mobile   │ │ POS      │ │ IoT Dash │ │ Andon  │ │
│  │ (React)  │ │ (Flutter)│ │ (Tablet) │ │ (Web)    │ │ (TV)   │ │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘ └───┬────┘ │
└───────┼────────────┼────────────┼────────────┼───────────┼──────┘
        └────────────┴──────┬─────┴────────────┴───────────┘
                            ▼
┌──────────────────────────────────────────────────────────────────┐
│              API GATEWAY (Kong/Nginx)                             │
│  Auth JWT │ Rate Limit │ Tenant Resolve │ Tracing │ Logging       │
└──────────────────────────┬───────────────────────────────────────┘
                           ▼
┌──────────────────────────────────────────────────────────────────┐
│           APPLICATION LAYER (Go Modular Monolith)                 │
│                                                                   │
│  ┌───────────────────── FARM DOMAIN ─────────────────────┐       │
│  │ erp/inventory │ erp/production │ erp/finance          │       │
│  │ crm/crm │ pos/pos │ ecommerce │ logistics             │       │
│  │ iot/iot │ forecast │ reporting │ KPI/Analytics        │       │
│  └───────────────────────────────────────────────────────┘       │
│                                                                   │
│  ┌───────────────────── FACTORY DOMAIN ──────────────────┐       │
│  │ mes/production │ mes/bom │ mes/routing │ mes/workorder │       │
│  │ qms/quality │ qms/spc │ qms/capa │ qms/andon          │       │
│  │ costing/standard │ costing/actual │ costing/variance  │       │
│  │ scada/gateway │ scada/plc │ scada/hmi                │       │
│  │ maintenance/pm │ maintenance/cm │ maintenance/spare   │       │
│  └───────────────────────────────────────────────────────┘       │
│                                                                   │
│  ┌───────────────────── SHARED CORE ─────────────────────┐       │
│  │ auth │ tenant │ user │ notification │ document │ workflow│      │
│  └───────────────────────────────────────────────────────┘       │
└──────────────────────────┬───────────────────────────────────────┘
                           ▼
┌──────────────────────────────────────────────────────────────────┐
│                  INFRASTRUCTURE LAYER                             │
│  PostgreSQL │ Redis │ Kafka │ ClickHouse │ MQTT │ MinIO          │
│  Elasticsearch │ ML Python │ OPC-UA │ Modbus TCP │ Prometheus   │
└──────────────────────────────────────────────────────────────────┘
```

## 3.2 วัตถุประสงค์ (Objectives)

### 3.2.1 วัตถุประสงค์ร่วม (Shared)
1. **รวมระบบ**: รวมทุกกระบวนการใน 1 แพลตฟอร์ม
2. **ลดต้นทุน**: SaaS ลดต้นทุน 80%
3. **Multi-tenant**: 1 ระบบให้บริการหลายองค์กร
4. **ขยายได้**: Modular Monolith → Microservices
5. **Real-time**: ข้อมูลทันที
6. **Compliance**: GAP/GMP/ISO/อย.

### 3.2.2 วัตถุประสงค์ฟาร์มเห็ด
1. **เพิ่มผลผลิต**: IoT + AI เพิ่ม 30%
2. **ลดของเสีย**: Traceability 100%
3. **ควบคุมคุณภาพ**: เกรด A/B/C อัตโนมัติ
4. **วางแผน**: AI Forecast ผลผลิต

### 3.2.3 วัตถุประสงค์โรงงาน SME
1. **เพิ่ม OEE**: จาก 60% → 85%
2. **ลดต้นทุน**: Costing ต่อหน่วยแม่นยำ
3. **ลดของเสีย**: SPC ควบคุมกระบวนการ
4. **Compliance**: ISO 9001/14001
5. **MES**: ติดตามการผลิต real-time

## 3.3 กลุ่มเป้าหมาย (Target Users)

| กลุ่ม | ลักษณะ | Module หลัก |
|---|---|---|
| **ฟาร์มเห็ดเล็ก** | 1-10 โรง | Inventory, POS, Reporting |
| **ฟาร์มเห็ดกลาง** | 10-50 โรง | + Production, CRM, IoT |
| **สหกรณ์การเกษตร** | 50+ โรง | + AI, Finance, Multi-farm |
| **โรงงานแปรรูปเห็ด** | GMP/อย. | + Quality, Batch, GMP |
| **โรงงาน SME เล็ก** | 1-5 สายการผลิต | MES, BOM, QC |
| **โรงงาน SME กลาง** | 5-20 สายการผลิต | + SPC, OEE, Costing |
| **โรงงาน OEM** | รับจ้างผลิต | + Multi-customer, Traceability |

## 3.4 ความรู้พื้นฐาน (Prerequisites)

| ด้าน | รายการ |
|---|---|
| **ภาษา** | Go 1.21+, SQL, JavaScript/TypeScript, Python (ML) |
| **Framework** | Gin, GORM/sqlx, React, Flutter |
| **Database** | PostgreSQL 15, Redis 7, ClickHouse |
| **Messaging** | Kafka, MQTT, gRPC |
| **DevOps** | Docker, Kubernetes, GitHub Actions |
| **Architecture** | Clean Architecture, DDD, CQRS, Event-Driven |
| **Domain** | ERP, MES, SCADA, GMP, ISO 9001 |
| **Security** | OWASP Top 10, JWT, RBAC |

## 3.5 เนื้อหาโดยย่อ

| Module | วัตถุประสงค์ | ประโยชน์ |
|---|---|---|
| **erp/inventory** | สต็อกฟาร์ม | ลดของเสีย |
| **erp/production** | ผลิตก้อนเชื้อ | เพิ่มผลผลิต |
| **erp/finance** | บัญชี | รู้กำไร |
| **crm/crm** | ลูกค้า | เพิ่มยอดขาย |
| **pos/pos** | ขายหน้าร้าน | ตัดสต็อก |
| **ecommerce** | ขายออนไลน์ | เข้าถึงกว้าง |
| **logistics** | ขนส่ง | ส่งตรงเวลา |
| **iot/iot** | เซ็นเซอร์ | ควบคุมโรงเพาะ |
| **forecast** | AI | วางแผน |
| **reporting** | รายงาน | ตัดสินใจ |
| **mes/production** | MES | OEE ↑ |
| **mes/bom** | BOM | ต้นทุนแม่น |
| **qms/quality** | QC | ลดของเสีย |
| **qms/spc** | SPC | ควบคุมกระบวนการ |
| **costing** | ต้นทุน | กำไรจริง |
| **scada** | SCADA | ควบคุมเครื่องจักร |
| **maintenance** | ซ่อมบำรุง | ลด downtime |

---

# 4. การออกแบบ Workflow

## 4.1 Workflow ฟาร์มเห็ด: วงจรชีวิตก้อนเชื้อ → ขาย

```
┌──────────────────────────────────────────────────────────────────┐
│              MUSHROOM FARM LIFECYCLE WORKFLOW                     │
└──────────────────────────────────────────────────────────────────┘

[1] จัดซื้อวัสดุ         [2] ผลิตก้อนเชื้อ      [3] บ่มเชื้อ
    │                       │                     │
    ▼                       ▼                     ▼
┌─────────┐           ┌─────────────┐       ┌─────────────┐
│Procure- │           │ Production  │       │ Incubation  │
│ment     │──────────▶│ Order       │──────▶│ (21-30 วัน) │
└─────────┘           └─────────────┘       └──────┬──────┘
    │                       │                      │
    │ PR→PO→GR              │ BOM→WO→QC            │
    ▼                       ▼                      ▼
┌─────────┐           ┌─────────────┐       ┌─────────────┐
│Supplier │           │ Inventory   │       │ Stock       │
│Payment  │           │ (Raw Mat)   │       │ (Incubated) │
└─────────┘           └─────────────┘       └──────┬──────┘
                                                    │
[4] เปิดดอก              [5] เก็บเกี่ยว        [6] แพ็ค/ขาย
    │                       │                     │
    ▼                       ▼                     ▼
┌─────────────┐       ┌─────────────┐       ┌─────────────┐
│ IoT Control │       │ Harvest     │       │ QC + Grade  │
│ (Temp/Hum)  │──────▶│ Record      │──────▶│ A/B/C       │
└─────────────┘       └─────────────┘       └──────┬──────┘
                                                    │
                       ┌────────────────────────────┼──────────┐
                       ▼                            ▼          ▼
                 ┌──────────┐              ┌──────────┐ ┌────────┐
                 │ POS      │              │E-commerce│ │Direct  │
                 └────┬─────┘              └────┬─────┘ └───┬────┘
                      └──────────┬──────────────┴───────────┘
                                 ▼
                          ┌─────────────┐
                          │ Logistics   │
                          └──────┬──────┘
                                 ▼
                          ┌─────────────┐
                          │ Finance     │
                          └──────┬──────┘
                                 ▼
                          ┌─────────────┐
                          │ Reporting   │
                          └─────────────┘
```

## 4.2 Workflow โรงงาน SME: Order-to-Delivery

```
┌──────────────────────────────────────────────────────────────────┐
│           FACTORY SME ORDER-TO-DELIVERY WORKFLOW                  │
└──────────────────────────────────────────────────────────────────┘

[1] รับออเดอร์           [2] วางแผนผลิต        [3] จัดซื้อวัสดุ
    │                       │                     │
    ▼                       ▼                     ▼
┌─────────┐           ┌─────────────┐       ┌─────────────┐
│Sales    │           │ MRP         │       │ Purchase    │
│Order    │──────────▶│ Calculation │──────▶│ Order       │
└─────────┘           └─────────────┘       └──────┬──────┘
    │                       │                      │
    │                       │ BOM Explosion        │
    ▼                       ▼                      ▼
┌─────────┐           ┌─────────────┐       ┌─────────────┐
│Contract │           │ Work Order  │       │ Goods       │
│Review   │           │ (WO)        │       │ Receipt     │
└─────────┘           └──────┬──────┘       └──────┬──────┘
                             │                      │
[4] ผลิต                    │                      │
    │                       ▼                      │
    ▼                 ┌─────────────┐              │
┌─────────────┐       │ MES         │              │
│ SCADA/PLC   │◀─────▶│ Execution   │◀─────────────┘
│ (Machine)   │       │ (Real-time) │
└──────┬──────┘       └──────┬──────┘
       │                     │
       │                     ▼
       │              ┌─────────────┐
       │              │ SPC/QC      │
       │              │ (In-process)│
       │              └──────┬──────┘
       │                     │
       ▼                     ▼
┌─────────────┐       ┌─────────────┐
│ OEE Track   │       │ Final QC    │
└─────────────┘       └──────┬──────┘
                             │
[5] บรรจุ/ส่ง                │
    │                       ▼
    ▼                 ┌─────────────┐
┌─────────────┐       │ Packing     │
│ Warehouse   │◀──────│ + Label     │
└──────┬──────┘       └──────┬──────┘
       │                     │
       ▼                     ▼
┌─────────────┐       ┌─────────────┐
│ Shipping    │       │ Invoice     │
└──────┬──────┘       └──────┬──────┘
       │                     │
       └──────────┬──────────┘
                  ▼
           ┌─────────────┐
           │ Costing     │
           │ (Actual)    │
           └──────┬──────┘
                  ▼
           ┌─────────────┐
           │ Variance    │
           │ Analysis    │
           └─────────────┘
```

## 4.3 Workflow OEE Monitoring (โรงงาน)

```
┌──────────────────────────────────────────────────────────────────┐
│                    OEE MONITORING WORKFLOW                        │
└──────────────────────────────────────────────────────────────────┘

Machine Sensors (PLC/SCADA)
      │
      ▼ (Modbus/OPC-UA)
┌─────────────────┐
│ SCADA Gateway   │
│ (Go Service)    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Kafka: machine.  │
│ telemetry        │
└────────┬────────┘
         │
         ├──────────────────────┐
         ▼                      ▼
┌─────────────────┐    ┌─────────────────┐
│ Time-series DB  │    │ OEE Calculator  │
│ (ClickHouse)    │    │ (Go Service)    │
└─────────────────┘    └────────┬────────┘
                                │
                    ┌───────────┼───────────┐
                    ▼           ▼           ▼
              ┌──────────┐ ┌──────────┐ ┌──────────┐
              │Availability│ │Performance│ │Quality  │
              │ = RT/PT   │ │ = AT/RT  │ │ = GT/AT  │
              └─────┬─────┘ └─────┬────┘ └─────┬────┘
                    └─────────────┼────────────┘
                                  ▼
                          ┌───────────────┐
                          │ OEE = A×P×Q   │
                          └───────┬───────┘
                                  │
                    ┌─────────────┼─────────────┐
                    ▼             ▼             ▼
              ┌──────────┐  ┌──────────┐  ┌──────────┐
              │Dashboard │  │Andon     │  │Alert     │
              │(Grafana) │  │(TV)      │  │(LINE)    │
              └──────────┘  └──────────┘  └──────────┘
```

## 4.4 Workflow Cross-Domain: ฟาร์ม → โรงงาน → ลูกค้า

```
┌──────────────────────────────────────────────────────────────────┐
│              FARM → FACTORY → CUSTOMER INTEGRATION                │
└──────────────────────────────────────────────────────────────────┘

┌──────────────┐                    ┌──────────────┐
│  FARM        │                    │  FACTORY     │
│  (เห็ดสด)     │                    │  (แปรรูป)     │
└──────┬───────┘                    └──────┬───────┘
       │                                    │
       │ 1. เก็บเกี่ยว                      │
       ▼                                    │
┌──────────────┐                            │
│ Harvest      │                            │
│ Record       │                            │
└──────┬───────┘                            │
       │                                    │
       │ 2. ส่งเห็ดสด                        │
       │    (Kafka: farm.harvest.completed) │
       └────────────────────────────────────┤
                                            ▼
                                    ┌──────────────┐
                                    │ Raw Material │
                                    │ Receipt      │
                                    └──────┬───────┘
                                           │
                                           │ 3. ผลิต
                                           ▼
                                    ┌──────────────┐
                                    │ Production   │
                                    │ (MES)        │
                                    └──────┬───────┘
                                           │
                                           │ 4. QC
                                           ▼
                                    ┌──────────────┐
                                    │ Quality      │
                                    │ Check        │
                                    └──────┬───────┘
                                           │
                                           │ 5. บรรจุ
                                           ▼
                                    ┌──────────────┐
                                    │ Finished     │
                                    │ Goods        │
                                    └──────┬───────┘
                                           │
                                           │ 6. ส่งลูกค้า
                                           ▼
                                    ┌──────────────┐
                                    │ Customer     │
                                    │ Delivery     │
                                    └──────────────┘
```

---

# 5. Case Study

## 5.1 Case Study ฟาร์มเห็ด: "ภูเขียว" จ.ชัยภูมิ

### ปัญหาก่อนใช้ระบบ
- ฟาร์มขนาด 30 โรงเพาะ
- จดบันทึกในสมุด → ข้อมูลหาย, คำนวณต้นทุนผิด
- ขาดสต็อกบ่อย → ผลผลิตลด 20%
- ส่งของช้า → ลูกค้ายกเลิก 15%
- ไม่รู้ต้นทุนจริงต่อก้อน → ตั้งราคาผิด

### แนวทางแก้ไข

| ปัญหา | Module | ผลลัพธ์ |
|---|---|---|
| บันทึกในสมุด | erp/inventory + Mobile | Real-time |
| ขาดสต็อก | Reorder + AI Forecast | ลด 90% |
| ส่งช้า | logistics + IoT | ตรงเวลา 98% |
| ไม่รู้ต้นทุน | erp/finance + Costing | รู้ต้นทุน/ก้อน |
| ลูกค้าหนี | crm + POS | ยอดซ้ำ +40% |

### ผลลัพธ์ 6 เดือน
- ผลผลิต **+35%** (500 → 675 กก./วัน)
- ต้นทุน **-22%** (45 → 35 บาท/กก.)
- กำไร **+60%** (150K → 240K บาท/เดือน)
- ลูกค้าประจำ **×3**

---

## 5.2 Case Study โรงงาน SME: "เมืองเลยฟู้ดส์" แปรรูปเห็ด

### ปัญหาก่อนใช้ระบบ
- โรงงานขนาดกลาง 8 สายการผลิต
- OEE เพียง 58%
- ของเสีย 12% ต่อล็อต
- ต้นทุนต่อหน่วยไม่แม่นยำ (±15%)
- GMP/อย. เอกระจัดกระจาย
- Downtime เครื่องจักรสูง 20 ชม./เดือน

### แนวทางแก้ไข

| ปัญหา | Module | ผลลัพธ์ |
|---|---|---|
| OEE ต่ำ | mes/production + scada | OEE 58% → 84% |
| ของเสียสูง | qms/spc + qms/quality | ของเสีย 12% → 3% |
| ต้นทุนไม่แม่น | costing/standard + actual | ±15% → ±2% |
| GMP | qms/capa + audit_trail | ผ่าน GMP 3 เดือน |
| Downtime | maintenance/pm | 20 → 4 ชม./เดือน |

### ผลลัพธ์ 6 เดือน
- OEE **+45%** (58% → 84%)
- ของเสีย **-75%** (12% → 3%)
- ต้นทุน **-18%**
- กำไร **+55%**
- Downtime **-80%**

---

## 5.3 Case Study Cross-Domain: "สหกรณ์เห็ดภาคอีสาน" → "โรงงานแปรรูป"

### สถานการณ์
- สหกรณ์ 200 ฟาร์ม ส่งเห็ดสดให้โรงงานแปรรูป
- ปัญหา: traceability, ของเสีย, การวางแผน

### แนวทางแก้ไข

| ระดับ | Module | ผลลัพธ์ |
|---|---|---|
| ฟาร์ม | iot + forecast | พยากรณ์ผลผลิต |
| สหกรณ์ | erp/inventory + logistics | รวมสต็อก |
| โรงงาน | mes + qms | รับวัตถุดิบ |
| ทั้งระบบ | Cross-domain Kafka | Real-time sync |

### ผลลัพธ์
- Traceability 100% (farm → factory → customer)
- ของเสีย **-60%**
- วางแผนผลผลิตแม่นยำ **85%**
- กำไรสหกรณ์ **+40%**

---

## 5.4 Case Study: โรงงาน OEM "ไทยพรีซิชั่น"

### ปัญหา
- รับงาน OEM จาก 5 ลูกค้า
- ต้องแยกต้นทุนต่อลูกค้า
- Traceability ต่อ lot
- ISO 9001 audit

### แนวทางแก้ไข

| ปัญหา | Module | ผลลัพธ์ |
|---|---|---|
| ต้นทุนต่อลูกค้า | costing + multi-customer | แม่นยำ |
| Traceability | mes + qms/batch | 100% |
| ISO 9001 | qms + document | ผ่าน audit |
| OEE | scada + mes | +30% |

---

# 6. โครงสร้างโฟลเดอร์

## 6.1 Unified Folder Structure

```
smart-erp-platform/
├── cmd/
│   ├── api/main.go                  # API entry
│   ├── worker/main.go               # Background worker
│   ├── scada-gateway/main.go        # SCADA gateway
│   └── migrator/main.go             # DB migration
├── internal/
│   ├── modules/
│   │   │
│   │   ├── farm/                    # 🌱 FARM DOMAIN
│   │   │   ├── erp/
│   │   │   │   ├── inventory/       # ✅ พร้อม
│   │   │   │   ├── production/      # ⏳
│   │   │   │   └── finance/         # ⏳
│   │   │   ├── crm/crm/
│   │   │   ├── pos/pos/
│   │   │   ├── ecommerce/ecommerce/
│   │   │   ├── logistics/logistics/
│   │   │   ├── iot/iot/
│   │   │   ├── forecast/forecast/
│   │   │   └── reporting/reporting/
│   │   │
│   │   ├── factory/                 # 🏭 FACTORY DOMAIN
│   │   │   ├── mes/
│   │   │   │   ├── production/      # MES Core
│   │   │   │   ├── bom/             # Bill of Materials
│   │   │   │   ├── routing/         # Routing
│   │   │   │   ├── workorder/       # Work Order
│   │   │   │   └── scheduling/      # Scheduling
│   │   │   ├── qms/
│   │   │   │   ├── quality/         # QC
│   │   │   │   ├── spc/             # Statistical Process Control
│   │   │   │   ├── capa/            # CAPA
│   │   │   │   └── andon/           # Andon System
│   │   │   ├── costing/
│   │   │   │   ├── standard/        # Standard Cost
│   │   │   │   ├── actual/          # Actual Cost
│   │   │   │   └── variance/        # Variance Analysis
│   │   │   ├── scada/
│   │   │   │   ├── gateway/         # SCADA Gateway
│   │   │   │   ├── plc/             # PLC Integration
│   │   │   │   ├── opcua/           # OPC-UA
│   │   │   │   └── hmi/             # HMI
│   │   │   ├── maintenance/
│   │   │   │   ├── pm/              # Preventive
│   │   │   │   ├── cm/              # Corrective
│   │   │   │   └── spare/           # Spare Parts
│   │   │   └── warehouse/           # Factory Warehouse
│   │   │
│   │   └── shared/                  # 🤝 SHARED CORE
│   │       ├── auth/
│   │       ├── tenant/
│   │       ├── user/
│   │       ├── notification/
│   │       ├── document/
│   │       ├── workflow/
│   │       └── audit/
│   │
│   ├── shared/                      # Shared utilities
│   │   ├── config/
│   │   ├── database/
│   │   ├── logger/
│   │   ├── middleware/
│   │   ├── errors/
│   │   ├── events/
│   │   └── utils/
│   │
│   └── platform/                    # Platform services
│       ├── auth/
│       ├── tenant/
│       ├── storage/
│       └── messaging/
├── pkg/                             # Public packages
│   ├── uuid/
│   ├── decimal/
│   ├── validator/
│   └── oee/                         # OEE calculator
├── migrations/
│   ├── farm/
│   └── factory/
├── deployments/
│   ├── docker/
│   ├── k8s/
│   └── terraform/
├── docs/
├── tests/
├── scripts/
├── .env.example
├── docker-compose.yml
├── Makefile
├── go.mod
└── README.md
```

## 6.2 Module Structure มาตรฐาน (Farm + Factory)

```
internal/modules/{domain}/{prefix}/{module}/
├── domain/                          # 🎯 Business Logic
│   ├── entity/
│   ├── value_object/
│   ├── repository/
│   ├── service/
│   ├── event/
│   └── errors/
├── application/                     # 🔄 Use Cases
│   ├── command/
│   ├── query/
│   ├── dto/
│   └── port/
├── infrastructure/                  # 🔧 Implementation
│   ├── persistence/
│   │   ├── postgres/
│   │   ├── redis/
│   │   ├── clickhouse/
│   │   └── timeseries/
│   ├── messaging/
│   │   ├── kafka_producer.go
│   │   ├── kafka_consumer.go
│   │   └── mqtt_client.go
│   ├── scada/                       # สำหรับ factory เท่านั้น
│   │   ├── opcua_client.go
│   │   └── modbus_client.go
│   └── external/
├── interfaces/                      # 🌐 API Layer
│   ├── http/
│   │   ├── handler/
│   │   ├── request/
│   │   ├── response/
│   │   └── routes.go
│   ├── grpc/
│   └── middleware/
├── module.go
└── README.md
```

## 6.3 ตัวอย่างโครงสร้าง MES Production

```
internal/modules/factory/mes/production/
├── domain/
│   ├── entity/
│   │   ├── production_order.go      # Aggregate Root
│   │   ├── work_order.go
│   │   ├── operation.go
│   │   ├── work_center.go
│   │   └── machine.go
│   ├── value_object/
│   │   ├── order_status.go
│   │   ├── operation_status.go
│   │   ├── oee.go                   # OEE Value Object
│   │   └── cycle_time.go
│   ├── repository/
│   │   ├── production_order_repository.go
│   │   ├── work_order_repository.go
│   │   └── work_center_repository.go
│   ├── service/
│   │   ├── production_service.go
│   │   ├── oee_calculator.go        # Domain Service
│   │   └── scheduling_service.go
│   ├── event/
│   │   ├── production_started.go
│   │   ├── production_completed.go
│   │   └── oee_updated.go
│   └── errors/
│       └── errors.go
├── application/
│   ├── command/
│   │   ├── create_production_order.go
│   │   ├── start_production.go
│   │   ├── record_output.go
│   │   └── complete_production.go
│   ├── query/
│   │   ├── get_oee.go
│   │   └── get_production_status.go
│   └── dto/
├── infrastructure/
│   ├── persistence/
│   │   ├── postgres/
│   │   ├── clickhouse/              # Time-series OEE
│   │   └── redis/
│   ├── messaging/
│   │   ├── kafka_producer.go
│   │   └── kafka_consumer.go
│   └── scada/
│       ├── opcua_client.go
│       └── modbus_client.go
├── interfaces/
│   ├── http/
│   │   ├── handler/
│   │   └── routes.go
│   └── websocket/
│       └── realtime_oee.go          # Real-time OEE
└── module.go
```

---

# 7. หลักการทำงาน (Concept)

## 7.1 Clean Architecture — Dual Domain

```
┌─────────────────────────────────────────────────────────────────┐
│  INTERFACES (HTTP/gRPC/WebSocket/MQTT)                           │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  APPLICATION (Use Cases, DTO, Ports)                       │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  DOMAIN (Entity, VO, Repository, Domain Service)    │  │  │
│  │  │  ┌───────────────────────────────────────────────┐  │  │  │
│  │  │  │  Business Rules (Pure Go, ไม่มี dependency)    │  │  │  │
│  │  │  └───────────────────────────────────────────────┘  │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────┘  │
│  INFRASTRUCTURE (DB, Cache, MQ, SCADA, External)                 │
└─────────────────────────────────────────────────────────────────┘
        ▲                                          ▲
        │  Dependency Rule:                        │
        │  Outer → Inner (ห้ามย้อนกลับ)              │
        │  Domain ห้าม import อะไรนอกจาก stdlib      │
```

## 7.2 Dependency Rule

| Layer | Import ได้ | Import ไม่ได้ |
|---|---|---|
| **Domain** | stdlib, pkg | application, infrastructure, interfaces |
| **Application** | domain, pkg | infrastructure, interfaces |
| **Infrastructure** | domain, application, pkg | interfaces |
| **Interfaces** | ทุก layer | - |

## 7.3 Multi-tenant + Multi-domain Pattern

```go
// ทุก Repository ต้องรับทั้ง tenant_id และ domain
type ProductionOrderRepository interface {
    FindByID(ctx context.Context, tenantID, id string) (*ProductionOrder, error)
    FindAll(ctx context.Context, tenantID string, filter Filter) ([]*ProductionOrder, error)
    Save(ctx context.Context, tenantID string, order *ProductionOrder) error
}

// Domain Context — แยก farm vs factory
type DomainContext struct {
    TenantID  string
    Domain    string  // "farm" หรือ "factory"
    UserID    string
    Role      string
}
```

## 7.4 Event-Driven Cross-Domain

```go
// Farm → Factory Event
type HarvestCompletedEvent struct {
    TenantID    string
    FarmID      string
    BatchID     string
    ProductType string
    Quantity    decimal.Decimal
    Quality     string  // A, B, C
    HarvestedAt time.Time
}

// Factory → Farm Event
type RawMaterialReceivedEvent struct {
    TenantID    string
    FactoryID   string
    BatchID     string
    Quantity    decimal.Decimal
    ReceivedAt  time.Time
}
```

## 7.5 OEE Calculation Pattern

```go
// OEE Value Object
type OEE struct {
    Availability decimal.Decimal  // Availability = Run Time / Planned Production Time
    Performance  decimal.Decimal  // Performance = (Ideal Cycle Time × Total Count) / Run Time
    Quality      decimal.Decimal  // Quality = Good Count / Total Count
    Value        decimal.Decimal  // OEE = Availability × Performance × Quality
}

// Calculate OEE
// คำนวณ OEE
func CalculateOEE(runTime, plannedTime, idealCycleTime, totalCount, goodCount decimal.Decimal) OEE {
    // Availability
    // ความพร้อมใช้งาน
    availability := runTime.Div(plannedTime)

    // Performance
    // ประสิทธิภาพ
    performance := idealCycleTime.Mul(totalCount).Div(runTime)

    // Quality
    // คุณภาพ
    quality := goodCount.Div(totalCount)

    // OEE
    // OEE รวม
    oee := availability.Mul(performance).Mul(quality)

    return OEE{
        Availability: availability,
        Performance:  performance,
        Quality:      quality,
        Value:        oee,
    }
}
```

## 7.6 SCADA Integration Pattern

```go
// SCADA Gateway — อ่านข้อมูลจาก PLC
type SCADAGateway interface {
    // ReadTag อ่านค่าจาก PLC tag
    ReadTag(ctx context.Context, tag string) (interface{}, error)

    // WriteTag เขียนค่าลง PLC tag
    WriteTag(ctx context.Context, tag string, value interface{}) error

    // SubscribeTags subscribe tags เพื่อรับค่า real-time
    SubscribeTags(ctx context.Context, tags []string, callback func(tag string, value interface{})) error
}

// OPC-UA Implementation
type OPCUAGateway struct {
    client *opcua.Client
}

func (g *OPCUAGateway) ReadTag(ctx context.Context, tag string) (interface{}, error) {
    // อ่านค่าจาก OPC-UA server
    nodeID := fmt.Sprintf("ns=2;s=%s", tag)
    return g.client.ReadNode(ctx, nodeID)
}
```

## 7.7 SPC (Statistical Process Control) Pattern

```go
// SPC Chart — ควบคุมกระบวนการด้วยสถิติ
type SPCChart struct {
    TenantID    string
    ProcessID   string
    Metric      string  // เช่น "diameter", "weight"
    Samples     []Sample
    UCL         decimal.Decimal  // Upper Control Limit
    LCL         decimal.Decimal  // Lower Control Limit
    CL          decimal.Decimal  // Center Line (Mean)
}

// CalculateControlLimits คำนวณ UCL/LCL
func (s *SPCChart) CalculateControlLimits() {
    // ค่าเฉลี่ย
    mean := s.calculateMean()

    // ค่าเบี่ยงเบนมาตรฐาน
    stdDev := s.calculateStdDev()

    // UCL = Mean + 3σ
    // LCL = Mean - 3σ
    threeSigma := stdDev.Mul(decimal.NewFromInt(3))
    s.UCL = mean.Add(threeSigma)
    s.LCL = mean.Sub(threeSigma)
    s.CL = mean
}

// IsOutOfControl ตรวจสอบว่ากระบวนการหลุดควบคุมหรือไม่
func (s *SPCChart) IsOutOfControl(sample Sample) bool {
    return sample.Value.GreaterThan(s.UCL) || sample.Value.LessThan(s.LCL)
}
```

---

# 8. Workflow และ Dataflow

## 8.1 Dataflow: Farm — Create Product

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

## 8.2 Dataflow: Factory — Production Execution

```
┌──────────────────────────────────────────────────────────────────┐
│              FACTORY PRODUCTION EXECUTION DATAFLOW                │
└──────────────────────────────────────────────────────────────────┘

[1] Sales Order → MRP
    └── คำนวณความต้องการวัสดุ

[2] Work Order Creation
    ├── BOM Explosion
    ├── Routing Selection
    └── Work Center Assignment

[3] MES Execution
    ├── Operator Login (Andon)
    ├── Scan Barcode (WO)
    ├── Machine Start (SCADA)
    └── Real-time Telemetry

[4] SCADA/PLC Data
    ├── Cycle Time
    ├── Good Count
    ├── Reject Count
    └── Downtime Reason

[5] Kafka: machine.telemetry
         │
         ├──────────────────────┐
         ▼                      ▼
[6a] ClickHouse            [6b] OEE Calculator
    Time-series                ├── Availability
                               ├── Performance
                               ├── Quality
                               └── OEE Value
                                        │
                                        ▼
                               [7] OEE Events
                                   ├── Kafka: oee.updated
                                   ├── WebSocket: real-time
                                   └── Andon Display

[8] SPC Check
    ├── Sample Collection
    ├── Control Chart
    └── Alert ถ้าหลุดควบคุม

[9] Production Completion
    ├── Final QC
    ├── Packing
    └── Finished Goods

[10] Costing
    ├── Standard Cost
    ├── Actual Cost
    └── Variance Analysis

[11] Reporting
    ├── OEE Dashboard
    ├── Cost Report
    └── Quality Report
```

## 8.3 Dataflow: Cross-Domain Farm → Factory

```
┌──────────────────────────────────────────────────────────────────┐
│                  CROSS-DOMAIN DATAFLOW                            │
└──────────────────────────────────────────────────────────────────┘

┌──────────────┐                    ┌──────────────┐
│  FARM        │                    │  FACTORY     │
│  Domain      │                    │  Domain      │
└──────┬───────┘                    └──────┬───────┘
       │                                    │
       │ 1. Harvest Completed              │
       ▼                                    │
┌──────────────┐                            │
│ Kafka:       │                            │
│ farm.harvest.│                            │
│ completed    │                            │
└──────┬───────┘                            │
       │                                    │
       │ 2. Event consumed                  │
       └────────────────────────────────────┤
                                            ▼
                                    ┌──────────────┐
                                    │ Raw Material │
                                    │ Receipt      │
                                    └──────┬───────┘
                                           │
                                           │ 3. Production
                                           ▼
                                    ┌──────────────┐
                                    │ Kafka:       │
                                    │ factory.     │
                                    │ production.  │
                                    │ completed    │
                                    └──────┬───────┘
                                           │
                                           │ 4. Event consumed
                                           ▼
                                    ┌──────────────┐
                                    │ Farm         │
                                    │ Dashboard    │
                                    │ (Traceability)│
                                    └──────────────┘

┌──────────────────────────────────────────────────────────────────┐
│  TRACEABILITY CHAIN: Farm → Factory → Customer                    │
│                                                                   │
│  Batch ID: FARM-2024-001                                          │
│    ├── Farm: ภูเขียว, วันที่เก็บ 2024-01-15                       │
│    ├── Factory: เมืองเลยฟู้ดส์, วันที่รับ 2024-01-16              │
│    ├── Production: WO-2024-001, วันที่ผลิต 2024-01-17            │
│    ├── QC: Passed, 2024-01-17                                    │
│    └── Customer: บริษัท ABC, วันที่ส่ง 2024-01-18                │
└──────────────────────────────────────────────────────────────────┘
```

## 8.4 Dataflow: IoT → Alert → Action (Farm)

```
┌──────────────────────────────────────────────────────────────────┐
│                    IoT DATAFLOW (FARM)                            │
└──────────────────────────────────────────────────────────────────┘

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

## 8.5 Dataflow: SCADA → OEE (Factory)

```
┌──────────────────────────────────────────────────────────────────┐
│                  SCADA → OEE DATAFLOW                             │
└──────────────────────────────────────────────────────────────────┘

[1] PLC (Siemens/Mitsubishi)
    ├── Machine State: RUNNING
    ├── Cycle Count: 1500
    ├── Good Count: 1450
    ├── Reject Count: 50
    └── Downtime: 15 min
         │
         ▼ OPC-UA / Modbus TCP
[2] SCADA Gateway (Go)
    ├── Poll every 1s
    ├── Parse tags
    └── Publish to Kafka
         │
         ▼
[3] Kafka: machine.telemetry
         │
         ▼
[4] OEE Calculator (Go Service)
    ├── Calculate Availability
    │   └── RT / PT
    ├── Calculate Performance
    │   └── (ICT × TC) / RT
    ├── Calculate Quality
    │   └── GC / TC
    └── Calculate OEE
        └── A × P × Q
         │
         ▼
[5] Kafka: oee.updated
         │
         ├──────────────────────┐
         ▼                      ▼
[6a] ClickHouse            [6b] WebSocket
    Time-series OEE            Real-time push
         │                      │
         ▼                      ▼
[7a] Grafana               [7b] Andon TV
    Dashboard                  Display
```

---

# 9. Code Template

## 9.1 Shared: Tenant Middleware

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
    TenantIDKey ContextKey = "tenant_id" // รหัสองค์กร
    UserIDKey   ContextKey = "user_id"   // รหัสผู้ใช้
    RoleKey     ContextKey = "role"      // บทบาท
    DomainKey   ContextKey = "domain"    // farm หรือ factory
)

// TenantMiddleware extracts tenant_id from JWT and injects into context.
// TenantMiddleware ดึง tenant_id จาก JWT และใส่ใน context
func TenantMiddleware(jwtSecret string) gin.HandlerFunc {
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
            return []byte(jwtSecret), nil
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
        domain, _ := claims["domain"].(string) // "farm" หรือ "factory"

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
        ctx = context.WithValue(ctx, DomainKey, domain)

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

// GetDomain extracts domain from context.
// GetDomain ดึง domain จาก context
func GetDomain(c *gin.Context) string {
    if v, ok := c.Request.Context().Value(DomainKey).(string); ok {
        return v
    }
    return ""
}
```

## 9.2 Farm: Product Entity (Aggregate Root)

```go
// internal/modules/farm/erp/inventory/domain/entity/product.go
package entity

import (
    "errors"
    "time"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)

// Product is the Aggregate Root for farm inventory.
// Product เป็น Aggregate Root สำหรับสินค้าคงคลังฟาร์ม
type Product struct {
    ID          string          // รหัสสินค้า
    TenantID    string          // รหัสองค์กร
    SKU         string          // รหัส SKU
    Name        string          // ชื่อสินค้า
    Description string          // คำอธิบาย
    Category    string          // หมวดหมู่ (เห็ดสด/เห็ดแห้ง/ก้อนเชื้อ)
    Unit        string          // หน่วยนับ
    CostPrice   decimal.Decimal // ราคาทุน
    SalePrice   decimal.Decimal // ราคาขาย
    MinStock    decimal.Decimal // สต็อกขั้นต่ำ
    MaxStock    decimal.Decimal // สต็อกสูงสุด
    ShelfLife   int             // อายุการเก็บ (วัน)
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

// CalculateMargin calculates the profit margin.
// CalculateMargin คำนวณกำไรขั้นต้น
func (p *Product) CalculateMargin() decimal.Decimal {
    if p.CostPrice.IsZero() {
        return decimal.Zero
    }
    return p.SalePrice.Sub(p.CostPrice).Div(p.SalePrice).Mul(decimal.NewFromInt(100))
}

// IsLowStock checks if the stock is below minimum.
// IsLowStock ตรวจสอบว่าสต็อกต่ำกว่าขั้นต่ำหรือไม่
func (p *Product) IsLowStock(currentStock decimal.Decimal) bool {
    return currentStock.LessThanOrEqual(p.MinStock)
}
```

## 9.3 Factory: Production Order Entity

```go
// internal/modules/factory/mes/production/domain/entity/production_order.go
package entity

import (
    "errors"
    "time"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)

// ProductionOrder is the Aggregate Root for factory production.
// ProductionOrder เป็น Aggregate Root สำหรับการผลิตโรงงาน
type ProductionOrder struct {
    ID           string              // รหัสใบสั่งผลิต
    TenantID     string              // รหัสองค์กร
    OrderNumber  string              // เลขที่ใบสั่งผลิต
    ProductID    string              // รหัสสินค้า
    BOMID        string              // รหัส BOM
    Quantity     decimal.Decimal     // จำนวนที่ผลิต
    Unit         string              // หน่วยนับ
    Status       OrderStatus         // สถานะ
    Priority     Priority            // ความสำคัญ
    PlannedStart time.Time           // วางแผนเริ่ม
    PlannedEnd   time.Time           // วางแผนสิ้นสุด
    ActualStart  *time.Time          // เริ่มจริง
    ActualEnd    *time.Time          // สิ้นสุดจริง
    WorkOrders   []*WorkOrder        // ใบสั่งงาน
    CreatedAt    time.Time           // วันที่สร้าง
    UpdatedAt    time.Time           // วันที่แก้ไข
    Version      int                 // Optimistic Lock
}

// OrderStatus represents the status of a production order.
// OrderStatus แทนสถานะของใบสั่งผลิต
type OrderStatus string

const (
    OrderStatusDraft       OrderStatus = "draft"       // แบบร่าง
    OrderStatusPlanned     OrderStatus = "planned"     // วางแผนแล้ว
    OrderStatusReleased    OrderStatus = "released"    // ปล่อยผลิต
    OrderStatusInProgress  OrderStatus = "in_progress" // กำลังผลิต
    OrderStatusCompleted   OrderStatus = "completed"   // เสร็จ
    OrderStatusCancelled   OrderStatus = "cancelled"   // ยกเลิก
)

// Priority represents the priority of a production order.
// Priority แทนความสำคัญของใบสั่งผลิต
type Priority string

const (
    PriorityLow    Priority = "low"    // ต่ำ
    PriorityNormal Priority = "normal" // ปกติ
    PriorityHigh   Priority = "high"   // สูง
    PriorityUrgent Priority = "urgent" // เร่งด่วน
)

// NewProductionOrder creates a new ProductionOrder.
// NewProductionOrder สร้าง ProductionOrder ใหม่
func NewProductionOrder(
    tenantID, orderNumber, productID, bomID string,
    quantity decimal.Decimal, unit string,
    plannedStart, plannedEnd time.Time,
) (*ProductionOrder, error) {
    // Validate tenant ID
    // ตรวจสอบรหัสองค์กร
    if tenantID == "" {
        return nil, errors.New("tenant_id is required")
    }

    // Validate order number
    // ตรวจสอบเลขที่ใบสั่งผลิต
    if orderNumber == "" {
        return nil, errors.New("order_number is required")
    }

    // Validate quantity
    // ตรวจสอบจำนวน
    if quantity.LessThanOrEqual(decimal.Zero) {
        return nil, errors.New("quantity must be positive")
    }

    // Validate dates
    // ตรวจสอบวันที่
    if plannedEnd.Before(plannedStart) {
        return nil, errors.New("planned_end must be after planned_start")
    }

    now := time.Now()
    return &ProductionOrder{
        ID:           uuid.NewString(),
        TenantID:     tenantID,
        OrderNumber:  orderNumber,
        ProductID:    productID,
        BOMID:        bomID,
        Quantity:     quantity,
        Unit:         unit,
        Status:       OrderStatusDraft,
        Priority:     PriorityNormal,
        PlannedStart: plannedStart,
        PlannedEnd:   plannedEnd,
        CreatedAt:    now,
        UpdatedAt:    now,
        Version:      1,
    }, nil
}

// Release releases the production order for execution.
// Release ปล่อยใบสั่งผลิตเพื่อดำเนินการ
func (po *ProductionOrder) Release() error {
    // Check status
    // ตรวจสอบสถานะ
    if po.Status != OrderStatusPlanned {
        return errors.New("only planned orders can be released")
    }

    po.Status = OrderStatusReleased
    po.UpdatedAt = time.Now()
    po.Version++
    return nil
}

// Start starts the production order.
// Start เริ่มการผลิต
func (po *ProductionOrder) Start() error {
    // Check status
    // ตรวจสอบสถานะ
    if po.Status != OrderStatusReleased {
        return errors.New("only released orders can be started")
    }

    now := time.Now()
    po.Status = OrderStatusInProgress
    po.ActualStart = &now
    po.UpdatedAt = now
    po.Version++
    return nil
}

// Complete completes the production order.
// Complete เสร็จสิ้นการผลิต
func (po *ProductionOrder) Complete() error {
    // Check status
    // ตรวจสอบสถานะ
    if po.Status != OrderStatusInProgress {
        return errors.New("only in-progress orders can be completed")
    }

    now := time.Now()
    po.Status = OrderStatusCompleted
    po.ActualEnd = &now
    po.UpdatedAt = now
    po.Version++
    return nil
}

// CalculateLeadTime calculates the actual lead time in hours.
// CalculateLeadTime คำนวณเวลานำจริงเป็นชั่วโมง
func (po *ProductionOrder) CalculateLeadTime() decimal.Decimal {
    if po.ActualStart == nil || po.ActualEnd == nil {
        return decimal.Zero
    }
    duration := po.ActualEnd.Sub(*po.ActualStart)
    return decimal.NewFromFloat(duration.Hours())
}
```

## 9.4 Factory: OEE Value Object

```go
// internal/modules/factory/mes/production/domain/value_object/oee.go
package value_object

import (
    "github.com/shopspring/decimal"
)

// OEE represents Overall Equipment Effectiveness.
// OEE แทนประสิทธิภาพโดยรวมของเครื่องจักร
type OEE struct {
    Availability decimal.Decimal // ความพร้อมใช้งาน (0-1)
    Performance  decimal.Decimal // ประสิทธิภาพ (0-1)
    Quality      decimal.Decimal // คุณภาพ (0-1)
    Value        decimal.Decimal // OEE รวม (0-1)
}

// NewOEE creates a new OEE value object.
// NewOEE สร้าง OEE value object ใหม่
func NewOEE(availability, performance, quality decimal.Decimal) OEE {
    // Clamp values to 0-1
    // จำกัดค่าให้อยู่ 0-1
    availability = clamp01(availability)
    performance = clamp01(performance)
    quality = clamp01(quality)

    // OEE = A × P × Q
    // OEE = A × P × Q
    value := availability.Mul(performance).Mul(quality)

    return OEE{
        Availability: availability,
        Performance:  performance,
        Quality:      quality,
        Value:        value,
    }
}

// CalculateAvailability calculates availability.
// CalculateAvailability คำนวณความพร้อมใช้งาน
// Availability = Run Time / Planned Production Time
// ความพร้อมใช้งาน = เวลาที่เดินเครื่อง / เวลาที่วางแผนผลิต
func CalculateAvailability(runTime, plannedTime decimal.Decimal) decimal.Decimal {
    if plannedTime.IsZero() {
        return decimal.Zero
    }
    return runTime.Div(plannedTime)
}

// CalculatePerformance calculates performance.
// CalculatePerformance คำนวณประสิทธิภาพ
// Performance = (Ideal Cycle Time × Total Count) / Run Time
// ประสิทธิภาพ = (เวลารอบมาตรฐาน × จำนวนผลผลิต) / เวลาที่เดินเครื่อง
func CalculatePerformance(idealCycleTime, totalCount, runTime decimal.Decimal) decimal.Decimal {
    if runTime.IsZero() {
        return decimal.Zero
    }
    return idealCycleTime.Mul(totalCount).Div(runTime)
}

// CalculateQuality calculates quality.
// CalculateQuality คำนวณคุณภาพ
// Quality = Good Count / Total Count
// คุณภาพ = จำนวนดี / จำนวนทั้งหมด
func CalculateQuality(goodCount, totalCount decimal.Decimal) decimal.Decimal {
    if totalCount.IsZero() {
        return decimal.Zero
    }
    return goodCount.Div(totalCount)
}

// AsPercentage returns OEE as percentage.
// AsPercentage คืนค่า OEE เป็นเปอร์เซ็นต์
func (o OEE) AsPercentage() decimal.Decimal {
    return o.Value.Mul(decimal.NewFromInt(100))
}

// WorldClass checks if OEE is world-class (≥ 85%).
// WorldClass ตรวจสอบว่า OEE ระดับโลกหรือไม่ (≥ 85%)
func (o OEE) WorldClass() bool {
    return o.Value.GreaterThanOrEqual(decimal.NewFromFloat(0.85))
}

// clamp01 clamps a decimal to [0, 1].
// clamp01 จำกัด decimal ให้อยู่ [0, 1]
func clamp01(d decimal.Decimal) decimal.Decimal {
    zero := decimal.Zero
    one := decimal.NewFromInt(1)
    if d.LessThan(zero) {
        return zero
    }
    if d.GreaterThan(one) {
        return one
    }
    return d
}
```

## 9.5 Factory: OEE Calculator Service

```go
// internal/modules/factory/mes/production/domain/service/oee_calculator.go
package service

import (
    "context"
    "time"

    "github.com/shopspring/decimal"
    "github.com/yourorg/smart-erp/internal/modules/factory/mes/production/domain/value_object"
)

// MachineTelemetry represents raw machine data.
// MachineTelemetry แทนข้อมูลดิบจากเครื่องจักร
type MachineTelemetry struct {
    TenantID        string          // รหัสองค์กร
    MachineID       string          // รหัสเครื่องจักร
    Timestamp       time.Time       // เวลาที่บันทึก
    State           string          // สถานะ: RUNNING, IDLE, DOWN
    RunTime         decimal.Decimal // เวลาที่เดินเครื่อง (นาที)
    PlannedTime     decimal.Decimal // เวลาที่วางแผนผลิต (นาที)
    IdealCycleTime  decimal.Decimal // เวลารอบมาตรฐาน (นาที/ชิ้น)
    TotalCount      decimal.Decimal // จำนวนผลผลิตทั้งหมด
    GoodCount       decimal.Decimal // จำนวนผลผลิตดี
    RejectCount     decimal.Decimal // จำนวนของเสีย
    DowntimeMinutes decimal.Decimal // เวลาหยุด (นาที)
}

// OEECalculator calculates OEE from telemetry.
// OEECalculator คำนวณ OEE จาก telemetry
type OEECalculator struct {
    // Dependencies
    // Dependencies
}

// NewOEECalculator creates a new OEE calculator.
// NewOEECalculator สร้าง OEE calculator ใหม่
func NewOEECalculator() *OEECalculator {
    return &OEECalculator{}
}

// Calculate calculates OEE from telemetry.
// Calculate คำนวณ OEE จาก telemetry
func (c *OEECalculator) Calculate(ctx context.Context, t MachineTelemetry) value_object.OEE {
    // Step 1: Calculate Availability
    // ขั้นตอนที่ 1: คำนวณความพร้อมใช้งาน
    availability := value_object.CalculateAvailability(t.RunTime, t.PlannedTime)

    // Step 2: Calculate Performance
    // ขั้นตอนที่ 2: คำนวณประสิทธิภาพ
    performance := value_object.CalculatePerformance(t.IdealCycleTime, t.TotalCount, t.RunTime)

    // Step 3: Calculate Quality
    // ขั้นตอนที่ 3: คำนวณคุณภาพ
    quality := value_object.CalculateQuality(t.GoodCount, t.TotalCount)

    // Step 4: Create OEE value object
    // ขั้นตอนที่ 4: สร้าง OEE value object
    return value_object.NewOEE(availability, performance, quality)
}

// CalculateBatch calculates OEE for multiple machines.
// CalculateBatch คำนวณ OEE สำหรับหลายเครื่อง
func (c *OEECalculator) CalculateBatch(ctx context.Context, telemetries []MachineTelemetry) map[string]value_object.OEE {
    results := make(map[string]value_object.OEE)
    for _, t := range telemetries {
        results[t.MachineID] = c.Calculate(ctx, t)
    }
    return results
}
```

## 9.6 Factory: SPC Chart

```go
// internal/modules/factory/qms/spc/domain/value_object/spc_chart.go
package value_object

import (
    "math"
    "time"

    "github.com/shopspring/decimal"
)

// Sample represents a single measurement sample.
// Sample แทนตัวอย่างการวัด 1 ครั้ง
type Sample struct {
    Timestamp time.Time       // เวลาที่วัด
    Value     decimal.Decimal // ค่าที่วัด
    Operator  string          // ผู้วัด
}

// SPCChart represents a Statistical Process Control chart.
// SPCChart แทนแผนภูมิควบคุมกระบวนการเชิงสถิติ
type SPCChart struct {
    TenantID  string          // รหัสองค์กร
    ProcessID string          // รหัสกระบวนการ
    Metric    string          // ชื่อตัวชี้วัด (เช่น "diameter")
    Unit      string          // หน่วย
    Samples   []Sample        // ตัวอย่างทั้งหมด
    UCL       decimal.Decimal // Upper Control Limit
    LCL       decimal.Decimal // Lower Control Limit
    CL        decimal.Decimal // Center Line (Mean)
    USL       decimal.Decimal // Upper Specification Limit
    LSL       decimal.Decimal // Lower Specification Limit
    UpdatedAt time.Time       // เวลาที่อัปเดต
}

// NewSPCChart creates a new SPC chart.
// NewSPCChart สร้าง SPC chart ใหม่
func NewSPCChart(tenantID, processID, metric, unit string, usl, lsl decimal.Decimal) *SPCChart {
    return &SPCChart{
        TenantID:  tenantID,
        ProcessID: processID,
        Metric:    metric,
        Unit:      unit,
        USL:       usl,
        LSL:       lsl,
        Samples:   []Sample{},
        UpdatedAt: time.Now(),
    }
}

// AddSample adds a sample and recalculates control limits.
// AddSample เพิ่มตัวอย่างและคำนวณ control limits ใหม่
func (c *SPCChart) AddSample(s Sample) {
    c.Samples = append(c.Samples, s)
    c.calculateControlLimits()
    c.UpdatedAt = time.Now()
}

// calculateControlLimits calculates UCL, LCL, CL.
// calculateControlLimits คำนวณ UCL, LCL, CL
func (c *SPCChart) calculateControlLimits() {
    if len(c.Samples) < 2 {
        return
    }

    // Calculate mean
    // คำนวณค่าเฉลี่ย
    sum := decimal.Zero
    for _, s := range c.Samples {
        sum = sum.Add(s.Value)
    }
    mean := sum.Div(decimal.NewFromInt(int64(len(c.Samples))))
    c.CL = mean

    // Calculate standard deviation
    // คำนวณค่าเบี่ยงเบนมาตรฐาน
    varianceSum := decimal.Zero
    for _, s := range c.Samples {
        diff := s.Value.Sub(mean)
        varianceSum = varianceSum.Add(diff.Mul(diff))
    }
    variance := varianceSum.Div(decimal.NewFromInt(int64(len(c.Samples) - 1)))
    stdDev := decimal.NewFromFloat(math.Sqrt(variance.InexactFloat64()))

    // UCL = Mean + 3σ
    // LCL = Mean - 3σ
    threeSigma := stdDev.Mul(decimal.NewFromInt(3))
    c.UCL = mean.Add(threeSigma)
    c.LCL = mean.Sub(threeSigma)
}

// IsOutOfControl checks if a sample is out of control.
// IsOutOfControl ตรวจสอบว่าตัวอย่างหลุดควบคุมหรือไม่
func (c *SPCChart) IsOutOfControl(s Sample) bool {
    return s.Value.GreaterThan(c.UCL) || s.Value.LessThan(c.LCL)
}

// IsOutOfSpec checks if a sample is out of specification.
// IsOutOfSpec ตรวจสอบว่าตัวอย่างหลุดข้อกำหนดหรือไม่
func (c *SPCChart) IsOutOfSpec(s Sample) bool {
    return s.Value.GreaterThan(c.USL) || s.Value.LessThan(c.LSL)
}

// ProcessCapability calculates Cp and Cpk.
// ProcessCapability คำนวณ Cp และ Cpk
func (c *SPCChart) ProcessCapability() (cp, cpk decimal.Decimal) {
    if len(c.Samples) < 2 {
        return decimal.Zero, decimal.Zero
    }

    // Calculate std dev
    // คำนวณค่าเบี่ยงเบนมาตรฐาน
    varianceSum := decimal.Zero
    for _, s := range c.Samples {
        diff := s.Value.Sub(c.CL)
        varianceSum = varianceSum.Add(diff.Mul(diff))
    }
    variance := varianceSum.Div(decimal.NewFromInt(int64(len(c.Samples) - 1)))
    stdDev := decimal.NewFromFloat(math.Sqrt(variance.InexactFloat64()))

    if stdDev.IsZero() {
        return decimal.Zero, decimal.Zero
    }

    // Cp = (USL - LSL) / (6σ)
    // Cp = (USL - LSL) / (6σ)
    cp = c.USL.Sub(c.LSL).Div(stdDev.Mul(decimal.NewFromInt(6)))

    // Cpk = min((USL - Mean) / 3σ, (Mean - LSL) / 3σ)
    // Cpk = min((USL - Mean) / 3σ, (Mean - LSL) / 3σ)
    threeSigma := stdDev.Mul(decimal.NewFromInt(3))
    cpu := c.USL.Sub(c.CL).Div(threeSigma)
    cpl := c.CL.Sub(c.LSL).Div(threeSigma)
    if cpu.LessThan(cpl) {
        cpk = cpu
    } else {
        cpk = cpl
    }

    return cp, cpk
}
```

## 9.7 Cross-Domain Event

```go
// internal/shared/events/cross_domain.go
package events

import (
    "time"

    "github.com/shopspring/decimal"
)

// HarvestCompletedEvent is published when a farm completes harvest.
// HarvestCompletedEvent ถูก publish เมื่อฟาร์มเก็บเกี่ยวเสร็จ
type HarvestCompletedEvent struct {
    EventID     string          `json:"event_id"`      // รหัส event
    TenantID    string          `json:"tenant_id"`     // รหัสองค์กร
    FarmID      string          `json:"farm_id"`       // รหัสฟาร์ม
    BatchID     string          `json:"batch_id"`      // รหัส批次
    ProductType string          `json:"product_type"`  // ประเภทสินค้า
    Quantity    decimal.Decimal `json:"quantity"`      // จำนวน
    Unit        string          `json:"unit"`          // หน่วย
    Quality     string          `json:"quality"`       // เกรด (A/B/C)
    HarvestedAt time.Time       `json:"harvested_at"`  // เวลาที่เก็บ
    Timestamp   time.Time       `json:"timestamp"`     // เวลาที่ publish
}

// RawMaterialReceivedEvent is published when factory receives raw material.
// RawMaterialReceivedEvent ถูก publish เมื่อโรงงานรับวัตถุดิบ
type RawMaterialReceivedEvent struct {
    EventID       string          `json:"event_id"`        // รหัส event
    TenantID      string          `json:"tenant_id"`       // รหัสองค์กร
    FactoryID     string          `json:"factory_id"`      // รหัสโรงงาน
    SourceBatchID string          `json:"source_batch_id"` // รหัส批次ต้นทาง
    Quantity      decimal.Decimal `json:"quantity"`        // จำนวน
    Unit          string          `json:"unit"`            // หน่วย
    ReceivedAt    time.Time       `json:"received_at"`     // เวลาที่รับ
    Timestamp     time.Time       `json:"timestamp"`       // เวลาที่ publish
}

// ProductionCompletedEvent is published when factory completes production.
// ProductionCompletedEvent ถูก publish เมื่อโรงงานผลิตเสร็จ
type ProductionCompletedEvent struct {
    EventID       string          `json:"event_id"`        // รหัส event
    TenantID      string          `json:"tenant_id"`       // รหัสองค์กร
    FactoryID     string          `json:"factory_id"`      // รหัสโรงงาน
    ProductionID  string          `json:"production_id"`   // รหัสการผลิต
    SourceBatchID string          `json:"source_batch_id"` // รหัส批次ต้นทาง
    OutputBatchID string          `json:"output_batch_id"` // รหัส批次ปลายทาง
    Quantity      decimal.Decimal `json:"quantity"`        // จำนวน
    Unit          string          `json:"unit"`            // หน่วย
    Quality       string          `json:"quality"`         // เกรด
    CompletedAt   time.Time       `json:"completed_at"`    // เวลาที่เสร็จ
    Timestamp     time.Time       `json:"timestamp"`       // เวลาที่ publish
}
```

## 9.8 Module Wire-up

```go
// internal/modules/factory/mes/production/module.go
package production

import (
    "database/sql"

    "github.com/gin-gonic/gin"
    "github.com/yourorg/smart-erp/internal/modules/factory/mes/production/application/command"
    "github.com/yourorg/smart-erp/internal/modules/factory/mes/production/domain/service"
    "github.com/yourorg/smart-erp/internal/modules/factory/mes/production/infrastructure/persistence/postgres"
    "github.com/yourorg/smart-erp/internal/modules/factory/mes/production/interfaces/http/handler"
    "github.com/yourorg/smart-erp/internal/shared/events"
)

// Module represents the MES production module.
// Module แทนโมดูลการผลิต MES
type Module struct {
    // Domain Services
    // Domain Services
    OEECalculator *service.OEECalculator

    // Use Cases
    // Use Cases
    CreateOrderUC *command.CreateProductionOrderUseCase
    StartOrderUC  *command.StartProductionUseCase
    CompleteUC    *command.CompleteProductionUseCase

    // Handlers
    // Handlers
    ProductionHandler *handler.ProductionHandler
}

// NewModule initializes the MES production module.
// NewModule เริ่มต้นโมดูลการผลิต MES
func NewModule(db *sql.DB, eventBus events.EventBus) *Module {
    // Infrastructure layer
    // Infrastructure layer
    orderRepo := postgres.NewProductionOrderRepository(db)

    // Domain services
    // Domain services
    oeeCalc := service.NewOEECalculator()

    // Application layer
    // Application layer
    createUC := command.NewCreateProductionOrderUseCase(orderRepo, eventBus)
    startUC := command.NewStartProductionUseCase(orderRepo, eventBus)
    completeUC := command.NewCompleteProductionUseCase(orderRepo, eventBus)

    // Interfaces layer
    // Interfaces layer
    prodHandler := handler.NewProductionHandler(createUC, startUC, completeUC)

    return &Module{
        OEECalculator:     oeeCalc,
        CreateOrderUC:     createUC,
        StartOrderUC:      startUC,
        CompleteUC:        completeUC,
        ProductionHandler: prodHandler,
    }
}

// RegisterRoutes registers HTTP routes for the module.
// RegisterRoutes ลงทะเบียน HTTP routes สำหรับโมดูล
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
    production := router.Group("/mes/production")
    {
        production.POST("/orders", m.ProductionHandler.CreateOrder)
        production.POST("/orders/:id/start", m.ProductionHandler.StartOrder)
        production.POST("/orders/:id/complete", m.ProductionHandler.CompleteOrder)
        production.GET("/orders/:id", m.ProductionHandler.GetOrder)
        production.GET("/oee/:machine_id", m.ProductionHandler.GetOEE)
    }
}
```

## 9.9 Main Entry Point (Unified)

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

    // Farm modules
    farmInventory "github.com/yourorg/smart-erp/internal/modules/farm/erp/inventory"
    farmCRM "github.com/yourorg/smart-erp/internal/modules/farm/crm/crm"
    farmIoT "github.com/yourorg/smart-erp/internal/modules/farm/iot/iot"

    // Factory modules
    factoryMES "github.com/yourorg/smart-erp/internal/modules/factory/mes/production"
    factoryQMS "github.com/yourorg/smart-erp/internal/modules/factory/qms/quality"
    factoryCosting "github.com/yourorg/smart-erp/internal/modules/factory/costing/standard"

    // Shared
    "github.com/yourorg/smart-erp/internal/shared/config"
    "github.com/yourorg/smart-erp/internal/shared/events"
    "github.com/yourorg/smart-erp/internal/shared/middleware"
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

    // Initialize FARM modules
    // เริ่มต้นโมดูล FARM
    inventoryModule := farmInventory.NewModule(db, eventBus)
    crmModule := farmCRM.NewModule(db, eventBus)
    iotModule := farmIoT.NewModule(db, eventBus)

    // Initialize FACTORY modules
    // เริ่มต้นโมดูล FACTORY
    mesModule := factoryMES.NewModule(db, eventBus)
    qmsModule := factoryQMS.NewModule(db, eventBus)
    costingModule := factoryCosting.NewModule(db, eventBus)

    // Setup router
    // ตั้งค่า router
    router := gin.Default()
    router.Use(middleware.TenantMiddleware(cfg.JWTSecret))

    // Register FARM routes
    // ลงทะเบียน routes FARM
    farmAPI := router.Group("/api/v1/farm")
    inventoryModule.RegisterRoutes(farmAPI.Group("/inventory"))
    crmModule.RegisterRoutes(farmAPI.Group("/crm"))
    iotModule.RegisterRoutes(farmAPI.Group("/iot"))

    // Register FACTORY routes
    // ลงทะเบียน routes FACTORY
    factoryAPI := router.Group("/api/v1/factory")
    mesModule.RegisterRoutes(factoryAPI)
    qmsModule.RegisterRoutes(factoryAPI.Group("/qms"))
    costingModule.RegisterRoutes(factoryAPI.Group("/costing"))

    // Health check
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "ok",
            "domain": "unified",
            "modules": []string{
                "farm/inventory", "farm/crm", "farm/iot",
                "factory/mes", "factory/qms", "factory/costing",
            },
        })
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

    log.Printf("unified ERP server started on port %s", cfg.Port)

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

## 10.1 Farm Domain Checklist

| # | Module | สถานะ | ความสำคัญ |
|---|---|---|---|
| 1 | farm/erp/inventory | ✅ | 🔴 สูง |
| 2 | farm/erp/production | ⏳ | 🔴 สูง |
| 3 | farm/erp/finance | ⏳ | 🔴 สูง |
| 4 | farm/crm/crm | ⏳ | 🟡 กลาง |
| 5 | farm/pos/pos | ⏳ | 🟡 กลาง |
| 6 | farm/iot/iot | ⏳ | 🟡 กลาง |
| 7 | farm/ecommerce | ⏳ | 🟢 ต่ำ |
| 8 | farm/logistics | ⏳ | 🟢 ต่ำ |
| 9 | farm/forecast | ⏳ | 🟢 ต่ำ |
| 10 | farm/reporting | ⏳ | 🟡 กลาง |

## 10.2 Factory Domain Checklist

| # | Module | สถานะ | ความสำคัญ |
|---|---|---|---|
| 1 | factory/mes/production | ⏳ | 🔴 สูง |
| 2 | factory/mes/bom | ⏳ | 🔴 สูง |
| 3 | factory/mes/routing | ⏳ | 🔴 สูง |
| 4 | factory/mes/workorder | ⏳ | 🔴 สูง |
| 5 | factory/mes/scheduling | ⏳ | 🟡 กลาง |
| 6 | factory/qms/quality | ⏳ | 🔴 สูง |
| 7 | factory/qms/spc | ⏳ | 🟡 กลาง |
| 8 | factory/qms/capa | ⏳ | 🟡 กลาง |
| 9 | factory/qms/andon | ⏳ | 🟡 กลาง |
| 10 | factory/costing/standard | ⏳ | 🔴 สูง |
| 11 | factory/costing/actual | ⏳ | 🔴 สูง |
| 12 | factory/costing/variance | ⏳ | 🟡 กลาง |
| 13 | factory/scada/gateway | ⏳ | 🔴 สูง |
| 14 | factory/scada/opcua | ⏳ | 🟡 กลาง |
| 15 | factory/maintenance/pm | ⏳ | 🟡 กลาง |
| 16 | factory/maintenance/cm | ⏳ | 🟢 ต่ำ |
| 17 | factory/maintenance/spare | ⏳ | 🟢 ต่ำ |
| 18 | factory/warehouse | ⏳ | 🟡 กลาง |

## 10.3 Module Completion Checklist (ใช้ทุก Module)

| # | รายการ | สถานะ |
|---|---|---|
| 1 | Domain Entity ครบ (Aggregate Root + VO) | ☐ |
| 2 | Repository Interface ใน domain | ☐ |
| 3 | Repository Implementation | ☐ |
| 4 | Use Case (Command + Query) | ☐ |
| 5 | DTO (Request + Response) | ☐ |
| 6 | HTTP Handler | ☐ |
| 7 | Routes Registration | ☐ |
| 8 | Tenant Middleware | ☐ |
| 9 | Domain Isolation (farm/factory) | ☐ |
| 10 | Auth Middleware | ☐ |
| 11 | Validation | ☐ |
| 12 | Error Handling | ☐ |
| 13 | Domain Events | ☐ |
| 14 | Kafka Producer/Consumer | ☐ |
| 15 | Unit Tests (≥80%) | ☐ |
| 16 | Integration Tests | ☐ |
| 17 | API Documentation | ☐ |
| 18 | Migration Scripts | ☐ |
| 19 | Audit Trail | ☐ |
| 20 | Logging | ☐ |
| 21 | Metrics (Prometheus) | ☐ |
| 22 | SCADA Integration (factory) | ☐ |
| 23 | OEE Tracking (factory) | ☐ |
| 24 | SPC Chart (factory) | ☐ |

---

# 11. Security Code

## 11.1 Security Checklist (Extended)

| # | รายการ | Farm | Factory |
|---|---|---|---|
| 1 | **Authentication** | JWT | JWT + MFA |
| 2 | **Authorization** | RBAC | RBAC + ABAC |
| 3 | **Tenant Isolation** | ✅ | ✅ |
| 4 | **Domain Isolation** | ✅ | ✅ |
| 5 | **Input Validation** | ✅ | ✅ |
| 6 | **SQL Injection** | Parameterized | Parameterized |
| 7 | **XSS** | Sanitize | Sanitize |
| 8 | **CSRF** | Token | Token |
| 9 | **Rate Limiting** | ✅ | ✅ + Per-machine |
| 10 | **HTTPS** | บังคับ | บังคับ |
| 11 | **Password Hashing** | bcrypt 12 | bcrypt 12 |
| 12 | **Secret Management** | Vault | Vault |
| 13 | **Audit Log** | ✅ | ✅ + GMP |
| 14 | **Data Encryption** | AES-256 | AES-256 |
| 15 | **PII Protection** | ✅ | ✅ |
| 16 | **Dependency Scan** | Trivy | Trivy |
| 17 | **SCADA Security** | - | OPC-UA Security |
| 18 | **PLC Access Control** | - | Role-based |
| 19 | **Network Segmentation** | - | VLAN + Firewall |
| 20 | **OT/IT Security** | - | DMZ |

## 11.2 RBAC + ABAC

```go
// internal/shared/middleware/rbac.go
package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// Permission represents a permission.
// Permission แทนสิทธิ์
type Permission string

const (
    // Farm permissions
    // สิทธิ์ฟาร์ม
    PermFarmInventoryRead   Permission = "farm:inventory:read"
    PermFarmInventoryWrite  Permission = "farm:inventory:write"
    PermFarmProductionRead  Permission = "farm:production:read"
    PermFarmProductionWrite Permission = "farm:production:write"

    // Factory permissions
    // สิทธิ์โรงงาน
    PermFactoryMESRead      Permission = "factory:mes:read"
    PermFactoryMESWrite     Permission = "factory:mes:write"
    PermFactoryQMSRead      Permission = "factory:qms:read"
    PermFactoryQMSWrite     Permission = "factory:qms:write"
    PermFactorySCADARead    Permission = "factory:scada:read"
    PermFactorySCADAWrite   Permission = "factory:scada:write"
    PermFactoryCostingRead  Permission = "factory:costing:read"
    PermFactoryCostingWrite Permission = "factory:costing:write"
)

// RolePermissions maps roles to permissions.
// RolePermissions map roles กับ permissions
var RolePermissions = map[string][]Permission{
    "farm_owner": {
        PermFarmInventoryRead, PermFarmInventoryWrite,
        PermFarmProductionRead, PermFarmProductionWrite,
    },
    "farm_staff": {
        PermFarmInventoryRead, PermFarmProductionRead,
    },
    "factory_owner": {
        PermFactoryMESRead, PermFactoryMESWrite,
        PermFactoryQMSRead, PermFactoryQMSWrite,
        PermFactoryCostingRead, PermFactoryCostingWrite,
    },
    "factory_operator": {
        PermFactoryMESRead, PermFactoryMESWrite,
        PermFactoryQMSRead,
    },
    "factory_engineer": {
        PermFactoryMESRead,
        PermFactorySCADARead, PermFactorySCADAWrite,
    },
    "admin": {
        // ทุก permission
    },
}

// RequirePermission checks if user has required permission.
// RequirePermission ตรวจสอบว่าผู้ใช้มี permission ที่ต้องการหรือไม่
func RequirePermission(perm Permission) gin.HandlerFunc {
    return func(c *gin.Context) {
        role := GetRole(c)
        perms := RolePermissions[role]
        for _, p := range perms {
            if p == perm {
                c.Next()
                return
            }
        }
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
            "error": "insufficient permissions",
            "required": perm,
        })
    }
}
```

## 11.3 SCADA Security

```go
// internal/modules/factory/scada/gateway/security.go
package gateway

import (
    "crypto/tls"
    "crypto/x509"
    "errors"
    "os"
)

// SCADASecurityConfig holds security configuration for SCADA.
// SCADASecurityConfig เก็บ configuration ความปลอดภัยสำหรับ SCADA
type SCADASecurityConfig struct {
    // TLS Configuration
    // Configuration TLS
    CertFile   string
    KeyFile    string
    CAFile     string
    ServerName string

    // Authentication
    // การยืนยันตัวตน
    Username string
    Password string

    // Access Control
    // การควบคุมการเข้าถึง
    AllowedTags []string
}

// NewTLSConfig creates a TLS configuration for OPC-UA.
// NewTLSConfig สร้าง TLS configuration สำหรับ OPC-UA
func NewTLSConfig(cfg SCADASecurityConfig) (*tls.Config, error) {
    // Load client certificate
    // โหลด certificate ของ client
    cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
    if err != nil {
        return nil, errors.New("failed to load client certificate")
    }

    // Load CA certificate
    // โหลด CA certificate
    caCert, err := os.ReadFile(cfg.CAFile)
    if err != nil {
        return nil, errors.New("failed to load CA certificate")
    }

    caCertPool := x509.NewCertPool()
    if !caCertPool.AppendCertsFromPEM(caCert) {
        return nil, errors.New("failed to parse CA certificate")
    }

    return &tls.Config{
        Certificates: []tls.Certificate{cert},
        RootCAs:      caCertPool,
        ServerName:   cfg.ServerName,
        MinVersion:   tls.VersionTLS12,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
        },
    }, nil
}

// IsTagAllowed checks if a tag is allowed for access.
// IsTagAllowed ตรวจสอบว่า tag ได้รับอนุญาตให้เข้าถึงหรือไม่
func IsTagAllowed(tag string, allowed []string) bool {
    for _, t := range allowed {
        if t == tag {
            return true
        }
    }
    return false
}
```

---

# 12. Load Test

## 12.1 Load Test Plan (Extended)

| รายการ | Farm | Factory |
|---|---|---|
| **เครื่องมือ** | k6 | k6 + JMeter |
| **เป้าหมาย** | 1,000 concurrent | 500 concurrent |
| **RPS** | 5,000 | 2,000 |
| **Response Time** | P95 < 200ms | P95 < 300ms |
| **Error Rate** | < 0.1% | < 0.1% |
| **Duration** | 30 นาที | 30 นาที |
| **SCADA Polling** | - | 1,000 tags/s |
| **OEE Calculation** | - | 10,000/min |
| **MQTT Messages** | 10,000/s | 5,000/s |

## 12.2 k6 Script — Farm API

```javascript
// tests/load/farm_api_test.js
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');

export const options = {
  stages: [
    { duration: '2m', target: 100 },
    { duration: '5m', target: 500 },
    { duration: '2m', target: 1000 },
    { duration: '5m', target: 1000 },
    { duration: '2m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<200'],
    errors: ['rate<0.001'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const TOKEN = __ENV.JWT_TOKEN;

export default function () {
  // Farm: Create product
  const createPayload = JSON.stringify({
    sku: `MUSH-${__VU}-${__ITER}`,
    name: `เห็ดนางฟ้า ${__VU}`,
    unit: 'kg',
    cost_price: 35,
    sale_price: 80,
  });

  const createRes = http.post(`${BASE_URL}/api/v1/farm/inventory/products`, createPayload, {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${TOKEN}`,
    },
  });

  check(createRes, {
    'create status 201': (r) => r.status === 201,
    'create response time < 200ms': (r) => r.timings.duration < 200,
  }) || errorRate.add(1);

  sleep(1);
}
```

## 12.3 k6 Script — Factory MES API

```javascript
// tests/load/factory_mes_test.js
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

const errorRate = new Rate('errors');
const oeeLatency = new Trend('oee_latency');

export const options = {
  stages: [
    { duration: '2m', target: 50 },
    { duration: '5m', target: 200 },
    { duration: '2m', target: 500 },
    { duration: '5m', target: 500 },
    { duration: '2m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<300'],
    errors: ['rate<0.001'],
    oee_latency: ['p(95)<100'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const TOKEN = __ENV.JWT_TOKEN;

export default function () {
  // Factory: Create production order
  const orderPayload = JSON.stringify({
    order_number: `WO-${__VU}-${__ITER}`,
    product_id: 'PROD-001',
    bom_id: 'BOM-001',
    quantity: 1000,
    unit: 'pcs',
    planned_start: '2024-01-15T08:00:00Z',
    planned_end: '2024-01-15T17:00:00Z',
  });

  const orderRes = http.post(`${BASE_URL}/api/v1/factory/mes/production/orders`, orderPayload, {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${TOKEN}`,
    },
  });

  check(orderRes, {
    'order status 201': (r) => r.status === 201,
  }) || errorRate.add(1);

  // Factory: Get OEE
  const oeeRes = http.get(`${BASE_URL}/api/v1/factory/mes/production/oee/MACHINE-001`, {
    headers: { 'Authorization': `Bearer ${TOKEN}` },
  });

  check(oeeRes, {
    'oee status 200': (r) => r.status === 200,
  }) || errorRate.add(1);

  oeeLatency.add(oeeRes.timings.duration);

  sleep(1);
}
```

## 12.4 SCADA Load Test

```javascript
// tests/load/scada_test.js
import { check } from 'k6';
import mqtt from 'k6/x/mqtt';

export const options = {
  vus: 100,
  duration: '5m',
};

const BROKER = __ENV.MQTT_BROKER || 'tcp://localhost:1883';
const TOPIC = 'factory/telemetry';

export default function () {
  const client = mqtt.connect(BROKER);

  // Publish telemetry
  const payload = JSON.stringify({
    machine_id: `MACHINE-${__VU}`,
    timestamp: Date.now(),
    state: 'RUNNING',
    run_time: 480,
    planned_time: 480,
    ideal_cycle_time: 0.5,
    total_count: 1500,
    good_count: 1450,
    reject_count: 50,
  });

  client.publish(TOPIC, payload);
  check(client, {
    'published': (c) => c.connected,
  });

  client.end();
}
```

## 12.5 รัน Load Test

```bash
# Farm API
BASE_URL=https://api.farm.com JWT_TOKEN=$TOKEN \
  k6 run tests/load/farm_api_test.js

# Factory MES API
BASE_URL=https://api.factory.com JWT_TOKEN=$TOKEN \
  k6 run tests/load/factory_mes_test.js

# SCADA MQTT
MQTT_BROKER=tcp://mqtt.factory.com:1883 \
  k6 run tests/load/scada_test.js

# รันทุก test
make load-test-all
```

## 12.6 Performance Budget

| Metric | Farm Target | Factory Target |
|---|---|---|
| API P95 | < 200ms | < 300ms |
| API P99 | < 500ms | < 800ms |
| DB Query P95 | < 50ms | < 100ms |
| OEE Calc | - | < 100ms |
| MQTT Latency | < 100ms | < 50ms |
| SCADA Poll | - | < 1s |
| WebSocket Push | < 200ms | < 100ms |
| Concurrent Users | 1,000 | 500 |
| RPS | 5,000 | 2,000 |

---

# 13. สรุป

## 13.1 ประโยชน์ที่ได้รับ

| # | ประโยชน์ | Farm | Factory |
|---|---|---|---|
| 1 | **รวมศูนย์** | ✅ | ✅ |
| 2 | **ลดต้นทุน** | 80% | 60% |
| 3 | **เพิ่มผลผลิต** | +30% | +45% (OEE) |
| 4 | **มาตรฐาน** | GAP/GMP | ISO 9001 |
| 5 | **ขยายได้** | ✅ | ✅ |
| 6 | **Multi-tenant** | ✅ | ✅ |
| 7 | **Real-time** | IoT | SCADA/MES |
| 8 | **ตัดสินใจดี** | AI Forecast | OEE/SPC |

## 13.2 ข้อควรระวัง

| # | ข้อควรระวัง | แนวทาง |
|---|---|---|
| 1 | ข้อมูลรั่วระหว่าง tenant | บังคับ tenant_id |
| 2 | ปนกัน farm/factory | Domain Isolation |
| 3 | Performance ตก | Load Test + Index |
| 4 | SCADA security | OPC-UA Security |
| 5 | PLC access | Role-based |
| 6 | Network segmentation | VLAN + DMZ |
| 7 | Migration ยาก | Version + Rollback |
| 8 | Downtime | Zero-downtime deploy |

## 13.3 ข้อดี

- ✅ Clean Architecture → ทดสอบง่าย
- ✅ DDD → Business logic ชัด
- ✅ Modular → แยกทีม
- ✅ Event-Driven → Loose coupling
- ✅ Multi-tenant → ต้นทุนต่ำ
- ✅ Dual-Domain → ขยายตลาด
- ✅ Cross-Domain → Traceability

## 13.4 ข้อเสีย

- ❌ ซับซ้อนกว่า Monolith ปกติ
- ❌ ต้องทีมมีประสบการณ์
- ❌ ต้นทุนเริ่มต้นสูง
- ❌ Setup นาน
- ❌ ต้องเข้าใจทั้ง 2 domain

## 13.5 ข้อห้าม

| # | ข้อห้าม | เหตุผล |
|---|---|---|
| 1 | ❌ import infrastructure ใน domain | ผิด Dependency Rule |
| 2 | ❌ query โดยไม่มี tenant_id | ข้อมูลรั่ว |
| 3 | ❌ ข้าม Domain (farm → factory ตรง) | ต้องผ่าน Event |
| 4 | ❌ commit secret | Security |
| 5 | ❌ SELECT * | Performance |
| 6 | ❌ ข้าม Layer | Architecture |
| 7 | ❌ deploy โดยไม่ test | Quality |
| 8 | ❌ hardcode config | ยืดหยุ่น |
| 9 | ❌ ใช้ float กับเงิน | ใช้ decimal |
| 10 | ❌ SCADA ไม่มี security | OT Security |

## 13.6 ตัวอย่างโค้ดที่รันได้จริง

### Docker Compose (Unified)

```yaml
# docker-compose.yml
version: '3.9'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: smart_erp
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

  clickhouse:
    image: clickhouse/clickhouse-server:23.8
    ports:
      - "8123:8123"
      - "9000:9000"

  mosquitto:
    image: eclipse-mosquitto:2
    ports:
      - "1883:1883"
      - "9001:9001"
    volumes:
      - ./mosquitto.conf:/mosquitto/config/mosquitto.conf

  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://erp_user:erp_pass@postgres:5432/smart_erp?sslmode=disable
      REDIS_URL: redis://redis:6379
      KAFKA_BROKERS: kafka:9092
      CLICKHOUSE_URL: clickhouse://clickhouse:9000
      MQTT_BROKER: tcp://mosquitto:1883
      JWT_SECRET: your-secret-key
    depends_on:
      - postgres
      - redis
      - kafka
      - clickhouse
      - mosquitto

volumes:
  postgres_data:
```

### รันระบบ

```bash
# Clone
git clone https://github.com/yourorg/smart-erp.git
cd smart-erp

# Start infrastructure
docker-compose up -d postgres redis kafka clickhouse mosquitto

# Run migrations
make migrate-up

# Start API
go run cmd/api/main.go

# Test Farm
curl -X POST http://localhost:8080/api/v1/farm/inventory/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"sku":"MUSH-001","name":"เห็ดนางฟ้า","unit":"kg","cost_price":35,"sale_price":80}'

# Test Factory
curl -X POST http://localhost:8080/api/v1/factory/mes/production/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"order_number":"WO-001","product_id":"PROD-001","bom_id":"BOM-001","quantity":1000,"unit":"pcs","planned_start":"2024-01-15T08:00:00Z","planned_end":"2024-01-15T17:00:00Z"}'
```

---

# 14. คู่มือการทดสอบ

## 14.1 Test Strategy

```
┌─────────────────────────────────────────────────────────────────┐
│                    TEST PYRAMID                                  │
│                                                                  │
│                      ┌─────────┐                                 │
│                      │   E2E   │  ← 10%                          │
│                      │  Tests  │                                 │
│                      └─────────┘                                 │
│                   ┌───────────────┐                              │
│                   │ Integration   │  ← 20%                       │
│                   │ Tests         │                              │
│                   └───────────────┘                              │
│              ┌─────────────────────────┐                         │
│              │     Unit Tests          │  ← 70%                  │
│              │  (Domain, Use Case)     │                         │
│              └─────────────────────────┘                         │
└─────────────────────────────────────────────────────────────────┘
```

## 14.2 Unit Test — Farm Product

```go
// internal/modules/farm/erp/inventory/domain/entity/product_test.go
package entity_test

import (
    "testing"

    "github.com/shopspring/decimal"
    "github.com/stretchr/testify/assert"
    "github.com/yourorg/smart-erp/internal/modules/farm/erp/inventory/domain/entity"
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
            }
        })
    }
}
```

## 14.3 Unit Test — Factory OEE

```go
// internal/modules/factory/mes/production/domain/value_object/oee_test.go
package value_object_test

import (
    "testing"

    "github.com/shopspring/decimal"
    "github.com/stretchr/testify/assert"
    "github.com/yourorg/smart-erp/internal/modules/factory/mes/production/domain/value_object"
)

func TestCalculateOEE(t *testing.T) {
    tests := []struct {
        name           string
        runTime        decimal.Decimal
        plannedTime    decimal.Decimal
        idealCycleTime decimal.Decimal
        totalCount     decimal.Decimal
        goodCount      decimal.Decimal
        expectedOEE    string
    }{
        {
            name:           "world class OEE",
            runTime:        decimal.NewFromInt(450),
            plannedTime:    decimal.NewFromInt(480),
            idealCycleTime: decimal.NewFromFloat(0.5),
            totalCount:     decimal.NewFromInt(900),
            goodCount:      decimal.NewFromInt(890),
            // A = 450/480 = 0.9375
            // P = (0.5 × 900)/450 = 1.0
            // Q = 890/900 = 0.9889
            // OEE = 0.9375 × 1.0 × 0.9889 = 0.9271
            expectedOEE: "0.927",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            availability := value_object.CalculateAvailability(tt.runTime, tt.plannedTime)
            performance := value_object.CalculatePerformance(tt.idealCycleTime, tt.totalCount, tt.runTime)
            quality := value_object.CalculateQuality(tt.goodCount, tt.totalCount)

            oee := value_object.NewOEE(availability, performance, quality)

            assert.InDelta(t, 0.927, oee.Value.InexactFloat64(), 0.001)
            assert.True(t, oee.WorldClass())
        })
    }
}
```

## 14.4 Integration Test — Cross-Domain

```go
// tests/integration/cross_domain_test.go
package integration

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/yourorg/smart-erp/internal/modules/farm/erp/inventory/domain/entity"
    "github.com/yourorg/smart-erp/internal/modules/factory/mes/production/domain/entity"
    "github.com/yourorg/smart-erp/internal/shared/events"
)

func TestFarmToFactoryFlow(t *testing.T) {
    ctx := context.Background()

    // 1. Farm harvests
    // 1. ฟาร์มเก็บเกี่ยว
    harvestEvent := events.HarvestCompletedEvent{
        EventID:     "evt-001",
        TenantID:    "tenant-1",
        FarmID:      "farm-001",
        BatchID:     "BATCH-2024-001",
        ProductType: "mushroom_fresh",
        Quantity:    decimal.NewFromInt(500),
        Unit:        "kg",
        Quality:     "A",
        HarvestedAt: time.Now(),
        Timestamp:   time.Now(),
    }

    // Publish to Kafka
    // Publish ไป Kafka
    err := eventBus.Publish(ctx, "farm.harvest.completed", harvestEvent)
    require.NoError(t, err)

    // 2. Factory receives
    // 2. โรงงานรับ
    // Wait for event
    // รอ event
    time.Sleep(2 * time.Second)

    // Verify factory received
    // ตรวจสอบว่าโรงงานรับ
    received := factoryService.GetReceivedBatch(ctx, "BATCH-2024-001")
    assert.NotNil(t, received)
    assert.Equal(t, 500.0, received.Quantity.InexactFloat64())

    // 3. Factory produces
    // 3. โรงงานผลิต
    productionOrder, err := entity.NewProductionOrder(
        "tenant-1", "WO-2024-001", "PROD-001", "BOM-001",
        decimal.NewFromInt(1000), "pcs",
        time.Now(), time.Now().Add(8*time.Hour),
    )
    require.NoError(t, err)

    // 4. Verify traceability
    // 4. ตรวจสอบ traceability
    trace := tracingService.GetTrace(ctx, "BATCH-2024-001")
    assert.NotNil(t, trace)
    assert.Equal(t, "farm-001", trace.FarmID)
    assert.Equal(t, "factory-001", trace.FactoryID)
}
```

## 14.5 Checklist Test

| # | ประเภท | รายการ | เครื่องมือ |
|---|---|---|---|
| 1 | Unit | Domain logic (Farm) | Go testing |
| 2 | Unit | Domain logic (Factory) | Go testing |
| 3 | Unit | OEE Calculation | Go testing |
| 4 | Unit | SPC Chart | Go testing |
| 5 | Unit | Use Case | Go testing + mock |
| 6 | Integration | Repository + DB | Testcontainers |
| 7 | Integration | Kafka | Testcontainers |
| 8 | Integration | SCADA/OPC-UA | Mock server |
| 9 | Integration | Cross-Domain | Testcontainers |
| 10 | E2E | Farm API flow | Go + httptest |
| 11 | E2E | Factory API flow | Go + httptest |
| 12 | Load | Farm API | k6 |
| 13 | Load | Factory API | k6 |
| 14 | Load | SCADA MQTT | k6 |
| 15 | Security | OWASP Top 10 | OWASP ZAP |
| 16 | Security | SCADA security | Custom |
| 17 | Tenant | Data isolation | Custom |
| 18 | Domain | Farm/Factory isolation | Custom |
| 19 | Contract | API contract | Pact |
| 20 | Chaos | Failure injection | Chaos Mesh |

---

# 15. คู่มือการใช้งาน

## 15.1 การเริ่มต้นใช้งาน

### 15.1.1 สำหรับฟาร์มเห็ด

1. **สมัครสมาชิก** ที่ https://app.smart-erp.com/register?domain=farm
2. **ยืนยันอีเมล**
3. **สร้างองค์กร** (Tenant)
4. **เพิ่มโรงเพาะ** ใน Settings → Farms
5. **เพิ่มสินค้า** ใน Inventory → Products
6. **เริ่มบันทึก** การผลิต/ขาย

### 15.1.2 สำหรับโรงงาน SME

1. **สมัครสมาชิก** ที่ https://app.smart-erp.com/register?domain=factory
2. **ยืนยันอีเมล**
3. **สร้างองค์กร** (Tenant)
4. **ตั้งค่าสายการผลิต** ใน Settings → Production Lines
5. **เพิ่ม BOM** ใน MES → BOM
6. **เพิ่ม Work Center** ใน MES → Work Centers
7. **เริ่มผลิต** ผ่าน Work Order

### 15.1.3 สำหรับ Admin (Dual-Domain)

```bash
# 1. Login (Farm)
curl -X POST https://api.smart-erp.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"farm@example.com","password":"password","domain":"farm"}'

# 2. Login (Factory)
curl -X POST https://api.smart-erp.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"factory@example.com","password":"password","domain":"factory"}'

# 3. ใช้ token
export TOKEN="eyJ..."

# 4. Farm: Create product
curl -X POST https://api.smart-erp.com/api/v1/farm/inventory/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"sku":"MUSH-001","name":"เห็ดนางฟ้า","unit":"kg","cost_price":35,"sale_price":80}'

# 5. Factory: Create production order
curl -X POST https://api.smart-erp.com/api/v1/factory/mes/production/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"order_number":"WO-001","product_id":"PROD-001","bom_id":"BOM-001","quantity":1000,"unit":"pcs","planned_start":"2024-01-15T08:00:00Z","planned_end":"2024-01-15T17:00:00Z"}'
```

## 15.2 เมนูหลัก

### Farm Domain

| เมนู | หน้าที่ |
|---|---|
| **Farm Dashboard** | ภาพรวมฟาร์ม |
| **Inventory** | สต็อกก้อนเชื้อ/เห็ด |
| **Production** | ผลิตก้อนเชื้อ |
| **Sales** | ขาย/ออเดอร์ |
| **CRM** | ลูกค้า |
| **IoT** | เซ็นเซอร์ |
| **Forecast** | AI พยากรณ์ |
| **Reports** | รายงาน |

### Factory Domain

| เมนู | หน้าที่ |
|---|---|
| **Factory Dashboard** | ภาพรวมโรงงาน |
| **MES** | ใบสั่งผลิต |
| **BOM** | สูตรการผลิต |
| **Work Orders** | ใบสั่งงาน |
| **QMS** | ควบคุมคุณภาพ |
| **SPC** | แผนภูมิควบคุม |
| **Costing** | ต้นทุน |
| **SCADA** | ควบคุมเครื่องจักร |
| **Maintenance** | ซ่อมบำรุง |
| **Andon** | แจ้งเตือน |
| **Reports** | รายงาน |

## 15.3 การจัดการผู้ใช้

### Farm Roles

| Role | สิทธิ์ |
|---|---|
| **Farm Owner** | ทุกอย่าง |
| **Farm Manager** | ดู/อนุมัติ |
| **Farm Staff** | บันทึกข้อมูล |
| **Farm Viewer** | ดูอย่างเดียว |

### Factory Roles

| Role | สิทธิ์ |
|---|---|
| **Factory Owner** | ทุกอย่าง |
| **Production Manager** | MES, QMS |
| **Operator** | MES, Andon |
| **Engineer** | SCADA, Maintenance |
| **Quality** | QMS, SPC |
| **Viewer** | ดูอย่างเดียว |

---

# 16. คู่มือการบำรุงรักษา

## 16.1 Daily Maintenance

| เวลา | งาน | เครื่องมือ |
|---|---|---|
| 06:00 | Health check (Farm + Factory) | `/health` |
| 07:00 | ดู error log | Grafana/Loki |
| 08:00 | ตรวจสอบ KPI | Grafana |
| 09:00 | ตรวจสอบ SCADA connection | SCADA Gateway |
| 12:00 | ตรวจสอบ disk usage | `df -h` |
| 15:00 | ตรวจสอบ OEE | Dashboard |
| 18:00 | Backup database | pg_dump |
| 23:00 | ตรวจสอบ slow query | pg_stat_statements |

## 16.2 Weekly Maintenance

- ✅ Update dependencies (Dependabot)
- ✅ Review security alerts
- ✅ Check certificate expiry
- ✅ Review capacity planning
- ✅ Test backup restore
- ✅ Review SCADA logs
- ✅ Calibrate sensors (Farm)
- ✅ Verify PLC firmware (Factory)

## 16.3 Monthly Maintenance

- ✅ Patch OS
- ✅ Rotate secrets
- ✅ Review access logs
- ✅ Update documentation
- ✅ Load test
- ✅ SCADA security audit
- ✅ OEE calibration
- ✅ SPC chart review

## 16.4 Backup Strategy

```bash
# Daily backup (Farm + Factory)
pg_dump -h $DB_HOST -U $DB_USER -d smart_erp -F c -f backup_$(date +%Y%m%d).dump

# Upload to S3
aws s3 cp backup_$(date +%Y%m%d).dump s3://smart-erp-backup/

# ClickHouse backup (time-series)
clickhouse-backup create

# Retention: 7 daily, 4 weekly, 12 monthly
```

## 16.5 Monitoring

| Metric | Farm Alert | Factory Alert |
|---|---|---|
| CPU | > 80% (5 min) | > 80% (5 min) |
| Memory | > 85% | > 85% |
| Disk | > 90% | > 90% |
| Error Rate | > 1% | > 1% |
| P95 Latency | > 500ms | > 800ms |
| DB Connections | > 80% pool | > 80% pool |
| SCADA Poll | - | > 5s |
| MQTT Lag | > 1s | > 500ms |
| OEE Calc | - | > 200ms |
| Sensor Offline | > 5 min | - |

---

# 17. คู่มือการขยาย/แก้ไข

## 17.1 การเพิ่ม Farm Module ใหม่

```bash
# 1. สร้างโครงสร้าง
mkdir -p internal/modules/farm/{prefix}/{module}/{domain,application,infrastructure,interfaces}

# 2. สร้างไฟล์พื้นฐาน
touch internal/modules/farm/{prefix}/{module}/module.go
touch internal/modules/farm/{prefix}/{module}/domain/entity/{aggregate}.go
touch internal/modules/farm/{prefix}/{module}/domain/repository/{aggregate}_repository.go
touch internal/modules/farm/{prefix}/{module}/application/command/create_{aggregate}.go
touch internal/modules/farm/{prefix}/{module}/infrastructure/persistence/postgres/{aggregate}_repo_impl.go
touch internal/modules/farm/{prefix}/{module}/interfaces/http/handler/{aggregate}_handler.go

# 3. Register ใน main.go
```

## 17.2 การเพิ่ม Factory Module ใหม่

```bash
# 1. สร้างโครงสร้าง
mkdir -p internal/modules/factory/{prefix}/{module}/{domain,application,infrastructure,interfaces}

# 2. สำหรับ SCADA module
mkdir -p internal/modules/factory/scada/{module}/{domain,application,infrastructure,interfaces}
touch internal/modules/factory/scada/{module}/infrastructure/opcua_client.go
touch internal/modules/factory/scada/{module}/infrastructure/modbus_client.go

# 3. Register ใน main.go
```

## 17.3 การเพิ่ม Field ใหม่

```go
// 1. เพิ่มใน Entity
type Product struct {
    // ... existing
    NewField string `json:"new_field"`
}

// 2. เพิ่มใน Migration
// migrations/farm/20240115120000_add_new_field.up.sql
ALTER TABLE products ADD COLUMN new_field VARCHAR(255);

// 3. เพิ่มใน DTO
type CreateProductRequest struct {
    NewField string `json:"new_field"`
}

// 4. อัปเดต Repository
```

## 17.4 การเพิ่ม SCADA Protocol

```go
// internal/modules/factory/scada/gateway/modbus_client.go
package gateway

import (
    "context"
    "fmt"

    "github.com/goburrow/modbus"
)

// ModbusClient wraps Modbus TCP client.
// ModbusClient ครอบ Modbus TCP client
type ModbusClient struct {
    handler *modbus.TCPClientHandler
    client  modbus.Client
}

// NewModbusClient creates a new Modbus client.
// NewModbusClient สร้าง Modbus client ใหม่
func NewModbusClient(address string, port int) (*ModbusClient, error) {
    handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%d", address, port))
    if err := handler.Connect(); err != nil {
        return nil, err
    }

    return &ModbusClient{
        handler: handler,
        client:  modbus.NewClient(handler),
    }, nil
}

// ReadHoldingRegisters reads holding registers.
// ReadHoldingRegisters อ่าน holding registers
func (c *ModbusClient) ReadHoldingRegisters(ctx context.Context, address, quantity uint16) ([]byte, error) {
    return c.client.ReadHoldingRegisters(address, quantity)
}
```

## 17.5 การ Scale Out

```
Phase 1: Vertical Scaling
  └── เพิ่ม CPU/RAM

Phase 2: Read Replica
  ├── Master (Write)
  └── Replica (Read) × N

Phase 3: Sharding by tenant_id
  ├── Shard 1: Farm tenants A-M
  ├── Shard 2: Farm tenants N-Z
  ├── Shard 3: Factory tenants A-M
  └── Shard 4: Factory tenants N-Z

Phase 4: Microservices
  ├── Farm Services
  │   ├── Inventory Service
  │   ├── Production Service
  │   └── IoT Service
  └── Factory Services
      ├── MES Service
      ├── QMS Service
      ├── SCADA Service
      └── Costing Service

Phase 5: Edge Computing (Factory)
  ├── Edge Gateway (on-premise)
  ├── Local SCADA
  └── Cloud Sync
```

---

# 18. Git Flow

## 18.1 Branching Model (Extended)

```
┌─────────────────────────────────────────────────────────────────┐
│                      GIT FLOW (DUAL-DOMAIN)                      │
└─────────────────────────────────────────────────────────────────┘

main (production)
  │
  ├─── hotfix/urgent-fix ──────────────┐
  │                                     │
  └─── develop (staging)               │
        │                               │
        ├─── feature/farm-inventory ───┤
        │                               │
        ├─── feature/factory-mes ──────┤
        │                               │
        ├─── feature/factory-scada ────┤
        │                               │
        ├─── release/v1.2.0 ───────────┘
        │
        └─── bugfix/fix-oee-calc ──────┘
```

## 18.2 Branch Types

| Branch | จาก | Merge ไป | ชื่อ |
|---|---|---|---|
| **main** | - | - | `main` |
| **develop** | main | main | `develop` |
| **feature** | develop | develop | `feature/{domain}-{module}-{desc}` |
| **release** | develop | main + develop | `release/v{X.Y.Z}` |
| **hotfix** | main | main + develop | `hotfix/{desc}` |
| **bugfix** | develop | develop | `bugfix/{domain}-{desc}` |

## 18.3 Commit Convention

```
<type>(<domain>/<scope>): <subject>

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
feat(factory/mes): add OEE calculation service

- Add OEE value object
- Add OEECalculator domain service
- Add OEE API endpoint
- Add unit tests for OEE

Closes #456
```

```
feat(farm/inventory): add product creation API

- Add CreateProductUseCase
- Add ProductRepository implementation
- Add HTTP handler with validation

Closes #123
```

## 18.4 Workflow

```bash
# 1. Start feature (Farm)
git checkout develop
git pull origin develop
git checkout -b feature/farm-inventory-create-product

# 2. Start feature (Factory)
git checkout develop
git pull origin develop
git checkout -b feature/factory-mes-oee-calc

# 3. Develop
git add .
git commit -m "feat(farm/inventory): add product entity"
git commit -m "feat(factory/mes): add OEE calculator"

# 4. Push
git push origin feature/farm-inventory-create-product

# 5. Create PR
gh pr create --base develop \
  --title "feat(farm/inventory): product creation" \
  --body "..."
```

---

# 19. Code Review & PR

## 19.1 PR Template (Extended)

```markdown
<!-- .github/pull_request_template.md -->
## 📋 Description
<!-- อธิบายการเปลี่ยนแปลง -->

## 🎯 Domain
- [ ] 🌱 Farm
- [ ] 🏭 Factory
- [ ] 🤝 Shared

## 📦 Module
<!-- ระบุ module ที่แก้ไข -->

## 🔄 Type of Change
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
- [ ] Cross-domain tests
- [ ] SCADA tests (factory)
- [ ] Manual testing

## 📸 Screenshots (if applicable)

## ✅ Checklist
### Architecture
- [ ] Code follows Clean Architecture
- [ ] Dependency Rule respected
- [ ] Domain Isolation (farm/factory) maintained
- [ ] No cross-domain direct calls

### Security
- [ ] Tenant isolation verified
- [ ] Input validation added
- [ ] No SQL injection
- [ ] No secrets committed
- [ ] SCADA security (factory)

### Quality
- [ ] Self-review done
- [ ] Comments added (2 languages)
- [ ] Documentation updated
- [ ] No new warnings
- [ ] Tests added (≥80% coverage)
- [ ] All tests pass

### Performance
- [ ] No N+1 queries
- [ ] Indexes considered
- [ ] Load test (if needed)

## 🔍 Reviewer Notes
```

## 19.2 Code Review Checklist (Extended)

| # | หัวข้อ | รายการ |
|---|---|---|
| 1 | **Architecture** | ถูก Layer? Dependency Rule? |
| 2 | **Domain Isolation** | Farm/Factory แยก? |
| 3 | **Cross-Domain** | ผ่าน Event เท่านั้น? |
| 4 | **Tenant** | มี tenant_id ทุก query? |
| 5 | **Security** | Validate input? SQL injection? |
| 6 | **SCADA** | Security? OPC-UA? |
| 7 | **OEE** | สูตรถูก? Clamp 0-1? |
| 8 | **SPC** | Control limits ถูก? |
| 9 | **Error Handling** | Handle ครบ? |
| 10 | **Testing** | Coverage ≥ 80%? |
| 11 | **Performance** | N+1? Index? |
| 12 | **Naming** | ชื่อสื่อความหมาย? |
| 13 | **Comments** | 2 ภาษา? |
| 14 | **Documentation** | อัปเดต? |

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

---

# 20. CI/CD

## 20.1 Pipeline Overview (Extended)

```
┌─────────────────────────────────────────────────────────────────┐
│                    CI/CD PIPELINE (DUAL-DOMAIN)                  │
└─────────────────────────────────────────────────────────────────┘

Push/PR
   │
   ▼
┌──────────┐
│  Lint    │───▶ golangci-lint
└────┬─────┘
     ▼
┌──────────┐
│  Test    │───▶ go test -race -cover
│  Farm    │
└────┬─────┘
     ▼
┌──────────┐
│  Test    │───▶ go test -race -cover
│  Factory │
└────┬─────┘
     ▼
┌──────────┐
│  Test    │───▶ Cross-domain tests
│  Cross   │
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

  test-farm:
    name: Test Farm
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
      - name: Run farm tests
        run: go test -v -race -coverprofile=coverage-farm.out ./internal/modules/farm/...

  test-factory:
    name: Test Factory
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_PASSWORD: test
          POSTGRES_DB: testdb
        ports: ['5432:5432']
      clickhouse:
        image: clickhouse/clickhouse-server:23.8
        ports: ['8123:8123', '9000:9000']
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      - name: Run factory tests
        run: go test -v -race -coverprofile=coverage-factory.out ./internal/modules/factory/...

  test-cross-domain:
    name: Test Cross-Domain
    runs-on: ubuntu-latest
    needs: [test-farm, test-factory]
    services:
      kafka:
        image: confluentinc/cp-kafka:7.5.0
        ports: ['9092:9092']
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      - name: Run cross-domain tests
        run: go test -v -race ./tests/integration/...

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
    needs: [lint, test-farm, test-factory, test-cross-domain, security]
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

## 20.3 Dockerfile (Unified)

```dockerfile
# Multi-stage build
# Build แบบ multi-stage

# Stage 1: Builder
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.version=$(git describe --tags --always)" \
    -o /app/api ./cmd/api

# Stage 2: Runtime
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

RUN adduser -D -u 1000 appuser

WORKDIR /app

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

### 21.1.2 รวบรวมข้อมูล (Collect Data)

| แหล่งข้อมูล | รายละเอียด |
|---|---|
| **Logs** | Application, Access, SCADA |
| **Metrics** | CPU, Memory, DB, OEE |
| **Traces** | Distributed tracing (Jaeger) |
| **Events** | Deployment, Config changes |
| **Users** | Reports จากผู้ใช้ |
| **SCADA** | PLC logs, Machine state |

### 21.1.3 ระบุสาเหตุที่เป็นไปได้ (Identify Possible Causes)

**Fishbone Diagram (Factory Example):**

```
                    ┌── คน ──┐
                    │ Operator ใหม่
                    │ ไม่มี training
                    └────────┘
┌── วิธี ──┐        │        ┌── เครื่อง ──┐
│ ไม่มี SOP │        │        │ PLC error  │
│ ขั้นตอนผิด│────────┼────────│ Sensor เสีย │
└──────────┘        │        └────────────┘
                    │
              ┌─────┴─────┐
              │  ปัญหา:   │
              │ OEE ต่ำ   │
              └─────┬─────┘
                    │
┌── วัสดุ ──┐        │        ┌── สภาพแวดล้อม ──┐
│ วัตถุดิบ  │        │        │ อุณหภูมิสูง     │
│ ไม่คงที่  │────────┼────────│ ฝุ่น            │
└──────────┘        │        └──────────────────┘
                    │
              ┌── การวัด ──┐
              │ Sensor ไม่ calibrate
              │ ข้อมูลผิด
              └────────────┘
```

### 21.1.4 หา "สาเหตุหลัก" (Find the Root Cause)

**ใช้ 5 Whys (Factory Example):**

```
1. ทำไม OEE ต่ำ?
   → เพราะ Availability ต่ำ

2. ทำไม Availability ต่ำ?
   → เพราะเครื่องจักรหยุดบ่อย

3. ทำไมเครื่องจักรหยุดบ่อย?
   → เพราะ Sensor เสีย

4. ทำไม Sensor เสีย?
   → เพราะไม่ได้ calibrate

5. ทำไมไม่ได้ calibrate?
   → เพราะไม่มี PM schedule (Root Cause)
```

### 21.1.5 วางแผนและแก้ไข (Implement Solution)

| ระดับ | การแก้ไข |
|---|---|
| **Immediate** | เปลี่ยน Sensor ทันที |
| **Short-term** | เพิ่ม PM schedule |
| **Long-term** | IoT sensor monitoring |
| **Preventive** | Predictive maintenance |

### 21.1.6 ติดตามผล (Monitor)

- ✅ OEE กลับมา > 85%
- ✅ เพิ่ม alert สำหรับ sensor
- ✅ PM schedule ครบ
- ✅ Predictive maintenance

## 21.2 RCA Template

```markdown
# Root Cause Analysis Report

## 1. Incident Summary
- **Date**: 2024-01-15 14:30
- **Domain**: Factory
- **Duration**: 2 ชั่วโมง
- **Severity**: P1 (Critical)
- **Impact**: OEE ลด 20%

## 2. Timeline
| เวลา | เหตุการณ์ |
|---|---|
| 14:30 | Alert: OEE < 70% |
| 14:35 | ทีมเริ่ม investigate |
| 14:45 | พบ Sensor เสีย |
| 15:00 | เปลี่ยน Sensor |
| 15:30 | Calibrate |
| 16:30 | OEE กลับปกติ |

## 3. Root Cause
ไม่มี PM schedule สำหรับ sensor calibration

## 4. 5 Whys
1. ทำไม OEE ต่ำ? → Availability ต่ำ
2. ทำไม Availability ต่ำ? → เครื่องหยุดบ่อย
3. ทำไมเครื่องหยุดบ่อย? → Sensor เสีย
4. ทำไม Sensor เสีย? → ไม่ได้ calibrate
5. ทำไมไม่ได้ calibrate? → ไม่มี PM schedule

## 5. Action Items
| # | การแก้ไข | ผู้รับผิดชอบ | กำหนดเสร็จ |
|---|---|---|---|
| 1 | เปลี่ยน Sensor | Engineer | 15 ม.ค. |
| 2 | เพิ่ม PM schedule | Maintenance | 16 ม.ค. |
| 3 | IoT monitoring | IT | 20 ม.ค. |
| 4 | Predictive ML | Data | 30 ม.ค. |

## 6. Lessons Learned
- ต้องมี PM schedule
- ต้องมี sensor monitoring
- ต้องมี predictive maintenance

## 7. Prevention
- Automated PM scheduling
- IoT sensor monitoring
- Predictive maintenance ML
```

---

# 22. SME Factory Extension

## 22.1 MES Module

### 22.1.1 MES Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                    MES ARCHITECTURE                              │
└─────────────────────────────────────────────────────────────────┘

┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   ERP        │────▶│   MES        │────▶│   SCADA      │
│  (Planning)  │     │  (Execution) │     │  (Control)   │
└──────────────┘     └──────────────┘     └──────────────┘
       │                    │                    │
       │                    │                    │
       ▼                    ▼                    ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│ Sales Order  │     │ Work Order   │     │ PLC/CNC      │
│ MRP          │     │ Routing      │     │ Sensor       │
│ BOM          │     │ OEE          │     │ Actuator     │
└──────────────┘     └──────────────┘     └──────────────┘
```

### 22.1.2 MES Entities

```go
// internal/modules/factory/mes/production/domain/entity/work_order.go
package entity

import (
    "errors"
    "time"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)

// WorkOrder represents a work order in MES.
// WorkOrder แทนใบสั่งงานใน MES
type WorkOrder struct {
    ID              string          // รหัสใบสั่งงาน
    TenantID        string          // รหัสองค์กร
    ProductionOrder string          // รหัสใบสั่งผลิต
    OperationNumber int             // ลำดับการทำงาน
    WorkCenterID    string          // รหัสศูนย์งาน
    MachineID       string          // รหัสเครื่องจักร
    OperatorID      string          // รหัสผู้ปฏิบัติงาน
    Status          WorkOrderStatus // สถานะ
    PlannedQty      decimal.Decimal // จำนวนที่วางแผน
    CompletedQty    decimal.Decimal // จำนวนที่เสร็จ
    RejectedQty     decimal.Decimal // จำนวนที่เสีย
    ScrapQty        decimal.Decimal // จำนวนเศษ
    PlannedStart    time.Time       // วางแผนเริ่ม
    PlannedEnd      time.Time       // วางแผนสิ้นสุด
    ActualStart     *time.Time      // เริ่มจริง
    ActualEnd       *time.Time      // สิ้นสุดจริง
    StandardTime    decimal.Decimal // เวลามาตรฐาน (นาที)
    ActualTime      decimal.Decimal // เวลาจริง (นาที)
    CreatedAt       time.Time       // วันที่สร้าง
    UpdatedAt       time.Time       // วันที่แก้ไข
    Version         int             // Optimistic Lock
}

// WorkOrderStatus represents the status of a work order.
// WorkOrderStatus แทนสถานะของใบสั่งงาน
type WorkOrderStatus string

const (
    WOStatusPending    WorkOrderStatus = "pending"    // รอดำเนินการ
    WOStatusInProgress WorkOrderStatus = "in_progress" // กำลังดำเนินการ
    WOStatusPaused     WorkOrderStatus = "paused"     // หยุดชั่วคราว
    WOStatusCompleted  WorkOrderStatus = "completed"  // เสร็จ
    WOStatusCancelled  WorkOrderStatus = "cancelled"  // ยกเลิก
)

// NewWorkOrder creates a new WorkOrder.
// NewWorkOrder สร้าง WorkOrder ใหม่
func NewWorkOrder(
    tenantID, productionOrder string, operationNumber int,
    workCenterID, machineID string,
    plannedQty decimal.Decimal,
    plannedStart, plannedEnd time.Time,
) (*WorkOrder, error) {
    if tenantID == "" {
        return nil, errors.New("tenant_id is required")
    }
    if plannedQty.LessThanOrEqual(decimal.Zero) {
        return nil, errors.New("planned_qty must be positive")
    }

    now := time.Now()
    return &WorkOrder{
        ID:              uuid.NewString(),
        TenantID:        tenantID,
        ProductionOrder: productionOrder,
        OperationNumber: operationNumber,
        WorkCenterID:    workCenterID,
        MachineID:       machineID,
        Status:          WOStatusPending,
        PlannedQty:      plannedQty,
        CompletedQty:    decimal.Zero,
        RejectedQty:     decimal.Zero,
        ScrapQty:        decimal.Zero,
        PlannedStart:    plannedStart,
        PlannedEnd:      plannedEnd,
        StandardTime:    decimal.Zero,
        ActualTime:      decimal.Zero,
        CreatedAt:       now,
        UpdatedAt:       now,
        Version:         1,
    }, nil
}

// Start starts the work order.
// Start เริ่มใบสั่งงาน
func (wo *WorkOrder) Start(operatorID string) error {
    if wo.Status != WOStatusPending {
        return errors.New("only pending work orders can be started")
    }

    now := time.Now()
    wo.Status = WOStatusInProgress
    wo.OperatorID = operatorID
    wo.ActualStart = &now
    wo.UpdatedAt = now
    wo.Version++
    return nil
}

// RecordOutput records production output.
// RecordOutput บันทึกผลผลิต
func (wo *WorkOrder) RecordOutput(goodQty, rejectedQty, scrapQty decimal.Decimal) error {
    if wo.Status != WOStatusInProgress {
        return errors.New("only in-progress work orders can record output")
    }

    wo.CompletedQty = wo.CompletedQty.Add(goodQty)
    wo.RejectedQty = wo.RejectedQty.Add(rejectedQty)
    wo.ScrapQty = wo.ScrapQty.Add(scrapQty)
    wo.UpdatedAt = time.Now()
    wo.Version++
    return nil
}

// Complete completes the work order.
// Complete เสร็จสิ้นใบสั่งงาน
func (wo *WorkOrder) Complete() error {
    if wo.Status != WOStatusInProgress {
        return errors.New("only in-progress work orders can be completed")
    }

    now := time.Now()
    wo.Status = WOStatusCompleted
    wo.ActualEnd = &now
    if wo.ActualStart != nil {
        wo.ActualTime = decimal.NewFromFloat(now.Sub(*wo.ActualStart).Minutes())
    }
    wo.UpdatedAt = now
    wo.Version++
    return nil
}

// CalculateYield calculates the yield percentage.
// CalculateYield คำนวณเปอร์เซ็นต์ผลได้
func (wo *WorkOrder) CalculateYield() decimal.Decimal {
    total := wo.CompletedQty.Add(wo.RejectedQty).Add(wo.ScrapQty)
    if total.IsZero() {
        return decimal.Zero
    }
    return wo.CompletedQty.Div(total).Mul(decimal.NewFromInt(100))
}
```

## 22.2 QMS Module

### 22.2.1 Quality Control Entity

```go
// internal/modules/factory/qms/quality/domain/entity/quality_check.go
package entity

import (
    "errors"
    "time"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)

// QualityCheck represents a quality check.
// QualityCheck แทนการตรวจสอบคุณภาพ
type QualityCheck struct {
    ID            string          // รหัสการตรวจ
    TenantID      string          // รหัสองค์กร
    WorkOrderID   string          // รหัสใบสั่งงาน
    CheckType     CheckType       // ประเภทการตรวจ
    SampleSize    int             // ขนาดตัวอย่าง
    DefectCount   int             // จำนวนของเสีย
    DefectRate    decimal.Decimal // อัตราของเสีย
    Result        CheckResult     // ผลการตรวจ
    InspectorID   string          // รหัสผู้ตรวจ
    Notes         string          // หมายเหตุ
    Attachments   []string        // ไฟล์แนบ
    CheckedAt     time.Time       // เวลาที่ตรวจ
    CreatedAt     time.Time       // วันที่สร้าง
    UpdatedAt     time.Time       // วันที่แก้ไข
}

// CheckType represents the type of quality check.
// CheckType แทนประเภทการตรวจสอบ
type CheckType string

const (
    CheckTypeIncoming  CheckType = "incoming"  // ตรวจรับ
    CheckTypeInProcess CheckType = "in_process" // ระหว่างผลิต
    CheckTypeFinal     CheckType = "final"     // ตรวจสุดท้าย
    CheckTypeOutgoing  CheckType = "outgoing"  // ตรวจก่อนส่ง
)

// CheckResult represents the result of a quality check.
// CheckResult แทนผลการตรวจสอบ
type CheckResult string

const (
    CheckResultPass    CheckResult = "pass"    // ผ่าน
    CheckResultFail    CheckResult = "fail"    // ไม่ผ่าน
    CheckResultPartial CheckResult = "partial" // ผ่านบางส่วน
    CheckResultRework  CheckResult = "rework"  // ต้องแก้ไข
)

// NewQualityCheck creates a new QualityCheck.
// NewQualityCheck สร้าง QualityCheck ใหม่
func NewQualityCheck(
    tenantID, workOrderID string,
    checkType CheckType,
    sampleSize, defectCount int,
    inspectorID string,
) (*QualityCheck, error) {
    if tenantID == "" {
        return nil, errors.New("tenant_id is required")
    }
    if sampleSize <= 0 {
        return nil, errors.New("sample_size must be positive")
    }
    if defectCount < 0 || defectCount > sampleSize {
        return nil, errors.New("defect_count must be between 0 and sample_size")
    }

    defectRate := decimal.NewFromInt(int64(defectCount)).
        Div(decimal.NewFromInt(int64(sampleSize))).
        Mul(decimal.NewFromInt(100))

    result := CheckResultPass
    if defectCount > 0 {
        if defectRate.GreaterThan(decimal.NewFromInt(10)) {
            result = CheckResultFail
        } else {
            result = CheckResultPartial
        }
    }

    now := time.Now()
    return &QualityCheck{
        ID:          uuid.NewString(),
        TenantID:    tenantID,
        WorkOrderID: workOrderID,
        CheckType:   checkType,
        SampleSize:  sampleSize,
        DefectCount: defectCount,
        DefectRate:  defectRate,
        Result:      result,
        InspectorID: inspectorID,
        CheckedAt:   now,
        CreatedAt:   now,
        UpdatedAt:   now,
    }, nil
}
```

## 22.3 Costing Module

### 22.3.1 Standard Cost Entity

```go
// internal/modules/factory/costing/standard/domain/entity/standard_cost.go
package entity

import (
    "errors"
    "time"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)

// StandardCost represents the standard cost of a product.
// StandardCost แทนต้นทุนมาตรฐานของสินค้า
type StandardCost struct {
    ID              string          // รหัสต้นทุนมาตรฐาน
    TenantID        string          // รหัสองค์กร
    ProductID       string          // รหัสสินค้า
    Version         string          // เวอร์ชัน
    EffectiveDate   time.Time       // วันที่มีผล
    MaterialCost    decimal.Decimal // ต้นทุนวัสดุ
    LaborCost       decimal.Decimal // ต้นทุนแรงงาน
    OverheadCost    decimal.Decimal // ต้นทุนค่าใช้จ่ายการผลิต
    TotalCost       decimal.Decimal // ต้นทุนรวม
    Currency        string          // สกุลเงิน
    CreatedAt       time.Time       // วันที่สร้าง
    UpdatedAt       time.Time       // วันที่แก้ไข
    CreatedBy       string          // ผู้สร้าง
}

// NewStandardCost creates a new StandardCost.
// NewStandardCost สร้าง StandardCost ใหม่
func NewStandardCost(
    tenantID, productID, version string,
    materialCost, laborCost, overheadCost decimal.Decimal,
    currency string,
) (*StandardCost, error) {
    if tenantID == "" {
        return nil, errors.New("tenant_id is required")
    }
    if productID == "" {
        return nil, errors.New("product_id is required")
    }
    if materialCost.IsNegative() || laborCost.IsNegative() || overheadCost.IsNegative() {
        return nil, errors.New("costs cannot be negative")
    }

    totalCost := materialCost.Add(laborCost).Add(overheadCost)

    now := time.Now()
    return &StandardCost{
        ID:            uuid.NewString(),
        TenantID:      tenantID,
        ProductID:     productID,
        Version:       version,
        EffectiveDate: now,
        MaterialCost:  materialCost,
        LaborCost:     laborCost,
        OverheadCost:  overheadCost,
        TotalCost:     totalCost,
        Currency:      currency,
        CreatedAt:     now,
        UpdatedAt:     now,
    }, nil
}
```

### 22.3.2 Variance Analysis

```go
// internal/modules/factory/costing/variance/domain/service/variance_calculator.go
package service

import (
    "github.com/shopspring/decimal"
)

// Variance represents cost variance.
// Variance แทนผลต่างต้นทุน
type Variance struct {
    StandardCost decimal.Decimal // ต้นทุนมาตรฐาน
    ActualCost   decimal.Decimal // ต้นทุนจริง
    Variance     decimal.Decimal // ผลต่าง
    VariancePct  decimal.Decimal // ผลต่าง %
    IsFavorable  bool            // เป็นประโยชน์หรือไม่
}

// CalculateMaterialVariance calculates material cost variance.
// CalculateMaterialVariance คำนวณผลต่างต้นทุนวัสดุ
func CalculateMaterialVariance(
    standardQty, standardPrice decimal.Decimal,
    actualQty, actualPrice decimal.Decimal,
) Variance {
    standardCost := standardQty.Mul(standardPrice)
    actualCost := actualQty.Mul(actualPrice)
    variance := actualCost.Sub(standardCost)

    varPct := decimal.Zero
    if !standardCost.IsZero() {
        varPct = variance.Div(standardCost).Mul(decimal.NewFromInt(100))
    }

    // Favorable if actual < standard
    // เป็นประโยชน์ถ้าจริง < มาตรฐาน
    isFavorable := variance.LessThan(decimal.Zero)

    return Variance{
        StandardCost: standardCost,
        ActualCost:   actualCost,
        Variance:     variance,
        VariancePct:  varPct,
        IsFavorable:  isFavorable,
    }
}

// CalculateLaborVariance calculates labor cost variance.
// CalculateLaborVariance คำนวณผลต่างต้นทุนแรงงาน
func CalculateLaborVariance(
    standardHours, standardRate decimal.Decimal,
    actualHours, actualRate decimal.Decimal,
) Variance {
    standardCost := standardHours.Mul(standardRate)
    actualCost := actualHours.Mul(actualRate)
    variance := actualCost.Sub(standardCost)

    varPct := decimal.Zero
    if !standardCost.IsZero() {
        varPct = variance.Div(standardCost).Mul(decimal.NewFromInt(100))
    }

    isFavorable := variance.LessThan(decimal.Zero)

    return Variance{
        StandardCost: standardCost,
        ActualCost:   actualCost,
        Variance:     variance,
        VariancePct:  varPct,
        IsFavorable:  isFavorable,
    }
}
```

## 22.4 Maintenance Module

### 22.4.1 Preventive Maintenance

```go
// internal/modules/factory/maintenance/pm/domain/entity/pm_schedule.go
package entity

import (
    "errors"
    "time"

    "github.com/google/uuid"
)

// PMSchedule represents a preventive maintenance schedule.
// PMSchedule แทนตารางบำรุงรักษาเชิงป้องกัน
type PMSchedule struct {
    ID            string        // รหัสตาราง PM
    TenantID      string        // รหัสองค์กร
    MachineID     string        // รหัสเครื่องจักร
    TaskName      string        // ชื่องาน
    Frequency     Frequency     // ความถี่
    IntervalDays  int           // ระยะห่าง (วัน)
    LastDoneAt    *time.Time    // ครั้งสุดท้าย
    NextDueAt     time.Time     // ครั้งถัดไป
    AssignedTo    string        // ผู้รับผิดชอบ
    Checklist     []string      // รายการตรวจ
    Status        PMStatus      // สถานะ
    CreatedAt     time.Time     // วันที่สร้าง
    UpdatedAt     time.Time     // วันที่แก้ไข
}

// Frequency represents the frequency of PM.
// Frequency แทนความถี่ของ PM
type Frequency string

const (
    FrequencyDaily   Frequency = "daily"   // รายวัน
    FrequencyWeekly  Frequency = "weekly"  // รายสัปดาห์
    FrequencyMonthly Frequency = "monthly" // รายเดือน
    FrequencyQuarterly Frequency = "quarterly" // รายไตรมาส
    FrequencyYearly  Frequency = "yearly"  // รายปี
)

// PMStatus represents the status of PM.
// PMStatus แทนสถานะของ PM
type PMStatus string

const (
    PMStatusScheduled PMStatus = "scheduled" // กำหนดการ
    PMStatusDue       PMStatus = "due"       // ครบกำหนด
    PMStatusOverdue   PMStatus = "overdue"   // เกินกำหนด
    PMStatusCompleted PMStatus = "completed" // เสร็จ
)

// NewPMSchedule creates a new PM schedule.
// NewPMSchedule สร้างตาราง PM ใหม่
func NewPMSchedule(
    tenantID, machineID, taskName string,
    frequency Frequency, intervalDays int,
    assignedTo string,
) (*PMSchedule, error) {
    if tenantID == "" {
        return nil, errors.New("tenant_id is required")
    }
    if machineID == "" {
        return nil, errors.New("machine_id is required")
    }
    if intervalDays <= 0 {
        return nil, errors.New("interval_days must be positive")
    }

    now := time.Now()
    nextDue := now.AddDate(0, 0, intervalDays)

    return &PMSchedule{
        ID:           uuid.NewString(),
        TenantID:     tenantID,
        MachineID:    machineID,
        TaskName:     taskName,
        Frequency:    frequency,
        IntervalDays: intervalDays,
        NextDueAt:    nextDue,
        AssignedTo:   assignedTo,
        Status:       PMStatusScheduled,
        CreatedAt:    now,
        UpdatedAt:    now,
    }, nil
}

// Complete marks the PM as completed and schedules next.
// Complete ทำเครื่องหมาย PM เสร็จและกำหนดครั้งถัดไป
func (pm *PMSchedule) Complete() error {
    if pm.Status == PMStatusCompleted {
        return errors.New("PM already completed")
    }

    now := time.Now()
    pm.LastDoneAt = &now
    pm.NextDueAt = now.AddDate(0, 0, pm.IntervalDays)
    pm.Status = PMStatusScheduled
    pm.UpdatedAt = now
    return nil
}

// CheckDue updates the status based on due date.
// CheckDue อัปเดตสถานะตามวันครบกำหนด
func (pm *PMSchedule) CheckDue() {
    now := time.Now()
    if now.After(pm.NextDueAt) {
        pm.Status = PMStatusOverdue
    } else if now.AddDate(0, 0, 3).After(pm.NextDueAt) {
        pm.Status = PMStatusDue
    }
}
```

---

# 23. Dual-Domain Architecture

## 23.1 Domain Isolation Strategy

```
┌─────────────────────────────────────────────────────────────────┐
│                    DUAL-DOMAIN ISOLATION                         │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────┐     ┌─────────────────────────┐
│      FARM DOMAIN        │     │     FACTORY DOMAIN       │
│                         │     │                          │
│  ┌─────────────────┐   │     │  ┌─────────────────┐    │
│  │ erp/inventory   │   │     │  │ mes/production  │    │
│  ├─────────────────┤   │     │  ├─────────────────┤    │
│  │ erp/production  │   │     │  │ mes/bom         │    │
│  ├─────────────────┤   │     │  ├─────────────────┤    │
│  │ crm/crm         │   │     │  │ qms/quality     │    │
│  ├─────────────────┤   │     │  ├─────────────────┤    │
│  │ iot/iot         │   │     │  │ scada/gateway   │    │
│  ├─────────────────┤   │     │  ├─────────────────┤    │
│  │ forecast        │   │     │  │ costing         │    │
│  └─────────────────┘   │     │  └─────────────────┘    │
│                         │     │                          │
│  Database: farm_db      │     │  Database: factory_db    │
│  Kafka: farm.*          │     │  Kafka: factory.*        │
└────────────┬────────────┘     └────────────┬────────────┘
             │                               │
             │      SHARED KAFKA TOPICS      │
             │  ┌─────────────────────────┐  │
             └──│ cross-domain.events     │──┘
                │ - farm.harvest.completed│
                │ - factory.production.   │
                │   completed             │
                │ - factory.raw.received  │
                └─────────────────────────┘
```

## 23.2 Shared Core

```go
// internal/modules/shared/core/module.go
package core

import (
    "database/sql"

    "github.com/gin-gonic/gin"
    "github.com/yourorg/smart-erp/internal/modules/shared/auth"
    "github.com/yourorg/smart-erp/internal/modules/shared/tenant"
    "github.com/yourorg/smart-erp/internal/modules/shared/user"
    "github.com/yourorg/smart-erp/internal/modules/shared/notification"
    "github.com/yourorg/smart-erp/internal/modules/shared/document"
    "github.com/yourorg/smart-erp/internal/modules/shared/workflow"
    "github.com/yourorg/smart-erp/internal/modules/shared/audit"
)

// SharedCore represents shared modules across domains.
// SharedCore แทนโมดูลที่ใช้ร่วมกันระหว่าง domains
type SharedCore struct {
    Auth         *auth.Module
    Tenant       *tenant.Module
    User         *user.Module
    Notification *notification.Module
    Document     *document.Module
    Workflow     *workflow.Module
    Audit        *audit.Module
}

// NewSharedCore initializes shared modules.
// NewSharedCore เริ่มต้นโมดูลที่ใช้ร่วมกัน
func NewSharedCore(db *sql.DB) *SharedCore {
    return &SharedCore{
        Auth:         auth.NewModule(db),
        Tenant:       tenant.NewModule(db),
        User:         user.NewModule(db),
        Notification: notification.NewModule(db),
        Document:     document.NewModule(db),
        Workflow:     workflow.NewModule(db),
        Audit:        audit.NewModule(db),
    }
}

// RegisterRoutes registers shared routes.
// RegisterRoutes ลงทะเบียน routes ที่ใช้ร่วมกัน
func (s *SharedCore) RegisterRoutes(router *gin.RouterGroup) {
    s.Auth.RegisterRoutes(router.Group("/auth"))
    s.Tenant.RegisterRoutes(router.Group("/tenants"))
    s.User.RegisterRoutes(router.Group("/users"))
    s.Notification.RegisterRoutes(router.Group("/notifications"))
    s.Document.RegisterRoutes(router.Group("/documents"))
    s.Workflow.RegisterRoutes(router.Group("/workflows"))
    s.Audit.RegisterRoutes(router.Group("/audit"))
}
```

## 23.3 Cross-Domain Event Bus

```go
// internal/shared/events/cross_domain_bus.go
package events

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/segmentio/kafka-go"
)

// CrossDomainEventBus handles cross-domain events.
// CrossDomainEventBus จัดการ events ระหว่าง domains
type CrossDomainEventBus struct {
    writer *kafka.Writer
    reader *kafka.Reader
}

// NewCrossDomainEventBus creates a new cross-domain event bus.
// NewCrossDomainEventBus สร้าง event bus ระหว่าง domains ใหม่
func NewCrossDomainEventBus(brokers []string) *CrossDomainEventBus {
    return &CrossDomainEventBus{
        writer: &kafka.Writer{
            Addr:     kafka.TCP(brokers...),
            Balancer: &kafka.LeastBytes{},
        },
    }
}

// PublishFarmEvent publishes a farm event.
// PublishFarmEvent publish event จากฟาร์ม
func (b *CrossDomainEventBus) PublishFarmEvent(ctx context.Context, eventType string, event interface{}) error {
    topic := fmt.Sprintf("farm.%s", eventType)
    return b.publish(ctx, topic, event)
}

// PublishFactoryEvent publishes a factory event.
// PublishFactoryEvent publish event จากโรงงาน
func (b *CrossDomainEventBus) PublishFactoryEvent(ctx context.Context, eventType string, event interface{}) error {
    topic := fmt.Sprintf("factory.%s", eventType)
    return b.publish(ctx, topic, event)
}

// SubscribeFarmEvents subscribes to farm events.
// SubscribeFarmEvents subscribe events จากฟาร์ม
func (b *CrossDomainEventBus) SubscribeFarmEvents(ctx context.Context, eventType string, handler func([]byte) error) error {
    topic := fmt.Sprintf("farm.%s", eventType)
    return b.subscribe(ctx, topic, handler)
}

// SubscribeFactoryEvents subscribes to factory events.
// SubscribeFactoryEvents subscribe events จากโรงงาน
func (b *CrossDomainEventBus) SubscribeFactoryEvents(ctx context.Context, eventType string, handler func([]byte) error) error {
    topic := fmt.Sprintf("factory.%s", eventType)
    return b.subscribe(ctx, topic, handler)
}

func (b *CrossDomainEventBus) publish(ctx context.Context, topic string, event interface{}) error {
    payload, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("marshal event: %w", err)
    }

    return b.writer.WriteMessages(ctx, kafka.Message{
        Topic: topic,
        Key:   []byte(time.Now().Format(time.RFC3339Nano)),
        Value: payload,
    })
}

func (b *CrossDomainEventBus) subscribe(ctx context.Context, topic string, handler func([]byte) error) error {
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{"localhost:9092"},
        Topic:   topic,
        GroupID: fmt.Sprintf("%s-group", topic),
    })
    defer reader.Close()

    for {
        msg, err := reader.ReadMessage(ctx)
        if err != nil {
            return err
        }
        if err := handler(msg.Value); err != nil {
            // Log error
            continue
        }
    }
}
```

## 23.4 Cross-Domain Service

```go
// internal/modules/factory/mes/production/application/service/cross_domain_service.go
package service

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/yourorg/smart-erp/internal/shared/events"
)

// CrossDomainService handles cross-domain integration.
// CrossDomainService จัดการ integration ระหว่าง domains
type CrossDomainService struct {
    eventBus *events.CrossDomainEventBus
}

// NewCrossDomainService creates a new cross-domain service.
// NewCrossDomainService สร้าง cross-domain service ใหม่
func NewCrossDomainService(eventBus *events.CrossDomainEventBus) *CrossDomainService {
    return &CrossDomainService{eventBus: eventBus}
}

// StartConsuming starts consuming cross-domain events.
// StartConsuming เริ่ม consume events ระหว่าง domains
func (s *CrossDomainService) StartConsuming(ctx context.Context) error {
    // Subscribe to farm harvest events
    // Subscribe events การเก็บเกี่ยวจากฟาร์ม
    go func() {
        err := s.eventBus.SubscribeFarmEvents(ctx, "harvest.completed", func(payload []byte) error {
            var event events.HarvestCompletedEvent
            if err := json.Unmarshal(payload, &event); err != nil {
                return fmt.Errorf("unmarshal harvest event: %w", err)
            }

            // Process: create raw material receipt in factory
            // ประมวลผล: สร้างใบรับวัตถุดิบในโรงงาน
            return s.handleHarvestCompleted(ctx, event)
        })
        if err != nil {
            // Log error
        }
    }()

    return nil
}

// handleHarvestCompleted handles harvest completed event.
// handleHarvestCompleted จัดการ event การเก็บเกี่ยวเสร็จ
func (s *CrossDomainService) handleHarvestCompleted(ctx context.Context, event events.HarvestCompletedEvent) error {
    // 1. Create raw material receipt
    // 1. สร้างใบรับวัตถุดิบ
    // 2. Update inventory
    // 2. อัปเดตสต็อก
    // 3. Notify factory
    // 3. แจ้งโรงงาน
    return nil
}
```

---

# 24. Cross-Domain Integration

## 24.1 Integration Patterns

```
┌─────────────────────────────────────────────────────────────────┐
│              CROSS-DOMAIN INTEGRATION PATTERNS                   │
└─────────────────────────────────────────────────────────────────┘

Pattern 1: Event-Driven (Recommended)
┌──────────┐     ┌──────────┐     ┌──────────┐
│  Farm    │────▶│  Kafka   │────▶│ Factory  │
│  Event   │     │  Topic   │     │ Handler  │
└──────────┘     └──────────┘     └──────────┘

Pattern 2: API Gateway
┌──────────┐     ┌──────────┐     ┌──────────┐
│  Farm    │────▶│  API     │────▶│ Factory  │
│  Client  │     │  Gateway │     │  API     │
└──────────┘     └──────────┘     └──────────┘

Pattern 3: Shared Database (Anti-pattern)
┌──────────┐     ┌──────────┐     ┌──────────┐
│  Farm    │────▶│  Shared  │◀────│ Factory  │
│  Service │     │  DB      │     │ Service  │
└──────────┘     └──────────┘     └──────────┘
                 ❌ ไม่แนะนำ

Pattern 4: Saga Pattern (Complex flows)
┌──────────┐     ┌──────────┐     ┌──────────┐
│  Farm    │────▶│  Saga    │────▶│ Factory  │
│  Step 1  │     │Orchestr. │     │ Step 2   │
└──────────┘     └──────────┘     └──────────┘
                       │
                       ▼
                 ┌──────────┐
                 │Compensate│
                 │on failure│
                 └──────────┘
```

## 24.2 Traceability Chain

```go
// internal/modules/shared/tracing/traceability.go
package tracing

import (
    "context"
    "time"

    "github.com/shopspring/decimal"
)

// TraceRecord represents a traceability record.
// TraceRecord แทนบันทึกการติดตามย้อนกลับ
type TraceRecord struct {
    ID          string          `json:"id"`           // รหัสบันทึก
    TenantID    string          `json:"tenant_id"`    // รหัสองค์กร
    BatchID     string          `json:"batch_id"`     // รหัส批次
    Stage       string          `json:"stage"`        // ขั้นตอน
    Location    string          `json:"location"`     // สถานที่
    FarmID      string          `json:"farm_id"`      // รหัสฟาร์ม
    FactoryID   string          `json:"factory_id"`   // รหัสโรงงาน
    CustomerID  string          `json:"customer_id"`  // รหัสลูกค้า
    ProductID   string          `json:"product_id"`   // รหัสสินค้า
    Quantity    decimal.Decimal `json:"quantity"`     // จำนวน
    Unit        string          `json:"unit"`         // หน่วย
    Quality     string          `json:"quality"`      // เกรด
    OperatorID  string          `json:"operator_id"`  // รหัสผู้ปฏิบัติ
    Timestamp   time.Time       `json:"timestamp"`    // เวลาที่บันทึก
    Metadata    map[string]interface{} `json:"metadata"` // ข้อมูลเพิ่มเติม
}

// TraceabilityService provides traceability operations.
// TraceabilityService ให้บริการ operations สำหรับ traceability
type TraceabilityService struct {
    repo TraceRepository
}

// TraceRepository defines the repository interface.
// TraceRepository กำหนด interface ของ repository
type TraceRepository interface {
    Save(ctx context.Context, record *TraceRecord) error
    FindByBatchID(ctx context.Context, tenantID, batchID string) ([]*TraceRecord, error)
    FindByFarmID(ctx context.Context, tenantID, farmID string) ([]*TraceRecord, error)
    FindByFactoryID(ctx context.Context, tenantID, factoryID string) ([]*TraceRecord, error)
    FindByCustomerID(ctx context.Context, tenantID, customerID string) ([]*TraceRecord, error)
}

// NewTraceabilityService creates a new traceability service.
// NewTraceabilityService สร้าง traceability service ใหม่
func NewTraceabilityService(repo TraceRepository) *TraceabilityService {
    return &TraceabilityService{repo: repo}
}

// RecordHarvest records a harvest event.
// RecordHarvest บันทึก event การเก็บเกี่ยว
func (s *TraceabilityService) RecordHarvest(ctx context.Context, record *TraceRecord) error {
    record.Stage = "harvest"
    record.Timestamp = time.Now()
    return s.repo.Save(ctx, record)
}

// RecordProduction records a production event.
// RecordProduction บันทึก event การผลิต
func (s *TraceabilityService) RecordProduction(ctx context.Context, record *TraceRecord) error {
    record.Stage = "production"
    record.Timestamp = time.Now()
    return s.repo.Save(ctx, record)
}

// RecordDelivery records a delivery event.
// RecordDelivery บันทึก event การส่งของ
func (s *TraceabilityService) RecordDelivery(ctx context.Context, record *TraceRecord) error {
    record.Stage = "delivery"
    record.Timestamp = time.Now()
    return s.repo.Save(ctx, record)
}

// GetFullTrace returns the full traceability chain.
// GetFullTrace คืนค่า chain การติดตามย้อนกลับทั้งหมด
func (s *TraceabilityService) GetFullTrace(ctx context.Context, tenantID, batchID string) ([]*TraceRecord, error) {
    return s.repo.FindByBatchID(ctx, tenantID, batchID)
}
```

## 24.3 Cross-Domain Dashboard

```go
// internal/modules/shared/dashboard/cross_domain.go
package dashboard

import (
    "context"

    "github.com/shopspring/decimal"
)

// CrossDomainDashboard provides cross-domain KPIs.
// CrossDomainDashboard ให้ KPI ระหว่าง domains
type CrossDomainDashboard struct {
    // Dependencies
}

// FarmToFactoryMetrics represents farm-to-factory metrics.
// FarmToFactoryMetrics แทน metrics จากฟาร์มถึงโรงงาน
type FarmToFactoryMetrics struct {
    TotalHarvested     decimal.Decimal // ผลผลิตรวม
    TotalReceived      decimal.Decimal // รับรวม
    TotalProduced      decimal.Decimal // ผลิตรวม
    TotalDelivered     decimal.Decimal // ส่งรวม
    YieldRate          decimal.Decimal // อัตราผลได้
    LeadTimeAvg        decimal.Decimal // เวลานำเฉลี่ย (วัน)
    TraceabilityRate   decimal.Decimal // อัตรา traceability
}

// GetFarmToFactoryMetrics returns farm-to-factory metrics.
// GetFarmToFactoryMetrics คืนค่า metrics จากฟาร์มถึงโรงงาน
func (d *CrossDomainDashboard) GetFarmToFactoryMetrics(ctx context.Context, tenantID string) (*FarmToFactoryMetrics, error) {
    // Query farm harvest
    // Query factory receipt
    // Query factory production
    // Query delivery
    // Calculate metrics
    return &FarmToFactoryMetrics{
        YieldRate:        decimal.NewFromFloat(92.5),
        LeadTimeAvg:      decimal.NewFromFloat(3.2),
        TraceabilityRate: decimal.NewFromInt(100),
    }, nil
}

// FactoryToCustomerMetrics represents factory-to-customer metrics.
// FactoryToCustomerMetrics แทน metrics จากโรงงานถึงลูกค้า
type FactoryToCustomerMetrics struct {
    TotalProduced     decimal.Decimal // ผลิตรวม
    TotalDelivered    decimal.Decimal // ส่งรวม
    OnTimeDelivery    decimal.Decimal // ส่งตรงเวลา
    QualityRate       decimal.Decimal // อัตราคุณภาพ
    CustomerSatisfaction decimal.Decimal // ความพึงพอใจลูกค้า
    ReturnRate        decimal.Decimal // อัตราการคืน
}

// GetFactoryToCustomerMetrics returns factory-to-customer metrics.
// GetFactoryToCustomerMetrics คืนค่า metrics จากโรงงานถึงลูกค้า
func (d *CrossDomainDashboard) GetFactoryToCustomerMetrics(ctx context.Context, tenantID string) (*FactoryToCustomerMetrics, error) {
    return &FactoryToCustomerMetrics{
        OnTimeDelivery:       decimal.NewFromFloat(98.5),
        QualityRate:          decimal.NewFromFloat(99.2),
        CustomerSatisfaction: decimal.NewFromFloat(4.7),
        ReturnRate:           decimal.NewFromFloat(0.8),
    }, nil
}
```

## 24.4 Unified KPI

```go
// internal/modules/shared/kpi/unified_kpi.go
package kpi

import (
    "github.com/shopspring/decimal"
)

// UnifiedKPI represents unified KPIs across domains.
// UnifiedKPI แทน KPI รวมระหว่าง domains
type UnifiedKPI struct {
    // Farm KPIs
    // KPI ฟาร์ม
    FarmYield           decimal.Decimal // ผลผลิตฟาร์ม
    FarmCostPerKg       decimal.Decimal // ต้นทุนต่อ กก.
    FarmQualityRate     decimal.Decimal // อัตราคุณภาพ
    FarmOnTimeDelivery  decimal.Decimal // ส่งตรงเวลา

    // Factory KPIs
    // KPI โรงงาน
    FactoryOEE          decimal.Decimal // OEE
    FactoryYield        decimal.Decimal // อัตราผลได้
    FactoryDefectRate   decimal.Decimal // อัตราของเสีย
    FactoryOnTimeDelivery decimal.Decimal // ส่งตรงเวลา

    // Cross-Domain KPIs
    // KPI ระหว่าง domains
    FarmToFactoryYield  decimal.Decimal // ผลได้ ฟาร์ม → โรงงาน
    FarmToFactoryLeadTime decimal.Decimal // เวลานำ
    EndToEndTraceability decimal.Decimal // Traceability

    // Financial KPIs
    // KPI การเงิน
    TotalRevenue        decimal.Decimal // รายได้รวม
    TotalCost           decimal.Decimal // ต้นทุนรวม
    GrossMargin         decimal.Decimal // กำไรขั้นต้น
    NetMargin           decimal.Decimal // กำไรสุทธิ
}

// CalculateUnifiedKPI calculates unified KPIs.
// CalculateUnifiedKPI คำนวณ KPI รวม
func CalculateUnifiedKPI(farmData, factoryData, financialData map[string]decimal.Decimal) *UnifiedKPI {
    return &UnifiedKPI{
        FarmYield:           farmData["yield"],
        FarmCostPerKg:       farmData["cost_per_kg"],
        FarmQualityRate:     farmData["quality_rate"],
        FarmOnTimeDelivery:  farmData["on_time_delivery"],
        FactoryOEE:          factoryData["oee"],
        FactoryYield:        factoryData["yield"],
        FactoryDefectRate:   factoryData["defect_rate"],
        FactoryOnTimeDelivery: factoryData["on_time_delivery"],
        FarmToFactoryYield:  farmData["yield"].Mul(factoryData["yield"]),
        TotalRevenue:        financialData["revenue"],
        TotalCost:           financialData["cost"],
        GrossMargin:         financialData["revenue"].Sub(financialData["cost"]),
    }
}
```

---

# 📎 ภาคผนวก

## A. เอกสารอ้างอิง

- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design - Eric Evans](https://domainlanguage.com/ddd/)
- [MES Best Practices - MESA International](https://www.mesa.org/)
- [ISA-95 Standard](https://www.isa.org/standards-and-publications/isa-standards/isa-95-standard)
- [OEE Foundation](https://www.oee.com/)
- [Go Standard Project Layout](https://github.com/golang-standards/project-layout)
- [12 Factor App](https://12factor.net/)

## B. เครื่องมือที่ใช้

| ประเภท | Farm | Factory |
|---|---|---|
| **ภาษา** | Go 1.21+ | Go 1.21+ |
| **Web Framework** | Gin | Gin |
| **Database** | PostgreSQL 15 | PostgreSQL 15 |
| **Cache** | Redis 7 | Redis 7 |
| **Message Queue** | Kafka | Kafka |
| **Time-series** | ClickHouse | ClickHouse |
| **Search** | Elasticsearch | Elasticsearch |
| **IoT Protocol** | MQTT | MQTT |
| **Industrial Protocol** | - | OPC-UA, Modbus |
| **Container** | Docker | Docker |
| **Orchestration** | Kubernetes | Kubernetes |
| **CI/CD** | GitHub Actions | GitHub Actions |
| **Monitoring** | Prometheus + Grafana | Prometheus + Grafana |
| **Tracing** | Jaeger | Jaeger |
| **Logging** | Loki | Loki |
| **Testing** | Testify, Testcontainers | Testify, Testcontainers |
| **Load Test** | k6 | k6 + JMeter |

## C. ติดต่อ

- **Repository**: https://github.com/yourorg/smart-erp
- **Documentation**: https://docs.smart-erp.com
- **Support**: support@smart-erp.com

---

**เวอร์ชัน**: 2.0 (Extended Edition)  
**วันที่**: 2024-01-15  
**ผู้จัดทำ**: Architecture Team  
**ประเภทเอกสาร**: Application Design Document — Dual-Domain (Farm + Factory)  
**ขอบเขต**: ERP+SaaS สำหรับฟาร์มเห็ดอัจฉริยะ และโรงงาน SME

---

*เอกสารนี้เป็นต้นแบบสำหรับให้ AI สร้าง/แก้ไขโปรแกรมภาษา Go ทุกประเภท ใช้โครงสร้าง Clean Architecture + DDD พร้อม Layer ที่ชัดเจน, Cross-cutting Concerns ครบ, และ Pattern แยก Modules รองรับทั้ง Farm Domain และ Factory Domain พร้อม Cross-Domain Integration ผ่าน Event-Driven Architecture*