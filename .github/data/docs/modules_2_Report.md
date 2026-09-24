 ---

## สร้าง `module/job` (modules_2.md) เสร็จสมบูรณ์

โครงสร้างที่สร้างตาม spec ใช้ package `com.icmon.module.job`:

| Layer | ไฟล์ที่สร้าง |
|---|---|
| **domain/** | `TJob`, `TJobService`, `TJobPartSales`, `TJobServiceCarSymptom`, `TJobDiagTroubleCode`, `TJobStatusHistory` |
| **domain/enums/** | `JobStatus` (14 สถานะ) |
| **domain/valueobjects/** | `JobNumber`, `Mileage` |
| **infrastructure/entity/** | 6 entities (`JobEntity`, `JobServiceEntity`, `JobPartSalesEntity`, `JobSymptomEntity`, `JobDiagCodeEntity`, `JobStatusHistoryEntity`) |
| **infrastructure/repository/** | 6 JPA repositories + 3 impl wrappers |
| **infrastructure/mapper/** | `JobMapper`, `JobServiceMapper`, `JobPartSalesMapper` |
| **infrastructure/cache/** | `JobCacheService` (`@Cacheable`/`@CachePut`/`@CacheEvict`), `JobStatusCacheService` |
| **application/interfaces/** | `JobService`, `JobStatusService`, `JobPartSaleService` |
| **application/impl/** | `JobServiceImpl`, `JobStatusServiceImpl`, `JobPartSaleServiceImpl` |
| **application/usecase/** | 8 use cases |
| **presentation/controller/** | `JobController` (10 endpoints + Rate Limit), `JobServiceController`, `JobPartController` |
| **presentation/dto/** | 5 request DTOs + 5 response DTOs |
| **presentation/validator/** | `JobValidator` |
| **db/migration/** | V2__job_card_schema.sql (6 ตาราง + Auto-generate job_no Trigger) |

Made changes.