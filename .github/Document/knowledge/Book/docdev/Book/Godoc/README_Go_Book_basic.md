ขอรับผม! จัดการเรียบเรียง **คู่มือเรียน Go (12 บทแรก) + โครงสร้างโปรเจกต์ Golangiot** รวมเป็นเล่มเดียว ภาษาไทย พร้อมตัวอย่างโค้ด เรียบร้อยแล้วครับ

---

# คู่มือเรียนภาษา Go ฉบับสมบูรณ์ + โปรเจกต์ REST API ระดับ Production

## สารบัญ

**ส่วนที่ 1: พื้นฐาน Go**
- 01. เริ่มต้นกับ Go
- 02. พื้นฐานภาษา Go
- 03. Strings และการจัดรูปแบบ
- 04. การควบคุมการทำงาน
- 05. ฟังก์ชัน
- 06. Packages, Modules และ Workspaces
- 07. โครงสร้างข้อมูลหลัก
- 08. Structs และ Methods
- 09. พอยน์เตอร์
- 10. Interfaces (Polymorphism ของ Go)
- 11. การจัดการ Error
- 12. Generics (Go 1.18+)

**ส่วนที่ 2: การเขียนโปรแกรมแบบ Concurrent**
- 13. Goroutines และ Channels
- 14. รูปแบบการเขียน Concurrent ขั้นสูง

**ส่วนที่ 3: การจัดการไฟล์**
- 15. การจัดการไฟล์

**ส่วนที่ 4: โปรเจกต์ REST API ระดับ Production (Golangiot)**
- โครงสร้างโปรเจกต์
- การติดตั้งและรัน
- ฟีเจอร์หลัก
- การต่อ Database, Redis, Queue
- Middleware และ Security
- Deployment

---

# ส่วนที่ 1: พื้นฐาน Go

## 01. เริ่มต้นกับ Go

### Go คืออะไร?
Go (หรือ Golang) คือภาษาโปรแกรมมิ่งที่พัฒนาโดย Google ในปี 2007 เพื่อแก้ปัญหาการเขียนโปรแกรมยุคใหม่

### ทำไมต้อง Go?
- **ความเร็ว**: โค้ดทำงานเร็วเทียบเท่า C/C++
- **Concurrency**: รองรับงานพร้อมกันด้วย Goroutines
- **ความเรียบง่าย**: ไวยากรณ์กระชับ เรียนรู้ง่าย

### ติดตั้ง Go

**Windows:**
```bash
# ดาวน์โหลด installer จาก https://go.dev/dl/
# ติดตั้งที่ C:\Go\
go version
```

**Mac:**
```bash
brew install go
go version
```

**Linux:**
```bash
sudo apt-get install golang-go
# หรือดาวน์โหลด binary
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
```

### VS Code Setup
ติดตั้ง extensions:
- **Go** โดย Google
- **Go Test Explorer**

### รันโปรแกรม Go
```bash
# รันโดยตรง
go run main.go

# Build เป็น executable
go build main.go
./main  # หรือ main.exe (Windows)
```

---

## 02. พื้นฐานภาษา Go

### Hello World
```go
package main

import "fmt"

func main() {
    fmt.Println("สวัสดีโลก!")
}
```

### ตัวแปร (Variables)
```go
package main

import "fmt"

func main() {
    // แบบเต็มรูปแบบ
    var name string = "สมชาย"
    
    // แบบ推断ประเภท
    var age = 30
    
    // ประกาศโดยไม่กำหนดค่า (zero value)
    var salary float64  // = 0.0
    
    // Short declaration (ใช้ใน function เท่านั้น)
    city := "กรุงเทพฯ"
    
    // ประกาศหลายตัว
    var x, y int = 10, 20
    a, b := "hello", true
    
    fmt.Println(name, age, salary, city, x, y, a, b)
}
```

### ค่าคงที่ (Constants) และ iota
```go
package main

import "fmt"

func main() {
    const pi = 3.14159
    const appName = "MyApp"
    
    // iota - auto-increment
    const (
        Sunday = iota  // 0
        Monday         // 1
        Tuesday        // 2
    )
    
    fmt.Println(Sunday, Monday, Tuesday)
}
```

### ชนิดข้อมูล (Data Types)
```go
package main

import "fmt"

func main() {
    // Integer
    var i int = 42
    var i8 int8 = 127
    
    // Float
    var f32 float32 = 3.14
    var f64 float64 = 3.1415926535
    
    // String
    var s string = "Hello"
    
    // Boolean
    var b bool = true
    
    fmt.Printf("Type: %T, Value: %v\n", i, i)
}
```

### การแปลงชนิดข้อมูล (Type Conversion)
```go
package main

import "fmt"

func main() {
    var x int = 10
    var y float64 = float64(x)  // แปลง int -> float64
    
    var a float64 = 3.14
    var b int = int(a)  // ตัดเศษทิ้ง = 3
    
    fmt.Println(y, b)
}
```

### ตัวดำเนินการ (Operators)
```go
package main

import "fmt"

func main() {
    // Arithmetic
    a, b := 10, 3
    fmt.Println(a+b, a-b, a*b, a/b, a%b)
    
    // Comparison
    fmt.Println(a == b, a != b, a > b, a < b)
    
    // Logical
    x, y := true, false
    fmt.Println(x && y, x || y, !x)
}
```

