# 📚 PART 8, 12, 13, 14 — INTERFACE · TESTING · DEVOPS · SECURITY

> **ขนาด**: ใหญ่มาก — 4 PART ใหญ่
> **PART 8**: Interface Layer Deep Dive
> **PART 12**: Testing Strategy
> **PART 13**: Deployment & DevOps
> **PART 14**: Security Deep Dive
> **เป้าหมาย**: ครอบคลุมตั้งแต่วิธีรับ request ไปจนถึง deploy production อย่างปลอดภัย

---
---

# 🌐 PART 8 — INTERFACE LAYER DEEP DIVE

> **ขนาด**: ใหญ่ — แยก 9 ตอนย่อย
> **Part 8A**: Interface Layer Architecture
> **Part 8B**: HTTP Handlers & Routes
> **Part 8C**: Middleware Stack
> **Part 8D**: Request/Response Patterns
> **Part 8E**: WebSocket Hub
> **Part 8F**: MQTT Gateway
> **Part 8G**: CLI & Scheduler Entry Points
> **Part 8H**: API Versioning & Deprecation
> **Part 8I**: Error Handling & Observability

---

## 🅰️ PART 8A — INTERFACE LAYER ARCHITECTURE

### A.1 Position

```
┌─────────────────────────────────────────────────────────────┐
│  ★★★ INTERFACE LAYER ★★★  ← PART 8 นี้                     │
│  HTTP (Gin) · WebSocket · MQTT Gateway · CLI · Scheduler    │
│  ────────────────────────────────────────────────────────   │
│  ✓ รับ input จากภายนอก → map เป็น Use Case Input            │
│  ✓ Map Use Case Output → HTTP response                       │
│  ✓ Middleware: auth, tenant, rate-limit, CORS, logging       │
│  ✗ ไม่มี business logic                                     │
│  ✗ ไม่เรียก repository โดยตรง (ต้องผ่าน use case)           │
├─────────────────────────────────────────────────────────────┤
│                   Application Layer                         │
│  Use Cases · DTOs · Ports                                   │
└─────────────────────────────────────────────────────────────┘
```

### A.2 Sub-Folders

```
internal/modules/{{module}}/interfaces/
├── http/
│   ├── {{aggregate}}_handler.go
│   ├── routes.go
│   └── dto.go
├── websocket/
│   ├── hub.go
│   ├── client.go
│   └── handler.go
├── mqtt/
│   ├── gateway.go
│   └── topics.go
└── middleware/
    └── {{module}}_middleware.go
```

---

## 🅱️ PART 8B — HTTP HANDLERS & ROUTES

### B.1 Handler Template

```go
// internal/modules/device/interfaces/http/device_handler.go
package http

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    "icmongolang/internal/modules/device/application"
    apperrors "icmongolang/internal/shared/application"
)

type DeviceHandler struct {
    registerUC   *application.RegisterDeviceUseCase
    updateUC     *application.UpdateDeviceUseCase
    deleteUC     *application.DeleteDeviceUseCase
    provisionUC  *application.ProvisionDeviceUseCase
    sendCmdUC    *application.SendCommandUseCase
    getUC        *application.GetDeviceUseCase
    listUC       *application.ListDevicesUseCase
    telemetryUC  *application.IngestTelemetryUseCase
    logger       application.Logger
}

func NewDeviceHandler(
    registerUC *application.RegisterDeviceUseCase,
    updateUC *application.UpdateDeviceUseCase,
    deleteUC *application.DeleteDeviceUseCase,
    provisionUC *application.ProvisionDeviceUseCase,
    sendCmdUC *application.SendCommandUseCase,
    getUC *application.GetDeviceUseCase,
    listUC *application.ListDevicesUseCase,
    telemetryUC *application.IngestTelemetryUseCase,
    logger application.Logger,
) *DeviceHandler {
    return &DeviceHandler{
        registerUC: registerUC, updateUC: updateUC, deleteUC: deleteUC,
        provisionUC: provisionUC, sendCmdUC: sendCmdUC,
        getUC: getUC, listUC: listUC, telemetryUC: telemetryUC,
        logger: logger,
    }
}

// ============================================================
// POST /api/v1/devices
// ============================================================
func (h *DeviceHandler) Register(c *gin.Context) {
    var req RegisterDeviceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.badRequest(c, "VALIDATION_ERROR", err.Error())
        return
    }

    tenantID := c.MustGet("tenant_id").(uuid.UUID)
    userID := c.MustGet("user_id").(uuid.UUID)
    idemKey := c.GetHeader("Idempotency-Key")

    customerID, _ := uuid.Parse(req.CustomerID)
    siteID, _ := uuid.Parse(req.SiteID)
    var zoneID, modelID uuid.UUID
    if req.ZoneID != "" {
        zoneID, _ = uuid.Parse(req.ZoneID)
    }
    if req.ModelID != "" {
        modelID, _ = uuid.Parse(req.ModelID)
    }

    out, err := h.registerUC.Execute(c.Request.Context(), application.RegisterDeviceInput{
        TenantID:       tenantID,
        ActorID:        userID,
        CustomerID:     customerID,
        SiteID:         siteID,
        ZoneID:         zoneID,
        SerialNo:       req.SerialNo,
        Name:           req.Name,
        Type:           req.Type,
        Protocol:       req.Protocol,
        ModelID:        modelID,
        Tags:           req.Tags,
        IdempotencyKey: idemKey,
    })
    if err != nil {
        h.respondError(c, err)
        return
    }

    c.Header("Location", "/api/v1/devices/"+out.ID.String())
    c.JSON(http.StatusCreated, gin.H{"data": out})
}

// ============================================================
// GET /api/v1/devices
// ============================================================
func (h *DeviceHandler) List(c *gin.Context) {
    tenantID := c.MustGet("tenant_id").(uuid.UUID)

    var req ListDevicesQuery
    if err := c.ShouldBindQuery(&req); err != nil {
        h.badRequest(c, "VALIDATION_ERROR", err.Error())
        return
    }

    // Parse optional filters
    var siteID *uuid.UUID
    if req.SiteID != "" {
        if id, err := uuid.Parse(req.SiteID); err == nil {
            siteID = &id
        }
    }

    out, err := h.listUC.Execute(c.Request.Context(), application.ListDevicesInput{
        TenantID: tenantID,
        SiteID:   siteID,
        Status:   req.Status,
        Type:     req.Type,
        Search:   req.Search,
        Page:     req.Page,
        PageSize: req.PageSize,
    })
    if err != nil {
        h.respondError(c, err)
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": out})
}

// ============================================================
// GET /api/v1/devices/:id
// ============================================================
func (h *DeviceHandler) Get(c *gin.Context) {
    tenantID := c.MustGet("tenant_id").(uuid.UUID)
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        h.badRequest(c, "INVALID_ID", "invalid device ID")
        return
    }
    out, err := h.getUC.Execute(c.Request.Context(), application.GetDeviceInput{
        TenantID: tenantID,
        DeviceID: id,
    })
    if err != nil {
        h.respondError(c, err)
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": out})
}

// ============================================================
// PUT /api/v1/devices/:id
// ============================================================
func (h *DeviceHandler) Update(c *gin.Context) {
    tenantID := c.MustGet("tenant_id").(uuid.UUID)
    userID := c.MustGet("user_id").(uuid.UUID)
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        h.badRequest(c, "INVALID_ID", "invalid device ID")
        return
    }
    var req UpdateDeviceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.badRequest(c, "VALIDATION_ERROR", err.Error())
        return
    }
    out, err := h.updateUC.Execute(c.Request.Context(), application.UpdateDeviceInput{
        TenantID: tenantID,
        DeviceID: id,
        Name:     req.Name,
        Tags:     req.Tags,
        ActorID:  userID,
    })
    if err != nil {
        h.respondError(c, err)
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": out})
}

// ============================================================
// DELETE /api/v1/devices/:id
// ============================================================
func (h *DeviceHandler) Delete(c *gin.Context) {
    tenantID := c.MustGet("tenant_id").(uuid.UUID)
    userID := c.MustGet("user_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    if err := h.deleteUC.Execute(c.Request.Context(), application.DeleteDeviceInput{
        TenantID: tenantID,
        DeviceID: id,
        ActorID:  userID,
    }); err != nil {
        h.respondError(c, err)
        return
    }
    c.Status(http.StatusNoContent)
}

// ============================================================
// POST /api/v1/devices/:id/provision
// ============================================================
func (h *DeviceHandler) Provision(c *gin.Context) {
    tenantID := c.MustGet("tenant_id").(uuid.UUID)
    userID := c.MustGet("user_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    out, err := h.provisionUC.Execute(c.Request.Context(), application.ProvisionDeviceInput{
        TenantID: tenantID,
        DeviceID: id,
        ActorID:  userID,
    })
    if err != nil {
        h.respondError(c, err)
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": out})
}

// ============================================================
// POST /api/v1/devices/:id/commands
// ============================================================
func (h *DeviceHandler) SendCommand(c *gin.Context) {
    tenantID := c.MustGet("tenant_id").(uuid.UUID)
    userID := c.MustGet("user_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    var req SendCommandRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.badRequest(c, "VALIDATION_ERROR", err.Error())
        return
    }

    out, err := h.sendCmdUC.Execute(c.Request.Context(), application.SendCommandInput{
        TenantID:     tenantID,
        DeviceID:     id,
        Command:      req.Command,
        Payload:      req.Payload,
        Priority:     req.Priority,
        ExpiresInSec: req.ExpiresInSec,
        IssuedBy:     userID,
    })
    if err != nil {
        h.respondError(c, err)
        return
    }
    c.JSON(http.StatusAccepted, gin.H{"data": out})
}

// ============================================================
// Helpers
// ============================================================
func (h *DeviceHandler) badRequest(c *gin.Context, code, msg string) {
    c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
        "error": gin.H{"code": code, "message": msg},
    })
}

func (h *DeviceHandler) respondError(c *gin.Context, err error) {
    status := apperrors.HTTPStatus(err)
    resp := apperrors.ToResponse(err)
    resp.Error.RequestID = c.GetString("request_id")
    c.JSON(status, resp)
}
```

### B.2 Routes Registration

