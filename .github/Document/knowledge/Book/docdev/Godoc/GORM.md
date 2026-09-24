# Golng GORM ORM คืออะไร  
------------------------
1. GORM CRUD
2. Funtion sessionFactory 
3. Cache Query Cache 
  3.1. Cache option 1.quit delete key 2.option time lite
4. Query Cache Queue Processor:
5. Queue Processor Funtion
6. SQL queue translate rollback commit
7. kafka message queue
8. go live monitoring
------------------------
1.สร้างบทนำ
2.สร้างบทนิยาม
3.สร้างบทหัวข้อ
5.ออกแบบคู่มือ
6.ออกแบบ workflow
7.code template เพือนำไปใช้งาน
------------------------
## บทที่ 48: GORM – ORM ทรงพลังสำหรับ Go

### 48.1 บทนำ

การจัดการฐานข้อมูลเป็นหัวใจสำคัญของแอปพลิเคชันส่วนใหญ่ การเขียน SQL ด้วยมือนั้นให้ความยืดหยุ่นสูง แต่ก็อาจทำให้โค้ดยุ่งเหยิง เสี่ยงต่อข้อผิดพลาด และต้องดูแลการแปลงข้อมูลระหว่าง Go struct กับตารางด้วยตนเอง GORM (Go Object Relational Mapping) เป็น ORM ที่ได้รับความนิยมสูงสุดในภาษา Go ช่วยลดความซับซ้อนเหล่านี้ ด้วยการแมป struct กับตารางโดยอัตโนมัติ พร้อมฟังก์ชัน CRUD ที่ใช้งานง่าย การจัดการ transaction แบบมีระบบ และยังมีฟีเจอร์ขั้นสูงอย่าง query cache, queue processor ที่ช่วยให้แอปพลิเคชันทำงานได้เร็วและมั่นคงยิ่งขึ้น

บทนี้จะพาคุณสำรวจ GORM ตั้งแต่พื้นฐาน CRUD ไปจนถึงการออกแบบ session factory, การใช้ query cache, การประมวลผลคิว, และการจัดการ transaction แบบ rollback/commit อย่างเป็นระบบ พร้อมตัวอย่างโค้ดและ workflow ที่นำไปใช้ได้จริง

---

### 48.2 บทนิยาม

#### 48.2.1 CRUD
**CRUD** คือชุดการดำเนินการพื้นฐานสี่ประการในการจัดการข้อมูล:
- **C**reate – เพิ่มข้อมูลใหม่
- **R**ead – อ่าน/ดึงข้อมูล
- **U**pdate – แก้ไขข้อมูลที่มีอยู่
- **D**elete – ลบข้อมูล

ใน GORM การดำเนินการเหล่านี้ทำผ่าน `*gorm.DB` object และ method ต่างๆ เช่น `Create`, `First`, `Find`, `Save`, `Update`, `Delete`

#### 48.2.2 ORM
**ORM (Object-Relational Mapping)** คือเทคนิคที่ช่วยแปลงข้อมูลระหว่างฐานข้อมูลเชิงสัมพันธ์กับโครงสร้างข้อมูลเชิงวัตถุ (object) ในภาษาโปรแกรม โดย GORM จะทำหน้าที่:
- แปลง Go struct ให้เป็น SQL และแปลงผลลัพธ์ SQL กลับเป็น Go struct
- จัดการความสัมพันธ์ระหว่างตาราง (has one, has many, belongs to, many to many)
- รองรับ migration, transaction, hooks, และ caching

#### 48.2.3 GORM
**GORM** คือ ORM library สำหรับ Go ที่ออกแบบมาให้ใช้งานง่าย กระชับ และมีประสิทธิภาพสูง โดยมีคุณสมบัติเด่น:
- Full-featured ORM (CRUD, associations, hooks, transactions)
- ใช้งานได้กับฐานข้อมูลหลายชนิด (MySQL, PostgreSQL, SQLite, SQL Server, ClickHouse)
- มี chainable API ที่อ่านง่าย
- รองรับ eager loading (Preload)
- มี auto-migration
- สามารถขยายฟังก์ชันผ่าน plugins

