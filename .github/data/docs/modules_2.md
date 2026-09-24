**โมดูลที่ 2: 🚗 Job Card Management (การจัดการใบงานซ่อมรถ)**

โครงสร้างที่ครบถ้วนของโมดูล Job Card ตามข้อกำหนดของคุณ ได้แก่:
- ออกแบบฐานข้อมูลที่เกี่ยวข้องทั้งหมด (ครบทุกตาราง)
- ระบบ Cache (Redis) สำหรับข้อมูลใบงานและสถานะ
- ระบบ Rate Limit สำหรับ API แต่ละ Endpoint
- REST API Routing พร้อมคำอธิบายฟังก์ชันภาษาไทย/อังกฤษแยกบรรทัด

---

## 📁 โครงสร้างโมดูล Job Card (`modules/job`)

```
modules/job/
├── application/
│   ├── interfaces/
│   │   ├── JobService.java
│   │   ├── JobStatusService.java
│   │   └── JobPartSaleService.java
│   ├── impl/
│   │   ├── JobServiceImpl.java
│   │   ├── JobStatusServiceImpl.java
│   │   └── JobPartSaleServiceImpl.java
│   └── usecase/
│       ├── CreateJobUseCase.java
│       ├── UpdateJobUseCase.java
│       ├── GetJobUseCase.java
│       ├── DeleteJobUseCase.java
│       ├── ChangeJobStatusUseCase.java
│       ├── AddJobServiceUseCase.java
│       ├── AddJobPartUseCase.java
│       └── ListJobsByStatusUseCase.java
├── domain/
│   ├── TJob.java
│   ├── TJobService.java
│   ├── TJobPartSales.java
│   ├── TJobServiceCarSymptom.java
│   ├── TJobDiagTroubleCode.java
│   ├── TJobStatusHistory.java
│   ├── enums/
│   │   └── JobStatus.java         // OPEN, IN_PROGRESS, QUOTATION_PENDING, QUOTATION_APPROVED, PART_PICKING, REPAIR_IN_PROGRESS, REPAIR_DONE, INVOICE_PENDING, INVOICE_CREATED, PAYMENT_RECEIVED, CLOSED, CANCELLED, ON_HOLD, WAITING_PARTS
│   └── valueobjects/
│       ├── JobNumber.java
│       └── Mileage.java
├── infrastructure/
│   ├── repository/
│   │   ├── JobRepository.java
│   │   ├── JobServiceRepository.java
│   │   ├── JobPartSalesRepository.java
│   │   ├── JobSymptomRepository.java
│   │   ├── JobDiagCodeRepository.java
│   │   ├── JobStatusHistoryRepository.java
│   │   └── impl/
│   │       ├── JobRepositoryImpl.java
│   │       ├── JobServiceRepositoryImpl.java
│   │       └── JobPartSalesRepositoryImpl.java
│   ├── cache/                              // ⬅️ ระบบ Cache สำหรับ Job
│   │   ├── JobCacheService.java
│   │   └── JobStatusCacheService.java
│   ├── entity/
│   │   ├── JobEntity.java
│   │   ├── JobServiceEntity.java
│   │   ├── JobPartSalesEntity.java
│   │   ├── JobSymptomEntity.java
│   │   ├── JobDiagCodeEntity.java
│   │   └── JobStatusHistoryEntity.java
│   └── mapper/
│       ├── JobMapper.java
│       ├── JobServiceMapper.java
│       └── JobPartSalesMapper.java
└── presentation/
    ├── controller/
    │   ├── JobController.java              // CRUD + Status + Report
    │   ├── JobServiceController.java       // Manage service items
    │   └── JobPartController.java          // Manage part sales
    ├── dto/
    │   ├── request/
    │   │   ├── JobCreateRequestDTO.java
    │   │   ├── JobUpdateRequestDTO.java
    │   │   ├── JobStatusChangeRequestDTO.java
    │   │   ├── JobServiceRequestDTO.java
    │   │   └── JobPartRequestDTO.java
    │   └── response/
    │       ├── JobResponseDTO.java
    │       ├── JobDetailResponseDTO.java
    │       └── JobStatusHistoryDTO.java
    └── validator/
        └── JobValidator.java
```

---

## 🗄️ ออกแบบฐานข้อมูล (Database Design) สำหรับโมดูล Job Card

