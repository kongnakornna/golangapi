# Module 5: Delivery Layer (Handlers, Middleware, DTOs, Router)

## สำหรับโฟลเดอร์ `internal/delivery/rest/`

ไฟล์ที่เกี่ยวข้อง:
- `internal/delivery/rest/handler/auth_handler.go`
- `internal/delivery/rest/handler/user_handler.go`
- `internal/delivery/rest/handler/health_handler.go`
- `internal/delivery/rest/middleware/auth.go`
- `internal/delivery/rest/middleware/cors.go`
- `internal/delivery/rest/middleware/logger.go`
- `internal/delivery/rest/middleware/rate_limit.go`
- `internal/delivery/rest/middleware/security.go`
- `internal/delivery/rest/middleware/monitoring.go`
- `internal/delivery/rest/dto/auth_dto.go`
- `internal/delivery/rest/dto/user_dto.go`
- `internal/delivery/rest/dto/error_dto.go`
- `internal/delivery/rest/router.go`

---

## หลักการ (Concept)

### คืออะไร?
Delivery layer เป็นชั้นที่อยู่ด้านนอกสุดของ Clean Architecture ทำหน้าที่รับ request จากผู้ใช้ (HTTP, gRPC, CLI) แปลงข้อมูล, เรียกใช้ usecase, และส่ง response กลับ โดยไม่มีการประมวลผลทางธุรกิจใดๆ

### มีกี่แบบ?
1. **HTTP/REST handlers** – รับ JSON, เรียก usecase, ส่ง JSON response
2. **Middleware** – ทำงานก่อน/หลัง handlers (authentication, logging, CORS, rate limiting)
3. **WebSocket handlers** – จัดการ real-time connections
4. **gRPC services** – สำหรับ internal microservices
5. **CLI commands** – สำหรับ admin tasks (migrate, seed)

### ใช้อย่างไร / นำไปใช้กรณีไหน
- Handler: แปลง HTTP request → usecase input, usecase output → HTTP response
- Middleware: ตรวจสอบ token, log request, จำกัด rate, เพิ่ม security headers
- DTO: กำหนดโครงสร้าง JSON request/response (แยกจาก entity model)
- Router: จับคู่ path กับ handler และ middleware

### ทำไมต้องใช้
- แยกการรับ/ส่งข้อมูลออกจาก business logic
- เปลี่ยนจาก REST เป็น GraphQL ได้โดยไม่ต้องแก้ usecase
- จัดการ cross-cutting concerns (logging, auth) เป็น centralized

### ประโยชน์ที่ได้รับ
- เปลี่ยนรูปแบบ API (REST → gRPC) โดยไม่กระทบ usecase
- ทดสอบ handler แบบ integration ได้ง่าย
- middleware reuse

### ข้อควรระวัง
- handler ควรสั้น (แค่ binding, validation, call usecase, response)
- อย่าใส่ business logic ใน handler
- DTO ควรแยกจาก entity model เพื่อป้องกันข้อมูล泄露 (password hash)

### ข้อดี
- แยกชัดเจน, ยืดหยุ่น, middleware จัดการ统一

### ข้อเสีย
- มีไฟล์จำนวนมาก (handler, dto, middleware แต่ละตัว)
- อาจมีการ mapping ซ้ำซ้อน (entity → dto)

### ข้อห้าม
- ห้ามเรียก repository โดยตรงจาก handler
- ห้ามทำ business logic (if-else ที่เกี่ยวกับธุรกิจ) ใน handler
- ห้ามใช้ entity model เป็น request DTO ถ้ามี field ที่ไม่ต้องการให้ client ส่งมา

---

## โค้ดที่รันได้จริง

### 1. DTOs (Data Transfer Objects) – `dto/`

#### `dto/auth_dto.go`

