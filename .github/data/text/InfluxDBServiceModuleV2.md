# เอกสารฉบับสมบูรณ์: InfluxDB Service Module + IoT Stack (MQTT, Socket.IO, Redis, PostgreSQL) ร่วมกับ LLM & AI Agent

---

## สารบัญ

1. [ภาพรวมสถาปัตยกรรม IoT + LLM/AI Agent](#1-ภาพรวมสถาปัตยกรรม-iot---llmai-agent)
2. [โครงสร้างโปรเจกต์ที่ขยายเพิ่ม](#2-โครงสร้างโปรเจกต์ที่ขยายเพิ่ม)
3. [โมดูล MQTT](#3-โมดูล-mqtt)
4. [โมดูล Socket.IO](#4-โมดูล-socketio)
5. [โมดูล Redis (Cache & Stream)](#5-โมดูล-redis-cache--stream)
6. [โมดูล PostgreSQL (Alarm & History)](#6-โมดูล-postgresql-alarm--history)
7. [โมดูล LLM + AI Agent](#7-โมดูล-llm--ai-agent)
8. [การเชื่อมต่อทุกอย่างเข้าด้วยกัน](#8-การเชื่อมต่อทุกอย่างเข้าด้วยกัน)
9. [Dependency Injection (Wire)](#9-dependency-injection-wire)
10. [การ Deploy ด้วย Docker Compose](#10-การ-deploy-ด้วย-docker-compose)
11. [สรุป](#11-สรุป)

---

## 1. ภาพรวมสถาปัตยกรรม IoT + LLM / AI Agent

สถาปัตยกรรมนี้ขยายจาก DDD + Clean Architecture เดิม เพื่อรองรับระบบ IoT ที่เชื่อมต่อกับ LLM และ AI Agent โดยมีส่วนประกอบหลักดังนี้:

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              IoT Devices / Sensors                              │
│                                   (MQTT Publisher)                              │
└─────────────────────────────────────┬───────────────────────────────────────────┘
                                      │ MQTT
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              MQTT Broker (EMQX/Mosquitto)                       │
└─────────────────────────────────────┬───────────────────────────────────────────┘
                                      │
┌─────────────────────────────────────▼───────────────────────────────────────────┐
│                                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                         MQTT Consumer (Worker Pool)                      │   │
│  │                    - Subscribe to topics                                 │   │
│  │                    - Parse & validate data                               │   │
│  │                    - Publish to Redis Stream / Channel                   │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                      │                                          │
│                                      ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                              Redis                                       │   │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────────┐ │   │
│  │  │   Stream/Queue   │  │     Cache       │  │  Session / Rate Limit   │ │   │
│  │  │  (IoT Data)      │  │  (LLM Response) │  │                         │ │   │
│  │  └─────────────────┘  └─────────────────┘  └─────────────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                      │                                          │
│                                      ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                         InfluxDB (Time-Series)                          │   │
│  │                    - Store sensor data                                   │   │
│  │                    - Query for trends & anomalies                       │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                      │                                          │
│                                      ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                      Alarm Engine (Rules Engine)                        │   │
│  │                    - Evaluate conditions                                 │   │
│  │                    - Trigger alerts                                      │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                      │                                          │
│                                      ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                         PostgreSQL                                       │   │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────────┐ │   │
│  │  │  Alarm History  │  │   Device Registry│  │  LLM Interaction Log   │ │   │
│  │  └─────────────────┘  └─────────────────┘  └─────────────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                      │                                          │
│                                      ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                      LLM / AI Agent Service                             │   │
│  │  ┌─────────────────────────────────────────────────────────────────┐   │   │
│  │  │  • Anomaly Detection  • Predictive Maintenance                  │   │   │
│  │  │  • Natural Language Query  • Automated Response                 │   │   │
│  │  └─────────────────────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                      │                                          │
│                                      ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                    Socket.IO (Real-time Dashboard)                      │   │
│  │               - Push real-time updates to clients                       │   │
│  │               - Bidirectional communication                            │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. โครงสร้างโปรเจกต์ที่ขยายเพิ่ม

```
icmongolang/
├── cmd/
│   └── api/
│       ├── main.go                 # Entrypoint
│       └── bootstrap/
│           └── bootstrap.go        # Application bootstrap
│
├── internal/
│   ├── core/                       # Core business logic
│   │   ├── domain/
│   │   │   ├── metric.go           # Metric entity
│   │   │   ├── metric_repository.go
│   │   │   ├── device.go           # Device entity
│   │   │   ├── alarm.go            # Alarm entity
│   │   │   ├── llm_request.go      # LLM request/response
│   │   │   └── errors.go
│   │   │
│   │   ├── application/
│   │   │   ├── command/
│   │   │   │   ├── write_metric.go
│   │   │   │   ├── process_iot_data.go
│   │   │   │   ├── trigger_alarm.go
│   │   │   │   └── query_llm.go
│   │   │   ├── query/
│   │   │   │   ├── query_metric.go
│   │   │   │   ├── query_device.go
│   │   │   │   └── query_alarm.go
│   │   │   └── interfaces.go
│   │   │
│   │   └── ports/
│   │       ├── metric_port.go
│   │       ├── device_port.go
│   │       ├── alarm_port.go
│   │       ├── cache_port.go
│   │       └── llm_port.go
│   │
│   ├── adapters/
│   │   ├── in/                     # Inbound adapters
│   │   │   ├── http/
│   │   │   │   ├── handler.go
│   │   │   │   └── dto.go
│   │   │   ├── websocket/          # Socket.IO
│   │   │   │   ├── handler.go
│   │   │   │   └── events.go
│   │   │   └── mqtt/               # MQTT Consumer
│   │   │       ├── consumer.go
│   │   │       ├── handler.go
│   │   │       └── worker_pool.go
│   │   │
│   │   └── out/                    # Outbound adapters
│   │       ├── influxdb/
│   │       │   ├── client.go
│   │       │   ├── metric_repository.go
│   │       │   └── config.go
│   │       ├── postgres/
│   │       │   ├── client.go
│   │       │   ├── device_repository.go
│   │       │   ├── alarm_repository.go
│   │       │   └── models/
│   │       │       └── *.go       # GORM models
│   │       ├── redis/
│   │       │   ├── client.go
│   │       │   ├── cache_repository.go
│   │       │   └── stream_publisher.go
│   │       └── llm/
│   │           ├── client.go      # OpenAI/Anthropic/etc.
│   │           ├── agent.go       # AI Agent logic
│   │           └── prompt.go      # Prompt templates
│   │
│   ├── di/                         # Dependency Injection
│   │   └── wire.go
│   │
│   └── pkg/                        # Shared utilities
│       ├── logger/
│       ├── errors/
│       └── context/
│
├── configs/
│   └── config.yaml
│
├── deployments/
│   └── docker-compose.yml
│
├── go.mod
└── go.sum
```

---

## 3. โมดูล MQTT

MQTT เป็นโปรโตคอลหลักสำหรับรับข้อมูลจากอุปกรณ์ IoT ใช้ไลบรารี **Eclipse Paho** 

### 3.1 MQTT Client (Infrastructure)

```go
// internal/adapters/in/mqtt/client.go
package mqtt

import (
    "context"
    "crypto/tls"
    "fmt"
    "time"

    mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Config holds MQTT connection configuration
type Config struct {
    Broker     string        `yaml:"broker"`
    ClientID   string        `yaml:"client_id"`
    Username   string        `yaml:"username"`
    Password   string        `yaml:"password"`
    Topics     []string      `yaml:"topics"`
    QOS        byte          `yaml:"qos"`
    Timeout    time.Duration `yaml:"timeout"`
    KeepAlive  time.Duration `yaml:"keep_alive"`
    WorkerSize int           `yaml:"worker_size"`
}

// Client wraps the Paho MQTT client
type Client struct {
    client   mqtt.Client
    config   *Config
    onMessage MessageHandler
}

type MessageHandler func(ctx context.Context, topic string, payload []byte) error

// NewClient creates a new MQTT client
func NewClient(config *Config, handler MessageHandler) (*Client, error) {
    opts := mqtt.NewClientOptions()
    opts.AddBroker(config.Broker)
    opts.SetClientID(config.ClientID)
    opts.SetUsername(config.Username)
    opts.SetPassword(config.Password)
    opts.SetKeepAlive(config.KeepAlive)
    opts.SetPingTimeout(config.Timeout)
    opts.SetConnectTimeout(config.Timeout)
    opts.SetAutoReconnect(true)
    opts.SetMaxReconnectInterval(10 * time.Second)
    opts.SetCleanSession(true)

    // TLS configuration
    opts.SetTLSConfig(&tls.Config{
        InsecureSkipVerify: false,
    })

    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        return nil, fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
    }

    return &Client{
        client:   client,
        config:   config,
        onMessage: handler,
    }, nil
}

// Subscribe subscribes to topics with worker pool
func (c *Client) Subscribe(ctx context.Context, topics []string) error {
    // Create worker pool for message processing
    pool := NewWorkerPool(c.config.WorkerSize, c.onMessage)
    pool.Start(ctx)

    for _, topic := range topics {
        token := c.client.Subscribe(topic, c.config.QOS, func(client mqtt.Client, msg mqtt.Message) {
            // Offload to worker pool
            pool.Submit(msg.Topic(), msg.Payload())
        })
        if token.Wait() && token.Error() != nil {
            return fmt.Errorf("failed to subscribe to %s: %w", topic, token.Error())
        }
    }

    return nil
}

// Publish publishes a message to a topic
func (c *Client) Publish(topic string, payload []byte) error {
    token := c.client.Publish(topic, c.config.QOS, false, payload)
    if token.Wait() && token.Error() != nil {
        return fmt.Errorf("failed to publish: %w", token.Error())
    }
    return nil
}

// Close closes the MQTT client
func (c *Client) Close() {
    if c.client != nil && c.client.IsConnected() {
        c.client.Disconnect(250)
    }
}
```

### 3.2 MQTT Worker Pool

```go
// internal/adapters/in/mqtt/worker_pool.go
package mqtt

import (
    "context"
    "sync"
)

// WorkerPool processes MQTT messages concurrently
type WorkerPool struct {
    size       int
    handler    MessageHandler
    taskChan   chan *MessageTask
    wg         sync.WaitGroup
    ctx        context.Context
    cancelFunc context.CancelFunc
}

type MessageTask struct {
    Topic   string
    Payload []byte
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(size int, handler MessageHandler) *WorkerPool {
    return &WorkerPool{
        size:     size,
        handler:  handler,
        taskChan: make(chan *MessageTask, size*100), // Buffer
    }
}

// Start starts the worker pool
func (p *WorkerPool) Start(ctx context.Context) {
    p.ctx, p.cancelFunc = context.WithCancel(ctx)

    for i := 0; i < p.size; i++ {
        p.wg.Add(1)
        go p.worker(i)
    }
}

// Submit submits a message to the pool
func (p *WorkerPool) Submit(topic string, payload []byte) {
    select {
    case p.taskChan <- &MessageTask{Topic: topic, Payload: payload}:
    default:
        // Channel full - log and drop or block
        // Consider using a non-blocking send with retry
    }
}

func (p *WorkerPool) worker(id int) {
    defer p.wg.Done()

    for {
        select {
        case <-p.ctx.Done():
            return
        case task := <-p.taskChan:
            if err := p.handler(p.ctx, task.Topic, task.Payload); err != nil {
                // Log error
            }
        }
    }
}

// Stop stops the worker pool
func (p *WorkerPool) Stop() {
    if p.cancelFunc != nil {
        p.cancelFunc()
    }
    p.wg.Wait()
    close(p.taskChan)
}
```

### 3.3 MQTT Message Handler (Use Case)

```go
// internal/adapters/in/mqtt/handler.go
package mqtt

import (
    "context"
    "encoding/json"
    "log/slog"
    "time"

    "your-project/internal/core/application/command"
)

// IoTData represents incoming IoT data
type IoTData struct {
    DeviceID    string                 `json:"device_id"`
    Measurement string                 `json:"measurement"`
    Tags        map[string]string      `json:"tags"`
    Fields      map[string]interface{} `json:"fields"`
    Timestamp   *time.Time             `json:"timestamp,omitempty"`
}

// MessageHandlerImpl handles MQTT messages
type MessageHandlerImpl struct {
    writeMetricHandler *command.WriteMetricHandler
    processIoTHandler  *command.ProcessIoTDataHandler
}

func NewMessageHandlerImpl(
    writeMetricHandler *command.WriteMetricHandler,
    processIoTHandler *command.ProcessIoTDataHandler,
) *MessageHandlerImpl {
    return &MessageHandlerImpl{
        writeMetricHandler: writeMetricHandler,
        processIoTHandler:  processIoTHandler,
    }
}

// Handle processes incoming MQTT message
func (h *MessageHandlerImpl) Handle(ctx context.Context, topic string, payload []byte) error {
    slog.Info("received MQTT message", "topic", topic, "size", len(payload))

    var data IoTData
    if err := json.Unmarshal(payload, &data); err != nil {
        return err
    }

    // Process IoT data - write to InfluxDB and trigger AI analysis
    cmd := command.ProcessIoTDataCommand{
        DeviceID:    data.DeviceID,
        Measurement: data.Measurement,
        Tags:        data.Tags,
        Fields:      data.Fields,
        Timestamp:   data.Timestamp,
    }

    return h.processIoTHandler.Handle(ctx, cmd)
}
```

---

## 4. โมดูล Socket.IO

Socket.IO ให้การสื่อสารแบบ real-time ระหว่าง server และ clients (dashboard, mobile apps)

### 4.1 Socket.IO Server

```go
// internal/adapters/in/websocket/handler.go
package websocket

import (
    "context"
    "encoding/json"
    "log/slog"
    "net/http"

    socketio "github.com/googollee/go-socket.io"
    "github.com/googollee/go-socket.io/engineio"
    "github.com/googollee/go-socket.io/engineio/transport"
    "github.com/googollee/go-socket.io/engineio/transport/websocket"
)

// SocketIOServer manages Socket.IO connections
type SocketIOServer struct {
    server   *socketio.Server
    hub      *Hub
    onEvent  EventHandler
}

type EventHandler func(ctx context.Context, event string, data interface{}) error

// NewSocketIOServer creates a new Socket.IO server
func NewSocketIOServer(handler EventHandler) *SocketIOServer {
    server := socketio.NewServer(&engineio.Options{
        Transports: []transport.Transport{
            &websocket.Transport{
                CheckOrigin: func(r *http.Request) bool {
                    return true // Configure properly in production
                },
            },
        },
    })

    s := &SocketIOServer{
        server:  server,
        hub:     NewHub(),
        onEvent: handler,
    }

    s.setupHandlers()
    return s
}

func (s *SocketIOServer) setupHandlers() {
    // Connection handler
    s.server.OnConnect("/", func(conn socketio.Conn) error {
        slog.Info("client connected", "id", conn.ID())
        s.hub.AddClient(conn)
        return nil
    })

    // Disconnection handler
    s.server.OnDisconnect("/", func(conn socketio.Conn, reason string) {
        slog.Info("client disconnected", "id", conn.ID(), "reason", reason)
        s.hub.RemoveClient(conn)
    })

    // Custom event: iot_data
    s.server.OnEvent("/", "iot_data", func(conn socketio.Conn, data interface{}) {
        if err := s.onEvent(context.Background(), "iot_data", data); err != nil {
            slog.Error("failed to process iot_data event", "error", err)
        }
    })

    // Custom event: query_llm
    s.server.OnEvent("/", "query_llm", func(conn socketio.Conn, data interface{}) {
        if err := s.onEvent(context.Background(), "query_llm", data); err != nil {
            slog.Error("failed to process query_llm event", "error", err)
        }
    })
}

// ServeHTTP serves Socket.IO requests
func (s *SocketIOServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.server.ServeHTTP(w, r)
}

// Broadcast sends a message to all connected clients
func (s *SocketIOServer) Broadcast(event string, data interface{}) error {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return err
    }
    s.server.BroadcastToNamespace("/", event, string(jsonData))
    return nil
}

// Close closes the Socket.IO server
func (s *SocketIOServer) Close() error {
    return s.server.Close()
}

// Hub manages connected clients
type Hub struct {
    clients map[string]socketio.Conn
    mu      sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients: make(map[string]socketio.Conn),
    }
}

func (h *Hub) AddClient(conn socketio.Conn) {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.clients[conn.ID()] = conn
}

func (h *Hub) RemoveClient(conn socketio.Conn) {
    h.mu.Lock()
    defer h.mu.Unlock()
    delete(h.clients, conn.ID())
}

func (h *Hub) GetClients() []socketio.Conn {
    h.mu.RLock()
    defer h.mu.RUnlock()
    clients := make([]socketio.Conn, 0, len(h.clients))
    for _, c := range h.clients {
        clients = append(clients, c)
    }
    return clients
}
```

### 4.2 Socket.IO Events

```go
// internal/adapters/in/websocket/events.go
package websocket

import (
    "context"
    "encoding/json"
    "log/slog"

    "your-project/internal/core/application/command"
    "your-project/internal/core/application/query"
)

// EventHandlerImpl handles Socket.IO events
type EventHandlerImpl struct {
    queryLLMHandler      *command.QueryLLMHandler
    queryMetricHandler   *query.QueryMetricHandler
    writeMetricHandler   *command.WriteMetricHandler
}

func NewEventHandlerImpl(
    queryLLMHandler *command.QueryLLMHandler,
    queryMetricHandler *query.QueryMetricHandler,
    writeMetricHandler *command.WriteMetricHandler,
) *EventHandlerImpl {
    return &EventHandlerImpl{
        queryLLMHandler:    queryLLMHandler,
        queryMetricHandler: queryMetricHandler,
        writeMetricHandler: writeMetricHandler,
    }
}

// Handle processes Socket.IO events
func (h *EventHandlerImpl) Handle(ctx context.Context, event string, data interface{}) error {
    slog.Info("socket event received", "event", event)

    switch event {
    case "iot_data":
        return h.handleIoTData(ctx, data)
    case "query_llm":
        return h.handleQueryLLM(ctx, data)
    default:
        return nil
    }
}

func (h *EventHandlerImpl) handleIoTData(ctx context.Context, data interface{}) error {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return err
    }

    var req struct {
        DeviceID    string                 `json:"device_id"`
        Measurement string                 `json:"measurement"`
        Tags        map[string]string      `json:"tags"`
        Fields      map[string]interface{} `json:"fields"`
    }
    if err := json.Unmarshal(jsonData, &req); err != nil {
        return err
    }

    cmd := command.WriteMetricCommand{
        Measurement: req.Measurement,
        Tags:        req.Tags,
        Fields:      req.Fields,
        Timestamp:   nil,
    }

    return h.writeMetricHandler.Handle(ctx, cmd)
}

func (h *EventHandlerImpl) handleQueryLLM(ctx context.Context, data interface{}) error {
    // Process LLM query and return result via Socket.IO
    // Implementation details in LLM module
    return nil
}
```

---

## 5. โมดูล Redis (Cache & Stream)

Redis ใช้สำหรับ caching, rate limiting, session management, และ message streaming

### 5.1 Redis Client

```go
// internal/adapters/out/redis/client.go
package redis

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

// Client wraps go-redis client
type Client struct {
    rdb *redis.Client
    cfg *Config
}

type Config struct {
    Addr     string        `yaml:"addr"`
    Password string        `yaml:"password"`
    DB       int           `yaml:"db"`
    PoolSize int           `yaml:"pool_size"`
    Timeout  time.Duration `yaml:"timeout"`
}

// NewClient creates a new Redis client
func NewClient(cfg *Config) (*Client, error) {
    rdb := redis.NewClient(&redis.Options{
        Addr:         cfg.Addr,
        Password:     cfg.Password,
        DB:           cfg.DB,
        PoolSize:     cfg.PoolSize,
        DialTimeout:  cfg.Timeout,
        ReadTimeout:  cfg.Timeout,
        WriteTimeout: cfg.Timeout,
    })

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := rdb.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("failed to connect to Redis: %w", err)
    }

    return &Client{
        rdb: rdb,
        cfg: cfg,
    }, nil
}

// Close closes the Redis connection
func (c *Client) Close() error {
    return c.rdb.Close()
}

// --- Cache Operations ---

// Set stores a value in cache with TTL
func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return c.rdb.Set(ctx, key, data, ttl).Err()
}

// Get retrieves a value from cache
func (c *Client) Get(ctx context.Context, key string, dest interface{}) error {
    data, err := c.rdb.Get(ctx, key).Bytes()
    if err != nil {
        return err
    }
    return json.Unmarshal(data, dest)
}

// Delete removes a key from cache
func (c *Client) Delete(ctx context.Context, key string) error {
    return c.rdb.Del(ctx, key).Err()
}

// --- Stream Operations ---

// PublishToStream publishes a message to a Redis stream
func (c *Client) PublishToStream(ctx context.Context, stream string, values map[string]interface{}) error {
    return c.rdb.XAdd(ctx, &redis.XAddArgs{
        Stream: stream,
        Values: values,
    }).Err()
}

// ConsumeStream consumes messages from a Redis stream
func (c *Client) ConsumeStream(ctx context.Context, stream, group, consumer string, count int64) ([]redis.XMessage, error) {
    // Create consumer group if not exists
    if err := c.rdb.XGroupCreateMkStream(ctx, stream, group, "0").Err(); err != nil {
        if err.Error() != "BUSYGROUP Consumer Group name already exists" {
            return nil, err
        }
    }

    return c.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
        Group:    group,
        Consumer: consumer,
        Streams:  []string{stream, ">"},
        Count:    count,
        Block:    0,
    }).Result()
}

// AckStream acknowledges a message in a stream
func (c *Client) AckStream(ctx context.Context, stream, group string, ids []string) error {
    return c.rdb.XAck(ctx, stream, group, ids...).Err()
}

// --- Rate Limiting ---

// RateLimit checks if a key has exceeded rate limit
func (c *Client) RateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
    pipe := c.rdb.Pipeline()
    incr := pipe.Incr(ctx, key)
    pipe.Expire(ctx, key, window)
    _, err := pipe.Exec(ctx)
    if err != nil {
        return false, err
    }
    return incr.Val() > int64(limit), nil
}
```

### 5.2 Redis Cache Repository (Port Implementation)

```go
// internal/adapters/out/redis/cache_repository.go
package redis

import (
    "context"
    "time"

    "your-project/internal/core/domain"
)

// CacheRepository implements domain.CacheRepository
type CacheRepository struct {
    client *Client
}

func NewCacheRepository(client *Client) *CacheRepository {
    return &CacheRepository{client: client}
}

func (r *CacheRepository) SetLLMResponse(ctx context.Context, query string, response *domain.LLMResponse) error {
    key := "llm:" + query
    return r.client.Set(ctx, key, response, 1*time.Hour)
}

func (r *CacheRepository) GetLLMResponse(ctx context.Context, query string) (*domain.LLMResponse, error) {
    key := "llm:" + query
    var resp domain.LLMResponse
    if err := r.client.Get(ctx, key, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

func (r *CacheRepository) SetDeviceStatus(ctx context.Context, deviceID string, status *domain.DeviceStatus) error {
    key := "device:" + deviceID
    return r.client.Set(ctx, key, status, 5*time.Minute)
}

func (r *CacheRepository) GetDeviceStatus(ctx context.Context, deviceID string) (*domain.DeviceStatus, error) {
    key := "device:" + deviceID
    var status domain.DeviceStatus
    if err := r.client.Get(ctx, key, &status); err != nil {
        return nil, err
    }
    return &status, nil
}

func (r *CacheRepository) PublishIoTData(ctx context.Context, data *domain.IoTData) error {
    values := map[string]interface{}{
        "device_id":   data.DeviceID,
        "measurement": data.Measurement,
        "fields":      data.Fields,
        "timestamp":   data.Timestamp.Format(time.RFC3339),
    }
    return r.client.PublishToStream(ctx, "iot:stream", values)
}

func (r *CacheRepository) CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
    return r.client.RateLimit(ctx, key, limit, window)
}
```

---

## 6. โมดูล PostgreSQL (Alarm & History)

PostgreSQL ใช้เก็บข้อมูลเชิง relational เช่น device registry, alarm history, user data

### 6.1 PostgreSQL Client with GORM

```go
// internal/adapters/out/postgres/client.go
package postgres

import (
    "fmt"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// Config holds PostgreSQL configuration
type Config struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    Database string `yaml:"database"`
    SSLMode  string `yaml:"ssl_mode"`
}

// Client wraps GORM
type Client struct {
    db *gorm.DB
}

// NewClient creates a new PostgreSQL client
func NewClient(cfg *Config) (*Client, error) {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
    })
    if err != nil {
        return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
    }

    // Auto-migrate models
    if err := db.AutoMigrate(&Device{}, &AlarmLog{}, &LLMInteraction{}); err != nil {
        return nil, fmt.Errorf("failed to migrate: %w", err)
    }

    return &Client{db: db}, nil
}

// Close closes the database connection
func (c *Client) Close() error {
    sqlDB, err := c.db.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}

// GetDB returns the underlying GORM DB
func (c *Client) GetDB() *gorm.DB {
    return c.db
}
```

### 6.2 GORM Models

```go
// internal/adapters/out/postgres/models/device.go
package models

import (
    "time"

    "gorm.io/gorm"
)

// Device represents an IoT device
type Device struct {
    ID          string         `gorm:"primaryKey;type:varchar(50)" json:"id"`
    Name        string         `gorm:"type:varchar(100)" json:"name"`
    Type        string         `gorm:"type:varchar(50)" json:"type"`
    Location    string         `gorm:"type:varchar(100)" json:"location"`
    Status      string         `gorm:"type:varchar(20);default:'active'" json:"status"`
    Metadata    map[string]interface{} `gorm:"type:jsonb" json:"metadata"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// AlarmLog represents an alarm/alert record
type AlarmLog struct {
    ID          string         `gorm:"primaryKey;type:varchar(50)" json:"id"`
    DeviceID    string         `gorm:"type:varchar(50);index" json:"device_id"`
    Measurement string         `gorm:"type:varchar(100)" json:"measurement"`
    Condition   string         `gorm:"type:text" json:"condition"`
    Severity    string         `gorm:"type:varchar(20);default:'warning'" json:"severity"`
    Message     string         `gorm:"type:text" json:"message"`
    TriggeredAt time.Time      `gorm:"index" json:"triggered_at"`
    ResolvedAt  *time.Time     `json:"resolved_at,omitempty"`
    ResolvedBy  string         `gorm:"type:varchar(50)" json:"resolved_by,omitempty"`
    CreatedAt   time.Time      `json:"created_at"`
}

// LLMInteraction represents an interaction with LLM
type LLMInteraction struct {
    ID          string         `gorm:"primaryKey;type:varchar(50)" json:"id"`
    DeviceID    string         `gorm:"type:varchar(50);index" json:"device_id"`
    Query       string         `gorm:"type:text" json:"query"`
    Response    string         `gorm:"type:text" json:"response"`
    Model       string         `gorm:"type:varchar(50)" json:"model"`
    TokensUsed  int            `json:"tokens_used"`
    Duration    int64          `json:"duration_ms"`
    Success     bool           `json:"success"`
    Error       string         `gorm:"type:text" json:"error,omitempty"`
    CreatedAt   time.Time      `gorm:"index" json:"created_at"`
}
```

### 6.3 Device Repository Implementation

```go
// internal/adapters/out/postgres/device_repository.go
package postgres

import (
    "context"
    "fmt"

    "gorm.io/gorm"

    "your-project/internal/core/domain"
    "your-project/internal/adapters/out/postgres/models"
)

// DeviceRepository implements domain.DeviceRepository
type DeviceRepository struct {
    db *gorm.DB
}

func NewDeviceRepository(client *Client) *DeviceRepository {
    return &DeviceRepository{db: client.GetDB()}
}

func (r *DeviceRepository) Create(ctx context.Context, device *domain.Device) error {
    model := &models.Device{
        ID:       device.ID,
        Name:     device.Name,
        Type:     device.Type,
        Location: device.Location,
        Status:   device.Status,
        Metadata: device.Metadata,
    }
    return r.db.WithContext(ctx).Create(model).Error
}

func (r *DeviceRepository) GetByID(ctx context.Context, id string) (*domain.Device, error) {
    var model models.Device
    if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
        return nil, err
    }
    return &domain.Device{
        ID:       model.ID,
        Name:     model.Name,
        Type:     model.Type,
        Location: model.Location,
        Status:   model.Status,
        Metadata: model.Metadata,
    }, nil
}

func (r *DeviceRepository) Update(ctx context.Context, device *domain.Device) error {
    return r.db.WithContext(ctx).Model(&models.Device{}).
        Where("id = ?", device.ID).
        Updates(map[string]interface{}{
            "name":      device.Name,
            "type":      device.Type,
            "location":  device.Location,
            "status":    device.Status,
            "metadata":  device.Metadata,
        }).Error
}

func (r *DeviceRepository) Delete(ctx context.Context, id string) error {
    return r.db.WithContext(ctx).Delete(&models.Device{}, "id = ?", id).Error
}

func (r *DeviceRepository) List(ctx context.Context, filter *domain.DeviceFilter) ([]*domain.Device, error) {
    query := r.db.WithContext(ctx).Model(&models.Device{})

    if filter.Type != "" {
        query = query.Where("type = ?", filter.Type)
    }
    if filter.Status != "" {
        query = query.Where("status = ?", filter.Status)
    }
    if filter.Limit > 0 {
        query = query.Limit(filter.Limit)
    }
    if filter.Offset > 0 {
        query = query.Offset(filter.Offset)
    }

    var models []models.Device
    if err := query.Find(&models).Error; err != nil {
        return nil, err
    }

    devices := make([]*domain.Device, len(models))
    for i, m := range models {
        devices[i] = &domain.Device{
            ID:       m.ID,
            Name:     m.Name,
            Type:     m.Type,
            Location: m.Location,
            Status:   m.Status,
            Metadata: m.Metadata,
        }
    }
    return devices, nil
}
```

### 6.4 Alarm Repository Implementation

```go
// internal/adapters/out/postgres/alarm_repository.go
package postgres

import (
    "context"
    "time"

    "gorm.io/gorm"

    "your-project/internal/core/domain"
    "your-project/internal/adapters/out/postgres/models"
)

// AlarmRepository implements domain.AlarmRepository
type AlarmRepository struct {
    db *gorm.DB
}

func NewAlarmRepository(client *Client) *AlarmRepository {
    return &AlarmRepository{db: client.GetDB()}
}

func (r *AlarmRepository) Save(ctx context.Context, alarm *domain.Alarm) error {
    model := &models.AlarmLog{
        ID:          alarm.ID,
        DeviceID:    alarm.DeviceID,
        Measurement: alarm.Measurement,
        Condition:   alarm.Condition,
        Severity:    alarm.Severity,
        Message:     alarm.Message,
        TriggeredAt: alarm.TriggeredAt,
        ResolvedAt:  alarm.ResolvedAt,
        ResolvedBy:  alarm.ResolvedBy,
    }
    return r.db.WithContext(ctx).Create(model).Error
}

func (r *AlarmRepository) GetHistory(ctx context.Context, filter *domain.AlarmFilter) ([]*domain.Alarm, error) {
    query := r.db.WithContext(ctx).Model(&models.AlarmLog{})

    if filter.DeviceID != "" {
        query = query.Where("device_id = ?", filter.DeviceID)
    }
    if filter.Severity != "" {
        query = query.Where("severity = ?", filter.Severity)
    }
    if !filter.Start.IsZero() {
        query = query.Where("triggered_at >= ?", filter.Start)
    }
    if !filter.End.IsZero() {
        query = query.Where("triggered_at <= ?", filter.End)
    }
    if filter.Limit > 0 {
        query = query.Limit(filter.Limit)
    }
    query = query.Order("triggered_at DESC")

    var models []models.AlarmLog
    if err := query.Find(&models).Error; err != nil {
        return nil, err
    }

    alarms := make([]*domain.Alarm, len(models))
    for i, m := range models {
        alarms[i] = &domain.Alarm{
            ID:          m.ID,
            DeviceID:    m.DeviceID,
            Measurement: m.Measurement,
            Condition:   m.Condition,
            Severity:    m.Severity,
            Message:     m.Message,
            TriggeredAt: m.TriggeredAt,
            ResolvedAt:  m.ResolvedAt,
            ResolvedBy:  m.ResolvedBy,
        }
    }
    return alarms, nil
}

func (r *AlarmRepository) Resolve(ctx context.Context, id, resolvedBy string) error {
    now := time.Now()
    return r.db.WithContext(ctx).Model(&models.AlarmLog{}).
        Where("id = ?", id).
        Updates(map[string]interface{}{
            "resolved_at": now,
            "resolved_by": resolvedBy,
        }).Error
}
```

---

## 7. โมดูล LLM + AI Agent

LLM module ทำหน้าที่เป็น AI Agent สำหรับวิเคราะห์ข้อมูล IoT, ตรวจจับ anomalies, และตอบคำถาม

### 7.1 LLM Client

```go
// internal/adapters/out/llm/client.go
package llm

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

// Config holds LLM configuration
type Config struct {
    Provider   string        `yaml:"provider"` // openai, anthropic, local
    APIKey     string        `yaml:"api_key"`
    BaseURL    string        `yaml:"base_url"`
    Model      string        `yaml:"model"`
    MaxTokens  int           `yaml:"max_tokens"`
    Temperature float64      `yaml:"temperature"`
    Timeout    time.Duration `yaml:"timeout"`
}

// Client interfaces with LLM providers
type Client struct {
    config *Config
    http   *http.Client
}

// NewClient creates a new LLM client
func NewClient(cfg *Config) *Client {
    return &Client{
        config: cfg,
        http: &http.Client{
            Timeout: cfg.Timeout,
        },
    }
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
    Model    string        `json:"model"`
    Messages []ChatMessage `json:"messages"`
    MaxTokens int          `json:"max_tokens,omitempty"`
    Temperature float64    `json:"temperature,omitempty"`
}

// ChatMessage represents a message in chat
type ChatMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
    ID      string `json:"id"`
    Choices []struct {
        Message ChatMessage `json:"message"`
    } `json:"choices"`
    Usage struct {
        PromptTokens     int `json:"prompt_tokens"`
        CompletionTokens int `json:"completion_tokens"`
        TotalTokens      int `json:"total_tokens"`
    } `json:"usage"`
}

// Chat sends a chat completion request
func (c *Client) Chat(ctx context.Context, messages []ChatMessage) (*ChatResponse, error) {
    reqBody := ChatRequest{
        Model:       c.config.Model,
        Messages:    messages,
        MaxTokens:   c.config.MaxTokens,
        Temperature: c.config.Temperature,
    }

    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/v1/chat/completions", bytes.NewReader(jsonData))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

    resp, err := c.http.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("LLM API error: %s - %s", resp.Status, string(body))
    }

    var result ChatResponse
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    return &result, nil
}
```

### 7.2 AI Agent

```go
// internal/adapters/out/llm/agent.go
package llm

import (
    "context"
    "encoding/json"
    "fmt"
    "log/slog"
    "time"

    "your-project/internal/core/domain"
)

// Agent represents an AI agent for IoT analysis
type Agent struct {
    llmClient *Client
    cache     domain.CacheRepository
    prompt    *PromptEngine
}

// NewAgent creates a new AI agent
func NewAgent(llmClient *Client, cache domain.CacheRepository) *Agent {
    return &Agent{
        llmClient: llmClient,
        cache:     cache,
        prompt:    NewPromptEngine(),
    }
}

// AnalyzeIoTData analyzes IoT data for anomalies and insights
func (a *Agent) AnalyzeIoTData(ctx context.Context, data *domain.IoTData) (*domain.AgentResult, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("analysis:%s:%s", data.DeviceID, data.Measurement)
    var cached domain.AgentResult
    if err := a.cache.GetLLMResponse(ctx, cacheKey); err == nil {
        return &cached, nil
    }

    // Build prompt
    prompt := a.prompt.BuildAnalysisPrompt(data)

    // Call LLM
    messages := []ChatMessage{
        {Role: "system", Content: "You are an IoT data analyst AI agent. Analyze sensor data and provide insights."},
        {Role: "user", Content: prompt},
    }

    resp, err := a.llmClient.Chat(ctx, messages)
    if err != nil {
        return nil, err
    }

    if len(resp.Choices) == 0 {
        return nil, fmt.Errorf("no response from LLM")
    }

    // Parse response
    result := &domain.AgentResult{
        Analysis:   resp.Choices[0].Message.Content,
        TokensUsed: resp.Usage.TotalTokens,
        Timestamp:  time.Now(),
    }

    // Cache result
    _ = a.cache.SetLLMResponse(ctx, cacheKey, &domain.LLMResponse{
        Query:     cacheKey,
        Response:  result.Analysis,
        Tokens:    result.TokensUsed,
        Timestamp: result.Timestamp,
    })

    return result, nil
}

// PredictMaintenance predicts when maintenance is needed
func (a *Agent) PredictMaintenance(ctx context.Context, deviceID string, metrics []*domain.Metric) (*domain.PredictionResult, error) {
    prompt := a.prompt.BuildPredictionPrompt(deviceID, metrics)

    messages := []ChatMessage{
        {Role: "system", Content: "You are a predictive maintenance AI agent. Analyze sensor data and predict maintenance needs."},
        {Role: "user", Content: prompt},
    }

    resp, err := a.llmClient.Chat(ctx, messages)
    if err != nil {
        return nil, err
    }

    if len(resp.Choices) == 0 {
        return nil, fmt.Errorf("no response from LLM")
    }

    return &domain.PredictionResult{
        DeviceID:    deviceID,
        Prediction:  resp.Choices[0].Message.Content,
        Confidence:  0.85, // Placeholder - should be parsed from response
        Timestamp:   time.Now(),
    }, nil
}

// NaturalLanguageQuery handles natural language queries about IoT data
func (a *Agent) NaturalLanguageQuery(ctx context.Context, query string, deviceID string) (*domain.LLMResponse, error) {
    // Check cache
    cacheKey := fmt.Sprintf("nlq:%s:%s", deviceID, query)
    var cached domain.LLMResponse
    if err := a.cache.GetLLMResponse(ctx, cacheKey); err == nil {
        return &cached, nil
    }

    prompt := a.prompt.BuildNLQPrompt(query, deviceID)

    messages := []ChatMessage{
        {Role: "system", Content: "You are an IoT data assistant. Answer questions about IoT devices and sensor data."},
        {Role: "user", Content: prompt},
    }

    resp, err := a.llmClient.Chat(ctx, messages)
    if err != nil {
        return nil, err
    }

    if len(resp.Choices) == 0 {
        return nil, fmt.Errorf("no response from LLM")
    }

    result := &domain.LLMResponse{
        Query:     query,
        Response:  resp.Choices[0].Message.Content,
        Tokens:    resp.Usage.TotalTokens,
        Timestamp: time.Now(),
    }

    // Cache result
    _ = a.cache.SetLLMResponse(ctx, cacheKey, result)

    return result, nil
}
```

### 7.3 Prompt Engine

```go
// internal/adapters/out/llm/prompt.go
package llm

import (
    "fmt"
    "strings"
    "time"

    "your-project/internal/core/domain"
)

// PromptEngine builds prompts for LLM
type PromptEngine struct{}

func NewPromptEngine() *PromptEngine {
    return &PromptEngine{}
}

// BuildAnalysisPrompt builds a prompt for IoT data analysis
func (p *PromptEngine) BuildAnalysisPrompt(data *domain.IoTData) string {
    fields := make([]string, 0, len(data.Fields))
    for k, v := range data.Fields {
        fields = append(fields, fmt.Sprintf("%s: %v", k, v))
    }

    return fmt.Sprintf(`Analyze the following IoT sensor data:

Device: %s
Measurement: %s
Timestamp: %s
Data: %s

Please provide:
1. Any anomalies or outliers detected
2. Insights about the data patterns
3. Recommended actions if any

Be concise and specific.`, data.DeviceID, data.Measurement, data.Timestamp.Format(time.RFC3339), strings.Join(fields, ", "))
}

// BuildPredictionPrompt builds a prompt for predictive maintenance
func (p *PromptEngine) BuildPredictionPrompt(deviceID string, metrics []*domain.Metric) string {
    dataPoints := make([]string, 0, len(metrics))
    for _, m := range metrics {
        for k, v := range m.Fields {
            dataPoints = append(dataPoints, fmt.Sprintf("%s: %v at %s", k, v, m.Timestamp.Format(time.RFC3339)))
        }
    }

    return fmt.Sprintf(`Predict maintenance needs for device %s based on historical data:

%s

Please predict:
1. When maintenance might be needed
2. What type of maintenance
3. Confidence level

Consider trends and patterns in the data.`, deviceID, strings.Join(dataPoints, "\n"))
}

// BuildNLQPrompt builds a prompt for natural language query
func (p *PromptEngine) BuildNLQPrompt(query, deviceID string) string {
    return fmt.Sprintf(`Answer the following question about IoT device data:

Device: %s
Question: %s

Provide a clear, helpful response based on typical IoT sensor data patterns.
If you don't have specific data, provide general guidance.`, deviceID, query)
}
```

---

## 8. การเชื่อมต่อทุกอย่างเข้าด้วยกัน

### 8.1 Process IoT Data Use Case

```go
// internal/core/application/command/process_iot_data.go
package command

import (
    "context"
    "log/slog"
    "time"

    "your-project/internal/core/domain"
)

// ProcessIoTDataCommand represents the command to process IoT data
type ProcessIoTDataCommand struct {
    DeviceID    string
    Measurement string
    Tags        map[string]string
    Fields      map[string]interface{}
    Timestamp   *time.Time
}

// ProcessIoTDataHandler handles IoT data processing
type ProcessIoTDataHandler struct {
    metricRepo   domain.MetricRepository
    cacheRepo    domain.CacheRepository
    deviceRepo   domain.DeviceRepository
    alarmRepo    domain.AlarmRepository
    llmAgent     *llm.Agent
    socketServer *websocket.SocketIOServer
}

func NewProcessIoTDataHandler(
    metricRepo domain.MetricRepository,
    cacheRepo domain.CacheRepository,
    deviceRepo domain.DeviceRepository,
    alarmRepo domain.AlarmRepository,
    llmAgent *llm.Agent,
    socketServer *websocket.SocketIOServer,
) *ProcessIoTDataHandler {
    return &ProcessIoTDataHandler{
        metricRepo:   metricRepo,
        cacheRepo:    cacheRepo,
        deviceRepo:   deviceRepo,
        alarmRepo:    alarmRepo,
        llmAgent:     llmAgent,
        socketServer: socketServer,
    }
}

// Handle processes incoming IoT data
func (h *ProcessIoTDataHandler) Handle(ctx context.Context, cmd ProcessIoTDataCommand) error {
    slog.Info("processing IoT data", "device", cmd.DeviceID, "measurement", cmd.Measurement)

    // 1. Write to InfluxDB
    metric, err := domain.NewMetric(cmd.Measurement, cmd.Tags, cmd.Fields, *cmd.Timestamp)
    if err != nil {
        return err
    }
    if err := h.metricRepo.Write(ctx, metric); err != nil {
        return err
    }

    // 2. Publish to Redis Stream for other consumers
    iotData := &domain.IoTData{
        DeviceID:    cmd.DeviceID,
        Measurement: cmd.Measurement,
        Tags:        cmd.Tags,
        Fields:      cmd.Fields,
        Timestamp:   *cmd.Timestamp,
    }
    if err := h.cacheRepo.PublishIoTData(ctx, iotData); err != nil {
        slog.Error("failed to publish to Redis stream", "error", err)
    }

    // 3. Check for anomalies using AI Agent (async)
    go func() {
        result, err := h.llmAgent.AnalyzeIoTData(context.Background(), iotData)
        if err != nil {
            slog.Error("AI analysis failed", "error", err)
            return
        }

        // Broadcast analysis result via Socket.IO
        if h.socketServer != nil {
            _ = h.socketServer.Broadcast("analysis_result", result)
        }

        // Check if analysis indicates alarm condition
        if strings.Contains(strings.ToLower(result.Analysis), "alert") ||
           strings.Contains(strings.ToLower(result.Analysis), "warning") {
            // Create alarm
            alarm := &domain.Alarm{
                ID:          generateID(),
                DeviceID:    cmd.DeviceID,
                Measurement: cmd.Measurement,
                Condition:   "AI detected anomaly",
                Severity:    "warning",
                Message:     result.Analysis,
                TriggeredAt: time.Now(),
            }
            _ = h.alarmRepo.Save(context.Background(), alarm)
        }
    }()

    // 4. Update device cache
    status := &domain.DeviceStatus{
        DeviceID:    cmd.DeviceID,
        LastReading: *cmd.Timestamp,
        Status:      "active",
    }
    _ = h.cacheRepo.SetDeviceStatus(ctx, cmd.DeviceID, status)

    return nil
}
```

### 8.2 Query LLM Use Case

```go
// internal/core/application/command/query_llm.go
package command

import (
    "context"

    "your-project/internal/core/domain"
    "your-project/internal/adapters/out/llm"
)

// QueryLLMCommand represents a command to query LLM
type QueryLLMCommand struct {
    Query    string
    DeviceID string
    Context  map[string]interface{}
}

// QueryLLMHandler handles LLM queries
type QueryLLMHandler struct {
    llmAgent     *llm.Agent
    cacheRepo    domain.CacheRepository
    socketServer *websocket.SocketIOServer
}

func NewQueryLLMHandler(
    llmAgent *llm.Agent,
    cacheRepo domain.CacheRepository,
    socketServer *websocket.SocketIOServer,
) *QueryLLMHandler {
    return &QueryLLMHandler{
        llmAgent:     llmAgent,
        cacheRepo:    cacheRepo,
        socketServer: socketServer,
    }
}

// Handle processes an LLM query
func (h *QueryLLMHandler) Handle(ctx context.Context, cmd QueryLLMCommand) (*domain.LLMResponse, error) {
    // Check rate limit
    limited, err := h.cacheRepo.CheckRateLimit(ctx, "llm:"+cmd.DeviceID, 10, 1*time.Minute)
    if err != nil {
        return nil, err
    }
    if limited {
        return nil, domain.ErrRateLimitExceeded
    }

    // Query LLM
    resp, err := h.llmAgent.NaturalLanguageQuery(ctx, cmd.Query, cmd.DeviceID)
    if err != nil {
        return nil, err
    }

    // Broadcast response via Socket.IO
    if h.socketServer != nil {
        _ = h.socketServer.Broadcast("llm_response", resp)
    }

    return resp, nil
}
```

---

## 9. Dependency Injection (Wire)

```go
// internal/di/wire.go
//go:build wireinject
// +build wireinject

package di

import (
    "context"

    "github.com/google/wire"

    "your-project/internal/adapters/in/http"
    "your-project/internal/adapters/in/mqtt"
    "your-project/internal/adapters/in/websocket"
    "your-project/internal/adapters/out/influxdb"
    "your-project/internal/adapters/out/postgres"
    "your-project/internal/adapters/out/redis"
    "your-project/internal/adapters/out/llm"
    "your-project/internal/core/application/command"
    "your-project/internal/core/application/query"
    "your-project/internal/core/domain"
)

type App struct {
    HTTPServer      *http.Server
    SocketServer    *websocket.SocketIOServer
    MQTTClient      *mqtt.Client
    InfluxClient    *influxdb.Client
    PostgresClient  *postgres.Client
    RedisClient     *redis.Client
    LLMClient       *llm.Client
}

func InitializeApp(ctx context.Context, config *Config) (*App, error) {
    wire.Build(
        // Infrastructure - InfluxDB
        influxdb.NewClient,
        influxdb.NewMetricRepository,

        // Infrastructure - PostgreSQL
        postgres.NewClient,
        postgres.NewDeviceRepository,
        postgres.NewAlarmRepository,

        // Infrastructure - Redis
        redis.NewClient,
        redis.NewCacheRepository,

        // Infrastructure - LLM
        llm.NewClient,
        llm.NewAgent,

        // Infrastructure - MQTT
        mqtt.NewMessageHandlerImpl,
        mqtt.NewClient,

        // Infrastructure - Socket.IO
        websocket.NewEventHandlerImpl,
        websocket.NewSocketIOServer,

        // Domain - Repository interface bindings
        wire.Bind(new(domain.MetricRepository), new(*influxdb.MetricRepository)),
        wire.Bind(new(domain.DeviceRepository), new(*postgres.DeviceRepository)),
        wire.Bind(new(domain.AlarmRepository), new(*postgres.AlarmRepository)),
        wire.Bind(new(domain.CacheRepository), new(*redis.CacheRepository)),

        // Application - Commands
        command.NewWriteMetricHandler,
        command.NewProcessIoTDataHandler,
        command.NewQueryLLMHandler,

        // Application - Queries
        query.NewQueryMetricHandler,

        // HTTP Handlers
        http.NewMetricHandler,
        http.NewServer,

        // App struct
        wire.Struct(new(App), "*"),
    )
    return &App{}, nil
}
```

---

## 10. การ Deploy ด้วย Docker Compose

```yaml
# deployments/docker-compose.yml
version: '3.8'

services:
  # MQTT Broker
  emqx:
    image: emqx/emqx:latest
    container_name: emqx
    ports:
      - "1883:1883"      # MQTT
      - "8083:8083"      # WebSocket
      - "8084:8084"      # WSS
      - "18083:18083"    # Dashboard
    environment:
      - EMQX_NAME=emqx
      - EMQX_HOST=127.0.0.1
    volumes:
      - emqx-data:/opt/emqx/data
      - emqx-log:/opt/emqx/log
    healthcheck:
      test: ["CMD", "emqx", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # InfluxDB
  influxdb:
    image: influxdb:2.7
    container_name: influxdb
    ports:
      - "8086:8086"
    environment:
      - DOCKER_INFLUXDB_INIT_MODE=setup
      - DOCKER_INFLUXDB_INIT_USERNAME=admin
      - DOCKER_INFLUXDB_INIT_PASSWORD=password
      - DOCKER_INFLUXDB_INIT_ORG=myorg
      - DOCKER_INFLUXDB_INIT_BUCKET=mybucket
      - DOCKER_INFLUXDB_INIT_ADMIN_TOKEN=my-super-secret-token
    volumes:
      - influxdb-data:/var/lib/influxdb2
    healthcheck:
      test: ["CMD", "influx", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # PostgreSQL
  postgres:
    image: postgres:15-alpine
    container_name: postgres
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=appuser
      - POSTGRES_PASSWORD=apppass
      - POSTGRES_DB=icmongolang
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U appuser"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis
  redis:
    image: redis:7-alpine
    container_name: redis
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Application
  app:
    build:
      context: ..
      dockerfile: deployments/Dockerfile
    container_name: icmongolang
    ports:
      - "8080:8080"
      - "8081:8081"  # Socket.IO
    depends_on:
      emqx:
        condition: service_healthy
      influxdb:
        condition: service_healthy
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    environment:
      - MQTT_BROKER=tcp://emqx:1883
      - INFLUXDB_URL=http://influxdb:8086
      - INFLUXDB_TOKEN=my-super-secret-token
      - INFLUXDB_ORG=myorg
      - INFLUXDB_BUCKET=mybucket
      - POSTGRES_HOST=postgres
      - POSTGRES_USER=appuser
      - POSTGRES_PASSWORD=apppass
      - POSTGRES_DB=icmongolang
      - REDIS_ADDR=redis:6379
      - LLM_API_KEY=${LLM_API_KEY}
    volumes:
      - ../configs:/app/configs

  # Prometheus
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

  # Grafana
  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    depends_on:
      - prometheus
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-data:/var/lib/grafana

volumes:
  emqx-data:
  emqx-log:
  influxdb-data:
  postgres-data:
  redis-data:
  grafana-data:
```

---

## 11. สรุป

เอกสารนี้ได้ขยายสถาปัตยกรรม DDD + Clean Architecture เดิมให้รองรับระบบ IoT ที่ครบวงจรด้วยส่วนประกอบเพิ่มเติม:

| ส่วนประกอบ | บทบาท | การใช้งาน |
|-----------|-------|-----------|
| **MQTT** | รับข้อมูลจาก IoT devices | Eclipse Paho Go client  |
| **Socket.IO** | Real-time communication | go-socket.io library  |
| **Redis** | Cache, Stream, Rate Limit | go-redis client |
| **PostgreSQL** | Device registry, Alarm history | GORM + PostgreSQL |
| **LLM + AI Agent** | Anomaly detection, Predictive maintenance | OpenAI/Anthropic API |
| **InfluxDB** | Time-series data storage | Existing module |

### กระแสข้อมูลหลัก

1. **IoT → MQTT → Worker Pool** → รับข้อมูลจากอุปกรณ์
2. **Worker → InfluxDB + Redis Stream** → เก็บข้อมูลและส่งต่อ
3. **Redis Stream → AI Agent** → วิเคราะห์ข้อมูลด้วย LLM
4. **AI Agent → Alarm Engine** → ตรวจจับ anomalies
5. **Alarm → PostgreSQL** → บันทึกประวัติ
6. **Socket.IO → Dashboard** → ส่งข้อมูล real-time

### การทำงานร่วมกับ LLM

- **Anomaly Detection**: AI Agent วิเคราะห์ข้อมูล sensor และแจ้งเตือนเมื่อพบความผิดปกติ
- **Predictive Maintenance**: พยากรณ์ความต้องการบำรุงรักษาจากแนวโน้มข้อมูล
- **Natural Language Query**: ผู้ใช้สามารถถามคำถามเกี่ยวกับอุปกรณ์ IoT ด้วยภาษาธรรมชาติ
- **Automated Response**: AI Agent สามารถตอบสนองอัตโนมัติเมื่อพบเหตุการณ์

---
 
1.C:\github\icmongolang\docs\InfluxDBServiceModuleV2.md
2.C:\github\icmongolang\docs\InfluxDBServiceModule.md
นำมา ปรับใช้กับ ของเดิม  โดยต้องไม่กระทบการทำงาน เดิม
ให้ ทำเอกสาร แผนงานออกมาก่อน
ต้องการ คือ
1.ระบบ Settings
2.ระบบแสดงข้อมูล Read time dashboard monitoring
3.ระบบ Alearm  email sms line discord io device  on dashboard  
4.ระบบ managemet Auto control mamnagermt  สั่ง งาน ตั่งค้า อุปกรณ์
5.ระบบ Loging  แบบทั้วไป และ  Smart Contract Blockchain
6.API  Management
7.ระบบ diagram workflow  เหมือน Node-red

