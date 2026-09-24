# Plan: แปลง NestJS Settings Module → Go (`internal/modules/settings`)

## ภาพรวมการแก้ไข

| # | Layer | ไฟล์ | สิ่งที่ทำ |
|---|-------|------|-----------|
| 0 | convention | `internal/modules/settings/*` | กติกากลาง: envelope, error map, query parse, sort whitelist, auth group |
| 1 | presenter | `internal/modules/settings/presenter/presenters.go` | request DTO + `ListResult` meta wrapper |
| 2 | repository | `internal/modules/settings/repository/pg_repository.go` | struct + generic `ListPaginate` / CRUD helpers / sort whitelist (โค้ดเต็มใน Section 2) |
| 3 | repository | `repository/pg_repository_{master,device,schedule,integration,alarm}.go` | spec ต่อ resource — port SQL จาก function ตามตาราง inventory |
| 4 | interface | `settings/handler.go`, `settings/usecase.go`, `settings/repository.go` | interface กลาง module root package |
| 5 | usecase | `usecase/usecase.go` + `_master/_device/_schedule/_integration/_alarm/_utility` | validate → repo → return; special flows (mqttdata/sendemail) |
| 6 | delivery | `delivery/http/handlers.go` + `handlers_*.go` | decode → validate → usecase → envelope; ~200 handlers ตามตาราง |
| 7 | delivery | `delivery/http/routes.go` | `MapSettingRoute` — path+verb ตรง NestJS ทุก endpoint, JWT middleware ทั้งหมด |
| 8 | wiring | `internal/server/handlers.go` | สร้าง repo/uc/handler + mount `/api/settings` |
| 9 | test | `*_test.go` ใน settings module | unit tests (decision: **เขียน**) |

---

## เป้าหมาย

Port settings module จาก NestJS (`C:\github\gistdaapi\src\modules\settings`) → Go module ใหม่ตาม pattern ของ `icmongolang`:

- **1:1 ทุก endpoint** (~200) — path + HTTP verb เหมือนเดิมทั้งหมด (รวม GET-based delete/create)
- **Response envelope ใหม่** — ใช้ `pkg/responses` มาตรฐานของโปรเจกต์ (`{is_success, data}` / `{is_success:false, error:{...}}`), HTTP status จริง
- **Auth ทุก route** — JWT middleware (`Verifier/Authenticator/CurrentUser/ActiveUser`) ภายใต้ `/api/settings`
- **Models ไม่สร้างใหม่** — ใช้ `internal/models` เดิม (`SdIotSetting`, `SdIotLocation`, ..., `SdDashboardConfig`, `SdScheduleProcessLog`, `SdAlarmProcessLog*`)
- **Unit tests: เขียน** (Section Unit Tests)

### Deviations ที่ตกลงกันแล้ว (บันทึกชัดเจน)

| เรื่อง | NestJS เดิม | Go ใหม่ |
|--------|-------------|---------|
| Response shape | `{statusCode,code,payload,message,message_th}` HTTP 200 เสมอ | envelope `pkg/responses` + HTTP status จริง |
| Auth | comment ทิ้ง = public | JWT ทุก route |
| List ว่าง | `code:400, payload:null` | 200 + `items:[], total:0` |
| `GET devicetypeall` ประกาศซ้ำ 2 รอบ | register ซ้ำ | register ครั้งเดียว (chi panic ถ้าซ้ำ) |
| `GET listsensor` (controller เรียก `location_list_paginate` ตอน query data = bug copy-paste) | ได้ข้อมูล location ผิด type | ใช้ sensor query (`sensor_list_paginate` SQL) — ตั้งใจแก้ |
| `sendemail` / `testgemail` | service จริงถูก comment / hardcode credential ในไฟล์ | `testgemail`: พยายาม SMTP verify+send จริงด้วย `net/smtp` อ่าน credential จาก env (`GMAIL_USERNAME`, `GMAIL_APP_PASSWORD`); `sendemail`: return success payload ตรงต้นฉบับ (ต้นฉบับไม่ได้ส่งจริง) |
| createddate/updateddate format | `convertTZ + timeConvertermas` string เดิม | `time.Time` marshal JSON = RFC3339 อัตโนมัติ (ไม่ทำ converter) |

### แหล่งอ้างอิงต้นฉบับ

- Controller: `C:\github\gistdaapi\src\modules\settings\settings.controller.ts` (27,601 บรรทัด)
- Service: `C:\github\gistdaapi\src\modules\settings\settings.service.ts` (36,186 บรรทัด)
- DTO: `...\settings\dto\*.ts` (27 ไฟล์) — field ของ create/update body
- Helper: `@src/helpers/iot.helper.ts` (`list_sensor_type` / `list_sensor_type_th`), `@src/helpers/format.helper.ts` (`convertSortInput`)

> **กฎการตั้งชื่อ method:** ชื่อ handler/usecase method = PascalCase ของ route path (เช่น `/listsetting` → `ListSetting()`, `/listdevicepageactive1` → `ListDevicePageActive1()`) กรณี path เดียวกันมี 2 verb → suffix `ViaGet` / ตัวหลักเป็น POST (เช่น `CreateScheduleDeviceViaGet()` และ `CreateScheduleDevice()`) ตารางใน Section 6 ระบุครบทุกแถว

---

## Section 0: Conventions กลาง

**Location:** ทุกไฟล์ใหม่ใน `internal/modules/settings/`

### 0.1 Response envelope

```go
// success
render.Respond(w, r, responses.CreateSuccessResponse(data))
// error
render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusBadRequest, "msg")))
```

### 0.2 Error mapping

| Case | HTTP | วิธี |
|------|------|------|
| body/query decode fail, validate fail, sort whitelist reject, param แปลงไม่ได้ | 400 | `httpErrors.NewError(400, msg)` |
| id ไม่พบ (get/update/delete) | 404 | `httpErrors.NewError(404, "...not found")` |
| DB error / MQTT error | 500 | ส่ง error เดิมผ่าน `CreateErrorResponse` (log full error, body generic) |
| list ว่าง | 200 | `items: [], total: 0` |

### 0.3 Query parse + pagination

ทุก GET list: `page` (default 1), `pageSize` (default 1000, cap 5000), `sort=field-ASC|DESC`, `keyword`, และ filter เฉพาะ resource — parse ครั้งเดียวด้วย `repo.NewListQuery(r)` (Section 2)

### 0.4 Sort whitelist

field ไม่อยู่ใน `spec.SortCols` และ `sort != ""` → 400 `invalid sort option` (แทน `convertSortInput`)

