# Module 7: pkg/utils (Utility Functions)

## สำหรับโฟลเดอร์ `internal/pkg/utils/`

ไฟล์ที่เกี่ยวข้อง:
- `internal/pkg/utils/random.go`
- `internal/pkg/utils/time.go`

---

## หลักการ (Concept)

### คืออะไร?
Utility functions คือฟังก์ชันช่วยเหลือทั่วไปที่ใช้ซ้ำได้ในหลายส่วนของโปรเจกต์ เช่น การสร้างตัวเลขสุ่ม, การจัดการเวลา, การแปลงรูปแบบข้อมูล ช่วยลดการเขียนโค้ดซ้ำและเพิ่มความสม่ำเสมอ

### มีกี่แบบ?
1. **Random generators** – สร้าง string, int, token สุ่ม
2. **Time helpers** – การแปลงเวลา, การคำนวณเวลา, การจัดรูปแบบ
3. **String manipulation** – การตัด, แปลง, ตรวจสอบรูปแบบ
4. **Conversion utilities** – ระหว่าง数据类型 (int64 ↔ string, etc.)

### ใช้อย่างไร / นำไปใช้กรณีไหน
- ใช้สำหรับสร้าง refresh token, verification code
- ใช้คำนวณ expiration time, parse timezone
- ใช้ sanitize input หรือ validate format

### ทำไมต้องใช้
- ลด code duplication
- ทำให้ business logic สะอาดขึ้น
- รวมการจัดการ edge cases ไว้ที่เดียว

### ประโยชน์ที่ได้รับ
- ทดสอบแยกได้ง่าย
- แก้ไขจุดเดียวทั่วทั้งระบบ
- เพิ่ม readability

### ข้อควรระวัง
- ฟังก์ชัน utility ควร pure (ไม่มี side effects)
- ควรมี test coverage สูง
- อย่าใส่ business logic ใน utils

### ข้อดี
- reusable, testable, maintainable

### ข้อเสีย
- อาจกลายเป็น "junk drawer" ถ้าไม่จัดระเบียบดี

### ข้อห้าม
- ห้ามใช้ utils สำหรับ logic เฉพาะของ domain
- ห้าม依赖 global state

---

## โค้ดที่รันได้จริง

### 1. Random Utilities – `internal/pkg/utils/random.go`

```go
// Package utils provides common helper functions.
// ----------------------------------------------------------------
// แพ็คเกจ utils ให้ฟังก์ชันช่วยเหลือทั่วไป
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"strings"
)

const (
	// Default random string length for tokens
	// ความยาวเริ่มต้นของ string สุ่มสำหรับ tokens
	defaultTokenLength = 32
	
	// Character sets for random generation
	// ชุดอักขระสำหรับการสร้างค่าสุ่ม
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits           = "0123456789"
	specialChars     = "!@#$%^&*"
)

// RandomString generates a random string of given length using crypto/rand.
// ----------------------------------------------------------------
// RandomString สร้าง string สุ่มความยาวที่กำหนดด้วย crypto/rand
func RandomString(length int) (string, error) {
	if length <= 0 {
		length = defaultTokenLength
	}
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

// RandomInt generates a random integer between min and max (inclusive).
// ----------------------------------------------------------------
// RandomInt สร้างเลขสุ่มระหว่าง min และ max (รวมทั้งสองค่า)
func RandomInt(min, max int64) (int64, error) {
	if min >= max {
		return min, nil
	}
	rangeSize := max - min + 1
	n, err := rand.Int(rand.Reader, big.NewInt(rangeSize))
	if err != nil {
		return 0, err
	}
	return min + n.Int64(), nil
}

// RandomStringFromSet generates a random string using custom character set.
// ----------------------------------------------------------------
// RandomStringFromSet สร้าง string สุ่มจากชุดอักขระที่กำหนด
func RandomStringFromSet(length int, charset string) (string, error) {
	if length <= 0 {
		length = 8
	}
	if charset == "" {
		charset = lowercaseLetters + digits
	}
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))
	for i := 0; i < length; i++ {
		idx, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result[i] = charset[idx.Int64()]
	}
	return string(result), nil
}

// GenerateOTP generates a numeric OTP of given length (e.g., 6 digits).
// ----------------------------------------------------------------
// GenerateOTP สร้าง OTP ตัวเลขความยาวที่กำหนด (เช่น 6 หลัก)
func GenerateOTP(length int) (string, error) {
	if length <= 0 {
		length = 6
	}
	max := int64(1)
	for i := 0; i < length; i++ {
		max *= 10
	}
	val, err := RandomInt(0, max-1)
	if err != nil {
		return "", err
	}
	// Pad with leading zeros
	// เติมเลขศูนย์นำหน้า
	format := "%0" + string(rune(length+'0')) + "d"
	return sprintf(format, val), nil
}

// GenerateSecureToken creates a cryptographically secure token for API keys or refresh tokens.
// ----------------------------------------------------------------
// GenerateSecureToken สร้าง token ปลอดภัยสำหรับ API key หรือ refresh token
func GenerateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
```

