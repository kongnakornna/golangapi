# 🏗️ โครงสร้าง Modules ERP+SaaS (จัดเรียงใหม่ตาม Prefix)

## 📐 หลักการออกแบบ

| หลักการ | คำอธิบาย |
|---|---|
| **Clean Architecture** | แยกชั้น `domain → application → infrastructure → interfaces` |
| **DDD** | แต่ละโมดูลมี Aggregate Root, Value Object, Repository |
| **Multi-tenant** | ทุก layer ต้องมี `tenant_id` |
| **Modular Monolith** | แยกโมดูลชัดเจน แต่ deploy รวมกันได้ |
| **Event-Driven** | สื่อสารระหว่างโมดูลผ่าน Kafka |

---

## 📂 โครงสร้าง Prefix Modules

```
internal/modules/
│
├── erp/                          # 🏢 กลุ่ม ERP
│   ├── inventory/                # 1.1 ✅ พร้อมใช้งาน
│   ├── production/               # 1.2 ⏳ กำลังพัฒนา
│   └── finance/                  # 1.3 ⏳ กำลังพัฒนา
│
├── crm/                          # 🤝 กลุ่ม CRM
│   └── crm/                      # 1.4 ⏳ กำลังพัฒนา
│
├── pos/                          # 🛒 กลุ่ม POS
│   └── pos/                      # 1.5 ⏳ กำลังพัฒนา
│
├── ecommerce/                    # 🌐 กลุ่ม E-commerce
│   └── ecommerce/                # 1.6 ⏳ กำลังพัฒนา
│
├── logistics/                    # 🚚 กลุ่ม Logistics
│   └── logistics/                # 1.7 ⏳ กำลังพัฒนา
│
├── iot/                          # 📡 กลุ่ม IoT/QR
│   └── iot/                      # 1.8 ⏳ กำลังพัฒนา
│
├── forecast/                     # 🤖 กลุ่ม AI
│   └── forecast/                 # 1.9 ⏳ กำลังพัฒนา
│
└── reporting/                    # 📊 กลุ่มรายงาน
    └── reporting/                # 1.10 ⏳ กำลังพัฒนา
```

---

## 1️⃣ `erp/inventory` – สินค้าคงคลัง ✅

```
internal/modules/erp/inventory/
├── domain/
│   ├── entity/
│   │   ├── product.go              # Aggregate Root
│   │   ├── warehouse.go
│   │   ├── stock_movement.go
│   │   ├── stock_balance.go
│   │   └── unit_of_measure.go
│   ├── value_object/
│   │   ├── sku.go
│   │   ├── quantity.go
│   │   ├── stock_status.go
│   │   └── movement_type.go
│   ├── repository/
│   │   ├── product_repository.go
│   │   ├── warehouse_repository.go
│   │   └── stock_repository.go
│   ├── service/
│   │   └── inventory_domain_service.go
│   └── errors/errors.go
├── application/
│   ├── create_product.go
│   ├── adjust_stock.go
│   ├── transfer_stock.go
│   ├── reserve_stock.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   │   ├── product_repo_impl.go
│   │   ├── stock_repo_impl.go
│   │   └── models.go
│   ├── persistence/redis/stock_cache.go
│   └── messaging/kafka_producer.go
├── interfaces/
│   ├── http/
│   │   ├── product_handler.go
│   │   ├── stock_handler.go
│   │   └── routes.go
│   └── middleware/tenant.go
└── module.go
```

---

## 2️⃣ `erp/production` – การผลิต ⏳

```
internal/modules/erp/production/
├── domain/
│   ├── entity/
│   │   ├── production_order.go     # Aggregate Root
│   │   ├── bom.go
│   │   ├── work_center.go
│   │   ├── routing.go
│   │   └── quality_check.go
│   ├── value_object/
│   │   ├── order_status.go
│   │   ├── production_qty.go
│   │   └── oee.go
│   ├── repository/
│   │   ├── production_order_repository.go
│   │   ├── bom_repository.go
│   │   └── work_center_repository.go
│   ├── service/
│   │   └── production_domain_service.go
│   └── errors/errors.go
├── application/
│   ├── create_production_order.go
│   ├── release_order.go
│   ├── record_progress.go
│   ├── complete_order.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/redis/production_cache.go
│   └── messaging/kafka_producer.go
├── interfaces/
│   ├── http/
│   └── middleware/tenant.go
└── module.go
```

---

## 3️⃣ `erp/finance` – การเงิน ⏳

```
internal/modules/erp/finance/
├── domain/
│   ├── entity/
│   │   ├── invoice.go              # Aggregate Root
│   │   ├── payment.go
│   │   ├── journal_entry.go
│   │   ├── account.go
│   │   └── tax_invoice.go
│   ├── value_object/
│   │   ├── money.go
│   │   ├── invoice_status.go
│   │   ├── tax_id.go
│   │   └── vat_rate.go
│   ├── repository/
│   │   ├── invoice_repository.go
│   │   ├── payment_repository.go
│   │   └── journal_repository.go
│   ├── service/
│   │   └── finance_domain_service.go
│   └── errors/errors.go
├── application/
│   ├── create_invoice.go
│   ├── record_payment.go
│   ├── issue_tax_invoice.go
│   ├── close_period.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/redis/finance_cache.go
│   └── messaging/kafka_producer.go
├── interfaces/
│   ├── http/
│   └── middleware/tenant.go
└── module.go
```

---

## 4️⃣ `crm` – CRM ⏳

```
internal/modules/crm/crm/
├── domain/
│   ├── entity/
│   │   ├── lead.go                 # Aggregate Root
│   │   ├── opportunity.go
│   │   ├── activity.go
│   │   ├── campaign.go
│   │   └── customer_360.go
│   ├── value_object/
│   │   ├── lead_status.go
│   │   ├── opportunity_stage.go
│   │   ├── source.go
│   │   └── score.go
│   ├── repository/
│   │   ├── lead_repository.go
│   │   ├── opportunity_repository.go
│   │   └── activity_repository.go
│   ├── service/
│   │   └── crm_domain_service.go
│   └── errors/errors.go
├── application/
│   ├── create_lead.go
│   ├── convert_lead.go
│   ├── create_opportunity.go
│   ├── log_activity.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/redis/crm_cache.go
│   ├── search/elasticsearch/lead_indexer.go
│   └── messaging/kafka_producer.go
├── interfaces/
│   ├── http/
│   └── middleware/tenant.go
└── module.go
```

---

## 5️⃣ `pos` – POS ⏳

```
internal/modules/pos/pos/
├── domain/
│   ├── entity/
│   │   ├── sale.go                 # Aggregate Root
│   │   ├── terminal.go
│   │   ├── shift.go
│   │   ├── receipt.go
│   │   └── payment_method.go
│   ├── value_object/
│   │   ├── sale_status.go
│   │   ├── payment_type.go
│   │   └── discount.go
│   ├── repository/
│   │   ├── sale_repository.go
│   │   ├── terminal_repository.go
│   │   └── shift_repository.go
│   ├── service/
│   │   └── pos_domain_service.go
│   └── errors/errors.go
├── application/
│   ├── open_shift.go
│   ├── create_sale.go
│   ├── void_sale.go
│   ├── close_shift.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/redis/pos_cache.go
│   └── messaging/kafka_producer.go
├── interfaces/
│   ├── http/
│   └── middleware/tenant.go
└── module.go
```

---

## 6️⃣ `ecommerce` – E-commerce ⏳

```
internal/modules/ecommerce/ecommerce/
├── domain/
│   ├── entity/
│   │   ├── cart.go                 # Aggregate Root
│   │   ├── order.go
│   │   ├── product_listing.go
│   │   ├── shipping.go
│   │   └── promotion.go
│   ├── value_object/
│   │   ├── order_status.go
│   │   ├── cart_item.go
│   │   └── channel.go
│   ├── repository/
│   │   ├── cart_repository.go
│   │   ├── order_repository.go
│   │   └── listing_repository.go
│   ├── service/
│   │   └── ecommerce_domain_service.go
│   └── errors/errors.go
├── application/
│   ├── add_to_cart.go
│   ├── checkout.go
│   ├── confirm_payment.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/redis/cart_cache.go
│   └── messaging/kafka_producer.go
├── interfaces/
│   ├── http/
│   └── middleware/tenant.go
└── module.go
```

