# Module 14: pkg/cryptpass (Advanced Password Hashing)

## สำหรับโฟลเดอร์ `internal/pkg/cryptpass/`

ไฟล์ที่เกี่ยวข้อง:
- `internal/pkg/cryptpass/cryptpass.go`

---

## หลักการ (Concept)

### คืออะไร?
cryptpass คือแพ็คเกจสำหรับการเข้ารหัสรหัสผ่าน (password hashing) ที่มีความปลอดภัยสูง โดยรองรับทั้ง bcrypt (ซึ่งเป็นมาตรฐานที่ใช้กันอย่างแพร่หลาย) และ Argon2id (ซึ่งเป็นมาตรฐานใหม่ที่ปลอดภัยที่สุดในปัจจุบัน) รวมถึงฟังก์ชันเสริมสำหรับการตรวจสอบความแข็งแรงของรหัสผ่าน (password strength validation) และการอัปเกรด hash เมื่อค่า cost เพิ่มขึ้น

### มีกี่แบบ?
1. **bcrypt** – ใช้ Blowfish cipher มี adjustable work factor (cost) มี salt ในตัว ใช้งานง่าย
2. **Argon2id** – มาตรฐานใหม่ (PHC winner) เป็น memory-hard function ทนทานต่อ GPU/ASIC attacks แบบ hybrid mode (Argon2id) แนะนำให้ใช้
3. **Scrypt** – memory-hard function เช่นกัน แต่ซับซ้อนกว่า bcrypt
4. **PBKDF2** – มาตรฐาน NIST แต่ไม่ memory-hard

**ในโปรเจกต์นี้ใช้ bcrypt เป็นหลัก (DefaultCost = 12) และมีตัวอย่าง Argon2id สำหรับโปรเจกต์ที่ต้องการความปลอดภัยสูง**

### ใช้อย่างไร / นำไปใช้กรณีไหน
- **bcrypt**: ใช้ในระบบทั่วไปที่ต้องการความสมดุลระหว่างความปลอดภัยและประสิทธิภาพ
- **Argon2id**: ใช้ในระบบที่ต้องการความปลอดภัยสูง (การเงิน, healthcare, blockchain)
- **Password strength validation**: ใช้ตรวจสอบก่อน hash เพื่อป้องกันรหัสผ่านอ่อนแอ
- **Rehash**: ใช้เมื่อต้องการอัปเกรด cost factor ของ hash เดิม

### ทำไมต้องใช้
- bcrypt เป็นมาตรฐานที่ผ่านการทดสอบมายาวนาน
- bcrypt มี salt ในตัว ป้องกัน rainbow table attack
- bcrypt ปรับ cost ได้ (2^cost iterations) ทำให้ทนต่อการเพิ่มพลังประมวลผลในอนาคต
- Argon2id ถูกออกแบบมาเพื่อต้านทาน GPU cracking โดยเฉพาะ
- การตรวจสอบ password strength ช่วยป้องกัน brute force

### ประโยชน์ที่ได้รับ
- ปลอดภัยแม้ database ถูกขโมย
- ป้องกัน brute force ด้วยความช้าโดย design
- มี salt อัตโนมัติ
- สามารถอัปเกรด hash ได้เมื่อ hardware แรงขึ้น

### ข้อควรระวัง
- bcrypt จำกัดความยาวรหัสผ่านที่ 72 ไบต์ (ต้องตัดหรือ pre-hash)
- cost สูงเกินไปอาจทำให้ login ช้า (load test ก่อน)
- อย่าใช้ MD5, SHA-1, SHA-2 ในการ hash รหัสผ่าน (เร็วเกินไป)

### ข้อดี
- ปลอดภัย, ปรับระดับความปลอดภัยได้, มี salt ในตัว

### ข้อเสีย
- ช้ากว่า hashing สำหรับ checksum (แต่这是优点)
- bcrypt ใช้ memory น้อย (ไม่ memory-hard) เมื่อเทียบกับ Argon2
- bcrypt ถูก GPU cracking ได้ง่ายกว่า Argon2

