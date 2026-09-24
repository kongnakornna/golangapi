
---

# 🔧 PART 7 — INFRASTRUCTURE LAYER DEEP DIVE

> **ขนาด**: ใหญ่มาก — แยก 9 ตอนย่อย
> **Part 7A**: Infrastructure Architecture
> **Part 7B**: PostgreSQL Repositories
> **Part 7C**: Redis Cache & Rate Limit
> **Part 7D**: Kafka Producers & Consumers
> **Part 7E**: MQTT Subscriber & Publisher
> **Part 7F**: InfluxDB Repositories
> **Part 7G**: Elasticsearch Indexers
> **Part 7H**: External Services (LLM, Payment, Email)
> **Part 7I**: Schedulers & Background Jobs

---

## 🅰️ PART 7A — INFRASTRUCTURE ARCHITECTURE

### A.1 Position

```
┌─────────────────────────────────────────────────────────────┐
│                      Domain Layer                           │
│  Repos (interface) · Ports (EventProducer, Cache, ...)      │
├─────────────────────────────────────────────────────────────┤
│  ★★★ INFRASTRUCTURE LAYER ★★★  ← PART 7 นี้                │
│  Implement interfaces จาก Domain + Application              │
│  ────────────────────────────────────────────────────────   │
│  ✓ Postgres · Redis · Kafka · InfluxDB · ES · MQTT · LLM   │
│  ✓ Implement repository interfaces                          │
│  ✓ Map model ↔ entity                                       │
│  ✓ Handle external service integration                      │
├─────────────────────────────────────────────────────────────┤
│                    Interface Layer                          │
│  HTTP handlers · MQTT gateway · Workers                     │
└─────────────────────────────────────────────────────────────┘
```

### A.2 Infrastructure Sub-Folders

```
internal/modules/{{module}}/infrastructure/
├── persistence/
│   ├── postgres/
│   │   ├── {{aggregate}}_repo_impl.go
│   │   ├── read_model_repo_impl.go
│   │   ├── models.go              # GORM models (prefix)
│   │   └── mappers.go             # model ↔ entity
│   ├── redis/
│   │   ├── {{aggregate}}_cache.go
│   │   └── keys.go                # Redis key namespace
│   └── influxdb/
│       └── {{metric}}_repo.go
├── messaging/
│   ├── kafka_producer.go
│   └── consumers/
│       ├── {{topic}}_consumer.go
│       └── consumer_group.go
├── search/elasticsearch/
│   └── {{aggregate}}_indexer.go
├── iot/
│   ├── mqtt_subscriber.go
│   ├── mqtt_publisher.go
│   └── protocol/
│       ├── modbus.go
│       ├── lorawan.go
│       └── zigbee.go
├── services/
│   ├── llm/ollama.go
│   ├── email/smtp.go
│   ├── sms/twilio.go
│   ├── payment/stripe.go
│   └── blockchain/ethereum.go
└── scheduler/
    └── {{job}}_job.go
```

### A.3 Rule: Infrastructure Imports

```
✅ Infrastructure imports:
   - domain (entities, repos, events, errors)
   - application (ports, use cases)
   - gorm, sarama, redis, mqtt, influx, elastic
   - stdlib, uuid, decimal

❌ Infrastructure NEVER imports:
   - github.com/gin-gonic/gin (ยกเว้น interfaces layer)
   - handler / routes (แยกชั้น)
```

---

## 🅱️ PART 7B — POSTGRESQL REPOSITORIES

### B.1 GORM Model Convention

```go
// internal/modules/device/infrastructure/persistence/postgres/models.go
package postgres

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/datatypes"
)

// DeviceModel — GORM model พร้อม prefix
type DeviceModel struct {
    ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID   uuid.UUID      `gorm:"type:uuid;not null;index"`
    CustomerID uuid.UUID      `gorm:"type:uuid;not null;index"`
    SiteID     uuid.UUID      `gorm:"type:uuid;not null;index"`
    ZoneID     uuid.UUID      `gorm:"type:uuid;index"`
    ModelID    uuid.UUID      `gorm:"type:uuid;index"`

    SerialNo     string `gorm:"type:varchar(100);uniqueIndex;not null"`
    Name         string `gorm:"type:varchar(255);not null"`
    Type         string `gorm:"type:varchar(50);not null;index"`
    Protocol     string `gorm:"type:varchar(50);not null"`
    Status       string `gorm:"type:varchar(30);not null;index;default:'OFFLINE'"`

    LastSeenAt   *time.Time
    OfflineReason string `gorm:"type:varchar(100)"`
    FaultCode    string `gorm:"type:varchar(50)"`
    FaultMessage string `gorm:"type:text"`

    MQTTClientID    string `gorm:"type:varchar(255);index"`
    DeviceTokenHash string `gorm:"type:varchar(255)"`
    FirmwareVersion string `gorm:"type:varchar(50)"`
    FirmwareChannel string `gorm:"type:varchar(30);default:'stable'"`
    ProvisionedAt   *time.Time

    Capabilities datatypes.JSON `gorm:"type:jsonb"`
    Config       datatypes.JSON `gorm:"type:jsonb"`
    Tags         datatypes.JSON `gorm:"type:jsonb"`
    Metadata     datatypes.JSON `gorm:"type:jsonb"`

    CreatedBy uuid.UUID `gorm:"type:uuid;not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (DeviceModel) TableName() string { return "device_devices" }

// ---- Indexes via migration ----
// CREATE INDEX idx_device_devices_tenant_status ON device_devices(tenant_id, status);
// CREATE INDEX idx_device_devices_last_seen ON device_devices(last_seen_at);
```

### B.2 Repository Impl Pattern

```go
// internal/modules/device/infrastructure/persistence/postgres/device_repo_impl.go
package postgres

import (
    "context"
    "errors"
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"

    "icmongolang/internal/modules/device/domain/entity"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
    "icmongolang/internal/modules/device/domain/repository"
    vo "icmongolang/internal/modules/device/domain/value_object"
    sharedpostgres "icmongolang/internal/shared/infrastructure/postgres"
)

type deviceRepoImpl struct {
    db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) repository.DeviceRepository {
    return &deviceRepoImpl{db: db}
}

// ============================================================
// CREATE
// ============================================================
func (r *deviceRepoImpl) Save(ctx context.Context, d *entity.Device) error {
    db := r.txOrDB(ctx)
    m := toDeviceModel(d)
    return db.WithContext(ctx).Create(m).Error
}

// ============================================================
// READ
// ============================================================
func (r *deviceRepoImpl) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Device, error) {
    var m DeviceModel
    err := r.txOrDB(ctx).WithContext(ctx).
        Where("tenant_id = ? AND id = ?", tenantID, id).
        First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrDeviceNotFound
    }
    if err != nil {
        return nil, err
    }
    return toDeviceEntity(&m), nil
}

func (r *deviceRepoImpl) FindBySerial(ctx context.Context, tenantID uuid.UUID, serial vo.DeviceSerial) (*entity.Device, error) {
    var m DeviceModel
    err := r.txOrDB(ctx).WithContext(ctx).
        Where("tenant_id = ? AND serial_no = ?", tenantID, serial.String()).
        First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrDeviceNotFound
    }
    if err != nil {
        return nil, err
    }
    return toDeviceEntity(&m), nil
}

func (r *deviceRepoImpl) FindBySite(ctx context.Context, tenantID, siteID uuid.UUID) ([]*entity.Device, error) {
    var models []DeviceModel
    err := r.txOrDB(ctx).WithContext(ctx).
        Where("tenant_id = ? AND site_id = ?", tenantID, siteID).
        Order("created_at DESC").
        Find(&models).Error
    if err != nil {
        return nil, err
    }
    return toDeviceEntities(models), nil
}

func (r *deviceRepoImpl) FindByStatus(ctx context.Context, tenantID uuid.UUID, status vo.DeviceStatus) ([]*entity.Device, error) {
    var models []DeviceModel
    err := r.txOrDB(ctx).WithContext(ctx).
        Where("tenant_id = ? AND status = ?", tenantID, string(status)).
        Find(&models).Error
    if err != nil {
        return nil, err
    }
    return toDeviceEntities(models), nil
}

