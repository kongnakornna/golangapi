# ตัวอย่างโค้ดเต็มรูปแบบสำหรับ GO Shop API
## โครงสร้างโปรเจกต์
```
isekai-shop/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── core/
│   │   ├── domain/
│   │   │   ├── user.go
│   │   │   ├── item.go
│   │   │   ├── player.go
│   │   │   └── inventory.go
│   │   ├── ports/
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   └── services/
│   │       ├── auth_service.go
│   │       ├── user_service.go
│   │       ├── item_service.go
│   │       └── player_service.go
│   ├── adapters/
│   │   ├── repositories/
│   │   │   ├── user_repo.go
│   │   │   ├── item_repo.go
│   │   │   ├── player_repo.go
│   │   │   └── gorm.go
│   │   ├── handlers/
│   │   │   ├── auth_handler.go
│   │   │   ├── user_handler.go
│   │   │   ├── item_handler.go
│   │   │   └── player_handler.go
│   │   ├── middlewares/
│   │   │   ├── auth.go
│   │   │   ├── logger.go
│   │   │   ├── rate_limit.go
│   │   │   └── recover.go
│   │   └── routers/
│   │       └── router.go
│   └── pkg/
│       ├── config/
│       │   └── config.go
│       ├── logger/
│       │   └── logger.go
│       ├── cache/
│       │   └── redis.go
│       └── utils/
│           └── jwt.go
├── migrations/
├── configs/
│   └── app.env
├── go.mod
└── Makefile
```

## 1. การตั้งค่า Configuration (`pkg/config/config.go`)
```go
package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort    string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	RedisAddr     string
	JWTSecret     string
	JWTExpiration int // hours
}

func LoadConfig() *Config {
	err := godotenv.Load("configs/app.env")
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "postgres"),
		DBName:        getEnv("DB_NAME", "isekai"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:     getEnv("JWT_SECRET", "mysecretkey"),
		JWTExpiration: getEnvAsInt("JWT_EXPIRATION", 24),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valStr := getEnv(key, ")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return fallback
}
```

## 2. Domain Models (`internal/core/domain/`)
### user.go
```go
package domain

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"` // hashed
	Email    string `gorm:"uniqueIndex;not null"`
	Role     string `gorm:"default:'user'"`
	Player   *Player
}
```

### item.go
```go
package domain

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	Name        string  `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null"`
	Stock       int     `gorm:"default:0"`
}
```

### player.go
```go
package domain

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	UserID    uint   `gorm:"uniqueIndex;not null"`
	Coins     float64 `gorm:"default:0"`
	Inventory []Inventory
}

type Inventory struct {
	gorm.Model
	PlayerID uint `gorm:"index"`
	ItemID   uint `gorm:"index"`
	Quantity int  `gorm:"default:0"`
}

type PurchaseHistory struct {
	gorm.Model
	PlayerID uint
	ItemID   uint
	Quantity int
	Price    float64 // price at time of purchase
	Type     string  // "buy" or "sell"
}
```

## 3. Repository Interfaces (`internal/core/ports/repository.go`)
```go
package ports

import (
	"context"
	"isekai-shop/internal/core/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]domain.User, error)
}

type ItemRepository interface {
	Create(ctx context.Context, item *domain.Item) error
	GetByID(ctx context.Context, id uint) (*domain.Item, error)
	Update(ctx context.Context, item *domain.Item) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]domain.Item, int64, error)
}

type PlayerRepository interface {
	Create(ctx context.Context, player *domain.Player) error
	GetByUserID(ctx context.Context, userID uint) (*domain.Player, error)
	Update(ctx context.Context, player *domain.Player) error
	AddCoins(ctx context.Context, playerID uint, amount float64) error
	GetInventory(ctx context.Context, playerID uint) ([]domain.Inventory, error)
	AddToInventory(ctx context.Context, playerID, itemID uint, quantity int) error
	RemoveFromInventory(ctx context.Context, playerID, itemID uint, quantity int) error
	RecordPurchase(ctx context.Context, history *domain.PurchaseHistory) error
	// Transaction support
	WithTransaction(ctx context.Context, fn func(tx interface{}) error) error
}
```

## 4. Repository Implementations (`internal/adapters/repositories/`)
### gorm.go (ฐานข้อมูลและการทำ Transaction)
```go
package repositories

