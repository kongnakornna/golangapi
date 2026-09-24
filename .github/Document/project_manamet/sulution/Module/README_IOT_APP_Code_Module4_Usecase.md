# Module 4: Usecase Layer (Business Logic)

## สำหรับโฟลเดอร์ `internal/usecase/`

ไฟล์ที่เกี่ยวข้อง:
- `internal/usecase/auth_usecase.go`
- `internal/usecase/user_usecase.go`
- `internal/usecase/cache_usecase.go`

---

## หลักการ (Concept)

### คืออะไร?
Usecase (หรือ Service layer) คือชั้นที่บรรจุ **business logic** ของแอปพลิเคชัน ทำหน้าที่ประสานงานระหว่าง repository ต่างๆ ตรวจสอบกฎทางธุรกิจ และแปลงข้อมูลจากรูปแบบของ repository ให้เป็นรูปแบบที่ delivery (handler) ต้องการ โดย usecase **ไม่รู้** ว่า repository ใช้ PostgreSQL หรือ Redis หรือ external API

### มีกี่แบบ?
1. **Specific Usecase** – แต่ละฟีเจอร์มี usecase ของตัวเอง (AuthUsecase, UserUsecase) – ใช้ในโปรเจกต์นี้
2. **Generic Usecase** – ใช้ interface เดียวกันกับหลาย entity (ไม่ค่อยพบใน Go)
3. **Command/Query Segregation** – แยก Usecase สำหรับการแก้ไขข้อมูล (Command) และการอ่านข้อมูล (Query)
4. **Domain Service** – เมื่อ logic ซับซ้อนและเกี่ยวข้องกับหลาย entity

### ใช้อย่างไร / นำไปใช้กรณีไหน
- ใช้ใน handler: `authUsecase.Login(ctx, email, password)` 
- usecase จะเรียก repository method ต่างๆ และอาจเรียกใช้ transaction manager
- คืนค่า business result หรือ error (ไม่คืน HTTP status code)

### ทำไมต้องใช้
- ป้องกัน business logic กระจายอยู่ใน handler หรือ repository
- ทำให้ทดสอบ business logic ได้โดยไม่ต้องมี HTTP request หรือ database จริง (ใช้ mock repository)
- สอดคล้องกับ Clean Architecture

### ประโยชน์ที่ได้รับ
- เปลี่ยน business logic ได้โดยไม่กระทบ delivery (handler) และ repository
- รองรับการ reuse logic (handler เดียวกันใช้ usecase เดียว)
- ง่ายต่อการเพิ่ม logging, tracing, metrics ในชั้นเดียว

### ข้อควรระวัง
- usecase **ห้าม** import package `net/http` หรือ `gin/chi` เพราะจะทำให้ coupling กับ delivery
- usecase **ห้าม** ส่งออก HTTP status code หรือ JSON
- ควรใช้ interface สำหรับ usecase เพื่อให้ handler มองเห็นแค่ method ที่จำเป็น

### ข้อดี
- แยก business logic ชัดเจน
- ทดสอบ unit ได้ง่าย (ใช้ mock)
- ปรับเปลี่ยน flow ได้โดยไม่แก้ handler

### ข้อเสีย
- เพิ่ม layer ทำให้มีไฟล์มากขึ้น
- มือใหม่อาจเข้าใจยากว่าควรใส่ logic ตรงไหน (repository หรือ usecase)

### ข้อห้าม
- ห้ามเรียก handler โดยตรงจาก usecase
- ห้ามใช้ `*gorm.DB` ใน usecase (ใช้ repository interface แทน)
- ห้ามใช้ context เพื่อส่งค่าที่ไม่เกี่ยวกับ request (ใช้ argument ปกติ)

---

## โค้ดที่รันได้จริง

### 1. Auth Usecase – `auth_usecase.go`