---

## 03. Strings และการจัดรูปแบบ

### ฟังก์ชันจัดการ String
```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    s := "Hello, Go!"
    
    fmt.Println(len(s))              // ความยาว
    fmt.Println(strings.ToUpper(s))  // HELLO, GO!
    fmt.Println(strings.Contains(s, "Go"))    // true
    fmt.Println(strings.Split(s, ", "))  // [Hello Go!]
    fmt.Println(strings.Join([]string{"A", "B"}, "-"))  // A-B
}
```

### Rune vs Byte (สำคัญสำหรับภาษาไทย)
```go
package main

import "fmt"

func main() {
    text := "ภาษาไทย"
    
    // byte (ASCII, 1 byte) - ใช้ไม่ได้กับภาษาไทย
    fmt.Println("As bytes:")
    for i := 0; i < len(text); i++ {
        fmt.Printf("%c ", text[i])  // ผลลัพธ์ผิด
    }
    
    // rune (Unicode, 1-4 bytes) - ใช้กับภาษาไทยได้
    fmt.Println("\n\nAs runes:")
    for i, r := range text {
        fmt.Printf("Index: %d, Rune: %c, Unicode: %U\n", i, r, r)
    }
    
    // แปลง string to rune slice
    runes := []rune(text)
    fmt.Println("\nLength in runes:", len(runes))  // ภาษาไทย = 4 ตัว
}
```

### การจัดรูปแบบ String (Formatting)
```go
package main

import "fmt"

func main() {
    name := "สมชาย"
    age := 30
    score := 95.5
    
    // Printf (print formatted)
    fmt.Printf("ชื่อ: %s, อายุ: %d ปี\n", name, age)
    fmt.Printf("คะแนน: %.2f\n", score)
    
    // Sprintf (return string)
    message := fmt.Sprintf("สวัสดี %s คุณอายุ %d", name, age)
    fmt.Println(message)
    
    // Format verbs ที่ใช้บ่อย
    // %v - default format
    // %T - type
    // %d - integer
    // %s - string
    // %f - float
    // %t - boolean
}
```

### Multi-line Strings
```go
package main

import "fmt"

func main() {
    // Backtick (raw string) - ขึ้นบรรทัดใหม่ได้
    html := `
        <html>
            <body>
                <h1>Hello</h1>
            </body>
        </html>
    `
    
    fmt.Println(html)
}
```

---

## 04. การควบคุมการทำงาน

### if / else
```go
package main

import "fmt"

func main() {
    score := 75
    
    if score >= 80 {
        fmt.Println("Grade A")
    } else if score >= 70 {
        fmt.Println("Grade B")
    } else {
        fmt.Println("Grade C")
    }
    
    // Short statement in if
    if num := 10; num%2 == 0 {
        fmt.Println("Even number")
    }
}
```

### switch
```go
package main

import "fmt"

func main() {
    day := "Monday"
    
    switch day {
    case "Monday", "Tuesday":
        fmt.Println("Start of week")
    case "Friday":
        fmt.Println("TGIF")
    default:
        fmt.Println("Mid week")
    }
    
    // Type switch
    var x interface{} = "Hello"
    switch v := x.(type) {
    case int:
        fmt.Println("Integer:", v)
    case string:
        fmt.Println("String:", v)
    default:
        fmt.Println("Unknown type")
    }
}
```

### for loops
```go
package main

import "fmt"

func main() {
    // Traditional for loop
    for i := 0; i < 5; i++ {
        fmt.Println(i)
    }
    
    // While-style
    count := 0
    for count < 3 {
        fmt.Println("Count:", count)
        count++
    }
    
    // Range กับ slice
    numbers := []int{10, 20, 30}
    for index, value := range numbers {
        fmt.Printf("Index: %d, Value: %d\n", index, value)
    }
    
    // Range กับ map
    scores := map[string]int{"สมชาย": 85, "สมหญิง": 90}
    for name, score := range scores {
        fmt.Printf("%s ได้ %d คะแนน\n", name, score)
    }
}
```

### break, continue, defer
```go
package main

import "fmt"

func main() {
    // break - ออกจาก loop
    for i := 0; i < 10; i++ {
        if i == 5 {
            break
        }
        fmt.Println(i)
    }
    
    // continue - ข้ามรอบนั้น
    for i := 0; i < 5; i++ {
        if i%2 == 0 {
            continue
        }
        fmt.Println("Odd:", i)
    }
    
    // defer - ทำงานตอนจบ function (LIFO)
    defer fmt.Println("World")
    fmt.Println("Hello")
    // ผลลัพธ์: Hello World
}
```

---

## 05. ฟังก์ชัน

### ฟังก์ชันพื้นฐาน
```go
package main

import "fmt"

// Basic function
func greet() {
    fmt.Println("Hello!")
}

// Function with parameters and return
func add(x int, y int) int {
    return x + y
}

// Same type parameters
func multiply(x, y int) int {
    return x * y
}

func main() {
    greet()
    result := add(5, 3)
    fmt.Println(result)
}
```

