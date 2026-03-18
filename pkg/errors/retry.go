package errors

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// RetryConfig การกำหนดค่าสำหรับการลองใหม่
type RetryConfig struct {
	MaxAttempts     int              // จำนวนครั้งที่ลองสูงสุด
	InitialDelay    time.Duration    // หน่วงเริ่มต้น
	MaxDelay        time.Duration    // หน่วงสูงสุด
	Multiplier      float64          // ตัวคูณหน่วง
	RandomizeFactor float64          // ปัจจัยสุ่ม (ระหว่าง 0-1)
	RetryIf         func(error) bool // ฟังก์ชันตรวจสอบว่าควรลองใหม่หรือไม่
}

// DefaultRetryConfig ค่ากำหนดเริ่มต้นสำหรับการลองใหม่
var DefaultRetryConfig = RetryConfig{
	MaxAttempts:     3,
	InitialDelay:    100 * time.Millisecond,
	MaxDelay:        10 * time.Second,
	Multiplier:      2.0,
	RandomizeFactor: 0.1,
	RetryIf:         IsRetryable,
}

// RetryableFunc ประเภทฟังก์ชันที่สามารถลองใหม่ได้
type RetryableFunc func() error

// RetryableWithContextFunc ประเภทฟังก์ชันที่สามารถลองใหม่ได้พร้อม context
type RetryableWithContextFunc func(context.Context) error

// Retry ดำเนินการฟังก์ชันพร้อมการลองใหม่
func Retry(fn RetryableFunc, config *RetryConfig) error {
	return RetryWithContext(context.Background(), func(ctx context.Context) error {
		return fn()
	}, config)
}

// RetryWithContext ดำเนินการฟังก์ชันพร้อม context และการลองใหม่
func RetryWithContext(ctx context.Context, fn RetryableWithContextFunc, config *RetryConfig) error {
	if config == nil {
		config = &DefaultRetryConfig
	}

	var lastErr error

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		// ตรวจสอบว่า context ถูกยกเลิกหรือไม่
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context ถูกยกเลิก: %w", err)
		}

		// ดำเนินการฟังก์ชัน
		err := fn(ctx)
		if err == nil {
			return nil // สำเร็จ
		}

		lastErr = err

		// ตรวจสอบว่าควรลองใหม่หรือไม่
		if config.RetryIf != nil && !config.RetryIf(err) {
			return err // ข้อผิดพลาดที่ไม่สามารถลองใหม่ได้
		}

		// ถ้าเป็นครั้งสุดท้าย ให้คืนค่าข้อผิดพลาด
		if attempt == config.MaxAttempts-1 {
			break
		}

		// คำนวณเวลาหน่วง
		delay := calculateDelay(attempt, config)

		// รอหรือจนกว่า context จะถูกยกเลิก
		select {
		case <-time.After(delay):
			// ดำเนินการลองครั้งถัดไป
		case <-ctx.Done():
			return fmt.Errorf("context ถูกยกเลิกระหว่างรอลองใหม่: %w", ctx.Err())
		}
	}

	return &RetryError{
		LastError: lastErr,
		Attempts:  config.MaxAttempts,
	}
}

// calculateDelay คำนวณเวลาหน่วงสำหรับการลองใหม่
func calculateDelay(attempt int, config *RetryConfig) time.Duration {
	// การหน่วงแบบ exponential backoff
	delay := float64(config.InitialDelay) * math.Pow(config.Multiplier, float64(attempt))

	// เพิ่มการสุ่ม (jitter)
	if config.RandomizeFactor > 0 {
		randomFactor := 1.0 + (rand.Float64()*2-1)*config.RandomizeFactor
		delay *= randomFactor
	}

	// จำกัดไม่ให้เกินค่าสูงสุด
	if delay > float64(config.MaxDelay) {
		delay = float64(config.MaxDelay)
	}

	return time.Duration(delay)
}

// RetryError ข้อผิดพลาดจากการลองใหม่
type RetryError struct {
	LastError error
	Attempts  int
}

func (e *RetryError) Error() string {
	return fmt.Sprintf("ดำเนินการล้มเหลวหลังจากลอง %d ครั้ง: %v", e.Attempts, e.LastError)
}

func (e *RetryError) Unwrap() error {
	return e.LastError
}