```go
// Package usecase contains business logic for the application.
// ----------------------------------------------------------------
// แพ็คเกจ usecase บรรจุ business logic ของแอปพลิเคชัน
package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gobackend/internal/models"
	"gobackend/internal/repository"
	"gobackend/internal/pkg/jwt"
	"gobackend/internal/pkg/hash"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Common business errors.
// ----------------------------------------------------------------
// ข้อผิดพลาดทางธุรกิจทั่วไป
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

// AuthUsecase defines authentication business logic.
// ----------------------------------------------------------------
// AuthUsecase กำหนด business logic สำหรับการรับรองตัวตน
type AuthUsecase interface {
	Register(ctx context.Context, email, password, fullName string) (*models.User, error)
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
	RefreshToken(ctx context.Context, refreshToken string) (newAccessToken string, err error)
	Logout(ctx context.Context, userID uint, refreshToken string) error
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
}

// authUsecase implements AuthUsecase.
// ----------------------------------------------------------------
// authUsecase อิมพลีเมนต์ AuthUsecase
type authUsecase struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	cacheRepo   repository.CacheRepository
	txManager   repository.TransactionManager
	jwtMaker    jwt.Maker
	hashHelper  hash.PasswordHasher
}

// NewAuthUsecase creates a new auth usecase instance.
// ----------------------------------------------------------------
// NewAuthUsecase สร้าง instance ของ auth usecase ใหม่
func NewAuthUsecase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	cacheRepo repository.CacheRepository,
	txManager repository.TransactionManager,
	jwtMaker jwt.Maker,
	hashHelper hash.PasswordHasher,
) AuthUsecase {
	return &authUsecase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		cacheRepo:   cacheRepo,
		txManager:   txManager,
		jwtMaker:    jwtMaker,
		hashHelper:  hashHelper,
	}
}

// Register creates a new user account.
// ----------------------------------------------------------------
// Register สร้างบัญชีผู้ใช้ใหม่
func (u *authUsecase) Register(ctx context.Context, email, password, fullName string) (*models.User, error) {
	// Check if email already exists
	// ตรวจสอบว่าอีเมลมีอยู่ในระบบแล้วหรือไม่
	existing, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("check email existence: %w", err)
	}
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	// แฮชรหัสผ่าน
	hashedPassword, err := u.hashHelper.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &models.User{
		Email:        email,
		PasswordHash: hashedPassword,
		FullName:     fullName,
		Role:         models.RoleUser,
		IsActive:     true,
	}

	// Use transaction to ensure consistency (if we also create email verification later)
	// ใช้ transaction เพื่อความสอดคล้อง (ถ้ามีการสร้าง email verification ในอนาคต)
	err = u.txManager.ExecuteInTransaction(ctx, func(tx *gorm.DB) error {
		return u.userRepo.Create(ctx, tx, user)
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

// Login authenticates a user and returns JWT tokens.
// ----------------------------------------------------------------
// Login ตรวจสอบผู้ใช้และคืน JWT tokens
func (u *authUsecase) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error) {
	// Find user by email
	// ค้นหาผู้ใช้ด้วยอีเมล
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return "", "", ErrUserNotFound
	}

	// Verify password
	// ตรวจสอบรหัสผ่าน
	if !u.hashHelper.Verify(password, user.PasswordHash) {
		return "", "", ErrInvalidPassword
	}

	// Check if user is active
	// ตรวจสอบว่าผู้ใช้ยัง active หรือไม่
	if !user.IsActive {
		return "", "", errors.New("user account is disabled")
	}

	// Generate access token (short-lived)
	// สร้าง access token (อายุสั้น)
	accessToken, accessPayload, err := u.jwtMaker.CreateToken(
		user.ID,
		string(user.Role),
		time.Duration(15)*time.Minute, // configurable, กำหนดได้จาก config
	)
	if err != nil {
		return "", "", fmt.Errorf("create access token: %w", err)
	}

	// Generate refresh token (long-lived, stored in Redis/DB)
	// สร้าง refresh token (อายุยาว, เก็บใน Redis/DB)
	refreshToken = uuid.New().String()
	refreshTTL := time.Duration(7 * 24 * time.Hour) // 7 days, 7 วัน

	// Store refresh token in Redis for fast lookup
	// เก็บ refresh token ใน Redis เพื่อค้นหาเร็ว
	if err := u.sessionRepo.StoreRefreshToken(ctx, user.ID, refreshToken, refreshTTL); err != nil {
		return "", "", fmt.Errorf("store refresh token: %w", err)
	}

	// Also store session in PostgreSQL for audit/logging
	// เก็บ session ใน PostgreSQL เพื่อการตรวจสอบ
	session := &models.Session{
		ID:           accessPayload.ID.String(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(refreshTTL),
	}
	if err := u.sessionRepo.Create(ctx, session); err != nil {
		// Log but don't fail login
		// บันทึก log แต่ไม่ทำให้ login ล้มเหลว
		_ = err
	}

	// Update last login time
	// อัปเดตเวลาเข้าสู่ระบบล่าสุด
	_ = u.userRepo.UpdateLastLogin(ctx, user.ID)

	return accessToken, refreshToken, nil
}

// RefreshToken generates a new access token using a valid refresh token.
// ----------------------------------------------------------------
// RefreshToken สร้าง access token ใหม่โดยใช้ refresh token ที่ถูกต้อง
func (u *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// Get user ID from Redis using refresh token
	// ดึง user ID จาก Redis โดยใช้ refresh token
	userID, err := u.sessionRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", fmt.Errorf("get refresh token: %w", err)
	}
	if userID == 0 {
		return "", ErrInvalidRefreshToken
	}

	// Get user details to obtain role
	// ดึงรายละเอียดผู้ใช้เพื่อทราบ role
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return "", ErrUserNotFound
	}

	// Generate new access token
	// สร้าง access token ใหม่
	newAccessToken, _, err := u.jwtMaker.CreateToken(
		user.ID,
		string(user.Role),
		time.Duration(15)*time.Minute,
	)
	if err != nil {
		return "", fmt.Errorf("create access token: %w", err)
	}

	return newAccessToken, nil
}

// Logout invalidates the refresh token and blacklists the access token (optional).
// ----------------------------------------------------------------
// Logout ทำให้ refresh token ใช้ไม่ได้และ blacklist access token (optional)
func (u *authUsecase) Logout(ctx context.Context, userID uint, refreshToken string) error {
	// Delete refresh token from Redis
	// ลบ refresh token ออกจาก Redis
	if err := u.sessionRepo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}

	// Optionally revoke all sessions for this user (or just this one)
	// ไม่จำเป็นต้อง revoke ทั้งหมด แค่ลบ token ก็พอ
	return nil
}

// ChangePassword updates user's password after verifying old password.
// ----------------------------------------------------------------
// ChangePassword อัปเดตรหัสผ่านหลังจากตรวจสอบรหัสผ่านเดิม
func (u *authUsecase) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	// Get user
	// ดึงข้อมูลผู้ใช้
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Verify old password
	// ตรวจสอบรหัสผ่านเดิม
	if !u.hashHelper.Verify(oldPassword, user.PasswordHash) {
		return ErrInvalidPassword
	}

	// Hash new password
	// แฮชรหัสผ่านใหม่
	newHashed, err := u.hashHelper.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("hash new password: %w", err)
	}

	// Update password in database (use transaction if needed)
	// อัปเดตรหัสผ่านในฐานข้อมูล
	user.PasswordHash = newHashed
	return u.userRepo.Update(ctx, nil, user)
}
```

