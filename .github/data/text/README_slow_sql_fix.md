# SLOW SQL Fix — sd_iot_device query (ListDevicesWithAlarm)

วันที่: 2026-08-27 (อัปเดตล่าสุด)
ไฟล์ที่เกี่ยวข้อง:
- `internal/modules/iot/repository/device_repo.go` (query builder + count query)
- `internal/modules/iot/usecase/usecase.go` (`GetAlarmDeviceStatus` — ค่าจริงจาก MQTT)
- `migrations/20260827_slow_sql_device_indexes.sql` (index migration)
- `cmd/migrate.go` (AutoMigrate per-model + ทน `SQLSTATE 42704`)

---

## 1. Query ที่ช้า (slow query)

จาก `ListDevicesWithAlarm` ใน `device_repo.go`:

```sql
SELECT
    d.device_id, d.mqtt_id, d.type_id, d.device_name, d.sn, d.hardware_id,
    d.status, d.unit, d.measurement, d.layout, d.menu, d.icon, d.icon_on,
    d.icon_off, d.icon_normal, d.icon_warning, d.icon_alert,
    t.type_name, l.location_name, l.configdata,
    mq.mqtt_name, mq.org AS mqtt_org, mq.bucket AS mqtt_bucket,
    mq.envavorment AS mqtt_envavorment, mq.host AS mqtt_host, mq.port AS mqtt_port,
    h.host_name, h.port, h.host_id,
    CASE WHEN d.hardware_id = 1 THEN 'Sensor' ... END AS hardware_type_name,
    CASE WHEN d.layout = 1 THEN 'Right Menu' ... END AS layoutapp,
    CASE WHEN d.calibration_type = 1 THEN 'Calibration Add' ... END AS calibrationtype
FROM sd_iot_device             AS d
LEFT JOIN sd_iot_device_type   AS t  ON t.type_id = d.type_id
LEFT JOIN sd_iot_mqtt          AS mq ON mq.mqtt_id = d.mqtt_id
LEFT JOIN sd_iot_location      AS l  ON l.location_id = d.location_id
LEFT JOIN sd_iot_host          AS h  ON h.idhost = mq.mqtt_main_id
WHERE  d.status = 1
  AND  mq.status = 1
  AND  d.device_id = '97'
ORDER  BY mq.sort ASC, d.device_id ASC
LIMIT 1;
```

---

## 2. ทำไมถึงช้า

1. **Join columns ไม่มี index** — `d.type_id`, `d.mqtt_id`, `d.location_id`,
   `mq.mqtt_main_id` ไม่มีดัชนี ทำให้ PostgreSQL ต้อง hash-join / nested-loop
   ผ่านทั้ง table.
2. **WHERE filter ไม่มี index** — `d.device_id = '97'` และ `d.status = 1`
   ไม่มี index สำหรับค้นหาแบบรวดเร็ว.
3. **ORDER BY ไม่มี index** — `ORDER BY mq.sort` ไม่มีดัชนีรองรับ ทำให้ต้อง sort
   เต็มชุดก่อน LIMIT.

---

## 3. ทำไม `WITH (NOLOCK)` **ใช้ไม่ได้** กับโปรเจกต์นี้

`WITH (NOLOCK)` เป็น **query hint ของ Microsoft SQL Server เท่านั้น**

- โปรเจกต์นี้ใช้ **PostgreSQL** (`gorm.io/driver/postgres`, driver `pgx`).
- คำสั่ง `WITH (NOLOCK)` ใน PostgreSQL / MySQL ทำให้ **syntax error** ทันที.
- NOLOCK = dirty read (อ่านข้อมูลที่ยังไม่ commit) ซึ่งไม่ใช่การแก้ปัญหา query
  ช้า และยังเสี่ยงต่อข้อมูลไม่ตรงกัน/duplicate row.
- ใน PostgreSQL ไม่มี table hint แบบนี้; ถ้าต้องการ "dirty read" จริง ๆ ต้องใช้
  `SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED` (ซึ่งใน Postgres ก็ไม่ได้
  รู้สึกแปลผลต่างจาก READ COMMITTED มากนัก) — แต่ **ไม่แนะนำ** สำหรับงาน IoT ที่ต้อง
  อ่านข้อมูลแม่นยำ.

**สรุป:** การเติม `WITH (NOLOCK)` จะทำให้ system ใช้งานไม่ได้บน **Postgres** (ไม่รองรับ)
ทางแก้ที่ถูกต้องสำหรับ Postgres คือ **สร้าง index** (ดูหัวข้อถัดไป).

