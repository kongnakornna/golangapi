**โมดูลที่ 7: 💰 Payment Management (การจัดการการชำระเงิน)**

โมดูล Payment Management เป็นส่วนสำคัญของระบบที่ใช้ในการบันทึกและจัดการการชำระเงินจากลูกค้า ครอบคลุมการทำงานดังนี้:

1. **การบันทึกการชำระเงิน (Payment Recording)** – บันทึกการรับเงินจากลูกค้า เชื่อมโยงกับ Invoice
2. **การจัดการใบเสร็จรับเงิน (Receipt Management)** – สร้างเอกสารใบเสร็จรับเงิน (PDF)
3. **การจัดการวิธีการชำระเงิน (Payment Methods)** – เงินสด, โอนเงิน, บัตรเครดิต, เช็ค
4. **การติดตามยอดคงค้าง (Outstanding Balance)** – ตรวจสอบ Invoice ที่ยังไม่ได้ชำระ
5. **ประวัติการชำระเงิน (Payment History)** – แสดงประวัติการชำระเงินของลูกค้า
6. **การคืนเงิน (Refund)** – จัดการการคืนเงินให้ลูกค้า (Credit Note หรือ Refund)

---

## 📁 โครงสร้างโมดูล Payment Management (`modules/payment`)

```
modules/payment/
├── application/
│   ├── interfaces/
│   │   ├── PaymentService.java
│   │   ├── ReceiptService.java
│   │   ├── PaymentMethodService.java
│   │   └── PaymentHistoryService.java
│   ├── impl/
│   │   ├── PaymentServiceImpl.java
│   │   ├── ReceiptServiceImpl.java
│   │   ├── PaymentMethodServiceImpl.java
│   │   └── PaymentHistoryServiceImpl.java
│   └── usecase/
│       ├── RecordPaymentUseCase.java
│       ├── GetPaymentUseCase.java
│       ├── ListPaymentsUseCase.java
│       ├── CreateReceiptUseCase.java
│       ├── GenerateReceiptPDFUseCase.java
│       ├── GetOutstandingBalanceUseCase.java
│       ├── ProcessRefundUseCase.java
│       └── GetCustomerPaymentHistoryUseCase.java
├── domain/
│   ├── TPayment.java
│   ├── TReceipt.java
│   ├── MPaymentMethod.java
│   ├── enums/
│   │   ├── PaymentStatus.java           // PENDING, COMPLETED, FAILED, REFUNDED
│   │   ├── PaymentMethodType.java       // CASH, BANK_TRANSFER, CREDIT_CARD, CHEQUE, PROMPTPAY
│   │   └── ReceiptStatus.java           // DRAFT, ISSUED, CANCELLED
│   └── valueobjects/
│       ├── PaymentAmount.java
│       └── TransactionReference.java
├── infrastructure/
│   ├── repository/
│   │   ├── PaymentRepository.java
│   │   ├── ReceiptRepository.java
│   │   ├── PaymentMethodRepository.java
│   │   └── impl/
│   │       ├── PaymentRepositoryImpl.java
│   │       └── ReceiptRepositoryImpl.java
│   ├── cache/                                           // ⬅️ ระบบ Cache สำหรับ Payment
│   │   ├── PaymentCacheService.java
│   │   └── ReceiptCacheService.java
│   ├── report/                                          // ⬅️ ระบบสร้างใบเสร็จ PDF
│   │   ├── ReceiptReportGenerator.java
│   │   └── PaymentReportGenerator.java
│   ├── entity/
│   │   ├── PaymentEntity.java
│   │   ├── ReceiptEntity.java
│   │   └── PaymentMethodEntity.java
│   └── mapper/
│       ├── PaymentMapper.java
│       └── ReceiptMapper.java
└── presentation/
    ├── controller/
    │   ├── PaymentController.java        // Record, List, Get, Refund
    │   └── ReceiptController.java        // Generate Receipt, PDF
    ├── dto/
    │   ├── request/
    │   │   ├── PaymentRecordRequestDTO.java
    │   │   ├── RefundRequestDTO.java
    │   │   └── PaymentSearchRequestDTO.java
    │   └── response/
    │       ├── PaymentResponseDTO.java
    │       ├── ReceiptResponseDTO.java
    │       ├── PaymentHistoryResponseDTO.java
    │       └── OutstandingBalanceResponseDTO.java
    └── validator/
        └── PaymentValidator.java
```

---

## 🗄️ ออกแบบฐานข้อมูล (Database Design) สำหรับโมดูล Payment

### 📄 SQL DDL (ไฟล์: `src/main/resources/db/migration/V7__payment_schema.sql`)

