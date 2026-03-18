# บทความภาษาไทย: การออกแบบสถาปัตยกรรมซอฟต์แวร์ด้วย Golang และ LangChain

## สารบัญ
1. [บทนำ](#1-บทนำ-introduction)
2. [บทนิยาม](#2-บทนิยาม-definitions)
3. [สถาปัตยกรรม Golang แบบ DDD และ Clean Architecture](#3-สถาปัตยกรรม-golang-แบบ-ddd-และ-clean-architecture)
4. [สถาปัตยกรรม LangChain FastAPI](#4-สถาปัตยกรรม-langchain-fastapi)
5. [เวิร์กโฟลว์การพัฒนา](#5-เวิร์กโฟลว์การพัฒนา-development-workflow)
6. [ผังการไหลของข้อมูล](#6-ผังการไหลของข้อมูล-data-flow-architecture)
7. [แม่แบบโค้ด](#7-แม่แบบโค้ด-code-templates)

---

## 1. บทนำ (Introduction)

### 1.1 Golang กับ Domain-Driven Design
Golang หรือ Go เป็นภาษาโปรแกรมมิ่งที่พัฒนาโดย Google มีจุดเด่นด้านประสิทธิภาพสูง การทำงานพร้อมกัน (Concurrency) ที่ยอดเยี่ยม และความเรียบง่ายของภาษา การนำ Domain-Driven Design (DDD) มาประยุกต์ใช้กับ Golang ร่วมกับ Clean Architecture ช่วยให้การพัฒนาซอฟต์แวร์มีโครงสร้างที่ชัดเจน แยกส่วนการทำงานออกจากกันอย่างเป็นระบบ และง่ายต่อการบำรุงรักษา

**ประโยชน์ของ DDD + Clean Architecture ใน Golang:**
- แยก Business Logic ออกจาก Infrastructure
- ทำให้ Testing ทำได้ง่ายขึ้น
- ระบบมีความยืดหยุ่นสูง สามารถเปลี่ยนเทคโนโลยีได้โดยไม่กระทบ Core Business
- ทีมพัฒนาสามารถทำงานแบบขนานได้ (Parallel Development)

### 1.2 LangChain กับ FastAPI
LangChain เป็นเฟรมเวิร์กสำหรับพัฒนาแอปพลิเคชันที่ใช้ Large Language Models (LLMs) ร่วมกับ FastAPI ซึ่งเป็นเว็บเฟรมเวิร์กความเร็วสูงของ Python การรวมกันนี้ช่วยให้สร้าง AI-powered applications ที่มีประสิทธิภาพ

**ประโยชน์ของ LangChain + FastAPI:**
- รองรับการทำงานกับ LLM หลายรูปแบบ (OpenAI, Hugging Face, etc.)
- มีเครื่องมือสำหรับสร้าง Chains, Agents, และ Memory
- FastAPI ให้ประสิทธิภาพสูงและรองรับ async/operation
- มีระบบ Documentation อัตโนมัติ (Swagger/OpenAPI)

---

## 2. บทนิยาม (Definitions)

### 2.1 คำศัพท์สำคัญใน DDD

| คำศัพท์ | คำอธิบาย |
|--------|----------|
| **Domain** | ขอบเขตของปัญหาที่ระบบต้องการแก้ไข |
| **Entity** | วัตถุที่มี Identity ไม่ซ้ำกันและเปลี่ยนแปลงสถานะได้ |
| **Value Object** | วัตถุที่ไม่มี Identity ค่าของมันคือตัวตนของมัน |
| **Aggregate** | กลุ่มของ Entities และ Value Objects ที่ทำงานร่วมกัน |
| **Repository** | ตัวกลางสำหรับเข้าถึงข้อมูล (Persistence) |
| **Domain Service** | Business logic ที่ไม่อยู่ใน Entity หรือ Value Object |
| **Application Service** | จัดการ Use cases และประสานงานระหว่าง Domain |
| **DTO (Data Transfer Object)** | วัตถุสำหรับส่งข้อมูลระหว่าง layers |

### 2.2 คำศัพท์ใน LangChain

| คำศัพท์ | คำอธิบาย |
|--------|----------|
| **Chain** | ลำดับของการเรียกใช้ LLM หรือ tools |
| **Agent** | ระบบที่สามารถตัดสินใจและเลือกใช้ tools |
| **Tool** | ฟังก์ชันที่ Agent สามารถเรียกใช้ได้ |
| **Memory** | การจดจำประวัติการสนทนา |
| **Embedding** | การแปลงข้อความเป็นเวกเตอร์ |
| **Vector Store** | ฐานข้อมูลสำหรับเก็บ embeddings |

---

## 3. สถาปัตยกรรม Golang แบบ DDD และ Clean Architecture

### 3.1 โครงสร้างโปรเจค Golang (Template)

```
project-root/
├── api/                          # เอกสาร API
│   └── app/
│       ├── docs.go
│       └── swagger.json
│
├── cmd/                          # จุดเริ่มต้นโปรแกรม
│   └── app/
│       ├── grpc.go                # gRPC server entry
│       ├── http.go                 # HTTP server entry
│       ├── worker.go               # Background worker entry
│       └── main.go
│
├── configs/                       # ไฟล์การตั้งค่า
│   ├── config.go
│   ├── config.yaml
│   └── config.example.yaml
│
├── deploy/                        # Deployment configs
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   └── k8s/
│       ├── deployment.yaml
│       └── service.yaml
│
├── internal/                      # โค้ดภายใน (ไม่ถูก import จากภายนอก)
│   │
│   ├── core/                      # Domain Layer
│   │   ├── auth/
│   │   │   ├── entity/            # Domain entities
│   │   │   │   ├── user.go
│   │   │   │   ├── token.go
│   │   │   │   └── role.go
│   │   │   │
│   │   │   ├── repository/        # Repository interfaces
│   │   │   │   ├── user_repository.go
│   │   │   │   └── token_repository.go
│   │   │   │
│   │   │   ├── service/           # Domain services
│   │   │   │   ├── auth_service.go
│   │   │   │   └── password_service.go
│   │   │   │
│   │   │   ├── dto/               # Data Transfer Objects
│   │   │   │   ├── auth_request.go
│   │   │   │   └── auth_response.go
│   │   │   │
│   │   │   ├── model/             # Database models (GORM)
│   │   │   │   ├── user_model.go
│   │   │   │   └── token_model.go
│   │   │   │
│   │   │   ├── routes/            # Route definitions
│   │   │   │   └── auth_routes.go
│   │   │   │
│   │   │   ├── handler/           # HTTP handlers (Application Layer)
│   │   │   │   ├── auth_handler.go
│   │   │   │   └── auth_handler_test.go
│   │   │   │
│   │   │   └── auth.go            # Module setup
│   │   │
│   │   ├── user/
│   │   │   ├── entity/
│   │   │   ├── repository/
│   │   │   ├── service/
│   │   │   ├── dto/
│   │   │   ├── model/
│   │   │   ├── handler/
│   │   │   └── user.go
│   │   │
│   │   ├── iot/
│   │   │   ├── entity/
│   │   │   │   ├── device.go
│   │   │   │   ├── sensor.go
│   │   │   │   └── telemetry.go
│   │   │   ├── repository/
│   │   │   ├── service/
│   │   │   ├── handler/
│   │   │   └── iot.go
│   │   │
│   │   └── health/
│   │       ├── handler/
│   │       │   └── health_handler.go
│   │       └── health.go
│   │
│   ├── platform/                   # Infrastructure Layer
│   │   ├── config/
│   │   │   └── config_loader.go
│   │   │
│   │   ├── db/                     # Database connections
│   │   │   ├── postgres/
│   │   │   │   ├── connection.go
│   │   │   │   └── migration.go
│   │   │   └── redis/
│   │   │       ├── connection.go
│   │   │       └── cache.go
│   │   │
│   │   ├── cache/                  # Cache implementations
│   │   │   ├── redis_cache.go
│   │   │   └── cache_interface.go
│   │   │
│   │   ├── queue/                   # Message queue
│   │   │   ├── redis_queue.go
│   │   │   ├── worker_pool.go
│   │   │   └── dead_letter.go
│   │   │
│   │   └── logger/                  # Logging
│   │       ├── logger.go
│   │       └── middleware.go
│   │
│   └── transport/                   # Transport Layer
│       ├── middleware/
│       │   ├── auth_middleware.go
│       │   ├── logging_middleware.go
│       │   ├── recovery_middleware.go
│       │   ├── rate_limit.go
│       │   ├── cors_middleware.go
│       │   └── security_headers.go
│       │
│       ├── httpx/                    # HTTP utilities
│       │   ├── response.go
│       │   ├── request.go
│       │   └── validator.go
│       │
│       └── utils/
│           ├── context.go
│           └── pagination.go
│
├── pkg/                             # Reusable packages
│   ├── cache/
│   │   ├── cache_strategy.go
│   │   └── cache_errors.go
│   │
│   ├── errors/
│   │   ├── errors.go
│   │   ├── app_errors.go
│   │   └── error_handler.go
│   │
│   ├── jwt/
│   │   ├── jwt_service.go
│   │   └── jwt_claims.go
│   │
│   ├── logger/
│   │   ├── logger.go
│   │   └── fields.go
│   │
│   ├── queue/
│   │   ├── queue_interface.go
│   │   ├── redis_queue.go
│   │   └── job.go
│   │
│   ├── transaction/
│   │   └── transaction_manager.go
│   │
│   ├── oauth/
│   │   ├── google_oauth.go
│   │   └── oauth_interface.go
│   │
│   ├── grpc/
│   │   ├── client/
│   │   └── server/
│   │
│   └── utils/
│       ├── crypto.go
│       ├── time.go
│       └── convert.go
│
├── migrations/                       # Database migrations
│   ├── 001_create_users_table.sql
│   ├── 002_create_tokens_table.sql
│   └── 003_create_devices_table.sql
│
├── scripts/                          # Utility scripts
│   ├── setup.sh
│   ├── migrate.sh
│   └── seed.sh
│
├── test/                             # Integration tests
│   ├── integration/
│   └── fixtures/
│
├── .air.toml                         # Hot reload config
├── .env.example
├── go.mod
├── Makefile
└── README.md
```

### 3.2 ตัวอย่างโค้ดสำคัญ

#### Entity (internal/core/auth/entity/user.go)
```go
package entity

import (
    "time"
    "google/uuid"
)

type User struct {
    ID        uuid.UUID
    Email     string
    Password  string
    FirstName string
    LastName  string
    Role      Role
    IsActive  bool
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
}

type Role string

const (
    RoleAdmin Role = "admin"
    RoleUser  Role = "user"
    RoleGuest Role = "guest"
)

func (u *User) IsAdmin() bool {
    return u.Role == RoleAdmin
}

func (u *User) UpdateProfile(firstName, lastName string) {
    u.FirstName = firstName
    u.LastName = lastName
    u.UpdatedAt = time.Now()
}

func (u *User) Deactivate() {
    u.IsActive = false
    now := time.Now()
    u.DeletedAt = &now
}
```

#### Repository Interface (internal/core/auth/repository/user_repository.go)
```go
package repository

import (
    "context"
    "google/uuid"
    "your-project/internal/core/auth/entity"
)

type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id uuid.UUID) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
    FindByEmail(ctx context.Context, email string) (*entity.User, error)
    FindAll(ctx context.Context, limit, offset int) ([]*entity.User, int64, error)
    Exists(ctx context.Context, email string) (bool, error)
}
```

#### Repository Implementation (internal/platform/db/postgres/user_repository_impl.go)
```go
package postgres

import (
    "context"
    "errors"
    "time"
    
    "google/uuid"
    "gorm.io/gorm"
    "your-project/internal/core/auth/entity"
    "your-project/internal/core/auth/repository"
    "your-project/internal/core/auth/model"
)

type userRepositoryImpl struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
    return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
    userModel := &model.UserModel{
        ID:        user.ID,
        Email:     user.Email,
        Password:  user.Password,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        Role:      string(user.Role),
        IsActive:  user.IsActive,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }
    
    return r.db.WithContext(ctx).Create(userModel).Error
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
    var userModel model.UserModel
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&userModel).Error
    
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    return userModel.ToEntity(), nil
}
```

#### Domain Service (internal/core/auth/service/auth_service.go)
```go
package service

import (
    "context"
    "time"
    
    "google/uuid"
    "golang.org/x/crypto/bcrypt"
    "your-project/internal/core/auth/entity"
    "your-project/internal/core/auth/repository"
    "your-project/pkg/errors"
    "your-project/pkg/jwt"
)

type AuthService interface {
    Register(ctx context.Context, req *RegisterRequest) (*entity.User, error)
    Login(ctx context.Context, email, password string) (*AuthResponse, error)
    ValidateToken(ctx context.Context, token string) (*jwt.Claims, error)
    RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)
    Logout(ctx context.Context, userID uuid.UUID) error
}

type authService struct {
    userRepo    repository.UserRepository
    tokenRepo   repository.TokenRepository
    jwtService  jwt.JWTService
    cache       Cache
}

func NewAuthService(
    userRepo repository.UserRepository,
    tokenRepo repository.TokenRepository,
    jwtService jwt.JWTService,
    cache Cache,
) AuthService {
    return &authService{
        userRepo:   userRepo,
        tokenRepo:  tokenRepo,
        jwtService: jwtService,
        cache:      cache,
    }
}

func (s *authService) Register(ctx context.Context, req *RegisterRequest) (*entity.User, error) {
    // Check if user exists
    exists, err := s.userRepo.Exists(ctx, req.Email)
    if err != nil {
        return nil, errors.Wrap(err, "failed to check user existence")
    }
    if exists {
        return nil, errors.NewConflictError("user already exists")
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, errors.Wrap(err, "failed to hash password")
    }
    
    // Create user entity
    now := time.Now()
    user := &entity.User{
        ID:        uuid.New(),
        Email:     req.Email,
        Password:  string(hashedPassword),
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Role:      entity.RoleUser,
        IsActive:  true,
        CreatedAt: now,
        UpdatedAt: now,
    }
    
    // Save to database
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, errors.Wrap(err, "failed to create user")
    }
    
    return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
    // Get user
    user, err := s.userRepo.FindByEmail(ctx, email)
    if err != nil {
        return nil, errors.Wrap(err, "failed to find user")
    }
    if user == nil {
        return nil, errors.NewUnauthorizedError("invalid credentials")
    }
    
    // Check password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, errors.NewUnauthorizedError("invalid credentials")
    }
    
    // Generate tokens
    accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email, string(user.Role))
    if err != nil {
        return nil, errors.Wrap(err, "failed to generate access token")
    }
    
    refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
    if err != nil {
        return nil, errors.Wrap(err, "failed to generate refresh token")
    }
    
    // Save refresh token
    if err := s.tokenRepo.SaveRefreshToken(ctx, user.ID, refreshToken); err != nil {
        return nil, errors.Wrap(err, "failed to save refresh token")
    }
    
    return &AuthResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        TokenType:    "Bearer",
        ExpiresIn:    s.jwtService.GetAccessTokenExpiry(),
    }, nil
}
```

#### Handler (internal/core/auth/handler/auth_handler.go)
```go
package handler

import (
    "net/http"
    
    "gin-gonic/gin"
    "your-project/internal/core/auth/dto"
    "your-project/internal/core/auth/service"
    "your-project/internal/transport/httpx"
    "your-project/pkg/logger"
)

type AuthHandler struct {
    authService service.AuthService
    logger      logger.Logger
}

func NewAuthHandler(authService service.AuthService, logger logger.Logger) *AuthHandler {
    return &AuthHandler{
        authService: authService,
        logger:      logger,
    }
}

// Register godoc
// @Summary Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration request"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        httpx.RespondError(c, httpx.ErrBadRequest(err.Error()))
        return
    }
    
    if err := req.Validate(); err != nil {
        httpx.RespondError(c, httpx.ErrBadRequest(err.Error()))
        return
    }
    
    user, err := h.authService.Register(c.Request.Context(), &req)
    if err != nil {
        h.logger.Error("failed to register user", "error", err)
        httpx.RespondError(c, err)
        return
    }
    
    httpx.RespondSuccess(c, http.StatusCreated, dto.ToRegisterResponse(user))
}

// Login godoc
// @Summary User login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.LoginResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
    var req dto.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        httpx.RespondError(c, httpx.ErrBadRequest(err.Error()))
        return
    }
    
    authResp, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
    if err != nil {
        httpx.RespondError(c, err)
        return
    }
    
    httpx.RespondSuccess(c, http.StatusOK, dto.ToLoginResponse(authResp))
}
```

#### DTO (internal/core/auth/dto/auth_request.go)
```go
package dto

import (
    "regexp"
    "strings"
    
    "go-playground/validator/v10"
)

type RegisterRequest struct {
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=8"`
    FirstName string `json:"first_name" binding:"required"`
    LastName  string `json:"last_name" binding:"required"`
}

func (r *RegisterRequest) Validate() error {
    validate := validator.New()
    
    // Custom validation
    _ = validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
        password := fl.Field().String()
        // At least one uppercase, one lowercase, one number
        hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
        hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
        hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
        
        return hasUpper && hasLower && hasNumber
    })
    
    return validate.Struct(r)
}

func (r *RegisterRequest) Sanitize() {
    r.Email = strings.TrimSpace(strings.ToLower(r.Email))
    r.FirstName = strings.TrimSpace(r.FirstName)
    r.LastName = strings.TrimSpace(r.LastName)
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

func (r *LoginRequest) Sanitize() {
    r.Email = strings.TrimSpace(strings.ToLower(r.Email))
}
```

#### Middleware (internal/transport/middleware/auth_middleware.go)
```go
package middleware

import (
    "net/http"
    "strings"
    
    "gin-gonic/gin"
    "your-project/pkg/jwt"
    "your-project/internal/transport/httpx"
)

type AuthMiddleware struct {
    jwtService jwt.JWTService
}

func NewAuthMiddleware(jwtService jwt.JWTService) *AuthMiddleware {
    return &AuthMiddleware{
        jwtService: jwtService,
    }
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractToken(c)
        if token == " {
            httpx.RespondError(c, httpx.ErrUnauthorized("missing authorization token"))
            c.Abort()
            return
        }
        
        claims, err := m.jwtService.ValidateToken(token)
        if err != nil {
            httpx.RespondError(c, httpx.ErrUnauthorized("invalid token"))
            c.Abort()
            return
        }
        
        // Set user info in context
        c.Set("user_id", claims.UserID)
        c.Set("user_email", claims.Email)
        c.Set("user_role", claims.Role)
        
        c.Next()
    }
}

func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("user_role")
        if !exists {
            httpx.RespondError(c, httpx.ErrForbidden("access denied"))
            c.Abort()
            return
        }
        
        for _, role := range roles {
            if role == userRole {
                c.Next()
                return
            }
        }
        
        httpx.RespondError(c, httpx.ErrForbidden("insufficient permissions"))
        c.Abort()
    }
}

func extractToken(c *gin.Context) string {
    // Check Authorization header
    authHeader := c.GetHeader("Authorization")
    if authHeader != " {
        parts := strings.Split(authHeader, " ")
        if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
            return parts[1]
        }
    }
    
    // Check cookie
    token, err := c.Cookie("access_token")
    if err == nil {
        return token
    }
    
    return "
}
```

#### Rate Limiting Middleware (internal/transport/middleware/rate_limit.go)
```go
package middleware

