**โมดูลที่ 12: 🌏 Multi-Language (i18n) - การรองรับหลายภาษา**

โมดูล Multi-Language (i18n) เป็นระบบที่รองรับการใช้งานหลายภาษาสำหรับผู้ใช้ในประเทศต่างๆ ครอบคลุมการทำงานดังนี้:

1. **การจัดการภาษา (Language Management)** – รองรับ 18 ภาษา
2. **การสลับภาษา (Language Switching)** – ผู้ใช้สามารถเปลี่ยนภาษาได้
3. **เทมเพลตข้อความ (Message Templates)** – ข้อความในระบบรองรับหลายภาษา
4. **การแปลเนื้อหา (Content Translation)** – แปลเนื้อหาในเอกสาร PDF, อีเมล
5. **การจัดการ Resource Bundle** – ใช้ Spring i18n และไฟล์ properties
6. **การแปลข้อมูล Master Data** – ข้อมูลหลัก (ประเภท, หมวดหมู่) รองรับหลายภาษา

---

## 📁 โครงสร้างโมดูล Multi-Language (`modules/i18n`)

```
modules/i18n/
├── application/
│   ├── interfaces/
│   │   ├── LanguageService.java
│   │   ├── MessageService.java
│   │   ├── TranslationService.java
│   │   └── LocaleResolverService.java
│   ├── impl/
│   │   ├── LanguageServiceImpl.java
│   │   ├── MessageServiceImpl.java
│   │   ├── TranslationServiceImpl.java
│   │   └── LocaleResolverServiceImpl.java
│   └── usecase/
│       ├── GetMessageUseCase.java
│       ├── GetLanguagesUseCase.java
│       ├── SwitchLanguageUseCase.java
│       ├── TranslateTextUseCase.java
│       ├── GetResourceBundleUseCase.java
│       └── SaveTranslationUseCase.java
├── domain/
│   ├── MLanguage.java
│   ├── MTranslation.java
│   ├── enums/
│   │   └── LanguageCode.java            // th, en, zh, ja, ko, es, fr, de, etc.
│   └── valueobjects/
│       ├── LocaleValue.java
│       └── MessageKey.java
├── infrastructure/
│   ├── repository/
│   │   ├── LanguageRepository.java
│   │   ├── TranslationRepository.java
│   │   └── impl/
│   │       ├── LanguageRepositoryImpl.java
│   │       └── TranslationRepositoryImpl.java
│   ├── cache/                                           // ⬅️ ระบบ Cache สำหรับ i18n
│   │   ├── MessageCacheService.java
│   │   ├── LanguageCacheService.java
│   │   └── TranslationCacheService.java
│   ├── resource/                                        // ⬅️ Resource Bundle
│   │   ├── MessageSourceConfig.java
│   │   ├── CustomMessageSource.java
│   │   └── DatabaseMessageSource.java
│   ├── resolver/                                        // ⬅️ Locale Resolver
│   │   ├── CustomLocaleResolver.java
│   │   └── AcceptHeaderLocaleResolver.java
│   ├── entity/
│   │   ├── LanguageEntity.java
│   │   └── TranslationEntity.java
│   └── mapper/
│       ├── LanguageMapper.java
│       └── TranslationMapper.java
└── presentation/
    ├── controller/
    │   ├── LanguageController.java       // Language Management APIs
    │   └── TranslationController.java    // Translation Management APIs
    ├── dto/
    │   ├── request/
    │   │   ├── LanguageSwitchRequestDTO.java
    │   │   └── TranslationUpdateRequestDTO.java
    │   └── response/
    │       ├── LanguageResponseDTO.java
    │       ├── MessageResponseDTO.java
    │       └── TranslationResponseDTO.java
    └── validator/
        └── TranslationValidator.java
```

---

## 🗄️ ออกแบบฐานข้อมูล (Database Design) สำหรับโมดูล Multi-Language

### 📄 SQL DDL (ไฟล์: `src/main/resources/db/migration/V12__i18n_schema.sql`)