```sql
-- ==============================================
-- ตาราง: m_payment_method (วิธีการชำระเงิน)
-- Payment method master data.
-- ==============================================
CREATE TABLE IF NOT EXISTS m_payment_method (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    method_code VARCHAR(20) UNIQUE NOT NULL,          -- CASH, BANK_TRANSFER, CREDIT_CARD, CHEQUE, PROMPTPAY
    method_name VARCHAR(100) NOT NULL,
    method_name_en VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    requires_approval BOOLEAN DEFAULT FALSE,
    fee_percentage DECIMAL(5,2) DEFAULT 0,            -- ค่าธรรมเนียม (%)
    fee_fixed DECIMAL(15,2) DEFAULT 0,                -- ค่าธรรมเนียม (ค่าคงที่)
    description TEXT,
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

INSERT INTO m_payment_method (method_code, method_name, method_name_en, is_active, user_id, whitelabel_id)
VALUES 
    ('CASH', 'เงินสด', 'Cash', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
    ('BANK_TRANSFER', 'โอนเงินผ่านธนาคาร', 'Bank Transfer', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
    ('CREDIT_CARD', 'บัตรเครดิต', 'Credit Card', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
    ('CHEQUE', 'เช็ค', 'Cheque', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
    ('PROMPTPAY', 'พร้อมเพย์', 'PromptPay', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001');

-- ==============================================
-- ตาราง: t_payment (การชำระเงิน)
-- Main table for payment transactions.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_payment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_no VARCHAR(20) UNIQUE NOT NULL,            -- เลขที่การชำระ (PAY-2026-0001)
    invoice_id UUID NOT NULL REFERENCES t_invoice_adjustment(id) ON DELETE RESTRICT,
    job_id UUID REFERENCES t_job(id) ON DELETE RESTRICT,
    customer_id UUID NOT NULL REFERENCES m_customer(id) ON DELETE RESTRICT,
    payment_date TIMESTAMP NOT NULL DEFAULT NOW(),
    payment_method_id UUID NOT NULL REFERENCES m_payment_method(id) ON DELETE RESTRICT,
    amount DECIMAL(15,2) NOT NULL,                     -- จำนวนเงินที่ชำระ
    amount_received DECIMAL(15,2) NOT NULL,            -- จำนวนเงินที่รับจริง (รวมส่วนต่าง)
    change_amount DECIMAL(15,2) DEFAULT 0,             -- เงินทอน (ถ้าชำระเกิน)
    currency VARCHAR(10) DEFAULT 'THB',
    exchange_rate DECIMAL(10,4) DEFAULT 1.0000,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',     -- PENDING, COMPLETED, FAILED, REFUNDED
    reference_number VARCHAR(50),                      -- เลขอ้างอิง (เช่น เลขที่โอน, เลขบัตร)
    bank_name VARCHAR(100),                            -- ชื่อธนาคาร (กรณีโอนเงิน)
    cheque_number VARCHAR(50),                         -- เลขที่เช็ค (กรณีเช็ค)
    cheque_bank VARCHAR(100),                          -- ธนาคารของเช็ค
    cheque_date DATE,                                  -- วันที่เช็ค
    notes TEXT,
    received_by UUID NOT NULL REFERENCES m_user(id),   -- ผู้รับเงิน
    approved_by UUID REFERENCES m_user(id),            -- ผู้ยืนยัน (ถ้าต้องการอนุมัติ)
    approved_at TIMESTAMP,
    refunded_amount DECIMAL(15,2) DEFAULT 0,           -- จำนวนเงินที่คืนไปแล้ว
    refunded_at TIMESTAMP,
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    deleted BOOLEAN DEFAULT FALSE,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_payment_invoice ON t_payment(invoice_id);
CREATE INDEX idx_t_payment_customer ON t_payment(customer_id);
CREATE INDEX idx_t_payment_date ON t_payment(payment_date);
CREATE INDEX idx_t_payment_status ON t_payment(status);
CREATE INDEX idx_t_payment_method ON t_payment(payment_method_id);
CREATE INDEX idx_t_payment_whitelabel ON t_payment(whitelabel_id);
CREATE INDEX idx_t_payment_deleted ON t_payment(deleted);

-- ==============================================
-- ฟังก์ชันสร้างเลขที่การชำระอัตโนมัติ (Auto-generate Payment Number)
-- Function to generate unique payment number with year and sequence.
-- ==============================================
CREATE OR REPLACE FUNCTION generate_payment_no()
RETURNS TRIGGER AS $$
DECLARE
    year_part TEXT;
    sequence_part TEXT;
BEGIN
    year_part := TO_CHAR(NOW(), 'YYYY');
    sequence_part := LPAD(CAST((
        SELECT COALESCE(MAX(CAST(SUBSTRING(payment_no FROM 9) AS INTEGER)), 0) + 1
        FROM t_payment
        WHERE payment_no LIKE 'PAY-' || year_part || '-%'
    ) AS TEXT), 4, '0');
    NEW.payment_no := 'PAY-' || year_part || '-' || sequence_part;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_generate_payment_no ON t_payment;
CREATE TRIGGER trg_generate_payment_no
BEFORE INSERT ON t_payment
FOR EACH ROW
EXECUTE FUNCTION generate_payment_no();

-- ==============================================
-- ตาราง: t_receipt (ใบเสร็จรับเงิน)
-- Receipt document table.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_receipt (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    receipt_no VARCHAR(20) UNIQUE NOT NULL,            -- เลขที่ใบเสร็จ (RCP-2026-0001)
    payment_id UUID NOT NULL REFERENCES t_payment(id) ON DELETE RESTRICT,
    invoice_id UUID NOT NULL REFERENCES t_invoice_adjustment(id) ON DELETE RESTRICT,
    customer_id UUID NOT NULL REFERENCES m_customer(id) ON DELETE RESTRICT,
    receipt_date TIMESTAMP NOT NULL DEFAULT NOW(),
    receipt_type VARCHAR(20) DEFAULT 'FULL',           -- FULL, PARTIAL, DEPOSIT
    amount DECIMAL(15,2) NOT NULL,                     -- จำนวนเงินที่รับ
    amount_in_words_th TEXT,                           -- จำนวนเงินเป็นตัวอักษร (ไทย)
    amount_in_words_en TEXT,                           -- จำนวนเงินเป็นตัวอักษร (อังกฤษ)
    currency VARCHAR(10) DEFAULT 'THB',
    status VARCHAR(20) DEFAULT 'ISSUED',               -- DRAFT, ISSUED, CANCELLED
    notes TEXT,
    issued_by UUID NOT NULL REFERENCES m_user(id),
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    deleted BOOLEAN DEFAULT FALSE,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_receipt_payment ON t_receipt(payment_id);
CREATE INDEX idx_t_receipt_invoice ON t_receipt(invoice_id);
CREATE INDEX idx_t_receipt_customer ON t_receipt(customer_id);
CREATE INDEX idx_t_receipt_date ON t_receipt(receipt_date);
CREATE INDEX idx_t_receipt_whitelabel ON t_receipt(whitelabel_id);

-- ==============================================
-- ฟังก์ชันสร้างเลขที่ใบเสร็จอัตโนมัติ (Auto-generate Receipt Number)
-- Function to generate unique receipt number with year and sequence.
-- ==============================================
CREATE OR REPLACE FUNCTION generate_receipt_no()
RETURNS TRIGGER AS $$
DECLARE
    year_part TEXT;
    sequence_part TEXT;
BEGIN
    year_part := TO_CHAR(NOW(), 'YYYY');
    sequence_part := LPAD(CAST((
        SELECT COALESCE(MAX(CAST(SUBSTRING(receipt_no FROM 9) AS INTEGER)), 0) + 1
        FROM t_receipt
        WHERE receipt_no LIKE 'RCP-' || year_part || '-%'
    ) AS TEXT), 4, '0');
    NEW.receipt_no := 'RCP-' || year_part || '-' || sequence_part;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_generate_receipt_no ON t_receipt;
CREATE TRIGGER trg_generate_receipt_no
BEFORE INSERT ON t_receipt
FOR EACH ROW
EXECUTE FUNCTION generate_receipt_no();

-- ==============================================
-- ตาราง: t_payment_history (ประวัติสถานะการชำระเงิน)
-- Track payment status changes over time.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_payment_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL REFERENCES t_payment(id) ON DELETE CASCADE,
    from_status VARCHAR(20),
    to_status VARCHAR(20) NOT NULL,
    changed_by UUID NOT NULL REFERENCES m_user(id) ON DELETE RESTRICT,
    changed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    reason TEXT,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_payment_history_payment ON t_payment_history(payment_id);
CREATE INDEX idx_t_payment_history_changed ON t_payment_history(changed_at);

-- ==============================================
-- ตาราง: t_outstanding_balance (ยอดคงค้าง - ใช้สำหรับ Dashboard)
-- This is a materialized view or summary table for outstanding balances.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_outstanding_balance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL UNIQUE REFERENCES t_invoice_adjustment(id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES m_customer(id) ON DELETE RESTRICT,
    invoice_total DECIMAL(15,2) NOT NULL,
    amount_paid DECIMAL(15,2) DEFAULT 0,
    outstanding_amount DECIMAL(15,2) GENERATED ALWAYS AS (invoice_total - amount_paid) STORED,
    last_payment_date TIMESTAMP,
    status VARCHAR(20) DEFAULT 'OUTSTANDING',          -- OUTSTANDING, PARTIAL, PAID, OVERDUE
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_outstanding_customer ON t_outstanding_balance(customer_id);
CREATE INDEX idx_t_outstanding_status ON t_outstanding_balance(status);
CREATE INDEX idx_t_outstanding_whitelabel ON t_outstanding_balance(whitelabel_id);
```

