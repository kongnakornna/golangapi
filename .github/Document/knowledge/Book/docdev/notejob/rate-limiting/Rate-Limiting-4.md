ด้านล่างคือตัวอย่างโค้ด **การ return JSON error response เมื่อ JWT invalid** ใช้กับ Gin ทั้งแบบ `appleboy/gin-jwt` และ custom middleware ที่เขียนเอง ใช้ได้ทันที. [dev](https://dev.to/leapcell/jwt-in-action-secure-authentication-authorization-in-go-jde)

***

## ตัวอย่าง 1: ใช้ `appleboy/gin-jwt` (return 401 JSON)

ไฟล์ `main.go`:

```go
// main.go
// ตัวอย่าง: คืน JSON error response เมื่อ JWT ผิดพลาด (invalid / malformed / expired)

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// JSON Error Response ใช้สำหรับ Unauthorized
type ErrorResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ฟังก์ชันตั้งค่า response ตอน JWT ไม่ถูกต้อง
// ใช้ตั้งค่าทุกกรณี:
//  - ไม่มี token
//  - token ไม่ถูกต้อง / malformed
//  - signature ผิด
//  - expired / not yet valid
func unauthorizedHandler(c *gin.Context, code int, message string) {
	c.Header("Content-Type", "application/json")

	// ตัวอย่าง JSON response
	resp := ErrorResponse{
		Code:    code,
		Status:  "error",
		Message: message,
	}

	if gin.IsDebugging() {
		resp.Error = "Please check Authorization header: Bearer <JWT>"
	}

	c.JSON(code, resp)
	c.Abort() // จบ middleware chain
}

// Handler ที่ต้องใช้ JWT
func securedHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Access granted",
	})
}

func loginHandler(c *gin.Context) {
	// ตัวอย่าง login ง่ายๆ
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

func main() {
	r := gin.New()

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
		// ใช้ Unauthorized ตั้งค่า response JSON ตอน JWT invalid
		Unauthorized: unauthorizedHandler,
		Authorizator: func(userId string, c *gin.Context) bool {
			return true
		},
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT middleware setup failed: ", err)
	}

	// routes
	r.POST("/login", authMiddleware.LoginHandler)

	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc()) // ใช้ middleware ตรวจสอบ JWT
	{
		auth.GET("/profile", securedHandler)
	}

	// รัน server
	r.Run(":8080")
}
```

เมื่อ JWT ผิดพลาด ทุก request จะได้เช่น:

```http
HTTP 401 Unauthorized
Content-Type: application/json

{
  "code": 401,
  "status": "error",
  "message": "invalid or expired jwt",
  "error": "Please check Authorization header: Bearer <JWT>"
}
```

***

## ตัวอย่าง 2: Custom JWT middleware เขียนเอง แล้ว return JSON

ไฟล์ `custom_jwt.go`:

```go
// custom_jwt.go
// ตัวอย่าง custom middleware สำหรับ JWT แบบเขียนเอง
// ใช้ github.com/golang-jwt/jwt/v4 (ถ้าต้องการ ติดตั้ง: go get github.com/golang-jwt/jwt/v4)
package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ฟังก์ชัน verify JWT แบบง่าย (ตัวอย่าง ใช้ HSM / secret key)
// จริงๆ ควรใช้ library อย่าง golang-jwt/jwt
func verifyToken(tokenString string) (string, error) {
	// ตัวอย่างง่ายๆ: ตรวจสอบ format อย่างคร่าวๆ
	if !strings.HasPrefix(tokenString, "ey") {
		return "", errors.New("invalid token format")
	}

	// สมมุติว่า token ถูกต้อง แล้วคืน user ID
	return "user-123", nil // ตัวอย่าง user ID จาก claims
}

// JWTMiddleware แบบ custom ที่ return JSON ถ้า JWT ผิดพลาด
func JWTMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ดึง token จาก header: Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Status:  "error",
				Message: "Authorization header required",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Status:  "error",
				Message: "Authorization must be Bearer <JWT>",
			})
			c.Abort()
			return
		}

		token := parts [dev](https://dev.to/leapcell/jwt-in-action-secure-authentication-authorization-in-go-jde)

		// ตรวจสอบ token จริงๆ ควรใช้ jwt.Parse หรือ library จริง
		userId, err := verifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Status:  "error",
				Message: "Invalid or expired JWT token",
				Error:   err.Error(),
			})
			c.Abort()
			return
		}

		// ถ้าถูกต้อง ตั้ง user ID ลง context
		c.Set("user_id", userId)

		// ไปต่อ middleware / handler ต่อไป
		c.Next()
	}
}

func main() {
	r := gin.Default()
	secret := "my-secret-key"

	// ใช้ custom JWT middleware
	r.Use(JWTMiddleware(secret))

	r.GET("/protected", func(c *gin.Context) {
		userId, _ := c.Get("user_id")
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Access granted",
			"user_id": userId,
		})
	})

	r.Run(":8080")
}
```

***

## สรุปแนวทาง

- ใช้ `c.JSON(code, ErrorResponse{...}); c.Abort()` เมื่อ JWT ไม่ถูกต้อง ทั้งแบบ `gin-jwt` (ผ่าน `Unauthorized`) หรือ custom middleware. [tillitsdone](https://tillitsdone.com/blogs/gin-jwt-authentication-guide/)
- ใช้ `http.StatusUnauthorized` (401) สำหรับกรณี JWT ไม่ถูกต้อง / หมดอายุ ที่ common ที่สุด.  

 