# Plan: ประยุกต์ InfluxDBServiceModule(V1+V2) เข้ากับ icmongolang เดิม — ขยายแพลตฟอร์ม IoT 7 ระบบ

> ส่วนนี้เป็นการวางแผน **เอกสารแผนงาน (Plan)** ก่อนเริ่มเขียนโค้ด ตามที่ user ร้องขอ ("ให้ ทำเอกสาร แผนงานออกมาก่อน") ยังไม่มีการแก้โค้ดใดๆ

---

## 0. วัตถุประสงค์และขอบเขต

เอกสาร `docs/InfluxDBServiceModule.md` (DDD + Clean Architecture สำหรับ InfluxDB module) และ
`docs/InfluxDBServiceModuleV2.md` (ขยายเป็น IoT Stack เต็ม: MQTT, Socket.IO, Redis, PostgreSQL, LLM/AI Agent)
จะถูก **ประยุกต์/ปรับใช้ออกแบบ** กับโปรเจกต์ `icmongolang` ที่มีอยู่แล้ว โดย **ไม่ทำให้การทำงานเดิมเสียหาย**
(ทุก module เดิม: auth, users, iot, mqtt, alarm, dashboard, settings, ... ต้องยัง build / ทำงานเหมือนเดิม)

โปรเจกต์เดิมมีกรอบที่ต่างจากเอกสารทั้งสองในหลาย ๆ จุด เราจะ **map แนวคิด** ของเอกสารเข้ากับกรอบเดิม ไม่ใช่คัดลอกโครงสร้างไฟล์จากเอกสารแบบตรงตัว:

| ด้าน | เอกสาร V1/V2 (สมมติ) | icmongolang เดิม (กรอบจริง) | แนวทางประยุกต์ |
|------|---------------------|---------------------------|----------------|
| Architecture | DDD + Clean (core/application/command/query + ports/adapters) | Clean 3-layer + Delivery: `models` / `repository` / `usecase` / `delivery/http` ต่อ module | ใช้กรอบ module เดิมเป็นหลัก (เดินตาม `internal/modules/*`) |
| DI | Google Wire (`internal/di/wire.go`) | manual wiring ที่ `internal/server/handlers.go` | ต่อ Dependency ผ่าน constructor ใน `handlers.go` (ไม่เพิ่ม Wire) |
| InfluxDB write/query | `internal/adapters/out/influxdb` | `pkg/influxdb` (client wrapper) + `internal/modules/influxdb` | reuse `pkg/influxdb`/module ที่มีอยู่ |
| MQTT subscribe/publish | `internal/adapters/in/mqtt` | `pkg/mqtt` + `internal/modules/mqtt` + `iot.StartIngest` (subscribe `#`) | reuse ของเดิม |
| WebSocket/Socket.IO | go-socket.io server | `pkg/websocket` + `internal/modules/websocket` (Hub/Client) | reuse ของเดิม |
| Alarm rules | `Alarm Engine` (domain) | `internal/modules/alarm/processor` + `pkg/helpers/iot.go` | reuse + ขยายการ dispatch |
| Notification (email/sms/line/...) | — (V2 แค่ Alarm engine) | มี **Settings module** (`SdIotEmail/Sms/Line/Telegram/Nodered`) + `internal/modules/email` | reuse และเชื่อม dispatch |
| Auth | JWT (ในเอกสารมี) | JWT mw: `Verifier/Authenticator/CurrentUser/ActiveUser` | reuse ของเดิม |

### เป้าหมาย 7 ระบบ (ตาม user)

1. **ระบบ Settings** — ส่วนใหญ่ **ทำเสร็จแล้ว** ใน `internal/modules/settings` (ดู `docs/plan/2026-08-22-settings-module.md`) → ต้อง **ตรวจ/ปิดช่องว่าง** ไม่สร้างซ้ำ
2. **ระบบแสดงข้อมูล Real-time Dashboard Monitoring** — reuse `iot` endpoint (`/monitordevicegroup`, `/monitordevicechart`) + `pkg/websocket` push → เสริม real-time + aggregate query จาก InfluxDB
3. **ระบบ Alarm — email / sms / line / discord / io device บน dashboard** — ขยาย alarm dispatch ให้เชื่อม `Settings` (Email/Sms/Line) + **เพิ่ม Discord** + **IO device control** + push บน dashboard ผ่าน websocket
4. **ระบบ Management — Auto control / สั่งงาน / ตั้งค่าอุปกรณ์** — reuse `DeviceControl`/`controls` + Settings schedule → เสริม auto-control rule (V2 AI agent + alarm event control)
5. **ระบบ Logging แบบทั่วไป + Smart Contract Blockchain** — logging ทั่วไป (reuse logger/zap + DB log) + **เพิ่ม blockchain audit log module** (แบบ new)
6. **ระบบ API Management** — reuse Settings `SdIotApi` + Swagger + add API key/rate limit (Redis) ตาม V2 section 5
7. **ระบบ Diagram workflow เหมือน Node-RED** — reuse Settings `SdIotNodered` + **เพิ่ม flow-engine** (node/edge graph พร้อม runtime) ใหม่

