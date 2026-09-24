# คู่มือการใช้งานโมดูล Kafka (Async Order Processing)

โมดูล `internal/modules/kafka` + `pkg/kafka` เป็นตัวอย่างระบบ **Asynchronous Processing** ที่ใช้ **Apache Kafka** เป็น message broker: รับคำสั่งซื้อ (Order) ผ่าน HTTP/WebSocket → ส่งเข้า Kafka topic → consumer ดึงออกมาประมวลผล → อัปเดตสถานะในฐานข้อมูล → broadcast ผลลัพธ์แบบ real-time ผ่าน WebSocket

สร้างด้วย library **`github.com/IBM/sarama`** (Apache Kafka client สำหรับ Go)

**สถาปัตยกรรมแบบครบวงจร:**

```
Client
 │ POST /orders (HTTP)          WebSocket /ws
 ▼                                ▼
┌───────────────────────────────────────────────┐
│ delivery/http/order  (CreateOrder)             │
│ delivery/ws  (ServeWS → HandleWebSocketMessage)│
└───────────────────────────────────────────────┘
        │  สร้าง Order + ส่ง Kafka message
        ▼
┌───────────────────────┐
│ pkg/kafka/Producer     │ ──▶ Kafka Topic (icmon-events) ──▶ Consumer Group
│ sarama.SyncProducer    │                                     │
└───────────────────────┘                                     ▼
                                                     ┌──────────────────────┐
                                                     │ pkg/kafka/Consumer    │
                                                     │ ProcessOrderMessage   │
                                                     └──────────────────────┘
                                                             │
                                                             ▼
                                                     ┌──────────────────────┐
                                                     │ usecase: อัปเดต Order  │
                                                     │ PENDING → PROCESSED   │
                                                     └──────────────────────┘
                                                             │
                                                             ▼ (WebSocket)
                                                     hub.BroadcastMessage("order_updated", ...)
```

> **สำคัญ:** โมดูลนี้เป็น **service แยกต่างหาก** — **ไม่ได้ถูก register ใน main API server** (`internal/server/handlers.go` ไม่มี route ของ Kafka) ต้องรันจาก `cmd/kafka` เอง

---

## 1. ตารางไฟล์

| ไฟล์ | บทบาท |
|------|-------|
| `cmd/kafka/main.go` | Main — wire ทุกอย่าง (DB, producer, consumer, hub, HTTP/WS) แล้วรัน server |
| `cmd/kafka/kafka.go` | Cobra subcommand `kafka` (เหมือน main แต่เรียกผ่าน CLI `icmongolang kafka`) |
| `internal/modules/kafka/delivery/http/handler.go` | REST handler `POST /orders` |
| `internal/modules/kafka/delivery/ws/handler.go` | WebSocket handler `GET /ws` (รับคำสั่งซื้อผ่าน WS) |
| `internal/modules/kafka/usecase/order_usecase.go` | Business logic: สร้าง order, publish, process, broadcast |
| `internal/modules/kafka/repository/order_repo.go` | GORM repository ของตาราง `orders` |
| `internal/modules/kafka/models/order.go` | `Order`, `CreateOrderRequest` |
| `pkg/kafka/producer.go` | `Producer` — ส่งข้อความเข้า Kafka (sync) |
| `pkg/kafka/consumer.go` | `Consumer` — รับข้อความจาก Kafka (consumer group) |
| `pkg/kafka/message.go` | `OrderMessage` — รูปแบบข้อความใน Kafka |
| `pkg/kafka/config.go` | `LoadConfig()` อ่าน env `KAFKA_*` (fallback) |
| `migrations/20250619_kafka_tables.sql` | ตาราง `orders` |
| `internal/modules/kafka/ws/handler.go` | ⚠️ ไฟล์ซ้ำของ delivery/http handler (package `http`, ไม่ได้ถูกใช้) |

---

## 2. การตั้งค่า Configuration

### 2.1 ตัวแปรสภาพแวดล้อม (`.env`)

```env
KAFKA_BROKERS=localhost:9092        # broker address (support หลายตัวคั่นด้วย comma)
KAFKA_TOPIC=icmon-events            # ชื่อ topic
KAFKA_GROUP_ID=icmon-consumer-group # consumer group id
KAFKA_TIMEOUT=30                    # network timeout หน่วยวินาที (default 30)
```

> หมายเหตุ: `pkg/kafka/config.go` มี `LoadConfig()` แยกที่อ่าน env เดียวกัน แต่ service จริงใช้ config ผ่าน `cfg.Kafka.*` (จาก `config.GetCfg()`)