### 📄 SQL DDL (ไฟล์: `src/main/resources/db/migration/V2__job_card_schema.sql`)

```sql
-- ==============================================
-- ตาราง: t_job (ใบงานหลัก)
-- Main table for repair work orders.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_job (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_no VARCHAR(20) UNIQUE NOT NULL,                 -- เลขที่ใบงาน (เช่น JOB-2026-0001)
    customer_id UUID NOT NULL REFERENCES m_customer(id) ON DELETE RESTRICT,
    car_id UUID NOT NULL REFERENCES m_car(id) ON DELETE RESTRICT,
    mechanic_id UUID NOT NULL REFERENCES m_staff(id) ON DELETE RESTRICT,
    status VARCHAR(30) NOT NULL DEFAULT 'OPEN',         -- OPEN, IN_PROGRESS, QUOTATION_PENDING, etc.
    start_date TIMESTAMP NOT NULL DEFAULT NOW(),
    end_date TIMESTAMP,
    symptom TEXT,                                       -- อาการเบื้องต้น / Initial symptom.
    diagnosis_note TEXT,                                -- หมายเหตุการวินิจฉัย / Diagnosis note.
    mileage INTEGER,                                    -- ระยะทาง (กม.) / Mileage (km).
    estimated_cost DECIMAL(15,2),                       -- ค่าใช้จ่ายโดยประมาณ / Estimated cost.
    actual_cost DECIMAL(15,2),                          -- ค่าใช้จ่ายจริง / Actual cost.
    priority VARCHAR(20) DEFAULT 'NORMAL',              -- NORMAL, URGENT, EMERGENCY
    -- Audit Fields (from GenericBusinessEntity)
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    deleted BOOLEAN DEFAULT FALSE,
    user_id UUID NOT NULL,                              -- ผู้สร้าง (พนักงานหน้าร้าน)
    whitelabel_id UUID NOT NULL                         -- บริษัท/สาขา (Multi-tenancy)
);

CREATE INDEX idx_t_job_customer ON t_job(customer_id);
CREATE INDEX idx_t_job_car ON t_job(car_id);
CREATE INDEX idx_t_job_mechanic ON t_job(mechanic_id);
CREATE INDEX idx_t_job_status ON t_job(status);
CREATE INDEX idx_t_job_whitelabel ON t_job(whitelabel_id);
CREATE INDEX idx_t_job_deleted ON t_job(deleted);

-- ==============================================
-- ตาราง: t_job_service (รายการบริการที่กำหนดในใบงาน)
-- Services/repairs specified in the job card.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_job_service (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES t_job(id) ON DELETE CASCADE,
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

CREATE INDEX idx_t_job_service_job ON t_job_service(job_id);
CREATE INDEX idx_t_job_service_whitelabel ON t_job_service(whitelabel_id);

-- ==============================================
-- ตาราง: t_job_part_sales (รายการอะไหล่ที่ขาย/ใช้ในใบงาน)
-- Parts used/sold in the job card.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_job_part_sales (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES t_job(id) ON DELETE CASCADE,
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

CREATE INDEX idx_t_job_part_sales_job ON t_job_part_sales(job_id);
CREATE INDEX idx_t_job_part_sales_whitelabel ON t_job_part_sales(whitelabel_id);

-- ==============================================
-- ตาราง: t_job_service_car_symptom (อาการของรถที่แจ้งเข้ามา)
-- Vehicle symptoms reported by customer.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_job_service_car_symptom (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES t_job(id) ON DELETE CASCADE,
    symptom_code VARCHAR(20),                           -- รหัสอาการ
    symptom_description TEXT NOT NULL,                  -- คำอธิบายอาการ
    severity VARCHAR(20) DEFAULT 'MEDIUM',              -- LOW, MEDIUM, HIGH, CRITICAL
    reported_by VARCHAR(100),                           -- ผู้แจ้ง (ลูกค้าหรือพนักงาน)
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_job_symptom_job ON t_job_service_car_symptom(job_id);

-- ==============================================
-- ตาราง: t_job_diag_trouble_code (รหัสข้อผิดพลาดจากการวินิจฉัย - OBD2/Diagnostic)
-- Diagnostic trouble codes (OBD2, etc.) found during inspection.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_job_diag_trouble_code (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES t_job(id) ON DELETE CASCADE,
    trouble_code VARCHAR(20) NOT NULL,                  -- เช่น P0300, P0420
    description TEXT,                                   -- คำอธิบายโค้ด
    system VARCHAR(50),                                 -- ระบบที่เกี่ยวข้อง (Engine, ABS, etc.)
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_job_diag_code_job ON t_job_diag_trouble_code(job_id);

-- ==============================================
-- ตาราง: t_job_status_history (ประวัติการเปลี่ยนสถานะใบงาน)
-- Track job status changes over time.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_job_status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES t_job(id) ON DELETE CASCADE,
    from_status VARCHAR(30),
    to_status VARCHAR(30) NOT NULL,
    changed_by UUID NOT NULL REFERENCES m_user(id) ON DELETE RESTRICT,
    changed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    reason TEXT,                                        -- เหตุผลในการเปลี่ยนสถานะ
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_job_status_history_job ON t_job_status_history(job_id);
CREATE INDEX idx_t_job_status_history_changed ON t_job_status_history(changed_at);

-- ==============================================
-- ฟังก์ชันสร้างเลขที่ใบงานอัตโนมัติ (Auto-generate Job Number)
-- Function to generate unique job number with year and sequence.
-- ==============================================
CREATE OR REPLACE FUNCTION generate_job_no()
RETURNS TRIGGER AS $$
DECLARE
    year_part TEXT;
    sequence_part TEXT;
BEGIN
    year_part := TO_CHAR(NOW(), 'YYYY');
    sequence_part := LPAD(CAST((
        SELECT COALESCE(MAX(CAST(SUBSTRING(job_no FROM 9) AS INTEGER)), 0) + 1
        FROM t_job
        WHERE job_no LIKE 'JOB-' || year_part || '-%'
    ) AS TEXT), 4, '0');
    NEW.job_no := 'JOB-' || year_part || '-' || sequence_part;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- เรียกใช้ฟังก์ชันนี้ก่อน INSERT ทุกครั้ง / Apply before INSERT.
CREATE TRIGGER trg_generate_job_no
BEFORE INSERT ON t_job
FOR EACH ROW
EXECUTE FUNCTION generate_job_no();
```