import (
	"context"
	"isekai-shop/internal/core/ports"
	"gorm.io/gorm"
)

type GormTransactionManager struct {
	db *gorm.DB
}

func NewGormTransactionManager(db *gorm.DB) ports.TransactionManager {
	return &GormTransactionManager{db: db}
}

func (tm *GormTransactionManager) WithTransaction(ctx context.Context, fn func(tx interface{}) error) error {
	return tm.db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
```

### user_repo.go (บางส่วน)
```go
package repositories

import (
	"context"
	"isekai-shop/internal/core/domain"
	"isekai-shop/internal/core/ports"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
// ... methods อื่น ๆ
```

## 5. Service Layer (Use Cases) (`internal/core/services/`)
### auth_service.go
```go
package services

import (
	"context"
	"errors"
	"isekai-shop/internal/core/domain"
	"isekai-shop/internal/core/ports"
	"isekai-shop/pkg/utils"
	"time"
)

type AuthService struct {
	userRepo ports.UserRepository
	playerRepo ports.PlayerRepository
	jwtSecret string
	jwtExp    int
}

func NewAuthService(userRepo ports.UserRepository, playerRepo ports.PlayerRepository, jwtSecret string, jwtExp int) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		playerRepo: playerRepo,
		jwtSecret: jwtSecret,
		jwtExp: jwtExp,
	}
}

func (s *AuthService) Register(ctx context.Context, username, password, email string) (*domain.User, error) {
	// ตรวจสอบ username ซ้ำ
	existing, _ := s.userRepo.GetByUsername(ctx, username)
	if existing != nil {
		return nil, errors.New("username already exists")
	}
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		Username: username,
		Password: hashed,
		Email:    email,
		Role:     "user",
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	// สร้าง player ให้ผู้ใช้ใหม่
	player := &domain.Player{
		UserID: user.ID,
		Coins:  0,
	}
	if err := s.playerRepo.Create(ctx, player); err != nil {
		// อาจ rollback? ในที่นี้สมมติว่าไม่เป็นไร
	}
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return ", errors.New("invalid credentials")
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return ", errors.New("invalid credentials")
	}
	// สร้าง JWT token
	token, err := utils.GenerateJWT(user.ID, user.Role, s.jwtSecret, time.Duration(s.jwtExp)*time.Hour)
	if err != nil {
		return ", err
	}
	return token, nil
}
```

### player_service.go (จัดการเหรียญ, ซื้อขาย)
```go
package services

import (
	"context"
	"errors"
	"isekai-shop/internal/core/domain"
	"isekai-shop/internal/core/ports"
)

type PlayerService struct {
	playerRepo ports.PlayerRepository
	itemRepo   ports.ItemRepository
	txManager  ports.TransactionManager
}

func NewPlayerService(playerRepo ports.PlayerRepository, itemRepo ports.ItemRepository, txManager ports.TransactionManager) *PlayerService {
	return &PlayerService{
		playerRepo: playerRepo,
		itemRepo:   itemRepo,
		txManager:  txManager,
	}
}

func (s *PlayerService) BuyItem(ctx context.Context, userID uint, itemID uint, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	return s.txManager.WithTransaction(ctx, func(tx interface{}) error {
		// ดึง player
		player, err := s.playerRepo.GetByUserID(ctx, userID)
		if err != nil {
			return err
		}
		// ดึง item
		item, err := s.itemRepo.GetByID(ctx, itemID)
		if err != nil {
			return err
		}
		totalCost := item.Price * float64(quantity)
		if player.Coins < totalCost {
			return errors.New("insufficient coins")
		}
		// หักเหรียญ
		if err := s.playerRepo.AddCoins(ctx, player.ID, -totalCost); err != nil {
			return err
		}
		// เพิ่ม inventory
		if err := s.playerRepo.AddToInventory(ctx, player.ID, itemID, quantity); err != nil {
			return err
		}
		// บันทึกประวัติ
		history := &domain.PurchaseHistory{
			PlayerID: player.ID,
			ItemID:   itemID,
			Quantity: quantity,
			Price:    item.Price,
			Type:     "buy",
		}
		return s.playerRepo.RecordPurchase(ctx, history)
	})
}

