mqtt3.controller.ts   + mqtt3.service.ts  
 - หลักการทำงาน (Concept) 
  - ออกแบบ workflow
    - วาดรูป dataflow สร้าง รูปแบบ dataflow เหมือนจริง ลักษณะ flowchart   เพื่ออธิบายกระบวนการ ทำความเข้าใจ
    - พร้อมอธิบาย แบบ ละเอียด 
    - คอมเม้น code ภาษาไทย และ ภาษาอังถถษ อธิบาย การทำงาน แต่ละจุด
    - ยกตัวอย่างการใช้งานจริง หรือ กรณีศึกษา แนวทางแก้ไขปัญหา ที่อาจจะเกิดขึ้น  
   - Check list  module การทำงาน 

โครงสร้าง Code

icmongolang/
├── cmd/
│   ├── api/
│   │   └── main.go
│   ├── initdata.go
│   ├── migrate.go
│   ├── root.go
│   ├── serve.go
│   └── worker.go
├── config/
│   ├── config-local.yml
│   ├── config-prod.yml
│   └── config.go
├── docdev/
├── docs/
├── internal/
│   ├── models/ 
│   │   └── device.go
│   ├── device/   
│   │   ├── delivery/
│   │   │   ├── http/
│   │   │   │   ├── handler.go 
│   │   │   │   └── routes.go
│   │   ├── presenter/ 
│   │   │   └── presenter.go
│   │   ├── distributor/ 
│   │   │   └── distributor.go
│   │   ├── processor/ 
│   │   │   └── processor.go
│   │   ├── repository/ 
│   │   │   ├── pg_repository.go
│   │   │   └── redis_repository.go
│   │   ├── usecase/ 
│   │   │   └── usecase.go
│   │   │  
│   │   ├─ handlers.go
│   │   ├─ pg_repository.go
│   │   ├─ redis_repository.go
│   │   └─ usecase.go
├───── pkg/
│       ├── db/ 
│       │   ├── postgres.go/ 
│       │   │     └── db_conn.go 
│       │   └── redis/ 
│       │         └── redis_conn.go 
│       ├── cryptpass/
│       │   └── password.go
│       ├── hash/
│       │   └── bcrypt.go
│       ├── jwt/
│       │   ├── token.go
│       ├── influxdb/
│       │   └── client.go
│       ├── mqtt/
│       │   └── client.go
│       ├── logger/
│       │   └── zap_logger.go
│       ├── redis/
│       │   ├── cache.go
│       │   ├── client.go
│       │   └── refresh_store.go
│       ├── utils/
│       │   ├── random.go
│       │   └── time.go
│       └── validator/
│           └── custom_validator.go

สร้าง modules  iotservice ภาษา golang ตามโครงสร้าง icmongolang/  ใช้โคร้างนี้ ทำให้ครบ ตาม mqtt3.controller.ts   + mqtt3.service.ts    แปลง จาก typescript ไป go lang   
- วัตุประสงค์
  - แก้ repo ดึง today_status ,  deleted_flag
  - MSRTR2547-2:[handler+service+repo] edit GetEmployeeValidation
  -ข้อควรระวัง  
    - ไม่แก้ โครงสร้างเดิม
   -ข้อห้าม 
     ห้าม แก้ไข funvtion โครงสร้างเดิม ที่มรการใช้งาน ส่วนอื่นๆ

- หลักการทำงาน (Concept)  
  
📁file /Path:
    - Path
	- file
📥 Input:
	 - ตัวอย่างโค้ด 
	   - รูปแบบ

🔁 Flow:- 
	 - Workflow และ Dataflow  
	 
📥 Output:
	- สรุป 

🧪 Unit: unit test & swager 
   - Check list unit
   - Check list review Code PR 
   
- cmd
   - Gen swager 
      swag init -g cmd/api/main.go 
   - Gen mock data
      mockery --all      
   - unit test 
      go test ./...
	  
------------------------

## บทสรุปสำหรับผู้บริหาร

งานนี้เป็นการแปลงโมดูล MQTT3 จาก NestJS (TypeScript) ไปเป็นภาษา Go ตามสถาปัตยกรรมของโปรเจกต์ `icmongolang` โดยมีเป้าหมายเพื่อ:

- รองรับการรับ/ส่งข้อมูลจากอุปกรณ์ IoT ผ่าน MQTT
- จัดการ Alarm และการแจ้งเตือน (Email, Line, SMS, Telegram)
- ดึงข้อมูลแบบ Real‑time พร้อม Cache (Redis)
- Query ข้อมูลย้อนหลังจาก InfluxDB
- จัดการ Schedule (งานตามเวลา) และ Log กระบวนการ

โดย **ไม่แก้ไข** โครงสร้างเดิมของระบบ (เพิ่ม module ใหม่เท่านั้น) และต้องสามารถทำงานร่วมกับส่วนอื่น ๆ ที่มีอยู่แล้ว (PostgreSQL, Redis, InfluxDB, MQTT client เดิม)

---

## 1. หลักการทำงาน (Concept) และ Workflow

### 1.1 สถาปัตยกรรมโดยรวม

```
[Device IoT] --- MQTT ---> [MQTT Broker]
                               |
                               v
                       [MQTT3 Service (Go)]
                               |
                +--------------+--------------+
                |              |              |
                v              v              v
           [Redis]       [PostgreSQL]    [InfluxDB]
           (Cache)       (config, log)   (time‑series)
```

- **MQTT3 Service** ทำหน้าที่เป็นทั้ง subscriber และ publisher
- ใช้ **Request‑Response Pattern** พร้อม `MqttRequestManager` เพื่อรอรับ message เฉพาะ topic ที่ต้องการแบบ timeout
- **Cache** ทุก topic ที่ subscribe ไว้ใน Redis (TTL กำหนดได้) เพื่อลดการเรียก MQTT ซ้ำ
- **Alarm Logic** แยกอยู่ใน `iothelper` เพื่อคำนวณสถานะ (Warning, Critical, Recovery) และสร้างข้อความแจ้งเตือน
- **Schedule** ตรวจสอบวันและเวลา (ตาม config ใน DB) แล้วส่งคำสั่งควบคุมอุปกรณ์อัตโนมัติ

### 1.2 Data Flow ของการดึงข้อมูลจาก MQTT (getDataFromTopic)

```mermaid
graph TD
    A[HTTP Request: /mqtt3/gettopicdata?topic=xxx] --> B[Controller]
    B --> C{Usecase: GetTopicData}
    C --> D[Check Redis Cache]
    D -->|Hit| E[Return Cache]
    D -->|Miss| F[MQTT Request Manager]
    F --> G[Subscribe to topic]
    G --> H[Wait for message (timeout)]
    H -->|Received| I[Unsubscribe]
    I --> J[Parse payload]
    J --> K[Store in Redis]
    K --> L[Return to client]
```

### 1.3 Data Flow ของ Alarm (Device Control & Notification)

