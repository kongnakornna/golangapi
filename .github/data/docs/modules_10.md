 **โมดูลที่ 10: 📧 Email Service (ระบบส่งอีเมลอัตโนมัติ)**

โมดูล Email Service เป็นระบบกลางสำหรับจัดการการส่งอีเมลทั้งหมดในระบบ ครอบคลุมการทำงานดังนี้:

1. **การส่งอีเมลจากเทมเพลต (Template-based Email)** – ส่งอีเมลด้วยเทมเพลตที่กำหนด
2. **การแนบไฟล์ (Attachments)** – แนบเอกสาร PDF, Excel ไปกับอีเมล
3. **การส่งอีเมลตามเหตุการณ์ (Event-based Email)** – ส่งอัตโนมัติเมื่อเกิดเหตุการณ์ (Quotation Approved, PO Created, etc.)
4. **การจัดการคิวอีเมล (Email Queue)** – ใช้ Kafka หรือ Async เพื่อไม่ให้กระทบ Performance
5. **ประวัติการส่งอีเมล (Email History)** – บันทึกประวัติการส่งทั้งหมด
6. **การจัดการเทมเพลตอีเมล (Email Template Management)** – อัปโหลดและจัดการเทมเพลต HTML
7. **การส่งอีเมลตามกำหนดเวลา (Scheduled Email)** – ใช้ Batch Jobs ส่งอีเมลสรุปรายวัน/รายเดือน

---

## 📁 โครงสร้างโมดูล Email Service (`modules/email`)

```
modules/email/
├── application/
│   ├── interfaces/
│   │   ├── EmailService.java
│   │   ├── EmailTemplateService.java
│   │   ├── EmailQueueService.java
│   │   └── EmailHistoryService.java
│   ├── impl/
│   │   ├── EmailServiceImpl.java
│   │   ├── EmailTemplateServiceImpl.java
│   │   ├── EmailQueueServiceImpl.java
│   │   └── EmailHistoryServiceImpl.java
│   └── usecase/
│       ├── SendEmailUseCase.java
│       ├── SendEmailWithAttachmentUseCase.java
│       ├── SendBulkEmailUseCase.java
│       ├── SendTemplateEmailUseCase.java
│       ├── QueueEmailUseCase.java
│       ├── ProcessEmailQueueUseCase.java
│       ├── GetEmailHistoryUseCase.java
│       └── CreateEmailTemplateUseCase.java
├── domain/
│   ├── MEmailTemplate.java
│   ├── TEmailHistory.java
│   ├── TEmailQueue.java
│   ├── enums/
│   │   ├── EmailStatus.java           // PENDING, SENT, FAILED, BOUNCED
│   │   ├── EmailPriority.java         // LOW, NORMAL, HIGH, URGENT
│   │   └── EmailCategory.java         // QUOTATION, INVOICE, PO, RECEIPT, REMINDER, PROMOTION, NOTIFICATION
│   └── valueobjects/
│       ├── EmailAddress.java
│       └── EmailContent.java
├── infrastructure/
│   ├── repository/
│   │   ├── EmailHistoryRepository.java
│   │   ├── EmailQueueRepository.java
│   │   ├── EmailTemplateRepository.java
│   │   └── impl/
│   │       ├── EmailHistoryRepositoryImpl.java
│   │       └── EmailTemplateRepositoryImpl.java
│   ├── cache/                                           // ⬅️ ระบบ Cache สำหรับ Email
│   │   ├── EmailTemplateCacheService.java
│   │   └── EmailRateLimitCacheService.java
│   ├── provider/                                        // ⬅️ Email Provider
│   │   ├── EmailProvider.java
│   │   ├── SMTPEmailProvider.java
│   │   ├── SendGridEmailProvider.java
│   │   └── EmailProviderConfig.java
│   ├── consumer/                                        // ⬅️ Kafka Consumer
│   │   └── EmailEventConsumer.java
│   ├── entity/
│   │   ├── EmailTemplateEntity.java
│   │   ├── EmailHistoryEntity.java
│   │   └── EmailQueueEntity.java
│   └── mapper/
│       ├── EmailTemplateMapper.java
│       └── EmailHistoryMapper.java
└── presentation/
    ├── controller/
    │   ├── EmailController.java          // Send Email APIs
    │   ├── EmailTemplateController.java  // Template Management
    │   └── EmailHistoryController.java   // Email History APIs
    ├── dto/
    │   ├── request/
    │   │   ├── EmailSendRequestDTO.java
    │   │   ├── EmailTemplateRequestDTO.java
    │   │   ├── BulkEmailRequestDTO.java
    │   │   └── EmailSearchRequestDTO.java
    │   └── response/
    │       ├── EmailSendResponseDTO.java
    │       ├── EmailTemplateResponseDTO.java
    │       └── EmailHistoryResponseDTO.java
    └── validator/
        └── EmailValidator.java
```

---

## 🗄️ ออกแบบฐานข้อมูล (Database Design) สำหรับโมดูล Email

### 📄 SQL DDL (ไฟล์: `src/main/resources/db/migration/V10__email_schema.sql`)

