```go
rateLimitConfig := middleware.RateLimitConfig{
	RequestsPerSecond: 5,     // 5 requests per second | 5 คำขอต่อวินาที
	Burst:             10,    // Burst up to 10 | อนุญาตให้ส่งสูงสุด 10 คำขอในครั้งเดียว
	CleanupInterval:   10 * time.Minute, // Clean up every 10 minutes | ทำความสะอาดทุก 10 นาที
}
```

## **คำอธิบายเพิ่มเติม (Additional explanation):**

```go
rateLimitConfig := middleware.RateLimitConfig{
	// RequestsPerSecond: จำนวนคำขอสูงสุดที่อนุญาตต่อวินาที
	// Maximum number of requests allowed per second
	RequestsPerSecond: 5,     // 5 requests per second | 5 คำขอต่อวินาที
	
	// Burst: จำนวนคำขอสูงสุดที่อนุญาตให้ส่งพร้อมกัน (กระชุ)
	// Maximum number of requests allowed to burst at once
	Burst:             10,    // Burst up to 10 | อนุญาตให้ส่งสูงสุด 10 คำขอในครั้งเดียว
	
	// CleanupInterval: ระยะเวลาในการทำความสะอาดข้อมูล IP ที่ไม่ใช้งาน
	// Interval for cleaning up inactive IP data
	CleanupInterval:   10 * time.Minute, // Clean up every 10 minutes | ทำความสะอาดทุก 10 นาที
}
```

## **ตัวอย่างการใช้งานจริง (Usage example):**

```go
// สำหรับ Login endpoint (เข้มงวดมาก/Very strict)
loginRateLimit := middleware.RateLimitConfig{
	RequestsPerSecond: 3,     // 3 requests per second | 3 คำขอต่อวินาที
	Burst:             5,     // Burst up to 5 | กระชุได้สูงสุด 5 ครั้ง
	CleanupInterval:   5 * time.Minute, // Clean up every 5 minutes | ทำความสะอาดทุก 5 นาที
}

// สำหรับ API ทั่วไป (ปานกลาง / Moderate)
apiRateLimit := middleware.RateLimitConfig{
	RequestsPerSecond: 20,    // 20 requests per second | 20 คำขอต่อวินาที
	Burst:             40,    // Burst up to 40 | กระชุได้สูงสุด 40 ครั้ง
	CleanupInterval:   10 * time.Minute, // Clean up every 10 minutes | ทำความสะอาดทุก 10 นาที
}

// สำหรับ Public endpoints (หลวม/Loose)
publicRateLimit := middleware.RateLimitConfig{
	RequestsPerSecond: 50,    // 50 requests per second | 50 คำขอต่อวินาที
	Burst:             100,   // Burst up to 100 | กระชุได้สูงสุด 100 ครั้ง
	CleanupInterval:   15 * time.Minute, // Clean up every 15 minutes | ทำความสะอาดทุก 15 นาที
}
```