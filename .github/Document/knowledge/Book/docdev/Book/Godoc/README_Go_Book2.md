# คู่มือ Golang สำหรับการใช้งานจริง เชิงลึก

## สารบัญ
1. [Project Structure & Architecture](#1-project-structure--architecture)
2. [Database & ORM](#2-database--orm)
3. [Web Framework & API Development](#3-web-framework--api-development)
4. [Middleware & Authentication](#4-middleware--authentication)
5. [Logging & Monitoring](#5-logging--monitoring)
6. [Configuration Management](#6-configuration-management)
7. [Testing Strategies](#7-testing-strategies)
8. [Docker & Deployment](#8-docker--deployment)
9. [Performance Optimization](#9-performance-optimization)
10. [Security Best Practices](#10-security-best-practices)
11. [Microservices Patterns](#11-microservices-patterns)
12. [Real-world Case Study](#12-real-world-case-study)

---

## 1. Project Structure & Architecture

### Clean Architecture Implementation

```
project/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point
├── internal/
│   ├── domain/                     # Enterprise business rules
│   │   ├── entity/
│   │   │   └── user.go
│   │   ├── repository/
│   │   │   └── user_repository.go  # Interface
│   │   └── service/
│   │       └── user_service.go     # Interface
│   ├── usecase/                    # Application business rules
│   │   └── user_usecase.go
│   ├── infrastructure/              # Frameworks, drivers
│   │   ├── database/
│   │   │   ├── mysql.go
│   │   │   ├── redis.go
│   │   │   └── mongodb.go
│   │   ├── cache/
│   │   │   └── redis_cache.go
│   │   ├── queue/
│   │   │   └── rabbitmq.go
│   │   └── http/
│   │       ├── router.go
│   │       └── middleware.go
│   ├── interface/                  # Adapters
│   │   ├── handler/
│   │   │   ├── user_handler.go
│   │   │   └── response.go
│   │   └── repository/
│   │       └── user_repository_impl.go
│   └── pkg/                        # Shared packages
│       ├── logger/
│       ├── errors/
│       ├── validator/
│       └── utils/
├── migrations/
│   └── *.sql
├── scripts/
│   └── build.sh
├── api/
│   └── openapi.yaml
├── configs/
│   ├── config.yaml
│   ├── config.dev.yaml
│   └── config.prod.yaml
├── deployments/
│   ├── docker-compose.yaml
│   └── kubernetes/
├── .env.example
├── go.mod
├── go.sum
└── Makefile
```

### Domain Layer Implementation

```go
// internal/domain/entity/user.go
package entity

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID        uuid.UUID `json:"id"`
    Email     string    `json:"email"`
    Password  string    `json:"-"` // Never expose password
    Name      string    `json:"name"`
    Role      Role      `json:"role"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type Role string

const (
    RoleUser  Role = "user"
    RoleAdmin Role = "admin"
)

// Domain errors
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrUserAlreadyExist = errors.New("user already exists")
    ErrInvalidPassword  = errors.New("invalid password")
)

// Business logic methods
func (u *User) IsAdmin() bool {
    return u.Role == RoleAdmin
}

func (u *User) Validate() error {
    if u.Email == "" {
        return errors.New("email is required")
    }
    if !isValidEmail(u.Email) {
        return errors.New("invalid email format")
    }
    if len(u.Password) < 8 {
        return errors.New("password must be at least 8 characters")
    }
    return nil
}

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
```

### Repository Interface

```go
// internal/domain/repository/user_repository.go
package repository

import (
    "context"
    "your-project/internal/domain/entity"
)

type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id uuid.UUID) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
    FindByEmail(ctx context.Context, email string) (*entity.User, error)
    FindAll(ctx context.Context, filter UserFilter) ([]*entity.User, int64, error)
    ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type UserFilter struct {
    Role   *entity.Role
    Search string
    Limit  int
    Offset int
    OrderBy string
}
```

### Use Case Implementation

```go
// internal/usecase/user_usecase.go
package usecase

import (
    "context"
    "your-project/internal/domain/entity"
    "your-project/internal/domain/repository"
    "github.com/google/uuid"
)

type UserUseCase interface {
    Register(ctx context.Context, req RegisterRequest) (*entity.User, error)
    Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
    GetProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error)
    UpdateProfile(ctx context.Context, userID uuid.UUID, req UpdateProfileRequest) (*entity.User, error)
    ListUsers(ctx context.Context, req ListUsersRequest) (*ListUsersResponse, error)
}

type userUseCase struct {
    userRepo repository.UserRepository
    cache    CacheService
    tokenGen TokenGenerator
    logger   Logger
}

func NewUserUseCase(
    userRepo repository.UserRepository,
    cache CacheService,
    tokenGen TokenGenerator,
    logger Logger,
) UserUseCase {
    return &userUseCase{
        userRepo: userRepo,
        cache:    cache,
        tokenGen: tokenGen,
        logger:   logger,
    }
}

func (uc *userUseCase) Register(ctx context.Context, req RegisterRequest) (*entity.User, error) {
    // Validate request
    if err := req.Validate(); err != nil {
        return nil, err
    }
    
    // Check if user exists
    exists, err := uc.userRepo.ExistsByEmail(ctx, req.Email)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, entity.ErrUserAlreadyExist
    }
    
    // Create user entity
    user := &entity.User{
        ID:       uuid.New(),
        Email:    req.Email,
        Name:     req.Name,
        Role:     entity.RoleUser,
    }
    
    // Hash password
    if err := user.HashPassword(); err != nil {
        return nil, err
    }
    
    // Save to database
    if err := uc.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    // Cache user (optional)
    uc.cache.Set(ctx, cacheKey(user.ID), user, 1*time.Hour)
    
    return user, nil
}

func (uc *userUseCase) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
    // Get user by email
    user, err := uc.userRepo.FindByEmail(ctx, req.Email)
    if err != nil {
        if errors.Is(err, entity.ErrUserNotFound) {
            return nil, entity.ErrInvalidPassword
        }
        return nil, err
    }
    
    // Check password
    if !user.CheckPassword(req.Password) {
        return nil, entity.ErrInvalidPassword
    }
    
    // Generate tokens
    accessToken, err := uc.tokenGen.GenerateAccessToken(user)
    if err != nil {
        return nil, err
    }
    
    refreshToken, err := uc.tokenGen.GenerateRefreshToken(user)
    if err != nil {
        return nil, err
    }
    
    return &AuthResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        User:         user,
    }, nil
}
```

---

## 2. Database & ORM

### GORM with Advanced Features

```go
// internal/infrastructure/database/gorm.go
package database

import (
    "fmt"
    "time"
    "gorm.io/gorm"
    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "gorm.io/plugin/soft_delete"
)

type Config struct {
    Driver          string
    Host            string
    Port            int
    Username        string
    Password        string
    Database        string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    LogMode         bool
}

func NewGormConnection(cfg Config) (*gorm.DB, error) {
    var dialector gorm.Dialector
    
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
    
    switch cfg.Driver {
    case "mysql":
        dialector = mysql.Open(dsn)
    case "postgres":
        dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
            cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port)
        dialector = postgres.Open(dsn)
    default:
        return nil, fmt.Errorf("unsupported driver: %s", cfg.Driver)
    }
    
    db, err := gorm.Open(dialector, &gorm.Config{
        SkipDefaultTransaction: true,
        PrepareStmt:            true,
        Logger:                 logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, err
    }
    
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
    sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
    sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
    
    return db, nil
}

// Advanced Model with Soft Delete
type BaseModel struct {
    ID        uuid.UUID           `gorm:"type:char(36);primary_key"`
    CreatedAt time.Time          `gorm:"index"`
    UpdatedAt time.Time          `gorm:"index"`
    DeletedAt soft_delete.DeletedAt `gorm:"index"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
    if b.ID == uuid.Nil {
        b.ID = uuid.New()
    }
    return nil
}

// User Model with Relations
type UserModel struct {
    BaseModel
    Email          string         `gorm:"uniqueIndex;size:255;not null"`
    Password       string         `gorm:"size:255;not null"`
    Name           string         `gorm:"size:255;not null"`
    Role           string         `gorm:"size:50;default:user"`
    Profile        *Profile       `gorm:"foreignKey:UserID"`
    Orders         []Order        `gorm:"foreignKey:UserID"`
    Sessions       []Session      `gorm:"foreignKey:UserID"`
}

type Profile struct {
    ID        uuid.UUID `gorm:"type:char(36);primary_key"`
    UserID    uuid.UUID `gorm:"type:char(36);uniqueIndex"`
    Avatar    string    `gorm:"size:500"`
    Bio       string    `gorm:"type:text"`
    Phone     string    `gorm:"size:20"`
    Address   string    `gorm:"type:text"`
}

// Repository Implementation with GORM
type userRepositoryImpl struct {
    db     *gorm.DB
    cache  CacheService
}

func NewUserRepository(db *gorm.DB, cache CacheService) repository.UserRepository {
    return &userRepositoryImpl{
        db:    db,
        cache: cache,
    }
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
    model := r.toModel(user)
    
    err := r.db.WithContext(ctx).
        Transaction(func(tx *gorm.DB) error {
            // Create user
            if err := tx.Create(model).Error; err != nil {
                return err
            }
            
            // Create profile if exists
            if model.Profile != nil {
                if err := tx.Create(model.Profile).Error; err != nil {
                    return err
                }
            }
            
            return nil
        })
    
    if err != nil {
        return err
    }
    
    return r.cache.Set(ctx, cacheKey(user.ID), user, 1*time.Hour)
}

func (r *userRepositoryImpl) FindAll(ctx context.Context, filter repository.UserFilter) ([]*entity.User, int64, error) {
    var models []*UserModel
    var total int64
    
    query := r.db.WithContext(ctx).Model(&UserModel{})
    
    // Apply filters
    if filter.Role != nil {
        query = query.Where("role = ?", *filter.Role)
    }
    
    if filter.Search != "" {
        query = query.Where("name LIKE ? OR email LIKE ?", 
            "%"+filter.Search+"%", 
            "%"+filter.Search+"%")
    }
    
    // Get total count
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Apply pagination
    if filter.Limit > 0 {
        query = query.Limit(filter.Limit)
    }
    if filter.Offset > 0 {
        query = query.Offset(filter.Offset)
    }
    
    // Apply ordering
    if filter.OrderBy != "" {
        query = query.Order(filter.OrderBy)
    } else {
        query = query.Order("created_at DESC")
    }
    
    // Preload relations
    query = query.Preload("Profile")
    
    if err := query.Find(&models).Error; err != nil {
        return nil, 0, err
    }
    
    users := make([]*entity.User, len(models))
    for i, model := range models {
        users[i] = r.toEntity(model)
    }
    
    return users, total, nil
}
```

### Database Migration with Goose

```go
// migrations/001_create_users_table.up.sql
CREATE TABLE IF NOT EXISTS users (
    id CHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_users_email (email),
    INDEX idx_users_role (role),
    INDEX idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- migrations/001_create_users_table.down.sql
DROP TABLE IF EXISTS users;
```

### Connection Pooling & Circuit Breaker

```go
// internal/infrastructure/database/pool.go
package database

import (
    "context"
    "database/sql"
    "github.com/sony/gobreaker"
)

type DatabasePool struct {
    db     *sql.DB
    cb     *gobreaker.CircuitBreaker
    config PoolConfig
}

type PoolConfig struct {
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    ConnMaxIdleTime time.Duration
}

func NewDatabasePool(driver, dsn string, config PoolConfig) (*DatabasePool, error) {
    db, err := sql.Open(driver, dsn)
    if err != nil {
        return nil, err
    }
    
    db.SetMaxOpenConns(config.MaxOpenConns)
    db.SetMaxIdleConns(config.MaxIdleConns)
    db.SetConnMaxLifetime(config.ConnMaxLifetime)
    db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
    
    // Circuit breaker settings
    cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
        Name:        "database",
        MaxRequests: 3,
        Interval:    10 * time.Second,
        Timeout:     30 * time.Second,
        ReadyToTrip: func(counts gobreaker.Counts) bool {
            failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
            return counts.Requests >= 10 && failureRatio >= 0.6
        },
    })
    
    return &DatabasePool{
        db: db,
        cb: cb,
        config: config,
    }, nil
}

func (p *DatabasePool) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    result, err := p.cb.Execute(func() (interface{}, error) {
        return p.db.QueryContext(ctx, query, args...)
    })
    
    if err != nil {
        return nil, err
    }
    
    return result.(*sql.Rows), nil
}
```

---

## 3. Web Framework & API Development

### Gin Framework Advanced Setup

```go
// internal/infrastructure/http/router.go
package http

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/gin-contrib/gzip"
    "github.com/gin-contrib/requestid"
    "github.com/ulule/limiter/v3"
    ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
    "github.com/ulule/limiter/v3/drivers/store/memory"
)

type Router struct {
    engine *gin.Engine
    config RouterConfig
}

type RouterConfig struct {
    Mode           string
    TrustedProxies []string
    RateLimit      RateLimitConfig
    CORS           CORSConfig
}

func NewRouter(config RouterConfig) *Router {
    // Set gin mode
    gin.SetMode(config.Mode)
    
    engine := gin.New()
    
    // Trust proxies
    if len(config.TrustedProxies) > 0 {
        engine.SetTrustedProxies(config.TrustedProxies)
    }
    
    // Middlewares
    engine.Use(gin.Recovery())
    engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
        Output:    os.Stdout,
        SkipPaths: []string{"/health"},
    }))
    engine.Use(requestid.New())
    engine.Use(gzip.Gzip(gzip.DefaultCompression))
    
    // CORS
    engine.Use(cors.New(cors.Config{
        AllowOrigins:     config.CORS.AllowOrigins,
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))
    
    // Rate Limiter
    rate, err := limiter.NewRateFromFormatted(config.RateLimit.Format)
    if err != nil {
        panic(err)
    }
    store := memory.NewStore()
    rateLimiter := ginlimiter.NewMiddleware(limiter.New(store, rate))
    engine.Use(rateLimiter)
    
    return &Router{
        engine: engine,
        config: config,
    }
}

func (r *Router) SetupRoutes(handlers *Handlers) {
    // Health check
    r.engine.GET("/health", handlers.HealthCheck)
    
    // API v1
    v1 := r.engine.Group("/api/v1")
    {
        // Public routes
        auth := v1.Group("/auth")
        {
            auth.POST("/register", handlers.User.Register)
            auth.POST("/login", handlers.User.Login)
            auth.POST("/refresh", handlers.User.RefreshToken)
        }
        
        // Protected routes
        protected := v1.Group("/")
        protected.Use(AuthMiddleware())
        {
            // User routes
            users := protected.Group("/users")
            {
                users.GET("/me", handlers.User.GetProfile)
                users.PUT("/me", handlers.User.UpdateProfile)
                users.DELETE("/me", handlers.User.DeleteAccount)
                
                // Admin only
                admin := users.Group("/")
                admin.Use(RoleMiddleware("admin"))
                {
                    admin.GET("/", handlers.User.ListUsers)
                    admin.GET("/:id", handlers.User.GetUser)
                    admin.PUT("/:id/role", handlers.User.UpdateUserRole)
                    admin.DELETE("/:id", handlers.User.DeleteUser)
                }
            }
            
            // Orders
            orders := protected.Group("/orders")
            {
                orders.POST("/", handlers.Order.CreateOrder)
                orders.GET("/", handlers.Order.ListOrders)
                orders.GET("/:id", handlers.Order.GetOrder)
                orders.PUT("/:id/cancel", handlers.Order.CancelOrder)
            }
        }
    }
}
```

### Handler Implementation with Validation

```go
// internal/interface/handler/user_handler.go
package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
)

type UserHandler struct {
    userUseCase usecase.UserUseCase
    validator   *validator.Validate
    logger      Logger
}

func NewUserHandler(userUseCase usecase.UserUseCase, logger Logger) *UserHandler {
    return &UserHandler{
        userUseCase: userUseCase,
        validator:   validator.New(),
        logger:      logger,
    }
}

// Request/Response DTOs
type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required,min=2,max=100"`
}

type RegisterResponse struct {
    ID    string `json:"id"`
    Email string `json:"email"`
    Name  string `json:"name"`
}

func (h *UserHandler) Register(c *gin.Context) {
    var req RegisterRequest
    
    // Bind JSON
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Code:    "INVALID_REQUEST",
            Message: err.Error(),
        })
        return
    }
    
    // Validate
    if err := h.validator.Struct(req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Code:    "VALIDATION_ERROR",
            Message: formatValidationError(err),
        })
        return
    }
    
    // Call use case
    user, err := h.userUseCase.Register(c.Request.Context(), usecase.RegisterRequest{
        Email:    req.Email,
        Password: req.Password,
        Name:     req.Name,
    })
    
    if err != nil {
        h.logger.Error("failed to register user", 
            zap.Error(err), 
            zap.String("email", req.Email))
        
        // Handle specific errors
        switch {
        case errors.Is(err, entity.ErrUserAlreadyExist):
            c.JSON(http.StatusConflict, ErrorResponse{
                Code:    "USER_EXISTS",
                Message: "User already exists",
            })
            return
        default:
            c.JSON(http.StatusInternalServerError, ErrorResponse{
                Code:    "INTERNAL_ERROR",
                Message: "Failed to register user",
            })
            return
        }
    }
    
    c.JSON(http.StatusCreated, RegisterResponse{
        ID:    user.ID.String(),
        Email: user.Email,
        Name:  user.Name,
    })
}

// Response wrapper
type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorDetail `json:"error,omitempty"`
    Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorDetail struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details interface{} `json:"details,omitempty"`
}

