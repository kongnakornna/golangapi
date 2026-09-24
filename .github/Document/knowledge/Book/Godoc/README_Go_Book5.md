# Golangiot - Production-Ready Go REST API คู่มือฉบับสมบูรณ์

## สารบัญ
1. [ภาพรวมสถาปัตยกรรม](#1-ภาพรวมสถาปัตยกรรม)
2. [โครงสร้างโปรเจค](#2-โครงสร้างโปรเจค)
3. [การติดตั้งและ Setup](#3-การติดตั้งและ-setup)
4. [Core Layer Implementation](#4-core-layer-implementation)
5. [Repository Layer](#5-repository-layer)
6. [Service Layer](#6-service-layer)
7. [Handler Layer](#7-handler-layer)
8. [Middleware](#8-middleware)
9. [Cache Layer](#9-cache-layer)
10. [Message Queue](#10-message-queue)
11. [Authentication & Authorization](#11-authentication--authorization)
12. [Health Monitoring](#12-health-monitoring)
13. [Configuration Management](#13-configuration-management)
14. [Testing](#14-testing)
15. [Deployment](#15-deployment)

---

## 1. ภาพรวมสถาปัตยกรรม

Golangiot ใช้สถาปัตยกรรมแบบ Clean Architecture พร้อม三层分层 (Three-Layer Architecture):

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Server (Chi)                       │
├─────────────────────────────────────────────────────────────┤
│                     Middleware Layer                         │
│  (Auth, Logging, Rate Limit, CORS, Security Headers)       │
├─────────────────────────────────────────────────────────────┤
│                      Handler Layer                           │
│           (Request Validation, Response Format)             │
├─────────────────────────────────────────────────────────────┤
│                      Service Layer                           │
│              (Business Logic, Orchestration)                │
├─────────────────────────────────────────────────────────────┤
│                    Repository Layer                          │
│          (Database Operations, Cache, Queue)                │
├─────────────────────────────────────────────────────────────┤
│                    Infrastructure                            │
│        (PostgreSQL, Redis, GORM, slog)                      │
└─────────────────────────────────────────────────────────────┘
```

---

## 2. โครงสร้างโปรเจค

```
golangiot/
├── cmd/
│   └── api/
│       └── main.go                    # Application entry point
├── internal/
│   ├── domain/
│   │   ├── entity/
│   │   │   ├── user.go               # User entity
│   │   │   └── token.go              # Token entity
│   │   ├── repository/
│   │   │   ├── user_repository.go    # User repository interface
│   │   │   └── token_repository.go   # Token repository interface
│   │   └── service/
│   │       ├── auth_service.go       # Auth service interface
│   │       └── user_service.go       # User service interface
│   ├── repository/
│   │   ├── postgres/
│   │   │   ├── user_repository.go    # User repository implementation
│   │   │   ├── token_repository.go   # Token repository implementation
│   │   │   └── db.go                 # Database connection
│   │   └── redis/
│   │       ├── cache.go              # Cache implementation
│   │       └── queue.go              # Message queue implementation
│   ├── service/
│   │   ├── auth_service.go           # Auth service implementation
│   │   └── user_service.go           # User service implementation
│   ├── handler/
│   │   ├── auth_handler.go           # Auth HTTP handlers
│   │   ├── user_handler.go           # User HTTP handlers
│   │   ├── health_handler.go         # Health check handlers
│   │   └── response.go               # Response formatters
│   └── middleware/
│       ├── auth.go                   # JWT authentication
│       ├── logger.go                 # Request logging
│       ├── rate_limit.go             # Rate limiting
│       ├── security.go               # Security headers
│       └── context.go                # Request context
├── pkg/
│   ├── logger/
│   │   └── logger.go                 # Structured logging
│   ├── cache/
│   │   └── cache.go                  # Cache interface
│   ├── queue/
│   │   └── queue.go                  # Queue interface
│   ├── errors/
│   │   └── errors.go                 # Custom errors
│   └── validator/
│       └── validator.go              # Input validation
├── configs/
│   ├── config.yaml                   # Default configuration
│   ├── config.go                     # Configuration loader
│   └── config.dev.yaml               # Development config
├── migrations/
│   ├── 001_create_users_table.sql
│   └── 002_create_tokens_table.sql
├── scripts/
│   ├── seed.go                       # Database seeder
│   └── migrate.sh                    # Migration script
├── test/
│   ├── integration/
│   └── e2e/
├── deployments/
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   └── k8s/
│       └── deployment.yaml
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## 3. การติดตั้งและ Setup

### 3.1 Prerequisites

```bash
# Required
Go 1.21+
PostgreSQL 14+
Redis 7+

# Optional
Docker & Docker Compose
Make
golangci-lint
```

### 3.2 Initial Setup

```bash
# Clone repository
git clone https://github.com/yourusername/golangiot.git
cd golangiot

# Copy environment file
cp .env.example .env

# Install dependencies
go mod download

# Run database migrations
make migrate-up

# Seed database
go run scripts/seed.go

# Run application
make run
```

### 3.3 Environment Variables (.env.example)

```env
# Application
APP_NAME=golangiot
APP_ENV=development
APP_DEBUG=true
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=golangiot
DB_SSL_MODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-secret-key-change-in-production
JWT_ACCESS_EXPIRE=15m
JWT_REFRESH_EXPIRE=168h

# Rate Limit
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

---

## 4. Core Layer Implementation

### 4.1 Main Entry Point (cmd/api/main.go)

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/cors"
    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"

    "golangiot/configs"
    "golangiot/internal/handler"
    "golangiot/internal/middleware"
    "golangiot/internal/repository/postgres"
    "golangiot/internal/repository/redis"
    "golangiot/internal/service"
    "golangiot/pkg/logger"
)

func main() {
    // Load configuration
    cfg, err := configs.LoadConfig("configs")
    if err != nil {
        panic(fmt.Sprintf("Failed to load config: %v", err))
    }

    // Initialize logger
    log, err := logger.NewLogger(&logger.Config{
        Level:  cfg.Log.Level,
        Format: cfg.Log.Format,
        Output: cfg.Log.Output,
    })
    if err != nil {
        panic(fmt.Sprintf("Failed to initialize logger: %v", err))
    }
    defer log.Sync()

    log.Info("Starting application", "name", cfg.App.Name, "env", cfg.App.Env)

    // Initialize database
    db, err := postgres.NewDatabase(&cfg.Database)
    if err != nil {
        log.Fatal("Failed to connect to database", "error", err)
    }
    log.Info("Database connected")

    // Initialize Redis client
    redisClient, err := redis.NewRedisClient(&cfg.Redis)
    if err != nil {
        log.Fatal("Failed to connect to Redis", "error", err)
    }
    log.Info("Redis connected")

    // Initialize repositories
    userRepo := postgres.NewUserRepository(db)
    tokenRepo := postgres.NewTokenRepository(db)
    
    // Initialize cache
    cache := redis.NewCache(redisClient)
    
    // Initialize queue
    queue := redis.NewQueue(redisClient)

    // Initialize services
    authService := service.NewAuthService(userRepo, tokenRepo, cache, cfg.JWT)
    userService := service.NewUserService(userRepo, cache, queue)

    // Initialize handlers
    authHandler := handler.NewAuthHandler(authService, log)
    userHandler := handler.NewUserHandler(userService, log)
    healthHandler := handler.NewHealthHandler(db, redisClient, log)

    // Setup router
    router := chi.NewRouter()

    // Setup middleware
    setupMiddleware(router, cfg, log)

    // Setup routes
    setupRoutes(router, authHandler, userHandler, healthHandler, cfg, log)

    // Create HTTP server
    server := &http.Server{
        Addr:         fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
        Handler:      router,
        ReadTimeout:  cfg.HTTP.ReadTimeout,
        WriteTimeout: cfg.HTTP.WriteTimeout,
        IdleTimeout:  cfg.HTTP.IdleTimeout,
    }

    // Graceful shutdown
    go func() {
        log.Info("Server started", "addr", server.Addr)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal("Server failed", "error", err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Info("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.GracefulTimeout)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Error("Server forced to shutdown", "error", err)
    }

    log.Info("Server exited")
}

func setupMiddleware(router *chi.Mux, cfg *configs.Config, log logger.Logger) {
    // Recovery middleware
    router.Use(middleware.Recoverer)
    
    // Request ID middleware
    router.Use(middleware.RequestID)
    
    // Real IP middleware
    router.Use(middleware.RealIP)
    
    // Logger middleware
    router.Use(middleware.NewLogger(log))
    
    // CORS middleware
    router.Use(cors.Handler(cors.Options{
        AllowedOrigins:   cfg.CORS.AllowedOrigins,
        AllowedMethods:   cfg.CORS.AllowedMethods,
        AllowedHeaders:   cfg.CORS.AllowedHeaders,
        ExposedHeaders:   cfg.CORS.ExposedHeaders,
        AllowCredentials: cfg.CORS.AllowCredentials,
        MaxAge:           cfg.CORS.MaxAge,
    }))
    
    // Security headers middleware
    router.Use(middleware.SecurityHeaders)
    
    // Rate limit middleware
    if cfg.RateLimit.Enabled {
        limiter := middleware.NewRateLimiter(cfg.RateLimit.Requests, cfg.RateLimit.Window)
        router.Use(limiter.Middleware)
    }
    
    // Request context middleware
    router.Use(middleware.RequestContext)
}

func setupRoutes(
    router *chi.Mux,
    authHandler *handler.AuthHandler,
    userHandler *handler.UserHandler,
    healthHandler *handler.HealthHandler,
    cfg *configs.Config,
    log logger.Logger,
) {
    // Health check routes (no auth required)
    router.Get("/health", healthHandler.Health)
    router.Get("/health/detailed", healthHandler.DetailedHealth)
    router.Get("/ready", healthHandler.Ready)
    router.Get("/live", healthHandler.Live)

    // Public routes
    router.Group(func(r chi.Router) {
        r.Post("/api/v1/auth/register", authHandler.Register)
        r.Post("/api/v1/auth/login", authHandler.Login)
        r.Post("/api/v1/auth/refresh", authHandler.RefreshToken)
    })

    // Protected routes
    router.Group(func(r chi.Router) {
        // Apply authentication middleware
        authMiddleware := middleware.NewAuthMiddleware(cfg.JWT, log)
        r.Use(authMiddleware.Authenticate)
        
        // Auth routes
        r.Post("/api/v1/auth/logout", authHandler.Logout)
        r.Get("/api/v1/auth/me", authHandler.Me)
        
        // User routes
        r.Get("/api/v1/users", userHandler.List)
        r.Get("/api/v1/users/{id}", userHandler.GetByID)
        r.Put("/api/v1/users/{id}", userHandler.Update)
        r.Delete("/api/v1/users/{id}", userHandler.Delete)
        
        // Admin only routes
        r.Group(func(r chi.Router) {
            r.Use(authMiddleware.RequireRole("admin"))
            r.Post("/api/v1/users", userHandler.Create)
        })
    })
}
```

### 4.2 Configuration Management (configs/config.go)

```go
package configs

import (
    "fmt"
    "time"
    
    "github.com/spf13/viper"
    "github.com/caarlos0/env/v8"
)

type Config struct {
    App       AppConfig       `mapstructure:"app"`
    HTTP      HTTPConfig      `mapstructure:"http"`
    Database  DatabaseConfig  `mapstructure:"database"`
    Redis     RedisConfig     `mapstructure:"redis"`
    JWT       JWTConfig       `mapstructure:"jwt"`
    CORS      CORSConfig      `mapstructure:"cors"`
    RateLimit RateLimitConfig `mapstructure:"rate_limit"`
    Log       LogConfig       `mapstructure:"log"`
}

type AppConfig struct {
    Name    string `mapstructure:"name" env:"APP_NAME"`
    Version string `mapstructure:"version" env:"APP_VERSION"`
    Env     string `mapstructure:"env" env:"APP_ENV" envDefault:"development"`
    Debug   bool   `mapstructure:"debug" env:"APP_DEBUG" envDefault:"false"`
}

type HTTPConfig struct {
    Host            string        `mapstructure:"host" env:"APP_HOST" envDefault:"0.0.0.0"`
    Port            int           `mapstructure:"port" env:"APP_PORT" envDefault:"8080"`
    ReadTimeout     time.Duration `mapstructure:"read_timeout" envDefault:"30s"`
    WriteTimeout    time.Duration `mapstructure:"write_timeout" envDefault:"30s"`
    IdleTimeout     time.Duration `mapstructure:"idle_timeout" envDefault:"60s"`
    GracefulTimeout time.Duration `mapstructure:"graceful_timeout" envDefault:"15s"`
}

type DatabaseConfig struct {
    Host            string        `mapstructure:"host" env:"DB_HOST" envDefault:"localhost"`
    Port            int           `mapstructure:"port" env:"DB_PORT" envDefault:"5432"`
    User            string        `mapstructure:"user" env:"DB_USER" envDefault:"postgres"`
    Password        string        `mapstructure:"password" env:"DB_PASSWORD"`
    DBName          string        `mapstructure:"dbname" env:"DB_NAME" envDefault:"golangiot"`
    SSLMode         string        `mapstructure:"ssl_mode" env:"DB_SSL_MODE" envDefault:"disable"`
    MaxOpenConns    int           `mapstructure:"max_open_conns" envDefault:"100"`
    MaxIdleConns    int           `mapstructure:"max_idle_conns" envDefault:"10"`
    ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" envDefault:"1h"`
}

type RedisConfig struct {
    Host     string `mapstructure:"host" env:"REDIS_HOST" envDefault:"localhost"`
    Port     int    `mapstructure:"port" env:"REDIS_PORT" envDefault:"6379"`
    Password string `mapstructure:"password" env:"REDIS_PASSWORD"`
    DB       int    `mapstructure:"db" env:"REDIS_DB" envDefault:"0"`
    PoolSize int    `mapstructure:"pool_size" envDefault:"10"`
}

type JWTConfig struct {
    Secret        string        `mapstructure:"secret" env:"JWT_SECRET"`
    AccessExpire  time.Duration `mapstructure:"access_expire" env:"JWT_ACCESS_EXPIRE" envDefault:"15m"`
    RefreshExpire time.Duration `mapstructure:"refresh_expire" env:"JWT_REFRESH_EXPIRE" envDefault:"168h"`
    Issuer        string        `mapstructure:"issuer" envDefault:"golangiot"`
}

type CORSConfig struct {
    AllowedOrigins   []string `mapstructure:"allowed_origins"`
    AllowedMethods   []string `mapstructure:"allowed_methods"`
    AllowedHeaders   []string `mapstructure:"allowed_headers"`
    ExposedHeaders   []string `mapstructure:"exposed_headers"`
    AllowCredentials bool     `mapstructure:"allow_credentials"`
    MaxAge           int      `mapstructure:"max_age"`
}

type RateLimitConfig struct {
    Enabled  bool          `mapstructure:"enabled" env:"RATE_LIMIT_ENABLED" envDefault:"true"`
    Requests int           `mapstructure:"requests" env:"RATE_LIMIT_REQUESTS" envDefault:"100"`
    Window   time.Duration `mapstructure:"window" env:"RATE_LIMIT_WINDOW" envDefault:"1m"`
}

type LogConfig struct {
    Level  string `mapstructure:"level" env:"LOG_LEVEL" envDefault:"info"`
    Format string `mapstructure:"format" env:"LOG_FORMAT" envDefault:"json"`
    Output string `mapstructure:"output" env:"LOG_OUTPUT" envDefault:"stdout"`
}

func LoadConfig(path string) (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(path)
    viper.AddConfigPath(".")
    
    // Environment variables
    viper.AutomaticEnv()
    viper.SetEnvPrefix("APP")
    
    // Default values
    setDefaults()
    
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, fmt.Errorf("error reading config file: %w", err)
        }
    }
    
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, fmt.Errorf("error unmarshaling config: %w", err)
    }
    
    // Override with environment variables
    if err := env.Parse(&config); err != nil {
        return nil, fmt.Errorf("error parsing env: %w", err)
    }
    
    return &config, nil
}

func setDefaults() {
    viper.SetDefault("app.env", "development")
    viper.SetDefault("http.port", 8080)
    viper.SetDefault("http.read_timeout", "30s")
    viper.SetDefault("database.max_open_conns", 100)
    viper.SetDefault("rate_limit.requests", 100)
    viper.SetDefault("rate_limit.window", "1m")
    viper.SetDefault("log.level", "info")
    viper.SetDefault("log.format", "json")
}
```

### 4.3 Configuration File (configs/config.yaml)

```yaml
app:
  name: "golangiot"
  version: "1.0.0"
  env: "development"
  debug: true

http:
  host: "0.0.0.0"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "60s"
  graceful_timeout: "15s"

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "golangiot"
  ssl_mode: "disable"
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: "1h"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10

jwt:
  secret: "your-secret-key-change-in-production"
  access_expire: "15m"
  refresh_expire: "168h" # 7 days
  issuer: "golangiot"

cors:
  allowed_origins:
    - "http://localhost:3000"
    - "http://localhost:8080"
  allowed_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allowed_headers:
    - "Accept"
    - "Authorization"
    - "Content-Type"
    - "X-CSRF-Token"
  exposed_headers:
    - "Link"
  allow_credentials: true
  max_age: 300

rate_limit:
  enabled: true
  requests: 100
  window: "1m"

log:
  level: "info"
  format: "json"
  output: "stdout"
```

---

## 5. Domain Layer (Entities & Interfaces)

### 5.1 User Entity (internal/domain/entity/user.go)

```go
package entity

import (
    "time"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
    Email     string    `gorm:"uniqueIndex;not null" json:"email"`
    Password  string    `gorm:"not null" json:"-"`
    Name      string    `gorm:"not null" json:"name"`
    Role      string    `gorm:"default:'user'" json:"role"`
    Active    bool      `gorm:"default:true" json:"active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type UserRole string

const (
    RoleUser  UserRole = "user"
    RoleAdmin UserRole = "admin"
)

func (u *User) HashPassword() error {
    hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashed)
    return nil
}

func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    return err == nil
}

func (u *User) IsAdmin() bool {
    return u.Role == string(RoleAdmin)
}

type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required,min=2,max=100"`
}

type UpdateUserRequest struct {
    Name  string `json:"name" validate:"omitempty,min=2,max=100"`
    Email string `json:"email" validate:"omitempty,email"`
}

type UserResponse struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    Role      string    `json:"role"`
    Active    bool      `json:"active"`
    CreatedAt time.Time `json:"created_at"`
}
```

### 5.2 Token Entity (internal/domain/entity/token.go)

```go
package entity

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Token struct {
    ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
    UserID    string    `gorm:"not null;index" json:"user_id"`
    Token     string    `gorm:"uniqueIndex;not null" json:"token"`
    Type      TokenType `gorm:"not null" json:"type"`
    ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
    Revoked   bool      `gorm:"default:false" json:"revoked"`
    CreatedAt time.Time `json:"created_at"`
}

type TokenType string

const (
    TokenTypeAccess  TokenType = "access"
    TokenTypeRefresh TokenType = "refresh"
)

type JWTClaims struct {
    UserID string   `json:"user_id"`
    Email  string   `json:"email"`
    Role   string   `json:"role"`
    jwt.RegisteredClaims
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int64  `json:"expires_in"`
    User         UserResponse `json:"user"`
}

type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token" validate:"required"`
}
```

---

## 6. Repository Layer

### 6.1 User Repository Interface (internal/domain/repository/user_repository.go)

```go
package repository

import (
    "context"
    "golangiot/internal/domain/entity"
)

type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id string) error
    FindByID(ctx context.Context, id string) (*entity.User, error)
    FindByEmail(ctx context.Context, email string) (*entity.User, error)
    FindAll(ctx context.Context, filter *UserFilter) ([]*entity.User, int64, error)
    ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type UserFilter struct {
    Page   int    `json:"page"`
    Limit  int    `json:"limit"`
    Search string `json:"search"`
    Role   string `json:"role"`
    Active *bool  `json:"active"`
}
```

### 6.2 User Repository Implementation (internal/repository/postgres/user_repository.go)

```go
package postgres

import (
    "context"
    "errors"
    "gorm.io/gorm"
    
    "golangiot/internal/domain/entity"
    "golangiot/internal/domain/repository"
)

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
    if err := user.HashPassword(); err != nil {
        return err
    }
    
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
    return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
    return r.db.WithContext(ctx).Delete(&entity.User{}, "id = ?", id).Error
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
    var user entity.User
    err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
    var user entity.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) FindAll(ctx context.Context, filter *repository.UserFilter) ([]*entity.User, int64, error) {
    var users []*entity.User
    var total int64
    
    query := r.db.WithContext(ctx).Model(&entity.User{})
    
    if filter.Search != "" {
        query = query.Where("name ILIKE ? OR email ILIKE ?", 
            "%"+filter.Search+"%", "%"+filter.Search+"%")
    }
    
    if filter.Role != "" {
        query = query.Where("role = ?", filter.Role)
    }
    
    if filter.Active != nil {
        query = query.Where("active = ?", *filter.Active)
    }
    
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    offset := (filter.Page - 1) * filter.Limit
    if err := query.Offset(offset).
        Limit(filter.Limit).
        Order("created_at DESC").
        Find(&users).Error; err != nil {
        return nil, 0, err
    }
    
    return users, total, nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).Model(&entity.User{}).
        Where("email = ?", email).
        Count(&count).Error
    return count > 0, err
}
```

### 6.3 Token Repository (internal/repository/postgres/token_repository.go)

```go
package postgres

import (
    "context"
    "time"
    "gorm.io/gorm"
    
    "golangiot/internal/domain/entity"
    "golangiot/internal/domain/repository"
)

type tokenRepository struct {
    db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) repository.TokenRepository {
    return &tokenRepository{db: db}
}

func (r *tokenRepository) Create(ctx context.Context, token *entity.Token) error {
    return r.db.WithContext(ctx).Create(token).Error
}

func (r *tokenRepository) Update(ctx context.Context, token *entity.Token) error {
    return r.db.WithContext(ctx).Save(token).Error
}

func (r *tokenRepository) FindByToken(ctx context.Context, token string) (*entity.Token, error) {
    var t entity.Token
    err := r.db.WithContext(ctx).Where("token = ? AND revoked = ?", token, false).
        First(&t).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &t, nil
}

func (r *tokenRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
    return r.db.WithContext(ctx).Model(&entity.Token{}).
        Where("user_id = ?", userID).
        Update("revoked", true).Error
}

func (r *tokenRepository) DeleteExpiredTokens(ctx context.Context) error {
    return r.db.WithContext(ctx).
        Where("expires_at < ?", time.Now()).
        Delete(&entity.Token{}).Error
}
```

---

## 7. Service Layer

### 7.1 Auth Service (internal/service/auth_service.go)

```go
package service

import (
    "context"
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    
    "golangiot/configs"
    "golangiot/internal/domain/entity"
    "golangiot/internal/domain/repository"
    "golangiot/pkg/cache"
)

type AuthService struct {
    userRepo  repository.UserRepository
    tokenRepo repository.TokenRepository
    cache     cache.Cache
    jwtCfg    configs.JWTConfig
}

func NewAuthService(
    userRepo repository.UserRepository,
    tokenRepo repository.TokenRepository,
    cache cache.Cache,
    jwtCfg configs.JWTConfig,
) *AuthService {
    return &AuthService{
        userRepo:  userRepo,
        tokenRepo: tokenRepo,
        cache:     cache,
        jwtCfg:    jwtCfg,
    }
}

func (s *AuthService) Register(ctx context.Context, req *entity.CreateUserRequest) (*entity.UserResponse, error) {
    // Check if user exists
    exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, errors.New("user already exists")
    }
    
    // Create user
    user := &entity.User{
        Email:    req.Email,
        Password: req.Password,
        Name:     req.Name,
        Role:     string(entity.RoleUser),
        Active:   true,
    }
    
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    return &entity.UserResponse{
        ID:        user.ID,
        Email:     user.Email,
        Name:      user.Name,
        Role:      user.Role,
        Active:    user.Active,
        CreatedAt: user.CreatedAt,
    }, nil
}

func (s *AuthService) Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error) {
    // Find user by email
    user, err := s.userRepo.FindByEmail(ctx, req.Email)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("invalid credentials")
    }
    
    // Check password
    if !user.CheckPassword(req.Password) {
        return nil, errors.New("invalid credentials")
    }
    
    // Check if user is active
    if !user.Active {
        return nil, errors.New("user account is deactivated")
    }
    
    // Generate tokens
    accessToken, expiresIn, err := s.generateAccessToken(user)
    if err != nil {
        return nil, err
    }
    
    refreshToken, err := s.generateRefreshToken(user)
    if err != nil {
        return nil, err
    }
    
    return &entity.LoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    expiresIn,
        User: entity.UserResponse{
            ID:        user.ID,
            Email:     user.Email,
            Name:      user.Name,
            Role:      user.Role,
            Active:    user.Active,
            CreatedAt: user.CreatedAt,
        },
    }, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*entity.LoginResponse, error) {
    // Validate refresh token
    token, err := s.tokenRepo.FindByToken(ctx, refreshToken)
    if err != nil {
        return nil, err
    }
    if token == nil {
        return nil, errors.New("invalid refresh token")
    }
    
    // Check if token is expired
    if time.Now().After(token.ExpiresAt) {
        return nil, errors.New("refresh token expired")
    }
    
    // Get user
    user, err := s.userRepo.FindByID(ctx, token.UserID)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("user not found")
    }
    
    // Revoke old refresh token
    if err := s.tokenRepo.RevokeAllUserTokens(ctx, user.ID); err != nil {
        return nil, err
    }
    
    // Generate new tokens
    accessToken, expiresIn, err := s.generateAccessToken(user)
    if err != nil {
        return nil, err
    }
    
    newRefreshToken, err := s.generateRefreshToken(user)
    if err != nil {
        return nil, err
    }
    
    return &entity.LoginResponse{
        AccessToken:  accessToken,
        RefreshToken: newRefreshToken,
        ExpiresIn:    expiresIn,
        User: entity.UserResponse{
            ID:        user.ID,
            Email:     user.Email,
            Name:      user.Name,
            Role:      user.Role,
            Active:    user.Active,
            CreatedAt: user.CreatedAt,
        },
    }, nil
}

func (s *AuthService) Logout(ctx context.Context, userID string) error {
    return s.tokenRepo.RevokeAllUserTokens(ctx, userID)
}

func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*entity.JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &entity.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(s.jwtCfg.Secret), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*entity.JWTClaims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}

