# golangapi

> Production-ready Go RESTful API boilerplate with Chi, GORM, PostgreSQL, Redis, and enterprise features

## 🌟 Features

### 🚀 Core Features
- **🏭 Clean Architecture** - Three-layer architecture (Repository/Service/Handler) with comprehensive dependency injection
- **🔒 JWT Authentication** - Complete authentication system with access/refresh tokens and token blacklisting
- **👥 User Management** - Full CRUD operations with basic role-based protection (admin-only create)
- **📝 Structured Logging** - Unified logging via `pkg/logger` (built on Go's slog) with trace/request context
- **🚫 Rate Limiting** - IP-based request throttling with automatic cleanup
- **📊 Health Monitoring** - Comprehensive health checks with dependency monitoring
- **🌐 Redis Cache** - Production-ready caching layer with TTL management and object serialization
- **📦 Message Queue** - Redis-based pub/sub messaging with worker pools and dead letter queue support
- **💼 Transaction Management** - GORM transaction manager with nested transaction support
- **🛡️ Security** - Multiple security layers including CORS, security headers, and input validation

### 🛠️ Middleware Stack
- **Request Context** - Trace IDs, request IDs, and user context propagation
- **Security Headers** - CSP, HSTS, X-Frame-Options, XSS Protection
- **CORS Handling** - Configurable cross-origin resource sharing
- **Panic Recovery** - Application-level panic handling with graceful error responses
- **Request Logging** - Structured request/response logging with performance metrics
- **Authentication** - JWT middleware with role-based route protection
- **Input Validation** - Comprehensive request validation using go-playground/validator

### 📈 Health & Monitoring
- **Health Endpoints** - Basic, detailed, readiness, and liveness probes
- **Dependency Checks** - Database and Redis connection status via `/health/detailed`
- **System Metrics (Optional)** - Runtime/memory snapshot handler available for wiring
- **Performance Tracking (Optional)** - Monitoring middleware + metrics handler available
- **Kubernetes Ready** - `/ready` and `/live` probes

## Directory Structure

Design reference:

- [go project layout](https://kongnakornna/goreststarterio)
- [go modules layout](https://go.dev/doc/modules/layout)

```md
project-root/
├── api/                          # API related files
│   └── app/                      # API app docs
│       └── docs.go               # docs.go
│       └── swagger.json          # Swagger documentation
├── cmd/                          # Main program entry
│   └── app/                      # Application
│       └── main.go               # Program entry point
├── configs/                      # Configuration files (optimized as single source)
├── deploy/                       # Deployment configurations (simplified to essential scripts)
│   └── docker/                  
│   └── k8s/                     
├── internal/                     # Internal application code
│   ├── apps/                     # Application composition roots
│   │   └── app/                  # Current app
│   │       ├── bootstrap/        # App bootstrap (DI, startup)
│   │       └── router/           # App router
│   ├── core/                     # Domain modules (auth/user/health)
│   ├── platform/                 # Infrastructure (config/db)
│   └── transport/                # HTTP transport helpers (httpx/middleware)
├── migrations/                   # Database migration files (version controlled)
├── pkg/                          # External packages (independent reusable components)
│   ├── cache/                    # Redis cache and strategies
│   ├── errors/                   # Custom error handling package
│   ├── jwt/                      # JWT helpers
│   ├── logger/                   # Structured logger
│   ├── queue/                    # Redis queue helpers
│   ├── transaction/              # Transaction manager
│   └── utils                     # Common utility functions
├── scripts/                      # Development and deployment scripts (simplified workflow)
├── .air.toml                     # Development hot-reload configuration
├── go.mod                        # Go module definition
└── README.md                     # Project documentation
```

## 🚀 Quick Start

## 📚 Documentation
See `docs/README.md` for architecture, development, maintenance, and deployment guides.

### Prerequisites
- **Go 1.25+** - [Install Go](https://golang.org/doc/install)
- **PostgreSQL 12+** - [Install PostgreSQL](https://postgresql.org/download/)
- **Redis 6+** - [Install Redis](https://redis.io/download)

### Installation

```bash
# Clone the repository
git clone https://golangapi.git
cd golangapi

# Install dependencies
Setup 
go mod tidy
update
go mod download

# Copy and configure the config file
cp configs/config.example.yaml configs/config.yaml
# Edit configs/config.yaml with your database and Redis settings
```

### Development Mode

```bash
# Run development server (with auto-reload)
./scripts/dev.sh

# Or run directly
go run cmd/app/main.go
```

### Access API Documentation

After starting the service, visit **http://localhost:7001/swagger** to view the interactive API documentation.

## 📚 API Endpoints

# http://localhost:7001/api/v1/users

```bash
http://localhost:7001/api/v1/health
http://localhost:7001/health
http://localhost:7001/health/detailed
http://localhost:7001/ready
http://localhost:7001/live
http://localhost:7001/api/v1/auth/login
http://localhost:7001/api/v1/auth/refresh
http://localhost:7001/api/v1/account/logout
http://localhost:7001/api/v1/users
http://localhost:7001/api/v1/users/{id} 
http://localhost:7001/version
http://localhost:7001/status
http://localhost:7001/api/v1/account/logout 

```

### 🏥 Health Check Endpoints
- `GET /health` - Basic health check with uptime
- `GET /health/detailed` - Detailed health check (includes DB and Redis status)
- `GET /ready` - Kubernetes readiness probe
- `GET /live` - Kubernetes liveness probe

### 🔐 Authentication Endpoints (Public)
- `POST /api/v1/auth/login` - User authentication
- `POST /api/v1/auth/refresh` - Refresh JWT token

### 🔒 Account Management Endpoints (Protected)
- `POST /api/v1/account/logout` - User logout (invalidates tokens)

### 👥 User Management Endpoints (Protected)
- `GET /api/v1/users` - List users with pagination
- `POST /api/v1/users` - Create new user (Admin only)
- `GET /api/v1/users/{id}` - Get user details by ID
- `PUT /api/v1/users/{id}` - Update user information
- `DELETE /api/v1/users/{id}` - Delete user

### 📊 System Endpoints
- `GET /version` - API version information
- `GET /status` - Service status

## ⚙️ Configuration

### Configuration Files
The application uses YAML configuration files with environment variable override support:

```bash
# Primary configuration
configs/config.yaml          # Main configuration (create from example)
configs/config.example.yaml  # Example configuration template
configs/config.production.yaml  # Production-specific overrides
```

The config file path can be overridden with `CONFIG_PATH` (defaults to `configs/config.yaml`).

### Environment Variables
All configuration values can be overridden using environment variables with `APP_` prefix:

```bash
# Config Path
CONFIG_PATH=configs/config.yaml

# Server Configuration
APP_SERVER_PORT=7001
APP_SERVER_TIMEOUT=30s
APP_SERVER_READ_TIMEOUT=15s
APP_SERVER_WRITE_TIMEOUT=15s

# Database Configuration
APP_DB_HOST=localhost
APP_DB_PORT=5432
APP_DB_USERNAME=postgres
APP_DB_PASSWORD=your-password
APP_DB_NAME=myapp
APP_DB_SSLMODE=disable
APP_DB_MAX_OPEN_CONNS=20
APP_DB_MAX_IDLE_CONNS=5
APP_DB_CONN_MAX_LIFETIME=1h

# Redis Configuration
APP_REDIS_ENABLED=true
APP_REDIS_HOST=localhost
APP_REDIS_PORT=6379
APP_REDIS_PASSWORD=""
APP_REDIS_DB=0

# JWT Configuration
APP_JWT_SECRET=your-secure-secret-key-change-in-production
APP_JWT_ACCESS_TOKEN_EXP=24h
APP_JWT_REFRESH_TOKEN_EXP=168h
APP_JWT_ISSUER=golangapi

# Logging Configuration
APP_LOG_LEVEL=info
APP_LOG_FILE=logs/app.log
APP_LOG_CONSOLE=true
```

### Configuration Structure
```yaml
app:
  server:
    port: 7001
    timeout: 30s
    read_timeout: 15s
    write_timeout: 15s
  database:
    driver: postgres
    host: localhost
    port: 5432
    # ... (see config.example.yaml for full structure)
```

## 🚀 Build and Deploy

### Build Binary

```bash
# Build for current platform
go build -o app cmd/app/main.go

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o app cmd/app/main.go

# Build with version info
go build -ldflags="-s -w" -o app cmd/app/main.go
```

### Docker Deployment

```bash
# Build Docker image
docker build -t golangapi -f deploy/docker/Dockerfile .

# Run with Docker Compose (includes PostgreSQL and Redis)
cd deploy/docker
docker compose up -d

# Run container only (requires external database and Redis)
docker run -p 7001:7001 \
  -e APP_DB_HOST=your-db-host \
  -e APP_REDIS_HOST=your-redis-host \
  golangapi
```

### Kubernetes Deployment

```bash
# Deploy to Kubernetes
kubectl apply -f deploy/k8s/

# Check deployment status
kubectl get pods -l app=golangapi
kubectl logs -f deployment/golangapi
```

## 🧪 Testing

### Run Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage report in browser
go tool cover -html=coverage.out

# Run tests with race detection
go test -race ./...

# Run tests with verbose output
go test -v ./...
```

### Test Categories

```bash
# Unit tests (service layer)
go test ./internal/core/user/service/

# Integration tests (if available)
go test -tags=integration ./...

# Benchmark tests
go test -bench=. ./...
```

## 🛠️ Technology Stack

### Core Framework & Libraries
- **Web Framework**: `chi/v5` - Lightweight, fast HTTP router with middleware support
- **ORM**: `GORM v1.30.0` - Feature-rich ORM with auto-migration and relations
- **Database Driver**: `gorm.io/driver/postgres` - PostgreSQL driver for GORM
- **Cache**: `redis/go-redis/v9` - Redis client with pipeline and pub/sub support

### Authentication & Security
- **JWT**: `golang-jwt/jwt/v5` - JSON Web Token implementation
- **Password Hashing**: `golang.org/x/crypto/bcrypt` - Secure password hashing
- **Input Validation**: `go-playground/validator/v10` - Struct validation with tags
- **Rate Limiting**: `golang.org/x/time/rate` - Token bucket rate limiting

### Configuration & Utilities
- **Configuration**: `spf13/viper` - Configuration management (YAML, ENV, JSON)
- **Logging**: `pkg/logger` - Structured logging (based on Go slog)
- **Testing**: `stretchr/testify` - Testing toolkit with assertions and mocks

### Documentation & Development
- **API Documentation**: `swaggo/swag` - Swagger/OpenAPI 3.0 documentation generator
- **HTTP Swagger UI**: `swaggo/http-swagger/v2` - Swagger UI integration

## 🌟 Architecture & Design

### Clean Architecture Implementation
- **Handler Layer** - HTTP request handling and response formatting
- **Service Layer** - Business logic and transaction management
- **Repository Layer** - Data access and database operations
- **Dependency Injection** - Interface-based design with comprehensive DI container

### Key Design Patterns
- **Repository Pattern** - Abstract data access layer
- **Service Pattern** - Encapsulated business logic
- **Middleware Chain** - Composable request processing
- **Factory Pattern** - Component initialization
- **Observer Pattern** - Configuration watching and hot reload

### Security Features
- **JWT Authentication** - Stateless authentication with token blacklisting
- **Role-Based Access Control** - Admin/User role separation
- **Security Headers** - CSP, HSTS, X-Frame-Options, XSS protection
- **Input Validation** - Request validation with custom error messages
- **Rate Limiting** - IP-based request throttling
- **Password Security** - bcrypt hashing with configurable cost

### Performance Features
- **Connection Pooling** - Database and Redis connection management
- **Caching Layer** - Redis-based caching with TTL management
- **Structured Logging** - High-performance logging with context
- **Graceful Shutdown** - Zero-downtime deployments
- **Health Checks** - Kubernetes-ready probes

## 📝 Usage Examples

### Authentication Flow
```bash
# Login
curl -X POST http://localhost:7001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "password"}'

# Use the returned token for authenticated requests
curl -X GET http://localhost:7001/api/v1/users \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### User Management
```bash
# Create user (Admin only)
curl -X POST http://localhost:7001/api/v1/users \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password", "role": "user"}'

# Get user list with pagination
curl -X GET "http://localhost:7001/api/v1/users?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Health Monitoring
```bash
# Basic health check
curl http://localhost:7001/health

# Detailed health with dependencies
curl http://localhost:7001/health/detailed

# Readiness probe
curl http://localhost:7001/ready

# Liveness probe
curl http://localhost:7001/live
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Go Project Layout](https://kongnakornna/goreststarterio) - Standard Go project structure
- [Chi Router](https://go-chi/chi) - Lightweight HTTP router
- [GORM](https://gorm.io/) - The fantastic ORM library for Golang



# บทความภาษาไทย: การออกแบบสถาปัตยกรรมซอฟต์แวร์ด้วย Golang และ LangChain

## สารบัญ
1. [บทนำ](#1-บทนำ-introduction)
2. [บทนิยาม](#2-บทนิยาม-definitions)
3. [สถาปัตยกรรม Golang แบบ DDD และ Clean Architecture](#3-สถาปัตยกรรม-golang-แบบ-ddd-และ-clean-architecture)
4. [สถาปัตยกรรม LangChain FastAPI](#4-สถาปัตยกรรม-langchain-fastapi)
5. [เวิร์กโฟลว์การพัฒนา](#5-เวิร์กโฟลว์การพัฒนา-development-workflow)
6. [ผังการไหลของข้อมูล](#6-ผังการไหลของข้อมูล-data-flow-architecture)
7. [แม่แบบโค้ด](#7-แม่แบบโค้ด-code-templates)

---

## 1. บทนำ (Introduction)

### 1.1 Golang กับ Domain-Driven Design
Golang หรือ Go เป็นภาษาโปรแกรมมิ่งที่พัฒนาโดย Google มีจุดเด่นด้านประสิทธิภาพสูง การทำงานพร้อมกัน (Concurrency) ที่ยอดเยี่ยม และความเรียบง่ายของภาษา การนำ Domain-Driven Design (DDD) มาประยุกต์ใช้กับ Golang ร่วมกับ Clean Architecture ช่วยให้การพัฒนาซอฟต์แวร์มีโครงสร้างที่ชัดเจน แยกส่วนการทำงานออกจากกันอย่างเป็นระบบ และง่ายต่อการบำรุงรักษา

**ประโยชน์ของ DDD + Clean Architecture ใน Golang:**
- แยก Business Logic ออกจาก Infrastructure
- ทำให้ Testing ทำได้ง่ายขึ้น
- ระบบมีความยืดหยุ่นสูง สามารถเปลี่ยนเทคโนโลยีได้โดยไม่กระทบ Core Business
- ทีมพัฒนาสามารถทำงานแบบขนานได้ (Parallel Development)

### 1.2 LangChain กับ FastAPI
LangChain เป็นเฟรมเวิร์กสำหรับพัฒนาแอปพลิเคชันที่ใช้ Large Language Models (LLMs) ร่วมกับ FastAPI ซึ่งเป็นเว็บเฟรมเวิร์กความเร็วสูงของ Python การรวมกันนี้ช่วยให้สร้าง AI-powered applications ที่มีประสิทธิภาพ

**ประโยชน์ของ LangChain + FastAPI:**
- รองรับการทำงานกับ LLM หลายรูปแบบ (OpenAI, Hugging Face, etc.)
- มีเครื่องมือสำหรับสร้าง Chains, Agents, และ Memory
- FastAPI ให้ประสิทธิภาพสูงและรองรับ async/operation
- มีระบบ Documentation อัตโนมัติ (Swagger/OpenAPI)

---

## 2. บทนิยาม (Definitions)

### 2.1 คำศัพท์สำคัญใน DDD

| คำศัพท์ | คำอธิบาย |
|--------|----------|
| **Domain** | ขอบเขตของปัญหาที่ระบบต้องการแก้ไข |
| **Entity** | วัตถุที่มี Identity ไม่ซ้ำกันและเปลี่ยนแปลงสถานะได้ |
| **Value Object** | วัตถุที่ไม่มี Identity ค่าของมันคือตัวตนของมัน |
| **Aggregate** | กลุ่มของ Entities และ Value Objects ที่ทำงานร่วมกัน |
| **Repository** | ตัวกลางสำหรับเข้าถึงข้อมูล (Persistence) |
| **Domain Service** | Business logic ที่ไม่อยู่ใน Entity หรือ Value Object |
| **Application Service** | จัดการ Use cases และประสานงานระหว่าง Domain |
| **DTO (Data Transfer Object)** | วัตถุสำหรับส่งข้อมูลระหว่าง layers |

### 2.2 คำศัพท์ใน LangChain

| คำศัพท์ | คำอธิบาย |
|--------|----------|
| **Chain** | ลำดับของการเรียกใช้ LLM หรือ tools |
| **Agent** | ระบบที่สามารถตัดสินใจและเลือกใช้ tools |
| **Tool** | ฟังก์ชันที่ Agent สามารถเรียกใช้ได้ |
| **Memory** | การจดจำประวัติการสนทนา |
| **Embedding** | การแปลงข้อความเป็นเวกเตอร์ |
| **Vector Store** | ฐานข้อมูลสำหรับเก็บ embeddings |

---

## 3. สถาปัตยกรรม Golang แบบ DDD และ Clean Architecture

### 3.1 โครงสร้างโปรเจค Golang (Template)

```
project-root/
├── api/                          # เอกสาร API
│   └── app/
│       ├── docs.go
│       └── swagger.json
│
├── cmd/                          # จุดเริ่มต้นโปรแกรม
│   └── app/
│       ├── grpc.go                # gRPC server entry
│       ├── http.go                 # HTTP server entry
│       ├── worker.go               # Background worker entry
│       └── main.go
│
├── configs/                       # ไฟล์การตั้งค่า
│   ├── config.go
│   ├── config.yaml
│   └── config.example.yaml
│
├── deploy/                        # Deployment configs
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   └── k8s/
│       ├── deployment.yaml
│       └── service.yaml
│
├── internal/                      # โค้ดภายใน (ไม่ถูก import จากภายนอก)
│   │
│   ├── core/                      # Domain Layer
│   │   ├── auth/
│   │   │   ├── entity/            # Domain entities
│   │   │   │   ├── user.go
│   │   │   │   ├── token.go
│   │   │   │   └── role.go
│   │   │   │
│   │   │   ├── repository/        # Repository interfaces
│   │   │   │   ├── user_repository.go
│   │   │   │   └── token_repository.go
│   │   │   │
│   │   │   ├── service/           # Domain services
│   │   │   │   ├── auth_service.go
│   │   │   │   └── password_service.go
│   │   │   │
│   │   │   ├── dto/               # Data Transfer Objects
│   │   │   │   ├── auth_request.go
│   │   │   │   └── auth_response.go
│   │   │   │
│   │   │   ├── model/             # Database models (GORM)
│   │   │   │   ├── user_model.go
│   │   │   │   └── token_model.go
│   │   │   │
│   │   │   ├── routes/            # Route definitions
│   │   │   │   └── auth_routes.go
│   │   │   │
│   │   │   ├── handler/           # HTTP handlers (Application Layer)
│   │   │   │   ├── auth_handler.go
│   │   │   │   └── auth_handler_test.go
│   │   │   │
│   │   │   └── auth.go            # Module setup
│   │   │
│   │   ├── user/
│   │   │   ├── entity/
│   │   │   ├── repository/
│   │   │   ├── service/
│   │   │   ├── dto/
│   │   │   ├── model/
│   │   │   ├── handler/
│   │   │   └── user.go
│   │   │
│   │   ├── iot/
│   │   │   ├── entity/
│   │   │   │   ├── device.go
│   │   │   │   ├── sensor.go
│   │   │   │   └── telemetry.go
│   │   │   ├── repository/
│   │   │   ├── service/
│   │   │   ├── handler/
│   │   │   └── iot.go
│   │   │
│   │   └── health/
│   │       ├── handler/
│   │       │   └── health_handler.go
│   │       └── health.go
│   │
│   ├── platform/                   # Infrastructure Layer
│   │   ├── config/
│   │   │   └── config_loader.go
│   │   │
│   │   ├── db/                     # Database connections
│   │   │   ├── postgres/
│   │   │   │   ├── connection.go
│   │   │   │   └── migration.go
│   │   │   └── redis/
│   │   │       ├── connection.go
│   │   │       └── cache.go
│   │   │
│   │   ├── cache/                  # Cache implementations
│   │   │   ├── redis_cache.go
│   │   │   └── cache_interface.go
│   │   │
│   │   ├── queue/                   # Message queue
│   │   │   ├── redis_queue.go
│   │   │   ├── worker_pool.go
│   │   │   └── dead_letter.go
│   │   │
│   │   └── logger/                  # Logging
│   │       ├── logger.go
│   │       └── middleware.go
│   │
│   └── transport/                   # Transport Layer
│       ├── middleware/
│       │   ├── auth_middleware.go
│       │   ├── logging_middleware.go
│       │   ├── recovery_middleware.go
│       │   ├── rate_limit.go
│       │   ├── cors_middleware.go
│       │   └── security_headers.go
│       │
│       ├── httpx/                    # HTTP utilities
│       │   ├── response.go
│       │   ├── request.go
│       │   └── validator.go
│       │
│       └── utils/
│           ├── context.go
│           └── pagination.go
│
├── pkg/                             # Reusable packages
│   ├── cache/
│   │   ├── cache_strategy.go
│   │   └── cache_errors.go
│   │
│   ├── errors/
│   │   ├── errors.go
│   │   ├── app_errors.go
│   │   └── error_handler.go
│   │
│   ├── jwt/
│   │   ├── jwt_service.go
│   │   └── jwt_claims.go
│   │
│   ├── logger/
│   │   ├── logger.go
│   │   └── fields.go
│   │
│   ├── queue/
│   │   ├── queue_interface.go
│   │   ├── redis_queue.go
│   │   └── job.go
│   │
│   ├── transaction/
│   │   └── transaction_manager.go
│   │
│   ├── oauth/
│   │   ├── google_oauth.go
│   │   └── oauth_interface.go
│   │
│   ├── grpc/
│   │   ├── client/
│   │   └── server/
│   │
│   └── utils/
│       ├── crypto.go
│       ├── time.go
│       └── convert.go
│
├── migrations/                       # Database migrations
│   ├── 001_create_users_table.sql
│   ├── 002_create_tokens_table.sql
│   └── 003_create_devices_table.sql
│
├── scripts/                          # Utility scripts
│   ├── setup.sh
│   ├── migrate.sh
│   └── seed.sh
│
├── test/                             # Integration tests
│   ├── integration/
│   └── fixtures/
│
├── .air.toml                         # Hot reload config
├── .env.example
├── go.mod
├── Makefile
└── README.md
```

### 3.2 ตัวอย่างโค้ดสำคัญ

#### Entity (internal/core/auth/entity/user.go)
```go
package entity

import (
    "time"
    "google/uuid"
)

type User struct {
    ID        uuid.UUID
    Email     string
    Password  string
    FirstName string
    LastName  string
    Role      Role
    IsActive  bool
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
}

type Role string

const (
    RoleAdmin Role = "admin"
    RoleUser  Role = "user"
    RoleGuest Role = "guest"
)

func (u *User) IsAdmin() bool {
    return u.Role == RoleAdmin
}

func (u *User) UpdateProfile(firstName, lastName string) {
    u.FirstName = firstName
    u.LastName = lastName
    u.UpdatedAt = time.Now()
}

func (u *User) Deactivate() {
    u.IsActive = false
    now := time.Now()
    u.DeletedAt = &now
}
```

#### Repository Interface (internal/core/auth/repository/user_repository.go)
```go
package repository

import (
    "context"
    "google/uuid"
    "your-project/internal/core/auth/entity"
)

type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id uuid.UUID) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
    FindByEmail(ctx context.Context, email string) (*entity.User, error)
    FindAll(ctx context.Context, limit, offset int) ([]*entity.User, int64, error)
    Exists(ctx context.Context, email string) (bool, error)
}
```

#### Repository Implementation (internal/platform/db/postgres/user_repository_impl.go)
```go
package postgres

import (
    "context"
    "errors"
    "time"
    
    "google/uuid"
    "gorm.io/gorm"
    "your-project/internal/core/auth/entity"
    "your-project/internal/core/auth/repository"
    "your-project/internal/core/auth/model"
)

type userRepositoryImpl struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
    return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
    userModel := &model.UserModel{
        ID:        user.ID,
        Email:     user.Email,
        Password:  user.Password,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        Role:      string(user.Role),
        IsActive:  user.IsActive,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }
    
    return r.db.WithContext(ctx).Create(userModel).Error
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
    var userModel model.UserModel
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&userModel).Error
    
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    return userModel.ToEntity(), nil
}
```

#### Domain Service (internal/core/auth/service/auth_service.go)
```go
package service

import (
    "context"
    "time"
    
    "google/uuid"
    "golang.org/x/crypto/bcrypt"
    "your-project/internal/core/auth/entity"
    "your-project/internal/core/auth/repository"
    "your-project/pkg/errors"
    "your-project/pkg/jwt"
)

type AuthService interface {
    Register(ctx context.Context, req *RegisterRequest) (*entity.User, error)
    Login(ctx context.Context, email, password string) (*AuthResponse, error)
    ValidateToken(ctx context.Context, token string) (*jwt.Claims, error)
    RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)
    Logout(ctx context.Context, userID uuid.UUID) error
}

type authService struct {
    userRepo    repository.UserRepository
    tokenRepo   repository.TokenRepository
    jwtService  jwt.JWTService
    cache       Cache
}

func NewAuthService(
    userRepo repository.UserRepository,
    tokenRepo repository.TokenRepository,
    jwtService jwt.JWTService,
    cache Cache,
) AuthService {
    return &authService{
        userRepo:   userRepo,
        tokenRepo:  tokenRepo,
        jwtService: jwtService,
        cache:      cache,
    }
}

func (s *authService) Register(ctx context.Context, req *RegisterRequest) (*entity.User, error) {
    // Check if user exists
    exists, err := s.userRepo.Exists(ctx, req.Email)
    if err != nil {
        return nil, errors.Wrap(err, "failed to check user existence")
    }
    if exists {
        return nil, errors.NewConflictError("user already exists")
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, errors.Wrap(err, "failed to hash password")
    }
    
    // Create user entity
    now := time.Now()
    user := &entity.User{
        ID:        uuid.New(),
        Email:     req.Email,
        Password:  string(hashedPassword),
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Role:      entity.RoleUser,
        IsActive:  true,
        CreatedAt: now,
        UpdatedAt: now,
    }
    
    // Save to database
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, errors.Wrap(err, "failed to create user")
    }
    
    return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
    // Get user
    user, err := s.userRepo.FindByEmail(ctx, email)
    if err != nil {
        return nil, errors.Wrap(err, "failed to find user")
    }
    if user == nil {
        return nil, errors.NewUnauthorizedError("invalid credentials")
    }
    
    // Check password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, errors.NewUnauthorizedError("invalid credentials")
    }
    
    // Generate tokens
    accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email, string(user.Role))
    if err != nil {
        return nil, errors.Wrap(err, "failed to generate access token")
    }
    
    refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
    if err != nil {
        return nil, errors.Wrap(err, "failed to generate refresh token")
    }
    
    // Save refresh token
    if err := s.tokenRepo.SaveRefreshToken(ctx, user.ID, refreshToken); err != nil {
        return nil, errors.Wrap(err, "failed to save refresh token")
    }
    
    return &AuthResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        TokenType:    "Bearer",
        ExpiresIn:    s.jwtService.GetAccessTokenExpiry(),
    }, nil
}
```

#### Handler (internal/core/auth/handler/auth_handler.go)
```go
package handler

import (
    "net/http"
    
    "gin-gonic/gin"
    "your-project/internal/core/auth/dto"
    "your-project/internal/core/auth/service"
    "your-project/internal/transport/httpx"
    "your-project/pkg/logger"
)

type AuthHandler struct {
    authService service.AuthService
    logger      logger.Logger
}

func NewAuthHandler(authService service.AuthService, logger logger.Logger) *AuthHandler {
    return &AuthHandler{
        authService: authService,
        logger:      logger,
    }
}

// Register godoc
// @Summary Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration request"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        httpx.RespondError(c, httpx.ErrBadRequest(err.Error()))
        return
    }
    
    if err := req.Validate(); err != nil {
        httpx.RespondError(c, httpx.ErrBadRequest(err.Error()))
        return
    }
    
    user, err := h.authService.Register(c.Request.Context(), &req)
    if err != nil {
        h.logger.Error("failed to register user", "error", err)
        httpx.RespondError(c, err)
        return
    }
    
    httpx.RespondSuccess(c, http.StatusCreated, dto.ToRegisterResponse(user))
}

// Login godoc
// @Summary User login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.LoginResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
    var req dto.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        httpx.RespondError(c, httpx.ErrBadRequest(err.Error()))
        return
    }
    
    authResp, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
    if err != nil {
        httpx.RespondError(c, err)
        return
    }
    
    httpx.RespondSuccess(c, http.StatusOK, dto.ToLoginResponse(authResp))
}
```

#### DTO (internal/core/auth/dto/auth_request.go)
```go
package dto

import (
    "regexp"
    "strings"
    
    "go-playground/validator/v10"
)

type RegisterRequest struct {
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=8"`
    FirstName string `json:"first_name" binding:"required"`
    LastName  string `json:"last_name" binding:"required"`
}

func (r *RegisterRequest) Validate() error {
    validate := validator.New()
    
    // Custom validation
    _ = validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
        password := fl.Field().String()
        // At least one uppercase, one lowercase, one number
        hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
        hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
        hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
        
        return hasUpper && hasLower && hasNumber
    })
    
    return validate.Struct(r)
}

func (r *RegisterRequest) Sanitize() {
    r.Email = strings.TrimSpace(strings.ToLower(r.Email))
    r.FirstName = strings.TrimSpace(r.FirstName)
    r.LastName = strings.TrimSpace(r.LastName)
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

func (r *LoginRequest) Sanitize() {
    r.Email = strings.TrimSpace(strings.ToLower(r.Email))
}
```

#### Middleware (internal/transport/middleware/auth_middleware.go)
```go
package middleware

import (
    "net/http"
    "strings"
    
    "gin-gonic/gin"
    "your-project/pkg/jwt"
    "your-project/internal/transport/httpx"
)

type AuthMiddleware struct {
    jwtService jwt.JWTService
}

func NewAuthMiddleware(jwtService jwt.JWTService) *AuthMiddleware {
    return &AuthMiddleware{
        jwtService: jwtService,
    }
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractToken(c)
        if token == "" {
            httpx.RespondError(c, httpx.ErrUnauthorized("missing authorization token"))
            c.Abort()
            return
        }
        
        claims, err := m.jwtService.ValidateToken(token)
        if err != nil {
            httpx.RespondError(c, httpx.ErrUnauthorized("invalid token"))
            c.Abort()
            return
        }
        
        // Set user info in context
        c.Set("user_id", claims.UserID)
        c.Set("user_email", claims.Email)
        c.Set("user_role", claims.Role)
        
        c.Next()
    }
}

func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("user_role")
        if !exists {
            httpx.RespondError(c, httpx.ErrForbidden("access denied"))
            c.Abort()
            return
        }
        
        for _, role := range roles {
            if role == userRole {
                c.Next()
                return
            }
        }
        
        httpx.RespondError(c, httpx.ErrForbidden("insufficient permissions"))
        c.Abort()
    }
}

func extractToken(c *gin.Context) string {
    // Check Authorization header
    authHeader := c.GetHeader("Authorization")
    if authHeader != "" {
        parts := strings.Split(authHeader, " ")
        if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
            return parts[1]
        }
    }
    
    // Check cookie
    token, err := c.Cookie("access_token")
    if err == nil {
        return token
    }
    
    return ""
}
```

#### Rate Limiting Middleware (internal/transport/middleware/rate_limit.go)
```go
package middleware

import (
    "net/http"
    "sync"
    "time"
    
    "gin-gonic/gin"
    "redis/go-redis/v9"
    "your-project/internal/transport/httpx"
)

type RateLimiter struct {
    redisClient *redis.Client
    limit       int
    window      time.Duration
    mu          sync.Mutex
}

func NewRateLimiter(redisClient *redis.Client, limit int, window time.Duration) *RateLimiter {
    return &RateLimiter{
        redisClient: redisClient,
        limit:       limit,
        window:      window,
    }
}

func (rl *RateLimiter) Limit() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get client IP
        clientIP := c.ClientIP()
        
        // Create key
        key := "rate_limit:" + clientIP
        
        // Use Redis for distributed rate limiting
        val, err := rl.redisClient.Incr(c.Request.Context(), key).Result()
        if err != nil {
            c.Next() // Fail open
            return
        }
        
        if val == 1 {
            // Set expiration on first request
            rl.redisClient.Expire(c.Request.Context(), key, rl.window)
        }
        
        if val > int64(rl.limit) {
            httpx.RespondError(c, httpx.ErrTooManyRequests("rate limit exceeded"))
            c.Abort()
            return
        }
        
        // Set headers
        c.Header("X-RateLimit-Limit", string(rl.limit))
        c.Header("X-RateLimit-Remaining", string(rl.limit-int(val)))
        
        c.Next()
    }
}
```

#### Logger Package (pkg/logger/logger.go)
```go
package logger

import (
    "context"
    "log/slog"
    "os"
    "time"
    
    "google/uuid"
)

type Logger interface {
    Debug(msg string, args ...any)
    Info(msg string, args ...any)
    Warn(msg string, args ...any)
    Error(msg string, args ...any)
    With(args ...any) Logger
    WithContext(ctx context.Context) Logger
}

type logger struct {
    slog *slog.Logger
}

type contextKey string

const (
    RequestIDKey contextKey = "request_id"
    TraceIDKey   contextKey = "trace_id"
    UserIDKey    contextKey = "user_id"
)

func New(level string) Logger {
    var slogLevel slog.Level
    switch level {
    case "debug":
        slogLevel = slog.LevelDebug
    case "info":
        slogLevel = slog.LevelInfo
    case "warn":
        slogLevel = slog.LevelWarn
    case "error":
        slogLevel = slog.LevelError
    default:
        slogLevel = slog.LevelInfo
    }
    
    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slogLevel,
        ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
            if a.Key == slog.TimeKey {
                a.Value = slog.StringValue(time.Now().Format(time.RFC3339Nano))
            }
            return a
        },
    })
    
    return &logger{
        slog: slog.New(handler),
    }
}

func (l *logger) Debug(msg string, args ...any) {
    l.slog.Debug(msg, args...)
}

func (l *logger) Info(msg string, args ...any) {
    l.slog.Info(msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
    l.slog.Warn(msg, args...)
}

func (l *logger) Error(msg string, args ...any) {
    l.slog.Error(msg, args...)
}

func (l *logger) With(args ...any) Logger {
    return &logger{
        slog: l.slog.With(args...),
    }
}

func (l *logger) WithContext(ctx context.Context) Logger {
    args := []any{}
    
    if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
        args = append(args, "request_id", requestID)
    }
    
    if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
        args = append(args, "trace_id", traceID)
    }
    
    if userID, ok := ctx.Value(UserIDKey).(string); ok {
        args = append(args, "user_id", userID)
    }
    
    return l.With(args...)
}
```

#### Error Handling (pkg/errors/errors.go)
```go
package errors

import (
    "fmt"
    "net/http"
)

type ErrorType string

const (
    ErrorTypeNotFound      ErrorType = "NOT_FOUND"
    ErrorTypeUnauthorized  ErrorType = "UNAUTHORIZED"
    ErrorTypeForbidden     ErrorType = "FORBIDDEN"
    ErrorTypeBadRequest    ErrorType = "BAD_REQUEST"
    ErrorTypeConflict      ErrorType = "CONFLICT"
    ErrorTypeInternal      ErrorType = "INTERNAL"
    ErrorTypeTooMany       ErrorType = "TOO_MANY_REQUESTS"
    ErrorTypeValidation    ErrorType = "VALIDATION"
)

type AppError struct {
    Type    ErrorType
    Message string
    Code    string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %s: %v", e.Type, e.Message, e.Err)
    }
    return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
    return e.Err
}

func (e *AppError) StatusCode() int {
    switch e.Type {
    case ErrorTypeNotFound:
        return http.StatusNotFound
    case ErrorTypeUnauthorized:
        return http.StatusUnauthorized
    case ErrorTypeForbidden:
        return http.StatusForbidden
    case ErrorTypeBadRequest:
        return http.StatusBadRequest
    case ErrorTypeConflict:
        return http.StatusConflict
    case ErrorTypeTooMany:
        return http.StatusTooManyRequests
    case ErrorTypeValidation:
        return http.StatusUnprocessableEntity
    default:
        return http.StatusInternalServerError
    }
}

func NewNotFoundError(message string) *AppError {
    return &AppError{
        Type:    ErrorTypeNotFound,
        Message: message,
    }
}

func NewUnauthorizedError(message string) *AppError {
    return &AppError{
        Type:    ErrorTypeUnauthorized,
        Message: message,
    }
}

func NewForbiddenError(message string) *AppError {
    return &AppError{
        Type:    ErrorTypeForbidden,
        Message: message,
    }
}

func NewBadRequestError(message string) *AppError {
    return &AppError{
        Type:    ErrorTypeBadRequest,
        Message: message,
    }
}

func NewConflictError(message string) *AppError {
    return &AppError{
        Type:    ErrorTypeConflict,
        Message: message,
    }
}

func NewInternalError(err error) *AppError {
    return &AppError{
        Type:    ErrorTypeInternal,
        Message: "internal server error",
        Err:     err,
    }
}

func Wrap(err error, message string) *AppError {
    if err == nil {
        return nil
    }
    
    if appErr, ok := err.(*AppError); ok {
        return &AppError{
            Type:    appErr.Type,
            Message: message,
            Err:     appErr,
        }
    }
    
    return &AppError{
        Type:    ErrorTypeInternal,
        Message: message,
        Err:     err,
    }
}
```

#### Cache Implementation (pkg/cache/redis_cache.go)
```go
package cache

import (
    "context"
    "encoding/json"
    "time"
    
    "redis/go-redis/v9"
)

type Cache interface {
    Get(ctx context.Context, key string, dest interface{}) error
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
    Incr(ctx context.Context, key string) (int64, error)
    Expire(ctx context.Context, key string, ttl time.Duration) error
    GetOrSet(ctx context.Context, key string, dest interface{}, ttl time.Duration, fn func() (interface{}, error)) error
}

type redisCache struct {
    client *redis.Client
}

func NewRedisCache(client *redis.Client) Cache {
    return &redisCache{
        client: client,
    }
}

func (c *redisCache) Get(ctx context.Context, key string, dest interface{}) error {
    val, err := c.client.Get(ctx, key).Bytes()
    if err == redis.Nil {
        return nil
    }
    if err != nil {
        return err
    }
    
    return json.Unmarshal(val, dest)
}

func (c *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    
    return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
    return c.client.Del(ctx, key).Err()
}

func (c *redisCache) Exists(ctx context.Context, key string) (bool, error) {
    n, err := c.client.Exists(ctx, key).Result()
    return n > 0, err
}

func (c *redisCache) Incr(ctx context.Context, key string) (int64, error) {
    return c.client.Incr(ctx, key).Result()
}

func (c *redisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
    return c.client.Expire(ctx, key, ttl).Err()
}

func (c *redisCache) GetOrSet(ctx context.Context, key string, dest interface{}, ttl time.Duration, fn func() (interface{}, error)) error {
    // Try to get from cache
    err := c.Get(ctx, key, dest)
    if err == nil {
        return nil
    }
    
    // If not found, execute function
    val, err := fn()
    if err != nil {
        return err
    }
    
    // Store in cache
    if err := c.Set(ctx, key, val, ttl); err != nil {
        return err
    }
    
    // Convert to dest
    data, err := json.Marshal(val)
    if err != nil {
        return err
    }
    
    return json.Unmarshal(data, dest)
}
```

#### Queue Implementation (pkg/queue/redis_queue.go)
```go
package queue

import (
    "context"
    "encoding/json"
    "time"
    
    "redis/go-redis/v9"
    "google/uuid"
)

type Job struct {
    ID        string          `json:"id"`
    Type      string          `json:"type"`
    Payload   json.RawMessage `json:"payload"`
    Attempts  int             `json:"attempts"`
    MaxAttempts int           `json:"max_attempts"`
    CreatedAt time.Time       `json:"created_at"`
    UpdatedAt time.Time       `json:"updated_at"`
}

type Queue interface {
    Push(ctx context.Context, queue string, job *Job) error
    Pop(ctx context.Context, queue string) (*Job, error)
    Ack(ctx context.Context, queue string, jobID string) error
    Nack(ctx context.Context, queue string, jobID string) error
    Len(ctx context.Context, queue string) (int64, error)
}

type redisQueue struct {
    client *redis.Client
}

func NewRedisQueue(client *redis.Client) Queue {
    return &redisQueue{
        client: client,
    }
}

func (q *redisQueue) Push(ctx context.Context, queue string, job *Job) error {
    if job.ID == "" {
        job.ID = uuid.New().String()
    }
    if job.CreatedAt.IsZero() {
        job.CreatedAt = time.Now()
    }
    job.UpdatedAt = time.Now()
    
    data, err := json.Marshal(job)
    if err != nil {
        return err
    }
    
    // Push to main queue
    return q.client.RPush(ctx, "queue:"+queue, data).Err()
}

func (q *redisQueue) Pop(ctx context.Context, queue string) (*Job, error) {
    // Pop from main queue and push to processing queue atomically
    result, err := q.client.LPop(ctx, "queue:"+queue).Bytes()
    if err == redis.Nil {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    var job Job
    if err := json.Unmarshal(result, &job); err != nil {
        return nil, err
    }
    
    // Store in processing queue with TTL
    processingKey := "processing:" + queue + ":" + job.ID
    if err := q.client.Set(ctx, processingKey, result, 5*time.Minute).Err(); err != nil {
        return nil, err
    }
    
    return &job, nil
}

func (q *redisQueue) Ack(ctx context.Context, queue string, jobID string) error {
    // Remove from processing queue
    return q.client.Del(ctx, "processing:"+queue+":"+jobID).Err()
}

func (q *redisQueue) Nack(ctx context.Context, queue string, jobID string) error {
    // Get job from processing queue
    data, err := q.client.GetDel(ctx, "processing:"+queue+":"+jobID).Bytes()
    if err != nil {
        return err
    }
    
    var job Job
    if err := json.Unmarshal(data, &job); err != nil {
        return err
    }
    
    // Increment attempts
    job.Attempts++
    job.UpdatedAt = time.Now()
    
    if job.Attempts >= job.MaxAttempts {
        // Move to dead letter queue
        return q.client.RPush(ctx, "dead:"+queue, data).Err()
    }
    
    // Push back to main queue
    return q.Push(ctx, queue, &job)
}

func (q *redisQueue) Len(ctx context.Context, queue string) (int64, error) {
    return q.client.LLen(ctx, "queue:"+queue).Result()
}
```

#### Worker Pool (internal/platform/queue/worker_pool.go)
```go
package queue

import (
    "context"
    "sync"
    "time"
    
    "your-project/pkg/logger"
    "your-project/pkg/queue"
)

type WorkerFunc func(ctx context.Context, job *queue.Job) error

type WorkerPool struct {
    queue      queue.Queue
    workers    int
    jobQueue   string
    handler    WorkerFunc
    logger     logger.Logger
    wg         sync.WaitGroup
    stopChan   chan struct{}
}

func NewWorkerPool(
    queue queue.Queue,
    workers int,
    jobQueue string,
    handler WorkerFunc,
    logger logger.Logger,
) *WorkerPool {
    return &WorkerPool{
        queue:    queue,
        workers:  workers,
        jobQueue: jobQueue,
        handler:  handler,
        logger:   logger,
        stopChan: make(chan struct{}),
    }
}

func (wp *WorkerPool) Start(ctx context.Context) {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker(i)
    }
}

func (wp *WorkerPool) Stop() {
    close(wp.stopChan)
    wp.wg.Wait()
}

func (wp *WorkerPool) worker(id int) {
    defer wp.wg.Done()
    
    ctx := context.Background()
    
    for {
        select {
        case <-wp.stopChan:
            return
        default:
            job, err := wp.queue.Pop(ctx, wp.jobQueue)
            if err != nil {
                wp.logger.Error("failed to pop job", "error", err, "worker", id)
                time.Sleep(time.Second)
                continue
            }
            
            if job == nil {
                time.Sleep(100 * time.Millisecond)
                continue
            }
            
            // Process job
            if err := wp.handler(ctx, job); err != nil {
                wp.logger.Error("job processing failed", 
                    "error", err, 
                    "job_id", job.ID,
                    "worker", id,
                )
                
                // Negative acknowledgment
                if err := wp.queue.Nack(ctx, wp.jobQueue, job.ID); err != nil {
                    wp.logger.Error("failed to nack job", "error", err)
                }
            } else {
                // Positive acknowledgment
                if err := wp.queue.Ack(ctx, wp.jobQueue, job.ID); err != nil {
                    wp.logger.Error("failed to ack job", "error", err)
                }
            }
        }
    }
}
```

#### JWT Service (pkg/jwt/jwt_service.go)
```go
package jwt

import (
    "time"
    
    "golang-jwt/jwt/v5"
    "google/uuid"
)

type Claims struct {
    UserID    uuid.UUID `json:"user_id"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    jwt.RegisteredClaims
}

type JWTService interface {
    GenerateAccessToken(userID uuid.UUID, email, role string) (string, error)
    GenerateRefreshToken(userID uuid.UUID) (string, error)
    ValidateToken(token string) (*Claims, error)
    GetAccessTokenExpiry() int64
}

type jwtService struct {
    secretKey     []byte
    accessExpiry  time.Duration
    refreshExpiry time.Duration
    issuer        string
}

func NewJWTService(secretKey string, accessExpiry, refreshExpiry time.Duration, issuer string) JWTService {
    return &jwtService{
        secretKey:     []byte(secretKey),
        accessExpiry:  accessExpiry,
        refreshExpiry: refreshExpiry,
        issuer:        issuer,
    }
}

func (s *jwtService) GenerateAccessToken(userID uuid.UUID, email, role string) (string, error) {
    now := time.Now()
    claims := &Claims{
        UserID: userID,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpiry)),
            IssuedAt:  jwt.NewNumericDate(now),
            NotBefore: jwt.NewNumericDate(now),
            Issuer:    s.issuer,
            Subject:   userID.String(),
            ID:        uuid.New().String(),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.secretKey)
}

func (s *jwtService) GenerateRefreshToken(userID uuid.UUID) (string, error) {
    now := time.Now()
    claims := &jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshExpiry)),
        IssuedAt:  jwt.NewNumericDate(now),
        NotBefore: jwt.NewNumericDate(now),
        Issuer:    s.issuer,
        Subject:   userID.String(),
        ID:        uuid.New().String(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.secretKey)
}

func (s *jwtService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return s.secretKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, jwt.ErrTokenInvalidClaims
}

func (s *jwtService) GetAccessTokenExpiry() int64 {
    return int64(s.accessExpiry.Seconds())
}
```

#### Transaction Manager (pkg/transaction/transaction_manager.go)
```go
package transaction

import (
    "context"
    "database/sql"
    
    "gorm.io/gorm"
)

type TransactionManager interface {
    Begin(ctx context.Context) (context.Context, error)
    Commit(ctx context.Context) error
    Rollback(ctx context.Context) error
    GetDB(ctx context.Context) *gorm.DB
}

type transactionManager struct {
    db *gorm.DB
}

type txKey struct{}

func NewTransactionManager(db *gorm.DB) TransactionManager {
    return &transactionManager{
        db: db,
    }
}

func (tm *transactionManager) Begin(ctx context.Context) (context.Context, error) {
    tx := tm.db.Begin()
    if tx.Error != nil {
        return ctx, tx.Error
    }
    
    return context.WithValue(ctx, txKey{}, tx), nil
}

func (tm *transactionManager) Commit(ctx context.Context) error {
    tx, ok := ctx.Value(txKey{}).(*gorm.DB)
    if !ok {
        return sql.ErrTxDone
    }
    
    return tx.Commit().Error
}

func (tm *transactionManager) Rollback(ctx context.Context) error {
    tx, ok := ctx.Value(txKey{}).(*gorm.DB)
    if !ok {
        return sql.ErrTxDone
    }
    
    return tx.Rollback().Error
}

func (tm *transactionManager) GetDB(ctx context.Context) *gorm.DB {
    tx, ok := ctx.Value(txKey{}).(*gorm.DB)
    if ok {
        return tx
    }
    
    return tm.db.WithContext(ctx)
}
```

#### Main Application (cmd/app/main.go)
```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "gin-gonic/gin"
    "redis/go-redis/v9"
    "gorm.io/gorm"
    
    "your-project/configs"
    "your-project/internal/platform/config"
    "your-project/internal/platform/db/postgres"
    "your-project/internal/platform/logger"
    "your-project/internal/platform/cache"
    "your-project/internal/platform/queue"
    "your-project/internal/transport/middleware"
    "your-project/internal/core/auth/handler"
    "your-project/internal/core/auth/repository"
    "your-project/internal/core/auth/service"
    authRepo "your-project/internal/platform/db/postgres/repository"
    "your-project/pkg/jwt"
    "your-project/pkg/transaction"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    // Initialize logger
    appLogger := logger.New(cfg.LogLevel)
    appLogger.Info("Starting application", "version", cfg.Version)
    
    // Initialize database
    db, err := postgres.NewConnection(cfg.Database)
    if err != nil {
        appLogger.Error("Failed to connect to database", "error", err)
        os.Exit(1)
    }
    
    // Initialize Redis
    redisClient := redis.NewClient(&redis.Options{
        Addr:     cfg.Redis.Addr,
        Password: cfg.Redis.Password,
        DB:       cfg.Redis.DB,
    })
    
    if err := redisClient.Ping(context.Background()).Err(); err != nil {
        appLogger.Error("Failed to connect to Redis", "error", err)
        os.Exit(1)
    }
    
    // Initialize components
    cacheService := cache.NewRedisCache(redisClient)
    queueService := queue.NewRedisQueue(redisClient)
    jwtService := jwt.NewJWTService(
        cfg.JWT.Secret,
        cfg.JWT.AccessExpiry,
        cfg.JWT.RefreshExpiry,
        cfg.JWT.Issuer,
    )
    txManager := transaction.NewTransactionManager(db)
    
    // Initialize repositories
    userRepo := authRepo.NewUserRepository(db)
    tokenRepo := authRepo.NewTokenRepository(db)
    
    // Initialize services
    authService := service.NewAuthService(
        userRepo,
        tokenRepo,
        jwtService,
        cacheService,
    )
    
    // Initialize handlers
    authHandler := handler.NewAuthHandler(authService, appLogger)
    
    // Initialize middleware
    authMiddleware := middleware.NewAuthMiddleware(jwtService)
    rateLimiter := middleware.NewRateLimiter(redisClient, 100, time.Minute)
    
    // Setup router
    router := gin.New()
    router.Use(
        gin.Recovery(),
        middleware.LoggingMiddleware(appLogger),
        middleware.SecurityHeaders(),
        middleware.CORS(cfg.CORS),
        rateLimiter.Limit(),
    )
    
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // API v1 routes
    v1 := router.Group("/api/v1")
    {
        // Public routes
        auth := v1.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)
            auth.POST("/refresh", authHandler.Refresh)
        }
        
        // Protected routes
        protected := v1.Group("")
        protected.Use(authMiddleware.RequireAuth())
        {
            protected.POST("/auth/logout", authHandler.Logout)
            protected.GET("/auth/me", authHandler.GetMe)
        }
        
        // Admin routes
        admin := v1.Group("/admin")
        admin.Use(authMiddleware.RequireAuth())
        admin.Use(authMiddleware.RequireRole("admin"))
        {
            // Admin routes here
        }
    }
    
    // Start server
    srv := &http.Server{
        Addr:    ":" + cfg.Port,
        Handler: router,
    }
    
    go func() {
        appLogger.Info("Server starting", "port", cfg.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            appLogger.Error("Server failed", "error", err)
            os.Exit(1)
        }
    }()
    
    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    appLogger.Info("Shutting down server...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        appLogger.Error("Server forced to shutdown", "error", err)
    }
    
    appLogger.Info("Server exited")
}
```

### 3.3 Configuration (configs/config.yaml)
```yaml
server:
  port: 8080
  mode: development
  version: 1.0.0

database:
  driver: postgres
  host: localhost
  port: 5432
  username: postgres
  password: postgres
  database: myapp
  sslmode: disable
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 1h

redis:
  addr: localhost:6379
  password: ""
  db: 0
  pool_size: 10

jwt:
  secret: your-secret-key
  access_expiry: 15m
  refresh_expiry: 168h
  issuer: myapp

cors:
  allowed_origins:
    - http://localhost:3000
  allowed_methods:
    - GET
    - POST
    - PUT
    - DELETE
    - OPTIONS
  allowed_headers:
    - Content-Type
    - Authorization

log:
  level: info
  format: json
```

### 3.4 Makefile
```makefile
.PHONY: build run test migrate clean docker-up docker-down

build:
    go build -o bin/app cmd/app/main.go

run:
    go run cmd/app/main.go

test:
    go test -v ./...

test-coverage:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out

migrate-up:
    migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/myapp?sslmode=disable" up

migrate-down:
    migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/myapp?sslmode=disable" down

migrate-create:
    migrate create -ext sql -dir migrations -seq $(name)

docker-up:
    docker-compose up -d

docker-down:
    docker-compose down

lint:
    golangci-lint run

swagger:
    swag init -g cmd/app/main.go -o api/app

clean:
    rm -rf bin/
    rm -f coverage.out
```

---

## 4. สถาปัตยกรรม LangChain FastAPI

### 4.1 โครงสร้างโปรเจค LangChain FastAPI

```
langchain-project/
├── app/
│   ├── __init__.py
│   ├── main.py                    # Entry point
│   │
│   ├── core/                       # Domain Layer
│   │   ├── __init__.py
│   │   ├── entities/
│   │   │   ├── __init__.py
│   │   │   ├── document.py
│   │   │   ├── conversation.py
│   │   │   └── message.py
│   │   │
│   │   ├── interfaces/
│   │   │   ├── __init__.py
│   │   │   ├── repository.py
│   │   │   └── llm_service.py
│   │   │
│   │   └── value_objects/
│   │       ├── __init__.py
│   │       └── embedding.py
│   │
│   ├── application/                 # Application Layer
│   │   ├── __init__.py
│   │   ├── services/
│   │   │   ├── __init__.py
│   │   │   ├── chat_service.py
│   │   │   ├── document_service.py
│   │   │   ├── embedding_service.py
│   │   │   └── agent_service.py
│   │   │
│   │   ├── dtos/
│   │   │   ├── __init__.py
│   │   │   ├── chat_dto.py
│   │   │   └── document_dto.py
│   │   │
│   │   └── use_cases/
│   │       ├── __init__.py
│   │       ├── chat_use_case.py
│   │       └── qa_use_case.py
│   │
│   ├── infrastructure/              # Infrastructure Layer
│   │   ├── __init__.py
│   │   ├── database/
│   │   │   ├── __init__.py
│   │   │   ├── postgres.py
│   │   │   ├── redis.py
│   │   │   └── repositories/
│   │   │       ├── __init__.py
│   │   │       ├── document_repository.py
│   │   │       └── conversation_repository.py
│   │   │
│   │   ├── llm/
│   │   │   ├── __init__.py
│   │   │   ├── openai_service.py
│   │   │   ├── huggingface_service.py
│   │   │   └── chain_factory.py
│   │   │
│   │   ├── vector_store/
│   │   │   ├── __init__.py
│   │   │   ├── chroma_store.py
│   │   │   ├── pinecone_store.py
│   │   │   └── qdrant_store.py
│   │   │
│   │   ├── cache/
│   │   │   ├── __init__.py
│   │   │   └── redis_cache.py
│   │   │
│   │   └── queue/
│   │       ├── __init__.py
│   │       ├── redis_queue.py
│   │       └── tasks.py
│   │
│   ├── api/                          # Interface Layer
│   │   ├── __init__.py
│   │   ├── dependencies/
│   │   │   ├── __init__.py
│   │   │   └── container.py
│   │   │
│   │   ├── middleware/
│   │   │   ├── __init__.py
│   │   │   ├── auth.py
│   │   │   ├── logging.py
│   │   │   └── rate_limit.py
│   │   │
│   │   ├── routes/
│   │   │   ├── __init__.py
│   │   │   ├── chat_routes.py
│   │   │   ├── document_routes.py
│   │   │   └── health_routes.py
│   │   │
│   │   └── models/
│   │       ├── __init__.py
│   │       ├── request_models.py
│   │       └── response_models.py
│   │
│   ├── config/
│   │   ├── __init__.py
│   │   ├── settings.py
│   │   └── logging_config.py
│   │
│   └── utils/
│       ├── __init__.py
│       ├── errors.py
│       ├── logger.py
│       └── helpers.py
│
├── migrations/
│   ├── versions/
│   └── env.py
│
├── tests/
│   ├── __init__.py
│   ├── unit/
│   ├── integration/
│   └── conftest.py
│
├── scripts/
│   ├── seed_data.py
│   └── migrate_db.py
│
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yml
│
├── .env.example
├── .gitignore
├── requirements.txt
├── requirements-dev.txt
├── alembic.ini
└── README.md
```

### 4.2 ตัวอย่างโค้ด LangChain

#### Entity (app/core/entities/document.py)
```python
from dataclasses import dataclass
from datetime import datetime
from typing import Optional, Dict, Any
from uuid import UUID, uuid4

@dataclass
class Document:
    id: UUID
    title: str
    content: str
    metadata: Dict[str, Any]
    embedding: Optional[list[float]]
    created_at: datetime
    updated_at: datetime
    user_id: Optional[UUID]
    
    @classmethod
    def create(cls, title: str, content: str, user_id: Optional[UUID] = None) -> "Document":
        now = datetime.utcnow()
        return cls(
            id=uuid4(),
            title=title,
            content=content,
            metadata={},
            embedding=None,
            created_at=now,
            updated_at=now,
            user_id=user_id
        )
    
    def update_embedding(self, embedding: list[float]) -> None:
        self.embedding = embedding
        self.updated_at = datetime.utcnow()
    
    def add_metadata(self, key: str, value: Any) -> None:
        self.metadata[key] = value
        self.updated_at = datetime.utcnow()
```

#### Service (app/application/services/chat_service.py)
```python
from typing import List, Optional, Dict, Any
from uuid import UUID
from langchain.chains import ConversationChain
from langchain.memory import ConversationBufferMemory
from langchain.schema import BaseMessage, HumanMessage, AIMessage

from app.core.entities.conversation import Conversation, Message
from app.core.interfaces.repository import ConversationRepository
from app.core.interfaces.llm_service import LLMService
from app.application.dtos.chat_dto import ChatRequest, ChatResponse
from app.infrastructure.cache.redis_cache import RedisCache


class ChatService:
    def __init__(
        self,
        conversation_repo: ConversationRepository,
        llm_service: LLMService,
        cache: RedisCache
    ):
        self.conversation_repo = conversation_repo
        self.llm_service = llm_service
        self.cache = cache
        self.chains: Dict[str, ConversationChain] = {}
    
    async def chat(
        self,
        request: ChatRequest,
        user_id: Optional[UUID] = None
    ) -> ChatResponse:
        # Get or create conversation
        if request.conversation_id:
            conversation = await self.conversation_repo.get_by_id(
                request.conversation_id
            )
        else:
            conversation = Conversation.create(user_id=user_id)
        
        # Create message
        user_message = Message(
            role="user",
            content=request.message,
            metadata=request.metadata
        )
        conversation.add_message(user_message)
        
        # Get or create chain
        chain = await self._get_or_create_chain(conversation.id)
        
        # Generate response
        response = await chain.acall(inputs={"input": request.message})
        
        # Create AI message
        ai_message = Message(
            role="assistant",
            content=response["response"],
            metadata={"tokens": response.get("tokens", {})}
        )
        conversation.add_message(ai_message)
        
        # Save conversation
        await self.conversation_repo.save(conversation)
        
        return ChatResponse(
            conversation_id=conversation.id,
            message=ai_message.content,
            metadata=ai_message.metadata,
            created_at=ai_message.created_at
        )
    
    async def _get_or_create_chain(self, conversation_id: UUID) -> ConversationChain:
        # Check cache
        cache_key = f"chain:{conversation_id}"
        cached_chain = await self.cache.get(cache_key)
        if cached_chain:
            return cached_chain
        
        # Get conversation history
        conversation = await self.conversation_repo.get_by_id(conversation_id)
        
        # Create memory with history
        memory = ConversationBufferMemory()
        for msg in conversation.messages[-10:]:  # Last 10 messages
            if msg.role == "user":
                memory.chat_memory.add_user_message(msg.content)
            else:
                memory.chat_memory.add_ai_message(msg.content)
        
        # Create chain
        chain = ConversationChain(
            llm=self.llm_service.get_llm(),
            memory=memory,
            verbose=True
        )
        
        # Cache chain
        await self.cache.set(cache_key, chain, ttl=3600)
        
        return chain
    
    async def get_conversation_history(
        self,
        conversation_id: UUID,
        limit: int = 50
    ) -> List[Dict[str, Any]]:
        conversation = await self.conversation_repo.get_by_id(conversation_id)
        
        return [
            {
                "role": msg.role,
                "content": msg.content,
                "created_at": msg.created_at.isoformat(),
                "metadata": msg.metadata
            }
            for msg in conversation.messages[-limit:]
        ]
```

#### Document Service with LangChain (app/application/services/document_service.py)
```python
from typing import List, Optional, Dict, Any
from uuid import UUID
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain.embeddings import OpenAIEmbeddings
from langchain.vectorstores import Chroma
from langchain.schema import Document as LangChainDocument

from app.core.entities.document import Document
from app.core.interfaces.repository import DocumentRepository
from app.infrastructure.vector_store.chroma_store import ChromaStore


class DocumentService:
    def __init__(
        self,
        document_repo: DocumentRepository,
        vector_store: ChromaStore,
        embeddings: OpenAIEmbeddings
    ):
        self.document_repo = document_repo
        self.vector_store = vector_store
        self.embeddings = embeddings
        self.text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=1000,
            chunk_overlap=200,
            length_function=len,
            separators=["\n\n", "\n", " ", ""]
        )
    
    async def create_document(
        self,
        title: str,
        content: str,
        user_id: Optional[UUID] = None,
        metadata: Optional[Dict[str, Any]] = None
    ) -> Document:
        # Create document entity
        document = Document.create(
            title=title,
            content=content,
            user_id=user_id
        )
        
        if metadata:
            document.metadata.update(metadata)
        
        # Split document into chunks
        chunks = self.text_splitter.split_text(content)
        
        # Create LangChain documents
        lc_documents = [
            LangChainDocument(
                page_content=chunk,
                metadata={
                    "document_id": str(document.id),
                    "chunk_index": i,
                    "title": title,
                    **metadata
                }
            )
            for i, chunk in enumerate(chunks)
        ]
        
        # Generate embeddings and store in vector store
        await self.vector_store.add_documents(lc_documents)
        
        # Save to database
        await self.document_repo.save(document)
        
        return document
    
    async def search_similar(
        self,
        query: str,
        limit: int = 5,
        user_id: Optional[UUID] = None
    ) -> List[Dict[str, Any]]:
        # Search in vector store
        results = await self.vector_store.similarity_search(
            query=query,
            k=limit,
            filter={"user_id": str(user_id)} if user_id else None
        )
        
        return [
            {
                "content": doc.page_content,
                "metadata": doc.metadata,
                "score": score
            }
            for doc, score in results
        ]
    
    async def get_document_with_context(
        self,
        document_id: UUID,
        query: str
    ) -> Dict[str, Any]:
        # Get document from database
        document = await self.document_repo.get_by_id(document_id)
        
        # Search for relevant chunks
        similar = await self.vector_store.similarity_search(
            query=query,
            k=3,
            filter={"document_id": str(document_id)}
        )
        
        context = "\n\n".join([doc.page_content for doc, _ in similar])
        
        return {
            "document": document,
            "context": context,
            "relevant_chunks": [
                {
                    "content": doc.page_content,
                    "score": score
                }
                for doc, score in similar
            ]
        }
```

#### Agent Service (app/application/services/agent_service.py)
```python
from typing import List, Dict, Any, Optional
from langchain.agents import Tool, AgentExecutor, create_react_agent
from langchain.prompts import PromptTemplate
from langchain.chains import LLMMathChain
from langchain.tools import DuckDuckGoSearchRun

from app.core.interfaces.llm_service import LLMService
from app.application.services.document_service import DocumentService


class AgentService:
    def __init__(
        self,
        llm_service: LLMService,
        document_service: DocumentService
    ):
        self.llm_service = llm_service
        self.document_service = document_service
        self.agent_executor = self._create_agent()
    
    def _create_agent(self) -> AgentExecutor:
        # Create tools
        tools = [
            Tool(
                name="Search",
                func=DuckDuckGoSearchRun().run,
                description="Search the internet for current information"
            ),
            Tool(
                name="Calculator",
                func=LLMMathChain.from_llm(
                    self.llm_service.get_llm()
                ).run,
                description="Perform mathematical calculations"
            ),
            Tool(
                name="DocumentSearch",
                func=self._search_documents,
                description="Search through internal documents"
            )
        ]
        
        # Create prompt
        prompt = PromptTemplate.from_template("""
        You are a helpful AI assistant with access to various tools.
        Use the following tools to help answer the user's question:
        
        {tools}
        
        Question: {input}
        
        Thought process: {agent_scratchpad}
        """)
        
        # Create agent
        agent = create_react_agent(
            llm=self.llm_service.get_llm(),
            tools=tools,
            prompt=prompt
        )
        
        return AgentExecutor(
            agent=agent,
            tools=tools,
            verbose=True,
            handle_parsing_errors=True,
            max_iterations=5
        )
    
    async def _search_documents(self, query: str) -> str:
        """Search through internal documents"""
        results = await self.document_service.search_similar(
            query=query,
            limit=3
        )
        
        if not results:
            return "No relevant documents found."
        
        return "\n\n".join([
            f"[Document: {r['metadata'].get('title', 'Unknown')}]\n{r['content']}"
            for r in results
        ])
    
    async def run_agent(
        self,
        query: str,
        user_id: Optional[str] = None
    ) -> Dict[str, Any]:
        """Run the agent with a query"""
        result = await self.agent_executor.arun(
            input=query,
            user_id=user_id
        )
        
        return {
            "response": result,
            "agent_used": True
        }
```

#### DTOs (app/application/dtos/chat_dto.py)
```python
from pydantic import BaseModel, Field
from typing import Optional, Dict, Any
from uuid import UUID
from datetime import datetime


class ChatRequest(BaseModel):
    message: str = Field(..., min_length=1, max_length=4000)
    conversation_id: Optional[UUID] = None
    metadata: Optional[Dict[str, Any]] = Field(default_factory=dict)
    stream: bool = False
    
    class Config:
        json_schema_extra = {
            "example": {
                "message": "What is the capital of France?",
                "metadata": {"source": "web"}
            }
        }


class ChatResponse(BaseModel):
    conversation_id: UUID
    message: str
    metadata: Dict[str, Any]
    created_at: datetime


class ConversationSummary(BaseModel):
    id: UUID
    title: Optional[str]
    message_count: int
    last_message_at: datetime
    created_at: datetime
```

#### Vector Store (app/infrastructure/vector_store/chroma_store.py)
```python
from typing import List, Optional, Dict, Any
import chromadb
from chromadb.config import Settings
from langchain.vectorstores import Chroma
from langchain.embeddings import OpenAIEmbeddings
from langchain.schema import Document


class ChromaStore:
    def __init__(
        self,
        collection_name: str,
        persist_directory: str,
        embeddings: OpenAIEmbeddings
    ):
        self.client = chromadb.Client(Settings(
            chroma_db_impl="duckdb+parquet",
            persist_directory=persist_directory
        ))
        
        self.vectorstore = Chroma(
            collection_name=collection_name,
            embedding_function=embeddings,
            client=self.client,
            persist_directory=persist_directory
        )
    
    async def add_documents(
        self,
        documents: List[Document],
        ids: Optional[List[str]] = None
    ) -> List[str]:
        """Add documents to vector store"""
        return self.vectorstore.add_documents(
            documents=documents,
            ids=ids
        )
    
    async def similarity_search(
        self,
        query: str,
        k: int = 4,
        filter: Optional[Dict[str, Any]] = None
    ) -> List[tuple[Document, float]]:
        """Search for similar documents"""
        results = self.vectorstore.similarity_search_with_score(
            query=query,
            k=k,
            filter=filter
        )
        return results
    
    async def delete_document(self, document_id: str) -> None:
        """Delete document by ID"""
        self.vectorstore.delete(filter={"document_id": document_id})
    
    async def get_collection_stats(self) -> Dict[str, Any]:
        """Get collection statistics"""
        collection = self.client.get_collection(
            self.vectorstore._collection.name
        )
        return {
            "count": collection.count(),
            "name": collection.name,
            "metadata": collection.metadata
        }
```

#### LLM Service (app/infrastructure/llm/openai_service.py)
```python
from typing import Optional, Dict, Any
from langchain.chat_models import ChatOpenAI
from langchain.llms import OpenAI
from langchain.schema import BaseMessage, HumanMessage, AIMessage
from tenacity import retry, stop_after_attempt, wait_exponential

from app.core.interfaces.llm_service import LLMService
from app.config.settings import Settings


class OpenAIService(LLMService):
    def __init__(self, settings: Settings):
        self.settings = settings
        self.chat_model = ChatOpenAI(
            model=settings.OPENAI_MODEL,
            temperature=settings.OPENAI_TEMPERATURE,
            max_tokens=settings.OPENAI_MAX_TOKENS,
            openai_api_key=settings.OPENAI_API_KEY
        )
        self.llm = OpenAI(
            model=settings.OPENAI_MODEL,
            temperature=settings.OPENAI_TEMPERATURE,
            openai_api_key=settings.OPENAI_API_KEY
        )
    
    @retry(
        stop=stop_after_attempt(3),
        wait=wait_exponential(multiplier=1, min=4, max=10)
    )
    async def generate_response(
        self,
        messages: List[BaseMessage],
        temperature: Optional[float] = None
    ) -> str:
        """Generate response from messages"""
        response = await self.chat_model.agenerate([messages])
        return response.generations[0][0].text
    
    @retry(
        stop=stop_after_attempt(3),
        wait=wait_exponential(multiplier=1, min=4, max=10)
    )
    async def generate_completion(
        self,
        prompt: str,
        max_tokens: Optional[int] = None
    ) -> str:
        """Generate completion from prompt"""
        response = await self.llm.agenerate([prompt])
        return response.generations[0][0].text
    
    async def stream_response(
        self,
        messages: List[BaseMessage]
    ) -> Any:
        """Stream response"""
        async for chunk in self.chat_model.astream(messages):
            yield chunk
    
    def get_llm(self):
        """Get LangChain LLM instance"""
        return self.chat_model
    
    def get_embeddings(self):
        """Get embeddings model"""
        from langchain.embeddings import OpenAIEmbeddings
        return OpenAIEmbeddings(
            openai_api_key=self.settings.OPENAI_API_KEY
        )
```

#### API Routes (app/api/routes/chat_routes.py)
```python
from fastapi import APIRouter, Depends, HTTPException, BackgroundTasks
from typing import List
from uuid import UUID

from app.api.dependencies.container import get_chat_service
from app.api.models.request_models import ChatRequestModel
from app.api.models.response_models import ChatResponseModel, ConversationModel
from app.application.services.chat_service import ChatService
from app.api.middleware.auth import get_current_user

router = APIRouter(prefix="/chat", tags=["chat"])


@router.post("/", response_model=ChatResponseModel)
async def chat(
    request: ChatRequestModel,
    background_tasks: BackgroundTasks,
    chat_service: ChatService = Depends(get_chat_service),
    user_id: Optional[UUID] = Depends(get_current_user)
):
    """
    Send a message to the chat
    """
    try:
        response = await chat_service.chat(
            request=request,
            user_id=user_id
        )
        
        # Background task for analytics
        background_tasks.add_task(
            log_conversation,
            conversation_id=response.conversation_id,
            user_id=user_id
        )
        
        return response
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/conversations/{conversation_id}", response_model=ConversationModel)
async def get_conversation(
    conversation_id: UUID,
    chat_service: ChatService = Depends(get_chat_service),
    user_id: UUID = Depends(get_current_user)
):
    """
    Get conversation history
    """
    messages = await chat_service.get_conversation_history(
        conversation_id=conversation_id
    )
    return ConversationModel(
        id=conversation_id,
        messages=messages
    )


@router.delete("/conversations/{conversation_id}")
async def delete_conversation(
    conversation_id: UUID,
    chat_service: ChatService = Depends(get_chat_service),
    user_id: UUID = Depends(get_current_user)
):
    """
    Delete a conversation
    """
    await chat_service.delete_conversation(conversation_id)
    return {"message": "Conversation deleted"}


async def log_conversation(conversation_id: UUID, user_id: Optional[UUID]):
    """Background task for logging"""
    # Implement logging logic
    pass
```

#### Main Application (app/main.py)
```python
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
import uvicorn

from app.api.routes import chat_routes, document_routes, health_routes
from app.api.middleware.logging import LoggingMiddleware
from app.api.middleware.rate_limit import RateLimitMiddleware
from app.config.settings import settings
from app.utils.errors import AppException, app_exception_handler
from app.infrastructure.database.postgres import init_db
from app.infrastructure.cache.redis import redis_client


def create_app() -> FastAPI:
    """Application factory"""
    app = FastAPI(
        title="LangChain API",
        version=settings.VERSION,
        docs_url="/api/docs" if settings.ENVIRONMENT == "development" else None,
        redoc_url="/api/redoc" if settings.ENVIRONMENT == "development" else None
    )
    
    # Setup middleware
    app.add_middleware(
        CORSMiddleware,
        allow_origins=settings.CORS_ORIGINS,
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )
    app.add_middleware(LoggingMiddleware)
    app.add_middleware(RateLimitMiddleware)
    
    # Setup exception handlers
    app.add_exception_handler(AppException, app_exception_handler)
    
    # Include routers
    app.include_router(health_routes.router, prefix="/api")
    app.include_router(chat_routes.router, prefix="/api/v1")
    app.include_router(document_routes.router, prefix="/api/v1")
    
    @app.on_event("startup")
    async def startup():
        """Startup tasks"""
        await init_db()
        await redis_client.initialize()
    
    @app.on_event("shutdown")
    async def shutdown():
        """Shutdown tasks"""
        await redis_client.close()
    
    return app


app = create_app()

if __name__ == "__main__":
    uvicorn.run(
        "app.main:app",
        host=settings.HOST,
        port=settings.PORT,
        reload=settings.ENVIRONMENT == "development"
    )
```

#### Settings (app/config/settings.py)
```python
from pydantic_settings import BaseSettings
from typing import List, Optional
from functools import lru_cache


class Settings(BaseSettings):
    # App
    APP_NAME: str = "LangChain API"
    VERSION: str = "1.0.0"
    ENVIRONMENT: str = "development"
    DEBUG: bool = False
    HOST: str = "0.0.0.0"
    PORT: int = 8000
    
    # Database
    DATABASE_URL: str = "postgresql://user:pass@localhost:5432/langchain"
    DATABASE_POOL_SIZE: int = 20
    
    # Redis
    REDIS_URL: str = "redis://localhost:6379/0"
    REDIS_MAX_CONNECTIONS: int = 10
    
    # OpenAI
    OPENAI_API_KEY: str
    OPENAI_MODEL: str = "gpt-3.5-turbo"
    OPENAI_TEMPERATURE: float = 0.7
    OPENAI_MAX_TOKENS: int = 1000
    
    # Vector Store
    VECTOR_STORE_TYPE: str = "chroma"
    CHROMA_PERSIST_DIR: str = "./chroma_db"
    PINECONE_API_KEY: Optional[str] = None
    PINECONE_ENVIRONMENT: Optional[str] = None
    
    # Security
    SECRET_KEY: str
    JWT_ALGORITHM: str = "HS256"
    ACCESS_TOKEN_EXPIRE_MINUTES: int = 30
    CORS_ORIGINS: List[str] = ["http://localhost:3000"]
    
    # Rate Limiting
    RATE_LIMIT_REQUESTS: int = 100
    RATE_LIMIT_PERIOD: int = 60  # seconds
    
    # Logging
    LOG_LEVEL: str = "INFO"
    LOG_FORMAT: str = "json"
    
    class Config:
        env_file = ".env"
        case_sensitive = True


@lru_cache()
def get_settings() -> Settings:
    return Settings()


settings = get_settings()
```

#### Docker Compose (docker/docker-compose.yml)
```yaml
version: '3.8'

services:
  api:
    build:
      context: ..
      dockerfile: docker/Dockerfile
    ports:
      - "8000:8000"
    environment:
      - DATABASE_URL=postgresql://postgres:postgres@db:5432/langchain
      - REDIS_URL=redis://redis:6379/0
      - OPENAI_API_KEY=${OPENAI_API_KEY}
    volumes:
      - ../app:/app/app
      - chroma_data:/app/chroma_db
    depends_on:
      - db
      - redis
    networks:
      - langchain-network
    command: uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload

  worker:
    build:
      context: ..
      dockerfile: docker/Dockerfile
    environment:
      - DATABASE_URL=postgresql://postgres:postgres@db:5432/langchain
      - REDIS_URL=redis://redis:6379/0
      - OPENAI_API_KEY=${OPENAI_API_KEY}
    depends_on:
      - db
      - redis
    networks:
      - langchain-network
    command: python -m app.worker

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=langchain
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - langchain-network

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - langchain-network

  chroma:
    image: chromadb/chroma:latest
    ports:
      - "8001:8000"
    volumes:
      - chroma_data:/chroma/chroma
    networks:
      - langchain-network

volumes:
  postgres_data:
  redis_data:
  chroma_data:

networks:
  langchain-network:
    driver: bridge
```

---

## 5. เวิร์กโฟลว์การพัฒนา (Development Workflow)

### 5.1 Golang Development Workflow

```mermaid
graph TD
    A[Write Code] --> B[Run Tests]
    B --> C{Tests Pass?}
    C -->|No| A
    C -->|Yes| D[Run Linter]
    D --> E{Lint Pass?}
    E -->|No| A
    E -->|Yes| F[Build Binary]
    F --> G[Run Migrations]
    G --> H[Start Application]
    H --> I[Test API]
    I --> J{All Good?}
    J -->|No| A
    J -->|Yes| K[Create PR]
```

**ขั้นตอนการพัฒนา:**
1. **Local Development**: ใช้ Air สำหรับ hot-reload
   ```bash
   air
   ```

2. **Testing**:
   ```bash
   make test              # Run all tests
   make test-coverage     # Check coverage
   ```

3. **Database Migration**:
   ```bash
   make migrate-create name=add_users_table
   make migrate-up
   ```

4. **Generate Swagger Docs**:
   ```bash
   make swagger
   ```

5. **Build and Run**:
   ```bash
   make build
   ./bin/app
   ```

### 5.2 LangChain FastAPI Development Workflow

```mermaid
graph TD
    A[Define Requirements] --> B[Create/Update Models]
    B --> C[Write Business Logic]
    C --> D[Create API Endpoints]
    D --> E[Write Tests]
    E --> F{Tests Pass?}
    F -->|No| C
    F -->|Yes| G[Run Locally]
    G --> H{Test API}
    H -->|Issues| C
    H -->|OK| I[Create PR]
```

**ขั้นตอนการพัฒนา:**
1. **Setup Environment**:
   ```bash
   python -m venv venv
   source venv/bin/activate
   pip install -r requirements-dev.txt
   ```

2. **Run Locally**:
   ```bash
   uvicorn app.main:app --reload
   ```

3. **Database Migrations**:
   ```bash
   alembic revision --autogenerate -m "message"
   alembic upgrade head
   ```

4. **Testing**:
   ```bash
   pytest
   pytest --cov=app tests/
   ```

5. **Docker Development**:
   ```bash
   docker-compose up -d
   docker-compose logs -f api
   ```

---

## 6. ผังการไหลของข้อมูล (Data Flow Architecture)

### 6.1 Golang Clean Architecture Data Flow

```mermaid
graph LR
    Client[Client] -->|HTTP Request| Router[Router]
    Router -->|Route| Middleware[Middleware Stack]
    Middleware -->|Auth/Rate Limit| Handler[Handler]
    Handler -->|DTO| Service[Service Layer]
    Service -->|Entity| Repository[Repository Interface]
    Repository -->|Implementation| DB[(Database)]
    
    Service -->|Cache Key| Cache[Redis Cache]
    Service -->|Job| Queue[Message Queue]
    
    Cache -->|Cached Data| Service
    Queue -->|Async Job| Worker[Worker Pool]
    Worker -->|Process| Service
    
    Service -->|Business Logic| Domain[Domain Layer]
    Domain -->|Events| EventBus[Event Bus]
    EventBus -->|Handle| EventHandler[Event Handler]
```

**ลำดับการไหลของข้อมูล:**
1. Client ส่ง HTTP Request มา
2. Router จับคู่กับ Handler ที่เหมาะสม
3. Middleware ทำการตรวจสอบ (Auth, Rate Limit, Logging)
4. Handler รับ Request และแปลงเป็น DTO
5. Service Layer ประมวลผล Business Logic
6. Repository Interface ถูกเรียกเพื่อ Data
7. Cache ถูกตรวจสอบก่อน Query Database
8. Database ส่งข้อมูลกลับ
9. Service สร้าง Events (ถ้ามี)
10. Response ถูกส่งกลับไปยัง Client

### 6.2 LangChain FastAPI Data Flow

```mermaid
graph TD
    User[User] -->|Query| API[FastAPI Endpoint]
    API -->|Validate| DTO[Request DTO]
    DTO -->|Call| Service[Service Layer]
    
    Service -->|Check| Cache[(Redis Cache)]
    Cache -->|Cache Hit| Service
    
    Service -->|Cache Miss| Chain[LangChain Chain]
    
    Chain -->|Retrieve| VectorStore[(Vector Store)]
    Chain -->|LLM Call| LLM[OpenAI/LLM]
    Chain -->|Tool Use| Tools[Tools/Search]
    
    VectorStore -->|Similar Docs| Chain
    LLM -->|Response| Chain
    Tools -->|Results| Chain
    
    Chain -->|Final Response| Service
    Service -->|Store| DB[(PostgreSQL)]
    Service -->|Cache| Cache
    
    Service -->|Response DTO| API
    API -->|JSON| User
    
    Service -->|Async Task| Queue[Message Queue]
    Queue -->|Process| Worker[Background Worker]
    Worker -->|Update| DB
```

**อธิบายการไหล:**
1. **Request Phase**: User ส่ง query → API validation → Service layer
2. **Cache Phase**: Service ตรวจ cache → ถ้าเจอ return ทันที
3. **Processing Phase**: ถ้าไม่เจอ cache → LangChain chain ทำงาน
4. **Retrieval Phase**: Chain query vector store → get relevant documents
5. **Generation Phase**: Chain เรียก LLM with context → generate response
6. **Post-processing**: Response → cache → database → return to user

### 6.3 การไหลของข้อมูลแบบ Streaming

```mermaid
sequenceDiagram
    participant User
    participant API
    participant Service
    participant Chain
    participant LLM
    
    User->>API: POST /chat/stream
    API->>Service: stream_chat(request)
    Service->>Chain: astream(messages)
    
    loop For each token
        Chain->>LLM: get_next_token()
        LLM-->>Chain: token
        Chain-->>Service: token
        Service-->>API: token
        API-->>User: SSE: token
    end
    
    Chain-->>Service: complete
    Service-->>API: done
    API-->>User: SSE: [DONE]
```

---

## 7. แม่แบบโค้ด (Code Templates)

### 7.1 Golang Template Files

#### Repository Template (internal/core/[module]/repository/repository.go)
```go
package repository

import (
    "context"
    "google/uuid"
    "your-project/internal/core/[module]/entity"
)

type [Module]Repository interface {
    Create(ctx context.Context, entity *entity.[Entity]) error
    Update(ctx context.Context, entity *entity.[Entity]) error
    Delete(ctx context.Context, id uuid.UUID) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.[Entity], error)
    FindAll(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]*entity.[Entity], int64, error)
}
```

#### Service Template (internal/core/[module]/service/service.go)
```go
package service

import (
    "context"
    "google/uuid"
    "your-project/internal/core/[module]/entity"
    "your-project/internal/core/[module]/repository"
    "your-project/pkg/logger"
)

type [Module]Service interface {
    Create(ctx context.Context, req *CreateRequest) (*entity.[Entity], error)
    GetByID(ctx context.Context, id uuid.UUID) (*entity.[Entity], error)
    Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*entity.[Entity], error)
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, filter map[string]interface{}, page, pageSize int) (*ListResponse, error)
}

type [module]Service struct {
    repo   repository.[Module]Repository
    cache  cache.Cache
    logger logger.Logger
}

func New[Module]Service(
    repo repository.[Module]Repository,
    cache cache.Cache,
    logger logger.Logger,
) [Module]Service {
    return &[module]Service{
        repo:   repo,
        cache:  cache,
        logger: logger,
    }
}

func (s *[module]Service) Create(ctx context.Context, req *CreateRequest) (*entity.[Entity], error) {
    // Validation
    if err := req.Validate(); err != nil {
        return nil, errors.NewBadRequestError(err.Error())
    }
    
    // Business logic
    entity := &entity.[Entity]{
        ID:        uuid.New(),
        // Map fields from req
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    // Save
    if err := s.repo.Create(ctx, entity); err != nil {
        s.logger.Error("failed to create [module]", "error", err)
        return nil, errors.Wrap(err, "failed to create [module]")
    }
    
    // Invalidate cache
    cacheKey := fmt.Sprintf("[module]:%s", entity.ID)
    _ = s.cache.Delete(ctx, cacheKey)
    
    return entity, nil
}
```

#### Handler Template (internal/core/[module]/handler/handler.go)
```go
package handler

import (
    "net/http"
    
    "gin-gonic/gin"
    "google/uuid"
    
    "your-project/internal/core/[module]/dto"
    "your-project/internal/core/[module]/service"
    "your-project/internal/transport/httpx"
    "your-project/pkg/logger"
)

type [Module]Handler struct {
    service service.[Module]Service
    logger  logger.Logger
}

func New[Module]Handler(service service.[Module]Service, logger logger.Logger) *[Module]Handler {
    return &[Module]Handler{
        service: service,
        logger:  logger,
    }
}

// Create godoc
// @Summary Create new [module]
// @Tags [module]
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateRequest true "Create request"
// @Success 201 {object} dto.[Module]Response
// @Failure 400 {object} httpx.ErrorResponse
// @Router /api/v1/[module] [post]
func (h *[Module]Handler) Create(c *gin.Context) {
    var req dto.CreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        httpx.RespondError(c, httpx.ErrBadRequest(err.Error()))
        return
    }
    
    entity, err := h.service.Create(c.Request.Context(), &req)
    if err != nil {
        httpx.RespondError(c, err)
        return
    }
    
    httpx.RespondSuccess(c, http.StatusCreated, dto.ToResponse(entity))
}

// GetByID godoc
// @Summary Get [module] by ID
// @Tags [module]
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Module ID"
// @Success 200 {object} dto.[Module]Response
// @Failure 404 {object} httpx.ErrorResponse
// @Router /api/v1/[module]/{id} [get]
func (h *[Module]Handler) GetByID(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        httpx.RespondError(c, httpx.ErrBadRequest("invalid id format"))
        return
    }
    
    entity, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        httpx.RespondError(c, err)
        return
    }
    
    httpx.RespondSuccess(c, http.StatusOK, dto.ToResponse(entity))
}
```

#### DTO Template (internal/core/[module]/dto/dto.go)
```go
package dto

import (
    "time"
    
    "go-playground/validator/v10"
    "google/uuid"
    "your-project/internal/core/[module]/entity"
)

type CreateRequest struct {
    // Add fields
    Name  string `json:"name" binding:"required"`
    Field string `json:"field" binding:"required"`
}

func (r *CreateRequest) Validate() error {
    validate := validator.New()
    return validate.Struct(r)
}

type UpdateRequest struct {
    Name  string `json:"name" binding:"omitempty"`
    Field string `json:"field" binding:"omitempty"`
}

func (r *UpdateRequest) Validate() error {
    validate := validator.New()
    return validate.Struct(r)
}

type [Module]Response struct {
    ID        uuid.UUID `json:"id"`
    Name      string    `json:"name"`
    Field     string    `json:"field"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func ToResponse(entity *entity.[Entity]) *[Module]Response {
    return &[Module]Response{
        ID:        entity.ID,
        Name:      entity.Name,
        Field:     entity.Field,
        CreatedAt: entity.CreatedAt,
        UpdatedAt: entity.UpdatedAt,
    }
}

type ListResponse struct {
    Items      []*[Module]Response `json:"items"`
    Total      int64                `json:"total"`
    Page       int                  `json:"page"`
    PageSize   int                  `json:"page_size"`
    TotalPages int                  `json:"total_pages"`
}
```

### 7.2 LangChain FastAPI Template Files

#### Service Template (app/application/services/base_service.py)
```python
from typing import TypeVar, Generic, List, Optional, Dict, Any
from uuid import UUID
from abc import ABC, abstractmethod

from app.core.entities.base import BaseEntity
from app.core.interfaces.repository import BaseRepository
from app.infrastructure.cache.redis_cache import RedisCache

T = TypeVar('T', bound=BaseEntity)


class BaseService(Generic[T], ABC):
    def __init__(
        self,
        repository: BaseRepository[T],
        cache: Optional[RedisCache] = None
    ):
        self.repository = repository
        self.cache = cache
    
    async def get_by_id(self, id: UUID) -> Optional[T]:
        # Check cache first
        if self.cache:
            cached = await self.cache.get(f"{self.entity_name}:{id}")
            if cached:
                return cached
        
        # Get from database
        entity = await self.repository.get_by_id(id)
        
        # Cache if found
        if entity and self.cache:
            await self.cache.set(
                f"{self.entity_name}:{id}",
                entity,
                ttl=300
            )
        
        return entity
    
    async def create(self, data: Dict[str, Any]) -> T:
        entity = self.entity_class.create(**data)
        created = await self.repository.create(entity)
        
        # Invalidate list cache
        if self.cache:
            await self.cache.delete_pattern(f"{self.entity_name}:list:*")
        
        return created
    
    async def update(self, id: UUID, data: Dict[str, Any]) -> Optional[T]:
        entity = await self.get_by_id(id)
        if not entity:
            return None
        
        for key, value in data.items():
            if hasattr(entity, key):
                setattr(entity, key, value)
        
        updated = await self.repository.update(entity)
        
        # Invalidate caches
        if self.cache:
            await self.cache.delete(f"{self.entity_name}:{id}")
            await self.cache.delete_pattern(f"{self.entity_name}:list:*")
        
        return updated
    
    async def delete(self, id: UUID) -> bool:
        result = await self.repository.delete(id)
        
        # Invalidate caches
        if self.cache and result:
            await self.cache.delete(f"{self.entity_name}:{id}")
            await self.cache.delete_pattern(f"{self.entity_name}:list:*")
        
        return result
    
    async def list(
        self,
        skip: int = 0,
        limit: int = 100,
        **filters
    ) -> List[T]:
        cache_key = f"{self.entity_name}:list:{skip}:{limit}:{hash(frozenset(filters.items()))}"
        
        # Check cache
        if self.cache:
            cached = await self.cache.get(cache_key)
            if cached:
                return cached
        
        # Get from database
        items = await self.repository.list(skip=skip, limit=limit, **filters)
        
        # Cache results
        if self.cache:
            await self.cache.set(cache_key, items, ttl=60)
        
        return items
```

#### Chain Template (app/infrastructure/llm/chains/qa_chain.py)
```python
from typing import List, Optional, Dict, Any
from langchain.chains import RetrievalQA
from langchain.prompts import PromptTemplate
from langchain.callbacks import AsyncCallbackHandler
from langchain.schema import BaseRetriever

from app.infrastructure.llm.base import BaseLLMService
from app.infrastructure.vector_store.base import BaseVectorStore


class QAChain:
    """Template for QA chain with retrieval"""
    
    def __init__(
        self,
        llm_service: BaseLLMService,
        vector_store: BaseVectorStore,
        retriever_kwargs: Optional[Dict[str, Any]] = None
    ):
        self.llm_service = llm_service
        self.vector_store = vector_store
        self.retriever_kwargs = retriever_kwargs or {"k": 4}
        
        # Create retriever
        self.retriever = self.vector_store.as_retriever(
            **self.retriever_kwargs
        )
        
        # Create prompt template
        self.prompt = PromptTemplate(
            template="""
            Use the following pieces of context to answer the question at the end.
            If you don't know the answer, just say that you don't know, don't try to make up an answer.
            
            Context: {context}
            
            Question: {question}
            
            Answer: """,
            input_variables=["context", "question"]
        )
        
        # Create chain
        self.chain = RetrievalQA.from_chain_type(
            llm=self.llm_service.get_llm(),
            chain_type="stuff",
            retriever=self.retriever,
            chain_type_kwargs={
                "prompt": self.prompt,
                "verbose": True
            },
            return_source_documents=True
        )
    
    async def ainvoke(
        self,
        query: str,
        callbacks: Optional[List[AsyncCallbackHandler]] = None
    ) -> Dict[str, Any]:
        """Async invoke chain"""
        result = await self.chain.ainvoke(
            {"query": query},
            callbacks=callbacks
        )
        
        return {
            "answer": result["result"],
            "sources": [
                {
                    "content": doc.page_content,
                    "metadata": doc.metadata
                }
                for doc in result["source_documents"]
            ]
        }
```

#### API Route Template (app/api/routes/module_routes.py)
```python
from fastapi import APIRouter, Depends, HTTPException, Query
from typing import List, Optional
from uuid import UUID

from app.api.dependencies.container import get_service
from app.api.models.request_models import CreateRequest, UpdateRequest
from app.api.models.response_models import EntityResponse, ListResponse
from app.application.services.base_service import BaseService
from app.api.middleware.auth import get_current_user, require_permissions

router = APIRouter(prefix="/module", tags=["module"])


@router.post("/", response_model=EntityResponse, status_code=201)
async def create_entity(
    request: CreateRequest,
    service: BaseService = Depends(get_service),
    user_id: UUID = Depends(get_current_user)
):
    """
    Create a new entity
    """
    try:
        entity = await service.create(request.dict())
        return EntityResponse.from_entity(entity)
    except Exception as e:
        raise HTTPException(status_code=400, detail=str(e))


@router.get("/{entity_id}", response_model=EntityResponse)
async def get_entity(
    entity_id: UUID,
    service: BaseService = Depends(get_service),
    user_id: UUID = Depends(get_current_user)
):
    """
    Get entity by ID
    """
    entity = await service.get_by_id(entity_id)
    if not entity:
        raise HTTPException(status_code=404, detail="Entity not found")
    
    return EntityResponse.from_entity(entity)


@router.get("/", response_model=ListResponse)
async def list_entities(
    skip: int = Query(0, ge=0),
    limit: int = Query(100, ge=1, le=1000),
    search: Optional[str] = None,
    service: BaseService = Depends(get_service),
    user_id: UUID = Depends(get_current_user)
):
    """
    List entities with pagination
    """
    filters = {}
    if search:
        filters["search"] = search
    
    items, total = await service.list(
        skip=skip,
        limit=limit,
        **filters
    )
    
    return ListResponse(
        items=[EntityResponse.from_entity(item) for item in items],
        total=total,
        skip=skip,
        limit=limit
    )


@router.put("/{entity_id}", response_model=EntityResponse)
async def update_entity(
    entity_id: UUID,
    request: UpdateRequest,
    service: BaseService = Depends(get_service),
    user_id: UUID = Depends(get_current_user),
    _: bool = Depends(require_permissions(["admin"]))
):
    """
    Update entity (admin only)
    """
    entity = await service.update(
        entity_id,
        request.dict(exclude_unset=True)
    )
    if not entity:
        raise HTTPException(status_code=404, detail="Entity not found")
    
    return EntityResponse.from_entity(entity)


@router.delete("/{entity_id}", status_code=204)
async def delete_entity(
    entity_id: UUID,
    service: BaseService = Depends(get_service),
    user_id: UUID = Depends(get_current_user),
    _: bool = Depends(require_permissions(["admin"]))
):
    """
    Delete entity (admin only)
    """
    deleted = await service.delete(entity_id)
    if not deleted:
        raise HTTPException(status_code=404, detail="Entity not found")
```

## บทสรุป

การออกแบบสถาปัตยกรรมซอฟต์แวร์ด้วย Golang และ LangChain FastAPI ตามหลักการ Domain-Driven Design (DDD) และ Clean Architecture ช่วยให้ระบบมีโครงสร้างที่ชัดเจน แยกส่วนการทำงานออกจากกัน ทำให้:

1. **Maintainability**: ง่ายต่อการบำรุงรักษาและพัฒนาเพิ่มเติม
2. **Testability**: สามารถทดสอบแต่ละส่วนได้อย่างอิสระ
3. **Scalability**: ขยายระบบได้ง่ายตามความต้องการ
4. **Flexibility**: เปลี่ยนเทคโนโลยีได้โดยไม่กระทบ Business Logic
5. **Team Development**: ทีมสามารถทำงานแบบ Parallel ได้

### การเลือกใช้ระหว่าง Golang และ LangChain FastAPI

| พิจารณา | Golang | LangChain FastAPI |
|--------|--------|-------------------|
| **Performance** | สูงมาก | ปานกลาง |
| **Concurrency** | ยอดเยี่ยม (Goroutines) | ดี (Async/Await) |
| **AI/ML Features** | จำกัด | ยอดเยี่ยม |
| **Ecosystem** | ดีสำหรับ Backend | ดีสำหรับ Data/ML |
| **Learning Curve** | ปานกลาง | ต่ำ |

**แนะนำให้ใช้ Golang เมื่อ:**
- ต้องการ performance สูง
- ระบบมี traffic มาก
- ต้องการการทำงานแบบ concurrent
- ต้องการ type safety

**แนะนำให้ใช้ LangChain FastAPI เมื่อ:**
- ต้องการ integrate กับ LLMs
- ต้องการพัฒนา AI applications
- ต้องการ prototyping ที่รวดเร็ว
- ทีมมีประสบการณ์ Python

### แหล่งข้อมูลเพิ่มเติม

1. **Golang Resources**:
   - [Official Go Documentation](https://golang.org/doc/)
   - [Clean Architecture in Go](https://bxcodec/go-clean-arch)
   - [GORM Documentation](https://gorm.io/docs/)

2. **LangChain Resources**:
   - [LangChain Documentation](https://python.langchain.com/docs)
   - [FastAPI Documentation](https://fastapi.tiangolo.com/)
   - [OpenAI API Documentation](https://platform.openai.com/docs)

3. **DDD Resources**:
   - [Domain-Driven Design in Go](https://marcusolsson/goddd)
   - [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

**จัดทำโดย**: ทีมพัฒนาซอฟต์แวร์  
**ปรับปรุงล่าสุด**: 2024  
**เวอร์ชัน**: 1.0.0



# คู่มือการใช้งาน Repository (Repo)

## โครงสร้างโปรเจกต์และการจัดระเบียบโมดูล
จุดเริ่มต้นของ service อยู่ที่ `cmd/app/main.go`  
การประกอบแอปพลิเคชันอยู่ที่ `internal/apps/app/bootstrap/`  
การประกอบเส้นทาง (routing) อยู่ที่ `internal/apps/app/router/`  
โมดูลหลัก (domain) อยู่ที่ `internal/core/` (เช่น `auth`, `user`, `health`)  
แต่ละโมดูลจะมีโครงสร้างย่อย เช่น `handler`, `service`, `repository`, `dto`, `model`, `routes`  
โครงสร้างพื้นฐาน (infrastructure) อยู่ที่ `internal/platform/` (เช่น `config`, `db`)  
เครื่องมือสำหรับ transport layer อยู่ที่ `internal/transport/` (เช่น `httpx`, `middleware`)  
ส่วนประกอบที่สามารถนำมาใช้ซ้ำได้ (reusable components) อยู่ที่ `pkg/` (cache, jwt, logger, queue, transaction, utils)  
เอกสาร API ที่ถูกสร้างจะอยู่ที่ `api/app/`  
การตั้งค่าอยู่ที่ `configs/`  
SQL migrations อยู่ที่ `migrations/`  
ทรัพยากรสำหรับ deploy อยู่ที่ `deploy/`  
สคริปต์ต่างๆ อยู่ที่ `scripts/`

## คำสั่ง build, test และพัฒนา
- `./scripts/dev.sh` ใช้สร้าง Swagger และใช้ `air` สำหรับ hot reload (ดู `.air.toml`)
- `go run cmd/app/main.go` สำหรับรัน service โดยตรง
- `./scripts/swagger.sh` ใช้สร้าง Swagger ใหม่ไปยัง `api/app/`
- `./scripts/build.sh` ใช้ build binary สำหรับหลาย ๆ แพลตฟอร์มไปยัง `build/`
- `go test ./...` ใช้รัน tests ทั้งหมด; สามารถเพิ่ม `-race` หรือ `-coverprofile=coverage.out`
- `docker compose -f deploy/docker/docker-compose.yaml up -d` ใช้เริ่ม PostgreSQL และ Redis

## รูปแบบการเขียนโค้ดและการตั้งชื่อ
ใช้ `gofmt` (ย่อหน้า (indent) ด้วย tab)  
ชื่อ package เป็นตัวพิมพ์เล็ก, ชื่อไฟล์ใช้ snake_case (เช่น `user_service.go`, `auth_handler.go`)  
identifier ที่ export ใช้ PascalCase, ที่ไม่ export ใช้ camelCase  
เมื่อเพิ่ม type ใหม่ ควรใส่ไว้ในโมดูลหลักที่เกี่ยวข้องก่อน และรักษาการแบ่งชั้น (layer) handler/service/repository ให้ชัดเจน

## ข้อกำหนดการทดสอบ
tests จะอยู่ไดเรกทอรีเดียวกับโค้ด, ชื่อไฟล์ `*_test.go`, ชื่อฟังก์ชัน `TestXxx`  
โปรเจกต์ใช้ `testify` (ดูตัวอย่างได้ที่ `internal/core/user/service/user_service_test.go`)  
เมื่อแก้ไข logic ของ service สามารถรัน: `go test ./internal/core/user/service/`

## ข้อกำหนดการ commit และ PR
ข้อความ commit ให้ใช้รูปแบบ prefix ตามที่เคยใช้ (เช่น `feat: ...`, `fix: ...`)  
PR ควรมีคำอธิบายสั้น ๆ, คำสั่งสำหรับทดสอบ, และระบุว่ามีการอัปเดตการตั้งค่า/migration/Swagger หรือไม่ (หากมีการเปลี่ยนแปลงที่ `api/app/` กรุณาระบุด้วย)

## การตั้งค่าและความปลอดภัย
สำหรับการพัฒนาในเครื่อง ให้คัดลอก `configs/config.example.yaml` เป็น `configs/config.yaml` และสามารถใช้ environment variable ที่ขึ้นต้นด้วย `APP_` เพื่อแทนที่ค่าได้  
อย่า commit คีย์สำคัญ (secret); สำหรับ production ให้ใช้ `configs/config.production.yaml` หรือ environment variable ในการแทนที่  
หากมีการแก้ไข route หรือโครงสร้าง request/response กรุณารัน `./scripts/swagger.sh` เพื่ออัปเดตเอกสาร


# การตั้ง ชื่อ Branch  ทีม ใหญ่   branch หลักดังนี้:​

- main: เก็บโค้ดที่พร้อมใช้งานใน production เท่านั้น
- develop: เป็น branch หลักสำหรับการพัฒนา feature ต่างๆ
- feature/xxx: branch สำหรับพัฒนา feature ใหม่แต่ละตัว
- release/xxx: branch สำหรับเตรียมการ release เวอร์ชันใหม่
- hotfix/xxx: branch สำหรับแก้ไข bug เร่งด่วนใน production​

-  สำหรับ **ทีมขนาดใหญ่ (Large Team)** การตั้งชื่อ Branch ต้องเน้นความ **"ตรวจสอบได้ (Traceability)"** และ **"ความเป็นมาตรฐาน (Standardization)"** เพื่อให้รู้ว่าใครทำอะไร เกี่ยวข้องกับงานไหน และสถานะเป็นอย่างไร

นี่คือมาตรฐานการตั้งชื่อ Branch ที่นิยมใช้ในองค์กรใหญ่และทำงานร่วมกับระบบ Ticket (Jira, Trello, ClickUp):

### รูปแบบหลัก (Pattern)

ใช้เครื่องหมาย `/` แบ่งประเภท และ `-` คั่นคำในชื่อ (Kebab-case)
> `ประเภท/รหัสงาน-คำอธิบายสั้นๆ`

***

### 1. Feature Branches (`feature/xxx`)

ใช้สำหรับพัฒนาฟีเจอร์ใหม่ แตกจาก `develop`

* **Pattern:** `feature/<TICKET-ID>-<short-description>`
* **ความสำคัญ:** ทีมใหญ่ต้องระบุ `Ticket ID` เสมอ เพื่อให้ระบบ Automation (เช่น Jira/GitLab) ลิงก์โค้ดเข้ากับ Task งานอัตโนมัติ
* **ตัวอย่าง:**
    * `feature/JIRA-123-login-screen` (มี Ticket ID ชัดเจน)
    * `feature/AUTH-456-google-oauth`
    * `feature/PAY-789-payment-gateway`


### 2. Release Branches (`release/xxx`)

ใช้สำหรับเตรียมเวอร์ชันใหม่ แตกจาก `develop` เพื่อทำ Final Test

* **Pattern:** `release/v<MAJOR.MINOR.PATCH>` (ตามหลัก Semantic Versioning)
* **ความสำคัญ:** ห้ามใช้ชื่อเล่น (เช่น release/summer-update) ต้องใช้ตัวเลขเวอร์ชันเท่านั้น
* **ตัวอย่าง:**
    * `release/v1.0.0` (เวอร์ชันแรก)
    * `release/v1.2.0` (เพิ่มฟีเจอร์ใหม่)
    * `release/v2.0.0-rc1` (Release Candidate 1)


### 3. Hotfix Branches (`hotfix/xxx`)

ใช้แก้บั๊กเร่งด่วนบน Production แตกจาก `main`

```
*   **Pattern:** `hotfix/<TICKET-ID>-<short-description>` หรือ `hotfix/v<VERSION>-<description>`
```

* **ความสำคัญ:** ต้องระบุสิ่งที่แก้ชัดเจน เพราะเป็น Branch ที่ซีเรียสที่สุด
* **ตัวอย่าง:**
    * `hotfix/v1.0.1-fix-login-crash` (ระบุเวอร์ชันที่จะแก้)
    * `hotfix/PROD-99-fix-memory-leak` (อิงตาม Ticket แจ้งปัญหา)


### 4. (เสริม) Bugfix Branches (`bugfix/xxx`)

สำหรับทีมใหญ่ มักแยก "บั๊กทั่วไป" (บน develop) ออกจาก "บั๊กเร่งด่วน" (hotfix บน main)

* **Pattern:** `bugfix/<TICKET-ID>-<description>`
* **ความสำคัญ:** ใช้แก้บั๊กที่เจอระหว่าง Test ในช่วง Develop (ยังไม่ขึ้น Production)
* **ตัวอย่าง:**
    * `bugfix/QA-55-fix-button-color`
    * `bugfix/JIRA-124-typo-correction`


### 5. (เสริม) Chore/Refactor (`chore/xxx`)

สำหรับงานที่ไม่เกี่ยวกับ Business Logic โดยตรง เช่น อัปเกรด Library หรือจัดระเบียบโค้ด

* **ตัวอย่าง:**
    * `chore/update-nestjs-v10`
    * `refactor/clean-up-user-service`

***

### ข้อตกลงร่วมกัน (Best Practices)

1. **ตัวพิมพ์เล็กทั้งหมด (Lowercase):** ป้องกันปัญหา Case Sensitivity บน Windows/Mac (`Feature/Login` ❌ -> `feature/login` ✅)
2. **ใช้ขีดกลาง (Hyphen):** อ่านง่ายกว่า Underscore (`user_login` -> `user-login`)
3. **ห้ามใช้ชื่อบุคคล:** เช่น `feature/somchai-login` เพราะงานหนึ่งอาจทำหลายคน หรือมีการส่งต่องาน ให้ใช้ Ticket ID แทน
4. **ห้ามตั้งชื่อซ้ำ:** เมื่อ Merge และลบ Branch แล้ว ชื่อเดิมสามารถนำมาใช้ใหม่ได้ แต่ไม่แนะนำเพื่อป้องกันความสับสนใน History

### ตัวอย่างภาพรวมใน Project

| ประเภท | ชื่อ Branch | Ticket ที่เกี่ยวข้อง |
| :-- | :-- | :-- |
| **Main** | `main` | - |
| **Develop** | `develop` | - |
| **Feature** | `feature/CART-101-add-to-cart` | Ticket: CART-101 |
| **Bugfix** | `bugfix/CART-102-fix-cart-total` | Ticket: CART-102 |
| **Release** | `release/v1.5.0` | - |
| **Hotfix** | `hotfix/v1.5.1-emergency-fix` | Ticket: INC-001 |
 

# การตั้ง ชื่อ Branch  ทีม ใหญ่   branch หลักดังนี้:​

main: เก็บโค้ดที่พร้อมใช้งานใน production เท่านั้น
develop: เป็น branch หลักสำหรับการพัฒนา feature ต่างๆ
feature/xxx: branch สำหรับพัฒนา feature ใหม่แต่ละตัว
release/xxx: branch สำหรับเตรียมการ release เวอร์ชันใหม่
hotfix/xxx: branch สำหรับแก้ไข bug เร่งด่วนใน production​
enhancement/xxx: branch สำหรับแก้ไข enhancement ใน production​
 

สำหรับ **ทีมขนาดใหญ่ (Large Team)** การตั้งชื่อ Branch ต้องเน้นความ **"ตรวจสอบได้ (Traceability)"** และ **"ความเป็นมาตรฐาน (Standardization)"** เพื่อให้รู้ว่าใครทำอะไร เกี่ยวข้องกับงานไหน และสถานะเป็นอย่างไร

นี่คือมาตรฐานการตั้งชื่อ Branch ที่นิยมใช้ในองค์กรใหญ่และทำงานร่วมกับระบบ Ticket (Jira, Trello, ClickUp):

### รูปแบบหลัก (Pattern)

ใช้เครื่องหมาย `/` แบ่งประเภท และ `-` คั่นคำในชื่อ (Kebab-case)
> `ประเภท/รหัสงาน-คำอธิบายสั้นๆ`

***

### 1. Feature Branches (`feature/xxx`)

ใช้สำหรับพัฒนาฟีเจอร์ใหม่ แตกจาก `develop`

* **Pattern:** `feature/<TICKET-ID>-<short-description>`
* **ความสำคัญ:** ทีมใหญ่ต้องระบุ `Ticket ID` เสมอ เพื่อให้ระบบ Automation (เช่น Jira/GitLab) ลิงก์โค้ดเข้ากับ Task งานอัตโนมัติ
* **ตัวอย่าง:**
    * `feature/JIRA-123-login-screen` (มี Ticket ID ชัดเจน)
    * `feature/AUTH-456-google-oauth`
    * `feature/PAY-789-payment-gateway`


### 2. Release Branches (`release/xxx`)

ใช้สำหรับเตรียมเวอร์ชันใหม่ แตกจาก `develop` เพื่อทำ Final Test

* **Pattern:** `release/v<MAJOR.MINOR.PATCH>` (ตามหลัก Semantic Versioning)
* **ความสำคัญ:** ห้ามใช้ชื่อเล่น (เช่น release/summer-update) ต้องใช้ตัวเลขเวอร์ชันเท่านั้น
* **ตัวอย่าง:**
    * `release/v1.0.0` (เวอร์ชันแรก)
    * `release/v1.2.0` (เพิ่มฟีเจอร์ใหม่)
    * `release/v2.0.0-rc1` (Release Candidate 1)


### 3. Hotfix Branches (`hotfix/xxx`)

ใช้แก้บั๊กเร่งด่วนบน Production แตกจาก `main`

```
*   **Pattern:** `hotfix/<TICKET-ID>-<short-description>` หรือ `hotfix/v<VERSION>-<description>`
```

* **ความสำคัญ:** ต้องระบุสิ่งที่แก้ชัดเจน เพราะเป็น Branch ที่ซีเรียสที่สุด
* **ตัวอย่าง:**
    * `hotfix/v1.0.1-fix-login-crash` (ระบุเวอร์ชันที่จะแก้)
    * `hotfix/PROD-99-fix-memory-leak` (อิงตาม Ticket แจ้งปัญหา)


### 4. (เสริม) Bugfix Branches (`bugfix/xxx`)

สำหรับทีมใหญ่ มักแยก "บั๊กทั่วไป" (บน develop) ออกจาก "บั๊กเร่งด่วน" (hotfix บน main)

* **Pattern:** `bugfix/<TICKET-ID>-<description>`
* **ความสำคัญ:** ใช้แก้บั๊กที่เจอระหว่าง Test ในช่วง Develop (ยังไม่ขึ้น Production)
* **ตัวอย่าง:**
    * `bugfix/QA-55-fix-button-color`
    * `bugfix/JIRA-124-typo-correction`


### 5. (เสริม) Chore/Refactor (`chore/xxx`)

สำหรับงานที่ไม่เกี่ยวกับ Business Logic โดยตรง เช่น อัปเกรด Library หรือจัดระเบียบโค้ด

* **ตัวอย่าง:**
    * `chore/update-nestjs-v10`
    * `refactor/clean-up-user-service`

***

### ข้อตกลงร่วมกัน (Best Practices)

1. **ตัวพิมพ์เล็กทั้งหมด (Lowercase):** ป้องกันปัญหา Case Sensitivity บน Windows/Mac (`Feature/Login` ❌ -> `feature/login` ✅)
2. **ใช้ขีดกลาง (Hyphen):** อ่านง่ายกว่า Underscore (`user_login` -> `user-login`)
3. **ห้ามใช้ชื่อบุคคล:** เช่น `feature/somchai-login` เพราะงานหนึ่งอาจทำหลายคน หรือมีการส่งต่องาน ให้ใช้ Ticket ID แทน
4. **ห้ามตั้งชื่อซ้ำ:** เมื่อ Merge และลบ Branch แล้ว ชื่อเดิมสามารถนำมาใช้ใหม่ได้ แต่ไม่แนะนำเพื่อป้องกันความสับสนใน History

### ตัวอย่างภาพรวมใน Project

| ประเภท | ชื่อ Branch | Ticket ที่เกี่ยวข้อง |
| :-- | :-- | :-- |
| **Main** | `main` | - |
| **Develop** | `develop` | - |
| **Feature** | `feature/CART-101-add-to-cart` | Ticket: CART-101 |
| **Bugfix** | `bugfix/CART-102-fix-cart-total` | Ticket: CART-102 |
| **Release** | `release/v1.5.0` | - |
| **Hotfix** | `hotfix/v1.5.1-emergency-fix` | Ticket: INC-001 |



# Enchant Code **"Enhancement Code"** (การปรับปรุงโค้ดให้ดีขึ้น) หรืออาจเป็นการเขียนโค้ดให้ **"Clean" (Clean Code)** และมีคุณภาพสูง

ในบริบทของ Git Branching สำหรับทีมใหญ่ มักจะมีการแยกประเภท Branch สำหรับ "งานปรับปรุง" ออกจาก "ฟีเจอร์ใหม่" เพื่อให้ง่ายต่อการตรวจสอบ (Code Review) และการจัดลำดับความสำคัญ (Prioritization) ดังนี้ครับ:

- ในบริบทของ Git Branching สำหรับทีมใหญ่ มักจะมีการแยกประเภท Branch สำหรับ "งานปรับปรุง" ออกจาก "ฟีเจอร์ใหม่" เพื่อให้ง่ายต่อการตรวจสอบ (Code Review) และการจัดลำดับความสำคัญ (Prioritization) ดังนี้ครับ:


### **1. Enhancement Branch (`enhancement/xxx`)**

ใช้สำหรับงานที่ **"ไม่ใช่ฟีเจอร์ใหม่ แต่ทำให้ระบบดีขึ้น"** เช่น การปรับจูน Performance, ปรับปรุง UI/UX เล็กน้อย หรือเพิ่ม Logging

* **ความแตกต่างจาก Feature:** Feature คือสิ่งที่ User "เห็นและใช้งานได้ใหม่" ส่วน Enhancement มักเป็นการปรับปรุงของเดิมให้ดีกว่าเดิม
* **Base Branch:** แตกออกจาก `develop`
* **Pattern:** `enhancement/<TICKET-ID>-<description>`
* **ตัวอย่าง:**
    * `enhancement/CART-105-improve-loading-speed` (ทำให้โหลดเร็วขึ้น)
    * `enhancement/UI-202-adjust-button-shadow` (ปรับเงาปุ่มให้สวยขึ้น)

***

### **2. Refactor Branch (`refactor/xxx`)**

ใช้สำหรับงาน **"รื้อโครงสร้างโค้ด (Refactoring)"** โดยที่ **"ผลลัพธ์การทำงานต้องเหมือนเดิม"** (User ไม่เห็นความเปลี่ยนแปลง แต่โค้ดอ่านง่ายขึ้น บำรุงรักษาง่ายขึ้น)

* **ความสำคัญ:** ทีมใหญ่แยก Branch นี้ออกมาเพื่อบอก Reviewer ว่า *"ไม่ต้อง Test ฟังก์ชันนะ แค่ดู Logic ว่าเขียนดีขึ้นไหม"*
* **Base Branch:** แตกออกจาก `develop`
* **Pattern:** `refactor/<TICKET-ID>-<scope>`
* **ตัวอย่าง:**
    * `refactor/USER-300-clean-auth-service` (จัดระเบียบโค้ดใน Auth Service)
    * `refactor/CORE-404-remove-unused-imports` (ลบโค้ดที่ไม่ได้ใช้ออก)

***

### **3. Chore Branch (`chore/xxx`)**

ใช้สำหรับงาน **"งานบ้าน/งานจุกจิก"** ที่ไม่กระทบ Code หลัก เช่น อัปเกรด Library, แก้ไฟล์ Config, หรือเขียน Document

* **Base Branch:** แตกออกจาก `develop`
* **Pattern:** `chore/<TICKET-ID>-<task>`
* **ตัวอย่าง:**
    * `chore/DEVOPS-500-update-nestjs-v10` (อัปเกรด Version Framework)
    * `chore/DOC-101-update-readme` (แก้ไฟล์คู่มือ)

***

### **สรุปตารางเปรียบเทียบ "Enchant" (Improvement) Branches**

| ประเภท Branch | ความหมาย | ผลกระทบต่อ User | ต้องเขียน Test เพิ่มไหม? |
| :-- | :-- | :-- | :-- |
| `feature/xxx` | ของใหม่ | **เห็น** และใช้งานได้ | ✅ ต้องมี Unit/E2E Test |
| `enhancement/xxx` | ของเดิมที่ดีขึ้น | **เห็น** (ทำงานดีขึ้น/เร็วขึ้น) | ✅ อาจต้องแก้ Test เดิม |
| `refactor/xxx` | จัดระเบียบโค้ด | **ไม่เห็น** (ทำงานเหมือนเดิมเป๊ะ) | ❌ ไม่ควรแก้ Test (ถ้า Logic ไม่เปลี่ยน) |
| `chore/xxx` | งานจุกจิก/Config | **ไม่เห็น** | ❌ ไม่ต้องเขียน Test |

### **คำแนะนำเพิ่มเติมสำหรับทีมใหญ่**

ถ้าทีมของคุณต้องการเน้นเรื่อง "Enchant Code" (ทำให้โค้ดดูขลัง/เทพขึ้น 🧙‍♂️) แนะนำให้เพิ่มขั้นตอน **Automated Code Quality Check** ใน Pipeline ก่อน Merge:

1. **Linting:** บังคับใช้กฎการเขียนโค้ด (ESLint/Prettier)
2. **SonarQube:** สแกนหา "Code Smell" หรือจุดที่เขียนไม่ดี
3. **Commit Message Lint:** บังคับให้ใส่ชื่อ Branch ประเภทนี้ใน Commit Message เช่น `refactor: clean up user service` เพื่อให้ Log อ่านง่าย



# แนะนำรูปแบบการตั้งชื่อ branch สำหรับฟีเจอร์ขนาดเล็ก

สำหรับการตั้งชื่อ Branch ที่เป็น **"ฟีเจอร์ขนาดเล็ก" (Small Feature)** หรือการปรับแก้เล็กๆ น้อยๆ ที่อาจจะไม่ถึงขั้นเรียกว่าเป็น "Full Feature" ใหญ่ๆ นั้น หลักการสำคัญคือ **"อย่าสร้าง Prefix ใหม่เยอะเกินความจำเป็น"** เพื่อไม่ให้ทีมสับสนครับ

นี่คือ 3 แนวทางที่แนะนำ เรียงตามความนิยมและความเหมาะสม:

### 1. ใช้ `feature/` เหมือนเดิม (แนะนำสูงสุด ⭐️)

แม้จะเป็นงานเล็กๆ ก็ควรนับเป็น Feature เพื่อความสม่ำเสมอ (Consistency) ในการตั้งชื่อและการค้นหา

* **หลักการ:** ใช้ Pattern เดิมแต่เน้นคำอธิบายที่ **กระชับ** และเจาะจง
* **Pattern:** `feature/<TICKET>-<specific-action>`
* **ตัวอย่าง:**
    * `feature/CART-201-add-delete-btn` (แค่เพิ่มปุ่มลบปุ่มเดียว)
    * `feature/USER-305-change-font-size` (เปลี่ยนขนาดฟอนต์)
    * `feature/AUTH-102-hide-password` (ซ่อนรหัสผ่าน)


### 2. ใช้ `tweak/` (สำหรับงานจุกจิก/ปรับแต่ง)

ถ้าทีมรู้สึกว่าคำว่า `feature` ดูยิ่งใหญ่ไปสำหรับงานแก้สีปุ่ม หรือขยับ Layout นิดหน่อย การใช้ `tweak` จะสื่อความหมายได้ดีกว่าว่า "เป็นการปรับแต่งเล็กน้อย"

* **ความหมาย:** การบิด/ดัดแปลง/ปรับแต่ง (ไม่ใช่การสร้างใหม่)

```
*   **Pattern:** `tweak/<TICKET>-<description>`
```

* **ตัวอย่าง:**
    * `tweak/UI-501-adjust-padding` (ขยับช่องว่าง)
    * `tweak/UX-112-wording-change` (แก้คำผิด/แก้ข้อความ)
    * `tweak/CSS-99-dark-mode-color` (ปรับสี Dark mode นิดหน่อย)


### 3. ใช้ `ui/` หรือ `ux/` (เน้นงานหน้าบ้านโดยเฉพาะ)

สำหรับทีมที่มี Frontend หรือ Designer แยกชัดเจน อาจใช้ Prefix นี้เพื่อบอกว่าเป็นงานที่ไม่กระทบ Logic หรือ Database เลย

* **ความหมาย:** งานที่เกี่ยวกับหน้าตาล้วนๆ

```
*   **Pattern:** `ui/<TICKET>-<description>`
```

* **ตัวอย่าง:**
    * `ui/HOME-404-hero-banner` (เปลี่ยนรูป Banner)
    * `ui/MENU-202-icon-update` (เปลี่ยนไอคอน)

***

### 💡 เคล็ดลับ: การตั้งชื่อเมื่อเป็น "ส่วนย่อย" ของฟีเจอร์ใหญ่

บางครั้งฟีเจอร์เล็กๆ นั้นเป็นส่วนหนึ่งของโปรเจกต์ใหญ่ (เช่น ทำระบบ Login แต่แยก Branch มาทำแค่ปุ่ม "Lium Password")
แนะนำให้ใช้ **"Slash Grouping"** (การใช้ / ซ้อนกัน) เพื่อจัดกลุ่มใน Git GUI (เช่น SourceTree, GitKraken จะเห็นเป็นโฟลเดอร์)

* **Pattern:** `feature/<MAJOR-FEATURE>/<SUB-TASK>`
* **ตัวอย่าง:**
    * `feature/login/forget-password-link`
    * `feature/login/google-button`
    * `feature/dashboard/chart-widget`
    * `feature/dashboard/export-csv-btn`

**ข้อดี:** เวลาดูในโปรแกรม Git จะเห็น Branch ถูกจัดกลุ่มเป็นระเบียบ ไม่สะเปะสะปะครับ

### สรุปคำแนะนำ

* ถ้าอยากให้เป็นระเบียบที่สุด -> ใช้ **`feature/`** ทั้งหมด (Work สุดสำหรับ Automation)
* ถ้าอยากแยกความรู้สึกของงาน -> ใช้ **`tweak/`** สำหรับงานแก้เล็กๆ น้อยๆ
* ถ้าเป็นงานส่วนย่อยของงานใหญ่ -> ใช้ **`feature/main/sub`**

# ---------------------- 
 
# จะตั้งเกณฑ์อะไรให้ reviewer ตรวจโค้ดอย่างมีประสิทธิภาพ

การกำหนด **"เกณฑ์ (Criteria)"** ให้ Reviewer ใช้ตรวจงาน เป็นสิ่งสำคัญมาก เพื่อไม่ให้การตรวจกลายเป็นแค่เรื่อง "รสนิยมส่วนตัว" นี่คือ Checklist ที่แบ่งตามความสำคัญ (Must/Should/Could) ที่คุณสามารถนำไปแปะไว้ใน Template ของ Pull Request ได้เลยครับ

### 1. Functionality (ใช้งานได้จริงไหม?) — **[MUST]**

*ด่านแรกที่สำคัญที่สุด ถ้าไม่ผ่านข้อนี้ ห้าม Approve เด็ดขาด*

* [ ] **Requirements Met:** โค้ดทำงานตรงตาม Ticket หรือ User Story หรือไม่?
* [ ] **Edge Cases:** รองรับกรณีแปลกๆ หรือยัง? (เช่น ค่า NULL, ข้อมูลว่าง, เน็ตหลุด, User กรอกภาษาต่างดาว)
* [ ] **No Regression:** การแก้นี้ไปทำลายฟีเจอร์เก่าที่เคยทำงานได้หรือไม่?


### 2. Unit Test Quality (เทสมีคุณภาพไหม?) — **[MUST]**

*อ้างอิงจากบทความ EPT: ใช้หลักการ 5 ข้อมาจับ*

* [ ] **Arrange-Act-Assert:** เขียนเทสชัดเจนไหม? (เตรียมของ -> ทำ -> ตรวจผล)
* [ ] **Isolation:** เทสนี้ "แยกขาด" จริงไหม? (ต้องไม่มีการต่อ Database จริง, ไม่เรียก API ภายนอกจริง ต้องใช้ Mock เท่านั้น)
* [ ] **Coverage:** ครอบคลุมทั้งเคส "ปกติ" (Happy Path) และเคส "Error" (Unhappy Path) หรือยัง?
* [ ] **Readability:** ชื่อ Test Function อ่านแล้วรู้เรื่องทันทีไหมว่าเทสอะไร? (เช่น `should_throw_error_when_password_too_short` ไม่ใช่ `test_fail_1`)


### 3. Code Quality \& Readability (อ่านรู้เรื่องไหม?) — **[SHOULD]**

*เน้นความยั่งยืน (Maintainability) เพื่อให้คนอื่นมาแก้ต่อได้*

* [ ] **Naming:** ชื่อตัวแปร/ฟังก์ชัน สื่อความหมายชัดเจน ไม่ใช้ชื่อย่อที่รู้กันเอง (เช่น `x`, `data`, `temp` ❌ -> `userList`, `totalPrice` ✅)
* [ ] **Complexity:** ฟังก์ชันยาวเกินไปไหม? (ถ้าเกิน 20-30 บรรทัด ควรแยกฟังก์ชัน)
* [ ] **DRY (Don't Repeat Yourself):** มีการ Copy-Paste โค้ดเดิมซ้ำๆ ไหม? (ถ้ามีควรยุบเป็นฟังก์ชันกลาง)
* [ ] **Comments:** คอมเมนต์อธิบาย "ทำไม" (Why) ไม่ใช่อธิบายว่า "ทำอะไร" (What) (โค้ดที่ดีควรอธิบายตัวเองได้อยู่แล้ว)


### 4. Security \& Performance (ปลอดภัยและเร็วไหม?) — **[SHOULD]**

* [ ] **SQL Injection:** มีการต่อ String ใน SQL ตรงๆ ไหม? (ต้องใช้ ORM หรือ Parameterized Query)
* [ ] **Sensitive Data:** เผลอ Hardcode รหัสผ่านหรือ API Key ลงไปในโค้ดหรือเปล่า? ☠️
* [ ] **N+1 Problem:** มีการ Loop เรียก Database ทีละ row ไหม? (ถ้ามีควรแก้เป็น Query ทีเดียว)


### 5. Style \& Housekeeping (เรื่องจุกจิก) — **[COULD]**

*เรื่องพวกนี้ควรให้ "ระบบอัตโนมัติ (Linter/Prettier)" จัดการแทนคน*

* [ ] **Formatting:** เว้นวรรค, ปีกกา, ย่อหน้า ตรงตามมาตรฐานทีมไหม?
* [ ] **Unused Code:** มีตัวแปรที่ประกาศไว้แต่ไม่ได้ใช้ หรือ `console.log` ที่ลืมลบไหม?

***

### ตัวอย่าง Comment ที่ดี vs ไม่ดี (สำหรับ Reviewer)

| ❌ ไม่ดี (Vague/Rude) | ✅ ดี (Specific/Constructive) |
| :-- | :-- |
| "โค้ดแย่มาก ไปแก้มาใหม่" | "ฟังก์ชันนี้ดูซับซ้อนไปนิด ลองแยก Logic ส่วนคำนวณภาษีออกมาเป็นอีกฟังก์ชันดีไหมครับ? จะได้เทสง่ายขึ้น" |
| "ชื่อตัวแปรงง" | "ตัวแปร `d` สื่อความหมายไม่ชัดเจน แนะนำให้เปลี่ยนเป็น `transactionDate` เพื่อให้อ่านง่ายขึ้นครับ" |
| "ทำไมไม่เขียนเทส?" | "รบกวนเพิ่ม Unit Test สำหรับเคสที่ 'User ไม่ได้ Login' ด้วยครับ เพื่อให้ Coverage ครอบคลุมตามเกณฑ์" |

### คำแนะนำสำหรับการนำไปใช้

ให้สร้าง **Pull Request Template** ใน GitLab/GitHub แล้วใส่ Checklist นี้ลงไป เพื่อให้ Developer ต้องติ๊กยืนยันก่อนกดส่ง PR ครับ

```markdown
## Reviewer Checklist
- [ ] Functionality: ทำงานถูกต้องตาม Requirement และรองรับ Edge Cases
- [ ] Testing: มี Unit Test ครอบคลุม (Arrange, Act, Assert, Isolation)
- [ ] Readability: ตั้งชื่อตัวแปรชัดเจน ไม่ซับซ้อนเกินไป
- [ ] Security: ไม่มี Sensitive Data หรือช่องโหว่
```

# สรุป Code Review Workflow แบบเข้าใจง่าย

1) Code Assistant — เขียนโค้ดด้วยเครื่องมือช่วย
• ใช้ AI/Code Assistant ช่วยเขียนโค้ดให้เร็วขึ้น
• ตรวจสไตล์โค้ดเบื้องต้นก่อนส่งขึ้นระบบ
เริ่มต้นให้โค้ดพร้อมตรวจ
⸻
2) Pull Request — ส่งโค้ดให้ทีมตรวจ
• เปิด Pull Request พร้อมคำอธิบาย
• แนบรายการเปลี่ยนแปลง, Issue ที่เกี่ยวข้อง
• Reviewer สามารถเข้ามาตรวจได้ทันที
เป็นขั้นตอนขอให้ทีมช่วยตรวจ
⸻
3) CI Pipeline — ตรวจอัตโนมัติ
• ระบบรัน Test / Lint / Build อัตโนมัติ
• ถ้าไม่ผ่าน จะต้องแก้ไขก่อนเข้าสู่ Code Review
ป้องกันโค้ดเสียตั้งแต่ต้นทาง
⸻
4) Code Review — ทีมตรวจโค้ด
• Reviewer ตรวจคุณภาพโค้ด ความถูกต้อง ความเหมาะสม
• อาจมีข้อเสนอแนะหรือให้แก้ไข
ให้แน่ใจว่าคุณภาพได้มาตรฐานทีม
⸻
5) Code Update — แก้ไขตามรีวิว
• ปรับตามคำแนะนำ
• Push โค้ดขึ้นมาใหม่ให้ระบบ rerun
รอบแก้ไขเพื่อให้โค้ดสมบูรณ์ขึ้น
⸻
6) Deploy — ปล่อยสู่ระบบ
• เมื่อรีวิวผ่านและ CI ผ่าน
• โค้ดถูก merge และนำไป deploy
ก้าวสู่ production อย่างมั่นใจ
⸻
เหมาะกับใคร?
• Developer ทุกระดับ
• ทีมที่ต้องการ workflow ที่ชัดเจน
• ผู้เริ่มต้นที่อยากเรียนรู้ขั้นตอน CI/CD + Code Review
• Tech Lead ที่ต้องการอธิบาย Flow ให้ทีมใหม่เข้าใจเร็ว

นี่คือสรุป **Code Review Workflow แบบเข้าใจง่าย** สำหรับทีมที่ต้องการกระบวนการที่มีคุณภาพและรวดเร็ว ตั้งแต่เริ่มเขียนจนถึงขึ้นระบบ โดยเน้นความชัดเจนในแต่ละขั้นตอน:

### 1. Code Assistant — เขียนโค้ดด้วยตัวช่วย

* **หน้าที่:** Developer
* **การกระทำ:**
    * ใช้ AI หรือ Code Assistant ช่วยร่างโค้ดและตรวจสอบ Syntax เบื้องต้น
    * รัน Unit Test (local) ด้วยตัวเองให้ผ่านทุกข้อ
    * จัด Format โค้ดให้สวยงามตามมาตรฐานทีม (Linting)
* **เป้าหมาย:** ส่งโค้ดที่ "สะอาด" และ "ทำงานได้" เข้าสู่ระบบ ลดภาระคนตรวจ


### 2. Pull Request (PR) — ส่งโค้ดให้ทีมตรวจ

* **หน้าที่:** Developer
* **การกระทำ:**
    * สร้าง PR/MR เข้า Branch หลัก (เช่น `develop`)
    * **สำคัญ:** เขียนคำอธิบาย PR ให้ชัดเจน (ทำอะไร? เพื่อแก้ Ticket ไหน? มีผลกระทบอะไร?)
    * แนบรูปภาพหรือผลเทสประกอบถ้ามี
* **เป้าหมาย:** แจ้งทีมว่า "งานเสร็จแล้ว ช่วยมาดูหน่อย"


### 3. CI Pipeline — ตรวจอัตโนมัติ (ด่านหน้า)

* **หน้าที่:** ระบบอัตโนมัติ (System)
* **การกระทำ:**
    * ทันทีที่เปิด PR ระบบจะรัน Test, Lint, และ Build Docker Image
    * เช็ค Code Coverage (ต้องผ่านเกณฑ์ที่ตั้งไว้)
    * *ถ้าไม่ผ่าน:* ระบบจะ Block ไม่ให้ Merge และแจ้งเตือน Developer ให้ไปแก้ก่อน
* **เป้าหมาย:** คัดกรอง Error พื้นฐานออกไป ไม่ให้เสียเวลาคนตรวจ


### 4. Code Review — ทีมตรวจคุณภาพ

* **หน้าที่:** Reviewer (Senior/Lead/Peer)
* **การกระทำ:**
    * อ่าน Logic ว่าถูกต้องและปลอดภัยไหม
    * เช็คความอ่านง่าย (Readability) และการตั้งชื่อตัวแปร
    * ดูว่า Unit Test ครอบคลุมและแยกส่วน (Isolation) จริงหรือไม่
    * ให้ Comment แนะนำจุดที่ควรปรับ (เน้นสร้างสรรค์ ไม่ใช่จับผิด)
* **เป้าหมาย:** รักษามาตรฐานทีม และแชร์ความรู้ร่วมกัน


### 5. Code Update — แก้ไขตามรีวิว

* **หน้าที่:** Developer
* **การกระทำ:**
    * อ่าน Comment และแก้ไขโค้ดตามคำแนะนำ
    * Push โค้ดใหม่ขึ้นไป (ระบบ CI จะรันซ้ำอัตโนมัติ)
    * Reply ใน PR ว่า "แก้แล้วครับ" เพื่อเรียกคนตรวจกลับมาดู
* **เป้าหมาย:** ปรับปรุงงานให้สมบูรณ์ที่สุดก่อนรวมร่าง


### 6. Deploy — ปล่อยสู่ระบบ

* **หน้าที่:** Team Lead / System
* **การกระทำ:**
    * เมื่อ Reviewer กด **Approve** และ CI เป็นสีเขียว ✅
    * กดปุ่ม **Merge** โค้ดเข้าสู่ Branch หลัก
    * ระบบ CD (Continuous Deployment) จะนำโค้ดใหม่ไป Deploy ขึ้น Server (Staging/Production) ทันที
* **เป้าหมาย:** ส่งมอบคุณค่า (Value) ให้ผู้ใช้งานอย่างรวดเร็วและมั่นใจ

***

### แผนภาพสรุป (Workflow Diagram)

```mermaid
graph TD
    A[1. Code Assistant<br/>(Dev เขียนโค้ด)] -->|Push Code| B[3. CI Pipeline<br/>(Auto Test/Lint)]
    B -->|❌ ไม่ผ่าน| A
    B -->|✅ ผ่าน| C[2. Pull Request<br/>(ขอตรวจงาน)]
    C --> D[4. Code Review<br/>(ทีมช่วยดู)]
    D -->|❌ ขอแก้| E[5. Code Update<br/>(Dev แก้ตามสั่ง)]
    E -->|Push ใหม่| B
    D -->|✅ Approve| F[6. Deploy<br/>(Merge & Run)]
```

**เหมาะกับใคร?**

* **Developer ทุกระดับ:** เพื่อให้รู้หน้าที่ตัวเองในแต่ละขั้นตอน
* **ทีมที่ต้องการมาตรฐาน:** ลดการ "หลุด" ของบั๊กง่ายๆ
* **Tech Lead:** ใช้สอนน้องใหม่ในทีมให้เข้าใจ Flow การทำงานร่วมกันได้อย่างรวดเร็ว

# เกณฑ์สำคัญ 10 ข้อสำหรับการประเมินโค้ดโดย reviewer

**เกณฑ์สำคัญ 10 ข้อสำหรับการประเมินโค้ด (Code Review Criteria)**
เพื่อช่วยให้ Reviewer ตรวจงานได้อย่างมีทิศทาง ลดความขัดแย้ง และยกระดับคุณภาพซอฟต์แวร์ นี่คือ Checklist 10 ข้อที่ครอบคลุมทั้ง Functionality, Quality และ Testing ครับ:

***

### หมวดที่ 1: ความถูกต้องและการใช้งาน (Does it work?)

**1. ตรงตาม Requirement (Correctness):**

* โค้ดทำงานถูกต้องตาม Ticket/User Story หรือไม่?
* Logic การคำนวณหรือ Flow การทำงานถูกต้องตาม Business Rule ไหม?
* *คำถาม:* "ถ้า User ทำตาม Step 1-2-3 ผลลัพธ์ออกมาถูกเป๊ะไหม?"

**2. รองรับ Edge Cases (Robustness):**

* ทดสอบกรณี "ข้อมูลแปลกๆ" หรือยัง? (เช่น ข้อมูลเป็น Null, Array ว่าง, User กรอก Emoji, เน็ตหลุดกลางทาง)
* มีการจัดการ Error (Error Handling) ที่เหมาะสมหรือไม่? ไม่ใช่แค่ `try-catch` ทิ้งไว้เฉยๆ

**3. ความปลอดภัย (Security):**

* มีการตรวจสอบข้อมูลนำเข้า (Input Validation) ไหม?
* มีช่องโหว่พื้นฐานหรือไม่? (เช่น SQL Injection, XSS, หรือเผลอ Hardcode Password/API Key ลงไปในโค้ด)

***

### หมวดที่ 2: คุณภาพโค้ด (Is it clean?)

**4. อ่านง่ายและสื่อความหมาย (Readability \& Naming):**

* ชื่อตัวแปร ฟังก์ชัน และคลาส สื่อความหมายชัดเจนหรือไม่? (เช่น `d` ❌ vs `daysSinceLastLogin` ✅)
* โครงสร้างโค้ดซับซ้อนเกินไปไหม? (Cyclomatic Complexity) ถ้าอ่านแล้วต้องขมวดคิ้วเกิน 3 วิ แสดงว่าควรแก้

**5. ไม่ทำงานซ้ำซ้อน (DRY - Don't Repeat Yourself):**

* มีการ Copy-Paste Logic เดิมไปแปะหลายที่ไหม?
* ถ้ามี ควรยุบรวมเป็น Function กลาง หรือ Component ที่ใช้ร่วมกันได้

**6. ประสิทธิภาพ (Performance):**

* มีการ Loop ซ้อน Loop (O(n^2)) โดยไม่จำเป็นไหม?
* มีการ Query Database ใน Loop (N+1 Problem) หรือไม่?
* มีการโหลดข้อมูลมาเยอะเกินความจำเป็นไหม? (เช่น `SELECT *` แต่ใช้แค่ 2 fields)

**7. สไตล์และมาตรฐาน (Coding Standard):**

* การจัด Format (เว้นวรรค, ย่อหน้า) ตรงตามมาตรฐานทีมหรือ Linter ไหม?
* โครงสร้าง Folder/File ถูกต้องตาม Architecture ของโปรเจกต์ไหม?

***

### หมวดที่ 3: การทดสอบและการดูแลรักษา (Can we maintain it?)

**8. การทดสอบ (Test Coverage \& Quality):**

* มี Unit Test ครอบคลุม Logic ใหม่หรือไม่?
* Test เขียนตามหลัก Arrange-Act-Assert และ Isolation (Mock dependency) หรือไม่?
* Test เคส Unhappy Path (กรณี Error) ด้วยหรือเปล่า?

**9. ไม่กระทบของเดิม (No Regression):**

* การแก้นี้ไปทำให้ฟีเจอร์เก่าที่เคยดีอยู่...พังไหม?
* ควรเช็คว่ามีการแก้ไขไฟล์ที่ไม่เกี่ยวข้องโดยไม่ตั้งใจหรือไม่

**10. เอกสารประกอบ (Documentation):**

* ถ้ามีการแก้ API มีการอัปเดต Swagger/Postman หรือยัง?
* ถ้า Logic ซับซ้อนมาก มี Comment อธิบาย "Why" (ทำไมถึงเขียนแบบนี้) ไว้ไหม?

***

### 💡 เคล็ดลับสำหรับ Reviewer

* **Be Constructive:** วิจารณ์ที่ "โค้ด" ไม่ใช่ "คน" (เช่น "ตรงนี้อาจทำให้ช้า" แทน "ทำไมเขียนแบบนี้")
* **Nitpicks:** เรื่องเล็กน้อย (เช่น ลืมลบ console.log) ให้ระบุว่าเป็น "Nitpick" (แก้ก็ดี ไม่แก้ก็ได้) เพื่อไม่ให้ผู้ถูกตรวจรู้สึกกดดัน
* **Approve with Comments:** ถ้ามีแก้เล็กน้อย ให้ Approve ไปเลยแต่ฝากแก้ด้วย เพื่อไม่ให้งานสะดุด
 

 
# ตัวอย่างคำถามที่ reviewer ควรถามเมื่อรีวิวโค้ด

คำถามที่ดีคือเครื่องมือที่ทรงพลังที่สุดของ Reviewer ครับ เพราะมันช่วย "กระตุ้นให้คิด" มากกว่า "ออกคำสั่ง" ทำให้ Developer ไม่รู้สึกเหมือนโดนจับผิด

นี่คือตัวอย่างคำถามที่ Reviewer ควรใช้ แบ่งตามสถานการณ์ เพื่อให้ได้โค้ดที่มีคุณภาพและทีมมีความสุข:

### 1. เมื่อสงสัยใน Logic หรือความซับซ้อน (Complexity)

แทนที่จะบอกว่า "เขียนงงมาก" ให้ถามว่า:

* "ช่วยอธิบาย Flow ตรงนี้ให้ฟังหน่อยได้ไหมครับ ว่ามันทำงานยังไง?" (ให้เขาเล่า Logic เอง อาจจะเจอจุดผิดเอง)
* "ถ้า Input เป็นค่า [X] หรือ [Null] ฟังก์ชันนี้จะยังทำงานถูกไหม?" (ชวนคิดเรื่อง Edge Cases)
* "ตรงนี้ถ้าเราแตกเป็นฟังก์ชันย่อยออกมา จะทำให้อ่านง่ายขึ้นไหม หรือมีเหตุผลอะไรที่ต้องรวมไว้ที่เดียว?"
* "มีวิธีอื่นที่เขียนสั้นกว่านี้ไหม หรือแบบนี้คือดีที่สุดแล้วในมุมมองของคุณ?"


### 2. เมื่อกังวลเรื่อง Performance (Performance)

แทนที่จะบอกว่า "ช้าแน่ๆ" ให้ถามว่า:

* "ถ้าข้อมูลใน Database โตขึ้นเป็นแสน record ตรงนี้จะมีปัญหาไหม?"
* "เราจำเป็นต้อง Loop ตรงนี้ทุกรอบไหม หรือ Cache ไว้ได้?"
* "Query นี้มีโอกาสเกิด N+1 ปัญหาไหมครับ?"


### 3. เมื่อโค้ดดูไม่ปลอดภัย (Security)

แทนที่จะบอกว่า "ไม่ปลอดภัย" ให้ถามว่า:

* "เรามั่นใจได้ยังไงว่า Input ตัวนี้ปลอดภัยจากการถูก Hack (เช่น SQL Injection)?"
* "ถ้า User คนอื่นมาเรียก API นี้ เขาจะเห็นข้อมูลของคนอื่นไหม?"
* "เราควรซ่อน Sensitive Data ตรงนี้ใน Log ไหมครับ?"


### 4. เมื่ออยากให้เพิ่ม Test (Testing)

แทนที่จะบอกว่า "ไปเขียนเทสมา" ให้ถามว่า:

* "เราจะมั่นใจได้ยังไงว่า Logic นี้ทำงานถูก ถ้าในอนาคตมีคนมาแก้โค้ดบรรทัดนี้?"
* "มี Test Case ไหนที่ครอบคลุมกรณี Error นี้หรือยังครับ?"
* "ส่วนนี้ Mock dependency ไว้หรือยัง หรือว่าต่อ Database จริง?"


### 5. เมื่อดูแล้ว "ดีแล้ว" แต่อยากแนะแนวทาง (Suggestion)

* "อันนี้ Logic ดีแล้วครับ แต่ถ้าใช้ [Library X / Function Y] อาจจะประหยัดบรรทัดได้อีก สนใจลองดูไหม?"
* "ตั้งชื่อตัวแปรแบบนี้ก็โอเคครับ แต่ถ้าเปลี่ยนเป็น [ชื่อใหม่] จะสื่อความหมายชัดกว่าไหม?"

***

### 💡 Tip: เทคนิคการตั้งคำถามที่ดี

1. **ถาม "Why" ไม่ใช่ "What":** ถามหาเหตุผล ("ทำไมถึงเลือกวิธีนี้?") แทนที่จะถามว่าทำอะไร
2. **ใช้ "เรา" แทน "คุณ":** "ตรงนี้ **เรา** จะปรับให้เร็วขึ้นได้ไหม?" (ให้ความรู้สึกเป็นทีมเดียวกัน)
3. **เสนอทางเลือก:** "ถ้าลองทำแบบ A หรือ B คิดว่าแบบไหนดีกว่ากันครับ?"

- การใช้คำถามแบบนี้จะเปลี่ยนบรรยากาศ Code Review จากการ **"สอบสวน"** ให้กลายเป็นการ **"ปรึกษาหารือ"** (Discussion) ซึ่งดีต่อสุขภาพจิตของทีมมาก 



# คำถามเฉพาะสำหรับรีวิวความปลอดภัยของโค้ด

สำหรับเรื่องความปลอดภัย (Security) ซึ่งเป็นจุดตายที่สำคัญมาก นี่คือชุดคำถาม **"Security-Focused Questions"** ที่ Reviewer ควรใช้จี้จุด Developer โดยอ้างอิงจากมาตรฐาน OWASP และ Secure Coding Practice:

### 1. Input Validation (ข้อมูลขาเข้า)

*สำคัญที่สุด เพราะ 80% ของการโดน Hack มาจากช่องนี้*

* "ตัวแปร `userInput` นี้ เรามีการ Validate หรือ Sanitize ก่อนนำไปใช้ไหมครับ? (ป้องกัน XSS/Injection)"
* "ถ้าผมส่งค่า `null`, ค่าติดลบ หรือ String ยาว 1 ล้านตัวอักษรเข้ามา ระบบจะพังไหม?"
* "ตรงนี้รับ File Upload มีการเช็ค Mime-Type และนามสกุลไฟล์จริงๆ ไหม? (ไม่ใช่แค่เช็คชื่อไฟล์)"


### 2. Authentication \& Authorization (ใครเป็นใคร ทำอะไรได้บ้าง)

* "API Endpoint นี้ มีการเช็ค Permission ไหมว่า User คนนี้มีสิทธิ์เรียกจริงๆ? (ป้องกัน IDOR)"
* "ถ้าผมเปลี่ยน `userId` ใน URL เป็นของคนอื่น ผมจะเห็นข้อมูลของเขาไหม?"
* "ทำไมเราต้องส่ง `password` หรือ `token` กลับไปใน Response ด้วยครับ? (ควรเอาออก)"


### 3. Data Protection (การปกป้องข้อมูล)

* "Log บรรทัดนี้ มีการปริ้นท์ข้อมูลส่วนตัว (PII) เช่น บัตรประชาชน หรือเบอร์โทร ลงไปไหม?"
* "ค่า `API_KEY` นี้ Hardcode ไว้ในโค้ดหรือเปล่า? (ควรย้ายไป Environment Variable)"
* "ข้อมูลที่ส่งออกไปหน้าบ้าน มีการ Encrypt หรือ Hash ข้อมูลสำคัญหรือยัง?"


### 4. Dependency \& Libraries (ของแถม)

* "Library ตัวนี้ที่เอามาใช้ เป็นเวอร์ชันล่าสุดที่มี Patch ความปลอดภัยหรือยัง?"
* "เราจำเป็นต้องใช้ Library ใหญ่ขนาดนี้ เพื่อทำฟังก์ชันแค่นี้จริงเหรอ? (ลด Attack Surface)"


### 5. Logic \& Business Rules (ตรรกะทางธุรกิจ)

* "ใน Flow การจ่ายเงิน ถ้าเน็ตหลุดตอนตัดบัตรเครดิตเสร็จ แต่ยังไม่บันทึก Database เงินจะหายไหม?"
* "มี Rate Limit ป้องกันคนยิง API รัวๆ หรือยัง?"

***

### 💡 Tip: เทคนิคการรีวิว Security ให้ไม่เครียด

ใช้คำถามเชิงสมมติสถานการณ์ (Hypothetical Scenarios) จะช่วยให้ Developer เห็นภาพความเสี่ยงได้ชัดขึ้นครับ เช่น:
> *"ถ้าสมมติผมเป็น Hacker แล้วลองส่ง Script `<script>alert('hack')</script>` เข้ามาที่ช่อง Comment นี้ หน้าเว็บของ User คนอื่นจะเด้ง Popup ไหมครับ?"*

คำถามแบบนี้จะทำให้ Developer รีบกลับไปเช็คโค้ดตัวเองทันทีโดยที่เราไม่ต้องด่าเลยครับ 😉

# ## 

# รายการเช็คลิสต์ตรวจหาช่องโหว่ทั่วไปในโค้ด

นี่คือรายการเช็คลิสต์ (Checklist) ที่สรุปมาจากมาตรฐาน **OWASP Top 10** และ **CWE Top 25** เพื่อใช้ตรวจหาช่องโหว่ทั่วไปในโค้ดได้อย่างครอบคลุม:

### 1. การตรวจสอบข้อมูลขาเข้า (Input Validation)

*สาเหตุอันดับ 1 ของการถูกแฮก (เช่น SQL Injection, XSS)*

* [ ] **Type Checking:** ตัวแปรรับค่าถูกประเภทไหม? (เช่น รับตัวเลขต้องเป็น Int ไม่ใช่ String)
* [ ] **Length Check:** จำกัดความยาว Input หรือยัง? (ป้องกัน Buffer Overflow)
* [ ] **Allowlist:** ตรวจสอบค่าที่ยอมรับเท่านั้น (White-listing) แทนที่จะแบนค่าที่ห้าม (Black-listing)
* [ ] **Sanitization:** ลบอักขระพิเศษที่อันตรายออกก่อนนำไปใช้ (เช่น `<script>`, `'`, `--`)


### 2. การยืนยันตัวตนและสิทธิ์ (Authentication \& Authorization)

* [ ] **Broken Access Control:** User A สามารถแก้ URL เพื่อดูข้อมูล User B ได้ไหม? (IDOR)
* [ ] **Permission Check:** ทุก API Endpoint มีการเช็ค Role/Permission ก่อนทำงานเสมอไหม?
* [ ] **No Hardcoded Credential:** ไม่มี Username/Password หรือ API Key ฝังอยู่ในโค้ด
* [ ] **Session Management:** Session ID ถูกสร้างใหม่ทุกครั้งที่ Login ไหม? (ป้องกัน Session Fixation)


### 3. การจัดการข้อมูลสำคัญ (Sensitive Data Exposure)

* [ ] **Encryption:** รหัสผ่านถูก Hash ด้วย Algorithm ที่ปลอดภัย (เช่น Argon2, Bcrypt) หรือไม่? (ห้ามใช้ MD5/SHA1)
* [ ] **No Logging Secrets:** มั่นใจว่าไม่มีการ `console.log` หรือเขียน Log ข้อมูลบัตรเครดิต/รหัสผ่าน
* [ ] **HTTPS Only:** บังคับใช้ HTTPS เท่านั้น ไม่มีการส่งข้อมูลผ่าน HTTP ธรรมดา


### 4. ความปลอดภัยของโค้ดและไลบรารี (Vulnerable Components)

* [ ] **Outdated Libraries:** ไลบรารีที่ใช้ (npm, pip, maven) เป็นเวอร์ชันล่าสุดที่แพทช์ช่องโหว่แล้วหรือยัง?
* [ ] **Unused Code:** ลบโค้ดเก่าที่ไม่ได้ใช้ (Dead Code) ออกเพื่อลดช่องทางโจมตี


### 5. การจัดการข้อผิดพลาด (Error Handling)

* [ ] **Generic Error Message:** เมื่อระบบ Error ต้องไม่โชว์ Stack Trace หรือข้อมูล Database ให้ User เห็น (โชว์แค่ "เกิดข้อผิดพลาด กรุณาติดต่อแอดมิน")
* [ ] **Fail Safe:** ถ้าระบบล่ม ต้องล่มในสถานะที่ "ปิด" (Deny Access) ไม่ใช่ "เปิด" ให้เข้าได้ทุกคน

***

### 🛠 เครื่องมือช่วยสแกนอัตโนมัติ (แนะนำให้ใช้เสริม)

* **SonarQube:** ตรวจคุณภาพโค้ดและช่องโหว่พื้นฐาน
* **OWASP ZAP:** สแกนเว็บหาช่องโหว่แบบ Black-box
* **Snyk / npm audit:** ตรวจหา Library ที่มีช่องโหว่

**คำแนะนำ:** นำ Checklist นี้ไปใส่ใน Pull Request Template หมวด "Security Review" เพื่อเตือนสติ Developer ก่อน Merge ครับ

### คู่มือเกณฑ์การทำ Code Review
---------------
####  1. หลักการพื้นฐาน

### 1.1 จุดประสงค์ของ Code Review
- **ปรับปรุงคุณภาพโค้ด** - ทำให้โค้ดอ่านง่าย, บำรุงรักษาง่าย
- **แบ่งปันความรู้** - เพิ่มความเข้าใจในระบบ across team
- **ตรวจสอบความถูกต้อง** - ลด bugs และ security issues
- **รักษามาตรฐาน** - ให้โค้ดสอดคล้องกับ coding standards ของทีม

### 1.2 Mindset ที่ควรมี
- **เป็นผู้ช่วย ไม่ใช่ผู้พิพากษา** - ใช้คำถามมากกว่าการสั่ง
- **เคารพผู้เขียน** - โฟกัสที่โค้ด ไม่ใช่ตัวบุคคล
- **เปิดใจรับฟัง** - ทั้งผู้ review และผู้เขียน
- **เน้นการเรียนรู้** - ทุก review เป็นโอกาสในการเรียนรู้

## 2. เกณฑ์การประเมินหลัก

### 2.1 ความถูกต้อง (Correctness)
- [ ] **Logic ถูกต้อง** - โค้ดทำงานตาม requirement จริงหรือไม่
- [ ] **Edge cases** - จัดการกับ corner cases อย่างเหมาะสม
- [ ] **Error handling** - มีการจัดการข้อผิดพลาดที่ครอบคลุม
- [ ] **Business logic** - ตรงตาม business rules

### 2.2 ความปลอดภัย (Security)
- [ ] **Input validation** - validate input ทุกช่องทาง
- [ ] **Authentication/Authorization** - ตรวจสอบสิทธิ์อย่างเหมาะสม
- [ ] **Data protection** - ไม่ expose sensitive data
- [ ] **SQL injection/XSS** - ป้องกัน vulnerability พื้นฐาน

### 2.3 ประสิทธิภาพ (Performance)
- [ ] **Algorithm efficiency** - ใช้ algorithm ที่เหมาะสม
- [ ] **Database queries** - query มีประสิทธิภาพ (ใช้ index, ไม่มี N+1)
- [ ] **Memory usage** - ไม่มี memory leak
- [ ] **Response time** - อยู่ในเกณฑ์ที่ยอมรับได้

### 2.4 การทดสอบ (Testing)
- [ ] **Unit tests** - มี test coverage ที่เหมาะสม
- [ ] **Test cases** - ครอบคลุมทั้ง happy path และ edge cases
- [ ] **Test readability** - test อ่านเข้าใจง่าย
- [ ] **Integration tests** - (ถ้าจำเป็น) สำหรับ critical flows

### 2.5 การออกแบบ (Design)
- [ ] **Separation of concerns** - แต่ละ module/class มีหน้าที่ชัดเจน
- [ ] **SOLID principles** - ใช้ principles พื้นฐานอย่างเหมาะสม
- [ ] **Design patterns** - ใช้ patterns เมื่อเหมาะสม (แต่ไม่ over-engineer)
- [ ] **Dependencies** - การพึ่งพาระหว่าง components มีเหตุผล

## 3. คุณภาพของโค้ด

### 3.1 การอ่านเข้าใจ (Readability)
- [ ] **การตั้งชื่อ** - ตัวแปร, ฟังก์ชัน, class ชื่อสื่อความหมาย
- [ ] **ความซับซ้อน** - ฟังก์ชันไม่ยาวเกินไป (แนะนำ < 20-30 lines)
- [ ] **Comment** - มี comment เมื่อจำเป็น (อธิบาย why ไม่ใช่ what)
- [ ] **Consistency** - สอดคล้องกับ style ของโปรเจค

### 3.2 การบำรุงรักษา (Maintainability)
- [ ] **Duplication** - ไม่มี code ซ้ำซ้อน (DRY principle)
- [ ] **Complexity** - cyclomatic complexity ไม่สูงเกินไป
- [ ] **Modularity** - แยกส่วนที่ reuse ได้
- [ ] **Configuration** - hard-coded values น้อยที่สุด

### 3.3 Coding Standards
- [ ] **Formatting** - ตาม convention ของภาษาและทีม
- [ ] **Language features** - ใช้ features ใหม่เมื่อเหมาะสม
- [ ] **Best practices** - ตาม community standards
- [ ] **Linting rules** - ผ่าน linting rules ที่กำหนด

## 4. กระบวนการ Review

### 4.1 สำหรับผู้ Review
**ควรทำ:**
- Review ภายใน timeline ที่กำหนด (แนะนำ < 24 ชม.)
- ให้ feedback ที่เป็น constructive
- ชื่นชมจุดที่ดีของโค้ด
- ถามคำถามเมื่อไม่เข้าใจ
- ตรวจสอบทั้ง implementation และ tests

**ไม่ควรทำ:**
- คิดแทนผู้เขียน (micromanage)
- จับผิดเรื่อง formatting เล็กน้อย
- บังคับให้ refactor โดยไม่มีเหตุผลสำคัญ
- ล่าช้าโดยไม่จำเป็น

### 4.2 สำหรับผู้ขอ Review
**ควรเตรียม:**
- คำอธิบายการเปลี่ยนแปลง (PR description)
- Link ไปถึง requirement/ticket
- Test instructions (ถ้าจำเป็น)
- ระบุจุดที่ต้องการ feedback เป็นพิเศษ

**ควรตอบสนอง:**
- พิจารณา feedback อย่างเปิดใจ
- อธิบาย reasoning เมื่อไม่เห็นด้วย
- ขอ clarification เมื่อไม่เข้าใจ feedback
- ขอบคุณ reviewer

## 5. Checklist อย่างรวดเร็ว

### 5.1 Pre-review
- [ ] โค้ด compile/run ได้
- [ ] Tests ผ่านทั้งหมด
- [ ] CI/CD pipeline ผ่าน
- [ ] Documentation updated (ถ้าจำเป็น)

### 5.2 During Review
- [ ] Security issues
- [ ] Performance problems
- [ ] Major bugs
- [ ] Architectural concerns
- [ ] Test coverage

### 5.3 Before Approval
- [ ] All comments addressed/resolved
- [ ] No regression introduced
- [ ] Meets acceptance criteria
- [ ] Ready for deployment

## 6. การให้ Feedback ที่มีประสิทธิภาพ

### 6.1 เทคนิคการเขียน comment
- **ใช้คำถาม**: "คุณคิดว่า approach นี้จะทำงานกับ case X อย่างไร?"
- **ให้ตัวอย่าง**: "อาจลองใช้วิธีนี้: `const result = data.filter(x => x.active)`"
- **อ้างอิง standards**: "ตาม style guide เราใช้ camelCase สำหรับตัวแปร"
- **ระบุระดับความสำคัญ**: 
  - MUST (critical): ต้องแก้ก่อน merge
  - SHOULD (important): แนะนำให้แก้
  - COULD (minor): optional improvement

### 6.2 การจัดการความเห็นต่าง
1. **อภิปรายด้วยข้อมูล** - อ้างอิง facts, metrics, benchmarks
2. **ปรึกษา third party** - ขอความเห็นจาก senior/tech lead เมื่อจำเป็น
3. **พิจารณา trade-offs** - ทุกทางเลือกมีข้อดีข้อเสีย
4. **ตัดสินใจและเดินต่อไป** - ไม่ติดอยู่กับ perfectionism

## 7. Metrics และการปรับปรุง

### 7.1 วัดประสิทธิภาพ
- **Cycle time** - เวลาจากเปิด PR ถึง merge
- **Review depth** - จำนวน comments และ discussion
- **Defect rate** - bugs ที่พบหลัง deployment
- **Team satisfaction** - survey ความพึงพอใจ

### 7.2 การปรับปรุงกระบวนการ
- **Regular retrospectives** - พูดคุยปรับปรุงกระบวนการ
- **Pair reviewing** - สำหรับ complex changes
- **Rotation** - หมุนเวียน reviewer
- **Training** - แบ่งปัน best practices

---

### สรุป
#### Code review ที่ดีคือ **การสนทนาระหว่าง professionals** เพื่อสร้างโค้ดที่ดีที่สุดเท่าที่จะเป็นไปได้ เน้นที่:
1. **คุณภาพ** มากกว่าความเร็ว
2. **การเรียนรู้** มากกว่าการจับผิด
3. **ความร่วมมือ** มากกว่าการแข่งขัน

#### **Remember:** เราทุกคนอยู่ในทีมเดียวกัน เป้าหมายคือสร้าง software ที่ดี ไม่ใชว่าใครเก่งกว่าใคร
---------------
# Code Review Criteria Handbook

## 1. Fundamental Principles

### 1.1 Purpose of Code Review
- **Improve code quality** - Make code readable, maintainable
- **Share knowledge** - Increase system understanding across the team
- **Verify correctness** - Reduce bugs and security issues
- **Maintain standards** - Ensure code aligns with team coding standards

### 1.2 Recommended Mindset
- **Be a helper, not a judge** - Use questions more than commands
- **Respect the author** - Focus on the code, not the person
- **Be open-minded** - Both reviewers and authors should listen
- **Focus on learning** - Every review is a learning opportunity

## 2. Primary Evaluation Criteria

### 2.1 Correctness
- [ ] **Logic is correct** - Does code work according to requirements?
- [ ] **Edge cases** - Handles corner cases appropriately
- [ ] **Error handling** - Comprehensive error management
- [ ] **Business logic** - Follows business rules correctly

### 2.2 Security
- [ ] **Input validation** - Validate all input channels
- [ ] **Authentication/Authorization** - Proper permission checking
- [ ] **Data protection** - No sensitive data exposure
- [ ] **SQL injection/XSS** - Prevent basic vulnerabilities

### 2.3 Performance
- [ ] **Algorithm efficiency** - Uses appropriate algorithms
- [ ] **Database queries** - Efficient queries (use indexes, no N+1)
- [ ] **Memory usage** - No memory leaks
- [ ] **Response time** - Within acceptable limits

### 2.4 Testing
- [ ] **Unit tests** - Appropriate test coverage
- [ ] **Test cases** - Covers both happy path and edge cases
- [ ] **Test readability** - Tests are easy to understand
- [ ] **Integration tests** - (If needed) for critical flows

### 2.5 Design
- [ ] **Separation of concerns** - Each module/class has clear responsibility
- [ ] **SOLID principles** - Appropriately applies basic principles
- [ ] **Design patterns** - Uses patterns when appropriate (not over-engineered)
- [ ] **Dependencies** - Reasonable dependencies between components

## 3. Code Quality

### 3.1 Readability
- [ ] **Naming** - Variables, functions, classes have meaningful names
- [ ] **Complexity** - Functions not too long (recommended < 20-30 lines)
- [ ] **Comments** - Comments when necessary (explain why, not what)
- [ ] **Consistency** - Consistent with project style

### 3.2 Maintainability
- [ ] **Duplication** - No code duplication (DRY principle)
- [ ] **Complexity** - Cyclomatic complexity not too high
- [ ] **Modularity** - Separates reusable components
- [ ] **Configuration** - Minimal hard-coded values

### 3.3 Coding Standards
- [ ] **Formatting** - Follows language and team conventions
- [ ] **Language features** - Uses new features when appropriate
- [ ] **Best practices** - Follows community standards
- [ ] **Linting rules** - Passes defined linting rules

## 4. Review Process

### 4.1 For Reviewers
**Do:**
- Review within specified timeline (recommended < 24 hours)
- Provide constructive feedback
- Praise good aspects of the code
- Ask questions when unclear
- Review both implementation and tests

**Don't:**
- Micromanage implementation
- Nitpick minor formatting issues
- Force refactoring without important reasons
- Delay unnecessarily

### 4.2 For Authors
**Prepare:**
- Change description (PR description)
- Links to requirements/tickets
- Test instructions (if needed)
- Specify areas needing special feedback

**Respond by:**
- Considering feedback with open mind
- Explaining reasoning when disagreeing
- Requesting clarification when unclear
- Thanking reviewers

## 5. Quick Checklist

### 5.1 Pre-review
- [ ] Code compiles/runs
- [ ] All tests pass
- [ ] CI/CD pipeline passes
- [ ] Documentation updated (if needed)

### 5.2 During Review
- [ ] Security issues
- [ ] Performance problems
- [ ] Major bugs
- [ ] Architectural concerns
- [ ] Test coverage

### 5.3 Before Approval
- [ ] All comments addressed/resolved
- [ ] No regression introduced
- [ ] Meets acceptance criteria
- [ ] Ready for deployment

## 6. Effective Feedback Techniques

### 6.1 Comment Writing Techniques
- **Use questions**: "How would this approach work with case X?"
- **Provide examples**: "Could try this approach: `const result = data.filter(x => x.active)`"
- **Reference standards**: "According to our style guide, we use camelCase for variables"
- **Specify priority level**:
  - MUST (critical): Must fix before merge
  - SHOULD (important): Recommended to fix
  - COULD (minor): Optional improvement

### 6.2 Handling Disagreements
1. **Discuss with data** - Reference facts, metrics, benchmarks
2. **Consult third party** - Seek opinion from senior/tech lead when needed
3. **Consider trade-offs** - Every option has pros and cons
4. **Decide and move on** - Avoid perfectionism paralysis

## 7. Metrics and Improvement

### 7.1 Performance Measurement
- **Cycle time** - Time from PR opening to merge
- **Review depth** - Number of comments and discussions
- **Defect rate** - Bugs found after deployment
- **Team satisfaction** - Satisfaction survey

### 7.2 Process Improvement
- **Regular retrospectives** - Discuss process improvements
- **Pair reviewing** - For complex changes
- **Rotation** - Rotate reviewers regularly
- **Training** - Share best practices

---

## Summary

Good code review is **a conversation between professionals** to create the best possible code. Focus on:
1. **Quality** over speed
2. **Learning** over fault-finding
3. **Collaboration** over competition

#### **Remember:** We're all on the same team. The goal is to build great software, not to prove who is better.
---------------



# Prompt Engineering: โครงสร้างและการเขียน

## 1. ภาษาไทย

### โครงสร้างพื้นฐานของ Prompt ที่ดี
```
บทบาท/บทบาทสมมติ + ภารกิจ + ข้อกำหนดเฉพาะ + รูปแบบผลลัพธ์ + เงื่อนไขเพิ่มเติม
```

### องค์ประกอบสำคัญ
1. **บทบาท (Role)**
   ```
   "คุณเป็นผู้เชี่ยวชาญด้านการตลาดดิจิทัล..."
   "ในฐานะครูสอนวิทยาศาสตร์..."
   ```

2. **ภารกิจ (Task)**
   ```
   "เขียนเนื้อหาโพสต์ Facebook เกี่ยวกับ..."
   "วิเคราะห์ข้อมูลต่อไปนี้และสรุปประเด็นหลัก..."
   ```

3. **บริบท (Context)**
   ```
   "สำหรับธุรกิจร้านกาแฟขนาดเล็ก..."
   "เพื่อใช้สอนนักเรียนชั้นมัธยมศึกษาปีที่ 3..."
   ```

4. **รายละเอียดและข้อกำหนด (Specifications)**
   ```
   "ความยาวประมาณ 300 คำ"
   "ใช้ภาษาที่เป็นทางการ"
   "ระบุแหล่งที่มา 3 แหล่ง"
   ```

5. **รูปแบบผลลัพธ์ (Output Format)**
   ```
   "จัดรูปแบบเป็น bullet points"
   "สรุปในตาราง"
   "เขียนเป็นเรียงความ 5 ย่อหน้า"
   ```

### ตัวอย่าง Prompt ภาษาไทย
```
"ในฐานะนักโภชนาการ กรุณาอธิบายประโยชน์ของอาหารเมดิเตอร์เรเนียนสำหรับผู้สูงอายุ 
โดยเน้นที่ผลต่อสุขภาพหัวใจ ความยาวประมาณ 400 คำ ใช้ภาษาที่เข้าใจง่าย 
และสรุปเป็นข้อๆ 5 ข้อท้ายบทความ"
```

### เคล็ดลับการเขียน
- ระบุความชัดเจนมากกว่าเป็นทั่วไป
- ให้ตัวอย่างถ้าต้องการรูปแบบเฉพาะ
- กำหนดขอบเขตและข้อจำกัด
- ทดลองปรับปรุง prompt หลายครั้ง

## 2. ภาษาอังกฤษ

### Basic Prompt Structure
```
Role + Task + Context + Specifications + Output Format
```

### Key Components
1. **Role Definition**
   ```
   "You are an expert in digital marketing..."
   "As a data scientist with 10 years of experience..."
   ```

2. **Clear Task**
   ```
   "Write a product description for..."
   "Analyze the following dataset and identify trends..."
   ```

3. **Context Provision**
   ```
   "For a startup targeting Gen Z consumers..."
   "In an academic research context..."
   ```

4. **Detailed Specifications**
   ```
   "Use simple language suitable for beginners"
   "Include 5 key takeaways"
   "Limit to 500 words"
   ```

5. **Output Format**
   ```
   "Format as a JSON object"
   "Create a markdown table"
   "Structure as an executive summary"
   ```

### Example English Prompt
```
"As a financial analyst, create an investment risk assessment for renewable energy stocks. 
Consider market volatility, regulatory changes, and technological disruption. 
Present in a structured report with: 1) Executive summary, 2) Risk categories, 
3) Mitigation strategies, 4) Recommendations. Use professional tone and include data points where relevant."
```

### Prompt Writing Techniques
1. **Zero-shot Prompting**
   ```
   "Translate this paragraph to French."
   ```

2. **Few-shot Prompting** (providing examples)
   ```
   "Example 1: [input] -> [output]
    Example 2: [input] -> [output]
    Now process this: [new input]"
   ```

3. **Chain-of-Thought Prompting**
   ```
   "Explain your reasoning step by step."
   "Let's think through this problem systematically."
   ```

### Best Practices
- Be specific and unambiguous
- Use delimiters for complex inputs
- Specify the desired length and depth
- Iterate and refine based on results
- Break complex tasks into subtasks

### Advanced Techniques
- **Temperature setting** (for creativity vs. consistency)
- **System prompts** for setting behavior parameters
- **Template prompts** for reproducible results
- **Meta-prompts** for generating better prompts

ทั้งสองภาษามีหลักการเดียวกัน แต่ต้องคำนึงถึงลักษณะเฉพาะของภาษาและบริบทวัฒนธรรมในการเขียน prompt ที่มีประสิทธิภาพ****