type Meta struct {
    Total   int64 `json:"total,omitempty"`
    Limit   int   `json:"limit,omitempty"`
    Offset  int   `json:"offset,omitempty"`
    Version string `json:"version,omitempty"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Success: true,
        Data:    data,
    })
}

func PaginatedResponse(c *gin.Context, data interface{}, total int64, limit, offset int) {
    c.JSON(http.StatusOK, Response{
        Success: true,
        Data:    data,
        Meta: &Meta{
            Total:  total,
            Limit:  limit,
            Offset: offset,
        },
    })
}
```

---

## 4. Middleware & Authentication

### JWT Authentication Middleware

```go
// internal/infrastructure/http/middleware/auth.go
package middleware

import (
    "context"
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
    tokenService TokenService
    userRepo     repository.UserRepository
    logger       Logger
}

func NewAuthMiddleware(tokenService TokenService, userRepo repository.UserRepository, logger Logger) *AuthMiddleware {
    return &AuthMiddleware{
        tokenService: tokenService,
        userRepo:     userRepo,
        logger:       logger,
    }
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract token
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, ErrorResponse{
                Code:    "MISSING_TOKEN",
                Message: "Authorization header is required",
            })
            c.Abort()
            return
        }
        
        // Parse token
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, ErrorResponse{
                Code:    "INVALID_TOKEN_FORMAT",
                Message: "Authorization header must be Bearer token",
            })
            c.Abort()
            return
        }
        
        token := parts[1]
        
        // Validate token
        claims, err := m.tokenService.ValidateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, ErrorResponse{
                Code:    "INVALID_TOKEN",
                Message: err.Error(),
            })
            c.Abort()
            return
        }
        
        // Get user from database
        userID, err := uuid.Parse(claims.UserID)
        if err != nil {
            c.JSON(http.StatusUnauthorized, ErrorResponse{
                Code:    "INVALID_USER",
                Message: "Invalid user ID in token",
            })
            c.Abort()
            return
        }
        
        user, err := m.userRepo.FindByID(c.Request.Context(), userID)
        if err != nil {
            m.logger.Error("failed to find user", zap.Error(err), zap.String("user_id", claims.UserID))
            c.JSON(http.StatusUnauthorized, ErrorResponse{
                Code:    "USER_NOT_FOUND",
                Message: "User not found",
            })
            c.Abort()
            return
        }
        
        // Set user in context
        c.Set("user", user)
        c.Set("user_id", user.ID)
        c.Set("claims", claims)
        
        c.Next()
    }
}

func (m *AuthMiddleware) Authorize(roles ...entity.Role) gin.HandlerFunc {
    return func(c *gin.Context) {
        user, exists := c.Get("user")
        if !exists {
            c.JSON(http.StatusUnauthorized, ErrorResponse{
                Code:    "UNAUTHORIZED",
                Message: "User not authenticated",
            })
            c.Abort()
            return
        }
        
        currentUser := user.(*entity.User)
        
        // Check if user has required role
        allowed := false
        for _, role := range roles {
            if currentUser.Role == role {
                allowed = true
                break
            }
        }
        
        if !allowed {
            c.JSON(http.StatusForbidden, ErrorResponse{
                Code:    "FORBIDDEN",
                Message: "Insufficient permissions",
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// Token Service Implementation
type TokenService interface {
    GenerateAccessToken(user *entity.User) (string, error)
    GenerateRefreshToken(user *entity.User) (string, error)
    ValidateToken(token string) (*Claims, error)
}

type JWTTokenService struct {
    secretKey     []byte
    accessExpiry  time.Duration
    refreshExpiry time.Duration
    issuer        string
}

type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func NewJWTTokenService(config TokenConfig) *JWTTokenService {
    return &JWTTokenService{
        secretKey:     []byte(config.SecretKey),
        accessExpiry:  config.AccessExpiry,
        refreshExpiry: config.RefreshExpiry,
        issuer:        config.Issuer,
    }
}

func (s *JWTTokenService) GenerateAccessToken(user *entity.User) (string, error) {
    claims := &Claims{
        UserID: user.ID.String(),
        Email:  user.Email,
        Role:   string(user.Role),
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessExpiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    s.issuer,
            Subject:   user.ID.String(),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.secretKey)
}

func (s *JWTTokenService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return s.secretKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}
```

### Rate Limiting Middleware

```go
// internal/infrastructure/http/middleware/rate_limit.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis_rate/v9"
    "github.com/redis/go-redis/v9"
    "golang.org/x/time/rate"
)

type RateLimiter struct {
    limiter *redis_rate.Limiter
    config  RateLimitConfig
}

type RateLimitConfig struct {
    RequestsPerSecond int
    Burst             int
    EnableRedis       bool
    RedisURL          string
}

func NewRateLimiter(config RateLimitConfig) (*RateLimiter, error) {
    if config.EnableRedis {
        rdb := redis.NewClient(&redis.Options{
            Addr: config.RedisURL,
        })
        limiter := redis_rate.NewLimiter(rdb)
        return &RateLimiter{
            limiter: limiter,
            config:  config,
        }, nil
    }
    
    return &RateLimiter{
        config: config,
    }, nil
}

func (r *RateLimiter) Limit(limitKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        if r.limiter != nil {
            // Redis-based rate limiting
            res, err := r.limiter.Allow(c.Request.Context(), limitKey, redis_rate.PerSecond(r.config.RequestsPerSecond))
            if err != nil {
                c.Next()
                return
            }
            
            c.Header("X-RateLimit-Limit", r.config.RequestsPerSecond)
            c.Header("X-RateLimit-Remaining", res.Remaining)
            c.Header("X-RateLimit-Reset", res.ResetAfter.String())
            
            if res.Allowed == 0 {
                c.JSON(http.StatusTooManyRequests, ErrorResponse{
                    Code:    "RATE_LIMIT_EXCEEDED",
                    Message: "Rate limit exceeded",
                })
                c.Abort()
                return
            }
        } else {
            // In-memory rate limiting
            limiter := rate.NewLimiter(rate.Limit(r.config.RequestsPerSecond), r.config.Burst)
            if !limiter.Allow() {
                c.JSON(http.StatusTooManyRequests, ErrorResponse{
                    Code:    "RATE_LIMIT_EXCEEDED",
                    Message: "Rate limit exceeded",
                })
                c.Abort()
                return
            }
        }
        
        c.Next()
    }
}
```

---

## 5. Logging & Monitoring

### Structured Logging with Zap

```go
// internal/pkg/logger/logger.go
package logger

import (
    "os"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
    Level      string
    Format     string // json or console
    OutputPath string
    ErrorPath  string
    MaxSize    int // MB
    MaxBackups int
    MaxAge     int // days
    Compress   bool
}

type Logger interface {
    Debug(msg string, fields ...zap.Field)
    Info(msg string, fields ...zap.Field)
    Warn(msg string, fields ...zap.Field)
    Error(msg string, fields ...zap.Field)
    Fatal(msg string, fields ...zap.Field)
    With(fields ...zap.Field) Logger
    Sync() error
}

type zapLogger struct {
    logger *zap.Logger
}

func NewLogger(config Config) (Logger, error) {
    var level zapcore.Level
    if err := level.UnmarshalText([]byte(config.Level)); err != nil {
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
    if config.Format == "json" {
        encoder = zapcore.NewJSONEncoder(encoderConfig)
    } else {
        encoder = zapcore.NewConsoleEncoder(encoderConfig)
    }
    
    // Log rotation
    writer := &lumberjack.Logger{
        Filename:   config.OutputPath,
        MaxSize:    config.MaxSize,
        MaxBackups: config.MaxBackups,
        MaxAge:     config.MaxAge,
        Compress:   config.Compress,
    }
    
    // Error writer
    errorWriter := &lumberjack.Logger{
        Filename:   config.ErrorPath,
        MaxSize:    config.MaxSize,
        MaxBackups: config.MaxBackups,
        MaxAge:     config.MaxAge,
        Compress:   config.Compress,
    }
    
    core := zapcore.NewTee(
        zapcore.NewCore(encoder, zapcore.AddSync(writer), level),
        zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), zapcore.ErrorLevel),
        zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
    )
    
    logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
    
    return &zapLogger{logger: logger}, nil
}

func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
    l.logger.Debug(msg, fields...)
}

func (l *zapLogger) Info(msg string, fields ...zap.Field) {
    l.logger.Info(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...zap.Field) {
    l.logger.Warn(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...zap.Field) {
    l.logger.Error(msg, fields...)
}

func (l *zapLogger) Fatal(msg string, fields ...zap.Field) {
    l.logger.Fatal(msg, fields...)
}

func (l *zapLogger) With(fields ...zap.Field) Logger {
    return &zapLogger{logger: l.logger.With(fields...)}
}

func (l *zapLogger) Sync() error {
    return l.logger.Sync()
}

// Request logging middleware
func RequestLogger(logger Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        query := c.Request.URL.RawQuery
        
        c.Next()
        
        latency := time.Since(start)
        status := c.Writer.Status()
        
        fields := []zap.Field{
            zap.String("method", c.Request.Method),
            zap.String("path", path),
            zap.String("query", query),
            zap.Int("status", status),
            zap.Duration("latency", latency),
            zap.String("ip", c.ClientIP()),
            zap.String("user_agent", c.Request.UserAgent()),
            zap.String("request_id", c.GetString("request_id")),
        }
        
        if userID, exists := c.Get("user_id"); exists {
            fields = append(fields, zap.String("user_id", userID.(string)))
        }
        
        if status >= 500 {
            logger.Error("request failed", fields...)
        } else if status >= 400 {
            logger.Warn("request client error", fields...)
        } else {
            logger.Info("request completed", fields...)
        }
    }
}
```

### Metrics with Prometheus

```go
// internal/infrastructure/metrics/prometheus.go
package metrics

import (
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    // HTTP metrics
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request latency",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
    
    // Business metrics
    activeUsers = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_users",
            Help: "Number of active users",
        },
    )
    
    ordersTotal = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "orders_total",
            Help: "Total number of orders",
        },
    )
    
    // Database metrics
    dbQueryDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "db_query_duration_seconds",
            Help:    "Database query duration",
            Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
        },
        []string{"query", "table"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
    prometheus.MustRegister(activeUsers)
    prometheus.MustRegister(ordersTotal)
    prometheus.MustRegister(dbQueryDuration)
}

func PrometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start).Seconds()
        status := c.Writer.Status()
        method := c.Request.Method
        endpoint := c.FullPath()
        
        httpRequestsTotal.WithLabelValues(method, endpoint, strconv.Itoa(status)).Inc()
        httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
    }
}

func MetricsHandler() gin.HandlerFunc {
    h := promhttp.Handler()
    return func(c *gin.Context) {
        h.ServeHTTP(c.Writer, c.Request)
    }
}
```

### Distributed Tracing with OpenTelemetry

```go
// internal/infrastructure/tracing/otel.go
package tracing

import (
    "context"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
    "go.opentelemetry.io/otel/trace"
)

type TracingConfig struct {
    ServiceName    string
    ServiceVersion string
    Environment    string
    JaegerEndpoint string
    SampleRate     float64
}

func InitTracing(config TracingConfig) (trace.TracerProvider, error) {
    // Create Jaeger exporter
    exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.JaegerEndpoint)))
    if err != nil {
        return nil, err
    }
    
    // Create resource
    res, err := resource.New(context.Background(),
        resource.WithAttributes(
            semconv.ServiceNameKey.String(config.ServiceName),
            semconv.ServiceVersionKey.String(config.ServiceVersion),
            semconv.DeploymentEnvironmentKey.String(config.Environment),
        ),
    )
    if err != nil {
        return nil, err
    }
    
    // Create trace provider
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(res),
        sdktrace.WithSampler(sdktrace.TraceIDRatioBased(config.SampleRate)),
    )
    
    otel.SetTracerProvider(tp)
    otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{},
        propagation.Baggage{},
    ))
    
    return tp, nil
}

// Tracing middleware for Gin
func TracingMiddleware(serviceName string) gin.HandlerFunc {
    tracer := otel.Tracer(serviceName)
    
    return func(c *gin.Context) {
        // Extract context from headers
        ctx := otel.GetTextMapPropagator().Extract(
            c.Request.Context(),
            propagation.HeaderCarrier(c.Request.Header),
        )
        
        // Start span
        ctx, span := tracer.Start(ctx, c.Request.URL.Path,
            trace.WithSpanKind(trace.SpanKindServer),
        )
        defer span.End()
        
        // Set attributes
        span.SetAttributes(
            semconv.HTTPMethodKey.String(c.Request.Method),
            semconv.HTTPURLKey.String(c.Request.URL.String()),
            semconv.HTTPTargetKey.String(c.Request.URL.Path),
            semconv.HTTPUserAgentKey.String(c.Request.UserAgent()),
        )
        
        // Add to context
        c.Request = c.Request.WithContext(ctx)
        
        c.Next()
        
        // Set status
        span.SetAttributes(semconv.HTTPStatusCodeKey.Int(c.Writer.Status()))
        if c.Writer.Status() >= 500 {
            span.SetStatus(codes.Error, "server error")
        }
    }
}
```

---

## 6. Configuration Management

### Configuration with Viper

```go
// internal/infrastructure/config/config.go
package config

import (
    "fmt"
    "github.com/spf13/viper"
    "github.com/joho/godotenv"
)

type Config struct {
    App      AppConfig      `mapstructure:"app"`
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    Log      LogConfig      `mapstructure:"log"`
    Metrics  MetricsConfig  `mapstructure:"metrics"`
}

type AppConfig struct {
    Name    string `mapstructure:"name"`
    Version string `mapstructure:"version"`
    Env     string `mapstructure:"env"`
    Debug   bool   `mapstructure:"debug"`
}

type ServerConfig struct {
    Port         int           `mapstructure:"port"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
    IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type DatabaseConfig struct {
    Driver          string        `mapstructure:"driver"`
    Host            string        `mapstructure:"host"`
    Port            int           `mapstructure:"port"`
    Username        string        `mapstructure:"username"`
    Password        string        `mapstructure:"password"`
    Database        string        `mapstructure:"database"`
    MaxOpenConns    int           `mapstructure:"max_open_conns"`
    MaxIdleConns    int           `mapstructure:"max_idle_conns"`
    ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type JWTConfig struct {
    Secret        string        `mapstructure:"secret"`
    AccessExpiry  time.Duration `mapstructure:"access_expiry"`
    RefreshExpiry time.Duration `mapstructure:"refresh_expiry"`
    Issuer        string        `mapstructure:"issuer"`
}

func LoadConfig(configPath string) (*Config, error) {
    // Load .env file if exists
    _ = godotenv.Load()
    
    v := viper.New()
    
    // Set defaults
    v.SetDefault("app.env", "development")
    v.SetDefault("server.port", 8080)
    v.SetDefault("server.read_timeout", "30s")
    v.SetDefault("server.write_timeout", "30s")
    
    // Config file
    v.SetConfigFile(configPath)
    v.SetConfigType("yaml")
    
    // Environment variables
    v.AutomaticEnv()
    v.SetEnvPrefix("APP")
    
    // Replace dots with underscores for env vars
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    
    // Read config file
    if err := v.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    
    var config Config
    if err := v.Unmarshal(&config); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }
    
    // Validate config
    if err := config.Validate(); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }
    
    return &config, nil
}

