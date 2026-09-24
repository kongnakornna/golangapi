# SLOW SQL Fix — MQTT write path + index hygiene (iot_data / activity_log / device_status)

วันที่: 2026-08-27 (อัปเดตล่าสุด)
ขอบเขต: วิเคราะห์ query ทั้งโปรเจกต์ พร้อมทดสอบกับ DB จริง (PostgreSQL 15.19)
เอกสารนี้ต่อจาก `README_slow_sql_fix.md` ซึ่งแก้ **device query** ไปแล้ว
เอกสารนี้ครอบคลุม **ทุก slow SQL ที่เหลือ** ในเส้นทางรับข้อมูล MQTT + ปัญหา index ซ้ำซ้อน

ไฟล์ที่เกี่ยวข้อง:
- `internal/modules/iot/usecase/usecase.go` — `ProcessMqttData` (บรรทัด 2104), `UpdateDeviceStatus` (2001), `logActivity` (2556)
- `internal/modules/iot/repository/iot_data_repo.go`
- INSERT / COUNT / DELETE (บรรทัด 36, 72, 81)
- `internal/modules/iot/repository/activity_log_repo.go:23` — INSERT
- `internal/modules/iot/repository/device_status_repo.go` — SELECT (26), Save/Upsert (34)
- `internal/modules/iot/repository/device_repo.go:436` — COUNT ของ ListDevicesWithAlarm
- หลักฐาน: `docs/SLOW_SQL.txt` (GORM slow-log ≥ 200ms)

---

## 1. ภาพรวม: ปัญหาใหญ่ที่สุดคือ "write cost ต่อ 1 MQTT message"

ทุกครั้งที่ device ส่ง MQTT 1 ข้อความ `ProcessMqttData` จะรัน **5 statements DB แยกกัน**
(แต่ละตัว auto-commit แยก, ไม่ batch กัน) ตามหลักฐานใน `docs/SLOW_SQL.txt`:

```text
device_repo.go:436      SELECT COUNT(*) FROM sd_iot_device ... d.device_id = '56'      [325ms]
activity_log_repo.go:23 INSERT INTO "activity_log" (...) VALUES ('DATA_RECEIVED','75',...) [243ms]
iot_data_repo.go:37     INSERT INTO "iot_data" (...) VALUES ('68', ...)                 [284ms]
device_status_repo.go:26 SELECT * FROM "device_status" WHERE "deviceId" = '70' ...     [217ms]
device_status_repo.go:34 UPDATE "device_status" SET "deviceId"='139', ...,"lastData"='...' WHERE "id"=2205 [204ms]
```

ใน log 4 วินาที (18:46:51–55) มี count query รัน **มากกว่า 20 ครั้ง** บน device ต่าง ๆ —
นั่นคือโถ้ง **N+1**: วนลูป query ทีละ device แทนที่จะรวมเป็นชุดเดียว.

ในขณะที่แต่ละ statement รันตัวเดียวแบบ idle จะเร็ว (ดู §3) แต่พอรันพร้อมกัน ~20+ statements/วิ
+ autovacuum + WAL commit ทุกตัวแยก → เกิด contention → แต่ละตัวช้าเป็น 200–800ms.

---

## 2. Root causes (เรียงตามผลกระทบ)

### 2.1 N+1: `ListDevicesWithAlarm` รันต่อ 1 message (usecase.go:2113)
เรียกครั้งละ device เดียว ทั้งที่ข้อมูล (MqttStatusDataName / Measurement / bucket)
Device เปลี่ยนน้อยมาก — ควร cache ไว้ใน memory (map[deviceID]DeviceAlarmListItem) แล้ว refresh เป็นระยะ
แทนการ query. สังเกต log: count query รับ `device_id` ต่างกันเป็นสิบ ๆ ครั้งในไม่กี่วินาที.

### 2.2 `device_status` ถูก UPDATE ทุกคอลัมน์ทุก message (device_status_repo.go:34 = gorm Save)
```sql
UPDATE "device_status" SET "deviceId"='139',"isOnline"=true,"lastSeen"=...,"lastData"='<json>',
  "batteryLevel"=NULL,"signalStrength"=NULL,"temperature"=NULL,...,"customFields"=NULL,
  "createdAt"=...,"updatedAt"=...
WHERE "id" = 2205
```
`Save(status)` เขียน **ทุกฟิลด์** (รวม jsonb ทั้งหมดเป็น NULL) ทุกครั้ง → rewrite ทั้ง row + WAL ใหญ่.
ทางแก้: ใช้ `UPDATE ... SET "lastSeen"=NOW(), "lastData"=$1, "isOnline"=true WHERE "deviceId"=$2`
(เฉพาะคอลัมน์ที่เปลี่ยน) — หรือ `INSERT ... ON CONFLICT ("deviceId") DO UPDATE` แบบระบุคอลัมน์.
อีกจุด: `Save` ต้องผ่าน `GetByDeviceID` ก่อน (SELECT) จึงเป็น 2 statements; UPSERT ตัวเดียวจบ.

