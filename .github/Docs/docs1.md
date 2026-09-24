We'll implement the missing services in Go following the existing Clean Architecture structure. We'll extend `internal/modules/iot/usecase` with device management and data processing, add new repository interfaces for new entities, and provide HTTP handlers and WebSocket support.
NodeJS type script  convert to golang   GORM entities


### โฟลเดอร์หลัก (Modules)
โปรเจกต์ใช้ **Clean Architecture** 3-layer + Delivery:
| Layer | ตำแหน่ง | หน้าที่ |
|-------|---------|--------|
| **Model** | `internal/models/` | Entity (GORM) – `User`, `Session`, `VerificationToken` |
| **Repository** | `internal/repository/` | อ่าน/เขียน DB และ Redis ผ่าน interface |
| **Usecase** | `internal/usecase/` | Business logic: hash, JWT, email queue, validation |
| **Delivery** | `internal/delivery/rest/` | HTTP handlers, middleware, DTO, router |
| **Worker** | `internal/delivery/worker/` | Background job สำหรับส่งอีเมล |
 
```
  api/
    ├── cmd/                     # Cobra CLI (serve, migrate, initdata, worker)
    │   ├──apir/
    │   │  └── main.go
    │   ├── initdata.go             
    │   ├── root.go
    │   ├── serve.go
    │   └── worker.go        
    ├── config/                  # Viper config (YAML + env)
    │   ├── config.default.yml
    │   ├── config.dev.yml
    │   └── config.go  
    ├── internal/                # โค้ดส่วนตัว (ไม่ถูก import จากภายนอก)
    │   ├── iot/                 # shared packages (jwt, redis, email, logger, hash, utils)
    │   │   ├── delivery/        #  HTTP handlers, middleware, dto, router
    │   │   │     └── http/
    │   │   │          ├── handler.go
    │   │   │          └── routes.go
    │   │   ├── worker/        
    │   │   │     └──worker.go
    │   │   ├── helper/          #helper
    │   │   │     └──alarm.go
    │   │   ├── models/          # GORM entities
    │   │   │     ├── alarm.go
    │   │   │     ├── common.go
    │   │   │     └── device_type.go
    │   │   ├── presenter/          
    │   │   │     └── presenter.go
    │   │   ├── repository/         
    │   │   │     ├── alarm_log_repo.go
    │   │   │     ├── device_repo.go
    │   │   │     └── schedule_repo.go
    │   │   ├── usecase/          # GORM entities
    │   │   │     └── usecase.go
    │   │   └── iot.go
    │   ├── repository/          # interfaces + impl (postgres, redis)
    │   │       ├── pg.go
    │   │       └── redis.go
    │   ├── server/           
    │   │       ├── handlers.go
    │   │       └── server.go
    │   ├── usecase/             # business logic
    │   │       └── usecase.go
    │   ├── pg_repository.go   
    │   ├── redis_repository.go   
    │   ├── usecase.go   
    │   ├── pkg/                 #  shared packages (jwt, redis, email, logger, hash, utils)
    │   │   ├── db/              #  DB postgres,redis
    │   │   │   ├── postgres/ 
    │   │   │   │     └── db_conn.go
    │   │   │   └── redis/ 
    │   │   │          └── redis_conn.go
    │   │   ├── helpers/
    │   │   │   ├── iot.go          # Alarm logic (สมบูรณ์)
    │   │   │   └── format.go       # ฟังก์ชันช่วยเหลือ (time, string, random)
    │   │   ├── mqtt/
    │   │   │   └── client.go       # MQTT client พร้อม GetDataFromTopic
    │   │   ├── influxdb/
    │   │   │   └── client.go       # InfluxDB client
    │   │   └── httpErrors/
    │   │       └── httpErrors.go   # httpErrors
    ├── migrations/              # raw SQL (optional)
    ├── docker-compose.dev.yml   # Postgres + Redis + MailHog
    ├── Dockerfile.dev / .air.toml
    └── go.mod
```
---

## 1. New Repository Interfaces & Implementations

We'll create repositories for `DeviceStatus`, `DeviceConfig`, `IotData`, `ActivityLog`, `CommandLog`, and `DeviceAlert`. Place them in `internal/modules/iot/repository/`.

