# 📘 สรุปเอกสารภาพรวมระบบ ERP+SaaS (ฉบับเข้าใจง่าย)

## 🎯 ระบบนี้คืออะไร

**ERP+SaaS** คือระบบบริหารจัดการธุรกิจแบบครบวงจรที่ให้บริการผ่านอินเทอร์เน็ต ผู้ใช้ไม่ต้องติดตั้งโปรแกรมเอง ไม่ต้องซื้อเซิร์ฟเวอร์ จ่ายตามการใช้งานจริง เปรียบเสมือน **"ศูนย์กลางธุรกิจ"** ที่รวมทุกอย่างไว้ในที่เดียว ตั้งแต่ขายสินค้า ผลิต จัดส่ง ไปจนถึงบัญชีและภาษี

---

## 🏢 ใครเหมาะกับระบบนี้

| ประเภท | ตัวอย่าง |
|---|---|
| 🏭 โรงงานผลิต | ขนาดกลาง-เล็ก |
| 🛒 ร้านค้าปลีก/ค้าส่ง | ทุกระดับ |
| 🚚 ธุรกิจขนส่ง | โลจิสติกส์ |
| 🌾 เกษตรกร | และเกษตรแปรรูป |
| 🍽️ ร้านอาหาร | และร้านกาแฟ |

---

## 📦 โมดูลหลักของระบบ (1.1 – 1.14)

### กลุ่มบริหารธุรกิจ
| ลำดับ | โมดูล | คำอธิบาย |
|---|---|---|
| 1.1 | Business Canvas | เครื่องมือวางแผนธุรกิจ 9 ช่อง |
| 1.2 | ERP | สินค้าคงคลัง การผลิต การเงิน บุคลากร จัดซื้อ |
| 1.3 | CRM | จัดการลูกค้าตั้งแต่สนใจจนถึงหลังการขาย |
| 1.4 | KPI | ตัวชี้วัด เช่น OEE, OTIF, CAC, LTV |
| 1.6 | Online Marketing | SEO, SEM, โซเชียล, อีเมล, คอนเทนต์ |
| 1.11 | POS | ขายหน้าร้าน ออกใบเสร็จ รับชำระ |
| 1.12 | E-commerce | ขายออนไลน์ เชื่อมระบบหลังบ้าน |
| 1.13 | บัญชีและภาษี | GL, AR, AP, VAT, e-Tax Invoice |

### กลุ่มโลจิสติกส์และติดตาม
| ลำดับ | โมดูล | คำอธิบาย |
|---|---|---|
| 1.5 | Supply Chain | วางแผน จัดหา ผลิต ส่งมอบ |
| 1.7 | Logistics SaaS | จัดการขนส่ง ติดตามพัสดุ |
| 1.8 | IoT | เซนเซอร์ → ส่งข้อมูล → แดชบอร์ด |
| 1.10 | QR Code ติดตามการผลิต | สแกน → บันทึกสถานะเรียลไทม์ |

### กลุ่มวิเคราะห์และเทคโนโลยี
| ลำดับ | โมดูล | คำอธิบาย |
|---|---|---|
| 1.9 | AI พยากรณ์การผลิต | คาดการณ์ความต้องการล่วงหน้า |
| 1.14 | GAP | มาตรฐานเกษตรดี ตรวจสอบย้อนกลับ |

---

## 🔄 ภาพรวมการทำงาน

```
ลูกค้าสั่งซื้อ → ยืนยันออเดอร์ → วางแผนผลิต → พิมพ์ QR ติดงาน
      → พนักงานสแกน QR → ข้อมูลเข้าระบบ → AI วิเคราะห์
      → ผู้บริหารดูแดชบอร์ด
```

ทุกขั้นตอนเชื่อมต่อกันอัตโนมัติ ข้อมูลไหลจากหน้างานถึงผู้บริหารโดยไม่ต้องกรอกซ้ำ

---

## ⭐ จุดเด่นของระบบ

| จุดเด่น | อธิบาย |
|---|---|
| 🎯 Modular | เลือกเปิดเฉพาะโมดูลที่ใช้ |
| 🎯 Scalable | รองรับลูกค้าหลายรายไม่จำกัด |
| 🎯 Secure | ข้อมูลแต่ละบริษัทแยกจากกัน 100% |
| 🎯 Observable | มีระบบติดตามและแจ้งเตือนครบ |
| 🎯 AI-Ready | วิเคราะห์และพยากรณ์อัตโนมัติ |