### 2.2 ค่าเริ่มต้นใน `config/config.default.yml`

```yaml
kafka:
  brokers:
    - "localhost:9092"
  topic: "icmon-events"
  group_id: "icmon-consumer-group"
  timeout: 30
```

### 2.3 Docker

`docker-compose.yml` มี **apache/kafka:3.7.1** (KRaft mode — ไม่ต้องใช้ ZooKeeper) map port `9092:9092` พร้อม `kafka-exporter` port `9308` สำหรับ Prometheus:

```bash
docker compose up -d kafka
```

---

## 3. การรัน Service

### วิธีที่ 1 — รันตรง
```bash
go run ./cmd/kafka
```
### วิธีที่ 2 — ผ่าน Cobra CLI
```bash
go run ./main.go kafka        # หรือ go run ./cmd kafka
```

- **HTTP port:** ค่า `SERVER_PORT` จาก config (default `5051`)
- **Swagger UI:** `http://localhost:5051/swagger/index.html`
- **WebSocket:** `ws://localhost:5051/ws`

ตอนรัน service จะสร้าง:
1. PostgreSQL connection (เก็บตาราง `orders`)
2. Kafka **Producer** (sync)
3. Kafka **Consumer** (consumer group) — เริ่ม consume ทันที
4. WebSocket **Hub** — ใช้ broadcast สถานะ order

---

## 4. API Reference

### 4.1 `POST /orders` — สร้างคำสั่งซื้อ (Async ผ่าน Kafka)

- ตอบกลับ **HTTP 202 Accepted** ทันที — การประมวลผลเกิดขึ้นแบบ async โดย consumer

**Request body:**

```json
{
  "product_id": "PLC-S7-1500",
  "quantity": 5
}
```

| Field | ประเภท | บังคับ | เงื่อนไข |
|-------|--------|-------|----------|
| `product_id` | string | ✅ | ต้องไม่ว่าง |
| `quantity` | int | ✅ | ต้อง ≥ 1 |

```bash
curl -X POST http://localhost:5051/orders \
  -H "Content-Type: application/json" \
  -d '{ "product_id": "PLC-S7-1500", "quantity": 5 }'
```

**Response 202:**

```json
{
  "order_id": "3f2b...-uuid",
  "status": "PENDING",
  "message": "Order accepted for processing"
}
```

**ลำดับการทำงาน (ดู `usecase.go:41`):**
1. สร้าง `Order` สถานะ `PENDING` (uuid ใหม่)
2. สร้าง `OrderMessage` แล้ว `producer.PublishMessage(orderID, msg)` → Kafka topic
3. บันทึก order ลง PostgreSQL
4. ตอบ 202

---

### 4.2 `GET /ws` — WebSocket (ส่งคำสั่งซื้อ + รับ real-time status)

- Client เชื่อมต่อแล้วส่ง JSON ตามรูปแบบ `CreateOrderRequest` ได้เลย

**Query params:**

| Param | คำอธิบาย | ค่าเริ่มต้น |
|-------|----------|------------|
| `user_id` | ระบุตัวผู้ส่ง | `anonymous` |
| `room` | ให้ client เข้าร่วม room ทันทีที่เชื่อม | — |

**ตัวอย่าง (JavaScript):**

```js
const ws = new WebSocket("ws://localhost:5051/ws?user_id=user-001&room=orders");
ws.onopen = () => {
  ws.send(JSON.stringify({ product_id: "PLC-S7-1500", quantity: 3 }));
};
ws.onmessage = (evt) => console.log(JSON.parse(evt.data));
```

**Server ตอบกลับ client คนนั้น (ตอนสร้าง order สำเร็จ):**

```json
{
  "event": "order_created",
  "order_id": "3f2b...-uuid",
  "status": "PENDING"
}
```

**ทุก client (ตอน consumer ประมวลผลเสร็จ):**

```json
{
  "event": "order_updated",
  "order_id": "3f2b...-uuid",
  "status": "PROCESSED",
  "timestamp": "2026-08-14T10:30:00+07:00"
}
```

**ลำดับการทำงาน (ดู `usecase.go:94`):**
1. client ส่งข้อความ → `HandleWebSocketMessage` decode เป็น `CreateOrderRequest`
2. สร้าง order `PENDING` + publish ไป Kafka + บันทึก DB
3. ส่ง `order_created` กลับเฉพาะ client นั้น (`client.SendMessage`)
4. ต่อมา consumer รับข้อความ → process → `BroadcastOrderStatus` → **ทุก client** ได้รับ `order_updated`

