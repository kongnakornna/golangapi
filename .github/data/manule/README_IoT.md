# คู่มือการใช้งานโมดูล IoT (MQTT v3 / Device Monitoring & Control)

โมดูล `internal/modules/iot` เป็นระบบ **IoT Platform** หลักของ `icmongolang`: จัดการอุปกรณ์ (device) ที่เชื่อมต่อผ่าน MQTT broker, เก็บข้อมูลอนุกรมเวลา (time-series) ลง PostgreSQL + InfluxDB, ตรวจสถานะ/สัญญาณเตือน (alarm) แบบ real-time, สร้างกราฟ monitor, สั่งงานอุปกรณ์ (control), และส่งออกข้อมูล (export)

**สถาปัตยกรรม:**

```
               ┌──────────────────────┐
               │   MQTT Broker        │
               │   (mosquitto:1883)   │
               └──────────┬───────────┘
                          │  subscribe "#" (StartIngest)
                          ▼
              ┌──────────────────────────┐
              │   iot usecase (mqtt3UC)  │
              │  ProcessMqttData /       │
              │  Alarm / DeviceControl   │
              └──────┬──────────┬────────┘
                     │          │
        ┌────────────┴───┐  ┌───┴────────────┐
        │   PostgreSQL   │  │   Redis Cache  │
        │ iot_data,      │  │ mqtt_topic,    │
        │ device_status, │  │ mqtt_payload,  │
        │ device_config, │  │ iot_group, ... │
        │ activity_log   │  └───────────────┘
        └────────────┬───┘
                     │
        ┌────────────┴───┐
        │   InfluxDB     │  ← sensor charts (QueryFilterData)
        └────────────────┘
```

> **สำคัญ:** โมดูลนี้ถูก register ใน main API server ที่ base path **`/api/iot`** (`internal/server/handlers.go:216-252`)

---

## 1. ตารางไฟล์

| ไฟล์ | บทบาท |
|------|-------|
| `delivery/http/routes.go` | ลงทะเบียน route `/api/iot/*` |
| `delivery/http/handler.go` | HTTP handlers (874 บรรทัด) + Redis cache |
| `presenter/presenter.go` | Request/Response DTOs |
| `usecase/usecase.go` | Business logic หลัก (2013 บรรทัด) — ตรวจ alarm, monitor group, chart, ingest |
| `iothelper/alarm.go` | ตรรกะประเมิน alarm (Thai/English message) |
| `repository/*_repo.go` | GORM repositories (device, iot_data, device_status, device_config, ...) |
| `models/*.go` | ตารางทั้งหมด (device, sensor_data, alarm, notification, air, schedule, ...) |

---

## 2. การตั้งค่า Configuration

โมดูล IoT **ไม่มี config section แยก** — ใช้ config เดียวกับ MQTT + InfluxDB + Redis

### 2.1 `config/config.default.yml`
```yaml
mqtt:
  broker: "mqtt://localhost:1883"    # ⚠️ ควรใช้ tcp:// หรือ ssl://
  client_id: "icmongolang_client"
  username: ""
  password: ""
  qos: 1
  connection_timeout: 30

influxdb:
  url: "http://localhost:9087"
  token: ""
  org: "cmon_org"
  bucket: "AIRCOM1"
  timeout: 30
  username: ""
  password: ""

server:
  BaseUrl: "http://localhost:5000/api"   # ใช้สร้าง control/graph URL ใน monitor group
```

---

## 3. การเชื่อมต่อและ Wiring

`internal/server/handlers.go:216-252`:

1. สร้าง repositories: `DeviceRepository`, `AlarmLogRepository`, `DeviceStatusRepository`, `DeviceConfigRepository`, `IotDataRepository`, `ActivityLogRepository`, `CommandLogRepository`, `DeviceAlertRepository`
2. `iotUC := iotUsecase.NewMQTT3UseCase(deviceRepo, alarmLogRepo, mqttClient, redisClient, influxClient, logger, cfg, deviceStatusRepo, deviceConfigRepo, iotDataRepo, activityLogRepo, commandLogRepo, deviceAlertRepo)`
3. **ถ้า MQTT client เชื่อมต่ออยู่** → `iotUC.StartIngest(context.Background())` — เริ่ม subscribe `#` ทันที (รับข้อมูลทุก topic)
4. `iotHandler := iotHttp.NewMQTT3Handler(iotUC, logger, redisCache)`
5. `iotHttp.MapMQTT3Routes(apiRouter, iotHandler, mw)`

