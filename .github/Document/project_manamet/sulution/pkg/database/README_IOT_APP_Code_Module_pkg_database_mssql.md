# Module 17: pkg/mssql (Microsoft SQL Server Client & Repository)

## สำหรับโฟลเดอร์ `internal/pkg/mssql/` และ `internal/repository/`

ไฟล์ที่เกี่ยวข้อง:
- `internal/pkg/mssql/client.go`
- `internal/pkg/mssql/repository.go`
- `internal/repository/mssql_user_repo.go`
- `migrations/mssql/` (สำหรับ migration files เฉพาะ SQL Server)

---

## หลักการ (Concept)

### Microsoft SQL Server คืออะไร?

Microsoft SQL Server เป็นระบบจัดการฐานข้อมูลเชิงสัมพันธ์ (RDBMS) แบบ enterprise ที่พัฒนาโดย Microsoft รองรับการทำงานปริมาณสูง มีความปลอดภัยระดับสูง และมีฟีเจอร์เฉพาะ เช่น **Window Functions**, **Common Table Expressions (CTEs)**, **Full-Text Search**, **Always On Availability Groups** สำหรับ high availability โดยสามารถเชื่อมต่อผ่าน GORM ซึ่งมี driver สำหรับ SQL Server โดยเฉพาะ

### มีกี่แบบ? (Use Cases ในระบบ Enterprise)

| การใช้งาน | เหมาะกับ MSSQL? | เหตุผล |
|-----------|-----------------|---------|
| **Enterprise applications** | ✅ เหมาะมาก | รองรับ ACID, transactions, high concurrency |
| **Legacy system integration** | ✅ เหมาะมาก | องค์กรขนาดใหญ่มักมี MSSQL เป็นฐานข้อมูลหลัก |
| **Reporting & analytics** | ✅ เหมาะ | SSRS, SSIS, และ SQL Server Analysis Services |
| **High availability** | ✅ เหมาะ | Always On, failover clustering |
| **Real‑time IoT data** | ⚠️ ปานกลาง | PostgreSQL หรือ MongoDB เหมาะกว่า |
| **Open‑source projects** | ❌ ไม่เหมาะ | มี licensing cost |

### GORM SQL Server Driver

GORM รองรับ SQL Server ผ่าน driver `gorm.io/driver/sqlserver` ซึ่ง底层使用 `github.com/microsoft/go-mssqldb` (TDS protocol) [reference:0]。Driver นี้รองรับทั้ง Windows Authentication และ SQL Server Authentication รวมถึงการเชื่อมต่อผ่าน named instances

### ใช้อย่างไร / นำไปใช้กรณีไหน

1. **Enterprise data storage** – เก็บข้อมูลผู้ใช้, สินค้า, คำสั่งซื้อ ในระบบที่มีความต้องการความปลอดภัยสูง
2. **Legacy system migration** – เชื่อมต่อ Go backend กับ MSSQL เดิมขององค์กร
3. **Reporting** – ใช้ MSSQL View และ Stored Procedures สำหรับสร้างรายงาน
4. **High‑throughput OLTP** – รองรับธุรกรรมปริมาณสูงด้วย locking และ isolation levels ที่หลากหลาย

### ทำไมต้องใช้ (แทน PostgreSQL)

| คุณสมบัติ | PostgreSQL | SQL Server |
|-----------|------------|-------------|
| **License** | Open‑source (free) | Commercial (ต้องจ่าย license) |
| **Windows integration** | Limited | Native (Active Directory, ฯลฯ) |
| **Full‑text search** | Good | Excellent (built‑in) |
| **Reporting tools** | Limited (第三方) | SSRS, SSIS, SSAS |
| **High availability** | 第三方 (repmgr, Patroni) | Always On, clustering |
| **GORM support** | ✅ excellent | ✅ good |

**เมื่อใช้ MSSQL:** องค์กรมี license อยู่แล้ว, ต้องการ integration กับ ecosystem Microsoft, หรือมี legacy database ที่ต้องเชื่อมต่อ

### ประโยชน์ที่ได้รับ

- **Windows Authentication** – ใช้ Active Directory ได้โดยไม่ต้องจัดการ passwords แยก
- **Stored Procedures & Functions** – เรียกใช้ business logic ที่写ใน SQL Server โดยตรง
- **Transaction isolation levels** – `READ UNCOMMITTED` (NOLOCK) สำหรับ reporting queries [reference:1]
- **Linked Servers** – query ข้อมูลข้าม database servers
- **Full‑text indexing** – ค้นหาข้อความแบบรวดเร็ว

