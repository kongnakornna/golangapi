# คู่มือ Golang สำหรับการใช้งานจริง (Production-Ready)

## สารบัญ
1. [การออกแบบสถาปัตยกรรม](#1-การออกแบบสถาปัตยกรรม)
2. [การจัดการ Project Structure](#2-การจัดการ-project-structure)
3. [การจัดการ Configuration](#3-การจัดการ-configuration)
4. [Database Integration](#4-database-integration)
5. [Logging และ Monitoring](#5-logging-และ-monitoring)
6. [Middleware และ Authentication](#6-middleware-และ-authentication)
7. [Testing Strategy](#7-testing-strategy)
8. [Performance Optimization](#8-performance-optimization)
9. [Deployment และ Containerization](#9-deployment-และ-containerization)
10. [Security Best Practices](#10-security-best-practices)

---

## 1. การออกแบบสถาปัตยกรรม

### Clean Architecture Implementation

```go
// domain/entity/todo.go
package entity

import "time"

type Todo struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Status    TodoStatus `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type TodoStatus string

const (
    TodoStatusPending    TodoStatus = "pending"
    TodoStatusInProgress TodoStatus = "in_progress"
    TodoStatusCompleted  TodoStatus = "completed"
    TodoStatusCancelled  TodoStatus = "cancelled"
)

// domain/repository/todo_repository.go
package repository

import (
    "context"
    "your-project/domain/entity"
)

type TodoRepository interface {
    Create(ctx context.Context, todo *entity.Todo) error
    Update(ctx context.Context, todo *entity.Todo) error
    Delete(ctx context.Context, id string) error
    FindByID(ctx context.Context, id string) (*entity.Todo, error)
    FindAll(ctx context.Context, filter *TodoFilter) ([]*entity.Todo, int64, error)
}

// domain/service/todo_service.go
package service

import (
    "context"
    "errors"
    "your-project/domain/entity"
    "your-project/domain/repository"
)

type TodoService struct {
    repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) *TodoService {
    return &TodoService{repo: repo}
}

func (s *TodoService) CreateTodo(ctx context.Context, title, content string) (*entity.Todo, error) {
    if title == "" {
        return nil, errors.New("title is required")
    }
    
    todo := &entity.Todo{
        ID:        generateID(),
        Title:     title,
        Content:   content,
        Status:    entity.TodoStatusPending,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    if err := s.repo.Create(ctx, todo); err != nil {
        return nil, err
    }
    
    return todo, nil
}

// usecase/todo_usecase.go
package usecase

import (
    "context"
    "your-project/domain/entity"
)

type TodoUseCase interface {
    CreateTodo(ctx context.Context, req *CreateTodoRequest) (*CreateTodoResponse, error)
    GetTodo(ctx context.Context, id string) (*GetTodoResponse, error)
    UpdateTodoStatus(ctx context.Context, id string, status string) error
    ListTodos(ctx context.Context, req *ListTodoRequest) (*ListTodoResponse, error)
}

type todoUseCase struct {
    todoService *service.TodoService
    logger      Logger
    cache       Cache
}

func NewTodoUseCase(
    todoService *service.TodoService,
    logger Logger,
    cache Cache,
) TodoUseCase {
    return &todoUseCase{
        todoService: todoService,
        logger:      logger,
        cache:       cache,
    }
}

func (uc *todoUseCase) CreateTodo(ctx context.Context, req *CreateTodoRequest) (*CreateTodoResponse, error) {
    // Validation
    if err := req.Validate(); err != nil {
        return nil, err
    }
    
    // Business logic
    todo, err := uc.todoService.CreateTodo(ctx, req.Title, req.Content)
    if err != nil {
        uc.logger.Error("failed to create todo", "error", err)
        return nil, err
    }
    
    // Invalidate cache
    uc.cache.Delete(ctx, "todos:list")
    
    return &CreateTodoResponse{
        ID:      todo.ID,
        Title:   todo.Title,
        Status:  string(todo.Status),
        CreatedAt: todo.CreatedAt,
    }, nil
}
```

### Dependency Injection ด้วย Wire

```go
// wire.go
//go:build wireinject
// +build wireinject

package main

import (
    "github.com/google/wire"
    "your-project/internal/handler"
    "your-project/internal/repository"
    "your-project/internal/service"
    "your-project/internal/usecase"
)

func InitializeServer() (*Server, error) {
    wire.Build(
        // Database
        NewDatabase,
        
        // Repository
        repository.NewTodoRepository,
        repository.NewUserRepository,
        
        // Service
        service.NewTodoService,
        service.NewAuthService,
        
        // UseCase
        usecase.NewTodoUseCase,
        usecase.NewAuthUseCase,
        
        // Handler
        handler.NewTodoHandler,
        handler.NewAuthHandler,
        
        // Server
        NewServer,
    )
    return nil, nil
}
```

---

## 2. การจัดการ Project Structure

### Standard Project Layout

```
your-project/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── worker/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── entity/
│   │   ├── repository/
│   │   └── service/
│   ├── usecase/
│   ├── handler/
│   │   ├── http/
│   │   └── grpc/
│   ├── repository/
│   │   ├── postgres/
│   │   ├── mongodb/
│   │   └── redis/
│   └── middleware/
├── pkg/
│   ├── logger/
│   ├── errors/
│   ├── validator/
│   ├── cache/
│   └── metrics/
├── api/
│   ├── openapi/
│   │   └── swagger.yaml
│   └── proto/
│       └── todo.proto
├── configs/
│   ├── config.yaml
│   └── config.go
├── scripts/
│   ├── migration/
│   └── seed/
├── test/
│   ├── integration/
│   └── e2e/
├── deployments/
│   ├── docker/
│   └── k8s/
├── docs/
├── go.mod
├── go.sum
├── Makefile
└── .env.example
```

### Makefile สำหรับ Automation

```makefile
.PHONY: help build run test lint migrate docker-build

help:
	@echo "Available commands:"
	@echo "  make build      - Build the application"
	@echo "  make run        - Run the application"
	@echo "  make test       - Run tests"
	@echo "  make lint       - Run linter"
	@echo "  make migrate    - Run database migrations"

build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/worker cmd/worker/main.go

run:
	go run cmd/api/main.go

test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

migrate-up:
	migrate -path scripts/migration -database "postgresql://localhost:5432/db" up

migrate-down:
	migrate -path scripts/migration -database "postgresql://localhost:5432/db" down

docker-build:
	docker build -t your-app:latest -f deployments/docker/Dockerfile .

docker-run:
	docker-compose -f deployments/docker/docker-compose.yml up

gen-proto:
	protoc --go_out=. --go-grpc_out=. api/proto/*.proto

swagger:
	swag init -g cmd/api/main.go -o api/openapi
```

---

## 3. การจัดการ Configuration

### Multi-Environment Configuration

```go
// configs/config.go
package config

import (
    "fmt"
    "time"
    "github.com/spf13/viper"
    "github.com/caarlos0/env/v8"
)

type Config struct {
    App      AppConfig      `mapstructure:"app"`
    HTTP     HTTPConfig     `mapstructure:"http"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    Log      LogConfig      `mapstructure:"log"`
}

type AppConfig struct {
    Name    string `mapstructure:"name" env:"APP_NAME"`
    Version string `mapstructure:"version" env:"APP_VERSION"`
    Env     string `mapstructure:"env" env:"APP_ENV" envDefault:"development"`
    Debug   bool   `mapstructure:"debug" env:"APP_DEBUG"`
}

type HTTPConfig struct {
    Host            string        `mapstructure:"host" env:"HTTP_HOST" envDefault:"0.0.0.0"`
    Port            int           `mapstructure:"port" env:"HTTP_PORT" envDefault:"8080"`
    ReadTimeout     time.Duration `mapstructure:"read_timeout" env:"HTTP_READ_TIMEOUT" envDefault:"30s"`
    WriteTimeout    time.Duration `mapstructure:"write_timeout" env:"HTTP_WRITE_TIMEOUT" envDefault:"30s"`
    IdleTimeout     time.Duration `mapstructure:"idle_timeout" env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
    MaxHeaderBytes  int           `mapstructure:"max_header_bytes" envDefault:"1048576"`
    GracefulTimeout time.Duration `mapstructure:"graceful_timeout" envDefault:"15s"`
}

type DatabaseConfig struct {
    Driver          string        `mapstructure:"driver" env:"DB_DRIVER" envDefault:"postgres"`
    Host            string        `mapstructure:"host" env:"DB_HOST" envDefault:"localhost"`
    Port            int           `mapstructure:"port" env:"DB_PORT" envDefault:"5432"`
    Username        string        `mapstructure:"username" env:"DB_USERNAME"`
    Password        string        `mapstructure:"password" env:"DB_PASSWORD"`
    Database        string        `mapstructure:"database" env:"DB_DATABASE"`
    SSLMode         string        `mapstructure:"ssl_mode" env:"DB_SSL_MODE" envDefault:"disable"`
    MaxOpenConns    int           `mapstructure:"max_open_conns" envDefault:"100"`
    MaxIdleConns    int           `mapstructure:"max_idle_conns" envDefault:"10"`
    ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" envDefault:"1h"`
}

type JWTConfig struct {
    Secret        string        `mapstructure:"secret" env:"JWT_SECRET"`
    AccessExpire  time.Duration `mapstructure:"access_expire" env:"JWT_ACCESS_EXPIRE" envDefault:"15m"`
    RefreshExpire time.Duration `mapstructure:"refresh_expire" env:"JWT_REFRESH_EXPIRE" envDefault:"7d"`
    Issuer        string        `mapstructure:"issuer" env:"JWT_ISSUER" envDefault:"your-app"`
}

func LoadConfig(path string) (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(path)
    viper.AddConfigPath(".")
    viper.AddConfigPath("./configs")
    
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
}

// configs/config.yaml
/*
app:
  name: "todo-api"
  version: "1.0.0"
  env: "development"
  debug: true

http:
  host: "0.0.0.0"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"
  graceful_timeout: "15s"

database:
  driver: "postgres"
  host: "localhost"
  port: 5432
  username: "postgres"
  password: "postgres"
  database: "todo_db"
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
  secret: "your-secret-key"
  access_expire: "15m"
  refresh_expire: "168h" # 7 days

log:
  level: "debug"
  format: "json"
  output: "stdout"
*/
```

---

## 4. Database Integration

### PostgreSQL with GORM และ Connection Pool

```go
// internal/repository/postgres/database.go
package postgres

import (
    "context"
    "fmt"
    "time"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

type Database struct {
    DB *gorm.DB
}

func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode)
    
    gormConfig := &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
    }
    
    db, err := gorm.Open(postgres.Open(dsn), gormConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to connect database: %w", err)
    }
    
    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get sql.DB: %w", err)
    }
    
    // Connection pool configuration
    sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
    sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
    sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
    
    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := sqlDB.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    return &Database{DB: db}, nil
}

// Repository Implementation with Transaction Support
type todoRepository struct {
    db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) repository.TodoRepository {
    return &todoRepository{db: db}
}

func (r *todoRepository) Create(ctx context.Context, todo *entity.Todo) error {
    return r.db.WithContext(ctx).Create(todo).Error
}

func (r *todoRepository) Update(ctx context.Context, todo *entity.Todo) error {
    return r.db.WithContext(ctx).Save(todo).Error
}

func (r *todoRepository) FindByID(ctx context.Context, id string) (*entity.Todo, error) {
    var todo entity.Todo
    err := r.db.WithContext(ctx).
        Where("id = ? AND deleted_at IS NULL", id).
        First(&todo).Error
    
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    
    return &todo, nil
}

func (r *todoRepository) FindAll(ctx context.Context, filter *repository.TodoFilter) ([]*entity.Todo, int64, error) {
    var todos []*entity.Todo
    var total int64
    
    query := r.db.WithContext(ctx).Model(&entity.Todo{}).
        Where("deleted_at IS NULL")
    
    if filter.Status != "" {
        query = query.Where("status = ?", filter.Status)
    }
    
    if filter.Search != "" {
        query = query.Where("title LIKE ? OR content LIKE ?", 
            "%"+filter.Search+"%", "%"+filter.Search+"%")
    }
    
    // Get total count
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Pagination
    offset := (filter.Page - 1) * filter.Limit
    if err := query.Offset(offset).
        Limit(filter.Limit).
        Order("created_at DESC").
        Find(&todos).Error; err != nil {
        return nil, 0, err
    }
    
    return todos, total, nil
}

// Transaction Helper
func WithTransaction(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
    tx := db.Begin()
    if tx.Error != nil {
        return tx.Error
    }
    
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)
        }
    }()
    
    if err := fn(tx); err != nil {
        tx.Rollback()
        return err
    }
    
    return tx.Commit().Error
}
```

### Redis for Caching

```go
// pkg/cache/redis.go
package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "github.com/go-redis/redis/v8"
)

type RedisCache struct {
    client *redis.Client
}

func NewRedisCache(cfg *config.RedisConfig) (*RedisCache, error) {
    client := redis.NewClient(&redis.Options{
        Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
        Password:     cfg.Password,
        DB:           cfg.DB,
        PoolSize:     cfg.PoolSize,
        MinIdleConns: 10,
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
    })
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("failed to connect redis: %w", err)
    }
    
    return &RedisCache{client: client}, nil
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

func (c *RedisCache) Delete(ctx context.Context, pattern string) error {
    keys, err := c.client.Keys(ctx, pattern).Result()
    if err != nil {
        return err
    }
    
    if len(keys) > 0 {
        return c.client.Del(ctx, keys...).Err()
    }
    
    return nil
}

// Cache Aside Pattern
func (uc *todoUseCase) GetTodo(ctx context.Context, id string) (*entity.Todo, error) {
    cacheKey := fmt.Sprintf("todo:%s", id)
    
    // Try cache first
    var todo entity.Todo
    err := uc.cache.Get(ctx, cacheKey, &todo)
    if err == nil {
        return &todo, nil
    }
    
    // Cache miss, get from database
    todo, err := uc.todoRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    if todo == nil {
        return nil, ErrNotFound
    }
    
    // Store in cache
    go uc.cache.Set(context.Background(), cacheKey, todo, 10*time.Minute)
    
    return todo, nil
}
```

---

## 5. Logging และ Monitoring

### Structured Logging with Zap

```go
// pkg/logger/logger.go
package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

type Logger interface {
    Debug(msg string, fields ...zap.Field)
    Info(msg string, fields ...zap.Field)
    Warn(msg string, fields ...zap.Field)
    Error(msg string, fields ...zap.Field)
    Fatal(msg string, fields ...zap.Field)
    With(fields ...zap.Field) Logger
    Sync() error
}

type ZapLogger struct {
    logger *zap.Logger
    sugar  *zap.SugaredLogger
}

func NewLogger(cfg *config.LogConfig) (Logger, error) {
    var level zapcore.Level
    switch cfg.Level {
    case "debug":
        level = zapcore.DebugLevel
    case "info":
        level = zapcore.InfoLevel
    case "warn":
        level = zapcore.WarnLevel
    case "error":
        level = zapcore.ErrorLevel
    default:
        level = zapcore.InfoLevel
    }
    
    encoderConfig := zapcore.EncoderConfig{
        TimeKey:        "timestamp",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "caller",
        MessageKey:     "message",
        StacktraceKey:  "stacktrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.LowercaseLevelEncoder,
        EncodeTime:     zapcore.ISO8601TimeEncoder,
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
    }
    
    var encoder zapcore.Encoder
    if cfg.Format == "json" {
        encoder = zapcore.NewJSONEncoder(encoderConfig)
    } else {
        encoder = zapcore.NewConsoleEncoder(encoderConfig)
    }
    
    core := zapcore.NewCore(
        encoder,
        zapcore.AddSync(os.Stdout),
        level,
    )
    
    logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
    
    return &ZapLogger{logger: logger, sugar: logger.Sugar()}, nil
}

func (l *ZapLogger) Debug(msg string, fields ...zap.Field) {
    l.logger.Debug(msg, fields...)
}

func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
    l.logger.Info(msg, fields...)
}

func (l *ZapLogger) Warn(msg string, fields ...zap.Field) {
    l.logger.Warn(msg, fields...)
}

func (l *ZapLogger) Error(msg string, fields ...zap.Field) {
    l.logger.Error(msg, fields...)
}

func (l *ZapLogger) Fatal(msg string, fields ...zap.Field) {
    l.logger.Fatal(msg, fields...)
}

func (l *ZapLogger) With(fields ...zap.Field) Logger {
    return &ZapLogger{logger: l.logger.With(fields...)}
}

func (l *ZapLogger) Sync() error {
    return l.logger.Sync()
}

// Usage with Context
type contextKey string

const loggerKey contextKey = "logger"

func WithLogger(ctx context.Context, logger Logger) context.Context {
    return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) Logger {
    if logger, ok := ctx.Value(loggerKey).(Logger); ok {
        return logger
    }
    return defaultLogger
}
```

### Prometheus Metrics

```go
// pkg/metrics/prometheus.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
    // HTTP metrics
    HTTPRequestsTotal    *prometheus.CounterVec
    HTTPRequestDuration  *prometheus.HistogramVec
    HTTPRequestsInFlight *prometheus.GaugeVec
    
    // Business metrics
    TodoCreatedTotal    prometheus.Counter
    TodoCompletedTotal  prometheus.Counter
    
    // Database metrics
    DatabaseQueriesTotal    *prometheus.CounterVec
    DatabaseQueryDuration   *prometheus.HistogramVec
    
    // Cache metrics
    CacheHitsTotal    prometheus.Counter
    CacheMissesTotal  prometheus.Counter
}

func NewMetrics() *Metrics {
    return &Metrics{
        HTTPRequestsTotal: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Name: "http_requests_total",
                Help: "Total number of HTTP requests",
            },
            []string{"method", "path", "status"},
        ),
        
        HTTPRequestDuration: promauto.NewHistogramVec(
            prometheus.HistogramOpts{
                Name:    "http_request_duration_seconds",
                Help:    "HTTP request duration in seconds",
                Buckets: prometheus.DefBuckets,
            },
            []string{"method", "path"},
        ),
        
        HTTPRequestsInFlight: promauto.NewGaugeVec(
            prometheus.GaugeOpts{
                Name: "http_requests_in_flight",
                Help: "Number of HTTP requests currently in flight",
            },
            []string{"method"},
        ),
        
        TodoCreatedTotal: promauto.NewCounter(
            prometheus.CounterOpts{
                Name: "todo_created_total",
                Help: "Total number of todos created",
            },
        ),
        
        TodoCompletedTotal: promauto.NewCounter(
            prometheus.CounterOpts{
                Name: "todo_completed_total",
                Help: "Total number of todos completed",
            },
        ),
        
        CacheHitsTotal: promauto.NewCounter(
            prometheus.CounterOpts{
                Name: "cache_hits_total",
                Help: "Total number of cache hits",
            },
        ),
        
        CacheMissesTotal: promauto.NewCounter(
            prometheus.CounterOpts{
                Name: "cache_misses_total",
                Help: "Total number of cache misses",
            },
        ),
    }
}