func (s *AuthService) generateAccessToken(user *entity.User) (string, int64, error) {
    expiresAt := time.Now().Add(s.jwtCfg.AccessExpire)
    expiresIn := int64(s.jwtCfg.AccessExpire.Seconds())
    
    claims := &entity.JWTClaims{
        UserID: user.ID,
        Email:  user.Email,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    s.jwtCfg.Issuer,
            ID:        uuid.New().String(),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(s.jwtCfg.Secret))
    if err != nil {
        return "", 0, err
    }
    
    return tokenString, expiresIn, nil
}

func (s *AuthService) generateRefreshToken(user *entity.User) (string, error) {
    expiresAt := time.Now().Add(s.jwtCfg.RefreshExpire)
    
    tokenString := uuid.New().String()
    
    token := &entity.Token{
        ID:        uuid.New().String(),
        UserID:    user.ID,
        Token:     tokenString,
        Type:      entity.TokenTypeRefresh,
        ExpiresAt: expiresAt,
        Revoked:   false,
    }
    
    if err := s.tokenRepo.Create(context.Background(), token); err != nil {
        return "", err
    }
    
    return tokenString, nil
}
```

### 7.2 User Service (internal/service/user_service.go)

```go
package service

import (
    "context"
    "errors"
    "fmt"
    "time"
    
    "golangiot/internal/domain/entity"
    "golangiot/internal/domain/repository"
    "golangiot/pkg/cache"
    "golangiot/pkg/queue"
)