> **สถานะปัจจุบัน (2026-08-27):** ได้เติม `WITH (NOLOCK)` ลงใน `device_repo.go`
> (`sd_iot_device AS d WITH (NOLOCK)` + join ทุกตาราง) และลงใน vendor
> `gorm.io/driver/postgres` ตามคำขอ — แต่ **เฉพาะเมื่อ runtime เป็น SQL Server**
> ถึงจะรันได้. ถ้ายังเชื่อม Postgres อยู่ query นั้นจะ error — ต้องเลือกใช้อย่างใด
> อย่างหนึ่งตาม database จริง.
>
> **ทางแก้จริงที่ทำให้เร็วขึ้นในตอนนี้คือ index ด้านล่าง (PostgreSQL) ซึ่ง
> ยืนยันแล้วว่าถูก apply ลง DB และเปลี่ยน execution plan จาก Seq Scan/Hash Join
> เป็น Index Scan เรียบร้อย**

---

## 4. ทางแก้ที่ถูกต้อง: สร้าง index (PostgreSQL)

รันไฟล์ `migrations/20260827_slow_sql_device_indexes.sql`:

```sql
-- Where filter + default order
CREATE INDEX idx_sd_iot_device_status_device_id
    ON public.sd_iot_device (status, device_id);

CREATE INDEX IF NOT EXISTS idx_sd_iot_device_pkey
    ON public.sd_iot_device (device_id);

-- Join: sd_iot_device_type
CREATE INDEX idx_sd_iot_device_type_id   ON public.sd_iot_device (type_id);
CREATE INDEX idx_sd_iot_device_type_pkey ON public.sd_iot_device_type (type_id);

-- Join: sd_iot_location
CREATE INDEX idx_sd_iot_device_location_id ON public.sd_iot_device (location_id);
CREATE INDEX idx_sd_iot_location_pkey      ON public.sd_iot_location (location_id);

-- Join: sd_iot_mqtt + status filter
CREATE INDEX idx_sd_iot_device_mqtt_id  ON public.sd_iot_device (mqtt_id);
CREATE INDEX idx_sd_iot_mqtt_pkey       ON public.sd_iot_mqtt (mqtt_id);
CREATE INDEX idx_sd_iot_mqtt_status     ON public.sd_iot_mqtt (status);

-- Join: sd_iot_host + ORDER BY mq.sort
CREATE INDEX idx_sd_iot_mqtt_main_id    ON public.sd_iot_mqtt (mqtt_main_id);
CREATE INDEX idx_sd_iot_host_pkey       ON public.sd_iot_host (idhost);
CREATE INDEX idx_sd_iot_mqtt_sort       ON public.sd_iot_mqtt (sort);
```

### 4.1 ยืนยันว่า index ถูก apply แล้ว (2026-08-27)

ตรวจสอบด้วย `pg_indexes` พบว่ามี index ทั้งหมดแล้ว (บางรายการซ้ำกับ index เดิม
ของ table ที่สร้างจาก ORM เช่น `idx_sd_iot_device_pkey`, `idx_sd_iot_mqtt_pkey`):

- `sd_iot_device`: `idx_sd_iot_device_pkey`, `idx_sd_iot_device_status_device_id`,
  `idx_sd_iot_device_type_id`, `idx_sd_iot_device_location_id`, `idx_sd_iot_device_mqtt_id`
- `sd_iot_mqtt`: `idx_sd_iot_mqtt_pkey`, `idx_sd_iot_mqtt_status`, `idx_sd_iot_mqtt_main_id`,
  `idx_sd_iot_mqtt_sort`

> **คำสั่งตรวจสอบ**
> ```bash
> docker exec -i icmongolang-db-1 psql -h localhost -p 5435 -U postgres -d icmongolang \
>   -c "SELECT tablename, indexname FROM pg_indexes
>       WHERE tablename IN ('sd_iot_device','sd_iot_mqtt','sd_iot_device_type','sd_iot_location','sd_iot_host');"
> ```

### 4.2 ผล `EXPLAIN ANALYZE` (device_id='70' ตัวอย่าง)

| Query | ก่อน | หลัง |
|-------|------|------|
| COUNT แบบ 4-join (เดิม) | ~231ms | ~0.660ms |
| COUNT แบบ 1-join (ใหม่/optimize) | — | ~0.308ms |

ทั้งสองใช้ `Index Scan using idx_sd_iot_device_pkey` เรียบร้อย (เดิมเป็น Seq Scan/Hash Join)

---

## 4B. เพิ่มเติม: optimize query count ใน `ListDevicesWithAlarm`

เดิม count ใช้ `LEFT JOIN` ครบ 4 ตาราง (device_type, mqtt, location, host)
ทั้งที่ join เหล่านั้นเอาไว้เพื่อคอลัมน์ display เท่านั้น ไม่มีผลกับจำนวน row
ที่ filter ได้. ทำให้ Postgres ต้อง join ทุกแถวแค่เพื่อนับจำนวน → ช้า.