```mermaid
graph TD
    A[MQTT Message on DATA topic] --> B[MQTT Client callback]
    B --> C[Alarm Processor]
    C --> D[Read device config from DB]
    D --> E[Calculate alarm status via iothelper]
    E --> F{AlarmStatusSet?}
    F -->|Warning/Critical| G[Check cooldown (Redis/DB Log)]
    G -->|Can send| H[Send Email/Line/SMS/Telegram]
    G -->|Can control| I[Publish CONTROL topic]
    F -->|Recovery| J[Send recovery notification]
    F -->|Normal| K[No action]
```

### 1.4 โครงสร้างโมดูล Go ที่จะสร้าง

```
internal/modules/iot/
├── delivery/
│   └── http/
│       ├── handler.go      // MQTT3Handler (GET/POST endpoints)
│       └── routes.go       // Register routes
├── usecase/
│   └── usecase.go          // MQTT3UseCase interface & implementation
├── repository/
│   ├── device_repo.go      // ดึงข้อมูล device, mqtt config จาก PostgreSQL
│   ├── alarm_log_repo.go   // บันทึก alarm log, schedule log
│   └── schedule_repo.go    // ดึง schedule และ device mapping
├── presenter/
│   └── presenter.go        // Request/Response structs
└── iothelper/
    └── alarm.go            // ฟังก์ชัน AlarmDetailValidate (แปลงจาก iot.helper.ts)
pkg/mqtt/
└── manager.go              // MqttRequestManager (request‑response pattern)
```

---

## 2. การแปลงฟังก์ชันหลักจาก TypeScript → Go

เนื่องจากมีโค้ดจำนวนมาก จะขอยกตัวอย่างเฉพาะฟังก์ชันสำคัญที่แสดงแนวทางการแปลง พร้อมคำอธิบาย

### 2.1 MqttRequestManager (request‑response pattern)

**TypeScript** (`mqtt.helper.ts`):
```typescript
export class MqttRequestManager {
  async getDataFromTopic(topic: string, timeoutMs = 10000): Promise<any> { ... }
}
```

**Go** (`pkg/mqtt/manager.go`):
```go
package mqtt

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "time"
    mqtt "github.com/eclipse/paho.mqtt.golang"
)

type pendingRequest struct {
    resolve   chan<- interface{}
    reject    chan<- error
    timeout   *time.Timer
    topic     string
}

type RequestManager struct {
    client            mqtt.Client
    mu                sync.Mutex
    pending           map[string][]*pendingRequest
    subscriptionCount map[string]int
    logger            logger.Logger
}

func NewRequestManager(client mqtt.Client, log logger.Logger) *RequestManager {
    rm := &RequestManager{
        client:            client,
        pending:           make(map[string][]*pendingRequest),
        subscriptionCount: make(map[string]int),
        logger:            log,
    }
    client.AddRoute("#", rm.handleIncomingMessage)
    return rm
}

func (rm *RequestManager) GetDataFromTopic(ctx context.Context, topic string, timeout time.Duration) (interface{}, error) {
    // Implement similar logic as TypeScript version
    // Use channel + timeout, manage subscription count
}
```

### 2.2 AlarmDetailValidate (iothelper)

**TypeScript** (`iot.helper.ts`):
มีฟังก์ชัน `processAlarmDetail` ที่ใช้เช็ค hardware_id, ค่า sensor, min/max, status_warning/alert และคืน `AlarmDetailResult`

**Go** (`internal/modules/iot/iothelper/alarm.go`):
```go
package iothelper

type AlarmDetailDto struct {
    HardwareID        interface{}
    ValueData         interface{}
    Max               interface{}
    Min               interface{}
    StatusAlert       interface{}
    StatusWarning     interface{}
    RecoveryWarning   interface{}
    RecoveryAlert     interface{}
    DeviceName        string
    ActionName        string
    MqttName          string
    MqttControlOn     string
    MqttControlOff    string
    CountAlarm        interface{}
    Event             interface{}
    Unit              string
}

type AlarmDetailResult struct {
    Status             int
    AlarmStatusSet     int
    Title              string
    Subject            string
    Content            string
    DataAlarm          interface{}
    EventControl       int
    MessageMqttControl string
    // ... other fields
}

func AlarmDetailValidate(dto AlarmDetailDto) AlarmDetailResult {
    // แปลง logic จาก processAlarmDetail
    // ใช้ switch case ตาม hardware_id, เปรียบเทียบค่า, สร้างข้อความ
}
```

### 2.3 MQTT3 Service – การเชื่อมต่อและ subscribe

**TypeScript** (`mqtt3.service.ts`):
- มี `initializeMqttClient`, `subscribeToTopicWithResponse`, `getMqttTopicData` ฯลฯ

**Go** – ใช้ MQTT client เดิมที่มีใน `pkg/mqtt/client.go` แต่เพิ่ม `RequestManager` เข้าไป

```go
// internal/modules/iot/usecase/usecase.go
type MQTT3UseCase interface {
    GetTopicData(ctx context.Context, topic string) ([]byte, error)
    DeviceControl(ctx context.Context, topic, message string) error
    GetDeviceList(ctx context.Context, req *presenter.DeviceListRequest) (*presenter.DeviceListResponse, error)
    GetAlarmDeviceStatus(ctx context.Context, req *presenter.AlarmStatusRequest) (*presenter.AlarmStatusResponse, error)
    // ... other methods
}

type mqtt3UseCase struct {
    mqttClient     mqtt.Client
    reqManager     *mqtt.RequestManager
    redisClient    *redis.Client
    db             *gorm.DB
    influxClient   *influxdb.InfluxClient
    logger         logger.Logger
}

func (u *mqtt3UseCase) GetTopicData(ctx context.Context, topic string) ([]byte, error) {
    // 1. ลองจาก Redis cache
    cacheKey := "mqtt_topic:" + topic
    if cached, err := u.redisClient.Get(ctx, cacheKey).Bytes(); err == nil {
        return cached, nil
    }
    // 2. ใช้ RequestManager ขอข้อมูล
    data, err := u.reqManager.GetDataFromTopic(ctx, topic, 10*time.Second)
    if err != nil {
        return nil, err
    }
    var payload []byte
    switch v := data.(type) {
    case string:
        payload = []byte(v)
    case []byte:
        payload = v
    default:
        payload, _ = json.Marshal(v)
    }
    // 3. เก็บ cache (TTL 30 วินาที)
    u.redisClient.Set(ctx, cacheKey, payload, 30*time.Second)
    return payload, nil
}
```

### 2.4 Controller – HTTP endpoints

**TypeScript** (`mqtt3.controller.ts`):
มี endpoints:
- `GET /mqtt3/checkconnection`
- `GET /mqtt3/topic`
- `GET /mqtt3/control`
- `GET /mqtt3/listdevicepage`
- `GET /mqtt3/devicebuckets`
- `GET /mqtt3/sensercharts`
- `GET /mqtt3/alarmdevicestatus`
- `GET /mqtt3/monitordevicegroup`
- อีกหลายตัว