type UserService struct {
    userRepo repository.UserRepository
    cache    cache.Cache
    queue    queue.Queue
}

func NewUserService(
    userRepo repository.UserRepository,
    cache cache.Cache,
    queue queue.Queue,
) *UserService {
    return &UserService{
        userRepo: userRepo,
        cache:    cache,
        queue:    queue,
    }
}

func (s *UserService) Create(ctx context.Context, req *entity.CreateUserRequest) (*entity.UserResponse, error) {
    // Check if user exists
    exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, errors.New("user already exists")
    }
    
    // Create user
    user := &entity.User{
        Email:    req.Email,
        Password: req.Password,
        Name:     req.Name,
        Role:     string(entity.RoleUser),
        Active:   true,
    }
    
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    // Publish event
    event := map[string]interface{}{
        "event":   "user.created",
        "user_id": user.ID,
        "email":   user.Email,
        "time":    time.Now(),
    }
    s.queue.Publish(ctx, "user.events", event)
    
    return &entity.UserResponse{
        ID:        user.ID,
        Email:     user.Email,
        Name:      user.Name,
        Role:      user.Role,
        Active:    user.Active,
        CreatedAt: user.CreatedAt,
    }, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*entity.UserResponse, error) {
    // Try cache first
    cacheKey := fmt.Sprintf("user:%s", id)
    var user entity.UserResponse
    
    err := s.cache.Get(ctx, cacheKey, &user)
    if err == nil {
        return &user, nil
    }
    
    // Get from database
    userEntity, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if userEntity == nil {
        return nil, errors.New("user not found")
    }
    
    response := &entity.UserResponse{
        ID:        userEntity.ID,
        Email:     userEntity.Email,
        Name:      userEntity.Name,
        Role:      userEntity.Role,
        Active:    userEntity.Active,
        CreatedAt: userEntity.CreatedAt,
    }
    
    // Store in cache
    go s.cache.Set(context.Background(), cacheKey, response, 5*time.Minute)
    
    return response, nil
}