### การคืนค่าหลายค่า (Multiple Returns)
```go
package main

import (
    "fmt"
    "errors"
)

func getPerson(id int) (string, int, error) {
    if id == 1 {
        return "สมชาย", 30, nil
    }
    return "", 0, errors.New("ไม่พบข้อมูล")
}

func main() {
    name, age, err := getPerson(1)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("ชื่อ: %s, อายุ: %d\n", name, age)
}
```

### Variadic Functions (รับพารามิเตอร์ไม่จำกัด)
```go
package main

import "fmt"

func sum(nums ...int) int {
    total := 0
    for _, num := range nums {
        total += num
    }
    return total
}

func main() {
    fmt.Println(sum(1, 2, 3))        // 6
    fmt.Println(sum(10, 20, 30, 40)) // 100
    
    // ส่ง slice
    numbers := []int{1, 2, 3, 4, 5}
    fmt.Println(sum(numbers...))
}
```

### Closure (ฟังก์ชันที่จดจำตัวแปร)
```go
package main

import "fmt"

func counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

func main() {
    c1 := counter()
    fmt.Println(c1()) // 1
    fmt.Println(c1()) // 2
    fmt.Println(c1()) // 3
    
    c2 := counter()
    fmt.Println(c2()) // 1 (แยก instance)
}
```

### Higher-order Functions
```go
package main

import "fmt"

// รับฟังก์ชันเป็นพารามิเตอร์
func applyOperation(x, y int, operation func(int, int) int) int {
    return operation(x, y)
}

// คืนค่าฟังก์ชัน
func getOperation(op string) func(int, int) int {
    switch op {
    case "add":
        return func(a, b int) int { return a + b }
    case "multiply":
        return func(a, b int) int { return a * b }
    default:
        return nil
    }
}

func main() {
    sum := applyOperation(5, 3, func(a, b int) int {
        return a + b
    })
    fmt.Println("Sum:", sum)
    
    multiply := getOperation("multiply")
    fmt.Println("Multiply:", multiply(4, 5))
}
```

---

## 06. Packages, Modules และ Workspaces

### สร้างโมดูลใหม่
```bash
go mod init myapp
```

### ไฟล์ go.mod
```go
module myapp

go 1.21

require (
    github.com/gin-gonic/gin v1.9.0
)
```

### การ Import Packages
```go
package main

import (
    "fmt"
    "math/rand"
    "strings"
    
    "myapp/utils"      // local package
    "github.com/gin-gonic/gin"  // external package
)

func main() {
    fmt.Println(rand.Intn(100))
    fmt.Println(strings.ToUpper("hello"))
}
```

### การสร้าง Package ของตัวเอง
```go
// utils/math.go
package utils

func Add(a, b int) int {
    return a + b
}

// main.go
package main

import (
    "fmt"
    "myapp/utils"
)

func main() {
    fmt.Println(utils.Add(5, 3))
}
```

### Go Workspace (Go 1.18+)
```bash
# สร้าง workspace สำหรับหลายโมดูล
mkdir myworkspace
cd myworkspace
go work init ./module1 ./module2
```

---

## 07. โครงสร้างข้อมูลหลัก

### Arrays (ขนาดคงที่)
```go
package main

import "fmt"

func main() {
    var arr1 [3]int                    // [0,0,0]
    arr2 := [3]int{1, 2, 3}           // [1,2,3]
    arr3 := [...]int{4, 5, 6, 7}      // infer size
    
    arr2[0] = 10
    fmt.Println(arr2[0])
    
    // Iterate
    for i, v := range arr2 {
        fmt.Printf("Index: %d, Value: %d\n", i, v)
    }
}
```

### Slices (ขนาดยืดหยุ่น)
```go
package main

import "fmt"

func main() {
    // สร้าง slice
    var s1 []int                    // nil slice
    s2 := []int{1, 2, 3}           // literal
    s3 := make([]int, 5)            // length 5, capacity 5
    s4 := make([]int, 3, 10)        // length 3, capacity 10
    
    // Append
    s2 = append(s2, 4, 5)           // [1,2,3,4,5]
    
    // Slicing operations
    slice := []int{1, 2, 3, 4, 5}
    fmt.Println(slice[1:4])          // [2,3,4]
    fmt.Println(slice[:3])           // [1,2,3]
    fmt.Println(slice[2:])           // [3,4,5]
    
    // Copy
    src := []int{1, 2, 3}
    dst := make([]int, len(src))
    copy(dst, src)
}
```

### Maps
```go
package main

import "fmt"

func main() {
    // สร้าง map
    m1 := make(map[string]int)      // empty map
    m2 := map[string]int{
        "apple": 5,
        "banana": 3,
    }
    
    // CRUD Operations
    m2["orange"] = 7
    m2["apple"] = 10
    
    // Read
    value := m2["apple"]
    value2, exists := m2["grape"]   // exists = false
    
    // Delete
    delete(m2, "banana")
    
    // Iterate
    for key, value := range m2 {
        fmt.Printf("%s: %d\n", key, value)
    }
}
```

