
---

# 🔒 PART 14 — SECURITY DEEP DIVE

> **ขนาด**: ใหญ่ — แยก 8 ตอนย่อย
> **Part 14A**: Security Architecture & Threat Model
> **Part 14B**: Authentication (JWT + Refresh + MFA)
> **Part 14C**: Authorization (RBAC + ABAC + RLS)
> **Part 14D**: Encryption (at rest + in transit)
> **Part 14E**: Input Validation & OWASP Top 10
> **Part 14F**: Device Security (IoT-specific)
> **Part 14G**: PDPA / GDPR Compliance
> **Part 14H**: Security Operations (SIEM, Pen Test, Incident)

---

## 🅰️ PART 14A — SECURITY ARCHITECTURE & THREAT MODEL

### A.1 Defense in Depth

```
┌─────────────────────────────────────────────────────────────┐
│  Layer 1: Edge (Cloudflare)                                 │
│  - DDoS protection · WAF · Rate limit · Bot detection       │
├─────────────────────────────────────────────────────────────┤
│  Layer 2: Network                                           │
│  - VPC · Security Groups · Network Policies · mTLS          │
├─────────────────────────────────────────────────────────────┤
│  Layer 3: Application                                       │
│  - Auth · RBAC · Input validation · Rate limit · CSP        │
├─────────────────────────────────────────────────────────────┤
│  Layer 4: Data                                              │
│  - Encryption at rest · RLS · Audit log · Backup            │
├─────────────────────────────────────────────────────────────┤
│  Layer 5: Monitoring                                        │
│  - SIEM · Anomaly detection · Alert · Incident response     │
└─────────────────────────────────────────────────────────────┘
```

### A.2 STRIDE Threat Model

| Threat | ตัวอย่าง | Mitigation |
|:---|:---|:---|
| **Spoofing** | ปลอม device token | mTLS, token rotation, HMAC |
| **Tampering** | แก้ payload ระหว่างทาง | TLS, signature |
| **Repudiation** | ปฏิเสธว่าไม่ได้ทำ | Audit log + signature |
| **Information Disclosure** | leak PII | Encryption, RLS, PDPA |
| **Denial of Service** | flood API | Rate limit, CDN, WAF |
| **Elevation of Privilege** | user → admin | RBAC, validation |

### A.3 Attack Surface

| Surface | Risk | Control |
|:---|:---|:---|
| **REST API** | High | Auth + rate limit + WAF |
| **WebSocket** | Medium | JWT + origin check |
| **MQTT broker** | High | mTLS + ACL + token |
| **Admin CLI** | Medium | SSH key + IP allowlist |
| **Database** | High | VPC only + TLS + RLS |
| **Kafka** | Medium | SASL + TLS + ACL |
| **S3 backups** | Medium | KMS + versioning + MFA delete |

---

## 🅱️ PART 14B — AUTHENTICATION

### B.1 JWT Structure

```json
// Header
{
    "alg": "HS256",  // หรือ RS256 สำหรับ production
    "typ": "JWT",
    "kid": "key-2026-01"
}

// Payload
{
    "iss": "icmongolang",
    "sub": "user-uuid",
    "aud": "icmon-api",
    "exp": 1737100800,
    "iat": 1737099900,
    "nbf": 1737099900,
    "jti": "token-uuid",
    "tenant_id": "tenant-uuid",
    "email": "user@example.com",
    "role": "ADMIN",
    "scopes": ["devices:read", "devices:write"],
    "amr": ["pwd"],
    "session_id": "session-uuid"
}
```

### B.2 JWT Service (Production: RS256)

```go
// internal/modules/auth/application/jwt_service.go
package application

import (
    "crypto/rsa"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

type JWTService struct {
    privateKey *rsa.PrivateKey
    publicKey  *rsa.PublicKey
    issuer     string
    accessTTL  time.Duration
    refreshTTL time.Duration
}

type Claims struct {
    TenantID  uuid.UUID `json:"tenant_id"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    Scopes    []string  `json:"scopes"`
    SessionID uuid.UUID `json:"session_id"`
    jwt.RegisteredClaims
}

func NewJWTService(privateKeyPEM, publicKeyPEM []byte, issuer string, accessTTL, refreshTTL time.Duration) (*JWTService, error) {
    privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
    if err != nil {
        return nil, err
    }
    publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyPEM)
    if err != nil {
        return nil, err
    }
    return &JWTService{
        privateKey: privateKey, publicKey: publicKey,
        issuer: issuer, accessTTL: accessTTL, refreshTTL: refreshTTL,
    }, nil
}

func (s *JWTService) IssueAccessToken(userID uuid.UUID, tenantID uuid.UUID, email, role string, scopes []string, sessionID uuid.UUID) (string, error) {
    now := time.Now().UTC()
    claims := Claims{
        TenantID:  tenantID,
        Email:     email,
        Role:      role,
        Scopes:    scopes,
        SessionID: sessionID,
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    s.issuer,
            Subject:   userID.String(),
            Audience:  jwt.ClaimStrings{"icmon-api"},
            ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
            IssuedAt:  jwt.NewNumericDate(now),
            NotBefore: jwt.NewNumericDate(now),
            ID:        uuid.NewString(),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    token.Header["kid"] = "key-2026-01"
    return token.SignedString(s.privateKey)
}

func (s *JWTService) Validate(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }
        return s.publicKey, nil
    }, jwt.WithIssuer(s.issuer), jwt.WithAudience("icmon-api"), jwt.WithExpirationRequired())

    if err != nil {
        return nil, err
    }
    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }
    return claims, nil
}
```

### B.3 Refresh Token (Rotating)

```go
// internal/modules/auth/domain/entity/refresh_token.go
package entity

type RefreshToken struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    TenantID  uuid.UUID
    TokenHash string
    SessionID uuid.UUID
    ExpiresAt time.Time
    RevokedAt *time.Time
    ReplacedBy *uuid.UUID  // token rotation chain
    UserAgent string
    IP        string
    CreatedAt time.Time
}

