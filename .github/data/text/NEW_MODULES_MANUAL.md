# 📘 คู่มือโมดูลใหม่ — New Modules Manual

## สารบัญ
1. [Job Card Management](#1-job-card-management)
2. [Customer Management](#2-customer-management)
3. [Quotation Management](#3-quotation-management)
4. [Purchase Order Management](#4-purchase-order-management)
5. [Payment & Receipt Management](#5-payment--receipt-management)
6. [Dashboard & Reports](#6-dashboard--reports)
7. [Document Management](#7-document-management)
8. [Email Service](#8-email-service)
9. [Batch Job Management](#9-batch-job-management)
10. [I18n / Translation](#10-i18n--translation)
11. [Web Order System (WOS)](#11-web-order-system-wos)
12. [Auth Enhancements](#12-auth-enhancements)

---

## 1. Job Card Management

**โมดูล:** `internal/modules/job/`
**ตาราง:** `t_job`, `t_job_service`, `t_job_part_sales`, `t_job_service_car_symptom`, `t_job_diag_trouble_code`, `t_job_status_history`

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/jobs` | สร้างใบงานใหม่ |
| GET | `/api/jobs/{id}` | ดึงข้อมูลใบงาน |
| GET | `/api/jobs` | รายการใบงาน (มี filter) |
| PUT | `/api/jobs/{id}` | แก้ไขใบงาน |
| DELETE | `/api/jobs/{id}` | ลบใบงาน (soft delete) |
| PUT | `/api/jobs/{id}/status` | เปลี่ยนสถานะใบงาน |
| POST | `/api/jobs/{id}/services` | เพิ่มบริการในใบงาน |
| POST | `/api/jobs/{id}/parts` | เพิ่มอะไหล่ในใบงาน |
| GET | `/api/jobs/{id}/history` | ประวัติการเปลี่ยนสถานะ |
| GET | `/api/jobs/report/{id}` | รายงาน PDF |

### สถานะ Job

`OPEN → IN_PROGRESS → QUOTATION_PENDING → QUOTATION_APPROVED → PART_PICKING → REPAIR_IN_PROGRESS → REPAIR_DONE → INVOICE_PENDING → INVOICE_CREATED → PAYMENT_RECEIVED → CLOSED`

Special: `CANCELLED`, `ON_HOLD`, `WAITING_PARTS`

---

## 2. Customer Management

**โมดูล:** `internal/modules/customer/`
**ตาราง:** `m_customer`, `m_car`, `m_car_service_history`

### API Endpoints — Customer

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/customers` | สร้างลูกค้า |
| GET | `/api/customers/{id}` | ดึงข้อมูลลูกค้า |
| PUT | `/api/customers/{id}` | แก้ไขลูกค้า |
| DELETE | `/api/customers/{id}` | ลบลูกค้า (soft delete) |
| POST | `/api/customers/search` | ค้นหาลูกค้า |
| GET | `/api/customers/phone/{phone}` | ค้นหาลูกค้าด้วยเบอร์โทร |
| GET | `/api/customers/{id}/history` | ประวัติการซ่อม |
| POST | `/api/customers/record-visit` | บันทึกการมาใช้บริการ |

### API Endpoints — Car

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/cars` | เพิ่มรถยนต์ |
| GET | `/api/cars/{id}` | ดึงข้อมูลรถยนต์ |
| PUT | `/api/cars/{id}` | แก้ไขรถยนต์ |
| DELETE | `/api/cars/{id}` | ลบรถยนต์ |
| GET | `/api/cars/customer/{customerId}` | รถยนต์ทั้งหมดของลูกค้า |
| GET | `/api/cars/plate/{licensePlate}` | ค้นหาด้วยทะเบียน |
| POST | `/api/cars/search` | ค้นหารถยนต์ |

---

## 3. Quotation Management

**โมดูล:** `internal/modules/quotation/`
**ตาราง:** `t_quotation`, `t_quotation_part`, `t_quotation_service`, `t_quotation_status_history`

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/quotations` | สร้างใบเสนอราคา |
| GET | `/api/quotations/{id}` | ดึงใบเสนอราคา |
| GET | `/api/quotations` | รายการใบเสนอราคา |
| PUT | `/api/quotations/{id}` | แก้ไขใบเสนอราคา |
| DELETE | `/api/quotations/{id}` | ลบ (soft delete) |
| PUT | `/api/quotations/{id}/approve` | อนุมัติใบเสนอราคา |
| PUT | `/api/quotations/{id}/reject` | ปฏิเสธ |
| GET | `/api/quotations/{id}/pdf` | สร้าง PDF |
| GET | `/api/quotations/{id}/history` | ประวัติสถานะ |
| GET | `/api/quotations/job/{jobId}` | ดึงตาม Job ID |
| POST | `/api/quotations/parts` | เพิ่มอะไหล่ |
| PUT | `/api/quotations/parts/{id}` | แก้ไขอะไหล่ |
| DELETE | `/api/quotations/parts/{id}` | ลบอะไหล่ |
| GET | `/api/quotations/parts/quotation/{qId}` | รายการอะไหล่ทั้งหมด |
| POST | `/api/quotations/services` | เพิ่มบริการ |
| PUT | `/api/quotations/services/{id}` | แก้ไขบริการ |
| DELETE | `/api/quotations/services/{id}` | ลบบริการ |
| GET | `/api/quotations/services/quotation/{qId}` | รายการบริการทั้งหมด |

### สถานะ Quotation

```
DRAFT → PENDING → APPROVED → CONVERTED (to PO)
DRAFT → PENDING → REJECTED
DRAFT → PENDING → EXPIRED
```

---

## 4. Purchase Order Management

**โมดูล:** `internal/modules/purchaseorder/`
**ตาราง:** `t_po_header`, `t_po_detail`, `t_po_status_history`

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/purchase-orders` | สร้าง PO |
| GET | `/api/purchase-orders/{id}` | ดึง PO พร้อมรายการ |
| GET | `/api/purchase-orders` | รายการ PO |
| PUT | `/api/purchase-orders/{id}` | แก้ไข PO |
| DELETE | `/api/purchase-orders/{id}` | ลบ |
| PUT | `/api/purchase-orders/{id}/send` | ส่ง PO |
| PUT | `/api/purchase-orders/{id}/confirm` | ยืนยัน PO |
| PUT | `/api/purchase-orders/{id}/receive` | รับสินค้า |
| PUT | `/api/purchase-orders/{id}/cancel` | ยกเลิก PO |
| GET | `/api/purchase-orders/{id}/pdf` | PDF |
| GET | `/api/purchase-orders/{id}/history` | ประวัติสถานะ |
| GET | `/api/purchase-orders/{id}/suggestions` | แนะนำอะไหล่ตาม Job |
| POST | `/api/purchase-orders/from-quotation` | สร้าง PO จาก Quotation |

### สถานะ PO

```
DRAFT → SENT → CONFIRMED → SHIPPED → RECEIVED → CANCELLED
```

---

## 5. Payment & Receipt Management

**โมดูล:** `internal/modules/payment/`
**ตาราง:** `t_payment`, `t_receipt`, `t_payment_history`, `t_outstanding_balance`

### API Endpoints — Payment

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/payments` | บันทึกการชำระเงิน |
| GET | `/api/payments/{id}` | ดึงข้อมูล payment |
| GET | `/api/payments/invoice/{id}` | ตาม Invoice ID |
| POST | `/api/payments/search` | ค้นหา |
| GET | `/api/payments/outstanding/{customerId}` | ยอดค้างชำระ |
| GET | `/api/payments/history/{customerId}` | ประวัติการชำระ |
| POST | `/api/payments/{id}/refund` | คืนเงิน |
| PUT | `/api/payments/{id}/cancel` | ยกเลิก |

### API Endpoints — Receipt

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/receipts/{id}` | ดึงใบเสร็จ |
| GET | `/api/receipts/payment/{paymentId}` | ตาม payment ID |
| GET | `/api/receipts/{id}/pdf` | PDF |
| PUT | `/api/receipts/{id}/cancel` | ยกเลิกใบเสร็จ |

---

## 6. Dashboard & Reports

**โมดูล:** `internal/modules/dashboard/`

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/dashboard/stats` | สถิติรวม (Job/Payment/Customer count) |
| GET | `/api/dashboard/revenue` | รายได้ตามช่วงเวลา |
| GET | `/api/dashboard/top-parts` | อะไหล่ที่ขายดีที่สุด |
| GET | `/api/dashboard/job-status` | สรุปสถานะ Job |

---

## 7. Document Management

**โมดูล:** `internal/modules/document/`
**ตาราง:** `document`

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/documents/upload` | อัปโหลดไฟล์ |
| GET | `/api/documents/{id}` | ดาวน์โหลดไฟล์ |
| GET | `/api/documents` | รายการไฟล์ |
| DELETE | `/api/documents/{id}` | ลบไฟล์ |

---

## 8. Email Service

**โมดูล:** `internal/modules/email/`
**ตาราง:** `email_log`, `email_config`

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/email/send` | ส่งอีเมล |
| GET | `/api/email/logs` | ประวัติการส่ง |
| GET | `/api/email/logs/{id}` | log รายการ |
| GET | `/api/email/config` | ดู config ปัจจุบัน |
| PUT | `/api/email/config` | อัปเดต SMTP config |

---

## 9. Batch Job Management

**โมดูล:** `internal/modules/batch/`
**ตาราง:** `t_batch_job`, `t_batch_job_log`

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/batch/jobs` | สร้าง batch job |
| GET | `/api/batch/jobs/{id}` | ดึง job |
| GET | `/api/batch/jobs` | รายการ jobs |
| PUT | `/api/batch/jobs/{id}` | แก้ไข |
| DELETE | `/api/batch/jobs/{id}` | ลบ |
| POST | `/api/batch/jobs/{id}/run` | รันทันที |
| GET | `/api/batch/jobs/{id}/logs` | logs ของ job |

---

## 10. I18n / Translation

**โมดูล:** `internal/modules/i18n/`
**ตาราง:** `m_translation`

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/i18n/{locale}` |  translations ทั้งหมดตาม locale |
| GET | `/api/i18n/{locale}/{key}` | ค่าเฉพาะ key |
| POST | `/api/i18n` | สร้าง translation |
| PUT | `/api/i18n/{id}` | แก้ไข translation |
| DELETE | `/api/i18n/{id}` | ลบ translation |

---

## 11. Web Order System (WOS)

**โมดูล:** `internal/modules/wos/`
**ตาราง:** `wos_order`

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/wos/orders` | สร้าง order ออนไลน์ |
| GET | `/api/wos/orders/{id}` | ดึง order |
| GET | `/api/wos/orders` | รายการ orders |
| PUT | `/api/wos/orders/{id}/status` | อัปเดตสถานะ |

---

## 12. Auth Enhancements

**โมดูล:** `internal/modules/auth/`
**ตาราง:** `refresh_tokens`

เพิ่มเติมจากระบบ auth เดิม:
- `AuthUseCaseI` interface ที่ wrap `users.UserUseCaseI`
- `AuthPgRepository` สำหรับจัดการ refresh tokens ใน PostgreSQL
- แยก auth presenter (DTOs) จาก users presenter

### API Endpoints (คงเดิม)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/auth/login` | ล็อกอินด้วย username |
| POST | `/auth/signin` | ล็อกอินด้วย email |
| GET | `/auth/refresh` | ขอ token ใหม่ |
| GET | `/auth/logout` | ออกจากระบบ |
| GET | `/auth/logoutall` | ออกจากระบบทุก session |
| GET | `/auth/publickey` | ดึง public key |
| GET | `/auth/verifyemail` | ยืนยันอีเมล |
| POST | `/auth/forgotpassword` | ลืมรหัสผ่าน |
| PATCH | `/auth/resetpassword` | ตั้งรหัสผ่านใหม่ |

---

## ไฟล์ Migration

| File | Description |
|------|-------------|
| `migrations/20260712_new_modules_schema.sql` | DDL สำหรับทุกตารางใหม่ + Triggers |
| `migrations/20260712_seed_data.sql` | ข้อมูลตัวอย่าง |

## วิธีรัน Migration

```bash
# รันผ่าน GORM AutoMigrate (ใช้กับ Go run)
go run cmd/api/main.go migrate

# หรือรัน SQL โดยตรง
psql -h localhost -U postgres -d db -f migrations/20260712_new_modules_schema.sql
psql -h localhost -U postgres -d db -f migrations/20260712_seed_data.sql
```