```go
// internal/modules/device/interfaces/http/routes.go
package http

import (
    "time"

    "github.com/gin-gonic/gin"
    sharedmw "icmongolang/internal/shared/interfaces/http/middleware"
)

type Handlers struct {
    Device    *DeviceHandler
    Shadow    *ShadowHandler
    Telemetry *TelemetryHandler
    AlertRule *AlertRuleHandler
}

type RoutesConfig struct {
    Auth       gin.HandlerFunc
    Tenant     gin.HandlerFunc
    RateLimit  gin.HandlerFunc
    RequireRole func(roles ...string) gin.HandlerFunc
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handlers, cfg RoutesConfig) {
    g := rg.Group("/devices")
    g.Use(cfg.Auth, cfg.Tenant)

    // CRUD
    g.POST("", cfg.RequireRole("ADMIN", "MANAGER"), h.Device.Register)
    g.GET("", h.Device.List)
    g.GET("/:id", h.Device.Get)
    g.PUT("/:id", cfg.RequireRole("ADMIN", "MANAGER"), h.Device.Update)
    g.DELETE("/:id", cfg.RequireRole("ADMIN"), h.Device.Delete)

    // Provisioning
    g.POST("/:id/provision", cfg.RequireRole("ADMIN", "MANAGER"), h.Device.Provision)

    // Commands
    g.POST("/:id/commands", cfg.RateLimit, h.Device.SendCommand)
    g.GET("/:id/commands", h.Device.ListCommands)

    // Telemetry
    g.POST("/:id/telemetry", h.Telemetry.Ingest)
    g.GET("/:id/telemetry", h.Telemetry.Query)
    g.GET("/:id/telemetry/latest", h.Telemetry.Latest)

    // Shadow
    g.GET("/:id/shadow", h.Shadow.Get)
    g.PATCH("/:id/shadow/desired", h.Shadow.UpdateDesired)
    g.DELETE("/:id/shadow/desired", h.Shadow.ClearDesired)

    // Alert Rules
    ar := g.Group("/alert-rules")
    ar.GET("", h.AlertRule.List)
    ar.POST("", cfg.RequireRole("ADMIN", "MANAGER"), h.AlertRule.Create)
    ar.PUT("/:id", cfg.RequireRole("ADMIN", "MANAGER"), h.AlertRule.Update)
    ar.DELETE("/:id", cfg.RequireRole("ADMIN"), h.AlertRule.Delete)

    // Alert Events
    ae := g.Group("/alert-events")
    ae.GET("", h.AlertRule.ListEvents)
    ae.POST("/:id/ack", h.AlertRule.Acknowledge)
}
```

### B.3 Main Router Setup

```go
// internal/shared/interfaces/http/router.go
package http

func NewRouter(
    cfg Config,
    logger application.Logger,
    authSvc middleware.AuthService,
    deviceHandlers *devicehttp.Handlers,
    customerHandlers *customerhttp.Handlers,
    packageHandlers *packagehttp.Handlers,
    erpHandlers *erphttp.Handlers,
    logisticsHandlers *logisticshttp.Handlers,
    reportHandlers *reporthttp.Handlers,
) *gin.Engine {
    if cfg.Env == "production" {
        gin.SetMode(gin.ReleaseMode)
    }
    r := gin.New()

    // Global middleware
    r.Use(
        middleware.Recovery(logger),
        middleware.RequestID(),
        middleware.Logger(logger),
        middleware.CORS(cfg.CORS),
        middleware.SecurityHeaders(),
        middleware.PrometheusMetrics(),
    )

    // Health endpoints (no auth)
    r.GET("/health", healthHandler)
    r.GET("/readyz", readyzHandler(cfg))
    r.GET("/metrics", gin.WrapH(promhttp.Handler()))

    // API v1
    v1 := r.Group("/api/v1")

    // Public routes
    auth := v1.Group("/auth")
    {
        auth.POST("/login", authHandler.Login)
        auth.POST("/refresh", authHandler.Refresh)
    }

    // Protected routes
    protected := v1.Group("")
    protected.Use(
        middleware.Auth(authSvc),
        middleware.Tenant(),
        middleware.RateLimitGlobal(cfg.RateLimit),
    )

    // Register module routes
    devicehttp.RegisterRoutes(protected, deviceHandlers, devicehttp.RoutesConfig{
        Auth:        middleware.Auth(authSvc),
        Tenant:      middleware.Tenant(),
        RateLimit:   middleware.RateLimitPerTenant(100, time.Minute),
        RequireRole: middleware.RequireRole,
    })
    customerhttp.RegisterRoutes(protected, customerHandlers, ...)
    packagehttp.RegisterRoutes(protected, packageHandlers, ...)
    erphttp.RegisterRoutes(protected, erpHandlers, ...)
    logisticshttp.RegisterRoutes(protected, logisticsHandlers, ...)
    reporthttp.RegisterRoutes(protected, reportHandlers, ...)

    return r
}
```

### B.4 DTOs (HTTP Layer)

```go
// internal/modules/device/interfaces/http/dto.go
package http

// ============================================================
// REQUEST DTOs
// ============================================================
type RegisterDeviceRequest struct {
    CustomerID string   `json:"customer_id" binding:"required,uuid4"`
    SiteID     string   `json:"site_id"     binding:"required,uuid4"`
    ZoneID     string   `json:"zone_id"     binding:"omitempty,uuid4"`
    ModelID    string   `json:"model_id"    binding:"omitempty,uuid4"`
    SerialNo   string   `json:"serial_no"   binding:"required,min=3,max=100"`
    Name       string   `json:"name"        binding:"required,min=1,max=255"`
    Type       string   `json:"type"        binding:"required,oneof=SENSOR ACTUATOR GATEWAY CAMERA METER TRACKER"`
    Protocol   string   `json:"protocol"    binding:"required,oneof=MQTT MODBUS LORAWAN ZIGBEE HTTP COAP"`
    Tags       []string `json:"tags"        binding:"omitempty,max=10,dive,min=1,max=30"`
}

type UpdateDeviceRequest struct {
    Name string   `json:"name" binding:"omitempty,min=1,max=255"`
    Tags []string `json:"tags" binding:"omitempty,max=10,dive,min=1,max=30"`
}

type SendCommandRequest struct {
    Command      string                 `json:"command"        binding:"required,min=1,max=50"`
    Payload      map[string]interface{} `json:"payload"`
    Priority     int                    `json:"priority"       binding:"omitempty,min=1,max=5"`
    ExpiresInSec int                    `json:"expires_in_sec" binding:"omitempty,min=1,max=3600"`
}

type ListDevicesQuery struct {
    SiteID   string `form:"site_id"   binding:"omitempty,uuid4"`
    Status   *string `form:"status"    binding:"omitempty,oneof=ONLINE OFFLINE FAULT MAINTENANCE PROVISIONED"`
    Type     *string `form:"type"      binding:"omitempty,oneof=SENSOR ACTUATOR GATEWAY CAMERA METER TRACKER"`
    Search   string `form:"search"    binding:"omitempty,max=100"`
    Page     int    `form:"page"      binding:"omitempty,min=1"`
    PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
}

type UpdateShadowDesiredRequest struct {
    Desired map[string]interface{} `json:"desired" binding:"required"`
}

// ============================================================
// RESPONSE DTOs
// ============================================================
type DeviceResponse struct {
    ID             string               `json:"id"`
    TenantID       string               `json:"tenant_id"`
    CustomerID     string               `json:"customer_id"`
    SiteID         string               `json:"site_id"`
    ZoneID         string               `json:"zone_id,omitempty"`
    SerialNo       string               `json:"serial_no"`
    Name           string               `json:"name"`
    Type           string               `json:"type"`
    Protocol       string               `json:"protocol"`
    Status         string               `json:"status"`
    Firmware       string               `json:"firmware,omitempty"`
    LastSeenAt     *time.Time           `json:"last_seen_at,omitempty"`
    ProvisionedAt  *time.Time           `json:"provisioned_at,omitempty"`
    Tags           []string             `json:"tags,omitempty"`
    CreatedAt      time.Time            `json:"created_at"`
    UpdatedAt      time.Time            `json:"updated_at"`
}
```

---

## 🅲 PART 8C — MIDDLEWARE STACK

### C.1 Middleware Order (สำคัญมาก)

```
1. Recovery          ← ป้องกัน panic → 500
2. RequestID         ← inject X-Request-ID
3. Logger            ← log request/response
4. CORS              ← จัดการ cross-origin
5. SecurityHeaders   ← X-Frame-Options, HSTS, ...
6. PrometheusMetrics ← track metrics
7. RateLimitGlobal   ← per-IP / per-tenant
8. Auth              ← JWT validation
9. Tenant            ← resolve tenant จาก header/token
10. RequireRole      ← RBAC
11. RateLimitPerRoute ← per-endpoint
12. Handler          ← business logic
```

### C.2 RequestID Middleware

```go
// internal/shared/interfaces/http/middleware/request_id.go
package middleware

import "github.com/gin-gonic/gin"

const (
    HeaderRequestID = "X-Request-ID"
    ContextRequestID = "request_id"
)

func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        rid := c.GetHeader(HeaderRequestID)
        if rid == "" {
            rid = uuid.NewString()
        }
        c.Set(ContextRequestID, rid)
        c.Header(HeaderRequestID, rid)
        c.Next()
    }
}
```

### C.3 Logger Middleware

```go
// internal/shared/interfaces/http/middleware/logger.go
package middleware

import (
    "time"
    "github.com/gin-gonic/gin"
    "icmongolang/internal/shared/application"
)

func Logger(logger application.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        query := c.Request.URL.RawQuery

        c.Next()

        latency := time.Since(start)
        status := c.Writer.Status()
        tenantID := c.GetString("tenant_id")
        userID := c.GetString("user_id")
        requestID := c.GetString("request_id")

        fields := []application.Field{
            {Key: "method", Value: c.Request.Method},
            {Key: "path", Value: path},
            {Key: "query", Value: query},
            {Key: "status", Value: status},
            {Key: "latency_ms", Value: latency.Milliseconds()},
            {Key: "ip", Value: c.ClientIP()},
            {Key: "request_id", Value: requestID},
            {Key: "tenant_id", Value: tenantID},
            {Key: "user_id", Value: userID},
        }

        if status >= 500 {
            logger.Error("request failed", fields...)
        } else if status >= 400 {
            logger.Warn("request error", fields...)
        } else {
            logger.Info("request", fields...)
        }
    }
}
```

