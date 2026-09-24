# 📘 เล่มที่ 7: โครงการ IoT ขนาดใหญ่ (Go IoT Platform)

---

## 📖 บทนำ

ในยุคที่ทุกสิ่งเชื่อมต่อถึงกัน (Internet of Things) การพัฒนาแพลตฟอร์มที่สามารถรองรับอุปกรณ์จำนวนมากและจัดการข้อมูลแบบ Real-time เป็นความท้าทายที่สำคัญ Go เป็นภาษาที่เหมาะสำหรับงานประเภทนี้ด้วยประสิทธิภาพสูงและการจัดการ Concurrency ที่ยอดเยี่ยม

เล่มนี้จะพาคุณสร้าง **Go IoT Platform** — แพลตฟอร์ม IoT ระดับ Production ที่พร้อมใช้งานจริง โดยอ้างอิงจากโปรเจกต์ Open Source [go-iot-platform](https://github.com/iot-ecology/go-iot-platform)

### ภาพรวมของ Go IoT Platform

```mermaid
flowchart TB
    subgraph Devices["อุปกรณ์ IoT"]
        D1[ESP8266/ESP32]
        D2[Raspberry Pi]
        D3[Arduino]
        D4[Industrial Sensors]
    end
    
    subgraph Gateway["MQTT Gateway"]
        MQTT[MQTT Broker<br/>EMQX/Mosquitto]
    end
    
    subgraph Backend["Go Backend Services"]
        DM[Device Manager<br/>go-iot]
        MQ[Message Queue<br/>RabbitMQ]
        DS[Data Storage<br/>TimescaleDB]
        AE[Alert Engine]
        API[REST API]
    end
    
    subgraph Frontend["Frontend"]
        DASH[Admin Dashboard<br/>Vue.js]
        MOB[Mobile App]
    end
    
    D1 --> MQTT
    D2 --> MQTT
    D3 --> MQTT
    D4 --> MQTT
    
    MQTT --> DM
    MQTT --> MQ
    
    DM --> DS
    MQ --> DS
    DM --> AE
    AE --> DASH
    
    API --> DS
    DASH --> API
    MOB --> API
```

---

## สารบัญ เล่มที่ 7

### บทที่ 1: ภาพรวมของ Go IoT Platform
- 1.1 Go IoT Platform คืออะไร?
- 1.2 คุณสมบัติหลัก

### บทที่ 2: สถาปัตยกรรมระบบ
- 2.1 go-iot — MQTT Client Management Service
- 2.2 go-iot-mq — Rabbit Message Queue Service
- 2.3 iot-go-project — Management Backend Service
- 2.4 ant-vue — Admin Dashboard

### บทที่ 3: MQTT Client Management
- 3.1 การ maintain stable connections
- 3.2 การจัดการ MQTT Client จำนวนมาก
- 3.3 การเพิ่มและลบ MQTT Client

### บทที่ 4: Data Storage
- 4.1 การออกแบบฐานข้อมูลสำหรับ IoT
- 4.2 การเก็บข้อมูลที่รายงานจาก MQTT
- 4.3 การ query ข้อมูลย้อนหลัง

### บทที่ 5: Alarm Analysis
- 5.1 Real-time Monitoring
- 5.2 การตั้งค่า Alarm
- 5.3 การแจ้งเตือนเมื่อเกิดเหตุการณ์

### บทที่ 6: Data Visualization
- 6.1 การแสดงข้อมูลแบบ Real-time
- 6.2 Dashboard Design
- 6.3 การสร้างกราฟและ图表

### บทที่ 7: Offline Computing
- 7.1 การประมวลผลข้อมูลย้อนหลัง
- 7.2 Batch Processing
- 7.3 Data Analytics

### บทที่ 8: การติดตั้งและ Deployment
- 8.1 Deployment Guide
- 8.2 การตั้งค่า Environment
- 8.3 การรันใน Production

### บทที่ 9: การมีส่วนร่วมและเอกสาร
- 9.1 การรายงาน Issues
- 9.2 การส่ง Pull Requests
- 9.3 การปรับปรุง Documentation

---

## บทที่ 1: ภาพรวมของ Go IoT Platform

---

### 1.1 Go IoT Platform คืออะไร?

**Go IoT Platform** คือแพลตฟอร์ม Open Source สำหรับการจัดการอุปกรณ์ IoT ที่เขียนด้วยภาษา Go มีเป้าหมายเพื่อให้:

- **รองรับอุปกรณ์จำนวนมาก** — สามารถจัดการ MQTT Client ได้หลายพันตัว
- **ทำงานแบบ Real-time** — รับและประมวลผลข้อมูลทันที
- **มีความเสถียรสูง** — ออกแบบมาสำหรับ Production
- **ขยายได้ง่าย** — สถาปัตยกรรมแบบ Microservices

#### โครงสร้างหลัก

```
go-iot-platform/
├── go-iot/              # MQTT Client Management Service
│   ├── cmd/
│   ├── internal/
│   └── pkg/
├── go-iot-mq/           # RabbitMQ Message Queue Service
│   ├── cmd/
│   └── internal/
├── iot-go-project/      # Management Backend Service
│   ├── cmd/
│   └── internal/
└── ant-vue/             # Admin Dashboard (Vue.js)
    ├── src/
    └── public/
```

### 1.2 คุณสมบัติหลัก

| คุณสมบัติ | คำอธิบาย |
|-----------|----------|
| **MQTT Client Management** | จัดการการเชื่อมต่อ MQTT, Reconnect, Authentication |
| **Data Storage** | เก็บข้อมูลใน TimescaleDB (Time-series Database) |
| **Alarm Analysis** | วิเคราะห์และแจ้งเตือนแบบ Real-time |
| **Data Visualization** | แสดงข้อมูลผ่าน Dashboard |
| **Offline Computing** | ประมวลผลข้อมูลย้อนหลัง |
| **REST API** | API สำหรับการจัดการอุปกรณ์และข้อมูล |
| **Multi-tenant** | รองรับหลายองค์กร/โปรเจกต์ |
| **Scalable** | ขยายแนวนอน (Horizontal Scaling) |

---

## บทที่ 2: สถาปัตยกรรมระบบ

---

### 2.1 go-iot — MQTT Client Management Service

**go-iot** เป็นบริการหลักที่รับผิดชอบการจัดการ MQTT Client ทั้งหมด

#### โครงสร้าง

```go
// go-iot/internal/mqtt/manager.go
package mqtt

import (
    "context"
    "crypto/tls"
    "fmt"
    "log"
    "sync"
    "time"

    mqtt "github.com/eclipse/paho.mqtt.golang"
)

type ClientManager struct {
    clients      map[string]*ManagedClient
    mu           sync.RWMutex
    config       Config
    messageCh    chan Message
    ctx          context.Context
    cancel       context.CancelFunc
    wg           sync.WaitGroup
}

type ManagedClient struct {
    ID         string
    Client     mqtt.Client
    Options    *mqtt.ClientOptions
    Topics     []string
    Status     ClientStatus
    LastPing   time.Time
    CreatedAt  time.Time
    UpdatedAt  time.Time
    Metadata   map[string]interface{}
    mu         sync.RWMutex
}

type ClientStatus string

const (
    StatusConnected   ClientStatus = "connected"
    StatusConnecting  ClientStatus = "connecting"
    StatusDisconnected ClientStatus = "disconnected"
    StatusError       ClientStatus = "error"
)

type Config struct {
    BrokerURL      string
    ConnectTimeout time.Duration
    KeepAlive      int
    MaxReconnects  int
    ReconnectDelay time.Duration
    TLSConfig      *tls.Config
}

func NewClientManager(cfg Config) *ClientManager {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &ClientManager{
        clients:   make(map[string]*ManagedClient),
        config:    cfg,
        messageCh: make(chan Message, 1000),
        ctx:       ctx,
        cancel:    cancel,
    }
}

func (m *ClientManager) Start() error {
    log.Println("Starting MQTT Client Manager")
    
    // Start message processor
    m.wg.Add(1)
    go m.processMessages()
    
    // Start health checker
    m.wg.Add(1)
    go m.healthChecker()
    
    return nil
}

func (m *ClientManager) Stop() {
    m.cancel()
    
    // Disconnect all clients
    m.mu.RLock()
    for _, client := range m.clients {
        if client.Client != nil && client.Client.IsConnected() {
            client.Client.Disconnect(250)
        }
    }
    m.mu.RUnlock()
    
    m.wg.Wait()
    log.Println("MQTT Client Manager stopped")
}

func (m *ClientManager) RegisterClient(id string, topics []string, metadata map[string]interface{}) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    if _, exists := m.clients[id]; exists {
        return fmt.Errorf("client %s already registered", id)
    }
    
    // Create client options
    opts := mqtt.NewClientOptions()
    opts.AddBroker(m.config.BrokerURL)
    opts.SetClientID(id)
    opts.SetKeepAlive(time.Duration(m.config.KeepAlive) * time.Second)
    opts.SetAutoReconnect(true)
    opts.SetConnectRetry(true)
    opts.SetConnectRetryInterval(m.config.ReconnectDelay)
    
    if m.config.TLSConfig != nil {
        opts.SetTLSConfig(m.config.TLSConfig)
    }
    
    // Set callbacks
    opts.SetOnConnectHandler(func(client mqtt.Client) {
        m.handleConnect(id)
        m.subscribeTopics(id, topics)
    })
    
    opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
        m.handleDisconnect(id, err)
    })
    
    // Create client
    client := mqtt.NewClient(opts)
    
    managed := &ManagedClient{
        ID:        id,
        Client:    client,
        Options:   opts,
        Topics:    topics,
        Status:    StatusConnecting,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Metadata:  metadata,
    }
    
    m.clients[id] = managed
    
    // Connect
    token := client.Connect()
    token.Wait()
    
    if token.Error() != nil {
        managed.Status = StatusError
        return token.Error()
    }
    
    log.Printf("Client %s registered and connected", id)
    return nil
}

func (m *ClientManager) subscribeTopics(id string, topics []string) {
    m.mu.RLock()
    managed, exists := m.clients[id]
    m.mu.RUnlock()
    
    if !exists || managed.Client == nil {
        return
    }
    
    for _, topic := range topics {
        token := managed.Client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
            m.onMessage(id, msg.Topic(), msg.Payload())
        })
        token.Wait()
        if token.Error() != nil {
            log.Printf("Failed to subscribe to %s: %v", topic, token.Error())
        }
    }
}

func (m *ClientManager) UnregisterClient(id string) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    managed, exists := m.clients[id]
    if !exists {
        return fmt.Errorf("client %s not found", id)
    }
    
    if managed.Client != nil && managed.Client.IsConnected() {
        managed.Client.Disconnect(250)
    }
    
    delete(m.clients, id)
    log.Printf("Client %s unregistered", id)
    
    return nil
}

func (m *ClientManager) Publish(id string, topic string, qos byte, retained bool, payload []byte) error {
    m.mu.RLock()
    managed, exists := m.clients[id]
    m.mu.RUnlock()
    
    if !exists {
        return fmt.Errorf("client %s not found", id)
    }
    
    if managed.Client == nil || !managed.Client.IsConnected() {
        return fmt.Errorf("client %s is not connected", id)
    }
    
    token := managed.Client.Publish(topic, qos, retained, payload)
    token.Wait()
    return token.Error()
}

func (m *ClientManager) onMessage(id, topic string, payload []byte) {
    // Send to message channel
    select {
    case m.messageCh <- Message{
        ClientID: id,
        Topic:    topic,
        Payload:  payload,
        Timestamp: time.Now(),
    }:
    default:
        log.Printf("Message channel full, dropping message from %s", id)
    }
}

func (m *ClientManager) processMessages() {
    defer m.wg.Done()
    
    for {
        select {
        case <-m.ctx.Done():
            return
        case msg := <-m.messageCh:
            m.handleMessage(msg)
        }
    }
}

func (m *ClientManager) handleMessage(msg Message) {
    log.Printf("Received from %s: topic=%s, size=%d bytes", 
        msg.ClientID, msg.Topic, len(msg.Payload))
    
    // Store in database
    // Trigger alerts
    // Process data
}

func (m *ClientManager) healthChecker() {
    defer m.wg.Done()
    
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-m.ctx.Done():
            return
        case <-ticker.C:
            m.checkHealth()
        }
    }
}

func (m *ClientManager) checkHealth() {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    now := time.Now()
    for id, managed := range m.clients {
        if managed.Status == StatusConnected {
            // Check if client is still connected
            if managed.Client == nil || !managed.Client.IsConnected() {
                managed.Status = StatusDisconnected
                log.Printf("Client %s disconnected", id)
                continue
            }
            
            // Check last ping
            if now.Sub(managed.LastPing) > 2*time.Duration(m.config.KeepAlive)*time.Second {
                log.Printf("Client %s seems dead (no ping)", id)
                managed.Status = StatusError
            }
        }
    }
}

func (m *ClientManager) handleConnect(id string) {
    m.mu.RLock()
    managed, exists := m.clients[id]
    m.mu.RUnlock()
    
    if exists {
        managed.mu.Lock()
        managed.Status = StatusConnected
        managed.LastPing = time.Now()
        managed.UpdatedAt = time.Now()
        managed.mu.Unlock()
        
        log.Printf("Client %s connected", id)
    }
}

func (m *ClientManager) handleDisconnect(id string, err error) {
    m.mu.RLock()
    managed, exists := m.clients[id]
    m.mu.RUnlock()
    
    if exists {
        managed.mu.Lock()
        managed.Status = StatusDisconnected
        managed.UpdatedAt = time.Now()
        managed.mu.Unlock()
        
        log.Printf("Client %s disconnected: %v", id, err)
    }
}
```

#### API Handler

```go
// go-iot/internal/api/handler.go
package api

import (
    "encoding/json"
    "net/http"
    "strconv"

    "go-iot/internal/mqtt"
    "go-iot/internal/storage"
)

type Handler struct {
    manager *mqtt.ClientManager
    storage *storage.Storage
}

func NewHandler(manager *mqtt.ClientManager, storage *storage.Storage) *Handler {
    return &Handler{
        manager: manager,
        storage: storage,
    }
}

func (h *Handler) RegisterClient(w http.ResponseWriter, r *http.Request) {
    var req struct {
        ID       string                 `json:"id"`
        Topics   []string               `json:"topics"`
        Metadata map[string]interface{} `json:"metadata"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    if req.ID == "" {
        respondError(w, http.StatusBadRequest, "Client ID is required")
        return
    }
    
    if err := h.manager.RegisterClient(req.ID, req.Topics, req.Metadata); err != nil {
        respondError(w, http.StatusConflict, err.Error())
        return
    }
    
    respondSuccess(w, http.StatusCreated, map[string]interface{}{
        "id": req.ID,
        "status": "registered",
    })
}

