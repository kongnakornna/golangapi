# 📄 เอกสารโครงการระบบบริหารจัดการอู่ซ่อมรถ ฉบับสมบูรณ์  
**Auto Repair Shop Management System – Complete Project Documentation**

---

| รายการ | รายละเอียด |
|--------|-----------|
| **ชื่อโครงการ** | ระบบบริหารจัดการอู่ซ่อมรถ (Auto Repair Shop Management System) |
| **ผู้เขียน** | Kongnakorn Jantakun |
| **วันที่** | 2026-07-04 |
| **เวอร์ชัน** | 3.0 (Final) |
| **สถานะ** | ฉบับสมบูรณ์ (Complete) |

---

## 📋 สารบัญ (Table of Contents)

1. [บทนำและภาพรวมระบบ](#1-บทนำและภาพรวมระบบ)
2. [สถาปัตยกรรมระบบ](#2-สถาปัตยกรรมระบบ)
3. [โมดูลที่ 1: Authentication & Permission](#3-โมดูลที่-1-authentication--permission)
4. [โมดูลที่ 2: Job Card Management](#4-โมดูลที่-2-job-card-management)
5. [โมดูลที่ 3: Customer Management](#5-โมดูลที่-3-customer-management)
6. [โมดูลที่ 4: Quotation](#6-โมดูลที่-4-quotation)
7. [โมดูลที่ 5: Purchase Order](#7-โมดูลที่-5-purchase-order)
8. [โมดูลที่ 6: Inventory Management](#8-โมดูลที่-6-inventory-management)
9. [โมดูลที่ 7: Payment Management](#9-โมดูลที่-7-payment-management)
10. [โมดูลที่ 8: Dashboard & Reports](#10-โมดูลที่-8-dashboard--reports)
11. [โมดูลที่ 9: Document Management](#11-โมดูลที่-9-document-management)
12. [โมดูลที่ 10: Email Service](#12-โมดูลที่-10-email-service)
13. [โมดูลที่ 11: Batch Jobs](#13-โมดูลที่-11-batch-jobs)
14. [โมดูลที่ 12: Multi-Language (i18n)](#14-โมดูลที่-12-multi-language-i18n)
15. [โมดูลที่ 13: IoT & GPS Tracking](#15-โมดูลที่-13-iot--gps-tracking)
16. [โมดูลที่ 14: Web Order System (WOS)](#16-โมดูลที่-14-web-order-system-wos)
17. [สรุปฐานข้อมูลทั้งหมด (Consolidated Database)](#17-สรุปฐานข้อมูลทั้งหมด-consolidated-database)
18. [สรุป API ทั้งหมด (Consolidated API)](#18-สรุป-api-ทั้งหมด-consolidated-api)
19. [สรุป Redis Cache Keys](#19-สรุป-redis-cache-keys)
20. [สรุป Rate Limit Policy](#20-สรุป-rate-limit-policy)
21. [ภาคผนวก: การติดตั้งและการใช้งาน](#21-ภาคผนวก-การติดตั้งและการใช้งาน)

---

## 1. บทนำและภาพรวมระบบ

### 1.1 วัตถุประสงค์
ระบบบริหารจัดการอู่ซ่อมรถ (Auto Repair Shop Management System) ถูกพัฒนาขึ้นเพื่อเพิ่มประสิทธิภาพในการดำเนินงานของศูนย์บริการหรืออู่ซ่อมรถ ครอบคลุมตั้งแต่การรับรถเข้าซ่อม การวินิจฉัยปัญหา การเสนอราคา การสั่งซื้ออะไหล่ การเบิกจ่ายสินค้าคงคลัง การออกใบแจ้งหนี้ และการติดตามประวัติการซ่อมบำรุงของลูกค้าและยานพาหนะ

ระบบถูกออกแบบให้มีความยืดหยุ่น รองรับการขยายตัวในอนาคต โดยใช้สถาปัตยกรรมแบบ **Domain‑Driven Design (DDD)** ผสานกับ **Clean Architecture** และ **Event‑Driven** เพื่อแยกความรับผิดชอบและเพิ่มความสามารถในการบำรุงรักษา

### 1.2 เทคโนโลยีหลัก (Tech Stack)

| หมวดหมู่ | เทคโนโลยี | เวอร์ชัน |
|---------|-----------|----------|
| **ภาษา** | Java | 17+ / 21 |
| **Framework** | Spring Boot | 3.4.1 |
| **ORM** | Spring Data JPA (Hibernate) | - |
| **ฐานข้อมูลหลัก** | PostgreSQL | 15+ |
| **Cache** | Redis | 7+ |
| **Message Queue** | Apache Kafka | 3.4+ |
| **Logging & Monitoring** | ELK Stack, Grafana, Micrometer | - |
| **การจัดการเอกสาร** | JasperReports (PDF), Apache POI (Excel) | - |
| **IoT** | MQTT (EMQX), InfluxDB | - |
| **Documentation** | Swagger/OpenAPI 3.0 | - |
| **Build Tool** | Maven | 3.8+ |

### 1.3 กลุ่มผู้ใช้งาน

| กลุ่มผู้ใช้งาน | บทบาทหน้าที่ |
|--------------|-------------|
| **พนักงานหน้าร้าน (Service Advisor)** | รับรถ, สร้าง Job Card, ออก Quotation, ปิดงาน |
| **ช่างเทคนิค (Mechanic)** | วินิจฉัย, ซ่อมแซม, อัปเดตสถานะงาน |
| **พนักงานคลังสินค้า (Store Keeper)** | จัดการสินค้าคงคลัง, เบิกจ่ายอะไหล่, รับสินค้า |
| **ฝ่ายจัดซื้อ (Purchasing)** | สร้าง Purchase Order, ติดตาม Supplier |
| **ฝ่ายบัญชี/การเงิน (Finance)** | ออก Invoice, รับชำระเงิน, จัดการเอกสารปรับปรุง |
| **ผู้ดูแลระบบ (Admin)** | จัดการผู้ใช้, สิทธิ์การใช้งาน, ข้อมูลพื้นฐาน |
| **ผู้บริหาร (Executive)** | ดู Dashboard และรายงานสรุป |

---

## 2. สถาปัตยกรรมระบบ

### 2.1 แผนภาพสถาปัตยกรรม (Layered + Event‑Driven)

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              EXTERNAL SYSTEMS                                    │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                             API GATEWAY / LOAD BALANCER                          │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                             SPRING BOOT APPLICATION                              │
│  ┌───────────────────────────────────────────────────────────────────────────┐  │
│  │                           CONTROLLER LAYER                                 │  │
│  └───────────────────────────────────────────────────────────────────────────┘  │
│                                      │                                           │
│  ┌───────────────────────────────────▼───────────────────────────────────────┐  │
│  │                             SERVICE LAYER                                  │  │
│  └───────────────────────────────────┬───────────────────────────────────────┘  │
│                                      │                                           │
│  ┌───────────────────────────────────▼───────────────────────────────────────┐  │
│  │                            REPOSITORY LAYER                                │  │
│  └───────────────────────────────────┬───────────────────────────────────────┘  │
└────────────────────────────────────┼──────────────────────────────────────────┘
                                     │
                 ┌───────────────────┼───────────────────┐
                 ▼                   ▼                   ▼
        ┌──────────────┐    ┌──────────────┐    ┌──────────────┐
        │  PostgreSQL  │    │    Redis     │    │    Kafka     │
        │  (Primary DB)│    │   (Cache)    │    │ (Event Bus)  │
        └──────────────┘    └──────────────┘    └──────┬───────┘
                                                        │
        ┌───────────────────────────────────────────────┼───────────────────┐
        ▼                                               ▼                   ▼
 ┌─────────────┐                                ┌─────────────┐      ┌─────────────┐
 │ Elasticsearch│                                │  InfluxDB   │      │   Grafana   │
 │   (ELK)     │                                │  (IoT Data) │      │  (Metrics)  │
 └─────────────┘                                └─────────────┘      └─────────────┘
```

### 2.2 หลักการออกแบบ

1. **Separation of Concerns**: แต่ละชั้นมีหน้าที่เฉพาะ
2. **Dependency Inversion**: Domain ไม่ขึ้นกับ Infrastructure
3. **Code Reuse**: Generic Structures ลดการเขียนโค้ดซ้ำ
4. **Testability**: สถาปัตยกรรมเอื้อต่อการสร้าง Unit Test

---

## 3. โมดูลที่ 1: Authentication & Permission

### 3.1 Database Tables

| ตาราง | คำอธิบาย |
|-------|----------|
| `m_user` | ผู้ใช้ระบบ (Username, Email, PasswordHash, Status) |
| `m_role` | บทบาท (ADMIN, MANAGER, USER) |
| `m_permission` | สิทธิ์การใช้งาน (JOB_CREATE, INVENTORY_READ, etc.) |
| `m_user_role` | เชื่อม User-Role |
| `m_role_permission` | เชื่อม Role-Permission |
| `m_user_token` | Token ที่ออกให้ (JWT, Refresh) |

### 3.2 API Endpoints

| Method | Path | คำอธิบาย | Rate Limit |
|--------|------|----------|------------|
| POST | `/auth/login` | เข้าสู่ระบบ | 5/300s |
| POST | `/auth/logout` | ออกจากระบบ | 10/60s |
| POST | `/auth/refresh` | ต่ออายุ Token | 20/3600s |
| POST | `/auth/register` | ลงทะเบียนผู้ใช้ใหม่ | 3/3600s |
| GET | `/auth/me` | ข้อมูลผู้ใช้ปัจจุบัน | 50/60s |

### 3.3 Redis Cache Keys

| Cache Key | TTL | คำอธิบาย |
|-----------|-----|----------|
| `user-permissions:{userId}` | 1 ชม. | สิทธิ์ของผู้ใช้ |
| `user-sessions:{userId}` | 15 นาที | Session ผู้ใช้ |
| `token-blacklist:{token}` | 1 วัน | Token ที่ถูกเพิกถอน |

---

## 4. โมดูลที่ 2: Job Card Management

### 4.1 Database Tables

| ตาราง | คำอธิบาย |
|-------|----------|
| `t_job` | ใบงานหลัก (job_no, customer_id, car_id, mechanic_id, status) |
| `t_job_service` | รายการบริการในใบงาน |
| `t_job_part_sales` | อะไหล่ที่ขายในใบงาน |
| `t_job_service_car_symptom` | อาการรถ |
| `t_job_diag_trouble_code` | รหัสข้อผิดพลาด (OBD2) |
| `t_job_status_history` | ประวัติการเปลี่ยนสถานะ |

### 4.2 Job Statuses

`OPEN → IN_PROGRESS → QUOTATION_PENDING → QUOTATION_APPROVED → PART_PICKING → REPAIR_IN_PROGRESS → REPAIR_DONE → INVOICE_PENDING → INVOICE_CREATED → PAYMENT_RECEIVED → CLOSED`

### 4.3 API Endpoints

| Method | Path | คำอธิบาย | Rate Limit |
|--------|------|----------|------------|
| POST | `/jobs/create` | สร้างใบงานใหม่ | 30/60s |
| GET | `/jobs/list` | รายการใบงาน | 50/60s |
| GET | `/jobs/{id}` | ดูรายละเอียด | 100/60s |
| PUT | `/jobs/{id}/status` | เปลี่ยนสถานะ | 60/60s |
| GET | `/jobs/report/{id}` | PDF ใบงาน | 20/300s |

### 4.4 Redis Cache Keys

| Cache Key | TTL | คำอธิบาย |
|-----------|-----|----------|
| `jobs:{jobId}` | 30 นาที | ข้อมูลใบงาน |
| `job_status_summary:{whitelabelId}` | 5 นาที | สรุปสถานะใบงาน |

---

## 5. โมดูลที่ 3: Customer Management

### 5.1 Database Tables

| ตาราง | คำอธิบาย |
|-------|----------|
| `m_customer` | ข้อมูลลูกค้า (ชื่อ, เบอร์โทร, อีเมล, ที่อยู่, TIN) |
| `m_car` | ข้อมูลรถยนต์ (ทะเบียน, ยี่ห้อ, รุ่น, ปี, VIN) |
| `m_car_service_history` | ประวัติการซ่อมบำรุงของรถ |

### 5.2 API Endpoints

| Method | Path | คำอธิบาย | Rate Limit |
|--------|------|----------|------------|
| POST | `/customers/create` | สร้างลูกค้าใหม่ | 20/60s |
| GET | `/customers/{id}` | ดูข้อมูลลูกค้า | 100/60s |
| GET | `/customers/phone/{phone}` | ค้นหาด้วยเบอร์โทร | 60/60s |
| POST | `/customers/search` | ค้นหาลูกค้า | 50/60s |
| POST | `/cars/create` | เพิ่มรถใหม่ | 30/60s |
| GET | `/cars/plate/{licensePlate}` | ค้นหาด้วยทะเบียน | 60/60s |

### 5.3 Redis Cache Keys

| Cache Key | TTL | คำอธิบาย |
|-----------|-----|----------|
| `customers:{customerId}` | 1 ชม. | ข้อมูลลูกค้า |
| `customer_phone:{phone}` | 1 ชม. | ลูกค้าตามเบอร์โทร |
| `cars:{carId}` | 1 ชม. | ข้อมูลรถยนต์ |
| `car_plate:{licensePlate}` | 1 ชม. | รถตามทะเบียน |

---

## 6. โมดูลที่ 4: Quotation

### 6.1 Database Tables

| ตาราง | คำอธิบาย |
|-------|----------|
| `t_quotation` | ใบเสนอราคา (quotation_no, job_id, status, subtotal, tax, total) |
| `t_quotation_part` | อะไหล่ในใบเสนอราคา |
| `t_quotation_service` | บริการในใบเสนอราคา |
| `t_quotation_status_history` | ประวัติการเปลี่ยนสถานะ |

### 6.2 API Endpoints

| Method | Path | คำอธิบาย | Rate Limit |
|--------|------|----------|------------|
| POST | `/quotations/create` | สร้างใบเสนอราคา | 20/60s |
| GET | `/quotations/{id}` | ดูรายละเอียด | 100/60s |
| GET | `/quotations/job/{jobId}` | ดึงตาม Job | 80/60s |
| PUT | `/quotations/{id}/approve` | อนุมัติ | 20/60s |
| GET | `/quotations/{id}/pdf` | PDF ใบเสนอราคา | 15/300s |

### 6.3 Redis Cache Keys

| Cache Key | TTL | คำอธิบาย |
|-----------|-----|----------|
| `quotations:{quotationId}` | 30 นาที | ข้อมูลใบเสนอราคา |
| `quotation_job:{jobId}` | 30 นาที | Quotation ตาม Job |

---

## 7. โมดูลที่ 5: Purchase Order

### 7.1 Database Tables

| ตาราง | คำอธิบาย |
|-------|----------|
| `t_purchase_order_header` | ใบสั่งซื้อ (po_no, supplier_id, status, total) |
| `t_purchase_order_detail` | รายการใบสั่งซื้อ |
| `t_purchase_order_status_history` | ประวัติสถานะ PO |

### 7.2 API Endpoints

| Method | Path | คำอธิบาย | Rate Limit |
|--------|------|----------|------------|
| POST | `/purchase-orders/create` | สร้าง PO | 20/60s |
| POST | `/purchase-orders/from-quotation/{quotationId}` | สร้างจาก Quotation | 15/60s |
| POST | `/{id}/send` | ส่งให้ Supplier | 10/60s |
| POST | `/{id}/receive` | รับสินค้า | 15/60s |
| GET | `/{id}/pdf` | PDF PO | 15/300s |

### 7.3 Redis Cache Keys

| Cache Key | TTL | คำอธิบาย |
|-----------|-----|----------|
| `purchase_orders:{poId}` | 30 นาที | ข้อมูล PO |
| `po_quotation:{quotationId}` | 30 นาที | PO ตาม Quotation |
| `po_suggestion:{jobId}` | 15 นาที | แนะนำ PO |

---

## 8. โมดูลที่ 6: Inventory Management

### 8.1 Database Tables

| ตาราง | คำอธิบาย |
|-------|----------|
| `m_part_master` | อะไหล่หลัก (part_code, part_name, stock_quantity, reorder_level) |
| `m_stock_location` | ตำแหน่งจัดเก็บ |
| `t_inventory` | การเคลื่อนไหวสินค้า (RECEIVE, ISSUE, ADJUSTMENT) |
| `t_inventory_adjustment_header/detail` | การปรับปรุงสต็อก |
| `t_stocktake_header/detail` | การตรวจนับสต็อก |
| `t_part_picking_request/detail` | คำขอเบิกอะไหล่ |

### 8.2 API Endpoints

| Method | Path | คำอธิบาย | Rate Limit |
|--------|------|----------|------------|
| POST | `/inventory/receive` | รับสินค้าเข้า | 20/60s |
| POST | `/inventory/issue` | เบิกจ่ายสินค้า | 30/60s |
| GET | `/inventory/part/{partId}` | ดู Inventory อะไหล่ | 80/60s |
| POST | `/part-picking/create` | สร้างคำขอเบิก | 30/60s |
| PUT | `/part-picking/{id}/confirm` | ยืนยันการเบิก | 20/60s |
| GET | `/parts/low-stock` | อะไหล่ต่ำกว่าเกณฑ์ | 30/60s |

### 8.3 Redis Cache Keys

| Cache Key | TTL | คำอธิบาย |
|-----------|-----|----------|
| `parts:{partId}` | 1 ชม. | ข้อมูลอะไหล่ |
| `part_code:{partCode}` | 1 ชม. | อะไหล่ตามรหัส |
| `stock_summary:{partId}` | 15 นาที | สรุปสต็อก |

---

## 9. โมดูลที่ 7: Payment Management

### 9.1 Database Tables

| ตาราง | คำอธิบาย |
|-------|----------|
| `m_payment_method` | วิธีการชำระเงิน (CASH, BANK_TRANSFER, CREDIT_CARD) |
| `t_payment` | การชำระเงิน (payment_no, invoice_id, amount, status) |
| `t_receipt` | ใบเสร็จรับเงิน |
| `t_payment_history` | ประวัติสถานะการชำระ |
| `t_outstanding_balance` | ยอดคงค้าง (Dashboard) |

### 9.2 API Endpoints

| Method | Path | คำอธิบาย | Rate Limit |
|--------|------|----------|------------|
| POST | `/payments/record` | บันทึกการชำระเงิน | 20/60s |
| GET | `/payments/invoice/{invoiceId}` | ดึงตาม Invoice | 60/60s |
| GET | `/payments/outstanding/{customerId}` | ยอดคงค้าง | 40/60s |
| POST | `/{id}/refund` | คืนเงิน | 10/3600s |
| GET | `/receipts/{id}/pdf` | PDF ใบเสร็จ | 15/300s |

### 9.3 Redis Cache Keys

| Cache Key | TTL | คำอธิบาย |
|-----------|-----|----------|
| `payments:{paymentId}` | 1 ชม. | ข้อมูลการชำระเงิน |
| `payment_invoice:{invoiceId}` | 1 ชม. | Payment ตาม Invoice |
| `receipts:{receiptId}` | 1 ชม. | ข้อมูลใบเสร็จ |

---

## 10. โมดูลที่ 8: Dashboard & Reports

### 10.1 Database Views

| View | คำอธิบาย |
|------|----------|
| `v_dashboard_sales_overview` | ภาพรวมยอดขาย |
| `v_dashboard_job_status` | สรุปสถานะใบงาน |
| `v_dashboard_inventory_overview` | ภาพรวมสินค้าคงคลัง |
| `v_dashboard_top_parts` | อะไหล่ขายดี |
| `v_dashboard_financial_summary` | สรุปการเงิน |

### 10.2 API Endpoints

| Method | Path | คำอธิบาย | Rate Limit |
|--------|------|----------|------------|
| GET | `/dashboard/overview` | ภาพรวม Dashboard | 30/60s |
| GET | `/dashboard/sales` | ภาพรวมยอดขาย | 30/60s |
| GET | `/dashboard/job-status` | สรุปสถานะใบงาน | 30/60s |
| GET | `/dashboard/inventory` | ภาพรวมสินค้าคงคลัง | 30/60s |
| GET | `/dashboard/top-parts` | อะไหล่ขายดี | 20/60s |
| GET | `/dashboard/revenue` | รายได้แยกช่วงเวลา | 25/60s |
| POST | `/reports/daily` | รายงานรายวัน | 10/300s |
| POST | `/reports/monthly` | รายงานรายเดือน | 10/300s |
| POST | `/reports/yearly` | รายงานรายปี | 5/3600s |

### 10.3 Redis Cache Keys

| Cache Key | TTL | คำอธิบาย |
|-----------|-----|----------|
| `dashboard_overview:{whitelabelId}` | 5 นาที | ภาพรวม Dashboard |
| `sales_by_period:{whitelabelId}:{period}` | 5 นาที | ยอดขายแยกช่วง |
| `top_parts:{whitelabelId}` | 5 นาที | อะไหล่ขายดี |

---

## 11. โมดูลที่ 9: Document Management

### 11.1 Database Tables

| ตาราง | คำอธิบาย |
|-------|----------|
| `m_document_template` | เทมเพลตเอกสาร (JasperReports .jrxml) |
| `t_document` | เอกสารที่สร้างแล้ว (PDF, Excel) |
| `t_document_history` | ประวัติเอกสาร |
| `t_ocr_result` | ผลลัพธ์ OCR |

### 11.2 API Endpoints

| Method | Path | คำอธิบาย | Rate Limit |
|--------|------|----------|------------|
| POST | `/documents/generate` | สร้างเอกสารจากเทมเพลต | 15/300s |
| POST | `/documents/upload` | อัปโหลดเอกสาร | 20/300s |
| GET | `/documents/{id}/download` | ดาวน์โหลดเอกสาร | 30/60s |
| POST | `/documents/ocr` | กระบวนการ OCR | 10/60s |
| POST | `/templates/upload` | อัปโหลดเทมเพลต | 10/300s |

### 11.3 Redis Cache Keys

| Cache Key | TTL | คำอธิบาย |
|-----------|-----|----------|
| `documents:{documentId}` | 1 ชม. | ข้อมูลเอกสาร |
| `document_ref:{refType}:{refId}` | 1 ชม. | เอกสารตาม Reference |
| `templates:{templateCode}` | 2 ชม. | เทมเพลตเอกสาร |

---

## 12. โมดูลที่ 10: Email Service

### 12.1 Database Tables

| ตาราง | คำอธิบาย |
|-------|----------|
| `m_email_template` | เทมเพลตอีเมล (HTML, Text) |
| `t_email_history` | ประวัติการส่งอีเมล |
| `t_email_queue` | คิวอีเมล (สำหรับ Async) |

### 12.2 API Endpoints

| Method | Path | คำอธิบาย | Rate Limit |
|--------|------|----------|------------|
| POST | `/email/send` | ส่งอีเมล | 20/60s |
| POST | `/email/send-template` | ส่งอีเมลจากเทมเพลต | 25/60s |
| POST | `/email/bulk` | ส่งอีเมลจำนวนมาก | 5/300s |
| GET | `/email/status/{emailId}` | สถานะการส่ง | 50/60s |
| POST | `/email/resend/{emailId}` | ส่งอีเมลซ้ำ | 10/3600s |

### 12.3 Redis Cache Keys

| Cache Key | TTL | คำอธิบาย |
|-----------|-----|----------|
| `email_templates:{templateCode}:{lang}` | 2 ชม. | เทมเพลตอีเมล |

---

## 13. โมดูลที่ 11: Batch Jobs

### 13.1 Database Tables

| ตาราง | คำอธิบาย |
|-------|----------|
| `m_batch_job` | งาน Batch (job_code, cron_expression, enabled) |
| `t_batch_job_history` | ประวัติการรันงาน |

### 13.2 6 Batch Jobs

| Job | Cron | เวลา | หน้าที่ |
|-----|------|------|--------|
| batch001 | `0 30 6 ? * *` | 06:30 | ส่งอีเมลแจ้งเตือนรายวัน |
| batch002 | `0 45 6 ? * *` | 06:45 | สร้างรายงานประจำวัน |
| batch003 | `0 30 6 ? * *` | 06:30 | อัปเดตสถานะงานค้าง |
| batch004 | `0 0 3 ? * *` | 03:00 | ล้างข้อมูล/ซิงค์ฐานข้อมูล |
| batch005 | `0 0/30 * * * ?` | ทุก 30 นาที | ซิงค์ข้อมูล Realtime |
| batch006 | `0 30 6 ? * *` | 06:30 | ส่งสรุปยอดขาย |

### 13.3 API Endpoints

| Method | Path | คำอธิบาย | Rate Limit |
|--------|------|----------|------------|
| GET | `/batch-jobs/list` | รายการงาน Batch | 20/60s |
| GET | `/{jobCode}/status` | สถานะงาน | 30/60s |
| POST | `/{jobCode}/trigger` | สั่งรันงาน | 5/3600s |
| GET | `/{jobCode}/history` | ประวัติการรัน | 20/60s |

### 13.4 Redis Cache Keys

| Cache Key | TTL | คำอธิบาย |
|-----------|-----|----------|
| `batch_jobs:{jobCode}` | 5 นาที | ข้อมูลงาน Batch |
| `batch_lock:{jobCode}` | 5 นาที | Distributed Lock (Redis) |

---

## 14. โมดูลที่ 12: Multi-Language (i18n)

### 14.1 18 ภาษาที่รองรับ

| รหัส | ภาษา | สัญลักษณ์ |
|------|------|-----------|
| `th` | ภาษาไทย | 🇹🇭 |
| `en` | English | 🇬🇧 |
| `zh` | 中文 | 🇨🇳 |
| `ja` | 日本語 | 🇯🇵 |
| `ko` | 한국어 | 🇰🇷 |
| `es` | Español | 🇪🇸 |
| `fr` | Français | 🇫🇷 |
| `de` | Deutsch | 🇩🇪 |
| `it` | Italiano | 🇮🇹 |
| `pt` | Português | 🇵🇹 |
| `ru` | Русский | 🇷🇺 |
| `ar` | العربية | 🇸🇦 (RTL) |
| `hi` | हिन्दी | 🇮🇳 |
| `id` | Bahasa Indonesia | 🇮🇩 |
| `ms` | Bahasa Melayu | 🇲🇾 |
| `vi` | Tiếng Việt | 🇻🇳 |
| `my` | မြန်မာဘာသာ | 🇲🇲 |
| `lo` | ພາສາລາວ | 🇱🇦 |

### 14.2 Database Tables

| ตาราง | คำอธิบาย |
|-------|----------|
| `m_language` | ข้อมูลภาษา (language_code, locale, is_rtl) |
| `m_translation` | ข้อความที่แปลแล้ว (message_key, language_code, message_text) |

### 14.3 API Endpoints

| Method | Path | คำอธิบาย | Rate Limit |
|--------|------|----------|------------|
| GET | `/languages` | รายการภาษาที่รองรับ | 50/60s |
| GET | `/languages/current` | ภาษาปัจจุบัน | 100/60s |
| POST | `/languages/switch` | สลับภาษา | 20/60s |
| GET | `/languages/messages/{lang}` | ข้อความทั้งหมดของภาษา | 30/60s |

### 14.4 Redis Cache Keys

| Cache Key | TTL | คำอธิบาย |
|-----------|-----|----------|
| `messages:{messageKey}:{lang}` | 1 ชม. | ข้อความที่แปลแล้ว |
| `messages_all:{lang}` | 1 ชม. | ข้อความทั้งหมดของภาษา |
| `languages:{languageCode}` | 2 ชม. | ข้อมูลภาษา |

---

## 15. โมดูลที่ 13: IoT & GPS Tracking

### 15.1 Database Tables

| ตาราง | คำอธิบาย |
|-------|----------|
| `m_iot_device` | อุปกรณ์ IoT (device_id, status, last_location) |
| `t_gps_data` | ข้อมูล GPS (ละติจูด, ลองจิจูด, ความเร็ว) |
| `t_device_history` | ประวัติการทำงานของอุปกรณ์ |
| `t_device_access_log` | บันทึกการเข้าถึงอุปกรณ์ |
| `m_geofence` | รั้วรอบขอบ (CIRCLE, POLYGON) |
| `t_geofence_alert` | การแจ้งเตือน Geofence |
| `t_auto_report` | รายงานอัตโนมัติ |

### 15.2 MQTT Topics

| Topic | คำอธิบาย |
|-------|----------|
| `gps/data` | ข้อมูล GPS จากอุปกรณ์ |
| `device/status` | สถานะอุปกรณ์ (ONLINE, OFFLINE) |
| `device/alert` | การแจ้งเตือนจากอุปกรณ์ |

### 15.3 API Endpoints

| Method | Path | คำอธิบาย | Rate Limit |
|--------|------|----------|------------|
| POST | `/iot/devices/register` | ลงทะเบียนอุปกรณ์ | 10/300s |
| GET | `/iot/devices/{id}` | ดูข้อมูลอุปกรณ์ | 100/60s |
| GET | `/iot/devices/{id}/location` | ตำแหน่งล่าสุด | 60/60s |
| GET | `/iot/devices/{id}/history` | ประวัติตำแหน่ง | 30/60s |
| POST | `/iot/mqtt/publish` | ส่งข้อความ MQTT | 20/60s |
| POST | `/iot/geofences` | สร้าง Geofence | 10/300s |
| GET | `/iot/geofences/alerts` | รายการแจ้งเตือน | 30/60s |

### 15.4 Redis Cache Keys

| Cache Key | TTL | คำอธิบาย |
|-----------|-----|----------|
| `iot_devices:{deviceId}` | 1 นาที | ข้อมูลอุปกรณ์ IoT |
| `iot_device_identifier:{identifier}` | 1 นาที | อุปกรณ์ตาม Hardware ID |
| `device_location:{deviceId}` | 30 วินาที | ตำแหน่งล่าสุด |

---

## 16. โมดูลที่ 14: Web Order System (WOS)

### 16.1 Database Tables

| ตาราง | คำอธิบาย |
|-------|----------|
| `m_catalogue_category` | หมวดหมู่สินค้าในแคตตาล็อก |
| `m_catalogue_item` | สินค้าในแคตตาล็อก (เชื่อมโยงกับ `m_part_master`) |
| `m_sales_price` | ราคาขาย (รองรับหลาย Tier) |
| `m_promotion` | โปรโมชัน/ส่วนลด |
| `t_shopping_cart` | ตะกร้าสินค้า (Session-based) |
| `t_shopping_cart_item` | รายการในตะกร้า |
| `t_web_order` | ออเดอร์ออนไลน์ |
| `t_web_order_item` | รายการออเดอร์ |
| `t_web_order_status_history` | ประวัติสถานะออเดอร์ |

### 16.2 Order Statuses

`PENDING → CONFIRMED → PROCESSING → SHIPPED → DELIVERED`

### 16.3 API Endpoints

| Method | Path | คำอธิบาย | Rate Limit |
|--------|------|----------|------------|
| GET | `/wos/catalogue` | รายการสินค้า | 100/60s |
| GET | `/wos/catalogue/{id}` | ดูสินค้า | 200/60s |
| GET | `/wos/catalogue/search` | ค้นหาสินค้า | 80/60s |
| POST | `/wos/cart/add` | เพิ่มสินค้าในตะกร้า | 50/60s |
| GET | `/wos/cart` | ดูตะกร้าสินค้า | 100/60s |
| POST | `/wos/cart/promotion` | ใช้โค้ดส่วนลด | 20/60s |
| POST | `/wos/orders` | สั่งซื้อ | 20/60s |
| GET | `/wos/orders/{id}` | ดูออเดอร์ | 100/60s |
| GET | `/wos/orders/list` | ประวัติการสั่งซื้อ | 50/60s |
| PUT | `/wos/orders/{id}/status` | เปลี่ยนสถานะ | 20/60s |

### 16.4 Redis Cache Keys

| Cache Key | TTL | คำอธิบาย |
|-----------|-----|----------|
| `catalogue:{itemId}` | 1 ชม. | สินค้าแคตตาล็อก |
| `catalogue_code:{itemCode}` | 1 ชม. | สินค้าตามรหัส |
| `catalogue_category:{categoryId}` | 1 ชม. | สินค้าตามหมวดหมู่ |
| `cart:{cartId}` | 30 นาที | ตะกร้าสินค้า |
| `orders:{orderId}` | 30 นาที | ข้อมูลออเดอร์ |

---

## 17. สรุปฐานข้อมูลทั้งหมด (Consolidated Database)

### 17.1 Master Data Tables (34 ตาราง)

| ลำดับ | ตาราง | โมดูล |
|-------|-------|-------|
| 1 | `m_user` | Auth |
| 2 | `m_role` | Auth |
| 3 | `m_permission` | Auth |
| 4 | `m_user_role` | Auth |
| 5 | `m_role_permission` | Auth |
| 6 | `m_user_token` | Auth |
| 7 | `m_customer` | Customer |
| 8 | `m_car` | Customer |
| 9 | `m_supplier` | Master |
| 10 | `m_part_master` | Inventory |
| 11 | `m_service` | Master |
| 12 | `m_category` | Master |
| 13 | `m_stock_location` | Inventory |
| 14 | `m_payment_method` | Payment |
| 15 | `m_currency` | Master |
| 16 | `m_exchange_rate` | Master |
| 17 | `m_country` | Master |
| 18 | `m_province` | Master |
| 19 | `m_city` | Master |
| 20 | `m_shop_profile` | Master |
| 21 | `m_staff` | Staff |
| 22 | `m_iot_device` | IoT |
| 23 | `m_geofence` | IoT |
| 24 | `m_document_template` | Document |
| 25 | `m_email_template` | Email |
| 26 | `m_language` | i18n |
| 27 | `m_translation` | i18n |
| 28 | `m_batch_job` | Batch |
| 29 | `m_catalogue_category` | WOS |
| 30 | `m_catalogue_item` | WOS |
| 31 | `m_sales_price` | WOS |
| 32 | `m_promotion` | WOS |
| 33 | `m_user_menu` | Auth |
| 34 | `m_user_job_role` | Auth |

### 17.2 Transaction Tables (44 ตาราง)

| ลำดับ | ตาราง | โมดูล |
|-------|-------|-------|
| 1-6 | `t_job`, `t_job_service`, `t_job_part_sales`, `t_job_symptom`, `t_job_diag_code`, `t_job_status_history` | Job |
| 7-10 | `t_quotation`, `t_quotation_part`, `t_quotation_service`, `t_quotation_status_history` | Quotation |
| 11-13 | `t_purchase_order_header`, `t_purchase_order_detail`, `t_po_status_history` | Purchase |
| 14-16 | `t_invoice_adjustment`, `t_invoice_adjustment_part`, `t_invoice_adjustment_service` | Invoice |
| 17-21 | `t_payment`, `t_receipt`, `t_payment_history`, `t_outstanding_balance`, `t_received_amount` | Payment |
| 22-28 | `t_inventory`, `t_inventory_adjustment_header`, `t_inventory_adjustment_detail`, `t_stocktake_header`, `t_stocktake_detail`, `t_part_picking_request`, `t_part_picking_detail` | Inventory |
| 29-31 | `t_document`, `t_document_history`, `t_ocr_result` | Document |
| 32-34 | `t_email_history`, `t_email_queue` | Email |
| 35 | `t_batch_job_history` | Batch |
| 36-40 | `t_gps_data`, `t_device_history`, `t_device_access_log`, `t_geofence_alert`, `t_auto_report` | IoT |
| 41-44 | `t_shopping_cart`, `t_shopping_cart_item`, `t_web_order`, `t_web_order_item`, `t_web_order_status_history` | WOS |

### 17.3 Dashboard Views (9 Views)

| ลำดับ | View | โมดูล |
|-------|------|-------|
| 1 | `v_dashboard_sales_overview` | Dashboard |
| 2 | `v_dashboard_job_status` | Dashboard |
| 3 | `v_dashboard_inventory_overview` | Dashboard |
| 4 | `v_dashboard_top_parts` | Dashboard |
| 5 | `v_dashboard_service_category` | Dashboard |
| 6 | `v_dashboard_financial_summary` | Dashboard |
| 7 | `v_dashboard_revenue_by_period` | Dashboard |
| 8 | `v_available_languages` | i18n |
| 9 | `v_email_analytics` | Email |

---

## 18. สรุป API ทั้งหมด (Consolidated API)

### 18.1 สรุปจำนวน Endpoints แยกตามโมดูล

| โมดูล | จำนวน Endpoints |
|-------|-----------------|
| Authentication | 7 |
| Job Card | 10 |
| Customer | 7 |
| Car | 7 |
| Quotation | 11 |
| Purchase Order | 12 |
| Inventory | 5 |
| Part Master | 7 |
| Part Picking | 7 |
| Payment | 8 |
| Receipt | 4 |
| Dashboard | 9 |
| Reports | 4 |
| Email | 5 |
| Batch Jobs | 7 |
| IoT Devices | 9 |
| MQTT | 1 |
| Geofence | 5 |
| WOS Catalogue | 6 |
| WOS Cart | 6 |
| WOS Orders | 8 |
| **รวม** | **~155 Endpoints** |

---

## 19. สรุป Redis Cache Keys (41 Keys)

| Cache Name | Key Pattern | TTL | โมดูล |
|------------|-------------|-----|-------|
| `user-permissions` | `{userId}` | 1 ชม. | Auth |
| `user-sessions` | `{userId}` | 15 นาที | Auth |
| `token-blacklist` | `{token}` | 1 วัน | Auth |
| `jobs` | `{jobId}` | 30 นาที | Job |
| `job_status_summary` | `{whitelabelId}` | 5 นาที | Job |
| `customers` | `{customerId}` | 1 ชม. | Customer |
| `customer_phone` | `{phone}` | 1 ชม. | Customer |
| `cars` | `{carId}` | 1 ชม. | Car |
| `car_plate` | `{licensePlate}` | 1 ชม. | Car |
| `quotations` | `{quotationId}` | 30 นาที | Quotation |
| `quotation_job` | `{jobId}` | 30 นาที | Quotation |
| `purchase_orders` | `{poId}` | 30 นาที | Purchase |
| `po_quotation` | `{quotationId}` | 30 นาที | Purchase |
| `po_suggestion` | `{jobId}` | 15 นาที | Purchase |
| `parts` | `{partId}` | 1 ชม. | Inventory |
| `part_code` | `{partCode}` | 1 ชม. | Inventory |
| `stock_summary` | `{partId}` | 15 นาที | Inventory |
| `payments` | `{paymentId}` | 1 ชม. | Payment |
| `payment_invoice` | `{invoiceId}` | 1 ชม. | Payment |
| `receipts` | `{receiptId}` | 1 ชม. | Receipt |
| `receipt_payment` | `{paymentId}` | 1 ชม. | Receipt |
| `dashboard_overview` | `{whitelabelId}` | 5 นาที | Dashboard |
| `sales_by_period` | `{whitelabelId}:{period}` | 5 นาที | Dashboard |
| `top_parts` | `{whitelabelId}` | 5 นาที | Dashboard |
| `documents` | `{documentId}` | 1 ชม. | Document |
| `document_ref` | `{refType}:{refId}` | 1 ชม. | Document |
| `templates` | `{templateCode}` | 2 ชม. | Document |
| `email_templates` | `{templateCode}:{lang}` | 2 ชม. | Email |
| `messages` | `{messageKey}:{lang}` | 1 ชม. | i18n |
| `messages_all` | `{lang}` | 1 ชม. | i18n |
| `languages` | `{languageCode}` | 2 ชม. | i18n |
| `languages_active` | `'all'` | 2 ชม. | i18n |
| `batch_jobs` | `{jobCode}` | 5 นาที | Batch |
| `batch_lock` | `{jobCode}` | 5 นาที | Batch |
| `iot_devices` | `{deviceId}` | 1 นาที | IoT |
| `iot_device_identifier` | `{identifier}` | 1 นาที | IoT |
| `device_location` | `{deviceId}` | 30 วินาที | IoT |
| `catalogue` | `{itemId}` | 1 ชม. | WOS |
| `catalogue_code` | `{itemCode}` | 1 ชม. | WOS |
| `catalogue_category` | `{categoryId}` | 1 ชม. | WOS |
| `cart` | `{cartId}` | 30 นาที | WOS |
| `orders` | `{orderId}` | 30 นาที | WOS |

---

## 20. สรุป Rate Limit Policy

| ประเภท Operation | Limit | Duration | คำอธิบาย |
|------------------|-------|----------|----------|
| **Login** | 5 | 300s | ป้องกัน Brute Force |
| **Register** | 3 | 3600s | ป้องกันการสมัครซ้ำ |
| **Create (POST)** | 15-30 | 60s | จำกัดการสร้างข้อมูล |
| **Update (PUT)** | 15-20 | 60s | จำกัดการแก้ไขข้อมูล |
| **Delete (DELETE)** | 5-10 | 3600s | ป้องกันการลบโดยไม่ตั้งใจ |
| **Read (GET)** | 50-200 | 60s | อนุญาตให้อ่านบ่อย |
| **PDF Generation** | 15 | 300s | การสร้าง PDF ใช้ทรัพยากรมาก |
| **Bulk Email** | 5 | 300s | ป้องกัน Spam |
| **OCR** | 10 | 60s | OCR ใช้ทรัพยากรมาก |
| **Batch Trigger** | 5 | 3600s | ป้องกันการรันงานซ้ำ |
| **MQTT Publish** | 20 | 60s | จำกัดการส่ง MQTT |

---

## 21. ภาคผนวก: การติดตั้งและการใช้งาน

### 21.1 ข้อกำหนดเบื้องต้น
- Java 21+
- Maven 3.8+
- Docker & Docker Compose
- Git

### 21.2 ขั้นตอนการติดตั้ง

```bash
# 1. Clone Repository
git clone https://github.com/your-org/auto-repair-system.git
cd auto-repair-system

# 2. กำหนดค่าตัวแปรสภาพแวดล้อม
cp .env.example .env
# แก้ไขไฟล์ .env ให้ตรงกับสภาพแวดล้อมจริง

# 3. เริ่มต้นฐานข้อมูลด้วย Docker Compose
docker-compose up -d

# 4. รันแอปพลิเคชัน
mvn spring-boot:run -Dspring-boot.run.profiles=dev

# 5. เข้าถึง Swagger UI
# http://localhost:5000/api/swagger-ui.html
```

### 21.3 การใช้งานเบื้องต้น
1. **เข้าสู่ระบบ** ด้วย `admin/admin`
2. **สร้าง Job Card** – ระบุลูกค้า, รถยนต์, อาการ, และช่าง
3. **สร้าง Quotation** – เพิ่มบริการและอะไหล่
4. **อนุมัติ Quotation** – ระบบสร้าง PO และเริ่มเบิกอะไหล่
5. **ดำเนินการซ่อม** – ช่างอัปเดตสถานะงาน
6. **ออก Invoice** – เมื่อซ่อมเสร็จ
7. **รับชำระเงิน** – บันทึกการชำระและออกใบเสร็จ

---

## ✅ สรุป

เอกสารฉบับนี้ครอบคลุมการออกแบบระบบบริหารจัดการอู่ซ่อมรถอย่างครบถ้วน ประกอบด้วย:

- **14 โมดูลหลัก** – ครบทุกฟังก์ชันการทำงาน
- **34 ตาราง Master Data** – ข้อมูลหลัก
- **44 ตาราง Transaction** – ประวัติการทำธุรกรรม
- **9 Dashboard Views** – สำหรับการวิเคราะห์
- **~155 API Endpoints** – พร้อม Rate Limit
- **41 Redis Cache Keys** – พร้อม TTL

ระบบพร้อมใช้งานเป็นเทมเพลตสำหรับพัฒนาระบบจริง และสามารถปรับแต่งให้เหมาะกับความต้องการเฉพาะของแต่ละองค์กร

---

**ผู้เขียน:** Kongnakorn Jantakun  
**วันที่:** 2026-07-04  
**เวอร์ชัน:** 3.0 (Final)  
**สถานะ:** ฉบับสมบูรณ์ ✅