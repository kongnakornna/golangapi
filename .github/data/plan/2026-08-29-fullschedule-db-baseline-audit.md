# แผน: ปรับโครงสร้าง Database โมดูล fullschedule — Baseline Schema + Audit + Constraint

## ภาพรวมการแก้ไข

| # | Layer | ไฟล์ | สิ่งที่แก้ |
|---|-------|------|-----------|
| 1 | migration (ใหม่) | `migrations/20260833_fullschedule_baseline_audit.sql` | รวม schema `fs_*` ทั้งหมดเป็นไฟล์ baseline เดียว (idempotent) + เพิ่ม audit columns / FK / unique / check constraint / index |
| 2 | model | `internal/modules/fullschedule/model.go` | เพิ่ม field audit (`UpdatedBy`, `Version`, ...) ตรงกับคอลัมน์ใหม่ใน `Schedule` / `Group` / `Zone` / `Area` |
| 3 | presenter | `internal/modules/fullschedule/presenter/presenters.go` | เพิ่ม field audit ใน `ScheduleResponse` / `GroupResponse` / ... (ถ้าต้องส่งออก) |
| 4 | handlers | `internal/modules/fullschedule/delivery/http/handlers.go`, `master_handlers.go` | mapping audit fields ใน create/update/response |
| 5 | Go build/vet/test | — | ตรวจสอบ compile + vet + test ผ่านทั้งหมด |

---

## เหตุผล / เป้าหมาย

- ปัจจุบัน schema `fs_*` ถูกกระจายอยู่ในหลายไฟล์ migration ที่รันด้วยมือ (`20260829`, `20260830`, `20260831`, `20260832`) ทำให้ไม่เห็น "โครงสร้างเดียวจบ" และยังขาด
  - audit columns: `updated_by`, `version` (optimistic locking), `deleted_at` ที่ยังไม่ครบทุกตาราง
  - check constraint บังคับค่าที่ถูกต้อง (`status`, `mode`, `event`, `event_action`, วันในสัปดาห์ `0/1`)
  - `NOT NULL`/default ที่สอดคล้องกับประเภทข้อมูลใหม่ (status int8 1/0/3)
- เป้าหมายคือสร้างไฟล์ **baseline** ไฟล์เดียวที่อัปเดต schema จริงแบบ idempotent (รันซ้ำได้ ไม่ทำลายข้อมูลเดิม) และใช้เป็น "source of truth" ของโครงสร้างตาราง `fs_*`

---

## 1) migration ใหม่: `migrations/20260833_fullschedule_baseline_audit.sql`

**Location:** สร้างไฟล์ใหม่ `migrations/20260833_fullschedule_baseline_audit.sql`

**From → To:** เปลี่ยนจากการ "แบ่ง schema เป็นหลายไฟล์" → เป็น "ไฟล์ baseline เดียวที่อัปเดตโครงสร้างครบทุกตาราง fs_* โดยรันซ้ำได้ (idempotent)" ใช้หลักการ `CREATE TABLE IF NOT EXISTS`, `ALTER TABLE ... ADD COLUMN IF NOT EXISTS`, `CREATE INDEX IF NOT EXISTS`, `DO $$` สำหรับ check constraint (ป้องกัน error เมื่อ constraint มีอยู่แล้ว)

**Change detail — แบ่งเป็น 6 ตาราง ดังนี้**

### 1.1 `fs_schedule` — คอลัมน์ + constraint + index

`ALTER TABLE fs_schedule ADD COLUMN IF NOT EXISTS ...`:
- `updated_by UUID` (audit)
- `version BIGINT NOT NULL DEFAULT 1` (optimistic locking)
- ตั้งค่า default คอลัมน์เดิมที่ยังไม่มี cast: (เนื่องจากการแก้ไปก่อนหน้านี้ `status` เป็น int8 แล้ว ไม่ต้องแตะ type)