> ต่างจากโมดูล MQTT มาตรฐาน: โมดูล IoT **ลง route เสมอ** แม้ MQTT จะไม่เชื่อมต่อ (เพียงแต่จะไม่ StartIngest)

---

## 4. API Reference

ทุก route อยู่ใต้ **`/api/iot`** และถูก wrap ด้วย `mw.RateLimit()` (default 10 req/s, burst 20) ทุกตัว

### 4.1 ตารางรวม endpoints

| Method | Path | Auth | คำอธิบาย |
|--------|------|------|----------|
| GET | `/api/iot/topic` | ❌ | อ่านข้อมูลสดจาก topic MQTT |
| GET | `/api/iot/topicdevicechart` | ❌ | กราฟ + ค่าล่าสุดของ topic |
| GET | `/api/iot/controls` | ❌ | สั่งงานอุปกรณ์ (GET variant) |
| GET | `/api/iot/monitordevicegroup` | ❌ | หน้าจอ monitor รวม (กลุ่มอุปกรณ์) |
| GET | `/api/iot/monitordevicechart` | ❌ | กราฟ monitor ของ measurement |
| GET | `/api/iot/device` | ❌ | รายการอุปกรณ์ (paged) |
| GET | `/api/iot/devicebuckets` | ❌ | อุปกรณ์ใน bucket |
| GET | `/api/iot/sensercharts` | ❌ | กราฟ sensor จาก InfluxDB |
| GET | `/api/iot/locationdevice` | ❌ | อุปกรณ์ตาม location |
| GET | `/api/iot/devicesensercharts` | ❌ | alias ของ sensercharts |
| GET | `/api/iot/alarmdevicestatus` | ❌ | สถานะอุปกรณ์ + สัญญาณเตือน |
| GET | `/api/iot/alarmdevicestatuscontrol` | ❌ | เหมือนข้างบน (ฝั่ง control) |
| GET | `/api/iot/devicestatus` | ❌ | สถานะ online/active ของอุปกรณ์ |
| GET | `/api/iot/deviceconfig` | ❌ | config ของอุปกรณ์ |
| GET | `/api/iot/deviceiotdata` | ❌ | ข้อมูล iot_data (paged) |
| GET | `/api/iot/devicestats` | ❌ | สถิติข้อมูลอุปกรณ์ |
| GET | `/api/iot/devicedataexport` | ❌ | export ข้อมูล (csv/json) |
| POST | `/api/iot/control` | ✅ JWT | สั่งงานอุปกรณ์ (POST, body JSON) |
| PUT | `/api/iot/devicestatus` | ✅ JWT | อัปเดตสถานะอุปกรณ์ |
| PUT | `/api/iot/updatedeviceconfig` | ✅ JWT | อัปเดต config อุปกรณ์ |
| DELETE | `/api/iot/devicedatacleanup` | ✅ JWT | ลบข้อมูลเก่า |

---

### 4.2 `GET /api/iot/topic?topic=<topic>&delcache=1`

อ่านข้อมูลสดจาก MQTT topic (มี Redis cache)

| param | บังคับ | หมายเหตุ |
|-------|--------|----------|
| `topic` | ✅ | เช่น `CMONBUGKET01/DATA` |
| `delcache` | - | `1` = ลบ cache ก่อนดึงข้อมูลใหม่ |

**Response 200:**
```json
{
  "topic": "CMONBUGKET01/DATA",
  "payload": "25.5,60,1013",
  "from": "mqtt",          // "cache" | "mqtt" | "error" | "fallback" | "cache_fallback"
  "cache": false,
  "timestamp": "2026-06-08 12:36:58",
  "mqtt_connected": true,
  "cache_enabled": true
}
```

