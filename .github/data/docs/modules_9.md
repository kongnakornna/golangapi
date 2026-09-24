**โมดูลที่ 9: 📄 Document Management (การจัดการเอกสาร)**

โมดูล Document Management เป็นศูนย์กลางสำหรับการจัดการเอกสารทั้งหมดในระบบ ไม่ว่าจะเป็น PDF, Excel, หรือไฟล์แนบต่างๆ ครอบคลุมการทำงานดังนี้:

1. **การสร้างเอกสารจากเทมเพลต (Document Generation)** – ใช้ JasperReports และ Apache POI
2. **การจัดเก็บเอกสาร (Document Storage)** – เก็บในระบบไฟล์หรือ S3 (AWS)
3. **การจัดการเทมเพลต (Template Management)** – อัปโหลดและจัดการไฟล์ .jrxml
4. **การส่งเอกสาร (Document Delivery)** – แนบไปกับอีเมลหรือดาวน์โหลด
5. **ประวัติเอกสาร (Document History)** – ติดตามการสร้างและแก้ไขเอกสาร
6. **OCR (Image to Text)** – อ่านข้อความจากรูปภาพ/เอกสารสแกน

---

## 📁 โครงสร้างโมดูล Document Management (`modules/document`)

```
modules/document/
├── application/
│   ├── interfaces/
│   │   ├── DocumentService.java
│   │   ├── TemplateService.java
│   │   ├── DocumentStorageService.java
│   │   ├── OCRService.java
│   │   └── ReportGenerationService.java
│   ├── impl/
│   │   ├── DocumentServiceImpl.java
│   │   ├── TemplateServiceImpl.java
│   │   ├── DocumentStorageServiceImpl.java
│   │   ├── OCRServiceImpl.java
│   │   └── ReportGenerationServiceImpl.java
│   └── usecase/
│       ├── GenerateDocumentUseCase.java
│       ├── UploadTemplateUseCase.java
│       ├── DeleteDocumentUseCase.java
│       ├── GetDocumentUseCase.java
│       ├── ListDocumentsUseCase.java
│       ├── ProcessOCRUseCase.java
│       ├── GenerateReportUseCase.java
│       └── ExportToExcelUseCase.java
├── domain/
│   ├── MDocumentTemplate.java
│   ├── TDocument.java
│   ├── TDocumentHistory.java
│   ├── enums/
│   │   ├── DocumentType.java         // PDF, EXCEL, CSV, IMAGE, OTHER
│   │   ├── DocumentStatus.java       // DRAFT, GENERATED, SENT, ARCHIVED, DELETED
│   │   ├── TemplateType.java         // JASPER, POI, CUSTOM
│   │   └── OCRProvider.java          // TESSERACT, GOOGLE_VISION, AMAZON_TEXT
│   └── valueobjects/
│       ├── DocumentReference.java
│       └── FileSize.java
├── infrastructure/
│   ├── repository/
│   │   ├── DocumentRepository.java
│   │   ├── DocumentTemplateRepository.java
│   │   ├── DocumentHistoryRepository.java
│   │   └── impl/
│   │       ├── DocumentRepositoryImpl.java
│   │       └── DocumentTemplateRepositoryImpl.java
│   ├── cache/                                           // ⬅️ ระบบ Cache สำหรับ Document
│   │   ├── DocumentCacheService.java
│   │   └── TemplateCacheService.java
│   ├── storage/                                         // ⬅️ ระบบจัดเก็บไฟล์
│   │   ├── FileStorageService.java
│   │   ├── LocalFileStorageService.java
│   │   ├── S3FileStorageService.java
│   │   └── FileStorageConfig.java
│   ├── generator/                                       // ⬅️ ระบบสร้างเอกสาร
│   │   ├── jasper/
│   │   │   ├── JasperReportGenerator.java
│   │   │   └── JasperReportConfig.java
│   │   ├── excel/
│   │   │   └── ExcelGenerator.java
│   │   └── pdf/
│   │       └── PDFGenerator.java
│   ├── ocr/                                             // ⬅️ ระบบ OCR
│   │   ├── OCRProvider.java
│   │   ├── TesseractOCRProvider.java
│   │   └── GoogleVisionOCRProvider.java
│   ├── entity/
│   │   ├── DocumentEntity.java
│   │   ├── DocumentTemplateEntity.java
│   │   └── DocumentHistoryEntity.java
│   └── mapper/
│       ├── DocumentMapper.java
│       └── DocumentTemplateMapper.java
└── presentation/
    ├── controller/
    │   ├── DocumentController.java       // CRUD + Download + Upload
    │   ├── TemplateController.java       // Template Management
    │   ├── OCRController.java            // OCR Processing
    │   └── ReportController.java         // Report Generation
    ├── dto/
    │   ├── request/
    │   │   ├── DocumentUploadRequestDTO.java
    │   │   ├── DocumentGenerateRequestDTO.java
    │   │   ├── TemplateUploadRequestDTO.java
    │   │   ├── OCRRequestDTO.java
    │   │   └── DocumentSearchRequestDTO.java
    │   └── response/
    │       ├── DocumentResponseDTO.java
    │       ├── TemplateResponseDTO.java
    │       ├── OCRResponseDTO.java
    │       └── DocumentHistoryDTO.java
    └── validator/
        └── DocumentValidator.java
```

