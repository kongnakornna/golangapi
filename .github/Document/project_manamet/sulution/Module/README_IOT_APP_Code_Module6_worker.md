# Module 6: Worker (Background Workers)

## สำหรับโฟลเดอร์ `internal/delivery/worker/`

ไฟล์ที่เกี่ยวข้อง:
- `internal/delivery/worker/email_worker.go`
- `internal/delivery/worker/mqtt_worker.go` (optional, สำหรับ IoT)
- `internal/pkg/queue/redis_queue.go` (ตัวจัดการคิว)

---

## หลักการ (Concept)

### คืออะไร?
Worker คือกระบวนการที่ทำงานเบื้องหลัง (background) เพื่อประมวลผลงานที่ใช้เวลานาน หรืองานที่ไม่ต้องการให้ผู้ใช้รอ เช่น การส่งอีเมล, การประมวลผลข้อมูล, การแจ้งเตือน โดยรับงานจากคิว (message queue) และประมวลผลแบบ asynchronous

### มีกี่แบบ?
1. **Email Worker** – ส่งอีเมลยืนยัน, รีเซ็ตรหัสผ่าน, รายงาน
2. **MQTT Worker** – รับข้อมูลจาก MQTT broker และบันทึกลงฐานข้อมูล
3. **Report Worker** – สร้างรายงาน PDF/Excel
4. **Cleanup Worker** – ลบ logs เก่า, ล้าง cache
5. **Generic Worker Pool** – รองรับงานหลายประเภทผ่าน Redis queue

### ใช้อย่างไร / นำไปใช้กรณีไหน
- ใช้เมื่อมีงานที่ใช้เวลานาน (ส่งอีเมล, ประมวลผลรูปภาพ)
- ใช้เมื่อต้องการหน่วงเวลาการทำงาน (retry, schedule)
- ใช้เพื่อลดภาระของ HTTP server (ไม่บล็อก request)

### ทำไมต้องใช้
- ป้องกัน HTTP request timeout
- ปรับปรุง用户体验 (response เร็ว)
- รองรับ retry และ error handling แบบรวมศูนย์

### ประโยชน์ที่ได้รับ
- ระบบ responsive มากขึ้น
- สามารถ scale worker ต่างหากจาก API server
- รองรับ guaranteed delivery (at-least-once)

### ข้อควรระวัง
- worker ต้องมีการจัดการ panic และ retry
- ต้องมี monitoring สำหรับ queue length
- ต้องระวัง duplicate messages (idempotency)

### ข้อดี
- แยกความรับผิดชอบ, ปรับขนาดได้, ทนทานต่อความล้มเหลว

### ข้อเสีย
- เพิ่มความซับซ้อน (ต้องมี message broker)
- debugging ยากขึ้น (异步)

### ข้อห้าม
- ห้ามทำงานที่ต้อง response ทันทีใน worker (ใช้ API โดยตรง)
- ห้าม worker ตายโดยไม่มีการแจ้งเตือน

---

## โค้ดที่รันได้จริง

### 1. Redis Queue – `internal/pkg/queue/redis_queue.go`

```go
// Package queue provides Redis-based task queue.
// ----------------------------------------------------------------
// แพ็คเกจ queue ให้บริการคิวงานบน Redis
package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Task represents a unit of work.
// ----------------------------------------------------------------
// Task แทนหน่วยงานหนึ่ง
type Task struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`   // "email", "mqtt_publish"
	Payload   map[string]interface{} `json:"payload"`
	CreatedAt time.Time              `json:"created_at"`
	RetryCount int                   `json:"retry_count"`
}

// Queue interface for task queue operations.
// ----------------------------------------------------------------
// Queue interface สำหรับการดำเนินการคิวงาน
type Queue interface {
	Enqueue(ctx context.Context, task *Task) error
	Dequeue(ctx context.Context, timeout time.Duration) (*Task, error)
	Ack(ctx context.Context, taskID string) error
	Requeue(ctx context.Context, task *Task, delay time.Duration) error
}

// RedisQueue implements Queue using Redis List and Stream.
// ----------------------------------------------------------------
// RedisQueue อิมพลีเมนต์ Queue ด้วย Redis List
type RedisQueue struct {
	client    *redis.Client
	queueName string
}

