# Module 12: pkg/jwt (JWT Authentication)

## สำหรับโฟลเดอร์ `internal/pkg/jwt/`

ไฟล์ที่เกี่ยวข้อง:
- `internal/pkg/jwt/maker.go`
- `internal/pkg/jwt/payload.go`
- `internal/pkg/jwt/rsa_maker.go`

---

## หลักการ (Concept)

### คืออะไร?
JWT (JSON Web Token) เป็นมาตรฐาน (RFC 7519) สำหรับสร้าง token ที่มีข้อมูล (claims) ฝังอยู่และสามารถตรวจสอบความถูกต้องได้โดยใช้ลายเซ็นดิจิทัล โดยไม่ต้องเก็บ session บนเซิร์ฟเวอร์ (stateless authentication)

### มีกี่แบบ?
1. **HMAC (HS256)** – ใช้ secret key ร่วมกัน (symmetric) เร็ว ง่าย แต่แจกจ่าย secret ลำบาก
2. **RSA (RS256)** – ใช้ private key สำหรับ sign, public key สำหรับ verify (asymmetric) ปลอดภัยกว่า เหมาะกับ microservices
3. **ECDSA (ES256)** – เช่น RSA แต่ใช้คีย์สั้นกว่า เร็ว เหมาะกับอุปกรณ์ IoT

### ใช้อย่างไร / นำไปใช้กรณีไหน
- ใช้ RS256 ใน production เพราะแยก public/private key
- private key เก็บเฉพาะ authorization server
- public key แจกจ่ายให้ service อื่นเพื่อ verify token
- access token อายุสั้น (15 นาที), refresh token อายุยาว (7 วัน)

### ทำไมต้องใช้
- Stateless → ไม่ต้องเก็บ session ใน database
- รองรับการขยายแนวนอน (scale out) ได้ดี
- สามารถเก็บข้อมูล user id, role ใน token

### ประโยชน์ที่ได้รับ
- ลด load database
- รองรับ microservices architecture
- token ตรวจสอบได้โดยไม่ต้องเรียก central service

### ข้อควรระวัง
- payload ของ JWT เป็นแค่ base64 encoded (ไม่เข้ารหัส) ห้ามเก็บ secret
- token ขนาดใหญ่กว่า session ID (ส่งผลต่อ bandwidth)
- revoke token ทำได้ยาก (ต้องใช้ blacklist)

### ข้อดี
- Stateless, scalable, รองรับหลาย platform

### ข้อเสีย
- Revoke ยาก, token มีขนาดใหญ่, ไม่สามารถ invalidate ได้ทันที

### ข้อห้าม
- ห้ามเก็บ password หรือข้อมูลอ่อนไหวใน payload
- ห้ามใช้ HS256 ถ้าต้องแจกจ่าย key ให้หลาย service
- ห้ามตั้ง expiry นานเกินไป (access token ไม่ควรเกิน 1 ชั่วโมง)

---

## โค้ดที่รันได้จริง

### 1. Maker Interface – `maker.go`

```go
// Package jwt provides JWT creation and verification using RSA256.
// ----------------------------------------------------------------
// แพ็คเกจ jwt ให้บริการสร้างและตรวจสอบ JWT ด้วย RSA256
package jwt

import (
	"errors"
	"time"
)

// Common JWT errors.
// ----------------------------------------------------------------
// ข้อผิดพลาด JWT ที่พบบ่อย
var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

// Maker is an interface for managing JWT tokens.
// ----------------------------------------------------------------
// Maker คือ interface สำหรับจัดการ JWT tokens
type Maker interface {
	// CreateToken creates a new JWT token for given user and role.
	// ----------------------------------------------------------------
	// CreateToken สร้าง JWT token ใหม่สำหรับผู้ใช้และบทบาทที่กำหนด
	CreateToken(userID uint, role string, duration time.Duration) (string, *Payload, error)

	// VerifyToken validates the token and returns its payload.
	// ----------------------------------------------------------------
	// VerifyToken ตรวจสอบ token และคืนค่า payload
	VerifyToken(token string) (*Payload, error)
}
```

### 2. Payload – `payload.go`

```go
package jwt

import (
	"time"

	"github.com/google/uuid"
)

// Payload contains the JWT claims.
// ----------------------------------------------------------------
// Payload บรรจุ claims ของ JWT
type Payload struct {
	ID        uuid.UUID `json:"id"`         // unique token ID (jti)
	UserID    uint      `json:"user_id"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload.
// ----------------------------------------------------------------
// NewPayload สร้าง payload ใหม่สำหรับ token
func NewPayload(userID uint, role string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Payload{
		ID:        tokenID,
		UserID:    userID,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

// Valid checks if the payload is not expired (required for jwt.Claims interface).
// ----------------------------------------------------------------
// Valid ตรวจสอบว่า payload ยังไม่หมดอายุ (required สำหรับ jwt.Claims interface)
func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
```

### 3. RSA Maker – `rsa_maker.go`

```go
package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// RSAMaker implements Maker using RSA256 algorithm.
// ----------------------------------------------------------------
// RSAMaker อิมพลีเมนต์ Maker ด้วยอัลกอริทึม RSA256
type RSAMaker struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewRSAMaker creates a new RSAMaker from PEM-encoded keys (base64 or raw).
// ----------------------------------------------------------------
// NewRSAMaker สร้าง RSAMaker ใหม่จาก PEM keys (base64 หรือ raw)
func NewRSAMaker(privateKeyPEM, publicKeyPEM string) (*RSAMaker, error) {
	// Parse private key
	// แปลง private key
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to decode private key PEM")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS1
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
	}
	rsaPrivate, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not RSA")
	}

	// Parse public key
	// แปลง public key
	blockPub, _ := pem.Decode([]byte(publicKeyPEM))
	if blockPub == nil {
		return nil, fmt.Errorf("failed to decode public key PEM")
	}
	publicKey, err := x509.ParsePKIXPublicKey(blockPub.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	rsaPublic, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not RSA")
	}

	return &RSAMaker{
		privateKey: rsaPrivate,
		publicKey:  rsaPublic,
	}, nil
}