**Go** (`internal/modules/iot/delivery/http/handler.go`):
```go
package http

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/render"
)

type MQTT3Handler struct {
    uc     usecase.MQTT3UseCase
    logger logger.Logger
}

func (h *MQTT3Handler) GetTopicData(w http.ResponseWriter, r *http.Request) {
    topic := r.URL.Query().Get("topic")
    if topic == "" {
        render.Render(w, r, ErrBadRequest("topic is required"))
        return
    }
    data, err := h.uc.GetTopicData(r.Context(), topic)
    if err != nil {
        render.Render(w, r, ErrInternal(err))
        return
    }
    render.JSON(w, r, map[string]interface{}{
        "status": 1,
        "data":   string(data),
    })
}

func (h *MQTT3Handler) DeviceControl(w http.ResponseWriter, r *http.Request) {
    topic := r.URL.Query().Get("topic")
    message := r.URL.Query().Get("message")
    if topic == "" || message == "" {
        render.Render(w, r, ErrBadRequest("topic and message are required"))
        return
    }
    if err := h.uc.DeviceControl(r.Context(), topic, message); err != nil {
        render.Render(w, r, ErrInternal(err))
        return
    }
    render.JSON(w, r, map[string]string{"status": "ok"})
}
```

### 2.5 การดึงข้อมูลอุปกรณ์ + MQTT + InfluxDB (listdevicepage)

ฟังก์ชันที่ซับซ้อนที่สุดคือ `listdevicepage` ซึ่งต้อง:
- ดึงรายการ device จาก PostgreSQL (พร้อม join mqtt, location)
- สำหรับแต่ละ device: ดึงข้อมูลล่าสุดจาก MQTT (หรือ cache), คำนวณ alarm ผ่าน iothelper, ดึง chart data จาก InfluxDB

แนวทางใน Go:
- ใช้ goroutine concurrent ดึงข้อมูล MQTT หลาย device พร้อมกัน (limit concurrency)
- จัดการ cache ที่ Redis สำหรับ MQTT data และ计算结果
- ใช้ `influxClient` ที่มีอยู่แล้ว query ข้อมูล chart

ตัวอย่าง pseudocode:
```go
func (u *mqtt3UseCase) GetDeviceListPage(ctx context.Context, req *DeviceListRequest) (*DeviceListPageResponse, error) {
    // 1. Query devices from DB (ด้วย GORM)
    var devices []models.Device
    u.db.Preload("Mqtt").Preload("Location").Find(&devices)

    // 2. สร้าง worker pool (semaphore)
    sem := make(chan struct{}, 5)
    var wg sync.WaitGroup
    resultCh := make(chan DeviceDetail, len(devices))

    for _, dev := range devices {
        wg.Add(1)
        go func(d models.Device) {
            defer wg.Done()
            sem <- struct{}{}
            defer func() { <-sem }()
            detail := u.fetchDeviceDetail(ctx, d)
            resultCh <- detail
        }(dev)
    }
    go func() {
        wg.Wait()
        close(resultCh)
    }()
    // 3. รวบรวมผลลัพธ์
    var details []DeviceDetail
    for d := range resultCh {
        details = append(details, d)
    }
    return &DeviceListPageResponse{Data: details}, nil
}
```

---

## 3. การจัดการ Cache และ Log

- **Redis Cache**: ใช้ client เดิม `pkg/db/redis` และ helper `pkg/redis/cache.go` (เช่น `GetCacheData`, `SetCacheData`)
- **Alarm Log**: สร้าง repository ใหม่ `internal/modules/iot/repository/alarm_log_repo.go` สำหรับบันทึก log ลง PostgreSQL (ตาราง `sd_alarm_process_log`, `sd_alarm_process_log_email` ฯลฯ)
- **Schedule Log**: ใช้ `scheduleprocesslogRepository` แนวเดียวกัน

---

## 4. หน่วยทดสอบ (Unit Test) และ Swagger

### 4.1 Unit Test

สร้างไฟล์ `*_test.go` สำหรับแต่ละ layer:
- `usecase_test.go` – mock MQTT client, Redis, DB
- `handler_test.go` – test HTTP endpoints ด้วย `httptest`

ใช้ `mockery` สร้าง mock จาก interface (ตามที่มีใน `internal/mocks`)

ตัวอย่าง:
```go
// internal/modules/iot/usecase/usecase_test.go
func TestGetTopicData_CacheHit(t *testing.T) {
    mockRedis := new(mocks.RedisClient)
    mockRedis.On("Get", mock.Anything, "mqtt_topic:test").Return(redis.NewStringResult(`"cached"`, nil))
    uc := NewMQTT3UseCase(mockMQTT, mockRedis, mockDB, mockInflux, logger)
    data, err := uc.GetTopicData(context.Background(), "test")
    assert.NoError(t, err)
    assert.Equal(t, []byte(`"cached"`), data)
}
```

### 4.2 Swagger

ใช้ `swaggo/swag`  annotate handler:
```go
// GetTopicData godoc
// @Summary      Get live or cached MQTT topic data
// @Tags         mqtt3
// @Param        topic query string true "MQTT topic"
// @Success      200 {object} map[string]interface{}
// @Router       /mqtt3/topic [get]
func (h *MQTT3Handler) GetTopicData(w http.ResponseWriter, r *http.Request) { ... }
```
จากนั้นรัน `swag init -g cmd/api/main.go` เพื่อสร้าง docs

---

## 5. Checklist การพัฒนา

- [x] สร้างโครงสร้างโฟลเดอร์ `internal/modules/iot`
- [x] สร้าง `iothelper/alarm.go` (แปลงจาก iot.helper.ts)
- [x] สร้าง `pkg/mqtt/manager.go` (MqttRequestManager)
- [x] สร้าง `usecase/usecase.go` พร้อม interface และ implementation
- [x] สร้าง `delivery/http/handler.go` และ `routes.go`
- [x] สร้าง `presenter` structs สำหรับ request/response
- [x] สร้าง `repository` สำหรับ device, alarm log, schedule
- [x] เพิ่ม routes หลักใน `server/handlers.go`
- [ ] เพิ่ม unit test (ครอบคลุมฟังก์ชันสำคัญ)
- [ ] ทดสอบ integration กับ MQTT broker จริง
- [ ] ตรวจสอบ performance (cache, concurrent)

---

## 6. ข้อควรระวังและแนวทางแก้ปัญหา