func (c *Config) Validate() error {
    if c.App.Name == "" {
        return errors.New("app name is required")
    }
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return errors.New("invalid server port")
    }
    if c.Database.Driver == "" {
        return errors.New("database driver is required")
    }
    return nil
}

// DSN for database connection
func (c *DatabaseConfig) DSN() string {
    switch c.Driver {
    case "mysql":
        return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            c.Username, c.Password, c.Host, c.Port, c.Database)
    case "postgres":
        return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
            c.Host, c.Port, c.Username, c.Password, c.Database)
    default:
        return ""
    }
}
```

### Configuration Files

```yaml
# configs/config.yaml
app:
  name: "myapp"
  version: "1.0.0"
  env: "development"
  debug: true

server:
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "120s"

database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  username: "root"
  password: "${DB_PASSWORD}"
  database: "myapp"
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
  secret: "${JWT_SECRET}"
  access_expiry: "15m"
  refresh_expiry: "7d"
  issuer: "myapp"

log:
  level: "info"
  format: "json"
  output_path: "/var/log/myapp/app.log"
  error_path: "/var/log/myapp/error.log"
  max_size: 100
  max_backups: 3
  max_age: 28
  compress: true

metrics:
  enabled: true
  path: "/metrics"
  port: 9090
```

---

## 7. Testing Strategies

### Unit Testing with Testify

```go
// internal/usecase/user_usecase_test.go
package usecase_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/suite"
)

