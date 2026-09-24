# Module 15: Database PostgreSQL (GORM & Migrations)

## สำหรับโฟลเดอร์ `internal/repository/`, `migrations/`, และ `cmd/`

ไฟล์ที่เกี่ยวข้อง:
- `internal/repository/pg_repository.go`
- `internal/repository/user_repo.go` (ขยาย)
- `internal/repository/session_repo.go` (ขยาย)
- `migrations/000001_create_users_table.up.sql`
- `migrations/000001_create_users_table.down.sql`
- `cmd/migrate.go`

---

## หลักการ (Concept)

### คืออะไร?
PostgreSQL เป็นฐานข้อมูลเชิงสัมพันธ์ (RDBMS) แบบ open-source ที่มีประสิทธิภาพสูง รองรับการทำงานพร้อมกัน (concurrency) ได้ดี มีฟีเจอร์ที่ทันสมัย เช่น JSONB, full-text search, และ partial indexes เหมาะสำหรับระบบ production ทุกระดับ

### ทำไมต้องใช้ PostgreSQL?
| คุณสมบัติ | PostgreSQL | MySQL |
|-----------|------------|-------|
| **JSON support** | JSONB (binary, indexed) | JSON (text only) |
| **Concurrency** | MVCC, better for high write | InnoDB, good but different |
| **Extensions** | PostGIS, TimescaleDB, etc. | Limited |
| **Standard compliance** | Highest | Moderate |
| **Advanced indexes** | Partial, expression, GIN | Basic B-tree/哈希 |

### GORM คืออะไร?
GORM (Go Object Relational Mapping) คือ ORM (Object-Relational Mapping) ที่ได้รับความนิยมสูงสุดใน Go ช่วยแปลง struct → ตาราง database และ method call → SQL โดยอัตโนมัติ

#### มีกี่แบบ? (Database Connection Methods)

| วิธี | ข้อดี | ข้อเสีย | ใช้เมื่อ |
|-----|------|--------|---------|
| **GORM + standard driver** | ง่าย, รองรับ connection pool | ไม่สามารถควบคุม low-level ได้ | แอปทั่วไป |
| **GORM + pgx** | รองรับ PostgreSQL เฉพาะ, fast | ตั้งค่าซับซ้อนกว่า | ต้องการ performance สูง |
| **Raw database/sql** | ควบคุม SQL 100% | ต้องเขียน SQL เอง, mapping ด้วยมือ | Query ซับซ้อนมาก |
| **sqlc** | type-safe, generate code | ต้องเขียน SQL ก่อน, ไม่อัตโนมัติ | ทีมที่ถนัด SQL |

### ใช้อย่างไร / นำไปใช้กรณีไหน
- ใช้ GORM สำหรับ CRUD ปกติ (users, sessions, sensor logs)
- ใช้ raw SQL หรือ `db.Raw()` สำหรับ query ที่ซับซ้อน (reports, analytics)
- ใช้ AutoMigrate สำหรับ development แต่ไม่แนะนำ production (ควรใช้ migration files)

### ประโยชน์ที่ได้รับ
- **ความสอดคล้องของข้อมูล** (ACID: Atomicity, Consistency, Isolation, Durability)
- **รองรับการทำงานพร้อมกัน** (concurrent access) โดยไม่ corrupt ข้อมูล
- **ความเร็ว** ด้วย indexing, connection pooling, และ query optimization
- **ความปลอดภัย** มีการจัดการ user roles, SSL/TLS, row-level security
- **ขยายแนวนอนได้** (replication, partitioning, read replicas)

### ข้อควรระวัง
- GORM v2 ต้องสร้าง `sql.DB` ก่อนแล้วส่งให้ GORM (ไม่สามารถใช้ direct string ได้)
- AutoMigrate **ไม่ลบ** ฟิลด์หรือเปลี่ยนชนิดข้อมูลเมื่อ struct เปลี่ยน
- `db.Create()` โดยไม่มี transaction จะช้า (GORM ทำ transaction อัตโนมัติ)
- อย่าใช้ `db.Where(&User{Email: ""})` เพราะ GORM จะแปล empty string เป็น `WHERE email = ''`
- ระวังการรัน migration ใน production โดยไม่มีการ backup