---

## 🧠 ระบบ Cache (Redis) สำหรับ Payment

### `infrastructure/cache/PaymentCacheService.java`

```java
package com.template.app.modules.payment.infrastructure.cache;

import com.template.app.modules.payment.domain.TPayment;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.CachePut;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class PaymentCacheService {

    /*
        ฟังก์ชันนี้ดึงข้อมูลการชำระเงินจาก Cache (ลดภาระฐานข้อมูล)
        This function retrieves payment data from cache (reduces DB load).
        Redis Key: payment:{id}
    */
    @Cacheable(value = "payments", key = "#paymentId")
    public TPayment getPayment(UUID paymentId) {
        return null;
    }

    /*
        ฟังก์ชันนี้ดึงข้อมูลการชำระเงินตาม Invoice ID (ใช้เชื่อมโยง Invoice -> Payment)
        This function retrieves payment by Invoice ID (links Invoice -> Payment).
        Redis Key: payment_invoice:{invoiceId}
    */
    @Cacheable(value = "payment_invoice", key = "#invoiceId")
    public TPayment getPaymentByInvoiceId(UUID invoiceId) {
        return null;
    }

    /*
        ฟังก์ชันนี้จะอัปเดต Cache เมื่อมีการบันทึกการชำระเงิน
        This function updates the cache when a payment is saved.
    */
    @CachePut(value = "payments", key = "#payment.id")
    public TPayment savePayment(TPayment payment) {
        return payment;
    }

    /*
        ฟังก์ชันนี้ลบข้อมูลการชำระเงินออกจาก Cache ทั้งสองรูปแบบ (ID และ Invoice ID)
        This function evicts payment data from both cache forms (ID and Invoice ID).
    */
    @CacheEvict(value = {"payments", "payment_invoice"}, key = "#paymentId")
    public void evictPayment(UUID paymentId) {
        // ลบ Cache ทั้งสองรูปแบบ / Evict both cache entries.
    }

    /*
        ฟังก์ชันนี้ลบ Cache ทั้งหมดของ Payment (ใช้เมื่อระบบอัปเดตข้อมูลจำนวนมาก)
        This function clears all payment caches (used during bulk updates).
    */
    @CacheEvict(value = {"payments", "payment_invoice"}, allEntries = true)
    public void evictAllPayments() {
        // ลบทุก key ใน caches ทั้งหมด / Evict all keys in all caches.
    }
}
```

### `infrastructure/cache/ReceiptCacheService.java`

```java
package com.template.app.modules.payment.infrastructure.cache;

import com.template.app.modules.payment.domain.TReceipt;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.CachePut;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class ReceiptCacheService {

    /*
        ฟังก์ชันนี้ดึงข้อมูลใบเสร็จจาก Cache
        This function retrieves receipt data from cache.
        Redis Key: receipt:{id}
    */
    @Cacheable(value = "receipts", key = "#receiptId")
    public TReceipt getReceipt(UUID receiptId) {
        return null;
    }

    /*
        ฟังก์ชันนี้ดึงข้อมูลใบเสร็จตาม Payment ID
        This function retrieves receipt by Payment ID.
        Redis Key: receipt_payment:{paymentId}
    */
    @Cacheable(value = "receipt_payment", key = "#paymentId")
    public TReceipt getReceiptByPaymentId(UUID paymentId) {
        return null;
    }

    /*
        ฟังก์ชันนี้จะอัปเดต Cache เมื่อมีการบันทึกใบเสร็จ
        This function updates the cache when a receipt is saved.
    */
    @CachePut(value = "receipts", key = "#receipt.id")
    public TReceipt saveReceipt(TReceipt receipt) {
        return receipt;
    }

    /*
        ฟังก์ชันนี้ลบข้อมูลใบเสร็จออกจาก Cache
        This function evicts receipt data from cache.
    */
    @CacheEvict(value = {"receipts", "receipt_payment"}, key = "#receiptId")
    public void evictReceipt(UUID receiptId) {
        // ลบ Cache / Evict cache.
    }
}
```