### C.4 CORS Middleware

```go
// internal/shared/interfaces/http/middleware/cors.go
package middleware

func CORS(cfg CORSConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.GetHeader("Origin")
        if origin == "" || isAllowedOrigin(origin, cfg.AllowedOrigins) {
            c.Header("Access-Control-Allow-Origin", origin)
            c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
            c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Tenant-ID, X-Request-ID, Idempotency-Key")
            c.Header("Access-Control-Expose-Headers", "X-Request-ID, X-RateLimit-Limit, X-RateLimit-Remaining")
            c.Header("Access-Control-Allow-Credentials", "true")
            c.Header("Access-Control-Max-Age", "86400")
        }
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        c.Next()
    }
}

func isAllowedOrigin(origin string, allowed []string) bool {
    for _, a := range allowed {
        if a == "*" || a == origin {
            return true
        }
    }
    return false
}
```

### C.5 Auth Middleware

```go
// internal/shared/interfaces/http/middleware/auth.go
package middleware

type AuthService interface {
    ValidateToken(ctx context.Context, token string) (*AuthClaims, error)
}

type AuthClaims struct {
    UserID   uuid.UUID
    TenantID uuid.UUID
    Email    string
    Role     string
    Scopes   []string
    ExpiresAt time.Time
}

func Auth(svc AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        header := c.GetHeader("Authorization")
        if header == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{"code": "MISSING_TOKEN", "message": "authorization header required"},
            })
            return
        }
        if !strings.HasPrefix(header, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{"code": "INVALID_TOKEN_FORMAT", "message": "bearer token required"},
            })
            return
        }
        token := strings.TrimPrefix(header, "Bearer ")
        claims, err := svc.ValidateToken(c.Request.Context(), token)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{"code": "INVALID_TOKEN", "message": err.Error()},
            })
            return
        }
        c.Set("user_id", claims.UserID)
        c.Set("tenant_id", claims.TenantID)
        c.Set("user_email", claims.Email)
        c.Set("user_role", claims.Role)
        c.Set("scopes", claims.Scopes)
        c.Next()
    }
}
```

### C.6 Tenant Middleware

```go
// internal/shared/interfaces/http/middleware/tenant.go
package middleware

func Tenant() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Priority: header > JWT claim
        headerTenant := c.GetHeader("X-Tenant-ID")
        claimTenant := c.GetString("tenant_id")

        if headerTenant == "" && claimTenant == "" {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
                "error": gin.H{"code": "TENANT_REQUIRED", "message": "tenant context required"},
            })
            return
        }

        // ตรวจสอบว่า tenant header ตรงกับ JWT
        if headerTenant != "" && claimTenant != "" && headerTenant != claimTenant {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "error": gin.H{"code": "TENANT_MISMATCH", "message": "tenant header does not match token"},
            })
            return
        }

        tid := headerTenant
        if tid == "" {
            tid = claimTenant
        }
        tenantUUID, err := uuid.Parse(tid)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
                "error": gin.H{"code": "INVALID_TENANT", "message": "invalid tenant ID"},
            })
            return
        }
        c.Set("tenant_id", tenantUUID)
        c.Next()
    }
}
```

### C.7 RequireRole Middleware

```go
// internal/shared/interfaces/http/middleware/role.go
package middleware

func RequireRole(roles ...string) gin.HandlerFunc {
    allowed := make(map[string]bool, len(roles))
    for _, r := range roles {
        allowed[r] = true
    }
    return func(c *gin.Context) {
        role := c.GetString("user_role")
        if !allowed[role] {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "error": gin.H{
                    "code":    "FORBIDDEN",
                    "message": "insufficient permissions",
                    "details": gin.H{"required": roles, "current": role},
                },
            })
            return
        }
        c.Next()
    }
}
```

### C.8 Rate Limit Middleware (per route)

```go
// internal/shared/interfaces/http/middleware/rate_limit.go
package middleware

import (
    "fmt"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "icmongolang/internal/shared/infrastructure/redis"
)

func RateLimitPerTenant(limiter *redis.RateLimiter, limit int, window time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        tenantID := c.GetString("tenant_id")
        key := fmt.Sprintf("rl:tenant:%s:%s", tenantID, c.FullPath())

        result, err := limiter.Allow(c.Request.Context(), key, limit, window)
        if err != nil {
            c.Next() // fail-open
            return
        }

        c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", result.Limit))
        c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", result.Remaining))
        c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(result.ResetAfter).Unix()))

        if !result.Allowed {
            c.Header("Retry-After", fmt.Sprintf("%d", int(result.ResetAfter.Seconds())))
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "error": gin.H{
                    "code":    "RATE_LIMIT_EXCEEDED",
                    "message": "too many requests",
                },
            })
            return
        }
        c.Next()
    }
}
```

### C.9 Recovery Middleware

```go
// internal/shared/interfaces/http/middleware/recovery.go
package middleware

func Recovery(logger application.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if r := recover(); r != nil {
                stack := debug.Stack()
                requestID := c.GetString("request_id")
                logger.Error("panic recovered",
                    application.Field{Key: "request_id", Value: requestID},
                    application.Field{Key: "path", Value: c.Request.URL.Path},
                    application.Field{Key: "error", Value: fmt.Sprintf("%v", r)},
                    application.Field{Key: "stack", Value: string(stack)},
                )
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                    "error": gin.H{
                        "code":       "INTERNAL_ERROR",
                        "message":    "an unexpected error occurred",
                        "request_id": requestID,
                    },
                })
            }
        }()
        c.Next()
    }
}
```

### C.10 Security Headers

```go
// internal/shared/interfaces/http/middleware/security_headers.go
package middleware

func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Header("Content-Security-Policy", "default-src 'self'")
        c.Header("Permissions-Policy", "geolocation=(self), microphone=()")
        c.Next()
    }
}
```

---

## 🅳 PART 8D — REQUEST/RESPONSE PATTERNS

### D.1 Standard Response Format

```json
// Success (single)
{
    "data": { ... }
}

// Success (list)
{
    "data": {
        "items": [ ... ],
        "pagination": {
            "page": 1,
            "page_size": 20,
            "total": 150,
            "total_pages": 8
        }
    }
}

// Error
{
    "error": {
        "code": "VALIDATION_ERROR",
        "message": "invalid serial number format",
        "details": {
            "field": "serial_no",
            "value": "abc"
        },
        "request_id": "550e8400-e29b-41d4-a716-446655440000",
        "timestamp": "2026-01-15T10:00:00Z"
    }
}

// Async (accepted)
{
    "data": {
        "id": "...",
        "status": "PENDING",
        "resource_url": "/api/v1/jobs/..."
    }
}
```

### D.2 HTTP Status Codes Convention

| Status | ใช้เมื่อ |
|:---|:---|
| **200 OK** | GET, PUT, PATCH สำเร็จ |
| **201 Created** | POST สำเร็จ (พร้อม Location header) |
| **202 Accepted** | Async operation queued |
| **204 No Content** | DELETE สำเร็จ, action ที่ไม่ return body |
| **400 Bad Request** | Validation error |
| **401 Unauthorized** | ไม่มี token / token หมดอายุ |
| **403 Forbidden** | ไม่มีสิทธิ์ |
| **404 Not Found** | ไม่พบ resource |
| **409 Conflict** | Duplicate, version mismatch |
| **422 Unprocessable** | Business rule violation |
| **429 Too Many Requests** | Rate limit |
| **500 Internal Server Error** | Bug, unexpected |
| **503 Service Unavailable** | Device offline, DB down |

### D.3 Pagination Pattern

```go
type Pagination struct {
    Page       int   `json:"page"`
    PageSize   int   `json:"page_size"`
    Total      int64 `json:"total"`
    TotalPages int   `json:"total_pages"`
    HasNext    bool  `json:"has_next"`
    HasPrev    bool  `json:"has_prev"`
}

func NewPagination(page, pageSize int, total int64) Pagination {
    totalPages := int(total) / pageSize
    if int(total)%pageSize != 0 {
        totalPages++
    }
    return Pagination{
        Page:       page,
        PageSize:   pageSize,
        Total:      total,
        TotalPages: totalPages,
        HasNext:    page < totalPages,
        HasPrev:    page > 1,
    }
}
```

### D.4 Filtering & Sorting

```go
// Query parameters:
//   ?status=ACTIVE&status=PENDING        (multi-value)
//   ?created_from=2026-01-01T00:00:00Z
//   ?created_to=2026-01-31T23:59:59Z
//   ?search=hello
//   ?sort_by=created_at&sort_order=desc
//   ?page=1&page_size=20

// Whitelist sort fields เพื่อความปลอดภัย
var allowedSortFields = map[string]bool{
    "created_at": true,
    "updated_at": true,
    "name":       true,
    "status":     true,
}
```

### D.5 ETag & Conditional GET

```go
func (h *DeviceHandler) GetWithETag(c *gin.Context) {
    tenantID := c.MustGet("tenant_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    // 1. Load resource
    device, err := h.getUC.Execute(c.Request.Context(), application.GetDeviceInput{
        TenantID: tenantID,
        DeviceID: id,
    })
    if err != nil {
        h.respondError(c, err)
        return
    }

    // 2. Compute ETag
    etag := fmt.Sprintf(`W/"%x"`, md5.Sum([]byte(fmt.Sprintf("%s-%d",
        device.ID, device.UpdatedAt.UnixNano()))))

    // 3. Check If-None-Match
    if c.GetHeader("If-None-Match") == etag {
        c.Status(http.StatusNotModified)
        return
    }
    c.Header("ETag", etag)
    c.Header("Cache-Control", "private, max-age=60")
    c.JSON(http.StatusOK, gin.H{"data": device})
}
```

---

## 🅴 PART 8E — WEBSOCKET HUB

### E.1 Hub Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                    WebSocket Hub                             │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Tenant A     │  │ Tenant B     │  │ Tenant C     │      │
│  │ ┌──┐ ┌──┐    │  │ ┌──┐         │  │ ┌──┐ ┌──┐    │      │
│  │ │C1│ │C2│    │  │ │C3│         │  │ │C4│ │C5│    │      │
│  │ └──┘ └──┘    │  │ └──┘         │  │ └──┘ └──┘    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                              │
│  Kafka Consumer ──► BroadcastToTenant ──► Client channels   │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### E.2 Hub Implementation

