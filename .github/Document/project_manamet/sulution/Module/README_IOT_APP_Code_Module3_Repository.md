# Module 3: Repository Layer

## สำหรับโฟลเดอร์ `internal/repository/`

ไฟล์ที่เกี่ยวข้อง:
- `internal/repository/user_repo.go`
- `internal/repository/session_repo.go`
- `internal/repository/redis_repo.go`
- `internal/repository/pg_repository.go` (base transaction manager)

---

## หลักการ (Concept)

### คืออะไร?
Repository Pattern คือตัวกลาง (abstraction) ระหว่าง business logic (usecase) และแหล่งข้อมูล (database, cache, external API) โดยกำหนด interface สำหรับการเข้าถึงข้อมูล และมี implementation ที่เป็นรูปธรรม (PostgreSQL, Redis) แยกออกจากกัน

### มีกี่แบบ?
1. **Specific Repository** – แต่ละ entity มี interface ของตัวเอง (UserRepository, SessionRepository) – ใช้ในโปรเจกต์นี้
2. **Generic Repository** – interface เดียวที่ใช้กับ entity ใดก็ได้ (ใช้ reflection หรือ interface{})
3. **Transaction Repository** – repository ที่มี method สำหรับ transaction (Begin, Commit, Rollback)
4. **Cached Repository** – decorator ที่เพิ่ม cache layer ให้กับ repository หลัก

### ใช้อย่างไร / นำไปใช้กรณีไหน
- ใช้ interface เพื่อกำหนด method ที่ usecase จะเรียก
- implementation จริงใช้ GORM สำหรับ PostgreSQL และ go-redis สำหรับ Redis
- usecase ไม่รู้ว่าข้อมูลมาจาก DB หรือ Cache
- เหมาะกับระบบที่ต้องเปลี่ยนแหล่งข้อมูลบ่อย หรือต้องการ mock สำหรับ unit test

### ทำไมต้องใช้
- แยก business logic ออกจาก data access code
- ทดสอบ usecase ได้ง่ายโดยใช้ mock repository
- สามารถเปลี่ยนจาก PostgreSQL เป็น MongoDB ได้โดยไม่ต้องแก้ usecase

### ประโยชน์ที่ได้รับ
- ลด dependency coupling
- โค้ดสะอาดขึ้น (Clean Architecture)
- รองรับ caching, logging, monitoring ได้โดยไม่แก้ usecase

### ข้อควรระวัง
- repository ควรคืนค่าเป็น model ของ domain (internal/models) ไม่ใช่ DTO
- repository ไม่ควรมี business logic (เช่น การตรวจสอบว่า email ซ้ำ ควรอยู่ใน usecase)
- ระวังเรื่อง transaction: ถ้าต้องการ atomic operation ควรส่ง transaction object (`*gorm.DB`) เข้าไปใน method

### ข้อดี
- ทดสอบง่าย, เปลี่ยนแหล่งข้อมูลได้, แยกความรับผิดชอบชัดเจน

### ข้อเสีย
- มีโค้ดเพิ่มขึ้น (interface + implementation หลายตัว)
- อาจเพิ่มความซับซ้อนในโปรเจกต์เล็ก

### ข้อห้าม
- ห้ามเรียก repository โดยตรงจาก handler (ต้องผ่าน usecase)
- ห้ามใช้ repository ใน repository อื่น (ควรใช้ service หรือ usecase แทน)
- ห้ามใส่ context ลงใน struct repository (ควรส่งผ่าน method argument)

---

## โค้ดที่รันได้จริง

### 1. Base Repository (Transaction Manager) – `pg_repository.go`

```go
// Package repository provides data access layer interfaces and implementations.
// ----------------------------------------------------------------
// แพ็คเกจ repository ให้บริการชั้นการเข้าถึงข้อมูล (อินเทอร์เฟซและอิมพลีเมนต์)
package repository

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
// Automatically commits if no error, otherwise rolls back.
// ----------------------------------------------------------------
// ExecuteInTransaction รันฟังก์ชันที่กำหนดภายใน transaction
// Commit อัตโนมัติถ้าไม่มี error มิฉะนั้น Rollback
func (m *GormTransactionManager) ExecuteInTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx, err := m.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			m.Rollback(tx)
			panic(r)
		}
	}()
	
	if err := fn(tx); err != nil {
		m.Rollback(tx)
		return err
	}
	return m.Commit(tx)
}
```