```sql
-- ==============================================
-- ตาราง: m_language (ข้อมูลภาษา)
-- Language master data.
-- ==============================================
CREATE TABLE IF NOT EXISTS m_language (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    language_code VARCHAR(10) UNIQUE NOT NULL,      -- th, en, zh, ja, ko, etc.
    language_name VARCHAR(50) NOT NULL,             -- ภาษาไทย, English, 中文, 日本語, 한국어
    language_name_en VARCHAR(50),                   -- Thai, English, Chinese, Japanese, Korean
    flag_emoji VARCHAR(10),                         -- 🇹🇭, 🇬🇧, 🇨🇳, 🇯🇵, 🇰🇷
    is_rtl BOOLEAN DEFAULT FALSE,                   -- Right-to-Left (อาหรับ, ฮีบรู)
    is_active BOOLEAN DEFAULT TRUE,
    is_default BOOLEAN DEFAULT FALSE,
    sort_order INTEGER DEFAULT 0,
    locale VARCHAR(20),                             -- th_TH, en_US, zh_CN, ja_JP, ko_KR
    date_format VARCHAR(50),                        -- รูปแบบวันที่
    time_format VARCHAR(50),                        -- รูปแบบเวลา
    number_format VARCHAR(50),                      -- รูปแบบตัวเลข
    currency_symbol VARCHAR(10),                    -- สัญลักษณ์สกุลเงิน
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_m_language_code ON m_language(language_code);
CREATE INDEX idx_m_language_active ON m_language(is_active);

-- ==============================================
-- ข้อมูลเริ่มต้น: 18 ภาษา
-- Initial data: 18 languages.
-- ==============================================
INSERT INTO m_language (language_code, language_name, language_name_en, flag_emoji, is_rtl, is_default, sort_order, locale, date_format, currency_symbol, user_id, whitelabel_id)
VALUES 
('th', 'ภาษาไทย', 'Thai', '🇹🇭', false, true, 1, 'th_TH', 'dd/MM/yyyy', '฿', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('en', 'English', 'English', '🇬🇧', false, false, 2, 'en_US', 'MM/dd/yyyy', '$', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('zh', '中文', 'Chinese', '🇨🇳', false, false, 3, 'zh_CN', 'yyyy/MM/dd', '¥', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('ja', '日本語', 'Japanese', '🇯🇵', false, false, 4, 'ja_JP', 'yyyy/MM/dd', '¥', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('ko', '한국어', 'Korean', '🇰🇷', false, false, 5, 'ko_KR', 'yyyy/MM/dd', '₩', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('es', 'Español', 'Spanish', '🇪🇸', false, false, 6, 'es_ES', 'dd/MM/yyyy', '€', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('fr', 'Français', 'French', '🇫🇷', false, false, 7, 'fr_FR', 'dd/MM/yyyy', '€', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('de', 'Deutsch', 'German', '🇩🇪', false, false, 8, 'de_DE', 'dd.MM.yyyy', '€', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('it', 'Italiano', 'Italian', '🇮🇹', false, false, 9, 'it_IT', 'dd/MM/yyyy', '€', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('pt', 'Português', 'Portuguese', '🇵🇹', false, false, 10, 'pt_PT', 'dd/MM/yyyy', '€', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('ru', 'Русский', 'Russian', '🇷🇺', false, false, 11, 'ru_RU', 'dd.MM.yyyy', '₽', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('ar', 'العربية', 'Arabic', '🇸🇦', true, false, 12, 'ar_SA', 'dd/MM/yyyy', '﷼', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('hi', 'हिन्दी', 'Hindi', '🇮🇳', false, false, 13, 'hi_IN', 'dd/MM/yyyy', '₹', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('id', 'Bahasa Indonesia', 'Indonesian', '🇮🇩', false, false, 14, 'id_ID', 'dd/MM/yyyy', 'Rp', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('ms', 'Bahasa Melayu', 'Malay', '🇲🇾', false, false, 15, 'ms_MY', 'dd/MM/yyyy', 'RM', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('vi', 'Tiếng Việt', 'Vietnamese', '🇻🇳', false, false, 16, 'vi_VN', 'dd/MM/yyyy', '₫', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('my', 'မြန်မာဘာသာ', 'Burmese', '🇲🇲', false, false, 17, 'my_MM', 'dd/MM/yyyy', 'Ks', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('lo', 'ພາສາລາວ', 'Lao', '🇱🇦', false, false, 18, 'lo_LA', 'dd/MM/yyyy', '₭', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001');

-- ==============================================
-- ตาราง: m_translation (ข้อความที่แปลแล้ว)
-- Translation data for messages.
-- ==============================================
CREATE TABLE IF NOT EXISTS m_translation (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_key VARCHAR(255) NOT NULL,              -- รหัสข้อความ (เช่น job.status.open)
    language_code VARCHAR(10) NOT NULL REFERENCES m_language(language_code) ON DELETE CASCADE,
    message_text TEXT NOT NULL,                     -- ข้อความที่แปลแล้ว
    context VARCHAR(100),                           -- บริบท (เช่น UI, EMAIL, PDF)
    description TEXT,                               -- คำอธิบายเพิ่มเติม
    version INTEGER DEFAULT 1,
    is_approved BOOLEAN DEFAULT TRUE,
    approved_by UUID REFERENCES m_user(id),
    approved_at TIMESTAMP,
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL,
    -- UNIQUE constraint
    UNIQUE (message_key, language_code, context)
);

CREATE INDEX idx_m_translation_key ON m_translation(message_key);
CREATE INDEX idx_m_translation_language ON m_translation(language_code);
CREATE INDEX idx_m_translation_context ON m_translation(context);
CREATE INDEX idx_m_translation_whitelabel ON m_translation(whitelabel_id);

-- ==============================================
-- ข้อมูลเริ่มต้น: ข้อความพื้นฐาน
-- Initial data: Basic messages.
-- ==============================================
INSERT INTO m_translation (message_key, language_code, message_text, context, user_id, whitelabel_id)
VALUES 
('job.status.open', 'th', 'เปิดใบงาน', 'UI', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('job.status.open', 'en', 'Open', 'UI', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('job.status.in_progress', 'th', 'กำลังดำเนินการ', 'UI', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('job.status.in_progress', 'en', 'In Progress', 'UI', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('quotation.approved', 'th', 'อนุมัติใบเสนอราคาแล้ว', 'UI', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('quotation.approved', 'en', 'Quotation Approved', 'UI', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001');

-- ==============================================
-- VIEW: v_available_languages (ภาษาที่ใช้งานได้)
-- Available languages view.
-- ==============================================
CREATE OR REPLACE VIEW v_available_languages AS
SELECT 
    language_code,
    language_name,
    language_name_en,
    flag_emoji,
    is_rtl,
    is_default,
    locale,
    date_format,
    currency_symbol
FROM m_language
WHERE is_active = true;
```

