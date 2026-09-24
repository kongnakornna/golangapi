# 📋 Prompts สำหรับดำเนินการตามเอกสารโมดูล Batch Jobs (ฉบับสมบูรณ์)

ตามเอกสาร `modules_11_Batch.md` ที่คุณได้จัดเตรียมไว้ ข้าพเจ้าจะดำเนินการตามคำสั่งของคุณอย่างครบถ้วน โดยแบ่งเป็น 8 ขั้นตอนตามที่ระบุ พร้อมสร้าง Prompt ที่สามารถนำไปใช้ต่อได้ทันที

---

## 🎯 สรุปคำสั่ง (Prompts) ทั้ง 8 ข้อ

### Prompt 1: ดำเนินการตามเอกสาร ให้สร้าง Folder, File และเขียน Code ตามเอกสาร

```
1. ดำเนินการตามเอกสาร modules_11_Batch.md ให้สมบูรณ์ โดยให้ทำการสร้าง Folder และ File ตามโครงสร้างที่กำหนดไว้ในเอกสาร และเขียน Code ตามตัวอย่างที่ให้มาในแต่ละส่วน ให้ครบทุกไฟล์

2. กำหนดให้มี 7 งาน (batch001 - batch007) โดย batch007 เป็นงานแจ้งเตือนสต็อกต่ำที่เชื่อมต่อกับโมดูล Inventory

3. ทุกไฟล์ต้องมีคำอธิบายภาษาไทยและภาษาอังกฤษ (Bilingual Comments) ตามมาตรฐานของโปรเจกต์

4. ตรวจสอบให้แน่ใจว่า Dependency Injection ถูกต้อง ไม่มี Circular Dependency

5. สร้าง Database Schema ตาม V11__batch_schema.sql และเพิ่มข้อมูลเริ่มต้น (Initial Data) ทั้ง 7 งาน
```

---

### Prompt 2: สร้าง Folder และ File ตามโครงสร้าง และเขียน Code ตามตัวอย่าง

```
1. สร้างโครงสร้าง Folder ตามนี้:
src/main/java/com/icmon/module/batch/
├── application/
│   ├── interfaces/
│   │   ├── BatchJobService.java
│   │   ├── BatchJobScheduler.java
│   │   └── BatchJobHistoryService.java
│   ├── impl/
│   │   ├── BatchJobServiceImpl.java
│   │   ├── BatchJobSchedulerImpl.java
│   │   └── BatchJobHistoryServiceImpl.java
│   └── usecase/
│       ├── ExecuteBatchJobUseCase.java
│       ├── GetBatchJobStatusUseCase.java
│       ├── ListBatchJobsUseCase.java
│       ├── TriggerBatchJobManuallyUseCase.java
│       ├── StopBatchJobUseCase.java
│       ├── GetBatchJobHistoryUseCase.java
│       └── RetryFailedBatchJobUseCase.java
├── domain/
│   ├── MBatchJob.java
│   ├── TBatchJobHistory.java
│   ├── enums/
│   │   ├── BatchJobStatus.java
│   │   ├── BatchJobType.java
│   │   └── BatchJobPriority.java
│   └── valueobjects/
│       ├── CronExpression.java
│       └── JobExecutionTime.java
├── infrastructure/
│   ├── repository/
│   │   ├── BatchJobRepository.java
│   │   ├── BatchJobHistoryRepository.java
│   │   └── impl/
│   │       ├── BatchJobRepositoryImpl.java
│   │       └── BatchJobHistoryRepositoryImpl.java
│   ├── cache/
│   │   ├── BatchJobCacheService.java
│   │   └── BatchJobLockCacheService.java
│   ├── scheduler/
│   │   ├── BatchSchedulerConfig.java
│   │   ├── Job001DailyNotification.java
│   │   ├── Job002DailyReport.java
│   │   ├── Job003UpdatePendingStatus.java
│   │   ├── Job004CleanupAndSync.java
│   │   ├── Job005RealtimeSync.java
│   │   ├── Job006DailySalesSummary.java
│   │   └── Job007LowStockAlert.java
│   ├── lock/
│   │   └── DistributedLockService.java
│   ├── executor/
│   │   ├── BatchJobExecutor.java
│   │   └── ThreadPoolConfig.java
│   ├── entity/
│   │   ├── BatchJobEntity.java
│   │   └── BatchJobHistoryEntity.java
│   └── mapper/
│       ├── BatchJobMapper.java
│       └── BatchJobHistoryMapper.java
└── presentation/
    ├── controller/
    │   ├── BatchJobController.java
    │   └── BatchJobHistoryController.java
    ├── dto/
    │   ├── request/
    │   │   ├── TriggerJobRequestDTO.java
    │   │   └── BatchJobSearchRequestDTO.java
    │   └── response/
    │       ├── BatchJobResponseDTO.java
    │       ├── BatchJobHistoryResponseDTO.java
    │       └── BatchJobStatusResponseDTO.java
    └── validator/
        └── BatchJobValidator.java

2. เขียน Code ให้ครบทุกไฟล์ตามตัวอย่างในเอกสาร พร้อม Bilingual Comments

3. ให้เพิ่ม batch007 (Low Stock Alert) ที่เชื่อมต่อกับ Inventory Module ตามตัวอย่างที่ให้ไว้
```

