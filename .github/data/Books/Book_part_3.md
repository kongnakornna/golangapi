## 📚 เล่มที่ 3: พื้นฐานภาษาและโครงสร้างข้อมูล (ตอนที่ 2)

เล่มนี้จะพาคุณเข้าสู่หัวใจสำคัญของภาษา Go นั่นคือการจัดระเบียบโค้ดด้วยแพคเกจ การสร้างชนิดข้อมูลเอง การเพิ่มพฤติกรรมด้วยเมธอด รวมถึงการเข้าใจพอยน์เตอร์และอินเทอร์เฟซ ซึ่งเป็นรากฐานของการเขียนโค้ดที่ยืดหยุ่นและนำกลับมาใช้ใหม่ได้

---

### บทที่ 11: แพคเกจและการนำเข้า (Packages & Imports)

ในโลกของ Go โค้ดทุกชิ้นจะอยู่ใน **แพคเกจ (Package)** ซึ่งเปรียบเสมือนกล่องใส่ฟังก์ชัน ตัวแปร และชนิดข้อมูลต่างๆ ที่เกี่ยวข้องกัน

**กฎสำคัญ:**
- โปรแกรม Go ทุกตัวต้องมีแพคเกจหลักชื่อ `package main` เพื่อบอกให้ Compiler รู้ว่าต้องสร้างไฟล์执行ได้
- แพคเกจอื่นๆ ใช้สำหรับสร้างไลบรารีหรือส่วนประกอบที่นำกลับมาใช้ใหม่ได้

**การสร้างแพคเกจของตัวเอง:**
สมมติเราสร้างโฟลเดอร์ `mathops` และไฟล์ `mathops.go` ข้างใน:

```go
// mathops/mathops.go
package mathops // ตั้งชื่อแพคเกจ

// ฟังก์ชันที่ขึ้นต้นด้วยตัวพิมพ์ใหญ่ จะถูกมองเห็นจากภายนอก (Exported)
func Add(a, b int) int {
    return a + b
}

// ฟังก์ชันที่ขึ้นต้นด้วยตัวพิมพ์เล็ก จะมองเห็นเฉพาะภายในแพคเกจนี้เท่านั้น (Un-exported)
func subtract(a, b int) int {
    return a - b
}
```

**การนำเข้าแพคเกจ (Import) ในไฟล์ main.go:**
```go
// main.go
package main

import (
    "fmt"
    "mathops" // ถ้าใช้ Go Modules จะอิงตามชื่อโมดูล เช่น "myproject/mathops"
)

func main() {
    result := mathops.Add(10, 5)
    fmt.Println(result) // 15

    // mathops.subtract(10, 5) // ERROR! ไม่สามารถเรียกใช้ฟังก์ชันที่ขึ้นต้นด้วยตัวพิมพ์เล็กได้
}
```

**การตั้งชื่อเล่นให้แพคเกจ (Alias):**
```go
import (
    mo "mathops" // ตั้งชื่อเล่นเป็น mo
)

func main() {
    result := mo.Add(10, 5)
}
```

**การนำเข้าเฉพาะฟังก์ชัน (ไม่แนะนำ แต่มีให้ใช้):**
```go
import . "mathops"

func main() {
    result := Add(10, 5) // เรียกใช้ได้โดยไม่ต้องนำหน้า mathops.
}
```

---

### บทที่ 12: การเริ่มต้นทำงานของแพคเกจ (Package Initialization)

เมื่อโปรแกรมเริ่มทำงาน Go จะเรียกใช้ฟังก์ชัน `init()` ในแต่ละแพคเกจก่อนที่จะเรียก `main()` เสมอ โดย `init()` จะถูกเรียกตามลำดับการนำเข้า

**ลักษณะของ `init()`:**
- ไม่มีพารามิเตอร์
- ไม่มีค่าที่คืน
- ถูกเรียกโดยอัตโนมัติ (เราไม่สามารถเรียกมันเองได้)
- ใช้สำหรับตั้งค่าเริ่มต้น เช่น เชื่อมต่อฐานข้อมูล อ่านไฟล์ Config

