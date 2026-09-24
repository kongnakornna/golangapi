# คู่มือการใช้งานโมดูล WebSocket (Real-time)

โมดูล `internal/modules/websocket` + `pkg/websocket` ทำหน้าที่ส่งข้อมูล **แบบ Real-time** ไปยัง Client ผ่านโปรโตคอล **WebSocket** ทำงานร่วมกับ **Redis Queue** (pub/sub ของข้อความ) และ **PostgreSQL** (เก็บ history) เพื่อรองรับ 2 สถาปัตยกรรมหลัก:

1. **Standalone WebSocket Service** (`cmd/websocket/main.go`) — บริการแยกตัว เปิด port ของตัวเอง มีครบทั้ง WebSocket + REST history/rooms + auth ผ่านตาราง `ws_sessions`
2. **ฝังใน Main API Server** (`/api/ws`) — ใช้สำหรับรับ real-time จาก **MQTT** (topic `BAACTW05/DATA` → ห้อง `BAACTW05`) แล้ว broadcast ไปยัง client ที่ join ห้องนั้น

รูปแบบที่ Client ได้รับเสมอคือ frame JSON:
```json
{ "event": "<event_name>", "data": <payload> }
```

---

## 1. สถาปัตยกรรม (Architecture)

```
 Client (browser / mobile / app)
        │  WebSocket (ws://)
        ▼
┌─────────────────────────────────────────────────────────┐
│ delivery/ws/handler.go                                   │
│   ServeHTTP → upgrade HTTP→WS → authenticate → readPump   │
└─────────────────────────────────────────────────────────┘
        ▼
┌─────────────────────────────────────────────────────────┐
│ pkg/websocket/hub.go  (Hub)                              │
│   clients / topics / rooms + broadcast                   │
└──────────┬─────────────────────────────┬────────────────┘
           │                             │
           ▼                             ▼
┌──────────────────────┐      ┌──────────────────────────────┐
│ Redis Queue           │      │ PostgreSQL                   │
│ (internal/modules/    │      │ ws_messages (history)        │
│  queue/manager.go)    │      │ ws_sessions (WS auth token)  │
└──────────────────────┘      └──────────────────────────────┘
           │
           ▼
 MQTT service (Broadcaster) → hub.BroadcastToRoom("BAACTW05", "mqtt", data)
```

### 1.1 ตารางไฟล์

| ไฟล์ | บทบาท |
|------|-------|
| `cmd/websocket/main.go` | Main ของ Standalone WebSocket Service (wire ทั้งหมด) |
| `internal/modules/websocket/delivery/ws/handler.go` | ตัวจัดการ WebSocket (upgrade + รับ/ส่ง frame) |
| `internal/modules/websocket/delivery/http/handler.go` | REST endpoints (`/api/ws/messages`, `/api/ws/rooms`) |
| `internal/modules/websocket/usecase/ws_usecase.go` | Business logic — auth, save history, publish queue, broadcast |
| `internal/modules/websocket/repository/ws_repo.go` | Interface ของ repository |
| `internal/modules/websocket/repository/postgres/ws_repo_pg.go` | Repository ใช้ `sql.DB` กับตาราง `ws_messages` / `ws_sessions` |
| `internal/modules/websocket/models/ws_models.go` | `WSMessage`, `Session` |
| `internal/modules/websocket/presenter/ws_presenter.go` | DTO สำหรับ REST endpoints |
| `internal/modules/websocket/interface.go` | `Broadcaster` interface (ให้ MQTT ใช้ broadcast ผ่าน Hub) |
| `pkg/websocket/hub.go` | Hub — เก็บ client/topic/room + broadcast ไปทุกทิศทาง |
| `pkg/websocket/client.go` | Client — goroutine `WritePump` ส่งข้อความ + ping |
| `pkg/websocket/message.go` | รูปแบบ envelope ที่ client ส่งเข้า |
| `pkg/websocket/metrics.go` | Prometheus metrics ของ WebSocket |
| `internal/modules/queue/manager.go` | Redis Queue — Publish/Subscribe/retry/dead-letter |
| `migrations/20250619_websocket_tables.sql` | ตาราง `ws_messages` + `ws_sessions` |

---

## 2. สองโหมดการใช้งาน (Deployment Modes)

### 2.1 โหมด A — Standalone Service (`cmd/websocket`)

ใช้รันเป็น service แยก (เช่นใน Docker) เปิด endpoints ครบ:

