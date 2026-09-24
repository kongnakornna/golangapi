# 📘 เล่มที่ 2: โครงสร้างข้อมูลขั้นสูงและการจัดการ JSON

## เริ่มทำเนื้อหาทั้งเล่ม


# เรียน Go อย่างมืออาชีพ
## เล่มที่ 2: โครงสร้างข้อมูลขั้นสูงและการจัดการ JSON

---

# สารบัญ เล่มที่ 2

### บทที่ 1: Structs
- 1.1 Introduction to Structs
- 1.2 Struct กับ Pointer
- 1.3 Struct Embedding และ Anonymous Field
- 1.4 Struct และ Recursive
- 1.5 Embedded Struct — คอนเซปและการประยุกต์ใช้

### บทที่ 2: Methods
- 2.1 The Need of Method
- 2.2 Method Declaration
- 2.3 Method by Value vs by Reference
- 2.4 Method Best Practice และ Syntax Sugar
- 2.5 Method with Nil Receiver
- 2.6 Method Values

### บทที่ 3: Interfaces
- 3.1 Interfaces — คอนเซปพื้นฐาน
- 3.2 Interface Satisfaction — การทำให้ Interface เป็นจริง
- 3.3 Interface Value
- 3.4 Empty Interface
- 3.5 Type Assertion — ตอนที่ 1 และ 2
- 3.6 Type Assertion without Panic
- 3.7 Type Switches (Discriminated Union)
- 3.8 flag.Interface และการประยุกต์ใช้
- 3.9 sort.Interface — การเรียงลำดับด้วย Interface
- 3.10 error Interface
- 3.11 Interface ที่มี Methods ซ้ำกัน

### บทที่ 4: JSON
- 4.1 JSON Introduction
- 4.2 JSON Unmarshal — แปลง JSON เป็น Struct
- 4.3 JSON Marshal — แปลง Struct เป็น JSON
- 4.4 ใช้งาน JSON Unmarshal กับ HTTP GET
- 4.5 JSON Encoder/Decoder
- 4.6 Struct Tags สำหรับ JSON

### บทที่ 5: Templates
- 5.1 Text Template Introduction
- 5.2 Text Template — Actions และฟีเจอร์เพิ่มเติม
- 5.3 HTML Template

### บทที่ 6: Reflection
- 6.1 คอนเซปของ Reflection
- 6.2 Report Application — โครงสร้างโปรเจกต์
- 6.3 การเข้าถึง Struct Name
- 6.4 การเข้าถึง Struct Value
- 6.5 การเข้าถึง Struct Tag
- 6.6 Multiple Tags
- 6.7 การแก้ไข Struct Value
- 6.8 Reflection สรุป

### บทที่ 7: การจัดการ Error ขั้นสูง
- 7.1 Project Overview — Validation Library
- 7.2 Initial Project with Simple Length Validation
- 7.3 Add More Validations Chain
- 7.4 Refactored into Package
- 7.5 First Look at `errors.Is` และ `Unwrap` Interface
- 7.6 Multi-layer Wrapping and Unwrapping Error
- 7.7 Tough Life before `errors.As`
- 7.8 Type Assertion with `errors.As`
- 7.9 Refactored Custom `As` Function
- 7.10 Add More Error Types

---

# บทที่ 1: Structs

---

## 1.1 Introduction to Structs

**Struct** คือ collection ของ fields ที่สามารถมีชนิดข้อมูลต่างกันได้ ใช้สำหรับจัดกลุ่มข้อมูลที่เกี่ยวข้องกัน

### รูปแบบ

```go
type StructName struct {
    Field1 Type1
    Field2 Type2
    Field3 Type3
}
```

### ตัวอย่าง

```go
package main

import "fmt"

// ประกาศ struct
type Person struct {
    Name    string
    Age     int
    Height  float64
    IsActive bool
}

func main() {
    // สร้าง instance ของ struct (แบบที่ 1)
    var p1 Person
    p1.Name = "สมชาย"
    p1.Age = 30
    p1.Height = 175.5
    p1.IsActive = true
    fmt.Println(p1)

    // แบบที่ 2: ระบุ fields
    p2 := Person{
        Name:    "สมหญิง",
        Age:     25,
        Height:  165.0,
        IsActive: true,
    }
    fmt.Println(p2)

    // แบบที่ 3: ระบุเฉพาะค่าตามลำดับ fields
    p3 := Person{"วิชัย", 40, 180.0, false}
    fmt.Println(p3)

    // เข้าถึง field
    fmt.Println(p2.Name)  // สมหญิง

    // เปลี่ยนค่า
    p2.Age = 26
    fmt.Println(p2.Age)  // 26
}
```

### Struct ที่ซ้อน struct

```go
package main

import "fmt"

type Address struct {
    Street string
    City   string
    Zip    string
}

type Employee struct {
    Name    string
    Age     int
    Address Address  // struct ซ้อน struct
}

func main() {
    emp := Employee{
        Name: "สมชาย",
        Age:  30,
        Address: Address{
            Street: "123 ถนนสุขุมวิท",
            City:   "กรุงเทพฯ",
            Zip:    "10110",
        },
    }

    fmt.Println(emp.Name)                  // สมชาย
    fmt.Println(emp.Address.City)          // กรุงเทพฯ
}
```

---

## 1.2 Struct กับ Pointer

การใช้ pointer กับ struct ช่วยให้สามารถแก้ไข struct โดยไม่ต้อง copy ทั้ง struct

### ตัวอย่าง

```go
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

// รับ pointer เพื่อแก้ไข struct
func updateAge(p *Person, newAge int) {
    p.Age = newAge
}

// รับ pointer เพื่อเปลี่ยนชื่อ
func rename(p *Person, newName string) {
    p.Name = newName
}

func main() {
    person := Person{Name: "สมชาย", Age: 30}

    fmt.Println("ก่อน:", person)  // {สมชาย 30}

    // ส่ง pointer
    updateAge(&person, 35)
    rename(&person, "สมชาย ใจดี")

    fmt.Println("หลัง:", person)  // {สมชาย ใจดี 35}
}
```

### การสร้าง struct ด้วย new

```go
func main() {
    // new คืนค่า pointer ไปที่ zero value
    p := new(Person)
    fmt.Println(p)   // &{ 0}
    fmt.Println(*p)  // { 0}

    p.Name = "สมชาย"
    p.Age = 30
    fmt.Println(p)   // &{สมชาย 30}
}
```

### การใช้ & กับ struct literal

```go
func main() {
    // &struct literal คืน pointer ทันที
    p := &Person{Name: "สมชาย", Age: 30}
    fmt.Println(p)   // &{สมชาย 30}

    // ไม่ต้องใช้ *p เมื่อเข้าถึง field (Go ทำ auto-dereference)
    p.Age = 35
    fmt.Println(p.Age)  // 35
}
```

---

## 1.3 Struct Embedding และ Anonymous Field

### Anonymous Field (ฟิลด์ไม่มีชื่อ)

```go
package main

import "fmt"

type Person struct {
    string  // anonymous field (ชนิด string)
    int     // anonymous field (ชนิด int)
}

func main() {
    p := Person{"สมชาย", 30}
    fmt.Println(p.string)  // สมชาย
    fmt.Println(p.int)     // 30

    // ระบุชื่อฟิลด์
    p2 := Person{string: "สมหญิง", int: 25}
    fmt.Println(p2)  // {สมหญิง 25}
}
```

### Struct Embedding (ฝัง struct)

```go
package main

import "fmt"

type Address struct {
    Street string
    City   string
}

type Person struct {
    Name    string
    Age     int
    Address  // ฝัง struct (no field name)
}

func main() {
    p := Person{
        Name: "สมชาย",
        Age:  30,
        Address: Address{
            Street: "123 ถนนสุขุมวิท",
            City:   "กรุงเทพฯ",
        },
    }

    // เข้าถึง field ของ Address โดยตรง (เหมือนเป็น field ของ Person)
    fmt.Println(p.Street)  // 123 ถนนสุขุมวิท
    fmt.Println(p.City)    // กรุงเทพฯ

    // ยังสามารถเข้าถึงผ่าน Address ได้
    fmt.Println(p.Address.Street)  // 123 ถนนสุขุมวิท
}
```

### Method Inheritance ผ่าน Embedding

```go
package main

import "fmt"

type Animal struct {
    Name string
}

func (a Animal) Speak() string {
    return "เสียงสัตว์"
}

type Dog struct {
    Animal  // ฝัง Animal
    Breed string
}

func main() {
    dog := Dog{
        Animal: Animal{Name: "ทองทอง"},
        Breed:  "โกลเด้น",
    }

    // Dog มี method Speak จาก Animal
    fmt.Println(dog.Speak())  // เสียงสัตว์

    // สามารถ override method ได้
    // (จะได้เรียนในหัวข้อ Methods)
}
```

---

## 1.4 Struct และ Recursive

Struct สามารถอ้างอิงตัวเองได้ (recursive struct) เหมาะสำหรับโครงสร้างข้อมูลแบบ Tree

### ตัวอย่าง

```go
package main

import "fmt"

// Node ใน Binary Tree
type TreeNode struct {
    Value int
    Left  *TreeNode  // ชี้ไปที่ Node อื่น (ใช้ pointer)
    Right *TreeNode
}

// Linked List Node
type ListNode struct {
    Value int
    Next  *ListNode
}

func main() {
    // สร้าง Binary Tree
    //       1
    //      / \
    //     2   3
    //    / \
    //   4   5
    root := &TreeNode{
        Value: 1,
        Left: &TreeNode{
            Value: 2,
            Left:  &TreeNode{Value: 4},
            Right: &TreeNode{Value: 5},
        },
        Right: &TreeNode{
            Value: 3,
        },
    }

    // พิมพ์แบบ inorder
    fmt.Print("Inorder: ")
    printInorder(root)
    fmt.Println()
}

func printInorder(node *TreeNode) {
    if node == nil {
        return
    }
    printInorder(node.Left)
    fmt.Printf("%d ", node.Value)
    printInorder(node.Right)
}
```

### แบบฝึกหัด Recursive Struct