func (t *RefreshToken) IsValid() bool {
    return t.RevokedAt == nil && time.Now().UTC().Before(t.ExpiresAt)
}

func (t *RefreshToken) Revoke(replacedBy uuid.UUID) {
    now := time.Now().UTC()
    t.RevokedAt = &now
    t.ReplacedBy = &replacedBy
}
```

```go
// internal/modules/auth/application/refresh_token.go
func (uc *RefreshTokenUseCase) Execute(ctx context.Context, in RefreshInput) (*RefreshOutput, error) {
    // 1. Hash token
    hash := sha256.Sum256([]byte(in.RefreshToken))
    hashHex := hex.EncodeToString(hash[:])

    // 2. Find token
    existing, err := uc.repo.FindByHash(ctx, hashHex)
    if err != nil {
        return nil, domainerrors.ErrInvalidToken
    }
    if !existing.IsValid() {
        // Reuse detection: ถ้า token ถูก revoke แล้ว → revoke ทั้ง session (token theft)
        if existing.RevokedAt != nil {
            _ = uc.repo.RevokeAllBySession(ctx, existing.SessionID)
            uc.logger.Warn("refresh token reuse detected",
                application.Field{Key: "user_id", Value: existing.UserID},
                application.Field{Key: "session_id", Value: existing.SessionID},
            )
        }
        return nil, domainerrors.ErrRefreshTokenRevoked
    }

    // 3. Load user
    user, err := uc.userRepo.FindByID(ctx, existing.UserID)
    if err != nil {
        return nil, domainerrors.ErrNotFound
    }
    if user.Status != "ACTIVE" {
        return nil, domainerrors.ErrUserDisabled
    }

    // 4. Issue new tokens
    newAccess, _ := uc.jwt.IssueAccessToken(user.ID, existing.TenantID, user.Email, user.Role, user.Scopes, existing.SessionID)
    newRefresh, _ := uc.generateRefreshToken()

    // 5. Rotate: revoke old, save new
    newRefreshHash := sha256.Sum256([]byte(newRefresh))
    newToken := &entity.RefreshToken{
        ID:        uuid.New(),
        UserID:    user.ID,
        TenantID:  existing.TenantID,
        TokenHash: hex.EncodeToString(newRefreshHash[:]),
        SessionID: existing.SessionID,
        ExpiresAt: time.Now().UTC().Add(uc.jwt.RefreshTTL()),
        UserAgent: in.UserAgent,
        IP:        in.IP,
    }
    if err := uc.repo.Rotate(ctx, existing, newToken); err != nil {
        return nil, err
    }

    return &RefreshOutput{
        AccessToken:  newAccess,
        RefreshToken: newRefresh,
        ExpiresAt:    time.Now().UTC().Add(uc.jwt.AccessTTL()),
    }, nil
}
```

### B.4 Password Hashing (bcrypt with cost)

```go
// internal/modules/auth/infrastructure/crypto/password.go
package crypto

