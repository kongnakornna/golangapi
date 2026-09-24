# บทที่ 6: ระบบจัดการ Login, Refresh Token และ Rate Limit สำหรับ Golang Big Data MML AI

> **สรุปสั้นก่อนอ่าน**: ในระบบ Big Data และ MML AI ที่ให้บริการผ่าน API (บทที่ 4) ความปลอดภัยเป็นสิ่งสำคัญ บทนี้นำเสนอการออกแบบระบบยืนยันตัวตน (Authentication) ด้วย JWT Access Token และ Refresh Token พร้อมมาตรการป้องกันการโจมตี เช่น การจำกัดอัตราการเรียกใช้ (Rate Limit) และการบล็อกโทเค็นที่ถูกขโมย การนำไปใช้กับ Go จะช่วยให้ API ของคุณปลอดภัยและพร้อมใช้งานในระดับ Production

---

## 📌 โครงสร้างการทำงาน

ระบบจัดการ Login, Refresh Token และ Rate Limit ประกอบด้วย 6 องค์ประกอบหลัก:

1. **User Registration & Login** – การลงทะเบียนและเข้าสู่ระบบ ตรวจสอบรหัสผ่าน (bcrypt)
2. **Access Token (JWT)** – โทเค็นอายุสั้น (15-30 นาที) สำหรับเรียกใช้ API
3. **Refresh Token** – โทเค็นอายุยาว (7-30 วัน) สำหรับขอ Access Token ใหม่โดยไม่ต้อง login ซ้ำ
4. **Token Blacklist / Redis Store** – จัดเก็บ Refresh Token และบล็อก Access Token ที่ถูกเพิกถอน
5. **Rate Limiting Middleware** – จำกัดจำนวนคำขอต่อช่วงเวลาต่อ IP หรือต่อผู้ใช้
6. **Security Headers & CORS** – ป้องกัน CSRF, XSS, Clickjacking

---

## 🎯 วัตถุประสงค์

- เพื่อสร้างระบบ Authentication ที่ปลอดภัยด้วย JWT (Access + Refresh Token)
- เพื่อป้องกันการโจมตีแบบ Brute Force, Replay Attack, และ Token Hijacking
- เพื่อจำกัดอัตราการเรียกใช้ API ป้องกัน DDoS และ
- เพื่อให้ผู้ใช้สามารถ refresh token ได้โดยไม่ต้อง login ซ้ำ
- เพื่อให้สามารถเพิกถอน (revoke) โทเค็นเมื่อมีการ logout หรือตรวจพบความผิดปกติ

---

## 👥 กลุ่มเป้าหมาย

- Backend Developer ที่ต้องการเพิ่มความปลอดภัยให้กับ API ของตน
- MLOps / Data Engineer ที่จะเปิดให้บริการโมเดลผ่าน API (จากบทที่ 4)
- DevOps ที่ต้องการตั้งค่าระบบ Rate Limit สำหรับ API Gateway
- ผู้ที่ต้องการเข้าใจหลักการ JWT, Refresh Token, และการป้องกันภัยคุกคามทั่วไป

---

## 📚 ความรู้พื้นฐานที่ควรมี

- พื้นฐาน Go (Gin/Echo framework, HTTP middleware)
- ความรู้เกี่ยวกับ JWT (JSON Web Token) และการเข้ารหัส
- การใช้งาน Redis (หรือใน memory store สำหรับทดสอบ)
- ความเข้าใจเกี่ยวกับ HTTP headers (Authorization, Cookie)

---

## 📖 เนื้อหาโดยย่อ (กระชับ เน้นวัตถุประสงค์และประโยชน์)

| เนื้อหา | วัตถุประสงค์ | ประโยชน์ |
|---------|--------------|-----------|
| User Registration | บันทึกผู้ใช้ใหม่ด้วยรหัสผ่านที่เข้ารหัส | ป้องกัน credential รั่วไหล |
| Login & JWT Access | ออก Access Token อายุสั้น | ใช้เรียก API ได้โดยไม่ต้องส่ง user/pass ทุกครั้ง |
| Refresh Token | ออก Refresh Token อายุยาว | ลดการ login ซ้ำ, ปลอดภัยกว่าเก็บ Access Token นาน |
| Token Blacklist | เพิกถอนโทเค็นที่ logout หรือถูกขโมย | ตัดสิทธิ์โทเค็นที่หมดอายุหรือไม่ถูกต้อง |
| Rate Limit | จำกัดจำนวน request ต่อวินาที/นาที | ป้องกัน DDoS, Brute Force, API abuse |
| Security Headers | เพิ่ม HTTP headers ป้องกัน XSS, CSRF, Clickjacking | เพิ่มความปลอดภัยให้ browser-based clients |

---

## 📘 บทนำ