```go
package main

import "fmt"

// Recursive struct สำหรับ แฟมิลีทรี
type FamilyMember struct {
    Name     string
    Children []*FamilyMember
}

func main() {
    // สร้าง Family Tree
    family := &FamilyMember{
        Name: "ปู่",
        Children: []*FamilyMember{
            {
                Name: "พ่อ",
                Children: []*FamilyMember{
                    {Name: "ลูกชายคนโต"},
                    {Name: "ลูกสาวคนเล็ก"},
                },
            },
            {
                Name: "อา",
                Children: []*FamilyMember{
                    {Name: "ลูกอา"},
                },
            },
        },
    }

    printFamily(family, 0)
}

func printFamily(member *FamilyMember, level int) {
    indent := ""
    for i := 0; i < level; i++ {
        indent += "  "
    }
    fmt.Println(indent + member.Name)
    for _, child := range member.Children {
        printFamily(child, level+1)
    }
}
```

---

## 1.5 Embedded Struct — คอนเซปและการประยุกต์ใช้

Embedded Struct (การฝัง struct) เป็นวิธีที่ Go ใช้เพื่อการ **composition** แทนการ inheritance (สืบทอด)

### การเปรียบเทียบ Inheritance vs Composition

**Inheritance (ในภาษา OOP แบบดั้งเดิม):**

```java
// Java-style inheritance
class Animal {
    void speak() { ... }
}
class Dog extends Animal {
    void bark() { ... }
}
```

**Composition (ใน Go):**

```go
// Go-style composition
type Animal struct {
    Name string
}

func (a Animal) Speak() string {
    return "เสียงสัตว์"
}

type Dog struct {
    Animal  // ฝัง Animal
    Breed string
}

func (d Dog) Bark() string {
    return "โฮ่ง!"
}
```

### ตัวอย่าง Embedded Struct ระดับโปรเจกต์

```go
package main

import "fmt"
import "time"

// Base Model
type Model struct {
    ID        uint
    CreatedAt time.Time
    UpdatedAt time.Time
}

// User inherits from Model
type User struct {
    Model      // ฝัง Model
    Username string
    Email    string
}

// Product inherits from Model
type Product struct {
    Model       // ฝัง Model
    Name     string
    Price    float64
    InStock  bool
}

func main() {
    user := User{
        Model: Model{
            ID:        1,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        },
        Username: "สมชาย",
        Email:    "somchai@example.com",
    }

    // เข้าถึง field ของ Model ได้โดยตรง
    fmt.Printf("User ID: %d\n", user.ID)
    fmt.Printf("Username: %s\n", user.Username)
    fmt.Printf("Created: %v\n", user.CreatedAt)

    product := Product{
        Model: Model{
            ID:        100,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        },
        Name:    "แล็ปท็อป",
        Price:   25000.0,
        InStock: true,
    }

    fmt.Printf("Product ID: %d\n", product.ID)
    fmt.Printf("Product Name: %s\n", product.Name)
}
```

---

## สรุปบทที่ 1

ในบทนี้คุณได้เรียนรู้:

✅ การประกาศและใช้งาน Struct
✅ Struct กับ Pointer
✅ Anonymous Field
✅ Struct Embedding (การฝัง struct)
✅ Recursive Struct
✅ Composition แทน Inheritance

---

# บทที่ 2: Methods

---

## 2.1 The Need of Method

Method คือฟังก์ชันที่ **ผูกอยู่กับ type** โดยเฉพาะ (เรียกว่า **receiver**)

### ทำไมต้องมี Method?

1. **การจัดระเบียบโค้ด** — ฟังก์ชันที่เกี่ยวข้องกับ type ใด type หนึ่งถูกจัดกลุ่มไว้ด้วยกัน
2. **การเข้าถึงข้อมูล** — Method สามารถเข้าถึง field ของ receiver ได้โดยตรง
3. **การเขียนโค้ดที่อ่านง่าย** — `person.GetFullName()` อ่านง่ายกว่า `getFullName(person)`
4. **การ implement interface** — Method จำเป็นสำหรับการ implement interface

### เปรียบเทียบ Function vs Method

```go
package main

import "fmt"

type Rectangle struct {
    Width  float64
    Height float64
}

// Function (ไม่ผูกกับ type)
func areaFunc(r Rectangle) float64 {
    return r.Width * r.Height
}

// Method (ผูกกับ type Rectangle)
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}

    // เรียกใช้ function
    fmt.Println(areaFunc(rect))  // 50

    // เรียกใช้ method (อ่านง่ายกว่า)
    fmt.Println(rect.Area())     // 50
}
```

---

## 2.2 Method Declaration

### รูปแบบ

```go
func (receiver ReceiverType) MethodName(params) ReturnType {
    // body
}
```

### ตัวอย่าง

```go
package main

import "fmt"
import "math"

type Circle struct {
    Radius float64
}

// Method of Circle
func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Circumference() float64 {
    return 2 * math.Pi * c.Radius
}

func (c Circle) Diameter() float64 {
    return 2 * c.Radius
}

type Rectangle struct {
    Width  float64
    Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

func main() {
    circle := Circle{Radius: 5}
    fmt.Printf("Area: %.2f\n", circle.Area())
    fmt.Printf("Circumference: %.2f\n", circle.Circumference())

    rect := Rectangle{Width: 10, Height: 5}
    fmt.Printf("Area: %.2f\n", rect.Area())
    fmt.Printf("Perimeter: %.2f\n", rect.Perimeter())
}
```

### Method บน Type ที่ไม่ใช่ Struct

Method สามารถประกาศบน **type ที่ถูกสร้างขึ้นมาใหม่** (ไม่ใช่ built-in type)

```go
package main

import "fmt"

// สร้าง type ใหม่จาก int
type MyInt int

// Method บน MyInt
func (m MyInt) IsEven() bool {
    return int(m)%2 == 0
}

func (m MyInt) Double() MyInt {
    return m * 2
}

// สร้าง type ใหม่จาก string
type MyString string

func (s MyString) Reverse() string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func main() {
    var x MyInt = 42
    fmt.Println(x.IsEven())   // true
    fmt.Println(x.Double())   // 84

    s := MyString("Hello, Go!")
    fmt.Println(s.Reverse())  // !oG ,olleH
}
```

---

## 2.3 Method by Value vs by Reference

### Value Receiver (รับค่า)

Method ทำงานกับ **copy** ของ receiver — การแก้ไขไม่ส่งผลถึงต้นฉบับ

```go
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

// Value receiver: ไม่เปลี่ยนแปลงต้นฉบับ
func (p Person) UpdateAgeValue(newAge int) {
    p.Age = newAge
    fmt.Println("ใน method:", p.Age)
}

func main() {
    p := Person{Name: "สมชาย", Age: 30}
    fmt.Println("ก่อน:", p.Age)  // 30

    p.UpdateAgeValue(35)
    fmt.Println("หลัง:", p.Age)  // 30 (ไม่เปลี่ยนแปลง)
}
```

### Pointer Receiver (รับ pointer)

Method ทำงานกับ **ต้นฉบับ** — การแก้ไขส่งผลถึงต้นฉบับ

```go
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

// Pointer receiver: เปลี่ยนแปลงต้นฉบับ
func (p *Person) UpdateAgePointer(newAge int) {
    p.Age = newAge
    fmt.Println("ใน method:", p.Age)
}

// Pointer receiver สำหรับการแก้ไขหลาย fields
func (p *Person) SetNameAndAge(name string, age int) {
    p.Name = name
    p.Age = age
}

func main() {
    p := Person{Name: "สมชาย", Age: 30}
    fmt.Println("ก่อน:", p)  // {สมชาย 30}

    p.UpdateAgePointer(35)
    fmt.Println("หลัง:", p)  // {สมชาย 35}

    p.SetNameAndAge("สมชาย ใจดี", 40)
    fmt.Println("หลังเปลี่ยน:", p)  // {สมชาย ใจดี 40}
}
```

### ข้อควรจำ

| Receiver Type | ข้อดี | ข้อเสีย |
|---------------|------|--------|
| **Value** | ไม่แก้ไขต้นฉบับ, ปลอดภัย | เสีย performance เมื่อ struct ใหญ่ |
| **Pointer** | แก้ไขต้นฉบับได้, ประสิทธิภาพดี | ระวังการแก้ไขโดยไม่ตั้งใจ |

```go
// กฎทั่วไป:
// - ใช้ Value receiver ถ้า method ไม่ได้แก้ไข receiver
// - ใช้ Pointer receiver ถ้า method แก้ไข receiver หรือ struct ใหญ่
```

---

## 2.4 Method Best Practice และ Syntax Sugar

### การเลือก Value vs Pointer Receiver

```go
package main

import "fmt"

type BankAccount struct {
    Balance float64
    Owner   string
}

// ✅ ใช้ Value receiver (อ่านอย่างเดียว)
func (b BankAccount) GetBalance() float64 {
    return b.Balance
}

// ✅ ใช้ Pointer receiver (แก้ไข)
func (b *BankAccount) Deposit(amount float64) {
    b.Balance += amount
}

func (b *BankAccount) Withdraw(amount float64) error {
    if amount > b.Balance {
        return fmt.Errorf("ยอดเงินไม่พอ: %.2f > %.2f", amount, b.Balance)
    }
    b.Balance -= amount
    return nil
}

// ✅ Pointer receiver ถึงแม้จะไม่ได้แก้ไข (struct ใหญ่)
func (b *BankAccount) DisplayInfo() {
    fmt.Printf("เจ้าของ: %s, ยอดเงิน: %.2f\n", b.Owner, b.Balance)
}

func main() {
    account := BankAccount{Balance: 1000, Owner: "สมชาย"}

    account.Deposit(500)
    account.Withdraw(200)
    account.DisplayInfo()  // เจ้าของ: สมชาย, ยอดเงิน: 1300.00
}
```

### Syntax Sugar (Go จะทำ auto-conversion ให้)

```go
func main() {
    p := Person{Name: "สมชาย", Age: 30}

    // เรียก method ที่รับ pointer โดยใช้ value (Go ทำ auto-conversion)
    p.UpdateAgePointer(35)  // Go แปลงเป็น (&p).UpdateAgePointer(35)

    // เรียก method ที่รับ value โดยใช้ pointer
    pPtr := &p
    pPtr.GetBalance()  // Go แปลงเป็น (*pPtr).GetBalance()
}
```

---

## 2.5 Method with Nil Receiver

Method สามารถทำงานกับ **nil receiver** ได้

### ตัวอย่าง

