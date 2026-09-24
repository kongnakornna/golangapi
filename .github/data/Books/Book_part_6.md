## 📚 เล่มที่ 6: สู่การเป็นนักพัฒนา Go มืออาชีพ

เล่มนี้จะเติมเต็มทักษะที่จำเป็นสำหรับนักพัฒนา Go ในระดับ Production ไม่ว่าจะเป็นการเขียน Benchmark เพื่อวัดประสิทธิภาพ, การสร้าง HTTP Client ที่ยืดหยุ่น, การใช้ Context เพื่อควบคุมการทำงาน, การเขียนโค้ดแบบ Generic ที่นำกลับมาใช้ใหม่ได้, และการทำความเข้าใจว่า Go จัดการกับแนวคิด OOP อย่างไร นอกจากนี้ยังมีคำแนะนำการออกแบบโค้ดที่ดีและชีทสรุปสำหรับการอ้างอิงด่วน

---

### บทที่ 34: การวัดประสิทธิภาพ (Benchmarks)

การ Benchmark คือการทดสอบประสิทธิภาพของโค้ดเพื่อวัดว่าใช้เวลาทำงานนานแค่ไหน, ใช้หน่วยความจำเท่าไหร่, หรือทำงานได้กี่รอบต่อวินาที Go มีเครื่องมือในตัวที่ทำให้การ Benchmark ง่ายมาก

**การเขียน Benchmark:**
กฎคล้ายกับการเขียน Unit Test แต่ชื่อฟังก์ชันต้องขึ้นต้นด้วย `Benchmark` และรับพารามิเตอร์ `*testing.B`

```go
// math.go
package math

func Sum(numbers []int) int {
    total := 0
    for _, n := range numbers {
        total += n
    }
    return total
}
```

```go
// math_test.go
package math

import "testing"

// Benchmark ทั่วไป
func BenchmarkSum(b *testing.B) {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // b.N คือจำนวนรอบที่ Go จะเรียกให้เหมาะสมอัตโนมัติ
    for i := 0; i < b.N; i++ {
        Sum(numbers)
    }
}

// Benchmark ที่ต้องเตรียมข้อมูลก่อน (ใช้ ResetTimer)
func BenchmarkSumWithPrep(b *testing.B) {
    // เตรียมข้อมูล (ไม่นับเวลา)
    numbers := make([]int, 1000)
    for i := range numbers {
        numbers[i] = i
    }
    
    b.ResetTimer() // รีเซ็ตเวลาเริ่มนับใหม่
    
    for i := 0; i < b.N; i++ {
        Sum(numbers)
    }
}

// Benchmark ที่วัดหน่วยความจำ (ใช้ b.ReportAllocs())
func BenchmarkSumMemory(b *testing.B) {
    b.ReportAllocs() // รายงานการจัดสรรหน่วยความจำ
    numbers := []int{1, 2, 3, 4, 5}
    for i := 0; i < b.N; i++ {
        Sum(numbers)
    }
}
```

**การรัน Benchmark:**
```bash
# รัน Benchmark ทั้งหมด
go test -bench=.

# รัน Benchmark เฉพาะฟังก์ชันที่ขึ้นต้นด้วย BenchmarkSum
go test -bench=BenchmarkSum

# รันหลายรอบและแสดงรายละเอียดหน่วยความจำ
go test -bench=. -benchmem -count=5

# ผลลัพธ์ตัวอย่าง:
# BenchmarkSum-8        1000000000         0.235 ns/op        0 B/op        0 allocs/op
# - 8 = จำนวน CPU
# - 1000000000 = จำนวนรอบที่รัน
# - 0.235 ns/op = เวลาเฉลี่ยต่อรอบ
# - 0 B/op = หน่วยความจำที่ใช้ต่อรอบ
# - 0 allocs/op = จำนวนการจัดสรรหน่วยความจำต่อรอบ
```

**การเปรียบเทียบประสิทธิภาพ (ใช้ benchstat):**
```bash
# 安裝 benchstat
go get golang.org/x/perf/cmd/benchstat

# เก็บผลลัพธ์
go test -bench=. -count=10 > old.txt
# แก้ไขโค้ดแล้วรันอีกครั้ง
go test -bench=. -count=10 > new.txt

# เปรียบเทียบ
benchstat old.txt new.txt
```