```go
// internal/modules/realtime/interfaces/websocket/hub.go
package websocket

import (
    "encoding/json"
    "sync"

    "github.com/gorilla/websocket"
)

type Hub struct {
    mu      sync.RWMutex
    tenants map[string]map[*Client]bool  // tenantID → clients
    rooms   map[string]map[*Client]bool  // roomID → clients
    logger  application.Logger
}

func NewHub(logger application.Logger) *Hub {
    return &Hub{
        tenants: make(map[string]map[*Client]bool),
        rooms:   make(map[string]map[*Client]bool),
        logger:  logger,
    }
}

func (h *Hub) Register(c *Client) {
    h.mu.Lock()
    defer h.mu.Unlock()
    if h.tenants[c.TenantID] == nil {
        h.tenants[c.TenantID] = make(map[*Client]bool)
    }
    h.tenants[c.TenantID][c] = true
    h.logger.Info("ws client registered",
        application.Field{Key: "tenant", Value: c.TenantID},
        application.Field{Key: "user", Value: c.UserID},
    )
}

func (h *Hub) Unregister(c *Client) {
    h.mu.Lock()
    defer h.mu.Unlock()
    if m, ok := h.tenants[c.TenantID]; ok {
        if _, exists := m[c]; exists {
            delete(m, c)
            close(c.send)
        }
        if len(m) == 0 {
            delete(h.tenants, c.TenantID)
        }
    }
}

// BroadcastToTenant ส่งให้ทุก client ของ tenant
func (h *Hub) BroadcastToTenant(tenantID string, event Event) {
    data, _ := json.Marshal(event)
    h.mu.RLock()
    defer h.mu.RUnlock()
    for c := range h.tenants[tenantID] {
        select {
        case c.send <- data:
        default:
            // buffer เต็ม — drop
        }
    }
}

// BroadcastToRoom ส่งให้ client ที่ subscribe room (เช่น device/:id)
func (h *Hub) BroadcastToRoom(roomID string, event Event) {
    data, _ := json.Marshal(event)
    h.mu.RLock()
    defer h.mu.RUnlock()
    for c := range h.rooms[roomID] {
        select {
        case c.send <- data:
        default:
        }
    }
}

func (h *Hub) JoinRoom(c *Client, roomID string) {
    h.mu.Lock()
    defer h.mu.Unlock()
    if h.rooms[roomID] == nil {
        h.rooms[roomID] = make(map[*Client]bool)
    }
    h.rooms[roomID][c] = true
}

func (h *Hub) LeaveRoom(c *Client, roomID string) {
    h.mu.Lock()
    defer h.mu.Unlock()
    if m, ok := h.rooms[roomID]; ok {
        delete(m, c)
    }
}

func (h *Hub) CountByTenant(tenantID string) int {
    h.mu.RLock()
    defer h.mu.RUnlock()
    return len(h.tenants[tenantID])
}
```

### E.3 Client Implementation

```go
// internal/modules/realtime/interfaces/websocket/client.go
package websocket

import (
    "context"
    "encoding/json"
    "time"

    "github.com/gorilla/websocket"
    "github.com/google/uuid"
)

const (
    writeWait      = 10 * time.Second
    pongWait       = 60 * time.Second
    pingPeriod     = (pongWait * 9) / 10
    maxMessageSize = 8192
)

type Client struct {
    ID       uuid.UUID
    TenantID string
    UserID   string
    Conn     *websocket.Conn
    send     chan []byte
    hub      *Hub
    logger   application.Logger
    subs     map[string]bool // rooms ที่ subscribe
}

type Event struct {
    Type      string          `json:"type"`
    Data      json.RawMessage `json:"data"`
    Timestamp time.Time       `json:"timestamp"`
    TenantID  string          `json:"tenant_id,omitempty"`
}

func NewClient(conn *websocket.Conn, tenantID, userID string, hub *Hub, logger application.Logger) *Client {
    return &Client{
        ID:       uuid.New(),
        TenantID: tenantID,
        UserID:   userID,
        Conn:     conn,
        send:     make(chan []byte, 256),
        hub:      hub,
        logger:   logger,
        subs:     make(map[string]bool),
    }
}

// ReadPump — อ่านข้อความจาก client (subscribe/unsubscribe)
func (c *Client) ReadPump(ctx context.Context) {
    defer func() {
        c.hub.Unregister(c)
        c.Conn.Close()
    }()

    c.Conn.SetReadLimit(maxMessageSize)
    c.Conn.SetReadDeadline(time.Now().Add(pongWait))
    c.Conn.SetPongHandler(func(string) error {
        c.Conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        _, message, err := c.Conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                c.logger.Warn("ws read error", application.Field{Key: "error", Value: err})
            }
            return
        }
        c.handleMessage(message)
    }
}

// WritePump — ส่งข้อความไป client
func (c *Client) WritePump(ctx context.Context) {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.Conn.Close()
    }()

    for {
        select {
        case <-ctx.Done():
            c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
            return
        case message, ok := <-c.send:
            c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
                return
            }
        case <-ticker.C:
            c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

type clientMessage struct {
    Action string `json:"action"` // subscribe, unsubscribe, ping
    Room   string `json:"room,omitempty"`
}

func (c *Client) handleMessage(data []byte) {
    var msg clientMessage
    if err := json.Unmarshal(data, &msg); err != nil {
        return
    }
    switch msg.Action {
    case "subscribe":
        c.hub.JoinRoom(c, msg.Room)
        c.subs[msg.Room] = true
    case "unsubscribe":
        c.hub.LeaveRoom(c, msg.Room)
        delete(c.subs, msg.Room)
    case "ping":
        c.send <- []byte(`{"type":"pong"}`)
    }
}
```

### E.4 WebSocket Handler

```go
// internal/modules/realtime/interfaces/websocket/handler.go
package websocket

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

type Handler struct {
    hub      *Hub
    authSvc  middleware.AuthService
    logger   application.Logger
    upgrader websocket.Upgrader
}

func NewHandler(hub *Hub, authSvc middleware.AuthService, cfg WSConfig, logger application.Logger) *Handler {
    return &Handler{
        hub:     hub,
        authSvc: authSvc,
        logger:  logger,
        upgrader: websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
            CheckOrigin: func(r *http.Request) bool {
                origin := r.Header.Get("Origin")
                return isAllowedOrigin(origin, cfg.AllowedOrigins)
            },
        },
    }
}

// Handle — GET /ws?token=<jwt>
func (h *Handler) Handle(c *gin.Context) {
    // Authenticate via query param (browsers ไม่ส่ง header ใน WS)
    token := c.Query("token")
    if token == "" {
        // fallback: Authorization header
        token = extractBearer(c.GetHeader("Authorization"))
    }
    if token == "" {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token required"})
        return
    }
    claims, err := h.authSvc.ValidateToken(c.Request.Context(), token)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
        return
    }

    conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        h.logger.Error("ws upgrade failed", application.Field{Key: "error", Value: err})
        return
    }

    client := NewClient(conn, claims.TenantID.String(), claims.UserID.String(), h.hub, h.logger)
    h.hub.Register(client)

    // ส่ง welcome
    welcome, _ := json.Marshal(Event{
        Type:      "connected",
        Timestamp: time.Now().UTC(),
        Data:      json.RawMessage(`{"status":"connected"}`),
    })
    client.send <- welcome

    ctx, cancel := context.WithCancel(c.Request.Context())
    defer cancel()

    go client.WritePump(ctx)
    go client.ReadPump(ctx)

    // Block จนกว่า client จะ disconnect
    <-ctx.Done()
}
```

### E.5 Kafka → WS Bridge

```go
// internal/modules/realtime/infrastructure/messaging/consumers/ws_broadcast.go
package consumers

type WSBroadcastConsumer struct {
    hub    *websocket.Hub
    logger application.Logger
}

func (c *WSBroadcastConsumer) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
    tenantID := extractHeader(msg, "tenant_id")
    if tenantID == "" {
        return nil
    }

    event := websocket.Event{
        Type:      msg.Topic,
        Data:      json.RawMessage(msg.Value),
        Timestamp: time.Now().UTC(),
        TenantID:  tenantID,
    }
    c.hub.BroadcastToTenant(tenantID, event)
    return nil
}

func extractHeader(msg *sarama.ConsumerMessage, key string) string {
    for _, h := range msg.Headers {
        if string(h.Key) == key {
            return string(h.Value)
        }
    }
    return ""
}
```

---

## 🅵 PART 8F — MQTT GATEWAY

### F.1 Gateway Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    MQTT Broker (EMQX / Mosquitto)           │
│                                                             │
│  iot/+/telemetry    iot/+/heartbeat    iot/+/status         │
│  iot/+/fault        device/+/shadow/update                   │
│                                                             │
└──────────────────────────┬──────────────────────────────────┘
                           │ subscribe
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                  MQTT Gateway (cmd/mqtt-gateway)            │
│                                                             │
│  ┌────────────────┐  ┌────────────────┐  ┌──────────────┐  │
│  │ Subscriber     │  │ Parser Registry│  │ Authenticator│  │
│  │ (paho)         │──│ (modbus/lora)  │──│ (device token)│  │
│  └────────────────┘  └────────────────┘  └──────────────┘  │
│           │                                                │
│           ▼                                                │
│  ┌────────────────────────────────────────────────────┐    │
│  │           Use Cases (IngestTelemetry, etc.)        │    │
│  └────────────────────────────────────────────────────┘    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### F.2 Topic Convention

```
iot/{serial}/telemetry        # sensor data (QoS 1)
iot/{serial}/heartbeat        # keep-alive (QoS 0)
iot/{serial}/status           # device status (QoS 1)
iot/{serial}/fault            # fault report (QoS 1)
iot/{serial}/command          # command from cloud (QoS 1)
iot/{serial}/command/ack      # command ack (QoS 1)
device/{serial}/shadow/update # shadow sync (QoS 1)
device/{serial}/shadow/desired
gateway/{gateway_id}/status   # gateway online/offline
```

### F.3 Gateway Main