---

## 7️⃣ `logistics` – Logistics ⏳

```
internal/modules/logistics/logistics/
├── domain/
│   ├── entity/
│   │   ├── shipment.go             # Aggregate Root
│   │   ├── carrier.go
│   │   ├── route.go
│   │   ├── tracking_event.go
│   │   └── delivery_proof.go
│   ├── value_object/
│   │   ├── shipment_status.go
│   │   ├── address.go
│   │   └── tracking_number.go
│   ├── repository/
│   │   ├── shipment_repository.go
│   │   ├── carrier_repository.go
│   │   └── tracking_repository.go
│   ├── service/
│   │   └── logistics_domain_service.go
│   └── errors/errors.go
├── application/
│   ├── create_shipment.go
│   ├── assign_carrier.go
│   ├── update_tracking.go
│   ├── confirm_delivery.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/redis/tracking_cache.go
│   └── messaging/kafka_producer.go
├── interfaces/
│   ├── http/
│   └── middleware/tenant.go
└── module.go
```

---

## 8️⃣ `iot` – IoT / QR ⏳

```
internal/modules/iot/iot/
├── domain/
│   ├── entity/
│   │   ├── device.go               # Aggregate Root
│   │   ├── sensor_reading.go
│   │   ├── qr_code.go
│   │   ├── scan_event.go
│   │   └── alert_rule.go
│   ├── value_object/
│   │   ├── device_status.go
│   │   ├── sensor_type.go
│   │   ├── qr_payload.go
│   │   └── alert_level.go
│   ├── repository/
│   │   ├── device_repository.go
│   │   ├── reading_repository.go
│   │   └── qr_repository.go
│   ├── service/
│   │   └── iot_domain_service.go
│   └── errors/errors.go
├── application/
│   ├── register_device.go
│   ├── ingest_reading.go
│   ├── generate_qr.go
│   ├── scan_qr.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/timeseries/
│   ├── persistence/redis/device_cache.go
│   ├── messaging/mqtt/consumer.go
│   └── messaging/kafka_producer.go
├── interfaces/
│   ├── http/
│   ├── websocket/realtime.go
│   └── middleware/tenant.go
└── module.go
```

---

## 9️⃣ `forecast` – AI Forecast ⏳

```
internal/modules/forecast/forecast/
├── domain/
│   ├── entity/
│   │   ├── forecast_model.go       # Aggregate Root
│   │   ├── prediction.go
│   │   ├── training_data.go
│   │   └── model_metric.go
│   ├── value_object/
│   │   ├── model_type.go
│   │   ├── forecast_horizon.go
│   │   ├── confidence.go
│   │   └── accuracy.go
│   ├── repository/
│   │   ├── model_repository.go
│   │   └── prediction_repository.go
│   ├── service/
│   │   └── forecast_domain_service.go
│   └── errors/errors.go
├── application/
│   ├── train_model.go
│   ├── predict_demand.go
│   ├── retrain_model.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/timeseries/features.go
│   ├── ml/python_bridge.go
│   └── messaging/kafka_producer.go
├── interfaces/
│   ├── http/
│   └── middleware/tenant.go
└── module.go
```

---

## 🔟 `reporting` – รายงาน / KPI ⏳

```
internal/modules/reporting/reporting/
├── domain/
│   ├── entity/
│   │   ├── report.go               # Aggregate Root
│   │   ├── dashboard.go
│   │   ├── kpi_definition.go
│   │   ├── kpi_value.go
│   │   └── schedule.go
│   ├── value_object/
│   │   ├── report_type.go
│   │   ├── kpi_code.go
│   │   ├── period.go
│   │   └── export_format.go
│   ├── repository/
│   │   ├── report_repository.go
│   │   ├── dashboard_repository.go
│   │   └── kpi_repository.go
│   ├── service/
│   │   └── reporting_domain_service.go
│   └── errors/errors.go
├── application/
│   ├── generate_report.go
│   ├── create_dashboard.go
│   ├── calculate_kpi.go
│   ├── schedule_report.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/clickhouse/
│   ├── persistence/redis/report_cache.go
│   └── messaging/kafka_consumer.go
├── interfaces/
│   ├── http/
│   └── middleware/tenant.go
└── module.go
```

---

## 📊 สรุปภาพรวมทั้งหมด

| # | Prefix | โมดูล | Aggregate Root | สถานะ |
|---|---|---|---|---|
| 1 | `erp` | inventory | Product, Warehouse | ✅ |
| 2 | `erp` | production | ProductionOrder, BOM | ⏳ |
| 3 | `erp` | finance | Invoice, Payment | ⏳ |
| 4 | `crm` | crm | Lead, Opportunity | ⏳ |
| 5 | `pos` | pos | Sale, Shift | ⏳ |
| 6 | `ecommerce` | ecommerce | Cart, Order | ⏳ |
| 7 | `logistics` | logistics | Shipment, Carrier | ⏳ |
| 8 | `iot` | iot | Device, QRCode | ⏳ |
| 9 | `forecast` | forecast | ForecastModel | ⏳ |
| 10 | `reporting` | reporting | Report, Dashboard | ⏳ |

---

## 🔗 การสื่อสารระหว่างโมดูล (Event-Driven)

```
┌──────────────────┐   order.placed    ┌──────────────────┐
│ ecommerce/ecommerce│ ───────────────▶│ erp/inventory    │
└──────────────────┘                   └──────────────────┘
                                              │
                                              │ stock.reserved
                                              ▼
┌──────────────────┐  production.done  ┌──────────────────┐
│ erp/production   │ ───────────────▶│ erp/inventory    │
└──────────────────┘                   └──────────────────┘
                                              │
                                              │ stock.updated
                                              ▼
┌──────────────────┐  invoice.issued   ┌──────────────────┐
│ erp/finance      │ ◀───────────────│ pos/pos          │
└──────────────────┘                   └──────────────────┘
       │
       │ payment.received
       ▼
┌──────────────────┐   kpi.updated     ┌──────────────────┐
│ reporting        │ ◀───────────────│ ทุกโมดูล          │
└──────────────────┘                   └──────────────────┘
```

---

## 🎯 หลักการสำคัญที่ทุกโมดูลต้องมี

| หลักการ | รายละเอียด |
|---|---|
| **tenant_id** | ทุก entity, repository, handler |
| **domain/errors** | Error ของแต่ละโมดูล |
| **module.go** | Entry point wire dependency |
| **dto.go** | แยก request/response |
| **middleware/tenant.go** | inject tenant จาก JWT |
| **kafka_producer** | ส่ง event เมื่อมีการเปลี่ยนแปลง |
| **repository interface** | อยู่ใน domain, implement อยู่ใน infrastructure |

---

**เวอร์ชัน**: 1.0
**ประเภทเอกสาร**: โครงสร้างโมดูลสำหรับทีมพัฒนา
**อ้างอิง**: เอกสารภาพรวมระบบ ERP+SaaS



# ERP มีอะไรบ้าง — ภาพรวมทั้งหมด

**ERP (Enterprise Resource Planning)** คือระบบที่รวมทุกกระบวนการทางธุรกิจไว้ในแพลตฟอร์มเดียว ต่อไปนี้คือ **โมดูลทั้งหมด** ที่ ERP มาตรฐานทั่วไปมี และผมจะเชื่อมโยงกับ **ระบบฟาร์มเห็ดอัจฉริยะ** ที่คุณกำลังออกแบบด้วย

---

## 1. ภาพรวมโมดูล ERP ทั้งหมด

