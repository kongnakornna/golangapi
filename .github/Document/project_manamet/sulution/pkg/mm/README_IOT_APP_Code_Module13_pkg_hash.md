# Module 13: pkg/hash (Password Hashing)

## สำหรับโฟลเดอร์ `internal/pkg/hash/`

ไฟล์ที่เกี่ยวข้อง:
- `internal/pkg/hash/bcrypt.go`

---

## หลักการ (Concept)

### คืออะไร?
Password hashing คือกระบวนการแปลงรหัสผ่านข้อความธรรมดา (plain text) ให้เป็นค่าแฮช (hash) ที่ไม่สามารถย้อนกลับไปหารหัสผ่านเดิมได้ เพื่อเก็บในฐานข้อมูลอย่างปลอดภัย ป้องกันการถูกขโมยรหัสผ่านเมื่อ database ถูกโจมตี

### มีกี่แบบ?
1. **bcrypt** – ปัจจุบันนิยมใช้มากที่สุด, ปรับค่า cost ได้, ทนต่อ brute force ด้วยความช้าโดย design
2. **scrypt** – ใช้ memory-hard function เหมาะกับ environment ที่มี RAM จำกัด
3. **argon2** – มาตรฐานใหม่ (2015), ปลอดภัยที่สุด แต่ซับซ้อนกว่า
4. **PBKDF2** – มาตรฐาน NIST, ใช้ได้ แต่ช้ากว่า bcrypt เมื่อ cost สูง

### ใช้อย่างไร / นำไปใช้กรณีไหน
- ใช้ `bcrypt.GenerateFromPassword` เมื่อสร้างผู้ใช้ใหม่หรือเปลี่ยนรหัสผ่าน
- ใช้ `bcrypt.CompareHashAndPassword` เมื่อตรวจสอบรหัสผ่านตอน login
- ปรับ cost ให้เหมาะสม (ค่าเริ่มต้น 10, production อาจใช้ 12-14)

### ทำไมต้องใช้
- ห้ามเก็บรหัสผ่านใน plain text เด็ดขาด (GDPR, PCI-DSS กำหนด)
- การ hash ช่วยป้องกันข้อมูลรั่วไหล
- bcrypt มี salt ในตัว ป้องกัน rainbow table attack

### ประโยชน์ที่ได้รับ
- ปลอดภัย แม้ database ถูกขโมย
- ป้องกัน brute force เพราะ hash ช้า (故意设计)
- มี salt อัตโนมัติ

### ข้อควรระวัง
- cost สูงเกินไปจะทำให้ login ช้า (CPU bound)
- hash ไม่สามารถ reverse กลับเป็นรหัสผ่านเดิมได้ (ถ้าลืมรหัส ต้อง reset)
- bcrypt จำกัดความยาวรหัสผ่านที่ 72 ไบต์ (ควรตัดก่อน)

### ข้อดี
- ปลอดภัย, มี salt ในตัว, ปรับความเร็วได้

### ข้อเสีย
- ช้ากว่า hashing สำหรับ checksum (แต่这是优点)
- ไม่เหมาะกับระบบที่ต้อง validate password เร็วมาก

### ข้อห้าม
- ห้ามใช้ MD5, SHA-1, SHA-256 (เร็วเกินไป) สำหรับรหัสผ่าน
- ห้ามเก็บรหัสผ่านใน plain text
- ห้ามใช้ hash โดยไม่ใส่ salt

---

## โค้ดที่รันได้จริง

### ไฟล์ `internal/pkg/hash/bcrypt.go`

```go
// Package hash provides password hashing using bcrypt.
// ----------------------------------------------------------------
// แพ็คเกจ hash ให้บริการ hashing รหัสผ่านด้วย bcrypt
package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// DefaultCost is the bcrypt cost factor used for hashing.
// Higher cost = more secure but slower.
// ----------------------------------------------------------------
// DefaultCost คือค่า cost factor ของ bcrypt ที่ใช้ในการแฮช
// ค่า越高 = ปลอดภัยมากขึ้น แต่ช้าลง
const DefaultCost = 12

// PasswordHasher defines interface for password hashing operations.
// ----------------------------------------------------------------
// PasswordHasher กำหนด interface สำหรับการแฮชรหัสผ่าน
type PasswordHasher interface {
	// Hash hashes a plain password and returns the hashed string.
	// ----------------------------------------------------------------
	// Hash แฮชรหัสผ่านธรรมดาและคืนค่า string ที่ถูกแฮช
	Hash(password string) (string, error)

	// Verify checks if a plain password matches a hash.
	// ----------------------------------------------------------------
	// Verify ตรวจสอบว่ารหัสผ่านธรรมดาตรงกับแฮชหรือไม่
	Verify(password, hash string) bool
}

// BcryptHasher implements PasswordHasher using bcrypt.
// ----------------------------------------------------------------
// BcryptHasher อิมพลีเมนต์ PasswordHasher ด้วย bcrypt
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher creates a new bcrypt hasher with default cost.
// ----------------------------------------------------------------
// NewBcryptHasher สร้าง bcrypt hasher ใหม่พร้อมค่า cost เริ่มต้น
func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{cost: DefaultCost}
}

// NewBcryptHasherWithCost creates a bcrypt hasher with custom cost.
// ----------------------------------------------------------------
// NewBcryptHasherWithCost สร้าง bcrypt hasher พร้อมค่า cost ที่กำหนดเอง
func NewBcryptHasherWithCost(cost int) *BcryptHasher {
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}
	if cost > bcrypt.MaxCost {
		cost = bcrypt.MaxCost
	}
	return &BcryptHasher{cost: cost}
}

// Hash hashes a plain password.
// ----------------------------------------------------------------
// Hash แฮชรหัสผ่านธรรมดา
func (h *BcryptHasher) Hash(password string) (string, error) {
	// Truncate password to 72 bytes (bcrypt limit)
	// ตัดรหัสผ่านให้เหลือ 72 ไบต์ (ข้อจำกัดของ bcrypt)
	if len(password) > 72 {
		password = password[:72]
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// Verify checks if plain password matches the hash.
// ----------------------------------------------------------------
// Verify ตรวจสอบว่ารหัสผ่านธรรมดาตรงกับแฮชหรือไม่
func (h *BcryptHasher) Verify(password, hash string) bool {
	// Truncate password to 72 bytes
	// ตัดรหัสผ่านให้เหลือ 72 ไบต์
	if len(password) > 72 {
		password = password[:72]
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HashPassword is a convenience function using default cost.
// ----------------------------------------------------------------
// HashPassword เป็นฟังก์ชันสะดวกใช้ค่า cost เริ่มต้น
func HashPassword(password string) (string, error) {
	hasher := NewBcryptHasher()
	return hasher.Hash(password)
}

// VerifyPassword is a convenience function.
// ----------------------------------------------------------------
// VerifyPassword เป็นฟังก์ชันสะดวก
func VerifyPassword(password, hash string) bool {
	hasher := NewBcryptHasher()
	return hasher.Verify(password, hash)
}
```