### ข้อควรระวัง

- **License cost** – MSSQL Standard edition มีค่าใช้จ่าย ต่อ core
- **Password with special characters** – ถ้ารหัสผ่านมี `@` ต้อง escape หรือใช้ `fmt.Sprintf` แทนการ拼接 string [reference:2]
- **Connection string format** – ใช้ `server=host;user id=user;password=pass;database=db` (ไม่ใช่แบบ URL)
- **Port** – Default port คือ 1433 แต่สามารถระบุใน connection string ได้
- **Encryption** – ต้องระบุ `encrypt=disable` สำหรับ development หรือ `encrypt=true` สำหรับ production พร้อม certificate

### ข้อดี
- Enterprise features, high security, Windows integration, SSRS/SSIS

### ข้อเสีย
- License cost, ไม่ open‑source, resource hungry, community support น้อยกว่า PostgreSQL

### ข้อห้าม
- ห้ามใช้ `SELECT * FROM table (NOLOCK)` โดยไม่เข้าใจผลกระทบ (อาจ读到 dirty data)
- ห้ามใช้ `AutoMigrate` ใน production (ต้องใช้ migration tools)
- ห้ามใช้ Windows Authentication บน Linux container (ต้องใช้ SQL Server Authentication)
- ห้ามละเลยการตั้งค่า connection pool limits

---

## การออกแบบ Workflow และ Dataflow

### Workflow: การเชื่อมต่อ MSSQL ผ่าน GORM

```mermaid
flowchart TB
    Start[Start Application] --> LoadConfig[Load DB Config]
    LoadConfig --> BuildConnString{Build Connection String}
    BuildConnString --> |ODBC style| ConnStrODBC[server=host;user id=user;password=pass;database=db]
    BuildConnString --> |URL style| ConnStrURL[sqlserver://user:pass@host:port?database=db]
    ConnStrODBC --> GormOpen[gorm.Open with sqlserver.Open]
    ConnStrURL --> GormOpen
    GormOpen --> SetPool[SetMaxOpenConns, SetMaxIdleConns, SetConnMaxLifetime]
    SetPool --> Ping[Ping database]
    Ping -->|success| StoreDB[Store *gorm.DB in DI container]
    Ping -->|failed| Retry{Retry?}
    Retry -->|yes| Wait[Wait 1-5s] --> BuildConnString
    Retry -->|no| Exit[Exit with fatal error]
    StoreDB --> Ready[Application ready]
```

**รูปที่ 21:** ขั้นตอนการสร้าง connection ไปยัง SQL Server ผ่าน GORM พร้อม connection pool configuration

### Workflow: Repository Pattern สำหรับ MSSQL

```mermaid
flowchart LR
    Handler[HTTP Handler] --> Usecase[Usecase Layer]
    Usecase -->|calls interface| UserRepo[UserRepository interface]
    UserRepo -->|implements| MSSQLRepo[MSSQL UserRepository]
    MSSQLRepo --> GORM[GORM DB]
    GORM --> MSSQL[(SQL Server)]
```

**รูปที่ 22:** การทำงานของ Repository Pattern ที่แยก interface ออกจาก implementation สำหรับ SQL Server

---

## ตัวอย่างโค้ดที่รันได้จริง

### 1. MSSQL Client – `client.go`

