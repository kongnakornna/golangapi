# Module 10: pkg/redis (Redis Client & Utilities)

## สำหรับโฟลเดอร์ `internal/pkg/redis/`

ไฟล์ที่เกี่ยวข้อง:
- `internal/pkg/redis/client.go`
- `internal/pkg/redis/cache.go`
- `internal/pkg/redis/refresh_store.go`
- `internal/pkg/redis/blacklist.go`
- `internal/pkg/redis/lock.go`

---

## หลักการ (Concept)

### คืออะไร?
Redis เป็น in-memory data store ที่ใช้เป็น cache, session store, message broker, และ distributed lock ในระบบ backend ช่วยเพิ่มประสิทธิภาพและลดภาระของฐานข้อมูลหลัก

### มีกี่แบบ?
1. **Cache repository** – เก็บข้อมูลที่อ่านบ่อย (user profile, config)
2. **Refresh token store** – เก็บ refresh token แบบมี TTL
3. **Token blacklist** – เก็บ JWT ที่ถูก revoke จนกว่าจะหมดอายุ
4. **Distributed lock** – ป้องกันงานซ้ำใน multi-instance environment
5. **Pub/Sub** – ส่งข้อความระหว่าง services (real-time notifications)

### ใช้อย่างไร / นำไปใช้กรณีไหน
- Cache: ใช้สำหรับข้อมูลที่ไม่เปลี่ยนแปลงบ่อย (user, sensor config)
- Session: เก็บ refresh token, user session
- Blacklist: เก็บ JWT ที่ logout แล้ว
- Lock: ป้องกัน cron job ทำงานซ้ำ

### ทำไมต้องใช้
- Redis เร็วกว่า PostgreSQL 10-100 เท่า
- ลด load บนฐานข้อมูลหลัก
- รองรับการ distributed environment

### ประโยชน์ที่ได้รับ
- response time ลดลง
- scalability ดีขึ้น
- มี data structures หลากหลาย (string, hash, list, set, sorted set)

### ข้อควรระวัง
- Redis เป็น in-memory → ข้อมูลหายได้ถ้าไม่ตั้งค่า persistence
- ต้องกำหนด maxmemory และ eviction policy
- ระวัง cache stampede effect

### ข้อดี
- เร็วมาก, รองรับ distributed lock, TTL อัตโนมัติ

### ข้อเสีย
- ข้อมูลไม่ถาวร (ต้อง backup)
- เพิ่ม complexity อีก一层

### ข้อห้าม
- ห้ามใช้ Redis เป็น primary database สำหรับข้อมูลสำคัญ
- ห้ามเก็บข้อมูลใหญ่เกินไป (อาจ overflow memory)

---

## โค้ดที่รันได้จริง

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
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// NewClient creates a new Redis client with default config.
// ----------------------------------------------------------------
// NewClient สร้าง Redis client ใหม่พร้อมค่ากำหนดเริ่มต้น
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

// Set stores a value with TTL.
// ----------------------------------------------------------------
// Set เก็บค่าใน Redis พร้อม TTL
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

// Get retrieves and unmarshals a value.
// ----------------------------------------------------------------
// Get ดึงและ unmarshal ค่าจาก Redis
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil // not found
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
// Increment เพิ่มค่า counter แบบอะตอม
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

// RefreshSession represents refresh token data.
// ----------------------------------------------------------------
// RefreshSession แทนข้อมูล refresh token
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

// Delete removes refresh token.
// ----------------------------------------------------------------
// Delete ลบ refresh token
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

// Add adds a JWT ID to blacklist with TTL equal to remaining token expiry.
// ----------------------------------------------------------------
// Add เพิ่ม JWT ID ลง blacklist พร้อม TTL เท่ากับเวลาที่เหลือของ token
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
// ----------------------------------------------------------------
// DistributedLock ให้บริการ distributed lock บน Redis
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
// ----------------------------------------------------------------
// WithLock ทำงานฟังก์ชันโดยถือ lock อยู่
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

### 6. ตัวอย่างการใช้งานรวม – `example_test.go`

```go
package redis_test

import (
	"context"
	"time"
	"gobackend/internal/pkg/redis"
)

func Example() {
	client, _ := redis.NewClient(redis.Config{
		Addr: "localhost:6379",
	})
	defer client.Close()
	
	ctx := context.Background()
	
	// Cache
	cache := redis.NewCache(client)
	cache.Set(ctx, "user:1", map[string]string{"name": "John"}, 5*time.Minute)
	
	var user map[string]string
	cache.Get(ctx, "user:1", &user)
	
	// Refresh store
	refreshStore := redis.NewRefreshStore(client)
	session := &redis.RefreshSession{UserID: 1, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)}
	refreshStore.Create(ctx, "token123", session)
	
	// Blacklist
	blacklist := redis.NewBlacklist(client)
	blacklist.Add(ctx, "jti123", 15*time.Minute)
	
	// Lock
	lock := redis.NewDistributedLock(client)
	acquired, _ := lock.Acquire(ctx, "job:cleanup", 30*time.Second)
	if acquired {
		// do work
		lock.Release(ctx, "job:cleanup")
	}
}
```

---

## วิธีใช้งาน module นี้

1. ติดตั้ง dependency:
   ```bash
   go get github.com/redis/go-redis/v9
   ```
2. วางไฟล์ทั้งหมดใน `internal/pkg/redis/`
3. สร้าง client ใน `main.go`:
   ```go
   redisClient, err := redis.NewClient(redis.Config{
       Addr: os.Getenv("REDIS_ADDR"),
       Password: os.Getenv("REDIS_PASSWORD"),
       DB: 0,
   })
   ```
4. สร้าง repositories:
   ```go
   cacheRepo := redis.NewCache(redisClient)
   refreshStore := redis.NewRefreshStore(redisClient)
   blacklist := redis.NewBlacklist(redisClient)
   lock := redis.NewDistributedLock(redisClient)
   ```
5. Inject เข้า usecase หรือ handler

---

## ตารางสรุป Utilities

| ชื่อ | ฟังก์ชันหลัก | การใช้งาน |
|------|-------------|-----------|
| `CacheRepository` | Set, Get, Delete, Exists | แคชข้อมูลผู้ใช้, ค่า config |
| `RefreshStore` | Create, Get, Delete | จัดการ refresh token |
| `Blacklist` | Add, IsBlacklisted | เพิกถอน JWT |
| `DistributedLock` | Acquire, Release, WithLock | ป้องกันงานซ้ำ |

---

## แบบฝึกหัดท้าย module (3 ข้อ)

1. เพิ่ม method `GetMultiple` ใน `CacheRepository` ที่รับหลาย keys และคืน map ของ results โดยใช้ MGET
2. Implement `RefreshStore.Rotate` ที่สร้าง refresh token ใหม่และลบอันเก่าใน atomic operation (ใช้ pipeline)
3. ปรับปรุง `DistributedLock` ให้มี `Extend` method ที่เพิ่ม TTL ให้ lock ที่ถืออยู่ (ป้องกันงานที่ใช้เวลานาน)

---

## แหล่งอ้างอิง

- [go-redis documentation](https://redis.uptrace.dev/)
- [Redis patterns: caching](https://redis.io/docs/latest/develop/use/patterns/caching/)
- [Distributed locks with Redis](https://redis.io/docs/latest/develop/use/patterns/distributed-locks/)

---

**หมายเหตุ:** module นี้ครบถ้วนสำหรับ `pkg/redis` หากต้องการ module ถัดไป (เช่น `pkg/email`, `pkg/jwt`, `pkg/hash`) โปรดแจ้ง