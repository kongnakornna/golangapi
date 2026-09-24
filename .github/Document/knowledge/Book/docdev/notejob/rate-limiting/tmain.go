// main_official.go
// ใช้ official library สำหรับ Token Bucket Rate Limiting
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/time/rate" // Official rate limiter package
)

// RateLimiterWithContext - ตัวอย่างการใช้งาน rate.Limiter แบบเต็มรูปแบบ
// Official Limiter ใช้ Token Bucket algorithm รองรับทั้ง blocking และ non-blocking
type RateLimiterWithContext struct {
	limiter *rate.Limiter
	name    string
}

// NewRateLimiterWithContext - สร้าง limiter ใหม่
// Parameters:
//   - rps: float64 - อัตรา request ต่อวินาที (Rate Limit)
//   - burst: int - ขนาด burst สูงสุด (ความจุ bucket)
func NewRateLimiterWithContext(name string, rps float64, burst int) *RateLimiterWithContext {
	// rate.Limit(rps) แปลง float64 เป็น type Limit
	// burst คือความจุของ bucket ที่จะเก็บ token สูงสุด
	limiter := rate.NewLimiter(rate.Limit(rps), burst)
	return &RateLimiterWithContext{
		limiter: limiter,
		name:    name,
	}
}

// Allow - Non-blocking check (ไม่รอ ถ้าไม่มี token ให้ปฏิเสธทันที)
func (rl *RateLimiterWithContext) Allow() bool {
	// limiter.Allow() เป็น non-blocking
	// ถ้ามี token: ใช้ token 1 ใบ และ return true
	// ถ้าไม่มี token: return false ทันที
	return rl.limiter.Allow()
}

// AllowN - Non-blocking check สำหรับหลาย token ต่อ request
func (rl *RateLimiterWithContext) AllowN(n int) bool {
	// ใช้ AllowN กรณี request หนึ่งต้องการ consume หลาย token
	return rl.limiter.AllowN(time.Now(), n)
}

// Wait - Blocking check (รอจนกว่าจะมี token)
func (rl *RateLimiterWithContext) Wait(ctx context.Context) error {
	// limiter.Wait จะ block การทำงานจนกว่าจะมี token พร้อม
	// ใช้ ctx เพื่อ support timeout หรือ cancellation
	return rl.limiter.Wait(ctx)
}

// ReserveAndDelay - ดูเวลาที่ต้องรอแล้วค่อยรอเอง
func (rl *RateLimiterWithContext) ReserveAndDelay() time.Duration {
	// Reserve จอง token ล่วงหน้า โดยไม่ต้องใช้ token ทันที
	reservation := rl.limiter.Reserve()
	if !reservation.OK() {
		return 0 // ไม่สามารถจองได้
	}
	// Delay คือระยะเวลาที่ต้องรอจนกว่าจะได้ token
	delay := reservation.Delay()
	time.Sleep(delay) // รอตาม delay ที่คำนวณได้
	return delay
}

func main() {
	fmt.Println("=== Official Token Bucket Rate Limiter Demo ===")
	fmt.Println("Limit: 3 requests/sec, Burst capacity: 5 tokens")
	fmt.Println()

	// 1. สร้าง limiter: อนุญาต 3 requests ต่อวินาที, burst สูงสุด 5 tokens
	limiter := NewRateLimiterWithContext("API-Limiter", 3.0, 5)

	// 2. ทดสอบการ Allow() (non-blocking)
	fmt.Println("--- Non-blocking Test (Allow) ---")
	for i := 1; i <= 8; i++ {
		if limiter.Allow() {
			fmt.Printf("Request %d: ✅ ALLOWED (token consumed)\n", i)
		} else {
			fmt.Printf("Request %d: ❌ REJECTED (no token available)\n", i)
		}
		time.Sleep(200 * time.Millisecond)
	}

	time.Sleep(2 * time.Second) // รอ token เติม
	fmt.Println("\n--- After waiting 2 seconds (tokens replenished) ---")
	for i := 1; i <= 3; i++ {
		if limiter.Allow() {
			fmt.Printf("Request %d: ✅ ALLOWED (token replenished)\n", i)
		}
	}

	// 3. ทดสอบการ Wait() (blocking)
	fmt.Println("\n--- Blocking Test (Wait with Context) ---")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 1; i <= 6; i++ {
		start := time.Now()
		err := limiter.Wait(ctx)
		elapsed := time.Since(start)

		if err != nil {
			log.Printf("Request %d: ❌ WAIT FAILED: %v\n", i, err)
			break
		}
		fmt.Printf("Request %d: ✅ PROCESSED after waiting %v\n", i, elapsed)
	}

	// 4. ทดสอบการใช้ AllowN (bulk token consumption)
	fmt.Println("\n--- Bulk Token Consumption Test (AllowN) ---")
	bulkLimiter := rate.NewLimiter(1, 3) // 1 request/sec, burst 3 tokens

	// พยายาม consume 5 tokens พร้อมกัน
	if bulkLimiter.AllowN(time.Now(), 5) {
		fmt.Println("✅ Consumed 5 tokens successfully")
	} else {
		fmt.Println("❌ Cannot consume 5 tokens at once (only 3 tokens available)")
	}

	// 5. ทดสอบการใช้ Reserve (pre-booking)
	fmt.Println("\n--- Reserve (Pre-booking) Test ---")
	reserveLimiter := rate.NewLimiter(2, 2) // 2 request/sec, burst 2 tokens

	// จอง token ล่วงหน้า
	r := reserveLimiter.Reserve()
	if r.OK() {
		fmt.Printf("Token reserved, delay required: %v\n", r.Delay())
		time.Sleep(r.Delay())
		fmt.Println("✅ Token obtained after delay")
		r.Cancel() // ยกเลิกการจองถ้าไม่ต้องการใช้
	}

	fmt.Println("\n=== Demo Complete ===")
}

// ผลลัพธ์ตัวอย่าง:
// Request 1: ✅ ALLOWED (token consumed)
// Request 2: ✅ ALLOWED (token consumed)
// Request 3: ✅ ALLOWED (token consumed)
// Request 4: ❌ REJECTED (no token available)
// ...