---

### Prompt 3: สร้าง Unit Test และ Run Test ให้ผ่านทุก Test Case

```
1. สร้าง Unit Test สำหรับโมดูล Batch Jobs ให้ครอบคลุมทุกฟังก์ชันหลัก:

src/test/java/com/icmon/module/batch/
├── domain/
│   └── MBatchJobTest.java
├── application/
│   ├── BatchJobServiceImplTest.java
│   └── BatchJobHistoryServiceImplTest.java
├── infrastructure/
│   ├── BatchJobExecutorTest.java
│   ├── BatchSchedulerConfigTest.java
│   └── BatchJobLockCacheServiceTest.java
└── presentation/
    └── BatchJobControllerTest.java

2. Test Cases ที่ต้องมี:
   - ✅ การสร้างงาน Batch ใหม่
   - ✅ การอัปเดตงาน Batch
   - ✅ การเปิด/ปิดใช้งานงาน
   - ✅ การสั่งรันงานด้วยตนเอง (Manual Trigger)
   - ✅ การหยุดงานที่กำลังทำงาน
   - ✅ การดึงสถานะงาน
   - ✅ การดึงประวัติการรันงาน
   - ✅ การ Retry งานที่ล้มเหลว
   - ✅ การทำงานของ Distributed Lock
   - ✅ การทำงานของ batch007 (Low Stock Alert)

3. ใช้ Mockito สำหรับ Mock Dependencies
4. ใช้ AssertJ สำหรับ Assertions
5. ใช้ @DisplayName สำหรับตั้งชื่อ Test แบบภาษาไทยและอังกฤษ
6. Run Test และปรับปรุงจนผ่านทุก Test Case (Coverage ≥ 80%)
```

---

### Prompt 4: เพิ่ม Demo Data และ Run Test ให้ผ่าน