**ตัวอย่าง:**
```go
package database

import "fmt"

var connectionString string

// ฟังก์ชัน init จะถูกเรียกโดยอัตโนมัติเมื่อแพคเกจนี้ถูกนำเข้า
func init() {
    fmt.Println("Initializing database package...")
    connectionString = "user:pass@tcp(localhost:3306)/dbname"
    // ทำงานตั้งค่าอื่นๆ เพิ่มเติม...
}

func GetConnection() string {
    return connectionString
}
```

**ลำดับการทำงานในโปรแกรมที่มีหลายแพคเกจ:**
```
1. ตัวแปรระดับแพคเกจ (Package-level variables) ถูกกำหนดค่า
2. ฟังก์ชัน init() ของแพคเกจที่ถูกนำเข้ามาทำงาน (เรียงตามลำดับการนำเข้า)
3. ฟังก์ชัน init() ของแพคเกจ main ทำงาน
4. ฟังก์ชัน main() ทำงานเป็นลำดับสุดท้าย
```

---

### บทที่ 13: การสร้างชนิดข้อมูลใหม่ (Types)

Go อนุญาตให้เราสร้างชนิดข้อมูลใหม่จากชนิดข้อมูลที่มีอยู่แล้ว โดยใช้คีย์เวิร์ด `type` ซึ่งช่วยให้โค้ดอ่านง่ายขึ้นและมีความหมายมากกว่าใช้ชนิดดั้งเดิมโดยตรง

**รูปแบบ:**
```go
type MyInt int          // สร้างชนิดใหม่ชื่อ MyInt ที่มีพื้นฐานเป็น int
type MyString string    // สร้างชนิดใหม่ชื่อ MyString ที่มีพื้นฐานเป็น string
```

**ตัวอย่างการนำไปใช้:**
```go
package main

import "fmt"

// สร้างชนิดใหม่สำหรับอุณหภูมิ
type Celsius float64
type Fahrenheit float64

// ฟังก์ชันแปลงองศา
func CToF(c Celsius) Fahrenheit {
    return Fahrenheit(c*9/5 + 32)
}

func FToC(f Fahrenheit) Celsius {
    return Celsius((f - 32) * 5 / 9)
}

func main() {
    var tempC Celsius = 25.0
    var tempF Fahrenheit = CToF(tempC)
    
    fmt.Printf("%.2f°C = %.2f°F\n", tempC, tempF)
    // แม้ Celsius และ Fahrenheit จะมีพื้นฐานเป็น float64 เหมือนกัน
    // แต่เราสามารถนำมาเปรียบเทียบกันตรงๆ ไม่ได้ (ต้องแปลงก่อน)
}
```

**การใช้ Struct (โครงสร้าง):**
`struct` เป็นชนิดข้อมูลที่ใช้รวมข้อมูลหลายชนิดเข้าเป็นตัวเดียวกัน เหมือนตารางข้อมูล

```go
// สร้างชนิดข้อมูลใหม่ชื่อ Person
type Person struct {
    Name    string
    Age     int
    Address string
}

func main() {
    // วิธีสร้าง struct
    p1 := Person{"Somchai", 30, "Bangkok"}          // วิธีที่ 1: เรียงตามลำดับ
    p2 := Person{Name: "Somsri", Age: 25}           // วิธีที่ 2: ระบุชื่อฟิลด์ (ไม่ต้องเรียง)
    p3 := Person{}                                   // วิธีที่ 3: ค่าว่าง (ชื่อ = "", อายุ = 0)

    p1.Age = 31                                      // เปลี่ยนค่า
    fmt.Println(p1.Name, p1.Age)
}
```

---

### บทที่ 14: เมธอด (Methods)

เมธอดคือฟังก์ชันที่ผูกติดกับชนิดข้อมูลใดชนิดหนึ่งโดยเฉพาะ (คล้ายกับคลาสในภาษา OOP อื่นๆ) แต่ใน Go เราเพิ่มเมธอดให้กับ **ชนิดใดก็ได้** ไม่ใช่แค่ struct

**รูปแบบ:**
```go
func (ตัวรับ Receiver) ชื่อเมธอด(พารามิเตอร์) ชนิดที่คืน {
    // คำสั่ง
}
```