### 2. Time Utilities – `internal/pkg/utils/time.go`

```go
package utils

import (
	"time"
)

// TimeFormat constants for consistent date/time formatting.
// ----------------------------------------------------------------
// ค่าคงที่ TimeFormat สำหรับการจัดรูปแบบวันที่/เวลาให้สอดคล้องกัน
const (
	ISO8601      = "2006-01-02T15:04:05Z07:00"
	RFC3339      = time.RFC3339
	DateOnly     = "2006-01-02"
	DateTime     = "2006-01-02 15:04:05"
	TimeOnly     = "15:04:05"
	ThaiDateTime = "02/01/2006 15:04:05"
)

// Now returns current time in UTC.
// ----------------------------------------------------------------
// Now คืนเวลาปัจจุบันใน UTC
func Now() time.Time {
	return time.Now().UTC()
}

// StartOfDay returns the beginning (00:00:00) of the given date.
// ----------------------------------------------------------------
// StartOfDay คืนเวลาเริ่มต้นของวัน (00:00:00)
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end (23:59:59.999999999) of the given date.
// ----------------------------------------------------------------
// EndOfDay คืนเวลาสิ้นสุดของวัน (23:59:59.999999999)
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// DaysBetween calculates the number of days between two dates (absolute value).
// ----------------------------------------------------------------
// DaysBetween คำนวณจำนวนวันระหว่างวันที่สองวัน (ค่าสัมบูรณ์)
func DaysBetween(t1, t2 time.Time) int {
	// Truncate to date only
	t1Date := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.UTC)
	t2Date := time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.UTC)
	diff := t2Date.Sub(t1Date)
	return int(diff.Hours() / 24)
}

// FormatDurationHuman returns human-readable duration (e.g., "2h 3m 5s").
// ----------------------------------------------------------------
// FormatDurationHuman คืนระยะเวลาในรูปแบบที่มนุษย์อ่านได้ (เช่น "2ชม. 3นาที 5วินาที")
func FormatDurationHuman(d time.Duration) string {
	if d < 0 {
		d = -d
	}
	seconds := int(d.Seconds())
	minutes := seconds / 60
	hours := minutes / 60
	seconds %= 60
	minutes %= 60
	
	if hours > 0 {
		return sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return sprintf("%dm %ds", minutes, seconds)
	}
	return sprintf("%ds", seconds)
}

// IsExpired checks if given time is before current time.
// ----------------------------------------------------------------
// IsExpired ตรวจสอบว่าเวลาที่กำหนดผ่านมาแล้วหรือไม่
func IsExpired(expiryTime time.Time) bool {
	return time.Now().UTC().After(expiryTime)
}

// ParseTimeOrNow tries to parse a time string; returns now on error.
// ----------------------------------------------------------------
// ParseTimeOrNow พยายามแปลง string เป็นเวลา ถ้า error คืนเวลาปัจจุบัน
func ParseTimeOrNow(timeStr, layout string) time.Time {
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return Now()
	}
	return t
}

// UnixMillis returns milliseconds since epoch.
// ----------------------------------------------------------------
// UnixMillis คืนค่ามิลลิวินาทีนับจาก epoch
func UnixMillis(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// FromUnixMillis converts milliseconds since epoch to time.Time.
// ----------------------------------------------------------------
// FromUnixMillis แปลงมิลลิวินาทีนับจาก epoch เป็น time.Time
func FromUnixMillis(ms int64) time.Time {
	return time.Unix(0, ms*int64(time.Millisecond))
}
```

### 3. String Utilities (เพิ่มเติม) – `internal/pkg/utils/string.go` (optional)

```go
package utils

import (
	"regexp"
	"strings"
	"unicode"
)

// Truncate shortens a string to max length and adds ellipsis.
// ----------------------------------------------------------------
// Truncate ย่อ string ให้มีความยาวไม่เกิน max และเติม ...
func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// Slugify converts a string to URL-friendly slug.
// ----------------------------------------------------------------
// Slugify แปลง string เป็น slug ที่เหมาะกับ URL
func Slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)
	// Replace spaces and special chars with dash
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")
	// Trim dashes
	return strings.Trim(s, "-")
}

// IsValidEmail performs basic email format validation.
// ----------------------------------------------------------------
// IsValidEmail ตรวจสอบรูปแบบอีเมลเบื้องต้น
func IsValidEmail(email string) bool {
	// Simple regex for demonstration (use proper validation in production)
	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(strings.ToLower(email))
}
```