```go
// Package mssql provides Microsoft SQL Server client and utilities using GORM.
// ----------------------------------------------------------------
// แพ็คเกจ mssql ให้บริการ Microsoft SQL Server client และ utilities ด้วย GORM
package mssql

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Config holds MSSQL connection settings.
// ----------------------------------------------------------------
// Config เก็บค่ากำหนดการเชื่อมต่อ MSSQL
type Config struct {
	Server   string // hostname or IP, e.g., "localhost" or "192.168.1.100"
	Port     int    // default 1433
	User     string // SQL Server authentication username
	Password string // SQL Server authentication password
	Database string // database name
	// Optional: Windows Authentication (use "trusted_connection=yes" instead of user/password)
	TrustedConnection bool   // use Windows Authentication (default false)
	Encrypt           bool   // enable encryption (default false for dev, true for prod)
	AppName           string // application name for tracing
	// Connection pool settings
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultConfig returns recommended config for development.
// ----------------------------------------------------------------
// DefaultConfig คืนค่า config ที่แนะนำสำหรับ development
func DefaultConfig() *Config {
	return &Config{
		Server:            "localhost",
		Port:              1433,
		User:              "sa",
		Password:          "YourStrong!Passw0rd",
		Database:          "cmon_db",
		TrustedConnection: false,
		Encrypt:           false,
		MaxOpenConns:      50,
		MaxIdleConns:      10,
		ConnMaxLifetime:   30 * time.Minute,
		ConnMaxIdleTime:   5 * time.Minute,
	}
}

// ConnectionString returns ODBC-style connection string for SQL Server.
// This format handles special characters in password correctly.
// ----------------------------------------------------------------
// ConnectionString คืน connection string แบบ ODBC สำหรับ SQL Server
// รูปแบบนี้จัดการอักขระพิเศษในรหัสผ่านได้อย่างถูกต้อง
func (c *Config) ConnectionString() string {
	if c.TrustedConnection {
		// Windows Authentication (trusted connection)
		// การรับรองตัวตนแบบ Windows
		return fmt.Sprintf(
			"server=%s;port=%d;database=%s;trusted_connection=yes;encrypt=%s",
			c.Server, c.Port, c.Database, c.encryptParam(),
		)
	}
	// SQL Server Authentication
	// การรับรองตัวตนแบบ SQL Server
	return fmt.Sprintf(
		"server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=%s",
		c.Server, c.User, c.Password, c.Port, c.Database, c.encryptParam(),
	)
}

// URLConnectionString returns URL-style connection string.
// May have issues with special characters in password.
// ----------------------------------------------------------------
// URLConnectionString คืน connection string แบบ URL
// อาจมีปัญหากับอักขระพิเศษในรหัสผ่าน
func (c *Config) URLConnectionString() string {
	encrypt := "disable"
	if c.Encrypt {
		encrypt = "true"
	}
	return fmt.Sprintf(
		"sqlserver://%s:%s@%s:%d?database=%s&encrypt=%s",
		c.User, c.Password, c.Server, c.Port, c.Database, encrypt,
	)
}

// encryptParam returns "true" if encryption is enabled, otherwise "disable".
// ----------------------------------------------------------------
// encryptParam คืน "true" ถ้าเปิด encryption, มิฉะนั้น "disable"
func (c *Config) encryptParam() string {
	if c.Encrypt {
		return "true"
	}
	return "disable"
}

// Client wraps GORM DB instance with connection management.
// ----------------------------------------------------------------
// Client ห่อหุ้ม GORM DB instance พร้อมการจัดการการเชื่อมต่อ
type Client struct {
	DB     *gorm.DB
	config *Config
}

// NewClient creates a new MSSQL client with connection pool.
// ----------------------------------------------------------------
// NewClient สร้าง MSSQL client ใหม่พร้อม connection pool
func NewClient(ctx context.Context, cfg *Config, logLevel gormlogger.LogLevel) (*Client, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// Configure GORM logger
	// กำหนดค่า GORM logger
	gormLogger := gormlogger.New(
		log.New(logWriter{}, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Open connection with GORM
	// เปิด connection ด้วย GORM
	dsn := cfg.ConnectionString()
	gormDB, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: false,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQL Server: %w", err)
	}

	// Get underlying sql.DB for connection pool configuration
	// ดึง sql.DB สำหรับการกำหนดค่า connection pool
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Configure connection pool
	// กำหนดค่า connection pool
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}

	// Test connection
	// ทดสอบการเชื่อมต่อ
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping SQL Server: %w", err)
	}

	return &Client{
		DB:     gormDB,
		config: cfg,
	}, nil
}

// Close gracefully closes the database connection.
// ----------------------------------------------------------------
// Close ปิดการเชื่อมต่อฐานข้อมูลอย่างนุ่มนวล
func (c *Client) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// logWriter adapts standard log for GORM.
// ----------------------------------------------------------------
// logWriter ปรับ log มาตรฐานสำหรับ GORM
type logWriter struct{}

func (l logWriter) Write(p []byte) (n int, err error) {
	log.Print(string(p))
	return len(p), nil
}
```

### 2. Generic MSSQL Repository – `repository.go`

