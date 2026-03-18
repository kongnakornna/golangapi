# Full Code Examples for All Topics

This document provides concise code examples for each topic listed. The examples are primarily in Go (for networking, Go basics, OOP, SOLID) and SQL (for database operations). Each snippet is self-contained and demonstrates the core concept.

---

## หมวดที่ 1: พื้นฐานเครือข่ายและโปรโตคอล

### 1.1 ความรู้เบื้องต้นเกี่ยวกับเครือข่าย (Basic Network)
ไม่มีโค้ดเฉพาะ แต่เราจะใช้ Go `net` package เพื่อแสดงการทำงานกับ IP และ Port

### 1.2 เปรียบเทียบ TCP และ UDP

#### TCP Server และ Client (Go)

**TCP Server:**
```go
package main

import (
    "fmt"
    "net"
)

func main() {
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        panic(err)
    }
    defer listener.Close()
    fmt.Println("TCP server listening on :8080")

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println(err)
            continue
        }
        go handleTCPConnection(conn)
    }
}

func handleTCPConnection(conn net.Conn) {
    defer conn.Close()
    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("Received: %s\n", string(buf[:n]))
    conn.Write([]byte("Message received"))
}
```

**TCP Client:**
```go
package main

import (
    "fmt"
    "net"
)

func main() {
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    
    _, err = conn.Write([]byte("Hello TCP Server"))
    if err != nil {
        panic(err)
    }
    
    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Response: %s\n", string(buf[:n]))
}
```

#### UDP Server และ Client (Go)

**UDP Server:**
```go
package main

import (
    "fmt"
    "net"
)

func main() {
    addr, err := net.ResolveUDPAddr("udp", ":8081")
    if err != nil {
        panic(err)
    }
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    fmt.Println("UDP server listening on :8081")

    buf := make([]byte, 1024)
    for {
        n, clientAddr, err := conn.ReadFromUDP(buf)
        if err != nil {
            fmt.Println(err)
            continue
        }
        fmt.Printf("Received from %v: %s\n", clientAddr, string(buf[:n]))
        conn.WriteToUDP([]byte("Message received"), clientAddr)
    }
}
```

**UDP Client:**
```go
package main

import (
    "fmt"
    "net"
)

func main() {
    addr, err := net.ResolveUDPAddr("udp", "localhost:8081")
    if err != nil {
        panic(err)
    }
    conn, err := net.DialUDP("udp", nil, addr)
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    
    _, err = conn.Write([]byte("Hello UDP Server"))
    if err != nil {
        panic(err)
    }
    
    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Response: %s\n", string(buf[:n]))
}
```

### 1.3 โปรโตคอล HTTP

#### HTTP Server (Go)
```go
package main

import (
    "fmt"
    "net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func main() {
    http.HandleFunc("/", helloHandler)
    http.ListenAndServe(":8080", nil)
}
```

#### HTTP Client (Go)
```go
package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
)

func main() {
    resp, err := http.Get("http://localhost:8080/world")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### 1.4 โปรโตคอล MQTT สำหรับ IoT

ใช้ไลบรารี `paho.mqtt.golang`

```go
package main

import (
    "fmt"
    "time"
    mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
    opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
    opts.SetClientID("go_publisher")
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }

    // Publish
    token := client.Publish("test/topic", 0, false, "Hello MQTT")
    token.Wait()

    // Subscribe
    client.Subscribe("test/topic", 0, func(client mqtt.Client, msg mqtt.Message) {
        fmt.Printf("Received: %s from topic: %s\n", msg.Payload(), msg.Topic())
    })

    time.Sleep(2 * time.Second)
    client.Disconnect(250)
}
```

### 1.5 โปรโตคอล SNMP สำหรับจัดการเครือข่าย

ใช้ไลบรารี `gosnmp`

```go
package main

import (
    "fmt"
    "log"
    "github.com/gosnmp/gosnmp"
)