```
1. เพิ่ม Demo Data สำหรับโมดูล Batch Jobs ลงใน Database:

-- เพิ่มข้อมูลงานทั้ง 7 ใน m_batch_job (ถ้ายังไม่มี)
INSERT INTO m_batch_job (job_code, job_name, job_type, description, cron_expression, enabled, user_id, whitelabel_id)
VALUES 
('batch001', 'ส่งอีเมลแจ้งเตือนรายวัน', 'EMAIL', 'ส่งอีเมลแจ้งเตือนรายวัน', '0 30 6 ? * *', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('batch002', 'สร้างรายงานประจำวัน', 'REPORT', 'สร้างรายงานสรุปประจำวัน', '0 45 6 ? * *', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('batch003', 'อัปเดตสถานะงานค้าง', 'UPDATE', 'ตรวจสอบและอัปเดตสถานะงานที่ค้างนานเกินไป', '0 30 6 ? * *', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('batch004', 'ล้างข้อมูล/ซิงค์ฐานข้อมูล', 'CLEANUP', 'ล้างข้อมูลเก่าและซิงค์ฐานข้อมูล', '0 0 3 ? * *', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('batch005', 'ซิงค์ข้อมูล Realtime', 'SYNC', 'ซิงค์ข้อมูล Realtime กับระบบภายนอก', '0 0/30 * * * ?', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('batch006', 'ส่งสรุปยอดขาย', 'SUMMARY', 'ส่งสรุปยอดขายประจำวัน', '0 30 6 ? * *', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('batch007', 'แจ้งเตือนสต็อกต่ำ', 'INVENTORY', 'ตรวจสอบสินค้าต่ำสต็อกและส่งแจ้งเตือน', '0 30 6 ? * *', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001');

-- เพิ่มประวัติการรันตัวอย่าง
INSERT INTO t_batch_job_history (job_code, started_at, finished_at, status, result_summary, records_processed, duration_ms, whitelabel_id)
VALUES 
('batch001', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day' + INTERVAL '2 minutes', 'COMPLETED', 'Sent 10 notification emails', 10, 120000, '00000000-0000-0000-0000-000000000001'),
('batch006', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day' + INTERVAL '3 minutes', 'COMPLETED', 'Sales summary sent to 5 recipients', 5, 180000, '00000000-0000-0000-0000-000000000001');

2. สร้าง Demo Data สำหรับ Inventory Alert History
INSERT INTO t_inventory_alert_history (alert_date, part_id, part_code, part_name, current_stock, reorder_level, reorder_quantity, alert_sent, resolved, whitelabel_id)
VALUES 
(CURRENT_DATE, (SELECT id FROM m_part_master WHERE part_code='OIL-FILTER-001'), 'OIL-FILTER-001', 'ไส้กรองน้ำมันเครื่อง', 5, 10, 20, true, false, '00000000-0000-0000-0000-000000000001');

3. Run Test ทั้งหมดอีกครั้งเพื่อยืนยันว่าผ่านทุก Test Case
```

---

### Prompt 5: เขียน Manual แยกแต่ละ Module และตัวอย่างการใช้งาน

```
1. เขียน Manual สำหรับโมดูล Batch Jobs (แยกไฟล์) ให้ครบถ้วน ประกอบด้วย:
   - ภาพรวมและฟังก์ชันหลัก
   - โครงสร้างโมดูล
   - การติดตั้งและ Configuration
   - API Endpoints พร้อมตัวอย่าง Request/Response
   - การใช้งาน Batch Jobs (Cron Expression)
   - การ Manual Trigger
   - การตรวจสอบสถานะและประวัติ
   - การแก้ไขปัญหาเบื้องต้น

2. ตัวอย่างการใช้งานแต่ละ API:
   - GET /api/v1/batch-jobs
   - GET /api/v1/batch-jobs/{jobCode}/status
   - POST /api/v1/batch-jobs/{jobCode}/trigger
   - POST /api/v1/batch-jobs/{jobCode}/stop
   - PUT /api/v1/batch-jobs/{jobCode}/toggle
   - GET /api/v1/batch-jobs/{jobCode}/history
   - GET /api/v1/batch-jobs/history/all

3. ตัวอย่างการใช้งาน batch007 (Low Stock Alert):
   - ตั้งค่า Cron Expression
   - การทำงานร่วมกับ Inventory Module
   - การรับอีเมลแจ้งเตือน
```

---

### Prompt 6: Update README.md ให้ครบถ้วน

```
1. อัปเดต README.md ใน Root Project ให้รวมข้อมูลโมดูล Batch Jobs:
   - เพิ่มโมดูล Batch Jobs ในตารางสรุปโมดูล
   - เพิ่มคำอธิบายฟังก์ชันหลัก
   - เพิ่ม Batch Schedule ตาราง (7 งาน)
   - เพิ่มตัวอย่างการใช้งานเบื้องต้น

2. เขียนตัวอย่างการใช้งานของ Batch Jobs:
   - การเรียกดูรายการงานทั้งหมด
   - การสั่งรันงานด้วยตนเอง
   - การตรวจสอบสถานะงาน
   - การดูประวัติการรันงาน

3. อัปเดตตารางสรุป ณ ตอนนี้มี 16 โมดูล
```