### ข้อดี
- ใช้ struct ในการ映射ตาราง → type-safe, ลด human error
- มี built-in connection pooling
- รองรับ associations, preloading, และ hooks

### ข้อเสีย
- ORM overhead สำหรับ query ที่ซับซ้อน
- GORM ไม่ generate efficient SQL ในทุกกรณี
- AutoMigrate ไม่ใช่เครื่องมือจัดการ schema สำหรับ production

### ข้อห้าม
- **ห้าม** เก็บ password ใน plain text (ต้อง hash ด้วย bcrypt)
- **ห้าม** ลืมตั้งค่า connection pool limits (max_open_conns)
- **ห้าม** ใช้ AutoMigrate ใน production ถ้ามีข้อมูลสำคัญ
- **ห้าม** ละเลยการทำ backup ก่อน migrate

---

## การออกแบบ Workflow และ Dataflow

### Workflow: การเชื่อมต่อ PostgreSQL ผ่าน GORM

```mermaid
flowchart TB
    Start[Start Application] --> LoadConfig[Load DB Config from config.yml/env]
    LoadConfig --> CreateSQLDB[sql.Open - create connection pool]
    CreateSQLDB --> GormOpen[gorm.Open - wrap with GORM]
    GormOpen --> SetPoolConfig[SetMaxOpenConns, SetMaxIdleConns, SetConnMaxLifetime]
    SetPoolConfig --> Ping[Ping database]
    Ping -->|success| StoreDB[Store *gorm.DB globally/in DI container]
    Ping -->|failed| Retry{Retry?}
    Retry -->|yes| Wait[Wait 1-5s] --> CreateSQLDB
    Retry -->|no| Exit[Exit with fatal error]
    StoreDB --> Ready[Application ready]
```

**รูปที่ 16:** ขั้นตอนการสร้าง connection ไปยัง PostgreSQL ผ่าน GORM พร้อม connection pool configuration

### Workflow: Database Migration

```mermaid
flowchart TB
    Start[Command: migrate up] --> LoadConfig[Load DB config]
    LoadConfig --> CreateMigrate[golang-migrate instance]
    CreateMigrate --> ReadFiles[Read migration files from migrations/ folder]
    ReadFiles --> CheckVersion[Check current version in schema_migrations table]
    CheckVersion --> ApplyPending[Apply pending .up.sql files in order]
    ApplyPending -->|success| UpdateVersion[Update schema_migrations]
    ApplyPending -->|failure| Rollback[Rollback last migration?]
    Rollback -->|down file exists| ApplyDown[Apply .down.sql]
    UpdateVersion --> Done[Done]
```

**รูปที่ 17:** กระบวนการทำ migration แบบ version-controlled โดยใช้ golang-migrate

### Dataflow: Repository Pattern

```mermaid
flowchart LR
    Handler[HTTP Handler] --> Usecase[Usecase Layer]
    Usecase -->|calls interface| UserRepo[UserRepository interface]
    UserRepo -->|implements| PGUserRepo[PostgreSQL UserRepository]
    PGUserRepo --> GORM[GORM DB]
    GORM --> PostgreSQL[(PostgreSQL)]
```

**รูปที่ 18:** การทำงานของ Repository Pattern ที่แยก interface ออกจาก implementation

---

## ตัวอย่างโค้ดที่รันได้จริง

### 1. การตั้งค่า Connection และ Connection Pool – `internal/repository/pg_repository.go`

