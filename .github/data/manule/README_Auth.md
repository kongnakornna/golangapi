# คู่มือการใช้งานโมดูล Auth (Authentication & Authorization)

โมดูล `internal/modules/auth` + `internal/modules/users` เป็นระบบยืนยันตัวตน (JWT RS256) ของ `icmongolang`: login/signin, refresh token, logout, verify email, forgot/reset password และการตรวจสิทธิ์ผ่าน middleware chain

**ระบบ token:**
- **Access token** — RS256, หมดอายุตาม `AccessTokenExpireDuration` (default 3600 นาที), ใช้กับ endpoint ทั่วไป (`Verifier(true)`)
- **Refresh token** — RS256, หมดอายุตาม `RefreshTokenExpireDuration` (default 1440 นาที), ใช้กับ `/auth/refresh`, `/auth/logout`, `/auth/logoutall` (`Verifier(false)`) และ **ถูกเก็บใน Redis** (เซต `RefreshToken:{userID}`) เพื่อ rotation/revoke

> **สถาปัตยกรรมสำคัญ:** `authUseCase`/`authPgRepo` เป็นแค่ delegate ที่ถูก discard (`_ = ...`) ใน `handlers.go:147-150` — **โค้ดจริงทั้งหมดอยู่ใน users module** (handler ของ auth เรียก `usersUC` โดยตรง)

---

## 1. ตารางไฟล์

| ไฟล์ | บทบาท |
|------|-------|
| `internal/modules/auth/delivery/http/routes.go` | Route `/api/auth/*` |
| `internal/modules/auth/delivery/http/handlers.go` | HTTP handlers (Login, Refresh, Verify, Reset, ...) |
| `internal/modules/auth/presenter/presenters.go` | DTOs (นิยามแต่ handler ไม่ใช้ตรงๆ) |
| `internal/modules/auth/usecase/usecase.go` | Delegate → usersUC |
| `internal/modules/users/usecase/usecase.go` | **logic จริง**: createToken, SignIn, Refresh, Verify, ForgotPassword, ResetPassword |
| `internal/modules/users/repository/pg_repository.go` | GORM + Redis repo |
| `internal/middleware/jwtauth.go` | Verifier/Authenticator/CurrentUser/ActiveUser/SuperUser |
| `pkg/jwt/token.go` | สร้าง/parse JWT RS256 |
| `pkg/cryptpass/password.go` | bcrypt hash/compare |
| `internal/modules/users/distributor/` | Enqueue email task (asynq) |

---

## 2. การตั้งค่า Configuration

### 2.1 JWT — `config/config.default.yml`
```yaml
jwt:
  SecretKey: ""
  Issuer: go-dev
  AccessTokenExpireDuration: 3600    # ⚠️ หน่วยเป็น นาที
  AccessTokenPrivateKey: ""          # base64 PEM (RS256)
  AccessTokenPublicKey: ""
  RefreshTokenExpireDuration: 1440   # นาที
  RefreshTokenPrivateKey: ""
  RefreshTokenPublicKey: ""
```

> `config.dev.yml` มี RSA key pairs + `firstSuperUser` (root@gmail.com / root / root_password) ให้ใช้งาน dev

### 2.2 Redis (สำหรับ refresh token + user cache)
```yaml
redis:
  Addr: localhost:6379
  Db: 0
```

### 2.3 Email (สำหรับ verify/reset)
```yaml
email:
  From: "kongnakornna@gmail.com"
  VerificationSubject: "Your account verification code"
  ResetSubject: "Your account password reset token"
```

---

## 3. API Reference

ทุก route อยู่ใต้ **`/api/auth`** ยกเว้นกลุ่ม register/login ของ users module

### 3.1 ตารางรวม endpoints

| Method | Path | Auth | คำอธิบาย |
|--------|------|------|----------|
| POST | `/api/auth/login` | ❌ (rate limit 10/s, burst 20) | login ด้วย username **หรือ** email |
| POST | `/api/auth/signin` | ❌ | login ด้วย email |
| GET | `/api/auth/publickey` | ❌ | ดู public key (access + refresh) |
| GET | `/api/auth/verifyemail?code=` | ❌ | ยืนยันอีเมลด้วย code |
| POST | `/api/auth/forgotpassword` | ❌ | ขอ reset password (email) |
| PATCH | `/api/auth/resetpassword` | ❌ | ตั้งรหัสใหม่ด้วย token |
| GET | `/api/auth/refresh` | ✅ **refresh token** | หมุน refresh token → ได้คู่ token ใหม่ |
| GET | `/api/auth/logout` | ✅ refresh token | revoke refresh token ตัวเดียว |
| GET | `/api/auth/logoutall` | ✅ refresh token | revoke refresh token ทั้งหมดของ user |
| POST | `/api/register` | ❌ | register ผู้ใช้ใหม่ (users module) |
| POST | `/api/users` | ❌ | สร้าง user |
| POST | `/api/login` / `/api/signin` | ❌ | login (users module variant) |
| GET | `/api/user/me` | ✅ access | ข้อมูลผู้ใช้ตัวเอง |

