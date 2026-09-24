ขอเสนอ Golang Tutorial แบบเข้มข้น ครอบคลุมตั้งแต่พื้นฐานจนถึง advanced ครับ

## 1. พื้นฐานภาษา Go

### โครงสร้างโปรแกรมพื้นฐาน
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### ตัวแปรและการประกาศ
```go
// วิธีประกาศตัวแปร
var name string = "John"
var age int = 30
var isActive bool = true

// Short declaration (ใช้บ่อย)
message := "Hello"  // type inference
x, y := 10, 20

// Constant
const PI = 3.14159
const (
    StatusOK = 200
    StatusNotFound = 404
)
```

### Data Types พื้นฐาน
```go
// Numbers
var i int = 42
var f float64 = 3.14
var b byte = 255  // uint8

// String
var s string = "Go"

// Boolean
var flag bool = true

// Array (fixed size)
var arr [3]int = [3]int{1, 2, 3}

// Slice (dynamic)
var slice []int = []int{1, 2, 3, 4, 5}
slice = append(slice, 6)

// Map
var m map[string]int = make(map[string]int)
m["key"] = 100

// Struct
type Person struct {
    Name string
    Age  int
}
p := Person{Name: "Alice", Age: 25}
```

## 2. Control Flow

### Conditionals
```go
// if-else
if x > 10 {
    fmt.Println("x is greater than 10")
} else if x == 10 {
    fmt.Println("x equals 10")
} else {
    fmt.Println("x is less than 10")
}

// if with short statement
if err := process(); err != nil {
    fmt.Println("Error:", err)
}

// switch (ไม่ต้องมี break)
switch day {
case "Monday":
    fmt.Println("Start of week")
case "Friday":
    fmt.Println("TGIF!")
default:
    fmt.Println("Other day")
}

// switch without condition
switch {
case score >= 90:
    grade = "A"
case score >= 80:
    grade = "B"
default:
    grade = "C"
}
```

### Loops
```go
// for loop (แบบดั้งเดิม)
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// while-style
count := 0
for count < 5 {
    fmt.Println(count)
    count++
}

// infinite loop
for {
    // break to exit
}

// range loop
numbers := []int{1, 2, 3, 4, 5}
for index, value := range numbers {
    fmt.Printf("Index: %d, Value: %d\n", index, value)
}

// range with map
for key, value := range m {
    fmt.Printf("%s: %d\n", key, value)
}
```

## 3. Functions

### Function Basics
```go
// Basic function
func add(x int, y int) int {
    return x + y
}

// Multiple return values
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Named return values
func getCoordinates() (x, y int) {
    x = 10
    y = 20
    return  // naked return
}

// Variadic functions
func sum(numbers ...int) int {
    total := 0
    for _, n := range numbers {
        total += n
    }
    return total
}

// Function as value
var fn func(int, int) int = add

// Closure
func counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}
```

## 4. Methods และ Interfaces

### Methods
```go
type Rectangle struct {
    Width  float64
    Height float64
}

// Value receiver
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Pointer receiver (สามารถ modify struct)
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}
```

### Interfaces
```go
// Interface definition
type Shape interface {
    Area() float64
    Perimeter() float64
}

//  implement interface โดยอัตโนมัติ
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

// Empty interface (รับค่าได้ทุกประเภท)
func printAnything(v interface{}) {
    fmt.Println(v)
}

// Type assertion
func process(i interface{}) {
    // Type assertion
    if val, ok := i.(string); ok {
        fmt.Println("String:", val)
    }
    
    // Type switch
    switch v := i.(type) {
    case int:
        fmt.Println("Integer:", v)
    case string:
        fmt.Println("String:", v)
    default:
        fmt.Println("Unknown type")
    }
}
```

## 5. Error Handling

### Error Pattern
```go
// Custom error
type ValidationError struct {
    Field string
    Value interface{}
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %v", e.Field, e.Value)
}

// Error handling pattern
func validateUser(user User) error {
    if user.Name == "" {
        return ValidationError{Field: "name", Value: user.Name}
    }
    if user.Age < 0 {
        return ValidationError{Field: "age", Value: user.Age}
    }
    return nil
}

// Panic and Recover
func safeDivision(a, b int) (result int) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
            result = 0
        }
    }()
    
    if b == 0 {
        panic("division by zero")
    }
    return a / b
}
```

## 6. Concurrency

### Goroutines
```go
// Basic goroutine
go func() {
    fmt.Println("Running in goroutine")
}()

// Wait for goroutines
func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Printf("Worker %d started\n", id)
    time.Sleep(time.Second)
    fmt.Printf("Worker %d finished\n", id)
}

func main() {
    var wg sync.WaitGroup
    
    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }
    
    wg.Wait()
}
```

### Channels
```go
// Unbuffered channel
ch := make(chan int)

// Send and receive
go func() {
    ch <- 42  // send
}()
value := <-ch  // receive

// Buffered channel
ch := make(chan string, 3)
ch <- "message1"
ch <- "message2"

// Channel with range
ch := make(chan int)
go func() {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch)
}()
for value := range ch {
    fmt.Println(value)
}

// Select statement
select {
case msg1 := <-ch1:
    fmt.Println("Received from ch1:", msg1)
case msg2 := <-ch2:
    fmt.Println("Received from ch2:", msg2)
case <-time.After(1 * time.Second):
    fmt.Println("Timeout")
default:
    fmt.Println("No channels ready")
}

// Worker pool pattern
func workerPool(jobs <-chan int, results chan<- int) {
    for job := range jobs {
        results <- job * 2
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)
    
    // Start workers
    for w := 0; w < 3; w++ {
        go workerPool(jobs, results)
    }
    
    // Send jobs
    for j := 0; j < 10; j++ {
        jobs <- j
    }
    close(jobs)
    
    // Collect results
    for r := 0; r < 10; r++ {
        <-results
    }
}
```

