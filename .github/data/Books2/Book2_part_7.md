# 📘 เล่มที่ 4: การพัฒนาเว็บและเครือข่าย — บทที่ 5: MQTT และ IoT

---

## บทที่ 5: MQTT และ IoT

---

### 5.1 MQTT คืออะไร?

#### MQTT Protocol

**MQTT (Message Queuing Telemetry Transport)** เป็นโปรโตคอลการสื่อสารแบบ **lightweight** ที่ออกแบบมาสำหรับ:
- อุปกรณ์ IoT (Internet of Things)
- เซ็นเซอร์และอุปกรณ์ที่มีทรัพยากรจำกัด
- การสื่อสารแบบ Real-time
- เครือข่ายที่มีแบนด์วิดท์ต่ำหรือไม่เสถียร

#### คุณสมบัติสำคัญ

| คุณสมบัติ | คำอธิบาย |
|-----------|----------|
| **Publish/Subscribe** | ผู้ส่งและผู้รับไม่รู้จักกันโดยตรง |
| **Lightweight** | Header ขนาดเล็ก (2 ไบต์) |
| **QoS (Quality of Service)** | 3 ระดับ: 0, 1, 2 |
| **Last Will** | แจ้งเตือนเมื่อ Client หลุด |
| **Retained Message** | เก็บข้อความล่าสุด |
| **Secure** | รองรับ TLS/SSL |

#### MQTT Architecture

```mermaid
flowchart LR
    subgraph Publishers
        S1[Temperature Sensor]
        S2[Humidity Sensor]
        S3[Light Sensor]
    end
    
    subgraph Broker
        B[MQTT Broker<br/>Mosquitto/EMQX]
    end
    
    subgraph Subscribers
        D1[Database]
        D2[Dashboard]
        D3[Mobile App]
        D4[Alert System]
    end
    
    S1 -->|publish: sensors/temp| B
    S2 -->|publish: sensors/humidity| B
    S3 -->|publish: sensors/light| B
    
    B -->|subscribe: sensors/#| D1
    B -->|subscribe: sensors/#| D2
    B -->|subscribe: sensors/temp| D3
    B -->|subscribe: sensors/temp| D4
```

#### Quality of Service (QoS)

| QoS Level | คำอธิบาย | ใช้เมื่อ |
|-----------|----------|---------|
| **0 (At most once)** | ส่งครั้งเดียว ไม่รับประกัน | ข้อมูลที่ไม่สำคัญ |
| **1 (At least once)** | ส่งซ้ำจนกว่าจะได้รับ ACK | ข้อมูลสำคัญปานกลาง |
| **2 (Exactly once)** | รับประกันว่าส่งถึงครั้งเดียว | ข้อมูลที่สำคัญมาก |

#### MQTT vs HTTP

| คุณสมบัติ | MQTT | HTTP |
|-----------|------|------|
| Protocol | Publish/Subscribe | Request/Response |
| Header Size | 2 ไบต์ | > 200 ไบต์ |
| Bandwidth | ต่ำ | สูง |
| Latency | ต่ำ | ปานกลาง |
| Battery Usage | ต่ำ | สูง |
| Real-time | ดีเยี่ยม | ปานกลาง |
| Use Case | IoT, Sensor | Web API |

---

### 5.2 การเชื่อมต่อ MQTT Broker

#### การติดตั้ง MQTT Broker

**1. Mosquitto (เปิด Source)**

```bash
# Ubuntu/Debian
sudo apt-get install mosquitto mosquitto-clients

# macOS
brew install mosquitto

# Windows (ด้วย Chocolatey)
choco install mosquitto

# หรือใช้ Docker
docker run -d -p 1883:1883 -p 9001:9001 eclipse-mosquitto
```

**2. EMQX (Enterprise Ready)**

```bash
# Docker
docker run -d --name emqx -p 1883:1883 -p 8083:8083 -p 8084:8084 -p 8883:8883 -p 18083:18083 emqx/emqx:latest

# Dashboard: http://localhost:18083
# Username: admin / Password: public
```

**3. Cloud MQTT Broker**

- **HiveMQ Cloud**: https://www.hivemq.com/cloud/
- **AWS IoT Core**: Amazon Web Services
- **Azure IoT Hub**: Microsoft Azure
- **EMQX Cloud**: https://www.emqx.com/cloud

#### การติดตั้ง MQTT Client Library

```bash
go get github.com/eclipse/paho.mqtt.golang
```

#### การเชื่อมต่อ MQTT Broker

```go
package main

import (
    "fmt"
    "log"
    "time"

    mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
    // 1. สร้าง Client Options
    opts := mqtt.NewClientOptions()
    opts.AddBroker("tcp://localhost:1883")
    opts.SetClientID("go-client-1")
    opts.SetUsername("admin")
    opts.SetPassword("password")
    opts.SetCleanSession(true)
    opts.SetAutoReconnect(true)
    opts.SetConnectRetry(true)
    opts.SetConnectRetryInterval(5 * time.Second)
    opts.SetKeepAlive(60 * time.Second)
    opts.SetPingTimeout(10 * time.Second)
    opts.SetConnectionLostHandler(onConnectionLost)
    opts.SetOnConnectHandler(onConnect)
    
    // 2. สร้าง Client
    client := mqtt.NewClient(opts)
    
    // 3. เชื่อมต่อ
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        log.Fatal("Failed to connect:", token.Error())
    }
    
    fmt.Println("Connected to MQTT Broker")
    
    // 4. ส่งข้อความ
    publish(client)
    
    // 5. รอ
    time.Sleep(5 * time.Second)
    
    // 6. ยกเลิกการเชื่อมต่อ
    client.Disconnect(250)
    fmt.Println("Disconnected")
}

func onConnect(client mqtt.Client) {
    fmt.Println("Connected to MQTT Broker")
}

func onConnectionLost(client mqtt.Client, err error) {
    fmt.Printf("Connection lost: %v\n", err)
}

func publish(client mqtt.Client) {
    token := client.Publish("test/topic", 1, false, "Hello MQTT")
    token.Wait()
    if token.Error() != nil {
        log.Println("Publish error:", token.Error())
    } else {
        fmt.Println("Message published")
    }
}
```

#### การตั้งค่า TLS/SSL

```go
func main() {
    opts := mqtt.NewClientOptions()
    opts.AddBroker("ssl://localhost:8883")
    opts.SetClientID("secure-client")
    
    // TLS Configuration
    tlsConfig := &tls.Config{
        InsecureSkipVerify: false,
        MinVersion:         tls.VersionTLS12,
    }
    
    // Load CA Certificate
    caCert, err := os.ReadFile("ca.crt")
    if err == nil {
        caPool := x509.NewCertPool()
        caPool.AppendCertsFromPEM(caCert)
        tlsConfig.RootCAs = caPool
    }
    
    // Load Client Certificate (ถ้ามี)
    cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
    if err == nil {
        tlsConfig.Certificates = []tls.Certificate{cert}
    }
    
    opts.SetTLSConfig(tlsConfig)
    
    // WebSocket + TLS
    opts.AddBroker("wss://localhost:8083/mqtt")
    
    client := mqtt.NewClient(opts)
    // ... connect
}
```

---

### 5.3 การส่งและรับข้อความ MQTT

#### Publish (ส่งข้อความ)

```go
package main

import (
    "encoding/json"
    "fmt"
    "time"

    mqtt "github.com/eclipse/paho.mqtt.golang"
)

type SensorData struct {
    SensorID  string    `json:"sensor_id"`
    Value     float64   `json:"value"`
    Unit      string    `json:"unit"`
    Timestamp time.Time `json:"timestamp"`
}

func publishSensorData(client mqtt.Client) {
    data := SensorData{
        SensorID:  "sensor-001",
        Value:     25.5,
        Unit:      "°C",
        Timestamp: time.Now(),
    }
    
    payload, err := json.Marshal(data)
    if err != nil {
        log.Println("JSON marshal error:", err)
        return
    }
    
    token := client.Publish("sensors/temperature", 1, false, payload)
    token.Wait()
    
    if token.Error() == nil {
        fmt.Println("Sensor data published:", string(payload))
    }
}
```

#### Subscribe (รับข้อความ)

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
    // ... connect
    
    // Subscribe
    if token := client.Subscribe("sensors/#", 1, messageHandler); token.Wait() && token.Error() != nil {
        log.Fatal("Subscribe error:", token.Error())
    }
    
    fmt.Println("Subscribed to sensors/#")
    
    // Keep alive
    select {}
}

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
    topic := msg.Topic()
    payload := msg.Payload()
    
    fmt.Printf("Received message\n")
    fmt.Printf("  Topic: %s\n", topic)
    fmt.Printf("  QoS: %d\n", msg.Qos())
    fmt.Printf("  Payload: %s\n", string(payload))
    
    // Parse JSON
    var data SensorData
    if err := json.Unmarshal(payload, &data); err != nil {
        log.Println("JSON unmarshal error:", err)
        return
    }
    
    fmt.Printf("  Sensor: %s, Value: %.2f%s\n", 
        data.SensorID, data.Value, data.Unit)
}
```

#### ระบบ MQTT Manager

```go
// mqtt/manager.go
package mqtt