---

### Prompt 7: Update Swagger ให้ครบถ้วน

```
1. อัปเดต Swagger Configuration ในโมดูล Batch Jobs:
   - เพิ่ม @Tag สำหรับ Batch Jobs
   - เพิ่ม @Operation สำหรับทุก API Endpoint
   - เพิ่ม @Parameter สำหรับ Request Parameters
   - เพิ่ม @Schema สำหรับ DTOs

2. ตรวจสอบให้แน่ใจว่า Swagger UI แสดง API ทั้งหมดของ Batch Jobs อย่างถูกต้อง

3. เพิ่มตัวอย่าง Request/Response ใน Swagger (ใช้ @ExampleObject)
```

---

### Prompt 8: ทดสอบ API และแก้ไขปัญหา

```
1. ทดสอบ API ทั้งหมดของโมดูล Batch Jobs ด้วยเครื่องมือ (Postman / cURL):

   ทดสอบ GET /api/v1/batch-jobs
   curl -X GET "http://localhost:5000/api/v1/batch-jobs" -H "Authorization: Bearer {token}"

   ทดสอบ GET /api/v1/batch-jobs/batch001/status
   curl -X GET "http://localhost:5000/api/v1/batch-jobs/batch001/status" -H "Authorization: Bearer {token}"

   ทดสอบ POST /api/v1/batch-jobs/batch001/trigger
   curl -X POST "http://localhost:5000/api/v1/batch-jobs/batch001/trigger" -H "Authorization: Bearer {token}" -H "Content-Type: application/json" -d '{"note":"Manual trigger for testing"}'

   ทดสอบ PUT /api/v1/batch-jobs/batch001/toggle?enabled=false
   curl -X PUT "http://localhost:5000/api/v1/batch-jobs/batch001/toggle?enabled=false" -H "Authorization: Bearer {token}"

   ทดสอบ GET /api/v1/batch-jobs/batch001/history
   curl -X GET "http://localhost:5000/api/v1/batch-jobs/batch001/history" -H "Authorization: Bearer {token}"

2. หากพบปัญหา ให้ทำการแก้ไขและบันทึกไว้

3. Update เอกสารให้ครบถ้วน (Manual, README, Swagger)
```

---

## 📋 สรุป Prompts ทั้งหมด

| Prompt | หัวข้อ | คำอธิบาย |
|--------|-------|----------|
| Prompt 1 | ดำเนินการตามเอกสาร | สร้าง Folder, File และ Code ทั้งหมด |
| Prompt 2 | สร้างโครงสร้างตามเอกสาร | สร้าง Folder/File และ Code ตามตัวอย่าง |
| Prompt 3 | สร้าง Unit Test | สร้าง Test และ Run ให้ผ่านทุก Test Case |
| Prompt 4 | เพิ่ม Demo Data | เพิ่มข้อมูลตัวอย่างและ Run Test ผ่าน |
| Prompt 5 | เขียน Manual | เขียน Manual แยกแต่ละ Module |
| Prompt 6 | Update README.md | Update README และตัวอย่างการใช้งาน |
| Prompt 7 | Update Swagger | Update Swagger ให้ครบถ้วน |
| Prompt 8 | ทดสอบ API | ทดสอบ API และแก้ไขปัญหา |

---

## 🚀 วิธีเริ่มดำเนินการ

คุณสามารถเริ่มต้นด้วย Prompt 1 ทันที โดยข้าพเจ้าจะสร้างไฟล์ทั้งหมดตามเอกสารให้คุณ หรือคุณจะเลือกเริ่มจากข้อใดก็ได้ แล้วข้าพเจ้าจะดำเนินการตามลำดับ

**กรุณาระบุว่าต้องการเริ่มจากข้อใด หรือให้ดำเนินการทั้งหมดตามลำดับ**

```
เริ่มดำเนินการตาม Prompt 1 - 8 ตามลำดับ
```

หรือ

```
ขอเริ่มที่ Prompt 3 (Unit Test) ก่อน
```

 