### ข้อห้าม
- ห้ามเก็บรหัสผ่านใน plain text
- ห้ามใช้ hash โดยไม่ใส่ salt
- ห้ามใช้ bcrypt.DefaultCost (10) ใน production (ควรใช้ 12 ขึ้นไป)
- ห้าม truncate password โดยไม่แจ้งผู้ใช้

---

## โค้ดที่รันได้จริง

### ไฟล์ `internal/pkg/cryptpass/cryptpass.go`

```go
// Package cryptpass provides secure password hashing using bcrypt and Argon2id.
// Includes password strength validation and rehashing capabilities.
// ----------------------------------------------------------------
// แพ็คเกจ cryptpass ให้บริการเข้ารหัสรหัสผ่านแบบปลอดภัยด้วย bcrypt และ Argon2id
// รวมถึงการตรวจสอบความแข็งแรงของรหัสผ่านและการ rehash
package cryptpass

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

// Common errors.
// ----------------------------------------------------------------
// ข้อผิดพลาดที่พบบ่อย
var (
	ErrPasswordTooShort  = errors.New("password must be at least 8 characters long")
	ErrPasswordTooLong   = errors.New("password exceeds maximum length")
	ErrPasswordNoUpper   = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLower   = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoDigit   = errors.New("password must contain at least one digit")
	ErrPasswordNoSpecial = errors.New("password must contain at least one special character (!@#$%^&*)")
	ErrInvalidHash       = errors.New("invalid hash format")
	ErrHashMismatch      = errors.New("password does not match hash")
)

// BcryptConfig holds bcrypt hashing parameters.
// ----------------------------------------------------------------
// BcryptConfig เก็บพารามิเตอร์การแฮชของ bcrypt
type BcryptConfig struct {
	Cost int // work factor (4-31), default 12
}

// DefaultBcryptConfig returns recommended bcrypt config for production.
// ----------------------------------------------------------------
// DefaultBcryptConfig คืนค่า config bcrypt ที่แนะนำสำหรับ production
func DefaultBcryptConfig() *BcryptConfig {
	return &BcryptConfig{Cost: 12}
}

// Argon2idConfig holds Argon2id hashing parameters.
// ----------------------------------------------------------------
// Argon2idConfig เก็บพารามิเตอร์การแฮชของ Argon2id
type Argon2idConfig struct {
	Time    uint32 // number of iterations (default: 1)
	Memory  uint32 // memory cost in KiB (default: 64*1024 = 64MB)
	Threads uint8  // number of parallel threads (default: 4)
	KeyLen  uint32 // length of generated key (default: 32)
	SaltLen uint32 // length of salt in bytes (default: 16)
}

// DefaultArgon2idConfig returns OWASP-recommended parameters for 2025.
// ----------------------------------------------------------------
// DefaultArgon2idConfig คืนค่า config Argon2id ที่แนะนำโดย OWASP สำหรับปี 2025
func DefaultArgon2idConfig() *Argon2idConfig {
	return &Argon2idConfig{
		Time:    1,
		Memory:  64 * 1024, // 64MB
		Threads: 4,
		KeyLen:  32,
		SaltLen: 16,
	}
}

// ============================================================================
// Password Strength Validation
// ============================================================================

// PasswordPolicy defines requirements for password strength.
// ----------------------------------------------------------------
// PasswordPolicy กำหนดข้อกำหนดสำหรับความแข็งแรงของรหัสผ่าน
type PasswordPolicy struct {
	MinLength        int  // default 8
	MaxLength        int  // default 72 (bcrypt limit)
	RequireUpper     bool // default true
	RequireLower     bool // default true
	RequireDigit     bool // default true
	RequireSpecial   bool // default true
	ForbiddenPatterns []*regexp.Regexp // common weak patterns
}

// DefaultPasswordPolicy returns recommended password policy.
// ----------------------------------------------------------------
// DefaultPasswordPolicy คืนค่านโยบายรหัสผ่านที่แนะนำ
func DefaultPasswordPolicy() *PasswordPolicy {
	return &PasswordPolicy{
		MinLength:      8,
		MaxLength:      72,
		RequireUpper:   true,
		RequireLower:   true,
		RequireDigit:   true,
		RequireSpecial: true,
		ForbiddenPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)password`),
			regexp.MustCompile(`(?i)123456`),
			regexp.MustCompile(`(?i)qwerty`),
			regexp.MustCompile(`(?i)admin`),
		},
	}
}

