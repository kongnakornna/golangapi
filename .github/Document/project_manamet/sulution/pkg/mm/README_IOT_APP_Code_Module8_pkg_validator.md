# Module 8: pkg/validator (Custom Validator)

## สำหรับโฟลเดอร์ `internal/pkg/validator/`

ไฟล์ที่เกี่ยวข้อง:
- `internal/pkg/validator/custom_validator.go`

---

## หลักการ (Concept)

### คืออะไร?
Validator คือไลบรารีหรือฟังก์ชันที่ใช้ตรวจสอบความถูกต้องของข้อมูลที่รับมาจากผู้ใช้ (request body, query parameters) ก่อนนำไปประมวลผลทางธุรกิจ ช่วยลดข้อผิดพลาดและเพิ่มความปลอดภัย

### มีกี่แบบ?
1. **Built-in validators** – เช่น `required`, `email`, `min`, `max`, `len` จาก `go-playground/validator`
2. **Custom validators** – กำหนดเองตามความต้องการของระบบ (เช่น ตรวจสอบรหัสผ่านมีความแข็งแรง, รูปแบบเบอร์โทรศัพท์)
3. **Cross-field validators** – ตรวจสอบความสัมพันธ์ระหว่างฟิลด์ (เช่น `password` กับ `confirm_password`)
4. **Database validators** – ตรวจสอบ uniqueness หรือ existence ในฐานข้อมูล (อาจเรียก repository)

### ใช้อย่างไร / นำไปใช้กรณีไหน
- ใช้ใน handler layer ก่อนเรียก usecase
- ประกาศ tag `validate:"required,email"` ใน struct field
- ลงทะเบียน custom validator ก่อนใช้

### ทำไมต้องใช้
- ป้องกันข้อมูลไม่ถูกต้องเข้าสู่ระบบ (defense in depth)
- ลด boilerplate code สำหรับการตรวจสอบ
- ให้ error message ที่ชัดเจนแก่ client

### ประโยชน์ที่ได้รับ
- ความปลอดภัยสูงขึ้น
- ลดภาระการตรวจสอบใน usecase
- มาตรฐานเดียวกันทั้งระบบ

### ข้อควรระวัง
- validation ควรอยู่ที่ delivery (handler) เท่านั้น ไม่ซ้ำใน usecase
- ระวัง performance ถ้ามี custom validator ที่ซับซ้อน
- อย่าใช้ validator สำหรับ business logic ที่ซับซ้อน

### ข้อดี
- สะอาด,  reusable,  ประกาศด้วย tag

### ข้อเสีย
- ต้องเรียน syntax tags
- custom validator ต้องเขียนและลงทะเบียน

### ข้อห้าม
- ห้ามใช้ validator เพื่อตรวจสอบข้อมูลที่ต้องเรียกฐานข้อมูล (ควรทำใน usecase)
- ห้าม skip validation ใน production

---

## โค้ดที่รันได้จริง

### ไฟล์ `internal/pkg/validator/custom_validator.go`

