### 📘 Golang Module Master Template (icmongolang)

> **เอกสารต้นแบบสำหรับให้ AI สร้าง/แก้ไขโปรแกรมภาษา Go ทุกประเภท**  
> ใช้โครงสร้าง **Clean Architecture + DDD** พร้อม Layer ที่ชัดเจน, Cross-cutting Concerns ครบ, และ Pattern ที่นำกลับมาใช้ซ้ำได้
>
> **วิธีใช้:** แทนที่ `{{module_name}}` ด้วยชื่อโมดูลจริง (lowercase, เช่น `pdpa`, `users`, `orders`) และ `{{Entity}}`, `{{Action}}` ตามบริบท

---

## สารบัญ

1. [ปรัชญาและหลักการ](#1-ปรัชญาและหลักการ)
2. [โครงสร้างมาตรฐานของโมดูล](#2-โครงสร้างมาตรฐานของโมดูล)
3. [Layer Templates (ฟอร์มเปล่า)](#3-layer-templates)
   - 3.1 [Domain Layer](#31-domain-layer)
   - 3.2 [Application Layer](#32-application-layer)
   - 3.3 [Infrastructure Layer](#33-infrastructure-layer)
   - 3.4 [Interface Layer](#34-interface-layer)
4. [Cross-cutting Concerns](#4-cross-cutting-concerns)
5. [Database & Migrations](#5-database--migrations)
6. [Bootstrap / Entry Points](#6-bootstrap--entry-points)
7. [Reference Implementation (PDPA)](#7-reference-implementation-pdpa)
8. [Prompt Template สำหรับ AI](#8-prompt-template-สำหรับ-ai)
9. [Checklist ก่อน Commit](#9-checklist-ก่อน-commit)

---

## 1. ปรัชญาและหลักการ

### 1.1 Clean Architecture (Robert C. Martin)

แบ่งโค้ดเป็นชั้นๆ แยก **Business Rules** ออกจาก **Technology**:  
```
┌───────────────────────────────────────────┐
│  Interface Layer   (HTTP, WS, CLI)        │  ← รู้จักโลกภายนอก
├───────────────────────────────────────────┤
│  Infrastructure    (DB, Kafka, ES, Redis) │  ← Adapter
├───────────────────────────────────────────┤
│  Application       (Use Cases)            │  ← Business Flow
├───────────────────────────────────────────┤
│  Domain            (Entities, VOs)        │  ← ❤️ หัวใจ (ไม่รู้จักใคร)
└───────────────────────────────────────────┘
```

### 1.2 Domain-Driven Design

- **Entity** – มี identity (UUID), เปลี่ยนแปลงสถานะได้
- **Value Object** – ไม่มี identity, immutable, เทียบด้วยค่า
- **Aggregate Root** – ประตูเดียวในการเข้าถึง entity ย่อย
- **Repository** – interface อยู่ใน Domain, implementation อยู่ใน Infrastructure
- **Domain Service** – logic ที่ไม่ผูกกับ entity ตัวใดตัวหนึ่ง
- **Domain Error** – error ที่ธุรกิจกำหนด (ไม่ใช่ error ของ framework)

### 1.3 กฎการพึ่งพา (Dependency Rule)

| Layer | รู้จัก (import ได้) | ห้ามรู้จัก |
| :--- | :--- | :--- |
| Domain | stdlib, uuid | ทุก layer อื่น, gorm, gin, sarama |
| Application | Domain | Infrastructure, Interface |
| Infrastructure | Domain, Application | Interface |
| Interface | ทุก layer | – |

### 1.4 เมื่อไหร่ใช้ / ไม่ใช้

**✅ ควรใช้ Clean + DDD เมื่อ:**
- Business logic ซับซ้อน, มีหลาย state transitions
- ต้องสลับ technology (PostgreSQL ↔ MongoDB, Kafka ↔ RabbitMQ)
- ทีม > 2 คน, โมดูล > 5
- ต้องการ test โดยไม่ต้อง spin infra

**❌ ไม่ควรใช้ เมื่อ:**
- CRUD ง่ายๆ ไม่มี business rules
- Prototype / PoC
- Script ขนาดเล็ก

---

## 2. โครงสร้างมาตรฐานของโมดูล

```
internal/modules/{{module_name}}/
│
├── domain/                                    # 🏛️ DOMAIN LAYER (ไม่มี dependency ภายนอก)
│   ├── entity/
│   │   └── {{entity}}.go                     # Aggregate Root + Entity ย่อย
│   ├── value_object/
│   │   ├── {{vo_name}}.go                    # Immutable value objects
│   │   └── status.go                         # Enum-like status
│   ├── repository/
│   │   └── {{entity}}_repository.go          # Interface เท่านั้น
│   ├── service/
│   │   ├── {{domain}}_service.go  # Stateless domain logic
│   │   └── {{port}}_port.go       # Outbound port interface (เช่น Hasher, Blockchain)
│   ├── event/
│   │   └── event.go                          # Domain events
│   └── errors/
│       └── errors.go                         # Sentinel errors
│
├── application/                                # 🎯 APPLICATION LAYER (orchestration)
│   ├── {{verb}}_{{entity}}.go                # UseCase ละไฟล์
│   ├── dto.go                                # Request/Response DTO
│   └── mappers.go                            # Entity ↔ DTO
│
├── infrastructure/                             # 🔧 INFRASTRUCTURE LAYER (adapters)
│   ├── persistence/
│   │   ├── postgres/
│   │   │   ├── {{entity}}_repo_impl.go
│   │   │   └── models.go                     # GORM models (มี prefix)
│   │   └── redis/
│   │       └── {{entity}}_cache.go
│   ├── messaging/
│   │   ├── kafka_producer.go
│   │   └── consumers/
│   │       ├── {{topic}}_consumer.go
│   │       └── consumer_group.go
│   ├── search/elasticsearch/
│   │   └── {{entity}}_indexer.go
│   ├── services/                              # Outbound adapters ภายนอก
│   │   ├── email/smtp.go
│   │   ├── llm/openai.go
│   │   ├── blockchain/ethereum.go
│   │   ├── jwt/jwt_maker.go
│   │   └── hash/bcrypt_hasher.go
│   └── scheduler/
│       └── {{job}}_job.go
│
├── interfaces/                                 # 🌐 INTERFACE LAYER (inbound)
│   ├── http/
│   │   ├── {{entity}}_handler.go
│   │   ├── routes.go
│   │   └── dto.go                            # HTTP-specific DTO (ถ้าแยกจาก application)
│   ├── websocket/
│   │   └── hub.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   ├── logging.go
│   │   └── rate_limit.go
│   └── localization/
│       ├── i18n.go
│       └── locales/{en,th}.toml
│
└── module.go                                   # Composition Root (Wire-up ทั้งหมด)
```

### 2.1 Entry Points (cmd/)

```
cmd/
├── api/main.go              # REST + WebSocket server
├── scheduler/main.go        # Cron jobs
├── migrate/main.go          # DB migrations
└── workers/
    ├── {{topic}}/main.go    # Kafka consumer per topic
    └── ...
```

### 2.2 Shared Package (pkg/)

ใช้ได้เฉพาะ package เหล่านี้:
```
pkg/
├── cryptpass/    ├── jwt/         ├── kafka/       ├── logger/
├── db/           ├── elasticsearch├── mqtt/        ├── llm/
├── emailTemplates/├── helpers/    ├── report/      ├── responses/
├── httpErrors/   ├── influxdb/    ├── secureRandom/├── sendEmail/
├── http-swagger/ ├── transaction/ ├── utils/       ├── vectordb/
└── websocket/
```

---

## 3. Layer Templates

### 3.1 Domain Layer

#### 3.1.1 Entity

**`domain/entity/{{entity}}.go`**
```go
package entity

import (
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/{{module_name}}/domain/errors"
	valueobject "icmongolang/internal/modules/{{module_name}}/domain/value_object"
)

// {{Entity}} – Aggregate Root
type {{Entity}} struct {
	ID        uuid.UUID              `json:"id"`
	// ... fields
	Status    valueobject.{{Status}} `json:"status"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// Constructor – บังคับ invariants ตอนสร้าง
func New{{Entity}}(/* required args */) *{{Entity}} {
	now := time.Now()
	return &{{Entity}}{
		ID:        uuid.New(),
		Status:    valueobject.{{Status}}Active,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Behavior methods – เปลี่ยน state ผ่าน method เท่านั้น
func (e *{{Entity}}) Activate() error {
	if e.Status == valueobject.{{Status}}Active {
		return domainerrors.Err{{Entity}}AlreadyActive
	}
	e.Status = valueobject.{{Status}}Active
	e.UpdatedAt = time.Now()
	return nil
}

// Query methods
func (e *{{Entity}}) IsActive() bool {
	return e.Status == valueobject.{{Status}}Active
}
```

#### 3.1.2 Value Object

**`domain/value_object/{{vo_name}}.go`**
```go
package valueobject

type {{Status}} string

const (
	{{Status}}Active   {{Status}} = "ACTIVE"
	{{Status}}Inactive {{Status}} = "INACTIVE"
)

func (s {{Status}}) IsValid() bool {
	switch s {
	case {{Status}}Active, {{Status}}Inactive:
		return true
	}
	return false
}

func (s {{Status}}) String() string { return string(s) }
```

#### 3.1.3 Repository Interface

**`domain/repository/{{entity}}_repository.go`**
```go
package repository

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/{{module_name}}/domain/entity"
)

type {{Entity}}Repository interface {
	Save(ctx context.Context, e *entity.{{Entity}}) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.{{Entity}}, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.{{Entity}}, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
```

#### 3.1.4 Domain Service / Port

**`domain/service/{{port}}_port.go`** (Outbound Port)
```go
package service

import "context"

// {{Port}} – outbound dependency ที่ Domain ต้องการ
type {{Port}} interface {
	DoSomething(ctx context.Context, input string) (string, error)
}
```

#### 3.1.5 Domain Errors

**`domain/errors/errors.go`**
```go
package domainerrors

import "errors"

var (
	Err{{Entity}}NotFound        = errors.New("{{entity}} not found")
	Err{{Entity}}AlreadyActive   = errors.New("{{entity}} already active")
	ErrInvalid{{Field}}          = errors.New("invalid {{field}}")
	Err{{Action}}NotAllowed      = errors.New("{{action}} not allowed")
)
```

---

### 3.2 Application Layer

#### 3.2.1 Use Case

**`application/{{verb}}_{{entity}}.go`**
```go
package application

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/{{module_name}}/domain/entity"
	domainerrors "icmongolang/internal/modules/{{module_name}}/domain/errors"
	"icmongolang/internal/modules/{{module_name}}/domain/repository"
	"icmongolang/internal/modules/{{module_name}}/domain/service"
)

// {{Verb}}{{Entity}}UseCase
type {{Verb}}{{Entity}}UseCase struct {
	repo         repository.{{Entity}}Repository
	auditRepo    repository.AuditRepository
	port         service.{{Port}}
	// ... dependencies อื่น
}

func New{{Verb}}{{Entity}}UseCase(
	repo repository.{{Entity}}Repository,
	auditRepo repository.AuditRepository,
	port service.{{Port}},
) *{{Verb}}{{Entity}}UseCase {
	return &{{Verb}}{{Entity}}UseCase{
		repo:      repo,
		auditRepo: auditRepo,
		port:      port,
	}
}

// Input DTO – ใช้เฉพาะ field ที่จำเป็น
type {{Verb}}{{Entity}}Input struct {
	UserID    uuid.UUID
	// ... fields
	IPAddress string
	UserAgent string
}

// Execute – point of entry เดียว
func (uc *{{Verb}}{{Entity}}UseCase) Execute(ctx context.Context, input {{Verb}}{{Entity}}Input) error {
	// 1. Validate input
	if input.UserID == uuid.Nil {
		return domainerrors.ErrInvalidUserID
	}

	// 2. Load aggregate
	e, err := uc.repo.FindByID(ctx, input.UserID)
	if err != nil {
		return err
	}

	// 3. Call domain behavior
	if err := e.{{Action}}(); err != nil {
		return err
	}

	// 4. Persist
	if err := uc.repo.Save(ctx, e); err != nil {
		return err
	}

	// 5. Side effects (audit, event, cache)
	audit := entity.NewAuditTrail(&input.UserID, "{{ACTION}}", map[string]interface{}{
		"entity_id": e.ID.String(),
	})
	_ = uc.auditRepo.Save(ctx, audit)

	return nil
}
```

#### 3.2.2 DTO

**`application/dto.go`**
```go
package application

import "time"

// Output DTO – ส่งกลับให้ Interface Layer
type {{Entity}}Response struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
```

---

### 3.3 Infrastructure Layer

#### 3.3.1 GORM Model (มี Prefix)

**`infrastructure/persistence/postgres/models.go`**
```go
package postgres

import (
	"time"

	"github.com/google/uuid"
)

// TableName ใช้ prefix {{module_name}}_ เสมอ
type {{Entity}}Model struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	// ... columns
	Status    string    `gorm:"type:varchar(20);not null;index"`
	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time `gorm:"default:now()"`
}

func ({{Entity}}Model) TableName() string { return "{{module_name}}_{{entity}}s" }
```

#### 3.3.2 Repository Implementation

**`infrastructure/persistence/postgres/{{entity}}_repo_impl.go`**
```go
package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"icmongolang/internal/modules/{{module_name}}/domain/entity"
	domainerrors "icmongolang/internal/modules/{{module_name}}/domain/errors"
)

type {{entity}}RepoImpl struct{ db *gorm.DB }

func New{{Entity}}Repository(db *gorm.DB) *{{entity}}RepoImpl {
	return &{{entity}}RepoImpl{db: db}
}

func (r *{{entity}}RepoImpl) Save(ctx context.Context, e *entity.{{Entity}}) error {
	m := &{{Entity}}Model{
		ID:     e.ID,
		Status: string(e.Status),
		// map fields
	}
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *{{entity}}RepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.{{Entity}}, error) {
	var m {{Entity}}Model
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, domainerrors.Err{{Entity}}NotFound
	}
	if err != nil {
		return nil, err
	}
	return r.toEntity(&m), nil
}

func (r *{{entity}}RepoImpl) toEntity(m *{{Entity}}Model) *entity.{{Entity}} {
	return &entity.{{Entity}}{
		ID:     m.ID,
		Status: valueobject.{{Status}}(m.Status),
	}
}
```

#### 3.3.3 Kafka Producer

**`infrastructure/messaging/kafka_producer.go`**
```go
package messaging

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
)

type Producer interface {
	Publish(ctx context.Context, topic string, key string, payload interface{}) error
}

type kafkaProducer struct{ p sarama.SyncProducer }

func NewKafkaProducer(brokers []string) (Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5

	p, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}
	return &kafkaProducer{p: p}, nil
}

func (k *kafkaProducer) Publish(ctx context.Context, topic, key string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, _, err = k.p.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
		Timestamp: time.Now(),
	})
	return err
}
```

#### 3.3.4 Kafka Consumer

**`infrastructure/messaging/consumers/{{topic}}_consumer.go`**
```go
package consumers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type {{Topic}}Handler interface {
	Handle(ctx context.Context, payload map[string]interface{}) error
}

type {{Topic}}Consumer struct {
	handler {{Topic}}Handler
}

func New{{Topic}}Consumer(h {{Topic}}Handler) *{{Topic}}Consumer {
	return &{{Topic}}Consumer{handler: h}
}

func (c *{{Topic}}Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (c *{{Topic}}Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *{{Topic}}Consumer) ConsumeClaim(s sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var payload map[string]interface{}
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			log.Printf("[{{topic}}] unmarshal error: %v", err)
			continue
		}
		if err := c.handler.Handle(context.Background(), payload); err != nil {
			log.Printf("[{{topic}}] handle error: %v", err)
		}
		s.MarkMessage(msg, "")
	}
	return nil
}
```

#### 3.3.5 Redis Cache

**`infrastructure/persistence/redis/{{entity}}_cache.go`**
```go
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type {{Entity}}Cache interface {
	Set(ctx context.Context, id uuid.UUID, val string) error
	Get(ctx context.Context, id uuid.UUID) (string, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type redis{{Entity}}Cache struct {
	client *redis.Client
	ttl    time.Duration
}

func New{{Entity}}Cache(client *redis.Client, ttl time.Duration) *redis{{Entity}}Cache {
	return &redis{{Entity}}Cache{client: client, ttl: ttl}
}

func (c *redis{{Entity}}Cache) key(id uuid.UUID) string {
	return fmt.Sprintf("{{module_name}}:{{entity}}:%s", id)
}

func (c *redis{{Entity}}Cache) Set(ctx context.Context, id uuid.UUID, val string) error {
	return c.client.Set(ctx, c.key(id), val, c.ttl).Err()
}

func (c *redis{{Entity}}Cache) Get(ctx context.Context, id uuid.UUID) (string, error) {
	return c.client.Get(ctx, c.key(id)).Result()
}

func (c *redis{{Entity}}Cache) Delete(ctx context.Context, id uuid.UUID) error {
	return c.client.Del(ctx, c.key(id)).Err()
}
```

#### 3.3.6 Elasticsearch Indexer

**`infrastructure/search/elasticsearch/{{entity}}_indexer.go`**
```go
package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8"
)

type {{Entity}}Indexer interface {
	Index(ctx context.Context, doc interface{}) error
	Delete(ctx context.Context, id string) error
}

type es{{Entity}}Indexer struct {
	client *elasticsearch.Client
	index  string
}

func New{{Entity}}Indexer(client *elasticsearch.Client, index string) *es{{Entity}}Indexer {
	return &es{{Entity}}Indexer{client: client, index: index}
}

func (i *es{{Entity}}Indexer) Index(ctx context.Context, doc interface{}) error {
	data, _ := json.Marshal(doc)
	_, err := i.client.Index(i.index, bytes.NewReader(data))
	return err
}

func (i *es{{Entity}}Indexer) Delete(ctx context.Context, id string) error {
	_, err := i.client.Delete(i.index, id)
	return err
}
```

---

### 3.4 Interface Layer

#### 3.4.1 HTTP Handler

**`interfaces/http/{{entity}}_handler.go`**
```go
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"icmongolang/internal/modules/{{module_name}}/application"
)

type {{Entity}}Handler struct {
	createUC *application.Create{{Entity}}UseCase
	getUC    *application.Get{{Entity}}UseCase
}

func New{{Entity}}Handler(
	createUC *application.Create{{Entity}}UseCase,
	getUC *application.Get{{Entity}}UseCase,
) *{{Entity}}Handler {
	return &{{Entity}}Handler{createUC: createUC, getUC: getUC}
}

func (h *{{Entity}}Handler) Create(c *gin.Context) {
	uid, ok := c.MustGet("user_id").(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		// bind fields
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := application.Create{{Entity}}Input{
		UserID:    uid,
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}
	if err := h.createUC.Execute(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}
```

#### 3.4.2 Routes

**`interfaces/http/routes.go`**
```go
package http

import "github.com/gin-gonic/gin"

type Handlers struct {
	{{Entity}} *{{Entity}}Handler
}

func RegisterRoutes(r *gin.RouterGroup, h *Handlers, auth gin.HandlerFunc) {
	g := r.Group("/{{module_name}}")
	g.Use(auth)

	g.POST("/{{entity}}", h.{{Entity}}.Create)
	g.GET("/{{entity}}/:id", h.{{Entity}}.Get)
}
```

#### 3.4.3 Middleware

**`interfaces/middleware/auth.go`**
```go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		uid, err := uuid.Parse(claims["sub"].(string))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid sub"})
			return
		}
		c.Set("user_id", uid)
		c.Next()
	}
}
```

---

## 4. Cross-cutting Concerns

### 4.1 WebSocket Hub

**`interfaces/websocket/hub.go`**
```go
package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	UserID string
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
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = true
			h.mu.Unlock()
		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.Send)
			}
			h.mu.Unlock()
		case msg := <-h.broadcast:
			h.mu.RLock()
			for c := range h.clients {
				select {
				case c.Send <- msg:
				default:
					close(c.Send)
					delete(h.clients, c)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) SendToUser(userID string, event interface{}) {
	data, _ := json.Marshal(event)
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		if c.UserID == userID {
			select {
			case c.Send <- data:
			default:
			}
		}
	}
}

func (h *Hub) HandleWS(c *gin.Context, userID string) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("ws upgrade: %v", err)
		return
	}
	client := &Client{Hub: h, Conn: conn, Send: make(chan []byte, 256), UserID: userID}
	h.register <- client
	go client.writePump()
	go client.readPump()
}

func (c *Client) writePump() {
	t := time.NewTicker(30 * time.Second)
	defer func() { t.Stop(); c.Conn.Close() }()
	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-t.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() { c.Hub.unregister <- c; c.Conn.Close() }()
	for {
		if _, _, err := c.Conn.ReadMessage(); err != nil {
			return
		}
	}
}
```

### 4.2 LLM Client (OpenAI-compatible)

**`infrastructure/services/llm/openai.go`**
```go
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type Client interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type openaiClient struct {
	apiKey, model, url string
	http               *http.Client
}

func NewOpenAIClient(apiKey, model, url string) Client {
	return &openaiClient{apiKey: apiKey, model: model, url: url, http: &http.Client{}}
}

func (c *openaiClient) Generate(ctx context.Context, prompt string) (string, error) {
	body := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var out struct {
		Choices []struct {
			Message struct{ Content string } `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", nil
	}
	return out.Choices[0].Message.Content, nil
}
```

### 4.3 Blockchain (EVM)

**`infrastructure/services/blockchain/ethereum.go`**
```go
package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"icmongolang/internal/modules/{{module_name}}/domain/service"
)

type ethClient struct {
	client          *ethclient.Client
	priv            *ecdsa.PrivateKey
	from            common.Address
	contractAddress common.Address
	chainID         *big.Int
	enabled         bool
}

func NewEthereumClient(rpcURL, privHex, contractAddr string, chainID int64, enabled bool) (service.BlockchainService, error) {
	if !enabled {
		return &ethClient{enabled: false}, nil
	}
	c, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	priv, err := crypto.HexToECDSA(privHex)
	if err != nil {
		return nil, err
	}
	pub := priv.Public().(*ecdsa.PublicKey)
	from := crypto.PubkeyToAddress(*pub)
	return &ethClient{
		client:          c,
		priv:            priv,
		from:            from,
		contractAddress: common.HexToAddress(contractAddr),
		chainID:         big.NewInt(chainID),
		enabled:         true,
	}, nil
}

func (e *ethClient) RecordHash(ctx context.Context, data string) (string, error) {
	if !e.enabled {
		return "mock_tx", nil
	}
	nonce, err := e.client.PendingNonceAt(ctx, e.from)
	if err != nil {
		return "", err
	}
	hash := crypto.Keccak256Hash([]byte(data)).Bytes()
	tx := types.NewTransaction(nonce, e.contractAddress, big.NewInt(0), 200000, big.NewInt(25e9), hash)
	signed, err := types.SignTx(tx, types.NewEIP155Signer(e.chainID), e.priv)
	if err != nil {
		return "", err
	}
	if err := e.client.SendTransaction(ctx, signed); err != nil {
		return "", err
	}
	return signed.Hash().Hex(), nil
}

var _ = fmt.Sprintf // กัน unused
```

### 4.4 Scheduler (Cron)

**`infrastructure/scheduler/{{job}}_job.go`**
```go
package scheduler

import (
	"context"
	"log"
	"time"
)

type {{Job}}Job struct {
	// dependencies
}

func New{{Job}}Job(/* deps */) *{{Job}}Job {
	return &{{Job}}Job{}
}

func (j *{{Job}}Job) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	log.Println("[{{job}}] starting")
	if err := j.execute(ctx); err != nil {
		log.Printf("[{{job}}] error: %v", err)
	}
	log.Println("[{{job}}] done")
}

func (j *{{Job}}Job) execute(ctx context.Context) error {
	// implement
	return nil
}
```

---

## 5. Database & Migrations

### 5.1 Naming Convention

| สิ่ง | รูปแบบ | ตัวอย่าง |
| :--- | :--- | :--- |
| Table | `{{module_name}}_<plural>` | `pdpa_consents`, `users_accounts` |
| Column | snake_case | `user_id`, `created_at` |
| PK | `id UUID` | – |
| FK | `<entity>_id` | `user_id` |
| Index | `idx_<table>_<column>` | `idx_pdpa_consents_user_id` |
| Unique | `uq_<table>_<column>` | `uq_pdpa_purposes_code` |
| Migration file | `YYYYMMDD_<module>_<desc>.sql` | `20240101_pdpa_init.sql` |

### 5.2 ทำไมต้องใช้ Prefix `{{module_name}}_`

- ✅ **จัดหมวดหมู่ชัดเจน** – แยกตาราง module นี้ออกจาก module อื่น
- ✅ **Backup / Migrate เฉพาะเจาะจง** – DBA สามารถ backup เฉพาะตารางของ module ได้
- ✅ **ป้องกัน Name Collision** – เช่น `logs` ของระบบ general กับ `pdpa_audit_trails`
- ✅ **เวลา Migrate ฐานข้อมูล** – ง่ายต่อการ DROP ทั้ง module

### 5.3 Migration Template

**`migrations/YYYYMMDD_{{module_name}}_init.sql`**
```sql
-- ============================================================
-- {{module_name}} module — initial schema
-- Prefix: {{module_name}}_ เพื่อแยกจากโมดูลอื่น
-- ============================================================

CREATE TABLE IF NOT EXISTS {{module_name}}_{{entity}}s (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    -- columns...
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_{{module_name}}_{{entity}}s_user_id
    ON {{module_name}}_{{entity}}s (user_id);
CREATE INDEX IF NOT EXISTS idx_{{module_name}}_{{entity}}s_status
    ON {{module_name}}_{{entity}}s (status);
```

---

## 6. Bootstrap / Entry Points

### 6.1 Module Composition Root

**`module.go`** (ใช้ Wire-up ทุกอย่างในโมดูล)
```go
package {{module_name}}

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"icmongolang/internal/modules/{{module_name}}/application"
	"icmongolang/internal/modules/{{module_name}}/domain/service"
	"icmongolang/internal/modules/{{module_name}}/infrastructure/messaging"
	"icmongolang/internal/modules/{{module_name}}/infrastructure/persistence/postgres"
	redisrepo "icmongolang/internal/modules/{{module_name}}/infrastructure/persistence/redis"
	httpiface "icmongolang/internal/modules/{{module_name}}/interfaces/http"
)

type Dependencies struct {
	DB      *gorm.DB
	Redis   *redis.Client
	Producer messaging.Producer
}

func Init(router *gin.RouterGroup, deps Dependencies, auth gin.HandlerFunc) {
	// Repositories
	{{entity}}Repo := postgres.New{{Entity}}Repository(deps.DB)
	auditRepo := postgres.NewAuditRepository(deps.DB)

	// Cache
	cache := redisrepo.New{{Entity}}Cache(deps.Redis, 0)

	// Domain services
	policy := service.New{{Domain}}Service()

	// Use cases
	createUC := application.NewCreate{{Entity}}UseCase({{entity}}Repo, auditRepo, policy)

	// Handlers
	handler := httpiface.New{{Entity}}Handler(createUC)

	// Routes
	httpiface.RegisterRoutes(router, &httpiface.Handlers{
		{{Entity}}: handler,
	}, auth)
}
```

### 6.2 API Entry Point

**`cmd/api/main.go`**
```go
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"icmongolang/internal/modules/{{module_name}}"
	"icmongolang/internal/modules/{{module_name}}/infrastructure/messaging"
	"icmongolang/internal/modules/{{module_name}}/interfaces/middleware"
)

func main() {
	_ = godotenv.Load()

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ADDR")})

	producer, err := messaging.NewKafkaProducer([]string{os.Getenv("KAFKA_BROKERS")})
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	api := r.Group("/api/v1")

	auth := middleware.Auth(os.Getenv("JWT_SECRET"))

	{{module_name}}.Init(api, {{module_name}}.Dependencies{
		DB:       db,
		Redis:    rdb,
		Producer: producer,
	}, auth)

	log.Println("listening :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
```

### 6.3 Worker Entry Point

**`cmd/workers/{{topic}}/main.go`**
```go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"

	"icmongolang/internal/modules/{{module_name}}/infrastructure/messaging/consumers"
)

func main() {
	_ = godotenv.Load()

	brokers := []string{os.Getenv("KAFKA_BROKERS")}
	cfg := sarama.NewConfig()
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup(brokers, "{{module_name}}-{{topic}}-group", cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer group.Close()

	handler := consumers.New{{Topic}}Consumer(/* usecase */)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			if err := group.Consume(ctx, []string{"{{topic}}"}, handler); err != nil {
				log.Printf("consume error: %v", err)
			}
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("shutting down")
}
```

---

## 7. Reference Implementation (PDPA)

> PDPA module ในเวอร์ชันเต็มเป็น **ตัวอย่างอ้างอิง** ของ template นี้  
> โครงสร้างทั้งหมดสอดคล้องกับหัวข้อ 2–6 แล้ว

### 7.1 Mapping ตาราง ↔ โครงสร้าง

| Domain | Application | Infrastructure | Interface |
| :--- | :--- | :--- | :--- |
| `ConsentLog`, `DSARRequest`, `AuditTrail`, `UserAccountStatus` | `RecordConsent`, `RevokeConsent`, `SubmitDSAR`, `ImmediateDeletion`, `AutoDeleteExpiredConsents` | `consent_repo_impl.go`, `dsar_repo_impl.go`, `consent_cache.go`, `kafka_producer.go`, `consent_indexer.go`, `consent_cleanup_job.go` | `consent_handler.go`, `dsar_handler.go`, `deletion_handler.go`, `routes.go`, `hub.go`, `i18n.go` |
| VO: `ConsentPurpose`, `ConsentStatus`, `DSARType`, `DSARStatus`, `AccountStatus` | Ports: `BlockchainService`, `LocalizationService`, `EmailSender`, `LLMClient` | Scheduler: `ConsentCleanupJob` | Middleware: Auth, CORS, Rate limit |

### 7.2 PDPA Table Schema (ตัวอย่าง)

```sql
-- ใช้ prefix pdpa_ ทุกตาราง
pdpa_policies              -- เวอร์ชันนโยบาย
pdpa_purposes              -- วัตถุประสงค์ (NECESSARY, ANALYTICS, MARKETING)
pdpa_consents              -- บันทึกความยินยอม
pdpa_user_requests         -- DSAR
pdpa_user_account_statuses -- สถานะบัญชี
pdpa_audit_trails          -- Audit log
pdpa_request_responses     -- ไฟล์ตอบกลับ
```

### 7.3 PDPA Business Flows

| Flow | Use Case | Side Effects |
| :--- | :--- | :--- |
| ผู้ใช้ให้ความยินยอม | `RecordConsent` | DB, Redis, Kafka, Audit, Blockchain |
| ผู้ใช้เพิกถอน | `RevokeConsent` | DB, Redis, Kafka, Audit |
| ผู้ใช้ยื่น DSAR | `SubmitDSAR` | DB, OTP, Kafka (`pdpa.dsar.request`), Audit |
| ระบบประมวลผล DSAR (async) | `ProcessDSAR` | DB, WebSocket broadcast, Email, LLM |
| ลบทันที (หลังยืนยัน) | `ImmediateDeletion` | DB hard delete consent/DSAR, Anonymize user, Redis, ES, Blockchain |
| ลบอัตโนมัติ (Scheduler) | `AutoDeleteExpiredConsents` | เหมือนด้านบน + Kafka `pdpa.data.deleted` |
| รับ event บัญชี | `HandleAccountEvent` | DB update status, Audit |

---

## 8. Prompt Template สำหรับ AI

ใช้ prompt นี้เมื่อสั่งให้ AI สร้างหรือแก้ไข module ในโปรเจกต์ `icmongolang`

````
สร้าง/แก้ไขโมดูลในโปรเจกต์ icmongolang ตาม Template_Module.md

## ข้อมูลนำเข้า
- module_name: <ชื่อโมดูล lowercase เช่น pdpa, users, orders>
- entity: <ชื่อ entity หลัก เช่น Consent, User, Order>
- actions: <รายการ use cases ที่ต้องการ เช่น Create, Get, Delete>
- มี cross-cutting concerns: <Kafka/ES/Redis/LLM/Blockchain/WebSocket/Email – เลือกเฉพาะที่ใช้>
- database prefix: <module_name>_ (default = module_name)
- retention: <ถ้ามี>

## กฎการสร้าง
1. โครงสร้างตามหัวข้อ 2 ของ Template_Module.md
2. Domain Layer:
   - Entity มี constructor + behavior methods (ห้าม set field ตรง)
   - Value Object เป็น type + const + IsValid()
   - Repository เป็น interface
   - Errors เป็น sentinel errors
3. Application Layer:
   - 1 use case ต่อ 1 ไฟล์, มี Input/Output DTO
   - Execute() เป็น entry เดียว, return error
   - เรียก domain behavior, persist, side effects
4. Infrastructure Layer:
   - GORM models ใช้ TableName() + prefix {{module_name}}_
   - Repository impl map model ↔ entity
   - Kafka producer/consumer ใช้ interface
5. Interface Layer:
   - Gin handler ดึง user_id จาก context
   - Route ลงทะเบียนใน routes.go
   - Middleware auth ตรวจ JWT
6. Import package ต้องมาจาก pkg/ หรือ module ภายในเท่านั้น
   ห้ามใช้ package นอกเหนือจาก: pkg/{cryptpass, db, elasticsearch, emailTemplates,
   helpers, http-swagger, httpErrors, influxdb, jwt, kafka, llm, logger, mqtt,
   report, responses, secureRandom, sendEmail, transaction, utils, vectordb, websocket}
7. ถ้า module ใหม่:
   - สร้าง migration YYYYMMDD_<module>_init.sql
   - Wire-up ใน module.go + cmd/api/main.go
8. ถ้า module เดิม:
   - แก้ไขเฉพาะไฟล์ในโฟลเดอร์ของ module นั้น
9. รัน `go build ./...` และ `go test ./...` ให้ผ่าน
10. แสดงเฉพาะไฟล์ที่สร้าง/แก้ + อธิบายสั้นๆ ต่อไฟล์

## Output ที่ต้องการ
1. Tree ของไฟล์ทั้งหมด
2. โค้ดแต่ละไฟล์ (ครบ ไม่ตัด)
3. Migration SQL
4. ตัวอย่าง .env
5. วิธีรัน
````

### 8.1 Prompt ย่อย สำหรับงานเฉพาะทาง

**เพิ่ม Kafka Consumer ใหม่:**
```
เพิ่ม Kafka consumer ใน module {{module_name}}
- topic: {{topic}}
- payload: {{payload schema}}
- business logic: {{อธิบาย}}
- สร้างที่ infrastructure/messaging/consumers/{{topic}}_consumer.go
- สร้าง cmd/workers/{{topic}}/main.go
```

**เพิ่ม Use Case:**
```
เพิ่ม use case {{Verb}}{{Entity}} ใน module {{module_name}}
- Input: {{fields}}
- Business rules: {{อธิบาย}}
- Side effects: {{audit/kafka/cache/blockchain/ws}}
- สร้างที่ application/{{verb}}_{{entity}}.go
- อัปเดต module.go และ handler
```

**เพิ่ม Value Object:**
```
เพิ่ม value object {{VO}} ใน module {{module_name}}
- ค่าที่เป็นไปได้: {{list}}
- Validation rules: {{อธิบาย}}
```

---

## 9. Checklist ก่อน Commit

### 9.1 Domain Layer
- [ ] Entity ทุกตัวมี constructor (`New{{Entity}}`)
- [ ] Entity ไม่มี setter ตรง – เปลี่ยน state ผ่าน behavior method
- [ ] Value Object ทุกตัวมี `IsValid()` หรือ validation method
- [ ] Repository เป็น interface เท่านั้น
- [ ] Domain error เป็น sentinel error
- [ ] Domain **ไม่ import** gorm, gin, sarama, redis

### 9.2 Application Layer
- [ ] Use case ละ 1 ไฟล์, ชื่อ `{{verb}}_{{entity}}.go`
- [ ] Input/Output DTO แยกชัดเจน
- [ ] Execute() return error, ไม่ panic
- [ ] ไม่มี SQL/HTTP ใน use case

### 9.3 Infrastructure Layer
- [ ] GORM model มี `TableName()` + prefix
- [ ] Repository impl map model ↔ entity ถูกต้อง
- [ ] Kafka message ใช้ JSON
- [ ] Redis key มี namespace (`{{module_name}}:{{entity}}:<id>`)
- [ ] ES index มีชื่อสอดคล้อง

### 9.4 Interface Layer
- [ ] Handler ดึง `user_id` จาก context เท่านั้น
- [ ] Route ลงทะเบียนครบ
- [ ] Error response เป็น JSON ที่สอดคล้องกัน
- [ ] Rate limit / auth middleware ถูก apply

### 9.5 Build & Test
- [ ] `go build ./...` ผ่าน
- [ ] `go vet ./...` ผ่าน
- [ ] `go test ./...` ผ่าน
- [ ] ไม่มี import จาก package ที่ไม่อยู่ใน whitelist

### 9.6 Migration
- [ ] ไฟล์ชื่อ `YYYYMMDD_<module>_<desc>.sql`
- [ ] ทุกตารางมี prefix
- [ ] มี index สำหรับ query ที่ใช้บ่อย
- [ ] มี FK constraint

### 9.7 Documentation
- [ ] อัปเดต README ของ module
- [ ] เพิ่มตัวอย่าง .env ถ้ามี env ใหม่
- [ ] อัปเดต docker-compose ถ้ามี service ใหม่

---

## ภาคผนวก

### A. Dependency Matrix

| Layer | stdlib | uuid | gorm | gin | sarama | redis | es | jwt |
| :--- | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: |
| Domain | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Application | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Infrastructure | ✅ | ✅ | ✅ | ❌ | ✅ | ✅ | ✅ | ✅ |
| Interface | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ | ❌ | ✅ |

### B. Kafka Topic Convention

```
<module>.<entity>.<action>
เช่น:
  pdpa.consent.log        # consent granted
  pdpa.consent.revoked    # consent revoked
  pdpa.dsar.request       # DSAR submitted
  pdpa.data.deleted       # data deleted
  pdpa.account.event      # account status changed
  pdpa.email.send         # email queue
  pdpa.llm.analyze        # LLM analysis request
```

### C. Error Code Convention (HTTP)

| Domain Error | HTTP Status |
| :--- | :-: |
| `ErrNotFound` | 404 |
| `ErrAlreadyExists` | 409 |
| `ErrInvalid*` | 400 |
| `ErrUnauthorized` | 401 |
| `ErrForbidden` | 403 |
| `ErrRateLimit` | 429 |

### D. Env Naming

```env
# Database
DB_DSN=
DB_MAX_OPEN=25
DB_MAX_IDLE=10

# Redis
REDIS_ADDR=
REDIS_PASSWORD=
REDIS_DB=0

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_GROUP_ID={{module_name}}-group

# Elasticsearch
ELASTICSEARCH_URL=
ELASTICSEARCH_INDEX={{module_name}}_*

# JWT
JWT_SECRET=
JWT_TTL=24h

# LLM
LLM_API_KEY=
LLM_MODEL=gpt-4
LLM_API_URL=

# Blockchain
BLOCKCHAIN_ENABLED=false
BLOCKCHAIN_RPC_URL=
BLOCKCHAIN_CHAIN_ID=5
BLOCKCHAIN_PRIVATE_KEY=
BLOCKCHAIN_CONTRACT_ADDRESS=

# Module-specific
{{MODULE}}_RETENTION_YEARS=1
{{MODULE}}_DEFAULT_LANG=en
```

---

> **สรุป:** เอกสารนี้เป็น **master template** ที่รวบรวม pattern ทั้งหมดที่ใช้ในโปรเจกต์ `icmongolang` ไว้ในที่เดียว  
> ใช้แทนที่ `{{module_name}}` / `{{Entity}}` / `{{Action}}` เพื่อสร้าง module ใหม่  
> ใช้ section 8 เป็น prompt ให้ AI สร้างโค้ดตาม pattern ที่กำหนด  
> ใช้ section 9 เป็น checklist ก่อน merge