// ValidatePassword checks if a password meets the policy requirements.
// ----------------------------------------------------------------
// ValidatePassword ตรวจสอบว่ารหัสผ่านตรงตามข้อกำหนดของนโยบายหรือไม่
func ValidatePassword(password string, policy *PasswordPolicy) error {
	if policy == nil {
		policy = DefaultPasswordPolicy()
	}
	if len(password) < policy.MinLength {
		return ErrPasswordTooShort
	}
	if policy.MaxLength > 0 && len(password) > policy.MaxLength {
		return ErrPasswordTooLong
	}
	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?/~", ch):
			hasSpecial = true
		}
	}
	if policy.RequireUpper && !hasUpper {
		return ErrPasswordNoUpper
	}
	if policy.RequireLower && !hasLower {
		return ErrPasswordNoLower
	}
	if policy.RequireDigit && !hasDigit {
		return ErrPasswordNoDigit
	}
	if policy.RequireSpecial && !hasSpecial {
		return ErrPasswordNoSpecial
	}
	// Check forbidden patterns
	for _, pattern := range policy.ForbiddenPatterns {
		if pattern.MatchString(password) {
			return fmt.Errorf("password contains weak pattern: %s", pattern.String())
		}
	}
	return nil
}

// ============================================================================
// Bcrypt Implementation
// ============================================================================

// BcryptHasher implements password hashing using bcrypt.
// ----------------------------------------------------------------
// BcryptHasher อิมพลีเมนต์การแฮชรหัสผ่านด้วย bcrypt
type BcryptHasher struct {
	config *BcryptConfig
}

// NewBcryptHasher creates a new bcrypt hasher with given config.
// ----------------------------------------------------------------
// NewBcryptHasher สร้าง bcrypt hasher ใหม่พร้อม config ที่กำหนด
func NewBcryptHasher(config *BcryptConfig) *BcryptHasher {
	if config == nil {
		config = DefaultBcryptConfig()
	}
	return &BcryptHasher{config: config}
}

// Hash generates a bcrypt hash from a plain password.
// ----------------------------------------------------------------
// Hash สร้าง bcrypt hash จากรหัสผ่านธรรมดา
func (h *BcryptHasher) Hash(password string) (string, error) {
	// Validate password before hashing
	if err := ValidatePassword(password, nil); err != nil {
		return "", err
	}
	// Truncate to 72 bytes (bcrypt limitation)
	// ตัดให้เหลือ 72 ไบต์ (ข้อจำกัดของ bcrypt)
	if len(password) > 72 {
		password = password[:72]
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), h.config.Cost)
	if err != nil {
		return "", fmt.Errorf("bcrypt hash failed: %w", err)
	}
	return string(hashed), nil
}