---

### บทที่ 35: สร้าง HTTP Client

การเรียก API จากภายนอกเป็นสิ่งที่ทำบ่อยมาก Go มี `net/http` สำหรับการสร้าง HTTP Client ที่ทรงพลังและยืดหยุ่น

**HTTP Client พื้นฐาน:**
```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
)

func main() {
    // GET Request
    resp, err := http.Get("https://api.github.com/users/golang")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("Status:", resp.StatusCode)
    fmt.Println("Body:", string(body))

    // POST Request (ส่ง JSON)
    user := map[string]string{"name": "สมชาย", "email": "somchai@mail.com"}
    jsonData, _ := json.Marshal(user)

    resp, err = http.Post(
        "https://httpbin.org/post",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    // อ่าน Response
    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    fmt.Println("Response:", result)
}
```

**การสร้าง HTTP Client แบบกำหนดเอง (Custom Client):**
```go
// สร้าง Client ที่มี Timeout และ Proxy
client := &http.Client{
    Timeout: 10 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:    100,
        IdleConnTimeout: 90 * time.Second,
        DisableCompression: false,
    },
}

// ใช้ Client ในการ Request
req, _ := http.NewRequest("GET", "https://api.example.com/data", nil)
req.Header.Set("Authorization", "Bearer token123")
req.Header.Set("User-Agent", "MyApp/1.0")

resp, err := client.Do(req)
if err != nil {
    panic(err)
}
defer resp.Body.Close()
```

**การจัดการ Response และ Error อย่างมืออาชีพ:**
```go
type APIResponse struct {
    Status  string          `json:"status"`
    Data    json.RawMessage `json:"data"` // เก็บไว้ก่อนเพื่อ Decode ทีหลัง
}

func CallAPI(url string, target interface{}) error {
    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Get(url)
    if err != nil {
        return fmt.Errorf("请求失败: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("สถานะผิดปกติ: %d", resp.StatusCode)
    }

    // Decode JSON
    if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
        return fmt.Errorf(" decode JSON ผิดพลาด: %w", err)
    }
    return nil
}

func main() {
    var data map[string]interface{}
    if err := CallAPI("https://api.github.com/users/golang", &data); err != nil {
        fmt.Println("เกิดข้อผิดพลาด:", err)
        return
    }
    fmt.Println("ชื่อ:", data["login"])
}
```

---

### บทที่ 36: การวิเคราะห์โปรไฟล์ (Program Profiling)

Profiling คือการวิเคราะห์ว่าโปรแกรมใช้ทรัพยากร (CPU, Memory) อย่างไร เพื่อหาจุดที่ทำให้ช้าหรือกินหน่วยความจำมากเกินไป Go มีเครื่องมือ Profiling ในตัวที่ทรงพลัง

**การสร้าง CPU Profile:**
```go
package main

import (
    "os"
    "runtime/pprof"
)

func heavyFunction() {
    // ฟังก์ชันที่ใช้ CPU มาก
    sum := 0
    for i := 0; i < 100000000; i++ {
        sum += i
    }
}

func main() {
    // สร้างไฟล์สำหรับเก็บ CPU Profile
    f, _ := os.Create("cpu.prof")
    defer f.Close()

    // เริ่ม CPU Profiling
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    // เรียกฟังก์ชันที่ต้องการวัด
    heavyFunction()
}
```

**การสร้าง Memory Profile:**
```go
func main() {
    // สร้าง Memory Profile
    f, _ := os.Create("mem.prof")
    defer f.Close()

    // เรียกฟังก์ชันที่ใช้หน่วยความจำมาก
    allocateMemory()

    // เขียน Memory Profile
    pprof.WriteHeapProfile(f)
}

func allocateMemory() {
    data := make([]byte, 100*1024*1024) // 100 MB
    _ = data
}
```

**การวิเคราะห์ Profile ด้วย `go tool pprof`:**
```bash
# วิเคราะห์ CPU Profile
go tool pprof cpu.prof

# วิเคราะห์ Memory Profile
go tool pprof mem.prof

# เปิด Web UI (ต้องมี Graphviz)
go tool pprof -http=:8080 cpu.prof

# คำสั่งใน pprof interactive:
# (pprof) top10     # แสดงฟังก์ชันที่ใช้ CPU มากที่สุด 10 อันดับ
# (pprof) list heavyFunction  # แสดงโค้ดพร้อมเวลา
# (pprof) web       # สร้างกราฟเป็น SVG
```