### 2. User Repository – `user_repo.go`

```go
package repository

import (
	"context"
	"errors"
	"gobackend/internal/models"
	"gorm.io/gorm"
)

// UserRepository defines database operations for users.
// ----------------------------------------------------------------
// UserRepository กำหนดการดำเนินการฐานข้อมูลสำหรับผู้ใช้
type UserRepository interface {
	Create(ctx context.Context, tx *gorm.DB, user *models.User) error
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, tx *gorm.DB, user *models.User) error
	Delete(ctx context.Context, tx *gorm.DB, id uint) error
	List(ctx context.Context, limit, offset int) ([]models.User, int64, error)
	UpdateLastLogin(ctx context.Context, userID uint) error
}

// userRepository implements UserRepository using GORM.
// ----------------------------------------------------------------
// userRepository อิมพลีเมนต์ UserRepository ด้วย GORM
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance.
// ----------------------------------------------------------------
// NewUserRepository สร้าง instance ของ user repository ใหม่
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
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

// Create inserts a new user into the database.
// ----------------------------------------------------------------
// Create เพิ่มผู้ใช้ใหม่ลงฐานข้อมูล
func (r *userRepository) Create(ctx context.Context, tx *gorm.DB, user *models.User) error {
	db := r.getDB(tx)
	return db.WithContext(ctx).Create(user).Error
}

// FindByID retrieves a user by primary key.
// ----------------------------------------------------------------
// FindByID ค้นหาผู้ใช้ด้วยรหัสหลัก
func (r *userRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail retrieves a user by email address (unique).
// ----------------------------------------------------------------
// FindByEmail ค้นหาผู้ใช้ด้วยอีเมล (unique)
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update modifies an existing user.
// ----------------------------------------------------------------
// Update แก้ไขผู้ใช้ที่มีอยู่
func (r *userRepository) Update(ctx context.Context, tx *gorm.DB, user *models.User) error {
	db := r.getDB(tx)
	return db.WithContext(ctx).Save(user).Error
}

// Delete removes a user (soft delete if DeletedAt field exists).
// ----------------------------------------------------------------
// Delete ลบผู้ใช้ (soft delete ถ้ามีฟิลด์ DeletedAt)
func (r *userRepository) Delete(ctx context.Context, tx *gorm.DB, id uint) error {
	db := r.getDB(tx)
	return db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// List returns paginated list of users and total count.
// ----------------------------------------------------------------
// List คืนค่ารายชื่อผู้ใช้แบบแบ่งหน้าและจำนวนทั้งหมด
func (r *userRepository) List(ctx context.Context, limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	
	query := r.db.WithContext(ctx).Model(&models.User{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	err := query.Limit(limit).Offset(offset).Order("id ASC").Find(&users).Error
	return users, total, err
}

// UpdateLastLogin updates the last_login_at field.
// ----------------------------------------------------------------
// UpdateLastLogin อัปเดตฟิลด์ last_login_at
func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("last_login_at", gorm.Expr("NOW()")).Error
}
```

### 3. Session Repository (PostgreSQL + Redis) – `session_repo.go`