func (s *UserService) Update(ctx context.Context, id string, req *entity.UpdateUserRequest) (*entity.UserResponse, error) {
    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("user not found")
    }
    
    if req.Name != "" {
        user.Name = req.Name
    }
    
    if req.Email != "" {
        // Check if email is taken
        exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
        if err != nil {
            return nil, err
        }
        if exists && user.Email != req.Email {
            return nil, errors.New("email already taken")
        }
        user.Email = req.Email
    }
    
    if err := s.userRepo.Update(ctx, user); err != nil {
        return nil, err
    }
    
    // Invalidate cache
    cacheKey := fmt.Sprintf("user:%s", id)
    go s.cache.Delete(context.Background(), cacheKey)
    
    return &entity.UserResponse{
        ID:        user.ID,
        Email:     user.Email,
        Name:      user.Name,
        Role:      user.Role,
        Active:    user.Active,
        CreatedAt: user.CreatedAt,
    }, nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        return err
    }
    if user == nil {
        return errors.New("user not found")
    }
    
    if err := s.userRepo.Delete(ctx, id); err != nil {
        return err
    }
    
    // Invalidate cache
    cacheKey := fmt.Sprintf("user:%s", id)
    go s.cache.Delete(context.Background(), cacheKey)
    
    return nil
}