| ปัญหาที่อาจเกิด | แนวทางแก้ไข |
|----------------|--------------|
| MQTT connection หลุดบ่อย | ใช้ `AutoReconnect=true` ใน Paho client และมี health check endpoint |
| การ subscribe topic เยอะเกินไป | ใช้ `RequestManager` ที่ subscribe แบบ temporary แล้ว unsubscribe ทันที (ลดภาระ broker) |
| Alarm ส่งซ้ำถี่เกิน | ใช้ Redis เก็บ timestamp ของการส่งล่าสุด และ cooldown period (time_life) |
| การ query InfluxDB ช้า | ใช้ Redis cache ผลลัพธ์ของ query (TTL ตามต้องการ) และ pagination |
| Goroutine leak | ใช้ `context.WithTimeout` และปิด channel อย่างถูกต้อง |

---

## 7. สรุป

การแปลง MQTT3 module จาก TypeScript ไป Go สามารถทำได้โดย:
1. ย้าย logic การจัดการ MQTT request‑response ไปอยู่ใน `pkg/mqtt/manager.go`
2. ย้าย logic การคำนวณ alarm ไปอยู่ใน `internal/modules/iot/iothelper`
3. สร้าง usecase ใหม่ที่ใช้ dependency เดิม (mqtt client, redis, db, influxdb)
4. สร้าง HTTP handler ตาม endpoints ที่มี
5. เพิ่ม routes ใน `server/handlers.go`

โดยไม่มีการแก้ไขโค้ดเดิมของระบบ (อาศัยการเพิ่ม module ใหม่และใช้ interface ที่มีอยู่) จึงปลอดภัยและสามารถ deploy พร้อมกับระบบปัจจุบันได้

**ไฟล์ที่ต้องสร้างใหม่** ทั้งหมดจะอยู่ใน `internal/modules/iot` และ `pkg/mqtt/manager.go` พร้อมกับเพิ่ม routes หนึ่งบรรทัดใน `server/handlers.go` เท่านั้น


----------------------

สร้าง modules  iotservice ภาษา golang ตามโครงสร้าง icmongolang/ สร้าง api   ตามนี้   @Get('topic')
@Get('control')
    @Get('device')
	@Get('listdevicepage')
    @Get('devicebuckets')
    @Get('sensercharts')
    @Get('senserdatachart')
    @Get('senserdata')
    @Get('devicelist')
	@Get('locationdevice') 
	@Get('devicesensercharts')
    @Get('alarmdevicestatus')
    @Get('alarmdevicestatuscontrol')
    @Get('monitordevicegroup')
    @Get('monitordevicechart')
	
----------------------

เราจะสร้าง **MQTT3 Service** ใน Go ตามสถาปัตยกรรม Clean Architecture ของโปรเจกต์ `icmongolang` โดยมี endpoints ครบตามที่กำหนด

---

## โครงสร้างโฟลเดอร์ที่จะสร้าง

```
internal/modules/iot/
├── delivery/
│   └── http/
│       ├── handler.go
│       └── routes.go
├── usecase/
│   └── usecase.go
├── repository/
│   ├── device_repo.go
│   ├── alarm_log_repo.go
│   └── schedule_repo.go
├── presenter/
│   └── presenter.go
├── iothelper/
│   └── alarm.go
pkg/mqtt/
└── manager.go
```

## 1. iothelper/alarm.go – ตรรกะการคำนวณ Alarm และข้อความ

