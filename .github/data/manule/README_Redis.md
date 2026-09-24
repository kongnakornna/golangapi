# คู่มือการใช้งาน Redis (Cache, Session & Queue)

Redis ในโปรเจกต์ `icmongolang` ถูกใช้หลายบทบาท:
1. **Cache** — user cache, MQTT/iot topic cache, device list, chart
2. **Session/Token** — refresh token ของระบบ auth
3. **Task queue (asynq)** — worker ส่งอีเมล (DB index 1)
4. **Custom queue** — ระบบ queue ของ websocket hub / kafka demo (Redis list)

มี **3 ชุด config** ในไฟล์เดียว: `redis` (DB 0), `taskRedis` (DB 1) และ `RedisHelper` (hardcode localhost:6380)

---

## 1. การตั้งค่า Configuration

### 1.1 `config/config.default.yml`
```yaml
redis:
  Addr: localhost:6379
  Password:
  Db: 0
  MinIdleConns: 200
  PoolSize: 12000
  PoolTimeout: 240

taskRedis:
  Addr: localhost:6379
  Db: 1
  PoolTimeout: 240
```

### 1.2 Docker Compose (`docker-compose.yml:55-61`)
```yaml
redis:
  image: redis:7-alpine
  expose:
    - 6579       # ⚠️ มีผลเฉพาะเอกสาร — container ยังฟังที่ 6379
  volumes:
    - app-redis-data:/data
```

> ⚠️ `expose: 6579` ไม่ได้เปลี่ยน port ที่ redis ฟัง — container ยังเป็น 6379 (exporter ยังเชื่อม `redis:6379`)

---

## 2. การเชื่อมต่อจาก Go

### 2.1 `pkg/db/redis/redis_conn.go` — `NewRedis(cfg)`
```go
redis.NewClient(&redis.Options{
    Addr:         cfg.Redis.Addr,
    MinIdleConns: cfg.Redis.MinIdleConns,
    PoolSize:     cfg.Redis.PoolSize,
    PoolTimeout:  time.Duration(cfg.Redis.PoolTimeout) * time.Second,
    Password:     cfg.Redis.Password,
    DB:           cfg.Redis.Db,
})
```
Library: `github.com/redis/go-redis/v9`

### 2.2 `RedisCache` — interface สำหรับ MQTT/iot handler
```go
type Cache interface {
    Get(ctx context.Context, key string, dst interface{}) error
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}
// NewCache(client) → RedisCache (JSON marshal/unmarshal)
```

### 2.3 `RedisHelper` (⚠️ singleton hardcode)
`GetRedisHelper()` ใช้ `Addr: localhost:6380` **hardcode** — ไม่ได้อ่าน config; มี method: `Set/Get/Delete/Exists/TTL/Keys/FlushAll` (op timeout 5s)

---

## 3. Redis Keys จำแนกตามการใช้งาน

### 3.1 Cache keys
| Key pattern | TTL | ใช้ที่ |
|-------------|-----|-------|
| `sd_user:{id}` | 3600s | user cache (auth/users) |
| `mqtt_topic:<topic>` | 10s | MQTT/iot topic data |
| `mqtt_payload:<bucket>` | 30-60s | iot monitor/alarm |
| `mqtt_device_list:<md5>` | 15min | iot monitor group |
| `iot_group:<md5>` | 5min | iot monitor group |
| `mqtt_chart:<md5>` | 45s | iot charts |
| `mqtt_topic_data:<md5>` | 60s | MQTT gettopicdata |
| `device_control:<md5>` / `device_ctrl:<t>:<m>` | 5s | MQTT device control |

### 3.2 Token (auth)
| Key | ประเภท | หมายเหตุ |
|-----|--------|----------|
| `RefreshToken:{userID}` | SET | refresh tokens ของ user |

### 3.3 asynq task queue (DB 1)
| Key | หมายเหตุ |
|-----|----------|
| `asynq:{queue}:pending/active` | งานรอ/กำลังทำ |
| `asynq:delayed` / `asynq:retry` / `asynq:dead` | ดีเลย์ / retry / dead letter |

### 3.4 Custom queue module (websocket/kafka demo)
| Key | หมายเหตุ |
|-----|----------|
| `queue:<topic>` | list งาน (LPush/BRPop) |
| `dead_letter:<topic>` | งานที่ failed |
| `delayed_queue` | sorted set งานดีเลย์ |

---

## 4. คำสั่งใช้งาน

```bash
docker compose up -d redis
```

### เข้า Redis-cli
```bash
docker exec -it $(docker compose ps -q redis) redis-cli
```

### เช็คคีย์
```bash
redis-cli -n 0 keys '*'        # cache (DB 0)
redis-cli -n 1 keys 'asynq:*'  # task queue (DB 1)
redis-cli -n 0 TTL sd_user:xxx
```

### ล้าง cache / ข้อมูลทั้งหมด
```bash
redis-cli FLUSHALL             # ⚠️ กระทบทุก DB รวมถึง session
redis-cli -n 0 FLUSHDB
```

---

## 5. Edge Cases & ข้อควรระวัง

1. **flush Redis = logout ทุกคน + cache ว่าง** — refresh token อยู่ใน Redis; ล้างแล้วต้อง login ใหม่
2. **RedisHelper hardcode port 6380** — ถ้าไม่ได้รัน redis ที่ 6380 จะใช้ไม่ได้
3. **DB 0 vs DB 1** — cache/refresh ใช้ 0, asynq ใช้ 1; ระวัง command ที่ flush ทั้ง server
4. **Pool ขนาดใหญ่** — `PoolSize: 12000` เหมาะกับ concurrent สูง; ใช้ MinIdleConns 200
5. **worker ต้องชี้ DB 1 เดียวกับ API** — ไม่งั้น enqueue กับ consume คนละคิว
6. **ไม่มี password ใน default config** — เหมาะกับ internal network; ใส่ `redis.Password` ถ้าเซิร์ฟเวอร์มี auth

---

## 6. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| `dial tcp localhost:6379` fail | redis ไม่รัน → `docker compose up -d redis` |
| refresh token หาย / ถูกบังคับ logout | Redis ถูก flush หรือ restart (ใน memory) — ใช้ volume `app-redis-data` |
| asynq ไม่ process | ดู `taskRedis.Db` (ต้อง 1); worker รันไหม |
| งานค้างใน `asynq:dead` | ตรวจ worker log; SMTP error ตอนส่งเมล |
| `RedisHelper` connect fail | ต้องรัน redis ที่ 6380 (hardcode) |