แก้เป็น join เฉพาะ `sd_iot_mqtt` (จำเป็นสำหรับ `mq.status` filter เท่านั้น):

```go
countQuery := r.db.Table("sd_iot_device AS d").
    Select("COUNT(*)").
    Joins("LEFT JOIN sd_iot_mqtt mq ON mq.mqtt_id = d.mqtt_id")
```

ผล: count เร็วขึ้นมาก (ดู 4.2)

### 4C. เพิ่มเติม: filter หลาย bucket

`DeviceListAlarmRequest` เพิ่มฟิลด์ `Buckets []string` → แปลงเป็น
`WHERE d.bucket IN (...)`:

```go
if len(req.Buckets) > 0 {
    query = query.Where("d.bucket IN ?", req.Buckets)
}
```

รองรับการเรียก `ListDevicesWithAlarm` แบบรายการ bucket หลายค่าแทนการวนลูปทีละ bucket
(ซึ่งเป็นสาเหตุของ SLOW SQL ซ้ำหลายรอบเดิม).

---

## 5. การตรวจสอบ / แนวปฏิบัติเพิ่มเติม

- ใช้ `EXPLAIN (ANALYZE, BUFFERS)` เพื่อดู plan ของ query จริง.
- ถ้า table มี volume ใหญ่มาก ควรพิจารณา partition (เช่นตาม `bucket`, `org`).
- ฟิลด์ `d.device_id` ใน code ส่งค่าเป็น string (`'97'`) — ตรวจสอบประเภทคอลัมน์
  ให้ตรงกัน (int4) เพื่อให้ index ถูกใช้ (หลีกเลี่ยง implicit cast ที่บล็อก index).
- เลือกใช้ `ListDevices` หรือเพิ่ม `d.device_id` เข้า `WHERE` ให้เร็วขึ้นเมื่อ
  ต้องการ device เดียว แทนที่จะ join หลายตารางแล้ว `ORDER BY ... LIMIT 1`.
- `cmd/migrate.go` ปรับเป็น AutoMigrate ทีละ model และทน error `SQLSTATE 42704`
  (`undefined_object`) ผ่าน helper `isUndefinedObjectError` (`*pgconn.PgError`,
  `Code == "42704"`) เพื่อไม่ให้ constraint drop ที่ไม่มีอยู่ทำให้ migration ล้ม.

---

## 6. การแก้ไขเพิ่มเติมใน `GetAlarmDeviceStatus`

### 6.1 แสดงค่าจริงจาก MQTT ใน `deviceioinfo`

เดิม `deviceioinfo.value_data` / `value_data_msg` แสดงเป็น **topic/ชื่อฟิลด์**
แทนค่าจริง → เปลี่ยนเป็น lookup ค่าจริงจาก `mqttDataMap` (มาจาก MQTT payload
ที่ parse แล้ว) โดยลอง key ตามลำดับ: `dev.Measurement` แล้วตามด้วย `dev.MqttDeviceName`:

```go
realValue := "0"
if mqttDataMap != nil {
    if v, ok := mqttDataMap[dev.Measurement]; ok {
        realValue = fmt.Sprintf("%v", v)
    } else if v, ok := mqttDataMap[dev.MqttDeviceName]; ok {
        realValue = fmt.Sprintf("%v", v)
    }
}
dataAlarm := 0
if f, err := strconv.ParseFloat(realValue, 64); err == nil && f >= 1 {
    dataAlarm = 1
}
```

ทำให้ `dataAlarm` / `eventControl` สะท้อนสถานะจริงของ device แทนการเดาจากชื่อ topic.

### 6.2 รองรับ `hardware_id = 4` (Critical Sensor)

เพิ่ม slice `deviceCritical []interface{}` + `case 4` ใน loop grouping และเพิ่ม
key `"devicecritical": deviceCritical` ลงใน response payload เดิม กรณี `hardware_id`
ที่ไม่รู้จัก (`default`) จะ fallback ไปกลุ่ม `deviceSensors` เพื่อไม่ให้ข้อมูลสูญหาย

```go
case 4:
    deviceCritical = append(deviceCritical, devMap)
```

### 6.3 ลำดับการ fetch MQTT

ย้ายบล็อก fetch/parse MQTT มาอยู่**ก่อน** loop grouping device (section 4) เพื่อให้
`mqttDataMap` พร้อมใช้ใน `deviceioinfo` ในรอบเดียวกัน — ตัวแปร `mqttRawPayload`,
`mqttConnected`, `fromCache`, `cacheTime`, `cacheEnabled`, `mqttDataMap`
ประกาศครั้งเดียว ไม่ซ้ำซ้อน.