// CreateToken generates a new RSA256 JWT token.
// ----------------------------------------------------------------
// CreateToken สร้าง JWT token ใหม่ด้วย RSA256
func (m *RSAMaker) CreateToken(userID uint, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, role, duration)
	if err != nil {
		return "", nil, err
	}

	claims := jwt.MapClaims{
		"jti":      payload.ID.String(),
		"user_id":  payload.UserID,
		"role":     payload.Role,
		"exp":      payload.ExpiredAt.Unix(),
		"iat":      payload.IssuedAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(m.privateKey)
	if err != nil {
		return "", nil, err
	}
	return signedToken, payload, nil
}

// VerifyToken validates the token and returns its payload.
// ----------------------------------------------------------------
// VerifyToken ตรวจสอบ token และคืนค่า payload
func (m *RSAMaker) VerifyToken(tokenString string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrInvalidToken
		}
		return m.publicKey, nil
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// Extract claims
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}
	role, _ := claims["role"].(string)
	jtiStr, _ := claims["jti"].(string)
	exp, _ := claims["exp"].(float64)
	iat, _ := claims["iat"].(float64)

	tokenID, err := uuid.Parse(jtiStr)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &Payload{
		ID:        tokenID,
		UserID:    uint(userIDFloat),
		Role:      role,
		IssuedAt:  time.Unix(int64(iat), 0),
		ExpiredAt: time.Unix(int64(exp), 0),
	}, nil
}
```

### 4. ตัวอย่างการสร้าง RSA Keys และใช้งาน

**สร้าง private/public key (ใช้ openssl):**
```bash
# Generate private key
openssl genrsa -out private.pem 2048

# Extract public key
openssl rsa -in private.pem -pubout -out public.pem

# Convert to base64 for environment variables
cat private.pem | base64 | tr -d '\n'
cat public.pem | base64 | tr -d '\n'
```

**ตัวอย่างการใช้งานใน `auth_usecase.go`:**
```go
func (u *authUsecase) Login(ctx context.Context, email, password string) (string, string, error) {
    // ... validate user ...
    
    // Create access token (15 minutes)
    accessToken, _, err := u.jwtMaker.CreateToken(user.ID, string(user.Role), 15*time.Minute)
    if err != nil {
        return "", "", err
    }
    
    // Create refresh token (7 days)
    refreshToken := uuid.New().String()
    // store refresh token in Redis...
    
    return accessToken, refreshToken, nil
}
```

**ตัวอย่างการตรวจสอบ token ใน middleware:**
```go
func JWTAuth(maker jwt.Maker) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            tokenString := extractBearerToken(r)
            payload, err := maker.VerifyToken(tokenString)
            if err != nil {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            ctx := context.WithValue(r.Context(), "user_id", payload.UserID)
            ctx = context.WithValue(ctx, "role", payload.Role)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

---

## วิธีใช้งาน module นี้

1. ติดตั้ง dependencies:
   ```bash
   go get github.com/golang-jwt/jwt/v5
   go get github.com/google/uuid
   ```
2. วางไฟล์ `maker.go`, `payload.go`, `rsa_maker.go` ใน `internal/pkg/jwt/`
3. สร้าง JWT maker ใน `main.go`:
   ```go
   privateKeyBase64 := os.Getenv("JWT_PRIVATE_KEY_BASE64")
   publicKeyBase64 := os.Getenv("JWT_PUBLIC_KEY_BASE64")
   privateKeyBytes, _ := base64.StdEncoding.DecodeString(privateKeyBase64)
   publicKeyBytes, _ := base64.StdEncoding.DecodeString(publicKeyBase64)
   jwtMaker, err := jwt.NewRSAMaker(string(privateKeyBytes), string(publicKeyBytes))
   ```
4. Inject เข้า `authUsecase` และ `middleware`

---

## ตารางสรุปฟังก์ชันหลัก

| ฟังก์ชัน | อินพุต | เอาต์พุต | ใช้เมื่อ |
|----------|--------|----------|---------|
| `CreateToken` | userID, role, duration | token string, payload, error | login, refresh |
| `VerifyToken` | token string | payload, error | middleware auth |
| `NewPayload` | userID, role, duration | payload, error | ภายใน CreateToken |

---

## แบบฝึกหัดท้าย module (3 ข้อ)

1. เพิ่ม method `CreateTokenWithCustomClaims` ที่รองรับ claims เพิ่มเติม (เช่น `device_id`)
2. Implement `HS256Maker` สำหรับ environment ที่ไม่ต้องการ RSA (development) โดยใช้ secret key
3. เพิ่ม `RefreshToken` function ที่สร้าง refresh token แบบ UUID และเก็บ payload ใน Redis (เชื่อมกับ pkg/redis)

---

## แหล่งอ้างอิง

- [JWT RFC 7519](https://tools.ietf.org/html/rfc7519)
- [golang-jwt/jwt documentation](https://github.com/golang-jwt/jwt)
- [RS256 vs HS256](https://auth0.com/blog/rs256-vs-hs256-whats-the-difference/)

---

**หมายเหตุ:** module นี้ครบถ้วนสำหรับ `pkg/jwt` หากต้องการ module เพิ่มเติม (เช่น `pkg/hash`) โปรดแจ้ง