---

## 🧠 ระบบ Cache (Redis) สำหรับ Multi-Language

### `infrastructure/cache/MessageCacheService.java`

```java
package com.template.app.modules.i18n.infrastructure.cache;

import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

@Service
public class MessageCacheService {

    /*
        ฟังก์ชันนี้ดึงข้อความที่แปลแล้วจาก Cache ตาม message_key และภาษา
        This function retrieves translated messages from cache by message_key and language.
        Redis Key: message:{messageKey}:{languageCode}
    */
    @Cacheable(value = "messages", key = "#messageKey + ':' + #languageCode")
    public String getMessage(String messageKey, String languageCode) {
        return null;
    }

    /*
        ฟังก์ชันนี้ดึงข้อความที่แปลแล้วทั้งหมดของภาษาที่ระบุ (ใช้สำหรับโหลดครั้งแรก)
        This function retrieves all translated messages for a language (used for initial load).
        Redis Key: messages_all:{languageCode}
    */
    @Cacheable(value = "messages_all", key = "#languageCode")
    public java.util.Map<String, String> getAllMessages(String languageCode) {
        return null;
    }

    /*
        ฟังก์ชันนี้ลบข้อความที่แปลแล้วออกจาก Cache เมื่อมีการอัปเดต
        This function evicts translated messages from cache when updated.
    */
    @CacheEvict(value = {"messages", "messages_all"}, key = "#messageKey + ':' + #languageCode")
    public void evictMessage(String messageKey, String languageCode) {
        // ลบ Cache / Evict cache.
    }

    /*
        ฟังก์ชันนี้ลบ Cache ทั้งหมดของข้อความ
        This function clears all message caches.
    */
    @CacheEvict(value = {"messages", "messages_all"}, allEntries = true)
    public void evictAllMessages() {
        // ลบทุก key / Evict all keys.
    }
}
```

### `infrastructure/cache/LanguageCacheService.java`

```java
package com.template.app.modules.i18n.infrastructure.cache;

import com.template.app.modules.i18n.domain.MLanguage;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.CachePut;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class LanguageCacheService {

    /*
        ฟังก์ชันนี้ดึงข้อมูลภาษาจาก Cache
        This function retrieves language data from cache.
        Redis Key: language:{languageCode}
    */
    @Cacheable(value = "languages", key = "#languageCode")
    public MLanguage getLanguage(String languageCode) {
        return null;
    }

    /*
        ฟังก์ชันนี้ดึงรายการภาษาทั้งหมดที่ใช้งานอยู่
        This function retrieves all active languages.
        Redis Key: languages_active
    */
    @Cacheable(value = "languages_active", key = "'all'")
    public List<MLanguage> getActiveLanguages() {
        return null;
    }

    /*
        ฟังก์ชันนี้จะอัปเดต Cache เมื่อมีการบันทึกภาษา
        This function updates the cache when a language is saved.
    */
    @CachePut(value = "languages", key = "#language.languageCode")
    public MLanguage saveLanguage(MLanguage language) {
        return language;
    }

    /*
        ฟังก์ชันนี้ลบข้อมูลภาษาออกจาก Cache
        This function evicts language data from cache.
    */
    @CacheEvict(value = {"languages", "languages_active"}, key = "#languageCode")
    public void evictLanguage(String languageCode) {
        // ลบ Cache / Evict cache.
    }

    /*
        ฟังก์ชันนี้ลบ Cache ทั้งหมดของภาษา
        This function clears all language caches.
    */
    @CacheEvict(value = {"languages", "languages_active"}, allEntries = true)
    public void evictAllLanguages() {
        // ลบทุก key / Evict all keys.
    }
}
```