import (
    "golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

func HashPassword(password string) (string, error) {
    if len(password) < 12 {
        return "", domainerrors.ErrWeakPassword
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

func VerifyPassword(hash, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
```

### B.5 MFA (TOTP)

```go
// internal/modules/auth/infrastructure/mfa/totp.go
package mfa

import (
    "github.com/pquerna/otp/totp"
)

func GenerateSecret(email string) (string, string, error) {
    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      "icmongolang",
        AccountName: email,
    })
    if err != nil {
        return "", "", err
    }
    return key.Secret(), key.URL(), nil
}

func Validate(secret, code string) bool {
    return totp.Validate(code, secret)
}
```

### B.6 Login Use Case (Complete)

```go
// internal/modules/auth/application/login.go
func (uc *LoginUseCase) Execute(ctx context.Context, in LoginInput) (*LoginOutput, error) {
    // 1. Rate limit per email/IP (ป้องกัน brute force)
    allowed, err := uc.rateLimiter.Allow(ctx, "login:"+in.Email, 5, 15*time.Minute)
    if err != nil || !allowed {
        return nil, domainerrors.ErrRateLimitExceeded
    }

    // 2. Find user
    user, err := uc.userRepo.FindByEmail(ctx, in.TenantID, in.Email)
    if err != nil {
        // Timing attack protection: hash dummy
        _ = bcrypt.CompareHashAndPassword([]byte("$2a$12$dummy"), []byte(in.Password))
        return nil, domainerrors.ErrInvalidCredentials
    }

    // 3. Check status
    if user.Status != "ACTIVE" {
        return nil, domainerrors.ErrUserDisabled
    }

    // 4. Check locked
    if user.LockedUntil != nil && time.Now().UTC().Before(*user.LockedUntil) {
        return nil, domainerrors.ErrAccountLocked
    }

    // 5. Verify password
    if err := uc.passwordSvc.Verify(user.PasswordHash, in.Password); err != nil {
        user.IncrementFailedAttempts()
        if user.FailedAttempts >= 5 {
            user.Lock(15 * time.Minute)
        }
        _ = uc.userRepo.Update(ctx, user)
        return nil, domainerrors.ErrInvalidCredentials
    }

    // 6. Check MFA
    if user.MFAEnabled {
        if in.MFACode == "" {
            return &LoginOutput{MFARequired: true}, nil
        }
        if !mfa.Validate(user.MFASecret, in.MFACode) {
            return nil, domainerrors.ErrInvalidMFACode
        }
    }

    // 7. Reset failed attempts
    user.ResetFailedAttempts()
    user.LastLoginAt = ptr(time.Now().UTC())
    _ = uc.userRepo.Update(ctx, user)

    // 8. Create session
    session := &entity.Session{
        ID:       uuid.New(),
        UserID:   user.ID,
        TenantID: user.TenantID,
        IP:       in.IP,
        UserAgent: in.UserAgent,
        CreatedAt: time.Now().UTC(),
        ExpiresAt: time.Now().UTC().Add(uc.jwt.RefreshTTL()),
    }
    _ = uc.sessionRepo.Save(ctx, session)

    // 9. Issue tokens
    accessToken, _ := uc.jwt.IssueAccessToken(user.ID, user.TenantID, user.Email, user.Role, user.Scopes, session.ID)
    refreshToken, _ := uc.generateRefresh()
    refreshHash := sha256.Sum256([]byte(refreshToken))

    rt := &entity.RefreshToken{
        ID:        uuid.New(),
        UserID:    user.ID,
        TenantID:  user.TenantID,
        TokenHash: hex.EncodeToString(refreshHash[:]),
        SessionID: session.ID,
        ExpiresAt: time.Now().UTC().Add(uc.jwt.RefreshTTL()),
        UserAgent: in.UserAgent,
        IP:        in.IP,
    }
    _ = uc.refreshRepo.Save(ctx, rt)

    // 10. Audit
    _ = uc.auditRepo.Save(ctx, entity.NewAuditTrail(user.ID, "LOGIN_SUCCESS", map[string]interface{}{
        "ip": in.IP, "user_agent": in.UserAgent,
    }))

    return &LoginOutput{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresAt:    time.Now().UTC().Add(uc.jwt.AccessTTL()),
        User:         ToUserDTO(user),
    }, nil
}
```

---

## 🅲 PART 14C — AUTHORIZATION

### C.1 RBAC Model

```go
// internal/modules/auth/domain/entity/role.go
package entity

type Role string

const (
    RoleAdmin   Role = "ADMIN"
    RoleManager Role = "MANAGER"
    RoleUser    Role = "USER"
    RoleViewer  Role = "VIEWER"
)

type Permission string

const (
    PermDeviceRead   Permission = "devices:read"
    PermDeviceWrite  Permission = "devices:write"
    PermDeviceDelete Permission = "devices:delete"
    PermCustomerRead Permission = "customers:read"
    PermCustomerWrite Permission = "customers:write"
    PermOrderRead    Permission = "orders:read"
    PermOrderWrite   Permission = "orders:write"
    PermInvoiceRead  Permission = "invoices:read"
    PermInvoiceWrite Permission = "invoices:write"
    PermReportRead   Permission = "reports:read"
    PermAdminAll     Permission = "*:*"
)

var rolePermissions = map[Role][]Permission{
    RoleAdmin:   {PermAdminAll},
    RoleManager: {
        PermDeviceRead, PermDeviceWrite,
        PermCustomerRead, PermCustomerWrite,
        PermOrderRead, PermOrderWrite,
        PermInvoiceRead, PermInvoiceWrite,
        PermReportRead,
    },
    RoleUser: {
        PermDeviceRead, PermDeviceWrite,
        PermCustomerRead,
        PermOrderRead,
        PermReportRead,
    },
    RoleViewer: {
        PermDeviceRead, PermCustomerRead, PermOrderRead, PermInvoiceRead, PermReportRead,
    },
}

func (r Role) HasPermission(p Permission) bool {
    perms, ok := rolePermissions[r]
    if !ok {
        return false
    }
    for _, perm := range perms {
        if perm == PermAdminAll || perm == p {
            return true
        }
    }
    return false
}
```

### C.2 Permission Middleware

```go
// internal/shared/interfaces/http/middleware/permission.go
package middleware

func RequirePermission(perm entity.Permission) gin.HandlerFunc {
    return func(c *gin.Context) {
        role := entity.Role(c.GetString("user_role"))
        if !role.HasPermission(perm) {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "error": gin.H{
                    "code":    "FORBIDDEN",
                    "message": "insufficient permissions",
                    "details": gin.H{"required": string(perm), "role": string(role)},
                },
            })
            return
        }
        c.Next()
    }
}

// Usage
g.POST("/devices", middleware.RequirePermission(entity.PermDeviceWrite), h.Register)
```

### C.3 Row-Level Security (PostgreSQL RLS)

```sql
-- Enable RLS
ALTER TABLE device_devices ENABLE ROW LEVEL SECURITY;

-- Policy: tenant isolation
CREATE POLICY tenant_isolation ON device_devices
    USING (tenant_id = current_setting('app.current_tenant')::uuid)
    WITH CHECK (tenant_id = current_setting('app.current_tenant')::uuid);

-- Force RLS (บังคับแม้แต่ owner)
ALTER TABLE device_devices FORCE ROW LEVEL SECURITY;

-- Set tenant context ที่ middleware
```

```go
// internal/shared/infrastructure/postgres/rls.go
package postgres

func SetTenantContext(ctx context.Context, db *gorm.DB, tenantID uuid.UUID) error {
    return db.WithContext(ctx).Exec("SET app.current_tenant = ?", tenantID.String()).Error
}
```

```go
// Middleware: set tenant ก่อน query
func (m *TenantMiddleware) WithRLS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tenantID := extractTenant(r)
        ctx := context.WithValue(r.Context(), "tenant_id", tenantID)
        // Set RLS context
        _ = m.db.Exec("SET LOCAL app.current_tenant = ?", tenantID).Error
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### C.4 Attribute-Based Access Control (ABAC)

```go
// internal/modules/device/application/abac.go
type Policy interface {
    Evaluate(ctx context.Context, principal Principal, resource Resource, action string) (bool, error)
}

type Principal struct {
    UserID   uuid.UUID
    TenantID uuid.UUID
    Role     string
    Scopes   []string
    Attributes map[string]interface{}
}

type Resource struct {
    Type      string
    ID        uuid.UUID
    TenantID  uuid.UUID
    OwnerID   uuid.UUID
    Attributes map[string]interface{}
}

// DevicePolicy: ตรวจว่าผู้ใช้มีสิทธิ์กับ device
func (p *DevicePolicy) Evaluate(ctx context.Context, pr Principal, res Resource, action string) (bool, error) {
    // 1. Tenant isolation
    if pr.TenantID != res.TenantID {
        return false, nil
    }
    // 2. Admin ผ่าน
    if pr.Role == "ADMIN" {
        return true, nil
    }
    // 3. Manager: only if device in their site
    if pr.Role == "MANAGER" {
        siteID, _ := pr.Attributes["site_id"].(uuid.UUID)
        deviceSite, _ := res.Attributes["site_id"].(uuid.UUID)
        return siteID == deviceSite, nil
    }
    // 4. User: only own devices
    if pr.Role == "USER" {
        return res.OwnerID == pr.UserID, nil
    }
    return false, nil
}
```