```
┌─────────────────────────────────────────────────────────────────────┐
│                         ERP CORE MODULES                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │   Finance    │  │   Inventory  │  │  Production  │              │
│  │   & Accounting│  │  & Warehouse │  │  & Manufacturing│           │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
│                                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │  Procurement │  │     Sales    │  │     HRM      │              │
│  │  & Purchasing│  │   & Order    │  │  & Payroll   │              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
│                                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │     CRM      │  │   Supply     │  │   Quality    │              │
│  │              │  │    Chain     │  │  Management  │              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
│                                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │   Project    │  │    Asset     │  │   Reporting  │              │
│  │  Management  │  │  Management  │  │  & Analytics │              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 2. โมดูลหลักของ ERP (11 โมดูล)

### 2.1 Finance & Accounting (การเงินและบัญชี)

| ระบบย่อย | หน้าที่ | ใช้ในฟาร์มเห็ด |
|:---|:---|:---|
| **General Ledger** | บัญชีแยกประเภท | บันทึกรายรับ-รายจ่ายทั้งหมด |
| **Accounts Payable** | เจ้าหนี้การค้า | จ่ายค่าปุ๋ย, ค่าไฟ, ค่าแรง |
| **Accounts Receivable** | ลูกหนี้การค้า | เก็บเงินจากลูกค้า |
| **Cash Management** | กระแสเงินสด | ตรวจสอบเงินสดในมือ |
| **Fixed Assets** | สินทรัพย์ถาวร | โรงเพาะ, เครื่องจักร |
| **Tax Management** | ภาษี | VAT, ภาษีเงินได้ |
| **Cost Accounting** | บัญชีต้นทุน | ต้นทุนต่อก้อนเชื้อ |
| **Budgeting** | งบประมาณ | วางแผนรายจ่าย |

**ตัวอย่างในฟาร์มเห็ด:**
```
- คิดต้นทุนต่อก้อนเชื้อ = (ค่าวัสดุ + ค่าแรง + ค่าไฟ + ค่าน้ำ) / จำนวนก้อน
- คิดต้นทุนต่อกิโลกรัม = ต้นทุนรวม / ผลผลิตทั้งหมด
- คิดกำไรต่อช่องทาง = รายได้ - ต้นทุน - ค่าขนส่ง
```

---

### 2.2 Inventory & Warehouse (สินค้าคงคลังและคลังสินค้า)

| ระบบย่อย | หน้าที่ | ใช้ในฟาร์มเห็ด |
|:---|:---|:---|
| **Stock Management** | จัดการสต็อก | นับก้อนเชื้อ, เห็ดพร้อมขาย |
| **Warehouse Management** | จัดการคลัง | แยกโซนห้องเย็น, ห้องเก็บ |
| **Lot/Batch Tracking** | ติดตามล็อต | ติดตามก้อนเชื้อแต่ละรุ่น |
| **Serial Number** | หมายเลขเครื่อง | ติดตามอุปกรณ์ |
| **Barcode/QR/RFID** | สแกน | สแกนก้อนเชื้อ |
| **Reorder Point** | จุดสั่งซื้อซ้ำ | เตือนเมื่อสต็อกต่ำ |
| **Multi-Location** | หลายคลัง | หลายโรงเพาะ |
| **Stock Adjustment** | ปรับสต็อก | เมื่อของเสียหาย |

**ตัวอย่างในฟาร์มเห็ด:**
```
- สต็อกก้อนเชื้อ: 5,000 ก้อน (พร้อมเปิดดอก 2,000, บ่ม 3,000)
- สต็อกเห็ดสด: 500 กก. (ต้องขายภายใน 3 วัน)
- สต็อกวัสดุ: ขี้เลื่อย 2 ตัน, หัวเชื้อ 50 ขวด
- จุดสั่งซื้อซ้ำ: ขี้เลื่อย < 500 กก. → สั่งซื้อ
```

---

### 2.3 Production & Manufacturing (การผลิต)

| ระบบย่อย | หน้าที่ | ใช้ในฟาร์มเห็ด |
|:---|:---|:---|
| **BOM (Bill of Materials)** | สูตรการผลิต | สูตรก้อนเชื้อ |
| **Production Planning** | วางแผนผลิต | แผนเพาะรายสัปดาห์ |
| **Work Order** | ใบสั่งผลิต | ใบสั่งทำก้อนเชื้อ |
| **Routing** | เส้นทางการผลิต | ขั้นตอนการเพาะ |
| **Shop Floor Control** | ควบคุมหน้างาน | ติดตามการผลิต |
| **Capacity Planning** | วางแผนกำลังผลิต | จำนวนโรงเพาะ |
| **MRP (Material Requirements Planning)** | วางแผนวัสดุ | ต้องการขี้เลื่อยเท่าไร |
| **Scheduling** | จัดตาราง | ตารางเก็บเกี่ยว |

**ตัวอย่างในฟาร์มเห็ด:**
```
BOM สำหรับก้อนเชื้อ 1,000 ก้อน:
- ขี้เลื่อย 800 กก.
- รำละเอียด 100 กก.
- ปุ๋ยยูเรีย 5 กก.
- หัวเชื้อ 20 ขวด
- ถุงพลาสติก 1,000 ใบ
- ค่าแรง 10 ชม.
- ค่าไฟ 50 kWh
```

---

### 2.4 Procurement & Purchasing (จัดซื้อ)

| ระบบย่อย | หน้าที่ | ใช้ในฟาร์มเห็ด |
|:---|:---|:---|
| **Purchase Requisition** | ใบขอซื้อ | ขอซื้อขี้เลื่อย |
| **Purchase Order** | ใบสั่งซื้อ | สั่งซื้อวัสดุ |
| **Supplier Management** | จัดการซัพพลายเออร์ | ฐานข้อมูลผู้ขาย |
| **Quotation** | ใบเสนอราคา | เปรียบเทียบราคา |
| **Goods Receipt** | รับสินค้า | ตรวจรับวัสดุ |
| **Invoice Matching** | กระทบใบเสร็จ | ตรวจสอบใบแจ้งหนี้ |
| **Payment** | จ่ายเงิน | จ่ายค่าซัพพลายเออร์ |

**ตัวอย่างในฟาร์มเห็ด:**
```
- ซัพพลายเออร์ขี้เลื่อย: 3 ราย, ราคา 2-2.5 บาท/กก.
- ซัพพลายเออร์หัวเชื้อ: 2 ราย, ราคา 50-60 บาท/ขวด
- ซัพพลายเออร์ถุงพลาสติก: 1 ราย, ราคา 0.5 บาท/ใบ
```

---

### 2.5 Sales & Order Management (การขายและคำสั่งซื้อ)

| ระบบย่อย | หน้าที่ | ใช้ในฟาร์มเห็ด |
|:---|:---|:---|
| **Quotation** | ใบเสนอราคา | เสนอราคาให้ลูกค้า |
| **Sales Order** | ใบสั่งขาย | ออเดอร์จากลูกค้า |
| **Delivery** | ใบส่งของ | ส่งเห็ดให้ลูกค้า |
| **Invoice** | ใบแจ้งหนี้ | เรียกเก็บเงิน |
| **Credit Note** | ใบลดหนี้ | คืนสินค้า |
| **Price List** | ราคาขาย | ราคาตามช่องทาง |
| **Discount** | ส่วนลด | โปรโมชั่น |
| **Sales Commission** | ค่านายหน้า | ค่านายหน้าเซลล์ |

**ตัวอย่างในฟาร์มเห็ด:**
```
ช่องทางขาย:
- ตลาดสด: 80 บาท/กก.
- ห้างสรรพสินค้า: 120 บาท/กก.
- ออนไลน์ (Shopee/Lazada): 150 บาท/กก.
- ส่งโรงงาน: 60 บาท/กก.
- ร้านอาหาร: 100 บาท/กก.
```

---

### 2.6 CRM (Customer Relationship Management)

| ระบบย่อย | หน้าที่ | ใช้ในฟาร์มเห็ด |
|:---|:---|:---|
| **Customer Database** | ฐานข้อมูลลูกค้า | ประวัติลูกค้า |
| **Contact Management** | จัดการผู้ติดต่อ | เบอร์, อีเมล |
| **Lead Management** | จัดการลีด | ลูกค้าใหม่ |
| **Opportunity** | โอกาสขาย | ติดตามดีล |
| **Campaign** | แคมเปญ | โปรโมชั่น |
| **Customer Service** | บริการลูกค้า | รับเรื่องร้องเรียน |
| **Loyalty Program** | สะสมแต้ม | ลูกค้าประจำ |
| **Feedback** | ความคิดเห็น | รีวิว |

**ตัวอย่างในฟาร์มเห็ด:**
```
- ลูกค้า A: ร้านอาหาร, สั่ง 50 กก./สัปดาห์
- ลูกค้า B: ห้างสรรพสินค้า, สั่ง 200 กก./สัปดาห์
- ลูกค้า C: ขายออนไลน์, สั่ง 5 กก./ครั้ง
- Campaign: ซื้อ 10 กก. แถม 1 กก.
```

---

### 2.7 Supply Chain Management (ห่วงโซ่อุปทาน)

| ระบบย่อย | หน้าที่ | ใช้ในฟาร์มเห็ด |
|:---|:---|:---|
| **Demand Planning** | วางแผนอุปสงค์ | พยากรณ์ความต้องการ |
| **Supply Planning** | วางแผนอุปทาน | วางแผนจัดหา |
| **Logistics** | โลจิสติกส์ | ขนส่ง |
| **Transportation** | การขนส่ง | เส้นทาง |
| **Fleet Management** | จัดการยานพาหนะ | รถขนส่ง |
| **Route Optimization** | เส้นทางที่เหมาะสม | ประหยัดน้ำมัน |
| **Cold Chain** | ห่วงโซ่ความเย็น | รักษาอุณหภูมิ |
| **Track & Trace** | ติดตาม | GPS |

**ตัวอย่างในฟาร์มเห็ด:**
```
เส้นทางขนส่ง:
- ฟาร์ม → ตลาดสด (20 กม., 30 นาที)
- ฟาร์ม → ห้าง (50 กม., 1 ชม.)
- ฟาร์ม → โรงงาน (100 กม., 2 ชม.)
อุณหภูมิระหว่างขนส่ง: 2-8°C
```

---

### 2.8 HRM & Payroll (ทรัพยากรบุคคลและเงินเดือน)

| ระบบย่อย | หน้าที่ | ใช้ในฟาร์มเห็ด |
|:---|:---|:---|
| **Employee Database** | ฐานข้อมูลพนักงาน | ประวัติพนักงาน |
| **Recruitment** | สรรหา | รับสมัคร |
| **Attendance** | เวลาทำงาน | เช็คอิน-เอาท์ |
| **Leave Management** | ลา | ลาป่วย, ลากิจ |
| **Payroll** | เงินเดือน | คำนวณเงินเดือน |
| **Performance** | ประเมินผล | KPI พนักงาน |
| **Training** | อบรม | อบรมพนักงาน |
| **Benefits** | สวัสดิการ | ประกัน, โบนัส |

**ตัวอย่างในฟาร์มเห็ด:**
```
พนักงาน:
- ผู้จัดการฟาร์ม: 1 คน, 30,000 บาท/เดือน
- พนักงานเพาะ: 5 คน, 12,000 บาท/เดือน
- พนักงานเก็บเกี่ยว: 10 คน, 350 บาท/วัน
- พนักงานขนส่ง: 2 คน, 15,000 บาท/เดือน
- เซลล์: 2 คน, 15,000 + คอมมิชชั่น
```

---

### 2.9 Quality Management (การจัดการคุณภาพ)

| ระบบย่อย | หน้าที่ | ใช้ในฟาร์มเห็ด |
|:---|:---|:---|
| **Quality Control** | ควบคุมคุณภาพ | ตรวจสอบเห็ด |
| **Quality Assurance** | ประกันคุณภาพ | มาตรฐาน |
| **Inspection** | ตรวจสอบ | ตรวจรับวัสดุ |
| **Non-Conformance** | ของไม่ได้มาตรฐาน | เห็ดเสีย |
| **CAPA** | แก้ไขและป้องกัน | แก้ปัญหา |
| **Certification** | ใบรับรอง | GAP, Organic |
| **Traceability** | ติดตามย้อนกลับ | ตรวจสอบที่มา |
| **Audit** | ตรวจสอบ | ตรวจประเมิน |

**ตัวอย่างในฟาร์มเห็ด:**
```
มาตรฐาน:
- GAP (Good Agricultural Practice)
- Organic Thailand
- อย. (อาหารและยา)
- HACCP

