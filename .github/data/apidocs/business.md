# Business Modules API — Quotation / Purchase Order / Payment / Receipt / Job / Customer / Car

> เจ้าของ: icmongolang · Framework: Go + Chi · Auth: **JWT Bearer (ทุก endpoint)**
> Middleware: `Verifier(true)` → `Authenticator()` → `CurrentUser()` → `ActiveUser()`
> ไม่มี endpoint public ในกลุ่มนี้

---

## Response Envelope

```json
{ "data": { ... }, "error": null, "is_success": true }
```
Error:
```json
{ "data": null, "error": {"status":400,"status_text":"...","msg":"..."}, "is_success": false }
```

> ⚠️ **JSON key convention ไม่สม่ำเสมอ**: quotation / purchaseorder / payment / receipt ใช้ **snake_case**; job / customer / car ใช้ **camelCase**

---

## 1. Quotation — `/api/quotation`

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| GET | `/quotation` | `GetMulti` | Bearer |
| POST | `/quotation` | `Create` | Bearer |
| GET | `/quotation/{id}` | `Get` | Bearer |
| PUT | `/quotation/{id}` | `Update` | Bearer |
| DELETE | `/quotation/{id}` | `Delete` | Bearer |
| GET | `/quotation/{id}/pdf` | `GetPDF` | Bearer |

### POST `/quotation` — Create
Body (`QuotationCreate`); `quotation_no`, `job_id`, `customer_id`, `expiry_date` required:
```json
{
  "quotation_no": "QTN-2024-0001",
  "job_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "customer_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "quotation_date": "2024-01-15T10:00:00Z",
  "expiry_date": "2024-02-15T10:00:00Z",
  "subtotal": 1000.0,
  "tax_rate": 7.0,
  "tax_amount": 70.0,
  "discount_type": "PERCENTAGE",
  "discount_value": 100.0,
  "total": 970.0,
  "currency": "THB",
  "exchange_rate": 1.0,
  "notes": "",
  "terms_and_conditions": ""
}
```

### GET `/quotation?page=1&per_page=10`
Query: `page` (default 1), `per_page`/`limit` (default 10, max 100)
Response (`PaginatedQuotationResponse` — field `items`):
```json
{
  "data": { "items": [ { "...QuotationResponse" } ], "total": 42, "page": 1, "per_page": 10, "total_pages": 5 },
  "error": null, "is_success": true
}
```

### GET `/quotation/{id}`
Response `QuotationResponse`: `id`, `quotation_no`, `job_id`, `customer_id`, `quotation_date`, `expiry_date`, `status`, `subtotal`, `tax_rate`, `tax_amount`, `discount_type`, `discount_value`, `total`, `amount_in_words_th`, `amount_in_words_en`, `currency`, `exchange_rate`, `notes`, `terms_and_conditions`, `approved_by`, `approved_at`, `rejected_reason`, `converted_to_po`, `user_id`, `whitelabel_id`, `created_at`, `updated_at`

### PUT `/quotation/{id}` — partial update
Body `QuotationUpdate` (ทุก field pointer):
```json
{ "expiry_date": "2024-03-15T10:00:00Z", "notes": "แก้ไขเงื่อนไข", "total": 990.0 }
```

---

## 2. Purchase Order — `/api/purchase-orders`

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| GET | `/purchase-orders` | `List` | Bearer |
| POST | `/purchase-orders` | `Create` | Bearer |
| GET | `/purchase-orders/suggestions/{jobId}` | `GetSuggestions` | Bearer |
| POST | `/purchase-orders/from-quotation/{quotationId}` | `CreateFromQuotation` | Bearer |
| GET | `/purchase-orders/{id}` | `GetByID` | Bearer |
| PUT | `/purchase-orders/{id}` | `Update` | Bearer |
| DELETE | `/purchase-orders/{id}` | `Delete` | Bearer |
| POST | `/purchase-orders/{id}/send` | `Send` | Bearer |
| PUT | `/purchase-orders/{id}/confirm` | `Confirm` | Bearer |
| POST | `/purchase-orders/{id}/receive` | `Receive` | Bearer |
| PUT | `/purchase-orders/{id}/cancel` | `Cancel` | Bearer |
| GET | `/purchase-orders/{id}/pdf` | `GetPDF` | Bearer |
| GET | `/purchase-orders/{id}/history` | `GetStatusHistory` | Bearer |

