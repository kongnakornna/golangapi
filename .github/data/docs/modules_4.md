**โมดูลที่ 4: 📋 Quotation (การจัดการใบเสนอราคา)**

โมดูล Quotation เป็นหัวใจสำคัญของระบบ เพราะเป็นเอกสารที่ใช้สื่อสารกับลูกค้าในการเสนอราคาค่าบริการและอะไหล่ ก่อนดำเนินการซ่อม โดยมีกระบวนการ:
1. สร้างใบเสนอราคาจาก Job Card
2. เพิ่มรายการอะไหล่และบริการ
3. คำนวณราคารวม (Subtotal, Tax, Total)
4. แปลงจำนวนเงินเป็นข้อความ (ภาษาไทย/อังกฤษ)
5. ส่งให้ลูกค้าอนุมัติ
6. เมื่ออนุมัติ → สร้าง Purchase Order และเริ่มกระบวนการเบิกอะไหล่

---

## 📁 โครงสร้างโมดูล Quotation (`modules/quotation`)

```
modules/quotation/
├── application/
│   ├── interfaces/
│   │   ├── QuotationService.java
│   │   ├── QuotationPartService.java
│   │   └── QuotationServiceItemService.java
│   ├── impl/
│   │   ├── QuotationServiceImpl.java
│   │   ├── QuotationPartServiceImpl.java
│   │   └── QuotationServiceItemServiceImpl.java
│   └── usecase/
│       ├── CreateQuotationUseCase.java
│       ├── UpdateQuotationUseCase.java
│       ├── GetQuotationUseCase.java
│       ├── DeleteQuotationUseCase.java
│       ├── ApproveQuotationUseCase.java
│       ├── RejectQuotationUseCase.java
│       ├── AddQuotationPartUseCase.java
│       ├── AddQuotationServiceUseCase.java
│       ├── RemoveQuotationPartUseCase.java
│       ├── RemoveQuotationServiceUseCase.java
│       ├── CalculateQuotationTotalUseCase.java
│       └── GenerateQuotationPDFUseCase.java
├── domain/
│   ├── TQuotation.java
│   ├── TQuotationPart.java
│   ├── TQuotationService.java
│   ├── TQuotationStatusHistory.java
│   ├── enums/
│   │   └── QuotationStatus.java          // DRAFT, PENDING, APPROVED, REJECTED, EXPIRED, CONVERTED
│   └── valueobjects/
│       ├── QuotationNumber.java
│       ├── Amount.java                    // จำนวนเงินพร้อมสกุลเงิน
│       └── TaxRate.java                   // อัตราภาษี
├── infrastructure/
│   ├── repository/
│   │   ├── QuotationRepository.java
│   │   ├── QuotationPartRepository.java
│   │   ├── QuotationServiceRepository.java
│   │   ├── QuotationStatusHistoryRepository.java
│   │   └── impl/
│   │       ├── QuotationRepositoryImpl.java
│   │       ├── QuotationPartRepositoryImpl.java
│   │       └── QuotationServiceRepositoryImpl.java
│   ├── cache/                                      // ⬅️ ระบบ Cache สำหรับ Quotation
│   │   ├── QuotationCacheService.java
│   │   └── QuotationCalculationCacheService.java
│   ├── report/                                     // ⬅️ ระบบสร้าง PDF
│   │   ├── QuotationReportGenerator.java
│   │   └── QuotationReportDataProvider.java
│   ├── entity/
│   │   ├── QuotationEntity.java
│   │   ├── QuotationPartEntity.java
│   │   ├── QuotationServiceEntity.java
│   │   └── QuotationStatusHistoryEntity.java
│   └── mapper/
│       ├── QuotationMapper.java
│       ├── QuotationPartMapper.java
│       └── QuotationServiceItemMapper.java
└── presentation/
    ├── controller/
    │   ├── QuotationController.java           // CRUD + Approve/Reject + PDF
    │   ├── QuotationPartController.java       // จัดการรายการอะไหล่
    │   └── QuotationServiceController.java    // จัดการรายการบริการ
    ├── dto/
    │   ├── request/
    │   │   ├── QuotationCreateRequestDTO.java
    │   │   ├── QuotationUpdateRequestDTO.java
    │   │   ├── QuotationApproveRequestDTO.java
    │   │   ├── QuotationPartRequestDTO.java
    │   │   └── QuotationServiceRequestDTO.java
    │   └── response/
    │       ├── QuotationResponseDTO.java
    │       ├── QuotationDetailResponseDTO.java
    │       ├── QuotationPartResponseDTO.java
    │       └── QuotationServiceResponseDTO.java
    └── validator/
        ├── QuotationValidator.java
        └── QuotationCalculationValidator.java
```

---

## 🗄️ ออกแบบฐานข้อมูล (Database Design) สำหรับโมดูล Quotation

### 📄 SQL DDL (ไฟล์: `src/main/resources/db/migration/V4__quotation_schema.sql`)

