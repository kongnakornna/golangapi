## 📚 เล่มที่ 4: การจัดการโปรเจกต์และโครงสร้างข้อมูลขั้นสูง

เล่มนี้จะพาคุณเข้าสู่โลกของการจัดการโปรเจกต์สมัยใหม่ด้วย Go Modules การเขียนทดสอบเพื่อการันตีคุณภาพโค้ด และทำความรู้จักโครงสร้างข้อมูลที่ใช้กันทุกวันในชีวิตจริง ไม่ว่าจะเป็น Array, Slice, Map รวมถึงการจัดการข้อผิดพลาดอย่างมืออาชีพ

---

### บทที่ 17: Go Modules - การจัดการโปรเจกต์สมัยใหม่

ก่อนยุค Go Modules (สมัยใช้ GOPATH) การเขียน Go นั้นค่อนข้างยุ่งยากเพราะโค้ดทั้งหมดต้องอยู่ใต้โฟลเดอร์ `src` เดียวกัน Go Modules เกิดขึ้นในเวอร์ชัน 1.11 เพื่อแก้ปัญหานี้ ทำให้เราสามารถสร้างโปรเจกต์ได้ทุกที่ในเครื่อง

**การเริ่มต้นโปรเจกต์ใหม่ด้วย Go Modules:**
```bash
# สร้างโฟลเดอร์โปรเจกต์
mkdir myproject
cd myproject

# เริ่มต้นโมดูล (module) ตั้งชื่อตามเส้นทางที่ใช้ดึงจากเว็บ
go mod init github.com/username/myproject
# หรือถ้าเป็นโปรเจกต์ส่วนตัวในเครื่อง ตั้งชื่อง่ายๆ เช่น go mod init myproject
```

เมื่อรันคำสั่งนี้ จะเกิดไฟล์ `go.mod` ขึ้นมา ซึ่งเป็นหัวใจของระบบ Modules มีลักษณะดังนี้:
```go
module github.com/username/myproject

go 1.21 // เวอร์ชันของ Go ที่ใช้
```

**การเพิ่ม Dependency (แพคเกจภายนอก):**
เวลาเรานำเข้าแพคเกจจากภายนอกในโค้ด เช่น:
```go
import "github.com/gin-gonic/gin"
```
แล้วรัน `go mod tidy` คำสั่งนี้จะ:
1. ดาวน์โหลดแพคเกจดังกล่าวมาไว้ในเครื่อง
2. เพิ่มข้อมูลลงใน `go.mod` พร้อมเวอร์ชันล่าสุด
3. สร้างไฟล์ `go.sum` เก็บ Checksum เพื่อยืนยันความถูกต้องของไฟล์ที่ดาวน์โหลด

**คำสั่งสำคัญที่ควรรู้:**
- `go mod init <ชื่อ>` : สร้างโมดูลใหม่
- `go mod tidy` : ดาวน์โหลด dependencies ที่ขาดหาย และลบตัวที่ไม่ใช้ออก
- `go mod vendor` : คัดลอกโค้ด dependencies ทั้งหมดมาไว้ในโฟลเดอร์ `vendor/` (ใช้ในกรณีที่ต้องอัปโหลดทุกอย่างไปบน Server ปิด)
- `go get <package>@<version>` : ดาวน์โหลดหรืออัปเดตแพคเกจเป็นเวอร์ชันเฉพาะ (เช่น `go get github.com/gin-gonic/gin@v1.9.0`)

---

### บทที่ 18: Go Module Proxies

โดยปกติ `go get` จะไปดึงโค้ดจาก GitHub หรือ Git โดยตรง แต่ในบางประเทศหรือเครือข่ายองค์กรที่เชื่อมต่อช้า Google จึงสร้าง **Go Proxy** ขึ้นมาเพื่อทำหน้าที่เป็นตัว Cache กลาง ช่วยให้ดาวน์โหลดเร็วขึ้นและป้องกันการหายไปของโค้ด (ถ้าเจ้าของลบ Repo ทิ้ง)

**Proxy หลักที่ Go ใช้คือ:** `https://proxy.golang.org`

**การตั้งค่า Proxy:**
```bash
# ดูค่าปัจจุบัน
go env GOPROXY

# ตั้งให้ใช้ Proxy ของ Google (ค่าเริ่มต้น)
go env -w GOPROXY=https://proxy.golang.org,direct

# ถ้าอยู่ในองค์กรที่ใช้ Private Repo ให้ใช้ proxy ของตัวเอง
go env -w GOPROXY=https://mycorp-proxy.com,https://proxy.golang.org,direct
```
- `direct` หมายถึง ถ้า Proxy หาไม่เจอ ให้ไปดึงจาก Git โดยตรง
- ถ้าต้องการปิด Proxy ใช้คำสั่ง: `go env -w GOPROXY=direct`