```go
// internal/modules/iot/iothelper/alarm.go
package iothelper

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type AlarmDetailDto struct {
	HardwareID        interface{}
	ValueData         interface{}
	Max               interface{}
	Min               interface{}
	StatusAlert       interface{}
	StatusWarning     interface{}
	RecoveryWarning   interface{}
	RecoveryAlert     interface{}
	DeviceName        string
	ActionName        string
	MqttName          string
	MqttControlOn     string
	MqttControlOff    string
	CountAlarm        interface{}
	Event             interface{}
	Unit              string
}

type AlarmDetailResult struct {
	Status             int    `json:"status"`
	StatusControl      int    `json:"status_control"`
	AlarmTypeId        int    `json:"alarm_type_id"`
	TypeId             int    `json:"type_id"`
	HardwareId         int    `json:"hardware_id"`
	AlarmStatusSet     int    `json:"alarm_status_set"`
	Title              string `json:"title"`
	Subject            string `json:"subject"`
	Content            string `json:"content"`
	DataAlarm          int    `json:"data_alarm"`
	EventControl       int    `json:"event_control"`
	MessageMqttControl string `json:"message_mqtt_control"`
	SensorData         interface{} `json:"sensor_data"`
	CountAlarm         int    `json:"count_alarm"`
	Unit               string `json:"unit"`
	Timestamp          string `json:"timestamp"`
}

var thaiMessages = map[string]string{
	"warning":          "คำเตือน มีความผิดปกติ",
	"critical":         "ภาวะวิกฤตต้องแก้ไขทันที",
	"recoveryWarning":  "คืนสู่ภาวะปกติ (คำเตือน)",
	"recoveryCritical": "คืนสู่ภาวะปกติ (วิกฤต)",
	"normal":           "ปกติ",
	"normal2":          "ปกติ",
	"criticalMax":      "วิกฤต มีค่าสูงเกินกำหนด",
	"criticalMin":      "วิกฤต มีค่าต่ำกว่ากำหนด",
	"normal3":          "ปกติ",
}

var englishMessages = map[string]string{
	"warning":          "Warning",
	"critical":         "Critical",
	"recoveryWarning":  "Recovery Warning",
	"recoveryCritical": "Recovery Critical",
	"normal":           "Normal",
	"normal2":          "Normal",
	"criticalMax":      "Critical! Maximum limit.",
	"criticalMin":      "Critical! Minimum limit",
	"normal3":          "Normal",
}

func toInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	case string:
		i, _ := strconv.Atoi(val)
		return i
	default:
		return 0
	}
}

func toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	default:
		return 0
	}
}

func normalizeSensorValue(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
		up := strings.ToUpper(val)
		if up == "ON" {
			return 1
		}
		if up == "OFF" {
			return 0
		}
		return val
	default:
		return val
	}
}

func processAlarmDetail(dto AlarmDetailDto, messages map[string]string) AlarmDetailResult {
	hardwareID := toInt(dto.HardwareID)
	typeID := hardwareID

	sensorValue := normalizeSensorValue(dto.ValueData)
	maxVal := toFloat(dto.Max)
	minVal := toFloat(dto.Min)
	statusAlert := toInt(dto.StatusAlert)
	statusWarning := toInt(dto.StatusWarning)
	recoveryWarning := toInt(dto.RecoveryWarning)
	recoveryAlert := toInt(dto.RecoveryAlert)
	countAlarm := toInt(dto.CountAlarm)
	event := toInt(dto.Event)

	unit := dto.Unit
	mqttName := dto.MqttName
	deviceName := dto.DeviceName
	alarmActionName := dto.ActionName
	mqttControlOn := dto.MqttControlOn
	mqttControlOff := dto.MqttControlOff
	valueAlarm := dto.ValueData
	valueRelay := dto.ValueData
	valueControlRelay := dto.ValueData

	var sensorData interface{}
	var valueData interface{}

	switch hardwareID {
	case 1:
		sensorData = dto.ValueData
		valueData = dto.ValueData
	case 2:
		if toInt(dto.ValueData) == 1 {
			sensorData = 1
			valueData = 1
			sensorValue = 1
		} else {
			sensorData = toInt(dto.ValueData)
			valueData = toInt(dto.ValueData)
			sensorValue = toInt(dto.ValueData)
		}
	case 3:
		sensorData = toInt(dto.ValueData)
		valueData = dto.ValueData
		sensorValue = dto.ValueData
	case 4:
		sensorData = dto.ValueData
		valueData = dto.ValueData
	default:
		sensorData = toInt(dto.ValueData)
		valueData = dto.ValueData
	}

	alarmStatusSet := 999
	dataAlarm := 0
	eventControl := event
	messageMqttControl := mqttControlOff
	if event == 1 {
		messageMqttControl = mqttControlOn
	}
	status := 5
	title := messages["normal"]
	subject := messages["normal"]
	content := messages["normal"] + " "

	if hardwareID == 3 && (sensorValue == 1 || sensorValue == 0 || sensorValue == "ON" || sensorValue == "OFF" || sensorValue == "on" || sensorValue == "off") {
		alarmStatusSet = 999
		title = messages["normal"]
		subject = messages["normal"]
		content = fmt.Sprintf("%s %v %s", messages["normal"], sensorValue, unit)
		status = 5
	} else if hardwareID == 4 && sensorValue != 1 {
		alarmStatusSet = 2
		title = messages["critical"]
		subject = fmt.Sprintf("%s %s %s : %v %s", mqttName, messages["critical"], deviceName, sensorValue, unit)
		content = fmt.Sprintf("%s %s %s : %s :%v %s", mqttName, alarmActionName, messages["critical"], deviceName, sensorValue, unit)
		dataAlarm = statusWarning
		status = 2
	} else if hardwareID == 4 && sensorValue == 1 {
		alarmStatusSet = 999
		title = messages["normal2"]
		subject = messages["normal2"]
		content = fmt.Sprintf("%s %v %s", messages["normal2"], sensorValue, unit)
		status = 5
	} else if maxVal != 0 && toFloat(sensorValue) >= maxVal && (hardwareID == 1 || hardwareID == 2) {
		alarmStatusSet = 2
		title = messages["criticalMax"]
		subject = fmt.Sprintf("%s %s %s : %v %s", mqttName, messages["criticalMax"], deviceName, sensorValue, unit)
		content = fmt.Sprintf("%s %s %s : %s :%v %s", mqttName, alarmActionName, messages["criticalMax"], deviceName, sensorValue, unit)
		dataAlarm = statusWarning
		status = 2
	} else if minVal != 0 && toFloat(sensorValue) <= minVal && (hardwareID == 1 || hardwareID == 2) {
		alarmStatusSet = 1
		title = messages["criticalMin"]
		subject = fmt.Sprintf("%s %s %s : %v %s", mqttName, messages["criticalMin"], deviceName, sensorValue, unit)
		content = fmt.Sprintf("%s %s %s : %s :%v %s", mqttName, alarmActionName, messages["criticalMin"], deviceName, sensorValue, unit)
		dataAlarm = statusWarning
		status = 1
	} else if hardwareID == 1 && statusWarning > 0 && toFloat(sensorValue) >= float64(statusWarning) && toFloat(sensorValue) < float64(statusAlert) {
		alarmStatusSet = 1
		title = messages["warning"]
		subject = fmt.Sprintf("%s %s : %s : %v %s", mqttName, messages["warning"], deviceName, sensorValue, unit)
		content = fmt.Sprintf("%s %s %s: %s :%v %s", mqttName, alarmActionName, messages["warning"], deviceName, sensorValue, unit)
		dataAlarm = statusWarning
		status = 1
	} else if hardwareID == 1 && statusAlert > 0 && toFloat(sensorValue) >= float64(statusAlert) {
		alarmStatusSet = 2
		title = messages["critical"]
		subject = fmt.Sprintf("%s %s : %s :%v %s", mqttName, messages["critical"], deviceName, sensorValue, unit)
		content = fmt.Sprintf("%s %s %s: %s :%v %s", mqttName, alarmActionName, messages["critical"], deviceName, sensorValue, unit)
		dataAlarm = statusAlert
		status = 2
	} else if toInt(valueAlarm) == 0 && (hardwareID == 2 || hardwareID == 3 || hardwareID == 4) {
		isCritical := hardwareID == 4
		if isCritical {
			alarmStatusSet = 2
			title = messages["critical"]
		} else {
			alarmStatusSet = 1
			title = messages["warning"]
		}
		subject = fmt.Sprintf("%s %s : %s : %v %s", mqttName, title, deviceName, sensorValue, unit)
		content = fmt.Sprintf("%s %s %s: %s :%v %s", mqttName, alarmActionName, title, deviceName, sensorValue, unit)
		if isCritical {
			dataAlarm = statusAlert
		} else {
			dataAlarm = statusWarning
		}
		status = 2
		if !isCritical {
			status = 1
		}
	} else if countAlarm >= 1 && recoveryWarning > 0 && toFloat(sensorValue) <= float64(recoveryWarning) && (hardwareID == 1 || hardwareID == 2) {
		alarmStatusSet = 3
		title = messages["recoveryWarning"]
		subject = fmt.Sprintf("%s %s : %s :%v %s", mqttName, messages["recoveryWarning"], deviceName, sensorValue, unit)
		content = fmt.Sprintf("%s %s %s: %s :%v %s", mqttName, alarmActionName, messages["recoveryWarning"], deviceName, sensorValue, unit)
		dataAlarm = recoveryWarning
		eventControl = 0
		if event == 1 {
			eventControl = 1
		}
		if event == 1 {
			messageMqttControl = mqttControlOff
		} else {
			messageMqttControl = mqttControlOn
		}
		status = 3
	} else if countAlarm >= 1 && recoveryAlert > 0 && toFloat(sensorValue) <= float64(recoveryAlert) && (hardwareID == 1 || hardwareID == 2) {
		alarmStatusSet = 4
		title = fmt.Sprintf("%s %s", mqttName, messages["recoveryCritical"])
		subject = fmt.Sprintf("%s %s :%s :%v %s", mqttName, messages["recoveryCritical"], deviceName, sensorValue, unit)
		content = fmt.Sprintf("%s %s %s :%s :%v %s", mqttName, alarmActionName, messages["recoveryCritical"], deviceName, sensorValue, unit)
		dataAlarm = recoveryAlert
		eventControl = 0
		if event == 1 {
			eventControl = 1
		}
		if event == 1 {
			messageMqttControl = mqttControlOff
		} else {
			messageMqttControl = mqttControlOn
		}
		status = 4
	} else if countAlarm >= 1 && toInt(valueAlarm) >= 1 && (hardwareID == 2 || hardwareID == 3 || hardwareID == 4) {
		alarmStatusSet = 4
		title = fmt.Sprintf("%s %s", mqttName, messages["recoveryCritical"])
		subject = fmt.Sprintf("%s %s :%s :%v %s", mqttName, messages["recoveryCritical"], deviceName, sensorValue, unit)
		content = fmt.Sprintf("%s %s %s :%s :%v %s", mqttName, alarmActionName, messages["recoveryCritical"], deviceName, sensorValue, unit)
		dataAlarm = recoveryAlert
		eventControl = 0
		if event == 1 {
			eventControl = 1
		}
		if event == 1 {
			messageMqttControl = mqttControlOff
		} else {
			messageMqttControl = mqttControlOn
		}
		status = 4
	} else {
		alarmStatusSet = 999
		title = messages["normal3"]
		subject = messages["normal3"]
		content = messages["normal"] + " "
		dataAlarm = 0
		status = 5
	}

	return AlarmDetailResult{
		Status:             status,
		StatusControl:      status,
		AlarmTypeId:        hardwareID,
		TypeId:             typeID,
		HardwareId:         hardwareID,
		AlarmStatusSet:     alarmStatusSet,
		Title:              title,
		Subject:            subject,
		Content:            content,
		DataAlarm:          dataAlarm,
		EventControl:       eventControl,
		MessageMqttControl: messageMqttControl,
		SensorData:         sensorData,
		CountAlarm:         countAlarm,
		Unit:               unit,
		Timestamp:          time.Now().Format("2006-01-02 15:04:05"),
	}
}

func AlarmDetailValidate(dto AlarmDetailDto) AlarmDetailResult {
	return processAlarmDetail(dto, thaiMessages)
}

func AlarmDetailValidateEn(dto AlarmDetailDto) AlarmDetailResult {
	return processAlarmDetail(dto, englishMessages)
}

func AlarmDetailValidateTh(dto AlarmDetailDto) AlarmDetailResult {
	return processAlarmDetail(dto, thaiMessages)
}
```