#### 48.2.4 Transaction
**Transaction (ทรานแซคชัน)** คือกลุ่มของคำสั่ง SQL ที่ต้องสำเร็จทั้งหมดหรือไม่ทำเลย (All-or-Nothing) โดยมีคุณสมบัติ ACID:
- **Atomicity**: งานทั้งหมดสำเร็จหรือล้มเหลวพร้อมกัน
- **Consistency**: ข้อมูลคงความถูกต้องตามกฎธุรกิจ
- **Isolation**: ทรานแซคชันที่ทำงานพร้อมกันไม่รบกวนกัน
- **Durability**: เมื่อ commit ข้อมูลจะถูกบันทึกถาวร

GORM รองรับ transaction ผ่าน `db.Transaction()` หรือ `db.Begin()` / `Commit()` / `Rollback()`

#### 48.2.5 Cache
**Cache (แคช)** คือการเก็บข้อมูลชั่วคราวเพื่อลดการเข้าถึงฐานข้อมูลซ้ำๆ เพิ่มความเร็วในการตอบสนอง GORM มีกลไก:
- **First-level cache**: ภายใน session เดียว (ไม่เปิดเผยให้ผู้ใช้จัดการ)
- **Second-level cache**: ระดับ application, ต้องใช้ plugin เช่น `gorm-cache`
- **Query cache**: เก็บผลลัพธ์ของคำสั่ง `Find`, `First` เพื่อใช้ซ้ำ

#### 48.2.6 Queue Processor
**Queue Processor** เป็นกลไกในการประมวลผลงานที่อาจใช้เวลานานหรือเกิดบ่อยแบบ asynchronous โดย GORM สามารถทำงานร่วมกับระบบคิว (เช่น Redis, RabbitMQ) เพื่อแยกการรับ request ออกจากการทำงานหนัก ทำให้ระบบตอบสนองได้เร็วขึ้น

---

### 48.3 หัวข้อหลัก

1. **GORM CRUD** – การสร้าง, อ่าน, อัปเดต, ลบข้อมูลพื้นฐาน
2. **SessionFactory** – การสร้างและจัดการ `*gorm.DB` session สำหรับแยก context
3. **Query Cache** – การแคชผลลัพธ์ query เพื่อลดภาระฐานข้อมูล
4. **Cache Query Queue Processor** – การใช้คิวเพื่อประมวลผล query แบบ asynchronous และแคช
5. **Queue Processor Function** – การออกแบบฟังก์ชันสำหรับประมวลผลคิว
6. **SQL Queue Translate Rollback Commit** – การแปลง SQL เป็นคิวและจัดการ transaction ในระบบคิว

---

### 48.4 คู่มือการใช้งาน GORM

#### 48.4.1 การติดตั้งและการตั้งค่าพื้นฐาน

```go
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres // หรือ driver อื่นตามที่ใช้
```

**ตัวอย่างการเชื่อมต่อ PostgreSQL**

```go
package main

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
)

func main() {
    dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Bangkok"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect database")
    }

    // ใช้ db ต่อไป
}
```

#### 48.4.2 กำหนด Model (Entity)

```go
type User struct {
    ID        uint           `gorm:"primaryKey"`
    Name      string         `gorm:"size:100;not null"`
    Email     string         `gorm:"uniqueIndex;size:100;not null"`
    Age       int            `gorm:"default:0"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

#### 48.4.3 CRUD Operations

**Create**

```go
user := User{Name: "สมชาย", Email: "somchai@example.com", Age: 30}
result := db.Create(&user) // สร้าง record ใหม่
fmt.Println(user.ID)       // คืนค่า ID ที่ถูกสร้าง
fmt.Println(result.Error)  // error ถ้ามี
```

**Read**

```go
// ดึง record แรกที่ตรงเงื่อนไข
var user User
db.First(&user, 1)                 // by primary key
db.First(&user, "email = ?", "somchai@example.com")

// ดึงทั้งหมด
var users []User
db.Find(&users)

// พร้อมเงื่อนไข
db.Where("age > ?", 20).Find(&users)
db.Where(&User{Name: "สมชาย"}).Find(&users)
```

**Update**