---

## 🧠 ระบบ Cache (Redis) สำหรับ Job Card

### `infrastructure/cache/JobCacheService.java`

```java
package com.template.app.modules.job.infrastructure.cache;

import com.template.app.modules.job.domain.TJob;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.CachePut;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.UUID;

@Service
public class JobCacheService {

    /*
        ฟังก์ชันนี้ดึงข้อมูลใบงานจาก Cache หากมีข้อมูลอยู่จะไม่เรียกฐานข้อมูล (ลดภาระ DB)
        This function retrieves job data from cache. If present, it avoids hitting the database (reduces DB load).
        ใช้ Redis Key: job:{id}
        Redis Key: job:{id}
    */
    @Cacheable(value = "jobs", key = "#jobId")
    public TJob getJob(UUID jobId) {
        // คืนค่า null เพื่อให้ Spring รู้ว่าต้องไปดึงจาก Repository
        // Return null so Spring knows to fetch from Repository (if not cached).
        return null;
    }

    /*
        ฟังก์ชันนี้จะอัปเดตข้อมูลใน Cache เมื่อมีการเรียกใช้ฟังก์ชัน saveJob
        This function updates the cache when saveJob is called.
    */
    @CachePut(value = "jobs", key = "#job.id")
    public TJob saveJob(TJob job) {
        // Spring จะอัปเดต Cache ด้วยค่าที่คืนจาก Method
        // Spring updates the cache with the returned value.
        return job;
    }

    /*
        ฟังก์ชันนี้จะลบข้อมูลใบงานออกจาก Cache เมื่อมีการลบหรือปิดใบงาน
        This function evicts the job from cache when the job is deleted or closed.
    */
    @CacheEvict(value = "jobs", key = "#jobId")
    public void evictJob(UUID jobId) {
        // ลบ Cache ทิ้ง / Evict the cache entry.
    }

    /*
        ฟังก์ชันนี้ลบ Cache ของใบงานทั้งหมด (ใช้เมื่อมีการปิดระบบครั้งใหญ่)
        This function clears all job caches (used during major system updates).
    */
    @CacheEvict(value = "jobs", allEntries = true)
    public void evictAllJobs() {
        // ลบทุก cache key ที่มี prefix jobs:* / Evict all keys under jobs:*.
    }
}
```