func (r *deviceRepoImpl) FindByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) ([]*entity.Device, error) {
    var models []DeviceModel
    err := r.txOrDB(ctx).WithContext(ctx).
        Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
        Find(&models).Error
    if err != nil {
        return nil, err
    }
    return toDeviceEntities(models), nil
}

// ============================================================
// UPDATE — optimistic locking via updated_at
// ============================================================
func (r *deviceRepoImpl) Update(ctx context.Context, d *entity.Device) error {
    db := r.txOrDB(ctx)
    now := time.Now().UTC()

    result := db.WithContext(ctx).
        Model(&DeviceModel{}).
        Where("tenant_id = ? AND id = ?", d.TenantID, d.ID).
        Updates(map[string]interface{}{
            "name":             d.Name,
            "status":           string(d.Status),
            "firmware_version": d.FirmwareVersion,
            "firmware_channel": d.FirmwareChannel,
            "last_seen_at":     d.LastSeenAt,
            "offline_reason":   d.OfflineReason,
            "fault_code":       d.FaultCode,
            "fault_message":    d.FaultMessage,
            "mqtt_client_id":   d.MQTTClientID,
            "device_token_hash": d.DeviceTokenHash,
            "provisioned_at":   d.ProvisionedAt,
            "config":           d.Config,
            "tags":             d.Tags,
            "updated_at":       now,
        })
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return domainerrors.ErrDeviceNotFound
    }
    return nil
}

// ============================================================
// DELETE
// ============================================================
func (r *deviceRepoImpl) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
    return r.txOrDB(ctx).WithContext(ctx).
        Where("tenant_id = ? AND id = ?", tenantID, id).
        Delete(&DeviceModel{}).Error
}

// ============================================================
// EXISTENCE
// ============================================================
func (r *deviceRepoImpl) ExistsBySerial(ctx context.Context, tenantID uuid.UUID, serial vo.DeviceSerial) (bool, error) {
    var count int64
    err := r.txOrDB(ctx).WithContext(ctx).
        Model(&DeviceModel{}).
        Where("tenant_id = ? AND serial_no = ?", tenantID, serial.String()).
        Count(&count).Error
    return count > 0, err
}

// ============================================================
// BULK
// ============================================================
func (r *deviceRepoImpl) FindStale(ctx context.Context, tenantID uuid.UUID, threshold time.Duration, limit int) ([]*entity.Device, error) {
    cutoff := time.Now().UTC().Add(-threshold)
    var models []DeviceModel
    err := r.txOrDB(ctx).WithContext(ctx).
        Where("tenant_id = ? AND status = ? AND last_seen_at < ?",
            tenantID, string(vo.DeviceStatusOnline), cutoff).
        Limit(limit).
        Find(&models).Error
    if err != nil {
        return nil, err
    }
    return toDeviceEntities(models), nil
}

func (r *deviceRepoImpl) UpdateBatch(ctx context.Context, devices []*entity.Device) error {
    if len(devices) == 0 {
        return nil
    }
    return r.txOrDB(ctx).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        for _, d := range devices {
            result := tx.Model(&DeviceModel{}).
                Where("tenant_id = ? AND id = ?", d.TenantID, d.ID).
                Updates(map[string]interface{}{
                    "status":         string(d.Status),
                    "offline_reason": d.OfflineReason,
                    "updated_at":     time.Now().UTC(),
                })
            if result.Error != nil {
                return result.Error
            }
        }
        return nil
    })
}

// ============================================================
// Helpers
// ============================================================
func (r *deviceRepoImpl) txOrDB(ctx context.Context) *gorm.DB {
    if tx := sharedpostgres.TxFromContext(ctx); tx != nil {
        return tx
    }
    return r.db
}
```

### B.3 Mappers (Model ↔ Entity)

```go
// internal/modules/device/infrastructure/persistence/postgres/mappers.go
package postgres

import (
    "encoding/json"
    "icmongolang/internal/modules/device/domain/entity"
    vo "icmongolang/internal/modules/device/domain/value_object"
)

func toDeviceModel(d *entity.Device) *DeviceModel {
    caps, _ := json.Marshal(d.Capabilities)
    cfg, _ := json.Marshal(d.Config)
    tags, _ := json.Marshal(d.Tags)
    meta, _ := json.Marshal(d.Metadata)

    return &DeviceModel{
        ID:              d.ID,
        TenantID:        d.TenantID,
        CustomerID:      d.CustomerID,
        SiteID:          d.SiteID,
        ZoneID:          d.ZoneID,
        ModelID:         d.ModelID,
        SerialNo:        d.SerialNo.String(),
        Name:            d.Name,
        Type:            string(d.Type),
        Protocol:        string(d.Protocol),
        Status:          string(d.Status),
        LastSeenAt:      d.LastSeenAt,
        OfflineReason:   d.OfflineReason,
        FaultCode:       d.FaultCode,
        FaultMessage:    d.FaultMessage,
        MQTTClientID:    d.MQTTClientID,
        DeviceTokenHash: d.DeviceTokenHash,
        FirmwareVersion: d.FirmwareVersion,
        FirmwareChannel: d.FirmwareChannel,
        ProvisionedAt:   d.ProvisionedAt,
        Capabilities:    caps,
        Config:          cfg,
        Tags:            tags,
        Metadata:        meta,
        CreatedBy:       d.CreatedBy,
        CreatedAt:       d.CreatedAt,
        UpdatedAt:       d.UpdatedAt,
    }
}

func toDeviceEntity(m *DeviceModel) *entity.Device {
    var caps []vo.DeviceCapability
    var cfg map[string]interface{}
    var tags []string
    var meta map[string]interface{}
    _ = json.Unmarshal(m.Capabilities, &caps)
    _ = json.Unmarshal(m.Config, &cfg)
    _ = json.Unmarshal(m.Tags, &tags)
    _ = json.Unmarshal(m.Metadata, &meta)

    return &entity.Device{
        ID:              m.ID,
        TenantID:        m.TenantID,
        CustomerID:      m.CustomerID,
        SiteID:          m.SiteID,
        ZoneID:          m.ZoneID,
        ModelID:         m.ModelID,
        SerialNo:        vo.DeviceSerial(m.SerialNo),
        Name:            m.Name,
        Type:            vo.DeviceType(m.Type),
        Protocol:        vo.Protocol(m.Protocol),
        Status:          vo.DeviceStatus(m.Status),
        LastSeenAt:      m.LastSeenAt,
        OfflineReason:   m.OfflineReason,
        FaultCode:       m.FaultCode,
        FaultMessage:    m.FaultMessage,
        MQTTClientID:    m.MQTTClientID,
        DeviceTokenHash: m.DeviceTokenHash,
        FirmwareVersion: m.FirmwareVersion,
        FirmwareChannel: m.FirmwareChannel,
        ProvisionedAt:   m.ProvisionedAt,
        Capabilities:    caps,
        Config:          cfg,
        Tags:            tags,
        Metadata:        meta,
        CreatedBy:       m.CreatedBy,
        CreatedAt:       m.CreatedAt,
        UpdatedAt:       m.UpdatedAt,
    }
}

func toDeviceEntities(models []DeviceModel) []*entity.Device {
    out := make([]*entity.Device, 0, len(models))
    for i := range models {
        out = append(out, toDeviceEntity(&models[i]))
    }
    return out
}
```

### B.4 Test (Repository)

```go
// internal/modules/device/infrastructure/persistence/postgres/device_repo_test.go
//go:build integration
package postgres_test

func TestDeviceRepository_SaveAndFind(t *testing.T) {
    db := setupTestDB(t)
    repo := postgres.NewDeviceRepository(db)

    device, err := entity.NewDevice(
        uuid.New(), uuid.New(), uuid.New(), uuid.Nil,
        "SN-TEST-001", vo.DeviceTypeSensor, vo.ProtocolMQTT, "Test Device", uuid.New(),
    )
    require.NoError(t, err)

    // Save
    require.NoError(t, repo.Save(context.Background(), device))

    // Find
    found, err := repo.FindByID(context.Background(), device.TenantID, device.ID)
    require.NoError(t, err)
    assert.Equal(t, device.ID, found.ID)
    assert.Equal(t, device.SerialNo, found.SerialNo)
}
```

---

## 🅲 PART 7C — REDIS CACHE & RATE LIMIT

### C.1 Redis Client Wrapper

```go
// internal/shared/infrastructure/redis/client.go
package redis