### Deep Copy vs Shallow Copy
```go
package main

import "fmt"

func main() {
    // Shallow copy (ชี้ไปที่เดียวกัน)
    original := []int{1, 2, 3}
    shallow := original
    shallow[0] = 99
    fmt.Println(original)  // [99,2,3] เปลี่ยนด้วย
    
    // Deep copy (แยกกัน)
    src := []int{1, 2, 3}
    dst := make([]int, len(src))
    copy(dst, src)
    dst[0] = 99
    fmt.Println(src)  // [1,2,3] ไม่เปลี่ยน
    
    // Map deep copy
    originalMap := map[string]int{"a": 1, "b": 2}
    copyMap := make(map[string]int)
    for k, v := range originalMap {
        copyMap[k] = v
    }
}
```

---

## 08. Structs และ Methods

### Structs พื้นฐาน
```go
package main

import "fmt"

type Person struct {
    Name    string
    Age     int
    Address string
}

// Struct with tags (สำหรับ JSON)
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    p1 := Person{"สมชาย", 30, "กรุงเทพฯ"}
    p2 := Person{Name: "สมหญิง", Age: 25}
    
    fmt.Println(p1.Name)
    p2.Age = 26
}
```

### Embedding (การสืบทอดแบบ Go)
```go
package main

import "fmt"

type Animal struct {
    Name string
}

func (a Animal) Speak() {
    fmt.Println(a.Name, "makes sound")
}

type Dog struct {
    Animal  // Embedding
    Breed string
}

func (d Dog) Speak() {
    fmt.Println(d.Name, "barks!")
}

func main() {
    dog := Dog{
        Animal: Animal{Name: "Max"},
        Breed:  "Golden",
    }
    
    dog.Speak()           // Override
    dog.Animal.Speak()    // Call parent
}
```

### Methods (Value vs Pointer Receivers)
```go
package main

import "fmt"

type Counter struct {
    count int
}

// Value receiver (ไม่เปลี่ยนค่าเดิม)
func (c Counter) Value() int {
    return c.count
}

// Pointer receiver (เปลี่ยนค่าเดิม)
func (c *Counter) Increment() {
    c.count++
}

func main() {
    c := Counter{count: 0}
    
    c.Increment()
    fmt.Println(c.count)  // 1
}
```

### การแปลงเป็น JSON (Marshalling/Unmarshalling)
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
}

func main() {
    // struct -> JSON
    p := Product{ID: 1, Name: "Laptop", Price: 25000.50, InStock: true}
    jsonData, _ := json.Marshal(p)
    fmt.Println(string(jsonData))
    
    // JSON -> struct
    jsonString := `{"id":2,"name":"Mouse","price":500,"in_stock":false}`
    var product Product
    json.Unmarshal([]byte(jsonString), &product)
    fmt.Printf("%+v\n", product)
}
```

---

## 09. พอยน์เตอร์

### พื้นฐานพอยน์เตอร์
```go
package main

import "fmt"

func main() {
    x := 10
    p := &x  // p is pointer to x
    
    fmt.Println("Value of x:", x)      // 10
    fmt.Println("Address of x:", &x)   // 0xc0000140a0
    fmt.Println("Value pointed by p:", *p)  // 10 (dereference)
    
    // เปลี่ยนค่าผ่านพอยน์เตอร์
    *p = 20
    fmt.Println("New x:", x)  // 20
}
```

### Pointer Receivers vs Value Receivers
```go
package main

import "fmt"

type Rectangle struct {
    Width, Height float64
}

// Value receiver (ไม่เปลี่ยน struct เดิม)
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Pointer receiver (เปลี่ยน struct เดิม)
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    rect.Scale(2)
    fmt.Println(rect)  // {20, 10}
}
```

**เมื่อไรควรใช้ Pointer Receivers:**
1. ต้องการ modify receiver
2. struct มีขนาดใหญ่ (avoid copying)
3. ต้องการ nil value

---

## 10. Interfaces (Polymorphism ของ Go)

### Empty Interface (any)
```go
package main

import "fmt"

func printAnything(v interface{}) {
    fmt.Printf("Value: %v, Type: %T\n", v, v)
}

func main() {
    var anything interface{}
    
    anything = 42
    printAnything(anything)
    
    anything = "hello"
    printAnything(anything)
    
    // Go 1.18+ use 'any' alias
    var x any = "hello"
    fmt.Println(x)
}
```

### Type Assertions
```go
package main

import "fmt"

func main() {
    var i interface{} = "hello"
    
    // Safe assertion (with ok)
    s, ok := i.(string)
    if ok {
        fmt.Println("String:", s)
    }
    
    num, ok := i.(int)
    if !ok {
        fmt.Println("Not an integer")
    }
}
```

### การกำหนด Interface และ Implement
```go
package main

import "fmt"

// กำหนด interface
type Speaker interface {
    Speak() string
    GetName() string
}

// Dog implement Speaker
type Dog struct {
    Name string
}

func (d Dog) Speak() string {
    return "Woof!"
}

func (d Dog) GetName() string {
    return d.Name
}

// Cat implement Speaker
type Cat struct {
    Name string
}

func (c Cat) Speak() string {
    return "Meow!"
}

func (c Cat) GetName() string {
    return c.Name
}

// ฟังก์ชันที่รับ interface
func makeSound(s Speaker) {
    fmt.Printf("%s says: %s\n", s.GetName(), s.Speak())
}