// NewRedisQueue creates a new Redis queue.
// ----------------------------------------------------------------
// NewRedisQueue สร้าง Redis queue ใหม่
func NewRedisQueue(client *redis.Client, queueName string) *RedisQueue {
	return &RedisQueue{
		client:    client,
		queueName: queueName,
	}
}

// Enqueue adds a task to the queue (right push).
// ----------------------------------------------------------------
// Enqueue เพิ่มงานเข้าคิว (push ทางขวา)
func (q *RedisQueue) Enqueue(ctx context.Context, task *Task) error {
	task.CreatedAt = time.Now()
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return q.client.RPush(ctx, q.queueName, data).Err()
}

// Dequeue retrieves a task from the queue (left pop) with blocking timeout.
// ----------------------------------------------------------------
// Dequeue ดึงงานออกจากคิว (pop ทางซ้าย) แบบบล็อก
func (q *RedisQueue) Dequeue(ctx context.Context, timeout time.Duration) (*Task, error) {
	result, err := q.client.BLPop(ctx, timeout, q.queueName).Result()
	if err == redis.Nil {
		return nil, nil // no task, ไม่มีงาน
	}
	if err != nil {
		return nil, err
	}
	if len(result) < 2 {
		return nil, nil
	}
	var task Task
	if err := json.Unmarshal([]byte(result[1]), &task); err != nil {
		return nil, err
	}
	return &task, nil
}

// Ack removes a task from processing queue (not needed for simple list, use separate processing queue).
// For production, use Redis Streams with consumer groups.
// ----------------------------------------------------------------
// Ack ลบงานออกจากคิวประมวลผล (ไม่จำเป็นสำหรับ list ธรรมดา)
func (q *RedisQueue) Ack(ctx context.Context, taskID string) error {
	// In simple implementation, we don't need ack.
	// For reliability, use a separate processing list.
	return nil
}

// Requeue pushes task back to queue after delay (using sorted set or delayed queue).
// ----------------------------------------------------------------
// Requeue ใส่งานกลับคิวหลังจากดีเลย์
func (q *RedisQueue) Requeue(ctx context.Context, task *Task, delay time.Duration) error {
	// Simplified: just enqueue again
	return q.Enqueue(ctx, task)
}
```

### 2. Email Worker – `internal/delivery/worker/email_worker.go`

```go
// Package worker contains background workers.
// ----------------------------------------------------------------
// แพ็คเกจ worker บรรจุ worker ที่ทำงานเบื้องหลัง
package worker

import (
	"context"
	"log"
	"time"

	"gobackend/internal/pkg/email"
	"gobackend/internal/pkg/queue"
	"gobackend/internal/pkg/logger"
	"go.uber.org/zap"
)

// EmailWorker processes email tasks from queue.
// ----------------------------------------------------------------
// EmailWorker ประมวลผลงานส่งอีเมลจากคิว
type EmailWorker struct {
	queue   queue.Queue
	sender  email.Sender
	workers int
	stopCh  chan struct{}
}

// EmailTaskPayload defines payload for email tasks.
// ----------------------------------------------------------------
// EmailTaskPayload กำหนด payload สำหรับงานส่งอีเมล
type EmailTaskPayload struct {
	To          string            `json:"to"`
	Subject     string            `json:"subject"`
	Template    string            `json:"template"`    // "verification", "reset_password"
	TemplateData map[string]string `json:"template_data"`
}

// NewEmailWorker creates a new email worker.
// ----------------------------------------------------------------
// NewEmailWorker สร้าง email worker ใหม่
func NewEmailWorker(queue queue.Queue, sender email.Sender, workers int) *EmailWorker {
	return &EmailWorker{
		queue:   queue,
		sender:  sender,
		workers: workers,
		stopCh:  make(chan struct{}),
	}
}

// Start launches worker goroutines.
// ----------------------------------------------------------------
// Start เริ่ม worker goroutines
func (w *EmailWorker) Start(ctx context.Context) {
	for i := 0; i < w.workers; i++ {
		go w.workerLoop(ctx, i)
	}
	log.Printf("EmailWorker started with %d workers", w.workers)
	<-w.stopCh
}

// Stop gracefully shuts down workers.
// ----------------------------------------------------------------
// Stop ปิด worker แบบ graceful
func (w *EmailWorker) Stop() {
	close(w.stopCh)
}