```go
package main

import "fmt"

type LinkedList struct {
    Value int
    Next  *LinkedList
}

// Method บน pointer receiver ที่รองรับ nil
func (list *LinkedList) Print() {
    if list == nil {
        fmt.Println("(empty list)")
        return
    }

    current := list
    for current != nil {
        fmt.Printf("%d -> ", current.Value)
        current = current.Next
    }
    fmt.Println("nil")
}

// Append ต่อท้าย
func (list *LinkedList) Append(value int) *LinkedList {
    if list == nil {
        return &LinkedList{Value: value}
    }

    current := list
    for current.Next != nil {
        current = current.Next
    }
    current.Next = &LinkedList{Value: value}
    return list
}

func main() {
    var list *LinkedList

    // nil receiver
    list.Print()  // (empty list)

    // Append กับ nil
    list = list.Append(10)
    list = list.Append(20)
    list = list.Append(30)

    list.Print()  // 10 -> 20 -> 30 -> nil
}
```

---

## 2.6 Method Values

Method Values คือการเก็บ method ไว้ในตัวแปรเพื่อเรียกใช้ทีหลัง

### ตัวอย่าง

```go
package main

import "fmt"

type Calculator struct {
    Value int
}

func (c *Calculator) Add(x int) {
    c.Value += x
}

func (c *Calculator) Multiply(x int) {
    c.Value *= x
}

func (c *Calculator) Reset() {
    c.Value = 0
}

func main() {
    calc := &Calculator{Value: 10}

    // Method Values
    addFn := calc.Add
    mulFn := calc.Multiply

    addFn(5)        // calc.Value = 15
    mulFn(2)        // calc.Value = 30

    fmt.Println(calc.Value)  // 30

    // ใช้ method value ใน slice
    operations := []func(int){
        calc.Add,
        calc.Multiply,
        calc.Reset,
        calc.Add,
    }

    calc.Reset()
    for _, op := range operations {
        op(5)
    }
    fmt.Println(calc.Value)  // 5 (Reset -> Add 5)
}
```

---

## สรุปบทที่ 2

ในบทนี้คุณได้เรียนรู้:

✅ การประกาศ Method
✅ Value Receiver vs Pointer Receiver
✅ Method Best Practice
✅ Method with Nil Receiver
✅ Method Values

---

# บทที่ 3: Interfaces

---

## 3.1 Interfaces — คอนเซปพื้นฐาน

**Interface** ใน Go คือชุดของ method signatures (method ที่ต้องมี) Type ใดก็ตามที่มี method ตรงตาม interface จะถูกเรียกว่า implement interface นั้น

### รูปแบบ

```go
type InterfaceName interface {
    Method1(param1 Type1) ReturnType1
    Method2(param2 Type2) ReturnType2
}
```

### ตัวอย่าง

```go
package main

import "fmt"
import "math"

// ประกาศ interface
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Circle implements Shape
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

// Rectangle implements Shape
type Rectangle struct {
    Width  float64
    Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// ฟังก์ชันที่รับ interface
func printShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func main() {
    circle := Circle{Radius: 5}
    rect := Rectangle{Width: 10, Height: 5}

    printShapeInfo(circle)  // Area: 78.54, Perimeter: 31.42
    printShapeInfo(rect)    // Area: 50.00, Perimeter: 30.00
}
```

---

## 3.2 Interface Satisfaction — การทำให้ Interface เป็นจริง

Interface ใน Go เป็น **implicit** — ไม่ต้องประกาศว่า implement interface ไหน เพียงแค่มี method ครบก็ถือว่า implement

### ตัวอย่าง

```go
package main

import "fmt"

// Interface
type Writer interface {
    Write(data string) error
}

// Logger implement Writer โดยอัตโนมัติ (ไม่ต้องประกาศ)
type Logger struct {
    messages []string
}

func (l *Logger) Write(data string) error {
    l.messages = append(l.messages, data)
    fmt.Println("Logged:", data)
    return nil
}

// FileWriter ก็ implement Writer
type FileWriter struct {
    Filename string
}

func (f *FileWriter) Write(data string) error {
    fmt.Printf("Writing to %s: %s\n", f.Filename, data)
    return nil
}

// ฟังก์ชันที่รับ Writer interface
func processData(w Writer, data string) {
    w.Write(data)
}

func main() {
    logger := &Logger{}
    fileWriter := &FileWriter{Filename: "log.txt"}

    // ทั้ง logger และ fileWriter เป็น Writer
    processData(logger, "Hello")
    processData(fileWriter, "World")
}
```

---

## 3.3 Interface Value

Interface value ประกอบด้วย:
1. **Type** — ชนิดข้อมูลของ value ที่เก็บ
2. **Value** — ค่าจริงที่เก็บ

### ตัวอย่าง

```go
package main

import "fmt"

type Animal interface {
    Speak() string
}

type Dog struct{}

func (d Dog) Speak() string {
    return "โฮ่ง!"
}

type Cat struct{}

func (c Cat) Speak() string {
    return "เหมียว!"
}

func describe(i Animal) {
    fmt.Printf("Type: %T, Value: %v, Speak: %s\n", i, i, i.Speak())
}

func main() {
    var animal Animal

    animal = Dog{}
    describe(animal)  // Type: main.Dog, Value: {}, Speak: โฮ่ง!

    animal = Cat{}
    describe(animal)  // Type: main.Cat, Value: {}, Speak: เหมียว!

    // nil interface
    var nilAnimal Animal
    fmt.Printf("nilAnimal: %T, %v\n", nilAnimal, nilAnimal)
    // nilAnimal: <nil>, <nil>
}
```

---

## 3.4 Empty Interface

`interface{}` (หรือ `any` ใน Go 1.18+) คือ interface ที่ไม่มี method — ทุก type implement interface นี้

### ตัวอย่าง

```go
package main

import "fmt"

func printAnything(v interface{}) {
    fmt.Printf("Type: %T, Value: %v\n", v, v)
}

func main() {
    // empty interface สามารถเก็บอะไรก็ได้
    var anything interface{}

    anything = 42
    fmt.Println(anything)  // 42

    anything = "hello"
    fmt.Println(anything)  // hello

    anything = []int{1, 2, 3}
    fmt.Println(anything)  // [1 2 3]

    anything = struct{ Name string }{"สมชาย"}
    fmt.Println(anything)  // {สมชาย}

    // ใช้กับฟังก์ชัน
    printAnything(42)
    printAnything("Hello")
    printAnything([]float64{1.1, 2.2})
}
```

---

## 3.5 Type Assertion — ตอนที่ 1 และ 2

Type Assertion ใช้เพื่อดึงค่าเดิมจาก interface

### รูปแบบ

```go
value := interfaceVariable.(Type)      // panic ถ้า type ไม่ตรง
value, ok := interfaceVariable.(Type)  // safe
```

### ตัวอย่าง

```go
package main

import "fmt"

func main() {
    var i interface{} = "Hello, Go!"

    // Type assertion (unsafe)
    s := i.(string)
    fmt.Println(s)  // Hello, Go!

    // Type assertion (safe)
    if value, ok := i.(string); ok {
        fmt.Println("เป็น string:", value)
    } else {
        fmt.Println("ไม่ใช่ string")
    }

    // ลองกับ type ที่ไม่ใช่
    if value, ok := i.(int); ok {
        fmt.Println("เป็น int:", value)
    } else {
        fmt.Println("ไม่ใช่ int")
    }

    // panic ถ้า type ไม่ตรง
    // n := i.(int)  // panic: interface conversion: interface {} is string, not int
}
```

---

## 3.6 Type Assertion without Panic

การใช้ type assertion แบบปลอดภัย (comma-ok idiom)

### ตัวอย่าง

```go
package main

import "fmt"

func main() {
    var data interface{} = "some string"

    // Safe type assertion
    if s, ok := data.(string); ok {
        fmt.Println("String length:", len(s))
    } else {
        fmt.Println("Not a string")
    }

    // ตรวจสอบหลาย type
    checkType(data)  // string: some string

    checkType(42)    // int: 42
    checkType(true)  // bool: true
    checkType(3.14)  // float64: 3.14
    checkType([]int{1, 2})  // unknown type: []int
}

func checkType(v interface{}) {
    switch v := v.(type) {
    case string:
        fmt.Println("string:", v)
    case int:
        fmt.Println("int:", v)
    case bool:
        fmt.Println("bool:", v)
    case float64:
        fmt.Println("float64:", v)
    default:
        fmt.Printf("unknown type: %T\n", v)
    }
}
```

---

## 3.7 Type Switches (Discriminated Union)

Type Switch คือ switch ที่ใช้ตรวจสอบ type ของ interface

### ตัวอย่าง

```go
package main

import "fmt"

// ใช้ type switch เพื่อทำงานต่างกันตาม type
func processValue(v interface{}) {
    switch v := v.(type) {
    case nil:
        fmt.Println("nil value")
    case int:
        fmt.Printf("int: %d, คูณ 2 = %d\n", v, v*2)
    case string:
        fmt.Printf("string: %s, length: %d\n", v, len(v))
    case bool:
        fmt.Printf("bool: %t\n", v)
    case []int:
        sum := 0
        for _, n := range v {
            sum += n
        }
        fmt.Printf("[]int: %v, sum: %d\n", v, sum)
    case map[string]int:
        fmt.Printf("map: %v\n", v)
    default:
        fmt.Printf("unknown type: %T, value: %v\n", v, v)
    }
}

func main() {
    processValue(42)
    processValue("hello")
    processValue(true)
    processValue(3.14)
    processValue([]int{1, 2, 3, 4, 5})
    processValue(map[string]int{"a": 1, "b": 2})
    processValue(nil)
}
```

---

## 3.8 flag.Interface และการประยุกต์ใช้

`flag` package มี `FlagSet.Var` ที่ใช้ interface ในการกำหนด flag

### ตัวอย่าง

```go
package main

import (
    "flag"
    "fmt"
    "strings"
)

// สร้าง type ที่ implement flag.Value
type StringList []string

func (s *StringList) String() string {
    return strings.Join(*s, ",")
}

func (s *StringList) Set(value string) error {
    *s = append(*s, value)
    return nil
}

func main() {
    var names StringList

    // ใช้ flag.Var กับ interface
    flag.Var(&names, "name", "ชื่อ (สามารถใช้หลายครั้ง)")

    flag.Parse()

    fmt.Println("ชื่อที่ระบุ:")
    for i, name := range names {
        fmt.Printf("%d. %s\n", i+1, name)
    }
}

// การใช้งาน:
// go run main.go -name สมชาย -name สมหญิง -name วิชัย
```

---

## 3.9 sort.Interface — การเรียงลำดับด้วย Interface

`sort.Interface` เป็น interface ที่ใช้ในการเรียงลำดับ

```go
// sort.Interface มี 3 methods:
type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}
```

### ตัวอย่าง

