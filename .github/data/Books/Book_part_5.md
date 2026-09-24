## 📚 เล่มที่ 5: การพัฒนาแอปพลิเคชันเชิงปฏิบัติ

เนื้อหาในเล่มนี้จะเปลี่ยนคุณจากผู้ที่เขียนโค้ดพื้นฐาน มาเป็นผู้ที่สามารถสร้างแอปพลิเคชันที่พร้อมใช้งานจริงได้ ด้วยการเรียนรู้การจัดการข้อมูลรูปแบบต่างๆ การสร้าง HTTP Server การทำงานแบบ Concurrent การจัดการ Log และการตั้งค่า Configuration อย่างมืออาชีพ

---

### บทที่ 24: ฟังก์ชันนิรนาม (Anonymous Functions) และ Closure

**ฟังก์ชันนิรนาม** คือฟังก์ชันที่ไม่มีชื่อ ใช้สำหรับงานที่ต้องการฟังก์ชันเฉพาะครั้ง เช่น ส่งต่อให้ฟังก์ชันอื่น หรือกำหนดเป็นตัวแปร

```go
package main

import "fmt"

func main() {
    // ประกาศฟังก์ชันนิรนามและเก็บไว้ในตัวแปร
    greeting := func(name string) string {
        return "สวัสดี " + name
    }

    fmt.Println(greeting("สมชาย")) // สวัสดี สมชาย

    // เรียกใช้ทันที (Immediately Invoked Function Expression - IIFE)
    result := func(a, b int) int {
        return a + b
    }(5, 3)
    fmt.Println(result) // 8
}
```

**Closure (การปิดล้อมตัวแปร):**
Closure คือฟังก์ชันนิรนามที่สามารถ **เข้าถึงและจดจำตัวแปร** ที่อยู่ภายนอกขอบเขตของมันได้ แม้ตัวแปรนั้นจะอยู่นอกฟังก์ชันก็ตาม

```go
func counter() func() int {
    count := 0 // ตัวแปรนี้จะถูก "ปิดล้อม" ไว้
    return func() int {
        count++ // ฟังก์ชันนี้สามารถเข้าถึง count ได้
        return count
    }
}

func main() {
    c1 := counter()
    fmt.Println(c1()) // 1
    fmt.Println(c1()) // 2
    fmt.Println(c1()) // 3

    c2 := counter() // สร้างตัวนับใหม่ แยกจากกัน
    fmt.Println(c2()) // 1
}
```

**การใช้งานจริง:** ใช้ใน HTTP Middleware, การจัดการ State ใน Goroutine, หรือการทำ Function Factory

---

### บทที่ 25: การจัดการข้อมูล JSON และ XML

JSON และ XML เป็นรูปแบบข้อมูลที่ใช้กันแพร่หลายในการสื่อสารระหว่างระบบ (API) Go มีแพคเกจ `encoding/json` และ `encoding/xml` สำหรับการแปลงข้อมูลระหว่าง Struct กับรูปแบบเหล่านี้โดยตรง

**การทำงานกับ JSON:**
```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
)

// ใช้ Struct Tags เพื่อกำหนดชื่อฟิลด์ใน JSON
type Person struct {
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Email   string `json:"email,omitempty"` // omitempty = ข้ามถ้าค่าว่าง
    Password string `json:"-"`              // ข้ามฟิลด์นี้ไปเลย
}

func main() {
    // แปลง Struct -> JSON (Marshal)
    p := Person{Name: "สมชาย", Age: 30, Password: "secret"}
    jsonData, err := json.Marshal(p)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(jsonData)) 
    // {"name":"สมชาย","age":30}  (Password ถูกลบ, Email ไม่มี)

    // แปลง JSON -> Struct (Unmarshal)
    jsonStr := `{"name":"สมศรี","age":25,"email":"somsri@mail.com"}`
    var person Person
    err = json.Unmarshal([]byte(jsonStr), &person)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%+v\n", person) // {Name:สมศรี Age:25 Email:somsri@mail.com Password:}
    
    // จัดการ JSON ที่ไม่รู้โครงสร้างล่วงหน้า (ใช้ map)
    var result map[string]interface{}
    json.Unmarshal([]byte(`{"status":"ok","code":200}`), &result)
    fmt.Println(result["status"]) // ok
}
```