### 2.3 Index ซ้ำซ้อน / ตาย (dead index) บนตาราง write-heavy — ภาษีติดทุก INSERT/UPDATE
ตารางมี **คอลัมน์ทั้งแบบ camelCase (TypeORM เดิม) และ snake_case (GORM ใหม่)** อยู่ด้วยกัน.
ตัวโปรแกรมปัจจุบัน (model GORM) เขียนไปที่คอลัมน์ camelCase เท่านั้น แต่ index บางตัวชี้ไป
คอลัมน์ snake_case ที่ไม่เคยมีข้อมูล → เป็น "dead index" ที่เสียพื้นที่ + โตทุกครั้งที่เขียน... (จริง ๆ แล้ว
มันถูก maintain เฉพาะคอลัมน์ที่เขียน; dead index ที่ชี้คอลัมน์ NULL ยังถูก insert ค่า NULL ลงไปทุกครั้ง).

ผลตรวจสอบจริง (COUNT คอลัมน์ที่ App เขียน/ไม่เขียน):

| ตาราง | คอลัมน์ (snake_case) | จำนวน row ที่มีค่า | สรุป |
|---|---|---|---|
| iot_data | `device_id`, `created_at` | **0** | App เขียน `deviceId`/`createdAt` → index snake_case ตาย |
| device_status | `device_id`, `last_seen` | **0** | App เขียน `deviceId`/`lastSeen` → index snake_case ตาย |
| activity_log | `deviceId` (camelCase) | **0** | App เขียน `device_id` → index camelCase ตาย |
| activity_log | `device_id` (snake_case) | **1,590,032** | คอลัมน์จริงที่เขียน แต่ **ไม่มี index** |

index ที่ซ้ำ/ตาย (ขนาดจริงจาก `pg_stat_user_indexes`):

**iot_data (1.58M rows, ~958MB) — 7 indexes:**
| index | คอลัมน์ | สถานะ | ขนาด |
|---|---|---|---|
| `IDX_f53156e53d01d570cda5115755` | ("deviceId","timestamp") | ใช้งาน | 64MB |
| `IDX_456a81db6d3c5b9c6d3d527e51` | ("deviceId","createdAt") | ใช้งาน (count) | 63MB |
| `idx_iot_data_device_id_timestamp` | (device_id,"timestamp") | **ตาย (คอลัมน์ device_id ไม่เคยเขียน)** | 48MB |
| `idx_iot_data_timestamp` | ("timestamp") | **ซ้ำกับ `IDX_9e0f...` (คอลัมน์เดียวกัน)** | 34MB |
| `IDX_9e0f5fe78daae7b8b9a7e844ce` | ("timestamp") | ใช้งาน | 34MB |
| `PK_7f96308951e1b8b75ec7d0b11ce` | (id) | จำเป็น | 34MB |
| `idx_iot_data_device_id_created_at` | (device_id, created_at) | **ตาย** | 16MB |

**activity_log (1.59M rows, ~416MB) — 5 indexes:**
| index | คอลัมน์ | สถานะ | ขนาด |
|---|---|---|---|
| `IDX_8875d63c11f98bb45d158a3116` | (type,"timestamp") | ใช้งาน | 61MB |
| `IDX_fc47cfd7808aac4bda9b4b9b15` | ("userId","timestamp") | **ตาย (`"userId"` ไม่เคยเขียน, คอลัมน์จริงคือ user_id)** | 48MB |
| `IDX_df6dcc35f8fc10ec9bc797e00e` | ("deviceId","timestamp") | **ตาย (`"deviceId"` ไม่เคยเขียน)** | 48MB |
| `PK_067d761e2956b77b14e534fd6f1` | (id) | จำเป็น | 34MB |
| `IDX_13f3cf247b11fa7fa38be36f01` | ("timestamp") | ใช้งาน | 34MB |
| *(ไม่มี)* | (device_id, "timestamp") | **ขาด! — คอลัมน์ device_id มี 1.59M ค่า** | — |