// Middleware for metrics
func (m *Metrics) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Track in-flight requests
        m.HTTPRequestsInFlight.WithLabelValues(r.Method).Inc()
        defer m.HTTPRequestsInFlight.WithLabelValues(r.Method).Dec()
        
        start := time.Now()
        
        // Create wrapped response writer to capture status code
        rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        next.ServeHTTP(rw, r)
        
        duration := time.Since(start).Seconds()
        
        m.HTTPRequestsTotal.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(rw.statusCode)).Inc()
        m.HTTPRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
    })
}
```

---

## 6. Middleware และ Authentication

### JWT Authentication

```go
// internal/middleware/auth.go
package middleware

import (
    "context"
    "net/http"
    "strings"
    "github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
    jwtSecret []byte
    logger    logger.Logger
}

func NewAuthMiddleware(cfg *config.JWTConfig, logger logger.Logger) *AuthMiddleware {
    return &AuthMiddleware{
        jwtSecret: []byte(cfg.Secret),
        logger:    logger,
    }
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            m.unauthorized(w, "missing authorization header")
            return
        }
        
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            m.unauthorized(w, "invalid authorization header format")
            return
        }
        
        tokenString := parts[1]
        claims, err := m.validateToken(tokenString)
        if err != nil {
            m.logger.Error("token validation failed", "error", err)
            m.unauthorized(w, "invalid or expired token")
            return
        }
        
        ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
        ctx = context.WithValue(ctx, "user_email", claims.Email)
        ctx = context.WithValue(ctx, "user_roles", claims.Roles)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func (m *AuthMiddleware) validateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return m.jwtSecret, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}