func (h *Handler) UnregisterClient(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    if id == "" {
        respondError(w, http.StatusBadRequest, "Client ID is required")
        return
    }
    
    if err := h.manager.UnregisterClient(id); err != nil {
        respondError(w, http.StatusNotFound, err.Error())
        return
    }
    
    respondSuccess(w, http.StatusOK, map[string]interface{}{
        "id": id,
        "status": "unregistered",
    })
}

func (h *Handler) GetClients(w http.ResponseWriter, r *http.Request) {
    // Implementation to get list of clients
    // ...
}

func (h *Handler) GetClientStatus(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    if id == "" {
        respondError(w, http.StatusBadRequest, "Client ID is required")
        return
    }
    
    // Get client status
    // ...
}
```

---

### 2.2 go-iot-mq — Rabbit Message Queue Service

**go-iot-mq** ทำหน้าที่เป็น Message Queue ระหว่างบริการต่างๆ

```go
// go-iot-mq/internal/rabbitmq/queue.go
package rabbitmq

import (
    "context"
    "encoding/json"
    "log"
    "time"

    amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
    ID        string      `json:"id"`
    Type      string      `json:"type"`
    Payload   interface{} `json:"payload"`
    Timestamp time.Time   `json:"timestamp"`
}

type Queue struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    config  Config
    ctx     context.Context
    cancel  context.CancelFunc
}

type Config struct {
    URL          string
    ExchangeName string
    QueueName    string
    RoutingKey   string
    Prefetch     int
}

func NewQueue(cfg Config) *Queue {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &Queue{
        config: cfg,
        ctx:    ctx,
        cancel: cancel,
    }
}

func (q *Queue) Connect() error {
    conn, err := amqp.Dial(q.config.URL)
    if err != nil {
        return err
    }
    q.conn = conn
    
    channel, err := conn.Channel()
    if err != nil {
        return err
    }
    q.channel = channel
    
    // Set QoS (prefetch count)
    if err := channel.Qos(q.config.Prefetch, 0, false); err != nil {
        return err
    }
    
    // Declare exchange
    if err := channel.ExchangeDeclare(
        q.config.ExchangeName,
        "topic",
        true,
        false,
        false,
        false,
        nil,
    ); err != nil {
        return err
    }
    
    // Declare queue
    _, err = channel.QueueDeclare(
        q.config.QueueName,
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return err
    }
    
    // Bind queue to exchange
    if err := channel.QueueBind(
        q.config.QueueName,
        q.config.RoutingKey,
        q.config.ExchangeName,
        false,
        nil,
    ); err != nil {
        return err
    }
    
    log.Println("Connected to RabbitMQ")
    return nil
}

