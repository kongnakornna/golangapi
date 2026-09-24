# 📘 เล่มที่ 3: การทำงานพร้อมกัน — หัวใจของ Go

## เริ่มทำเนื้อหาทั้งเล่ม


# เรียน Go อย่างมืออาชีพ
## เล่มที่ 3: การทำงานพร้อมกัน — หัวใจของ Go

---

# สารบัญ เล่มที่ 3

### บทที่ 1: Goroutines
- 1.1 Goroutine — คืออะไรและทำงานอย่างไร
- 1.2 Goroutine Example — Countdown Server
- 1.3 Initial Countdown Server
- 1.4 Countdown Server Client และ Complete Server
- 1.5 การตรวจสอบ Goroutines
- 1.6 การทำความเข้าใจ Goroutine อย่างลึกซึ้ง

### บทที่ 2: Channels
- 2.1 Why Channel?
- 2.2 คอนเซปของ Channels
- 2.3 Unbuffered Channels — คอนเซป
- 2.4 Unbuffered Channels — ตัวอย่าง
- 2.5 Unidirectional Channels Type
- 2.6 Buffered Channels — คอนเซป
- 2.7 Buffered Channels — ตัวอย่าง
- 2.8 Channel Synchronization
- 2.9 Channel Directions
- 2.10 Closing Channels
- 2.11 Range over Channels

### บทที่ 3: Select
- 3.1 Select — การเลือก Channel Event
- 3.2 Randomly Select Channel Event
- 3.3 Default Select
- 3.4 Timeouts
- 3.5 Non-Blocking Channel Operations

### บทที่ 4: Synchronization Primitives
- 4.1 WaitGroups — การรอให้ Goroutine ทำงานเสร็จ
- 4.2 Mutexes — ป้องกัน Data Race
- 4.3 sync.Mutex — Multiple Readers, Single Writer
- 4.4 Atomic Counters
- 4.5 One Time Initialization with `sync.Once`
- 4.6 Detecting Data Race

### บทที่ 5: Patterns
- 5.1 Worker Pools
- 5.2 Rate Limiting
- 5.3 Timers และ Tickers
- 5.4 Stateful Goroutines
- 5.5 Looping Goroutine — Project Overview
- 5.6 Project Baseline — Parse JSON
- 5.7 Project Baseline — Download Image
- 5.8 Project Baseline — Save Image
- 5.9 Project Baseline — Analyze Image
- 5.10 Project Baseline — Download Images in Parallel
- 5.11 Problem with Invalid URL และการแก้ด้วย WaitGroup
- 5.12 A Light Touch on Unsafe Counter
- 5.13 Project Baseline — Limited Number of Goroutine
- 5.14 Challenge Question

### บทที่ 6: Context และการยกเลิก
- 6.1 Context Overview
- 6.2 Hello Context — `WithTimeout`
- 6.3 Context — `WithCancel`
- 6.4 Context — `WithValue`
- 6.5 การจัดการ System Interrupt (Ctrl+C)
- 6.6 Cancel All Goroutines

### บทที่ 7: Race Condition และแนวทางแก้
- 7.1 Race Condition — คืออะไร
- 7.2 Go Mantra Solution — "Don’t communicate by sharing memory; share memory by communicating"
- 7.3 Go Mantra Solution — แบบฝึกหัดและเฉลย
- 7.4 Data Race Detection — `go run -race`

### บทที่ 8: รู้จัก Go Scheduler
- 8.1 Basic OS — Process, Process State, PCB
- 8.2 Basic OS — Thread
- 8.3 The Big Picture of Go Scheduler
- 8.4 GOMAXPROCS Experiment
- 8.5 Scheduling when Goroutine is Blocked
- 8.6 Stealing Work

---

# บทที่ 1: Goroutines

---

## 1.1 Goroutine — คืออะไรและทำงานอย่างไร

**Goroutine** คือฟังก์ชันที่ทำงานพร้อมกัน (concurrently) ในภาษา Go เป็น lightweight thread ที่ถูกจัดการโดย Go runtime

### จุดเด่นของ Goroutine

| คุณสมบัติ | คำอธิบาย |
|-----------|----------|
| **Lightweight** | ใช้ memory เพียง ~2KB ต่อ goroutine (เทียบกับ thread ~1MB) |
| **Concurrency** | สามารถมี goroutine นับหมื่นตัวพร้อมกัน |
| **Simple** | ใช้แค่คำว่า `go` นำหน้าฟังก์ชัน |
| **Managed by Go** | Go scheduler จัดการเอง ไม่ต้องใช้ OS thread โดยตรง |

### การสร้าง Goroutine

```go
package main

import (
    "fmt"
    "time"
)

func sayHello() {
    fmt.Println("Hello from goroutine!")
}

func main() {
    // สร้าง goroutine ด้วยคำว่า go
    go sayHello()

    // Goroutine แบบ anonymous
    go func() {
        fmt.Println("Anonymous goroutine")
    }()

    // รอให้ goroutine ทำงาน
    time.Sleep(100 * time.Millisecond)

    fmt.Println("Main function")
}
```

### ตัวอย่างการทำงานพร้อมกัน

```go
package main

import (
    "fmt"
    "time"
)

func printNumbers(prefix string) {
    for i := 1; i <= 5; i++ {
        fmt.Printf("%s: %d\n", prefix, i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // เริ่ม goroutine สองตัวทำงานพร้อมกัน
    go printNumbers("Goroutine A")
    go printNumbers("Goroutine B")

    // รอให้ทำงานเสร็จ
    time.Sleep(800 * time.Millisecond)
    fmt.Println("Done!")
}
```

---

## 1.2 Goroutine Example — Countdown Server

มาสร้างโปรเจกต์ที่ใช้ Goroutine เพื่อสร้าง Countdown Server

### โครงสร้างโปรเจกต์

```
countdown/
├── go.mod
├── server/
│   └── server.go
└── client/
    └── client.go
```

### ภาพรวมการทำงาน

1. Server รับ connection จาก client
2. แต่ละ client จะมี goroutine ของตัวเอง
3. Server นับถอยหลังและส่งข้อมูลกลับไป

---

## 1.3 Initial Countdown Server

```go
// countdown/server/server.go
package main

import (
    "fmt"
    "net"
    "strconv"
    "time"
)

func handleConnection(conn net.Conn) {
    defer conn.Close()

    fmt.Printf("Client connected: %s\n", conn.RemoteAddr())

    for i := 10; i >= 0; i-- {
        msg := strconv.Itoa(i) + "\n"
        conn.Write([]byte(msg))
        time.Sleep(1 * time.Second)
    }

    conn.Write([]byte("Blast off!\n"))
    fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr())
}

func main() {
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("Error starting server:", err)
        return
    }
    defer listener.Close()

    fmt.Println("Countdown server started on :8080")

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }

        // แต่ละ client มี goroutine ของตัวเอง
        go handleConnection(conn)
    }
}
```

---

## 1.4 Countdown Server Client และ Complete Server

### Client

```go
// countdown/client/client.go
package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
)

func main() {
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Error connecting:", err)
        return
    }
    defer conn.Close()

    fmt.Println("Connected to countdown server")

    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        line := scanner.Text()
        fmt.Println(line)

        if line == "Blast off!" {
            break
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading:", err)
    }

    fmt.Println("Disconnected")
}
```

### Complete Server (เพิ่มฟีเจอร์)