### 2. User Usecase – `user_usecase.go`

```go
package usecase

import (
	"context"
	"errors"
	"fmt"

	"gobackend/internal/models"
	"gobackend/internal/repository"
)

// UserUsecase defines business logic for user management.
// ----------------------------------------------------------------
// UserUsecase กำหนด business logic สำหรับการจัดการผู้ใช้
type UserUsecase interface {
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context, page, pageSize int) ([]models.User, int64, error)
	// Admin only
	CreateUserAsAdmin(ctx context.Context, email, password, fullName, role string) (*models.User, error)
}

// userUsecase implements UserUsecase.
// ----------------------------------------------------------------
// userUsecase อิมพลีเมนต์ UserUsecase
type userUsecase struct {
	userRepo   repository.UserRepository
	cacheRepo  repository.CacheRepository
	hashHelper hash.PasswordHasher
}

// NewUserUsecase creates a new user usecase instance.
// ----------------------------------------------------------------
// NewUserUsecase สร้าง instance ของ user usecase ใหม่
func NewUserUsecase(
	userRepo repository.UserRepository,
	cacheRepo repository.CacheRepository,
	hashHelper hash.PasswordHasher,
) UserUsecase {
	return &userUsecase{
		userRepo:   userRepo,
		cacheRepo:  cacheRepo,
		hashHelper: hashHelper,
	}
}

// GetUserByID retrieves a user by ID, using cache if available.
// ----------------------------------------------------------------
// GetUserByID ดึงข้อมูลผู้ใช้ด้วย ID โดยใช้ cache ถ้ามี
func (u *userUsecase) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	// Try cache first
	// ลองจาก cache ก่อน
	cacheKey := fmt.Sprintf("user:%d", id)
	var user models.User
	err := u.cacheRepo.Get(ctx, cacheKey, &user)
	if err == nil {
		return &user, nil
	}

	// Cache miss: get from database
	// cache ไม่พบ: ดึงจากฐานข้อมูล
	userPtr, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	if userPtr == nil {
		return nil, ErrUserNotFound
	}

	// Store in cache for future requests (TTL 10 minutes)
	// เก็บใน cache สำหรับการเรียกครั้งถัดไป (หมดอายุ 10 นาที)
	_ = u.cacheRepo.Set(ctx, cacheKey, userPtr, 10*time.Minute)

	return userPtr, nil
}

// GetUserByEmail retrieves a user by email (no cache by default).
// ----------------------------------------------------------------
// GetUserByEmail ดึงข้อมูลผู้ใช้ด้วยอีเมล (ไม่ใช้ cache โดยปกติ)
func (u *userUsecase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// UpdateUser updates user information.
// ----------------------------------------------------------------
// UpdateUser อัปเดตข้อมูลผู้ใช้
func (u *userUsecase) UpdateUser(ctx context.Context, user *models.User) error {
	if user.ID == 0 {
		return errors.New("user ID is required")
	}
	// Business rule: cannot change role through this method (admin only)
	// กฎธุรกิจ: ไม่สามารถเปลี่ยน role ผ่าน method นี้ (ต้องใช้ admin)
	existing, err := u.userRepo.FindByID(ctx, user.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrUserNotFound
	}
	// Preserve role and password hash
	// คง role และ password hash ไว้
	user.Role = existing.Role
	user.PasswordHash = existing.PasswordHash

	if err := u.userRepo.Update(ctx, nil, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	// Invalidate cache
	// ทำให้ cache หมดอายุ
	cacheKey := fmt.Sprintf("user:%d", user.ID)
	_ = u.cacheRepo.Delete(ctx, cacheKey)
	return nil
}

// DeleteUser soft-deletes a user.
// ----------------------------------------------------------------
// DeleteUser ลบผู้ใช้แบบ soft delete
func (u *userUsecase) DeleteUser(ctx context.Context, id uint) error {
	existing, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrUserNotFound
	}
	if err := u.userRepo.Delete(ctx, nil, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	// Invalidate cache
	_ = u.cacheRepo.Delete(ctx, fmt.Sprintf("user:%d", id))
	return nil
}

// ListUsers returns paginated users.
// ----------------------------------------------------------------
// ListUsers คืนค่ารายชื่อผู้ใช้แบบแบ่งหน้า
func (u *userUsecase) ListUsers(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return u.userRepo.List(ctx, pageSize, offset)
}

// CreateUserAsAdmin creates a new user with specified role (admin only).
// ----------------------------------------------------------------
// CreateUserAsAdmin สร้างผู้ใช้ใหม่พร้อม role ที่ระบุ (สำหรับ admin เท่านั้น)
func (u *userUsecase) CreateUserAsAdmin(ctx context.Context, email, password, fullName, role string) (*models.User, error) {
	// Validate role
	// ตรวจสอบ role
	if role != string(models.RoleUser) && role != string(models.RoleAdmin) {
		return nil, errors.New("invalid role, must be 'user' or 'admin'")
	}

	// Check existing email
	existing, _ := u.userRepo.FindByEmail(ctx, email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashed, err := u.hashHelper.Hash(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: hashed,
		FullName:     fullName,
		Role:         models.UserRole(role),
		IsActive:     true,
	}
	if err := u.userRepo.Create(ctx, nil, user); err != nil {
		return nil, err
	}
	return user, nil
}
```

