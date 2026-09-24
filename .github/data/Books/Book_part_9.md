## 📚 เล่มที่ 9: การผสานระบบภายนอกและคุณลักษณะเสริม

เล่มนี้จะสอนคุณถึงการเชื่อมต่อและใช้งานระบบภายนอกที่จำเป็นสำหรับแอปพลิเคชันสมัยใหม่ ไม่ว่าจะเป็นการทำ Cache ด้วย Redis, การส่งข้อความผ่าน RabbitMQ และ MQTT, การเก็บข้อมูลแบบ Time-Series ด้วย InfluxDB, การสื่อสารแบบ Real-time ด้วย WebSocket, และการแจ้งเตือนผ่านช่องทางต่างๆ

---

### บทที่ 52: Redis สำหรับ Cache และ Message Queue

Redis เป็น In-memory Database ที่รวดเร็ว นิยมใช้ทำ Cache, Session Management, Message Queue, และ Rate Limiting

**การติดตั้ง:**
```bash
go get github.com/go-redis/redis/v8
```

**การเชื่อมต่อและใช้งานพื้นฐาน:**
```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
    // เชื่อมต่อ Redis
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // ไม่มีรหัสผ่าน
        DB:       0,  // ใช้ database 0
    })

    // ทดสอบการเชื่อมต่อ
    pong, err := rdb.Ping(ctx).Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("เชื่อมต่อ Redis สำเร็จ:", pong)

    // 1. SET / GET
    err = rdb.Set(ctx, "name", "สมชาย", 10*time.Minute).Err()
    if err != nil {
        panic(err)
    }
    val, err := rdb.Get(ctx, "name").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("name:", val)

    // 2. ทำงานกับ Struct (ใช้ JSON)
    type User struct {
        ID   int
        Name string
        Age  int
    }
    user := User{ID: 1, Name: "สมศรี", Age: 25}
    rdb.Set(ctx, "user:1", user, 0)

    var result User
    err = rdb.Get(ctx, "user:1").Scan(&result)
    if err != nil {
        panic(err)
    }
    fmt.Printf("User: %+v\n", result)

    // 3. ตรวจสอบว่ามีคีย์หรือไม่
    exists, _ := rdb.Exists(ctx, "name").Result()
    fmt.Println("มีคีย์ name?", exists == 1)

    // 4. ตั้งค่าเวลาให้หมดอายุ
    rdb.Expire(ctx, "name", 5*time.Second)

    // 5. ลบข้อมูล
    rdb.Del(ctx, "name")
}
```

**การใช้ Redis เป็น Cache สำหรับ HTTP:**
```go
type CacheMiddleware struct {
    rdb *redis.Client
    ttl time.Duration
}

func NewCacheMiddleware(rdb *redis.Client, ttl time.Duration) *CacheMiddleware {
    return &CacheMiddleware{rdb: rdb, ttl: ttl}
}

func (m *CacheMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        key := r.URL.Path + r.URL.RawQuery
        cached, err := m.rdb.Get(ctx, key).Result()
        if err == nil {
            // พบใน Cache
            w.Header().Set("X-Cache", "HIT")
            w.Write([]byte(cached))
            return
        }

        // ไม่พบใน Cache เรียก Handler จริง
        rw := &responseWriter{w, nil}
        next.ServeHTTP(rw, r)

        // เก็บผลลัพธ์ใน Cache
        if rw.statusCode == http.StatusOK {
            m.rdb.Set(ctx, key, string(rw.body), m.ttl)
        }
    })
}

type responseWriter struct {
    http.ResponseWriter
    body       []byte
    statusCode int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    rw.body = append(rw.body, b...)
    return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
    rw.statusCode = statusCode
    rw.ResponseWriter.WriteHeader(statusCode)
}
```