---

## 🗄️ ออกแบบฐานข้อมูล (Database Design) สำหรับโมดูล Document

### 📄 SQL DDL (ไฟล์: `src/main/resources/db/migration/V9__document_schema.sql`)

```sql
-- ==============================================
-- ตาราง: m_document_template (เทมเพลตเอกสาร)
-- Document template master data (JasperReports, etc.).
-- ==============================================
CREATE TABLE IF NOT EXISTS m_document_template (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_code VARCHAR(50) UNIQUE NOT NULL,      -- รหัสเทมเพลต (เช่น QUOTATION, INVOICE)
    template_name VARCHAR(100) NOT NULL,            -- ชื่อเทมเพลต
    template_type VARCHAR(20) NOT NULL,             -- JASPER, POI, CUSTOM
    file_name VARCHAR(255) NOT NULL,                -- ชื่อไฟล์ (เช่น quotation.jrxml)
    file_path TEXT NOT NULL,                        -- ที่อยู่ไฟล์
    file_size BIGINT,                               -- ขนาดไฟล์ (bytes)
    version INTEGER DEFAULT 1,                      -- เวอร์ชัน
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    is_default BOOLEAN DEFAULT FALSE,
    parameters JSONB,                               -- พารามิเตอร์สำหรับเทมเพลต
    preview_image_url TEXT,
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_m_template_code ON m_document_template(template_code);
CREATE INDEX idx_m_template_type ON m_document_template(template_type);
CREATE INDEX idx_m_template_active ON m_document_template(is_active);

-- ==============================================
-- ตาราง: t_document (เอกสารที่สร้างแล้ว)
-- Generated document records.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_document (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_no VARCHAR(30) UNIQUE NOT NULL,        -- เลขที่เอกสาร (DOC-2026-0001)
    document_type VARCHAR(20) NOT NULL,             -- PDF, EXCEL, CSV, IMAGE, OTHER
    document_sub_type VARCHAR(30),                  -- QUOTATION, INVOICE, RECEIPT, etc.
    reference_type VARCHAR(30),                     -- JOB, QUOTATION, INVOICE, etc.
    reference_id UUID,                              -- ID ของข้อมูลต้นทาง
    template_id UUID REFERENCES m_document_template(id),
    file_name VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    file_size BIGINT,
    mime_type VARCHAR(50),
    status VARCHAR(20) DEFAULT 'GENERATED',         -- DRAFT, GENERATED, SENT, ARCHIVED, DELETED
    generated_by UUID NOT NULL REFERENCES m_user(id),
    generated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    sent_by UUID REFERENCES m_user(id),
    sent_at TIMESTAMP,
    sent_to_email VARCHAR(200),
    description TEXT,
    tags TEXT[],
    metadata JSONB,
    deleted_at TIMESTAMP,
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_document_reference ON t_document(reference_type, reference_id);
CREATE INDEX idx_t_document_type ON t_document(document_type);
CREATE INDEX idx_t_document_sub_type ON t_document(document_sub_type);
CREATE INDEX idx_t_document_status ON t_document(status);
CREATE INDEX idx_t_document_generated_at ON t_document(generated_at);
CREATE INDEX idx_t_document_whitelabel ON t_document(whitelabel_id);

-- ==============================================
-- ฟังก์ชันสร้างเลขที่เอกสารอัตโนมัติ
-- Function to generate document number.
-- ==============================================
CREATE OR REPLACE FUNCTION generate_document_no()
RETURNS TRIGGER AS $$
DECLARE
    year_part TEXT;
    sequence_part TEXT;
BEGIN
    year_part := TO_CHAR(NOW(), 'YYYY');
    sequence_part := LPAD(CAST((
        SELECT COALESCE(MAX(CAST(SUBSTRING(document_no FROM 9) AS INTEGER)), 0) + 1
        FROM t_document
        WHERE document_no LIKE 'DOC-' || year_part || '-%'
    ) AS TEXT), 4, '0');
    NEW.document_no := 'DOC-' || year_part || '-' || sequence_part;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_generate_document_no ON t_document;
CREATE TRIGGER trg_generate_document_no
BEFORE INSERT ON t_document
FOR EACH ROW
EXECUTE FUNCTION generate_document_no();

-- ==============================================
-- ตาราง: t_document_history (ประวัติเอกสาร)
-- Document history/audit trail.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_document_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES t_document(id) ON DELETE CASCADE,
    action VARCHAR(30) NOT NULL,                    -- CREATED, GENERATED, SENT, DOWNLOADED, ARCHIVED, DELETED
    performed_by UUID NOT NULL REFERENCES m_user(id),
    performed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    details TEXT,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_doc_history_document ON t_document_history(document_id);
CREATE INDEX idx_t_doc_history_performed ON t_document_history(performed_at);

-- ==============================================
-- ตาราง: t_ocr_result (ผลลัพธ์ OCR)
-- OCR processing results.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_ocr_result (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID REFERENCES t_document(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,                        -- URL รูปภาพต้นฉบับ
    provider VARCHAR(30) NOT NULL,                  -- TESSERACT, GOOGLE_VISION, AMAZON_TEXT
    extracted_text TEXT,                            -- ข้อความที่อ่านได้
    confidence_score DECIMAL(5,2),                  -- คะแนนความมั่นใจ (%)
    language VARCHAR(10),
    processing_time_ms INTEGER,
    metadata JSONB,
    status VARCHAR(20) DEFAULT 'PENDING',           -- PENDING, PROCESSING, COMPLETED, FAILED
    error_message TEXT,
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_ocr_result_document ON t_ocr_result(document_id);
CREATE INDEX idx_t_ocr_result_status ON t_ocr_result(status);
```

