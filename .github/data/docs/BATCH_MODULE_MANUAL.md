# Batch Jobs Module Manual

## ภาพรวม (Overview)

โมดูล Batch Jobs เป็นระบบงานอัตโนมัติที่ทำงานตามตารางเวลาที่กำหนด (Cron) สำหรับงานประจำต่างๆ มีทั้งหมด 7 งาน:

| รหัสงาน | ชื่องาน | ประเภท | เวลาทำงาน | คำอธิบาย |
|---------|--------|--------|-----------|----------|
| batch001 | ส่งอีเมลแจ้งเตือนรายวัน | EMAIL | 06:30 น. | ส่งอีเมลแจ้งเตือนรายวันให้พนักงาน |
| batch002 | สร้างรายงานประจำวัน | REPORT | 06:45 น. | สร้างรายงานสรุปประจำวัน |
| batch003 | อัปเดตสถานะงานค้าง | UPDATE | 06:30 น. | ตรวจสอบและอัปเดตสถานะงานที่ค้าง |
| batch004 | ล้างข้อมูล/ซิงค์ฐานข้อมูล | CLEANUP | 03:00 น. | ล้างข้อมูลเก่าและซิงค์ฐานข้อมูล |
| batch005 | ซิงค์ข้อมูล Realtime | SYNC | ทุก 30 นาที | ซิงค์ข้อมูลกับระบบภายนอก |
| batch006 | ส่งสรุปยอดขาย | SUMMARY | 06:30 น. | ส่งสรุปยอดขายประจำวัน |
| batch007 | แจ้งเตือนสต็อกต่ำ | INVENTORY | 06:30 น. | ตรวจสอบสินค้าต่ำสต็อกและแจ้งเตือน |

## โครงสร้างโมดูล (Module Structure)

```
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
│   │   ├── BatchJobRepository.java (JPA)
│   │   ├── BatchJobHistoryRepository.java (JPA)
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
│   ├── lock/DistributedLockService.java
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
```

## API Endpoints

| Method | Path | Description | Rate Limit |
|--------|------|-------------|------------|
| GET | `/api/v1/batch-jobs` | แสดงรายการงานทั้งหมด | 20 ครั้ง/นาที |
| GET | `/api/v1/batch-jobs/{jobCode}/status` | ดึงสถานะงาน | 30 ครั้ง/นาที |
| POST | `/api/v1/batch-jobs/{jobCode}/trigger` | สั่งรันงาน | 5 ครั้ง/ชม. |
| POST | `/api/v1/batch-jobs/{jobCode}/stop` | หยุดงานที่กำลังรัน | 3 ครั้ง/ชม. |
| PUT | `/api/v1/batch-jobs/{jobCode}/toggle` | เปิด/ปิดใช้งานงาน | 10 ครั้ง/5 นาที |
| GET | `/api/v1/batch-jobs/{jobCode}/history` | ประวัติการรันของงาน | 20 ครั้ง/นาที |
| GET | `/api/v1/batch-jobs/history/all` | ประวัติการรันทั้งหมด | 15 ครั้ง/นาที |

## ตัวอย่างการใช้งาน (Usage Examples)

### 1. Login เพื่อขอ Token
```bash
TOKEN=$(curl -s -X POST "http://localhost:5000/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}' | jq -r '.accessToken')
```

### 2. ดูรายการงานทั้งหมด
```bash
curl -X GET "http://localhost:5000/api/v1/batch-jobs" \
  -H "Authorization: Bearer $TOKEN"
```

### 3. ดูสถานะงาน
```bash
curl -X GET "http://localhost:5000/api/v1/batch-jobs/batch001/status" \
  -H "Authorization: Bearer $TOKEN"
```

### 4. สั่งรันงานด้วยตนเอง
```bash
curl -X POST "http://localhost:5000/api/v1/batch-jobs/batch001/trigger" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"note":"Manual trigger for testing"}'
```

### 5. หยุดงานที่กำลังทำงาน
```bash
curl -X POST "http://localhost:5000/api/v1/batch-jobs/batch001/stop" \
  -H "Authorization: Bearer $TOKEN"
```

### 6. เปิด/ปิดใช้งานงาน
```bash
# ปิดใช้งาน
curl -X PUT "http://localhost:5000/api/v1/batch-jobs/batch001/toggle?enabled=false" \
  -H "Authorization: Bearer $TOKEN"

# เปิดใช้งาน
curl -X PUT "http://localhost:5000/api/v1/batch-jobs/batch001/toggle?enabled=true" \
  -H "Authorization: Bearer $TOKEN"
```