func main() {
    // SNMP GET example
    gosnmp.Default.Target = "192.168.1.1"
    gosnmp.Default.Community = "public"
    gosnmp.Default.Version = gosnmp.Version2c
    err := gosnmp.Default.Connect()
    if err != nil {
        log.Fatalf("Connect() err: %v", err)
    }
    defer gosnmp.Default.Conn.Close()

    oids := []string{"1.3.6.1.2.1.1.1.0"} // sysDescr
    result, err2 := gosnmp.Default.Get(oids)
    if err2 != nil {
        log.Fatalf("Get() err: %v", err2)
    }

    for _, variable := range result.Variables {
        fmt.Printf("oid: %s, value: %s\n", variable.Name, variable.Value)
    }
}
```

---

## หมวดที่ 2: สถาปัตยกรรมซอฟต์แวร์

### 2.1 สถาปัตยกรรมแบบ Monolithic

ตัวอย่างโครงสร้างโปรเจกต์ Go แบบ Monolithic (รวมทุกอย่างไว้ใน package เดียว)

```
myapp/
  main.go
  handlers/
    user.go
    product.go
  models/
    user.go
    product.go
  repository/
    user_repo.go
    product_repo.go
  service/
    user_service.go
    product_service.go
  go.mod
```

**main.go** (ตัวอย่างง่าย)
```go
package main

import (
    "fmt"
    "net/http"
    "myapp/handlers"
)

func main() {
    http.HandleFunc("/users", handlers.GetUsers)
    http.HandleFunc("/products", handlers.GetProducts)
    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", nil)
}
```

---

## หมวดที่ 3: พื้นฐานภาษา Go

### 3.1 แนะนำ Go และการติดตั้ง (ไม่มีโค้ด)

### 3.2 Go Package และการนำเข้า

```go
package main

import (
    "fmt"
    "math/rand"
)

func main() {
    fmt.Println("Random number:", rand.Intn(100))
}
```

### 3.3 ตัวแปรและชนิดข้อมูล (Variables)

```go
package main

import "fmt"

func main() {
    var name string = "Go"
    var version = 1.21
    isFun := true
    var x, y int = 1, 2
    const pi = 3.14

    fmt.Println(name, version, isFun, x, y, pi)
}
```

### 3.4 ตัวดำเนินการ (Operators)

```go
package main

import "fmt"

func main() {
    a, b := 10, 3
    fmt.Println("Arithmetic:", a+b, a-b, a*b, a/b, a%b)
    fmt.Println("Comparison:", a == b, a != b, a > b)
    fmt.Println("Logical:", (a > 5) && (b < 5), (a > 5) || (b > 5), !(a > 5))
    fmt.Println("Bitwise:", a&b, a|b, a^b, a<<1, a>>1)
}
```

### 3.5 โครงสร้างควบคุม (Control Flow)

```go
package main

import "fmt"

func main() {
    // if-else
    x := 10
    if x > 5 {
        fmt.Println("x > 5")
    } else {
        fmt.Println("x <= 5")
    }

    // switch
    switch x {
    case 1, 2:
        fmt.Println("one or two")
    case 10:
        fmt.Println("ten")
    default:
        fmt.Println("other")
    }

    // for loop (while style)
    sum := 0
    for sum < 10 {
        sum += 2
    }
    fmt.Println("sum:", sum)
}
```

### 3.6 ฟังก์ชัน (Function)

```go
package main

import "fmt"

func add(a, b int) int {
    return a + b
}

func swap(x, y string) (string, string) {
    return y, x
}

func main() {
    fmt.Println(add(5, 3))
    a, b := swap("hello", "world")
    fmt.Println(a, b)
}
```

### 3.7 การวนซ้ำและการ debug (Loop)

```go
package main

import "fmt"

func main() {
    // Classic for loop
    for i := 0; i < 5; i++ {
        fmt.Println(i)
    }

    // Range loop
    nums := []int{2, 4, 6}
    for idx, val := range nums {
        fmt.Printf("index: %d, value: %d\n", idx, val)
    }

    // Infinite loop with break
    count := 0
    for {
        fmt.Println("loop")
        count++
        if count == 3 {
            break
        }
    }

    // Debugging: print variables
    x := 42
    fmt.Printf("Debug: x = %d\n", x)
}
```

### 3.8 พอยน์เตอร์ (Pointers)

```go
package main

import "fmt"

func zeroVal(val int) {
    val = 0
}

func zeroPtr(ptr *int) {
    *ptr = 0
}

func main() {
    x := 5
    zeroVal(x)
    fmt.Println(x) // 5

    zeroPtr(&x)
    fmt.Println(x) // 0

    fmt.Println("pointer address:", &x)
}
```

### 3.9 อาร์เรย์และสไลซ์ (Array and Slice)

```go
package main