func (s *PlayerService) SellItem(ctx context.Context, userID uint, itemID uint, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	return s.txManager.WithTransaction(ctx, func(tx interface{}) error {
		player, err := s.playerRepo.GetByUserID(ctx, userID)
		if err != nil {
			return err
		}
		// ตรวจสอบว่ามีของพอขายไหม
		inv, err := s.playerRepo.GetInventory(ctx, player.ID)
		if err != nil {
			return err
		}
		var currentQty int
		for _, invItem := range inv {
			if invItem.ItemID == itemID {
				currentQty = invItem.Quantity
				break
			}
		}
		if currentQty < quantity {
			return errors.New("not enough items in inventory")
		}
		item, err := s.itemRepo.GetByID(ctx, itemID)
		if err != nil {
			return err
		}
		// ขายได้ราคาเต็ม หรืออาจลดราคาได้
		totalGain := item.Price * float64(quantity)
		// เพิ่มเหรียญ
		if err := s.playerRepo.AddCoins(ctx, player.ID, totalGain); err != nil {
			return err
		}
		// ลด inventory
		if err := s.playerRepo.RemoveFromInventory(ctx, player.ID, itemID, quantity); err != nil {
			return err
		}
		history := &domain.PurchaseHistory{
			PlayerID: player.ID,
			ItemID:   itemID,
			Quantity: quantity,
			Price:    item.Price,
			Type:     "sell",
		}
		return s.playerRepo.RecordPurchase(ctx, history)
	})
}
```

## 6. Handlers (`internal/adapters/handlers/`)
### auth_handler.go
```go
package handlers

