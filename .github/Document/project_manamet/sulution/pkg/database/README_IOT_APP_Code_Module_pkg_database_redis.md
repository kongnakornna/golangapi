# Module 10: pkg/redis (Redis Client & Utilities) – ฉบับสมบูรณ์

## สำหรับโฟลเดอร์ `internal/pkg/redis/`

ไฟล์ที่เกี่ยวข้อง:
- `internal/pkg/redis/client.go`
- `internal/pkg/redis/cache.go`
- `internal/pkg/redis/refresh_store.go`
- `internal/pkg/redis/blacklist.go`
- `internal/pkg/redis/lock.go`

---

## หลักการ (Concept)

### Redis คืออะไร?
Redis (REmote DIctionary Server) คือ in‑memory data store ที่ใช้เป็น cache, session store, message broker และ distributed lock มีความเร็วสูงมาก (microsecond latency) และรองรับ data structure หลากหลาย เช่น String, Hash, List, Set, Sorted Set, Stream

### มีกี่แบบ? (Use Cases)
1. **Cache Repository** – เก็บข้อมูลที่อ่านบ่อย (user profile, sensor config) เพื่อลดภาระ PostgreSQL
2. **Refresh Token Store** – เก็บ refresh token แบบมี TTL (7 วัน) รองรับการ revoke ทันที
3. **Token Blacklist** – เก็บ JWT ที่ถูก logout หรือ revoke จนกว่าจะหมดอายุ
4. **Distributed Lock** – ป้องกัน cron job หรือ operation ทำงานซ้ำเมื่อมีหลาย instance
5. **Pub/Sub** – ส่งข้อความระหว่าง services สำหรับ real‑time notification (WebSocket broadcast)
6. **Rate Limiting** – จำกัดจำนวน request ต่อ IP หรือ user (Token Bucket หรือ Sliding Window)

### ใช้อย่างไร / นำไปใช้กรณีไหน
- **Cache**: ใช้ `cache.Set()` และ `cache.Get()` พร้อม TTL (เช่น 10 นาที)
- **Refresh token**: ใช้ `refreshStore.Create()` เมื่อ login, `refreshStore.Get()` เมื่อ refresh, `refreshStore.Delete()` เมื่อ logout
- **Blacklist**: ใช้ `blacklist.Add()` เมื่อ logout, `blacklist.IsBlacklisted()` ใน middleware
- **Lock**: ใช้ `lock.Acquire()` ก่อนทำงาน critical section, `lock.Release()` หลังเสร็จ

### ทำไมต้องใช้ Redis
- **ความเร็วสูง** – in‑memory, ตอบสนองในระดับ microsecond (เร็วกว่า PostgreSQL 100‑1000 เท่า)
- **ลดภาระฐานข้อมูลหลัก** – แคชข้อมูลที่เรียกบ่อย
- **TTL อัตโนมัติ** – ไม่ต้องเขียน cleanup script
- **Atomic operations** – INCR, SETNX, ฯลฯ เหมาะกับ distributed lock และ rate limiting
- **Pub/Sub** – ง่ายต่อการ broadcast ข้อความระหว่างหลาย instances

### ประโยชน์ที่ได้รับ
- Response time ลดลงอย่างมาก
- รองรับ high concurrency (Redis เป็น single‑threaded, แต่ non‑blocking I/O)
- Horizontal scaling ง่าย (Redis Cluster, Sentinel)
- ลด load บน PostgreSQL ทำให้ database รองรับการเขียนได้ดีขึ้น

### ข้อควรระวัง
- Redis เป็น in‑memory → ถ้า restart โดยไม่ตั้งค่า persistence จะสูญเสียข้อมูล (แต่ cache และ refresh token สามารถ rebuild ได้)
- ต้องตั้งค่า `maxmemory` และ `maxmemory-policy` (เช่น `allkeys-lru`) เพื่อป้องกัน OOM
- ระวัง cache stampede effect (หลาย request พร้อมกันที่ cache miss) – ใช้ `singleflight` หรือ distributed lock ป้องกัน
- Refresh token ที่เก็บใน Redis หาก Redis crash ผู้ใช้จะต้อง login ใหม่ (เป็นที่ยอมรับได้)

### ข้อดี
- เร็วมาก, รองรับ data structures หลากหลาย, atomic operations, TTL, Pub/Sub

### ข้อเสีย
- ข้อมูลไม่ถาวร (ถ้าไม่ตั้งค่า AOF / RDB)
- ใช้ memory (ต้อง monitor)
- Complexity เพิ่มขึ้น (ต้องดูแลอีกหนึ่ง component)