### C.5 API Key Authentication (สำหรับ external API)

```go
// internal/modules/auth/domain/entity/api_key.go
type APIKey struct {
    ID         uuid.UUID
    TenantID   uuid.UUID
    Name       string
    KeyHash    string  // SHA-256
    Prefix     string  // เช่น "icmon_sk_xxx" สำหรับ lookup
    Scopes     []string
    IPAllowlist []string
    RateLimit  int
    ExpiresAt  *time.Time
    LastUsedAt *time.Time
    RevokedAt  *time.Time
    CreatedBy  uuid.UUID
    CreatedAt  time.Time
}

func GenerateAPIKey() (plain, hash, prefix string) {
    random := make([]byte, 32)
    rand.Read(random)
    plain = "icmon_sk_" + base64.URLEncoding.EncodeToString(random)
    prefix = plain[:16]
    h := sha256.Sum256([]byte(plain))
    hash = hex.EncodeToString(h[:])
    return
}
```

```go
// API Key middleware
func APIKeyAuth(svc APIKeyService) gin.HandlerFunc {
    return func(c *gin.Context) {
        header := c.GetHeader("X-API-Key")
        if header == "" {
            c.Next() // fallback ไป JWT
            return
        }
        key, err := svc.Validate(c.Request.Context(), header)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "invalid API key"})
            return
        }
        c.Set("api_key_id", key.ID)
        c.Set("tenant_id", key.TenantID)
        c.Set("scopes", key.Scopes)
        c.Next()
    }
}
```

---

## 🅳 PART 14D — ENCRYPTION

### D.1 Encryption at Rest

| Component | Method | Key Management |
|:---|:---|:---|
| **PostgreSQL** | Transparent Data Encryption (TDE) | KMS |
| **S3 Backup** | SSE-KMS | AWS KMS |
| **InfluxDB** | Disk encryption | KMS |
| **Kafka** | Log encryption | KMS |
| **Sensitive fields** | Application-level (pgcrypto) | HashiCorp Vault |

### D.2 Column-Level Encryption (Sensitive PII)

```go
// internal/shared/infrastructure/crypto/aes.go
package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "io"
)

type Encryptor struct {
    gcm cipher.AEAD
}

func NewEncryptor(key []byte) (*Encryptor, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    return &Encryptor{gcm: gcm}, nil
}

func (e *Encryptor) Encrypt(plaintext []byte) (string, error) {
    nonce := make([]byte, e.gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }
    ciphertext := e.gcm.Seal(nonce, nonce, plaintext, nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (e *Encryptor) Decrypt(encoded string) ([]byte, error) {
    ciphertext, err := base64.StdEncoding.DecodeString(encoded)
    if err != nil {
        return nil, err
    }
    nonceSize := e.gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, errors.New("ciphertext too short")
    }
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    return e.gcm.Open(nil, nonce, ciphertext, nil)
}
```

### D.3 Encryption in Transit

```
- TLS 1.3 only (disable 1.0, 1.1, 1.2)
- HSTS with preload
- Certificate rotation อัตโนมัติ (Let's Encrypt / cert-manager)
- mTLS สำหรับ service-to-service
- MQTT over TLS (port 8883)
- Kafka over TLS + SASL
```

### D.4 TLS Config

```go
// internal/shared/infrastructure/tls/config.go
package tls

func ServerConfig(certFile, keyFile string) (*tls.Config, error) {
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        return nil, err
    }
    return &tls.Config{
        Certificates: []tls.Certificate{cert},
        MinVersion:   tls.VersionTLS13,
        CipherSuites: []uint16{
            tls.TLS_AES_128_GCM_SHA256,
            tls.TLS_AES_256_GCM_SHA384,
            tls.TLS_CHACHA20_POLY1305_SHA256,
        },
        CurvePreferences: []tls.CurveID{
            tls.X25519,
            tls.CurveP256,
        },
        PreferServerCipherSuites: true,
    }, nil
}
```

### D.5 Password Hashing Summary

| Use Case | Algorithm | Params |
|:---|:---|:---|
| **User password** | bcrypt | cost 12 |
| **API key** | SHA-256 | (cryptographic hash) |
| **Device token** | bcrypt | cost 10 |
| **JWT signing** | RS256 | RSA 4096 |
| **Sensitive field** | AES-256-GCM | random nonce |
| **Integrity** | HMAC-SHA256 | 32-byte key |

---

## 🅴 PART 14E — INPUT VALIDATION & OWASP TOP 10

### E.1 Input Validation Layers

```
1. HTTP Handler: binding tags (gin)
2. Application: DTO validation
3. Domain: invariant (constructor + behavior)
4. Infrastructure: parameterized query (GORM)
```

### E.2 OWASP Top 10 (2021) — Protection

| # | Risk | Mitigation |
|:-:|:---|:---|
| **A01** | Broken Access Control | RBAC + RLS + tenant middleware |
| **A02** | Cryptographic Failures | TLS 1.3, AES-256, bcrypt |
| **A03** | Injection | Parameterized queries, no string concat |
| **A04** | Insecure Design | Threat modeling, DDD invariants |
| **A05** | Security Misconfiguration | IaC, CIS benchmarks, helm lint |
| **A06** | Vulnerable Components | Dependabot, Snyk, go mod verify |
| **A07** | Authentication Failures | JWT + MFA + rate limit + rotation |
| **A08** | Data Integrity Failures | Signed tokens, HMAC webhooks |
| **A09** | Logging Failures | Structured logs, SIEM, alerts |
| **A10** | SSRF | Allowlist URLs, no user-controlled URLs |