> ข้อกำหนดสำคัญ: แต่ละระบบ = module ใหม่/ขยาย อยู่ใต้ `internal/modules/*` ใช้กรอบเดิม สร้างผ่าน constructor ใน `internal/server/handlers.go` วงจร build/test ต้องเขียวทุก phase

---

## 1. ภาพรวมโมดูลที่เกี่ยวข้อง (ปัจจุบัน ↔ เพิ่ม/ขยาย)

| Module | สถานะ | บทบาทใน 7 ระบบ |
|--------|-------|----------------|
| `internal/modules/settings` | ✅ มีอยู่ | ระบบ 1/3/4/6/7 (data: email/sms/line/telegram/nodered/api/device/schedule/alarm...) |
| `internal/modules/iot` | ✅ มีอยู่ | ระบบ 2/3/4 (monitor, control, ingester) |
| `internal/modules/influxdb` + `pkg/influxdb` | ✅ มีอยู่ | ระบบ 2 (time-series read) + aggregate |
| `internal/modules/mqtt` + `pkg/mqtt` | ✅ มีอยู่ | ระบบ 3/4 (publish control, receive data) |
| `internal/modules/websocket` + `pkg/websocket` | ✅ มีอยู่ | ระบบ 2/3 (real-time push dashboard) |
| `internal/modules/alarm` + `processor` | ✅ มีอยู่ | ระบบ 3 (rule engine) |
| `internal/modules/email` | ✅ มีอยู่ | ระบบ 3 (email dispatch) |
| `pkg/llm` | ✅ มีอยู่ | ระบบ 4 (AI agent assist — reuse V2) |
| `internal/modules/notifier` | 🔶 **ใหม่** | ระบบ 3 — dispatcher: email/sms/line/discord/io-device |
| `internal/modules/control` (หรือขยาย iot) | 🔶 **ใหม่/ขยาย** | ระบบ 4 — auto-control manager |
| `internal/modules/auditlog` | 🔶 **ใหม่** | ระบบ 5 — general log + blockchain anchor |
| `internal/modules/apimanager` | 🔶 **ใหม่** | ระบบ 6 — API key/rate limit/usage |
| `internal/modules/flowengine` | 🔶 **ใหม่** | ระบบ 7 — Node-RED-like workflow diagram + runtime |

> 🔶 = usually put into a dedicated new package directory following the same pattern.

---

## 2. สถาปัตยกรรมรวมที่ประยุกต์ (จาก V2 diagram → กรอบเดิม)

```
IoT Devices/Sensors ──MQTT──▶ MQTT Broker (pkg/mqtt + module/mqtt)
                                   │
      ┌────────────────────────────▼────────────────────────────┐
      │  iot.StartIngest (subscribe "#")  +  module/mqtt handler  │  (ofเดิม)
      └────────────────────────────┬────────────────────────────┘
        write                       │
        ▼                           ▼
   pkg/influxdb (time-series)   Redis (cache/rate-limit/stream)     (ของเดิม + V2)
        │                           │
        ▼                           ▼
   Alarm Rule Engine (module/alarm/processor)                      (ของเดิม)
        │
        ├─▶ notifier (NEW) ─ email / sms / line / discord / io-device ─▶ Settings config
        ├─▶ wsHub push ─▶ dashboard (real-time)                        (ของเดิม websocket)
        └─▶ control (NEW) auto-control ─▶ MQTT publish ─▶ device      (ระบบ 4)

   Settings module ──▶ API Manager (NEW ระบบ 6) ──▶ Audit log + Blockchain (NEW ระบบ 5)
   Flow Engine (NEW ระบบ 7) ──▶ trigger node/edge ──▶ notifier/control/mqtt
```