```go
// cmd/mqtt-gateway/main.go
package main

func main() {
    cfg := config.Load()
    logger := logger.New(cfg.LogLevel)

    // 1. Infrastructure
    db := postgres.NewDB(cfg.DB)
    redisClient := redis.NewClient(cfg.Redis)
    influxClient := influxdb.NewClient(cfg.Influx)
    mqttClient, err := mqtt.NewClient(cfg.MQTT, logger)
    if err != nil {
        log.Fatal(err)
    }
    defer mqttClient.Disconnect()

    // 2. Repositories
    deviceRepo := devicepg.NewDeviceRepository(db)
    telemetryRepo := deviceinflux.NewTelemetryRepository(influxClient, cfg.Influx.Bucket)
    auditRepo := auditpg.NewAuditRepository(db)

    // 3. Use Cases
    ingestUC := application.NewIngestTelemetryUseCase(
        deviceRepo, telemetryRepo, automationSvc, wsHub, producer,
    )

    // 4. Subscriber
    subscriber := iot.NewMQTTSubscriber(mqttClient, ingestUC, logger)
    if err := subscriber.Start(context.Background()); err != nil {
        log.Fatal(err)
    }

    // 5. Publisher (listen to command events → publish to MQTT)
    publisher := iot.NewMQTTPublisher(mqttClient)
    cmdConsumer := consumers.NewCommandPublishConsumer(publisher, deviceRepo)
    consumerGroup, _ := kafka.NewConsumerGroup(
        cfg.Kafka.Brokers,
        "mqtt-command-publisher",
        []string{"icmon.device.command.issued"},
        cmdConsumer.Handle,
        logger,
    )
    go consumerGroup.Run(context.Background())

    // 6. Graceful shutdown
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    <-sigCh

    logger.Info("shutting down MQTT gateway")
    consumerGroup.Close()
    subscriber.Stop()
    mqttClient.Disconnect()
}
```

### F.4 Device Authentication (MQTT)

```go
// internal/shared/infrastructure/mqtt/auth.go
package mqtt

// AuthenticateDevice — ตรวจ device token จาก MQTT username/password
func AuthenticateDevice(ctx context.Context, repo repository.DeviceRepository, username, password string) (*entity.Device, error) {
    // username = device serial หรือ client ID
    device, err := repo.FindBySerial(ctx, uuid.Nil, vo.DeviceSerial(username))
    if err != nil {
        return nil, domainerrors.ErrDeviceNotFound
    }
    // เทียบ password กับ device token hash
    if err := bcrypt.CompareHashAndPassword([]byte(device.DeviceTokenHash), []byte(password)); err != nil {
        return nil, domainerrors.ErrInvalidToken
    }
    return device, nil
}
```

### F.5 MQTT ACL (EMQX-style)

```json
// iot/{serial}/telemetry — device publish only
{
    "topic": "iot/{serial}/telemetry",
    "permission": "allow",
    "action": "publish",
    "username": "{serial}"
}
// iot/{serial}/command — device subscribe only
{
    "topic": "iot/{serial}/command",
    "permission": "allow",
    "action": "subscribe",
    "username": "{serial}"
}
```

---

## 🅶 PART 8G — CLI & SCHEDULER ENTRY POINTS

### G.1 Entry Points

```
cmd/
├── api/main.go                 # REST + WebSocket
├── scheduler/main.go           # Cron jobs
├── migrate/main.go             # DB migrations
├── mqtt-gateway/main.go        # MQTT bridge
├── cli/main.go                 # Admin CLI
└── workers/
    ├── device-telemetry/main.go
    ├── ai-inference/main.go
    ├── notification/main.go
    ├── erp-sync/main.go
    └── outbox-relay/main.go
```

### G.2 API Main

```go
// cmd/api/main.go
package main

func main() {
    cfg := config.Load()
    logger := logger.New(cfg.LogLevel, cfg.Env)

    // --- Infrastructure ---
    db, err := postgres.NewDB(cfg.DB)
    if err != nil { log.Fatal(err) }
    redisClient := redis.NewClient(cfg.Redis)
    kafkaProducer, _ := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Prefix, logger)
    influxClient := influxdb.NewClient(cfg.Influx)
    esClient, _ := es.NewClient(cfg.ES.URLs, cfg.ES.Prefix)

    // --- Auth service ---
    authSvc := auth.NewService(cfg.JWT, redisClient)

    // --- Wire modules ---
    deviceModule := devicemodule.New(devicemodule.Deps{
        DB: db, Redis: redisClient, Kafka: kafkaProducer, Influx: influxClient, ES: esClient,
        Logger: logger, AuthSvc: authSvc,
    })
    customerModule := customermodule.New(...)
    packageModule := packagemodule.New(...)
    erpModule := erpmodule.New(...)
    logisticsModule := logisticsmodule.New(...)
    reportModule := reportmodule.New(...)
    wsHub := websocket.NewHub(logger)

    // --- Kafka → WS Bridge ---
    wsConsumer := consumers.NewWSBroadcastConsumer(wsHub, logger)
    wsGroup, _ := kafka.NewConsumerGroup(
        cfg.Kafka.Brokers, "api-ws-broadcaster",
        []string{
            "icmon.device.telemetry.ingested",
            "icmon.device.command.acked",
            "icmon.device.device.online",
            "icmon.device.device.offline",
            "icmon.logistics.temperature.breached",
        }, wsConsumer.Handle, logger,
    )
    go wsGroup.Run(context.Background())

    // --- HTTP Router ---
    router := sharedhttp.NewRouter(
        cfg, logger, authSvc,
        deviceModule.HTTPHandlers(),
        customerModule.HTTPHandlers(),
        packageModule.HTTPHandlers(),
        erpModule.HTTPHandlers(),
        logisticsModule.HTTPHandlers(),
        reportModule.HTTPHandlers(),
    )

    // --- WebSocket route ---
    wsHandler := websocket.NewHandler(wsHub, authSvc, cfg.WS, logger)
    router.GET("/ws", wsHandler.Handle)

    // --- HTTP Server ---
    srv := &http.Server{
        Addr:    ":" + cfg.Port,
        Handler: router,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 30 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    go func() {
        logger.Info("api server started", application.Field{Key: "port", Value: cfg.Port})
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    // --- Graceful shutdown ---
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    logger.Info("shutting down")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    _ = srv.Shutdown(ctx)
    _ = wsGroup.Close()
    _ = kafkaProducer.Close()
    _ = redisClient.Close()
    _ = db.Close()
    logger.Info("server stopped")
}
```

### G.3 Admin CLI

```go
// cmd/cli/main.go
package main

import (
    "github.com/spf13/cobra"
)

func main() {
    rootCmd := &cobra.Command{Use: "icmon", Short: "icmongolang CLI"}

    // Migrate
    migrateCmd := &cobra.Command{
        Use:   "migrate [up|down|status]",
        Short: "Database migrations",
        Run:   runMigrate,
    }
    migrateCmd.Flags().Int("steps", 0, "number of steps")

    // Seed
    seedCmd := &cobra.Command{
        Use:   "seed [minimal|standard|demo]",
        Short: "Seed fixtures",
        Run:   runSeed,
    }

    // Device
    deviceCmd := &cobra.Command{Use: "device", Short: "Device management"}
    deviceCmd.AddCommand(
        &cobra.Command{Use: "list", Run: runDeviceList},
        &cobra.Command{Use: "info <serial>", Run: runDeviceInfo},
        &cobra.Command{Use: "rotate-token <serial>", Run: runDeviceRotateToken},
    )

    // Tenant
    tenantCmd := &cobra.Command{Use: "tenant", Short: "Tenant management"}
    tenantCmd.AddCommand(
        &cobra.Command{Use: "create <slug>", Run: runTenantCreate},
        &cobra.Command{Use: "list", Run: runTenantList},
    )

    rootCmd.AddCommand(migrateCmd, seedCmd, deviceCmd, tenantCmd)

    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

---

## 🅷 PART 8H — API VERSIONING & DEPRECATION

### H.1 URL Versioning

```
/api/v1/devices
/api/v2/devices   (future)
```

### H.2 Deprecation Headers

```go
func Deprecated(sunset time.Time, replacement string) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Deprecation", "true")
        c.Header("Sunset", sunset.Format(http.TimeFormat))
        if replacement != "" {
            c.Header("Link", fmt.Sprintf(`<%s>; rel="successor-version"`, replacement))
        }
        c.Next()
    }
}

// Usage
v1 := r.Group("/api/v1")
v1.Use(Deprecated(time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC), "/api/v2"))
```

### H.3 Breaking Change Policy

| ประเภท | ตัวอย่าง | Version |
|:---|:---|:---|
| **Non-breaking** | เพิ่ม field, endpoint ใหม่ | same version |
| **Breaking** | ลบ field, เปลี่ยน type, เปลี่ยน semantic | bump major |
| **Deprecation** | mark deprecated, แจ้งล่วงหน้า 6 เดือน | same + header |
| **Sunset** | ลบ version เก่า | after 12 เดือน |

---

## 🅸 PART 8I — ERROR HANDLING & OBSERVABILITY

### I.1 Panic Recovery → 500 + Stack

```go
// ดูที่ C.9 Recovery Middleware
```

### I.2 Structured Logging

```json
{
    "level": "info",
    "timestamp": "2026-01-15T10:00:00.123Z",
    "message": "request",
    "method": "POST",
    "path": "/api/v1/devices",
    "status": 201,
    "latency_ms": 45,
    "ip": "203.0.113.42",
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "tenant_id": "11111111-...",
    "user_id": "22222222-...",
    "trace_id": "abc123"
}
```

### I.3 Prometheus Metrics

```go
// internal/shared/interfaces/http/middleware/metrics.go
package middleware

var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{Name: "http_requests_total"},
        []string{"method", "path", "status", "tenant"},
    )
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )
    httpInFlight = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "http_requests_in_flight",
    })
)

func PrometheusMetrics() gin.HandlerFunc {
    return func(c *gin.Context) {
        httpInFlight.Inc()
        defer httpInFlight.Dec()

        start := time.Now()
        c.Next()

        status := strconv.Itoa(c.Writer.Status())
        path := c.FullPath()
        if path == "" {
            path = "unknown"
        }
        tenant := c.GetString("tenant_id")

        httpRequestsTotal.WithLabelValues(c.Request.Method, path, status, tenant).Inc()
        httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(time.Since(start).Seconds())
    }
}
```

### I.4 Distributed Tracing (OpenTelemetry)

```go
// internal/shared/infrastructure/tracing/otel.go
package tracing