| Path | Method | ความหมาย |
|------|--------|----------|
| `/ws` | GET | WebSocket endpoint (upgrade) |
| `/api/ws/messages` | POST | Broadcast ข้อความไปห้อง |
| `/api/ws/messages` | GET | ดึง history ข้อความ |
| `/api/ws/rooms` | GET | รายชื่อห้องที่ active |
| `/api/ws/rooms/{room}/stats` | GET | จำนวน client ในห้อง |
| `/health` | GET | Health check |

รันได้ด้วย:
```bash
go run ./cmd/websocket
# หรือ
go build -o websocket.exe ./cmd/websocket && ./websocket.exe
```

**Environment ที่ต้องมี:**
```env
DB_DSN=postgres://...           # PostgreSQL (เก็บ history + session)
REDIS_ADDR=localhost:6379       # Redis (queue)
WEBSOCKET_PORT=8080             # port ของ service (default 8080)
```

### 2.2 โหมด B — ฝังใน Main API Server

ใน `internal/server/handlers.go`:
```go
wsHub := websocket.NewHub(queue.NewQueue(), logger)
go wsHub.Run()
...
wsHandler := wsdelivery.NewWsHandler(wsHub, nil, logger)  // usecase = nil
apiRouter.Get("/ws", wsHandler.ServeHTTP)                 // → /api/ws
```

**ความต่างสำคัญจากโหมด A:**
- ใช้ `usecase = nil` → **ข้ามการตรวจ auth** (ทุกคนเป็น `"anonymous"`) และ **ไม่บันทึก history**
- `wsHub` ถูกส่งให้ MQTT ใช้เป็น `Broadcaster` → เมื่อ MQTT มีข้อความเข้า จะ broadcast ให้ client ในห้อง
- อยู่ที่ path `/api/ws` (ภายใต้ prefix `/api`)

---

## 3. แนวคิดหลัก: Topic กับ Room

Hub เก็บ client 2 แบบ subscription:

| | Room | Topic |
|---|------|-------|
| กลไก | ในหน่วยความจำ (in-memory) | ผ่าน Redis Queue (subscribe) |
| Command | `join_room` / `leave_room` | `subscribe` / `unsubscribe` |
| Broadcast | `BroadcastToRoom(room, event, data)` | `BroadcastToTopic(topic, event, data)` |
| ใช้กับ | MQTT real-time (ห้อง = ส่วนแรกของ topic เช่น `BAACTW05`) | ระบบ pub/sub ที่ต้องข้าม instance ผ่าน Redis |
| ผูกกับ MQTT | ใช่ — MQTT topic `BAACTW05/DATA` → broadcast เข้า room `BAACTW05` | ไม่โดยตรง |

> **ตัวอย่าง MQTT:** ข้อความเข้าที่ MQTT topic `BAACTW05/DATA` → `extractRoomFromTopic()` ตัดส่วนหน้า `BAACTW05` → `BroadcastToRoom("BAACTW05", "mqtt", {topic, payload})` → client ที่ join ห้อง `BAACTW05` ได้รับ frame:
> ```json
> { "event": "mqtt", "data": { "topic": "BAACTW05/DATA", "payload": "..." } }
> ```

---

## 4. WebSocket Protocol

### 4.1 เชื่อมต่อ

```
GET /ws?token=<WS_TOKEN>
```
หรือส่ง token ใน header `Authorization: Bearer <WS_TOKEN>`

- **โหมด A (Standalone):** token ต้องตรงกับแถวในตาราง `ws_sessions` (`token` + ยังไม่หมดอายุ `expires_at > NOW()`) ถ้าไม่ผ่าน → HTTP 401
- **โหมด B (Main API):** ข้ามการตรวจ → `userID = "anonymous"`

### 4.2 Frame ที่ Client ส่งเข้า (envelope)

| type | field ที่ใช้ | ความหมาย |
|------|-------------|----------|
| `subscribe` | `topic` | เริ่มรับข้อความจาก queue topic |
| `unsubscribe` | `topic` | เลิกรับจาก topic |
| `join_room` | `room` | เข้าห้อง (รับ broadcast ของห้อง) |
| `leave_room` | `room` | ออกจากห้อง |
| `message` | `topic` หรือ `room` + `payload` | ส่งข้อความ (publish queue + บันทึก history + broadcast) |

```json
{ "type": "subscribe", "topic": "BAACTW05/DATA" }
{ "type": "unsubscribe", "topic": "BAACTW05/DATA" }
{ "type": "join_room", "room": "BAACTW05" }
{ "type": "leave_room", "room": "BAACTW05" }
{ "type": "message", "room": "BAACTW05", "payload": { "temp": 27.5 } }
```