กฎสำคัญ: **ทุกทางเพิ่มใหม่ไม่แตะ interface/พฤติกรรมที่ของเดิมใช้อยู่** — เพิ่ม constructor param โดยให้ nil-safe / optional เพื่อไม่มีผลกับเส้นทางเดิม (ดู Section 8 backward-compat)

---

## 3. ระบบ 1 — Settings

**สถานะ:** ส่วนใหญ่เสร็จจาก `2026-08-22-settings-module.md` (endpoint ~200 ตัว, envelope, JWT, generic ListPaginate) แล้ว ดู `internal/modules/settings/*`

**Action (ตรวจ + ปิดช่องว่าง):**
1. ผูก settings repo เข้ากับ **notifier/control** (ระบบ 3/4) เพื่ออ่าน config email/sms/line/nodered/device/token
2. เพิ่ม Discord config ถ้าจำเป็น: ตรวจ `SdIotTelegram` รูปแบบใกล้เคียง → ต่อจาก Settings แบบ `SdIotDiscord` (table `sd_iot_discord`) ถ้ายังไม่มี — **เป็น optional**; ถ้าไม่ต้องการตารางใหม่ เก็บ discord webhook ใน `SdIotApi`/`SdIotHost` ที่มีอยู่แล้ว
3. ไม่สร้าง endpoint ซ้ำ

---

## 4. ระบบ 2 — Real-time Dashboard Monitoring

**ของเดิม:** `iot` มี `/monitordevicegroup`, `/monitordevicechart` แล้ว; `pkg/websocket` มี Hub/Broadcaster

**เพิ่ม:**
1. **Realtime push**: ให้ ingester/`ProcessMqttData` หลัง write เรียก `wsHub.Broadcast("monitor", payload)` (ของเดิม wsHub รับ `Broadcaster` interface แล้วใน mqtt module → ทำแบบเดียวกัน)
2. **Aggregate query ใหม่** ใน `internal/modules/influxdb` (หรือ iot): ย้ายแนวคิด `QueryAggregate` จาก V1 (mean/sum/max/min/count, group by time/tag) → endpoint `GET /api/influxdb/aggregate` หรือ `/api/iot/...`
3. **Read-time metric**: reuse `pkg/influxdb` read path; ตรวจ Redis cache TTL ตาม V2 (reuse `redisCache.NewCache`)

---

## 5. ระบบ 3 — Alarm: email / sms / line / discord / io-device บน dashboard

**ของเดิม:** `module/alarm/processor` (rule), `module/alarm` (alarm device/event), `module/mqtt` publish control, `SdIotEmail/Sms/Line/Telegram` config ใน settings

**เพิ่ม — `internal/modules/notifier` (ใหม่):**
```
Notifier Service
  └─ Dispatch(ctx, AlarmResult, device) 
       ├─ EmailChannel      → net/smtp (จาก SdIotEmail config / env)
       ├─ SmsChannel        → SMS gateway (SdIotSms config / token)
       ├─ LineChannel       → LINE Notify/Messaging API (SdIotLine token)
       ├─ DiscordChannel    → Discord webhook (NEW; config จาก settings api/host)
       └─ IODeviceChannel   → publish MQTT control (mqtt_control_on/off) ตาม processor.EventControl
```
- hook เข้ากับ alarm trigger path: หลัง `processor.AlarmDetailValidate` ได้ `AlarmResult.Status != 5` → `Notifier.Dispatch`
- **Push ขึ้น dashboard**: ผล dispatch + alarm → `wsHub.Broadcast("alarm", ...)` (ระบบ 2 websocket)
- เก็บประวัติ: reuse `SdAlarmProcessLog{Email,Line,Sms,Telegram,...}` ของ settings (ดู models) หรือตารางใหม่ `notifier_log`

**ของเดิมหา hook ที่จุดไหน:** ต้อง trace จุดที่ alarm engine ทำงานจริง (อาจอยู่ใน `iot/usecase` หรือ `alarm`/`processor`) — เปิดจุดนั้นแล้วเรียก notifier + ws push ด้วย **nil-safe** (Section 8)

---

## 6. ระบบ 4 — Management: Auto control / สั่งงาน / ตั้งค่าอุปกรณ์

**ของเดิม:** `DeviceControl`/`controls` (iot) + `DeviceActionUser`, schedule ใน settings (`sd_iot_schedule`/`SdIotScheduleDevice`)