```go
// Package repository provides database connection and transaction management.
// ----------------------------------------------------------------
// แพ็คเกจ repository ให้บริการการเชื่อมต่อฐานข้อมูลและการจัดการ transaction
package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"gobackend/internal/config"
	"gobackend/internal/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormDB wraps GORM DB instance with connection pool configuration.
// ----------------------------------------------------------------
// GormDB ห่อหุ้ม GORM DB instance พร้อมการกำหนดค่า connection pool
type GormDB struct {
	DB *gorm.DB
}

// NewPostgresConnection creates and configures PostgreSQL connection with GORM.
// ----------------------------------------------------------------
// NewPostgresConnection สร้างและกำหนดค่า PostgreSQL connection ด้วย GORM
func NewPostgresConnection(cfg *config.DatabaseConfig) (*GormDB, error) {
	// Build DSN (Data Source Name) for PostgreSQL
	// สร้าง DSN สำหรับ PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Bangkok",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	// Configure GORM logger
	// กำหนดค่า GORM logger
	gormLogLevel := gormlogger.Silent
	if cfg.Env == "development" {
		gormLogLevel = gormlogger.Info
	}
	gormLogger := gormlogger.New(
		log.New(logger.LogWriter(), "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormLogLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  cfg.Env == "development",
		},
	)

	// Open connection with GORM
	// เปิด connection ด้วย GORM
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: false,  // Keep default transaction for write safety
		PrepareStmt:            true,   // Enable prepared statements for performance
		TranslateError:         true,   // Translate database errors to GORM errors
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Get underlying sql.DB for connection pool configuration
	// ดึง sql.DB สำหรับการกำหนดค่า connection pool
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Configure connection pool
	// กำหนดค่า connection pool
	// MaxOpenConns: maximum number of open connections (避免超过数据库限制)
	// แนะนำให้ตั้งค่าไม่เกิน 70-80% ของ max_connections ของ PostgreSQL
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	if cfg.MaxOpenConns == 0 {
		sqlDB.SetMaxOpenConns(50) // default, ค่าเริ่มต้น
	}

	// MaxIdleConns: maximum number of idle connections
	// จำนวน connection ที่ idle สูงสุด (ควรน้อยกว่า MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	if cfg.MaxIdleConns == 0 {
		sqlDB.SetMaxIdleConns(25) // default, ค่าเริ่มต้น
	}

	// ConnMaxLifetime: maximum lifetime of a connection (prevents long-lived connections)
	// อายุสูงสุดของ connection (ป้องกัน connection ที่มีอายุยาวเกินไป)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// ConnMaxIdleTime: maximum idle time before closing idle connections
	// เวลาที่ connection idle ก่อนปิด (ลดทรัพยากร)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	// Test connection
	// ทดสอบ connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	logger.Info("PostgreSQL connection established",
		zap.String("host", cfg.Host),
		zap.String("database", cfg.DBName),
		zap.Int("max_open_conns", sqlDB.Stats().MaxOpenConnections),
	)

	return &GormDB{DB: gormDB}, nil
}

// Close gracefully closes the database connection.
// ----------------------------------------------------------------
// Close ปิดการเชื่อมต่อฐานข้อมูลอย่างนุ่มนวล
func (g *GormDB) Close() error {
	sqlDB, err := g.DB.DB()
	if err != nil {
		return err
	}
	logger.Info("Closing PostgreSQL connection")
	return sqlDB.Close()
}

// GetDB returns the underlying GORM DB instance.
// ----------------------------------------------------------------
// GetDB คืน GORM DB instance ด้านล่าง
func (g *GormDB) GetDB() *gorm.DB {
	return g.DB
}

// LoggerWriter adapts logger for GORM.
// ----------------------------------------------------------------
// LoggerWriter ปรับ logger สำหรับ GORM
func (loggerWriter) Write(p []byte) (n int, err error) {
	logger.Debug(string(p))
	return len(p), nil
}
```

### 2. Transaction Manager (Supporting Nested Transactions) – `internal/repository/transaction.go`