import "fmt"

func main() {
    // Array
    var arr [3]int = [3]int{1, 2, 3}
    fmt.Println(arr)

    // Slice
    slice := []int{4, 5, 6}
    slice = append(slice, 7)
    fmt.Println(slice)

    // Make slice
    s2 := make([]int, 5, 10)
    fmt.Println(len(s2), cap(s2))

    // Slicing
    sub := slice[1:3]
    fmt.Println(sub)
}
```

### 3.10 แมป (Map)

```go
package main

import "fmt"

func main() {
    // Declaration
    var m map[string]int
    m = make(map[string]int)
    m["one"] = 1
    m["two"] = 2
    fmt.Println(m)

    // Map literal
    colors := map[string]string{
        "red":   "#ff0000",
        "green": "#00ff00",
    }
    fmt.Println(colors["red"])

    // Check existence
    value, ok := colors["blue"]
    fmt.Println(value, ok)

    // Delete
    delete(colors, "red")
    fmt.Println(colors)
}
```

### 3.11 โครงสร้าง (Struct)

```go
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func (p Person) SayHello() {
    fmt.Printf("Hello, my name is %s\n", p.Name)
}

func (p *Person) HaveBirthday() {
    p.Age++
}

func main() {
    p1 := Person{Name: "Alice", Age: 30}
    fmt.Println(p1)
    p1.SayHello()
    p1.HaveBirthday()
    fmt.Println(p1.Age)

    p2 := Person{Name: "Bob"}
    p2.Age = 25
    fmt.Println(p2)
}
```

### 3.12 อินเทอร์เฟซ (Interface)

```go
package main

import "fmt"

type Greeter interface {
    Greet() string
}

type Dog struct{ Name string }
func (d Dog) Greet() string {
    return "Woof! I'm " + d.Name
}

type Cat struct{ Name string }
func (c Cat) Greet() string {
    return "Meow! I'm " + c.Name
}

func sayHello(g Greeter) {
    fmt.Println(g.Greet())
}

func main() {
    dog := Dog{Name: "Buddy"}
    cat := Cat{Name: "Kitty"}
    sayHello(dog)
    sayHello(cat)
}
```

### 3.13 เจเนอริก (Generics)

```go
package main

import "fmt"

func Map[T any, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

func main() {
    nums := []int{1, 2, 3}
    squares := Map(nums, func(x int) int { return x * x })
    fmt.Println(squares) // [1 4 3]

    strs := []string{"a", "b", "c"}
    lengths := Map(strs, func(s string) int { return len(s) })
    fmt.Println(lengths) // [1 1 1]
}
```

### 3.14 การเขียนโปรแกรม concurrent: Goroutines, Channels, Mutex

#### Goroutines
```go
package main

import (
    "fmt"
    "time"
)

func say(s string) {
    for i := 0; i < 3; i++ {
        time.Sleep(100 * time.Millisecond)
        fmt.Println(s)
    }
}

func main() {
    go say("world")
    say("hello")
}
```

#### Channels
```go
package main

import "fmt"

func sum(s []int, c chan int) {
    sum := 0
    for _, v := range s {
        sum += v
    }
    c <- sum
}

func main() {
    s := []int{7, 2, 8, -9, 4, 0}
    c := make(chan int)
    go sum(s[:len(s)/2], c)
    go sum(s[len(s)/2:], c)
    x, y := <-c, <-c
    fmt.Println(x, y, x+y)
}
```

#### Buffered Channel
```go
ch := make(chan int, 2)
ch <- 1
ch <- 2
fmt.Println(<-ch)
fmt.Println(<-ch)
```

#### Select
```go
package main

import "time"
import "fmt"

func main() {
    c1 := make(chan string)
    c2 := make(chan string)

    go func() {
        time.Sleep(1 * time.Second)
        c1 <- "one"
    }()
    go func() {
        time.Sleep(2 * time.Second)
        c2 <- "two"
    }()

    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-c1:
            fmt.Println("received", msg1)
        case msg2 := <-c2:
            fmt.Println("received", msg2)
        }
    }
}
```

#### Mutex
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

func (c *Counter) Inc() {
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
    var wg sync.WaitGroup
    counter := Counter{}

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Inc()
        }()
    }
    wg.Wait()
    fmt.Println("Final counter:", counter.Value())
}
```

---

