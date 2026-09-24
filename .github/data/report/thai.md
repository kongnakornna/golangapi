# รายงานทางเทคนิคโครงการ ICMON Go — ฉบับสมบูรณ์ (ภาษาไทย)

> **เวอร์ชัน:** 1.0.0  
> **โมดูล:** `icmongolang`  
> **Go เวอร์ชัน:** 1.25.0  
> **ลิขสิทธิ์:** github.com/kongnakornna  
> **วันที่:** กรกฎาคม 2026

---

## สารบัญ

1. [ภาพรวมโครงการ](#1-ภาพรวมโครงการ)
2. [สถาปัตยกรรมระบบ](#2-สถาปัตยกรรมระบบ)
3. [เทคโนโลยีและ dependencies](#3-เทคโนโลยีและ-dependencies)
4. [การตั้งค่าระบบ](#4-การตั้งค่าระบบ)
5. [จุดเริ่มต้นโปรแกรมและคำสั่ง CLI](#5-จุดเริ่มต้นโปรแกรมและคำสั่ง-cli)
6. [โครงสร้างไดเรกทอรี](#6-โครงสร้างไดเรกทอรี)
7. [โครงสร้างฐานข้อมูล](#7-โครงสร้างฐานข้อมูล)
8. [เส้นทาง API](#8-เส้นทาง-api)
9. [Middleware Pipeline](#9-middleware-pipeline)
10. [รายละเอียดโมดูล](#10-รายละเอียดโมดูล)
    - 10.1 โมดูล Users
    - 10.2 โมดูล Auth
    - 10.3 โมดูล Items
    - 10.4 โมดูล IoT
    - 10.5 โมดูล Alarm
    - 10.6 โมดูล Kafka
    - 10.7 โมดูล WebSocket
    - 10.8 โมดูล MQTT
    - 10.9 Distributor / Processor / Worker
    - 10.10 โมดูล Template
11. [มาตรการรักษาความปลอดภัย](#11-มาตรการรักษาความปลอดภัย)
12. [การปรับปรุงประสิทธิภาพ](#12-การปรับปรุงประสิทธิภาพ)
13. [การตรวจสอบและ Metrics](#13-การตรวจสอบและ-metrics)
14. [การติดตั้งและการใช้งาน](#14-การติดตั้งและการใช้งาน)
15. [สถิติโครงการ](#15-สถิติโครงการ)

---

## 1. ภาพรวมโครงการ

**ICMON** (Industrial IoT Control & Monitoring) เป็นระบบแบ็คเอนด์ที่พัฒนาด้วยภาษา Go มีหน้าที่หลัก:

- ตรวจสอบอุปกรณ์และเก็บข้อมูล IoT ผ่าน **MQTT**, **InfluxDB** (ฐานข้อมูลอนุกรมเวลา) และ **PostgreSQL**
- สตรีมข้อมูลแบบเรียลไทม์ผ่าน **WebSocket**
- ประมวลผลแบบ event-driven ผ่าน **Apache Kafka**
- ประมวลผลงานเบื้องหลังผ่าน **Asynq** (Redis queue)
- จัดการผู้ใช้แบบ **RBAC** (5 บทบาท: SUPERADMIN, ADMIN, EDITOR, MONITOR, USER)
- ยืนยันตัวตนด้วย JWT access + refresh tokens (RS256)
- ยืนยันอีเมล, รีเซ็ตรหัสผ่านผ่าน SMTP พร้อมเทมเพลต HTML
- การตรวจสอบด้วย Prometheus และ custom JSON metrics
- การจำกัดอัตราการเรียกใช้ต่อ IP

**กรณีการใช้งานหลัก:**
- เก็บข้อมูลเซนเซอร์ IoT ในโรงงานอุตสาหกรรมและแจ้งเตือน
- ระบบควบคุมคุณภาพอากาศ (mods, periods, warnings)
- การตั้งเวลาอุปกรณ์และการแจ้งเตือน
- แสดงผลแดชบอร์ดแบบเรียลไทม์ผ่าน WebSocket
- บันทึก Audit log ของกิจกรรมทั้งหมดในระบบ

---

## 2. สถาปัตยกรรมระบบ

### 2.1 Clean Architecture (แนวคิด DDD)

โปรเจกต์ใช้ **Clean Architecture** มี 3 ชั้นหลัก + ชั้น delivery:

```
┌─────────────────────────────────────────────────────────────┐
│                    Delivery Layer                            │
│  (HTTP handlers, WebSocket, Kafka consumers, CLI commands)  │
├─────────────────────────────────────────────────────────────┤
│                    Usecase Layer                             │
│  (Business logic, orchestration, validation)                │
├─────────────────────────────────────────────────────────────┤
│                    Repository Layer                          │
│  (Data access — PostgreSQL, Redis, InfluxDB)                │
├─────────────────────────────────────────────────────────────┤
│                    Model Layer                               │
│  (Domain entities, database models, DTOs)                   │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 Dependency Injection

การไหลของ dependency จะเข้าสู่ด้านใน — delivery layer ขึ้นกับ usecase, usecase ขึ้นกับ repository interface โดย concrete implementation จะถูก inject ที่ composition root (การตั้งค่า server)

### 2.3 การไหลของการสื่อสาร

```
┌──────────┐     ┌──────────┐     ┌───────────┐
│  Client   │────▶│  Router  │────▶│ Middleware  │
└──────────┘     └──────────┘     └───────────┘
                                         │
                                         ▼
                                  ┌───────────┐
                                  │  Handler   │
                                  └───────────┘
                                         │
                                         ▼
                                  ┌───────────┐
                                  │  Usecase   │
                                  └───────────┘
                                    │       │
                                    ▼       ▼
                              ┌────────┐ ┌────────┐
                              │ Repo   │ │  Infra │
                              │ (PG)   │ │ (Redis,│
                              │        │ │ Influx,│
                              │        │ │ MQTT)  │
                              └────────┘ └────────┘
```

### 2.4 การเชื่อมต่อระบบภายนอก

```
┌────────────────────────────────────────────────────────────────────┐
│                         icmongolang                                 │
│                                                                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────────────────┐  │
│  │PostgreSQL │  │  Redis   │  │ InfluxDB │  │       MQTT        │  │
│  │ (หลัก)    │  │(sessions,│  │(sensor   │  │(device telemetry) │  │
│  │           │  │  queue)  │  │  data)   │  │                   │  │
│  └──────────┘  └──────────┘  └──────────┘  └───────────────────┘  │
│                                                                     │
│  ┌──────────┐  ┌──────────┐  ┌─────────────────────────────────┐  │
│  │  Kafka   │  │  SMTP    │  │      Prometheus                  │  │
│  │(orders)  │  │(email)   │  │   (/metrics, /apimetric)         │  │
│  └──────────┘  └──────────┘  └─────────────────────────────────┘  │
└────────────────────────────────────────────────────────────────────┘
```

---

## 3. เทคโนโลยีและ dependencies

### 3.1 แกนหลัก

| ไลบรารี | เวอร์ชัน | วัตถุประสงค์ |
|----------|---------|-------------|
| Go | 1.25.0 | ภาษา |
| `github.com/go-chi/chi/v5` | v5.2.5 | HTTP router |
| `github.com/gorilla/mux` | v1.8.1 | HTTP router (Kafka/WS services) |
| `github.com/spf13/cobra` | v1.10.2 | CLI framework |
| `github.com/spf13/viper` | v1.21.0 | การจัดการ config |
| `github.com/swaggo/swag` | v1.16.6 | สร้าง Swagger docs |
| `github.com/swaggo/http-swagger` | v1.3.4 | Swagger UI |
| `go.uber.org/zap` | v1.27.1 | รองรับ structured logging |

### 3.2 ฐานข้อมูลและ Cache

| ไลบรารี | เวอร์ชัน | วัตถุประสงค์ |
|----------|---------|-------------|
| `gorm.io/gorm` | v1.25.8 | ORM |
| `gorm.io/driver/postgres` | v1.5.7 | PostgreSQL driver |
| `github.com/redis/go-redis/v9` | v9.20.0 | Redis client |
| `github.com/influxdata/influxdb-client-go/v2` | v2.14.0 | InfluxDB client |

### 3.3 Messaging และ Streaming

| ไลบรารี | เวอร์ชัน | วัตถุประสงค์ |
|----------|---------|-------------|
| `github.com/IBM/sarama` | v1.42.1 | Kafka client |
| `github.com/eclipse/paho.mqtt.golang` | v1.5.1 | MQTT client |
| `github.com/gorilla/websocket` | v1.5.3 | WebSocket |
| `github.com/hibiken/asynq` | v0.26.0 | Task queue |

### 3.4 ความปลอดภัยและ Validation

| ไลบรารี | เวอร์ชัน | วัตถุประสงค์ |
|----------|---------|-------------|
| `github.com/golang-jwt/jwt/v4` | v4.5.2 | JWT (RS256) |
| `github.com/go-playground/validator/v10` | v10.30.2 | ตรวจสอบ input |
| `golang.org/x/crypto` | v0.53.0 | bcrypt การเข้ารหัสรหัสผ่าน |
| `golang.org/x/time` | v0.14.0 | Rate limiting |

### 3.5 อีเมล

| ไลบรารี | เวอร์ชัน | วัตถุประสงค์ |
|----------|---------|-------------|
| `gopkg.in/gomail.v2` | - | ส่งอีเมลผ่าน SMTP |
| `github.com/matcornic/hermes/v2` | v2.1.0 | เทมเพลตอีเมล HTML |

### 3.6 การทดสอบ

| ไลบรารี | เวอร์ชัน | วัตถุประสงค์ |
|----------|---------|-------------|
| `github.com/stretchr/testify` | v1.11.1 | การ assert ในการทดสอบ |
| `github.com/alicebob/miniredis/v2` | - | Redis จำลอง |
| `github.com/prometheus/client_golang` | v1.23.2 | Prometheus metrics |

---

## 4. การตั้งค่าระบบ

### 4.1 โครงสร้าง Config

การตั้งค่าทั้งหมดจัดการผ่าน **Viper** รองรับทั้งไฟล์ YAML และ environment variables

```go
type Config struct {
    Server         ServerConfig
    Postgres       PostgresConfig
    Redis          RedisConfig
    Jwt            JwtConfig
    FirstSuperUser FirstSuperUserConfig
    Logger         Logger
    SmtpEmail      SmtpEmailConfig
    Email          EmailConfig
    TaskRedis      TaskRedisConfig
    MQTT           MQTTConfig
    InfluxDB       InfluxDBConfig
    Kafka          KafkaConfig
}
```

### 4.2 พารามิเตอร์สำคัญ

| พารามิเตอร์ | ค่าเริ่มต้น | คำอธิบาย |
|-------------|---------|-------------|
| `server.port` | 5000 | พอร์ต HTTP server |
| `server.mode` | Development | โหมด (Development/Production) |
| `server.timezone` | Asia/Bangkok | เขตเวลาของระบบ |
| `server.baseUrl` | http://localhost:5000/api | Base URL ของ API |
| `postgres.dbname` | icmon | ชื่อฐานข้อมูล |
| `redis.addr` | localhost:6379 | ที่อยู่ Redis |
| `influxdb.url` | http://localhost:9087 | URL InfluxDB |
| `influxdb.bucket` | icmon | Bucket InfluxDB |
| `kafka.brokers` | [localhost:9092] | รายการ Kafka broker |
| `mqtt.broker` | tcp://localhost:1891 | MQTT broker |
| `jwt.accessTokenExpireDuration` | 30m | อายุ access token |
| `jwt.refreshTokenExpireDuration` | 168h | อายุ refresh token |

### 4.3 ไฟล์ Config

- `config/config.default.yml` — config ฐานสำหรับทุก environment
- `config/config.dev.yml` — ค่าเฉพาะ development

Environment variables จะ map ไปยัง config keys โดยอัตโนมัติผ่าน `BindEnvs()`

---

## 5. จุดเริ่มต้นโปรแกรมและคำสั่ง CLI

### 5.1 จุดเริ่มต้นโปรแกรม

โปรเจกต์มี **4 จุดเริ่มต้น** สำหรับวัตถุประสงค์ต่างกัน:

| จุดเริ่มต้น | ไฟล์ | วัตถุประสงค์ |
|-------------|------|-------------|
| **Main API** | `cmd/api/main.go` | HTTP server หลัก (พอร์ต 5000) |
| **Kafka Service** | `cmd/kafka/main.go` | บริการ Kafka + WebSocket (พอร์ต 5051) |
| **WebSocket Service** | `cmd/websocket/main.go` | บริการ WebSocket + Redis Queue (พอร์ต 8080) |
| **CLI Root** | `cmd/root.go` | คำสั่ง Cobra (serve, migrate, worker, initdata) |

### 5.2 คำสั่ง CLI

| คำสั่ง | ไฟล์ | คำอธิบาย |
|---------|------|-------------|
| `serve` | `cmd/serve.go` | เริ่ม HTTP server พร้อมบริการทั้งหมด (Postgres, Redis, MQTT, InfluxDB, auto-migration) |
| `migrate` | `cmd/migrate.go` | รัน GORM AutoMigrate (~180 models) พร้อม skip-list |
| `worker` | `cmd/worker.go` | เริ่ม Asynq background task processor |
| `initdata` | `cmd/initdata.go` | สร้าง superuser เริ่มต้นจาก config |
| `kafka` | `cmd/kafka/kafka.go` | เริ่ม Kafka producer + consumer + WebSocket |

### 5.3 Global Flags

- `--config` — พาธไปยังไฟล์ config (ค่าเริ่มต้น: `$HOME/.go-base.yaml`)

---

## 6. โครงสร้างไดเรกทอรี

```
icmongolang/
├── main.go                          # จุดเข้าโปรแกรม (เรียก cmd.Execute())
│
├── cmd/                             # จุดเริ่มต้นและคำสั่ง CLI (9 ไฟล์)
│   ├── api/main.go                  # จุดเข้า API server
│   ├── kafka/main.go                # จุดเข้า Kafka server
│   ├── websocket/main.go            # จุดเข้า WebSocket server
│   ├── root.go                      # คำสั่ง root ของ Cobra
│   ├── serve.go                     # คำสั่ง `serve`
│   ├── migrate.go                   # คำสั่ง `migrate`
│   ├── worker.go                    # คำสั่ง `worker`
│   ├── initdata.go                  # คำสั่ง `initdata`
│   └── kafka/kafka.go               # คำสั่ง `kafka`
│
├── config/                          # ไฟล์ตั้งค่าระบบ
│   ├── config.go                    # โครงสร้าง Config และตัวโหลด
│   ├── config.default.yml           # ค่าตั้งต้น
│   └── config.dev.yml               # ค่าสำหรับ development
│
├── internal/                        # แกนกลางของแอปพลิเคชัน (159 ไฟล์ Go)
│   ├── alarm/                       # โมดูลจัดการ Alarm
│   ├── auth/                        # โมดูลยืนยันตัวตน
│   ├── delivery/rest/               # REST router และ handlers
│   ├── distributor/                 # กระจายงาน Asynq
│   ├── influxdb/                    # ตัวปรับต่อ InfluxDB
│   ├── iot/                         # โมดูลข้อมูล IoT (ใหญ่ที่สุด)
│   ├── items/                       # โมดูล CRUD สินค้า
│   ├── kafka/                       # โมดูล Kafka order
│   ├── middleware/                   # HTTP middleware ทั้งหมด
│   ├── models/                      # Model ฐานข้อมูลทั้งหมด
│   ├── mqtt/                        # โมดูล MQTT
│   ├── processor/                   # ประมวลผลงาน
│   ├── queue/                       # ส่วนที่เป็นนามธรรมของ queue
│   ├── repository/                  # Interface ของ repository
│   ├── server/                      # การตั้งค่า server และ lifecycle
│   ├── template/                    # โมดูลเทมเพลต
│   ├── usecase/                     # Interface ของ usecase
│   ├── users/                       # โมดูลจัดการผู้ใช้
│   ├── websocket/                   # โมดูล WebSocket
│   └── worker/                      # Asynq worker
│
├── pkg/                             | ไลบรารี共用 (30 ไฟล์ Go)
│   ├── cryptpass/                   # การเข้ารหัสรหัสผ่าน bcrypt
│   ├── db/postgres/                 # การเชื่อมต่อ PostgreSQL
│   ├── db/redis/                    # การเชื่อมต่อ Redis
│   ├── helpers/                     | helper การจัดรูปแบบและ IoT alarm
│   ├── influxdb/                    # InfluxDB client wrapper
│   ├── jwt/                         # การจัดการ JWT (RS256)
│   ├── kafka/                       # Kafka producer/consumer
│   ├── logger/                      # Zap structured logger
│   ├── mqtt/                        # MQTT client
│   ├── responses/                   | มาตรฐานการตอบสนอง API
│   ├── utils/                       # Validator และ form utilities
│   └── websocket/                   # WebSocket hub และ client
│
├── db/                              # ไฟล์ SQL และ database dump
│   ├── icmon.sql                    # database dump เต็ม (20K+ บรรทัด)
│   ├── public.sql                   # public schema
│   └── sd_iot_device.sql            # ตารางอุปกรณ์ IoT
│
├── migrations/                      # การย้ายฐานข้อมูล Manual SQL
│   ├── icmon.sql, public.sql, db.sql
│   ├── sd_user_access_menu.sql
│   ├── 20250619_websocket_tables.sql
│   └── 20250619_kafka_tables.sql
│
├── docs/                            # เอกสารประกอบ
│   ├── swagger.yaml                 # สเปค Swagger API (2,901 บรรทัด)
│   ├── docs.md                      # แผนภาพสถาปัตยกรรม
│   └── report/                      # รายงานโครงการ
│       ├── english.md               # ฉบับภาษาอังกฤษ
│       └── thai.md                  # ฉบับภาษาไทย
│
├── Makefile                         # ระบบ build อัตโนมัติ
└── README.md, README_RUN.md         # คู่มือโครงการ
```

---

## 7. โครงสร้างฐานข้อมูล

### 7.1 PostgreSQL

ฐานข้อมูล: `icmon`  
Database dump เต็ม: `db/icmon.sql` (20,052 บรรทัด, Navicat Premium dump)

**ตารางแบ่งตามหมวดหมู่ (จาก 180+ models):**

| หมวดหมู่ | ตาราง | คำอธิบาย |
|----------|--------|-------------|
| **หลัก** | `item`, `migrations` | ตารางอ้างอิงพื้นฐาน |
| **ผู้ใช้** | `users`, `sd_users`, `sd_user_roles`, `sd_user_roles_access`, `sd_user_roles_permision`, `sd_user_access_menu`, `sd_user_log`, `sd_user_log_type`, `sd_user_file` | จัดการผู้ใช้และ RBAC |
| **กิจกรรม** | `activity_log`, `sd_activity_log`, `sd_activity_type_log`, `sd_audit_log`, `sd_module_log`, `sd_device_log`, `sd_mqtt_log` | บันทึกการใช้งานทั้งหมด |
| **อุปกรณ์** | `devices`, `sd_iot_device`, `sd_iot_device_type`, `sd_iot_device_action`, `sd_iot_device_action_log`, `sd_iot_device_action_user`, `sd_iot_device_alarm_action`, `sd_device_alert`, `sd_device_config`, `sd_device_status`, `sd_device_category`, `sd_device_group`, `sd_device_member`, `sd_device_notification_config`, `sd_device_schedule`, `sd_device_status_history` | จัดการอุปกรณ์ |
| **IoT** | `iot_data`, `sd_iot_data`, `sd_iot_group`, `sd_iot_host`, `sd_iot_location`, `sd_iot_sensor`, `sd_iot_type`, `sd_iot_setting`, `sd_iot_api`, `sd_iot_token`, `sd_iot_email`, `sd_iot_line`, `sd_iot_sms`, `sd_iot_telegram`, `sd_iot_nodered`, `sd_iot_mqtt`, `sd_iot_influxdb`, `sd_iot_schedule`, `sd_iot_schedule_device` | การตั้งค่าและข้อมูล IoT |
| **Alarm** | `sd_iot_alarm_device`, `sd_iot_alarm_device_event`, `sd_alarm_process_log`, `sd_alarm_process_log_email`, `sd_alarm_process_log_line`, `sd_alarm_process_log_mqtt`, `sd_alarm_process_log_sms`, `sd_alarm_process_log_telegram`, `sd_alarm_process_log_temp` | จัดการแจ้งเตือน |
| **ควบคุมอากาศ** | `air_controls`, `air_mods`, `air_periods`, `air_setting_warnings`, `air_warnings` (และรุ่นที่มี `sd_` พร้อม device maps) | ควบคุมคุณภาพอากาศ |
| **การแจ้งเตือน** | `noti_notifications`, `notification_logs`, `notification_groups`, `notification_devices`, `notification_types`, `sd_notification_channel`, `sd_notification_condition`, `sd_notification_log`, `sd_notification_type`, `sd_channel_template`, `sd_group_notification_config` | ระบบแจ้งเตือน |
| **Dashboard** | `sd_dashboard_config` (alias: `tnb`) | การตั้งค่า dashboard |
| **API Keys** | `sd_api_keys` | จัดการคีย์ API |
| **ระบบ** | `sd_system_settings` | การตั้งค่าระบบ |
| **รายงาน** | `sd_report_data`, `sd_schedule_process_log` | รายงานข้อมูล |
| **ตารางงาน** | `sd_device_schedule`, `sd_iot_schedule`, `sd_iot_schedule_device` | การตั้งเวลา |
| **คำสั่ง** | `command_log` | บันทึกคำสั่ง |
| **Kafka** | `orders` | เหตุการณ์ order ของ Kafka |

### 7.2 Enums

```go
type UserRoleEnum string
const (
    SUPERADMIN UserRoleEnum = "SUPERADMIN"
    ADMIN      UserRoleEnum = "ADMIN"
    EDITOR     UserRoleEnum = "EDITOR"
    MONITOR    UserRoleEnum = "MONITOR"
    USER       UserRoleEnum = "USER"
)

type UserUsertypeEnum string
// therapist, supervisor, superadmin, system, admin, support, enduser
```

### 7.3 ส่วนขยายฐานข้อมูลที่สำคัญ

- `uuid-ossp` — สร้าง UUID
- `pgcrypto` — ฟังก์ชันการเข้ารหัส

### 7.4 กลยุทธ์การย้ายฐานข้อมูล

**สองแนวทาง:**
1. **GORM AutoMigrate** — ใช้เมื่อเริ่มระบบหรือผ่านคำสั่ง `migrate` สำหรับ 180+ GORM models มี skip-list สำหรับตารางที่มี composite primary key หรือ schema ซับซ้อน
2. **Manual SQL** — ไฟล์ใน `migrations/` สำหรับการเปลี่ยนแปลง schema ที่ซับซ้อน

---

## 8. เส้นทาง API

### 8.1 Base Path

เส้นทาง API ทั้งหมดอยู่ภายใต้ `/api` (ยกเว้น Kafka และ WebSocket services ที่ใช้ `/`)

### 8.2 เส้นทางสาธารณะ

| Method | Path | คำอธิบาย |
|--------|------|-------------|
| GET | `/metrics` | Prometheus metrics |
| GET | `/apimetric` | Custom JSON metrics |
| GET | `/health` | ตรวจสอบสถานะ |
| GET | `/api/ping` | Ping/pong |
| POST | `/api/auth/login` | เข้าสู่ระบบ |
| POST | `/api/auth/signin` | เข้าสู่ระบบ |
| GET | `/api/auth/publickey` | ดึง public key ของ JWT |
| GET | `/api/auth/verifyemail` | ยืนยันอีเมล |
| POST | `/api/auth/forgotpassword` | ขอรีเซ็ตรหัสผ่าน |
| PATCH | `/api/auth/resetpassword` | รีเซ็ตรหัสผ่าน |
| POST | `/api/register` | ลงทะเบียนผู้ใช้ใหม่ |
| POST | `/api/users` | สร้างผู้ใช้ |
| POST | `/api/signin` | เข้าสู่ระบบด้วยอีเมล |
| POST | `/api/login` | เข้าสู่ระบบด้วยชื่อผู้ใช้ |

### 8.3 เส้นทางที่ต้องยืนยันตัวตน (Authenticated)

| Method | Path | คำอธิบาย |
|--------|------|-------------|
| GET | `/api/auth/refresh` | ขอ access token ใหม่ |
| GET | `/api/auth/logout` | ออกจากระบบเซสชันปัจจุบัน |
| GET | `/api/auth/logoutall` | ออกจากระบบทุกเซสชัน |
| GET | `/api/user/me` | ดูโปรไฟล์ตัวเอง |
| PUT | `/api/user/me` | แก้ไขโปรไฟล์ตัวเอง |
| PATCH | `/api/user/me/updatepass` | เปลี่ยนรหัสผ่านตัวเอง |
| GET | `/api/item/` | ดูรายการ items |
| POST | `/api/item/` | สร้าง item |
| GET | `/api/item/{id}` | ดู item ตาม ID |
| DELETE | `/api/item/{id}` | ลบ item |
| PUT | `/api/item/{id}` | แก้ไข item |

### 8.4 เส้นทางที่ต้องเป็น SuperUser/Admin

| Method | Path | คำอธิบาย |
|--------|------|-------------|
| GET | `/api/user/` | ดูรายชื่อผู้ใช้ทั้งหมด |
| POST | `/api/user/` | สร้างผู้ใช้ (admin) |
| GET | `/api/user/{id}` | ดูผู้ใช้ตาม ID |
| PUT | `/api/user/{id}` | แก้ไขผู้ใช้ |
| DELETE | `/api/user/{id}` | ลบผู้ใช้ |
| PATCH | `/api/user/{id}/role` | เปลี่ยนบทบาทผู้ใช้ |
| PATCH | `/api/user/{id}/updatepass` | เปลี่ยนรหัสผ่านผู้ใช้ใดๆ |
| GET | `/api/user/{id}/logoutall` | บังคับออกจากระบบทุกเซสชัน |

### 8.5 เส้นทาง Kafka Service (พอร์ต 5051)

| Method | Path | คำอธิบาย |
|--------|------|-------------|
| POST | `/orders` | สร้าง order (ส่ง Kafka message) |
| GET | `/ws` | WebSocket endpoint |
| GET | `/health` | ตรวจสอบสถานะ |

---

## 9. Middleware Pipeline

### 9.1 ลำดับการทำงาน

คำขอจะผ่าน middleware chain ดังนี้ (กำหนดใน `internal/delivery/rest/router.go`):

```
1. panic recovery (Recoverer)
2. กำหนด request ID (RequestID)
3. ระบุ client IP จริง (RealIP)
4. กำหนด timeout การประมวลผล (Timeout)
5. เพิ่ม CORS headers (CORS)
6. เพิ่ม security headers (CSP, HSTS, XSS, NoSniff, FrameDeny, ฯลฯ)
7. จำกัดอัตราการเรียกใช้ต่อ IP (10 req/s, burst 20)
8. บันทึก log request/response
9. ตรวจสอบ Prometheus metrics
10. route handler (JWT auth → handler → response)
```

### 9.2 ส่วนประกอบของ Middleware

| Middleware | ไฟล์ | คำอธิบาย |
|-----------|------|-------------|
| **CORS** | `middleware/cors.go` | CORS ตาม environment (development: localhost; production: BaseUrl) |
| **Security** | `middleware/security.go` | Security headers: CSP, HSTS, XSS, NoSniff, FrameDeny, Referrer, Permissions-Policy, Cache-Control |
| **Rate Limit** | `middleware/rate_limit.go` | จำกัดอัตราการเรียกใช้ต่อ IP (10 req/s, burst 20, cleanup ทุก 30s) |
| **Logging** | `middleware/logging.go` | บันทึก method, path, status code (ไม่มี request body) |
| **Monitoring** | `middleware/monitoring.go` | Prometheus `http_requests_total`, `http_request_duration_seconds` + `/apimetric` |
| **JWT Auth** | `middleware/jwtauth.go` | ตรวจสอบ access/refresh tokens (RS256), context ผู้ใช้ปัจจุบัน, ตรวจสอบ super user, ตรวจสอบ active user |
| **Prometheus** | `middleware/prometheus.go` | System metrics: goroutines, memory, GC, CPU, network, Redis connection stats |

---

## 10. รายละเอียดโมดูล

### 10.1 โมดูล Users (`internal/modules/users/`)

**ชั้นต่างๆ:**
- `models/user.go` — เอนทิตี้ผู้ใช้
- `repository/pg_repository.go` — CRUD PostgreSQL (GORM), จัดการ session
- `usecase/usecase.go` — ลงทะเบียน, เข้าสู่ระบบ, จัดการโปรไฟล์, เปลี่ยนรหัสผ่าน, ลืม/รีเซ็ตรหัสผ่าน, ยืนยันอีเมล, จัดการบทบาท
- `delivery/http/handlers.go` — HTTP handlers สำหรับ endpoints ผู้ใช้ทั้งหมด
- `distributor/distributor.go` — กระจายงาน Asynq (งานอีเมล)
- `processor/processor.go` — ประมวลผลงาน Asynq (ส่งอีเมล)
- `worker/worker.go` — นิยามงาน worker

**ความสามารถหลัก:**
- CRUD เต็มรูปแบบพร้อม RBAC
- ขั้นตอนยืนยันอีเมลด้วย token
- ขั้นตอนลืม/รีเซ็ตรหัสผ่าน
- จัดการ session (refresh tokens ใน Redis)
- ส่งอีเมลแบบ async ผ่าน Asynq

### 10.2 โมดูล Auth (`internal/modules/auth/`)

**ชั้นต่างๆ:**
- `delivery/http/handlers.go` — เข้าสู่ระบบ, refresh token, ออกจากระบบ, public key, ยืนยันอีเมล
- `route.go` — นิยามเส้นทาง

**ความสามารถหลัก:**
- JWT access + refresh token (กุญแจ RSA RS256)
- Refresh tokens เก็บใน Redis
- Public key endpoint สำหรับการตรวจสอบภายนอก
- Token rotation เมื่อ refresh

### 10.3 โมดูล Items (`internal/modules/items/`)

**ชั้นต่างๆ:**
- `repository/` — CRUD repository มาตรฐาน
- `usecase/` — Business logic
- `delivery/http/` — HTTP handlers

โมดูล CRUD อย่างง่าย — ใช้เป็นรูปแบบอ้างอิงสำหรับการจัดการทรัพยากรพื้นฐาน

### 10.4 โมดูล IoT (`internal/modules/iot/`) — โมดูลที่ใหญ่ที่สุด

**ชั้นต่างๆ:**
- `models/` — โมเดลโดเมนคู่ขนาน (Device, IotData, AirControl, AirMod, AirPeriod, AirWarning, ฯลฯ)
- `repository/iot_data_repo.go` — คำสั่งค้นหาข้อมูล IoT พร้อมแบ่งหน้า, ตรวจสอบ device
- `repository/device_repo.go` — CRUD อุปกรณ์พร้อมการกรองแบบไดนามิก
- `usecase/usecase.go` — Business logic ที่ซับซ้อน (สถานะอุปกรณ์, alarm, ควบคุมอากาศ, รวมข้อมูล, dashboard)
- `presenter/` — การจัดรูปแบบ response

**ความสามารถหลัก:**
- จัดการอุปกรณ์ด้วยการกรองคอลัมน์แบบไดนามิก (ป้องกัน SQL injection ด้วย whitelist)
- ดึงข้อมูลเซนเซอร์ IoT พร้อมแบ่งหน้าที่ DB level
- ระบบควบคุมอากาศ (mods, periods, warnings พร้อม device mapping)
- การประเมินกฏ alarm และสร้างการแจ้งเตือน
- ติดตามประวัติสถานะอุปกรณ์
- คำสั่งค้นหาแบบรวมสำหรับ dashboard

### 10.5 โมดูล Alarm (`internal/modules/alarm/`)

**ชั้นต่างๆ:**
- `repository/alarm_log_repo.go` — คำสั่งค้นหา alarm log พร้อมการกรองแบบไดนามิก (ป้องกัน SQL injection)
- `usecase/` — Logic การประมวลผล alarm
- `delivery/` — HTTP handlers

**ความสามารถหลัก:**
- ดึงประวัติ alarm พร้อมกรองหลายช่องทาง
- รองรับช่องทางแจ้งเตือน email, line, SMS, MQTT, Telegram
- ประวัติ alarm แบบแบ่งหน้าและกรองได้

### 10.6 โมดูล Kafka (`internal/modules/kafka/`)

**ชั้นต่างๆ:**
- `models/` — โมเดลเหตุการณ์ Order
- `repository/` — Kafka repository
- `usecase/` — Logic การประมวลผล Order
- `delivery/` — HTTP + WebSocket handlers

**ความสามารถหลัก:**
- สร้าง/บริโภคเหตุการณ์ Order ผ่าน Apache Kafka (sarama)
- สตรีมสถานะ Order แบบเรียลไทม์ผ่าน WebSocket
- บริการอิสระ (รันบนพอร์ต 5051, ใช้ gorilla/mux)

### 10.7 โมดูล WebSocket (`internal/modules/websocket/`)

**ชั้นต่างๆ:**
- `models/` — ประเภทข้อความ WebSocket
- `repository/` — WebSocket session repository
- `usecase/` — Logic จัดการข้อความ
- `delivery/` — จัดการเชื่อมต่อ WebSocket

**ความสามารถหลัก:**
- การสื่อสารสองทางแบบเรียลไทม์
- จัดการ session และกระจายข้อความ
- รันแบบอิสระ (พอร์ต 8080) หรือแบบฝัง

### 10.8 โมดูล MQTT (`internal/modules/mqtt/`)

**ชั้นต่างๆ:**
- `delivery/` — MQTT message handlers
- `presenter/` — จัดรูปแบบข้อมูล
- `usecase/` — ประมวลผลข้อความ MQTT

**ความสามารถหลัก:**
- สื่อสารกับอุปกรณ์ IoT ผ่านโปรโตคอล MQTT
- Pipeline การสมัครรับและประมวลผลข้อความ
- การเชื่อมต่อใหม่แบบ exponential backoff

### 10.9 Distributor / Processor / Worker

**Distributor** (`internal/distributor/`):
- อินเทอร์เฟซกระจายงาน Asynq
- สร้างงานสำหรับยืนยันอีเมล, รีเซ็ตรหัสผ่าน

**Processor** (`internal/processor/`):
- ประมวลผลงาน Asynq (ส่งอีเมล)

**Worker** (`internal/worker/`):
- การตั้งค่า Asynq worker server
- Mux สำหรับการ route งาน
- การปิดระบบอย่างนุ่มนวล

### 10.10 โมดูล Template (`internal/template/`)

โมดูลเทมเพลตอ้างอิงสำหรับสร้างโมดูลโดเมนใหม่ มีโครงสร้างสำเร็จรูป:
- Delivery (HTTP handlers, routes)
- Repository (PostgreSQL + Redis)
- Usecase (business logic)
- Models (domain entities)
- Presenter (การจัดรูปแบบ response)

---

## 11. มาตรการรักษาความปลอดภัย

### 11.1 การยืนยันตัวตนและการอนุญาต

| มาตรการ | รายละเอียดการใช้งาน |
|---------|---------------|
| **JWT RS256** | access + refresh tokens แบบ RSA หมดอายุได้ตาม config |
| **Bcrypt passwords** | รหัสผ่านทั้งหมดเข้ารหัสด้วย bcrypt (`pkg/cryptpass`) |
| **Refresh token rotation** | token เก่าถูกยกเลิกเมื่อ refresh |
| **Redis-backed sessions** | Refresh tokens เก็บใน Redis พร้อม TTL |
| **SuperUser middleware** | ควบคุมการเข้าถึง endpoints สำหรับผู้ดูแลตามบทบาท |
| **Active user check** | ผู้ใช้ที่ถูกปิดใช้งานไม่สามารถเรียก API ได้ |

### 11.2 การป้องกัน SQL Injection

| จุดที่แก้ไข | ตำแหน่ง |
|-----|----------|
| whitelist ชื่อคอลัมน์ | `internal/modules/alarm/repository/alarm_log_repo.go` |
| whitelist ชื่อคอลัมน์ + ORDER BY | `internal/modules/iot/repository/device_repo.go` |

Dynamic `.Where()` clauses และ ORDER BY parameters ถูกตรวจสอบกับรายการฟิลด์ที่อนุญาตก่อนส่งไปยัง GORM

### 11.3 การป้องกัน CORS

- **Development**: อนุญาต origins localhost (มีและไม่มีพอร์ต)
- **Production**: จำกัดเฉพาะ `config.Server.BaseUrl` เท่านั้น
- ไม่มี wildcard (`*`) origin เด็ดขาด

### 11.4 การจำกัดอัตราการเรียกใช้

- **อัตรา**: 10 คำขอต่อวินาทีต่อ IP (จากเดิม 2000)
- **Burst**: 20 คำขอ
- **Cleanup**: IP ที่ไม่ได้ใช้จะถูกลบทุก 30 วินาที
- **แหล่งที่มา IP**: ใช้เฉพาะ `RemoteAddr` (ป้องกันการปลอมแปลง `X-Forwarded-For`)

### 11.5 Security Headers (CSP, HSTS, ฯลฯ)

ทุก response มี headers:
- `Content-Security-Policy`
- `Strict-Transport-Security`
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 0` (แนวทางสมัยใหม่)
- `Referrer-Policy`
- `Permissions-Policy`
- `Cache-Control`

### 11.6 การป้องกัน PII

| สิ่งที่แก้ไข | ตำแหน่ง |
|----------------|-------|
| ลบ request body ออกจาก logs | `internal/middleware/logging.go` |
| ลบ email/password ออกจาก auth logs | `internal/modules/auth/delivery/http/handlers.go` |
| ลบ PII ออกจาก user usecase logs | `internal/modules/users/usecase/usecase.go` |

### 11.7 Timeouts และ Context

| คอมโพเนนต์ | Timeout | ตำแหน่ง |
|-----------|---------|-------|
| Process timeout | ตั้งค่าได้ (ค่าเริ่มต้น 60s) | `internal/delivery/rest/router.go` |
| InfluxDB queries | 30 วินาที | `pkg/influxdb/client.go` |
| Redis operations | 5 วินาทีต่อ operation | `pkg/db/redis/redis_conn.go` |
| Server shutdown | 5 วินาที (graceful) | `internal/server/server.go` |

---

## 12. การปรับปรุงประสิทธิภาพ

### 12.1 การแบ่งหน้าฐานข้อมูล

| ก่อน | หลัง |
|--------|-------|
| In-memory: ดึง 10,000 records → ตัดใน Go | DB-level: `LIMIT ? OFFSET ?` |
| สิ้นเปลืองหน่วยความจำและ CPU ทุกครั้งที่แบ่งหน้า | การแบ่งหน้าที่มีประสิทธิภาพใช้ index |

ดำเนินการใน `internal/modules/iot/usecase/usecase.go` ผ่าน `CountByDeviceID()` + `GetByDeviceID(limit, offset)`

### 12.2 การปรับแต่ง `structToMap`

| ก่อน | หลัง |
|--------|-------|
| `json.Marshal` แล้ว `json.Unmarshal` เป็น map[string]interface{} | การแปลงด้วย reflection โดยตรง |
| จองหน่วยความจำสองครั้งต่อการเรียก | รอบเดียว ไม่ต้องมี JSON กลาง |

ดำเนินการใน `internal/modules/iot/usecase/usecase.go`

### 12.3 ลบความหน่วงเทียม

ลบ `time.Sleep(2 * time.Second)` ออกจาก Kafka order processing (`internal/modules/kafka/usecase/order_usecase.go`) — เป็น simulation ใน development ที่ขัดขวาง throughput ใน production

### 12.4 Context Timeouts

- คำสั่งค้นหา InfluxDB ทั้งหมดมี timeout 30 วินาที
- Redis operations ทั้งหมดมี timeout 5 วินาทีต่อ operation

### 12.5 การปรับอัตราการจำกัด

ลดอัตราจาก 2000 req/s เหลือ 10 req/s (ค่าเริ่มต้นที่เหมาะสม) — ป้องกัน DoS โดยไม่ได้ตั้งใจและสอดคล้องกับมาตรฐาน production

---

## 13. การตรวจสอบและ Metrics

### 13.1 Prometheus Metrics (`/metrics`)

| Metric | ประเภท | คำอธิบาย |
|--------|------|-------------|
| `http_requests_total` | Counter | จำนวน HTTP requests ทั้งหมด (labels: method, path, status) |
| `http_request_duration_seconds` | Histogram | การกระจายระยะเวลาของ request |

### 13.2 Custom JSON Metrics (`/apimetric`)

| Metric | คำอธิบาย |
|--------|-------------|
| `go_goroutines` | จำนวน goroutine ปัจจุบัน |
| `go_memory_alloc` | หน่วยความจำที่จัดสรรอยู่ (bytes) |
| `go_gc_count` | จำนวนรอบ GC ที่เสร็จสมบูรณ์ |
| `go_cpu_usage` | เปอร์เซ็นต์การใช้งาน CPU |
| `memory_usage_percent` | เปอร์เซ็นต์การใช้หน่วยความจำ |
| `network_in_bytes_total` | ปริมาณข้อมูลเครือข่ายขาเข้า |
| `network_out_bytes_total` | ปริมาณข้อมูลเครือข่ายขาออก |
| `redis_connections` | การเชื่อมต่อ Redis ที่ใช้งานอยู่ |
| `redis_hits_per_second` | อัตราการ cache hit ของ Redis |
| `redis_misses_per_second` | อัตราการ cache miss ของ Redis |
| `redis_max_connections` | การเชื่อมต่อสูงสุดของ Redis |

### 13.3 การตรวจสอบสถานะ

`GET /api/health` ส่งคืนสถานะเซิร์ฟเวอร์ สามารถใช้โดย load balancers และระบบตรวจสอบ

---

## 14. การติดตั้งและการใช้งาน

### 14.1 ความต้องการระบบ

- Go 1.25+
- PostgreSQL
- Redis
- (ไม่จำเป็น) InfluxDB, MQTT broker, Kafka, SMTP server

### 14.2 เริ่มต้นใช้งาน

```bash
# 1. เริ่ม dependencies (PostgreSQL, Redis)
docker-compose up -d

# 2. รัน migrations
go run cmd/api/main.go migrate

# 3. เริ่มเซิร์ฟเวอร์
go run cmd/api/main.go serve

# 4. สร้าง superuser เริ่มต้น
go run cmd/api/main.go initdata
```

### 14.3 คำสั่ง Make

```bash
make clean     # go clean -cache + go clean -modcache
make tidy      # go mod tidy
make download  # go mod download
make verify    # go mod verify
make migrate   # go run cmd/api/main.go migrate
make swag      # swag init -g cmd/api/main.go
make vendor    # go mod vendor
make test      # go test ./...
make run       # ติดตั้งครบ + air (hot-reload)
```

### 14.4 การรันบริการแยก

```bash
# Kafka + WebSocket service (พอร์ต 5051)
go run cmd/kafka/main.go

# WebSocket + Redis Queue service (พอร์ต 8080)
go run cmd/websocket/main.go

# Background task worker
go run cmd/api/main.go worker
```

### 14.5 ข้อมูลผู้ใช้เริ่มต้น

- **อีเมล:** root@gmail.com
- **รหัสผ่าน:** root_password

### 14.6 Swagger UI

เข้าถึงได้ที่: `http://localhost:5000/swagger/index.html`

---

## 15. สถิติโครงการ

| หมวดหมู่ | จำนวน |
|----------|-------|
| **ไฟล์ Go ทั้งหมด** | 198 |
| **ไฟล์ `cmd/`** | 9 |
| **ไฟล์ `internal/`** | 159 |
| **ไฟล์ `pkg/`** | 30 |
| **Model ฐานข้อมูล** | 180+ |
| **เส้นทาง API** | 30+ |
| **คอมโพเนนต์ Middleware** | 10 |
| **ขนาด Schema ฐานข้อมูล** | 20,052+ บรรทัด SQL |
| **ขนาด Swagger spec** | 2,901 บรรทัด YAML |
| **Dependencies ภายนอก** | 34 รายการ |

---

*สิ้นสุดรายงานทางเทคนิค — ฉบับภาษาไทย*
