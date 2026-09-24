การ **จัดการ error เมื่อ JWT token ไม่ถูกต้องใน middleware** ทำได้ดีที่สุดโดย:

1. ตั้ง `Unauthorized` ของ gin‑jwt ให้ return  response ชัดเจน  
2. จบ middleware ด้วย `c.Abort()` เพื่อไม่ให้ไปต่อ handler  
3. ใช้ `Binding errors` หรือ `c.Error(...)` สำหรับ custom error handling  

ด้านล่างคือตัวอย่างโค้ดแบบใช้ `appleboy/gin-jwt` (Gin) แบบเต็ม พร้อมจัดการทุก error case ที่เกี่ยวกับ JWT. [pkg.go](https://pkg.go.dev/github.com/appleboy/gin-jwt/v2)

***

## ตัวอย่างโค้ด: จัดการ error JWT ใน Gin middleware

ไฟล์ `main.go`:

```go
// main.go
// ตัวอย่าง: จัดการ error ทุกกรณีของ JWT token ใน Gin middleware
// ใช้ github.com/appleboy/gin-jwt/v2

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// Response สำหรับ error handling
type ErrorResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ฟังก์ชัน Unauthorized ใช้ตั้งค่า response เมื่อ JWT ผิดพลาด
// ใช้ handler นี้จัดการทุกกรณี:
//  - ไม่มี token
//  - token  grammar ผิด (malformed)
//  - signature ผิด (invalid signature)
//  - expired / not yet valid
func unauthorizedHandler(c *gin.Context, code int, message string) {
	// ใช้ gin.Error เพื่อให้ error ถูกจัดการโดย recovery / logging
	c.Error(fmt.Errorf("JWT unauthorized: %d - %s", code, message))

	c.JSON(code, ErrorResponse{
		Code:    code,
		Status:  "error",
		Message: message,
	})
	c.Abort() // จบ middleware chain ทันที
}

// Authenticator ใช้ตั้งค่าตอน login
func loginHandler(c *gin.Context) {
	// ตัวอย่าง login ง่ายๆ (ไม่ใช่ production)
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "admin" && password == "admin" {
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

// Handler ที่ต้องมี JWT ถูกต้องเท่านั้น
func securedHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)

	// ตัวอย่าง: ใช้ id หรือ email เป็น user ID
	userId, exists := claims["id"]
	if !exists {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    500,
			Status:  "error",
			Message: "Missing user id in token claims",
		})
		c.Abort()
		return
	}

	uid, ok := userId.(string)
	if !ok {
		uid = fmt.Sprintf("%v", userId)
	}

	c.JSON(200, gin.H{
		"status":   "success",
		"message":  "Access granted",
		"user_id":  uid,
		"jti":      claims["jti"],
		"exp_time": time.Unix(int64(claims["exp"].(float64)), 0).Format(time.RFC3339),
	})
}

func main() {
	r := gin.New()

	// ตั้งค่า JWT middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:         "test-realm",
		Key:           []byte("secret key"), // ใช้ env จริง
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			if (userId == "admin" && password == "admin") ||
			   (userId == "test" && password == "test") {
				return userId, true
			}
			return userId, false
		},
		// ใช้ Unauthorized จัดการ error ทุกกรณีของ JWT
		Unauthorized: unauthorizedHandler,
		// ตัวอย่าง: ใช้ต่อไปได้หากไม่อยาก restrict role
		Authorizator: func(userId string, c *gin.Context) bool {
			return true // ปล่อยทุก user ที่ auth ผ่าน
		},
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT middleware setup failed: ", err)
	}

	// ตั้งค่า routes
	r.POST("/login", authMiddleware.LoginHandler)

	// กลุ่มที่ต้อง.JWT auth
	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc()) // ใช้ middleware ตรวจสอบ JWT
	{
		auth.GET("/profile", securedHandler)
		auth.GET("/refresh", authMiddleware.RefreshHandler)
	}

	// ตัวอย่าง custom error handling (optional)
	r.Use(func(c *gin.Context) {
		c.Next()

		// ถ้า handler หรือ middleware ใช้ c.Error(...) จะมี error list
		if len(c.Errors) > 0 {
			// ตัวอย่าง: ตั้งค่า HTTP status ทั่วไปตาม error แรก
			firstErr := c.Errors[0]

			if gin.IsDebugging() {
				c.JSON(firstErr.Meta.(int), gin.H{
					"status":  "error",
					"message": firstErr.Error(),
					"code":    firstErr.Meta,
				})
			}
		}
	})

	// รัน
	r.Run(":8080")
}
```

***

## กรณี error ต่างๆ ที่ gin‑jwt จัดการ / สามารถส่งให้ `Unauthorized`

- ❌ ไม่มี token หรือไม่มี header `Authorization`  
- ❌ 	token ไม่ใช่รูปแบบ “Bearer JWT” (malformed)  
- ❌ 	signature ไม่ถูกต้อง  
- ❌ 	token หมดอายุ (`exp`) หรือ `nbf` ยังไม่ถึงเวลา  
- ❌ 	token ถูกตั้งค่า claims ผิด (เช่น realm ผิด, audience ผิด ถ้าตั้งค่า.) [dev](https://dev.to/leapcell/jwt-in-action-secure-authentication-authorization-in-go-jde)

gin‑jwt จะ **เรียก `Unauthorized(...)` ทุกกรณีเหล่านี้** ทำให้คุณเขียนโค้ด handle error ครั้งเดียว ไม่ต้องแยก `if/else` ในทุก handler.

***

## แนวทางดีที่สุด (Best practices)

- ใช้ `Unauthorized` ตั้ง code ที่เหมาะสม:  
  - `401 Unauthorized` ถ้าไม่มี token / token ไม่ถูกต้อง  
  - `403 Forbidden` ถ้า JWT ถูกต้องแต่ user ไม่ได้สิทธิ์ (ใช้ `Authorizator`). [stackoverflow](https://stackoverflow.com/questions/69948784/how-to-handle-errors-in-gin-middleware)
- ใช้ `c.JSON(code, ...); c.Abort()` ทุกครั้งใน middleware ถ้าพบ error สำคัญ → ไม่ต้องให้ไปถึง handler.  
- ใช้ `c.Error(...)` ถ้าอยากให้ error ถูกจัดการโดย custom error middleware หรือ logging.

 