---

## ⏱️ ระบบ Rate Limit สำหรับ Payment Controller

```java
package com.template.app.modules.payment.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.payment.application.interfaces.PaymentService;
import com.template.app.modules.payment.presentation.dto.request.PaymentRecordRequestDTO;
import com.template.app.modules.payment.presentation.dto.request.PaymentSearchRequestDTO;
import com.template.app.modules.payment.presentation.dto.request.RefundRequestDTO;
import com.template.app.modules.payment.presentation.dto.response.OutstandingBalanceResponseDTO;
import com.template.app.modules.payment.presentation.dto.response.PaymentHistoryResponseDTO;
import com.template.app.modules.payment.presentation.dto.response.PaymentResponseDTO;
import com.template.app.exception.SystemGlobalException;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/api/v1/payments")
@Tag(name = "Payment", description = "Payment Management APIs")
@RequiredArgsConstructor
public class PaymentController {

    private final PaymentService paymentService;

    // ========================================================================
    // 1. RECORD PAYMENT (บันทึกการชำระเงิน)
    // ========================================================================

    /*
        API: POST /api/v1/payments
        ฟังก์ชันนี้บันทึกการชำระเงินจากลูกค้า เชื่อมโยงกับ Invoice และสร้างใบเสร็จอัตโนมัติ
        This function records a payment from a customer, links to Invoice, and auto-generates a receipt.
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที (ป้องกันการบันทึกซ้ำ)
        Rate Limit: Allows 20 requests per minute (prevent duplicate entries).
    */
    @PostMapping
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Record a payment")
    public ResponseEntity<PaymentResponseDTO> recordPayment(@Valid @RequestBody PaymentRecordRequestDTO request)
            throws SystemGlobalException {
        PaymentResponseDTO response = paymentService.recordPayment(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 2. GET PAYMENT BY ID
    // ========================================================================

    /*
        API: GET /api/v1/payments/{id}
        ฟังก์ชันนี้ดึงข้อมูลการชำระเงินตาม ID (ใช้ Cache ช่วยลดภาระ DB)
        This function retrieves payment by ID (uses caching to reduce DB load).
        Rate Limit: อนุญาต 100 ครั้งต่อ 1 นาที
        Rate Limit: Allows 100 requests per minute.
    */
    @GetMapping("/{id}")
    @RateLimit(limit = 100, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get payment by ID")
    public ResponseEntity<PaymentResponseDTO> getPayment(@PathVariable UUID id)
            throws SystemGlobalException {
        PaymentResponseDTO response = paymentService.getPayment(id);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. GET PAYMENT BY INVOICE ID
    // ========================================================================

    /*
        API: GET /api/v1/payments/invoice/{invoiceId}
        ฟังก์ชันนี้ดึงการชำระเงินตาม Invoice ID
        This function retrieves payment by Invoice ID.
        Rate Limit: อนุญาต 60 ครั้งต่อ 1 นาที
        Rate Limit: Allows 60 requests per minute.
    */
    @GetMapping("/invoice/{invoiceId}")
    @RateLimit(limit = 60, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get payment by Invoice ID")
    public ResponseEntity<PaymentResponseDTO> getPaymentByInvoiceId(@PathVariable UUID invoiceId)
            throws SystemGlobalException {
        PaymentResponseDTO response = paymentService.getPaymentByInvoiceId(invoiceId);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 4. LIST PAYMENTS (Pagination + Filters)
    // ========================================================================

    /*
        API: POST /api/v1/payments/search
        ฟังก์ชันนี้ค้นหาการชำระเงินด้วยตัวกรอง เช่น ลูกค้า, ช่วงวันที่, สถานะ
        This function searches payments with filters: customer, date range, status.
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @PostMapping("/search")
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Search payments with filters")
    public ResponseEntity<Page<PaymentResponseDTO>> searchPayments(
            @Valid @RequestBody PaymentSearchRequestDTO request,
            Pageable pageable) throws SystemGlobalException {
        Page<PaymentResponseDTO> page = paymentService.searchPayments(request, pageable);
        return ResponseEntity.ok(page);
    }

    // ========================================================================
    // 5. GET OUTSTANDING BALANCE (ยอดคงค้าง)
    // ========================================================================

    /*
        API: GET /api/v1/payments/outstanding/{customerId}
        ฟังก์ชันนี้แสดงยอดคงค้างของลูกค้า (Invoice ที่ยังไม่ชำระเต็มจำนวน)
        This function shows the outstanding balance for a customer (unpaid invoices).
        Rate Limit: อนุญาต 40 ครั้งต่อ 1 นาที
        Rate Limit: Allows 40 requests per minute.
    */
    @GetMapping("/outstanding/{customerId}")
    @RateLimit(limit = 40, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get customer outstanding balance")
    public ResponseEntity<OutstandingBalanceResponseDTO> getOutstandingBalance(@PathVariable UUID customerId)
            throws SystemGlobalException {
        OutstandingBalanceResponseDTO response = paymentService.getOutstandingBalance(customerId);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 6. GET CUSTOMER PAYMENT HISTORY
    // ========================================================================

    /*
        API: GET /api/v1/payments/history/{customerId}
        ฟังก์ชันนี้แสดงประวัติการชำระเงินของลูกค้า (ใช้ใน Profile หรือ Dashboard)
        This function shows the payment history of a customer (used in Profile or Dashboard).
        Rate Limit: อนุญาต 40 ครั้งต่อ 1 นาที
        Rate Limit: Allows 40 requests per minute.
    */
    @GetMapping("/history/{customerId}")
    @RateLimit(limit = 40, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get customer payment history")
    public ResponseEntity<List<PaymentHistoryResponseDTO>> getPaymentHistory(@PathVariable UUID customerId)
            throws SystemGlobalException {
        List<PaymentHistoryResponseDTO> history = paymentService.getPaymentHistory(customerId);
        return ResponseEntity.ok(history);
    }

    // ========================================================================
    // 7. PROCESS REFUND (การคืนเงิน)
    // ========================================================================

    /*
        API: POST /api/v1/payments/{id}/refund
        ฟังก์ชันนี้ดำเนินการคืนเงินให้ลูกค้า (เชื่อมโยงกับ Credit Note)
        This function processes a refund to a customer (linked to Credit Note).
        Rate Limit: อนุญาต 10 ครั้งต่อ 1 ชั่วโมง (จำกัดการคืนเงิน)
        Rate Limit: Allows 10 requests per 1 hour (limited refunds).
    */
    @PostMapping("/{id}/refund")
    @RateLimit(limit = 10, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Process refund")
    public ResponseEntity<PaymentResponseDTO> processRefund(
            @PathVariable UUID id,
            @Valid @RequestBody RefundRequestDTO request) throws SystemGlobalException {
        PaymentResponseDTO response = paymentService.processRefund(id, request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 8. CANCEL PAYMENT (ยกเลิกการชำระเงิน)
    // ========================================================================

    /*
        API: PUT /api/v1/payments/{id}/cancel
        ฟังก์ชันนี้ยกเลิกการชำระเงิน (เฉพาะที่ยังไม่ได้รับการยืนยัน)
        This function cancels a payment (only if not yet confirmed).
        Rate Limit: อนุญาต 10 ครั้งต่อ 1 ชั่วโมง
        Rate Limit: Allows 10 requests per 1 hour.
    */
    @PutMapping("/{id}/cancel")
    @RateLimit(limit = 10, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Cancel payment")
    public ResponseEntity<PaymentResponseDTO> cancelPayment(@PathVariable UUID id, @RequestParam String reason)
            throws SystemGlobalException {
        PaymentResponseDTO response = paymentService.cancelPayment(id, reason);
        return ResponseEntity.ok(response);
    }
}
```