**เพิ่ม — Auto Control Manager (ใหม่หรือขยาย `iot`):**
1. **auto-control rule**: จาก `AlarmResult.EventControl` + `messageMqttControl` (ที่มีใน processor) → publish MQTT เพื่อเปิด/ปิดอุปกรณ์อัตโนมัติ (ต่างจากของเดิมที่ส่งแค่ event เปิด) — ใช้งานร่วมกับ notifier.IODeviceChannel
2. **Manual/schedule control**: reuse settings schedule + `DeviceActionUser`; เพิ่ม endpoint execute control ทันที (มี `/controls` แล้ว)
3. **ตั้งค่าอุปกรณ์**: reuse `updatedeviceconfig`/`devicestatus` (iot) + settings `update` device; ไม่สร้างซ้ำ
4. (optional) **AI assist** reuse `pkg/llm` ตาม V2 — แนะนำ action เมื่อ alarm; เป็น async ไม่ block

---

## 7. ระบบ 5 — Logging ทั่วไป + Smart Contract Blockchain

### 5a. Logging ทั่วไป
- **Logging ระบบ**: มี `pkg/logger` (zap) แล้ว — ใช้ต่อ ไม่แก้
- **Log การกระทำ (audit)**: เพิ่ม **`internal/modules/auditlog`** (ใหม่): บันทึก action ของ admin/ผู้ใช้ (create/update/delete/settings/control) ลงตาราง `audit_logs` (user, action, resource, before/after, timestamp, request_id)

### 5b. Smart Contract Blockchain (เพิ่มใหม่ จะไม่ทำ offline โดย default — design)
- **วิธีเจ้าใจปลอดภัย-ไม่ break**: ใช้ **blockchain anchoring** แบบ off-chain + hash chain:
  - ทุก audit log มี `hash = SHA256(previousHash + payload + timestamp + nonce)`
  - เก็บ `previousHash` ต่อกันเป็น **hash chain** (verifiable, tamper-evident)
  - optional: anchor root hash ไป external blockchain (e.g. Ethereum/other) ผ่าน channel/provider interface — default เป็น no-op provider ดังนั้น **ไม่มีผลต่อการทำงานเดิม**
- **Model**: `BlockchainEntry{ id, chain_id, payload, hash, prev_hash, anchored_at, tx_hash (nullable) }`
- **Endpoints**: `POST /api/auditlog`, `GET /api/blockchain/verify/{id}` (ตรวจ hash chain), `GET /api/blockchain/chain`
- ไม่แนะนำให้เขียน smart contract จริงบน chain ในการ iteration แรก (ค่าใช้จ่าย/complexity) — ทำ hash-chain anchor ก่อน แล้ว tie เข้า external ledger ภายหลังผ่าน interface

---

## 8. ระบบ 6 — API Management

**ของเดิม:** Settings มี `SdIotApi` (api_name, host, port, token_value) + Swagger (`/swagger/`), JWT middleware, middleware มี `MonitoringMiddleware`/metrics

**เพิ่ม — `internal/modules/apimanager` (ใหม่):**
1. **API Key management**: CRUD ต่อ `SdIotApi` + สร้าง key ปิด `token_value`; middleware `APIKeyAuth` (อ่าน `X-API-Key`) → ตรวจ Redis `pkg/db/redis` (ตาม V2 `CheckRateLimit`)
2. **Rate limit**: reuse `pkg/db/redis` rate limit (V2 section 5.1) middleware ต่อ API-key/route
3. **Usage log**: บันทึก request count/latency/status ต่อ key → dashboard ระบบ 2
4. **Doc auto**: reuse Swagger; เพิ่ม `@securityDefinitions.apikey APIKeyAuth` ที่ `main.go` (optional)
- เปิด middleware เป็น **optional per-route** เพื่อไม่ให้ API ของ module เดิมถูกบังคับโดยไม่ตั้งใจ (Section 8)

---

## 9. ระบบ 7 — Diagram Workflow เหมือน Node-RED

**ของเดิม:** Settings มี `SdIotNodered` (nodered_name, host, port, routing, client_id, grant_type, scope, username, password) — เป็น config การ connect ไป Node-RED ภายนอก