**การใช้ Redis เป็น Message Queue (Pub/Sub):**
```go
// Publisher
func publishMessage(rdb *redis.Client, channel, message string) error {
    return rdb.Publish(ctx, channel, message).Err()
}

// Subscriber
func subscribeMessages(rdb *redis.Client, channel string) {
    pubsub := rdb.Subscribe(ctx, channel)
    defer pubsub.Close()

    ch := pubsub.Channel()
    for msg := range ch {
        fmt.Printf("ได้รับข้อความจาก %s: %s\n", msg.Channel, msg.Payload)
    }
}

func main() {
    rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

    // เริ่ม Subscriber
    go subscribeMessages(rdb, "notifications")

    // ส่งข้อความ
    publishMessage(rdb, "notifications", "Hello from Go!")
    time.Sleep(2 * time.Second)
}
```

**การใช้ Redis สำหรับ Rate Limiting:**
```go
func rateLimit(rdb *redis.Client, key string, limit int, window time.Duration) bool {
    // ใช้ Redis INCR และ EXPIRE
    count, err := rdb.Incr(ctx, key).Result()
    if err != nil {
        return false
    }
    if count == 1 {
        rdb.Expire(ctx, key, window)
    }
    return count <= int64(limit)
}

// ใช้ใน HTTP Middleware
func RateLimitMiddleware(rdb *redis.Client, limit int, window time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip := r.RemoteAddr
            key := fmt.Sprintf("ratelimit:%s", ip)
            if !rateLimit(rdb, key, limit, window) {
                http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

---

### บทที่ 53: RabbitMQ – Message Broker มาตรฐานองค์กร

RabbitMQ เป็น Message Broker ที่ใช้กันอย่างแพร่หลายในระบบองค์กร รองรับรูปแบบการส่งข้อความที่หลากหลาย

**การติดตั้ง:**
```bash
go get github.com/streadway/amqp
```

**การเชื่อมต่อและใช้งานพื้นฐาน:**
```go
package main

import (
    "fmt"
    "log"
    "github.com/streadway/amqp"
)

func main() {
    // เชื่อมต่อ RabbitMQ
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    // สร้าง Channel
    ch, err := conn.Channel()
    if err != nil {
        panic(err)
    }
    defer ch.Close()

    // ประกาศ Queue
    q, err := ch.QueueDeclare(
        "task_queue", // ชื่อ Queue
        true,         // durable (เก็บไว้แม้ server restart)
        false,        // auto-delete
        false,        // exclusive
        false,        // no-wait
        nil,          // arguments
    )
    if err != nil {
        panic(err)
    }

    // ส่งข้อความ (Producer)
    body := "Hello RabbitMQ!"
    err = ch.Publish(
        "",          // exchange
        q.Name,      // routing key
        false,       // mandatory
        false,       // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(body),
            DeliveryMode: amqp.Persistent, // เก็บใน disk
        })
    if err != nil {
        panic(err)
    }
    fmt.Println("ส่งข้อความ:", body)

    // รับข้อความ (Consumer)
    msgs, err := ch.Consume(
        q.Name, // queue
        "",     // consumer
        true,   // auto-ack (ตอบรับอัตโนมัติ)
        false,  // exclusive
        false,  // no-local
        false,  // no-wait
        nil,    // args
    )
    if err != nil {
        panic(err)
    }

    go func() {
        for d := range msgs {
            log.Printf("ได้รับข้อความ: %s", d.Body)
        }
    }()

    // รอ
    fmt.Scanln()
}
```

**การใช้งาน RabbitMQ ในโปรเจกต์จริง (Producer/Consumer แยกกัน):**

*Producer (sender.go):*
```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/streadway/amqp"
)

type Order struct {
    ID     string `json:"id"`
    UserID string `json:"user_id"`
    Total  int    `json:"total"`
}

func sendOrder(order Order) error {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        return err
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        return err
    }
    defer ch.Close()

    // ประกาศ Exchange (Fanout = ส่งให้ทุก Queue)
    err = ch.ExchangeDeclare(
        "orders",   // name
        "fanout",   // type
        true,       // durable
        false,      // auto-deleted
        false,      // internal
        false,      // no-wait
        nil,        // arguments
    )
    if err != nil {
        return err
    }

    body, _ := json.Marshal(order)
    return ch.Publish(
        "orders", // exchange
        "",       // routing key
        false,    // mandatory
        false,    // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
}