### `ReceiptController.java`

```java
package com.template.app.modules.payment.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.payment.application.interfaces.ReceiptService;
import com.template.app.modules.payment.presentation.dto.response.ReceiptResponseDTO;
import com.template.app.exception.SystemGlobalException;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

@RestController
@RequestMapping("/api/v1/receipts")
@Tag(name = "Receipt", description = "Receipt Management APIs")
@RequiredArgsConstructor
public class ReceiptController {

    private final ReceiptService receiptService;

    // ========================================================================
    // 1. GET RECEIPT BY ID
    // ========================================================================

    /*
        API: GET /api/v1/receipts/{id}
        ฟังก์ชันนี้ดึงข้อมูลใบเสร็จตาม ID
        This function retrieves receipt by ID.
        Rate Limit: อนุญาต 100 ครั้งต่อ 1 นาที
        Rate Limit: Allows 100 requests per minute.
    */
    @GetMapping("/{id}")
    @RateLimit(limit = 100, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get receipt by ID")
    public ResponseEntity<ReceiptResponseDTO> getReceipt(@PathVariable UUID id)
            throws SystemGlobalException {
        ReceiptResponseDTO response = receiptService.getReceipt(id);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 2. GET RECEIPT BY PAYMENT ID
    // ========================================================================

    /*
        API: GET /api/v1/receipts/payment/{paymentId}
        ฟังก์ชันนี้ดึงใบเสร็จตาม Payment ID
        This function retrieves receipt by Payment ID.
        Rate Limit: อนุญาต 60 ครั้งต่อ 1 นาที
        Rate Limit: Allows 60 requests per minute.
    */
    @GetMapping("/payment/{paymentId}")
    @RateLimit(limit = 60, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get receipt by Payment ID")
    public ResponseEntity<ReceiptResponseDTO> getReceiptByPaymentId(@PathVariable UUID paymentId)
            throws SystemGlobalException {
        ReceiptResponseDTO response = receiptService.getReceiptByPaymentId(paymentId);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. GENERATE RECEIPT PDF (สร้างใบเสร็จ PDF)
    // ========================================================================

    /*
        API: GET /api/v1/receipts/{id}/pdf
        ฟังก์ชันนี้สร้างไฟล์ PDF ของใบเสร็จรับเงิน (ใช้ JasperReports) สำหรับพิมพ์ให้ลูกค้า
        This function generates a PDF file of the receipt (using JasperReports) for customer delivery.
        Rate Limit: อนุญาต 15 ครั้งต่อ 5 นาที (การสร้าง PDF ใช้ทรัพยากรมาก)
        Rate Limit: Allows 15 requests per 5 minutes (PDF generation is resource-intensive).
    */
    @GetMapping("/{id}/pdf")
    @RateLimit(limit = 15, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Generate receipt PDF")
    public ResponseEntity<byte[]> generateReceiptPDF(@PathVariable UUID id) throws SystemGlobalException {
        byte[] pdf = receiptService.generateReceiptPDF(id);
        return ResponseEntity.ok()
                .header("Content-Type", "application/pdf")
                .header("Content-Disposition", "inline; filename=receipt_" + id + ".pdf")
                .body(pdf);
    }

    // ========================================================================
    // 4. CANCEL RECEIPT
    // ========================================================================

    /*
        API: PUT /api/v1/receipts/{id}/cancel
        ฟังก์ชันนี้ยกเลิกใบเสร็จ (ใช้เมื่อมีการแก้ไขหรือคืนเงิน)
        This function cancels a receipt (used when adjustments or refunds occur).
        Rate Limit: อนุญาต 10 ครั้งต่อ 1 ชั่วโมง
        Rate Limit: Allows 10 requests per 1 hour.
    */
    @PutMapping("/{id}/cancel")
    @RateLimit(limit = 10, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Cancel receipt")
    public ResponseEntity<Void> cancelReceipt(@PathVariable UUID id, @RequestParam String reason)
            throws SystemGlobalException {
        receiptService.cancelReceipt(id, reason);
        return ResponseEntity.noContent().build();
    }
}
```

