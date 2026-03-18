package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Message ข้อความในคิว
type Message struct {
	ID         string          `json:"id"`
	Topic      string          `json:"topic"`
	Payload    json.RawMessage `json:"payload"`
	Timestamp  time.Time       `json:"timestamp"`
	Retries    int             `json:"retries"`
	MaxRetries int             `json:"max_retries"`
}

// Handler ตัวจัดการข้อความ 
type Handler func(ctx context.Context, msg *Message) error

// Queue อินเทอร์เฟซของคิว
type Queue interface {
	// Publish เผยแพร่ข้อความ
	Publish(ctx context.Context, topic string, payload interface{}) error
	// Subscribe สมัครรับข้อมูลหัวข้อ
	Subscribe(ctx context.Context, topic string, handler Handler) error
	// PublishDelayed เผยแพร่ข้อความแบบหน่วงเวลา
	PublishDelayed(ctx context.Context, topic string, payload interface{}, delay time.Duration) error
	// Close ปิดคิว
	Close() error
}

// RedisQueue การ implement คิวด้วย Redis
type RedisQueue struct {
	client     *redis.Client
	handlers   map[string][]Handler
	mu         sync.RWMutex
	workerPool chan struct{}
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

// NewRedisQueue สร้าง RedisQueue ใหม่
func NewRedisQueue(client *redis.Client, maxWorkers int) Queue {
	ctx, cancel := context.WithCancel(context.Background())

	rq := &RedisQueue{
		client:     client,
		handlers:   make(map[string][]Handler),
		workerPool: make(chan struct{}, maxWorkers),
		ctx:        ctx,
		cancel:     cancel,
	}

	// เติมพูลของ worker
	for i := 0; i < maxWorkers; i++ {
		rq.workerPool <- struct{}{}
	}

	// เริ่มตัวจัดการข้อความแบบหน่วงเวลา
	rq.wg.Add(1)
	go rq.processDelayedMessages()

	return rq
}

// Publish เผยแพร่ข้อความ
func (rq *RedisQueue) Publish(ctx context.Context, topic string, payload interface{}) error {
	// แปลง payload เป็น JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ไม่สามารถแปลง payload เป็น JSON ได้: %w", err)
	}

	// สร้างข้อความ
	msg := &Message{
		ID:         generateMessageID(),
		Topic:      topic,
		Payload:    data,
		Timestamp:  time.Now(),
		Retries:    0,
		MaxRetries: 3,
	}

	// แปลงข้อความเป็น JSON
	msgData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("ไม่สามารถแปลงข้อความเป็น JSON ได้: %w", err)
	}

	// เผยแพร่ไปยัง Redis
	key := fmt.Sprintf("queue:%s", topic)
	if err := rq.client.LPush(ctx, key, msgData).Err(); err != nil {
		return fmt.Errorf("ไม่สามารถเผยแพร่ข้อความได้: %w", err)
	}

	return nil
}

// Subscribe สมัครรับข้อมูลหัวข้อ
func (rq *RedisQueue) Subscribe(ctx context.Context, topic string, handler Handler) error {
	// ลงทะเบียนตัวจัดการ
	rq.mu.Lock()
	rq.handlers[topic] = append(rq.handlers[topic], handler)
	rq.mu.Unlock()

	// เริ่ม consumer
	rq.wg.Add(1)
	go rq.consume(topic)

	return nil
}

// consume บริโภคข้อความ
func (rq *RedisQueue) consume(topic string) {
	defer rq.wg.Done()

	key := fmt.Sprintf("queue:%s", topic)

	for {
		select {
		case <-rq.ctx.Done():
			return
		default:
			// ดึงข้อความจากคิว (บล็อค 1 วินาที)
			result, err := rq.client.BRPop(rq.ctx, time.Second, key).Result()
			if err != nil {
				if err == redis.Nil {
					continue // หมดเวลา รอต่อไป
				}
				// บันทึกข้อผิดพลาดและดำเนินการต่อ
				continue
			}

			if len(result) < 2 {
				continue
			}

			// ขอ token จากพูล worker
			<-rq.workerPool

			// ประมวลผลข้อความแบบไม่พร้อมกัน
			go func(data string) {
				defer func() {
					rq.workerPool <- struct{}{} // คืน token
				}()

				// แปลง JSON กลับเป็นข้อความ
				var msg Message
				if err := json.Unmarshal([]byte(data), &msg); err != nil {
					return
				}

				// ประมวลผลข้อความ
				rq.processMessage(&msg)
			}(result[1])
		}
	}
}