---

### 4.3 `GET /health` — Health check

```bash
curl http://localhost:5051/health
# OK
```

---

## 5. วงจรการประมวลผล (Processing Pipeline)

```
POST /orders (หรือ WS)
   │
   ├─ 1. Producer.PublishMessage(key=orderID, msg=OrderMessage)
   │        Topic: icmon-events
   │
   ├─ 2. repo.Create(order)  →  INSERT orders (status=PENDING)
   │
   └─ (async) ──────────────────────────────────────────────
        Consumer (consumer group "icmon-consumer-group")
           │ ConsumeClaim → json.Unmarshal → handler
           ▼
        usecase.ProcessOrderMessage(msg)
           │ 1. FindByID(orderID) — ถ้าไม่เจอ → error
           │ 2. ถ้า status != PENDING → return (ข้าม, ถือว่าประมวลผลแล้ว)
           │ 3. status = PROCESSED, Update DB
           ▼
        BroadcastOrderStatus → hub.BroadcastMessage("order_updated", ...)
```

**รูปแบบข้อความใน Kafka (`OrderMessage`):**

```json
{
  "order_id": "3f2b...-uuid",
  "product_id": "PLC-S7-1500",
  "quantity": 5,
  "created_at": "2026-08-14T10:29:59Z"
}
```

- **Key** ของข้อความ = `order_id` → ข้อความของ order เดียวกันจะไป partition เดียวกันเสมอ (เรียงลำดับ)

---

## 6. Producer กับ Consumer (ดู `pkg/kafka`)

### 6.1 Producer (Sync — `sarama.SyncProducer`)

```go
producer, err := kafka.NewProducer(brokers, topic, timeoutSec)
defer producer.Close()
err := producer.PublishMessage(orderID, msg)   // marshals เป็น JSON แล้ว SendMessage
```

การตั้งค่า:

| ตัวเลือก | ค่า | ความหมาย |
|---------|-----|----------|
| `Producer.RequiredAcks` | `WaitForAll` | ต้องได้รับ ack จากทุก replica (แข็งแรงที่สุด) |
| `Producer.Retry.Max` | 5 | ส่งใหม่สูงสุด 5 ครั้งถ้าล้มเหลว |
| `Producer.Return.Successes` | true | จำเป็นสำหรับ SyncProducer |
| `Net.*Timeout` | จาก config (default 30s) | timeout ของ dial/read/write |

**ข้อควรรู้:** SyncProducer จะ **block จนกว่าจะรู้ผล** (partition + offset) — เหมาะกับงานที่ต้องยืนยันความสำเร็จ

### 6.2 Consumer (Consumer Group — `sarama.ConsumerGroup`)

```go
consumer, err := kafka.NewConsumer(brokers, groupID, topic, handler, timeoutSec)
ctx, cancel := context.WithCancel(context.Background())
consumer.Start(ctx)   // goroutine วน Consume ต่อเนื่อง
consumer.Wait()       // รอให้หยุด
consumer.Close()
```

การตั้งค่า:

| ตัวเลือก | ค่า | ความหมาย |
|---------|-----|----------|
| `Rebalance.Strategy` | `RoundRobin` | แบ่ง partition ให้สมาชิกใน group แบบ round-robin |
| `Offsets.Initial` | `OffsetNewest` | ⚠️ อ่านเฉพาะข้อความใหม่ — ข้อความที่ produce ก่อน consumer start จะถูกข้าม |
| `Consumer.Return.Errors` | true | ส่ง error ไปยัง channel |

**Retry behavior:**
- ถ้า `json.Unmarshal` ผิด → log แล้ว `MarkMessage` (ข้ามข้อความนั้น)
- ถ้า `handler` คืน error → log แล้ว **ไม่ Mark** (ข้อความจะถูก redeliver เมื่อ rebalance/restart ครั้งหน้า)
- ถ้า handler สำเร็จ → `MarkMessage` (commit offset)

---

## 7. ตารางข้อมูล (PostgreSQL)

### `orders` (จาก `migrations/20250619_kafka_tables.sql`)

| Column | Type | หมายเหตุ |
|--------|------|----------|
| `id` | UUID PK | ใช้ `gen_random_uuid()` ฝั่ง app (ไม่ใช่ DB default) |
| `product_id` | VARCHAR(255) NOT NULL | |
| `quantity` | INT NOT NULL | |
| `status` | VARCHAR(50) NOT NULL | default `PENDING`; มีค่า `PROCESSED` เมื่อ process เสร็จ |
| `created_at` | TIMESTAMPTZ | default NOW() |
| `updated_at` | TIMESTAMPTZ | default NOW() |