import (
    "net/http"
    "sync"
    "time"
    
    "gin-gonic/gin"
    "github.com/redis/go-redis/v9"
    "your-project/internal/transport/httpx"
)

type RateLimiter struct {
    redisClient *redis.Client
    limit       int
    window      time.Duration
    mu          sync.Mutex
}

func NewRateLimiter(redisClient *redis.Client, limit int, window time.Duration) *RateLimiter {
    return &RateLimiter{
        redisClient: redisClient,
        limit:       limit,
        window:      window,
    }
}

func (rl *RateLimiter) Limit() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get client IP
        clientIP := c.ClientIP()
        
        // Create key
        key := "rate_limit:" + clientIP
        
        // Use Redis for distributed rate limiting
        val, err := rl.redisClient.Incr(c.Request.Context(), key).Result()
        if err != nil {
            c.Next() // Fail open
            return
        }
        
        if val == 1 {
            // Set expiration on first request
            rl.redisClient.Expire(c.Request.Context(), key, rl.window)
        }
        
        if val > int64(rl.limit) {
            httpx.RespondError(c, httpx.ErrTooManyRequests("rate limit exceeded"))
            c.Abort()
            return
        }
        
        // Set headers
        c.Header("X-RateLimit-Limit", string(rl.limit))
        c.Header("X-RateLimit-Remaining", string(rl.limit-int(val)))
        
        c.Next()
    }
}
```

#### Logger Package (pkg/logger/logger.go)
```go
package logger

import (
    "context"
    "log/slog"
    "os"
    "time"
    
    "google/uuid"
)

type Logger interface {
    Debug(msg string, args ...any)
    Info(msg string, args ...any)
    Warn(msg string, args ...any)
    Error(msg string, args ...any)
    With(args ...any) Logger
    WithContext(ctx context.Context) Logger
}

type logger struct {
    slog *slog.Logger
}

type contextKey string

const (
    RequestIDKey contextKey = "request_id"
    TraceIDKey   contextKey = "trace_id"
    UserIDKey    contextKey = "user_id"
)

func New(level string) Logger {
    var slogLevel slog.Level
    switch level {
    case "debug":
        slogLevel = slog.LevelDebug
    case "info":
        slogLevel = slog.LevelInfo
    case "warn":
        slogLevel = slog.LevelWarn
    case "error":
        slogLevel = slog.LevelError
    default:
        slogLevel = slog.LevelInfo
    }
    
    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slogLevel,
        ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
            if a.Key == slog.TimeKey {
                a.Value = slog.StringValue(time.Now().Format(time.RFC3339Nano))
            }
            return a
        },
    })
    
    return &logger{
        slog: slog.New(handler),
    }
}