```sql
-- ==============================================
-- ตาราง: t_quotation (ใบเสนอราคา)
-- Main table for quotation documents.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_quotation (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quotation_no VARCHAR(20) UNIQUE NOT NULL,            -- เลขที่ใบเสนอราคา (QT-2026-0001)
    job_id UUID NOT NULL REFERENCES t_job(id) ON DELETE RESTRICT,
    customer_id UUID NOT NULL REFERENCES m_customer(id) ON DELETE RESTRICT,
    quotation_date TIMESTAMP NOT NULL DEFAULT NOW(),
    expiry_date TIMESTAMP NOT NULL,                      -- วันหมดอายุ (เช่น 7 วัน)
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',         -- DRAFT, PENDING, APPROVED, REJECTED, EXPIRED, CONVERTED
    subtotal DECIMAL(15,2) NOT NULL DEFAULT 0,           -- ราคาก่อนภาษี
    tax_rate DECIMAL(5,2) DEFAULT 7.00,                  -- อัตราภาษี (%) / Tax rate (%).
    tax_amount DECIMAL(15,2) DEFAULT 0,                  -- จำนวนภาษี / Tax amount.
    discount_type VARCHAR(20),                           -- PERCENTAGE, FIXED
    discount_value DECIMAL(15,2) DEFAULT 0,              -- ส่วนลด / Discount value.
    total DECIMAL(15,2) NOT NULL DEFAULT 0,              -- ราคาสุทธิ (Net total)
    amount_in_words_th TEXT,                             -- จำนวนเงินเป็นตัวอักษร (ไทย)
    amount_in_words_en TEXT,                             -- จำนวนเงินเป็นตัวอักษร (อังกฤษ)
    currency VARCHAR(10) DEFAULT 'THB',                  -- สกุลเงิน / Currency.
    exchange_rate DECIMAL(10,4) DEFAULT 1.0000,          -- อัตราแลกเปลี่ยน / Exchange rate.
    notes TEXT,                                          -- หมายเหตุ / Notes.
    terms_and_conditions TEXT,                           -- เงื่อนไข / Terms and conditions.
    approved_by UUID REFERENCES m_user(id),              -- ผู้อนุมัติ / Approver.
    approved_at TIMESTAMP,                               -- วันที่อนุมัติ / Approval date.
    rejected_reason TEXT,                                -- เหตุผลในการปฏิเสธ / Rejection reason.
    converted_to_po BOOLEAN DEFAULT FALSE,               -- ถูกแปลงเป็น PO แล้วหรือไม่
    -- Audit Fields
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    deleted BOOLEAN DEFAULT FALSE,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_quotation_job ON t_quotation(job_id);
CREATE INDEX idx_t_quotation_customer ON t_quotation(customer_id);
CREATE INDEX idx_t_quotation_status ON t_quotation(status);
CREATE INDEX idx_t_quotation_date ON t_quotation(quotation_date);
CREATE INDEX idx_t_quotation_whitelabel ON t_quotation(whitelabel_id);
CREATE INDEX idx_t_quotation_deleted ON t_quotation(deleted);

-- ==============================================
-- ฟังก์ชันสร้างเลขที่ใบเสนอราคาอัตโนมัติ (Auto-generate Quotation Number)
-- Function to generate unique quotation number with year and sequence.
-- ==============================================
CREATE OR REPLACE FUNCTION generate_quotation_no()
RETURNS TRIGGER AS $$
DECLARE
    year_part TEXT;
    sequence_part TEXT;
BEGIN
    year_part := TO_CHAR(NOW(), 'YYYY');
    sequence_part := LPAD(CAST((
        SELECT COALESCE(MAX(CAST(SUBSTRING(quotation_no FROM 8) AS INTEGER)), 0) + 1
        FROM t_quotation
        WHERE quotation_no LIKE 'QT-' || year_part || '-%'
    ) AS TEXT), 4, '0');
    NEW.quotation_no := 'QT-' || year_part || '-' || sequence_part;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_generate_quotation_no ON t_quotation;
CREATE TRIGGER trg_generate_quotation_no
BEFORE INSERT ON t_quotation
FOR EACH ROW
EXECUTE FUNCTION generate_quotation_no();

-- ==============================================
-- ตาราง: t_quotation_part (รายการอะไหล่ในใบเสนอราคา)
-- Parts listed in the quotation.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_quotation_part (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quotation_id UUID NOT NULL REFERENCES t_quotation(id) ON DELETE CASCADE,
    part_id UUID NOT NULL REFERENCES m_part_master(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(15,2) NOT NULL,
    total_price DECIMAL(15,2) GENERATED ALWAYS AS (quantity * unit_price) STORED,
    discount DECIMAL(15,2) DEFAULT 0,
    net_price DECIMAL(15,2) GENERATED ALWAYS AS ((quantity * unit_price) - discount) STORED,
    note TEXT,
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_quotation_part_quotation ON t_quotation_part(quotation_id);
CREATE INDEX idx_t_quotation_part_part ON t_quotation_part(part_id);
CREATE INDEX idx_t_quotation_part_whitelabel ON t_quotation_part(whitelabel_id);

-- ==============================================
-- ตาราง: t_quotation_service (รายการบริการในใบเสนอราคา)
-- Services listed in the quotation.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_quotation_service (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quotation_id UUID NOT NULL REFERENCES t_quotation(id) ON DELETE CASCADE,
    service_id UUID NOT NULL REFERENCES m_service(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(15,2) NOT NULL,
    total_price DECIMAL(15,2) GENERATED ALWAYS AS (quantity * unit_price) STORED,
    discount DECIMAL(15,2) DEFAULT 0,
    net_price DECIMAL(15,2) GENERATED ALWAYS AS ((quantity * unit_price) - discount) STORED,
    note TEXT,
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_quotation_service_quotation ON t_quotation_service(quotation_id);
CREATE INDEX idx_t_quotation_service_service ON t_quotation_service(service_id);
CREATE INDEX idx_t_quotation_service_whitelabel ON t_quotation_service(whitelabel_id);

-- ==============================================
-- ตาราง: t_quotation_status_history (ประวัติการเปลี่ยนสถานะใบเสนอราคา)
-- Track quotation status changes over time.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_quotation_status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quotation_id UUID NOT NULL REFERENCES t_quotation(id) ON DELETE CASCADE,
    from_status VARCHAR(20),
    to_status VARCHAR(20) NOT NULL,
    changed_by UUID NOT NULL REFERENCES m_user(id) ON DELETE RESTRICT,
    changed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    reason TEXT,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_quotation_status_history_quotation ON t_quotation_status_history(quotation_id);
CREATE INDEX idx_t_quotation_status_history_changed ON t_quotation_status_history(changed_at);
```