func (q *Queue) Publish(message Message) error {
    body, err := json.Marshal(message)
    if err != nil {
        return err
    }
    
    return q.channel.Publish(
        q.config.ExchangeName,
        q.config.RoutingKey,
        false,
        false,
        amqp.Publishing{
            ContentType:  "application/json",
            Body:         body,
            Timestamp:    message.Timestamp,
            DeliveryMode: amqp.Persistent,
        },
    )
}

func (q *Queue) Consume(handler func(Message) error) error {
    deliveries, err := q.channel.Consume(
        q.config.QueueName,
        "",
        false,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return err
    }
    
    go func() {
        for delivery := range deliveries {
            var msg Message
            if err := json.Unmarshal(delivery.Body, &msg); err != nil {
                log.Printf("Failed to unmarshal message: %v", err)
                delivery.Ack(false)
                continue
            }
            
            if err := handler(msg); err != nil {
                log.Printf("Handler error: %v", err)
                delivery.Nack(false, true)
                continue
            }
            
            delivery.Ack(false)
        }
    }()
    
    return nil
}

func (q *Queue) Close() {
    q.cancel()
    if q.channel != nil {
        q.channel.Close()
    }
    if q.conn != nil {
        q.conn.Close()
    }
    log.Println("RabbitMQ connection closed")
}
```

---

### 2.3 iot-go-project — Management Backend Service

**iot-go-project** เป็น Backend Service สำหรับจัดการข้อมูลและ API

```go
// iot-go-project/internal/service/device.go
package service

import (
    "context"
    "time"

    "iot-go-project/internal/model"
    "iot-go-project/internal/repository"
)

type DeviceService struct {
    repo     *repository.DeviceRepository
    mqClient *mqtt.ClientManager
}

func NewDeviceService(repo *repository.DeviceRepository, mqClient *mqtt.ClientManager) *DeviceService {
    return &DeviceService{
        repo:     repo,
        mqClient: mqClient,
    }
}

func (s *DeviceService) CreateDevice(ctx context.Context, device *model.Device) error {
    device.CreatedAt = time.Now()
    device.UpdatedAt = time.Now()
    device.Status = "pending"
    
    if err := s.repo.Create(ctx, device); err != nil {
        return err
    }
    
    // Register with MQTT
    topics := []string{
        fmt.Sprintf("devices/%s/telemetry", device.ID),
        fmt.Sprintf("devices/%s/status", device.ID),
        fmt.Sprintf("devices/%s/command", device.ID),
    }
    
    return s.mqClient.RegisterClient(device.ID, topics, map[string]interface{}{
        "name": device.Name,
        "type": device.Type,
    })
}

func (s *DeviceService) GetDevice(ctx context.Context, id string) (*model.Device, error) {
    return s.repo.FindByID(ctx, id)
}

func (s *DeviceService) ListDevices(ctx context.Context, filters map[string]interface{}) ([]*model.Device, error) {
    return s.repo.List(ctx, filters)
}

func (s *DeviceService) UpdateDevice(ctx context.Context, id string, updates map[string]interface{}) (*model.Device, error) {
    device, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Apply updates
    if name, ok := updates["name"].(string); ok {
        device.Name = name
    }
    if status, ok := updates["status"].(string); ok {
        device.Status = status
    }
    if config, ok := updates["config"].(map[string]interface{}); ok {
        device.Config = config
    }
    
    device.UpdatedAt = time.Now()
    
    if err := s.repo.Update(ctx, device); err != nil {
        return nil, err
    }
    
    // Send updated config to device
    s.mqClient.Publish(device.ID, fmt.Sprintf("devices/%s/config", device.ID), 1, false, 
        []byte(device.Config.String()))
    
    return device, nil
}

func (s *DeviceService) SendCommand(ctx context.Context, deviceID string, command string, params map[string]interface{}) error {
    payload := map[string]interface{}{
        "command": command,
        "params":  params,
        "timestamp": time.Now(),
    }
    
    data, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    
    return s.mqClient.Publish(deviceID, fmt.Sprintf("devices/%s/command", deviceID), 1, false, data)
}

func (s *DeviceService) GetTelemetry(ctx context.Context, deviceID string, from, to time.Time) ([]*model.Telemetry, error) {
    return s.repo.GetTelemetry(ctx, deviceID, from, to)
}

func (s *DeviceService) GetLatestTelemetry(ctx context.Context, deviceID string) (*model.Telemetry, error) {
    return s.repo.GetLatestTelemetry(ctx, deviceID)
}
```

---

### 2.4 ant-vue — Admin Dashboard

**ant-vue** เป็น Frontend Dashboard ที่ใช้ Vue.js และ Ant Design

#### โครงสร้าง Vue Project

```
ant-vue/
├── src/
│   ├── api/
│   │   ├── device.js
│   │   ├── telemetry.js
│   │   └── alert.js
│   ├── components/
│   │   ├── DeviceList.vue
│   │   ├── DeviceDetail.vue
│   │   ├── Dashboard.vue
│   │   ├── TelemetryChart.vue
│   │   └── AlertList.vue
│   ├── views/
│   │   ├── Dashboard.vue
│   │   ├── Devices.vue
│   │   ├── Alerts.vue
│   │   └── Settings.vue
│   ├── router/
│   │   └── index.js
│   ├── store/
│   │   ├── device.js
│   │   └── alert.js
│   └── App.vue
└── package.json
```

#### Dashboard Component

```vue
<!-- ant-vue/src/components/Dashboard.vue -->
<template>
  <div class="dashboard">
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <a-icon type="wifi" />
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.totalDevices }}</div>
              <div class="stat-label">Total Devices</div>
            </div>
          </div>
        </a-card>
      </a-col>
      
      <a-col :span="6">
        <a-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon online">
              <a-icon type="check-circle" />
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.onlineDevices }}</div>
              <div class="stat-label">Online Devices</div>
            </div>
          </div>
        </a-card>
      </a-col>
      
      <a-col :span="6">
        <a-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon alert">
              <a-icon type="bell" />
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.activeAlerts }}</div>
              <div class="stat-label">Active Alerts</div>
            </div>
          </div>
        </a-card>
      </a-col>
      
      <a-col :span="6">
        <a-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon data">
              <a-icon type="database" />
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.totalMessages }}</div>
              <div class="stat-label">Total Messages</div>
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>
    
    <!-- Charts -->
    <a-row :gutter="16">
      <a-col :span="16">
        <a-card title="Device Telemetry">
          <telemetry-chart :data="chartData" />
        </a-card>
      </a-col>
      
      <a-col :span="8">
        <a-card title="Recent Alerts">
          <alert-list :alerts="recentAlerts" />
        </a-card>
      </a-col>
    </a-row>
    
    <!-- Device List -->
    <a-row>
      <a-col :span="24">
        <a-card title="Devices">
          <device-list :devices="devices" />
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script>
import { mapState, mapActions } from 'vuex'
import TelemetryChart from './TelemetryChart.vue'
import AlertList from './AlertList.vue'
import DeviceList from './DeviceList.vue'

export default {
  name: 'Dashboard',
  
  components: {
    TelemetryChart,
    AlertList,
    DeviceList,
  },
  
  data() {
    return {
      stats: {
        totalDevices: 0,
        onlineDevices: 0,
        activeAlerts: 0,
        totalMessages: 0,
      },
      chartData: [],
      recentAlerts: [],
      devices: [],
    }
  },
  
  mounted() {
    this.fetchData()
    this.startPolling()
  },
  
  methods: {
    ...mapActions(['fetchDevices', 'fetchAlerts']),
    
    async fetchData() {
      try {
        const [stats, devices, alerts, telemetry] = await Promise.all([
          this.$api.getStats(),
          this.$api.getDevices(),
          this.$api.getAlerts({ limit: 10 }),
          this.$api.getTelemetry({ limit: 100 }),
        ])
        
        this.stats = stats
        this.devices = devices
        this.recentAlerts = alerts
        this.chartData = telemetry
      } catch (error) {
        console.error('Failed to fetch data:', error)
      }
    },
    
    startPolling() {
      setInterval(() => {
        this.fetchData()
      }, 30000) // 30 seconds
    },
  },
}
</script>

<style scoped>
.dashboard {
  padding: 24px;
}

.stat-card {
  margin-bottom: 16px;
}

.stat-content {
  display: flex;
  align-items: center;
}