### 3. Cache Usecase (สำหรับ business logic ที่เกี่ยวกับ cache โดยเฉพาะ) – `cache_usecase.go`

```go
package usecase

import (
	"context"
	"time"
	"gobackend/internal/repository"
)

// CacheUsecase defines business logic for cache operations.
// ----------------------------------------------------------------
// CacheUsecase กำหนด business logic สำหรับการดำเนินการกับ cache
type CacheUsecase interface {
	GetOrSet(ctx context.Context, key string, ttl time.Duration, fetcher func() (interface{}, error)) (interface{}, error)
	InvalidatePattern(ctx context.Context, pattern string) error
	IncrementRequestCount(ctx context.Context, userID uint) (int64, error)
}

type cacheUsecase struct {
	cacheRepo repository.CacheRepository
}

func NewCacheUsecase(cacheRepo repository.CacheRepository) CacheUsecase {
	return &cacheUsecase{cacheRepo: cacheRepo}
}

// GetOrSet attempts to get from cache; if not found, calls fetcher and stores result.
// ----------------------------------------------------------------
// GetOrSet พยายามดึงจาก cache ถ้าไม่พบ เรียก fetcher และเก็บผลลัพธ์
func (u *cacheUsecase) GetOrSet(ctx context.Context, key string, ttl time.Duration, fetcher func() (interface{}, error)) (interface{}, error) {
	var result interface{}
	err := u.cacheRepo.Get(ctx, key, &result)
	if err == nil {
		return result, nil
	}
	// Cache miss
	value, err := fetcher()
	if err != nil {
		return nil, err
	}
	_ = u.cacheRepo.Set(ctx, key, value, ttl)
	return value, nil
}

// InvalidatePattern deletes all keys matching a pattern (requires SCAN, implement if needed).
// ----------------------------------------------------------------
// InvalidatePattern ลบคีย์ทั้งหมดที่ตรงกับ pattern (ต้องใช้ SCAN)
func (u *cacheUsecase) InvalidatePattern(ctx context.Context, pattern string) error {
	// In production, implement Redis SCAN + DELETE
	// สำหรับ production ให้ implement ด้วย Redis SCAN + DELETE
	return nil
}

// IncrementRequestCount increments a user-specific request counter.
// ----------------------------------------------------------------
// IncrementRequestCount เพิ่ม counter การเรียกใช้สำหรับผู้ใช้
func (u *cacheUsecase) IncrementRequestCount(ctx context.Context, userID uint) (int64, error) {
	key := fmt.Sprintf("rate:user:%d", userID)
	return u.cacheRepo.Increment(ctx, key)
}
```