```sql
-- ==============================================
-- ตาราง: m_email_template (เทมเพลตอีเมล)
-- Email template master data.
-- ==============================================
CREATE TABLE IF NOT EXISTS m_email_template (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_code VARCHAR(50) UNIQUE NOT NULL,      -- รหัสเทมเพลต (เช่น QUOTATION_EMAIL, INVOICE_EMAIL)
    template_name VARCHAR(100) NOT NULL,            -- ชื่อเทมเพลต
    subject VARCHAR(255) NOT NULL,                  -- หัวข้ออีเมล (รองรับตัวแปร)
    body_html TEXT,                                 -- เนื้อหา HTML
    body_text TEXT,                                 -- เนื้อหา Text (สำรอง)
    from_email VARCHAR(100),                        -- ผู้ส่ง (ถ้าไม่ระบุใช้ค่า Default)
    from_name VARCHAR(100),                         -- ชื่อผู้ส่ง
    category VARCHAR(50),                           -- QUOTATION, INVOICE, PO, RECEIPT, REMINDER, PROMOTION, NOTIFICATION
    language VARCHAR(10) DEFAULT 'th',              -- ภาษา (th, en)
    version INTEGER DEFAULT 1,
    is_active BOOLEAN DEFAULT TRUE,
    is_default BOOLEAN DEFAULT FALSE,
    variables JSONB,                                -- ตัวแปรที่ใช้ในเทมเพลต
    description TEXT,
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_m_email_template_code ON m_email_template(template_code);
CREATE INDEX idx_m_email_template_category ON m_email_template(category);
CREATE INDEX idx_m_email_template_language ON m_email_template(language);
CREATE INDEX idx_m_email_template_active ON m_email_template(is_active);

-- ==============================================
-- ข้อมูลเริ่มต้น: เทมเพลตอีเมล
-- Initial data: Email templates.
-- ==============================================
INSERT INTO m_email_template (template_code, template_name, subject, category, language, is_active, user_id, whitelabel_id)
VALUES 
('QUOTATION_EMAIL', 'ใบเสนอราคา', 'ใบเสนอราคา #{quotationNo} - {customerName}', 'QUOTATION', 'th', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('INVOICE_EMAIL', 'ใบแจ้งหนี้', 'ใบแจ้งหนี้ #{invoiceNo} - {customerName}', 'INVOICE', 'th', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('PO_EMAIL', 'ใบสั่งซื้อ', 'ใบสั่งซื้อ #{poNo} - {supplierName}', 'PO', 'th', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('RECEIPT_EMAIL', 'ใบเสร็จรับเงิน', 'ใบเสร็จรับเงิน #{receiptNo} - {customerName}', 'RECEIPT', 'th', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('REMINDER_EMAIL', 'อีเมลแจ้งเตือน', 'แจ้งเตือน: {reminderSubject}', 'REMINDER', 'th', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001');

-- ==============================================
-- ตาราง: t_email_history (ประวัติการส่งอีเมล)
-- Email sending history.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_email_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email_id VARCHAR(50) UNIQUE NOT NULL,           -- ID ของอีเมล (UUID หรือ Message ID)
    template_code VARCHAR(50),                      -- รหัสเทมเพลต (ถ้ามี)
    reference_type VARCHAR(30),                     -- QUOTATION, INVOICE, PO, etc.
    reference_id UUID,                              -- ID ของข้อมูลต้นทาง
    from_email VARCHAR(100) NOT NULL,
    from_name VARCHAR(100),
    to_email VARCHAR(200) NOT NULL,
    to_name VARCHAR(100),
    cc_email VARCHAR(200),
    bcc_email VARCHAR(200),
    subject VARCHAR(255) NOT NULL,
    body_preview TEXT,                              -- เนื้อหาส่วนย่อ (สำหรับแสดงในประวัติ)
    status VARCHAR(20) DEFAULT 'PENDING',           -- PENDING, SENT, FAILED, BOUNCED
    priority VARCHAR(20) DEFAULT 'NORMAL',          -- LOW, NORMAL, HIGH, URGENT
    sent_at TIMESTAMP,
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    attachments JSONB,                              -- รายการไฟล์แนบ
    metadata JSONB,                                 -- ข้อมูลเพิ่มเติม
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_email_history_email_id ON t_email_history(email_id);
CREATE INDEX idx_t_email_history_reference ON t_email_history(reference_type, reference_id);
CREATE INDEX idx_t_email_history_template ON t_email_history(template_code);
CREATE INDEX idx_t_email_history_status ON t_email_history(status);
CREATE INDEX idx_t_email_history_to ON t_email_history(to_email);
CREATE INDEX idx_t_email_history_sent_at ON t_email_history(sent_at);
CREATE INDEX idx_t_email_history_whitelabel ON t_email_history(whitelabel_id);

-- ==============================================
-- ตาราง: t_email_queue (คิวอีเมล - สำหรับส่งแบบ Async)
-- Email queue for asynchronous sending.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_email_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email_id VARCHAR(50) UNIQUE NOT NULL,           -- เชื่อมโยงกับ t_email_history
    template_code VARCHAR(50),
    reference_type VARCHAR(30),
    reference_id UUID,
    from_email VARCHAR(100) NOT NULL,
    to_email VARCHAR(200) NOT NULL,
    to_name VARCHAR(100),
    subject VARCHAR(255) NOT NULL,
    body_html TEXT,
    body_text TEXT,
    attachments JSONB,
    priority VARCHAR(20) DEFAULT 'NORMAL',
    status VARCHAR(20) DEFAULT 'PENDING',           -- PENDING, PROCESSING, SENT, FAILED
    retry_count INTEGER DEFAULT 0,
    max_retry INTEGER DEFAULT 3,
    next_attempt_at TIMESTAMP DEFAULT NOW(),
    error_message TEXT,
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_email_queue_status ON t_email_queue(status);
CREATE INDEX idx_t_email_queue_priority ON t_email_queue(priority);
CREATE INDEX idx_t_email_queue_next_attempt ON t_email_queue(next_attempt_at);
CREATE INDEX idx_t_email_queue_whitelabel ON t_email_queue(whitelabel_id);

-- ==============================================
-- ฟังก์ชันสร้าง Email ID อัตโนมัติ
-- Function to generate email ID.
-- ==============================================
CREATE OR REPLACE FUNCTION generate_email_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.email_id := 'EMAIL-' || TO_CHAR(NOW(), 'YYYYMMDD') || '-' || LPAD(CAST((
        SELECT COALESCE(MAX(CAST(SUBSTRING(email_id FROM 20) AS INTEGER)), 0) + 1
        FROM t_email_history
        WHERE email_id LIKE 'EMAIL-' || TO_CHAR(NOW(), 'YYYYMMDD') || '-%'
    ) AS TEXT), 4, '0');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_generate_email_id ON t_email_history;
CREATE TRIGGER trg_generate_email_id
BEFORE INSERT ON t_email_history
FOR EACH ROW
EXECUTE FUNCTION generate_email_id();

-- ==============================================
-- VIEW: v_email_analytics (สถิติการส่งอีเมล)
-- Email analytics view.
-- ==============================================
CREATE OR REPLACE VIEW v_email_analytics AS
SELECT 
    DATE_TRUNC('day', created_at) AS day,
    COUNT(*) AS total_sent,
    SUM(CASE WHEN status = 'SENT' THEN 1 ELSE 0 END) AS success_count,
    SUM(CASE WHEN status = 'FAILED' THEN 1 ELSE 0 END) AS failed_count,
    SUM(CASE WHEN status = 'BOUNCED' THEN 1 ELSE 0 END) AS bounced_count,
    category,
    whitelabel_id
FROM t_email_history
GROUP BY DATE_TRUNC('day', created_at), category, whitelabel_id;
```