### `device_status_repo.go`

```go
package repository

import (
	"context"
	"icmongolang/internal/modules/iot/models"
	"gorm.io/gorm"
)

type DeviceStatusRepository interface {
	GetByDeviceID(ctx context.Context, deviceID string) (*models.DeviceStatus, error)
	Upsert(ctx context.Context, status *models.DeviceStatus) error
	UpdateLastSeen(ctx context.Context, deviceID string) error
}

type deviceStatusRepo struct {
	db *gorm.DB
}

func NewDeviceStatusRepository(db *gorm.DB) DeviceStatusRepository {
	return &deviceStatusRepo{db: db}
}

func (r *deviceStatusRepo) GetByDeviceID(ctx context.Context, deviceID string) (*models.DeviceStatus, error) {
	var status models.DeviceStatus
	err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).First(&status).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &status, err
}

func (r *deviceStatusRepo) Upsert(ctx context.Context, status *models.DeviceStatus) error {
	return r.db.WithContext(ctx).Save(status).Error
}

func (r *deviceStatusRepo) UpdateLastSeen(ctx context.Context, deviceID string) error {
	return r.db.WithContext(ctx).Model(&models.DeviceStatus{}).
		Where("device_id = ?", deviceID).
		Update("last_seen", gorm.Expr("NOW()")).Error
}
```

### `device_config_repo.go`

```go
package repository

import (
	"context"
	"icmongolang/internal/modules/iot/models"
	"gorm.io/gorm"
)

type DeviceConfigRepository interface {
	GetByDeviceID(ctx context.Context, deviceID string) (*models.DeviceConfig, error)
	Upsert(ctx context.Context, config *models.DeviceConfig) error
}

type deviceConfigRepo struct {
	db *gorm.DB
}

func NewDeviceConfigRepository(db *gorm.DB) DeviceConfigRepository {
	return &deviceConfigRepo{db: db}
}

func (r *deviceConfigRepo) GetByDeviceID(ctx context.Context, deviceID string) (*models.DeviceConfig, error) {
	var config models.DeviceConfig
	err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &config, err
}

func (r *deviceConfigRepo) Upsert(ctx context.Context, config *models.DeviceConfig) error {
	return r.db.WithContext(ctx).Save(config).Error
}
```

### `iot_data_repo.go` (if not already present)

```go
package repository

import (
	"context"
	"icmongolang/internal/modules/iot/models"
	"gorm.io/gorm"
)

type IotDataRepository interface {
	Create(ctx context.Context, data *models.IotData) error
	GetLatest(ctx context.Context, deviceID string) (*models.IotData, error)
	GetByDateRange(ctx context.Context, deviceID string, start, end interface{}) ([]models.IotData, error)
	GetByDeviceID(ctx context.Context, deviceID string, limit int) ([]models.IotData, error)
	DeleteOlderThan(ctx context.Context, cutoff interface{}) (int64, error)
}

type iotDataRepo struct {
	db *gorm.DB
}

func NewIotDataRepository(db *gorm.DB) IotDataRepository {
	return &iotDataRepo{db: db}
}

func (r *iotDataRepo) Create(ctx context.Context, data *models.IotData) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *iotDataRepo) GetLatest(ctx context.Context, deviceID string) (*models.IotData, error) {
	var data models.IotData
	err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Order("timestamp DESC").
		First(&data).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &data, err
}

func (r *iotDataRepo) GetByDateRange(ctx context.Context, deviceID string, start, end interface{}) ([]models.IotData, error) {
	var list []models.IotData
	err := r.db.WithContext(ctx).
		Where("device_id = ? AND timestamp BETWEEN ? AND ?", deviceID, start, end).
		Order("timestamp ASC").
		Find(&list).Error
	return list, err
}

func (r *iotDataRepo) GetByDeviceID(ctx context.Context, deviceID string, limit int) ([]models.IotData, error) {
	var list []models.IotData
	err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Order("timestamp DESC").
		Limit(limit).
		Find(&list).Error
	return list, err
}

func (r *iotDataRepo) DeleteOlderThan(ctx context.Context, cutoff interface{}) (int64, error) {
	res := r.db.WithContext(ctx).
		Where("timestamp < ?", cutoff).
		Delete(&models.IotData{})
	return res.RowsAffected, res.Error
}
```

