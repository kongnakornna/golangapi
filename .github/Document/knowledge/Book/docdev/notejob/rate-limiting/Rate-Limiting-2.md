การใช้ **rate limit แยกตาม User ID ใน Gin middleware** ทำได้ง่าย โดยใช้ `golang.org/x/time/rate` แล้วดึง `userID` จาก JWT / session / header แล้วเอา `userId` ไปเป็น key ของ token bucket ต่อ user. [github](https://github.com/ljahier/gin-ratelimit)

ด้านล่างคือตัวอย่างโค้ด “เต็มรูปแบบ” ใช้ได้ทันที.

***

## ตัวอย่างโค้ด: rate limit แยกตาม User ID (TokenBucket + Gin)

ไฟล์ `main.go`:

```go
// main.go
// Rate limit แยกตาม User ID ใช้กับ Gin
// ใช้ token bucket ต่อ user (user‑based rate limit)

package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Response สำหรับ API
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	UserID  string `json:"user_id,omitempty"`
}

// UserIdRateLimiter สำหรับ rate limit ต่อ user ID
// โครงสร้าง:
//  - ใช้ map จัดเก็บ rate.Limiter ต่อ user ID
//  - ใช้ sync.RWMutex ป้องกัน race condition
type UserIdRateLimiter struct {
	limiterMap map[string]*rate.Limiter
	rate       rate.Limit // requests per second
	burst      int
	mu         sync.RWMutex
}

// NewUserIdRateLimiter สร้าง rate limiter ใหม่ต่อ user
// rps = จำนวน request ต่อวินาทีต่อ user, burst = ความจุ bucket สูงสุด
func NewUserIdRateLimiter(rps float64, burst int) *UserIdRateLimiter {
	return &UserIdRateLimiter{
		limiterMap: make(map[string]*rate.Limiter),
		rate:       rate.Limit(rps),
		burst:      burst,
	}
}

// getLimiter ดึง token bucket สำหรับ user ID นี้
func (ul *UserIdRateLimiter) getLimiter(userId string) *rate.Limiter {
	ul.mu.RLock()
	limiter, exists := ul.limiterMap[userId]
	ul.mu.RUnlock()

	if exists {
		return limiter
	}

	ul.mu.Lock()
	defer ul.mu.Unlock()

	if limiter, exists = ul.limiterMap[userId]; exists {
		return limiter
	}

	// สร้าง token bucket ใหม่ต่อ user ID
	limiter = rate.NewLimiter(ul.rate, ul.burst)
	ul.limiterMap[userId] = limiter
	return limiter
}

// Allow ใช้ตรวจสอบว่า user นี้สามารถทำ request ได้หรือไม่
func (ul *UserIdRateLimiter) Allow(userId string) bool {
	return ul.getLimiter(userId).Allow()
}

// Middleware สำหรับ Gin ใช้แยกตาม User ID
// ต้องมีการ authentication ก่อน แล้วตั้ง userId ลงใน context
func (ul *UserIdRateLimiter) Middleware(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ตัวอย่างการดึง User ID:
		//  - จาก JWT: c.MustGet("user_id").(string)
		//  - จาก header/API Key mapping: c.GetHeader("X-User-ID")
		// วิธีนี้แล้วแต่ design ระบบของคุณ
		userId, exists := c.Get("user_id")
		if !exists {
			// ถ้ายังไม่ได้ authenticate มาก่อน ให้ตั้งเป็น anonymous หรือ reject
			// ตัวอย่างนี้ ใช้ anonymous แค่สำหรับ demo
			userId = "anonymous"
		}

		// แปลงเป็น string (ถ้าจำเป็น)
		uid, ok := userId.(string)
		if !ok {
			uid = fmt.Sprintf("%v", userId) // ถ้าเป็นแบบอื่น เช่น int
		}

		// ตรวจสอบ rate limit ต่อ user ID
		if !ul.Allow(uid) {
			c.Header("Content-Type", "application/json")
			c.Header("Retry-After", "1")
			c.Header("X-RateLimit-Limit", fmt.Sprint(int(ul.rate)))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", fmt.Sprint(time.Now().Add(time.Second).Unix()))

			c.JSON(http.StatusTooManyRequests, Response{
				Status:  "error",
				Message: "Too many requests. Please try again later.",
				UserID:  uid,
			})
			c.Abort()
			return
		}

		// ถ้าผ่าน ให้ต่อไปยัง handler ถัดไป
		next(c)
	}
}

// ตัวอย่าง handler ที่ใช้ user ID (จำเป็นต้องมี authentication middleware ก่อน)
func dataHandler(c *gin.Context) {
	// สมมุติว่า authentication middleware ตั้งค่า user_id ลงใน context ไว้แล้ว
	userId, _ := c.Get("user_id")
	uid, _ := userId.(string)

	c.JSON(200, Response{
		Status:  "success",
		Message: "Data retrieved successfully",
		UserID:  uid,
	})
}

// ตัวอย่าง login / authenticate จำลอง
// จริงๆ ควรใช้ JWT หรือ session จริง
func loginHandler(c *gin.Context) {
	username := c.PostForm("username")

	// ตัวอย่าง: ใช้ username เป็น user ID
	// จริงควรใช้ user.ID จาก DB
	c.Set("user_id", username)

	c.JSON(200, Response{
		Status:  "success",
		Message: "Login successful",
		UserID:  username,
	})
}

func main() {
	r := gin.Default()

	// ตั้งค่า rate limiter ต่อ user ID
	// ตัวอย่าง: 10 req/sec ต่อ user, burst 20 requests
	userLimiter := NewUserIdRateLimiter(10.0, 20)

	// ตัวอย่าง route สำหรับ login (ไม่ถูก rate limit ต่อ user ยังไม่เข้าระบบ)
	r.POST("/login", loginHandler)

	// กลุ่ม route ที่ใช้ user ID และ rate limit
	api := r.Group("/api")
	// ตัวอย่าง middleware จำลอง authentication
	api.Use(func(c *gin.Context) {
		// ตัวอย่าง: สมมุติว่า client ต้อง login มาก่อน
		// จริงควรใช้ JWT middleware แล้วตั้งค่า user_id ลงใน context จริง

		// ตัวอย่างนี้ ใช้ header X-User-ID จำลอง (DEMO ONLY)
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}
		c.Set("user_id", userID)
		c.Next()
	})

	// ใช้ rate limit middleware ต่อ user ID กับ /api ทั้งหมด
	api.Use(func(c *gin.Context) {
		userLimiter.Middleware(func(c *gin.Context) {
			c.Next() // ต่อไปยัง handler จริง
		})(c) // ต้องใช้ pattern นี้เพราะ Middleware คืน gin.HandlerFunc
	})

	api.GET("/data", dataHandler)

	// รัน server
	r.Run(":8080")
}
```

***

## วิธีใช้จริง (production‑style)

- ใส่ **authentication middleware** ก่อน แล้วตั้ง `c.Set("user_id", user.ID)` ลงใน `gin.Context`. [github](https://github.com/JGLTechnologies/gin-rate-limit)
- ใช้ `userLimiter.Middleware(...)` สำหรับ route ที่ต้องการ rate limit ต่อ user.  
- ใช้ `rate.NewLimiter(rps, burst)` เป็น **Token Bucket ต่อ user ID** ทำให้แต่ละ user ถูกจำกัดแยกกัน ไม่ชนกัน. [linkedin](https://www.linkedin.com/pulse/improve-performance-security-go-rate-limiting-aslam-mulla-vs2ff)

***

## ตัวอย่างการทดสอบ

```bash
# 1. ตั้งค่า user ID ด้วย header จำลอง
curl -H "X-User-ID: user-123" http://localhost:8080/api/data
curl -H "X-User-ID: user-456" http://localhost:8080/api/data

# ถ้ารันเร็วเกิน 10 req/sec ต่อ user จะได้
# HTTP 429 Too Many Requests
```
 