**ตัวอย่างกับ Struct:**
```go
package main

import "fmt"

type Rectangle struct {
    Width  float64
    Height float64
}

// เมธอดที่ผูกกับ Rectangle (รับค่าเป็นสำเนา - Value Receiver)
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// เมธอดที่เปลี่ยนค่าของ struct (ต้องใช้ Pointer Receiver)
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    area := rect.Area()
    fmt.Println("พื้นที่:", area) // 50

    rect.Scale(2)
    fmt.Println("หลังจากขยาย:", rect.Width, rect.Height) // 20, 10
}
```

**เมธอดกับชนิดข้อมูลอื่นๆ (ไม่ใช่ struct):**
```go
type MyInt int

// เมธอดตรวจสอบว่าเป็นเลขคู่หรือไม่
func (m MyInt) IsEven() bool {
    return int(m)%2 == 0
}

func main() {
    num := MyInt(10)
    fmt.Println(num.IsEven()) // true
}
```

**ความแตกต่างระหว่าง Value Receiver กับ Pointer Receiver:**
- **Value Receiver (`func (r Rectangle)`)** : รับสำเนา เปลี่ยนค่าในเมธอดจะไม่กระทบตัวต้นฉบับ
- **Pointer Receiver (`func (r *Rectangle)`)** : รับ Address ของตัวแปร เปลี่ยนค่าจะกระทบตัวต้นฉบับ เหมาะกับการแก้ไขข้อมูล

---

### บทที่ 15: พอยน์เตอร์ (Pointer)

พอยน์เตอร์คือตัวแปรที่เก็บ **ที่อยู่หน่วยความจำ (Memory Address)** ของตัวแปรอื่น แทนที่จะเก็บค่าตรงๆ

**ทำไมต้องใช้พอยน์เตอร์?**
- เพื่อให้ฟังก์ชันสามารถแก้ไขค่าของตัวแปรต้นฉบับได้ (Pass by Reference)
- เพื่อประหยัดหน่วยความจำ เมื่อส่ง Struct ตัวใหญ่ๆ แทนที่จะคัดลอกทั้งก้อน ก็แค่ส่งพอยน์เตอร์ (8 ไบต์) ไปเลย

**รูปแบบและการใช้งาน:**
```go
package main

import "fmt"

func main() {
    x := 10
    var ptr *int = &x // ptr ชี้ไปที่ x (เก็บ Address ของ x)
    
    fmt.Println("ค่าของ x:", x)           // 10
    fmt.Println("Address ของ x:", &x)     // 0xc00001a0b0
    fmt.Println("ค่าใน ptr:", ptr)        // 0xc00001a0b0 (เหมือน Address ของ x)
    fmt.Println("ค่าที่ ptr ชี้ไป:", *ptr) // 10 (Dereferencing)

    // เปลี่ยนค่าผ่านพอยน์เตอร์
    *ptr = 20
    fmt.Println("ค่า x ใหม่:", x) // 20
}
```

**พอยน์เตอร์ในพารามิเตอร์ของฟังก์ชัน:**
```go
func updateValue(val *int) {
    *val = 100 // เปลี่ยนค่าตัวแปรต้นฉบับ
}

func main() {
    num := 5
    updateValue(&num) // ส่ง Address ไป
    fmt.Println(num) // 100
}
```

**ค่าว่างของพอยน์เตอร์: `nil`**
```go
var p *int
fmt.Println(p) // <nil>
if p == nil {
    fmt.Println("พอยน์เตอร์นี้ยังไม่ได้ชี้ไปที่ไหน")
}
```

**ข้อควรระวัง:** ไม่ควรใช้พอยน์เตอร์กับชนิดข้อมูลเล็กๆ (int, bool, string) เพราะเสียเวลาในการเข้าถึง Address มากกว่าเปล่าๆ ใช้กับ struct หรือ slice ขนาดใหญ่จะคุ้มกว่า

---

### บทที่ 16: อินเทอร์เฟซ (Interfaces)