import (
    "encoding/json"
    "fmt"
    "log"
    "sync"
    "time"

    mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MessageHandler func(topic string, payload []byte)

type Manager struct {
    client   mqtt.Client
    handlers map[string]MessageHandler
    mu       sync.RWMutex
    config   Config
}

type Config struct {
    Broker    string
    ClientID  string
    Username  string
    Password  string
    QOS       byte
    KeepAlive int
    TLS       bool
}

func NewManager(cfg Config) *Manager {
    m := &Manager{
        handlers: make(map[string]MessageHandler),
        config:   cfg,
    }
    
    opts := mqtt.NewClientOptions()
    opts.AddBroker(cfg.Broker)
    opts.SetClientID(cfg.ClientID)
    opts.SetUsername(cfg.Username)
    opts.SetPassword(cfg.Password)
    opts.SetKeepAlive(time.Duration(cfg.KeepAlive) * time.Second)
    opts.SetAutoReconnect(true)
    opts.SetCleanSession(true)
    opts.SetOnConnectHandler(m.onConnect)
    opts.SetConnectionLostHandler(m.onConnectionLost)
    
    if cfg.TLS {
        tlsConfig := &tls.Config{
            InsecureSkipVerify: false,
            MinVersion:         tls.VersionTLS12,
        }
        opts.SetTLSConfig(tlsConfig)
    }
    
    m.client = mqtt.NewClient(opts)
    return m
}

func (m *Manager) Connect() error {
    token := m.client.Connect()
    if token.Wait() && token.Error() != nil {
        return token.Error()
    }
    log.Println("MQTT Connected")
    return nil
}

func (m *Manager) Disconnect() {
    m.client.Disconnect(250)
    log.Println("MQTT Disconnected")
}

func (m *Manager) onConnect(client mqtt.Client) {
    log.Println("MQTT Connected")
    m.mu.RLock()
    defer m.mu.RUnlock()
    for topic := range m.handlers {
        m.subscribe(topic)
    }
}

func (m *Manager) onConnectionLost(client mqtt.Client, err error) {
    log.Printf("MQTT Connection Lost: %v", err)
}

func (m *Manager) Publish(topic string, qos byte, retained bool, payload interface{}) error {
    var data []byte
    switch v := payload.(type) {
    case string:
        data = []byte(v)
    case []byte:
        data = v
    default:
        var err error
        data, err = json.Marshal(payload)
        if err != nil {
            return err
        }
    }
    
    token := m.client.Publish(topic, qos, retained, data)
    token.Wait()
    return token.Error()
}

func (m *Manager) Subscribe(topic string, handler MessageHandler) error {
    m.mu.Lock()
    m.handlers[topic] = handler
    m.mu.Unlock()
    
    if m.client.IsConnected() {
        return m.subscribe(topic)
    }
    return nil
}

func (m *Manager) subscribe(topic string) error {
    token := m.client.Subscribe(topic, m.config.QOS, m.createHandler(topic))
    token.Wait()
    if token.Error() != nil {
        return token.Error()
    }
    log.Printf("Subscribed to: %s", topic)
    return nil
}

func (m *Manager) Unsubscribe(topic string) error {
    m.mu.Lock()
    delete(m.handlers, topic)
    m.mu.Unlock()
    
    token := m.client.Unsubscribe(topic)
    token.Wait()
    return token.Error()
}

func (m *Manager) createHandler(topic string) mqtt.MessageHandler {
    return func(client mqtt.Client, msg mqtt.Message) {
        m.mu.RLock()
        handler, ok := m.handlers[topic]
        m.mu.RUnlock()
        
        if ok {
            handler(msg.Topic(), msg.Payload())
        }
    }
}

// Wildcard subscription
func (m *Manager) SubscribeWildcard(pattern string, handler MessageHandler) error {
    return m.Subscribe(pattern, handler)
}

// ตัวอย่างการใช้งาน
func (m *Manager) PublishJSON(topic string, data interface{}) error {
    return m.Publish(topic, m.config.QOS, false, data)
}
```

---

### 5.4 Go IoT Platform — ภาพรวม

#### สถาปัตยกรรม IoT Platform

```mermaid
flowchart TB
    subgraph Devices
        D1[ESP8266]
        D2[Raspberry Pi]
        D3[Arduino]
        D4[Industrial Sensor]
    end
    
    subgraph Gateway
        G[MQTT Broker<br/>EMQX/Mosquitto]
    end
    
    subgraph Go Services
        GS1[Device Manager<br/>go-iot]
        GS2[Data Processor<br/>go-processor]
        GS3[Alert Engine<br/>go-alert]
        GS4[API Gateway<br/>go-api]
    end
    
    subgraph Storage
        S1[(TimeSeries DB<br/>TimescaleDB)]
        S2[(Cache<br/>Redis)]
    end
    
    subgraph Frontend
        F1[Dashboard<br/>React/Vue]
        F2[Mobile App]
    end
    
    D1 --> G
    D2 --> G
    D3 --> G
    D4 --> G
    
    G --> GS1
    G --> GS2
    
    GS1 --> S1
    GS2 --> S1
    GS2 --> GS3
    
    GS3 -->|Alert| F2
    
    GS4 --> S1
    GS4 --> S2
    F1 --> GS4
```

#### โครงสร้างโปรเจกต์ IoT Platform

```
go-iot-platform/
├── cmd/
│   ├── broker/
│   │   └── main.go
│   ├── processor/
│   │   └── main.go
│   └── api/
│       └── main.go
├── internal/
│   ├── device/
│   │   ├── manager.go
│   │   ├── registry.go
│   │   └── types.go
│   ├── mqtt/
│   │   ├── client.go
│   │   └── handler.go
│   ├── storage/
│   │   ├── timescale.go
│   │   └── redis.go
│   ├── alert/
│   │   ├── engine.go
│   │   └── rule.go
│   └── api/
│       ├── handler.go
│       └── middleware.go
├── pkg/
│   ├── models/
│   │   ├── device.go
│   │   ├── telemetry.go
│   │   └── alert.go
│   └── utils/
│       └── logger.go
├── config/
│   ├── config.go
│   └── config.yaml
├── web/
│   ├── dashboard/
│   └── static/
├── docker-compose.yaml
├── go.mod
└── Makefile
```

---

### 5.5 MQTT Client Management

#### Device Manager

```go
// internal/device/manager.go
package device

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "time"

    "go-iot-platform/internal/mqtt"
    "go-iot-platform/internal/storage"
)

type Device struct {
    ID         string            `json:"id"`
    Name       string            `json:"name"`
    Type       string            `json:"type"`
    Location   string            `json:"location"`
    Status     string            `json:"status"` // online, offline, error
    LastSeen   time.Time         `json:"last_seen"`
    CreatedAt  time.Time         `json:"created_at"`
    UpdatedAt  time.Time         `json:"updated_at"`
    Metadata   map[string]string `json:"metadata"`
    Config     map[string]interface{} `json:"config"`
}

type Telemetry struct {
    DeviceID  string                 `json:"device_id"`
    Timestamp time.Time              `json:"timestamp"`
    Type      string                 `json:"type"`
    Values    map[string]interface{} `json:"values"`
}

type DeviceManager struct {
    devices    map[string]*Device
    mu         sync.RWMutex
    mqttClient *mqtt.Manager
    storage    *storage.TimescaleDB
    redis      *storage.Redis
    ctx        context.Context
    cancel     context.CancelFunc
}

func NewDeviceManager(mqttClient *mqtt.Manager, db *storage.TimescaleDB, redis *storage.Redis) *DeviceManager {
    ctx, cancel := context.WithCancel(context.Background())
    dm := &DeviceManager{
        devices:    make(map[string]*Device),
        mqttClient: mqttClient,
        storage:    db,
        redis:      redis,
        ctx:        ctx,
        cancel:     cancel,
    }
    
    return dm
}

func (dm *DeviceManager) Start() error {
    // Subscribe to device topics
    topics := []string{
        "devices/+/telemetry",
        "devices/+/status",
        "devices/+/event",
        "devices/+/config/response",
    }
    
    for _, topic := range topics {
        if err := dm.mqttClient.Subscribe(topic, dm.handleMessage); err != nil {
            return fmt.Errorf("failed to subscribe to %s: %w", topic, err)
        }
    }
    
    // Start health check
    go dm.healthCheck()
    
    return nil
}

func (dm *DeviceManager) Stop() {
    dm.cancel()
}

func (dm *DeviceManager) handleMessage(topic string, payload []byte) {
    parts := strings.Split(topic, "/")
    if len(parts) < 3 {
        return
    }
    
    deviceID := parts[1]
    msgType := parts[2]
    
    switch msgType {
    case "telemetry":
        dm.handleTelemetry(deviceID, payload)
    case "status":
        dm.handleStatus(deviceID, payload)
    case "event":
        dm.handleEvent(deviceID, payload)
    case "config":
        dm.handleConfigResponse(deviceID, payload)
    }
}

func (dm *DeviceManager) handleTelemetry(deviceID string, payload []byte) {
    var telemetry Telemetry
    if err := json.Unmarshal(payload, &telemetry); err != nil {
        log.Printf("Failed to parse telemetry: %v", err)
        return
    }
    
    // Update last seen
    dm.updateDeviceStatus(deviceID, "online")
    
    // Store telemetry
    if err := dm.storage.InsertTelemetry(telemetry); err != nil {
        log.Printf("Failed to store telemetry: %v", err)
        return
    }
    
    // Cache latest value
    dm.redis.SetTelemetry(deviceID, telemetry)
    
    log.Printf("Telemetry from %s: %v", deviceID, telemetry.Values)
}