**การ Profile แบบ HTTP (ใช้ net/http/pprof):**
```go
import _ "net/http/pprof"

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // โค้ดแอปพลิเคชันของคุณ...
}
```
จากนั้นเข้าถึง `http://localhost:6060/debug/pprof/` เพื่อดูข้อมูลแบบ Real-time

---

### บทที่ 37: การจัดการ Context

Context เป็นแพคเกจที่ใช้ในการส่งค่า, การยกเลิกการทำงาน (Cancellation), และการกำหนด Deadline ระหว่าง Goroutine ต่างๆ ในโปรแกรม

**การสร้าง Context พื้นฐาน:**
```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    // 1. Background Context (Context เริ่มต้น)
    ctx := context.Background()

    // 2. TODO Context (ใช้เมื่อยังไม่แน่ใจว่าจะใช้ Context แบบไหน)
    ctx2 := context.TODO()

    // 3. Context ที่มีค่า
    ctxWithValue := context.WithValue(ctx, "userID", 123)

    // 4. Context ที่สามารถยกเลิกได้ (Cancel)
    ctxWithCancel, cancel := context.WithCancel(ctx)
    defer cancel() // อย่าลืมเรียก cancel เพื่อปล่อยทรัพยากร

    // 5. Context ที่มี Deadline
    ctxWithDeadline, cancel2 := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
    defer cancel2()

    // 6. Context ที่มี Timeout
    ctxWithTimeout, cancel3 := context.WithTimeout(ctx, 3*time.Second)
    defer cancel3()
}
```

**การใช้ Context เพื่อยกเลิกการทำงาน:**
```go
func doWork(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("หยุดทำงาน: ", ctx.Err())
            return
        default:
            fmt.Println("กำลังทำงาน...")
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    go doWork(ctx)

    time.Sleep(3 * time.Second)
    fmt.Println("จบโปรแกรม")
}
```

**การใช้ Context ส่งค่า (Value):**
```go
type contextKey string

func processRequest(ctx context.Context) {
    if userID, ok := ctx.Value("userID").(int); ok {
        fmt.Println("User ID:", userID)
    }
    if requestID, ok := ctx.Value("requestID").(string); ok {
        fmt.Println("Request ID:", requestID)
    }
}

func main() {
    ctx := context.Background()
    ctx = context.WithValue(ctx, "userID", 1001)
    ctx = context.WithValue(ctx, "requestID", "req-abc123")
    processRequest(ctx)
}
```

**การประยุกต์ใช้ Context ใน HTTP Server:**
```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    // ใช้ ctx ในการเรียกฟังก์ชันอื่นๆ เพื่อให้สามารถยกเลิกได้
    result, err := longRunningTask(ctx)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, result)
}

func longRunningTask(ctx context.Context) (string, error) {
    select {
    case <-ctx.Done():
        return "", ctx.Err()
    case <-time.After(5 * time.Second):
        return "ทำงานเสร็จ", nil
    }
}
```

---

### บทที่ 38: Generics - การเขียนโค้ดแบบยืดหยุ่น

Generics ถูกเพิ่มใน Go 1.18 ทำให้เราสามารถเขียนฟังก์ชันและโครงสร้างข้อมูลที่ทำงานกับหลายชนิดข้อมูลได้ โดยไม่ต้องเขียนซ้ำ