func main() {
    var speaker Speaker
    
    speaker = Dog{Name: "Max"}
    makeSound(speaker)
    
    speaker = Cat{Name: "Luna"}
    makeSound(speaker)
}
```

### Dependency Injection Pattern
```go
package main

import "fmt"

// Repository interface
type UserRepository interface {
    GetUser(id int) (string, error)
}

// Service depends on interface
type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

// Concrete implementation
type InMemoryUserRepo struct {
    users map[int]string
}

func (r *InMemoryUserRepo) GetUser(id int) (string, error) {
    if name, exists := r.users[id]; exists {
        return name, nil
    }
    return "", fmt.Errorf("user not found")
}

func main() {
    repo := &InMemoryUserRepo{users: map[int]string{1: "John"}}
    service := NewUserService(repo)
    name, _ := service.repo.GetUser(1)
    fmt.Println("User:", name)
}
```

---

## 11. การจัดการ Error

### error Type พื้นฐาน
```go
package main

import (
    "errors"
    "fmt"
)

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func main() {
    result, err := divide(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Result:", result)
}
```

### Custom Errors
```go
package main

import (
    "fmt"
    "time"
)

type ValidationError struct {
    Field   string
    Value   interface{}
    Message string
    Time    time.Time
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("Validation error at %s: field '%s' - %s",
        e.Time.Format("2006-01-02 15:04:05"),
        e.Field,
        e.Message,
    )
}

func validateAge(age int) error {
    if age < 0 {
        return ValidationError{
            Field:   "age",
            Value:   age,
            Message: "age cannot be negative",
            Time:    time.Now(),
        }
    }
    return nil
}

func main() {
    err := validateAge(-5)
    if err != nil {
        fmt.Println(err)
    }
}
```

### Panic และ Recover
```go
package main

import "fmt"

func riskyOperation() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
        }
    }()
    
    panic("something went wrong")
}

func main() {
    riskyOperation()
    fmt.Println("Program continues after panic recovery")
}
```

### แนวปฏิบัติที่ดีในการจัดการ Error
```go
// 1. จัดการ error ทุกครั้ง (อย่าใช้ _)
// result, _ := divide(10, 0)  // DON'T DO THIS

// 2. Return errors early
func process(data string) error {
    if data == "" {
        return fmt.Errorf("empty data")
    }
    return nil
}

// 3. ใช้ custom error types
type NotFoundError struct {
    Resource string
    ID       int
}

// 4. Wrap errors with context
func fetchUser(id int) error {
    err := dbQuery(id)
    if err != nil {
        return fmt.Errorf("fetching user %d: %w", id, err)
    }
    return nil
}
```

---

## 12. Go Generics (Go 1.18+)

### Generic Functions พื้นฐาน
```go
package main

import "fmt"

// Basic generic function
func Print[T any](value T) {
    fmt.Println(value)
}

// Generic function with constraints
func Sum[T int | float64](numbers []T) T {
    var total T
    for _, n := range numbers {
        total += n
    }
    return total
}

func main() {
    Print(42)
    Print("hello")
    
    ints := []int{1, 2, 3, 4, 5}
    floats := []float64{1.1, 2.2, 3.3}
    
    fmt.Println("Sum ints:", Sum(ints))
    fmt.Println("Sum floats:", Sum(floats))
}
```

### Generic Structs
```go
package main

import "fmt"

type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}

func main() {
    intStack := Stack[int]{}
    intStack.Push(10)
    intStack.Push(20)
    
    val, _ := intStack.Pop()
    fmt.Println("Popped:", val)
    
    stringStack := Stack[string]{}
    stringStack.Push("hello")
    stringStack.Push("world")
}
```

### Constraints แบบกำหนดเอง
```go
package main

import "fmt"

type Number interface {
    ~int | ~int64 | ~float64 | ~float32
}

type Ordered interface {
    ~int | ~int64 | ~float64 | ~string
}

func Max[T Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

func main() {
    fmt.Println("Max:", Max(10, 20))
    fmt.Println("Max string:", Max("apple", "banana"))
}
```

---

# ส่วนที่ 2: การเขียนโปรแกรมแบบ Concurrent

## 13. Goroutines และ Channels

### Goroutines
```go
package main

import (
    "fmt"
    "time"
)

func printNumbers() {
    for i := 1; i <= 5; i++ {
        time.Sleep(100 * time.Millisecond)
        fmt.Printf("%d ", i)
    }
}

func main() {
    go printNumbers()
    go printNumbers()
    
    time.Sleep(2 * time.Second)
    fmt.Println("\nDone!")
}
```

### Unbuffered Channels
```go
package main

import "fmt"

func main() {
    ch := make(chan int)
    
    go func() {
        ch <- 42  // ส่งค่า
    }()
    
    value := <-ch  // รับค่า
    fmt.Println("Received:", value)
}
```

### Buffered Channels
```go
package main

import "fmt"

func main() {
    // Buffered channel (capacity 2)
    ch := make(chan string, 2)
    
    ch <- "first"
    ch <- "second"
    
    fmt.Println(<-ch) // first
    fmt.Println(<-ch) // second
}
```

### Select Statement
```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "from ch1"
    }()
    
    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "from ch2"
    }()
    
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println(msg1)
        case msg2 := <-ch2:
            fmt.Println(msg2)
        case <-time.After(3 * time.Second):
            fmt.Println("Timeout")
        }
    }
}
```

### WaitGroups
```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    fmt.Printf("Worker %d starting\n", id)
    time.Sleep(time.Second)
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup
    
    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }
    
    wg.Wait()
    fmt.Println("All workers completed")
}
```

### Mutex (ป้องกัน Race Condition)
```go
package main