---

## 🔧 Spring Configuration สำหรับ i18n

### `infrastructure/resource/MessageSourceConfig.java`

```java
package com.template.app.modules.i18n.infrastructure.resource;

import org.springframework.context.MessageSource;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.support.ReloadableResourceBundleMessageSource;
import org.springframework.web.servlet.LocaleResolver;
import org.springframework.web.servlet.i18n.CookieLocaleResolver;
import org.springframework.web.servlet.i18n.LocaleChangeInterceptor;

import java.util.Locale;

@Configuration
public class MessageSourceConfig {

    /*
        ฟังก์ชันนี้กำหนดค่า MessageSource สำหรับการโหลดไฟล์ properties (รองรับหลายภาษา)
        This function configures MessageSource for loading properties files (multi-language support).
        ไฟล์ properties อยู่ใน: src/main/resources/i18n/messages_{language}.properties
        Properties files are in: src/main/resources/i18n/messages_{language}.properties
    */
    @Bean
    public MessageSource messageSource() {
        ReloadableResourceBundleMessageSource messageSource = new ReloadableResourceBundleMessageSource();
        messageSource.setBasename("classpath:i18n/messages");
        messageSource.setDefaultEncoding("UTF-8");
        messageSource.setCacheSeconds(3600); // Cache 1 ชั่วโมง / Cache 1 hour.
        messageSource.setFallbackToSystemLocale(false);
        messageSource.setUseCodeAsDefaultMessage(true);
        return messageSource;
    }

    /*
        ฟังก์ชันนี้กำหนดค่า LocaleResolver สำหรับการจัดเก็บภาษาที่ผู้ใช้เลือก (ใช้ Cookie)
        This function configures LocaleResolver for storing user-selected language (uses Cookie).
    */
    @Bean
    public LocaleResolver localeResolver() {
        CookieLocaleResolver localeResolver = new CookieLocaleResolver();
        localeResolver.setCookieName("lang");
        localeResolver.setCookieMaxAge(3600 * 24 * 30); // 30 วัน / 30 days.
        localeResolver.setDefaultLocale(Locale.forLanguageTag("th"));
        return localeResolver;
    }

    /*
        ฟังก์ชันนี้กำหนดค่า LocaleChangeInterceptor สำหรับการสลับภาษาผ่าน Query Parameter ?lang=th
        This function configures LocaleChangeInterceptor for language switching via Query Parameter ?lang=th.
    */
    @Bean
    public LocaleChangeInterceptor localeChangeInterceptor() {
        LocaleChangeInterceptor interceptor = new LocaleChangeInterceptor();
        interceptor.setParamName("lang");
        return interceptor;
    }
}
```

---

## ⏱️ API Controller สำหรับ Multi-Language

