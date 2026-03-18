# Go-RestAPI

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

- [go project layout](https://"github.com/kongnakornna/golangapigoreststarterio)
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
See `docs/README.md` (Chinese) for architecture, development, maintenance, and deployment guides.

### Prerequisites
- **Go 1.25+** - [Install Go](https://golang.org/doc/install)
- **PostgreSQL 12+** - [Install PostgreSQL](https://postgresql.org/download/)
- **Redis 6+** - [Install Redis](https://redis.io/download)

### Installation

```bash
# Clone the repository
git clone https://github.com/kongnakornna/golangapi.git
cd go-rest-starter

# Install dependencies
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
APP_REDIS_PASSWORD="
APP_REDIS_DB=0

# JWT Configuration
APP_JWT_SECRET=your-secure-secret-key-change-in-production
APP_JWT_ACCESS_TOKEN_EXP=24h
APP_JWT_REFRESH_TOKEN_EXP=168h
APP_JWT_ISSUER=go-rest-starter

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
docker build -t go-rest-starter -f deploy/docker/Dockerfile .

# Run with Docker Compose (includes PostgreSQL and Redis)
cd deploy/docker
docker compose up -d

# Run container only (requires external database and Redis)
docker run -p 7001:7001 \
  -e APP_DB_HOST=your-db-host \
  -e APP_REDIS_HOST=your-redis-host \
  go-rest-starter
```

### Kubernetes Deployment

```bash
# Deploy to Kubernetes
kubectl apply -f deploy/k8s/

# Check deployment status
kubectl get pods -l app=go-rest-starter
kubectl logs -f deployment/go-rest-starter
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
- **Cache**: `github.com/redis/go-redis/v9` - Redis client with pipeline and pub/sub support

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

- [Go Project Layout](https://"github.com/kongnakornna/golangapigoreststarterio) - Standard Go project structure
- [Chi Router](https://github.com/go-chi/chi) - Lightweight HTTP router
- [GORM](https://gorm.io/) - The fantastic ORM library for Golang
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
