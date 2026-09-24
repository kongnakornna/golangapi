# คู่มือการใช้งานโมดูล Email

โมดูล `internal/modules/email` เป็นระบบจัดการอีเมล: ส่งอีเมล, ดู log การส่ง, ดู/แก้ SMTP config โดยทุก endpoint ต้องยืนยันตัวตน (JWT Bearer token) + ผู้ใช้ต้อง active

> **⚠️ ข้อควรทราบก่อนใช้งาน:** โมดูลนี้ในโค้ดปัจจุบันยังเป็น **โครงสร้างครึ่งเดียว (half-implemented)** — `POST /send` ยังไม่ได้ส่งอีเมลจริง (แค่เขียน log), `GET /logs/{id}` และ `PUT /config` ยังเป็น stub ที่คืน `"ok"` เสมอ การส่งอีเมลจริงเกิดขึ้นผ่าน **worker service (asynq)** ในส่วนของ users module (verification/reset password) — ดูรายละเอียดในส่วนข้อจำกัดด้านล่าง

**สถาปัตยกรรม:**

```
Client
 │ JWT Bearer token (ต้องเป็น active user)
 ▼
┌──────────────────────────────────────┐
│ delivery/http/handlers.go             │
│  SendEmail / ListEmailLogs / GetConfig│
└──────────────────┬───────────────────┘
                   ▼
┌──────────────────────────────────────┐
│ usecase/usecase.go (emailUseCase)     │
│  SendEmail → เขียน email_log          │
│  GetConfig / UpdateConfig            │
└──────────────────┬───────────────────┘
                   ▼
┌──────────────────────────────────────┐
│ repository/pg_repository.go          │
│  GORM → email_log, email_config      │
└──────────────────────────────────────┘
```

---

## 1. ตารางไฟล์

| ไฟล์ | บทบาท |
|------|-------|
| `delivery/http/routes.go` | ลงทะเบียน route `/api/email/*` |
| `delivery/http/handlers.go` | HTTP handlers |
| `presenter/presenters.go` | DTOs (ใช้ใน swagger เป็นหลัก) |
| `usecase/usecase.go` | Business logic (embed generic usecase) |
| `repository/pg_repository.go` | GORM repository |
| `handler.go` / `usecase.go` / `pg_repository.go` | Interfaces ที่ root ของ module |
| `pkg/sendEmail/sendEmail.go` | ตัวส่งอีเมลจริง (gomail) — ใช้โดย worker ไม่ใช่ module email |
| `pkg/emailTemplates/` | สร้าง HTML template (hermes) |

---

## 2. การตั้งค่า Configuration

### 2.1 SMTP — `config/config.default.yml`
```yaml
smtpEmail:
  Host: "smtp.icmon.io"
  Port: 2525
  User: ""
  Password: ""
  UseTls: true
  UseSsl: false
  Timeout: 30
```

### 2.2 Template/branding — `config/config.default.yml`
```yaml
email:
  From: "kongnakornna@gmail.com"
  Name: "Go restapi"
  Link: "https://icmongolang"
  LogoLink: "https://media.vov.vn/...apple-logo...png"
  Copyright: "Copyright © 2026 example inc. All rights reserved."
  VerificationSubject: "Your account verification code"
  ResetSubject: "Your account password reset token"
```

> `config.dev.yml` ใช้ `User: icmon`, `Password: icmon`

---

## 3. API Reference

ทุก route อยู่ใต้ **`/api/email`** และใช้ middleware chain: `Verifier(true)` → `Authenticator()` → `CurrentUser()` → `ActiveUser()`

### 3.1 ตารางรวม endpoints

| Method | Path | คำอธิบาย |
|--------|------|----------|
| POST | `/api/email/send` | ส่งอีเมล |
| GET | `/api/email/logs` | รายการ email log |
| GET | `/api/email/logs/{id}` | รายละเอียด email log ตาม id |
| GET | `/api/email/config` | ดู SMTP config |
| PUT | `/api/email/config` | อัปเดต SMTP config |

---

### 3.2 `POST /api/email/send`

**Request:**
```json
{
  "to": "recipient@example.com",
  "cc": "",
  "bcc": "",
  "subject": "Test email",
  "body": "<p>Hello from ICMON</p>"
}
```

| field | type | บังคับ | validation |
|-------|------|--------|------------|
| `to` | string | ✅ | `required,email` |
| `cc` | string | - | |
| `bcc` | string | - | |
| `subject` | string | ✅ | `required` |
| `body` | string | ✅ | `required` |

**Response 200:**
```json
{ "data": "Email sent", "is_success": true }
```

> ⚠️ จริงๆ แล้ว **ไม่มีการส่งอีเมล** — โค้ดใน `usecase/usecase.go:33-44` มี TODO `// TODO: delegate to pkg/sendEmail` แล้วแค่เขียน `email_log` ด้วย `Status:"sent"` ลง DB (`cc`/`bcc` ถูกอ่านแต่ไม่ถูกส่งต่อ)

---

### 3.3 `GET /api/email/logs`

รายการ email log — ใช้ `GetMulti(50, 0)` **hardcode** (ignore query param `limit`/`offset`)

**Response 200:** (serialize `models.EmailLog` โดยตรง → ใช้ Go field names)
```json
{
  "data": [
    {
      "ID": 1,
      "To": "recipient@example.com",
      "Cc": null,
      "Bcc": null,
      "Subject": "Test email",
      "Body": "<p>Hello</p>",
      "Status": "sent",
      "ErrorMessage": null,
      "SentAt": null,
      "CreatedAt": "2026-06-08T12:00:00+07:00"
    }
  ],
  "is_success": true
}
```

---

