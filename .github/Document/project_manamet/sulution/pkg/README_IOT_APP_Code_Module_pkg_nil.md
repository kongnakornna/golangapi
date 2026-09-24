นี่คือ **Module การจัดการ Nil ใน Go** สำหรับโฟลเดอร์ `pkg/nil/` ตาม Template ที่ให้มา

---

# Module nil: pkg/nil

## สำหรับโฟลเดอร์ `pkg/nil/`

ไฟล์ที่เกี่ยวข้อง:
- `safe.go` - ฟังก์ชันช่วยจัดการ nil
- `safe_test.go` - หน่วยทดสอบ
- `example_test.go` - ตัวอย่างการใช้งาน

---

## หลักการ (Concept)

### Nil คืออะไร?
ใน Go `nil` เป็นค่าศูนย์ (zero value) สำหรับ pointer, interface, slice, map, channel, และ function การเข้าถึงผ่าน nil pointer ทำให้เกิด **panic: runtime error**

### มีกี่แบบ?

| ประเภท | ค่า nil | panic เมื่อ |
|--------|---------|-------------|
| `*T` (pointer) | `nil` | เรียก method หรือเข้าถึง field |
| `interface{}` | `nil` | เรียก method |
| `[]T` (slice) | `nil` | append ใช้ได้, index ใช้ไม่ได้ |
| `map[K]V` | `nil` | อ่านค่าได้ (คืน zero), เขียน panic |
| `chan T` | `nil` | send/receive ตายถาวร |
| `func()` | `nil` | เรียกใช้ panic |

**ข้อห้ามสำคัญ:** ห้ามใช้ Bucket Pattern ร่วมกับ Time Series Collections เพราะจะลดประสิทธิภาพ

### ใช้อย่างไร / นำไปใช้กรณีไหน
- ก่อน dereference pointer ทุกครั้ง
- ก่อนเรียก method บน interface
- ก่อนอ่าน/เขียน map ที่อาจเป็น nil
- ก่อนส่ง/รับ channel
- เมื่อรับค่าจาก external system (JSON, DB, API)

### ประโยชน์ที่ได้รับ
- ป้องกัน panic ใน production
- โค้ด robust ขึ้น
- Debug ง่ายขึ้น (error แทน panic)

### ข้อควรระวัง
- `nil slice` ใช้ `append` ได้ แต่ `nil map` ใช้เขียนไม่ได้
- `nil interface` ≠ `nil concrete type`
- การ over-check nil ทำให้โค้ดรก

### ข้อดี
- ควบคุมการทำงานได้
- ป้องกัน crash
- ใช้ pattern เดียวกันทั้ง project

### ข้อเสีย
- เพิ่ม boilerplate code
- ลืมตรวจได้ง่าย
- ส่งผลต่อ performance เล็กน้อย

### ข้อห้าม
- ห้าม ignore nil check สำหรับ pointer ที่รับจาก external
- ห้าม assume ว่า struct field ไม่เป็น nil
- ห้ามใช้ `nil` เป็น valid value โดยไม่มีเอกสาร

---

## การออกแบบ Workflow และ Dataflow

```
[External Input] → [Nil Checker] → [Safe Accessor] → [Result/Error]
       ↓                  ↓                ↓
    JSON/DB           IsNil()          GetOrZero()
    API/User          IsZero()         GetOrDefault()
                     Validate()        SafeCall()
```

---

## ตัวอย่างโค้ดที่รันได้จริง