### ข้อห้าม
- ห้ามใช้ Redis เป็น primary database สำหรับข้อมูลสำคัญ (ต้องมี persistence และ backup)
- ห้ามเก็บข้อมูลใหญ่เกินไป (เช่น blob, ไฟล์) เพราะจะกิน memory
- ห้ามใช้ KEYS command ใน production (ใช้ SCAN แทน)
- ห้ามใช้ FLUSHALL บน production โดยไม่คิด

---

## โค้ดที่รันได้จริง (ครบทุกไฟล์)

### 1. Redis Client – `client.go`

```go
// Package redis provides Redis client and utilities for caching, sessions, and distributed locks.
// ----------------------------------------------------------------
// แพ็คเกจ redis ให้บริการ Redis client และ utilities สำหรับ cache, session, และ distributed lock
package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config holds Redis connection settings.
// ----------------------------------------------------------------
// Config เก็บค่ากำหนดการเชื่อมต่อ Redis
type Config struct {
	Addr         string        // host:port, e.g., "localhost:6379"
	Password     string        // password (if any)
	DB           int           // database number (0-15)
	PoolSize     int           // maximum number of socket connections
	MinIdleConns int           // minimum number of idle connections
	DialTimeout  time.Duration // timeout for establishing connection
	ReadTimeout  time.Duration // timeout for read operations
	WriteTimeout time.Duration // timeout for write operations
}

// NewClient creates a new Redis client with connection pool.
// ----------------------------------------------------------------
// NewClient สร้าง Redis client ใหม่พร้อม connection pool
func NewClient(cfg Config) (*redis.Client, error) {
	opt := &redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
	// Apply defaults if not set
	// ใช้ค่าเริ่มต้นถ้ายังไม่ได้ตั้งค่า
	if opt.PoolSize == 0 {
		opt.PoolSize = 10
	}
	if opt.MinIdleConns == 0 {
		opt.MinIdleConns = 2
	}
	if opt.DialTimeout == 0 {
		opt.DialTimeout = 5 * time.Second
	}
	if opt.ReadTimeout == 0 {
		opt.ReadTimeout = 3 * time.Second
	}
	if opt.WriteTimeout == 0 {
		opt.WriteTimeout = 3 * time.Second
	}

	client := redis.NewClient(opt)

	// Verify connection
	// ตรวจสอบการเชื่อมต่อ
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return client, nil
}
```

### 2. Cache Repository – `cache.go`

```go
package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheRepository defines generic cache operations.
// ----------------------------------------------------------------
// CacheRepository กำหนดการดำเนินการ cache ทั่วไป
type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Increment(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error
}

// RedisCache implements CacheRepository using go-redis.
// ----------------------------------------------------------------
// RedisCache อิมพลีเมนต์ CacheRepository ด้วย go-redis
type RedisCache struct {
	client *redis.Client
}

// NewCache creates a new cache repository.
// ----------------------------------------------------------------
// NewCache สร้าง cache repository ใหม่
func NewCache(client *redis.Client) CacheRepository {
	return &RedisCache{client: client}
}

// Set stores a value with TTL (marshal to JSON automatically).
// ----------------------------------------------------------------
// Set เก็บค่าใน Redis พร้อม TTL (แปลงเป็น JSON อัตโนมัติ)
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

// Get retrieves and unmarshals a value from Redis.
// Returns nil error but dest unchanged if key not found.
// ----------------------------------------------------------------
// Get ดึงและ unmarshal ค่าจาก Redis
// คืน error เป็น nil แต่ dest ไม่เปลี่ยนถ้าไม่พบ key
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil // not found, ไม่พบ
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Delete removes a key.
// ----------------------------------------------------------------
// Delete ลบคีย์
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Exists checks if key exists.
// ----------------------------------------------------------------
// Exists ตรวจสอบว่าคีย์มีอยู่หรือไม่
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	return n > 0, err
}

// Increment atomically increments a counter.
// ----------------------------------------------------------------
// Increment เพิ่มค่า counter แบบอะตอม (ใช้สำหรับ rate limiting)
func (c *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

// Expire sets TTL on an existing key.
// ----------------------------------------------------------------
// Expire กำหนด TTL ให้กับคีย์ที่มีอยู่
func (c *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.client.Expire(ctx, key, ttl).Err()
}
```

### 3. Refresh Token Store – `refresh_store.go`

```go
package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// RefreshSession represents data stored with a refresh token.
// ----------------------------------------------------------------
// RefreshSession แทนข้อมูลที่เก็บร่วมกับ refresh token
type RefreshSession struct {
	UserID    uint      `json:"user_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RefreshStore manages refresh tokens in Redis.