type Claims struct {
    UserID    string   `json:"user_id"`
    Email     string   `json:"email"`
    Roles     []string `json:"roles"`
    jwt.RegisteredClaims
}

// Role-based Authorization
func (m *AuthMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            userRoles, ok := r.Context().Value("user_roles").([]string)
            if !ok {
                m.forbidden(w, "user roles not found")
                return
            }
            
            if !hasAnyRole(userRoles, roles) {
                m.forbidden(w, "insufficient permissions")
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

// Request ID Middleware
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = generateRequestID()
        }
        
        w.Header().Set("X-Request-ID", requestID)
        ctx := context.WithValue(r.Context(), "request_id", requestID)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// CORS Middleware
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")
            
            // Check if origin is allowed
            allowed := false
            for _, o := range allowedOrigins {
                if o == "*" || o == origin {
                    allowed = true
                    break
                }
            }
            
            if allowed {
                w.Header().Set("Access-Control-Allow-Origin", origin)
                w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
                w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
                w.Header().Set("Access-Control-Allow-Credentials", "true")
            }
            
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

// Rate Limiting Middleware
type RateLimiter struct {
    store    *redis.Client
    requests int
    window   time.Duration
}

func NewRateLimiter(client *redis.Client, requests int, window time.Duration) *RateLimiter {
    return &RateLimiter{
        store:    client,
        requests: requests,
        window:   window,
    }
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        key := fmt.Sprintf("rate_limit:%s", r.RemoteAddr)
        
        ctx := context.Background()
        count, err := rl.store.Incr(ctx, key).Result()
        if err != nil {
            rl.logger.Error("rate limiter error", "error", err)
            next.ServeHTTP(w, r)
            return
        }
        
        if count == 1 {
            rl.store.Expire(ctx, key, rl.window)
        }
        
        if count > int64(rl.requests) {
            w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.requests))
            w.Header().Set("X-RateLimit-Remaining", "0")
            w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(rl.window).Unix(), 10))
            http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
            return
        }
        
        w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.requests))
        w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(rl.requests-int(count)))
        
        next.ServeHTTP(w, r)
    })
}
```

---

## 7. Testing Strategy

### Unit Tests with Testify

```go
// internal/usecase/todo_usecase_test.go
package usecase_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)

