# คู่มือการใช้งาน Worker Service (Asynq Task Processor)

Worker เป็น **service แยกต่างหาก** ในโปรเจกต์ `icmongolang` — ประมวลผลงานเบื้องหลัง (background tasks) แบบ asynchronous ด้วยไลบรารี **`github.com/hibiken/asynq`** (Redis-backed task queue)

**งานปัจจุบัน:** ส่งอีเมล (`task:send_email`) — ใช้สำหรับ verification email และ password reset email จาก users module

**สถาปัตยกรรม:**

```
┌─────────────────────────────┐         ┌─────────────────────┐
│  API Server (producer)      │         │  Worker (consumer)   │
│  users usecase              │         │  cmd/worker.go       │
│  DistributeTaskSendEmail    │         │  internal/worker/    │
└──────────────┬──────────────┘         │  asynq.Server        │
               │ enqueue JSON           └──────────┬──────────┘
               ▼                                   │
        ┌──────────────────┐     dequeue          │
        │   Redis (DB 1)   │ ◄────────────────────┘
        │  asynq lists/    │
        │  zsets           │
        └──────────────────┘
```

> **สำคัญ:** worker ต้องรันคู่กับ API server — API server เป็นคน *enqueue* งาน ส่วน worker เป็นคน *consume* ถ้ารันแค่ API โดยไม่มี worker งานจะค้างใน Redis

---

## 1. ตารางไฟล์

| ไฟล์ | บทบาท |
|------|-------|
| `cmd/worker.go` | Cobra subcommand `worker` — รัน TaskProcessor |
| `internal/worker/worker.go` | `TaskProcessor` — asynq server + queue config + task registration |
| `internal/modules/users/worker.go` | Task type constant `task:send_email` + `PayloadSendEmail` |
| `internal/modules/users/processor/processor.go` | `ProcessTaskSendEmail` — รับงานแล้วส่งอีเมล |
| `internal/modules/users/distributor/distributor.go` | ฝั่ง enqueue (API server) |
| `internal/distributor/distributor.go` | Base asynq client builder |
| `pkg/sendEmail/sendEmail.go` | ตัวส่งอีเมลจริง (gomail) |
| `pkg/emailTemplates/` | สร้าง HTML template (hermes) |

---

## 2. การตั้งค่า Configuration

### 2.1 `config/config.default.yml`
```yaml
taskRedis:
  Addr: localhost:6379
  Db: 1            # ⚠️ ใช้ Redis DB 1 (ต่างจาก redis หลัก DB 0)
  PoolTimeout: 240 # วินาที → ใช้เป็น DialTimeout/ReadTimeout/WriteTimeout

smtpEmail:
  Host: "smtp.icmon.io"
  Port: 2525
  User: ""
  Password: ""
  UseTls: true
  UseSsl: false
  Timeout: 30

email:
  From: "kongnakornna@gmail.com"
  VerificationSubject: "Your account verification code"
  ResetSubject: "Your account password reset token"
```

### 2.2 Env vars (`config/config.go` BindEnvs)
`TASK_REDIS_ADDR`, `TASK_REDIS_DB`, `TASK_REDIS_POOLTIMEOUT`

---

## 3. การรัน Service

### วิธีที่ 1 — รันตรง
```bash
go run main.go worker
```

### วิธีที่ 2 — Docker
```bash
docker compose up -d worker
```
docker-compose: `wait-for-it -w db:5435 -w redis:6379 -- sleep 10; go run main.go worker`

---

## 4. Worker Internals

### 4.1 Queue config (`internal/worker/worker.go:37-49`)
| ตั้งค่า | ค่า |
|--------|-----|
| Queues | `critical` (weight 10), `default` (weight 5) |
| Concurrency | ไม่ตั้ง → `runtime.NumCPU()` |
| ErrorHandler | print `Err/Type/Payload/Msg: process task failed` ไป stdout |
| Logger | app logger (`pkg/logger`) |

### 4.2 Task ที่ลงทะเบียน (`worker.go:66`)
```go
mux.HandleFunc(users.TaskSendEmail, userRedisTaskProcessor.ProcessTaskSendEmail)
```

### 4.3 Payload — `PayloadSendEmail` (`users/worker.go:18-24`)
| Go field | JSON tag |
|----------|----------|
| `From` | `from` |
| `To` | `to` |
| `Subject` | `subject` |
| `BodyHtml` | `bodyHtml` |
| `BodyPlain` | `bodyPlain` |

