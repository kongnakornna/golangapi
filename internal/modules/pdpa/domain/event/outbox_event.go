package event

import (
	"time"

	"github.com/google/uuid"
)

// OutboxEventStatus แสดงสถานะของ event ใน outbox
// OutboxEventStatus represents the status of an event in the outbox.
type OutboxEventStatus int8

const (
	// OutboxEventStatusPending คือ event ที่รอ dispatch
	// OutboxEventStatusPending means the event is pending dispatch.
	OutboxEventStatusPending OutboxEventStatus = 1

	// OutboxEventStatusProcessing คือ event ที่กำลังถูกประมวลผล
	// OutboxEventStatusProcessing means the event is being processed.
	OutboxEventStatusProcessing OutboxEventStatus = 2

	// OutboxEventStatusSuccess คือ event ที่ dispatch สำเร็จ
	// OutboxEventStatusSuccess means the event was dispatched successfully.
	OutboxEventStatusSuccess OutboxEventStatus = 3

	// OutboxEventStatusFailed คือ event ที่ dispatch ล้มเหลว
	// OutboxEventStatusFailed means the event failed and will be sent to the DLQ.
	OutboxEventStatusFailed OutboxEventStatus = 4
)

// IsValid ตรวจสอบว่าสถานะเป็นค่าที่รองรับหรือไม่
// IsValid reports whether the status is a supported value.
func (s OutboxEventStatus) IsValid() bool {
	switch s {
	case OutboxEventStatusPending,
		OutboxEventStatusProcessing,
		OutboxEventStatusSuccess,
		OutboxEventStatusFailed:
		return true
	default:
		return false
	}
}

// String คืนค่าสถานะในรูปแบบตัวอักษร
// String returns the status in uppercase form, or "UNKNOWN" when invalid.
func (s OutboxEventStatus) String() string {
	switch s {
	case OutboxEventStatusPending:
		return "PENDING"
	case OutboxEventStatusProcessing:
		return "PROCESSING"
	case OutboxEventStatusSuccess:
		return "SUCCESS"
	case OutboxEventStatusFailed:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
}

// OutboxEvent คือ event ที่ถูกเก็บไว้ใน outbox เพื่อรอการ dispatch ไปยัง Kafka
// OutboxEvent is an event stored in the outbox waiting to be dispatched to Kafka.
type OutboxEvent struct {
	// ID คือ identifier ของ event
	// ID is the event identifier.
	ID uuid.UUID

	// Topic คือ topic ปลายทางใน Kafka
	// Topic is the destination Kafka topic.
	Topic string

	// Key คือ partition key สำหรับการ publish
	// Key is the partition key used when publishing.
	Key string

	// Payload คือข้อมูล JSON ของ event
	// Payload is the JSON payload of the event.
	Payload []byte

	// Status คือสถานะปัจจุบันของ event
	// Status is the current status of the event.
	Status OutboxEventStatus

	// Attempts คือจำนวนครั้งที่พยายาม dispatch
	// Attempts is the number of dispatch attempts.
	Attempts int

	// CreatedAt คือเวลาที่ event ถูกสร้าง
	// CreatedAt is when the event was created.
	CreatedAt time.Time

	// DispatchedAt คือเวลาที่ event ถูก dispatch สำเร็จ
	// DispatchedAt is when the event was successfully dispatched. Nil when not yet dispatched.
	DispatchedAt *time.Time
}