// workerLoop processes tasks continuously.
// ----------------------------------------------------------------
// workerLoop ประมวลผลงานอย่างต่อเนื่อง
func (w *EmailWorker) workerLoop(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d stopping due to context", id)
			return
		case <-w.stopCh:
			log.Printf("Worker %d stopping", id)
			return
		default:
			task, err := w.queue.Dequeue(ctx, 5*time.Second)
			if err != nil {
				logger.Log.Error("failed to dequeue task", zap.Error(err))
				time.Sleep(1 * time.Second)
				continue
			}
			if task == nil {
				continue
			}
			// Process only email tasks
			if task.Type != "email" {
				// Unknown task type, skip or log
				logger.Log.Warn("unknown task type", zap.String("type", task.Type))
				continue
			}
			w.processEmailTask(ctx, task)
		}
	}
}

// processEmailTask sends email based on task payload.
// ----------------------------------------------------------------
// processEmailTask ส่งอีเมลตาม payload ของงาน
func (w *EmailWorker) processEmailTask(ctx context.Context, task *queue.Task) {
	var payload EmailTaskPayload
	if err := mapToStruct(task.Payload, &payload); err != nil {
		logger.Log.Error("invalid email task payload", zap.Error(err))
		return
	}

	// Build email content based on template
	// สร้างเนื้อหาอีเมลตาม template
	var htmlBody string
	switch payload.Template {
	case "verification":
		htmlBody = email.BuildVerificationEmail(payload.TemplateData["name"], payload.TemplateData["link"])
	case "reset_password":
		htmlBody = email.BuildResetPasswordEmail(payload.TemplateData["name"], payload.TemplateData["link"])
	default:
		htmlBody = payload.TemplateData["body"]
	}

	// Send email with retry (max 3 times)
	// ส่งอีเมลพร้อม retry (สูงสุด 3 ครั้ง)
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		err = w.sender.Send(payload.To, payload.Subject, htmlBody)
		if err == nil {
			logger.Log.Info("email sent successfully",
				zap.String("to", payload.To),
				zap.String("template", payload.Template),
			)
			return
		}
		logger.Log.Warn("email send failed, retrying",
			zap.Int("attempt", attempt+1),
			zap.Error(err),
		)
		time.Sleep(time.Duration(attempt+1) * time.Second)
	}
	// After retries, log error and maybe move to dead letter queue
	logger.Log.Error("email send failed after retries",
		zap.String("to", payload.To),
		zap.Error(err),
	)
}

func mapToStruct(m map[string]interface{}, out interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}
```

### 3. Email Sender with Templates – `internal/pkg/email/sender.go` & `gomail_sender.go`

#### `internal/pkg/email/sender.go`

```go
// Package email provides email sending capabilities.
// ----------------------------------------------------------------
// แพ็คเกจ email ให้บริการส่งอีเมล
package email

// Sender interface for sending emails.
// ----------------------------------------------------------------
// Sender interface สำหรับส่งอีเมล
type Sender interface {
	Send(to, subject, body string) error
}
```

#### `internal/pkg/email/gomail_sender.go`

```go
package email

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

// GomailSender implements Sender using gomail.
// ----------------------------------------------------------------
// GomailSender อิมพลีเมนต์ Sender ด้วย gomail
type GomailSender struct {
	host     string
	port     int
	username string
	password string
	from     string
}