**ฟังก์ชัน Generic พื้นฐาน:**
```go
package main

import "fmt"

// ฟังก์ชัน Generic ที่รับชนิด T ใดก็ได้
func PrintSlice[T any](s []T) {
    for _, v := range s {
        fmt.Print(v, " ")
    }
    fmt.Println()
}

// ฟังก์ชัน Generic ที่มี Constraint (ต้องเป็นชนิดที่เปรียบเทียบได้)
func FindMax[T comparable](a, b T) T {
    // comparable = ชนิดที่ใช้ == และ != ได้
    if a == b {
        return a
    }
    // แต่ถ้าต้องการเปรียบเทียบมากกว่า/น้อยกว่า ต้องใช้ ordered
    return b
}

// ใช้ Constraint แบบ Custom
type Number interface {
    int | int64 | float64
}

func Sum[T Number](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}

func main() {
    // เรียกใช้ฟังก์ชัน Generic
    PrintSlice([]int{1, 2, 3, 4})
    PrintSlice([]string{"a", "b", "c"})
    PrintSlice([]float64{1.1, 2.2, 3.3})

    fmt.Println(Sum([]int{1, 2, 3}))        // 6
    fmt.Println(Sum([]float64{1.1, 2.2}))   // 3.3
}
```

**Struct Generic:**
```go
// Generic Struct
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
    last := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return last, true
}

func main() {
    intStack := Stack[int]{}
    intStack.Push(10)
    intStack.Push(20)
    val, _ := intStack.Pop()
    fmt.Println(val) // 20

    stringStack := Stack[string]{}
    stringStack.Push("hello")
    stringStack.Push("world")
    val2, _ := stringStack.Pop()
    fmt.Println(val2) // world
}
```

**การใช้ Constraint แบบซับซ้อน (Type Sets):**
```go
// Constraint ที่รวมหลายชนิด
type Numeric interface {
    ~int | ~int64 | ~float64 | ~float32
}

// ฟังก์ชันที่ทำงานกับ Numeric
func Multiply[T Numeric](a, b T) T {
    return a * b
}

// Constraint ที่ต้องการเมธอด
type Stringer interface {
    String() string
}

func PrintString[T Stringer](v T) {
    fmt.Println(v.String())
}
```

---

### บทที่ 39: Go กับกระบวนทัศน์ OOP?

Go ไม่ใช่ภาษา OOP แบบคลาสสิก (ไม่มี Class, Inheritance) แต่มีแนวคิดที่คล้ายคลึงผ่าน Struct และ Interface ช่วยให้เขียนโค้ดที่เป็นระเบียบและยืดหยุ่นได้

**1. Encapsulation (การห่อหุ้ม):**
ใน Go ใช้ **ตัวพิมพ์ใหญ่** สำหรับ Public, **ตัวพิมพ์เล็ก** สำหรับ Private
```go
package user

type User struct {
    Name  string // Public (มองเห็นจากภายนอก)
    email string // Private (มองเห็นเฉพาะใน package)
}

func NewUser(name, email string) *User {
    return &User{Name: name, email: email}
}

func (u *User) GetEmail() string { // Public method
    return u.email
}
```

**2. Composition (การประกอบ) แทน Inheritance:**
Go ใช้การฝัง (Embedding) Struct เพื่อ reuse โค้ด แทนการสืบทอด
```go
type Person struct {
    Name string
    Age  int
}

func (p Person) Greet() string {
    return "สวัสดี " + p.Name
}

// Employee ฝัง Person (Composition)
type Employee struct {
    Person      // ฝัง (Embedding)
    Salary float64
    Position string
}

func main() {
    emp := Employee{
        Person:  Person{Name: "สมชาย", Age: 30},
        Salary:  50000,
        Position: "Developer",
    }
    
    // สามารถเรียกใช้เมธอดของ Person ได้โดยตรง
    fmt.Println(emp.Greet()) // สวัสดี สมชาย
    fmt.Println(emp.Name)   // สมชาย
}
```

**3. Polymorphism (พหุสัณฐาน) ผ่าน Interface:**
```go
type Shape interface {
    Area() float64
}

type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return 3.14159 * c.Radius * c.Radius
}

func PrintArea(s Shape) {
    fmt.Printf("พื้นที่: %.2f\n", s.Area())
}

func main() {
    rect := Rectangle{10, 5}
    circle := Circle{7}
    
    PrintArea(rect)   // พื้นที่: 50.00
    PrintArea(circle) // พื้นที่: 153.94
}
```

**4. Dependency Injection (DI) ใน Go:**
```go
type Database interface {
    Query(sql string) ([]string, error)
}

type MySQLDB struct{}

func (m MySQLDB) Query(sql string) ([]string, error) {
    return []string{"data1", "data2"}, nil
}

type Service struct {
    db Database
}

func NewService(db Database) *Service {
    return &Service{db: db}
}

func (s *Service) GetData() []string {
    result, _ := s.db.Query("SELECT * FROM users")
    return result
}
```