func InitTracer(cfg Config) (func(context.Context) error, error) {
    exporter, err := otlptracehttp.New(context.Background(),
        otlptracehttp.WithEndpoint(cfg.Endpoint),
        otlptracehttp.WithInsecure(),
    )
    if err != nil {
        return nil, err
    }
    tp := tracesdk.NewTracerProvider(
        tracesdk.WithBatcher(exporter),
        tracesdk.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String(cfg.ServiceName),
            semconv.ServiceVersionKey.String(cfg.Version),
        )),
        tracesdk.WithSampler(tracesdk.TraceIDRatioBased(cfg.SampleRate)),
    )
    otel.SetTracerProvider(tp)
    return tp.Shutdown, nil
}
```

### I.5 Health & Readiness

```go
// internal/shared/interfaces/http/health.go
package http

type HealthHandler struct {
    db       *sql.DB
    redis    *redis.Client
    kafka    sarama.Client
    influx   influxdb2.Client
}

func (h *HealthHandler) Health(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().UTC()})
}

func (h *HealthHandler) Ready(c *gin.Context) {
    checks := map[string]string{}

    // DB
    if err := h.db.PingContext(c.Request.Context()); err != nil {
        checks["postgres"] = "down: " + err.Error()
    } else {
        checks["postgres"] = "ok"
    }

    // Redis
    if err := h.redis.Ping(c.Request.Context()).Err(); err != nil {
        checks["redis"] = "down: " + err.Error()
    } else {
        checks["redis"] = "ok"
    }

    // Kafka
    if h.kafka != nil && len(h.kafka.Brokers()) > 0 {
        checks["kafka"] = "ok"
    } else {
        checks["kafka"] = "down"
    }

    // Influx
    if err := h.influx.Ping(c.Request.Context()); err != nil {
        checks["influxdb"] = "down: " + err.Error()
    } else {
        checks["influxdb"] = "ok"
    }

    status := http.StatusOK
    for _, v := range checks {
        if strings.HasPrefix(v, "down") {
            status = http.StatusServiceUnavailable
            break
        }
    }
    c.JSON(status, gin.H{"checks": checks, "time": time.Now().UTC()})
}
```

---

## 🎯 PART 8 — SUMMARY

| Deliverable | Count | Location |
|---|:-:|---|
| **HTTP Handlers** | 120+ | `interfaces/http/` |
| **Route Groups** | 15 | `interfaces/http/routes.go` |
| **Middleware** | 10 | `shared/interfaces/http/middleware/` |
| **WebSocket Hub + Client** | 2 | `interfaces/websocket/` |
| **MQTT Gateway** | 1 | `cmd/mqtt-gateway/` |
| **Entry Points** | 9 | `cmd/` |
| **Metrics** | 20+ | Prometheus |

---
---

# 🧪 PART 12 — TESTING STRATEGY

> **ขนาด**: ใหญ่ — แยก 8 ตอนย่อย
> **Part 12A**: Testing Pyramid & Strategy
> **Part 12B**: Unit Tests (Domain + Application)
> **Part 12C**: Integration Tests (DB, Kafka, Redis)
> **Part 12D**: E2E / API Tests (Postman + Newman)
> **Part 12E**: Contract Tests (Consumer-Driven)
> **Part 12F**: Load & Performance Tests
> **Part 12G**: Test Fixtures & Builders
> **Part 12H**: CI Test Pipeline

---

## 🅰️ PART 12A — TESTING PYRAMID & STRATEGY

### A.1 Testing Pyramid

```
                       ▲
                      ╱ ╲
                     ╱   ╲       E2E / UI (5%)
                    ╱─────╲      ~50 tests · slow
                   ╱       ╲
                  ╱─────────╲    Integration (25%)
                 ╱           ╲   ~250 tests · medium
                ╱─────────────╲
               ╱               ╲  Unit (70%)
              ╱─────────────────╲ ~700 tests · fast
             ╱___________________╲
```

### A.2 Coverage Targets

| Layer | Coverage Target | เหตุผล |
|:---|:---|:---|
| **Domain** | ≥ 95% | Business logic ต้องแม่นยำ |
| **Application** | ≥ 85% | Orchestration + error handling |
| **Infrastructure** | ≥ 70% | Repos + adapters |
| **Interface** | ≥ 75% | Handlers + middleware |
| **Overall** | ≥ 80% | Industry standard |

### A.3 Test Types & Tools

| Type | Tool | Speed | Runs When |
|:---|:---|:---|:---|
| Unit | `testing` + `testify` | Fast (<1s) | Every commit |
| Integration | `testcontainers-go` | Medium (5-30s) | PR |
| E2E | Postman + Newman | Slow (1-5min) | PR + Nightly |
| Contract | `pact-go` | Medium | PR |
| Load | `k6` + `vegeta` | Very slow (10-60min) | Weekly |
| Fuzz | Go 1.18+ fuzzing | Slow | Nightly |

---

## 🅱️ PART 12B — UNIT TESTS

### B.1 Domain Entity Test (Table-Driven)

```go
// internal/modules/device/domain/entity/device_test.go
package entity_test

import (
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "icmongolang/internal/modules/device/domain/entity"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
    vo "icmongolang/internal/modules/device/domain/value_object"
)

func TestNewDevice(t *testing.T) {
    tenantID := uuid.New()
    customerID := uuid.New()
    siteID := uuid.New()
    actorID := uuid.New()

    tests := []struct {
        name       string
        tenantID   uuid.UUID
        siteID     uuid.UUID
        serial     string
        dtype      vo.DeviceType
        wantErr    error
    }{
        {
            name:     "valid sensor device",
            tenantID: tenantID, siteID: siteID,
            serial:   "SN-TEST-001",
            dtype:    vo.DeviceTypeSensor,
            wantErr:  nil,
        },
        {
            name:     "empty tenant",
            tenantID: uuid.Nil, siteID: siteID,
            serial:   "SN-TEST-002",
            dtype:    vo.DeviceTypeSensor,
            wantErr:  domainerrors.ErrInvalidTenant,
        },
        {
            name:     "invalid serial format",
            tenantID: tenantID, siteID: siteID,
            serial:   "abc", // lowercase
            dtype:    vo.DeviceTypeSensor,
            wantErr:  domainerrors.ErrInvalidSerial,
        },
        {
            name:     "invalid device type",
            tenantID: tenantID, siteID: siteID,
            serial:   "SN-TEST-003",
            dtype:    "INVALID",
            wantErr:  domainerrors.ErrInvalidDeviceType,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            d, err := entity.NewDevice(
                tt.tenantID, customerID, tt.siteID, uuid.Nil,
                tt.serial, tt.dtype, vo.ProtocolMQTT, "Test", actorID,
            )
            if tt.wantErr != nil {
                require.Error(t, err)
                assert.ErrorIs(t, err, tt.wantErr)
                assert.Nil(t, d)
            } else {
                require.NoError(t, err)
                assert.NotNil(t, d)
                assert.NotEqual(t, uuid.Nil, d.ID)
                assert.Equal(t, vo.DeviceStatusOffline, d.Status)
                assert.Len(t, d.PullEvents(), 1) // DeviceCreated
            }
        })
    }
}

func TestDevice_MarkOnline(t *testing.T) {
    d := newTestDevice(t)

    before := time.Now().UTC()
    d.Heartbeat()
    after := time.Now().UTC()

    assert.Equal(t, vo.DeviceStatusOnline, d.Status)
    require.NotNil(t, d.LastSeenAt)
    assert.True(t, d.LastSeenAt.After(before) || d.LastSeenAt.Equal(before))
    assert.True(t, d.LastSeenAt.Before(after) || d.LastSeenAt.Equal(after))

    events := d.PullEvents()
    assert.Len(t, events, 2) // Created + Online
    assert.Equal(t, "device.device.online", events[1].EventName())
}

func TestDevice_IssueCommand(t *testing.T) {
    d := newTestDevice(t)
    d.Heartbeat() // ต้อง online ก่อน

    actor := uuid.New()
    cmd, err := d.IssueCommand("on", map[string]interface{}{"duration": 300}, 3, actor, time.Minute)
    require.NoError(t, err)
    assert.NotNil(t, cmd)
    assert.Equal(t, vo.CommandStatusPending, cmd.Status)
    assert.Equal(t, 3, cmd.Priority)

    // Test offline
    d.MarkOfflineIfStale(0) // force
    _, err = d.IssueCommand("on", nil, 3, actor, time.Minute)
    assert.ErrorIs(t, err, domainerrors.ErrDeviceOffline)
}

func newTestDevice(t *testing.T) *entity.Device {
    d, err := entity.NewDevice(
        uuid.New(), uuid.New(), uuid.New(), uuid.Nil,
        "SN-TEST-001", vo.DeviceTypeActuator, vo.ProtocolMQTT, "Test", uuid.New(),
    )
    require.NoError(t, err)
    d.PullEvents() // clear
    return d
}
```

### B.2 Value Object Test (Fuzz)

```go
// internal/modules/device/domain/value_object/device_serial_test.go
package valueobject_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    vo "icmongolang/internal/modules/device/domain/value_object"
)