เกรดเห็ด:
- เกรด A: ดอกใหญ่, สีสวย, ไม่ช้ำ
- เกรด B: ดอกกลาง, สีพอใช้
- เกรด C: ดอกเล็ก, ช้ำ, ส่งโรงงาน
```

---

### 2.10 Asset Management (การจัดการสินทรัพย์)

| ระบบย่อย | หน้าที่ | ใช้ในฟาร์มเห็ด |
|:---|:---|:---|
| **Asset Register** | ทะเบียนสินทรัพย์ | โรงเพาะ, เครื่องจักร |
| **Depreciation** | ค่าเสื่อมราคา | คำนวณค่าเสื่อม |
| **Maintenance** | บำรุงรักษา | ซ่อมบำรุง |
| **Preventive Maintenance** | บำรุงรักษาเชิงป้องกัน | ตามกำหนด |
| **Work Order** | ใบสั่งซ่อม | สั่งซ่อม |
| **Asset Tracking** | ติดตามสินทรัพย์ | ตำแหน่ง |
| **Disposal** | จำหน่าย | ขายทิ้ง |

**ตัวอย่างในฟาร์มเห็ด:**
```
สินทรัพย์:
- โรงเพาะ 5 โรง: 500,000 บาท/โรง
- เครื่องพ่นหมอก 10 เครื่อง: 5,000 บาท/เครื่อง
- พัดลม 20 ตัว: 2,000 บาท/ตัว
- รถขนส่ง 2 คัน: 800,000 บาท/คัน
- ห้องเย็น 1 ห้อง: 300,000 บาท