**device_status (143 rows) — 10 indexes บนตารางเล็ก:**
- duplicate จริง: `IDX_55fcb62dbc5432c6c793b9a796` (UNIQUE "deviceId")
  และ `UQ_55fcb62dbc5432c6c793b9a7967` (UNIQUE CONSTRAINT "deviceId") — key เดียวกัน 2 อัน
- dead: `idx_device_status_device_id`, `idx_device_status_is_active`,
  `idx_device_status_is_online`, `idx_device_status_last_seen` (คอลัมน์ snake_case ไม่เคยเขียน)

### 2.4 COUNT บน `iot_data` ช้า (iot_data_repo.go:72 — ใช้ทุกหน้า ListIotData)
```sql
SELECT COUNT(*) FROM iot_data WHERE "deviceId" = '73';
-- EXPLAIN ANALYZE: Index Only Scan "IDX_456a81db6d3c5b9c6d3d527e51"
--   Heap Fetches: 13,372 / 58,566 rows
--   Execution Time: 1.07–2.55 วินาที
```
count ต้องไล่อ่าน **ทุก index entry ของ device นั้น** (~58k row สำหรับ device 73) และยังต้อง
heap fetch เพราะ pages ที่โดน insert ล่าสุดยังไม่ all-visible → ช้า. ยิ่ง device เก็บข้อมูลนาน
(ทุก 10 วิ × หลายเดือน) ยิ่งแย่.

### 2.5 Aggregation ทั้ง table ช้า (สำหรับ dashboard/report ถ้าค้นแบบ group-by)
```sql
SELECT "deviceId", COUNT(*) ... FROM iot_data GROUP BY "deviceId" -- 21 วินาที (Seq Scan 1.58M rows)
```
การ aggregate เหนือ `data` (jsonb) หรือ group-by ทั้งตารางจะช้ามาก — ต้อง scope ให้แคบด้วย
device + ช่วงเวลา (range บน `timestamp` ซึ่งมี index) ก่อนเสมอ.

---

## 3. ผลทดสอบจริง (EXPLAIN ANALYZE / \timing, DB จริง idle)

| Query | เวลา (idle) | หมายเหตุ |
|---|---|---|
| INSERT iot_data 1 row | 2.85 ms | ต้อง maintain 7 indexes |
| INSERT activity_log 1 row | 2.68 ms | ต้อง maintain 5 indexes |
| UPDATE device_status all-fields | 204–311 ms (under live load) | full-row rewrite, เห็นใน SLOW_SQL.txt |
| SELECT COUNT(*) iot_data per device | **1.07–2.55 s** | Index Only Scan + ~13k Heap Fetches |
| SELECT iot_data per device LIMIT 50 | 4.5 ms | OK (ใช้ index) |
| SELECT iot_data latest per device | 0.56 ms | OK |
| SELECT iot_data date-range 100 rows | 64 ms | OK |
| SELECT activity_log by device (recent 50) | 0.58 ms | OK หน้าแรก; deep offset จะแย่ |
| SELECT COUNT(*) device join mqtt (หลังมี index) | 0.13–0.66 ms | แก้แล้วจาก migration ก่อนหน้า |
| GROUP BY deviceId ทั้ง iot_data | **21 s** | Seq Scan ทั้ง 1.58M rows |

ข้อสังเกต: แต่ละ statement idle เร็ว แต่ **ภายใต้ MQTT load จริง** (จาก SLOW_SQL.txt)
ตัวเดียวกันใช้เวลา 243–799ms — ยืนยันว่าเป็นปัญหา contention/volume ไม่ใช่ plan ของ statement เดี่ยว
ยกเว้น COUNT ที่ช้าโดย nature.

---

## 4. แนวทางแก้ไข (implementable)

### 4.1 ลด round-trip ต่อ message (ผลกระทบสูงสุด)
- **Cache device lookup**: เก็บ `map[string]DeviceAlarmListItem` (MqttStatusDataName, Measurement,
  bucket) ไว้ใน unit ของ UseCase แล้ว refresh (เช่น expired 60–300s หรือ trigger แบบ push เมื่อ device เปลี่ยน)
  → ตัด `ListDevicesWithAlarm` ออกจากเส้นทาง hot per message. ถ้าอยากทำแบบมีหลักฐาน ให้แยก
  เป็น `GetDeviceAlarmItemCached(ctx, deviceID)` + fallback query เมื่อ miss.