### ตัวอย่างการใช้งานใน `auth_usecase.go`

```go
package usecase

import (
	"gobackend/internal/pkg/hash"
)

type authUsecase struct {
	userRepo   repository.UserRepository
	hasher     hash.PasswordHasher
	// ... other fields
}

func NewAuthUsecase(..., hasher hash.PasswordHasher) AuthUsecase {
	return &authUsecase{
		hasher: hasher,
		// ...
	}
}

func (u *authUsecase) Register(ctx context.Context, email, password, fullName string) error {
	hashedPassword, err := u.hasher.Hash(password)
	if err != nil {
		return err
	}
	user := &models.User{
		Email:        email,
		PasswordHash: hashedPassword,
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
	// proceed to create tokens...
}
```

### ตัวอย่างการทดสอบ (test file)

```go
// hash/bcrypt_test.go
package hash_test

import (
	"testing"
	"gobackend/internal/pkg/hash"
)

func TestBcryptHasher(t *testing.T) {
	hasher := hash.NewBcryptHasher()
	
	password := "MySecurePassword123!"
	
	hashed, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}
	
	if !hasher.Verify(password, hashed) {
		t.Error("Verify returned false for correct password")
	}
	
	if hasher.Verify("wrong", hashed) {
		t.Error("Verify returned true for wrong password")
	}
}
```

---

## วิธีใช้งาน module นี้

1. ติดตั้ง dependency:
   ```bash
   go get golang.org/x/crypto/bcrypt
   ```
2. วาง `bcrypt.go` ใน `internal/pkg/hash/`
3. สร้าง hasher instance ใน `main.go`:
   ```go
   hasher := hash.NewBcryptHasherWithCost(12)
   ```
4. Inject เข้า usecase ที่ต้องการ

---

## ตารางสรุป bcrypt Cost และเวลาที่ใช้ (โดยประมาณ)

| Cost | เวลาต่อ 1 hash (CPU) | เวลา brute force (8 char, 95 chars) |
|------|---------------------|--------------------------------------|
| 10   | ~80ms               | ~2.5 ปี                              |
| 12   | ~250ms              | ~8 ปี                                |
| 14   | ~800ms              | ~25 ปี                               |

> หมายเหตุ: ค่า cost ที่เหมาะสมสำหรับ production คือ 12 (balance ระหว่างความปลอดภัยและ用户体验)

---

## แบบฝึกหัดท้าย module (3 ข้อ)

1. เพิ่ม method `NeedsRehash(hash string) bool` ที่ตรวจสอบว่า hash ใช้ cost ต่ำกว่าค่า default หรือไม่ (เพื่อ upgrade cost เมื่อ login)
2. Implement `Argon2Hasher` ที่ใช้ argon2id (ศึกษา package `golang.org/x/crypto/argon2`)
3. สร้างฟังก์ชัน `ValidatePasswordStrength(password string) error` ที่ตรวจสอบความยาว, ตัวพิมพ์ใหญ่, ตัวเลข, อักขระพิเศษ ก่อน hash

---

## แหล่งอ้างอิง

- [bcrypt package documentation](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [Why bcrypt?](https://security.stackexchange.com/questions/4781/do-any-security-experts-recommend-bcrypt-for-password-storage)
- [OWASP Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)

---

**หมายเหตุ:** module นี้ครบถ้วนสำหรับ `pkg/hash` หากต้องการ module เพิ่มเติม (เช่น `pkg/validator`, `pkg/logger`) ได้ดำเนินการไปแล้ว โปรดแจ้งหากต้องการ module อื่น หรือสรุปเนื้อหาทั้งหมด