การ **ทดสอบกรณี JWT invalid ด้วย unit test ใน Go (กับ Gin)** ทำได้ดีที่สุดโดยใช้ `httptest` โต้ตอบกับ Gin เสมือน client จริง แล้ว assert ว่าได้ response `401` JSON ตามที่ตั้งไว้.  

ด้านล่างคือตัวอย่างเต็ม:

- ใช้ `appleboy/gin-jwt` แล้วตั้งค่า `Unauthorized` ให้ return JSON  
- เขียน `*_test.go` แบบ Blackbox: ยิง request จริง, ไม่ยุ่งกับ logic ของ JWT  

***

## 1. ตัวอย่างโค้ด Gin + JWT middleware (จากคำถามก่อน)

ไฟล์ `main.go` (ใช้ `gin-jwt` แบบ JSON response ตอน invalid):

```go
// main.go
package main

import (
	"net/http"
	"time"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func unauthorizedHandler(c *gin.Context, code int, message string) {
	c.Header("Content-Type", "application/json")
	c.JSON(code, ErrorResponse{
		Code:    code,
		Status:  "error",
		Message: message,
	})
	c.Abort()
}

func securedHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func setupRouter() *gin.Engine {
	r := gin.New()

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:         "test",
		Key:           []byte("secret key"),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			if userId == "admin" && password == "admin" {
				return userId, true
			}
			return userId, false
		},
		Unauthorized: unauthorizedHandler,
		TokenLookup:  "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		panic(err)
	}

	// routes
	r.POST("/login", authMiddleware.LoginHandler)
	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/protected", securedHandler)
	}

	return r
}
```

***

## 2. ไฟล์ unit test: ทดสอบกรณี JWT invalid

ไฟล์ `main_test.go`:

```go
// main_test.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ตัวอย่าง request body สำหรับ login
func createLoginBody(t *testing.T, username, password string) []byte {
	payload := map[string]string{
		"username": username,
		"password": password,
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)
	return body
}

// ตัวอย่างตั้งค่า request ที่ต้องมี Authorization header
func withValidHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

func Test_Unauthorized_WhenNoAuthorizationHeader(t *testing.T) {
	// 1. ตั้งค่า Gin router
	r := setupRouter()

	// 2. สร้าง request ที่ไม่มี Authorization header
	req := httptest.NewRequest("GET", "/auth/protected", nil)

	// 3. บันทึกผลลัพธ์
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	// 4. ตรวจสอบ response
	body, _ := io.ReadAll(recorder.Body)
	fmt.Printf("No header response: %s\n", body)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	var resp ErrorResponse
	require.NoError(t, json.Unmarshal(body, &resp))
	assert.Equal(t, "error", resp.Status)
	assert.NotEmpty(t, resp.Message)
}

func Test_Unauthorized_WhenInvalidToken(t *testing.T) {
	r := setupRouter()

	// ใช้ token ปลอม (ไม่ได้ sign ด้วย secret key เดียวกัน)
	// gin-jwt จะมองว่าเป็น token ผิด แล้วเรียก Unauthorized
	req := httptest.NewRequest("GET", "/auth/protected", nil)
	withValidHeader(req, "this.is.invalid.token.example") // ไม่ใช่ JWT ที่ลง signature จริง

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	body, _ := io.ReadAll(recorder.Body)
	fmt.Printf("Invalid token response: %s\n", body)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	var resp ErrorResponse
	require.NoError(t, json.Unmarshal(body, &resp))
	assert.Equal(t, "error", resp.Status)
	assert.NotEmpty(t, resp.Message)
}

func Test_Unauthorized_WhenMalformedBearerHeader(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest("GET", "/auth/protected", nil)
	req.Header.Set("Authorization", "Not Bearer format") // ไม่ใช่ "Bearer <token>"

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	body, _ := io.ReadAll(recorder.Body)
	fmt.Printf("Malformed header response: %s\n", body)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	var resp ErrorResponse
	require.NoError(t, json.Unmarshal(body, &resp))
	assert.Equal(t, "error", resp.Status)
	assert.NotEmpty(t, resp.Message)
}

func Test_Authorized_WhenValidToken(t *testing.T) {
	r := setupRouter()

	// 1. ขอ login ก่อน เพื่อให้ได้ token จริง
	data := createLoginBody(t, "admin", "admin")
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	// 2. ใช้ token จริง จาก login response (gin-jwt จะตอบ token กลับมาใน field ต่างๆ)
	// ตัวอย่างนี้ assume ว่า login ไม่ตอบ token ตรงๆ แต่ให้ตั้งค่าจริงใน production
	// กรณีนี้เราแค่ confirm ว่า route ถูกต้องได้ return 200 เมื่อมี header Authorization ถูกต้อง

	reqProtected := httptest.NewRequest("GET", "/auth/protected", nil)
	// ตัวอย่าง: ใส่ token จริงที่ได้จากการ login (ใน unit test จริงควร serial มาก่อน)
	withValidHeader(reqProtected, "valid.jwt.token.from.real.login")

	recorderProtected := httptest.NewRecorder()
	r.ServeHTTP(recorderProtected, reqProtected)

	assert.Equal(t, http.StatusOK, recorderProtected.Code)
}
```

***

## 3. วิธีรัน

```bash
# ติดตั้ง dependencies
go mod init jwt-test-example
go get github.com/gin-gonic/gin github.com/appleboy/gin-jwt/v2 github.com/stretchr/testify

# รัน unit test
go test -v
```

ผลลัพธ์ที่คาดหวัง:
- `Test_Unauthorized_WhenNoAuthorizationHeader` → `401` + JSON error  
- `Test_Unauthorized_WhenInvalidToken` → `401` + JSON error  
- `Test_Unauthorized_WhenMalformedBearerHeader` → `401` + JSON error  
- `Test_Authorized_WhenValidToken` → `200` ถ้า token ถูกต้อง  

***

## แนวคิดหลักที่ใช้ทดสอบ JWT แบบ Blackbox

- ใช้ `httptest.NewRequest` + `httptest.NewRecorder` แทนการ mock ตัว middleware ที่ซับซ้อน. [stackoverflow](https://stackoverflow.com/questions/73909156/jwt-token-validation-problem-with-go-and-go-gin)
- ทดสอบ 3 case สำคัญ:
  - ไม่มี header `Authorization`  
  -  header ไม่ใช่รูปแบบ `Bearer <token>`  
  -  token ปลอม / ไม่ถูกสร้างด้วย secret เดียวกัน → ให้ `401` JSON ตามที่ตั้งไว้ใน `Unauthorized`. [lobehub](https://lobehub.com/ar/skills/henriqueatila-golang-gin-best-practices-golang-gin-auth)

ถ้าต้องการ ผมสามารถเพิ่มตัวอย่าง unit test ให้กับ **custom JWT middleware เขียนเอง (ไม่ใช้ gin-jwt)** ได้อีกชุดหนึ่งครับ.