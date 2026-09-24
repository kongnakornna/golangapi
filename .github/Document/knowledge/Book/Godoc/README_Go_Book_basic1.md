# เอกสารประกอบการเรียนรู้ภาษา Go (Golang)
## หัวข้อ: การพัฒนาโปรแกรมด้วย Go ตั้งแต่พื้นฐานจนถึงเว็บแอปพลิเคชัน

---

## 1. บทนำ

ภาษา Go หรือ Golang เป็นภาษาโปรแกรมมิ่งที่พัฒนาโดย Google ในปี 2007 และเปิดตัวสู่สาธารณะในปี 2009 ถูกออกแบบโดย Robert Griesemer, Rob Pike และ Ken Thompson เพื่อแก้ไขข้อจำกัดของภาษาexisting ภาษา โดยมุ่งเน้นที่ความเรียบง่าย ประสิทธิภาพสูง และการรองรับการทำงานแบบ concurrent ได้อย่างยอดเยี่ยม

Go ได้รับความนิยมอย่างรวดเร็วในวงการพัฒนา ซอฟต์แวร์เนื่องจากมีจุดเด่นที่สำคัญ:
- **ความเรียบง่าย**: ไวยากรณ์กระชับ เรียนรู้ง่าย
- **ประสิทธิภาพสูง**: Compiled language ที่ทำงานเร็วเทียบเท่า C/C++
- **Concurrency แบบเนทีฟ**: Goroutine และ Channel ทำให้การเขียนโปรแกรม concurrent ง่ายขึ้น
- **Garbage Collection**: จัดการหน่วยความจำอัตโนมัติ
- **Standard Library ที่ทรงพลัง**: มี package ให้ใช้งานครอบคลุมความต้องการ

---

## 2. บทนิยาม

| คำศัพท์ | นิยาม |
|--------|-------|
| **Go (Golang)** | ภาษาโปรแกรมมิ่งที่พัฒนาโดย Google มีจุดเด่นด้านประสิทธิภาพและความเรียบง่าย |
| **Goroutine** | หน่วยการทำงานย่อยที่เบา (lightweight thread) จัดการโดย Go runtime |
| **Channel** | กลไกในการสื่อสารระหว่าง goroutine แบบ type-safe |
| **Package** | กลุ่มของโค้ด Go ที่รวมกันเพื่อให้สามารถนำไป reuse ได้ |
| **Module** | กลุ่มของ package ที่ versioned together สำหรับการจัดการ dependency |
| **Struct** | ชนิดข้อมูลแบบ composite ที่รวม fields หลายๆ type เข้าด้วยกัน |
| **Interface** | ชุดของ method signatures ที่กำหนดพฤติกรรมของ type |
| **Pointer** | ตัวแปรที่เก็บ memory address ของตัวแปรอื่น |
| **Slice** | Dynamic array ที่สามารถขยายขนาดได้ |
| **Map** | ข้อมูลแบบ key-value store |
| **Defer** | คำสั่งที่ใช้หน่วงเวลาการทำงานของ function จนกว่า function ที่เรียกจะจบการทำงาน |
| **Go Module** | ระบบจัดการ dependency อย่างเป็นทางการของ Go |

---

## 3. บทหัวข้อ: เนื้อหาการเรียนรู้แบ่งตามส่วน

### ส่วนที่ 1: Introduction to Go (11 บทเรียน)
- แนะนำภาษา Go
- การติดตั้ง Go และ VS Code
- โปรแกรมแรก Hello, World!
- โครงสร้างโปรแกรม Go
- ตัวแปรและ Dot Notation
- การรันโปรแกรม Eliza

### ส่วนที่ 2: Variables and Scope (12 บทเรียน)
- การประกาศ和使用ตัวแปร
- เกมทายตัวเลข
- ขอบเขตการใช้งานของตัวแปร (Scope)

### ส่วนที่ 3: Console Input and Strings (10 บทเรียน)
- การรับข้อมูลผ่าน Console
- การตรวจจับการกดแป้นพิมพ์
- String Interpolation

### ส่วนที่ 4: Types, Pointers, and Data Structures (18 บทเรียน)
- ชนิดข้อมูลพื้นฐานและชนิดข้อมูลรวม
- Pointers, Slices, Maps
- Functions, Channels, Interfaces
- Composition และ Exported vs Unexported

### ส่วนที่ 5: Loops and Control Flow (15 บทเรียน)
- For loop แบบต่างๆ
- While loop และ Infinite loop
- การใช้ Debugger