### 3.4 `GET /api/email/logs/{id}` — **STUB**

Parse `id` แล้ว **discard ทันที** → คืน `{"data":"ok","is_success":true}` เสมอ

---

### 3.5 `GET /api/email/config`

คืน `*models.EmailConfig` (row `is_active=true` อันแรก) — serialize ด้วย Go field names:

```json
{
  "data": {
    "ID": 1,
    "SmtpHost": "smtp.gmail.com",
    "SmtpPort": 587,
    "SmtpUser": "noreply@example.com",
    "SmtpPass": "app_password",
    "FromEmail": "noreply@example.com",
    "FromName": "ICMON System",
    "IsActive": true,
    "CreatedAt": "...",
    "UpdatedAt": "..."
  },
  "is_success": true
}
```

### 3.6 `PUT /api/email/config` — **STUB**

**ไม่ทำอะไรเลย** — ไม่ decode body, ไม่เรียก usecase → คืน `{"data":"ok","is_success":true}` เสมอ (method `UpdateConfig` ใน usecase/repo มีอยู่แต่ **reachable ไม่ได้ผ่าน HTTP**)

---

## 4. โครงสร้างฐานข้อมูล

### `email_log` (`migrations/20260712_new_modules_schema.sql:492`)
| column | type | หมายเหตุ |
|--------|------|----------|
| `id` | SERIAL | PK |
| `to` | VARCHAR(255) | NOT NULL |
| `cc` | VARCHAR(500) | |
| `bcc` | VARCHAR(500) | |
| `subject` | VARCHAR(500) | NOT NULL |
| `body` | TEXT | NOT NULL |
| `status` | VARCHAR(20) | DEFAULT `'pending'` |
| `error_message` | TEXT | |
| `sent_at` | TIMESTAMP | |
| `created_at` | TIMESTAMP | DEFAULT NOW() |

### `email_config` (`20260712_new_modules_schema.sql:505`)
| column | type |
|--------|------|
| `id` | SERIAL PK |
| `smtp_host` | VARCHAR(255) NOT NULL |
| `smtp_port` | INTEGER NOT NULL |
| `smtp_user` | VARCHAR(255) NOT NULL |
| `smtp_pass` | VARCHAR(255) NOT NULL |
| `from_email` | VARCHAR(255) NOT NULL |
| `from_name` | VARCHAR(255) NOT NULL |
| `is_active` | BOOLEAN DEFAULT TRUE |
| `created_at` / `updated_at` | TIMESTAMP |

Seed data: `migrations/20260712_seed_data.sql:12` — Gmail SMTP `noreply@example.com` / `app_password`

---

## 5. ระบบส่งอีเมลจริง (ใช้โดย worker)

### `pkg/sendEmail/sendEmail.go`
- Library: `gopkg.in/gomail.v2`
- Interface: `SendEmail(ctx, from, to, subject, bodyHtml, bodyPlain) error`
- ส่ง HTML + plain-text alternative
- Dialer จาก `cfg.SmtpEmail.{Host,Port,User,Password}`
- ⚠️ `d.TLSConfig = &tls.Config{InsecureSkipVerify: true}` **hardcode** — `UseTls`/`UseSsl` ไม่ถูกอ่านจริง
- Timeout = `cfg.SmtpEmail.Timeout` (default 30s); error: `"smtp timeout after %d seconds"`

### `pkg/emailTemplates/` (hermes v2)
- `GenerateVerificationCodeTemplate(ctx, name, link)` → ปุ่มสีเขียว `#22BC66` "Confirm your account"
- `GeneratePasswordResetTemplate(ctx, name, link)` → ปุ่มสีแดง `#DC4D2F` "Reset your password"
- ⚠️ template error ถูก swallow (return `("", "", nil)`)

### เส้นทางส่งจริง
`users module` → สร้าง template → enqueue asynq task `task:send_email` → **worker** consume → `sendEmail.SendEmail` → SMTP
(ดูคู่มือ `README_Worker.md`)

---

## 6. Edge Cases & ข้อควรระวัง

1. **`POST /send` ยังไม่ส่งอีเมลจริง** — แค่เขียน log; ถ้าต้องการส่งจริงต้องใช้ users verification/reset flow ผ่าน worker
2. **`GET /logs/{id}` และ `PUT /config` เป็น stub** — คืน `"ok"` เสมอ
3. **`GET /logs` ignore limit/offset** — hardcode (50, 0)
4. **JSON field names ต่างจาก swagger** — จริงๆ serialize ใช้ Go field names (`SmtpHost`, `SmtpPass`, `ErrorMessage`, `SentAt`) ไม่ใช่ camelCase ตามที่ swagger/postman เขียนไว้
5. **ทุก endpoint ต้อง JWT + active user** — token ไม่มี → 401; user inactive → 403
6. **`cc`/`bcc` ถูกอ่านแต่ไม่ถูกใช้** ใน `SendEmail`
7. **GetConfig คืน row `is_active=true` อันแรก** — ถ้าไม่มีเลย → 404/error

---

## 7. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| 401 เรียก API ไม่ได้ | ไม่มี/ผิด Bearer token |
| 403 | user ไม่ active (`status != 1`) |
| ส่งแล้ว log มีแต่ไม่มีอีเมลเข้าคนรับ | เป็นการทำงานที่คาดไว้ — module email ยังไม่ส่งจริง; ใช้ flow verification/reset (worker) แทน |
| Config field เป็น PascalCase | เป็นพฤติกรรมปัจจุบันของโค้ด — ดู field ที่ชื่อตาม struct Go |
| ส่งอีเมลจาก worker timeout | SMTP host/port/credential ผิด หรือ timeout น้อยไป (`smtpEmail.Timeout`) |