```java
package com.template.app.modules.i18n.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.i18n.application.interfaces.LanguageService;
import com.template.app.modules.i18n.application.interfaces.MessageService;
import com.template.app.modules.i18n.presentation.dto.request.LanguageSwitchRequestDTO;
import com.template.app.modules.i18n.presentation.dto.response.LanguageResponseDTO;
import com.template.app.modules.i18n.presentation.dto.response.MessageResponseDTO;
import com.template.app.exception.SystemGlobalException;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Locale;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/languages")
@Tag(name = "Multi-Language (i18n)", description = "Internationalization and Language Management APIs")
@RequiredArgsConstructor
public class LanguageController {

    private final LanguageService languageService;
    private final MessageService messageService;

    // ========================================================================
    // 1. GET ALL LANGUAGES (รายการภาษาทั้งหมด)
    // ========================================================================

    /*
        API: GET /api/v1/languages
        ฟังก์ชันนี้แสดงรายการภาษาที่รองรับทั้งหมด (ใช้ Cache)
        This function lists all supported languages (uses caching).
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @GetMapping
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get all supported languages")
    public ResponseEntity<List<LanguageResponseDTO>> getLanguages() throws SystemGlobalException {
        List<LanguageResponseDTO> languages = languageService.getAllLanguages();
        return ResponseEntity.ok(languages);
    }

    // ========================================================================
    // 2. GET CURRENT LANGUAGE (ภาษาปัจจุบัน)
    // ========================================================================

    /*
        API: GET /api/v1/languages/current
        ฟังก์ชันนี้ดึงภาษาที่ผู้ใช้กำลังใช้งานอยู่
        This function retrieves the user's current language.
        Rate Limit: อนุญาต 100 ครั้งต่อ 1 นาที
        Rate Limit: Allows 100 requests per minute.
    */
    @GetMapping("/current")
    @RateLimit(limit = 100, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get current language")
    public ResponseEntity<LanguageResponseDTO> getCurrentLanguage() throws SystemGlobalException {
        LanguageResponseDTO response = languageService.getCurrentLanguage();
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. SWITCH LANGUAGE (สลับภาษา)
    // ========================================================================

    /*
        API: POST /api/v1/languages/switch
        ฟังก์ชันนี้สลับภาษาให้กับผู้ใช้ปัจจุบัน (จะถูกบันทึกใน Cookie)
        This function switches the language for the current user (stored in Cookie).
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @PostMapping("/switch")
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Switch language")
    public ResponseEntity<LanguageResponseDTO> switchLanguage(
            @Valid @RequestBody LanguageSwitchRequestDTO request,
            Locale locale) throws SystemGlobalException {
        LanguageResponseDTO response = languageService.switchLanguage(request.getLanguageCode());
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 4. GET MESSAGES (ข้อความทั้งหมดของภาษา)
    // ========================================================================

    /*
        API: GET /api/v1/languages/messages/{languageCode}
        ฟังก์ชันนี้ดึงข้อความทั้งหมดของภาษาที่ระบุ (ใช้สำหรับ Frontend i18n)
        This function retrieves all messages for a language (used for Frontend i18n).
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที
        Rate Limit: Allows 30 requests per minute.
    */
    @GetMapping("/messages/{languageCode}")
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get all messages for language")
    public ResponseEntity<Map<String, String>> getMessages(@PathVariable String languageCode)
            throws SystemGlobalException {
        Map<String, String> messages = messageService.getAllMessages(languageCode);
        return ResponseEntity.ok(messages);
    }

    // ========================================================================
    // 5. GET SINGLE MESSAGE (ข้อความเดียว)
    // ========================================================================

    /*
        API: GET /api/v1/languages/message
        ฟังก์ชันนี้ดึงข้อความเดียวตาม message_key และภาษา (ใช้ Cache)
        This function retrieves a single message by message_key and language (uses caching).
        Rate Limit: อนุญาต 100 ครั้งต่อ 1 นาที
        Rate Limit: Allows 100 requests per minute.
    */
    @GetMapping("/message")
    @RateLimit(limit = 100, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get single message by key")
    public ResponseEntity<MessageResponseDTO> getMessage(
            @RequestParam String messageKey,
            @RequestParam(required = false) String languageCode) throws SystemGlobalException {
        MessageResponseDTO response = messageService.getMessage(messageKey, languageCode);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 6. GET DEFAULT LANGUAGE (ภาษาเริ่มต้น)
    // ========================================================================

    /*
        API: GET /api/v1/languages/default
        ฟังก์ชันนี้ดึงภาษาเริ่มต้นของระบบ
        This function retrieves the system's default language.
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @GetMapping("/default")
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get default language")
    public ResponseEntity<LanguageResponseDTO> getDefaultLanguage() throws SystemGlobalException {
        LanguageResponseDTO response = languageService.getDefaultLanguage();
        return ResponseEntity.ok(response);
    }
}
```

### `TranslationController.java` (จัดการข้อความที่แปลแล้ว)