import (
    "context"
    "time"
    "github.com/redis/go-redis/v9"
)

type Client struct {
    rdb *redis.Client
}

func NewClient(addr, password string, db int) *Client {
    rdb := redis.NewClient(&redis.Options{
        Addr:         addr,
        Password:     password,
        DB:           db,
        PoolSize:     100,
        MinIdleConns: 10,
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
    })
    return &Client{rdb: rdb}
}

func (c *Client) Ping(ctx context.Context) error {
    return c.rdb.Ping(ctx).Err()
}

func (c *Client) Close() error { return c.rdb.Close() }

func (c *Client) Raw() *redis.Client { return c.rdb }
```

### C.2 Cache Adapter (implement `application.Cache`)

```go
// internal/shared/infrastructure/redis/cache.go
package redis

import (
    "context"
    "errors"
    "time"

    "github.com/redis/go-redis/v9"
    "icmongolang/internal/shared/application"
)

type Cache struct {
    client *Client
}

func NewCache(client *Client) application.Cache {
    return &Cache{client: client}
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
    val, err := c.client.Raw().Get(ctx, key).Bytes()
    if errors.Is(err, redis.Nil) {
        return nil, nil // not found
    }
    return val, err
}

func (c *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
    return c.client.Raw().Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Delete(ctx context.Context, key string) error {
    return c.client.Raw().Del(ctx, key).Err()
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
    n, err := c.client.Raw().Exists(ctx, key).Result()
    return n > 0, err
}

// DeleteByPattern — ระวัง: ใช้ SCAN ไม่ใช่ KEYS
func (c *Cache) DeleteByPattern(ctx context.Context, pattern string) (int64, error) {
    var cursor uint64
    var deleted int64
    for {
        keys, next, err := c.client.Raw().Scan(ctx, cursor, pattern, 100).Result()
        if err != nil {
            return deleted, err
        }
        if len(keys) > 0 {
            n, err := c.client.Raw().Del(ctx, keys...).Result()
            if err != nil {
                return deleted, err
            }
            deleted += n
        }
        cursor = next
        if cursor == 0 {
            break
        }
    }
    return deleted, nil
}
```

### C.3 Redis Key Namespace

```go
// internal/shared/infrastructure/redis/keys.go
package redis

import "fmt"

type KeyBuilder struct{ prefix string }

func NewKeyBuilder(prefix string) *KeyBuilder {
    return &KeyBuilder{prefix: prefix}
}

func (k *KeyBuilder) Device(tenantID, deviceID string) string {
    return fmt.Sprintf("%s:device:%s:%s", k.prefix, tenantID, deviceID)
}

func (k *KeyBuilder) DeviceList(tenantID string, page int) string {
    return fmt.Sprintf("%s:devices:list:%s:page:%d", k.prefix, tenantID, page)
}

func (k *KeyBuilder) Customer(tenantID, customerID string) string {
    return fmt.Sprintf("%s:customer:%s:%s", k.prefix, tenantID, customerID)
}

func (k *KeyBuilder) Session(userID, sessionID string) string {
    return fmt.Sprintf("%s:session:%s:%s", k.prefix, userID, sessionID)
}

func (k *KeyBuilder) RateLimit(resource, key string) string {
    return fmt.Sprintf("%s:rl:%s:%s", k.prefix, resource, key)
}

func (k *KeyBuilder) Idempotency(key string) string {
    return fmt.Sprintf("%s:idem:%s", k.prefix, key)
}
```

### C.4 Rate Limiter (Sliding Window)

```go
// internal/shared/infrastructure/redis/rate_limit.go
package redis

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

type RateLimiter struct {
    client *Client
}

func NewRateLimiter(client *Client) *RateLimiter {
    return &RateLimiter{client: client}
}

type RateLimitResult struct {
    Allowed    bool
    Limit      int
    Remaining  int
    ResetAfter time.Duration
}

// Allow — sliding window
func (r *RateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (*RateLimitResult, error) {
    now := time.Now()
    windowStart := now.Add(-window)

    pipe := r.client.Raw().Pipeline()

    // 1. ลบ entry เก่า
    pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.UnixNano()))

    // 2. นับ
    countCmd := pipe.ZCard(ctx, key)

    // 3. เพิ่ม entry ใหม่
    pipe.ZAdd(ctx, key, redis.Z{
        Score:  float64(now.UnixNano()),
        Member: now.UnixNano(),
    })

    // 4. Set TTL
    pipe.Expire(ctx, key, window)

    if _, err := pipe.Exec(ctx); err != nil {
        return nil, err
    }

    count := countCmd.Val()
    result := &RateLimitResult{
        Limit:     limit,
        Remaining: limit - int(count) - 1,
        ResetAfter: window,
    }
    if int(count) >= limit {
        result.Allowed = false
        result.Remaining = 0
    } else {
        result.Allowed = true
    }
    return result, nil
}
```

### C.5 Rate Limit Middleware (ใช้ใน Interface Layer)

```go
// internal/modules/device/interfaces/middleware/rate_limit.go
package middleware

func RateLimit(limiter *redis.RateLimiter, limit int, window time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        tenantID := c.GetString("tenant_id")
        userID := c.GetString("user_id")
        key := fmt.Sprintf("rl:%s:%s:%s", tenantID, userID, c.FullPath())

        result, err := limiter.Allow(c.Request.Context(), key, limit, window)
        if err != nil {
            c.Next()
            return
        }

        c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", result.Limit))
        c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", result.Remaining))
        c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(result.ResetAfter).Unix()))

        if !result.Allowed {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "error": gin.H{
                    "code":    "RATE_LIMIT_EXCEEDED",
                    "message": "too many requests",
                },
            })
            return
        }
        c.Next()
    }
}
```

### C.6 Distributed Lock

```go
// internal/shared/infrastructure/redis/lock.go
package redis

type Lock struct {
    client *Client
}

type LockHandle struct {
    key   string
    token string
    ttl   time.Duration
}

// Acquire — SET NX EX
func (l *Lock) Acquire(ctx context.Context, key string, ttl time.Duration) (*LockHandle, error) {
    token := uuid.NewString()
    ok, err := l.client.Raw().SetNX(ctx, key, token, ttl).Result()
    if err != nil {
        return nil, err
    }
    if !ok {
        return nil, errors.New("lock already held")
    }
    return &LockHandle{key: key, token: token, ttl: ttl}, nil
}

// Release — Lua script (ตรวจ token ก่อนลบ)
var releaseScript = redis.NewScript(`
    if redis.call("GET", KEYS[1]) == ARGV[1] then
        return redis.call("DEL", KEYS[1])
    else
        return 0
    end
`)

func (l *Lock) Release(ctx context.Context, h *LockHandle) error {
    return releaseScript.Run(ctx, l.client.Raw(), []string{h.key}, h.token).Err()
}
```

---

## 🅳 PART 7D — KAFKA PRODUCERS & CONSUMERS

### D.1 Producer (outbox pattern)

ดูตัวอย่างที่ PART 6H H.4 — ใช้ `OutboxProducer` + `OutboxRelay` เพื่อ consistency

### D.2 Consumer Group Pattern

```go
// internal/shared/infrastructure/kafka/consumer.go
package kafka

import (
    "context"
    "fmt"
    "log"

    "github.com/IBM/sarama"
)

type ConsumerHandler func(ctx context.Context, msg *sarama.ConsumerMessage) error

type ConsumerGroup struct {
    group   sarama.ConsumerGroup
    topics  []string
    handler ConsumerHandler
    logger  application.Logger
}

func NewConsumerGroup(brokers []string, groupID string, topics []string, handler ConsumerHandler, logger application.Logger) (*ConsumerGroup, error) {
    cfg := sarama.NewConfig()
    cfg.Version = sarama.V3_0_0_0
    cfg.Consumer.Return.Errors = true
    cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
    cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

    group, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
    if err != nil {
        return nil, fmt.Errorf("create consumer group: %w", err)
    }
    return &ConsumerGroup{
        group:   group,
        topics:  topics,
        handler: handler,
        logger:  logger,
    }, nil
}