ค่าเสื่อมราคา:
- โรงเพาะ: 20 ปี
- เครื่องจักร: 5 ปี
- รถ: 10 ปี
```

---

### 2.11 Reporting & Analytics (รายงานและการวิเคราะห์)

| ระบบย่อย | หน้าที่ | ใช้ในฟาร์มเห็ด |
|:---|:---|:---|
| **Financial Reports** | รายงานการเงิน | งบกำไรขาดทุน |
| **Sales Reports** | รายงานการขาย | ยอดขายรายวัน |
| **Inventory Reports** | รายงานสต็อก | สต็อกคงเหลือ |
| **Production Reports** | รายงานการผลิต | ผลผลิต |
| **HR Reports** | รายงาน HR | เงินเดือน |
| **Custom Reports** | รายงานกำหนดเอง | ตามต้องการ |
| **Dashboard** | แดชบอร์ด | ภาพรวม |
| **KPI** | ตัวชี้วัด | วัดผล |

**ตัวอย่างในฟาร์มเห็ด:**
```
รายงานประจำวัน:
- ผลผลิต: 500 กก.
- ยอดขาย: 45,000 บาท
- ต้นทุน: 25,000 บาท
- กำไร: 20,000 บาท
- สต็อกคงเหลือ: 300 กก.
```

---

## 3. โมดูลเสริม (Advanced Modules)

### 3.1 Business Intelligence (BI)

| ระบบย่อย | หน้าที่ |
|:---|:---|
| **Data Warehouse** | คลังข้อมูล |
| **ETL** | ดึง-แปลง-โหลด |
| **OLAP** | วิเคราะห์หลายมิติ |
| **Data Mining** | ขุดข้อมูล |
| **Predictive Analytics** | วิเคราะห์เชิงพยากรณ์ |
| **Dashboard** | แดชบอร์ด |
| **Self-Service BI** | BI ด้วยตนเอง |

### 3.2 E-Commerce Integration

| ระบบย่อย | หน้าที่ |
|:---|:---|
| **Online Store** | ร้านค้าออนไลน์ |
| **Marketplace** | ตลาดออนไลน์ |
| **Payment Gateway** | ชำระเงิน |
| **Shipping Integration** | ขนส่ง |
| **Product Sync** | ซิงค์สินค้า |
| **Order Sync** | ซิงค์ออเดอร์ |

### 3.3 IoT Integration

| ระบบย่อย | หน้าที่ |
|:---|:---|
| **Sensor Data** | ข้อมูลเซ็นเซอร์ |
| **Real-time Monitoring** | ตรวจสอบเรียลไทม์ |
| **Automation** | อัตโนมัติ |
| **Predictive Maintenance** | บำรุงรักษาเชิงพยากรณ์ |
| **Digital Twin** | แบบจำลองเสมือน |

### 3.4 Mobile & Field Service

| ระบบย่อย | หน้าที่ |
|:---|:---|
| **Mobile App** | แอปมือถือ |
| **Field Service** | บริการภาคสนาม |
| **Route Management** | จัดการเส้นทาง |
| **Mobile POS** | ขายหน้างาน |
| **Offline Mode** | ทำงานออฟไลน์ |

### 3.5 Compliance & Governance

| ระบบย่อย | หน้าที่ |
|:---|:---|
| **Regulatory Compliance** | ข้อกำหนดกฎหมาย |
| **Risk Management** | จัดการความเสี่ยง |
| **Internal Audit** | ตรวจสอบภายใน |
| **Document Management** | จัดการเอกสาร |
| **Workflow** | ขั้นตอนการทำงาน |

---

## 4. การเชื่อมโยง ERP กับระบบฟาร์มเห็ด

```
┌─────────────────────────────────────────────────────────────────────┐
│                    MUSHROOM FARM ERP ARCHITECTURE                    │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐        │
│  │     IoT      │────→│  Production  │────→│   Inventory  │        │
│  │  (Sensors)   │     │  (Planning)  │     │   (Stock)    │        │
│  └──────────────┘     └──────────────┘     └──────────────┘        │
│         │                    │                    │                 │
│         ↓                    ↓                    ↓                 │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐        │
│  │     AI       │────→│   Quality    │────→│    Sales     │        │
│  │ (Analytics)  │     │  (Control)   │     │   (Order)    │        │
│  └──────────────┘     └──────────────┘     └──────────────┘        │
│         │                    │                    │                 │
│         ↓                    ↓                    ↓                 │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐        │
│  │   Finance    │←────│    Cost      │←────│     CRM      │        │
│  │ (Accounting) │     │  (Calculate) │     │  (Customer)  │        │
│  └──────────────┘     └──────────────┘     └──────────────┘        │
│         │                    │                    │                 │
│         ↓                    ↓                    ↓                 │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐        │
│  │  Procurement │     │  Supply      │     │   HRM        │        │
│  │ (Purchasing) │     │  Chain       │     │  (Payroll)   │        │
│  └──────────────┘     └──────────────┘     └──────────────┘        │
│         │                    │                    │                 │
│         └────────────────────┼────────────────────┘                 │
│                              ↓                                       │
│                    ┌──────────────┐                                 │
│                    │  Reporting   │                                 │
│                    │  & KPI       │                                 │
│                    └──────────────┘                                 │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 5. ตารางสรุปโมดูล ERP ทั้งหมด

| # | โมดูล | หน้าที่หลัก | ใช้ในฟาร์มเห็ด |
|:---|:---|:---|:---|
| 1 | **Finance & Accounting** | การเงิน บัญชี ต้นทุน | คิดต้นทุนต่อก้อน |
| 2 | **Inventory & Warehouse** | สต็อก คลัง | นับก้อนเชื้อ |
| 3 | **Production & Manufacturing** | การผลิต BOM | สูตรก้อนเชื้อ |
| 4 | **Procurement & Purchasing** | จัดซื้อ | ซื้อขี้เลื่อย |
| 5 | **Sales & Order** | ขาย ออเดอร์ | ใบสั่งขาย |
| 6 | **CRM** | ลูกค้า | ประวัติลูกค้า |
| 7 | **Supply Chain** | ห่วงโซ่อุปทาน | ขนส่ง |
| 8 | **HRM & Payroll** | พนักงาน เงินเดือน | ค่าแรง |
| 9 | **Quality Management** | คุณภาพ | เกรดเห็ด |
| 10 | **Asset Management** | สินทรัพย์ | โรงเพาะ |
| 11 | **Reporting & Analytics** | รายงาน | แดชบอร์ด |
| 12 | **BI** | วิเคราะห์ | พยากรณ์ |
| 13 | **E-Commerce** | ออนไลน์ | Shopee/Lazada |
| 14 | **IoT Integration** | IoT | เซ็นเซอร์ |
| 15 | **Mobile** | มือถือ | แอป |
| 16 | **Compliance** | กฎหมาย | GAP, อย. |

---

## 6. ERP ที่นิยมใช้ (ตัวอย่าง)

### 6.1 ระดับโลก

| ERP | จุดเด่น | เหมาะกับ |
|:---|:---|:---|
| **SAP S/4HANA** | ครบถ้วน, องค์กรใหญ่ | บริษัทขนาดใหญ่ |
| **Oracle NetSuite** | Cloud, ยืดหยุ่น | ธุรกิจทุกขนาด |
| **Microsoft Dynamics 365** | รวมกับ Office | องค์กรทั่วไป |
| **Infor CloudSuite** | เฉพาะอุตสาหกรรม | ผลิต, อาหาร |
| **Odoo** | Open Source, ยืดหยุ่น | SMEs |
| **ERPNext** | Open Source, ฟรี | SMEs |

### 6.2 สำหรับฟาร์มเห็ด

| ERP | เหมาะกับ | ราคา |
|:---|:---|:---|
| **Odoo** | ฟาร์มขนาดกลาง | ฟรี - 50,000 บาท/ปี |
| **ERPNext** | ฟาร์มขนาดเล็ก | ฟรี |
| **NetSuite** | ฟาร์มขนาดใหญ่ | 500,000+ บาท/ปี |
| **SAP B1** | ฟาร์มขนาดกลาง-ใหญ่ | 300,000+ บาท/ปี |

---

## 7. สรุป

**ERP ประกอบด้วย 11 โมดูลหลัก + 5 โมดูลเสริม** รวมทั้งหมด **16 โมดูล** ที่ทำงานร่วมกัน:

1. **การเงินและบัญชี** – หัวใจของ ERP
2. **สินค้าคงคลัง** – ติดตามของ
3. **การผลิต** – ผลิตสินค้า
4. **จัดซื้อ** – ซื้อวัตถุดิบ
5. **การขาย** – ขายสินค้า
6. **CRM** – ดูแลลูกค้า
7. **ห่วงโซ่อุปทาน** – ขนส่ง
8. **HRM** – ดูแลพนักงาน
9. **คุณภาพ** – มาตรฐาน
10. **สินทรัพย์** – ทรัพย์สิน
11. **รายงาน** – ข้อมูล
12. **BI** – วิเคราะห์
13. **E-Commerce** – ออนไลน์
14. **IoT** – เซ็นเซอร์
15. **Mobile** – มือถือ
16. **Compliance** – กฎหมาย

 # การออกแบบ Module ของ CRM + KPI สำหรับ 4 กลุ่มธุรกิจ (การผลิต การบริการ การรับงานภาครัฐ การเกษตร)

## ภาพรวมโครงสร้าง Module