// Mock Repository
type MockTodoRepository struct {
    mock.Mock
}

func (m *MockTodoRepository) Create(ctx context.Context, todo *entity.Todo) error {
    args := m.Called(ctx, todo)
    return args.Error(0)
}

func (m *MockTodoRepository) FindByID(ctx context.Context, id string) (*entity.Todo, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.Todo), args.Error(1)
}

// Test Suite
type TodoUseCaseTestSuite struct {
    suite.Suite
    useCase    usecase.TodoUseCase
    mockRepo   *MockTodoRepository
    mockLogger *MockLogger
    mockCache  *MockCache
}

func (s *TodoUseCaseTestSuite) SetupTest() {
    s.mockRepo = new(MockTodoRepository)
    s.mockLogger = new(MockLogger)
    s.mockCache = new(MockCache)
    
    todoService := service.NewTodoService(s.mockRepo)
    s.useCase = usecase.NewTodoUseCase(todoService, s.mockLogger, s.mockCache)
}

func (s *TodoUseCaseTestSuite) TestCreateTodo_Success() {
    // Arrange
    ctx := context.Background()
    req := &usecase.CreateTodoRequest{
        Title:   "Test Todo",
        Content: "Test Content",
    }
    
    s.mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Todo")).
        Return(nil)
    
    s.mockCache.On("Delete", ctx, "todos:list").
        Return(nil)
    
    // Act
    resp, err := s.useCase.CreateTodo(ctx, req)
    
    // Assert
    assert.NoError(s.T(), err)
    assert.NotEmpty(s.T(), resp.ID)
    assert.Equal(s.T(), req.Title, resp.Title)
    assert.Equal(s.T(), "pending", resp.Status)
    
    s.mockRepo.AssertExpectations(s.T())
    s.mockCache.AssertExpectations(s.T())
}