```go
// อัปเดต single column
db.Model(&user).Update("Name", "สมชาย ใหม่")

// อัปเดตหลาย columns ด้วย struct (ไม่สนใจ zero values)
db.Model(&user).Updates(User{Name: "สมชาย ใหม่", Age: 31})

// อัปเดตหลาย columns ด้วย map
db.Model(&user).Updates(map[string]interface{}{"name": "สมชาย ใหม่", "age": 31})
```

**Delete**

```go
// soft delete (ถ้ามี gorm.DeletedAt)
db.Delete(&user, 1)

// hard delete
db.Unscoped().Delete(&user, 1)
```

#### 48.4.4 SessionFactory

Session factory เป็นรูปแบบการสร้าง `*gorm.DB` ที่มี configuration คงที่ (เช่น logging, connection pool) และสามารถสร้าง session ใหม่สำหรับแต่ละ request หรือ transaction

**ตัวอย่าง session factory**

```go
package db

import (
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "log"
    "time"
)

type SessionFactory struct {
    db *gorm.DB
}

func NewSessionFactory(dsn string) (*SessionFactory, error) {
    // ตั้งค่า logger และ connection pool
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold: time.Second,
            LogLevel:      logger.Info,
            Colorful:      true,
        },
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: newLogger,
        NowFunc: func() time.Time { return time.Now().UTC() },
    })
    if err != nil {
        return nil, err
    }

    // ตั้งค่า connection pool
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    return &SessionFactory{db: db}, nil
}

// NewSession คืนค่า session ใหม่สำหรับการทำ transaction หรือ query แบบแยก context
func (sf *SessionFactory) NewSession() *gorm.DB {
    return sf.db.Session(&gorm.Session{})
}
```

**การใช้งาน**

```go
factory, _ := db.NewSessionFactory(dsn)

// สร้าง session ใหม่สำหรับ request นี้
session := factory.NewSession()
var user User
session.First(&user, 1)

// เมื่อต้องการ transaction
session.Transaction(func(tx *gorm.DB) error {
    // ใช้ tx แทน session
    return nil
})
```

#### 48.4.5 Query Cache

GORM เองไม่มี query cache ในตัว แต่สามารถใช้ plugin หรือจัดการเองผ่าน Redis หรือ memory cache

**ตัวอย่างการใช้ Redis cache กับ GORM (แบบง่าย)**

```go
import (
    "context"
    "encoding/json"
    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"
    "time"
)

type CachedDB struct {
    db    *gorm.DB
    cache *redis.Client
}

func (c *CachedDB) FirstWithCache(dest interface{}, conds ...interface{}) error {
    // สร้าง cache key จากเงื่อนไข
    key := fmt.Sprintf("query:%v", conds)

    // พยายามอ่านจาก cache
    val, err := c.cache.Get(context.Background(), key).Result()
    if err == nil {
        // พบใน cache
        return json.Unmarshal([]byte(val), dest)
    }

    // ไม่พบใน cache, query ฐานข้อมูล
    if err := c.db.First(dest, conds...).Error; err != nil {
        return err
    }

    // บันทึกผลลัพธ์ลง cache (serialize)
    data, _ := json.Marshal(dest)
    c.cache.Set(context.Background(), key, data, 5*time.Minute)
    return nil
}
```

#### 48.4.6 Cache Query Queue Processor

เมื่อมี query จำนวนมากที่ต้องการ cache แบบ asynchronous (เช่น การ pre-cache ข้อมูลที่คาดว่าจะถูกเรียกบ่อย) เราสามารถใช้ queue processor ในการรับ query, ประมวลผล, และเก็บผลลัพธ์ลง cache

**โครงสร้าง**

```
[Client Request] → [API] → (อ่าน cache) → ถ้ามี → ส่งกลับ
                     ↓ ถ้าไม่มี
              [Queue (Redis)] → [Worker] → Query DB → Store Cache → Notify
```

**ตัวอย่างการใช้ Redis เป็น queue และ worker**