---

## 🏢 แนวคิด Multi-tenant (อธิบายง่าย)

เปรียบเหมือน **"อพาร์ตเมนต์"**
- อาคารเดียวกัน (ระบบเดียว)
- แต่ละห้อง (แต่ละบริษัท) มีข้อมูลของตัวเอง
- ไม่มีทางเห็นข้อมูลห้องอื่น
- แต่ใช้โครงสร้างพื้นฐานร่วมกัน ช่วยลดต้นทุน

---

## 🔐 ความปลอดภัยและความเป็นส่วนตัว

- 🔐 ข้อมูลแต่ละบริษัทแยกขาดจากกัน
- 🔑 ต้องเข้าสู่ระบบด้วยรหัสผ่านที่เข้ารหัส
- 🛡️ มีการตรวจสอบสิทธิ์ทุกครั้งที่เข้าถึงข้อมูล
- 📝 บันทึกทุกการเปลี่ยนแปลงเพื่อตรวจสอบย้อนหลัง
- 🚦 จำกัดจำนวนการเรียกใช้งาน ป้องกันการโจมตี

---

## 📊 แผนการพัฒนา

| โมดูล | สถานะ |
|---|---|
| ระบบเข้าสู่ระบบและแยกบริษัท | ✅ พร้อมใช้งาน |
| ERP – สินค้าคงคลัง | ✅ พร้อมใช้งาน |
| ERP – การผลิต | ⏳ กำลังพัฒนา |
| ERP – การเงิน | ⏳ กำลังพัฒนา |
| CRM | ⏳ กำลังพัฒนา |
| POS | ⏳ กำลังพัฒนา |
| E-commerce | ⏳ กำลังพัฒนา |
| Logistics | ⏳ กำลังพัฒนา |
| IoT / QR | ⏳ กำลังพัฒนา |
| AI Forecast | ⏳ กำลังพัฒนา |
| รายงาน / KPI | ⏳ กำลังพัฒนา |

---

## 🏭 ตัวอย่างความสำเร็จ: โรงงาน ABC

| ก่อนใช้ระบบ | หลังใช้ระบบ 6 เดือน |
|---|---|
| ผลิตเกินความต้องการ 30% | ✅ ของเสียลดจาก 15% → 5% |
| ของเสีย 15% | ✅ ส่งมอบตรงเวลาเพิ่มจาก 70% → 95% |
| ส่งมอบตรงเวลาเพียง 70% | ✅ กำไรเพิ่ม 25% |

---

## ✅ ประโยชน์ที่ธุรกิจจะได้รับ

- ✅ ลดต้นทุนด้านซอฟต์แวร์และฮาร์ดแวร์
- ✅ ขยายธุรกิจได้รวดเร็ว
- ✅ ข้อมูลรวมศูนย์ วิเคราะห์ได้ทันที
- ✅ AI ช่วยพยากรณ์แม่นยำ
- ✅ ตัดสินใจบนข้อมูลจริง ไม่เดา

---

## ⚠️ ข้อควรระวัง

- ⚠️ ต้องแยกข้อมูลลูกค้าแต่ละรายอย่างเข้มงวด
- ⚠️ ต้องออกแบบระบบรองรับการเติบโต
- ⚠️ ค่าใช้จ่ายคลาวด์ต้องบริหารให้ดี

---

## 📋 ข้อดีและข้อเสีย

**ข้อดี**
- Modular, ขยายได้, คุ้มค่า

**ข้อเสีย**
- ซับซ้อนในการตรวจสอบปัญหาข้ามบริษัท
- ต้องมีทีมดูแลระบบที่แข็งแรง

**ข้อห้าม**
- ❌ ห้ามเข้าถึงข้อมูลโดยไม่ระบุว่าของบริษัทไหน
- ❌ ห้ามบันทึกข้อมูลอ่อนไหวลงใน log
- ❌ ห้ามเก็บรหัสลับไว้ในโค้ด

---

## 🚀 คู่มือการใช้งาน (ภาพรวม)

1. สมัครใช้งาน → ได้ที่อยู่ระบบเฉพาะบริษัท
2. เข้าสู่ระบบ → รับสิทธิ์การใช้งาน
3. ตั้งค่าผู้ใช้และสิทธิ์
4. เปิดโมดูลที่ต้องการ
5. นำเข้าข้อมูลเริ่มต้น
6. เริ่มใช้งานได้ทันที