```go
// countdown/server/server_complete.go
package main

import (
    "bufio"
    "fmt"
    "net"
    "strconv"
    "strings"
    "sync"
    "time"
)

type CountdownServer struct {
    clients    map[net.Conn]bool
    mutex      sync.Mutex
    countdowns map[net.Conn]int
}

func NewCountdownServer() *CountdownServer {
    return &CountdownServer{
        clients:    make(map[net.Conn]bool),
        countdowns: make(map[net.Conn]int),
    }
}

func (s *CountdownServer) addClient(conn net.Conn) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    s.clients[conn] = true
    s.countdowns[conn] = 10
}

func (s *CountdownServer) removeClient(conn net.Conn) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    delete(s.clients, conn)
    delete(s.countdowns, conn)
}

func (s *CountdownServer) handleConnection(conn net.Conn) {
    defer conn.Close()
    defer s.removeClient(conn)

    s.addClient(conn)
    fmt.Printf("Client connected: %s (total: %d)\n",
        conn.RemoteAddr(), len(s.clients))

    writer := bufio.NewWriter(conn)
    reader := bufio.NewReader(conn)

    for {
        // อ่านคำสั่งจาก client
        msg, err := reader.ReadString('\n')
        if err != nil {
            break
        }

        msg = strings.TrimSpace(msg)

        if msg == "start" {
            s.startCountdown(conn, writer)
        } else if strings.HasPrefix(msg, "set ") {
            parts := strings.Split(msg, " ")
            if len(parts) == 2 {
                if n, err := strconv.Atoi(parts[1]); err == nil && n > 0 && n <= 60 {
                    s.mutex.Lock()
                    s.countdowns[conn] = n
                    s.mutex.Unlock()
                    writer.WriteString("OK Countdown set to " + parts[1] + "\n")
                    writer.Flush()
                }
            }
        } else {
            writer.WriteString("Unknown command: " + msg + "\n")
            writer.Flush()
        }
    }

    fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr())
}

func (s *CountdownServer) startCountdown(conn net.Conn, writer *bufio.Writer) {
    s.mutex.Lock()
    count := s.countdowns[conn]
    s.mutex.Unlock()

    for i := count; i >= 0; i-- {
        msg := strconv.Itoa(i) + "\n"
        writer.WriteString(msg)
        writer.Flush()
        time.Sleep(1 * time.Second)
    }

    writer.WriteString("Blast off!\n")
    writer.Flush()
}

func main() {
    server := NewCountdownServer()

    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("Error starting server:", err)
        return
    }
    defer listener.Close()

    fmt.Println("Countdown server started on :8080")
    fmt.Println("Commands: start, set N (1-60)")

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }

        go server.handleConnection(conn)
    }
}
```

---

## 1.5 การตรวจสอบ Goroutines

Go มีเครื่องมือสำหรับตรวจสอบ goroutine ที่กำลังทำงาน

### วิธีที่ 1: `runtime.NumGoroutine()`

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    fmt.Println("Goroutines:", runtime.NumGoroutine())

    go func() {
        time.Sleep(5 * time.Second)
    }()

    fmt.Println("Goroutines after:", runtime.NumGoroutine())

    time.Sleep(1 * time.Second)
    fmt.Println("Goroutines:", runtime.NumGoroutine())
}
```

### วิธีที่ 2: `pprof` — Profiling

```go
package main

import (
    "fmt"
    "net/http"
    _ "net/http/pprof"
    "runtime"
    "time"
)

func leakyFunction() {
    for {
        time.Sleep(100 * time.Millisecond)
        // Simulate work
    }
}

func main() {
    // Start pprof server
    go func() {
        fmt.Println("pprof server on :6060")
        http.ListenAndServe(":6060", nil)
    }()

    // Start 100 goroutines
    for i := 0; i < 100; i++ {
        go leakyFunction()
    }

    fmt.Println("Goroutines:", runtime.NumGoroutine())

    // Keep running
    select {}
}
```

เข้าถึง pprof ที่: http://localhost:6060/debug/pprof/goroutine

### วิธีที่ 3: `go tool trace`

```go
package main

import (
    "fmt"
    "os"
    "runtime/trace"
    "time"
)

func main() {
    // Create trace file
    f, err := os.Create("trace.out")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    // Start tracing
    trace.Start(f)
    defer trace.Stop()

    // Run your program
    go func() {
        for i := 0; i < 5; i++ {
            fmt.Println("Goroutine:", i)
            time.Sleep(100 * time.Millisecond)
        }
    }()

    time.Sleep(1 * time.Second)
}

// ดู trace: go tool trace trace.out
```

---

## 1.6 การทำความเข้าใจ Goroutine อย่างลึกซึ้ง

### Goroutine vs OS Thread

| คุณสมบัติ | Goroutine | OS Thread |
|-----------|-----------|-----------|
| Memory | ~2KB | ~1MB |
| Creation | เร็วมาก | ช้า |
| Context Switch | เร็ว (user space) | ช้า (kernel space) |
| Scheduling | Go Scheduler | OS Scheduler |
| จำนวน | หมื่น/แสนตัว | จำกัด |

### Goroutine Lifecycle

```
     ┌─────────────────────────────────────────────────┐
     │               GOROUTINE LIFECYCLE               │
     └─────────────────────────────────────────────────┘

    [Created] ──→ [Runnable] ──→ [Running] ──→ [Done]
         │              │              │
         │              │              │
         └──────────────┼──────────────┘
                        │
                        ▼
                  [Blocked]
                        │
                        ▼
                  [Waiting]
                        │
                        ▼
                  [Done]
```

### ตัวอย่าง: Goroutine ที่ถูก Block

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func blockingFunction() {
    // Simulate blocking I/O
    time.Sleep(10 * time.Second)
}

func cpuIntensive() {
    sum := 0
    for i := 0; i < 1000000000; i++ {
        sum += i
    }
    fmt.Println("CPU intensive done:", sum)
}

func main() {
    fmt.Println("Initial goroutines:", runtime.NumGoroutine())

    // เริ่ม goroutine ที่ block
    go blockingFunction()

    // เริ่ม goroutine ที่ใช้ CPU
    go cpuIntensive()

    time.Sleep(100 * time.Millisecond)
    fmt.Println("During execution:", runtime.NumGoroutine())

    // สังเกตว่า goroutine ถูก schedule อย่างไร
    time.Sleep(2 * time.Second)
    fmt.Println("Final goroutines:", runtime.NumGoroutine())
}
```

---

## สรุปบทที่ 1

ในบทนี้คุณได้เรียนรู้:

✅ Goroutine คืออะไรและทำงานอย่างไร
✅ การสร้าง Goroutine ด้วย `go`
✅ Countdown Server ที่ใช้ Goroutine
✅ การตรวจสอบ Goroutines ด้วยเครื่องมือต่างๆ
✅ ความแตกต่างระหว่าง Goroutine และ OS Thread
✅ Goroutine Lifecycle

---

# บทที่ 2: Channels

---

## 2.1 Why Channel?

**Channel** คือกลไกในการสื่อสารระหว่าง Goroutines

### ทำไมต้องใช้ Channel?

1. **การสื่อสารที่ปลอดภัย** — ไม่ต้องใช้ shared memory
2. **Synchronization** — Goroutine รอกันและกันผ่าน channel
3. **ง่ายต่อการเข้าใจ** — code อ่านง่ายกว่าใช้ mutex
4. **Go mantra:** *"Don't communicate by sharing memory; share memory by communicating"*

### เปรียบเทียบการสื่อสาร

**การใช้ shared memory (ไม่แนะนำ):**

```go
var counter int
var mutex sync.Mutex

func increment() {
    mutex.Lock()
    counter++
    mutex.Unlock()
}
```

**การใช้ Channel (แนะนำ):**

```go
ch := make(chan int)

func increment() {
    ch <- 1
}

// In main
go increment()
value := <-ch
```

---

## 2.2 คอนเซปของ Channels

### รูปแบบ

```go
// สร้าง channel
ch := make(chan Type)

// ส่งค่าเข้า channel
ch <- value

// รับค่าจาก channel
value := <-ch
```

### ตัวอย่างพื้นฐาน

```go
package main

import "fmt"

func main() {
    // สร้าง channel ที่ส่ง int
    ch := make(chan int)

    // Goroutine ส่งค่าเข้า channel
    go func() {
        ch <- 42
    }()

    // รับค่าจาก channel
    value := <-ch
    fmt.Println(value)  // 42
}
```

### Channel มีทิศทาง

```go
// สองทิศทาง (รับและส่ง)
ch := make(chan int)

// ส่งอย่างเดียว
sendOnly := make(chan<- int)

// รับอย่างเดียว
receiveOnly := make(<-chan int)
```

---

## 2.3 Unbuffered Channels — คอนเซป

**Unbuffered Channel** คือ channel ที่ไม่มี buffer — การส่งจะ block จนกว่ามีการรับ

### การทำงาน

```
Goroutine A           Goroutine B
    │                       │
    │   ch <- 42            │
    │   (block)             │
    │                       │
    │                       │  value := <-ch
    │                       │  (รับค่า)
    │   (unblock)           │
    ▼                       ▼
```

### ตัวอย่าง

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch := make(chan string)

    // Goroutine ส่งค่า
    go func() {
        fmt.Println("Sending message...")
        ch <- "Hello"
        fmt.Println("Message sent!")
    }()

    // รอ 2 วินาที แล้วรับค่า
    time.Sleep(2 * time.Second)
    msg := <-ch
    fmt.Println("Received:", msg)
}
```

---

## 2.4 Unbuffered Channels — ตัวอย่าง

### ตัวอย่าง: Worker-Pool พื้นฐาน

```go
package main