import (
    "fmt"
    "sync"
)

type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func main() {
    counter := Counter{}
    var wg sync.WaitGroup
    
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }
    
    wg.Wait()
    fmt.Println("Counter value:", counter.value)
}
```

---

## 14. รูปแบบการเขียน Concurrent ขั้นสูง

### Worker Pools
```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Job struct {
    ID int
}

func worker(id int, jobs <-chan Job, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job.ID)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    const numJobs = 20
    const numWorkers = 5
    
    jobs := make(chan Job, numJobs)
    var wg sync.WaitGroup
    
    // Start workers
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, &wg)
    }
    
    // Send jobs
    for j := 1; j <= numJobs; j++ {
        jobs <- Job{ID: j}
    }
    close(jobs)
    
    wg.Wait()
}
```

### Pipelines
```go
package main

import "fmt"

func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

func main() {
    numbers := generate(1, 2, 3, 4, 5)
    squared := square(numbers)
    
    for result := range squared {
        fmt.Println(result)
    }
}
```

### Context สำหรับการยกเลิกการทำงาน
```go
package main

import (
    "context"
    "fmt"
    "time"
)

func worker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d stopped\n", id)
            return
        default:
            fmt.Printf("Worker %d working...\n", id)
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    
    go worker(ctx, 1)
    go worker(ctx, 2)
    
    time.Sleep(2 * time.Second)
    cancel() // ยกเลิกการทำงานทั้งหมด
    time.Sleep(1 * time.Second)
    
    // With timeout
    ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel2()
    
    select {
    case <-time.After(3 * time.Second):
        fmt.Println("Completed")
    case <-ctx2.Done():
        fmt.Println("Timeout:", ctx2.Err())
    }
}
```

---

# ส่วนที่ 3: การจัดการไฟล์

## 15. การจัดการไฟล์

### อ่าน/เขียนไฟล์
```go
package main

import (
    "fmt"
    "io/ioutil"
    "os"
)

func main() {
    // เขียนไฟล์
    data := []byte("Hello, Go!\nSecond line")
    err := ioutil.WriteFile("example.txt", data, 0644)
    if err != nil {
        panic(err)
    }
    
    // อ่านไฟล์
    content, err := ioutil.ReadFile("example.txt")
    if err != nil {
        panic(err)
    }
    fmt.Println(string(content))
    
    // ต่อท้ายไฟล์
    f, err := os.OpenFile("example.txt", os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    
    f.WriteString("\nAppended line")
}
```

### สร้าง/ลบไฟล์และไดเรกทอรี
```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // สร้างไฟล์
    file, err := os.Create("newfile.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    // ตรวจสอบว่าไฟล์มีอยู่หรือไม่
    if _, err := os.Stat("newfile.txt"); err == nil {
        fmt.Println("File exists")
    }
    
    // ลบไฟล์
    os.Remove("newfile.txt")
    
    // สร้างไดเรกทอรี
    os.Mkdir("mydir", 0755)
    os.MkdirAll("path/to/nested/dir", 0755)
    
    // ลบไดเรกทอรี
    os.Remove("mydir")
    os.RemoveAll("path")
}
```

### อ่านไฟล์ทีละบรรทัดด้วย bufio
```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    file, err := os.Open("large.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        fmt.Println(scanner.Text())
    }
}
```

### การทำงานกับ JSON/YAML
```go
package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
)

type Config struct {
    AppName string `json:"app_name"`
    Port    int    `json:"port"`
    Debug   bool   `json:"debug"`
}