**Private Repository:** ถ้ามี Git Repo ส่วนตัวที่ต้องใช้เป็น Dependency ให้ตั้งตัวแปร `GOPRIVATE` เพื่อบอก Go ว่าอย่าไปหาใน Proxy:
```bash
go env -w GOPRIVATE=github.com/mycompany/*
```

---

### บทที่ 19: การทดสอบหน่วย (Unit Tests)

การทดสอบเป็นสิ่งสำคัญที่ช่วยให้มั่นใจว่าโค้ดของคุณทำงานถูกต้อง Go มีเครื่องมือทดสอบในตัวที่ใช้งานง่าย ไม่ต้องติดตั้งเพิ่ม

**กฎการตั้งชื่อไฟล์ทดสอบ:**
- ไฟล์ต้องลงท้ายด้วย `_test.go` เช่น `math_test.go`
- ฟังก์ชันทดสอบต้องขึ้นต้นด้วยคำว่า `Test` และรับพารามิเตอร์ `*testing.T`

**ตัวอย่างการเขียนทดสอบ:**
สมมติเรามีไฟล์ `math.go`:
```go
package math

func Add(a, b int) int {
    return a + b
}

func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("หารด้วย 0 ไม่ได้")
    }
    return a / b, nil
}
```
ให้สร้างไฟล์ `math_test.go`:
```go
package math

import "testing"

// ทดสอบฟังก์ชัน Add
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    expected := 5

    if result != expected {
        t.Errorf("Add(2, 3) ได้ %d แต่คาดหวัง %d", result, expected)
    }
}

// ทดสอบฟังก์ชัน Divide (รวมถึงกรณี Error)
func TestDivide(t *testing.T) {
    result, err := Divide(10, 2)
    if err != nil {
        t.Errorf("ไม่ควรมี Error: %v", err)
    }
    if result != 5 {
        t.Errorf("Divide(10,2) ได้ %d แต่คาดหวัง 5", result)
    }

    // ทดสอบกรณีหารด้วย 0
    _, err = Divide(10, 0)
    if err == nil {
        t.Error("ควรมี Error แต่ไม่เกิด")
    }
}
```

**Table-Driven Tests (รูปแบบที่นิยมใน Go):** ใช้ทดสอบหลายเคสด้วยกัน
```go
func TestAddTable(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"บวกเลขบวก", 1, 2, 3},
        {"บวกเลขลบ", -1, -2, -3},
        {"บวกกับศูนย์", 5, 0, 5},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("ได้ %d แต่คาดหวัง %d", result, tt.expected)
            }
        })
    }
}
```
**รันการทดสอบ:**
- `go test` : รันเทสทั้งหมดในโฟลเดอร์ปัจจุบัน
- `go test -v` : แสดงผลแบบละเอียด
- `go test -run TestAdd` : รันเฉพาะฟังก์ชันที่ชื่อขึ้นต้นด้วย TestAdd

---

### บทที่ 20: อาเรย์ (Arrays)

อาเรย์คือโครงสร้างข้อมูลที่เก็บชุดข้อมูลชนิดเดียวกัน โดยมีขนาด **คงที่ (Fixed Size)** ซึ่งกำหนดไว้ตั้งแต่ตอนประกาศ

**ลักษณะเด่น:**
- ขนาดตายตัว เปลี่ยนแปลงไม่ได้
- เก็บข้อมูลเรียงต่อเนื่องในหน่วยความจำ
- การเข้าถึงข้อมูลรวดเร็วมาก (ใช้ดัชนี Index)

**การประกาศและการใช้งาน:**
```go
package main

import "fmt"

func main() {
    // วิธีที่ 1: ประกาศแล้วกำหนดขนาด
    var numbers [5]int        // อาเรย์ขนาด 5 เก็บ int (ค่าเริ่มต้นเป็น 0 ทั้งหมด)
    numbers[0] = 10
    numbers[1] = 20
    fmt.Println(numbers)      // [10 20 0 0 0]

    // วิธีที่ 2: ประกาศพร้อมกำหนดค่าเริ่มต้น (ใช้เครื่องหมาย ... เพื่อให้ Compiler นับขนาดให้)
    fruits := [...]string{"แอปเปิ้ล", "กล้วย", "ส้ม"}
    fmt.Println("ขนาด:", len(fruits)) // 3

    // วิธีที่ 3: กำหนดค่าบางตำแหน่ง
    scores := [5]int{0: 100, 4: 50} // ตำแหน่ง 0 = 100, ตำแหน่ง 4 = 50

    // การเข้าถึงข้อมูล
    fmt.Println(fruits[1]) // กล้วย

    // การวนลูป
    for i := 0; i < len(scores); i++ {
        fmt.Printf("scores[%d] = %d\n", i, scores[i])
    }

    // การใช้ Range
    for index, value := range fruits {
        fmt.Printf("ดัชนี %d = %s\n", index, value)
    }
}
```