func TestDeviceSerial(t *testing.T) {
    tests := []struct {
        in      string
        wantErr bool
    }{
        {"SN-001", false},
        {"ABC123", false},
        {"A-B-C", false},
        {"", true},
        {"ab", true},           // lowercase
        {"!!!", true},          // special
        {"A", true},            // too short
        {string(make([]byte, 200)), true}, // too long
    }
    for _, tt := range tests {
        t.Run(tt.in, func(t *testing.T) {
            _, err := vo.NewDeviceSerial(tt.in)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

// Go 1.18+ fuzzing
func FuzzDeviceSerial(f *testing.F) {
    f.Add("SN-001")
    f.Add("ABC")
    f.Add("")
    f.Fuzz(func(t *testing.T, s string) {
        serial, err := vo.NewDeviceSerial(s)
        if err == nil {
            // ถ้า parse ผ่าน ต้อง IsValid() = true
            assert.True(t, serial.IsValid())
            // ต้อง normalize ถูก
            assert.Equal(t, serial.String(), strings.ToUpper(strings.TrimSpace(s)))
        }
    })
}
```

### B.3 Use Case Test (Mock Repo)

```go
// internal/modules/device/application/register_device_test.go
package application_test

import (
    "context"
    "testing"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"

    "icmongolang/internal/modules/device/application"
    "icmongolang/internal/modules/device/application/mocks"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
)

func TestRegisterDeviceUseCase_Execute(t *testing.T) {
    ctx := context.Background()
    tenantID := uuid.New()
    actorID := uuid.New()
    customerID := uuid.New()
    siteID := uuid.New()

    t.Run("success", func(t *testing.T) {
        repo := mocks.NewDeviceRepository(t)
        txMgr := mocks.NewTransactionManager(t)
        producer := mocks.NewEventProducer(t)
        quotaSvc := service.NewQuotaService()

        // Expectations
        repo.On("ExistsBySerial", mock.Anything, tenantID, mock.Anything).Return(false, nil)
        txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
            fn := args.Get(1).(func(context.Context) error)
            _ = fn(ctx)
        })
        repo.On("Save", mock.Anything, mock.Anything).Return(nil)
        producer.On("PublishBatch", mock.Anything, mock.Anything).Return(nil)

        uc := application.NewRegisterDeviceUseCase(repo, txMgr, producer, quotaSvc, nil, nil)

        out, err := uc.Execute(ctx, application.RegisterDeviceInput{
            TenantID:   tenantID,
            ActorID:    actorID,
            CustomerID: customerID,
            SiteID:     siteID,
            SerialNo:   "SN-NEW-001",
            Name:       "New Device",
            Type:       "SENSOR",
            Protocol:   "MQTT",
        })

        require.NoError(t, err)
        assert.NotEqual(t, uuid.Nil, out.ID)
        assert.Equal(t, "SN-NEW-001", out.SerialNo)
    })

    t.Run("duplicate serial", func(t *testing.T) {
        repo := mocks.NewDeviceRepository(t)
        repo.On("ExistsBySerial", mock.Anything, tenantID, mock.Anything).Return(true, nil)

        uc := application.NewRegisterDeviceUseCase(repo, nil, nil, nil, nil, nil)

        _, err := uc.Execute(ctx, application.RegisterDeviceInput{
            TenantID: tenantID, ActorID: actorID,
            SerialNo: "SN-DUP-001", Type: "SENSOR", Protocol: "MQTT",
        })

        assert.ErrorIs(t, err, domainerrors.ErrSerialAlreadyExists)
    })

    t.Run("invalid tenant", func(t *testing.T) {
        uc := application.NewRegisterDeviceUseCase(nil, nil, nil, nil, nil, nil)
        _, err := uc.Execute(ctx, application.RegisterDeviceInput{
            TenantID: uuid.Nil, ActorID: actorID, SerialNo: "SN-001",
        })
        assert.ErrorIs(t, err, domainerrors.ErrInvalidTenant)
    })
}
```

### B.4 Mock Generation (mockery)

```yaml
# .mockery.yaml
with-expecter: true
dir: "{{.InterfaceDir}}/mocks"
outpkg: mocks
packages:
  icmongolang/internal/modules/device/domain/repository:
    interfaces:
      DeviceRepository:
      TelemetryRepository:
      CommandRepository:
  icmongolang/internal/shared/application:
    interfaces:
      EventProducer:
      Cache:
      Logger:
      TransactionManager:
```

```bash
# Generate
mockery --config .mockery.yaml

# หรือ go:generate
//go:generate mockery --name=DeviceRepository --with-expecter
```

---

## 🅲 PART 12C — INTEGRATION TESTS

### C.1 Testcontainers Setup

```go
// test/integration/testenv.go
package integration

import (
    "context"
    "testing"

    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/testcontainers/testcontainers-go/modules/redis"
    "github.com/testcontainers/testcontainers-go/modules/kafka"
    "github.com/testcontainers/testcontainers-go/wait"
)

type TestEnv struct {
    Postgres *postgres.PostgresContainer
    Redis    *redis.RedisContainer
    Kafka    *kafka.KafkaContainer
    DB       *gorm.DB
    RedisCli *redis.Client
}

func NewTestEnv(t *testing.T) *TestEnv {
    ctx := context.Background()

    // Postgres
    pgContainer, err := postgres.Run(ctx,
        "postgres:16-alpine",
        postgres.WithDatabase("test"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).WithStartupTimeout(30*time.Second)),
    )
    require.NoError(t, err)

    // Redis
    redisContainer, err := redis.Run(ctx, "redis:7-alpine")
    require.NoError(t, err)

    // Kafka
    kafkaContainer, err := kafka.Run(ctx,
        "confluentinc/cp-kafka:latest",
        kafka.WithClusterID("test-cluster"),
    )
    require.NoError(t, err)

    // Connect DB
    dsn, _ := pgContainer.ConnectionString(ctx, "sslmode=disable")
    db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

    // Run migrations
    runMigrations(t, db)

    t.Cleanup(func() {
        pgContainer.Terminate(ctx)
        redisContainer.Terminate(ctx)
        kafkaContainer.Terminate(ctx)
    })

    return &TestEnv{Postgres: pgContainer, Redis: redisContainer, Kafka: kafkaContainer, DB: db}
}
```

### C.2 Repository Integration Test

```go
// internal/modules/device/infrastructure/persistence/postgres/device_repo_test.go
//go:build integration

package postgres_test

import (
    "context"
    "testing"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "icmongolang/internal/modules/device/domain/entity"
    devicepg "icmongolang/internal/modules/device/infrastructure/persistence/postgres"
    vo "icmongolang/internal/modules/device/domain/value_object"
    "icmongolang/test/integration"
)

func TestDeviceRepository_CRUD(t *testing.T) {
    env := integration.NewTestEnv(t)
    repo := devicepg.NewDeviceRepository(env.DB)
    ctx := context.Background()

    tenantID := uuid.New()
    device, err := entity.NewDevice(
        tenantID, uuid.New(), uuid.New(), uuid.Nil,
        "SN-INT-001", vo.DeviceTypeSensor, vo.ProtocolMQTT, "Test", uuid.New(),
    )
    require.NoError(t, err)

    // --- Save ---
    require.NoError(t, repo.Save(ctx, device))

    // --- FindByID ---
    found, err := repo.FindByID(ctx, tenantID, device.ID)
    require.NoError(t, err)
    assert.Equal(t, device.ID, found.ID)
    assert.Equal(t, device.SerialNo, found.SerialNo)
    assert.Equal(t, device.Type, found.Type)

    // --- FindBySerial ---
    found2, err := repo.FindBySerial(ctx, tenantID, device.SerialNo)
    require.NoError(t, err)
    assert.Equal(t, device.ID, found2.ID)

    // --- Update ---
    device.Name = "Updated Name"
    device.Status = vo.DeviceStatusOnline
    require.NoError(t, repo.Update(ctx, device))

    updated, _ := repo.FindByID(ctx, tenantID, device.ID)
    assert.Equal(t, "Updated Name", updated.Name)
    assert.Equal(t, vo.DeviceStatusOnline, updated.Status)

    // --- Delete ---
    require.NoError(t, repo.Delete(ctx, tenantID, device.ID))
    _, err = repo.FindByID(ctx, tenantID, device.ID)
    assert.ErrorIs(t, err, domainerrors.ErrDeviceNotFound)
}

func TestDeviceRepository_TenantIsolation(t *testing.T) {
    env := integration.NewTestEnv(t)
    repo := devicepg.NewDeviceRepository(env.DB)
    ctx := context.Background()

    tenantA := uuid.New()
    tenantB := uuid.New()

    dA, _ := entity.NewDevice(tenantA, uuid.New(), uuid.New(), uuid.Nil,
        "SN-A-001", vo.DeviceTypeSensor, vo.ProtocolMQTT, "A", uuid.New())
    dB, _ := entity.NewDevice(tenantB, uuid.New(), uuid.New(), uuid.Nil,
        "SN-B-001", vo.DeviceTypeSensor, vo.ProtocolMQTT, "B", uuid.New())

    require.NoError(t, repo.Save(ctx, dA))
    require.NoError(t, repo.Save(ctx, dB))

    // Tenant A ไม่เห็น B
    _, err := repo.FindByID(ctx, tenantA, dB.ID)
    assert.ErrorIs(t, err, domainerrors.ErrDeviceNotFound)

    // Tenant B ไม่เห็น A
    _, err = repo.FindByID(ctx, tenantB, dA.ID)
    assert.ErrorIs(t, err, domainerrors.ErrDeviceNotFound)
}
```

### C.3 Kafka Integration Test

```go
// internal/modules/device/infrastructure/messaging/consumers/telemetry_consumer_test.go
//go:build integration

func TestTelemetryConsumer_EndToEnd(t *testing.T) {
    env := integration.NewTestEnv(t)
    ctx := context.Background()

    // Setup producer
    producer, _ := kafka.NewProducer([]string{env.Kafka.Broker()}, "test", nil)

    // Setup consumer
    consumed := make(chan []byte, 1)
    handler := func(ctx context.Context, msg *sarama.ConsumerMessage) error {
        consumed <- msg.Value
        return nil
    }
    group, _ := kafka.NewConsumerGroup(
        []string{env.Kafka.Broker()},
        "test-group",
        []string{"test.telemetry"},
        handler, nil,
    )
    go group.Run(ctx)
    defer group.Close()

    // Publish
    payload := []byte(`{"device_id":"...","metrics":[{"metric":"temp","value":25}]}`)
    err := producer.Publish(ctx, &testEvent{payload})
    require.NoError(t, err)

    // Wait for consumption
    select {
    case msg := <-consumed:
        assert.Equal(t, payload, msg)
    case <-time.After(10 * time.Second):
        t.Fatal("message not consumed in time")
    }
}
```

### C.4 Redis Cache Test

```go
//go:build integration
func TestCache_SetGetDelete(t *testing.T) {
    env := integration.NewTestEnv(t)
    cache := redis.NewCache(env.RedisCli)
    ctx := context.Background()

    key := "test:key:1"
    value := []byte(`{"hello":"world"}`)

    // Set
    require.NoError(t, cache.Set(ctx, key, value, time.Minute))

    // Get
    got, err := cache.Get(ctx, key)
    require.NoError(t, err)
    assert.Equal(t, value, got)

    // Exists
    ok, _ := cache.Exists(ctx, key)
    assert.True(t, ok)

    // Delete
    require.NoError(t, cache.Delete(ctx, key))
    got, _ = cache.Get(ctx, key)
    assert.Nil(t, got)
}
```

---

## 🅳 PART 12D — E2E / API TESTS

### D.1 Newman CI Pipeline

ดูรายละเอียดใน PART 10F — GitHub Actions workflow

### D.2 E2E Smoke Script

ดูรายละเอียดใน PART 9D.2 — `test/e2e/smoke.sh`

### D.3 Playwright (สำหรับ UI)

```typescript
// test/e2e/ui/dashboard.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    // Login
    await page.goto('/login');
    await page.fill('[name=email]', 'admin@demo.local');
    await page.fill('[name=password]', 'Password123!');
    await page.click('[type=submit]');
    await expect(page).toHaveURL('/dashboard');
  });

  test('shows KPI cards', async ({ page }) => {
    await expect(page.locator('[data-testid=kpi-mrr]')).toBeVisible();
    await expect(page.locator('[data-testid=kpi-customers]')).toBeVisible();
    await expect(page.locator('[data-testid=kpi-devices]')).toBeVisible();
  });

  test('realtime telemetry via WebSocket', async ({ page }) => {
    await page.goto('/devices/f0000000-0000-0000-0000-000000000001');
    const value = page.locator('[data-testid=telemetry-temperature]');
    const initial = await value.textContent();

    // รอ 10 วินาที ให้ WS อัปเดต
    await page.waitForTimeout(10000);
    const updated = await value.textContent();
    expect(updated).not.toBe(initial);
  });
});
```

---

## 🅴 PART 12E — CONTRACT TESTS

### E.1 Pact (Consumer-Driven Contract)

```go
// test/contract/pact/device_consumer_test.go
package pact