func (dm *DeviceManager) handleStatus(deviceID string, payload []byte) {
    var status struct {
        Status  string `json:"status"`
        Message string `json:"message,omitempty"`
    }
    if err := json.Unmarshal(payload, &status); err != nil {
        log.Printf("Failed to parse status: %v", err)
        return
    }
    
    dm.updateDeviceStatus(deviceID, status.Status)
}

func (dm *DeviceManager) handleEvent(deviceID string, payload []byte) {
    var event struct {
        Event   string                 `json:"event"`
        Data    map[string]interface{} `json:"data"`
        Time    time.Time              `json:"time"`
    }
    if err := json.Unmarshal(payload, &event); err != nil {
        log.Printf("Failed to parse event: %v", err)
        return
    }
    
    log.Printf("Event from %s: %s", deviceID, event.Event)
}

func (dm *DeviceManager) handleConfigResponse(deviceID string, payload []byte) {
    // Process config response
    log.Printf("Config response from %s: %s", deviceID, string(payload))
}

func (dm *DeviceManager) updateDeviceStatus(deviceID, status string) {
    dm.mu.Lock()
    defer dm.mu.Unlock()
    
    if device, ok := dm.devices[deviceID]; ok {
        device.Status = status
        device.LastSeen = time.Now()
        dm.storage.UpdateDeviceStatus(deviceID, status)
    }
}

func (dm *DeviceManager) RegisterDevice(device *Device) error {
    dm.mu.Lock()
    defer dm.mu.Unlock()
    
    device.CreatedAt = time.Now()
    device.UpdatedAt = time.Now()
    device.Status = "pending"
    
    dm.devices[device.ID] = device
    
    // Save to database
    if err := dm.storage.CreateDevice(device); err != nil {
        return err
    }
    
    // Send configuration to device
    dm.sendConfig(device.ID, device.Config)
    
    return nil
}

func (dm *DeviceManager) sendConfig(deviceID string, config map[string]interface{}) {
    topic := fmt.Sprintf("devices/%s/config", deviceID)
    payload, _ := json.Marshal(config)
    
    dm.mqttClient.Publish(topic, 1, false, payload)
}

func (dm *DeviceManager) GetDevice(id string) (*Device, error) {
    dm.mu.RLock()
    defer dm.mu.RUnlock()
    
    if device, ok := dm.devices[id]; ok {
        return device, nil
    }
    
    return dm.storage.GetDevice(id)
}

func (dm *DeviceManager) ListDevices() ([]*Device, error) {
    dm.mu.RLock()
    defer dm.mu.RUnlock()
    
    devices := make([]*Device, 0, len(dm.devices))
    for _, d := range dm.devices {
        devices = append(devices, d)
    }
    return devices, nil
}

func (dm *DeviceManager) healthCheck() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-dm.ctx.Done():
            return
        case <-ticker.C:
            dm.checkDeviceHealth()
        }
    }
}

func (dm *DeviceManager) checkDeviceHealth() {
    dm.mu.RLock()
    defer dm.mu.RUnlock()
    
    timeout := 5 * time.Minute
    now := time.Now()
    
    for id, device := range dm.devices {
        if now.Sub(device.LastSeen) > timeout {
            if device.Status != "offline" {
                device.Status = "offline"
                log.Printf("Device %s is offline", id)
                
                // Emit alert
                dm.emitAlert(Alert{
                    DeviceID: id,
                    Type:     "device_offline",
                    Message:  fmt.Sprintf("Device %s has been offline for %v", id, now.Sub(device.LastSeen)),
                    Severity: "warning",
                })
            }
        }
    }
}
```

---

### 5.6 Data Storage

#### TimescaleDB Integration

```go
// internal/storage/timescale.go
package storage

import (
    "database/sql"
    "fmt"
    "time"

    _ "github.com/jackc/pgx/v5/stdlib"
)

type TimescaleDB struct {
    db *sql.DB
}

