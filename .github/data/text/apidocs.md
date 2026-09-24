cdh# ICMON API Specification

> Base URL: `/api`
> Framework: Go + Chi
> Auth: JWT (RS256) Bearer tokens

---

## Contents

1. [Authentication `/auth`](#1-authentication)
2. [Users `/user` `/register` `/users`](#2-users)
3. [Items `/item`](#3-items)
4. [IoT Devices `/iot`](#4-iot-devices)
5. [MQTT `/mqtt`](#5-mqtt)
6. [InfluxDB `/influx`](#6-influxdb)
7. [Alarm Validation `/alarm`](#7-alarm-validation)
8. [Orders (Kafka) `/orders`](#8-orders-kafka)
9. [WebSocket `/ws`](#9-websocket)
10. [System & Health](#10-system--health)
11. [Common Response Envelope](#11-common-response-envelope)
12. [Authentication & Roles](#12-authentication--roles)
13. [Settings `/settings`](#13-settings--apisetttings)
14. [Business Modules](#14-business-modules)
15. [Infrastructure & Utility Modules](#15-infrastructure--utility-modules)

---

## 1. Authentication

### POST `/auth/login`
Login with username + password. Rate-limited (50 req/s, burst 100).  
**Auth:** Public  
**Request:**
```json
{ "username": "kongnakornna", "password": "password" }
```
**Response 200:**
```json
{ "access_token": "eyJ...", "refresh_token": "eyJ...", "token_type": "Bearer" }
```
**Error:** `400` `401` `404`

### POST `/auth/signin`
Login with email + password. Rate-limited (50 req/s, burst 100).  
**Auth:** Public  
**Request:**
```json
{ "email": "kongnakornna@gmail.com", "password": "password" }
```
**Response 200:** Token object (same as login)  
**Error:** `400` `401` `404`

### GET `/auth/publickey`
Get RSA public keys for JWT verification.  
**Auth:** Public  
**Response 200:**
```json
{
  "public_key_access_token": "-----BEGIN PUBLIC KEY-----\n...",
  "public_key_refresh_token": "-----BEGIN PUBLIC KEY-----\n..."
}
```
**Error:** `400` `422`

### GET `/auth/verifyemail`
Verify email using code.  
**Auth:** Public  
**Query:** `code` (string, required)  
**Response 200:** Success string  
**Error:** `400`

### POST `/auth/forgotpassword`
Send reset code to email.  
**Auth:** Public  
**Request:**
```json
{ "email": "kongnakornna@gmail.com" }
```
**Response 200:** Success string  
**Error:** `400` `422`

### PATCH `/auth/resetpassword`
Reset password using code.  
**Auth:** Public  
**Query:** `code` (string, required)  
**Request:**
```json
{ "new_password": "password", "confirm_password": "password" }
```
**Response 200:** Success string  
**Error:** `400` `422`

### GET `/auth/refresh`
Get new access token from refresh token.  
**Auth:** Bearer (refresh token)  
**Response 200:** Token object  
**Error:** `400` `422`

### GET `/auth/logout`
Remove current refresh token from DB.  
**Auth:** Bearer  
**Error:** `400` `422`

### GET `/auth/logoutall`
Logout all sessions for current user.  
**Auth:** Bearer  
**Error:** `400` `422`

---

## 2. Users

### POST `/register`
Register new user (public).  
**Auth:** Public  
**Request:**
```json
{
  "username": "kongnakornna", "password": "password",
  "confirm_password": "password", "email": "kongnakornna@gmail.com",
  "firstname": "Kongnakorn", "lastname": "Jantakun",
  "fullname": "Kongnakorn Jantakun", "phone_number": "021234567",
  "mobile_number": "0812345678", "line_id": "kongnakorn_line",
  "location_id": "loc_001", "role_id": 2
}
```
**Response 200:** Created user object  
**Error:** `400` `422`

### POST `/users`
Create user (public, no auth).  
**Auth:** Public  
**Request:** Same as `/register`  
**Error:** `400` `422`

### POST `/login`
Login with username + password (root-level variant).  
**Auth:** Public  
**Request:**
```json
{ "username": "kongnakornna", "password": "password" }
```
**Response 200:** Token object

### POST `/signin`
Login with email + password (root-level variant).  
**Auth:** Public  
**Request:**
```json
{ "email": "kongnakornna@gmail.com", "password": "password" }
```
**Response 200:** Token object

### GET `/user/me`
Get current user profile.  
**Auth:** Bearer  
**Response 200:**
```json
{
  "is_success": true,
  "data": {
    "id": "405d33ab-...", "email": "kongnakornna@gmail.com",
    "username": "kongnakornna", "role_id": 2,
    "firstname": "Kongnakorn", "lastname": "Jantakun",
    "fullname": "Kongnakorn Jantakun", "mobile_number": "0955088091",
    "phone_number": "0955088091", "line_id": "kongnakornna",
    "location_id": "Bangkok", "status": 1,
    "is_superuser": true, "verified": true,
    "last_sign_in": "2026-05-25T00:44:14+07:00",
    "createddate": "2026-05-25T00:44:14+07:00",
    "updateddate": "2026-05-25T00:44:14+07:00"
  }
}
```
**Error:** `401`

### PUT `/user/me`
Update current user profile.  
**Auth:** Bearer  
**Request:**
```json
{
  "firstname": "Kongnakorn", "lastname": "Jantakun",
  "fullname": "Kongnakorn Jantakun", "phone_number": "021234567",
  "mobile_number": "0812345678", "line_id": "kongnakorn_line",
  "location_id": "loc_001"
}
```
**Error:** `400`

### PATCH `/user/me/updatepass`
Change current user password.  
**Auth:** Bearer  
**Request:**
```json
{ "old_password": "oldpass", "new_password": "newpass", "confirm_password": "newpass" }
```
**Error:** `400`

### GET `/user/`
List users (admin).  
**Auth:** Bearer (SuperUser)  
**Query:** `limit` (int), `offset` (int), `email`, `username`, `status`, `role_id`  
**Error:** `400`

### POST `/user/`
Create user (admin).  
**Auth:** Bearer (SuperUser)  
**Request:** Same as `/register`

### GET `/user/{id}`
Get user by ID.  
**Auth:** Bearer  
**Error:** `400` `404`

### PUT `/user/{id}`
Update user (admin).  
**Auth:** Bearer (SuperUser)  
**Error:** `400`

### DELETE `/user/{id}`
Delete user (admin).  
**Auth:** Bearer (SuperUser)  
**Error:** `400` `404`

### PATCH `/user/{id}/role`
Update user role (admin).  
**Auth:** Bearer (SuperUser)  
**Request:**
```json
{ "role_id": 1 }
```
**Error:** `400`

### PATCH `/user/{id}/updatepass`
Update user password (admin).  
**Auth:** Bearer (SuperUser)  
**Request:**
```json
{ "old_password": "oldpass", "new_password": "newpass", "confirm_password": "newpass" }
```
**Error:** `400`

### GET `/user/{id}/logoutall`
Force logout all sessions (admin).  
**Auth:** Bearer (SuperUser)  
**Error:** `400`

---

## 3. Items

### GET `/item/`
List items with pagination.  
**Auth:** Bearer  
**Query:** `page` (int, default 1), `per_page` (int, default 10)  
**Response 200:**
```json
{
  "is_success": true,
  "data": {
    "items": [{ "id": "uuid", "title": "...", "description": "...", "owner_id": "uuid" }],
    "page": 1, "per_page": 10, "total": 100, "total_pages": 10
  }
}
```
**Error:** `400` `401` `422`

### POST `/item/`
Create item.  
**Auth:** Bearer  
**Request:**
```json
{ "title": "item title", "description": "item description" }
```
**Response 200:**
```json
{ "is_success": true, "data": { "id": "uuid", "title": "...", "description": "...", "owner_id": "uuid" } }
```
**Error:** `400` `401` `422`

### GET `/item/{id}`
Get item by ID.  
**Auth:** Bearer  
**Error:** `400` `401` `403` `404` `422`

### PUT `/item/{id}`
Update item.  
**Auth:** Bearer  
**Request:**
```json
{ "title": "new title", "description": "new description" }
```
**Error:** `400` `401` `403` `404` `422`

### DELETE `/item/{id}`
Delete item.  
**Auth:** Bearer  
**Error:** `400` `401` `403` `404` `422`

---

## 4. IoT Devices

### GET `/iot/topic`
Get live/cached MQTT topic data. Subscribes, waits 1 msg, unsubscribes. Cache 60s.  
**Auth:** Public  
**Query:** `topic` (string, required), `delcache` (string, "1" to refresh)  
**Response 200:**
```json
{
  "topic": "AIRCOM4/DATA", "payload": "...", "from": "mqtt",
  "timestamp": "2026-06-08T12:36:58+07:00",
  "mqtt_connected": true, "cache_enabled": true, "cache_hit": false,
  "data_length": 128, "fetch_duration_ms": 145
}
```
**Error:** `400` `500`

### GET `/iot/topicdevicechart`
Time-series + latest MQTT payload. Redis cache (45s / 10s).  
**Auth:** Public  
**Query:** `bucket` (default AIRCOM1), `topic` (required), `measurement` (default temperature), `field` (default value), `start` (default -10m), `stop` (default now()), `limit` (default 100), `delcache`  
**Error:** `400` `500`

### GET `/iot/controls`
Send control command via GET.  
**Auth:** Public  
**Query:** `topic` (required, e.g. BAACTW05/CONTROL), `message` (required, ON/OFF/1/0)  
**Error:** `400` `500`

### POST `/iot/control`
Send control command.  
**Auth:** Bearer  
**Request:**
```json
{ "topic": "CMONBUGKET01/CONTROL", "message": "ON" }
```
**Error:** `400` `500`

### GET `/iot/monitordevicegroup`
Get device groups for monitoring dashboard.  
**Auth:** Public  
**Query:** `bucket` (required), `location_id`, `hardware_id`, `lang`, `delcache`  
**Error:** `500`

### GET `/iot/monitordevicechart`
Get chart data for monitoring dashboard.  
**Auth:** Public  
**Query:** `bucket` (required), `measurement` (required), `field`, `start`, `stop`, `limit`  
**Error:** `500`

### GET `/iot/device`
List devices with pagination.  
**Auth:** Public  
**Query:** `page` (default 1), `pageSize` (default 10), `bucket`, `hardware_id`, `type_id`, `keyword`, `lang`  
**Response 200:**
```json
{
  "is_success": true,
  "data": {
    "data": [{
      "device_id": 1, "device_name": "Sensor-01", "type_name": "Temperature",
      "value_data": "25.5", "unit": "°C", "status": 1,
      "alarm_title": "", "status_warning": "", "status_alert": "",
      "recovery_warning": "", "recovery_alert": "",
      "icon": "", "color_normal": "", "color_warning": "", "color_alert": ""
    }],
    "total": 100, "page": 1
  }
}
```
**Error:** `500`

### GET `/iot/devicebuckets`
Get devices by bucket.  
**Auth:** Public  
**Query:** `bucket` (required)  
**Response 200:**
```json
{ "is_success": true, "data": { "bucket": "AIRCOM1", "devices": [...] } }
```
**Error:** `400` `500`

### GET `/iot/sensercharts`
Aggregated time-series chart data from InfluxDB.  
**Auth:** Public  
**Query:** `bucket` (required), `measurement` (required), `field` (default value), `start`, `stop`, `limit` (default 100)  
**Response 200:**
```json
{
  "is_success": true,
  "data": { "data": [25.5, 26.1, 24.8], "date": ["2026-06-08T12:00:00Z", "2026-06-08T12:01:00Z"], "cache": "redis" }
}
```
**Error:** `400` `500`

### GET `/iot/devicesensercharts`
Alias for `/iot/sensercharts`. Same params & response.  
**Auth:** Public

### GET `/iot/locationdevice`
Get devices by location ID.  
**Auth:** Public  
**Query:** `location_id` (int, required)  
**Error:** `400` `500`

### GET `/iot/alarmdevicestatus`
Get alarm status of devices.  
**Auth:** Public  
**Query:** `bucket`, `measurement`, `device_id`, `type_id`, `hardware_id`, `page` (default 1), `pageSize` (default 1000)  
**Error:** `500`

### GET `/iot/alarmdevicestatuscontrol`
Get alarm status with control info.  
**Auth:** Public  
**Query:** `bucket`, `device_id`, `type_id`, `hardware_id`  
**Error:** `500`

### GET `/iot/devicestatus`
Get current device status.  
**Auth:** Public  
**Query:** `deviceId` (string, required)  
**Response 200:**
```json
{
  "is_success": true,
  "data": {
    "deviceId": "1", "isOnline": true, "isActive": true,
    "lastSeen": "2026-06-08T12:00:00Z", "batteryLevel": 85,
    "signalStrength": -65, "firmwareVersion": "v2.1.0",
    "location": { "lat": 13.75, "lng": 100.5 },
    "lastData": { "temperature": 25.5 }, "uptime": "72h30m"
  }
}
```
**Error:** `400` `500`

### PUT `/iot/devicestatus`
Update device status.  
**Auth:** Bearer  
**Request:** Arbitrary JSON (must include `deviceId`)
```json
{ "deviceId": "1", "isOnline": true, "isActive": true }
```
**Error:** `400` `500`

### GET `/iot/deviceconfig`
Get device configuration.  
**Auth:** Public  
**Query:** `deviceId` (string, required)  
**Response 200:** Device config JSON  
**Error:** `400` `500`

### PUT `/iot/updatedeviceconfig`
Update device config (merge).  
**Auth:** Bearer  
**Request:** Arbitrary JSON (must include `deviceId`)
```json
{ "deviceId": "1", "thresholds": { "temp_min": 10, "temp_max": 40 } }
```
**Error:** `400` `500`

### GET `/iot/deviceiotdata`
List IoT data with pagination.  
**Auth:** Public  
**Query:** `deviceId` (required), `page`, `limit`, `startDate` (RFC3339), `endDate` (RFC3339)  
**Response 200:**
```json
{
  "is_success": true,
  "data": {
    "data": [{ "id": 1, "device_id": "1", "data": { "temperature": 25.5 }, "timestamp": "2026-06-08T12:00:00Z", "location": null, "metadata": null }],
    "pagination": { "page": 1, "limit": 20, "total": 100, "pages": 5 }
  }
}
```
**Error:** `400` `500`

### GET `/iot/devicestats`
Get device statistics.  
**Auth:** Public  
**Query:** `deviceId` (string, required)  
**Response 200:**
```json
{
  "is_success": true,
  "data": {
    "count": 1500,
    "firstRecord": "2026-01-01T00:00:00Z",
    "lastRecord": "2026-06-08T12:00:00Z",
    "dataPoints": { "temperature": { "min": 10, "max": 40 } }
  }
}
```
**Error:** `400` `500`

### GET `/iot/devicedataexport`
Export IoT data as CSV/JSON file.  
**Auth:** Public  
**Query:** `deviceId` (required), `startDate` (required, RFC3339), `endDate` (required, RFC3339), `format` (csv/json, default json)  
**Response 200:** File download  
**Error:** `400` `500`

### DELETE `/iot/devicedatacleanup`
Delete old IoT data.  
**Auth:** Bearer  
**Query:** `days` (int, default 30)  
**Error:** `500`

---

## 5. MQTT

### GET `/mqtt/subscriptions`
List subscribed MQTT topics.  
**Auth:** Public  
**Response 200:**
```json
{ "topics": ["BAACTW01/DATA", "BAACTW02/DATA"], "count": 2 }
```
**Error:** `500`

### GET `/mqtt/status`
Get MQTT connection status.  
**Auth:** Public  
**Response 200:**
```json
{ "connected": true, "broker": "tcp://localhost:1883", "client_id": "icmon-client", "timestamp": "2026-06-08T12:36:58+07:00" }
```
**Error:** `500`

### GET `/mqtt/gettopicdata`
Subscribe, wait 1 msg, unsubscribe. Cache 60s.  
**Auth:** Public  
**Query:** `topic` (string, required)  
**Response 200:**
```json
{
  "statuscode": 200, "code": 200, "topic": "AIRCOM4/DATA", "payload": "...",
  "from": "mqtt", "timestamp": "2026-06-08T12:36:58+07:00",
  "message": "live data", "mqtt_connected": true,
  "cache_enabled": true, "cache_hit": false,
  "data_length": 128, "fetch_duration_ms": 145
}
```
**Error:** `400` `500`

### GET `/mqtt/devicecontrol`
Send control and wait for DATA topic response.  
**Auth:** Public  
**Query:** `topic` (required, e.g. BAACTW02/CONTROL), `message` (required, ON/OFF/1/0)  
**Response 200:**
```json
{
  "statuscode": 200, "code": 200,
  "topic_control": "BAACTW02/CONTROL", "topic_data": "BAACTW02/DATA",
  "message_sent": "ON", "payload": "25.5,60,1013",
  "data": ["25.5", "60", "1013"], "status": 1, "status_msg": "ON",
  "timestamp": "2026-06-08T18:30:00+07:00",
  "message": "Control sent to BAACTW02/CONTROL, response received",
  "message_th": "ส่งคำสั่งไปยัง BAACTW02/CONTROL และได้รับข้อมูลตอบกลับ",
  "from": "mqtt", "fetch_duration_ms": 245
}
```
**Error:** `400` `500`

### POST `/mqtt/publish`
Publish MQTT message.  
**Auth:** Bearer  
**Request:**
```json
{ "topic": "test/topic", "qos": 1, "retained": false, "payload": { "key": "value" } }
```
**Response 200:**
```json
{ "success": true, "message": "Message published to test/topic" }
```
**Error:** `400` `500`

### POST `/mqtt/subscribe`
Persistent MQTT subscription. Messages via WebSocket.  
**Auth:** Bearer  
**Request:**
```json
{ "topic": "test/topic", "qos": 1 }
```
**Response 200:**
```json
{ "topic": "test/topic", "qos": 1, "subscribed_at": "2026-06-08T12:36:58+07:00", "message": "Subscribed successfully" }
```
**Error:** `400` `500`

### POST `/mqtt/unsubscribe`
Unsubscribe from MQTT topic.  
**Auth:** Bearer  
**Request:**
```json
{ "topic": "test/topic" }
```
**Response 200:**
```json
{ "success": true, "topic": "test/topic", "message": "Unsubscribed successfully" }
```
**Error:** `400` `500`

---

## 6. InfluxDB

### POST `/influx/write`
Write data point to InfluxDB.  
**Auth:** Bearer  
**Request:**
```json
{ "measurement": "temperature", "fields": { "value": 25.5 }, "tags": { "device_id": "sensor-01" } }
```
**Error:** `400` `500`

### GET `/influx/query`
Query raw time-series via query params.  
**Auth:** Bearer  
**Query:** `measurement` (required), `field` (required), `start`, `stop` (default now()), `limit` (default 1000), `offset` (default 0)  
**Response 200:**
```json
[
  { "time": "2026-06-08T12:00:00Z", "value": 25.5 },
  { "time": "2026-06-08T12:01:00Z", "value": 26.1 }
]
```
**Error:** `400` `500`

### POST `/influx/devicechart`
Query device chart data.  
**Auth:** Bearer  
**Request:**
```json
{ "measurement": "temperature", "field": "value", "bucket": "AIRCOM1", "start": "-1h", "stop": "now()", "limit": 10000, "offset": 1 }
```
**Response 200:** Array of `{ time, value }`  
**Error:** `400` `500`

### POST `/influx/filters`
Query filtered data points.  
**Auth:** Bearer  
**Request:** Same as `/influx/devicechart`  
**Response 200:** Array of `{ time, value }`  
**Error:** `400` `500`

### POST `/influx/statistics`
Compute statistics over time-series data.  
**Auth:** Bearer  
**Request:**
```json
{ "measurement": "temperature", "bucket": "AIRCOM1", "field": "value", "aggregate": "mean", "start": "-15s", "stop": "now()", "windowPeriod": "15s", "percentile": 95 }
```
**Response 200:**
```json
{ "aggregate": "mean", "value": 25.5, "data": [{ "time": "2026-06-08T12:00:00Z", "value": 25.5 }] }
```
**Error:** `400` `500`

---

## 7. Alarm Validation

### POST `/alarm/validate`
Validate alarm status (Thai).  
**Auth:** Bearer  
**Request:**
```json
{
  "actionName": "action1", "deviceName": "Sensor-01", "hardwareID": 1, "unit": "°C",
  "valueData": 25.5, "sensorValueData": 25.5, "min": 10, "max": 30,
  "valueAlarm": 28, "statusAlert": 1, "statusWarning": 0,
  "recoveryAlert": 0, "recoveryWarning": 0, "countAlarm": 3,
  "mqttName": "sensor/mqtt", "mqttControlOn": "ON", "mqttControlOff": "OFF",
  "event": 0, "valueRelay": 0, "valueControlRelay": 0
}
```
**Response 200:**
```json
{
  "alarmActionName": "action1", "title": "แจ้งเตือน: อุณหภูมิสูง",
  "subject": "การแจ้งเตือนจาก Sensor-01",
  "content": "ค่า temperature อยู่ที่ 25.5 °C ซึ่งเกินค่าที่กำหนด",
  "lang": "th", "deviceNameStr": "Sensor-01", "deviceNameVal": "Sensor-01",
  "mqttControlOnVal": "ON", "mqttControlOffVal": "OFF",
  "status": 1, "statusAlertVal": 1, "statusWarningVal": 0,
  "recoveryAlertVal": 0, "recoveryWarningVal": 0,
  "typeId": 1, "alarmTypeId": 1, "hardwareId": 1, "alarmStatusSet": 1,
  "unit": "°C", "timestamp": "2026-06-08T12:00:00Z",
  "messageMqttControl": "...", "eventControl": 0, "eventVal": 0,
  "dataAlarm": 0, "dataAlarmRaw": 0, "sensorData": null, "sensorValue": null
}
```
**Error:** `400`

### POST `/alarm/validate/en`
Validate alarm (English).  
**Auth:** Bearer  
**Request/Response:** Same structure, `"lang": "en"`  
**Error:** `400`

### POST `/alarm/validate/th`
Validate alarm (Thai, explicit).  
**Auth:** Bearer  
**Request:** Same as `/alarm/validate`  
**Error:** `400`

---

## 8. Orders (Kafka)

> Note: Served by a separate `cmd/kafka` binary, not the main API server.

### POST `/orders`
Create order (async via Kafka).  
**Auth:** Bearer  
**Request:**
```json
{ "product_id": "prod-001", "quantity": 2 }
```
**Response 202:**
```json
{ "order_id": "uuid", "status": "pending", "message": "Order accepted for processing" }
```
**Error:** `400` `500`

---

## 9. WebSocket

### GET `/ws`
WebSocket upgrade for real-time streaming.  
**Auth:** Token (query param `token` or `Authorization` header)

**Client → Server:**
```json
{ "type": "subscribe",   "topic": "BAACTW01/DATA" }
{ "type": "unsubscribe", "topic": "BAACTW01/DATA" }
{ "type": "message",     "payload": { "text": "hello" }, "room": "general" }
```

**Server → Client:**
```json
{ "type": "message",      "topic": "BAACTW01/DATA", "payload": "25.5",  "timestamp": "..." }
{ "type": "notification", "title": "...",           "body": "...",      "timestamp": "..." }
{ "type": "error",        "message": "..." }
```

---

## 10. System & Health

### GET `/ping`
Health ping.  
**Auth:** Public  
**Response 200:**
```json
{ "status": "ok", "timestamp": "2026-07-11 00:00:00 UTC" }
```

### GET `/health`
Detailed health check.  
**Auth:** Public  
**Response 200:**
```json
{ "status": "up", "checks": { "database": true, "redis": true, "influxdb": true, "mqtt": true, "websocket": true } }
```

### GET `/metrics`
Prometheus metrics (text format).  
**Auth:** Public

### GET `/apimetric`
Custom JSON metrics.  
**Auth:** Public  
**Response 200:**
```json
{ "active_requests": 5, "total_requests": 15000, "uptime_seconds": 3600 }
```

### GET `/swagger/*`
Swagger UI.  
**Auth:** Public

---

## 11. Common Response Envelope

**Success (most endpoints):**
```json
{ "is_success": true, "data": { ... } }
```

**Error (auth, users, items, IoT):**
```json
{ "is_success": false, "error": { "status": 404, "statusText": "not_found", "msg": "resource not found" } }
```

**Error (alarm, influx, mqtt):**
```json
{ "error": "error description" }
```

---

## 12. Authentication & Roles

```
Authorization: Bearer <access_token>
```

Tokens are JWT (RS256). Obtain from `POST /api/auth/login` or `POST /api/auth/signin`.

### Token Lifetimes
| Token | Duration | Env Key |
|-------|----------|---------|
| Access Token | ~24h | `JWT_ACCESS_TOKEN_EXPIRE_DURATION` |
| Refresh Token | ~24h | `JWT_REFRESH_TOKEN_EXPIRE_DURATION` |

### Role Hierarchy
| Role | Level | Description |
|------|-------|-------------|
| SUPERADMIN | 1 | Full system access |
| ADMIN | 2 | Admin access |
| EDITOR | 3 | Can edit resources |
| MONITOR | 4 | Read-only monitoring |
| USER | 5 | Basic authenticated user |

### Middleware Chain (Protected Routes)
1. `Verifier` — validates JWT signature
2. `Authenticator` — checks token validity / blacklist
3. `CurrentUser` — loads user from DB
4. `ActiveUser` — verifies user is active
5. `SuperUser` — checks SUPERADMIN role (admin-only routes)

---

**Endpoints (core): 72 total** — 9 Auth · 15 Users · 5 Items · 21 IoT · 7 MQTT · 5 InfluxDB · 3 Alarm · 1 Orders · 1 WS · 5 System

> ⚠️ 260+ รายการ ด้านล่างนี้เป็น module ที่เพิ่มเติม (Settings / Business / Infra) — รายละเอียดเต็มทุก endpoint อยู่ในโฟลเดอร์ [`docs/apidocs/`](./apidocs/README.md) และฉบับรวม [`docs/apidocs_full.md`](./apidocs_full.md)

---

## 13. Settings — `/api/settings`

> Auth ทุก endpoint (JWT Bearer) · รายละเอียดเต็ม: [`docs/apidocs/settings.md`](./apidocs/settings.md)

### Master
- **Setting:** `GET /settings/listsetting`, `GET /settings/settingall`, `POST /settings/createsetting`, `POST /settings/updatesetting`, `GET|DELETE /settings/deletesetting`
- **Location:** `GET /settings/listlocation`, `GET /settings/locationall`, `POST /settings/createlocation`, `POST /settings/updatelocation`, `GET /settings/deletelocation`
- **Type:** `GET /settings/listtype`, `GET /settings/typeall`, `POST /settings/createtype`, `POST /settings/updatetype`, `GET /settings/deletetype`
- **DeviceType:** `GET /settings/listdevicetype`, `GET /settings/devicetypeall`, `GET /settings/devicetypeallcontrol`, `POST /settings/createdevicetype`, `POST /settings/updatedevicetype`, `GET /settings/deletedevicetype`
- **Group:** `GET /settings/lisgroup`, `GET /settings/lisgroupall`, `GET /settings/listgrouppage`, `POST /settings/creategroup`, `POST /settings/updategroup`, `GET /settings/deletegroup`
- **Sensor:** `GET /settings/listsensor`, `GET /settings/sensorall`, `POST /settings/createsensor`, `POST /settings/updatesensor`, `GET /settings/deletesensor`

### Device & Schedule
- **Device:** `GET /settings/deviceall`, `GET /settings/listdevicepage`, `GET /settings/listdevicepagess`, `GET /settings/listdevicepageactive1`, `GET /settings/listdevicepageactive`, `GET /settings/listdevicepageall`, `GET /settings/listdevicepagesensor`, `GET /settings/listdevicepageallactive`, `GET /settings/listdevicepageallactiveschedule`, `GET /settings/deviceeditget`, `GET /settings/devicedetail`, `GET /settings/devicedelete`, `GET /settings/deletedevice`, `POST /settings/createdevice`, `POST /settings/updatedevice`, `POST /settings/updatestatusdeviceid`, `POST /settings/deviceactionuser`
- **Schedule:** `GET /settings/listscheduledevice`, `GET /settings/findscheduledevicechk`, `GET /settings/schedulelist`, `GET /settings/scheduleall`, `GET /settings/listschedulepage`, `GET /settings/scheduledevicepage`, `GET /settings/listdevicescheduledata`, `GET /settings/createscheduledevice`, `GET /settings/deletescheduledevice`, `GET /settings/deletedeviceschedule`, `GET /settings/deletedeviceandschedule`, `POST /settings/createschedule`, `POST /settings/createscheduledevice`, `POST /settings/updateschedule`, `POST /settings/updateschedulestatus`, `POST /settings/updatescheduledaystatus`, `GET /settings/deleteschedule`

### Integration
- **MQTT:** `GET /settings/lismqtt`, `/settings/listmqtt`, `/settings/lismqttall`, `/settings/getmqttdetail`, `/settings/listmqttpaginate`, `/settings/listmqttpaginateactive`, `/settings/listmqttdevicepaginate`, `/settings/mqttdelete`, `/settings/deletemqtt`, `POST /settings/createmqtt`, `POST /settings/updatemqtt`, `POST /settings/updatemqttstatus`, `POST /settings/mqtttsort`
- **MQTTHost:** `GET /settings/listmqtthost`, `/settings/mqtthostall`, `POST /settings/createmqtthost`, `/settings/updatemqtthost`, `/settings/updatemqtthoststatus`, `GET /settings/deletemqtthost`
- **API:** `GET /settings/apiall`, `/settings/listapipage`, `POST /settings/createapi`, `/settings/updateapi`, `GET /settings/deleteapi`
- **Email:** `GET /settings/listemail`, `/settings/emailall`, `POST /settings/createemail`, `/settings/updateemail`, `/settings/updateemailstatus`, `GET /settings/deleteemail`
- **Host:** `GET /settings/hostall`, `/settings/listhostpage`, `POST /settings/createhost`, `/settings/updatehost`, `GET /settings/deletehost`
- **InfluxDB:** `GET /settings/influxdball`, `/settings/listinfluxdbpage`, `POST /settings/createinfluxdb`, `/settings/updateinfluxdb`, `/settings/updateinfluxdbstatus`, `GET /settings/deleteinfluxdb`
- **Line:** `GET /settings/lineall`, `/settings/listlinepage`, `POST /settings/createline`, `/settings/updateline`, `/settings/updatelinestatus`, `GET /settings/deleteline`
- **NodeRED:** `GET /settings/noderedall`, `/settings/listnoderedpaginate`, `POST /settings/createnodered`, `/settings/updatenodered`, `/settings/updatenoderedstatus`, `GET /settings/deletenodered`
- **SMS:** `GET /settings/smsall`, `/settings/listsmspage`, `POST /settings/createsms`, `/settings/updatesms`, `/settings/updatesmsstatus`, `GET /settings/deletesms`
- **Token:** `GET /settings/tokenall`, `/settings/tokensmspage`, `POST /settings/createtoken`, `/settings/updatetoken`, `GET /settings/deletetoken`
- **Telegram:** `POST /settings/createtelegram`, `/settings/updatetelegram`, `GET /settings/deletetelegram`

### Dashboard Config
`POST /settings/dashboardconfig`, `GET /settings/dashboardconfig_1`, `GET /settings/dashboardconfig`, `GET /settings/dashboardconfig/search`, `GET/PATCH/DELETE /settings/dashboardconfig/{id}`

### Alarm Device & Logs
- **Alarm links:** `GET /settings/listalarmdevicepage`, `/settings/listalarmdeviceactivepage`, `/settings/listalarmeventdevicepage`, `/settings/listalarmeventdevicecontrolpage`, `/settings/activealarmdevicepage`, `/settings/activealarmeventdeviceeventpage`, `/settings/alarmdevice`, `/settings/alarmdevicestatus`, `/settings/deviceactivemqttalarm`, `/settings/devicealarm`, `/settings/createalarmdevice`, `/settings/deletealarmdevice`, `/settings/createalarmeventdevice`, `/settings/deletealarmeventdevice`, `/settings/deletearmdevice`, `/settings/deletearmdevicev2`, `POST /settings/createalarmDevice`, `/settings/createalarmdevicepaginate`, `/settings/updatealarmdevice`, `/settings/updatealarmstatus`, `/settings/createdevicealarmaction`
- **Monitors:** `GET /settings/listdevicealarm`, `/settings/listdevicealarmairV1`, `/settings/listdevicealarmair`, `/settings/listdevicealarmall`, `/settings/_listdevicealarmfan`, `/settings/listdevicealarmfan`, `/settings/listdevicealarmlimit`, `/settings/_devicemonitor`, `/settings/devicemonitor`, `/settings/devicemonitors`
- **Logs:** `GET /settings/scheduleproces`, `/settings/scheduleprocesslog`, `/settings/scheduleprocesslogpaginate`, `/settings/mqtterrorlogpaginate`, `/settings/alarmlogpaginate`, `/settings/alarmlogpaginateemail`, `/settings/alarmlogpaginateline`, `/settings/alarmlogpaginatesms`, `/settings/alarmlogpaginatetelegram`, `/settings/alarmlogpaginatecontrols`, `/settings/alarmlogpaginatecontrol`

---

## 14. Business Modules

> Auth ทุก endpoint · รายละเอียดเต็ม: [`docs/apidocs/business.md`](./apidocs/business.md)

### Quotation — `/api/quotation`
| Method | Path |
|--------|------|
| GET/POST | `/quotation` |
| GET/PUT/DELETE | `/quotation/{id}` |
| GET | `/quotation/{id}/pdf` |

### Purchase Order — `/api/purchase-orders`
| Method | Path |
|--------|------|
| GET/POST | `/purchase-orders` |
| GET | `/purchase-orders/suggestions/{jobId}` |
| POST | `/purchase-orders/from-quotation/{quotationId}` |
| GET/PUT/DELETE | `/purchase-orders/{id}` |
| POST | `/purchase-orders/{id}/send` |
| PUT | `/purchase-orders/{id}/confirm` |
| POST | `/purchase-orders/{id}/receive` |
| PUT | `/purchase-orders/{id}/cancel` |
| GET | `/purchase-orders/{id}/pdf` · `/purchase-orders/{id}/history` |

### Payment & Receipt — `/api/payments` `/api/receipts`
| Method | Path |
|--------|------|
| POST/GET | `/payments` |
| POST | `/payments/search` |
| GET | `/payments/outstanding/{customerId}` · `/payments/history/{customerId}` · `/payments/{id}` · `/payments/invoice/{invoiceId}` |
| POST | `/payments/{id}/refund` |
| PUT | `/payments/{id}/cancel` |
| GET | `/receipts/{id}` · `/receipts/{id}/pdf` · `/receipts/payment/{paymentId}` |
| PUT | `/receipts/{id}/cancel` |

### Job — `/api/job`
| Method | Path |
|--------|------|
| GET/POST | `/job` |
| GET/PUT/DELETE | `/job/{id}` |
| PUT | `/job/{id}/status` |
| GET | `/job/{id}/history` · `/job/{id}/report` · `/job/{id}/pdf` · `/job/{id}/picking/pdf` · `/job/{id}/delivery/pdf` |
| POST/GET | `/job/{id}/services` · `/job/{id}/parts` |

### Customer & Car — `/api/customer` `/api/car`
| Method | Path |
|--------|------|
| GET/POST | `/customer` |
| GET/PUT/DELETE | `/customer/{id}` |
| GET/POST | `/car` |
| GET/PUT/DELETE | `/car/{id}` |

---

## 15. Infrastructure & Utility Modules

> ส่วนใหญ่ auth ยกเว้นที่ระบุ public · รายละเอียดเต็ม: [`docs/apidocs/infra.md`](./apidocs/infra.md)

### Dashboard — `/api/dashboard`
`GET /dashboard/stats`, `GET /dashboard/revenue`, `GET /dashboard/top-parts`, `GET /dashboard/job-status`

### Document — `/api/documents`
`POST|GET /documents/`, `GET|DELETE /documents/{id}`

### Email — `/api/email`
`POST /email/send`, `GET /email/logs`, `GET /email/logs/{id}`, `GET|PUT /email/config`

### Batch — `/api/batch`
`POST|GET /batch/jobs`, `GET|PUT|DELETE /batch/jobs/{id}`, `POST /batch/jobs/{id}/run`, `GET /batch/jobs/{id}/logs`

### I18n — `/api/i18n`
`GET|POST /i18n/translations`, `GET|PUT|DELETE /i18n/translations/{key}`

### WOS — `/api/wos`
`POST|GET /wos/orders`, `GET /wos/orders/{id}`, `PUT /wos/orders/{id}/status`

### Report — `/reports` (ไม่มี `/api`)
`GET /reports/daily-sales/pdf`, `/reports/inventory-summary/pdf`, `/reports/customer-list/pdf`, `/reports/invoice/pdf`, `/reports/credit-note/pdf`, `/reports/debit-note/pdf`, `POST /reports/token`
> PDF routes รองรับ full token หรือ short-lived `?token=`; `/reports/token` ใช้ full token เท่านั้น

### VectorData — `/api/vectordata`
`GET /vectordata/health` (**public**), `POST /vectordata/create-index`, `POST|GET /vectordata/documents`, `GET|PUT|DELETE /vectordata/documents/{id}`, `POST /vectordata/search`, `POST /vectordata/seed`

### Elasticsearch — `/api/elasticsearch`
`GET /elasticsearch/health` (**public**), `POST /elasticsearch/create-index`, `POST /elasticsearch/embeddings`, `POST /elasticsearch/vectors`, `POST /elasticsearch/search`