func main() {
    order := Order{ID: "ORD-001", UserID: "USER-123", Total: 5000}
    if err := sendOrder(order); err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("ส่ง order สำเร็จ")
    }
}
```

*Consumer (receiver.go):*
```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "github.com/streadway/amqp"
)

func main() {
    conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
    defer conn.Close()
    ch, _ := conn.Channel()
    defer ch.Close()

    // ประกาศ Exchange เดียวกัน
    ch.ExchangeDeclare("orders", "fanout", true, false, false, false, nil)

    // สร้าง Queue ชั่วคราว (เฉพาะ consumer นี้)
    q, _ := ch.QueueDeclare(
        "",    // ชื่อว่าง = สร้างชื่ออัตโนมัติ
        false, // durable
        false, // auto-delete
        true,  // exclusive
        false, // no-wait
        nil,
    )

    // ผูก Queue กับ Exchange
    ch.QueueBind(q.Name, "", "orders", false, nil)

    msgs, _ := ch.Consume(q.Name, "", true, false, false, false, nil)

    for d := range msgs {
        var order Order
        json.Unmarshal(d.Body, &order)
        log.Printf("ได้รับ Order: %+v", order)
        // ประมวลผล Order (ส่งอีเมล, อัปเดต DB ฯลฯ)
    }
}
```

---

### บทที่ 54: MQTT สำหรับ IoT และระบบเรียลไทม์

MQTT เป็นโปรโตคอลน้ำหนักเบาสำหรับ IoT (Internet of Things) เหมาะสำหรับการส่งข้อมูลระหว่างอุปกรณ์ที่มีทรัพยากรจำกัด

**การติดตั้ง:**
```bash
go get github.com/eclipse/paho.mqtt.golang
```

**การใช้งานพื้นฐาน:**
```go
package main

import (
    "fmt"
    "time"
    mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
    fmt.Printf("ได้รับข้อความ: %s จาก topic: %s\n", msg.Payload(), msg.Topic())
}

func main() {
    // ตั้งค่า MQTT Client
    opts := mqtt.NewClientOptions()
    opts.AddBroker("tcp://broker.emqx.io:1883")
    opts.SetClientID("go_mqtt_client")
    opts.SetUsername("")
    opts.SetPassword("")

    // ตั้งค่า Callback
    opts.SetDefaultPublishHandler(messagePubHandler)

    // เชื่อมต่อ
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }

    // Subscribe topic
    topic := "sensor/temperature"
    token := client.Subscribe(topic, 1, nil)
    token.Wait()
    fmt.Printf("Subscribed to topic: %s\n", topic)

    // Publish ข้อความ
    go func() {
        for i := 0; i < 10; i++ {
            payload := fmt.Sprintf(`{"temperature": %.2f, "timestamp": %d}`, 25.5+float64(i), time.Now().Unix())
            token := client.Publish(topic, 1, false, payload)
            token.Wait()
            fmt.Printf("ส่งข้อความ: %s\n", payload)
            time.Sleep(2 * time.Second)
        }
        client.Disconnect(250)
    }()

    // รอ
    time.Sleep(30 * time.Second)
}
```

**การใช้งานในระบบ IoT (Sensor Data Collector):**
```go
type SensorData struct {
    DeviceID    string  `json:"device_id"`
    Temperature float64 `json:"temperature"`
    Humidity    float64 `json:"humidity"`
    Timestamp   int64   `json:"timestamp"`
}

func handleSensorData(client mqtt.Client, msg mqtt.Message) {
    var data SensorData
    if err := json.Unmarshal(msg.Payload(), &data); err != nil {
        log.Println("Parse error:", err)
        return
    }
    log.Printf("Device: %s, Temp: %.2f, Humidity: %.2f", data.DeviceID, data.Temperature, data.Humidity)
    // บันทึกไปฐานข้อมูล
}
```

---

### บทที่ 55: InfluxDB – Time‑Series Database

InfluxDB เป็นฐานข้อมูลที่ออกแบบมาเฉพาะสำหรับข้อมูลที่มีเวลา (Time-Series Data) เช่น Metrics, Logs, Sensor Data

**การติดตั้ง:**
```bash
go get github.com/influxdata/influxdb-client-go/v2
```

**การใช้งานพื้นฐาน:**
```go
package main

