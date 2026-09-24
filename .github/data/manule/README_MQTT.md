# คู่มือการใช้งานโมดูล MQTT (Message Broker / Device Control)

โมดูล `internal/modules/mqtt` + `pkg/mqtt` เป็นระบบเชื่อมต่อ **MQTT broker** เพื่อรับ/ส่งข้อความระหว่างอุปกรณ์ IoT กับระบบ: พับลิชข้อความ, subscribe topic, อ่านข้อมูลสดจากอุปกรณ์ (request-response), สั่งงานอุปกรณ์ (device control), broadcast ข้อความไป WebSocket แบบ real-time, บันทึกข้อความลง InfluxDB และเขียน alarm log ลง PostgreSQL

สร้างด้วย library **`github.com/eclipse/paho.mqtt.golang`** (Paho MQTT client สำหรับ Go)

**สถาปัตยกรรมแบบครบวงจร:**

```
                 ┌───────────────────────────────┐
                 │   MQTT Broker (mosquitto:1883) │
                 └───────────────┬───────────────┘
                                 │ subscribe / publish
                 ┌───────────────▼───────────────┐
                 │       pkg/mqtt/Client          │  ← Paho MQTT + requestManager
                 │  (Connect / Publish / Subscribe)│
                 │  (GetDataFromTopic / RequestData)│
                 └───────────────┬───────────────┘
                                 │
                 ┌───────────────▼───────────────┐
                 │   usecase (mqttUseCase)        │
                 │  Publish / Subscribe / GetTopic│
                 │  DeviceControl / SaveAlarmLog  │
                 └───────────────┬───────────────┘
                     ┌───────────┼───────────────┐
                     ▼           ▼               ▼
              ┌────────────┐ ┌────────┐  ┌──────────────┐
              │ WebSocket  │ │InfluxDB│  │ AlarmLog (DB)│
              │Broadcast   │ │        │  │              │
              │ room=parts0│ │ mqtt_  │  │ subject/device│
              │            │ │messages│  │ status       │
              └────────────┘ └────────┘  └──────────────┘
```

**จุดสำคัญของโมดูล:**
- ได้รับการ **register ใน main API server** (`internal/server/handlers.go` — อยู่ที่ base path `/api/mqtt`) เฉพาะเมื่อ MQTT client เชื่อมต่อสำเร็จ
- แยก endpoint เป็น **2 กลุ่ม**: 4 อันแรก (GET) ใช้แค่ `RateLimit()` — **ไม่ต้อง login**; ส่วน POST `/publish`, `/subscribe`, `/unsubscribe` ต้องใช้ **JWT Bearer token** (verifier → authenticator → currentUser → activeUser)
- มี **Redis cache** ทั้งใน handler และ usecase ลดการเรียก MQTT ซ้ำ

---

## 1. ตารางไฟล์

| ไฟล์ | บทบาท |
|------|-------|
| `pkg/mqtt/client.go` | `Client` interface + `mqttClient` (Paho MQTT wrapper) + `requestManager` (ระบบ request-response + ref-count subscription) |
| `internal/modules/mqtt/usecase/usecase.go` | Business logic หลัก: `MQTTUseCase` interface, subscribe handler, device control, alarm log, save ลง InfluxDB |
| `internal/modules/mqtt/presenter/presenter.go` | DTO: `PublishRequest`, `SubscribeRequest`, `MQTTTopicDataResponse`, `DeviceControlRequest/Response` เป็นต้น |
| `internal/modules/mqtt/delivery/http/handler.go` | HTTP handlers + Redis cache (`GetTopicData`, `DeviceControl`, `Publish`, ...) |
| `internal/modules/mqtt/delivery/http/routes.go` | ลงทะเบียน route `/api/mqtt/*` (ใช้ `sync.Once` + panic guard) |
| `mqtt/mosquitto.conf` | Config ของ Mosquitto broker (Docker) |

---

## 2. การตั้งค่า Configuration

### 2.1 ตัวแปรใน `config/config.default.yml`