import (
    "fmt"
    "time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job)
        time.Sleep(1 * time.Second)
        results <- job * 2
    }
}

func main() {
    const numJobs = 5
    const numWorkers = 3

    jobs := make(chan int)
    results := make(chan int)

    // เริ่ม workers
    for w := 1; w <= numWorkers; w++ {
        go worker(w, jobs, results)
    }

    // ส่งงาน
    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    // รับผลลัพธ์
    for r := 1; r <= numJobs; r++ {
        result := <-results
        fmt.Println("Result:", result)
    }
}
```

### ตัวอย่าง: การซิงโครไนซ์ด้วย Channel

```go
package main

import "fmt"

func printNumbers(ch chan bool) {
    for i := 1; i <= 5; i++ {
        fmt.Println(i)
    }
    ch <- true  // ส่งสัญญาณว่าเสร็จแล้ว
}

func main() {
    done := make(chan bool)

    go printNumbers(done)

    // รอให้ goroutine เสร็จ
    <-done
    fmt.Println("Done!")
}
```

---

## 2.5 Unidirectional Channels Type

Channel ที่มีทิศทางเดียว

### ตัวอย่าง

```go
package main

import "fmt"

// รับเฉพาะการส่ง
func sendOnly(ch chan<- int, values []int) {
    for _, v := range values {
        ch <- v
    }
    close(ch)
}

// รับเฉพาะการรับ
func receiveOnly(ch <-chan int) {
    for v := range ch {
        fmt.Println("Received:", v)
    }
}

func main() {
    ch := make(chan int)

    go sendOnly(ch, []int{1, 2, 3, 4, 5})

    receiveOnly(ch)
}
```

### ประโยชน์

1. **ความปลอดภัย** — ป้องกันการส่ง/รับผิดทิศทาง
2. **เอกสาร** — บอกเจตนาของฟังก์ชัน
3. **Type safety** — compiler ตรวจสอบให้

---

## 2.6 Buffered Channels — คอนเซป

**Buffered Channel** มี buffer ที่เก็บค่าได้ โดยไม่ต้องรอ receiver

### การทำงาน

```
Goroutine A           Buffer            Goroutine B
    │                    │                    │
    │   ch <- 42         │                    │
    │───────────────────▶│   [42]             │
    │   (ไม่ block)       │                    │
    │                    │                    │
    │   ch <- 43         │                    │
    │───────────────────▶│   [42, 43]         │
    │   (ไม่ block)       │                    │
    │                    │                    │
    │                    │   value := <-ch    │
    │                    │◀───────────────────│
    │                    │   [43]             │
    │                    │                    │
    │                    │   value := <-ch    │
    │                    │◀───────────────────│
    │                    │   []               │
    ▼                    ▼                    ▼
```

### ตัวอย่าง

```go
package main

import "fmt"

func main() {
    // Buffered channel ขนาด 2
    ch := make(chan string, 2)

    // ส่งค่าโดยไม่ block (ยังมีที่ว่าง)
    ch <- "Hello"
    ch <- "World"

    // รับค่า
    fmt.Println(<-ch)  // Hello
    fmt.Println(<-ch)  // World

    // ถ้าส่งเกิน buffer จะ block
    // ch <- "Extra"  // block! (ไม่มี receiver)
}
```

---

## 2.7 Buffered Channels — ตัวอย่าง

### ตัวอย่าง: Producer-Consumer

```go
package main

import (
    "fmt"
    "time"
)

func producer(ch chan<- int) {
    for i := 1; i <= 10; i++ {
        fmt.Printf("Producing: %d\n", i)
        ch <- i
        time.Sleep(100 * time.Millisecond)
    }
    close(ch)
}

func consumer(ch <-chan int, id int) {
    for v := range ch {
        fmt.Printf("Consumer %d: %d\n", id, v)
        time.Sleep(300 * time.Millisecond)
    }
}

func main() {
    // Buffered channel ขนาด 3
    ch := make(chan int, 3)

    go producer(ch)
    go consumer(ch, 1)
    go consumer(ch, 2)

    time.Sleep(5 * time.Second)
}
```

### ตัวอย่าง: การใช้ Buffer เพื่อลด Latency

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // Unbuffered
    unbuffered := make(chan int)
    start := time.Now()

    go func() {
        for i := 0; i < 1000; i++ {
            unbuffered <- i
        }
        close(unbuffered)
    }()

    for range unbuffered {
        // Process
    }
    fmt.Println("Unbuffered time:", time.Since(start))

    // Buffered (ขนาด 100)
    buffered := make(chan int, 100)
    start = time.Now()

    go func() {
        for i := 0; i < 1000; i++ {
            buffered <- i
        }
        close(buffered)
    }()

    for range buffered {
        // Process
    }
    fmt.Println("Buffered time:", time.Since(start))
}
```

---

## 2.8 Channel Synchronization

การใช้ Channel เพื่อ synchronize goroutines

### ตัวอย่าง: ใช้ Channel แทน WaitGroup

```go
package main

import (
    "fmt"
    "time"
)

func worker(id int, done chan bool) {
    fmt.Printf("Worker %d starting\n", id)
    time.Sleep(time.Duration(id) * time.Second)
    fmt.Printf("Worker %d done\n", id)
    done <- true
}

func main() {
    const numWorkers = 5
    done := make(chan bool, numWorkers)

    for i := 1; i <= numWorkers; i++ {
        go worker(i, done)
    }

    // รอให้ทุก worker เสร็จ
    for i := 0; i < numWorkers; i++ {
        <-done
    }

    fmt.Println("All workers done!")
}
```

### ตัวอย่าง: การส่งสัญญาณ Start/Stop

```go
package main

import (
    "fmt"
    "time"
)

func worker(stop <-chan bool) {
    for {
        select {
        case <-stop:
            fmt.Println("Worker stopped")
            return
        default:
            fmt.Println("Working...")
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    stop := make(chan bool)

    go worker(stop)

    time.Sleep(3 * time.Second)
    stop <- true

    time.Sleep(1 * time.Second)
    fmt.Println("Main done")
}
```

---

## 2.9 Channel Directions

Channel สามารถระบุทิศทางในการประกาศ

### ตัวอย่าง

```go
package main

import "fmt"

// send-only channel
func producer(out chan<- int) {
    for i := 0; i < 5; i++ {
        out <- i
    }
    close(out)
}

// receive-only channel
func consumer(in <-chan int) {
    for v := range in {
        fmt.Println(v)
    }
}

// send and receive channel
func process(in <-chan int, out chan<- int) {
    for v := range in {
        out <- v * 2
    }
    close(out)
}

func main() {
    ch1 := make(chan int)
    ch2 := make(chan int)

    go producer(ch1)
    go process(ch1, ch2)
    consumer(ch2)
}
```

---

## 2.10 Closing Channels

### การปิด Channel

```go
ch := make(chan int)
close(ch)
```

### การตรวจสอบว่า Channel ปิดแล้ว

```go
// วิธีที่ 1: ใช้ ok
value, ok := <-ch
if !ok {
    fmt.Println("Channel closed")
}

// วิธีที่ 2: ใช้ range
for value := range ch {
    fmt.Println(value)
}
// loop ออกเมื่อ channel ปิด
```

### ตัวอย่าง

```go
package main

import "fmt"

func main() {
    ch := make(chan int, 3)

    // ส่งค่า
    ch <- 1
    ch <- 2
    ch <- 3

    close(ch)

    // รับค่า
    fmt.Println(<-ch)  // 1
    fmt.Println(<-ch)  // 2
    fmt.Println(<-ch)  // 3

    // รับจาก channel ที่ปิดแล้ว
    v, ok := <-ch
    fmt.Printf("%d, %t\n", v, ok)  // 0, false

    // range จะออกเมื่อ channel ปิด
    ch2 := make(chan int, 3)
    ch2 <- 10
    ch2 <- 20
    ch2 <- 30
    close(ch2)

    for v := range ch2 {
        fmt.Println(v)
    }
}
```

---

## 2.11 Range over Channels

การใช้ `range` กับ channel เพื่อรับค่าจนกว่า channel จะปิด

### ตัวอย่าง

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch := make(chan int)

    go func() {
        for i := 1; i <= 5; i++ {
            ch <- i
            time.Sleep(100 * time.Millisecond)
        }
        close(ch)
    }()

    // รับค่าจนกว่า channel จะปิด
    for v := range ch {
        fmt.Println("Received:", v)
    }

    fmt.Println("Channel closed!")
}
```

### ตัวอย่าง: Multiple Producers

```go
package main