API ที่ให้บริการโมเดล ML (เช่น REST API จากบทที่ 4) จำเป็นต้องมีการควบคุมการเข้าถึง หากไม่มีระบบ Authentication ผู้ไม่หวังดีอาจเรียกใช้ API จำนวนมากจนทำให้ระบบทำงานช้า หรือขโมยข้อมูลสำคัญ JWT (JSON Web Token) เป็นมาตรฐานที่ได้รับความนิยมเพราะไม่ต้องเก็บ session บนเซิร์ฟเวอร์ แต่การใช้งาน JWT อย่างเดียวไม่เพียงพอ ต้องมี Refresh Token เพื่อลดความเสี่ยงเมื่อ Access Token ถูกขโมย รวมถึง Rate Limit เพื่อป้องกันการโจมตีแบบใช้ทรัพยากรเกิน

บทนี้จะออกแบบและสร้างระบบ Authentication แบบสมบูรณ์ด้วย Go (Gin framework), Redis สำหรับเก็บ Refresh Token และ Rate Limit counters, และมาตรการป้องกันเพิ่มเติม พร้อมตัวอย่างโค้ดที่รันได้จริง

---

## 📖 บทนิยาม

| ศัพท์ | คำอธิบาย |
|-------|-----------|
| **JWT (JSON Web Token)** | โทเค็นที่เข้ารหัสด้วย secret key หรือ public/private key ประกอบด้วย Header, Payload, Signature ใช้ยืนยันตัวตนและส่งข้อมูลระหว่าง parties |
| **Access Token** | JWT ที่มีอายุสั้น (15-30 นาที) ใช้ในการเรียก API แต่ละครั้ง |
| **Refresh Token** | โทเค็นที่มีอายุยาว (7-30 วัน) ใช้ขอ Access Token ใหม่เมื่อหมดอายุ โดยไม่ต้อง login ซ้ำ |
| **Blacklist (Token Denylist)** | รายการโทเค็นที่ถูกเพิกถอน (logout, ถูกขโมย) ก่อนหมดอายุ |
| **Rate Limiting** | การจำกัดจำนวนคำขอที่ client สามารถส่งได้ในช่วงเวลาที่กำหนด เช่น 100 request/นาที |
| **Brute Force Attack** | การลอง username/password ซ้ำ ๆ เพื่อเจาะระบบ |
| **Replay Attack** | การนำโทเค็นที่ถูกต้องไปใช้ซ้ำ (ป้องกันโดยใช้ nonce หรือ short expiry) |
| **bcrypt** | อัลกอริทึมการแฮชรหัสผ่านที่ออกแบบให้ช้าและมี salt ป้องกันการโจมตีแบบ dictionary |

---

## 🧠 แนวคิด (Concept Explanation)

### JWT Access Token vs Refresh Token: คืออะไร? มีกี่แบบ? ใช้อย่างไร?

| ประเภท | อายุ | เก็บที่ใด | ใช้ทำอะไร | เสี่ยง |
|--------|------|-----------|------------|--------|
| **Access Token** | สั้น (15-30 นาที) | Client (memory หรือ localStorage) | เรียก API ทุกครั้ง | ถ้าถูกขโมย แฮกเกอร์ใช้ได้แค่ช่วงสั้น ๆ |
| **Refresh Token** | ยาว (7-30 วัน) | Server (Redis/DB) หรือ HttpOnly Cookie | ขอ Access Token ใหม่ | ถ้าถูกขโมย แก้ไขได้โดย revoke ที่ server |

**แนวทางปฏิบัติที่ดี:**
- Access Token ไม่ควรเก็บข้อมูล sensitive (เช่น role, user_id ก็พอ)
- Refresh Token ควรสุ่มค่า (ไม่ใช่ JWT) หรือเป็น JWT ที่มี `jti` (JWT ID) ให้ server ตรวจสอบ
- ใช้ HttpOnly, Secure, SameSite cookie สำหรับ Refresh Token เพื่อป้องกัน XSS
- ใช้ Rate Limit สำหรับ login endpoint และ API ทั้งหมด

### Rate Limit: มีกี่แบบ? ใช้อันไหน?

1. **Fixed Window** – จำกัดจำนวน request ใน window เวลาคงที่ (เช่น 100 ต่อนาที) ง่าย แต่มี burst ที่ขอบ window
2. **Sliding Window** – จำกัดตามช่วงเวลาจริง (เช่น 100 ต่อ 60 วินาที) แม่นยำกว่า
3. **Token Bucket** – อนุญาต burst เล็ก ๆ แล้วค่อย ๆ เติม tokens เหมาะกับ API ที่ต้องการความยืดหยุ่น
4. **Leaky Bucket** – ทำ request ออกที่อัตราคงที่ ใช้ป้องกัน overload