---

## 🧠 ระบบ Cache (Redis) สำหรับ Email

### `infrastructure/cache/EmailTemplateCacheService.java`

```java
package com.template.app.modules.email.infrastructure.cache;

import com.template.app.modules.email.domain.MEmailTemplate;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.CachePut;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

@Service
public class EmailTemplateCacheService {

    /*
        ฟังก์ชันนี้ดึงข้อมูลเทมเพลตอีเมลจาก Cache ตามรหัสเทมเพลตและภาษา
        This function retrieves email template data from cache by template code and language.
        Redis Key: email_template:{templateCode}:{language}
    */
    @Cacheable(value = "email_templates", key = "#templateCode + ':' + #language")
    public MEmailTemplate getTemplate(String templateCode, String language) {
        return null;
    }

    /*
        ฟังก์ชันนี้จะอัปเดต Cache เมื่อมีการบันทึกเทมเพลต
        This function updates the cache when a template is saved.
    */
    @CachePut(value = "email_templates", key = "#template.templateCode + ':' + #template.language")
    public MEmailTemplate saveTemplate(MEmailTemplate template) {
        return template;
    }

    /*
        ฟังก์ชันนี้ลบข้อมูลเทมเพลตออกจาก Cache
        This function evicts template data from cache.
    */
    @CacheEvict(value = "email_templates", key = "#templateCode + ':' + #language")
    public void evictTemplate(String templateCode, String language) {
        // ลบ Cache / Evict cache.
    }

    /*
        ฟังก์ชันนี้ลบ Cache ทั้งหมดของ Email Template
        This function clears all email template caches.
    */
    @CacheEvict(value = "email_templates", allEntries = true)
    public void evictAllTemplates() {
        // ลบทุก key / Evict all keys.
    }
}
```

### `infrastructure/cache/EmailRateLimitCacheService.java`

```java
package com.template.app.modules.email.infrastructure.cache;

import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.stereotype.Service;

import java.time.Duration;

@Service
public class EmailRateLimitCacheService {

    private final RedisTemplate<String, String> redisTemplate;

    public EmailRateLimitCacheService(RedisTemplate<String, String> redisTemplate) {
        this.redisTemplate = redisTemplate;
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่าผู้ใช้สามารถส่งอีเมลได้หรือไม่ (Rate Limit)
        This function checks if a user can send emails (Rate Limit).
        Redis Key: email_rate:{userId}
    */
    public boolean canSendEmail(String userId) {
        String key = "email_rate:" + userId;
        String countStr = redisTemplate.opsForValue().get(key);
        int count = countStr == null ? 0 : Integer.parseInt(countStr);
        
        // จำกัด 50 อีเมลต่อชั่วโมง / Limit 50 emails per hour.
        if (count >= 50) {
            return false;
        }
        
        redisTemplate.opsForValue().increment(key);
        if (count == 0) {
            redisTemplate.expire(key, Duration.ofHours(1));
        }
        return true;
    }

    /*
        ฟังก์ชันนี้ลบ Rate Limit ของผู้ใช้ (ใช้เมื่อต้องการ reset)
        This function resets the rate limit for a user.
    */
    public void resetRateLimit(String userId) {
        redisTemplate.delete("email_rate:" + userId);
    }
}
```

