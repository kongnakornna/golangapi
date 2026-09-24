# Module 2: Models (Entity Models)

## สำหรับโฟลเดอร์ `internal/models/`

ไฟล์ที่เกี่ยวข้อง:
- `internal/models/base.go`
- `internal/models/user.go`
- `internal/models/session.go`
- `internal/models/verification.go`

---

## หลักการ (Concept)

### คืออะไร?
Models คือโครงสร้างข้อมูล (struct) ที่แทน entity ในฐานข้อมูล หรือข้อมูลที่ใช้ในการสื่อสารระหว่าง layers (DTO) โดยปกติจะสอดคล้องกับตารางใน PostgreSQL และใช้ GORM tags สำหรับ mapping

### มีกี่แบบ?
1. **Entity Model** – สอดคล้องกับตาราง DB โดยตรง (user, session)
2. **DTO (Data Transfer Object)** – ใช้รับ/ส่งข้อมูลระหว่าง API (Request/Response)
3. **Embedded Model** – struct ที่ถูกแทรกใน model อื่น (เช่น BaseModel)
4. **Enum-like Model** – ใช้ iota สำหรับ status constants

### ใช้อย่างไร / นำไปใช้กรณีไหน
- ใช้ GORM annotations (`gorm:"column:name;type:..."`) เพื่อกำหนด schema
- ใช้ JSON tags (`json:"field_name"`) สำหรับ serialization
- ใช้ Validator tags (`validate:"required,email"`) สำหรับ input validation

### ทำไมต้องใช้
- จัดระเบียบโครงสร้างข้อมูลให้เป็นหนึ่งเดียว
- ช่วยให้ GORM สร้างตารางอัตโนมัติ (AutoMigrate)
- แยก business entity ออกจาก database details

### ประโยชน์ที่ได้รับ
- Type safety ใน Go (ไม่ต้องใช้ map[string]interface{})
- ลด boilerplate code สำหรับ CRUD
- รองรับความสัมพันธ์ระหว่างตาราง (Relationships: BelongsTo, HasMany)

### ข้อควรระวัง
- ห้ามเก็บ password plain text (ต้อง hashed)
- ใช้ pointer type สำหรับ nullable fields (`*time.Time` แทน `time.Time`)
- ระวัง zero values (0, "", false) vs null

### ข้อดี
- ชัดเจน, ตรวจสอบได้ตอน compile
- รองรับ GORM hooks (BeforeCreate, AfterUpdate)

### ข้อเสีย
- ต้องเปลี่ยนแปลง struct เมื่อ schema เปลี่ยน
- อาจมีหลาย struct ที่คล้ายกัน (entity vs response DTO)

### ข้อห้าม
- ห้ามใช้ model สำหรับ business logic (ควรอยู่ใน usecase)
- ห้าม serialize model ที่มี password ไปเป็น JSON

---

## โค้ดที่รันได้จริง

### ไฟล์ `internal/models/base.go`

```go
// Package models defines data structures for database entities and DTOs.
// ----------------------------------------------------------------
// แพ็คเกจ models กำหนดโครงสร้างข้อมูลสำหรับ entity ในฐานข้อมูลและ DTO
package models

import (
	"time"
)

// BaseModel provides common fields for all database entities.
// Includes ID, created_at, updated_at, and soft delete support.
// ----------------------------------------------------------------
// BaseModel ให้ฟิลด์ทั่วไปสำหรับ entity ในฐานข้อมูลทั้งหมด
// ประกอบด้วย ID, created_at, updated_at และรองรับการลบแบบ soft delete
type BaseModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"` // soft delete, soft delete
}

// IsDeleted returns true if the record is soft-deleted.
// ----------------------------------------------------------------
// IsDeleted คืนค่า true ถ้าเรกคอร์ดถูกลบแบบ soft delete
func (b *BaseModel) IsDeleted() bool {
	return b.DeletedAt != nil && !b.DeletedAt.IsZero()
}
```

### ไฟล์ `internal/models/user.go`

```go
package models

import (
	"time"
)