```
┌─────────────────────────────────────────────────────┐
│           CORE MODULES (ใช้ร่วมทุกธุรกิจ)            │
│  ┌──────────┬──────────┬──────────┬──────────┐      │
│  │ Customer │ Contact  │ Activity │ Task &   │      │
│  │ Master   │ Mgmt     │ Log      │ Calendar │      │
│  └──────────┴──────────┴──────────┴──────────┘      │
│  ┌──────────┬──────────┬──────────┬──────────┐      │
│  │ Document │ Report & │ User &   │ Notifica-│      │
│  │ Mgmt     │ Dashboard│ Role     │ tion     │      │
│  └──────────┴──────────┴──────────┴──────────┘      │
├─────────────────────────────────────────────────────┤
│         COMPLIANCE LAYER (GRP/GMP)                   │
│  ┌──────────┬──────────┬──────────┬──────────┐      │
│  │ Audit    │ Batch    │ Quality  │ Approval │      │
│  │ Trail    │ Tracking │ Control  │ Workflow │      │
│  └──────────┴──────────┴──────────┴──────────┘      │
├─────────────────────────────────────────────────────┤
│      INDUSTRY EXTENSION MODULES (เฉพาะธุรกิจ)        │
│  ┌──────────┬──────────┬──────────┬──────────┐      │
│  │Manufactu-│ Service  │Government│Agriculture│     │
│  │ring      │          │Contract  │          │      │
│  └──────────┴──────────┴──────────┴──────────┘      │
├─────────────────────────────────────────────────────┤
│           KPI & ANALYTICS MODULES                    │
│  ┌──────────┬──────────┬──────────┬──────────┐      │
│  │ Sales    │ Marketing│ Service  │ Executive│      │
│  │ KPI      │ KPI      │ KPI      │ KPI      │      │
│  └──────────┴──────────┴──────────┴──────────┘      │
└─────────────────────────────────────────────────────┘
```

---

## ชั้นที่ 1: Core Modules (ใช้ร่วมทุกธุรกิจ)

### 1.1 Customer Master Module (โมดูลข้อมูลหลักลูกค้า)

**วัตถุประสงค์**: เก็บข้อมูลหลักของลูกค้า/หน่วยงานแบบรวมศูนย์ ใช้ร่วมทุกธุรกิจ

**ฟิลด์หลัก**:
| ฟิลด์ | ประเภท | คำอธิบาย |
|---|---|---|
| Customer ID | Auto | รหัสลูกค้าอัตโนมัติ |
| Customer Type | Dropdown | นิติบุคคล/บุคคล/หน่วยงานรัฐ/สหกรณ์การเกษตร |
| Industry Segment | Dropdown | การผลิต/บริการ/ภาครัฐ/การเกษตร |
| Tax ID | Text | เลขประจำตัวผู้เสียภาษี |
| Address | Address | ที่อยู่หลายแห่ง |
| Credit Terms | Number | เงื่อนไขเครดิต |
| Compliance Flag | Boolean | อยู่ภายใต้ GRP/GMP หรือไม่ |
| Status | Dropdown | Active/Inactive/Blacklist |

**KPI ที่เกี่ยวข้อง**: จำนวนลูกค้าใหม่, อัตราการเติบโตของฐานลูกค้า, ลูกค้าที่อยู่ภายใต้ compliance

---

### 1.2 Contact Management Module (โมดูลจัดการผู้ติดต่อ)

**วัตถุประสงค์**: บันทึกผู้ติดต่อหลายคนต่อหนึ่งลูกค้า พร้อมบทบาทและอิทธิพลในการตัดสินใจ

**ฟิลด์หลัก**:
| ฟิลด์ | ประเภท | คำอธิบาย |
|---|---|---|
| Contact ID | Auto | รหัสผู้ติดต่อ |
| Customer Link | Relation | เชื่อมกับ Customer Master |
| Position | Text | ตำแหน่ง |
| Decision Role | Dropdown | ผู้ตัดสินใจ/ผู้มีอิทธิพล/ผู้ใช้งาน/ผู้จัดซื้อ |
| Department | Text | หน่วยงาน |
| Phone/Email | Text | ช่องทางติดต่อ |
| Last Contact Date | Date | ติดต่อครั้งล่าสุด |

**KPI ที่เกี่ยวข้อง**: ความครอบคลุมของผู้ติดต่อ, อัตราการเข้าถึงผู้ตัดสินใจ

---

### 1.3 Activity Log Module (โมดูลบันทึกกิจกรรม)

**วัตถุประสงค์**: บันทึกทุกปฏิสัมพันธ์กับลูกค้า (โทร, เยี่ยม, อีเมล, ประชุม)

**ฟิลด์หลัก**:
| ฟิลด์ | ประเภท | คำอธิบาย |
|---|---|---|
| Activity ID | Auto | รหัสกิจกรรม |
| Activity Type | Dropdown | โทร/เยี่ยม/อีเมล/ประชุม/ส่งเอกสาร |
| Related To | Relation | ลูกค้า/โอกาสขาย/ใบงาน |
| Date & Time | DateTime | วันเวลาที่ทำ |
| Duration | Number | ระยะเวลา (นาที) |
| Outcome | Text | ผลลัพธ์ |
| Next Action | Text | การดำเนินการถัดไป |
| Attachment | File | ไฟล์แนบ |

**KPI ที่เกี่ยวข้อง**: จำนวนกิจกรรมต่อวัน, อัตรากิจกรรม→โอกาสขาย

---

### 1.4 Task & Calendar Module (โมดูลงานและปฏิทิน)

**วัตถุประสงค์**: จัดการงานที่ต้องทำ กำหนดเวลา และติดตามความคืบหน้า

**ฟิลด์หลัก**:
| ฟิลด์ | ประเภท | คำอธิบาย |
|---|---|---|
| Task ID | Auto | รหัสงาน |
| Subject | Text | หัวข้องาน |
| Assigned To | Relation | ผู้รับผิดชอบ |
| Due Date | Date | กำหนดส่ง |
| Priority | Dropdown | สูง/กลาง/ต่ำ |
| Status | Dropdown | ใหม่/กำลังทำ/เสร็จ/เลยกำหนด |
| Reminder | Number | แจ้งเตือนล่วงหน้า (ชั่วโมง) |

**KPI ที่เกี่ยวข้อง**: อัตรางานเสร็จตรงเวลา, งานเลยกำหนด

---

### 1.5 Document Management Module (โมดูลจัดการเอกสาร)

**วัตถุประสงค์**: เก็บเอกสารที่เกี่ยวข้องกับลูกค้า/สัญญา/คุณสมบัติ

**ฟิลด์หลัก**:
| ฟิลด์ | ประเภท | คำอธิบาย |
|---|---|---|
| Document ID | Auto | รหัสเอกสาร |
| Document Type | Dropdown | สัญญา/ใบเสนอราคา/ใบกำกับ/ใบรับรอง |
| Related To | Relation | ลูกค้า/โอกาสขาย/โครงการ |
| Version | Number | เวอร์ชัน |
| Expiry Date | Date | วันหมดอายุ |
| File | File | ไฟล์เอกสาร |

**KPI ที่เกี่ยวข้อง**: เอกสารหมดอายุ, อัตราความครบถ้วนของเอกสาร

---

### 1.6 Report & Dashboard Module (โมดูลรายงานและแดชบอร์ด)

**วัตถุประสงค์**: แสดง KPI ตามบทบาทผู้ใช้

**แดชบอร์ดแยกตามบทบาท**:
- แดชบอร์ดผู้บริหาร: ภาพรวมทุกธุรกิจ
- แดชบอร์ดผู้จัดการขาย: ผลงานทีม
- แดชบอร์ด frontline: งานประจำวัน

---

### 1.7 User & Role Module (โมดูลผู้ใช้และสิทธิ์)

**วัตถุประสงค์**: จัดการผู้ใช้ สิทธิ์การเข้าถึง และการมองเห็นข้อมูล

**ระดับสิทธิ์**:
| ระดับ | สิทธิ์ |
|---|---|
| Admin | ทุกอย่าง |
| Manager | ดูข้อมูลทีม, อนุมัติ |
| User | ดูข้อมูลตัวเอง |
| Auditor | ดูอย่างเดียว, เข้าถึง audit trail |
| Quality | เข้าถึงข้อมูล compliance |

---

### 1.8 Notification Module (โมดูลแจ้งเตือน)

**วัตถุประสงค์**: แจ้งเตือนอัตโนมัติตามเงื่อนไข

**ประเภทการแจ้งเตือน**:
- แจ้งเตือนงานใกล้ครบกำหนด
- แจ้งเตือนวันปิดประมูล
- แจ้งเตือนสัญญาใกล้หมดอายุ
- แจ้งเตือน SLA ใกล้ละเมิด
- แจ้งเตือน compliance ไม่ครบ

---

## ชั้นที่ 2: Compliance Layer (GRP/GMP)

### 2.1 Audit Trail Module (โมดูลบันทึกร่องรอยการตรวจสอบ)