### 0.5 Auth group

```go
router.Route("/settings", func(r chi.Router) {
    r.Use(mw.Verifier(true))
    r.Use(mw.Authenticator())
    r.Use(mw.CurrentUser())
    r.Use(mw.ActiveUser())
    // ... routes ทั้งหมด
})
```

### 0.6 Partial update

update ทุก resource: build `map[string]interface{}` เฉพาะ field ที่ส่งค่ามา (non-zero) + `updateddate: time.Now()` แล้ว `UpdateFields(...)` — เหมือน pattern `if (dto.field)` ของ NestJS

---

## Section 1: presenter

**Location:** `internal/modules/settings/presenter/presenters.go` (ไฟล์ใหม่)

**From → To:** ไม่มี presenter settings → สร้างใหม่

**Change detail:**

```go
package presenter

// ---- list meta ----
type ListResult struct {
	Items      []map[string]interface{} `json:"items"`
	Page       int                      `json:"page"`
	PageSize   int                      `json:"pageSize"`
	Total      int64                    `json:"total"`
	TotalPages int64                    `json:"totalPages"`
}

func NewListResult(items []map[string]interface{}, total int64, page, pageSize int) *ListResult {
	tp := total / int64(pageSize)
	if total%int64(pageSize) > 0 {
		tp++
	}
	if tp < 1 {
		tp = 1
	}
	return &ListResult{Items: items, Page: page, PageSize: pageSize, Total: total, TotalPages: tp}
}

// ---- setting ----
type SettingCreateRequest struct {
	LocationID    int    `json:"location_id" validate:"required"`
	SettingTypeID int    `json:"setting_type_id" validate:"required"`
	SettingName   string `json:"setting_name" validate:"required"`
	SN            string `json:"sn" validate:"required"`
	Status        int    `json:"status"`
}

type SettingUpdateRequest struct {
	SettingID     *int    `json:"setting_id"`
	LocationID    *int    `json:"location_id"`
	SettingTypeID *int    `json:"setting_type_id"`
	SettingName   *string `json:"setting_name"`
	SN            *string `json:"sn"`
	Status        *int    `json:"status"`
}

// ---- utility ----
type SendEmailResult struct {
	Success   bool   `json:"success"`
	Code      int    `json:"code"`
	Error     string `json:"error,omitempty"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	Content   string `json:"content"`
	Message   string `json:"message"`
	MessageTh string `json:"message_th"`
}

