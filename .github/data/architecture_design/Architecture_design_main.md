# 🗺️ สรุปทุก Modules และการเชื่อมโยงกัน — `icmongolang IoT Platform`

> **เอกสารนี้เป็น Master Summary** ของทั้ง 40 modules + แสดงความสัมพันธ์ทุกรูปแบบ (sync/async/data/event)

---

## 📑 สารบัญ

1. [ภาพรวมระบบ](#1-ภาพรวมระบบ)
2. [ทะเบียน Modules ทั้ง 40](#2-ทะเบียน-modules-ทั้ง-40)
3. [Bounded Context Map](#3-bounded-context-map)
4. [Dependency Graph (Layer 1-5)](#4-dependency-graph-layer-1-5)
5. [การเชื่อมโยง 4 รูปแบบ](#5-การเชื่อมโยง-4-รูปแบบ)
6. [Data Flow หลัก 8 เส้นทาง](#6-data-flow-หลัก-8-เส้นทาง)
7. [Event Matrix (Kafka)](#7-event-matrix-kafka)
8. [Sync API Matrix](#8-sync-api-matrix)
9. [Shared Data (DB/Redis/Influx/ES)](#9-shared-data)
10. [Cross-Cutting Concerns](#10-cross-cutting-concerns)
11. [Deployment Topology](#11-deployment-topology)
12. [Roadmap & Production Order](#12-roadmap--production-order)

---

## 1. ภาพรวมระบบ

### 1.1 แพลตฟอร์มทำอะไร

```
┌─────────────────────────────────────────────────────────────────┐
│              icmongolang IoT Platform (40 modules)              │
├─────────────────────────────────────────────────────────────────┤
│  🎯 Business Layer                                              │
│  ┌───────┬───────┬───────┬───────┬───────┬───────┬───────┐    │
│  │customer│package│  erp  │  crm  │logistics│report│ pdpa  │    │
│  └───────┴───────┴───────┴───────┴───────┴───────┴───────┘    │
├─────────────────────────────────────────────────────────────────┤
│  ⚙️ IoT Core Layer                                              │
│  ┌───────┬───────┬───────┬───────┬───────┬───────┐            │
│  │device │  iot  │ mqtt  │ kafka │influx │  ws   │            │
│  └───────┴───────┴───────┴───────┴───────┴───────┘            │
├─────────────────────────────────────────────────────────────────┤
│  🔧 Infrastructure Layer                                        │
│  ┌───────┬───────┬───────┬───────┬───────┬───────┐            │
│  │ auth  │ users │payment│notifier│ alarm │ audit │            │
│  └───────┴───────┴───────┴───────┴───────┴───────┘            │
├─────────────────────────────────────────────────────────────────┤
│  📦 Shared Layer                                                │
│  ┌───────┬───────┬───────┬───────┬───────┬───────┐            │
│  │  i18n │ api-mgr│ flow │control│schedule│ queue│            │
│  └───────┴───────┴───────┴───────┴───────┴───────┘            │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 Core Capabilities

| Capability | Modules ที่ให้บริการ |
| :--- | :--- |
| **IoT Ingestion** | `device`, `iot`, `mqtt`, `kafka`, `influxdb` |
| **Real-time Control** | `device`, `control`, `flowengine`, `websocket` |
| **Business Ops** | `customer`, `package`, `erp`, `crm`, `logistics` |
| **Analytics** | `report`, `dashboard`, `vectordata`, `elasticsearch` |
| **Compliance** | `pdpa`, `auditlog`, `document` |
| **Security** | `auth`, `users`, `apimanager`, `settings` |

---

## 2. ทะเบียน Modules ทั้ง 40

### 2.1 โมดูลใหม่ — Greenfield (7)

| # | Module | Context | Priority | Status | LOC (est.) |
| :-: | :--- | :--- | :-: | :-: | :-: |
| 01 | `customer` | Customer Mgmt | P0 | ✅ Done | 2,500 |
| 02 | `package` | Subscription | P0 | ✅ Done | 2,800 |
| 03 | `device` | IoT Core | P0 | ✅ Done | 4,500 |
| 04 | `erp` | Enterprise Resource | P1 | ⏳ | 3,800 |
| 05 | `crm` | Sales & Support | P1 | ⏳ | 3,200 |
| 06 | `logistics` | Supply Chain | P2 | ⏳ | 2,800 |
| 07 | `report` | Analytics | P0 | 🟡 Upgrade | 3,500 |

### 2.2 โมดูลเดิม — Core Infrastructure (8)

| # | Module | สถานะ | สิ่งที่เพิ่ม | Priority |
| :-: | :--- | :-: | :--- | :-: |
| 08 | `auth` | 🟡 Upgrade | tenant, SSO, MFA, session | P0 |
| 09 | `users` | 🟡 Upgrade | multi-tenant, RBAC+ABAC | P0 |
| 10 | `settings` | 🟡 Upgrade | per-tenant, hierarchical | P0 |
| 11 | `payment` | 🟡 Upgrade | subscription, QR, multi-currency | P0 |
| 12 | `notifier` | 🟡 Upgrade | multi-channel (LINE, Push) | P0 |
| 13 | `alarm` | 🟡 Upgrade | rule engine, auto-command | P0 |
| 14 | `pdpa` | 🟡 Upgrade | IoT consent, retention | P1 |
| 15 | `auditlog` | 🟡 Upgrade | ES search, tamper-proof | P1 |

### 2.3 โมดูลเดิม — Business (9)

| # | Module | สถานะ | เชื่อมกับ |
| :-: | :--- | :-: | :--- |
| 16 | `dashboard` | 🟡 Upgrade | device, report, WS |
| 17 | `items` | 🟡 Upgrade | erp, purchaseorder |
| 18 | `purchaseorder` | 🟡 Upgrade | erp, items |
| 19 | `quotation` | 🟡 Upgrade | crm, erp |
| 20 | `wos` | 🟡 Upgrade | device, logistics |
| 21 | `document` | 🟡 Upgrade | crm, erp, pdpa |
| 22 | `email` | 🟡 Upgrade | notifier, report |
| 23 | `batch` | 🟡 Upgrade | device (telemetry batch) |
| 24 | `job` | 🟡 Upgrade | scheduler |

### 2.4 โมดูลเดิม — Integration (8)

| # | Module | สถานะ | บทบาท |
| :-: | :--- | :-: | :--- |
| 25 | `iot` | 🟡 Deprecate | migrate → device |
| 26 | `mqtt` | 🟡 Upgrade | multi-protocol, TLS, auth |
| 27 | `kafka` | 🟡 Upgrade | IoT topics, DLQ, schema |
| 28 | `influxdb` | 🟡 Upgrade | retention, downsampling |
| 29 | `elasticsearch` | 🟡 Upgrade | IoT index, ILM |
| 30 | `websocket` | 🟡 Upgrade | multi-tenant hub, channels |
| 31 | `realtime` | 🟡 Merge | รวมกับ websocket |
| 32 | `vectordata` | 🟡 Upgrade | embedding pipeline |

### 2.5 โมดูลเดิม — Utility (8)

| # | Module | สถานะ | บทบาท |
| :-: | :--- | :-: | :--- |
| 33 | `i18n` | 🟡 Upgrade | IoT terms, per-tenant locale |
| 34 | `apimanager` | 🟡 Upgrade | API key, rate limit, usage |
| 35 | `flowengine` | 🟡 Upgrade | visual flow, sub-flow |
| 36 | `control` | 🟡 Upgrade | command templates, batch |
| 37 | `fullschedule` | 🟡 Upgrade | calendar, holiday-aware |
| 38 | `queue` | 🟡 Upgrade | priority, delayed, DLQ |
| 39 | `notifier` | (ดู #12) | – |
| 40 | `worker` | 🟡 Upgrade | task orchestration |

### 2.6 สรุปตัวเลข

| กลุ่ม | จำนวน | Done | Upgrade | Greenfield |
| :--- | :-: | :-: | :-: | :-: |
| Greenfield | 7 | 3 | 1 | 3 |
| Core | 8 | 0 | 8 | 0 |
| Business | 9 | 0 | 9 | 0 |
| Integration | 8 | 0 | 8 | 0 |
| Utility | 8 | 0 | 8 | 0 |
| **รวม** | **40** | **3** | **34** | **3** |

---

## 3. Bounded Context Map

### 3.1 Strategic Context Map

```
                    ┌──────────────────────────────────────┐
                    │         Core Domain (สำคัญสุด)       │
                    │  ┌──────────┐  ┌──────────┐         │
                    │  │  device  │  │ package  │         │
                    │  └──────────┘  └──────────┘         │
                    └──────────────────────────────────────┘
                                       ▲
                                       │
        ┌──────────────────────────────┼──────────────────────────────┐
        │                              │                              │
        ▼                              ▼                              ▼
┌──────────────┐              ┌──────────────┐              ┌──────────────┐
│   Customer   │              │   Business   │              │  Compliance  │
│   Context    │              │   Context    │              │   Context    │
├──────────────┤              ├──────────────┤              ├──────────────┤
│  customer    │              │ แ   │              │  pdpa        │
│              │              │  logistics   │              │  auditlog    │
│              │              │  quotation   │              │  document    │
│              │              │  purchase    │              │              │
│              │              │  wos         │              │              │
└──────────────┘              └──────────────┘              └──────────────┘
        │                              │                              │
        └──────────────────────────────┼──────────────────────────────┘
                                       │
                                       ▼
                    ┌──────────────────────────────────────┐
                    │       Supporting Subdomain           │
                    ├──────────────────────────────────────┤
                    │  auth, users, settings, payment      │
                    │  notifier, alarm, dashboard, report  │
                    │  items, batch, job                   │
                    └──────────────────────────────────────┘
                                       │
                                       ▼
                    ┌──────────────────────────────────────┐
                    │       Generic Subdomain              │
                    ├──────────────────────────────────────┤
                    │  i18n, apimanager, queue             │
                    │  flowengine, control, fullschedule   │
                    └──────────────────────────────────────┘

                    ┌──────────────────────────────────────┐
                    │       Technical Infrastructure       │
                    ├──────────────────────────────────────┤
                    │  iot, mqtt, kafka, influxdb          │
                    │  elasticsearch, websocket, realtime  │
                    │  vectordata                          │
                    └──────────────────────────────────────┘
```

### 3.2 Context Relationship Types

| Upstream (U) | Downstream (D) | Pattern | ตัวอย่าง |
| :--- | :--- | :--- | :--- |
| `customer` | `package` | Customer-Supplier | ลูกค้าสมัครแพ็กเกจ |
| `customer` | `crm` | Customer-Supplier | Lead → Customer |
| `customer` | `erp` | Conformist | ERP ใช้ customer_id |
| `customer` | `logistics` | Customer-Supplier | จัดส่งให้ลูกค้า |
| `customer` | `device` | Partnership | Device ผูก Site |
| `package` | `payment` | Partnership | Subscription → Payment |
| `package` | `device` | Customer-Supplier | Quota จำกัด device |
| `package` | `erp` | Partnership | Billing → Invoice |
| `device` | `alarm` | Published Language | Telemetry → Alert |
| `device` | `report` | Published Language | Telemetry → Analytics |
| `device` | `logistics` | Customer-Supplier | GPS tracker |
| `crm` | `quotation` | Partnership | Opportunity → Quote |
| `erp` | `logistics` | Partnership | Order → Shipment |
| `erp` | `payment` | Partnership | Invoice → Payment |
| `pdpa` | ทุกโมดูล | Conformist | Consent ทุกที่ |

---

## 4. Dependency Graph (Layer 1-5)

### 4.1 ภาพรวม 5 Layers

```
┌─────────────────────────────────────────────────────────────────┐
│ Layer 5 — User-Facing                                           │
│  dashboard · report · websocket · apimanager                    │
└─────────────────────────────────────────────────────────────────┘
                              ▲
┌─────────────────────────────────────────────────────────────────┐
│ Layer 4 — Business Domain                                       │
│  crm · erp · logistics · quotation · purchaseorder · wos        │
└─────────────────────────────────────────────────────────────────┘
                              ▲
┌─────────────────────────────────────────────────────────────────┐
│ Layer 3 — Core Domain                                           │
│  customer · package · device · payment · pdpa                   │
└─────────────────────────────────────────────────────────────────┘
                              ▲
┌─────────────────────────────────────────────────────────────────┐
│ Layer 2 — Platform Services                                     │
│  auth · users · settings · notifier · alarm · auditlog · job    │
│  i18n · email · document · batch · queue · flowengine           │
│  control · fullschedule                                         │
└─────────────────────────────────────────────────────────────────┘
                              ▲
┌─────────────────────────────────────────────────────────────────┐
│ Layer 1 — Infrastructure                                        │
│  kafka · mqtt · influxdb · elasticsearch · vectordata · realtime│
│  iot (deprecated)                                               │
└─────────────────────────────────────────────────────────────────┘
```

### 4.2 Dependency Rules

| Layer | ห้าม import | อนุญาต import |
| :--- | :--- | :--- |
| Layer 5 | – | L1-L4 |
| Layer 4 | L5 | L1-L3 |
| Layer 3 | L4, L5 | L1-L2 |
| Layer 2 | L3-L5 | L1 |
| Layer 1 | L2-L5 | stdlib เท่านั้น |

### 4.3 Dependency Matrix (Top 10 modules)

| Module | Depends on | Depended by |
| :--- | :--- | :--- |
| `customer` | auth, users, auditlog | package, crm, erp, logistics, device |
| `package` | customer, payment, notifier | device, erp, crm |
| `device` | customer, package, mqtt, kafka, influxdb, ws | alarm, report, logistics, dashboard |
| `erp` | customer, items, payment | logistics, report, wos |
| `crm` | customer, notifier, email | quotation, report |
| `logistics` | customer, erp, device, notifier | report, wos |
| `auth` | users, settings | ทุกโมดูล |
| `mqtt` | – | device, iot |
| `kafka` | – | ทุกโมดูลที่ publish |
| `influxdb` | – | device, report, dashboard |

---

## 5. การเชื่อมโยง 4 รูปแบบ

### 5.1 รูปแบบที่ 1 — Sync API (REST)

```
Module A ──HTTP/REST──► Module B
         (internal)
```

| จาก | ไป | Endpoint | ใช้เมื่อ |
| :--- | :--- | :--- | :--- |
| `device` | `customer` | `GET /internal/customers/:id` | ตรวจ customer active |
| `device` | `package` | `POST /internal/quota/check` | ตรวจโควตา |
| `erp` | `customer` | `GET /internal/customers/:id` | ดึงข้อมูลลูกค้า |
| `erp` | `items` | `GET /internal/items/:id` | ดึงข้อมูลสินค้า |
| `logistics` | `erp` | `GET /internal/orders/:id` | ดึงข้อมูล order |
| `report` | `device` | `GET /internal/devices/:id` | enrich report |

### 5.2 รูปแบบที่ 2 — Async Event (Kafka)

```
Module A ──publish──► Kafka Topic ──consume──► Module B, C, D
                    (fan-out)
```

ตัวอย่าง:
```
customer ──publish──► customer.customer.activated
                              │
                              ├──consume──► package (สร้าง subscription)
                              ├──consume──► crm (update lead)
                              └──consume──► notifier (ส่ง welcome email)
```

### 5.3 รูปแบบที่ 3 — Data Share (DB/Redis/Influx)

```
Module A ──read/write──► Shared DB/Redis/InfluxDB ◄──read/write── Module B
```

| Storage | ใช้โดย | ตัวอย่าง |
| :--- | :--- | :--- |
| **Postgres** | ทุกโมดูล | แต่ละโมดูลมี schema แยก (prefix) |
| **Redis** | device, auth, customer | cache, session, rate limit |
| **InfluxDB** | device (write), report/dashboard (read) | telemetry |
| **Elasticsearch** | device, auditlog, customer (write), report (read) | search |
| **S3** | document, report, pdpa | file storage |

### 5.4 รูปแบบที่ 4 — Realtime (WebSocket)

```
Module A ──broadcast──► WS Hub ──push──► Client (browser/app)
```

ตัวอย่าง:
```
device ──telemetry──► WS Hub ──push──► Frontend
alarm  ──alert──────► WS Hub ──push──► Frontend
crm    ──ticket─────► WS Hub ──push──► Agent Dashboard
```

---

## 6. Data Flow หลัก 8 เส้นทาง

### Flow 1: Customer Onboarding

```
┌──────────┐  1.signup   ┌──────────┐  2.create   ┌──────────┐
│  Visitor │────────────►│  auth    │────────────►│ customer │
└──────────┘             └──────────┘             └────┬─────┘
                                                       │
                                          3.publish: customer.customer.created
                                                       │
                          ┌────────────────────────────┼────────────────┐
                          ▼                            ▼                ▼
                   ┌──────────┐              ┌──────────┐      ┌──────────┐
                   │ package  │              │  crm     │      │notifier  │
                   │(default) │              │(lead→    │      │(welcome  │
                   │          │              │ customer)│      │ email)   │
                   └──────────┘              └──────────┘      └──────────┘
```

### Flow 2: Device Provisioning → Telemetry

```
┌──────────┐  1.register  ┌──────────┐
│  Admin   │─────────────►│  device  │
└──────────┘              └────┬─────┘
                               │
                   2.publish: device.device.registered
                               │
                               ▼
                          ┌──────────┐
                          │  alarm   │ (create default rules)
                          └──────────┘

┌──────────┐  3.MQTT       ┌──────────┐
│  Device  │─────────────►│   mqtt   │
└──────────┘               └────┬─────┘
                                │
                                ▼
                          ┌──────────┐
                          │  device  │◄── IngestTelemetryUseCase
                          └────┬─────┘
                               │
              ┌────────────────┼────────────────┐
              ▼                ▼                ▼
      ┌──────────┐     ┌──────────┐     ┌──────────┐
      │ InfluxDB │     │  Kafka   │     │  WS Hub  │
      │  (write) │     │(publish) │     │(broadcast│
      └──────────┘     └────┬─────┘     └──────────┘
                            │
              ┌─────────────┼─────────────┬─────────────┐
              ▼             ▼             ▼             ▼
        ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐
        │  alarm   │ │   AI     │ │  report  │ │logistics │
        │(check)   │ │(analyze) │ │(aggregate│ │(tracking)│
        └──────────┘ └──────────┘ └──────────┘ └──────────┘
```

### Flow 3: Automation Rule Fire → Command

```
┌──────────┐
│Telemetry │──► device.IngestTelemetryUseCase
└──────────┘         │
                     ▼
              ┌──────────────────┐
              │AutomationEngine  │
              │EvaluateAll()     │
              └────────┬─────────┘
                       │ rule fired
                       ▼
              ┌──────────────────┐
              │publish:          │
              │device.rule.fired │
              └────────┬─────────┘
                       │
                       ▼
              ┌──────────────────┐
              │  control module  │──► IssueCommandUseCase
              └────────┬─────────┘
                       │
              ┌────────┼────────┐
              ▼        ▼        ▼
        ┌──────┐ ┌──────┐ ┌──────┐
        │  DB  │ │ Kafk │ │ MQTT │──► Device
        └──────┘ └──────┘ └──────┘
```

### Flow 4: Subscription Lifecycle

```
┌──────────┐  1.subscribe  ┌──────────┐  2.publish  ┌──────────┐
│Customer  │──────────────►│ package  │────────────►│ payment  │
└──────────┘               └────┬─────┘             └────┬─────┘
                                │                        │
                                ▼                        ▼
                          ┌──────────┐             ┌──────────┐
                          │  Redis   │             │Checkout  │
                          │ (cache)  │             │ URL      │
                          └──────────┘             └────┬─────┘
                                                        │
                                             3.payment succeeded
                                                        │
                                                        ▼
                          ┌──────────────────────────────────────┐
                          │ publish: payment.payment.succeeded   │
                          └────────────┬─────────────────────────┘
                                       │
              ┌────────────────────────┼────────────────────────┐
              ▼                        ▼                        ▼
       ┌──────────┐            ┌──────────┐            ┌──────────┐
       │ package  │            │ device   │            │ notifier │
       │(activate)│            │(unlock   │            │(email    │
       │          │            │ quota)   │            │ receipt) │
       └──────────┘            └──────────┘            └──────────┘
```

### Flow 5: Order → Shipment

```
┌──────────┐  1.create   ┌──────────┐
│Customer  │────────────►│   erp    │──► Reserve Stock
└──────────┘  order      └────┬─────┘
                              │ 2.publish: erp.order.created
                              ▼
                        ┌──────────┐
                        │  items   │──► Deduct stock
                        └──────────┘
                              │
                   3.publish: erp.order.confirmed
                              │
                              ▼
                        ┌──────────┐
                        │logistics │──► Create shipment
                        └────┬─────┘
                              │ 4.publish: logistics.shipment.created
                              ▼
                        ┌──────────┐
                        │  device  │──► Attach tracker (if cold chain)
                        └──────────┘
                              │
                   5.publish: logistics.shipment.delivered
                              │
                              ▼
                        ┌──────────┐
                        │   erp    │──► Update order → DELIVERED
                        └────┬─────┘
                              │
                   6.publish: erp.invoice.issued
                              │
                              ▼
                        ┌──────────┐
                        │ payment  │──► Send payment request
                        └──────────┘
```

### Flow 6: Alarm → Notification

```
┌──────────┐  1.telemetry  ┌──────────┐
│  Device  │──────────────►│  device  │
└──────────┘               └────┬─────┘
                                │ 2.publish: device.telemetry.ingested
                                ▼
                          ┌──────────┐
                          │  alarm   │──► Check rules
                          └────┬─────┘
                               │ 3.alert triggered
                               ▼
                    publish: alarm.alert.triggered
                               │
              ┌────────────────┼────────────────┐
              ▼                ▼                ▼
       ┌──────────┐     ┌──────────┐     ┌──────────┐
       │notifier  │     │  device  │     │  WS Hub  │
       │(multi-   │     │(auto-cmd)│     │(realtime)│
       │ channel) │     │          │     │          │
       └────┬─────┘     └──────────┘     └──────────┘
            │
    ┌───────┼───────┬───────┐
    ▼       ▼       ▼       ▼
┌──────┐┌──────┐┌──────┐┌──────┐
│Email ││ SMS  ││ LINE ││ Push │
└──────┘└──────┘└──────┘└──────┘
```

### Flow 7: PDPA DSAR

```
┌──────────┐  1.DSAR     ┌──────────┐
│  User    │────────────►│   pdpa   │
└──────────┘  request    └────┬─────┘
                              │ 2.verify identity (OTP)
                              ▼
                        ┌──────────┐
                        │  auth    │
                        └────┬─────┘
                              │ 3.collect data
              ┌───────────────┼───────────────┬────────────┐
              ▼               ▼               ▼            ▼
       ┌──────────┐   ┌──────────┐    ┌──────────┐  ┌──────────┐
       │customer  │   │  erp     │    │  device  │  │   crm    │
       │(profile) │   │ (orders) │    │(telemetry│  │(tickets) │
       └──────────┘   └──────────┘    └──────────┘  └──────────┘
              │               │               │            │
              └───────────────┼───────────────┴────────────┘
                              ▼
                        ┌──────────┐
                        │ document │──► Generate PDF/JSON
                        └────┬─────┘
                              │ 4.send
                              ▼
                        ┌──────────┐
                        │notifier  │
                        └──────────┘
```

### Flow 8: Report Generation

```
┌──────────┐  1.request   ┌──────────┐
│  User    │─────────────►│  report  │
└──────────┘  report      └────┬─────┘
                               │ 2.query data sources
              ┌────────────────┼────────────────┐
              ▼                ▼                ▼
       ┌──────────┐    ┌──────────┐    ┌──────────┐
       │Postgres  │    │ InfluxDB │    │   ES     │
       │(master)  │    │(telemetry│    │ (logs)   │
       └────┬─────┘    └────┬─────┘    └────┬─────┘
            │               │               │
            └───────────────┼───────────────┘
                            ▼
                    ┌──────────────┐
                    │  Aggregate   │
                    └──────┬───────┘
                           │ 3.optional: AI insight
                           ▼
                    ┌──────────────┐
                    │  LLM (pkg)   │
                    └──────┬───────┘
                           │ 4.export
                           ▼
                    ┌──────────────┐
                    │  Exporter    │──► PDF / Excel / CSV
                    └──────┬───────┘
                           │ 5.upload
                           ▼
                    ┌──────────────┐
                    │     S3       │
                    └──────┬───────┘
                           │ 6.notify
                           ▼
                    ┌──────────────┐
                    │  notifier    │
                    └──────────────┘
```

---

## 7. Event Matrix (Kafka)

### 7.1 ตาราง Publish → Consume

| Topic | Publisher | Consumer(s) | Payload |
| :--- | :--- | :--- | :--- |
| **Customer** | | | |
| `customer.customer.created` | customer | package, crm, notifier | `{customer_id, tenant_id, code, name}` |
| `customer.customer.activated` | customer | package, device | `{customer_id}` |
| `customer.customer.suspended` | customer | device, notifier | `{customer_id, reason}` |
| `customer.customer.churned` | customer | package, device | `{customer_id, reason}` |
| `customer.package.assigned` | customer | package | `{customer_id, package_id}` |
| **Package** | | | |
| `package.subscription.created` | package | payment, customer, device | `{subscription_id, customer_id, package_id}` |
| `package.subscription.upgraded` | package | payment, device | `{subscription_id, from, to}` |
| `package.subscription.cancelled` | package | device, notifier | `{subscription_id, reason}` |
| `package.subscription.expired` | package | device, notifier | `{subscription_id}` |
| `package.quota.exceeded` | package | device, notifier | `{subscription_id, quota_type}` |
| **Device** | | | |
| `device.device.registered` | device | alarm, report | `{device_id, tenant_id, serial, type}` |
| `device.device.online` | device | dashboard | `{device_id, last_seen}` |
| `device.device.offline` | device | alarm, notifier, crm | `{device_id}` |
| `device.telemetry.ingested` | device | alarm, AI, report, logistics | `{device_id, metrics, timestamp}` |
| `device.command.issued` | device | control | `{command_id, device_id, type}` |
| `device.command.acked` | device | control, dashboard | `{command_id, device_id}` |
| `device.rule.fired` | device | control, notifier | `{rule_id, actions}` |
| **Alarm** | | | |
| `alarm.alert.triggered` | alarm | notifier, device, WS | `{alert_id, severity, device_id}` |
| **ERP** | | | |
| `erp.order.created` | erp | logistics, items, notifier | `{order_id, customer_id, total}` |
| `erp.order.confirmed` | erp | logistics | `{order_id}` |
| `erp.order.shipped` | erp | logistics, customer | `{order_id}` |
| `erp.invoice.issued` | erp | payment, notifier | `{invoice_id, total}` |
| `erp.invoice.paid` | erp | customer, report | `{invoice_id}` |
| `erp.stock.low` | erp | purchaseorder | `{product_id, quantity}` |
| **CRM** | | | |
| `crm.lead.created` | crm | notifier | `{lead_id}` |
| `crm.lead.converted` | crm | customer | `{lead_id, customer_id}` |
| `crm.ticket.opened` | crm | notifier, wos | `{ticket_id, priority}` |
| `crm.ticket.resolved` | crm | notifier | `{ticket_id}` |
| **Logistics** | | | |
| `logistics.shipment.created` | logistics | device, erp | `{shipment_id}` |
| `logistics.shipment.delivered` | logistics | erp, notifier | `{shipment_id}` |
| `logistics.temperature.breached` | logistics | alarm, notifier | `{shipment_id, temp}` |
| **Payment** | | | |
| `payment.payment.succeeded` | payment | package, erp, notifier | `{payment_id, reference_id}` |
| `payment.payment.failed` | payment | package, notifier | `{payment_id, reason}` |
| `payment.payment.refunded` | payment | erp, notifier | `{payment_id}` |
| **Alarm** | | | |
| `alarm.alert.resolved` | alarm | notifier, dashboard | `{alert_id}` |

### 7.2 Naming Convention

```
<context>.<aggregate>.<event>

Examples:
  customer.customer.created
  device.telemetry.ingested
  erp.order.shipped
  alarm.alert.triggered

Rules:
- context = bounded context
- aggregate = entity name (lowercase)
- event = past tense verb
- separator = dot (.)
```

### 7.3 Consumer Groups

| Consumer Group | Topic ที่ consume | Module |
| :--- | :--- | :--- |
| `package-group` | customer.customer.created, payment.* | package |
| `device-group` | customer.customer.activated, package.quota.exceeded | device |
| `alarm-group` | device.telemetry.ingested, logistics.temperature.breached | alarm |
| `notifier-group` | *.* (event ที่ต้องแจ้งเตือน) | notifier |
| `report-group` | device.telemetry.ingested, erp.invoice.paid | report |
| `logistics-group` | erp.order.confirmed, device.telemetry.ingested | logistics |
| `ai-group` | device.telemetry.ingested | (AI service) |

---

## 8. Sync API Matrix

### 8.1 Internal APIs (module → module)

| จาก | ไป | Endpoint | ใช้เมื่อ |
| :--- | :--- | :--- | :--- |
| `device` | `customer` | `GET /internal/customers/:id/status` | ตรวจ active |
| `device` | `package` | `POST /internal/quota/check-device` | ตรวจโควตา |
| `package` | `customer` | `GET /internal/customers/:id` | ดึงข้อมูล |
| `package` | `payment` | `POST /internal/payments` | สร้าง payment |
| `erp` | `customer` | `GET /internal/customers/:id` | ดึงข้อมูล |
| `erp` | `items` | `GET /internal/items/:id` | ดึงสินค้า |
| `erp` | `payment` | `POST /internal/payments` | สร้าง invoice |
| `logistics` | `erp` | `GET /internal/orders/:id` | ดึง order |
| `logistics` | `device` | `GET /internal/devices/:id` | ดึง tracker |
| `crm` | `customer` | `POST /internal/customers` | Convert lead |
| `report` | `device` | `GET /internal/devices/:id` | Enrich |
| `dashboard` | `device` | `GET /internal/devices/:id/latest` | Realtime |
| `alarm` | `device` | `POST /internal/devices/:id/command` | Auto-action |
| `pdpa` | ทุกโมดูล | `GET /internal/data/:user_id` | DSAR |

### 8.2 Public APIs (client → platform)

| Base Path | Module | Auth |
| :--- | :--- | :--- |
| `/api/v1/auth/*` | auth | Public (login) |
| `/api/v1/customers/*` | customer | JWT + tenant + RBAC |
| `/api/v1/packages/*` | package | JWT |
| `/api/v1/subscriptions/*` | package | JWT |
| `/api/v1/sites/*` | device | JWT + tenant |
| `/api/v1/devices/*` | device | JWT + tenant |
| `/api/v1/automation-rules/*` | device | JWT + tenant |
| `/api/v1/erp/*` | erp | JWT + tenant + RBAC |
| `/api/v1/crm/*` | crm | JWT + tenant + RBAC |
| `/api/v1/logistics/*` | logistics | JWT + tenant |
| `/api/v1/reports/*` | report | JWT + tenant |
| `/api/v1/dashboards/*` | dashboard | JWT + tenant |
| `/api/v1/ws` | websocket | JWT (query param) |

### 8.3 API Versioning

```
/api/v1/...  → current stable
/api/v2/...  → when breaking change needed
```

---

## 9. Shared Data

### 9.1 Postgres — Schema per Module

| Module | Tables | Prefix |
| :--- | :--- | :--- |
| customer | customer_customers, customer_contacts, customer_sites | `customer_` |
| package | package_packages, package_subscriptions | `package_` |
| device | device_sites, device_zones, device_devices, device_commands, device_automation_rules | `device_` |
| erp | erp_orders, erp_order_items, erp_invoices, erp_warehouses, erp_stock_items | `erp_` |
| crm | crm_leads, crm_opportunities, crm_tickets, crm_activities, crm_campaigns | `crm_` |
| logistics | logistics_shipments, logistics_routes, logistics_vehicles, logistics_drivers | `logistics_` |
| report | report_definitions, report_runs, report_dashboards | `report_` |
| auth | auth_sessions, auth_refresh_tokens, auth_mfa | `auth_` |
| users | users_users, users_roles, users_permissions | `users_` |
| payment | payment_payments, payment_transactions, payment_refunds | `payment_` |
| notifier | notifier_templates, notifier_logs, notifier_channels | `notifier_` |
| alarm | alarm_rules, alarm_events, alarm_escalations | `alarm_` |
| pdpa | pdpa_consents, pdpa_dsar, pdpa_audit | `pdpa_` |
| auditlog | audit_trails, audit_events | `audit_` |
| settings | settings_tenant_settings, settings_user_preferences | `settings_` |
| items | items_products, items_categories | `items_` |
| purchaseorder | po_purchase_orders, po_po_items | `po_` |
| quotation | quotation_quotes, quotation_quote_items | `quotation_` |
| wos | wos_work_orders, wos_tasks | `wos_` |
| document | document_documents, document_versions | `document_` |
| email | email_templates, email_queue | `email_` |
| batch | batch_jobs, batch_runs | `batch_` |
| job | job_jobs, job_executions | `job_` |
| i18n | i18n_translations, i18n_locales | `i18n_` |
| apimanager | apimanager_keys, apimanager_usage | `apimanager_` |
| flowengine | flowengine_flows, flowengine_runs | `flowengine_` |
| control | control_templates, control_batches | `control_` |
| fullschedule | fullschedule_schedules | `fullschedule_` |
| queue | queue_jobs, queue_dlq | `queue_` |

### 9.2 Redis — Key Namespaces

| Pattern | Module | TTL | Purpose |
| :--- | :--- | :-: | :--- |
| `customer:customer:{id}` | customer | 15m | Cache |
| `package:package:{id}` | package | 30m | Cache |
| `package:customer:{customer_id}` | package | 15m | Active sub lookup |
| `device:device:{id}` | device | 5m | Cache |
| `device:device:serial:{serial}` | device | 5m | Lookup by serial |
| `device:shadow:{id}` | device | 30s | Shadow cache |
| `device:rate:{device_id}:tel` | device | 1s | Rate limit |
| `auth:session:{session_id}` | auth | 24h | Session |
| `auth:refresh:{token}` | auth | 7d | Refresh token |
| `auth:jwt:blacklist:{jti}` | auth | = TTL | JWT blacklist |
| `apimanager:key:{key}` | apimanager | 5m | API key cache |
| `apimanager:usage:{key}:{date}` | apimanager | 24h | Usage counter |
| `notifier:rate:{channel}:{user}` | notifier | 1h | Rate limit |
| `alarm:cooldown:{rule_id}` | alarm | varies | Rule cooldown |
| `queue:dlq:{queue_name}` | queue | ∞ | Dead letter |

### 9.3 InfluxDB — Buckets & Measurements

| Bucket | Measurement | Writer | Reader | Retention |
| :--- | :--- | :--- | :--- | :-: |
| `telemetry` | `telemetry` | device | report, dashboard, AI | 90d |
| `telemetry_1h` | `telemetry` | downsampling task | report | 1y |
| `telemetry_1d` | `telemetry` | downsampling task | report | 3y |
| `energy` | `power` | device | report | 2y |
| `logistics` | `temperature` | logistics | report, alarm | 1y |

**Tags:** `device_id`, `site_id`, `tenant_id`, `metric`, `unit`, `quality`

### 9.4 Elasticsearch — Indices

| Index | Writer | Reader | Purpose |
| :--- | :--- | :--- | :--- |
| `customer_customers` | customer | crm, report | Search customers |
| `device_devices` | device | report, dashboard | Search devices |
| `device_telemetry_*` | device | report | Search telemetry (daily index) |
| `audit_trails_*` | auditlog | report, compliance | Search audit (monthly) |
| `erp_orders` | erp | report | Search orders |
| `crm_leads` | crm | report | Search leads |
| `pdpa_events` | pdpa | compliance | Search consent/DSAR |

**ILM Policy:** hot (7d) → warm (30d) → cold (90d) → delete (365d)

### 9.5 S3 Buckets

| Bucket | Module | Content |
| :--- | :--- | :--- |
| `icmongolang-documents` | document | PDFs, contracts |
| `icmongolang-reports` | report | Generated reports |
| `icmongolang-firmware` | device | Firmware binaries |
| `icmongolang-pdpa` | pdpa | DSAR exports |
| `icmongolang-backups` | infra | DB backups |

---

## 10. Cross-Cutting Concerns

### 10.1 Middleware Stack

```
Request ──► RequestID ──► Logging ──► CORS ──► Security ──► RateLimit ──► Auth(JWT) ──► Tenant ──► RBAC ──► Handler
                 │           │        │        │            │             │            │         │
                 │           │        │        │            │             │            │         │
              X-Req-ID    trace_id   headers  headers    redis count   verify JWT   X-Tenant  role check
```

### 10.2 Observability Stack

| Signal | เครื่องมือ | Module ที่เกี่ยวข้อง |
| :--- | :--- | :--- |
| **Metrics** | Prometheus + Grafana | ทุกโมดูล |
| **Logs** | Loki / ELK | ทุกโมดูล |
| **Traces** | Tempo / Jaeger | device, erp, logistics |
| **Profiling** | Pyroscope / pprof | – |
| **APM** | OpenTelemetry Collector | รวมทุก signal |
| **Alerts** | Alertmanager | – |

### 10.3 Security Layers

| Layer | Module | กลไก |
| :--- | :--- | :--- |
| **Authentication** | auth | JWT + refresh + MFA + SSO |
| **Authorization** | users | RBAC (role) + ABAC (attribute) |
| **Tenant Isolation** | ทุกโมดูล | `tenant_id` + RLS |
| **API Security** | apimanager | API key + rate limit + quota |
| **Data Encryption** | pdpa | AES-256 at rest + TLS in transit |
| **Audit** | auditlog | Every state change |
| **Input Validation** | ทุกโมดูล | Handler + Domain (defense in depth) |
| **SQL Injection** | ทุกโมดูล | Parameterized query (GORM/SQLAlchemy) |

### 10.4 Event-Driven Patterns

| Pattern | Module | ตัวอย่าง |
| :--- | :--- | :--- |
| **Pub/Sub** | ทุกโมดูล | Kafka events |
| **Event Sourcing** | (เฉพาะ) | auditlog, pdpa |
| **CQRS** | customer, device, erp | Read model + Write model |
| **Saga** | erp, logistics | Order → Shipment → Invoice → Payment |
| **Outbox** | ทุกโมดูล | Transactional outbox pattern |
| **DLQ** | kafka, queue | Failed message handling |
| **Idempotency** | payment, erp | Idempotency key |

### 10.5 Resilience Patterns

| Pattern | ใช้ที่ | กลไก |
| :--- | :--- | :--- |
| **Circuit Breaker** | device → customer | gobreaker |
| **Retry with Backoff** | Kafka consumer | exponential |
| **Timeout** | External calls | context.WithTimeout |
| **Bulkhead** | AI service | worker pool |
| **Rate Limiting** | API Gateway | Redis token bucket |
| **Fallback** | AI | cached response |

---

## 11. Deployment Topology

### 11.1 Services ที่ต้อง deploy

```
┌────────────────────────────────────────────────────────────┐
│                    Production Cluster                       │
├────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  api-gw-1..N │  │ worker-1..N  │  │scheduler-1..2│     │
│  │  (Gin/FastAPI)│ │  (Kafka)     │  │  (cron)      │     │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘     │
│         │                 │                  │             │
│  ┌──────┴─────────────────┴──────────────────┴──────┐     │
│  │              Internal Load Balancer              │     │
│  └──────────────────────┬───────────────────────────┘     │
│                         │                                  │
│         ┌───────────────┼───────────────┐                 │
│         ▼               ▼               ▼                 │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐          │
│  │ Postgres   │  │  Redis     │  │  Kafka     │          │
│  │ (Primary   │  │  Cluster   │  │  Cluster   │          │
│  │  + Replica)│  │            │  │            │          │
│  └────────────┘  └────────────┘  └────────────┘          │
│                                                             │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐          │
│  │ InfluxDB   │  │Elasticsearch│ │   MinIO    │          │
│  │ (Cluster)  │  │  Cluster    │ │  (S3)      │          │
│  └────────────┘  └────────────┘  └────────────┘          │
│                                                             │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐          │
│  │ MQTT       │  │Prometheus  │  │  Grafana   │          │
│  │ (EMQX)     │  │ + Alert    │  │  + Loki    │          │
│  └────────────┘  └────────────┘  └────────────┘          │
│                                                             │
└────────────────────────────────────────────────────────────┘
```

### 11.2 Process Types

| Process | บทบาท | Module ที่รัน |
| :--- | :--- | :--- |
| `api` | HTTP + WebSocket | ทุกโมดูล (routes) |
| `worker` | Kafka consumer | device, alarm, notifier, AI |
| `scheduler` | Cron jobs | device (heartbeat), package (expire), report |
| `mqtt-gateway` | MQTT bridge | device |
| `migrate` | DB migration | – |

### 11.3 Container Images

```
icmongolang/api:latest          → HTTP + WS server
icmongolang/worker-telemetry    → Ingest consumer
icmongolang/worker-alarm        → Alarm consumer
icmongolang/worker-notifier     → Notifier consumer
icmongolang/worker-ai           → AI consumer
icmongolang/scheduler           → Cron jobs
icmongolang/mqtt-gateway        → MQTT bridge
icmongolang/ollama-exporter     → Metrics exporter
```

---

## 12. Roadmap & Production Order

### 12.1 Phase Timeline

```
        Q1        Q2        Q3        Q4        Q5        Q6        Q7        Q8
       ──┴─────────┴─────────┴─────────┴─────────┴─────────┴─────────┴─────────┴──

Ph0    ████                                                                         Foundation
Ph1              ████████████                                                     MVP
Ph2                        ████████████████                                       Growth
Ph3                                      ████████████████                         Enterprise
Ph4                                                ████████████████               Scale
```

### 12.2 Module Production Order

| Phase | Month | Modules | Goal |
| :--- | :-: | :--- | :--- |
| **Ph0.1** | 1-2 | `customer` ✅, `package` ✅ | รากฐาน |
| **Ph0.2** | 3-4 | `device` ✅, `mqtt`, `influxdb` | IoT core |
| **Ph0.3** | 5 | `alarm`, `notifier`, `websocket` | Realtime |
| **Ph0.4** | 6 | `payment`, `auth` (upgrade), `users` (upgrade) | Billing + security |
| **Ph1.1** | 7-8 | `erp`, `items` | Business core |
| **Ph1.2** | 9-10 | `crm`, `quotation` | Sales |
| **Ph1.3** | 11-12 | `report` (upgrade), `dashboard` (upgrade) | Analytics |
| **Ph2.1** | 13-15 | `logistics`, `wos` | Supply chain |
| **Ph2.2** | 16-18 | `settings`, `auditlog`, `pdpa` (upgrade) | Governance |
| **Ph3** | 19-24 | ที่เหลือ + upgrade ที่เหลือ | ครบระบบ |

### 12.3 Milestone

| Milestone | ปลาย Phase | Deliverable |
| :--- | :--- | :--- |
| **M1** | Ph0 | ระบบพื้นฐานทำงานได้ end-to-end |
| **M2** | Ph1 | MVP พร้อม pilot 10 ราย |
| **M3** | Ph2 | เปิดขายจริง, revenue |
| **M4** | Ph3 | Enterprise-ready, 1000 customers |
| **M5** | Ph4 | Scale ได้ 100k devices |

---

## 📊 สรุปสุดท้าย

### ✅ สิ่งที่ได้

| หมวด | จำนวน |
| :--- | :-: |
| Modules ทั้งหมด | 40 |
| Modules Done (full doc) | 3 (customer, package, device) |
| Modules มี compact spec | 5 (erp, crm, logistics, report, auth) |
| Modules มี audit+target | 32 |
| Bounded Contexts | 8 |
| Kafka Topics | ~50 |
| Sync APIs (internal) | ~15 |
| Tables (Postgres) | ~80 |
| Redis Key Patterns | ~15 |
| InfluxDB Buckets | 5 |
| ES Indices | 7 |
| S3 Buckets | 5 |
| Container Images | 8 |

### 🎯 Key Insights

1. **`device` เป็นหัวใจ** — 5 cross-cutting (MQTT, Kafka, InfluxDB, WS, Cache)
2. **Event-Driven เป็นหลัก** — 50+ Kafka topics, fan-out pattern
3. **Sync API เฉพาะที่จำเป็น** — ส่วนใหญ่ใช้ async
4. **Multi-tenant ทุกที่** — `tenant_id` + RLS
5. **CQRS ในโมดูลสำคัญ** — customer, device, erp
6. **Layer dependency เข้มงวด** — Domain ไม่รู้จัก Infra

### 🚀 ขั้นตอนถัดไป

1. **Generate code** — ป้อน Module 01-03 specs ให้ AI → ได้ code
2. **ทำ Module 04-08** — compact spec พร้อมแล้ว
3. **Audit โมดูลเดิม** — ใช้ template เดียวกัน
4. **Deploy Phase 0** — customer + package + device
5. **Monitor + iterate** — 30 นาทีแรกหลัง deploy

---

> **เอกสารนี้สรุปทุก modules + การเชื่อมโยงครบ 12 มุมมอง** — ใช้เป็นแผนที่นำทางสำหรับทีม Development, DevOps, และ AI Generation

**พิมพ์คำสั่งถัดไปได้:**
- `"ทำ Module 04"` → ERP full deep dive
- `"ทำ Module 05"` → CRM
- `"ทำ Module 06"` → Logistics
- `"ทำ Module 07"` → Report upgrade (breaking change)
- `"ทำ Module 08"` → Auth upgrade
- `"Generate code ทั้งหมด"` → สร้าง prompt batch สำหรับ AI