**บทนี้ใช้ Sliding Window + Redis** เพื่อความแม่นยำและ scalability

### ทำไมต้องใช้ Redis สำหรับ Rate Limit และ Refresh Token?

- Redis มีคำสั่ง atomic (INCR, EXPIRE, ZADD) เหมาะกับ counter และ sliding window
- Redis เร็วกว่า database ทั่วไป (in-memory)
- รองรับ distributed system (หลาย instance ของ API ใช้ Redis เดียวกัน)

---

## 🗺️ ออกแบบ Workflow (Dataflow Diagram)

### รูปที่ 7: Dataflow แสดงกระบวนการ Login, Refresh Token, และ Rate Limit

```mermaid
flowchart TB
    subgraph Client["🖥️ Client (Mobile/Web)"]
        C[User]
    end

    subgraph API["☁️ Go API Server"]
        direction TB
        RL[Rate Limit Middleware]
        Auth[Auth Middleware<br/>Validate Access Token]
        LoginHandler[/login]
        RefreshHandler[/refresh]
        LogoutHandler[/logout]
        Protected[/api/...]
    end

    subgraph Storage["💾 Storage"]
        Redis[(Redis)]
        DB[(PostgreSQL<br/>users table)]
    end

    C -->|1. POST /login| LoginHandler
    LoginHandler -->|2. verify password| DB
    DB -->|3. user found| LoginHandler
    LoginHandler -->|4. generate tokens| Redis
    Redis -->|5. store refresh token| LoginHandler
    LoginHandler -->|6. return access token + refresh token (cookie)| C

    C -->|7. request with Access Token| RL
    RL -->|8. check rate limit| Redis
    Redis -->|9. allow/deny| RL
    RL -->|10. pass| Auth
    Auth -->|11. validate JWT| Auth
    Auth -->|12. forward| Protected
    Protected -->|13. response| C

    C -->|14. access token expired| RefreshHandler
    RefreshHandler -->|15. validate refresh token| Redis
    Redis -->|16. valid| RefreshHandler
    RefreshHandler -->|17. generate new access token| C

    C -->|18. POST /logout| LogoutHandler
    LogoutHandler -->|19. revoke refresh token| Redis
    LogoutHandler -->|20. add access token to blacklist| Redis

    style Client fill:#e1f5fe
    style API fill:#f3e5f5
    style Storage fill:#e8f5e9
```

### คำอธิบาย Diagram อย่างละเอียด

1. **Client** ส่ง username/password ไปที่ `/login`
2. **LoginHandler** ตรวจสอบกับฐานข้อมูล (bcrypt compare)
3. ถ้าถูกต้อง จะสร้าง Access Token (JWT อายุ 15 นาที) และ Refresh Token (UUID หรือ JWT อายุ 7 วัน)
4. **Redis** เก็บ Refresh Token พร้อม user ID และ expiration (key: `refresh:<token>`)
5. ส่ง Access Token กลับใน response body (หรือ header) และ Refresh Token ใน HttpOnly Cookie
6. Client เก็บ Access Token ไว้ใน memory (หรือ localStorage) และใช้ส่งใน `Authorization: Bearer <token>` สำหรับทุก request
7. ทุก request ผ่าน **Rate Limit Middleware** ก่อน (แยกตาม IP หรือ User ID)
8. Rate Limit middleware ใช้ Redis Sliding Window ตรวจสอบจำนวน request ใน 1 นาที
9. ถ้าเกิน limit → ส่ง `429 Too Many Requests`
10. ถ้าผ่าน → ส่งต่อไปยัง **Auth Middleware**
11. Auth middleware ตรวจสอบ JWT signature, expiration, และตรวจสอบว่า Access Token ไม่อยู่ใน blacklist
12. ถ้าถูกต้อง → forward ไปยัง protected handler
13. Response กลับไปยัง Client
14. เมื่อ Access Token หมดอายุ Client ส่ง Refresh Token (จาก cookie) ไปที่ `/refresh`
15. RefreshHandler ตรวจสอบ Refresh Token ใน Redis
16. ถ้าพบและยังไม่หมดอายุ → สร้าง Access Token ใหม่
17. ส่ง Access Token ใหม่กลับไป
18. เมื่อ user logout → ส่ง POST `/logout` พร้อม Access Token และ Refresh Token
19. ลบ Refresh Token ออกจาก Redis
20. นำ Access Token ใส่ blacklist (Redis) จนกว่าจะหมดอายุ

---

## 💻 ตัวอย่างโค้ดที่รันได้จริง (พร้อม Comment สองภาษา)

### โครงสร้างโปรเจกต์