func NewTimescaleDB(connString string) (*TimescaleDB, error) {
    db, err := sql.Open("pgx", connString)
    if err != nil {
        return nil, err
    }
    
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    // Create hypertable
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS telemetry (
            device_id VARCHAR(100) NOT NULL,
            timestamp TIMESTAMPTZ NOT NULL,
            metric_type VARCHAR(50) NOT NULL,
            value DOUBLE PRECISION NOT NULL,
            metadata JSONB
        );
        
        SELECT create_hypertable('telemetry', 'timestamp', if_not_exists => TRUE);
        
        CREATE INDEX IF NOT EXISTS idx_telemetry_device_time 
        ON telemetry (device_id, timestamp DESC);
    `)
    if err != nil {
        return nil, err
    }
    
    return &TimescaleDB{db: db}, nil
}

func (t *TimescaleDB) InsertTelemetry(telemetry Telemetry) error {
    for key, value := range telemetry.Values {
        query := `
            INSERT INTO telemetry (device_id, timestamp, metric_type, value, metadata)
            VALUES ($1, $2, $3, $4, $5)
        `
        metadata, _ := json.Marshal(map[string]interface{}{
            "type": telemetry.Type,
        })
        
        _, err := t.db.Exec(query, telemetry.DeviceID, telemetry.Timestamp, 
            key, value, metadata)
        if err != nil {
            return err
        }
    }
    return nil
}

func (t *TimescaleDB) QueryTelemetry(deviceID string, from, to time.Time, interval string) ([]TelemetryPoint, error) {
    query := `
        SELECT 
            time_bucket($1, timestamp) AS bucket,
            metric_type,
            AVG(value) AS avg_value,
            MAX(value) AS max_value,
            MIN(value) AS min_value
        FROM telemetry
        WHERE device_id = $2 
            AND timestamp BETWEEN $3 AND $4
        GROUP BY bucket, metric_type
        ORDER BY bucket DESC
    `
    
    rows, err := t.db.Query(query, interval, deviceID, from, to)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var results []TelemetryPoint
    for rows.Next() {
        var point TelemetryPoint
        err := rows.Scan(&point.Timestamp, &point.MetricType, 
            &point.Avg, &point.Max, &point.Min)
        if err != nil {
            return nil, err
        }
        results = append(results, point)
    }
    
    return results, nil
}
```

#### Redis Cache Integration

```go
// internal/storage/redis.go
package storage

import (
    "context"
    "encoding/json"
    "time"

    "github.com/redis/go-redis/v9"
)

type Redis struct {
    client *redis.Client
    ctx    context.Context
}

func NewRedis(addr, password string, db int) (*Redis, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })
    
    ctx := context.Background()
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, err
    }
    
    return &Redis{
        client: client,
        ctx:    ctx,
    }, nil
}

func (r *Redis) SetTelemetry(deviceID string, telemetry Telemetry) error {
    key := fmt.Sprintf("telemetry:latest:%s", deviceID)
    data, err := json.Marshal(telemetry)
    if err != nil {
        return err
    }
    
    return r.client.Set(r.ctx, key, data, 1*time.Hour).Err()
}

func (r *Redis) GetLatestTelemetry(deviceID string) (*Telemetry, error) {
    key := fmt.Sprintf("telemetry:latest:%s", deviceID)
    data, err := r.client.Get(r.ctx, key).Result()
    if err == redis.Nil {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    var telemetry Telemetry
    if err := json.Unmarshal([]byte(data), &telemetry); err != nil {
        return nil, err
    }
    
    return &telemetry, nil
}

func (r *Redis) SetDeviceConfig(deviceID string, config map[string]interface{}) error {
    key := fmt.Sprintf("device:config:%s", deviceID)
    data, err := json.Marshal(config)
    if err != nil {
        return err
    }
    return r.client.Set(r.ctx, key, data, 24*time.Hour).Err()
}

func (r *Redis) GetDeviceConfig(deviceID string) (map[string]interface{}, error) {
    key := fmt.Sprintf("device:config:%s", deviceID)
    data, err := r.client.Get(r.ctx, key).Result()
    if err == redis.Nil {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    var config map[string]interface{}
    if err := json.Unmarshal([]byte(data), &config); err != nil {
        return nil, err
    }
    return config, nil
}
```

---

### 5.7 Alarm Analysis

#### Alert Engine

```go
// internal/alert/engine.go
package alert

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "sync"
    "time"

    "go-iot-platform/internal/mqtt"
    "go-iot-platform/internal/storage"
)

type Alert struct {
    ID        string                 `json:"id"`
    DeviceID  string                 `json:"device_id"`
    Type      string                 `json:"type"`
    Message   string                 `json:"message"`
    Severity  string                 `json:"severity"` // info, warning, critical
    Status    string                 `json:"status"`   // pending, acknowledged, resolved
    CreatedAt time.Time              `json:"created_at"`
    UpdatedAt time.Time              `json:"updated_at"`
    Metadata  map[string]interface{} `json:"metadata"`
}

type AlertRule struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    DeviceID    string                 `json:"device_id,omitempty"` // ถ้าเป็น empty ให้ใช้กับทุก device
    Metric      string                 `json:"metric"`
    Condition   string                 `json:"condition"` // gt, lt, eq, neq, between
    Threshold   interface{}            `json:"threshold"`
    Threshold2  interface{}            `json:"threshold2,omitempty"`
    Severity    string                 `json:"severity"`
    Message     string                 `json:"message"`
    Enabled     bool                   `json:"enabled"`
    Cooldown    int                    `json:"cooldown"` // seconds
    Actions     []string               `json:"actions"` // email, webhook, sms
    CreatedAt   time.Time              `json:"created_at"`
}

type AlertEngine struct {
    rules      []AlertRule
    alerts     []Alert
    mu         sync.RWMutex
    mqttClient *mqtt.Manager
    storage    *storage.TimescaleDB
    ctx        context.Context
    cancel     context.CancelFunc
    cooldown   map[string]time.Time
}

func NewAlertEngine(mqttClient *mqtt.Manager, db *storage.TimescaleDB) *AlertEngine {
    ctx, cancel := context.WithCancel(context.Background())
    return &AlertEngine{
        rules:      make([]AlertRule, 0),
        alerts:     make([]Alert, 0),
        mqttClient: mqttClient,
        storage:    db,
        ctx:        ctx,
        cancel:     cancel,
        cooldown:   make(map[string]time.Time),
    }
}

func (ae *AlertEngine) Start() error {
    // Subscribe to telemetry for alert processing
    if err := ae.mqttClient.Subscribe("devices/+/telemetry", ae.processTelemetry); err != nil {
        return err
    }
    return nil
}

func (ae *AlertEngine) processTelemetry(topic string, payload []byte) {
    var telemetry storage.Telemetry
    if err := json.Unmarshal(payload, &telemetry); err != nil {
        return
    }
    
    ae.mu.RLock()
    defer ae.mu.RUnlock()
    
    for _, rule := range ae.rules {
        if !rule.Enabled {
            continue
        }
        
        if rule.DeviceID != "" && rule.DeviceID != telemetry.DeviceID {
            continue
        }
        
        value, ok := telemetry.Values[rule.Metric]
        if !ok {
            continue
        }
        
        if ae.checkRule(rule, value) {
            ae.triggerAlert(rule, telemetry)
        }
    }
}

func (ae *AlertEngine) checkRule(rule AlertRule, value interface{}) bool {
    val := toFloat64(value)
    
    switch rule.Condition {
    case "gt":
        threshold := toFloat64(rule.Threshold)
        return val > threshold
    case "gte":
        threshold := toFloat64(rule.Threshold)
        return val >= threshold
    case "lt":
        threshold := toFloat64(rule.Threshold)
        return val < threshold
    case "lte":
        threshold := toFloat64(rule.Threshold)
        return val <= threshold
    case "eq":
        threshold := toFloat64(rule.Threshold)
        return val == threshold
    case "neq":
        threshold := toFloat64(rule.Threshold)
        return val != threshold
    case "between":
        t1 := toFloat64(rule.Threshold)
        t2 := toFloat64(rule.Threshold2)
        return val >= t1 && val <= t2
    default:
        return false
    }
}

func (ae *AlertEngine) triggerAlert(rule AlertRule, telemetry storage.Telemetry) {
    // Check cooldown
    key := fmt.Sprintf("%s:%s", rule.ID, telemetry.DeviceID)
    if lastTime, ok := ae.cooldown[key]; ok {
        if time.Since(lastTime).Seconds() < float64(rule.Cooldown) {
            return
        }
    }
    
    alert := Alert{
        ID:        uuid.New().String(),
        DeviceID:  telemetry.DeviceID,
        Type:      rule.Metric,
        Message:   rule.Message,
        Severity:  rule.Severity,
        Status:    "pending",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Metadata: map[string]interface{}{
            "rule_id":   rule.ID,
            "rule_name": rule.Name,
            "value":     telemetry.Values[rule.Metric],
        },
    }
    
    ae.mu.Lock()
    ae.alerts = append(ae.alerts, alert)
    ae.cooldown[key] = time.Now()
    ae.mu.Unlock()
    
    log.Printf("[ALERT] %s: %s", alert.Severity, alert.Message)
    
    // Send alert via MQTT
    payload, _ := json.Marshal(alert)
    ae.mqttClient.Publish("alerts/new", 1, false, payload)
    
    // Execute actions
    ae.executeActions(alert, rule.Actions)
}

func (ae *AlertEngine) executeActions(alert Alert, actions []string) {
    for _, action := range actions {
        switch action {
        case "email":
            ae.sendEmailAlert(alert)
        case "webhook":
            ae.sendWebhookAlert(alert)
        case "sms":
            ae.sendSMSAlert(alert)
        }
    }
}

func (ae *AlertEngine) AddRule(rule AlertRule) {
    ae.mu.Lock()
    defer ae.mu.Unlock()
    rule.CreatedAt = time.Now()
    ae.rules = append(ae.rules, rule)
}

func (ae *AlertEngine) GetAlerts(deviceID string, limit int) ([]Alert, error) {
    // Query from database
    return ae.storage.GetAlerts(deviceID, limit)
}

func (ae *AlertEngine) AckAlert(id string) error {
    ae.mu.Lock()
    defer ae.mu.Unlock()
    
    for i, alert := range ae.alerts {
        if alert.ID == id {
            ae.alerts[i].Status = "acknowledged"
            ae.alerts[i].UpdatedAt = time.Now()
            return nil
        }
    }
    return fmt.Errorf("alert not found")
}

// sendEmailAlert, sendWebhookAlert, sendSMSAlert implementations
```

---

### 5.8 Data Visualization

#### Real-time Dashboard

```go
// internal/api/handler.go
package api

import (
    "encoding/json"
    "net/http"
    "time"

    "go-iot-platform/internal/device"
    "go-iot-platform/internal/storage"
)

type DashboardHandler struct {
    deviceManager *device.DeviceManager
    storage       *storage.TimescaleDB
}

func (h *DashboardHandler) GetDashboardData(w http.ResponseWriter, r *http.Request) {
    // Get device status summary
    devices, _ := h.deviceManager.ListDevices()
    
    online := 0
    offline := 0
    for _, d := range devices {
        if d.Status == "online" {
            online++
        } else {
            offline++
        }
    }
    
    // Get latest telemetry
    latestData := make([]device.Telemetry, 0)
    for _, d := range devices {
        if data, err := h.storage.GetLatestTelemetry(d.ID); err == nil && data != nil {
            latestData = append(latestData, *data)
        }
    }
    
    response := map[string]interface{}{
        "status": map[string]interface{}{
            "total_devices":   len(devices),
            "online_devices":  online,
            "offline_devices": offline,
        },
        "latest_telemetry": latestData,
        "timestamp":        time.Now(),
    }
    
    json.NewEncoder(w).Encode(response)
}

func (h *DashboardHandler) GetDeviceData(w http.ResponseWriter, r *http.Request) {
    deviceID := r.PathValue("id")
    from := r.URL.Query().Get("from")
    to := r.URL.Query().Get("to")
    interval := r.URL.Query().Get("interval")
    
    if interval == "" {
        interval = "1h"
    }
    
    // Parse time
    var fromTime, toTime time.Time
    if from != "" {
        fromTime, _ = time.Parse(time.RFC3339, from)
    } else {
        fromTime = time.Now().Add(-24 * time.Hour)
    }
    
    if to != "" {
        toTime, _ = time.Parse(time.RFC3339, to)
    } else {
        toTime = time.Now()
    }
    
    // Query data
    data, err := h.storage.QueryTelemetry(deviceID, fromTime, toTime, interval)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(data)
}
```

#### WebSocket for Real-time Updates

```go
// internal/api/websocket.go
package api

import (
    "encoding/json"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type WSClient struct {
    Conn   *websocket.Conn
    Send   chan []byte
    Device string
}

type WSHub struct {
    clients    map[*WSClient]bool
    broadcast  chan []byte
    register   chan *WSClient
    unregister chan *WSClient
    mu         sync.Mutex
}

var hub = &WSHub{
    clients:    make(map[*WSClient]bool),
    broadcast:  make(chan []byte),
    register:   make(chan *WSClient),
    unregister: make(chan *WSClient),
}

func init() {
    go hub.run()
}

func (h *WSHub) run() {
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
            h.mu.Lock()
            for client := range h.clients {
                select {
                case client.Send <- message:
                default:
                    close(client.Send)
                    delete(h.clients, client)
                }
            }
            h.mu.Unlock()
        }
    }
}

func (h *WSHub) Broadcast(data interface{}) {
    message, _ := json.Marshal(data)
    h.broadcast <- message
}

func (h *WSHub) BroadcastToDevice(deviceID string, data interface{}) {
    message, _ := json.Marshal(data)
    
    h.mu.Lock()
    defer h.mu.Unlock()
    
    for client := range h.clients {
        if client.Device == deviceID || client.Device == "" {
            select {
            case client.Send <- message:
            default:
                close(client.Send)
                delete(h.clients, client)
            }
        }
    }
}

func ServeWS(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, "WebSocket upgrade failed", http.StatusBadRequest)
        return
    }
    
    client := &WSClient{
        Conn:   conn,
        Send:   make(chan []byte, 256),
        Device: r.URL.Query().Get("device"),
    }
    
    hub.register <- client
    
    // Write goroutine
    go func() {
        defer conn.Close()
        for message := range client.Send {
            if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
                return
            }
        }
    }()
    
    // Read goroutine
    go func() {
        defer func() {
            hub.unregister <- client
            conn.Close()
        }()
        
        for {
            _, message, err := conn.ReadMessage()
            if err != nil {
                return
            }
            
            // Handle client messages
            handleClientMessage(client, message)
        }
    }()
}

func handleClientMessage(client *WSClient, message []byte) {
    // Process client message
    var msg map[string]interface{}
    if err := json.Unmarshal(message, &msg); err != nil {
        return
    }
    
    // Echo back or process
    hub.BroadcastToDevice(client.Device, map[string]interface{}{
        "type":    "response",
        "message": "Message received",
        "data":    msg,
    })
}
```

---

## 📝 สรุปบทที่ 5: MQTT และ IoT

### สิ่งที่เราเรียนรู้

| หัวข้อ | รายละเอียด |
|--------|------------|
| **MQTT Protocol** | Publish/Subscribe, QoS, Broker |
| **MQTT Client** | การเชื่อมต่อ, Publish, Subscribe |
| **Device Manager** | การจัดการอุปกรณ์, สถานะ, Config |
| **Data Storage** | TimescaleDB, Redis Cache |
| **Alert Engine** | Rules, Threshold, Notifications |
| **WebSocket** | Real-time Updates, Dashboard |

### โครงสร้าง Complete IoT Platform

```yaml
# docker-compose.yaml
services:
  mosquitto:
    image: eclipse-mosquitto
    ports:
      - "1883:1883"
      - "9001:9001"
  
  timescaledb:
    image: timescale/timescaledb:latest-pg14
    environment:
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
  
  redis:
    image: redis:alpine
    ports:
      - "6379:6379"
  
  go-processor:
    build: .
    environment:
      MQTT_BROKER: mosquitto:1883
      DB_CONN: postgres://postgres:password@timescaledb:5432/postgres
      REDIS_ADDR: redis:6379
    depends_on:
      - mosquitto
      - timescaledb
      - redis
  
  api:
    build: ./api
    ports:
      - "8080:8080"
    depends_on:
      - go-processor
```

### การทดสอบ

```bash
# 1. เริ่มระบบ
docker-compose up -d

# 2. จำลองการส่งข้อมูลจาก Sensor
mosquitto_pub -t "devices/sensor-001/telemetry" \
  -m '{"device_id":"sensor-001","timestamp":"2024-01-01T00:00:00Z","values":{"temperature":28.5,"humidity":65}}'

# 3. ดูข้อมูลผ่าน API
curl http://localhost:8080/api/dashboard

# 4. WebSocket Dashboard
# เปิด Browser และเชื่อมต่อ ws://localhost:8080/ws
```

---
## Example

# 🚀 โค้ดเต็มรูปแบบสำหรับโมดูล IoT (Full Implementation)


## 📁 โครงสร้างไฟล์ที่จะสร้าง/เพิ่ม

```
icmongolang/
├── pkg/
│   └── mqtt/                          # [สร้างใหม่] MQTT Client Wrapper
│       ├── client.go
│       ├── manager.go
│       └── options.go
├── internal/modules/iot/
│   ├── models/
│   │   ├── alert_rule.go              # [สร้างใหม่] ตารางเก็บกฎแจ้งเตือน
│   │   └── telemetry.go               # [สร้างใหม่] Struct สำหรับรับข้อมูลจาก MQTT
│   ├── repository/
│   │   └── telemetry_repo.go          # [สร้างใหม่] Query ข้อมูลแบบ Temporal
│   ├── iothelper/
│   │   ├── alarm.go                   # [มีอยู่แล้ว - ขยาย]
│   │   └── alert_engine.go            # [สร้างใหม่] Rule Engine หลัก
│   ├── usecase/
│   │   ├── usecase.go                 # [มีอยู่แล้ว - เพิ่มฟังก์ชัน]
│   │   ├── device_manager.go          # [สร้างใหม่] จัดการสถานะ Device แบบ Real-time
│   │   └── telemetry_usecase.go       # [สร้างใหม่] ประมวลผล Telemetry
│   ├── delivery/http/
│   │   ├── handler.go                 # [มีอยู่แล้ว - เพิ่ม endpoint]
│   │   ├── websocket_handler.go       # [สร้างใหม่] เชื่อมต่อ WebSocket กับ Hub
│   │   └── routes.go                  # [มีอยู่แล้ว - เพิ่ม route]
└── cmd/
    └── mqtt-worker/                   # [สร้างใหม่] Microservice รับ MQTT
        └── main.go
```

---

## 1. MQTT Client Wrapper (`pkg/mqtt`)

### `pkg/mqtt/options.go`
```go
package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"time"
)

type Config struct {
	Broker    string `mapstructure:"broker"`     // tcp://localhost:1883
	ClientID  string `mapstructure:"client_id"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	QOS       byte   `mapstructure:"qos"`        // 0, 1, 2
	KeepAlive int    `mapstructure:"keep_alive"` // seconds

	// TLS
	TLSEnabled      bool   `mapstructure:"tls_enabled"`
	CAFile          string `mapstructure:"ca_file"`
	ClientCertFile  string `mapstructure:"client_cert_file"`
	ClientKeyFile   string `mapstructure:"client_key_file"`
	InsecureSkipVerify bool `mapstructure:"insecure_skip_verify"`

	// Retry
	ConnectRetryInterval time.Duration `mapstructure:"connect_retry_interval"`
	AutoReconnect        bool          `mapstructure:"auto_reconnect"`
	CleanSession         bool          `mapstructure:"clean_session"`
}