func (l *logger) Debug(msg string, args ...any) {
    l.slog.Debug(msg, args...)
}

func (l *logger) Info(msg string, args ...any) {
    l.slog.Info(msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
    l.slog.Warn(msg, args...)
}

func (l *logger) Error(msg string, args ...any) {
    l.slog.Error(msg, args...)
}

func (l *logger) With(args ...any) Logger {
    return &logger{
        slog: l.slog.With(args...),
    }
}

func (l *logger) WithContext(ctx context.Context) Logger {
    args := []any{}
    
    if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
        args = append(args, "request_id", requestID)
    }
    
    if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
        args = append(args, "trace_id", traceID)
    }
    
    if userID, ok := ctx.Value(UserIDKey).(string); ok {
        args = append(args, "user_id", userID)
    }
    
    return l.With(args...)
}
```

#### Error Handling (pkg/errors/errors.go)
```go
package errors

import (
    "fmt"
    "net/http"
)

type ErrorType string

const (
    ErrorTypeNotFound      ErrorType = "NOT_FOUND"
    ErrorTypeUnauthorized  ErrorType = "UNAUTHORIZED"
    ErrorTypeForbidden     ErrorType = "FORBIDDEN"
    ErrorTypeBadRequest    ErrorType = "BAD_REQUEST"
    ErrorTypeConflict      ErrorType = "CONFLICT"
    ErrorTypeInternal      ErrorType = "INTERNAL"
    ErrorTypeTooMany       ErrorType = "TOO_MANY_REQUESTS"
    ErrorTypeValidation    ErrorType = "VALIDATION"
)

type AppError struct {
    Type    ErrorType
    Message string
    Code    string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %s: %v", e.Type, e.Message, e.Err)
    }
    return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
    return e.Err
}

func (e *AppError) StatusCode() int {
    switch e.Type {
    case ErrorTypeNotFound:
        return http.StatusNotFound
    case ErrorTypeUnauthorized:
        return http.StatusUnauthorized
    case ErrorTypeForbidden:
        return http.StatusForbidden
    case ErrorTypeBadRequest:
        return http.StatusBadRequest
    case ErrorTypeConflict:
        return http.StatusConflict
    case ErrorTypeTooMany:
        return http.StatusTooManyRequests
    case ErrorTypeValidation:
        return http.StatusUnprocessableEntity
    default:
        return http.StatusInternalServerError
    }
}

func NewNotFoundError(message string) *AppError {
    return &AppError{
        Type:    ErrorTypeNotFound,
        Message: message,
    }
}

func NewUnauthorizedError(message string) *AppError {
    return &AppError{
        Type:    ErrorTypeUnauthorized,
        Message: message,
    }
}

func NewForbiddenError(message string) *AppError {
    return &AppError{
        Type:    ErrorTypeForbidden,
        Message: message,
    }
}

func NewBadRequestError(message string) *AppError {
    return &AppError{
        Type:    ErrorTypeBadRequest,
        Message: message,
    }
}

func NewConflictError(message string) *AppError {
    return &AppError{
        Type:    ErrorTypeConflict,
        Message: message,
    }
}

func NewInternalError(err error) *AppError {
    return &AppError{
        Type:    ErrorTypeInternal,
        Message: "internal server error",
        Err:     err,
    }
}

func Wrap(err error, message string) *AppError {
    if err == nil {
        return nil
    }
    
    if appErr, ok := err.(*AppError); ok {
        return &AppError{
            Type:    appErr.Type,
            Message: message,
            Err:     appErr,
        }
    }
    
    return &AppError{
        Type:    ErrorTypeInternal,
        Message: message,
        Err:     err,
    }
}
```

#### Cache Implementation (pkg/cache/redis_cache.go)
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
    Incr(ctx context.Context, key string) (int64, error)
    Expire(ctx context.Context, key string, ttl time.Duration) error
    GetOrSet(ctx context.Context, key string, dest interface{}, ttl time.Duration, fn func() (interface{}, error)) error
}

type redisCache struct {
    client *redis.Client
}

func NewRedisCache(client *redis.Client) Cache {
    return &redisCache{
        client: client,
    }
}

func (c *redisCache) Get(ctx context.Context, key string, dest interface{}) error {
    val, err := c.client.Get(ctx, key).Bytes()
    if err == redis.Nil {
        return nil
    }
    if err != nil {
        return err
    }
    
    return json.Unmarshal(val, dest)
}

func (c *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    
    return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
    return c.client.Del(ctx, key).Err()
}

func (c *redisCache) Exists(ctx context.Context, key string) (bool, error) {
    n, err := c.client.Exists(ctx, key).Result()
    return n > 0, err
}

func (c *redisCache) Incr(ctx context.Context, key string) (int64, error) {
    return c.client.Incr(ctx, key).Result()
}

func (c *redisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
    return c.client.Expire(ctx, key, ttl).Err()
}

func (c *redisCache) GetOrSet(ctx context.Context, key string, dest interface{}, ttl time.Duration, fn func() (interface{}, error)) error {
    // Try to get from cache
    err := c.Get(ctx, key, dest)
    if err == nil {
        return nil
    }
    
    // If not found, execute function
    val, err := fn()
    if err != nil {
        return err
    }
    
    // Store in cache
    if err := c.Set(ctx, key, val, ttl); err != nil {
        return err
    }
    
    // Convert to dest
    data, err := json.Marshal(val)
    if err != nil {
        return err
    }
    
    return json.Unmarshal(data, dest)
}
```

#### Queue Implementation (pkg/queue/redis_queue.go)
```go
package queue

import (
    "context"
    "encoding/json"
    "time"
    
    "github.com/redis/go-redis/v9"
    "google/uuid"
)

type Job struct {
    ID        string          `json:"id"`
    Type      string          `json:"type"`
    Payload   json.RawMessage `json:"payload"`
    Attempts  int             `json:"attempts"`
    MaxAttempts int           `json:"max_attempts"`
    CreatedAt time.Time       `json:"created_at"`
    UpdatedAt time.Time       `json:"updated_at"`
}

type Queue interface {
    Push(ctx context.Context, queue string, job *Job) error
    Pop(ctx context.Context, queue string) (*Job, error)
    Ack(ctx context.Context, queue string, jobID string) error
    Nack(ctx context.Context, queue string, jobID string) error
    Len(ctx context.Context, queue string) (int64, error)
}

type redisQueue struct {
    client *redis.Client
}

func NewRedisQueue(client *redis.Client) Queue {
    return &redisQueue{
        client: client,
    }
}

func (q *redisQueue) Push(ctx context.Context, queue string, job *Job) error {
    if job.ID == " {
        job.ID = uuid.New().String()
    }
    if job.CreatedAt.IsZero() {
        job.CreatedAt = time.Now()
    }
    job.UpdatedAt = time.Now()
    
    data, err := json.Marshal(job)
    if err != nil {
        return err
    }
    
    // Push to main queue
    return q.client.RPush(ctx, "queue:"+queue, data).Err()
}

func (q *redisQueue) Pop(ctx context.Context, queue string) (*Job, error) {
    // Pop from main queue and push to processing queue atomically
    result, err := q.client.LPop(ctx, "queue:"+queue).Bytes()
    if err == redis.Nil {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    var job Job
    if err := json.Unmarshal(result, &job); err != nil {
        return nil, err
    }
    
    // Store in processing queue with TTL
    processingKey := "processing:" + queue + ":" + job.ID
    if err := q.client.Set(ctx, processingKey, result, 5*time.Minute).Err(); err != nil {
        return nil, err
    }
    
    return &job, nil
}

func (q *redisQueue) Ack(ctx context.Context, queue string, jobID string) error {
    // Remove from processing queue
    return q.client.Del(ctx, "processing:"+queue+":"+jobID).Err()
}

func (q *redisQueue) Nack(ctx context.Context, queue string, jobID string) error {
    // Get job from processing queue
    data, err := q.client.GetDel(ctx, "processing:"+queue+":"+jobID).Bytes()
    if err != nil {
        return err
    }
    
    var job Job
    if err := json.Unmarshal(data, &job); err != nil {
        return err
    }
    
    // Increment attempts
    job.Attempts++
    job.UpdatedAt = time.Now()
    
    if job.Attempts >= job.MaxAttempts {
        // Move to dead letter queue
        return q.client.RPush(ctx, "dead:"+queue, data).Err()
    }
    
    // Push back to main queue
    return q.Push(ctx, queue, &job)
}

func (q *redisQueue) Len(ctx context.Context, queue string) (int64, error) {
    return q.client.LLen(ctx, "queue:"+queue).Result()
}
```

#### Worker Pool (internal/platform/queue/worker_pool.go)
```go
package queue

import (
    "context"
    "sync"
    "time"
    
    "your-project/pkg/logger"
    "your-project/pkg/queue"
)

type WorkerFunc func(ctx context.Context, job *queue.Job) error

type WorkerPool struct {
    queue      queue.Queue
    workers    int
    jobQueue   string
    handler    WorkerFunc
    logger     logger.Logger
    wg         sync.WaitGroup
    stopChan   chan struct{}
}

func NewWorkerPool(
    queue queue.Queue,
    workers int,
    jobQueue string,
    handler WorkerFunc,
    logger logger.Logger,
) *WorkerPool {
    return &WorkerPool{
        queue:    queue,
        workers:  workers,
        jobQueue: jobQueue,
        handler:  handler,
        logger:   logger,
        stopChan: make(chan struct{}),
    }
}

func (wp *WorkerPool) Start(ctx context.Context) {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker(i)
    }
}

func (wp *WorkerPool) Stop() {
    close(wp.stopChan)
    wp.wg.Wait()
}

func (wp *WorkerPool) worker(id int) {
    defer wp.wg.Done()
    
    ctx := context.Background()
    
    for {
        select {
        case <-wp.stopChan:
            return
        default:
            job, err := wp.queue.Pop(ctx, wp.jobQueue)
            if err != nil {
                wp.logger.Error("failed to pop job", "error", err, "worker", id)
                time.Sleep(time.Second)
                continue
            }
            
            if job == nil {
                time.Sleep(100 * time.Millisecond)
                continue
            }
            
            // Process job
            if err := wp.handler(ctx, job); err != nil {
                wp.logger.Error("job processing failed", 
                    "error", err, 
                    "job_id", job.ID,
                    "worker", id,
                )
                
                // Negative acknowledgment
                if err := wp.queue.Nack(ctx, wp.jobQueue, job.ID); err != nil {
                    wp.logger.Error("failed to nack job", "error", err)
                }
            } else {
                // Positive acknowledgment
                if err := wp.queue.Ack(ctx, wp.jobQueue, job.ID); err != nil {
                    wp.logger.Error("failed to ack job", "error", err)
                }
            }
        }
    }
}
```