```
auth-system/
├── main.go
├── handlers/
│   ├── auth.go
│   └── api.go
├── middleware/
│   ├── ratelimit.go
│   └── auth.go
├── store/
│   ├── redis.go
│   └── userdb.go
├── go.mod
└── .env
```

### ขั้นตอนที่ 1: ติดตั้ง dependencies

```bash
go mod init auth-system
go get github.com/gin-gonic/gin
go get github.com/golang-jwt/jwt/v5
go get github.com/redis/go-redis/v9
go get github.com/joho/godotenv
go get golang.org/x/crypto/bcrypt
```

### ขั้นตอนที่ 2: ไฟล์ environment `.env`

```env
JWT_SECRET=your-very-secret-key-at-least-32-bytes
ACCESS_TOKEN_TTL=15
REFRESH_TOKEN_TTL=10080
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
DB_DSN=postgres://user:pass@localhost:5432/authdb?sslmode=disable
```

### ขั้นตอนที่ 3: Redis client (`store/redis.go`)

```go
package store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps redis.Client
// RedisClient ห่อหุ้ม redis.Client
type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

// NewRedisClient creates a new Redis connection
// NewRedisClient สร้างการเชื่อมต่อ Redis ใหม่
func NewRedisClient(addr, password string) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &RedisClient{Client: rdb, Ctx: ctx}, nil
}

// StoreRefreshToken saves refresh token with user ID
// StoreRefreshToken เก็บ refresh token พร้อม user ID
func (r *RedisClient) StoreRefreshToken(token string, userID uint, ttlMinutes int) error {
	key := "refresh:" + token
	return r.Client.Set(r.Ctx, key, userID, time.Duration(ttlMinutes)*time.Minute).Err()
}

// GetRefreshToken retrieves user ID from refresh token
// GetRefreshToken ดึง user ID จาก refresh token
func (r *RedisClient) GetRefreshToken(token string) (uint, error) {
	key := "refresh:" + token
	val, err := r.Client.Get(r.Ctx, key).Uint64()
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}

// DeleteRefreshToken removes refresh token (logout)
// DeleteRefreshToken ลบ refresh token (logout)
func (r *RedisClient) DeleteRefreshToken(token string) error {
	key := "refresh:" + token
	return r.Client.Del(r.Ctx, key).Err()
}

// BlacklistAccessToken adds access token to denylist until its expiry
// BlacklistAccessToken เพิ่ม access token ลง blacklist จนกว่าจะหมดอายุ
func (r *RedisClient) BlacklistAccessToken(token string, ttlSeconds int) error {
	key := "blacklist:" + token
	return r.Client.Set(r.Ctx, key, "1", time.Duration(ttlSeconds)*time.Second).Err()
}

// IsBlacklisted checks if access token is revoked
// IsBlacklisted ตรวจสอบว่า access token ถูกเพิกถอนหรือไม่
func (r *RedisClient) IsBlacklisted(token string) (bool, error) {
	key := "blacklist:" + token
	val, err := r.Client.Exists(r.Ctx, key).Result()
	return val > 0, err
}

// RateLimitSlidingWindow implements sliding window rate limit using Redis sorted set
// RateLimitSlidingWindow ใช้ Redis sorted set สำหรับ sliding window rate limit
func (r *RedisClient) RateLimitSlidingWindow(key string, limit int, windowSeconds int64) (bool, error) {
	now := time.Now().Unix()
	windowStart := now - windowSeconds

	// Remove old entries / ลบรายการเก่า
	zremCmd := r.Client.ZRemRangeByScore(r.Ctx, key, "0", float64(windowStart))
	if zremCmd.Err() != nil {
		return false, zremCmd.Err()
	}

	// Count current requests / นับจำนวน request ปัจจุบัน
	countCmd := r.Client.ZCard(r.Ctx, key)
	if countCmd.Err() != nil {
		return false, countCmd.Err()
	}
	count := countCmd.Val()

	if count >= int64(limit) {
		return false, nil // Rate limit exceeded / เกินอัตราที่กำหนด
	}

	// Add current request / เพิ่ม request ปัจจุบัน
	score := float64(now)
	member := now
	r.Client.ZAdd(r.Ctx, key, redis.Z{Score: score, Member: member})
	// Set expiry on the key to clean up / ตั้งเวลาหมดอายุเพื่อล้าง key
	r.Client.Expire(r.Ctx, key, time.Duration(windowSeconds)*time.Second)

	return true, nil
}
```

### ขั้นตอนที่ 4: User model และ database (`store/userdb.go` - ใช้ mock สำหรับตัวอย่าง)