import (
    "fmt"
    "sync"
)

func producer(id int, ch chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for i := 1; i <= 3; i++ {
        ch <- id*10 + i
    }
}

func main() {
    ch := make(chan int, 10)
    var wg sync.WaitGroup

    // เริ่ม 3 producers
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go producer(i, ch, &wg)
    }

    // เมื่อ producers เสร็จ ให้ปิด channel
    go func() {
        wg.Wait()
        close(ch)
    }()

    // รับค่าจากทุก producer
    for v := range ch {
        fmt.Println("Received:", v)
    }

    fmt.Println("All done!")
}
```

---

## สรุปบทที่ 2

ในบทนี้คุณได้เรียนรู้:

✅ Channel คืออะไรและทำไมต้องใช้
✅ Unbuffered Channels
✅ Buffered Channels
✅ Unidirectional Channels
✅ Channel Synchronization
✅ Closing Channels
✅ Range over Channels

---

# บทที่ 3: Select

---

## 3.1 Select — การเลือก Channel Event

`select` ใช้เพื่อรอการทำงานของหลาย channel พร้อมกัน

### รูปแบบ

```go
select {
case <-ch1:
    // เมื่อ ch1 มีค่า
case value := <-ch2:
    // เมื่อ ch2 มีค่า
case ch3 <- value:
    // เมื่อส่งค่าเข้า ch3 ได้
default:
    // ถ้าไม่มี case ไหนทำงาน
}
```

### ตัวอย่างพื้นฐาน

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
        ch1 <- "จาก ch1"
    }()

    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "จาก ch2"
    }()

    // select จะเลือก case ที่พร้อมทำงาน
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println(msg1)
        case msg2 := <-ch2:
            fmt.Println(msg2)
        }
    }
}
```

---

## 3.2 Randomly Select Channel Event

`select` จะสุ่มเลือกเมื่อหลาย case พร้อมทำงาน

### ตัวอย่าง

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
        for {
            ch1 <- "Ping!"
            time.Sleep(100 * time.Millisecond)
        }
    }()

    go func() {
        for {
            ch2 <- "Pong!"
            time.Sleep(100 * time.Millisecond)
        }
    }()

    for i := 0; i < 10; i++ {
        select {
        case msg := <-ch1:
            fmt.Println(msg)
        case msg := <-ch2:
            fmt.Println(msg)
        }
        time.Sleep(50 * time.Millisecond)
    }
}
```

---

## 3.3 Default Select

`default` ใน select จะทำงานเมื่อไม่มี case ไหนพร้อม

### ตัวอย่าง

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch := make(chan int)

    // Non-blocking receive
    select {
    case v := <-ch:
        fmt.Println("Received:", v)
    default:
        fmt.Println("No value available")
    }

    // Non-blocking send
    select {
    case ch <- 42:
        fmt.Println("Sent 42")
    default:
        fmt.Println("Cannot send (channel full or no receiver)")
    }

    // Loop with default
    ticker := time.NewTicker(500 * time.Millisecond)
    done := time.After(3 * time.Second)

    for {
        select {
        case t := <-ticker.C:
            fmt.Println("Tick at", t)
        case <-done:
            fmt.Println("Done!")
            return
        default:
            fmt.Println("Waiting...")
            time.Sleep(100 * time.Millisecond)
        }
    }
}
```

---

## 3.4 Timeouts

การกำหนด timeout ด้วย `time.After`

### ตัวอย่าง

```go
package main

import (
    "fmt"
    "time"
)

func worker(ch chan string) {
    time.Sleep(2 * time.Second)
    ch <- "Done!"
}

func main() {
    ch := make(chan string)

    go worker(ch)

    // รอด้วย timeout
    select {
    case result := <-ch:
        fmt.Println("Result:", result)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout: worker took too long")
    }

    // ตัวอย่างอื่น: ใช้กับ channel ที่ไม่มี sender
    ch2 := make(chan int)

    select {
    case v := <-ch2:
        fmt.Println("Received:", v)
    case <-time.After(3 * time.Second):
        fmt.Println("Timeout: no data received")
    }
}
```

---

## 3.5 Non-Blocking Channel Operations

การใช้ select เพื่อทำ non-blocking operations

### ตัวอย่าง

```go
package main

import "fmt"

func main() {
    ch := make(chan int, 2)

    // Non-blocking send
    select {
    case ch <- 1:
        fmt.Println("Sent 1")
    default:
        fmt.Println("Channel full")
    }

    select {
    case ch <- 2:
        fmt.Println("Sent 2")
    default:
        fmt.Println("Channel full")
    }

    // ลองส่งเพิ่ม (channel full)
    select {
    case ch <- 3:
        fmt.Println("Sent 3")
    default:
        fmt.Println("Cannot send (channel full)")
    }

    // Non-blocking receive
    select {
    case v := <-ch:
        fmt.Println("Received:", v)
    default:
        fmt.Println("No data")
    }

    // รับที่เหลือ
    fmt.Println(<-ch)  // 2
}
```

---

## สรุปบทที่ 3

ในบทนี้คุณได้เรียนรู้:

✅ Select — การเลือก Channel Event
✅ Random Select
✅ Default Select
✅ Timeouts
✅ Non-Blocking Channel Operations

---

# บทที่ 4: Synchronization Primitives

---

## 4.1 WaitGroups — การรอให้ Goroutine ทำงานเสร็จ

`sync.WaitGroup` ใช้รอให้หลาย goroutine ทำงานเสร็จ

### รูปแบบ

```go
var wg sync.WaitGroup

// เพิ่ม counter
wg.Add(1)

// ลด counter เมื่อเสร็จ
defer wg.Done()

// รอจนกว่า counter เป็น 0
wg.Wait()
```

### ตัวอย่าง

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()  // เมื่อ worker เสร็จ

    fmt.Printf("Worker %d starting\n", id)
    time.Sleep(time.Duration(id) * 500 * time.Millisecond)
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup

    // เริ่ม 5 workers
    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }

    // รอทุก worker
    wg.Wait()
    fmt.Println("All workers done!")
}
```

### ตัวอย่าง: การใช้ Add ก่อนสร้าง Goroutine

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup

    // Add ก่อนสร้าง goroutine
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Goroutine %d\n", id)
        }(i)
    }

    wg.Wait()
    fmt.Println("All done!")
}
```

---

## 4.2 Mutexes — ป้องกัน Data Race

`sync.Mutex` ใช้ป้องกันการเข้าถึงข้อมูลพร้อมกัน

### ตัวอย่างที่ไม่มี Mutex (มี Data Race)

```go
package main

import (
    "fmt"
    "sync"
)

var counter int

func increment(wg *sync.WaitGroup) {
    defer wg.Done()
    for i := 0; i < 1000; i++ {
        counter++  // Data race!
    }
}

func main() {
    var wg sync.WaitGroup

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go increment(&wg)
    }

    wg.Wait()
    fmt.Println("Counter:", counter)  // อาจไม่ใช่ 100000
}
```

### ตัวอย่างที่มี Mutex

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

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

func main() {
    counter := &Counter{}
    var wg sync.WaitGroup

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < 1000; j++ {
                counter.Increment()
            }
        }()
    }

    wg.Wait()
    fmt.Println("Counter:", counter.Value())  // 100000
}
```

---

## 4.3 sync.Mutex — Multiple Readers, Single Writer

`sync.RWMutex` รองรับการอ่านพร้อมกันหลายตัว แต่เขียนได้ครั้งละหนึ่งตัว

### ตัวอย่าง

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type SafeMap struct {
    mu   sync.RWMutex
    data map[string]string
}

func NewSafeMap() *SafeMap {
    return &SafeMap{
        data: make(map[string]string),
    }
}

func (m *SafeMap) Get(key string) string {
    m.mu.RLock()         // ล็อกอ่าน
    defer m.mu.RUnlock()
    return m.data[key]
}

func (m *SafeMap) Set(key, value string) {
    m.mu.Lock()          // ล็อกเขียน
    defer m.mu.Unlock()
    m.data[key] = value
}

func (m *SafeMap) Delete(key string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    delete(m.data, key)
}

func main() {
    sm := NewSafeMap()

    // Readers
    for i := 0; i < 5; i++ {
        go func(id int) {
            for {
                value := sm.Get("key")
                if value != "" {
                    fmt.Printf("Reader %d: %s\n", id, value)
                }
                time.Sleep(100 * time.Millisecond)
            }
        }(i)
    }

    // Writer
    go func() {
        for i := 0; i < 10; i++ {
            sm.Set("key", fmt.Sprintf("value-%d", i))
            fmt.Printf("Writer: value-%d\n", i)
            time.Sleep(500 * time.Millisecond)
        }
    }()

    time.Sleep(5 * time.Second)
}
```