### POST `/purchase-orders` — Create (with items)
`supplier_id` + `items[].part_id`/`quantity_ordered`/`unit_price` required:
```json
{
  "quotation_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "job_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "supplier_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "expected_delivery_date": "2026-07-20T00:00:00Z",
  "currency": "THB", "exchange_rate": 1.0, "shipping_cost": 0,
  "payment_terms": "Net 30", "delivery_address": "123/4 ถนนสุขุมวิท",
  "notes": "", "terms_and_conditions": "", "tax_rate": 7.0,
  "discount_type": "PERCENTAGE", "discount_value": 0,
  "items": [
    { "part_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6", "quantity_ordered": 10, "unit_price": 150.0, "discount": 0, "note": "" }
  ]
}
```
Handler sets `status: "DRAFT"`; `subtotal: 0`, `total: 0` (คำนวณใน usecase)

### GET `/purchase-orders?page=1&per_page=10&supplier_id&status&start_date&end_date`
Query: `page`, `per_page` (max 100), `supplier_id`, `status` (DRAFT/SENT/CONFIRMED/SHIPPED/RECEIVED/CANCELLED), `start_date`, `end_date`
Response (`PaginatedPurchaseOrdersResponse` — field `items`)

### GET `/purchase-orders/{id}`
Response `PurchaseOrderResponse` (รวม `items[]`):
```json
{
  "data": {
    "id": "...", "po_no": "PO-2026-0001", "quotation_id": "...", "job_id": "...", "supplier_id": "...",
    "po_date": "...", "expected_delivery_date": "...", "actual_delivery_date": null,
    "status": "DRAFT", "subtotal": 1500.0, "tax_rate": 7.0, "tax_amount": 105.0,
    "discount_type": "PERCENTAGE", "discount_value": 0, "total": 1605.0,
    "currency": "THB", "exchange_rate": 1.0, "shipping_cost": 0,
    "payment_terms": "Net 30", "delivery_address": "123/4", "notes": "", "terms_and_conditions": "",
    "sent_at": null, "confirmed_at": null, "received_by": null,
    "created_at": "...", "updated_at": null, "user_id": "...", "whitelabel_id": "...",
    "items": [
      { "id": "...", "part_id": "...", "quantity_ordered": 10, "quantity_received": 0,
        "unit_price": 150.0, "total_price": 1500.0, "discount": 0, "net_price": 1500.0, "note": "" }
    ]
  },
  "error": null, "is_success": true
}
```

### POST `/purchase-orders/{id}/receive`
Body (`PurchaseOrderReceiveRequest`):
```json
{ "items": [ { "detail_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6", "received_quantity": 5 } ] }
```

### PUT `/purchase-orders/{id}/cancel`
Body (`PurchaseOrderStatusRequest`):
```json
{ "reason": "ขอเปลี่ยนสถานะใบสั่งซื้อ" }
```

### GET `/purchase-orders/suggestions/{jobId}`
Response `[]PurchaseOrderSuggestionDTO`: `part_id`, `part_name`, `part_code`, `suggested_qty`, `current_stock`, `unit_price`, `from_quotation`

### GET `/purchase-orders/{id}/history`
Response `[]PurchaseOrderStatusHistoryResponse`: `id`, `po_header_id`, `from_status`, `to_status`, `changed_by`, `changed_at`, `reason`

---

## 3. Payment — `/api/payments` + Receipt `/api/receipts`

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| POST | `/payments` | `Create` | Bearer |
| GET | `/payments` | `List` | Bearer |
| POST | `/payments/search` | `Search` | Bearer |
| GET | `/payments/outstanding/{customerId}` | `GetOutstanding` | Bearer |
| GET | `/payments/history/{customerId}` | `GetHistory` | Bearer |
| GET | `/payments/{id}` | `Get` | Bearer |
| POST | `/payments/{id}/refund` | `Refund` | Bearer |
| PUT | `/payments/{id}/cancel` | `Cancel` | Bearer |
| GET | `/payments/invoice/{invoiceId}` | `GetByInvoice` | Bearer |
| GET | `/receipts/{id}` | `GetReceipt` | Bearer |
| GET | `/receipts/{id}/pdf` | `GetReceiptPDF` | Bearer |
| PUT | `/receipts/{id}/cancel` | `CancelReceipt` | Bearer |
| GET | `/receipts/payment/{paymentId}` | `GetReceiptByPayment` | Bearer |

> ⚠️ `/payments/{id}` ลงทะเบียนก่อน `/payments/invoice/{invoiceId}` — chi match ตามลำดับ

### POST `/payments`
Body (`PaymentRecordRequest`; `invoice_id`, `payment_method_id`, `amount`, `amount_received` required):
```json
{
  "invoice_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "payment_method_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "amount": 1500.00, "amount_received": 2000.00, "change_amount": 500.00,
  "currency": "THB", "exchange_rate": 1.0,
  "reference_number": "TRF-2026-0001", "bank_name": "ธนาคารกรุงเทพ",
  "cheque_number": "CHQ-001", "cheque_bank": "ธนาคารกสิกรไทย", "cheque_date": "2026-01-15",
  "notes": "ชำระค่าซ่อมบำรุง"
}
```
Handler sets `status: "PENDING"`, `received_by`/`user_id` จาก user ที่ login