// ----------------------------------------------------------------
// RefreshStore จัดการ refresh tokens ใน Redis
type RefreshStore interface {
	Create(ctx context.Context, token string, session *RefreshSession) error
	Get(ctx context.Context, token string) (*RefreshSession, error)
	Delete(ctx context.Context, token string) error
}

// redisRefreshStore implements RefreshStore.
// ----------------------------------------------------------------
// redisRefreshStore อิมพลีเมนต์ RefreshStore
type redisRefreshStore struct {
	client *redis.Client
}

// NewRefreshStore creates a new refresh token store.
// ----------------------------------------------------------------
// NewRefreshStore สร้าง refresh token store ใหม่
func NewRefreshStore(client *redis.Client) RefreshStore {
	return &redisRefreshStore{client: client}
}

// Create stores a refresh token with TTL.
// Key format: "refresh:{token}"
// ----------------------------------------------------------------
// Create เก็บ refresh token พร้อม TTL
func (s *redisRefreshStore) Create(ctx context.Context, token string, session *RefreshSession) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	ttl := time.Until(session.ExpiresAt)
	return s.client.Set(ctx, "refresh:"+token, data, ttl).Err()
}

// Get retrieves refresh session by token.
// ----------------------------------------------------------------
// Get ดึง refresh session ด้วย token
func (s *redisRefreshStore) Get(ctx context.Context, token string) (*RefreshSession, error) {
	data, err := s.client.Get(ctx, "refresh:"+token).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var session RefreshSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// Delete removes refresh token from Redis.
// ----------------------------------------------------------------
// Delete ลบ refresh token ออกจาก Redis
func (s *redisRefreshStore) Delete(ctx context.Context, token string) error {
	return s.client.Del(ctx, "refresh:"+token).Err()
}
```

### 4. Token Blacklist – `blacklist.go`

```go
package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Blacklist manages revoked JWT tokens.
// ----------------------------------------------------------------
// Blacklist จัดการ JWT tokens ที่ถูกเพิกถอน
type Blacklist interface {
	Add(ctx context.Context, jti string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, jti string) (bool, error)
}

// redisBlacklist implements Blacklist.
// ----------------------------------------------------------------
// redisBlacklist อิมพลีเมนต์ Blacklist
type redisBlacklist struct {
	client *redis.Client
}

// NewBlacklist creates a new blacklist.
// ----------------------------------------------------------------
// NewBlacklist สร้าง blacklist ใหม่
func NewBlacklist(client *redis.Client) Blacklist {
	return &redisBlacklist{client: client}
}

// Add adds a JWT ID (jti) to blacklist with TTL equal to remaining token expiry.
// ----------------------------------------------------------------
// Add เพิ่ม JWT ID (jti) ลง blacklist พร้อม TTL เท่ากับเวลาที่เหลือของ token
func (b *redisBlacklist) Add(ctx context.Context, jti string, ttl time.Duration) error {
	return b.client.Set(ctx, "blacklist:"+jti, "1", ttl).Err()
}

// IsBlacklisted checks if JWT ID is blacklisted.
// ----------------------------------------------------------------
// IsBlacklisted ตรวจสอบว่า JWT ID ถูก blacklist หรือไม่
func (b *redisBlacklist) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	_, err := b.client.Get(ctx, "blacklist:"+jti).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
```

### 5. Distributed Lock – `lock.go`

```go
package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// DistributedLock provides a simple Redis-based distributed lock.
// Uses SET NX EX to acquire, and DEL to release.
// ----------------------------------------------------------------
// DistributedLock ให้บริการ distributed lock บน Redis
// ใช้ SET NX EX เพื่อขอ lock, และ DEL เพื่อปล่อย lock
type DistributedLock struct {
	client *redis.Client
}

// NewDistributedLock creates a new lock instance.
// ----------------------------------------------------------------
// NewDistributedLock สร้าง distributed lock ใหม่
func NewDistributedLock(client *redis.Client) *DistributedLock {
	return &DistributedLock{client: client}
}

// Acquire tries to acquire a lock with given TTL.
// Returns true if lock acquired, false otherwise.
// ----------------------------------------------------------------
// Acquire พยายามขอ lock ด้วย TTL ที่กำหนด
// คืน true ถ้าได้ lock, false ถ้าไม่ได้
func (l *DistributedLock) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	ok, err := l.client.SetNX(ctx, "lock:"+key, "locked", ttl).Result()
	return ok, err
}

// Release releases the lock.
// ----------------------------------------------------------------
// Release ปล่อย lock
func (l *DistributedLock) Release(ctx context.Context, key string) error {
	return l.client.Del(ctx, "lock:"+key).Err()
}

