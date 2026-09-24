## 📚 เล่มที่ 7: เครื่องมือและไลบรารียอดนิยม

เล่มนี้จะพาคุณรู้จักกับไลบรารีและเครื่องมือที่ขาดไม่ได้ในการพัฒนาแอปพลิเคชัน Go สมัยใหม่ ตั้งแต่ Web Framework, การจัดการ Configuration, CLI Tool, Logging ระดับสูง, ไปจนถึง ORM และการส่งอีเมล

---

### บทที่ 43: chi, viper, cobra, zap และเครื่องมือสำคัญ

#### 43.1 chi - HTTP Router ที่เบาและเร็ว

`chi` เป็น lightweight HTTP router ที่มีประสิทธิภาพสูง เหมาะสำหรับการสร้าง REST API รองรับ Middleware, Sub-routing, และการจับพารามิเตอร์ในเส้นทางได้อย่างยืดหยุ่น

**การติดตั้ง:**
```bash
go get github.com/go-chi/chi/v5
```

**การใช้งานพื้นฐาน:**
```go
package main

import (
    "fmt"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    r := chi.NewRouter()

    // ใช้ Middleware พื้นฐาน
    r.Use(middleware.Logger)    // บันทึก Log การ Request
    r.Use(middleware.Recoverer) // จัดการ Panic

    // กำหนด Routes
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("สวัสดีจาก chi!"))
    })

    // Route พร้อมพารามิเตอร์
    r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
        userID := chi.URLParam(r, "id")
        fmt.Fprintf(w, "ผู้ใช้ ID: %s", userID)
    })

    // Sub-routing
    r.Route("/api/v1", func(r chi.Router) {
        r.Get("/products", listProducts)
        r.Post("/products", createProduct)
        r.Get("/products/{id}", getProduct)
    })

    http.ListenAndServe(":8080", r)
}
```

**การสร้าง Middleware แบบกำหนดเอง:**
```go
// Middleware ตรวจสอบ Authentication
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "ไม่มี Token", http.StatusUnauthorized)
            return
        }
        // ตรวจสอบ Token...
        next.ServeHTTP(w, r)
    })
}

// ใช้งาน
r.Use(AuthMiddleware)
```

---

#### 43.2 viper - การจัดการ Configuration ระดับสูง

`viper` เป็นไลบรารีสำหรับจัดการ Configuration ที่รองรับหลายรูปแบบ ทั้ง JSON, YAML, TOML, ENV, และ Command-line flags

**การติดตั้ง:**
```bash
go get github.com/spf13/viper
```

**การใช้งานพื้นฐาน:**
```go
package main

import (
    "fmt"
    "github.com/spf13/viper"
)

func main() {
    // ตั้งค่า Viper
    viper.SetConfigName("config")   // ชื่อไฟล์ (ไม่ต้องมีนามสกุล)
    viper.SetConfigType("yaml")     // ประเภทไฟล์
    viper.AddConfigPath(".")        // หาไฟล์ในโฟลเดอร์ปัจจุบัน
    viper.AddConfigPath("/etc/myapp/") // หาในโฟลเดอร์อื่นด้วย

    // อ่านไฟล์ Config
    if err := viper.ReadInConfig(); err != nil {
        fmt.Println("ไม่พบไฟล์ config:", err)
    }

    // อ่านค่าจาก Environment Variable (ทับค่าในไฟล์ได้)
    viper.AutomaticEnv()
    viper.SetEnvPrefix("MYAPP") // MYAPP_PORT

    // ดึงค่า
    port := viper.GetInt("server.port")
    host := viper.GetString("server.host")
    debug := viper.GetBool("debug")

    fmt.Printf("Server: %s:%d (debug=%v)\n", host, port, debug)

    // ตั้งค่า Default
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("server.host", "localhost")

    // อ่านค่าทั้งหมดเป็น Struct
    type Config struct {
        Server struct {
            Host string `mapstructure:"host"`
            Port int    `mapstructure:"port"`
        } `mapstructure:"server"`
        Database struct {
            Host     string `mapstructure:"host"`
            User     string `mapstructure:"user"`
            Password string `mapstructure:"password"`
        } `mapstructure:"database"`
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", config)
}
```

**ไฟล์ config.yaml ตัวอย่าง:**
```yaml
server:
  host: 0.0.0.0
  port: 8080

database:
  host: localhost
  user: admin
  password: secret123
```

---

#### 43.3 cobra - การสร้าง CLI Application

`cobra` เป็นไลบรารีที่ใช้สร้าง Command Line Interface (CLI) ที่ซับซ้อน คล้ายกับ `kubectl`, `docker`, `git` ใช้สำหรับสร้างเครื่องมือที่มีคำสั่งย่อย (subcommands)