type UserUseCaseSuite struct {
    suite.Suite
    useCase   usecase.UserUseCase
    mockRepo  *MockUserRepository
    mockCache *MockCacheService
    mockToken *MockTokenGenerator
}

func (s *UserUseCaseSuite) SetupTest() {
    s.mockRepo = new(MockUserRepository)
    s.mockCache = new(MockCacheService)
    s.mockToken = new(MockTokenGenerator)
    
    s.useCase = usecase.NewUserUseCase(
        s.mockRepo,
        s.mockCache,
        s.mockToken,
        logger.NewNoOpLogger(),
    )
}

func (s *UserUseCaseSuite) TestRegister_Success() {
    // Arrange
    req := usecase.RegisterRequest{
        Email:    "test@example.com",
        Password: "password123",
        Name:     "Test User",
    }
    
    s.mockRepo.On("ExistsByEmail", mock.Anything, req.Email).Return(false, nil)
    s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
    s.mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
    
    // Act
    user, err := s.useCase.Register(context.Background(), req)
    
    // Assert
    assert.NoError(s.T(), err)
    assert.NotNil(s.T(), user)
    assert.Equal(s.T(), req.Email, user.Email)
    assert.Equal(s.T(), req.Name, user.Name)
    assert.NotEmpty(s.T(), user.Password) // Password should be hashed
    assert.NotEqual(s.T(), req.Password, user.Password)
    
    s.mockRepo.AssertExpectations(s.T())
}