## หมวดที่ 4: การเขียนโปรแกรมเชิงวัตถุ (OOP)

### 4.1 ความหมายและประโยชน์ของ OOP (ไม่มีโค้ด)

### 4.2 รู้จัก UML Diagram เบื้องต้น (ไม่มีโค้ด)

### 4.3 เสาหลักของ OOP (Pillars of OOP) ใน Go

#### Encapsulation
```go
package person

type Person struct {
    name string // unexported field
    age  int
}

func NewPerson(name string, age int) *Person {
    return &Person{name: name, age: age}
}

func (p *Person) GetName() string {
    return p.name
}

func (p *Person) SetName(name string) {
    p.name = name
}
```

#### Inheritance (via Embedding)
```go
package main

import "fmt"

type Animal struct {
    Name string
}

func (a Animal) Speak() {
    fmt.Println("I am an animal")
}

type Dog struct {
    Animal // embedded
    Breed  string
}

func main() {
    d := Dog{Animal: Animal{Name: "Buddy"}, Breed: "Lab"}
    d.Speak() // method promoted from Animal
    fmt.Println(d.Name)
}
```

#### Polymorphism (via Interfaces)
```go
package main

import "fmt"

type Speaker interface {
    Speak() string
}

type Cat struct{ Name string }
func (c Cat) Speak() string { return "Meow" }

type Robot struct{ Model string }
func (r Robot) Speak() string { return "Beep boop" }

func main() {
    var s Speaker
    s = Cat{Name: "Kitty"}
    fmt.Println(s.Speak())
    s = Robot{Model: "R2D2"}
    fmt.Println(s.Speak())
}
```

#### Abstraction
```go
package main

import "fmt"

// Database interface abstracts the underlying database
type Database interface {
    Save(data string)
}

type MySQL struct{}
func (m MySQL) Save(data string) {
    fmt.Println("Saving to MySQL:", data)
}

type Redis struct{}
func (r Redis) Save(data string) {
    fmt.Println("Saving to Redis:", data)
}

func StoreData(db Database, data string) {
    db.Save(data)
}

func main() {
    StoreData(MySQL{}, "hello")
    StoreData(Redis{}, "world")
}
```

### 4.4 ความสัมพันธ์ระหว่าง Objects

#### Association (uses-a)
```go
type Driver struct {
    Name string
}

type Car struct {
    Model string
}

func (c Car) Drive(d Driver) {
    fmt.Println(d.Name, "drives", c.Model)
}
```

#### Aggregation (has-a) – weak relationship
```go
type Team struct {
    Players []Player
}

type Player struct {
    Name string
}
```

#### Composition (part-of) – strong relationship
```go
type House struct {
    Rooms []Room // Rooms are created with House, destroyed with House
}

type Room struct {
    Name string
}
```

### 4.5 หลักการ SOLID (SOLID Principles) ใน Go

#### S: Single Responsibility Principle
```go
// Bad: one struct with multiple responsibilities
type User struct {
    Name string
    Age  int
}
func (u User) Save() { /* save to DB */ }
func (u User) SendEmail() { /* send email */ }

// Good: separate concerns
type User struct { Name string; Age int }
type UserRepository struct {}
func (r UserRepository) Save(u User) {}
type EmailService struct {}
func (e EmailService) SendEmail(u User) {}
```

#### O: Open/Closed Principle
```go
// Open for extension, closed for modification
type Notifier interface {
    Notify(message string)
}

type EmailNotifier struct{}
func (e EmailNotifier) Notify(msg string) {
    fmt.Println("Email:", msg)
}

type SMSNotifier struct{}
func (s SMSNotifier) Notify(msg string) {
    fmt.Println("SMS:", msg)
}

func SendAlert(n Notifier, msg string) {
    n.Notify(msg)
}
// New notifiers can be added without changing SendAlert
```

#### L: Liskov Substitution Principle
```go
type Bird interface {
    Fly() error
}

type Eagle struct{}
func (e Eagle) Fly() error { return nil }

type Penguin struct{}
func (p Penguin) Fly() error { return errors.New("can't fly") }
// Penguin violates LSP because it cannot truly substitute Bird.
// Better to have separate interfaces for flyable and non-flyable birds.
```