```go
package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"gobackend/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// SessionRepository defines operations for user sessions (refresh tokens).
// ----------------------------------------------------------------
// SessionRepository กำหนดการดำเนินการสำหรับ session ของผู้ใช้ (refresh token)
type SessionRepository interface {
	// PostgreSQL operations
	Create(ctx context.Context, session *models.Session) error
	FindByRefreshToken(ctx context.Context, token string) (*models.Session, error)
	Revoke(ctx context.Context, sessionID string) error
	RevokeAllUserSessions(ctx context.Context, userID uint) error
	
	// Redis operations (fast lookup)
	StoreRefreshToken(ctx context.Context, userID uint, token string, ttl time.Duration) error
	GetRefreshToken(ctx context.Context, token string) (uint, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}

// sessionRepository implements SessionRepository using PostgreSQL and Redis.
// ----------------------------------------------------------------
// sessionRepository อิมพลีเมนต์ SessionRepository ด้วย PostgreSQL และ Redis
type sessionRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewSessionRepository creates a new session repository.
// ----------------------------------------------------------------
// NewSessionRepository สร้าง session repository ใหม่
func NewSessionRepository(db *gorm.DB, redisClient *redis.Client) SessionRepository {
	return &sessionRepository{
		db:    db,
		redis: redisClient,
	}
}

// Create saves a session to PostgreSQL.
// ----------------------------------------------------------------
// Create บันทึก session ลง PostgreSQL
func (r *sessionRepository) Create(ctx context.Context, session *models.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// FindByRefreshToken retrieves session from PostgreSQL by refresh token.
// ----------------------------------------------------------------
// FindByRefreshToken ค้นหา session จาก PostgreSQL ด้วย refresh token
func (r *sessionRepository) FindByRefreshToken(ctx context.Context, token string) (*models.Session, error) {
	var session models.Session
	err := r.db.WithContext(ctx).Where("refresh_token = ?", token).First(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// Revoke marks a session as revoked.
// ----------------------------------------------------------------
// Revoke กำหนดสถานะ session เป็นถูกเพิกถอน
func (r *sessionRepository) Revoke(ctx context.Context, sessionID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Session{}).
		Where("id = ?", sessionID).
		Update("revoked_at", now).Error
}

// RevokeAllUserSessions revokes all sessions for a user.
// ----------------------------------------------------------------
// RevokeAllUserSessions เพิกถอน session ทั้งหมดของผู้ใช้
func (r *sessionRepository) RevokeAllUserSessions(ctx context.Context, userID uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Session{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
}

// StoreRefreshToken stores refresh token in Redis with TTL.
// Key format: "refresh:{token}" -> user_id
// ----------------------------------------------------------------
// StoreRefreshToken เก็บ refresh token ใน Redis พร้อม TTL
func (r *sessionRepository) StoreRefreshToken(ctx context.Context, userID uint, token string, ttl time.Duration) error {
	data, _ := json.Marshal(map[string]uint{"user_id": userID})
	return r.redis.Set(ctx, "refresh:"+token, data, ttl).Err()
}

// GetRefreshToken retrieves user ID from Redis by refresh token.
// ----------------------------------------------------------------
// GetRefreshToken ค้นหา user ID จาก Redis ด้วย refresh token
func (r *sessionRepository) GetRefreshToken(ctx context.Context, token string) (uint, error) {
	data, err := r.redis.Get(ctx, "refresh:"+token).Bytes()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	var result map[string]uint
	if err := json.Unmarshal(data, &result); err != nil {
		return 0, err
	}
	return result["user_id"], nil
}

// DeleteRefreshToken removes refresh token from Redis.
// ----------------------------------------------------------------
// DeleteRefreshToken ลบ refresh token ออกจาก Redis
func (r *sessionRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.redis.Del(ctx, "refresh:"+token).Err()
}
```

### 4. Redis Cache Repository – `redis_repo.go`

```go
package repository

import (
	"context"
	"encoding/json"
	"time"
	"github.com/redis/go-redis/v9"
)

// CacheRepository defines operations for generic caching.
// ----------------------------------------------------------------
// CacheRepository กำหนดการดำเนินการสำหรับการแคชทั่วไป
type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Increment(ctx context.Context, key string) (int64, error)
	SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) // for distributed lock
}

// redisCacheRepository implements CacheRepository using go-redis.
// ----------------------------------------------------------------
// redisCacheRepository อิมพลีเมนต์ CacheRepository ด้วย go-redis
type redisCacheRepository struct {
	client *redis.Client
}

// NewCacheRepository creates a new cache repository.
// ----------------------------------------------------------------
// NewCacheRepository สร้าง cache repository ใหม่
func NewCacheRepository(client *redis.Client) CacheRepository {
	return &redisCacheRepository{client: client}
}

// Set stores a value in Redis with TTL.
// ----------------------------------------------------------------
// Set เก็บค่าใน Redis พร้อม TTL
func (r *redisCacheRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

// Get retrieves and unmarshals a value from Redis.
// ----------------------------------------------------------------
// Get ดึงและแปลงค่าจาก Redis
func (r *redisCacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Delete removes a key from Redis.
// ----------------------------------------------------------------
// Delete ลบคีย์ออกจาก Redis
func (r *redisCacheRepository) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in Redis.
// ----------------------------------------------------------------
// Exists ตรวจสอบว่าคีย์มีอยู่ใน Redis หรือไม่
func (r *redisCacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	return n > 0, err
}

// Increment atomically increments a counter.
// ----------------------------------------------------------------
// Increment เพิ่มค่า counter แบบอะตอม
func (r *redisCacheRepository) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// SetNX sets a key only if it does not exist (for distributed locks).
// ----------------------------------------------------------------
// SetNX กำหนดคีย์เฉพาะเมื่อยังไม่มี (ใช้สำหรับ distributed lock)
func (r *redisCacheRepository) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	return r.client.SetNX(ctx, key, data, ttl).Result()
}
```