---

## ⏱️ ระบบ Rate Limit สำหรับ Job Controller

```java
package com.template.app.modules.job.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.job.application.interfaces.JobService;
import com.template.app.modules.job.presentation.dto.request.JobCreateRequestDTO;
import com.template.app.modules.job.presentation.dto.request.JobStatusChangeRequestDTO;
import com.template.app.modules.job.presentation.dto.response.JobResponseDTO;
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
@RequestMapping("/api/v1/jobs")
@Tag(name = "Job Card", description = "Job Card Management APIs")
@RequiredArgsConstructor
public class JobController {

    private final JobService jobService;

    // ========================================================================
    // 1. CREATE JOB
    // ========================================================================

    /*
        API: POST /api/v1/jobs
        ฟังก์ชันนี้สร้างใบงานใหม่ โดยรับข้อมูลลูกค้า รถยนต์ และช่างที่รับผิดชอบ
        This function creates a new job card by receiving customer, vehicle, and assigned mechanic data.
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที (ป้องกันการสร้างใบงานซ้ำ/สแปม)
        Rate Limit: Allows 30 requests per minute (prevent duplicate/spam job creation).
    */
    @PostMapping
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Create a new job card")
    public ResponseEntity<JobResponseDTO> createJob(@Valid @RequestBody JobCreateRequestDTO request) 
            throws SystemGlobalException {
        JobResponseDTO response = jobService.createJob(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 2. GET JOB BY ID
    // ========================================================================

    /*
        API: GET /api/v1/jobs/{id}
        ฟังก์ชันนี้ดึงข้อมูลใบงานตาม ID (ใช้ Cache ช่วยลดภาระฐานข้อมูล)
        This function retrieves a job by its ID (uses caching to reduce database load).
        Rate Limit: อนุญาต 100 ครั้งต่อ 1 นาที (การเรียกดูข้อมูลบ่อยๆ)
        Rate Limit: Allows 100 requests per minute (frequent data retrieval).
    */
    @GetMapping("/{id}")
    @RateLimit(limit = 100, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get job by ID")
    public ResponseEntity<JobResponseDTO> getJob(@PathVariable UUID id) throws SystemGlobalException {
        JobResponseDTO response = jobService.getJob(id);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. LIST JOBS WITH PAGINATION
    // ========================================================================

    /*
        API: GET /api/v1/jobs
        ฟังก์ชันนี้แสดงรายการใบงานแบบแบ่งหน้า พร้อมตัวกรองตามสถานะและช่วงวันที่
        This function lists jobs with pagination, filtering by status and date range.
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @GetMapping
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "List jobs with pagination")
    public ResponseEntity<Page<JobResponseDTO>> listJobs(
            @RequestParam(required = false) String status,
            @RequestParam(required = false) String startDate,
            @RequestParam(required = false) String endDate,
            Pageable pageable) throws SystemGlobalException {
        Page<JobResponseDTO> page = jobService.listJobs(status, startDate, endDate, pageable);
        return ResponseEntity.ok(page);
    }

    // ========================================================================
    // 4. UPDATE JOB
    // ========================================================================

    /*
        API: PUT /api/v1/jobs/{id}
        ฟังก์ชันนี้แก้ไขข้อมูลใบงาน (ยกเว้นสถานะ) เช่น อาการ, หมายเหตุ, ระยะทาง
        This function updates job details (except status) such as symptoms, notes, mileage.
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @PutMapping("/{id}")
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Update job details")
    public ResponseEntity<JobResponseDTO> updateJob(@PathVariable UUID id, 
                                                    @Valid @RequestBody JobUpdateRequestDTO request) 
            throws SystemGlobalException {
        JobResponseDTO response = jobService.updateJob(id, request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 5. CHANGE JOB STATUS
    // ========================================================================

    /*
        API: PUT /api/v1/jobs/{id}/status
        ฟังก์ชันนี้เปลี่ยนสถานะใบงาน (เช่น OPEN -> IN_PROGRESS) และบันทึกประวัติการเปลี่ยนแปลง
        This function changes the job status (e.g., OPEN -> IN_PROGRESS) and records the change history.
        Rate Limit: อนุญาต 60 ครั้งต่อ 1 นาที (การอัปเดตสถานะบ่อยครั้งระหว่างทำงาน)
        Rate Limit: Allows 60 requests per minute (frequent status updates during work).
    */
    @PutMapping("/{id}/status")
    @RateLimit(limit = 60, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Change job status")
    public ResponseEntity<JobResponseDTO> changeStatus(@PathVariable UUID id,
                                                       @Valid @RequestBody JobStatusChangeRequestDTO request) 
            throws SystemGlobalException {
        JobResponseDTO response = jobService.changeStatus(id, request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 6. DELETE JOB (SOFT DELETE)
    // ========================================================================

    /*
        API: DELETE /api/v1/jobs/{id}
        ฟังก์ชันนี้ลบใบงานแบบ Soft Delete (เปลี่ยนสถานะเป็น DELETED และซ่อนจากรายการ)
        This function soft-deletes the job (marks as DELETED and hides from lists).
        Rate Limit: อนุญาต 10 ครั้งต่อ 1 ชั่วโมง (ป้องกันการลบข้อมูลโดยไม่ตั้งใจ)
        Rate Limit: Allows 10 requests per 1 hour (prevent accidental deletions).
    */
    @DeleteMapping("/{id}")
    @RateLimit(limit = 10, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Delete job (soft delete)")
    public ResponseEntity<Void> deleteJob(@PathVariable UUID id) throws SystemGlobalException {
        jobService.deleteJob(id);
        return ResponseEntity.noContent().build();
    }

    // ========================================================================
    // 7. GENERATE JOB REPORT (PDF)
    // ========================================================================

    /*
        API: GET /api/v1/jobs/report/{id}
        ฟังก์ชันนี้สร้าง PDF รายละเอียดใบงาน (ใช้ JasperReports) สำหรับพิมพ์/ส่งให้ลูกค้า
        This function generates a PDF job report (using JasperReports) for printing or customer delivery.
        Rate Limit: อนุญาต 20 ครั้งต่อ 5 นาที (การสร้าง PDF ใช้ทรัพยากรมาก)
        Rate Limit: Allows 20 requests per 5 minutes (PDF generation is resource-intensive).
    */
    @GetMapping("/report/{id}")
    @RateLimit(limit = 20, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Generate job report PDF")
    public ResponseEntity<byte[]> generateJobReport(@PathVariable UUID id) throws SystemGlobalException {
        byte[] pdf = jobService.generateJobReport(id);
        return ResponseEntity.ok()
                .header("Content-Type", "application/pdf")
                .header("Content-Disposition", "inline; filename=job_" + id + ".pdf")
                .body(pdf);
    }

    // ========================================================================
    // 8. GET JOB STATUS HISTORY
    // ========================================================================

    /*
        API: GET /api/v1/jobs/{id}/history
        ฟังก์ชันนี้ดึงประวัติการเปลี่ยนสถานะของใบงาน
        This function retrieves the status change history of a job.
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @GetMapping("/{id}/history")
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get job status history")
    public ResponseEntity<List<JobStatusHistoryDTO>> getStatusHistory(@PathVariable UUID id) 
            throws SystemGlobalException {
        List<JobStatusHistoryDTO> history = jobService.getStatusHistory(id);
        return ResponseEntity.ok(history);
    }

    // ========================================================================
    // 9. ADD SERVICE ITEM TO JOB
    // ========================================================================

    /*
        API: POST /api/v1/jobs/{id}/services
        ฟังก์ชันนี้เพิ่มรายการบริการให้กับใบงาน (เช่น เปลี่ยนถ่ายน้ำมัน, เช็คระบบเบรก)
        This function adds a service item to the job (e.g., oil change, brake system check).
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที
        Rate Limit: Allows 30 requests per minute.
    */
    @PostMapping("/{id}/services")
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Add service to job")
    public ResponseEntity<JobServiceResponseDTO> addService(@PathVariable UUID id,
                                                            @Valid @RequestBody JobServiceRequestDTO request) 
            throws SystemGlobalException {
        JobServiceResponseDTO response = jobService.addService(id, request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 10. ADD PART TO JOB
    // ========================================================================

    /*
        API: POST /api/v1/jobs/{id}/parts
        ฟังก์ชันนี้เพิ่มอะไหล่ที่ใช้ในงานซ่อมให้กับใบงาน (เชื่อมโยงไปยัง Inventory)
        This function adds parts used in the repair to the job (links to Inventory).
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที
        Rate Limit: Allows 30 requests per minute.
    */
    @PostMapping("/{id}/parts")
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Add part to job")
    public ResponseEntity<JobPartResponseDTO> addPart(@PathVariable UUID id,
                                                      @Valid @RequestBody JobPartRequestDTO request) 
            throws SystemGlobalException {
        JobPartResponseDTO response = jobService.addPart(id, request);
        return ResponseEntity.ok(response);
    }
}
```