**เพิ่ม — `internal/modules/flowengine` (ใหม่):** สร้าง flow editor + runtime ภายใน (Node-RED-like):
1. **Flow model**: `Flow{ id, name, nodes[jsonb], edges[jsonb], enabled }` — node = trigger/action (mqtt-in, alarm, delay, email/line/discord/io-device, control, http), edge = wire ต่อ node
2. **CRUD + validate**: validate graph (cycle, ไร้ dangling edge), save เป็น jsonb
3. **Runtime**: รับ event (mqtt message / alarm / schedule) → ท่องกราฟ (BFS/topological) execute node → เรียก notifier/control/mqtt
4. **Visual**: ฝั่ง frontend ดึง `/api/flowengine/flow/:id` ส่ง nodes/edges ให้วาด diagram (Node-RED-like) — backend จัด graph schema ให้
5. optional: proxy/call ไป Node-RED จริงผ่าน `SdIotNodered` config (reuse settings)

---

## 10. Backward-Compatibility & ไม่กระทบของเดิม (ข้อบังคับ)

1. **ทุก module ใหม่** อยู่ใน `internal/modules/*` — ไม่แก้ไฟล์ของ module เดิม (ยกเว้น wiring กลาง)
2. **Wiring ที่ `internal/server/handlers.go`**: เพิ่ม constructor/mount **ต่อท้าย** บรรทัดเดิม; ของใหม่ต้อง nil-safe — ถ้า dependency (mqttClient/influxClient/wsHub) เป็น nil ให้ **skip + log warn** เหมือน pattern ที่มีอยู่แล้วในไฟล์ (เช่น block ES/VectorData)
3. **ไม่ลบ/ไม่เปลี่ยน signature ซ้าย** ของ constructor เก่า; ถ้าต้องส่งค่าใหม่ ให้เพิ่ม param ตัวท้าย หรือสร้างตัวสร้างใหม่
4. **ไม่เปลี่ยน interface `pkg/websocket` Hub/Broadcaster** เดิม
5. **Migration**: ถ้าต้องเพิ่งตารางใหม่ (audit_logs, blockchain_entry, flow, flow_node/log, notifier_log, sd_iot_discord) — เข้า migration list ของโปรเจกต์ ตาม pattern ที่ใช้ (`cmd/migrate.go` + models) เป็น additive เท่านั้น (ไม่มี drop/alter ของเดิม)
6. **Config**: ถ้าต้อง env ใหม่ เพิ่มลง `config/config.go` struct + `config.default.yml` เป็น optional default — ไม่เปลี่ยน env เดิม
7. **ทุก phase**: `go build ./...` + `go vet ./...` + `go test ./...` ต้องผ่าน (ถือเป็นเกท) — และ smoke endpoint เดิมยัง work
8. **Redis key**: ตั้ง prefix ใหม่ (เช่น `notifier:`, `flow:`, `apikey:`, `audit:`) ไม่ชน key ของเดิม (`device:`, `llm:`, `get_device_data_ALL`)

---

## 11. Roadmap / Phases (build green ทุก phase)

| Phase | รายการ | ระบบ |
|-------|--------|------|
| 0 | Skeleton: สร้าง dir/code ของ module ใหม่ + wiring nil-safe + stub (return ErrNotImplemented) → build green | 3,4,5,6,7 |
| 1 | ตรวจ Settings (ระบบ 1) — scan ช่องว่าง, ผูก config | 1 |
| 2 | Aggregate InfluxDB + realtime ws push (dashboard) | 2 |
| 3 | Notifier: email/sms/line/discord/io-device + hook จุด trigger alarm + ws push | 3 |
| 4 | Auto-control manager + schedule execute + (optional AI assist) | 4 |
| 5 | AuditLog + Blockchain hash-chain anchor | 5 |
| 6 | API Manager: key/rate-limit/usage + optional middleware | 6 |
| 7 | Flow Engine: model/CRUD/validate/runtime + wire to notifier/control | 7 |
| 8 | รวม wiring, docs, verify ทั้งหมด | ทั้งหมด |

ทุก phase จบด้วย `go build ./... && go vet ./...` + `go test ./internal/modules/<new>/...`

---

## 12. Unit Tests (ตามกรอบ tdd — decision: เขียน)