**การติดตั้ง:**
```bash
go get github.com/spf13/cobra
```

**การใช้งานพื้นฐาน:**
```go
package main

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "mycli",
    Short: "เครื่องมือ CLI ตัวอย่าง",
    Long:  "นี่คือเครื่องมือ CLI ที่สร้างด้วย cobra",
}

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "แสดงเวอร์ชัน",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("mycli v1.0.0")
    },
}

var greetCmd = &cobra.Command{
    Use:   "greet [name]",
    Short: "ทักทายผู้ใช้",
    Args:  cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        name := args[0]
        fmt.Printf("สวัสดี %s!\n", name)
    },
}

// คำสั่งที่มี flag
var configCmd = &cobra.Command{
    Use:   "config",
    Short: "จัดการ Configuration",
    Run: func(cmd *cobra.Command, args []string) {
        configFile, _ := cmd.Flags().GetString("file")
        fmt.Println("ไฟล์ config:", configFile)
    },
}

func main() {
    // เพิ่ม Flags ให้กับ configCmd
    configCmd.Flags().StringP("file", "f", "config.yaml", "ไฟล์ config ที่ใช้")

    rootCmd.AddCommand(versionCmd)
    rootCmd.AddCommand(greetCmd)
    rootCmd.AddCommand(configCmd)

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
```

**การใช้งานจริง:**
```bash
# รันคำสั่ง
go run main.go version          # mycli v1.0.0
go run main.go greet สมชาย      # สวัสดี สมชาย!
go run main.go config -f my.yaml # ไฟล์ config: my.yaml
```

---

#### 43.4 zap - Logging ระดับสูง

`zap` โดย Uber เป็นไลบรารี Logging ที่เร็วที่สุดตัวหนึ่งใน Go รองรับ Structured Logging และ Log หลายระดับ

**การติดตั้ง:**
```bash
go get go.uber.org/zap
```

**การใช้งานพื้นฐาน:**
```go
package main

import (
    "go.uber.org/zap"
)

func main() {
    // 1. ใช้ Production Config (รวดเร็ว, JSON format)
    logger, _ := zap.NewProduction()
    defer logger.Sync()

    // 2. ใช้ Development Config (อ่านง่าย, สีสัน)
    // logger, _ := zap.NewDevelopment()

    // Log ระดับต่างๆ
    logger.Info("แอปพลิเคชันเริ่มทำงาน",
        zap.String("version", "1.0.0"),
        zap.Int("port", 8080),
    )

    logger.Debug("Debug ข้อมูล", zap.Any("data", map[string]int{"a": 1}))
    logger.Warn("คำเตือน", zap.String("message", "ทรัพยากรใกล้หมด"))
    logger.Error("เกิดข้อผิดพลาด", zap.Error(fmt.Errorf("connection failed")))

    // การสร้าง Logger แบบกำหนดเอง
    customLogger, _ := zap.NewProduction(
        zap.WithCaller(true),      // แสดง caller
        zap.AddStacktrace(zap.ErrorLevel), // แสดง stack trace เมื่อ Error
    )
    customLogger.Info("Custom Logger")

    // สร้าง Sugar Logger (ใช้งานง่ายกว่า)
    sugar := logger.Sugar()
    sugar.Infow("Sugar Logger",
        "name", "สมชาย",
        "age", 30,
    )
    sugar.Infof("ชื่อ %s อายุ %d", "สมศรี", 25)
}
```

**ผลลัพธ์ (Production - JSON):**
```json
{"level":"info","ts":1678901234.567,"caller":"main.go:12","msg":"แอปพลิเคชันเริ่มทำงาน","version":"1.0.0","port":8080}
{"level":"warn","ts":1678901234.568,"caller":"main.go:16","msg":"คำเตือน","message":"ทรัพยากรใกล้หมด"}
```

**การใช้งานกับ HTTP Server:**
```go
var logger *zap.Logger

func main() {
    logger, _ = zap.NewProduction()
    defer logger.Sync()

    http.HandleFunc("/hello", helloHandler)
    http.ListenAndServe(":8080", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    logger.Info("Request received",
        zap.String("method", r.Method),
        zap.String("path", r.URL.Path),
        zap.String("remote", r.RemoteAddr),
    )
    w.Write([]byte("Hello"))
}
```

---

### บทที่ 44: GORM – ORM ทรงพลังสำหรับ Go

GORM (Go Object-Relational Mapping) เป็นไลบรารีที่ช่วยให้การทำงานกับฐานข้อมูลสะดวกขึ้น โดยแปลงข้อมูลในฐานข้อมูลให้เป็น Struct และในทางกลับกัน