func (s *UserUseCaseSuite) TestRegister_UserAlreadyExists() {
    // Arrange
    req := usecase.RegisterRequest{
        Email:    "existing@example.com",
        Password: "password123",
        Name:     "Existing User",
    }
    
    s.mockRepo.On("ExistsByEmail", mock.Anything, req.Email).Return(true, nil)
    
    // Act
    user, err := s.useCase.Register(context.Background(), req)
    
    // Assert
    assert.Error(s.T(), err)
    assert.Nil(s.T(), user)
    assert.Equal(s.T(), entity.ErrUserAlreadyExist, err)
    
    s.mockRepo.AssertExpectations(s.T())
    s.mockRepo.AssertNotCalled(s.T(), "Create")
}

func (s *UserUseCaseSuite) TestLogin_Success() {
    // Arrange
    req := usecase.LoginRequest{
        Email:    "test@example.com",
        Password: "password123",
    }
    
    user := &entity.User{
        ID:       uuid.New(),
        Email:    req.Email,
        Password: "$2a$10$hashedpassword", // Hashed version of "password123"
        Name:     "Test User",
        Role:     entity.RoleUser,
    }
    
    s.mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(user, nil)
    s.mockToken.On("GenerateAccessToken", user).Return("access-token", nil)
    s.mockToken.On("GenerateRefreshToken", user).Return("refresh-token", nil)
    
    // Mock password verification
    user.Password = hashPassword("password123") // Helper to hash
    
    // Act
    resp, err := s.useCase.Login(context.Background(), req)
    
    // Assert
    assert.NoError(s.T(), err)
    assert.NotNil(s.T(), resp)
    assert.Equal(s.T(), "access-token", resp.AccessToken)
    assert.Equal(s.T(), "refresh-token", resp.RefreshToken)
    assert.Equal(s.T(), user, resp.User)
}