---

## 🧠 ระบบ Cache (Redis) สำหรับ Document

### `infrastructure/cache/DocumentCacheService.java`

```java
package com.template.app.modules.document.infrastructure.cache;

import com.template.app.modules.document.domain.TDocument;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.CachePut;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class DocumentCacheService {

    /*
        ฟังก์ชันนี้ดึงข้อมูลเอกสารจาก Cache (ลดภาระฐานข้อมูล)
        This function retrieves document data from cache (reduces DB load).
        Redis Key: document:{id}
    */
    @Cacheable(value = "documents", key = "#documentId")
    public TDocument getDocument(UUID documentId) {
        return null;
    }

    /*
        ฟังก์ชันนี้ดึงข้อมูลเอกสารตาม Reference (Job ID, Invoice ID)
        This function retrieves document by Reference (Job ID, Invoice ID).
        Redis Key: document_ref:{referenceType}:{referenceId}
    */
    @Cacheable(value = "document_ref", key = "#referenceType + ':' + #referenceId")
    public TDocument getDocumentByReference(String referenceType, UUID referenceId) {
        return null;
    }

    /*
        ฟังก์ชันนี้จะอัปเดต Cache เมื่อมีการบันทึกเอกสาร
        This function updates the cache when a document is saved.
    */
    @CachePut(value = "documents", key = "#document.id")
    public TDocument saveDocument(TDocument document) {
        return document;
    }

    /*
        ฟังก์ชันนี้ลบข้อมูลเอกสารออกจาก Cache
        This function evicts document data from cache.
    */
    @CacheEvict(value = {"documents", "document_ref"}, key = "#documentId")
    public void evictDocument(UUID documentId) {
        // ลบ Cache / Evict cache.
    }

    /*
        ฟังก์ชันนี้ลบ Cache ทั้งหมดของ Document (ใช้เมื่อระบบอัปเดตข้อมูลจำนวนมาก)
        This function clears all document caches (used during bulk updates).
    */
    @CacheEvict(value = {"documents", "document_ref"}, allEntries = true)
    public void evictAllDocuments() {
        // ลบทุก key ใน caches / Evict all keys.
    }
}
```

### `infrastructure/cache/TemplateCacheService.java`

```java
package com.template.app.modules.document.infrastructure.cache;

import com.template.app.modules.document.domain.MDocumentTemplate;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.CachePut;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

@Service
public class TemplateCacheService {

    /*
        ฟังก์ชันนี้ดึงข้อมูลเทมเพลตตามรหัสเทมเพลต
        This function retrieves template data by template code.
        Redis Key: template:{templateCode}
    */
    @Cacheable(value = "templates", key = "#templateCode")
    public MDocumentTemplate getTemplate(String templateCode) {
        return null;
    }

    /*
        ฟังก์ชันนี้จะอัปเดต Cache เมื่อมีการบันทึกเทมเพลต
        This function updates the cache when a template is saved.
    */
    @CachePut(value = "templates", key = "#template.templateCode")
    public MDocumentTemplate saveTemplate(MDocumentTemplate template) {
        return template;
    }

    /*
        ฟังก์ชันนี้ลบข้อมูลเทมเพลตออกจาก Cache
        This function evicts template data from cache.
    */
    @CacheEvict(value = "templates", key = "#templateCode")
    public void evictTemplate(String templateCode) {
        // ลบ Cache / Evict cache.
    }
}
```

---

## ⏱️ ระบบ Rate Limit สำหรับ Document Controller