**ตรรกะ** (`usecase.go:137-218`):
- cache key `mqtt_topic:<topic>`, TTL **10s**
- ถ้า `delcache=1` → ลบ cache ก่อน
- MQTT fetch timeout **5s**; ถ้า MQTT ล้มเหลวแต่มี cache → `from="cache_fallback"`
- payload parse: JSON → object; มี comma → `[]string`; ไม่มี → string

---

### 4.3 `GET /api/iot/topicdevicechart`

กราฟ (45s cache) + payload ล่าสุด (10s cache) ของ topic เดียว

| param | บังคับ | ค่า default |
|-------|--------|-------------|
| `bucket` | - | `AIRCOM1` |
| `topic` | ✅ | — |
| `measurement` | - | `temperature` |
| `field` | - | `value` |
| `start` | - | `-10m` |
| `stop` | - | (none) |
| `limit` | - | `100` |
| `delcache` | - | (none) |

**Response keys:** `topic`, `chart{data, date}`, `latest_payload`, `latest_from`, `cache`, `latest_error?`

---

### 4.4 `GET /api/iot/controls?topic=<topic>&message=<message>`

**GET variant** ของการสั่งงาน — publish ข้อความไป topic ด้วย QoS 1 (ไม่ต้อง login)

**Response 200:**
```json
{ "status": "ok", "statusCode": "200" }
```

### 4.5 `POST /api/iot/control` (ต้อง auth)

**Request (JSON):**
```json
{ "topic": "CMONBUGKET01/CONTROL", "message": "ON" }
```

**Response 200:**
```json
{ "status": "ok" }
```

> ทั้งสอง variant ใช้ `mqttClient.Publish(topic, QOS=1, retained=false, message)` (`usecase.go:221-228`)

---

### 4.6 `GET /api/iot/monitordevicegroup`

หน้าจอ monitor รวมที่สมบูรณ์ที่สุด — กลุ่มอุปกรณ์ตาม `hardware_id` พร้อมค่าล่าสุด, สถานะ alarm, layout

| param | บังคับ | หมายเหตุ |
|-------|--------|----------|
| `bucket` | ✅ | |
| `location_id` | - | กรองตาม location |
| `hardware_id` | - | กรองตามชนิดฮาร์ดแวร์ |
| `lang` | - | `en` (default) / `th` |
| `delcache` | - | `1` = ล้าง cache ทั้งหมด |

**Response keys:** `bucket`, `timestamp`, `device_count`, `layout`, `layout_name`, `group_name`, `device_type`, `data[groups]`, `mqtt_connected`, `mqtt_raw_payload`, `cache_used`

**ตรรกะ** (`usecase.go:658-1017`):
- timeout รวม **15s** + panic recover
- **hardware_id → group:** 1=Sensor, 2=IO Sensor, 3=IO Control, 4=Critical Sensor
- **layout:** 1=Right Menu, 2=Card, 3=Left Menu, อื่นๆ=Footer Menu
- cache: `mqtt_device_list:<md5>` TTL **15 นาที**, `iot_group:<md5>` TTL **5 นาที**, payload `mqtt_payload:<bucket>` TTL 30s
- max 1,000 devices (เกิน → error)
- device แต่ละตัว enrich ถึง **30 fields** รวม `value_data` (พร้อม calibration add/subtract สำหรับ hardware_id=1), ประเมิน alarm ผ่าน `iothelper.AlarmDetailValidateEn/Th`
- ใช้ `cfg.Server.BaseUrl` สร้าง control URL (`/iot/controls`) และ graph URL (`/iot/monitordevicechart`)

---

### 4.7 `GET /api/iot/monitordevicechart`

| param | บังคับ | default |
|-------|--------|---------|
| `bucket` | ✅ | — |
| `measurement` | - | `temperature` |
| `field` | - | `value` |
| `start` | - | `-10m` |
| `stop` | - | `now()` |
| `limit` | - | `100` |

