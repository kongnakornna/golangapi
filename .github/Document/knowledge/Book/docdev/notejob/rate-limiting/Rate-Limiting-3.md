ด้านล่างคือตัวอย่างโค้ด **“extract user ID จาก JWT token แล้วใส่ลงใน Gin context”** แบบใช้ได้จริง ใช้ library ยอดนิยม [`appleboy/gin-jwt`](https://github.com/appleboy/gin-jwt) ซึ่งใช้กันเยอะใน Gin ecosystem. [github](https://github.com/appleboy/gin-jwt)

***

## ตัวอย่าง 1: ใช้ `appleboy/gin-jwt` ดึง user ID จาก JWT แล้วใส่ลง context

ไฟล์ `main.go`:

```go
// main.go
// ตัวอย่าง extract user ID จาก JWT token แล้วใช้ใน Gin middleware
// ใช้ library: appleboy/gin-jwt

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// ตัวอย่าง response
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	UserID  string `json:"user_id,omitempty"`
}

// ตัวอย่าง handler ที่ใช้ user ID จาก JWT
func securedHandler(c *gin.Context) {
	// ใช้ gin-jwt ดึง claims (payload ของ JWT token)
	claims := jwt.ExtractClaims(c)

	// สมมุติว่าตอนสร้าง token ได้ใส่ "id" หรือ "user_id" ลง claims
	userId, exists := claims["id"]
	if !exists {
		userId = "anonymous"
	}

	// แปลงเป็น string
	uid, ok := userId.(string)
	if !ok {
		uid = fmt.Sprintf("%v", userId) // ถ้าไม่ใช่ string
	}

	// ตัวอย่าง: ใส่ user ID ลง context สำหรับ middleware ต่อไป (เช่น rate limit)
	c.Set("user_id", uid)

	c.JSON(200, Response{
		Status:  "success",
		Message: "This is a secured endpoint",
		UserID:  uid,
	})
}

// ตัวอย่าง login handler (ใช้ gin-jwt built‑in)
func loginHandler(c *gin.Context) {
	// ตัวอย่างการ authenticate จริงควรเช็ค DB จริง
	// ตัวอย่างนี้ ใช้ username/password แบบ hardcode demo
	username := c.PostForm("username")
	password := c.PostForm("password")

	// ตัวอย่าง logic ง่ายๆ (ไม่ใช่ best practice จริง)
	if (username == "admin" && password == "admin") || (username == "test" && password == "test") {
		// gin-jwt จะ generate JWT ให้ตาม config
		// ถ้าใช้ gin-jwt อย่างถูกต้อง ไม่ต้องเขียน encode ด้วยตัวเอง
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Login successful",
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"status":  "error",
		"message": "Invalid credentials",
	})
}

func main() {
	r := gin.New()

	// ตั้งค่า JWT middleware (gin-jwt)
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		// ชื่อ realm (ใช้แสดงใน error / log)
		Realm: "test zone",

		// secret key สำหรับ sign JWT (ควรอ่านจาก env จริง)
		Key: []byte("secret key"),

		// ระยะเวลา token หมดอายุ
		Timeout:  time.Hour,
		MaxRefresh: time.Hour,

		// ฟังก์ชันตรวจสอบ credential ตอน login
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			if (userId == "admin" && password == "admin") || (userId == "test" && password == "test") {
				return userId, true // ใช้ userId เป็น claim หลัก
			}
			return userId, false
		},

		// ฟังก์ชันตรวจสอบว่า user นี้มีสิทธิ์ใช้ endpoint หรือไม่ (optional)
		Authorizator: func(userId string, c *gin.Context) bool {
			// ตัวอย่าง: ให้เฉพาะ admin เข้า
			// ถ้าไม่ต้องการ restrict ให้ return true เลย
			return userId == "admin"
		},

		// ฟังก์ชันที่เรียกเมื่อไม่ผ่าน auth
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},

		// บอกว่า token ดึงมาจากไหน (header:Authorization)
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",

		// ใช้ time.Now ปกติ
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT middleware setup failed: ", err)
	}

	// ตั้งค่า routes
	r.POST("/login", authMiddleware.LoginHandler)

	// กลุ่ม route ที่ใช้ JWT auth
	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc()) // ใช้ middleware ตรวจสอบ JWT
	{
		auth.GET("/hello", securedHandler)          // handler ใช้ user ID จาก claims
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	r.Run(":8080")
}
```

***

## วิธีการ extract user ID แบบละเอียด

1. **ตอน login**  
   - ใช้ `authMiddleware.LoginHandler` ของ gin‑jwt.  
   - ฟังก์ชัน `Authenticator` คืน `string userId` ค่าที่จะใส่ลง JWT claims (เช่น `id` หรือ `user_id`).  

2. **ตอนใช้ endpoint ที่ต้องการ JWT**  
   - ใช้ `authMiddleware.MiddlewareFunc()` วางเป็น Gin middleware.  
   - ใช้ `jwt.ExtractClaims(c)` ดึง payload ของ JWT ออกมาเป็น `map[string]interface{}`.  
   - ดึงค่า `claims["id"]` หรือ `claims["user_id"]` แล้วเอาไปตั้ง `c.Set("user_id", ...)` สำหรับ middleware ต่อไป เช่น rate limit แยกตาม user ID. [pkg.go](https://pkg.go.dev/github.com/orange-health/gin-jwt)

***

## ตัวอย่างการใช้ร่วมกับ rate limit middleware

หลังจาก extract user ID แล้ว คุณสามารถใช้โครงสร้างแบบนี้:

```go
auth.GET("/data", func(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userId, _ := claims["id"].(string)
	c.Set("user_id", userId)

	// ต่อไป middleware rate limit ต่อ user ID ที่เคยตั้งไว้จะใช้ c.Get("user_id") ได้
	c.JSON(200, gin.H{"user_id": userId, "message": "OK"})
})
```
 