// IsRetryable ตรวจสอบว่าข้อผิดพลาดสามารถลองใหม่ได้หรือไม่
func IsRetryable(err error) bool {
	// หากเป็นข้อผิดพลาดที่เกี่ยวกับ HTTP status code ให้พิจารณาตาม status code
	// 5xx, 429, 408 สามารถลองใหม่ได้
	// ที่นี่เขียนแบบง่าย ในทางปฏิบัติควรตรวจสอบตามประเภทข้อผิดพลาดจริง

	// ค่าเริ่มต้นคือไม่สามารถลองใหม่ได้
	return false
}

// ExponentialBackoff ลองใหม่แบบ exponential backoff
func ExponentialBackoff(fn RetryableFunc) error {
	config := &RetryConfig{
		MaxAttempts:     5,
		InitialDelay:    100 * time.Millisecond,
		MaxDelay:        30 * time.Second,
		Multiplier:      2.0,
		RandomizeFactor: 0.2,
		RetryIf:         IsRetryable,
	}
	return Retry(fn, config)
}

// LinearBackoff ลองใหม่แบบ linear backoff
func LinearBackoff(fn RetryableFunc) error {
	config := &RetryConfig{
		MaxAttempts:     3,
		InitialDelay:    1 * time.Second,
		MaxDelay:        5 * time.Second,
		Multiplier:      1.0,
		RandomizeFactor: 0,
		RetryIf:         IsRetryable,
	}
	return Retry(fn, config)
}

// RetryWithFixedDelay ลองใหม่แบบหน่วงเวลาคงที่
func RetryWithFixedDelay(fn RetryableFunc, delay time.Duration, attempts int) error {
	config := &RetryConfig{
		MaxAttempts:     attempts,
		InitialDelay:    delay,
		MaxDelay:        delay,
		Multiplier:      1.0,
		RandomizeFactor: 0,
		RetryIf:         IsRetryable,
	}
	return Retry(fn, config)
}

// CircuitBreaker วงจรตัด (Circuit Breaker)
type CircuitBreaker struct {
	maxFailures      int
	resetTimeout     time.Duration
	halfOpenRequests int

	failures        int
	lastFailureTime time.Time
	state           CircuitState
}

// CircuitState สถานะของวงจรตัด
type CircuitState int

const (
	// StateClosed ปิด (ปกติ)
	StateClosed CircuitState = iota
	// StateOpen เปิด (ตัดวงจร)
	StateOpen
	// StateHalfOpen เปิดครึ่ง (พยายามฟื้นคืน)
	StateHalfOpen
)

// NewCircuitBreaker สร้างวงจรตัดใหม่
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:      maxFailures,
		resetTimeout:     resetTimeout,
		halfOpenRequests: 1,
		state:            StateClosed,
	}
}

// Execute ดำเนินการฟังก์ชันพร้อมการป้องกันด้วยวงจรตัด
func (cb *CircuitBreaker) Execute(fn RetryableFunc) error {
	// ตรวจสอบสถานะวงจรตัด
	if cb.state == StateOpen {
		if time.Since(cb.lastFailureTime) < cb.resetTimeout {
			return &CircuitOpenError{
				ResetAt: cb.lastFailureTime.Add(cb.resetTimeout),
			}
		}
		// พยายามเข้าสู่สถานะ half-open
		cb.state = StateHalfOpen
		cb.halfOpenRequests = 1
	}

	// ดำเนินการฟังก์ชัน
	err := fn()

	if err != nil {
		cb.recordFailure()
		return err
	}

	cb.recordSuccess()
	return nil
}

// recordFailure บันทึกความล้มเหลว
func (cb *CircuitBreaker) recordFailure() {
	cb.failures++
	cb.lastFailureTime = time.Now()

	if cb.failures >= cb.maxFailures {
		cb.state = StateOpen
	}
}

// recordSuccess บันทึกความสำเร็จ
func (cb *CircuitBreaker) recordSuccess() {
	if cb.state == StateHalfOpen {
		cb.halfOpenRequests--
		if cb.halfOpenRequests <= 0 {
			cb.state = StateClosed
			cb.failures = 0
		}
	}
}

// CircuitOpenError ข้อผิดพลาดเมื่อวงจรตัดเปิดอยู่
type CircuitOpenError struct {
	ResetAt time.Time
}

func (e *CircuitOpenError) Error() string {
	return fmt.Sprintf("วงจรตัดเปิดอยู่ จะรีเซ็ตที่ %v", e.ResetAt)
}