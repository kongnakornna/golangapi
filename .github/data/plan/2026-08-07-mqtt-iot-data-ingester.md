# Plan: MQTT → iot_data Background Ingester

## ภาพรวมการแก้ไข

| # | Layer | ไฟล์ | สิ่งที่แก้ |
|---|-------|------|-----------|
| 1 | cmd | `cmd/migrate.go` | แก้ skip key จาก `fmt.Sprintf("%T", m)` เป็น package path (`reflect.TypeOf` + `PkgPath`) → `iot_data` และ models ของ iot module อีก 7 ตัว ไม่โดน skip โดยไม่ตั้งใจ |
| 2 | repo | `internal/modules/iot/repository/device_repo.go` | เพิ่ม filter `MqttDeviceName` ใน `DeviceListAlarmRequest` + `applyFilters` |
| 3 | usecase | `internal/modules/iot/usecase/usecase.go` | เพิ่ม `buildDataMap` (parse JSON object / comma), `ProcessMqttData` รองรับ 2 format, เพิ่ม `StartIngest` + `resolveDeviceID` + cache, เพิ่ม field/mutex/interface method |
| 4 | server | `internal/server/handlers.go` | เรียก `iotUC.StartIngest` หลังสร้าง iotUC |

---

## Section 1 — แก้ migrate skip bug (`cmd/migrate.go`)

**Location:** `cmd/migrate.go:63-84` (skip map), `cmd/migrate.go:210-217` (filter loop), import block `cmd/migrate.go:3-14`

**From → to:** ตอนนี้ skip key เป็น `"*models.X"` จาก `fmt.Sprintf("%T", m)` ซึ่งแสดงแค่ package name (`models`) — ทั้ง `internal/models` และ `internal/modules/iot/models` ต่างมี package name `models` จึงทำให้ `&iotmodels.IotData{}` ได้ค่า `%T` = `*models.IotData` เท่ากับ internal/models → โดน skip map กลืนไปด้วย (ผลกระทบเดียวกันกับ `iotmodels.ActivityLog, Device, DeviceCategory, DeviceGroup, DeviceNotificationConfig, DeviceSchedule, DeviceStatusHistory` อีก 7 ตัว) → `iot_data` table ไม่เคยถูกสร้าง

เปลี่ยนเป็น key ที่ใช้ import path เต็ม (`PkgPath + "." + Name`) เพื่อแยก package ทั้งสองออกจากกันอย่างชัดเจน

**Change detail:**

แก้ import block:
```go
import (
	"fmt"
	"icmongolang/config"
	iotmodels "icmongolang/internal/modules/iot/models"
	kafkamodels "icmongolang/internal/modules/kafka/models"
	"icmongolang/internal/models"
	"icmongolang/pkg/db/postgres"
	"icmongolang/pkg/logger"
	"reflect"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)
```

แก้ skip map (เฉพาะ key — ค่า `true` คงเดิม, เพิ่ม comment ชี้ว่าเป็น internal/models โดยเฉพาะ):
```go
	skipModels := map[string]bool{
		// Skip models with Device relation (because Device has Mqtt/Location issues)
		"icmongolang/internal/models.Device":                   true,
		"icmongolang/internal/models.DeviceNotificationConfig": true,
		"icmongolang/internal/models.DeviceSchedule":           true,
		"icmongolang/internal/models.DeviceStatusHistory":      true,
		"icmongolang/internal/models.SdActivityLog":            true, // has ActivityType relation
		"icmongolang/internal/models.DeviceGroupMember":        true, // may also have Device relation
		"icmongolang/internal/models.NotificationCondition":    true,
		"icmongolang/internal/models.NotificationLog":          true,
		"icmongolang/internal/models.ReportData":               true,
		"icmongolang/internal/models.SensorData":               true,
		"icmongolang/internal/models.DeviceGroup":              true, // might have nested issues
		"icmongolang/internal/models.DeviceCategory":           true,
		"icmongolang/internal/models.DeviceStatus":             true, // if it has Device relation
		"icmongolang/internal/models.DeviceConfig":             true,
		"icmongolang/internal/models.CommandLog":               true,
		"icmongolang/internal/models.DeviceAlert":              true,
		"icmongolang/internal/models.IotData":                  true,
		"icmongolang/internal/models.ActivityLog":              true,
		// Add any other models that cause "invalid field" errors here
	}
```

แก้ filter loop (lines 210-217):
```go
	// Filter out skipped models
	var toMigrate []interface{}
	for _, m := range allModels {
		t := reflect.TypeOf(m)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		typeName := t.PkgPath() + "." + t.Name()
		if skipModels[typeName] {
			log.Warnf("Skipping model %s due to known issue", typeName)
			continue
		}
		toMigrate = append(toMigrate, m)
	}
```