---

## 🧩 ตัวอย่าง Domain Entity และ Enum

### `domain/enums/JobStatus.java`

```java
package com.template.app.modules.job.domain.enums;

/*
    Enum สำหรับสถานะของใบงานซ่อมรถ
    Enum for job card statuses.
*/
public enum JobStatus {
    OPEN,               // เปิดใบงาน / Job opened.
    IN_PROGRESS,        // กำลังดำเนินการ / In progress.
    QUOTATION_PENDING,  // รอการอนุมัติใบเสนอราคา / Waiting for quotation approval.
    QUOTATION_APPROVED, // อนุมัติใบเสนอราคาแล้ว / Quotation approved.
    PART_PICKING,       // กำลังเบิกอะไหล่ / Part picking in progress.
    REPAIR_IN_PROGRESS, // กำลังซ่อม / Repair in progress.
    REPAIR_DONE,        // ซ่อมเสร็จ / Repair completed.
    INVOICE_PENDING,    // รอออกใบแจ้งหนี้ / Invoice pending.
    INVOICE_CREATED,    // ออกใบแจ้งหนี้แล้ว / Invoice created.
    PAYMENT_RECEIVED,   // รับชำระเงินแล้ว / Payment received.
    CLOSED,             // ปิดงาน / Job closed.
    CANCELLED,          // ยกเลิก / Cancelled.
    ON_HOLD,            // ระงับชั่วคราว / On hold.
    WAITING_PARTS       // รออะไหล่ / Waiting for parts.
}
```