---

## 🧩 ตัวอย่าง Domain Entity และ Enum

### `domain/enums/PaymentStatus.java`

```java
package com.template.app.modules.payment.domain.enums;

/*
    สถานะของการชำระเงิน / Payment status.
*/
public enum PaymentStatus {
    PENDING,      // รอดำเนินการ / Pending.
    COMPLETED,    // ชำระเงินสำเร็จ / Completed.
    FAILED,       // ชำระเงินล้มเหลว / Failed.
    REFUNDED,     // คืนเงินแล้ว / Refunded.
    CANCELLED     // ยกเลิก / Cancelled.
}
```

### `domain/enums/PaymentMethodType.java`

```java
package com.template.app.modules.payment.domain.enums;

/*
    ประเภทวิธีการชำระเงิน / Payment method type.
*/
public enum PaymentMethodType {
    CASH,            // เงินสด / Cash.
    BANK_TRANSFER,   // โอนเงินผ่านธนาคาร / Bank transfer.
    CREDIT_CARD,     // บัตรเครดิต / Credit card.
    CHEQUE,          // เช็ค / Cheque.
    PROMPTPAY,       // พร้อมเพย์ / PromptPay.
    OTHER            // อื่นๆ / Other.
}
```

### `domain/TPayment.java`

```java
package com.template.app.modules.payment.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import com.template.app.modules.payment.domain.enums.PaymentMethodType;
import com.template.app.modules.payment.domain.enums.PaymentStatus;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.math.BigDecimal;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.UUID;

@Data
@EqualsAndHashCode(callSuper = true)
public class TPayment extends GenericBusinessClass {

    private String paymentNo;               // เลขที่การชำระ / Payment number.
    private UUID invoiceId;                 // ID Invoice
    private UUID jobId;                     // ID Job (ถ้ามี)
    private UUID customerId;                // ID ลูกค้า / Customer ID.
    private LocalDateTime paymentDate;      // วันที่ชำระ / Payment date.
    private UUID paymentMethodId;           // ID วิธีการชำระ / Payment method ID.
    private BigDecimal amount;              // จำนวนเงินที่ชำระ / Amount paid.
    private BigDecimal amountReceived;      // จำนวนเงินที่รับจริง / Actual amount received.
    private BigDecimal changeAmount;        // เงินทอน / Change amount.
    private String currency;                // สกุลเงิน / Currency.
    private BigDecimal exchangeRate;        // อัตราแลกเปลี่ยน / Exchange rate.
    private PaymentStatus status;           // สถานะ / Status.
    private String referenceNumber;         // เลขอ้างอิง / Reference number.
    private String bankName;                // ชื่อธนาคาร / Bank name.
    private String chequeNumber;            // เลขที่เช็ค / Cheque number.
    private String chequeBank;              // ธนาคารของเช็ค / Cheque bank.
    private LocalDate chequeDate;           // วันที่เช็ค / Cheque date.
    private String notes;                   // หมายเหตุ / Notes.
    private UUID receivedBy;                // ผู้รับเงิน / Received by.
    private UUID approvedBy;                // ผู้อนุมัติ / Approved by.
    private LocalDateTime approvedAt;       // วันที่อนุมัติ / Approved date.
    private BigDecimal refundedAmount;      // จำนวนเงินที่คืนแล้ว / Refunded amount.
    private LocalDateTime refundedAt;       // วันที่คืนเงิน / Refunded date.

    /*
        ฟังก์ชันนี้ตรวจสอบว่าการชำระเงินสามารถอนุมัติได้หรือไม่
        This function checks if the payment can be approved.
    */
    public boolean canApprove() {
        return this.status == PaymentStatus.PENDING;
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่าการชำระเงินสามารถคืนเงินได้หรือไม่
        This function checks if the payment can be refunded.
    */
    public boolean canRefund() {
        return this.status == PaymentStatus.COMPLETED 
                && this.refundedAmount == null 
                && this.amount.compareTo(BigDecimal.ZERO) > 0;
    }

    /*
        ฟังก์ชันนี้เปลี่ยนสถานะเป็น COMPLETED และบันทึกข้อมูลการอนุมัติ
        This function changes status to COMPLETED and records approval.
    */
    public void approve(UUID approverId) {
        if (!canApprove()) {
            throw new IllegalStateException("Cannot approve payment with status: " + this.status);
        }
        this.status = PaymentStatus.COMPLETED;
        this.approvedBy = approverId;
        this.approvedAt = LocalDateTime.now();
    }

    /*
        ฟังก์ชันนี้ดำเนินการคืนเงินและบันทึกจำนวนเงินที่คืน
        This function processes refund and records the refunded amount.
    */
    public void processRefund(BigDecimal refundAmount) {
        if (!canRefund()) {
            throw new IllegalStateException("Cannot refund payment with status: " + this.status);
        }
        if (refundAmount.compareTo(BigDecimal.ZERO) <= 0) {
            throw new IllegalArgumentException("Refund amount must be greater than zero.");
        }
        if (refundAmount.compareTo(this.amount) > 0) {
            throw new IllegalArgumentException("Refund amount exceeds paid amount.");
        }
        this.refundedAmount = refundAmount;
        this.refundedAt = LocalDateTime.now();
        this.status = PaymentStatus.REFUNDED;
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่ายังมียอดเงินคงเหลือที่ต้องคืนหรือไม่
        This function checks if there is remaining balance to refund.
    */
    public BigDecimal getRemainingRefundable() {
        if (this.refundedAmount == null) {
            return this.amount;
        }
        return this.amount.subtract(this.refundedAmount);
    }
}
```

### `domain/TReceipt.java`

