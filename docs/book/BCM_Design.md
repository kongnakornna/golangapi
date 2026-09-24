# 🏗️ ออกแบบระบบแพลตฟอร์มให้บริการ IoT (Smart Farm / Smart Building)
### ตามโครงสร้าง `icmongolang` — Clean Architecture + DDD

ออกแบบ ระบบ  ให้บริการ IoT  เช่น ระบบ Smart farm , Smart building   
1.ระบบลูกค้า 
3.Package และ บริการ  
4.ระบบ ERP IoT Sulution 
5.ระบบ CRM  IoT Sulution
6.ระบบ จัดการอุกรณ์ IoT ,AI ,Automation
7.ระบบ บริหารจัดการ Logistic IoT 
8.ระบบรายงาน 

ตามโครงสร้าง แต่ละ Module
1. [ภาพรวมระบบ](#1-ภาพรวมระบบ)
2. [โครงสร้าง Module](#2-โครงสร้าง-module)
3. [Bounded Contexts และ Context Map](#3-bounded-contexts-และ-context-map)
4. [Ubiquitous Language](#4-ubiquitous-language)
5. [Domain Layer](#5-domain-layer)
   - 5.1 Entities (Aggregates)
   - 5.2 Value Objects
   - 5.3 Repository Interfaces
   - 5.4 Domain Services
   - 5.5 Domain Errors
   - 5.6 Invariants
6. [Application Layer](#6-application-layer)
   - 6.1 Use Cases
   - 6.2 DTOs
7. [Infrastructure Layer](#7-infrastructure-layer)
   - 7.1 Repository Implementations
   - 7.2 Kafka Consumers
   - 7.3 WebSocket Broadcaster
   - 7.4 JWT, Bcrypt, Rate Limit
8. [Interface Layer](#8-interface-layer)
   - 8.1 HTTP Handlers
   - 8.2 Routes
   - 8.3 Middleware
9. [Database Migrations](#9-database-migrations)
10. [System Flow](#10-system-flow)
11. [Workflow Diagram](#11-workflow-diagram)
12. [การติดตั้งและใช้งาน](#12-การติดตั้งและใช้งาน)
13. [Business Model](#13-business-model)
14. [ภาคผนวก](#14-ภาคผนวก)
15. [Prompt สำหรับการขยายระบบในอนาคต](#15-prompt-สำหรับการขยายระบบในอนาคต)
16. [การตั้งชื่อตาราง Database (Prefix)](#16-การตั้งชื่อตาราง-database-prefix)
17. [DDD Validation Checklist](#17-ddd-validation-checklist)

ออกแบบ ตามโครงสร้าง
---

## 1. ภาพรวมสถาปัตยกรรม

ระบบแพลตฟอร์ม IoT ให้บริการแบบ **multi-tenant SaaS** สำหรับ Smart Farm และ Smart Building โดยใช้ 7 โมดูลธุรกิจหลัก + cross-cutting services ต่อยอดจากโครงสร้างเดิมของ `icmongolang` (มีอยู่แล้ว: `customer`, `items`, `payment`, `report`, `iot`, `notifier`, `auditlog`)

```
┌──────────────────────────────────────────────────────────────────┐
│                    IoT Service Platform (SaaS)                    │
│                                                                    │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐          │
│  │ Customer │  │ Package  │  │   ERP    │  │   CRM    │          │
│  │ ระบบลูกค้า│  │& Service │  │          │  │          │          │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘          │
│       │             │             │             │                 │
│  ┌────┴─────────────┴─────────────┴─────────────┴─────┐          │
│  │            Shared Domain Services                   │          │
│  │  Tenant · Billing · Notification · Audit · IAM      │          │
│  └────┬─────────────┬─────────────┬─────────────┬─────┘          │
│       │             │             │             │                 │
│  ┌────┴─────┐  ┌────┴─────┐  ┌────┴─────┐  ┌────┴─────┐          │
│  │ IoT      │  │ IoT      │  │ Report   │  │ Realtime │          │
│  │ Device   │  │ Logistics│  │ & Analytic│ │ (WS/MQTT)│          │
│  │ Mgmt     │  │          │  │          │  │          │          │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘          │
└──────────────────────────────────────────────────────────────────┘
         ▲                    ▲                     ▲
         │                    │                     │
   MQTT / Kafka          REST / gRPC          WebSocket
         │                    │                     │
   [IoT Devices]      [Mobile/Web App]      [Dashboard]
```

**Multi-tenant strategy:** ทุกตารางมี `tenant_id UUID NOT NULL` + Row-Level Security ที่ระดับ Repository

---

## 2. โครงสร้างโมดูลทั้งหมด (7 โมดูลหลัก)

```
internal/modules/
├── customer/          # 1. ระบบลูกค้า (Customer 360)
├── package/           # 2. Package & Service Catalog + Subscription
├── erp/               # 3. ERP (Inventory, Procurement, Finance)
├── crm/               # 4. CRM (Lead, Opportunity, Ticket)
├── iotdevice/         # 5. IoT Device Management
├── iotlogistics/      # 6. IoT Logistics (Asset tracking)
└── report/            # 7. Reporting & Analytics (extend ของเดิม)
```

---

## 3. โมดูลที่ 1 — ระบบลูกค้า (`customer`)

### 3.1 โครงสร้าง
```
internal/modules/customer/
├── domain/
│   ├── entity/
│   │   ├── customer.go            # Aggregate Root
│   │   ├── contact.go             # ลูกค้าติดต่อ
│   │   ├── site.go                # สถานที่ติดตั้ง (farm/building)
│   │   └── contract.go            # สัญญาบริการ
│   ├── value_object/
│   │   ├── customer_type.go       # INDIVIDUAL | CORPORATE | GOVERNMENT
│   │   ├── customer_status.go     # LEAD | ACTIVE | SUSPENDED | CHURNED
│   │   └── tax_id.go              # เลขผู้เสียภาษี + validation
│   ├── repository/
│   │   ├── customer_repository.go
│   │   ├── site_repository.go
│   │   └── contract_repository.go
│   ├── service/
│   │   └── customer_domain_service.go  # ตรวจเครดิต, KYC
│   └── errors/errors.go
├── application/
│   ├── create_customer.go
│   ├── update_customer.go
│   ├── onboard_customer.go        # multi-step: KYC → contract → site
│   ├── suspend_customer.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   │   ├── customer_repo_impl.go
│   │   ├── site_repo_impl.go
│   │   └── models.go
│   ├── persistence/redis/customer_cache.go
│   ├── search/elasticsearch/customer_indexer.go
│   └── messaging/kafka_producer.go
├── interfaces/
│   ├── http/
│   │   ├── customer_handler.go
│   │   ├── site_handler.go
│   │   └── routes.go
│   └── middleware/tenant.go       # inject tenant_id
└── module.go
```

### 3.2 Database
```sql
-- migrations/20260101_customer_init.sql
CREATE TABLE customer_customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code VARCHAR(30) UNIQUE NOT NULL,       -- CUS-2026-0001
    type VARCHAR(20) NOT NULL,              -- INDIVIDUAL/CORPORATE/GOVERNMENT
    name VARCHAR(255) NOT NULL,
    tax_id VARCHAR(20),
    email VARCHAR(255),
    phone VARCHAR(30),
    address JSONB,
    status VARCHAR(20) DEFAULT 'LEAD',
    credit_limit NUMERIC(15,2) DEFAULT 0,
    created_by UUID,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_customer_customers_tenant ON customer_customers(tenant_id);
CREATE INDEX idx_customer_customers_status ON customer_customers(tenant_id, status);

CREATE TABLE customer_sites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customer_customers(id) ON DELETE CASCADE,
    site_type VARCHAR(30),          -- FARM | BUILDING | FACTORY | WAREHOUSE
    name VARCHAR(255) NOT NULL,
    geo_lat NUMERIC(10,7),
    geo_lng NUMERIC(10,7),
    address JSONB,
    area_size NUMERIC(12,2),
    metadata JSONB,                 -- farm: crop type; building: floors
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE customer_contracts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customer_customers(id),
    package_id UUID NOT NULL,
    contract_no VARCHAR(50) UNIQUE NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    status VARCHAR(20) DEFAULT 'DRAFT',
    signed_at TIMESTAMP,
    document_url TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 3.3 Kafka Topics
```
customer.created         # ใหม่ → CRM + ERP
customer.updated         # sync ทุกโมดูล
customer.status.changed  # → package (suspend service), notifier
customer.churned         # → iotlogistics (คืนอุปกรณ์)
```

### 3.4 HTTP Routes
```go
// interfaces/http/routes.go
g := r.Group("/customers")
g.Use(auth, tenant.Inject())
g.POST   ("/",              h.Create)        // POST /api/v1/customers
g.GET    ("/",              h.List)          // list + filter
g.GET    ("/:id",           h.Get)
g.PUT    ("/:id",           h.Update)
g.POST   ("/:id/onboard",   h.Onboard)       // start KYC flow
g.POST   ("/:id/suspend",   h.Suspend)
g.GET    ("/:id/sites",     h.ListSites)
g.POST   ("/:id/sites",     h.CreateSite)
g.GET    ("/:id/contracts", h.ListContracts)
```

---

## 4. โมดูลที่ 2 — Package & Service (`package`)

> ⚠️ ชื่อ Go package ต้องไม่ชน keyword → ใช้ `packagecatalog` เป็น Go package name แต่โฟลเดอร์ชื่อ `package`

### 4.1 โครงสร้าง
```
internal/modules/package/
├── domain/
│   ├── entity/
│   │   ├── service_package.go       # แพ็กเกจ (Smart Farm Pro, Building Basic)
│   │   ├── service_item.go          # บริการย่อยในแพ็กเกจ
│   │   ├── subscription.go          # การสมัครใช้ของลูกค้า
│   │   ├── usage_record.go          # บันทึกการใช้งาน (quota)
│   │   └── quota.go                 # VO ของโควต้า
│   ├── value_object/
│   │   ├── billing_cycle.go         # MONTHLY | YEARLY | USAGE_BASED
│   │   ├── package_category.go      # SMART_FARM | SMART_BUILDING | MIXED
│   │   └── subscription_status.go
│   ├── repository/
│   │   ├── package_repository.go
│   │   ├── subscription_repository.go
│   │   └── usage_repository.go
│   └── service/entitlement_service.go  # ตรวจสิทธิ์ตาม package
├── application/
│   ├── create_package.go
│   ├── publish_package.go
│   ├── subscribe_package.go         # เรียก payment + iotdevice
│   ├── upgrade_subscription.go
│   ├── cancel_subscription.go
│   ├── record_usage.go              # รับ event จาก MQTT/Kafka
│   └── check_entitlement.go         # เรียกจาก middleware
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/redis/quota_cache.go   # counter โควต้าต่อ tenant
│   ├── messaging/kafka/
│   │   ├── usage_producer.go
│   │   └── consumers/usage_consumer.go
│   └── scheduler/quota_reset_job.go
├── interfaces/http/
│   ├── package_handler.go
│   ├── subscription_handler.go
│   └── routes.go
└── module.go
```

### 4.2 Database
```sql
-- migrations/20260102_package_init.sql
CREATE TABLE package_packages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID,                 -- NULL = platform-level package
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(30) NOT NULL,  -- SMART_FARM | SMART_BUILDING
    description TEXT,
    billing_cycle VARCHAR(20) NOT NULL,
    base_price NUMERIC(15,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'THB',
    status VARCHAR(20) DEFAULT 'DRAFT',  -- DRAFT|PUBLISHED|DEPRECATED
    metadata JSONB,                 -- features flag
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

CREATE TABLE package_items (           -- บริการย่อย
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    package_id UUID NOT NULL REFERENCES package_packages(id) ON DELETE CASCADE,
    service_code VARCHAR(50) NOT NULL,   -- DEVICE_QUOTA, DATA_STORAGE, API_CALL
    unit VARCHAR(20),                    -- device, GB, request
    quota NUMERIC(15,2),
    overage_price NUMERIC(15,4),
    UNIQUE(package_id, service_code)
);

CREATE TABLE package_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL,
    package_id UUID NOT NULL,
    contract_id UUID,
    status VARCHAR(20) DEFAULT 'PENDING',  -- PENDING|ACTIVE|SUSPENDED|CANCELLED
    start_date DATE NOT NULL,
    end_date DATE,
    next_billing_date DATE,
    auto_renew BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_pkg_subs_customer ON package_subscriptions(customer_id, status);

CREATE TABLE package_usage_records (
    id BIGSERIAL PRIMARY KEY,
    subscription_id UUID NOT NULL REFERENCES package_subscriptions(id),
    service_code VARCHAR(50) NOT NULL,
    quantity NUMERIC(15,4) NOT NULL,
    recorded_at TIMESTAMP DEFAULT NOW(),
    period_start DATE NOT NULL,     -- เพื่อ group by cycle
    idempotency_key VARCHAR(100) UNIQUE
);
CREATE INDEX idx_pkg_usage_sub_period ON package_usage_records(subscription_id, period_start);
```

### 4.3 Kafka Topics
```
package.published           # → catalog cache invalidate
package.subscription.created# → payment (issue invoice) + iotdevice (provision)
package.subscription.updated# → iotdevice (update quota) + notifier
package.usage.recorded      # ← iotdevice, mqtt gateway
package.quota.exceeded      # → notifier + CRM (upsell ticket)
```

### 4.4 HTTP Routes
```
GET    /api/v1/packages
POST   /api/v1/packages                    (admin)
POST   /api/v1/packages/:id/publish
GET    /api/v1/subscriptions
POST   /api/v1/subscriptions               {customer_id, package_id}
POST   /api/v1/subscriptions/:id/upgrade
POST   /api/v1/subscriptions/:id/cancel
GET    /api/v1/subscriptions/:id/usage
POST   /api/v1/subscriptions/:id/check-entitlement
```

---

## 5. โมดูลที่ 3 — ERP (`erp`)

### 5.1 โครงสร้าง
```
internal/modules/erp/
├── domain/
│   ├── entity/
│   │   ├── product.go              # สินค้า/อุปกรณ์ขาย
│   │   ├── inventory_item.go       # stock
│   │   ├── purchase_order.go       # PO
│   │   ├── sales_order.go          # SO
│   │   ├── invoice.go              # ใบแจ้งหนี้
│   │   ├── vendor.go
│   │   └── warehouse.go
│   ├── value_object/
│   │   ├── money.go                # currency + amount
│   │   ├── uom.go                  # unit of measure
│   │   └── doc_status.go           # DRAFT|APPROVED|POSTED|VOID
│   ├── repository/
│   ├── service/
│   │   ├── inventory_service.go    # ตรวจ stock, จอง
│   │   └── pricing_service.go
│   └── errors/
├── application/
│   ├── create_po.go
│   ├── approve_po.go
│   ├── receive_goods.go            # GRN
│   ├── issue_invoice.go
│   ├── post_payment.go
│   ├── stock_adjust.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── services/pdf/invoice_pdf.go
│   ├── messaging/kafka/accounting_producer.go
│   └── scheduler/low_stock_alert_job.go
├── interfaces/http/
│   ├── product_handler.go
│   ├── inventory_handler.go
│   ├── po_handler.go
│   ├── invoice_handler.go
│   └── routes.go
└── module.go
```

### 5.2 Database (ย่อ)
```sql
-- migrations/20260103_erp_init.sql
CREATE TABLE erp_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(50),          -- SENSOR | GATEWAY | ACTUATOR | ACCESSORY
    uom VARCHAR(20) DEFAULT 'pcs',
    cost_price NUMERIC(15,2),
    sell_price NUMERIC(15,2),
    is_iot_device BOOLEAN DEFAULT FALSE,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE erp_warehouses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(30) UNIQUE NOT NULL,
    name VARCHAR(255),
    address JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE erp_inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES erp_products(id),
    warehouse_id UUID NOT NULL REFERENCES erp_warehouses(id),
    qty_on_hand NUMERIC(15,2) DEFAULT 0,
    qty_reserved NUMERIC(15,2) DEFAULT 0,
    qty_available NUMERIC(15,2) GENERATED ALWAYS AS (qty_on_hand - qty_reserved) STORED,
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(product_id, warehouse_id)
);

CREATE TABLE erp_purchase_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    po_no VARCHAR(30) UNIQUE NOT NULL,
    vendor_id UUID,
    status VARCHAR(20) DEFAULT 'DRAFT',
    total_amount NUMERIC(15,2),
    expected_at DATE,
    created_by UUID,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE erp_po_lines (...);
CREATE TABLE erp_invoices (...);
CREATE TABLE erp_payments (...);
CREATE TABLE erp_stock_movements (        -- audit trail
    id BIGSERIAL PRIMARY KEY,
    product_id UUID, warehouse_id UUID,
    movement_type VARCHAR(20),            -- IN|OUT|ADJUST|TRANSFER
    qty NUMERIC(15,2),
    ref_type VARCHAR(30), ref_id UUID,    -- PO, SO, ...
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 5.3 Kafka Topics
```
erp.po.created / erp.po.approved / erp.po.received
erp.inventory.low_stock          # → notifier
erp.invoice.issued               # → payment gateway
erp.payment.posted               # → accounting, report
erp.product.created              # → iotdevice (device model registry)
```

---

## 6. โมดูลที่ 4 — CRM (`crm`)

### 6.1 โครงสร้าง
```
internal/modules/crm/
├── domain/
│   ├── entity/
│   │   ├── lead.go                 # ผู้สนใจ
│   │   ├── opportunity.go          # โอกาสขาย
│   │   ├── pipeline.go             # sales pipeline stage
│   │   ├── activity.go             # call/email/meeting
│   │   ├── ticket.go               # support ticket
│   │   └── campaign.go
│   ├── value_object/
│   │   ├── lead_source.go          # WEB|REFERRAL|EVENT|ADS
│   │   ├── opportunity_stage.go    # NEW|QUALIFIED|PROPOSAL|WON|LOST
│   │   └── ticket_priority.go
│   ├── repository/
│   ├── service/
│   │   ├── lead_scoring_service.go
│   │   └── sla_service.go
│   └── errors/
├── application/
│   ├── create_lead.go
│   ├── convert_lead.go             # → customer + opportunity
│   ├── move_opportunity_stage.go
│   ├── create_ticket.go
│   ├── assign_ticket.go
│   ├── resolve_ticket.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/redis/lead_assign_cache.go
│   ├── services/llm/ticket_classifier.go   # AI จัดหมวดหมู่ ticket
│   └── messaging/kafka/
├── interfaces/http/
│   ├── lead_handler.go
│   ├── opportunity_handler.go
│   ├── ticket_handler.go
│   └── routes.go
└── module.go
```

### 6.2 Database (ย่อ)
```sql
-- migrations/20260104_crm_init.sql
CREATE TABLE crm_leads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    source VARCHAR(30),
    name VARCHAR(255), email VARCHAR(255), phone VARCHAR(30),
    company VARCHAR(255),
    interest_category VARCHAR(30),   -- SMART_FARM | SMART_BUILDING
    score INT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'NEW',-- NEW|CONTACTED|QUALIFIED|CONVERTED|LOST
    assigned_to UUID,
    converted_customer_id UUID,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE crm_opportunities (...);
CREATE TABLE crm_activities (...);
CREATE TABLE crm_tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL,
    subject VARCHAR(255),
    description TEXT,
    category VARCHAR(50),            # AI-classified
    priority VARCHAR(20) DEFAULT 'NORMAL',
    status VARCHAR(20) DEFAULT 'OPEN',
    sla_due_at TIMESTAMP,
    assigned_to UUID,
    resolved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE crm_campaigns (...);
```

### 6.3 Kafka Topics
```
crm.lead.created            # ← website, ads
crm.lead.converted          # → customer, package (proposal)
crm.opportunity.won         # → package, erp
crm.ticket.created          # → notifier
crm.ticket.sla.breach       # → notifier (escalate)
crm.ticket.resolved         # → report
```

---

## 7. โมดูลที่ 5 — IoT Device Management (`iotdevice`)

### 7.1 โครงสร้าง
```
internal/modules/iotdevice/
├── domain/
│   ├── entity/
│   │   ├── device.go               # อุปกรณ์ (Aggregate Root)
│   │   ├── device_model.go         # รุ่น
│   │   ├── device_group.go         # จัดกลุ่ม (farm zone, floor)
│   │   ├── firmware.go             # firmware version
│   │   ├── telemetry.go            # ข้อมูลที่ส่งเข้ามา
│   │   ├── command.go              # คำสั่งควบคุม
│   │   └── alert_rule.go           # กฎแจ้งเตือน
│   ├── value_object/
│   │   ├── device_status.go        # PROVISIONED|ONLINE|OFFLINE|FAULT
│   │   ├── protocol.go             # MQTT|HTTP|LORAWAN|MODBUS|ZIGBEE
│   │   ├── device_type.go          # SENSOR|GATEWAY|ACTUATOR|CAMERA
│   │   └── capability.go
│   ├── repository/
│   │   ├── device_repository.go
│   │   ├── telemetry_repository.go # (InfluxDB)
│   │   └── command_repository.go
│   ├── service/
│   │   ├── provisioning_service.go
│   │   ├── health_monitor_service.go
│   │   └── alert_evaluator.go
│   └── errors/
├── application/
│   ├── register_device.go
│   ├── provision_device.go         # gen cert/token, assign to site
│   ├── send_command.go
│   ├── ingest_telemetry.go         # ← MQTT consumer (hot path)
│   ├── evaluate_alert.go
│   ├── ota_update_firmware.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/influxdb/telemetry_repo.go
│   ├── persistence/redis/device_state_cache.go
│   ├── messaging/mqtt/
│   │   ├── mqtt_broker.go
│   │   └── topic_router.go
│   ├── messaging/kafka/
│   │   ├── telemetry_producer.go
│   │   └── consumers/telemetry_consumer.go
│   ├── search/elasticsearch/device_indexer.go
│   └── scheduler/offline_detector_job.go
├── interfaces/
│   ├── http/device_handler.go
│   ├── http/telemetry_handler.go
│   ├── websocket/live_telemetry_hub.go
│   └── routes.go
└── module.go
```

### 7.2 Database
```sql
-- migrations/20260105_iotdevice_init.sql
CREATE TABLE iotdevice_models (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vendor VARCHAR(100),
    model_no VARCHAR(100) UNIQUE,
    device_type VARCHAR(30),
    protocol VARCHAR(30),
    capabilities JSONB,          -- [{"metric":"temp","unit":"C"}]
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE iotdevice_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    site_id UUID NOT NULL,               -- customer_sites.id
    model_id UUID REFERENCES iotdevice_models(id),
    serial_no VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255),
    status VARCHAR(20) DEFAULT 'PROVISIONED',
    protocol VARCHAR(30),
    mqtt_client_id VARCHAR(100),
    device_token_hash VARCHAR(255),      -- bcrypt ของ device token
    firmware_version VARCHAR(50),
    last_seen_at TIMESTAMP,
    installed_at TIMESTAMP,
    group_id UUID,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_iot_devices_tenant_customer ON iotdevice_devices(tenant_id, customer_id);
CREATE INDEX idx_iot_devices_site ON iotdevice_devices(site_id);
CREATE INDEX idx_iot_devices_status ON iotdevice_devices(status, last_seen_at);

CREATE TABLE iotdevice_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id UUID NOT NULL,
    name VARCHAR(100),
    group_type VARCHAR(30),       -- ZONE | FLOOR | ROOM | BARN
    metadata JSONB
);

CREATE TABLE iotdevice_commands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL,
    command VARCHAR(100) NOT NULL,
    payload JSONB,
    status VARCHAR(20) DEFAULT 'PENDING',   -- PENDING|SENT|ACK|FAILED
    issued_by UUID,
    issued_at TIMESTAMP DEFAULT NOW(),
    acked_at TIMESTAMP
);
CREATE INDEX idx_iot_cmd_device ON iotdevice_commands(device_id, issued_at DESC);

CREATE TABLE iotdevice_alert_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id UUID,
    device_id UUID,
    metric VARCHAR(50) NOT NULL,
    operator VARCHAR(10) NOT NULL,        -- >, <, >=, ==
    threshold NUMERIC(15,4),
    severity VARCHAR(20),                 -- INFO|WARN|CRITICAL
    actions JSONB,                        -- [{"type":"notify","channel":"email"}]
    is_active BOOLEAN DEFAULT TRUE
);
```

**Telemetry เก็บใน InfluxDB** (time-series) — `iotdevice/persistence/influxdb/`:
```
measurement: device_telemetry
tags:        tenant_id, device_id, site_id, metric
fields:      value, quality
time:        event timestamp
```

### 7.3 MQTT Topic Convention
```
iot/{tenant_id}/{device_id}/telemetry       # device → platform
iot/{tenant_id}/{device_id}/status          # LWT / online-offline
iot/{tenant_id}/{device_id}/cmd             # platform → device
iot/{tenant_id}/{device_id}/cmd/ack         # device → platform
iot/{tenant_id}/{device_id}/ota             # firmware
```

### 7.4 Kafka Topics
```
iot.device.registered / iot.device.provisioned / iot.device.offline
iot.telemetry.raw                  # hot path
iot.telemetry.aggregated           # หลัง window 5s
iot.alert.triggered                # → notifier + crm (ticket)
iot.command.sent / iot.command.acked
iot.firmware.updated
```

### 7.5 HTTP Routes
```
POST   /api/v1/devices                         register
POST   /api/v1/devices/:id/provision
GET    /api/v1/devices?site_id=&status=
GET    /api/v1/devices/:id
POST   /api/v1/devices/:id/commands            {command, payload}
GET    /api/v1/devices/:id/telemetry?from=&to=
GET    /api/v1/sites/:site_id/devices
WS     /api/v1/devices/:id/live                WebSocket stream
POST   /api/v1/alert-rules
```

---

## 8. โมดูลที่ 6 — IoT Logistics (`iotlogistics`)

> บริหาร **อุปกรณ์ IoT** ในมุม logistics: shipment, installation, RMA, maintenance, spare parts

### 8.1 โครงสร้าง
```
internal/modules/iotlogistics/
├── domain/
│   ├── entity/
│   │   ├── shipment.go             # การขนส่งอุปกรณ์ไปติดตั้ง
│   │   ├── shipment_item.go
│   │   ├── installation_job.go     # งานติดตั้ง
│   │   ├── technician.go           # ช่าง
│   │   ├── maintenance_schedule.go # ตารางบำรุงรักษา
│   │   ├── work_order.go           # ใบสั่งงาน (ซ่อม/บำรุง)
│   │   ├── rma.go                  # เคลมสินค้า
│   │   └── spare_part.go
│   ├── value_object/
│   │   ├── shipment_status.go      # PENDING|IN_TRANSIT|DELIVERED|INSTALLED
│   │   ├── job_status.go           # SCHEDULED|IN_PROGRESS|DONE|FAILED
│   │   ├── route_optimizer.go      # VO สำหรับเส้นทาง
│   │   └── gps_location.go
│   ├── repository/
│   ├── service/
│   │   ├── route_planning_service.go
│   │   ├── sla_calculator.go
│   │   └── geo_fence_service.go
│   └── errors/
├── application/
│   ├── create_shipment.go
│   ├── dispatch_shipment.go
│   ├── confirm_delivery.go
│   ├── assign_installation_job.go
│   ├── complete_installation.go    # → iotdevice.provision + package.activate
│   ├── schedule_maintenance.go
│   ├── create_rma.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── services/maps/google_maps.go
│   ├── services/llm/route_advisor.go
│   ├── messaging/kafka/
│   └── scheduler/pm_due_job.go
├── interfaces/http/
│   ├── shipment_handler.go
│   ├── installation_handler.go
│   ├── maintenance_handler.go
│   └── routes.go
└── module.go
```

### 8.2 Database (ย่อ)
```sql
-- migrations/20260106_iotlogistics_init.sql
CREATE TABLE iotlogistics_shipments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shipment_no VARCHAR(30) UNIQUE NOT NULL,
    customer_id UUID NOT NULL,
    site_id UUID NOT NULL,
    from_warehouse_id UUID,          -- erp_warehouses
    carrier VARCHAR(100),
    tracking_no VARCHAR(100),
    status VARCHAR(30) DEFAULT 'PENDING',
    scheduled_at TIMESTAMP,
    delivered_at TIMESTAMP,
    gps_last JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE iotlogistics_shipment_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shipment_id UUID REFERENCES iotlogistics_shipments(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,       -- erp_products
    device_id UUID,                 -- iotdevice_devices (หลัง provision)
    qty NUMERIC(15,2),
    serials TEXT[]
);

CREATE TABLE iotlogistics_installation_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shipment_id UUID,
    site_id UUID NOT NULL,
    technician_id UUID,
    job_type VARCHAR(30),           -- INSTALL | MAINTENANCE | REPAIR | REMOVAL
    status VARCHAR(20) DEFAULT 'SCHEDULED',
    scheduled_at TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    checklist JSONB,
    photos TEXT[],
    signature_url TEXT,
    notes TEXT
);

CREATE TABLE iotlogistics_technicians (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,          -- อ้าง users module
    skill_set TEXT[],               -- ['SMART_FARM','SOLAR','NETWORK']
    zone VARCHAR(100),
    gps_last JSONB,
    is_available BOOLEAN DEFAULT TRUE
);

CREATE TABLE iotlogistics_maintenance_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL,
    pm_type VARCHAR(30),            -- MONTHLY|QUARTERLY|ANNUAL
    next_due_at TIMESTAMP,
    last_done_at TIMESTAMP,
    assigned_to UUID
);

CREATE TABLE iotlogistics_rma (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rma_no VARCHAR(30) UNIQUE NOT NULL,
    device_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    reason TEXT,
    status VARCHAR(20) DEFAULT 'OPEN', -- OPEN|SHIPPED|RECEIVED|REPAIRED|REPLACED
    replacement_device_id UUID,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 8.3 Kafka Topics
```
iotlogistics.shipment.created / dispatched / delivered
iotlogistics.installation.scheduled / started / completed
iotlogistics.maintenance.due       # → notifier + crm (ticket)
iotlogistics.rma.created / closed  # → erp (inventory adjust)
```

---

## 9. โมดูลที่ 7 — Reporting & Analytics (`report` — extend)

### 9.1 โครงสร้าง (ต่อจากของเดิม)
```
internal/modules/report/
├── domain/
│   ├── entity/
│   │   ├── report_definition.go
│   │   ├── report_schedule.go
│   │   └── kpi_snapshot.go
│   ├── value_object/
│   │   ├── report_type.go          # FINANCIAL|OPERATIONAL|IOT|SALES
│   │   └── export_format.go        # PDF|XLSX|CSV|JSON
│   └── repository/
├── application/
│   ├── generate_report.go
│   ├── schedule_report.go
│   ├── export_report.go
│   ├── get_dashboard_kpi.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/clickhouse/     # OLAP (optional)
│   ├── renderer/pdf_renderer.go    # ใช้ pkg/report
│   ├── renderer/xlsx_renderer.go
│   ├── templates/                  # HTML/Go template ต่อ report
│   └── scheduler/report_scheduler_job.go
├── interfaces/http/
│   ├── report_handler.go
│   ├── dashboard_handler.go
│   └── routes.go
└── module.go
```

### 9.2 Reports ที่ต้องมี

| Report | แหล่งข้อมูล | กลุ่มเป้าหมาย |
|---|---|---|
| Customer 360 | customer + crm + package | Sales, CS |
| Revenue by Package | package + erp | Exec |
| Device Uptime & Health | iotdevice + influxdb | Ops |
| Telemetry Analytics (Smart Farm) | influxdb (temp, humidity, soil) | ลูกค้าเกษตร |
| Energy Analytics (Smart Building) | influxdb (kWh, power) | ลูกค้าอาคาร |
| Logistics SLA | iotlogistics | Ops |
| Maintenance Due | iotlogistics | Service |
| Ticket SLA | crm | CS Manager |
| Inventory Aging | erp | WH |
| Usage vs Quota | package | Sales/CS |

### 9.3 Database
```sql
-- migrations/20260107_report_init.sql
CREATE TABLE report_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE,
    name VARCHAR(255),
    report_type VARCHAR(30),
    query_config JSONB,
    template_path VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE report_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id UUID NOT NULL REFERENCES report_definitions(id),
    cron_expr VARCHAR(50),
    recipients TEXT[],
    format VARCHAR(10),
    filters JSONB,
    last_run_at TIMESTAMP,
    next_run_at TIMESTAMP
);

CREATE TABLE report_kpi_snapshots (
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID,
    kpi_code VARCHAR(50),
    value NUMERIC(20,4),
    dimensions JSONB,
    snapshot_date DATE NOT NULL,
    UNIQUE(tenant_id, kpi_code, snapshot_date)
);
```

### 9.4 Kafka Topics
```
report.generated            # หลังสร้างเสร็จ → sendEmail
report.kpi.snapshot.updated # ← scheduler รายวัน
```

---

## 10. Cross-cutting Concerns

### 10.1 Tenant Isolation
ทุกโมดูลต้องมี middleware แทรก `tenant_id` และ Repository กรองทุก query:
```go
// internal/middleware/tenant.go
func Inject() gin.HandlerFunc {
    return func(c *gin.Context) {
        tid := c.GetHeader("X-Tenant-ID")
        if tid == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "missing tenant"})
            return
        }
        id, err := uuid.Parse(tid)
        if err != nil {
            c.AbortWithStatusJSON(400, gin.H{"error": "invalid tenant"})
            return
        }
        c.Set("tenant_id", id)
        c.Next()
    }
}
```

### 10.2 Event-driven Integration Matrix

| Producer | Topic | Consumer |
|---|---|---|
| `crm` | `crm.opportunity.won` | `package`, `erp` |
| `customer` | `customer.created` | `crm`, `package`, `erp` |
| `package` | `package.subscription.created` | `payment`, `iotdevice`, `notifier` |
| `erp` | `erp.invoice.issued` | `payment`, `report` |
| `iotdevice` | `iot.alert.triggered` | `notifier`, `crm` |
| `iotdevice` | `iot.telemetry.aggregated` | `report`, `package` (usage) |
| `iotlogistics` | `iotlogistics.installation.completed` | `iotdevice` (activate), `package` (start billing) |
| `iotlogistics` | `iotlogistics.maintenance.due` | `notifier`, `crm` |
| `report` | `report.generated` | `notifier` (`sendEmail`) |

### 10.3 Kafka Topic Convention (ทั้งแพลตฟอร์ม)
```
<module>.<entity>.<action>
ตัวอย่าง:
  customer.created
  package.subscription.created
  erp.invoice.issued
  crm.ticket.created
  iot.telemetry.raw
  iotlogistics.installation.completed
  report.generated
```

### 10.4 Shared Packages ที่ใช้
```
pkg/helpers, pkg/db, pkg/transaction, pkg/logger, pkg/jwt,
pkg/httpErrors, pkg/responses, pkg/kafka, pkg/mqtt, pkg/influxdb,
pkg/elasticsearch, pkg/llm, pkg/sendEmail, pkg/emailTemplates,
pkg/report, pkg/secureRandom, pkg/cryptpass, pkg/utils
```

---

## 11. Bootstrap — Composition Root

**`cmd/api/main.go`** (ย่อ):
```go
func main() {
    // ... init db, redis, kafka, mqtt, influx
    api := r.Group("/api/v1")
    auth := middleware.Auth(cfg.JWT.Secret)
    tenant := middleware.Inject()

    customer.Init(api, customer.Deps{DB: db, Producer: kp}, auth, tenant)
    pkgcatalog.Init(api, pkgcatalog.Deps{DB: db, Redis: rdb, Producer: kp}, auth, tenant)
    erp.Init(api, erp.Deps{DB: db}, auth, tenant)
    crm.Init(api, crm.Deps{DB: db, Producer: kp, LLM: llmClient}, auth, tenant)
    iotdevice.Init(api, iotdevice.Deps{DB: db, Influx: influx, MQTT: mqtt, WS: hub}, auth, tenant)
    iotlogistics.Init(api, iotlogistics.Deps{DB: db, Maps: mapsClient}, auth, tenant)
    report.Init(api, report.Deps{DB: db, Renderer: pdfRenderer}, auth, tenant)

    r.Run(":8080")
}
```

---

## 12. Migration Files สรุป

| ลำดับ | ไฟล์ | โมดูล |
|---|---|---|
| 1 | `20260101_customer_init.sql` | customer |
| 2 | `20260102_package_init.sql` | package |
| 3 | `20260103_erp_init.sql` | erp |
| 4 | `20260104_crm_init.sql` | crm |
| 5 | `20260105_iotdevice_init.sql` | iotdevice |
| 6 | `20260106_iotlogistics_init.sql` | iotlogistics |
| 7 | `20260107_report_init.sql` | report |
| 8 | `20260108_seed_iot_packages.sql` | seed ข้อมูลแพ็กเกจ Smart Farm/Building |

---

## 13. Prompt Template (ย่อ) สำหรับ AI สร้างโค้ด

```
สร้างโมดูล <module_name> ในโปรเจกต์ icmongolang ตาม Template_Module.md

Input:
- module_name: <iotdevice|package|erp|crm|iotlogistics|...>
- entity: <Device|Subscription|PurchaseOrder|Ticket|Shipment>
- actions: <Create, Get, Update, SendCommand, ...>
- cross-cutting: <Kafka|MQTT|InfluxDB|ES|WebSocket|LLM>
- DB prefix: <module_name>_
- Multi-tenant: yes (tenant_id ทุกตาราง)

กฎเพิ่มเติม (นอกเหนือจาก Template_Module.md):
1. ทุกตารางต้องมี tenant_id UUID NOT NULL + index
2. Repository ทุก query ต้องมี WHERE tenant_id = ?
3. Kafka topic ใช้ pattern <module>.<entity>.<action>
4. MQTT topic ใช้ pattern iot/{tenant_id}/{device_id}/{action}
5. Telemetry เก็บใน InfluxDB, metadata ใน PostgreSQL
6. Use case ที่ข้ามโมดูล → ผ่าน Kafka เท่านั้น (ห้าม import ตรง)
7. Entity ที่เป็น Aggregate Root ต้องมี behavior method + constructor
8. ทุก endpoint ต้องผ่าน auth + tenant.Inject() middleware
9. Run: go build ./... && go test ./...
```

---

## 14. Checklist สรุปต่อโมดูล

| # | โมดูล | ตาราง | Kafka | MQTT | Influx | WS | LLM |
|---|---|:-:|:-:|:-:|:-:|:-:|:-:|
| 1 | customer | 3 | ✅ | – | – | – | – |
| 2 | package | 4 | ✅ | – | – | – | – |
| 3 | erp | 8 | ✅ | – | – | – | – |
| 4 | crm | 5 | ✅ | – | – | – | ✅ |
| 5 | iotdevice | 5 | ✅ | ✅ | ✅ | ✅ | – |
| 6 | iotlogistics | 6 | ✅ | – | – | ✅ | ✅ |
| 7 | report | 3 | ✅ | – | (query) | ✅ | – |

---

ต้องการให้ผม **ลงลึกโมดูลใดโมดูลหนึ่ง** (เช่น สร้างโค้ดเต็มของ `iotdevice` พร้อม handler + usecase + repository + test) หรือ **generate migration SQL ของทั้ง 7 โมดูล** ให้ครบเลยไหมครับ?