**ข้อควรจำ:** อาเรย์จะถูก **คัดลอกทั้งก้อน** เมื่อส่งให้ฟังก์ชัน (Pass by Value) ถ้าต้องการให้ฟังก์ชันแก้ไขอาเรย์ต้นฉบับ ต้องใช้พอยน์เตอร์
```go
func updateArray(arr *[3]int) {
    (*arr)[0] = 999
}
```

---

### บทที่ 21: สไลซ์ (Slices)

สไลซ์คือส่วนขยายของอาเรย์ที่มีขนาด **ยืดหยุ่น เปลี่ยนแปลงได้** (Dynamic Array) นี่คือโครงสร้างข้อมูลที่ใช้กันมากที่สุดใน Go

**สไลซ์ประกอบด้วย 3 ส่วน:**
1. พอยน์เตอร์ชี้ไปที่อาเรย์ต้นแบบ (Underlying Array)
2. ความยาว (Length) : จำนวนสมาชิกในสไลซ์
3. ความจุ (Capacity) : จำนวนสมาชิกสูงสุดที่อาเรย์ต้นแบบรองรับ โดยไม่ต้องขยาย

**การสร้างสไลซ์:**
```go
package main

import "fmt"

func main() {
    // วิธีที่ 1: ใช้ฟังก์ชัน make
    slice1 := make([]int, 5)    // length=5, capacity=5
    slice2 := make([]int, 3, 5) // length=3, capacity=5

    // วิธีที่ 2: ประกาศแบบ literal (เหมือน Array แต่ไม่ใส่ขนาด)
    fruits := []string{"แอปเปิ้ล", "กล้วย", "ส้ม"}

    // วิธีที่ 3: ตัดมาจากอาเรย์ (Slicing)
    arr := [5]int{1, 2, 3, 4, 5}
    slice3 := arr[1:4] // ได้ [2 3 4] (ดัชนี 1 ถึง 3)

    // การเพิ่มสมาชิก (Append)
    numbers := []int{1, 2, 3}
    numbers = append(numbers, 4, 5) // [1 2 3 4 5]
    numbers = append(numbers, []int{6, 7}...) // ใช้ ... เพื่อแตก slice

    // การตรวจสอบ length และ capacity
    fmt.Printf("len=%d cap=%d %v\n", len(numbers), cap(numbers), numbers)
}
```

**การทำงานของ Append และการขยาย Capacity:**
เมื่อ append แล้วพื้นที่ไม่พอ Go จะสร้างอาเรย์ใหม่ที่มีขนาดเป็น 2 เท่า (โดยประมาณ) และคัดลอกข้อมูลเก่าไปไว้ ทำให้เกิดการจัดสรรหน่วยความจำใหม่
```go
s := []int{1, 2, 3}
fmt.Println(cap(s)) // 3
s = append(s, 4)
fmt.Println(cap(s)) // 6 (ขยายเป็น 2 เท่า)
```

**การรวม (Copy) สไลซ์:**
```go
src := []int{1, 2, 3}
dst := make([]int, len(src))
copy(dst, src) // คัดลอก src ไป dst
```

---

### บทที่ 22: แมพ (Maps)

Map คือโครงสร้างข้อมูลแบบคู่คีย์-ค่า (Key-Value) เหมือนพจนานุกรม ใช้ค้นหาข้อมูลได้รวดเร็วด้วยคีย์ (Key)

**การประกาศและการใช้งาน:**
```go
package main

import "fmt"

func main() {
    // วิธีที่ 1: ใช้ make
    scores := make(map[string]int) // คีย์เป็น string, ค่าเป็น int
    scores["Somchai"] = 85
    scores["Somsri"] = 92

    // วิธีที่ 2: ประกาศพร้อมค่าเริ่มต้น (Literal)
    fruits := map[string]string{
        "A": "แอปเปิ้ล",
        "B": "กล้วย",
        "C": "ส้ม",
    }

    // การเข้าถึงข้อมูล
    fmt.Println(scores["Somchai"]) // 85

    // การตรวจสอบว่ามีคีย์นี้หรือไม่ (Comma-ok idiom)
    value, exists := scores["Somsri"]
    if exists {
        fmt.Println("คะแนน Somsri คือ", value)
    } else {
        fmt.Println("ไม่พบชื่อนี้")
    }

    // การลบข้อมูล
    delete(fruits, "B") // ลบคีย์ "B" ออก

    // การวนลูป
    for key, value := range scores {
        fmt.Printf("%s ได้ %d คะแนน\n", key, value)
    }
}
```