### `domain/TJob.java`

```java
package com.template.app.modules.job.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import com.template.app.modules.job.domain.enums.JobStatus;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.UUID;

@Data
@EqualsAndHashCode(callSuper = true)
public class TJob extends GenericBusinessClass {

    private String jobNo;               // เลขที่ใบงาน / Job number.
    private UUID customerId;            // ID ลูกค้า / Customer ID.
    private UUID carId;                 // ID รถยนต์ / Vehicle ID.
    private UUID mechanicId;            // ID ช่างผู้รับผิดชอบ / Assigned mechanic ID.
    private JobStatus status;           // สถานะ / Status.
    private LocalDateTime startDate;    // วันที่เริ่ม / Start date.
    private LocalDateTime endDate;      // วันที่เสร็จ / Completion date.
    private String symptom;             // อาการ / Symptom.
    private String diagnosisNote;       // หมายเหตุวินิจฉัย / Diagnosis note.
    private Integer mileage;            // ระยะทาง (กม.) / Mileage (km).
    private BigDecimal estimatedCost;   // ค่าใช้จ่ายโดยประมาณ / Estimated cost.
    private BigDecimal actualCost;      // ค่าใช้จ่ายจริง / Actual cost.
    private String priority;            // ความสำคัญ: NORMAL, URGENT, EMERGENCY / Priority level.

    /*
        ฟังก์ชันนี้ใช้เปลี่ยนสถานะใบงาน พร้อมตรวจสอบความถูกต้องของลำดับสถานะ
        This function changes the job status, validating the transition order.
    */
    public void changeStatus(JobStatus newStatus) {
        // ตรวจสอบว่าไม่สามารถเปลี่ยนจาก CLOSED หรือ CANCELLED ได้
        // Cannot transition from CLOSED or CANCELLED.
        if (this.status == JobStatus.CLOSED || this.status == JobStatus.CANCELLED) {
            throw new IllegalStateException("Cannot change status of closed or cancelled job.");
        }
        this.status = newStatus;
        // ถ้าสถานะเป็น CLOSED หรือ REPAIR_DONE ให้บันทึกเวลา endDate / Set endDate for terminal statuses.
        if (newStatus == JobStatus.CLOSED || newStatus == JobStatus.REPAIR_DONE) {
            this.endDate = LocalDateTime.now();
        }
    }

    /*
        ฟังก์ชันนี้ใช้ตรวจสอบว่าใบงานสามารถสร้าง Quotation ได้หรือไม่
        This function checks if a quotation can be created for this job.
    */
    public boolean canCreateQuotation() {
        return this.status == JobStatus.OPEN || 
               this.status == JobStatus.IN_PROGRESS || 
               this.status == JobStatus.QUOTATION_PENDING;
    }

    /*
        ฟังก์ชันนี้ใช้คำนวณค่าใช้จ่ายรวมจากรายการบริการและอะไหล่
        This function calculates the total cost from service and part items.
    */
    public BigDecimal calculateTotalCost() {
        // TODO: รวมค่าบริการ + ค่าอะไหล่ / Aggregate service + part costs.
        return BigDecimal.ZERO;
    }
}
```