func (c *ConsumerGroup) Run(ctx context.Context) error {
    handler := &consumerGroupHandler{
        handler: c.handler,
        logger:  c.logger,
    }
    for {
        select {
        case <-ctx.Done():
            return nil
        default:
            if err := c.group.Consume(ctx, c.topics, handler); err != nil {
                c.logger.Error("consume error", application.Field{Key: "error", Value: err})
                return err
            }
        }
    }
}

func (c *ConsumerGroup) Close() error {
    return c.group.Close()
}

type consumerGroupHandler struct {
    handler ConsumerHandler
    logger  application.Logger
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        if err := h.handler(session.Context(), msg); err != nil {
            h.logger.Error("message handler failed",
                application.Field{Key: "topic", Value: msg.Topic},
                application.Field{Key: "error", Value: err},
            )
            // ไม่ mark — จะ retry
            continue
        }
        session.MarkMessage(msg, "")
    }
    return nil
}
```

### D.3 Telemetry Consumer

```go
// internal/modules/device/infrastructure/messaging/consumers/telemetry_consumer.go
package consumers

type TelemetryConsumer struct {
    ingestUC *application.IngestTelemetryUseCase
    logger   application.Logger
}

func NewTelemetryConsumer(uc *application.IngestTelemetryUseCase, logger application.Logger) *TelemetryConsumer {
    return &TelemetryConsumer{ingestUC: uc, logger: logger}
}

func (c *TelemetryConsumer) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
    var payload application.IngestTelemetryInput
    if err := json.Unmarshal(msg.Value, &payload); err != nil {
        c.logger.Error("unmarshal telemetry", application.Field{Key: "error", Value: err})
        return nil // skip bad message
    }
    if err := c.ingestUC.Execute(ctx, payload); err != nil {
        return fmt.Errorf("ingest telemetry: %w", err)
    }
    return nil
}

func (c *TelemetryConsumer) Register(brokers []string, logger application.Logger) (*kafka.ConsumerGroup, error) {
    return kafka.NewConsumerGroup(
        brokers,
        "device-telemetry-consumer",
        []string{"icmon.device.telemetry.ingested"},
        c.Handle,
        logger,
    )
}
```

### D.4 Event → WS Broadcaster Consumer

```go
// internal/modules/realtime/infrastructure/messaging/consumers/ws_broadcast_consumer.go
package consumers

type WSBroadcastConsumer struct {
    hub *websocket.Hub
}

func (c *WSBroadcastConsumer) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
    // Extract tenant จาก header
    var tenantID string
    for _, h := range msg.Headers {
        if string(h.Key) == "tenant_id" {
            tenantID = string(h.Value)
            break
        }
    }
    if tenantID == "" {
        return nil
    }

    c.hub.BroadcastToTenant(tenantID, websocket.WSEvent{
        Type: msg.Topic,
        Data: json.RawMessage(msg.Value),
    })
    return nil
}
```

### D.5 Consumer Registry

```go
// cmd/workers/device-telemetry/main.go
package main

func main() {
    // ... setup deps
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Telemetry consumer
    telConsumer := consumers.NewTelemetryConsumer(ingestUC, logger)
    telGroup, err := telConsumer.Register(cfg.KafkaBrokers, logger)
    if err != nil {
        log.Fatal(err)
    }
    go func() {
        if err := telGroup.Run(ctx); err != nil {
            log.Fatal(err)
        }
    }()

    // WS broadcast consumer
    wsConsumer := consumers.NewWSBroadcastConsumer(hub)
    wsGroup, _ := kafka.NewConsumerGroup(
        cfg.KafkaBrokers,
        "ws-broadcast-consumer",
        []string{
            "icmon.device.telemetry.ingested",
            "icmon.device.command.issued",
            "icmon.device.command.acked",
            "icmon.logistics.temperature.breached",
        },
        wsConsumer.Handle,
        logger,
    )
    go func() {
        if err := wsGroup.Run(ctx); err != nil {
            log.Fatal(err)
        }
    }()

    // Graceful shutdown
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    <-sigCh
    cancel()
    telGroup.Close()
    wsGroup.Close()
}
```

---

## 🅴 PART 7E — MQTT SUBSCRIBER & PUBLISHER

### E.1 MQTT Client Wrapper

```go
// internal/shared/infrastructure/mqtt/client.go
package mqtt

import (
    mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
    client mqtt.Client
    logger application.Logger
}

func NewClient(cfg Config, logger application.Logger) (*Client, error) {
    opts := mqtt.NewClientOptions().
        AddBroker(cfg.Broker).
        SetClientID(cfg.ClientID).
        SetUsername(cfg.Username).
        SetPassword(cfg.Password).
        SetAutoReconnect(true).
        SetConnectRetry(true).
        SetConnectRetryInterval(5 * time.Second).
        SetMaxReconnectInterval(1 * time.Minute).
        SetKeepAlive(30 * time.Second).
        SetPingTimeout(10 * time.Second).
        SetCleanSession(false).
        SetOrderMatters(false)

    // TLS
    if cfg.UseTLS {
        tlsConfig := &tls.Config{
            InsecureSkipVerify: cfg.InsecureSkipVerify,
        }
        opts.SetTLSConfig(tlsConfig)
    }

    // Last Will
    opts.SetWill("gateway/status", "offline", 1, true)

    c := mqtt.NewClient(opts)
    if token := c.Connect(); token.Wait() && token.Error() != nil {
        return nil, token.Error()
    }

    return &Client{client: c, logger: logger}, nil
}

func (c *Client) Subscribe(topic string, qos byte, handler mqtt.MessageHandler) error {
    token := c.client.Subscribe(topic, qos, handler)
    token.Wait()
    return token.Error()
}

func (c *Client) Publish(topic string, qos byte, retained bool, payload []byte) error {
    token := c.client.Publish(topic, qos, retained, payload)
    token.Wait()
    return token.Error()
}

func (c *Client) Disconnect() {
    c.client.Disconnect(250)
}
```

### E.2 MQTT Subscriber (Telemetry)

```go
// internal/modules/device/infrastructure/iot/mqtt_subscriber.go
package iot

type MQTTSubscriber struct {
    client   *mqtt.Client
    ingestUC *application.IngestTelemetryUseCase
    logger   application.Logger
}

func NewMQTTSubscriber(client *mqtt.Client, uc *application.IngestTelemetryUseCase, logger application.Logger) *MQTTSubscriber {
    return &MQTTSubscriber{client: client, ingestUC: uc, logger: logger}
}

func (s *MQTTSubscriber) Start(ctx context.Context) error {
    // Subscribe หลาย topic
    topics := map[string]byte{
        "iot/+/telemetry":       1,  // QoS 1
        "iot/+/heartbeat":       0,
        "iot/+/status":          1,
        "iot/+/fault":           1,
        "device/+/shadow/update": 1,
    }
    for topic, qos := range topics {
        if err := s.client.Subscribe(topic, qos, s.handle(ctx, topic)); err != nil {
            return fmt.Errorf("subscribe %s: %w", topic, err)
        }
        s.logger.Info("subscribed", application.Field{Key: "topic", Value: topic})
    }
    return nil
}

func (s *MQTTSubscriber) handle(ctx context.Context, topic string) mqtt.MessageHandler {
    return func(_ mqtt.Client, msg mqtt.Message) {
        // Extract device serial จาก topic: iot/{serial}/telemetry
        parts := strings.Split(msg.Topic(), "/")
        if len(parts) < 3 {
            return
        }
        serial := parts[1]
        action := parts[2]

        switch action {
        case "telemetry":
            s.handleTelemetry(ctx, serial, msg.Payload())
        case "heartbeat":
            s.handleHeartbeat(ctx, serial)
        case "status":
            s.handleStatus(ctx, serial, msg.Payload())
        case "fault":
            s.handleFault(ctx, serial, msg.Payload())
        }
    }
}

func (s *MQTTSubscriber) handleTelemetry(ctx context.Context, serial string, payload []byte) {
    var p struct {
        Timestamp time.Time              `json:"ts"`
        Metrics   []struct {
            Name  string  `json:"metric"`
            Value float64 `json:"value"`
            Unit  string  `json:"unit"`
        } `json:"metrics"`
    }
    if err := json.Unmarshal(payload, &p); err != nil {
        s.logger.Warn("invalid telemetry payload", application.Field{Key: "error", Value: err})
        return
    }

    // Convert to use case input
    // (ต้อง resolve deviceID จาก serial ก่อน — ใช้ cache หรือ repo)
    // ...
    _ = s.ingestUC.Execute(ctx, application.IngestTelemetryInput{
        // DeviceID, Metrics, Timestamp
    })
}
```

### E.3 MQTT Publisher (Command → Device)

```go
// internal/modules/device/infrastructure/iot/mqtt_publisher.go
package iot