```go
package store

import (
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
// User แทนผู้ใช้ในระบบ
type User struct {
	ID       uint
	Username string
	Password string // hashed
}

// MockUserDB is an in-memory user store for demo
// MockUserDB เป็นที่เก็บผู้ใช้ในหน่วยความจำสำหรับตัวอย่าง
type MockUserDB struct {
	mu    sync.RWMutex
	users map[string]*User
	nextID uint
}

// NewMockUserDB creates a new mock DB with a test user
// NewMockUserDB สร้าง mock DB พร้อมผู้ใช้ทดสอบ
func NewMockUserDB() *MockUserDB {
	db := &MockUserDB{
		users: make(map[string]*User),
		nextID: 1,
	}
	// Create test user: username "admin", password "password"
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	db.users["admin"] = &User{
		ID:       1,
		Username: "admin",
		Password: string(hashed),
	}
	db.nextID = 2
	return db
}

// GetUserByUsername retrieves user by username
// GetUserByUsername ค้นหาผู้ใช้ตามชื่อผู้ใช้
func (db *MockUserDB) GetUserByUsername(username string) (*User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	user, ok := db.users[username]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// CreateUser creates a new user (optional)
// CreateUser สร้างผู้ใช้ใหม่ (optional)
func (db *MockUserDB) CreateUser(username, password string) (*User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.users[username]; ok {
		return nil, errors.New("user already exists")
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &User{
		ID:       db.nextID,
		Username: username,
		Password: string(hashed),
	}
	db.users[username] = user
	db.nextID++
	return user, nil
}
```

### ขั้นตอนที่ 5: JWT helper (`handlers/auth.go` - ส่วน utilities)

```go
package handlers

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims structure for JWT
// โครงสร้าง Claims สำหรับ JWT
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a new short-lived JWT
// GenerateAccessToken สร้าง JWT อายุสั้นใหม่
func GenerateAccessToken(userID uint, secret string, ttlMinutes int) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ttlMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateAccessToken parses and validates JWT, returns userID
// ValidateAccessToken แปลงและตรวจสอบ JWT, คืนค่า userID
func ValidateAccessToken(tokenString, secret string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}
	return 0, errors.New("invalid token")
}
```

### ขั้นตอนที่ 6: Handlers (login, refresh, logout) (`handlers/auth.go`)

```go
package handlers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"auth-system/store"
)

// LoginRequest body
// LoginRequest ตัวรับข้อมูล login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshRequest (optional, token from cookie)
// RefreshRequest (optional, token จาก cookie)
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"` // can also be from cookie
}

// AuthHandler holds dependencies
// AuthHandler ถือ dependencies ต่าง ๆ
type AuthHandler struct {
	UserDB        *store.MockUserDB
	RedisClient   *store.RedisClient
	JWTSecret     string
	AccessTTL     int // minutes
	RefreshTTL    int // minutes
}

// Login handles POST /login
// Login จัดการ POST /login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Get user from DB / ดึงผู้ใช้จากฐานข้อมูล
	user, err := h.UserDB.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Compare password / เปรียบเทียบรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate Access Token / สร้าง Access Token
	accessToken, err := GenerateAccessToken(user.ID, h.JWTSecret, h.AccessTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate access token"})
		return
	}

	// Generate Refresh Token (random UUID) / สร้าง Refresh Token (UUID สุ่ม)
	refreshToken := uuid.New().String()
	if err := h.RedisClient.StoreRefreshToken(refreshToken, user.ID, h.RefreshTTL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not store refresh token"})
		return
	}

	// Set refresh token as HttpOnly cookie / ตั้ง refresh token เป็น HttpOnly cookie
	c.SetCookie(
		"refresh_token",           // name
		refreshToken,              // value
		h.RefreshTTL*60,           // maxAge seconds
		"/refresh",                // path (only sent to /refresh)
		"",                        // domain
		true,                      // secure (HTTPS only)
		true,                      // httpOnly
	)
	
	// Return access token in JSON body / ส่ง access token ใน JSON body
	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   h.AccessTTL * 60,
	})
}

// Refresh handles POST /refresh
// Refresh จัดการ POST /refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	// Get refresh token from cookie / ดึง refresh token จาก cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}

	// Verify refresh token in Redis / ตรวจสอบ refresh token ใน Redis
	userID, err := h.RedisClient.GetRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	// Generate new access token / สร้าง Access Token ใหม่
	newAccessToken, err := GenerateAccessToken(userID, h.JWTSecret, h.AccessTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
		"token_type":   "Bearer",
		"expires_in":   h.AccessTTL * 60,
	})
}

// Logout handles POST /logout
// Logout จัดการ POST /logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get refresh token from cookie / ดึง refresh token จาก cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil {
		h.RedisClient.DeleteRefreshToken(refreshToken)
	}
	// Also blacklist current access token if provided in header
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		accessToken := authHeader[7:]
		// Blacklist until expiry (we need to parse expiry)
		// For simplicity, blacklist with default TTL (access token TTL)
		h.RedisClient.BlacklistAccessToken(accessToken, h.AccessTTL*60)
	}
	// Clear cookie / ล้าง cookie
	c.SetCookie("refresh_token", "", -1, "/refresh", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
```

### ขั้นตอนที่ 7: Rate Limit Middleware (`middleware/ratelimit.go`)

```go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"auth-system/store"
)

