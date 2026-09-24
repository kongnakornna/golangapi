# คู่มือการใช้งาน WebSocket + REST API

## 📌 บทนำ

ระบบนี้ให้บริการ **WebSocket** สำหรับการส่งข้อความแบบ real-time และ **REST API** สำหรับการจัดการข้อความ ห้อง และสถานะ โดยทำงานบนพอร์ต `8081` (สามารถปรับได้ใน config)

---

## 🚀 การรัน WebSocket Server

```bash
# รันด้วย Air (hot-reload)
air

# หรือรันด้วย go run
go run cmd/websocket/main.go

# หรือ build แล้วรัน
go build -o websocket-server.exe cmd/websocket/main.go
./websocket-server.exe
```

เมื่อรันสำเร็จจะเห็นข้อความ:
```
WebSocket + REST server listening on :8081
```

---

## 🔗 REST API Endpoints

ฐาน URL: `http://localhost:8081/api/ws`

| Method | Endpoint | คำอธิบาย |
|--------|----------|----------|
| POST | `/messages` | ส่งข้อความไปยังห้อง (room) |
| GET | `/messages` | ดึงประวัติข้อความในห้อง |
| GET | `/rooms` | รายการห้องที่มีผู้ใช้งาน |
| GET | `/rooms/{room}/stats` | สถิติจำนวนผู้ใช้ในห้อง |

---

### 1. ส่งข้อความไปยังห้อง

**POST** `/api/ws/messages`

**Request Body (JSON):**
```json
{
  "room": "sensor",
  "event": "control",
  "data": {
    "command": "ON",
    "device": "fan_01"
  }
}
```

**Response:**
```json
{
  "status": "sent"
}
```

> ข้อความนี้จะถูก broadcast ไปยังทุก client ที่อยู่ในห้อง `sensor` ทันที

---

### 2. ดึงประวัติข้อความในห้อง

**GET** `/api/ws/messages?room=sensor&limit=10&offset=0`

**Query Parameters:**
- `room` (required): ชื่อห้อง
- `limit` (optional, default=20): จำนวนข้อความที่ต้องการ
- `offset` (optional, default=0): เริ่มต้นที่ index ไหน

**Response:**
```json
[
  {
    "id": 1,
    "room": "sensor",
    "event": "mqtt",
    "data": {
      "topic": "sensor/temperature",
      "payload": "25.5"
    },
    "sender": "api",
    "created_at": "2026-06-20T10:00:00Z"
  }
]
```

---

### 3. รายการห้องที่มีผู้ใช้งาน

**GET** `/api/ws/rooms`

**Response:**
```json
{
  "rooms": ["sensor", "alarm", "control"],
  "count": 3
}
```

---

### 4. สถิติจำนวนผู้ใช้ในห้อง

**GET** `/api/ws/rooms/{room}/stats`

**ตัวอย่าง:** `/api/ws/rooms/sensor/stats`

**Response:**
```json
{
  "room": "sensor",
  "clients": 5
}
```

---

## 🔌 WebSocket Connection

### URL สำหรับเชื่อมต่อ

```
ws://localhost:8081/ws?user_id={user_id}&room={room}
```

**Query Parameters:**
- `user_id` (optional): ระบุตัวตนของผู้ใช้ (ค่าเริ่มต้น `anonymous`)
- `room` (optional): ห้องเริ่มต้นที่ต้องการเข้าร่วม (ค่าเริ่มต้น `""`)

**ตัวอย่าง:**
```
ws://localhost:8081/ws?user_id=device_01&room=sensor
```

---

### การส่งข้อความจาก Client ไปยัง Server

Client สามารถส่งข้อความในรูปแบบ JSON ไปยัง Server ได้:

#### 1. เข้าร่วมห้อง (Join Room)

```json
{
  "event": "join",
  "room": "alarm"
}
```

หลังจากส่งข้อความนี้ Client จะย้ายไปอยู่ในห้อง `alarm` และจะได้รับข้อความทั้งหมดที่ broadcast ในห้องนี้

#### 2. ส่งข้อความทั่วไป

```json
{
  "event": "ping",
  "data": "hello from client"
}
```

หากไม่ระบุ `room` ระบบจะใช้ห้องปัจจุบันของ Client อัตโนมัติ

---

### การรับข้อความจาก Server

Server จะส่งข้อความในรูปแบบ:

```json
{
  "room": "sensor",
  "event": "mqtt",
  "data": {
    "topic": "sensor/temperature",
    "payload": "25.5"
  },
  "sender": "mqtt",
  "timestamp": 1718875200
}
```

หรือเมื่อมีการส่งผ่าน REST API:

```json
{
  "room": "sensor",
  "event": "control",
  "data": {
    "command": "ON",
    "device": "fan_01"
  },
  "sender": "api",
  "timestamp": 1718875300
}
```

---

## 🔄 การทำงานร่วมกับ MQTT

เมื่อ MQTT client ได้รับข้อมูลจาก topic ใด ๆ ระบบจะ:

1. แยก `room` จากชื่อ topic (ใช้ส่วนแรกของ topic เช่น `sensor/temperature` → room = `sensor`)
2. Broadcast ข้อมูลไปยัง WebSocket clients ที่อยู่ในห้องนั้น
3. บันทึกข้อมูลลง InfluxDB (ถ้ามี)

**รูปแบบข้อความที่ได้รับจาก WebSocket เมื่อมี MQTT data:**

```json
{
  "room": "sensor",
  "event": "mqtt",
  "data": {
    "topic": "sensor/temperature",
    "payload": "25.5"
  },
  "sender": "mqtt",
  "timestamp": 1718875400
}
```

---

## 🛠️ การทดสอบด้วย Tools

### 1. ทดสอบ WebSocket ด้วย `wscat` (command line)

ติดตั้ง:
```bash
npm install -g wscat
```

เชื่อมต่อ:
```bash
wscat -c "ws://localhost:8081/ws?user_id=test&room=sensor"
```

ส่งข้อความ join room:
```
{"event":"join","room":"sensor"}
```

ส่งข้อความ:
```
{"event":"ping","data":"Hello"}
```

---

### 2. ทดสอบ WebSocket ด้วย Postman

1. เปิด Postman → New → WebSocket Request
2. ใส่ URL: `ws://localhost:8081/ws?user_id=test&room=sensor`
3. กด Connect
4. ไปที่ tab "Message" พิมพ์ JSON แล้วกด Send

---

### 3. ทดสอบ REST API ด้วย `curl`

**ส่งข้อความไปยังห้อง:**
```bash
curl -X POST http://localhost:8081/api/ws/messages \
  -H "Content-Type: application/json" \
  -d '{"room":"sensor","event":"test","data":"hello"}'
```

**ดึงประวัติข้อความ:**
```bash
curl "http://localhost:8081/api/ws/messages?room=sensor&limit=5"
```

**ดูรายการห้อง:**
```bash
curl http://localhost:8081/api/ws/rooms
```

**ดูสถิติห้อง:**
```bash
curl http://localhost:8081/api/ws/rooms/sensor/stats
```

---

## 📝 สรุป

| การทำงาน | วิธีใช้ |
|----------|--------|
| เชื่อมต่อ WebSocket | `ws://localhost:8081/ws?user_id=xxx&room=xxx` |
| ส่งข้อความผ่าน WebSocket | ส่ง JSON `{"event":"...","data":...}` |
| เข้าร่วมห้องผ่าน WebSocket | ส่ง `{"event":"join","room":"..."}` |
| ส่งข้อความผ่าน REST API | `POST /api/ws/messages` |
| ดูประวัติข้อความ | `GET /api/ws/messages?room=...` |
| ดูรายการห้อง | `GET /api/ws/rooms` |
| ดูจำนวนคนในห้อง | `GET /api/ws/rooms/{room}/stats` |

---

## ❗ ข้อควรระวัง

- WebSocket และ REST API ใช้พอร์ตเดียวกัน (`8081`)
- หากต้องการให้ WebSocket บันทึกข้อความลง Database ต้องแน่ใจว่า `wsUsecase` ถูกส่งเข้าไปใน `wsdelivery.NewWsHandler` แทนค่า `nil`
- ข้อมูล MQTT จะถูก broadcast ไปยัง room ที่แยกจาก topic โดยอัตโนมัติ (ใช้ส่วนแรกของ topic)
- ถ้าไม่มี room ใน query parameter เมื่อเชื่อมต่อ Client จะอยู่ในห้องว่าง (`""`) และจะไม่ได้รับข้อความใด ๆ จนกว่าจะ `join` ห้อง

---

✨ **พร้อมใช้งานแล้ว!**