```go
package nil

import (
	"reflect"
)

// SafeDeref คืนค่า pointer value หรือ zero value ถ้า nil
func SafeDeref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// SafeDerefWithDefault คืนค่า pointer value หรือ default value ถ้า nil
func SafeDerefWithDefault[T any](p *T, defaultValue T) T {
	if p == nil {
		return defaultValue
	}
	return *p
}

// IsNilOrZero ตรวจสอบว่าเป็น nil หรือ zero value
func IsNilOrZero(v interface{}) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return rv.IsNil()
	case reflect.String, reflect.Int, reflect.Int64, reflect.Float64, reflect.Bool:
		return reflect.DeepEqual(v, reflect.Zero(rv.Type()).Interface())
	}
	return false
}

// SafeMapGet อ่านค่า map อย่างปลอดภัย
func SafeMapGet[K comparable, V any](m map[K]V, key K) (V, bool) {
	if m == nil {
		var zero V
		return zero, false
	}
	val, ok := m[key]
	return val, ok
}

// SafeMapSet เขียนค่า map อย่างปลอดภัย (สร้าง map ถ้า nil)
func SafeMapSet[K comparable, V any](m map[K]V, key K, value V) map[K]V {
	if m == nil {
		m = make(map[K]V)
	}
	m[key] = value
	return m
}

// SafeCall เรียก function อย่างปลอดภัย
func SafeCall(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = &PanicError{Value: r}
		}
	}()
	if fn == nil {
		return &NilFunctionError{}
	}
	fn()
	return nil
}

// PanicError error จาก panic
type PanicError struct {
	Value interface{}
}

func (e *PanicError) Error() string {
	return "panic: " + toString(e.Value)
}

// NilFunctionError function เป็น nil
type NilFunctionError struct{}

func (e *NilFunctionError) Error() string {
	return "function is nil"
}

func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return "unknown"
}

// Chain เป็น safe chain สำหรับ pointer fields
type Chain[T any] struct {
	value *T
	err   error
}

// Of เริ่ม chain จาก pointer
func Of[T any](p *T) *Chain[T] {
	return &Chain[T]{value: p}
}

// Get รับ field โดยใช้ function selector
func (c *Chain[T]) Get(selector func(*T) interface{}) *Chain[interface{}] {
	if c.err != nil || c.value == nil {
		return &Chain[interface{}]{err: c.err}
	}
	val := selector(c.value)
	return &Chain[interface{}]{value: &val}
}

// Value คืนค่าหรือ panic ถ้ามี error
func (c *Chain[T]) Value() T {
	if c.err != nil {
		panic(c.err)
	}
	if c.value == nil {
		var zero T
		return zero
	}
	return *c.value
}

// OrZero คืนค่า zero value ถ้า nil หรือ error
func (c *Chain[T]) OrZero() T {
	if c.err != nil || c.value == nil {
		var zero T
		return zero
	}
	return *c.value
}
```

---

## วิธีใช้งาน module นี้

```go
package main

import (
	"fmt"
	"yourmodule/pkg/nil"
)

type User struct {
	Name  string
	Age   *int
	Address *struct {
		City string
	}
}

func main() {
	// 1. SafeDeref
	var age *int
	fmt.Println(nil.SafeDeref(age)) // 0

	defaultAge := 30
	fmt.Println(nil.SafeDerefWithDefault(age, defaultAge)) // 30

	// 2. IsNilOrZero
	fmt.Println(nil.IsNilOrZero(""))     // true
	fmt.Println(nil.IsNilOrZero(0))      // true
	fmt.Println(nil.IsNilOrZero([]int{})) // false (empty slice not nil)

	// 3. SafeMapGet/Set
	var m map[string]int
	val, ok := nil.SafeMapGet(m, "key")
	fmt.Println(val, ok) // 0 false

	m = nil.SafeMapSet(m, "key", 100)
	fmt.Println(m["key"]) // 100

	// 4. Safe Call
	err := nil.SafeCall(nil)
	fmt.Println(err) // function is nil

	// 5. Chain pattern
	user := &User{Name: "John"}
	result := nil.Of(user).Get(func(u *User) interface{} {
		return u.Address
	}).Get(func(a interface{}) interface{} {
		if addr, ok := a.(*struct{ City string }); ok && addr != nil {
			return addr.City
		}
		return nil
	}).OrZero()
	fmt.Println(result) // "" (zero string)
}
```

---

## การติดตั้ง

```bash
go get -u github.com/yourcompany/ccc/pkg/nil
```

หรือใช้ local module:
```go
// go.mod
module github.com/yourcompany/ccc

require (
    // no external deps for nil package
)
```

---

## การตั้งค่า configuration

ไม่มี configuration พิเศษ แต่แนะนำให้ตั้ง `nilcheck` linter:

```yaml
# .golangci.yml
linters-settings:
  nilnil:
    checked-types:
      - ptr
      - func
      - iface
  nilness:
    enabled: true
```

---

## การรวมกับ GORM

```go
package repository

import (
	"gorm.io/gorm"
	"yourproject/pkg/nil"
)

type UserRepo struct {
	db *gorm.DB
}

func (r *UserRepo) GetUser(id uint) (*User, error) {
	var user User
	err := r.db.First(&user, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil // return nil, not error
	}
	return &user, err
}

func (r *UserRepo) GetUserName(id uint) string {
	user, _ := r.GetUser(id)
	// Safe access ผ่าน pkg/nil
	name := nil.SafeDerefWithDefault(&user.Name, "Unknown")
	
	// หรือใช้ chain
	age := nil.Of(user).
		Get(func(u *User) interface{} { return u.Age }).
		OrZero().(int)
	
	fmt.Printf("Name: %s, Age: %d\n", name, age)
	return name
}

// GORM Callback ที่ปลอดภัย
func (r *UserRepo) BeforeSave(tx *gorm.DB) error {
	// ตรวจสอบ tx ไม่เป็น nil
	if tx == nil {
		return fmt.Errorf("transaction is nil")
	}
	// หรือใช้ nil.SafeCall
	return nil.SafeCall(func() {
		// logic ที่อาจ panic
	})
}
```