```go
package mssql

import (
	"context"

	"gorm.io/gorm"
)

// Repository defines generic CRUD operations for any entity.
// ----------------------------------------------------------------
// Repository กำหนดการดำเนินการ CRUD ทั่วไปสำหรับ entity ใดๆ
type Repository[T any] interface {
	Create(ctx context.Context, tx *gorm.DB, entity *T) error
	FindByID(ctx context.Context, id interface{}) (*T, error)
	Update(ctx context.Context, tx *gorm.DB, entity *T) error
	Delete(ctx context.Context, tx *gorm.DB, id interface{}) error
	List(ctx context.Context, limit, offset int) ([]T, int64, error)
}

// GenericRepository implements Repository with GORM.
// ----------------------------------------------------------------
// GenericRepository อิมพลีเมนต์ Repository ด้วย GORM
type GenericRepository[T any] struct {
	db *gorm.DB
}

// NewGenericRepository creates a new generic repository.
// ----------------------------------------------------------------
// NewGenericRepository สร้าง generic repository ใหม่
func NewGenericRepository[T any](db *gorm.DB) *GenericRepository[T] {
	return &GenericRepository[T]{db: db}
}

// getDB returns transaction if provided, otherwise default db.
// ----------------------------------------------------------------
// getDB คืนค่า transaction ถ้ามี หรือ db ปกติ
func (r *GenericRepository[T]) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

// Create inserts a new entity.
// ----------------------------------------------------------------
// Create เพิ่ม entity ใหม่
func (r *GenericRepository[T]) Create(ctx context.Context, tx *gorm.DB, entity *T) error {
	db := r.getDB(tx)
	return db.WithContext(ctx).Create(entity).Error
}

// FindByID retrieves an entity by primary key.
// ----------------------------------------------------------------
// FindByID ดึง entity ด้วย primary key
func (r *GenericRepository[T]) FindByID(ctx context.Context, id interface{}) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Update modifies an existing entity.
// ----------------------------------------------------------------
// Update แก้ไข entity ที่มีอยู่
func (r *GenericRepository[T]) Update(ctx context.Context, tx *gorm.DB, entity *T) error {
	db := r.getDB(tx)
	return db.WithContext(ctx).Save(entity).Error
}

// Delete removes an entity by primary key.
// ----------------------------------------------------------------
// Delete ลบ entity ด้วย primary key
func (r *GenericRepository[T]) Delete(ctx context.Context, tx *gorm.DB, id interface{}) error {
	db := r.getDB(tx)
	return db.WithContext(ctx).Delete(new(T), id).Error
}

// List returns paginated list of entities.
// ----------------------------------------------------------------
// List คืนค่ารายการ entity แบบแบ่งหน้า
func (r *GenericRepository[T]) List(ctx context.Context, limit, offset int) ([]T, int64, error) {
	var entities []T
	var total int64

	query := r.db.WithContext(ctx).Model(new(T))
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Limit(limit).Offset(offset).Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, total, nil
}
```

### 3. Transaction Manager – `transaction.go`

```go
package mssql

import (
	"context"

	"gorm.io/gorm"
)

// TransactionManager defines methods for managing database transactions.
// ----------------------------------------------------------------
// TransactionManager กำหนด method สำหรับจัดการ transaction ของฐานข้อมูล
type TransactionManager interface {
	Begin(ctx context.Context) (*gorm.DB, error)
	Commit(tx *gorm.DB) error
	Rollback(tx *gorm.DB) error
	ExecuteInTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}

// GormTransactionManager implements TransactionManager using GORM.
// ----------------------------------------------------------------
// GormTransactionManager อิมพลีเมนต์ TransactionManager ด้วย GORM
type GormTransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new transaction manager.
// ----------------------------------------------------------------
// NewTransactionManager สร้าง transaction manager ใหม่
func NewTransactionManager(db *gorm.DB) TransactionManager {
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
// ----------------------------------------------------------------
// ExecuteInTransaction รันฟังก์ชันที่กำหนดภายใน transaction
func (m *GormTransactionManager) ExecuteInTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return m.db.WithContext(ctx).Transaction(fn)
}
```

### 4. User Repository for MSSQL – `internal/repository/mssql_user_repo.go`