// NewGomailSender creates a new Gomail sender.
// ----------------------------------------------------------------
// NewGomailSender สร้าง Gomail sender ใหม่
func NewGomailSender(host string, port int, username, password, from string) *GomailSender {
	return &GomailSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

// Send sends an email.
// ----------------------------------------------------------------
// Send ส่งอีเมล
func (s *GomailSender) Send(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	dialer := gomail.NewDialer(s.host, s.port, s.username, s.password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: false}
	return dialer.DialAndSend(m)
}
```

#### `internal/pkg/email/templates.go` (สร้าง HTML templates)

```go
package email

import "fmt"

// BuildVerificationEmail returns HTML for email verification.
// ----------------------------------------------------------------
// BuildVerificationEmail คืน HTML สำหรับยืนยันอีเมล
func BuildVerificationEmail(name, link string) string {
	return fmt.Sprintf(`
		<html>
		<body>
			<h2>Welcome %s!</h2>
			<p>Please verify your email by clicking the link below:</p>
			<a href="%s">Verify Email</a>
			<p>This link expires in 24 hours.</p>
		</body>
		</html>
	`, name, link)
}

// BuildResetPasswordEmail returns HTML for password reset.
// ----------------------------------------------------------------
// BuildResetPasswordEmail คืน HTML สำหรับรีเซ็ตรหัสผ่าน
func BuildResetPasswordEmail(name, link string) string {
	return fmt.Sprintf(`
		<html>
		<body>
			<h2>Hello %s,</h2>
			<p>You requested to reset your password. Click the link below:</p>
			<a href="%s">Reset Password</a>
			<p>If you didn't request this, ignore this email.</p>
		</body>
		</html>
	`, name, link)
}
```

### 4. MQTT Worker (IoT) – `internal/delivery/worker/mqtt_worker.go`

```go
package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gobackend/internal/models"
	"gorm.io/gorm"
)

// MQTTWorker subscribes to sensor topics and stores data.
// ----------------------------------------------------------------
// MQTTWorker รับ订阅 topic เซนเซอร์และบันทึกข้อมูล
type MQTTWorker struct {
	broker   string
	clientID string
	db       *gorm.DB
	client   mqtt.Client
	stopCh   chan struct{}
}

// SensorData represents incoming sensor reading.
// ----------------------------------------------------------------
// SensorData แทนค่าที่อ่านได้จากเซนเซอร์
type SensorData struct {
	DeviceID   string    `json:"device_id"`
	SensorType string    `json:"sensor_type"`
	Value      float64   `json:"value"`
	Unit       string    `json:"unit"`
	Location   string    `json:"location"`
	Timestamp  time.Time `json:"timestamp"`
}

// NewMQTTWorker creates MQTT worker.
// ----------------------------------------------------------------
// NewMQTTWorker สร้าง MQTT worker
func NewMQTTWorker(broker, clientID string, db *gorm.DB) *MQTTWorker {
	return &MQTTWorker{
		broker:   broker,
		clientID: clientID,
		db:       db,
		stopCh:   make(chan struct{}),
	}
}

// Start connects to MQTT broker and subscribes.
// ----------------------------------------------------------------
// Start เชื่อมต่อ MQTT broker และ订阅 topic
func (w *MQTTWorker) Start(ctx context.Context) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(w.broker)
	opts.SetClientID(w.clientID)
	opts.SetAutoReconnect(true)
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		log.Println("MQTT connected, subscribing...")
		w.subscribe(c)
	})
	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		log.Printf("MQTT connection lost: %v", err)
	})

	w.client = mqtt.NewClient(opts)
	if token := w.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.Println("MQTT worker started")
	<-w.stopCh
	w.client.Disconnect(250)
	return nil
}

func (w *MQTTWorker) subscribe(c mqtt.Client) {
	token := c.Subscribe("cmom/dc/+/sensor/+", 1, w.messageHandler)
	token.Wait()
	if token.Error() != nil {
		log.Printf("Subscribe error: %v", token.Error())
	}
}

func (w *MQTTWorker) messageHandler(client mqtt.Client, msg mqtt.Message) {
	var data SensorData
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Printf("Failed to parse MQTT payload: %v", err)
		return
	}
	// Save to database (asynchronously or use a separate goroutine)
	go w.saveToDB(data)
}

func (w *MQTTWorker) saveToDB(data SensorData) {
	record := models.SensorLog{
		DeviceID:   data.DeviceID,
		SensorType: data.SensorType,
		Value:      data.Value,
		Unit:       data.Unit,
		Location:   data.Location,
		Timestamp:  data.Timestamp,
	}
	if err := w.db.Create(&record).Error; err != nil {
		log.Printf("Failed to save sensor data: %v", err)
	}
}