**การทำงานกับ XML (คล้ายกัน):**
```go
import "encoding/xml"

type Product struct {
    XMLName xml.Name `xml:"product"`
    ID      int      `xml:"id,attr"`  // attribute
    Name    string   `xml:"name"`
    Price   float64  `xml:"price"`
}

func main() {
    p := Product{ID: 1, Name: "โน้ตบุ๊ค", Price: 25000.0}
    xmlData, _ := xml.MarshalIndent(p, "", "  ")
    fmt.Println(string(xmlData))
    // <product id="1">
    //   <name>โน้ตบุ๊ค</name>
    //   <price>25000</price>
    // </product>
}
```

---

### บทที่ 26: พื้นฐานการสร้าง HTTP Server

Go มีแพคเกจ `net/http` ในตัวที่ทรงพลังมาก เราสามารถสร้าง HTTP Server ที่ทำงานได้จริงด้วยโค้ดไม่กี่บรรทัด

**HTTP Server พื้นฐาน:**
```go
package main

import (
    "fmt"
    "net/http"
)

// Handler ฟังก์ชันที่รับ Request และส่ง Response
func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "สวัสดีจาก Go Server!")
}

func main() {
    // กำหนดเส้นทาง (Route)
    http.HandleFunc("/", helloHandler)
    http.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {
        name := r.URL.Query().Get("name") // รับ Query Parameter
        fmt.Fprintf(w, "สวัสดี %s", name)
    })

    // เริ่ม Server ที่พอร์ต 8080
    fmt.Println("Server เริ่มที่ http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

**การแยก HTTP Method (GET, POST):**
```go
func userHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        fmt.Fprintln(w, "GET: ดึงข้อมูลผู้ใช้")
    case http.MethodPost:
        // อ่าน Body
        var data map[string]interface{}
        json.NewDecoder(r.Body).Decode(&data)
        fmt.Fprintf(w, "POST: รับข้อมูล %v", data)
    default:
        http.Error(w, "Method ไม่รองรับ", http.StatusMethodNotAllowed)
    }
}

func main() {
    http.HandleFunc("/user", userHandler)
    http.ListenAndServe(":8080", nil)
}
```

**การใช้ ServeMux (Router) แยกต่างหาก:**
```go
mux := http.NewServeMux()
mux.HandleFunc("/api/v1/users", getUsers)

server := &http.Server{
    Addr:    ":8080",
    Handler: mux,
}
server.ListenAndServe()
```

---

### บทที่ 27: Enum, Iota และ Bitmask

Go ไม่มี Enum ในตัวโดยตรง แต่ใช้ `iota` ซึ่งเป็นตัวนับอัตโนมัติในบล็อก `const` เพื่อสร้างค่าคงที่เรียงลำดับ

**การใช้ iota สร้าง Enum:**
```go
package main

import "fmt"

// สร้าง Enum สำหรับสถานะ
type Status int

const (
    Pending Status = iota // 0
    Active                // 1
    Inactive              // 2
    Suspended             // 3
)

func (s Status) String() string {
    names := []string{"รอดำเนินการ", "ใช้งานอยู่", "ไม่ใช้งาน", "ถูกระงับ"}
    return names[s]
}

func main() {
    status := Active
    fmt.Println(status) // ใช้งานอยู่
}
```

**การใช้ iota กับ Bitmask (การเก็บหลายสถานะในตัวแปรเดียว):**
```go
type Permission int

const (
    Read Permission = 1 << iota // 1  (0001)
    Write                       // 2  (0010)
    Execute                     // 4  (0100)
    Admin                       // 8  (1000)
)