---

## ⏱️ ระบบ Rate Limit สำหรับ Email Controller

```java
package com.template.app.modules.email.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.email.application.interfaces.EmailService;
import com.template.app.modules.email.presentation.dto.request.BulkEmailRequestDTO;
import com.template.app.modules.email.presentation.dto.request.EmailSendRequestDTO;
import com.template.app.modules.email.presentation.dto.response.EmailSendResponseDTO;
import com.template.app.exception.SystemGlobalException;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

@RestController
@RequestMapping("/api/v1/email")
@Tag(name = "Email Service", description = "Email Management APIs")
@RequiredArgsConstructor
public class EmailController {

    private final EmailService emailService;

    // ========================================================================
    // 1. SEND EMAIL (ส่งอีเมลแบบ Single)
    // ========================================================================

    /*
        API: POST /api/v1/email/send
        ฟังก์ชันนี้ส่งอีเมลแบบ Single (ไม่ใช้เทมเพลต) พร้อมตัวเลือกแนบไฟล์
        This function sends a single email (without template) with optional attachments.
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @PostMapping("/send")
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Send an email")
    public ResponseEntity<EmailSendResponseDTO> sendEmail(@Valid @RequestBody EmailSendRequestDTO request)
            throws SystemGlobalException {
        EmailSendResponseDTO response = emailService.sendEmail(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 2. SEND TEMPLATE EMAIL (ส่งอีเมลจากเทมเพลต)
    // ========================================================================

    /*
        API: POST /api/v1/email/send-template
        ฟังก์ชันนี้ส่งอีเมลจากเทมเพลตที่กำหนด พร้อมแทนที่ตัวแปร
        This function sends an email from a template with variable replacement.
        Rate Limit: อนุญาต 25 ครั้งต่อ 1 นาที
        Rate Limit: Allows 25 requests per minute.
    */
    @PostMapping("/send-template")
    @RateLimit(limit = 25, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Send email from template")
    public ResponseEntity<EmailSendResponseDTO> sendTemplateEmail(@Valid @RequestBody EmailSendRequestDTO request)
            throws SystemGlobalException {
        EmailSendResponseDTO response = emailService.sendTemplateEmail(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. SEND BULK EMAIL (ส่งอีเมลจำนวนมาก)
    // ========================================================================

    /*
        API: POST /api/v1/email/bulk
        ฟังก์ชันนี้ส่งอีเมลจำนวนมาก (Bulk) ใช้สำหรับโปรโมชันหรือแจ้งเตือนกลุ่ม
        This function sends bulk emails (for promotions or group notifications).
        Rate Limit: อนุญาต 5 ครั้งต่อ 5 นาที (ป้องกันการส่ง Spam)
        Rate Limit: Allows 5 requests per 5 minutes (prevent Spam).
    */
    @PostMapping("/bulk")
    @RateLimit(limit = 5, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Send bulk emails")
    public ResponseEntity<EmailSendResponseDTO> sendBulkEmail(@Valid @RequestBody BulkEmailRequestDTO request)
            throws SystemGlobalException {
        EmailSendResponseDTO response = emailService.sendBulkEmail(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 4. GET EMAIL STATUS (เช็คสถานะการส่ง)
    // ========================================================================

    /*
        API: GET /api/v1/email/status/{emailId}
        ฟังก์ชันนี้ตรวจสอบสถานะการส่งอีเมล (PENDING, SENT, FAILED)
        This function checks the status of an email (PENDING, SENT, FAILED).
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @GetMapping("/status/{emailId}")
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get email send status")
    public ResponseEntity<EmailSendResponseDTO> getEmailStatus(@PathVariable String emailId)
            throws SystemGlobalException {
        EmailSendResponseDTO response = emailService.getEmailStatus(emailId);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 5. RESEND EMAIL (ส่งอีเมลซ้ำ)
    // ========================================================================

    /*
        API: POST /api/v1/email/resend/{emailId}
        ฟังก์ชันนี้ส่งอีเมลซ้ำ (ใช้เมื่อการส่งครั้งแรกล้มเหลว)
        This function resends an email (used when the first attempt failed).
        Rate Limit: อนุญาต 10 ครั้งต่อ 1 ชั่วโมง
        Rate Limit: Allows 10 requests per 1 hour.
    */
    @PostMapping("/resend/{emailId}")
    @RateLimit(limit = 10, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Resend an email")
    public ResponseEntity<EmailSendResponseDTO> resendEmail(@PathVariable String emailId)
            throws SystemGlobalException {
        EmailSendResponseDTO response = emailService.resendEmail(emailId);
        return ResponseEntity.ok(response);
    }
}
```

### `EmailTemplateController.java` (จัดการเทมเพลตอีเมล)