### ส่วนที่ 6: Conditionals and Select (11 บทเรียน)
- if-else statements
- switch statement
- select statement สำหรับ channels

### ส่วนที่ 7: Operators (10 บทเรียน)
- โอเปอเรเตอร์และลำดับความสำคัญ
- Modulus, Relational, Conditional Operators
- Short Circuit Evaluation

### ส่วนที่ 8: Strings (9 บทเรียน)
- การจัดการ string
- strings package
- การแปลงตัวพิมพ์

### ส่วนที่ 9: Web Application Development (14 บทเรียน)
- การสร้างเว็บแอปพลิเคชันด้วย Go
- Serving HTML และ JSON
- การพัฒนาเกมเป่ายิ้งฉุบบนเว็บ

---

## 4. โครงสร้างโปรแกรม Go

### 4.1 โครงสร้างพื้นฐานของโปรแกรม Go

```go
// 1. package declaration - ทุกไฟล์ Go ต้องมี
package main

// 2. import statements - นำเข้า packages ที่ต้องการใช้
import (
    "fmt"
    "strings"
)

// 3. constant declarations
const Pi = 3.14159

// 4. variable declarations
var message string = "Hello, Go!"

// 5. type declarations
type Person struct {
    Name string
    Age  int
}

// 6. function declarations
func main() {
    // 7. function body - จุดเริ่มต้นของโปรแกรม
    fmt.Println(message)
    
    // 8. calling other functions
    greet("World")
}

func greet(name string) {
    fmt.Printf("Hello, %s!\n", name)
}
```

### 4.2 ตัวอย่างโปรแกรมแรก: Hello, World!

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### 4.3 โครงสร้างโปรแกรมที่ซับซ้อนขึ้น

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

// ประกาศ struct สำหรับเก็บข้อมูล
type Student struct {
    ID   int
    Name string
    GPA  float64
}

// method ของ struct
func (s Student) Display() {
    fmt.Printf("ID: %d, Name: %s, GPA: %.2f\n", s.ID, s.Name, s.GPA)
}

// interface definition
type Greeter interface {
    Greet() string
}

// implement interface
func (s Student) Greet() string {
    return fmt.Sprintf("Hello, I'm %s", s.Name)
}

func main() {
    // การใช้งานตัวแปร
    var name string = "John Doe"
    age := 25  // short declaration
    
    fmt.Printf("Name: %s, Age: %d\n", name, age)
    
    // การใช้งาน array และ slice
    numbers := []int{1, 2, 3, 4, 5}
    for i, num := range numbers {
        fmt.Printf("numbers[%d] = %d\n", i, num)
    }
    
    // การใช้งาน map
    scores := map[string]int{
        "Math":    90,
        "Science": 85,
        "English": 88,
    }
    
    for subject, score := range scores {
        fmt.Printf("%s: %d\n", subject, score)
    }
    
    // การใช้งาน struct
    student := Student{
        ID:   1001,
        Name: "Alice Wonderland",
        GPA:  3.75,
    }
    student.Display()
    
    // การใช้ interface
    var g Greeter = student
    fmt.Println(g.Greet())
    
    // การรับ input จาก console
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter your name: ")
    input, _ := reader.ReadString('\n')
    input = strings.TrimSpace(input)
    fmt.Printf("Hello, %s!\n", input)
}
```

---

## 5. Workflow การพัฒนาโปรแกรมด้วย Go

### 5.1 Workflow ทั่วไป

```
┌─────────────────────────────────────────────────────────────────┐
│                     Workflow การพัฒนา Go                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐ │
│  │ 1. Setup │───▶│ 2. Write │───▶│ 3. Build │───▶│ 4. Test  │ │
│  │  环境设置  │    │   โค้ด    │    │  คอมไพล์  │    │   ทดสอบ  │ │
│  └──────────┘    └──────────┘    └──────────┘    └──────────┘ │
│       │              │              │              │           │
│       ▼              ▼              ▼              ▼           │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐ │
│  │ 5. Debug │───▶│ 6. Run   │───▶│ 7. Deploy│───▶│ 8. Maintain│ │
│  │  แก้ไขบั๊ก  │    │  รันโปรแกรม│    │  นำไปใช้  │    │   บำรุงรักษา│ │
│  └──────────┘    └──────────┘    └──────────┘    └──────────┘ │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 5.2 ขั้นตอนการพัฒนาโดยละเอียด