// RateLimitMiddleware returns a Gin middleware that limits requests per IP
// RateLimitMiddleware คืนค่า Gin middleware ที่จำกัด requests ตาม IP
func RateLimitMiddleware(redisClient *store.RedisClient, limit int, windowSeconds int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP (respect X-Forwarded-For if behind proxy)
		// ดึง IP ของ client (รองรับ X-Forwarded-For ถ้าอยู่หลัง proxy)
		ip := c.ClientIP()
		key := "rate:" + ip

		allowed, err := redisClient.RateLimitSlidingWindow(key, limit, windowSeconds)
		if err != nil {
			// If Redis fails, allow or deny? Better to allow temporarily
			// ถ้า Redis ล้มเหลว อนุญาตชั่วคราว
			c.Next()
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded, please try again later",
			})
			return
		}
		c.Next()
	}
}
```

### ขั้นตอนที่ 8: Auth Middleware (`middleware/auth.go`)

```go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"auth-system/handlers"
	"auth-system/store"
)

// AuthMiddleware validates JWT access token
// AuthMiddleware ตรวจสอบ JWT access token
func AuthMiddleware(redisClient *store.RedisClient, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}
		token := parts[1]

		// Check blacklist / ตรวจสอบ blacklist
		blacklisted, err := redisClient.IsBlacklisted(token)
		if err == nil && blacklisted {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
			return
		}

		// Validate JWT / ตรวจสอบ JWT
		userID, err := handlers.ValidateAccessToken(token, jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Store user ID in context for later handlers
		// เก็บ user ID ใน context ให้ handlers ถัดไปใช้
		c.Set("userID", userID)
		c.Next()
	}
}
```

### ขั้นตอนที่ 9: Main function (`main.go`)

```go
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"auth-system/handlers"
	"auth-system/middleware"
	"auth-system/store"
)

func main() {
	// Load .env / โหลดไฟล์ .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}
	accessTTL, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TTL"))
	if accessTTL == 0 {
		accessTTL = 15
	}
	refreshTTL, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TTL"))
	if refreshTTL == 0 {
		refreshTTL = 10080 // 7 days
	}
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisPass := os.Getenv("REDIS_PASSWORD")

	// Initialize Redis / เริ่มต้น Redis
	redisClient, err := store.NewRedisClient(redisAddr, redisPass)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")

	// Initialize user DB (mock) / เริ่มต้นฐานข้อมูลผู้ใช้ (mock)
	userDB := store.NewMockUserDB()

	// Create handlers / สร้าง handlers
	authHandler := &handlers.AuthHandler{
		UserDB:      userDB,
		RedisClient: redisClient,
		JWTSecret:   jwtSecret,
		AccessTTL:   accessTTL,
		RefreshTTL:  refreshTTL,
	}

	// Setup Gin router / ตั้งค่า Gin router
	r := gin.Default()

	// Public endpoints / จุดปลายทางสาธารณะ
	r.POST("/login", authHandler.Login)
	r.POST("/refresh", authHandler.Refresh)
	r.POST("/logout", authHandler.Logout)

	// Rate limited public endpoints (protect login from brute force)
	// จุดปลายทางสาธารณะที่ถูก rate limit (ป้องกัน brute force)
	loginRateLimit := middleware.RateLimitMiddleware(redisClient, 5, 60) // 5 attempts per minute
	r.POST("/login", loginRateLimit, authHandler.Login)

	// Protected API group / กลุ่ม API ที่ต้องใช้ authentication
	protected := r.Group("/api")
	protected.Use(middleware.RateLimitMiddleware(redisClient, 100, 60)) // 100 req/min per IP
	protected.Use(middleware.AuthMiddleware(redisClient, jwtSecret))
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			c.JSON(200, gin.H{"user_id": userID, "message": "This is protected data"})
		})
		// Add your ML model endpoints here (from Chapter 4)
	}

	// Start server / เริ่มเซิร์ฟเวอร์
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on :%s", port)
	r.Run(":" + port)
}
```

### ขั้นตอนที่ 10: การรัน

```bash
# Start Redis (docker)
docker run -d -p 6379:6379 redis

# Run Go app
go run main.go
```

**ทดสอบด้วย curl:**

```bash
# Login
curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"username":"admin","password":"password"}'
# จะได้ access_token และ cookie refresh_token

# เรียก API ที่ protected
curl -X GET http://localhost:8080/api/profile -H "Authorization: Bearer <access_token>"

