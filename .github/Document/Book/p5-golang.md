## การวิเคราะห์และแผนพัฒนาโมดูล MQTT และ InfluxDB ด้วย Go

จากโครงสร้างโปรเจกต์ `icmongolang` และไฟล์ตัวอย่างที่ให้มา (Go backend พร้อม MQTT client, usecase, handler) รวมถึงฟังก์ชันเสริมจาก TypeScript (`MqttRequestManager`, `getDataFromTopic`, cache, alarm validation) เราจะพัฒนาโมดูล MQTT และ InfluxDB ให้สมบูรณ์ยิ่งขึ้น โดยเน้น:

1. **MQTT Request-Response Pattern** – รองรับการส่งคำขอและรอรับ response ตาม pattern ที่มีใน TypeScript (subscribe ชั่วคราว, timeout, unsubscribe อัตโนมัติ)
2. **Cache Layer** – สำหรับลดการเรียก MQTT ซ้ำ (Redis cache)
3. **InfluxDB Integration** – บันทึกข้อมูล MQTT และ alarm logs ลง InfluxDB (time-series)
4. **HTTP API** – เพิ่ม endpoint สำหรับดึงข้อมูล MQTT แบบ request-response (คล้าย `/v1/mqtt2/topic`) และจัดการ cache
5. **Alarm Processing** – ปรับใช้ logic การประเมิน alarm จาก `iot.helper.ts` มาเป็น Go (ใช้ struct และฟังก์ชัน)

---

## 1. ปรับปรุง MQTT Client ให้รองรับ Request-Response

ไฟล์ `pkg/mqtt/client.go` มี `RequestData` และ `requestManager` อยู่แล้ว แต่ยังไม่สมบูรณ์ (handler ยังไม่จัดการ pending ตาม topic ที่ response) และยังไม่มีฟังก์ชัน `GetDataFromTopic` แบบ subscribe ชั่วคราว.

### เพิ่มเมธอด `GetDataFromTopic` ใน `Client` interface

```go
// pkg/mqtt/client.go (เพิ่มใน interface)
type Client interface {
    // ... existing methods ...
    // GetDataFromTopic subscribes to a topic, waits for a single message, then unsubscribes.
    GetDataFromTopic(ctx context.Context, topic string, timeout time.Duration) ([]byte, error)
}
```

### Implement `GetDataFromTopic` โดยใช้ channel และ context

```go
// pkg/mqtt/client.go
func (c *mqttClient) GetDataFromTopic(ctx context.Context, topic string, timeout time.Duration) ([]byte, error) {
    if c.client == nil || !c.client.IsConnected() {
        return nil, fmt.Errorf("MQTT client not connected")
    }
    
    result := make(chan []byte, 1)
    errCh := make(chan error, 1)
    
    // subscribe temporarily
    token := c.client.Subscribe(topic, 0, func(cl mqtt.Client, msg mqtt.Message) {
        select {
        case result <- msg.Payload():
        default:
        }
    })
    if token.Wait() && token.Error() != nil {
        return nil, token.Error()
    }
    
    // ensure unsubscribe after
    defer c.client.Unsubscribe(topic)
    
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    case <-time.After(timeout):
        return nil, fmt.Errorf("timeout waiting for message on %s", topic)
    case payload := <-result:
        return payload, nil
    }
}
```

> **หมายเหตุ**: ฟังก์ชันนี้จะ subscribe, รอรับ message แรก, แล้ว unsubscribe ทันที. เหมาะสำหรับ request-response แบบครั้งเดียว.

---

## 2. เพิ่ม HTTP Handler สำหรับดึง MQTT Data พร้อม Cache (คล้าย TypeScript)

สร้าง handler ใหม่ใน `internal/modules/mqtt/delivery/http/handler.go` (หรือแยกไฟล์) สำหรับ endpoint:

```
GET /mqtt/topic/data?topic=BAACTW05/DATA
```

### presenter เพิ่ม struct response

```go
// internal/modules/mqtt/presenter/presenter.go
type MQTTTopicDataResponse struct {
    StatusCode int         `json:"statuscode"`
    Code       int         `json:"code"`
    Topic      string      `json:"topic"`
    Payload    interface{} `json:"payload"`
    From       string      `json:"from"` // "cache" or "mqtt"
    Timestamp  string      `json:"timestamp"`
    Message    string      `json:"message"`
}
```

### handler method

