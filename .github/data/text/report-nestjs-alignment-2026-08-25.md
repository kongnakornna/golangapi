# รายงานผลการดำเนินงาน: จัดแนว Go port กับ NestJS (โมดูล settings)

**วันที่:** อังคารที่ 25 สิงหาคม 2026
**เวลารายงาน:** 10:13 น.
**โปรเจกต์:** icmongolang (`C:\github\icmongolang`)
**ระบบอ้างอิง:** gistdaapi (`C:\github\gistdaapi`) — NestJS, ไฟล์หลัก `src/modules/settings/settings.controller.ts` และ `settings.service.ts`

---

## ขอบเขตงาน

จัดแนว endpoint ของโมดูล `settings` ฝั่ง Go ให้สอดคล้องกับพฤติกรรมจริงของระบบอ้างอิง NestJS ครอบคลุม 5 หมวด:

1. Device list variants
2. Alarm-device paginate family + alarm/process log params
3. Schedule process log + MQTT error log paginate
4. Dashboard config upsert semantics
5. Master create/update dup-check rules

**นโยบายที่ตกลงใช้:** ส่วนที่ NestJS เสีย (อ้าง alias/คอลัมน์ที่ไม่มีอยู่จน SQL error) จะไม่ถูกทำซ้ำใน Go แต่จะรักษา "เจตนา" ของ filter/sort นั้นผ่านทางเลือกที่ทำงานได้จริง

---

## 1. Device list variants — เสร็จสิ้น

จัดแนว `listdevicepage`, `listdeviceactive`, `listdeviceactive1`, `listdevicepagesensor`, `deviceall*`, `allactiveschedule` ให้ตรงกับ NestJS: selects, joins, filters, sort whitelist/default sort

## 2. Alarm-device paginate + alarm/process log params — เสร็จสิ้น

- ตระกูล alarm-device paginate ทั้งหมด
- Log variants (email / line / sms / telegram / control) ใช้ helper `alarmLogVariantSpec(table)` ร่วมกัน:
  - INNER joins เหมือน NestJS
  - default sort `al.createddate DESC, al.updateddate DESC`
  - mapping `type_id_log`: `1=email_alarm`, `2=line_alarm`, `3=telegram_alarm`, `4=sms_alarm`, `5=nonc_alarm`

## 3. Schedule process log + MQTT error log paginate — เสร็จสิ้น

### ScheduleProcessLogSpec (`internal/modules/settings/specs_alarm_logs.go`)
- alias `sl` บน `sd_schedule_process_log`; selects: sl.id, sl.schedule_id, sl.schedule_event_start, sl.day, sl.doday, sl.dotime, sl.schedule_event, sl.device_status + device block (`deviceActiveSelects`) ที่ join จาก d/t/mq/l
- Joins: `LEFT JOIN sd_iot_device d ON d.device_id = sl.device_id`, `sd_iot_device_type t`, `sd_iot_mqtt mq`, `sd_iot_location l`
- Filters ฝั่ง `d.*` ที่ NestJS ใช้ได้จริง: keyword→d.device_name (LIKE), org, bucket, type_id, status_warning, recovery_warning, status_alert, recovery_alert, time_life, period
- Intent-parity filters ฝั่ง `sl.*`: schedule_id, device_id, day, doday, dotime, schedule_event, device_status, status (NestJS อ้าง `st.*` / `l.schedule_id` / `l.device_id` ซึ่งเป็น SQL error)
- **เหตุผลที่ default sort เป็น `sl.id ASC`:** ORDER BY ของ NestJS อ้าง alias `alarm.` ที่ไม่มีอยู่ — endpoint นี้ error ทุกครั้งตอน sort; Go จึงใช้ sort เชิง deterministic แทน

### MqttErrorLogSpec
- แก้จาก `sd_mqtt_log` (spec เดิมเดาผิด) เป็น `alarmLogVariantSpec("sd_alarm_process_log_mqtt")` ตรงตาม service `mqttlogpaginate` ของ NestJS (controller route `/mqtterrorlogpaginate` → `settingsService.mqttlogpaginate`)
- ยืนยันแล้วว่า selects ของ NestJS รวม `al.subject AS subject` / `al.content AS content` และ spec ร่วมครอบคลุมครบ

## 4. Dashboard config upsert — เสร็จสิ้น

**ไฟล์:** `internal/modules/settings/delivery/http/handlers_dashboard.go`

