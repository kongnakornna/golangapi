# README — Settings Module (`internal/modules/settings`)

> Port 1:1 จาก NestJS settings service → Go (chi + GORM + PostgreSQL) ด้วยสถาปัตยกรรม 3 ชั้น
> ทุก path/verb ของ endpoint ตรงกับ NestJS controller เดิมทั้งหมด (~185 routes) รวมถึง GET-deletes แบบ legacy

---

## 1. ภาพรวม

Settings module เป็นโมดูลที่ใหญ่ที่สุดของระบบ ดูแล master data, device, schedule,
integration configs (mqtt/email/line/sms/telegram/influxdb/nodered ฯลฯ), dashboard config,
alarm links, monitor views และ process logs

จุดเด่นของการออกแบบคือ **"spec-driven querying"** — แทนที่จะเขียน SQL ต่อ endpoint,
แต่ละ resource นิยาม query ไว้เป็น `ListSpec` (declarative) แล้ว handler/usecase/repository
ชุดเดียวกัน reuse กับทุก resource ผ่าน generic operations

## 2. โครงสร้างไฟล์

```
internal/modules/settings/
├── handler.go               # Interface `Handlers` — HTTP boundary (1 method ต่อ 1 endpoint)
├── usecase.go               # Interface `SettingsUseCaseI` — business boundary
├── repository.go            # Interface `SettingsRepositoryI` — storage boundary
├── types.go                 # Filter, ListSpec, ListQuery, ParseSort, EscapeLike
├── spec_with.go             # FixedFilter + ListSpec.With() (variant ของ spec)
├── specs_master.go          # Spec ฝั่ง master (setting/location/type/devicetype/group/sensor)
├── specs_device_schedule.go # Device (+active variants), Schedule, ScheduleDeviceJoin
├── specs_integration.go     # mqtt/mqtthost/api/email/host/influxdb/line/nodered/sms/token/telegram/dashboardconfig
├── specs_alarm_logs.go      # alarm link specs + alarm/schedule/mqtt log specs
├── types_test.go            # Unit test ของ ParseSort / EscapeLike / With()
├── presenter/
│   └── presenters.go        # Row, ListResult, SendEmailResult, MqttDataResult, Affected
├── usecase/
│   ├── usecase.go           # CreateSettingsUseCaseI — generic ops + special flows
│   ├── usecase_utility.go   # smtpVerifyAndSend, sendMail, sensorTypeRows
│   └── usecase_test.go
├── repository/
│   └── pg_repository.go     # GORM implementation ของ SettingsRepositoryI
└── delivery/http/
    ├── routes.go            # MapSettingRoute — ลงทะเบียนทุก endpoint ใต้ /api/settings
    ├── handlers.go          # settingsHandler struct + shared helpers (listQueryFrom, decodeBody, ...)
    ├── handler_helpers.go   # Generic handler factories (list/all/create/update/delete)
    ├── handlers_master.go       # setting/location/type/devicetype/group/sensor
    ├── handlers_device.go       # device CRUD + list variants
    ├── handlers_schedule.go     # schedule + schedule-device link
    ├── handlers_integration.go  # mqtt/mqtthost/api/email/host/influxdb/line/nodered/sms/token/telegram
    ├── handlers_dashboard.go    # dashboard config
    ├── handlers_alarm.go        # alarm device/event links + actions
    ├── handlers_monitor.go      # device alarm views / monitors
    ├── handlers_utility.go      # sensortype / testgemail / sendemail / mqttdata
    └── handlers_test.go         # handler tests + TestRoutesRegisterWithoutPanic
```

## 3. Request Flow

```
chi router (delivery/http/routes.go)
   │  middleware: Verifier(true) → Authenticator → CurrentUser → ActiveUser
   ▼
Handlers interface (handler.go)          ← HTTP boundary
   ▼  generic factories (handler_helpers.go)
SettingsUseCaseI (usecase/)              ← business boundary, logging, special flows
   ▼
SettingsRepositoryI (repository/)        ← storage boundary
   ▼
GORM → PostgreSQL
   ▼
presenter.ListResult / responses envelope → JSON
```

Dependency Injection อยู่ที่ `internal/server/handlers.go` (ช่วงบรรทัด ~328):

```go
settingsRepo := settingsRepoPkg.CreateSettingsPgRepository(db)
settingsUC := settingsUCPkg.CreateSettingsUseCaseI(settingsRepo, mqttClient, redisCache, logger)
settingsHandler := settingsHttp.CreateSettingsHandler(settingsUC, cfg, logger)
settingsHttp.MapSettingRoute(apiRouter, settingsHandler, mw)
```