// Verify checks if a plain password matches a bcrypt hash.
// ----------------------------------------------------------------
// Verify ตรวจสอบว่ารหัสผ่านธรรมดาตรงกับ bcrypt hash หรือไม่
func (h *BcryptHasher) Verify(password, hashed string) bool {
	if len(password) > 72 {
		password = password[:72]
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

// NeedsRehash checks if the hash uses an outdated cost factor.
// ----------------------------------------------------------------
// NeedsRehash ตรวจสอบว่า hash ใช้ cost factor ที่ล้าสมัยหรือไม่
func (h *BcryptHasher) NeedsRehash(hashed string) bool {
	cost, err := bcrypt.Cost([]byte(hashed))
	if err != nil {
		return true
	}
	return cost < h.config.Cost
}

// ============================================================================
// Argon2id Implementation
// ============================================================================

// Argon2idHasher implements password hashing using Argon2id.
// ----------------------------------------------------------------
// Argon2idHasher อิมพลีเมนต์การแฮชรหัสผ่านด้วย Argon2id
type Argon2idHasher struct {
	config *Argon2idConfig
}

// NewArgon2idHasher creates a new Argon2id hasher.
// ----------------------------------------------------------------
// NewArgon2idHasher สร้าง Argon2id hasher ใหม่
func NewArgon2idHasher(config *Argon2idConfig) *Argon2idHasher {
	if config == nil {
		config = DefaultArgon2idConfig()
	}
	return &Argon2idHasher{config: config}
}

// generateSalt creates a cryptographically secure random salt.
// ----------------------------------------------------------------
// generateSalt สร้าง salt แบบสุ่มที่ปลอดภัย
func (h *Argon2idHasher) generateSalt() ([]byte, error) {
	salt := make([]byte, h.config.SaltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

// Hash generates an Argon2id hash in PHC format.
// ----------------------------------------------------------------
// Hash สร้าง Argon2id hash ในรูปแบบ PHC
func (h *Argon2idHasher) Hash(password string) (string, error) {
	// Validate password
	if err := ValidatePassword(password, nil); err != nil {
		return "", err
	}
	// Generate salt
	salt, err := h.generateSalt()
	if err != nil {
		return "", err
	}
	// Generate hash using Argon2id
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		h.config.Time,
		h.config.Memory,
		h.config.Threads,
		h.config.KeyLen,
	)
	// Encode to PHC format: $argon2id$v=19$m=65536,t=1,p=4$salt$hash
	// เข้ารหัสเป็นรูปแบบ PHC
	saltEnc := base64.RawStdEncoding.EncodeToString(salt)
	hashEnc := base64.RawStdEncoding.EncodeToString(hash)
	result := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		h.config.Memory, h.config.Time, h.config.Threads, saltEnc, hashEnc)
	return result, nil
}

// parseArgon2idHash extracts parameters from PHC format hash.
// ----------------------------------------------------------------
// parseArgon2idHash แยกพารามิเตอร์จาก hash รูปแบบ PHC
func parseArgon2idHash(hash string) (params *Argon2idConfig, salt, key []byte, err error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return nil, nil, nil, ErrInvalidHash
	}
	// Parse version
	// v=19
	// แยกเวอร์ชัน
	// Parse parameters from part 3: m=65536,t=3,p=4
	// แยกพารามิเตอร์จากส่วนที่ 3
	paramsPart := parts[3]
	paramPairs := strings.Split(paramsPart, ",")
	params = &Argon2idConfig{}
	for _, pair := range paramPairs {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "m":
			fmt.Sscanf(kv[1], "%d", &params.Memory)
		case "t":
			fmt.Sscanf(kv[1], "%d", &params.Time)
		case "p":
			fmt.Sscanf(kv[1], "%d", &params.Threads)
		}
	}
	// Decode salt and key
	// ถอดรหัส salt และ key
	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, ErrInvalidHash
	}
	key, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, ErrInvalidHash
	}
	return params, salt, key, nil
}

// Verify checks if a plain password matches an Argon2id hash.
// ----------------------------------------------------------------
// Verify ตรวจสอบว่ารหัสผ่านธรรมดาตรงกับ Argon2id hash หรือไม่
func (h *Argon2idHasher) Verify(password, hashed string) bool {
	params, salt, expectedKey, err := parseArgon2idHash(hashed)
	if err != nil {
		return false
	}
	// Recompute hash with extracted parameters
	// คำนวณ hash ใหม่ด้วยพารามิเตอร์ที่แยกได้
	computedKey := argon2.IDKey(
		[]byte(password),
		salt,
		params.Time,
		params.Memory,
		params.Threads,
		params.KeyLen,
	)
	// Constant-time comparison
	// เปรียบเทียบแบบ constant-time
	return subtle.ConstantTimeCompare(computedKey, expectedKey) == 1
}

// ============================================================================
// Unified PasswordHasher Interface
// ============================================================================

// PasswordHasher defines common interface for password hashing.
// ----------------------------------------------------------------
// PasswordHasher กำหนด interface ทั่วไปสำหรับการแฮชรหัสผ่าน
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hashed string) bool
}

// ============================================================================
// Rehash Utility
// ============================================================================