### 4.3 Frame ที่ Server ส่งกลับ (Broadcast)

```json
{ "event": "message", "data": {...} }
{ "event": "mqtt", "data": { "topic": "...", "payload": "..." } }
```

### 4.4 กลไก Keep-alive

- Server ส่ง **Ping** ทุก `54s` (`pingPeriod = pongWait × 9/10`)
- Client ที่เชื่อมตาม protocol จะตอบ **Pong** อัตโนมัติ (browser ทำเอง)
- ถ้าไม่มี Pong กลับภายใน `60s` (`pongWait`) → ตัดการเชื่อมต่อ
- อ่านข้อมูลได้สูงสุด `512 KB` ต่อ frame (เกิน = ตัด)

---

## 5. REST API (เฉพาะ Standalone Service)

### 5.1 `POST /api/ws/messages` — Broadcast ข้อความไปห้อง

**Request body:**

```json
{
  "room": "BAACTW01",
  "event": "message",
  "data": { "status": "running", "temp": 27.5 }
}
```

| Field | ประเภท | บังคับ | คำอธิบาย |
|-------|--------|-------|----------|
| `room` | string | ✅ | ชื่อห้องที่จะ broadcast |
| `event` | string | ✅ | ชื่อ event ใน frame ที่ client ได้รับ |
| `data` | any | ❌ | payload ที่จะแนบ (JSON ใดก็ได้) |

**Response 200:**

```json
{ "status": "sent" }
```

> หมายเหตุ: `SendMessage` เรียก `BroadcastToRoom` จากนั้น `go SaveMessage(...)` แบบ async ด้วย payload ว่าง `[]byte{}` และ sender `"api"` — มีผลเฉพาะในโหมดที่มี usecase

### 5.2 `GET /api/ws/messages?room=&limit=` — ประวัติข้อความ

| Parameter | คำอธิบาย | ค่าเริ่มต้น |
|-----------|----------|------------|
| `room` | หัวข้อ/ห้องที่ต้องการ history (ตรงกับ column `topic`) | — |
| `limit` | จำนวนสูงสุด | 20 |

**Response 200** (เรียงจากใหม่ไปเก่า `ORDER BY sent_at DESC`):

```json
[
  {
    "id": "2f0c...-uuid",
    "topic": "BAACTW01",
    "payload": { "temp": 27.5 },
    "sender_id": "anonymous",
    "sent_at": "2026-08-14T10:30:00+07:00"
  }
]
```

### 5.3 `GET /api/ws/rooms` — รายชื่อห้อง active

**Response 200:**

```json
{ "rooms": ["BAACTW01", "BAACTW05"], "count": 2 }
```

### 5.4 `GET /api/ws/rooms/{room}/stats` — จำนวน client ในห้อง

**Response 200:**

```json
{ "room": "BAACTW01", "clients": 2 }
```

---

## 6. การใช้งานจากฝั่ง Client (ตัวอย่าง)

### 6.1 JavaScript (browser)

```js
// โหมด A: ต้องส่ง token จริงจาก ws_sessions
const ws = new WebSocket("ws://localhost:8080/ws?token=<WS_TOKEN>");

// โหมด B: เชื่อมต่อกับ main API server (ไม่ต้อง token)
// const ws = new WebSocket("ws://localhost:5000/api/ws");

ws.onopen = () => {
  ws.send(JSON.stringify({ type: "join_room", room: "BAACTW05" }));
  ws.send(JSON.stringify({ type: "subscribe", topic: "BAACTW05/DATA" }));
};

ws.onmessage = (evt) => {
  const frame = JSON.parse(evt.data);
  console.log(frame.event, frame.data);
  // เช่น event "mqtt" → ข้อมูลเรียลไทม์จากอุปกรณ์
};

ws.onclose = () => console.log("disconnected");
```

### 6.2 ตรวจสอบด้วย `websocat` / `wscat`

```bash
wscat -c "ws://localhost:8080/ws?token=$WS_TOKEN"
> {"type":"join_room","room":"BAACTW01"}
< {"event":"mqtt","data":{"topic":"BAACTW01/DATA","payload":"27.5"}}
```

### 6.3 Go client (ฝั่ง API อื่นต้องการ broadcast)

```go
import "icmongolang/pkg/websocket"

// ผ่าน Broadcaster interface
var b websocket.Broadcaster = hub
b.BroadcastToRoom("BAACTW01", "alert", map[string]interface{}{"level": "high"})
b.BroadcastMessage("announcement", "system restart in 5 min") // ไปทุก client
```

---