## 4. ชั้น Contract (root package)

| ไฟล์ | สิ่งที่ประกาศ | หมายเหตุ |
|---|---|---|
| `handler.go` | `Handlers` | method ชื่อเดียวกับ path (PascalCase) คืน `http.HandlerFunc`; มี comment กำกับ route เดิม |
| `usecase.go` | `SettingsUseCaseI` | generic ops + special flows (`MqttData`, `SendEmailStub`, `TestGmailConnection`, `SensorTypes`) |
| `repository.go` | `SettingsRepositoryI` | row payload เป็น `map[string]interface{}` keyed ด้วย alias ของ SELECT |

### Query engine types (`types.go`, `spec_with.go`)

```go
type Filter struct { Column, Key, Op string }        // Op "" = equality, "LIKE" = ILIKE '%..%'
type FixedFilter struct { Column, Value string }     // WHERE แบบ unconditional

type ListSpec struct {
    Table        string
    Selects      []string
    Joins        []string
    Filters      []Filter        // optional — apply เมื่อ query param มีค่า
    FixedFilters []FixedFilter   // apply ทุกครั้ง
    SortCols     map[string]string // whitelist: param name → column
    DefaultSort  string
}

func (s ListSpec) With(fixed ...FixedFilter) ListSpec // คืน copy (ไม่ mutate base)
```

- `ListQuery` — page/pageSize/sort + raw query values; `.Value(key)` อ่านค่าแบบ trim และ nil-safe
- `ParseSort(sort, allowed, fallback)` — รับ `"field-ASC|DESC"` ตรวจกับ whitelist เท่านั้น (กัน SQL injection ฝั่ง ORDER BY) ไม่ผ่าน → `ErrInvalidSort`
- `EscapeLike(s)` — escape `\ % _` ก่อนใส่ pattern ILIKE

## 5. Specs Catalog (23 specs)

| ไฟล์ | Spec | ตาราง |
|---|---|---|
| specs_master.go | `SettingSpec` | sd_iot_setting (+JOIN location, type) |
| | `LocationSpec` | sd_iot_location |
| | `TypeSpec` | sd_iot_type (+JOIN group) |
| | `DeviceTypeSpec` | sd_iot_device_type |
| | `GroupSpec` | sd_iot_group |
| | `SensorSpec` | sd_iot_sensor (+JOIN setting/type/location) |
| specs_device_schedule.go | `DeviceSpec` | sd_iot_device (+JOIN type/location/mqtt/host) |
| | `DeviceActive1Spec`, `DeviceActiveSpec`, `DeviceAllActiveSpec`, `DeviceAllActiveSchedSpec` | variant ของ DeviceSpec ผ่าน `.With(status=1)` |
| | `ScheduleSpec` | sd_iot_schedule (+JOIN device) |
| | `ScheduleDeviceJoinSpec` | sd_iot_schedule_device (+JOIN schedule/device) |
| specs_integration.go | `MqttSpec` | sd_iot_mqtt |
| | `MqttHostSpec` | sd_mqtt_host |
| | `ApiSpec` | sd_iot_api |
| | `EmailSpec` | sd_iot_email |
| | `HostSpec` | sd_iot_host |
| | `InfluxdbSpec` | sd_iot_influxdb |
| | `LineSpec` | sd_iot_line |
| | `NoderedSpec` | sd_iot_nodered |
| | `SmsSpec` | sd_iot_sms |
| | `TokenSpec` | sd_iot_token |
| | `TelegramSpec` | sd_iot_telegram |
| | `DashboardConfigSpec` | sd_dashboard_config |
| specs_alarm_logs.go | `AlarmDeviceJoinSpec` / `AlarmEventDeviceJoinSpec` | sd_iot_alarm_device(_event) (+JOIN alarm_action/device) — สร้างจาก factory `alarmLinkSpec()` |
| | `AlarmProcessLogSpec` + Email/Line/Sms/Telegram variants | sd_alarm_process_log* — สร้างจาก factory `alarmLogSpec()` |
| | `ScheduleProcessLogSpec` | sd_schedule_process_log |
| | `MqttErrorLogSpec` | sd_mqtt_log |

## 6. Delivery Layer (`delivery/http`)

### routes.go — `MapSettingRoute(router chi.Router, h settings.Handlers, mw)`

ลงทะเบียนทั้งหมด **185 routes** ใต้ `/api/settings` พร้อม middleware JWT:

