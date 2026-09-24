# IoT & MQTT Modules API

## IoT (MQTT3) — `/api/iot`

> Rate-limited (ทุก route ใช้ `mw.RateLimit()`). ส่วนใหญ่เป็น public GET; บาง route ต้อง auth

### Public endpoints

| Method | Path | Handler |
|--------|------|---------|
| GET | `/iot/topic` | `GetTopicData` |
| GET | `/iot/topicdevicechart` | `GetTopicDataDeviceChart` |
| GET | `/iot/controls` | `DeviceControls` |
| GET | `/iot/monitordevicegroup` | `GetMonitorDeviceGroup` |
| GET | `/iot/monitordevicechart` | `GetMonitorDeviceChart` |
| GET | `/iot/device` | `GetDeviceList` |
| GET | `/iot/devicebuckets` | `GetDeviceBuckets` |
| GET | `/iot/sensercharts` | `GetSenserCharts` |
| GET | `/iot/locationdevice` | `GetDeviceByLocation` |
| GET | `/iot/devicesensercharts` | `GetDeviceSenserCharts` |
| GET | `/iot/alarmdevicestatus` | `GetAlarmDeviceStatus` |
| GET | `/iot/alarmdevicestatuscontrol` | `GetAlarmDeviceStatusControl` |
| GET | `/iot/devicemqtt` | `DeviceMqtt` |
| GET | `/iot/devicestatus` | `GetDeviceStatus` |
| GET | `/iot/deviceconfig` | `GetDeviceConfig` |
| GET | `/iot/deviceiotdata` | `ListIotData` |
| GET | `/iot/devicestats` | `GetDeviceStats` |
| GET | `/iot/devicedataexport` | `ExportData` |

### Auth endpoints

| Method | Path | Handler |
|--------|------|---------|
| POST | `/iot/control` | `DeviceControl` |
| PUT | `/iot/devicestatus` | `UpdateDeviceStatus` |
| PUT | `/iot/updatedeviceconfig` | `UpdateDeviceConfig` |
| DELETE | `/iot/devicedatacleanup` | `CleanupOldData` |

### GET `/api/iot/devicemqtt`

> สร้างขึ้นเพื่อให้ output ตรงกับ NestJS `mqtt2.controller.ts` endpoint `devicemqtt`
> Response wrap ใน `{ code, payload, message, message_th }`

**Query params:**
| Param | Type | Default | คำอธิบาย |
|-------|------|---------|----------|
| `bucket` | string | `""` | ชื่อ bucket (InfluxDB/MQTT) |
| `page` | int | `1` | หน้า |
| `pageSize` | int | `10000000` | จำนวน/หน้า |
| `lang` | string | `en` | `en` หรือ `th` |
| `device_id` | string | — | filter device |
| `mqtt_id` | string | — | filter mqtt |
| `type_id` | int | — | filter type |
| `hardware_id` | int | — | filter hardware |
| `keyword` | string | — | ค้นหา |
| `deletecache` | int | `0` | `1` = บังคับ refresh จาก MQTT (ลบ Redis cache) |

**Response `payload`:**

```jsonc
{
  "timestamps": "2006-01-02 15:04:05",
  "lang": "en",
  "page": 1,
  "currentPage": 1,
  "pageSize": 10000000,
  "totalPages": 1,
  "total": 0,
  "cache": "cache",            // หรือ "no cache"
  "device_count": 0,
  "connectionMqtt": { "isConnected": true, "connected": true, "status": 1, "msg": "...", "url": "tcp://..." },
  "device": [
    {
      "device_id": "...",
      "tiime": "...",
      "system_name": "...",          // location_name
      "location_name": "...",        // mqtt_name
      "zone_name": "...",            // type_name
      "device_name": "...",
      "hardware_type": "...",        // hardware_type_name (Sensor/IO Sensor/IO Control/Critical Sensor)
      "bucket": "...",               // mqtt_bucket
      "measurement": "...",
      "sensor_type": "Sensor",
      "devicedata": "12.34 ppm",     // hardware_id==1 : value+unit ; >1 : "ON"/"OFF"
      "alarm_title": "...",
      "alarm_detail": "...",
      "notification_status": "...",
      "notification_alarm_status": "...",
      "mqtt_data": "...",            // mqtt_data_value
      "mqtt_control": "...",         // mqtt_data_control
      "control": [],                 // hardware_id==1 ; >1 : "/iot/controls?topic=&message="
      "sensercharts": "/iot/sensercharts?bucket=&measurement=",
      "from": "cache",               // หรือ "mqtt"
      "ttl": 60,
      "cachelift": "cache",          // หรือ "no cache"
      "hardware_id": 1,
      "type_id": 0,
      "mqtt_id": 0,
      "type_name": "...",
      "location_name_lbl": "...",
      "unit": "...",
      "status": 1,
      "value_data": "12.34",
      "icon_access": "...",
      "alarm_subject": "...",
      "alarm_status": "..."
    }
  ]
}
```

### GET `/api/iot/alarmdevicestatus`

> เดิมคืน `{ statuscode, status, Mqttstatus, payload:{ checkConnectionMqtt, mqttrs, mqttname, bucket, time, mqttdata, deviceioinfo, devicesensor, deviceio, devicecritical, cache, chart } }`
> ปัจจุบันเพิ่ม field mqtt2-compatible: `lang`, `page`, `currentPage`, `pageSize`, `total`, `device_count`, `device` (array rss_data-style เดียวกับ `devicemqtt`)

**Required query:** `bucket` (ถ้าไม่มี → error "bucket is required"); `page`/`pageSize` ต้อง > 0

### GET `/api/iot/monitordevicegroup`
Required query: `bucket` (ถ้าไม่มี → error "bucket is required")
Query: `location_id`, `hardware_id`, `lang` (en/th), `page`, `pageSize`

### GET `/api/iot/device`
Query (ดู handler `GetDeviceList`): `bucket`, `page`, `pageSize`, `keyword`, `location_id`, `type_id`, `hardware_id` — Response: `DeviceDetailResponse[]`

---

## MQTT Gateway — `/api/mqtt`

> Rate-limited; 4 public GET, 3 auth POST

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| GET | `/mqtt/subscriptions` | `Subscriptions` | **Public** |
| GET | `/mqtt/status` | `Status` | **Public** |
| GET | `/mqtt/gettopicdata` | `GetTopicData` | **Public** |
| GET | `/mqtt/devicecontrol` | `DeviceControl` | **Public** |
| POST | `/mqtt/publish` | `Publish` | Auth |
| POST | `/mqtt/subscribe` | `Subscribe` | Auth |
| POST | `/mqtt/unsubscribe` | `Unsubscribe` | Auth |

---

## InfluxDB — `/api/influx`

> Mount แบบ conditional (เฉพาะเมื่อ influx client ไม่ nil). ทุก route auth

| Method | Path | Handler |
|--------|------|---------|
| POST | `/influx/write` | `WriteData` |
| GET | `/influx/query` | `QueryGetFilter` |
| POST | `/influx/devicechart` | `Querydevicechart` |
| POST | `/influx/filters` | `QueryFilters` |
| POST | `/influx/statistics` | `QueryStatistics` |

---

## Alarm Validation — `/api/alarm`

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| POST | `/alarm/validate` | `ValidateAlarm` | Auth |
| POST | `/alarm/validate/en` | `ValidateAlarmEn` | Auth |
| POST | `/alarm/validate/th` | `ValidateAlarmTh` | Auth |