```yaml
mqtt:
  # broker: "mqtt://172.17.0.1:1883"
  broker: "mqtt://localhost:1883"   # ⚠️ ดูหมายเหตุเรื่อง scheme ด้านล่าง
  client_id: "icmongolang_client"
  username: ""                       # ถ้า broker ต้อง auth
  password: ""
  qos: 1
  connection_timeout: 30             # วินาที
```

### 2.2 โครงสร้าง config (`config/config.go:118`)

| field (yml) | field (Go) | ความหมาย |
|-------------|-----------|-----------|
| `broker` | `Broker` | URL broker เช่น `tcp://localhost:1883` |
| `client_id` | `ClientID` | ชื่อ client (ถ้าเว้น → อัตโนมัติ `go-client-<unixnano>`) |
| `username` / `password` | `Username` / `Password` | Credential broker |
| `qos` | `QOS` | QoS default |
| `connection_timeout` | `ConnectionTimeout` | วินาที ใช้เป็น timeout ในการเชื่อม (default 30s ถ้าเป็น 0) |

> **⚠️ หมายเหตุเรื่อง broker scheme:** ใน `internal/server/server.go` มี warning — ต้องใช้ `tcp://` (หรือ `ssl://` สำหรับ TLS) **ไม่ใช้ `mqtt://`** เพราะ Paho จะ parse scheme ไม่ตรงกับที่คาดไว้

### 2.3 Connection tuning ใน `pkg/mqtt/client.go` (`New`)

ค่า Paho options ที่ตั้งไว้:
- `CleanSession(true)` — เริ่ม session ใหม่ทุกครั้งที่ connect
- `AutoReconnect(true)` + `ConnectRetry(true)` — reconnect อัตโนมัติ, `ConnectRetryInterval=5s`, `MaxReconnectInterval=10s`
- `ResumeSubs(false)` — ไม่ resume subscription หลัง reconnect (subscription ต้องลงใหม่เอง)
- `KeepAlive=30s`, `PingTimeout=15s`, `WriteTimeout=15s`
- `TLS InsecureSkipVerify=true` — เปิดไว้สำหรับ self-signed cert ในการทดสอบภายใน

---

## 3. การเชื่อมต่อและ Wiring

### 3.1 ตอน Server เริ่มทำงาน

`internal/server/server.go`:
1. สร้าง client: `mqttClient := mqtt.New(&cfg.MQTT, log)` — ยังไม่เชื่อมต่อ
2. เรียก `connectMQTTWithRetry(mqttClient, connTimeout, 3, log)` — พยายามเชื่อม **3 ครั้ง** โดยรอ backoff `attempt × 1s` ระหว่างครั้ง
3. `internal/server/handlers.go:255-280` — ถ้า `mqttClient.IsConnected()` จริง → สร้าง usecase ด้วย `NewMQTTUseCaseWithClientAndWS(mqttClient, cfg, logger, wsHub, db, alarmLogRepo, influxDBClient, redisCache)` → สร้าง handler + ลง route
4. ถ้าเชื่อมไม่สำเร็จ → log warning `"MQTT client not connected – MQTT routes skipped"` และ **ไม่ลง route `/api/mqtt` เลย**

> นี่คือเหตุผลที่ timezone/timeout ต่างกัน: ตัว usecase จะรอ connection ต่อ (waitForConnection) แต่ถ้าตอน boot ไม่เคยเชื่อมสำเร็จ routes จะไม่มีอยู่จริง

### 3.2 ตัวสร้าง UseCase

| Constructor | ใช้เมื่อ | มีของเพิ่ม |
|-------------|---------|-----------|
| `NewMQTTUseCaseWithClientAndWS(client, cfg, log, wsBroadcaster, db, alarmLogRepo, influxClient, cache)` | main server | WebSocket broadcast + InfluxDB + alarm log + cache |
| `NewMQTTUseCaseWithClient(client, cfg, log)` | แบบง่าย | มีแค่ client + config + logger |

ทั้งสอง panic/error ถ้า client ยังไม่เชื่อมต่อ (`client.IsConnected() == false`).

