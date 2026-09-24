# Access Management Module 

Module ที่ 3 : การจัดการสิทธิ์การเข้าถึง (Access Management Module)

- `internal/modules/auth`
- `docs/modules/3-Access_Management_Module.md`

> **เป้าหมาย:** ให้ระบบจัดการสิทธิ์การเข้าถึง (Authentication / Authorization) ครบวงจร ได้แก่ การเข้าสู่ระบบด้วยชื่อผู้ใช้หรืออีเมล การออกและต่ออายุ Token แบบ Access Token + Refresh Token การออกจากระบบ การตรวจสอบสิทธิ์ผ่าน Middleware การยืนยันอีเมล การขอรีเซ็ตรหัสผ่านและการรีเซ็ตรหัสผ่าน
>
> **เหมาะสำหรับ:** ทีมพัฒนาที่ต้องการทำความเข้าใจสถาปัตยกรรม ฟีเจอร์ และระบบไหลของข้อมูล (Data Flow) ของโมดูลการเข้าถึงระบบบนพื้นฐานโครงสร้างแบบ Layered Architecture (Deliverable → Use Case → Repository → Model)

---

## สารบัญ

- [1. ภาพรวมระบบ](#1-ภาพรวมระบบ)
- [2. โครงสร้าง Module](#2-โครงสร้าง-module)
- [3. DOMAIN LAYER](#3-domain-layer)
- [4. APPLICATION LAYER](#4-application-layer)
- [5. INFRASTRUCTURE LAYER](#5-infrastructure-layer)
- [6. INTERFACE LAYER](#6-interface-layer)
- [7. Database Migrations](#7-database-migrations)
- [8. System Flow](#8-system-flow)
- [9. ภาคผนวก](#9-ภาคผนวก)

---

## 1. ภาพรวมระบบ

โมดูลนี้กระจุกอยู่ที่ `internal/modules/auth` โดยทำหน้าที่เป็น "ศูนย์กลางการพิสูจน์ตัวตน" ของทั้งระบบ ประกอบด้วย 2 กลไกหลัก คือ

1. **SignIn (เข้าสู่ระบบด้วยอีเมล + รหัสผ่าน)** — นักพัฒนาใช้งานผ่าน `authHandler` และ `authUseCase`
2. **Login (เข้าสู่ระบบด้วยชื่อผู้ใช้ + รหัสผ่าน)** — นักพัฒนาใช้งานผ่าน `authHandler` และ `authUseCase` โดยเรียกผ่าน `usersUC.SignInByUsername`

เมื่อพิสูจน์ตัวตนสำเร็จ ระบบจะออก **Access Token** และ **Refresh Token** (รูปแบบ JWT) พร้อมทั้งเซ็ต **Cookie** ประเภท `HttpOnly` + `SameSite=Lax` ให้เบราว์เซอร์ เพื่อให้ล็อกอินครั้งต่อไปส่ง Token ไปยังเซิร์ฟเวอร์ได้อย่างปลอดภัย

สำหรับเส้นทาง (Route) สาธารณะ เช่น การเข้าสู่ระบบ การยืนยันอีเมล และการรีเซ็ตรหัสผ่าน ระบบจะใช้ **Rate Limiting** ในระดับ Configurable (`MapAuthRoute`) เพื่อจำกัดจำนวนคำขอต่อวินาที ส่วนเส้นทางที่ต้องมีสิทธิ์ เช่น การรีเฟรช Token และการออกจากระบบ จะถูกห่อด้วย Middleware `Verifier` + `Authenticator`

✅ **ฟีเจอร์หลัก**

- เข้าสู่ระบบด้วยชื่อผู้ใช้และรหัสผ่าน (`POST /auth/login`)
- เข้าสู่ระบบด้วยอีเมลและรหัสผ่าน (`POST /auth/signin`)
- ออกรหัสสาธารณะ (Public Key) สำหรับตรวจสอบ JWT (`GET /auth/publickey`)
- ต่ออายุ Token (`GET /auth/refresh`) — ต้องใช้ Bearer Token
- ออกจากระบบ (`GET /auth/logout`) — ต้องใช้ Bearer Token
- ออกจากระบบทุกอุปกรณ์ (`GET` และ `POST /auth/logoutall`) — ต้องใช้ Bearer Token
- ยืนยันอีเมลด้วยรหัสยืนยัน (`GET /auth/verifyemail`)
- ขอรีเซ็ตรหัสผ่าน (`POST /auth/forgotpassword`)
- รีเซ็ตรหัสผ่าน (`PATCH /auth/resetpassword`)

---

## 2. โครงสร้าง Module

โครงสร้างจริงในโค้ด (`internal/modules/auth`) มีการแบ่งเลเยอร์ดังนี้

```
internal/modules/auth
├── Auth_User_Management_Module.md   // เอกสารกำกับโมดูลเดิม (อ้างอิงการออกแบบ)
├── delivery/http
│   ├── handler.go                   // การจับคู่ HTTP Handler → Use Case
│   ├── handlers.go                  // ตัวจัดการ HTTP (Login, SignIn, Refresh, ...)
│   ├── routes.go                    // ตารางเส้นทาง /auth
│   └── router.go                    // การลงทะเบียน Router หลัก + ค่า Rate Limit
├── model                           // Entity สำหรับ GORM (เช่น RefreshToken)
├── presenter/presenters.go         // DTO สำหรับรับ-ส่งข้อมูลผ่าน HTTP
├── repository
│   └── pg_repository.go            // Data Access ด้วย GORM
└── usecase
    └── usecase.go                  // Business Logic / Orchestration
```

> หมายเหตุ: เอกสารกำกับโมดูลเดิม (`Auth_User_Management_Module.md`) แสดงโครงสร้าง `domain/...` ที่เป็นเป้าหมายการออกแบบ ซึ่งอาจไม่ตรงกับโครงสร้างจริงในโค้ด เอกสารฉบับนี้ถือตามโครงสร้างจริงในโค้ดเป็นหลัก

---

## 3. DOMAIN LAYER

### 3.1 Entities

- **`RefreshToken`** (Model สำหรับ GORM ในโฟลเดอร์ `model`) — ใช้จัดเก็บ Refresh Token แต่ละชุดที่ระบบออกให้ผู้ใช้

| Field      | Type          | หมายเหตุ                                               |
|------------|---------------|--------------------------------------------------------|
| `ID`       | UUID (PK)     | สร้างจาก `gen_random_uuid()`                            |
| `UserID`   | UUID (FK)     | อ้างอิงตาราง `user` พร้อม `ON DELETE CASCADE`           |
| `Token`    | VARCHAR(512)  | ข้อความ Token                                           |
| `ExpiresAt`| TIMESTAMP     | เวลาหมดอายุ                                             |
| `CreatedAt`| TIMESTAMP     | เวลาสร้าง                                               |
| `Revoked`  | BOOLEAN       | สถานะเพิกถอน Token                                      |

### 3.2 Value Objects

- **ชนิดข้อมูลโครงสร้างทางธุรกิจ**: Access Token + Refresh Token ถูกมองเป็น "คู่ของข้อมูลประจำตัว" (`TokenResponse`) ที่ประกอบด้วย access_token, refresh_token, token_type และ expires_in
- ค่าคงที่สำหรับนิยาม Token — ระบบใช้ Public Key จาก `cfg.Jwt.AccessTokenPublicKey` และ `cfg.Jwt.RefreshTokenPublicKey` สำหรับตรวจสอบ/ถอดรหัส JWT

### 3.3 Repository Interfaces

- Interface `AuthRepositoryI` (ที่ root `internal/modules/auth`) กำหนดสัญญาณ Data Access สำหรับ auth module โดยมีเมทอดหลัก 4 รายการ
  - `StoreRefreshToken(ctx, req)` — บันทึก Refresh Token ใหม่
  - `FindRefreshToken(ctx, token)` — ค้นหา Refresh Token ด้วยข้อความ Token
  - `DeleteRefreshToken(ctx, token)` — ลบ Refresh Token
  - `DeleteAllUserRefreshTokens(ctx, userID, token)` — ลบ Refresh Token ทั้งหมดของ User ยกเว้น Token ปัจจุบัน

### 3.4 Domain Services

- ไม่มี Domain Service ภายใน auth module โดยตรง การดำเนินการทางธุรกิจหลัก (ตรวจสอบข้อมูลประจำตัว, ออก Token, ตรวจสอบสิทธิ์) ถูกจัดวางไว้ที่เลเยอร์ Use Case ผ่านการเรียก `usersUC` (จากโมดูลผู้ใช้) ร่วมกับ `authUseCase` เอง

### 3.5 Domain Errors

- โมดูลนี้ใช้นิยาม Error Code รหัสประจำโมดูล: `07_0x02` และ `07_0x03` เพื่อให้ระบุต้นตอของข้อผิดพลาดได้ชัดเจน เช่น การตรวจสอบข้อมูลประจำตัวล้มเหลว หรือการยืนยันข้อมูลไม่ถูกต้อง (ดูตาราง Error Code ในภาคผนวก)

---

## 4. APPLICATION LAYER

### 4.1 Use Cases

- **`authUseCase`** — โครงสร้างที่ประกอบด้วย `usersUC users.UserUseCaseI`, `cfg`, และ `logger` สร้างได้จาก `CreateAuthUseCaseI(usersUC, cfg, logger)`
- ตารางการแมปคำสั่งจากชั้น Delivery ไปยัง Use Case:

| HTTP / ฟีเจอร์          | เมทอดที่ถูกเรียก                          |
|--------------------------|------------------------------------------|
| `SignIn`                 | `usersUC.SignIn(credentials)`            |
| `Login`                  | `usersUC.SignInByUsername(username, password)` |
| `Refresh`                | `usersUC.Refresh(...)`                  |
| `Logout`                 | `usersUC.Logout(...)`                   |
| `LogoutAll`              | `usersUC.ParseIdFromRefreshToken(...)` แล้วเรียก `LogoutAll(...)` |
| `VerifyEmail`            | `usersUC.Verify(code)`                  |
| `ForgotPassword`         | `usersUC.ForgotPassword(email)`         |
| `ResetPassword`          | `usersUC.ResetPassword(resetToken, newPassword, confirmPassword)` |

### 4.2 DTOs

- **`TokenResponse`** — โครงสร้างตอบกลับสำหรับ Token

| Field          | Type   | หมายเหตุ                    |
|----------------|--------|-----------------------------|
| `access_token` | string | Access Token (JWT)          |
| `refresh_token`| string | Refresh Token (JWT)         |
| `token_type`   | string | ค่า: `bearer`                |
| `expires_in`   | int64  | จำนวนวินาทีก่อนหมดอายุ       |

- **`SignInRequest`** — ข้อมูลสำหรับเข้าสู่ระบบด้วยอีเมล

| Field    | Type   | Validation          |
|----------|--------|---------------------|
| `Email`  | string | `required,email`    |
| `Password`| string| `required,min=8`    |

- **`LoginRequest`** — ข้อมูลสำหรับเข้าสู่ระบบด้วยชื่อผู้ใช้

| Field      | Type   | Validation          |
|------------|--------|---------------------|
| `Username` | string | `required`          |
| `Password` | string | `required,min=8`    |

- **`RefreshTokenRequest`** — ข้อมูลสำหรับต่ออายุ Token

| Field          | Type   | Validation      |
|----------------|--------|-----------------|
| `RefreshToken` | string | `required`      |

- **`ForgotPasswordRequest`** — ข้อมูลสำหรับขอรีเซ็ตรหัสผ่าน

| Field | Type   | Validation       |
|-------|--------|------------------|
| `Email`| string| `required,email` |

- **`ResetPasswordRequest`** — ข้อมูลสำหรับรีเซ็ตรหัสผ่าน

| Field            | Type   | Validation                  |
|------------------|--------|-----------------------------|
| `NewPassword`    | string | `required,min=8` (json `new_password`)    |
| `ConfirmPassword`| string | `required,min=8` (json `confirm_password`)|

- **`UserSignIn`** — Struct สำหรับผูกข้อมูล Body บน Handler `SignIn` (อีเมล + รหัสผ่าน) ตาม Swagger

---

## 5. INFRASTRUCTURE LAYER

- **`authPgRepo`** — Repository ตัวจริงสำหรับการเข้าถึงฐานข้อมูลด้วย GORM ประกอบด้วย `db *gorm.DB` สร้างได้จาก `CreateAuthPgRepository(db)`
- การทำงานร่วมกับตารางผ่าน Model ของ GORM:
  - `StoreRefreshToken` — สร้างระเบียน `models.RefreshToken{UserID, Token, ExpiresAt}` ลงตาราง `refresh_tokens`
  - `FindRefreshToken` — ค้นหาโดย `WHERE token = ?`
  - `DeleteRefreshToken` — ลบ Refresh Token
  - `DeleteAllUserRefreshTokens` — ลบ Refresh Token ทั้งหมดของ User (ยกเว้น Token ปัจจุบัน) เพื่อรองรับฟีเจอร์ Logout All Devices
- ตัวชี้วัด (Router/Rate Limit) อยู่ในเลเยอร์ Delivery (`router.go`) โดยใช้ค่าคอนฟิกจาก `MapAuthRoute`

---

## 6. INTERFACE LAYER

### 6.1 HTTP Handlers

Handler ทั้งหมดอยู่ที่ `delivery/http/handlers.go` (ตรวจสอบครบจนถึงบรรทัดที่ 589) รายละเอียดดังนี้

| Handler            | HTTP Method / Path      | ลักษณะการทำงาน                                      |
|--------------------|-------------------------|-----------------------------------------------------|
| `Login`            | POST `/auth/login`      | ตรวจสอบชื่อผู้ใช้ + รหัสผ่าน ผ่าน `usersUC.SignInByUsername` ตอบกลับ JSON `{access_token, refresh_token, token_type:"bearer", expires_in}` พร้อม Cookie `HttpOnly` + `SameSite=Lax` |
| `SignIn`           | POST `/auth/signin`     | ตรวจสอบอีเมล + รหัสผ่าน ผ่าน `usersUC.SignIn` โดยผูก Body ด้วย `presenter.UserSignIn` |
| `RefreshToken`     | GET `/auth/refresh`     | ต่ออายุ Token จาก Refresh Token ที่ได้ ใช้ `@Security BearerAuth` ตอบกลับ `presenter.Token` |
| `GetPublicKey`     | GET `/auth/publickey`   | ถอดรหัส Base64 จาก `cfg.Jwt.AccessTokenPublicKey` และ `cfg.Jwt.RefreshTokenPublicKey` แล้วตอบกลับ |
| `VerifyEmail`      | GET `/auth/verifyemail`| ใช้ค่า query `code` เรียก `usersUC.Verify(code)` เพื่อยืนยันอีเมล |
| `ForgotPassword`   | POST `/auth/forgotpassword` | เรียก `usersUC.ForgotPassword(email)` ตอบกลับข้อความ "You will receive a reset email if user with that email exist" |
| `ResetPassword`    | PATCH `/auth/resetpassword`  | ผูก Body `presenter.ResetPassword` + อ่าน query `code` เรียก `usersUC.ResetPassword(...)` ตอบกลับ "Password data updated successfully, please re-login" |
| `Logout`           | GET `/auth/logout`      | ยกเลิก Refresh Token ปัจจุบัน (ต้องใช้ Bearer Token) |
| `LogoutAllToken`   | GET+POST `/auth/logoutall`  | ยกเลิก Refresh Token ทุกชุดของ User (ต้องใช้ Bearer Token) |

### 6.2 Routes

จาก `delivery/http/routes.go` (ตรวจสอบ L1–52 และยืนยันรายละเอียดก่อนหน้า) เส้นทาง `/auth` แบ่งเป็น 2 กลุ่ม

**กลุ่ม Public — ใช้ Rate Limiting ตาม `MapAuthRoute` (10 req/s, burst 20, cleanup 15 นาที)**

| Method | Path                  | Handler        |
|--------|-----------------------|----------------| 
| POST   | `/auth/login`         | `handler.Login`    |
| POST   | `/auth/signin`        | `handler.SignIn`   |
| GET    | `/auth/publickey`     | `handler.GetPublicKey` |
| GET    | `/auth/verifyemail`   | `handler.VerifyEmail` |
| POST   | `/auth/forgotpassword`| `handler.ForgotPassword` |
| PATCH  | `/auth/resetpassword` | `handler.ResetPassword` |

**กลุ่ม Protected — ใช้ `mw.Verifier(false)` + `mw.Authenticator()`**

| Method | Path                  | Handler        |
|--------|-----------------------|----------------|
| GET    | `/auth/refresh`       | `handler.RefreshToken` |
| GET    | `/auth/logout`        | `handler.Logout` |
| GET    | `/auth/logoutall`     | `handler.LogoutAllToken` |
| POST   | `/auth/logoutall`     | `handler.LogoutAllToken` |

นอกจากนี้ `authHandler.Routes` ยังมีเส้นทางแยกสำหรับ `POST /login` และ `POST /signin` ด้วยค่า Rate Limit 50 req/s, burst 100 โดยที่ `signup`, `logout`, `refresh` ถูก Comment ไว้

> ข้อสังเกต: ระบบไม่มีเส้นทาง `/auth/register` และ `/auth/verify` แต่มี `/register` อยู่ระดับบนสุด

### 6.3 Middleware

- **Rate Limiter** — ใช้กับกลุ่ม Public เพื่อป้องกันการยิงคำขอจำนวนมาก (10 req/s burst 20 สำหรับกลุ่ม `/auth` หลัก และ 50 req/s burst 100 สำหรับ `POST /login`, `POST /signin` ผ่าน `authHandler.Routes`)
- **`mw.Verifier(false)`** — ตรวจสอบโครงสร้าง/ลายเซ็นของ Token
- **`mw.Authenticator()`** — ยืนยันตัวตนผู้ใช้จาก Token ก่อนเข้าถึงเส้นทางที่ต้องมีสิทธิ์

---

## 7. Database Migrations

| ไฟล์                                 | เนื้อหาโดยสังเขป                                                              |
|--------------------------------------|-------------------------------------------------------------------------------|
| `migrations/20260712_new_modules_schema.sql` | สร้างตาราง `refresh_tokens`: `id` UUID PK `gen_random_uuid()`, `user_id` UUID NOT NULL REFERENCES `"user"(id)` ON DELETE CASCADE, `token` VARCHAR(512) NOT NULL, `expires_at` TIMESTAMP NOT NULL, `created_at` TIMESTAMP NOT NULL DEFAULT NOW(), `revoked` BOOLEAN NOT NULL DEFAULT FALSE พร้อม Index `idx_refresh_tokens_token` และ `idx_refresh_tokens_user` |
| `migrations/sd_user_access_menu.sql` | Navicat dump — ตาราง `public.sd_user_access_menu`: `user_access_id` int8 PK (Default nextval), `user_type_id` int8, `menu_id` int8, `parent_id` int8 (ยังไม่มีข้อมูล) |
| `migrations/20260712_seed_data.sql`  | Seed ข้อมูล `email_config` (เช่น smtp.gmail.com) สำหรับส่งอีเมล และ `m_translation` (i18n ไทย/อังกฤษ) |
| `internal/template/migration.sql`    | แม่แบบข้อตกลงการสร้าง Migration — ใช้ `pgcrypto` extension, `id` UUID PK `gen_random_uuid()`, `created_at`/`updated_at` TIMESTAMPTZ DEFAULT NOW(), `deleted_at` TIMESTAMPTZ, `status` INTEGER DEFAULT 1 |

---

## 8. System Flow

### 8.1 เข้าสู่ระบบ (Login / SignIn)

1. Client ส่งคำขอ `POST /auth/login` (username) หรือ `POST /auth/signin` (email) พร้อมรหัสผ่าน
2. Rate Limiter ตรวจสอบจำนวนคำขอ
3. `handler.Login` / `handler.SignIn` เรียก Use Case → `usersUC.SignInByUsername` / `usersUC.SignIn`
4. เมื่อข้อมูลถูกต้อง ระบบออก Access Token + Refresh Token
5. ตอบกลับ JSON `{access_token, refresh_token, token_type, expires_in}` พร้อมเซ็ต Cookie `HttpOnly` + `SameSite=Lax`
6. `StoreRefreshToken` บันทึก Refresh Token ลงตาราง `refresh_tokens`

### 8.2 ต่ออายุ Token (Refresh)

1. Client ส่ง `GET /auth/refresh` พร้อม Bearer Token
2. ผ่าน `mw.Verifier(false)` + `mw.Authenticator()`
3. `handler.RefreshToken` เรียก `usersUC.Refresh(...)` เพื่อตรวจสอบและออก Token ชุดใหม่
4. ตอบกลับ `presenter.Token`

### 8.3 ออกจากระบบ (Logout / Logout All)

1. Client ส่ง `GET /auth/logout` หรือ `GET/POST /auth/logoutall` พร้อม Bearer Token
2. ผ่าน Middleware ตรวจสอบสิทธิ์
3. `handler.Logout` ลบ Refresh Token ปัจจุบันด้วย `DeleteRefreshToken`
4. `handler.LogoutAllToken` ระบุ User จาก RefreshToken `ParseIdFromRefreshToken` แล้วลบ Token ทั้งหมดด้วย `DeleteAllUserRefreshTokens`

### 8.4 ขอรีเซ็ตรหัสผ่าน (Forgot Password)

1. Client ส่ง `POST /auth/forgotpassword` พร้อมอีเมล
2. `handler.ForgotPassword` เรียก `usersUC.ForgotPassword(email)`
3. ระบบตอบกลับข้อความกลาง "You will receive a reset email if user with that email exist" (ไม่รั่วไหลข้อมูลว่าอีเมลมีอยู่ในระบบหรือไม่)

### 8.5 รีเซ็ตรหัสผ่าน (Reset Password)

1. Client ส่ง `PATCH /auth/resetpassword` พร้อม Body `{new_password, confirm_password}` และ query `code`
2. `handler.ResetPassword` เรียก `usersUC.ResetPassword(resetToken, new, confirm)`
3. เมื่อสำเร็จ ตอบกลับ "Password data updated successfully, please re-login"

### 8.6 ยืนยันอีเมล (Verify Email)

1. Client ส่ง `GET /auth/verifyemail?code=...`
2. `handler.VerifyEmail` เรียก `usersUC.Verify(code)`
3. เมื่อสำเร็จ อีเมลของผู้ใช้ถูกยืนยัน

---

## 9. ภาคผนวก

### 9.1 Error Codes ของโมดูล

| Code      | หมวดหมู่        | หมายเหตุ                                         |
|-----------|----------------|--------------------------------------------------|
| `07_0x02` | Authentication | ข้อผิดพลาดที่เกี่ยวกับการตรวจสอบสิทธิ์การเข้าถึง                      |
| `07_0x03` | Authentication | ข้อผิดพลาดที่เกี่ยวกับการตรวจสอบสิทธิ์การเข้าถึง (กลุ่มรอง) |

### 9.2 ตาราง Endpoint ↔ ฟังก์ชัน (อ้างอิง Swagger)

| Path                  | Method | Summary ตาม Swagger (โดยสังเขป)            | Body ตาม Swagger      |
|-----------------------|--------|------------------------------------------|-----------------------|
| `/auth/login`         | POST   | "User login with username and password" / 200 "Tokens returned successfully" | `presenter.SignIn` (username-based) |
| `/auth/signin`        | POST   | เข้าสู่ระบบด้วย email + password          | `presenter.UserSignIn`|
| `/auth/publickey`     | GET    | ส่ง Public Key สำหรับตรวจสอบ JWT          | -                     |
| `/auth/refresh`       | GET    | Refresh Token (`@Security BearerAuth`)    | -                     |
| `/auth/logout`        | GET    | ออกจากระบบ                                | -                     |
| `/auth/logoutall`     | GET/POST | ออกจากระบบทุกอุปกรณ์                     | -                     |
| `/auth/verifyemail`   | GET    | ยืนยันอีเมลด้วย `code`                     | -                     |
| `/auth/forgotpassword`| POST   | ขอรีเซ็ตรหัสผ่านด้วยอีเมล                  | -                     |
| `/auth/resetpassword` | PATCH  | รีเซ็ตรหัสผ่านด้วย `code` + `new_password`/`confirm_password` | `presenter.ResetPassword` |

### 9.3 ข้อสังเกตความไม่ตรงกันระหว่าง Swagger และ Handler (Discrepancy)

- Swagger `/auth/login` ระบุ Body เป็น `presenter.SignIn` (แบบ **username-based**) ขณะที่ description ของ `/auth/signin` ระบุ email+password และ Body เป็น `presenter.UserSignIn`
- ในระดับ Handler โค้ดจริงผูกไว้ชัดเจน: `Login` → `usersUC.SignInByUsername` และ `SignIn` → `usersUC.SignIn` ซึ่งสอดคล้องกับรูปแบบ "username login" และ "email login" ตามลำดับ
- จึงควรยึดตาม Handler เป็นหลัก และใช้ Swagger เป็นเอกสารประกอบการเรียกใช้ API เท่านั้น

### 9.4 เอกสารอ้างอิง

- `internal/modules/auth/delivery/http/handlers.go`
- `internal/modules/auth/delivery/http/routes.go`
- `internal/modules/auth/delivery/http/router.go`
- `internal/modules/auth/delivery/http/handler.go`
- `internal/modules/auth/presenter/presenters.go`
- `internal/modules/auth/usecase/usecase.go`
- `internal/modules/auth/repository/pg_repository.go`
- `internal/modules/auth/model`
- `migrations/20260712_new_modules_schema.sql`
- `migrations/sd_user_access_menu.sql`
- `migrations/20260712_seed_data.sql`
- `docs/swagger/swagger.yaml` (ส่วน `/auth` + definitions)