// WithLock executes a function while holding a lock.
// Automatically acquires the lock and releases it after fn finishes.
// ----------------------------------------------------------------
// WithLock ทำงานฟังก์ชันโดยถือ lock อยู่
// ขอ lock อัตโนมัติและปล่อยเมื่อ fn ทำงานเสร็จ
func (l *DistributedLock) WithLock(ctx context.Context, key string, ttl time.Duration, fn func() error) error {
	acquired, err := l.Acquire(ctx, key, ttl)
	if err != nil {
		return err
	}
	if !acquired {
		return nil // skip, ไม่ได้ lock
	}
	defer l.Release(ctx, key)
	return fn()
}
```

---

## วิธีใช้งาน module นี้

### 1. ติดตั้ง dependency
```bash
go get github.com/redis/go-redis/v9
```

### 2. วางไฟล์ทั้งหมดใน `internal/pkg/redis/`

### 3. สร้าง client และ repositories ใน `main.go`
```go
import "gobackend/internal/pkg/redis"

func main() {
    // Connect to Redis
    redisClient, err := redis.NewClient(redis.Config{
        Addr:     os.Getenv("REDIS_ADDR"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0,
        PoolSize: 20,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer redisClient.Close()

    // Create repositories
    cacheRepo := redis.NewCache(redisClient)
    refreshStore := redis.NewRefreshStore(redisClient)
    blacklist := redis.NewBlacklist(redisClient)
    lock := redis.NewDistributedLock(redisClient)

    // Inject into usecases or handlers...
}
```

### 4. ตัวอย่างการใช้งานจริง

**Cache user profile (ใน user_usecase.go):**
```go
func (u *userUsecase) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)
    var user models.User
    err := u.cacheRepo.Get(ctx, cacheKey, &user)
    if err == nil {
        return &user, nil
    }
    // cache miss
    userPtr, err := u.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    u.cacheRepo.Set(ctx, cacheKey, userPtr, 10*time.Minute)
    return userPtr, nil
}
```

**Refresh token management (ใน auth_usecase.go):**
```go
// Login: create refresh token
refreshToken := uuid.New().String()
session := &redis.RefreshSession{
    UserID:    user.ID,
    Role:      string(user.Role),
    CreatedAt: time.Now(),
    ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
}
refreshStore.Create(ctx, refreshToken, session)

// Refresh: validate refresh token
session, err := refreshStore.Get(ctx, refreshToken)
if err != nil || session == nil {
    return "", ErrInvalidRefreshToken
}
```

**Distributed lock for cron job (ใน scheduler):**
```go
err := lock.WithLock(ctx, "job:daily_report", 5*time.Minute, func() error {
    return generateAndSendReport()
})
if err != nil {
    log.Println("failed to acquire lock, skipping")
}
```

---

## ตารางสรุป Redis Key Patterns

| Purpose | Key Pattern | TTL | Example |
|---------|-------------|-----|---------|
| User cache | `user:{user_id}` | 10 min | `user:123` |
| Refresh token | `refresh:{token}` | 7 days | `refresh:abc-123` |
| Blacklist | `blacklist:{jti}` | token expiry | `blacklist:abc-123` |
| Distributed lock | `lock:{resource}` | task duration | `lock:job:cleanup` |
| Rate limit | `rate:{ip}:{endpoint}` | 1 min | `rate:192.168.1.1:login` |

---

## แบบฝึกหัดท้าย module (3 ข้อ)

1. **เพิ่ม method `GetMultiple`** ใน `CacheRepository` ที่รับหลาย keys และคืน map ของ results โดยใช้ MGET เพื่อลด round‑trip
2. **Implement `RefreshStore.Rotate`** ที่สร้าง refresh token ใหม่และลบอันเก่าใน atomic operation (ใช้ pipeline หรือ Lua script)
3. **ปรับปรุง `DistributedLock`** ให้มี `Extend` method ที่เพิ่ม TTL ให้ lock ที่ถืออยู่ (ป้องกันงานที่ใช้เวลานานเกิน TTL)

---

## แหล่งอ้างอิง

- [go-redis documentation](https://redis.uptrace.dev/)
- [Redis patterns: caching](https://redis.io/docs/latest/develop/use/patterns/caching/)
- [Distributed locks with Redis](https://redis.io/docs/latest/develop/use/patterns/distributed-locks/)
- [Redis Pub/Sub](https://redis.io/docs/latest/develop/interact/pubsub/)

---

**หมายเหตุ:** module นี้ครบถ้วนสำหรับ `pkg/redis` หากต้องการ module อื่น (เช่น `pkg/email`, `pkg/jwt`, `pkg/hash`) โปรดแจ้ง หรือจะขอให้สรุปเนื้อหาทั้งหมดของระบบ gobackend ก็ได้