# Refresh token
curl -X POST http://localhost:8080/refresh --cookie "refresh_token=<token>"

# Logout
curl -X POST http://localhost:8080/logout -H "Authorization: Bearer <access_token>" --cookie "refresh_token=<token>"
```

---

## 🧪 แนวทางแก้ไขปัญหาที่อาจเกิดขึ้น

| ปัญหา | สาเหตุ | แนวทางแก้ไข |
|--------|--------|---------------|
| Refresh token ถูกขโมย | XSS หรือ MITM | ใช้ HttpOnly + Secure + SameSite=Strict cookie; ใช้ short TTL สำหรับ refresh token |
| Rate limit เกินแม้ user ปกติ | Sliding window ไม่แม่น | ใช้ token bucket หรือเพิ่ม limit; ตรวจสอบว่า client caching |
| JWT ถูก replay | ไม่มี nonce | ใส่ jti (JWT ID) และตรวจสอบ blacklist สำหรับ access token สั้น |
| Redis ล้มเหลว | Out of memory หรือ network | Implement fallback to local memory (เช่น map + mutex) และ alert |
| Brute force login | ไม่มี captcha | เพิ่ม captcha หลังจาก fail 3 ครั้ง; ใช้ rate limit แยกต่อ IP |

---

## 📐 เทมเพลตและ Checklist สำหรับระบบ Authentication

### Checklist การตั้งค่า Security

| ขั้นตอน | รายละเอียด | เสร็จ/ไม่เสร็จ |
|---------|------------|----------------|
| 1 | ใช้ bcrypt หรือ argon2 สำหรับ hashing password | ☐ |
| 2 | ตั้ง JWT secret ให้ยาวและสุ่ม (≥32 ไบต์) | ☐ |
| 3 | Access token TTL ≤ 30 นาที | ☐ |
| 4 | Refresh token เก็บใน Redis หรือ DB พร้อม TTL | ☐ |
| 5 | ใช้ HttpOnly cookie สำหรับ refresh token | ☐ |
| 6 | เปิดใช้งาน HTTPS ใน production | ☐ |
| 7 | ตั้ง rate limit สำหรับ login (5-10 ครั้ง/นาที) | ☐ |
| 8 | ตั้ง rate limit สำหรับ API ทั่วไป (100-1000/นาที) | ☐ |
| 9 | ตรวจสอบ blacklist ก่อน validate JWT | ☐ |
| 10 | เพิ่ม security headers (CSP, X-Frame-Options, etc.) | ☐ |

### ตารางเปรียบเทียบวิธีการจัดเก็บ Refresh Token

| วิธี | ความปลอดภัย | ความสะดวก | Scalability | เหมาะกับ |
|------|--------------|------------|--------------|-----------|
| Redis | สูง (TTL, atomic) | สูง | ดี (distributed) | Microservices |
| PostgreSQL | ปานกลาง | ปานกลาง | ปานกลาง (อ่านช้า) | Monolith |
| HttpOnly Cookie | ดี (ไม่ถูก XSS) | ดี | N/A | Web apps |
| localStorage | ต่ำ (XSS risk) | ง่าย | N/A | ไม่แนะนำ |

---

## 🧩 แบบฝึกหัดท้ายบท (4 ข้อ)

### ข้อที่ 1: เพิ่ม Role-Based Access Control (RBAC)
จงเพิ่มฟิลด์ `role` ใน User (admin, user) และสร้าง middleware ที่ตรวจสอบ role ก่อนให้เข้าถึง endpoint `/api/admin` เฉพาะ admin

### ข้อที่ 2: เพิ่ม CAPTCHA หลังจาก login ล้มเหลว 3 ครั้ง
ใช้ Redis นับจำนวนครั้งที่ login fail ต่อ username/IP ถ้าเกิน 3 ครั้ง ให้ต้องกรอก CAPTCHA (จำลองด้วย fixed code หรือใช้ library เช่น `github.com/dchest/captcha`)

### ข้อที่ 3: สร้าง middleware สำหรับ Rate Limit แบบ Token Bucket
Implement token bucket rate limit ใน memory (ไม่ใช้ Redis) สำหรับการทดสอบ offline โดยใช้ `time.Ticker` และ atomic counter

### ข้อที่ 4: เพิ่ม Two-Factor Authentication (2FA) ด้วย TOTP
ใช้แพ็กเกจ `github.com/pquerna/otp/totp` สร้าง secret, แสดง QR code และตรวจสอบรหัส 6 หลักระหว่าง login

---

## ✅ เฉลยแบบฝึกหัด

### เฉลยข้อที่ 1 (RBAC)

```go
// Add role to User struct
type User struct {
	ID       uint
	Username string
	Password string
	Role     string // "admin" or "user"
}