```go
r.Use(mw.Verifier(true)); r.Use(mw.Authenticator())
r.Use(mw.CurrentUser());  r.Use(mw.ActiveUser())
```

### handler_helpers.go — Generic Handler Factories

| Factory | ใช้ทำ | พฤติกรรม |
|---|---|---|
| `listHandler(spec)` | `GET /list*` | paginated → `presenter.ListResult` |
| `allHandler(spec)` | `GET /*all` | ไม่มี LIMIT/OFFSET (PageSize=0) |
| `createHandler(table, cols, defaults)` | `POST /create*` | filter body ด้วย whitelist cols + apply defaults (เช่น `status:1`) |
| `updateBodyHandler(table, idCol, cols)` | `POST /update*` | partial update + stamp `updateddate`, id จาก body |
| `deleteGetHandler(table, idCol, param)` | `GET|DELETE /delete*` | delete ตาม query param → `{affected:n}` |

Helpers อื่น: `listQueryFrom` (parse page/pageSize/sort/filters), `decodeBody`, `filterKeys`,
`coerceID` (string ตัวเลข → int ให้ PG เทียบ int column), response helpers
(`ok/fail/badRequest/notFound/affected`) ซึ่ง wrap ด้วย `pkg/responses`

### handlers_*.go — Pattern ต่อ domain

- ประกาศ **column whitelist** ต่อตาราง (มาจาก internal/models เดิม) เช่น `settingCols`
- Method ละ 1 endpoint → เรียก factory + spec/table ที่เหมาะสม
- มี Swagger annotation (swag) ครบทุก method: `@Summary/@Param/@Success/@Router`
- Endpoint variants ที่ต่างกันแค่ fixed filter ใช้ spec variant เดียวกัน เช่น
  `ListDevicePageActive()` → `listHandler(settings.DeviceActiveSpec)`

### Handler tests (`handlers_test.go`)

- ทดสอบ helper functions (`filterKeys`, `coerceID`, `toStr`)
- `TestRoutesRegisterWithoutPanic` — ลงทะเบียน `MapSettingRoute` บน chi router จริง เพื่อ verify ครบ/ไม่ panic

## 7. Usecase Layer (`usecase/`)

`CreateSettingsUseCaseI(repo, mqttClient, cache, logger)` คืน `settings.SettingsUseCaseI`

**Generic operations** — delegate ไป repo + log error ทุก failure path:
`ListPaginate`, `RowsAll`, `CreateRow`, `GetWhere`, `ExistsWhere`, `UpdateFields`, `DeleteRow`

**Special flows:**

| Flow | Route | รายละเอียด |
|---|---|---|
| `MqttData` | GET /mqttdata | Redis cache key `get_device_data_ALL<topic>` (TTL 30s) → miss แล้วรอ MQTT topic 10s; payload ไม่ใช่ JSON จะ wrap เป็น `{"raw": ...}`; response บอกแหล่งข้อมูล `getdataFrom: Cache \| MQTT` |
| `TestGmailConnection` | GET /testgemail | ทดสอบ smtp.gmail.com port 465 (implicit TLS) แล้ว fallback 587 (STARTTLS); credential จาก env `GMAIL_USERNAME/GMAIL_APP_PASSWORD/GMAIL_TEST_TO` |
| `SendEmailStub` | GET /sendemail | parity กับ NestJS (ฝั่งเดิม nodemailer ถูก comment) — คืน success shape เดิม ไม่ส่งจริง |
| `SensorTypes` | GET /sensortype | static EN/TH list (port iot.helper) |

`usecase_utility.go`: `smtpVerifyAndSend` (dial+auth+sendMail ตาม port), `sendMail`
(สร้าง MIME message เอง), `sensorTypeRows` (ข้อมูลประเภท sensor EN/TH)

Tests: mock repository + fake cache/logger ใน `usecase_test.go`

## 8. Repository Layer (`repository/pg_repository.go`)

`SettingsPgRepo{DB *gorm.DB}` — implement ทุก method ของ `SettingsRepositoryI` แบบ generic:

- `listBase()` — ประกอบ Table + Joins + optional Filters (ILIKE + EscapeLike สำหรับ LIKE) + FixedFilters (`=`)
- `orderOf()` — validate sort ผ่าน `ParseSort` whitelist; ไม่ผ่าน → `ErrInvalidSort`
- `ListPaginate()` — clamp `page≥1`, `pageSize∈[1,5000]` (default 1000) → COUNT + SELECT LIMIT/OFFSET
- `RowsAll()` — ไม่มี LIMIT/OFFSET เมื่อ `PageSize<=0`
- `CreateMap / GetWhere / ExistsWhere / UpdateFieldsMap / DeleteWhere` — CRUD แบบ map-driven
  (ทุก query bind ผ่าน placeholder `?` เสมอ)