---

### 3.2 `POST /api/auth/login`

รับ **JSON หรือ form-data** — ใช้ได้ทั้ง username และ email ใน field `username` (repo fallback: email ก่อน แล้ว username)

**Request:**
```json
{ "username": "admin", "password": "yourpassword" }
```
หรือ form: `-F "username=admin" -F "password=yourpassword"`

**Response 200:** (raw map — ไม่ wrap ใน `data`)
```json
{
  "access_token": "<jwt>",
  "refresh_token": "<jwt>",
  "token_type": "bearer",
  "expires_in": 216000
}
```

- `expires_in` = `AccessTokenExpireDuration × 60` (วินาที)
- password ต้อง min 8 ตัว
- มี rate limit 2 ชั้น: handler-internal 50/s burst 100 + route 10/s burst 20

### 3.3 `POST /api/auth/signin`

เหมือน login แต่ใช้ field `email`:
```json
{ "email": "admin@example.com", "password": "yourpassword" }
```
Response เหมือน login

---

### 3.4 `GET /api/auth/refresh` (ต้องใช้ refresh token)

ใช้ `Authorization: Bearer <refresh_token>` → ตรวจ `SIsMember` ใน Redis → **ลบ token เก่า (rotation)** → สร้างคู่ใหม่ → บันทึก token ใหม่

**Response 200:**
```json
{
  "access_token": "<new>",
  "refresh_token": "<new>",
  "token_type": "bearer"
}
```
> ⚠️ ไม่มี `expires_in` ใน response นี้

### 3.5 `GET /api/auth/logout` / `GET /api/auth/logoutall`

- `logout`: `Srem` refresh token ตัวเดียวที่ส่งมา → **HTTP 200 ว่าง**
- `logoutall`: `Delete` ทั้ง key `RefreshToken:{userID}` → **HTTP 200 ว่าง**

### 3.6 `GET /api/auth/verifyemail?code=<code>`

ยืนยันอีเมลด้วย code (32 hex chars จาก `secureRandom.RandomHex(16)`)
```json
{ "data": "Email verified successfully", "is_success": true }
```
- ถ้า verified ไปแล้ว → 401 `user_already_verified`

### 3.7 `POST /api/auth/forgotpassword`

```json
{ "email": "user@example.com" }
```
- ต้อง verified ก่อน (ไม่งั้น 401 `user_not_verified`)
- สร้าง reset token (16-byte hex) + `password_reset_at = now + 15 นาที`
- **Enqueue email ไป worker** (asynq `task:send_email`, queue `critical`, MaxRetry 10, delay 10s) พร้อมลิงก์ `http://localhost:5000/auth/resetpassword?code=<token>`

**Response:**
```json
{ "data": "You will receive a reset email if user with that email exist", "is_success": true }
```

### 3.8 `PATCH /api/auth/resetpassword?code=<token>`

```json
{ "new_password": "newpass123", "confirm_password": "newpass123" }
```
- เช็ค password ตรงกัน + token ยังไม่หมดอายุ (`GetByResetTokenResetAt`)
- bcrypt hash ใหม่ → ล้าง token → ล้าง cache + refresh token set
- Response: `{"data":"Password data updated successfully, please re-login","is_success":true}`

---

## 4. JWT Claims & Middleware

### 4.1 Claims (`pkg/jwt/token.go:14-18`)
```go
type AuthClaims struct {
    jwt.RegisteredClaims
    Email string `json:"email"`
    Id    string `json:"id"`
}
```
RegisteredClaims: `iat`, `iss` (Issuer), `exp`, `sub` (= user id), `nbf`

### 4.2 Middleware chain (`internal/middleware/jwtauth.go`)

| Middleware | หน้าที่ |
|------------|---------|
| `Verifier(requireAccessToken bool)` | อ่าน `Authorization: Bearer <token>`; `true`→ใช้ AccessTokenPublicKey, `false`→ใช้ RefreshTokenPublicKey; เก็บ token/id/email/error ลง context (**ยังไม่ reject**) |
| `Authenticator()` | reject 401 ถ้ามี error ใน context |
| `CurrentUser()` | parse id → uuid → โหลด `SdUser` เต็มจาก DB ลง context |
| `ActiveUser()` | reject 403 `inactive_user` ถ้า `Status != 1` |
| `SuperUser()` | reject 403 `not_enough_privileges` ถ้าไม่ใช่ superuser |