func (s *UserService) List(ctx context.Context, filter *repository.UserFilter) ([]*entity.UserResponse, int64, error) {
    users, total, err := s.userRepo.FindAll(ctx, filter)
    if err != nil {
        return nil, 0, err
    }
    
    responses := make([]*entity.UserResponse, len(users))
    for i, user := range users {
        responses[i] = &entity.UserResponse{
            ID:        user.ID,
            Email:     user.Email,
            Name:      user.Name,
            Role:      user.Role,
            Active:    user.Active,
            CreatedAt: user.CreatedAt,
        }
    }
    
    return responses, total, nil
}
```

---

## 8. Handler Layer

### 8.1 Response Handler (internal/handler/response.go)

```go
package handler

import (
    "encoding/json"
    "net/http"
)

type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
    Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
    Page       int   `json:"page"`
    Limit      int   `json:"limit"`
    Total      int64 `json:"total"`
    TotalPages int64 `json:"total_pages"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func Success(w http.ResponseWriter, data interface{}) {
    JSON(w, http.StatusOK, Response{
        Success: true,
        Data:    data,
    })
}

func Created(w http.ResponseWriter, data interface{}) {
    JSON(w, http.StatusCreated, Response{
        Success: true,
        Data:    data,
    })
}

func Error(w http.ResponseWriter, status int, message string) {
    JSON(w, status, Response{
        Success: false,
        Error:   message,
    })
}

func BadRequest(w http.ResponseWriter, message string) {
    Error(w, http.StatusBadRequest, message)
}

func Unauthorized(w http.ResponseWriter, message string) {
    Error(w, http.StatusUnauthorized, message)
}

func Forbidden(w http.ResponseWriter, message string) {
    Error(w, http.StatusForbidden, message)
}

func NotFound(w http.ResponseWriter, message string) {
    Error(w, http.StatusNotFound, message)
}

func InternalError(w http.ResponseWriter, message string) {
    Error(w, http.StatusInternalServerError, message)
}

func PaginatedResponse(w http.ResponseWriter, data interface{}, page, limit int, total int64) {
    totalPages := (total + int64(limit) - 1) / int64(limit)
    
    JSON(w, http.StatusOK, Response{
        Success: true,
        Data:    data,
        Meta: &Meta{
            Page:       page,
            Limit:      limit,
            Total:      total,
            TotalPages: totalPages,
        },
    })
}
```

### 8.2 Auth Handler (internal/handler/auth_handler.go)