---

### บทที่ 40: การอัปเกรดหรือดาวน์เกรดเวอร์ชัน Go

การจัดการเวอร์ชัน Go เป็นสิ่งสำคัญในการพัฒนา โดยเฉพาะเมื่อต้องทำงานกับโปรเจกต์หลายตัวที่ใช้ Go คนละเวอร์ชัน

**ตรวจสอบเวอร์ชันปัจจุบัน:**
```bash
go version
# go version go1.21.0 linux/amd64
```

**การอัปเกรด Go (MacOS ด้วย Homebrew):**
```bash
brew update
brew upgrade go
```

**การอัปเกรด Go (Windows/Linux ด้วยตัวติดตั้ง):**
1. ดาวน์โหลดเวอร์ชันใหม่จาก https://go.dev/dl/
2. ติดตั้งทับของเดิม (Go จะอัปเดต PATH ให้อัตโนมัติ)

**การใช้ Go Version Manager (gvm) - Linux/Mac:**
```bash
# ติดตั้ง gvm
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)

# ติดตั้ง Go เวอร์ชันต่างๆ
gvm install go1.20
gvm install go1.21

# เปลี่ยนเวอร์ชันที่ใช้
gvm use go1.21 --default

# ดูรายการเวอร์ชันที่ติดตั้ง
gvm list
```

**การอัปเกรดเวอร์ชันใน go.mod:**
```go
module myproject

go 1.21 // แก้ไขตรงนี้ให้เป็นเวอร์ชันที่ต้องการ
```

**การจัดการ Dependencies เวอร์ชันต่างๆ:**
```bash
# อัปเดตแพคเกจทั้งหมดเป็นเวอร์ชันล่าสุด
go get -u ./...

# อัปเดตเฉพาะแพคเกจที่ต้องการ
go get -u github.com/gin-gonic/gin

# ดูเวอร์ชันของแพคเกจที่ใช้
go list -m -versions github.com/gin-gonic/gin
```

---

### บทที่ 41: คำแนะนำในการออกแบบโค้ดที่ดี

การเขียนโค้ดให้ดีไม่ใช่แค่ให้ทำงานได้ แต่ต้องอ่านง่าย, บำรุงรักษาง่าย, และมีประสิทธิภาพ

**1. ตั้งชื่อให้มีความหมาย:**
```go
// ไม่ดี
func calc(a, b int) int { ... }

// ดี
func calculateTotalPrice(basePrice, taxRate int) int { ... }
```

**2. ใช้ Receiver Names ที่สั้น:**
```go
// ดี
func (c *Client) GetUser(id int) User { ... }
func (s *Server) Start() error { ... }

// ไม่ดี
func (client *Client) GetUser(id int) User { ... }
```

**3. จัดการ Error ทันทีและไม่กลืน Error:**
```go
// ไม่ดี
data, _ := ioutil.ReadFile("config.json")

// ดี
data, err := ioutil.ReadFile("config.json")
if err != nil {
    return fmt.Errorf("อ่านไฟล์ config ไม่ได้: %w", err)
}
```

**4. ใช้ Interface ที่เล็กที่สุด (Interface Segregation):**
```go
// ไม่ดี (ใหญ่เกินไป)
type Database interface {
    Insert(data interface{}) error
    Update(id int, data interface{}) error
    Delete(id int) error
    Query(sql string) ([]interface{}, error)
}

// ดี (แยกตามการใช้งาน)
type Reader interface {
    Query(sql string) ([]interface{}, error)
}

type Writer interface {
    Insert(data interface{}) error
    Update(id int, data interface{}) error
    Delete(id int) error
}
```

**5. หลีกเลี่ยง Global Variables:**
```go
// ไม่ดี
var db *sql.DB

func init() { db, _ = sql.Open("mysql", "...") }

// ดี (Inject Dependency)
type App struct {
    db *sql.DB
}

func NewApp(db *sql.DB) *App {
    return &App{db: db}
}
```