Chain มาตรฐานสำหรับ protected routes: `Verifier(true) → Authenticator → CurrentUser → ActiveUser` (และ `SuperUser` สำหรับ admin)

---

## 5. Token Lifecycle (Redis)

| Redis key | ประเภท | ใช้ที่ |
|-----------|--------|-------|
| `RefreshToken:{userID}` | Set (refresh tokens) | signin `Sadd` / refresh `SIsMember`+`Srem`+`Sadd` / logout `Srem` / logoutall `Delete` |
| `sd_user:{id}` | String JSON, TTL 3600s | user cache |

---

## 6. รูปแบบ Response

- **Success:** `{"data": ..., "is_success": true}`
- **Error:** `{"data": null, "error": {"status": <int>, "statusText": "<text>", "msg": "<msg>"}, "is_success": false}`

### Error statuses (`pkg/httpErrors/httpErrors.go`)
| statusText | HTTP |
|------------|------|
| `wrong_password` | 401 |
| `token_not_found` / `invalid_jwt_token` / `invalid_jwt_claims` | 401 |
| `not_found_refresh_token_redis` | 401 |
| `user_already_verified` / `user_not_verified` | 401 |
| `inactive_user` | 403 |
| `not_enough_privileges` | 403 |
| `not_found` | 404 |
| `validation` | 422 |

---

## 7. โครงสร้างฐานข้อมูล

### `sd_user` (auth columns)
| column | type |
|--------|------|
| `email` | text NOT NULL |
| `username` | text NOT NULL |
| `password` | text NOT NULL (bcrypt hash) |
| `verified` | bool default false |
| `verification_code` | varchar(64) |
| `password_reset_token` | varchar(64) |
| `password_reset_at` | timestamptz |
| `is_superuser` | bool default false |
| `status` | int2 (1=active, 0=inactive) |

### `refresh_tokens` (PG table — **มีแต่ไม่ได้ใช้จริง**)
มีตาราง `refresh_tokens` ใน migration (`20260712_new_modules_schema.sql:13`) แต่ระบบจริงเก็บ refresh token ใน **Redis** ไม่ใช่ตารางนี้ (`authPgRepo` เป็น dead code)

---

## 8. Edge Cases & ข้อควรระวัง

1. **Refresh token เก็บใน Redis ไม่ใช่ PostgreSQL** — ถ้า Redis flush จะถูกบังคับ logout ทุกคน (token ใน Redis หาย)
2. **`expires_in` เป็นวินาที** แต่ config duration เป็น **นาที**
3. **Refresh token rotation** — token เก่าใช้ไม่ได้ทันทีหลัง refresh (ถูก `Srem` ออก)
4. **Access token ไม่มี blacklist** — logout revoke เฉพาะ refresh token; access token ที่ยังไม่หมดอายุใช้ได้ต่อ
5. **ลิงก์ verify/reset hardcode `http://localhost:5000`** — ไม่ได้ใช้ `cfg.Server.BaseUrl` (ใช้ไม่ได้ถ้า deploy ต่าง domain)
6. **`/api/register`, `/api/users`, `/api/login`, `/api/signin` เป็น public** (users module)
7. **reset token หมดอายุ 15 นาที**
8. **password hash ใช้ bcrypt (DefaultCost)**
9. **Rate limit login** — 429 พร้อม header `X-RateLimit-Limit`, `Retry-After: 60` และ message ไทย `"ส่งคำขอถี่เกินไป กรุณาลองใหม่ภายหลัง"`
10. **Global security headers** (`internal/middleware/security.go`): `X-Frame-Options: DENY`, `X-Content-Type-Options: nosniff`, CSP, `Cache-Control: no-store` สำหรับ `/api/`

---

## 9. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| 401 `token_not_found` | ไม่ส่ง Bearer header |
| 401 `invalid_jwt_token` | token หมดอายุ/ถูกแก้/key ไม่ตรง (เช็ค `jwt.*Key` ใน config) |
| 401 `not_found_refresh_token_redis` | refresh token ถูกรีโวค/rotation ไปแล้ว หรือ Redis ถูก flush |
| 403 `inactive_user` | user `status != 1` |
| 403 `not_enough_privileges` | ต้องเป็น superuser |
| login 429 | ถี่เกินไป — รอ `Retry-After` |
| `expires_in` ดูแปลก (มาก) | เป็นวินาที = นาที×60 |
| verify/reset email ไม่ส่ง | worker ไม่รัน หรือ SMTP config ผิด (ดู README_Worker.md) |
