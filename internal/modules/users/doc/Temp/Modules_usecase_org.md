# เอกสารประกอบโมดูล User

## สารบัญ
1. [โครงสร้างโฟลเดอร์](#1-โครงสร้างโฟลเดอร์-folder-structure)
2. [หลักการทำงาน (Concept)](#2-หลักการทำงาน-concept)
3. [Workflow และ Dataflow](#3-workflow-และ-dataflow)
4. [คู่มือการทดสอบ](#4-คู่มือการทดสอบ)
5. [คู่มือการใช้งาน](#5-คู่มือการใช้งาน)
6. [คู่มือการบำรุงรักษา](#6-คู่มือการบำรุงรักษา)
7. [คู่มือการขยาย/แก้ไข](#7-คู่มือการขยายหรือแก้ไข)
8. [Check List Test Module](#8-checklist-test-module)

---

## 1. โครงสร้างโฟลเดอร์ (Folder Structure)

```text
icmongolang/
│
├── cmd/                                    # Command line entry point
│   └── main.go
│
├── config/                                 # Configuration management
│   └── config.go
│
├── internal/                               # Private application code
│   │
│   ├── middleware/                         # HTTP middleware layer / ระดับกลาง HTTP
│   │   ├── jwtauth.go                     # JWT authentication logic / ตรรกะตรวจสอบสิทธิ์ JWT
│   │   └── middleware.go                  # Core middleware manager / ตัวจัดการมิดเดิลแวร์หลัก
│   │
│   ├── models/                             # Data models / โมเดลข้อมูล
│   │   └── sd_user.go                     # SdUser entity definition / นิยามเอนทิตีผู้ใช้
│   │
│   └── users/                              # User module root / โมดูลหลักของผู้ใช้
│       │
│       ├── handler.go                      # HTTP handler interface definition
│       │                                   # นิยามอินเทอร์เฟซสำหรับตัวจัดการ HTTP
│       │
│       ├── pg_repository.go                # PostgreSQL repository interface
│       │                                   # อินเทอร์เฟซรีโพสิทอรี PostgreSQL
│       │
│       ├── redis_repository.go             # Redis repository interface
│       │                                   # อินเทอร์เฟซรีโพสิทอรี Redis
│       │
│       ├── usecase.go                      # Use case interface definition
│       │                                   # นิยามอินเทอร์เฟซยูสเคส
│       │
│       ├── worker.go                       # Async task definitions / นิยามงาน异步
│       │
│       ├── delivery/                       # Delivery layer / ชั้นการนำส่งข้อมูล
│       │   └── http/
│       │       ├── handlers.go             # HTTP handlers implementation
│       │       │                           # การ implement ตัวจัดการ HTTP
│       │       └── routes.go               # Route registration / การลงทะเบียนเส้นทาง
│       │
│       ├── distributor/                    # Task distributor / ตัวกระจายงาน
│       │   └── distributor.go              # Redis task distribution logic
│       │                                   # ตรรกะการกระจายงานผ่าน Redis
│       │
│       ├── presenter/                      # Data presentation layer
│       │   └── presenters.go               # Request/Response DTOs
│       │                                   # DTO สำหรับรับ/ส่งข้อมูล
│       │
│       ├── processor/                      # Task processor / ตัวประมวลผลงาน
│       │   └── processor.go                # Redis task processing logic
│       │                                   # ตรรกะการประมวลผลงานผ่าน Redis
│       │
│       ├── repository/                     # Repository implementation
│       │   ├── pg_repository.go            # PostgreSQL implementation
│       │   │                               # การ implement PostgreSQL
│       │   └── redis_repository.go         # Redis implementation
│       │                                   # การ implement Redis
│       │
│       └── usecase/                        # Business logic layer
│           └── usecase.go                  # Use case implementation
│                                           # การ implement ยูสเคส
│
├── pkg/                                    # Public shared packages
│   ├── cryptpass/                          # Password hashing / การเข้ารหัสรหัสผ่าน
│   ├── emailTemplates/                     # Email template generator
│   ├── httpErrors/                         # HTTP error definitions
│   ├── jwt/                                # JWT management
│   ├── logger/                             # Logging utility
│   ├── responses/                          # API response formatter
│   ├── secureRandom/                       # Secure random generator
│   └── utils/                              # Utility functions
│
└── migrations/                             # Database migrations
    └── *.sql
```

---

## 2. หลักการทำงาน (Concept)

### คืออะไร? (What is it?)
โมดูล User เป็นระบบจัดการผู้ใช้ที่สมบูรณ์ (Complete User Management System) ที่รองรับการ:
- ลงทะเบียนและยืนยันตัวตน (Registration & Email Verification)
- เข้าสู่ระบบด้วย JWT (JWT-based Authentication)
- จัดการโปรไฟล์ (Profile Management)
- เปลี่ยนรหัสผ่าน (Password Change)
- ลืมรหัสผ่าน (Forgot/Reset Password)
- สิทธิ์ผู้ดูแลระบบ (Admin/Superuser Privileges)
- ระบบออกจากระบบทั้งหมด (Global Logout)

### มีกี่แบบ? (How many types?)
ระบบแบ่งออกเป็น **3 ระดับสิทธิ์ (3 Access Levels)**:

| ระดับ | ชื่อ | สิทธิ์ |
|------|------|--------|
| 1 | Super Admin | ทำได้ทุกอย่าง |
| 2 | Admin | จัดการผู้ใช้ได้ ยกเว้น Super Admin |
| 3 | User | จัดการเฉพาะโปรไฟล์ตัวเอง |

### ข้อห้ามสำคัญ (Critical Prohibitions)
1. **ห้ามเก็บรหัสผ่านแบบ plain text** - ต้องผ่านการ hash ด้วย bcrypt เสมอ
2. **ห้ามส่ง JWT token ผ่าน URL** - ใช้ Header `Authorization: Bearer <token>` เท่านั้น
3. **ห้าม hardcode secret key** - ต้องอ่านจาก environment variable
4. **ห้าม bypass middleware** - ทุก API ที่ต้องใช้สิทธิ์ต้องผ่าน middleware
5. **ห้ามเก็บ refresh token ใน local storage** - ความเสี่ยง XSS สูง

---

## 3. Workflow และ Dataflow

### 3.1 Workflow การลงทะเบียน (Registration Workflow)

```
[Client]                    [Handler]                   [UseCase]                  [Repository]              [Redis]              [Email Worker]
   │                           │                           │                           │                         │                      │
   │  POST /register          │                           │                           │                         │                      │
   │  {email, password, ...}  │                           │                           │                         │                      │
   │──────────────────────────>│                           │                           │                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │  Validate request        │                           │                         │                      │
   │                           │  (check password match)   │                           │                         │                      │
   │                           │──────────────────────────>│                           │                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │  Hash password            │                         │                      │
   │                           │                           │  (bcrypt)                 │                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │  Create user              │                         │                      │
   │                           │                           │──────────────────────────>│                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │                           │  INSERT INTO sd_users  │                      │
   │                           │                           │                           │  (status=1, verified=0)│                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │<──────────────────────────│                         │                      │
   │                           │                           │  user created            │                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │  Generate verification   │                         │                      │
   │                           │                           │  code (random hex 16)    │                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │  Update verification_code│                         │                      │
   │                           │                           │──────────────────────────>│                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │  Distribute email task   │                         │                      │
   │                           │                           │─────────────────────────────────────────────────────>│                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │                           │                         │  Enqueue task        │
   │                           │                           │                           │                         │  (Queue: critical)   │
   │                           │                           │                           │                         │                      │
   │                           │                           │                           │                         │──────────────────────>│
   │                           │                           │                           │                         │                      │
   │                           │<──────────────────────────│                           │                         │                      │
   │                           │  UserResponse (no token)  │                           │                         │                      │
   │                           │                           │                           │                         │                      │
   │<──────────────────────────│                           │                           │                         │                      │
   │  201 Created             │                           │                           │                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │                           │                         │  [Async] Send email  │
   │                           │                           │                           │                         │                      │
```

### 3.2 Workflow การยืนยันอีเมล (Email Verification Workflow)

```
[Client]                    [Handler]                   [UseCase]                  [Repository]              [Redis]
   │                           │                           │                           │                         │
   │  GET /verifyemail?code=xxx│                           │                           │                         │
   │──────────────────────────>│                           │                           │                         │
   │                           │                           │                           │                         │
   │                           │  Verify(verificationCode) │                           │                         │
   │                           │──────────────────────────>│                           │                         │
   │                           │                           │                           │                         │
   │                           │                           │  GetByVerificationCode()  │                         │
   │                           │                           │──────────────────────────>│                         │
   │                           │                           │                           │                         │
   │                           │                           │                           │  SELECT * FROM sd_users│
   │                           │                           │                           │  WHERE verification_code│
   │                           │                           │                           │                         │
   │                           │                           │<──────────────────────────│                         │
   │                           │                           │  user found              │                         │
   │                           │                           │                           │                         │
   │                           │                           │  Check if already verified│                         │
   │                           │                           │  -> if yes: return error │                         │
   │                           │                           │                           │                         │
   │                           │                           │  UpdateVerification()     │                         │
   │                           │                           │  (clear code, verified=1) │                         │
   │                           │                           │──────────────────────────>│                         │
   │                           │                           │                           │                         │
   │                           │                           │  Delete Redis cache       │                         │
   │                           │                           │─────────────────────────────────────────────────────>│
   │                           │                           │                           │                         │
   │                           │<──────────────────────────│                           │                         │
   │                           │  nil error               │                           │                         │
   │                           │                           │                           │                         │
   │<──────────────────────────│                           │                           │                         │
   │  Success (redirect or JSON)                           │                           │                         │
   │                           │                           │                           │                         │
```

### 3.3 Workflow การลืมรหัสผ่าน (Forgot/Reset Password Workflow)

```
[Client]                    [Handler]                   [UseCase]                  [Repository]              [Redis]              [Email Worker]
   │                           │                           │                           │                         │                      │
   │  POST /forgotpassword     │                           │                           │                         │                      │
   │  {email}                  │                           │                           │                         │                      │
   │──────────────────────────>│                           │                           │                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │  ForgotPassword(email)    │                           │                         │                      │
   │                           │──────────────────────────>│                           │                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │  GetByEmail()             │                         │                      │
   │                           │                           │──────────────────────────>│                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │<──────────────────────────│                         │                      │
   │                           │                           │  user found              │                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │  Check if verified        │                         │                      │
   │                           │                           │  -> if not: return error │                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │  Generate reset token     │                         │                      │
   │                           │                           │  (random hex 16)          │                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │  UpdatePasswordReset()    │                         │                      │
   │                           │                           │  (token, expires+15min)   │                         │                      │
   │                           │                           │──────────────────────────>│                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │  Delete Redis cache       │                         │                      │
   │                           │                           │─────────────────────────────────────────────────────>│                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │  Distribute reset email   │                         │                      │
   │                           │                           │─────────────────────────────────────────────────────────────────────────>│
   │                           │                           │                           │                         │                      │
   │                           │<──────────────────────────│                           │                         │                      │
   │                           │  nil error               │                           │                         │                      │
   │                           │                           │                           │                         │                      │
   │<──────────────────────────│                           │                           │                         │                      │
   │  200 OK (success)        │                           │                           │                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │                           │                         │                      │
   │  POST /resetpassword      │                           │                           │                         │                      │
   │  {token, new_password,    │                           │                           │                         │                      │
   │   confirm_password}       │                           │                           │                         │                      │
   │──────────────────────────>│                           │                           │                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │  ResetPassword(token,     │                           │                         │                      │
   │                           │    newPassword, confirm)  │                           │                         │                      │
   │                           │──────────────────────────>│                           │                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │  GetByResetTokenResetAt() │                         │                      │
   │                           │                           │  (check not expired)      │                         │                      │
   │                           │                           │──────────────────────────>│                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │<──────────────────────────│                         │                      │
   │                           │                           │  user found              │                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │  Hash new password        │                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │  UpdatePasswordResetToken()│                        │                      │
   │                           │                           │  (clear token, update pw) │                         │                      │
   │                           │                           │──────────────────────────>│                         │                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │  Delete Redis cache       │                         │                      │
   │                           │                           │─────────────────────────────────────────────────────>│                      │
   │                           │                           │                           │                         │                      │
   │                           │                           │  Delete refresh tokens    │                         │                      │
   │                           │                           │─────────────────────────────────────────────────────>│                      │
   │                           │                           │                           │                         │                      │
   │                           │<──────────────────────────│                           │                         │                      │
   │                           │  nil error               │                           │                         │                      │
   │                           │                           │                           │                         │                      │
   │<──────────────────────────│                           │                           │                         │                      │
   │  200 OK (success)        │                           │                           │                         │                      │
   │                           │                           │                           │                         │                      │
```

### 3.4 Workflow การเข้าสู่ระบบ (Login Workflow)

```
[Client]                    [Handler]                   [UseCase]                  [Repository]              [Redis]
   │                           │                           │                           │                         │
   │  POST /signin             │                           │                           │                         │
   │  {email, password}        │                           │                           │                         │
   │──────────────────────────>│                           │                           │                         │
   │                           │                           │                           │                         │
   │                           │  SignIn(email, password)  │                           │                         │
   │                           │──────────────────────────>│                           │                         │
   │                           │                           │                           │                         │
   │                           │                           │  GetByEmail()             │                         │
   │                           │                           │──────────────────────────>│                         │
   │                           │                           │                           │                         │
   │                           │                           │<──────────────────────────│                         │
   │                           │                           │  user found              │                         │
   │                           │                           │                           │                         │
   │                           │                           │  Compare password (bcrypt)│                         │
   │                           │                           │  -> if fail: return error│                         │
   │                           │                           │                           │                         │
   │                           │                           │  Create Access Token      │                         │
   │                           │                           │  (RS256, 15min expiry)    │                         │
   │                           │                           │                           │                         │
   │                           │                           │  Create Refresh Token     │                         │
   │                           │                           │  (RS256, 7d expiry)       │                         │
   │                           │                           │                           │                         │
   │                           │                           │  Save refresh token       │                         │
   │                           │                           │  to Redis Set             │                         │
   │                           │                           │─────────────────────────────────────────────────────>│
   │                           │                           │                           │                         │
   │                           │<──────────────────────────│                           │                         │
   │                           │  accessToken, refreshToken│                           │                         │
   │                           │                           │                           │                         │
   │<──────────────────────────│                           │                           │                         │
   │  {access_token, refresh_token}                         │                           │                         │
   │                           │                           │                           │                         │
```

### 3.5 Dataflow Diagram (ระดับข้อมูล)

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                                    DATA FLOW                                         │
└─────────────────────────────────────────────────────────────────────────────────────┘

    HTTP Request                 Handler                   UseCase
         │                           │                        │
         ▼                           ▼                        ▼
┌─────────────────┐         ┌─────────────────┐     ┌─────────────────┐
│  JSON Request   │         │   Validator     │     │  Business Logic │
│  (Presenter)    │ ──────► │  (utils)        │ ──► │  (Service)      │
└─────────────────┘         └─────────────────┘     └────────┬────────┘
                                                              │
                    ┌─────────────────────────────────────────┼─────────────────────────────────────────┐
                    │                                         │                                         │
                    ▼                                         ▼                                         ▼
         ┌─────────────────┐                       ┌─────────────────┐                       ┌─────────────────┐
         │  PostgreSQL     │                       │     Redis       │                       │   Email Queue   │
         │  (Primary DB)   │                       │   (Cache)       │                       │   (Asynq)       │
         ├─────────────────┤                       ├─────────────────┤                       ├─────────────────┤
         │ • sd_users      │                       │ • user:{id}     │                       │ • Task: email   │
         │ • CRUD ops      │                       │ • refresh:{id}  │                       │ • Retry logic   │
         │ • Indexes       │                       │ • TTL: 1 hour   │                       │ • Queue: crit   │
         └────────┬────────┘                       └────────┬────────┘                       └────────┬────────┘
                  │                                         │                                         │
                  ▼                                         ▼                                         ▼
         ┌─────────────────┐                       ┌─────────────────┐                       ┌─────────────────┐
         │  Model Mapping  │                       │  Cache Inval    │                       │  Email Sender   │
         │  (SdUser)       │                       │  on Update      │                       │  (SMTP)         │
         └─────────────────┘                       └─────────────────┘                       └─────────────────┘
```

---

## 4. คู่มือการทดสอบ (Testing Guide)

### 4.1 Environment Setup

```bash
# 1. Start dependencies
docker-compose up -d postgres redis

# 2. Run migrations
make migrate-up

# 3. Run tests
go test ./internal/modules/users/... -v

# 4. Run with coverage
go test ./internal/modules/users/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 4.2 Test Checklist

| # | Test Case | Method | Endpoint | Expected Result |
|---|-----------|--------|----------|-----------------|
| **Register API** |
| 1 | Register success | POST | /register | 200, user data returned |
| 2 | Register duplicate email | POST | /register | 400, error message |
| 3 | Register password mismatch | POST | /register | 400, validation error |
| 4 | Register invalid email | POST | /register | 400, validation error |
| **SignIn API** |
| 5 | SignIn success | POST | /signin | 200, tokens returned |
| 6 | SignIn wrong password | POST | /signin | 401, unauthorized |
| 7 | SignIn non-existent user | POST | /signin | 404, not found |
| **User Management (Admin)** |
| 8 | Create user (admin) | POST | /user | 200, user created |
| 9 | Get user by ID | GET | /user/{id} | 200, user data |
| 10 | Get user list | GET | /user?limit=10&offset=0 | 200, user array |
| 11 | Update user | PUT | /user/{id} | 200, updated user |
| 12 | Delete user | DELETE | /user/{id} | 200, deleted user |
| **Profile Management (Self)** |
| 13 | Get my profile | GET | /user/me | 200, user data |
| 14 | Update my profile | PUT | /user/me | 200, updated user |
| 15 | Update my password | PATCH | /user/me/updatepass | 200, success |
| **Email Verification** |
| 16 | Verify email with valid code | GET | /verifyemail?code=xxx | 200, success |
| 17 | Verify email with expired code | GET | /verifyemail?code=xxx | 400, error |
| 18 | Verify already verified | GET | /verifyemail?code=xxx | 400, already verified |
| **Forgot/Reset Password** |
| 19 | Forgot password valid email | POST | /forgotpassword | 200, success |
| 20 | Forgot password invalid email | POST | /forgotpassword | 404, not found |
| 21 | Reset password valid token | POST | /resetpassword | 200, success |
| 22 | Reset password expired token | POST | /resetpassword | 400, error |
| **Token Management** |
| 23 | Refresh token | POST | /refresh | 200, new tokens |
| 24 | Logout | POST | /logout | 200, success |
| 25 | Logout all (admin) | GET | /user/{id}/logoutall | 200, success |
| **Authorization** |
| 26 | Access admin API as normal user | GET | /user | 403, forbidden |
| 27 | Access without token | GET | /user/me | 401, unauthorized |

### 4.3 API Test Commands (cURL)

```bash
# Register
curl -X POST http://localhost:5000/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "confirm_password": "password123",
    "role_id": 2,
    "firstname": "Test",
    "lastname": "User"
  }'

# Sign In
curl -X POST http://localhost:5000/signin \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# Get My Profile (with token)
curl -X GET http://localhost:5000/user/me \
  -H "Authorization: Bearer <access_token>"

# Update My Profile
curl -X PUT http://localhost:5000/user/me \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "firstname": "UpdatedName",
    "phone_number": "0812345678"
  }'

# Change Password
curl -X PATCH http://localhost:5000/user/me/updatepass \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "password123",
    "new_password": "newpassword123",
    "confirm_password": "newpassword123"
  }'

# Forgot Password
curl -X POST http://localhost:5000/forgotpassword \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com"
  }'

# Reset Password
curl -X POST http://localhost:5000/resetpassword \
  -H "Content-Type: application/json" \
  -d '{
    "reset_token": "<token_from_email>",
    "new_password": "newpassword123",
    "confirm_password": "newpassword123"
  }'

# Refresh Token
curl -X POST http://localhost:5000/refresh \
  -H "Authorization: Bearer <refresh_token>"

# Logout
curl -X POST http://localhost:5000/logout \
  -H "Authorization: Bearer <refresh_token>"
```

---

## 5. คู่มือการใช้งาน (User Guide)

### 5.1 API Endpoints Summary

| Method | Endpoint | Auth Required | Role | Description |
|--------|----------|---------------|------|-------------|
| POST | /register | ❌ | Public | Register new user |
| POST | /signin | ❌ | Public | Login |
| POST | /refresh | ❌ | Public | Refresh access token |
| POST | /logout | ✅ | User | Logout |
| POST | /forgotpassword | ❌ | Public | Request password reset |
| POST | /resetpassword | ❌ | Public | Reset password with token |
| GET | /verifyemail | ❌ | Public | Verify email with code |
| GET | /user/me | ✅ | User | Get own profile |
| PUT | /user/me | ✅ | User | Update own profile |
| PATCH | /user/me/updatepass | ✅ | User | Change own password |
| POST | /user | ✅ | Super Admin | Create user |
| GET | /user | ✅ | Super Admin | List users |
| GET | /user/{id} | ✅ | Super Admin | Get user by ID |
| PUT | /user/{id} | ✅ | Super Admin | Update user |
| DELETE | /user/{id} | ✅ | Super Admin | Delete user |
| PATCH | /user/{id}/role | ✅ | Super Admin | Update user role |
| PATCH | /user/{id}/updatepass | ✅ | Super Admin | Update user password |
| GET | /user/{id}/logoutall | ✅ | Super Admin | Force logout all sessions |

### 5.2 Request/Response Examples

#### Register Request
```json
{
  "email": "john.doe@example.com",
  "password": "SecurePass123",
  "confirm_password": "SecurePass123",
  "role_id": 2,
  "firstname": "John",
  "lastname": "Doe",
  "fullname": "John Doe",
  "mobile_number": "0812345678",
  "phone_number": "021234567",
  "line_id": "johndoe_line",
  "location_id": "loc_001"
}
```

#### Register Response
```json
{
  "success": true,
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "john.doe@example.com",
    "role_id": 2,
    "firstname": "John",
    "lastname": "Doe",
    "fullname": "John Doe",
    "mobile_number": "0812345678",
    "phone_number": "021234567",
    "line_id": "johndoe_line",
    "location_id": "loc_001",
    "status": 1,
    "is_superuser": false,
    "verified": false,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

#### SignIn Response
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJSUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJSUzI1NiIs...",
    "token_type": "Bearer"
  }
}
```

---

## 6. คู่มือการบำรุงรักษา (Maintenance Guide)

### 6.1 Daily Maintenance Tasks

```bash
# 1. Check Redis memory usage
redis-cli INFO memory

# 2. Check Asynq queue size
redis-cli LLEN asynq:queues:critical

# 3. Monitor failed jobs
redis-cli ZCOUNT asynq:failed "(" "+inf"

# 4. Check database connections
psql -c "SELECT count(*) FROM pg_stat_activity;"

# 5. View logs
docker-compose logs --tail=100 app
```

### 6.2 Weekly Maintenance Tasks

```bash
# 1. Clean expired refresh tokens (Redis)
redis-cli --scan --pattern "RefreshToken:*" | xargs redis-cli DEL

# 2. Analyze database
psql -c "ANALYZE sd_users;"

# 3. Check for unverified users (>7 days)
psql -c "SELECT COUNT(*) FROM sd_users WHERE verified=false AND created_date < NOW() - INTERVAL '7 days';"

# 4. Backup Redis
redis-cli SAVE
```

### 6.3 Troubleshooting Common Issues

| Issue | Possible Cause | Solution |
|-------|---------------|----------|
| Token validation failed | Wrong public key | Check JWT config in config.yaml |
| Email not sent | SMTP configuration | Verify SMTP credentials |
| Redis connection refused | Redis down | Restart Redis: `docker restart redis` |
| Slow query | Missing index | Run: `CREATE INDEX CONCURRENTLY...` |
| High memory usage | Too many cached users | Adjust TTL or implement LRU |

---

## 7. คู่มือการขยายหรือแก้ไข (Extension Guide)

### 7.1 การเพิ่ม API ใหม่ (Adding New API)

**ขั้นตอนที่ 1: เพิ่ม Presenter**

```go
// File: internal/modules/users/presenter/presenters.go
// Thai: เพิ่ม DTO สำหรับ API ใหม่
// English: Add DTO for new API

type UserChangeStatus struct {
    Status int16 `json:"status" validate:"required,oneof=0 1"`
    // Thai: สถานะผู้ใช้ (0= inactive, 1= active)
    // English: User status (0= inactive, 1= active)
}
```

**ขั้นตอนที่ 2: เพิ่ม Interface**

```go
// File: internal/modules/users/usecase.go
// Thai: เพิ่ม method ใน UserUseCaseI
// English: Add method to UserUseCaseI

type UserUseCaseI interface {
    // ... existing methods ...
    
    // ChangeStatus เปลี่ยนสถานะผู้ใช้
    // ChangeStatus changes user status (active/inactive)
    ChangeStatus(ctx context.Context, id uuid.UUID, status int16) (*models.SdUser, error)
}
```

**ขั้นตอนที่ 3: Implement UseCase**

```go
// File: internal/modules/users/usecase/usecase.go
// Thai: Implement logic ใน usecase
// English: Implement logic in usecase

func (u *userUseCase) ChangeStatus(ctx context.Context, id uuid.UUID, status int16) (*models.SdUser, error) {
    // Thai: ตรวจสอบว่าสถานะถูกต้อง
    // English: Validate status
    if status != 0 && status != 1 {
        return nil, httpErrors.ErrValidation(errors.New("status must be 0 or 1"))
    }
    
    // Thai: อัปเดตสถานะ
    // English: Update status
    updatedUser, err := u.Update(ctx, id, map[string]interface{}{"status": status})
    if err != nil {
        return nil, err
    }
    
    // Thai: ล้าง cache
    // English: Clear cache
    _ = u.redisRepo.Delete(ctx, u.GenerateRedisUserKey(id))
    
    return updatedUser, nil
}
```

**ขั้นตอนที่ 4: เพิ่ม Handler**

```go
// File: internal/modules/users/delivery/http/handlers.go
// Thai: เพิ่ม handler method
// English: Add handler method

// ChangeStatus - PATCH /user/{id}/status (admin only)
// Thai: เปลี่ยนสถานะผู้ใช้ (admin เท่านั้น)
// English: Change user status (admin only)
func (h *userHandler) ChangeStatus() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Thai: ดึง ID จาก URL
        // English: Get ID from URL
        id, err := uuid.Parse(chi.URLParam(r, "id"))
        if err != nil {
            render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
            return
        }
        
        // Thai: อ่าน request body
        // English: Read request body
        req := new(presenter.UserChangeStatus)
        if err := json.NewDecoder(r.Body).Decode(req); err != nil {
            render.Render(w, r, responses.CreateErrorResponse(err))
            return
        }
        
        // Thai: เรียก usecase
        // English: Call usecase
        updatedUser, err := h.usersUC.ChangeStatus(r.Context(), id, req.Status)
        if err != nil {
            render.Render(w, r, responses.CreateErrorResponse(err))
            return
        }
        
        // Thai: ส่ง response
        // English: Send response
        render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedUser)))
    }
}
```

**ขั้นตอนที่ 5: เพิ่ม Route**

```go
// File: internal/modules/users/delivery/http/routes.go
// Thai: เพิ่ม route ใน MapUserRoute
// English: Add route in MapUserRoute

func MapUserRoute(router *chi.Mux, h users.Handlers, mw *middleware.MiddlewareManager) {
    // ... existing routes ...
    
    router.Route("/user", func(r chi.Router) {
        r.Group(func(r chi.Router) {
            r.Use(mw.Verifier(true))
            r.Use(mw.Authenticator())
            r.Use(mw.CurrentUser())
            r.Use(mw.ActiveUser())
            
            r.Route("/{id}", func(r chi.Router) {
                r.Group(func(r chi.Router) {
                    r.Use(mw.SuperUser())
                    // ... existing routes ...
                    r.Patch("/status", h.ChangeStatus())  // เพิ่ม/Add this line
                })
            })
        })
    })
}
```

**ขั้นตอนที่ 6: อัปเดต Interface Handler**

```go
// File: internal/modules/users/handler.go
// Thai: เพิ่ม method ใน Handlers interface
// English: Add method to Handlers interface

type Handlers interface {
    // ... existing methods ...
    ChangeStatus() http.HandlerFunc  // เพิ่ม/Add this line
}
```

### 7.2 การเพิ่ม Field ใน User Model

**ขั้นตอนที่ 1: อัปเดต Model**

```go
// File: internal/models/sd_user.go
// Thai: เพิ่ม field ใหม่ใน SdUser struct
// English: Add new field to SdUser struct

type SdUser struct {
    // ... existing fields ...
    
    // Thai: รหัสอ้างอิงพนักงาน
    // English: Employee reference code
    EmployeeCode *string `json:"employee_code,omitempty" gorm:"column:employee_code;type:varchar(50)"`
}
```

**ขั้นตอนที่ 2: สร้าง Migration**

```sql
-- File: migrations/xxxxxx_add_employee_code.sql
-- Thai: เพิ่มคอลัมน์ employee_code
-- English: Add employee_code column

ALTER TABLE sd_users ADD COLUMN employee_code VARCHAR(50);

CREATE INDEX idx_sd_users_employee_code ON sd_users(employee_code);
```

**ขั้นตอนที่ 3: อัปเดต Presenter**

```go
// File: internal/modules/users/presenter/presenters.go
// Thai: เพิ่ม field ใน UserCreate และ UserUpdate
// English: Add field to UserCreate and UserUpdate

type UserCreate struct {
    // ... existing fields ...
    EmployeeCode string `json:"employee_code,omitempty" validate:"omitempty,max=50"`
}

type UserUpdate struct {
    // ... existing fields ...
    EmployeeCode *string `json:"employee_code,omitempty"`
}

type UserResponse struct {
    // ... existing fields ...
    EmployeeCode *string `json:"employee_code,omitempty"`
}
```

**ขั้นตอนที่ 4: อัปเดต Mapping Functions**

```go
// File: internal/modules/users/delivery/http/handlers.go
// Thai: แก้ไข mapModel และ mapModelResponse
// English: Update mapModel and mapModelResponse

func mapModel(req *presenter.UserCreate) *models.SdUser {
    return &models.SdUser{
        // ... existing fields ...
        EmployeeCode: stringPtr(req.EmployeeCode),
    }
}

func mapModelResponse(user *models.SdUser) *presenter.UserResponse {
    return &presenter.UserResponse{
        // ... existing fields ...
        EmployeeCode: user.EmployeeCode,
    }
}
```

**ขั้นตอนที่ 5: อัปเดต Update Logic**

```go
// File: internal/modules/users/delivery/http/handlers.go
// Thai: เพิ่ม employee_code ใน Update handler
// English: Add employee_code to Update handler

func (h *userHandler) Update() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // ... existing code ...
        
        values := make(map[string]interface{})
        // ... existing fields ...
        if req.EmployeeCode != nil {
            values["employee_code"] = *req.EmployeeCode
        }
        
        // ... rest of the code ...
    }
}
```

---

## 8. Checklist Test Module

### 8.1 ไฟล์และโครงสร้าง (Files and Structure)

| # | File Path | Status | Description |
|---|-----------|--------|-------------|
| 1 | `internal/modules/users/handler.go` | ✅ | HTTP handler interface |
| 2 | `internal/modules/users/pg_repository.go` | ✅ | PostgreSQL interface |
| 3 | `internal/modules/users/redis_repository.go` | ✅ | Redis interface |
| 4 | `internal/modules/users/usecase.go` | ✅ | UseCase interface |
| 5 | `internal/modules/users/worker.go` | ✅ | Async task definitions |
| 6 | `internal/modules/users/delivery/http/handlers.go` | ✅ | HTTP handlers implementation |
| 7 | `internal/modules/users/delivery/http/routes.go` | ✅ | Route registration |
| 8 | `internal/modules/users/distributor/distributor.go` | ✅ | Task distributor |
| 9 | `internal/modules/users/presenter/presenters.go` | ✅ | DTO definitions |
| 10 | `internal/modules/users/processor/processor.go` | ✅ | Task processor |
| 11 | `internal/modules/users/repository/pg_repository.go` | ✅ | PostgreSQL implementation |
| 12 | `internal/modules/users/repository/redis_repository.go` | ✅ | Redis implementation |
| 13 | `internal/modules/users/usecase/usecase.go` | ✅ | Business logic |

### 8.2 ฟังก์ชันที่ต้องมี (Required Functions)

| # | Function | File | Status |
|---|----------|------|--------|
| 1 | `Register()` | handlers.go | ✅ |
| 2 | `Create()` | handlers.go | ✅ |
| 3 | `Get()` | handlers.go | ✅ |
| 4 | `GetMulti()` | handlers.go | ✅ |
| 5 | `Update()` | handlers.go | ✅ |
| 6 | `Delete()` | handlers.go | ✅ |
| 7 | `UpdatePassword()` | handlers.go | ✅ |
| 8 | `UpdateRole()` | handlers.go | ✅ |
| 9 | `Me()` | handlers.go | ✅ |
| 10 | `UpdateMe()` | handlers.go | ✅ |
| 11 | `UpdatePasswordMe()` | handlers.go | ✅ |
| 12 | `LogoutAllAdmin()` | handlers.go | ✅ |

### 8.3 UseCase Functions

| # | Function | File | Status |
|---|----------|------|--------|
| 1 | `CreateUser()` | usecase.go | ✅ |
| 2 | `SignIn()` | usecase.go | ✅ |
| 3 | `IsActive()` | usecase.go | ✅ |
| 4 | `IsSuper()` | usecase.go | ✅ |
| 5 | `CreateSuperUserIfNotExist()` | usecase.go | ✅ |
| 6 | `UpdatePassword()` | usecase.go | ✅ |
| 7 | `ParseIdFromRefreshToken()` | usecase.go | ✅ |
| 8 | `Refresh()` | usecase.go | ✅ |
| 9 | `Logout()` | usecase.go | ✅ |
| 10 | `LogoutAll()` | usecase.go | ✅ |
| 11 | `Verify()` | usecase.go | ✅ |
| 12 | `ForgotPassword()` | usecase.go | ✅ |
| 13 | `ResetPassword()` | usecase.go | ✅ |

### 8.4 Repository Functions

| # | Function | File | Status |
|---|----------|------|--------|
| 1 | `GetByEmail()` | pg_repository.go | ✅ |
| 2 | `UpdatePassword()` | pg_repository.go | ✅ |
| 3 | `UpdateVerificationCode()` | pg_repository.go | ✅ |
| 4 | `UpdateVerification()` | pg_repository.go | ✅ |
| 5 | `GetByVerificationCode()` | pg_repository.go | ✅ |
| 6 | `UpdatePasswordReset()` | pg_repository.go | ✅ |
| 7 | `GetByResetToken()` | pg_repository.go | ✅ |
| 8 | `GetByResetTokenResetAt()` | pg_repository.go | ✅ |
| 9 | `UpdatePasswordResetToken()` | pg_repository.go | ✅ |

### 8.5 API Endpoints Checklist

| # | Method | Endpoint | Auth | Status |
|---|--------|----------|------|--------|
| 1 | POST | /register | ❌ | ✅ |
| 2 | POST | /signin | ❌ | ✅ |
| 3 | POST | /refresh | ❌ | ✅ |
| 4 | POST | /logout | ✅ | ✅ |
| 5 | POST | /forgotpassword | ❌ | ✅ |
| 6 | POST | /resetpassword | ❌ | ✅ |
| 7 | GET | /verifyemail | ❌ | ✅ |
| 8 | GET | /user/me | ✅ | ✅ |
| 9 | PUT | /user/me | ✅ | ✅ |
| 10 | PATCH | /user/me/updatepass | ✅ | ✅ |
| 11 | POST | /user | ✅ (Admin) | ✅ |
| 12 | GET | /user | ✅ (Admin) | ✅ |
| 13 | GET | /user/{id} | ✅ (Admin) | ✅ |
| 14 | PUT | /user/{id} | ✅ (Admin) | ✅ |
| 15 | DELETE | /user/{id} | ✅ (Admin) | ✅ |
| 16 | PATCH | /user/{id}/role | ✅ (Admin) | ✅ |
| 17 | PATCH | /user/{id}/updatepass | ✅ (Admin) | ✅ |
| 18 | GET | /user/{id}/logoutall | ✅ (Admin) | ✅ |

### 8.6 การเพิ่มฟังก์ชันใหม่ (Adding New Functions)

**วิธีการ (How to):**

1. **เพิ่ม Presenter** - สร้าง DTO ใน `presenter/presenters.go`
2. **เพิ่ม Interface** - เพิ่ม method ใน `usecase.go` และ `handler.go`
3. **Implement UseCase** - เขียน business logic ใน `usecase/usecase.go`
4. **Implement Handler** - สร้าง HTTP handler ใน `delivery/http/handlers.go`
5. **Register Route** - เพิ่ม route ใน `delivery/http/routes.go`
6. **Update Repository (ถ้าจำเป็น)** - เพิ่ม method ใน `pg_repository.go`

**ข้อควรระวัง (Precautions):**
- ตรวจสอบการ validate input ก่อนทุกครั้ง
- ใช้ context timeout สำหรับ long-running operations
- ล้าง Redis cache เมื่อมีการ update/delete
- ตรวจสอบสิทธิ์ (authorization) สำหรับ admin endpoints
- Log error ทุกครั้งที่เกิด panic หรือ unexpected error

---

## 9. คอมเมนต์โค้ด (Code Comments)

### ตัวอย่างการคอมเมนต์ (Comment Example)

```go
// CreateUser - สร้างผู้ใช้ใหม่พร้อมส่งอีเมลยืนยัน
// CreateUser - Creates a new user and sends verification email
// Parameters:
//   - ctx: Context สำหรับ timeout และ cancellation / Context for timeout and cancellation
//   - exp: ข้อมูลผู้ใช้ที่ต้องการสร้าง / User data to create
//   - confirmPassword: รหัสผ่านยืนยัน / Password confirmation
// Returns:
//   - *models.SdUser: ข้อมูลผู้ใช้ที่สร้างแล้ว / Created user data
//   - error: ข้อผิดพลาดที่เกิดขึ้น / Error if any
// Usage:
//   - ใช้สำหรับ API /register และ /user (admin) / Used for /register and /user (admin) APIs
// Note:
//   - รหัสผ่านจะถูกเข้ารหัสด้วย bcrypt ก่อนบันทึก / Password is hashed with bcrypt before saving
//   - อีเมลจะถูกแปลงเป็น lowercase / Email is converted to lowercase
//   - จะส่งอีเมลยืนยันเฉพาะเมื่อ verified = false / Verification email sent only when verified = false
func (u *userUseCase) CreateUser(ctx context.Context, exp *models.SdUser, confirmPassword string) (*models.SdUser, error) {
    // ตรวจสอบว่ารหัสผ่านตรงกัน / Check if passwords match
    if exp.Password != confirmPassword {
        return nil, httpErrors.ErrValidation(errors.New("password do not match"))
    }
    return u.Create(ctx, exp)
}
```

---

## 10. สรุป (Summary)

โมดูล User เป็นระบบที่สมบูรณ์สำหรับการจัดการผู้ใช้ รองรับ:
- ✅ Authentication ด้วย JWT (RS256)
- ✅ Authorization 3 ระดับ (User, Admin, Super Admin)
- ✅ Email Verification ผ่าน Asynq Queue
- ✅ Forgot/Reset Password
- ✅ Redis Cache สำหรับลดภาระ Database
- ✅ Refresh Token Management
- ✅ Global Logout
- ✅ ระบบ Logging และ Error Handling

**ข้อดี (Advantages):**
- Clean Architecture แยกชั้นชัดเจน
- Scalable ด้วย Async Task Queue
- High Performance ด้วย Redis Cache
- Secure ด้วย JWT และ bcrypt

**ข้อเสีย (Disadvantages):**
- Complexity สูง เหมาะกับ project ขนาดกลาง-ใหญ่
- ต้องพึ่งพา Redis และ Asynq