type MQTTPublisher struct {
    client *mqtt.Client
}

func (p *MQTTPublisher) PublishCommand(ctx context.Context, cmd *entity.Command, deviceSerial string) error {
    topic := fmt.Sprintf("iot/%s/command", deviceSerial)
    payload := map[string]interface{}{
        "command_id": cmd.ID.String(),
        "command":    cmd.Type,
        "payload":    cmd.Payload,
        "expires_at": cmd.ExpiresAt,
        "issued_at":  cmd.IssuedAt,
    }
    data, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    return p.client.Publish(topic, 1, false, data)
}

func (p *MQTTPublisher) PublishShadowDesired(ctx context.Context, deviceSerial string, desired map[string]interface{}, version int) error {
    topic := fmt.Sprintf("device/%s/shadow/desired", deviceSerial)
    payload := map[string]interface{}{
        "desired": desired,
        "version": version,
    }
    data, _ := json.Marshal(payload)
    return p.client.Publish(topic, 1, false, data)
}
```

### E.4 Protocol Parsers (Modbus, LoRaWAN)

```go
// internal/modules/device/infrastructure/iot/protocol/parser.go
package protocol

type Parser interface {
    Parse(payload []byte) ([]TelemetryReading, error)
    Protocol() string
}

type TelemetryReading struct {
    Metric string
    Value  float64
    Unit   string
}

// Registry
var parsers = map[string]Parser{}

func Register(p Parser) {
    parsers[p.Protocol()] = p
}

func Get(protocol string) (Parser, bool) {
    p, ok := parsers[protocol]
    return p, ok
}
```

```go
// internal/modules/device/infrastructure/iot/protocol/modbus.go
package protocol

type ModbusParser struct{}

func (p *ModbusParser) Protocol() string { return "MODBUS" }

func (p *ModbusParser) Parse(payload []byte) ([]TelemetryReading, error) {
    if len(payload) < 6 {
        return nil, errors.New("modbus payload too short")
    }
    // Parse Modbus RTU/TCP frame
    // [slave_id][func_code][reg_addr_hi][reg_addr_lo][data...][crc]
    // ...
    return []TelemetryReading{}, nil
}

// internal/modules/device/infrastructure/iot/protocol/lorawan.go
type LoRaWANParser struct {
    // decoder ตาม payload spec
}

func (p *LoRaWANParser) Protocol() string { return "LORAWAN" }

func (p *LoRaWANParser) Parse(payload []byte) ([]TelemetryReading, error) {
    // Cayenne LPP format
    // ...
    return []TelemetryReading{}, nil
}

// init — register parsers
func init() {
    Register(&ModbusParser{})
    Register(&LoRaWANParser{})
}
```

---

## 🅵 PART 7F — INFLUXDB REPOSITORIES

### F.1 InfluxDB Client Wrapper

```go
// internal/shared/infrastructure/influxdb/client.go
package influxdb

import (
    "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"
)

type Client struct {
    client   influxdb2.Client
    writeAPI api.WriteAPIBlocking
    queryAPI api.QueryAPI
    bucket   string
    org      string
}

func NewClient(cfg Config) *Client {
    c := influxdb2.NewClient(cfg.URL, cfg.Token)
    return &Client{
        client:   c,
        writeAPI: c.WriteAPIBlocking(cfg.Org, cfg.Bucket),
        queryAPI: c.QueryAPI(cfg.Org),
        bucket:   cfg.Bucket,
        org:      cfg.Org,
    }
}

func (c *Client) Close() {
    c.client.Close()
}

func (c *Client) Ping(ctx context.Context) error {
    return c.client.Ping(ctx)
}
```

### F.2 Telemetry Repository

```go
// internal/modules/device/infrastructure/persistence/influxdb/telemetry_repo.go
package influxdb

import (
    "context"
    "fmt"
    "time"

    "github.com/influxdata/influxdb-client-go/v2/api/write"
    "github.com/google/uuid"

    "icmongolang/internal/modules/device/domain/entity"
    "icmongolang/internal/modules/device/domain/repository"
    sharedinflux "icmongolang/internal/shared/infrastructure/influxdb"
)

type TelemetryRepoImpl struct {
    client *sharedinflux.Client
    bucket string
}

func NewTelemetryRepository(client *sharedinflux.Client, bucket string) repository.TelemetryRepository {
    return &TelemetryRepoImpl{client: client, bucket: bucket}
}

// ============================================================
// WRITE
// ============================================================
func (r *TelemetryRepoImpl) Write(ctx context.Context, t *entity.Telemetry) error {
    point := write.NewPoint(
        "device_telemetry",
        map[string]string{
            "tenant_id": t.TenantID.String(),
            "device_id": t.DeviceID.String(),
            "site_id":   t.SiteID.String(),
            "metric":    t.Metric,
            "unit":      t.Unit,
            "quality":   t.Quality,
        },
        map[string]interface{}{
            "value": t.Value,
        },
        t.Timestamp,
    )
    return r.client.WriteAPI().WritePoint(ctx, point)
}

func (r *TelemetryRepoImpl) WriteBatch(ctx context.Context, batch []*entity.Telemetry) error {
    points := make([]*write.Point, 0, len(batch))
    for _, t := range batch {
        points = append(points, write.NewPoint(
            "device_telemetry",
            map[string]string{
                "tenant_id": t.TenantID.String(),
                "device_id": t.DeviceID.String(),
                "site_id":   t.SiteID.String(),
                "metric":    t.Metric,
                "unit":      t.Unit,
                "quality":   t.Quality,
            },
            map[string]interface{}{"value": t.Value},
            t.Timestamp,
        ))
    }
    return r.client.WriteAPI().WritePoint(ctx, points...)
}

// ============================================================
// QUERY
// ============================================================
func (r *TelemetryRepoImpl) Query(ctx context.Context, q repository.TelemetryQuery) (*repository.TelemetryResult, error) {
    interval := q.Interval
    if interval == "" {
        interval = "1m"
    }

    flux := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r._measurement == "device_telemetry")
            |> filter(fn: (r) => r.tenant_id == "%s")
            |> filter(fn: (r) => r.device_id == "%s")
            |> filter(fn: (r) => r._field == "value")
            |> aggregateWindow(every: %s, fn: mean, createEmpty: false)
            |> keep(columns: ["_time", "_value", "metric", "unit"])
    `,
        r.bucket,
        q.From.Format(time.RFC3339),
        q.To.Format(time.RFC3339),
        q.TenantID.String(),
        q.DeviceID.String(),
        interval,
    )

    result, err := r.client.QueryAPI().Query(ctx, flux)
    if err != nil {
        return nil, fmt.Errorf("influx query: %w", err)
    }
    defer result.Close()

    series := make(map[string][]repository.TelemetryPoint)
    for result.Next() {
        rec := result.Record()
        metric, _ := rec.ValueByKey("metric").(string)
        unit, _ := rec.ValueByKey("unit").(string)
        val, _ := rec.Value().(float64)

        series[metric] = append(series[metric], repository.TelemetryPoint{
            Timestamp: rec.Time(),
            Value:     val,
            Unit:      unit,
        })
    }
    if result.Err() != nil {
        return nil, result.Err()
    }

    return &repository.TelemetryResult{
        DeviceID: q.DeviceID,
        Series:   series,
    }, nil
}

func (r *TelemetryRepoImpl) Latest(ctx context.Context, tenantID, deviceID uuid.UUID, metric string) (*entity.Telemetry, error) {
    flux := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: -24h)
            |> filter(fn: (r) => r._measurement == "device_telemetry")
            |> filter(fn: (r) => r.tenant_id == "%s")
            |> filter(fn: (r) => r.device_id == "%s")
            |> filter(fn: (r) => r.metric == "%s")
            |> filter(fn: (r) => r._field == "value")
            |> last()
    `, r.bucket, tenantID.String(), deviceID.String(), metric)

    result, err := r.client.QueryAPI().Query(ctx, flux)
    if err != nil {
        return nil, err
    }
    defer result.Close()

    if !result.Next() {
        return nil, nil
    }
    rec := result.Record()
    val, _ := rec.Value().(float64)
    unit, _ := rec.ValueByKey("unit").(string)

    return &entity.Telemetry{
        TenantID:  tenantID,
        DeviceID:  deviceID,
        Metric:    metric,
        Value:     val,
        Unit:      unit,
        Timestamp: rec.Time(),
    }, nil
}