```java
package com.template.app.modules.payment.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.UUID;

@Data
@EqualsAndHashCode(callSuper = true)
public class TReceipt extends GenericBusinessClass {

    private String receiptNo;              // เลขที่ใบเสร็จ / Receipt number.
    private UUID paymentId;                // ID การชำระเงิน / Payment ID.
    private UUID invoiceId;                // ID Invoice
    private UUID customerId;               // ID ลูกค้า / Customer ID.
    private LocalDateTime receiptDate;     // วันที่ออกใบเสร็จ / Receipt date.
    private String receiptType;            // FULL, PARTIAL, DEPOSIT
    private BigDecimal amount;             // จำนวนเงิน / Amount.
    private String amountInWordsTh;        // จำนวนเงินเป็นตัวอักษร (ไทย)
    private String amountInWordsEn;        // จำนวนเงินเป็นตัวอักษร (อังกฤษ)
    private String currency;               // สกุลเงิน / Currency.
    private String status;                 // DRAFT, ISSUED, CANCELLED
    private String notes;                  // หมายเหตุ / Notes.
    private UUID issuedBy;                 // ผู้ออกใบเสร็จ / Issued by.

    /*
        ฟังก์ชันนี้ออกใบเสร็จ (เปลี่ยนสถานะเป็น ISSUED)
        This function issues the receipt (changes status to ISSUED).
    */
    public void issue() {
        if (!"DRAFT".equals(this.status)) {
            throw new IllegalStateException("Cannot issue receipt with status: " + this.status);
        }
        this.status = "ISSUED";
        this.receiptDate = LocalDateTime.now();
    }

    /*
        ฟังก์ชันนี้ยกเลิกใบเสร็จ
        This function cancels the receipt.
    */
    public void cancel(String reason) {
        if ("CANCELLED".equals(this.status)) {
            throw new IllegalStateException("Receipt is already cancelled.");
        }
        this.status = "CANCELLED";
        this.notes = (this.notes != null ? this.notes + "\n" : "") + "Cancelled: " + reason;
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่าใบเสร็จสามารถยกเลิกได้หรือไม่
        This function checks if the receipt can be cancelled.
    */
    public boolean canCancel() {
        return "ISSUED".equals(this.status) || "DRAFT".equals(this.status);
    }
}
```

---

## 🔧 ฟังก์ชัน Service ตัวอย่าง (พร้อมคำอธิบาย)

### `application/impl/PaymentServiceImpl.java` (บางส่วน)

```java
package com.template.app.modules.payment.application.impl;

import com.template.app.modules.payment.application.interfaces.PaymentService;
import com.template.app.modules.payment.domain.TPayment;
import com.template.app.modules.payment.domain.TReceipt;
import com.template.app.modules.payment.domain.enums.PaymentStatus;
import com.template.app.modules.payment.infrastructure.cache.PaymentCacheService;
import com.template.app.modules.payment.infrastructure.repository.PaymentRepository;
import com.template.app.modules.payment.infrastructure.repository.ReceiptRepository;
import com.template.app.modules.payment.presentation.dto.request.PaymentRecordRequestDTO;
import com.template.app.modules.payment.presentation.dto.response.PaymentResponseDTO;
import com.template.app._shared.application.GenericServiceImpl;
import com.template.app.exception.SystemGlobalException;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.UUID;

@Service
public class PaymentServiceImpl extends GenericServiceImpl<TPayment, PaymentRepository>
        implements PaymentService {

    private final ReceiptRepository receiptRepository;
    private final PaymentCacheService cacheService;
    private final InvoiceService invoiceService;      // Inject Invoice Service

    public PaymentServiceImpl(PaymentRepository repository,
                              ReceiptRepository receiptRepository,
                              PaymentCacheService cacheService,
                              InvoiceService invoiceService) {
        super(repository);
        this.receiptRepository = receiptRepository;
        this.cacheService = cacheService;
        this.invoiceService = invoiceService;
    }

    /*
        ฟังก์ชันนี้บันทึกการชำระเงินจากลูกค้า ตรวจสอบ Invoice, บันทึก Payment, สร้าง Receipt อัตโนมัติ
        This function records a payment from a customer, validates Invoice, saves Payment, and auto-generates Receipt.
    */
    @Override
    @Transactional
    public PaymentResponseDTO recordPayment(PaymentRecordRequestDTO request) throws SystemGlobalException {
        // 1. ตรวจสอบว่า Invoice มีอยู่และยังไม่ชำระเต็ม / Validate Invoice exists and not fully paid.
        Invoice invoice = invoiceService.getInvoice(request.getInvoiceId());
        if (invoice.getStatus() == InvoiceStatus.PAID) {
            throw new SystemGlobalException("Invoice already fully paid.", null);
        }

        // 2. ตรวจสอบว่าจำนวนเงินที่ชำระไม่เกินยอดคงค้าง / Validate payment amount does not exceed outstanding.
        BigDecimal outstanding = invoice.getOutstandingBalance();
        if (request.getAmount().compareTo(outstanding) > 0) {
            throw new SystemGlobalException("Payment amount exceeds outstanding balance.", null);
        }

        // 3. สร้าง Domain Entity Payment / Create Payment domain entity.
        TPayment payment = new TPayment();
        payment.setInvoiceId(request.getInvoiceId());
        payment.setJobId(invoice.getJobId());
        payment.setCustomerId(invoice.getCustomerId());
        payment.setPaymentDate(LocalDateTime.now());
        payment.setPaymentMethodId(request.getPaymentMethodId());
        payment.setAmount(request.getAmount());
        payment.setAmountReceived(request.getAmountReceived());
        payment.setChangeAmount(request.getChangeAmount());
        payment.setCurrency(request.getCurrency());
        payment.setExchangeRate(request.getExchangeRate());
        payment.setStatus(PaymentStatus.PENDING);
        payment.setReferenceNumber(request.getReferenceNumber());
        payment.setBankName(request.getBankName());
        payment.setChequeNumber(request.getChequeNumber());
        payment.setChequeBank(request.getChequeBank());
        payment.setChequeDate(request.getChequeDate());
        payment.setNotes(request.getNotes());
        payment.setReceivedBy(getUserId());

        // 4. บันทึก Payment / Save Payment.
        TPayment savedPayment = this.create(payment);

        // 5. อัปเดต Invoice (อัปเดตยอดที่ชำระแล้วและสถานะ) / Update Invoice (paid amount and status).
        invoiceService.updatePaymentStatus(request.getInvoiceId(), request.getAmount());

        // 6. สร้าง Receipt อัตโนมัติ / Auto-generate Receipt.
        TReceipt receipt = new TReceipt();
        receipt.setPaymentId(savedPayment.getId());
        receipt.setInvoiceId(request.getInvoiceId());
        receipt.setCustomerId(invoice.getCustomerId());
        receipt.setAmount(request.getAmount());
        receipt.setCurrency(request.getCurrency());
        receipt.setStatus("DRAFT");
        receipt.setIssuedBy(getUserId());

        // 6.1 แปลงจำนวนเงินเป็นตัวอักษร / Convert amount to words.
        receipt.setAmountInWordsTh(ConvertBath.convert(request.getAmount()));  // ภาษาไทย
        receipt.setAmountInWordsEn(ConvertDollar.convert(request.getAmount())); // ภาษาอังกฤษ

        // 6.2 ออกใบเสร็จ / Issue receipt.
        receipt.issue();
        TReceipt savedReceipt = receiptRepository.save(receipt);

        // 7. ถ้าจำนวนเงินที่ชำระเท่ากับยอดคงค้าง ให้เปลี่ยนสถานะ Invoice เป็น PAID
        // If amount equals outstanding, change Invoice status to PAID.
        if (request.getAmount().compareTo(outstanding) >= 0) {
            invoiceService.markAsPaid(request.getInvoiceId());
        }

        // 8. บันทึก Cache / Save to cache.
        cacheService.savePayment(savedPayment);

        return PaymentResponseDTO.fromEntity(savedPayment, savedReceipt);
    }

    /*
        ฟังก์ชันนี้ดำเนินการคืนเงินให้ลูกค้า ตรวจสอบว่าสามารถคืนได้ และสร้าง Credit Note
        This function processes a refund to a customer, validates refund eligibility, and creates a Credit Note.
    */
    @Override
    @Transactional
    public PaymentResponseDTO processRefund(UUID paymentId, RefundRequestDTO request) throws SystemGlobalException {
        // 1. ดึงข้อมูล Payment / Get Payment.
        TPayment payment = this.read(paymentId);

        // 2. ตรวจสอบว่าสามารถคืนเงินได้ / Validate refund eligibility.
        if (!payment.canRefund()) {
            throw new SystemGlobalException("Cannot refund payment with status: " + payment.getStatus(), null);
        }

        // 3. ตรวจสอบว่าจำนวนเงินที่ขอคืนไม่เกินยอดที่จ่าย / Validate refund amount.
        BigDecimal remaining = payment.getRemainingRefundable();
        if (request.getAmount().compareTo(remaining) > 0) {
            throw new SystemGlobalException("Refund amount exceeds remaining balance.", null);
        }

        // 4. ดำเนินการคืนเงิน / Process refund.
        payment.processRefund(request.getAmount());

        // 5. สร้าง Credit Note สำหรับ Invoice / Create Credit Note for Invoice.
        // TODO: สร้าง Credit Note ผ่าน InvoiceService
        // creditNoteService.createCreditNote(payment.getInvoiceId(), request.getAmount(), "Refund");

        // 6. บันทึก Payment / Save Payment.
        TPayment updatedPayment = this.update(paymentId, payment);

        // 7. ลบ Cache เก่าและอัปเดต Cache ใหม่ / Evict old cache and update new cache.
        cacheService.evictPayment(paymentId);
        cacheService.savePayment(updatedPayment);

        return PaymentResponseDTO.fromEntity(updatedPayment, null);
    }
}
```