| Module | Test |
|--------|------|
| `notifier` | email/sms/line/discord formatter; io-device → mqtt payload mapping; nil-config fallback |
| `auditlog` | hash-chain: hash เปลี่ยนเมื่อ payload เปลี่ยน; prev_hash ต่อเนื่อง; verify |
| `blockchain` | anchor provider stub; tx_hash null path |
| `apimanager` | key generate/validate; rate-limit middleware (Redis mock ผ่าน miniredis มีใน go.mod) |
| `flowengine` | validate: cycle/dangling; BFS execute order; node action dispatch |
| `influxdb` (aggregate) | flux builder (mean/sum/group by) — ตาม V1 |

---

## 13. Verification

1. `go build ./...`, `go vet ./...`, `go test ./...` ผ่านทุก phase
2. Run server → smoke endpoint เดิม (login, `/api/iot/monitordevicegroup`, `/api/settings/listsetting`) ยัง return เดิม
3. Smoke ของใหม่: publish MQTT data → ดู alarm trigger → notifier log + ws push; create flow → trigger → execute; register api-key → เรียกด้วย `X-API-Key`
4. ตรวจตารางใหม่ถูก migrate (additive)

---

## 14. สรุป

เอกสารทั้งสองถูกประยุกต์เป็น **แผนงาน 7 ระบบ** บนกรอบเดิมของ icmongolang — ระบบ 1 ทำเสร็จแล้ว (Settings) เหลือตรวจ/ผูก; ระบบ 2 ใช้ของเดิม + เสริม aggregate/ws; ระบบ 3/4 ขยาย alarm→notifier/control; ระบบ 5-7 เป็น module ใหม่ (auditlog+blockchain, apimanager, flowengine) ทั้งหมดเดินตาม convention เดิมและไม่กระทบเส้นทางเดิม (nil-safe wiring + build green ทุก phase)

---

## 15. บันทึกความคืบหน้า (Progress Log)

> อัปเดตล่าสุด: 2026-08-28 — เซสชันนี้เสร็จ Phase 0, 2, 3, 4, 5, 6, 7, 8 แล้ว (แผนครบถ้วน)

### เสร็จแล้ว ✅

**Phase 0 — Skeleton + wiring (4 module ใหม่)**
- `internal/modules/notifier/` — dispatch email/sms/line/discord/io (delivery/http + usecase + presenter)
- `internal/modules/auditlog/` — audit + blockchain hash-chain (skeleton)
- `internal/modules/apimanager/` — API key/rate-limit/usage (skeleton)
- `internal/modules/flowengine/` — Node-RED-like flow (skeleton)
- Wire nil-safe ใน `internal/server/handlers.go` (block "15b") ≈ line 357+
- **หมายเหตุ:** restore `docs/docs.go` จาก git — ไฟล์ swagger-gen ถูกเขียนทับด้วย markdown (ของเดิมเสีย ไม่ใช่จากงานนี้)

**Phase 2 — Realtime dashboard push**
- `internal/modules/realtime/` (ใหม่): `POST /api/realtime/publish`, `GET /api/realtime/health`
- usecase รับ `pkg/websocket.Broadcaster` (wsHub) nil-safe; broadcast ไป room หรือ global
- unit tests ผ่าน (room/global/nil-broadcaster/invalid)
- Aggregate InfluxDB มีอยู่แล้ว (`pkg/influxdb.CalculateStatistics` + `/api/influx/statistics`) ไม่ต้องสร้างซ้ำ

**Phase 3 — Notifier จริง**
- `internal/modules/notifier/repository/pg_repository.go` — อ่าน channel config จาก DB (`sd_iot_email/sms/line/nodered`) + `Repository` interface
- `internal/modules/notifier/provider/provider.go` — 5 channel จริง:
  - `email` → reuse `pkg/sendEmail` (SMTP) — เปิดเมื่อ `appCfg != nil`
  - `sms` → HTTP SMS gateway (host/apikey/originator)
  - `line` → LINE Notify API (access token)
  - `discord` → Discord webhook POST
  - `io` → publish MQTT control (`mqtt.Client`)
- usecase: ดึง config → send ผ่าน provider → คืนผล; nil-safe ทุกจุด
- wiring: `NewNotifierUseCase(notifierRepo, cfg, mqttClient, logger)` ใน `handlers.go`
- unit tests ผ่าน (invalid/unknown-channel/email-not-configured-graceful/health)