```go
// โครงสร้างงาน
type CacheJob struct {
    QueryKey  string        `json:"query_key"`
    TableName string        `json:"table_name"`
    Conditions []interface{} `json:"conditions"`
    TTL       time.Duration `json:"ttl"`
}

// Producer: ส่งงานเข้าระบบ queue
func (c *CachedDB) QueueQuery(ctx context.Context, job CacheJob) error {
    data, _ := json.Marshal(job)
    return c.redis.LPush(ctx, "cache_queue", data).Err()
}

// Worker: รับงานจาก queue และประมวลผล
func (c *CachedDB) StartCacheWorker(ctx context.Context) {
    for {
        result, err := c.redis.BRPop(ctx, 0, "cache_queue").Result()
        if err != nil {
            continue
        }
        var job CacheJob
        json.Unmarshal([]byte(result[1]), &job)

        // Query ฐานข้อมูลตามเงื่อนไข
        var dest interface{}
        // สมมติว่าทราบ type จาก table name
        switch job.TableName {
        case "users":
            var users []User
            c.db.Where(job.Conditions...).Find(&users)
            dest = users
        // ...
        }

        // เก็บลง cache
        data, _ := json.Marshal(dest)
        c.redis.Set(ctx, job.QueryKey, data, job.TTL)
    }
}
```

#### 48.4.7 Queue Processor Function

ฟังก์ชันสำหรับประมวลผลคิวควรถูกออกแบบให้ทำงานแบบ concurrent, มีการ retry, และจัดการ error อย่างเหมาะสม

**เทมเพลต worker**

```go
type WorkerPool struct {
    workers int
    jobs    chan CacheJob
    wg      sync.WaitGroup
    db      *gorm.DB
    redis   *redis.Client
}

func NewWorkerPool(workers int, db *gorm.DB, redis *redis.Client) *WorkerPool {
    return &WorkerPool{
        workers: workers,
        jobs:    make(chan CacheJob, 100),
        db:      db,
        redis:   redis,
    }
}

func (wp *WorkerPool) Start(ctx context.Context) {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker(ctx)
    }
}

func (wp *WorkerPool) worker(ctx context.Context) {
    defer wp.wg.Done()
    for {
        select {
        case job := <-wp.jobs:
            wp.processJob(job)
        case <-ctx.Done():
            return
        }
    }
}

func (wp *WorkerPool) processJob(job CacheJob) {
    // implement query และ cache logic
    // มี retry และ error logging
}
```

#### 48.4.8 SQL Queue Translate Rollback Commit

ในระบบที่มีการประมวลผลแบบคิว งานบางอย่างอาจเกี่ยวข้องกับการเปลี่ยนแปลงฐานข้อมูล (เช่น การอัปเดตสถานะ) ซึ่งต้องมี transaction เพื่อรักษาความถูกต้อง หลักการคือ:

- **Translate SQL to Queue**: แทนที่จะ execute SQL ทันที ให้สร้าง job ที่มีข้อมูลเพียงพอที่จะ execute SQL ในภายหลัง
- **Rollback**: ถ้า job ถูกประมวลผลไม่สำเร็จหลัง commit ไปแล้ว อาจต้องมีการชดเชย (compensating action) เนื่องจากฐานข้อมูลไม่สามารถ rollback ได้ง่ายเมื่อ transaction จบแล้ว
- **Commit**: เมื่อ job ทำงานสำเร็จ ให้ commit การเปลี่ยนแปลงในฐานข้อมูล

**ตัวอย่างการใช้ queue ร่วมกับ transaction (Outbox Pattern)**

```go
type OutboxMessage struct {
    ID          uint
    AggregateID string
    EventType   string
    Payload     string
    Status      string // pending, processed, failed
    CreatedAt   time.Time
}

// ใน transaction หลัก
func (s *Service) CreateOrder(ctx context.Context, order *Order) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 1. บันทึก order
        if err := tx.Create(order).Error; err != nil {
            return err
        }

        // 2. ใส่ OutboxMessage
        msg := OutboxMessage{
            AggregateID: fmt.Sprintf("%d", order.ID),
            EventType:   "order.created",
            Payload:     `{"order_id":` + fmt.Sprintf("%d", order.ID) + `}`,
            Status:      "pending",
        }
        if err := tx.Create(&msg).Error; err != nil {
            return err
        }
        // transaction commit จะบันทึกทั้ง order และ outbox message
        return nil
    })
}

// Worker: อ่าน outbox messages และส่งไปยัง queue
func (s *Service) ProcessOutbox(ctx context.Context) {
    var messages []OutboxMessage
    s.db.Where("status = ?", "pending").Find(&messages)
    for _, msg := range messages {
        // ส่งไปยัง message broker
        if err := s.queue.Publish(msg.EventType, msg.Payload); err == nil {
            s.db.Model(&msg).Update("status", "processed")
        } else {
            // retry logic
        }
    }
}
```