#### JWT Service (pkg/jwt/jwt_service.go)
```go
package jwt

import (
    "time"
    
    "golang-jwt/jwt/v5"
    "google/uuid"
)

type Claims struct {
    UserID    uuid.UUID `json:"user_id"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    jwt.RegisteredClaims
}

type JWTService interface {
    GenerateAccessToken(userID uuid.UUID, email, role string) (string, error)
    GenerateRefreshToken(userID uuid.UUID) (string, error)
    ValidateToken(token string) (*Claims, error)
    GetAccessTokenExpiry() int64
}

type jwtService struct {
    secretKey     []byte
    accessExpiry  time.Duration
    refreshExpiry time.Duration
    issuer        string
}

func NewJWTService(secretKey string, accessExpiry, refreshExpiry time.Duration, issuer string) JWTService {
    return &jwtService{
        secretKey:     []byte(secretKey),
        accessExpiry:  accessExpiry,
        refreshExpiry: refreshExpiry,
        issuer:        issuer,
    }
}

func (s *jwtService) GenerateAccessToken(userID uuid.UUID, email, role string) (string, error) {
    now := time.Now()
    claims := &Claims{
        UserID: userID,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpiry)),
            IssuedAt:  jwt.NewNumericDate(now),
            NotBefore: jwt.NewNumericDate(now),
            Issuer:    s.issuer,
            Subject:   userID.String(),
            ID:        uuid.New().String(),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.secretKey)
}

func (s *jwtService) GenerateRefreshToken(userID uuid.UUID) (string, error) {
    now := time.Now()
    claims := &jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshExpiry)),
        IssuedAt:  jwt.NewNumericDate(now),
        NotBefore: jwt.NewNumericDate(now),
        Issuer:    s.issuer,
        Subject:   userID.String(),
        ID:        uuid.New().String(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.secretKey)
}

func (s *jwtService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return s.secretKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, jwt.ErrTokenInvalidClaims
}

func (s *jwtService) GetAccessTokenExpiry() int64 {
    return int64(s.accessExpiry.Seconds())
}
```

#### Transaction Manager (pkg/transaction/transaction_manager.go)
```go
package transaction

import (
    "context"
    "database/sql"
    
    "gorm.io/gorm"
)

type TransactionManager interface {
    Begin(ctx context.Context) (context.Context, error)
    Commit(ctx context.Context) error
    Rollback(ctx context.Context) error
    GetDB(ctx context.Context) *gorm.DB
}

type transactionManager struct {
    db *gorm.DB
}

type txKey struct{}

func NewTransactionManager(db *gorm.DB) TransactionManager {
    return &transactionManager{
        db: db,
    }
}

func (tm *transactionManager) Begin(ctx context.Context) (context.Context, error) {
    tx := tm.db.Begin()
    if tx.Error != nil {
        return ctx, tx.Error
    }
    
    return context.WithValue(ctx, txKey{}, tx), nil
}

func (tm *transactionManager) Commit(ctx context.Context) error {
    tx, ok := ctx.Value(txKey{}).(*gorm.DB)
    if !ok {
        return sql.ErrTxDone
    }
    
    return tx.Commit().Error
}

func (tm *transactionManager) Rollback(ctx context.Context) error {
    tx, ok := ctx.Value(txKey{}).(*gorm.DB)
    if !ok {
        return sql.ErrTxDone
    }
    
    return tx.Rollback().Error
}

func (tm *transactionManager) GetDB(ctx context.Context) *gorm.DB {
    tx, ok := ctx.Value(txKey{}).(*gorm.DB)
    if ok {
        return tx
    }
    
    return tm.db.WithContext(ctx)
}
```

#### Main Application (cmd/app/main.go)
```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "gin-gonic/gin"
    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"
    
    "your-project/configs"
    "your-project/internal/platform/config"
    "your-project/internal/platform/db/postgres"
    "your-project/internal/platform/logger"
    "your-project/internal/platform/cache"
    "your-project/internal/platform/queue"
    "your-project/internal/transport/middleware"
    "your-project/internal/core/auth/handler"
    "your-project/internal/core/auth/repository"
    "your-project/internal/core/auth/service"
    authRepo "your-project/internal/platform/db/postgres/repository"
    "your-project/pkg/jwt"
    "your-project/pkg/transaction"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    // Initialize logger
    appLogger := logger.New(cfg.LogLevel)
    appLogger.Info("Starting application", "version", cfg.Version)
    
    // Initialize database
    db, err := postgres.NewConnection(cfg.Database)
    if err != nil {
        appLogger.Error("Failed to connect to database", "error", err)
        os.Exit(1)
    }
    
    // Initialize Redis
    redisClient := redis.NewClient(&redis.Options{
        Addr:     cfg.Redis.Addr,
        Password: cfg.Redis.Password,
        DB:       cfg.Redis.DB,
    })
    
    if err := redisClient.Ping(context.Background()).Err(); err != nil {
        appLogger.Error("Failed to connect to Redis", "error", err)
        os.Exit(1)
    }
    
    // Initialize components
    cacheService := cache.NewRedisCache(redisClient)
    queueService := queue.NewRedisQueue(redisClient)
    jwtService := jwt.NewJWTService(
        cfg.JWT.Secret,
        cfg.JWT.AccessExpiry,
        cfg.JWT.RefreshExpiry,
        cfg.JWT.Issuer,
    )
    txManager := transaction.NewTransactionManager(db)
    
    // Initialize repositories
    userRepo := authRepo.NewUserRepository(db)
    tokenRepo := authRepo.NewTokenRepository(db)
    
    // Initialize services
    authService := service.NewAuthService(
        userRepo,
        tokenRepo,
        jwtService,
        cacheService,
    )
    
    // Initialize handlers
    authHandler := handler.NewAuthHandler(authService, appLogger)
    
    // Initialize middleware
    authMiddleware := middleware.NewAuthMiddleware(jwtService)
    rateLimiter := middleware.NewRateLimiter(redisClient, 100, time.Minute)
    
    // Setup router
    router := gin.New()
    router.Use(
        gin.Recovery(),
        middleware.LoggingMiddleware(appLogger),
        middleware.SecurityHeaders(),
        middleware.CORS(cfg.CORS),
        rateLimiter.Limit(),
    )
    
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // API v1 routes
    v1 := router.Group("/api/v1")
    {
        // Public routes
        auth := v1.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)
            auth.POST("/refresh", authHandler.Refresh)
        }
        
        // Protected routes
        protected := v1.Group(")
        protected.Use(authMiddleware.RequireAuth())
        {
            protected.POST("/auth/logout", authHandler.Logout)
            protected.GET("/auth/me", authHandler.GetMe)
        }
        
        // Admin routes
        admin := v1.Group("/admin")
        admin.Use(authMiddleware.RequireAuth())
        admin.Use(authMiddleware.RequireRole("admin"))
        {
            // Admin routes here
        }
    }
    
    // Start server
    srv := &http.Server{
        Addr:    ":" + cfg.Port,
        Handler: router,
    }
    
    go func() {
        appLogger.Info("Server starting", "port", cfg.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            appLogger.Error("Server failed", "error", err)
            os.Exit(1)
        }
    }()
    
    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    appLogger.Info("Shutting down server...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        appLogger.Error("Server forced to shutdown", "error", err)
    }
    
    appLogger.Info("Server exited")
}
```

### 3.3 Configuration (configs/config.yaml)
```yaml
server:
  port: 8080
  mode: development
  version: 1.0.0

database:
  driver: postgres
  host: localhost
  port: 5432
  username: postgres
  password: postgres
  database: myapp
  sslmode: disable
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 1h

redis:
  addr: localhost:6379
  password: "
  db: 0
  pool_size: 10

jwt:
  secret: your-secret-key
  access_expiry: 15m
  refresh_expiry: 168h
  issuer: myapp

cors:
  allowed_origins:
    - http://localhost:3000
  allowed_methods:
    - GET
    - POST
    - PUT
    - DELETE
    - OPTIONS
  allowed_headers:
    - Content-Type
    - Authorization

log:
  level: info
  format: json
```

### 3.4 Makefile
```makefile
.PHONY: build run test migrate clean docker-up docker-down

build:
	go build -o bin/app cmd/app/main.go

run:
	go run cmd/app/main.go

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

migrate-up:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/myapp?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/myapp?sslmode=disable" down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

lint:
	golangci-lint run

swagger:
	swag init -g cmd/app/main.go -o api/app

clean:
	rm -rf bin/
	rm -f coverage.out
```

---

## 4. สถาปัตยกรรม LangChain FastAPI

### 4.1 โครงสร้างโปรเจค LangChain FastAPI

```
langchain-project/
├── app/
│   ├── __init__.py
│   ├── main.py                    # Entry point
│   │
│   ├── core/                       # Domain Layer
│   │   ├── __init__.py
│   │   ├── entities/
│   │   │   ├── __init__.py
│   │   │   ├── document.py
│   │   │   ├── conversation.py
│   │   │   └── message.py
│   │   │
│   │   ├── interfaces/
│   │   │   ├── __init__.py
│   │   │   ├── repository.py
│   │   │   └── llm_service.py
│   │   │
│   │   └── value_objects/
│   │       ├── __init__.py
│   │       └── embedding.py
│   │
│   ├── application/                 # Application Layer
│   │   ├── __init__.py
│   │   ├── services/
│   │   │   ├── __init__.py
│   │   │   ├── chat_service.py
│   │   │   ├── document_service.py
│   │   │   ├── embedding_service.py
│   │   │   └── agent_service.py
│   │   │
│   │   ├── dtos/
│   │   │   ├── __init__.py
│   │   │   ├── chat_dto.py
│   │   │   └── document_dto.py
│   │   │
│   │   └── use_cases/
│   │       ├── __init__.py
│   │       ├── chat_use_case.py
│   │       └── qa_use_case.py
│   │
│   ├── infrastructure/              # Infrastructure Layer
│   │   ├── __init__.py
│   │   ├── database/
│   │   │   ├── __init__.py
│   │   │   ├── postgres.py
│   │   │   ├── redis.py
│   │   │   └── repositories/
│   │   │       ├── __init__.py
│   │   │       ├── document_repository.py
│   │   │       └── conversation_repository.py
│   │   │
│   │   ├── llm/
│   │   │   ├── __init__.py
│   │   │   ├── openai_service.py
│   │   │   ├── huggingface_service.py
│   │   │   └── chain_factory.py
│   │   │
│   │   ├── vector_store/
│   │   │   ├── __init__.py
│   │   │   ├── chroma_store.py
│   │   │   ├── pinecone_store.py
│   │   │   └── qdrant_store.py
│   │   │
│   │   ├── cache/
│   │   │   ├── __init__.py
│   │   │   └── redis_cache.py
│   │   │
│   │   └── queue/
│   │       ├── __init__.py
│   │       ├── redis_queue.py
│   │       └── tasks.py
│   │
│   ├── api/                          # Interface Layer
│   │   ├── __init__.py
│   │   ├── dependencies/
│   │   │   ├── __init__.py
│   │   │   └── container.py
│   │   │
│   │   ├── middleware/
│   │   │   ├── __init__.py
│   │   │   ├── auth.py
│   │   │   ├── logging.py
│   │   │   └── rate_limit.py
│   │   │
│   │   ├── routes/
│   │   │   ├── __init__.py
│   │   │   ├── chat_routes.py
│   │   │   ├── document_routes.py
│   │   │   └── health_routes.py
│   │   │
│   │   └── models/
│   │       ├── __init__.py
│   │       ├── request_models.py
│   │       └── response_models.py
│   │
│   ├── config/
│   │   ├── __init__.py
│   │   ├── settings.py
│   │   └── logging_config.py
│   │
│   └── utils/
│       ├── __init__.py
│       ├── errors.py
│       ├── logger.py
│       └── helpers.py
│
├── migrations/
│   ├── versions/
│   └── env.py
│
├── tests/
│   ├── __init__.py
│   ├── unit/
│   ├── integration/
│   └── conftest.py
│
├── scripts/
│   ├── seed_data.py
│   └── migrate_db.py
│
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yml
│
├── .env.example
├── .gitignore
├── requirements.txt
├── requirements-dev.txt
├── alembic.ini
└── README.md
```

### 4.2 ตัวอย่างโค้ด LangChain

#### Entity (app/core/entities/document.py)
```python
from dataclasses import dataclass
from datetime import datetime
from typing import Optional, Dict, Any
from uuid import UUID, uuid4