---

## การใช้งานจริง

### 1. JSON Unmarshal กับ optional fields

```go
type Config struct {
	Timeout  *int    `json:"timeout"`   // optional
	CacheDir *string `json:"cache_dir"` // optional
}

func LoadConfig(data []byte) Config {
	var cfg Config
	json.Unmarshal(data, &cfg)
	
	return Config{
		Timeout:  cfg.Timeout,
		CacheDir: cfg.CacheDir,
	}
}

func (c *Config) GetTimeout() int {
	return nil.SafeDerefWithDefault(c.Timeout, 30) // default 30 sec
}
```

### 2. API Response Handler

```go
type APIResponse struct {
	Data  *UserData `json:"data"`
	Error *string   `json:"error"`
}

func HandleResponse(resp *APIResponse) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}
	
	// ตรวจ error ก่อน
	if errMsg := nil.SafeDeref(resp.Error); errMsg != "" {
		return fmt.Errorf("api error: %s", errMsg)
	}
	
	// ใช้ chain สำหรับ nested data
	userID := nil.Of(resp.Data).
		Get(func(d *UserData) interface{} { return d.ID }).
		OrZero().(int)
	
	return nil
}
```

### 3. Middleware สำหรับ HTTP

```go
func NilRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("recovered from nil panic: %v", rec)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"internal server error"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
```

---

## ตารางสรุป Components

| Component | Type | Purpose |
|-----------|------|---------|
| `SafeDeref[T]` | Function | อ่านค่า pointer อย่างปลอดภัย |
| `SafeDerefWithDefault[T]` | Function | อ่านค่า pointer พร้อม default |
| `IsNilOrZero` | Function | ตรวจจับ nil หรือ zero value |
| `SafeMapGet[K,V]` | Function | อ่านค่า map ปลอดภัย |
| `SafeMapSet[K,V]` | Function | เขียนค่า map (auto-create) |
| `SafeCall` | Function | เรียก function ปลอดภัย |
| `Chain[T]` | Struct | Method chaining สำหรับ nested pointer |
| `PanicError` | Type | Error จาก panic |
| `NilFunctionError` | Type | Error จาก nil function |

---

## แบบฝึกหัดท้าย module (5 ข้อ)

### ข้อ 1
จงเขียนฟังก์ชัน `MergeMaps` ที่รับ map สองตัว (อาจเป็น nil) และคืน map ใหม่ที่รวมค่าทั้งสอง (ถ้าคีย์ซ้ำให้ใช้ค่าจาก map ที่สอง)

<details><summary>เฉลย</summary>

```go
func MergeMaps[K comparable, V any](m1, m2 map[K]V) map[K]V {
    result := make(map[K]V)
    for k, v := range m1 {
        result[k] = v
    }
    for k, v := range m2 {
        result[k] = v
    }
    return result
}
```
</details>

### ข้อ 2
สร้าง `SafeSliceGet` ที่อ่านค่าจาก slice อย่างปลอดภัย (คืน zero และ false ถ้า index out of range หรือ slice nil)

### ข้อ 3
ใช้ Chain pattern อ่าน `user.Profile.Address.ZipCode` จาก struct ที่มี pointer ทุก layer โดยปลอดภัย

### ข้อ 4
เขียน `OrElse` method สำหรับ `Chain[T]` ที่รับ fallback function และ execute เฉพาะเมื่อ value เป็น nil

### ข้อ 5
สร้าง `Must` function ที่รับ `(value T, err error)` และ panic ถ้า err ไม่เป็น nil (ใช้สำหรับ init phase)

---

## แหล่งอ้างอิง

1. [Go Spec: Nil](https://go.dev/ref/spec#Nil)
2. [Effective Go: Nil](https://go.dev/doc/effective_go#nil)
3. [GORM: Error Handling](https://gorm.io/docs/error_handling.html)
4. [Uber Go Style Guide: Nil Slices](https://github.com/uber-go/guide/blob/master/style.md#nil-slices)
5. [Dave Cheney: Why nil is not nil](https://dave.cheney.net/2017/08/09/why-is-a-nil-error-not-equal-to-nil)

---

**หมายเหตุ:** module นี้ครบถ้วนสำหรับ `pkg/nil` สำหรับระบบ ccc หากต้องการ module เพิ่มเติม (เช่น `pkg/errors`, `pkg/optional`) สามารถสร้างจาก template เดียวกันได้