```go
// Package repository provides MSSQL-specific implementations.
// ----------------------------------------------------------------
// แพ็คเกจ repository ให้บริการ implementation เฉพาะของ MSSQL
package repository

import (
	"context"
	"errors"

	"gobackend/internal/models"
	"gobackend/internal/pkg/mssql"
	"gorm.io/gorm"
)

// MSSQLUserRepository implements UserRepository for SQL Server.
// ----------------------------------------------------------------
// MSSQLUserRepository อิมพลีเมนต์ UserRepository สำหรับ SQL Server
type MSSQLUserRepository struct {
	db  *gorm.DB
	gen *mssql.GenericRepository[models.User]
}

// NewMSSQLUserRepository creates a new MSSQL user repository.
// ----------------------------------------------------------------
// NewMSSQLUserRepository สร้าง MSSQL user repository ใหม่
func NewMSSQLUserRepository(db *gorm.DB) *MSSQLUserRepository {
	return &MSSQLUserRepository{
		db:  db,
		gen: mssql.NewGenericRepository[models.User](db),
	}
}

// Create inserts a new user.
// ----------------------------------------------------------------
// Create เพิ่มผู้ใช้ใหม่
func (r *MSSQLUserRepository) Create(ctx context.Context, tx *gorm.DB, user *models.User) error {
	return r.gen.Create(ctx, tx, user)
}

// FindByID retrieves a user by ID.
// ----------------------------------------------------------------
// FindByID ดึงผู้ใช้ด้วย ID
func (r *MSSQLUserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	return r.gen.FindByID(ctx, id)
}

// FindByEmail retrieves a user by email using MSSQL-specific query.
// SQL Server uses = instead of ILIKE (case-insensitive depends on collation).
// ----------------------------------------------------------------
// FindByEmail ดึงผู้ใช้ด้วยอีเมล (ใช้ query เฉพาะของ MSSQL)
func (r *MSSQLUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	// SQL Server collation determines case sensitivity
	// For case-insensitive search, use COLLATE SQL_Latin1_General_CP1_CI_AS
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user.
// ----------------------------------------------------------------
// Update อัปเดตผู้ใช้ที่มีอยู่
func (r *MSSQLUserRepository) Update(ctx context.Context, tx *gorm.DB, user *models.User) error {
	return r.gen.Update(ctx, tx, user)
}

// Delete soft-deletes a user (if DeletedAt field exists).
// ----------------------------------------------------------------
// Delete ลบผู้ใช้แบบ soft delete (ถ้ามีฟิลด์ DeletedAt)
func (r *MSSQLUserRepository) Delete(ctx context.Context, tx *gorm.DB, id uint) error {
	return r.gen.Delete(ctx, tx, id)
}

// List returns paginated users with MSSQL-specific pagination.
// SQL Server uses OFFSET/FETCH for pagination.
// ----------------------------------------------------------------
// List คืนค่ารายชื่อผู้ใช้แบบแบ่งหน้าด้วย OFFSET/FETCH (SQL Server style)
func (r *MSSQLUserRepository) List(ctx context.Context, limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Paginated query using OFFSET/FETCH
	err := r.db.WithContext(ctx).
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

// RawSQLExample demonstrates executing raw SQL for complex queries.
// Useful for MSSQL-specific features like Window Functions, CTEs.
// ----------------------------------------------------------------
// RawSQLExample แสดงการ execute raw SQL สำหรับ query ที่ซับซ้อน
// มีประโยชน์สำหรับฟีเจอร์เฉพาะของ MSSQL เช่น Window Functions, CTEs
func (r *MSSQLUserRepository) RawSQLExample(ctx context.Context) ([]models.User, error) {
	var users []models.User
	// Example: Using ROW_NUMBER() for pagination (alternative to OFFSET/FETCH)
	sql := `
		SELECT * FROM (
			SELECT *, ROW_NUMBER() OVER (ORDER BY id) AS row_num
			FROM users
			WHERE is_active = 1
		) AS numbered
		WHERE row_num BETWEEN @start AND @end
	`
	err := r.db.WithContext(ctx).Raw(sql, map[string]interface{}{
		"start": 1,
		"end":   10,
	}).Scan(&users).Error
	return users, err
}
```

### 5. MSSQL Migration Example – `migrations/mssql/000001_create_users_table.up.sql`

```sql
-- Create users table for SQL Server
-- สร้างตาราง users สำหรับ SQL Server
CREATE TABLE users (
    id INT IDENTITY(1,1) PRIMARY KEY,
    email NVARCHAR(255) NOT NULL UNIQUE,
    password_hash NVARCHAR(255) NOT NULL,
    full_name NVARCHAR(255) NULL,
    role NVARCHAR(20) NOT NULL DEFAULT 'user',
    is_active BIT NOT NULL DEFAULT 1,
    last_login_at DATETIME2 NULL,
    created_at DATETIME2 NOT NULL DEFAULT GETUTCDATE(),
    updated_at DATETIME2 NOT NULL DEFAULT GETUTCDATE(),
    deleted_at DATETIME2 NULL
);

-- Create indexes for performance
-- สร้าง indexes เพื่อประสิทธิภาพ
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Create trigger to auto-update updated_at
-- สร้าง trigger สำหรับอัปเดต updated_at อัตโนมัติ
CREATE TRIGGER trg_users_updated_at
ON users
AFTER UPDATE
AS
BEGIN
    SET NOCOUNT ON;
    UPDATE users SET updated_at = GETUTCDATE()
    FROM users u
    INNER JOIN inserted i ON u.id = i.id;
END;
```