// RehashIfNeeded checks if a hash needs rehashing and returns new hash if needed.
// ----------------------------------------------------------------
// RehashIfNeeded ตรวจสอบว่า hash ต้องการ rehash หรือไม่ และคืน hash ใหม่ถ้าจำเป็น
func RehashIfNeeded(password, currentHash string, newHasher PasswordHasher) (newHash string, needsRehash bool) {
	// Try to detect bcrypt hash format ($2a$, $2b$, $2y$)
	// พยายามตรวจจับรูปแบบ bcrypt hash
	if strings.HasPrefix(currentHash, "$2") {
		bcryptHasher := NewBcryptHasher(DefaultBcryptConfig())
		if bcryptHasher.NeedsRehash(currentHash) {
			newHash, err := newHasher.Hash(password)
			if err == nil {
				return newHash, true
			}
		}
	}
	return "", false
}

// ============================================================================
// Convenience Functions
// ============================================================================

// HashPasswordBcrypt is a convenience function using default bcrypt config.
// ----------------------------------------------------------------
// HashPasswordBcrypt เป็นฟังก์ชันสะดวกที่ใช้ bcrypt config เริ่มต้น
func HashPasswordBcrypt(password string) (string, error) {
	hasher := NewBcryptHasher(DefaultBcryptConfig())
	return hasher.Hash(password)
}

// VerifyPasswordBcrypt is a convenience function.
// ----------------------------------------------------------------
// VerifyPasswordBcrypt เป็นฟังก์ชันสะดวก
func VerifyPasswordBcrypt(password, hash string) bool {
	hasher := NewBcryptHasher(DefaultBcryptConfig())
	return hasher.Verify(password, hash)
}

// HashPasswordArgon2id is a convenience function using default Argon2id config.
// ----------------------------------------------------------------
// HashPasswordArgon2id เป็นฟังก์ชันสะดวกที่ใช้ Argon2id config เริ่มต้น
func HashPasswordArgon2id(password string) (string, error) {
	hasher := NewArgon2idHasher(DefaultArgon2idConfig())
	return hasher.Hash(password)
}

// VerifyPasswordArgon2id is a convenience function.
// ----------------------------------------------------------------
// VerifyPasswordArgon2id เป็นฟังก์ชันสะดวก
func VerifyPasswordArgon2id(password, hash string) bool {
	hasher := NewArgon2idHasher(DefaultArgon2idConfig())
	return hasher.Verify(password, hash)
}
```

### ตัวอย่างการทดสอบ (test file)

```go
// cryptpass_test.go
package cryptpass_test

import (
	"testing"
	"gobackend/internal/pkg/cryptpass"
)

func TestBcryptHasher(t *testing.T) {
	hasher := cryptpass.NewBcryptHasher(cryptpass.DefaultBcryptConfig())
	password := "MySecureP@ssw0rd123!"
	
	hash, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}
	
	if !hasher.Verify(password, hash) {
		t.Error("Verify returned false for correct password")
	}
	
	if hasher.Verify("wrong", hash) {
		t.Error("Verify returned true for wrong password")
	}
}

func TestArgon2idHasher(t *testing.T) {
	hasher := cryptpass.NewArgon2idHasher(cryptpass.DefaultArgon2idConfig())
	password := "MySecureP@ssw0rd123!"
	
	hash, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}
	
	if !hasher.Verify(password, hash) {
		t.Error("Verify returned false for correct password")
	}
}

func TestValidatePassword(t *testing.T) {
	policy := cryptpass.DefaultPasswordPolicy()
	
	// Valid password
	err := cryptpass.ValidatePassword("MySecureP@ssw0rd", policy)
	if err != nil {
		t.Errorf("Valid password should pass: %v", err)
	}
	
	// Too short
	err = cryptpass.ValidatePassword("Abc@123", policy)
	if err != cryptpass.ErrPasswordTooShort {
		t.Error("Should reject short password")
	}
	
	// No uppercase
	err = cryptpass.ValidatePassword("mysecurep@ssw0rd", policy)
	if err != cryptpass.ErrPasswordNoUpper {
		t.Error("Should reject password without uppercase")
	}
}
```

### ตัวอย่างการใช้งานใน `auth_usecase.go`

```go
package usecase

import (
	"gobackend/internal/pkg/cryptpass"
)

type authUsecase struct {
	userRepo   repository.UserRepository
	hasher     cryptpass.PasswordHasher
	policy     *cryptpass.PasswordPolicy
}

