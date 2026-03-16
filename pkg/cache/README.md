# คอมโพเนนต์แคช (pkg/cache)

## ออกแบบ: Redis เท่านั้น

แพ็คเกจนี้รองรับเฉพาะ Redis เป็นพื้นที่จัดเก็บแคชเท่านั้น ไม่มีการนำเสนอการใช้งานแคชในหน่วยความจำ โดยพิจารณาจากปัจจัยหลักดังนี้:

1. **หลีกเลี่ยงการเติบโตของหน่วยความจำที่ไม่สามารถควบคุมได้**
   - แคชภายในโพรเซสมักมีความเสี่ยงต่อการขยายตัวของหน่วยความจำและการรั่วไหล
   - ความจุไม่สามารถจัดเรียงอย่างเคร่งครัดกับโควต้าทรัพยากรของคอนเทนเนอร์/โหนดได้

2. **ความสอดคล้องแบบกระจาย**
   - ภายใต้การปรับใช้หลายอินสแตนซ์ แคชภายในเครื่องอาจทำให้ข้อมูลไม่สอดคล้องกัน
   - Redis จัดเตรียมแคชแบบรวมศูนย์ ซึ่งง่ายต่อการจัดการความสอดคล้อง

3. **เป็นมิตรต่อการตรวจสอบและการดำเนินงาน**
   - Redis มีระบบการแสดงผลและการแจ้งเตือนที่สมบูรณ์แบบ
   - สะดวกในการกำหนดขีดจำกัดหน่วยความจำและนโยบายการคัดออก

4. **ความสามารถในการกู้คืน**
   - Redis รองรับการคงอยู่ของข้อมูลและการกู้คืนหลังจากรีสตาร์ท
   - ไม่จำเป็นต้องมีกระบวนการ预热แคชเพิ่มเติม

## วิธีการใช้งาน

```go
// เริ่มต้นแคช (ขึ้นอยู่กับ Redis)
cacheOpts := cache.Options{
    RedisAddress:      "localhost:6379",
    RedisPassword:     "password",
    RedisDB:           0,
    DefaultExpiration: 10 * time.Minute,
    CleanupInterval:   5 * time.Minute,
}

cacheInstance, err := cache.NewCache(cacheOpts)
if err != nil {
    // หาก Redis ไม่พร้อมใช้งาน สามารถลดระดับเป็น Noop ได้
    cacheInstance = cache.NewNoop()
}

// Set / Get
_ = cacheInstance.Set(ctx, "key", []byte("value"), 5*time.Minute)
value, err := cacheInstance.Get(ctx, "key")

// SetObject / GetObject
user := &User{ID: 1, Name: "John"}
_ = cacheInstance.SetObject(ctx, "user:1", user, 10*time.Minute)

var cachedUser User
err = cacheInstance.GetObject(ctx, "user:1", &cachedUser)
```

## กลยุทธ์แคช

ภายในแพ็คเกจมีกลยุทธ์ดังต่อไปนี้ (ดูได้ที่ `pkg/cache/strategies.go`):

- **Cache-Aside**: เมื่อแคชไม่ตรง จะโหลดจากแหล่งข้อมูลต้นทางและเขียนลงแคช
- **SingleFlight**: คำขอพร้อมกันสำหรับคีย์เดียวกันจะเรียกแหล่งข้อมูลต้นทางเพียงครั้งเดียว

## การกำหนดค่า

การกำหนดค่า Redis มาจาก `app.redis`:

```yaml
app:
  redis:
    enabled: true
    host: localhost
    port: 6379
    password: ""
    db: 0
```

การแทนที่ด้วยตัวแปรสภาพแวดล้อมใช้คำนำหน้า `APP_` (เช่น `APP_REDIS_HOST`, `APP_REDIS_ENABLED`)

## ลักษณะการลดระดับ

- `cache.NewCache` ต้องการที่อยู่ Redis มิฉะนั้นจะส่งคืนข้อผิดพลาด
- หาก Redis ถูกปิดใช้งานหรือการเชื่อมต่อล้มเหลวขณะเริ่มต้นแอปพลิเคชัน จะใช้ `cache.NewNoop()`
- การใช้งาน `Noop` สำหรับ `Get/GetObject` จะส่งคืน `cache.ErrNotFound` ส่วนการดำเนินการเขียนจะไม่มีผลข้างเคียง