อินเทอร์เฟซคือข้อตกลง (Contract) ที่บอกว่า **"ถ้าชนิดข้อมูลนี้มีเมธอดเหล่านี้ ฉันจะถือว่ามันเป็นประเภทนี้"** ช่วยให้เราเขียนโค้ดที่ยืดหยุ่นและสามารถเปลี่ยนการทำงานได้โดยไม่ต้องแก้ไขตรรกะหลัก

**รูปแบบ:**
```go
type ชื่อInterface interface {
    ชื่อเมธอด1(พารามิเตอร์) ชนิดที่คืน
    ชื่อเมธอด2(พารามิเตอร์) ชนิดที่คืน
}
```

**ตัวอย่างคลาสสิก:**
```go
package main

import "fmt"

// ประกาศอินเทอร์เฟซ
type Animal interface {
    Speak() string
    Move() string
}

// ชนิดข้อมูล Dog
type Dog struct{}

func (d Dog) Speak() string {
    return "โฮ่ง โฮ่ง!"
}
func (d Dog) Move() string {
    return "วิ่ง 4 ขา"
}

// ชนิดข้อมูล Bird
type Bird struct{}

func (b Bird) Speak() string {
    return "จาบ จาบ!"
}
func (b Bird) Move() string {
    return "บิน"
}

// ฟังก์ชันที่รับอินเทอร์เฟซ Animal (ใช้กับได้ทุกชนิดที่ทำตามข้อตกลง)
func DescribeAnimal(a Animal) {
    fmt.Println("เสียง:", a.Speak())
    fmt.Println("การเคลื่อนไหว:", a.Move())
}

func main() {
    dog := Dog{}
    bird := Bird{}

    DescribeAnimal(dog)   // ใช้ได้!
    DescribeAnimal(bird)  // ใช้ได้!
}
```

**อินเทอร์เฟซแบบ Empty Interface (`interface{}` หรือ `any`):**
อินเทอร์เฟซที่ไม่มีเมธอดใดๆ ทุกชนิดข้อมูลใน Go จะเป็นไปตามอินเทอร์เฟซนี้โดยอัตโนมัติ (เหมือน Object ใน Java)

```go
var anything interface{} // หรือ var anything any
anything = 10
anything = "Hello"
anything = Dog{}

fmt.Println(anything) // เก็บอะไรก็ได้
```

**Type Assertion (การตรวจสอบชนิดข้อมูลจริง):**
เมื่อเรามีตัวแปรประเภทอินเทอร์เฟซ เราสามารถตรวจสอบได้ว่าข้างในมันคือชนิดอะไร

```go
func main() {
    var a Animal = Dog{}
    
    // ตรวจสอบว่าเป็น Dog หรือไม่
    dog, ok := a.(Dog)
    if ok {
        fmt.Println("มันคือ Dog:", dog.Speak())
    }

    // ตรวจสอบหลายชนิดด้วย switch
    switch v := a.(type) {
    case Dog:
        fmt.Println("เป็น Dog", v)
    case Bird:
        fmt.Println("เป็น Bird", v)
    default:
        fmt.Println("ชนิดอื่น ๆ")
    }
}
```

**ประโยชน์ของอินเทอร์เฟซ:**
- เขียนฟังก์ชันที่รับได้หลายชนิด (Polymorphism)
- แยกการทำงานออกจากการนำไปใช้ (Decoupling) เหมาะกับการเขียน Test
- โค้ดยืดหยุ่น เปลี่ยนแปลงได้ง่าย

---

### 📌 สรุปเล่มที่ 3
ในเล่มนี้คุณได้เรียนรู้หัวใจสำคัญของการจัดโครงสร้างโปรแกรม:
✅ การสร้างและนำเข้าแพคเกจ (Public/Private ด้วยตัวพิมพ์ใหญ่-เล็ก)
✅ ฟังก์ชัน `init()` สำหรับการตั้งค่าเริ่มต้นอัตโนมัติ
✅ การสร้างชนิดข้อมูลใหม่ด้วย `type` และ Struct
✅ การเพิ่มพฤติกรรมด้วยเมธอด (Value vs Pointer Receiver)
✅ การใช้พอยน์เตอร์เพื่อส่งต่อข้อมูลแบบ Reference
✅ การเขียนโค้ดที่ยืดหยุ่นด้วยอินเทอร์เฟซและ Type Assertion

---
 