- `POST /dashboardconfig` เดิมเป็น insert ธรรมดา → เปลี่ยนเป็น **upsert ตาม (name, location_id)** ตาม `createDashboardConfig` ของ NestJS:
  - พบแถวเดิม → update `config_data` = body.config + `updated_date` = now แล้ว query ซ้ำส่งแถวล่าสุดกลับ
  - ไม่พบ → insert `{location_id, name, config_data, created_date, updated_date}`
- แก้ mapping body key ให้ถูกต้องตาม DTO ของ NestJS (`dashboardConfig.dto.ts`): body รับ key `config` → column `config_data`

## 5. Master create dup-checks — เสร็จสิ้น

**ไฟล์หลัก:** `internal/modules/settings/delivery/http/handler_helpers.go` (+ call sites ทุก handler)

- เพิ่ม type `dupCheck{BodyKey, Column, Field}` รองรับแบบ variadic ใน generic `createHandler`
- ก่อน insert: ถ้ามีแถวใดในตารางมีค่าตรงกับที่ส่งมา → reject ด้วย HTTP 422 ข้อความ `"The <field> <value> duplicate this data cannot create."` (mirror ข้อความ NestJS)
- การ check เป็น equality parameterized ผ่าน `FixedFilter` ของ ListSpec (query LIMIT 1)

### ตาราง mapping ตาม endpoint ของ NestJS

| Endpoint | ตรวจซ้ำที่ column | หมายเหตุ |
|---|---|---|
| createsetting | sn | ข้อความ NestJS สะกด "createddate." — Go ใช้ข้อความมาตรฐาน |
| createlocation | ipaddress | |
| createtype | type_name | |
| createdevicetype | type_name | |
| createsensor | sensor_name | |
| creategroup | group_name | |
| createmqtt | mqtt_name + bucket | ตรวจ 2 ค่าตามลำดับ |
| createmqtthost | hostname | |
| createapi | api_name | |
| createemail | email_name | |
| createhost | host_name | |
| createline | line_name | |
| createnodered | nodered_name | |
| createsms | sms_name | |
| createtoken | token_name | |
|createtelegram | telegram_name | |
| createdevice | sn | |
| createschedule | schedule_name | |
| createinfluxdb | — | NestJS comment dup-check ทิ้งไว้ → จงใจไม่ใส่ |

- Update endpoints ไม่มี dup-check ทั้งสองฝั่ง (ยืนยันจาก `update_type` ฯลฯ ของ NestJS) — จึงไม่เพิ่ม

---

## ไฟล์ที่แก้ไขทั้งหมด

| ไฟล์ | สิ่งที่ทำ |
|---|---|
| `internal/modules/settings/specs_alarm_logs.go` | spec schedule process log + mqtt error log |
| `internal/modules/settings/delivery/http/handlers_dashboard.go` | upsert dashboardconfig |
| `internal/modules/settings/delivery/http/handler_helpers.go` | dupCheck type + logic ใน createHandler |
| `internal/modules/settings/delivery/http/handlers_master.go` | wire dup-checks 6 master |
| `internal/modules/settings/delivery/http/handlers_integration.go` | wire dup-checks integration masters |
| `internal/modules/settings/delivery/http/handlers_device.go` | dup-check sn |
| `internal/modules/settings/delivery/http/handlers_schedule.go` | dup-check schedule_name |

(ไฟล์ device-list/alarm-log จากงานช่วงต้นถูกรวมอยู่ใน specs/handlers ชุดเดียวกัน)

## การตรวจสอบคุณภาพ

ตรวจล่าสุด วันที่ 25 สิงหาคม 2026 เวลา 10:13 น.

| รายการ | คำสั่ง | ผลลัพธ์ |
|---|---|---|
| Build | `go build ./...` | ผ่าน |
| Static analysis | `go vet ./...` | ผ่าน |
| Unit tests | `go test ./internal/modules/settings/... -count=1` | ผ่านทั้งหมด (settings, delivery/http, usecase) |
| Format | `gofmt -l` บนไฟล์ที่แก้ | เรียบร้อย |

> หมายเหตุ: ผล `gofmt -l .` ทั้ง repo มีรายชื่อไฟล์ legacy/vendor จำนวนมากที่เป็นของเดิมใน repo ไม่เกี่ยวข้องกับงานครั้งนี้