- cache key `mqtt_chart:<md5>`, TTL **45s**; `cache_delete=1` ล้าง
- Timestamp แปลงเป็น Asia/Bangkok รูปแบบ `2006-01-02:15:04:05` (colon separator)
- ถ้า InfluxDB error → `{data:[], date:[], cache:"error", error:...}` แต่ **HTTP 200** (ไม่ error)

---

### 4.8 `GET /api/iot/device` — รายการอุปกรณ์

| param | default | หมายเหตุ |
|-------|---------|----------|
| `page` | 1 | |
| `pageSize` | **1000** | |
| `bucket` | - | filter |
| `hardware_id` | - | filter |
| `type_id` | - | อ่านแต่ **ไม่ได้ใช้** ในการกรอง |
| `keyword` | - | filter |
| `lang` | - | อ่านแต่ไม่ได้ใช้ |

**Response 200:**
```json
{
  "data": [ { "device_id": "...", "device_name": "...", "value_data": "...", "status": 1, ... } ],
  "total": 42,
  "page": 1,
  "pageSize": 1000,
  "totalPage": 1
}
```

**`DeviceDetailResponse` fields:** `device_id`, `device_name`, `type_name`, `value_data`, `unit`, `status`, `alarm_title`, `status_warning`, `status_alert`, `recovery_warning`, `recovery_alert`, `icon`, `color_normal`, `color_warning`, `color_alert`

> ⚠️ `TypeName` ถูกตั้งให้ว่างเสมอใน usecase (`usecase.go:231-255`)

---

### 4.9 `GET /api/iot/devicestatus?deviceId=<id>` + `PUT /api/iot/devicestatus`

**GET:** คืนสถานะ online/active ของอุปกรณ์
```json
{
  "deviceId": "...",
  "isOnline": true,
  "isActive": true,
  "lastSeen": "2026-06-08T12:36:58+07:00",
  "batteryLevel": 85,
  "signalStrength": 3,
  "firmwareVersion": "v1.2",
  "location": {...},
  "lastData": {...},
  "uptime": "3h20m"
}
```
- ถ้ายังไม่มี record → สร้าง default `{IsOnline:false, IsActive:true}`
- `IsOnline` = `LastSeen` ภายใน 15 นาทีที่ผ่านมา
- `Uptime` = duration ตั้งแต่ `FirstSeen`

**PUT (ต้อง auth):** body JSON ต้องมี `deviceId`; รองรับ `battery`, `signal`, `firmware`, `location` → อัปเดต `LastSeen=now`, `LastData=json(body)`, `IsOnline=true` → response `{"status":"updated"}`

---

### 4.10 `GET /api/iot/deviceconfig?deviceId=<id>` + `PUT /api/iot/updatedeviceconfig`

**GET:** คืน config ของอุปกรณ์; ถ้าไม่มี → สร้าง default:
```json
{
  "general":   { "deviceName": "", "timezone": "Asia/Bangkok", "location": {"lat":0,"lng":0,"address":""} },
  "reporting": { "enabled": true, "interval": 300, "format": "json" },
  "thresholds":{ "temperature": {"min":15,"max":40}, "humidity": {"min":30,"max":80} },
  "alerts":    { "enabled": true, "email": [], "sms": [] }
}
```

**PUT (ต้อง auth):** deep-merge กับ config เดิม (recursive) แล้ว Upsert; **`deviceId` key ถูกลบออกก่อน merge** → response `{"status":"updated"}`

---

### 4.11 `GET /api/iot/deviceiotdata?deviceId=<id>&page=1&limit=50&startDate=&endDate=`

รายการ `iot_data` แบบ paged

**Response:**
```json
{
  "data": [ { "id": 1, "device_id": "...", "data": {...}, "timestamp": "...", "location": null, "metadata": null } ],
  "pagination": { "page": 1, "limit": 50, "total": 120, "pages": 3 }
}
```
- `startDate`/`endDate` ใช้ `time.Parse(time.RFC3339)` — วันที่ format ผิดจะถูก **ignore เงียบๆ**

### 4.12 `GET /api/iot/devicestats?deviceId=<id>`