.stat-icon {
  font-size: 32px;
  margin-right: 16px;
  color: #1890ff;
}

.stat-icon.online {
  color: #52c41a;
}

.stat-icon.alert {
  color: #faad14;
}

.stat-icon.data {
  color: #722ed1;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
}

.stat-label {
  font-size: 14px;
  color: #8c8c8c;
}
</style>
```

---

## บทที่ 3: MQTT Client Management

---

### 3.1 การ maintain stable connections

#### การจัดการ Reconnection

```go
// go-iot/internal/mqtt/reconnect.go
package mqtt

import (
    "math"
    "time"
)

type ReconnectStrategy struct {
    MaxAttempts    int
    InitialDelay   time.Duration
    MaxDelay       time.Duration
    BackoffFactor  float64
}

func (s *ReconnectStrategy) GetDelay(attempt int) time.Duration {
    if attempt >= s.MaxAttempts {
        return s.MaxDelay
    }
    
    delay := float64(s.InitialDelay) * math.Pow(s.BackoffFactor, float64(attempt))
    if delay > float64(s.MaxDelay) {
        delay = float64(s.MaxDelay)
    }
    
    return time.Duration(delay)
}

func (m *ClientManager) reconnectClient(id string) error {
    m.mu.RLock()
    managed, exists := m.clients[id]
    m.mu.RUnlock()
    
    if !exists {
        return fmt.Errorf("client %s not found", id)
    }
    
    strategy := ReconnectStrategy{
        MaxAttempts:   10,
        InitialDelay:  1 * time.Second,
        MaxDelay:      60 * time.Second,
        BackoffFactor: 1.5,
    }
    
    for attempt := 0; attempt < strategy.MaxAttempts; attempt++ {
        delay := strategy.GetDelay(attempt)
        time.Sleep(delay)
        
        if managed.Client != nil && managed.Client.IsConnected() {
            return nil
        }
        
        log.Printf("Attempting to reconnect client %s (attempt %d/%d)", id, attempt+1, strategy.MaxAttempts)
        
        token := managed.Client.Connect()
        token.Wait()
        
        if token.Error() == nil {
            log.Printf("Client %s reconnected successfully", id)
            managed.Status = StatusConnected
            return nil
        }
    }
    
    return fmt.Errorf("failed to reconnect client %s after %d attempts", id, strategy.MaxAttempts)
}
```

---

### 3.2 การจัดการ MQTT Client จำนวนมาก

#### Connection Pool Management

```go
// go-iot/internal/mqtt/pool.go
package mqtt

import (
    "sync"
    "time"
)

type ConnectionPool struct {
    maxConnections int
    connections    []*ManagedClient
    mu             sync.RWMutex
    notifyChan     chan *ManagedClient
}

func NewConnectionPool(maxConn int) *ConnectionPool {
    return &ConnectionPool{
        maxConnections: maxConn,
        connections:    make([]*ManagedClient, 0, maxConn),
        notifyChan:     make(chan *ManagedClient, 100),
    }
}

func (p *ConnectionPool) Add(client *ManagedClient) error {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if len(p.connections) >= p.maxConnections {
        return fmt.Errorf("connection pool is full (max: %d)", p.maxConnections)
    }
    
    p.connections = append(p.connections, client)
    p.notifyChan <- client
    
    log.Printf("Added client %s to pool. Pool size: %d/%d", client.ID, len(p.connections), p.maxConnections)
    return nil
}

func (p *ConnectionPool) Remove(id string) error {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    for i, client := range p.connections {
        if client.ID == id {
            p.connections = append(p.connections[:i], p.connections[i+1:]...)
            log.Printf("Removed client %s from pool. Pool size: %d/%d", id, len(p.connections), p.maxConnections)
            return nil
        }
    }
    
    return fmt.Errorf("client %s not found in pool", id)
}

func (p *ConnectionPool) Get(id string) (*ManagedClient, bool) {
    p.mu.RLock()
    defer p.mu.RUnlock()
    
    for _, client := range p.connections {
        if client.ID == id {
            return client, true
        }
    }
    return nil, false
}

func (p *ConnectionPool) GetAll() []*ManagedClient {
    p.mu.RLock()
    defer p.mu.RUnlock()
    
    result := make([]*ManagedClient, len(p.connections))
    copy(result, p.connections)
    return result
}

func (p *ConnectionPool) Stats() PoolStats {
    p.mu.RLock()
    defer p.mu.RUnlock()
    
    stats := PoolStats{
        Total: len(p.connections),
        Max:   p.maxConnections,
    }
    
    for _, client := range p.connections {
        if client.Status == StatusConnected {
            stats.Connected++
        } else if client.Status == StatusDisconnected {
            stats.Disconnected++
        } else if client.Status == StatusError {
            stats.Error++
        }
    }
    
    return stats
}

type PoolStats struct {
    Total       int
    Connected   int
    Disconnected int
    Error       int
    Max         int
}
```

---

### 3.3 การเพิ่มและลบ MQTT Client

#### Client Lifecycle Management

```go
// go-iot/internal/mqtt/lifecycle.go
package mqtt

import (
    "context"
    "time"
)

type ClientLifecycle struct {
    manager *ClientManager
    pool    *ConnectionPool
    ctx     context.Context
    cancel  context.CancelFunc
}

func NewClientLifecycle(manager *ClientManager, pool *ConnectionPool) *ClientLifecycle {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &ClientLifecycle{
        manager: manager,
        pool:    pool,
        ctx:     ctx,
        cancel:  cancel,
    }
}

func (l *ClientLifecycle) Start() {
    go l.monitorClients()
}

func (l *ClientLifecycle) Stop() {
    l.cancel()
}

func (l *ClientLifecycle) monitorClients() {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-l.ctx.Done():
            return
        case <-ticker.C:
            l.cleanupDeadClients()
        }
    }
}

func (l *ClientLifecycle) cleanupDeadClients() {
    clients := l.pool.GetAll()
    now := time.Now()
    
    for _, client := range clients {
        // Check for clients that have been disconnected too long
        if client.Status == StatusDisconnected {
            if now.Sub(client.UpdatedAt) > 5*time.Minute {
                log.Printf("Cleaning up dead client: %s", client.ID)
                l.manager.UnregisterClient(client.ID)
                l.pool.Remove(client.ID)
            }
        }
    }
}

func (l *ClientLifecycle) CreateClient(config ClientConfig) (*ManagedClient, error) {
    // Validate
    if config.ID == "" {
        return nil, fmt.Errorf("client ID is required")
    }
    
    if _, exists := l.pool.Get(config.ID); exists {
        return nil, fmt.Errorf("client %s already exists", config.ID)
    }
    
    // Create client
    client, err := l.manager.createClient(config)
    if err != nil {
        return nil, err
    }
    
    // Add to pool
    if err := l.pool.Add(client); err != nil {
        return nil, err
    }
    
    return client, nil
}

func (l *ClientLifecycle) DeleteClient(id string) error {
    // Remove from pool
    if err := l.pool.Remove(id); err != nil {
        return err
    }
    
    // Unregister from manager
    return l.manager.UnregisterClient(id)
}