---

## 🔧 การดูแลรักษาระบบ

| รอบ | กิจกรรม |
|---|---|
| รายวัน | ตรวจสอบ log และแจ้งเตือน |
| รายสัปดาห์ | ตรวจสอบการสำรองข้อมูล |
| รายเดือน | อัปเดตความปลอดภัย |
| รายไตรมาส | ทดสอบแผนกู้คืนระบบ |

---

## 📈 การขยายและแก้ไขระบบ

- เพิ่มโมดูลใหม่ → สร้างส่วนงานใหม่แยกออกมา
- เพิ่มฟิลด์ข้อมูล → อัปเดตแบบไม่กระทบของเดิม
- เพิ่มลูกค้าใหม่ → ไม่ต้องแก้โค้ด (รองรับหลายบริษัทโดยออกแบบ)

---

## 👥 กระบวนการทำงานร่วมกัน

### การแบ่งสายงาน
- `main` → ระบบจริงสำหรับลูกค้า
- `develop` → ระบบทดสอบ
- `feature/*` → งานที่กำลังพัฒนา
- `hotfix/*` → แก้ไขเร่งด่วน

### การตรวจสอบก่อนรวมงาน
- โค้ดผ่านการตรวจสอบคุณภาพ
- มีการทดสอบครอบคลุม
- ไม่มีรหัสลับหลุดในโค้ด
- มีความคิดเห็นสองภาษา
- อัปเดตเอกสารประกอบ
- มีผู้ตรวจสอบอย่างน้อย 2 คน

---

## 🚢 การส่งมอบและติดตั้งอัตโนมัติ

- ทุกครั้งที่มีการรวมโค้ด → ระบบทดสอบอัตโนมัติ
- ผ่านการทดสอบ → ส่งขึ้นระบบจริง
- ใช้วิธี Blue-Green → เปลี่ยนเวอร์ชันโดยไม่หยุดให้บริการ

---

## 🔍 การวิเคราะห์สาเหตุของปัญหา (RCA)

### ขั้นตอน
1. **ระบุปัญหา** – เช่น "ระบบ POS ล่มทุกวันเวลา 14:00 น."
2. **รวบรวมข้อมูล** – log, รายงานผู้ใช้, ไทม์ไลน์
3. **ระบุสาเหตุที่เป็นไปได้** – ฐานข้อมูลเต็ม? หน่วยความจำรั่ว?
4. **หาสาเหตุหลัก** – ใช้คำถาม "ทำไม" 5 ครั้ง
5. **วางแผนแก้ไข** – แก้ที่ต้นเหตุ ไม่ใช่แค่อาการ
6. **ติดตามผล** – ดูแดชบอร์ดและทบทวนรายสัปดาห์

### ตัวอย่าง 5 Whys
1. ทำไม POS ล่ม? → ฐานข้อมูล timeout
2. ทำไม timeout? → connection pool เต็ม
3. ทำไมเต็ม? → query ค้าง
4. ทำไมค้าง? → ขาด index
5. ทำไมขาด index? → **ไม่มีกระบวนการตรวจสอบ migration** ← สาเหตุหลัก

---

## 🏗️ โครงสร้าง Modules Golang (ตัวอย่าง)

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

---

## 🎯 บทสรุป

**ERP+SaaS** คือแพลตฟอร์มที่ช่วยธุรกิจ SME
- ลดต้นทุน
- เพิ่มประสิทธิภาพ
- ตัดสินใจบนข้อมูลจริง
- พร้อม AI ช่วยพยากรณ์
- ขยายธุรกิจได้ไม่จำกัด

เอกสารนี้เป็น **ภาพรวมสำหรับผู้บริหารและทีมธุรกิจ** หากต้องการรายละเอียดทางเทคนิค สามารถดูได้จากเอกสารต้นแบบฉบับเต็ม

---

## 📌 วิธีใช้เอกสารนี้

```
"จากเอกสารภาพรวมนี้ จงอธิบายโมดูล [ชื่อโมดูล] 
ให้ผู้บริหารเข้าใจง่าย ภายใน 1 หน้า 
พร้อมประโยชน์ทางธุรกิจและตัวอย่างการใช้งาน"
```

---

**เวอร์ชัน**: 1.0
**ประเภทเอกสาร**: ภาพรวมสำหรับผู้บริหาร ทีมธุรกิจ และผู้ใช้งาน
**ไม่ลงรายละเอียดทางเทคนิค**