```go
// Package validator provides custom validation functions for go-playground/validator.
// ----------------------------------------------------------------
// แพ็คเกจ validator ให้ฟังก์ชันตรวจสอบแบบกำหนดเองสำหรับ go-playground/validator
package validator

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// CustomValidator wraps the validator instance with custom tags.
// ----------------------------------------------------------------
// CustomValidator ห่อหุ้ม instance ของ validator พร้อม custom tags
type CustomValidator struct {
	Validate *validator.Validate
}

// NewCustomValidator creates a new validator instance and registers custom validations.
// ----------------------------------------------------------------
// NewCustomValidator สร้าง instance validator ใหม่และลงทะเบียน validations แบบกำหนดเอง
func NewCustomValidator() *CustomValidator {
	v := validator.New()
	cv := &CustomValidator{Validate: v}
	cv.registerCustomValidations()
	return cv
}

// registerCustomValidations registers all custom validation functions.
// ----------------------------------------------------------------
// registerCustomValidations ลงทะเบียนฟังก์ชันตรวจสอบแบบกำหนดเองทั้งหมด
func (cv *CustomValidator) registerCustomValidations() {
	// Register password strength validator
	// ลงทะเบียนตัวตรวจสอบความแข็งแรงของรหัสผ่าน
	cv.Validate.RegisterValidation("password_strength", cv.validatePasswordStrength)
	
	// Register Thai phone number validator
	// ลงทะเบียนตัวตรวจสอบเบอร์โทรศัพท์ไทย
	cv.Validate.RegisterValidation("thai_phone", cv.validateThaiPhone)
	
	// Register username validator (alphanumeric + underscore, no spaces)
	// ลงทะเบียนตัวตรวจสอบชื่อผู้ใช้ (ตัวอักษร, ตัวเลข, ขีดล่าง, ไม่มีเว้นวรรค)
	cv.Validate.RegisterValidation("username", cv.validateUsername)
	
	// Register one of fields match validator (e.g., password == confirm_password)
	// ลงทะเบียนตัวตรวจสอบการจับคู่ระหว่างฟิลด์
	cv.Validate.RegisterValidation("eqfield_custom", cv.validateEqualField)
}

// validatePasswordStrength checks if password meets security requirements:
// - at least 8 characters
// - at least one uppercase letter
// - at least one lowercase letter
// - at least one digit
// - at least one special character (!@#$%^&*)
// ----------------------------------------------------------------
// validatePasswordStrength ตรวจสอบว่ารหัสผ่านตรงตามข้อกำหนดด้านความปลอดภัย:
// - อย่างน้อย 8 ตัวอักษร
// - อย่างน้อย 1 ตัวพิมพ์ใหญ่
// - อย่างน้อย 1 ตัวพิมพ์เล็ก
// - อย่างน้อย 1 ตัวเลข
// - อย่างน้อย 1 อักขระพิเศษ (!@#$%^&*)
func (cv *CustomValidator) validatePasswordStrength(fl validator.FieldLevel) bool {
	password, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	if len(password) < 8 {
		return false
	}
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	specialChars := "!@#$%^&*"
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasNumber = true
		case strings.ContainsRune(specialChars, ch):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}

// validateThaiPhone validates Thai mobile phone numbers.
// Formats accepted: 08xxxxxxxx, 09xxxxxxxx, 06xxxxxxxx (10 digits)
// ----------------------------------------------------------------
// validateThaiPhone ตรวจสอบเบอร์โทรศัพท์มือถือไทย
// รูปแบบที่รองรับ: 08xxxxxxxx, 09xxxxxxxx, 06xxxxxxxx (10 หลัก)
func (cv *CustomValidator) validateThaiPhone(fl validator.FieldLevel) bool {
	phone, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	// Remove spaces and hyphens
	// ลบช่องว่างและขีดกลาง
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	// Check length and prefix
	if len(cleaned) != 10 {
		return false
	}
	prefix := cleaned[:2]
	validPrefixes := []string{"08", "09", "06"}
	for _, p := range validPrefixes {
		if prefix == p {
			return true
		}
	}
	return false
}

// validateUsername checks for valid username:
// - length between 3 and 30 characters
// - only alphanumeric and underscore
// - cannot start or end with underscore
// ----------------------------------------------------------------
// validateUsername ตรวจสอบชื่อผู้ใช้:
// - ความยาวระหว่าง 3-30 ตัวอักษร
// - ประกอบด้วยตัวอักษร ตัวเลข และขีดล่างเท่านั้น
// - ห้ามขึ้นต้นหรือลงท้ายด้วยขีดล่าง
func (cv *CustomValidator) validateUsername(fl validator.FieldLevel) bool {
	username, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	if len(username) < 3 || len(username) > 30 {
		return false
	}
	// Regex: alphanumeric and underscore, not starting/ending with underscore
	regex := `^[a-zA-Z0-9][a-zA-Z0-9_]*[a-zA-Z0-9]$`
	matched, _ := regexp.MatchString(regex, username)
	return matched
}

// validateEqualField checks if field value equals another field's value.
// Use tag: validate:"eqfield_custom=ConfirmPassword"
// ----------------------------------------------------------------
// validateEqualField ตรวจสอบว่าค่าฟิลด์นี้เท่ากับอีกฟิลด์หนึ่ง
// ใช้ tag: validate:"eqfield_custom=ConfirmPassword"
func (cv *CustomValidator) validateEqualField(fl validator.FieldLevel) bool {
	field := fl.Field()
	param := fl.Param() // name of the other field
	otherField := fl.Parent().FieldByName(param)
	if !otherField.IsValid() {
		return false
	}
	return field.Interface() == otherField.Interface()
}

// ValidateStruct validates a struct using registered validators.
// Returns validation errors as map of field -> error message.
// ----------------------------------------------------------------
// ValidateStruct ตรวจสอบ struct โดยใช้ validators ที่ลงทะเบียนไว้
// คืนค่า validation errors เป็น map ของฟิลด์ -> ข้อความผิดพลาด
func (cv *CustomValidator) ValidateStruct(s interface{}) map[string]string {
	err := cv.Validate.Struct(s)
	if err == nil {
		return nil
	}
	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()
		param := err.Param()
		message := getErrorMessage(field, tag, param)
		errors[field] = message
	}
	return errors
}

// getErrorMessage returns human-readable error message for validation tag.
// ----------------------------------------------------------------
// getErrorMessage คืนข้อความ error ที่อ่านง่ายสำหรับ validation tag
func getErrorMessage(field, tag, param string) string {
	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be at least " + param + " characters"
	case "max":
		return field + " must be at most " + param + " characters"
	case "password_strength":
		return field + " must be at least 8 characters with uppercase, lowercase, number, and special character"
	case "thai_phone":
		return field + " must be a valid Thai mobile number (08/09/06 followed by 8 digits)"
	case "username":
		return field + " must be 3-30 characters, alphanumeric and underscore only, no leading/trailing underscore"
	case "eqfield_custom":
		return field + " does not match " + param
	default:
		return field + " is invalid"
	}
}
```