**การติดตั้ง:**
```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get -u gorm.io/driver/postgres
go get -u gorm.io/driver/sqlite
```

**การเชื่อมต่อฐานข้อมูล:**
```go
package main

import (
    "gorm.io/driver/mysql"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type User struct {
    gorm.Model        // ID, CreatedAt, UpdatedAt, DeletedAt
    Name     string
    Email    string `gorm:"uniqueIndex"`
    Age      int
    Active   bool
}

func main() {
    // เชื่อมต่อ SQLite
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("เชื่อมต่อฐานข้อมูลล้มเหลว")
    }

    // เชื่อมต่อ MySQL
    // dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    // db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

    // Auto Migrate (สร้างตารางอัตโนมัติ)
    db.AutoMigrate(&User{})
}
```

**CRUD พื้นฐาน:**
```go
// CREATE
user := User{Name: "สมชาย", Email: "somchai@mail.com", Age: 30}
result := db.Create(&user)
fmt.Println("ID:", user.ID, "Error:", result.Error)

// READ
var users []User
db.Find(&users) // ดึงทั้งหมด
db.First(&user, 1) // ดึง ID=1
db.First(&user, "name = ?", "สมชาย")
db.Where("age > ?", 25).Find(&users)

// UPDATE
db.Model(&user).Update("Age", 31)
db.Model(&user).Updates(User{Name: "สมชายใหม่", Age: 32})

// DELETE (Soft Delete - ใช้ DeletedAt)
db.Delete(&user, 1) // จะตั้งค่า DeletedAt แทนการลบจริง
```

**การใช้ GORM กับ Relations:**
```go
// One-to-One
type Profile struct {
    gorm.Model
    UserID uint
    Bio    string
}

type User struct {
    gorm.Model
    Name    string
    Profile Profile // One-to-One
}

// One-to-Many
type Post struct {
    gorm.Model
    Title  string
    UserID uint
    User   User // belongs to
}

type User struct {
    gorm.Model
    Name  string
    Posts []Post // has many
}

// Many-to-Many
type Product struct {
    gorm.Model
    Name    string
    Orders  []Order `gorm:"many2many:order_products;"`
}

// Preloading (Eager Loading)
db.Preload("Profile").First(&user, 1)
db.Preload("Posts").Preload("Posts.Tags").Find(&users)
```

**การใช้ Transactions:**
```go
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&User{Name: "สมชาย"}).Error; err != nil {
        return err // Rollback
    }
    if err := tx.Create(&Profile{UserID: 1, Bio: "Developer"}).Error; err != nil {
        return err // Rollback
    }
    return nil // Commit
})
```

---

### บทที่ 45: การส่งอีเมลด้วย gomail และ hermes

#### 45.1 gomail - การส่งอีเมลที่ง่ายและรวดเร็ว

`gomail` (หรือ `gopkg.in/mail.v2`) เป็นไลบรารีสำหรับส่งอีเมลแบบ SMTP รองรับ HTML, Attachments, และ CC/BCC

**การติดตั้ง:**
```bash
go get gopkg.in/mail.v2
```

**การใช้งานพื้นฐาน:**
```go
package main

import (
    "gopkg.in/mail.v2"
)

func main() {
    // สร้างข้อความ
    m := mail.NewMessage()

    // ตั้งค่า Sender และ Receiver
    m.SetHeader("From", "sender@example.com")
    m.SetHeader("To", "recipient@example.com")
    m.SetHeader("Cc", "cc@example.com")
    m.SetHeader("Subject", "สวัสดีจาก Go!")

    // เนื้อหา (Text + HTML)
    m.SetBody("text/plain", "สวัสดี นี่คือข้อความทดสอบ")
    m.AddAlternative("text/html", "<h1>สวัสดี</h1><p>นี่คือข้อความ <b>HTML</b></p>")

    // แนบไฟล์
    m.Attach("/path/to/file.pdf")

    // ตั้งค่า SMTP Server
    d := mail.NewDialer("smtp.gmail.com", 587, "your-email@gmail.com", "your-password")

    // ส่งอีเมล
    if err := d.DialAndSend(m); err != nil {
        panic(err)
    }
    fmt.Println("ส่งอีเมลสำเร็จ!")
}
```