### `activity_log_repo.go`

```go
package repository

import (
	"context"
	"icmongolang/internal/modules/iot/models"
	"gorm.io/gorm"
)

type ActivityLogRepository interface {
	Create(ctx context.Context, log *models.ActivityLog) error
}

type activityLogRepo struct {
	db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) ActivityLogRepository {
	return &activityLogRepo{db: db}
}

func (r *activityLogRepo) Create(ctx context.Context, log *models.ActivityLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
```

Similarly, you can create `command_log_repo.go` and `device_alert_repo.go` as needed.

---

## 2. Extended Usecase (Device Management & Data Processing)

We'll extend the existing `MQTT3UseCase` interface and implementation (already in `internal/modules/iot/usecase/usecase.go`). Add new methods for device status, config, data processing, etc.

### Update `usecase.go` – add to interface

```go
type MQTT3UseCase interface {
	// ... existing methods ...

	// Device management
	GetDeviceStatus(ctx context.Context, deviceID string) (*presenter.DeviceStatusResponse, error)
	UpdateDeviceStatus(ctx context.Context, deviceID string, data map[string]interface{}) error
	GetDeviceConfig(ctx context.Context, deviceID string) (*models.DeviceConfig, error)
	UpdateDeviceConfig(ctx context.Context, deviceID string, config map[string]interface{}) error

	// IoT data processing
	ProcessMqttData(ctx context.Context, deviceID string, rawData string) (*models.IotData, error)
	GetLatestData(ctx context.Context, deviceID string, limit int) ([]models.IotData, error)
	GetDataByDateRange(ctx context.Context, deviceID string, start, end time.Time) ([]models.IotData, error)
	CleanupOldData(ctx context.Context, days int) (int64, error)

	// Additional methods from IotsocketioService
	ListIotData(ctx context.Context, opts *presenter.IotDataListOptions) (*presenter.PaginatedIotData, error)
	GetDeviceStats(ctx context.Context, deviceID string) (*presenter.DeviceStats, error)
	ExportData(ctx context.Context, req *presenter.ExportRequest) ([]byte, string, error)
	// ... etc.
}
```

### Implement these methods in `mqtt3UseCase` struct

Add the new repositories to the struct:

```go
type mqtt3UseCase struct {
	// ... existing fields ...
	deviceStatusRepo repository.DeviceStatusRepository
	deviceConfigRepo repository.DeviceConfigRepository
	iotDataRepo      repository.IotDataRepository
	activityLogRepo  repository.ActivityLogRepository
	commandLogRepo   repository.CommandLogRepository
	deviceAlertRepo  repository.DeviceAlertRepository
}
```

Update the constructor `NewMQTT3UseCase` to accept these new dependencies.

Now implement the methods:

#### GetDeviceStatus

```go
func (u *mqtt3UseCase) GetDeviceStatus(ctx context.Context, deviceID string) (*presenter.DeviceStatusResponse, error) {
	status, err := u.deviceStatusRepo.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	if status == nil {
		// create default status
		status = &models.DeviceStatus{
			DeviceID:  deviceID,
			IsOnline:  false,
			IsActive:  true,
			LastSeen:  time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_ = u.deviceStatusRepo.Upsert(ctx, status)
	}
	// calculate online status
	fifteenMinAgo := time.Now().Add(-15 * time.Minute)
	isOnline := status.LastSeen.After(fifteenMinAgo)
	return &presenter.DeviceStatusResponse{
		DeviceID:      status.DeviceID,
		IsOnline:      isOnline,
		IsActive:      status.IsActive,
		LastSeen:      status.LastSeen,
		BatteryLevel:  status.BatteryLevel,
		SignalStrength: status.SignalStrength,
		FirmwareVersion: status.FirmwareVersion,
		Location:      status.Location,
		LastData:      status.LastData,
		Uptime:        calculateUptime(status.FirstSeen),
	}, nil
}

func calculateUptime(firstSeen *time.Time) string {
	if firstSeen == nil {
		return "0s"
	}
	dur := time.Since(*firstSeen)
	return dur.String() // or format nicely
}
```