func (s *TodoUseCaseTestSuite) TestCreateTodo_ValidationError() {
    // Arrange
    ctx := context.Background()
    req := &usecase.CreateTodoRequest{
        Title:   "", // Empty title
        Content: "Test Content",
    }
    
    // Act
    resp, err := s.useCase.CreateTodo(ctx, req)
    
    // Assert
    assert.Error(s.T(), err)
    assert.Nil(s.T(), resp)
    assert.Contains(s.T(), err.Error(), "title is required")
}

// Integration Tests
func TestTodoAPI_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // Setup test server
    server := setupTestServer(t, db)
    defer server.Close()
    
    // Test cases
    tests := []struct {
        name       string
        method     string
        path       string
        body       interface{}
        expected   int
    }{
        {
            name:     "Create Todo",
            method:   "POST",
            path:     "/todos",
            body:     map[string]string{"title": "Integration Test"},
            expected: http.StatusCreated,
        },
        {
            name:     "Get Todos",
            method:   "GET",
            path:     "/todos",
            expected: http.StatusOK,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var body io.Reader
            if tt.body != nil {
                jsonBody, _ := json.Marshal(tt.body)
                body = bytes.NewBuffer(jsonBody)
            }
            
            req, _ := http.NewRequest(tt.method, server.URL+tt.path, body)
            req.Header.Set("Content-Type", "application/json")
            
            client := &http.Client{}
            resp, err := client.Do(req)
            
            require.NoError(t, err)
            defer resp.Body.Close()
            
            assert.Equal(t, tt.expected, resp.StatusCode)
        })
    }
}