# 🏗️ ออกแบบ Modules ERP+SaaS ตามโครงสร้าง Golang

## 📐 หลักการออกแบบ (Design Principles)

| หลักการ | คำอธิบาย |
|---|---|
| **Clean Architecture** | แยกชั้น domain → application → infrastructure → interfaces |
| **DDD (Domain-Driven Design)** | แต่ละโมดูลมี Aggregate Root, Value Object, Repository |
| **Multi-tenant** | ทุกโมดูลต้องมี `tenant_id` ในทุก layer |
| **Modular Monolith** | แยกโมดูลชัดเจน แต่ deploy รวมกันได้ |
| **Event-Driven** | สื่อสารระหว่างโมดูลผ่าน Kafka |

---

## 1️⃣ ERP – สินค้าคงคลัง (Inventory) ✅ พร้อมใช้งาน

```
internal/modules/inventory/
├── domain/
│   ├── entity/
│   │   ├── product.go              # Aggregate Root - สินค้า
│   │   ├── warehouse.go            # คลังสินค้า
│   │   ├── stock_movement.go       # การเคลื่อนไหวสต็อก
│   │   ├── stock_balance.go        # ยอดคงเหลือ
│   │   └── unit_of_measure.go      # หน่วยนับ
│   ├── value_object/
│   │   ├── sku.go                  # รหัสสินค้า + validation
│   │   ├── quantity.go             # จำนวน + หน่วย
│   │   ├── stock_status.go         # IN_STOCK | LOW | OUT_OF_STOCK
│   │   └── movement_type.go        # IN | OUT | TRANSFER | ADJUST
│   ├── repository/
│   │   ├── product_repository.go
│   │   ├── warehouse_repository.go
│   │   └── stock_repository.go
│   ├── service/
│   │   └── inventory_domain_service.go  # ตรวจสต็อก, จองสินค้า
│   └── errors/errors.go
├── application/
│   ├── create_product.go
│   ├── adjust_stock.go
│   ├── transfer_stock.go
│   ├── reserve_stock.go            # จองสินค้าสำหรับออเดอร์
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   │   ├── product_repo_impl.go
│   │   ├── stock_repo_impl.go
│   │   └── models.go
│   ├── persistence/redis/stock_cache.go
│   └── messaging/kafka_producer.go # stock.updated, stock.low
├── interfaces/
│   ├── http/
│   │   ├── product_handler.go
│   │   ├── stock_handler.go
│   │   └── routes.go
│   └── middleware/tenant.go
└── module.go
```

---

## 2️⃣ ERP – การผลิต (Production) ⏳ กำลังพัฒนา

```
internal/modules/production/
├── domain/
│   ├── entity/
│   │   ├── production_order.go     # Aggregate Root - ใบสั่งผลิต
│   │   ├── bom.go                  # Bill of Materials
│   │   ├── work_center.go          # ศูนย์ผลิต
│   │   ├── routing.go              # ลำดับขั้นตอนผลิต
│   │   └── quality_check.go        # ตรวจสอบคุณภาพ
│   ├── value_object/
│   │   ├── order_status.go         # DRAFT | RELEASED | IN_PROGRESS | DONE
│   │   ├── production_qty.go       # จำนวนผลิต
│   │   └── oee.go                  # Overall Equipment Effectiveness
│   ├── repository/
│   │   ├── production_order_repository.go
│   │   ├── bom_repository.go
│   │   └── work_center_repository.go
│   ├── service/
│   │   └── production_domain_service.go  # วางแผนผลิต, คำนวณวัสดุ
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
│   └── messaging/kafka_producer.go # production.started, production.done
├── interfaces/
│   ├── http/
│   └── middleware/tenant.go
└── module.go
```

---

## 3️⃣ ERP – การเงิน (Finance) ⏳ กำลังพัฒนา