```go
package repository

import (
	"context"

	"gorm.io/gorm"
)

// TransactionManager defines methods for managing database transactions.
// ----------------------------------------------------------------
// TransactionManager กำหนด method สำหรับจัดการ transaction ของฐานข้อมูล
type TransactionManager interface {
	// Begin starts a new transaction
	// เริ่ม transaction ใหม่
	Begin(ctx context.Context) (*gorm.DB, error)

	// Commit commits the transaction
	// ยืนยัน transaction
	Commit(tx *gorm.DB) error

	// Rollback aborts the transaction
	// ยกเลิก transaction
	Rollback(tx *gorm.DB) error

	// ExecuteInTransaction runs the given function within a transaction.
	// Automatically commits if no error, otherwise rolls back.
	// Supports nested transactions via GORM savepoints.
	// ----------------------------------------------------------------
	// ExecuteInTransaction รันฟังก์ชันที่กำหนดภายใน transaction
	// Commit อัตโนมัติถ้าไม่มี error มิฉะนั้น Rollback
	// รองรับ nested transaction ผ่าน savepoints ของ GORM
	ExecuteInTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}

// GormTransactionManager implements TransactionManager using GORM.
// ----------------------------------------------------------------
// GormTransactionManager อิมพลีเมนต์ TransactionManager ด้วย GORM
type GormTransactionManager struct {
	db *gorm.DB
}

// NewGormTransactionManager creates a new transaction manager.
// ----------------------------------------------------------------
// NewGormTransactionManager สร้าง transaction manager ใหม่
func NewGormTransactionManager(db *gorm.DB) TransactionManager {
	return &GormTransactionManager{db: db}
}

// Begin starts a new transaction.
// ----------------------------------------------------------------
// Begin เริ่ม transaction ใหม่
func (m *GormTransactionManager) Begin(ctx context.Context) (*gorm.DB, error) {
	tx := m.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

// Commit commits the transaction.
// ----------------------------------------------------------------
// Commit ยืนยัน transaction
func (m *GormTransactionManager) Commit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// Rollback aborts the transaction.
// ----------------------------------------------------------------
// Rollback ยกเลิก transaction
func (m *GormTransactionManager) Rollback(tx *gorm.DB) error {
	return tx.Rollback().Error
}

// ExecuteInTransaction runs the given function within a transaction.
// Supports nested transactions via GORM's savepoint mechanism.
// ----------------------------------------------------------------
// ExecuteInTransaction รันฟังก์ชันที่กำหนดภายใน transaction
// รองรับ nested transaction ผ่าน savepoints ของ GORM
func (m *GormTransactionManager) ExecuteInTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	// GORM's Transaction method automatically creates savepoints for nested calls
	// เมธอด Transaction ของ GORM จะสร้าง savepoints สำหรับ nested calls โดยอัตโนมัติ
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
```

### 3. การใช้ Repository with Transaction (ตัวอย่างใน `user_repo.go`)

```go
// Create inserts a new user into database with optional transaction.
// If tx is provided, uses that transaction; otherwise uses default connection.
// ----------------------------------------------------------------
// Create เพิ่มผู้ใช้ใหม่ลงฐานข้อมูล พร้อม transaction (ถ้ามี)
func (r *userRepository) Create(ctx context.Context, tx *gorm.DB, user *models.User) error {
	db := r.getDB(tx)
	return db.WithContext(ctx).Create(user).Error
}

// getDB returns transaction if provided, otherwise default db.
// ----------------------------------------------------------------
// getDB คืนค่า transaction ถ้ามี หรือ db ปกติ
func (r *userRepository) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}
```

### 4. การทำ AutoMigrate (สำหรับ Development) – `cmd/migrate.go`

```go
// Package main provides CLI commands for database migrations.
// ----------------------------------------------------------------
// แพ็คเกจ main ให้บริการคำสั่ง CLI สำหรับการ migration ฐานข้อมูล
package main

import (
	"flag"
	"log"

	"gobackend/internal/config"
	"gobackend/internal/models"
	"gobackend/internal/repository"
)

func main() {
	// Parse command line arguments
	// แยกอาร์กิวเมนต์บรรทัดคำสั่ง
	cmd := flag.String("cmd", "up", "migration command: up, down, drop")
	flag.Parse()

	// Load configuration
	// โหลดค่ากำหนด
	cfg := config.Load()

	// Connect to database
	// เชื่อมต่อฐานข้อมูล
	dbConn, err := repository.NewPostgresConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	db := dbConn.GetDB()

	switch *cmd {
	case "up":
		// Auto migrate schema (development only)
		// สร้าง/อัปเดต schema (เฉพาะ development)
		log.Println("Running auto migration...")
		if err := db.AutoMigrate(
			&models.User{},
			&models.Session{},
			&models.Verification{},
		); err != nil {
			log.Fatalf("Auto migration failed: %v", err)
		}
		log.Println("Auto migration completed successfully")

	case "drop":
		// Drop all tables (dangerous!)
		// ลบทุกตาราง (อันตราย!)
		log.Println("Dropping all tables...")
		if err := db.Migrator().DropTable(
			&models.User{},
			&models.Session{},
			&models.Verification{},
		); err != nil {
			log.Fatalf("Drop tables failed: %v", err)
		}
		log.Println("All tables dropped")
	}
}
```