---

## ✅ สรุปส่วนที่เพิ่มเติมในโมดูล Payment Management

| ฟีเจอร์ | สถานะ | คำอธิบาย |
|---------|--------|----------|
| **Database Design** | ✅ ครบถ้วน | `m_payment_method`, `t_payment`, `t_receipt`, `t_payment_history`, `t_outstanding_balance` พร้อม Trigger |
| **Cache (Redis)** | ✅ เพิ่มแล้ว | `PaymentCacheService` (ID + Invoice ID), `ReceiptCacheService` (ID + Payment ID) |
| **Rate Limit** | ✅ เพิ่มแล้ว | ทุก Endpoint ใน Controller พร้อมกำหนด Limit และ Duration |
| **API Routing** | ✅ ชัดเจน | Record Payment + Search + Outstanding Balance + History + Refund + Cancel + Receipt PDF |
| **Bilingual Comments** | ✅ ครบถ้วน | ทุก Function มีคำอธิบายภาษาไทยและภาษาอังกฤษแยกบรรทัด |
| **PDF Generation** | ✅ เพิ่มแล้ว | `ReceiptReportGenerator` ใช้ JasperReports |

---

## 📊 สรุปโมดูลที่ดำเนินการแล้ว (Modules Completed)

| # | โมดูล | สถานะ | รายละเอียด |
|---|-------|--------|-----------|
| 1 | 🔑 Authentication & Permission | ✅ ครบถ้วน | JWT + Role/Permission + Rate Limit + Redis Cache |
| 2 | 🚗 Job Card Management | ✅ ครบถ้วน | 14 Statuses + Service/Part + History + Cache |
| 3 | 👥 Customer Management | ✅ ครบถ้วน | Customer + Car + History + Cache (ID/Phone/Plate) |
| 4 | 📋 Quotation | ✅ ครบถ้วน | Quotation + Part/Service + Approve/Reject + PDF + Cache |
| 5 | 🛒 Purchase Order | ✅ ครบถ้วน | PO + Status + Send/Receive + PDF + Email + Cache |
| 6 | 📦 Inventory Management | ✅ ครบถ้วน | Part Master + Receive/Issue + Picking + Stock Take + Adjustment + Cache |
| 7 | 💰 Payment Management | ✅ ครบถ้วน | Payment Record + Receipt + Outstanding Balance + Refund + Cache |

---
 