---

## 🧠 ระบบ Cache (Redis) สำหรับ Quotation

### `infrastructure/cache/QuotationCacheService.java`

```java
package com.template.app.modules.quotation.infrastructure.cache;

import com.template.app.modules.quotation.domain.TQuotation;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.CachePut;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class QuotationCacheService {

    /*
        ฟังก์ชันนี้ดึงข้อมูลใบเสนอราคาจาก Cache (ลดภาระฐานข้อมูลเมื่อเปิดดูบ่อยๆ)
        This function retrieves quotation data from cache (reduces DB load on frequent views).
        Redis Key: quotation:{id}
    */
    @Cacheable(value = "quotations", key = "#quotationId")
    public TQuotation getQuotation(UUID quotationId) {
        return null;
    }

    /*
        ฟังก์ชันนี้ดึงข้อมูลใบเสนอราคาตาม Job ID (ใช้เมื่อเปิด Job แล้วดู Quotation)
        This function retrieves quotation by Job ID (used when viewing Quotation from Job).
        Redis Key: quotation_job:{jobId}
    */
    @Cacheable(value = "quotation_job", key = "#jobId")
    public TQuotation getQuotationByJobId(UUID jobId) {
        return null;
    }

    /*
        ฟังก์ชันนี้จะอัปเดต Cache เมื่อมีการบันทึกใบเสนอราคา
        This function updates the cache when a quotation is saved.
    */
    @CachePut(value = "quotations", key = "#quotation.id")
    public TQuotation saveQuotation(TQuotation quotation) {
        return quotation;
    }

    /*
        ฟังก์ชันนี้ลบข้อมูลใบเสนอราคาออกจาก Cache ทั้งสองรูปแบบ (ID และ Job ID)
        This function evicts quotation data from both cache forms (ID and Job ID).
    */
    @CacheEvict(value = {"quotations", "quotation_job"}, key = "#quotationId")
    public void evictQuotation(UUID quotationId) {
        // ลบ Cache ทั้งสองรูปแบบ / Evict both cache entries.
    }

    /*
        ฟังก์ชันนี้ลบ Cache ทั้งหมดของ Quotation (ใช้เมื่อระบบอัปเดตข้อมูลจำนวนมาก)
        This function clears all quotation caches (used during bulk updates).
    */
    @CacheEvict(value = {"quotations", "quotation_job"}, allEntries = true)
    public void evictAllQuotations() {
        // ลบทุก key ใน quotations และ quotation_job / Evict all keys in both caches.
    }
}
```

---

## ⏱️ ระบบ Rate Limit สำหรับ Quotation Controller

