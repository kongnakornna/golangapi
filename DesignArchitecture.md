# คู่มือสร้าง API Service ด้วยภาษา Go ตามหลักการ Best Practices ฉบับสมบูรณ์

## สารบัญ
- [คู่มือสร้าง API Service ด้วยภาษา Go ตามหลักการ Best Practices ฉบับสมบูรณ์](#คู่มือสร้าง-api-service-ด้วยภาษา-go-ตามหลักการ-best-practices-ฉบับสมบูรณ์)
  - [สารบัญ](#สารบัญ)
  - [บทนำ](#บทนำ)
  - [พื้นฐานที่ต้องรู้ก่อนพัฒนา API](#พื้นฐานที่ต้องรู้ก่อนพัฒนา-api)
    - [ความรู้ด้านเครือข่ายเบื้องต้น](#ความรู้ด้านเครือข่ายเบื้องต้น)
    - [สถาปัตยกรรมแบบ Monolithic](#สถาปัตยกรรมแบบ-monolithic)
    - [พื้นฐานภาษา Go](#พื้นฐานภาษา-go)
    - [แนวคิด OOP และ SOLID](#แนวคิด-oop-และ-solid)
    - [พื้นฐาน SQL และฐานข้อมูลเชิงสัมพันธ์](#พื้นฐาน-sql-และฐานข้อมูลเชิงสัมพันธ์)
    - [พื้นฐาน Git และการจัดการเวอร์ชัน](#พื้นฐาน-git-และการจัดการเวอร์ชัน)
  - [การออกแบบ API Service](#การออกแบบ-api-service)
    - [หลักการออกแบบ API ที่ดีใน Go](#หลักการออกแบบ-api-ที่ดีใน-go)
    - [Domain-Driven Design (DDD) เบื้องต้น](#domain-driven-design-ddd-เบื้องต้น)
    - [Clean Architecture คืออะไร](#clean-architecture-คืออะไร)
    - [โครงสร้างโปรเจกต์ตาม Clean Architecture](#โครงสร้างโปรเจกต์ตาม-clean-architecture)
  - [ลงมือสร้างโปรเจกต์ Isekai Shop API](#ลงมือสร้างโปรเจกต์-isekai-shop-api)
    - [การตั้งค่าโปรเจกต์เบื้องต้น](#การตั้งค่าโปรเจกต์เบื้องต้น)
    - [การจัดการ Configuration](#การจัดการ-configuration)
    - [การเชื่อมต่อฐานข้อมูล PostgreSQL ด้วย GORM](#การเชื่อมต่อฐานข้อมูล-postgresql-ด้วย-gorm)
    - [การทำ Database Migration](#การทำ-database-migration)
    - [การตั้งค่า Logging ด้วย slog](#การตั้งค่า-logging-ด้วย-slog)
    - [การตั้งค่า Redis สำหรับ Cache และ Queue](#การตั้งค่า-redis-สำหรับ-cache-และ-queue)
    - [การสร้าง HTTP Server ด้วย Echo](#การสร้าง-http-server-ด้วย-echo)
    - [การทำ Middleware](#การทำ-middleware)
    - [ระบบ Authentication และ Authorization](#ระบบ-authentication-และ-authorization)
      - [JWT Authentication](#jwt-authentication)
      - [Google OAuth2.0](#google-oauth20)
    - [โมดูลผู้ใช้ (User Module)](#โมดูลผู้ใช้-user-module)
    - [โมดูลไอเทม (Item Module)](#โมดูลไอเทม-item-module)
      - [การเพิ่มไอเทม](#การเพิ่มไอเทม)
      - [การแสดงรายการไอเทม](#การแสดงรายการไอเทม)
      - [การแก้ไขและลบไอเทม](#การแก้ไขและลบไอเทม)
      - [การค้นหาและกรองข้อมูล](#การค้นหาและกรองข้อมูล)
      - [การแบ่งหน้า (Pagination)](#การแบ่งหน้า-pagination)
    - [โมดูลผู้เล่นและเหรียญ (Player \& Coin Module)](#โมดูลผู้เล่นและเหรียญ-player--coin-module)
      - [การเพิ่มเหรียญ](#การเพิ่มเหรียญ)
      - [การแสดงเหรียญของผู้เล่น](#การแสดงเหรียญของผู้เล่น)
      - [การซื้อและขายไอเทม](#การซื้อและขายไอเทม)
      - [ระบบคลังสินค้า (Inventory)](#ระบบคลังสินค้า-inventory)
      - [บันทึกประวัติการซื้อขาย](#บันทึกประวัติการซื้อขาย)
    - [การจัดการ Transaction](#การจัดการ-transaction)
    - [การจัดการ Error อย่างเหมาะสม](#การจัดการ-error-อย่างเหมาะสม)
    - [การทำ Health Check และ Monitoring](#การทำ-health-check-และ-monitoring)
    - [การทำ Rate Limiting](#การทำ-rate-limiting)
    - [การทำ Caching](#การทำ-caching)
    - [การทำ Message Queue](#การทำ-message-queue)
  - [การ deploy แอปพลิเคชัน](#การ-deploy-แอปพลิเคชัน)
    - [เตรียมแอปสำหรับ production](#เตรียมแอปสำหรับ-production)
    - [Deploy บน GCP Cloud Run](#deploy-บน-gcp-cloud-run)
  - [สรุปและแหล่งข้อมูลเพิ่มเติม](#สรุปและแหล่งข้อมูลเพิ่มเติม)

---

## บทนำ

การพัฒนา API ที่มีประสิทธิภาพและสามารถรองรับการขยายตัวในอนาคตต้องอาศัยความรู้รอบด้านตั้งแต่พื้นฐานเครือข่าย ภาษาโปรแกรมมิ่ง การออกแบบซอฟต์แวร์ และการเลือกใช้เครื่องมือที่เหมาะสม บทความนี้รวบรวมเนื้อหาตั้งแต่พื้นฐานจนถึงการสร้างโปรเจกต์จริงด้วยภาษา Go โดยเน้นหลักการ Best Practices และสถาปัตยกรรมที่ทันสมัย เช่น Clean Architecture และ Domain-Driven Design (DDD) พร้อมตัวอย่างโค้ดจากโปรเจกต์ Isekai Shop API ซึ่งเป็นร้านค้าออนไลน์จำลองที่ครอบคลุมฟีเจอร์ต่าง ๆ ตั้งแต่การจัดการผู้ใช้ ไอเทม การซื้อขาย การยืนยันตัวตนด้วย JWT และ OAuth2 รวมถึงการใช้ Redis สำหรับ Cache และ Queue การ deploy บน Cloud Run

เนื้อหาทั้งหมดถูกจัดเรียงอย่างเป็นระบบเพื่อให้ผู้อ่านสามารถทำตามได้ทีละขั้นตอน เหมาะสำหรับนักพัฒนาที่มีพื้นฐานการเขียนโปรแกรมมาบ้างแล้วและต้องการยกระดับทักษะการออกแบบซอฟต์แวร์ด้วย Go

---

## พื้นฐานที่ต้องรู้ก่อนพัฒนา API

### ความรู้ด้านเครือข่ายเบื้องต้น
- **TCP vs UDP**: TCP (Transmission Control Protocol) เป็นโปรโตคอลที่เชื่อถือได้ รับประกันการส่งข้อมูลตามลำดับ ส่วน UDP (User Datagram Protocol) เร็วกว่าแต่ไม่รับประกันความถูกต้อง เหมาะสำหรับการสตรีมหรือเกม
- **HTTP**: Hypertext Transfer Protocol เป็นโปรโตคอลหลักสำหรับเว็บและ REST API ทำงานบน TCP ใช้วิธีการร้องขอเช่น GET, POST, PUT, DELETE
- **MQTT**: โปรโตคอลแบบ publish-subscribe สำหรับ IoT น้ำหนักเบา ทำงานบน TCP
- **SNMP**: โปรโตคอลสำหรับจัดการอุปกรณ์เครือข่าย

### สถาปัตยกรรมแบบ Monolithic
Monolithic คือการรวมทุกส่วนของแอปพลิเคชันไว้ในโค้ดเบสเดียว เหมาะสำหรับแอปขนาดเล็ก พัฒนาง่าย แต่เมื่อใหญ่ขึ้นจะยุ่งยากในการขยายและบำรุงรักษา บทความนี้ใช้แนวทางนี้ในการเริ่มต้นก่อนที่จะแยกเป็น microservices ได้ในอนาคต

### พื้นฐานภาษา Go
Go เป็นภาษาเชิงระบบที่เรียบง่าย มีประสิทธิภาพสูง และเหมาะสำหรับการพัฒนา API เนื่องจากการจัดการ concurrence ที่ดีด้วย goroutine และ channel

หัวข้อสำคัญใน Go:
- **Package**: การจัดกลุ่มโค้ดและการนำเข้า
- **Variables, Operators**: การประกาศตัวแปรและตัวดำเนินการ
- **Control Flow**: if, for, switch
- **Function**: ฟังก์ชันสามารถคืนค่าได้หลายค่า
- **Loop**: Go มีแค่ `for` แต่ใช้ได้ทั้งแบบ while และ infinite loop (พร้อมเทคนิคการ debug)
- **Pointers**: การส่งค่าด้วย reference
- **Array, Slice**: อาร์เรย์ขนาดคงที่ และ slice ที่ยืดหยุ่น
- **Map**: โครงสร้างข้อมูลแบบ key-value
- **Struct**: การรวมฟิลด์เข้าด้วยกัน คล้ายคลาส
- **Interface**: การกำหนดพฤติกรรม โดยไม่ต้องสืบทอด
- **Generics**: เริ่มมีใน Go 1.18 ทำให้เขียนโค้ดที่ reuse ได้มากขึ้น
- **Goroutines**: หน่วยการทำงานพร้อมกันที่เบา
- **Channel**: การสื่อสารระหว่าง goroutine
- **Mutex**: การป้องกัน race condition

### แนวคิด OOP และ SOLID
แม้ Go ไม่สนับสนุน OOP เต็มรูปแบบ แต่เราสามารถจำลองแนวคิดผ่าน struct และ interface ได้
- **OOP คืออะไร**: การเขียนโปรแกรมเชิงวัตถุ เน้นการห่อหุ้มข้อมูลและพฤติกรรม
- **Pillars of OOP**: Encapsulation, Inheritance (ใน Go ใช้ composition), Polymorphism (ผ่าน interface), Abstraction
- **ความสัมพันธ์ระหว่าง Objects**: Association, Aggregation, Composition
- **SOLID Principles**:
  - Single Responsibility: หนึ่งคลาสควรมีหน้าที่เดียว
  - Open/Closed: เปิดให้ขยาย ปิดให้แก้ไข
  - Liskov Substitution:  subtype ต้องแทนที่ supertype ได้
  - Interface Segregation: แยก interface เฉพาะที่จำเป็น
  - Dependency Inversion: ขึ้นกับ abstraction ไม่ใช่ concretion

### พื้นฐาน SQL และฐานข้อมูลเชิงสัมพันธ์
- **SQL คืออะไร**: ภาษาในการจัดการฐานข้อมูลเชิงสัมพันธ์
- **Relationship**: One-to-One, One-to-Many, Many-to-Many
- **การติดตั้ง PostgreSQL ด้วย Docker**:
  ```bash
  docker run --name postgres -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres
  ```
- คำสั่งพื้นฐาน: INSERT, SELECT, WHERE, LIKE, AND/OR, ORDER BY, UPDATE, DELETE
- **Join**: INNER JOIN, LEFT JOIN เป็นต้น
- **Transaction**: การทำรายการหลายขั้นตอนให้สำเร็จหรือล้มเหลวพร้อมกัน

### พื้นฐาน Git และการจัดการเวอร์ชัน
- **Git คืออะไร**: ระบบควบคุมเวอร์ชันแบบกระจาย
- **Git Quick Start**: init, add, commit, branch, merge, remote
- **Git Flow**: branching model ที่มี main, develop, feature, release, hotfix

---

## การออกแบบ API Service

### หลักการออกแบบ API ที่ดีใน Go
- ใช้ RESTful concept (HTTP methods + resource)
- มี versioning (เช่น `/api/v1/`)
- ใช้ JSON ในการแลกเปลี่ยนข้อมูล
- มี error handling ที่ชัดเจน
- ใช้ status codes อย่างเหมาะสม
- มีการ validation input
- มีการ authentication และ authorization
- มี logging และ monitoring

### Domain-Driven Design (DDD) เบื้องต้น
DDD เป็นแนวทางการออกแบบซอฟต์แวร์โดยเน้นที่ domain หลักและภาษาที่ใช้ร่วมกันระหว่างทีมพัฒนาและผู้เชี่ยวชาญด้านธุรกิจ
- **Ubiquitous Language**: ภาษากลางที่ทุกฝ่ายใช้ตรงกัน
- **Bounded Context**: การแบ่ง domain ออกเป็นบริบทย่อย ๆ แต่ละบริบทมี model ของตนเอง
- **Context Mapping**: การเชื่อมโยงระหว่าง bounded context
- **Entity vs Value Object**: Entity มี identity (เช่น ผู้ใช้) ส่วน Value Object ไม่มี identity (เช่น ที่อยู่)

### Clean Architecture คืออะไร
Clean Architecture เป็นสถาปัตยกรรมที่แยกส่วนออกเป็นชั้น ๆ โดยมีกฎ Dependency Rule: ชั้นในไม่ควรพึ่งพาชั้นนอก ชั้นนอกพึ่งพาชั้นผ่าน interface
ชั้นต่าง ๆ:
- **Entities**: ข้อมูลและกฎธุรกิจทั่วไป
- **Use Cases**: กฎเฉพาะของแอปพลิเคชัน
- **Interface Adapters**: แปลงข้อมูลระหว่าง use cases และภายนอก (controllers, presenters)
- **Frameworks & Drivers**: ฐานข้อมูล, web framework, etc.

### โครงสร้างโปรเจกต์ตาม Clean Architecture
ตัวอย่างโครงสร้างสำหรับโปรเจกต์ Go:
```
shop/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── core/
│   │   ├── domain/          # Entities, value objects
│   │   ├── usecases/        # Business logic (interfaces)
│   │   └── ports/           # Interfaces for outer layers
│   ├── adapters/
│   │   ├── repositories/    # Database implementations
│   │   ├── handlers/        # HTTP handlers
│   │   ├── middlewares/     # HTTP middlewares
│   │   └── routers/         # Routing
│   └── pkg/                  # Shared packages (logger, config, etc.)
├── migrations/               # SQL migration files
├── configs/                  # Config files
├── go.mod
└── Makefile
```

---

## ลงมือสร้างโปรเจกต์ Isekai Shop API

เราจะสร้าง API สำหรับร้านค้าออนไลน์ในโลกอีกใบ (Isekai) ที่มีฟีเจอร์หลัก:
- จัดการผู้ใช้ (User)
- จัดการไอเทม (Item)
- จัดการผู้เล่นและเหรียญ (Player & Coin)
- ซื้อขายไอเทม
- ยืนยันตัวตนด้วย JWT และ Google OAuth2
- ใช้ Redis สำหรับ cache และ queue
- มี logging, monitoring, rate limiting

### การตั้งค่าโปรเจกต์เบื้องต้น
สร้างโฟลเดอร์และ initialize Go module:
```bash
mkdir isekai-shop
cd isekai-shop
go mod init github.com/yourname/isekai-shop
```

ติดตั้ง package ที่จำเป็น:
```bash
go get -u github.com/labstack/echo/v4
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/go-redis/redis/v8
go get -u github.com/golang-jwt/jwt/v5
go get -u golang.org/x/oauth2
go get -u github.com/google/uuid
go get -u github.com/spf13/viper
go get -u go.uber.org/zap # หรือใช้ slog ในตัว
```

### การจัดการ Configuration
ใช้ Viper อ่านค่าจากไฟล์ `.env` หรือ `config.yaml` สร้าง `configs/config.go`:
```go
package config

import "github.com/spf13/viper"

type Config struct {
    ServerPort string `mapstructure:"SERVER_PORT"`
    DBHost     string `mapstructure:"DB_HOST"`
    DBPort     string `mapstructure:"DB_PORT"`
    DBUser     string `mapstructure:"DB_USER"`
    DBPassword string `mapstructure:"DB_PASSWORD"`
    DBName     string `mapstructure:"DB_NAME"`
    RedisAddr  string `mapstructure:"REDIS_ADDR"`
    JWTSecret  string `mapstructure:"JWT_SECRET"`
    // ...
}

func LoadConfig(path string) (config Config, err error) {
    viper.AddConfigPath(path)
    viper.SetConfigName("app")
    viper.SetConfigType("env")
    viper.AutomaticEnv()
    err = viper.ReadInConfig()
    if err != nil {
        return
    }
    err = viper.Unmarshal(&config)
    return
}
```

### การเชื่อมต่อฐานข้อมูล PostgreSQL ด้วย GORM
สร้าง `internal/adapters/repositories/db.go`:
```go
package repositories

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "isekai-shop/internal/core/domain"
)

func NewGormDB(cfg config.Config) (*gorm.DB, error) {
    dsn := "host=" + cfg.DBHost + " user=" + cfg.DBUser + " password=" + cfg.DBPassword +
           " dbname=" + cfg.DBName + " port=" + cfg.DBPort + " sslmode=disable TimeZone=Asia/Bangkok"
    return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
```

### การทำ Database Migration
ใช้ GORM AutoMigrate หรือเครื่องมืออย่าง `golang-migrate` ในที่นี้ใช้ AutoMigrate ใน `main.go`:
```go
db.AutoMigrate(&domain.User{}, &domain.Item{}, &domain.Player{}, &domain.Inventory{}, &domain.PurchaseHistory{})
```

### การตั้งค่า Logging ด้วย slog
Go 1.21 มี `log/slog` ในตัว สร้างแพ็กเกจ `pkg/logger`:
```go
package logger

import (
    "log/slog"
    "os"
)

var Log *slog.Logger

func Init() {
    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
    Log = slog.New(handler)
}
```

### การตั้งค่า Redis สำหรับ Cache และ Queue
สร้าง `pkg/cache/redis.go`:
```go
package cache

import (
    "context"
    "github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis(addr string) {
    RedisClient = redis.NewClient(&redis.Options{
        Addr: addr,
    })
}
```
และสำหรับ queue อาจใช้ Redis pub/sub หรือสร้าง interface เผื่อเปลี่ยนได้

### การสร้าง HTTP Server ด้วย Echo
ใน `cmd/api/main.go`:
```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"

    "github.com/labstack/echo/v4"
    "isekai-shop/config"
    "isekai-shop/internal/adapters/handlers"
    "isekai-shop/internal/adapters/middlewares"
    "isekai-shop/internal/adapters/repositories"
    "isekai-shop/internal/core/usecases"
    "isekai-shop/pkg/logger"
)

func main() {
    // Load config
    cfg, err := config.LoadConfig("./configs")
    if err != nil {
        log.Fatal("cannot load config:", err)
    }

    // Init logger
    logger.Init()

    // Init database
    db, err := repositories.NewGormDB(cfg)
    if err != nil {
        logger.Log.Error("cannot connect to db", "error", err)
        os.Exit(1)
    }

    // Init Redis
    cache.InitRedis(cfg.RedisAddr)

    // Setup Echo
    e := echo.New()

    // Middlewares
    e.Use(middlewares.RequestLogger())
    e.Use(middlewares.Recover())
    e.Use(middlewares.SecurityHeaders())
    e.Use(middlewares.RateLimiter())

    // Dependency Injection
    userRepo := repositories.NewUserRepository(db)
    userUsecase := usecases.NewUserUsecase(userRepo)
    userHandler := handlers.NewUserHandler(userUsecase)

    authUsecase := usecases.NewAuthUsecase(userRepo, cfg.JWTSecret)
    authHandler := handlers.NewAuthHandler(authUsecase)

    // Routes
    api := e.Group("/api/v1")
    handlers.RegisterUserRoutes(api, userHandler)
    handlers.RegisterAuthRoutes(api, authHandler)

    // Start server
    go func() {
        if err := e.Start(":" + cfg.ServerPort); err != nil && err != http.ErrServerClosed {
            logger.Log.Error("shutting down the server", "error", err)
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)
    <-quit
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := e.Shutdown(ctx); err != nil {
        e.Logger.Fatal(err)
    }
}
```

### การทำ Middleware
ตัวอย่าง middleware บางส่วน:
- **RequestLogger**: บันทึก request/response
- **Recover**: จัดการ panic
- **SecurityHeaders**: เพิ่ม headers ด้านความปลอดภัย (CSP, HSTS, ฯลฯ)
- **RateLimiter**: จำกัด request ต่อ IP (ใช้ Redis หรือ in-memory)
- **AuthMiddleware**: ตรวจสอบ JWT token

สร้างไฟล์ `internal/adapters/middlewares/auth.go`:
```go
package middlewares

import (
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v5"
    "github.com/labstack/echo/v4"
)

func AuthMiddleware(secret string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            authHeader := c.Request().Header.Get("Authorization")
            if authHeader == "" {
                return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
            }
            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                return echo.NewHTTPError(http.StatusUnauthorized, "invalid token format")
            }
            token := parts[1]
            // Parse and validate token
            claims := jwt.MapClaims{}
            _, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
                return []byte(secret), nil
            })
            if err != nil {
                return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
            }
            c.Set("user_id", claims["user_id"])
            return next(c)
        }
    }
}
```

### ระบบ Authentication และ Authorization

#### JWT Authentication
- สร้าง usecase สำหรับ auth มีฟังก์ชัน Login, Register, RefreshToken
- ใช้ JWT สร้าง access token (อายุสั้น) และ refresh token (อายุยาว) เก็บ refresh token ในฐานข้อมูลหรือ Redis
- มี blacklist สำหรับ token ที่ถูก revoke

#### Google OAuth2.0
- สร้าง OAuth2.0 app ใน Google Cloud Platform
- ใช้ `golang.org/x/oauth2` และ `oauth2/google`
- Flow:
  1. ผู้ใช้คลิก login กับ Google
  2. redirect ไปยัง Google consent screen
  3. Google redirect กลับมาพร้อม code
  4. แลก code เป็น token
  5. ดึงข้อมูลผู้ใช้จาก Google แล้วสร้าง/อัปเดตผู้ใช้ในระบบ
  6. สร้าง JWT ให้กับผู้ใช้

### โมดูลผู้ใช้ (User Module)
- CRUD ผู้ใช้
- มี role (admin, user) เบื้องต้น
- admin เท่านั้นที่สร้างผู้ใช้ใหม่ (ตาม requirement)

### โมดูลไอเทม (Item Module)
#### การเพิ่มไอเทม
- เฉพาะ admin เท่านั้นที่เพิ่มได้
- รับข้อมูลผ่าน DTO (Data Transfer Object) แล้ว validate
- บันทึกลง DB

#### การแสดงรายการไอเทม
- GET `/items` คืนค่ารายการไอเทมทั้งหมดหรือตามเงื่อนไข

#### การแก้ไขและลบไอเทม
- PUT, PATCH, DELETE สำหรับ admin

#### การค้นหาและกรองข้อมูล
- ใช้ query parameters เช่น `?name=sword&minPrice=100`
- ใน repository layer ใช้ GORM สร้าง dynamic query

#### การแบ่งหน้า (Pagination)
- รับ `page` และ `pageSize` จาก query
- ใช้ `.Offset().Limit()` ของ GORM

### โมดูลผู้เล่นและเหรียญ (Player & Coin Module)
ผู้เล่นคือผู้ใช้ที่มีข้อมูลเพิ่มเติม เช่น เหรียญ คลังสินค้า

#### การเพิ่มเหรียญ
- admin สามารถเพิ่มเหรียญให้ผู้เล่นได้ (เช่น จากการซื้อจริง)

#### การแสดงเหรียญของผู้เล่น
- GET `/player/coins`

#### การซื้อและขายไอเทม
- ซื้อ: ตรวจสอบเหรียญเพียงพอ และไอเทมมี stock (ถ้ามี) หักเหรียญ เพิ่มใน inventory
- ขาย: ตรวจสอบว่ามีไอเทมใน inventory เพิ่มเหรียญ คืนไอเทม

#### ระบบคลังสินค้า (Inventory)
- ตาราง inventory เชื่อม user กับ item พร้อมจำนวน
- เมื่อซื้อ เพิ่มจำนวน เมื่อขาย ลดจำนวน

#### บันทึกประวัติการซื้อขาย
- ตาราง purchase_history เก็บข้อมูลการซื้อขาย (user, item, quantity, price, type)

### การจัดการ Transaction
เพื่อความถูกต้องของข้อมูลในการซื้อขาย (หลายตาราง) ต้องใช้ transaction:
```go
func (u *playerUsecase) BuyItem(ctx context.Context, userID, itemID uint, quantity int) error {
    return u.repo.WithTransaction(ctx, func(tx *gorm.DB) error {
        // 1. ดึงข้อมูลผู้เล่น (lock row)
        // 2. ตรวจสอบเหรียญ
        // 3. ดึงไอเทมและราคา
        // 4. คำนวณราคารวม
        // 5. อัปเดตเหรียญผู้เล่น
        // 6. เพิ่มใน inventory
        // 7. บันทึกประวัติ
        return nil
    })
}
```

### การจัดการ Error อย่างเหมาะสม
- สร้าง custom error type ที่มี status code และข้อความ
- ใน handler จับ error แล้วส่ง response ที่เหมาะสม
- ใช้ `echo.NewHTTPError` หรือ implement `HTTPError` interface

### การทำ Health Check และ Monitoring
- สร้าง endpoint `/health` สำหรับ basic check
- `/health/detailed` เช็คการเชื่อมต่อ DB, Redis
- `/ready`, `/live` สำหรับ Kubernetes

### การทำ Rate Limiting
- ใช้ middleware จำกัด request ต่อ IP
- อาจใช้ Redis เก็บ counter

### การทำ Caching
- ใช้ Redis cache สำหรับข้อมูลที่อ่านบ่อยไม่ค่อยเปลี่ยน เช่น รายการไอเทม
- สร้าง interface `Cache` และ implement ด้วย Redis

### การทำ Message Queue
- ใช้ Redis pub/sub สำหรับส่ง event ระหว่างส่วนต่าง ๆ
- สร้าง interface `Queue` มี method Publish, Subscribe
- มี Dead Letter Queue สำหรับข้อความที่ประมวลผลไม่สำเร็จ

---

## การ deploy แอปพลิเคชัน

### เตรียมแอปสำหรับ production
- ใช้ environment variables สำหรับ config
- build binary: `GOOS=linux go build -o app cmd/api/main.go`
- เขียน Dockerfile
```
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/api/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
COPY configs/ ./configs/
EXPOSE 8080
CMD ["./main"]
```

### Deploy บน GCP Cloud Run
1. สร้างโปรเจกต์บน GCP
2. เปิดใช้งาน Cloud Run และ Artifact Registry
3. สร้าง Docker image และ push ไปยัง Artifact Registry
4. Deploy บน Cloud Run โดยกำหนด environment variables และตั้งค่าให้เชื่อมต่อกับ Cloud SQL (PostgreSQL) และ Redis (อาจใช้ Memorystore)

---

## สรุปและแหล่งข้อมูลเพิ่มเติม

บทความนี้ได้นำเสนอตั้งแต่พื้นฐานที่จำเป็น การออกแบบระบบด้วย DDD และ Clean Architecture ไปจนถึงการลงมือเขียน API ด้วย Go พร้อมฟีเจอร์ต่าง ๆ ครบถ้วน การทำตามขั้นตอนจะช่วยให้คุณสามารถสร้าง API ที่มีคุณภาพ รองรับการขยาย และบำรุงรักษาง่าย

แหล่งข้อมูลเพิ่มเติม:
- [Official Go Documentation](https://go.dev/doc/)
- [Echo Framework](https://echo.labstack.com/)
- [GORM](https://gorm.io/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://domainlanguage.com/ddd/)

 