## 7. การทำงานภายใน Hub (ดู `pkg/websocket/hub.go`)

### 7.1 lifecycle ของ client

1. **`ServeHTTP`** (delivery/ws) → auth → `upgrader.Upgrade` → `hub.AddClient(conn, userID, room)`
2. `AddClient` = `Register` (เก็บใน `clients` map) + `JoinRoom` (ถ้าระบุ room)
3. เริ่ม `go client.WritePump()` — goroutine ส่งข้อความจาก channel `Send` ไป websocket + ส่ง Ping ตามรอบ
4. `readPump` วนอ่าน frame จาก client → ประมวลผลคำสั่ง (subscribe/join/message ฯลฯ)
5. เมื่ออ่าน error (client ปิด) → `hub.Unregister` → ลบออกจาก topics/rooms/clients → ปิด channel `Send`

### 7.2 การ broadcast (3 ระดับ)

| Method | กลุ่มเป้าหมาย |
|--------|---------------|
| `BroadcastToRoom(room, event, data)` | client ที่ join room |
| `BroadcastToTopic(topic, event, data)` | client ที่ subscribe topic |
| `BroadcastMessage(event, data)` | client ทั้งหมด (global) |

ทุก broadcast สร้าง frame `{"event":...,"data":...}` แล้วส่งผ่าน `client.safeSend()`

### 7.3 `safeSend` — ไม่ block และไม่ panic

- ใช้ `sync.Mutex` + flag `closed` ป้องกัน "send on closed channel"
- ถ้า channel เต็ม (buffer 256) หรือ client ปิด → ปล่อยข้อความนั้น (drop) แล้วนับ metric `ws_messages_dropped_total`

### 7.4 Queue subscription (อัตโนมัติ)

- เมื่อมี **client แรก** subscribe topic → Hub เริ่ม goroutine `subscribeQueueTopic` → `queue.Subscribe(ctx, topic, handler)` (BRPop จาก Redis)
- เมื่อได้รับข้อความจาก queue → `BroadcastToTopic(topic, "message", payload)` ไปทุก client ที่ subscribe
- เมื่อ **client สุดท้าย** unsubscribe → ยกเลิก subscriber (`context.CancelFunc`)

---

## 8. Redis Queue (ดู `internal/modules/queue/manager.go`)

โครงสร้าง key ใน Redis:
- `queue:<topic>` — รายการข้อความ (LPUSH / BRPop)
- `delayed_queue` — sorted set สำหรับข้อความหน่วงเวลา (ประมวลผลทุก 1 วินาที)
- `dead_letter:<topic>` — ข้อความที่ retry ครบ 3 ครั้งแล้ว

| Feature | รายละเอียด |
|---------|------------|
| Publish | `LPUSH queue:<topic>` พร้อม JSON payload |
| Consume | `BRPop` บล็อค 1 วินาที + worker pool (default 10 workers) |
| Retry | Exponential backoff `2s × retries` สูงสุด 3 ครั้ง |
| Dead letter | เกิน 3 ครั้ง → `dead_letter:<topic>` |
| Delayed | `ZADD delayed_queue` แล้ว republish เมื่อครบเวลา |
| Noop mode | ตั้ง `QUEUE_TYPE=noop` → `NoopQueue` (ไม่ทำอะไร) ใช้สำหรับโหมดลดทอนเมื่อไม่มี Redis |

---

## 9. Authentication (WS Sessions)

ตาราง `ws_sessions` (จาก migration):

```sql
CREATE TABLE ws_sessions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    TEXT NOT NULL,
    token      TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);
```

Usecase ตรวจผ่าน `repo.ValidateSession(ctx, token)`:
```sql
SELECT id, user_id, token, created_at, expires_at
FROM ws_sessions
WHERE token = $1 AND expires_at > NOW();
```
- เจอ → ใช้ `user_id` เป็น sender ของข้อความ
- ไม่เจอ/หมดอายุ → ตอบ 401 `"unauthorized"`

---

## 10. Metrics (Prometheus)

โมดูล expose metrics ผ่าน `/metrics` (ของ main API server) และ `/apimetric`:

| Metric | ประเภท | ความหมาย |
|--------|--------|----------|
| `ws_connected_clients` | Gauge | จำนวน client ที่เชื่อมอยู่ตอนนี้ |
| `ws_connections_total` | Counter | จำนวน connection ทั้งหมด (สะสม) |
| `ws_active_rooms` | Gauge | จำนวนห้องที่ active |
| `ws_active_topics` | Gauge | จำนวน topic ที่มี subscriber |
| `ws_messages_broadcast_total{type}` | Counter | ข้อความที่ broadcast (room/topic/global) |
| `ws_messages_received_total{type}` | Counter | frame ที่ client ส่งเข้า (subscribe/message/etc.) |
| `ws_messages_sent_total` | Counter | ข้อความที่เขียนไป client สำเร็จ |
| `ws_messages_dropped_total` | Counter | ข้อความที่ drop (buffer เต็ม / client ปิด) |