**Phase 4 — Auto-control manager**
- `internal/modules/control/` (ใหม่): ตาม pattern notifier
  - `presenter/` — `ExecuteRequest{DeviceID, Event, Topic/MqttControlOn/MqttControlOff override}`, `ExecuteResult`, `ScheduleRunRequest{ScheduleID, Force}`, `ScheduleRunResult`, `HealthResponse`
  - `repository/` — interface: `GetDeviceByID`, `GetActiveSchedules`, `GetScheduleByID`, `GetScheduleDevices`, `CreateScheduleLog`; PG impl ต่อ `sd_iot_device`/`sd_iot_schedule`/`sd_iot_schedule_device`/`sd_schedule_process_log`
  - `usecase/` — `ExecuteControl` (publish MQTT ตาม event: 1=ON→`MqttControlOn`, 0=OFF→`MqttControlOff`, topic จาก `MqttDataControl`/override), `RunSchedule` (จับคู่ weekday+start time หรือ Force; execute ต่อ device; บันทึก schedule log), `Health`; **ใช้ interface `Publisher` แคบ (Publish+IsConnected)** — `pkg/mqtt.Client` เป็นไปตาม เพิ่ม test ได้ง่าย (แก้ pain: `mqtt.Client` มี method `getTopic` unexported ทำให้ mock นอก package ไม่ได้)
  - `delivery/http/` — `POST /control/execute`, `POST /control/schedule/run`, `GET /control/health`
- wiring: block "15c" ใน `handlers.go` — `NewControlUseCase(controlRepo, mqttClient, logger)` nil-safe
- unit tests ผ่าน (invalid/nil-mqtt/not-connected/device-not-found/on-off mapping + override/schedule force/schedule specific/day-mismatch skip/schedule log/nil-repo/health)

**Phase 5 — AuditLog + Blockchain hash-chain**
- `internal/modules/auditlog/` ขยายจาก skeleton → จริง:
  - models (additive ใน `internal/models/models.go` + migrate): `AuditLogEntry` (table `audit_logs`: id/action/resource/resource_id/before/after/hash/prev_hash/nonce/created_at), `BlockchainAnchor` (table `blockchain_anchor`: chain_id/payload/hash/prev_hash/anchored_at/tx_hash)
  - `repository/` — interface: `CreateEntry`, `GetEntryByID`, `GetLastEntry`, `CreateAnchor`, `GetAnchorByChainID`, `CountEntries`; PG impl
  - `provider/` — `AnchorProvider` interface + `NewNoopProvider()` (default no-op → ไม่เขียน external ledger, tx_hash nil, ไม่กระทบของเดิม)
  - `usecase/` — `Create` (hash-chain: `hash=SHA256(prevHash+payload+timestamp+nonce)`, ดึง last ต่อ prev_hash, persist, anchor ผ่าน provider), `Verify` (recompute hash + ตรวจ tamper), `Chain` (length + head); nil-safe ทั้ง repo และ anchor
  - wiring: `NewAuditLogUseCase(auditlogRepo, NewNoopProvider(), logger)` ใน `handlers.go` block 15b
- unit tests ผ่าน (invalid/nil-repo graceful/chain continuity prev_hash/hash ต่างกัน/tamper verify fails/noop anchor tx_hash nil/replay anchor เขียน record/not-found/chain length/nil-repo chain)

**Phase 6 — API Manager (key + rate-limit)**
- `internal/modules/apimanager/` ขยายจาก skeleton → จริง:
  - `repository/` — interface + PG impl ผูกตาราง **`sd_api_key` เดิม** (additive, ไม่ alter — user เลือก) : `Create`/`GetByID`/`GetByKey`/`SetActive`/`IncrementUsage`
  - `usecase/` — `CreateKey` (generate plaintext `apk_...` + SHA256 hashed เก็บใน DB, คืน plaintext ครั้งเดียว), `RevokeKey` (set inactive), `CheckRateLimit` (fixed-window counter ผ่าน Cache, key `apikey:ratelimit:<id>`, default 100/min; nil-cache → allow), `Usage` (จาก usage_count/last_used_at); **interface `Cache` แคบ (Get/Set/Incr)** สำหรับ stub test
  - `internal/server/api_cache_adapter.go` — adapt `redisDb.Cache` (Get/Set) + `redisClient` (INCR) ให้ตรง interface แคบ
  - wiring: `NewAPIManagerUseCase(apimanagerRepo, &apiCacheAdapter{cache: redisCache, redis: redisClient}, logger)` ใน `handlers.go`