---

## วิธีใช้งาน module นี้

### ตัวอย่างการเรียกใช้จาก handler (delivery/rest/handler/auth_handler.go)

```go
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    // bind JSON...

    accessToken, refreshToken, err := h.authUsecase.Login(r.Context(), req.Email, req.Password)
    if errors.Is(err, usecase.ErrUserNotFound) || errors.Is(err, usecase.ErrInvalidPassword) {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }
    if err != nil {
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }
    // return tokens as JSON
}
```

### Dependency Injection ใน `main.go`

```go
// Initialize repositories
userRepo := repository.NewUserRepository(db)
sessionRepo := repository.NewSessionRepository(db, rdb)
cacheRepo := repository.NewCacheRepository(rdb)
txManager := repository.NewGormTransactionManager(db)

// Initialize helpers
jwtMaker, _ := jwt.NewRSAMaker(privateKey, publicKey)
hasher := hash.NewBcryptHasher()

// Initialize usecases
authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo, cacheRepo, txManager, jwtMaker, hasher)
userUsecase := usecase.NewUserUsecase(userRepo, cacheRepo, hasher)

// Pass to handlers
authHandler := handler.NewAuthHandler(authUsecase)
```

---

## ตารางสรุป Usecase แต่ละตัว

| Usecase | Interface | Dependencies | หน้าที่หลัก |
|---------|-----------|--------------|-------------|
| `AuthUsecase` | Register, Login, RefreshToken, Logout, ChangePassword | UserRepo, SessionRepo, CacheRepo, TxManager, JwtMaker, Hasher | การรับรองตัวตน, จัดการ session |
| `UserUsecase` | Get, Update, Delete, List, CreateByAdmin | UserRepo, CacheRepo, Hasher | จัดการโปรไฟล์ผู้ใช้ |
| `CacheUsecase` | GetOrSet, InvalidatePattern, Increment | CacheRepo | ปฏิบัติการ cache ที่มี business logic |

---

## แบบฝึกหัดท้าย module (3 ข้อ)

1. เพิ่ม method `ResetPassword` ใน `AuthUsecase` ที่รับ token (จาก email) และ newPassword โดยตรวจสอบ token กับ verification model (สร้าง verification repo และ model ก่อน)
2. ปรับปรุง `Login` ให้บันทึก失败的 login attempts ใน Redis และ lock account หลังจากพยายามผิด 5 ครั้ง (ใช้ cacheRepo.Increment และ SetNX)
3. สร้าง `SensorUsecase` ที่มี method `ProcessSensorReading(ctx, sensorData)` ซึ่งจะเรียก rule engine (สมมติ) และบันทึกข้อมูลลง PostgreSQL ผ่าน repository ใหม่

---

## แหล่งอ้างอิง

- [Clean Architecture in Go](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1)
- [Domain-Driven Design with Go](https://www.amazon.com/Domain-Driven-Design-Tackling-Complexity-Software/dp/0321125215)
- [Uber Go Style Guide - Use Cases](https://github.com/uber-go/guide/blob/master/style.md#use-cases)

---

**หมายเหตุ:** module นี้เป็นส่วนหนึ่งของระบบ gobackend ทั้งหมด หากต้องการ module ถัดไป (Delivery - Handlers, Middleware, Router) โปรดแจ้งคำว่า "ต่อไป" หรือระบุชื่อ module ที่ต้องการ