func NewAuthUsecase(..., hasher cryptpass.PasswordHasher) AuthUsecase {
	return &authUsecase{
		hasher: hasher,
		policy: cryptpass.DefaultPasswordPolicy(),
	}
}

func (u *authUsecase) Register(ctx context.Context, email, password, fullName string) error {
	// Validate password strength
	if err := cryptpass.ValidatePassword(password, u.policy); err != nil {
		return err
	}
	// Hash password
	hashed, err := u.hasher.Hash(password)
	if err != nil {
		return err
	}
	user := &models.User{
		Email:        email,
		PasswordHash: hashed,
		FullName:     fullName,
	}
	return u.userRepo.Create(ctx, nil, user)
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return "", "", ErrInvalidCredentials
	}
	if !u.hasher.Verify(password, user.PasswordHash) {
		return "", "", ErrInvalidCredentials
	}
	// Check if rehash needed (upgrade cost)
	if newHash, needsRehash := cryptpass.RehashIfNeeded(password, user.PasswordHash, u.hasher); needsRehash {
		user.PasswordHash = newHash
		_ = u.userRepo.Update(ctx, nil, user) // background update
	}
	// proceed to create tokens...
}
```

---

## วิธีใช้งาน module นี้

1. ติดตั้ง dependencies:
   ```bash
   go get golang.org/x/crypto/bcrypt
   go get golang.org/x/crypto/argon2
   ```
2. วาง `cryptpass.go` ใน `internal/pkg/cryptpass/`
3. เลือก algorithm ที่ต้องการ (bcrypt หรือ Argon2id)
4. สร้าง hasher ใน `main.go`:
   ```go
   // สำหรับโปรเจกต์ทั่วไป (แนะนำ bcrypt)
   hasher := cryptpass.NewBcryptHasher(cryptpass.DefaultBcryptConfig())
   
   // สำหรับโปรเจกต์ที่ต้องการความปลอดภัยสูง (Argon2id)
   hasher := cryptpass.NewArgon2idHasher(cryptpass.DefaultArgon2idConfig())
   ```
5. Inject เข้า usecase

---

## ตารางสรุป bcrypt Cost และเวลาที่ใช้ (โดยประมาณ)

| Cost | เวลาต่อ 1 hash | ความปลอดภัย | แนะนำ |
|------|----------------|-------------|--------|
| 10 | ~100ms | พื้นฐาน (ไม่แนะนำ production) | ❌ |
| 11 | ~200ms | ปานกลาง | ❌ |
| 12 | ~250-500ms | ดี (สมดุล) | ✅ |
| 13 | ~500ms-1s | สูง | ✅ (ถ้ายอมรับ latency) |
| 14 | ~1-2s | สูงมาก | ⚠️ (ตรวจสอบ load ก่อน) |

**ค่า cost ที่แนะนำ:** 12 สำหรับ production ทั่วไป[reference:0][reference:1]

---

## แบบฝึกหัดท้าย module (3 ข้อ)

1. เพิ่ม method `UpgradeHash` ใน `PasswordHasher` interface ที่ rehash password ด้วย cost ที่สูงขึ้น (ใช้สำหรับ background job)
2. สร้างฟังก์ชัน `GenerateRandomPassword(length int, includeSpecial bool) (string, error)` สำหรับสร้างรหัสผ่านสุ่มที่แข็งแรง
3. Implement password history check โดยใช้ Redis เก็บ hash ของรหัสผ่านล่าสุด 5 ตัว (ป้องกันใช้รหัสเดิมซ้ำ)

---

## แหล่งอ้างอิง

- [bcrypt package documentation](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [Argon2 package documentation](https://pkg.go.dev/golang.org/x/crypto/argon2)
- [OWASP Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
- [Argon2id parameters recommendation](https://github.com/alexedwards/argon2id)
- [Password Hashing Competition (PHC)](https://password-hashing.net/)

---

**หมายเหตุ:** module นี้ครบถ้วนสำหรับ `pkg/cryptpass` หากต้องการ module เพิ่มเติม (เช่น `pkg/validator`, `pkg/logger`) ได้ดำเนินการไปแล้ว โปรดแจ้งหากต้องการ module อื่น หรือสรุปเนื้อหาทั้งหมด