func (r *TelemetryRepoImpl) Aggregate(ctx context.Context, q repository.AggregateQuery) ([]repository.AggregatePoint, error) {
    fn := q.Function
    if fn == "" {
        fn = "mean"
    }
    flux := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r._measurement == "device_telemetry")
            |> filter(fn: (r) => r.tenant_id == "%s")
            |> filter(fn: (r) => r.device_id == "%s")
            |> filter(fn: (r) => r.metric == "%s")
            |> filter(fn: (r) => r._field == "value")
            |> aggregateWindow(every: %s, fn: %s, createEmpty: false)
    `, r.bucket, q.From.Format(time.RFC3339), q.To.Format(time.RFC3339),
        q.TenantID, q.DeviceID, q.Metric, q.Window, fn)

    result, err := r.client.QueryAPI().Query(ctx, flux)
    if err != nil {
        return nil, err
    }
    defer result.Close()

    var points []repository.AggregatePoint
    for result.Next() {
        rec := result.Record()
        val, _ := rec.Value().(float64)
        points = append(points, repository.AggregatePoint{
            Timestamp: rec.Time(),
            Value:     val,
        })
    }
    return points, result.Err()
}

func (r *TelemetryRepoImpl) DeleteOlderThan(ctx context.Context, before time.Time) (int64, error) {
    // ใช้ delete API
    // ...
    return 0, nil
}
```

### F.3 Retention Policies (Setup)

```go
// scripts/setup_influx.go
func setupInflux() {
    // สร้าง bucket + retention
    // - iot_telemetry (30d retention)
    // - iot_telemetry_archive (365d)
    // - iot_telemetry_5m (aggregated, 365d)
    // ใช้ InfluxDB API
}
```

---

## 🅶 PART 7G — ELASTICSEARCH INDEXERS

### G.1 ES Client Wrapper

```go
// internal/shared/infrastructure/elasticsearch/client.go
package elasticsearch

import (
    "github.com/elastic/go-elasticsearch/v8"
)

type Client struct {
    es     *elasticsearch.Client
    prefix string
}

func NewClient(urls []string, prefix string) (*Client, error) {
    es, err := elasticsearch.NewClient(elasticsearch.Config{
        Addresses: urls,
    })
    if err != nil {
        return nil, err
    }
    return &Client{es: es, prefix: prefix}, nil
}

func (c *Client) Index(name string) string {
    return c.prefix + "_" + name
}

func (c *Client) Raw() *elasticsearch.Client { return c.es }
```

### G.2 Device Indexer

```go
// internal/modules/device/infrastructure/search/elasticsearch/device_indexer.go
package elasticsearch

type DeviceIndexer struct {
    client *sharedelastic.Client
}

const deviceMapping = `{
    "mappings": {
        "properties": {
            "id":          {"type": "keyword"},
            "tenant_id":   {"type": "keyword"},
            "customer_id": {"type": "keyword"},
            "site_id":     {"type": "keyword"},
            "serial_no":   {"type": "keyword"},
            "name":        {"type": "text", "fields": {"keyword": {"type": "keyword"}}},
            "type":        {"type": "keyword"},
            "protocol":    {"type": "keyword"},
            "status":      {"type": "keyword"},
            "tags":        {"type": "keyword"},
            "last_seen_at":{"type": "date"},
            "created_at":  {"type": "date"}
        }
    }
}`

func (i *DeviceIndexer) EnsureIndex(ctx context.Context) error {
    index := i.client.Index("devices")
    res, err := i.client.Raw().Indices.Exists([]string{index})
    if err != nil {
        return err
    }
    defer res.Body.Close()
    if res.StatusCode == 200 {
        return nil // already exists
    }
    // Create
    createRes, err := i.client.Raw().Indices.Create(
        index,
        i.client.Raw().Indices.Create.WithBody(strings.NewReader(deviceMapping)),
    )
    if err != nil {
        return err
    }
    defer createRes.Body.Close()
    return nil
}

func (i *DeviceIndexer) Index(ctx context.Context, d *entity.Device) error {
    doc := map[string]interface{}{
        "id":           d.ID.String(),
        "tenant_id":    d.TenantID.String(),
        "customer_id":  d.CustomerID.String(),
        "site_id":      d.SiteID.String(),
        "serial_no":    d.SerialNo.String(),
        "name":         d.Name,
        "type":         string(d.Type),
        "protocol":     string(d.Protocol),
        "status":       string(d.Status),
        "tags":         d.Tags,
        "last_seen_at": d.LastSeenAt,
        "created_at":   d.CreatedAt,
    }
    body, _ := json.Marshal(doc)
    res, err := i.client.Raw().Index(
        i.client.Index("devices"),
        bytes.NewReader(body),
        i.client.Raw().Index.WithDocumentID(d.ID.String()),
        i.client.Raw().Index.WithContext(ctx),
    )
    if err != nil {
        return err
    }
    defer res.Body.Close()
    return nil
}

func (i *DeviceIndexer) Search(ctx context.Context, tenantID uuid.UUID, query string, from, size int) (*SearchResult, error) {
    q := map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "filter": []map[string]interface{}{
                    {"term": map[string]interface{}{"tenant_id": tenantID.String()}},
                },
                "must": []map[string]interface{}{
                    {"multi_match": map[string]interface{}{
                        "query":  query,
                        "fields": []string{"name^2", "serial_no^3", "tags"},
                    }},
                },
            },
        },
        "from": from,
        "size": size,
    }
    body, _ := json.Marshal(q)
    res, err := i.client.Raw().Search(
        i.client.Raw().Search.WithIndex(i.client.Index("devices")),
        i.client.Raw().Search.WithBody(bytes.NewReader(body)),
        i.client.Raw().Search.WithContext(ctx),
    )
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    var result SearchResult
    if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
        return nil, err
    }
    return &result, nil
}
```

### G.3 ES Consumer (sync from Kafka)

```go
// internal/modules/device/infrastructure/messaging/consumers/es_indexer_consumer.go
package consumers

type ESIndexerConsumer struct {
    indexer *elasticsearch.DeviceIndexer
    repo    repository.DeviceRepository
}

func (c *ESIndexerConsumer) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
    var evt map[string]interface{}
    if err := json.Unmarshal(msg.Value, &evt); err != nil {
        return nil
    }
    deviceID, err := uuid.Parse(evt["aggregate_id"].(string))
    if err != nil {
        return nil
    }
    tenantID, _ := uuid.Parse(evt["tenant_id"].(string))

    // Load from DB (read model)
    device, err := c.repo.FindByID(ctx, tenantID, deviceID)
    if err != nil {
        return nil
    }
    return c.indexer.Index(ctx, device)
}
```

---

## 🅷 PART 7H — EXTERNAL SERVICES

### H.1 LLM Service (Ollama)

```go
// internal/shared/infrastructure/services/llm/ollama.go
package llm

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type OllamaClient struct {
    baseURL string
    model   string
    http    *http.Client
}

func NewOllamaClient(baseURL, model string) *OllamaClient {
    return &OllamaClient{
        baseURL: baseURL,
        model:   model,
        http:    &http.Client{Timeout: 60 * time.Second},
    }
}

type GenerateRequest struct {
    Model   string                 `json:"model"`
    Prompt  string                 `json:"prompt"`
    Stream  bool                   `json:"stream"`
    Format  string                 `json:"format,omitempty"` // "json"
    Options map[string]interface{} `json:"options,omitempty"`
}

type GenerateResponse struct {
    Model     string `json:"model"`
    Response  string `json:"response"`
    Done      bool   `json:"done"`
    CreatedAt string `json:"created_at"`
}

func (c *OllamaClient) Generate(ctx context.Context, prompt string, format string) (string, error) {
    req := GenerateRequest{
        Model:  c.model,
        Prompt: prompt,
        Stream: false,
        Format: format,
        Options: map[string]interface{}{
            "temperature": 0.2,
            "top_p":       0.9,
            "num_predict": 500,
        },
    }
    body, _ := json.Marshal(req)
    httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/generate", bytes.NewReader(body))
    if err != nil {
        return "", err
    }
    httpReq.Header.Set("Content-Type", "application/json")

    resp, err := c.http.Do(httpReq)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return "", fmt.Errorf("ollama status %d", resp.StatusCode)
    }

    var out GenerateResponse
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
        return "", err
    }
    return out.Response, nil
}
```