type ClientConfig struct {
    ID        string
    Topics    []string
    Metadata  map[string]interface{}
    QoS       byte
    KeepAlive int
}
```

---

## บทที่ 4: Data Storage

---

### 4.1 การออกแบบฐานข้อมูลสำหรับ IoT

#### Schema Design

```sql
-- TimescaleDB Schema
-- devices table
CREATE TABLE IF NOT EXISTS devices (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    location JSONB,
    config JSONB,
    metadata JSONB,
    last_seen TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- telemetry table (hypertable)
CREATE TABLE IF NOT EXISTS telemetry (
    device_id VARCHAR(100) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    metric VARCHAR(50) NOT NULL,
    value DOUBLE PRECISION,
    string_value TEXT,
    metadata JSONB
);

SELECT create_hypertable('telemetry', 'timestamp', if_not_exists => TRUE);

-- Indexes
CREATE INDEX idx_telemetry_device_time ON telemetry (device_id, timestamp DESC);
CREATE INDEX idx_telemetry_metric ON telemetry (metric);

-- alerts table
CREATE TABLE IF NOT EXISTS alerts (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    device_id VARCHAR(100) NOT NULL,
    rule_id VARCHAR(100) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    message TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMPTZ
);

-- rules table
CREATE TABLE IF NOT EXISTS alert_rules (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    device_id VARCHAR(100),
    metric VARCHAR(50) NOT NULL,
    condition VARCHAR(20) NOT NULL,
    threshold JSONB NOT NULL,
    severity VARCHAR(20) NOT NULL,
    message TEXT NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    cooldown_seconds INTEGER DEFAULT 300,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

---

### 4.2 การเก็บข้อมูลที่รายงานจาก MQTT

#### Storage Implementation

```go
// iot-go-project/internal/storage/timescale.go
package storage

import (
    "context"
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
    
    return &TimescaleDB{db: db}, nil
}

func (t *TimescaleDB) InsertTelemetry(ctx context.Context, deviceID string, timestamp time.Time, 
    metrics map[string]interface{}, metadata map[string]interface{}) error {
    
    tx, err := t.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO telemetry (device_id, timestamp, metric, value, string_value, metadata)
        VALUES ($1, $2, $3, $4, $5, $6)
    `)
    if err != nil {
        return err
    }
    defer stmt.Close()
    
    metaJSON, _ := json.Marshal(metadata)
    
    for metric, value := range metrics {
        var floatVal sql.NullFloat64
        var stringVal sql.NullString
        
        switch v := value.(type) {
        case float64:
            floatVal.Valid = true
            floatVal.Float64 = v
        case int:
            floatVal.Valid = true
            floatVal.Float64 = float64(v)
        case string:
            stringVal.Valid = true
            stringVal.String = v
        case bool:
            stringVal.Valid = true
            stringVal.String = fmt.Sprintf("%t", v)
        default:
            stringVal.Valid = true
            stringVal.String = fmt.Sprintf("%v", v)
        }
        
        if _, err := stmt.ExecContext(ctx, deviceID, timestamp, metric, 
            floatVal, stringVal, metaJSON); err != nil {
            return err
        }
    }
    
    // Update device last_seen
    if _, err := tx.ExecContext(ctx, 
        "UPDATE devices SET last_seen = $1, updated_at = $1 WHERE id = $2",
        timestamp, deviceID); err != nil {
        return err
    }
    
    return tx.Commit()
}

func (t *TimescaleDB) QueryTelemetry(ctx context.Context, deviceID string, 
    metric string, from, to time.Time, interval string) ([]TelemetryPoint, error) {
    
    query := `
        SELECT 
            time_bucket($1, timestamp) AS bucket,
            AVG(value) AS avg_value,
            MAX(value) AS max_value,
            MIN(value) AS min_value,
            COUNT(*) AS sample_count
        FROM telemetry
        WHERE device_id = $2 
            AND metric = $3
            AND timestamp BETWEEN $4 AND $5
        GROUP BY bucket
        ORDER BY bucket ASC
    `
    
    rows, err := t.db.QueryContext(ctx, query, interval, deviceID, metric, from, to)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var points []TelemetryPoint
    for rows.Next() {
        var point TelemetryPoint
        err := rows.Scan(&point.Timestamp, &point.Avg, &point.Max, &point.Min, &point.Count)
        if err != nil {
            return nil, err
        }
        points = append(points, point)
    }
    
    return points, nil
}

type TelemetryPoint struct {
    Timestamp time.Time
    Avg       float64
    Max       float64
    Min       float64
    Count     int64
}
```

---

### 4.3 การ query ข้อมูลย้อนหลัง

#### Query Service

```go
// iot-go-project/internal/service/telemetry.go
package service

import (
    "context"
    "fmt"
    "time"

    "iot-go-project/internal/storage"
)

type TelemetryService struct {
    db *storage.TimescaleDB
}

func NewTelemetryService(db *storage.TimescaleDB) *TelemetryService {
    return &TelemetryService{db: db}
}

func (s *TelemetryService) GetTelemetry(ctx context.Context, req TelemetryRequest) (*TelemetryResponse, error) {
    var from, to time.Time
    var err error
    
    // Parse time range
    if req.From != "" {
        from, err = time.Parse(time.RFC3339, req.From)
        if err != nil {
            return nil, fmt.Errorf("invalid from time: %w", err)
        }
    } else {
        from = time.Now().Add(-24 * time.Hour)
    }
    
    if req.To != "" {
        to, err = time.Parse(time.RFC3339, req.To)
        if err != nil {
            return nil, fmt.Errorf("invalid to time: %w", err)
        }
    } else {
        to = time.Now()
    }
    
    // Determine interval
    interval := req.Interval
    if interval == "" {
        interval = s.determineInterval(from, to)
    }
    
    // Query data
    points, err := s.db.QueryTelemetry(ctx, req.DeviceID, req.Metric, from, to, interval)
    if err != nil {
        return nil, err
    }
    
    return &TelemetryResponse{
        DeviceID: req.DeviceID,
        Metric:   req.Metric,
        From:     from,
        To:       to,
        Interval: interval,
        Points:   points,
    }, nil
}

func (s *TelemetryService) determineInterval(from, to time.Time) string {
    duration := to.Sub(from)
    
    switch {
    case duration <= time.Hour:
        return "1m"   // 1 minute
    case duration <= 6*time.Hour:
        return "5m"   // 5 minutes
    case duration <= 24*time.Hour:
        return "30m"  // 30 minutes
    case duration <= 7*24*time.Hour:
        return "2h"   // 2 hours
    case duration <= 30*24*time.Hour:
        return "12h"  // 12 hours
    default:
        return "1d"   // 1 day
    }
}

type TelemetryRequest struct {
    DeviceID string
    Metric   string
    From     string
    To       string
    Interval string
}

type TelemetryResponse struct {
    DeviceID string
    Metric   string
    From     time.Time
    To       time.Time
    Interval string
    Points   []storage.TelemetryPoint
}
```

---

## บทที่ 5: Alarm Analysis

---

### 5.1 Real-time Monitoring

#### Alert Engine

```go
// iot-go-project/internal/alert/engine.go
package alert

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "sync"
    "time"

    "iot-go-project/internal/model"
    "iot-go-project/internal/storage"
)

type Engine struct {
    rules      map[string]*model.AlertRule
    alerts     []*model.Alert
    mu         sync.RWMutex
    storage    *storage.TimescaleDB
    mqClient   *mqtt.ClientManager
    notifyChan chan *model.Alert
    ctx        context.Context
    cancel     context.CancelFunc
    wg         sync.WaitGroup
}

func NewEngine(storage *storage.TimescaleDB, mqClient *mqtt.ClientManager) *Engine {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &Engine{
        rules:      make(map[string]*model.AlertRule),
        alerts:     make([]*model.Alert, 0),
        storage:    storage,
        mqClient:   mqClient,
        notifyChan: make(chan *model.Alert, 1000),
        ctx:        ctx,
        cancel:     cancel,
    }
}

func (e *Engine) Start() error {
    // Load rules from database
    if err := e.loadRules(); err != nil {
        return err
    }
    
    // Start alert processor
    e.wg.Add(1)
    go e.processAlerts()
    
    // Start notification worker
    e.wg.Add(1)
    go e.notifyWorker()
    
    log.Printf("Alert engine started with %d rules", len(e.rules))
    return nil
}

func (e *Engine) Stop() {
    e.cancel()
    e.wg.Wait()
    log.Println("Alert engine stopped")
}

func (e *Engine) loadRules() error {
    // Implementation to load from database
    return nil
}

func (e *Engine) ProcessTelemetry(deviceID string, timestamp time.Time, 
    metrics map[string]interface{}) {
    
    e.mu.RLock()
    defer e.mu.RUnlock()
    
    for _, rule := range e.rules {
        if !rule.Enabled {
            continue
        }
        
        if rule.DeviceID != "" && rule.DeviceID != deviceID {
            continue
        }
        
        value, ok := metrics[rule.Metric]
        if !ok {
            continue
        }
        
        if e.evaluateRule(rule, value) {
            e.triggerAlert(rule, deviceID, value, timestamp)
        }
    }
}

func (e *Engine) evaluateRule(rule *model.AlertRule, value interface{}) bool {
    // Implementation of rule evaluation
    // ...
    return false
}

func (e *Engine) triggerAlert(rule *model.AlertRule, deviceID string, 
    value interface{}, timestamp time.Time) {
    
    alert := &model.Alert{
        ID:        uuid.New().String(),
        DeviceID:  deviceID,
        RuleID:    rule.ID,
        Severity:  rule.Severity,
        Message:   rule.Message,
        Status:    "pending",
        Metadata: map[string]interface{}{
            "value": value,
            "rule":  rule.Name,
        },
        CreatedAt: timestamp,
    }
    
    e.mu.Lock()
    e.alerts = append(e.alerts, alert)
    e.mu.Unlock()
    
    // Store alert
    e.storage.InsertAlert(alert)
    
    // Send to notification channel
    select {
    case e.notifyChan <- alert:
    default:
        log.Printf("Notification channel full, dropping alert: %s", alert.ID)
    }
}