```go
// Package dto defines request and response structures for HTTP APIs.
// ----------------------------------------------------------------
// แพ็คเกจ dtd กำหนดโครงสร้าง request และ response สำหรับ HTTP APIs
package dto

// LoginRequest represents the login request body.
// ----------------------------------------------------------------
// LoginRequest แทนโครงสร้าง request สำหรับเข้าสู่ระบบ
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse represents the login response body.
// ----------------------------------------------------------------
// LoginResponse แทนโครงสร้าง response สำหรับเข้าสู่ระบบ
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"` // seconds, วินาที
}

// RefreshTokenRequest represents the refresh token request.
// ----------------------------------------------------------------
// RefreshTokenRequest แทนโครงสร้าง request สำหรับการขอ access token ใหม่
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse represents the refresh token response.
// ----------------------------------------------------------------
// RefreshTokenResponse แทนโครงสร้าง response สำหรับการขอ access token ใหม่
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

// RegisterRequest represents the registration request.
// ----------------------------------------------------------------
// RegisterRequest แทนโครงสร้าง request สำหรับการลงทะเบียน
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
}

// RegisterResponse represents the registration response.
// ----------------------------------------------------------------
// RegisterResponse แทนโครงสร้าง response สำหรับการลงทะเบียน
type RegisterResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

// ChangePasswordRequest represents the change password request.
// ----------------------------------------------------------------
// ChangePasswordRequest แทนโครงสร้าง request สำหรับเปลี่ยนรหัสผ่าน
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}
```

#### `dto/user_dto.go`

```go
package dto

// UserResponse represents user data returned to client.
// ----------------------------------------------------------------
// UserResponse แทนข้อมูลผู้ใช้ที่ส่งกลับไปยัง client
type UserResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// UpdateUserRequest represents request to update user profile.
// ----------------------------------------------------------------
// UpdateUserRequest แทนโครงสร้าง request สำหรับอัปเดตโปรไฟล์ผู้ใช้
type UpdateUserRequest struct {
	FullName string `json:"full_name" validate:"omitempty"`
	Email    string `json:"email" validate:"omitempty,email"`
}

// ListUsersRequest represents query parameters for user listing.
// ----------------------------------------------------------------
// ListUsersRequest แทน query parameters สำหรับรายการผู้ใช้
type ListUsersRequest struct {
	Page     int `json:"page" query:"page"`
	PageSize int `json:"page_size" query:"page_size"`
}