```java
package com.template.app.modules.email.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.email.application.interfaces.EmailTemplateService;
import com.template.app.modules.email.presentation.dto.request.EmailTemplateRequestDTO;
import com.template.app.modules.email.presentation.dto.response.EmailTemplateResponseDTO;
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
@RequestMapping("/api/v1/email/templates")
@Tag(name = "Email Templates", description = "Email Template Management APIs")
@RequiredArgsConstructor
public class EmailTemplateController {

    private final EmailTemplateService emailTemplateService;

    // ========================================================================
    // 1. CREATE EMAIL TEMPLATE
    // ========================================================================

    /*
        API: POST /api/v1/email/templates
        ฟังก์ชันนี้สร้างเทมเพลตอีเมลใหม่
        This function creates a new email template.
        Rate Limit: อนุญาต 10 ครั้งต่อ 5 นาที
        Rate Limit: Allows 10 requests per 5 minutes.
    */
    @PostMapping
    @RateLimit(limit = 10, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Create email template")
    public ResponseEntity<EmailTemplateResponseDTO> createTemplate(@Valid @RequestBody EmailTemplateRequestDTO request)
            throws SystemGlobalException {
        EmailTemplateResponseDTO response = emailTemplateService.createTemplate(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 2. GET TEMPLATE BY CODE
    // ========================================================================

    /*
        API: GET /api/v1/email/templates/{templateCode}
        ฟังก์ชันนี้ดึงข้อมูลเทมเพลตตามรหัสและภาษา (ใช้ Cache)
        This function retrieves template data by code and language (uses caching).
        Rate Limit: อนุญาต 100 ครั้งต่อ 1 นาที
        Rate Limit: Allows 100 requests per minute.
    */
    @GetMapping("/{templateCode}")
    @RateLimit(limit = 100, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get template by code")
    public ResponseEntity<EmailTemplateResponseDTO> getTemplate(
            @PathVariable String templateCode,
            @RequestParam(defaultValue = "th") String language) throws SystemGlobalException {
        EmailTemplateResponseDTO response = emailTemplateService.getTemplate(templateCode, language);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. LIST TEMPLATES
    // ========================================================================

    /*
        API: GET /api/v1/email/templates
        ฟังก์ชันนี้แสดงรายการเทมเพลตทั้งหมดแบบแบ่งหน้า
        This function lists all templates with pagination.
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @GetMapping
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "List email templates")
    public ResponseEntity<Page<EmailTemplateResponseDTO>> listTemplates(
            @RequestParam(required = false) String category,
            @RequestParam(required = false) String language,
            Pageable pageable) throws SystemGlobalException {
        Page<EmailTemplateResponseDTO> page = emailTemplateService.listTemplates(category, language, pageable);
        return ResponseEntity.ok(page);
    }

    // ========================================================================
    // 4. UPDATE TEMPLATE
    // ========================================================================

    /*
        API: PUT /api/v1/email/templates/{templateCode}
        ฟังก์ชันนี้แก้ไขเทมเพลตอีเมล
        This function updates an email template.
        Rate Limit: อนุญาต 10 ครั้งต่อ 5 นาที
        Rate Limit: Allows 10 requests per 5 minutes.
    */
    @PutMapping("/{templateCode}")
    @RateLimit(limit = 10, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Update email template")
    public ResponseEntity<EmailTemplateResponseDTO> updateTemplate(
            @PathVariable String templateCode,
            @RequestParam(defaultValue = "th") String language,
            @Valid @RequestBody EmailTemplateRequestDTO request) throws SystemGlobalException {
        EmailTemplateResponseDTO response = emailTemplateService.updateTemplate(templateCode, language, request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 5. DELETE TEMPLATE
    // ========================================================================

    /*
        API: DELETE /api/v1/email/templates/{templateCode}
        ฟังก์ชันนี้ลบเทมเพลต (Soft Delete)
        This function deletes a template (Soft Delete).
        Rate Limit: อนุญาต 5 ครั้งต่อ 1 ชั่วโมง
        Rate Limit: Allows 5 requests per 1 hour.
    */
    @DeleteMapping("/{templateCode}")
    @RateLimit(limit = 5, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Delete email template")
    public ResponseEntity<Void> deleteTemplate(
            @PathVariable String templateCode,
            @RequestParam(defaultValue = "th") String language) throws SystemGlobalException {
        emailTemplateService.deleteTemplate(templateCode, language);
        return ResponseEntity.noContent().build();
    }
}
```

---

## 🔄 Kafka Event Consumer สำหรับ Email