@dataclass
class Document:
    id: UUID
    title: str
    content: str
    metadata: Dict[str, Any]
    embedding: Optional[list[float]]
    created_at: datetime
    updated_at: datetime
    user_id: Optional[UUID]
    
    @classmethod
    def create(cls, title: str, content: str, user_id: Optional[UUID] = None) -> "Document":
        now = datetime.utcnow()
        return cls(
            id=uuid4(),
            title=title,
            content=content,
            metadata={},
            embedding=None,
            created_at=now,
            updated_at=now,
            user_id=user_id
        )
    
    def update_embedding(self, embedding: list[float]) -> None:
        self.embedding = embedding
        self.updated_at = datetime.utcnow()
    
    def add_metadata(self, key: str, value: Any) -> None:
        self.metadata[key] = value
        self.updated_at = datetime.utcnow()
```

#### Service (app/application/services/chat_service.py)
```python
from typing import List, Optional, Dict, Any
from uuid import UUID
from langchain.chains import ConversationChain
from langchain.memory import ConversationBufferMemory
from langchain.schema import BaseMessage, HumanMessage, AIMessage

from app.core.entities.conversation import Conversation, Message
from app.core.interfaces.repository import ConversationRepository
from app.core.interfaces.llm_service import LLMService
from app.application.dtos.chat_dto import ChatRequest, ChatResponse
from app.infrastructure.cache.redis_cache import RedisCache


class ChatService:
    def __init__(
        self,
        conversation_repo: ConversationRepository,
        llm_service: LLMService,
        cache: RedisCache
    ):
        self.conversation_repo = conversation_repo
        self.llm_service = llm_service
        self.cache = cache
        self.chains: Dict[str, ConversationChain] = {}
    
    async def chat(
        self,
        request: ChatRequest,
        user_id: Optional[UUID] = None
    ) -> ChatResponse:
        # Get or create conversation
        if request.conversation_id:
            conversation = await self.conversation_repo.get_by_id(
                request.conversation_id
            )
        else:
            conversation = Conversation.create(user_id=user_id)
        
        # Create message
        user_message = Message(
            role="user",
            content=request.message,
            metadata=request.metadata
        )
        conversation.add_message(user_message)
        
        # Get or create chain
        chain = await self._get_or_create_chain(conversation.id)
        
        # Generate response
        response = await chain.acall(inputs={"input": request.message})
        
        # Create AI message
        ai_message = Message(
            role="assistant",
            content=response["response"],
            metadata={"tokens": response.get("tokens", {})}
        )
        conversation.add_message(ai_message)
        
        # Save conversation
        await self.conversation_repo.save(conversation)
        
        return ChatResponse(
            conversation_id=conversation.id,
            message=ai_message.content,
            metadata=ai_message.metadata,
            created_at=ai_message.created_at
        )
    
    async def _get_or_create_chain(self, conversation_id: UUID) -> ConversationChain:
        # Check cache
        cache_key = f"chain:{conversation_id}"
        cached_chain = await self.cache.get(cache_key)
        if cached_chain:
            return cached_chain
        
        # Get conversation history
        conversation = await self.conversation_repo.get_by_id(conversation_id)
        
        # Create memory with history
        memory = ConversationBufferMemory()
        for msg in conversation.messages[-10:]:  # Last 10 messages
            if msg.role == "user":
                memory.chat_memory.add_user_message(msg.content)
            else:
                memory.chat_memory.add_ai_message(msg.content)
        
        # Create chain
        chain = ConversationChain(
            llm=self.llm_service.get_llm(),
            memory=memory,
            verbose=True
        )
        
        # Cache chain
        await self.cache.set(cache_key, chain, ttl=3600)
        
        return chain
    
    async def get_conversation_history(
        self,
        conversation_id: UUID,
        limit: int = 50
    ) -> List[Dict[str, Any]]:
        conversation = await self.conversation_repo.get_by_id(conversation_id)
        
        return [
            {
                "role": msg.role,
                "content": msg.content,
                "created_at": msg.created_at.isoformat(),
                "metadata": msg.metadata
            }
            for msg in conversation.messages[-limit:]
        ]
```

#### Document Service with LangChain (app/application/services/document_service.py)
```python
from typing import List, Optional, Dict, Any
from uuid import UUID
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain.embeddings import OpenAIEmbeddings
from langchain.vectorstores import Chroma
from langchain.schema import Document as LangChainDocument

from app.core.entities.document import Document
from app.core.interfaces.repository import DocumentRepository
from app.infrastructure.vector_store.chroma_store import ChromaStore


class DocumentService:
    def __init__(
        self,
        document_repo: DocumentRepository,
        vector_store: ChromaStore,
        embeddings: OpenAIEmbeddings
    ):
        self.document_repo = document_repo
        self.vector_store = vector_store
        self.embeddings = embeddings
        self.text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=1000,
            chunk_overlap=200,
            length_function=len,
            separators=["\n\n", "\n", " ", "]
        )
    
    async def create_document(
        self,
        title: str,
        content: str,
        user_id: Optional[UUID] = None,
        metadata: Optional[Dict[str, Any]] = None
    ) -> Document:
        # Create document entity
        document = Document.create(
            title=title,
            content=content,
            user_id=user_id
        )
        
        if metadata:
            document.metadata.update(metadata)
        
        # Split document into chunks
        chunks = self.text_splitter.split_text(content)
        
        # Create LangChain documents
        lc_documents = [
            LangChainDocument(
                page_content=chunk,
                metadata={
                    "document_id": str(document.id),
                    "chunk_index": i,
                    "title": title,
                    **metadata
                }
            )
            for i, chunk in enumerate(chunks)
        ]
        
        # Generate embeddings and store in vector store
        await self.vector_store.add_documents(lc_documents)
        
        # Save to database
        await self.document_repo.save(document)
        
        return document
    
    async def search_similar(
        self,
        query: str,
        limit: int = 5,
        user_id: Optional[UUID] = None
    ) -> List[Dict[str, Any]]:
        # Search in vector store
        results = await self.vector_store.similarity_search(
            query=query,
            k=limit,
            filter={"user_id": str(user_id)} if user_id else None
        )
        
        return [
            {
                "content": doc.page_content,
                "metadata": doc.metadata,
                "score": score
            }
            for doc, score in results
        ]
    
    async def get_document_with_context(
        self,
        document_id: UUID,
        query: str
    ) -> Dict[str, Any]:
        # Get document from database
        document = await self.document_repo.get_by_id(document_id)
        
        # Search for relevant chunks
        similar = await self.vector_store.similarity_search(
            query=query,
            k=3,
            filter={"document_id": str(document_id)}
        )
        
        context = "\n\n".join([doc.page_content for doc, _ in similar])
        
        return {
            "document": document,
            "context": context,
            "relevant_chunks": [
                {
                    "content": doc.page_content,
                    "score": score
                }
                for doc, score in similar
            ]
        }
```

#### Agent Service (app/application/services/agent_service.py)
```python
from typing import List, Dict, Any, Optional
from langchain.agents import Tool, AgentExecutor, create_react_agent
from langchain.prompts import PromptTemplate
from langchain.chains import LLMMathChain
from langchain.tools import DuckDuckGoSearchRun

from app.core.interfaces.llm_service import LLMService
from app.application.services.document_service import DocumentService


class AgentService:
    def __init__(
        self,
        llm_service: LLMService,
        document_service: DocumentService
    ):
        self.llm_service = llm_service
        self.document_service = document_service
        self.agent_executor = self._create_agent()
    
    def _create_agent(self) -> AgentExecutor:
        # Create tools
        tools = [
            Tool(
                name="Search",
                func=DuckDuckGoSearchRun().run,
                description="Search the internet for current information"
            ),
            Tool(
                name="Calculator",
                func=LLMMathChain.from_llm(
                    self.llm_service.get_llm()
                ).run,
                description="Perform mathematical calculations"
            ),
            Tool(
                name="DocumentSearch",
                func=self._search_documents,
                description="Search through internal documents"
            )
        ]
        
        # Create prompt
        prompt = PromptTemplate.from_template(""
        You are a helpful AI assistant with access to various tools.
        Use the following tools to help answer the user's question:
        
        {tools}
        
        Question: {input}
        
        Thought process: {agent_scratchpad}
        "")
        
        # Create agent
        agent = create_react_agent(
            llm=self.llm_service.get_llm(),
            tools=tools,
            prompt=prompt
        )
        
        return AgentExecutor(
            agent=agent,
            tools=tools,
            verbose=True,
            handle_parsing_errors=True,
            max_iterations=5
        )
    
    async def _search_documents(self, query: str) -> str:
        ""Search through internal documents""
        results = await self.document_service.search_similar(
            query=query,
            limit=3
        )
        
        if not results:
            return "No relevant documents found."
        
        return "\n\n".join([
            f"[Document: {r['metadata'].get('title', 'Unknown')}]\n{r['content']}"
            for r in results
        ])
    
    async def run_agent(
        self,
        query: str,
        user_id: Optional[str] = None
    ) -> Dict[str, Any]:
        ""Run the agent with a query""
        result = await self.agent_executor.arun(
            input=query,
            user_id=user_id
        )
        
        return {
            "response": result,
            "agent_used": True
        }
```

#### DTOs (app/application/dtos/chat_dto.py)
```python
from pydantic import BaseModel, Field
from typing import Optional, Dict, Any
from uuid import UUID
from datetime import datetime


class ChatRequest(BaseModel):
    message: str = Field(..., min_length=1, max_length=4000)
    conversation_id: Optional[UUID] = None
    metadata: Optional[Dict[str, Any]] = Field(default_factory=dict)
    stream: bool = False
    
    class Config:
        json_schema_extra = {
            "example": {
                "message": "What is the capital of France?",
                "metadata": {"source": "web"}
            }
        }