func DefaultConfig() Config {
	return Config{
		QOS:                  1,
		KeepAlive:            60,
		AutoReconnect:        true,
		CleanSession:         true,
		ConnectRetryInterval: 5 * time.Second,
		InsecureSkipVerify:   false,
	}
}

func (c *Config) TLSConfig() (*tls.Config, error) {
	if !c.TLSEnabled {
		return nil, nil
	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.InsecureSkipVerify,
		MinVersion:         tls.VersionTLS12,
	}
	if c.CAFile != "" {
		caCert, err := os.ReadFile(c.CAFile)
		if err != nil {
			return nil, err
		}
		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caPool
	}
	if c.ClientCertFile != "" && c.ClientKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(c.ClientCertFile, c.ClientKeyFile)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	return tlsConfig, nil
}
```

### `pkg/mqtt/manager.go`
```go
package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MessageHandler รับ topic และ payload
type MessageHandler func(topic string, payload []byte)

// Manager จัดการ MQTT Client และ Routing
type Manager struct {
	client   mqtt.Client
	opts     *mqtt.ClientOptions
	config   Config
	handlers map[string]MessageHandler // topic pattern -> handler
	mu       sync.RWMutex
	isConnected bool
	onConnectHooks []func()
}

// NewManager สร้าง MQTT Manager
func NewManager(cfg Config) (*Manager, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.Broker)
	opts.SetClientID(cfg.ClientID)
	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
	}
	if cfg.Password != "" {
		opts.SetPassword(cfg.Password)
	}
	opts.SetKeepAlive(time.Duration(cfg.KeepAlive) * time.Second)
	opts.SetCleanSession(cfg.CleanSession)
	opts.SetAutoReconnect(cfg.AutoReconnect)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(cfg.ConnectRetryInterval)
	opts.SetPingTimeout(10 * time.Second)

	if cfg.TLSEnabled {
		tlsCfg, err := cfg.TLSConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS config: %w", err)
		}
		opts.SetTLSConfig(tlsCfg)
	}

	mgr := &Manager{
		config:   cfg,
		opts:     opts,
		handlers: make(map[string]MessageHandler),
	}
	opts.SetOnConnectHandler(mgr.onConnect)
	opts.SetConnectionLostHandler(mgr.onConnectionLost)

	mgr.client = mqtt.NewClient(opts)
	return mgr, nil
}

// Connect เชื่อมต่อ MQTT Broker
func (m *Manager) Connect() error {
	token := m.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	m.isConnected = true
	log.Printf("[MQTT] Connected to %s", m.config.Broker)
	return nil
}

// Disconnect ยกเลิกการเชื่อมต่อ
func (m *Manager) Disconnect(quiesce uint) {
	m.client.Disconnect(quiesce)
	m.isConnected = false
	log.Println("[MQTT] Disconnected")
}

// IsConnected เช็คสถานะการเชื่อมต่อ
func (m *Manager) IsConnected() bool {
	return m.client.IsConnected()
}

// onConnect ถูกเรียกเมื่อเชื่อมต่อสำเร็จ (รวมถึง reconnect)
func (m *Manager) onConnect(client mqtt.Client) {
	m.isConnected = true
	log.Println("[MQTT] OnConnect: Connected (or reconnected)")
	m.mu.RLock()
	handlers := make(map[string]MessageHandler)
	for k, v := range m.handlers {
		handlers[k] = v
	}
	m.mu.RUnlock()

	// Resubscribe ทุก topic
	for topic, handler := range handlers {
		m.doSubscribe(topic, handler)
	}

	for _, hook := range m.onConnectHooks {
		hook()
	}
}

// onConnectionLost ถูกเรียกเมื่อ connection หลุด
func (m *Manager) onConnectionLost(client mqtt.Client, err error) {
	m.isConnected = false
	log.Printf("[MQTT] Connection Lost: %v", err)
}

// RegisterOnConnectHook เพิ่มฟังก์ชันที่จะเรียกเมื่อเชื่อมต่อ
func (m *Manager) RegisterOnConnectHook(hook func()) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onConnectHooks = append(m.onConnectHooks, hook)
}

// Publish ส่งข้อความ
func (m *Manager) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	var data []byte
	switch v := payload.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		var err error
		data, err = json.Marshal(v)
		if err != nil {
			return err
		}
	}
	token := m.client.Publish(topic, qos, retained, data)
	token.Wait()
	return token.Error()
}

// Subscribe ลงทะเบียนรับข้อความ
func (m *Manager) Subscribe(topic string, handler MessageHandler) error {
	m.mu.Lock()
	m.handlers[topic] = handler
	m.mu.Unlock()

	if m.client.IsConnected() {
		return m.doSubscribe(topic, handler)
	}
	return nil
}