**ผลลัพธ์ที่คาด:** log เปลี่ยนจาก `Migrating 89 models (skipped 18)` เป็น `Migrating 97 models (skipped 18)` และ `iot_data` (พร้อมอีก 7 tables ของ iot module) ถูกสร้างตอน restart

---

## Section 2 — เพิ่ม filter `MqttDeviceName` (`internal/modules/iot/repository/device_repo.go`)

**Location:** `device_repo.go:113-145` (struct), `device_repo.go:283-363` (applyFilters)

**From → to:** ตอนนี้ `DeviceListAlarmRequest` มี `DeviceID`, `MqttID` แต่ไม่มี filter ด้วยชื่อ MQTT ของ device (`mqtt_device_name`); ingester ต้องค้น device จาก topic prefix จึงต้องเพิ่ม filter นี้

**Change detail:**

เพิ่ม field ต่อจาก `MqttID` (บรรทัด ~115):
```go
type DeviceListAlarmRequest struct {
	DeviceID        string
	MqttID          string
	MqttDeviceName  string
	Keyword         string
	...
```

เพิ่ม filter ใน `applyFilters` ต่อจากบล็อก `MqttID` (บรรทัด ~295):
```go
		if req.MqttID != "" {
			query = query.Where("d.mqtt_id = ?", req.MqttID)
		}
		if req.MqttDeviceName != "" {
			query = query.Where("d.mqtt_device_name = ?", req.MqttDeviceName)
		}
```

---

## Section 3 — iot usecase: parse 2 format + ingester (`internal/modules/iot/usecase/usecase.go`)

### 3.1 Imports
**From → to:** เพิ่ม `"sync"` (ยังไม่มี) และ paho alias (เพราะชื่อ `mqtt` ถูกใช้โดย `icmongolang/pkg/mqtt` แล้ว)

```go
import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"icmongolang/config"
	"icmongolang/internal/modules/iot/iothelper"
	"icmongolang/internal/modules/iot/models"
	"icmongolang/internal/modules/iot/presenter"
	"icmongolang/internal/modules/iot/repository"
	"icmongolang/pkg/helpers"
	"icmongolang/pkg/influxdb"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/mqtt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/redis/go-redis/v9"
)
```

### 3.2 Interface `MQTT3UseCase`
**Location:** `usecase.go:26-59`

เพิ่ม method (หลัง `ExportData`, บรรทัด ~58):
```go
	StartIngest(ctx context.Context) error
```

### 3.3 Struct + Constructor
**Location:** `usecase.go:61-78` (struct), `usecase.go:96-112` (constructor)

เพิ่ม field (หลัง `deviceAlertRepo`):
```go
type mqtt3UseCase struct {
	...
	deviceAlertRepo  repository.DeviceAlertRepository

	deviceCacheMu   sync.Mutex
	deviceNameCache map[string]string
}
```

เพิ่ม init ใน constructor (ก่อน `return &mqtt3UseCase{...}`):
```go
	return &mqtt3UseCase{
		deviceRepo:       deviceRepo,
		...
		deviceAlertRepo:  deviceAlertRepo,
		deviceNameCache:  make(map[string]string),
	}
```

### 3.4 `buildDataMap` helper (package-level function ใหม่)
**Location:** เพิ่มใหม่ใกล้ๆ `ProcessMqttData` (ก่อนบรรทัด 1526)

```go
// buildDataMap แปลง rawData เป็น map[string]interface{}:
//   - ถ้า rawData เป็น JSON object → ใช้ key/value โดยตรง
//   - ถ้าไม่ใช่ JSON → split ด้วย comma แล้ว map index ตาม configMap (MqttStatusDataName)
// ทั้งสองกรณีเก็บค่าเดิมทั้งหมดไว้ใน key "raw"
func buildDataMap(rawData string, configMap map[string]string) map[string]interface{} {
	var jsonObj map[string]interface{}
	if err := json.Unmarshal([]byte(rawData), &jsonObj); err == nil {
		jsonObj["raw"] = rawData
		return jsonObj
	}

	parts := strings.Split(rawData, ",")
	dataMap := make(map[string]interface{}, len(parts)+1)
	for i, val := range parts {
		key := fmt.Sprintf("%d", i)
		if configMap != nil {
			if mapped, ok := configMap[key]; ok {
				key = mapped
			}
		}
		trimmed := strings.TrimSpace(val)
		if f, err := strconv.ParseFloat(trimmed, 64); err == nil {
			dataMap[key] = f
		} else {
			dataMap[key] = trimmed
		}
	}
	dataMap["raw"] = rawData
	return dataMap
}
```