### E.3 SQL Injection Prevention

```go
// ❌ NEVER do this
db.Raw("SELECT * FROM users WHERE email = '" + email + "'")

// ✅ Parameterized
db.Raw("SELECT * FROM users WHERE email = ?", email).Scan(&user)

// ✅ GORM
db.Where("email = ?", email).First(&user)

// ✅ Struct-based
db.Where(&User{Email: email}).First(&user)
```

### E.4 XSS Prevention

```go
// Server-side: escape output
import "html"
func escapeHTML(s string) string {
    return html.EscapeString(s)
}

// Gin JSON response: safe by default (application/json)
// HTML response: ใช้ template ที่ escape อัตโนมัติ
```

### E.5 CSRF Protection

```go
// API ใช้ JWT ใน Authorization header → ไม่ต้อง CSRF
// แต่ถ้าใช้ cookie → ต้อง CSRF token
func CSRFMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.Method == "GET" || c.Request.Method == "HEAD" {
            c.Next()
            return
        }
        // Only for cookie-based auth
        if c.GetHeader("Authorization") != "" {
            c.Next()
            return
        }
        token := c.GetHeader("X-CSRF-Token")
        cookieToken, _ := c.Cookie("csrf_token")
        if token == "" || token != cookieToken {
            c.AbortWithStatusJSON(403, gin.H{"error": "CSRF token mismatch"})
            return
        }
        c.Next()
    }
}
```

### E.6 SSRF Prevention

```go
// Allowlist internal services
var allowedHosts = map[string]bool{
    "api.internal.svc":    true,
    "webhook.stripe.com":  true,
    "api.openai.com":      true,
}

func validateURL(rawURL string) error {
    u, err := url.Parse(rawURL)
    if err != nil {
        return err
    }
    // Block internal IPs
    if isInternalIP(u.Hostname()) {
        return errors.New("internal IP not allowed")
    }
    // Only HTTPS
    if u.Scheme != "https" {
        return errors.New("HTTPS required")
    }
    // Allowlist (optional)
    if !allowedHosts[u.Host] {
        return errors.New("host not allowed")
    }
    return nil
}
```

### E.7 Security Headers (Full)

```go
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; connect-src 'self' wss://api.icmongolang.io")
        c.Header("Permissions-Policy", "geolocation=(self), camera=(), microphone=()")
        c.Header("Cross-Origin-Opener-Policy", "same-origin")
        c.Header("Cross-Origin-Resource-Policy", "same-origin")
        c.Next()
    }
}
```

### E.8 Request Size Limits

```go
r.Use(func(c *gin.Context) {
    c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20) // 10MB
    c.Next()
})
```

---

## 🅵 PART 14F — DEVICE SECURITY

### F.1 Device Provisioning Flow

```
1. Admin registers device (serial, type, site)
   ↓
2. Cloud generates device_token (32 bytes random)
   ↓
3. Hash token ด้วย bcrypt → เก็บใน DB
   ↓
4. Return plaintext token ให้ admin (ครั้งเดียว!)
   ↓
5. Admin flash token ลง device (secure element ถ้ามี)
   ↓
6. Device connect MQTT ด้วย username=serial, password=token
   ↓
7. Broker verify กับ cloud (auth webhook)
   ↓
8. Cloud คืน ACL + topic allowlist
   ↓
9. Device publish telemetry
```

### F.2 Device Authentication

```go
// MQTT auth webhook (EMQX)
func (h *MQTTAuthHandler) Authenticate(c *gin.Context) {
    var req struct {
        ClientID string `json:"clientid"`
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"result": "deny"})
        return
    }

    // Username = serial, password = token
    device, err := h.deviceRepo.FindBySerial(c.Request.Context(), uuid.Nil, vo.DeviceSerial(req.Username))
    if err != nil {
        c.JSON(200, gin.H{"result": "deny"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(device.DeviceTokenHash), []byte(req.Password)); err != nil {
        c.JSON(200, gin.H{"result": "deny"})
        return
    }

    if device.Status == vo.DeviceStatusFault {
        c.JSON(200, gin.H{"result": "deny", "reason": "device in fault"})
        return
    }

    // Success
    c.JSON(200, gin.H{
        "result":      "allow",
        "is_superuser": false,
        "client_attrs": gin.H{
            "tenant_id": device.TenantID.String(),
            "device_id": device.ID.String(),
        },
    })
}
```

### F.3 MQTT ACL (Topic Authorization)

```go
func (h *MQTTAuthHandler) ACL(c *gin.Context) {
    var req struct {
        Username string `json:"username"`
        Topic    string `json:"topic"`
        Action   string `json:"action"` // publish, subscribe
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(200, gin.H{"result": "deny"})
        return
    }

    // ตรวจว่า topic ตรงกับ pattern ที่อนุญาต
    // Publish: iot/{serial}/telemetry, iot/{serial}/heartbeat, ...
    // Subscribe: iot/{serial}/command, device/{serial}/shadow/desired
    switch req.Action {
    case "publish":
        allowed := strings.HasPrefix(req.Topic, "iot/"+req.Username+"/")
        if allowed {
            c.JSON(200, gin.H{"result": "allow"})
            return
        }
    case "subscribe":
        allowed := strings.HasPrefix(req.Topic, "iot/"+req.Username+"/command") ||
                   strings.HasPrefix(req.Topic, "device/"+req.Username+"/shadow")
        if allowed {
            c.JSON(200, gin.H{"result": "allow"})
            return
        }
    }
    c.JSON(200, gin.H{"result": "deny"})
}
```

### F.4 Device Token Rotation