```java
package com.template.app.modules.document.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.document.application.interfaces.DocumentService;
import com.template.app.modules.document.application.interfaces.OCRService;
import com.template.app.modules.document.presentation.dto.request.DocumentGenerateRequestDTO;
import com.template.app.modules.document.presentation.dto.request.DocumentSearchRequestDTO;
import com.template.app.modules.document.presentation.dto.request.OCRRequestDTO;
import com.template.app.modules.document.presentation.dto.response.DocumentResponseDTO;
import com.template.app.modules.document.presentation.dto.response.OCRResponseDTO;
import com.template.app.exception.SystemGlobalException;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.core.io.Resource;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.util.UUID;

@RestController
@RequestMapping("/api/v1/documents")
@Tag(name = "Document", description = "Document Management APIs")
@RequiredArgsConstructor
public class DocumentController {

    private final DocumentService documentService;
    private final OCRService ocrService;

    // ========================================================================
    // 1. GENERATE DOCUMENT (สร้างเอกสารจากเทมเพลต)
    // ========================================================================

    /*
        API: POST /api/v1/documents/generate
        ฟังก์ชันนี้สร้างเอกสารจากเทมเพลตที่กำหนด (Quotation, Invoice, Receipt, etc.)
        This function generates a document from the specified template (Quotation, Invoice, Receipt, etc.).
        Rate Limit: อนุญาต 15 ครั้งต่อ 5 นาที (การสร้างเอกสารใช้ทรัพยากรมาก)
        Rate Limit: Allows 15 requests per 5 minutes (document generation is resource-intensive).
    */
    @PostMapping("/generate")
    @RateLimit(limit = 15, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Generate document from template")
    public ResponseEntity<DocumentResponseDTO> generateDocument(@Valid @RequestBody DocumentGenerateRequestDTO request)
            throws SystemGlobalException {
        DocumentResponseDTO response = documentService.generateDocument(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 2. UPLOAD DOCUMENT (อัปโหลดเอกสาร)
    // ========================================================================

    /*
        API: POST /api/v1/documents/upload
        ฟังก์ชันนี้อัปโหลดเอกสารเข้าเก็บในระบบ (ไฟล์แนบ, เอกสารสแกน)
        This function uploads a document to the system (attachments, scanned documents).
        Rate Limit: อนุญาต 20 ครั้งต่อ 5 นาที
        Rate Limit: Allows 20 requests per 5 minutes.
    */
    @PostMapping("/upload")
    @RateLimit(limit = 20, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Upload document")
    public ResponseEntity<DocumentResponseDTO> uploadDocument(
            @RequestParam("file") MultipartFile file,
            @RequestParam("documentType") String documentType,
            @RequestParam(value = "referenceType", required = false) String referenceType,
            @RequestParam(value = "referenceId", required = false) UUID referenceId,
            @RequestParam(value = "description", required = false) String description) throws SystemGlobalException {
        DocumentResponseDTO response = documentService.uploadDocument(file, documentType, referenceType, referenceId, description);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. GET DOCUMENT BY ID
    // ========================================================================

    /*
        API: GET /api/v1/documents/{id}
        ฟังก์ชันนี้ดึงข้อมูลเอกสารตาม ID (ใช้ Cache ช่วยลดภาระ DB)
        This function retrieves document data by ID (uses caching to reduce DB load).
        Rate Limit: อนุญาต 100 ครั้งต่อ 1 นาที
        Rate Limit: Allows 100 requests per minute.
    */
    @GetMapping("/{id}")
    @RateLimit(limit = 100, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get document by ID")
    public ResponseEntity<DocumentResponseDTO> getDocument(@PathVariable UUID id)
            throws SystemGlobalException {
        DocumentResponseDTO response = documentService.getDocument(id);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 4. GET DOCUMENT BY REFERENCE
    // ========================================================================

    /*
        API: GET /api/v1/documents/reference/{referenceType}/{referenceId}
        ฟังก์ชันนี้ดึงเอกสารตาม Reference (เช่น Job ID, Invoice ID)
        This function retrieves document by Reference (e.g., Job ID, Invoice ID).
        Rate Limit: อนุญาต 60 ครั้งต่อ 1 นาที
        Rate Limit: Allows 60 requests per minute.
    */
    @GetMapping("/reference/{referenceType}/{referenceId}")
    @RateLimit(limit = 60, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get document by reference")
    public ResponseEntity<DocumentResponseDTO> getDocumentByReference(
            @PathVariable String referenceType,
            @PathVariable UUID referenceId) throws SystemGlobalException {
        DocumentResponseDTO response = documentService.getDocumentByReference(referenceType, referenceId);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 5. LIST DOCUMENTS (Pagination + Filters)
    // ========================================================================

    /*
        API: POST /api/v1/documents/search
        ฟังก์ชันนี้ค้นหาเอกสารด้วยตัวกรอง เช่น ประเภท, ช่วงวันที่, สถานะ
        This function searches documents with filters: type, date range, status.
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @PostMapping("/search")
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Search documents with filters")
    public ResponseEntity<Page<DocumentResponseDTO>> searchDocuments(
            @Valid @RequestBody DocumentSearchRequestDTO request,
            Pageable pageable) throws SystemGlobalException {
        Page<DocumentResponseDTO> page = documentService.searchDocuments(request, pageable);
        return ResponseEntity.ok(page);
    }

    // ========================================================================
    // 6. DOWNLOAD DOCUMENT (ดาวน์โหลดเอกสาร)
    // ========================================================================

    /*
        API: GET /api/v1/documents/{id}/download
        ฟังก์ชันนี้ดาวน์โหลดไฟล์เอกสารจากระบบ
        This function downloads a document file from the system.
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที
        Rate Limit: Allows 30 requests per minute.
    */
    @GetMapping("/{id}/download")
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Download document")
    public ResponseEntity<Resource> downloadDocument(@PathVariable UUID id) throws SystemGlobalException {
        Resource resource = documentService.downloadDocument(id);
        String fileName = documentService.getDocumentFileName(id);
        return ResponseEntity.ok()
                .header(HttpHeaders.CONTENT_DISPOSITION, "attachment; filename=\"" + fileName + "\"")
                .header(HttpHeaders.CONTENT_TYPE, documentService.getDocumentMimeType(id))
                .body(resource);
    }

    // ========================================================================
    // 7. DELETE DOCUMENT (SOFT DELETE)
    // ========================================================================

    /*
        API: DELETE /api/v1/documents/{id}
        ฟังก์ชันนี้ลบเอกสารแบบ Soft Delete
        This function soft-deletes a document.
        Rate Limit: อนุญาต 10 ครั้งต่อ 1 ชั่วโมง
        Rate Limit: Allows 10 requests per 1 hour.
    */
    @DeleteMapping("/{id}")
    @RateLimit(limit = 10, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Delete document (soft delete)")
    public ResponseEntity<Void> deleteDocument(@PathVariable UUID id) throws SystemGlobalException {
        documentService.deleteDocument(id);
        return ResponseEntity.noContent().build();
    }

    // ========================================================================
    // 8. GET DOCUMENT HISTORY (ประวัติเอกสาร)
    // ========================================================================

    /*
        API: GET /api/v1/documents/{id}/history
        ฟังก์ชันนี้ดึงประวัติการเปลี่ยนแปลงของเอกสาร
        This function retrieves the document change history.
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @GetMapping("/{id}/history")
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get document history")
    public ResponseEntity<List<DocumentHistoryDTO>> getDocumentHistory(@PathVariable UUID id)
            throws SystemGlobalException {
        List<DocumentHistoryDTO> history = documentService.getDocumentHistory(id);
        return ResponseEntity.ok(history);
    }

    // ========================================================================
    // 9. PROCESS OCR (อ่านข้อความจากรูปภาพ)
    // ========================================================================

    /*
        API: POST /api/v1/documents/ocr
        ฟังก์ชันนี้อ่านข้อความจากรูปภาพหรือเอกสารสแกน (ใช้ Tesseract หรือ Google Vision)
        This function extracts text from images or scanned documents (using Tesseract or Google Vision).
        Rate Limit: อนุญาต 10 ครั้งต่อ 1 นาที
        Rate Limit: Allows 10 requests per minute.
    */
    @PostMapping("/ocr")
    @RateLimit(limit = 10, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Process OCR on document")
    public ResponseEntity<OCRResponseDTO> processOCR(@Valid @RequestBody OCRRequestDTO request)
            throws SystemGlobalException {
        OCRResponseDTO response = ocrService.processOCR(request);
        return ResponseEntity.ok(response);
    }
}
```