#### UpdateDeviceStatus

```go
func (u *mqtt3UseCase) UpdateDeviceStatus(ctx context.Context, deviceID string, data map[string]interface{}) error {
	status, err := u.deviceStatusRepo.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return err
	}
	if status == nil {
		status = &models.DeviceStatus{DeviceID: deviceID}
	}
	status.LastSeen = time.Now()
	status.LastData = data
	status.IsOnline = true
	// update fields from data
	if battery, ok := data["battery"].(float64); ok {
		status.BatteryLevel = intPtr(int(battery))
	}
	if signal, ok := data["signal"].(float64); ok {
		status.SignalStrength = intPtr(int(signal))
	}
	if firmware, ok := data["firmware"].(string); ok {
		status.FirmwareVersion = &firmware
	}
	if loc, ok := data["location"].(map[string]interface{}); ok {
		status.Location = loc
	}
	status.UpdatedAt = time.Now()
	return u.deviceStatusRepo.Upsert(ctx, status)
}
```

#### GetDeviceConfig

```go
func (u *mqtt3UseCase) GetDeviceConfig(ctx context.Context, deviceID string) (*models.DeviceConfig, error) {
	cfg, err := u.deviceConfigRepo.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = &models.DeviceConfig{
			DeviceID: deviceID,
			Config:   getDefaultConfig(),
			Status:   "active",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_ = u.deviceConfigRepo.Upsert(ctx, cfg)
	}
	return cfg, nil
}
```

#### ProcessMqttData (core)

```go
func (u *mqtt3UseCase) ProcessMqttData(ctx context.Context, deviceID string, rawData string) (*models.IotData, error) {
	// Parse raw data (assumes CSV format)
	parts := strings.Split(rawData, ",")
	dataMap := make(map[string]interface{})
	// Example: map to fields (customize as needed)
	if len(parts) >= 1 {
		dataMap["temperature"] = parseFloat(parts[0])
	}
	// ... more parsing

	iotData := &models.IotData{
		DeviceID:  deviceID,
		Data:      dataMap,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}
	// save
	if err := u.iotDataRepo.Create(ctx, iotData); err != nil {
		return nil, err
	}
	// update device status
	_ = u.UpdateDeviceStatus(ctx, deviceID, dataMap)
	// log activity
	_ = u.logActivity(ctx, "DATA_RECEIVED", "Received data from "+deviceID, deviceID, nil)
	return iotData, nil
}

func (u *mqtt3UseCase) logActivity(ctx context.Context, typ, details, deviceID string, data interface{}) error {
	log := &models.ActivityLog{
		Type:      typ,
		DeviceID:  &deviceID,
		Details:   details,
		Severity:  "info",
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}
	if data != nil {
		jsonData, _ := json.Marshal(data)
		log.Data = jsonData
	}
	return u.activityLogRepo.Create(ctx, log)
}
```

#### Additional data listing methods (from IotsocketioService)

```go
func (u *mqtt3UseCase) ListIotData(ctx context.Context, opts *presenter.IotDataListOptions) (*presenter.PaginatedIotData, error) {
	// Build query with filters, pagination, sorting
	// Use u.iotDataRepo with custom query builder if needed
	// We'll implement a simple version using repository method
	// For advanced filters, add a method to iotDataRepo
	// For brevity, here's a basic implementation using existing GetByDeviceID
	var list []models.IotData
	var total int64
	// We'll create a new method in repo for paginated list
	// For now, let's assume we add List(ctx, opts) method in repo
	// Return result
	return &presenter.PaginatedIotData{Data: list, Total: total, Page: opts.Page, Limit: opts.Limit}, nil
}
```

You'd need to implement `List` in `IotDataRepository` with filtering and pagination.

---

## 3. Presenter DTOs (add to `presenter.go`)