```go
func (uc *RotateDeviceTokenUseCase) Execute(ctx context.Context, in RotateTokenInput) (*RotateTokenOutput, error) {
    device, err := uc.repo.FindByID(ctx, in.TenantID, in.DeviceID)
    if err != nil {
        return nil, err
    }

    // Generate new token
    newToken := generateSecureToken(32)
    hash, _ := bcrypt.GenerateFromPassword([]byte(newToken), 10)

    if err := device.RotateToken(string(hash), in.ActorID); err != nil {
        return nil, err
    }
    if err := uc.repo.Update(ctx, device); err != nil {
        return nil, err
    }

    // Publish event → disconnect device → บังคับ reconnect ด้วย token ใหม่
    _ = uc.producer.PublishBatch(ctx, device.PullEvents())

    return &RotateTokenOutput{
        NewToken: newToken,  // return ครั้งเดียว!
    }, nil
}
```

### F.5 Firmware Signing

```go
// Firmware upload + signature
type Firmware struct {
    ID          uuid.UUID
    Version     string
    Channel     string  // stable, beta, canary
    URL         string
    SHA256      string  // hash
    Signature   string  // RSA-4096 signature
    Size        int64
    ReleaseNotes string
    CreatedAt   time.Time
}

func (s *FirmwareService) Verify(fw *Firmware, content []byte) error {
    // 1. ตรวจ hash
    h := sha256.Sum256(content)
    if hex.EncodeToString(h[:]) != fw.SHA256 {
        return errors.New("hash mismatch")
    }
    // 2. ตรวจ signature
    h2 := sha256.Sum256(content)
    if err := rsa.VerifyPKCS1v15(s.publicKey, crypto.SHA256, h2[:], []byte(fw.Signature)); err != nil {
        return errors.New("signature invalid")
    }
    return nil
}
```

### F.6 IoT Security Best Practices

```
✓ Device token ≥ 256 bits random
✓ Token หมุนทุก 90 วัน
✓ mTLS สำหรับ device ที่รองรับ
✓ Secure boot + secure element (ถ้ามี)
✓ OTA update ต้อง signed
✓ Rate limit ต่อ device (กัน flood)
✓ Anomaly detection (unusual traffic)
✓ Firmware version pinning
✓ Network segmentation (IoT VLAN)
✓ Disable unused ports/services
```

---

## 🅶 PART 14G — PDPA / GDPR COMPLIANCE

### G.1 PDPA Requirements (ไทย)

| ข้อ | ความต้องการ | Implementation |
|:---|:---|:---|
| **Consent** | ขอความยินยอมก่อนเก็บข้อมูล | `pdpa_consents` table |
| **Purpose Limitation** | ใช้ตามวัตถุประสงค์ที่แจ้ง | Consent scope |
| **Data Minimization** | เก็บเท่าที่จำเป็น | Schema design |
| **Access Right** | ขอดูข้อมูลตัวเองได้ | `/api/v1/pdpa/my-data` |
| **Rectification** | แก้ไขข้อมูลได้ | `/api/v1/pdpa/correct` |
| **Erasure (Right to be Forgotten)** | ลบข้อมูลได้ | `/api/v1/pdpa/erase` |
| **Portability** | export ข้อมูลได้ | `/api/v1/pdpa/export` |
| **Breach Notification** | แจ้งภายใน 72 ชม. | Incident runbook |
| **DPO** | มี Data Protection Officer | (org) |

### G.2 Consent Management

```go
// internal/modules/pdpa/domain/entity/consent.go
type Consent struct {
    ID          uuid.UUID
    TenantID    uuid.UUID
    SubjectID   uuid.UUID  // customer หรือ user
    SubjectType string     // CUSTOMER, USER
    Purpose     string     // MARKETING, ANALYTICS, ESSENTIAL
    Granted     bool
    GrantedAt   time.Time
    RevokedAt   *time.Time
    IP          string
    UserAgent   string
    Version     string     // consent text version
    CreatedAt   time.Time
}
```

```go
// API
POST /api/v1/pdpa/consents           # ให้ consent
GET  /api/v1/pdpa/consents           # ดู consent ทั้งหมด
DELETE /api/v1/pdpa/consents/:purpose # ถอน consent
```

### G.3 Data Subject Access Request (DSAR)

```go
// POST /api/v1/pdpa/dsar
type DSARRequest struct {
    Type      string  // ACCESS, RECTIFY, ERASE, PORTABILITY
    SubjectID string
    Reason    string
}

type DSARService struct {
    userRepo    repository.UserRepository
    customerRepo repository.CustomerRepository
    orderRepo   repository.OrderRepository
    deviceRepo  repository.DeviceRepository
    auditRepo   repository.AuditRepository
}

func (s *DSARService) Access(ctx context.Context, subjectID uuid.UUID) (*DataPackage, error) {
    pkg := &DataPackage{
        ExportedAt: time.Now().UTC(),
        SubjectID:  subjectID,
    }
    // รวบรวมข้อมูลจากทุก module
    if user, err := s.userRepo.FindByID(ctx, subjectID); err == nil {
        pkg.User = toUserExport(user)
    }
    if customer, err := s.customerRepo.FindByUserID(ctx, subjectID); err == nil {
        pkg.Customer = toCustomerExport(customer)
    }
    // ...
    return pkg, nil
}

func (s *DSARService) Erase(ctx context.Context, subjectID uuid.UUID, reason string) error {
    // 1. Anonymize (ไม่ลบจริง เพื่อ integrity)
    // - user: email → hash, name → "Deleted User", phone → null
    // - customer: name → "Deleted", contacts → ลบ
    // - orders: เก็บไว้ (legal requirement) แต่ anonymize PII
    // 2. ลบข้อมูลที่ไม่จำเป็น
    // - device tokens
    // - sessions, refresh tokens
    // - consent ที่ไม่ active
    // 3. Audit
    return s.auditRepo.Save(ctx, entity.NewAuditTrail(subjectID, "DSAR_ERASE", map[string]interface{}{
        "reason": reason, "at": time.Now().UTC(),
    }))
}
```

### G.4 Data Retention Policy

