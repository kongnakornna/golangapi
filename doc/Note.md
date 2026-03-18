
Basic Network
 - TCP vs. UDP
 - TCP ?
 - UDP
 - What is HTTP ?
 -  mqtt
 - snmp
-Monolithic Architecture
Go พื้นฐาน 
- Go Package
- Variables
- Operators
- Control Flow
- Function
- Loop (Debugging included)
- Pointers
- Array and Slice
-  Map
-  Struct
-  Interface
- Generics
- Goroutines
- Channel
- Mutex
- OOP คืออะไร ?
- การอ่าน UML Diagram เบื้องต้น
- Pillars of OOP
- ความสัมพันธ์ ระหว่าง Objects
 -SOLID คือ
- SOLID Principles
- SQL คืออะไร ?
- Relationship
-  ติดตั้ง PostgreSQL บน Docker
- Insert
- Select
- Where
- Like
- And Or
- Order By
- Update
-  Delete
- Join
- Transaction

Git คืออะไร ?
- Git Quick Start
- Git Flow
- API Service ด้วยภาษา Go ในรูปแบบของ Best Practices
Design Sevive
Design ระบบด้วย Domain Driven Design (DDD)
- คำแนะนำ
- DDD คืออะไร ?
- Ubiquitous Language
- Bounded Context และ Context Mapping
- Entity vs. Value Object
- Isekai Shop API Code Repository
- Clean Architecture คืออะไร ?
- การวาง Project Structure
-การแยก Layers
- injection
- value object 

Go Project
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


- ORM 
- CRUD
- core
- auth modules
- health modules
- Logging modules 
- user modules
- dto
- handler
- service
- repository
- model
- router
- platform
Go Features
  Features
 Socket IO
 MQTT
 SNMP

-Package ที่ใช้
-Entities
- Config Setup
-  Package ที่ใช้
-  ตัวอย่าง Config File
-  Config Setup
-  Package ที่ใช้
-  Database Connecting
-  Database Migration โดยใช้ GORM
-  ตัวอย่าง Code
-  Deploying Overview
-  Resources เพิ่มเติม
- Deploy Application ขึ้น GCP Cloud Run  
- การทำ Serve HTTP ด้วย  Echo
- Package ที่ใช้
- Echo Server Initialization
- Gracefully Shutdown
- Middleware Initialization
- ตัวอย่าง Code
- Deploying Overview
-  Resources เพิ่มเติม
-  Deploy Application ขึ้น GCP Cloud Run  
- Item Listing 
- ตั้วอย่าง Code ของ Item Adding Migration
- เพิ่ม Items เข้า Database
- Item Listing
- Implement Custom Error
- Filtering
- Pagination
- Item Managing
- Item Managing Initialization
- Item Creating
- Item Editing
- Item Archiving
- Refactoring
-  ตัวอย่าง Code
- Deploying Overview
- Resources เพิ่มเติม
- Deploy Application ขึ้น GCP Cloud Run [Part 1]
- Deploy Application ขึ้น GCP Cloud Run [Part 2]
- มาลองสร้าง API Service ด้วยภาษา Go ในรูปแบบของ Best Practices
- Google oauth2
- Package ที่ใช้
- OAuth2.0 คืออะไร ?
- Cookie คืออะไร ?
- ทำการสร้าง Player Domain และ Admin Domain
-  GCP
-  ลงมือสร้าง Google OAuth2.0 App
-  OAuth2 Initialization
- Login และ Logout
- Authorization Middleware
- Player Coin 
-  Authorization Middleware
-  Player Coin Initialization
- Coin Adding
- Player's Coin Showing
- Buying and Selling Item in Item shop
- Overview
- Inventory Item Filling
- Inventory Item Removing
- Player's Item Counting
- Inventory Listing
-  Purchase History Recording
- Buying Item
- Selling Item
- สรุปภาพรวมของ Isekaishop Project
-  ตัวอย่าง Code
-  Deploying Overview
-  Deploy Application ขึ้น GCP Cloud Run  
- jwt
- logger
- queue
 - Message
 - NoopQueue
- transaction
- cache คือ
- redis