---

## 4.4 Atomic Counters

`sync/atomic` ใช้สำหรับตัวนับที่ปลอดภัย (ไม่มี Mutex)

### ตัวอย่าง

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

func main() {
    var counter int64
    var wg sync.WaitGroup

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < 1000; j++ {
                atomic.AddInt64(&counter, 1)
            }
        }()
    }

    wg.Wait()
    fmt.Println("Counter:", atomic.LoadInt64(&counter))  // 100000
}
```

### Atomic Operations

```go
package main

import (
    "fmt"
    "sync/atomic"
)

func main() {
    var value int64 = 0

    // Add
    atomic.AddInt64(&value, 10)
    fmt.Println("After add:", value)

    // Compare and Swap
    swapped := atomic.CompareAndSwapInt64(&value, 10, 20)
    fmt.Printf("CAS result: %t, value: %d\n", swapped, value)

    // Load and Store
    atomic.StoreInt64(&value, 100)
    fmt.Println("After store:", atomic.LoadInt64(&value))

    // Swap
    old := atomic.SwapInt64(&value, 200)
    fmt.Printf("Old: %d, New: %d\n", old, value)
}
```

---

## 4.5 One Time Initialization with sync.Once

`sync.Once` ใช้สำหรับการ initialize เพียงครั้งเดียว

### ตัวอย่าง

```go
package main

import (
    "fmt"
    "sync"
)

var once sync.Once
var config string

func initConfig() {
    fmt.Println("Initializing config...")
    config = "Config loaded!"
}

func getConfig() string {
    once.Do(initConfig)  // เรียกแค่ครั้งเดียว
    return config
}

func main() {
    var wg sync.WaitGroup

    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Goroutine %d: %s\n", id, getConfig())
        }(i)
    }

    wg.Wait()
}
```

### ตัวอย่าง: Singleton Pattern

```go
package main

import (
    "fmt"
    "sync"
)

type Database struct {
    Name string
}

var (
    instance *Database
    once     sync.Once
)

func GetDatabase() *Database {
    once.Do(func() {
        fmt.Println("Creating database connection...")
        instance = &Database{Name: "MyDB"}
    })
    return instance
}

func main() {
    var wg sync.WaitGroup

    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            db := GetDatabase()
            fmt.Printf("Goroutine %d: %s\n", id, db.Name)
        }(i)
    }

    wg.Wait()
}
```

---

## 4.6 Detecting Data Race

### วิธีตรวจสอบ Data Race

1. **`go run -race`** — ตรวจจับตอนรัน
2. **`go test -race`** — ตรวจจับตอน test
3. **`go build -race`** — สร้าง binary ที่มี race detection

### ตัวอย่าง

```go
// race.go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var counter int
    var wg sync.WaitGroup

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter++  // Data race!
        }()
    }

    wg.Wait()
    fmt.Println("Counter:", counter)
}
```

### การรัน

```bash
go run -race race.go
```

ผลลัพธ์:
```
==================
WARNING: DATA RACE
...
Found 1 data race(s)
```

### ตัวอย่างการตรวจจับด้วย test

```go
// counter_test.go
package main

import (
    "sync"
    "testing"
)

func TestCounter(t *testing.T) {
    var counter int
    var wg sync.WaitGroup

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter++
        }()
    }

    wg.Wait()
    t.Log("Counter:", counter)
}
```

```bash
go test -race counter_test.go
```

---

## สรุปบทที่ 4

ในบทนี้คุณได้เรียนรู้:

✅ WaitGroups — รอให้ Goroutine เสร็จ
✅ Mutexes — ป้องกัน Data Race
✅ RWMutex — Multiple Readers, Single Writer
✅ Atomic Counters
✅ sync.Once — One Time Initialization
✅ Data Race Detection

---

# บทที่ 5: Patterns

---

## 5.1 Worker Pools

Worker Pool คือการมี worker จำนวนจำกัดที่ทำงานจากงานใน queue

### ตัวอย่าง

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Job struct {
    ID    int
    Data  string
}

type Result struct {
    JobID int
    Data  string
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
    defer wg.Done()

    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job.ID)
        time.Sleep(1 * time.Second)  // Simulate work

        results <- Result{
            JobID: job.ID,
            Data:  fmt.Sprintf("Processed: %s", job.Data),
        }
    }
}

func main() {
    const numWorkers = 3
    const numJobs = 10

    jobs := make(chan Job, numJobs)
    results := make(chan Result, numJobs)

    var wg sync.WaitGroup

    // เริ่ม workers
    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go worker(i, jobs, results, &wg)
    }

    // ส่งงาน
    for i := 1; i <= numJobs; i++ {
        jobs <- Job{
            ID:   i,
            Data: fmt.Sprintf("Data-%d", i),
        }
    }
    close(jobs)

    // รอ workers เสร็จ
    go func() {
        wg.Wait()
        close(results)
    }()

    // รับผลลัพธ์
    for result := range results {
        fmt.Printf("Result: Job %d -> %s\n", result.JobID, result.Data)
    }
}
```

---

## 5.2 Rate Limiting

Rate Limiting คือการจำกัดอัตราการทำงาน

### ตัวอย่าง: Ticker-based Rate Limiting

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // จำกัดที่ 3 requests ต่อวินาที
    limiter := time.Tick(333 * time.Millisecond)

    requests := make([]int, 10)
    for i := range requests {
        <-limiter  // รอ
        fmt.Printf("Request %d at %s\n", i, time.Now().Format("15:04:05.000"))
    }
}
```

### ตัวอย่าง: Burst Limiter

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // Burst: อนุญาต 3 requests ทันที
    burstyLimiter := make(chan time.Time, 3)

    // เติม burst
    for i := 0; i < 3; i++ {
        burstyLimiter <- time.Now()
    }

    // ทุก 200ms เติม 1
    go func() {
        for t := range time.Tick(200 * time.Millisecond) {
            burstyLimiter <- t
        }
    }()

    // จำลอง 10 requests
    for i := 0; i < 10; i++ {
        <-burstyLimiter
        fmt.Printf("Request %d at %s\n", i, time.Now().Format("15:04:05.000"))
    }
}
```

---

## 5.3 Timers และ Tickers

### Timer — ทำงานครั้งเดียว

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // Timer 1: ใช้ time.After
    select {
    case <-time.After(2 * time.Second):
        fmt.Println("2 seconds passed")
    }

    // Timer 2: ใช้ time.NewTimer
    timer := time.NewTimer(1 * time.Second)
    <-timer.C
    fmt.Println("1 second passed")

    // Stop timer
    timer2 := time.NewTimer(5 * time.Second)
    go func() {
        <-timer2.C
        fmt.Println("Timer fired (should not happen)")
    }()
    if timer2.Stop() {
        fmt.Println("Timer stopped")
    }

    // Reset timer
    timer3 := time.NewTimer(1 * time.Second)
    timer3.Reset(2 * time.Second)
    <-timer3.C
    fmt.Println("Timer reset and fired")
}
```

### Ticker — ทำงานเป็นระยะ

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // Ticker ทุก 500ms
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()

    done := time.After(3 * time.Second)

    for {
        select {
        case t := <-ticker.C:
            fmt.Println("Tick at", t)
        case <-done:
            fmt.Println("Done!")
            return
        }
    }
}
```

---

## 5.4 Stateful Goroutines

การใช้ Goroutine เพื่อจัดการ state แทน Mutex

### ตัวอย่าง