```go
// internal/modules/mqtt/delivery/http/handler.go
func (h *MQTTHandler) GetTopicData(w http.ResponseWriter, r *http.Request) {
    topic := r.URL.Query().Get("topic")
    if topic == "" {
        render.Render(w, r, ErrInvalidRequestString("topic is required"))
        return
    }
    
    ctx := r.Context()
    // 1. check cache (ใช้ redis หรือ cache interface)
    cacheKey := "mqtt_topic_data:" + md5(topic)
    var cachedData []byte
    if err := h.cache.Get(ctx, cacheKey, &cachedData); err == nil && cachedData != nil {
        render.JSON(w, r, presenter.MQTTTopicDataResponse{
            StatusCode: http.StatusOK,
            Code:       200,
            Topic:      topic,
            Payload:    string(cachedData),
            From:       "cache",
            Timestamp:  time.Now().Format(time.RFC3339),
            Message:    "from cache",
        })
        return
    }
    
    // 2. get from MQTT (request-response)
    data, err := h.uc.GetTopicData(ctx, topic, 10*time.Second)
    if err != nil {
        h.logger.Errorf("GetTopicData error: %v", err)
        render.Render(w, r, ErrInternal(err))
        return
    }
    
    // 3. save cache (TTL 60 sec)
    _ = h.cache.Set(ctx, cacheKey, data, 60*time.Second)
    
    // 4. parse CSV/JSON as needed (like TypeScript version)
    payload := processPayload(data) // แยก comma หรือแปลง JSON
    
    render.JSON(w, r, presenter.MQTTTopicDataResponse{
        StatusCode: http.StatusOK,
        Code:       200,
        Topic:      topic,
        Payload:    payload,
        From:       "mqtt",
        Timestamp:  time.Now().Format(time.RFC3339),
        Message:    "live data",
    })
}

// processPayload ทำ类似 TypeScript: ถ้า JSON → Object.values().join(',') else string
func processPayload(data []byte) interface{} {
    str := strings.TrimSpace(string(data))
    // ถ้าเป็น JSON object ให้เอาค่า value มา join
    var jsonObj map[string]interface{}
    if err := json.Unmarshal(data, &jsonObj); err == nil {
        values := make([]string, 0, len(jsonObj))
        for _, v := range jsonObj {
            values = append(values, fmt.Sprintf("%v", v))
        }
        return strings.Join(values, ",")
    }
    // ถ้าเป็น array string ธรรมดา
    if strings.Contains(str, ",") {
        return strings.Split(str, ",")
    }
    return str
}
```

ต้องเพิ่ม `cache` dependency ใน `MQTTHandler` และ `usecase`.

---

## 3. โมดูล InfluxDB

สร้าง package `pkg/influxdb/client.go` สำหรับเชื่อมต่อ InfluxDB และฟังก์ชันเขียนข้อมูล.

### โครงสร้าง

```go
// pkg/influxdb/client.go
package influxdb

import (
    "context"
    "time"
    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"
)

type Client interface {
    WritePoint(measurement string, tags map[string]string, fields map[string]interface{}, t time.Time) error
    Close()
}

type influxClient struct {
    client   influxdb2.Client
    writeAPI api.WriteAPI
}

func New(cfg *config.InfluxDBConfig) (Client, error) {
    client := influxdb2.NewClient(cfg.URL, cfg.Token)
    // test connection
    _, err := client.Health(context.Background())
    if err != nil {
        return nil, err
    }
    writeAPI := client.WriteAPI(cfg.Org, cfg.Bucket)
    return &influxClient{
        client:   client,
        writeAPI: writeAPI,
    }, nil
}

func (ic *influxClient) WritePoint(measurement string, tags map[string]string, fields map[string]interface{}, t time.Time) error {
    p := influxdb2.NewPoint(measurement, tags, fields, t)
    ic.writeAPI.WritePoint(p)
    return nil
}

func (ic *influxClient) Close() {
    ic.writeAPI.Flush()
    ic.client.Close()
}
```

### เพิ่มใน usecase สำหรับบันทึก MQTT data และ alarm logs

ใน `internal/modules/mqtt/usecase/usecase.go` ปรับ constructor ให้รับ `influxClient` และมี method `SaveMQTTDataToInflux`.

