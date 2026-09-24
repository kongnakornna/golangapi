# 📚 Go Programming Mastery: Function 

> **สารบัญ (Table of Contents)**
> - [บทที่ 1: รู้จักกับภาษา Go](#บทที่-1-รู้จักกับภาษา-go)
> - [บทที่ 2: ตัวแปรและชนิดข้อมูล](#บทที่-2-ตัวแปรและชนิดข้อมูล)
> - [บทที่ 3: โครงสร้างควบคุม (Control Structures)](#บทที่-3-โครงสร้างควบคุม-control-structures)
> - [บทที่ 4: ฟังก์ชัน (Functions)](#บทที่-4-ฟังก์ชัน-functions)
> - [บทที่ 5: Array, Slice, และ Map](#บทที่-5-array-slice-และ-map)
> - [บทที่ 6: Struct และ Interface](#บทที่-6-struct-และ-interface)
> - [บทที่ 7: Concurrency (Goroutine & Channel)](#บทที่-7-concurrency-goroutine--channel)
> - [บทที่ 8: Error Handling](#บทที่-8-error-handling)
> - [บทที่ 9: Web API Development](#บทที่-9-web-api-development)
> - [บทที่ 10: Database Operations](#บทที่-10-database-operations)
> - [บทที่ 11: Middleware และ Rate Limiting](#บทที่-11-middleware-และ-rate-limiting)
> - [บทที่ 12: Testing และ Benchmarking](#บทที่-12-testing-และ-benchmarking)

---

# บทที่ 1: รู้จักกับภาษา Go

## 📌 สรุปสั้น (Executive Summary)

**Go** (หรือ Golang) เป็นภาษาโปรแกรมที่พัฒนาโดย Google ในปี 2007 และเปิดตัวในปี 2009 ออกแบบมาเพื่อแก้ปัญหาความซับซ้อนของภาษา C++ และ Java ในขณะที่ยังคงประสิทธิภาพสูง รองรับการทำงานพร้อมกัน (Concurrency) ได้ดีเยี่ยม

---

## 1.1 โครงสร้างการทำงาน (Architecture)

```
┌─────────────────────────────────────────────────────────────┐
│                    Go Programming Language                  │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   Compiler  │  │  Runtime    │  │  Garbage    │         │
│  │   (gc)      │  │  Scheduler  │  │  Collector  │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Standard Library (stdlib)               │   │
│  │  • net/http  • database/sql  • encoding/json        │   │
│  │  • crypto    • sync          • testing              │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Go Runtime Features                     │   │
│  │  • Goroutines (lightweight threads)                 │   │
│  │  • Channels (communication between goroutines)      │   │
│  │  • Defer (cleanup operations)                       │   │
│  │  • Panic/Recover (error handling)                   │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

## 1.2 วัตถุประสงค์ (Objectives)

| วัตถุประสงค์ | คำอธิบาย |
|-------------|----------|
| **ประสิทธิภาพสูง** | รวบรวมเป็น Binary File ที่ทำงานเร็วเทียบเท่า C/C++ |
| **พัฒนาง่าย** | ไวยากรณ์(clean syntax) เรียนรู้ง่าย |
| **Concurrency** | รองรับการทำงานพร้อมกันผ่าน Goroutine |
| **Scalability** | ขยายระบบได้ง่าย รองรับ Traffic สูง |
| **Cross-Platform** | รองรับ Windows, Linux, macOS, Docker |

---

## 1.3 กลุ่มเป้าหมาย (Target Audience)

| กลุ่ม | เหมาะสม | เหตุผล |
|------|---------|--------|
| **นักพัฒนาระบบคลาวด์** | ⭐⭐⭐⭐⭐ | Docker, Kubernetes เขียนด้วย Go |
| **นักพัฒนา Backend/API** | ⭐⭐⭐⭐⭐ | net/http แข็งแกร่ง, Performance ดี |
| **DevOps Engineer** | ⭐⭐⭐⭐⭐ | เครื่องมือ CLI, Infrastructure as Code |
| **นักพัฒนา Microservices** | ⭐⭐⭐⭐⭐ | ไฟล์ Binary เล็ก, Start-up เร็ว |
| **นักพัฒนาเริ่มต้น** | ⭐⭐⭐⭐ | ไวยากรณ์สั้น เรียนรู้ง่าย |
| **นักพัฒนา Mobile** | ⭐⭐ | มี Golang Mobile แต่ยังไม่แพร่หลาย |

---

## 1.4 ความรู้พื้นฐาน (Prerequisites)

```
ความรู้ที่ควรมีก่อนเรียน Go:
┌─────────────────────────────────────────────────────────────┐
│  ✅ พื้นฐานการเขียนโปรแกรม (Programming Fundamentals)       │
│     • ตัวแปร (Variables)                                    │
│     • เงื่อนไข (If-Else, Switch)                            │
│     • ลูป (Loops)                                           │
│     • ฟังก์ชัน (Functions)                                   │
├─────────────────────────────────────────────────────────────┤
│  ✅ ความเข้าใจเรื่อง Data Types                               │
│     • Integer, Float, String, Boolean                       │
│     • Array, List, Dictionary                               │
├─────────────────────────────────────────────────────────────┤
│  ✅ ความเข้าใจเรื่อง Pointer (เบื้องต้น)                      │
├─────────────────────────────────────────────────────────────┤
│  🔧 (ไม่จำเป็น) OOP เนื่องจาก Go ไม่มี Class แต่มี Struct     │
└─────────────────────────────────────────────────────────────┘
```

---

## 1.5 เนื้อหาโดยย่อ (Content Summary)

### 1.5.1 วัตถุประสงค์หลักของการใช้ Go

| การใช้งาน | ประโยชน์ | ตัวอย่างระบบ |
|-----------|----------|--------------|
| **Web Server/API** | รองรับ Concurrent สูง | Uber, Dropbox API |
| **CLI Tools** | Build ได้ไฟล์เดียว พกพาสะดวก | Docker CLI, kubectl |
| **Cloud Services** | ทำงานบน Container ได้ดี | Kubernetes, Terraform |
| **Networking** | net/http, websocket ทรงพลัง | Caddy, Traefik |
| **DevOps Tools** | เขียนได้เร็ว ทำงานไว | Prometheus, Grafana |

### 1.5.2 จุดเด่นของ Go (Advantages)

```
┌─────────────────────────────────────────────────────────────────┐
│  📌 จุดเด่นที่ทำให้ Go แตกต่างจากภาษาอื่น                         │
├─────────────────────────────────────────────────────────────────┤
│  1. Goroutine: ลูทน้ำหนักเบา เริ่มต้นใช้ RAM แค่ 2KB             │
│     ต่างจาก Thread ของ Java ที่ใช้ 1MB                          │
│                                                                  │
│  2. Channel: การสื่อสารระหว่าง Goroutine ที่ปลอดภัย              │
│     "Don't communicate by sharing memory; share memory by       │
│      communicating."                                            │
│                                                                  │
│  3. Garbage Collection: มีตัวเก็บขยะที่รวดเร็ว                    │
│                                                                  │
│  4. Fast Compilation: คอมไพล์เร็วมาก ต่างจาก C++/Java            │
│                                                                  │
│  5. Static Typing: ตรวจสอบชนิดข้อมูลตอนคอมไพล์ ลด Error         │
│                                                                  │
│  6. Built-in Concurrency: ไม่ต้องใช้ Library ภายนอก              │
└─────────────────────────────────────────────────────────────────┘
```

---

## 1.6 บทนำ (Introduction)

ภาษา Go ถูกสร้างขึ้นโดย **Robert Griesemer**, **Rob Pike**, และ **Ken Thompson** ที่ Google โดยมีแรงบันดาลใจจากการทำงานกับระบบขนาดใหญ่ที่ต้องรองรับการทำงานพร้อมกันนับล้านๆ ครั้ง ปัญหาที่พบในภาษา C++ และ Java คือ:

1. **Compilation ช้า** - รอคอมไพล์นานหลายนาที
2. **Concurrency ยาก** - การเขียนโปรแกรมแบบ Concurrent ซับซ้อน
3. **Dependency ยุ่งยาก** - การจัดการ Library หลายเวอร์ชัน
4. **Deployment ใหญ่** - ต้องติดตั้ง Runtime Environment

**Go จึงถูกออกแบบมาให้มี:**
- ไวยากรณ์ที่กระชับ (25 keywords เท่านั้น)
- รวบรวมเป็นไฟล์ Binary เดียว (Static Linking)
- Goroutine สำหรับ Concurrent Programming
- Built-in Testing และ Benchmarking
- Garbage Collection อัตโนมัติ

---

## 1.7 บทนิยาม (Definitions)

| ศัพท์เทคนิค | คำอธิบาย | ตัวอย่าง |
|------------|----------|----------|
| **Goroutine** | ฟังก์ชันที่ทำงานพร้อมกันแบบ lightweight (น้ำหนักเบา) | `go myFunction()` |
| **Channel** | ท่อส่งข้อมูลระหว่าง Goroutine | `ch := make(chan int)` |
| **Defer** | กำหนดให้ฟังก์ชันทำงานเมื่อฟังก์ชันหลักจบ | `defer file.Close()` |
| **Package** | กลุ่มของไฟล์ Go ที่ทำงานร่วมกัน | `package main` |
| **Module** | กลุ่มของ Packages ที่มี version ร่วมกัน | `go mod init myproject` |
| **Interface** | ชุดของ method signatures | `type Writer interface { Write([]byte) error }` |
| **Struct** | การรวมข้อมูลหลายชนิดเข้าด้วยกัน | `type User struct { Name string }` |
| **Pointer** | ตัวแปรที่เก็บ address ของตัวแปรอื่น | `var p *int` |
| **Zero Value** | ค่าเริ่มต้นของตัวแปรที่ยังไม่กำหนดค่า | `var i int` → `0` |

---

## 1.8 ออกแบบ Workflow และ Dataflow

### 1.8.1 การติดตั้งและเริ่มต้นโปรเจค Workflow

```
รูปที่ 1.1: Go Development Workflow

┌─────────────────────────────────────────────────────────────────────┐
│                    GO DEVELOPMENT WORKFLOW                          │
└─────────────────────────────────────────────────────────────────────┘

    [Start]
       │
       ▼
┌─────────────────┐
│ 1. Install Go   │  →  https://go.dev/dl/
│   ติดตั้ง Go      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ 2. Set GOPATH   │  →  ตั้งค่า GOPATH (Windows: %USERPROFILE%\go)
│   ตั้งค่า PATH    │      (Linux/macOS: ~/go)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ 3. Create       │  →  mkdir myproject && cd myproject
│   โปรเจคใหม่     │      go mod init github.com/username/myproject
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ 4. Write Code   │  →  สร้างไฟล์ main.go
│   เขียนโค้ด       │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ 5. Build/Run    │  →  go build    (สร้าง binary)
│   รันโปรแกรม      │      go run     (รันโดยไม่สร้าง binary)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ 6. Test         │  →  go test ./...
│   ทดสอบ          │      go test -v -cover
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ 7. Deploy       │  →  ./myproject (บน Linux/macOS)
│   ติดตั้งใช้งาน    │      myproject.exe (บน Windows)
└─────────────────┘
       │
       ▼
    [End]
```

### 1.8.2 Go Compilation Process (Dataflow)

```
รูปที่ 1.2: Go Compilation Dataflow

┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Source     │     │   Parser     │     │   AST        │
│   Code       │────▶│   วิเคราะห์   │────▶│   Abstract   │
│   main.go    │     │   syntax     │     │   Syntax Tree│
└──────────────┘     └──────────────┘     └──────┬───────┘
                                                   │
                                                   ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Machine    │     │   SSA        │     │   Type       │
│   Code       │◀────│   Static     │◀────│   Checking   │
│   .exe/.out  │     │   Single     │     │   ตรวจสอบชนิด │
│              │     │   Assignment │     │              │
└──────────────┘     └──────────────┘     └──────────────┘
        │
        ▼
┌──────────────┐
│   Runtime    │
│   Execution  │
│   ทำงานจริง   │
└──────────────┘
```

---

## 1.9 การติดตั้งและตั้งค่า (Installation & Setup)

### 1.9.1 Windows Installation

```powershell
# ขั้นตอนที่ 1: ดาวน์โหลด installer จาก https://go.dev/dl/
# เลือก: go1.22.x.windows-amd64.msi

# ขั้นตอนที่ 2: รัน installer และทำตามขั้นตอน
# ติดตั้งที่ C:\Go\

# ขั้นตอนที่ 3: ตรวจสอบการติดตั้ง
go version
# Expected output: go version go1.22.0 windows/amd64

# ขั้นตอนที่ 4: ตั้งค่า GOPATH (Environment Variable)
# เปิด System Properties → Environment Variables
# เพิ่ม: GOPATH = %USERPROFILE%\go

# ขั้นตอนที่ 5: ตรวจสอบ Environment
go env GOPATH
# Expected output: C:\Users\YourName\go
```

### 1.9.2 Linux/macOS Installation

```bash
# macOS (using Homebrew)
brew install go
go version

# Linux (Ubuntu/Debian)
sudo apt update
sudo apt install golang-go
go version

# หรือดาวน์โหลด manual
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### 1.9.3 VS Code Extensions

```json
{
  "recommendations": [
    "golang.go",           // Official Go extension
    "premparihar.gotestexplorer",  // Test Explorer
    "zxh404.vscode-proto3"  // Protocol Buffer support
  ]
}
```

---

## 1.10 โค้ดแรก: Hello World (Runnable Code Example)

```go
// File: main.go
// โปรแกรมแรกในภาษา Go - แสดงข้อความ "Hello, World!"
// First program in Go - Display "Hello, World!"

package main
// package main คือ package หลักที่ใช้สำหรับสร้าง executable
// package main is the main package for creating executable

import "fmt"
// import "fmt" - นำเข้า package สำหรับรับ/ส่งข้อมูล (format)
// import "fmt" - import package for formatted I/O

// ฟังก์ชัน main เป็นจุดเริ่มต้นของโปรแกรม (entry point)
// main function is the entry point of the program
func main() {
    // Println พิมพ์ข้อความและขึ้นบรรทัดใหม่
    // Println prints text with new line at the end
    fmt.Println("Hello, World!")
    fmt.Println("สวัสดีชาวโลก!")
    
    // ตัวอย่างการพิมพ์ตัวแปร
    // Example of printing variables
    name := "Gopher"
    age := 12
    
    // Printf พิมพ์แบบมี format (%s = string, %d = integer)
    // Printf prints with formatting (%s = string, %d = integer)
    fmt.Printf("ชื่อ: %s, อายุ: %d ปี\n", name, age)
    // Output: ชื่อ: Gopher, อายุ: 12 ปี
}
```

### วิธีการรัน (How to Run)

```bash
# วิธีที่ 1: รันโดยตรง (ไม่สร้างไฟล์)
# Method 1: Run directly (no binary file created)
go run main.go

# วิธีที่ 2: สร้าง binary แล้วรัน
# Method 2: Build binary then run
go build main.go
./main     # บน Linux/macOS
main.exe   # บน Windows

# วิธีที่ 3: ติดตั้งเป็นคำสั่งใน $GOPATH/bin
# Method 3: Install as command in $GOPATH/bin
go install
```

### Expected Output

```
Hello, World!
สวัสดีชาวโลก!
ชื่อ: Gopher, อายุ: 12 ปี
```

---

## 1.11 กรณีศึกษาและแนวทางแก้ไขปัญหา (Case Studies)

### กรณีศึกษาที่ 1: Uber ย้ายจาก Node.js มา Go

**ปัญหา (Problem):**
- ระบบ Geolocation ต้องจัดการ Request 1 ล้านครั้ง/วินาที
- Node.js มีปัญหาด้วย Memory และ CPU สูง
- Response Time ช้า (เฉลี่ย 200ms)

**แนวทางแก้ไข (Solution):**
```go
// โครงสร้างที่ Uber ใช้
// Structure that Uber uses

package geolocation

import (
    "sync"
    "time"
)

type LocationService struct {
    mu        sync.RWMutex
    locations map[string]*Location
}

func (s *LocationService) GetLocation(userID string) *Location {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    return s.locations[userID]
}

// ใช้ Goroutine จัดการ Request พร้อมกัน
// Use Goroutine to handle concurrent requests
func (s *LocationService) ProcessBatch(users []string) {
    var wg sync.WaitGroup
    // สร้าง worker pool
    // Create worker pool
    workers := 100
    ch := make(chan string, len(users))
    
    // สร้าง worker จำนวน 100 ตัว
    // Create 100 workers
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for userID := range ch {
                s.GetLocation(userID)
            }
        }()
    }
    
    // ส่งงานเข้า channel
    // Send jobs to channel
    for _, userID := range users {
        ch <- userID
    }
    close(ch)
    
    wg.Wait()
}
```

**ผลลัพธ์ (Result):**
- Response Time ลดลงเหลือ 10ms (95% faster)
- Memory Usage ลดลง 80%
- CPU Usage ลดลง 70%
- รองรับ Request ได้ 5 ล้านครั้ง/วินาที

---

### กรณีศึกษาที่ 2: Docker พัฒนาด้วย Go

**ปัญหา (Problem):**
- ต้องการเครื่องมือที่ทำงานได้ทุก Platform (Linux, Windows, macOS)
- ต้องการ Performance สูง ใกล้เคียง Native
- ต้องการ Static Binary ที่ไม่มี Dependency

**แนวทางแก้ไข (Solution):**
```go
// Docker CLI โครงสร้าง
// Docker CLI structure

package main

import (
    "flag"
    "fmt"
    "os"
    "os/exec"
)

func main() {
    // กำหนด command line flags
    // Define command line flags
    var (
        image = flag.String("image", "", "Docker image name")
        port  = flag.Int("port", 8080, "Container port")
    )
    flag.Parse()
    
    if *image == "" {
        fmt.Println("Error: --image is required")
        os.Exit(1)
    }
    
    // รัน container (ตัวอย่างจำลอง)
    // Run container (example simulation)
    cmd := exec.Command("docker", "run", "-d", "-p", 
        fmt.Sprintf("%d:80", *port), *image)
    
    output, err := cmd.Output()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("Container started: %s\n", string(output))
}
```

**ผลลัพธ์ (Result):**
- 1 Binary File ~20MB (เทียบกับ Docker ที่เป็น Python ก่อนหน้านี้)
- รันได้ทุก Platform โดยไม่ต้องติดตั้ง Runtime
- Startup Time < 1ms

---

## 1.12 ปัญหาที่อาจเกิดขึ้นและแนวทางแก้ไข (Troubleshooting)

### ปัญหาที่ 1: `go: go.mod file not found`

```bash
# สาเหตุ: ยังไม่ได้ initialize module
# Cause: Module not initialized

# แนวทางแก้ไข:
go mod init <project-name>

# ตัวอย่าง:
go mod init hello-world
```

### ปัญหาที่ 2: `package is not in GOROOT`

```bash
# สาเหตุ: import path ไม่ถูกต้อง
# Cause: Incorrect import path

# แนวทางแก้ไข:
# ตรวจสอบ go.mod ว่า module name ถูกต้อง
# Check go.mod for correct module name

# cat go.mod
module github.com/username/project

# แก้ไข import ให้ตรงกับ module name
# Fix import to match module name
import "github.com/username/project/mypackage"
```

### ปัญหาที่ 3: Build ช้ามาก

```bash
# สาเหตุ: Go Modules cache เสีย
# Cause: Corrupted Go Modules cache

# แนวทางแก้ไข:
go clean -modcache
go mod download
```

---

## 1.13 สรุป (Summary)

### ประโยชน์ที่ได้รับ (Benefits)

| ข้อ | ประโยชน์ | คำอธิบาย |
|-----|----------|----------|
| 1 | **Performance สูง** | เทียบเท่า C/C++ ในการประมวลผล |
| 2 | **Concurrency เด่น** | Goroutine จัดการงานพร้อมกันได้ดี |
| 3 | **Deploy ง่าย** | Binary เดียว ไม่ต้องติดตั้ง Runtime |
| 4 | **เรียนรู้ง่าย** | 25 keywords เท่านั้น |
| 5 | **Tooling ดี** | Built-in testing, formatting, profiling |
| 6 | **Cross-platform** | รองรับทุก OS ยอดนิยม |

### ข้อควรระวัง (Cautions)

| ข้อ | ข้อควรระวัง | คำอธิบาย |
|-----|-------------|----------|
| 1 | **ไม่มี Generics (ก่อน 1.18)** | เวอร์ชันเก่าต้องใช้ interface{} |
| 2 | **Error Handling ซ้ำซาก** | ต้องตรวจสอบ error ทุกครั้ง |
| 3 | **ไม่มี Inheritance** | ต้องใช้ Composition แทน |
| 4 | **Garbage Collection** | อาจเกิด STW (Stop The World) ในบางกรณี |
| 5 | **Module Management** | ต้องเข้าใจ Go Modules ให้ดี |

### ข้อดี (Advantages)

```
✅ ไวยากรณ์สั้น กระชับ เรียนรู้ง่าย
✅ รวบรวมเป็น Binary ไฟล์เดียว ไม่มี Dependency
✅ Goroutine จัดการ Concurrency ได้ดีกว่า Thread ทั่วไป
✅ Built-in Testing Framework
✅ Standard Library ครอบคลุม (net/http, crypto, encoding)
✅ Cross-compile รองรับทุก Platform
✅ Garbage Collection อัตโนมัติ
✅ Static Type ช่วยลด Error
✅ Fast Compilation
✅ Strong Community (Google, Uber, Dropbox, Docker)
```

### ข้อเสีย (Disadvantages)

```
❌ ไม่มี Generics (เพิ่งมีใน 1.18)
❌ Error Handling ซ้ำซาก (if err != nil)
❌ ไม่มี Inheritance (แต่ใช้ Composition แทน)
❌ Package Management เคยยุ่งยาก (ปัจจุบัน Go Modules ดีขึ้น)
❌ GUI Library ไม่แข็งแรง
❌ Mobile Development ยังไม่
```

### ข้อห้าม (Prohibitions)

| ข้อห้าม | เหตุผล | วิธีแก้ไข |
|---------|--------|----------|
| **ห้ามใช้ Goroutine โดยไม่ควบคุม** | จะทำให้ Goroutine leak | ใช้ WaitGroup หรือ Context |
| **ห้ามส่ง Channel โดยไม่ปิด** | Deadlock | ใช้ `close(ch)` เมื่อเสร็จ |
| **ห้าม ignore error** | อาจทำให้โปรแกรมพัง | ตรวจสอบ `if err != nil` ทุกครั้ง |
| **ห้ามใช้ Panic แทน Error** | ทำให้โปรแกรม Crash | ใช้ Error Return แทน |
| **ห้ามใช้ Global Variables** | Race condition | ใช้ Dependency Injection |

---

## 1.14 แบบฝึกหัดท้ายบท (Exercises)

### แบบฝึกหัดที่ 1: Hello World แบบมีชื่อ
**โจทย์:** เขียนโปรแกรม Go ที่รับชื่อผู้ใช้จาก命令行 (command line argument) แล้วแสดงข้อความ "Hello, [ชื่อ]!"

<details>
<summary>เฉลย (Solution)</summary>

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // รับ argument จาก command line
    // Get argument from command line
    args := os.Args
    
    if len(args) < 2 {
        fmt.Println("กรุณาใส่ชื่อของคุณ (Please enter your name)")
        fmt.Println("Usage: go run main.go <name>")
        return
    }
    
    name := args[1]
    fmt.Printf("Hello, %s!\n", name)
    fmt.Printf("สวัสดี, %s!\n", name)
}
```

</details>

---

### แบบฝึกหัดที่ 2: คำนวณพื้นที่สี่เหลี่ยม
**โจทย์:** เขียนฟังก์ชัน `calculateArea(width, height float64) float64` ที่คำนวณพื้นที่สี่เหลี่ยมผืนผ้า และแสดงผล

<details>
<summary>เฉลย (Solution)</summary>

```go
package main

import "fmt"

// calculateArea คำนวณพื้นที่สี่เหลี่ยมผืนผ้า
// calculateArea calculates rectangle area
func calculateArea(width, height float64) float64 {
    // สูตร: กว้าง x สูง
    // Formula: width x height
    return width * height
}

func main() {
    // ตัวอย่างการใช้งาน
    // Example usage
    width := 10.5
    height := 5.2
    
    area := calculateArea(width, height)
    
    fmt.Printf("สี่เหลี่ยมผืนผ้า กว้าง %.2f สูง %.2f\n", width, height)
    fmt.Printf("มีพื้นที่ %.2f ตารางหน่วย\n", area)
    
    // ทดสอบกับค่าอื่น
    // Test with other values
    fmt.Println("\nการทดสอบเพิ่มเติม (Additional tests):")
    fmt.Printf("3 x 4 = %.2f\n", calculateArea(3, 4))
    fmt.Printf("7.5 x 2.5 = %.2f\n", calculateArea(7.5, 2.5))
}
```

</details>

---

### แบบฝึกหัดที่ 3: แปลงหน่วยอุณหภูมิ
**โจทย์:** เขียนโปรแกรมแปลงอุณหภูมิจาก Celsius เป็น Fahrenheit (สูตร: °F = (°C × 9/5) + 32)

<details>
<summary>เฉลย (Solution)</summary>

```go
package main

import "fmt"

// celsiusToFahrenheit แปลง Celsius เป็น Fahrenheit
// celsiusToFahrenheit converts Celsius to Fahrenheit
func celsiusToFahrenheit(celsius float64) float64 {
    // สูตรการแปลง (Conversion formula)
    // °F = (°C × 9/5) + 32
    return (celsius * 9 / 5) + 32
}

// fahrenheitToCelsius แปลง Fahrenheit เป็น Celsius
// fahrenheitToCelsius converts Fahrenheit to Celsius
func fahrenheitToCelsius(fahrenheit float64) float64 {
    // สูตรการแปลง (Conversion formula)
    // °C = (°F - 32) × 5/9
    return (fahrenheit - 32) * 5 / 9
}

func main() {
    // ทดสอบการแปลง (Test conversion)
    celsius := 25.0
    fahrenheit := celsiusToFahrenheit(celsius)
    
    fmt.Printf("%.1f°C = %.1f°F\n", celsius, fahrenheit)
    
    // ทดสอบกลับทาง (Reverse test)
    backToCelsius := fahrenheitToCelsius(fahrenheit)
    fmt.Printf("%.1f°F = %.1f°C\n", fahrenheit, backToCelsius)
    
    // แสดงตารางเปรียบเทียบ
    // Show comparison table
    fmt.Println("\nตารางเปรียบเทียบอุณหภูมิ (Temperature Comparison Table):")
    fmt.Println("-------------------------------------")
    fmt.Println("  Celsius  |  Fahrenheit")
    fmt.Println("-------------------------------------")
    
    for c := -20.0; c <= 40; c += 10 {
        f := celsiusToFahrenheit(c)
        fmt.Printf("  %8.1f  |  %8.1f\n", c, f)
    }
}
```

**Expected Output:**
```
25.0°C = 77.0°F
77.0°F = 25.0°C

ตารางเปรียบเทียบอุณหภูมิ (Temperature Comparison Table):
-------------------------------------
  Celsius  |  Fahrenheit
-------------------------------------
    -20.0  |     -4.0
    -10.0  |     14.0
      0.0  |     32.0
     10.0  |     50.0
     20.0  |     68.0
     30.0  |     86.0
     40.0  |    104.0
```

</details>

---

### แบบฝึกหัดที่ 4: สร้างโปรแกรมทายตัวเลข
**โจทย์:** เขียนโปรแกรมสุ่มตัวเลข 1-100 ให้ผู้ใช้ทาย พร้อมบอกว่าสูงไป/ต่ำไป

<details>
<summary>เฉลย (Solution)</summary>

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

func main() {
    // กำหนด seed สำหรับการสุ่ม
    // Set seed for random number generation
    rand.Seed(time.Now().UnixNano())
    
    // สุ่มเลข 1-100
    // Generate random number 1-100
    secretNumber := rand.Intn(100) + 1
    
    var guess int
    attempts := 0
    
    fmt.Println("🎮 เกมทายตัวเลข 1-100 (Number Guessing Game)")
    fmt.Println("===========================================")
    
    for {
        fmt.Print("ป้อนตัวเลขที่คุณทาย (Enter your guess): ")
        fmt.Scan(&guess)
        attempts++
        
        if guess < secretNumber {
            fmt.Println("📈 ต่ำไป! (Too low!)")
        } else if guess > secretNumber {
            fmt.Println("📉 สูงไป! (Too high!)")
        } else {
            fmt.Printf("🎉 ถูกต้อง! ใช้ความพยายาม %d ครั้ง\n", attempts)
            fmt.Printf("🎉 Correct! You took %d attempts\n", attempts)
            break
        }
    }
}
```

</details>

---

### แบบฝึกหัดที่ 5: สร้างโปรแกรมบันทึกค่าใช้จ่าย
**โจทย์:** สร้างโปรแกรมบันทึกค่าใช้จ่ายประจำวัน สามารถเพิ่มรายการและแสดงยอดรวม

<details>
<summary>เฉลย (Solution)</summary>

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
)

// Expense โครงสร้างบันทึกค่าใช้จ่าย
// Expense structure for recording expenses
type Expense struct {
    Description string
    Amount      float64
}

func main() {
    // สร้าง slice สำหรับเก็บค่าใช้จ่าย
    // Create slice to store expenses
    expenses := make([]Expense, 0)
    scanner := bufio.NewScanner(os.Stdin)
    
    fmt.Println("💰 โปรแกรมบันทึกค่าใช้จ่าย (Expense Tracker)")
    fmt.Println("===========================================")
    
    for {
        fmt.Println("\nเลือกเมนู (Select menu):")
        fmt.Println("1. เพิ่มค่าใช้จ่าย (Add expense)")
        fmt.Println("2. แสดงรายการทั้งหมด (Show all expenses)")
        fmt.Println("3. แสดงยอดรวม (Show total)")
        fmt.Println("4. ออกจากโปรแกรม (Exit)")
        fmt.Print("เลือก (Choice): ")
        
        scanner.Scan()
        choice := strings.TrimSpace(scanner.Text())
        
        switch choice {
        case "1":
            // เพิ่มค่าใช้จ่าย
            // Add expense
            fmt.Print("รายการ (Description): ")
            scanner.Scan()
            desc := strings.TrimSpace(scanner.Text())
            
            fmt.Print("จำนวนเงิน (Amount): ")
            scanner.Scan()
            amountStr := strings.TrimSpace(scanner.Text())
            
            amount, err := strconv.ParseFloat(amountStr, 64)
            if err != nil {
                fmt.Println("❌ จำนวนเงินไม่ถูกต้อง (Invalid amount)")
                continue
            }
            
            expenses = append(expenses, Expense{
                Description: desc,
                Amount:      amount,
            })
            fmt.Printf("✅ เพิ่ม '%s' จำนวน %.2f บาท\n", desc, amount)
            
        case "2":
            // แสดงรายการทั้งหมด
            // Show all expenses
            if len(expenses) == 0 {
                fmt.Println("📭 ไม่มีรายการค่าใช้จ่าย (No expenses)")
                continue
            }
            
            fmt.Println("\n📋 รายการค่าใช้จ่ายทั้งหมด (All Expenses):")
            fmt.Println("-------------------------------------")
            for i, exp := range expenses {
                fmt.Printf("%d. %-20s %10.2f บาท\n", 
                    i+1, exp.Description, exp.Amount)
            }
            
        case "3":
            // แสดงยอดรวม
            // Show total
            total := 0.0
            for _, exp := range expenses {
                total += exp.Amount
            }
            fmt.Printf("\n💰 ยอดรวมทั้งหมด: %.2f บาท\n", total)
            fmt.Printf("💰 Total amount: %.2f THB\n", total)
            
        case "4":
            // ออกจากโปรแกรม
            // Exit program
            fmt.Println("👋 ขอบคุณที่ใช้บริการ (Thank you!)")
            return
            
        default:
            fmt.Println("❌ เมนูไม่ถูกต้อง (Invalid choice)")
        }
    }
}
```

</details>

---

## 1.15 แหล่งอ้างอิง (References)

| แหล่งข้อมูล | URL | คำอธิบาย |
|------------|-----|----------|
| **Official Go Website** | https://go.dev/ | เอกสารทางการ ดาวน์โหลด ติดตั้ง |
| **Go Tour** | https://tour.go.dev/ | เรียนรู้ Go แบบ Interacive |
| **Go by Example** | https://gobyexample.com/ | เรียนรู้ผ่านตัวอย่างโค้ด |
| **Effective Go** | https://go.dev/doc/effective_go | แนวทางการเขียน Go ที่ดี |
| **Go Playground** | https://go.dev/play/ | ทดลองรัน Go ออนไลน์ |
| **Go Blog** | https://go.dev/blog/ | บล็อกอัปเดตจากทีมพัฒนา |
| **Awesome Go** | https://awesome-go.com/ | รายการ Library/เครื่องมือ Go |
| **Go Forum** | https://forum.golangbridge.org/ | ชุมชนช่วยเหลือปัญหา Go |

---

**จบบทที่ 1** ✅

ในบทต่อไป เราจะเกี่ยวกับ:
- **บทที่ 2**: ตัวแปรและชนิดข้อมูล (Variables & Data Types)
- **บทที่ 3**: โครงสร้างควบคุม (Control Structures)
- **บทที่ 4**: ฟังก์ชัน (Functions)

**คำถามทบทวนตนเอง (Self-Review Questions):**
1. Go ถูกพัฒนาขึ้นที่บริษัทอะไร และปีอะไร? (Go was developed by which company and in which year?)
2. Goroutine แตกต่างจาก Thread ทั่วไปอย่างไร? (How is Goroutine different from normal Thread?)
3. ข้อดีหลัก 3 ข้อของ Go คืออะไร? (What are the 3 main advantages of Go?)
4. คำสั่ง `go mod init` ใช้ทำอะไร? (What does `go mod init` do?)
5. Go รองรับการทำงานบนระบบปฏิบัติการใดบ้าง? (Which operating systems does Go support?)

---