```java
package com.template.app.modules.i18n.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.i18n.application.interfaces.TranslationService;
import com.template.app.modules.i18n.presentation.dto.request.TranslationUpdateRequestDTO;
import com.template.app.modules.i18n.presentation.dto.response.TranslationResponseDTO;
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
@RequestMapping("/api/v1/translations")
@Tag(name = "Translations", description = "Translation Management APIs")
@RequiredArgsConstructor
public class TranslationController {

    private final TranslationService translationService;

    // ========================================================================
    // 1. GET TRANSLATIONS (ข้อความที่แปลแล้วทั้งหมด)
    // ========================================================================

    /*
        API: GET /api/v1/translations
        ฟังก์ชันนี้แสดงรายการข้อความที่แปลแล้วทั้งหมด แบบแบ่งหน้า
        This function lists all translated messages with pagination.
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที
        Rate Limit: Allows 30 requests per minute.
    */
    @GetMapping
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "List all translations")
    public ResponseEntity<Page<TranslationResponseDTO>> listTranslations(
            @RequestParam(required = false) String languageCode,
            @RequestParam(required = false) String messageKey,
            Pageable pageable) throws SystemGlobalException {
        Page<TranslationResponseDTO> page = translationService.listTranslations(languageCode, messageKey, pageable);
        return ResponseEntity.ok(page);
    }

    // ========================================================================
    // 2. UPDATE TRANSLATION (อัปเดตข้อความที่แปลแล้ว)
    // ========================================================================

    /*
        API: PUT /api/v1/translations/{id}
        ฟังก์ชันนี้แก้ไขข้อความที่แปลแล้ว (ต้องมีสิทธิ์ Admin เท่านั้น)
        This function updates a translated message (Admin only).
        Rate Limit: อนุญาต 10 ครั้งต่อ 5 นาที
        Rate Limit: Allows 10 requests per 5 minutes.
    */
    @PutMapping("/{id}")
    @RateLimit(limit = 10, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Update translation")
    public ResponseEntity<TranslationResponseDTO> updateTranslation(
            @PathVariable UUID id,
            @Valid @RequestBody TranslationUpdateRequestDTO request) throws SystemGlobalException {
        TranslationResponseDTO response = translationService.updateTranslation(id, request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. CREATE TRANSLATION (สร้างข้อความที่แปลแล้วใหม่)
    // ========================================================================

    /*
        API: POST /api/v1/translations
        ฟังก์ชันนี้สร้างข้อความที่แปลแล้วใหม่ (ต้องมีสิทธิ์ Admin เท่านั้น)
        This function creates a new translated message (Admin only).
        Rate Limit: อนุญาต 10 ครั้งต่อ 5 นาที
        Rate Limit: Allows 10 requests per 5 minutes.
    */
    @PostMapping
    @RateLimit(limit = 10, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Create translation")
    public ResponseEntity<TranslationResponseDTO> createTranslation(
            @Valid @RequestBody TranslationUpdateRequestDTO request) throws SystemGlobalException {
        TranslationResponseDTO response = translationService.createTranslation(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 4. DELETE TRANSLATION (ลบข้อความที่แปลแล้ว)
    // ========================================================================

    /*
        API: DELETE /api/v1/translations/{id}
        ฟังก์ชันนี้ลบข้อความที่แปลแล้ว (ต้องมีสิทธิ์ Admin เท่านั้น)
        This function deletes a translated message (Admin only).
        Rate Limit: อนุญาต 5 ครั้งต่อ 1 ชั่วโมง
        Rate Limit: Allows 5 requests per 1 hour.
    */
    @DeleteMapping("/{id}")
    @RateLimit(limit = 5, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Delete translation")
    public ResponseEntity<Void> deleteTranslation(@PathVariable UUID id) throws SystemGlobalException {
        translationService.deleteTranslation(id);
        return ResponseEntity.noContent().build();
    }

    // ========================================================================
    // 5. BULK IMPORT TRANSLATIONS (นำเข้าข้อความจำนวนมาก)
    // ========================================================================

    /*
        API: POST /api/v1/translations/bulk-import
        ฟังก์ชันนี้นำเข้าข้อความที่แปลแล้วจำนวนมาก (ใช้ JSON หรือ CSV)
        This function bulk imports translated messages (using JSON or CSV).
        Rate Limit: อนุญาต 3 ครั้งต่อ 1 ชั่วโมง
        Rate Limit: Allows 3 requests per 1 hour.
    */
    @PostMapping("/bulk-import")
    @RateLimit(limit = 3, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Bulk import translations")
    public ResponseEntity<BulkImportResponseDTO> bulkImportTranslations(
            @RequestBody BulkImportRequestDTO request) throws SystemGlobalException {
        BulkImportResponseDTO response = translationService.bulkImportTranslations(request);
        return ResponseEntity.ok(response);
    }
}
```

---

## 🧩 ตัวอย่าง Domain Entity และ Enum

### `domain/enums/LanguageCode.java`

