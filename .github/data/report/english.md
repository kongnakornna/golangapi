# ICMON Go Project — Complete Technical Report (English)

> **Version:** 1.0.0  
> **Module:** `icmongolang`  
> **Go Version:** 1.25.0  
> **License:** github.com/kongnakornna  
> **Date:** July 2026

---

## Table of Contents

1. [Project Overview](#1-project-overview)
2. [System Architecture](#2-system-architecture)
3. [Tech Stack & Dependencies](#3-tech-stack--dependencies)
4. [Configuration](#4-configuration)
5. [Entry Points & CLI Commands](#5-entry-points--cli-commands)
6. [Directory Structure](#6-directory-structure)
7. [Database Schema](#7-database-schema)
8. [API Routes](#8-api-routes)
9. [Middleware Pipeline](#9-middleware-pipeline)
10. [Module Breakdown](#10-module-breakdown)
    - 10.1 Users Module
    - 10.2 Auth Module
    - 10.3 Items Module
    - 10.4 IoT Module
    - 10.5 Alarm Module
    - 10.6 Kafka Module
    - 10.7 WebSocket Module
    - 10.8 MQTT Module
    - 10.9 Distributor / Processor / Worker
    - 10.10 Template Module
11. [Security Measures](#11-security-measures)
12. [Performance Optimizations](#12-performance-optimizations)
13. [Monitoring & Metrics](#13-monitoring--metrics)
14. [Setup & Deployment](#14-setup--deployment)
15. [Project Statistics](#15-project-statistics)

---

## 1. Project Overview

**ICMON** (Industrial IoT Control & Monitoring) is a backend system written in Go that provides:

- Device monitoring and IoT data collection via **MQTT**, **InfluxDB** (time-series), and **PostgreSQL**
- Real-time data streaming via **WebSocket**
- Event-driven processing via **Apache Kafka**
- Background task processing via **Asynq** (Redis-backed task queue)
- User management with **RBAC** (5 roles: SUPERADMIN, ADMIN, EDITOR, MONITOR, USER)
- JWT-based authentication with access + refresh tokens (RS256)
- Email verification, password reset via SMTP with HTML templates
- Prometheus metrics and custom JSON monitoring
- Per-IP rate limiting

**Primary use cases:**
- Industrial IoT sensor data collection and alerting
- Air control system management (mods, periods, warnings)
- Device scheduling and notification configuration
- Real-time dashboard data via WebSocket
- Audit logging of all system activities

---

## 2. System Architecture

### 2.1 Clean Architecture (DDD-inspired)

The project follows **Clean Architecture** with 3 core layers plus a delivery layer:

```
┌─────────────────────────────────────────────────────────────┐
│                    Delivery Layer                            │
│  (HTTP handlers, WebSocket, Kafka consumers, CLI commands)  │
├─────────────────────────────────────────────────────────────┤
│                    Usecase Layer                             │
│  (Business logic, orchestration, validation)               │
├─────────────────────────────────────────────────────────────┤
│                    Repository Layer                          │
│  (Data access — PostgreSQL, Redis, InfluxDB)                │
├─────────────────────────────────────────────────────────────┤
│                    Model Layer                               │
│  (Domain entities, database models, DTOs)                   │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 Dependency Injection

Dependencies flow inward — the delivery layer depends on usecases, which depend on repository interfaces. Concrete implementations are injected at composition root (server setup).

### 2.3 Communication Flow

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

### 2.4 External Service Integration

```
┌────────────────────────────────────────────────────────────────────┐
│                         icmongolang                                 │
│                                                                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────────────────┐  │
│  │PostgreSQL │  │  Redis   │  │ InfluxDB │  │       MQTT        │  │
│  │ (primary) │  │(sessions,│  │(sensor   │  │(device telemetry) │  │
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

## 3. Tech Stack & Dependencies

### 3.1 Core

| Dependency | Version | Purpose |
|-----------|---------|---------|
| Go | 1.25.0 | Language |
| `github.com/go-chi/chi/v5` | v5.2.5 | HTTP router |
| `github.com/gorilla/mux` | v1.8.1 | HTTP router (Kafka/WS services) |
| `github.com/spf13/cobra` | v1.10.2 | CLI framework |
| `github.com/spf13/viper` | v1.21.0 | Configuration management |
| `github.com/swaggo/swag` | v1.16.6 | Swagger doc generation |
| `github.com/swaggo/http-swagger` | v1.3.4 | Swagger UI |
| `go.uber.org/zap` | v1.27.1 | Structured logging |

### 3.2 Database & Cache

| Dependency | Version | Purpose |
|-----------|---------|---------|
| `gorm.io/gorm` | v1.25.8 | ORM |
| `gorm.io/driver/postgres` | v1.5.7 | PostgreSQL driver |
| `github.com/redis/go-redis/v9` | v9.20.0 | Redis client |
| `github.com/influxdata/influxdb-client-go/v2` | v2.14.0 | InfluxDB client |

### 3.3 Messaging & Streaming

| Dependency | Version | Purpose |
|-----------|---------|---------|
| `github.com/IBM/sarama` | v1.42.1 | Kafka client |
| `github.com/eclipse/paho.mqtt.golang` | v1.5.1 | MQTT client |
| `github.com/gorilla/websocket` | v1.5.3 | WebSocket |
| `github.com/hibiken/asynq` | v0.26.0 | Task queue |

### 3.4 Security & Validation

| Dependency | Version | Purpose |
|-----------|---------|---------|
| `github.com/golang-jwt/jwt/v4` | v4.5.2 | JWT (RS256) |
| `github.com/go-playground/validator/v10` | v10.30.2 | Input validation |
| `golang.org/x/crypto` | v0.53.0 | bcrypt password hashing |
| `golang.org/x/time` | v0.14.0 | Rate limiting |

### 3.5 Email

| Dependency | Version | Purpose |
|-----------|---------|---------|
| `gopkg.in/gomail.v2` | - | SMTP email |
| `github.com/matcornic/hermes/v2` | v2.1.0 | HTML email templates |

### 3.6 Testing

| Dependency | Version | Purpose |
|-----------|---------|---------|
| `github.com/stretchr/testify` | v1.11.1 | Test assertions |
| `github.com/alicebob/miniredis/v2` | - | Redis mock |
| `github.com/prometheus/client_golang` | v1.23.2 | Prometheus metrics |

---

## 4. Configuration

### 4.1 Config Structure

All configuration is managed through **Viper** with support for YAML files and environment variables.

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

### 4.2 Key Config Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `server.port` | 5000 | HTTP server port |
| `server.mode` | Development | Environment (Development/Production) |
| `server.timezone` | Asia/Bangkok | Application timezone |
| `server.baseUrl` | http://localhost:5000/api | Base API URL |
| `postgres.dbname` | icmon | Database name |
| `redis.addr` | localhost:6379 | Redis address |
| `influxdb.url` | http://localhost:9087 | InfluxDB URL |
| `influxdb.bucket` | icmon | InfluxDB bucket |
| `kafka.brokers` | [localhost:9092] | Kafka broker list |
| `mqtt.broker` | tcp://localhost:1891 | MQTT broker |
| `jwt.accessTokenExpireDuration` | 30m | Access token TTL |
| `jwt.refreshTokenExpireDuration` | 168h | Refresh token TTL |

### 4.3 Config Files

- `config/config.default.yml` — base config for all environments
- `config/config.dev.yml` — development overrides

Environment variables map to config keys automatically via `BindEnvs()`.

---

## 5. Entry Points & CLI Commands

### 5.1 Entry Points

The project has **4 entry points** serving different purposes:

| Entry Point | File | Purpose |
|-------------|------|---------|
| **Main API** | `cmd/api/main.go` | Primary HTTP server (Port 5000) |
| **Kafka Service** | `cmd/kafka/main.go` | Kafka + WebSocket service (Port 5051) |
| **WebSocket Service** | `cmd/websocket/main.go` | WebSocket + Redis Queue (Port 8080) |
| **CLI Root** | `cmd/root.go` | Cobra commands (serve, migrate, worker, initdata) |

### 5.2 CLI Commands

| Command | File | Description |
|---------|------|-------------|
| `serve` | `cmd/serve.go` | Start HTTP server with all services (Postgres, Redis, MQTT, InfluxDB, auto-migration) |
| `migrate` | `cmd/migrate.go` | Run GORM AutoMigrate (~180 models) with skip-list |
| `worker` | `cmd/worker.go` | Start Asynq background task processor |
| `initdata` | `cmd/initdata.go` | Create initial superuser from config |
| `kafka` | `cmd/kafka/kafka.go` | Start Kafka producer + consumer + WebSocket |

### 5.3 Global Flags

- `--config` — path to config file (default: `$HOME/.go-base.yaml`)

---

## 6. Directory Structure

```
icmongolang/
├── main.go                          # App entry (calls cmd.Execute())
│
├── cmd/                             # Entry points & CLI commands (9 files)
│   ├── api/main.go                  # API server entry
│   ├── kafka/main.go                # Kafka server entry
│   ├── websocket/main.go            # WebSocket server entry
│   ├── root.go                      # Cobra root command
│   ├── serve.go                     # `serve` command
│   ├── migrate.go                   # `migrate` command
│   ├── worker.go                    # `worker` command
│   ├── initdata.go                  # `initdata` command
│   └── kafka/kafka.go               # `kafka` command
│
├── config/                          # Configuration files
│   ├── config.go                    # Config struct & loader
│   ├── config.default.yml           # Default configuration
│   └── config.dev.yml               # Development overrides
│
├── internal/                        # Application core (159 Go files)
│   ├── alarm/                       # Alarm management module
│   ├── auth/                        # Authentication module
│   ├── delivery/rest/               # REST router & handlers
│   ├── distributor/                 # Asynq task distributor
│   ├── influxdb/                    # InfluxDB adapter
│   ├── iot/                         # IoT data module (largest)
│   ├── items/                       # Items CRUD module
│   ├── kafka/                       # Kafka order module
│   ├── middleware/                   # All HTTP middleware
│   ├── models/                      # All database models
│   ├── mqtt/                        # MQTT module
│   ├── processor/                   # Task processing
│   ├── queue/                       # Queue abstraction
│   ├── repository/                  # Repository interfaces
│   ├── server/                      # Server setup & lifecycle
│   ├── template/                    # Template module
│   ├── usecase/                     # Usecase interfaces
│   ├── users/                       # User management module
│   ├── websocket/                   # WebSocket module
│   └── worker/                      # Asynq worker
│
├── pkg/                             # Shared libraries (30 Go files)
│   ├── cryptpass/                   # bcrypt password hashing
│   ├── db/postgres/                 # PostgreSQL connection
│   ├── db/redis/                    # Redis connection
│   ├── helpers/                     # Format & IoT alarm helpers
│   ├── influxdb/                    # InfluxDB client wrapper
│   ├── jwt/                         # JWT (RS256) token management
│   ├── kafka/                       # Kafka producer/consumer
│   ├── logger/                      # Zap structured logger
│   ├── mqtt/                        # MQTT client
│   ├── responses/                   # Standard API responses
│   ├── utils/                       # Validator & form utilities
│   └── websocket/                   # WebSocket hub & client
│
├── db/                              # SQL files & dumps
│   ├── icmon.sql                    # Full database dump (20K+ lines)
│   ├── public.sql                   # Public schema
│   └── sd_iot_device.sql            # IoT device table
│
├── migrations/                      # Manual SQL migrations
│   ├── icmon.sql, public.sql, db.sql
│   ├── sd_user_access_menu.sql
│   ├── 20250619_websocket_tables.sql
│   └── 20250619_kafka_tables.sql
│
├── docs/                            # Documentation
│   ├── swagger.yaml                 # Swagger API spec (2,901 lines)
│   └── docs.md                      # Architecture flow diagrams
│
├── Makefile                         # Build automation
└── README.md, README_RUN.md         # Project guides
```

---

## 7. Database Schema

### 7.1 PostgreSQL

Database: `icmon`  
Full dump: `db/icmon.sql` (20,052 lines, Navicat Premium dump)

**Domain tables (some of 180+ models):**

| Category | Tables | Description |
|----------|--------|-------------|
| **Core** | `item`, `migrations` | Base reference tables |
| **Users** | `users`, `sd_users`, `sd_user_roles`, `sd_user_roles_access`, `sd_user_roles_permision`, `sd_user_access_menu`, `sd_user_log`, `sd_user_log_type`, `sd_user_file` | User management & RBAC |
| **Activity** | `activity_log`, `sd_activity_log`, `sd_activity_type_log`, `sd_audit_log`, `sd_module_log`, `sd_device_log`, `sd_mqtt_log` | All audit trails |
| **Devices** | `devices`, `sd_iot_device`, `sd_iot_device_type`, `sd_iot_device_action`, `sd_iot_device_action_log`, `sd_iot_device_action_user`, `sd_iot_device_alarm_action`, `sd_device_alert`, `sd_device_config`, `sd_device_status`, `sd_device_category`, `sd_device_group`, `sd_device_member`, `sd_device_notification_config`, `sd_device_schedule`, `sd_device_status_history` | Device management |
| **IoT** | `iot_data`, `sd_iot_data`, `sd_iot_group`, `sd_iot_host`, `sd_iot_location`, `sd_iot_sensor`, `sd_iot_type`, `sd_iot_setting`, `sd_iot_api`, `sd_iot_token`, `sd_iot_email`, `sd_iot_line`, `sd_iot_sms`, `sd_iot_telegram`, `sd_iot_nodered`, `sd_iot_mqtt`, `sd_iot_influxdb`, `sd_iot_schedule`, `sd_iot_schedule_device` | IoT configuration & data |
| **Alarms** | `sd_iot_alarm_device`, `sd_iot_alarm_device_event`, `sd_alarm_process_log`, `sd_alarm_process_log_email`, `sd_alarm_process_log_line`, `sd_alarm_process_log_mqtt`, `sd_alarm_process_log_sms`, `sd_alarm_process_log_telegram`, `sd_alarm_process_log_temp` | Alarm management |
| **Air Control** | `air_controls`, `air_mods`, `air_periods`, `air_setting_warnings`, `air_warnings` (with `sd_` prefixed variants and device maps) | Air quality control |
| **Notifications** | `noti_notifications`, `notification_logs`, `notification_groups`, `notification_devices`, `notification_types`, `sd_notification_channel`, `sd_notification_condition`, `sd_notification_log`, `sd_notification_type`, `sd_channel_template`, `sd_group_notification_config` | Notification system |
| **Dashboards** | `sd_dashboard_config` (alias: `tnb`) | Dashboard configuration |
| **API Keys** | `sd_api_keys` | API key management |
| **System** | `sd_system_settings` | System settings |
| **Reports** | `sd_report_data`, `sd_schedule_process_log` | Data reporting |
| **Schedules** | `sd_device_schedule`, `sd_iot_schedule`, `sd_iot_schedule_device` | Scheduling |
| **Commands** | `command_log` | Command audit log |
| **Kafka** | `orders` | Kafka order events |

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

### 7.3 Key Database Extensions

- `uuid-ossp` — UUID generation
- `pgcrypto` — Cryptographic functions

### 7.4 Migration Strategy

**Dual approach:**
1. **GORM AutoMigrate** — used at startup or via `migrate` command for the 180+ GORM models. A skip-list handles tables with composite primary keys or complex schemas.
2. **Manual SQL** — files in `migrations/` directory for complex schema changes.

---

## 8. API Routes

### 8.1 Base Path

All API routes are under `/api` (except Kafka and WebSocket services which use `/`).

### 8.2 Public Routes

| Method | Path | Description |
|--------|------|-------------|
| GET | `/metrics` | Prometheus metrics |
| GET | `/apimetric` | Custom JSON application metrics |
| GET | `/health` | Health check |
| GET | `/api/ping` | Ping/pong |
| POST | `/api/auth/login` | Login with credentials |
| POST | `/api/auth/signin` | Sign-in |
| GET | `/api/auth/publickey` | Get JWT public key |
| GET | `/api/auth/verifyemail` | Verify email address |
| POST | `/api/auth/forgotpassword` | Request password reset |
| PATCH | `/api/auth/resetpassword` | Reset password with token |
| POST | `/api/register` | Register new user |
| POST | `/api/users` | Create user |
| POST | `/api/signin` | Sign-in by email |
| POST | `/api/login` | Sign-in by username |

### 8.3 Protected Routes (Authenticated)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/auth/refresh` | Refresh access token |
| GET | `/api/auth/logout` | Logout current session |
| GET | `/api/auth/logoutall` | Logout all sessions |
| GET | `/api/user/me` | Get own profile |
| PUT | `/api/user/me` | Update own profile |
| PATCH | `/api/user/me/updatepass` | Change own password |
| GET | `/api/item/` | List items |
| POST | `/api/item/` | Create item |
| GET | `/api/item/{id}` | Get item by ID |
| DELETE | `/api/item/{id}` | Delete item |
| PUT | `/api/item/{id}` | Update item |

### 8.4 Protected Routes (SuperUser/Admin)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/user/` | List all users |
| POST | `/api/user/` | Create user (admin) |
| GET | `/api/user/{id}` | Get user by ID |
| PUT | `/api/user/{id}` | Update user |
| DELETE | `/api/user/{id}` | Delete user |
| PATCH | `/api/user/{id}/role` | Update user role |
| PATCH | `/api/user/{id}/updatepass` | Update any user's password |
| GET | `/api/user/{id}/logoutall` | Force logout all sessions |

### 8.5 Kafka Service Routes (Port 5051)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/orders` | Create order (produces Kafka message) |
| GET | `/ws` | WebSocket endpoint |
| GET | `/health` | Health check |

---

## 9. Middleware Pipeline

### 9.1 Execution Order

Requests pass through the following middleware chain (defined in `internal/delivery/rest/router.go`):

```
1. panic recovery (Recoverer)
2. request ID assignment (RequestID)
3. real client IP (RealIP)
4. process timeout (Timeout)
5. CORS headers (CORS)
6. security headers (CSP, HSTS, XSS, NoSniff, FrameDeny, etc.)
7. per-IP rate limiting (10 req/s, burst 20)
8. request/response logging
9. Prometheus metrics monitoring
10. route handler (JWT auth → handler → response)
```

### 9.2 Middleware Components

| Middleware | File | Description |
|-----------|------|-------------|
| **CORS** | `middleware/cors.go` | Environment-aware CORS (development: localhost; production: BaseUrl) |
| **Security** | `middleware/security.go` | Security headers: CSP, HSTS, XSS, NoSniff, FrameDeny, Referrer, Permissions-Policy, Cache-Control |
| **Rate Limit** | `middleware/rate_limit.go` | Per-IP rate limiter using `golang.org/x/time/rate` (10 req/s, burst 20, cleanup every 30s) |
| **Logging** | `middleware/logging.go` | Logs method, path, status code (no request body) |
| **Monitoring** | `middleware/monitoring.go` | Prometheus `http_requests_total`, `http_request_duration_seconds` + `/apimetric` JSON endpoint |
| **JWT Auth** | `middleware/jwtauth.go` | Access/refresh token verification (RS256), current user context, super user check, active user check |
| **Prometheus** | `middleware/prometheus.go` | System metrics: goroutines, memory, GC, CPU, network, Redis connection stats |

---

## 10. Module Breakdown

### 10.1 Users Module (`internal/modules/users/`)

**Layers:**
- `models/user.go` — User domain entity
- `repository/pg_repository.go` — PostgreSQL CRUD (GORM), session management
- `usecase/usecase.go` — Registration, sign-in, profile management, password change, forgot/reset password, email verification, role management
- `delivery/http/handlers.go` — HTTP handlers for all user endpoints
- `distributor/distributor.go` — Asynq task distribution (email tasks)
- `processor/processor.go` — Asynq task processing (email sending)
- `worker/worker.go` — Worker task definitions

**Key Features:**
- Full CRUD with Role-Based Access Control
- Email verification flow with token
- Forgot/reset password flow
- Session management (Redis-backed refresh tokens)
- Asynq-based async email delivery

### 10.2 Auth Module (`internal/modules/auth/`)

**Layers:**
- `delivery/http/handlers.go` — Login, sign-in, refresh token, logout, public key, email verify
- `route.go` — Route definitions

**Key Features:**
- JWT access + refresh token pair (RS256 RSA keys)
- Refresh tokens stored in Redis
- Public key endpoint for external verification
- Token rotation on refresh

### 10.3 Items Module (`internal/modules/items/`)

**Layers:**
- `repository/` — Standard CRUD repository
- `usecase/` — Business logic
- `delivery/http/` — HTTP handlers

Simple CRUD module — serves as a reference pattern for basic resource management.

### 10.4 IoT Module (`internal/modules/iot/`) — Largest Module

**Layers:**
- `models/` — Parallel domain models (Device, IotData, AirControl, AirMod, AirPeriod, AirWarning, etc.)
- `repository/iot_data_repo.go` — IoT data queries with pagination, device validation
- `repository/device_repo.go` — Device CRUD with dynamic filtering
- `usecase/usecase.go` — Complex business logic (device health, alarms, air control, data aggregation, dashboard)
- `presenter/` — Response formatting

**Key Features:**
- Device management with dynamic column filtering (SQL injection protected via whitelist)
- IoT sensor data retrieval with DB-level pagination
- Air control system (mods, periods, warnings with device mapping)
- Alarm rule evaluation and alert generation
- Device status history tracking
- Dashboard aggregation queries

### 10.5 Alarm Module (`internal/modules/alarm/`)

**Layers:**
- `repository/alarm_log_repo.go` — Alarm log queries with dynamic filtering (SQL injection protected)
- `usecase/` — Alarm processing logic
- `delivery/` — HTTP handlers

**Key Features:**
- Alarm log retrieval with multi-channel filtering
- Supports email, line, SMS, MQTT, Telegram alarm channels
- Paginated and filterable alarm history

### 10.6 Kafka Module (`internal/modules/kafka/`)

**Layers:**
- `models/` — Order event model
- `repository/` — Kafka repository
- `usecase/` — Order processing logic
- `delivery/` — HTTP + WebSocket handlers

**Key Features:**
- Order event production/consumption via Apache Kafka (sarama)
- Real-time order status streaming via WebSocket
- Independent service (runs on port 5051, gorilla/mux)

### 10.7 WebSocket Module (`internal/modules/websocket/`)

**Layers:**
- `models/` — WebSocket message types
- `repository/` — WebSocket session repository
- `usecase/` — Message handling logic
- `delivery/` — WebSocket connection management

**Key Features:**
- Real-time bidirectional communication
- Session management and message broadcasting
- Can run independently (port 8080) or embedded

### 10.8 MQTT Module (`internal/modules/mqtt/`)

**Layers:**
- `delivery/` — MQTT message handlers
- `presenter/` — Data formatting
- `usecase/` — MQTT message processing

**Key Features:**
- IoT device communication via MQTT protocol
- Message subscription and processing pipeline
- Reconnection with exponential backoff

### 10.9 Distributor / Processor / Worker

**Distributor** (`internal/distributor/`):
- Asynq task distribution interface
- Creates tasks for email verification, password reset

**Processor** (`internal/processor/`):
- Asynq task processing (email sending)

**Worker** (`internal/worker/`):
- Asynq worker server configuration
- Mux for task routing
- Graceful shutdown

### 10.10 Template Module (`internal/template/`)

A reference module template for creating new domain modules. Contains pre-built structure:
- Delivery (HTTP handlers, routes)
- Repository (PostgreSQL + Redis)
- Usecase (business logic)
- Models (domain entities)
- Presenter (response formatting)

---

## 11. Security Measures

### 11.1 Authentication & Authorization

| Measure | Implementation |
|---------|---------------|
| **JWT RS256** | RSA-based access + refresh tokens with configurable expiry |
| **Bcrypt passwords** | All passwords hashed via bcrypt (`pkg/cryptpass`) |
| **Refresh token rotation** | Old tokens invalidated on refresh |
| **Redis-backed sessions** | Refresh tokens stored in Redis with TTL |
| **SuperUser middleware** | Role-based access for admin endpoints |
| **Active user check** | Inactive users blocked from API access |

### 11.2 SQL Injection Prevention

| Fix | Location |
|-----|----------|
| Column name whitelist | `internal/modules/alarm/repository/alarm_log_repo.go` |
| Column name + ORDER BY whitelist | `internal/modules/iot/repository/device_repo.go` |

Dynamic `.Where()` clauses and ORDER BY parameters are validated against allowed field lists before being passed to GORM.

### 11.3 CORS Protection

- **Development**: Allows localhost origins (with and without port)
- **Production**: Restricted to `config.Server.BaseUrl` only
- No wildcard (`*`) origin ever allowed

### 11.4 Rate Limiting

- **Limit**: 10 requests per second per IP (down from 2000)
- **Burst**: 20 requests
- **Cleanup**: Idle IPs removed every 30 seconds
- **IP source**: Uses only `RemoteAddr` (prevents `X-Forwarded-For` spoofing)

### 11.5 Security Headers (CSP, HSTS, etc.)

All responses include:
- `Content-Security-Policy`
- `Strict-Transport-Security`
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 0` (modern approach)
- `Referrer-Policy`
- `Permissions-Policy`
- `Cache-Control`

### 11.6 PII Protection

| What was fixed | Where |
|----------------|-------|
| Removed request body from logs | `internal/middleware/logging.go` |
| Removed email/password from auth logs | `internal/modules/auth/delivery/http/handlers.go` |
| Removed PII from user usecase logs | `internal/modules/users/usecase/usecase.go` |

### 11.7 Timeouts & Context

| Component | Timeout | Where |
|-----------|---------|-------|
| Process timeout | Configurable (default 60s) | `internal/delivery/rest/router.go` |
| InfluxDB queries | 30 seconds | `pkg/influxdb/client.go` |
| Redis operations | 5 seconds per operation | `pkg/db/redis/redis_conn.go` |
| Server shutdown | 5 seconds graceful | `internal/server/server.go` |

---

## 12. Performance Optimizations

### 12.1 Database Pagination

| Before | After |
|--------|-------|
| In-memory: fetch 10,000 records → slice in Go | DB-level: `LIMIT ? OFFSET ?` |
| Wasted memory & CPU for every paginated query | Efficient index-based pagination |

Implemented in `internal/modules/iot/usecase/usecase.go` via `CountByDeviceID()` + `GetByDeviceID(limit, offset)`.

### 12.2 `structToMap` Optimization

| Before | After |
|--------|-------|
| `json.Marshal` then `json.Unmarshal` into map[string]interface{} | Direct reflection-based conversion |
| Double allocation per call | Single pass, no intermediate JSON |

Implemented in `internal/modules/iot/usecase/usecase.go`.

### 12.3 Removed Artificial Latency

Removed `time.Sleep(2 * time.Second)` from Kafka order processing (`internal/modules/kafka/usecase/order_usecase.go`) — was a development simulation that blocked production throughput.

### 12.4 Context Timeouts

- All InfluxDB queries now bounded by 30s context timeout
- All Redis operations bounded by 5s per-operation timeout

### 12.5 Rate Limit Tuning

Rate limit reduced from 2000 req/s to 10 req/s (sensible default) — prevents accidental DoS and aligns with production standards.

---

## 13. Monitoring & Metrics

### 13.1 Prometheus Metrics (`/metrics`)

| Metric | Type | Description |
|--------|------|-------------|
| `http_requests_total` | Counter | Total HTTP requests (labels: method, path, status) |
| `http_request_duration_seconds` | Histogram | Request duration distribution |

### 13.2 Custom JSON Metrics (`/apimetric`)

| Metric | Description |
|--------|-------------|
| `go_goroutines` | Current goroutine count |
| `go_memory_alloc` | Current memory allocation (bytes) |
| `go_gc_count` | Number of completed GC cycles |
| `go_cpu_usage` | CPU usage percentage |
| `memory_usage_percent` | Memory usage percentage |
| `network_in_bytes_total` | Total network input |
| `network_out_bytes_total` | Total network output |
| `redis_connections` | Active Redis connections |
| `redis_hits_per_second` | Redis cache hit rate |
| `redis_misses_per_second` | Redis cache miss rate |
| `redis_max_connections` | Redis max connections |

### 13.3 Health Check

`GET /api/health` returns server status and can be used by load balancers and monitoring systems.

---

## 14. Setup & Deployment

### 14.1 Prerequisites

- Go 1.25+
- PostgreSQL
- Redis
- (Optional) InfluxDB, MQTT broker, Kafka, SMTP server

### 14.2 Quick Start

```bash
# 1. Start dependencies (PostgreSQL, Redis)
docker-compose up -d

# 2. Run migrations
go run cmd/api/main.go migrate

# 3. Start server
go run cmd/api/main.go serve

# 4. Initialize superuser
go run cmd/api/main.go initdata
```

### 14.3 Available Make Commands

```bash
make clean     # go clean -cache + go clean -modcache
make tidy      # go mod tidy
make download  # go mod download
make verify    # go mod verify
make migrate   # go run cmd/api/main.go migrate
make swag      # swag init -g cmd/api/main.go
make vendor    # go mod vendor
make test      # go test ./...
make run       # Full setup + air (hot-reload)
```

### 14.4 Running Independent Services

```bash
# Kafka + WebSocket service (port 5051)
go run cmd/kafka/main.go

# WebSocket + Redis Queue service (port 8080)
go run cmd/websocket/main.go

# Background task worker
go run cmd/api/main.go worker
```

### 14.5 Default Credentials

- **Email:** root@gmail.com
- **Password:** root_password

### 14.6 Swagger UI

Available at: `http://localhost:5000/swagger/index.html`

---

## 15. Project Statistics

| Category | Count |
|----------|-------|
| **Total Go files** | 198 |
| **`cmd/` files** | 9 |
| **`internal/` files** | 159 |
| **`pkg/` files** | 30 |
| **Database models** | 180+ |
| **API routes** | 30+ |
| **Middleware components** | 10 |
| **Database schema size** | 20,052+ lines SQL |
| **Swagger spec size** | 2,901 lines YAML |
| **External dependencies** | 34 direct |

---

*End of Technical Report — English Version*