```go
package handler

import (
    "encoding/json"
    "net/http"
    
    "golangiot/internal/domain/entity"
    "golangiot/internal/service"
    "golangiot/pkg/logger"
    "golangiot/pkg/validator"
)

type AuthHandler struct {
    authService *service.AuthService
    logger      logger.Logger
    validator   *validator.Validator
}

func NewAuthHandler(authService *service.AuthService, logger logger.Logger) *AuthHandler {
    return &AuthHandler{
        authService: authService,
        logger:      logger,
        validator:   validator.NewValidator(),
    }
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req entity.CreateUserRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        BadRequest(w, "Invalid request body")
        return
    }
    
    if err := h.validator.Validate(&req); err != nil {
        BadRequest(w, err.Error())
        return
    }
    
    user, err := h.authService.Register(r.Context(), &req)
    if err != nil {
        h.logger.Error("Failed to register user", "error", err)
        BadRequest(w, err.Error())
        return
    }
    
    Created(w, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req entity.LoginRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        BadRequest(w, "Invalid request body")
        return
    }
    
    if err := h.validator.Validate(&req); err != nil {
        BadRequest(w, err.Error())
        return
    }
    
    response, err := h.authService.Login(r.Context(), &req)
    if err != nil {
        h.logger.Error("Failed to login", "error", err)
        Unauthorized(w, err.Error())
        return
    }
    
    Success(w, response)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
    var req entity.RefreshTokenRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        BadRequest(w, "Invalid request body")
        return
    }
    
    if err := h.validator.Validate(&req); err != nil {
        BadRequest(w, err.Error())
        return
    }
    
    response, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
    if err != nil {
        h.logger.Error("Failed to refresh token", "error", err)
        Unauthorized(w, err.Error())
        return
    }
    
    Success(w, response)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("user_id").(string)
    
    if err := h.authService.Logout(r.Context(), userID); err != nil {
        h.logger.Error("Failed to logout", "error", err)
        InternalError(w, "Failed to logout")
        return
    }
    
    Success(w, map[string]string{"message": "Logged out successfully"})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("user_id").(string)
    
    user, err := h.authService.GetUserByID(r.Context(), userID)
    if err != nil {
        h.logger.Error("Failed to get user", "error", err)
        NotFound(w, "User not found")
        return
    }
    
    Success(w, user)
}
```

### 8.3 User Handler (internal/handler/user_handler.go)

```go
package handler

import (
    "encoding/json"
    "net/http"
    "strconv"
    
    "github.com/go-chi/chi/v5"
    
    "golangiot/internal/domain/entity"
    "golangiot/internal/domain/repository"
    "golangiot/internal/service"
    "golangiot/pkg/logger"
    "golangiot/pkg/validator"
)

type UserHandler struct {
    userService *service.UserService
    logger      logger.Logger
    validator   *validator.Validator
}

func NewUserHandler(userService *service.UserService, logger logger.Logger) *UserHandler {
    return &UserHandler{
        userService: userService,
        logger:      logger,
        validator:   validator.NewValidator(),
    }
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req entity.CreateUserRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        BadRequest(w, "Invalid request body")
        return
    }
    
    if err := h.validator.Validate(&req); err != nil {
        BadRequest(w, err.Error())
        return
    }
    
    user, err := h.userService.Create(r.Context(), &req)
    if err != nil {
        h.logger.Error("Failed to create user", "error", err)
        BadRequest(w, err.Error())
        return
    }
    
    Created(w, user)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    user, err := h.userService.GetByID(r.Context(), id)
    if err != nil {
        h.logger.Error("Failed to get user", "error", err)
        NotFound(w, "User not found")
        return
    }
    
    Success(w, user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    var req entity.UpdateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        BadRequest(w, "Invalid request body")
        return
    }
    
    if err := h.validator.Validate(&req); err != nil {
        BadRequest(w, err.Error())
        return
    }
    
    user, err := h.userService.Update(r.Context(), id, &req)
    if err != nil {
        h.logger.Error("Failed to update user", "error", err)
        BadRequest(w, err.Error())
        return
    }
    
    Success(w, user)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    if err := h.userService.Delete(r.Context(), id); err != nil {
        h.logger.Error("Failed to delete user", "error", err)
        BadRequest(w, err.Error())
        return
    }
    
    Success(w, map[string]string{"message": "User deleted successfully"})
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    if page < 1 {
        page = 1
    }
    
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    if limit < 1 || limit > 100 {
        limit = 20
    }
    
    filter := &repository.UserFilter{
        Page:   page,
        Limit:  limit,
        Search: r.URL.Query().Get("search"),
        Role:   r.URL.Query().Get("role"),
    }
    
    users, total, err := h.userService.List(r.Context(), filter)
    if err != nil {
        h.logger.Error("Failed to list users", "error", err)
        InternalError(w, "Failed to list users")
        return
    }
    
    PaginatedResponse(w, users, page, limit, total)
}
```

---

## 9. Middleware

### 9.1 Auth Middleware (internal/middleware/auth.go)

```go
package middleware

import (
    "context"
    "net/http"
    "strings"
    
    "golangiot/configs"
    "golangiot/internal/service"
    "golangiot/pkg/logger"
)

type AuthMiddleware struct {
    authService *service.AuthService
    logger      logger.Logger
}

func NewAuthMiddleware(cfg configs.JWTConfig, logger logger.Logger) *AuthMiddleware {
    authService := service.NewAuthService(nil, nil, nil, cfg)
    return &AuthMiddleware{
        authService: authService,
        logger:      logger,
    }
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing authorization header", http.StatusUnauthorized)
            return
        }
        
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
            return
        }
        
        tokenString := parts[1]
        claims, err := m.authService.ValidateToken(r.Context(), tokenString)
        if err != nil {
            m.logger.Error("Token validation failed", "error", err)
            http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
            return
        }
        
        ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
        ctx = context.WithValue(ctx, "user_email", claims.Email)
        ctx = context.WithValue(ctx, "user_role", claims.Role)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func (m *AuthMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            userRole, ok := r.Context().Value("user_role").(string)
            if !ok {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            
            allowed := false
            for _, role := range roles {
                if userRole == role {
                    allowed = true
                    break
                }
            }
            
            if !allowed {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

### 9.2 Logger Middleware (internal/middleware/logger.go)

```go
package middleware

import (
    "net/http"
    "time"
    
    "github.com/go-chi/chi/v5/middleware"
    "golangiot/pkg/logger"
)

func NewLogger(log logger.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Wrap response writer to capture status code
            rw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
            
            next.ServeHTTP(rw, r)
            
            duration := time.Since(start)
            
            log.Info("HTTP Request",
                "method", r.Method,
                "path", r.URL.Path,
                "status", rw.Status(),
                "duration", duration,
                "ip", r.RemoteAddr,
                "user_agent", r.UserAgent(),
                "request_id", middleware.GetReqID(r.Context()),
            )
        })
    }
}
```

### 9.3 Rate Limit Middleware (internal/middleware/rate_limit.go)

```go
package middleware

import (
    "net/http"
    "sync"
    "time"
)