// UserRole defines user permission levels.
// ----------------------------------------------------------------
// UserRole กำหนดระดับสิทธิ์ของผู้ใช้
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

// User represents the application user entity.
// Stores authentication credentials and profile information.
// ----------------------------------------------------------------
// User แทน entity ผู้ใช้ของแอปพลิเคชัน
// เก็บข้อมูลรับรองและข้อมูลโปรไฟล์
type User struct {
	BaseModel
	Email        string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash string    `gorm:"not null;size:255" json:"-"` // "-" hides from JSON, ซ่อนจาก JSON
	FullName     string    `gorm:"size:255" json:"full_name"`
	Role         UserRole  `gorm:"type:varchar(20);default:'user'" json:"role"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	
	// Relationships (จะใช้ใน repository joins)
	// ความสัมพันธ์ (ใช้ในการ join ใน repository)
	Sessions     []Session     `gorm:"foreignKey:UserID" json:"-"`
	Verifications []Verification `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies custom table name for GORM.
// ----------------------------------------------------------------
// TableName กำหนดชื่อตารางสำหรับ GORM
func (User) TableName() string {
	return "users"
}

// IsAdmin checks if user has admin role.
// ----------------------------------------------------------------
// IsAdmin ตรวจสอบว่าผู้ใช้มีบทบาท admin หรือไม่
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// BeforeCreate GORM hook: set default role if not set.
// ----------------------------------------------------------------
// BeforeCreate ฟังก์ชันที่ GORM เรียกก่อนสร้างเรกคอร์ด: กำหนด role เริ่มต้น
func (u *User) BeforeCreate() error {
	if u.Role == "" {
		u.Role = RoleUser
	}
	return nil
}
```

### ไฟล์ `internal/models/session.go`

```go
package models

import (
	"time"
)

// Session represents a user refresh token session stored in Redis.
// This is used for JWT refresh token management.
// ----------------------------------------------------------------
// Session แทน session ของ refresh token ที่เก็บใน Redis
// ใช้สำหรับจัดการ JWT refresh token
type Session struct {
	ID           string    `gorm:"primaryKey;size:36" json:"id"` // UUID, UUID
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	RefreshToken string    `gorm:"uniqueIndex;not null;size:255" json:"-"` // hidden, ซ่อน
	UserAgent    string    `gorm:"size:255" json:"user_agent,omitempty"`
	ClientIP     string    `gorm:"size:45" json:"client_ip,omitempty"`
	ExpiresAt    time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty"`
}

// TableName specifies table name for sessions (PostgreSQL).
// ----------------------------------------------------------------
// TableName กำหนดชื่อตาราง sessions
func (Session) TableName() string {
	return "sessions"
}

// IsExpired checks if the session has expired.
// ----------------------------------------------------------------
// IsExpired ตรวจสอบว่า session หมดอายุหรือยัง
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsRevoked checks if the session has been revoked.
// ----------------------------------------------------------------
// IsRevoked ตรวจสอบว่า session ถูกเพิกถอนแล้วหรือไม่
func (s *Session) IsRevoked() bool {
	return s.RevokedAt != nil && !s.RevokedAt.IsZero()
}

// IsValid returns true if session is not expired and not revoked.
// ----------------------------------------------------------------
// IsValid คืนค่า true ถ้า session ยังไม่หมดอายุและไม่ถูกเพิกถอน
func (s *Session) IsValid() bool {
	return !s.IsExpired() && !s.IsRevoked()
}
```

### ไฟล์ `internal/models/verification.go`

```go
package models

import (
	"time"
)

// VerificationType defines the purpose of verification.
// ----------------------------------------------------------------
// VerificationType กำหนดวัตถุประสงค์ของการยืนยัน
type VerificationType string

const (
	VerificationEmail VerificationType = "email_verification"
	VerificationReset VerificationType = "password_reset"
)

// Verification represents email verification or password reset tokens.
// ----------------------------------------------------------------
// Verification แทน token สำหรับยืนยันอีเมลหรือรีเซ็ตรหัสผ่าน
type Verification struct {
	BaseModel
	UserID    uint             `gorm:"not null;index" json:"user_id"`
	Token     string           `gorm:"uniqueIndex;not null;size:255" json:"token"`
	Type      VerificationType `gorm:"type:varchar(50);not null" json:"type"`
	ExpiresAt time.Time        `gorm:"not null" json:"expires_at"`
	UsedAt    *time.Time       `json:"used_at,omitempty"`
	
	// Relationship
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies table name for verifications.
// ----------------------------------------------------------------
// TableName กำหนดชื่อตาราง verifications
func (Verification) TableName() string {
	return "verifications"
}

// IsExpired checks if the verification token has expired.
// ----------------------------------------------------------------
// IsExpired ตรวจสอบว่า token หมดอายุแล้วหรือไม่
func (v *Verification) IsExpired() bool {
	return time.Now().After(v.ExpiresAt)
}

// IsUsed checks if the token has already been used.
// ----------------------------------------------------------------
// IsUsed ตรวจสอบว่า token ถูกใช้ไปแล้วหรือไม่
func (v *Verification) IsUsed() bool {
	return v.UsedAt != nil && !v.UsedAt.IsZero()
}

// IsValid returns true if token is not expired and not used.
// ----------------------------------------------------------------
// IsValid คืนค่า true ถ้า token ยังไม่หมดอายุและยังไม่ถูกใช้
func (v *Verification) IsValid() bool {
	return !v.IsExpired() && !v.IsUsed()
}

// MarkUsed sets the token as used at current time.
// ----------------------------------------------------------------
// MarkUsed กำหนดว่า token ถูกใช้แล้ว ณ เวลาปัจจุบัน
func (v *Verification) MarkUsed() {
	now := time.Now()
	v.UsedAt = &now
}
```

---

## วิธีใช้งาน module นี้

1. วางไฟล์ทั้งหมดใน `internal/models/`
2. ใน `main.go` หรือ `migrate.go` เรียก `db.AutoMigrate(&models.User{}, &models.Session{}, &models.Verification{})`
3. ใช้ structs ใน repository layer:
   ```go
   var user models.User
   db.First(&user, 1)
   if user.IsAdmin() { ... }
   ```
4. ใช้สำหรับรับ JSON request (ถ้าใช้เป็น DTO ก็ได้ แต่ควรแยก DTO ใน `dto/` แทน)

---

## ตารางสรุป Model แต่ละตัว

| Model | ตาราง | ฟิลด์หลัก | ใช้สำหรับ |
|-------|-------|-----------|-----------|
| `User` | users | email, password_hash, role | ผู้ใช้ระบบ |
| `Session` | sessions | refresh_token, user_id, expires_at | จัดการ refresh token |
| `Verification` | verifications | token, type, expires_at | ยืนยันอีเมล / รีเซ็ตพาสเวิร์ด |
| `BaseModel` | (embedded) | id, created_at, updated_at, deleted_at | ฟิลด์พื้นฐานทุกตาราง |

---

## แบบฝึกหัดท้าย module (3 ข้อ)

1. เพิ่มฟิลด์ `PhoneNumber` ใน `User` และเพิ่ม validation tag `validate:"e164"` (สมมติ)
2. สร้าง Model `SensorLog` สำหรับเก็บประวัติเซนเซอร์ (fields: rack_id, sensor_type, value, unit, timestamp) พร้อม GORM tags
3. ปรับ `Session` ให้มี `RefreshTokenHash` แทน `RefreshToken` (เก็บ hash เพื่อความปลอดภัย) และเพิ่ม method `CheckToken(plain string) bool`

---

## แหล่งอ้างอิง

- [GORM Models documentation](https://gorm.io/docs/models.html)
- [GORM conventions](https://gorm.io/docs/conventions.html)
- [Go JSON tags](https://pkg.go.dev/encoding/json)

---

**หมายเหตุ:** module นี้เป็นส่วนหนึ่งของระบบ gobackend ทั้งหมด หากต้องการ module ถัดไป (Repository, Usecase, Delivery, ฯลฯ) โปรดแจ้งคำว่า "ต่อไป" หรือระบุชื่อ module ที่ต้องการ เช่น "Repository"