---

## 4. API Reference

ทุก route อยู่ใต้ **`/api/mqtt`** และถูก wrap ด้วย `mw.RateLimit()` ทุกตัว (ทั้ง GET และ POST)

### 4.1 ตารางรวม endpoints

| Method | Path | Auth | คำอธิบาย |
|--------|------|------|----------|
| `GET` | `/api/mqtt/subscriptions` | ❌ ไม่ต้อง | รายการ topic ที่ subscribe อยู่ |
| `GET` | `/api/mqtt/status` | ❌ ไม่ต้อง | สถานะการเชื่อมต่อกับ broker |
| `GET` | `/api/mqtt/gettopicdata?topic=...` | ❌ ไม่ต้อง | อ่านข้อมูลสดจาก topic (มี cache 60s) |
| `GET` | `/api/mqtt/devicecontrol?topic=...&message=...` | ❌ ไม่ต้อง | สั่งงานอุปกรณ์ (มี cache 5s) |
| `POST` | `/api/mqtt/publish` | ✅ JWT | พับลิชข้อความไป topic |
| `POST` | `/api/mqtt/subscribe` | ✅ JWT | subscribe topic แบบถาวร (แล้ว broadcast ไป WS) |
| `POST` | `/api/mqtt/unsubscribe` | ✅ JWT | ยกเลิก subscription |

---

### 4.2 `GET /api/mqtt/status`

ตรวจว่าตอนนี้ client เชื่อมต่อกับ broker หรือไม่

**Response 200:**
```json
{
  "connected": true,
  "status": "connected",
  "timestamp": "2026-06-08 12:36:58"
}
```

- `status` เป็น `"connected"` / `"disconnected"`
- `timestamp` เป็นเวลาปัจจุบันในโซน **Asia/Bangkok** (ผ่าน `helpers.GetTimeLocation()`)

---

### 4.3 `GET /api/mqtt/subscriptions`

รายการ topic ที่ subscribe อยู่ (จำใน memory ของ usecase)

**Response 200:**
```json
{
  "topics": ["BAACTW02/DATA", "AIRCOM4/DATA"],
  "count": 2
}
```

---

### 4.4 `POST /api/mqtt/publish`

พับลิช payload ไป topic ที่กำหนด (ผ่าน `usecase.Publish`)

**Request:**
```json
{
  "topic": "BAACTW02/CONTROL",
  "qos": 1,
  "retained": false,
  "payload": "ON"
}
```

| field | type | บังคับ | หมายเหตุ |
|-------|------|-------|----------|
| `topic` | string | ✅ | ชื่อ topic |
| `qos` | byte | - | 0–2 (`omitempty,min=0,max=2`) |
| `retained` | bool | - | เก็บข้อความสุดท้ายไว้ให้ subscriber ใหม่ |
| `payload` | any | ✅ | string / []byte / object (marshal JSON อัตโนมัติ) |

**Response 200:**
```json
{ "success": true, "message": "published successfully" }
```

**Errors:** `400` ถ้าไม่มี `topic` หรือ `payload`; `500` ถ้า publish ล้มเหลว

---

### 4.5 `POST /api/mqtt/subscribe`

subscribe topic แบบถาวร (persistent) — ข้อความที่เข้ามาจะถูก:
1. **broadcast ไป WebSocket** room = segment แรกของ topic (เช่น topic `BAACTW02/DATA` → room `BAACTW02`) ด้วย event name `"mqtt"`
2. **บันทึกไป InfluxDB** (measurement `mqtt_messages`)
3. ตรวจคำว่า `alarm` ใน topic (ยังเป็น TODO)

**Request:**
```json
{ "topic": "BAACTW02/DATA", "qos": 1 }
```

**Response 200:**
```json
{
  "topic": "BAACTW02/DATA",
  "qos": 1,
  "subscribed_at": "2026-06-08T12:36:58+07:00",
  "message": "subscribed successfully"
}
```

> ถ้า subscribe topic ซ้ำ (มีอยู่แล้ว) → return nil (ไม่ error, log warning)