### 3.5 `ProcessMqttData` — ใช้ `buildDataMap`
**Location:** `usecase.go:1546-1568`

**From → to:** ตอนนี้ parse แบบ comma อย่างเดียว inline; เปลี่ยนเป็นลอง JSON object ก่อน (ผ่าน helper)

แทนที่ block ข้อ 2-4 (บรรทัด 1546-1568):
```go
	// 2. แยก payload: ลอง JSON object ก่อน, fallback เป็น comma-separated
	dataMap := buildDataMap(rawData, configMap)
```
(block 5-9: `json.Marshal`, `IotData`, `Create`, `UpdateDeviceStatus`, `logActivity` ไม่เปลี่ยน)

### 3.6 `StartIngest` + helpers (methods ใหม่)
**Location:** เพิ่มใหม่หลัง `ProcessMqttData` (หลังบรรทัด 1596)

```go
// StartIngest เริ่ม ingester: subscribe "#" เพื่อรับทุก topic และ persist ลง iot_data
func (u *mqtt3UseCase) StartIngest(ctx context.Context) error {
	if u.mqttClient == nil || !u.mqttClient.IsConnected() {
		return fmt.Errorf("ingester: MQTT client not connected")
	}
	handler := func(client paho.Client, msg paho.Message) {
		topic := msg.Topic()
		payload := string(msg.Payload())
		deviceName := extractDeviceNameFromTopic(topic)
		if deviceName == "" {
			u.logger.Debugf("ingester: cannot extract device name from topic %s, ignoring", topic)
			return
		}
		deviceID, ok := u.resolveDeviceID(ctx, deviceName)
		if !ok {
			u.logger.Debugf("ingester: no device found for mqtt_device_name=%s (topic %s), skipping", deviceName, topic)
			return
		}
		go u.processIncoming(ctx, deviceID, payload)
	}
	if err := u.mqttClient.Subscribe("#", 0, handler); err != nil {
		return fmt.Errorf("ingester: subscribe # failed: %w", err)
	}
	u.logger.Info("Ingester subscribed to # for iot_data persistence")
	return nil
}

// extractDeviceNameFromTopic ตัด segment แรกของ topic เป็น device name (เช่น TEST01/DATA → TEST01)
func extractDeviceNameFromTopic(topic string) string {
	parts := strings.Split(strings.TrimPrefix(topic, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return ""
	}
	return parts[0]
}

// resolveDeviceID ค้น device_id จาก mqtt_device_name พร้อม cache ใน memory (thread-safe)
func (u *mqtt3UseCase) resolveDeviceID(ctx context.Context, deviceName string) (string, bool) {
	u.deviceCacheMu.Lock()
	defer u.deviceCacheMu.Unlock()
	if id, ok := u.deviceNameCache[deviceName]; ok {
		return id, true
	}
	result, err := u.deviceRepo.ListDevicesWithAlarm(ctx, &repository.DeviceListAlarmRequest{
		MqttDeviceName: deviceName,
		Page:           1,
		PageSize:       1,
	})
	if err != nil {
		u.logger.Errorf("ingester: device lookup error for %s: %v", deviceName, err)
		return "", false
	}
	if len(result.Items) == 0 {
		return "", false
	}
	id := strconv.Itoa(result.Items[0].DeviceID)
	u.deviceNameCache[deviceName] = id
	return id, true
}

// processIncoming จัดการ message นอก goroutine ของ paho พร้อม recover เพื่อไม่ให้ server crash
func (u *mqtt3UseCase) processIncoming(ctx context.Context, deviceID, payload string) {
	defer func() {
		if r := recover(); r != nil {
			u.logger.Errorf("ingester: panic processing device=%s payload=%q: %v", deviceID, payload, r)
		}
	}()
	if _, err := u.ProcessMqttData(ctx, deviceID, payload); err != nil {
		u.logger.Errorf("ingester: ProcessMqttData error for device=%s: %v", deviceID, err)
	}
}
```

**หมายเหตุ (ทำไม subscribe `#` ไม่ชนกับ requestManager):** `pkg/mqtt/client.go:55` ลง `client.AddRoute("#", rm.handleIncomingMessage)` ไว้เป็น **route** (ไม่ใช่ subscription); paho router fire ทุก handler ที่ match — ทั้ง routes และ subscriptions — ดังนั้น ingester subscribe `#` แยกได้ โดย `handleIncomingMessage` ยังทำงานเดิม (จัดการ pending request) ไม่กระทบ

---

## Section 4 — Wire ingester ใน server (`internal/server/handlers.go`)

**Location:** หลัง iotUC สร้างเสร็จ (หลังบรรทัด 240), import block (บรรทัด 3-5)