```go
type mqttUseCase struct {
    // ... existing fields
    influxClient influxdb.Client
}

func NewMQTTUseCaseWithClientAndInflux(
    client mqttPkg.Client,
    cfg *config.Config,
    log logger.Logger,
    wsHub *websocket.Hub,
    db *gorm.DB,
    alarmLogRepo AlarmLogRepository,
    influxClient influxdb.Client,
) (MQTTUseCase, error) {
    // ...
}

// เมื่อได้รับ message ใน Subscribe handler ให้บันทึกลง InfluxDB
func (u *mqttUseCase) handleMessage(msg mqtt.Message) {
    // ... existing broadcast to ws
    // บันทึกลง InfluxDB
    go u.saveToInflux(msg.Topic(), msg.Payload())
}

func (u *mqttUseCase) saveToInflux(topic string, payload []byte) {
    tags := map[string]string{"topic": topic}
    fields := map[string]interface{}{
        "payload": string(payload),
        "length":  len(payload),
    }
    // ถ้า payload เป็น JSON numeric ให้แตก fields อัตโนมัติ
    var data map[string]interface{}
    if err := json.Unmarshal(payload, &data); err == nil {
        for k, v := range data {
            if _, ok := v.(float64); ok {
                fields[k] = v
            }
        }
    }
    _ = u.influxClient.WritePoint("mqtt_messages", tags, fields, time.Now())
}
```

---

## 4. ปรับปรุง Alarm Logic (จาก iot.helper.ts มาเป็น Go)

สร้าง package `internal/modules/alarm/processor.go` สำหรับฟังก์ชัน `AlarmDetailValidate`.

```go
// internal/modules/alarm/processor.go
package alarm

type AlarmDetailDto struct {
    HardwareID        int     `json:"hardware_id"`
    ValueData         interface{} `json:"value_data"`
    ValueAlarm        interface{} `json:"value_alarm"`
    Max               *float64    `json:"max"`
    Min               *float64    `json:"min"`
    StatusAlert       float64     `json:"status_alert"`
    StatusWarning     float64     `json:"status_warning"`
    RecoveryWarning   float64     `json:"recovery_warning"`
    RecoveryAlert     float64     `json:"recovery_alert"`
    CountAlarm        float64     `json:"count_alarm"`
    Event             int         `json:"event"`
    Unit              string      `json:"unit"`
    DeviceName        string      `json:"device_name"`
    ActionName        string      `json:"action_name"`
    MqttControlOn     string      `json:"mqtt_control_on"`
    MqttControlOff    string      `json:"mqtt_control_off"`
    // ... fields
}

type AlarmResult struct {
    Status              int     `json:"status"` // 1 warning,2 critical,3 recovery warning,4 recovery critical,5 normal
    Title               string  `json:"title"`
    Subject             string  `json:"subject"`
    Content             string  `json:"content"`
    DataAlarm           float64 `json:"data_alarm"`
    EventControl        int     `json:"event_control"`
    MessageMqttControl  string  `json:"message_mqtt_control"`
}

func AlarmDetailValidate(dto AlarmDetailDto, lang string) AlarmResult {
    // implement logic based on hardware_id, value comparisons
    // similar to TypeScript version
    // return appropriate status and messages
}
```

แล้วเรียกใช้ใน MQTT handler เมื่อ topic มี keyword alarm.

---

## 5. การเชื่อมต่อทุกอย่างใน main.go

```go
// cmd/api/main.go
func main() {
    // ... load config, logger, db, redis, mqtt client
    influxClient, _ := influxdb.New(&cfg.InfluxDB)
    mqttClient := mqtt.New(&cfg.MQTT, logger)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := mqttClient.Connect(ctx); err != nil {
        logger.Fatal("mqtt connect error", err)
    }
    
    mqttUseCase, _ := usecase.NewMQTTUseCaseWithClientAndInflux(mqttClient, cfg, logger, wsHub, db, alarmLogRepo, influxClient)
    mqttHandler := http.NewMQTTHandler(mqttUseCase, logger)
    
    // routes
    r := chi.NewRouter()
    // ... middlewares
    r.Route("/mqtt", func(r chi.Router) {
        r.Get("/topic/data", mqttHandler.GetTopicData)
        r.Post("/publish", mqttHandler.Publish)
        // ... others
    })
    // start server
}
```

---

## สรุปสิ่งที่ต้องทำเพิ่ม

| Component | Action |
|-----------|--------|
| `pkg/mqtt/client.go` | เพิ่ม `GetDataFromTopic` method |
| `internal/modules/mqtt/delivery/http/handler.go` | เพิ่ม `GetTopicData` handler + cache |
| `internal/modules/mqtt/usecase/usecase.go` | เพิ่ม `GetTopicData` method และบันทึก InfluxDB |
| `pkg/influxdb/client.go` | สร้างใหม่ทั้งหมด |
| `internal/modules/alarm/processor.go` | ย้าย logic alarm จาก TypeScript |
| `config/config.go` | เพิ่ม struct สำหรับ InfluxDB |
| `cmd/api/main.go` | inject influx client |

นี่คือ blueprint สำหรับพัฒนาโมดูล MQTT และ InfluxDB ใน Go ตามโครงสร้างที่มีอยู่ ให้สามารถทำงานแบบ request-response, caching, และ time-series logging ได้ครบถ้วน.