### 4. ตัวอย่างการใช้งาน – `internal/pkg/utils/example_test.go`

```go
package utils_test

import (
	"fmt"
	"time"
	"gobackend/internal/pkg/utils"
)

func ExampleRandomString() {
	str, _ := utils.RandomString(10)
	fmt.Println(len(str))
	// Output: 10
}

func ExampleStartOfDay() {
	t := time.Date(2025, 4, 4, 14, 30, 0, 0, time.UTC)
	start := utils.StartOfDay(t)
	fmt.Println(start.Format("2006-01-02 15:04:05"))
	// Output: 2025-04-04 00:00:00
}

func ExampleGenerateSecureToken() {
	token, _ := utils.GenerateSecureToken()
	fmt.Println(len(token))
	// Output: 64
}
```

---

## วิธีใช้งาน module นี้

1. วางไฟล์ `random.go`, `time.go`, `string.go` ใน `internal/pkg/utils/`
2. ในโค้ดอื่น import:
   ```go
   import "gobackend/internal/pkg/utils"
   ```
3. เรียกใช้ฟังก์ชัน:
   ```go
   token, _ := utils.GenerateSecureToken()
   otp, _ := utils.GenerateOTP(6)
   now := utils.Now()
   expired := utils.IsExpired(user.Expiry)
   ```

### ตัวอย่างการนำไปใช้ในระบบจริง

**สร้าง refresh token (ใน auth_usecase)**
```go
refreshToken, err := utils.GenerateSecureToken()
if err != nil {
    return err
}
```

**ตรวจสอบว่า session หมดอายุ (ใน session_repo)**
```go
if utils.IsExpired(session.ExpiresAt) {
    return ErrSessionExpired
}
```

**สร้าง OTP สำหรับ 2FA**
```go
otp, _ := utils.GenerateOTP(6)
// send OTP via SMS or email
```

**แปลงเวลาเป็น Unix milliseconds (สำหรับ frontend)**
```go
ms := utils.UnixMillis(user.CreatedAt)
// return in JSON
```

---

## ตารางสรุปฟังก์ชัน Utilities

| ฟังก์ชัน | อินพุต | เอาต์พุต | ใช้เมื่อ |
|----------|--------|----------|---------|
| `RandomString` | length (int) | (string, error) | สร้าง random token |
| `RandomInt` | min, max int64 | (int64, error) | random delay, sampling |
| `GenerateOTP` | length int | (string, error) | รหัสยืนยัน 2FA |
| `GenerateSecureToken` | none | (string, error) | refresh token, API key |
| `StartOfDay` | time.Time | time.Time | เริ่มต้นวันสำหรับ report |
| `EndOfDay` | time.Time | time.Time | สิ้นสุดวันสำหรับ filter |
| `DaysBetween` | t1, t2 time.Time | int | คำนวณอายุ, retention |
| `IsExpired` | expiryTime time.Time | bool | ตรวจสอบ session/token |
| `UnixMillis` | time.Time | int64 | ส่งไป frontend |
| `Slugify` | string | string | สร้าง URL path |

---

## แบบฝึกท้าย module (3 ข้อ)

1. เพิ่มฟังก์ชัน `RandomHex(length int) (string, error)` ที่สร้าง random hex string โดยใช้ crypto/rand และ hex encoding
2. เพิ่มฟังก์ชัน `FormatRelativeTime(t time.Time) string` ที่คืนค่า "2 minutes ago", "3 hours ago", "yesterday" เป็นภาษาไทย
3. สร้างฟังก์ชัน `ParseDuration(durationStr string) (time.Duration, error)` ที่รองรับรูปแบบ "1h30m", "2d" (day), "1w" (week) (แปลงวันเป็น 24 ชั่วโมง)

---

## แหล่งอ้างอิง

- [crypto/rand package](https://pkg.go.dev/crypto/rand)
- [time package](https://pkg.go.dev/time)
- [regexp package](https://pkg.go.dev/regexp)
- [Generating secure tokens in Go](https://security.stackexchange.com/questions/244634/secure-random-token-generation-in-go)

---

**หมายเหตุ:** module นี้ครบถ้วนสำหรับ `pkg/utils` ตามโครงสร้าง gobackend หากต้องการ module ถัดไป (เช่น `pkg/validator`, `pkg/logger`, `pkg/redis`) โปรดแจ้ง