import (
	"net/http"
	"isekai-shop/internal/core/services"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	if err := c.Validate(&req); err != nil {
		return err
	}
	user, err := h.authService.Register(c.Request().Context(), req.Username, req.Password, req.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	token, err := h.authService.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
```

### player_handler.go
```go
package handlers

import (
	"net/http"
	"strconv"
	"isekai-shop/internal/core/services"
	"github.com/labstack/echo/v4"
)

type PlayerHandler struct {
	playerService *services.PlayerService
}

func NewPlayerHandler(playerService *services.PlayerService) *PlayerHandler {
	return &PlayerHandler{playerService: playerService}
}

func (h *PlayerHandler) GetCoins(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	player, err := h.playerService.GetByUserID(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "player not found")
	}
	return c.JSON(http.StatusOK, map[string]float64{"coins": player.Coins})
}

type BuyRequest struct {
	ItemID   uint `json:"item_id" validate:"required"`
	Quantity int  `json:"quantity" validate:"required,min=1"`
}

func (h *PlayerHandler) BuyItem(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	var req BuyRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	if err := c.Validate(&req); err != nil {
		return err
	}
	err := h.playerService.BuyItem(c.Request().Context(), userID, req.ItemID, req.Quantity)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func (h *PlayerHandler) SellItem(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	var req BuyRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	err := h.playerService.SellItem(c.Request().Context(), userID, req.ItemID, req.Quantity)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
```

## 7. Middlewares (`internal/adapters/middlewares/`)
### auth.go (ตามที่กล่าวไว้ก่อนหน้า)

### rate_limit.go (ใช้ Redis)
```go
package middlewares

import (
	"net/http"
	"time"
	"github.com/labstack/echo/v4"
	"github.com/go-redis/redis/v8"
	"context"
)

var rdb *redis.Client

func InitRateLimiter(redisAddr string) {
	rdb = redis.NewClient(&redis.Options{Addr: redisAddr})
}

func RateLimiter(limit int, window time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			key := "rate:" + ip
			ctx := context.Background()
			count, err := rdb.Incr(ctx, key).Result()
			if err != nil {
				return next(c) // ถ้า Redis ล้มเหลว ให้ผ่านไปก่อน
			}
			if count == 1 {
				rdb.Expire(ctx, key, window)
			}
			if count > int64(limit) {
				return echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
			}
			return next(c)
		}
	}
}
```

### logger.go (ใช้ slog)
```go
package middlewares

import (
	"time"
	"github.com/labstack/echo/v4"
	"isekai-shop/pkg/logger"
)

func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			stop := time.Now()
			logger.Log.Info("request",
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", c.Response().Status,
				"latency", stop.Sub(start).String(),
				"ip", c.RealIP(),
			)
			return err
		}
	}
}
```

## 8. Main (`cmd/api/main.go`)
```go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"isekai-shop/internal/adapters/handlers"
	"isekai-shop/internal/adapters/middlewares"
	"isekai-shop/internal/adapters/repositories"
	"isekai-shop/internal/core/domain"
	"isekai-shop/internal/core/services"
	"isekai-shop/pkg/config"
	"isekai-shop/pkg/logger"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// โหลด config
	cfg := config.LoadConfig()

	// ตั้งค่า logger
	logger.Init()

	// เชื่อมต่อ PostgreSQL
	dsn := "host=" + cfg.DBHost + " user=" + cfg.DBUser + " password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName + " port=" + cfg.DBPort + " sslmode=disable TimeZone=Asia/Bangkok"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.Error("failed to connect database", "error", err)
		os.Exit(1)
	}

	// Migrate schema
	db.AutoMigrate(&domain.User{}, &domain.Item{}, &domain.Player{}, &domain.Inventory{}, &domain.PurchaseHistory{})

	// เชื่อมต่อ Redis สำหรับ rate limit (และ cache)
	middlewares.InitRateLimiter(cfg.RedisAddr)

	// สร้าง Repositories
	userRepo := repositories.NewUserRepository(db)
	itemRepo := repositories.NewItemRepository(db)
	playerRepo := repositories.NewPlayerRepository(db)
	txManager := repositories.NewGormTransactionManager(db)

	// สร้าง Services
	authService := services.NewAuthService(userRepo, playerRepo, cfg.JWTSecret, cfg.JWTExpiration)
	playerService := services.NewPlayerService(playerRepo, itemRepo, txManager)
	itemService := services.NewItemService(itemRepo) // สมมติว่ามี

	// สร้าง Handlers
	authHandler := handlers.NewAuthHandler(authService)
	playerHandler := handlers.NewPlayerHandler(playerService)
	itemHandler := handlers.NewItemHandler(itemService)

	// สร้าง Echo instance
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middlewares
	e.Use(middlewares.RequestLogger())
	e.Use(middlewares.Recover())
	e.Use(middlewares.RateLimiter(100, time.Minute)) // 100 requests per minute

	// Routes
	api := e.Group("/api/v1")

	// Public routes
	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)
	api.GET("/items", itemHandler.List) // รายการไอเทม (public)

	// Protected routes
	userOnly := api.Group(")
	userOnly.Use(middlewares.AuthMiddleware(cfg.JWTSecret))
	userOnly.GET("/player/coins", playerHandler.GetCoins)
	userOnly.POST("/player/buy", playerHandler.BuyItem)
	userOnly.POST("/player/sell", playerHandler.SellItem)
	// admin routes (ตรวจสอบ role ใน handler)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

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

## 9. Utility Functions (`pkg/utils/`)
### jwt.go
```go
package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(userID uint, role, secret string, exp time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(exp).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
```

## 10. Environment File (`configs/app.env`)
```
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=isekai
REDIS_ADDR=localhost:6379
JWT_SECRET=mysecretkey
JWT_EXPIRATION=24
```

## การรันโปรเจกต์
1. ติดตั้ง Docker และรัน PostgreSQL, Redis:
   ```bash
   docker run --name postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=isekai -p 5432:5432 -d postgres
   docker run --name redis -p 6379:6379 -d redis   ```
2. รันคำสั่ง:
   ```bash
   go mod tidy
   go run cmd/api/main.go
   ```
3. ทดสอบ API ด้วย curl หรือ Postman
## หมายเหตุ
- โค้ดข้างต้นเป็นเพียงตัวอย่างเพื่อให้เห็นภาพรวม ยังขาดการจัดการ error ที่ละเอียด, validation, logging ที่สมบูรณ์, และการทดสอบ
- สำหรับฟีเจอร์อื่น ๆ เช่น OAuth2, caching, queue, สามารถเพิ่มเติมได้ตามต้องการ โดยใช้หลักการเดียวกัน
- การแบ่ง layer ตาม Clean Architecture ช่วยให้เปลี่ยน database หรือ framework ได้ง่าย
 