func (e *Engine) processAlerts() {
    defer e.wg.Done()
    
    // Subscribe to telemetry topic
    e.mqClient.Subscribe("devices/+/telemetry", func(client mqtt.Client, msg mqtt.Message) {
        var data struct {
            DeviceID  string                 `json:"device_id"`
            Timestamp time.Time              `json:"timestamp"`
            Values    map[string]interface{} `json:"values"`
        }
        
        if err := json.Unmarshal(msg.Payload(), &data); err != nil {
            log.Printf("Failed to unmarshal telemetry: %v", err)
            return
        }
        
        e.ProcessTelemetry(data.DeviceID, data.Timestamp, data.Values)
    })
    
    <-e.ctx.Done()
}
```

---

### 5.2 การตั้งค่า Alarm

#### Rule Management

```go
// iot-go-project/internal/alert/rule.go
package alert

import (
    "context"
    "fmt"

    "iot-go-project/internal/model"
)

func (e *Engine) CreateRule(ctx context.Context, rule *model.AlertRule) error {
    if rule.ID == "" {
        rule.ID = uuid.New().String()
    }
    rule.CreatedAt = time.Now()
    rule.UpdatedAt = time.Now()
    
    // Validate rule
    if err := e.validateRule(rule); err != nil {
        return err
    }
    
    // Save to database
    if err := e.storage.InsertRule(ctx, rule); err != nil {
        return err
    }
    
    // Add to memory
    e.mu.Lock()
    e.rules[rule.ID] = rule
    e.mu.Unlock()
    
    log.Printf("Alert rule created: %s", rule.ID)
    return nil
}

func (e *Engine) validateRule(rule *model.AlertRule) error {
    if rule.Name == "" {
        return fmt.Errorf("rule name is required")
    }
    
    if rule.Metric == "" {
        return fmt.Errorf("metric is required")
    }
    
    validConditions := []string{"gt", "gte", "lt", "lte", "eq", "neq", "between"}
    for _, c := range validConditions {
        if c == rule.Condition {
            goto valid
        }
    }
    return fmt.Errorf("invalid condition: %s", rule.Condition)
    
valid:
    if rule.Severity == "" {
        rule.Severity = "info"
    }
    
    validSeverity := []string{"info", "warning", "critical"}
    for _, s := range validSeverity {
        if s == rule.Severity {
            return nil
        }
    }
    
    return fmt.Errorf("invalid severity: %s", rule.Severity)
}

func (e *Engine) UpdateRule(ctx context.Context, id string, updates map[string]interface{}) (*model.AlertRule, error) {
    e.mu.RLock()
    rule, exists := e.rules[id]
    e.mu.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("rule %s not found", id)
    }
    
    // Apply updates
    // ...
    
    rule.UpdatedAt = time.Now()
    
    // Save to database
    if err := e.storage.UpdateRule(ctx, rule); err != nil {
        return nil, err
    }
    
    e.mu.Lock()
    e.rules[id] = rule
    e.mu.Unlock()
    
    return rule, nil
}

func (e *Engine) DeleteRule(ctx context.Context, id string) error {
    e.mu.Lock()
    defer e.mu.Unlock()
    
    if _, exists := e.rules[id]; !exists {
        return fmt.Errorf("rule %s not found", id)
    }
    
    delete(e.rules, id)
    return e.storage.DeleteRule(ctx, id)
}
```

---

### 5.3 การแจ้งเตือนเมื่อเกิดเหตุการณ์

#### Notification System

```go
// iot-go-project/internal/alert/notification.go
package alert

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/smtp"
    "time"

    "iot-go-project/internal/model"
)

type NotificationService struct {
    emailConfig  EmailConfig
    webhookURL   string
    httpClient   *http.Client
}

type EmailConfig struct {
    SMTPHost string
    SMTPPort string
    Username string
    Password string
    From     string
    To       []string
}

func NewNotificationService(emailConfig EmailConfig, webhookURL string) *NotificationService {
    return &NotificationService{
        emailConfig: emailConfig,
        webhookURL:  webhookURL,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (e *Engine) notifyWorker() {
    defer e.wg.Done()
    
    for {
        select {
        case <-e.ctx.Done():
            return
        case alert := <-e.notifyChan:
            e.sendNotifications(alert)
        }
    }
}

func (e *Engine) sendNotifications(alert *model.Alert) {
    log.Printf("Sending notification for alert: %s", alert.ID)
    
    // Send email
    if err := e.notificationService.SendEmail(alert); err != nil {
        log.Printf("Failed to send email: %v", err)
    }
    
    // Send webhook
    if err := e.notificationService.SendWebhook(alert); err != nil {
        log.Printf("Failed to send webhook: %v", err)
    }
}

func (n *NotificationService) SendEmail(alert *model.Alert) error {
    if len(n.emailConfig.To) == 0 {
        return nil
    }
    
    subject := fmt.Sprintf("[%s] Alert: %s", alert.Severity, alert.Message)
    body := fmt.Sprintf(`
Alert Details:
- ID: %s
- Device: %s
- Severity: %s
- Message: %s
- Time: %s
- Metadata: %v
`,
        alert.ID,
        alert.DeviceID,
        alert.Severity,
        alert.Message,
        alert.CreatedAt.Format(time.RFC3339),
        alert.Metadata,
    )
    
    msg := []byte(fmt.Sprintf(
        "From: %s\r\n"+
        "To: %s\r\n"+
        "Subject: %s\r\n"+
        "\r\n"+
        "%s\r\n",
        n.emailConfig.From,
        n.emailConfig.To[0],
        subject,
        body,
    ))
    
    auth := smtp.PlainAuth("", n.emailConfig.Username, n.emailConfig.Password, n.emailConfig.SMTPHost)
    addr := fmt.Sprintf("%s:%s", n.emailConfig.SMTPHost, n.emailConfig.SMTPPort)
    
    return smtp.SendMail(addr, auth, n.emailConfig.From, n.emailConfig.To, msg)
}

func (n *NotificationService) SendWebhook(alert *model.Alert) error {
    if n.webhookURL == "" {
        return nil
    }
    
    data, err := json.Marshal(alert)
    if err != nil {
        return err
    }
    
    resp, err := n.httpClient.Post(n.webhookURL, "application/json", bytes.NewReader(data))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return fmt.Errorf("webhook returned %s", resp.Status)
    }
    
    return nil
}
```

---

## บทที่ 6: Data Visualization

---

### 6.1 การแสดงข้อมูลแบบ Real-time

#### WebSocket Integration

```go
// iot-go-project/internal/websocket/hub.go
package websocket

import (
    "encoding/json"
    "log"
    "sync"
    "time"

    "github.com/gorilla/websocket"
)

type Hub struct {
    clients      map[*Client]bool
    broadcast    chan []byte
    register     chan *Client
    unregister   chan *Client
    mu           sync.RWMutex
}

type Client struct {
    Hub      *Hub
    Conn     *websocket.Conn
    Send     chan []byte
    DeviceID string
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
            log.Printf("WebSocket client connected: %s", client.DeviceID)
            
        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.Send)
            }
            h.mu.Unlock()
            log.Printf("WebSocket client disconnected: %s", client.DeviceID)
            
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

func (h *Hub) BroadcastTelemetry(deviceID string, data interface{}) {
    payload := map[string]interface{}{
        "type": "telemetry",
        "device_id": deviceID,
        "data": data,
        "timestamp": time.Now(),
    }
    
    msg, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Failed to marshal telemetry: %v", err)
        return
    }
    
    h.broadcast <- msg
}
```

---

### 6.2 Dashboard Design

#### Real-time Dashboard API

```go
// iot-go-project/internal/api/dashboard.go
package api

import (
    "encoding/json"
    "net/http"
    "time"
)

type DashboardAPI struct {
    deviceService *service.DeviceService
    alertService  *service.AlertService
    wsHub         *websocket.Hub
}

func (d *DashboardAPI) GetDashboardData(w http.ResponseWriter, r *http.Request) {
    // Get device stats
    devices, _ := d.deviceService.ListDevices(r.Context(), nil)
    
    online := 0
    offline := 0
    for _, device := range devices {
        if device.Status == "online" {
            online++
        } else {
            offline++
        }
    }
    
    // Get alert stats
    alerts, _ := d.alertService.GetAlerts(r.Context(), AlertFilter{
        Limit: 100,
    })
    
    pending := 0
    critical := 0
    for _, alert := range alerts {
        if alert.Status == "pending" {
            pending++
        }
        if alert.Severity == "critical" {
            critical++
        }
    }
    
    // Get total messages (last 24 hours)
    totalMessages, _ := d.deviceService.GetMessageCount(r.Context(), time.Now().Add(-24*time.Hour), time.Now())
    
    response := DashboardResponse{
        DeviceStats: DeviceStats{
            Total:   len(devices),
            Online:  online,
            Offline: offline,
        },
        AlertStats: AlertStats{
            Total:    len(alerts),
            Pending:  pending,
            Critical: critical,
        },
        MessageStats: MessageStats{
            Total: totalMessages,
        },
        Timestamp: time.Now(),
    }
    
    json.NewEncoder(w).Encode(response)
}