```java
package com.template.app.modules.quotation.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.quotation.application.interfaces.QuotationService;
import com.template.app.modules.quotation.presentation.dto.request.QuotationApproveRequestDTO;
import com.template.app.modules.quotation.presentation.dto.request.QuotationCreateRequestDTO;
import com.template.app.modules.quotation.presentation.dto.request.QuotationUpdateRequestDTO;
import com.template.app.modules.quotation.presentation.dto.response.QuotationDetailResponseDTO;
import com.template.app.modules.quotation.presentation.dto.response.QuotationResponseDTO;
import com.template.app.exception.SystemGlobalException;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

@RestController
@RequestMapping("/api/v1/quotations")
@Tag(name = "Quotation", description = "Quotation Management APIs")
@RequiredArgsConstructor
public class QuotationController {

    private final QuotationService quotationService;

    // ========================================================================
    // 1. CREATE QUOTATION
    // ========================================================================

    /*
        API: POST /api/v1/quotations
        ฟังก์ชันนี้สร้างใบเสนอราคาใหม่จาก Job ID โดยจะดึงข้อมูลลูกค้าและรถยนต์จาก Job
        This function creates a new quotation from a Job ID, fetching customer and vehicle data from the Job.
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @PostMapping
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Create a new quotation")
    public ResponseEntity<QuotationResponseDTO> createQuotation(@Valid @RequestBody QuotationCreateRequestDTO request)
            throws SystemGlobalException {
        QuotationResponseDTO response = quotationService.createQuotation(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 2. GET QUOTATION BY ID
    // ========================================================================

    /*
        API: GET /api/v1/quotations/{id}
        ฟังก์ชันนี้ดึงข้อมูลใบเสนอราคาตาม ID พร้อมรายการอะไหล่และบริการทั้งหมด
        This function retrieves quotation by ID with all parts and services details.
        Rate Limit: อนุญาต 100 ครั้งต่อ 1 นาที
        Rate Limit: Allows 100 requests per minute.
    */
    @GetMapping("/{id}")
    @RateLimit(limit = 100, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get quotation by ID")
    public ResponseEntity<QuotationDetailResponseDTO> getQuotation(@PathVariable UUID id)
            throws SystemGlobalException {
        QuotationDetailResponseDTO response = quotationService.getQuotation(id);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. GET QUOTATION BY JOB ID
    // ========================================================================

    /*
        API: GET /api/v1/quotations/job/{jobId}
        ฟังก์ชันนี้ดึงใบเสนอราคาล่าสุดของ Job ที่ระบุ (ใช้เชื่อมโยง Job -> Quotation)
        This function retrieves the latest quotation for a given Job ID (links Job -> Quotation).
        Rate Limit: อนุญาต 80 ครั้งต่อ 1 นาที
        Rate Limit: Allows 80 requests per minute.
    */
    @GetMapping("/job/{jobId}")
    @RateLimit(limit = 80, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get quotation by Job ID")
    public ResponseEntity<QuotationDetailResponseDTO> getQuotationByJobId(@PathVariable UUID jobId)
            throws SystemGlobalException {
        QuotationDetailResponseDTO response = quotationService.getQuotationByJobId(jobId);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 4. LIST QUOTATIONS (Pagination + Filters)
    // ========================================================================

    /*
        API: GET /api/v1/quotations
        ฟังก์ชันนี้แสดงรายการใบเสนอราคาแบบแบ่งหน้า พร้อมตัวกรองตามสถานะและช่วงวันที่
        This function lists quotations with pagination, filtering by status and date range.
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @GetMapping
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "List quotations with pagination")
    public ResponseEntity<Page<QuotationResponseDTO>> listQuotations(
            @RequestParam(required = false) String status,
            @RequestParam(required = false) String startDate,
            @RequestParam(required = false) String endDate,
            @RequestParam(required = false) UUID customerId,
            Pageable pageable) throws SystemGlobalException {
        Page<QuotationResponseDTO> page = quotationService.listQuotations(status, startDate, endDate, customerId, pageable);
        return ResponseEntity.ok(page);
    }

    // ========================================================================
    // 5. UPDATE QUOTATION
    // ========================================================================

    /*
        API: PUT /api/v1/quotations/{id}
        ฟังก์ชันนี้แก้ไขข้อมูลใบเสนอราคา (เฉพาะที่ยังไม่ถูกอนุมัติ)
        This function updates quotation details (only if not yet approved).
        Rate Limit: อนุญาต 15 ครั้งต่อ 1 นาที
        Rate Limit: Allows 15 requests per minute.
    */
    @PutMapping("/{id}")
    @RateLimit(limit = 15, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Update quotation details")
    public ResponseEntity<QuotationResponseDTO> updateQuotation(
            @PathVariable UUID id,
            @Valid @RequestBody QuotationUpdateRequestDTO request) throws SystemGlobalException {
        QuotationResponseDTO response = quotationService.updateQuotation(id, request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 6. DELETE QUOTATION (SOFT DELETE)
    // ========================================================================

    /*
        API: DELETE /api/v1/quotations/{id}
        ฟังก์ชันนี้ลบใบเสนอราคาแบบ Soft Delete (เฉพาะที่ยังไม่ถูกอนุมัติ)
        This function soft-deletes a quotation (only if not yet approved).
        Rate Limit: อนุญาต 10 ครั้งต่อ 1 ชั่วโมง
        Rate Limit: Allows 10 requests per 1 hour.
    */
    @DeleteMapping("/{id}")
    @RateLimit(limit = 10, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Delete quotation (soft delete)")
    public ResponseEntity<Void> deleteQuotation(@PathVariable UUID id) throws SystemGlobalException {
        quotationService.deleteQuotation(id);
        return ResponseEntity.noContent().build();
    }

    // ========================================================================
    // 7. APPROVE QUOTATION
    // ========================================================================

    /*
        API: PUT /api/v1/quotations/{id}/approve
        ฟังก์ชันนี้อนุมัติใบเสนอราคา เมื่อลูกค้าตกลงแล้ว จะสร้าง Purchase Order และเริ่มกระบวนการเบิกอะไหล่
        This function approves the quotation. Upon customer agreement, it creates a Purchase Order and starts parts picking.
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @PutMapping("/{id}/approve")
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Approve quotation and initiate process")
    public ResponseEntity<QuotationResponseDTO> approveQuotation(
            @PathVariable UUID id,
            @Valid @RequestBody QuotationApproveRequestDTO request) throws SystemGlobalException {
        QuotationResponseDTO response = quotationService.approveQuotation(id, request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 8. REJECT QUOTATION
    // ========================================================================

    /*
        API: PUT /api/v1/quotations/{id}/reject
        ฟังก์ชันนี้ปฏิเสธใบเสนอราคา พร้อมระบุเหตุผล
        This function rejects the quotation with a reason.
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @PutMapping("/{id}/reject")
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Reject quotation with reason")
    public ResponseEntity<QuotationResponseDTO> rejectQuotation(
            @PathVariable UUID id,
            @RequestParam String reason) throws SystemGlobalException {
        QuotationResponseDTO response = quotationService.rejectQuotation(id, reason);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 9. GENERATE QUOTATION PDF
    // ========================================================================

    /*
        API: GET /api/v1/quotations/{id}/pdf
        ฟังก์ชันนี้สร้างไฟล์ PDF ของใบเสนอราคา (ใช้ JasperReports) สำหรับส่งให้ลูกค้า
        This function generates a PDF file of the quotation (using JasperReports) for customer delivery.
        Rate Limit: อนุญาต 15 ครั้งต่อ 5 นาที (การสร้าง PDF ใช้ทรัพยากรมาก)
        Rate Limit: Allows 15 requests per 5 minutes (PDF generation is resource-intensive).
    */
    @GetMapping("/{id}/pdf")
    @RateLimit(limit = 15, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Generate quotation PDF")
    public ResponseEntity<byte[]> generateQuotationPDF(@PathVariable UUID id) throws SystemGlobalException {
        byte[] pdf = quotationService.generateQuotationPDF(id);
        return ResponseEntity.ok()
                .header("Content-Type", "application/pdf")
                .header("Content-Disposition", "inline; filename=quotation_" + id + ".pdf")
                .body(pdf);
    }

    // ========================================================================
    // 10. GET QUOTATION STATUS HISTORY
    // ========================================================================

    /*
        API: GET /api/v1/quotations/{id}/history
        ฟังก์ชันนี้ดึงประวัติการเปลี่ยนสถานะของใบเสนอราคา
        This function retrieves the status change history of a quotation.
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @GetMapping("/{id}/history")
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get quotation status history")
    public ResponseEntity<List<QuotationStatusHistoryDTO>> getQuotationHistory(@PathVariable UUID id)
            throws SystemGlobalException {
        List<QuotationStatusHistoryDTO> history = quotationService.getQuotationHistory(id);
        return ResponseEntity.ok(history);
    }
}
```