ด้วยวิธีนี้ แม้ worker จะล้มเหลว เราสามารถ retry ได้โดยไม่สูญเสียข้อมูล

---

### 48.5 การออกแบบ Workflow

#### 48.5.1 Workflow CRUD พื้นฐาน

```
[Request] → [Controller] → [Service] → [Repository] → [GORM] → [Database]
                                                              ↓
[Response] ← [Controller] ← [Service] ← [Repository] ← [GORM] ← [Result]
```

#### 48.5.2 Workflow Query Cache

```
[Client] → [API] → ตรวจสอบ cache → มี → [Response] (cache hit)
                       ↓ ไม่มี
                 Query Database → เก็บ cache → [Response] (cache miss)
```

#### 48.5.3 Workflow Cache Queue Processor

```
[API] → (ต้องการ cache ข้อมูล) → สร้าง CacheJob → Push to Queue
[Worker] → Pop from Queue → Query DB → Store Cache (TTL)
[API] → (request ครั้งถัดไป) → Read Cache → Response
```

#### 48.5.4 Workflow Transaction + Queue (Outbox)

```
1. Start Transaction
   ├─ Update Business Data (Order)
   └─ Insert Outbox Message
2. Commit Transaction
3. Worker Polls Outbox
   ├─ Publish Message to Queue
   └─ Update Outbox Status to 'processed'
4. Consumer processes message
```

---

### 48.6 Code Template สำหรับนำไปใช้

#### 48.6.1 โครงสร้างโปรเจกต์

```
project/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── core/
│   │   └── user/
│   │       ├── entity.go
│   │       ├── repository.go
│   │       ├── service.go
│   │       └── handler.go
│   ├── platform/
│   │   ├── db/
│   │   │   └── gorm/
│   │   │       ├── session_factory.go
│   │   │       └── user_repo.go
│   │   └── cache/
│   │       ├── redis.go
│   │       └── cache_worker.go
│   └── transport/
│       └── http/
│           └── handler.go
├── pkg/
│   └── queue/
│       └── redis_queue.go
├── go.mod
└── config.yaml
```

#### 48.6.2 Session Factory Template

```go
// internal/platform/db/gorm/session_factory.go
package gorm

import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "time"
)

type SessionFactory struct {
    db *gorm.DB
}

func NewSessionFactory(dsn string) (*SessionFactory, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        // กำหนด logger, nowFunc ตามต้องการ
    })
    if err != nil {
        return nil, err
    }

    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    return &SessionFactory{db: db}, nil
}

func (sf *SessionFactory) GetDB() *gorm.DB {
    return sf.db
}

func (sf *SessionFactory) NewSession() *gorm.DB {
    return sf.db.Session(&gorm.Session{})
}
```

#### 48.6.3 Repository Implementation with Cache and Queue

```go
// internal/platform/db/gorm/user_repo.go
package gorm

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"
    "your-project/internal/core/user"
)

type userRepository struct {
    db    *gorm.DB
    cache *redis.Client
    queue *redis.Client // same as cache for simplicity
}

func NewUserRepository(db *gorm.DB, redisClient *redis.Client) user.Repository {
    return &userRepository{
        db:    db,
        cache: redisClient,
        queue: redisClient,
    }
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*user.User, error) {
    // Try cache first
    key := fmt.Sprintf("user:%d", id)
    val, err := r.cache.Get(ctx, key).Result()
    if err == nil {
        var u user.User
        if err := json.Unmarshal([]byte(val), &u); err == nil {
            return &u, nil
        }
    }

    // Cache miss, query DB
    var u user.User
    if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
        return nil, err
    }

    // Store in cache asynchronously via queue
    job := CacheJob{
        QueryKey: key,
        Table:    "users",
        ID:       id,
        TTL:      10 * time.Minute,
    }
    data, _ := json.Marshal(job)
    r.queue.LPush(ctx, "cache_queue", data)

    return &u, nil
}

// ... other methods
```