## 2. presenter/presenter.go – structs สำหรับ API Request/Response

```go
// internal/modules/iot/presenter/presenter.go
package presenter

type DeviceListRequest struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
	Bucket     string `json:"bucket"`
	HardwareId int    `json:"hardware_id"`
	TypeId     int    `json:"type_id"`
	Keyword    string `json:"keyword"`
	Lang       string `json:"lang"`
}

type DeviceDetailResponse struct {
	DeviceId   int    `json:"device_id"`
	DeviceName string `json:"device_name"`
	TypeName   string `json:"type_name"`
	ValueData  string `json:"value_data"`
	Unit       string `json:"unit"`
	Status     int    `json:"status"`
	AlarmTitle string `json:"alarm_title"`
}

type TopicDataResponse struct {
	Topic   string      `json:"topic"`
	Payload interface{} `json:"payload"`
	From    string      `json:"from"`
	Cache   bool        `json:"cache"`
}

type ControlRequest struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

type SenserChartRequest struct {
	Bucket      string `json:"bucket"`
	Measurement string `json:"measurement"`
	Field       string `json:"field"`
	Start       string `json:"start"`
	Stop        string `json:"stop"`
	Limit       int    `json:"limit"`
}

type SenserChartResponse struct {
	Data  []float64 `json:"data"`
	Date  []string  `json:"date"`
	Cache string    `json:"cache"`
}
```

## 3. repository – การดึงข้อมูลจากฐานข้อมูล

```go
// internal/modules/iot/repository/device_repo.go
package repository

import (
	"icmongolang/internal/models"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	GetDeviceByID(id int) (*models.Device, error)
	GetDevicesByBucket(bucket string) ([]models.Device, error)
	ListDevices(filter map[string]interface{}, page, pageSize int) ([]models.Device, int64, error)
}

type deviceRepo struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepo{db: db}
}

func (r *deviceRepo) GetDeviceByID(id int) (*models.Device, error) {
	var device models.Device
	err := r.db.Preload("Mqtt").Preload("Location").First(&device, id).Error
	return &device, err
}

func (r *deviceRepo) GetDevicesByBucket(bucket string) ([]models.Device, error) {
	var devices []models.Device
	err := r.db.Where("bucket = ?", bucket).Preload("Mqtt").Find(&devices).Error
	return devices, err
}

func (r *deviceRepo) ListDevices(filter map[string]interface{}, page, pageSize int) ([]models.Device, int64, error) {
	var devices []models.Device
	query := r.db.Model(&models.Device{})
	for k, v := range filter {
		query = query.Where(k+" = ?", v)
	}
	var total int64
	query.Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Preload("Mqtt").Preload("Location").Find(&devices).Error
	return devices, total, err
}
```

## 4. usecase – business logic

