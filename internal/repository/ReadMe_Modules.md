# การแปลเอกสารประกอบโค้ด Redis Repository

## แปลภาษาไทยของคอมเมนต์

```go
package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"icmongolang/internal"
	"icmongolang/pkg/httpErrors"
	"github.com/redis/go-redis/v9"
)

type RedisRepo[M any] struct {
	RedisClient *redis.Client
}

/*************  ✨ Windsurf Command ⭐  *************/
// CreateRedisRepo คืนค่า instance ใหม่ของ RedisRepo พร้อมกับ Redis client ที่กำหนด
// ใช้สำหรับสร้าง Redis repository สำหรับ model type M ที่กำหนด
// instance ที่คืนกลับมาจะใช้ในการติดต่อกับ Redis สำหรับการดำเนินการ CRUD
//
/*******  a88129c8-1c2d-4c0c-a671-c80971efaaea  *******/
func CreateRedisRepo[M any](redisClient *redis.Client) RedisRepo[M] {
	return RedisRepo[M]{RedisClient: redisClient}
}

func CreateRedisRepository[M any](redisClient *redis.Client) internal.RedisRepository[M] {
	return &RedisRepo[M]{RedisClient: redisClient}
}

func (r *RedisRepo[M]) Create(ctx context.Context, key string, exp *M, seconds int) error {
	// แปลงข้อมูลเป็น JSON byte array
	objBytes, err := json.Marshal(exp)
	if err != nil {
		return httpErrors.ErrJson(err)
	}

	// เก็บข้อมูลลง Redis โดยกำหนดเวลาหมดอายุเป็นวินาที
	if err = r.RedisClient.Set(ctx, key, objBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		// TODO: ใช้ httpErrors
		return err
	}
	return nil
}

func (r *RedisRepo[M]) Get(ctx context.Context, key string) (*M, error) {
	// ดึงข้อมูลจาก Redis เป็น byte array
	objBytes, err := r.RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		// ถ้าไม่พบคีย์ ให้คืนค่า nil โดยไม่มี error
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var obj M

	// แปลง JSON กลับเป็น struct
	if err = json.Unmarshal(objBytes, &obj); err != nil {
		return nil, httpErrors.ErrJson(err)
	}

	return &obj, nil
}

func (r *RedisRepo[M]) Delete(ctx context.Context, key string) error {
	// ลบข้อมูลจาก Redis
	if err := r.RedisClient.Del(ctx, key).Err(); err != nil {
		// ถ้าไม่พบคีย์ให้ถือว่าสำเร็จ (ไม่ต้อง return error)
		if errors.Is(err, redis.Nil) {
			return nil
		}
		// TODO: ใช้ httpErrors
		return err
	}
	return nil
}

/*************  ✨ Windsurf Command ⭐  *************/
// Sadd เพิ่มค่า value ที่กำหนดเข้าไปใน Redis set ที่มี key ตามที่กำหนด
//
// ถ้าค่าซ้ำกันอยู่ใน set แล้ว ฟังก์ชันจะคืนค่า err เป็น nil
//
// ถ้าคำสั่ง Redis SADD ล้มเหลว ฟังก์ชันจะคืนค่า error ที่ถูกห่อหุ้มด้วย httpErrors.Error
//
// หมายเหตุ: error ที่คืนกลับมาอาจถูกห่อหุ้มด้วย httpErrors.Error 
// ถ้า error นั้นเกี่ยวข้องกับการสื่อสารกับ Redis
/*******  460b09a0-018c-41a6-9645-a6b1f395d74a  *******/
func (r *RedisRepo[M]) Sadd(ctx context.Context, key string, value string) error {
	if err := r.RedisClient.SAdd(ctx, key, value).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepo[M]) Sadds(ctx context.Context, key string, values []string) error {
	if err := r.RedisClient.SAdd(ctx, key, values).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepo[M]) Srem(ctx context.Context, key string, value string) error {
	if err := r.RedisClient.SRem(ctx, key, value).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepo[M]) SIsMember(ctx context.Context, key string, value string) (bool, error) {
	result := r.RedisClient.SIsMember(ctx, key, value)
	if result.Err() != nil {
		return false, result.Err()
	}
	return result.Val(), nil
}

// func (r *RedisRepo[M]) SMembers(ctx context.Context, key string) ([]string, error) {
// 	result := r.RedisClient.SPop(ctx, key)
// 	if result.Err() != nil {
// 		return nil, result.Err()
// 	}
// 	return result.Val(), nil
// }
```

## คำอธิบายฟังก์ชันต่างๆ

### 1. **CreateRedisRepo / CreateRedisRepository**
สร้าง instance ของ Redis repository สำหรับทำงานกับ model type M

### 2. **Create** (CRUD - Create)
บันทึกข้อมูลลง Redis โดยกำหนดเวลาหมดอายุ (TTL) เป็นวินาที

### 3. **Get** (CRUD - Read)
อ่านข้อมูลจาก Redis ตาม key ที่กำหนด

### 4. **Delete** (CRUD - Delete)
ลบข้อมูลจาก Redis ตาม key ที่กำหนด

### 5. **Sadd** (Set Operation)
เพิ่มสมาชิกเข้าไปใน Redis Set

### 6. **Sadds** (Set Operation - Multiple)
เพิ่มสมาชิกหลายรายการเข้าไปใน Redis Set

### 7. **Srem** (Set Operation - Remove)
ลบสมาชิกออกจาก Redis Set

### 8. **SIsMember** (Set Operation - Check)
ตรวจสอบว่าสมาชิกอยู่ใน Redis Set หรือไม่

## ตัวอย่างการใช้งาน

```go
// สร้าง User model
type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// สร้าง Redis repository สำหรับ User
userRepo := CreateRedisRepo[User](redisClient)

// สร้างข้อมูล (Create)
user := &User{ID: "1", Name: "John", Email: "john@example.com"}
err := userRepo.Create(ctx, "user:1", user, 3600) // หมดอายุใน 1 ชั่วโมง

// อ่านข้อมูล (Read)
userData, err := userRepo.Get(ctx, "user:1")

// ลบข้อมูล (Delete)
err = userRepo.Delete(ctx, "user:1")

// ใช้งาน Redis Set
// เพิ่มค่า
err = userRepo.Sadd(ctx, "user:roles:1", "admin")
err = userRepo.Sadds(ctx, "user:roles:1", []string{"editor", "viewer"})

// ตรวจสอบว่ามีค่าอยู่หรือไม่
isMember, err := userRepo.SIsMember(ctx, "user:roles:1", "admin")

// ลบค่า
err = userRepo.Srem(ctx, "user:roles:1", "viewer")
```

## ตารางสรุปการทำงาน

| ฟังก์ชัน | การทำงาน | Redis Command | การใช้งาน |
|---------|----------|---------------|----------|
| Create | บันทึกข้อมูล | SET | เก็บข้อมูล JSON พร้อม TTL |
| Get | ดึงข้อมูล | GET | ดึงข้อมูล JSON และแปลงเป็น struct |
| Delete | ลบข้อมูล | DEL | ลบข้อมูลตาม key |
| Sadd | เพิ่มสมาชิก | SADD | เพิ่มค่าใน Set |
| Sadds | เพิ่มหลายสมาชิก | SADD | เพิ่มหลายค่าใน Set |
| Srem | ลบสมาชิก | SREM | ลบค่าออกจาก Set |
| SIsMember | ตรวจสอบสมาชิก | SISMEMBER | ตรวจสอบว่าค่าอยู่ใน Set หรือไม่ |