```go
package main

import (
    "fmt"
    "sort"
)

type Person struct {
    Name string
    Age  int
}

// Sort by Age
type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }

// Sort by Name
type ByName []Person

func (n ByName) Len() int           { return len(n) }
func (n ByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByName) Less(i, j int) bool { return n[i].Name < n[j].Name }

func main() {
    people := []Person{
        {"สมชาย", 30},
        {"อานนท์", 25},
        {"วิชัย", 35},
        {"สมหญิง", 28},
    }

    fmt.Println("Original:", people)

    sort.Sort(ByAge(people))
    fmt.Println("Sorted by Age:", people)

    sort.Sort(ByName(people))
    fmt.Println("Sorted by Name:", people)

    // ใช้ sort.Slice (ง่ายกว่า)
    people2 := []Person{
        {"สมชาย", 30},
        {"อานนท์", 25},
        {"วิชัย", 35},
    }

    sort.Slice(people2, func(i, j int) bool {
        return people2[i].Age < people2[j].Age
    })
    fmt.Println("Sorted by Age (Slice):", people2)
}
```

---

## 3.10 error Interface

`error` เป็น interface ในตัวของ Go

```go
type error interface {
    Error() string
}
```

### ตัวอย่าง

```go
package main

import (
    "fmt"
)

// Custom error
type MyError struct {
    Code    int
    Message string
}

func (e MyError) Error() string {
    return fmt.Sprintf("Error #%d: %s", e.Code, e.Message)
}

// ฟังก์ชันที่คืน error
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, MyError{
            Code:    100,
            Message: "หารด้วยศูนย์ไม่ได้",
        }
    }
    return a / b, nil
}

func main() {
    result, err := divide(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
        // Type assertion เพื่อเข้าถึง field ของ custom error
        if myErr, ok := err.(MyError); ok {
            fmt.Printf("Error Code: %d\n", myErr.Code)
        }
        return
    }
    fmt.Println("Result:", result)
}
```

---

## 3.11 Interface ที่มี Methods ซ้ำกัน

ใน Go interface สามารถฝัง interface อื่นได้

### ตัวอย่าง

```go
package main

import "fmt"

// Basic interfaces
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

// Interface ที่ประกอบจากหลาย interface
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}

// ตัวอย่าง struct ที่ implement
type File struct {
    Name string
}

func (f File) Read(p []byte) (n int, err error) {
    fmt.Println("Reading from", f.Name)
    return len(p), nil
}

func (f File) Write(p []byte) (n int, err error) {
    fmt.Println("Writing to", f.Name)
    return len(p), nil
}

func (f File) Close() error {
    fmt.Println("Closing", f.Name)
    return nil
}

func main() {
    file := File{Name: "data.txt"}

    // file implements ReadWriteCloser
    var rw ReadWriteCloser = file
    rw.Read([]byte{})
    rw.Write([]byte{})
    rw.Close()
}
```

---

## สรุปบทที่ 3

ในบทนี้คุณได้เรียนรู้:

✅ Interface — คอนเซปพื้นฐาน
✅ Implicit Interface Satisfaction
✅ Interface Value
✅ Empty Interface (`interface{}` / `any`)
✅ Type Assertion และ Type Switch
✅ flag.Interface
✅ sort.Interface
✅ error Interface
✅ Interface Embedding

---

# บทที่ 4: JSON

---

## 4.1 JSON Introduction

**JSON (JavaScript Object Notation)** เป็นรูปแบบข้อมูลที่ใช้กัน广泛ใน Web API

### Struct Tags

Go ใช้ **struct tags** เพื่อควบคุมการแปลง JSON

```go
type Person struct {
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Email   string `json:"email,omitempty"`
    Address string `json:"address,omitempty"`
}
```

### Tags ที่สำคัญ

| Tag | ความหมาย |
|-----|----------|
| `json:"field"` | ใช้ชื่อ field ใน JSON |
| `json:"-"` | ข้าม field นี้ (ไม่แปลง) |
| `json:"field,omitempty"` | ข้ามถ้า value เป็น zero value |
| `json:",string"` | แปลงเป็น string ใน JSON |

---

## 4.2 JSON Unmarshal — แปลง JSON เป็น Struct

### ตัวอย่าง

```go
package main

import (
    "encoding/json"
    "fmt"
)

type Person struct {
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Email   string `json:"email,omitempty"`
    IsActive bool  `json:"is_active"`
}

func main() {
    // JSON string
    jsonData := `{
        "name": "สมชาย",
        "age": 30,
        "is_active": true
    }`

    // Unmarshal
    var person Person
    err := json.Unmarshal([]byte(jsonData), &person)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Printf("Name: %s\n", person.Name)
    fmt.Printf("Age: %d\n", person.Age)
    fmt.Printf("IsActive: %t\n", person.IsActive)
    fmt.Printf("Email: %q\n", person.Email)  // "" (zero value)
}
```

### การ Unmarshal ข้อมูลที่ซับซ้อน

```go
package main

import (
    "encoding/json"
    "fmt"
)

type Address struct {
    Street string `json:"street"`
    City   string `json:"city"`
    Zip    string `json:"zip"`
}

type Employee struct {
    ID       int      `json:"id"`
    Name     string   `json:"name"`
    Address  Address  `json:"address"`
    Skills   []string `json:"skills"`
    Projects []struct {
        Name string `json:"name"`
        Year int    `json:"year"`
    } `json:"projects"`
}

func main() {
    jsonData := `{
        "id": 1,
        "name": "สมชาย",
        "address": {
            "street": "123 ถนนสุขุมวิท",
            "city": "กรุงเทพฯ",
            "zip": "10110"
        },
        "skills": ["Go", "Python", "Docker"],
        "projects": [
            {"name": "Project A", "year": 2023},
            {"name": "Project B", "year": 2024}
        ]
    }`

    var emp Employee
    err := json.Unmarshal([]byte(jsonData), &emp)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Printf("ID: %d\n", emp.ID)
    fmt.Printf("Name: %s\n", emp.Name)
    fmt.Printf("Address: %s, %s, %s\n", emp.Address.Street, emp.Address.City, emp.Address.Zip)
    fmt.Printf("Skills: %v\n", emp.Skills)
    fmt.Printf("Projects: %v\n", emp.Projects)
}
```

---

## 4.3 JSON Marshal — แปลง Struct เป็น JSON

### ตัวอย่าง

```go
package main

import (
    "encoding/json"
    "fmt"
)

type Product struct {
    ID       int     `json:"id"`
    Name     string  `json:"name"`
    Price    float64 `json:"price"`
    InStock  bool    `json:"in_stock"`
    Category string  `json:"category,omitempty"`
    // ข้าม field นี้
    Internal string `json:"-"`
}

func main() {
    product := Product{
        ID:       1,
        Name:     "แล็ปท็อป",
        Price:    25000.0,
        InStock:  true,
        Category: "อุปกรณ์อิเล็กทรอนิกส์",
        Internal: "secret",
    }

    // Marshal
    jsonData, err := json.Marshal(product)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println(string(jsonData))
    // {"id":1,"name":"แล็ปท็อป","price":25000,"in_stock":true,"category":"อุปกรณ์อิเล็กทรอนิกส์"}

    // Marshal with indent (pretty print)
    jsonDataPretty, err := json.MarshalIndent(product, "", "  ")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println(string(jsonDataPretty))
}
```

### Omitempty

```go
package main

import (
    "encoding/json"
    "fmt"
)

type Person struct {
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Email   string `json:"email,omitempty"`   // ข้ามถ้า empty
    Phone   string `json:"phone,omitempty"`
    Address string `json:"address,omitempty"`
}

func main() {
    // มีเฉพาะ name และ age
    p1 := Person{Name: "สมชาย", Age: 30}
    data1, _ := json.Marshal(p1)
    fmt.Println(string(data1))  // {"name":"สมชาย","age":30}

    // มีทุก field
    p2 := Person{Name: "สมหญิง", Age: 25, Email: "somying@test.com", Phone: "0812345678"}
    data2, _ := json.Marshal(p2)
    fmt.Println(string(data2))  // {"name":"สมหญิง","age":25,"email":"somying@test.com","phone":"0812345678"}
}
```

---

## 4.4 ใช้งาน JSON Unmarshal กับ HTTP GET

### ตัวอย่าง

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type Post struct {
    UserID int    `json:"userId"`
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Body   string `json:"body"`
}

type Todo struct {
    UserID    int    `json:"userId"`
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}

func main() {
    // ดึงข้อมูลจาก API
    url := "https://jsonplaceholder.typicode.com/posts/1"

    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("HTTP Error:", err)
        return
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Read Error:", err)
        return
    }

    // Unmarshal JSON
    var post Post
    err = json.Unmarshal(body, &post)
    if err != nil {
        fmt.Println("JSON Error:", err)
        return
    }

    fmt.Printf("Post ID: %d\n", post.ID)
    fmt.Printf("Title: %s\n", post.Title)
    fmt.Printf("Body: %s\n", post.Body)

    // ตัวอย่างการดึงข้อมูลหลายๆ รายการ
    urlTodos := "https://jsonplaceholder.typicode.com/todos"

    respTodos, err := http.Get(urlTodos)
    if err != nil {
        fmt.Println("HTTP Error:", err)
        return
    }
    defer respTodos.Body.Close()

    bodyTodos, err := io.ReadAll(respTodos.Body)
    if err != nil {
        fmt.Println("Read Error:", err)
        return
    }

    var todos []Todo
    err = json.Unmarshal(bodyTodos, &todos)
    if err != nil {
        fmt.Println("JSON Error:", err)
        return
    }

    fmt.Printf("\nพบ %d todos\n", len(todos))
    for i, todo := range todos[:5] {
        fmt.Printf("%d. %s (completed: %t)\n", i+1, todo.Title, todo.Completed)
    }
}
```

---

## 4.5 JSON Encoder/Decoder

การใช้ Encoder/Decoder สำหรับ stream JSON

### ตัวอย่าง

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
)

type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    // --- Encoder ---
    persons := []Person{
        {"สมชาย", 30},
        {"สมหญิง", 25},
        {"วิชัย", 35},
    }

    // เขียน JSON ไปที่ stdout
    encoder := json.NewEncoder(os.Stdout)
    encoder.SetIndent("", "  ")

    for _, p := range persons {
        encoder.Encode(p)
    }

    // --- Decoder ---
    // อ่าน JSON จาก stdout (ในกรณีนี้เป็นตัวอย่าง)
    jsonData := `{"name":"สมชาย","age":30}`
    decoder := json.NewDecoder(os.Stdin)

    // จำลองการอ่านจาก stdin (ใช้ string จำลอง)
    // ในชีวิตจริงจะใช้ os.Stdin
    // decoder := json.NewDecoder(strings.NewReader(jsonData))

    var person Person
    err := decoder.Decode(&person)
    if err != nil {
        fmt.Println("Decode Error:", err)
        return
    }

    fmt.Printf("Decoded: %+v\n", person)
}
```