```go
// internal/modules/iot/usecase/usecase.go
package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"icmongolang/internal/modules/iot/iothelper"
	"icmongolang/internal/modules/iot/presenter"
	"icmongolang/internal/modules/iot/repository"
	"icmongolang/internal/models"
	"icmongolang/pkg/influxdb"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/mqtt"
	"icmongolang/pkg/redis"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MQTT3UseCase interface {
	GetTopicData(ctx context.Context, topic string) (*presenter.TopicDataResponse, error)
	DeviceControl(ctx context.Context, req *presenter.ControlRequest) error
	GetDeviceList(ctx context.Context, req *presenter.DeviceListRequest) ([]presenter.DeviceDetailResponse, int64, error)
	GetDeviceBuckets(ctx context.Context, bucket string) ([]presenter.DeviceDetailResponse, error)
	GetSenserCharts(ctx context.Context, req *presenter.SenserChartRequest) (*presenter.SenserChartResponse, error)
	// ... other methods
}

type mqtt3UseCase struct {
	deviceRepo   repository.DeviceRepository
	alarmLogRepo repository.AlarmLogRepository
	scheduleRepo repository.ScheduleRepository
	mqttClient   mqtt.Client
	reqManager   *mqtt.RequestManager
	redisClient  *redis.Client
	influxClient *influxdb.InfluxClient
	logger       logger.Logger
}

func NewMQTT3UseCase(
	deviceRepo repository.DeviceRepository,
	alarmLogRepo repository.AlarmLogRepository,
	scheduleRepo repository.ScheduleRepository,
	mqttClient mqtt.Client,
	redisClient *redis.Client,
	influxClient *influxdb.InfluxClient,
	log logger.Logger,
) MQTT3UseCase {
	return &mqtt3UseCase{
		deviceRepo:   deviceRepo,
		alarmLogRepo: alarmLogRepo,
		scheduleRepo: scheduleRepo,
		mqttClient:   mqttClient,
		reqManager:   mqtt.NewRequestManager(mqttClient, log),
		redisClient:  redisClient,
		influxClient: influxClient,
		logger:       log,
	}
}

func (u *mqtt3UseCase) GetTopicData(ctx context.Context, topic string) (*presenter.TopicDataResponse, error) {
	cacheKey := "mqtt_topic:" + topic
	cached, err := u.redisClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var payload interface{}
		json.Unmarshal(cached, &payload)
		return &presenter.TopicDataResponse{
			Topic:   topic,
			Payload: payload,
			From:    "cache",
			Cache:   true,
		}, nil
	}
	data, err := u.reqManager.GetDataFromTopic(ctx, topic, 10*time.Second)
	if err != nil {
		return nil, err
	}
	// store in cache (30 seconds)
	payloadBytes, _ := json.Marshal(data)
	u.redisClient.Set(ctx, cacheKey, payloadBytes, 30*time.Second)
	return &presenter.TopicDataResponse{
		Topic:   topic,
		Payload: data,
		From:    "mqtt",
		Cache:   false,
	}, nil
}

func (u *mqtt3UseCase) DeviceControl(ctx context.Context, req *presenter.ControlRequest) error {
	return u.mqttClient.Publish(req.Topic, 1, false, req.Message)
}

func (u *mqtt3UseCase) GetDeviceList(ctx context.Context, req *presenter.DeviceListRequest) ([]presenter.DeviceDetailResponse, int64, error) {
	filter := make(map[string]interface{})
	if req.Bucket != "" {
		filter["bucket"] = req.Bucket
	}
	if req.HardwareId != 0 {
		filter["hardware_id"] = req.HardwareId
	}
	devices, total, err := u.deviceRepo.ListDevices(filter, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}
	result := make([]presenter.DeviceDetailResponse, 0, len(devices))
	for _, dev := range devices {
		// ดึงข้อมูล MQTT ล่าสุด (ผ่าน cache)
		topicData, _ := u.GetTopicData(ctx, dev.MqttDataValue)
		var valueData string
		if topicData != nil {
			valueData = fmt.Sprintf("%v", topicData.Payload)
		}
		result = append(result, presenter.DeviceDetailResponse{
			DeviceId:   dev.DeviceID,
			DeviceName: dev.DeviceName,
			TypeName:   dev.TypeName,
			ValueData:  valueData,
			Unit:       dev.Unit,
			Status:     dev.Status,
		})
	}
	return result, total, nil
}

func (u *mqtt3UseCase) GetDeviceBuckets(ctx context.Context, bucket string) ([]presenter.DeviceDetailResponse, error) {
	devices, err := u.deviceRepo.GetDevicesByBucket(bucket)
	if err != nil {
		return nil, err
	}
	var result []presenter.DeviceDetailResponse
	for _, dev := range devices {
		result = append(result, presenter.DeviceDetailResponse{
			DeviceId:   dev.DeviceID,
			DeviceName: dev.DeviceName,
			TypeName:   dev.TypeName,
			Unit:       dev.Unit,
		})
	}
	return result, nil
}

func (u *mqtt3UseCase) GetSenserCharts(ctx context.Context, req *presenter.SenserChartRequest) (*presenter.SenserChartResponse, error) {
	// ใช้ influxClient.QueryFilterData
	params := influxdb.QueryParams{
		Measurement: req.Measurement,
		Field:       req.Field,
		Bucket:      req.Bucket,
		Start:       req.Start,
		Stop:        req.Stop,
		Limit:       req.Limit,
	}
	results, err := u.influxClient.QueryFilterData(params)
	if err != nil {
		return nil, err
	}
	var data []float64
	var date []string
	for _, r := range results {
		if val, ok := r["_value"].(float64); ok {
			data = append(data, val)
		}
		if t, ok := r["_time"].(time.Time); ok {
			date = append(date, t.Format("2006-01-02 15:04:05"))
		}
	}
	return &presenter.SenserChartResponse{
		Data:  data,
		Date:  date,
		Cache: "no cache",
	}, nil
}
```

## 5. handler – HTTP endpoints

```go
// internal/modules/iot/delivery/http/handler.go
package http

import (
	"net/http"
	"strconv"

	"icmongolang/internal/modules/iot/presenter"
	"icmongolang/internal/modules/iot/usecase"
	"icmongolang/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type MQTT3Handler struct {
	uc     usecase.MQTT3UseCase
	logger logger.Logger
}

func NewMQTT3Handler(uc usecase.MQTT3UseCase, log logger.Logger) *MQTT3Handler {
	return &MQTT3Handler{uc: uc, logger: log}
}

func (h *MQTT3Handler) GetTopicData(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		render.Render(w, r, ErrBadRequest("topic is required"))
		return
	}
	resp, err := h.uc.GetTopicData(r.Context(), topic)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

func (h *MQTT3Handler) DeviceControl(w http.ResponseWriter, r *http.Request) {
	var req presenter.ControlRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrBadRequest(err.Error()))
		return
	}
	if err := h.uc.DeviceControl(r.Context(), &req); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, map[string]string{"status": "ok"})
}

func (h *MQTT3Handler) GetDeviceList(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	req := presenter.DeviceListRequest{
		Page:       page,
		PageSize:   pageSize,
		Bucket:     r.URL.Query().Get("bucket"),
		HardwareId: func() int { v, _ := strconv.Atoi(r.URL.Query().Get("hardware_id")); return v }(),
		TypeId:     func() int { v, _ := strconv.Atoi(r.URL.Query().Get("type_id")); return v }(),
		Keyword:    r.URL.Query().Get("keyword"),
		Lang:       r.URL.Query().Get("lang"),
	}
	devices, total, err := h.uc.GetDeviceList(r.Context(), &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, map[string]interface{}{
		"data":  devices,
		"total": total,
		"page":  page,
	})
}

func (h *MQTT3Handler) GetDeviceBuckets(w http.ResponseWriter, r *http.Request) {
	bucket := r.URL.Query().Get("bucket")
	if bucket == "" {
		render.Render(w, r, ErrBadRequest("bucket is required"))
		return
	}
	devices, err := h.uc.GetDeviceBuckets(r.Context(), bucket)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, devices)
}

func (h *MQTT3Handler) GetSenserCharts(w http.ResponseWriter, r *http.Request) {
	req := presenter.SenserChartRequest{
		Bucket:      r.URL.Query().Get("bucket"),
		Measurement: r.URL.Query().Get("measurement"),
		Field:       r.URL.Query().Get("field"),
		Start:       r.URL.Query().Get("start"),
		Stop:        r.URL.Query().Get("stop"),
		Limit:       func() int { v, _ := strconv.Atoi(r.URL.Query().Get("limit")); return v }(),
	}
	if req.Bucket == "" || req.Measurement == "" {
		render.Render(w, r, ErrBadRequest("bucket and measurement are required"))
		return
	}
	resp, err := h.uc.GetSenserCharts(r.Context(), &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

// เพิ่ม handler สำหรับ endpoints อื่น ๆ ตามต้องการ (เช่น senserdatachart, devicelist, locationdevice, alarmdevicestatus, monitordevicegroup ฯลฯ)

func ErrBadRequest(msg string) render.Renderer {
	return &errResponse{HTTPStatusCode: http.StatusBadRequest, ErrorText: msg}
}
func ErrInternal(err error) render.Renderer {
	return &errResponse{HTTPStatusCode: http.StatusInternalServerError, ErrorText: err.Error()}
}
type errResponse struct {
	HTTPStatusCode int    `json:"-"`
	ErrorText      string `json:"error"`
}
func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
```