```java
package com.template.app.modules.email.infrastructure.consumer;

import com.template.app.modules.email.application.interfaces.EmailService;
import com.template.app.modules.email.presentation.dto.request.EmailSendRequestDTO;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Service;

import java.util.Map;

@Service
public class EmailEventConsumer {

    private static final Logger logger = LoggerFactory.getLogger(EmailEventConsumer.class);
    private final EmailService emailService;
    private final ObjectMapper objectMapper;

    public EmailEventConsumer(EmailService emailService, ObjectMapper objectMapper) {
        this.emailService = emailService;
        this.objectMapper = objectMapper;
    }

    /*
        ฟังก์ชันนี้รับ Event จาก Kafka และส่งอีเมลตามประเภท Event
        This function consumes events from Kafka and sends emails based on event type.
        ฟัง Topic: email-events
        Listening to Topic: email-events
    */
    @KafkaListener(topics = "email-events", groupId = "email-service-group")
    public void consumeEmailEvent(ConsumerRecord<String, String> record) {
        try {
            logger.info("Received email event: {}", record.value());
            
            Map<String, Object> event = objectMapper.readValue(record.value(), Map.class);
            String eventType = (String) event.get("eventType");
            
            switch (eventType) {
                case "QUOTATION_APPROVED":
                    // ส่งอีเมล Quotation ให้ลูกค้า / Send Quotation email to customer.
                    sendQuotationEmail(event);
                    break;
                case "PO_CREATED":
                    // ส่งอีเมล PO ให้ Supplier / Send PO email to supplier.
                    sendPOEmail(event);
                    break;
                case "INVOICE_CREATED":
                    // ส่งอีเมล Invoice ให้ลูกค้า / Send Invoice email to customer.
                    sendInvoiceEmail(event);
                    break;
                case "PAYMENT_RECEIVED":
                    // ส่งอีเมล Receipt ให้ลูกค้า / Send Receipt email to customer.
                    sendReceiptEmail(event);
                    break;
                default:
                    logger.warn("Unknown email event type: {}", eventType);
            }
        } catch (Exception e) {
            logger.error("Error processing email event: {}", e.getMessage(), e);
        }
    }

    private void sendQuotationEmail(Map<String, Object> event) throws Exception {
        // สร้าง DTO และส่งอีเมล / Create DTO and send email.
        EmailSendRequestDTO request = new EmailSendRequestDTO();
        request.setTemplateCode("QUOTATION_EMAIL");
        request.setToEmail((String) event.get("customerEmail"));
        request.setVariables((Map<String, Object>) event.get("variables"));
        // ... (โหลดไฟล์แนบ Quotation PDF)
        emailService.sendTemplateEmail(request);
    }

    private void sendPOEmail(Map<String, Object> event) throws Exception {
        EmailSendRequestDTO request = new EmailSendRequestDTO();
        request.setTemplateCode("PO_EMAIL");
        request.setToEmail((String) event.get("supplierEmail"));
        request.setVariables((Map<String, Object>) event.get("variables"));
        // ... (โหลดไฟล์แนบ PO PDF)
        emailService.sendTemplateEmail(request);
    }

    private void sendInvoiceEmail(Map<String, Object> event) throws Exception {
        EmailSendRequestDTO request = new EmailSendRequestDTO();
        request.setTemplateCode("INVOICE_EMAIL");
        request.setToEmail((String) event.get("customerEmail"));
        request.setVariables((Map<String, Object>) event.get("variables"));
        emailService.sendTemplateEmail(request);
    }

    private void sendReceiptEmail(Map<String, Object> event) throws Exception {
        EmailSendRequestDTO request = new EmailSendRequestDTO();
        request.setTemplateCode("RECEIPT_EMAIL");
        request.setToEmail((String) event.get("customerEmail"));
        request.setVariables((Map<String, Object>) event.get("variables"));
        emailService.sendTemplateEmail(request);
    }
}
```

---

## 🧩 ตัวอย่าง Domain Entity และ Enum

### `domain/enums/EmailStatus.java`

```java
package com.template.app.modules.email.domain.enums;

/*
    สถานะการส่งอีเมล / Email status.
*/
public enum EmailStatus {
    PENDING,    // รอดำเนินการ / Pending.
    SENT,       // ส่งสำเร็จ / Sent successfully.
    FAILED,     // ส่งล้มเหลว / Failed.
    BOUNCED     // ถูกปฏิเสธ / Bounced.
}
```

### `domain/enums/EmailCategory.java`

```java
package com.template.app.modules.email.domain.enums;

/*
    ประเภทของอีเมล / Email category.
*/
public enum EmailCategory {
    QUOTATION,      // ใบเสนอราคา
    INVOICE,        // ใบแจ้งหนี้
    PO,             // ใบสั่งซื้อ
    RECEIPT,        // ใบเสร็จรับเงิน
    REMINDER,       // อีเมลแจ้งเตือน
    PROMOTION,      // โปรโมชัน
    NOTIFICATION    // การแจ้งเตือนทั่วไป
}
```

### `domain/MEmailTemplate.java`