class ChatResponse(BaseModel):
    conversation_id: UUID
    message: str
    metadata: Dict[str, Any]
    created_at: datetime


class ConversationSummary(BaseModel):
    id: UUID
    title: Optional[str]
    message_count: int
    last_message_at: datetime
    created_at: datetime
```

#### Vector Store (app/infrastructure/vector_store/chroma_store.py)
```python
from typing import List, Optional, Dict, Any
import chromadb
from chromadb.config import Settings
from langchain.vectorstores import Chroma
from langchain.embeddings import OpenAIEmbeddings
from langchain.schema import Document


class ChromaStore:
    def __init__(
        self,
        collection_name: str,
        persist_directory: str,
        embeddings: OpenAIEmbeddings
    ):
        self.client = chromadb.Client(Settings(
            chroma_db_impl="duckdb+parquet",
            persist_directory=persist_directory
        ))
        
        self.vectorstore = Chroma(
            collection_name=collection_name,
            embedding_function=embeddings,
            client=self.client,
            persist_directory=persist_directory
        )
    
    async def add_documents(
        self,
        documents: List[Document],
        ids: Optional[List[str]] = None
    ) -> List[str]:
        ""Add documents to vector store""
        return self.vectorstore.add_documents(
            documents=documents,
            ids=ids
        )
    
    async def similarity_search(
        self,
        query: str,
        k: int = 4,
        filter: Optional[Dict[str, Any]] = None
    ) -> List[tuple[Document, float]]:
        ""Search for similar documents""
        results = self.vectorstore.similarity_search_with_score(
            query=query,
            k=k,
            filter=filter
        )
        return results
    
    async def delete_document(self, document_id: str) -> None:
        ""Delete document by ID""
        self.vectorstore.delete(filter={"document_id": document_id})
    
    async def get_collection_stats(self) -> Dict[str, Any]:
        ""Get collection statistics""
        collection = self.client.get_collection(
            self.vectorstore._collection.name
        )
        return {
            "count": collection.count(),
            "name": collection.name,
            "metadata": collection.metadata
        }
```

#### LLM Service (app/infrastructure/llm/openai_service.py)
```python
from typing import Optional, Dict, Any
from langchain.chat_models import ChatOpenAI
from langchain.llms import OpenAI
from langchain.schema import BaseMessage, HumanMessage, AIMessage
from tenacity import retry, stop_after_attempt, wait_exponential

from app.core.interfaces.llm_service import LLMService
from app.config.settings import Settings


class OpenAIService(LLMService):
    def __init__(self, settings: Settings):
        self.settings = settings
        self.chat_model = ChatOpenAI(
            model=settings.OPENAI_MODEL,
            temperature=settings.OPENAI_TEMPERATURE,
            max_tokens=settings.OPENAI_MAX_TOKENS,
            openai_api_key=settings.OPENAI_API_KEY
        )
        self.llm = OpenAI(
            model=settings.OPENAI_MODEL,
            temperature=settings.OPENAI_TEMPERATURE,
            openai_api_key=settings.OPENAI_API_KEY
        )
    
    @retry(
        stop=stop_after_attempt(3),
        wait=wait_exponential(multiplier=1, min=4, max=10)
    )
    async def generate_response(
        self,
        messages: List[BaseMessage],
        temperature: Optional[float] = None
    ) -> str:
        ""Generate response from messages""
        response = await self.chat_model.agenerate([messages])
        return response.generations[0][0].text
    
    @retry(
        stop=stop_after_attempt(3),
        wait=wait_exponential(multiplier=1, min=4, max=10)
    )
    async def generate_completion(
        self,
        prompt: str,
        max_tokens: Optional[int] = None
    ) -> str:
        ""Generate completion from prompt""
        response = await self.llm.agenerate([prompt])
        return response.generations[0][0].text
    
    async def stream_response(
        self,
        messages: List[BaseMessage]
    ) -> Any:
        ""Stream response""
        async for chunk in self.chat_model.astream(messages):
            yield chunk
    
    def get_llm(self):
        ""Get LangChain LLM instance""
        return self.chat_model
    
    def get_embeddings(self):
        ""Get embeddings model""
        from langchain.embeddings import OpenAIEmbeddings
        return OpenAIEmbeddings(
            openai_api_key=self.settings.OPENAI_API_KEY
        )
```

#### API Routes (app/api/routes/chat_routes.py)
```python
from fastapi import APIRouter, Depends, HTTPException, BackgroundTasks
from typing import List
from uuid import UUID

from app.api.dependencies.container import get_chat_service
from app.api.models.request_models import ChatRequestModel
from app.api.models.response_models import ChatResponseModel, ConversationModel
from app.application.services.chat_service import ChatService
from app.api.middleware.auth import get_current_user

router = APIRouter(prefix="/chat", tags=["chat"])


@router.post("/", response_model=ChatResponseModel)
async def chat(
    request: ChatRequestModel,
    background_tasks: BackgroundTasks,
    chat_service: ChatService = Depends(get_chat_service),
    user_id: Optional[UUID] = Depends(get_current_user)
):
    ""
    Send a message to the chat
    ""
    try:
        response = await chat_service.chat(
            request=request,
            user_id=user_id
        )
        
        # Background task for analytics
        background_tasks.add_task(
            log_conversation,
            conversation_id=response.conversation_id,
            user_id=user_id
        )
        
        return response
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/conversations/{conversation_id}", response_model=ConversationModel)
async def get_conversation(
    conversation_id: UUID,
    chat_service: ChatService = Depends(get_chat_service),
    user_id: UUID = Depends(get_current_user)
):
    ""
    Get conversation history
    ""
    messages = await chat_service.get_conversation_history(
        conversation_id=conversation_id
    )
    return ConversationModel(
        id=conversation_id,
        messages=messages
    )


@router.delete("/conversations/{conversation_id}")
async def delete_conversation(
    conversation_id: UUID,
    chat_service: ChatService = Depends(get_chat_service),
    user_id: UUID = Depends(get_current_user)
):
    ""
    Delete a conversation
    ""
    await chat_service.delete_conversation(conversation_id)
    return {"message": "Conversation deleted"}


async def log_conversation(conversation_id: UUID, user_id: Optional[UUID]):
    ""Background task for logging""
    # Implement logging logic
    pass
```

#### Main Application (app/main.py)
```python
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
import uvicorn

from app.api.routes import chat_routes, document_routes, health_routes
from app.api.middleware.logging import LoggingMiddleware
from app.api.middleware.rate_limit import RateLimitMiddleware
from app.config.settings import settings
from app.utils.errors import AppException, app_exception_handler
from app.infrastructure.database.postgres import init_db
from app.infrastructure.cache.redis import redis_client


def create_app() -> FastAPI:
    ""Application factory""
    app = FastAPI(
        title="LangChain API",
        version=settings.VERSION,
        docs_url="/api/docs" if settings.ENVIRONMENT == "development" else None,
        redoc_url="/api/redoc" if settings.ENVIRONMENT == "development" else None
    )
    
    # Setup middleware
    app.add_middleware(
        CORSMiddleware,
        allow_origins=settings.CORS_ORIGINS,
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )
    app.add_middleware(LoggingMiddleware)
    app.add_middleware(RateLimitMiddleware)
    
    # Setup exception handlers
    app.add_exception_handler(AppException, app_exception_handler)
    
    # Include routers
    app.include_router(health_routes.router, prefix="/api")
    app.include_router(chat_routes.router, prefix="/api/v1")
    app.include_router(document_routes.router, prefix="/api/v1")
    
    @app.on_event("startup")
    async def startup():
        ""Startup tasks""
        await init_db()
        await redis_client.initialize()
    
    @app.on_event("shutdown")
    async def shutdown():
        ""Shutdown tasks""
        await redis_client.close()
    
    return app


app = create_app()

if __name__ == "__main__":
    uvicorn.run(
        "app.main:app",
        host=settings.HOST,
        port=settings.PORT,
        reload=settings.ENVIRONMENT == "development"
    )
```

#### Settings (app/config/settings.py)
```python
from pydantic_settings import BaseSettings
from typing import List, Optional
from functools import lru_cache


class Settings(BaseSettings):
    # App
    APP_NAME: str = "LangChain API"
    VERSION: str = "1.0.0"
    ENVIRONMENT: str = "development"
    DEBUG: bool = False
    HOST: str = "0.0.0.0"
    PORT: int = 8000
    
    # Database
    DATABASE_URL: str = "postgresql://user:pass@localhost:5432/langchain"
    DATABASE_POOL_SIZE: int = 20
    
    # Redis
    REDIS_URL: str = "redis://localhost:6379/0"
    REDIS_MAX_CONNECTIONS: int = 10
    
    # OpenAI
    OPENAI_API_KEY: str
    OPENAI_MODEL: str = "gpt-3.5-turbo"
    OPENAI_TEMPERATURE: float = 0.7
    OPENAI_MAX_TOKENS: int = 1000
    
    # Vector Store
    VECTOR_STORE_TYPE: str = "chroma"
    CHROMA_PERSIST_DIR: str = "./chroma_db"
    PINECONE_API_KEY: Optional[str] = None
    PINECONE_ENVIRONMENT: Optional[str] = None
    
    # Security
    SECRET_KEY: str
    JWT_ALGORITHM: str = "HS256"
    ACCESS_TOKEN_EXPIRE_MINUTES: int = 30
    CORS_ORIGINS: List[str] = ["http://localhost:3000"]
    
    # Rate Limiting
    RATE_LIMIT_REQUESTS: int = 100
    RATE_LIMIT_PERIOD: int = 60  # seconds
    
    # Logging
    LOG_LEVEL: str = "INFO"
    LOG_FORMAT: str = "json"
    
    class Config:
        env_file = ".env"
        case_sensitive = True


@lru_cache()
def get_settings() -> Settings:
    return Settings()


settings = get_settings()
```

#### Docker Compose (docker/docker-compose.yml)
```yaml
version: '3.8'

services:
  api:
    build:
      context: ..
      dockerfile: docker/Dockerfile
    ports:
      - "8000:8000"
    environment:
      - DATABASE_URL=postgresql://postgres:postgres@db:5432/langchain
      - REDIS_URL=redis://redis:6379/0
      - OPENAI_API_KEY=${OPENAI_API_KEY}
    volumes:
      - ../app:/app/app
      - chroma_data:/app/chroma_db
    depends_on:
      - db
      - redis
    networks:
      - langchain-network
    command: uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload

  worker:
    build:
      context: ..
      dockerfile: docker/Dockerfile
    environment:
      - DATABASE_URL=postgresql://postgres:postgres@db:5432/langchain
      - REDIS_URL=redis://redis:6379/0
      - OPENAI_API_KEY=${OPENAI_API_KEY}
    depends_on:
      - db
      - redis
    networks:
      - langchain-network
    command: python -m app.worker

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=langchain
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - langchain-network

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - langchain-network

  chroma:
    image: chromadb/chroma:latest
    ports:
      - "8001:8000"
    volumes:
      - chroma_data:/chroma/chroma
    networks:
      - langchain-network

