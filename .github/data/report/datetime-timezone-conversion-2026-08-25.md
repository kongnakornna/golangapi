# รายงานการแปลง Timezone-Aware DateTime

> **วันที่:** 25 สิงหาคม 2026  
> **โมดูลที่แก้ไข:** settings, users, kafka, iot  
> **จุดประสงค์:** แปลง `time.Now()` ที่ใช้กับ datetime fields ให้เป็น timezone-aware ด้วย `helpers.GetTimeLocation()`

---

## สารบัญ

1. [ปัญหาเดิม](#1-ปัญหาเดิม)
2. [รูปแบบที่ใช้แก้ไข](#2-รูปแบบที่ใช้แก้ไข)
3. [รายละเอียดการแก้ไข](#3-รายละเอียดการแก้ไข)
4. [สิ่งที่ไม่แก้ไข](#4-สิ่งที่ไม่แก้ไข)
5. [Build Verification](#5-build-verification)

---

## 1. ปัญหาเดิม

ในโปรเจคใช้ `time.Now()` ตรงๆ สำหรับ datetime fields เช่น `updateddate`, `createddate`, `CreatedAt`, `UpdatedAt` ซึ่งจะใช้ **UTC timezone ของ server** ไม่ใช่ Asia/Bangkok ทำให้:

- เวลาใน DB ไม่ตรงกับเวลาจริงของผู้ใช้
- API response ส่งเวลา UTC กลับไป
- Log และ audit trail ไม่สอดคล้องกัน

---

## 2. รูปแบบที่ใช้แก้ไข

### 2.1 สำหรับ `time.Time` struct fields

```go
// BEFORE
CreatedAt: time.Now(),
UpdatedAt: time.Now(),

// AFTER
CreatedAt: time.Now().In(helpers.GetTimeLocation()),
UpdatedAt: time.Now().In(helpers.GetTimeLocation()),
```

### 2.2 สำหรับ map updates (GORM)

```go
// BEFORE
"updateddate": time.Now(),

// AFTER
"updateddate": time.Now().In(helpers.GetTimeLocation()),
```

### 2.3 สำหรับ string formatting

```go
// BEFORE
time.Now().Format("2006-01-02 15:04:05")

// AFTER
helpers.GetCurrentFullDatenow()
```

---

## 3. รายละเอียดการแก้ไข

### 3.1 `internal/modules/settings/usecase/usecase.go`

| บรรทัด | โค้ดเดิม | โค้ดใหม่ |
|--------|----------|----------|
| 95 | `"updateddate": time.Now()` | `"updateddate": time.Now().In(helpers.GetTimeLocation())` |
| 102 | `"updateddate": time.Now()` | `"updateddate": time.Now().In(helpers.GetTimeLocation())` |

**สิ่งที่เพิ่ม:** import `"icmongolang/pkg/helpers"`

---

### 3.2 `internal/modules/users/repository/pg_repository.go`

| บรรทัด | ฟังก์ชัน | โค้ดเดิม | โค้ดใหม่ |
|--------|----------|----------|----------|
| 246 | `UpdateActiveStatus` | `"updateddate": time.Now()` | `"updateddate": time.Now().In(helpers.GetTimeLocation())` |
| 260 | `UpdateAvatar` | `"updateddate": time.Now()` | `"updateddate": time.Now().In(helpers.GetTimeLocation())` |

**สิ่งที่เพิ่ม:** import `"icmongolang/pkg/helpers"`

---

### 3.3 `internal/modules/kafka/usecase/order_usecase.go`

| บรรทัด | ฟังก์ชัน | ฟิลด์ | โค้ดเดิม | โค้ดใหม่ |
|--------|----------|-------|----------|----------|
| 48 | `CreateOrder` | `CreatedAt` | `time.Now()` | `time.Now().In(helpers.GetTimeLocation())` |
| 49 | `CreateOrder` | `UpdatedAt` | `time.Now()` | `time.Now().In(helpers.GetTimeLocation())` |
| 85 | `ProcessOrderMessage` | `UpdatedAt` | `time.Now()` | `time.Now().In(helpers.GetTimeLocation())` |
| 106 | `HandleWebSocketMessage` | `CreatedAt` | `time.Now()` | `time.Now().In(helpers.GetTimeLocation())` |
| 107 | `HandleWebSocketMessage` | `UpdatedAt` | `time.Now()` | `time.Now().In(helpers.GetTimeLocation())` |
| 142 | `BroadcastOrderStatus` | `timestamp` | `time.Now()` | `time.Now().In(helpers.GetTimeLocation())` |

**สิ่งที่เพิ่ม:** import `"icmongolang/pkg/helpers"`

---

### 3.4 `internal/modules/iot/usecase/usecase.go`

| บรรทัด | ฟังก์ชัน | ฟิลด์ | โค้ดเดิม | โค้ดใหม่ |
|--------|----------|-------|----------|----------|
| 1655 | `UpdateDeviceStatus` | `UpdatedAt` | `time.Now()` | `time.Now().In(helpers.GetTimeLocation())` |
| 1694 | `GetDeviceConfig` | `CreatedAt` | `time.Now()` | `time.Now().In(helpers.GetTimeLocation())` |
| 1695 | `GetDeviceConfig` | `UpdatedAt` | `time.Now()` | `time.Now().In(helpers.GetTimeLocation())` |
| 1718 | `UpdateDeviceConfig` | `UpdatedAt` | `time.Now()` | `time.Now().In(helpers.GetTimeLocation())` |
| 1756 | `ProcessMqttData` | `Timestamp` | `time.Now()` | `time.Now().In(helpers.GetTimeLocation())` |
| 1757 | `ProcessMqttData` | `CreatedAt` | `time.Now()` | `time.Now().In(helpers.GetTimeLocation())` |
| 2181 | `logActivity` | `Timestamp` | `time.Now()` | `time.Now().In(helpers.GetTimeLocation())` |
| 2182 | `logActivity` | `CreatedAt` | `time.Now()` | `time.Now().In(helpers.GetTimeLocation())` |

**หมายเหตุ:** ไฟล์นี้ import `"icmongolang/pkg/helpers"` อยู่แล้ว

---

### 3.5 `internal/modules/iot/iothelper/alarm.go`

| บรรทัด | ฟังก์ชัน | ฟิลด์ | โค้ดเดิม | โค้ดใหม่ |
|--------|----------|-------|----------|----------|
| 324 | `processAlarmDetail` | `Timestamp` | `time.Now().Format("2006-01-02 15:04:05")` | `helpers.GetCurrentFullDatenow()` |

**สิ่งที่เพิ่ม:** import `"icmongolang/pkg/helpers"`  
**สิ่งที่ลบ:** import `"time"` (ไม่ใช้แล้ว)

---

## 4. สิ่งที่ไม่แก้ไข

รายการต่อไปนี้ **ไม่ต้องแก้** เพราะไม่เกี่ยวกับการแสดงเวลา:

| รูปแบบ | เหตุผล | ตัวอย่าง |
|--------|--------|----------|
| `time.Now().Add(...)` | สำหรับ duration/timeout | `10*time.Second`, `time.Minute*15` |
| `time.Now().Unix()` | สำหรับ Redis score | `float64(time.Now().Add(delay).Unix())` |
| GORM auto-fill `CreatedAt/UpdatedAt` | GORM จัดการเองผ่าน `NowFunc` | ไม่ต้องแก้ |
| Rate limiter `lastSeen` | ไม่เกี่ยวกับ display | `internal/middleware/rate_limit.go` |
| Monitoring `StartTime` | สำหรับคำนวณ duration | `internal/middleware/monitoring.go` |
| WebSocket `SetReadDeadline` | สำหรับ timeout | `internal/modules/websocket/` |
| InfluxDB `WritePoint` | ใช้ `time.Time` สำหรับ TSDB | `internal/modules/iot/usecase.go:1810` |

---

## 5. Build Verification

```bash
$ go build ./...
# (no output = success)
```

---

## สรุป

| รายการ | จำนวน |
|--------|-------|
| ไฟล์ที่แก้ไข | 5 |
| จุดที่แก้ไขทั้งหมด | 18 |
| ไฟล์ที่เพิ่ม import helpers | 4 |
| ไฟล์ที่ลบ import time | 1 |
| Build status | PASS |