- unit tests ผ่าน (invalid/nil-repo graceful/create stored hashed/plaintext≠hash/revoke not-found/success/nil-cache allowed/rate allowed/cache stub/usage nil-repo + with repo)

**Phase 7 — Flow Engine จริง (runtime + wire notifier/realtime)**
- model (additive ใน `internal/models/models.go` + migrate): `FlowDefinition` (table `flow_definitions`: id/name/nodes jsonb/edges jsonb/enabled/created_at/updated_at)
- `repository/` — `Repository` interface (`Create`/`GetByID`/`Update`) + PG impl เก็บ nodes/edges เป็น JSONB
- `usecase/` — แทนที่ skeleton (`ErrNotImplemented` เดิม) ด้วย:
  - `Create` — validate graph → persist (nil-safe repo)
  - `Get` — load by id (not-found → `ErrFlowNotFound`)
  - `Validate` — ตรวจ name/node id ซ้ำ/edge แหล่ง-ปลายทางมีจริง/self-loop + **cycle detection (Kahn's algorithm)**
  - `Trigger` — load flow → run เรียง topological order (กัน cycle) → dispatch ตาม node type: `email|sms|line|discord|io|notify` → **`notifierUC.Dispatch`**, `ws|realtime|monitor` → **`realtimeUC.Publish`**; nil dep → skip แบบ graceful แล้วยัง OK
  - **interface แคบ** `Notifier` (Dispatch) + `Realtime` (Publish) สำหรับ stub test
- wiring: `NewFlowEngineUseCase(flowengineRepo, notifierUC, realtimeUC, logger)` ใน `handlers.go` (ย้าย realtime block ขึ้นก่อน flowengine ให้ `realtimeUC` อยู่ใน scope)
- unit tests ผ่าน (validate invalid: nil/empty/duplicate/dangling/self-loop/cycle; validate valid; create persist/invalid raises/nil-repo graceful; get not-found; trigger dispatch email+ws; trigger nil-deps skip; trigger not-found)

**Phase 8 — รวม wiring + verify ทั้งหมด**
- `go build ./...` ✅
- `go vet ./...` ✅
- `go test ./internal/modules/...` ✅ (ทั้ง 6 module ใหม่: notifier/control/auditlog/apimanager/flowengine/realtime ผ่าน)
- `go test ./...` ✅ ยกเว้น `TestVerifyBrokerEnvOverride` ใน `config` (environment-sensitive: `MQTT_BROKER` ใน shell — มีอยู่ก่อน ไม่เกี่ยวกับโค้ดที่แก้)
- ตรวจ wiring: `handlers.go` block 15a/15b/15c — ทุก module ใหม่ถูก map routes ครบ 6 ระบบ, nil-safe, additive
- smoke บูต server: ต้องการ infra จริง (DB/MQTT/Redis) จึงทำได้เฉพาะใน env ที่มี; โครงสร้าง route/DI ผ่านการ build + vet ยืนยันแล้ว

### ตัดสินใจ/ข้าม ⏹️
- **Phase 3b (auto-hook alarm engine): CANCELLED** — user เลือก option 1: **ไม่แตะ alarm engine เดิม** (ใน `iot/usecase.go` 2684 บรรทัด) เพื่อรักษาข้อบังคับ "ต้องไม่กระทบการทำงานเดิม"
  - ผล: notifier ถูกเรียกผ่าน `POST /api/notifier/dispatch` และจะถูก trigger ผ่าน **Flow Engine (Phase 7)** alarm event → action nodes

### ยังเหลือ ⏳
- (ไม่มี — ครบทุก phase ของแผน 2026-08-28-applied-iot-platform-plan)

### วิธีทำต่อ (วันหลัง)
ครบทุก phase แล้ว (0, 2, 3, 4, 5, 6, 7, 8) ✅

ถ้าจะขยายต่อแนะนำ:
- **smoke จริง** บน env ที่มี DB/MQTT/Redis — ลง endpoint ใหม่แล้วตรวจตาราง migrate (`flow_definitions`, `audit_logs`, `blockchain_anchor`)
- **Flow Engine hook alarm event** (เดิมที่ Phase 3b ตัดสินใจไม่แตะ alarm engine) — ต่อ event → `triggerUC` ด้วย FlowNode type `alarm`/`mqtt-in`
- **API key middleware** แนบกับเส้นทางจริง (ตอนนี้ key สร้าง/revoke/rate-limit ผ่าน usecase แล้ว)