```java
package com.template.app.modules.i18n.domain.enums;

/*
    รหัสภาษาที่รองรับ / Supported language codes.
*/
public enum LanguageCode {
    TH("th", "ภาษาไทย", "Thai"),
    EN("en", "English", "English"),
    ZH("zh", "中文", "Chinese"),
    JA("ja", "日本語", "Japanese"),
    KO("ko", "한국어", "Korean"),
    ES("es", "Español", "Spanish"),
    FR("fr", "Français", "French"),
    DE("de", "Deutsch", "German"),
    IT("it", "Italiano", "Italian"),
    PT("pt", "Português", "Portuguese"),
    RU("ru", "Русский", "Russian"),
    AR("ar", "العربية", "Arabic"),
    HI("hi", "हिन्दी", "Hindi"),
    ID("id", "Bahasa Indonesia", "Indonesian"),
    MS("ms", "Bahasa Melayu", "Malay"),
    VI("vi", "Tiếng Việt", "Vietnamese"),
    MY("my", "မြန်မာဘာသာ", "Burmese"),
    LO("lo", "ພາສາລາວ", "Lao");

    private final String code;
    private final String name;
    private final String nameEn;

    LanguageCode(String code, String name, String nameEn) {
        this.code = code;
        this.name = name;
        this.nameEn = nameEn;
    }

    public String getCode() { return code; }
    public String getName() { return name; }
    public String getNameEn() { return nameEn; }

    public static LanguageCode fromCode(String code) {
        for (LanguageCode lang : values()) {
            if (lang.code.equalsIgnoreCase(code)) {
                return lang;
            }
        }
        return TH; // ค่าเริ่มต้น / Default.
    }
}
```

### `domain/MLanguage.java`

```java
package com.template.app.modules.i18n.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import lombok.Data;
import lombok.EqualsAndHashCode;

@Data
@EqualsAndHashCode(callSuper = true)
public class MLanguage extends GenericBusinessClass {

    private String languageCode;        // รหัสภาษา / Language code.
    private String languageName;        // ชื่อภาษา / Language name.
    private String languageNameEn;      // ชื่อภาษา (อังกฤษ) / Language name (English).
    private String flagEmoji;           // อีโมจิธง / Flag emoji.
    private Boolean isRtl;              // Right-to-Left / RTL.
    private Boolean isActive;           // ใช้งานอยู่ / Active.
    private Boolean isDefault;          // ค่าเริ่มต้น / Default.
    private Integer sortOrder;          // ลำดับการแสดง / Display order.
    private String locale;              // th_TH, en_US, etc.
    private String dateFormat;          // รูปแบบวันที่ / Date format.
    private String timeFormat;          // รูปแบบเวลา / Time format.
    private String numberFormat;        // รูปแบบตัวเลข / Number format.
    private String currencySymbol;      // สัญลักษณ์สกุลเงิน / Currency symbol.

    /*
        ฟังก์ชันนี้ตรวจสอบว่าเป็นภาษาเริ่มต้นหรือไม่
        This function checks if this is the default language.
    */
    public boolean isDefault() {
        return this.isDefault != null && this.isDefault;
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่าภาษานี้สามารถใช้งานได้หรือไม่
        This function checks if the language is active.
    */
    public boolean isUsable() {
        return this.isActive != null && this.isActive;
    }

    /*
        ฟังก์ชันนี้แปลง Locale String เป็น Java Locale object
        This function converts Locale String to Java Locale object.
    */
    public java.util.Locale toLocale() {
        if (this.locale == null) {
            return java.util.Locale.forLanguageTag(this.languageCode);
        }
        return java.util.Locale.forLanguageTag(this.locale.replace("_", "-"));
    }
}
```

### `domain/MTranslation.java`

```java
package com.template.app.modules.i18n.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.time.LocalDateTime;
import java.util.UUID;

@Data
@EqualsAndHashCode(callSuper = true)
public class MTranslation extends GenericBusinessClass {

    private String messageKey;          // รหัสข้อความ / Message key.
    private String languageCode;        // รหัสภาษา / Language code.
    private String messageText;         // ข้อความที่แปลแล้ว / Translated message.
    private String context;             // บริบท / Context.
    private String description;         // คำอธิบาย / Description.
    private Integer version;            // เวอร์ชัน / Version.
    private Boolean isApproved;         // อนุมัติแล้ว / Approved.
    private UUID approvedBy;            // ผู้อนุมัติ / Approved by.
    private LocalDateTime approvedAt;   // วันที่อนุมัติ / Approved date.

    /*
        ฟังก์ชันนี้บันทึกการอนุมัติข้อความ
        This function records message approval.
    */
    public void approve(UUID approverId) {
        this.isApproved = true;
        this.approvedBy = approverId;
        this.approvedAt = LocalDateTime.now();
        this.version = (this.version != null ? this.version : 0) + 1;
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่าข้อความได้รับการอนุมัติแล้วหรือไม่
        This function checks if the message is approved.
    */
    public boolean isApproved() {
        return this.isApproved != null && this.isApproved;
    }
}
```