type MqttDataResult struct {
	GetDataFrom string                   `json:"getdataFrom"`
	Payload     []map[string]interface{} `json:"payload"`
}
```

**Request structs ที่เหลือ** (create/update ต่อ resource): field = ชุดที่ DTO ฝั่ง NestJS ประกาศ (`dto/create-setting.dto.ts`, `create-location.dto.ts` = `location_name, ipaddress, location_detail, configdata, status`, `create_type.dto.ts` = `type_name, group_id, status`, `create_sensor.dto.ts`, `create_group.dto.ts`, `create-device.dto.ts`, `create_mqtt.dto.ts`, ...) map ตรง column ของ model ใน `internal/models/models.go` pointer-type (`*string/*int/*float64`) สำหรับ update เพื่อ partial update

---

## Section 2: repository — generic core

**Location:** `internal/modules/settings/repository/pg_repository.go` (ไฟล์ใหม่)

**From → To:** ไม่มี → generic `ListPaginate` + CRUD helpers ที่ทุก domain ใช้ร่วมกัน

**Change detail (โค้ดเต็ม):**

```go
package repository

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

var ErrInvalidSort = errors.New("invalid sort option")

type Filter struct {
	Column string // "s.setting_name"
	Key    string // ชื่อ query param: "keyword", "status", "setting_id", ...
	Op     string // "" (= ) หรือ "LIKE"
}

type ListSpec struct {
	Table       string            // "sd_iot_setting s"
	Selects     []string          // raw columns รวม column จาก join
	Joins       []string          // "LEFT JOIN sd_iot_location l ON l.location_id = s.location_id"
	Filters     []Filter
	SortCols    map[string]string // whitelist: "createddate" -> "s.createddate"
	DefaultSort string            // "s.setting_id ASC"
}

type ListQuery struct {
	Page     int
	PageSize int
	Sort     string
	values   map[string]string
}

func NewListQuery(page, pageSize int, sort string, values map[string]string) ListQuery {
	return ListQuery{Page: page, PageSize: pageSize, Sort: sort, values: values}
}

func (q ListQuery) Value(key string) string {
	if q.values == nil {
		return ""
	}
	return strings.TrimSpace(q.values[key])
}

func EscapeLike(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return r.Replace(s)
}

// ParseSort แปลง "createddate-DESC" -> "s.createddate DESC" ตาม whitelist;
// sort ว่าง -> fallback; field ไม่อยู่ใน whitelist -> ("" , false)
func ParseSort(sort string, allowed map[string]string, fallback string) (string, bool) {
	if strings.TrimSpace(sort) == "" {
		return fallback, true
	}
	parts := strings.SplitN(sort, "-", 2)
	if len(parts) != 2 {
		return "", false
	}
	col, ok := allowed[parts[0]]
	if !ok {
		return "", false
	}
	dir := strings.ToUpper(parts[1])
	if dir != "ASC" && dir != "DESC" {
		return "", false
	}
	return col + " " + dir, true
}

type SettingsPgRepo struct {
	DB *gorm.DB
}

func CreateSettingsPgRepository(db *gorm.DB) *SettingsPgRepo {
	return &SettingsPgRepo{DB: db}
}

func (r *SettingsPgRepo) listBase(ctx context.Context, spec ListSpec, q ListQuery) *gorm.DB {
	db := r.DB.WithContext(ctx).Table(spec.Table)
	for _, j := range spec.Joins {
		db = db.Joins(j)
	}
	for _, f := range spec.Filters {
		v := q.Value(f.Key)
		if v == "" {
			continue
		}
		if f.Op == "LIKE" {
			db = db.Where(f.Column+" ILIKE ?", "%"+EscapeLike(v)+"%")
		} else {
			db = db.Where(f.Column+" = ?", v)
		}
	}
	return db
}

func (r *SettingsPgRepo) ListPaginate(ctx context.Context, spec ListSpec, q ListQuery) ([]map[string]interface{}, int64, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 5000 {
		q.PageSize = 1000
	}
	order := spec.DefaultSort
	if q.Sort != "" {
		o, ok := ParseSort(q.Sort, spec.SortCols, spec.DefaultSort)
		if !ok {
			return nil, 0, ErrInvalidSort
		}
		order = o
	}

	var total int64
	if err := r.listBase(ctx, spec, q).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	items := []map[string]interface{}{}
	err := r.listBase(ctx, spec, q).
		Select(strings.Join(spec.Selects, ", ")).
		Order(order).
		Limit(q.PageSize).
		Offset((q.Page - 1) * q.PageSize).
		Find(&items).Error
	return items, total, err
}

// ---- CRUD helpers ----

func (r *SettingsPgRepo) CreateModel(ctx context.Context, model interface{}) error {
	return r.DB.WithContext(ctx).Create(model).Error
}

func (r *SettingsPgRepo) GetWhere(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return r.DB.WithContext(ctx).Where(query, args...).First(dest).Error
}

func (r *SettingsPgRepo) ExistsWhere(ctx context.Context, model interface{}, query string, args ...interface{}) (bool, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(model).Where(query, args...).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *SettingsPgRepo) UpdateFields(ctx context.Context, model interface{}, idCol string, idVal interface{}, fields map[string]interface{}) error {
	return r.DB.WithContext(ctx).Model(model).Where(idCol+" = ?", idVal).Updates(fields).Error
}

func (r *SettingsPgRepo) DeleteWhere(ctx context.Context, model interface{}, query string, args ...interface{}) (int64, error) {
	res := r.DB.WithContext(ctx).Where(query, args...).Delete(model)
	return res.RowsAffected, res.Error
}
```

---

## Section 3: repository — per-domain files + specs

**Location:** `repository/pg_repository_master.go`, `pg_repository_device.go`, `pg_repository_schedule.go`, `pg_repository_integration.go`, `pg_repository_alarm.go`

**From → To:** ไม่มี → method list ต่อ resource ที่ wrap `ListPaginate` ด้วย `ListSpec`

**Pattern (ตัวอย่าง setting — ที่เหลือทำตามตาราง inventory):**

```go
var settingSpec = ListSpec{
	Table: "sd_iot_setting s",
	Selects: []string{
		"s.setting_id AS setting_id", "s.location_id AS location_id",
		"s.setting_type_id AS setting_type_id", "s.setting_name AS setting_name",
		"s.sn AS sn", "l.location_name AS location_name", "t.type_name AS type_name",
		"s.status AS status", "s.createddate AS createddate", "s.updateddate AS updateddate",
	},
	Joins: []string{
		"LEFT JOIN sd_iot_location l ON l.location_id = s.location_id",
		"LEFT JOIN sd_iot_type t ON t.type_id = s.setting_type_id",
	},
	Filters: []Filter{
		{Column: "s.setting_name", Key: "keyword", Op: "LIKE"},
		{Column: "s.setting_id", Key: "setting_id"},
		{Column: "s.location_id", Key: "location_id"},
		{Column: "s.setting_type_id", Key: "setting_type_id"},
		{Column: "s.sn", Key: "sn"},
		{Column: "s.status", Key: "status"},
	},
	SortCols:    map[string]string{"setting_id": "s.setting_id", "setting_name": "s.setting_name", "createddate": "s.createddate"},
	DefaultSort: "s.setting_id ASC",
}

func (r *SettingsPgRepo) SettingListPaginate(ctx context.Context, q ListQuery) ([]map[string]interface{}, int64, error) {
	return r.ListPaginate(ctx, settingSpec, q)
}

func (r *SettingsPgRepo) SettingAll(ctx context.Context) ([]map[string]interface{}, error) {
	rows := []map[string]interface{}{}
	err := r.DB.WithContext(ctx).Table("sd_iot_setting").Find(&rows).Error
	return rows, err
}
```

**กฎการ port แต่ละ resource:** เปิด function ต้นฉบับตามชื่อในตาราง Section 6 → คัด `selects / joins / filters / sort whitelist` จาก QueryBuilder → เขียนเป็น `ListSpec`; create/update/delete ใช้ helpers จาก Section 2 + model ที่ mapping อยู่แล้วใน `internal/models/models.go`:

| Resource | Model | Table |
|---|---|---|
| setting | `models.SdIotSetting` | sd_iot_setting |
| location | `models.SdIotLocation` | sd_iot_location |
| type | `models.SdIotType` | sd_iot_type |
| devicetype | `models.SdIotDeviceType` | sd_iot_device_type |
| group | `models.SdIotGroup` | sd_iot_group |
| sensor | `models.SdIotSensor` | sd_iot_sensor |
| device | `models.SdIotDevice` | sd_iot_device |
| deviceaction/user/log/alarmaction | `SdIotDeviceAction` / `SdIotDeviceActionUser` / `SdIotDeviceActionLog` / `SdIotDeviceAlarmAction` | ตาม TableName() |
| schedule / scheduledevice | `SdIotSchedule` / `SdIotScheduleDevice` | sd_iot_schedule / sd_iot_schedule_device |
| mqtt / mqtthost | `SdIotMqtt` / `SdMqttHost` | sd_iot_mqtt / mqtthost table |
| api/email/host/influxdb/line/nodered/sms/token/telegram | `SdIotApi` / `SdIotEmail` / `SdIotHost` / `SdIotInfluxdb` / `SdIotLine` / `SdIotNodered` / `SdIotSms` / `SdIotToken` / `SdIotTelegram` | ตาม TableName() |
| alarmdevice / alarmevent | `SdIotAlarmDevice` / `SdIotAlarmDeviceEvent` | ตาม TableName() |
| logs | `SdScheduleProcessLog` / `SdAlarmProcessLog` + `...Email/Line/Mqtt/Sms/Telegram/Temp` / `SdMqttLog` | ตาม TableName() |
| dashboardconfig | `SdDashboardConfig` | ตาม TableName() |

---

## Section 4: interfaces + usecase

**Location:** `settings/usecase.go` + `settings/repository.go` (module root) + `usecase/*.go`

**From → To:** ไม่มี → interface กลาง + implementation

**Change detail:**

```go
// settings/repository.go
package settings

type SettingsRepositoryI interface {
	// generic core
	ListPaginate(ctx context.Context, spec ListSpecAlias, q ListQueryAlias) ([]map[string]interface{}, int64, error)
	CreateModel(ctx context.Context, model interface{}) error
	GetWhere(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExistsWhere(ctx context.Context, model interface{}, query string, args ...interface{}) (bool, error)
	UpdateFields(ctx context.Context, model interface{}, idCol string, idVal interface{}, fields map[string]interface{}) error
	DeleteWhere(ctx context.Context, model interface{}, query string, args ...interface{}) (int64, error)

	// typed list wrappers — ชุด method = ทุกแถว "list/all/paginate" ในตาราง Section 6
	// เช่น SettingListPaginate, LocationListPaginate, DeviceListPaginateAll, ...
}
```

> alias types `ListSpecAlias/ListQueryAlias` = type alias ชี้ `repository.ListSpec/ListQuery` (import cycle ไม่เกิดเพราะ repository package import settings ไม่ได้ — ย้าย `Filter/ListSpec/ListQuery/ParseSort/EscapeLike` ไปนิยามที่ `settings/types.go` แล้ว `repository` import `settings` แทน)

```go
// settings/usecase.go
package settings

type SettingsUseCaseI interface {
	// 1 method ต่อ 1 endpoint — ชื่อตามกฎ PascalCase(path) เดียวกับ Handlers
	// signature กลาง: func(ctx context.Context, ...) (result, error)
	ListSetting(ctx context.Context, q ListQueryAlias) (*presenter.ListResult, error)
	CreateSetting(ctx context.Context, req *presenter.SettingCreateRequest) (*models.SdIotSetting, error)
	DeleteSetting(ctx context.Context, settingID int64) error
	MqttData(ctx context.Context, topic string) (*presenter.MqttDataResult, error)
	SendEmailStub(ctx context.Context, to, subject, content string) *presenter.SendEmailResult
	TestGmailConnection(ctx context.Context) (*presenter.SendEmailResult, error)
	// ... ตามตาราง Section 6 ครบทุกแถว
}
```

**usecase struct + constructor (`usecase/usecase.go`):**

```go
type settingsUseCase struct {
	repo   settings.SettingsRepositoryI
	mqtt   mqtt.Client
	redis  *redis.Client
	cfg    *config.Config
	logger logger.Logger
}

func CreateSettingsUseCaseI(repo settings.SettingsRepositoryI, mqttClient mqtt.Client, redisClient *redis.Client, cfg *config.Config, log logger.Logger) settings.SettingsUseCaseI {
	return &settingsUseCase{repo: repo, mqtt: mqttClient, redis: redisClient, cfg: cfg, logger: log}
}
```

**Worked example — usecase (setting):**

```go
func (uc *settingsUseCase) ListSetting(ctx context.Context, q settings.ListQueryAlias) (*presenter.ListResult, error) {
	items, total, err := uc.repo.SettingListPaginate(ctx, q)
	if err != nil {
		return nil, err
	}
	return presenter.NewListResult(items, total, q.Page, q.PageSize), nil
}

func (uc *settingsUseCase) CreateSetting(ctx context.Context, req *presenter.SettingCreateRequest) (*models.SdIotSetting, error) {
	m := &models.SdIotSetting{
		LocationID:    req.LocationID,
		SettingTypeID: req.SettingTypeID,
		SettingName:   req.SettingName,
		SN:            req.SN,
		Status:        req.Status,
	}
	if m.Status == 0 {
		m.Status = 1
	}
	if err := uc.repo.CreateModel(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (uc *settingsUseCase) DeleteSetting(ctx context.Context, settingID int64) error {
	ok, err := uc.repo.ExistsWhere(ctx, &models.SdIotSetting{}, "setting_id = ?", settingID)
	if err != nil {
		return err
	}
	if !ok {
		return httpErrors.NewError(http.StatusNotFound, fmt.Sprintf("sd_iot_setting with setting_id %d not found", settingID))
	}
	_, err = uc.repo.DeleteWhere(ctx, &models.SdIotSetting{}, "setting_id = ?", settingID)
	return err
}
```

**Special — MqttData (`usecase/usecase_utility.go`):**

```go
func (uc *settingsUseCase) MqttData(ctx context.Context, topic string) (*presenter.MqttDataResult, error) {
	topic = strings.TrimSpace(topic)
	if topic == "" {
		return nil, httpErrors.NewError(http.StatusUnprocessableEntity, "mqttdata is null.")
	}
	cacheKey := "get_device_data_ALL" + topic
	if raw, err := uc.redis.Get(ctx, cacheKey).Bytes(); err == nil && len(raw) > 0 {
		var cached []map[string]interface{}
		if json.Unmarshal(raw, &cached) == nil {
			return &presenter.MqttDataResult{GetDataFrom: "Cache", Payload: cached}, nil
		}
	} else if err != nil && !errors.Is(err, redis.Nil) {
		uc.logger.Warnf("settings MqttData cache get: %v", err)
	}

	rawMQTT, err := uc.mqtt.GetDataFromTopic(ctx, topic, 10*time.Second)
	if err != nil {
		return nil, err
	}
	var payload []map[string]interface{}
	if len(rawMQTT) > 0 {
		_ = json.Unmarshal(rawMQTT, &payload)
	}
	_ = uc.redis.Set(ctx, cacheKey, payload, 30*time.Second)
	return &presenter.MqttDataResult{GetDataFrom: "MQTT", Payload: payload}, nil
}
```

**Special — TestGmailConnection / SendEmailStub (`usecase/usecase_utility.go`):**

```go
func (uc *settingsUseCase) SendEmailStub(_ context.Context, to, subject, content string) *presenter.SendEmailResult {
	if strings.TrimSpace(to) == "" {
		to = "cmoniots@gmail.com"
	}
	if subject == "" {
		subject = "CmonIoT test send email"
	}
	if content == "" {
		content = "test send email"
	}
	return &presenter.SendEmailResult{
		Success: true, Code: 200, To: to, Subject: subject, Content: content,
		Message: "finally send email", MessageTh: "finally ส่ง email",
	}
}

// TestGmailConnection — port จาก testGmailConnection (service): ลอง verify SMTP
// smtp.gmail.com พอร์ต 465 (TLS) แล้ว 587 (STARTTLS) ด้วย env GMAIL_USERNAME/GMAIL_APP_PASSWORD
// สำเร็จ → ส่งเมล์ทดสอบหา GMAIL_TEST_TO (default icmon0955@gmail.com) คืน success payload
// env ไม่มี / verify fail → คืน success=false + message (HTTP ยัง 200 เหมือนต้นฉบับ)
func (uc *settingsUseCase) TestGmailConnection(ctx context.Context) (*presenter.SendEmailResult, error) { ... }
```

implement ด้วย `net/smtp` + `crypto/tls` (ไม่เพิ่ม dependency), timeout 30s/context

---

## Section 5: delivery/http handlers

**Location:** `delivery/http/handlers.go` (struct+constructor) + `handlers_master.go`, `handlers_device.go`, `handlers_schedule.go`, `handlers_integration.go`, `handlers_alarm.go`, `handlers_utility.go`

**Change detail — struct + worked example:**

```go
package http

type settingsHandler struct {
	cfg    *config.Config
	uc     settings.SettingsUseCaseI
	logger logger.Logger
}

func CreateSettingsHandler(uc settings.SettingsUseCaseI, cfg *config.Config, logger logger.Logger) settings.Handlers {
	return &settingsHandler{cfg: cfg, uc: uc, logger: logger}
}

// helper กลาง — parse query params ทั้งหมดเป็น ListQuery
func listQueryFrom(r *http.Request) settings.ListQueryAlias {
	q := r.URL.Query()
	vals := map[string]string{}
	for k := range q {
		vals[k] = q.Get(k)
	}
	page, _ := strconv.Atoi(q.Get("page"))
	size, _ := strconv.Atoi(q.Get("pageSize"))
	return settings.NewListQueryAlias(page, size, q.Get("sort"), vals)
}

func (h *settingsHandler) ListSetting() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := h.uc.ListSetting(r.Context(), listQueryFrom(r))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(res))
	}
}

func (h *settingsHandler) CreateSetting() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.SettingCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(
				httpErrors.NewError(http.StatusBadRequest, err.Error())))
			return
		}
		result, err := h.uc.CreateSetting(r.Context(), req)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(result))
	}
}
```

- GET-delete handler: read `query.Get("<resource>_id")` → parse int → usecase delete → 200 `data: {"affected": n}`
- `Handlers` interface ใน `settings/handler.go`: 1 method ต่อ 1 endpoint ตามกฎตั้งชื่อ — compile จะบังคับครบทุก endpoint จากตาราง Section 6

---

## Section 6: routes — ตาราง endpoint ครบทุกตัว

**Location:** `delivery/http/routes.go` → `MapSettingRoute(router chi.Router, h settings.Handlers, mw *middleware.MiddlewareManager)`

**กติกาอ่านตาราง:** ทุกแถว = 1 route; "ต้นฉบับ controller→service" = function ที่ต้องเปิด port logic/SQL; ชื่อ Go method ตามกฎ PascalCase(path) — แถวที่ต่างจากกฎระบุชัด

### Master (handlers_master.go ↔ pg_repository_master.go)

| Verb+Path | Handler/Usecase | ต้นฉบับ controller → service |
|---|---|---|
| GET /listsetting | ListSetting | list_user_logs → setting_list_paginate |
| GET /settingall | SettingAll | setting_all → setting_all |
| POST /createsetting | CreateSetting | create_setting → create_setting |
| POST /updatesetting | UpdateSetting | update_setting → update_setting |
| GET /deletesetting | DeleteSettingViaGet | delete_setting → delete_setting |
| DELETE /deletesetting | DeleteSetting | delete_setting_del → delete_setting |
| GET /listlocation | ListLocation | list_location → location_list_paginate |
| GET /locationall | LocationAll | location_all → location_all |
| POST /createlocation | CreateLocation | create_location → create_location |
| POST /updatelocation | UpdateLocation | update_location → update_location |
| GET /deletelocation | DeleteLocation | delete_location → delete_location |
| GET /listtype | ListType | list_type → type_list_paginate |
| GET /typeall | TypeAll | list_type_all → type_all |
| POST /createtype | CreateType | create_type → create_type |
| POST /updatetype | UpdateType | update_type → update_type |
| GET /deletetype | DeleteType | delete_type → delete_type |
| GET /listdevicetype | ListDeviceTypePage | listdevicetype → devicetype_list_paginate |
| GET /devicetypeall | DeviceTypeAll | list_devicetype_all → devicetype_all *(register ครั้งเดียว)* |
| GET /devicetypeallcontrol | DeviceTypeAllControl | devicetypeallcontrol → devicetype_all_oi |
| POST /createdevicetype | CreateDeviceType | create_device_type → create_device_type |
| POST /updatedevicetype | UpdateDeviceType | update_device_type → update_device_type |
| GET /deletedevicetype | DeleteDeviceType | delete_device_type → delete_device_type |
| GET /lisgroup | ListGroup | list_group → group_list_paginate |
| GET /lisgroupall | GroupAll | list_group_all → group_all |
| GET /listgrouppage | ListGroupPage | list_group_page → group_list_paginate |
| POST /creategroup | CreateGroup | create_group → create_group |
| POST /updategroup | UpdateGroup | update_group → update_group |
| GET /deletegroup | DeleteGroup | delete_group → delete_group |
| GET /listsensor | ListSensor | list_sensor → sensor_list_paginate *(deviation: ใช้ sensor query)* |
| GET /sensorall | SensorAll | list_sensor_all → sensor_all |
| POST /createsensor | CreateSensor | create_sensor → create_sensor |
| POST /updatesensor | UpdateSensor | update_sensor → update_sensor |
| GET /deletesensor | DeleteSensor | delete_sensor → delete_sensor |

### Utility (handlers_utility.go)

| Verb+Path | Handler/Usecase | ต้นฉบับ |
|---|---|---|
| GET /sensortype | SensorType | sensortype → iothelper.list_sensor_type / _th (port static slice, query `lang=th` เลือกชุด TH) |
| GET /testgemail | TestGmailConnection | testgemail → testGmailConnection (SMTP verify จริง via env) |
| GET /sendemail | SendEmail | sendemail → sendEmail (stub success) |
| GET /mqttdata | MqttData | mqttdata → Cache + getdevicedataALL (Section 4 โค้ดเต็ม) |

### Device (handlers_device.go)

| Verb+Path | Handler/Usecase | ต้นฉบับ controller → service |
|---|---|---|
| GET /deviceall | ListDeviceAll | list_device_all → device_all |
| GET /listdevicepage | ListDevicePage | device_listdevicepage → device_list_paginate |
| GET /listdevicepagess | ListDevicePagess | device_list_paginate(ctrl) → device_list_paginate |
| GET /listdevicepageactive1 | ListDevicePageActive1 | device_list_paginate_active1 → device_list_paginate (active=1) |
| GET /listdevicepageactive | ListDevicePageActive | device_list_paginate_actives → device_list_paginate (active variant) |
| GET /listdevicepageall | ListDevicePageAll | device_list_paginate_all → device_list_paginate_all |
| GET /listdevicepagesensor | ListDevicePageSensor | list_device_page_ensor → device_list_paginate_all |
| GET /listdevicepageallactive | ListDevicePageAllActive | device_list_paginate_all_active → device_list_paginate_all_active |
| GET /listdevicepageallactiveschedule | ListDevicePageAllActiveSchedule | device_list_paginate_all_active_schedule → device_list_paginate_all_active_schedule |
| GET /deviceeditget | DeviceEditGet | deviceeditget → get_device |
| GET /devicedetail | DeviceDetail | devicedetail → get_device (join detail) |
| GET /devicedelete | DeviceDeleteCheck | deviceedelete → delete_device (check view) |
| GET /deletedevice | DeleteDevice | delete_device → delete_device |
| POST /createdevice | CreateDevice | create_device → create_device |
| POST /updatedevice | UpdateDevice | update_device → update_device |
| POST /updatestatusdeviceid | UpdateStatusDeviceId | update_status_deviceid → update_status_deviceid |
| POST /deviceactionuser | DeviceActionUser | Deviceactionuser → create_deviceactionuser |

### Schedule (handlers_schedule.go)

| Verb+Path | Handler/Usecase | ต้นฉบับ controller → service |
|---|---|---|
| GET /listscheduledevice | ListScheduleDevice | findscheduledevice → findscheduledevice (join sd_iot_schedule + sd_iot_device, filters วัน monday..sunday/start/event) |
| GET /findscheduledevicechk | FindScheduleDeviceChk | findscheduledevicechk → findscheduledevicechk |
| GET /schedulelist | ScheduleList | schedulelist → schedule list |
| GET /scheduleall | ScheduleAll | list_schedule_all → schedule_all |
| GET /listschedulepage | ListSchedulePage | schedule_list_page → schedule_list_paginate |
| GET /scheduledevicepage | ScheduleDevicePage | scheduledevicepage → schedule device paginate |
| GET /listdevicescheduledata | ListDeviceScheduleData | listdevicescheduledata → get_data_schedule_device |
| GET /createscheduledevice | CreateScheduleDeviceViaGet | create_schedule_device → createscheduledevice |
| GET /deletescheduledevice | DeleteScheduleDevices | delete_schedule_devices → delete_schedule_device |
| GET /deletedeviceschedule | DeleteDeviceSchedule | delete_device_schedule_id → removeScdeviceId |
| GET /deletedeviceandschedule | DeleteDeviceAndSchedule | delete_device_and_schedule → delete device + schedule |
| POST /createschedule | CreateSchedule | create_schedule → create_schedule |
| POST /createscheduledevice | CreateScheduleDevice | create_scheduleDevice → createscheduledevice |
| POST /updateschedule | UpdateSchedule | update_schedule → update_schedule |
| POST /updateschedulestatus | UpdateScheduleStatus | update_schedule_status → update_schedule_status |
| POST /updatescheduledaystatus | UpdateScheduleDayStatus | update_schedule_day_status → update_schedule_status_day |
| GET /deleteschedule | DeleteSchedule | delete_schedule → delete_schedule |

### Integration (handlers_integration.go)

| Verb+Path | Handler/Usecase | ต้นฉบับ controller → service |
|---|---|---|
| GET /lismqtt | ListMqtt | list_mqtt → mqtt_list_paginate |
| GET /listmqtt | ListMqttAlt | lists_mqtt → mqtt_list_paginate |
| GET /lismqttall | MqttAll | list_mqtt_all → mqtt_all |
| GET /getmqttdetail | GetMqttDetail | get_mqtt_detail → get_mqtt_detail |
| GET /listmqttpaginate | ListMqttPaginate | mqtt_list_paginate → mqtt_list_paginate |
| GET /listmqttpaginateactive | ListMqttPaginateActive | mqtt_list_active_paginate → mqtt_list_active_paginate |
| GET /listmqttdevicepaginate | ListMqttDevicePaginate | mqtt_list_device_paginate → mqtt_list_device_paginate |
| GET /mqttdelete | DeleteMqtt | mqttdelete → delete_mqtt |
| POST /createmqtt | CreateMqtt | create_mqtt → create_mqtt |
| POST /updatemqtt | UpdateMqtt | update_mqtt → update_mqtt |
| POST /updatemqttstatus | UpdateMqttStatus | update_mqtt_status → update_mqtt_status |
| POST /mqtttsort | UpdateMqtttSort | update_mqttt_sort → update_mqttt_sort |
| GET /listmqtthost | ListMqttHost | mqtthost_list_paginate → mqtthost_list_paginate |
| GET /mqtthostall | MqttHostAll | list_mqtthost_all → mqtthost_all |
| POST /createmqtthost | CreateMqttHost | create_mqtthost → create_mqtthost |
| POST /updatemqtthost | UpdateMqttHost | update_mqtthost → update_mqtthost |
| POST /updatemqtthoststatus | UpdateMqttHostStatus | update_mqtthoststatus → update_mqtthoststatus |
| GET /deletemqtthost | DeleteMqttHost | delete_emqtthost → delete_emqtthost |
| GET /apiall | ApiAll | list_api_all → api_all |
| GET /listapipage | ListApiPage | api_list_paginate → api_list_paginate |
| POST /createapi | CreateApi | create_api → create_api |
| POST /updateapi | UpdateApi | update_api → update_api |
| GET /deleteapi | DeleteApi | delete_api → delete_api |
| GET /listemail | ListEmail | email_list_paginate → email_list_paginate |
| GET /emailall | EmailAll | list_email_all → email_all |
| POST /createemail | CreateEmail | create_email → create_email |
| POST /updateemail | UpdateEmail | update_email → update_email |
| POST /updateemailstatus | UpdateEmailStatus | updateemailstatus → updateemailstatus |
| GET /deleteemail | DeleteEmail | delete_email → delete_email |
| GET /hostall | HostAll | list_host_all → host_all |
| GET /listhostpage | ListHostPage | host_list_paginate → host_list_paginate |
| POST /createhost | CreateHost | create_host → create_host |
| POST /updatehost | UpdateHost | update_host → update_host |
| GET /deletehost | DeleteHost | delete_host → delete_host |
| GET /influxdball | InfluxdbAll | list_influxdb_all → influxdb_all |
| GET /listinfluxdbpage | ListInfluxdbPage | influxdb_list_paginate → influxdb_list_paginate |
| POST /createinfluxdb | CreateInfluxdb | create_influxdb → create_influxdb |
| POST /updateinfluxdb | UpdateInfluxdb | update_influxdb → update_influxdb |
| POST /updateinfluxdbstatus | UpdateInfluxdbStatus | updateinfluxdbstatus → updateinfluxdbstatus |
| GET /deleteinfluxdb | DeleteInfluxdb | delete_influxdb → delete_influxdb |
| GET /lineall | LineAll | list_line_all → line_all |
| GET /listlinepage | ListLinePage | line_list_paginate → line_list_paginate |
| POST /createline | CreateLine | create_line → create_line |
| POST /updateline | UpdateLine | update_line → update_line |
| POST /updatelinestatus | UpdateLineStatus | updatelinestatus → updatelinestatus |
| GET /deleteline | DeleteLine | delete_line → delete_line |
| GET /noderedall | NoderedAll | list_nodered_all → nodered_all |
| GET /listnoderedpaginate | ListNoderedPaginate | nodered_list_paginate → nodered_list_paginate |
| POST /createnodered | CreateNodered | create_nodered → create_nodered |
| POST /updatenodered | UpdateNodered | update_nodered → update_nodered |
| POST /updatenoderedstatus | UpdateNoderedStatus | updatenoderedstatus → updatenoderedstatus |
| GET /deletenodered | DeleteNodered | delete_nodered → delete_nodered |
| GET /smsall | SmsAll | list_sms_all → sms_all |
| GET /listsmspage | ListSmsPage | sms_list_paginate → sms_list_paginate |
| POST /createsms | CreateSms | create_sms → create_sms |
| POST /updatesms | UpdateSms | update_sms → update_sms |
| POST /updatesmsstatus | UpdateSmsStatus | updatesmsstatus → updatesmsstatus |
| GET /deletesms | DeleteSms | delete_sms → delete_sms |
| GET /tokenall | TokenAll | list_token_all → token_all |
| GET /tokensmspage | ListTokenPage | token_list_paginate → token_list_paginate |
| POST /createtoken | CreateToken | create_token → create_token |
| POST /updatetoken | UpdateToken | update_token → update_token |
| GET /deletetoken | DeleteToken | delete_token → delete_token |
| POST /createtelegram | CreateTelegram | create_telegram → create_telegram |
| POST /updatetelegram | UpdateTelegram | update_telegram → update_telegram |
| GET /deletetelegram | DeleteTelegram | delete_telegram → delete_telegram |

### Dashboard config (handlers_integration.go)

| Verb+Path | Handler/Usecase | ต้นฉบับ controller → service |
|---|---|---|
| POST /dashboardconfig | CreateDashboardConfig | dashboardconfig_1 → createDashboardConfig |
| GET /dashboardconfig_1 | DashboardConfigByLocation | dashboardconfig_1 → find by location_id (404 ถ้าไม่ส่ง location_id) |
| GET /dashboardconfig | ListDashboardConfig | dashboardconfig → dashboardConfig list |
| GET /dashboardconfig/search | FindOrCreateDashboardConfig | findByCriteria → findOrCreateConfig(name, location_id) |
| GET /dashboardconfig/:id | GetDashboardConfig | findByCriteria → by id |
| PATCH /dashboardconfig/:id | UpdateDashboardConfig | updatedashboardconfig → updateDashboardConfig |
| DELETE /dashboardconfig/:id | RemoveDashboardConfig | remove → removeDashboardConfig |

### Alarm + monitor + logs (handlers_alarm.go)

| Verb+Path | Handler/Usecase | ต้นฉบับ controller → service |
|---|---|---|
| GET /listalarmdevicepage | ListAlarmDevicePage | list_alarm_device_page → alarm device paginate |
| GET /listalarmdeviceactivepage | ListAlarmDeviceActivePage | listalarmdeviceactivepage → alarm active paginate |
| GET /listalarmeventdevicepage | ListAlarmEventDevicePage | list_alarm_event_device_page → alarm event paginate |
| GET /listalarmeventdevicecontrolpage | ListAlarmEventDeviceControlPage | listalarmeventdevicecontrolpage → alarm event control paginate |
| GET /activealarmdevicepage | ActiveAlarmDevicePage | device_list_paginate_alarm_active → device_list_paginate_alarm_active |
| GET /activealarmeventdeviceeventpage | ActiveAlarmEventDeviceEventPage | device_event_list_paginate_alarm_active → device_event_list_paginate_alarm_active |
| GET /alarmdevice | AlarmDevice | alarm_device_paginate → alarm_device_paginate |
| GET /alarmdevicestatus | AlarmDeviceStatus | alarmdevicestatus → alarmdevicestatus |
| GET /deviceactivemqttalarm | DeviceActiveMqttAlarm | deviceactivemqttAlarm → live mqtt + alarm join |
| GET /devicealarm | DeviceAlarm | devicealarm → computed in controller (port ตรง) |
| GET /createalarmdevice | CreateAlarmDeviceViaGet | create_alarm_device → create_alarm_device |
| GET /deletealarmdevice | DeleteAlarmDevices | delete_alarm__devices → delete_alarm__devices |
| GET /createalarmeventdevice | CreateAlarmEventDeviceViaGet | create_alarm_event_device → create_alarm_event_device |
| GET /deletealarmeventdevice | DeleteAlarmEventDevices | delete_alarm_event_devices → delete_alarm_event_devices |
| GET /deletearmdevice | DeleteArmDevice | delete_armdevice → delete_armdevice |
| GET /deletearmdevicev2 | DeleteArmDeviceV2 | delete_armdevice_v2 → delete_armdevice_v2 |
| POST /createalarmDevice | CreateAlarmDevice | create_alarmDevice → create_alarm_device |
| POST /updatealarmdevice | UpdateAlarmDevice | update_alarm_device → update_alarm_device |
| POST /updatealarmstatus | UpdateAlarmStatus | updatealarmstatus → updatealarmstatus |
| POST /createdevicealarmaction | CreateDeviceAlarmAction | create_devicealarmaction → create_devicealarmaction |
| GET /listdevicealarm | ListDeviceAlarm | device_list_alarm → device_list_alarm |
| GET /listdevicealarmairV1 | ListDeviceAlarmAirV1 | device_list_alarm_air → device_list_alarm_air |
| GET /listdevicealarmair | ListDeviceAlarmAir | listdevicealarmair → air alarm join (aircontrol tables) |
| GET /listdevicealarmall | ListDeviceAlarmAll | deviceactivemqtttalarm → deviceactivemqtttalarm |
| GET /_listdevicealarmfan | UnderscoreListDeviceAlarmFan | _device_list_alarm_fan → _device_list_alarm_fan |
| GET /listdevicealarmfan | ListDeviceAlarmFan | device_list_alarm_fan → device_list_alarm_fan |
| GET /listdevicealarmlimit | ListDeviceAlarmLimit | device_list_alarm_limit → device_list_alarm_limit |
| GET /_devicemonitor | UnderscoreDeviceMonitor | _devicemonitor → monitor query |
| GET /devicemonitor | DeviceMonitor | devicemonitor → monitor query |
| GET /devicemonitors | DeviceMonitors | devicemonitors → monitor query |
| GET /scheduleproces | ScheduleProces | scheduleproces → schedule process run state |
| GET /scheduleprocesslog | ScheduleProcessLog | schedule_process_log_page → scheduleprocesslog paginate |
| GET /scheduleprocesslogpaginate | ScheduleProcessLogPaginate | scheduleprocesslogpaginate → scheduleprocesslog paginate |
| GET /mqtterrorlogpaginate | MqttErrorLogPaginate | mqtterrorlogpaginate → mqtterrorlog paginate |
| GET /alarmlogpaginate | AlarmLogPaginate | alarmlogpaginate → alarmprocesslog paginate |
| GET /alarmlogpaginateemail | AlarmLogPaginateEmail | alarmlogpaginateemail → alarmprocesslogemail paginate |
| GET /alarmlogpaginateline | AlarmLogPaginateLine | alarmlogpaginateline → alarmprocesslogline paginate |
| GET /alarmlogpaginatesms | AlarmLogPaginateSms | alarmlogpaginatesms → alarmprocesslogsms paginate |
| GET /alarmlogpaginatetelegram | AlarmLogPaginateTelegram | alarmlogpaginatetelegram → alarmprocesslogtelegram paginate |
| GET /alarmlogpaginatecontrols | AlarmLogPaginateControls | alarmlogpaginatecontrols → controls paginate |
| GET /alarmlogpaginatecontrol | AlarmLogPaginateControl | alarmlogpaginatecontrol → control paginate |

**Route registration pattern:**

```go
func MapSettingRoute(router chi.Router, h settings.Handlers, mw *middleware.MiddlewareManager) {
	router.Route("/settings", func(r chi.Router) {
		r.Use(mw.Verifier(true))
		r.Use(mw.Authenticator())
		r.Use(mw.CurrentUser())
		r.Use(mw.ActiveUser())

		// master
		r.Get("/listsetting", h.ListSetting())
		r.Get("/settingall", h.SettingAll())
		r.Post("/createsetting", h.CreateSetting())
		r.Post("/updatesetting", h.UpdateSetting())
		r.Get("/deletesetting", h.DeleteSettingViaGet())
		r.Delete("/deletesetting", h.DeleteSetting())
		// ... ตามตารางทุกแถว ...

		// static ก่อน param (chi resolve ถูกต้องอยู่แล้ว)
		r.Get("/dashboardconfig/search", h.FindOrCreateDashboardConfig())
		r.Get("/dashboardconfig/{id}", h.GetDashboardConfig())
	})
}
```

> หมายเหตุ: chi ใช้ `{id}` (ไม่ใช่ `:id`)

---

## Section 7: server wiring

**Location:** `internal/server/handlers.go`

**From → To:** ไม่มี settings wiring → เพิ่ม import + init + mount

**Change detail:**

1. imports เพิ่ม:
```go
settingsHttp "icmongolang/internal/modules/settings/delivery/http"
settingsPkg "icmongolang/internal/modules/settings"
settingsRepoPkg "icmongolang/internal/modules/settings/repository"
settingsUsecase "icmongolang/internal/modules/settings/usecase"
```

2. ใน `New(...)` หลังบรรทัด `redisCache := redisDb.NewCache(redisClient)` (ปัจจุบัน ~line 231):
```go
// --- Settings module ---
settingsRepoVar := settingsRepoPkg.CreateSettingsPgRepository(db)
settingsUC := settingsUsecase.CreateSettingsUseCaseI(settingsRepoVar, mqttClient, redisClient, cfg, logger)
settingsHandler := settingsHttp.CreateSettingsHandler(settingsUC, cfg, logger)
settingsHttp.MapSettingRoute(apiRouter, settingsHandler, mw)
logger.Info("✅ Settings routes registered")
```

---

## Section 8: ลำดับการ implement (phases — build ต้อง green ทุก phase)

1. **Skeleton:** types.go (ListSpec/Filter/ListQuery/alias) + repository core + interfaces + presenter + constructor + wiring โดยมี handler/usecase/repo ว่างเป็น stub ที่ return `ErrNotImplemented` → `go build ./...` ผ่าน
2. **Master domain:** setting/location/type/devicetype/group/sensor + tests
3. **Device**
4. **Schedule**
5. **Integration** (mqtt → mqtthost → api → email → host → influxdb → line → nodered → sms → token → telegram → dashboardconfig)
6. **Alarm + logs**
7. **Utility** (sensortype/testgemail/sendemail/mqttdata)

ทุก phase: `go build ./... && go vet ./...` + `go test ./internal/modules/settings/...`

---

## Unit Tests

**Decision: เขียน**

| File | Cases |
|------|-------|
| `repository/pg_repository_test.go` | (integration, ใช้ test Postgres จาก docker-compose dev) ListPaginate: count ถูกต้อง, LIKE escape `%`/`_`, sort whitelist invalid → ErrInvalidSort, offset/limit math, joined columns มีในผล |
| `types_test.go` | ParseSort: valid/missing-dir/bad-field/bad-direction/empty→fallback; EscapeLike |
| `usecase/usecase_test.go` | CreateSetting default status=1; DeleteSetting not found → 404 error; UpdateSetting partial fields map ถูก; SendEmailStub defaults; mock repo ด้วย fake struct |
| `delivery/http/handlers_test.go` | httptest+chi: token ไม่ผ่าน → 401 (mock mw หรือ route นอก group สำหรับ unit), body decode fail → 400 envelope `is_success=false`, ListSetting happy path envelope `is_success=true` + meta fields, DELETE deletesetting happy path `affected` |

## Verification

- `go build ./... && go vet ./...` ผ่านทุก phase
- `go test ./internal/modules/settings/...` ผ่าน
- Run server → smoke representative endpoints ต่อ domain ด้วย token จริง: `GET /api/settings/listsetting?page=1&pageSize=10`, `POST /api/settings/createsetting`, `GET /api/settings/listdevicepageall`, `GET /api/settings/listscheduledevice?schedule_id=&device_id=`, `GET /api/settings/mqttdata?mqttdata=<topic>` — เทียบ semantics กับ NestJS (envelope ต่างตาม design)