Check constraints (เพิ่มด้วย `DO $$ ... IF NOT EXISTS` เพื่อกัน error):
```sql
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_fs_schedule_status') THEN
        ALTER TABLE fs_schedule ADD CONSTRAINT chk_fs_schedule_status
            CHECK (status IN (0, 1, 3));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_fs_schedule_mode') THEN
        ALTER TABLE fs_schedule ADD CONSTRAINT chk_fs_schedule_mode
            CHECK (mode IN ('normal', 'full', 'batch'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_fs_schedule_event') THEN
        ALTER TABLE fs_schedule ADD CONSTRAINT chk_fs_schedule_event
            CHECK (event IS NULL OR event IN (0, 1, 2, 3));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_fs_schedule_event_action') THEN
        ALTER TABLE fs_schedule ADD CONSTRAINT chk_fs_schedule_event_action
            CHECK (event_action IN ('ON', 'OFF', 'start', 'stop'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_fs_schedule_weekday') THEN
        ALTER TABLE fs_schedule ADD CONSTRAINT chk_fs_schedule_weekday
            CHECK (sunday IN (0,1) AND monday IN (0,1) AND tuesday IN (0,1)
                AND wednesday IN (0,1) AND thursday IN (0,1) AND friday IN (0,1)
                AND saturday IN (0,1));
    END IF;
END $$;
```

Indexes (สร้างใหม่เฉพาะที่ยังไม่มี — ที่มีอยู่แล้วจาก migration เก่าคงไว้):
```sql
CREATE INDEX IF NOT EXISTS idx_fs_schedule_due
    ON fs_schedule (status, next_run_at) WHERE status = 1 AND deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_fs_schedule_mode_status ON fs_schedule (mode, status);
CREATE INDEX IF NOT EXISTS idx_fs_schedule_group_id ON fs_schedule (group_id);
CREATE INDEX IF NOT EXISTS idx_fs_schedule_zone_id  ON fs_schedule (zone_id);
CREATE INDEX IF NOT EXISTS idx_fs_schedule_area_id  ON fs_schedule (area_id);
CREATE UNIQUE INDEX IF NOT EXISTS uq_fs_schedule_legacy_id ON fs_schedule (legacy_id) WHERE legacy_id IS NOT NULL;
```

### 1.2 `fs_schedule_device` — audit + FK

`ALTER TABLE fs_schedule_device ADD COLUMN IF NOT EXISTS`:
- `updated_at TIMESTAMPTZ NOT NULL DEFAULT now()`
- `deleted_at TIMESTAMPTZ`

(มี FK → `fs_schedule(id) ON DELETE CASCADE` + UNIQUE `(schedule_id, device_id)` แล้วจาก migration เก่า — คงไว้ ไม่ต้องเพิ่มซ้ำ)

### 1.3 `fs_schedule_history` — คงโครงสร้างเดิม

- คงเดิมทั้งหมด (ไม่เพิ่ม audit เพราะ history เป็น log แบบ append-only)
- เติม index partial สำหรับ query worker (ถ้ายังไม่มี):
```sql
CREATE INDEX IF NOT EXISTS idx_fs_schedule_history_schedule_created
    ON fs_schedule_history (schedule_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_fs_schedule_history_status_created
    ON fs_schedule_history (status, created_at DESC);
```

### 1.4 `fs_schedule_settings` — คงโครงสร้างเดิม

- คงเดิม (มี UNIQUE `(schedule_id, key)` แล้ว)
- เพิ่ม index ถ้ายังไม่มี:
```sql
CREATE INDEX IF NOT EXISTS idx_fs_schedule_settings_schedule ON fs_schedule_settings (schedule_id);
```

### 1.5 `fs_groups` / `fs_zones` / `fs_areas` — audit

`ALTER TABLE ... ADD COLUMN IF NOT EXISTS` ทั้ง 3 ตาราง:
- `created_by UUID`
- `updated_by UUID`
- `version BIGINT NOT NULL DEFAULT 1`

(มี `created_at` / `updated_at` / `deleted_at` แล้ว)

Check constraint (ถ้ายังไม่มี) สำหรับ `status VARCHAR` ของ group/zone/area:
```sql
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_fs_groups_status') THEN
        ALTER TABLE fs_groups ADD CONSTRAINT chk_fs_groups_status CHECK (status IN ('active', 'inactive'));
    END IF;
END $$;
```