func main() {
    config := Config{
        AppName: "MyApp",
        Port:    8080,
        Debug:   true,
    }
    
    // struct -> JSON
    jsonData, _ := json.MarshalIndent(config, "", "  ")
    ioutil.WriteFile("config.json", jsonData, 0644)
    
    // JSON -> struct
    var loadedConfig Config
    jsonFile, _ := ioutil.ReadFile("config.json")
    json.Unmarshal(jsonFile, &loadedConfig)
    fmt.Printf("%+v\n", loadedConfig)
}
```

---

# ส่วนที่ 4: โปรเจกต์ REST API ระดับ Production (Golangiot)

## ภาพรวมโปรเจกต์

Golangiot คือเทมเพลตสำหรับสร้าง REST API ด้วย Go ที่พร้อมใช้งานจริงในระดับ Production มีฟีเจอร์ครบครันสำหรับระบบ IoT และ Web Application ทั่วไป

## 🌟 ฟีเจอร์หลัก

### 🚀 Core Features
- **Clean Architecture** - สถาปัตยกรรม 3 ชั้น (Repository/Service/Handler)
- **JWT Authentication** - ระบบยืนยันตัวตนด้วย access/refresh tokens
- **User Management** - CRUD ผู้ใช้พร้อม Role-based protection
- **Structured Logging** - การบันทึก Log แบบมีโครงสร้างด้วย slog
- **Rate Limiting** - จำกัดจำนวนคำขอตาม IP
- **Health Monitoring** - ตรวจสอบสถานะระบบและ dependencies
- **Redis Cache** - ระบบแคชพร้อม TTL management
- **Message Queue** - Redis-based pub/sub พร้อม worker pools
- **Transaction Management** - GORM transaction manager

### 🛠️ Middleware Stack
- Request Context (Trace IDs, Request IDs)
- Security Headers (CSP, HSTS, XSS Protection)
- CORS Handling
- Panic Recovery
- Request Logging
- JWT Authentication
- Input Validation

## โครงสร้างไดเรกทอรี

```
project-root/
├── api/                          # เอกสาร API (Swagger)
│   └── app/
│       └── docs.go
├── cmd/                          # จุดเริ่มต้นโปรแกรม
│   └── app/
│       └── main.go
├── configs/                      # ไฟล์กำหนดค่า
│   ├── config.yaml
│   ├── config.example.yaml
│   └── config.production.yaml
├── deploy/                       # Deployment configurations
│   ├── docker/
│   │   ├── docker-compose.yaml
│   │   └── Dockerfile
│   └── k8s/
│       └── deployment.yaml
├── internal/                     # โค้ดภายในโปรเจกต์
│   ├── apps/                     # Application composition
│   │   └── app/
│   │       ├── bootstrap/        # DI และการเริ่มต้นระบบ
│   │       └── router/           # Route definitions
│   ├── core/                     # Domain modules หลัก
│   │   ├── auth/                 # ระบบยืนยันตัวตน
│   │   ├── user/                 # จัดการผู้ใช้
│   │   ├── iot/                  # ระบบ IoT (device, sensor)
│   │   ├── email/                # ส่งอีเมล
│   │   └── health/               # ตรวจสอบสุขภาพระบบ
│   ├── platform/                 # Infrastructure Layer
│   │   ├── config/               # โหลดการตั้งค่า
│   │   ├── db/                   # Database connections
│   │   │   ├── postgres/
│   │   │   └── redis/
│   │   ├── cache/                # Cache implementations
│   │   ├── queue/                # Message queue
│   │   └── logger/               # Logging
│   └── transport/                # HTTP layer
│       ├── middleware/           # Middleware ต่างๆ
│       ├── httpx/                # HTTP utilities
│       └── utils/
├── migrations/                   # Database migration files
├── pkg/                          # External packages
│   ├── cache/
│   ├── errors/
│   ├── jwt/
│   ├── logger/
│   ├── queue/
│   └── transaction/
├── scripts/                      # Utility scripts
├── test/                         # Integration tests
├── .air.toml                     # Hot-reload configuration
├── go.mod
└── README.md
```

## การติดตั้งและรันโปรเจกต์

### ขั้นตอนการติดตั้ง

```bash
# 1. โคลนโปรเจกต์
git clone https://github.com/your-repo/golangiot.git
cd golangiot

# 2. ติดตั้ง dependencies
go mod download

# 3. คัดลอกไฟล์ config
cp configs/config.example.yaml configs/config.yaml

# 4. แก้ไข config.yaml ให้ตรงกับการตั้งค่าของคุณ
# - database host, port, user, password
# - redis host, port
# - jwt secret key

# 5. รันด้วย Docker Compose (รวม PostgreSQL และ Redis)
cd deploy/docker
docker-compose up -d

# 6. รัน migration
go run cmd/app/main.go migrate

# 7. รันโปรแกรม
go run cmd/app/main.go

# หรือ build แล้วรัน
go build -o app cmd/app/main.go
./app
```

### การใช้ Air (Hot-reload สำหรับพัฒนา)

```bash
# ติดตั้ง air
go install github.com/cosmtrek/air@latest

# รันด้วย air
air
```

## ไฟล์กำหนดค่า (config.yaml)

```yaml
server:
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: golangiot
  sslmode: disable
  max_open_conns: 100
  max_idle_conns: 10

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "your-secret-key"
  access_token_ttl: 15m
  refresh_token_ttl: 168h  # 7 days

rate_limit:
  requests_per_second: 100
  burst: 200

log:
  level: "info"  # debug, info, warn, error
  format: "json"  # json or text
```

## API Endpoints ตัวอย่าง

### Public Routes
```http
POST   /api/v1/auth/register     # ลงทะเบียนผู้ใช้ใหม่
POST   /api/v1/auth/login        # เข้าสู่ระบบ
POST   /api/v1/auth/refresh      # รับ access token ใหม่
GET    /health                   # ตรวจสอบสุขภาพระบบ
GET    /health/detailed          # ตรวจสอบละเอียด (DB, Redis)
```

### Protected Routes (ต้องใช้ JWT)
```http
GET    /api/v1/users             # รายชื่อผู้ใช้
GET    /api/v1/users/:id         # ดูข้อมูลผู้ใช้
PUT    /api/v1/users/:id         # แก้ไขข้อมูลผู้ใช้
DELETE /api/v1/users/:id         # ลบผู้ใช้