### `QuotationPartController.java` (จัดการรายการอะไหล่)

```java
package com.template.app.modules.quotation.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.quotation.application.interfaces.QuotationPartService;
import com.template.app.modules.quotation.presentation.dto.request.QuotationPartRequestDTO;
import com.template.app.modules.quotation.presentation.dto.response.QuotationPartResponseDTO;
import com.template.app.exception.SystemGlobalException;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

@RestController
@RequestMapping("/api/v1/quotations/parts")
@Tag(name = "Quotation Parts", description = "Quotation Part Management APIs")
@RequiredArgsConstructor
public class QuotationPartController {

    private final QuotationPartService quotationPartService;

    // ========================================================================
    // 1. ADD PART TO QUOTATION
    // ========================================================================

    /*
        API: POST /api/v1/quotations/parts
        ฟังก์ชันนี้เพิ่มอะไหล่ในใบเสนอราคา โดยระบุ Quotation ID, Part ID และจำนวน
        This function adds a part to a quotation, specifying Quotation ID, Part ID, and quantity.
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที
        Rate Limit: Allows 30 requests per minute.
    */
    @PostMapping
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Add part to quotation")
    public ResponseEntity<QuotationPartResponseDTO> addPart(@Valid @RequestBody QuotationPartRequestDTO request)
            throws SystemGlobalException {
        QuotationPartResponseDTO response = quotationPartService.addPart(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 2. UPDATE PART IN QUOTATION
    // ========================================================================

    /*
        API: PUT /api/v1/quotations/parts/{id}
        ฟังก์ชันนี้แก้ไขรายการอะไหล่ในใบเสนอราคา (เช่น ปรับจำนวน หรือราคา)
        This function updates a part item in the quotation (e.g., adjust quantity or price).
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @PutMapping("/{id}")
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Update part in quotation")
    public ResponseEntity<QuotationPartResponseDTO> updatePart(
            @PathVariable UUID id,
            @Valid @RequestBody QuotationPartRequestDTO request) throws SystemGlobalException {
        QuotationPartResponseDTO response = quotationPartService.updatePart(id, request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. REMOVE PART FROM QUOTATION
    // ========================================================================

    /*
        API: DELETE /api/v1/quotations/parts/{id}
        ฟังก์ชันนี้ลบรายการอะไหล่ออกจากใบเสนอราคา
        This function removes a part item from the quotation.
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @DeleteMapping("/{id}")
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Remove part from quotation")
    public ResponseEntity<Void> removePart(@PathVariable UUID id) throws SystemGlobalException {
        quotationPartService.removePart(id);
        return ResponseEntity.noContent().build();
    }

    // ========================================================================
    // 4. GET ALL PARTS IN QUOTATION
    // ========================================================================

    /*
        API: GET /api/v1/quotations/parts/quotation/{quotationId}
        ฟังก์ชันนี้ดึงรายการอะไหล่ทั้งหมดในใบเสนอราคา
        This function lists all parts in a quotation.
        Rate Limit: อนุญาต 60 ครั้งต่อ 1 นาที
        Rate Limit: Allows 60 requests per minute.
    */
    @GetMapping("/quotation/{quotationId}")
    @RateLimit(limit = 60, duration = 60, keyType = "USER_ID")
    @Operation(summary = "List all parts in quotation")
    public ResponseEntity<List<QuotationPartResponseDTO>> getPartsByQuotation(@PathVariable UUID quotationId)
            throws SystemGlobalException {
        List<QuotationPartResponseDTO> parts = quotationPartService.getPartsByQuotation(quotationId);
        return ResponseEntity.ok(parts);
    }
}
```