### GET `/payments?page=1&per_page=10`
Response (`PaginatedPaymentResponse` — **field `payments`**):
```json
{ "data": { "payments": [ { "...PaymentResponse" } ], "total": 12, "page": 1, "per_page": 10, "total_pages": 2 }, "error": null, "is_success": true }
```

### GET `/payments/{id}`
Response `PaymentResponse`: `id`, `payment_no`, `invoice_id`, `job_id`, `customer_id`, `payment_date`, `payment_method_id`, `amount`, `amount_received`, `change_amount`, `currency`, `exchange_rate`, `status`, `reference_number`, `bank_name`, `cheque_number`, `cheque_bank`, `cheque_date`, `notes`, `received_by`, `approved_by`, `approved_at`, `refunded_amount`, `refunded_at`, `created_at`

### POST `/payments/search`
Body (`PaymentSearchRequest`, ทุก field optional):
```json
{ "customer_id": "...", "invoice_id": null, "status": "COMPLETED",
  "payment_method_id": null, "date_from": "2026-01-01", "date_to": "2026-12-31",
  "page": 1, "per_page": 10 }
```

### GET `/payments/outstanding/{customerId}`
Response `[]OutstandingBalanceResponse`: `invoice_id`, `invoice_total`, `amount_paid`, `outstanding_amount`, `last_payment_date`, `status`

### POST `/payments/{id}/refund`
Body (`RefundRequest`):
```json
{ "amount": 500.00, "reason": "ลูกค้าขอยกเลิกบริการ" }
```

### PUT `/payments/{id}/cancel`
Query: `reason` (**required** — ถ้าไม่มีคืน 422)

### GET `/receipts/{id}`
Response `ReceiptResponse`: `id`, `receipt_no`, `payment_id`, `invoice_id`, `customer_id`, `receipt_date`, `receipt_type`, `amount`, `amount_in_words_th`, `amount_in_words_en`, `currency`, `status`, `notes`, `issued_by`, `created_at`

### PUT `/receipts/{id}/cancel`
Query: `reason` (**required**)

---

## 4. Job — `/api/job`

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| GET | `/job` | `List` | Bearer |
| POST | `/job` | `Create` | Bearer |
| GET | `/job/{id}` | `GetByID` | Bearer |
| PUT | `/job/{id}` | `Update` | Bearer |
| DELETE | `/job/{id}` | `Delete` | Bearer |
| PUT | `/job/{id}/status` | `ChangeStatus` | Bearer |
| GET | `/job/{id}/history` | `GetStatusHistory` | Bearer |
| GET | `/job/{id}/report` | `GetReport` | Bearer |
| GET | `/job/{id}/pdf` | `GetPDF` | Bearer |
| GET | `/job/{id}/picking/pdf` | `GetPickingPDF` | Bearer |
| GET | `/job/{id}/delivery/pdf` | `GetDeliveryPDF` | Bearer |
| POST | `/job/{id}/services` | `AddService` | Bearer |
| GET | `/job/{id}/services` | `GetServices` | Bearer |
| POST | `/job/{id}/parts` | `AddPart` | Bearer |
| GET | `/job/{id}/parts` | `GetParts` | Bearer |

> ⚠️ Job ใช้ **camelCase** JSON

### POST `/job`
Body (`JobCreate`; `jobNo`, `customerId`, `carId`, `mechanicId` required):
```json
{
  "jobNo": "JOB-2026-0001",
  "customerId": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "carId": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "mechanicId": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "symptom": "เครื่องยนต์สั่น", "diagnosisNote": "เช็คหัวเทียน",
  "mileage": 50000, "estimatedCost": 3000, "priority": "NORMAL"
}
```
`priority` default "NORMAL"

### GET `/job?page=1&per_page=10`
Response (`PaginatedJobResponse` — **field `jobs`**)

### GET `/job/{id}`
Response `JobResponse` (camelCase): `id`, `jobNo`, `customerId`, `carId`, `mechanicId`, `status`, `startDate`, `endDate`, `symptom`, `diagnosisNote`, `mileage`, `estimatedCost`, `actualCost`, `priority`, `userId`, `whitelabelId`, `createdAt`, `updatedAt`

### PUT `/job/{id}/status`
Body (`JobStatusChange`):
```json
{ "status": "COMPLETED", "reason": "ซ่อมเสร็จเรียบร้อย" }
```