func main() {
    // กำหนดสิทธิ์: อ่าน + เขียน
    var perms Permission = Read | Write
    fmt.Printf("สิทธิ์: %b\n", perms) // 0011

    // ตรวจสอบว่ามีสิทธิ์ Admin หรือไม่
    if perms&Admin != 0 {
        fmt.Println("มีสิทธิ์ Admin")
    } else {
        fmt.Println("ไม่มีสิทธิ์ Admin")
    }

    // เพิ่มสิทธิ์ Execute
    perms |= Execute
    fmt.Printf("สิทธิ์ใหม่: %b\n", perms) // 0111
}
```

---

### บทที่ 28: วันที่และเวลา

Go มีแพคเกจ `time` สำหรับจัดการวันที่และเวลาอย่างครบครัน

**การสร้างและการจัดการเวลา:**
```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // เวลาปัจจุบัน
    now := time.Now()
    fmt.Println("ปัจจุบัน:", now)

    // สร้างเวลาที่กำหนด
    birthDate := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)
    fmt.Println("วันเกิด:", birthDate)

    // แยกส่วนประกอบ
    year, month, day := now.Date()
    hour, min, sec := now.Clock()
    fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n", year, month, day, hour, min, sec)

    // การคำนวณเวลา
    future := now.Add(24 * time.Hour) // พรุ่งนี้
    fmt.Println("พรุ่งนี้:", future)

    // ตรวจสอบว่าระหว่างช่วงเวลาหรือไม่
    if now.After(birthDate) {
        fmt.Println("คุณอายุเกิน 0 ปีแล้ว")
    }

    // การ Format และ Parse
    layout := "2006-01-02 15:04:05"
    formatted := now.Format(layout)
    fmt.Println("รูปแบบ:", formatted)

    parsed, _ := time.Parse(layout, "2023-12-25 10:30:00")
    fmt.Println("Parse สำเร็จ:", parsed)
}
```
*หมายเหตุ:* Go ใช้ Reference Date `2006-01-02 15:04:05` (ตัวเลขเหล่านี้) ในการกำหนดรูปแบบ ไม่ใช่ `YYYY-MM-DD`

**การใช้ Timer และ Ticker:**
```go
// Timer: ทำงานครั้งเดียวหลังจากเวลาที่กำหนด
timer := time.NewTimer(2 * time.Second)
<-timer.C // รอจนกว่าครบ 2 วินาที
fmt.Println("ครบ 2 วินาทีแล้ว")

// Ticker: ทำงานซ้ำเป็นระยะ
ticker := time.NewTicker(1 * time.Second)
for i := 0; i < 5; i++ {
    <-ticker.C
    fmt.Println("ผ่านไป 1 วินาที")
}
ticker.Stop()
```

---

### บทที่ 29: การจัดเก็บข้อมูล: ไฟล์และฐานข้อมูล

**การอ่าน/เขียนไฟล์:**
```go
package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "os"
)

func main() {
    // เขียนไฟล์ (ทั้งไฟล์)
    data := []byte("สวัสดีโลก Go")
    err := ioutil.WriteFile("test.txt", data, 0644)
    if err != nil {
        panic(err)
    }

    // อ่านไฟล์ (ทั้งไฟล์)
    content, _ := ioutil.ReadFile("test.txt")
    fmt.Println(string(content))

    // อ่านไฟล์ทีละบรรทัด (เหมาะกับไฟล์ใหญ่)
    file, _ := os.Open("test.txt")
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        fmt.Println("บรรทัด:", scanner.Text())
    }

    // ต่อท้ายไฟล์
    f, _ := os.OpenFile("test.txt", os.O_APPEND|os.O_WRONLY, 0644)
    f.WriteString("\nเพิ่มบรรทัดใหม่")
    f.Close()
}
```

**การเชื่อมต่อฐานข้อมูล (SQL):**
Go รองรับ SQL ผ่าน `database/sql` และต้องใช้ Driver เฉพาะ เช่น `go-sqlite3`, `mysql`, `postgres`

```go
import (
    "database/sql"
    _ "github.com/go-sqlite3" // ใช้ _ เพื่อ import เฉพาะ init()
)