### H.2 AI Insight Service (ใช้ LLM)

```go
// internal/modules/report/infrastructure/services/ai/insight_service.go
package ai

type InsightService struct {
    llm    *llm.OllamaClient
    prompt *PromptTemplate
}

type InsightInput struct {
    TenantID   uuid.UUID
    DeviceID   uuid.UUID
    DeviceName string
    Metrics    []MetricSummary
    Period     string
}

type MetricSummary struct {
    Metric string
    Min    float64
    Max    float64
    Avg    float64
    Unit   string
}

type InsightOutput struct {
    Summary         string   `json:"summary"`
    Anomalies       []string `json:"anomalies"`
    Recommendations []string `json:"recommendations"`
    Severity        string   `json:"severity"` // INFO, WARNING, CRITICAL
    Confidence      float64  `json:"confidence"`
}

func (s *InsightService) Analyze(ctx context.Context, in InsightInput) (*InsightOutput, error) {
    prompt := s.prompt.Render(in)
    resp, err := s.llm.Generate(ctx, prompt, "json")
    if err != nil {
        return nil, err
    }
    var out InsightOutput
    if err := json.Unmarshal([]byte(resp), &out); err != nil {
        return nil, fmt.Errorf("parse LLM response: %w", err)
    }
    return &out, nil
}
```

### H.3 Payment Service (Stripe)

```go
// internal/shared/infrastructure/services/payment/stripe.go
package payment

import (
    "context"
    "github.com/stripe/stripe-go/v75"
    "github.com/stripe/stripe-go/v75/checkout/session"
    "github.com/stripe/stripe-go/v75/webhook"
)

type StripeClient struct {
    secretKey     string
    webhookSecret string
}

func NewStripeClient(secretKey, webhookSecret string) *StripeClient {
    stripe.Key = secretKey
    return &StripeClient{secretKey: secretKey, webhookSecret: webhookSecret}
}

type CheckoutRequest struct {
    TenantID   uuid.UUID
    CustomerID uuid.UUID
    Amount     float64
    Currency   string
    Reference  string
    ReturnURL  string
    Metadata   map[string]string
}

type CheckoutResult struct {
    ID          string
    CheckoutURL string
    ExpiresAt   time.Time
}

func (c *StripeClient) CreateCheckout(ctx context.Context, in CheckoutRequest) (*CheckoutResult, error) {
    params := &stripe.CheckoutSessionParams{
        Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
        LineItems: []*stripe.CheckoutSessionLineItemParams{
            {
                PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
                    Currency: stripe.String(strings.ToLower(in.Currency)),
                    ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
                        Name: stripe.String("Subscription " + in.Reference),
                    },
                    UnitAmount: stripe.Int64(int64(in.Amount * 100)),
                },
                Quantity: stripe.Int64(1),
            },
        },
        SuccessURL: stripe.String(in.ReturnURL + "?status=success"),
        CancelURL:  stripe.String(in.ReturnURL + "?status=cancelled"),
        Metadata:   in.Metadata,
    }
    sess, err := session.New(params)
    if err != nil {
        return nil, err
    }
    return &CheckoutResult{
        ID:          sess.ID,
        CheckoutURL: sess.URL,
        ExpiresAt:   time.Unix(sess.ExpiresAt, 0),
    }, nil
}

func (c *StripeClient) VerifyWebhook(payload []byte, signature string) (*stripe.Event, error) {
    return webhook.ConstructEvent(payload, signature, c.webhookSecret)
}
```

### H.4 Email Service (SMTP)

```go
// internal/shared/infrastructure/services/email/smtp.go
package email

import (
    "context"
    "fmt"
    "net/smtp"
    "strings"
)

type SMTPService struct {
    host     string
    port     int
    username string
    password string
    from     string
}

type Email struct {
    To      []string
    CC      []string
    Subject string
    Body    string
    IsHTML  bool
}

func (s *SMTPService) Send(ctx context.Context, e Email) error {
    headers := map[string]string{
        "From":         s.from,
        "To":           strings.Join(e.To, ", "),
        "Subject":      e.Subject,
        "MIME-Version": "1.0",
    }
    if e.IsHTML {
        headers["Content-Type"] = "text/html; charset=UTF-8"
    } else {
        headers["Content-Type"] = "text/plain; charset=UTF-8"
    }

    var msg strings.Builder
    for k, v := range headers {
        msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
    }
    msg.WriteString("\r\n")
    msg.WriteString(e.Body)

    auth := smtp.PlainAuth("", s.username, s.password, s.host)
    addr := fmt.Sprintf("%s:%d", s.host, s.port)
    return smtp.SendMail(addr, auth, s.from, e.To, []byte(msg.String()))
}
```

---

## 🅸 PART 7I — SCHEDULERS & BACKGROUND JOBS

### I.1 Scheduler Setup

```go
// internal/shared/infrastructure/scheduler/scheduler.go
package scheduler

import (
    "context"
    "github.com/robfig/cron/v3"
)

type Scheduler struct {
    cron   *cron.Cron
    logger application.Logger
}

func NewScheduler(logger application.Logger) *Scheduler {
    c := cron.New(cron.WithSeconds())
    return &Scheduler{cron: c, logger: logger}
}

func (s *Scheduler) Add(spec string, name string, job func(ctx context.Context) error) error {
    _, err := s.cron.AddFunc(spec, func() {
        ctx := context.Background()
        s.logger.Info("job started", application.Field{Key: "job", Value: name})
        if err := job(ctx); err != nil {
            s.logger.Error("job failed",
                application.Field{Key: "job", Value: name},
                application.Field{Key: "error", Value: err},
            )
        }
    })
    return err
}

func (s *Scheduler) Start() { s.cron.Start() }
func (s *Scheduler) Stop()  { s.cron.Stop() }
```

### I.2 Device Offline Detection Job

```go
// internal/modules/device/infrastructure/scheduler/offline_detector.go
package scheduler

type OfflineDetectorJob struct {
    repo      repository.DeviceRepository
    producer  application.EventProducer
    threshold time.Duration
    logger    application.Logger
}

func (j *OfflineDetectorJob) Run(ctx context.Context) error {
    devices, err := j.repo.FindStale(ctx, uuid.Nil, j.threshold, 1000)
    // Note: ต้องดึงทุก tenant — ใช้ system context
    if err != nil {
        return err
    }

    offline := make([]*entity.Device, 0)
    for _, d := range devices {
        if d.MarkOfflineIfStale(j.threshold) {
            offline = append(offline, d)
        }
    }
    if len(offline) == 0 {
        return nil
    }

    // Batch update
    if err := j.repo.UpdateBatch(ctx, offline); err != nil {
        return err
    }

    // Publish events
    for _, d := range offline {
        events := d.PullEvents()
        _ = j.producer.PublishBatch(ctx, events)
    }

    j.logger.Info("offline detection complete",
        application.Field{Key: "count", Value: len(offline)},
    )
    return nil
}

// Register: ทุก 1 นาที
// scheduler.Add("0 * * * * *", "device-offline-detector", job.Run)
```

### I.3 Subscription Expiry Job