## 6. routes – ลงทะเบียน routes

```go
// internal/modules/iot/delivery/http/routes.go
package http

import (
	"icmongolang/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func MapMQTT3Routes(router chi.Router, h *MQTT3Handler, mw *middleware.MiddlewareManager) {
	router.Route("/mqtt3", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())

			r.Get("/topic", h.GetTopicData)
			r.Post("/control", h.DeviceControl)
			r.Get("/device", h.GetDeviceList)
			r.Get("/devicebuckets", h.GetDeviceBuckets)
			r.Get("/sensercharts", h.GetSenserCharts)
			// เพิ่ม endpoints อื่น ๆ ตามต้องการ
		})
	})
}
```

## 7. pkg/mqtt/manager.go – RequestManager สำหรับ MQTT

```go
// pkg/mqtt/manager.go
package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type pendingRequest struct {
	resolve chan<- interface{}
	reject  chan<- error
	timeout *time.Timer
	topic   string
}

type RequestManager struct {
	client            mqtt.Client
	mu                sync.Mutex
	pending           map[string][]*pendingRequest
	subscriptionCount map[string]int
	logger            interface{ Debugf, Infof, Errorf }
}

func NewRequestManager(client mqtt.Client, log interface{}) *RequestManager {
	rm := &RequestManager{
		client:            client,
		pending:           make(map[string][]*pendingRequest),
		subscriptionCount: make(map[string]int),
		logger:            log,
	}
	client.AddRoute("#", rm.handleIncomingMessage)
	return rm
}

func (rm *RequestManager) handleIncomingMessage(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	rm.mu.Lock()
	defer rm.mu.Unlock()
	pendings, ok := rm.pending[topic]
	if !ok || len(pendings) == 0 {
		return
	}
	req := pendings[0]
	rm.pending[topic] = pendings[1:]
	if len(rm.pending[topic]) == 0 {
		delete(rm.pending, topic)
		rm.decrementSubscription(topic)
	}
	req.timeout.Stop()
	var result interface{}
	if err := json.Unmarshal(msg.Payload(), &result); err != nil {
		result = string(msg.Payload())
	}
	req.resolve <- result
}

func (rm *RequestManager) GetDataFromTopic(ctx context.Context, topic string, timeout time.Duration) (interface{}, error) {
	resultCh := make(chan interface{}, 1)
	errCh := make(chan error, 1)
	timer := time.NewTimer(timeout)
	rm.mu.Lock()
	if _, ok := rm.pending[topic]; !ok {
		rm.incrementSubscription(topic)
	}
	req := &pendingRequest{
		resolve: resultCh,
		reject:  errCh,
		timeout: timer,
		topic:   topic,
	}
	rm.pending[topic] = append(rm.pending[topic], req)
	rm.mu.Unlock()
	go func() {
		<-timer.C
		rm.mu.Lock()
		defer rm.mu.Unlock()
		list := rm.pending[topic]
		for i, r := range list {
			if r == req {
				rm.pending[topic] = append(list[:i], list[i+1:]...)
				break
			}
		}
		if len(rm.pending[topic]) == 0 {
			delete(rm.pending, topic)
			rm.decrementSubscription(topic)
		}
		errCh <- fmt.Errorf("timeout: no message from topic %s after %v", topic, timeout)
	}()
	select {
	case err := <-errCh:
		return nil, err
	case res := <-resultCh:
		return res, nil
	}
}

func (rm *RequestManager) incrementSubscription(topic string) {
	count := rm.subscriptionCount[topic]
	if count == 0 {
		rm.client.Subscribe(topic, 0, nil)
		rm.subscriptionCount[topic] = 1
	} else {
		rm.subscriptionCount[topic] = count + 1
	}
}

func (rm *RequestManager) decrementSubscription(topic string) {
	count := rm.subscriptionCount[topic]
	if count <= 1 {
		rm.client.Unsubscribe(topic)
		delete(rm.subscriptionCount, topic)
	} else {
		rm.subscriptionCount[topic] = count - 1
	}
}
```

## 8. เพิ่ม routes ใน server/handlers.go (ส่วนที่ต้องแก้ไข)

```go
// ในฟังก์ชัน New() ของ server/handlers.go ให้เพิ่ม
import iotHttp "icmongolang/internal/modules/iot/delivery/http"
iotUsecase "icmongolang/internal/modules/iot/usecase"
iotRepo "icmongolang/internal/modules/iot/repository"

// ... ข้างในหลังจากสร้าง dependencies
deviceRepo := iotRepo.NewDeviceRepository(db)
alarmLogRepo := iotRepo.NewAlarmLogRepository(db)
scheduleRepo := iotRepo.NewScheduleRepository(db)
iotUC := iotUsecase.NewMQTT3UseCase(deviceRepo, alarmLogRepo, scheduleRepo, mqttClient, redisClient, influxClient, logger)
iotHandler := iotHttp.NewMQTT3Handler(iotUC, logger)
iotHttp.MapMQTT3Routes(apiRouter, iotHandler, mw)
```

---

## สรุป

- **iothelper/alarm.go**: ตรรกะการแจ้งเตือน คืนสถานะ alarm, ข้อความ
- **presenter**: structs สำหรับ request/response
- **repository**: ดึง device, alarm log, schedule จาก PostgreSQL
- **usecase**: logic หลัก, ใช้ MQTT RequestManager, Redis cache, InfluxDB query
- **handler**: จัดการ HTTP request/response
- **routes**: ลงทะเบียน endpoints ภายใต้ `/api/mqtt3`
- **pkg/mqtt/manager.go**: จัดการ MQTT request‑response pattern

โค้ดข้างต้นเป็นตัวอย่างที่สามารถนำไปขยายให้ครบทุก endpoint ตามที่โจทย์กำหนด (เช่น `/listdevicepage`, `/alarmdevicestatus`, `/monitordevicegroup` ฯลฯ) โดยนำแนวทางเดียวกันมาใช้

**ข้อควรระวัง**: ต้องตรวจสอบว่า models.Device, models.Mqtt มี fields ตรงกับที่ใช้ใน repository (เช่น `MqttDataValue`, `DeviceName`, `TypeName`, `Unit`) หากไม่มีต้องเพิ่มหรือปรับให้เข้ากับโครงสร้างจริงของโปรเจกต์