### 4.4 `ProcessTaskSendEmail` (`processor/processor.go:29-39`)
1. `json.Unmarshal(task.Payload(), &payload)`
   - JSON ผิด → return `asynq.SkipRetry` (**ไม่ลองใหม่** — abort ทันที)
2. `emailSender.SendEmail(ctx, payload.From, payload.To, payload.Subject, payload.BodyHtml, payload.BodyPlain)`
3. error → return error → **asynq retry** (ตาม `MaxRetry` ของ task)
4. สำเร็จ → log `"Msg: email sended"`

### 4.5 งานที่ enqueue (producer — users module)
`DistributeTaskSendEmail` ใช้ตัวเลือก:
```go
asynq.MaxRetry(10),           // ลองใหม่สูงสุด 10 ครั้ง
asynq.ProcessIn(10*time.Second), // ดีเลย์ 10 วิ
asynq.Queue(worker.QueueCritical) // ไปคิว "critical"
```

จุดที่เรียก:
- **Verification email** (`users/usecase/usecase.go:258-269`) — ลิงก์ `http://localhost:5000/auth/verifyemail?code=...`
- **Password reset email** (`users/usecase/usecase.go:588-599`) — ลิงก์ `http://localhost:5000/auth/resetpassword?code=...`

### 4.6 Email sender (`pkg/sendEmail/sendEmail.go`)
- Library: `gopkg.in/gomail.v2`
- `SetHeader(From/To/Subject)`, `SetBody(text/html)`, `AddAlternative(text/plain)`
- Dialer: `cfg.SmtpEmail.{Host,Port,User,Password}`
- ⚠️ `TLS InsecureSkipVerify: true` hardcode — `UseTls`/`UseSsl` ไม่ถูกอ่าน
- Timeout: `cfg.SmtpEmail.Timeout` (default 30s), error `"smtp timeout after %d seconds"`

---

## 5. โครงสร้างข้อมูล (Redis)

**ไม่มีตาราง PostgreSQL** — asynq ใช้ Redis ทั้งหมด (lists/zsets/hashes ภายใต้ prefix `asynq:` ใน DB 1)

| Redis object | ความหมาย |
|--------------|----------|
| `asynq:{queue}:pending` | งานรอประมวลผล |
| `asynq:{queue}:active` | งานกำลังประมวลผล |
| `asynq:delayed` | งานดีเลย์ (ProcessIn) |
| `asynq:retry` | งานรอ retry |
| `asynq:dead` | งานที่หมด retry แล้ว (dead letter) |

---

## 6. Edge Cases & ข้อควรระวัง

1. **queue `low` ถูกนิยามแต่ไม่ได้ register** — task ที่ระบุ queue `low` จะถูก reject
2. **`template:send_notification` processor/distributor เป็น dead code** — ไม่มีใคร register ใน mux
3. **worker ขึ้นกับ Redis DB 1** — ต้องตรงกับ `taskRedis.Db` ของ API server (คนเดียวกัน)
4. **payload JSON ผิด → SkipRetry** (ไม่ลองใหม่) แต่ **ส่งอีเมลล้มเหลว → retry** สูงสุด 10 ครั้ง
5. **ลิงก์ในอีเมล hardcode `localhost:5000`** — ไม่ได้ใช้ config BaseUrl
6. **Concurrency = NumCPU** — ขนาดงานพร้อมกันจำกัดตาม CPU ของเครื่อง
7. **ต้องมีทั้ง producer และ consumer** — ถ้า worker ไม่รัน งานค้างใน Redis และจะถูกประมวลผลย้อนหลังเมื่อ worker กลับมา (asynq default retry-delay)

---

## 7. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| ไม่มีอีเมลส่ง | worker ไม่รัน → `docker compose up -d worker` หรือ `go run main.go worker` |
| งานค้างใน Redis ตลอด | `taskRedis.Addr/Db` ระหว่าง API และ worker ไม่ตรงกัน |
| อีเมล timeout | SMTP host/port/credential ผิด หรือ `smtpEmail.Timeout` น้อยเกินไป |
| งานตายใน dead letter | ลองครบ 10 ครั้งแล้ว (ดู `asynq:dead`) — ตรวจ SMTP log |
| worker fatal ตอน start | Redis ยังไม่พร้อม → ดู `taskRedis` config |
| อยากดูสถานะคิว | ใช้ asynq CLI / Redis: `keys asynq:*` ใน DB 1 |