```go
type DeviceStatusResponse struct {
	DeviceID        string                 `json:"deviceId"`
	IsOnline        bool                   `json:"isOnline"`
	IsActive        bool                   `json:"isActive"`
	LastSeen        time.Time              `json:"lastSeen"`
	BatteryLevel    *int                   `json:"batteryLevel,omitempty"`
	SignalStrength  *int                   `json:"signalStrength,omitempty"`
	FirmwareVersion *string                `json:"firmwareVersion,omitempty"`
	Location        map[string]interface{} `json:"location,omitempty"`
	LastData        map[string]interface{} `json:"lastData,omitempty"`
	Uptime          string                 `json:"uptime"`
}

type IotDataListOptions struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	DeviceID  string `json:"deviceId"`
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
	DataType  string `json:"dataType"`
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"`
}

type PaginatedIotData struct {
	Data       []models.IotData `json:"data"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"totalPages"`
}

type DeviceStats struct {
	Count        int                    `json:"count"`
	FirstRecord  *time.Time             `json:"firstRecord"`
	LastRecord   *time.Time             `json:"lastRecord"`
	DataPoints   map[string]interface{} `json:"dataPoints"`
}

type ExportRequest struct {
	DeviceID  string    `json:"deviceId"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Format    string    `json:"format"` // csv, json
}
```

---

## 4. HTTP Handlers (Delivery)

Create new handlers in `internal/modules/iot/delivery/http/handler.go` (or a separate file `device_handler.go`).

Add these methods to `MQTT3Handler` struct (or create a separate handler struct). We'll extend `MQTT3Handler`.

### Add new routes in `routes.go`

```go
func MapMQTT3Routes(router chi.Router, h *MQTT3Handler, mw *middleware.MiddlewareManager) {
	router.Route("/iot", func(r chi.Router) {
		// ... existing routes ...

		// Device status & config
		r.Get("/device/status", h.GetDeviceStatus)
		r.Put("/device/status", h.UpdateDeviceStatus) // via JSON body
		r.Get("/device/config", h.GetDeviceConfig)
		r.Put("/device/config", h.UpdateDeviceConfig)

		// IoT data
		r.Get("/device/data", h.ListIotData)
		r.Get("/device/data/latest", h.GetLatestData)
		r.Post("/device/data/process", h.ProcessMqttData) // for testing
		r.Delete("/device/data/cleanup", h.CleanupOldData)

		// WebSocket endpoint
		// handled separately (see below)
	})
}
```

### Implement handlers

```go
func (h *MQTT3Handler) GetDeviceStatus(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("deviceId")
	if deviceID == "" {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusBadRequest, "deviceId is required")))
		return
	}
	resp, err := h.uc.GetDeviceStatus(r.Context(), deviceID)
	if err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusInternalServerError, err.Error())))
		return
	}
	render.Render(w, r, responses.CreateSuccessResponse(resp))
}

func (h *MQTT3Handler) UpdateDeviceStatus(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
		return
	}
	deviceID, ok := req["deviceId"].(string)
	if !ok || deviceID == "" {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusBadRequest, "deviceId is required")))
		return
	}
	if err := h.uc.UpdateDeviceStatus(r.Context(), deviceID, req); err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusInternalServerError, err.Error())))
		return
	}
	render.Render(w, r, responses.CreateSuccessResponse(map[string]string{"status": "updated"}))
}

func (h *MQTT3Handler) GetDeviceConfig(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("deviceId")
	if deviceID == "" {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusBadRequest, "deviceId is required")))
		return
	}
	cfg, err := h.uc.GetDeviceConfig(r.Context(), deviceID)
	if err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusInternalServerError, err.Error())))
		return
	}
	render.Render(w, r, responses.CreateSuccessResponse(cfg))
}

func (h *MQTT3Handler) UpdateDeviceConfig(w http.ResponseWriter, r *http.Request) {
	// Similar to update status
}

func (h *MQTT3Handler) ListIotData(w http.ResponseWriter, r *http.Request) {
	// parse query params
	opts := &presenter.IotDataListOptions{
		Page:   atoi(r.URL.Query().Get("page")),
		Limit:  atoi(r.URL.Query().Get("limit")),
		DeviceID: r.URL.Query().Get("deviceId"),
	}
	if opts.Page <= 0 { opts.Page = 1 }
	if opts.Limit <= 0 { opts.Limit = 50 }
	// parse dates if provided
	// ...
	resp, err := h.uc.ListIotData(r.Context(), opts)
	if err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusInternalServerError, err.Error())))
		return
	}
	render.Render(w, r, responses.CreateSuccessResponse(resp))
}

func (h *MQTT3Handler) ProcessMqttData(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DeviceID string `json:"deviceId"`
		RawData  string `json:"rawData"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
		return
	}
	if req.DeviceID == "" || req.RawData == "" {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusBadRequest, "deviceId and rawData required")))
		return
	}
	data, err := h.uc.ProcessMqttData(r.Context(), req.DeviceID, req.RawData)
	if err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusInternalServerError, err.Error())))
		return
	}
	render.Render(w, r, responses.CreateSuccessResponse(data))
}