```go
package main

import (
    "fmt"
    "sync/atomic"
    "time"
)

type ReadOp struct {
    key  int
    resp chan int
}

type WriteOp struct {
    key  int
    val  int
    resp chan bool
}

func main() {
    var readOps uint64
    var writeOps uint64

    reads := make(chan ReadOp)
    writes := make(chan WriteOp)

    // Stateful goroutine
    go func() {
        var state = make(map[int]int)

        for {
            select {
            case read := <-reads:
                read.resp <- state[read.key]
            case write := <-writes:
                state[write.key] = write.val
                write.resp <- true
            }
        }
    }()

    // Readers
    for i := 0; i < 100; i++ {
        go func() {
            for {
                read := ReadOp{
                    key:  i % 10,
                    resp: make(chan int),
                }
                reads <- read
                <-read.resp
                atomic.AddUint64(&readOps, 1)
                time.Sleep(time.Millisecond)
            }
        }()
    }

    // Writers
    for i := 0; i < 10; i++ {
        go func() {
            for {
                write := WriteOp{
                    key:  i % 10,
                    val:  i,
                    resp: make(chan bool),
                }
                writes <- write
                <-write.resp
                atomic.AddUint64(&writeOps, 1)
                time.Sleep(time.Millisecond)
            }
        }()
    }

    time.Sleep(1 * time.Second)

    fmt.Printf("ReadOps: %d\n", atomic.LoadUint64(&readOps))
    fmt.Printf("WriteOps: %d\n", atomic.LoadUint64(&writeOps))
}
```

---

## 5.5 Looping Goroutine — Project Overview

โปรเจกต์นี้จะดาวน์โหลดรูปภาพจำนวนมากแบบ parallel

### โครงสร้างโปรเจกต์

```
downloader/
├── go.mod
├── main.go
├── image/
│   └── (downloaded images)
└── data/
    └── images.json
```

---

## 5.6 Project Baseline — Parse JSON

```go
// downloader/main.go (part 1)
package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
)

type ImageItem struct {
    ID      int    `json:"id"`
    URL     string `json:"url"`
    Title   string `json:"title"`
    Width   int    `json:"width"`
    Height  int    `json:"height"`
}

type ImageData struct {
    Images []ImageItem `json:"images"`
}

func loadImages(filename string) ([]ImageItem, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, err
    }

    var imageData ImageData
    err = json.Unmarshal(data, &imageData)
    if err != nil {
        return nil, err
    }

    return imageData.Images, nil
}

func main() {
    images, err := loadImages("data/images.json")
    if err != nil {
        fmt.Println("Error loading images:", err)
        return
    }

    fmt.Printf("Loaded %d images\n", len(images))

    for i, img := range images[:5] {
        fmt.Printf("%d: %s (%dx%d)\n", i+1, img.Title, img.Width, img.Height)
    }
}
```

---

## 5.7 Project Baseline — Download Image

```go
// downloader/main.go (part 2)
package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
)

func downloadImage(url, filename string) error {
    // สร้าง HTTP request
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // ตรวจสอบ status
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("bad status: %s", resp.Status)
    }

    // สร้างไฟล์
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    // Copy ข้อมูล
    _, err = io.Copy(file, resp.Body)
    return err
}

func downloadImageSequential(images []ImageItem) {
    for i, img := range images {
        filename := filepath.Join("images", fmt.Sprintf("%d_%s.jpg", img.ID, img.Title))

        fmt.Printf("Downloading %d/%d: %s\n", i+1, len(images), img.Title)

        err := downloadImage(img.URL, filename)
        if err != nil {
            fmt.Printf("  Error: %v\n", err)
            continue
        }

        fmt.Printf("  Saved to %s\n", filename)
    }
}
```

---

## 5.8 Project Baseline — Save Image

```go
// downloader/main.go (part 3)
package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func ensureDirectory(dir string) error {
    return os.MkdirAll(dir, 0755)
}

func getSafeFilename(title string, id int) string {
    // แทนที่อักขระที่ไม่ปลอดภัย
    safeTitle := title
    // ... (replace invalid characters)

    return filepath.Join("images", fmt.Sprintf("%d_%s.jpg", id, safeTitle))
}

func saveImages(images []ImageItem) error {
    // สร้าง directory
    if err := ensureDirectory("images"); err != nil {
        return err
    }

    for _, img := range images {
        filename := getSafeFilename(img.Title, img.ID)

        err := downloadImage(img.URL, filename)
        if err != nil {
            fmt.Printf("Error downloading %s: %v\n", img.Title, err)
            continue
        }

        fmt.Printf("Downloaded: %s\n", filename)
    }

    return nil
}
```

---

## 5.9 Project Baseline — Analyze Image

```go
// downloader/main.go (part 4)
package main

import (
    "fmt"
    "image"
    _ "image/jpeg"
    _ "image/png"
    "os"
    "path/filepath"
)

func analyzeImage(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    img, format, err := image.Decode(file)
    if err != nil {
        return err
    }

    bounds := img.Bounds()
    width := bounds.Dx()
    height := bounds.Dy()

    fmt.Printf("  Format: %s, Size: %dx%d\n", format, width, height)

    // ตรวจสอบขนาด
    if width > 2000 || height > 2000 {
        fmt.Printf("  Warning: Large image (%dx%d)\n", width, height)
    }

    return nil
}

func analyzeImages() {
    files, err := filepath.Glob("images/*.jpg")
    if err != nil {
        fmt.Println("Error listing files:", err)
        return
    }

    fmt.Printf("Analyzing %d images\n", len(files))

    for _, file := range files {
        fmt.Printf("Analyzing %s\n", filepath.Base(file))
        if err := analyzeImage(file); err != nil {
            fmt.Printf("  Error: %v\n", err)
        }
    }
}
```

---

## 5.10 Project Baseline — Download Images in Parallel

```go
// downloader/main.go (part 5)
package main

import (
    "fmt"
    "sync"
)

func downloadImageParallel(images []ImageItem, numWorkers int) {
    var wg sync.WaitGroup

    jobs := make(chan ImageItem, len(images))
    results := make(chan string, len(images))

    // เริ่ม workers
    for w := 0; w < numWorkers; w++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for img := range jobs {
                filename := getSafeFilename(img.Title, img.ID)
                fmt.Printf("Worker %d: downloading %s\n", workerID, img.Title)

                err := downloadImage(img.URL, filename)
                if err != nil {
                    results <- fmt.Sprintf("Error: %s - %v", img.Title, err)
                    continue
                }
                results <- fmt.Sprintf("Downloaded: %s", filename)
            }
        }(w)
    }

    // ส่งงาน
    for _, img := range images {
        jobs <- img
    }
    close(jobs)

    // รอ workers เสร็จ
    go func() {
        wg.Wait()
        close(results)
    }()

    // รับผลลัพธ์
    for result := range results {
        fmt.Println(result)
    }
}

func main() {
    // ... (load images)

    fmt.Println("=== Parallel Download ===")
    downloadImageParallel(images, 5)
}
```

---

## 5.11 Problem with Invalid URL และการแก้ด้วย WaitGroup

```go
package main

import (
    "fmt"
    "sync"
)

func downloadWithValidation(images []ImageItem, numWorkers int) {
    var wg sync.WaitGroup
    validImages := make(chan ImageItem, len(images))

    // ตรวจสอบ URL
    for _, img := range images {
        if isValidURL(img.URL) {
            validImages <- img
        } else {
            fmt.Printf("Invalid URL: %s\n", img.URL)
        }
    }
    close(validImages)

    // เริ่ม workers
    for w := 0; w < numWorkers; w++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for img := range validImages {
                // Download
                fmt.Printf("Worker %d: downloading %s\n", workerID, img.Title)
                // ...
            }
        }(w)
    }

    wg.Wait()
    fmt.Println("All downloads completed")
}

func isValidURL(url string) bool {
    // ตรวจสอบว่า URL ถูกต้อง
    return len(url) > 0 && (url[:4] == "http")
}
```

---

## 5.12 A Light Touch on Unsafe Counter

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

func main() {
    var counter int64
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()

            // Unsafe: ไม่มีการป้องกัน
            // counter++

            // Safe: ใช้ atomic
            atomic.AddInt64(&counter, 1)
        }()
    }

    wg.Wait()
    fmt.Println("Counter:", atomic.LoadInt64(&counter))
}
```

---

## 5.13 Project Baseline — Limited Number of Goroutine

```go
package main

import (
    "fmt"
    "sync"
)

type Semaphore struct {
    ch chan struct{}
}

func NewSemaphore(max int) *Semaphore {
    return &Semaphore{
        ch: make(chan struct{}, max),
    }
}

func (s *Semaphore) Acquire() {
    s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
    <-s.ch
}