func (m *Manager) doSubscribe(topic string, handler MessageHandler) error {
	token := m.client.Subscribe(topic, m.config.QOS, func(client mqtt.Client, msg mqtt.Message) {
		// เรียก handler ผ่าน goroutine เพื่อไม่ให้ block MQTT callback
		go handler(msg.Topic(), msg.Payload())
	})
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}
	log.Printf("[MQTT] Subscribed to: %s", topic)
	return nil
}

// Unsubscribe ยกเลิกการรับข้อความ
func (m *Manager) Unsubscribe(topic string) error {
	m.mu.Lock()
	delete(m.handlers, topic)
	m.mu.Unlock()

	token := m.client.Unsubscribe(topic)
	token.Wait()
	return token.Error()
}
```

---

## 2. Models (Entities)

### `internal/modules/iot/models/telemetry.go` (ใหม่)
```go
package models

import (
	"time"
)

// TelemetryPayload ตรงกับ JSON ที่ Device ส่งมา
type TelemetryPayload struct {
	DeviceID  string                 `json:"device_id"`
	Timestamp time.Time              `json:"timestamp"`
	Type      string                 `json:"type"` // เช่น "sensor", "status"
	Values    map[string]interface{} `json:"values"`
}

// TelemetryQueryResult ใช้สำหรับ Aggregated Query
type TelemetryAggregate struct {
	Bucket     time.Time `json:"bucket"`
	MetricType string    `json:"metric_type"`
	AvgValue   float64   `json:"avg_value"`
	MaxValue   float64   `json:"max_value"`
	MinValue   float64   `json:"min_value"`
	Count      int64     `json:"count"`
}
```

### `internal/modules/iot/models/alert_rule.go` (ใหม่)
```go
package models

import (
	"time"
	"gorm.io/gorm"
)

type AlertRule struct {
	ID          string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	DeviceID    string         `gorm:"type:varchar(100);index" json:"device_id"` // 空 = ทุกอุปกรณ์
	Metric      string         `gorm:"type:varchar(50);not null" json:"metric"`
	Condition   string         `gorm:"type:varchar(20);not null" json:"condition"` // gt, lt, eq, gte, lte, between, neq
	Threshold   float64        `json:"threshold"`
	Threshold2  float64        `json:"threshold2,omitempty"`
	Severity    string         `gorm:"type:varchar(20);default:'warning'" json:"severity"` // info, warning, critical
	Message     string         `gorm:"type:text" json:"message"`
	Cooldown    int            `gorm:"default:300" json:"cooldown"` // seconds
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	Actions     []string       `gorm:"type:jsonb" json:"actions"` // ["email", "webhook"]
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AlertRule) TableName() string {
	return "iot_alert_rules"
}
```

---

## 3. Repository

### `internal/modules/iot/repository/telemetry_repo.go` (ใหม่)
```go
package repository

import (
	"context"
	"time"

	"icmongolang/internal/modules/iot/models"
	"gorm.io/gorm"
)

type TelemetryRepository interface {
	Insert(ctx context.Context, data *models.IotData) error
	GetLatest(ctx context.Context, deviceID string) (*models.IotData, error)
	QueryAggregated(ctx context.Context, deviceID string, from, to time.Time, interval string) ([]models.TelemetryAggregate, error)
	QueryRaw(ctx context.Context, deviceID string, from, to time.Time, limit int) ([]models.IotData, error)
}

type telemetryRepo struct {
	db *gorm.DB
}

func NewTelemetryRepository(db *gorm.DB) TelemetryRepository {
	return &telemetryRepo{db: db}
}

// Insert ใช้กับตาราง iot_data ที่มีอยู่ (สมมติว่ามี field: device_id, timestamp, data_json)
func (r *telemetryRepo) Insert(ctx context.Context, data *models.IotData) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *telemetryRepo) GetLatest(ctx context.Context, deviceID string) (*models.IotData, error) {
	var data models.IotData
	err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Order("timestamp DESC").
		First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// QueryAggregated ใช้ PostgreSQL date_trunc หรือ TimescaleDB time_bucket
// ถ้าใช้ TimescaleDB: SELECT time_bucket($1, timestamp) AS bucket, ...
// ถ้าใช้ PostgreSQL ปกติ: SELECT date_trunc($1, timestamp) AS bucket, ...
func (r *telemetryRepo) QueryAggregated(ctx context.Context, deviceID string, from, to time.Time, interval string) ([]models.TelemetryAggregate, error) {
	// สมมติว่า iot_data มี field: device_id, timestamp, metric_name (varchar), value (double precision)
	// ถ้า data อยู่ใน JSONB ให้ใช้ jsonb_extract_path_text
	// เพื่อความง่าย สมมติว่ามีคอลัมน์ metric_name และ value
	query := `
		SELECT 
			date_trunc(?, timestamp) AS bucket,
			metric_name AS metric_type,
			AVG(value) AS avg_value,
			MAX(value) AS max_value,
			MIN(value) AS min_value,
			COUNT(*) AS count
		FROM iot_data
		WHERE device_id = ?
		  AND timestamp BETWEEN ? AND ?
		  AND metric_name IS NOT NULL
		GROUP BY bucket, metric_name
		ORDER BY bucket DESC
	`

	var results []models.TelemetryAggregate
	err := r.db.WithContext(ctx).Raw(query, interval, deviceID, from, to).Scan(&results).Error
	return results, err
}

func (r *telemetryRepo) QueryRaw(ctx context.Context, deviceID string, from, to time.Time, limit int) ([]models.IotData, error) {
	var data []models.IotData
	err := r.db.WithContext(ctx).
		Where("device_id = ? AND timestamp BETWEEN ? AND ?", deviceID, from, to).
		Order("timestamp DESC").
		Limit(limit).
		Find(&data).Error
	return data, err
}
```

---

## 4. IoT Helper (Alert Engine)

### `internal/modules/iot/iothelper/alert_engine.go` (ใหม่)
```go
package iothelper

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"icmongolang/internal/modules/iot/models"
)

// AlertEngine เป็น Rule Engine สำหรับประเมินเงื่อนไขและสร้าง Alert
type AlertEngine struct {
	rules      []models.AlertRule
	cooldown   map[string]time.Time // key: rule_id:device_id
	mu         sync.RWMutex
	onAlertFn  func(alert *models.AlarmLog) // callback เมื่อเกิด alert
}

func NewAlertEngine() *AlertEngine {
	return &AlertEngine{
		rules:    make([]models.AlertRule, 0),
		cooldown: make(map[string]time.Time),
	}
}

// LoadRules โหลดกฎจาก DB (เรียกตอนเริ่มต้น)
func (e *AlertEngine) LoadRules(rules []models.AlertRule) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules = rules
	log.Printf("[AlertEngine] Loaded %d rules", len(rules))
}

// SetAlertCallback ตั้ง callback เมื่อเกิด alert
func (e *AlertEngine) SetAlertCallback(fn func(alert *models.AlarmLog)) {
	e.onAlertFn = fn
}

// Evaluate ตรวจสอบ Telemetry กับทุก Rule
func (e *AlertEngine) Evaluate(deviceID string, values map[string]interface{}, timestamp time.Time) {
	e.mu.RLock()
	rules := make([]models.AlertRule, len(e.rules))
	copy(rules, e.rules)
	e.mu.RUnlock()

	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}
		if rule.DeviceID != "" && rule.DeviceID != deviceID {
			continue
		}

		val, ok := values[rule.Metric]
		if !ok {
			continue
		}

		// ตรวจสอบ Cooldown
		cooldownKey := fmt.Sprintf("%s:%s", rule.ID, deviceID)
		e.mu.RLock()
		lastTime, exists := e.cooldown[cooldownKey]
		e.mu.RUnlock()
		if exists {
			if time.Since(lastTime).Seconds() < float64(rule.Cooldown) {
				continue
			}
		}

		// ตรวจสอบเงื่อนไข
		triggered := e.checkCondition(rule, val)
		if triggered {
			// สร้าง Alert
			alert := &models.AlarmLog{
				// สมมติว่า AlarmLog มี fields: DeviceID, Type, Message, Severity, Status, CreatedAt
				DeviceID:  deviceID,
				Type:      rule.Metric,
				Message:   rule.Message,
				Severity:  rule.Severity,
				Status:    "pending",
				CreatedAt: timestamp,
				// Metadata: map[string]interface{}{"rule_id": rule.ID, "value": val},
			}
			// ใช้ json.RawMessage หรือ map ในกรณีที่ AlarmLog มี field Metadata
			// สมมติว่า AlarmLog มี Metadata *string หรือ json.RawMessage
			meta := map[string]interface{}{
				"rule_id": rule.ID,
				"value":   val,
			}
			metaBytes, _ := json.Marshal(meta)
			alert.Metadata = string(metaBytes) // ถ้าเป็น string

			// อัปเดต cooldown
			e.mu.Lock()
			e.cooldown[cooldownKey] = time.Now()
			e.mu.Unlock()

			log.Printf("[AlertEngine] TRIGGERED: %s | %s | %s", rule.Severity, deviceID, rule.Message)

			// เรียก callback
			if e.onAlertFn != nil {
				go e.onAlertFn(alert)
			}
		}
	}
}

// checkCondition เปรียบเทียบค่า
func (e *AlertEngine) checkCondition(rule models.AlertRule, val interface{}) bool {
	floatVal, ok := toFloat64(val)
	if !ok {
		return false
	}
	threshold := rule.Threshold
	threshold2 := rule.Threshold2

	switch rule.Condition {
	case "gt":
		return floatVal > threshold
	case "gte":
		return floatVal >= threshold
	case "lt":
		return floatVal < threshold
	case "lte":
		return floatVal <= threshold
	case "eq":
		return floatVal == threshold
	case "neq":
		return floatVal != threshold
	case "between":
		return floatVal >= threshold && floatVal <= threshold2
	default:
		return false
	}
}