# IoT endpoints
GET    /api/v1/devices           # รายการอุปกรณ์
POST   /api/v1/devices           # เพิ่มอุปกรณ์
GET    /api/v1/devices/:id/telemetry  # ดูข้อมูล telemetry
POST   /api/v1/devices/:id/command     # ส่งคำสั่งไปอุปกรณ์
```

## การเชื่อมต่อ Database

### PostgreSQL Connection
```go
// internal/platform/db/postgres/connection.go
package postgres

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func NewConnection(config Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
    
    return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
```

### Redis Connection
```go
// internal/platform/db/redis/connection.go
package redis

import (
    "github.com/go-redis/redis/v8"
)

func NewConnection(config Config) *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
        Password: config.Password,
        DB:       config.DB,
    })
}
```

## Middleware ตัวอย่าง

### JWT Authentication Middleware
```go
// internal/transport/middleware/authenticateJWT.go
func AuthenticateJWT(jwtService JWTService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                httpx.Error(w, http.StatusUnauthorized, "missing authorization header")
                return
            }
            
            tokenString := strings.TrimPrefix(authHeader, "Bearer ")
            claims, err := jwtService.ValidateToken(tokenString)
            if err != nil {
                httpx.Error(w, http.StatusUnauthorized, "invalid token")
                return
            }
            
            ctx := context.WithValue(r.Context(), "user", claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### Rate Limiting Middleware
```go
// internal/transport/middleware/rate_limit.go
func RateLimit(limiter *rate.Limiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                httpx.Error(w, http.StatusTooManyRequests, "too many requests")
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

## การใช้งาน Redis Cache

```go
// pkg/cache/redis_cache.go
type RedisCache struct {
    client *redis.Client
    ttl    time.Duration
}

func (c *RedisCache) Set(key string, value interface{}) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return c.client.Set(context.Background(), key, data, c.ttl).Err()
}

func (c *RedisCache) Get(key string, dest interface{}) error {
    data, err := c.client.Get(context.Background(), key).Bytes()
    if err != nil {
        return err
    }
    return json.Unmarshal(data, dest)
}
```

## การใช้งาน Message Queue

```go
// pkg/queue/redis_queue.go
type RedisQueue struct {
    client *redis.Client
}

func (q *RedisQueue) Publish(queue string, message interface{}) error {
    data, err := json.Marshal(message)
    if err != nil {
        return err
    }
    return q.client.LPush(context.Background(), queue, data).Err()
}

func (q *RedisQueue) Subscribe(queue string, handler func([]byte) error) {
    for {
        result, err := q.client.BRPop(context.Background(), 0, queue).Result()
        if err != nil {
            continue
        }
        
        if err := handler([]byte(result[1])); err != nil {
            // จัดการ error (อาจส่งไป dead letter queue)
            log.Printf("Error processing message: %v", err)
        }
    }
}
```

## การ Deploy

### Docker
```dockerfile
# deploy/docker/Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app cmd/app/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/app .
COPY --from=builder /app/configs ./configs
EXPOSE 8080
CMD ["./app"]
```

### Docker Compose
```yaml
# deploy/docker/docker-compose.yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
    environment:
      - DATABASE_HOST=postgres
      - REDIS_HOST=redis

  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: golangiot
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

### Kubernetes Deployment
```yaml
# deploy/k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: golangiot
spec:
  replicas: 3
  selector:
    matchLabels:
      app: golangiot
  template:
    metadata:
      labels:
        app: golangiot
    spec:
      containers:
      - name: app
        image: golangiot:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_HOST
          value: "postgres-service"
        - name: REDIS_HOST
          value: "redis-service"
---
apiVersion: v1
kind: Service
metadata:
  name: golangiot-service
spec:
  selector:
    app: golangiot
  ports:
  - port: 8080
    targetPort: 8080
  type: LoadBalancer
```

## คำสั่งที่มีประโยชน์

```bash
# รัน tests
go test ./...

# รัน tests with coverage
go test -cover ./...

# Build สำหรับ Linux
GOOS=linux GOARCH=amd64 go build -o app-linux cmd/app/main.go

# Build สำหรับ Windows
GOOS=windows GOARCH=amd64 go build -o app.exe cmd/app/main.go

# ตรวจสอบ race condition
go run -race cmd/app/main.go

# ดู dependencies
go mod graph

# อัปเดต dependencies
go get -u ./...
go mod tidy

# สร้าง Swagger docs
swag init -g cmd/app/main.go
```

---

## สรุป

คู่มือนี้ครอบคลุม:

1. **พื้นฐาน Go** - ตั้งแต่การติดตั้ง, ไวยากรณ์, ไปจนถึง Generics
2. **Concurrency** - Goroutines, Channels, Select, WaitGroups, Mutex
3. **การจัดการไฟล์** - อ่าน/เขียน, JSON/YAML
4. **โปรเจกต์จริง Golangiot** - REST API พร้อมใช้ระดับ Production

โปรเจกต์ Golangiot พร้อมให้คุณนำไปพัฒนาต่อยอดเป็นระบบ IoT, Web Application, หรือ Microservices ได้ทันที

---
 