// Mock implementations
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

func TestUserUseCase(t *testing.T) {
    suite.Run(t, new(UserUseCaseSuite))
}
```

### Integration Testing

```go
// tests/integration/user_api_test.go
package integration_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestUserAPI_Register(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // Setup router
    router := setupRouter(db)
    
    // Test cases
    tests := []struct {
        name       string
        request    RegisterRequest
        expectedStatus int
        expectedError string
    }{
        {
            name: "successful registration",
            request: RegisterRequest{
                Email:    "test@example.com",
                Password: "password123",
                Name:     "Test User",
            },
            expectedStatus: http.StatusCreated,
        },
        {
            name: "duplicate email",
            request: RegisterRequest{
                Email:    "test@example.com",
                Password: "password123",
                Name:     "Another User",
            },
            expectedStatus: http.StatusConflict,
            expectedError: "User already exists",
        },
        {
            name: "invalid email",
            request: RegisterRequest{
                Email:    "invalid-email",
                Password: "password123",
                Name:     "Test User",
            },
            expectedStatus: http.StatusBadRequest,
            expectedError: "invalid email format",
        },
        {
            name: "short password",
            request: RegisterRequest{
                Email:    "test2@example.com",
                Password: "short",
                Name:     "Test User",
            },
            expectedStatus: http.StatusBadRequest,
            expectedError: "password must be at least 8 characters",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Prepare request
            body, _ := json.Marshal(tt.request)
            req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            
            // Execute request
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            
            // Assert
            assert.Equal(t, tt.expectedStatus, w.Code)
            
            if tt.expectedError != "" {
                var response ErrorResponse
                err := json.Unmarshal(w.Body.Bytes(), &response)
                require.NoError(t, err)
                assert.Contains(t, response.Message, tt.expectedError)
            }
        })
    }
}

func setupTestDB(t *testing.T) *gorm.DB {
    // Create test database
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    require.NoError(t, err)
    
    // Run migrations
    err = db.AutoMigrate(&UserModel{}, &ProfileModel{})
    require.NoError(t, err)
    
    return db
}

func setupRouter(db *gorm.DB) *gin.Engine {
    // Setup repositories
    userRepo := repository.NewUserRepository(db)
    
    // Setup use cases
    userUseCase := usecase.NewUserUseCase(userRepo, nil, nil, nil)
    
    // Setup handlers
    userHandler := handler.NewUserHandler(userUseCase)
    
    // Setup router
    router := gin.New()
    api := router.Group("/api/v1")
    {
        auth := api.Group("/auth")
        {
            auth.POST("/register", userHandler.Register)
            auth.POST("/login", userHandler.Login)
        }
    }
    
    return router
}
```

### Benchmark Testing

```go
// internal/usecase/user_usecase_bench_test.go
package usecase_test

import (
    "context"
    "testing"
)

func BenchmarkUserUseCase_Register(b *testing.B) {
    // Setup
    mockRepo := new(MockUserRepository)
    mockCache := new(MockCacheService)
    mockToken := new(MockTokenGenerator)
    
    useCase := usecase.NewUserUseCase(mockRepo, mockCache, mockToken, nil)
    
    // Setup mock expectations
    mockRepo.On("ExistsByEmail", mock.Anything, mock.Anything).Return(false, nil)
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
    
    req := usecase.RegisterRequest{
        Email:    "bench@example.com",
        Password: "password123",
        Name:     "Benchmark User",
    }
    
    // Reset timer
    b.ResetTimer()
    
    // Run benchmark
    for i := 0; i < b.N; i++ {
        _, _ = useCase.Register(context.Background(), req)
    }
}

func BenchmarkDatabaseQuery(b *testing.B) {
    db := setupTestDB(b)
    defer cleanupTestDB(b, db)
    
    // Insert test data
    for i := 0; i < 1000; i++ {
        user := &UserModel{
            Email:    fmt.Sprintf("user%d@example.com", i),
            Password: "hashed",
            Name:     fmt.Sprintf("User %d", i),
        }
        db.Create(user)
    }
    
    repo := repository.NewUserRepository(db)
    
    b.ResetTimer()
    
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            _, _, _ = repo.FindAll(context.Background(), repository.UserFilter{
                Limit:  100,
                Offset: 0,
            })
        }
    })
}
```

---

## 8. Docker & Deployment

### Multi-stage Dockerfile

```dockerfile
# Dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=$(git describe --tags)" \
    -o /app/bin/api ./cmd/api

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create app user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy binary from builder
COPY --from=builder --chown=appuser:appgroup /app/bin/api .
COPY --chown=appuser:appgroup configs/config.yaml ./configs/

# Copy migrations
COPY --chown=appuser:appgroup migrations ./migrations

# Expose ports
EXPOSE 8080 9090

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

USER appuser

ENTRYPOINT ["./api"]
```

### Docker Compose for Development

```yaml
# docker-compose.yaml
version: '3.8'

services:
  api:
    build:
      context: .
      dockerfile: Dockerfile.dev
    ports:
      - "8080:8080"
      - "9090:9090"
    environment:
      - APP_ENV=development
      - DB_HOST=mysql
      - DB_PASSWORD=secret
      - REDIS_HOST=redis
      - JWT_SECRET=dev-secret-key
    volumes:
      - ./:/app
      - /app/tmp
    depends_on:
      - mysql
      - redis
      - jaeger
    networks:
      - app-network

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: secret
      MYSQL_DATABASE: myapp
      MYSQL_USER: appuser
      MYSQL_PASSWORD: apppass
    ports:
      - "3306:3306"
    volumes:
      - mysql-data:/var/lib/mysql
    networks:
      - app-network

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    networks:
      - app-network

  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686"
      - "6831:6831/udp"
    networks:
      - app-network

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9091:9090"
    volumes:
      - ./deployments/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus-data:/prometheus
    networks:
      - app-network

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-data:/var/lib/grafana
    networks:
      - app-network

volumes:
  mysql-data:
  redis-data:
  prometheus-data:
  grafana-data:

networks:
  app-network:
    driver: bridge
```

### Kubernetes Deployment

```yaml
# deployments/kubernetes/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp-api
  namespace: production
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp-api
  template:
    metadata:
      labels:
        app: myapp-api
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "9090"
        prometheus.io/path: "/metrics"
    spec:
      containers:
      - name: api
        image: myregistry/myapp:latest
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 9090
          name: metrics
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
  name: myapp-api
  namespace: production
spec:
  selector:
    app: myapp-api
  ports:
  - port: 80
    targetPort: 8080
    name: http
  - port: 9090
    targetPort: 9090
    name: metrics
  type: ClusterIP
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: myapp-api-hpa
  namespace: production
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: myapp-api
  minReplicas: 3
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