type RateLimiter struct {
    requests map[string][]time.Time
    mu       sync.RWMutex
    limit    int
    window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    rl := &RateLimiter{
        requests: make(map[string][]time.Time),
        limit:    limit,
        window:   window,
    }
    
    // Clean up old entries
    go rl.cleanup()
    
    return rl
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        
        if !rl.allow(ip) {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

func (rl *RateLimiter) allow(ip string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    now := time.Now()
    cutoff := now.Add(-rl.window)
    
    // Get existing requests
    requests := rl.requests[ip]
    
    // Filter out old requests
    valid := make([]time.Time, 0)
    for _, t := range requests {
        if t.After(cutoff) {
            valid = append(valid, t)
        }
    }
    
    if len(valid) >= rl.limit {
        return false
    }
    
    valid = append(valid, now)
    rl.requests[ip] = valid
    
    return true
}

func (rl *RateLimiter) cleanup() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        rl.mu.Lock()
        now := time.Now()
        cutoff := now.Add(-rl.window)
        
        for ip, requests := range rl.requests {
            valid := make([]time.Time, 0)
            for _, t := range requests {
                if t.After(cutoff) {
                    valid = append(valid, t)
                }
            }
            if len(valid) == 0 {
                delete(rl.requests, ip)
            } else {
                rl.requests[ip] = valid
            }
        }
        rl.mu.Unlock()
    }
}
```

### 9.4 Security Headers Middleware (internal/middleware/security.go)

```go
package middleware

import (
    "net/http"
)

func SecurityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Set security headers
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        
        next.ServeHTTP(w, r)
    })
}
```

### 9.5 Context Middleware (internal/middleware/context.go)

```go
package middleware

import (
    "context"
    "net/http"
    
    "github.com/google/uuid"
)

type contextKey string

const (
    RequestIDKey contextKey = "request_id"
    UserIDKey    contextKey = "user_id"
    UserEmailKey contextKey = "user_email"
    UserRoleKey  contextKey = "user_role"
)

func RequestContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Generate request ID if not present
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
        w.Header().Set("X-Request-ID", requestID)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func GetRequestID(ctx context.Context) string {
    if id, ok := ctx.Value(RequestIDKey).(string); ok {
        return id
    }
    return ""
}

func GetUserID(ctx context.Context) string {
    if id, ok := ctx.Value(UserIDKey).(string); ok {
        return id
    }
    return ""
}
```

---

## 10. Cache Layer (pkg/cache/cache.go)

```go
package cache

import (
    "context"
    "encoding/json"
    "time"
    
    "github.com/redis/go-redis/v9"
)

type Cache interface {
    Get(ctx context.Context, key string, dest interface{}) error
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
}

type RedisCache struct {
    client *redis.Client
}

func NewRedisCache(client *redis.Client) Cache {
    return &RedisCache{client: client}
}

func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
    val, err := c.client.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return ErrCacheMiss
        }
        return err
    }
    
    return json.Unmarshal([]byte(val), dest)
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    
    return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
    return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
    result, err := c.client.Exists(ctx, key).Result()
    if err != nil {
        return false, err
    }
    return result > 0, nil
}

var ErrCacheMiss = errors.New("cache: key not found")
```

---

## 11. Message Queue (pkg/queue/queue.go)

```go
package queue

import (
    "context"
    "encoding/json"
    "time"
    
    "github.com/redis/go-redis/v9"
)

type Queue interface {
    Publish(ctx context.Context, channel string, message interface{}) error
    Subscribe(ctx context.Context, channel string, handler MessageHandler) error
}

type MessageHandler func(ctx context.Context, message []byte) error

type RedisQueue struct {
    client *redis.Client
}

func NewRedisQueue(client *redis.Client) Queue {
    return &RedisQueue{client: client}
}

func (q *RedisQueue) Publish(ctx context.Context, channel string, message interface{}) error {
    data, err := json.Marshal(message)
    if err != nil {
        return err
    }
    
    return q.client.Publish(ctx, channel, data).Err()
}

func (q *RedisQueue) Subscribe(ctx context.Context, channel string, handler MessageHandler) error {
    pubsub := q.client.Subscribe(ctx, channel)
    defer pubsub.Close()
    
    ch := pubsub.Channel()
    
    for {
        select {
        case msg := <-ch:
            if err := handler(ctx, []byte(msg.Payload)); err != nil {
                // Log error but continue
                continue
            }
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

// Worker Pool Implementation
type WorkerPool struct {
    queue     Queue
    workers   int
    handler   MessageHandler
    ctx       context.Context
    cancel    context.CancelFunc
}

func NewWorkerPool(queue Queue, workers int, handler MessageHandler) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    return &WorkerPool{
        queue:   queue,
        workers: workers,
        handler: handler,
        ctx:     ctx,
        cancel:  cancel,
    }
}

func (p *WorkerPool) Start(channel string) {
    for i := 0; i < p.workers; i++ {
        go p.worker(channel)
    }
}

func (p *WorkerPool) worker(channel string) {
    p.queue.Subscribe(p.ctx, channel, p.handler)
}

func (p *WorkerPool) Stop() {
    p.cancel()
}
```

---

## 12. Health Monitoring (internal/handler/health_handler.go)

```go
package handler

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "time"
    
    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"
    
    "golangiot/pkg/logger"
)

type HealthHandler struct {
    db          *gorm.DB
    redisClient *redis.Client
    logger      logger.Logger
}

func NewHealthHandler(db *gorm.DB, redisClient *redis.Client, logger logger.Logger) *HealthHandler {
    return &HealthHandler{
        db:          db,
        redisClient: redisClient,
        logger:      logger,
    }
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
    Success(w, map[string]string{
        "status": "healthy",
        "time":   time.Now().UTC().Format(time.RFC3339),
    })
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
    // Check database connection
    sqlDB, err := h.db.DB()
    if err != nil || sqlDB.Ping() != nil {
        http.Error(w, "Database not ready", http.StatusServiceUnavailable)
        return
    }
    
    // Check Redis connection
    if err := h.redisClient.Ping(r.Context()).Err(); err != nil {
        http.Error(w, "Redis not ready", http.StatusServiceUnavailable)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "alive"})
}

func (h *HealthHandler) DetailedHealth(w http.ResponseWriter, r *http.Request) {
    health := map[string]interface{}{
        "status":    "healthy",
        "timestamp": time.Now().UTC().Format(time.RFC3339),
        "checks":    make(map[string]interface{}),
    }
    
    checks := health["checks"].(map[string]interface{})
    
    // Check database
    if err := h.checkDatabase(); err != nil {
        checks["database"] = map[string]interface{}{
            "status": "unhealthy",
            "error":  err.Error(),
        }
        health["status"] = "unhealthy"
    } else {
        checks["database"] = map[string]interface{}{
            "status": "healthy",
        }
    }
    
    // Check Redis
    if err := h.checkRedis(r.Context()); err != nil {
        checks["redis"] = map[string]interface{}{
            "status": "unhealthy",
            "error":  err.Error(),
        }
        health["status"] = "unhealthy"
    } else {
        checks["redis"] = map[string]interface{}{
            "status": "healthy",
        }
    }
    
    w.Header().Set("Content-Type", "application/json")
    if health["status"] == "unhealthy" {
        w.WriteHeader(http.StatusServiceUnavailable)
    } else {
        w.WriteHeader(http.StatusOK)
    }
    
    json.NewEncoder(w).Encode(health)
}