**การสร้างฟังก์ชันส่งอีเมลแบบ Reusable:**
```go
type EmailConfig struct {
    SMTPHost string
    SMTPPort int
    Username string
    Password string
}

type EmailMessage struct {
    To      []string
    Cc      []string
    Subject string
    Body    string
    HTMLBody string
    Attachments []string
}

func SendEmail(config EmailConfig, msg EmailMessage) error {
    m := mail.NewMessage()
    m.SetHeader("From", config.Username)
    m.SetHeader("To", msg.To...)
    if len(msg.Cc) > 0 {
        m.SetHeader("Cc", msg.Cc...)
    }
    m.SetHeader("Subject", msg.Subject)
    
    if msg.HTMLBody != "" {
        m.SetBody("text/html", msg.HTMLBody)
    } else {
        m.SetBody("text/plain", msg.Body)
    }

    for _, file := range msg.Attachments {
        m.Attach(file)
    }

    d := mail.NewDialer(config.SMTPHost, config.SMTPPort, config.Username, config.Password)
    return d.DialAndSend(m)
}
```

---

#### 45.2 hermes - การสร้างอีเมล HTML ที่สวยงาม

`hermes` เป็นไลบรารีสำหรับสร้างอีเมลที่มีรูปแบบสวยงาม รองรับ Responsive Design และเทมเพลตต่างๆ

**การติดตั้ง:**
```bash
go get github.com/matcornic/hermes/v2
```

**การใช้งาน:**
```go
package main

import (
    "fmt"
    "github.com/matcornic/hermes/v2"
)

func main() {
    h := hermes.Hermes{
        Product: hermes.Product{
            Name:      "My App",
            Link:      "https://myapp.com",
            Logo:      "https://myapp.com/logo.png",
            Copyright: "© 2024 My App",
        },
    }

    // สร้างอีเมลต้อนรับ
    email := hermes.Email{
        Body: hermes.Body{
            Name: "สมชาย",
            Intros: []string{
                "ยินดีต้อนรับสู่ My App!",
                "เราดีใจที่คุณเข้าร่วมกับเรา",
            },
            Actions: []hermes.Action{
                {
                    Instructions: "คลิกปุ่มด้านล่างเพื่อยืนยันอีเมล",
                    Button: hermes.Button{
                        Color: "#22BC66",
                        Text:  "ยืนยันอีเมล",
                        Link:  "https://myapp.com/verify?token=abc123",
                    },
                },
            },
            Outros: []string{
                "หากคุณมีข้อสงสัย ติดต่อเราได้ที่ support@myapp.com",
            },
        },
    }

    // สร้าง HTML
    htmlBody, err := h.GenerateHTML(email)
    if err != nil {
        panic(err)
    }

    // สร้าง Text (Plain)
    textBody, err := h.GeneratePlainText(email)
    if err != nil {
        panic(err)
    }

    fmt.Println("HTML:\n", htmlBody)
    fmt.Println("Text:\n", textBody)
}
```

**การใช้งานร่วมกับ gomail:**
```go
func SendWelcomeEmail(to string, name string, token string) error {
    h := hermes.Hermes{
        Product: hermes.Product{
            Name: "My App",
            Link: "https://myapp.com",
        },
    }

    email := hermes.Email{
        Body: hermes.Body{
            Name: name,
            Intros: []string{
                fmt.Sprintf("สวัสดี %s!", name),
                "ยินดีต้อนรับสู่ My App",
            },
            Actions: []hermes.Action{
                {
                    Instructions: "คลิกเพื่อยืนยันอีเมล",
                    Button: hermes.Button{
                        Color: "#22BC66",
                        Text:  "ยืนยัน",
                        Link:  fmt.Sprintf("https://myapp.com/verify?token=%s", token),
                    },
                },
            },
        },
    }

    html, _ := h.GenerateHTML(email)
    text, _ := h.GeneratePlainText(email)

    msg := EmailMessage{
        To:       []string{to},
        Subject:  "ยืนยันอีเมลของคุณ",
        Body:     text,
        HTMLBody: html,
    }

    config := EmailConfig{
        SMTPHost: "smtp.gmail.com",
        SMTPPort: 587,
        Username: "noreply@myapp.com",
        Password: "password",
    }

    return SendEmail(config, msg)
}
```

---

### 📌 สรุปเล่มที่ 7
ในเล่มนี้คุณได้เรียนรู้เครื่องมือและไลบรารียอดนิยมที่ใช้ในโลกจริง:
✅ **chi** - HTTP Router ที่รวดเร็วและยืดหยุ่น
✅ **viper** - การจัดการ Configuration จากหลายแหล่ง (ไฟล์, ENV, Flags)
✅ **cobra** - การสร้าง CLI Application ที่มี Subcommands
✅ **zap** - Structured Logging ที่มีความเร็วสูง
✅ **GORM** - ORM ที่ช่วยจัดการฐานข้อมูลแบบง่ายและทรงพลัง
✅ **gomail** - การส่งอีเมลผ่าน SMTP
✅ **hermes** - การสร้างอีเมล HTML ที่สวยงามและ Responsive