// toFloat64 แปลง interface{} เป็น float64
func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case int32:
		return float64(val), true
	case string:
		var f float64
		_, err := fmt.Sscanf(val, "%f", &f)
		return f, err == nil
	default:
		return 0, false
	}
}
```

---

## 5. Usecase

### `internal/modules/iot/usecase/device_manager.go` (ใหม่)
```go
package usecase

import (
	"context"
	"log"
	"sync"
	"time"

	"icmongolang/internal/modules/iot/models"
	"icmongolang/internal/modules/iot/repository"
)

// DeviceManager จัดการสถานะอุปกรณ์แบบ Real-time และ Health Check
type DeviceManager struct {
	repo       repository.DeviceRepo // สมมติว่ามี interface นี้
	statusRepo repository.DeviceStatusRepo
	devices    map[string]*models.Device // cache
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	onStatusChange func(deviceID, oldStatus, newStatus string)
}

func NewDeviceManager(deviceRepo repository.DeviceRepo, statusRepo repository.DeviceStatusRepo) *DeviceManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &DeviceManager{
		repo:       deviceRepo,
		statusRepo: statusRepo,
		devices:    make(map[string]*models.Device),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// SetStatusChangeCallback เรียกเมื่อสถานะเปลี่ยน
func (dm *DeviceManager) SetStatusChangeCallback(fn func(deviceID, old, new string)) {
	dm.onStatusChange = fn
}

// UpdateStatus อัปเดตสถานะ (เรียกจาก MQTT)
func (dm *DeviceManager) UpdateStatus(deviceID, status string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	device, exists := dm.devices[deviceID]
	if !exists {
		// โหลดจาก DB ก่อน
		d, err := dm.repo.GetByID(dm.ctx, deviceID)
		if err != nil {
			return err
		}
		device = d
		dm.devices[deviceID] = device
	}

	oldStatus := device.Status
	if oldStatus == status {
		return nil
	}

	device.Status = status
	device.LastSeen = time.Now()

	// อัปเดต DB
	if err := dm.statusRepo.UpdateStatus(dm.ctx, deviceID, status); err != nil {
		log.Printf("[DeviceManager] Failed to update status in DB: %v", err)
	}

	if dm.onStatusChange != nil {
		go dm.onStatusChange(deviceID, oldStatus, status)
	}
	return nil
}

// HealthCheck ตรวจสอบอุปกรณ์ที่หายไป
func (dm *DeviceManager) HealthCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-dm.ctx.Done():
			return
		case <-ticker.C:
			dm.checkOfflineDevices()
		}
	}
}

func (dm *DeviceManager) checkOfflineDevices() {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	timeout := 3 * time.Minute
	now := time.Now()

	for id, device := range dm.devices {
		if device.Status == "offline" {
			continue
		}
		if now.Sub(device.LastSeen) > timeout {
			log.Printf("[DeviceManager] Device %s is offline (last seen: %v)", id, device.LastSeen)
			// อัปเดตเป็น offline
			go dm.UpdateStatus(id, "offline")
		}
	}
}

// Stop หยุด HealthCheck
func (dm *DeviceManager) Stop() {
	dm.cancel()
}
```

### `internal/modules/iot/usecase/telemetry_usecase.go` (ใหม่)
```go
package usecase

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"icmongolang/internal/modules/iot/iothelper"
	"icmongolang/internal/modules/iot/models"
	"icmongolang/internal/modules/iot/repository"
)

// TelemetryUsecase ประมวลผลข้อมูลที่ได้รับจาก MQTT
type TelemetryUsecase struct {
	telemetryRepo repository.TelemetryRepository
	deviceManager *DeviceManager
	alertEngine   *iothelper.AlertEngine
	wsBroadcaster func(message interface{}) // ฟังก์ชัน Broadcast ไปยัง WebSocket
}

func NewTelemetryUsecase(
	telemetryRepo repository.TelemetryRepository,
	deviceManager *DeviceManager,
	alertEngine *iothelper.AlertEngine,
	wsBroadcaster func(message interface{}),
) *TelemetryUsecase {
	return &TelemetryUsecase{
		telemetryRepo: telemetryRepo,
		deviceManager: deviceManager,
		alertEngine:   alertEngine,
		wsBroadcaster: wsBroadcaster,
	}
}

// ProcessTelemetry รับ payload JSON และประมวลผล
func (u *TelemetryUsecase) ProcessTelemetry(payload []byte) {
	var telemetry models.TelemetryPayload
	if err := json.Unmarshal(payload, &telemetry); err != nil {
		log.Printf("[TelemetryUsecase] Invalid JSON: %v", err)
		return
	}

	if telemetry.DeviceID == "" {
		log.Println("[TelemetryUsecase] Missing device_id")
		return
	}
	if telemetry.Timestamp.IsZero() {
		telemetry.Timestamp = time.Now()
	}

	// 1. อัปเดตสถานะ Device
	u.deviceManager.UpdateStatus(telemetry.DeviceID, "online")

	// 2. บันทึก Telemetry ลง DB
	// แปลงเป็น models.IotData (สมมติว่า IotData มี DataJSON หรือ fields)
	iotData := &models.IotData{
		DeviceID:  telemetry.DeviceID,
		Timestamp: telemetry.Timestamp,
		// สมมติว่ามี field DataJSON (string) หรือ map
	}
	// ถ้า IotData มี DataJSON (string)
	if dataBytes, err := json.Marshal(telemetry.Values); err == nil {
		iotData.DataJSON = string(dataBytes)
	}
	if err := u.telemetryRepo.Insert(context.Background(), iotData); err != nil {
		log.Printf("[TelemetryUsecase] Failed to save: %v", err)
	}

	// 3. ตรวจสอบ Alert Rules
	u.alertEngine.Evaluate(telemetry.DeviceID, telemetry.Values, telemetry.Timestamp)

	// 4. Broadcast ไปยัง WebSocket (ถ้ามี)
	if u.wsBroadcaster != nil {
		u.wsBroadcaster(map[string]interface{}{
			"type":      "telemetry",
			"device_id": telemetry.DeviceID,
			"timestamp": telemetry.Timestamp,
			"values":    telemetry.Values,
		})
	}
}
```

---

## 6. WebSocket Handler (เชื่อมต่อกับ `pkg/websocket` ที่มีอยู่)

สมมติว่า `pkg/websocket` ที่มีอยู่มี `Hub` พร้อม `Broadcast` channel หรือ method

### `internal/modules/iot/delivery/http/websocket_handler.go` (ใหม่)
```go
package http