### Mutex และ Sync
```go
type SafeCounter struct {
    mu    sync.Mutex
    value int
}

func (c *SafeCounter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

// Once
var once sync.Once
once.Do(func() {
    fmt.Println("This runs only once")
})

// WaitGroup (shown above)
// Cond
// Pool
```

## 7. Packages และ Modules

### Module Management
```bash
# Initialize module
go mod init example.com/myproject

# Add dependencies
go get github.com/gin-gonic/gin

# Update dependencies
go get -u ./...

# Tidy modules
go mod tidy
```

### Package Structure
```
myproject/
├── go.mod
├── go.sum
├── main.go
├── internal/
│   └── auth/
│       └── auth.go
├── pkg/
│   └── utils/
│       └── utils.go
└── api/
    └── handler.go
```

## 8. Advanced Patterns

### Context
```go
func longRunningProcess(ctx context.Context, duration time.Duration) error {
    select {
    case <-time.After(duration):
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    err := longRunningProcess(ctx, 3*time.Second)
    if err != nil {
        fmt.Println("Process cancelled:", err)
    }
}
```

### Generics (Go 1.18+)
```go
// Generic function
func Map[T, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

// Generic type
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

// Constraint
type Number interface {
    int | int64 | float64
}

func Sum[T Number](numbers []T) T {
    var sum T
    for _, n := range numbers {
        sum += n
    }
    return sum
}
```

### Testing
```go
// main_test.go
func TestAdd(t *testing.T) {
    result := add(2, 3)
    expected := 5
    
    if result != expected {
        t.Errorf("add(2,3) = %d; want %d", result, expected)
    }
}

// Table-driven test
func TestDivide(t *testing.T) {
    tests := []struct {
        name     string
        a, b     float64
        expected float64
        hasError bool
    }{
        {"positive numbers", 10, 2, 5, false},
        {"divide by zero", 10, 0, 0, true},
        {"negative numbers", -10, 2, -5, false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := divide(tt.a, tt.b)
            if tt.hasError && err == nil {
                t.Error("expected error but got none")
            }
            if !tt.hasError && result != tt.expected {
                t.Errorf("got %f, want %f", result, tt.expected)
            }
        })
    }
}

// Benchmark
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        add(10, 20)
    }
}
```

## 9. Best Practices

1. **Code Organization**
   - Keep packages focused and small
   - Use internal packages for private code
   - Follow standard project layout

2. **Error Handling**
   - Always handle errors
   - Wrap errors with context
   - Never ignore errors

3. **Concurrency**
   - "Don't communicate by sharing memory; share memory by communicating"
   - Use goroutines for concurrent tasks
   - Be careful with goroutine leaks

4. **Performance**
   - Use pointers for large structs
   - Pre-allocate slices when size is known
   - Use sync.Pool for frequent allocations

5. **Naming Conventions**
   - CamelCase for exported identifiers
   - camelCase for unexported
   - Acronyms: URL, HTTP (not Url, Http)

## 10. Useful Commands

```bash
# Build
go build
go build -o myapp
go build -ldflags "-X main.version=1.0.0"

# Run
go run main.go
go run .

# Test
go test
go test -v
go test -cover
go test -bench=.

# Format
go fmt ./...

# Lint
go vet ./...
golangci-lint run

# Documentation
go doc fmt.Println
godoc -http=:6060
```

## Example: Complete REST API

```go
package main

import (
    "encoding/json"
    "net/http"
    "sync"
    "github.com/gorilla/mux"
)

type Todo struct {
    ID   string `json:"id"`
    Text string `json:"text"`
    Done bool   `json:"done"`
}

type TodoService struct {
    mu    sync.RWMutex
    items map[string]Todo
}

func NewTodoService() *TodoService {
    return &TodoService{
        items: make(map[string]Todo),
    }
}

func (s *TodoService) CreateTodo(w http.ResponseWriter, r *http.Request) {
    var todo Todo
    if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    s.mu.Lock()
    s.items[todo.ID] = todo
    s.mu.Unlock()
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(todo)
}

func (s *TodoService) GetTodos(w http.ResponseWriter, r *http.Request) {
    s.mu.RLock()
    todos := make([]Todo, 0, len(s.items))
    for _, todo := range s.items {
        todos = append(todos, todo)
    }
    s.mu.RUnlock()
    
    json.NewEncoder(w).Encode(todos)
}

func main() {
    router := mux.NewRouter()
    service := NewTodoService()
    
    router.HandleFunc("/todos", service.GetTodos).Methods("GET")
    router.HandleFunc("/todos", service.CreateTodo).Methods("POST")
    
    http.ListenAndServe(":8080", router)
}
```

นี่คือเนื้อหาครอบคลุมที่สำคัญของ Golang ครับ ควรฝึกปฏิบัติและทำโปรเจคจริงเพื่อความเข้าใจที่ลึกซึ้งยิ่งขึ้น