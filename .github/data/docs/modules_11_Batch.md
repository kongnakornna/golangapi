**โมดูลที่ 11: ⏱️ Batch Jobs (งานตามกำหนดเวลา)**

โมดูล Batch Jobs เป็นระบบงานอัตโนมัติที่ทำงานตามตารางเวลาที่กำหนด (Cron) สำหรับงานประจำต่างๆ ที่ต้องทำเป็นประจำโดยไม่ต้องอาศัยการแทรกแซงของมนุษย์ ครอบคลุมการทำงานดังนี้:

1. **งานที่ 1 (batch001):** ส่งอีเมลแจ้งเตือนรายวัน เวลา 06:30 น.
2. **งานที่ 2 (batch002):** สร้างรายงานประจำวัน เวลา 06:45 น.
3. **งานที่ 3 (batch003):** อัปเดตสถานะงานค้าง เวลา 06:30 น.
4. **งานที่ 4 (batch004):** ล้างข้อมูล/ซิงค์ฐานข้อมูล เวลา 03:00 น. (กลางคืน)
5. **งานที่ 5 (batch005):** ซิงค์ข้อมูล Realtime ทุก 30 นาที
6. **งานที่ 6 (batch006):** ส่งสรุปยอดขาย เวลา 06:30 น.

---

## 📁 โครงสร้างโมดูล Batch Jobs (`modules/batch`)