- **Batch insert `iot_data` + `activity_log`**: สะสม payload ใน buffer (slice/count≤N หรือเวลา≤T)
  แล้ว `INSERT` ทีละชุด (`CreateInBatches` / multi-row INSERT) — ตัด WAL commits หลายครั้งต่อวินาที.
  `activity_log` สำหรับ `DATA_RECEIVED` ต่อ message อาจลดความถี่ (เช่น log ทุก X วินาทีต่อ device) เพราะอ่านยากอยู่แล้ว.
- **UPSERT device_status แบบระบุคอลัมน์**:
  ```sql
  INSERT INTO "device_status" ("deviceId","isOnline","lastSeen","lastData","updatedAt")
  VALUES ($1,true,NOW(),$2,NOW())
  ON CONFLICT ("deviceId") DO UPDATE
    SET "isOnline"=true, "lastSeen"=NOW(), "lastData"=EXCLUDED."lastData", "updatedAt"=NOW();
  ```
  แทน `Save()` ที่เขียนทุกคอลัมน์ → ตัด SELECT ก่อน + เขียนเฉพาะที่เปลี่ยน.
- ใช้ **transaction เดียว** ครอบทั้ง message เพื่อให้ commits ต่อ message เหลือ 1 (หรือ batch).

### 4.2 ลบ index ซ้ำ/ตาย + สร้างที่ขาด (SQL ด้านล่าง) — ลดภาษี write ทุก statement
ลด index ที่ต้อง maintain ต่อ INSERT:
- iot_data: ลบ `idx_iot_data_device_id_timestamp`, `idx_iot_data_timestamp`,
  `idx_iot_data_device_id_created_at` (คืน ~98MB)
- activity_log: ลบ `IDX_fc47cfd7808aac4bda9b4b9b15`, `IDX_df6dcc35f8fc10ec9bc797e00e`
  (คืน ~96MB) และ **สร้าง index บนคอลัมน์จริง** `device_id`
- device_status: ลบ duplicate constraint/index `UQ_55fcb62dbc5432c6c793b9a7967`
  (เก็บ `IDX_55fc...` ไว้ตัวเดียว) + ลบ dead `idx_device_status_*` ทั้ง 4 ตัว

### 4.3 แก้ COUNT ช้า (iot_data_repo CountByDeviceID 1–2.5s)
- **Keyset pagination** แทน offset+count: ลูกค้าส่ง `timestamp`/`id` ของแถวสุดท้าย แล้ว query
  `WHERE "deviceId"=$1 AND timestamp < $2 ORDER BY timestamp DESC LIMIT $3` — ไม่ต้อง count เลย.
- ถ้าจำเป็นต้อง total: แสดง "n+ แถว" แบบประมาณ หรือ maintain counter แยก (trigger/ตาราง summary
  ต่อ device) ดีกว่าการ scan 58k index entries ทุกครั้ง.

### 4.4 Data lifecycle
- `iot_data` โต ~เท่าเดิมทุกครั้งที่ device ส่งทุก 10 วิ (58k rows/dev × ~30 devices). แนะนำ:
  - Partition by `timestamp` (หรือ `deviceId`) รายเดือน
  - scheduler `DeleteOlderThan` (CleanupOldData) รันทุกวัน พร้อม `VACUUM` ตาม
  - ควรเก็บ datapoint ลง InfluxDB (มีอยู่แล้ว) และลดการเก็บ raw ใน `iot_data` ให้เป็น abstract/
    aggregated แทนการ INSERT ทุกวินาที

### 4.5 config / autovacuum
- ตารางเขียนหนัก (`iot_data`, `activity_log`, `device_status`) มี dead tuple สูง:
  `device_status` dead_pct 73% (378/143) ก่อน autovacuum; `activity_log` เคยเห็น n_live_tup
  ต่างจาก COUNT จริง (bloat). ลองตั้งค่า per-table:
  ```sql
  ALTER TABLE device_status SET (autovacuum_vacuum_scale_factor = 0.05, autovacuum_vacuum_threshold = 50);
  ```
- ตรวจสอบ `work_mem` (4MB ปัจจุบัน), `shared_buffers` (128MB) สำหรับ container ขนาดนี้ —
  เพิ่ม work_mem ช่วย range/group ที่ใช้ hash.

### 4.6 query ที่ควรหลีกเลี่ยงใน API
- `GetByDateRange` (iot_data_repo.go:52) ดึง **ทุกแถว** ของช่วงโดยไม่มี LIMIT → เติม LIMIT + pagination.
- Aggregate / group-by ทั้งตาราง → จำกัดด้วย `WHERE "deviceId"=... AND timestamp BETWEEN ...`.