#### ขั้นตอนที่ 1: การติดตั้งและตั้งค่า Environment
```bash
# ดาวน์โหลดและติดตั้ง Go จาก https://golang.org/dl/
# ตรวจสอบการติดตั้ง
go version

# ตั้งค่า Go Module
go mod init myproject

# ติดตั้ง VS Code Extension
# - Go extension โดย Google
```

#### ขั้นตอนที่ 2: การเขียนโค้ด
```go
// main.go
package main

import "fmt"

func main() {
    fmt.Println("Starting application...")
}
```

#### ขั้นตอนที่ 3: การ Build และ Run
```bash
# รันโดยตรง
go run main.go

# Build เป็น executable
go build -o myapp main.go

# Build สำหรับ platform อื่น
GOOS=linux GOARCH=amd64 go build -o myapp-linux main.go
GOOS=windows GOARCH=amd64 go build -o myapp.exe main.go
GOOS=darwin GOARCH=amd64 go build -o myapp-mac main.go
```

#### ขั้นตอนที่ 4: การ Testing
```go
// main_test.go
package main

import "testing"

func TestGreeting(t *testing.T) {
    result := greeting("Go")
    expected := "Hello, Go!"
    if result != expected {
        t.Errorf("Expected %s but got %s", expected, result)
    }
}

// รัน test
// go test -v
```

#### ขั้นตอนที่ 5: การ Debugging
```bash
# ใช้ Delve debugger
dlv debug main.go

# หรือใช้ VS Code debugger
# สร้างไฟล์ .vscode/launch.json
```

#### ขั้นตอนที่ 6: การ Deploy
```dockerfile
# Dockerfile สำหรับ Go application
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o myapp .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/myapp .
CMD ["./myapp"]
```

---

## 6. Case Study: เกมทายตัวเลข (Guess the Number Game)

### 6.1 โจทย์และ Requirement
พัฒนาเกมทายตัวเลขที่คอมพิวเตอร์สุ่มตัวเลข 1-100 ให้ผู้เล่นทาย พร้อมให้คำแนะนำว่าสูงหรือต่ำกว่า

### 6.2 การออกแบบ (Design)

```
┌─────────────────────────────────────────────────────────────┐
│                   Guess the Number Game                      │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────┐      ┌─────────────────────────────────┐  │
│  │   Start     │─────▶│  Generate random number 1-100  │  │
│  └─────────────┘      └─────────────────────────────────┘  │
│                              │                             │
│                              ▼                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              While guess != secret                   │   │
│  │  ┌─────────────────────────────────────────────┐    │   │
│  │  │ 1. Prompt user for guess                     │    │   │
│  │  │ 2. Read input                                │    │   │
│  │  │ 3. Compare guess with secret                 │    │   │
│  │  │ 4. Give hint (higher/lower)                  │    │   │
│  │  │ 5. Increment attempts counter                │    │   │
│  │  └─────────────────────────────────────────────┘    │   │
│  └─────────────────────────────────────────────────────┘   │
│                              │                             │
│                              ▼                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │           Display "Congratulations!"                 │   │
│  │           Show number of attempts                    │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 6.3 โค้ดตัวอย่าง

```go
// guess_game.go
package main

import (
    "bufio"
    "fmt"
    "math/rand"
    "os"
    "strconv"
    "strings"
    "time"
)

// Game struct เก็บสถานะของเกม
type Game struct {
    secretNumber int
    maxAttempts  int
    attempts     int
    isGameOver   bool
}

// NewGame สร้างเกมใหม่
func NewGame(max int) *Game {
    rand.Seed(time.Now().UnixNano())
    return &Game{
        secretNumber: rand.Intn(100) + 1,
        maxAttempts:  max,
        attempts:     0,
        isGameOver:   false,
    }
}

// Play ดำเนินการเล่นเกม
func (g *Game) Play() {
    fmt.Println("===================================")
    fmt.Println("     Guess the Number Game!")
    fmt.Println("===================================")
    fmt.Printf("I'm thinking of a number between 1 and 100\n")
    fmt.Printf("You have %d attempts. Good luck!\n\n", g.maxAttempts)
    
    reader := bufio.NewReader(os.Stdin)
    
    for !g.isGameOver && g.attempts < g.maxAttempts {
        remaining := g.maxAttempts - g.attempts
        fmt.Printf("Attempts remaining: %d\n", remaining)
        fmt.Print("Enter your guess: ")
        
        input, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error reading input. Please try again.")
            continue
        }
        
        input = strings.TrimSpace(input)
        guess, err := strconv.Atoi(input)
        
        if err != nil {
            fmt.Println("Invalid input! Please enter a number.\n")
            continue
        }
        
        g.attempts++
        
        // ตรวจสอบคำตอบ
        if guess < g.secretNumber {
            fmt.Println("Too low! Try again.\n")
        } else if guess > g.secretNumber {
            fmt.Println("Too high! Try again.\n")
        } else {
            fmt.Println("\n===================================")
            fmt.Printf("🎉 Congratulations! You guessed it! 🎉\n")
            fmt.Printf("The number was %d\n", g.secretNumber)
            fmt.Printf("You took %d attempts\n", g.attempts)
            fmt.Println("===================================")
            g.isGameOver = true
        }
    }
    
    if !g.isGameOver {
        fmt.Println("\n===================================")
        fmt.Printf("😢 Game Over! The number was %d\n", g.secretNumber)
        fmt.Println("Better luck next time!")
        fmt.Println("===================================")
    }
}