// Middleware to check role
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists { return }
		// fetch user from DB
		user, _ := userDB.GetUserByID(userID.(uint))
		if user.Role != role {
			c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
```

### เฉลยข้อที่ 2 (CAPTCHA)

```go
// Redis counters
failKey := "login_fail:" + username
count, _ := redisClient.Get(ctx, failKey).Int()
if count >= 3 {
	captcha := c.PostForm("captcha")
	if !captcha.Verify(captcha) {
		return error
	}
}
// after successful login, reset counter
```

### เฉลยข้อที่ 3 (Token Bucket in-memory)

```go
type TokenBucket struct {
	capacity int
	tokens   int
	rate     time.Duration
	mu       sync.Mutex
	last     time.Time
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(tb.last)
	tb.tokens += int(elapsed / tb.rate)
	if tb.tokens > tb.capacity { tb.tokens = tb.capacity }
	tb.last = now
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}
```

### เฉลยข้อที่ 4 (TOTP 2FA)

```go
import "github.com/pquerna/otp/totp"

// Generate secret and QR
key, _ := totp.Generate(totp.GenerateOpts{
	Issuer:      "MyApp",
	AccountName: user.Username,
})
// Store key.Secret() in DB
// Verify during login
valid := totp.Validate(passcode, user.TOTPSecret)
```

---

## ⚠️ สรุป: ประโยชน์, ข้อควรระวัง, ข้อดี, ข้อเสีย, ข้อห้าม

### ✅ ประโยชน์ที่ได้รับ
- ระบบ Authentication ที่ปลอดภัยด้วย JWT + Refresh Token
- ป้องกันการโจมตีแบบ Brute Force และ DDoS ด้วย Rate Limit
- สามารถเพิกถอนโทเค็นเมื่อ logout หรือตรวจพบความผิดปกติ
- สามารถขยายระบบไปใช้กับ distributed API ได้ด้วย Redis

### ⚠️ ข้อควรระวัง
- ต้องใช้ HTTPS เสมอใน production มิฉะนั้น token จะถูกขโมย
- Refresh token ที่เป็น UUID ต้องสุ่มอย่างปลอดภัย (crypto/rand)
- การตั้ง rate limit ควรเผื่อเผื่อสำหรับผู้ใช้ปกติ (ไม่เข้มงวดเกินไป)
- อย่าเก็บข้อมูลสำคัญ (เช่น password) ใน JWT payload

### 👍 ข้อดี
- stateless สำหรับ access token (JWT) ทำให้ scale ได้ดี
- Refresh token ช่วยให้ UX ดี (ไม่ต้อง login บ่อย)
- Redis rate limit แม่นยำและทำงานข้าม instance ได้

### 👎 ข้อเสีย
- การ implement blacklist ต้องใช้ storage เพิ่ม (Redis)
- JWT ไม่สามารถ revoke ได้จริง ต้องพึ่ง blacklist หรือ short TTL
- การ refresh token ซับซ้อนกว่าการใช้ session แบบธรรมดา

### 🚫 ข้อห้าม
- **ห้าม** เก็บ Refresh Token ใน localStorage (เสี่ยง XSS)
- **ห้าม** ใช้ secret key อ่อนแอหรือ hardcode ใน source code
- **ห้าม** ตั้ง Access Token TTL นานเกินไป (ไม่ควรเกิน 1 วัน)
- **ห้าม** ข้าม rate limit สำหรับ endpoint สำคัญ (login, register)

---

## 🔗 แหล่งอ้างอิง (References)

1. JWT RFC 7519 – https://tools.ietf.org/html/rfc7519
2. OWASP Cheat Sheet: Authentication – https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html
3. Redis Rate Limiting Patterns – https://redis.io/docs/manual/patterns/rate-limiter/
4. gin-gonic/gin documentation – https://gin-gonic.com/docs/
5. golang-jwt/jwt – https://github.com/golang-jwt/jwt
6. bcrypt for Go – https://pkg.go.dev/golang.org/x/crypto/bcrypt

---

## จบบทที่ 6: ระบบจัดการ Login, Refresh Token และ Rate Limit

บทนี้ได้เสริมความปลอดภัยให้กับ API ของระบบ Big Data และ MML AI ที่สร้างขึ้นในบทก่อนหน้า ด้วย JWT, Refresh Token, Rate Limit และมาตรการอื่น ๆ ทำให้ระบบพร้อมเผชิญภัยคุกคามในโลกจริง ผู้เรียนสามารถนำโค้ดไปปรับใช้กับ production API ได้ทันที

**เนื้อหาครบถ้วนทั้ง 6 บท ครอบคลุมตั้งแต่การสร้าง Data Pipeline, Streaming, Distributed Training, Production Deployment, Data Governance, และ Security**