import (
	"encoding/json"
	"log"
	"net/http"

	"icmongolang/pkg/websocket" // สมมติว่ามี Hub อยู่แล้ว
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 在生产环境ควรจำกัด origin
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WSHandler struct {
	hub *websocket.Hub
}

func NewWSHandler(hub *websocket.Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

// ServeWS จัดการ WebSocket Connection
func (h *WSHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] Upgrade failed: %v", err)
		return
	}

	// สร้าง Client ตามโครงสร้างของ pkg/websocket
	// สมมติว่า pkg/websocket มี Client struct และ Hub มี Register/Unregister channel
	client := &websocket.Client{
		Hub:  h.hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	client.Hub.Register <- client

	// Goroutine สำหรับเขียน
	go func() {
		defer conn.Close()
		for {
			message, ok := <-client.Send
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}()

	// Goroutine สำหรับอ่าน (รับ message จาก client)
	go func() {
		defer func() {
			client.Hub.Unregister <- client
			conn.Close()
		}()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			// echo หรือประมวลผล message จาก client
			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err == nil {
				log.Printf("[WS] Received: %v", msg)
				// อาจจะตอบกลับ
			}
		}
	}()
}

// BroadcastTelemetry เป็นฟังก์ชันที่ใช้ส่งข้อมูลไปยังทุก Client
// ฟังก์ชันนี้จะถูกส่งเข้าไปใน TelemetryUsecase
func BroadcastTelemetry(hub *websocket.Hub, message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("[WS] Marshal error: %v", err)
		return
	}
	// สมมติว่า Hub มี Broadcast channel
	hub.Broadcast <- data
}
```

---

## 7. MQTT Worker (Microservice)

### `cmd/mqtt-worker/main.go` (ใหม่)
```go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"icmongolang/internal/modules/iot/iothelper"
	"icmongolang/internal/modules/iot/models"
	"icmongolang/internal/modules/iot/repository"
	"icmongolang/internal/modules/iot/usecase"
	"icmongolang/internal/modules/iot/delivery/http"
	"icmongolang/pkg/mqtt"
	"icmongolang/pkg/websocket"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.Println("[MQTT-Worker] Starting...")

	// 1. Load Config (ใช้ environment variables หรือ config file)
	broker := os.Getenv("MQTT_BROKER")
	if broker == "" {
		broker = "tcp://localhost:1883"
	}
	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		dbDSN = "host=localhost user=postgres password=password dbname=icmongo port=5432 sslmode=disable"
	}

	// 2. Connect Database
	db, err := gorm.Open(postgres.Open(dbDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}
	log.Println("[MQTT-Worker] DB connected")

	// 3. Auto-migrate (ถ้ายังไม่มีตาราง)
	db.AutoMigrate(&models.AlertRule{})

	// 4. Init Repositories
	deviceRepo := repository.NewDeviceRepo(db)          // สมมติว่ามี
	statusRepo := repository.NewDeviceStatusRepo(db)    // สมมติว่ามี
	telemetryRepo := repository.NewTelemetryRepository(db)
	alertRuleRepo := repository.NewAlertRuleRepo(db)    // สมมติว่ามี หรือสร้าง

	// 5. Init Alert Engine และโหลด Rules
	alertEngine := iothelper.NewAlertEngine()
	rules, _ := alertRuleRepo.GetAllEnabled(context.Background())
	alertEngine.LoadRules(rules)

	// 6. Init Device Manager
	deviceManager := usecase.NewDeviceManager(deviceRepo, statusRepo)
	go deviceManager.HealthCheck()

	// 7. Init WebSocket Hub (จาก pkg/websocket)
	hub := websocket.NewHub()
	go hub.Run()

	// 8. สร้าง Broadcaster function
	broadcaster := func(message interface{}) {
		// ใช้ method Broadcast ของ hub
		hub.Broadcast(message)
	}

	// 9. Init Telemetry Usecase
	telemetryUC := usecase.NewTelemetryUsecase(
		telemetryRepo,
		deviceManager,
		alertEngine,
		broadcaster,
	)

	// 10. ตั้งค่า Alert Callback (เมื่อเกิด Alert ให้ broadcast และบันทึก)
	alertEngine.SetAlertCallback(func(alert *models.AlarmLog) {
		// Broadcast Alert ไปยัง WebSocket
		broadcaster(map[string]interface{}{
			"type":     "alert",
			"device":   alert.DeviceID,
			"message":  alert.Message,
			"severity": alert.Severity,
			"time":     alert.CreatedAt,
		})
		// บันทึก Alert ลง DB (ใช้ AlarmLogRepo)
		if err := db.Create(alert).Error; err != nil {
			log.Printf("[MQTT-Worker] Failed to save alert: %v", err)
		}
	})

	// 11. Init MQTT Manager
	mqttCfg := mqtt.DefaultConfig()
	mqttCfg.Broker = broker
	mqttCfg.ClientID = "mqtt-worker-" + time.Now().Format("20060102150405")
	mqttCfg.QOS = 1

	mqttManager, err := mqtt.NewManager(mqttCfg)
	if err != nil {
		log.Fatalf("Failed to create MQTT manager: %v", err)
	}
	if err := mqttManager.Connect(); err != nil {
		log.Fatalf("Failed to connect MQTT: %v", err)
	}
	defer mqttManager.Disconnect(250)

	// 12. Subscribe Topics
	// Telemetry
	if err := mqttManager.Subscribe("devices/+/telemetry", func(topic string, payload []byte) {
		telemetryUC.ProcessTelemetry(payload)
	}); err != nil {
		log.Fatalf("Subscribe telemetry failed: %v", err)
	}

	// Status
	if err := mqttManager.Subscribe("devices/+/status", func(topic string, payload []byte) {
		// parse status
		var statusMsg struct {
			Status string `json:"status"`
		}
		if err := json.Unmarshal(payload, &statusMsg); err == nil {
			deviceID := extractDeviceIDFromTopic(topic)
			if deviceID != "" {
				deviceManager.UpdateStatus(deviceID, statusMsg.Status)
			}
		}
	}); err != nil {
		log.Fatalf("Subscribe status failed: %v", err)
	}

	// Event
	if err := mqttManager.Subscribe("devices/+/event", func(topic string, payload []byte) {
		log.Printf("[MQTT-Worker] Event received: %s", string(payload))
	}); err != nil {
		log.Fatalf("Subscribe event failed: %v", err)
	}

	log.Println("[MQTT-Worker] All subscriptions registered. Running...")

	// 13. Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("[MQTT-Worker] Shutting down...")
	deviceManager.Stop()
	log.Println("[MQTT-Worker] Stopped")
}

// helper: แยก deviceID จาก topic "devices/{id}/status"
func extractDeviceIDFromTopic(topic string) string {
	parts := strings.Split(topic, "/")
	if len(parts) >= 3 {
		return parts[1]
	}
	return ""
}
```

**อย่าลืมเพิ่ม import `"strings"` และ `"encoding/json"` ใน main.go**

---

## 8. Routes (เพิ่มใน `delivery/http/routes.go`)

```go
// เพิ่มใน routes.go ของ module iot
func SetupRoutes(router *gin.Engine, handler *Handler, wsHandler *WSHandler) {
	// ... routes เดิม ...

	// WebSocket
	router.GET("/ws", wsHandler.ServeWS)

	// Dashboard API (ใหม่)
	router.GET("/api/dashboard/summary", handler.GetDashboardSummary)
	router.GET("/api/devices/:id/telemetry", handler.GetDeviceTelemetry)
}
```

### เพิ่ม Handler สำหรับ Dashboard
```go
// ใน handler.go
func (h *Handler) GetDashboardSummary(c *gin.Context) {
	// เรียกใช้ Usecase หรือ Repository โดยตรง
	// ...
	c.JSON(200, gin.H{
		"total_devices":  100,
		"online_devices": 85,
		"offline_devices": 15,
		"latest_alerts":  []interface{}{},
	})
}

func (h *Handler) GetDeviceTelemetry(c *gin.Context) {
	deviceID := c.Param("id")
	from := c.Query("from")
	to := c.Query("to")
	// เรียก telemetryRepo.QueryAggregated
	// ...
	c.JSON(200, gin.H{"data": []interface{}{}})
}
```

---

## 9. Migration SQL (เพิ่มเติม)

### `migrations/20250620_alert_rules.sql`
```sql
CREATE TABLE IF NOT EXISTS iot_alert_rules (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    device_id VARCHAR(100),
    metric VARCHAR(50) NOT NULL,
    condition VARCHAR(20) NOT NULL,
    threshold DOUBLE PRECISION DEFAULT 0,
    threshold2 DOUBLE PRECISION DEFAULT 0,
    severity VARCHAR(20) DEFAULT 'warning',
    message TEXT,
    cooldown INT DEFAULT 300,
    enabled BOOLEAN DEFAULT TRUE,
    actions JSONB DEFAULT '[]'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_alert_rules_device ON iot_alert_rules(device_id);
CREATE INDEX idx_alert_rules_enabled ON iot_alert_rules(enabled);
```

---

## 10. Docker Compose (เพิ่ม MQTT-Worker)

```yaml
# docker-compose.yml (เพิ่ม service ใหม่)
services:
  mqtt-worker:
    build:
      context: .
      dockerfile: Dockerfile.worker  # หรือใช้ Dockerfile เดิมกับ target
    environment:
      MQTT_BROKER: "tcp://mosquitto:1883"
      DB_DSN: "host=timescaledb user=postgres password=password dbname=icmongo port=5432 sslmode=disable"
    depends_on:
      - mosquitto
      - timescaledb
    restart: unless-stopped
```

---

## สรุปสิ่งที่คุณได้เพิ่มเข้าไป

| ไฟล์ | หน้าที่ |
|------|--------|
| `pkg/mqtt/*` | MQTT Client Wrapper พร้อม reconnect, resubscribe |
| `models/alert_rule.go` | Entity สำหรับเก็บกฎ Alert |
| `models/telemetry.go` | Struct สำหรับรับ JSON จาก Device |
| `repository/telemetry_repo.go` | Query แบบ Aggregate (avg, max, min) |
| `iothelper/alert_engine.go` | Rule Engine ตรวจสอบเงื่อนไขและ Cooldown |
| `usecase/device_manager.go` | จัดการสถานะ Device + Health Check |
| `usecase/telemetry_usecase.go` | ประมวลผล Telemetry ครบวงจร |
| `delivery/http/websocket_handler.go` | เชื่อม WebSocket กับ Hub |
| `cmd/mqtt-worker/main.go` | Microservice รับ MQTT และประมวลผล |
| `migrations/20250620_alert_rules.sql` | สร้างตาราง Alert Rules |

---

## วิธีทดสอบ

1. **สร้างตารางใหม่**:
   ```bash
   go run cmd/migrate/main.go  # หรือใช้ migration tool
   ```

2. **รัน MQTT Worker**:
   ```bash
   go run cmd/mqtt-worker/main.go
   ```

3. **จำลองส่งข้อมูล**:
   ```bash
   mosquitto_pub -t "devices/sensor-001/telemetry" -m '{"device_id":"sensor-001","timestamp":"2026-09-07T10:00:00Z","values":{"temperature":35.5,"humidity":80}}'
   ```

4. **ดู WebSocket**:
   เปิด Browser แล้วเชื่อมต่อ `ws://localhost:8080/ws` (หรือพอร์ตตามที่ API ของคุณใช้)

5. **เพิ่ม Alert Rule** (ผ่าน SQL หรือ API):
   ```sql
   INSERT INTO iot_alert_rules (id, name, device_id, metric, condition, threshold, severity, message, cooldown)
   VALUES ('rule-1', 'Temp High', '', 'temperature', 'gt', 30, 'critical', 'Temperature exceeds 30°C', 60);
   ```
   เมื่อส่ง temperature 35.5 จะเกิด Alert ทันที

---
 

## บทถัดไป: REST API Design

ในบทที่ 6 เราจะเรียนรู้เกี่ยวกับ:
- การออกแบบ REST API อย่างมืออาชีพ
- Routing และ Handler ที่เป็นระเบียบ
- Middleware และ Error Handling
- การทำ API Documentation
- Security และ Authentication
 