func TestDeviceAPIConsumer(t *testing.T) {
    pact := &pactgo.Pact{
        Consumer: "web-dashboard",
        Provider: "device-api",
        Host:     "127.0.0.1",
        Port:     6666,
    }
    defer pact.Teardown()

    // Expect GET /api/v1/devices/:id
    pact.
        AddInteraction().
        Given("device SN-DEMO-001 exists").
        UponReceiving("a request for device details").
        WithRequest(dsl.Request{
            Method: "GET",
            Path:   dsl.String("/api/v1/devices/f0000000-0000-0000-0000-000000000001"),
            Headers: dsl.MapMatcher{
                "Authorization": dsl.String("Bearer ..."),
                "X-Tenant-ID":   dsl.String("11111111-..."),
            },
        }).
        WillRespondWith(dsl.Response{
            Status: 200,
            Body: dsl.MapMatcher{
                "data": dsl.MapMatcher{
                    "id":       dsl.String("f0000000-..."),
                    "serial_no": dsl.String("SN-DEMO-001"),
                    "name":     dsl.String("Test Device"),
                    "type":     dsl.String("SENSOR"),
                    "status":   dsl.String("ONLINE"),
                },
            },
        })

    // Verify
    require.NoError(t, pact.Verify(func() error {
        // call API จริง
        return nil
    }))
}
```

---

## 🅵 PART 12F — LOAD & PERFORMANCE TESTS

### F.1 k6 Script

```javascript
// test/load/k6/api_load.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '1m', target: 50 },   // ramp-up
        { duration: '5m', target: 200 },  // steady
        { duration: '1m', target: 0 },    // ramp-down
    ],
    thresholds: {
        http_req_duration: ['p(95)<500', 'p(99)<1000'],
        http_req_failed:   ['rate<0.01'],
    },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const TOKEN    = __ENV.ACCESS_TOKEN;
const TENANT   = __ENV.TENANT_ID;

export default function () {
    const headers = {
        'Authorization': `Bearer ${TOKEN}`,
        'X-Tenant-ID': TENANT,
        'Content-Type': 'application/json',
    };

    // List devices
    const r = http.get(`${BASE_URL}/api/v1/devices?page=1&page_size=20`, { headers });
    check(r, {
        'status 200':     (r) => r.status === 200,
        'has data':       (r) => JSON.parse(r.body).data.items.length > 0,
    });

    sleep(1);
}
```

### F.2 Vegeta (Constant Load)

```bash
# Test API ด้วย rate คงที่
echo "GET http://localhost:8080/api/v1/devices" | \
  vegeta attack -rate=100/s -duration=60s \
    -header="Authorization: Bearer $TOKEN" \
    -header="X-Tenant-ID: $TENANT" \
    | vegeta report

# ดู latency distribution
vegeta attack ... | vegeta plot > plot.html
```

### F.3 Load Test Scenarios

| Scenario | Load | Duration | Success Criteria |
|:---|:---|:---|:---|
| **Smoke** | 10 users | 1 min | p95 < 500ms |
| **Load** | 200 users | 5 min | p95 < 500ms, error < 1% |
| **Stress** | 500 → 2000 users | 15 min | ระบบไม่ crash |
| **Spike** | 0 → 1000 → 0 | 5 min | ฟื้นตัวใน 30s |
| **Soak** | 200 users | 4 hours | ไม่มี memory leak |
| **Telemetry** | 10k msg/s | 10 min | Kafka lag < 5s |
| **WS Broadcast** | 1000 concurrent WS | 10 min | p95 < 100ms |

### F.4 Telemetry Ingestion Load Test

```go
// test/load/telemetry_load.go
//go:build load
package main

func main() {
    // Simulate 10,000 devices ส่ง telemetry ทุก 5 วินาที
    devices := 10000
    interval := 5 * time.Second
    metricsPerDevice := 3

    producer := setupMQTTProducer()
    var wg sync.WaitGroup

    for i := 0; i < devices; i++ {
        wg.Add(1)
        go func(deviceID int) {
            defer wg.Done()
            ticker := time.NewTicker(interval)
            defer ticker.Stop()
            for range ticker.C {
                payload := generateTelemetry(deviceID, metricsPerDevice)
                topic := fmt.Sprintf("iot/SN-LOAD-%05d/telemetry", deviceID)
                _ = producer.Publish(topic, 1, false, payload)
            }
        }(i)
    }
    wg.Wait()
}
```

---

## 🅶 PART 12G — TEST FIXTURES & BUILDERS

ดูรายละเอียดทั้งหมดใน PART 9 — Sample Data & Fixtures (SQL + Go Builders + Loader)

### G.1 Test Fixture Pattern

```go
// internal/modules/device/application/testutil/builders.go
package testutil

type DeviceBuilder struct {
    tenantID uuid.UUID
    serial   string
    // ...
}

func NewDeviceBuilder() *DeviceBuilder {
    return &DeviceBuilder{
        tenantID: uuid.New(),
        serial:   "SN-TEST-001",
        // ...
    }
}

func (b *DeviceBuilder) WithTenant(id uuid.UUID) *DeviceBuilder {
    b.tenantID = id
    return b
}

func (b *DeviceBuilder) Build(t *testing.T) *entity.Device {
    d, err := entity.NewDevice(b.tenantID, ...)
    require.NoError(t, err)
    return d
}
```

---

## 🅷 PART 12H — CI TEST PIPELINE

### H.1 GitHub Actions (Complete)

```yaml
name: Test

on:
  push: { branches: [main, develop] }
  pull_request: { branches: [main, develop] }

jobs:
  # 1. Unit tests (fast, ทุก push)
  unit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.23' }
      - run: go mod download
      - name: Unit tests
        run: go test -short -race -coverprofile=coverage.out -covermode=atomic ./...
      - name: Coverage check
        run: |
          coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          if (( $(echo "$coverage < 80" | bc -l) )); then
            echo "Coverage $coverage% < 80%"
            exit 1
          fi
      - uses: codecov/codecov-action@v4
        with: { files: coverage.out }

  # 2. Lint & vet
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.23' }
      - run: go vet ./...
      - uses: golangci/golangci-lint-action@v6
        with: { version: v1.60 }

  # 3. Integration tests (testcontainers)
  integration:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.23' }
      - run: go test -tags=integration -timeout=15m ./test/integration/...

  # 4. E2E (Newman)
  e2e:
    runs-on: ubuntu-latest
    needs: [unit, integration]
    services:
      postgres: { image: postgres:16-alpine, env: {...} }
      redis:    { image: redis:7-alpine }
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.23' }
      - uses: actions/setup-node@v4
        with: { node-version: '20' }
      - run: npm i -g newman newman-reporter-htmlextra
      - run: go run ./cmd/migrate -action up
      - run: bash test/fixtures/load.sh
      - run: go build -o /tmp/api ./cmd/api && /tmp/api &
      - run: newman run postman/collection.json -e postman/env.dev.json --bail

  # 5. Security scan
  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: |
          go install github.com/securego/gosec/v2/cmd/gosec@latest
          gosec -fmt sarif -out results.sarif ./...
      - uses: github/codeql-action/upload-sarif@v3
        with: { sarif_file: results.sarif }
```

### H.2 Makefile Targets

```makefile
.PHONY: test
test: ## unit tests
	go test -short -race ./...

.PHONY: test-cover
test-cover: ## unit + coverage
	go test -short -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: test-integration
test-integration: ## integration tests (testcontainers)
	go test -tags=integration -timeout=15m ./test/integration/...

.PHONY: test-e2e
test-e2e: ## e2e via newman
	newman run postman/collection.json -e postman/env.dev.json --bail

.PHONY: test-load
test-load: ## load test (k6)
	k6 run test/load/k6/api_load.js

.PHONY: test-all
test-all: test test-integration test-e2e
	@echo "✓ all tests passed"

.PHONY: bench
bench: ## benchmarks
	go test -bench=. -benchmem ./...
```

---

## 🎯 PART 12 — SUMMARY

| Type | Count Target | Speed | Runs |
|---|:-:|:-:|:-:|
| **Unit** | 700+ | <1s | Every commit |
| **Integration** | 250+ | 5-30s | PR |
| **E2E API** | 50+ | 1-5min | PR |
| **E2E UI** | 20+ | 5-10min | Nightly |
| **Contract** | 15+ | <1min | PR |
| **Load** | 7 scenarios | 10-60min | Weekly |

---