```json
{ "count": 120, "firstRecord": "...", "lastRecord": "...", "dataPoints": {...} }
```
- ดึงล่าสุด 1000 records มาคำนวณ

### 4.13 `GET /api/iot/devicedataexport?deviceId=&startDate=&endDate=&format=csv|json`

Export ข้อมูล (ไม่ใช้ response envelope) — ตั้ง `Content-Type` + `Content-Disposition: attachment; filename=data.csv|data.json`

- `csv` → header `timestamp,device_id,data`, `text/csv`
- `json` (default) → `application/json`
- date format ผิด → 400

### 4.14 `DELETE /api/iot/devicedatacleanup?days=<n>` (ต้อง auth)

ลบ `iot_data` ที่เก่ากว่า `n` วัน (default **30**) → `{"deleted": <int64>}`

### 4.15 `GET /api/iot/sensercharts` / `devicesensercharts`

| param | บังคับ | หมายเหตุ |
|-------|--------|----------|
| `bucket` | ✅ | |
| `measurement` | ✅ | |
| `field` | - | |
| `start` / `stop` / `limit` | - | |

**Response:**
```json
{ "data": [25.5, 25.6, 25.4], "date": ["2026-06-08 12:36:58", ...], "cache": "no cache" }
```
- ใช้ `influxClient.QueryFilterData`, อ่าน `_value` + `_time` (format `2006-01-02 15:04:05`)

---

## 5. Data Ingest (StartIngest)

`usecase.go:1616-1640` — เมื่อ server เริ่มต้นและ MQTT เชื่อมต่อ:

1. `Subscribe("#", QOS 0)` — รับทุกข้อความ
2. แต่ละข้อความ → `extractDeviceNameFromTopic` (segment แรกของ topic) → หา device โดย `deviceNameCache` (ใน memory) + repo lookup ด้วย `mqtt_device_name`
3. `processIncoming` ใน goroutine (พร้อม panic recover)
4. `ProcessMqttData`: โหลด `MqttStatusDataName` mapping → `buildDataMap` (JSON object → direct; comma-split → map ตาม config; เก็บ `raw` เสมอ) → สร้าง `models.IotData` → อัปเดต device status → log activity `DATA_RECEIVED`

---

## 6. Alarm Logic (`iothelper/alarm.go`)

- `AlarmDetailValidate` (Thai) / `AlarmDetailValidateEn` (English) / `AlarmDetailValidateTh` (Thai)
- ผลลัพธ์ JSON keys: `status`, `status_control`, `alarm_type_id`, `type_id`, `hardware_id`, `alarm_status_set`, `title`, `subject`, `content`, `data_alarm`, `event_control`, `message_mqtt_control`, `sensor_data`, `count_alarm`, `unit`, `timestamp`
- ข้อความ Thai: `warning="คำเตือน มีความผิดปกติ"`, `critical="ภาวะวิกฤตต้องแก้ไขทันที"`, `recoveryWarning="คืนสู่ภาวะปกติ (คำเตือน)"`, `criticalMax="วิกฤต มีค่าสูงเกินกำหนด"`, `criticalMin="วิกฤต มีค่าต่ำกว่ากำหนด"`, `normal="ปกติ"`
- ตารางตัดสินใจ key ตาม `hardwareID`: 1=sensor, 2=IO sensor, 3=IO control, 4=critical sensor พร้อมค่า max/min/statusWarning/statusAlert/recovery thresholds

---

## 7. Redis Cache Keys

| Key | TTL | ใช้ที่ |
|-----|-----|-------|
| `mqtt_topic:<topic>` | 10s | GetTopicData |
| `mqtt_payload:<bucket>` | 30-60s | monitor group / alarm status |
| `mqtt_device_list:<md5>` | 15 min | monitor group |
| `iot_group:<md5>` | 5 min | monitor group |
| `mqtt_chart:<md5>` | 45s | monitor chart / topic chart |

---

## 8. โครงสร้างข้อมูลหลัก (PostgreSQL)