## 9. Presenter (`presenter/presenters.go`)

| Type | JSON | ใช้กับ |
|---|---|---|
| `Row` | object | row จาก list/all endpoints (free-form) |
| `ListResult` | `items,page,pageSize,total,totalPages` | ทุก paginated endpoint (normalize nil items, min page/pageSize, คำนวณ totalPages) |
| `SendEmailResult` | `success,code,error?,to,subject,content,message,message_th` | /sendemail, /testgemail |
| `MqttDataResult` | `getdataFrom,payload` | /mqttdata |
| `Affected` | `affected` | create/update/delete counts |

## 10. กลุ่ม Endpoint (สรุป 185 routes)

| กลุ่ม | จำนวน | ตัวอย่าง |
|---|---|---|
| Master: setting/location/type/devicetype/group/sensor | 33 | `/listsetting`, `/createsensor`, `/deletegroup` |
| Utility | 4 | `/sensortype`, `/testgemail`, `/sendemail`, `/mqttdata` |
| Device | 17 | `/deviceall`, `/listdevicepage*` (8 variants), `/createdevice`, `/deviceactionuser` |
| Schedule | 17 | `/scheduleall`, `/listschedulepage`, `/createscheduledevice`, `/updatescheduledaystatus` |
| Integration (mqtt, mqtthost, api, email, host, influxdb, line, nodered, sms, token, telegram) | 66 | `/lismqtt`, `/createmqtthost`, `/updateinfluxdbstatus` |
| Dashboard config | 7 | RESTful: POST/GET `/dashboardconfig`, `/search`, `/{id}`, PATCH/DELETE `/{id}` |
| Alarm device/event links + actions | 20 | `/listalarmdevicepage`, `/createalarmDevice`, `/createdevicealarmaction` |
| Device alarm views / monitors | 10 | `/listdevicealarmair`, `/devicemonitor(s)`, legacy `/_listdevicealarmfan`, `/_devicemonitor` |
| Process logs | 11 | `/scheduleprocesslogpaginate`, `/alarmlogpaginate{,email,line,sms,telegram,controls,control}` |

หมายเหตุ parity กับ NestJS:
- Delete ส่วนใหญ่เป็น **GET** ตาม controller เดิม (บาง route มีทั้ง GET + DELETE เช่น `/deletesetting`)
- `/devicetypeall` ประกาศซ้ำใน NestJS → register ครั้งเดียว
- Route static (`/dashboardconfig/search`) ต้องมาก่อน `/{id}`
- Path สะกดตามต้นฉบับแม้ typo เช่น `/lisgroup`, `/testgemail`, `/mqtttsort`

## 11. Convention ของ Query String (list endpoints)

```
?page=1&pageSize=1000&sort=createddate-DESC&keyword=xxx&status=1&location_id=3 ...
```

- `page` (default 1), `pageSize` (default 1000, max 5000)
- `sort=<field>-ASC|DESC` — field ต้องอยู่ใน `SortCols` whitelist ของ spec เท่านั้น
- `keyword` → ILIKE `%escaped%`; filter อื่น → equality
- Filter ที่ไม่ส่งค่ามาจะถูกข้าม (optional WHERE)

Response envelope มาตรฐานของ repo: success ผ่าน `responses.CreateSuccessResponse(data)`,
error ผ่าน `responses.CreateErrorResponse` (400/404/422/500 ตามกรณี)

## 12. วิธีเพิ่ม Resource ใหม่

1. เพิ่ม `XxxSpec = ListSpec{...}` ในไฟล์ `specs_*.go` ที่เหมาะสม (table/selects/joins/filters/sortCols/defaultSort)
2. เพิ่ม methods ใน interface `Handlers` (`handler.go`) — ชื่อ method = path แบบ PascalCase
3. ใน `delivery/http/handlers_<domain>.go`: ประกาศ column whitelist แล้ว implement ด้วย factory
   (`h.listHandler(spec)`, `h.createHandler(table, cols, defaults)` ฯลฯ) + Swagger annotation
4. ลงทะเบียน route ใน `routes.go` (ระวังลำดับ static ก่อน `/{id}`)
5. เพิ่ม test (types/handler/usecase) แล้วรัน `go test ./internal/modules/settings/...`