```
internal/modules/finance/
├── domain/
│   ├── entity/
│   │   ├── invoice.go              # Aggregate Root - ใบแจ้งหนี้
│   │   ├── payment.go              # การชำระเงิน
│   │   ├── journal_entry.go        # บันทึกบัญชี
│   │   ├── account.go              # ผังบัญชี
│   │   └── tax_invoice.go          # e-Tax Invoice
│   ├── value_object/
│   │   ├── money.go                # จำนวนเงิน + สกุล
│   │   ├── invoice_status.go       # DRAFT | ISSUED | PAID | OVERDUE
│   │   ├── tax_id.go               # เลขผู้เสียภาษี
│   │   └── vat_rate.go             # อัตรา VAT 7%
│   ├── repository/
│   │   ├── invoice_repository.go
│   │   ├── payment_repository.go
│   │   └── journal_repository.go
│   ├── service/
│   │   └── finance_domain_service.go  # คำนวณ VAT, GL posting
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
│   └── messaging/kafka_producer.go # invoice.issued, payment.received
├── interfaces/
│   ├── http/
│   └── middleware/tenant.go
└── module.go
```

---

## 4️⃣ CRM ⏳ กำลังพัฒนา

```
internal/modules/crm/
├── domain/
│   ├── entity/
│   │   ├── lead.go                 # Aggregate Root - ผู้สนใจ
│   │   ├── opportunity.go          # โอกาสขาย
│   │   ├── activity.go             # กิจกรรม (โทร, นัด, อีเมล)
│   │   ├── campaign.go             # แคมเปญการตลาด
│   │   └── customer_360.go         # มุมมองลูกค้ารอบด้าน
│   ├── value_object/
│   │   ├── lead_status.go          # NEW | CONTACTED | QUALIFIED | LOST
│   │   ├── opportunity_stage.go    # PROSPECTING | PROPOSAL | NEGOTIATION | WON
│   │   ├── source.go               # WEB | PHONE | REFERRAL | EVENT
│   │   └── score.go                # Lead scoring
│   ├── repository/
│   │   ├── lead_repository.go
│   │   ├── opportunity_repository.go
│   │   └── activity_repository.go
│   ├── service/
│   │   └── crm_domain_service.go   # Lead scoring, pipeline
│   └── errors/errors.go
├── application/
│   ├── create_lead.go
│   ├── convert_lead.go             # แปลงเป็น customer
│   ├── create_opportunity.go
│   ├── log_activity.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/redis/crm_cache.go
│   ├── search/elasticsearch/lead_indexer.go
│   └── messaging/kafka_producer.go # lead.created, deal.won
├── interfaces/
│   ├── http/
│   └── middleware/tenant.go
└── module.go
```

---

## 5️⃣ POS ⏳ กำลังพัฒนา

```
internal/modules/pos/
├── domain/
│   ├── entity/
│   │   ├── sale.go                 # Aggregate Root - การขาย
│   │   ├── terminal.go             # เครื่อง POS
│   │   ├── shift.go                # กะการทำงาน
│   │   ├── receipt.go              # ใบเสร็จ
│   │   └── payment_method.go       # วิธีชำระ (เงินสด/บัตร/QR)
│   ├── value_object/
│   │   ├── sale_status.go          # OPEN | COMPLETED | VOIDED | REFUNDED
│   │   ├── payment_type.go         # CASH | CARD | QR | CREDIT
│   │   └── discount.go             # ส่วนลด
│   ├── repository/
│   │   ├── sale_repository.go
│   │   ├── terminal_repository.go
│   │   └── shift_repository.go
│   ├── service/
│   │   └── pos_domain_service.go   # คำนวณเงินทอน, ปิดกะ
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
│   └── messaging/kafka_producer.go # sale.completed, shift.closed
├── interfaces/
│   ├── http/
│   └── middleware/tenant.go
└── module.go
```

---

## 6️⃣ E-commerce ⏳ กำลังพัฒนา

```
internal/modules/ecommerce/
├── domain/
│   ├── entity/
│   │   ├── cart.go                 # Aggregate Root - ตะกร้า
│   │   ├── order.go                # คำสั่งซื้อออนไลน์
│   │   ├── product_listing.go      # สินค้าบนเว็บ
│   │   ├── shipping.go             # การจัดส่ง
│   │   └── promotion.go            # โปรโมชัน
│   ├── value_object/
│   │   ├── order_status.go         # PENDING | PAID | SHIPPED | DELIVERED
│   │   ├── cart_item.go            # รายการในตะกร้า
│   │   └── channel.go              # WEB | MOBILE | MARKETPLACE
│   ├── repository/
│   │   ├── cart_repository.go
│   │   ├── order_repository.go
│   │   └── listing_repository.go
│   ├── service/
│   │   └── ecommerce_domain_service.go  # คำนวณราคา, โปรโมชัน
│   └── errors/errors.go
├── application/
│   ├── add_to_cart.go
│   ├── checkout.go
│   ├── confirm_payment.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/redis/cart_cache.go
│   └── messaging/kafka_producer.go # order.placed, order.paid
├── interfaces/
│   ├── http/
│   └── middleware/tenant.go
└── module.go
```