---

## 🔧 ฟังก์ชัน Service ตัวอย่าง (พร้อมคำอธิบาย)

### `application/impl/JobServiceImpl.java` (บางส่วน)

```java
package com.template.app.modules.job.application.impl;

import com.template.app.modules.job.application.interfaces.JobService;
import com.template.app.modules.job.domain.TJob;
import com.template.app.modules.job.domain.enums.JobStatus;
import com.template.app.modules.job.infrastructure.cache.JobCacheService;
import com.template.app.modules.job.infrastructure.repository.JobRepository;
import com.template.app.modules.job.presentation.dto.request.JobCreateRequestDTO;
import com.template.app.modules.job.presentation.dto.response.JobResponseDTO;
import com.template.app._shared.application.GenericServiceImpl;
import com.template.app.exception.SystemGlobalException;
import org.springframework.stereotype.Service;

@Service
public class JobServiceImpl extends GenericServiceImpl<TJob, JobRepository> implements JobService {

    private final JobCacheService jobCacheService;

    public JobServiceImpl(JobRepository repository, JobCacheService jobCacheService) {
        super(repository);
        this.jobCacheService = jobCacheService;
    }

    /*
        ฟังก์ชันนี้สร้างใบงานใหม่ โดยรับข้อมูล DTO และบันทึกลงฐานข้อมูล พร้อมอัปเดต Cache
        This function creates a new job by receiving DTO data, saves to DB, and updates the cache.
    */
    @Override
    public JobResponseDTO createJob(JobCreateRequestDTO request) throws SystemGlobalException {
        // 1. แปลง DTO เป็น Domain Entity / Convert DTO to Domain Entity.
        TJob newJob = new TJob();
        newJob.setCustomerId(request.getCustomerId());
        newJob.setCarId(request.getCarId());
        newJob.setMechanicId(request.getMechanicId());
        newJob.setStatus(JobStatus.OPEN);
        newJob.setSymptom(request.getSymptom());
        newJob.setMileage(request.getMileage());
        
        // 2. บันทึกผ่าน Generic Repository (จะสร้าง jobNo อัตโนมัติผ่าน Trigger) / Save via Generic Repository (jobNo auto-generated by Trigger).
        TJob savedJob = this.create(newJob);
        
        // 3. อัปเดต Cache ด้วย Job ที่สร้างแล้ว / Update cache with the created job.
        jobCacheService.saveJob(savedJob);
        
        // 4. แปลงเป็น DTO และส่งคืน / Convert to DTO and return.
        return JobResponseDTO.fromEntity(savedJob);
    }
}
```

---

## ✅ สรุปส่วนที่เพิ่มเติมในโมดูล Job Card

| ฟีเจอร์ | สถานะ | คำอธิบาย |
|---------|--------|----------|
| **Database Design** | ✅ ครบถ้วน | `t_job`, `t_job_service`, `t_job_part_sales`, `t_job_symptom`, `t_job_diag_code`, `t_job_status_history` พร้อม Trigger |
| **Cache (Redis)** | ✅ เพิ่มแล้ว | `JobCacheService` ใช้ `@Cacheable`, `@CachePut`, `@CacheEvict` |
| **Rate Limit** | ✅ เพิ่มแล้ว | ใช้ `@RateLimit` ใน Controller ทุก Endpoint พร้อมกำหนด Limit และ Duration ตามเหมาะสม |
| **API Routing** | ✅ ชัดเจน | CRUD + Status Change + Service/Part + Report + History |
| **Bilingual Comments** | ✅ ครบถ้วน | ทุก Function มีคำอธิบายภาษาไทยและภาษาอังกฤษแยกบรรทัด |

---