func (h *HealthHandler) checkDatabase() error {
    sqlDB, err := h.db.DB()
    if err != nil {
        return err
    }
    return sqlDB.Ping()
}

func (h *HealthHandler) checkRedis(ctx context.Context) error {
    return h.redisClient.Ping(ctx).Err()
}
```

---

## 13. Logger Implementation (pkg/logger/logger.go)

```go
package logger

import (
    "io"
    "log/slog"
    "os"
)

type Logger interface {
    Debug(msg string, args ...interface{})
    Info(msg string, args ...interface{})
    Warn(msg string, args ...interface{})
    Error(msg string, args ...interface{})
    Fatal(msg string, args ...interface{})
    With(args ...interface{}) Logger
    Sync() error
}

type Config struct {
    Level  string
    Format string
    Output string
}

type SlogLogger struct {
    logger *slog.Logger
}

func NewLogger(cfg *Config) (Logger, error) {
    var level slog.Level
    switch cfg.Level {
    case "debug":
        level = slog.LevelDebug
    case "info":
        level = slog.LevelInfo
    case "warn":
        level = slog.LevelWarn
    case "error":
        level = slog.LevelError
    default:
        level = slog.LevelInfo
    }
    
    var handler slog.Handler
    opts := &slog.HandlerOptions{Level: level}
    
    var output io.Writer
    if cfg.Output == "stdout" {
        output = os.Stdout
    } else {
        output = os.Stderr
    }
    
    if cfg.Format == "json" {
        handler = slog.NewJSONHandler(output, opts)
    } else {
        handler = slog.NewTextHandler(output, opts)
    }
    
    logger := slog.New(handler)
    
    return &SlogLogger{logger: logger}, nil
}

func (l *SlogLogger) Debug(msg string, args ...interface{}) {
    l.logger.Debug(msg, args...)
}

func (l *SlogLogger) Info(msg string, args ...interface{}) {
    l.logger.Info(msg, args...)
}

func (l *SlogLogger) Warn(msg string, args ...interface{}) {
    l.logger.Warn(msg, args...)
}

func (l *SlogLogger) Error(msg string, args ...interface{}) {
    l.logger.Error(msg, args...)
}

func (l *SlogLogger) Fatal(msg string, args ...interface{}) {
    l.logger.Error(msg, args...)
    os.Exit(1)
}

func (l *SlogLogger) With(args ...interface{}) Logger {
    return &SlogLogger{logger: l.logger.With(args...)}
}

func (l *SlogLogger) Sync() error {
    return nil
}
```

---

## 14. Database Migrations

### 14.1 Migration Files

```sql
-- migrations/001_create_users_table.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- migrations/002_create_tokens_table.sql
CREATE TABLE tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tokens_user_id ON tokens(user_id);
CREATE INDEX idx_tokens_token ON tokens(token);
CREATE INDEX idx_tokens_expires_at ON tokens(expires_at);

-- migrations/003_create_audit_logs_table.sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    resource VARCHAR(100),
    resource_id VARCHAR(100),
    old_data JSONB,
    new_data JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
```

### 14.2 Migration Script (scripts/migrate.sh)

```bash
#!/bin/bash

# Load environment variables
source .env

# Run migrations
migrate -path migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" up

# Seed database
go run scripts/seed.go
```

---

## 15. Makefile

```makefile
.PHONY: help build run test lint migrate-up migrate-down seed docker-build docker-run clean

help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make lint         - Run linter"
	@echo "  make migrate-up   - Run database migrations"
	@echo "  make migrate-down - Rollback database migrations"
	@echo "  make seed         - Seed database"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run with Docker Compose"
	@echo "  make clean        - Clean build artifacts"

build:
	go build -o bin/api cmd/api/main.go

run:
	go run cmd/api/main.go

test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

migrate-up:
	migrate -path migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" up

migrate-down:
	migrate -path migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" down

seed:
	go run scripts/seed.go

docker-build:
	docker build -t golangiot:latest -f deployments/docker/Dockerfile .

docker-run:
	docker-compose -f deployments/docker/docker-compose.yml up

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
```

---

## 16. Docker Deployment

### 16.1 Dockerfile (deployments/docker/Dockerfile)

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/api cmd/api/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary and configs
COPY --from=builder /app/bin/api /app/
COPY --from=builder /app/configs /app/configs

# Change ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run application
ENTRYPOINT ["/app/api"]
CMD ["--config", "/app/configs/config.yaml"]
```

### 16.2 Docker Compose (deployments/docker/docker-compose.yml)

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: golangiot
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
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

  api:
    build:
      context: ../..
      dockerfile: deployments/docker/Dockerfile
    ports:
      - "8080:8080"
    environment:
      APP_ENV: production
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: postgres
      DB_NAME: golangiot
      REDIS_HOST: redis
      REDIS_PORT: 6379
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    volumes:
      - ../../configs:/app/configs

volumes:
  postgres_data:
  redis_data:
```

---

## 17. Testing

### 17.1 Unit Test Example (internal/service/auth_service_test.go)

```go
package service_test

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    
    "golangiot/internal/domain/entity"
    "golangiot/internal/service"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.User), args.Error(1)
}

func TestAuthService_Register(t *testing.T) {
    mockUserRepo := new(MockUserRepository)
    mockTokenRepo := new(MockTokenRepository)
    mockCache := new(MockCache)
    
    authService := service.NewAuthService(mockUserRepo, mockTokenRepo, mockCache, configs.JWTConfig{})
    
    ctx := context.Background()
    req := &entity.CreateUserRequest{
        Email:    "test@example.com",
        Password: "password123",
        Name:     "Test User",
    }
    
    mockUserRepo.On("ExistsByEmail", ctx, req.Email).Return(false, nil)
    mockUserRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)
    
    user, err := authService.Register(ctx, req)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, user.ID)
    assert.Equal(t, req.Email, user.Email)
    assert.Equal(t, req.Name, user.Name)
    
    mockUserRepo.AssertExpectations(t)
}
```

---

## สรุป

Golangiot เป็น boilerplate ที่พร้อมใช้งานจริง (production-ready) ครบวงจรด้วย:

1. **Clean Architecture**: แยก层ชัดเจน (Handler → Service → Repository)
2. **Security**: JWT, rate limiting, security headers, input validation
3. **Performance**: Redis cache, connection pooling, worker pools
4. **Observability**: Structured logging (slog), health checks, metrics
5. **Database**: GORM with PostgreSQL, migrations, connection management
6. **Message Queue**: Redis pub/sub with worker pools
7. **Testing**: Unit tests, integration tests, mock support
8. **Deployment**: Docker, Docker Compose, Kubernetes ready