// Benchmark Tests
func BenchmarkCreateTodo(b *testing.B) {
    repo := &MockTodoRepository{}
    service := service.NewTodoService(repo)
    useCase := usecase.NewTodoUseCase(service, &MockLogger{}, &MockCache{})
    
    ctx := context.Background()
    req := &usecase.CreateTodoRequest{
        Title:   "Benchmark Todo",
        Content: "Benchmark Content",
    }
    
    repo.On("Create", ctx, mock.Anything).Return(nil)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = useCase.CreateTodo(ctx, req)
    }
}
```

### E2E Testing with Testcontainers

```go
// test/e2e/todo_test.go
package e2e_test

import (
    "context"
    "testing"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestE2E_TodoFlow(t *testing.T) {
    ctx := context.Background()
    
    // Start PostgreSQL container
    postgresContainer, err := postgres.Run(ctx,
        "postgres:15-alpine",
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("testuser"),
        postgres.WithPassword("testpass"),
    )
    require.NoError(t, err)
    defer postgresContainer.Terminate(ctx)
    
    // Get connection string
    connStr, err := postgresContainer.ConnectionString(ctx)
    require.NoError(t, err)
    
    // Start Redis container
    redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: testcontainers.ContainerRequest{
            Image:        "redis:7-alpine",
            ExposedPorts: []string{"6379/tcp"},
        },
        Started: true,
    })
    require.NoError(t, err)
    defer redisContainer.Terminate(ctx)
    
    redisHost, _ := redisContainer.Host(ctx)
    redisPort, _ := redisContainer.MappedPort(ctx, "6379")
    
    // Initialize application with test containers
    cfg := &config.Config{
        Database: config.DatabaseConfig{
            Driver:   "postgres",
            Host:     postgresContainer.Host,
            Port:     postgresContainer.MappedPort(ctx, "5432").Int(),
            Username: "testuser",
            Password: "testpass",
            Database: "testdb",
        },
        Redis: config.RedisConfig{
            Host: redisHost,
            Port: redisPort.Int(),
        },
    }
    
    app := setupApplication(t, cfg)
    defer app.Close()
    
    // E2E test flow
    t.Run("Create and Retrieve Todo", func(t *testing.T) {
        // Create todo
        createResp := createTodo(t, app.URL, "E2E Test Todo")
        assert.NotEmpty(t, createResp.ID)
        
        // Get todo
        todo := getTodo(t, app.URL, createResp.ID)
        assert.Equal(t, "E2E Test Todo", todo.Title)
        assert.Equal(t, "pending", todo.Status)
        
        // Update todo
        updateTodo(t, app.URL, createResp.ID, "completed")
        
        // Verify update
        updatedTodo := getTodo(t, app.URL, createResp.ID)
        assert.Equal(t, "completed", updatedTodo.Status)
        
        // Delete todo
        deleteTodo(t, app.URL, createResp.ID)
        
        // Verify deletion
        getTodoShouldFail(t, app.URL, createResp.ID)
    })
}
```

---

## 8. Performance Optimization

### Connection Pooling

```go
// Database connection pool optimization
func optimizeDatabasePool(db *sql.DB) {
    // Set maximum number of open connections
    db.SetMaxOpenConns(100)
    
    // Set maximum number of idle connections
    db.SetMaxIdleConns(50)
    
    // Set connection lifetime
    db.SetConnMaxLifetime(30 * time.Minute)
    
    // Set idle connection timeout
    db.SetConnMaxIdleTime(10 * time.Minute)
}