---

## 7️⃣ Logistics ⏳ กำลังพัฒนา

```
internal/modules/logistics/
├── domain/
│   ├── entity/
│   │   ├── shipment.go             # Aggregate Root - การจัดส่ง
│   │   ├── carrier.go              # ผู้ให้บริการขนส่ง
│   │   ├── route.go                # เส้นทาง
│   │   ├── tracking_event.go       # เหตุการณ์ติดตาม
│   │   └── delivery_proof.go       # หลักฐานการส่งมอบ
│   ├── value_object/
│   │   ├── shipment_status.go      # PICKED | IN_TRANSIT | DELIVERED | FAILED
│   │   ├── address.go              # ที่อยู่จัดส่ง
│   │   └── tracking_number.go      # เลขติดตาม
│   ├── repository/
│   │   ├── shipment_repository.go
│   │   ├── carrier_repository.go
│   │   └── tracking_repository.go
│   ├── service/
│   │   └── logistics_domain_service.go  # เลือก carrier, คำนวณเส้นทาง
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
│   └── messaging/kafka_producer.go # shipment.created, shipment.delivered
├── interfaces/
│   ├── http/
│   └── middleware/tenant.go
└── module.go
```

---

## 8️⃣ IoT / QR ⏳ กำลังพัฒนา

```
internal/modules/iot/
├── domain/
│   ├── entity/
│   │   ├── device.go               # Aggregate Root - อุปกรณ์
│   │   ├── sensor_reading.go       # ค่าที่อ่านได้
│   │   ├── qr_code.go              # QR Code ติดตามงาน
│   │   ├── scan_event.go           # เหตุการณ์สแกน
│   │   └── alert_rule.go           # กฎแจ้งเตือน
│   ├── value_object/
│   │   ├── device_status.go        # ONLINE | OFFLINE | ERROR
│   │   ├── sensor_type.go          # TEMP | HUMIDITY | PRESSURE | GPS
│   │   ├── qr_payload.go           # ข้อมูลใน QR
│   │   └── alert_level.go          # INFO | WARNING | CRITICAL
│   ├── repository/
│   │   ├── device_repository.go
│   │   ├── reading_repository.go
│   │   └── qr_repository.go
│   ├── service/
│   │   └── iot_domain_service.go   # ตรวจ threshold, สร้าง QR
│   └── errors/errors.go
├── application/
│   ├── register_device.go
│   ├── ingest_reading.go           # รับค่าจาก sensor
│   ├── generate_qr.go
│   ├── scan_qr.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/timeseries/     # Time-series DB สำหรับ readings
│   ├── persistence/redis/device_cache.go
│   ├── messaging/mqtt/consumer.go  # MQTT สำหรับ IoT
│   └── messaging/kafka_producer.go # device.alert, qr.scanned
├── interfaces/
│   ├── http/
│   ├── websocket/realtime.go       # แจ้งเตือนเรียลไทม์
│   └── middleware/tenant.go
└── module.go
```

---

## 9️⃣ AI Forecast ⏳ กำลังพัฒนา

```
internal/modules/forecast/
├── domain/
│   ├── entity/
│   │   ├── forecast_model.go       # Aggregate Root - โมเดล
│   │   ├── prediction.go           # ผลพยากรณ์
│   │   ├── training_data.go        # ข้อมูลฝึก
│   │   └── model_metric.go         # ตัวชี้วัดโมเดล
│   ├── value_object/
│   │   ├── model_type.go           # ARIMA | PROPHET | LSTM | XGBOOST
│   │   ├── forecast_horizon.go     # 7d | 30d | 90d
│   │   ├── confidence.go           # ความเชื่อมั่น
│   │   └── accuracy.go             # MAPE, RMSE
│   ├── repository/
│   │   ├── model_repository.go
│   │   └── prediction_repository.go
│   ├── service/
│   │   └── forecast_domain_service.go  # เลือกโมเดล, ประเมินผล
│   └── errors/errors.go
├── application/
│   ├── train_model.go
│   ├── predict_demand.go
│   ├── retrain_model.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/timeseries/features.go
│   ├── ml/python_bridge.go         # เรียก Python ML service
│   └── messaging/kafka_producer.go # forecast.generated
├── interfaces/
│   ├── http/
│   └── middleware/tenant.go
└── module.go
```