func (h *MQTT3Handler) CleanupOldData(w http.ResponseWriter, r *http.Request) {
	days := atoi(r.URL.Query().Get("days"))
	if days <= 0 { days = 30 }
	count, err := h.uc.CleanupOldData(r.Context(), days)
	if err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusInternalServerError, err.Error())))
		return
	}
	render.Render(w, r, responses.CreateSuccessResponse(map[string]int64{"deleted": count}))
}
```

---

## 5. WebSocket Support

We'll implement a WebSocket handler using `gorilla/websocket` and integrate it with the existing `MQTT3Handler` or as a separate package. We'll create a new struct `WebSocketHub` to manage clients and broadcast.

Create `internal/modules/iot/delivery/websocket/hub.go`:

```go
package websocket

import (
	"sync"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
	Key  string // optional identifier
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) BroadcastJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	h.broadcast <- data
	return nil
}
```

Then add WebSocket handler in `routes.go` or separate file:

```go
func (h *MQTT3Handler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Errorf("WebSocket upgrade error: %v", err)
		return
	}
	client := &websocket.Client{Conn: conn, Send: make(chan []byte)}
	h.hub.register <- client

	// read pump
	go func() {
		defer func() {
			h.hub.unregister <- client
			conn.Close()
		}()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			// handle incoming messages (e.g., heartbeat, set key)
			// parse JSON and act
			var data map[string]interface{}
			if err := json.Unmarshal(msg, &data); err == nil {
				if event, ok := data["event"].(string); ok {
					switch event {
					case "heartbeat":
						// respond
					case "set_key":
						if key, ok := data["key"].(string); ok {
							client.Key = key
						}
					}
				}
			}
		}
	}()

	// write pump
	go func() {
		defer conn.Close()
		for {
			select {
			case msg, ok := <-client.Send:
				if !ok {
					conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				conn.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}()
}
```

We need to add `hub` field to `MQTT3Handler` and initialize in constructor.

Then we can broadcast from usecase or handler when MQTT data arrives.

---

## 6. Dependency Injection (main.go / server setup)

You'll need to wire up the new repositories and usecase in your `main.go` or server initialization.

Example:

```go
// in server.go
db := ... // *gorm.DB
deviceStatusRepo := repository.NewDeviceStatusRepository(db)
deviceConfigRepo := repository.NewDeviceConfigRepository(db)
iotDataRepo := repository.NewIotDataRepository(db)
activityLogRepo := repository.NewActivityLogRepository(db)
// ... others

usecase := usecase.NewMQTT3UseCase(
	deviceRepo,
	alarmLogRepo,
	mqttClient,
	redisClient,
	influxClient,
	logger,
	cfg,
	// new dependencies:
	deviceStatusRepo,
	deviceConfigRepo,
	iotDataRepo,
	activityLogRepo,
	// ...
)

hub := websocket.NewHub()
go hub.Run()

handler := http.NewMQTT3Handler(usecase, logger, cache, hub) // update constructor

router := chi.NewRouter()
http.MapMQTT3Routes(router, handler, mw)
```

---

## 7. Summary

We've converted the remaining NestJS services to Go, extending the existing Clean Architecture:

- **Repositories** for new entities: `DeviceStatus`, `DeviceConfig`, `IotData`, `ActivityLog`.
- **Usecase** methods for device management, data processing, and data retrieval.
- **HTTP handlers** for REST endpoints (status, config, data, cleanup).
- **WebSocket hub** with broadcast and client management.

All code follows the project's existing patterns (interfaces, dependency injection, response rendering). You may need to adjust field names, JSON tags, and parsing logic to match your actual data schema and business rules.