### 5. Migration Files (golang-migrate) – `migrations/`

**migrations/000001_create_users_table.up.sql**
```sql
-- Create users table with proper constraints and indexes
-- ----------------------------------------------------------------
-- สร้างตาราง users พร้อม constraints และ indexes ที่เหมาะสม
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for performance
-- สร้าง indexes เพื่อประสิทธิภาพ
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Create trigger to auto-update updated_at
-- สร้าง trigger สำหรับอัปเดต updated_at อัตโนมัติ
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

**migrations/000001_create_users_table.down.sql**
```sql
-- Drop trigger and function
-- ลบ trigger และ function
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop table
-- ลบตาราง
DROP TABLE IF EXISTS users;
```

**migrations/000002_create_sessions_table.up.sql**
```sql
-- Create sessions table for refresh tokens
-- ----------------------------------------------------------------
-- สร้างตาราง sessions สำหรับ refresh tokens
CREATE TABLE IF NOT EXISTS sessions (
    id VARCHAR(36) PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token VARCHAR(255) NOT NULL UNIQUE,
    user_agent VARCHAR(255),
    client_ip VARCHAR(45),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
-- สร้าง indexes
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_refresh_token ON sessions(refresh_token);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
```

**migrations/000002_create_sessions_table.down.sql**
```sql
DROP TABLE IF EXISTS sessions;
```

---

## กรณีศึกษาและแนวทางแก้ไขปัญหา

### ปัญหา: Connection pool หมด (too many connections)

**อาการ:** หลังจาก deploy หลาย instances, PostgreSQL แสดง error `too many clients already`

**สาเหตุ:** แต่ละ instance ของ Go สร้าง connection pool ไปยัง PostgreSQL แบบไม่จำกัด หรือตั้งค่า `MaxOpenConns` สูงเกินไป

**แนวทางแก้ไข:**
1. ตั้งค่า `sqlDB.SetMaxOpenConns(50)` ต่อ instance
2. ใช้ PgBouncer เป็น connection pooler ระหว่าง app กับ PostgreSQL
3. คำนวณ: `max_connections = (max_open_conns_per_instance * number_of_instances) + 10` (สำรองสำหรับ admin tools)

**รหัสแก้ไข:**
```go
// ตั้งค่า connection pool limits
sqlDB.SetMaxOpenConns(50)  // แต่ละ instance ใช้ไม่เกิน 50 connections
sqlDB.SetMaxIdleConns(25)  // idle connections
sqlDB.SetConnMaxLifetime(5 * time.Minute)
```

### ปัญหา: Deadlock เมื่อใช้ transaction

**สาเหตุ:** การ lock rows ในลำดับที่ไม่สอดคล้องกัน (order mismatch) หรือ transaction ที่ใช้เวลานาน

**แนวทางแก้ไข:**
1. เข้าถึง tables ในลำดับเดียวกันทุก transaction
2. ใช้ `FOR UPDATE` เมื่อจำเป็นเท่านั้น
3. ทำให้ transaction สั้นที่สุด
4. ใช้ retry mechanism สำหรับ deadlock

**ตัวอย่างการ retry:**
```go
func (u *userUsecase) UpdateWithRetry(ctx context.Context, userID uint, fn func(*gorm.DB) error) error {
    for attempt := 0; attempt < 3; attempt++ {
        err := u.txManager.ExecuteInTransaction(ctx, fn)
        if err == nil {
            return nil
        }
        if !errors.Is(err, gorm.ErrDeadlock) {
            return err
        }
        time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
    }
    return errors.New("failed after 3 retries")
}
```

### ปัญหา: Migration ล้มเหลวใน production

**แนวทางแก้ไข:**
1. ทำ backup ก่อน migration เสมอ: `pg_dump mydb > backup.sql`
2. ทดสอบ migration บน staging environment ก่อน
3. ใช้ transaction ใน migration (ถ้า DDL รองรับ)
4. วางแผน zero-downtime migration: เพิ่ม column ก่อน, deploy code, แล้วค่อยลบ column

---

## ตารางสรุป PostgreSQL Connection Parameters

| พารามิเตอร์ | คำอธิบาย | ค่าเริ่มต้น | คำแนะนำ |
|-------------|----------|-----------|----------|
| `host` | PostgreSQL server address | localhost | Production ใช้ internal hostname |
| `port` | PostgreSQL port | 5432 | ค่า default |
| `user` | Database user | postgres | สร้าง user เฉพาะสำหรับแอป |
| `password` | User password | - | ใช้ environment variable |
| `dbname` | Database name | - | แยกตาม environment |
| `sslmode` | SSL mode | prefer | Production ใช้ `require` |
| `TimeZone` | Time zone | UTC | ใช้ `Asia/Bangkok` |
| `connect_timeout` | Connection timeout (seconds) | 10 | 5-10 seconds |

---

## วิธีการรันและทดสอบ

### 1. รัน PostgreSQL ด้วย Docker
```bash
docker run --name postgres-test \
  -e POSTGRES_USER=myuser \
  -e POSTGRES_PASSWORD=mypass \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 -d postgres:15-alpine
```

### 2. รัน migration (golang-migrate CLI)
```bash
# Install CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migration up
migrate -source file://migrations -database "postgres://myuser:mypass@localhost:5432/mydb?sslmode=disable" up

# Rollback migration down
migrate -source file://migrations -database "postgres://myuser:mypass@localhost:5432/mydb?sslmode=disable" down 1
```

### 3. ทดสอบ connection ใน Go
```go
func TestConnection(t *testing.T) {
    cfg := config.Load()
    dbConn, err := repository.NewPostgresConnection(&cfg.Database)
    require.NoError(t, err)
    defer dbConn.Close()

    var result int
    err = dbConn.GetDB().Raw("SELECT 1").Scan(&result).Error
    require.NoError(t, err)
    assert.Equal(t, 1, result)
}
```

---

## แบบฝึกหัดท้าย module (3 ข้อ)

1. **เพิ่มฟังก์ชัน `BatchInsert`** ใน repository ที่รับ slice ของ users และใช้ `CreateInBatches` ของ GORM เพื่อ insert หลาย record ครั้งเดียว (ประสิทธิภาพสูงขึ้น)
2. **สร้าง migration** สำหรับตาราง `sensor_logs` (id, device_id, sensor_type, value, unit, location, timestamp) พร้อม indexes ที่เหมาะสม
3. **Implement health check** ที่ ping PostgreSQL และ Redis พร้อม timeout (ใช้ `PingContext`) และ return status ใน `/health/ready`

---

## แหล่งอ้างอิง

- [GORM official documentation](https://gorm.io/docs/)
- [GORM PostgreSQL driver](https://github.com/go-gorm/postgres)
- [golang-migrate/migrate](https://github.com/golang-migrate/migrate)
- [PostgreSQL connection pooling with pgxpool](https://github.com/jackc/pgxpool)
- [PostgreSQL documentation](https://www.postgresql.org/docs/)

---

**หมายเหตุ:** module นี้เป็นส่วนหนึ่งของระบบ gobackend ทั้งหมด หากต้องการ module ถัดไป (เช่น `pkg/validator`, `pkg/logger`, `pkg/jwt`, `pkg/redis`) หรือต้องการปรับปรุงส่วนใดเพิ่มเติม โปรดแจ้ง