func downloadWithSemaphore(images []ImageItem, maxConcurrent int) {
    var wg sync.WaitGroup
    sem := NewSemaphore(maxConcurrent)

    for _, img := range images {
        wg.Add(1)
        go func(img ImageItem) {
            defer wg.Done()

            sem.Acquire()
            defer sem.Release()

            // Download
            fmt.Printf("Downloading: %s\n", img.Title)
            // ...
        }(img)
    }

    wg.Wait()
    fmt.Println("All downloads completed")
}
```

---

## 5.14 Challenge Question

### โจทย์

จงเขียนโปรแกรมที่:
1. มีตัวเลข 1-100
2. สร้าง goroutine 5 ตัว
3. แต่ละ goroutine นำตัวเลขมาบวกกัน
4. ใช้ channel และ WaitGroup

### เฉลย

```go
package main

import (
    "fmt"
    "sync"
)

func sumWorker(numbers []int, result chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()

    sum := 0
    for _, n := range numbers {
        sum += n
    }
    result <- sum
}

func main() {
    numbers := make([]int, 100)
    for i := range numbers {
        numbers[i] = i + 1
    }

    const numWorkers = 5
    chunkSize := len(numbers) / numWorkers

    var wg sync.WaitGroup
    results := make(chan int, numWorkers)

    for i := 0; i < numWorkers; i++ {
        start := i * chunkSize
        end := start + chunkSize
        if i == numWorkers-1 {
            end = len(numbers)
        }

        wg.Add(1)
        go sumWorker(numbers[start:end], results, &wg)
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    total := 0
    for sum := range results {
        total += sum
    }

    fmt.Println("Total sum:", total)
    fmt.Println("Expected:", 5050)
}
```

---

## สรุปบทที่ 5

ในบทนี้คุณได้เรียนรู้:

✅ Worker Pools
✅ Rate Limiting
✅ Timers และ Tickers
✅ Stateful Goroutines
✅ Image Downloader Project
✅ Parallel Download
✅ Semaphore Pattern
✅ Challenge Solution

---

# บทที่ 6: Context และการยกเลิก

---

## 6.1 Context Overview

**Context** ใช้สำหรับ:
1. การส่งค่า (values) ผ่าน chain
2. การยกเลิก (cancellation)
3. การกำหนด timeout

### รูปแบบ

```go
// สร้าง context
ctx := context.Background()
ctx := context.TODO()

// With cancellation
ctx, cancel := context.WithCancel(parent)

// With timeout
ctx, cancel := context.WithTimeout(parent, duration)

// With deadline
ctx, cancel := context.WithDeadline(parent, time)

// With value
ctx := context.WithValue(parent, key, value)
```

---

## 6.2 Hello Context — WithTimeout

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func longOperation(ctx context.Context) {
    select {
    case <-time.After(5 * time.Second):
        fmt.Println("Operation completed")
    case <-ctx.Done():
        fmt.Println("Operation cancelled:", ctx.Err())
    }
}

func main() {
    // Timeout 2 วินาที
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    go longOperation(ctx)

    time.Sleep(3 * time.Second)
}
```

---

## 6.3 Context — WithCancel

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

    // เริ่ม 3 workers
    for i := 1; i <= 3; i++ {
        go worker(ctx, i)
    }

    time.Sleep(2 * time.Second)

    fmt.Println("Cancelling all workers...")
    cancel()

    time.Sleep(1 * time.Second)
    fmt.Println("Done")
}
```

---

## 6.4 Context — WithValue

```go
package main

import (
    "context"
    "fmt"
)

type key string

const (
    RequestIDKey key = "request_id"
    UserIDKey    key = "user_id"
)

func handler(ctx context.Context) {
    // ดึงค่าจาก context
    reqID := ctx.Value(RequestIDKey)
    userID := ctx.Value(UserIDKey)

    if reqID != nil {
        fmt.Printf("Request ID: %s\n", reqID)
    }
    if userID != nil {
        fmt.Printf("User ID: %v\n", userID)
    }

    // เรียก function ถัดไป
    processRequest(ctx)
}

func processRequest(ctx context.Context) {
    reqID := ctx.Value(RequestIDKey)
    fmt.Printf("Processing request: %s\n", reqID)
}

func main() {
    ctx := context.Background()

    // เพิ่มค่าใน context
    ctx = context.WithValue(ctx, RequestIDKey, "req-12345")
    ctx = context.WithValue(ctx, UserIDKey, 42)

    handler(ctx)
}
```

---

## 6.5 การจัดการ System Interrupt (Ctrl+C)

```go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func worker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d shutting down\n", id)
            return
        default:
            fmt.Printf("Worker %d working...\n", id)
            time.Sleep(1 * time.Second)
        }
    }
}

func main() {
    // สร้าง context ที่สามารถ cancel ได้
    ctx, cancel := context.WithCancel(context.Background())

    // เริ่ม workers
    for i := 1; i <= 3; i++ {
        go worker(ctx, i)
    }

    // รอ signal
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

    fmt.Println("Press Ctrl+C to stop...")
    <-sigCh

    fmt.Println("Shutting down...")
    cancel()

    time.Sleep(2 * time.Second)
    fmt.Println("Done")
}
```

---

## 6.6 Cancel All Goroutines

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

type WorkerPool struct {
    ctx    context.Context
    cancel context.CancelFunc
    wg     sync.WaitGroup
}

func NewWorkerPool() *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    return &WorkerPool{
        ctx:    ctx,
        cancel: cancel,
    }
}

func (p *WorkerPool) Start(id int) {
    p.wg.Add(1)
    go func() {
        defer p.wg.Done()

        for {
            select {
            case <-p.ctx.Done():
                fmt.Printf("Worker %d: stopped\n", id)
                return
            default:
                fmt.Printf("Worker %d: working\n", id)
                time.Sleep(500 * time.Millisecond)
            }
        }
    }()
}

func (p *WorkerPool) Stop() {
    fmt.Println("Stopping all workers...")
    p.cancel()
    p.wg.Wait()
    fmt.Println("All workers stopped")
}

func main() {
    pool := NewWorkerPool()

    // เริ่ม 5 workers
    for i := 1; i <= 5; i++ {
        pool.Start(i)
    }

    time.Sleep(3 * time.Second)
    pool.Stop()
}
```

---

## สรุปบทที่ 6

ในบทนี้คุณได้เรียนรู้:

✅ Context Overview
✅ WithTimeout
✅ WithCancel
✅ WithValue
✅ System Interrupt Handling
✅ Cancel All Goroutines

---

# บทที่ 7: Race Condition และแนวทางแก้

---

## 7.1 Race Condition — คืออะไร

**Race Condition** เกิดขึ้นเมื่อ goroutine หลายตัวเข้าถึงข้อมูลเดียวกันพร้อมกัน

### ตัวอย่าง

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var counter int
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter++  // Data race!
        }()
    }

    wg.Wait()
    fmt.Println("Counter:", counter)  // อาจไม่ใช่ 1000
}
```

---

## 7.2 Go Mantra Solution

*"Don't communicate by sharing memory; share memory by communicating"*

### วิธีแก้ที่ 1: Use Channel

```go
package main

import "fmt"

func main() {
    ch := make(chan int, 1)
    ch <- 0

    var done = make(chan bool)

    for i := 0; i < 1000; i++ {
        go func() {
            val := <-ch
            ch <- val + 1
            if val+1 == 1000 {
                done <- true
            }
        }()
    }

    <-done
    fmt.Println("Counter:", <-ch)
}
```

### วิธีแก้ที่ 2: Use Mutex

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var counter int
    var mu sync.Mutex
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            mu.Lock()
            counter++
            mu.Unlock()
        }()
    }

    wg.Wait()
    fmt.Println("Counter:", counter)
}
```

### วิธีแก้ที่ 3: Use Atomic

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

func main() {
    var counter int64
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            atomic.AddInt64(&counter, 1)
        }()
    }

    wg.Wait()
    fmt.Println("Counter:", atomic.LoadInt64(&counter))
}
```

---

## 7.3 Go Mantra Solution — แบบฝึกหัดและเฉลย

### แบบฝึกหัด

```go
// จงแก้ไข Data Race ในโค้ดนี้
package main

import (
    "fmt"
    "sync"
)

var data = make(map[string]int)

func update(key string, value int) {
    data[key] = value
}

func read(key string) int {
    return data[key]
}

func main() {
    var wg sync.WaitGroup

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            key := fmt.Sprintf("key-%d", i%10)
            update(key, i)
            value := read(key)
            fmt.Printf("%s: %d\n", key, value)
        }(i)
    }

    wg.Wait()
}
```

### เฉลย

```go
package main

import (
    "fmt"
    "sync"
)

type SafeMap struct {
    mu   sync.RWMutex
    data map[string]int
}