func main() {
    // เชื่อมต่อ SQLite
    db, err := sql.Open("sqlite3", "myapp.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // สร้างตาราง
    createSQL := `CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        age INTEGER
    )`
    db.Exec(createSQL)

    // Insert ข้อมูล
    result, _ := db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", "สมชาย", 30)
    id, _ := result.LastInsertId()
    fmt.Println("新增 ID:", id)

    // Query ข้อมูล
    rows, _ := db.Query("SELECT id, name, age FROM users")
    defer rows.Close()
    for rows.Next() {
        var id int
        var name string
        var age int
        rows.Scan(&id, &name, &age)
        fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
    }
}
```

---

### บทที่ 30: การทำงานพร้อมกัน (Concurrency)

หัวใจสำคัญของ Go คือ **Goroutine** (เธรดน้ำหนักเบา) และ **Channel** (ช่องทางสื่อสาร) ทำให้การเขียนโปรแกรม Concurrent ง่ายและปลอดภัย

**Goroutine:** ใช้คีย์เวิร์ด `go` นำหน้าฟังก์ชัน
```go
func printNumbers() {
    for i := 0; i < 5; i++ {
        time.Sleep(100 * time.Millisecond)
        fmt.Print(i, " ")
    }
}

func main() {
    go printNumbers() // ทำงานแบบขนาน
    go printNumbers()
    time.Sleep(2 * time.Second) // รอให้ทำงานเสร็จ
    fmt.Println("จบ")
}
```

**Channel:** ใช้สื่อสารระหว่าง Goroutine
```go
func sum(nums []int, ch chan int) {
    total := 0
    for _, n := range nums {
        total += n
    }
    ch <- total // ส่งค่าผ่าน Channel
}

func main() {
    ch := make(chan int)
    go sum([]int{1, 2, 3, 4}, ch)
    go sum([]int{5, 6, 7, 8}, ch)

    result1 := <-ch // รับค่า
    result2 := <-ch
    fmt.Println(result1, result2) // 10, 26 (อาจสลับกัน)
}
```

**Buffered Channel:** มีพื้นที่เก็บข้อมูล
```go
ch := make(chan string, 2) // เก็บได้ 2 ค่า
ch <- "hello"
ch <- "world"
fmt.Println(<-ch) // hello
```

**Select:** รอหลาย Channel พร้อมกัน
```go
ch1 := make(chan string)
ch2 := make(chan string)

go func() { time.Sleep(1 * time.Second); ch1 <- "จาก ch1" }()
go func() { time.Sleep(2 * time.Second); ch2 <- "จาก ch2" }()

select {
case msg1 := <-ch1:
    fmt.Println(msg1)
case msg2 := <-ch2:
    fmt.Println(msg2)
}
```

**WaitGroup:** รอให้ Goroutine หลายตัวทำงานเสร็จ
```go
var wg sync.WaitGroup

func work(id int) {
    defer wg.Done() // เมื่อทำงานเสร็จ
    fmt.Println("Goroutine", id, "ทำงาน")
    time.Sleep(1 * time.Second)
}

func main() {
    for i := 0; i < 5; i++ {
        wg.Add(1) // เพิ่มตัวนับ
        go work(i)
    }
    wg.Wait() // รอจนครบ
    fmt.Println("ทุกตัวเสร็จสิ้น")
}
```

**Mutex:** ป้องกันการเข้าถึงข้อมูลพร้อมกัน
```go
var counter int
var mu sync.Mutex

func increment() {
    mu.Lock()
    counter++
    mu.Unlock()
}
```

---

### บทที่ 31: การบันทึกเหตุการณ์ (Logging)

การบันทึก Log ช่วยในการดีบักและตรวจสอบระบบ Go มีแพคเกจ `log` ในตัว และยังมีไลบรารีอื่นๆ เช่น `logrus`, `zap`

**การใช้ log แบบพื้นฐาน:**
```go
import "log"

func main() {
    // ตั้งค่าไฟล์ Log
    file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    log.SetOutput(file)
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

    log.Println("แอปพลิเคชันเริ่มทำงาน")
    log.Printf("ผู้ใช้ %s เข้าสู่ระบบ", "สมชาย")

    // Log ระดับ Error
    log.Fatal("เกิดข้อผิดพลาดร้ายแรง") // บันทึกและจบการทำงาน
}
```

**การสร้าง Logger แยกตามระดับ (ใช้ logrus):**
```go
import "github.com/sirupsen/logrus"

func main() {
    log := logrus.New()
    log.SetFormatter(&logrus.JSONFormatter{}) // ใช้ JSON format

    log.WithFields(logrus.Fields{
        "user": "สมชาย",
        "ip":   "192.168.1.1",
    }).Info("ผู้ใช้เข้าสู่ระบบ")
}
```

---

### บทที่ 32: เทมเพลต (Templates)

Go มีแพคเกจ `html/template` และ `text/template` สำหรับสร้างข้อความหรือ HTML แบบไดนามิก ปลอดภัยจาก XSS

**การใช้ text/template:**
```go
import (
    "os"
    "text/template"
)

type User struct {
    Name string
    Age  int
}

func main() {
    tmpl := "สวัสดี {{.Name}} อายุ {{.Age}} ปี"
    t, _ := template.New("greeting").Parse(tmpl)
    user := User{Name: "สมชาย", Age: 30}
    t.Execute(os.Stdout, user) // สวัสดี สมชาย อายุ 30 ปี
}
```

**การใช้ html/template (พร้อมการ Escape อัตโนมัติ):**
```go
import "html/template"

func main() {
    tmpl := `<h1>สวัสดี {{.}}</h1>`
    t, _ := template.New("html").Parse(tmpl)
    t.Execute(os.Stdout, "<script>alert('XSS')</script>") 
    // จะถูก Escape เป็น &lt;script&gt;alert(...) เพื่อความปลอดภัย
}
```

---

### บทที่ 33: การจัดการค่า Configuration

การตั้งค่าควรแยกออกจากโค้ด เพื่อให้ปรับเปลี่ยนได้ง่ายตามสภาพแวดล้อม (Dev/Prod) นิยมใช้ไฟล์ `.env`, `config.json`, `yaml` หรือ `viper`

**การอ่าน Environment Variables:**
```go
package main

import (
    "fmt"
    "os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // ค่าเริ่มต้น
    }
    fmt.Println("Server รันที่พอร์ต:", port)

    // ตั้งค่า Environment Variable ใน Go
    os.Setenv("DB_HOST", "localhost")
    dbHost := os.Getenv("DB_HOST")
    fmt.Println("DB Host:", dbHost)
}
```

**การอ่านไฟล์ JSON Config:**
```go
type Config struct {
    Server struct {
        Port string `json:"port"`
    } `json:"server"`
    Database struct {
        Host string `json:"host"`
        User string `json:"user"`
    } `json:"database"`
}