```go
// internal/modules/package/infrastructure/scheduler/expiry_job.go
package scheduler

type SubscriptionExpiryJob struct {
    repo     repository.SubscriptionRepository
    producer application.EventProducer
}

func (j *SubscriptionExpiryJob) Run(ctx context.Context) error {
    // หา subscription ที่ใกล้หมดอายุใน 3 วัน
    subs, err := j.repo.FindExpiringWithin(ctx, 3*24*time.Hour)
    if err != nil {
        return err
    }
    for _, sub := range subs {
        sub.MarkExpiring()
        _ = j.repo.Update(ctx, sub)
        _ = j.producer.PublishBatch(ctx, sub.PullEvents())
    }

    // หา subscription ที่หมดอายุแล้ว → EXPIRED
    expired, err := j.repo.FindExpired(ctx)
    if err != nil {
        return err
    }
    for _, sub := range expired {
        _ = sub.Expire()
        _ = j.repo.Update(ctx, sub)
        _ = j.producer.PublishBatch(ctx, sub.PullEvents())
    }
    return nil
}

// Register: ทุก 1 ชั่วโมง
// scheduler.Add("0 0 * * * *", "subscription-expiry", job.Run)
```

### I.4 KPI Snapshot Job

```go
// internal/modules/report/infrastructure/scheduler/kpi_snapshot_job.go
package scheduler

type KPISnapshotJob struct {
    kpiRepo    repository.KPIRepository
    metrics    MetricsProvider  // ดึงค่าจาก module ต่างๆ
    tenants    TenantLister
}

func (j *KPISnapshotJob) Run(ctx context.Context) error {
    tenants, err := j.tenants.ListActive(ctx)
    if err != nil {
        return err
    }
    today := time.Now().UTC().Truncate(24 * time.Hour)

    for _, tenantID := range tenants {
        kpis := []struct {
            Code  string
            Value float64
        }{
            {"MRR", j.metrics.GetMRR(ctx, tenantID)},
            {"ACTIVE_CUSTOMERS", j.metrics.GetActiveCustomers(ctx, tenantID)},
            {"DEVICE_ONLINE", j.metrics.GetOnlineDevices(ctx, tenantID)},
            {"DEVICE_UPTIME_PCT", j.metrics.GetDeviceUptime(ctx, tenantID)},
            {"ALERT_COUNT", j.metrics.GetAlertCount(ctx, tenantID)},
        }
        for _, kpi := range kpis {
            _ = j.kpiRepo.SaveSnapshot(ctx, repository.KPISnapshot{
                TenantID:     tenantID,
                KPICode:      kpi.Code,
                Value:        kpi.Value,
                SnapshotDate: today,
                ComputedAt:   time.Now().UTC(),
                Source:       "scheduler",
            })
        }
    }
    return nil
}

// Register: ทุกวัน 02:00
// scheduler.Add("0 0 2 * * *", "kpi-snapshot", job.Run)
```

### I.5 Scheduler Registry

```go
// cmd/scheduler/main.go
package main

func main() {
    // ... setup deps
    sched := scheduler.NewScheduler(logger)

    // Device
    offlineJob := devicescheduler.NewOfflineDetectorJob(deviceRepo, producer, 5*time.Minute, logger)
    _ = sched.Add("0 * * * * *", "device-offline-detector", offlineJob.Run) // ทุกนาที

    // Package
    expiryJob := packagescheduler.NewSubscriptionExpiryJob(subRepo, producer)
    _ = sched.Add("0 0 * * * *", "subscription-expiry", expiryJob.Run) // ทุกชั่วโมง

    // Report
    kpiJob := reportscheduler.NewKPISnapshotJob(kpiRepo, metricsProvider, tenantLister)
    _ = sched.Add("0 0 2 * * *", "kpi-snapshot", kpiJob.Run) // ทุกวัน 02:00

    // Cleanup
    cleanupJob := devicscheduler.NewCleanupJob(telemetryRepo, auditRepo)
    _ = sched.Add("0 0 3 * * *", "cleanup-old-data", cleanupJob.Run) // ทุกวัน 03:00

    sched.Start()
    defer sched.Stop()

    // Graceful shutdown
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    <-sigCh
}
```

---

## 🎯 PART 7 — SUMMARY

### Infrastructure Layer Deliverables

| Category | Count | Location |
|---|:-:|---|
| **Repository Impls** | 25+ | `infrastructure/persistence/postgres/` |
| **GORM Models** | 60+ | `infrastructure/persistence/postgres/models.go` |
| **Mappers** | 25+ | `infrastructure/persistence/postgres/mappers.go` |
| **InfluxDB Repos** | 3 | `infrastructure/persistence/influxdb/` |
| **Redis Adapters** | 4 | `infrastructure/persistence/redis/` |
| **Kafka Producers** | 1 | `infrastructure/messaging/` |
| **Kafka Consumers** | 10+ | `infrastructure/messaging/consumers/` |
| **MQTT Subscribers** | 3 | `infrastructure/iot/` |
| **MQTT Publishers** | 2 | `infrastructure/iot/` |
| **Protocol Parsers** | 5 | `infrastructure/iot/protocol/` |
| **ES Indexers** | 5 | `infrastructure/search/elasticsearch/` |
| **External Services** | 6 | `infrastructure/services/` |
| **Schedulers** | 8+ | `infrastructure/scheduler/` |

### Layer Dependency Check

```
✅ Infrastructure imports:
   - domain (entities, repos, events, errors, VOs)
   - application (ports, use cases)
   - gorm, sarama, redis, influx, elastic, mqtt
   - stdlib, uuid, decimal

✅ Infrastructure implements:
   - repository.DeviceRepository
   - repository.TelemetryRepository
   - application.EventProducer
   - application.Cache
   - application.Logger
   - application.MetricsRecorder
```

---

# 🏁 สรุป PART 5, 6, 7 — CORE LAYERS

```
┌─────────────────────────────────────────────────────────────────┐
│                  3 CORE LAYERS — OVERVIEW                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ★ PART 5 — DOMAIN LAYER                                        │
│  ─────────────────────────                                      │
│  • 18 Aggregates · 80+ Value Objects                            │
│  • 60+ Domain Events · 100+ Domain Errors                       │
│  • 25+ Repository Interfaces · 15+ Domain Services              │
│  • 20+ Specifications · 7 Factories                             │
│  • Pure Go — ห้าม import ORM/HTTP/MQTT                          │
│                                                                 │
│  ★ PART 6 — APPLICATION LAYER                                   │
│  ──────────────────────────                                     │
│  • 95 Command Use Cases · 64 Query Handlers                     │
│  • 240+ DTOs · 40 Mappers                                       │
│  • 8 Ports (EventProducer, Cache, Logger, ...)                  │
│  • Transaction + Idempotency + Orchestration                    │
│  • ไม่มี SQL/HTTP/MQTT โดยตรง                                   │
│                                                                 │
│  ★ PART 7 — INFRASTRUCTURE LAYER                                │
│  ─────────────────────────────                                  │
│  • 25+ Repository Impls · 60+ GORM Models                       │
│  • Kafka Producer + 10+ Consumers                               │
│  • MQTT Subscriber + Publisher + 5 Protocol Parsers             │
│  • InfluxDB Repos · Redis Adapters · ES Indexers                │
│  • 6 External Services (LLM, Payment, Email, SMS, ...)          │
│  • 8+ Schedulers                                                │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 📁 Final Folder Structure

```
internal/modules/device/
├── domain/                    ← PART 5
│   ├── entity/
│   ├── value_object/
│   ├── event/
│   ├── repository/            (interfaces)
│   ├── service/
│   ├── specification/
│   ├── factory/
│   └── errors/
├── application/               ← PART 6
│   ├── register_device.go
│   ├── ingest_telemetry.go
│   ├── dto.go
│   ├── mappers.go
│   ├── ports.go
│   └── query/
├── infrastructure/            ← PART 7
│   ├── persistence/
│   │   ├── postgres/
│   │   ├── redis/
│   │   └── influxdb/
│   ├── messaging/
│   ├── iot/
│   ├── search/
│   ├── services/
│   └── scheduler/
├── interfaces/                ← (PART 8 — Interface Layer)
│   ├── http/
│   ├── websocket/
│   └── mqtt/
└── module.go                  ← Composition Root (Wire-up)
```

### 📋 Response ถัดไปที่เป็นไปได้

| PART | เนื้อหา |
|:---|:---|
| **PART 8** | Interface Layer Deep Dive (HTTP handlers, WebSocket, MQTT Gateway, Middleware, Routes) |
| **PART 12** | Testing Strategy (Unit, Integration, E2E, Contract, Load) |
| **PART 13** | Deployment & DevOps (Docker, K8s, CI/CD, Monitoring) |
| **PART 14** | Security Deep Dive (Auth, PDPA, Encryption, Pen Test) |

 