// HTTP Client Pooling
var httpClient = &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
        TLSHandshakeTimeout: 10 * time.Second,
        DisableCompression:  false,
    },
    Timeout: 30 * time.Second,
}
```

### Memory Optimization

```go
// Object Pooling with sync.Pool
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 4096)
    },
}

func processRequest(data []byte) {
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf)
    
    // Use buf...
}

// String Builder Optimization
func buildSQLQuery(conditions []string) string {
    var builder strings.Builder
    builder.Grow(256) // Pre-allocate capacity
    
    builder.WriteString("SELECT * FROM todos WHERE 1=1")
    for _, cond := range conditions {
        builder.WriteString(" AND ")
        builder.WriteString(cond)
    }
    
    return builder.String()
}

// Reduce allocations with pointer receivers
type LargeStruct struct {
    Data [1024]byte
}

// Value receiver creates copy
func (l LargeStruct) Process() { /* ... */ }

// Pointer receiver avoids copy
func (l *LargeStruct) Process() { /* ... */ }
```

### Concurrent Processing

```go
// Worker Pool Pattern
type WorkerPool struct {
    workers   int
    jobQueue  chan Job
    wg        sync.WaitGroup
    ctx       context.Context
    cancel    context.CancelFunc
}

func NewWorkerPool(workers, queueSize int) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    return &WorkerPool{
        workers:  workers,
        jobQueue: make(chan Job, queueSize),
        ctx:      ctx,
        cancel:   cancel,
    }
}

func (p *WorkerPool) Start() {
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go p.worker()
    }
}

func (p *WorkerPool) worker() {
    defer p.wg.Done()
    for {
        select {
        case job := <-p.jobQueue:
            job.Execute()
        case <-p.ctx.Done():
            return
        }
    }
}

// Batch Processing
func batchProcess(items []Item, batchSize int, processor func([]Item) error) error {
    for i := 0; i < len(items); i += batchSize {
        end := i + batchSize
        if end > len(items) {
            end = len(items)
        }
        
        batch := items[i:end]
        if err := processor(batch); err != nil {
            return err
        }
    }
    return nil
}

// Parallel Processing with errgroup
func processParallel(items []Item) error {
    g, ctx := errgroup.WithContext(context.Background())
    
    // Set concurrency limit
    sem := make(chan struct{}, 10)
    
    for _, item := range items {
        item := item // capture loop variable
        g.Go(func() error {
            sem <- struct{}{}
            defer func() { <-sem }()
            
            return processItem(ctx, item)
        })
    }
    
    return g.Wait()
}
```

---

## 9. Deployment และ Containerization

### Multi-stage Docker Build

```dockerfile
# deployments/docker/Dockerfile
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
    -ldflags="-w -s -X main.version=$(git describe --tags)" \
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

### Docker Compose for Development

```yaml
# deployments/docker/docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: todo_db
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
      APP_ENV: development
      DB_HOST: postgres
      DB_PORT: 5432
      REDIS_HOST: redis
      REDIS_PORT: 6379
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    volumes:
      - ../../configs:/app/configs
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - api

volumes:
  postgres_data:
  redis_data:
```

### Kubernetes Deployment

```yaml
# deployments/k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: todo-api
  namespace: production
spec:
  replicas: 3
  selector:
    matchLabels:
      app: todo-api
  template:
    metadata:
      labels:
        app: todo-api
    spec:
      containers:
      - name: api
        image: your-registry/todo-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: APP_ENV
          value: "production"
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: password
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
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
  name: todo-api-service
  namespace: production
spec:
  selector:
    app: todo-api
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: todo-api-hpa
  namespace: production
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: todo-api
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

---

## 10. Security Best Practices

### Input Validation

```go
// pkg/validator/validator.go
package validator

import (
    "github.com/go-playground/validator/v10"
)

type CustomValidator struct {
    validator *validator.Validate
}

func NewValidator() *CustomValidator {
    v := validator.New()
    
    // Register custom validation
    v.RegisterValidation("status", validateStatus)
    v.RegisterValidation("password", validatePassword)
    
    return &CustomValidator{validator: v}
}

