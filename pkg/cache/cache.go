package cache

import (
	"context"
	"errors"
	"time"
)

// ErrNotFound ข้อผิดพลาดเมื่อไม่พบคีย์ในแคช
var ErrNotFound = errors.New("cache: ไม่พบคีย์")

// Cache กำหนดอินเทอร์เฟซของแคช
type Cache interface {
	// Get ดึงค่าจากแคช
	Get(ctx context.Context, key string) ([]byte, error)

	// Set เก็บค่าในแคช
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error

	// Delete ลบคีย์ที่ระบุออกจากแคช
	Delete(ctx context.Context, key string) error

	// Clear ล้างแคชทั้งหมด
	Clear(ctx context.Context) error

	// GetObject ดึงค่าและแปลงเป็นออบเจ็กต์ตามชนิดที่ระบุ
	GetObject(ctx context.Context, key string, value interface{}) error

	// SetObject แปลงออบเจ็กต์เป็นซีเรียลไลซ์แล้วเก็บในแคช
	SetObject(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

// Options ตัวเลือกการตั้งค่าแคช
type Options struct {
	// RedisAddress ที่อยู่ของ Redis
	RedisAddress string

	// RedisPassword รหัสผ่านของ Redis
	RedisPassword string

	// RedisDB หมายเลขฐานข้อมูล Redis
	RedisDB int

	// DefaultExpiration อายุเริ่มต้นของข้อมูลในแคช
	DefaultExpiration time.Duration

	// CleanupInterval ช่วงเวลาการล้างข้อมูลที่หมดอายุ
	CleanupInterval time.Duration
}

// NewCache สร้างอินสแตนซ์แคช (รองรับเฉพาะ Redis)
func NewCache(opts Options) (Cache, error) {
	if opts.RedisAddress == " {
		return nil, errors.New("ไม่ได้กำหนดที่อยู่ Redis เนื่องจากบริการแคชต้องใช้ Redis")
	}
	return newRedisCache(opts)
}
