# คู่มือการใช้งาน PostgreSQL (ฐานข้อมูลหลัก)

PostgreSQL เป็นฐานข้อมูลหลักของ `icmongolang` — เก็บข้อมูลทุกโมดูล (users, auth, email, iot device/data, kafka orders, websocket, batch, document, purchase order, ฯลฯ) ผ่าน GORM

**ในระบบนี้มีอยู่ 2 รูปแบบ:**
1. **PostgreSQL ใน Docker** — `docker compose up -d db` (ใช้ port `5435`)
2. **PostgreSQL ท้องถิ่น** — ผ่าน config `postgres.*` (default port 5432)

---

## 1. การตั้งค่า Configuration

### 1.1 `config/config.default.yml`
```yaml
postgres:
  Host: localhost
  Port: 5432
  User: "postgres"
  Password: "postgres"
  Dbname: icmongolang
  ConnectionTimeout: 10
```

### 1.2 Docker Compose (`docker-compose.yml:42-53`)
```yaml
db:
  image: postgres:15.2-alpine
  environment:
    - POSTGRES_USER=${POSTGRES_USER}
    - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
    - POSTGRES_DB=${POSTGRES_DBNAME}
  command: ["postgres", "-p", "5435"]   # ⚠️ port 5435
  ports:
    - "5435:5435"
  volumes:
    - app-postgres-data:/var/lib/postgresql/data/
```

> **ข้อสำคัญ:**  container ใช้ port **5435** (ไม่ใช่ 5432) — ทั้งจาก host และใน docker network (`wait-for-it db:5435`)

### 1.3 Env vars
`POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DBNAME`

---

## 2. การเชื่อมต่อจาก Go

`pkg/db/postgres/db_conn.go` — `NewPsqlDB(cfg)`:
```go
dsn := fmt.Sprintf(
    "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok connect_timeout=%d",
    cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Dbname, cfg.Postgres.Port, timeout,
)
```
- GORM driver `gorm.io/driver/postgres`
- `sslmode=disable` เสมอ
- `TimeZone=Asia/Bangkok`
- `connect_timeout` = `cfg.Postgres.ConnectionTimeout` (default 10s)
- `DisableForeignKeyConstraintWhenMigrating: true` — ตอน migrate GORM ไม่สร้าง FK constraint (จัดการเองผ่าน SQL)

---

## 3. ตารางที่สำคัญ

ตารางหลักๆ (จาก migration):
| ตาราง | โมดูล | หมายเหตุ |
|-------|-------|----------|
| `sd_user` | auth/users | ผู้ใช้ระบบ |
| `refresh_tokens` | auth | (มีแต่ระบบจริงใช้ Redis แทน) |
| `email_log` / `email_config` | email | log ส่งเมล + SMTP config |
| `orders` | kafka | ตัวอย่าง async processing |
| `ws_messages` / `ws_sessions` | websocket | ข้อความ + session |
| `sd_iot_device` | iot | อุปกรณ์ IoT |
| `iot_data` | iot | ข้อมูลอนุกรมเวลา |
| `device_status` / `device_config` | iot | สถานะ + config |
| `activity_log` / `command_log` | iot | audit + คำสั่ง |
| `t_job` / `m_customer` / `t_quotation` / `t_purchaseorder` | job/customer/quotation/purchaseorder | โมดูลธุรกิจใหม่ |

ไฟล์ migration หลัก: `migrations/` (db.sql, icmon.sql, 20260712_new_modules_schema.sql, 20250619_*.sql, ...)

---

## 4. คำสั่งใช้งาน

### รัน PostgreSQL
```bash
docker compose up -d db
```

### เข้า psql
```bash
# host
docker exec -it $(docker compose ps -q db) psql -U postgres -d icmongolang
# หรือผ่าน port 5435
psql -h localhost -p 5435 -U postgres -d icmongolang
```

### รัน migration
```bash
go run main.go migrate
```

### รัน seed data
```bash
go run main.go initdata
```

---

## 5. Edge Cases & ข้อควรระวัง

1. **Port 5435 (ไม่ใช่ 5432)** ใน Docker — ถ้า config ท้องถิ่นยังเป็น 5432 จะเชื่อมคนละตัว
2. **sslmode=disable** — ไม่เข้ารหัส; เหมาะกับ internal network เท่านั้น
3. **เวลาอยู่ในโซน Asia/Bangkok** — ผ่าน DSN `TimeZone`
4. **`initdata` สร้าง superuser เริ่มต้น** จาก config `firstSuperUser.*`
5. **Migrations บางไฟล์เป็น full dump ขนาดใหญ่** (db.sql ~174KB) — มีทั้ง schema + data
6. **Worker ขึ้นกับ db แม้ไม่ใช้จริง** — docker-compose worker `depends_on: db`
7. **volume `app-postgres-data`** — ลบ container ไม่ได้ลบข้อมูล; ล้างด้วย `docker compose down -v`

---

## 6. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| `failed to connect` ตอน boot | db ยังไม่พร้อม → wait-for-it รอ `db:5435` |
| connection timeout | `postgres.ConnectionTimeout` น้อยไป หรือ db ยังไม่ up |
| `password authentication failed` | `POSTGRES_USER/PASSWORD` ไม่ตรงกับ config |
| migrate fail FK | ดู migration SQL เอง (GORM disable FK เมื่อ migrate) |
| ข้อมูลหาย | เช็คว่าไม่ใช้ `docker compose down -v` (ลบ volume) |