| ตาราง | ความหมาย |
|-------|----------|
| `sd_iot_device` | อุปกรณ์หลัก (50+ คอลัมน์: mqtt_id, type_id, hardware_id, bucket, calibration, status, layout, alert_set, ...) |
| `iot_data` | ข้อมูลอนุกรมเวลา (jsonb, index `(device_id, timestamp)`) |
| `device_status` | สถานะ online/active/last_seen (unique device_id) |
| `device_config` | config JSON ของอุปกรณ์ (unique device_id) |
| `device_alert` | วงจรชีวิตการแจ้งเตือน (severity, resolved/acknowledged, escalation) |
| `sd_sensor_data` | sensor reading (legacy) |
| `activity_log` / `command_log` | audit + คำสั่งอุปกรณ์ |
| `sd_iot_location` / `sd_iot_mqtt` / `sd_iot_host` | location / broker / host |
| `sd_iot_device_type` / `sd_device_group` / `sd_device_member` | type / group |
| `sd_air_*` | โมเดล air quality (control, mod, period, warning) |
| `sd_notification_*` / `sd_channel_template` | ระบบแจ้งเตือน |
| `sd_iot_schedule` | ตารางเวลา |

> schema ทั้งหมดอยู่ใน migration dump: `migrations/db.sql` และ `migrations/icmon.sql` (ไม่มีไฟล์ migration แยกของ IoT)

---

## 9. Edge Cases & ข้อควรระวัง

1. **GET endpoints เกือบทั้งหมดไม่ต้อง auth** — มีแค่ POST `/control`, PUT `/devicestatus`, PUT `/updatedeviceconfig`, DELETE `/devicedatacleanup` ที่ต้อง JWT (แม้ Swagger จะแสดง BearerAuth บน GET ก็ตาม)
2. **fallback data** — `GetTopicData` คืน `"30.50,80.5,25.50,75.5"` (`from:"fallback"`) เมื่อ MQTT timeout; `GetMonitorDeviceChart` คืนกราฟว่าง HTTP 200 เมื่อ InfluxDB error
3. **hardware_id กำหนดกรุ๊ป** — 1=Sensor, 2=IO Sensor, 3=IO Control, 4=Critical Sensor; calibration เพิ่ม/ลดเฉพาะ hardware_id=1
4. **`type_id` และ `lang` ใน `/device` ถูกอ่านแต่ไม่ใช้กรอง**
5. **max 1,000 devices** ต่อ monitor group
6. **StartIngest เฉพาะเมื่อ MQTT เชื่อมต่อตอน boot** — ข้อมูลจาก devices จะถูกบันทึกเฉพาะเมื่อมีกระบวนการนี้รัน
7. **ค่า default device config** — timezone `Asia/Bangkok`, reporting interval 300s, temp threshold 15-40°C, humidity 30-80%
8. **Payload comma-separated** ต้อง map กับ `MqttStatusDataName` config ของ device ถึงจะแยกค่าได้ถูก
9. **`/api/iot` มี RateLimit 10 req/s, burst 20** ทุก route
10. **timestamp ใน chart ใช้ format ต่างกัน** — monitor chart ใช้ colon `2006-01-02:15:04:05`, อื่นใช้ space

---

## 10. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| ไม่มีข้อมูล iot_data ใหม่ | MQTT ไม่เชื่อมต่อตอน boot → ดู log `StartIngest`; เช็ค broker/`tcp://` scheme |
| device ไม่เจอ payload | device ต้องมี `mqtt_device_name` ตรงกับ segment แรกของ topic |
| ได้ `from:"fallback"` | MQTT fetch timeout 5s; เช็คอุปกรณ์ publish จริงไหม |
| กราฟว่าง `cache:"error"` | InfluxDB down หรือ bucket/measurement/field ผิด |
| monitor group ใช้ค่าที่ cache เก่า | ล้างด้วย `delcache=1` หรือรอ TTL (5-15 นาที) |
| ค่า value_data ไม่ตรงอุปกรณ์ | เช็ค `calibration_add`/`calibration_subtract` ของ device |
| สั่งงานไม่ได้ | ตรวจว่าอุปกรณ์อยู่ใน group control (hardware_id=3) และ topic ลงท้าย CONTROL |