## 9. Performance Optimization

### Connection Pooling & Caching

```go
// internal/infrastructure/cache/redis_cache.go
package cache

import (
    "context"
    "encoding/json"
    "time"
    "github.com/redis/go-redis/v9"
)

type RedisCache struct {
    client *redis.Client
    config RedisConfig
}

type RedisConfig struct {
    Host         string
    Port         int
    Password     string
    DB           int
    PoolSize     int
    MinIdleConns int
    MaxRetries   int
}

func NewRedisCache(config RedisConfig) (*RedisCache, error) {
    client := redis.NewClient(&redis.Options{
        Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
        Password:     config.Password,
        DB:           config.DB,
        PoolSize:     config.PoolSize,
        MinIdleConns: config.MinIdleConns,
        MaxRetries:   config.MaxRetries,
    })
    
    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, err
    }
    
    return &RedisCache{
        client: client,
        config: config,
    }, nil
}

func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
    data, err := c.client.Get(ctx, key).Bytes()
    if err != nil {
        return err
    }
    
    return json.Unmarshal(data, dest)
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    
    return c.client.Set(ctx, key, data, ttl).Err()
}

// Cache aside pattern
func (uc *userUseCase) GetUserWithCache(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    cacheKey := fmt.Sprintf("user:%s", id)
    
    // Try cache first
    var user entity.User
    err := uc.cache.Get(ctx, cacheKey, &user)
    if err == nil {
        return &user, nil
    }
    
    // Cache miss - get from DB
    user, err := uc.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    _ = uc.cache.Set(ctx, cacheKey, user, 1*time.Hour)
    
    return user, nil
}
```

### Query Optimization

```go
// Batch processing
func (r *userRepositoryImpl) CreateBatch(ctx context.Context, users []*entity.User) error {
    batchSize := 1000
    for i := 0; i < len(users); i += batchSize {
        end := i + batchSize
        if end > len(users) {
            end = len(users)
        }
        
        batch := users[i:end]
        models := make([]*UserModel, len(batch))
        for j, user := range batch {
            models[j] = r.toModel(user)
        }
        
        if err := r.db.WithContext(ctx).CreateInBatches(models, batchSize).Error; err != nil {
            return err
        }
    }
    return nil
}

// Efficient pagination using cursor
func (r *userRepositoryImpl) FindWithCursor(ctx context.Context, cursor string, limit int) ([]*entity.User, string, error) {
    var models []*UserModel
    
    query := r.db.WithContext(ctx).Model(&UserModel{})
    
    if cursor != "" {
        // Decode cursor (e.g., base64 encoded timestamp+id)
        lastID, lastCreatedAt, err := decodeCursor(cursor)
        if err == nil {
            query = query.Where("(created_at, id) > (?, ?)", lastCreatedAt, lastID)
        }
    }
    
    query = query.Order("created_at ASC, id ASC").Limit(limit)
    
    if err := query.Find(&models).Error; err != nil {
        return nil, "", err
    }
    
    users := make([]*entity.User, len(models))
    for i, model := range models {
        users[i] = r.toEntity(model)
    }
    
    // Generate next cursor
    var nextCursor string
    if len(models) == limit {
        last := models[len(models)-1]
        nextCursor = encodeCursor(last.ID, last.CreatedAt)
    }
    
    return users, nextCursor, nil
}
```

---

## 10. Security Best Practices

### Input Validation & Sanitization

```go
// internal/pkg/validator/validator.go
package validator

import (
    "regexp"
    "github.com/go-playground/validator/v10"
)

type CustomValidator struct {
    validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
    v := validator.New()
    
    // Register custom validations
    v.RegisterValidation("password", validatePassword)
    v.RegisterValidation("phone", validatePhone)
    v.RegisterValidation("username", validateUsername)
    
    return &CustomValidator{validator: v}
}

func validatePassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    
    // At least 8 characters
    if len(password) < 8 {
        return false
    }
    
    // At least one uppercase letter
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    // At least one lowercase letter
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    // At least one number
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
    // At least one special character
    hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)
    
    return hasUpper && hasLower && hasNumber && hasSpecial
}

// SQL Injection prevention (using parameterized queries)
func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
    var user UserModel
    
    // Using parameterized query - safe from SQL injection
    err := r.db.WithContext(ctx).
        Where("email = ?", email).
        First(&user).Error
    
    if err != nil {
        return nil, err
    }
    
    return r.toEntity(&user), nil
}
```

### Security Headers Middleware

```go
// internal/infrastructure/http/middleware/security.go
package middleware

import (
    "github.com/gin-gonic/gin"
)

func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Prevent MIME type sniffing
        c.Header("X-Content-Type-Options", "nosniff")
        
        // Enable XSS protection
        c.Header("X-XSS-Protection", "1; mode=block")
        
        // Prevent clickjacking
        c.Header("X-Frame-Options", "DENY")
        
        // HSTS
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        
        // Content Security Policy
        c.Header("Content-Security-Policy", "default-src 'self'")
        
        // Referrer Policy
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        
        c.Next()
    }
}
```

---

## 11. Microservices Patterns

### Service Discovery with Consul

```go
// internal/infrastructure/discovery/consul.go
package discovery

import (
    "github.com/hashicorp/consul/api"
)

type ServiceRegistry struct {
    client *api.Client
    config RegistryConfig
}

type RegistryConfig struct {
    Address     string
    ServiceName string
    ServiceID   string
    ServicePort int
    Tags        []string
}

func NewServiceRegistry(config RegistryConfig) (*ServiceRegistry, error) {
    client, err := api.NewClient(&api.Config{
        Address: config.Address,
    })
    if err != nil {
        return nil, err
    }
    
    return &ServiceRegistry{
        client: client,
        config: config,
    }, nil
}

func (r *ServiceRegistry) Register() error {
    registration := &api.AgentServiceRegistration{
        ID:      r.config.ServiceID,
        Name:    r.config.ServiceName,
        Port:    r.config.ServicePort,
        Tags:    r.config.Tags,
        Check: &api.AgentServiceCheck{
            HTTP:     fmt.Sprintf("http://localhost:%d/health", r.config.ServicePort),
            Interval: "10s",
            Timeout:  "3s",
            DeregisterCriticalServiceAfter: "30s",
        },
    }
    
    return r.client.Agent().ServiceRegister(registration)
}

func (r *ServiceRegistry) Deregister() error {
    return r.client.Agent().ServiceDeregister(r.config.ServiceID)
}

func (r *ServiceRegistry) Discover(serviceName string) ([]string, error) {
    services, _, err := r.client.Health().Service(serviceName, "", true, nil)
    if err != nil {
        return nil, err
    }
    
    addresses := make([]string, len(services))
    for i, service := range services {
        addresses[i] = fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port)
    }
    
    return addresses, nil
}
```

### Circuit Breaker Pattern

```go
// internal/infrastructure/http/client.go
package client

import (
    "github.com/sony/gobreaker"
    "github.com/afex/hystrix-go/hystrix"
)

type CircuitBreakerClient struct {
    client *http.Client
    cb     *gobreaker.CircuitBreaker
}

func NewCircuitBreakerClient() *CircuitBreakerClient {
    // Configure Hystrix
    hystrix.ConfigureCommand("external_service", hystrix.CommandConfig{
        Timeout:               1000,
        MaxConcurrentRequests: 100,
        ErrorPercentThreshold: 25,
        RequestVolumeThreshold: 20,
        SleepWindow:           5000,
    })
    
    return &CircuitBreakerClient{
        client: &http.Client{Timeout: 30 * time.Second},
    }
}

func (c *CircuitBreakerClient) CallExternalAPI(url string) (*http.Response, error) {
    var resp *http.Response
    var err error
    
    err = hystrix.Do("external_service", func() error {
        resp, err = c.client.Get(url)
        return err
    }, func(err error) error {
        // Fallback logic
        return errors.New("service unavailable, using cached data")
    })
    
    return resp, err
}
```