| Data Type | Retention | Reason |
|:---|:---|:---|
| **Telemetry** | 90 days (Pro), 365 days (Ent) | operational |
| **Audit logs** | 7 years | legal |
| **Invoices** | 10 years | tax law |
| **Sessions** | 7 days | security |
| **Refresh tokens** | 7 days (inactive) | security |
| **Consent records** | forever | compliance |
| **Deleted user PII** | immediately | PDPA |
| **Backups** | 30 days | DR |

### G.5 Data Processing Register (ROPA)

```yaml
# ตัวอย่าง Record of Processing Activities
activities:
  - name: "Customer Onboarding"
    purpose: "Provide IoT platform service"
    legal_basis: "Contract"
    data_categories:
      - "Name, email, phone"
      - "Tax ID"
      - "Address"
    retention: "Contract duration + 10 years"
    recipients: ["Payment processor", "Cloud provider"]
    transfers: ["AWS Singapore"]
    security: ["Encryption", "Access control", "Audit log"]

  - name: "Telemetry Analytics"
    purpose: "AI-powered insights"
    legal_basis: "Legitimate interest"
    data_categories:
      - "Device telemetry"
      - "Site information"
    retention: "90 days"
    recipients: ["LLM provider (Ollama - self-hosted)"]
    security: ["Encryption", "Anonymization where possible"]
```

### G.6 Breach Notification Procedure

```
1. Detect (alert, report, anomaly)
   ↓ (ภายใน 1 ชม.)
2. Triage (severity, scope, affected subjects)
   ↓ (ภายใน 24 ชม.)
3. Contain (stop bleeding, preserve evidence)
   ↓ (ภายใน 48 ชม.)
4. Assess (PDPA: แจ้งภายใน 72 ชม.)
   ↓
5. Notify:
   - PDPC (ถ้าจำเป็น)
   - Affected subjects (ถ้ามีความเสี่ยงสูง)
   - DPO, legal, PR
   ↓
6. Post-mortem + preventive actions
```

---

## 🅷 PART 14H — SECURITY OPERATIONS

### H.1 SIEM Setup

```
┌────────────────────────────────────────────────────────────┐
│                    SIEM (Elastic Security)                 │
│                                                            │
│  Sources:                                                  │
│  - App logs (Loki) ────────┐                               │
│  - Audit logs (Postgres)   │                               │
│  - K8s audit logs          ├──► Correlation ──► Alerts     │
│  - Network logs (VPC flow) │                               │
│  - MQTT broker logs        │                               │
│  - Cloudflare logs         ┘                               │
│                                                            │
│  Detections:                                               │
│  - Brute force (failed logins > 10/min)                    │
│  - Anomalous API calls (unusual patterns)                  │
│  - Data exfiltration (large exports)                       │
│  - Privilege escalation attempts                           │
│  - Device anomaly (unusual telemetry)                      │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

### H.2 Audit Log Schema

```sql
CREATE TABLE audit_trails (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    user_id UUID,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50),
    resource_id UUID,
    changes JSONB,          -- {before, after}
    ip VARCHAR(45),
    user_agent TEXT,
    request_id UUID,
    status VARCHAR(20),     -- SUCCESS, FAILED
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_audit_tenant_time ON audit_trails(tenant_id, created_at DESC);
CREATE INDEX idx_audit_user ON audit_trails(user_id, created_at DESC);
CREATE INDEX idx_audit_action ON audit_trails(action, created_at DESC);
CREATE INDEX idx_audit_resource ON audit_trails(resource_type, resource_id);
```

### H.3 Security Alerts

```yaml
# Prometheus alert rules
- alert: BruteForceLogin
  expr: rate(auth_login_failures_total[5m]) > 10
  for: 2m
  labels: { severity: critical }
  annotations:
    summary: "Possible brute force attack"

- alert: SuspiciousExport
  expr: increase(api_export_requests_total[1h]) > 100
  for: 5m
  labels: { severity: warning }
  annotations:
    summary: "Unusual data export volume"

- alert: PrivilegeEscalation
  expr: increase(auth_role_changes_total[1h]) > 5
  labels: { severity: critical }

- alert: AnomalousGeoAccess
  expr: count(count by (user_id, country) (auth_logins_total)) > 1
  for: 5m
  labels: { severity: warning }
```

### H.4 Penetration Testing

```
Scope:
- Web API (REST + WebSocket)
- Authentication flows
- Multi-tenant isolation
- File upload / download
- MQTT broker + device auth
- Admin panel

Frequency:
- Automated scan: weekly (OWASP ZAP, nuclei)
- Manual pen test: annually (external firm)
- Bug bounty: เริ่มหลัง launch 1 ปี

Tools:
- OWASP ZAP
- Burp Suite
- nuclei
- sqlmap
- nmap
```

### H.5 Vulnerability Management

```yaml
# Dependabot config
version: 2
updates:
  - package-ecosystem: gomod
    directory: /
    schedule: { interval: weekly }
    open-pull-requests-limit: 10
    reviewers: [security-team]
    labels: [dependencies, security]

  - package-ecosystem: docker
    directory: /
    schedule: { interval: weekly }

  - package-ecosystem: github-actions
    directory: /
    schedule: { interval: weekly }
```

```yaml
# Snyk scan (CI)
- name: Snyk scan
  uses: snyk/actions/golang@master
  env:
    SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}
  with:
    args: --severity-threshold=high
```

### H.6 Incident Response Plan

```
┌────────────────────────────────────────────────────────────┐
│                INCIDENT RESPONSE LIFECYCLE                 │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  1. PREPARATION                                            │
│     - Runbooks, on-call rotation, tools                    │
│     - Tabletop exercises                                   │
│                                                            │
│  2. DETECTION & ANALYSIS                                  │
│     - Alert → triage → severity (P0-P4)                    │
│     - P0: data breach, downtime > 15 min                   │
│     - P1: partial outage, security incident                │
│                                                            │
│  3. CONTAINMENT                                            │
│     - Short-term: isolate, block, disable                  │
│     - Long-term: patch, config change                      │
│                                                            │
│  4. ERADICATION                                            │
│     - Remove threat, patch vuln                            │
│                                                            │
│  5. RECOVERY                                               │
│     - Restore service, monitor                             │
│                                                            │
│  6. POST-INCIDENT                                          │
│     - Post-mortem (blameless) ภายใน 5 วัน                  │
│     - Action items + owner + due date                      │
│     - อัปเดต runbook                                       │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