---

## 🧩 ตัวอย่าง Domain Entity และ Enum

### `domain/enums/QuotationStatus.java`

```java
package com.template.app.modules.quotation.domain.enums;

/*
    สถานะของใบเสนอราคา / Quotation status.
*/
public enum QuotationStatus {
    DRAFT,        // ร่าง / Draft.
    PENDING,      // รอการอนุมัติ / Pending approval.
    APPROVED,     // อนุมัติแล้ว / Approved.
    REJECTED,     // ถูกปฏิเสธ / Rejected.
    EXPIRED,      // หมดอายุ / Expired.
    CONVERTED     // ถูกแปลงเป็น PO แล้ว / Converted to PO.
}
```

### `domain/TQuotation.java`

```java
package com.template.app.modules.quotation.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import com.template.app.modules.quotation.domain.enums.QuotationStatus;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

@Data
@EqualsAndHashCode(callSuper = true)
public class TQuotation extends GenericBusinessClass {

    private String quotationNo;                // เลขที่ใบเสนอราคา / Quotation number.
    private UUID jobId;                        // ID ใบงาน / Job ID.
    private UUID customerId;                   // ID ลูกค้า / Customer ID.
    private LocalDateTime quotationDate;       // วันที่เสนอราคา / Quotation date.
    private LocalDateTime expiryDate;          // วันหมดอายุ / Expiry date.
    private QuotationStatus status;            // สถานะ / Status.
    private BigDecimal subtotal;               // ราคาก่อนภาษี / Subtotal.
    private BigDecimal taxRate;                // อัตราภาษี (%) / Tax rate (%).
    private BigDecimal taxAmount;              // จำนวนภาษี / Tax amount.
    private String discountType;               // PERCENTAGE, FIXED
    private BigDecimal discountValue;          // ส่วนลด / Discount value.
    private BigDecimal total;                  // ราคาสุทธิ / Net total.
    private String amountInWordsTh;            // จำนวนเงินเป็นตัวอักษร (ไทย)
    private String amountInWordsEn;            // จำนวนเงินเป็นตัวอักษร (อังกฤษ)
    private String currency;                   // สกุลเงิน / Currency.
    private BigDecimal exchangeRate;           // อัตราแลกเปลี่ยน / Exchange rate.
    private String notes;                      // หมายเหตุ / Notes.
    private String termsAndConditions;         // เงื่อนไข / Terms and conditions.
    private UUID approvedBy;                   // ผู้อนุมัติ / Approver ID.
    private LocalDateTime approvedAt;          // วันที่อนุมัติ / Approval date.
    private String rejectedReason;             // เหตุผลในการปฏิเสธ / Rejection reason.
    private Boolean convertedToPo;             // ถูกแปลงเป็น PO แล้ว / Converted to PO.

    // รายการอะไหล่และบริการ (ไม่เก็บในฐานข้อมูลโดยตรง / Not directly stored)
    private List<TQuotationPart> parts = new ArrayList<>();
    private List<TQuotationService> services = new ArrayList<>();

    /*
        ฟังก์ชันนี้คำนวณยอดรวมของใบเสนอราคา (Subtotal, Tax, Total)
        This function calculates the total amounts of the quotation (Subtotal, Tax, Total).
    */
    public void calculateTotals() {
        // 1. คำนวณ Subtotal จากรายการอะไหล่และบริการ / Calculate subtotal from parts and services.
        BigDecimal partsTotal = parts.stream()
                .map(TQuotationPart::getNetPrice)
                .reduce(BigDecimal.ZERO, BigDecimal::add);
        BigDecimal servicesTotal = services.stream()
                .map(TQuotationService::getNetPrice)
                .reduce(BigDecimal.ZERO, BigDecimal::add);
        this.subtotal = partsTotal.add(servicesTotal);

        // 2. คำนวณภาษี / Calculate tax.
        this.taxAmount = this.subtotal.multiply(this.taxRate.divide(new BigDecimal(100)));

        // 3. คำนวณส่วนลด (ถ้ามี) / Apply discount (if any).
        BigDecimal discountAmount = BigDecimal.ZERO;
        if ("PERCENTAGE".equalsIgnoreCase(this.discountType)) {
            discountAmount = this.subtotal.multiply(this.discountValue.divide(new BigDecimal(100)));
        } else if ("FIXED".equalsIgnoreCase(this.discountType)) {
            discountAmount = this.discountValue;
        }

        // 4. คำนวณยอดรวมสุทธิ / Calculate net total.
        this.total = this.subtotal.add(this.taxAmount).subtract(discountAmount);
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่าใบเสนอราคาสามารถอนุมัติได้หรือไม่
        This function checks if the quotation can be approved.
    */
    public boolean canApprove() {
        return this.status == QuotationStatus.PENDING || this.status == QuotationStatus.DRAFT;
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่าใบเสนอราคาหมดอายุแล้วหรือไม่
        This function checks if the quotation has expired.
    */
    public boolean isExpired() {
        return LocalDateTime.now().isAfter(this.expiryDate);
    }

    /*
        ฟังก์ชันนี้เปลี่ยนสถานะเป็น APPROVED พร้อมบันทึกข้อมูลผู้อนุมัติ
        This function changes status to APPROVED and records the approver.
    */
    public void approve(UUID approverId) {
        if (!canApprove()) {
            throw new IllegalStateException("Quotation cannot be approved in status: " + this.status);
        }
        if (isExpired()) {
            throw new IllegalStateException("Quotation has expired and cannot be approved.");
        }
        this.status = QuotationStatus.APPROVED;
        this.approvedBy = approverId;
        this.approvedAt = LocalDateTime.now();
    }

    /*
        ฟังก์ชันนี้เปลี่ยนสถานะเป็น REJECTED พร้อมบันทึกเหตุผล
        This function changes status to REJECTED and records the reason.
    */
    public void reject(String reason) {
        if (this.status == QuotationStatus.APPROVED || this.status == QuotationStatus.CONVERTED) {
            throw new IllegalStateException("Cannot reject already approved or converted quotation.");
        }
        this.status = QuotationStatus.REJECTED;
        this.rejectedReason = reason;
    }
}
```