**วัตถุประสงค์**: บันทึกทุกการเปลี่ยนแปลงข้อมูลสำคัญ เพื่อให้ตรวจสอบย้อนกลับได้

**ข้อมูลที่บันทึก**:
| ฟิลด์ | คำอธิบาย |
|---|---|
| Timestamp | วันเวลาที่เปลี่ยนแปลง |
| User | ผู้เปลี่ยนแปลง |
| Field Changed | ฟิลด์ที่เปลี่ยน |
| Old Value | ค่าเดิม |
| New Value | ค่าใหม่ |
| Reason | เหตุผล (บังคับสำหรับฟิลด์ compliance) |
| IP Address | ที่อยู่ IP |

**KPI ที่เกี่ยวข้อง**: จำนวนการเปลี่ยนแปลงที่ไม่ได้รับอนุญาต, ความครบถ้วนของ audit trail

---

### 2.2 Batch Tracking Module (โมดูลติดตาม批次)

**วัตถุประสงค์**: เชื่อมข้อมูลลูกค้า/ข้อร้องเรียนกับ批次การผลิต

**ฟิลด์หลัก**:
| ฟิลด์ | ประเภท | คำอธิบาย |
|---|---|---|
| Batch ID | Text | รหัส批次 |
| Product | Relation | ผลิตภัณฑ์ |
| Production Date | Date | วันที่ผลิต |
| Expiry Date | Date | วันหมดอายุ |
| Release Status | Dropdown | รอตรวจ/ผ่าน/ไม่ผ่าน |
| QC Result | Text | ผลตรวจสอบคุณภาพ |
| Related Customer | Relation | ลูกค้าที่ได้รับ批次นี้ |

**KPI ที่เกี่ยวข้อง**: อัตราปล่อย批次ตรงเวลา, อัตราข้อร้องเรียนต่อ批次

---

### 2.3 Quality Control Module (โมดูลควบคุมคุณภาพ)

**วัตถุประสงค์**: จัดการข้อร้องเรียน, deviation, และ CAPA (Corrective and Preventive Action)

**ฟิลด์หลัก**:
| ฟิลด์ | ประเภท | คำอธิบาย |
|---|---|---|
| Complaint ID | Auto | รหัสข้อร้องเรียน |
| Customer | Relation | ลูกค้าที่ร้องเรียน |
| Batch Link | Relation | 批次ที่เกี่ยวข้อง |
| Severity | Dropdown | วิกฤต/สูง/กลาง/ต่ำ |
| Description | Text | รายละเอียด |
| Investigation | Text | ผลการสอบสวน |
| CAPA | Text | มาตรการแก้ไขและป้องกัน |
| Status | Dropdown | เปิด/กำลังสอบสวน/ปิด |

**KPI ที่เกี่ยวข้อง**: เวลาปิดข้อร้องเรียน, อัตราข้อร้องเรียนซ้ำ, อัตราปิด CAPA ตรงเวลา

---

### 2.4 Approval Workflow Module (โมดูลขั้นตอนอนุมัติ)

**วัตถุประสงค์**: กำหนดขั้นตอนอนุมัติสำหรับรายการที่ต้องควบคุม

**ตัวอย่างขั้นตอนอนุมัติ**:
- ส่วนลดเกินเพดาน → ผู้จัดการ → ผู้อำนวยการ
- ข้อร้องเรียนระดับวิกฤต → ผู้จัดการคุณภาพ → ผู้อำนวยการ
- สัญญาภาครัฐมูลค่าสูง → ฝ่ายกฎหมาย → ผู้บริหาร

**KPI ที่เกี่ยวข้อง**: เวลาอนุมัติเฉลี่ย, อัตราการปฏิเสธ

---

## ชั้นที่ 3: Industry Extension Modules

### 3.1 Manufacturing Module (โมดูลการผลิต)

**3.1.1 Sample-to-Sale Module**
| ฟิลด์ | คำอธิบาย |
|---|---|
| Sample ID | รหัสตัวอย่าง |
| Customer | ลูกค้า |
| Product | ผลิตภัณฑ์ |
| Sample Date | วันที่ส่งตัวอย่าง |
| Test Result | ผลทดสอบ |
| Conversion Status | สถานะแปลงเป็นโอกาสขาย |
| Cost | ต้นทุนตัวอย่าง |

**KPI**: อัตราแปลงตัวอย่าง, รอบเวลาตัวอย่าง→ปิดการขาย, ต้นทุนตัวอย่างต่อรายได้

**3.1.2 Product & Formula Module**
| ฟิลด์ | คำอธิบาย |
|---|---|
| Product ID | รหัสผลิตภัณฑ์ |
| Formula | สูตร |
| GMP Status | สถานะ GMP |
| Registration | ทะเบียน |
| Shelf Life | อายุการเก็บ |

**KPI**: อัตราผลิตภัณฑ์ผ่าน GMP, อายุทะเบียนคงเหลือ

---

### 3.2 Service Module (โมดูลการบริการ)

**3.2.1 Service Ticket Module**
| ฟิลด์ | คำอธิบาย |
|---|---|
| Ticket ID | รหัสใบงาน |
| Customer | ลูกค้า |
| Service Type | ประเภทบริการ |
| Priority | ความสำคัญ |
| Assigned To | ผู้รับผิดชอบ |
| SLA Deadline | กำหนด SLA |
| Status | สถานะ |
| Resolution | วิธีแก้ไข |

**KPI**: เวลาตอบสนองครั้งแรก, เวลาแก้ไข, อัตรา SLA, อัตรางานค้าง

**3.2.2 SLA Monitoring Module**
| ฟิลด์ | คำอธิบาย |
|---|---|
| SLA ID | รหัส SLA |
| Customer | ลูกค้า |
| Service Level | ระดับบริการ |
| Response Time | เวลาตอบสนองที่ตกลง |
| Resolution Time | เวลาแก้ไขที่ตกลง |
| Actual Time | เวลาจริง |
| Compliance | ตรงตาม SLA หรือไม่ |

**KPI**: อัตราบรรลุ SLA, อัตรา escalation

**3.2.3 Renewal & Health Score Module**
| ฟิลด์ | คำอธิบาย |
|---|---|
| Customer | ลูกค้า |
| Contract End | วันสิ้นสุดสัญญา |
| Usage Score | คะแนนการใช้งาน |
| Satisfaction Score | คะแนนความพึงพอใจ |
| Complaint Score | คะแนนข้อร้องเรียน |
| Health Score | คะแนนสุขภาพรวม |
| Renewal Probability | โอกาสต่อสัญญา |

**KPI**: อัตราต่อสัญญา, อัตราสูญเสียลูกค้า, NPS/CSAT

---

### 3.3 Government Contract Module (โมดูลรับงานภาครัฐ)

**3.3.1 Agency Profile Module**
| ฟิลด์ | คำอธิบาย |
|---|---|
| Agency ID | รหัสหน่วยงาน |
| Agency Name | ชื่อหน่วยงาน |
| Level | ระดับ (กระทรวง/กรม/ท้องถิ่น) |
| Qualification Status | สถานะคุณสมบัติ |
| Registration Expiry | วันหมดอายุทะเบียน |
| Key Contacts | ผู้ติดต่อหลัก |

**KPI**: อัตราคุณสมบัติยังใช้ได้, ความถี่การเข้าถึง

**3.3.2 Tender Pipeline Module**
| ฟิลด์ | คำอธิบาย |
|---|---|
| Tender ID | รหัสประมูล |
| Agency | หน่วยงาน |
| Project Name | ชื่อโครงการ |
| Budget | งบประมาณ |
| Announcement Date | วันประกาศ |
| Submission Deadline | วันยื่น |
| Stage | ขั้นตอน |
| Result | ผล |
| Contract Value | มูลค่าสัญญา |

**ขั้นตอน**: รวบรวมข้อมูล → ข้อเสนอ → ยื่นประมูล → รอผล → ปฏิบัติตามสัญญา

**KPI**: จำนวนการเข้าร่วมประมูล, อัตราชนะประมูล, รอบเวลาประมูลเฉลี่ย, อัตราพลาดประมูล