### POST `/job/{id}/services`
Body (`JobServiceRequest`; `quantity` default 1 ถ้า 0):
```json
{ "serviceId": "...", "quantity": 1, "unitPrice": 500.0, "discount": 0, "note": "เปลี่ยนน้ำมันเครื่อง" }
```

### POST `/job/{id}/parts`
Body (`JobPartRequest`):
```json
{ "partId": "...", "quantity": 2, "unitPrice": 150.0, "discount": 0, "note": "ผ้าเบรกหน้า" }
```

---

## 5. Customer & Car — `/api/customer` `/api/car`

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| GET | `/customer` | `GetMulti` | Bearer |
| POST | `/customer` | `Create` | Bearer |
| GET | `/customer/{id}` | `Get` | Bearer |
| PUT | `/customer/{id}` | `Update` | Bearer |
| DELETE | `/customer/{id}` | `Delete` | Bearer |
| GET | `/car` | `ListCars` | Bearer |
| POST | `/car` | `CreateCar` | Bearer |
| GET | `/car/{id}` | `GetCar` | Bearer |
| PUT | `/car/{id}` | `UpdateCar` | Bearer |
| DELETE | `/car/{id}` | `DeleteCar` | Bearer |

> ⚠️ Customer/Car ใช้ **camelCase** JSON

### POST `/customer`
Body (`CustomerCreate`; `customerCode`, `fullName`, `phoneNumber` required):
```json
{
  "customerCode": "CUST001", "fullName": "สมชาย ใจดี", "displayName": "คุณสมชาย",
  "customerType": "INDIVIDUAL", "taxId": "1234567890123", "email": "somchai@example.com",
  "phoneNumber": "0812345678", "secondaryPhone": null, "address": "123/4",
  "province": "กรุงเทพฯ", "city": "คลองเตย", "district": "คลองเตย",
  "postalCode": "10110", "country": "Thailand",
  "contactPerson": null, "contactPhone": null, "notes": null
}
```
Handler sets `user_id` จาก user ที่ login

### GET `/customer?page=1&per_page=10`
**Scoping:** ถ้า `user.IsSuperUser` → ลูกค้าทั้งหมด; ไม่ใช่ → เฉพาะของตัวเอง
Response (`PaginatedCustomersResponse` — **field `customers`**)

### GET `/customer/{id}`
Response `CustomerResponse` (camelCase): `id`, `customerCode`, `fullName`, `displayName`, `customerType`, `status`, `taxId`, `email`, `phoneNumber`, `secondaryPhone`, `address`, `province`, `city`, `district`, `postalCode`, `country`, `contactPerson`, `contactPhone`, `notes`, `lastVisitDate`, `totalVisitCount`, `totalSpent`, `userId`, `whitelabelId`, `createdAt`, `updatedAt`

### POST `/car`
Body (`CarCreate`; `customerId`, `licensePlate`, `brand`, `model` required):
```json
{
  "customerId": "...", "licensePlate": "กข1234", "province": "กรุงเทพฯ",
  "brand": "Toyota", "model": "Camry", "subModel": null, "year": 2020, "color": "ขาว",
  "engineNumber": null, "chassisNumber": null, "fuelType": "Gasoline",
  "transmissionType": "Automatic", "engineCc": 2500, "seatingCapacity": 5,
  "mileage": 30000, "lastServiceDate": null, "nextServiceMileage": null, "notes": null
}
```

### GET `/car?customerId&page&per_page`
Query: `customerId` (optional filter), `page`, `per_page`
Response (`PaginatedCarsResponse` — **field `cars`**)

### GET `/car/{id}`
Response `CarResponse`: `id`, `customerId`, `licensePlate`, `province`, `brand`, `model`, `subModel`, `year`, `color`, `engineNumber`, `chassisNumber`, `fuelType`, `transmissionType`, `engineCc`, `seatingCapacity`, `mileage`, `lastServiceDate`, `nextServiceMileage`, `notes`, `isActive`, `userId`, `whitelabelId`, `createdAt`, `updatedAt`

---

## หมายเหตุ
- **Pagination field ต่างกัน**: quotation/PO ใช้ `items`, payment ใช้ `payments`, job ใช้ `jobs`, customer ใช้ `customers`, car ใช้ `cars`
- รองรับ `page`/`per_page` หรือ `offset`/`limit` (ยกเว้น payment.List)
- `maxPerPage = 100`
- **Cancel** ที่อ่าน query `reason` (payment.cancel, receipt.cancel) → ต้องส่ง `reason`
- path param `{id}` เป็น UUID (invalid → 400)
- PDF endpoints คืน `application/pdf` inline (ไม่ผ่าน JSON envelope)