**รูปแบบ frame ที่ broadcast ไป WebSocket:**
```json
{
  "event": "mqtt",
  "data": {
    "topic": "BAACTW02/DATA",
    "payload": "25.5,60,1013"
  }
}
```

---

### 4.6 `POST /api/mqtt/unsubscribe`

**Request:**
```json
{ "topic": "BAACTW02/DATA" }
```

**Response 200:**
```json
{ "success": true, "topic": "BAACTW02/DATA", "message": "unsubscribed successfully" }
```

> ถ้า topic นั้นไม่ได้ subscribe อยู่ → return nil เหมือนกัน (ไม่ error)

---

### 4.7 `GET /api/mqtt/gettopicdata?topic=<topic>`

**Request-response:** subscribe topic → รอข้อความเดียว (max 10s) → unsubscribe → return ข้อมูล (แล้ว cache ไว้ 60s)

**Query params:**
| param | บังคับ | ตัวอย่าง |
|-------|--------|----------|
| `topic` | ✅ | `BAACTW05/DATA` |

**Response 200 (ข้อมูลสด `from: "mqtt"`):**
```json
{
  "statuscode": 200,
  "code": 200,
  "topic": "AIRCOM4/DATA",
  "payload": ["25.5", "60", "1013"],
  "from": "mqtt",
  "timestamp": "2026-06-08T12:36:58+07:00",
  "message": "live data",
  "mqtt_connected": true,
  "cache_enabled": true,
  "cache_hit": false,
  "data_length": 13,
  "fetch_duration_ms": 145
}
```

**Response 200 (จาก cache `from: "cache"`):** เหมือนข้างบน แต่ `from="cache"`, `cache_hit=true`, `message="from cache"`, `timestamp` = เวลาที่ cache ไว้

**พฤติกรรมพิเศษ:**
- `payload` ถูกแปลงผ่าน `processPayload()`: ถ้ามี comma → แยกเป็น `[]string`; ถ้าไม่มี → string ตรงๆ
- **Fallback** — ถ้า error เกี่ยวกับ `subscribe timeout` / `ResumeSubs` / `not currently connected` / `context deadline exceeded` → return 200 พร้อม payload fallback `"30.50,80.5,25.50,75.5"`, `from="fallback"`, `message="fallback (MQTT timeout)"` + `error_detail`
- **MQTT not connected** → return 500, `from="error"`, `message="MQTT client not connected to broker"`
- Cache key: `mqtt_topic_data:<md5(topic)>`, TTL **60 วินาที**

---

### 4.8 `GET /api/mqtt/devicecontrol?topic=<topic>&message=<message>`

**สั่งงานอุปกรณ์:** publish ข้อความไป topic (ปกติลงท้าย `CONTROL`) แล้วฟังคำตอบจาก topic `DATA` ที่สอดคล้องกัน (แทนที่ `CONTROL` → `DATA`; ถ้าไม่มีคำว่า CONTROL → ต่อท้าย `_DATA`)

**Query params:**
| param | บังคับ | ตัวอย่าง |
|-------|--------|----------|
| `topic` | ✅ | `BAACTW02/CONTROL` |
| `message` | ✅ | `ON`, `OFF`, `1`, `0`, `A1`...`G1` |

**Response 200:**
```json
{
  "statuscode": 200,
  "code": 200,
  "topic_control": "BAACTW02/CONTROL",
  "topic_data": "BAACTW02/DATA",
  "message_sent": "ON",
  "payload": "25.5,60,1013",
  "data": ["25.5", "60", "1013"],
  "status": 1,
  "status_msg": "ON",
  "timestamp": "2026-06-08T18:30:00+07:00",
  "message": "Control sent to BAACTW02/CONTROL, response received",
  "message_th": "ส่งคำสั่งไปยัง BAACTW02/CONTROL และได้รับข้อมูลตอบกลับ",
  "from": "mqtt",
  "fetch_duration_ms": 245
}
```