func (d *DashboardAPI) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
    upgrader := websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool { return true },
    }
    
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, "Failed to upgrade connection", http.StatusBadRequest)
        return
    }
    
    deviceID := r.URL.Query().Get("device")
    client := &websocket.Client{
        Hub:      d.wsHub,
        Conn:     conn,
        Send:     make(chan []byte, 256),
        DeviceID: deviceID,
    }
    
    d.wsHub.register <- client
    
    // Write pump
    go client.WritePump()
    
    // Read pump
    client.ReadPump()
}
```

---

### 6.3 การสร้างกราฟและ图表

#### Chart Data API

```go
// iot-go-project/internal/api/chart.go
package api

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"
)

func (d *DashboardAPI) GetChartData(w http.ResponseWriter, r *http.Request) {
    deviceID := r.PathValue("id")
    metric := r.URL.Query().Get("metric")
    timeRange := r.URL.Query().Get("range")
    
    if timeRange == "" {
        timeRange = "24h"
    }
    
    // Parse time range
    var from, to time.Time
    to = time.Now()
    
    switch timeRange {
    case "1h":
        from = to.Add(-1 * time.Hour)
    case "6h":
        from = to.Add(-6 * time.Hour)
    case "24h":
        from = to.Add(-24 * time.Hour)
    case "7d":
        from = to.Add(-7 * 24 * time.Hour)
    case "30d":
        from = to.Add(-30 * 24 * time.Hour)
    default:
        from = to.Add(-24 * time.Hour)
    }
    
    // Determine interval
    duration := to.Sub(from)
    interval := "1m"
    if duration > 24*time.Hour {
        interval = "1h"
    }
    if duration > 7*24*time.Hour {
        interval = "12h"
    }
    if duration > 30*24*time.Hour {
        interval = "1d"
    }
    
    // Query data
    points, err := d.telemetryService.GetTelemetry(r.Context(), service.TelemetryRequest{
        DeviceID: deviceID,
        Metric:   metric,
        From:     from.Format(time.RFC3339),
        To:       to.Format(time.RFC3339),
        Interval: interval,
    })
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Format for chart
    chartData := ChartData{
        Labels:    make([]string, len(points.Points)),
        Datasets: []Dataset{
            {
                Label:     metric,
                Data:      make([]float64, len(points.Points)),
                Fill:      false,
                Tension:   0.4,
                PointRadius: 2,
            },
        },
    }
    
    for i, point := range points.Points {
        chartData.Labels[i] = point.Timestamp.Format("15:04")
        chartData.Datasets[0].Data[i] = point.Avg
    }
    
    json.NewEncoder(w).Encode(chartData)
}

type ChartData struct {
    Labels   []string  `json:"labels"`
    Datasets []Dataset `json:"datasets"`
}

type Dataset struct {
    Label          string    `json:"label"`
    Data           []float64 `json:"data"`
    Fill           bool      `json:"fill"`
    Tension        float64   `json:"tension"`
    PointRadius    int       `json:"pointRadius"`
    BorderColor    string    `json:"borderColor,omitempty"`
    BackgroundColor string   `json:"backgroundColor,omitempty"`
}
```

---

## บทที่ 7: Offline Computing

---

### 7.1 การประมวลผลข้อมูลย้อนหลัง

#### Batch Processing

```go
// iot-go-project/internal/processor/batch.go
package processor

import (
    "context"
    "log"
    "sync"
    "time"

    "iot-go-project/internal/model"
    "iot-go-project/internal/storage"
)

type BatchProcessor struct {
    storage   *storage.TimescaleDB
    batchSize int
    interval  time.Duration
    wg        sync.WaitGroup
    ctx       context.Context
    cancel    context.CancelFunc
}

func NewBatchProcessor(storage *storage.TimescaleDB) *BatchProcessor {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &BatchProcessor{
        storage:   storage,
        batchSize: 1000,
        interval:  5 * time.Minute,
        ctx:       ctx,
        cancel:    cancel,
    }
}

func (p *BatchProcessor) Start() {
    p.wg.Add(1)
    go p.run()
    log.Println("Batch processor started")
}

func (p *BatchProcessor) Stop() {
    p.cancel()
    p.wg.Wait()
    log.Println("Batch processor stopped")
}

func (p *BatchProcessor) run() {
    defer p.wg.Done()
    
    ticker := time.NewTicker(p.interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-p.ctx.Done():
            return
        case <-ticker.C:
            p.processBatch()
        }
    }
}

func (p *BatchProcessor) processBatch() {
    log.Println("Processing batch...")
    
    // Get unprocessed data
    data, err := p.storage.GetUnprocessedTelemetry(p.ctx, p.batchSize)
    if err != nil {
        log.Printf("Failed to get unprocessed data: %v", err)
        return
    }
    
    if len(data) == 0 {
        return
    }
    
    // Process data
    for _, item := range data {
        p.processTelemetry(item)
    }
    
    // Mark as processed
    if err := p.storage.MarkProcessed(p.ctx, data); err != nil {
        log.Printf("Failed to mark as processed: %v", err)
    }
    
    log.Printf("Processed %d items", len(data))
}

func (p *BatchProcessor) processTelemetry(telemetry *model.Telemetry) {
    // Apply transformations
    // Calculate aggregates
    // Store results
}
```

---

### 7.2 Data Analytics

#### Analytics Service

```go
// iot-go-project/internal/analytics/service.go
package analytics

import (
    "context"
    "math"
    "time"

    "iot-go-project/internal/model"
    "iot-go-project/internal/storage"
)

type AnalyticsService struct {
    storage *storage.TimescaleDB
}

func NewAnalyticsService(storage *storage.TimescaleDB) *AnalyticsService {
    return &AnalyticsService{storage: storage}
}

func (s *AnalyticsService) CalculateStats(ctx context.Context, deviceID string, 
    metric string, from, to time.Time) (*DeviceStats, error) {
    
    points, err := s.storage.QueryTelemetry(ctx, deviceID, metric, from, to, "5m")
    if err != nil {
        return nil, err
    }
    
    if len(points) == 0 {
        return &DeviceStats{
            DeviceID: deviceID,
            Metric:   metric,
            From:     from,
            To:       to,
            Count:    0,
        }, nil
    }
    
    stats := &DeviceStats{
        DeviceID:   deviceID,
        Metric:     metric,
        From:       from,
        To:         to,
        Count:      len(points),
        Points:     points,
    }
    
    // Calculate statistics
    var sum, sumSq float64
    min := math.MaxFloat64
    max := -math.MaxFloat64
    
    for _, p := range points {
        val := p.Avg
        sum += val
        sumSq += val * val
        
        if val < min {
            min = val
        }
        if val > max {
            max = val
        }
    }
    
    avg := sum / float64(len(points))
    variance := (sumSq / float64(len(points))) - (avg * avg)
    stdDev := math.Sqrt(variance)
    
    stats.Min = min
    stats.Max = max
    stats.Avg = avg
    stats.StdDev = stdDev
    
    return stats, nil
}