---

## 11. ตารางข้อมูล (PostgreSQL)

### `ws_messages` — ประวัติข้อความ

| Column | Type | หมายเหตุ |
|--------|------|----------|
| `id` | UUID PK | default `gen_random_uuid()` |
| `topic` | TEXT | ห้อง/หัวข้อของข้อความ |
| `payload` | JSONB | เนื้อหาข้อความ |
| `sender_id` | TEXT | ผู้ส่ง (user_id หรือ "api"/"anonymous") |
| `sent_at` | TIMESTAMPTZ | เวลาส่ง (default NOW()) |
| `delivered_at` | TIMESTAMPTZ | เวลาส่งถึง (ยังไม่ถูกใช้ใน code) |

Index: `idx_ws_messages_topic`, `idx_ws_messages_sent_at`

---

## 12. Edge Cases / ข้อควรระวัง

1. **ใน main API server WS ไม่บันทึก history และไม่ตรวจ auth** — เพราะ `usecase = nil` (ทุกคนเป็น `anonymous`) ถ้าต้องการ history ต้องใช้ Standalone Service
2. **Room กับ Topic ต่างกัน** — `join_room` ไม่ได้ subscribe queue; `subscribe` ไม่ได้เข้าห้อง (ยกเว้น MQTT ที่ broadcast เข้า room เอง)
3. **Buffer ของ client = 256 ข้อความ** — client ที่รับช้า จะถูก drop ข้อความที่เกิน (นับใน `ws_messages_dropped_total`) ไม่ block ทั้งระบบ
4. **Frame จาก server ใช้ key `event`** แต่ frame จาก client ใช้ key `type` — อย่าสับสน
5. **`BroadcastToRoom` เจอห้องว่าง (ไม่มี client)** — จะส่งไปไม่มีใคร ไม่ error
6. **Message type ใน readPump** — ถ้าไม่มีทั้ง `topic` และ `room` → ข้าม (invalid); ถ้า usecase เป็น nil → ข้ามการประมวลผลทั้งหมด
7. **Ping** — browser ตอบ pong อัตโนมัติ; client อื่นต้อง implement pong เอง ไม่งั้นถูกตัดภายใน 60s
8. **Queue subscriber เป็นแบบ in-process** — ถ้ารันหลาย instance ต้องใช้ Redis กลาง; ข้อความที่ publish ตอนยังไม่มี subscriber จะยังค้างใน queue และถูกส่งเมื่อมี subscriber แรก (subscription เริ่มเมื่อ client subscribe เท่านั้น)
9. **`SendMessage` (REST) บันทึก payload ว่าง** — history ที่ได้จาก `POST /api/ws/messages` จะมี `payload: []` เพราะ handler ส่ง `[]byte{}` เข้า `SaveMessage`

---

## 13. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| `401 unauthorized` ต่อ `/ws` | token ไม่ตรงกับ `ws_sessions` หรือหมดอายุ (`expires_at <= NOW()`) — โหมด A เท่านั้น |
| เชื่อมได้แต่ไม่มีข้อมูล | ยังไม่ได้ `join_room` ให้ตรงชื่อห้องของ MQTT (เช่น `BAACTW05`) หรือยังไม่ได้ `subscribe` topic |
| ถูกตัด connection เป็นระยะ | client ไม่ตอบ Pong (ต้องตอบภายใน 60s); server ส่ง Ping ทุก 54s |
| ข้อความหายบางส่วน | buffer เต็ม (256) → drop; ตรวจ `ws_messages_dropped_total` หรือทำให้ client รับเร็วขึ้น |
| ข้อมูลจาก Redis queue ไม่มา | ยังไม่มี client subscribe topic นั้น (subscriber เริ่มเมื่อมีคน subscribe แรก); ตรวจ Redis `queue:<topic>` |
| `ws_messages` ว่าง | ใช้ main API server (ไม่บันทึก); ต้องใช้ Standalone Service หรือเขียนผ่าน usecase โดยตรง |
| REST `/api/ws/*` ไม่มี (404) | endpoints เหล่านี้มีเฉพาะใน Standalone Service เท่านั้น — main API server มีแค่ `/api/ws` |