func validateStatus(fl validator.FieldLevel) bool {
    status := fl.Field().String()
    validStatuses := []string{"pending", "in_progress", "completed", "cancelled"}
    
    for _, s := range validStatuses {
        if status == s {
            return true
        }
    }
    return false
}

func validatePassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    
    // Password must be at least 8 characters
    if len(password) < 8 {
        return false
    }
    
    // Must contain at least one number
    hasNumber := false
    // Must contain at least one uppercase
    hasUpper := false
    // Must contain at least one special character
    hasSpecial := false
    
    for _, char := range password {
        switch {
        case unicode.IsNumber(char):
            hasNumber = true
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }
    
    return hasNumber && hasUpper && hasSpecial
}

type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,password"`
    Name     string `json:"name" validate:"required,min=2,max=100"`
}

func (r *CreateUserRequest) Validate() error {
    return validator.New().Struct(r)
}
```

### SQL Injection Prevention

```go
// Always use parameterized queries
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
    var user entity.User
    
    // Good: Parameterized query
    err := r.db.WithContext(ctx).
        Where("email = ?", email).
        First(&user).Error
    
    // Bad: String concatenation (DO NOT DO THIS)
    // query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
    
    return &user, err
}

// For raw SQL queries
func (r *userRepository) ComplexQuery(ctx context.Context, userID string) error {
    // Good: Use placeholders
    query := `
        SELECT * FROM users 
        WHERE id = $1 
        AND status = $2
    `
    
    return r.db.WithContext(ctx).Raw(query, userID, "active").Scan(&users).Error
}
```

### Secret Management

```go
// pkg/secrets/vault.go
package secrets

import (
    "context"
    "github.com/hashicorp/vault/api"
)

type VaultSecretManager struct {
    client *api.Client
}

func NewVaultSecretManager(address, token string) (*VaultSecretManager, error) {
    config := &api.Config{
        Address: address,
    }
    
    client, err := api.NewClient(config)
    if err != nil {
        return nil, err
    }
    
    client.SetToken(token)
    
    return &VaultSecretManager{client: client}, nil
}

func (v *VaultSecretManager) GetSecret(ctx context.Context, path, key string) (string, error) {
    secret, err := v.client.Logical().ReadWithContext(ctx, path)
    if err != nil {
        return "", err
    }
    
    if secret == nil || secret.Data == nil {
        return "", errors.New("secret not found")
    }
    
    if value, ok := secret.Data[key].(string); ok {
        return value, nil
    }
    
    return "", errors.New("key not found in secret")
}

// Environment-based secrets for development
func getSecretFromEnv(key string) string {
    return os.Getenv(key)
}
```

### Security Headers Middleware

```go
// internal/middleware/security.go
package middleware

func SecurityHeadersMiddleware(next http.Handler) http.Handler {
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

### API Rate Limiting with Redis

```go
// pkg/ratelimit/sliding_window.go
package ratelimit

type SlidingWindowLimiter struct {
    client  *redis.Client
    limit   int
    window  time.Duration
}

func NewSlidingWindowLimiter(client *redis.Client, limit int, window time.Duration) *SlidingWindowLimiter {
    return &SlidingWindowLimiter{
        client: client,
        limit:  limit,
        window: window,
    }
}

func (l *SlidingWindowLimiter) Allow(ctx context.Context, key string) (bool, error) {
    now := time.Now().UnixMilli()
    windowStart := now - l.window.Milliseconds()
    
    script := `
        local key = KEYS[1]
        local now = tonumber(ARGV[1])
        local windowStart = tonumber(ARGV[2])
        local limit = tonumber(ARGV[3])
        
        -- Remove old entries
        redis.call('ZREMRANGEBYSCORE', key, 0, windowStart)
        
        -- Count current requests
        local current = redis.call('ZCARD', key)
        
        if current < limit then
            redis.call('ZADD', key, now, now)
            redis.call('EXPIRE', key, 60)
            return 1
        end
        
        return 0
    `
    
    result, err := l.client.Eval(ctx, script, []string{key}, now, windowStart, l.limit).Int()
    if err != nil {
        return false, err
    }
    
    return result == 1, nil
}
```

## สรุป

คู่มือนี้ครอบคลุมการพัฒนา Go สำหรับ production ตั้งแต่:
- **สถาปัตยกรรม**: Clean Architecture, DDD patterns
- **โครงสร้างโปรเจค**: Standard layout, Makefile automation
- **Configuration**: Multi-environment, Viper, env vars
- **Database**: Connection pooling, transactions, caching
- **Logging/Monitoring**: Structured logging, Prometheus metrics
- **Security**: JWT, rate limiting, input validation
- **Testing**: Unit, integration, E2E with testcontainers
- **Performance**: Pooling, concurrency, optimization
- **Deployment**: Docker, Kubernetes, CI/CD
- **Best Practices**: Error handling, code organization

ควรปรับใช้ตามความเหมาะสมของโปรเจค และทำการทดสอบอย่างละเอียดก่อนนำไปใช้งานจริง