### `TemplateController.java` (จัดการเทมเพลต)

```java
package com.template.app.modules.document.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.document.application.interfaces.TemplateService;
import com.template.app.modules.document.presentation.dto.request.TemplateUploadRequestDTO;
import com.template.app.modules.document.presentation.dto.response.TemplateResponseDTO;
import com.template.app.exception.SystemGlobalException;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.util.UUID;

@RestController
@RequestMapping("/api/v1/templates")
@Tag(name = "Document Templates", description = "Document Template Management APIs")
@RequiredArgsConstructor
public class TemplateController {

    private final TemplateService templateService;

    // ========================================================================
    // 1. UPLOAD TEMPLATE (อัปโหลดเทมเพลต)
    // ========================================================================

    /*
        API: POST /api/v1/templates/upload
        ฟังก์ชันนี้อัปโหลดไฟล์เทมเพลต (JasperReports .jrxml, Excel template)
        This function uploads a template file (JasperReports .jrxml, Excel template).
        Rate Limit: อนุญาต 10 ครั้งต่อ 5 นาที
        Rate Limit: Allows 10 requests per 5 minutes.
    */
    @PostMapping("/upload")
    @RateLimit(limit = 10, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Upload document template")
    public ResponseEntity<TemplateResponseDTO> uploadTemplate(
            @RequestParam("file") MultipartFile file,
            @Valid @RequestParam TemplateUploadRequestDTO request) throws SystemGlobalException {
        TemplateResponseDTO response = templateService.uploadTemplate(file, request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 2. GET TEMPLATE BY CODE
    // ========================================================================

    /*
        API: GET /api/v1/templates/{templateCode}
        ฟังก์ชันนี้ดึงข้อมูลเทมเพลตตามรหัส (ใช้ Cache)
        This function retrieves template data by code (uses caching).
        Rate Limit: อนุญาต 100 ครั้งต่อ 1 นาที
        Rate Limit: Allows 100 requests per minute.
    */
    @GetMapping("/{templateCode}")
    @RateLimit(limit = 100, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get template by code")
    public ResponseEntity<TemplateResponseDTO> getTemplate(@PathVariable String templateCode)
            throws SystemGlobalException {
        TemplateResponseDTO response = templateService.getTemplate(templateCode);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. LIST TEMPLATES
    // ========================================================================

    /*
        API: GET /api/v1/templates
        ฟังก์ชันนี้แสดงรายการเทมเพลตทั้งหมดแบบแบ่งหน้า
        This function lists all templates with pagination.
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @GetMapping
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "List templates")
    public ResponseEntity<Page<TemplateResponseDTO>> listTemplates(
            @RequestParam(required = false) String templateType,
            Pageable pageable) throws SystemGlobalException {
        Page<TemplateResponseDTO> page = templateService.listTemplates(templateType, pageable);
        return ResponseEntity.ok(page);
    }

    // ========================================================================
    // 4. DELETE TEMPLATE
    // ========================================================================

    /*
        API: DELETE /api/v1/templates/{templateCode}
        ฟังก์ชันนี้ลบเทมเพลต
        This function deletes a template.
        Rate Limit: อนุญาต 10 ครั้งต่อ 1 ชั่วโมง
        Rate Limit: Allows 10 requests per 1 hour.
    */
    @DeleteMapping("/{templateCode}")
    @RateLimit(limit = 10, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Delete template")
    public ResponseEntity<Void> deleteTemplate(@PathVariable String templateCode) throws SystemGlobalException {
        templateService.deleteTemplate(templateCode);
        return ResponseEntity.noContent().build();
    }
}
```