import (
    "context"
    "fmt"
    "time"
    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
    // สร้าง Client
    client := influxdb2.NewClient("http://localhost:8086", "my-token")
    defer client.Close()

    // สร้าง Write API
    writeAPI := client.WriteAPI("my-org", "my-bucket")

    // สร้าง Point (ข้อมูล)
    point := influxdb2.NewPoint(
        "temperature", // measurement
        map[string]string{"sensor": "sensor-001", "room": "living"}, // tags
        map[string]interface{}{"value": 25.5, "humidity": 60.0}, // fields
        time.Now(),
    )

    // เขียนข้อมูล
    writeAPI.WritePoint(point)
    writeAPI.Flush()

    // Query ข้อมูล
    queryAPI := client.QueryAPI("my-org")
    query := `from(bucket:"my-bucket")
              |> range(start: -1h)
              |> filter(fn: (r) => r._measurement == "temperature")
              |> filter(fn: (r) => r.sensor == "sensor-001")`

    result, err := queryAPI.Query(context.Background(), query)
    if err != nil {
        panic(err)
    }

    for result.Next() {
        record := result.Record()
        fmt.Printf("Time: %s, Value: %v\n", record.Time(), record.Value())
    }
}
```

**การใช้งานกับ HTTP Server (Metric Collection):**
```go
type MetricsCollector struct {
    client influxdb2.Client
    org    string
    bucket string
}

func NewMetricsCollector(url, token, org, bucket string) *MetricsCollector {
    return &MetricsCollector{
        client: influxdb2.NewClient(url, token),
        org:    org,
        bucket: bucket,
    }
}

func (m *MetricsCollector) RecordRequest(path string, method string, duration time.Duration, status int) {
    writeAPI := m.client.WriteAPI(m.org, m.bucket)
    point := influxdb2.NewPoint(
        "http_requests",
        map[string]string{
            "path":   path,
            "method": method,
            "status": fmt.Sprintf("%d", status),
        },
        map[string]interface{}{
            "duration_ms": duration.Milliseconds(),
            "count":       1,
        },
        time.Now(),
    )
    writeAPI.WritePoint(point)
    writeAPI.Flush()
}

// Middleware สำหรับบันทึก Metrics
func (m *MetricsCollector) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rw := &responseWriter{ResponseWriter: w}
        next.ServeHTTP(rw, r)
        m.RecordRequest(r.URL.Path, r.Method, time.Since(start), rw.statusCode)
    })
}
```

---

### บทที่ 56: WebSocket และ Socket.IO

WebSocket ช่วยให้การสื่อสารแบบ Real-time ระหว่าง Client และ Server เป็นไปได้อย่างมีประสิทธิภาพ

**การติดตั้ง WebSocket:**
```bash
go get github.com/gorilla/websocket
```

**WebSocket Server:**
```go
package main

import (
    "fmt"
    "net/http"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true }, // อนุญาตทุก Origin (เพื่อความง่าย)
}

// จัดการ Connection
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, "Cannot upgrade", http.StatusInternalServerError)
        return
    }
    defer conn.Close()

    fmt.Println("Client connected")

    // รับและส่งข้อความ
    for {
        // รับข้อความจาก Client
        messageType, msg, err := conn.ReadMessage()
        if err != nil {
            fmt.Println("Client disconnected")
            break
        }

        fmt.Printf("Received: %s\n", msg)

        // ส่งข้อความกลับ (Echo)
        err = conn.WriteMessage(messageType, []byte("Received: "+string(msg)))
        if err != nil {
            break
        }
    }
}

func main() {
    http.HandleFunc("/ws", handleWebSocket)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "index.html")
    })
    http.ListenAndServe(":8080", nil)
}
```

**WebSocket Chat Room แบบง่าย:**
```go
type ChatRoom struct {
    clients map[*websocket.Conn]bool
    join    chan *websocket.Conn
    leave   chan *websocket.Conn
    message chan []byte
}