**From → to:** ตอนนี้สร้าง iotUC แล้ว register routes เท่านั้น; เพิ่มการ start ingester โดยตรวจว่า mqttClient connected ก่อน (กรณีเดียวกับบล็อก section 11)

**Change detail:**

เพิ่ม import:
```go
import (
	"context"
	"net/http"
	"time"
	...
```

เพิ่มหลัง `iotUC := iotUsecase.NewMQTT3UseCase(...)` (บรรทัด 226-240) ก่อน `iotHandler`:
```go
	if mqttClient != nil && mqttClient.IsConnected() {
		if err := iotUC.StartIngest(context.Background()); err != nil {
			logger.Errorf("❌ Failed to start IoT ingester: %v", err)
		} else {
			logger.Info("✅ IoT MQTT ingester started (topic #)")
		}
	} else {
		logger.Warn("⚠️ MQTT client not connected – IoT MQTT ingester skipped")
	}
```

---

## Unit Tests

**ตัดสินใจ: เขียน**

สร้าง `internal/modules/iot/usecase/usecase_test.go` (package `usecase`) — ทดสอบเฉพาะ function บริสุทธิ์ (ไม่ต้อง DB/mocks):

| Test | Input | Expected |
|------|-------|----------|
| `TestBuildDataMap_JSONObject` | `{"temp":25.5,"hum":60}` | map มี `temp=25.5`, `hum=60`, `raw` = payload เดิม |
| `TestBuildDataMap_CommaWithConfigMap` | `"10,20,30"` + `{"0":"temp","1":"hum"}` | `temp=10.0`, `hum=20.0`, `"2"=30.0`, `raw` |
| `TestBuildDataMap_CommaWithoutConfig` | `"1.5,on"` + nil | `"0"=1.5`, `"1"="on"`, `raw` |
| `TestBuildDataMap_ScalarNonJSON` | `"23.5"` | `"0"=23.5`, `raw` |
| `TestExtractDeviceNameFromTopic` | `TEST01/DATA`, `/TEST01/DATA`, `DATA`, `""`, `"/"` | `TEST01`, `TEST01`, `DATA`, `""`, `""` |

รัน:
```bash
docker compose exec backend go test ./internal/modules/iot/usecase/
```

---

## Verification (manual integration)

จาก repo root (`C:\github\icmongolang`):

1. `docker compose restart backend`
2. Poll `http://localhost:5000/health` จนเป็น 200 (~90s: wait-for-it 15s + migrate + initdata + air build)
3. `docker compose logs backend` → ตรวจ
   - `Migrating 97 models (skipped 18)`
   - `✅ IoT MQTT ingester started (topic #)`
   - ไม่มี error/panic
4. ตรวจ table ถูกสร้าง:
   ```bash
   docker compose exec postgres psql -U icmongolang -d icmongolang -c "\dt iot_data"
   ```
5. insert device ทดสอบ (ดู schema ก่อน `\d sd_iot_device` / `\d sd_iot_mqtt`; ต้องมี NOT NULL columns ครบ, status=1 ทั้งคู่, `mqtt_id` เชื่อมกัน):
   ```sql
   INSERT INTO sd_iot_mqtt (mqtt_name, status) VALUES ('test-broker', 1) RETURNING mqtt_id;
   INSERT INTO sd_iot_device (device_id, device_name, mqtt_device_name, mqtt_id, status) VALUES (1, 'TEST01', 'TEST01', <mqtt_id>, 1);
   ```
6. signin → token (`curl.exe -d @.../signin.json`) → `POST /api/mqtt/publish`
   - payload JSON object: `{"topic":"TEST01/DATA","payload":"{\"temp\":25.5,\"hum\":60}"}`
   - payload comma: `{"topic":"TEST01/DATA","payload":"10,20,30"}`
7. ตรวจ row:
   ```bash
   docker compose exec postgres psql -U icmongolang -d icmongolang -c "SELECT id, device_id, data, timestamp FROM iot_data ORDER BY id DESC LIMIT 5;"
   ```
   → ต้องเจอ row ของ TEST01 (data มี `temp`/`hum` สำหรับ JSON, และ `"0"/"1"/"2"` สำหรับ comma)
8. publish topic ไม่รู้จัก `{"topic":"UNKNOWN01/DATA","payload":"5"}` → ไม่มี row ใหม่ + log `no device found for mqtt_device_name=UNKNOWN01`
9. cleanup: ลบไฟล์ throwaway ใน container `/app/tmp/iotmigrate_check{,2,3,4}.go` และบน host `%TEMP%\opencode\iotmigrate_check*.go`, `pub1-3.json` (ถ้าไม่ใช้)