### 7. ดูประวัติการรัน
```bash
# ประวัติของงานเดียว
curl -X GET "http://localhost:5000/api/v1/batch-jobs/batch001/history" \
  -H "Authorization: Bearer $TOKEN"

# ประวัติทั้งหมด
curl -X GET "http://localhost:5000/api/v1/batch-jobs/history/all" \
  -H "Authorization: Bearer $TOKEN"

# ประวัติทั้งหมดพร้อมกรองสถานะ
curl -X GET "http://localhost:5000/api/v1/batch-jobs/history/all?status=FAILED" \
  -H "Authorization: Bearer $TOKEN"
```

## ตัวอย่าง Response

### รายการงานทั้งหมด
```json
[
  {
    "id": "uuid",
    "jobCode": "batch001",
    "jobName": "ส่งอีเมลแจ้งเตือนรายวัน",
    "jobType": "EMAIL",
    "cronExpression": "0 30 6 ? * *",
    "enabled": true,
    "lastStatus": "COMPLETED",
    "lastRunTime": "2024-01-15T06:30:00",
    "nextRunTime": "2024-01-16T06:30:00",
    "totalRuns": 150
  }
]
```

### สถานะงาน
```json
{
  "jobCode": "batch001",
  "jobName": "ส่งอีเมลแจ้งเตือนรายวัน",
  "enabled": true,
  "lastStatus": "COMPLETED",
  "lastRunTime": "2024-01-15T06:30:00",
  "nextRunTime": "2024-01-16T06:30:00",
  "totalRuns": 150
}
```

### ประวัติการรัน
```json
{
  "content": [
    {
      "id": "uuid",
      "jobCode": "batch001",
      "startedAt": "2024-01-15T06:30:00",
      "finishedAt": "2024-01-15T06:32:00",
      "status": "COMPLETED",
      "resultSummary": "ส่งอีเมลแจ้งเตือน 10 ฉบับ",
      "recordsProcessed": 10,
      "durationMs": 120000,
      "triggerType": "SCHEDULED"
    }
  ],
  "totalElements": 1,
  "totalPages": 1
}
```

## การเชื่อมต่อกับโมดูลอื่น

### batch007 - Low Stock Alert (เชื่อมต่อ Inventory)

งาน batch007 จะตรวจสอบสินค้าที่ต่ำกว่า Reorder Level ในโมดูล Inventory และสร้าง Alert รวมทั้งส่งอีเมลแจ้งเตือน

**ขั้นตอนการทำงาน:**
1. เรียก `PartMasterRepository.findLowStockParts()` เพื่อหาสินค้าต่ำสต็อก
2. สร้าง `InventoryAlertEntity` สำหรับสินค้าที่ต่ำเกณฑ์
3. ส่งอีเมลสรุปรายการไปยังแผนกจัดซื้อ

### การตั้งค่า Cron Expression

| ตัวอย่าง | คำอธิบาย |
|----------|----------|
| `0 30 6 ? * *` | ทุกวันเวลา 06:30 น. |
| `0 45 6 ? * *` | ทุกวันเวลา 06:45 น. |
| `0 0 3 ? * *` | ทุกวันเวลา 03:00 น. |
| `0 0/30 * * * ?` | ทุก 30 นาที |
| `0 */5 * * * *` | ทุก 5 นาที |

## การแก้ไขปัญหาเบื้องต้น (Troubleshooting)

| ปัญหา | สาเหตุ | วิธีแก้ไข |
|-------|--------|-----------|
| งานไม่ทำงาน | Redis ไม่พร้อมใช้งาน | ตรวจสอบ `spring.redis.host` และ Redis service |
| งานไม่ทำงาน | `enabled = false` | ใช้ API toggle เพื่อเปิดใช้งาน |
| งานล้มเหลว | ขอมูลเชื่อมต่อภายนอกผิดพลาด | ตรวจสอบ External API endpoint |
| Lock ติดค้าง | เกิด Exception ก่อนปลด Lock | Redis key จะหมดอายุอัตโนมัติตาม timeout ที่ตั้งไว้ |

## การพัฒนาเพิ่มเติม (Development)

### การเพิ่มงาน Batch ใหม่
1. เพิ่ม job_type ใน `BatchJobType.java` enum
2. เพิ่มข้อมูลใน `m_batch_job` table
3. เพิ่ม method ใน `BatchJobExecutor.java`
4. เพิ่ม `@Scheduled` method ใน `BatchSchedulerConfig.java`

### การรัน Test
```bash
# รันเฉพาะ Batch Module Tests
mvn test -Dtest="com.icmon.module.batch.**"
```