### H.7 Compliance Audit Checklist

```
[ ] Access control matrix up to date
[ ] MFA enforced สำหรับ admin
[ ] Password policy (≥12 chars, complexity)
[ ] Session timeout (access 15 min, refresh 7 days)
[ ] Encryption at rest + in transit
[ ] Backup + tested restore
[ ] Audit logs ครบ 7 ปี
[ ] PDPA consent records
[ ] DSAR process documented
[ ] Breach notification procedure
[ ] Vendor risk assessment
[ ] Pen test report < 12 เดือน
[ ] Vulnerability remediation SLAs
[ ] Security awareness training
[ ] Incident response tested
[ ] Business continuity plan
```

---

## 🎯 PART 14 — SUMMARY

| Category | Count | Detail |
|---|:-:|---|
| **Auth Methods** | 3 | JWT + API Key + MQTT token |
| **Auth Layers** | 5 | Password + MFA + Session + Refresh + Device |
| **Authorization** | 3 | RBAC + ABAC + RLS |
| **Encryption** | 6 | TLS 1.3, AES-256-GCM, bcrypt, RS256, HMAC, KMS |
| **Compliance** | 2 | PDPA + GDPR |
| **Security Headers** | 10 | HSTS, CSP, X-Frame, ... |
| **OWASP Coverage** | 10/10 | ทุกข้อมี mitigation |
| **Runbooks** | 10+ | incidents + maintenance |
| **Audit Tables** | 3 | trails, consents, DSAR |
| **Pen Test** | 1/yr | + automated weekly |

---
---

# 🏁 สรุป PART 8, 12, 13, 14

```
┌─────────────────────────────────────────────────────────────────┐
│                   4 PARTS — FINAL OVERVIEW                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ★ PART 8 — INTERFACE LAYER                                     │
│  ─────────────────────────                                      │
│  • 120+ HTTP handlers · 10 middleware                           │
│  • WebSocket Hub + Client + Kafka bridge                        │
│  • MQTT Gateway (subscriber + publisher + auth)                 │
│  • 9 entry points (api, workers, scheduler, cli, ...)           │
│                                                                 │
│  ★ PART 12 — TESTING STRATEGY                                   │
│  ──────────────────────────                                     │
│  • 700+ unit · 250+ integration · 50+ E2E                       │
│  • testcontainers · mockery · k6 · pact                         │
│  • Coverage targets: Domain 95%, App 85%                        │
│  • CI pipeline (GitHub Actions)                                 │
│                                                                 │
│  ★ PART 13 — DEPLOYMENT & DEVOPS                                │
│  ─────────────────────────────                                  │
│  • Multi-stage Dockerfile · docker-compose · K8s               │
│  • CI/CD (GitLab, ArgoCD) · HPA · KEDA                         │
│  • Prometheus + Grafana + Loki + Jaeger                         │
│  • Backup + DR + Runbooks                                       │
│                                                                 │
│  ★ PART 14 — SECURITY                                           │
│  ────────────────────                                           │
│  • Auth (JWT + MFA + API Key) · RBAC + ABAC + RLS              │
│  • Encryption (TLS 1.3, AES-256, bcrypt, RS256)                 │
│  • OWASP Top 10 · Device security · mTLS                        │
│  • PDPA + GDPR · SIEM · Incident response                       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 📁 Complete Documentation Map

```
📚 icmongolang Documentation Suite
│
├── Architecture_design.md              (SRS + Roadmap)
├── Architecture_modules_9.md           (Sample Data & Fixtures)
│
├── PART_5_Domain_Layer.md              ✅ (in response)
├── PART_6_Application_Layer.md         ✅ (in response)
├── PART_7_Infrastructure_Layer.md      ✅ (in response)
├── PART_8_Interface_Layer.md           ✅ (in response นี้)
├── PART_9_Sample_Data.md               ✅ (PART 9A-F)
├── PART_10_Postman.md                  ✅ (PART 10A-G)
├── PART_11_Executive_Summary.md        ✅ (PART 11A-F)
├── PART_12_Testing.md                  ✅ (in response นี้)
├── PART_13_DevOps.md                   ✅ (in response นี้)
└── PART_14_Security.md                 ✅ (in response นี้)
```

### 🎯 Documentation Complete

```
╔═══════════════════════════════════════════════════════════════════╗
║                                                                   ║
║              ✅  icmongolang — ALL PARTS COMPLETE                ║
║                                                                   ║
║  📘 PART 1-4:   Architecture Design (SRS)                        ║
║  📘 PART 5:     Domain Layer Deep Dive                           ║
║  📘 PART 6:     Application Layer Deep Dive                      ║
║  📘 PART 7:     Infrastructure Layer Deep Dive                   ║
║  📘 PART 8:     Interface Layer Deep Dive                        ║
║  📘 PART 9:     Sample Data & Fixtures (A-F)                     ║
║  📘 PART 10:    Postman Collection (A-G)                         ║
║  📘 PART 11:    Executive Summary (A-F)                          ║
║  📘 PART 12:    Testing Strategy                                 ║
║  📘 PART 13:    Deployment & DevOps                              ║
║  📘 PART 14:    Security Deep Dive                               ║
║                                                                   ║
║  📦 40 Modules · 10 Bounded Contexts · 159 Use Cases             ║
║  🏗️  Clean Architecture + DDD + Event-Driven                     ║
║  🔒 Production-Ready: Security + Testing + DevOps                ║
║                                                                   ║
║  🎯 Ready for: Dev · Investment · Team Onboarding · Audit         ║
║                                                                   ║
╚═══════════════════════════════════════════════════════════════════╝
``` 