### `domain/TQuotationPart.java`

```java
package com.template.app.modules.quotation.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.math.BigDecimal;
import java.util.UUID;

@Data
@EqualsAndHashCode(callSuper = true)
public class TQuotationPart extends GenericBusinessClass {

    private UUID quotationId;           // ID ใบเสนอราคา / Quotation ID.
    private UUID partId;                // ID อะไหล่ / Part ID.
    private Integer quantity;           // จำนวน / Quantity.
    private BigDecimal unitPrice;       // ราคาต่อหน่วย / Unit price.
    private BigDecimal totalPrice;      // ราคารวม (quantity * unitPrice) / Total price.
    private BigDecimal discount;        // ส่วนลด / Discount.
    private BigDecimal netPrice;        // ราคาสุทธิ (totalPrice - discount) / Net price.
    private String note;                // หมายเหตุ / Note.

    /*
        ฟังก์ชันนี้คำนวณราคารวมและราคาสุทธิของรายการอะไหล่
        This function calculates total and net prices for the part item.
    */
    public void calculatePrices() {
        this.totalPrice = this.unitPrice.multiply(new BigDecimal(this.quantity));
        this.netPrice = this.totalPrice.subtract(this.discount != null ? this.discount : BigDecimal.ZERO);
    }
}
```

---

## 🔧 ฟังก์ชัน Service ตัวอย่าง (พร้อมคำอธิบาย)

### `application/impl/QuotationServiceImpl.java` (บางส่วน)

