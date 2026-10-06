// Package kafka ให้ producer wrapper สำหรับ publish ข้อความไปยัง Kafka
//
// ไฟล์นี้รวม API ของ sarama-based producer เข้าไว้ที่เดียว:
//   - PublishMessage  → auto JSON serialize (ใช้โดย PDPA outbox dispatcher)
//   - PublishRaw      → ส่ง raw bytes ไม่ serialize
//   - NewProducer     → รองรับ optional timeout (variadic)
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

// Producer ห่อ sarama.SyncProducer พร้อม helper สำหรับ publish JSON
//
// ใช้ idempotent producer + RequiredAcks=WaitForAll เพื่อให้แน่ใจว่า
// ข้อความถูกส่งไปยัง Kafka อย่างน้อยหนึ่งครั้ง (at-least-once) และไม่ซ้ำ
// เหมาะกับ outbox dispatcher ของ PDPA
type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewProducer สร้าง Producer ใหม่สำหรับ topic ที่กำหนด
//
// รองรับ optional timeout (variadic) เพื่อ backward-compat กับ
// cmd/kafka.go ที่เรียก NewProducer(brokers, topic, timeout)
//
// ถ้าไม่ส่ง timeout จะใช้ค่า default ของ sarama (10s)
func NewProducer(brokers []string, topic string, timeouts ...time.Duration) (*Producer, error) {
	cfg := sarama.NewConfig()

	// --- Reliability ---
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Idempotent = true
	cfg.Net.MaxOpenRequests = 1

	// --- Optional timeout override ---
	if len(timeouts) > 0 && timeouts[0] > 0 {
		d := timeouts[0]
		cfg.Producer.Timeout = d
		cfg.Net.DialTimeout = d
		cfg.Net.ReadTimeout = d
		cfg.Net.WriteTimeout = d
	}

	p, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, fmt.Errorf("create kafka producer: %w", err)
	}
	return &Producer{producer: p, topic: topic}, nil
}

// PublishMessage serialize msg เป็น JSON (หรือใช้ []byte/json.RawMessage ตรงๆ)
// แล้วส่งไปยัง topic ด้วย key ที่กำหนด
//
// รองรับ payload หลายรูปแบบ:
//   - []byte           → ส่งตรง
//   - json.RawMessage  → ส่งตรง
//   - อื่นๆ            → json.Marshal
//
// signature ตรงกับ interface ที่ PDPA scheduler ใช้:
//
//	type EventPublisher interface {
//	    PublishMessage(key string, msg interface{}) error
//	}
func (p *Producer) PublishMessage(key string, msg interface{}) error {
	var payload []byte
	switch v := msg.(type) {
	case []byte:
		payload = v
	case json.RawMessage:
		payload = v
	case nil:
		payload = nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("marshal message: %w", err)
		}
		payload = b
	}

	_, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(payload),
	})
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}
	return nil
}

// PublishRaw ส่ง raw bytes โดยไม่ serialize
//
// รับ ctx เพื่อให้ caller ยกเลิกได้ (sarama.SyncProducer.SendMessage
// ไม่ใช้ ctx โดยตรง แต่เราเช็ค ctx ก่อนส่งเพื่อ fail-fast)
func (p *Producer) PublishRaw(ctx context.Context, key string, payload []byte) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	_, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(payload),
	})
	if err != nil {
		return fmt.Errorf("send raw message: %w", err)
	}
	return nil
}

// PublishMessageContext เหมือน PublishMessage แต่รับ ctx
// (ใช้ในกรณีที่ caller ต้องการยกเลิกได้)
func (p *Producer) PublishMessageContext(ctx context.Context, key string, msg interface{}) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return p.PublishMessage(key, msg)
}

// Topic คืนค่า topic ปลายทางของ producer นี้
func (p *Producer) Topic() string {
	return p.topic
}

// Close ปิด producer และ flush ข้อความที่ค้างอยู่
func (p *Producer) Close() error {
	if p.producer == nil {
		return nil
	}
	return p.producer.Close()
}