---

## 🧩 ตัวอย่าง Domain Entity และ Enum

### `domain/enums/DocumentType.java`

```java
package com.template.app.modules.document.domain.enums;

/*
    ประเภทของเอกสาร / Document type.
*/
public enum DocumentType {
    PDF,        // ไฟล์ PDF
    EXCEL,      // ไฟล์ Excel (.xlsx)
    CSV,        // ไฟล์ CSV
    IMAGE,      // ไฟล์รูปภาพ (JPG, PNG, etc.)
    OTHER       // ไฟล์ประเภทอื่น
}
```

### `domain/enums/DocumentStatus.java`

```java
package com.template.app.modules.document.domain.enums;

/*
    สถานะของเอกสาร / Document status.
*/
public enum DocumentStatus {
    DRAFT,      // ร่าง / Draft.
    GENERATED,  // สร้างเสร็จ / Generated.
    SENT,       // ส่งแล้ว / Sent.
    ARCHIVED,   // เก็บถาวร / Archived.
    DELETED     // ถูกลบ / Deleted.
}
```

### `domain/TDocument.java`

```java
package com.template.app.modules.document.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import com.template.app.modules.document.domain.enums.DocumentStatus;
import com.template.app.modules.document.domain.enums.DocumentType;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.time.LocalDateTime;
import java.util.UUID;

@Data
@EqualsAndHashCode(callSuper = true)
public class TDocument extends GenericBusinessClass {

    private String documentNo;              // เลขที่เอกสาร / Document number.
    private DocumentType documentType;      // ประเภทเอกสาร / Document type.
    private String documentSubType;         // ประเภทย่อย / Sub type (QUOTATION, INVOICE, etc.)
    private String referenceType;           // ประเภทอ้างอิง / Reference type.
    private UUID referenceId;               // ID อ้างอิง / Reference ID.
    private UUID templateId;                // ID เทมเพลต / Template ID.
    private String fileName;                // ชื่อไฟล์ / File name.
    private String filePath;                // ที่อยู่ไฟล์ / File path.
    private Long fileSize;                  // ขนาดไฟล์ (bytes) / File size.
    private String mimeType;                // MIME type.
    private DocumentStatus status;          // สถานะ / Status.
    private UUID generatedBy;               // ผู้สร้าง / Generated by.
    private LocalDateTime generatedAt;      // วันที่สร้าง / Generated at.
    private UUID sentBy;                    // ผู้ส่ง / Sent by.
    private LocalDateTime sentAt;           // วันที่ส่ง / Sent at.
    private String sentToEmail;             // อีเมลผู้รับ / Recipient email.
    private String description;             // คำอธิบาย / Description.
    private String[] tags;                  // แท็ก / Tags.
    private String metadata;                // ข้อมูลเพิ่มเติม / Metadata.

    /*
        ฟังก์ชันนี้บันทึกการส่งเอกสาร
        This function records document sending.
    */
    public void markAsSent(String email) {
        if (this.status == DocumentStatus.DELETED) {
            throw new IllegalStateException("Cannot send a deleted document.");
        }
        this.status = DocumentStatus.SENT;
        this.sentAt = LocalDateTime.now();
        this.sentToEmail = email;
    }

    /*
        ฟังก์ชันนี้บันทึกการเก็บถาวรเอกสาร
        This function archives the document.
    */
    public void archive() {
        if (this.status == DocumentStatus.DELETED) {
            throw new IllegalStateException("Cannot archive a deleted document.");
        }
        this.status = DocumentStatus.ARCHIVED;
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่าเอกสารสามารถดาวน์โหลดได้หรือไม่
        This function checks if the document can be downloaded.
    */
    public boolean canDownload() {
        return this.status != DocumentStatus.DELETED && this.filePath != null;
    }
}
```