**ตาราง field ที่สำคัญ:**
| field | ความหมาย |
|-------|----------|
| `topic_control` | topic ที่ publish คำสั่ง |
| `topic_data` | topic ที่รอฟังข้อมูลตอบกลับ |
| `message_sent` | ข้อความที่ส่ง |
| `payload` | payload ตัวสุดท้ายจาก topic data (string) |
| `data` | payload ที่แยกด้วย comma (`strings.Split`) |
| `status` | `1` = ON, `0` = OFF |
| `status_msg` | `"ON"` / `"OFF"` |
| `from` | `"mqtt"` / `"cache"` / `"fallback"` |
| `error_detail` | มีเมื่อ fallback |

**ตรรกะ status:** message ที่ถือว่าเป็น `ON` = `1`, `on`, `a1`, `b1`, `c1`, `d1`, `e1`, `f1`, `g1` (case-insensitive) — นอกนั้นเป็น `OFF`

**พฤติกรรมพิเศษ:**
- **Cache:** key `device_ctrl:<topic>:<message>` TTL **5 วินาที** (handler ใช้ key แยก `device_control:<md5(topic:message)>` TTL 5s) — ถ้า cache hit → `from="cache"`, `cachelift` field
- **Fallback** — ถ้า waitForConnection / publish / ฟังตอบกลับ ผิดพลาด → return 200 พร้อม payload fallback `"0,0,0,0,0,0,0,0,0"`, `from="fallback"`, `message_th="fallback (MQTT หมดเวลาหรือผิดพลาด)"` + `error_detail`
- ใช้ context timeout **3 วินาที** ในการรอคำตอบจากอุปกรณ์

**ตัวอย่าง curl (จาก comment ใน `routes.go`):**
```bash
curl -H "Authorization: Bearer <token>" \
  "http://localhost:5000/api/mqtt/devicecontrol?topic=BAACTW02/CONTROL&message=ON"
```

---

## 5. การทำงานของ Pkg Client (`pkg/mqtt/client.go`)

### 5.1 `Client` interface

```go
type Client interface {
    Connect(ctx context.Context) error
    Disconnect(quiesce uint)
    Publish(topic string, qos byte, retained bool, payload interface{}) error
    Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error
    SubscribeMultiple(topics map[string]byte, callback mqtt.MessageHandler) error
    Unsubscribe(topics ...string) error
    IsConnected() bool
    RequestData(requestTopic string, responseTopic string, payload interface{}, timeout time.Duration) (interface{}, error)
    GetDataFromTopic(ctx context.Context, topic string, timeout time.Duration) ([]byte, error)
}
```

### 5.2 `GetDataFromTopic` — ขั้นตอนการอ่านข้อมูลสด

1. subscribe topic (QoS 0) ด้วย callback ที่ส่ง payload เข้า channel
2. รอ confirm จาก broker (สูงสุด 5s) — error ถ้า timeout
3. รอข้อความจริง โดย **race ระหว่าง**: context done / timeout / ได้ payload
4. `defer Unsubscribe` — ยกเลิก subscription ทุกครั้งหลังได้ข้อมูล (ไม่ว่าสำเร็จหรือไม่)

> นี่คือวิธีที่ `GetTopicData` และ `DeviceControl` ใช้ดึงข้อมูลจากอุปกรณ์

### 5.3 `requestManager` — request-response / ref-count subscription

- ลงทะเบียน global route `#` ตอน init (`client.AddRoute("#", rm.handleIncomingMessage)`)
- `getTopic(topic, timeoutMs)` — เพิ่ม pending request แล้วรอข้อความบน topic นั้น; ถ้าครบกำหนด → error `"timeout: no message from topic %s after %d ms"`
- `incrementSubscription` / `decrementSubscription` — นับจำนวน pending request ต่อ topic และทำ `Subscribe`/`Unsubscribe` จริงเมื่อนับเป็น 0↔1
- `RequestData(requestTopic, responseTopic, payload, timeout)` — publish ไป `requestTopic` (QoS 1) แล้วรอ response บน `waitTopic` (= `responseTopic` หรือ `requestTopic` ถ้าเว้น)
- response ถูก parse: ถ้า JSON → unmarshal เป็น object; ถ้าไม่ → เก็บเป็น string