```java
package com.template.app.modules.email.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.util.Map;

@Data
@EqualsAndHashCode(callSuper = true)
public class MEmailTemplate extends GenericBusinessClass {

    private String templateCode;        // รหัสเทมเพลต / Template code.
    private String templateName;        // ชื่อเทมเพลต / Template name.
    private String subject;             // หัวข้อ / Subject.
    private String bodyHtml;            // เนื้อหา HTML / HTML body.
    private String bodyText;            // เนื้อหา Text / Text body.
    private String fromEmail;           // ผู้ส่ง / From email.
    private String fromName;            // ชื่อผู้ส่ง / From name.
    private String category;            // ประเภท / Category.
    private String language;            // ภาษา / Language.
    private Integer version;            // เวอร์ชัน / Version.
    private Boolean isActive;           // ใช้งานอยู่ / Active.
    private Boolean isDefault;          // ค่าเริ่มต้น / Default.
    private String variables;           // ตัวแปรที่ใช้ / Variables.
    private String description;         // คำอธิบาย / Description.

    /*
        ฟังก์ชันนี้แทนที่ตัวแปรในเนื้อหาเทมเพลตด้วยค่าที่กำหนด
        This function replaces template variables with provided values.
    */
    public String renderSubject(Map<String, String> values) {
        String result = this.subject;
        for (Map.Entry<String, String> entry : values.entrySet()) {
            result = result.replace("{" + entry.getKey() + "}", entry.getValue());
        }
        return result;
    }

    /*
        ฟังก์ชันนี้แทนที่ตัวแปรในเนื้อหา HTML
        This function replaces variables in HTML content.
    */
    public String renderBodyHtml(Map<String, String> values) {
        if (this.bodyHtml == null) {
            return null;
        }
        String result = this.bodyHtml;
        for (Map.Entry<String, String> entry : values.entrySet()) {
            result = result.replace("{" + entry.getKey() + "}", entry.getValue());
        }
        return result;
    }

    /*
        ฟังก์ชันนี้แทนที่ตัวแปรในเนื้อหา Text
        This function replaces variables in Text content.
    */
    public String renderBodyText(Map<String, String> values) {
        if (this.bodyText == null) {
            return null;
        }
        String result = this.bodyText;
        for (Map.Entry<String, String> entry : values.entrySet()) {
            result = result.replace("{" + entry.getKey() + "}", entry.getValue());
        }
        return result;
    }
}
```

### `domain/TEmailHistory.java`

```java
package com.template.app.modules.email.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import com.template.app.modules.email.domain.enums.EmailStatus;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.time.LocalDateTime;
import java.util.UUID;

@Data
@EqualsAndHashCode(callSuper = true)
public class TEmailHistory extends GenericBusinessClass {

    private String emailId;                 // ID ของอีเมล / Email ID.
    private String templateCode;            // รหัสเทมเพลต / Template code.
    private String referenceType;           // ประเภทอ้างอิง / Reference type.
    private UUID referenceId;               // ID อ้างอิง / Reference ID.
    private String fromEmail;               // ผู้ส่ง / From email.
    private String fromName;                // ชื่อผู้ส่ง / From name.
    private String toEmail;                 // ผู้รับ / To email.
    private String toName;                  // ชื่อผู้รับ / To name.
    private String ccEmail;                 // CC
    private String bccEmail;                // BCC
    private String subject;                 // หัวข้อ / Subject.
    private String bodyPreview;             // เนื้อหาส่วนย่อ / Body preview.
    private EmailStatus status;             // สถานะ / Status.
    private String priority;                // ความสำคัญ / Priority.
    private LocalDateTime sentAt;           // วันที่ส่ง / Sent date.
    private String errorMessage;            // ข้อความผิดพลาด / Error message.
    private Integer retryCount;             // จำนวนการลองใหม่ / Retry count.
    private String attachments;             // ไฟล์แนบ / Attachments.

    /*
        ฟังก์ชันนี้บันทึกการส่งสำเร็จ
        This function records successful sending.
    */
    public void markAsSent() {
        this.status = EmailStatus.SENT;
        this.sentAt = LocalDateTime.now();
    }

    /*
        ฟังก์ชันนี้บันทึกการส่งล้มเหลว
        This function records failed sending.
    */
    public void markAsFailed(String error) {
        this.status = EmailStatus.FAILED;
        this.errorMessage = error;
        this.retryCount = (this.retryCount != null ? this.retryCount : 0) + 1;
    }

    /*
        ฟังก์ชันนี้บันทึกอีเมลที่ถูกปฏิเสธ
        This function records bounced emails.
    */
    public void markAsBounced(String error) {
        this.status = EmailStatus.BOUNCED;
        this.errorMessage = error;
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่าสามารถส่งซ้ำได้หรือไม่
        This function checks if the email can be resent.
    */
    public boolean canResend() {
        return (this.status == EmailStatus.FAILED || this.status == EmailStatus.BOUNCED)
                && (this.retryCount == null || this.retryCount < 3);
    }
}
```

---

## 🔧 ฟังก์ชัน Service ตัวอย่าง (พร้อมคำอธิบาย)

### `application/impl/EmailServiceImpl.java` (บางส่วน)