---

## 12. Real-world Case Study

### Complete E-commerce Order Service

```go
// internal/usecase/order_usecase.go
package usecase

import (
    "context"
    "fmt"
    "time"
    "github.com/google/uuid"
    "go.uber.org/zap"
)

type OrderUseCase struct {
    orderRepo     repository.OrderRepository
    userRepo      repository.UserRepository
    productRepo   repository.ProductRepository
    paymentClient PaymentClient
    inventorySvc  InventoryService
    notificationSvc NotificationService
    cache         CacheService
    logger        Logger
    txManager     TransactionManager
}

type CreateOrderRequest struct {
    UserID      uuid.UUID       `json:"user_id"`
    Items       []OrderItem     `json:"items"`
    ShippingAddress ShippingAddress `json:"shipping_address"`
    PaymentMethod string        `json:"payment_method"`
}

type OrderItem struct {
    ProductID uuid.UUID `json:"product_id"`
    Quantity  int       `json:"quantity"`
    Price     float64   `json:"price"`
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, req CreateOrderRequest) (*entity.Order, error) {
    // Start transaction
    tx, err := uc.txManager.Begin(ctx)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()
    
    // 1. Validate user
    user, err := uc.userRepo.FindByID(tx, req.UserID)
    if err != nil {
        return nil, fmt.Errorf("user not found: %w", err)
    }
    
    // 2. Validate inventory and calculate total
    var totalAmount float64
    orderItems := make([]*entity.OrderItem, 0, len(req.Items))
    
    for _, item := range req.Items {
        // Check product
        product, err := uc.productRepo.FindByID(tx, item.ProductID)
        if err != nil {
            return nil, fmt.Errorf("product %s not found: %w", item.ProductID, err)
        }
        
        // Check stock
        if product.Stock < item.Quantity {
            return nil, fmt.Errorf("insufficient stock for product %s", product.Name)
        }
        
        // Reserve stock
        if err := uc.inventorySvc.ReserveStock(tx, item.ProductID, item.Quantity); err != nil {
            return nil, fmt.Errorf("failed to reserve stock: %w", err)
        }
        
        // Calculate
        subtotal := product.Price * float64(item.Quantity)
        totalAmount += subtotal
        
        orderItems = append(orderItems, &entity.OrderItem{
            ProductID: item.ProductID,
            ProductName: product.Name,
            Quantity:  item.Quantity,
            Price:     product.Price,
            Subtotal:  subtotal,
        })
    }
    
    // 3. Create order
    order := &entity.Order{
        ID:              uuid.New(),
        UserID:          user.ID,
        OrderNumber:     generateOrderNumber(),
        Items:           orderItems,
        TotalAmount:     totalAmount,
        ShippingAddress: req.ShippingAddress,
        Status:          entity.OrderStatusPending,
        CreatedAt:       time.Now(),
    }
    
    if err := uc.orderRepo.Create(tx, order); err != nil {
        return nil, fmt.Errorf("failed to create order: %w", err)
    }
    
    // 4. Process payment (async with retry)
    paymentResult, err := uc.paymentClient.ProcessPayment(ctx, &PaymentRequest{
        OrderID:     order.ID,
        Amount:      totalAmount,
        Method:      req.PaymentMethod,
        UserID:      user.ID,
    })
    
    if err != nil {
        // Update order status to failed
        order.Status = entity.OrderStatusPaymentFailed
        order.ErrorMessage = err.Error()
        uc.orderRepo.Update(tx, order)
        
        // Release inventory
        for _, item := range req.Items {
            uc.inventorySvc.ReleaseStock(tx, item.ProductID, item.Quantity)
        }
        
        // Send notification
        uc.notificationSvc.SendPaymentFailedNotification(ctx, user.Email, order.ID)
        
        return nil, fmt.Errorf("payment failed: %w", err)
    }
    
    // 5. Update order with payment info
    order.PaymentID = paymentResult.TransactionID
    order.Status = entity.OrderStatusPaid
    
    if err := uc.orderRepo.Update(tx, order); err != nil {
        return nil, err
    }
    
    // 6. Confirm inventory
    for _, item := range req.Items {
        if err := uc.inventorySvc.ConfirmStock(tx, item.ProductID, item.Quantity); err != nil {
            uc.logger.Error("failed to confirm inventory", 
                zap.Error(err),
                zap.String("order_id", order.ID.String()))
        }
    }
    
    // 7. Commit transaction
    if err := tx.Commit(); err != nil {
        return nil, err
    }
    
    // 8. Async tasks
    go func() {
        // Send confirmation email
        uc.notificationSvc.SendOrderConfirmation(context.Background(), user.Email, order)
        
        // Update analytics
        uc.updateOrderAnalytics(order)
        
        // Cache order
        uc.cache.Set(context.Background(), fmt.Sprintf("order:%s", order.ID), order, 1*time.Hour)
    }()
    
    return order, nil
}

func generateOrderNumber() string {
    return fmt.Sprintf("ORD-%d%03d", time.Now().Unix(), rand.Intn(1000))
}
```

### Monitoring & Alerting Setup

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
    - static_configs:
        - targets: ['alertmanager:9093']

rule_files:
  - "alerts.yml"

scrape_configs:
  - job_name: 'myapp-api'
    static_configs:
      - targets: ['api:9090']
    metrics_path: /metrics
    
  - job_name: 'mysql'
    static_configs:
      - targets: ['mysql-exporter:9104']
      
  - job_name: 'redis'
    static_configs:
      - targets: ['redis-exporter:9121']
```

```yaml
# alerts.yml
groups:
  - name: myapp_alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate on {{ $labels.endpoint }}"
          
      - alert: SlowResponseTime
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Slow response time on {{ $labels.endpoint }}"
          
      - alert: DatabaseConnectionPoolExhausted
        expr: db_connections_active > db_connections_max * 0.9
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "Database connection pool is almost exhausted"
```

---

## สรุป

คู่มือนี้ครอบคลุมการพัฒนา Go สำหรับการใช้งานจริงในระดับ production ตั้งแต่:

1. **สถาปัตยกรรม**: Clean Architecture, Modular Design
2. **ฐานข้อมูล**: Connection Pooling, Migration, Optimization
3. **API**: RESTful, Validation, Documentation
4. **ความปลอดภัย**: Authentication, Authorization, Input Validation
5. **การ监控**: Logging, Metrics, Tracing
6. **การ部署**: Docker, Kubernetes, CI/CD
7. **ประสิทธิภาพ**: Caching, Query Optimization, Concurrency
8. **การทดสอบ**: Unit, Integration, Benchmark
9. **Microservices**: Service Discovery, Circuit Breaker

สำหรับการนำไปใช้งานจริง แนะนำให้:
- ปรับแต่งตามความต้องการของโปรเจค
- ใช้ Go modules สำหรับ dependency management
- ตั้งค่า CI/CD pipeline อัตโนมัติ
- มี monitoring และ alerting ที่ดี
- ทำ code review และ testing อย่างสม่ำเสมอ