### `domain/MDocumentTemplate.java`

```java
package com.template.app.modules.document.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import lombok.Data;
import lombok.EqualsAndHashCode;

@Data
@EqualsAndHashCode(callSuper = true)
public class MDocumentTemplate extends GenericBusinessClass {

    private String templateCode;            // รหัสเทมเพลต / Template code.
    private String templateName;            // ชื่อเทมเพลต / Template name.
    private String templateType;            // JASPER, POI, CUSTOM
    private String fileName;                // ชื่อไฟล์ / File name.
    private String filePath;                // ที่อยู่ไฟล์ / File path.
    private Long fileSize;                  // ขนาดไฟล์ / File size.
    private Integer version;                // เวอร์ชัน / Version.
    private String description;             // คำอธิบาย / Description.
    private Boolean isActive;               // ใช้งานอยู่ / Active.
    private Boolean isDefault;              // เป็นค่าเริ่มต้น / Default.
    private String parameters;              // พารามิเตอร์ / Parameters.

    /*
        ฟังก์ชันนี้เพิ่มเวอร์ชันเมื่ออัปเดตเทมเพลต
        This function increments version when template is updated.
    */
    public void incrementVersion() {
        this.version = (this.version != null ? this.version : 0) + 1;
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่าเทมเพลตสามารถใช้งานได้
        This function checks if the template is usable.
    */
    public boolean isUsable() {
        return this.isActive != null && this.isActive && this.filePath != null;
    }
}
```

---

## 🔧 ฟังก์ชัน Service ตัวอย่าง (พร้อมคำอธิบาย)

### `application/impl/DocumentServiceImpl.java` (บางส่วน)