---

## 4.6 Struct Tags สำหรับ JSON

### Tags ที่มีประโยชน์

```go
package main

import (
    "encoding/json"
    "fmt"
)

type User struct {
    ID        int    `json:"id"`
    Username  string `json:"username"`
    Password  string `json:"-"`                          // ข้าม field นี้
    Email     string `json:"email,omitempty"`           // ข้ามถ้า empty
    FullName  string `json:"full_name,omitempty"`       // เปลี่ยนชื่อ field
    CreatedAt string `json:"created_at,omitempty"`
    // แปลงเป็น string ใน JSON
    BigNumber int64  `json:"big_number,string,omitempty"`
    // ใช้ - ถ้าต้องการข้าม
    Temp      string `json:",omitempty"`                // ใช้ชื่อ field
}

func main() {
    user := User{
        ID:        1,
        Username:  "somchai",
        Password:  "secret123",
        Email:     "somchai@test.com",
        FullName:  "สมชาย ใจดี",
        CreatedAt: "2024-01-01",
        BigNumber: 9223372036854775807,
    }

    jsonData, _ := json.MarshalIndent(user, "", "  ")
    fmt.Println(string(jsonData))

    // JSON ที่ไม่มี email และ full_name
    user2 := User{
        ID:       2,
        Username: "somchai2",
    }
    jsonData2, _ := json.MarshalIndent(user2, "", "  ")
    fmt.Println(string(jsonData2))
}
```

---

## สรุปบทที่ 4

ในบทนี้คุณได้เรียนรู้:

✅ JSON และ Struct Tags
✅ JSON Unmarshal — แปลง JSON เป็น Struct
✅ JSON Marshal — แปลง Struct เป็น JSON
✅ HTTP GET + JSON
✅ JSON Encoder/Decoder
✅ Struct Tags ที่มีประโยชน์

---

# บทที่ 5: Templates

---

## 5.1 Text Template Introduction

`text/template` ใช้สร้างข้อความแบบ dynamic

### ตัวอย่าง

```go
package main

import (
    "os"
    "text/template"
)

func main() {
    // Template string
    tmpl := `Hello, {{.Name}}! You are {{.Age}} years old.`

    // สร้าง template
    t := template.Must(template.New("greeting").Parse(tmpl))

    // ข้อมูล
    data := struct {
        Name string
        Age  int
    }{
        Name: "สมชาย",
        Age:  30,
    }

    // Execute
    err := t.Execute(os.Stdout, data)
    if err != nil {
        panic(err)
    }
    // Output: Hello, สมชาย! You are 30 years old.
}
```

---

## 5.2 Text Template — Actions และฟีเจอร์เพิ่มเติม

### Template Actions

| Action | ความหมาย |
|--------|----------|
| `{{.}}` | ค่าปัจจุบัน |
| `{{.Field}}` | เข้าถึง field |
| `{{if .Condition}} ... {{end}}` | เงื่อนไข |
| `{{range .Items}} ... {{end}}` | ลูป |
| `{{with .Value}} ... {{end}}` | กำหนด context |
| `{{template "name" .}}` | เรียก sub-template |

### ตัวอย่าง

```go
package main

import (
    "os"
    "text/template"
)

func main() {
    // Template ที่ซับซ้อน
    tmpl := `
ข้อมูลผู้ใช้:
ชื่อ: {{.Name}}
อายุ: {{.Age}}
{{if .IsActive}}สถานะ: Active{{else}}สถานะ: Inactive{{end}}

ทักษะ:
{{range .Skills}}
  - {{.}}
{{else}}
  (ไม่มีทักษะ)
{{end}}

{{with .Address}}
ที่อยู่: {{.Street}}, {{.City}} {{.Zip}}
{{end}}
`

    type Address struct {
        Street string
        City   string
        Zip    string
    }

    type User struct {
        Name     string
        Age      int
        IsActive bool
        Skills   []string
        Address  Address
    }

    user := User{
        Name:     "สมชาย",
        Age:      30,
        IsActive: true,
        Skills:   []string{"Go", "Python", "Docker"},
        Address: Address{
            Street: "123 ถนนสุขุมวิท",
            City:   "กรุงเทพฯ",
            Zip:    "10110",
        },
    }

    t := template.Must(template.New("user").Parse(tmpl))
    t.Execute(os.Stdout, user)
}
```

### Functions ใน Template

```go
package main

import (
    "os"
    "text/template"
    "strings"
)

func main() {
    // สร้าง custom function
    funcMap := template.FuncMap{
        "upper": strings.ToUpper,
        "lower": strings.ToLower,
        "join":  strings.Join,
        "add": func(a, b int) int {
            return a + b
        },
    }

    tmpl := `
ชื่อ: {{.Name | upper}}
ชื่อเล็ก: {{.Name | lower}}

ข้อมูล: {{join .Info " | "}}
อายุ + 10: {{add .Age 10}}
`

    data := struct {
        Name string
        Age  int
        Info []string
    }{
        Name: "Somchai",
        Age:  30,
        Info: []string{"Go", "Developer", "Bangkok"},
    }

    t := template.Must(template.New("test").Funcs(funcMap).Parse(tmpl))
    t.Execute(os.Stdout, data)
}
```

---

## 5.3 HTML Template

`html/template` เหมือน `text/template` แต่ **escape HTML** อัตโนมัติเพื่อป้องกัน XSS

### ตัวอย่าง

```go
package main

import (
    "html/template"
    "os"
)

func main() {
    tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
</head>
<body>
    <h1>{{.Heading}}</h1>
    <p>{{.Content}}</p>

    <h2>รายการสินค้า:</h2>
    <ul>
    {{range .Items}}
        <li>{{.Name}} - ฿{{.Price}}</li>
    {{else}}
        <li>ไม่มีสินค้า</li>
    {{end}}
    </ul>

    {{if .ShowUser}}
    <div>ผู้ใช้: {{.User}}</div>
    {{end}}
</body>
</html>
    `

    type Item struct {
        Name  string
        Price float64
    }

    data := struct {
        Title    string
        Heading  string
        Content  string
        Items    []Item
        ShowUser bool
        User     string
    }{
        Title:   "ร้านค้าออนไลน์",
        Heading: "ยินดีต้อนรับ",
        Content: "สินค้าคุณภาพ ราคาประหยัด",
        Items: []Item{
            {"แล็ปท็อป", 25000.0},
            {"สมาร์ทโฟน", 15000.0},
            {"หูฟัง", 2000.0},
        },
        ShowUser: true,
        User:     "สมชาย",
    }

    t := template.Must(template.New("index.html").Parse(tmpl))
    t.Execute(os.Stdout, data)
}
```

### HTML Template ที่ safe

```go
package main

import (
    "html/template"
    "os"
)

func main() {
    // HTML จะถูก escape อัตโนมัติ
    tmpl := `<div>{{.}}</div>`

    // อันตราย: XSS
    dangerous := `<script>alert('XSS')</script>`

    t := template.Must(template.New("safe").Parse(tmpl))
    t.Execute(os.Stdout, dangerous)
    // Output: <div>&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;</div>

    // ถ้าต้องการไม่ escape ใช้ template.HTML
    safeTmpl := `<div>{{.Content | safe}}</div>`

    data := struct {
        Content template.HTML
    }{
        Content: template.HTML("<strong>ข้อความที่ปลอดภัย</strong>"),
    }

    t2 := template.Must(template.New("safe2").Parse(safeTmpl))
    t2.Execute(os.Stdout, data)
    // Output: <div><strong>ข้อความที่ปลอดภัย</strong></div>
}
```

---

## สรุปบทที่ 5

ในบทนี้คุณได้เรียนรู้:

✅ Text Template — การใช้งานพื้นฐาน
✅ Template Actions (if, range, with)
✅ Custom Functions ใน Template
✅ HTML Template และการป้องกัน XSS

---

# บทที่ 6: Reflection

---

## 6.1 คอนเซปของ Reflection

**Reflection** คือความสามารถของโปรแกรมในการตรวจสอบและแก้ไขโครงสร้างของตัวเองในระหว่างรันไทม์

### ฟังก์ชันสำคัญใน `reflect` Package

| ฟังก์ชัน | ความหมาย |
|---------|----------|
| `reflect.TypeOf(v)` | ได้ชนิดข้อมูลของ v |
| `reflect.ValueOf(v)` | ได้ค่าและ metadata ของ v |
| `v.Kind()` | ชนิดพื้นฐาน (struct, int, string, etc.) |
| `v.NumField()` | จำนวน field ของ struct |
| `v.Field(i)` | field ที่ index i |
| `v.FieldByName(name)` | field ตามชื่อ |
| `v.Type().Field(i)` | metadata ของ field |

### ตัวอย่าง

```go
package main

import (
    "fmt"
    "reflect"
)

type Person struct {
    Name string
    Age  int
}

func main() {
    p := Person{Name: "สมชาย", Age: 30}

    // Type
    t := reflect.TypeOf(p)
    fmt.Printf("Type: %s\n", t.Name())
    fmt.Printf("Kind: %s\n", t.Kind())

    // Value
    v := reflect.ValueOf(p)
    fmt.Printf("Value: %v\n", v)

    // Field information
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        fmt.Printf("Field %d: %s = %v (%s)\n", i, field.Name, value.Interface(), field.Type)
    }

    // เปลี่ยนค่า (ต้องใช้ pointer)
    p2 := &Person{Name: "สมหญิง", Age: 25}
    v2 := reflect.ValueOf(p2).Elem()
    v2.FieldByName("Name").SetString("สมหญิง ใจดี")
    v2.FieldByName("Age").SetInt(26)

    fmt.Printf("After modification: %+v\n", p2)  // {Name:สมหญิง ใจดี Age:26}
}
```

---

## 6.2 Report Application — โครงสร้างโปรเจกต์

สร้างโปรเจกต์ที่ใช้ reflection เพื่อสร้าง report จาก struct

```go
// report/main.go
package main

import (
    "fmt"
    "reflect"
    "strings"
)