---

## 🔧 ฟังก์ชัน Service ตัวอย่าง (พร้อมคำอธิบาย)

### `application/impl/MessageServiceImpl.java` (บางส่วน)

```java
package com.template.app.modules.i18n.application.impl;

import com.template.app.modules.i18n.application.interfaces.MessageService;
import com.template.app.modules.i18n.infrastructure.cache.MessageCacheService;
import com.template.app.modules.i18n.infrastructure.repository.TranslationRepository;
import com.template.app.modules.i18n.presentation.dto.response.MessageResponseDTO;
import com.template.app._shared.application.GenericServiceImpl;
import com.template.app.exception.SystemGlobalException;
import org.springframework.context.MessageSource;
import org.springframework.stereotype.Service;

import java.util.HashMap;
import java.util.List;
import java.util.Locale;
import java.util.Map;

@Service
public class MessageServiceImpl implements MessageService {

    private final MessageSource messageSource;
    private final TranslationRepository translationRepository;
    private final MessageCacheService cacheService;

    public MessageServiceImpl(MessageSource messageSource,
                              TranslationRepository translationRepository,
                              MessageCacheService cacheService) {
        this.messageSource = messageSource;
        this.translationRepository = translationRepository;
        this.cacheService = cacheService;
    }

    /*
        ฟังก์ชันนี้ดึงข้อความที่แปลแล้วจาก Cache หรือ Database (ใช้ Spring MessageSource)
        This function retrieves translated messages from Cache or Database (using Spring MessageSource).
    */
    @Override
    public MessageResponseDTO getMessage(String messageKey, String languageCode) throws SystemGlobalException {
        // 1. ตรวจสอบภาษา / Validate language.
        String lang = languageCode != null ? languageCode : "th";

        // 2. ลองดึงจาก Cache / Try to get from cache.
        String cachedMessage = cacheService.getMessage(messageKey, lang);
        if (cachedMessage != null) {
            return MessageResponseDTO.builder()
                    .messageKey(messageKey)
                    .messageText(cachedMessage)
                    .languageCode(lang)
                    .fromCache(true)
                    .build();
        }

        // 3. ถ้าไม่มีใน Cache ให้ใช้ Spring MessageSource / If not in cache, use Spring MessageSource.
        Locale locale = Locale.forLanguageTag(lang);
        String message = messageSource.getMessage(messageKey, null, messageKey, locale);

        // 4. ถ้า message เท่ากับ messageKey แสดงว่าไม่พบใน Resource Bundle
        //    ให้ค้นหาใน Database แทน / If message equals messageKey, it's not in Resource Bundle.
        if (message.equals(messageKey)) {
            // ค้นหาใน Database / Search in Database.
            message = translationRepository.findByMessageKeyAndLanguageCode(messageKey, lang)
                    .map(translation -> translation.getMessageText())
                    .orElse(messageKey);
        }

        // 5. บันทึกใน Cache / Store in cache.
        // cacheService.saveMessage(messageKey, lang, message);

        return MessageResponseDTO.builder()
                .messageKey(messageKey)
                .messageText(message)
                .languageCode(lang)
                .fromCache(false)
                .build();
    }

    /*
        ฟังก์ชันนี้ดึงข้อความทั้งหมดของภาษาที่ระบุ (ใช้สำหรับ Frontend)
        This function retrieves all messages for a language (used for Frontend).
    */
    @Override
    public Map<String, String> getAllMessages(String languageCode) throws SystemGlobalException {
        // 1. ลองดึงจาก Cache / Try to get from cache.
        Map<String, String> cachedMessages = cacheService.getAllMessages(languageCode);
        if (cachedMessages != null) {
            return cachedMessages;
        }

        // 2. ถ้าไม่มีใน Cache ให้ดึงจาก Database / If not in cache, fetch from Database.
        List<MTranslation> translations = translationRepository.findByLanguageCode(languageCode);
        Map<String, String> messages = new HashMap<>();
        for (MTranslation translation : translations) {
            messages.put(translation.getMessageKey(), translation.getMessageText());
        }

        // 3. บันทึกใน Cache / Store in cache.
        // cacheService.saveAllMessages(languageCode, messages);

        return messages;
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
| 11 | ⏱️ Batch Jobs | ✅ ครบถ้วน | 6 Scheduled Jobs + Distributed Lock + History + Manual Trigger + Cache |
| 12 | 🌏 Multi-Language (i18n) | ✅ ครบถ้วน | 18 Languages + Translation Management + Resource Bundle + Cache |

---
 