### ตัวอย่างการนำไปใช้ใน DTO (ภายใน `dto/auth_dto.go`)

```go
package dto

// RegisterRequest with custom validators
// ----------------------------------------------------------------
// RegisterRequest พร้อม custom validators
type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,password_strength"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield_custom=Password"`
	Username        string `json:"username" validate:"required,username"`
	Phone           string `json:"phone" validate:"omitempty,thai_phone"`
}
```

### ตัวอย่างการใช้งานใน Handler

```go
// ใน auth_handler.go
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}
	
	// ใช้ custom validator
	if errs := h.validator.ValidateStruct(req); len(errs) > 0 {
		respondValidationErrors(w, errs)
		return
	}
	
	// ... proceed to usecase
}
```

### การเริ่มต้นใน `main.go` หรือ `router.go`

```go
// สร้าง validator instance
customValidator := validator.NewCustomValidator()

// ส่งให้ handler
authHandler := handler.NewAuthHandler(authUsecase, customValidator)
```

---

## วิธีใช้งาน module นี้

1. ติดตั้ง dependency:
   ```bash
   go get github.com/go-playground/validator/v10
   ```
2. วาง `custom_validator.go` ใน `internal/pkg/validator/`
3. ใน handler struct ให้เพิ่ม field `validator *validator.CustomValidator`
4. เรียก `h.validator.ValidateStruct(req)` ก่อนส่งข้อมูลไป usecase

---

## ตารางสรุป Custom Validators

| ชื่อ tag | ตรวจสอบ | ตัวอย่างค่าที่ผ่าน | ตัวอย่างค่าที่ไม่ผ่าน |
|----------|---------|-------------------|----------------------|
| `password_strength` | ความแข็งแรงรหัสผ่าน | `P@ssw0rd` | `password123` |
| `thai_phone` | เบอร์โทรศัพท์ไทย | `0812345678` | `021234567` |
| `username` | ชื่อผู้ใช้ (3-30, alnum + _) | `john_doe` | `john!doe` |
| `eqfield_custom` | ค่าตรงกับอีกฟิลด์ | `password` = `confirm_password` | ต่างกัน |

---

## แบบฝึกหัดท้าย module (3 ข้อ)

1. เพิ่ม custom validator `no_profanity` ที่ตรวจสอบว่าข้อความไม่มีคำหยาบ (ใช้ list คำต้องห้าม)
2. สร้าง validator `date_after` สำหรับตรวจสอบว่า `end_date` อยู่หลัง `start_date`
3. ปรับปรุง `validateThaiPhone` ให้รองรับเบอร์ที่มีเครื่องหมาย `+66` นำหน้า (แปลงเป็น 0xx)

---

## แหล่งอ้างอิง

- [go-playground/validator documentation](https://pkg.go.dev/github.com/go-playground/validator/v10)
- [Custom validators example](https://github.com/go-playground/validator#custom-functions)
- [Validation best practices](https://www.alexedwards.net/blog/validation-tips)

---

**หมายเหตุ:** module นี้ครบถ้วนสำหรับ `pkg/validator` ตามโครงสร้าง gobackend หากต้องการ module ถัดไป (เช่น `pkg/logger`, `pkg/redis`, `pkg/email`, `pkg/jwt`) โปรดแจ้ง