func NewChatRoom() *ChatRoom {
    return &ChatRoom{
        clients: make(map[*websocket.Conn]bool),
        join:    make(chan *websocket.Conn),
        leave:   make(chan *websocket.Conn),
        message: make(chan []byte),
    }
}

func (c *ChatRoom) Run() {
    for {
        select {
        case conn := <-c.join:
            c.clients[conn] = true
            fmt.Println("New client joined")
        case conn := <-c.leave:
            delete(c.clients, conn)
            conn.Close()
            fmt.Println("Client left")
        case msg := <-c.message:
            for conn := range c.clients {
                conn.WriteMessage(websocket.TextMessage, msg)
            }
        }
    }
}

// ใช้ใน handler
func (c *ChatRoom) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    c.join <- conn
    defer func() { c.leave <- conn }()

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            break
        }
        c.message <- msg
    }
}
```

---

### บทที่ 57: การส่ง SMS และ LINE Notify

#### 57.1 การส่ง SMS (ตัวอย่าง Twilio)

**การติดตั้ง:**
```bash
go get github.com/twilio/twilio-go
```

**การใช้งาน:**
```go
package main

import (
    "fmt"
    "github.com/twilio/twilio-go"
    twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func sendSMS(to, from, body string) error {
    client := twilio.NewRestClientWithParams(twilio.ClientParams{
        Username: "ACxxxxxxxxxxxx", // Account SID
        Password: "xxxxxxxxxxxx",   // Auth Token
    })

    params := &twilioApi.CreateMessageParams{}
    params.SetTo(to)
    params.SetFrom(from)
    params.SetBody(body)

    resp, err := client.ApiV2010.CreateMessage(params)
    if err != nil {
        return err
    }
    fmt.Printf("SMS sent! SID: %s\n", *resp.Sid)
    return nil
}

func main() {
    sendSMS("+668xxxxxxx", "+1xxxxxxx", "สวัสดีจาก Go!")
}
```

#### 57.2 LINE Notify

LINE Notify เป็นบริการแจ้งเตือนฟรีผ่าน LINE

**การส่งข้อความ LINE Notify:**
```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

type LineNotify struct {
    Token string
}

func NewLineNotify(token string) *LineNotify {
    return &LineNotify{Token: token}
}

func (l *LineNotify) SendMessage(message string) error {
    url := "https://notify-api.line.me/api/notify"
    data := bytes.NewBufferString(fmt.Sprintf("message=%s", message))

    req, err := http.NewRequest("POST", url, data)
    if err != nil {
        return err
    }
    req.Header.Set("Authorization", "Bearer "+l.Token)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 200 {
        return fmt.Errorf("LINE Notify error: %s", string(body))
    }
    return nil
}

func main() {
    line := NewLineNotify("your-line-notify-token")
    line.SendMessage("🚀 แจ้งเตือน: Server เริ่มทำงานแล้ว!")
}
```

**ส่งข้อความพร้อมรูปภาพ:**
```go
func (l *LineNotify) SendMessageWithImage(message, imageURL string) error {
    url := "https://notify-api.line.me/api/notify"
    data := bytes.NewBufferString(fmt.Sprintf("message=%s&imageThumbnail=%s&imageFullsize=%s", message, imageURL, imageURL))
    // ... (เหมือนเดิม)
}
```

---

### บทที่ 58: Discord Webhook สำหรับแจ้งเตือน

Discord Webhook เป็นวิธีที่ง่ายในการส่งข้อความแจ้งเตือนไปยังช่อง Discord

**การใช้งาน:**
```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type DiscordWebhook struct {
    URL string
}

type DiscordMessage struct {
    Content   string         `json:"content,omitempty"`
    Username  string         `json:"username,omitempty"`
    AvatarURL string         `json:"avatar_url,omitempty"`
    Embeds    []DiscordEmbed `json:"embeds,omitempty"`
}

type DiscordEmbed struct {
    Title       string `json:"title,omitempty"`
    Description string `json:"description,omitempty"`
    Color       int    `json:"color,omitempty"` // สีแบบ Decimal
    URL         string `json:"url,omitempty"`
    Author      struct {
        Name    string `json:"name,omitempty"`
        URL     string `json:"url,omitempty"`
        IconURL string `json:"icon_url,omitempty"`
    } `json:"author,omitempty"`
    Fields []struct {
        Name   string `json:"name"`
        Value  string `json:"value"`
        Inline bool   `json:"inline,omitempty"`
    } `json:"fields,omitempty"`
    Timestamp string `json:"timestamp,omitempty"`
}

func NewDiscordWebhook(url string) *DiscordWebhook {
    return &DiscordWebhook{URL: url}
}

func (d *DiscordWebhook) SendMessage(msg DiscordMessage) error {
    data, err := json.Marshal(msg)
    if err != nil {
        return err
    }

    resp, err := http.Post(d.URL, "application/json", bytes.NewBuffer(data))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 204 {
        return fmt.Errorf("Discord error: status %d", resp.StatusCode)
    }
    return nil
}

func main() {
    webhook := NewDiscordWebhook("https://discord.com/api/webhooks/xxx/yyy")

    // ส่งข้อความธรรมดา
    webhook.SendMessage(DiscordMessage{
        Content: "🚨 แจ้งเตือน: มีผู้ใช้เข้าสู่ระบบ 100 คน",
        Username: "Alert Bot",
    })

    // ส่งข้อความพร้อม Embed
    embed := DiscordEmbed{
        Title:       "📊 รายงานรายวัน",
        Description: "สรุปข้อมูลของวันนี้",
        Color:       0x00ff00, // สีเขียว
        Fields: []struct {
            Name   string `json:"name"`
            Value  string `json:"value"`
            Inline bool   `json:"inline,omitempty"`
        }{
            {Name: "ผู้ใช้ใหม่", Value: "25 คน", Inline: true},
            {Name: "ยอดขาย", Value: "50,000 บาท", Inline: true},
            {Name: "ข้อผิดพลาด", Value: "0 ครั้ง", Inline: true},
        },
    }

    webhook.SendMessage(DiscordMessage{
        Embeds: []DiscordEmbed{embed},
    })
}
```

**การใช้งานในระบบ Monitoring:**
```go
type Alert struct {
    webhook *DiscordWebhook
}

func (a *Alert) SendError(err error, context map[string]interface{}) {
    embed := DiscordEmbed{
        Title:       "❌ เกิดข้อผิดพลาด",
        Description: err.Error(),
        Color:       0xff0000, // สีแดง
        Timestamp:   time.Now().Format(time.RFC3339),
    }

    // เพิ่ม context เป็น Fields
    for k, v := range context {
        embed.Fields = append(embed.Fields, struct {
            Name   string `json:"name"`
            Value  string `json:"value"`
            Inline bool   `json:"inline,omitempty"`
        }{
            Name:  k,
            Value: fmt.Sprintf("%v", v),
            Inline: true,
        })
    }

    a.webhook.SendMessage(DiscordMessage{
        Embeds: []DiscordEmbed{embed},
    })
}

// ใช้ใน HTTP Handler
func handler(w http.ResponseWriter, r *http.Request) {
    // ...
    if err := processRequest(r); err != nil {
        alert.SendError(err, map[string]interface{}{
            "path": r.URL.Path,
            "method": r.Method,
            "user_id": "user-123",
        })
        http.Error(w, "Internal Error", http.StatusInternalServerError)
        return
    }
}
```

---

### 📌 สรุปเล่มที่ 9
ในเล่มนี้คุณได้เรียนรู้การเชื่อมต่อกับระบบภายนอกต่างๆ:
✅ **Redis** - Cache, Message Queue, Rate Limiting
✅ **RabbitMQ** - Message Broker สำหรับระบบองค์กร
✅ **MQTT** - โปรโตคอลสำหรับ IoT และระบบเรียลไทม์
✅ **InfluxDB** - Time-Series Database สำหรับ Metrics
✅ **WebSocket** - การสื่อสารแบบ Real-time
✅ **SMS** และ **LINE Notify** - การแจ้งเตือนผ่านมือถือ
✅ **Discord Webhook** - การแจ้งเตือนผ่าน Discord

--- 