---

## 6. InfluxDB Integration

ทุกข้อความที่ subscribe มา (`SaveToInflux` ใน usecase) จะถูกเขียนด้วย `WritePoint("mqtt_messages", tags, fields, time.Now())`:

```go
tags  := map[string]string{"topic": topic}
fields := map[string]interface{}{
    "payload": string(payload),
    "length":  len(payload),
}
// ถ้า payload เป็น JSON → นำ key ที่เป็นตัวเลข (float64) ไปใส่ใน fields ด้วย
```

| measurement | tags | fields |
|-------------|------|--------|
| `mqtt_messages` | `topic` | `payload` (string), `length` (int) + key ตัวเลขจาก JSON payload |

> ต้องเปิดใช้ InfluxDB client ตอน wiring (`influxDBClient != nil`) ถึงจะทำงาน; ถ้า nil → ข้ามเงื่อนไขนี้ไป (`saveToInflux` return ทันที)

---

## 7. Alarm Log

`SaveAlarmLog(alarmStatus, subject, content, deviceID, valueData)` เขียนตาราง `alarm_logs` (ผ่าน `alarmLogRepo.Create`):

| field | ค่า |
|-------|-----|
| `alarm_status` | ตามที่ส่งเข้า |
| `subject` / `content` / `device_id` | ตามที่ส่งเข้า |
| `data_alarm` | `fmt.Sprintf("%v", valueData)` |
| `date` | `2006-01-02` (เวลาท้องถิ่น) |
| `time` | `15:04:05` |
| `created_at` / `updated_at` | `time.Now()` |

> ยังมี TODO: ใน subscribe handler มี `containsAlarmKeyword(topic)` (เช็คว่ามีคำ `alarm` หรือลงท้าย `/alarm`) แต่ยังไม่ได้เรียก `helpers.AlarmDetailValidate` / `SaveAlarmLog` จริง

---

## 8. Docker

### 8.1 Mosquitto broker (`docker-compose.yml`)

```yaml
mqtt:
  image: eclipse-mosquitto:2
  ports:
    - "1885:1883"
  volumes:
    - ./mqtt/mosquitto.conf:/mosquitto/config/mosquitto.conf:ro
    - app-mqtt-data:/mosquitto/data
```

- **Host port:** `1885` → container `1883` (ใช้เฉพาะ connect จาก host ไปยัง broker ใน container เช่น `mosquitto_sub -h localhost -p 1885`; ตัว app ใน container ใช้ `mqtt://mqtt:1883`)
- Config: `listener 1883`, `allow_anonymous true`, `persistence true`
- Volume: `app-mqtt-data` เก็บข้อมูล persistence
- **Env override:** `ICMON_MQTT_BROKER` (ตั้งใน `docker-compose.yml` เป็น `mqtt://mqtt:1883`); native run ใช้ค่าใน `config.default.yml` (`mqtt://localhost:1883`)

รัน:
```bash
docker compose up -d mqtt
```

### 8.2 Prometheus exporter

- Service `mosquitto-exporter` (image `sapcc/mosquitto-exporter`) ใช้ endpoint `tcp://mqtt:1883` bind port `9234` — สำหรับกราฟใน Prometheus/Grafana

---

## 9. ตัวอย่างการใช้งานครบวงจร

### 9.1 เช็คสถานะ broker (ไม่ต้อง auth)
```bash
curl -s http://localhost:5000/api/mqtt/status
# {"connected":true,"status":"connected","timestamp":"2026-06-08 12:36:58"}
```

### 9.2 ล็อกอินเพื่อเอา token
```bash
TOKEN=$(curl -s -X POST http://localhost:5000/api/auth/login \
  -F "username=admin" -F "password=yourpassword" | jq -r .token)
```