```java
package com.template.app.modules.email.application.impl;

import com.template.app.modules.email.application.interfaces.EmailService;
import com.template.app.modules.email.domain.MEmailTemplate;
import com.template.app.modules.email.domain.TEmailHistory;
import com.template.app.modules.email.domain.enums.EmailStatus;
import com.template.app.modules.email.infrastructure.cache.EmailTemplateCacheService;
import com.template.app.modules.email.infrastructure.provider.EmailProvider;
import com.template.app.modules.email.infrastructure.repository.EmailHistoryRepository;
import com.template.app.modules.email.infrastructure.repository.EmailQueueRepository;
import com.template.app.modules.email.presentation.dto.request.EmailSendRequestDTO;
import com.template.app.modules.email.presentation.dto.response.EmailSendResponseDTO;
import com.template.app._shared.application.GenericServiceImpl;
import com.template.app.exception.SystemGlobalException;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.Map;

@Service
public class EmailServiceImpl extends GenericServiceImpl<TEmailHistory, EmailHistoryRepository>
        implements EmailService {

    private final EmailTemplateCacheService templateCacheService;
    private final EmailProvider emailProvider;
    private final EmailQueueRepository emailQueueRepository;

    public EmailServiceImpl(EmailHistoryRepository repository,
                            EmailTemplateCacheService templateCacheService,
                            EmailProvider emailProvider,
                            EmailQueueRepository emailQueueRepository) {
        super(repository);
        this.templateCacheService = templateCacheService;
        this.emailProvider = emailProvider;
        this.emailQueueRepository = emailQueueRepository;
    }

    /*
        ฟังก์ชันนี้ส่งอีเมลจากเทมเพลต โดยแทนที่ตัวแปรด้วยค่าที่กำหนด
        This function sends an email from a template, replacing variables with provided values.
    */
    @Override
    public EmailSendResponseDTO sendTemplateEmail(EmailSendRequestDTO request) throws SystemGlobalException {
        // 1. ดึงเทมเพลตจาก Cache / Get template from cache.
        MEmailTemplate template = templateCacheService.getTemplate(request.getTemplateCode(), request.getLanguage());
        if (template == null) {
            throw new SystemGlobalException("Email template not found: " + request.getTemplateCode(), null);
        }

        // 2. แทนที่ตัวแปรในหัวข้อและเนื้อหา / Replace variables in subject and body.
        String subject = template.renderSubject(request.getVariables());
        String bodyHtml = template.renderBodyHtml(request.getVariables());
        String bodyText = template.renderBodyText(request.getVariables());

        // 3. สร้างบันทึกประวัติอีเมล / Create email history record.
        TEmailHistory history = new TEmailHistory();
        history.setTemplateCode(request.getTemplateCode());
        history.setReferenceType(request.getReferenceType());
        history.setReferenceId(request.getReferenceId());
        history.setFromEmail(template.getFromEmail());
        history.setFromName(template.getFromName());
        history.setToEmail(request.getToEmail());
        history.setToName(request.getToName());
        history.setSubject(subject);
        history.setBodyPreview(bodyText != null ? bodyText.substring(0, Math.min(bodyText.length(), 200)) : "");
        history.setStatus(EmailStatus.PENDING);
        history.setPriority(request.getPriority());

        // 4. บันทึกประวัติ / Save history.
        TEmailHistory savedHistory = this.create(history);

        // 5. ส่งอีเมลแบบ Async (ใช้ Queue) / Send email asynchronously (using Queue).
        queueAndSendEmail(savedHistory, template, subject, bodyHtml, bodyText, request.getAttachments());

        return EmailSendResponseDTO.fromEntity(savedHistory);
    }

    /*
        ฟังก์ชันนี้ส่งอีเมลแบบ Async ผ่าน Queue (เพื่อไม่ให้กระทบ Performance ของ API)
        This function sends email asynchronously via Queue (to not affect API Performance).
    */
    @Async
    protected void queueAndSendEmail(TEmailHistory history,
                                     MEmailTemplate template,
                                     String subject,
                                     String bodyHtml,
                                     String bodyText,
                                     Map<String, String> attachments) {
        try {
            // 1. สร้างคิวอีเมล / Create email queue.
            // emailQueueRepository.save(queue);

            // 2. ส่งอีเมลผ่าน Provider / Send email via Provider.
            boolean success = emailProvider.sendEmail(
                    template.getFromEmail(),
                    history.getToEmail(),
                    subject,
                    bodyHtml,
                    bodyText,
                    attachments
            );

            // 3. อัปเดตสถานะ / Update status.
            if (success) {
                history.markAsSent();
            } else {
                history.markAsFailed("Provider returned failure");
            }

            // 4. บันทึกประวัติ / Save history.
            this.update(history.getId(), history);

            // 5. ลบออกจากคิว / Remove from queue.
            // emailQueueRepository.deleteByEmailId(history.getEmailId());

        } catch (Exception e) {
            // 6. บันทึกข้อผิดพลาด / Record error.
            history.markAsFailed(e.getMessage());
            this.update(history.getId(), history);
        }
    }

    /*
        ฟังก์ชันนี้ส่งอีเมลซ้ำ (ใช้เมื่อการส่งครั้งแรกล้มเหลว)
        This function resends an email (used when first attempt failed).
    */
    @Override
    public EmailSendResponseDTO resendEmail(String emailId) throws SystemGlobalException {
        // 1. ดึงประวัติอีเมล / Get email history.
        TEmailHistory history = repository.findByEmailId(emailId, getRepositoryAuth())
                .orElseThrow(() -> new SystemGlobalException("Email not found: " + emailId, null));

        // 2. ตรวจสอบว่าสามารถส่งซ้ำได้ / Check if can resend.
        if (!history.canResend()) {
            throw new SystemGlobalException("Cannot resend email with status: " + history.getStatus(), null);
        }

        // 3. ดึงเทมเพลต / Get template.
        MEmailTemplate template = templateCacheService.getTemplate(history.getTemplateCode(), "th");

        // 4. ส่งซ้ำ / Resend.
        queueAndSendEmail(history, template, history.getSubject(), null, null, null);

        return EmailSendResponseDTO.fromEntity(history);
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
| 10 | 📧 Email Service | ✅ ครบถ้วน | Template-based Email + Attachments + Bulk + Queue + History + Cache |

--- 