**migrations/mssql/000001_create_users_table.down.sql**
```sql
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
DROP TABLE IF EXISTS users;
```

---

## วิธีใช้งาน module นี้

### 1. ติดตั้ง dependencies
```bash
go get gorm.io/gorm
go get gorm.io/driver/sqlserver
```

### 2. วางไฟล์ตามโครงสร้าง
```
internal/pkg/mssql/
├── client.go
├── repository.go
└── transaction.go

internal/repository/
└── mssql_user_repo.go

migrations/mssql/
└── *.up.sql
```

### 3. สร้าง client และ repositories ใน `main.go`

```go
import (
    "gobackend/internal/pkg/mssql"
    "gobackend/internal/repository"
)

func main() {
    // Create MSSQL client
    mssqlCfg := &mssql.Config{
        Server:            os.Getenv("MSSQL_HOST"),
        Port:              1433,
        User:              os.Getenv("MSSQL_USER"),
        Password:          os.Getenv("MSSQL_PASSWORD"),
        Database:          os.Getenv("MSSQL_DATABASE"),
        TrustedConnection: false,
        Encrypt:           true, // enable for production
        MaxOpenConns:      50,
        MaxIdleConns:      10,
    }
    
    mssqlClient, err := mssql.NewClient(context.Background(), mssqlCfg, gormlogger.Info)
    if err != nil {
        log.Fatal(err)
    }
    defer mssqlClient.Close()
    
    // Create repository
    userRepo := repository.NewMSSQLUserRepository(mssqlClient.DB)
    txManager := mssql.NewTransactionManager(mssqlClient.DB)
    
    // Inject into usecases...
}
```

### 4. ตัวอย่างการใช้งาน Usecase

```go
func (u *userUsecase) RegisterWithTransaction(ctx context.Context, email, password string) error {
    return u.txManager.ExecuteInTransaction(ctx, func(tx *gorm.DB) error {
        user := &models.User{
            Email:    email,
            Password: hashedPassword,
        }
        if err := u.userRepo.Create(ctx, tx, user); err != nil {
            return err
        }
        // Additional operations...
        return nil
    })
}
```

---

## ตารางสรุป MSSQL Components

| Component | หน้าที่ | ตัวอย่าง |
|-----------|--------|----------|
| `Client` | จัดการ connection pool | `mssql.NewClient()` |
| `GenericRepository[T]` | Generic CRUD | `Create()`, `FindByID()`, `List()` |
| `TransactionManager` | จัดการ transaction | `ExecuteInTransaction()` |
| `MSSQLUserRepository` | User-specific queries | `FindByEmail()`, `RawSQLExample()` |

---

## แบบฝึกหัดท้าย module (3 ข้อ)

1. เพิ่มฟังก์ชัน `BulkInsert` ใน `GenericRepository` ที่รับ slice ของ entities และใช้ `CreateInBatches` ของ GORM สำหรับ insert จำนวนมาก
2. สร้าง migration สำหรับตาราง `audit_logs` ที่บันทึกการเปลี่ยนแปลง (user_id, action, old_value, new_value, timestamp) พร้อม indexes ที่เหมาะสม
3. Implement repository method ที่ใช้ Common Table Expression (CTE) สำหรับ hierarchical data (เช่น organization chart)

---

## แหล่งอ้างอิง

- [GORM SQL Server Driver](https://github.com/go-gorm/sqlserver) [reference:3]
- [Microsoft go-mssqldb driver](https://github.com/microsoft/go-mssqldb)
- [GORM connection pool configuration](https://gorm.io/docs/generic_interface.html) [reference:4]
- [SQL Server connection strings](https://www.connectionstrings.com/sql-server/)

---

**หมายเหตุ:** module นี้ครบถ้วนสำหรับ `pkg/mssql` สำหรับระบบ gobackend หากต้องการ module เพิ่มเติม (เช่น `pkg/oracle`, `pkg/clickhouse`) โปรดแจ้ง