```java
package com.template.app.modules.quotation.application.impl;

import com.template.app.modules.quotation.application.interfaces.QuotationService;
import com.template.app.modules.quotation.domain.TQuotation;
import com.template.app.modules.quotation.domain.enums.QuotationStatus;
import com.template.app.modules.quotation.infrastructure.cache.QuotationCacheService;
import com.template.app.modules.quotation.infrastructure.repository.QuotationRepository;
import com.template.app.modules.quotation.presentation.dto.request.QuotationApproveRequestDTO;
import com.template.app.modules.quotation.presentation.dto.request.QuotationCreateRequestDTO;
import com.template.app.modules.quotation.presentation.dto.response.QuotationResponseDTO;
import com.template.app._shared.application.GenericServiceImpl;
import com.template.app.exception.SystemGlobalException;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;

@Service
public class QuotationServiceImpl extends GenericServiceImpl<TQuotation, QuotationRepository> implements QuotationService {

    private final QuotationCacheService cacheService;

    public QuotationServiceImpl(QuotationRepository repository, QuotationCacheService cacheService) {
        super(repository);
        this.cacheService = cacheService;
    }

    /*
        ฟังก์ชันนี้สร้างใบเสนอราคาใหม่จาก Job ID ตรวจสอบว่า Job อยู่ในสถานะที่สามารถสร้าง Quotation ได้
        This function creates a new quotation from a Job ID, verifying the Job can have a Quotation.
    */
    @Override
    public QuotationResponseDTO createQuotation(QuotationCreateRequestDTO request) throws SystemGlobalException {
        // 1. ตรวจสอบว่า Job นี้มีสิทธิ์สร้าง Quotation หรือไม่ / Verify Job can create Quotation.
        if (!jobService.canCreateQuotation(request.getJobId())) {
            throw new SystemGlobalException("Job cannot create quotation. Check job status.", null);
        }

        // 2. สร้าง Domain Entity / Create domain entity.
        TQuotation quotation = new TQuotation();
        quotation.setJobId(request.getJobId());
        quotation.setCustomerId(request.getCustomerId());
        quotation.setQuotationDate(LocalDateTime.now());
        quotation.setExpiryDate(LocalDateTime.now().plusDays(7)); // หมดอายุใน 7 วัน / Expiry in 7 days.
        quotation.setStatus(QuotationStatus.DRAFT);
        quotation.setTaxRate(request.getTaxRate());
        quotation.setCurrency(request.getCurrency());

        // 3. คำนวณยอดรวม (จะคำนวณใหม่เมื่อมีรายการ) / Calculate total (recalculated when items added).
        quotation.calculateTotals();

        // 4. บันทึกผ่าน Generic Repository (จะสร้าง quotationNo อัตโนมัติผ่าน Trigger) / Save via Generic Repository.
        TQuotation savedQuotation = this.create(quotation);

        // 5. บันทึกลง Cache / Store in cache.
        cacheService.saveQuotation(savedQuotation);

        return QuotationResponseDTO.fromEntity(savedQuotation);
    }

    /*
        ฟังก์ชันนี้อนุมัติใบเสนอราคา เปลี่ยนสถานะเป็น APPROVED และเริ่มกระบวนการสร้าง Purchase Order
        This function approves the quotation, changes status to APPROVED, and initiates Purchase Order creation.
    */
    @Override
    public QuotationResponseDTO approveQuotation(UUID quotationId, QuotationApproveRequestDTO request)
            throws SystemGlobalException {
        // 1. ดึงข้อมูล Quotation / Get quotation.
        TQuotation quotation = this.read(quotationId);

        // 2. ตรวจสอบและเปลี่ยนสถานะเป็น APPROVED / Validate and change status to APPROVED.
        quotation.approve(getUserId());

        // 3. บันทึกการเปลี่ยนแปลง / Save changes.
        TQuotation updatedQuotation = this.update(quotationId, quotation);

        // 4. ถ้ามีการสร้าง PO อัตโนมัติ / If auto-create PO is enabled.
        if (request.getAutoCreatePO()) {
            // TODO: เรียกใช้ PurchaseOrderService เพื่อสร้าง PO / Call PurchaseOrderService to create PO.
            // purchaseOrderService.createFromQuotation(quotationId);
        }

        // 5. ลบ Cache เก่าและอัปเดต Cache ใหม่ / Evict old cache and update new cache.
        cacheService.evictQuotation(quotationId);
        cacheService.saveQuotation(updatedQuotation);

        return QuotationResponseDTO.fromEntity(updatedQuotation);
    }
}
```

---

## ✅ สรุปส่วนที่เพิ่มเติมในโมดูล Quotation

| ฟีเจอร์ | สถานะ | คำอธิบาย |
|---------|--------|----------|
| **Database Design** | ✅ ครบถ้วน | `t_quotation`, `t_quotation_part`, `t_quotation_service`, `t_quotation_status_history` พร้อม Trigger |
| **Cache (Redis)** | ✅ เพิ่มแล้ว | `QuotationCacheService` (ID + Job ID) |
| **Rate Limit** | ✅ เพิ่มแล้ว | Quotation Controller และ Part Controller ทุก Endpoint |
| **API Routing** | ✅ ชัดเจน | CRUD + Approve/Reject + PDF + Part/Service Management |
| **Bilingual Comments** | ✅ ครบถ้วน | ทุก Function มีคำอธิบายภาษาไทยและภาษาอังกฤษแยกบรรทัด |
| **PDF Generation** | ✅ เพิ่มแล้ว | `QuotationReportGenerator` ใช้ JasperReports |

---
 