### 9.3 Publish ข้อความ
```bash
curl -s -X POST http://localhost:5000/api/mqtt/publish \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"topic":"BAACTW02/CONTROL","qos":1,"retained":false,"payload":"ON"}'
# {"success":true,"message":"published successfully"}
```

### 9.4 Subscribe + อ่านข้อมูลผ่าน WebSocket
```bash
curl -s -X POST http://localhost:5000/api/mqtt/subscribe \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"topic":"BAACTW02/DATA","qos":1}'
# {"topic":"BAACTW02/DATA","qos":1,"subscribed_at":"...","message":"subscribed successfully"}
```

แล้วเปิด WebSocket `ws://localhost:5000/ws` แล้ว join room `BAACTW02` (ตามที่ module WebSocket กำหนด) — เมื่อมีข้อความมา จะได้รับ:
```json
{"event":"mqtt","data":{"topic":"BAACTW02/DATA","payload":"25.5,60,1013"}}
```

### 9.5 อ่านข้อมูลสดจากอุปกรณ์ (request-response)
```bash
curl -s "http://localhost:5000/api/mqtt/gettopicdata?topic=BAACTW05/DATA"
# {"statuscode":200,...,"payload":["30.5","80.5","25.5","75.5"],"from":"mqtt","message":"live data",...}
```

### 9.6 สั่งงานอุปกรณ์ (device control)
```bash
curl -s "http://localhost:5000/api/mqtt/devicecontrol?topic=BAACTW02/CONTROL&message=ON"
# {"statuscode":200,...,"topic_control":"BAACTW02/CONTROL","topic_data":"BAACTW02/DATA",
#  "message_sent":"ON","payload":"25.5,60,1013","data":["25.5","60","1013"],
#  "status":1,"status_msg":"ON","from":"mqtt","fetch_duration_ms":245}
```

### 9.7 ตัวอย่าง Go

```go
package main

import (
	"context"
	"log"
	"time"

	"icmongolang/config"
	"icmongolang/pkg/logger"
	mqtt "icmongolang/pkg/mqtt"
)

func main() {
	cfg := &config.MQTTConfig{
		Broker:   "tcp://localhost:1883",
		ClientID: "demo-client",
		Username: "",
		Password: "",
	}
	client := mqtt.New(cfg, logger.GetLogger())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	defer client.Disconnect(250)

	if err := client.Publish("BAACTW02/CONTROL", 1, false, "ON"); err != nil {
		log.Fatalf("publish failed: %v", err)
	}

	data, err := client.GetDataFromTopic(ctx, "BAACTW02/DATA", 5*time.Second)
	if err != nil {
		log.Printf("no data: %v", err)
	} else {
		log.Printf("received: %s", data)
	}
}
```

---

## 10. Edge Cases & ข้อควรระวัง