// Reset เริ่มเกมใหม่
func (g *Game) Reset() {
    rand.Seed(time.Now().UnixNano())
    g.secretNumber = rand.Intn(100) + 1
    g.attempts = 0
    g.isGameOver = false
}

func main() {
    game := NewGame(10)
    game.Play()
    
    // ถามว่าจะเล่นอีกไหม
    fmt.Print("\nPlay again? (y/n): ")
    reader := bufio.NewReader(os.Stdin)
    answer, _ := reader.ReadString('\n')
    answer = strings.TrimSpace(strings.ToLower(answer))
    
    if answer == "y" || answer == "yes" {
        game.Reset()
        game.Play()
    }
    
    fmt.Println("Thanks for playing!")
}
```

### 6.4 การรันและทดสอบ

```bash
# รันเกม
go run guess_game.go

# ตัวอย่างผลลัพธ์
===================================
     Guess the Number Game!
===================================
I'm thinking of a number between 1 and 100
You have 10 attempts. Good luck!

Attempts remaining: 10
Enter your guess: 50
Too low! Try again.

Attempts remaining: 9
Enter your guess: 75
Too high! Try again.

Attempts remaining: 8
Enter your guess: 62
Too low! Try again.

Attempts remaining: 7
Enter your guess: 68

===================================
🎉 Congratulations! You guessed it! 🎉
The number was 68
You took 4 attempts
===================================
```

### 6.5 การทดสอบด้วย Unit Test

```go
// guess_game_test.go
package main

import (
    "testing"
)

func TestNewGame(t *testing.T) {
    game := NewGame(10)
    
    if game.secretNumber < 1 || game.secretNumber > 100 {
        t.Errorf("Secret number %d is out of range", game.secretNumber)
    }
    
    if game.maxAttempts != 10 {
        t.Errorf("Expected maxAttempts=10, got %d", game.maxAttempts)
    }
    
    if game.attempts != 0 {
        t.Errorf("Expected attempts=0, got %d", game.attempts)
    }
}

func TestGameReset(t *testing.T) {
    game := NewGame(10)
    oldSecret := game.secretNumber
    game.attempts = 5
    game.isGameOver = true
    
    game.Reset()
    
    if game.secretNumber == oldSecret {
        t.Error("Secret number should change after reset")
    }
    
    if game.attempts != 0 {
        t.Errorf("Expected attempts=0 after reset, got %d", game.attempts)
    }
    
    if game.isGameOver {
        t.Error("isGameOver should be false after reset")
    }
}

// Benchmark test
func BenchmarkNewGame(b *testing.B) {
    for i := 0; i < b.N; i++ {
        NewGame(10)
    }
}
```

---

## สรุป

ภาษา Go เป็นภาษาโปรแกรมมิ่งที่ทรงพลัง มีความเรียบง่ายและประสิทธิภาพสูง เหมาะสำหรับการพัฒนา:
- **Web Applications**: ด้วย net/http package ที่มีประสิทธิภาพ
- **Microservices**: ขนาดเล็ก ทำงานเร็ว รองรับ concurrent
- **CLI Tools**: compile เป็น standalone executable
- **Cloud & DevOps**: เครื่องมือยอดนิยมอย่าง Docker, Kubernetes เขียนด้วย Go

การเรียนภาษา Go ควรเริ่มจากพื้นฐานที่ถูกต้อง ทำความเข้าใจโครงสร้าง ตัวแปร ชนิดข้อมูล การควบคุม flow และค่อยๆ พัฒนาสู่การเขียนโปรแกรม concurrent และ web application เพื่อให้สามารถนำไปประยุกต์ใช้งานจริงได้อย่างมีประสิทธิภาพ