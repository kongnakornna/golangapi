# golangapi
> Production-ready Go RESTful API boilerplate with Chi, GORM, PostgreSQL, Redis, and enterprise features

## 🌟 Features

- D:\git\goreststarterio\deploy\docker>d
 docker-compose up -d

### 🚀 Core Features
- go install github.com/air-verse/air@latest
- go run cmd/app/main.go

```bash

 go install github.com/air-verse/air@latest
 go run cmd/app/main.go



# Run development server (with auto-reload)
./scripts/dev.sh

# Or run directly
go run cmd/app/main.go
```
- http://localhost:7001/swagger

### Access API Documentation

After starting the service, visit **http://localhost:7001/swagger** to view the interactive API documentation.

## 📚 API Endpoints

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

- [Go Project Layout](https://github.com/kongnakornna/goreststarterio) - Standard Go project structure
- [Chi Router](https://github.com/go-chi/chi) - Lightweight HTTP router
- [GORM](https://gorm.io/) - The fantastic ORM library for Golang


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
http://localhost:7001/api/v1/users/{id}
http://localhost:7001/api/v1/users/{id}
http://localhost:7001/version
http://localhost:7001/status
http://localhost:7001/api/v1/account/logout
http://localhost:7001/api/v1/account/logout
http://localhost:7001/api/v1/account/logout