func NewSafeMap() *SafeMap {
    return &SafeMap{
        data: make(map[string]int),
    }
}

func (m *SafeMap) Update(key string, value int) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.data[key] = value
}

func (m *SafeMap) Read(key string) int {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.data[key]
}

func main() {
    sm := NewSafeMap()
    var wg sync.WaitGroup

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            key := fmt.Sprintf("key-%d", i%10)
            sm.Update(key, i)
            value := sm.Read(key)
            fmt.Printf("%s: %d\n", key, value)
        }(i)
    }

    wg.Wait()
}
```

---

## 7.4 Data Race Detection

### การตรวจจับ

```bash
# รันด้วย race detector
go run -race main.go

# หรือ build ด้วย race detector
go build -race main.go

# หรือ test ด้วย race detector
go test -race ./...
```

### ตัวอย่างการตรวจจับ

```go
// race_demo.go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var counter int
    var wg sync.WaitGroup

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter++
        }()
    }

    wg.Wait()
    fmt.Println(counter)
}
```

```bash
$ go run -race race_demo.go

==================
WARNING: DATA RACE
Read at 0x00c0000b2008 by goroutine 7:
  main.main.func1()
      /path/race_demo.go:14 +0x3c

Previous write at 0x00c0000b2008 by goroutine 6:
  main.main.func1()
      /path/race_demo.go:14 +0x58

Goroutine 7 (running) created at:
  main.main()
      /path/race_demo.go:12 +0x8c
...
Found 1 data race(s)
```

---

## สรุปบทที่ 7

ในบทนี้คุณได้เรียนรู้:

✅ Race Condition คืออะไร
✅ วิธีแก้ด้วย Channel, Mutex, Atomic
✅ Go Mantra: "Don't communicate by sharing memory"
✅ Data Race Detection

---

# บทที่ 8: รู้จัก Go Scheduler

---

## 8.1 Basic OS — Process, Process State, PCB

### Process

Process คือโปรแกรมที่กำลังทำงาน มีทรัพยากรของตัวเอง

### Process States

```
       ┌──────────────┐
       │   New        │
       └──────┬───────┘
              ▼
       ┌──────────────┐
       │   Ready      │ ←─────┐
       └──────┬───────┘       │
              ▼               │
       ┌──────────────┐       │
       │   Running    │───────┤ (I/O, Time slice)
       └──────┬───────┘       │
              │               │
              ▼ (I/O, event)  │
       ┌──────────────┐       │
       │   Waiting    │───────┘ (I/O complete)
       └──────┬───────┘
              ▼
       ┌──────────────┐
       │   Terminated │
       └──────────────┘
```

### PCB (Process Control Block)

- Process ID
- Process State
- Program Counter
- Registers
- Memory Information

---

## 8.2 Basic OS — Thread

### Thread vs Process

| คุณสมบัติ | Process | Thread |
|-----------|---------|--------|
| Memory | แยกกัน | แชร์กัน |
| Context Switch | ช้า | เร็ว |
| Creation | ช้า | เร็ว |
| Communication | ยาก (IPC) | ง่าย (shared memory) |
| Resource | มาก | น้อย |

### Thread States

```
Running → Ready (time slice)
Running → Waiting (I/O)
Waiting → Ready (I/O complete)
```

---

## 8.3 The Big Picture of Go Scheduler

### M:N Scheduling

```
Goroutines (G) ←→ Logical Processors (P) ←→ OS Threads (M)
```

### Components

| Component | บทบาท |
|-----------|-------|
| **G (Goroutine)** | ฟังก์ชันที่ทำงานพร้อมกัน |
| **P (Processor)** | กำหนดจำนวน Goroutines ที่รันพร้อมกัน (GOMAXPROCS) |
| **M (Machine)** | OS Thread ที่รัน Goroutine |

### Scheduling

```
                    ┌─────────────────────┐
                    │   GOMAXPROCS = N    │
                    └─────────────────────┘
                              │
           ┌──────────────────┼──────────────────┐
           │                  │                  │
           ▼                  ▼                  ▼
    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
    │     P1      │    │     P2      │    │     P3      │
    │ Local Queue │    │ Local Queue │    │ Local Queue │
    │ G1 G2 G3    │    │ G4 G5 G6    │    │ G7 G8 G9    │
    └─────────────┘    └─────────────┘    └─────────────┘
           │                  │                  │
           ▼                  ▼                  ▼
    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
    │     M1      │    │     M2      │    │     M3      │
    │ (OS Thread) │    │ (OS Thread) │    │ (OS Thread) │
    └─────────────┘    └─────────────┘    └─────────────┘
```

---

## 8.4 GOMAXPROCS Experiment

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func cpuIntensive() {
    sum := 0
    for i := 0; i < 100000000; i++ {
        sum += i
    }
}

func main() {
    // ตรวจสอบค่าเริ่มต้น
    fmt.Println("Default GOMAXPROCS:", runtime.GOMAXPROCS(0))

    // ทดสอบกับ GOMAXPROCS ที่ต่างกัน
    for _, procs := range []int{1, 2, 4, 8} {
        runtime.GOMAXPROCS(procs)

        start := time.Now()
        var wg sync.WaitGroup

        for i := 0; i < 10; i++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                cpuIntensive()
            }()
        }

        wg.Wait()
        elapsed := time.Since(start)

        fmt.Printf("GOMAXPROCS=%d: %v\n", procs, elapsed)
    }
}
```

### ตั้งค่า GOMAXPROCS

```go
// ตั้งค่าที่ runtime
runtime.GOMAXPROCS(4)

// หรือใช้ environment variable
// export GOMAXPROCS=4
```

---

## 8.5 Scheduling when Goroutine is Blocked

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func blockingOperation() {
    fmt.Println("Blocking operation starting")
    time.Sleep(2 * time.Second)  // Simulate blocking I/O
    fmt.Println("Blocking operation done")
}

func cpuOperation(id int) {
    sum := 0
    for i := 0; i < 1000000; i++ {
        sum += i
    }
    fmt.Printf("CPU operation %d done\n", id)
}

func main() {
    runtime.GOMAXPROCS(1)  // Force 1 processor

    // Blocking operation
    go blockingOperation()

    // CPU operations
    for i := 0; i < 5; i++ {
        go cpuOperation(i)
    }

    time.Sleep(3 * time.Second)
}
```

---

## 8.6 Stealing Work

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func work(id int) {
    sum := 0
    for i := 0; i < 1000000; i++ {
        sum += i
    }
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    runtime.GOMAXPROCS(2)

    var wg sync.WaitGroup

    // สร้างงาน 10 งาน
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            work(id)
        }(i)
    }

    // ตัวอย่าง: P1 มีงานเยอะ, P2 มีงานน้อย
    // P2 จะขโมยงานจาก P1 (work stealing)

    wg.Wait()
}
```

### Work Stealing Diagram

```
P1 (งานเยอะ)          P2 (งานน้อย)
    │                      │
    │  Local Queue         │  Local Queue
    │  G1 G2 G3 G4 G5     │  G6
    │                      │
    │                      │  steal work!
    │  G3 G4 G5 ──────────│→ G1 G2
    │                      │
    ▼                      ▼
```

---

## สรุปบทที่ 8

ในบทนี้คุณได้เรียนรู้:

✅ Process vs Thread
✅ Process States และ PCB
✅ M:N Scheduling
✅ GOMAXPROCS
✅ Blocking Scheduling
✅ Work Stealing

---

# 🎯 แบบฝึกหัดทบทวน เล่มที่ 3

## แบบฝึกหัดที่ 1: Goroutine
เขียนโปรแกรมที่มี 2 Goroutine พิมพ์ข้อความสลับกัน

## แบบฝึกหัดที่ 2: Channel
ใช้ Channel ส่งข้อมูลระหว่าง Goroutines

## แบบฝึกหัดที่ 3: Select
ใช้ Select รับข้อมูลจาก 3 Channel

## แบบฝึกหัดที่ 4: WaitGroup
สร้าง Worker Pool 5 ตัว ทำงาน 20 งาน

## แบบฝึกหัดที่ 5: Mutex
ใช้ Mutex ป้องกัน Data Race ใน Counter

## แบบฝึกหัดที่ 6: Context
ใช้ Context เพื่อยกเลิกงานหลังจาก 3 วินาที

## แบบฝึกหัดที่ 7: Rate Limiting
จำกัดการทำงานที่ 10 requests ต่อวินาที

 