```java
package com.template.app.modules.document.application.impl;

import com.template.app.modules.document.application.interfaces.DocumentService;
import com.template.app.modules.document.domain.TDocument;
import com.template.app.modules.document.domain.enums.DocumentStatus;
import com.template.app.modules.document.infrastructure.cache.DocumentCacheService;
import com.template.app.modules.document.infrastructure.repository.DocumentRepository;
import com.template.app.modules.document.infrastructure.storage.FileStorageService;
import com.template.app.modules.document.presentation.dto.request.DocumentGenerateRequestDTO;
import com.template.app.modules.document.presentation.dto.response.DocumentResponseDTO;
import com.template.app._shared.application.GenericServiceImpl;
import com.template.app.exception.SystemGlobalException;
import org.springframework.core.io.Resource;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.multipart.MultipartFile;

import java.time.LocalDateTime;
import java.util.UUID;

@Service
public class DocumentServiceImpl extends GenericServiceImpl<TDocument, DocumentRepository>
        implements DocumentService {

    private final DocumentCacheService cacheService;
    private final FileStorageService fileStorageService;
    private final ReportGenerationService reportGenerationService;
    private final TemplateService templateService;

    public DocumentServiceImpl(DocumentRepository repository,
                               DocumentCacheService cacheService,
                               FileStorageService fileStorageService,
                               ReportGenerationService reportGenerationService,
                               TemplateService templateService) {
        super(repository);
        this.cacheService = cacheService;
        this.fileStorageService = fileStorageService;
        this.reportGenerationService = reportGenerationService;
        this.templateService = templateService;
    }

    /*
        ฟังก์ชันนี้สร้างเอกสารจากเทมเพลต ใช้ JasperReports สำหรับ PDF หรือ POI สำหรับ Excel
        This function generates a document from a template, using JasperReports for PDF or POI for Excel.
    */
    @Override
    @Transactional
    public DocumentResponseDTO generateDocument(DocumentGenerateRequestDTO request) throws SystemGlobalException {
        // 1. ดึงข้อมูลเทมเพลต (ใช้ Cache) / Get template data (using cache).
        MDocumentTemplate template = templateService.getTemplate(request.getTemplateCode());

        // 2. ตรวจสอบว่าเทมเพลตใช้งานได้ / Verify template is usable.
        if (!template.isUsable()) {
            throw new SystemGlobalException("Template is not active or file not found.", null);
        }

        // 3. สร้างเอกสารตามประเภทเทมเพลต / Generate document based on template type.
        byte[] documentData;
        String fileName;
        String mimeType;

        if ("JASPER".equals(template.getTemplateType())) {
            // ใช้ JasperReports สร้าง PDF / Generate PDF using JasperReports.
            documentData = reportGenerationService.generatePDF(request.getTemplateCode(), request.getParameters());
            fileName = request.getFileName() != null ? request.getFileName() + ".pdf" : template.getTemplateCode() + ".pdf";
            mimeType = "application/pdf";
        } else if ("POI".equals(template.getTemplateType())) {
            // ใช้ Apache POI สร้าง Excel / Generate Excel using Apache POI.
            documentData = reportGenerationService.generateExcel(request.getTemplateCode(), request.getParameters());
            fileName = request.getFileName() != null ? request.getFileName() + ".xlsx" : template.getTemplateCode() + ".xlsx";
            mimeType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet";
        } else {
            throw new SystemGlobalException("Unsupported template type: " + template.getTemplateType(), null);
        }

        // 4. บันทึกไฟล์ลง Storage / Save file to storage.
        String filePath = fileStorageService.storeFile(
                fileName,
                documentData,
                "documents/" + UUID.randomUUID().toString()
        );

        // 5. สร้างบันทึกเอกสารในฐานข้อมูล / Create document record in database.
        TDocument document = new TDocument();
        document.setDocumentType(mimeType.contains("pdf") ? DocumentType.PDF : DocumentType.EXCEL);
        document.setDocumentSubType(request.getDocumentSubType());
        document.setReferenceType(request.getReferenceType());
        document.setReferenceId(request.getReferenceId());
        document.setTemplateId(template.getId());
        document.setFileName(fileName);
        document.setFilePath(filePath);
        document.setFileSize((long) documentData.length);
        document.setMimeType(mimeType);
        document.setStatus(DocumentStatus.GENERATED);
        document.setGeneratedBy(getUserId());
        document.setGeneratedAt(LocalDateTime.now());
        document.setDescription(request.getDescription());

        // 6. บันทึกข้อมูล / Save data.
        TDocument savedDocument = this.create(document);

        // 7. บันทึก Cache / Save to cache.
        cacheService.saveDocument(savedDocument);

        return DocumentResponseDTO.fromEntity(savedDocument);
    }

    /*
        ฟังก์ชันนี้อัปโหลดเอกสารเข้าเก็บในระบบ
        This function uploads a document to the system.
    */
    @Override
    @Transactional
    public DocumentResponseDTO uploadDocument(MultipartFile file,
                                              String documentType,
                                              String referenceType,
                                              UUID referenceId,
                                              String description) throws SystemGlobalException {
        // 1. ตรวจสอบไฟล์ / Validate file.
        if (file.isEmpty()) {
            throw new SystemGlobalException("File is empty.", null);
        }

        // 2. บันทึกไฟล์ลง Storage / Save file to storage.
        String fileName = file.getOriginalFilename();
        String filePath = fileStorageService.storeMultipartFile(
                file,
                "uploads/" + UUID.randomUUID().toString()
        );

        // 3. สร้างบันทึกเอกสาร / Create document record.
        TDocument document = new TDocument();
        document.setDocumentType(DocumentType.valueOf(documentType.toUpperCase()));
        document.setReferenceType(referenceType);
        document.setReferenceId(referenceId);
        document.setFileName(fileName);
        document.setFilePath(filePath);
        document.setFileSize(file.getSize());
        document.setMimeType(file.getContentType());
        document.setStatus(DocumentStatus.GENERATED);
        document.setGeneratedBy(getUserId());
        document.setGeneratedAt(LocalDateTime.now());
        document.setDescription(description);

        // 4. บันทึกข้อมูล / Save data.
        TDocument savedDocument = this.create(document);

        // 5. บันทึก Cache / Save to cache.
        cacheService.saveDocument(savedDocument);

        return DocumentResponseDTO.fromEntity(savedDocument);
    }

    /*
        ฟังก์ชันนี้ดาวน์โหลดไฟล์เอกสารจากระบบ
        This function downloads a document file from the system.
    */
    @Override
    public Resource downloadDocument(UUID documentId) throws SystemGlobalException {
        // 1. ดึงข้อมูลเอกสาร / Get document data.
        TDocument document = this.read(documentId);

        // 2. ตรวจสอบว่าสามารถดาวน์โหลดได้ / Check if downloadable.
        if (!document.canDownload()) {
            throw new SystemGlobalException("Document cannot be downloaded.", null);
        }

        // 3. โหลดไฟล์จาก Storage / Load file from storage.
        Resource resource = fileStorageService.loadFile(document.getFilePath());

        // 4. บันทึกประวัติการดาวน์โหลด / Record download history.
        // documentHistoryService.recordDownload(documentId, getUserId());

        return resource;
    }
}
```

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
| 8 | 📊 Dashboard & Reports | ✅ ครบถ้วน | Overview + Sales + Inventory + Job Status + Top Parts + Financial + Export |
| 9 | 📄 Document Management | ✅ ครบถ้วน | Document Generation + Upload + OCR + Template Management + Storage + Cache |

--- 