volumes:
  postgres_data:
  redis_data:
  chroma_data:

networks:
  langchain-network:
    driver: bridge
```

---

## 5. เวิร์กโฟลว์การพัฒนา (Development Workflow)

### 5.1 Golang Development Workflow

```mermaid
graph TD
    A[Write Code] --> B[Run Tests]
    B --> C{Tests Pass?}
    C -->|No| A
    C -->|Yes| D[Run Linter]
    D --> E{Lint Pass?}
    E -->|No| A
    E -->|Yes| F[Build Binary]
    F --> G[Run Migrations]
    G --> H[Start Application]
    H --> I[Test API]
    I --> J{All Good?}
    J -->|No| A
    J -->|Yes| K[Create PR]
```

**ขั้นตอนการพัฒนา:**
1. **Local Development**: ใช้ Air สำหรับ hot-reload
   ```bash
   air
   ```

2. **Testing**:
   ```bash
   make test              # Run all tests
   make test-coverage     # Check coverage
   ```

3. **Database Migration**:
   ```bash
   make migrate-create name=add_users_table
   make migrate-up
   ```

4. **Generate Swagger Docs**:
   ```bash
   make swagger
   ```

5. **Build and Run**:
   ```bash
   make build
   ./bin/app
   ```

### 5.2 LangChain FastAPI Development Workflow

```mermaid
graph TD
    A[Define Requirements] --> B[Create/Update Models]
    B --> C[Write Business Logic]
    C --> D[Create API Endpoints]
    D --> E[Write Tests]
    E --> F{Tests Pass?}
    F -->|No| C
    F -->|Yes| G[Run Locally]
    G --> H{Test API}
    H -->|Issues| C
    H -->|OK| I[Create PR]
```

**ขั้นตอนการพัฒนา:**
1. **Setup Environment**:
   ```bash
   python -m venv venv
   source venv/bin/activate
   pip install -r requirements-dev.txt
   ```

2. **Run Locally**:
   ```bash
   uvicorn app.main:app --reload
   ```

3. **Database Migrations**:
   ```bash
   alembic revision --autogenerate -m "message"
   alembic upgrade head
   ```

4. **Testing**:
   ```bash
   pytest
   pytest --cov=app tests/
   ```

5. **Docker Development**:
   ```bash
   docker-compose up -d
   docker-compose logs -f api
   ```

---

## 6. ผังการไหลของข้อมูล (Data Flow Architecture)

### 6.1 Golang Clean Architecture Data Flow

```mermaid
graph LR
    Client[Client] -->|HTTP Request| Router[Router]
    Router -->|Route| Middleware[Middleware Stack]
    Middleware -->|Auth/Rate Limit| Handler[Handler]
    Handler -->|DTO| Service[Service Layer]
    Service -->|Entity| Repository[Repository Interface]
    Repository -->|Implementation| DB[(Database)]
    
    Service -->|Cache Key| Cache[Redis Cache]
    Service -->|Job| Queue[Message Queue]
    
    Cache -->|Cached Data| Service
    Queue -->|Async Job| Worker[Worker Pool]
    Worker -->|Process| Service
    
    Service -->|Business Logic| Domain[Domain Layer]
    Domain -->|Events| EventBus[Event Bus]
    EventBus -->|Handle| EventHandler[Event Handler]
```

**ลำดับการไหลของข้อมูล:**
1. Client ส่ง HTTP Request มา
2. Router จับคู่กับ Handler ที่เหมาะสม
3. Middleware ทำการตรวจสอบ (Auth, Rate Limit, Logging)
4. Handler รับ Request และแปลงเป็น DTO
5. Service Layer ประมวลผล Business Logic
6. Repository Interface ถูกเรียกเพื่อ存取 Data
7. Cache ถูกตรวจสอบก่อน Query Database
8. Database ส่งข้อมูลกลับ
9. Service สร้าง Events (ถ้ามี)
10. Response ถูกส่งกลับไปยัง Client

### 6.2 LangChain FastAPI Data Flow

```mermaid
graph TD
    User[User] -->|Query| API[FastAPI Endpoint]
    API -->|Validate| DTO[Request DTO]
    DTO -->|Call| Service[Service Layer]
    
    Service -->|Check| Cache[(Redis Cache)]
    Cache -->|Cache Hit| Service
    
    Service -->|Cache Miss| Chain[LangChain Chain]
    
    Chain -->|Retrieve| VectorStore[(Vector Store)]
    Chain -->|LLM Call| LLM[OpenAI/LLM]
    Chain -->|Tool Use| Tools[Tools/Search]
    
    VectorStore -->|Similar Docs| Chain
    LLM -->|Response| Chain
    Tools -->|Results| Chain
    
    Chain -->|Final Response| Service
    Service -->|Store| DB[(PostgreSQL)]
    Service -->|Cache| Cache
    
    Service -->|Response DTO| API
    API -->|JSON| User
    
    Service -->|Async Task| Queue[Message Queue]
    Queue -->|Process| Worker[Background Worker]
    Worker -->|Update| DB
```

**อธิบายการไหล:**
1. **Request Phase**: User ส่ง query → API validation → Service layer
2. **Cache Phase**: Service ตรวจ cache → ถ้าเจอ return ทันที
3. **Processing Phase**: ถ้าไม่เจอ cache → LangChain chain ทำงาน
4. **Retrieval Phase**: Chain query vector store → get relevant documents
5. **Generation Phase**: Chain เรียก LLM with context → generate response
6. **Post-processing**: Response → cache → database → return to user

### 6.3 การไหลของข้อมูลแบบ Streaming

```mermaid
sequenceDiagram
    participant User
    participant API
    participant Service
    participant Chain
    participant LLM
    
    User->>API: POST /chat/stream
    API->>Service: stream_chat(request)
    Service->>Chain: astream(messages)
    
    loop For each token
        Chain->>LLM: get_next_token()
        LLM-->>Chain: token
        Chain-->>Service: token
        Service-->>API: token
        API-->>User: SSE: token
    end
    
    Chain-->>Service: complete
    Service-->>API: done
    API-->>User: SSE: [DONE]
```

---

## 7. แม่แบบโค้ด (Code Templates)

### 7.1 Golang Template Files

#### Repository Template (internal/core/[module]/repository/repository.go)
```go
package repository

import (
    "context"
    "google/uuid"
    "your-project/internal/core/[module]/entity"
)

type [Module]Repository interface {
    Create(ctx context.Context, entity *entity.[Entity]) error
    Update(ctx context.Context, entity *entity.[Entity]) error
    Delete(ctx context.Context, id uuid.UUID) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.[Entity], error)
    FindAll(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]*entity.[Entity], int64, error)
}
```

#### Service Template (internal/core/[module]/service/service.go)
```go
package service

import (
    "context"
    "google/uuid"
    "your-project/internal/core/[module]/entity"
    "your-project/internal/core/[module]/repository"
    "your-project/pkg/logger"
)

type [Module]Service interface {
    Create(ctx context.Context, req *CreateRequest) (*entity.[Entity], error)
    GetByID(ctx context.Context, id uuid.UUID) (*entity.[Entity], error)
    Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*entity.[Entity], error)
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, filter map[string]interface{}, page, pageSize int) (*ListResponse, error)
}

type [module]Service struct {
    repo   repository.[Module]Repository
    cache  cache.Cache
    logger logger.Logger
}

func New[Module]Service(
    repo repository.[Module]Repository,
    cache cache.Cache,
    logger logger.Logger,
) [Module]Service {
    return &[module]Service{
        repo:   repo,
        cache:  cache,
        logger: logger,
    }
}

func (s *[module]Service) Create(ctx context.Context, req *CreateRequest) (*entity.[Entity], error) {
    // Validation
    if err := req.Validate(); err != nil {
        return nil, errors.NewBadRequestError(err.Error())
    }
    
    // Business logic
    entity := &entity.[Entity]{
        ID:        uuid.New(),
        // Map fields from req
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    // Save
    if err := s.repo.Create(ctx, entity); err != nil {
        s.logger.Error("failed to create [module]", "error", err)
        return nil, errors.Wrap(err, "failed to create [module]")
    }
    
    // Invalidate cache
    cacheKey := fmt.Sprintf("[module]:%s", entity.ID)
    _ = s.cache.Delete(ctx, cacheKey)
    
    return entity, nil
}
```

#### Handler Template (internal/core/[module]/handler/handler.go)
```go
package handler

import (
    "net/http"
    
    "gin-gonic/gin"
    "google/uuid"
    
    "your-project/internal/core/[module]/dto"
    "your-project/internal/core/[module]/service"
    "your-project/internal/transport/httpx"
    "your-project/pkg/logger"
)

type [Module]Handler struct {
    service service.[Module]Service
    logger  logger.Logger
}

func New[Module]Handler(service service.[Module]Service, logger logger.Logger) *[Module]Handler {
    return &[Module]Handler{
        service: service,
        logger:  logger,
    }
}

// Create godoc
// @Summary Create new [module]
// @Tags [module]
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateRequest true "Create request"
// @Success 201 {object} dto.[Module]Response
// @Failure 400 {object} httpx.ErrorResponse
// @Router /api/v1/[module] [post]
func (h *[Module]Handler) Create(c *gin.Context) {
    var req dto.CreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        httpx.RespondError(c, httpx.ErrBadRequest(err.Error()))
        return
    }
    
    entity, err := h.service.Create(c.Request.Context(), &req)
    if err != nil {
        httpx.RespondError(c, err)
        return
    }
    
    httpx.RespondSuccess(c, http.StatusCreated, dto.ToResponse(entity))
}

// GetByID godoc
// @Summary Get [module] by ID
// @Tags [module]
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Module ID"
// @Success 200 {object} dto.[Module]Response
// @Failure 404 {object} httpx.ErrorResponse
// @Router /api/v1/[module]/{id} [get]
func (h *[Module]Handler) GetByID(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        httpx.RespondError(c, httpx.ErrBadRequest("invalid id format"))
        return
    }
    
    entity, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        httpx.RespondError(c, err)
        return
    }
    
    httpx.RespondSuccess(c, http.StatusOK, dto.ToResponse(entity))
}
```

#### DTO Template (internal/core/[module]/dto/dto.go)
```go
package dto

import (
    "time"
    
    "go-playground/validator/v10"
    "google/uuid"
    "your-project/internal/core/[module]/entity"
)

type CreateRequest struct {
    // Add fields
    Name  string `json:"name" binding:"required"`
    Field string `json:"field" binding:"required"`
}

func (r *CreateRequest) Validate() error {
    validate := validator.New()
    return validate.Struct(r)
}

type UpdateRequest struct {
    Name  string `json:"name" binding:"omitempty"`
    Field string `json:"field" binding:"omitempty"`
}

func (r *UpdateRequest) Validate() error {
    validate := validator.New()
    return validate.Struct(r)
}