**6. ใช้ Context สำหรับ Request-scoped Data:**
```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    // ส่ง ctx ไปยังฟังก์ชันที่เรียก
    result := processWithContext(ctx)
}
```

**7. อย่า Panic ใน Library:**
```go
// ไม่ดี (panics)
func MustParseJSON(data []byte) interface{} {
    var result interface{}
    if err := json.Unmarshal(data, &result); err != nil {
        panic(err)
    }
    return result
}

// ดี (return error)
func ParseJSON(data []byte) (interface{}, error) {
    var result interface{}
    if err := json.Unmarshal(data, &result); err != nil {
        return nil, err
    }
    return result, nil
}
```

**8. ใช้ `defer` เพื่อ Cleanup:**
```go
func readFile(filename string) ([]byte, error) {
    f, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer f.Close() // รับรองว่าปิดไฟล์เสมอ
    
    return ioutil.ReadAll(f)
}
```

---

### บทที่ 42: ชีทสรุป (Cheatsheet)

**การประกาศตัวแปร:**
```go
var name string = "Somchai"     // ระบุชนิด
var age = 30                     // Type Inference
salary := 25000.5                // Short Declaration (ใช้บ่อย)
const PI = 3.14159               // ค่าคงที่
```

**โครงสร้างข้อมูล:**
```go
// Array (ขนาดตายตัว)
arr := [3]int{1, 2, 3}

// Slice (ขนาดยืดหยุ่น)
slice := []int{1, 2, 3}
slice = append(slice, 4)

// Map
m := map[string]int{"a": 1, "b": 2}
```

**โครงสร้างควบคุม:**
```go
// if-else
if x > 0 {
    fmt.Println("positive")
} else if x < 0 {
    fmt.Println("negative")
} else {
    fmt.Println("zero")
}

// switch
switch day {
case "Mon":
    fmt.Println("Monday")
default:
    fmt.Println("Other")
}

// for (Go มีแค่ for)
for i := 0; i < 10; i++ { }     // คลาสสิก
for i < 10 { }                  // เหมือน while
for { }                         // infinite loop
for i, v := range slice { }     // range
```

**ฟังก์ชัน:**
```go
func add(a, b int) int { return a + b }
func div(a, b int) (int, error) { return a / b, nil }
func (r Rectangle) Area() float64 { return r.W * r.H }
```

**Concurrency:**
```go
go function()                   // Goroutine
ch := make(chan int)            // Channel
ch <- value                     // ส่งค่า
value := <-ch                  // รับค่า
select { case <-ch: }           // Select
var wg sync.WaitGroup          // WaitGroup
mu.Lock(); mu.Unlock()          // Mutex
```

**Error Handling:**
```go
result, err := doSomething()
if err != nil {
    log.Fatal(err)
}
```

**Interface:**
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

**Testing:**
```go
func TestAdd(t *testing.T) {
    if result != expected {
        t.Errorf("got %d, want %d", result, expected)
    }
}

func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ { add() }
}
```

**Common Packages:**
- `fmt` : formatting & printing
- `io/ioutil` : file I/O
- `net/http` : HTTP client & server
- `encoding/json` : JSON
- `time` : date & time
- `sync` : concurrency (Mutex, WaitGroup)
- `context` : context management
- `testing` : testing & benchmarking
- `os` : OS operations
- `strings` : string manipulation

---

### 📌 สรุปเล่มที่ 6
ในเล่มนี้คุณได้ก้าวขึ้นสู่ระดับมืออาชีพ:
✅ การเขียน Benchmark และวัดประสิทธิภาพโค้ด
✅ การสร้าง HTTP Client ที่ยืดหยุ่นและจัดการ Error
✅ การ Profiling CPU และ Memory เพื่อวิเคราะห์ประสิทธิภาพ
✅ การใช้ Context เพื่อควบคุมการทำงานและการยกเลิก
✅ การเขียน Generic Code ที่ reuse ได้กับหลายชนิดข้อมูล
✅ การทำความเข้าใจ OOP ในแบบฉบับ Go (Encapsulation, Composition, Polymorphism)
✅ การจัดการเวอร์ชัน Go และการอัปเกรด
✅ คำแนะนำการออกแบบโค้ดที่ดี (Best Practices)
✅ ชีทสรุปสำหรับอ้างอิงด่วน

 