```
modules/batch/
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
│   │   ├── BatchJobStatus.java        // SCHEDULED, RUNNING, COMPLETED, FAILED, CANCELLED
│   │   ├── BatchJobType.java          // EMAIL, REPORT, SYNC, CLEANUP, UPDATE, SUMMARY
│   │   └── BatchJobPriority.java      // LOW, NORMAL, HIGH, CRITICAL
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
│   ├── cache/                                           // ⬅️ ระบบ Cache สำหรับ Batch
│   │   ├── BatchJobCacheService.java
│   │   └── BatchJobLockCacheService.java
│   ├── scheduler/                                       // ⬅️ Scheduler Core
│   │   ├── BatchSchedulerConfig.java
│   │   ├── Job001DailyNotification.java
│   │   ├── Job002DailyReport.java
│   │   ├── Job003UpdatePendingStatus.java
│   │   ├── Job004CleanupAndSync.java
│   │   ├── Job005RealtimeSync.java
│   │   └── Job006DailySalesSummary.java
│   ├── lock/                                            // ⬅️ Distributed Lock
│   │   └── DistributedLockService.java
│   ├── executor/                                        // ⬅️ Job Executor
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
    │   ├── BatchJobController.java      // Job Management APIs
    │   └── BatchJobHistoryController.java // History APIs
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

---

## 🗄️ ออกแบบฐานข้อมูล (Database Design) สำหรับโมดูล Batch Jobs

### 📄 SQL DDL (ไฟล์: `src/main/resources/db/migration/V11__batch_schema.sql`)

```sql
-- ==============================================
-- ตาราง: m_batch_job (ข้อมูลงานประจำ)
-- Batch job master data.
-- ==============================================
CREATE TABLE IF NOT EXISTS m_batch_job (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_code VARCHAR(30) UNIQUE NOT NULL,           -- รหัสงาน (เช่น batch001, batch002)
    job_name VARCHAR(100) NOT NULL,                 -- ชื่องาน
    job_type VARCHAR(30) NOT NULL,                  -- EMAIL, REPORT, SYNC, CLEANUP, UPDATE, SUMMARY
    description TEXT,                               -- คำอธิบายงาน
    cron_expression VARCHAR(100) NOT NULL,          -- Cron expression (เช่น 0 30 6 ? * *)
    enabled BOOLEAN DEFAULT TRUE,                   -- เปิดใช้งานหรือไม่
    max_retry INTEGER DEFAULT 3,                    -- จำนวนครั้งที่ลองใหม่
    retry_delay_ms INTEGER DEFAULT 60000,           -- หน่วงเวลาก่อนลองใหม่ (ms)
    timeout_seconds INTEGER DEFAULT 300,            -- หมดเวลา (วินาที)
    last_run_time TIMESTAMP,                        -- เวลาที่รันล่าสุด
    next_run_time TIMESTAMP,                        -- เวลาที่จะรันครั้งถัดไป
    total_runs INTEGER DEFAULT 0,                   -- จำนวนครั้งที่รันทั้งหมด
    last_status VARCHAR(20),                        -- สถานะล่าสุด
    parameters JSONB,                               -- พารามิเตอร์เพิ่มเติม
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_m_batch_job_code ON m_batch_job(job_code);
CREATE INDEX idx_m_batch_job_enabled ON m_batch_job(enabled);
CREATE INDEX idx_m_batch_job_type ON m_batch_job(job_type);

-- ==============================================
-- ข้อมูลเริ่มต้น: งานประจำทั้ง 6 งาน
-- Initial data: All 6 batch jobs.
-- ==============================================
INSERT INTO m_batch_job (job_code, job_name, job_type, description, cron_expression, enabled, user_id, whitelabel_id)
VALUES 
('batch001', 'ส่งอีเมลแจ้งเตือนรายวัน', 'EMAIL', 'ส่งอีเมลแจ้งเตือนรายวันให้กับพนักงาน', '0 30 6 ? * *', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('batch002', 'สร้างรายงานประจำวัน', 'REPORT', 'สร้างรายงานสรุปประจำวัน', '0 45 6 ? * *', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('batch003', 'อัปเดตสถานะงานค้าง', 'UPDATE', 'ตรวจสอบและอัปเดตสถานะงานที่ค้างนานเกินไป', '0 30 6 ? * *', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('batch004', 'ล้างข้อมูล/ซิงค์ฐานข้อมูล', 'CLEANUP', 'ล้างข้อมูลเก่าและซิงค์ฐานข้อมูล', '0 0 3 ? * *', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('batch005', 'ซิงค์ข้อมูล Realtime', 'SYNC', 'ซิงค์ข้อมูล Realtime กับระบบภายนอก', '0 0/30 * * * ?', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('batch006', 'ส่งสรุปยอดขาย', 'SUMMARY', 'ส่งสรุปยอดขายประจำวัน', '0 30 6 ? * *', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001');

-- ==============================================
-- ตาราง: t_batch_job_history (ประวัติการรันงาน)
-- Batch job execution history.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_batch_job_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_code VARCHAR(30) NOT NULL REFERENCES m_batch_job(job_code) ON DELETE CASCADE,
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMP,
    status VARCHAR(20) NOT NULL,                -- SCHEDULED, RUNNING, COMPLETED, FAILED, CANCELLED
    error_message TEXT,
    result_summary TEXT,                        -- สรุปผลการทำงาน
    records_processed INTEGER DEFAULT 0,        -- จำนวนรายการที่ประมวลผล
    duration_ms INTEGER,                        -- ระยะเวลาที่ใช้ (มิลลิวินาที)
    trigger_type VARCHAR(20) DEFAULT 'SCHEDULED', -- SCHEDULED, MANUAL, RETRY
    triggered_by UUID REFERENCES m_user(id),    -- ผู้ที่สั่งรัน (ถ้า manual)
    parameters JSONB,                           -- พารามิเตอร์ที่ใช้
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_batch_history_job ON t_batch_job_history(job_code);
CREATE INDEX idx_t_batch_history_status ON t_batch_job_history(status);
CREATE INDEX idx_t_batch_history_started ON t_batch_job_history(started_at);
CREATE INDEX idx_t_batch_history_whitelabel ON t_batch_job_history(whitelabel_id);
```

---

## 🧠 ระบบ Cache และ Distributed Lock สำหรับ Batch Jobs

### `infrastructure/cache/BatchJobLockCacheService.java`

```java
package com.template.app.modules.batch.infrastructure.cache;

import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.stereotype.Service;

import java.time.Duration;
import java.util.UUID;

@Service
public class BatchJobLockCacheService {

    private final RedisTemplate<String, String> redisTemplate;

    public BatchJobLockCacheService(RedisTemplate<String, String> redisTemplate) {
        this.redisTemplate = redisTemplate;
    }

    /*
        ฟังก์ชันนี้ใช้ Distributed Lock เพื่อป้องกันไม่ให้งานเดียวกันทำงานพร้อมกัน (ใช้ Redis)
        This function uses Distributed Lock to prevent the same job from running concurrently (uses Redis).
        Redis Key: batch_lock:{jobCode}
    */
    public boolean acquireLock(String jobCode, String instanceId, Duration timeout) {
        String key = "batch_lock:" + jobCode;
        Boolean success = redisTemplate.opsForValue()
                .setIfAbsent(key, instanceId, timeout);
        return Boolean.TRUE.equals(success);
    }

    /*
        ฟังก์ชันนี้ปลดล็อกเมื่อทำงานเสร็จ
        This function releases the lock when job is done.
    */
    public void releaseLock(String jobCode) {
        String key = "batch_lock:" + jobCode;
        redisTemplate.delete(key);
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่างานถูกล็อกอยู่หรือไม่
        This function checks if the job is locked.
    */
    public boolean isLocked(String jobCode) {
        String key = "batch_lock:" + jobCode;
        return Boolean.TRUE.equals(redisTemplate.hasKey(key));
    }
}
```

### `infrastructure/cache/BatchJobCacheService.java`

```java
package com.template.app.modules.batch.infrastructure.cache;

import com.template.app.modules.batch.domain.MBatchJob;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.CachePut;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

@Service
public class BatchJobCacheService {

    /*
        ฟังก์ชันนี้ดึงข้อมูลงาน Batch จาก Cache
        This function retrieves batch job data from cache.
        Redis Key: batch_job:{jobCode}
    */
    @Cacheable(value = "batch_jobs", key = "#jobCode")
    public MBatchJob getBatchJob(String jobCode) {
        return null;
    }

    /*
        ฟังก์ชันนี้จะอัปเดต Cache เมื่อมีการบันทึกงาน
        This function updates the cache when a job is saved.
    */
    @CachePut(value = "batch_jobs", key = "#job.jobCode")
    public MBatchJob saveBatchJob(MBatchJob job) {
        return job;
    }

    /*
        ฟังก์ชันนี้ลบข้อมูลงานออกจาก Cache
        This function evicts batch job data from cache.
    */
    @CacheEvict(value = "batch_jobs", key = "#jobCode")
    public void evictBatchJob(String jobCode) {
        // ลบ Cache / Evict cache.
    }

    /*
        ฟังก์ชันนี้ลบ Cache ทั้งหมดของ Batch Jobs
        This function clears all batch job caches.
    */
    @CacheEvict(value = "batch_jobs", allEntries = true)
    public void evictAllBatchJobs() {
        // ลบทุก key / Evict all keys.
    }
}
```

---

## ⏱️ Scheduler Configuration และงานทั้ง 6

### `infrastructure/scheduler/BatchSchedulerConfig.java`

```java
package com.template.app.modules.batch.infrastructure.scheduler;

import com.template.app.modules.batch.infrastructure.cache.BatchJobLockCacheService;
import com.template.app.modules.batch.infrastructure.executor.BatchJobExecutor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.context.annotation.Configuration;
import org.springframework.scheduling.annotation.EnableScheduling;
import org.springframework.scheduling.annotation.Scheduled;

import java.time.Duration;
import java.util.UUID;

@Configuration
@EnableScheduling
public class BatchSchedulerConfig {

    private static final Logger logger = LoggerFactory.getLogger(BatchSchedulerConfig.class);
    private final BatchJobExecutor batchJobExecutor;
    private final BatchJobLockCacheService lockCacheService;

    // ใช้ instanceId เพื่อระบุตัวตนของแต่ละ instance (กรณี Multi-instance)
    private final String instanceId = UUID.randomUUID().toString();

    public BatchSchedulerConfig(BatchJobExecutor batchJobExecutor,
                                BatchJobLockCacheService lockCacheService) {
        this.batchJobExecutor = batchJobExecutor;
        this.lockCacheService = lockCacheService;
    }

    // ========================================================================
    // batch001: ส่งอีเมลแจ้งเตือนรายวัน (06:30 น. ทุกวัน)
    // Job001: Send daily email notifications (06:30 AM daily)
    // ========================================================================
    @Scheduled(cron = "0 30 6 ? * *")
    public void executeBatch001() {
        executeWithLock("batch001", () -> batchJobExecutor.executeJob("batch001"));
    }

    // ========================================================================
    // batch002: สร้างรายงานประจำวัน (06:45 น. ทุกวัน)
    // Job002: Generate daily reports (06:45 AM daily)
    // ========================================================================
    @Scheduled(cron = "0 45 6 ? * *")
    public void executeBatch002() {
        executeWithLock("batch002", () -> batchJobExecutor.executeJob("batch002"));
    }

    // ========================================================================
    // batch003: อัปเดตสถานะงานค้าง (06:30 น. ทุกวัน)
    // Job003: Update pending job statuses (06:30 AM daily)
    // ========================================================================
    @Scheduled(cron = "0 30 6 ? * *")
    public void executeBatch003() {
        executeWithLock("batch003", () -> batchJobExecutor.executeJob("batch003"));
    }

    // ========================================================================
    // batch004: ล้างข้อมูล/ซิงค์ฐานข้อมูล (03:00 น. ทุกวัน)
    // Job004: Cleanup and sync database (03:00 AM daily)
    // ========================================================================
    @Scheduled(cron = "0 0 3 ? * *")
    public void executeBatch004() {
        executeWithLock("batch004", () -> batchJobExecutor.executeJob("batch004"));
    }

    // ========================================================================
    // batch005: ซิงค์ข้อมูล Realtime (ทุก 30 นาที)
    // Job005: Realtime data sync (every 30 minutes)
    // ========================================================================
    @Scheduled(cron = "0 0/30 * * * ?")
    public void executeBatch005() {
        executeWithLock("batch005", () -> batchJobExecutor.executeJob("batch005"));
    }

    // ========================================================================
    // batch006: ส่งสรุปยอดขาย (06:30 น. ทุกวัน)
    // Job006: Send daily sales summary (06:30 AM daily)
    // ========================================================================
    @Scheduled(cron = "0 30 6 ? * *")
    public void executeBatch006() {
        executeWithLock("batch006", () -> batchJobExecutor.executeJob("batch006"));
    }

    /*
        ฟังก์ชันนี้ใช้ Distributed Lock เพื่อป้องกันการทำงานซ้ำของงานเดียวกัน
        This function uses Distributed Lock to prevent duplicate execution of the same job.
    */
    private void executeWithLock(String jobCode, Runnable job) {
        // 1. พยายามล็อก (หมดเวลา 5 นาที) / Try to acquire lock (timeout 5 minutes).
        if (!lockCacheService.acquireLock(jobCode, instanceId, Duration.ofMinutes(5))) {
            logger.warn("Job {} is already running on another instance. Skipping.", jobCode);
            return;
        }

        try {
            logger.info("Executing job: {} on instance: {}", jobCode, instanceId);
            job.run();
        } catch (Exception e) {
            logger.error("Error executing job {}: {}", jobCode, e.getMessage(), e);
        } finally {
            // 2. ปลดล็อกเมื่อเสร็จ / Release lock when done.
            lockCacheService.releaseLock(jobCode);
            logger.info("Released lock for job: {}", jobCode);
        }
    }
}
```

---

## 🔧 ตัวอย่าง Implementation ของแต่ละงาน

### `infrastructure/executor/BatchJobExecutor.java`

```java
package com.template.app.modules.batch.infrastructure.executor;

import com.template.app.modules.batch.application.interfaces.BatchJobService;
import com.template.app.modules.batch.domain.MBatchJob;
import com.template.app.modules.batch.domain.TBatchJobHistory;
import com.template.app.modules.batch.domain.enums.BatchJobStatus;
import com.template.app.modules.batch.infrastructure.cache.BatchJobCacheService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;

@Service
public class BatchJobExecutor {

    private static final Logger logger = LoggerFactory.getLogger(BatchJobExecutor.class);
    private final BatchJobService batchJobService;
    private final BatchJobCacheService cacheService;

    public BatchJobExecutor(BatchJobService batchJobService, BatchJobCacheService cacheService) {
        this.batchJobService = batchJobService;
        this.cacheService = cacheService;
    }

    /*
        ฟังก์ชันนี้เป็นตัวกลางในการรันงาน Batch แต่ละงาน
        This function is the central executor for each batch job.
    */
    public void executeJob(String jobCode) {
        try {
            // 1. ดึงข้อมูลงาน / Get job data.
            MBatchJob job = batchJobService.getJob(jobCode);
            if (job == null || !job.getEnabled()) {
                logger.info("Job {} is disabled or not found. Skipping.", jobCode);
                return;
            }

            // 2. ตรวจสอบว่าควรจะรันตอนนี้หรือไม่ (ใช้ next_run_time) / Check if should run now.
            if (job.getNextRunTime() != null && job.getNextRunTime().isAfter(LocalDateTime.now())) {
                logger.info("Job {} is scheduled for {}. Skipping.", jobCode, job.getNextRunTime());
                return;
            }

            // 3. สร้างประวัติการรัน / Create execution history.
            TBatchJobHistory history = new TBatchJobHistory();
            history.setJobCode(jobCode);
            history.setStartedAt(LocalDateTime.now());
            history.setStatus(BatchJobStatus.RUNNING);
            history.setTriggerType("SCHEDULED");

            // 4. บันทึกประวัติก่อนรัน / Save history before running.
            TBatchJobHistory savedHistory = batchJobService.saveHistory(history);

            // 5. รันงานตามประเภท / Execute job based on type.
            boolean success = executeJobByType(job, savedHistory);

            // 6. อัปเดตประวัติหลังรันเสร็จ / Update history after execution.
            savedHistory.setFinishedAt(LocalDateTime.now());
            savedHistory.setStatus(success ? BatchJobStatus.COMPLETED : BatchJobStatus.FAILED);
            savedHistory.setDurationMs(
                (int) (savedHistory.getFinishedAt().toInstant().toEpochMilli() 
                - savedHistory.getStartedAt().toInstant().toEpochMilli())
            );
            batchJobService.updateHistory(savedHistory);

            // 7. อัปเดตข้อมูลงานหลัก / Update main job data.
            job.setLastRunTime(savedHistory.getStartedAt());
            job.setNextRunTime(calculateNextRunTime(job.getCronExpression()));
            job.setTotalRuns(job.getTotalRuns() + 1);
            job.setLastStatus(savedHistory.getStatus().name());
            batchJobService.updateJob(job);
            cacheService.saveBatchJob(job);

            logger.info("Job {} completed with status: {}", jobCode, savedHistory.getStatus());

        } catch (Exception e) {
            logger.error("Failed to execute job {}: {}", jobCode, e.getMessage(), e);
            // บันทึกข้อผิดพลาด / Record error.
            batchJobService.recordFailure(jobCode, e.getMessage());
        }
    }

    /*
        ฟังก์ชันนี้แยกประเภทงานและเรียกใช้งานตามประเภท
        This function routes to the appropriate job type implementation.
    */
    private boolean executeJobByType(MBatchJob job, TBatchJobHistory history) {
        switch (job.getJobType()) {
            case EMAIL:
                return executeEmailJob(job, history);
            case REPORT:
                return executeReportJob(job, history);
            case UPDATE:
                return executeUpdateJob(job, history);
            case CLEANUP:
                return executeCleanupJob(job, history);
            case SYNC:
                return executeSyncJob(job, history);
            case SUMMARY:
                return executeSummaryJob(job, history);
            default:
                logger.warn("Unknown job type: {}", job.getJobType());
                return false;
        }
    }

    // ========================================================================
    // batch001: EMAIL - ส่งอีเมลแจ้งเตือนรายวัน
    // ========================================================================
    private boolean executeEmailJob(MBatchJob job, TBatchJobHistory history) {
        logger.info("Executing Email Job: {}", job.getJobCode());
        // TODO: เรียกใช้ EmailService ส่งอีเมลแจ้งเตือน
        // emailService.sendDailyNotification();
        history.setResultSummary("Sent daily notification emails");
        history.setRecordsProcessed(100);
        return true;
    }

    // ========================================================================
    // batch002: REPORT - สร้างรายงานประจำวัน
    // ========================================================================
    private boolean executeReportJob(MBatchJob job, TBatchJobHistory history) {
        logger.info("Executing Report Job: {}", job.getJobCode());
        // TODO: เรียกใช้ ReportService สร้างรายงานประจำวัน
        // reportService.generateDailyReport();
        history.setResultSummary("Generated daily report");
        history.setRecordsProcessed(500);
        return true;
    }

    // ========================================================================
    // batch003: UPDATE - อัปเดตสถานะงานค้าง
    // ========================================================================
    private boolean executeUpdateJob(MBatchJob job, TBatchJobHistory history) {
        logger.info("Executing Update Job: {}", job.getJobCode());
        // TODO: เรียกใช้ JobService อัปเดตสถานะงานที่ค้างนานเกินไป
        // int updated = jobService.updatePendingJobs();
        int updated = 25;
        history.setResultSummary("Updated " + updated + " pending jobs");
        history.setRecordsProcessed(updated);
        return true;
    }

    // ========================================================================
    // batch004: CLEANUP - ล้างข้อมูล/ซิงค์ฐานข้อมูล
    // ========================================================================
    private boolean executeCleanupJob(MBatchJob job, TBatchJobHistory history) {
        logger.info("Executing Cleanup Job: {}", job.getJobCode());
        // TODO: ล้างข้อมูลเก่า (เช่น log, temp) และซิงค์ฐานข้อมูล
        // int deleted = cleanupService.cleanOldData();
        // syncService.syncDatabase();
        int deleted = 1500;
        history.setResultSummary("Cleaned up " + deleted + " old records and synced database");
        history.setRecordsProcessed(deleted);
        return true;
    }

    // ========================================================================
    // batch005: SYNC - ซิงค์ข้อมูล Realtime
    // ========================================================================
    private boolean executeSyncJob(MBatchJob job, TBatchJobHistory history) {
        logger.info("Executing Sync Job: {}", job.getJobCode());
        // TODO: ซิงค์ข้อมูลกับระบบภายนอก (Kafka, External API)
        // syncService.syncRealtimeData();
        int synced = 200;
        history.setResultSummary("Synced " + synced + " records in real-time");
        history.setRecordsProcessed(synced);
        return true;
    }

    // ========================================================================
    // batch006: SUMMARY - ส่งสรุปยอดขาย
    // ========================================================================
    private boolean executeSummaryJob(MBatchJob job, TBatchJobHistory history) {
        logger.info("Executing Summary Job: {}", job.getJobCode());
        // TODO: สรุปยอดขายและส่งอีเมล
        // summaryService.sendDailySalesSummary();
        history.setResultSummary("Sent daily sales summary");
        history.setRecordsProcessed(10);
        return true;
    }

    /*
        ฟังก์ชันนี้คำนวณเวลาที่จะรันครั้งถัดไปจาก Cron Expression
        This function calculates the next run time from Cron Expression.
    */
    private LocalDateTime calculateNextRunTime(String cronExpression) {
        // TODO: ใช้ CronParser หรือ Quartz เพื่อคำนวณ next run time
        // ตัวอย่าง: ใช้ Spring CronExpression
        try {
            org.springframework.scheduling.support.CronExpression cron = 
                org.springframework.scheduling.support.CronExpression.parse(cronExpression);
            java.time.ZonedDateTime next = cron.next(java.time.ZonedDateTime.now());
            return next != null ? next.toLocalDateTime() : null;
        } catch (Exception e) {
            logger.error("Error parsing cron expression: {}", cronExpression, e);
            return null;
        }
    }
}
```

---

## 🧩 ตัวอย่าง Domain Entity และ Enum

### `domain/enums/BatchJobStatus.java`

```java
package com.template.app.modules.batch.domain.enums;

/*
    สถานะการทำงานของ Batch Job / Batch job execution status.
*/
public enum BatchJobStatus {
    SCHEDULED,    // ถูกจัดตารางรอ / Scheduled.
    RUNNING,      // กำลังทำงาน / Running.
    COMPLETED,    // ทำงานสำเร็จ / Completed successfully.
    FAILED,       // ทำงานล้มเหลว / Failed.
    CANCELLED     // ถูกยกเลิก / Cancelled.
}
```

### `domain/enums/BatchJobType.java`

```java
package com.template.app.modules.batch.domain.enums;

/*
    ประเภทของ Batch Job / Batch job type.
*/
public enum BatchJobType {
    EMAIL,       // ส่งอีเมล / Send email.
    REPORT,      // สร้างรายงาน / Generate report.
    UPDATE,      // อัปเดตข้อมูล / Update data.
    CLEANUP,     // ล้างข้อมูล / Cleanup data.
    SYNC,        // ซิงค์ข้อมูล / Sync data.
    SUMMARY      // สรุปข้อมูล / Generate summary.
}
```

### `domain/MBatchJob.java`

```java
package com.template.app.modules.batch.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import com.template.app.modules.batch.domain.enums.BatchJobType;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.time.LocalDateTime;

@Data
@EqualsAndHashCode(callSuper = true)
public class MBatchJob extends GenericBusinessClass {

    private String jobCode;                 // รหัสงาน / Job code.
    private String jobName;                 // ชื่องาน / Job name.
    private BatchJobType jobType;           // ประเภทงาน / Job type.
    private String description;             // คำอธิบาย / Description.
    private String cronExpression;          // Cron expression.
    private Boolean enabled;                // เปิดใช้งาน / Enabled.
    private Integer maxRetry;               // จำนวนครั้งที่ลองใหม่ / Max retry.
    private Integer retryDelayMs;           // หน่วงเวลาก่อนลองใหม่ (ms) / Retry delay (ms).
    private Integer timeoutSeconds;         // หมดเวลา (วินาที) / Timeout (seconds).
    private LocalDateTime lastRunTime;      // เวลาที่รันล่าสุด / Last run time.
    private LocalDateTime nextRunTime;      // เวลาที่จะรันครั้งถัดไป / Next run time.
    private Integer totalRuns;              // จำนวนครั้งที่รันทั้งหมด / Total runs.
    private String lastStatus;              // สถานะล่าสุด / Last status.
    private String parameters;              // พารามิเตอร์ / Parameters.

    /*
        ฟังก์ชันนี้ตรวจสอบว่างานสามารถรันได้หรือไม่
        This function checks if the job can run.
    */
    public boolean canRun() {
        return this.enabled != null && this.enabled;
    }

    /*
        ฟังก์ชันนี้บันทึกการรันงานและคำนวณเวลาที่จะรันครั้งถัดไป
        This function records job execution and calculates next run time.
    */
    public void recordExecution(LocalDateTime startTime, String status) {
        this.lastRunTime = startTime;
        this.lastStatus = status;
        this.totalRuns = (this.totalRuns != null ? this.totalRuns : 0) + 1;
        // กำหนดเวลาที่จะรันครั้งถัดไป (จาก Cron) / Set next run time (from Cron).
        this.nextRunTime = calculateNextRunTime();
    }

    /*
        ฟังก์ชันนี้คำนวณเวลาที่จะรันครั้งถัดไปจาก Cron Expression
        This function calculates next run time from Cron Expression.
    */
    private LocalDateTime calculateNextRunTime() {
        // TODO: ใช้ CronParser เพื่อคำนวณ
        return null;
    }
}
```

---

## ⏱️ API Controller สำหรับ Batch Jobs

```java
package com.template.app.modules.batch.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.batch.application.interfaces.BatchJobService;
import com.template.app.modules.batch.presentation.dto.request.TriggerJobRequestDTO;
import com.template.app.modules.batch.presentation.dto.response.BatchJobHistoryResponseDTO;
import com.template.app.modules.batch.presentation.dto.response.BatchJobResponseDTO;
import com.template.app.modules.batch.presentation.dto.response.BatchJobStatusResponseDTO;
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
@RequestMapping("/api/v1/batch-jobs")
@Tag(name = "Batch Jobs", description = "Batch Job Management APIs")
@RequiredArgsConstructor
public class BatchJobController {

    private final BatchJobService batchJobService;

    // ========================================================================
    // 1. LIST ALL BATCH JOBS (แสดงงานทั้งหมด)
    // ========================================================================

    /*
        API: GET /api/v1/batch-jobs
        ฟังก์ชันนี้แสดงรายการงาน Batch ทั้งหมด พร้อมสถานะปัจจุบัน
        This function lists all batch jobs with current status.
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @GetMapping
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "List all batch jobs")
    public ResponseEntity<List<BatchJobResponseDTO>> listBatchJobs() throws SystemGlobalException {
        List<BatchJobResponseDTO> jobs = batchJobService.listAllJobs();
        return ResponseEntity.ok(jobs);
    }

    // ========================================================================
    // 2. GET JOB STATUS (ดึงสถานะงาน)
    // ========================================================================

    /*
        API: GET /api/v1/batch-jobs/{jobCode}/status
        ฟังก์ชันนี้ดึงสถานะปัจจุบันของงาน Batch ที่ระบุ (ใช้ Cache)
        This function retrieves the current status of a specified batch job (uses caching).
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที
        Rate Limit: Allows 30 requests per minute.
    */
    @GetMapping("/{jobCode}/status")
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get batch job status")
    public ResponseEntity<BatchJobStatusResponseDTO> getJobStatus(@PathVariable String jobCode)
            throws SystemGlobalException {
        BatchJobStatusResponseDTO response = batchJobService.getJobStatus(jobCode);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. TRIGGER JOB MANUALLY (สั่งรันงานด้วยตนเอง)
    // ========================================================================

    /*
        API: POST /api/v1/batch-jobs/{jobCode}/trigger
        ฟังก์ชันนี้สั่งให้รันงาน Batch ที่ระบุทันที (Manual Trigger)
        This function triggers a batch job to run immediately (Manual Trigger).
        Rate Limit: อนุญาต 5 ครั้งต่อ 1 ชั่วโมง (ป้องกันการรันซ้ำ)
        Rate Limit: Allows 5 requests per 1 hour (prevent duplicate runs).
    */
    @PostMapping("/{jobCode}/trigger")
    @RateLimit(limit = 5, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Trigger batch job manually")
    public ResponseEntity<BatchJobResponseDTO> triggerJob(
            @PathVariable String jobCode,
            @Valid @RequestBody TriggerJobRequestDTO request) throws SystemGlobalException {
        BatchJobResponseDTO response = batchJobService.triggerJob(jobCode, request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 4. STOP RUNNING JOB (หยุดงานที่กำลังทำงาน)
    // ========================================================================

    /*
        API: POST /api/v1/batch-jobs/{jobCode}/stop
        ฟังก์ชันนี้หยุดงาน Batch ที่กำลังทำงาน (ใช้ในกรณีฉุกเฉิน)
        This function stops a running batch job (used in emergencies).
        Rate Limit: อนุญาต 3 ครั้งต่อ 1 ชั่วโมง
        Rate Limit: Allows 3 requests per 1 hour.
    */
    @PostMapping("/{jobCode}/stop")
    @RateLimit(limit = 3, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Stop a running batch job")
    public ResponseEntity<BatchJobResponseDTO> stopJob(@PathVariable String jobCode)
            throws SystemGlobalException {
        BatchJobResponseDTO response = batchJobService.stopJob(jobCode);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 5. ENABLE/DISABLE JOB (เปิด/ปิดใช้งานงาน)
    // ========================================================================

    /*
        API: PUT /api/v1/batch-jobs/{jobCode}/toggle
        ฟังก์ชันนี้เปิดหรือปิดใช้งานงาน Batch ที่ระบุ
        This function enables or disables a batch job.
        Rate Limit: อนุญาต 10 ครั้งต่อ 5 นาที
        Rate Limit: Allows 10 requests per 5 minutes.
    */
    @PutMapping("/{jobCode}/toggle")
    @RateLimit(limit = 10, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Enable or disable batch job")
    public ResponseEntity<BatchJobResponseDTO> toggleJob(
            @PathVariable String jobCode,
            @RequestParam boolean enabled) throws SystemGlobalException {
        BatchJobResponseDTO response = batchJobService.toggleJob(jobCode, enabled);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 6. GET JOB HISTORY (ประวัติการรันงาน)
    // ========================================================================

    /*
        API: GET /api/v1/batch-jobs/{jobCode}/history
        ฟังก์ชันนี้ดึงประวัติการรันของงาน Batch ที่ระบุ แบบแบ่งหน้า
        This function retrieves execution history of a batch job with pagination.
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @GetMapping("/{jobCode}/history")
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get batch job execution history")
    public ResponseEntity<Page<BatchJobHistoryResponseDTO>> getJobHistory(
            @PathVariable String jobCode,
            @RequestParam(required = false) String status,
            @RequestParam(required = false) String startDate,
            @RequestParam(required = false) String endDate,
            Pageable pageable) throws SystemGlobalException {
        Page<BatchJobHistoryResponseDTO> page = batchJobService.getJobHistory(
                jobCode, status, startDate, endDate, pageable);
        return ResponseEntity.ok(page);
    }

    // ========================================================================
    // 7. GET ALL JOB HISTORIES (ประวัติทั้งหมด)
    // ========================================================================

    /*
        API: GET /api/v1/batch-jobs/history/all
        ฟังก์ชันนี้ดึงประวัติการรันของทุกงาน Batch แบบแบ่งหน้า
        This function retrieves execution history of all batch jobs with pagination.
        Rate Limit: อนุญาต 15 ครั้งต่อ 1 นาที
        Rate Limit: Allows 15 requests per minute.
    */
    @GetMapping("/history/all")
    @RateLimit(limit = 15, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get all batch job histories")
    public ResponseEntity<Page<BatchJobHistoryResponseDTO>> getAllHistories(
            @RequestParam(required = false) String jobCode,
            @RequestParam(required = false) String status,
            @RequestParam(required = false) String startDate,
            @RequestParam(required = false) String endDate,
            Pageable pageable) throws SystemGlobalException {
        Page<BatchJobHistoryResponseDTO> page = batchJobService.getAllHistories(
                jobCode, status, startDate, endDate, pageable);
        return ResponseEntity.ok(page);
    }
}
```

---
## 📁 โครงสร้างโมดูล Batch Jobs (ปรับปรุงใหม่)

```
modules/batch/
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
│   │   ├── BatchJobStatus.java        // SCHEDULED, RUNNING, COMPLETED, FAILED, CANCELLED
│   │   ├── BatchJobType.java          // EMAIL, REPORT, SYNC, CLEANUP, UPDATE, SUMMARY, INVENTORY
│   │   └── BatchJobPriority.java      // LOW, NORMAL, HIGH, CRITICAL
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
│   │   └── Job007LowStockAlert.java      // ⬅️ ใหม่: งานแจ้งเตือนสต็อกต่ำ
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
```

---

## 🗄️ ออกแบบฐานข้อมูล (Database Design) สำหรับโมดูล Batch Jobs

### 📄 SQL DDL (ไฟล์: `src/main/resources/db/migration/V11__batch_schema.sql`)

```sql
-- ==============================================
-- ตาราง: m_batch_job (ข้อมูลงานประจำ)
-- Batch job master data.
-- ==============================================
CREATE TABLE IF NOT EXISTS m_batch_job (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_code VARCHAR(30) UNIQUE NOT NULL,           -- รหัสงาน (batch001, batch002, ...)
    job_name VARCHAR(100) NOT NULL,                 -- ชื่องาน
    job_type VARCHAR(30) NOT NULL,                  -- EMAIL, REPORT, SYNC, CLEANUP, UPDATE, SUMMARY, INVENTORY
    description TEXT,                               -- คำอธิบายงาน
    cron_expression VARCHAR(100) NOT NULL,          -- Cron expression
    enabled BOOLEAN DEFAULT TRUE,                   -- เปิดใช้งานหรือไม่
    max_retry INTEGER DEFAULT 3,                    -- จำนวนครั้งที่ลองใหม่
    retry_delay_ms INTEGER DEFAULT 60000,           -- หน่วงเวลาก่อนลองใหม่ (ms)
    timeout_seconds INTEGER DEFAULT 300,            -- หมดเวลา (วินาที)
    last_run_time TIMESTAMP,                        -- เวลาที่รันล่าสุด
    next_run_time TIMESTAMP,                        -- เวลาที่จะรันครั้งถัดไป
    total_runs INTEGER DEFAULT 0,                   -- จำนวนครั้งที่รันทั้งหมด
    last_status VARCHAR(20),                        -- สถานะล่าสุด
    parameters JSONB,                               -- พารามิเตอร์เพิ่มเติม
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_m_batch_job_code ON m_batch_job(job_code);
CREATE INDEX idx_m_batch_job_enabled ON m_batch_job(enabled);
CREATE INDEX idx_m_batch_job_type ON m_batch_job(job_type);

-- ==============================================
-- ข้อมูลเริ่มต้น: งานประจำทั้ง 7 งาน
-- Initial data: All 7 batch jobs.
-- ==============================================
INSERT INTO m_batch_job (job_code, job_name, job_type, description, cron_expression, enabled, user_id, whitelabel_id)
VALUES 
('batch001', 'ส่งอีเมลแจ้งเตือนรายวัน', 'EMAIL', 'ส่งอีเมลแจ้งเตือนรายวันให้กับพนักงาน', '0 30 6 ? * *', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('batch002', 'สร้างรายงานประจำวัน', 'REPORT', 'สร้างรายงานสรุปประจำวัน', '0 45 6 ? * *', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('batch003', 'อัปเดตสถานะงานค้าง', 'UPDATE', 'ตรวจสอบและอัปเดตสถานะงานที่ค้างนานเกินไป', '0 30 6 ? * *', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('batch004', 'ล้างข้อมูล/ซิงค์ฐานข้อมูล', 'CLEANUP', 'ล้างข้อมูลเก่าและซิงค์ฐานข้อมูล', '0 0 3 ? * *', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('batch005', 'ซิงค์ข้อมูล Realtime', 'SYNC', 'ซิงค์ข้อมูล Realtime กับระบบภายนอก', '0 0/30 * * * ?', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('batch006', 'ส่งสรุปยอดขาย', 'SUMMARY', 'ส่งสรุปยอดขายประจำวัน', '0 30 6 ? * *', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001'),
('batch007', 'แจ้งเตือนสต็อกต่ำ', 'INVENTORY', 'ตรวจสอบสินค้าต่ำสต็อกและส่งแจ้งเตือน', '0 30 6 ? * *', true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001');

-- ==============================================
-- ตาราง: t_batch_job_history (ประวัติการรันงาน)
-- Batch job execution history.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_batch_job_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_code VARCHAR(30) NOT NULL REFERENCES m_batch_job(job_code) ON DELETE CASCADE,
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMP,
    status VARCHAR(20) NOT NULL,                -- SCHEDULED, RUNNING, COMPLETED, FAILED, CANCELLED
    error_message TEXT,
    result_summary TEXT,                        -- สรุปผลการทำงาน
    records_processed INTEGER DEFAULT 0,        -- จำนวนรายการที่ประมวลผล
    duration_ms INTEGER,                        -- ระยะเวลาที่ใช้ (มิลลิวินาที)
    trigger_type VARCHAR(20) DEFAULT 'SCHEDULED', -- SCHEDULED, MANUAL, RETRY
    triggered_by UUID REFERENCES m_user(id),    -- ผู้ที่สั่งรัน (ถ้า manual)
    parameters JSONB,                           -- พารามิเตอร์ที่ใช้
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_batch_history_job ON t_batch_job_history(job_code);
CREATE INDEX idx_t_batch_history_status ON t_batch_job_history(status);
CREATE INDEX idx_t_batch_history_started ON t_batch_job_history(started_at);
CREATE INDEX idx_t_batch_history_whitelabel ON t_batch_job_history(whitelabel_id);
```

---

## 🧠 ระบบ Cache และ Distributed Lock สำหรับ Batch Jobs

### `infrastructure/cache/BatchJobLockCacheService.java`

```java
package com.icmon.module.batch.infrastructure.cache;

import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.stereotype.Service;

import java.time.Duration;
import java.util.UUID;

@Service
public class BatchJobLockCacheService {

    private final RedisTemplate<String, String> redisTemplate;

    public BatchJobLockCacheService(RedisTemplate<String, String> redisTemplate) {
        this.redisTemplate = redisTemplate;
    }

    /*
        ฟังก์ชันนี้ใช้ Distributed Lock เพื่อป้องกันไม่ให้งานเดียวกันทำงานพร้อมกัน (ใช้ Redis)
        This function uses Distributed Lock to prevent the same job from running concurrently (uses Redis).
        Redis Key: batch_lock:{jobCode}
    */
    public boolean acquireLock(String jobCode, String instanceId, Duration timeout) {
        String key = "batch_lock:" + jobCode;
        Boolean success = redisTemplate.opsForValue()
                .setIfAbsent(key, instanceId, timeout);
        return Boolean.TRUE.equals(success);
    }

    /*
        ฟังก์ชันนี้ปลดล็อกเมื่อทำงานเสร็จ
        This function releases the lock when job is done.
    */
    public void releaseLock(String jobCode) {
        String key = "batch_lock:" + jobCode;
        redisTemplate.delete(key);
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่างานถูกล็อกอยู่หรือไม่
        This function checks if the job is locked.
    */
    public boolean isLocked(String jobCode) {
        String key = "batch_lock:" + jobCode;
        return Boolean.TRUE.equals(redisTemplate.hasKey(key));
    }
}
```

### `infrastructure/cache/BatchJobCacheService.java`

```java
package com.icmon.module.batch.infrastructure.cache;

import com.icmon.module.batch.domain.MBatchJob;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.CachePut;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

@Service
public class BatchJobCacheService {

    @Cacheable(value = "batch_jobs", key = "#jobCode")
    public MBatchJob getBatchJob(String jobCode) {
        return null;
    }

    @CachePut(value = "batch_jobs", key = "#job.jobCode")
    public MBatchJob saveBatchJob(MBatchJob job) {
        return job;
    }

    @CacheEvict(value = "batch_jobs", key = "#jobCode")
    public void evictBatchJob(String jobCode) {
        // ลบ Cache / Evict cache.
    }

    @CacheEvict(value = "batch_jobs", allEntries = true)
    public void evictAllBatchJobs() {
        // ลบทุก key / Evict all keys.
    }
}
```

---

## ⏱️ Scheduler Configuration และงานทั้ง 7

### `infrastructure/scheduler/BatchSchedulerConfig.java`

```java
package com.icmon.module.batch.infrastructure.scheduler;

import com.icmon.module.batch.infrastructure.cache.BatchJobLockCacheService;
import com.icmon.module.batch.infrastructure.executor.BatchJobExecutor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.context.annotation.Configuration;
import org.springframework.scheduling.annotation.EnableScheduling;
import org.springframework.scheduling.annotation.Scheduled;

import java.time.Duration;
import java.util.UUID;

@Configuration
@EnableScheduling
public class BatchSchedulerConfig {

    private static final Logger logger = LoggerFactory.getLogger(BatchSchedulerConfig.class);
    private final BatchJobExecutor batchJobExecutor;
    private final BatchJobLockCacheService lockCacheService;

    // ใช้ instanceId เพื่อระบุตัวตนของแต่ละ instance (กรณี Multi-instance)
    private final String instanceId = UUID.randomUUID().toString();

    public BatchSchedulerConfig(BatchJobExecutor batchJobExecutor,
                                BatchJobLockCacheService lockCacheService) {
        this.batchJobExecutor = batchJobExecutor;
        this.lockCacheService = lockCacheService;
    }

    // ========================================================================
    // batch001: ส่งอีเมลแจ้งเตือนรายวัน (06:30 น. ทุกวัน)
    // ========================================================================
    @Scheduled(cron = "0 30 6 ? * *")
    public void executeBatch001() {
        executeWithLock("batch001", () -> batchJobExecutor.executeJob("batch001"));
    }

    // ========================================================================
    // batch002: สร้างรายงานประจำวัน (06:45 น. ทุกวัน)
    // ========================================================================
    @Scheduled(cron = "0 45 6 ? * *")
    public void executeBatch002() {
        executeWithLock("batch002", () -> batchJobExecutor.executeJob("batch002"));
    }

    // ========================================================================
    // batch003: อัปเดตสถานะงานค้าง (06:30 น. ทุกวัน)
    // ========================================================================
    @Scheduled(cron = "0 30 6 ? * *")
    public void executeBatch003() {
        executeWithLock("batch003", () -> batchJobExecutor.executeJob("batch003"));
    }

    // ========================================================================
    // batch004: ล้างข้อมูล/ซิงค์ฐานข้อมูล (03:00 น. ทุกวัน)
    // ========================================================================
    @Scheduled(cron = "0 0 3 ? * *")
    public void executeBatch004() {
        executeWithLock("batch004", () -> batchJobExecutor.executeJob("batch004"));
    }

    // ========================================================================
    // batch005: ซิงค์ข้อมูล Realtime (ทุก 30 นาที)
    // ========================================================================
    @Scheduled(cron = "0 0/30 * * * ?")
    public void executeBatch005() {
        executeWithLock("batch005", () -> batchJobExecutor.executeJob("batch005"));
    }

    // ========================================================================
    // batch006: ส่งสรุปยอดขาย (06:30 น. ทุกวัน)
    // ========================================================================
    @Scheduled(cron = "0 30 6 ? * *")
    public void executeBatch006() {
        executeWithLock("batch006", () -> batchJobExecutor.executeJob("batch006"));
    }

    // ========================================================================
    // batch007: แจ้งเตือนสต็อกต่ำ (06:30 น. ทุกวัน) - ใหม่!
    // ========================================================================
    @Scheduled(cron = "0 30 6 ? * *")
    public void executeBatch007() {
        executeWithLock("batch007", () -> batchJobExecutor.executeJob("batch007"));
    }

    /*
        ฟังก์ชันนี้ใช้ Distributed Lock เพื่อป้องกันการทำงานซ้ำของงานเดียวกัน
        This function uses Distributed Lock to prevent duplicate execution of the same job.
    */
    private void executeWithLock(String jobCode, Runnable job) {
        // 1. พยายามล็อก (หมดเวลา 5 นาที) / Try to acquire lock (timeout 5 minutes).
        if (!lockCacheService.acquireLock(jobCode, instanceId, Duration.ofMinutes(5))) {
            logger.warn("Job {} is already running on another instance. Skipping.", jobCode);
            return;
        }

        try {
            logger.info("Executing job: {} on instance: {}", jobCode, instanceId);
            job.run();
        } catch (Exception e) {
            logger.error("Error executing job {}: {}", jobCode, e.getMessage(), e);
        } finally {
            // 2. ปลดล็อกเมื่อเสร็จ / Release lock when done.
            lockCacheService.releaseLock(jobCode);
            logger.info("Released lock for job: {}", jobCode);
        }
    }
}
```

---

## 🔧 ตัวอย่าง Implementation ของงานที่ 7 (batch007)

### `infrastructure/executor/BatchJobExecutor.java` (เพิ่มส่วนใหม่)

```java
package com.icmon.module.batch.infrastructure.executor;

import com.icmon.module.batch.application.interfaces.BatchJobService;
import com.icmon.module.batch.domain.MBatchJob;
import com.icmon.module.batch.domain.TBatchJobHistory;
import com.icmon.module.batch.domain.enums.BatchJobStatus;
import com.icmon.module.batch.infrastructure.cache.BatchJobCacheService;
import com.icmon.module.inventory.infrastructure.repository.PartMasterRepository;
import com.icmon.module.inventory.infrastructure.repository.InventoryAlertRepository;
import com.icmon.module.inventory.infrastructure.entity.PartMasterEntity;
import com.icmon.module.inventory.infrastructure.entity.InventoryAlertEntity;
import com.icmon.module.email.application.interfaces.EmailService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

@Service
public class BatchJobExecutor {

    private static final Logger logger = LoggerFactory.getLogger(BatchJobExecutor.class);
    private final BatchJobService batchJobService;
    private final BatchJobCacheService cacheService;

    // เพิ่ม Dependency สำหรับ Inventory
    private final PartMasterRepository partMasterRepository;
    private final InventoryAlertRepository alertRepository;
    private final EmailService emailService;

    public BatchJobExecutor(BatchJobService batchJobService,
                            BatchJobCacheService cacheService,
                            PartMasterRepository partMasterRepository,
                            InventoryAlertRepository alertRepository,
                            EmailService emailService) {
        this.batchJobService = batchJobService;
        this.cacheService = cacheService;
        this.partMasterRepository = partMasterRepository;
        this.alertRepository = alertRepository;
        this.emailService = emailService;
    }

    public void executeJob(String jobCode) {
        try {
            MBatchJob job = batchJobService.getJob(jobCode);
            if (job == null || !job.getEnabled()) {
                logger.info("Job {} is disabled or not found. Skipping.", jobCode);
                return;
            }

            if (job.getNextRunTime() != null && job.getNextRunTime().isAfter(LocalDateTime.now())) {
                logger.info("Job {} is scheduled for {}. Skipping.", jobCode, job.getNextRunTime());
                return;
            }

            TBatchJobHistory history = new TBatchJobHistory();
            history.setJobCode(jobCode);
            history.setStartedAt(LocalDateTime.now());
            history.setStatus(BatchJobStatus.RUNNING);
            history.setTriggerType("SCHEDULED");

            TBatchJobHistory savedHistory = batchJobService.saveHistory(history);

            boolean success = executeJobByType(job, savedHistory);

            savedHistory.setFinishedAt(LocalDateTime.now());
            savedHistory.setStatus(success ? BatchJobStatus.COMPLETED : BatchJobStatus.FAILED);
            savedHistory.setDurationMs(
                (int) (savedHistory.getFinishedAt().toInstant().toEpochMilli() 
                - savedHistory.getStartedAt().toInstant().toEpochMilli())
            );
            batchJobService.updateHistory(savedHistory);

            job.setLastRunTime(savedHistory.getStartedAt());
            job.setNextRunTime(calculateNextRunTime(job.getCronExpression()));
            job.setTotalRuns(job.getTotalRuns() + 1);
            job.setLastStatus(savedHistory.getStatus().name());
            batchJobService.updateJob(job);
            cacheService.saveBatchJob(job);

            logger.info("Job {} completed with status: {}", jobCode, savedHistory.getStatus());

        } catch (Exception e) {
            logger.error("Failed to execute job {}: {}", jobCode, e.getMessage(), e);
            batchJobService.recordFailure(jobCode, e.getMessage());
        }
    }

    private boolean executeJobByType(MBatchJob job, TBatchJobHistory history) {
        switch (job.getJobType()) {
            case EMAIL:
                return executeEmailJob(job, history);
            case REPORT:
                return executeReportJob(job, history);
            case UPDATE:
                return executeUpdateJob(job, history);
            case CLEANUP:
                return executeCleanupJob(job, history);
            case SYNC:
                return executeSyncJob(job, history);
            case SUMMARY:
                return executeSummaryJob(job, history);
            case INVENTORY:  // ✅ ใหม่: งานคลังสินค้า
                return executeInventoryJob(job, history);
            default:
                logger.warn("Unknown job type: {}", job.getJobType());
                return false;
        }
    }

    // ========================================================================
    // batch007: INVENTORY - แจ้งเตือนสต็อกต่ำ
    // ========================================================================
    private boolean executeInventoryJob(MBatchJob job, TBatchJobHistory history) {
        logger.info("📦 Executing Inventory Job: {}", job.getJobCode());
        try {
            UUID whitelabelId = getDefaultWhitelabelId();
            LocalDate today = LocalDate.now();
            
            // 1. ค้นหาสินค้าที่ต่ำกว่า Reorder Level
            List<PartMasterEntity> lowStockParts = partMasterRepository.findLowStockParts();
            
            if (lowStockParts.isEmpty()) {
                history.setResultSummary("No low stock items found.");
                history.setRecordsProcessed(0);
                logger.info("✅ [BATCH007] No low stock items found.");
                return true;
            }

            int alertCount = 0;
            int emailSentCount = 0;

            // 2. สร้าง Alert History
            for (PartMasterEntity part : lowStockParts) {
                boolean alreadyAlerted = alertRepository.existsByAlertDateAndPartId(today, part.getId());
                if (!alreadyAlerted) {
                    InventoryAlertEntity alert = new InventoryAlertEntity();
                    alert.setAlertDate(today);
                    alert.setPartId(part.getId());
                    alert.setPartCode(part.getPartCode());
                    alert.setPartName(part.getPartName());
                    alert.setCurrentStock(part.getStockQuantity());
                    alert.setReorderLevel(part.getReorderLevel());
                    alert.setReorderQuantity(part.getReorderQuantity());
                    alert.setAlertSent(false);
                    alert.setResolved(false);
                    alert.setWhitelabelId(whitelabelId);
                    alertRepository.save(alert);
                    alertCount++;
                }
            }

            logger.info("📊 [BATCH007] Found {} low stock parts, created {} new alerts.", 
                        lowStockParts.size(), alertCount);

            // 3. ส่งอีเมลสรุป (ถ้ามี Alert)
            if (alertCount > 0) {
                sendLowStockEmail(lowStockParts);
                emailSentCount = 1;
                
                // อัปเดตว่าแจ้งเตือนแล้ว
                alertRepository.findUnresolvedAlertsByDate(today, whitelabelId)
                        .forEach(alert -> {
                            alert.setAlertSent(true);
                            alert.setAlertSentAt(LocalDateTime.now());
                            alertRepository.save(alert);
                        });
            }

            history.setResultSummary("Found " + lowStockParts.size() + " low stock items, created " + 
                                     alertCount + " alerts, sent " + emailSentCount + " email.");
            history.setRecordsProcessed(lowStockParts.size());
            
            logger.info("✅ [BATCH007] Low stock alert job completed.");
            return true;

        } catch (Exception e) {
            logger.error("❌ [BATCH007] Failed to execute inventory job: {}", e.getMessage(), e);
            history.setResultSummary("Failed: " + e.getMessage());
            return false;
        }
    }

    // ========================================================================
    // ส่งอีเมลแจ้งเตือนสินค้าต่ำสต็อก
    // ========================================================================
    private void sendLowStockEmail(List<PartMasterEntity> lowStockParts) {
        try {
            StringBuilder html = new StringBuilder();
            html.append("<h2>⚠️ แจ้งเตือนสินค้าต่ำสต็อก (Low Stock Alert)</h2>");
            html.append("<p>รายการสินค้าที่ต่ำกว่าเกณฑ์ (Reorder Level):</p>");
            html.append("<table border='1' cellpadding='5' style='border-collapse:collapse;'>");
            html.append("<tr style='background:#4A90E2;color:#fff;'>");
            html.append("<th>รหัส</th><th>ชื่อ</th><th>คงเหลือ</th><th>Reorder Level</th><th>แนะนำสั่งซื้อ</th></tr>");
            
            for (PartMasterEntity part : lowStockParts) {
                String color = part.getStockQuantity() < 5 ? "color:red;font-weight:bold;" : "";
                html.append(String.format(
                    "<tr><td>%s</td><td>%s</td><td style='%s'>%d</td><td>%d</td><td>%d</td></tr>",
                    part.getPartCode(),
                    part.getPartName(),
                    color,
                    part.getStockQuantity(),
                    part.getReorderLevel(),
                    part.getReorderQuantity()
                ));
            }
            html.append("</table>");
            html.append("<p>กรุณาสั่งซื้อสินค้าตามจำนวนที่แนะนำ</p>");
            html.append("<p><small>ส่งจากระบบอัตโนมัติ - ").append(LocalDateTime.now()).append("</small></p>");

            emailService.sendHtmlEmail(
                "purchasing@company.com",
                "[AUTO] แจ้งเตือนสินค้าต่ำสต็อก ประจำวันที่ " + LocalDate.now(),
                html.toString()
            );
            
            logger.info("📧 [BATCH007] Low stock alert email sent successfully.");

        } catch (Exception e) {
            logger.error("❌ [BATCH007] Failed to send low stock alert email: {}", e.getMessage());
        }
    }

    private UUID getDefaultWhitelabelId() {
        return UUID.fromString("00000000-0000-0000-0000-000000000001");
    }

    // ... เมธอดอื่นๆ (executeEmailJob, executeReportJob, etc.) เหมือนเดิม ...
    private LocalDateTime calculateNextRunTime(String cronExpression) {
        try {
            org.springframework.scheduling.support.CronExpression cron = 
                org.springframework.scheduling.support.CronExpression.parse(cronExpression);
            java.time.ZonedDateTime next = cron.next(java.time.ZonedDateTime.now());
            return next != null ? next.toLocalDateTime() : null;
        } catch (Exception e) {
            logger.error("Error parsing cron expression: {}", cronExpression, e);
            return null;
        }
    }
}
```

---
#  สร้าง database  พร้อม Demo data