1. **`mqtt://` scheme ใช้งานไม่ได้** — ใช้ `tcp://` หรือ `ssl://` ใน `broker` (มี warning ใน `server.go`)
2. **routes ไม่ถูกสร้างถ้าเชื่อมไม่สำเร็จตอน boot** — แม้ client จะ reconnect เองได้ (Paho AutoReconnect) แต่ route `/api/mqtt/*` ถูกตัดสินครั้งเดียวตอน `NewServer`; ถ้า MQTT down ตอน start → ต้อง restart server ถึงจะได้ routes
3. **Port 1885 (host) vs 1883 (container)** — ตัว app อ่าน broker จาก config yml (`mqtt://localhost:1883` สำหรับ native run) โดย env override ใช้ชื่อ `ICMON_MQTT_BROKER` (ไม่ใช่ `MQTT_BROKER` — กัน env global ที่ชี้ port ผิดจาก docker มาชน); ส่วน `1885` เป็นแค่ host port ที่ docker เปิดไว้ให้ connect ไปยัง broker ใน container ด้วย `mosquitto_sub -h localhost -p 1885` เท่านั้น
4. **`ResumeSubs(false)`** — หลัง reconnect ระบบจะไม่ resume subscription ด้วยตัวเอง; โมดูลจะต้อง subscribe ใหม่ (ในโค้ดนี้ถ้า device เริ่ม publish ระหว่างที่ระบบ offline ข้อมูลจะหาย)
5. **Fallback data หลอก** — `GetTopicData` คืน fallback `"30.50,80.5,25.50,75.5"` (HTTP 504 + `from:"fallback"`) และ `DeviceControl` คืน `"0,0,0,0,0,0,0,0,0"` (HTTP 200) เมื่อ MQTT timeout — อย่าลืมตรวจ field `from` / `error_detail` / HTTP status ก่อนนำไปใช้จริง
6. **Cache TTL สั้น** — `gettopicdata` cache 60s, `devicecontrol` cache 5s (ทั้ง 2 layer handler+usecase); หลังสั่งงานซ้ำใน TTL จะได้ผลจาก cache ไม่ใช่สดจากอุปกรณ์
7. **GET endpoints ไม่ต้อง auth** — `subscriptions`, `status`, `gettopicdata`, `devicecontrol` เปิดให้เข้าถึงได้เฉพาะ `RateLimit` เท่านั้น
8. **Payload ที่เป็น JSON** — publish จะ marshal เป็น JSON อัตโนมัติ; ฝั่ง subscribe ถ้า payload เป็น JSON ที่มีค่า numeric จะถูกแยกเป็น field เพิ่มใน InfluxDB (`mqtt_messages`)
9. **คำสั่ง ON ครอบคลุม A1–G1** — message `a1`..`g1` ถือเป็น ON; อุปกรณ์บางตัวใช้ `A1` สลับช่อง
10. **`sync.Once`** — `MapMQTTRoutes` ลง route ได้ครั้งเดียวตลอด process (ถ้าเรียกซ้ำจะถูก ignore)

---

## 11. Troubleshooting

| อาการ | สาเหตุ / วิธีแก้ |
|-------|-----------------|
| ไม่มี route `/api/mqtt/*` (404) | MQTT เชื่อมไม่สำเร็จตอน boot → ดู log `"MQTT client not connected – MQTT routes skipped"`; เช็ค broker / `broker` scheme / port |
| `MQTT connect failed` | broker ยังไม่รัน (`docker compose up -d mqtt`) หรือ URL ผิด (ใช้ `tcp://` ไม่ใช่ `mqtt://`) หรือ auth ไม่ถูกต้อง |
| `timeout waiting for message on topic` | อุปกรณ์ไม่ได้ publish ข้อมูลบน topic นั้น; เช็คผ่าน mosquitto subscriber ว่ามี message จริง |
| ได้ fallback data (`from:"fallback"`) | MQTT subscription ไม่เสถียร / เกิน timeout 10s (gettopicdata → HTTP 504) หรือ 3s (devicecontrol → HTTP 200) |
| WebSocket ไม่ได้รับ frame `mqtt` | ยังไม่ได้ POST `/subscribe` ก่อน หรือ room ไม่ตรง (room = segment แรกของ topic) |
| cache hit ทุกครั้ง ข้อมูลไม่สด | `gettopicdata` cache 60s, `devicecontrol` cache 5s — รอให้ TTL หมดหรือ restart server |
| ข้อมูลไม่อยู่ใน InfluxDB | ต้องส่ง `influxDBClient` ตอน wiring (`NewMQTTUseCaseWithClientAndWS`); ถ้า nil จะข้าม |

---

## 12. โครงสร้างข้อมูลที่เกี่ยวข้อง

- **InfluxDB measurement `mqtt_messages`** — ข้อความ MQTT ที่ subscribe มา (tags: `topic`, fields: `payload`, `length`, + key ตัวเลขจาก JSON)
- **PostgreSQL `alarm_logs`** — ผ่าน `SaveAlarmLog` (ยังมี TODO ใน subscribe handler)
- **Redis** — cache key `mqtt_topic_data:<md5(topic)>` (60s), `device_ctrl:<topic>:<message>` (5s), `device_control:<md5(topic:message)>` (5s)
- **WebSocket** — room จาก segment แรกของ topic, event name `"mqtt"`, data `{topic, payload}`