### 1.6 `fs_device_area` — audit

`ALTER TABLE fs_device_area ADD COLUMN IF NOT EXISTS`:
- `updated_at TIMESTAMPTZ NOT NULL DEFAULT now()`
- `deleted_at TIMESTAMPTZ`

(มี UNIQUE `(area_id, device_id)` แล้ว)

---

## 2) model: `internal/modules/fullschedule/model.go`

**Location:** struct `Schedule` (~บรรทัด `Status` ถึง `DeletedAt`), struct `Group` / `Zone` / `Area`

**From → To:** ยังไม่มี field audit → เพิ่ม field ใหม่ตรงกับคอลัมน์ที่เพิ่มใน migration

**Change detail:**

`Schedule` (ตาราง `fs_schedule`) เพิ่ม:
```go
UpdatedBy *uuid.UUID `gorm:"column:updated_by;type:uuid" json:"updated_by,omitempty"`
Version   int64      `gorm:"column:version;not null;default:1" json:"version"`
```

`Group` (fs_groups), `Zone` (fs_zones), `Area` (fs_areas) เพิ่ม (ทั้ง 3 ตัวเหมือนกัน):
```go
CreatedBy *uuid.UUID `gorm:"column:created_by;type:uuid" json:"created_by,omitempty"`
UpdatedBy *uuid.UUID `gorm:"column:updated_by;type:uuid" json:"updated_by,omitempty"`
Version   int64      `gorm:"column:version;not null;default:1" json:"version"`
```

---

## 3) presenter: `internal/modules/fullschedule/presenter/presenters.go`

**Location:** struct `ScheduleResponse` (~บรรทัด `CreatedBy` / `UpdatedAt`), และ struct ที่ map Group/Zone/Area ถ้ามี

**From → To:** ยังไม่มี → เพิ่ม field audit ใน response (เพื่อแสดงใน API)

**Change detail:**

`ScheduleResponse` เพิ่ม:
```go
UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
Version   int64      `json:"version"`
```

---

## 4) handlers: mapping audit fields

**Location:**
- `internal/modules/fullschedule/delivery/http/handlers.go` → `mapScheduleResponse` (~บรรทัด `UpdatedAt`)
- `master_handlers.go` → response mapping ของ group/zone/area (ถ้ามี)

**From → To:** mapping ปัจจุบันไม่มี audit → เพิ่มการ map `UpdatedBy` / `Version` จาก model เข้า presenter

**Change detail:**

ใน `mapScheduleResponse` เพิ่ม:
```go
UpdatedBy: s.UpdatedBy,
Version:   s.Version,
```
เช่นเดียวกับใน master response mapping (ถ้ามี field นี้ใน presenter นั้น)

> หมายเหตุ: เวลาสร้าง/แก้ ควร set `version` (และ `updated_by` จาก token user ถ้ามี context) — แต่ถ้ายังไม่มี user context ชัดเจน ให้อย่างน้อยตั้ง `version = existing.Version + 1` ในการ update (optimistic locking) และ `version = 1` ในการ create ผ่าน default ใน DB ก็พอ

---

## Unit Tests

**ไม่เขียน** — เพราะงานนี้เป็นการแก้ schema SQL (migration) และ field mapping ใน struct/model ที่ไม่มี test infra สำหรับ SQL ใน repo นี้ (ทุก package ปัจจุบันเป็น `[no test files]`) และการตรวจสอบหลักคือ `go build ./...`, `go vet ./...`, `go test ./...` ที่ต้องผ่าน

---

## ลำดับการ implement

1. สร้างไฟล์ migration ใหม่ (ข้อ 1) — baseline audit
2. แก้ model.go (ข้อ 2)
3. แก้ presenters.go (ข้อ 3)
4. แก้ handlers.go / master_handlers.go (ข้อ 4)
5. รัน `go build ./...`, `go vet ./...`, `go test ./...` — ตรวจผ่านทั้งหมด