**3.3.3 Personnel Transfer Module**
| ฟิลด์ | คำอธิบาย |
|---|---|
| Agency | หน่วยงาน |
| Previous Owner | ผู้รับผิดชอบเดิม |
| New Owner | ผู้รับผิดชอบใหม่ |
| Transfer Date | วันที่โอน |
| Handover Notes | บันทึกส่งมอบ |
| Relationship Status | สถานะความสัมพันธ์ |

**KPI**: เวลาฟื้นตัวการติดตามหลังเปลี่ยนคน

**3.3.4 Deadline Automation Module**
| ฟิลด์ | คำอธิบาย |
|---|---|
| Tender ID | รหัสประมูล |
| Deadline | วันปิด |
| Reminder Days | แจ้งเตือนล่วงหน้า |
| Responsible | ผู้รับผิดชอบ |
| Alert Status | สถานะแจ้งเตือน |

**KPI**: อัตราพลาดประมูล (เป้าหมาย 0)

---

### 3.4 Agriculture Module (โมดูลการเกษตร)

**3.4.1 Crop & Zone Module**
| ฟิลด์ | คำอธิบาย |
|---|---|
| Zone ID | รหัสพื้นที่ |
| Crop Type | ชนิดพืช |
| Planting Area | พื้นที่เพาะปลูก |
| Farmer Group | กลุ่มเกษตรกร |
| Tech Level | ระดับเทคโนโลยี |
| Planting Cycle | รอบเพาะปลูก |

**KPI**: พื้นที่ครอบคลุม, ผลตอบแทนต่อไร่

**3.4.2 Offline Order Module**
| ฟิลด์ | คำอธิบาย |
|---|---|
| Order ID | รหัสคำสั่งซื้อ |
| Customer | ลูกค้า |
| Product | ผลิตภัณฑ์ |
| Quantity | จำนวน |
| Order Date | วันที่สั่ง |
| Sync Status | สถานะซิงค์ |
| Sync Time | เวลาซิงค์ |

**KPI**: อัตราคำสั่งซื้อ offline, อัตราซิงค์สำเร็จ

**3.4.3 Technical Visit Module**
| ฟิลด์ | คำอธิบาย |
|---|---|
| Visit ID | รหัสการเยี่ยม |
| Customer | ลูกค้า |
| Crop | พืช |
| Visit Date | วันที่เยี่ยม |
| Technical Advice | คำแนะนำเทคนิค |
| Follow-up | การติดตาม |
| Outcome | ผลลัพธ์ |

**KPI**: ความถี่การเยี่ยม, อัตราแปลงจากการเยี่ยม

**3.4.4 Target & Achievement Module**
| ฟิลด์ | คำอธิบาย |
|---|---|
| Target ID | รหัสเป้าหมาย |
| Crop | พืช |
| Branch | สาขา |
| Salesperson | พนักงานขาย |
| Target Amount | เป้าหมาย |
| Actual Amount | ผลจริง |
| Achievement % | % ความสำเร็จ |
| Forecast | พยากรณ์ |

**KPI**: อัตราบรรลุเป้า, ความแม่นยำพยากรณ์ถ่วงน้ำหนัก

---

## ชั้นที่ 4: KPI & Analytics Modules

### 4.1 Sales KPI Module

| KPI | สูตร | ใช้กับ |
|---|---|---|
| อัตราแปลงโอกาสขาย | ปิดการขาย / โอกาสขาย × 100 | ทุกธุรกิจ |
| รอบเวลาขาย | วันปิด - วันสร้างโอกาส | ทุกธุรกิจ |
| มูลค่าสัญญาเฉลี่ย | มูลค่ารวม / จำนวนสัญญา | ทุกธุรกิจ |
| อัตราชนะประมูล | ชนะ / เข้าร่วม × 100 | ภาครัฐ |
| อัตราแปลงตัวอย่าง | ปิดการขาย / ตัวอย่าง × 100 | การผลิต |

### 4.2 Marketing KPI Module

| KPI | สูตร | ใช้กับ |
|---|---|---|
| อัตราการตอบสนอง | ตอบสนอง / ส่ง × 100 | ทุกธุรกิจ |
| ต้นทุนได้มาซึ่งลูกค้า | ต้นทุนการตลาด / ลูกค้าใหม่ | ทุกธุรกิจ |
| อัตราแปลงเป็นลูกค้า | ลูกค้าใหม่ / Leads × 100 | ทุกธุรกิจ |

### 4.3 Service KPI Module

| KPI | สูตร | ใช้กับ |
|---|---|---|
| เวลาตอบสนองครั้งแรก | เวลาตอบ - เวลารับเรื่อง | บริการ |
| อัตราบรรลุ SLA | ตรง SLA / ทั้งหมด × 100 | บริการ |
| CSAT | คะแนนความพึงพอใจเฉลี่ย | บริการ |
| NPS | ผู้แนะนำ - ผู้ไม่แนะนำ | บริการ |
| อัตราต่อสัญญา | ต่อสัญญา / หมดสัญญา × 100 | บริการ |

### 4.4 Executive KPI Module

| KPI | สูตร | ใช้กับ |
|---|---|---|
| CLV | มูลค่าลูกค้าตลอดชีพ | ทุกธุรกิจ |
| อัตรา合规批次 | ผ่าน GMP / ทั้งหมด × 100 | การผลิต |
| อัตราคุณสมบัติครอบคลุม | คุณสมบัติใช้ได้ / ทั้งหมด × 100 | ภาครัฐ |
| อัตราบรรลุพื้นที่ | พื้นที่จริง / เป้า × 100 | การเกษตร |

---

## การแมป Module กับ KPI ตามบทบาท

| บทบาท | Module ที่ใช้ | KPI ที่เห็น |
|---|---|---|
| **ผู้บริหาร** | Executive KPI, Report & Dashboard | CLV, อัตราชนะประมูล, อัตรา合规批次, อัตราต่อสัญญา |
| **ผู้จัดการขาย** | Sales KPI, Tender Pipeline, Target & Achievement | อัตราแปลง, รอบเวลาขาย, อัตราบรรลุเป้า |
| **ผู้จัดการบริการ** | Service KPI, SLA Monitoring, Renewal & Health | SLA, CSAT, อัตราต่อสัญญา |
| **ผู้จัดการคุณภาพ** | Quality Control, Batch Tracking, Audit Trail | เวลาปิดข้อร้องเรียน, อัตรา合规批次, CAPA |
| **frontline ขาย** | Activity Log, Task, Sample-to-Sale | จำนวนกิจกรรม, งานเสร็จตรงเวลา, อัตราแปลงตัวอย่าง |
| **frontline บริการ** | Service Ticket, Task | เวลาตอบสนอง, งานเสร็จ |
| **frontline เกษตร** | Offline Order, Technical Visit | อัตราซิงค์, ความถี่เยี่ยม |
| **Auditor** | Audit Trail, Document Management | ความครบถ้วน audit trail |

---

## สรุปโครงสร้าง Module

| ชั้น | Module | จำนวน |
|---|---|---|
| **Core** | Customer, Contact, Activity, Task, Document, Report, User, Notification | 8 |
| **Compliance** | Audit Trail, Batch Tracking, Quality Control, Approval Workflow | 4 |
| **Manufacturing** | Sample-to-Sale, Product & Formula | 2 |
| **Service** | Service Ticket, SLA Monitoring, Renewal & Health | 3 |
| **Government** | Agency Profile, Tender Pipeline, Personnel Transfer, Deadline Automation | 4 |
| **Agriculture** | Crop & Zone, Offline Order, Technical Visit, Target & Achievement | 4 |
| **KPI** | Sales, Marketing, Service, Executive | 4 |
| **รวม** | | **29 modules** |

**หลักการออกแบบ**: Core + Compliance ใช้ร่วมทุกธุรกิจ ส่วน Industry Extension เลือกเปิดตามธุรกิจที่เกี่ยวข้อง ลูกค้าที่ทำหลายธุรกิจสามารถเปิดหลาย extension พร้อมกันได้ โดยใช้ Customer Master ร่วมกัน