#### I: Interface Segregation Principle
```go
// Bad: fat interface
type Worker interface {
    Work()
    Eat()
}

// Good: segregated interfaces
type Workable interface {
    Work()
}
type Eatable interface {
    Eat()
}

type Human struct{}
func (h Human) Work() { fmt.Println("Working") }
func (h Human) Eat()  { fmt.Println("Eating") }

type Robot struct{}
func (r Robot) Work() { fmt.Println("Working") }
// Robot doesn't need Eat()
```

#### D: Dependency Inversion Principle
```go
// Depend on abstractions, not concretions
type MessageSender interface {
    Send(message string) error
}

type NotificationService struct {
    sender MessageSender // depends on interface
}

func (n NotificationService) Notify(msg string) {
    n.sender.Send(msg)
}

type EmailSender struct{}
func (e EmailSender) Send(msg string) error {
    fmt.Println("Sending email:", msg)
    return nil
}
```

---

## หมวดที่ 5: ฐานข้อมูล SQL

### 5.1 ความรู้เบื้องต้นเกี่ยวกับ SQL และความสัมพันธ์ (Relationship)
ไม่มีโค้ด

### 5.2 การติดตั้ง PostgreSQL บน Docker
```bash
docker run --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres
```

### 5.3 คำสั่ง SQL พื้นฐาน

#### สร้างตาราง (ตัวอย่าง)
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    age INT,
    email VARCHAR(255) UNIQUE
);
```

#### Insert
```sql
INSERT INTO users (name, age, email) VALUES ('Alice', 30, 'alice@example.com');
INSERT INTO users (name, age, email) VALUES ('Bob', 25, 'bob@example.com');
```

#### Select
```sql
SELECT * FROM users;
SELECT name, age FROM users WHERE age > 25;
```

#### Where
```sql
SELECT * FROM users WHERE name = 'Alice';
```

#### Like
```sql
SELECT * FROM users WHERE email LIKE '%@example.com';
```

#### And/Or
```sql
SELECT * FROM users WHERE age > 20 AND name LIKE 'A%';
SELECT * FROM users WHERE age < 20 OR age > 30;
```

#### Order By
```sql
SELECT * FROM users ORDER BY age DESC;
SELECT * FROM users ORDER BY name ASC;
```

### 5.4 การแก้ไขและลบข้อมูล: Update, Delete

#### Update
```sql
UPDATE users SET age = 31 WHERE name = 'Alice';
```

#### Delete
```sql
DELETE FROM users WHERE name = 'Bob';
```

### 5.5 การเชื่อมตารางด้วย Join

สมมติมีตาราง orders
```sql
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    product VARCHAR(255),
    amount DECIMAL
);
```

#### Inner Join
```sql
SELECT users.name, orders.product, orders.amount
FROM users
JOIN orders ON users.id = orders.user_id;
```

#### Left Join
```sql
SELECT users.name, orders.product
FROM users
LEFT JOIN orders ON users.id = orders.user_id;
```

### 5.6 การทำธุรกรรม (Transaction)

```sql
BEGIN;
INSERT INTO users (name, age, email) VALUES ('Charlie', 28, 'charlie@example.com');
UPDATE users SET age = 29 WHERE name = 'Charlie';
-- If everything is OK
COMMIT;
-- If error
ROLLBACK;
```

#### การใช้ Transaction ใน Go (database/sql)
```go
package main

import (
    "database/sql"
    "log"
    _ "github.com/lib/pq"
)

func main() {
    connStr := "user=postgres dbname=test password=mysecretpassword host=localhost sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }

    _, err = tx.Exec("INSERT INTO users (name, age, email) VALUES ($1, $2, $3)", "David", 22, "david@example.com")
    if err != nil {
        tx.Rollback()
        log.Fatal(err)
    }

    _, err = tx.Exec("UPDATE users SET age = $1 WHERE name = $2", 23, "David")
    if err != nil {
        tx.Rollback()
        log.Fatal(err)
    }

    err = tx.Commit()
    if err != nil {
        log.Fatal(err)
    }
}
```

---

**หมายเหตุ:** โค้ดทั้งหมดเป็นเพียงตัวอย่างเพื่อแสดงแนวคิด อาจต้องติดตั้งไลบรารีเพิ่มเติม (เช่น `go get github.com/eclipse/paho.mqtt.golang`, `go get github.com/gosnmp/gosnmp`, `go get github.com/lib/pq`) และต้องรัน PostgreSQL ก่อนทดสอบ SQL