type DeviceStats struct {
    DeviceID string
    Metric   string
    From     time.Time
    To       time.Time
    Count    int
    Min      float64
    Max      float64
    Avg      float64
    StdDev   float64
    Points   []storage.TelemetryPoint
}
```

---

## บทที่ 8: การติดตั้งและ Deployment

---

### 8.1 Deployment Guide

#### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: timescale/timescaledb:latest-pg14
    environment:
      POSTGRES_DB: iot
      POSTGRES_USER: iot
      POSTGRES_PASSWORD: iot_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U iot"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  emqx:
    image: emqx/emqx:latest
    ports:
      - "1883:1883"      # MQTT
      - "8883:8883"      # MQTT SSL
      - "8083:8083"      # WebSocket
      - "8084:8084"      # WebSocket SSL
      - "18083:18083"    # Dashboard
    environment:
      EMQX_NAME: emqx
    healthcheck:
      test: ["CMD", "emqx", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  rabbitmq:
    image: rabbitmq:3-management-alpine
    ports:
      - "5672:5672"
      - "15672:15672"
    environment:
      RABBITMQ_DEFAULT_USER: iot
      RABBITMQ_DEFAULT_PASS: iot_password
    healthcheck:
      test: ["CMD", "rabbitmq-diagnostics", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  go-iot:
    build:
      context: ./go-iot
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      MQTT_BROKER: tcp://emqx:1883
      MQTT_USERNAME: admin
      MQTT_PASSWORD: public
      DB_CONNECTION: postgres://iot:iot_password@postgres:5432/iot
      REDIS_ADDR: redis:6379
      RABBITMQ_URL: amqp://iot:iot_password@rabbitmq:5672
      LOG_LEVEL: info
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      emqx:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy
    restart: unless-stopped

  iot-api:
    build:
      context: ./iot-go-project
      dockerfile: Dockerfile
    ports:
      - "8081:8081"
    environment:
      DB_CONNECTION: postgres://iot:iot_password@postgres:5432/iot
      REDIS_ADDR: redis:6379
      MQTT_BROKER: tcp://emqx:1883
      RABBITMQ_URL: amqp://iot:iot_password@rabbitmq:5672
      API_PORT: 8081
      JWT_SECRET: your_jwt_secret
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      emqx:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy
    restart: unless-stopped

  ant-vue:
    build:
      context: ./ant-vue
      dockerfile: Dockerfile
    ports:
      - "8082:80"
    environment:
      API_URL: http://iot-api:8081
      WS_URL: ws://iot-api:8081/ws
    depends_on:
      - iot-api
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
```

#### Dockerfile

```dockerfile
# go-iot/Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o go-iot ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/go-iot .

EXPOSE 8080

CMD ["./go-iot"]
```

---

### 8.2 การตั้งค่า Environment

#### Environment Variables

```bash
# .env
# Database
DB_CONNECTION=postgres://iot:iot_password@postgres:5432/iot?sslmode=disable

# Redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0

# MQTT
MQTT_BROKER=tcp://emqx:1883
MQTT_USERNAME=admin
MQTT_PASSWORD=public
MQTT_CLIENT_ID=go-iot-service

# RabbitMQ
RABBITMQ_URL=amqp://iot:iot_password@rabbitmq:5672

# API
API_PORT=8081
API_HOST=0.0.0.0
JWT_SECRET=your_jwt_secret_here
JWT_EXPIRATION=24h

# CORS
CORS_ALLOWED_ORIGINS=*

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Monitoring
ENABLE_METRICS=true
METRICS_PORT=9090

# Alert
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=alerts@example.com
SMTP_PASSWORD=your_smtp_password
ALERT_EMAILS=admin@example.com

# Webhook
WEBHOOK_URL=https://hooks.slack.com/services/xxx

# Batch Processing
BATCH_INTERVAL=5m
BATCH_SIZE=1000
```

---

### 8.3 การรันใน Production

#### Kubernetes Deployment

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-iot
  namespace: iot
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-iot
  template:
    metadata:
      labels:
        app: go-iot
    spec:
      containers:
      - name: go-iot
        image: your-registry/go-iot:latest
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 9090
          name: metrics
        env:
        - name: DB_CONNECTION
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: connection
        - name: MQTT_BROKER
          value: tcp://emqx-service:1883
        - name: REDIS_ADDR
          value: redis-service:6379
        - name: RABBITMQ_URL
          valueFrom:
            secretKeyRef:
              name: rabbitmq-secret
              key: url
        - name: LOG_LEVEL
          value: "info"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: go-iot-service
  namespace: iot
spec:
  selector:
    app: go-iot
  ports:
  - name: http
    port: 8080
    targetPort: 8080
  - name: metrics
    port: 9090
    targetPort: 9090
  type: ClusterIP
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: go-iot-hpa
  namespace: iot
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: go-iot
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

#### Monitoring with Prometheus

```go
// iot-go-project/internal/metrics/metrics.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    DevicesConnected = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "iot_devices_connected",
            Help: "Number of connected devices",
        },
        []string{"type"},
    )
    
    TelemetryReceived = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "iot_telemetry_total",
            Help: "Total telemetry messages received",
        },
    )
    
    AlertsTriggered = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "iot_alerts_total",
            Help: "Total alerts triggered",
        },
        []string{"severity"},
    )
    
    MQTTErrors = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "iot_mqtt_errors_total",
            Help: "Total MQTT errors",
        },
    )
    
    ProcessingDuration = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name: "iot_processing_duration_seconds",
            Help: "Duration of telemetry processing",
            Buckets: prometheus.DefBuckets,
        },
    )
)

func RecordDeviceConnected(deviceType string) {
    DevicesConnected.WithLabelValues(deviceType).Inc()
}

func RecordTelemetryReceived() {
    TelemetryReceived.Inc()
}

func RecordAlert(severity string) {
    AlertsTriggered.WithLabelValues(severity).Inc()
}

func RecordMQTTError() {
    MQTTErrors.Inc()
}
```

---

## บทที่ 9: การมีส่วนร่วมและเอกสาร

---

### 9.1 การรายงาน Issues

#### Issue Template

```markdown
## Bug Report

### Description
A clear description of the bug

### Steps to Reproduce
1. Step 1
2. Step 2
3. Step 3

### Expected Behavior
What should happen

### Actual Behavior
What actually happens

### Environment
- OS: [e.g., Ubuntu 22.04]
- Go Version: [e.g., 1.21]
- Version: [e.g., v1.0.0]

### Additional Context
Any other information
```

### 9.2 การส่ง Pull Requests

#### PR Template

```markdown
## Description
Describe the changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added
- [ ] Integration tests added
- [ ] Manual testing done

## Checklist
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] Tests passing
- [ ] No new warnings

## Screenshots (if applicable)
```

### 9.3 การปรับปรุง Documentation

#### API Documentation

```go
// Swagger annotations
// @Summary Register new device
// @Description Register a new IoT device
// @Tags devices
// @Accept json
// @Produce json
// @Param device body CreateDeviceRequest true "Device info"
// @Success 201 {object} DeviceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /devices [post]
func (h *Handler) CreateDevice(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

---

## 📝 สรุปเล่มที่ 7

### สิ่งที่เราเรียนรู้

| หัวข้อ | รายละเอียด |
|--------|------------|
| **สถาปัตยกรรม** | Microservices, MQTT, Message Queue |
| **MQTT Management** | Connection, Reconnect, Pool |
| **Data Storage** | TimescaleDB, Time-series |
| **Alert Engine** | Rules, Evaluation, Notification |
| **Visualization** | WebSocket, Dashboard, Charts |
| **Batch Processing** | Offline computing, Analytics |
| **Deployment** | Docker, Kubernetes, Monitoring |

### สถาปัตยกรรมที่สมบูรณ์

```mermaid
flowchart TB
    subgraph "IoT Platform"
        M[Device Manager]
        Q[Message Queue]
        S[Storage]
        A[Alert Engine]
        V[Visualization]
        B[Batch Processor]
        API[REST API]
    end
    
    subgraph "Infrastructure"
        D[Docker]
        K[Kubernetes]
        P[Prometheus]
        G[Grafana]
    end
    
    M --> S
    M --> A
    M --> V
    Q --> S
    S --> B
    B --> V
    API --> S
    API --> V
    
    D --> K
    K --> M
    K --> Q
    K --> S
    K --> A
    K --> V
    K --> B
    K --> API
    
    P --> M
    P --> Q
    P --> S
    P --> A
    P --> V
    P --> B
    P --> API
    
    G --> P
```

---

**จบเล่มที่ 7: โครงการ IoT ขนาดใหญ่** 🎉

---

## 🏁 สรุปภาพรวมหนังสือชุด “เรียน Go อย่างมืออาชีพ”

| เล่ม | ชื่อ | หัวข้อหลัก |
|------|------|-----------|
| 1 | พื้นฐานภาษา Go | Variables, Functions, Pointers, Packages |
| 2 | โครงสร้างข้อมูลขั้นสูงและการจัดการ JSON | Structs, Interfaces, JSON, Templates |
| 3 | การทำงานพร้อมกัน — หัวใจของ Go | Goroutines, Channels, Select, Scheduler |
| 4 | การพัฒนาเว็บและเครือข่าย | HTTP, TCP, MQTT, Web App, REST API |
| 5 | การทดสอบและการปรับปรุงประสิทธิภาพ | Testing, Benchmark, Profiling, Modules |
| 6 | 10 โครงการฝึกปฏิบัติ | To-Do, Game, Weather, Chat, Blog, etc. |
| 7 | โครงการ IoT ขนาดใหญ่ | MQTT, Data Storage, Alarm, Visualization |

### ขอให้สนุกกับการเขียน Go! 🚀

---

**THE END** 🎉