# ICMON API Specification — ฉบับรวม (Full)

> Base URL: `/api` · Framework: Go + Chi · Auth: JWT (RS256) Bearer tokens
> นี่คือคู่มือฉบับรวมทุก module ในไฟล์เดียว สำหรับเวอร์ชันแยกตาม module ดู [`docs/apidocs/`](./apidocs/README.md)

## สารบัญ
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
11. [Settings `/settings`](#11-settings--apisetttings)
12. [Business Modules](#12-business-modules)
13. [Infrastructure & Utility](#13-infrastructure--utility)
14. [Common Response Envelope](#14-common-response-envelope)
15. [Authentication & Roles](#15-authentication--roles)

---

## 1. Authentication

### POST `/auth/login`
Login with username + password. Rate-limited (50 req/s, burst 100). **Public**
```json
{ "username": "kongnakornna", "password": "password" }
```
Response: `{ "access_token": "eyJ...", "refresh_token": "eyJ...", "token_type": "Bearer" }` · Error: `400` `401` `404`

### POST `/auth/signin`
Login with email + password. **Public**
```json
{ "email": "kongnakornna@gmail.com", "password": "password" }
```

### GET `/auth/publickey` (**public**)
Response: `{ "public_key_access_token": "...", "public_key_refresh_token": "..." }`

### GET `/auth/verifyemail` (**public**)
Query: `code` (required)

### POST `/auth/forgotpassword` (**public**)
Body: `{ "email": "..." }`

### PATCH `/auth/resetpassword` (**public**)
Body: `{ "code": "...", "password": "..." }`

### GET `/auth/refresh` · GET `/auth/logout` · GET|POST `/auth/logoutall` (auth)

---

## 2. Users

| Method | Path | Access |
|--------|------|--------|
| POST | `/register` | public |
| POST | `/signin` | public |
| POST | `/login` | public |
| GET | `/user/me` | auth |
| PUT | `/user/me` | auth |
| PATCH | `/user/me/updatepass` | auth |
| POST | `/user/me/avatar` | auth |
| GET | `/user/profile/{id}` | auth |
| GET/POST | `/user/` | auth + superuser |
| PATCH | `/user/{id}/role` | auth + superuser |
| GET | `/user/list` | auth + superuser |
| GET | `/user/statistics` | auth + superuser |
| GET | `/user/notify/{channel}` | auth + superuser |
| PATCH | `/user/{id}/activestatus` | auth + superuser |
| GET/DELETE/PUT | `/user/{id}` | auth (+superuser) |
| PATCH | `/user/{id}/updatepass` | auth + superuser |
| GET | `/user/{id}/logoutall` | auth + superuser |

### POST `/register`
Body: `{ "username", "email", "password" }`

### GET `/user/me`
Profile ของ user ที่ login

---

## 3. Items

| Method | Path | Access |
|--------|------|--------|
| GET/POST | `/item/` | auth |
| GET/DELETE/PUT | `/item/{id}` | auth |

---

## 4. IoT Devices

> Rate-limited. Public GET ทั้งหมด ยกเว้น 4 route auth

### Public
`GET /iot/topic`, `/iot/topicdevicechart`, `/iot/controls`, `/iot/monitordevicegroup`, `/iot/monitordevicechart`, `/iot/device`, `/iot/devicebuckets`, `/iot/sensercharts`, `/iot/locationdevice`, `/iot/devicesensercharts`, `/iot/alarmdevicestatus`, `/iot/alarmdevicestatuscontrol`, `/iot/devicemqtt`, `/iot/devicestatus`, `/iot/deviceconfig`, `/iot/deviceiotdata`, `/iot/devicestats`, `/iot/devicedataexport`

### Auth
`POST /iot/control`, `PUT /iot/devicestatus`, `PUT /iot/updatedeviceconfig`, `DELETE /iot/devicedatacleanup`

### GET `/iot/devicemqtt` (ใหม่ — output ตรงกับ NestJS mqtt2/devicemqtt)
Query: `bucket`, `page`(1), `pageSize`(10000000), `lang`(en|th), `device_id`, `mqtt_id`, `type_id`, `hardware_id`, `keyword`, `deletecache`(0|1)
Response `{ code, payload, message, message_th }` — `payload` มี `timestamps`, `lang`, `page`, `currentPage`, `pageSize`, `totalPages`, `total`, `cache`, `device_count`, `connectionMqtt`, `device[]` (rss_data-style: `device_id`, `device_name`, `hardware_type`, `bucket`, `measurement`, `sensor_type`, `devicedata`, `alarm_title`, `alarm_detail`, `mqtt_data`, `mqtt_control`, `control`, `sensercharts`, `value_data`, `icon_access`, `alarm_subject`, `alarm_status` ฯลฯ)

### GET `/iot/alarmdevicestatus`
Query: `bucket` (required), `page`, `pageSize` (>0), `lang`, `measurement`
Response `{ statuscode, status, Mqttstatus, payload, message, message_th }` — `payload` มี `checkConnectionMqtt`, `mqttrs`, `mqttname`, `bucket`, `time`, `mqttdata`, `deviceioinfo`, `devicesensor`, `deviceio`, `devicecritical`, `cache`, `chart` + **mqtt2-compatible**: `lang`, `page`, `currentPage`, `pageSize`, `total`, `device_count`, `device[]`

---

## 5. MQTT Gateway

| Method | Path | Access |
|--------|------|--------|
| GET | `/mqtt/subscriptions` | public |
| GET | `/mqtt/status` | public |
| GET | `/mqtt/gettopicdata` | public |
| GET | `/mqtt/devicecontrol` | public |
| POST | `/mqtt/publish` | auth |
| POST | `/mqtt/subscribe` | auth |
| POST | `/mqtt/unsubscribe` | auth |

---

## 6. InfluxDB

| Method | Path |
|--------|------|
| POST | `/influx/write` |
| GET | `/influx/query` |
| POST | `/influx/devicechart` |
| POST | `/influx/filters` |
| POST | `/influx/statistics` |

---

## 7. Alarm Validation

| Method | Path |
|--------|------|
| POST | `/alarm/validate` |
| POST | `/alarm/validate/en` |
| POST | `/alarm/validate/th` |

---

## 8. Orders (Kafka)

| Method | Path |
|--------|------|
| POST | `/orders` |

## 9. WebSocket

| Method | Path | Access |
|--------|------|--------|
| GET | `/api/ws` | public |

## 10. System & Health

| Method | Path | Access |
|--------|------|--------|
| GET | `/` | public |
| GET | `/health` | public |
| GET | `/api/ping` | public |
| GET | `/api/health` | public |
| GET | `/metrics` | public (Prometheus) |
| GET | `/apimetric` | public |
| GET | `/uploads/avatar/*` | public (static) |
| GET | `/swagger/*` | public (Swagger UI) |

---

## 11. Settings — `/api/settings`

> Auth ทุก endpoint · รายละเอียดเต็มทุก endpoint: [`docs/apidocs/settings.md`](./apidocs/settings.md)

**Master:** Setting · Location · Type · DeviceType · Group · Sensor (CRUD: `GET ·/listX`, `GET ·/Xall`, `POST ·/createX`, `POST ·/updateX`, `GET ·/deleteX`)

**Device & Schedule:** `GET /settings/deviceall`, `/settings/listdevicepage`, `/settings/listdevicedevicepagess`, `/settings/listdevicepageactive1`, `/settings/listdevicepageactive`, `/settings/listdevicepageall`, `/settings/listdevicepagesensor`, `/settings/listdevicepageallactive`, `/settings/listdevicepageallactiveschedule`, `/settings/deviceeditget`, `/settings/devicedetail`, `/settings/devicedelete`, `/settings/deletedevice`, `POST /settings/createdevice`, `/settings/updatedevice`, `/settings/updatestatusdeviceid`, `/settings/deviceactionuser` · Schedule: `GET /settings/listscheduledevice`, `/settings/findscheduledevicechk`, `/settings/schedulelist`, `/settings/scheduleall`, `/settings/listschedulepage`, `/settings/scheduledevicepage`, `/settings/listdevicescheduledata`, `*device*`, `POST /settings/createschedule`, `/settings/createscheduledevice`, `/settings/updateschedule`, `/settings/updateschedulestatus`, `/settings/updatescheduledaystatus`, `GET /settings/deleteschedule`

**Integration:** MQTT · MQTTHost · API · Email · Host · InfluxDB · Line · NodeRED · SMS · Token · Telegram (มี `GET ·/listX`, `GET ·/Xall`, `POST ·/createX`, `POST ·/updateX`, `POST ·/updateXstatus`, `GET ·/deleteX`)

**Dashboard Config:** `POST /settings/dashboardconfig`, `GET /settings/dashboardconfig_1`, `GET /settings/dashboardconfig`, `GET /settings/dashboardconfig/search`, `GET/PATCH/DELETE /settings/dashboardconfig/{id}`

**Alarm Device & Logs:** `GET /settings/listalarmdevicepage`, `/settings/alarmdevicestatus`, `/settings/devicealarm`, `POST /settings/createalarmDevice`, `/settings/createalarmdevicepaginate`, `/settings/updatealarmdevice`, `/settings/updatealarmstatus` ฯลฯ · Monitors: `/settings/listdevicealarm*`, `/settings/devicemonitor*` · Logs: `/settings/scheduleprocess*`, `/settings/mqtterrorlogpaginate`, `/settings/alarmlogpaginate*`

---

## 12. Business Modules

> Auth ทุก endpoint · รายละเอียดเต็มทุก endpoint + request/response: [`docs/apidocs/business.md`](./apidocs/business.md)

### Quotation — `/api/quotation`
`GET|POST /quotation`, `GET|PUT|DELETE /quotation/{id}`, `GET /quotation/{id}/pdf` · snake_case

### Purchase Order — `/api/purchase-orders`
`GET|POST /purchase-orders`, `GET /purchase-orders/suggestions/{jobId}`, `POST /purchase-orders/from-quotation/{quotationId}`, `GET|PUT|DELETE /purchase-orders/{id}`, `POST /purchase-orders/{id}/send`, `PUT /purchase-orders/{id}/confirm`, `POST /purchase-orders/{id}/receive`, `PUT /purchase-orders/{id}/cancel`, `GET /purchase-orders/{id}/pdf|history` · snake_case

### Payment & Receipt — `/api/payments` `/api/receipts`
`POST|GET /payments`, `POST /payments/search`, `GET /payments/outstanding/{customerId}`, `GET /payments/history/{customerId}`, `GET /payments/{id}`, `POST /payments/{id}/refund`, `PUT /payments/{id}/cancel`, `GET /payments/invoice/{invoiceId}` · `GET /receipts/{id}`, `GET /receipts/{id}/pdf`, `PUT /receipts/{id}/cancel`, `GET /receipts/payment/{paymentId}` · snake_case

### Job — `/api/job`
`GET|POST /job`, `GET|PUT|DELETE /job/{id}`, `PUT /job/{id}/status`, `GET /job/{id}/history|report|pdf|picking/pdf|delivery/pdf`, `POST|GET /job/{id}/services`, `POST|GET /job/{id}/parts` · **camelCase** (`jobNo`, `customerId`...)

### Customer & Car — `/api/customer` `/api/car`
`GET|POST /customer`, `GET|PUT|DELETE /customer/{id}`, `GET|POST /car`, `GET|PUT|DELETE /car/{id}` · **camelCase**

---

## 13. Infrastructure & Utility

> ส่วนใหญ่ auth ยกเว้นที่ระบุ public · รายละเอียดเต็ม: [`docs/apidocs/infra.md`](./apidocs/infra.md)

**Dashboard `/api/dashboard`:** `GET /dashboard/stats`, `/dashboard/revenue`, `/dashboard/top-parts`, `/dashboard/job-status`
**Document `/api/documents`:** `POST|GET /documents/`, `GET|DELETE /documents/{id}`
**Email `/api/email`:** `POST /email/send`, `GET /email/logs`, `GET /email/logs/{id}`, `GET|PUT /email/config`
**Batch `/api/batch`:** `POST|GET /batch/jobs`, `GET|PUT|DELETE /batch/jobs/{id}`, `POST /batch/jobs/{id}/run`, `GET /batch/jobs/{id}/logs`
**I18n `/api/i18n`:** `GET|POST /i18n/translations`, `GET|PUT|DELETE /i18n/translations/{key}`
**WOS `/api/wos`:** `POST|GET /wos/orders`, `GET /wos/orders/{id}`, `PUT /wos/orders/{id}/status`
**Report `/reports`:** `GET /reports/daily-sales/pdf|inventory-summary/pdf|customer-list/pdf|invoice/pdf|credit-note/pdf|debit-note/pdf`, `POST /reports/token` (PDF รองรับ `?token=` short-lived)
**VectorData `/api/vectordata`:** `GET /vectordata/health` (**public**), `POST /vectordata/create-index`, `POST|GET /vectordata/documents`, `GET|PUT|DELETE /vectordata/documents/{id}`, `POST /vectordata/search`, `POST /vectordata/seed`
**Elasticsearch `/api/elasticsearch`:** `GET /elasticsearch/health` (**public**), `POST /elasticsearch/create-index`, `POST /elasticsearch/embeddings`, `POST /elasticsearch/vectors`, `POST /elasticsearch/search`

---

## 14. Common Response Envelope

**Success (ส่วนใหญ่):** `{ "is_success": true, "data": { ... } }`
**Error (auth, users, items, IoT, business, infra):** `{ "is_success": false, "error": { "status": 404, "statusText": "not_found", "msg": "..." } }`
**Error (alarm, influx, mqtt, vectordata, elasticsearch):** `{ "error": "error description" }`
**IoT devicemqtt:** `{ "code": 200, "payload": { ... }, "message": "OK", "message_th": "OK" }`

---

## 15. Authentication & Roles

```
Authorization: Bearer <access_token>
```
Tokens เป็น JWT (RS256) ได้จาก `POST /api/auth/login` หรือ `POST /api/auth/signin`

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

## ฉบับแยกตาม module
| ไฟล์ | ครอบคลุม |
|------|---------|
| [`docs/apidocs/auth-core.md`](./apidocs/auth-core.md) | Auth / Users / Items / Orders / WS / Health |
| [`docs/apidocs/iot-mqtt.md`](./apidocs/iot-mqtt.md) | IoT / MQTT / InfluxDB / Alarm |
| [`docs/apidocs/settings.md`](./apidocs/settings.md) | Settings (ทุกสาขา) |
| [`docs/apidocs/business.md`](./apidocs/business.md) | Quotation / PO / Payment / Receipt / Job / Customer / Car |
| [`docs/apidocs/infra.md`](./apidocs/infra.md) | Dashboard / Document / Email / Batch / I18n / WOS / Report / VectorData / Elasticsearch |