type [Module]Response struct {
    ID        uuid.UUID `json:"id"`
    Name      string    `json:"name"`
    Field     string    `json:"field"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func ToResponse(entity *entity.[Entity]) *[Module]Response {
    return &[Module]Response{
        ID:        entity.ID,
        Name:      entity.Name,
        Field:     entity.Field,
        CreatedAt: entity.CreatedAt,
        UpdatedAt: entity.UpdatedAt,
    }
}

type ListResponse struct {
    Items      []*[Module]Response `json:"items"`
    Total      int64                `json:"total"`
    Page       int                  `json:"page"`
    PageSize   int                  `json:"page_size"`
    TotalPages int                  `json:"total_pages"`
}
```

### 7.2 LangChain FastAPI Template Files

#### Service Template (app/application/services/base_service.py)
```python
from typing import TypeVar, Generic, List, Optional, Dict, Any
from uuid import UUID
from abc import ABC, abstractmethod

from app.core.entities.base import BaseEntity
from app.core.interfaces.repository import BaseRepository
from app.infrastructure.cache.redis_cache import RedisCache

T = TypeVar('T', bound=BaseEntity)


class BaseService(Generic[T], ABC):
    def __init__(
        self,
        repository: BaseRepository[T],
        cache: Optional[RedisCache] = None
    ):
        self.repository = repository
        self.cache = cache
    
    async def get_by_id(self, id: UUID) -> Optional[T]:
        # Check cache first
        if self.cache:
            cached = await self.cache.get(f"{self.entity_name}:{id}")
            if cached:
                return cached
        
        # Get from database
        entity = await self.repository.get_by_id(id)
        
        # Cache if found
        if entity and self.cache:
            await self.cache.set(
                f"{self.entity_name}:{id}",
                entity,
                ttl=300
            )
        
        return entity
    
    async def create(self, data: Dict[str, Any]) -> T:
        entity = self.entity_class.create(**data)
        created = await self.repository.create(entity)
        
        # Invalidate list cache
        if self.cache:
            await self.cache.delete_pattern(f"{self.entity_name}:list:*")
        
        return created
    
    async def update(self, id: UUID, data: Dict[str, Any]) -> Optional[T]:
        entity = await self.get_by_id(id)
        if not entity:
            return None
        
        for key, value in data.items():
            if hasattr(entity, key):
                setattr(entity, key, value)
        
        updated = await self.repository.update(entity)
        
        # Invalidate caches
        if self.cache:
            await self.cache.delete(f"{self.entity_name}:{id}")
            await self.cache.delete_pattern(f"{self.entity_name}:list:*")
        
        return updated
    
    async def delete(self, id: UUID) -> bool:
        result = await self.repository.delete(id)
        
        # Invalidate caches
        if self.cache and result:
            await self.cache.delete(f"{self.entity_name}:{id}")
            await self.cache.delete_pattern(f"{self.entity_name}:list:*")
        
        return result
    
    async def list(
        self,
        skip: int = 0,
        limit: int = 100,
        **filters
    ) -> List[T]:
        cache_key = f"{self.entity_name}:list:{skip}:{limit}:{hash(frozenset(filters.items()))}"
        
        # Check cache
        if self.cache:
            cached = await self.cache.get(cache_key)
            if cached:
                return cached
        
        # Get from database
        items = await self.repository.list(skip=skip, limit=limit, **filters)
        
        # Cache results
        if self.cache:
            await self.cache.set(cache_key, items, ttl=60)
        
        return items
```

#### Chain Template (app/infrastructure/llm/chains/qa_chain.py)
```python
from typing import List, Optional, Dict, Any
from langchain.chains import RetrievalQA
from langchain.prompts import PromptTemplate
from langchain.callbacks import AsyncCallbackHandler
from langchain.schema import BaseRetriever

from app.infrastructure.llm.base import BaseLLMService
from app.infrastructure.vector_store.base import BaseVectorStore


class QAChain:
    ""Template for QA chain with retrieval""
    
    def __init__(
        self,
        llm_service: BaseLLMService,
        vector_store: BaseVectorStore,
        retriever_kwargs: Optional[Dict[str, Any]] = None
    ):
        self.llm_service = llm_service
        self.vector_store = vector_store
        self.retriever_kwargs = retriever_kwargs or {"k": 4}
        
        # Create retriever
        self.retriever = self.vector_store.as_retriever(
            **self.retriever_kwargs
        )
        
        # Create prompt template
        self.prompt = PromptTemplate(
            template=""
            Use the following pieces of context to answer the question at the end.
            If you don't know the answer, just say that you don't know, don't try to make up an answer.
            
            Context: {context}
            
            Question: {question}
            
            Answer: "",
            input_variables=["context", "question"]
        )
        
        # Create chain
        self.chain = RetrievalQA.from_chain_type(
            llm=self.llm_service.get_llm(),
            chain_type="stuff",
            retriever=self.retriever,
            chain_type_kwargs={
                "prompt": self.prompt,
                "verbose": True
            },
            return_source_documents=True
        )
    
    async def ainvoke(
        self,
        query: str,
        callbacks: Optional[List[AsyncCallbackHandler]] = None
    ) -> Dict[str, Any]:
        ""Async invoke chain""
        result = await self.chain.ainvoke(
            {"query": query},
            callbacks=callbacks
        )
        
        return {
            "answer": result["result"],
            "sources": [
                {
                    "content": doc.page_content,
                    "metadata": doc.metadata
                }
                for doc in result["source_documents"]
            ]
        }
```

#### API Route Template (app/api/routes/module_routes.py)
```python
from fastapi import APIRouter, Depends, HTTPException, Query
from typing import List, Optional
from uuid import UUID

from app.api.dependencies.container import get_service
from app.api.models.request_models import CreateRequest, UpdateRequest
from app.api.models.response_models import EntityResponse, ListResponse
from app.application.services.base_service import BaseService
from app.api.middleware.auth import get_current_user, require_permissions

router = APIRouter(prefix="/module", tags=["module"])


@router.post("/", response_model=EntityResponse, status_code=201)
async def create_entity(
    request: CreateRequest,
    service: BaseService = Depends(get_service),
    user_id: UUID = Depends(get_current_user)
):
    ""
    Create a new entity
    ""
    try:
        entity = await service.create(request.dict())
        return EntityResponse.from_entity(entity)
    except Exception as e:
        raise HTTPException(status_code=400, detail=str(e))


@router.get("/{entity_id}", response_model=EntityResponse)
async def get_entity(
    entity_id: UUID,
    service: BaseService = Depends(get_service),
    user_id: UUID = Depends(get_current_user)
):
    ""
    Get entity by ID
    ""
    entity = await service.get_by_id(entity_id)
    if not entity:
        raise HTTPException(status_code=404, detail="Entity not found")
    
    return EntityResponse.from_entity(entity)


@router.get("/", response_model=ListResponse)
async def list_entities(
    skip: int = Query(0, ge=0),
    limit: int = Query(100, ge=1, le=1000),
    search: Optional[str] = None,
    service: BaseService = Depends(get_service),
    user_id: UUID = Depends(get_current_user)
):
    ""
    List entities with pagination
    ""
    filters = {}
    if search:
        filters["search"] = search
    
    items, total = await service.list(
        skip=skip,
        limit=limit,
        **filters
    )
    
    return ListResponse(
        items=[EntityResponse.from_entity(item) for item in items],
        total=total,
        skip=skip,
        limit=limit
    )


@router.put("/{entity_id}", response_model=EntityResponse)
async def update_entity(
    entity_id: UUID,
    request: UpdateRequest,
    service: BaseService = Depends(get_service),
    user_id: UUID = Depends(get_current_user),
    _: bool = Depends(require_permissions(["admin"]))
):
    ""
    Update entity (admin only)
    ""
    entity = await service.update(
        entity_id,
        request.dict(exclude_unset=True)
    )
    if not entity:
        raise HTTPException(status_code=404, detail="Entity not found")
    
    return EntityResponse.from_entity(entity)


@router.delete("/{entity_id}", status_code=204)
async def delete_entity(
    entity_id: UUID,
    service: BaseService = Depends(get_service),
    user_id: UUID = Depends(get_current_user),
    _: bool = Depends(require_permissions(["admin"]))
):
    ""
    Delete entity (admin only)
    ""
    deleted = await service.delete(entity_id)
    if not deleted:
        raise HTTPException(status_code=404, detail="Entity not found")
```

## บทสรุป

การออกแบบสถาปัตยกรรมซอฟต์แวร์ด้วย Golang และ LangChain FastAPI ตามหลักการ Domain-Driven Design (DDD) และ Clean Architecture ช่วยให้ระบบมีโครงสร้างที่ชัดเจน แยกส่วนการทำงานออกจากกัน ทำให้:

1. **Maintainability**: ง่ายต่อการบำรุงรักษาและพัฒนาเพิ่มเติม
2. **Testability**: สามารถทดสอบแต่ละส่วนได้อย่างอิสระ
3. **Scalability**: ขยายระบบได้ง่ายตามความต้องการ
4. **Flexibility**: เปลี่ยนเทคโนโลยีได้โดยไม่กระทบ Business Logic
5. **Team Development**: ทีมสามารถทำงานแบบ Parallel ได้

### การเลือกใช้ระหว่าง Golang และ LangChain FastAPI

| พิจารณา | Golang | LangChain FastAPI |
|--------|--------|-------------------|
| **Performance** | สูงมาก | ปานกลาง |
| **Concurrency** | ยอดเยี่ยม (Goroutines) | ดี (Async/Await) |
| **AI/ML Features** | จำกัด | ยอดเยี่ยม |
| **Ecosystem** | ดีสำหรับ Backend | ดีสำหรับ Data/ML |
| **Learning Curve** | ปานกลาง | ต่ำ |

**แนะนำให้ใช้ Golang เมื่อ:**
- ต้องการ performance สูง
- ระบบมี traffic มาก
- ต้องการการทำงานแบบ concurrent
- ต้องการ type safety

**แนะนำให้ใช้ LangChain FastAPI เมื่อ:**
- ต้องการ integrate กับ LLMs
- ต้องการพัฒนา AI applications
- ต้องการ prototyping ที่รวดเร็ว
- ทีมมีประสบการณ์ Python

### แหล่งข้อมูลเพิ่มเติม

1. **Golang Resources**:
   - [Official Go Documentation](https://golang.org/doc/)
   - [Clean Architecture in Go](https://bxcodec/go-clean-arch)
   - [GORM Documentation](https://gorm.io/docs/)

2. **LangChain Resources**:
   - [LangChain Documentation](https://python.langchain.com/docs)
   - [FastAPI Documentation](https://fastapi.tiangolo.com/)
   - [OpenAI API Documentation](https://platform.openai.com/docs)

3. **DDD Resources**:
   - [Domain-Driven Design in Go](https://marcusolsson/goddd)
   - [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

**จัดทำโดย**: ทีมพัฒนาซอฟต์แวร์  
**ปรับปรุงล่าสุด**: 2024  
**เวอร์ชัน**: 1.0.0