func LoadConfig(path string) (*Config, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var config Config
    err = json.Unmarshal(data, &config)
    return &config, err
}
```

**การใช้ Viper (ไลบรารียอดนิยม):**
```go
import "github.com/spf13/viper"

func main() {
    viper.SetConfigName("config")   // config.json
    viper.SetConfigType("json")
    viper.AddConfigPath(".")

    viper.AutomaticEnv() // ใช้ Environment Variables ทับค่าในไฟล์ได้

    if err := viper.ReadInConfig(); err != nil {
        panic(err)
    }

    port := viper.GetString("server.port")
    dbHost := viper.GetString("database.host")
    fmt.Println(port, dbHost)
}
```

---

### 📌 สรุปเล่มที่ 5
ในเล่มนี้คุณได้เรียนรู้การพัฒนาแอปพลิเคชันในโลกจริง:
✅ ฟังก์ชันนิรนามและ Closure สำหรับการเขียนโค้ดที่ยืดหยุ่น
✅ การ Marshal/Unmarshal JSON และ XML
✅ การสร้าง HTTP Server และจัดการ Routing
✅ การสร้าง Enum และ Bitmask ด้วย iota
✅ การจัดการวันที่และเวลาอย่างครอบคลุม
✅ การอ่าน/เขียนไฟล์ และเชื่อมต่อฐานข้อมูล SQL
✅ การเขียนโปรแกรม Concurrent ด้วย Goroutine, Channel, WaitGroup, Mutex
✅ การบันทึก Log อย่างเป็นระบบ
✅ การสร้างเทมเพลตสำหรับข้อความและ HTML
✅ การจัดการ Configuration ด้วย Environment Variables และไฟล์

---
 