type User struct {
    ID       int    `report:"รหัสผู้ใช้"`
    Username string `report:"ชื่อผู้ใช้"`
    Email    string `report:"อีเมล"`
    Age      int    `report:"อายุ"`
    IsActive bool   `report:"สถานะ"`
}

type Product struct {
    ID    int     `report:"รหัสสินค้า"`
    Name  string  `report:"ชื่อสินค้า"`
    Price float64 `report:"ราคา"`
    Stock int     `report:"จำนวนคงเหลือ"`
}

func generateReport(data interface{}) {
    v := reflect.ValueOf(data)
    t := reflect.TypeOf(data)

    // ตรวจสอบว่าเป็น slice หรือ array
    if v.Kind() == reflect.Slice {
        // ถ้าเป็น slice แสดงรายการ
        for i := 0; i < v.Len(); i++ {
            item := v.Index(i)
            printItem(item, t.Elem())
            fmt.Println("---")
        }
    } else {
        printItem(v, t)
    }
}

func printItem(v reflect.Value, t reflect.Type) {
    for i := 0; i < v.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)

        // ใช้ tag "report" เป็นหัวข้อ
        label := field.Tag.Get("report")
        if label == "" {
            label = field.Name
        }

        fmt.Printf("%s: %v\n", label, value.Interface())
    }
}

func main() {
    users := []User{
        {ID: 1, Username: "somchai", Email: "somchai@test.com", Age: 30, IsActive: true},
        {ID: 2, Username: "somying", Email: "somying@test.com", Age: 25, IsActive: false},
    }

    fmt.Println("=== รายงานผู้ใช้ ===")
    generateReport(users)

    products := []Product{
        {ID: 1, Name: "แล็ปท็อป", Price: 25000.0, Stock: 10},
        {ID: 2, Name: "สมาร์ทโฟน", Price: 15000.0, Stock: 20},
    }

    fmt.Println("\n=== รายงานสินค้า ===")
    generateReport(products)
}
```

---

## 6.3 การเข้าถึง Struct Name

```go
package main

import (
    "fmt"
    "reflect"
)

type User struct {
    Name string
    Age  int
}

func getStructName(data interface{}) string {
    t := reflect.TypeOf(data)

    // ถ้าเป็น pointer ให้ Elem()
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }

    if t.Kind() == reflect.Struct {
        return t.Name()
    }

    if t.Kind() == reflect.Slice {
        return t.Elem().Name()
    }

    return "unknown"
}

func main() {
    user := User{Name: "สมชาย", Age: 30}
    users := []User{{"สมชาย", 30}, {"สมหญิง", 25}}

    fmt.Println(getStructName(user))   // User
    fmt.Println(getStructName(&user))  // User
    fmt.Println(getStructName(users))  // User
    fmt.Println(getStructName(42))     // unknown
}
```

---

## 6.4 การเข้าถึง Struct Value

```go
package main

import (
    "fmt"
    "reflect"
)

type Person struct {
    Name string
    Age  int
    City string
}

func inspectValue(data interface{}) {
    v := reflect.ValueOf(data)

    // ถ้าเป็น pointer ให้ Elem()
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }

    if v.Kind() != reflect.Struct {
        fmt.Println("Not a struct")
        return
    }

    fmt.Printf("Struct with %d fields:\n", v.NumField())

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        fmt.Printf("  Field %d: %v (type: %s, kind: %s)\n",
            i, field.Interface(),
            field.Type().Name(),
            field.Kind())
    }
}

func main() {
    p := Person{Name: "สมชาย", Age: 30, City: "กรุงเทพฯ"}
    inspectValue(p)

    fmt.Println()

    inspectValue(&p)
}
```

---

## 6.5 การเข้าถึง Struct Tag

```go
package main

import (
    "fmt"
    "reflect"
)

type Person struct {
    Name    string `json:"name" validate:"required" label:"ชื่อ"`
    Age     int    `json:"age" validate:"min=0,max=150" label:"อายุ"`
    Email   string `json:"email,omitempty" validate:"email" label:"อีเมล"`
    Address string `json:"address,omitempty" label:"ที่อยู่"`
}

func inspectTags(data interface{}) {
    t := reflect.TypeOf(data)

    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }

    if t.Kind() != reflect.Struct {
        fmt.Println("Not a struct")
        return
    }

    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        fmt.Printf("Field: %s\n", field.Name)
        fmt.Printf("  Type: %s\n", field.Type)
        fmt.Printf("  JSON: %s\n", field.Tag.Get("json"))
        fmt.Printf("  Validate: %s\n", field.Tag.Get("validate"))
        fmt.Printf("  Label: %s\n", field.Tag.Get("label"))
        fmt.Println()
    }
}

func main() {
    p := Person{Name: "สมชาย", Age: 30}
    inspectTags(p)
}
```

---

## 6.6 Multiple Tags

```go
package main

import (
    "fmt"
    "reflect"
    "strings"
)

type Config struct {
    Host     string `env:"HOST" default:"localhost" desc:"Server host"`
    Port     int    `env:"PORT" default:"8080" desc:"Server port"`
    Debug    bool   `env:"DEBUG" default:"false" desc:"Debug mode"`
    Database string `env:"DB_URL" default:"postgres://localhost:5432" desc:"Database URL"`
}

func parseTags(data interface{}) {
    t := reflect.TypeOf(data)
    v := reflect.ValueOf(data)

    if t.Kind() == reflect.Ptr {
        t = t.Elem()
        v = v.Elem()
    }

    if t.Kind() != reflect.Struct {
        return
    }

    fmt.Println("=== Configuration ===")

    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)

        env := field.Tag.Get("env")
        defaultVal := field.Tag.Get("default")
        desc := field.Tag.Get("desc")

        fmt.Printf("%s:\n", env)
        fmt.Printf("  Value: %v\n", value.Interface())
        fmt.Printf("  Default: %s\n", defaultVal)
        fmt.Printf("  Description: %s\n", desc)
        fmt.Println()
    }
}

func loadFromEnv(data interface{}) {
    // จำลองการ load จาก environment variables
    // ในชีวิตจริงจะใช้ os.Getenv()

    // สมมติว่าตั้งค่า: PORT=9000
    v := reflect.ValueOf(data).Elem()

    for i := 0; i < v.NumField(); i++ {
        field := v.Type().Field(i)
        env := field.Tag.Get("env")

        // จำลอง: ถ้า env เป็น PORT ตั้งค่าเป็น 9000
        if env == "PORT" {
            v.Field(i).SetInt(9000)
        }
        if env == "DEBUG" {
            v.Field(i).SetBool(true)
        }
    }
}

func main() {
    config := &Config{
        Host:     "localhost",
        Port:     8080,
        Debug:    false,
        Database: "postgres://localhost:5432",
    }

    parseTags(config)

    fmt.Println("=== Before load ===")
    fmt.Printf("%+v\n", config)

    loadFromEnv(config)

    fmt.Println("=== After load ===")
    fmt.Printf("%+v\n", config)
}
```

---

## 6.7 การแก้ไข Struct Value

```go
package main

import (
    "fmt"
    "reflect"
)

type Person struct {
    Name  string
    Age   int
    Score float64
    Tags  []string
}

func setField(data interface{}, fieldName string, value interface{}) error {
    v := reflect.ValueOf(data)

    // ต้องเป็น pointer
    if v.Kind() != reflect.Ptr {
        return fmt.Errorf("must be pointer")
    }

    v = v.Elem()

    // ตรวจสอบว่าเป็น struct
    if v.Kind() != reflect.Struct {
        return fmt.Errorf("must be struct")
    }

    // หา field
    field := v.FieldByName(fieldName)
    if !field.IsValid() {
        return fmt.Errorf("field %s not found", fieldName)
    }

    // ตรวจสอบว่าแก้ไขได้
    if !field.CanSet() {
        return fmt.Errorf("field %s cannot be set", fieldName)
    }

    // แปลงค่า
    val := reflect.ValueOf(value)

    // ตรวจสอบ type
    if val.Type() != field.Type() {
        return fmt.Errorf("type mismatch: expected %s, got %s",
            field.Type(), val.Type())
    }

    field.Set(val)
    return nil
}