Index: `idx_orders_status`

> ตรวจสอบว่า migration รันแล้วก่อน start service (ตารางต้องมีอยู่)

---

## 8. การเรียกใช้ component เอง (Go)

```go
import (
    "context"
    "icmongolang/pkg/kafka"
)

// Producer
producer, _ := kafka.NewProducer([]string{"localhost:9092"}, "icmon-events", 30)
producer.PublishMessage("order-123", kafka.OrderMessage{
    OrderID:   "order-123",
    ProductID: "PLC-S7-1500",
    Quantity:  5,
})

// Consumer
consumer, _ := kafka.NewConsumer([]string{"localhost:9092"}, "my-group", "icmon-events",
    func(msg *kafka.OrderMessage) error {
        // ประมวลผลข้อความ
        return nil
    }, 30)
ctx, cancel := context.WithCancel(context.Background())
consumer.Start(ctx)
consumer.Wait()
```

---

## 9. Edge Cases / ข้อควรระวัง

1. **โมดูลไม่เชื่อมกับ main API server** — `POST /orders`, `/ws` (kafka) มีเฉพาะเมื่อรัน `cmd/kafka` ที่ port 5051 เท่านั้น อย่าหา endpoint นี้ใน port 5000
2. **`OffsetNewest`** — ถ้า consumer ยังไม่ start ตอน producer ส่งข้อความ → ข้อความแรกๆ จะไม่ถูกประมวลผล (ถูกข้าม) ถ้าต้องการรับข้อความเก่าด้วย ต้องเปลี่ยนเป็น `OffsetOldest`
3. **สถานะซ้ำ (idempotency)** — `ProcessOrderMessage` ข้าม order ที่ status ≠ `PENDING` อยู่แล้ว → กันการประมวลผลซ้ำจากการ redeliver
4. **Order ที่ error จะติดค้าง PENDING** — ถ้า `FindByID` ไม่เจอ (เช่น producer ส่งก่อน consumer มีข้อมูล) ข้อความจะไม่ถูก commit และจะถูก redeliver ในภายหลัง
5. **SyncProducer block** — ถ้า Kafka down การ publish จะค้าง/error → `POST /orders` จะตอบ 500 (ทั้งนี้ producer ใช้ ack `WaitForAll` + retry 5 ครั้ง)
6. **Order ของ file ซ้ำ** — `internal/modules/kafka/ws/handler.go` เป็นสำเนาของ `delivery/http/handler.go` (package `http`) และไม่ได้ถูกเรียกใช้ — อย่าลืมลบถ้าเห็นว่าเป็นของเหลือ
7. **การ broadcast ไปทุก client** — `BroadcastOrderStatus` ใช้ `BroadcastMessage` = ส่งให้ **ทุก** client ที่เชื่อมต่อ (ไม่ filter ตาม room/order)
8. **Auth** — endpoints ของโมดูลนี้ **ไม่มีการตรวจ JWT** (ต่างจากโมดูลอื่น) — ควรใส่ middleware เองถ้าจะ deploy จริง

---

## 10. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| `failed to create producer` | Kafka ยังไม่ขึ้น / `KAFKA_BROKERS` ผิด; ตรวจ `docker compose ps` และ `curl` broker |
| `POST /orders` ตอบ 500 | Kafka down ตอน publish หรือ PostgreSQL error ตอน save |
| order ค้าง `PENDING` ตลอด | consumer ยังไม่ start ตอน publish (`OffsetNewest` ข้ามข้อความเก่า); ตรวจ log consumer; ตรวจว่า consumer group ไม่ได้ active อยู่ |
| ไม่ได้รับ `order_updated` ผ่าน WS | ยังไม่ได้เปิด `/ws` ไว้ก่อน; broadcast ส่งให้ client ที่เชื่อมอยู่ตอนนั้นเท่านั้น |
| ข้อความถูกประมวลผลซ้ำ | เป็น normal ของ redelivery (ไม่ commit ตอน error) — มี idempotency guard ด้วย status อยู่แล้ว |
| ได้รับเฉพาะข้อความใหม่ ไม่ได้ข้อความเก่า | เพราะ `OffsetNewest`; ตั้ง `Offsets.Initial = OffsetOldest` ถ้าต้องการ consume ตั้งแต่ต้น |
| ตาราง `orders` ไม่มี (query error) | ยังไม่ได้รัน migration `20250619_kafka_tables.sql` |