---

## วิธีใช้งาน module นี้

### ตัวอย่างการนำไปใช้ใน `usecase`

```go
// ใน auth_usecase.go
type AuthUsecase struct {
    userRepo    repository.UserRepository
    sessionRepo repository.SessionRepository
    cacheRepo   repository.CacheRepository
    txManager   repository.TransactionManager
}

func (u *AuthUsecase) Register(ctx context.Context, email, password string) error {
    // ใช้ transaction
    return u.txManager.ExecuteInTransaction(ctx, func(tx *gorm.DB) error {
        user := &models.User{Email: email, PasswordHash: hashed}
        if err := u.userRepo.Create(ctx, tx, user); err != nil {
            return err
        }
        // อาจสร้าง session หรือ verification token ด้วย
        return nil
    })
}

func (u *AuthUsecase) GetUserWithCache(ctx context.Context, id uint) (*models.User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)
    var user models.User
    err := u.cacheRepo.Get(ctx, cacheKey, &user)
    if err == nil {
        return &user, nil
    }
    userPtr, err := u.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    u.cacheRepo.Set(ctx, cacheKey, userPtr, 10*time.Minute)
    return userPtr, nil
}
```

### การ Inject dependencies ใน `main.go`

```go
db := initDB()
redisClient := initRedis()

userRepo := repository.NewUserRepository(db)
sessionRepo := repository.NewSessionRepository(db, redisClient)
cacheRepo := repository.NewCacheRepository(redisClient)
txManager := repository.NewGormTransactionManager(db)

// ส่งต่อไปยัง usecase
authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo, cacheRepo, txManager)
```

---

## ตารางสรุป Repository แต่ละตัว

| Repository | Interface | Implementation | ใช้สำหรับ |
|------------|-----------|----------------|-----------|
| `UserRepository` | CRUD, FindByEmail, List | PostgreSQL (GORM) | จัดการผู้ใช้ |
| `SessionRepository` | Create, Revoke, Redis store | PostgreSQL + Redis | จัดการ refresh token |
| `CacheRepository` | Set, Get, Delete, Exists, Incr | Redis | แคชทั่วไป, distributed lock |
| `TransactionManager` | Begin, Commit, Rollback | GORM | transaction |

---

## แบบฝึกหัดท้าย module (3 ข้อ)

1. เพิ่ม method `FindByRole` ใน `UserRepository` ที่คืนค่ารายชื่อผู้ใช้ตาม role (admin/user) พร้อม pagination
2. Implement `BlacklistToken` ใน `CacheRepository` สำหรับ JWT blacklist (ใช้ SetNX และ TTL เท่ากับเวลาเหลือของ token)
3. สร้าง `SensorRepository` สำหรับบันทึกค่าเซนเซอร์ (temperature, humidity) ลง PostgreSQL พร้อม method `InsertBatch` สำหรับ insert หลาย record ครั้งเดียว (ประสิทธิภาพ)

---

## แหล่งอ้างอิง

- [Repository Pattern in Go](https://dev.to/stevensunflash/using-repository-pattern-in-golang-3h77)
- [GORM documentation](https://gorm.io/docs/)
- [go-redis documentation](https://redis.uptrace.dev/)
- [Clean Architecture with Repository](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1)

---

**หมายเหตุ:** module นี้เป็นส่วนหนึ่งของระบบ gobackend ทั้งหมด หากต้องการ module ถัดไป (Usecase, Delivery, Middleware, ฯลฯ) โปรดแจ้งคำว่า "ต่อไป" หรือระบุชื่อ module ที่ต้องการ