---

## 🔟 รายงาน / KPI ⏳ กำลังพัฒนา

```
internal/modules/reporting/
├── domain/
│   ├── entity/
│   │   ├── report.go               # Aggregate Root - รายงาน
│   │   ├── dashboard.go            # แดชบอร์ด
│   │   ├── kpi_definition.go       # นิยาม KPI
│   │   ├── kpi_value.go            # ค่า KPI
│   │   └── schedule.go             # ตารางสร้างรายงาน
│   ├── value_object/
│   │   ├── report_type.go          # SALES | INVENTORY | PRODUCTION | FINANCE
│   │   ├── kpi_code.go             # OEE | OTIF | CAC | LTV
│   │   ├── period.go               # DAY | WEEK | MONTH | QUARTER | YEAR
│   │   └── export_format.go        # PDF | EXCEL | CSV | JSON
│   ├── repository/
│   │   ├── report_repository.go
│   │   ├── dashboard_repository.go
│   │   └── kpi_repository.go
│   ├── service/
│   │   └── reporting_domain_service.go  # คำนวณ KPI, aggregate data
│   └── errors/errors.go
├── application/
│   ├── generate_report.go
│   ├── create_dashboard.go
│   ├── calculate_kpi.go
│   ├── schedule_report.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/clickhouse/     # OLAP สำหรับ analytics
│   ├── persistence/redis/report_cache.go
│   └── messaging/kafka_consumer.go # รับ event จากโมดูลอื่น
├── interfaces/
│   ├── http/
│   └── middleware/tenant.go
└── module.go
```

---

## 📊 สรุปภาพรวมโมดูลทั้งหมด

| # | โมดูล | Aggregate Root | สถานะ |
|---|---|---|---|
| 1 | Inventory | Product, Warehouse | ✅ พร้อม |
| 2 | Production | ProductionOrder, BOM | ⏳ |
| 3 | Finance | Invoice, Payment | ⏳ |
| 4 | CRM | Lead, Opportunity | ⏳ |
| 5 | POS | Sale, Shift | ⏳ |
| 6 | E-commerce | Cart, Order | ⏳ |
| 7 | Logistics | Shipment, Carrier | ⏳ |
| 8 | IoT/QR | Device, QRCode | ⏳ |
| 9 | Forecast | ForecastModel | ⏳ |
| 10 | Reporting | Report, Dashboard | ⏳ |

---

## 🔗 การสื่อสารระหว่างโมดูล (Event-Driven)

```
┌─────────────┐     order.placed      ┌─────────────┐
│ E-commerce  │ ─────────────────────▶│  Inventory  │
└─────────────┘                       └─────────────┘
                                             │
                                             │ stock.reserved
                                             ▼
┌─────────────┐     production.done   ┌─────────────┐
│  Production │ ─────────────────────▶│  Inventory  │
└─────────────┘                       └─────────────┘
                                             │
                                             │ stock.updated
                                             ▼
┌─────────────┐     invoice.issued    ┌─────────────┐
│   Finance   │ ◀─────────────────────│     POS     │
└─────────────┘                       └─────────────┘
       │
       │ payment.received
       ▼
┌─────────────┐     kpi.updated       ┌─────────────┐
│  Reporting  │ ◀─────────────────────│  ทุกโมดูล    │
└─────────────┘                       └─────────────┘
```

---

## 🎯 หลักการสำคัญที่ทุกโมดูลต้องมี

| หลักการ | รายละเอียด |
|---|---|
| **tenant_id** | ทุก entity, repository, handler ต้องมี |
| **domain/errors** | แต่ละโมดูลมี error ของตัวเอง |
| **module.go** | จุด entry point สำหรับ wire dependency |
| **dto.go** | แยก request/response ชัดเจน |
| **middleware/tenant.go** | inject tenant จาก JWT |
| **kafka_producer** | ส่ง event ออกเมื่อมีการเปลี่ยนแปลง |
| **repository interface** | อยู่ใน domain, implementation อยู่ใน infrastructure |

---

**เวอร์ชัน**: 1.0
**ประเภทเอกสาร**: โครงสร้างโมดูลสำหรับทีมพัฒนา
**อ้างอิง**: เอกสารภาพรวมระบบ ERP+SaaS