// ListUsersResponse represents paginated user list.
// ----------------------------------------------------------------
// ListUsersResponse แทนรายชื่อผู้ใช้แบบแบ่งหน้า
type ListUsersResponse struct {
	Data       []UserResponse `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// CreateUserRequest (admin) represents request to create user with role.
// ----------------------------------------------------------------
// CreateUserRequest (admin) แทนโครงสร้าง request สำหรับสร้างผู้ใช้พร้อม role
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	Role     string `json:"role" validate:"omitempty,oneof=user admin"`
}
```

#### `dto/error_dto.go`

```go
package dto

// ErrorResponse represents standard error response.
// ----------------------------------------------------------------
// ErrorResponse แทนโครงสร้าง error response มาตรฐาน
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
	TraceID string `json:"trace_id,omitempty"`
}

// ValidationErrorResponse represents field-level validation errors.
// ----------------------------------------------------------------
// ValidationErrorResponse แทน error การ validate ตาม field
type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Fields  map[string]string `json:"fields"`
	TraceID string            `json:"trace_id,omitempty"`
}
```

---

### 2. Middleware – `middleware/`

#### `middleware/auth.go`

```go
// Package middleware provides HTTP middleware functions.
// ----------------------------------------------------------------
// แพ็คเกจ middleware ให้ฟังก์ชัน middleware สำหรับ HTTP
package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"gobackend/internal/pkg/jwt"
	"gobackend/internal/pkg/logger"
	"go.uber.org/zap"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
	TraceIDKey contextKey = "trace_id"
)

// JWTAuth validates JWT token and injects user info into context.
// ----------------------------------------------------------------
// JWTAuth ตรวจสอบ JWT token และเพิ่มข้อมูลผู้ใช้ใน context
func JWTAuth(maker jwt.Maker, blacklistCache interface{ Has(string) bool }) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			// ดึง header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondWithError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			// Extract Bearer token
			// แยก token แบบ Bearer
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				respondWithError(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}
			tokenString := parts[1]

			// Verify token
			// ตรวจสอบ token
			payload, err := maker.VerifyToken(tokenString)
			if err != nil {
				if errors.Is(err, jwt.ErrExpiredToken) {
					respondWithError(w, http.StatusUnauthorized, "token expired")
				} else {
					respondWithError(w, http.StatusUnauthorized, "invalid token")
				}
				return
			}

			// Check if token is blacklisted (optional)
			// ตรวจสอบว่า token ถูก blacklist หรือไม่
			if blacklistCache != nil && blacklistCache.Has(payload.ID.String()) {
				respondWithError(w, http.StatusUnauthorized, "token revoked")
				return
			}

			// Inject user info into request context
			// เพิ่มข้อมูลผู้ใช้ใน context ของ request
			ctx := context.WithValue(r.Context(), UserIDKey, payload.UserID)
			ctx = context.WithValue(ctx, RoleKey, payload.Role)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole ensures user has at least one of the allowed roles.
// ----------------------------------------------------------------
// RequireRole ตรวจสอบว่าผู้ใช้มี role หนึ่งในที่อนุญาตหรือไม่
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(RoleKey).(string)
			if !ok {
				respondWithError(w, http.StatusForbidden, "forbidden")
				return
			}
			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}
			respondWithError(w, http.StatusForbidden, "insufficient permissions")
		})
	}
}

// GetUserIDFromContext extracts user ID from request context.
// ----------------------------------------------------------------
// GetUserIDFromContext ดึง user ID จาก context ของ request
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	id, ok := ctx.Value(UserIDKey).(uint)
	return id, ok
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	http.Error(w, message, code)
}
```

#### `middleware/cors.go`

```go
package middleware

import "net/http"

// CORS handles Cross-Origin Resource Sharing.
// ----------------------------------------------------------------
// CORS จัดการ Cross-Origin Resource Sharing
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			// Allow if origin is in allowed list or if in development
			// อนุญาตถ้า origin อยู่ในรายการที่อนุญาต หรืออยู่ในโหมดพัฒนา
			allowed := false
			for _, o := range allowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}
			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
```

#### `middleware/logger.go`

```go
package middleware

import (
	"net/http"
	"time"

	"gobackend/internal/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestLogger logs each HTTP request with method, path, status, duration.
// ----------------------------------------------------------------
// RequestLogger บันทึก HTTP request แต่ละรายการพร้อม method, path, status, ระยะเวลา
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		traceID := uuid.New().String()

		// Add trace ID to context
		// เพิ่ม trace ID ใน context
		ctx := context.WithValue(r.Context(), TraceIDKey, traceID)
		r = r.WithContext(ctx)

		// Wrap response writer to capture status code
		// ห่อ response writer เพื่อจับ status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			duration := time.Since(start)
			logger.Log.Info("HTTP request",
				zap.String("trace_id", traceID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", wrapped.statusCode),
				zap.Duration("duration", duration),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)
		}()

		next.ServeHTTP(wrapped, r)
	})
}

// responseWriter captures status code.
// ----------------------------------------------------------------
// responseWriter จับ status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
```

#### `middleware/rate_limit.go`

```go
package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"gobackend/internal/pkg/logger"
	"go.uber.org/zap"
)

// RateLimiter implements token bucket rate limiting per IP.
// ----------------------------------------------------------------
// RateLimiter จำกัดอัตราการเรียกใช้แบบ token bucket ต่อ IP
type RateLimiter struct {
	requestsPerSec int
	burst          int
	clients        map[string]*clientBucket
	mu             sync.RWMutex
}

type clientBucket struct {
	tokens     int
	lastRefill time.Time
}

// NewRateLimiter creates a rate limiter.
// ----------------------------------------------------------------
// NewRateLimiter สร้าง rate limiter ใหม่
func NewRateLimiter(requestsPerSec, burst int) *RateLimiter {
	rl := &RateLimiter{
		requestsPerSec: requestsPerSec,
		burst:          burst,
		clients:        make(map[string]*clientBucket),
	}
	// Cleanup old entries every minute
	go rl.cleanup()
	return rl
}

// Middleware returns HTTP handler middleware.
// ----------------------------------------------------------------
// Middleware คืน middleware สำหรับ HTTP
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getRealIP(r)
		if !rl.allow(ip) {
			logger.Log.Warn("rate limit exceeded", zap.String("ip", ip))
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	bucket, exists := rl.clients[ip]
	if !exists {
		rl.clients[ip] = &clientBucket{
			tokens:     rl.burst - 1,
			lastRefill: now,
		}
		return true
	}

	// Refill tokens based on time passed
	// เติม token ตามเวลาที่ผ่านไป
	elapsed := now.Sub(bucket.lastRefill).Seconds()
	newTokens := int(elapsed * float64(rl.requestsPerSec))
	if newTokens > 0 {
		bucket.tokens += newTokens
		if bucket.tokens > rl.burst {
			bucket.tokens = rl.burst
		}
		bucket.lastRefill = now
	}

	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}
	return false
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		for ip, bucket := range rl.clients {
			if time.Since(bucket.lastRefill) > 10*time.Minute {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func getRealIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	return r.RemoteAddr
}
```

#### `middleware/security.go`

```go
package middleware

import "net/http"

// SecurityHeaders adds security-related HTTP headers.
// ----------------------------------------------------------------
// SecurityHeaders เพิ่ม HTTP headers ที่เกี่ยวกับความปลอดภัย
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent XSS attacks
		// ป้องกัน XSS
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		// Prevent MIME type sniffing
		// ป้องกัน MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Prevent clickjacking
		// ป้องกัน clickjacking
		w.Header().Set("X-Frame-Options", "DENY")
		// Enforce HTTPS (HSTS)
		// บังคับใช้ HTTPS
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		// Content Security Policy (adjust as needed)
		// นโยบายความปลอดภัยของเนื้อหา
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		// Referrer policy
		// นโยบาย Referrer
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		next.ServeHTTP(w, r)
	})
}
```

#### `middleware/monitoring.go` (Prometheus metrics)

```go
package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

// Metrics records request count and duration for Prometheus.
// ----------------------------------------------------------------
// Metrics บันทึกจำนวน request และระยะเวลาสำหรับ Prometheus
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &statusWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		duration := time.Since(start).Seconds()
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(wrapped.statusCode)).Inc()
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}
```

---

### 3. Handlers – `handler/`

#### `handler/auth_handler.go`

```go
// Package handler contains HTTP request handlers.
// ----------------------------------------------------------------
// แพ็คเกจ handler บรรจุตัวจัดการ request HTTP
package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"gobackend/internal/delivery/rest/dto"
	"gobackend/internal/usecase"
	"github.com/go-playground/validator/v10"
)

// AuthHandler handles authentication endpoints.
// ----------------------------------------------------------------
// AuthHandler จัดการ endpoints สำหรับการรับรองตัวตน
type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	validate    *validator.Validate
}

// NewAuthHandler creates a new auth handler.
// ----------------------------------------------------------------
// NewAuthHandler สร้าง auth handler ใหม่
func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		validate:    validator.New(),
	}
}

// Register handles POST /register.
// ----------------------------------------------------------------
// Register จัดการ POST /register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		respondValidationError(w, err)
		return
	}

	user, err := h.authUsecase.Register(r.Context(), req.Email, req.Password, req.FullName)
	if err != nil {
		if errors.Is(err, usecase.ErrEmailAlreadyExists) {
			respondError(w, http.StatusConflict, "email already exists")
			return
		}
		respondError(w, http.StatusInternalServerError, "registration failed")
		return
	}

	resp := dto.RegisterResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
	respondJSON(w, http.StatusCreated, resp)
}

// Login handles POST /login.
// ----------------------------------------------------------------
// Login จัดการ POST /login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		respondValidationError(w, err)
		return
	}

	accessToken, refreshToken, err := h.authUsecase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) || errors.Is(err, usecase.ErrInvalidPassword) {
			respondError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		respondError(w, http.StatusInternalServerError, "login failed")
		return
	}

	resp := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    15 * 60, // 15 minutes
	}
	respondJSON(w, http.StatusOK, resp)
}

// Refresh handles POST /refresh.
// ----------------------------------------------------------------
// Refresh จัดการ POST /refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.RefreshToken == "" {
		respondError(w, http.StatusBadRequest, "refresh token required")
		return
	}

	newAccessToken, err := h.authUsecase.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidRefreshToken) {
			respondError(w, http.StatusUnauthorized, "invalid refresh token")
			return
		}
		respondError(w, http.StatusInternalServerError, "refresh failed")
		return
	}

	resp := dto.RefreshTokenResponse{
		AccessToken: newAccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   15 * 60,
	}
	respondJSON(w, http.StatusOK, resp)
}

// Logout handles POST /logout.
// ----------------------------------------------------------------
// Logout จัดการ POST /logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	// Get refresh token from request body (optional)
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	if err := h.authUsecase.Logout(r.Context(), userID, body.RefreshToken); err != nil {
		// Log but don't fail
	}
	w.WriteHeader(http.StatusNoContent)
}
```

#### `handler/user_handler.go`

```go
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"gobackend/internal/delivery/rest/dto"
	"gobackend/internal/delivery/rest/middleware"
	"gobackend/internal/usecase"
	"github.com/go-chi/chi/v5"
)

// UserHandler handles user management endpoints.
// ----------------------------------------------------------------
// UserHandler จัดการ endpoints สำหรับการจัดการผู้ใช้
type UserHandler struct {
	userUsecase usecase.UserUsecase
	validate    *validator.Validate
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		validate:    validator.New(),
	}
}

// GetProfile handles GET /profile.
// ----------------------------------------------------------------
// GetProfile จัดการ GET /profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	user, err := h.userUsecase.GetUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get user")
		return
	}
	respondJSON(w, http.StatusOK, toUserResponse(user))
}

// UpdateProfile handles PUT /profile.
// ----------------------------------------------------------------
// UpdateProfile จัดการ PUT /profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}
	// Get existing user
	user, err := h.userUsecase.GetUserByID(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if err := h.userUsecase.UpdateUser(r.Context(), user); err != nil {
		respondError(w, http.StatusInternalServerError, "update failed")
		return
	}
	respondJSON(w, http.StatusOK, toUserResponse(user))
}

// ListUsers handles GET /admin/users (admin only).
// ----------------------------------------------------------------
// ListUsers จัดการ GET /admin/users (เฉพาะ admin)
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	users, total, err := h.userUsecase.ListUsers(r.Context(), page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list users")
		return
	}
	data := make([]dto.UserResponse, len(users))
	for i, u := range users {
		data[i] = toUserResponse(&u)
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	resp := dto.ListUsersResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
	respondJSON(w, http.StatusOK, resp)
}

// CreateUserByAdmin handles POST /admin/users (admin only).
// ----------------------------------------------------------------
// CreateUserByAdmin จัดการ POST /admin/users (เฉพาะ admin)
func (h *UserHandler) CreateUserByAdmin(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		respondValidationError(w, err)
		return
	}
	role := req.Role
	if role == "" {
		role = "user"
	}
	user, err := h.userUsecase.CreateUserAsAdmin(r.Context(), req.Email, req.Password, req.FullName, role)
	if err != nil {
		if errors.Is(err, usecase.ErrEmailAlreadyExists) {
			respondError(w, http.StatusConflict, "email already exists")
			return
		}
		respondError(w, http.StatusInternalServerError, "create failed")
		return
	}
	respondJSON(w, http.StatusCreated, toUserResponse(user))
}

func toUserResponse(u *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FullName:  u.FullName,
		Role:      string(u.Role),
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
```

#### `handler/health_handler.go`

```go
package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"
	"context"
)

// HealthHandler provides health check endpoints.
// ----------------------------------------------------------------
// HealthHandler ให้บริการ endpoints สำหรับตรวจสอบสุขภาพ
type HealthHandler struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewHealthHandler(db *gorm.DB, rdb *redis.Client) *HealthHandler {
	return &HealthHandler{db: db, rdb: rdb}
}

// LivenessProbe returns 200 if process is alive.
// ----------------------------------------------------------------
// LivenessProbe คืน 200 ถ้า process ยังทำงาน
func (h *HealthHandler) LivenessProbe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// ReadinessProbe checks dependencies (DB, Redis).
// ----------------------------------------------------------------
// ReadinessProbe ตรวจสอบ dependencies (DB, Redis)
func (h *HealthHandler) ReadinessProbe(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	status := map[string]string{
		"postgres": "ok",
		"redis":    "ok",
	}
	overall := true

	// Check PostgreSQL
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.PingContext(ctx) != nil {
		status["postgres"] = "unavailable"
		overall = false
	}
	// Check Redis
	if err := h.rdb.Ping(ctx).Err(); err != nil {
		status["redis"] = "unavailable"
		overall = false
	}

	w.Header().Set("Content-Type", "application/json")
	if !overall {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	json.NewEncoder(w).Encode(status)
}
```

#### Helper functions for responses

```go
// respondJSON writes JSON response.
// ----------------------------------------------------------------
// respondJSON เขียน response แบบ JSON
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError writes error response.
// ----------------------------------------------------------------
// respondError เขียน error response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, dto.ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
		Code:    status,
	})
}

// respondValidationError writes validation errors.
// ----------------------------------------------------------------
// respondValidationError เขียน validation error
func respondValidationError(w http.ResponseWriter, err error) {
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		respondError(w, http.StatusBadRequest, "validation failed")
		return
	}
	fields := make(map[string]string)
	for _, e := range ve {
		fields[e.Field()] = e.Tag()
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(dto.ValidationErrorResponse{
		Error:  "validation failed",
		Fields: fields,
	})
}
```

---

### 4. Router – `router.go`

```go
// Package rest sets up HTTP routes and middleware.
// ----------------------------------------------------------------
// แพ็คเกจ rest ตั้งค่า routes และ middleware HTTP
package rest

import (
	"net/http"

	"gobackend/internal/delivery/rest/handler"
	"gobackend/internal/delivery/rest/middleware"
	"gobackend/internal/pkg/jwt"
	"gobackend/internal/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware as chimw"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// RouterConfig holds dependencies for router setup.
// ----------------------------------------------------------------
// RouterConfig เก็บ dependencies สำหรับการตั้งค่า router
type RouterConfig struct {
	AuthHandler    *handler.AuthHandler
	UserHandler    *handler.UserHandler
	HealthHandler  *handler.HealthHandler
	JWTMaker       jwt.Maker
	RateLimiter    *middleware.RateLimiter
	AllowedOrigins []string
}

// NewRouter creates and configures the main router.
// ----------------------------------------------------------------
// NewRouter สร้างและกำหนดค่า router หลัก
func NewRouter(cfg RouterConfig) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	// middleware ระดับ global
	r.Use(chimw.Recoverer)                         // panic recovery
	r.Use(middleware.RequestLogger)                // request logging
	r.Use(middleware.SecurityHeaders)              // security headers
	r.Use(middleware.CORS(cfg.AllowedOrigins))     // CORS
	r.Use(middleware.Metrics)                      // Prometheus metrics
	if cfg.RateLimiter != nil {
		r.Use(cfg.RateLimiter.Middleware)          // rate limiting
	}

	// Health endpoints (no auth)
	// endpoints สำหรับตรวจสอบสุขภาพ (ไม่ต้อง auth)
	r.Get("/health/live", cfg.HealthHandler.LivenessProbe)
	r.Get("/health/ready", cfg.HealthHandler.ReadinessProbe)

	// Metrics endpoint for Prometheus
	// endpoint Metrics สำหรับ Prometheus
	r.Handle("/metrics", promhttp.Handler())

	// Public routes (no auth)
	// routes สาธารณะ (ไม่ต้อง auth)
	r.Post("/register", cfg.AuthHandler.Register)
	r.Post("/login", cfg.AuthHandler.Login)
	r.Post("/refresh", cfg.AuthHandler.Refresh)

	// Protected routes (require authentication)
	// routes ที่ต้องการการรับรองตัวตน
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTMaker, nil)) // nil blacklist for now
		r.Post("/logout", cfg.AuthHandler.Logout)
		r.Get("/profile", cfg.UserHandler.GetProfile)
		r.Put("/profile", cfg.UserHandler.UpdateProfile)
	})

	// Admin routes (require role admin)
	// routes สำหรับ admin (ต้องการ role admin)
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTMaker, nil))
		r.Use(middleware.RequireRole("admin"))
		r.Get("/admin/users", cfg.UserHandler.ListUsers)
		r.Post("/admin/users", cfg.UserHandler.CreateUserByAdmin)
	})

	return r
}
```

---

## วิธีใช้งาน module นี้

### ใน `main.go` หรือ `cmd/api/main.go`

```go
func main() {
    // Load config, init DB, Redis, JWT maker, rate limiter, usecases...
    // ...
    // Create handlers
    authHandler := handler.NewAuthHandler(authUsecase)
    userHandler := handler.NewUserHandler(userUsecase)
    healthHandler := handler.NewHealthHandler(db, redisClient)
    rateLimiter := middleware.NewRateLimiter(cfg.RateLimit.RequestsPerSecond, cfg.RateLimit.Burst)

    routerConfig := rest.RouterConfig{
        AuthHandler:    authHandler,
        UserHandler:    userHandler,
        HealthHandler:  healthHandler,
        JWTMaker:       jwtMaker,
        RateLimiter:    rateLimiter,
        AllowedOrigins: []string{"http://localhost:3000", "https://yourdomain.com"},
    }
    r := rest.NewRouter(routerConfig)

    srv := &http.Server{Addr: ":8080", Handler: r}
    log.Fatal(srv.ListenAndServe())
}
```

---

## ตารางสรุป Delivery Components

| Component | หน้าที่ | ตัวอย่าง |
|-----------|--------|----------|
| Handler | แปลง HTTP → usecase → HTTP | `AuthHandler.Login` |
| Middleware | Pre/post processing | `JWTAuth`, `RateLimiter`, `RequestLogger` |
| DTO | โครงสร้าง JSON request/response | `LoginRequest`, `UserResponse` |
| Router | จับคู่ route กับ handler + middleware | `chi.Mux` |

---

## แบบฝึกหัดท้าย module (3 ข้อ)

1. เพิ่ม endpoint `DELETE /admin/users/{id}` ใน `UserHandler` และ route สำหรับ admin เท่านั้น พร้อมเรียก usecase `DeleteUser`
2. สร้าง middleware `RequestID` ที่สร้าง UUID และเพิ่มใน response header `X-Request-Id`
3. ปรับปรุง `JWTAuth` middleware ให้ดึง public key แบบ dynamic จาก Redis (rotation) โดยไม่ต้อง restart server

---

## แหล่งอ้างอิง

- [Chi Router documentation](https://go-chi.io/)
- [Gorilla WebSocket (for real-time)](https://github.com/gorilla/websocket)
- [Go validation package](https://github.com/go-playground/validator)
- [Prometheus Go client](https://github.com/prometheus/client_golang)

---

**หมายเหตุ:** module นี้เป็นส่วนสุดท้ายของ Delivery layer ที่สมบูรณ์สำหรับระบบ gobackend หากต้องการ module เพิ่มเติม (เช่น worker, pkg components) โปรดแจ้งคำว่า "ต่อไป" หรือระบุชื่อ module