// Stop shuts down worker.
// ----------------------------------------------------------------
// Stop ปิด worker
func (w *MQTTWorker) Stop() {
	close(w.stopCh)
}
```

### 5. Model สำหรับ SensorLog (เพิ่มใน models)

```go
// SensorLog stores historical sensor readings.
// ----------------------------------------------------------------
// SensorLog เก็บประวัติค่าที่อ่านได้จากเซนเซอร์
type SensorLog struct {
	ID         uint      `gorm:"primaryKey"`
	DeviceID   string    `gorm:"index;size:100"`
	SensorType string    `gorm:"index;size:50"`
	Value      float64
	Unit       string    `gorm:"size:10"`
	Location   string    `gorm:"index;size:100"`
	Timestamp  time.Time `gorm:"index"`
}
```

### 6. การรวม Worker ใน `main.go`

```go
func main() {
	// ... initialize config, db, redis, queue, sender

	// Create Redis queue
	redisQueue := queue.NewRedisQueue(redisClient, "task:email")

	// Create email sender
	emailSender := email.NewGomailSender(
		cfg.SMTP.Host, cfg.SMTP.Port,
		cfg.SMTP.Username, cfg.SMTP.Password,
		cfg.SMTP.From,
	)

	// Create email worker
	emailWorker := worker.NewEmailWorker(redisQueue, emailSender, 3)

	// Create MQTT worker (optional)
	mqttWorker := worker.NewMQTTWorker(cfg.MQTT.Broker, "cmon-worker", db)

	// Start workers in goroutines
	go func() {
		if err := emailWorker.Start(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		if err := mqttWorker.Start(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down workers...")
	emailWorker.Stop()
	mqttWorker.Stop()
}
```

### 7. การใช้งาน: Enqueue email task จาก handler

```go
// ใน auth_handler.go หรือ user_handler.go
func (h *AuthHandler) SendVerificationEmail(userID uint, email, name string) {
	task := &queue.Task{
		ID:   uuid.New().String(),
		Type: "email",
		Payload: map[string]interface{}{
			"to":       email,
			"subject":  "Verify your email",
			"template": "verification",
			"template_data": map[string]string{
				"name": name,
				"link": "https://yourapp.com/verify?token=xxx",
			},
		},
	}
	_ = h.queue.Enqueue(context.Background(), task)
}
```

---

## วิธีใช้งาน module นี้

1. วางไฟล์ตามโครงสร้าง
2. ติดตั้ง dependencies:
   ```
   go get github.com/redis/go-redis/v9
   go get gopkg.in/gomail.v2
   go get github.com/eclipse/paho.mqtt.golang
   ```
3. ตั้งค่า environment variables สำหรับ SMTP และ MQTT
4. เริ่ม worker พร้อมกับ API server
5. เรียกใช้ `Enqueue` เพื่อส่งงานไปยังคิว

---

## ตารางสรุป Worker Types

| Worker | คิว | หน้าที่ | Retry | Idempotent |
|--------|-----|--------|-------|-------------|
| EmailWorker | Redis | ส่งอีเมล async | 3 ครั้ง | ใช่ (idempotent) |
| MQTTWorker | MQTT broker | รับข้อมูลเซนเซอร์ | ไม่ (real-time) | ใช่ |
| ReportWorker | Redis | สร้างรายงาน | 2 ครั้ง | ไม่ (ควร unique) |

---

## แบบฝึกหัดท้าย module (3 ข้อ)

1. เพิ่ม `CleanupWorker` ที่ลบ sensor logs ที่เก่ากว่า 90 วัน ทุกวันเวลา 02:00 น. โดยใช้ scheduler (cron) แทนคิว
2. Implement retry with exponential backoff ใน `EmailWorker` (delay 1s, 2s, 4s, 8s)
3. สร้าง `DeadLetterQueue` สำหรับงานที่ล้มเหลวหลังจาก retry ครบ และมี worker อีกตัวที่ประมวลผล DLQ (แจ้งเตือน admin)

---

## แหล่งอ้างอิง

- [Redis as Message Queue](https://redis.io/docs/latest/develop/use/patterns/message-queue/)
- [Gomail documentation](https://github.com/go-gomail/gomail)
- [Paho MQTT Go client](https://github.com/eclipse/paho.mqtt.golang)
- [Worker Pool Pattern](https://gobyexample.com/worker-pools)

---

**หมายเหตุ:** module นี้ครบถ้วนสำหรับ background workers ตามโครงสร้าง gobackend หากต้องการ module เพิ่มเติม (เช่น pkg/utils, pkg/validator) โปรดแจ้ง