---

## 5. SQL สำหรับ migration ใหม่ (เสนอ)

```sql
-- ============================================================
-- migrations/2026XXXX_cleanup_dead_indexes_write_path.sql
-- ============================================================

-- 1) iot_data: ลบ index ตาย + index ซ้ำ
DROP INDEX IF EXISTS idx_iot_data_device_id_timestamp;
DROP INDEX IF EXISTS idx_iot_data_timestamp;
DROP INDEX IF EXISTS idx_iot_data_device_id_created_at;

-- 2) activity_log: ลบ index ตาย (camelCase ไม่เคยเขียน) + สร้าง index บนคอลัมน์จริง
DROP INDEX IF EXISTS "IDX_fc47cfd7808aac4bda9b4b9b15";  -- ("userId","timestamp")
DROP INDEX IF EXISTS "IDX_df6dcc35f8fc10ec9bc797e00e";  -- ("deviceId","timestamp")
CREATE INDEX IF NOT EXISTS idx_activity_log_device_id_timestamp
    ON public.activity_log (device_id, "timestamp");

-- 3) device_status: ลบ duplicate constraint (เก็บ IDX_55fc ไว้) + ลบ dead index
ALTER TABLE public.device_status DROP CONSTRAINT IF EXISTS "UQ_55fcb62dbc5432c6c793b9a7967";
DROP INDEX IF EXISTS idx_device_status_device_id;
DROP INDEX IF EXISTS idx_device_status_is_active;
DROP INDEX IF EXISTS idx_device_status_is_online;
DROP INDEX IF EXISTS idx_device_status_last_seen;
```

> ระวัง: ตรวจสอบอีกครั้งก่อน drop ว่า **ไม่มี query/code ใด** อ้างคอลัมน์ snake_case เหล่านั้น
> (`grep -r "device_id" internal/modules/iot` พบเพียง model GORM ซึ่ง mapping ไป camelCase แล้ว;
> ถ้ามี raw SQL อ้าง snake_case ให้รักษา index ที่จำเป็นไว้).

---

## 6. วิธีตรวจสอบ / ยืนยันหลัง deploy

1. เปิด GORM slow-log (< 200ms) — ตัว current logger อยู่แล้ว (เห็นใน SLOW_SQL.txt).
2. ดู `pg_stat_activity` ขณะ MQTT วิ่งเต็ม format เพื่อเช็ค contention/lock wait:
   ```sql
   SELECT pid, state, wait_event_type, wait_event, now()-query_start AS dur, left(query,90)
   FROM pg_stat_activity WHERE state <> 'idle' ORDER BY 5 DESC;
   ```
3. ตรวจ dead tuples + bloat หลัง deploy 2-3 วัน:
   ```sql
   SELECT relname, n_live_tup, n_dead_tup, round(100.0*n_dead_tup/NULLIF(n_live_tup+n_dead_tup,0),1) AS dead_pct,
          last_autovacuum FROM pg_stat_user_tables
   WHERE relname IN ('iot_data','activity_log','device_status') ORDER BY n_dead_tup DESC;
   ```
4. เปรียบเทียบ throughput: นับ messages ต่อวินาทีที่ ingester รับได้ก่อน/หลัง
   (หรือดู `sum(rows)` ต่อวินาทีใน `pg_stat_user_tables.n_tup_ins`).

---

## 7. สรุป prioritisation

| ลำดับ | การแก้ไข | ผลกระทบ | ความเสี่ยง |
|---|---|---|---|
| 1 | Cache device lookup (ตัด N+1 per message) | ลด statements ต่อ message ทันที | ต่ำ (fallback query ได้) |
| 2 | UPSERT device_status เฉพาะคอลัมน์ | ลดการ rewrite row + WAL | ต่ำ |
| 3 | Batch insert iot_data/activity_log | ลด commit contention | กลาง (ต้องจัดการ delay/retry) |
| 4 | ลบ index ซ้ำ/ตาย + สร้างที่ขาด | ภาษา write cost ทุก statement, คืนพื้นที่ | ต่ำ-กลาง (ตรวจการอ้างอิงก่อน) |
| 5 | Keyset pagination แทน COUNT | กำจัด COUNT 1–2.5s | กลาง (ต้องเปลี่ยน API/client) |
| 6 | Partition + retention + autovacuum | สกัดการโตไม่จำกัด | กลาง |