#### 48.6.4 Queue Worker Template

```go
// internal/platform/cache/cache_worker.go
package cache

import (
    "context"
    "encoding/json"
    "time"
    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"
)

type CacheWorker struct {
    redis *redis.Client
    db    *gorm.DB
    stop  chan struct{}
}

type CacheJob struct {
    QueryKey string        `json:"query_key"`
    Table    string        `json:"table"`
    ID       uint          `json:"id"`
    TTL      time.Duration `json:"ttl"`
}

func NewCacheWorker(redis *redis.Client, db *gorm.DB) *CacheWorker {
    return &CacheWorker{
        redis: redis,
        db:    db,
        stop:  make(chan struct{}),
    }
}

func (w *CacheWorker) Start(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case <-w.stop:
            return
        default:
            // Pop job from queue
            result, err := w.redis.BRPop(ctx, 0, "cache_queue").Result()
            if err != nil {
                continue
            }

            var job CacheJob
            if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
                continue
            }

            // Execute query based on table
            var data interface{}
            switch job.Table {
            case "users":
                var user User
                w.db.First(&user, job.ID)
                data = user
            // add more cases
            }

            // Store in cache
            bytes, _ := json.Marshal(data)
            w.redis.Set(ctx, job.QueryKey, bytes, job.TTL)
        }
    }
}

func (w *CacheWorker) Stop() {
    close(w.stop)
}
```

#### 48.6.5 การใช้ใน main.go

```go
// cmd/api/main.go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "your-project/internal/platform/db/gorm"
    "your-project/internal/platform/cache"
    "your-project/internal/core/user"
    userHttp "your-project/internal/transport/http"
)

func main() {
    // Load config
    cfg := loadConfig()

    // Setup session factory
    sf, err := gorm.NewSessionFactory(cfg.Database.DSN)
    if err != nil {
        log.Fatal(err)
    }

    // Setup Redis client
    redisClient := redis.NewClient(&redis.Options{
        Addr: cfg.Redis.Addr,
    })

    // Start cache worker
    worker := cache.NewCacheWorker(redisClient, sf.GetDB())
    ctx, cancel := context.WithCancel(context.Background())
    go worker.Start(ctx)

    // Setup repositories, services, handlers
    userRepo := gorm.NewUserRepository(sf.GetDB(), redisClient)
    userService := user.NewService(userRepo)
    userHandler := userHttp.NewHandler(userService)

    // Setup HTTP server
    router := setupRouter(userHandler)

    srv := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    cancel() // stop worker
    worker.Stop()

    ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancelShutdown()
    if err := srv.Shutdown(ctxShutdown); err != nil {
        log.Fatal("Server shutdown error:", err)
    }
}
```

---

### 48.7 สรุป

GORM เป็น ORM ที่ทรงพลังและใช้งานง่าย ช่วยให้นักพัฒนา Go จัดการฐานข้อมูลได้อย่างมีประสิทธิภาพ ด้วยฟีเจอร์ครบครันตั้งแต่ CRUD พื้นฐานไปจนถึง advanced patterns เช่น session factory, query cache, และการผสานกับระบบคิว เพื่อเพิ่มความเร็วและความน่าเชื่อถือของแอปพลิเคชัน

- **Session Factory** ช่วยแยก session และจัดการ connection pool
- **Query Cache** ลดภาระฐานข้อมูลด้วยการเก็บผลลัพธ์ซ้ำ
- **Cache Queue Processor** ทำให้การ pre-cache เป็น asynchronous ไม่กระทบ request ปัจจุบัน
- **SQL Queue Translate Rollback Commit** ผสาน transaction กับระบบคิวเพื่อความถูกต้องของข้อมูล

การใช้ GORM ร่วมกับหลักการออกแบบที่ดีจะช่วยให้คุณสร้างระบบที่มีประสิทธิภาพสูง บำรุงรักษาง่าย และพร้อมขยายตามความต้องการของธุรกิจ