func main() {
    p := Person{Name: "สมชาย", Age: 30, Score: 75.5, Tags: []string{"go", "dev"}}

    fmt.Println("Before:", p)

    // แก้ไข field
    setField(&p, "Name", "สมชาย ใจดี")
    setField(&p, "Age", 35)
    setField(&p, "Score", 85.0)
    setField(&p, "Tags", []string{"go", "dev", "golang"})

    fmt.Println("After:", p)

    // Error handling
    err := setField(&p, "Unknown", 100)
    if err != nil {
        fmt.Println("Error:", err)
    }

    err = setField(p, "Name", "test") // ไม่ใช้ pointer
    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

---

## 6.8 Reflection สรุป

### ข้อดีของ Reflection

✅ ทำงานกับข้อมูลที่ไม่รู้ type ล่วงหน้า
✅ สร้าง Generic functions
✅ Framework, ORM, Serialization
✅ Code generation

### ข้อเสียของ Reflection

❌ ช้ากว่า (runtime overhead)
❌ ไม่ปลอดภัย (type checking ที่ compile time หายไป)
❌ โค้ดอ่านยากขึ้น

### เมื่อควรใช้ Reflection

| ควรใช้ | ไม่ควรใช้ |
|--------|-----------|
| JSON marshal/unmarshal | แก้ปัญหาแบบ static typing |
| ORM (GORM) | แก้ปัญหาเฉพาะ type |
| Template engine | สร้างฟังก์ชันที่รู้ type ล่วงหน้า |
| Dependency injection | เว้นแต่จำเป็นจริงๆ |
| Generic libraries | |

### ตัวอย่างการใช้ Reflection ใน Framework

```go
package main

import (
    "fmt"
    "reflect"
)

// ตัวอย่าง JSON encoder อย่างง่าย
func simpleJSON(v interface{}) string {
    val := reflect.ValueOf(v)

    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }

    if val.Kind() != reflect.Struct {
        return fmt.Sprintf("%v", v)
    }

    result := "{"
    for i := 0; i < val.NumField(); i++ {
        field := val.Type().Field(i)
        value := val.Field(i)

        if i > 0 {
            result += ","
        }

        jsonTag := field.Tag.Get("json")
        if jsonTag == "" || jsonTag == "-" {
            continue
        }

        result += fmt.Sprintf(`"%s":`, jsonTag)

        // จัดการ type
        switch value.Kind() {
        case reflect.String:
            result += fmt.Sprintf(`"%s"`, value.Interface())
        case reflect.Int, reflect.Int64:
            result += fmt.Sprintf("%d", value.Interface())
        case reflect.Float64:
            result += fmt.Sprintf("%f", value.Interface())
        case reflect.Bool:
            result += fmt.Sprintf("%t", value.Interface())
        default:
            result += fmt.Sprintf(`"%v"`, value.Interface())
        }
    }
    result += "}"
    return result
}

type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
    Stock int     `json:"stock,omitempty"`
}

func main() {
    p := Product{ID: 1, Name: "แล็ปท็อป", Price: 25000.0, Stock: 10}
    fmt.Println(simpleJSON(p))
}
```

---

## สรุปบทที่ 6

ในบทนี้คุณได้เรียนรู้:

✅ คอนเซปของ Reflection
✅ TypeOf และ ValueOf
✅ การเข้าถึง Struct Name, Value, Tag
✅ การใช้ Tags ในการเก็บ metadata
✅ การแก้ไข Struct Value ด้วย Reflection
✅ ข้อดี/ข้อเสียของ Reflection
✅ การนำ Reflection ไปประยุกต์ใช้

---

# บทที่ 7: การจัดการ Error ขั้นสูง

---

## 7.1 Project Overview — Validation Library

สร้าง Validation Library ที่ใช้ error handling ขั้นสูง

### โครงสร้างโปรเจกต์

```
validator/
├── go.mod
├── validator/
│   └── validator.go
└── main.go
```

---

## 7.2 Initial Project with Simple Length Validation

```go
// validator/validator.go
package validator

import (
    "fmt"
    "reflect"
    "strings"
)

type ValidationError struct {
    Field string
    Tag   string
    Value interface{}
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("Field '%s' failed validation '%s' (value: %v)",
        e.Field, e.Tag, e.Value)
}

func Validate(data interface{}) []error {
    var errors []error

    v := reflect.ValueOf(data)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }

    if v.Kind() != reflect.Struct {
        return errors
    }

    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        tag := t.Field(i).Tag.Get("validate")

        if tag == "" {
            continue
        }

        // ตรวจสอบ validation
        for _, rule := range strings.Split(tag, ",") {
            rule = strings.TrimSpace(rule)

            if strings.HasPrefix(rule, "min_len=") {
                // min_len=3
                min := 0
                fmt.Sscanf(rule, "min_len=%d", &min)

                if field.Kind() == reflect.String && len(field.String()) < min {
                    errors = append(errors, ValidationError{
                        Field: t.Field(i).Name,
                        Tag:   "min_len",
                        Value: field.String(),
                    })
                }
            }
        }
    }

    return errors
}
```

### main.go

```go
package main

import (
    "fmt"
    "validator/validator"
)

type User struct {
    Name  string `validate:"min_len=3"`
    Email string `validate:"min_len=5"`
    Age   int
}

func main() {
    user := User{
        Name:  "A",
        Email: "a@b",
        Age:   30,
    }

    errors := validator.Validate(user)

    if len(errors) > 0 {
        fmt.Println("Validation errors:")
        for _, err := range errors {
            fmt.Println("  -", err)
        }
    } else {
        fmt.Println("Validation passed!")
    }
}
```

---

## 7.3 Add More Validations Chain

```go
// validator/validator.go (เพิ่มเติม)
package validator

import (
    "fmt"
    "reflect"
    "regexp"
    "strconv"
    "strings"
)

type ValidationError struct {
    Field string
    Tag   string
    Value interface{}
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("Field '%s' failed validation '%s' (value: %v)",
        e.Field, e.Tag, e.Value)
}

// เพิ่ม Error Wrapping
func (e ValidationError) Unwrap() error {
    return fmt.Errorf("validation error")
}

func Validate(data interface{}) []error {
    var errors []error

    v := reflect.ValueOf(data)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }

    if v.Kind() != reflect.Struct {
        return errors
    }

    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        tag := t.Field(i).Tag.Get("validate")

        if tag == "" {
            continue
        }

        for _, rule := range strings.Split(tag, ",") {
            rule = strings.TrimSpace(rule)

            switch {
            case strings.HasPrefix(rule, "min_len="):
                min := 0
                fmt.Sscanf(rule, "min_len=%d", &min)
                if field.Kind() == reflect.String && len(field.String()) < min {
                    errors = append(errors, ValidationError{
                        Field: t.Field(i).Name,
                        Tag:   "min_len",
                        Value: field.String(),
                    })
                }

            case strings.HasPrefix(rule, "max_len="):
                max := 0
                fmt.Sscanf(rule, "max_len=%d", &max)
                if field.Kind() == reflect.String && len(field.String()) > max {
                    errors = append(errors, ValidationError{
                        Field: t.Field(i).Name,
                        Tag:   "max_len",
                        Value: field.String(),
                    })
                }

            case rule == "email":
                if field.Kind() == reflect.String {
                    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
                    if !emailRegex.MatchString(field.String()) {
                        errors = append(errors, ValidationError{
                            Field: t.Field(i).Name,
                            Tag:   "email",
                            Value: field.String(),
                        })
                    }
                }

            case strings.HasPrefix(rule, "min="):
                min := 0
                fmt.Sscanf(rule, "min=%d", &min)
                if field.Kind() == reflect.Int && int(field.Int()) < min {
                    errors = append(errors, ValidationError{
                        Field: t.Field(i).Name,
                        Tag:   "min",
                        Value: field.Int(),
                    })
                }

            case strings.HasPrefix(rule, "max="):
                max := 0
                fmt.Sscanf(rule, "max=%d", &max)
                if field.Kind() == reflect.Int && int(field.Int()) > max {
                    errors = append(errors, ValidationError{
                        Field: t.Field(i).Name,
                        Tag:   "max",
                        Value: field.Int(),
                    })
                }

            case rule == "required":
                if field.Kind() == reflect.String && field.String() == "" {
                    errors = append(errors, ValidationError{
                        Field: t.Field(i).Name,
                        Tag:   "required",
                        Value: "",
                    })
                }
                if field.Kind() == reflect.Int && field.Int() == 0 {
                    errors = append(errors, ValidationError{
                        Field: t.Field(i).Name,
                        Tag:   "required",
                        Value: 0,
                    })
                }
            }
        }
    }

    return errors
}
```

---

## 7.4 Refactored into Package

```go
// validator/validator.go (เวอร์ชันที่ refactored)
package validator

import (
    "fmt"
    "reflect"
    "regexp"
    "strconv"
    "strings"
)

// ValidationError struct
type ValidationError struct {
    Field string
    Tag   string
    Value interface{}
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("Field '%s' failed validation '%s' (value: %v)",
        e.Field, e.Tag, e.Value)
}

// Rule interface
type Rule interface {
    Validate(field reflect.Value, tag string) error
}

// Rule implementations
type MinLenRule struct{}

func (r MinLenRule) Validate(field reflect.Value, tag string) error {
    min := 0
    fmt.Sscanf(tag, "min_len=%d", &min)
    if field.Kind() == reflect.String && len(field.String()) < min {
        return ValidationError{
            Field: "unknown",
            Tag:   "min_len",
            Value: field.String(),
        }
    }
    return nil
}

type MaxLenRule struct{}

func (r MaxLenRule) Validate(field reflect.Value, tag string) error {
    max := 0
    fmt.Sscanf(tag, "max_len=%d", &max)
    if field.Kind() == reflect.String && len(field.String()) > max {
        return ValidationError{
            Field: "unknown",
            Tag:   "max_len",
            Value: field.String(),
        }
    }
    return nil
}

type EmailRule struct{}

func (r EmailRule) Validate(field reflect.Value, tag string) error {
    if field.Kind() == reflect.String {
        emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
        if !emailRegex.MatchString(field.String()) {
            return ValidationError{
                Field: "unknown",
                Tag:   "email",
                Value: field.String(),
            }
        }
    }
    return nil
}

// Registry
var ruleRegistry = map[string]Rule{
    "min_len": MinLenRule{},
    "max_len": MaxLenRule{},
    "email":   EmailRule{},
}

func Validate(data interface{}) []error {
    var errors []error

    v := reflect.ValueOf(data)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }

    if v.Kind() != reflect.Struct {
        return errors
    }

    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        fieldType := t.Field(i)
        tag := fieldType.Tag.Get("validate")

        if tag == "" {
            continue
        }

        for _, ruleTag := range strings.Split(tag, ",") {
            ruleTag = strings.TrimSpace(ruleTag)

            // หา rule
            ruleName := strings.SplitN(ruleTag, "=", 2)[0]
            rule, exists := ruleRegistry[ruleName]
            if !exists {
                continue
            }

            // Validate
            err := rule.Validate(field, ruleTag)
            if err != nil {
                if ve, ok := err.(ValidationError); ok {
                    ve.Field = fieldType.Name
                    errors = append(errors, ve)
                } else {
                    errors = append(errors, err)
                }
            }
        }
    }

    return errors
}
```

---

## 7.5 First Look at errors.Is และ Unwrap Interface

```go
package main

import (
    "errors"
    "fmt"
)

// Custom error types
type NotFoundError struct {
    Resource string
    ID       int
}

func (e NotFoundError) Error() string {
    return fmt.Sprintf("%s with ID %d not found", e.Resource, e.ID)
}

// Implement Unwrap interface
func (e NotFoundError) Unwrap() error {
    return errors.New("resource not found")
}

type DatabaseError struct {
    Operation string
    Err       error
}

func (e DatabaseError) Error() string {
    return fmt.Sprintf("database %s failed: %v", e.Operation, e.Err)
}

func (e DatabaseError) Unwrap() error {
    return e.Err
}

func main() {
    // สร้าง error chain
    baseErr := errors.New("connection timeout")
    dbErr := DatabaseError{Operation: "query", Err: baseErr}
    notFoundErr := NotFoundError{Resource: "User", ID: 123}
    
    wrappedErr := fmt.Errorf("API error: %w", dbErr)

    // ตรวจสอบด้วย errors.Is
    fmt.Println("errors.Is(wrappedErr, dbErr):", errors.Is(wrappedErr, dbErr))
    fmt.Println("errors.Is(wrappedErr, baseErr):", errors.Is(wrappedErr, baseErr))

    // Unwrap
    fmt.Println("\nUnwrapping:")
    err := wrappedErr
    for err != nil {
        fmt.Println("  -", err)
        err = errors.Unwrap(err)
    }

    // errors.Is กับ custom error
    fmt.Println("\nerrors.Is(notFoundErr, NotFoundError{}):",
        errors.Is(notFoundErr, NotFoundError{}))
}
```

---

## 7.6 Multi-layer Wrapping and Unwrapping Error

```go
package main

import (
    "errors"
    "fmt"
)

type Layer1Error struct {
    Msg string
    Err error
}

func (e Layer1Error) Error() string {
    return fmt.Sprintf("layer1: %s", e.Msg)
}

func (e Layer1Error) Unwrap() error {
    return e.Err
}

type Layer2Error struct {
    Msg string
    Err error
}

func (e Layer2Error) Error() string {
    return fmt.Sprintf("layer2: %s", e.Msg)
}

func (e Layer2Error) Unwrap() error {
    return e.Err
}

type Layer3Error struct {
    Msg string
}

func (e Layer3Error) Error() string {
    return fmt.Sprintf("layer3: %s", e.Msg)
}

func main() {
    // สร้าง error chain 3 layer
    base := Layer3Error{Msg: "database connection failed"}
    layer2 := Layer2Error{Msg: "query execution failed", Err: base}
    layer1 := Layer1Error{Msg: "API call failed", Err: layer2}

    err := layer1

    // Unwrap ทั้งหมด
    fmt.Println("Full error chain:")
    for err != nil {
        fmt.Printf("  - %v\n", err)
        err = errors.Unwrap(err)
    }

    // ตรวจสอบ
    fmt.Println("\nerrors.Is(layer1, Layer3Error{}):",
        errors.Is(layer1, Layer3Error{}))

    // สร้าง error ที่ไม่ใช่ chain
    simpleErr := errors.New("simple error")
    wrapped := fmt.Errorf("wrapped: %w", simpleErr)

    fmt.Println("\nerrors.Is(wrapped, simpleErr):", errors.Is(wrapped, simpleErr))
}
```

---

## 7.7 Tough Life before errors.As

```go
package main

import (
    "fmt"
    "strings"
)

type APIError struct {
    Code    int
    Message string
}

func (e APIError) Error() string {
    return fmt.Sprintf("API error %d: %s", e.Code, e.Message)
}

type BusinessError struct {
    Code    int
    Message string
}

func (e BusinessError) Error() string {
    return fmt.Sprintf("Business error %d: %s", e.Code, e.Message)
}

// ก่อนมี errors.As ต้องใช้ type assertion แบบ manual
func extractAPIError(err error) *APIError {
    if err == nil {
        return nil
    }

    // ตรวจสอบ error chain ด้วยตัวเอง
    current := err
    for current != nil {
        if apiErr, ok := current.(APIError); ok {
            return &apiErr
        }
        if apiErr, ok := current.(*APIError); ok {
            return apiErr
        }

        // ถ้าเป็น interface ที่มี Unwrap
        type unwrapper interface {
            Unwrap() error
        }

        if unw, ok := current.(unwrapper); ok {
            current = unw.Unwrap()
        } else {
            break
        }
    }
    return nil
}

func main() {
    apiErr := APIError{Code: 404, Message: "Not found"}
    wrapped := fmt.Errorf("handler error: %w", apiErr)

    // แบบ manual
    extracted := extractAPIError(wrapped)
    if extracted != nil {
        fmt.Printf("Found API error: %d - %s\n", extracted.Code, extracted.Message)
    }

    // เปรียบเทียบกับ errors.As (จะใช้ในหัวข้อถัดไป)
    fmt.Println("\nตัวอย่าง manual extraction ที่ยุ่งยาก...")
    fmt.Println("ต้อง implement Unwrap detection ด้วยตัวเอง")
}
```

---

## 7.8 Type Assertion with errors.As

### เปรียบเทียบ errors.Is vs errors.As

| ฟังก์ชัน | ตรวจสอบ | ตัวอย่าง |
|---------|---------|----------|
| `errors.Is` | error เท่ากัน | `errors.Is(err, ErrNotFound)` |
| `errors.As` | type ของ error | `errors.As(err, &apiErr)` |

### ตัวอย่าง

```go
package main

import (
    "errors"
    "fmt"
)

type NotFoundError struct {
    Resource string
    ID       int
}

func (e NotFoundError) Error() string {
    return fmt.Sprintf("%s with ID %d not found", e.Resource, e.ID)
}

type PermissionError struct {
    User  string
    Action string
}

func (e PermissionError) Error() string {
    return fmt.Sprintf("user %s has no permission for %s", e.User, e.Action)
}

type ValidationError struct {
    Field string
    Value string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Value)
}

func main() {
    // สร้าง error chain
    base := ValidationError{Field: "email", Value: "invalid"}
    perm := PermissionError{User: "somchai", Action: "delete"}
    notFound := NotFoundError{Resource: "User", ID: 123}

    // ใช้ errors.As
    var notFoundErr NotFoundError
    var permErr PermissionError
    var validationErr ValidationError

    err := fmt.Errorf("operation failed: %w", notFound)

    if errors.As(err, &notFoundErr) {
        fmt.Printf("Found NotFoundError: %s ID=%d\n",
            notFoundErr.Resource, notFoundErr.ID)
    }

    err2 := fmt.Errorf("operation failed: %w", perm)
    if errors.As(err2, &permErr) {
        fmt.Printf("Found PermissionError: user=%s action=%s\n",
            permErr.User, permErr.Action)
    }

    err3 := fmt.Errorf("operation failed: %w", base)
    if errors.As(err3, &validationErr) {
        fmt.Printf("Found ValidationError: field=%s value=%s\n",
            validationErr.Field, validationErr.Value)
    }
}
```

---

## 7.9 Refactored Custom As Function

```go
package main

import (
    "errors"
    "fmt"
)

type CustomError struct {
    Code    int
    Message string
}

func (e CustomError) Error() string {
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Custom As function
func customAs(err error, target interface{}) bool {
    if err == nil {
        return false
    }

    // ใช้ errors.As จริง
    return errors.As(err, target)
}

// Extension: AsError
func AsError[T error](err error) (T, bool) {
    var target T
    if errors.As(err, &target) {
        return target, true
    }
    return target, false
}

func main() {
    customErr := CustomError{Code: 1001, Message: "Invalid input"}
    wrapped := fmt.Errorf("wrapped: %w", customErr)

    // ใช้ customAs
    var target CustomError
    if customAs(wrapped, &target) {
        fmt.Printf("Found: %d - %s\n", target.Code, target.Message)
    }

    // ใช้ AsError (generic)
    if err, ok := AsError[CustomError](wrapped); ok {
        fmt.Printf("Found (generic): %d - %s\n", err.Code, err.Message)
    }

    // ถ้าไม่เจอ
    if _, ok := AsError[*CustomError](errors.New("other error")); !ok {
        fmt.Println("No CustomError found")
    }
}
```

---

## 7.10 Add More Error Types

```go
package main

import (
    "errors"
    "fmt"
)

// Error types
type TimeoutError struct {
    Duration string
}

func (e TimeoutError) Error() string {
    return fmt.Sprintf("operation timeout after %s", e.Duration)
}

type RetryableError struct {
    Attempts int
    Err      error
}

func (e RetryableError) Error() string {
    return fmt.Sprintf("retryable error after %d attempts: %v", e.Attempts, e.Err)
}

func (e RetryableError) Unwrap() error {
    return e.Err
}

type DataIntegrityError struct {
    Table  string
    Record string
}

func (e DataIntegrityError) Error() string {
    return fmt.Sprintf("data integrity error in %s: %s", e.Table, e.Record)
}

// Business logic
func fetchData(id int) error {
    // Simulate timeout
    if id == 0 {
        return TimeoutError{Duration: "5s"}
    }

    // Simulate retry
    if id == 1 {
        return RetryableError{
            Attempts: 3,
            Err:      errors.New("connection refused"),
        }
    }

    // Simulate data integrity
    if id == 2 {
        return DataIntegrityError{
            Table:  "users",
            Record: fmt.Sprintf("id=%d", id),
        }
    }

    return nil
}

func handleError(err error) {
    if err == nil {
        fmt.Println("Success!")
        return
    }

    // ตรวจสอบ type ด้วย errors.As
    var timeout TimeoutError
    if errors.As(err, &timeout) {
        fmt.Printf("Timeout: %s\n", timeout.Duration)
        return
    }

    var retry RetryableError
    if errors.As(err, &retry) {
        fmt.Printf("Retryable: attempts=%d, cause=%v\n",
            retry.Attempts, retry.Err)
        return
    }

    var integrity DataIntegrityError
    if errors.As(err, &integrity) {
        fmt.Printf("Integrity: table=%s, record=%s\n",
            integrity.Table, integrity.Record)
        return
    }

    fmt.Printf("Unknown error: %v\n", err)
}

func main() {
    fmt.Print("ID=0: ")
    handleError(fetchData(0))

    fmt.Print("ID=1: ")
    handleError(fetchData(1))

    fmt.Print("ID=2: ")
    handleError(fetchData(2))

    fmt.Print("ID=3: ")
    handleError(fetchData(3))
}
```

---

## สรุปบทที่ 7

ในบทนี้คุณได้เรียนรู้:

✅ Validation Library ด้วย Reflection
✅ Multiple Validations Chain
✅ Error Wrapping และ Unwrapping
✅ `errors.Is` และ `errors.As`
✅ Multi-layer Error Chain
✅ Custom As Function
✅ Error Types ในระดับโปรเจกต์

---

# 🎯 แบบฝึกหัดทบทวน เล่มที่ 2

## แบบฝึกหัดที่ 1: Struct
สร้าง struct `Book` และฟังก์ชันที่แสดงข้อมูล

## แบบฝึกหัดที่ 2: Method
เพิ่ม Method ให้ `Book` เพื่อคำนวณราคาหลังส่วนลด

## แบบฝึกหัดที่ 3: Interface
สร้าง interface `Printable` และ implement สำหรับ `Book`

## แบบฝึกหัดที่ 4: JSON
แปลง `Book` เป็น JSON และแปลงกลับ

## แบบฝึกหัดที่ 5: Template
สร้าง template สำหรับแสดงข้อมูล `Book`

## แบบฝึกหัดที่ 6: Reflection
เขียนฟังก์ชันที่พิมพ์ struct fields พร้อม tags

## แบบฝึกหัดที่ 7: Error Handling
สร้าง custom error สำหรับการค้นหาหนังสือไม่พบ

---

**จบ เล่มที่ 2: โครงสร้างข้อมูลขั้นสูงและการจัดการ JSON**

---