**ข้อควรระวัง:** Map ใน Go เป็นแบบไม่เรียงลำดับ ถ้าต้องการเรียงลำดับต้องสร้าง Slice ของคีย์ขึ้นมา แล้ว Sort
```go
import "sort"

keys := make([]string, 0, len(scores))
for k := range scores {
    keys = append(keys, k)
}
sort.Strings(keys)

for _, k := range keys {
    fmt.Println(k, scores[k])
}
```

---

### บทที่ 23: การจัดการข้อผิดพลาด (Errors)

ใน Go ไม่มีการใช้ Exception (try-catch) แบบภาษาอื่นๆ แต่ใช้การคืนค่าผิดพลาดเป็นค่าที่สองของฟังก์ชัน (`error`) ซึ่งเป็นอินเทอร์เฟซในตัวของ Go

**รูปแบบมาตรฐาน:**
```go
func doSomething() (result int, err error) {
    if เกิดข้อผิดพลาด {
        return 0, errors.New("ข้อความผิดพลาด")
    }
    return result, nil // nil = ไม่มีข้อผิดพลาด
}
```

**ตัวอย่างการใช้งาน:**
```go
package main

import (
    "errors"
    "fmt"
)

// ฟังก์ชันที่อาจมี Error
func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("ไม่สามารถหารด้วยศูนย์ได้")
    }
    return a / b, nil
}

func main() {
    result, err := Divide(10, 0)
    if err != nil {
        // จัดการ Error
        fmt.Println("เกิดข้อผิดพลาด:", err)
        return
    }
    fmt.Println("ผลลัพธ์:", result)
}
```

**การสร้าง Error แบบละเอียดด้วย `fmt.Errorf`:**
```go
func ReadFile(filename string) error {
    if filename == "" {
        return fmt.Errorf("ชื่อไฟล์ว่างเปล่า: %s", filename)
    }
    // ...
    return nil
}
```

**การประกาศ Error แบบกำหนดเอง (Custom Error):**
```go
// สร้างชนิด Error เอง
type ValidationError struct {
    Field string
    Value interface{}
    Msg   string
}

// ทำให้เป็นไปตามอินเทอร์เฟซ error
func (e ValidationError) Error() string {
    return fmt.Sprintf("ฟิลด์ %s มีค่า %v ผิดพลาด: %s", e.Field, e.Value, e.Msg)
}

// ใช้งาน
func ValidateAge(age int) error {
    if age < 0 {
        return ValidationError{Field: "age", Value: age, Msg: "อายุต้องไม่ติดลบ"}
    }
    return nil
}
```

**การใช้ `errors.Is` และ `errors.As` (Go 1.13+):** ใช้ตรวจสอบประเภท Error
```go
import "errors"

var ErrNotFound = errors.New("ไม่พบข้อมูล")

func FindUser(id int) error {
    return fmt.Errorf("ค้นหา id %d: %w", id, ErrNotFound) // %w = wrap error
}

func main() {
    err := FindUser(100)
    if errors.Is(err, ErrNotFound) {
        fmt.Println("ไม่พบผู้ใช้นี้")
    }
}
```

**การจัดการ Panic และ Recover (ใช้เฉพาะกรณีร้ายแรง):**
- `panic` : ทำให้โปรแกรมหยุดทำงานทันที (เหมือน Exception รุนแรง)
- `recover` : ใช้ใน `defer` เพื่อดักจับ Panic แล้วกลับมาให้โปรแกรมทำงานต่อได้
```go
defer func() {
    if r := recover(); r != nil {
        fmt.Println("โปรแกรมถูกบันทึกไว้:", r)
    }
}()
panic("เกิดปัญหาสำคัญ!")

### 📌 สรุปเล่มที่ 4
ในเล่มนี้คุณได้เรียนรู้เครื่องมือและโครงสร้างข้อมูลที่ใช้ในชีวิตจริง:
✅ การจัดการโปรเจกต์ด้วย Go Modules (`go.mod`, `go.sum`)
✅ การใช้ Go Proxy เพื่อเพิ่มความเร็วในการดาวน์โหลด
✅ การเขียน Unit Tests และ Table-Driven Tests
✅ Array (ขนาดคงที่) และ Slices (ขนาดยืดหยุ่น) พร้อมการ Append
✅ Maps สำหรับการจัดเก็บข้อมูลแบบ Key-Value
✅ การจัดการ Error อย่างเป็นระบบและเป็นมืออาชีพ

---