// processMessage ประมวลผลข้อความ
func (rq *RedisQueue) processMessage(msg *Message) {
	rq.mu.RLock()
	handlers := rq.handlers[msg.Topic]
	rq.mu.RUnlock()

	for _, handler := range handlers {
		ctx, cancel := context.WithTimeout(rq.ctx, 30*time.Second)
		err := handler(ctx, msg)
		cancel()

		if err != nil {
			// ประมวลผลล้มเหลว ลองใหม่
			if msg.Retries < msg.MaxRetries {
				msg.Retries++
				rq.retryMessage(msg)
			} else {
				// เกินจำนวนครั้งสูงสุด ส่งไปยัง dead letter queue
				rq.sendToDeadLetter(msg, err)
			}
		}
	}
}

// retryMessage ลองประมวลผลข้อความอีกครั้ง
func (rq *RedisQueue) retryMessage(msg *Message) {
	// คำนวณเวลาหน่วงในการลองใหม่ (exponential backoff)
	delay := time.Duration(msg.Retries) * time.Second * 2

	// เผยแพร่เป็นข้อความแบบหน่วงเวลา
	rq.PublishDelayed(rq.ctx, msg.Topic, msg, delay)
}

// sendToDeadLetter ส่งไปยัง dead letter queue
func (rq *RedisQueue) sendToDeadLetter(msg *Message, err error) {
	deadLetterKey := fmt.Sprintf("dead_letter:%s", msg.Topic)

	// เพิ่มข้อมูลข้อผิดพลาด
	type DeadLetterMessage struct {
		*Message
		Error    string    `json:"error"`
		FailedAt time.Time `json:"failed_at"`
	}

	dlMsg := &DeadLetterMessage{
		Message:  msg,
		Error:    err.Error(),
		FailedAt: time.Now(),
	}

	data, _ := json.Marshal(dlMsg)
	rq.client.LPush(rq.ctx, deadLetterKey, data)
}

// PublishDelayed เผยแพร่ข้อความแบบหน่วงเวลา
func (rq *RedisQueue) PublishDelayed(ctx context.Context, topic string, payload interface{}, delay time.Duration) error {
	// แปลง payload เป็น JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ไม่สามารถแปลง payload เป็น JSON ได้: %w", err)
	}

	// สร้างข้อความ
	msg := &Message{
		ID:        generateMessageID(),
		Topic:     topic,
		Payload:   data,
		Timestamp: time.Now(),
	}

	// แปลงข้อความเป็น JSON
	msgData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("ไม่สามารถแปลงข้อความเป็น JSON ได้: %w", err)
	}

	// เพิ่มไปยังคิวแบบหน่วงเวลา (ใช้ sorted set)
	score := float64(time.Now().Add(delay).Unix())
	key := "delayed_queue"
	if err := rq.client.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: msgData,
	}).Err(); err != nil {
		return fmt.Errorf("ไม่สามารถเผยแพร่ข้อความแบบหน่วงเวลาได้: %w", err)
	}

	return nil
}

// processDelayedMessages ประมวลผลข้อความแบบหน่วงเวลา
func (rq *RedisQueue) processDelayedMessages() {
	defer rq.wg.Done()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-rq.ctx.Done():
			return
		case <-ticker.C:
			// ดึงข้อความที่ครบกำหนด
			now := float64(time.Now().Unix())
			key := "delayed_queue"

			// ดึงข้อความที่ครบกำหนดทั้งหมด
			messages, err := rq.client.ZRangeByScore(rq.ctx, key, &redis.ZRangeBy{
				Min: "0",
				Max: fmt.Sprintf("%f", now),
			}).Result()

			if err != nil {
				continue
			}

			for _, msgData := range messages {
				// แปลง JSON กลับเป็นข้อความ
				var msg Message
				if err := json.Unmarshal([]byte(msgData), &msg); err != nil {
					continue
				}

				// เผยแพร่ไปยังคิวปกติ
				if err := rq.Publish(rq.ctx, msg.Topic, msg.Payload); err != nil {
					continue
				}

				// ลบออกจากคิวแบบหน่วงเวลา
				rq.client.ZRem(rq.ctx, key, msgData)
			}
		}
	}
}

// Close ปิดคิว
func (rq *RedisQueue) Close() error {
	rq.cancel()
	rq.wg.Wait()
	return nil
}

// generateMessageID สร้าง ID สำหรับข้อความ
func generateMessageID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Nanosecond())
}