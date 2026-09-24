# 🔌 PART 3  — MODULE: `device` (IoT Device Management)

> **ขนาด**: ใหญ่พิเศษ — แยก 4 ตอนย่อย
> **Part 3A**: Domain Layer (Entities + VOs + Shadow + Events + Errors + Services)
> **Part 3B**: Application Layer (ทุก Use Case + DTO + Mappers)
> **Part 3C**: Infrastructure (Postgres + InfluxDB + MQTT + Kafka + Redis + ES + AI)
> **Part 3D**: Interface + Wiring + Migration + Tests + cmd/workers

> **Pattern เฉพาะของ module นี้**:
> 1. **Device Shadow** (digital twin — desired/reported/delta)
> 2. **MQTT hot path** (telemetry ingest ที่ throughput สูง)
> 3. **InfluxDB time-series** (แยกจาก PostgreSQL)
> 4. **Command & Control** (state machine + ACK tracking)
> 5. **Automation Engine** (event-condition-action)
> 6. **AI Inference** (anomaly/forecast ผ่าน outbound port)
> 7. **Multi-protocol** (MQTT/LoRaWAN/Modbus/HTTP)

---

# 🅰️ PART 3A — DOMAIN LAYER

## A.1 โครงสร้าง Domain

```
internal/modules/device/domain/
├── entity/
│   ├── device.go                # Aggregate Root #1
│   ├── device_shadow.go         # Aggregate Root #2 (Digital Twin)
│   ├── command.go               # Aggregate Root #3
│   ├── alert_rule.go            # Aggregate Root #4
│   ├── automation.go            # Aggregate Root #5
│   ├── device_group.go          # Aggregate Root #6
│   ├── device_model.go          # Entity
│   ├── firmware.go              # Entity
│   ├── ai_model.go              # Entity
│   ├── telemetry.go             # Value Object (event-sourced)
│   └── alert_event.go           # Entity (triggered alert)
├── value_object/
│   ├── device_status.go
│   ├── device_type.go
│   ├── protocol.go
│   ├── metric.go
│   ├── metric_value.go
│   ├── capability.go
│   ├── command_status.go
│   ├── alert_severity.go
│   ├── automation_trigger.go
│   ├── device_serial.go
│   ├── device_token.go
│   └── shadow_state.go          # desired/reported/delta
├── repository/
│   ├── device_repository.go
│   ├── shadow_repository.go
│   ├── command_repository.go
│   ├── telemetry_repository.go  # InfluxDB
│   ├── alert_rule_repository.go
│   ├── alert_event_repository.go
│   ├── automation_repository.go
│   ├── device_group_repository.go
│   ├── device_model_repository.go
│   └── firmware_repository.go
├── service/
│   ├── provisioning_service.go
│   ├── health_monitor_service.go
│   ├── shadow_merge_service.go
│   ├── alert_evaluator.go
│   ├── automation_engine.go
│   ├── command_dispatcher.go
│   ├── telemetry_aggregator.go
│   └── port/
│       ├── mqtt_publisher_port.go
│       ├── ai_inference_port.go
│       ├── notifier_port.go
│       └── code_generator_port.go
├── event/
│   ├── device_events.go
│   ├── telemetry_events.go
│   ├── command_events.go
│   ├── alert_events.go
│   └── automation_events.go
└── errors/
    └── errors.go
```

---

## A.2 Value Objects

### `domain/value_object/device_status.go`
```go
package valueobject

type DeviceStatus string

const (
    DeviceStatusRegistered     DeviceStatus = "REGISTERED"      // สร้างแล้ว ยังไม่ provision
    DeviceStatusProvisioned    DeviceStatus = "PROVISIONED"     // มี token แล้ว ยังไม่ connect
    DeviceStatusOnline         DeviceStatus = "ONLINE"
    DeviceStatusOffline        DeviceStatus = "OFFLINE"
    DeviceStatusFault          DeviceStatus = "FAULT"
    DeviceStatusMaintenance    DeviceStatus = "MAINTENANCE"
    DeviceStatusDecommissioned DeviceStatus = "DECOMMISSIONED"
)

func (s DeviceStatus) IsValid() bool {
    switch s {
    case DeviceStatusRegistered, DeviceStatusProvisioned, DeviceStatusOnline,
        DeviceStatusOffline, DeviceStatusFault, DeviceStatusMaintenance,
        DeviceStatusDecommissioned:
        return true
    }
    return false
}

func (s DeviceStatus) CanTransitionTo(next DeviceStatus) bool {
    t := map[DeviceStatus][]DeviceStatus{
        DeviceStatusRegistered:     {DeviceStatusProvisioned, DeviceStatusDecommissioned},
        DeviceStatusProvisioned:    {DeviceStatusOnline, DeviceStatusOffline, DeviceStatus.Decommissioned},
        DeviceStatusOnline:         {DeviceStatusOffline, DeviceStatusFault, DeviceStatus.Maintenance, DeviceStatus.Decommissioned},
        DeviceStatusOffline:        {DeviceStatusOnline, DeviceStatus.Fault, DeviceStatus.Maintenance, DeviceStatus.Decommissioned},
        DeviceStatusFault:          {DeviceStatus.Online, DeviceStatus.Maintenance, DeviceStatus.Offline, DeviceStatus.Decommissioned},
        DeviceStatusMaintenance:    {DeviceStatus.Online, DeviceStatus.Offline, DeviceStatus.Decommissioned},
        DeviceStatusDecommissioned: {}, // terminal
    }
    for _, allowed := range t[s] {
        if allowed == next { return true }
    }
    return false
}

func (s DeviceStatus) IsTerminal() bool {
    return s == DeviceStatusDecommissioned
}

func (s DeviceStatus) CanReceiveCommand() bool {
    return s == DeviceStatusOnline || s == DeviceStatus.Provisioned
}

func (s DeviceStatus) CanIngestTelemetry() bool {
    return s != DeviceStatusDecommissioned
}

func (s DeviceStatus) String() string { return string(s) }
```

### `domain/value_object/device_type.go`
```go
package valueobject

type DeviceType string

const (
    DeviceTypeSensor      DeviceType = "SENSOR"
    DeviceTypeGateway     DeviceType = "GATEWAY"
    DeviceTypeActuator    DeviceType = "ACTUATOR"
    DeviceTypeCamera      DeviceType = "CAMERA"
    DeviceTypeController  DeviceType = "CONTROLLER"
    DeviceTypeMeter       DeviceType = "METER"
    DeviceTypeTracker     DeviceType = "TRACKER"
)

func (t DeviceType) IsValid() bool {
    switch t {
    case DeviceTypeSensor, DeviceTypeGateway, DeviceTypeActuator,
        DeviceTypeCamera, DeviceTypeController, DeviceTypeMeter, DeviceTypeTracker:
        return true
    }
    return false
}

func (t DeviceType) IsReadOnly() bool {
    return t == DeviceTypeSensor || t == DeviceTypeMeter || t == DeviceTypeTracker
}

func (t DeviceType) IsControllable() bool {
    return t == DeviceTypeActuator || t == DeviceTypeController || t == DeviceTypeGateway
}
```

### `domain/value_object/protocol.go`
```go
package valueobject

type Protocol string

const (
    ProtocolMQTT    Protocol = "MQTT"
    ProtocolHTTP    Protocol = "HTTP"
    ProtocolLoRaWAN Protocol = "LORAWAN"
    ProtocolModbus  Protocol = "MODBUS"
    ProtocolZigbee  Protocol = "ZIGBEE"
    ProtocolBLE     Protocol = "BLE"
    ProtocolNB_IoT  Protocol = "NBIOT"
)

func (p Protocol) IsValid() bool {
    switch p {
    case ProtocolMQTT, ProtocolHTTP, ProtocolLoRaWAN, ProtocolModbus,
        ProtocolZigbee, ProtocolBLE, ProtocolNB_IoT:
        return true
    }
    return false
}

func (p Protocol) IsPush() bool {
    // server push ได้หรือไม่
    return p == ProtocolMQTT || p == ProtocolHTTP
}

func (p Protocol) DefaultQoS() int {
    if p == ProtocolMQTT { return 1 }
    return 0
}
```

### `domain/value_object/metric.go`
```go
package valueobject

type Metric string

const (
    MetricTemperature   Metric = "temperature"
    MetricHumidity      Metric = "humidity"
    MetricSoilMoisture  Metric = "soil_moisture"
    MetricSoilPH        Metric = "soil_ph"
    MetricLight         Metric = "light"
    MetricCO2           Metric = "co2"
    MetricPM25          Metric = "pm25"
    MetricPower         Metric = "power"
    MetricEnergy        Metric = "energy"
    MetricVoltage       Metric = "voltage"
    MetricCurrent       Metric = "current"
    MetricWaterFlow     Metric = "water_flow"
    MetricWaterLevel    Metric = "water_level"
    MetricPressure      Metric = "pressure"
    MetricMotion        Metric = "motion"
    MetricDoorState     Metric = "door_state"
    MetricOccupancy     Metric = "occupancy"
    MetricBatteryLevel  Metric = "battery_level"
    MetricSignalRSSI    Metric = "signal_rssi"
)

var metricUnits = map[Metric]string{
    MetricTemperature:  "°C",
    MetricHumidity:     "%",
    MetricSoilMoisture: "%",
    MetricSoilPH:       "pH",
    MetricLight:        "lux",
    MetricCO2:          "ppm",
    MetricPM25:         "µg/m³",
    MetricPower:        "W",
    MetricEnergy:       "kWh",
    MetricVoltage:      "V",
    MetricCurrent:      "A",
    MetricWaterFlow:    "L/min",
    MetricWaterLevel:   "cm",
    MetricPressure:     "hPa",
    MetricMotion:       "bool",
    MetricDoorState:    "bool",
    MetricOccupancy:    "int",
    MetricBatteryLevel: "%",
    MetricSignalRSSI:   "dBm",
}

func (m Metric) IsValid() bool {
    _, ok := metricUnits[m]
    return ok
}

func (m Metric) Unit() string { return metricUnits[m] }

func (m Metric) String() string { return string(m) }
```

### `domain/value_object/metric_value.go`
```go
package valueobject

import (
    "math"
    "time"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
)

// MetricValue – ค่าของ metric หนึ่งค่าพร้อม metadata
type MetricValue struct {
    Metric    Metric
    Value     float64
    Quality   string    // GOOD | UNCERTAIN | BAD
    Timestamp time.Time
}

func NewMetricValue(metric Metric, value float64, ts time.Time) (MetricValue, error) {
    if !metric.IsValid() {
        return MetricValue{}, domainerrors.ErrInvalidMetric
    }
    if math.IsNaN(value) || math.IsInf(value, 0) {
        return MetricValue{}, domainerrors.ErrInvalidMetricValue
    }
    if ts.IsZero() {
        ts = time.Now()
    }
    return MetricValue{
        Metric: metric, Value: value,
        Quality: "GOOD", Timestamp: ts,
    }, nil
}

func (mv MetricValue) IsZero() bool { return mv.Metric == "" }
```

### `domain/value_object/capability.go`
```go
package valueobject

import domainerrors "icmongolang/internal/modules/device/domain/errors"

// Capability – ความสามารถของ device
type Capability struct {
    Type     string     // SENSOR | ACTUATOR | GATEWAY
    Metric   Metric     // สำหรับ SENSOR
    Command  string     // สำหรับ ACTUATOR (on/off/set)
    MinValue *float64
    MaxValue *float64
    Unit     string
}

type CapabilitySet []Capability

func (cs CapabilitySet) HasMetric(m Metric) bool {
    for _, c := range cs {
        if c.Metric == m { return true }
    }
    return false
}

func (cs CapabilitySet) HasCommand(cmd string) bool {
    for _, c := range cs {
        if c.Command == cmd { return true }
    }
    return false
}

func (cs CapabilitySet) Validate() error {
    if len(cs) == 0 {
        return domainerrors.ErrEmptyCapabilities
    }
    for _, c := range cs {
        if c.Type == "" {
            return domainerrors.ErrInvalidCapability
        }
        if c.Type == "SENSOR" && !c.Metric.IsValid() {
            return domainerrors.ErrInvalidMetric
        }
        if c.MinValue != nil && c.MaxValue != nil && *c.MinValue > *c.MaxValue {
            return domainerrors.ErrInvalidCapability
        }
    }
    return nil
}
```

### `domain/value_object/command_status.go`
```go
package valueobject

type CommandStatus string

const (
    CommandStatusPending   CommandStatus = "PENDING"    // สร้างแล้ว ยังไม่ส่ง
    CommandStatusQueued    CommandStatus = "QUEUED"     // device offline — queue ไว้
    CommandStatusSent      CommandStatus = "SENT"       // MQTT ส่งแล้ว
    CommandStatusAcked     CommandStatus = "ACKED"      // device ACK
    CommandStatusFailed    CommandStatus = "FAILED"     // device NACK / timeout
    CommandStatusTimeout   CommandStatus = "TIMEOUT"
    CommandStatusCancelled CommandStatus = "CANCELLED"
)

func (s CommandStatus) IsTerminal() bool {
    switch s {
    case CommandStatusAcked, CommandStatusFailed,
        CommandStatusTimeout, CommandStatusCancelled:
        return true
    }
    return false
}

func (s CommandStatus) IsPending() bool {
    return s == CommandStatusPending || s == CommandStatusQueued || s == CommandStatusSent
}

func (s CommandStatus) String() string { return string(s) }
```

### `domain/value_object/alert_severity.go`
```go
package valueobject

type AlertSeverity string

const (
    AlertSeverityInfo     AlertSeverity = "INFO"
    AlertSeverityWarn     AlertSeverity = "WARN"
    AlertSeverityError    AlertSeverity = "ERROR"
    AlertSeverityCritical AlertSeverity = "CRITICAL"
)

func (s AlertSeverity) IsValid() bool {
    switch s {
    case AlertSeverityInfo, AlertSeverityWarn, AlertSeverityError, AlertSeverityCritical:
        return true
    }
    return false
}

func (s AlertSeverity) Level() int {
    switch s {
    case AlertSeverityInfo:     return 1
    case AlertSeverityWarn:     return 2
    case AlertSeverityError:    return 3
    case AlertSeverityCritical: return 4
    }
    return 0
}
```

### `domain/value_object/device_serial.go`
```go
package valueobject

import (
    "regexp"
    "strings"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
)

// DeviceSerial – Serial number ของอุปกรณ์
type DeviceSerial string

var serialPattern = regexp.MustCompile(`^[A-Z0-9][A-Z0-9\-]{5,63}$`)

func NewDeviceSerial(raw string) (DeviceSerial, error) {
    s := strings.ToUpper(strings.TrimSpace(raw))
    if !serialPattern.MatchString(s) {
        return "", domainerrors.ErrInvalidSerial
    }
    return DeviceSerial(s), nil
}

func (s DeviceSerial) String() string { return string(s) }
```

### `domain/value_object/device_token.go`
```go
package valueobject

import (
    "crypto/rand"
    "encoding/base64"
    "strings"

    "golang.org/x/crypto/bcrypt"

    domainerrors "icmongolang/internal/modules/device/domain/errors"
)

// DeviceToken – ค่า token สำหรับ device auth (plaintext + hash)
type DeviceToken struct {
    Plaintext string // จะ return ครั้งเดียว ไม่เก็บใน DB
    Hash      string
}

// NewDeviceToken – สร้าง token ใหม่ (raw + bcrypt hash)
func NewDeviceToken() (DeviceToken, error) {
    raw := make([]byte, 32)
    if _, err := rand.Read(raw); err != nil {
        return DeviceToken{}, err
    }
    plain := base64.RawURLEncoding.EncodeToString(raw)
    hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
    if err != nil {
        return DeviceToken{}, err
    }
    return DeviceToken{Plaintext: plain, Hash: string(hash)}, nil
}

// HashOnly – สร้างจาก plaintext ที่ให้มา (ใช้ใน tests)
func HashOnly(plain string) (DeviceToken, error) {
    if len(plain) < 16 {
        return DeviceToken{}, domainerrors.ErrWeakDeviceToken
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
    if err != nil {
        return DeviceToken{}, err
    }
    return DeviceToken{Hash: string(hash)}, nil
}

func (t DeviceToken) Verify(plain string) bool {
    return bcrypt.CompareHashAndPassword([]byte(t.Hash), []byte(plain)) == nil
}

func (t DeviceToken) MQTTPrefix(clientID string) string {
    return "dev-" + strings.ToLower(clientID)
}
```

### `domain/value_object/shadow_state.go` ⭐ (Pattern พิเศษของ module นี้)

```go
package valueobject

import (
    "encoding/json"
    "time"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
)

// ShadowState – Digital Twin ของ device (desired / reported / delta)
// Pattern: AWS IoT Device Shadow / Azure Device Twin
// ใช้ optimistic locking ด้วย Version

type ShadowState struct {
    Desired  map[string]any `json:"desired"`
    Reported map[string]any `json:"reported"`
    Delta    map[string]any `json:"delta,omitempty"`  // คำนวณอัตโนมัติ
    Metadata ShadowMetadata `json:"metadata"`
    Version  int64          `json:"version"`
}

type ShadowMetadata struct {
    Desired  map[string]time.Time `json:"desired,omitempty"`
    Reported map[string]time.Time `json:"reported,omitempty"`
}

func NewShadowState() ShadowState {
    return ShadowState{
        Desired:  map[string]any{},
        Reported: map[string]any{},
        Delta:    map[string]any{},
        Metadata: ShadowMetadata{
            Desired:  map[string]time.Time{},
            Reported: map[string]time.Time{},
        },
        Version: 1,
    }
}

// UpdateDesired – ตั้ง desired state (user → device)
func (s *ShadowState) UpdateDesired(patch map[string]any) error {
    if len(patch) == 0 {
        return domainerrors.ErrEmptyPatch
    }
    now := time.Now()
    for k, v := range patch {
        s.Desired[k] = v
        s.Metadata.Desired[k] = now
    }
    s.recomputeDelta()
    s.Version++
    return nil
}

// UpdateReported – ตั้ง reported state (device → user)
func (s *ShadowState) UpdateReported(patch map[string]any) error {
    if len(patch) == 0 {
        return domainerrors.ErrEmptyPatch
    }
    now := time.Now()
    for k, v := range patch {
        s.Reported[k] = v
        s.Metadata.Reported[k] = now
    }
    s.recomputeDelta()
    s.Version++
    return nil
}

// ClearDelta – เมื่อ device sync desired ครบแล้ว
func (s *ShadowState) ClearDelta(keys []string) {
    for _, k := range keys {
        delete(s.Delta, k)
    }
}

// recomputeDelta – desired vs reported (คำนวณใหม่ทุกครั้ง)
func (s *ShadowState) recomputeDelta() {
    s.Delta = map[string]any{}
    for k, dv := range s.Desired {
        rv, ok := s.Reported[k]
        if !ok || !jsonEqual(dv, rv) {
            s.Delta[k] = dv
        }
    }
}

func jsonEqual(a, b any) bool {
    ab, _ := json.Marshal(a)
    bb, _ := json.Marshal(b)
    return string(ab) == string(bb)
}

func (s ShadowState) HasDelta() bool { return len(s.Delta) > 0 }

func (s ShadowState) IsInSync() bool { return len(s.Delta) == 0 }

// MergeDesired – patch แบบ merge (ไม่ replace)
func (s *ShadowState) MergeDesired(patch map[string]any) {
    for k, v := range patch {
        if sub, ok := v.(map[string]any); ok {
            if existing, ok := s.Desired[k].(map[string]any); ok {
                for sk, sv := range sub {
                    existing[sk] = sv
                }
                continue
            }
        }
        s.Desired[k] = v
    }
    s.recomputeDelta()
    s.Version++
}
```

### `domain/value_object/automation_trigger.go`
```go
package valueobject

import domainerrors "icmongolang/internal/modules/device/domain/errors"

type AutomationTriggerType string

const (
    AutomationTriggerTelemetry  AutomationTriggerType = "TELEMETRY"
    AutomationTriggerSchedule   AutomationTriggerType = "SCHEDULE"
    AutomationTriggerCommandAck AutomationTriggerType = "COMMAND_ACK"
    AutomationTriggerAlert      AutomationTriggerType = "ALERT"
    AutomationTriggerShadow     AutomationTriggerType = "SHADOW_UPDATE"
    AutomationTriggerManual     AutomationTriggerType = "MANUAL"
)

func (t AutomationTriggerType) IsValid() bool {
    switch t {
    case AutomationTriggerTelemetry, AutomationTriggerSchedule,
        AutomationTriggerCommandAck, AutomationTriggerAlert,
        AutomationTriggerShadow, AutomationTriggerManual:
        return true
    }
    return false
}

type AutomationTrigger struct {
    Type      AutomationTriggerType `json:"type"`
    DeviceID  *string               `json:"device_id,omitempty"`
    GroupID   *string               `json:"group_id,omitempty"`
    Metric    *Metric               `json:"metric,omitempty"`
    CronExpr  string                `json:"cron_expr,omitempty"`     // สำหรับ SCHEDULE
    AlertSev  *AlertSeverity        `json:"alert_severity,omitempty"`// สำหรับ ALERT
}

func (t AutomationTrigger) Validate() error {
    if !t.Type.IsValid() {
        return domainerrors.ErrInvalidTrigger
    }
    switch t.Type {
    case AutomationTriggerTelemetry:
        if t.Metric == nil {
            return domainerrors.ErrTriggerMetricRequired
        }
    case AutomationTriggerSchedule:
        if t.CronExpr == "" {
            return domainerrors.ErrTriggerCronRequired
        }
    }
    return nil
}

// Condition – เงื่อนไขสำหรับ evaluate
type Condition struct {
    Field    string `json:"field"`    // "temperature" | "battery_level" | "shadow.reported.fan"
    Operator string `json:"operator"` // ">" | "<" | ">=" | "<=" | "==" | "!=" | "in" | "contains"
    Value    any    `json:"value"`
}

func (c Condition) IsValid() bool {
    switch c.Operator {
    case ">", "<", ">=", "<=", "==", "!=", "in", "contains":
        return true
    }
    return false
}

// Action – สิ่งที่ทำเมื่อ condition ผ่าน
type ActionType string

const (
    ActionTypeSendCommand  ActionType = "SEND_COMMAND"
    ActionTypeSetShadow    ActionType = "SET_SHADOW"
    ActionTypeNotify       ActionType = "NOTIFY"
    ActionTypeWebhook      ActionType = "WEBHOOK"
    ActionTypeBroadcastWS  ActionType = "BROADCAST_WS"
    ActionTypeCreateTicket ActionType = "CREATE_TICKET"
    ActionTypeDelay        ActionType = "DELAY"
)

type Action struct {
    Type     ActionType     `json:"type"`
    DeviceID *string        `json:"device_id,omitempty"`
    Command  string         `json:"command,omitempty"`
    Payload  map[string]any `json:"payload,omitempty"`
    Channel  string         `json:"channel,omitempty"`   // email|line|sms
    Target   string         `json:"target,omitempty"`
    DelayMs  int            `json:"delay_ms,omitempty"`
}

func (a Action) IsValid() bool {
    switch a.Type {
    case ActionTypeSendCommand:
        return a.Command != ""
    case ActionTypeSetShadow:
        return len(a.Payload) > 0
    case ActionTypeNotify:
        return a.Channel != "" && a.Target != ""
    case ActionTypeWebhook:
        return a.Target != ""
    case ActionTypeBroadcastWS:
        return true
    case ActionTypeCreateTicket:
        return true
    case ActionTypeDelay:
        return a.DelayMs > 0
    }
    return false
}
```

---

## A.3 Domain Entities

### `domain/entity/device_model.go`
```go
package entity

import (
    "time"

    "github.com/google/uuid"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

// DeviceModel – รุ่นอุปกรณ์ (master data)
type DeviceModel struct {
    ID           uuid.UUID
    Vendor       string
    ModelNo      string
    Name         string
    DeviceType   valueobject.DeviceType
    Protocol     valueobject.Protocol
    Capabilities valueobject.CapabilitySet
    DefaultConfig map[string]any
    FirmwareChannel string // stable | beta
    CreatedAt    time.Time
}

func (m *DeviceModel) SupportsMetric(metric valueobject.Metric) bool {
    return m.Capabilities.HasMetric(metric)
}

func (m *DeviceModel) SupportsCommand(cmd string) bool {
    return m.Capabilities.HasCommand(cmd)
}
```

### `domain/entity/device.go` ⭐ (Aggregate Root หลัก)
```go
package entity

import (
    "strings"
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

// Device – Aggregate Root
// Invariants:
//   - serial unique ต่อ tenant
//   - provisioned → ต้องมี token hash + mqtt client id
//   - online → last_seen_at ต้อง recent (< 5 min)
//   - decommissioned → terminal state (no transitions)
type Device struct {
    ID              uuid.UUID
    TenantID        uuid.UUID
    CustomerID      uuid.UUID
    SiteID          uuid.UUID
    GroupID         *uuid.UUID
    ModelID         *uuid.UUID

    SerialNo        valueobject.DeviceSerial
    Name            string
    Type            valueobject.DeviceType
    Protocol        valueobject.Protocol
    Status          valueobject.DeviceStatus

    MQTTClientID    string
    DeviceTokenHash string
    FirmwareVersion string
    FirmwareID      *uuid.UUID

    LastSeenAt      *time.Time
    InstalledAt     *time.Time
    ProvisionedAt   *time.Time
    DecommissionedAt *time.Time

    Capabilities    valueobject.CapabilitySet
    Tags            []string
    Metadata        map[string]any

    CreatedBy       uuid.UUID
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

// NewDevice – Factory
func NewDevice(
    tenantID, customerID, siteID, actorID uuid.UUID,
    serial valueobject.DeviceSerial,
    typ valueobject.DeviceType,
    protocol valueobject.Protocol,
    name string,
) (*Device, error) {
    if tenantID == uuid.Nil {
        return nil, domainerrors.ErrInvalidTenantID
    }
    if customerID == uuid.Nil {
        return nil, domainerrors.ErrInvalidCustomerID
    }
    if siteID == uuid.Nil {
        return nil, domainerrors.ErrInvalidSiteID
    }
    if !typ.IsValid() {
        return nil, domainerrors.ErrInvalidDeviceType
    }
    if !protocol.IsValid() {
        return nil, domainerrors.ErrInvalidProtocol
    }
    name = strings.TrimSpace(name)
    if name == "" {
        return nil, domainerrors.ErrInvalidDeviceName
    }

    now := time.Now()
    return &Device{
        ID:         uuid.New(),
        TenantID:   tenantID,
        CustomerID: customerID,
        SiteID:     siteID,
        SerialNo:   serial,
        Name:       name,
        Type:       typ,
        Protocol:   protocol,
        Status:     valueobject.DeviceStatusRegistered,
        Tags:       []string{},
        Metadata:   map[string]any{},
        CreatedBy:  actorID,
        CreatedAt:  now,
        UpdatedAt:  now,
    }, nil
}

// --- Behavior: Provisioning ---

// Provision – สร้าง token + mqtt client id → status = PROVISIONED
func (d *Device) Provision(token valueobject.DeviceToken, clientID string) error {
    if d.Status != valueobject.DeviceStatusRegistered {
        return domainerrors.ErrDeviceAlreadyProvisioned
    }
    if token.Hash == "" {
        return domainerrors.ErrEmptyDeviceToken
    }
    if clientID == "" {
        return domainerrors.ErrEmptyMQTTClientID
    }
    d.DeviceTokenHash = token.Hash
    d.MQTTClientID = clientID
    now := time.Now()
    d.ProvisionedAt = &now
    d.Status = valueobject.DeviceStatusProvisioned
    d.UpdatedAt = now
    return nil
}

// --- Behavior: Connection state ---

func (d *Device) MarkOnline(ts time.Time) error {
    if d.Status == valueobject.DeviceStatusDecommissioned {
        return domainerrors.ErrDeviceDecommissioned
    }
    if d.Status != valueobject.DeviceStatusProvisioned &&
        d.Status != valueobject.DeviceStatusOffline &&
        d.Status != valueobject.DeviceStatusOnline &&
        d.Status != valueobject.DeviceStatusFault {
        return domainerrors.ErrInvalidStatusTransition
    }
    d.LastSeenAt = &ts
    d.Status = valueobject.DeviceStatusOnline
    d.UpdatedAt = ts
    return nil
}

func (d *Device) MarkOffline(reason string) error {
    if d.Status == valueobject.DeviceStatusDecommissioned {
        return nil
    }
    if d.Status == valueobject.DeviceStatusOffline {
        return nil
    }
    d.Status = valueobject.DeviceStatusOffline
    if d.Metadata == nil {
        d.Metadata = map[string]any{}
    }
    d.Metadata["offline_reason"] = reason
    d.Metadata["offline_at"] = time.Now()
    d.UpdatedAt = time.Now()
    return nil
}

func (d *Device) Heartbeat(ts time.Time) error {
    if d.Status == valueobject.DeviceStatusDecommissioned {
        return domainerrors.ErrDeviceDecommissioned
    }
    d.LastSeenAt = &ts
    if d.Status != valueobject.DeviceStatusOnline {
        d.Status = valueobject.DeviceStatusOnline
    }
    d.UpdatedAt = ts
    return nil
}

// --- Behavior: Fault & Maintenance ---

func (d *Device) ReportFault(code, message string) error {
    if d.Status == valueobject.DeviceStatusDecommissioned {
        return domainerrors.ErrDeviceDecommissioned
    }
    d.Status = valueobject.DeviceStatusFault
    if d.Metadata == nil {
        d.Metadata = map[string]any{}
    }
    d.Metadata["fault_code"] = code
    d.Metadata["fault_message"] = message
    d.Metadata["fault_at"] = time.Now()
    d.UpdatedAt = time.Now()
    return nil
}

func (d *Device) ClearFault() error {
    if d.Status != valueobject.DeviceStatusFault {
        return domainerrors.ErrDeviceNotInFault
    }
    d.Status = valueobject.DeviceStatusOffline
    if d.Metadata != nil {
        delete(d.Metadata, "fault_code")
        delete(d.Metadata, "fault_message")
    }
    d.UpdatedAt = time.Now()
    return nil
}

func (d *Device) EnterMaintenance(reason string) error {
    if d.Status == valueobject.DeviceStatusDecommissioned {
        return domainerrors.ErrDeviceDecommissioned
    }
    if !d.Status.CanTransitionTo(valueobject.DeviceStatusMaintenance) {
        return domainerrors.ErrInvalidStatusTransition
    }
    d.Status = valueobject.DeviceStatusMaintenance
    if d.Metadata == nil {
        d.Metadata = map[string]any{}
    }
    d.Metadata["maintenance_reason"] = reason
    d.UpdatedAt = time.Now()
    return nil
}

func (d *Device) ExitMaintenance() error {
    if d.Status != valueobject.DeviceStatusMaintenance {
        return domainerrors.ErrDeviceNotInMaintenance
    }
    d.Status = valueobject.DeviceStatusOffline
    d.UpdatedAt = time.Now()
    return nil
}

// --- Behavior: Capabilities ---

func (d *Device) SetCapabilities(caps valueobject.CapabilitySet) error {
    if err := caps.Validate(); err != nil {
        return err
    }
    d.Capabilities = caps
    d.UpdatedAt = time.Now()
    return nil
}

// --- Behavior: Firmware ---

func (d *Device) UpdateFirmware(version string, firmwareID uuid.UUID) error {
    if version == "" {
        return domainerrors.ErrInvalidFirmwareVersion
    }
    d.FirmwareVersion = version
    d.FirmwareID = &firmwareID
    d.UpdatedAt = time.Now()
    return nil
}

// --- Behavior: Decommission ---

func (d *Device) Decommission(reason string, actorID uuid.UUID) error {
    if d.Status == valueobject.DeviceStatusDecommissioned {
        return domainerrors.ErrDeviceAlreadyDecommissioned
    }
    now := time.Now()
    d.Status = valueobject.DeviceStatusDecommissioned
    d.DecommissionedAt = &now
    if d.Metadata == nil {
        d.Metadata = map[string]any{}
    }
    d.Metadata["decommission_reason"] = reason
    d.Metadata["decommission_by"] = actorID.String()
    d.UpdatedAt = now
    return nil
}

// --- Query methods ---

func (d *Device) IsOnline() bool        { return d.Status == valueobject.DeviceStatusOnline }
func (d *Device) IsOffline() bool       { return d.Status == valueobject.DeviceStatusOffline }
func (d *Device) IsFaulty() bool        { return d.Status == valueobject.DeviceStatusFault }
func (d *Device) IsProvisioned() bool   { return d.DeviceTokenHash != "" }
func (d *Device) IsDecommissioned() bool { return d.Status == valueobject.DeviceStatusDecommissioned }

func (d *Device) CanReceiveCommand() bool {
    return d.Status.CanReceiveCommand()
}

func (d *Device) HasMetricCapability(metric valueobject.Metric) bool {
    return d.Capabilities.HasMetric(metric)
}

func (d *Device) HasCommandCapability(cmd string) bool {
    return d.Capabilities.HasCommand(cmd)
}

// SecondsSinceLastSeen – สำหรับ offline detection
func (d *Device) SecondsSinceLastSeen() int64 {
    if d.LastSeenAt == nil {
        return -1
    }
    return int64(time.Since(*d.LastSeenAt).Seconds())
}

func (d *Device) IsStale(threshold time.Duration) bool {
    if d.LastSeenAt == nil {
        return true
    }
    return time.Since(*d.LastSeenAt) > threshold
}
```

### `domain/entity/device_shadow.go` ⭐ (Digital Twin Aggregate)
```go
package entity

import (
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

// DeviceShadow – Aggregate Root (Digital Twin)
// ใช้ optimistic locking ด้วย Version
type DeviceShadow struct {
    DeviceID  uuid.UUID
    TenantID  uuid.UUID
    State     valueobject.ShadowState
    CreatedAt time.Time
    UpdatedAt time.Time
}

func NewDeviceShadow(tenantID, deviceID uuid.UUID) *DeviceShadow {
    now := time.Now()
    return &DeviceShadow{
        DeviceID:  deviceID,
        TenantID:  tenantID,
        State:     valueobject.NewShadowState(),
        CreatedAt: now,
        UpdatedAt: now,
    }
}

// SetDesired – user ตั้ง desired state
func (s *DeviceShadow) SetDesired(patch map[string]any) error {
    if err := s.State.UpdateDesired(patch); err != nil {
        return err
    }
    s.UpdatedAt = time.Now()
    return nil
}

// MergeDesired – patch แบบ merge (nested)
func (s *DeviceShadow) MergeDesired(patch map[string]any) error {
    if len(patch) == 0 {
        return domainerrors.ErrEmptyPatch
    }
    s.State.MergeDesired(patch)
    s.UpdatedAt = time.Now()
    return nil
}

// UpdateReported – device ส่ง reported state มา
func (s *DeviceShadow) UpdateReported(patch map[string]any) error {
    if err := s.State.UpdateReported(patch); err != nil {
        return err
    }
    s.UpdatedAt = time.Now()
    return nil
}

// GetDelta – ค่าที่ต้อง sync ลง device
func (s *DeviceShadow) GetDelta() map[string]any {
    return s.State.Delta
}

// ClearDelta – เมื่อ device sync ครบแล้ว
func (s *DeviceShadow) ClearDelta(keys []string) {
    s.State.ClearDelta(keys)
    s.UpdatedAt = time.Now()
}

// OptimisticLock – ตรวจ version ก่อน save
func (s *DeviceShadow) CheckVersion(expected int64) error {
    if s.State.Version != expected {
        return domainerrors.ErrShadowVersionMismatch
    }
    return nil
}

func (s *DeviceShadow) IsInSync() bool { return s.State.IsInSync() }
func (s *DeviceShadow) HasDelta() bool { return s.State.HasDelta() }
func (s *DeviceShadow) Version() int64 { return s.State.Version }
```

### `domain/entity/command.go`
```go
package entity

import (
    "strings"
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

// Command – Aggregate Root
// Invariants:
//   - command name ต้องไม่ว่าง
//   - PENDING → SENT → ACKED/FAILED/TIMEOUT
//   - timeout: 30s (configurable)
type Command struct {
    ID          uuid.UUID
    TenantID    uuid.UUID
    DeviceID    uuid.UUID

    Command     string
    Payload     map[string]any
    Status      valueobject.CommandStatus
    Priority    int // 1=low, 5=high

    IssuedBy    uuid.UUID
    IssuedAt    time.Time
    SentAt      *time.Time
    AckedAt     *time.Time
    ExpiresAt   time.Time

    ErrorMsg    string
    ResultPayload map[string]any
    RetryCount  int
    MaxRetries  int
}

const DefaultCommandTimeout = 30 * time.Second

func NewCommand(
    tenantID, deviceID, actorID uuid.UUID,
    command string,
    payload map[string]any,
) (*Command, error) {
    command = strings.TrimSpace(command)
    if command == "" {
        return nil, domainerrors.ErrInvalidCommand
    }
    if len(command) > 100 {
        return nil, domainerrors.ErrCommandTooLong
    }
    now := time.Now()
    return &Command{
        ID:         uuid.New(),
        TenantID:   tenantID,
        DeviceID:   deviceID,
        Command:    command,
        Payload:    payload,
        Status:     valueobject.CommandStatusPending,
        Priority:   3,
        IssuedBy:   actorID,
        IssuedAt:   now,
        ExpiresAt:  now.Add(DefaultCommandTimeout),
        MaxRetries: 2,
    }, nil
}

func (c *Command) Queue(reason string) error {
    if c.Status != valueobject.CommandStatusPending {
        return domainerrors.ErrInvalidCommandState
    }
    c.Status = valueobject.CommandStatusQueued
    c.ErrorMsg = reason // "device_offline"
    return nil
}

func (c *Command) MarkSent(ts time.Time) error {
    if c.Status != valueobject.CommandStatusPending &&
        c.Status != valueobject.CommandStatusQueued {
        return domainerrors.ErrInvalidCommandState
    }
    c.Status = valueobject.CommandStatusSent
    c.SentAt = &ts
    return nil
}

// Ack – device ตอบรับ (สำเร็จ)
func (c *Command) Ack(result map[string]any) error {
    if !c.Status.IsPending() {
        return domainerrors.ErrCommandAlreadyTerminal
    }
    if time.Now().After(c.ExpiresAt) {
        return c.Timeout()
    }
    now := time.Now()
    c.Status = valueobject.CommandStatusAcked
    c.AckedAt = &now
    c.ResultPayload = result
    return nil
}

// Nack – device ปฏิเสธ (ล้มเหลว)
func (c *Command) Nack(errMsg string) error {
    if !c.Status.IsPending() {
        return domainerrors.ErrCommandAlreadyTerminal
    }
    now := time.Now()
    c.Status = valueobject.CommandStatusFailed
    c.AckedAt = &now
    c.ErrorMsg = errMsg
    return nil
}

func (c *Command) Timeout() error {
    if !c.Status.IsPending() {
        return domainerrors.ErrCommandAlreadyTerminal
    }
    now := time.Now()
    c.Status = valueobject.CommandStatusTimeout
    c.AckedAt = &now
    c.ErrorMsg = "timeout"
    return nil
}

func (c *Command) Cancel(reason string) error {
    if !c.Status.IsPending() {
        return domainerrors.ErrCommandAlreadyTerminal
    }
    now := time.Now()
    c.Status = valueobject.CommandStatusCancelled
    c.AckedAt = &now
    c.ErrorMsg = reason
    return nil
}

func (c *Command) CanRetry() bool {
    return c.RetryCount < c.MaxRetries &&
        (c.Status == valueobject.CommandStatusFailed || c.Status == valueobject.CommandStatusTimeout)
}

func (c *Command) IncrementRetry() {
    c.RetryCount++
}

func (c *Command) IsExpired(asOf time.Time) bool {
    return asOf.After(c.ExpiresAt) && c.Status.IsPending()
}

func (c *Command) IsTerminal() bool { return c.Status.IsTerminal() }
```

### `domain/entity/alert_rule.go`
```go
package entity

import (
    "strings"
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

// AlertRule – Aggregate Root
type AlertRule struct {
    ID         uuid.UUID
    TenantID   uuid.UUID
    Name       string
    DeviceID   *uuid.UUID
    GroupID    *uuid.UUID
    SiteID     *uuid.UUID

    Metric     valueobject.Metric
    Operator   string   // ">" | "<" | ">=" | "<=" | "==" | "!=" | "change" | "delta"
    Threshold  float64
    DurationSec int     // ต้องค้างอยู่นานเท่าไหร่ (กัน flapping)

    Severity   valueobject.AlertSeverity
    Actions    []AlertAction
    IsActive   bool
    CooldownSec int    // หลัง trigger แล้ว รอ cooldown ก่อน trigger อีก

    LastTriggeredAt *time.Time
    CreatedBy  uuid.UUID
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

type AlertAction struct {
    Type     string         `json:"type"`    // NOTIFY | WEBHOOK | COMMAND | TICKET
    Channel  string         `json:"channel,omitempty"`
    Target   string         `json:"target,omitempty"`
    Payload  map[string]any `json:"payload,omitempty"`
}

func NewAlertRule(
    tenantID uuid.UUID,
    name string,
    metric valueobject.Metric,
    operator string,
    threshold float64,
    severity valueobject.AlertSeverity,
) (*AlertRule, error) {
    name = strings.TrimSpace(name)
    if name == "" {
        return nil, domainerrors.ErrInvalidAlertRuleName
    }
    if !metric.IsValid() {
        return nil, domainerrors.ErrInvalidMetric
    }
    if !isValidOperator(operator) {
        return nil, domainerrors.ErrInvalidOperator
    }
    if !severity.IsValid() {
        return nil, domainerrors.ErrInvalidSeverity
    }
    now := time.Now()
    return &AlertRule{
        ID:          uuid.New(),
        TenantID:    tenantID,
        Name:        name,
        Metric:      metric,
        Operator:    operator,
        Threshold:   threshold,
        Severity:    severity,
        DurationSec: 0,
        CooldownSec: 300,
        IsActive:    true,
        Actions:     []AlertAction{},
        CreatedAt:   now,
        UpdatedAt:   now,
    }, nil
}

// BindToDevice – ผูก rule กับ device
func (r *AlertRule) BindToDevice(deviceID uuid.UUID) {
    r.DeviceID = &deviceID
    r.UpdatedAt = time.Now()
}

func (r *AlertRule) BindToGroup(groupID uuid.UUID) {
    r.GroupID = &groupID
    r.UpdatedAt = time.Now()
}

// Evaluate – ตรวจว่าค่าที่เข้ามา trigger หรือไม่
func (r *AlertRule) Evaluate(value float64) bool {
    if !r.IsActive {
        return false
    }
    switch r.Operator {
    case ">":      return value > r.Threshold
    case "<":      return value < r.Threshold
    case ">=":     return value >= r.Threshold
    case "<=":     return value <= r.Threshold
    case "==":     return value == r.Threshold
    case "!=":     return value != r.Threshold
    case "change": return true // ตรวจ diff ที่ infrastructure
    case "delta":  return true // ตรวจ delta ที่ infrastructure
    }
    return false
}

// IsInCooldown – ตรวจว่ายังอยู่ในช่วง cooldown หรือไม่
func (r *AlertRule) IsInCooldown(asOf time.Time) bool {
    if r.LastTriggeredAt == nil {
        return false
    }
    return asOf.Sub(*r.LastTriggeredAt) < time.Duration(r.CooldownSec)*time.Second
}

func (r *AlertRule) MarkTriggered(asOf time.Time) {
    r.LastTriggeredAt = &asOf
    r.UpdatedAt = asOf
}

func (r *AlertRule) Disable() {
    r.IsActive = false
    r.UpdatedAt = time.Now()
}

func (r *AlertRule) Enable() {
    r.IsActive = true
    r.UpdatedAt = time.Now()
}

func (r *AlertRule) MatchesDevice(deviceID uuid.UUID) bool {
    if r.DeviceID != nil && *r.DeviceID == deviceID {
        return true
    }
    return false
}

func isValidOperator(op string) bool {
    switch op {
    case ">", "<", ">=", "<=", "==", "!=", "change", "delta":
        return true
    }
    return false
}
```

### `domain/entity/alert_event.go`
```go
package entity

import (
    "time"

    "github.com/google/uuid"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

// AlertEvent – บันทึกเหตุการณ์แจ้งเตือน (event-sourced)
type AlertEvent struct {
    ID            int64
    TenantID      uuid.UUID
    RuleID        uuid.UUID
    DeviceID      uuid.UUID
    Metric        valueobject.Metric
    Value         float64
    Threshold     float64
    Operator      string
    Severity      valueobject.AlertSeverity
    Message       string
    Acknowledged  bool
    AckedBy       *uuid.UUID
    AckedAt       *time.Time
    ResolvedAt    *time.Time
    TriggeredAt   time.Time
}

func (e *AlertEvent) Acknowledge(actor uuid.UUID) {
    if e.Acknowledged {
        return
    }
    now := time.Now()
    e.Acknowledged = true
    e.AckedBy = &actor
    e.AckedAt = &now
}

func (e *AlertEvent) Resolve() {
    now := time.Now()
    e.ResolvedAt = &now
}
```

### `domain/entity/automation.go` ⭐ (Event-Condition-Action)
```go
package entity

import (
    "strings"
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

// Automation – Aggregate Root
// Pattern: ECA (Event-Condition-Action)
type Automation struct {
    ID          uuid.UUID
    TenantID    uuid.UUID
    Name        string
    Description string

    Trigger     valueobject.AutomationTrigger
    Conditions  []valueobject.Condition
    Actions     []valueobject.Action

    IsActive    bool
    Priority    int    // ยิ่งสูง ยิ่งรันก่อน
    RunCount    int64
    LastRunAt   *time.Time
    LastError   string

    CreatedBy   uuid.UUID
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

const maxActionsPerAutomation = 10
const maxConditionsPerAutomation = 10

func NewAutomation(
    tenantID, actorID uuid.UUID,
    name string,
    trigger valueobject.AutomationTrigger,
) (*Automation, error) {
    name = strings.TrimSpace(name)
    if name == "" {
        return nil, domainerrors.ErrInvalidAutomationName
    }
    if err := trigger.Validate(); err != nil {
        return nil, err
    }
    now := time.Now()
    return &Automation{
        ID:         uuid.New(),
        TenantID:   tenantID,
        Name:       name,
        Trigger:    trigger,
        Conditions: []valueobject.Condition{},
        Actions:    []valueobject.Action{},
        IsActive:   true,
        Priority:   5,
        CreatedBy:  actorID,
        CreatedAt:  now,
        UpdatedAt:  now,
    }, nil
}

func (a *Automation) AddCondition(c valueobject.Condition) error {
    if len(a.Conditions) >= maxConditionsPerAutomation {
        return domainerrors.ErrTooManyConditions
    }
    if !c.IsValid() {
        return domainerrors.ErrInvalidCondition
    }
    a.Conditions = append(a.Conditions, c)
    a.UpdatedAt = time.Now()
    return nil
}

func (a *Automation) AddAction(act valueobject.Action) error {
    if len(a.Actions) >= maxActionsPerAutomation {
        return domainerrors.ErrTooManyActions
    }
    if !act.IsValid() {
        return domainerrors.ErrInvalidAction
    }
    a.Actions = append(a.Actions, act)
    a.UpdatedAt = time.Now()
    return nil
}

// CanRun – ตรวจเงื่อนไขทั้งหมด (AND)
func (a *Automation) CanRun(facts map[string]any) bool {
    if !a.IsActive {
        return false
    }
    for _, c := range a.Conditions {
        if !evaluateCondition(c, facts) {
            return false
        }
    }
    return true
}

func (a *Automation) MarkRun(ts time.Time, err error) {
    a.LastRunAt = &ts
    a.RunCount++
    if err != nil {
        a.LastError = err.Error()
    } else {
        a.LastError = ""
    }
    a.UpdatedAt = ts
}

func (a *Automation) Disable() {
    a.IsActive = false
    a.UpdatedAt = time.Now()
}

func (a *Automation) Enable() {
    a.IsActive = true
    a.UpdatedAt = time.Now()
}

func evaluateCondition(c valueobject.Condition, facts map[string]any) bool {
    v, ok := facts[c.Field]
    if !ok {
        return false
    }
    switch c.Operator {
    case "==":
        return valuesEqual(v, c.Value)
    case "!=":
        return !valuesEqual(v, c.Value)
    case ">":
        return toFloat(v) > toFloat(c.Value)
    case "<":
        return toFloat(v) < toFloat(c.Value)
    case ">=":
        return toFloat(v) >= toFloat(c.Value)
    case "<=":
        return toFloat(v) <= toFloat(c.Value)
    case "in":
        if arr, ok := c.Value.([]any); ok {
            for _, item := range arr {
                if valuesEqual(v, item) {
                    return true
                }
            }
        }
        return false
    case "contains":
        s, ok := v.(string)
        if !ok { return false }
        sub, ok := c.Value.(string)
        return ok && strings.Contains(s, sub)
    }
    return false
}

func valuesEqual(a, b any) bool { return toFloat(a) == toFloat(b) && a == b }

func toFloat(v any) float64 {
    switch n := v.(type) {
    case float64:
        return n
    case float32:
        return float64(n)
    case int:
        return float64(n)
    case int64:
        return float64(n)
    case int32:
        return float64(n)
    }
    return 0
}
```

### `domain/entity/device_group.go`
```go
package entity

import (
    "strings"
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
)

// DeviceGroup – Aggregate Root (zone/floor/barn)
type DeviceGroup struct {
    ID        uuid.UUID
    TenantID  uuid.UUID
    SiteID    uuid.UUID
    ParentID  *uuid.UUID
    Name      string
    GroupType string // ZONE | FLOOR | BARN | ROOM | AREA
    Metadata  map[string]any
    CreatedAt time.Time
    UpdatedAt time.Time
}

func NewDeviceGroup(tenantID, siteID uuid.UUID, name, groupType string) (*DeviceGroup, error) {
    name = strings.TrimSpace(name)
    if name == "" {
        return nil, domainerrors.ErrInvalidGroupName
    }
    if groupType == "" {
        groupType = "ZONE"
    }
    now := time.Now()
    return &DeviceGroup{
        ID:        uuid.New(),
        TenantID:  tenantID,
        SiteID:    siteID,
        Name:      name,
        GroupType: groupType,
        Metadata:  map[string]any{},
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}

func (g *DeviceGroup) SetParent(parentID uuid.UUID) error {
    if parentID == g.ID {
        return domainerrors.ErrSelfReference
    }
    g.ParentID = &parentID
    g.UpdatedAt = time.Now()
    return nil
}
```

### `domain/entity/firmware.go`
```go
package entity

import (
    "time"

    "github.com/google/uuid"
)

type FirmwareChannel string

const (
    FirmwareChannelStable FirmwareChannel = "stable"
    FirmwareChannelBeta   FirmwareChannel = "beta"
)

// Firmware – Entity
type Firmware struct {
    ID          uuid.UUID
    ModelID     uuid.UUID
    Version     string
    Channel     FirmwareChannel
    FileURL     string
    FileSize    int64
    Checksum    string   // SHA-256
    ReleaseNotes string
    MinVersion  string   // version ขั้นต่ำที่ upgrade ได้
    IsActive    bool
    ReleasedAt  time.Time
    CreatedAt   time.Time
}
```

### `domain/entity/ai_model.go`
```go
package entity

import (
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
)

type AIModelType string

const (
    AIModelTypeAnomaly  AIModelType = "ANOMALY"
    AIModelTypeForecast AIModelType = "FORECAST"
    AIModelTypeClassify AIModelType = "CLASSIFY"
    AIModelTypeNL       AIModelType = "NL"  // natural language
)

type AIProvider string

const (
    AIProviderOllama AIProvider = "OLLAMA"
    AIProviderOpenAI AIProvider = "OPENAI"
    AIProviderCustom AIProvider = "CUSTOM"
)

// AIModel – Entity (metadata ของ AI model)
type AIModel struct {
    ID         uuid.UUID
    TenantID   *uuid.UUID // nil = shared platform model
    Name       string
    Type       AIModelType
    Provider   AIProvider
    Endpoint   string
    ModelKey   string   // เช่น "llama3:8b"
    Config     map[string]any
    IsActive   bool
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

func NewAIModel(name string, typ AIModelType, provider AIProvider) (*AIModel, error) {
    if name == "" {
        return nil, domainerrors.ErrInvalidAIModelName
    }
    now := time.Now()
    return &AIModel{
        ID:        uuid.New(),
        Name:      name,
        Type:      typ,
        Provider:  provider,
        Config:    map[string]any{},
        IsActive:  true,
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}
```

### `domain/entity/telemetry.go`
```go
package entity

import (
    "time"

    "github.com/google/uuid"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

// Telemetry – event-sourced, ไม่ mutate aggregate
// เก็บใน InfluxDB (ไม่ใช่ PostgreSQL)
type Telemetry struct {
    TenantID  uuid.UUID
    DeviceID  uuid.UUID
    SiteID    uuid.UUID
    Metrics   []valueobject.MetricValue
    Timestamp time.Time
    Source    string // mqtt | http | backfill
}

// TelemetryBatch – สำหรับ batch write
type TelemetryBatch struct {
    Items []*Telemetry
}

func (b *TelemetryBatch) Add(t *Telemetry) {
    b.Items = append(b.Items, t)
}

func (b *TelemetryBatch) Len() int { return len(b.Items) }

func (b *TelemetryBatch) Reset() { b.Items = b.Items[:0] }
```

---

## A.4 Repository Interfaces

### `domain/repository/device_repository.go`
```go
package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "icmongolang/internal/modules/device/domain/entity"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

type DeviceFilter struct {
    TenantID   uuid.UUID
    CustomerID *uuid.UUID
    SiteID     *uuid.UUID
    GroupID    *uuid.UUID
    ModelID    *uuid.UUID
    Statuses   []valueobject.DeviceStatus
    Types      []valueobject.DeviceType
    Protocol   *valueobject.Protocol
    Search     string
    Tags       []string
    Page       int
    PageSize   int
    SortBy     string
    SortOrder  string
}

type DeviceRepository interface {
    Save(ctx context.Context, d *entity.Device) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Device, error)
    FindBySerial(ctx context.Context, tenantID uuid.UUID, serial valueobject.DeviceSerial) (*entity.Device, error)
    FindByMQTTClientID(ctx context.Context, clientID string) (*entity.Device, error)
    List(ctx context.Context, f DeviceFilter) ([]*entity.Device, int64, error)
    ListBySite(ctx context.Context, tenantID, siteID uuid.UUID) ([]*entity.Device, error)
    ListByGroup(ctx context.Context, tenantID, groupID uuid.UUID) ([]*entity.Device, error)
    ListStale(ctx context.Context, threshold time.Time) ([]*entity.Device, error)
    CountByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) (int64, error)
    CountByStatus(ctx context.Context, tenantID uuid.UUID) (map[valueobject.DeviceStatus]int64, error)
    Delete(ctx context.Context, tenantID, id uuid.UUID) error
}
```

### `domain/repository/shadow_repository.go`
```go
package repository

import (
    "context"

    "github.com/google/uuid"
    "icmongolang/internal/modules/device/domain/entity"
)

type ShadowRepository interface {
    Save(ctx context.Context, s *entity.DeviceShadow) error
    FindByDeviceID(ctx context.Context, tenantID, deviceID uuid.UUID) (*entity.DeviceShadow, error)
    Delete(ctx context.Context, tenantID, deviceID uuid.UUID) error
    SaveWithVersion(ctx context.Context, s *entity.DeviceShadow, expectedVersion int64) error
}
```

### `domain/repository/command_repository.go`
```go
package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "icmongolang/internal/modules/device/domain/entity"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

type CommandFilter struct {
    TenantID uuid.UUID
    DeviceID *uuid.UUID
    Status   []valueobject.CommandStatus
    From     *time.Time
    To       *time.Time
    Page     int
    PageSize int
}

type CommandRepository interface {
    Save(ctx context.Context, c *entity.Command) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Command, error)
    List(ctx context.Context, f CommandFilter) ([]*entity.Command, int64, error)
    ListPending(ctx context.Context, deviceID uuid.UUID) ([]*entity.Command, error)
    ListExpired(ctx context.Context, asOf time.Time, limit int) ([]*entity.Command, error)
    ListQueuedByDevice(ctx context.Context, deviceID uuid.UUID) ([]*entity.Command, error)
}
```

### `domain/repository/telemetry_repository.go`
```go
package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "icmongolang/internal/modules/device/domain/entity"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

type TelemetryQuery struct {
    TenantID  uuid.UUID
    DeviceID  uuid.UUID
    Metrics   []valueobject.Metric
    From      time.Time
    To        time.Time
    Aggregate string // raw | 1m | 5m | 1h | 1d
    Limit     int
}

type TelemetryAggregate struct {
    Metric   valueobject.Metric
    Timestamp time.Time
    Min, Max, Avg, Sum, Count float64
}

type TelemetryRepository interface {
    Write(ctx context.Context, items ...*entity.Telemetry) error
    WriteBatch(ctx context.Context, batch *entity.TelemetryBatch) error
    Query(ctx context.Context, q TelemetryQuery) ([]*entity.Telemetry, error)
    QueryAggregated(ctx context.Context, q TelemetryQuery) ([]*TelemetryAggregate, error)
    LatestByDevice(ctx context.Context, tenantID, deviceID uuid.UUID) (map[valueobject.Metric]float64, error)
    DeleteOlderThan(ctx context.Context, asOf time.Time) error
}
```

### `domain/repository/alert_rule_repository.go`
```go
package repository

import (
    "context"

    "github.com/google/uuid"
    "icmongolang/internal/modules/device/domain/entity"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

type AlertRuleRepository interface {
    Save(ctx context.Context, r *entity.AlertRule) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.AlertRule, error)
    ListActive(ctx context.Context, tenantID uuid.UUID) ([]*entity.AlertRule, error)
    ListActiveByDevice(ctx context.Context, tenantID, deviceID uuid.UUID) ([]*entity.AlertRule, error)
    ListActiveByGroup(ctx context.Context, tenantID, groupID uuid.UUID) ([]*entity.AlertRule, error)
    ListByMetric(ctx context.Context, tenantID uuid.UUID, metric valueobject.Metric) ([]*entity.AlertRule, error)
    Delete(ctx context.Context, tenantID, id uuid.UUID) error
}
```

### `domain/repository/alert_event_repository.go`
```go
package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "icmongolang/internal/modules/device/domain/entity"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

type AlertEventFilter struct {
    TenantID   uuid.UUID
    DeviceID   *uuid.UUID
    RuleID     *uuid.UUID
    Severity   []valueobject.AlertSeverity
    Acknowledged *bool
    From       *time.Time
    To         *time.Time
    Page       int
    PageSize   int
}

type AlertEventRepository interface {
    Save(ctx context.Context, e *entity.AlertEvent) error
    FindByID(ctx context.Context, id int64) (*entity.AlertEvent, error)
    List(ctx context.Context, f AlertEventFilter) ([]*entity.AlertEvent, int64, error)
    CountUnacknowledged(ctx context.Context, tenantID uuid.UUID) (int64, error)
}
```

### `domain/repository/automation_repository.go`
```go
package repository

import (
    "context"

    "github.com/google/uuid"
    "icmongolang/internal/modules/device/domain/entity"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

type AutomationRepository interface {
    Save(ctx context.Context, a *entity.Automation) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Automation, error)
    ListActiveByTrigger(ctx context.Context, tenantID uuid.UUID, triggerType valueobject.AutomationTriggerType) ([]*entity.Automation, error)
    ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.Automation, error)
    Delete(ctx context.Context, tenantID, id uuid.UUID) error
}
```

### `domain/repository/device_group_repository.go`
```go
package repository

import (
    "context"

    "github.com/google/uuid"
    "icmongolang/internal/modules/device/domain/entity"
)

type DeviceGroupRepository interface {
    Save(ctx context.Context, g *entity.DeviceGroup) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.DeviceGroup, error)
    ListBySite(ctx context.Context, tenantID, siteID uuid.UUID) ([]*entity.DeviceGroup, error)
    Delete(ctx context.Context, tenantID, id uuid.UUID) error
}
```

### `domain/repository/device_model_repository.go`
```go
package repository

import (
    "context"

    "github.com/google/uuid"
    "icmongolang/internal/modules/device/domain/entity"
)

type DeviceModelRepository interface {
    Save(ctx context.Context, m *entity.DeviceModel) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.DeviceModel, error)
    FindByModelNo(ctx context.Context, modelNo string) (*entity.DeviceModel, error)
    List(ctx context.Context) ([]*entity.DeviceModel, error)
}
```

### `domain/repository/firmware_repository.go`
```go
package repository

import (
    "context"

    "github.com/google/uuid"
    "icmongolang/internal/modules/device/domain/entity"
)

type FirmwareRepository interface {
    Save(ctx context.Context, f *entity.Firmware) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Firmware, error)
    FindLatest(ctx context.Context, modelID uuid.UUID, channel entity.FirmwareChannel) (*entity.Firmware, error)
    ListByModel(ctx context.Context, modelID uuid.UUID) ([]*entity.Firmware, error)
}
```

---

## A.5 Domain Services

### `domain/service/provisioning_service.go`
```go
package service

import (
    "context"
    "fmt"
    "strings"

    "github.com/google/uuid"

    "icmongolang/internal/modules/device/domain/entity"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
    "icmongolang/internal/modules/device/domain/repository"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

// ProvisioningService – สร้าง token + MQTT client id + ผูก model capabilities
type ProvisioningService struct {
    modelRepo repository.DeviceModelRepository
    shadowRepo repository.ShadowRepository
}

func NewProvisioningService(
    modelRepo repository.DeviceModelRepository,
    shadowRepo repository.ShadowRepository,
) *ProvisioningService {
    return &ProvisioningService{modelRepo: modelRepo, shadowRepo: shadowRepo}
}

type ProvisionResult struct {
    Token        valueobject.DeviceToken
    MQTTClientID string
    TopicPrefix  string
}

func (s *ProvisioningService) Provision(
    ctx context.Context,
    d *entity.Device,
) (*ProvisionResult, error) {
    if d.IsProvisioned() {
        return nil, domainerrors.ErrDeviceAlreadyProvisioned
    }

    // 1. Generate token
    token, err := valueobject.NewDeviceToken()
    if err != nil {
        return nil, err
    }

    // 2. Generate MQTT client id
    clientID := generateClientID(d)

    // 3. ถ้ามี model → copy capabilities
    if d.ModelID != nil {
        model, err := s.modelRepo.FindByID(ctx, *d.ModelID)
        if err == nil && model != nil {
            _ = d.SetCapabilities(model.Capabilities)
        }
    }

    // 4. Provision aggregate
    if err := d.Provision(token, clientID); err != nil {
        return nil, err
    }

    // 5. Create initial shadow
    shadow := entity.NewDeviceShadow(d.TenantID, d.ID)
    if err := s.shadowRepo.Save(ctx, shadow); err != nil {
        return nil, err
    }

    return &ProvisionResult{
        Token:        token,
        MQTTClientID: clientID,
        TopicPrefix:  topicPrefix(d),
    }, nil
}

func generateClientID(d *entity.Device) string {
    // dev-{tenant8}-{serial8}
    return fmt.Sprintf("dev-%s-%s",
        strings.ToLower(d.TenantID.String()[:8]),
        strings.ToLower(string(d.SerialNo)[:min(8, len(d.SerialNo))]),
    )
}

func topicPrefix(d *entity.Device) string {
    return fmt.Sprintf("iot/%s/%s", d.TenantID.String(), d.ID.String())
}

func min(a, b int) int {
    if a < b { return a }
    return b
}
```

### `domain/service/health_monitor_service.go`
```go
package service

import (
    "time"

    "icmongolang/internal/modules/device/domain/entity"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

// HealthMonitorService – ตัดสิน online/offline จาก last_seen
type HealthMonitorService struct {
    OfflineThreshold   time.Duration
    FaultThreshold     time.Duration
}

func NewHealthMonitorService(offline, fault time.Duration) *HealthMonitorService {
    return &HealthMonitorService{OfflineThreshold: offline, FaultThreshold: fault}
}

type HealthStatus struct {
    DeviceID   string
    Status     valueobject.DeviceStatus
    LastSeen   *time.Time
    ShouldMark bool
    Reason     string
}

func (s *HealthMonitorService) Evaluate(d *entity.Device, asOf time.Time) HealthStatus {
    res := HealthStatus{DeviceID: d.ID.String(), LastSeen: d.LastSeenAt}

    if d.IsDecommissioned() || d.Status == valueobject.DeviceStatusMaintenance {
        res.Status = d.Status
        return res
    }

    if d.LastSeenAt == nil {
        res.Status = d.Status
        res.ShouldMark = d.Status == valueobject.DeviceStatusOnline
        res.Reason = "never_seen"
        return res
    }

    elapsed := asOf.Sub(*d.LastSeenAt)
    if elapsed > s.OfflineThreshold {
        res.Status = valueobject.DeviceStatusOffline
        res.ShouldMark = d.Status != valueobject.DeviceStatusOffline
        res.Reason = "heartbeat_timeout"
        return res
    }

    res.Status = valueobject.DeviceStatusOnline
    res.ShouldMark = d.Status != valueobject.DeviceStatusOnline && d.Status != valueobject.DeviceStatusFault
    res.Reason = "heartbeat_ok"
    return res
}
```

### `domain/service/shadow_merge_service.go`
```go
package service

import (
    "icmongolang/internal/modules/device/domain/entity"
)

// ShadowMergeService – merge reported state จาก device เข้า shadow
type ShadowMergeService struct{}

func NewShadowMergeService() *ShadowMergeService { return &ShadowMergeService{} }

// ApplyReported – เมื่อ device ส่ง reported state มา
func (s *ShadowMergeService) ApplyReported(shadow *entity.DeviceShadow, reported map[string]any) (deltaCleared []string, err error) {
    // เก็บ keys ที่จะ sync ครบ
    beforeDelta := copyMap(shadow.State.Delta)

    if err := shadow.UpdateReported(reported); err != nil {
        return nil, err
    }

    // เทียบ before/after – keys ที่ delta หายไป = sync แล้ว
    for k := range beforeDelta {
        if _, still := shadow.State.Delta[k]; !still {
            deltaCleared = append(deltaCleared, k)
        }
    }
    return deltaCleared, nil
}

func copyMap(m map[string]any) map[string]any {
    out := make(map[string]any, len(m))
    for k, v := range m { out[k] = v }
    return out
}
```

### `domain/service/alert_evaluator.go`
```go
package service

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"

    "icmongolang/internal/modules/device/domain/entity"
    "icmongolang/internal/modules/device/domain/repository"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

// AlertEvaluator – ตรวจ alert rules กับ telemetry ที่เข้ามา
type AlertEvaluator struct {
    ruleRepo  repository.AlertRuleRepository
    eventRepo repository.AlertEventRepository
}

func NewAlertEvaluator(
    ruleRepo repository.AlertRuleRepository,
    eventRepo repository.AlertEventRepository,
) *AlertEvaluator {
    return &AlertEvaluator{ruleRepo: ruleRepo, eventRepo: eventRepo}
}

type TriggeredAlert struct {
    Rule    *entity.AlertRule
    Event   *entity.AlertEvent
    Message string
}

// Evaluate – ตรวจ rules ที่ active กับ telemetry
func (e *AlertEvaluator) Evaluate(
    ctx context.Context,
    t *entity.Telemetry,
) ([]*TriggeredAlert, error) {
    rules, err := e.ruleRepo.ListActiveByDevice(ctx, t.TenantID, t.DeviceID)
    if err != nil {
        return nil, err
    }
    if len(rules) == 0 {
        return nil, nil
    }

    now := time.Now()
    var triggered []*TriggeredAlert

    for _, rule := range rules {
        if rule.IsInCooldown(now) {
            continue
        }
        for _, mv := range t.Metrics {
            if mv.Metric != rule.Metric {
                continue
            }
            if !rule.Evaluate(mv.Value) {
                continue
            }
            // Create alert event
            evt := &entity.AlertEvent{
                TenantID:    t.TenantID,
                RuleID:      rule.ID,
                DeviceID:    t.DeviceID,
                Metric:      mv.Metric,
                Value:       mv.Value,
                Threshold:   rule.Threshold,
                Operator:    rule.Operator,
                Severity:    rule.Severity,
                Message:     formatMessage(rule, mv),
                TriggeredAt: now,
            }
            if err := e.eventRepo.Save(ctx, evt); err != nil {
                continue
            }
            rule.MarkTriggered(now)
            _ = e.ruleRepo.Save(ctx, rule)

            triggered = append(triggered, &TriggeredAlert{
                Rule:    rule,
                Event:   evt,
                Message: evt.Message,
            })
        }
    }
    return triggered, nil
}

func formatMessage(rule *entity.AlertRule, mv valueobject.MetricValue) string {
    return fmt.Sprintf("[%s] %s = %.2f %s (threshold: %s %.2f)",
        rule.Severity, mv.Metric, mv.Value, mv.Metric.Unit(),
        rule.Operator, rule.Threshold)
}

var _ = uuid.Nil
```

### `domain/service/automation_engine.go`
```go
package service

import (
    "context"
    "log"

    "icmongolang/internal/modules/device/domain/entity"
    "icmongolang/internal/modules/device/domain/repository"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

// AutomationEngine – รัน ECA rules
type AutomationEngine struct {
    repo        repository.AutomationRepository
    actionExec  ActionExecutor
}

// ActionExecutor – outbound port ที่ infrastructure implement
type ActionExecutor interface {
    Execute(ctx context.Context, tenantID string, act valueobject.Action) error
}

func NewAutomationEngine(
    repo repository.AutomationRepository,
    executor ActionExecutor,
) *AutomationEngine {
    return &AutomationEngine{repo: repo, actionExec: executor}
}

// OnTelemetry – เรียกเมื่อ telemetry เข้ามา
func (e *AutomationEngine) OnTelemetry(ctx context.Context, t *entity.Telemetry) error {
    automations, err := e.repo.ListActiveByTrigger(ctx, t.TenantID, valueobject.AutomationTriggerTelemetry)
    if err != nil {
        return err
    }
    if len(automations) == 0 {
        return nil
    }

    facts := telemetryFacts(t)
    for _, a := range automations {
        if !matchTrigger(a.Trigger, t) {
            continue
        }
        if !a.CanRun(facts) {
            continue
        }
        if err := e.execute(ctx, a, t); err != nil {
            log.Printf("[automation] %s failed: %v", a.ID, err)
            a.MarkRun(t.Timestamp, err)
        } else {
            a.MarkRun(t.Timestamp, nil)
        }
        _ = e.repo.Save(ctx, a)
    }
    return nil
}

func (e *AutomationEngine) execute(ctx context.Context, a *entity.Automation, t *entity.Telemetry) error {
    for _, act := range a.Actions {
        if err := e.actionExec.Execute(ctx, t.TenantID.String(), act); err != nil {
            return err
        }
    }
    return nil
}

func matchTrigger(trigger valueobject.AutomationTrigger, t *entity.Telemetry) bool {
    if trigger.DeviceID != nil && *trigger.DeviceID != t.DeviceID.String() {
        return false
    }
    if trigger.Metric != nil {
        found := false
        for _, mv := range t.Metrics {
            if mv.Metric == *trigger.Metric {
                found = true
                break
            }
        }
        if !found {
            return false
        }
    }
    return true
}

func telemetryFacts(t *entity.Telemetry) map[string]any {
    facts := map[string]any{
        "device_id": t.DeviceID.String(),
        "site_id":   t.SiteID.String(),
        "timestamp": t.Timestamp.Unix(),
    }
    for _, mv := range t.Metrics {
        facts[string(mv.Metric)] = mv.Value
    }
    return facts
}
```

### `domain/service/command_dispatcher.go`
```go
package service

import (
    "context"
    "fmt"

    "icmongolang/internal/modules/device/domain/entity"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
    "icmongolang/internal/modules/device/domain/service/port"
)

// CommandDispatcher – ส่ง command ผ่าน MQTT
type CommandDispatcher struct {
    mqtt port.MQTTPublisherPort
}

func NewCommandDispatcher(mqtt port.MQTTPublisherPort) *CommandDispatcher {
    return &CommandDispatcher{mqtt: mqtt}
}

// Dispatch – publish command ไปยัง topic
func (d *CommandDispatcher) Dispatch(ctx context.Context, dev *entity.Device, cmd *entity.Command) error {
    if !dev.CanReceiveCommand() {
        return domainerrors.ErrDeviceCannotReceiveCommand
    }
    topic := fmt.Sprintf("iot/%s/%s/cmd", dev.TenantID.String(), dev.ID.String())
    payload := map[string]any{
        "command_id": cmd.ID.String(),
        "command":    cmd.Command,
        "payload":    cmd.Payload,
        "issued_at":  cmd.IssuedAt.Unix(),
        "expires_at": cmd.ExpiresAt.Unix(),
    }
    return d.mqtt.Publish(ctx, topic, cmd.ID.String(), payload)
}
```

### `domain/service/port/mqtt_publisher_port.go`
```go
package port

import "context"

// MQTTPublisherPort – outbound port สำหรับ MQTT publish
type MQTTPublisherPort interface {
    Publish(ctx context.Context, topic, key string, payload any) error
}
```

### `domain/service/port/ai_inference_port.go`
```go
package port

import "context"

// AIInferencePort – outbound port สำหรับ AI
type AIInferencePort interface {
    AnomalyDetect(ctx context.Context, model string, values map[string]float64) (*AnomalyResult, error)
    Forecast(ctx context.Context, model string, metric string, horizonHours int, history []float64) (*ForecastResult, error)
    Classify(ctx context.Context, model, text string) (string, float64, error)
}

type AnomalyResult struct {
    IsAnomaly  bool
    Score      float64
    Reason     string
}

type ForecastResult struct {
    Values     []float64
    Confidence float64
}
```

### `domain/service/port/notifier_port.go`
```go
package port

import "context"

type NotifierPort interface {
    SendAlert(ctx context.Context, tenantID, deviceID, severity, message string, recipients []string) error
    SendEmail(ctx context.Context, to, subject, body string) error
}
```

### `domain/service/port/code_generator_port.go`
```go
package port

import (
    "context"

    "github.com/google/uuid"
)

type CodeGeneratorPort interface {
    NextDeviceTag(ctx context.Context, tenantID uuid.UUID) (string, error)
}
```

---

## A.6 Domain Events

### `domain/event/device_events.go`
```go
package event

import (
    "time"

    "github.com/google/uuid"
)

const (
    TopicDeviceRegistered     = "device.registered"
    TopicDeviceProvisioned    = "device.provisioned"
    TopicDeviceOnline         = "device.online"
    TopicDeviceOffline        = "device.offline"
    TopicDeviceFault          = "device.fault"
    TopicDeviceDecommissioned = "device.decommissioned"
    TopicDeviceFirmwareUpdated = "device.firmware.updated"
)

type DeviceRegistered struct {
    EventID    uuid.UUID `json:"event_id"`
    DeviceID   uuid.UUID `json:"device_id"`
    TenantID   uuid.UUID `json:"tenant_id"`
    CustomerID uuid.UUID `json:"customer_id"`
    SiteID     uuid.UUID `json:"site_id"`
    SerialNo   string    `json:"serial_no"`
    Protocol   string    `json:"protocol"`
    OccurredAt time.Time `json:"occurred_at"`
}

type DeviceProvisioned struct {
    EventID      uuid.UUID `json:"event_id"`
    DeviceID     uuid.UUID `json:"device_id"`
    TenantID     uuid.UUID `json:"tenant_id"`
    MQTTClientID string    `json:"mqtt_client_id"`
    OccurredAt   time.Time `json:"occurred_at"`
}

type DeviceOnline struct {
    EventID    uuid.UUID `json:"event_id"`
    DeviceID   uuid.UUID `json:"device_id"`
    TenantID   uuid.UUID `json:"tenant_id"`
    OccurredAt time.Time `json:"occurred_at"`
}

type DeviceOffline struct {
    EventID    uuid.UUID `json:"event_id"`
    DeviceID   uuid.UUID `json:"device_id"`
    TenantID   uuid.UUID `json:"tenant_id"`
    Reason     string    `json:"reason"`
    OccurredAt time.Time `json:"occurred_at"`
}

type DeviceFault struct {
    EventID    uuid.UUID `json:"event_id"`
    DeviceID   uuid.UUID `json:"device_id"`
    TenantID   uuid.UUID `json:"tenant_id"`
    Code       string    `json:"code"`
    Message    string    `json:"message"`
    OccurredAt time.Time `json:"occurred_at"`
}
```

### `domain/event/telemetry_events.go`
```go
package event

import (
    "time"

    "github.com/google/uuid"
)

const (
    TopicTelemetryRaw        = "iot.telemetry.raw"
    TopicTelemetryAggregated = "iot.telemetry.aggregated"
    TopicTelemetryBatch      = "iot.telemetry.batch"
)

type TelemetryRaw struct {
    EventID    uuid.UUID           `json:"event_id"`
    TenantID   uuid.UUID           `json:"tenant_id"`
    DeviceID   uuid.UUID           `json:"device_id"`
    SiteID     uuid.UUID           `json:"site_id"`
    Metrics    []TelemetryMetric   `json:"metrics"`
    Timestamp  time.Time           `json:"timestamp"`
    Source     string              `json:"source"`
    OccurredAt time.Time           `json:"occurred_at"`
}

type TelemetryMetric struct {
    Metric  string  `json:"metric"`
    Value   float64 `json:"value"`
    Unit    string  `json:"unit,omitempty"`
    Quality string  `json:"quality,omitempty"`
}

type TelemetryAggregated struct {
    EventID    uuid.UUID `json:"event_id"`
    TenantID   uuid.UUID `json:"tenant_id"`
    DeviceID   uuid.UUID `json:"device_id"`
    CustomerID uuid.UUID `json:"customer_id"`
    Metric     string    `json:"metric"`
    Count      int       `json:"count"`
    Min, Max, Avg float64 `json:"min_max_avg"`
    Window     time.Time `json:"window"`
    OccurredAt time.Time `json:"occurred_at"`
}
```

### `domain/event/command_events.go`
```go
package event

import (
    "time"

    "github.com/google/uuid"
)

const (
    TopicCommandIssued = "device.command.issued"
    TopicCommandSent   = "device.command.sent"
    TopicCommandAcked  = "device.command.acked"
    TopicCommandFailed = "device.command.failed"
    TopicCommandTimeout = "device.command.timeout"
)

type CommandIssued struct {
    EventID   uuid.UUID `json:"event_id"`
    CommandID uuid.UUID `json:"command_id"`
    DeviceID  uuid.UUID `json:"device_id"`
    TenantID  uuid.UUID `json:"tenant_id"`
    Command   string    `json:"command"`
    IssuedBy  uuid.UUID `json:"issued_by"`
    OccurredAt time.Time `json:"occurred_at"`
}

type CommandAcked struct {
    EventID    uuid.UUID      `json:"event_id"`
    CommandID  uuid.UUID      `json:"command_id"`
    DeviceID   uuid.UUID      `json:"device_id"`
    TenantID   uuid.UUID      `json:"tenant_id"`
    Result     map[string]any `json:"result,omitempty"`
    OccurredAt time.Time      `json:"occurred_at"`
}

type CommandFailed struct {
    EventID    uuid.UUID `json:"event_id"`
    CommandID  uuid.UUID `json:"command_id"`
    DeviceID   uuid.UUID `json:"device_id"`
    TenantID   uuid.UUID `json:"tenant_id"`
    Reason     string    `json:"reason"`
    OccurredAt time.Time `json:"occurred_at"`
}
```

### `domain/event/alert_events.go`
```go
package event

import (
    "time"

    "github.com/google/uuid"
)

const (
    TopicAlertTriggered    = "device.alert.triggered"
    TopicAlertAcknowledged = "device.alert.acknowledged"
    TopicAlertResolved     = "device.alert.resolved"
)

type AlertTriggered struct {
    EventID    uuid.UUID `json:"event_id"`
    AlertID    int64     `json:"alert_id"`
    RuleID     uuid.UUID `json:"rule_id"`
    DeviceID   uuid.UUID `json:"device_id"`
    CustomerID uuid.UUID `json:"customer_id"`
    TenantID   uuid.UUID `json:"tenant_id"`
    Metric     string    `json:"metric"`
    Value      float64   `json:"value"`
    Threshold  float64   `json:"threshold"`
    Severity   string    `json:"severity"`
    Message    string    `json:"message"`
    OccurredAt time.Time `json:"occurred_at"`
}
```

### `domain/event/automation_events.go`
```go
package event

import (
    "time"

    "github.com/google/uuid"
)

const (
    TopicAutomationTriggered = "device.automation.triggered"
    TopicAutomationFailed    = "device.automation.failed"
)

type AutomationTriggered struct {
    EventID      uuid.UUID `json:"event_id"`
    AutomationID uuid.UUID `json:"automation_id"`
    TenantID     uuid.UUID `json:"tenant_id"`
    DeviceID     uuid.UUID `json:"device_id"`
    ActionsCount int       `json:"actions_count"`
    OccurredAt   time.Time `json:"occurred_at"`
}
```

---

## A.7 Domain Errors

### `domain/errors/errors.go`
```go
package domainerrors

import "errors"

var (
    // Device core
    ErrDeviceNotFound              = errors.New("device not found")
    ErrDeviceAlreadyProvisioned    = errors.New("device already provisioned")
    ErrDeviceNotProvisioned        = errors.New("device not provisioned")
    ErrDeviceDecommissioned        = errors.New("device is decommissioned")
    ErrDeviceAlreadyDecommissioned = errors.New("device already decommissioned")
    ErrDeviceCannotReceiveCommand  = errors.New("device cannot receive command")
    ErrDeviceNotInFault            = errors.New("device not in fault state")
    ErrDeviceNotInMaintenance      = errors.New("device not in maintenance state")
    ErrInvalidStatusTransition     = errors.New("invalid status transition")

    // Validation
    ErrInvalidTenantID      = errors.New("invalid tenant id")
    ErrInvalidCustomerID    = errors.New("invalid customer id")
    ErrInvalidSiteID        = errors.New("invalid site id")
    ErrInvalidDeviceName    = errors.New("invalid device name")
    ErrInvalidDeviceType    = errors.New("invalid device type")
    ErrInvalidSerial        = errors.New("invalid device serial")
    ErrInvalidProtocol      = errors.New("invalid protocol")
    ErrInvalidMetric        = errors.New("invalid metric")
    ErrInvalidMetricValue   = errors.New("invalid metric value")
    ErrInvalidCapability    = errors.New("invalid capability")
    ErrEmptyCapabilities    = errors.New("empty capabilities")
    ErrInvalidFirmwareVersion = errors.New("invalid firmware version")

    // Provisioning
    ErrEmptyDeviceToken    = errors.New("empty device token")
    ErrEmptyMQTTClientID   = errors.New("empty mqtt client id")
    ErrWeakDeviceToken     = errors.New("device token too weak")

    // Command
    ErrInvalidCommand          = errors.New("invalid command")
    ErrCommandTooLong          = errors.New("command too long")
    ErrCommandAlreadyTerminal  = errors.New("command already terminal")
    ErrInvalidCommandState     = errors.New("invalid command state")
    ErrCommandTimeout          = errors.New("command timeout")

    // Shadow
    ErrShadowNotFound         = errors.New("shadow not found")
    ErrShadowVersionMismatch  = errors.New("shadow version mismatch")
    ErrEmptyPatch             = errors.New("empty patch")

    // Alert
    ErrInvalidAlertRuleName = errors.New("invalid alert rule name")
    ErrInvalidOperator      = errors.New("invalid operator")
    ErrInvalidSeverity      = errors.New("invalid severity")
    ErrAlertRuleNotFound    = errors.New("alert rule not found")
    ErrAlertEventNotFound   = errors.New("alert event not found")

    // Automation
    ErrInvalidAutomationName = errors.New("invalid automation name")
    ErrInvalidTrigger        = errors.New("invalid trigger")
    ErrTriggerMetricRequired = errors.New("trigger metric required")
    ErrTriggerCronRequired   = errors.New("trigger cron required")
    ErrInvalidCondition      = errors.New("invalid condition")
    ErrInvalidAction         = errors.New("invalid action")
    ErrTooManyActions        = errors.New("too many actions")
    ErrTooManyConditions     = errors.New("too many conditions")
    ErrAutomationNotFound    = errors.New("automation not found")

    // Group
    ErrInvalidGroupName = errors.New("invalid group name")
    ErrSelfReference    = errors.New("cannot reference self")

    // AI
    ErrInvalidAIModelName = errors.New("invalid ai model name")
    ErrAIModelNotFound    = errors.New("ai model not found")
    ErrInferenceFailed    = errors.New("ai inference failed")

    // Infrastructure
    ErrPersistenceFailure   = errors.New("persistence failure")
    ErrMQTTPublishFailure   = errors.New("mqtt publish failure")
    ErrInfluxWriteFailure   = errors.New("influx write failure")
    ErrConcurrentModification = errors.New("concurrent modification")
)
```

---

## A.8 Unit Tests (Domain)

### `domain/entity/device_test.go`
```go
package entity_test

import (
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "icmongolang/internal/modules/device/domain/entity"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

func newTestDevice(t *testing.T) *entity.Device {
    serial, _ := valueobject.NewDeviceSerial("SN-TEST-0001")
    d, err := entity.NewDevice(
        uuid.New(), uuid.New(), uuid.New(), uuid.New(),
        serial,
        valueobject.DeviceTypeSensor,
        valueobject.ProtocolMQTT,
        "Test Sensor",
    )
    require.NoError(t, err)
    return d
}

func TestNewDevice_Success(t *testing.T) {
    d := newTestDevice(t)
    assert.Equal(t, valueobject.DeviceStatusRegistered, d.Status)
    assert.False(t, d.IsProvisioned())
}

func TestDevice_Lifecycle_RegisterToDecommission(t *testing.T) {
    d := newTestDevice(t)

    // Provision
    token, _ := valueobject.NewDeviceToken()
    require.NoError(t, d.Provision(token, "dev-test"))
    assert.Equal(t, valueobject.DeviceStatusProvisioned, d.Status)

    // Online
    require.NoError(t, d.MarkOnline(time.Now()))
    assert.True(t, d.IsOnline())

    // Fault
    require.NoError(t, d.ReportFault("E001", "overheat"))
    assert.True(t, d.IsFaulty())

    // Clear fault → offline
    require.NoError(t, d.ClearFault())
    assert.Equal(t, valueobject.DeviceStatusOffline, d.Status)

    // Decommission
    require.NoError(t, d.Decommission("end of life", uuid.New()))
    assert.True(t, d.IsDecommissioned())

    // Cannot mark online again
    err := d.MarkOnline(time.Now())
    assert.ErrorIs(t, err, domainerrors.ErrDeviceDecommissioned)
}

func TestDevice_Provision_Twice_ShouldFail(t *testing.T) {
    d := newTestDevice(t)
    token, _ := valueobject.NewDeviceToken()
    require.NoError(t, d.Provision(token, "dev-1"))
    err := d.Provision(token, "dev-2")
    assert.ErrorIs(t, err, domainerrors.ErrDeviceAlreadyProvisioned)
}

func TestDevice_StaleCheck(t *testing.T) {
    d := newTestDevice(t)
    token, _ := valueobject.NewDeviceToken()
    _ = d.Provision(token, "dev-1")
    ts := time.Now().Add(-10 * time.Minute)
    _ = d.MarkOnline(ts)
    assert.True(t, d.IsStale(5*time.Minute))
    assert.False(t, d.IsStale(15*time.Minute))
}
```

### `domain/entity/device_shadow_test.go`
```go
package entity_test

import (
    "testing"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "icmongolang/internal/modules/device/domain/entity"
)

func TestShadow_DesiredReportedDelta(t *testing.T) {
    s := entity.NewDeviceShadow(uuid.New(), uuid.New())

    // Set desired: fan=on, threshold=25
    require.NoError(t, s.SetDesired(map[string]any{
        "fan": "on", "threshold": 25.0,
    }))
    assert.Len(t, s.GetDelta(), 2)

    // Device reports fan=on (match), threshold=20 (mismatch)
    require.NoError(t, s.UpdateReported(map[string]any{
        "fan": "on", "threshold": 20.0,
    }))

    delta := s.GetDelta()
    assert.Len(t, delta, 1)
    assert.Equal(t, 25.0, delta["threshold"])
    assert.False(t, s.IsInSync())
}

func TestShadow_VersionIncrement(t *testing.T) {
    s := entity.NewDeviceShadow(uuid.New(), uuid.New())
    v0 := s.Version()
    _ = s.SetDesired(map[string]any{"x": 1})
    assert.Greater(t, s.Version(), v0)
}

func TestShadow_OptimisticLock(t *testing.T) {
    s := entity.NewDeviceShadow(uuid.New(), uuid.New())
    _ = s.SetDesired(map[string]any{"x": 1})
    v := s.Version()
    assert.NoError(t, s.CheckVersion(v))
    assert.Error(t, s.CheckVersion(v-1))
}
```

### `domain/value_object/device_token_test.go`
```go
package valueobject_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

func TestDeviceToken_GenerateAndVerify(t *testing.T) {
    token, err := valueobject.NewDeviceToken()
    require.NoError(t, err)
    assert.NotEmpty(t, token.Plaintext)
    assert.NotEmpty(t, token.Hash)

    assert.True(t, token.Verify(token.Plaintext))
    assert.False(t, token.Verify("wrong"))
}
```

---

# 🅱️ PART 3B — APPLICATION LAYER

## B.1 Use Cases

### `application/dto.go`
```go
package application

import "time"

// ============================================================
// DEVICE
// ============================================================

type RegisterDeviceInput struct {
    TenantID   string         `json:"-"`
    CustomerID string         `json:"customer_id" binding:"required"`
    SiteID     string         `json:"site_id" binding:"required"`
    GroupID    string         `json:"group_id,omitempty"`
    ModelID    string         `json:"model_id,omitempty"`
    SerialNo   string         `json:"serial_no" binding:"required"`
    Name       string         `json:"name" binding:"required"`
    Type       string         `json:"type" binding:"required"`
    Protocol   string         `json:"protocol" binding:"required"`
    Tags       []string       `json:"tags,omitempty"`
    Metadata   map[string]any `json:"metadata,omitempty"`
    ActorID    string         `json:"-"`
}

type ProvisionDeviceInput struct {
    TenantID    string `json:"-"`
    DeviceID    string `json:"-"`
    MQTTBroker  string `json:"mqtt_broker,omitempty"`
    ActorID     string `json:"-"`
}

type UpdateDeviceInput struct {
    TenantID   string          `json:"-"`
    DeviceID   string          `json:"-"`
    Name       *string         `json:"name,omitempty"`
    GroupID    *string         `json:"group_id,omitempty"`
    Tags       *[]string       `json:"tags,omitempty"`
    Metadata   *map[string]any `json:"metadata,omitempty"`
    ActorID    string          `json:"-"`
}

type DecommissionDeviceInput struct {
    TenantID string `json:"-"`
    DeviceID string `json:"-"`
    Reason   string `json:"reason" binding:"required"`
    ActorID  string `json:"-"`
}

type DeviceResponse struct {
    ID              string         `json:"id"`
    TenantID        string         `json:"tenant_id"`
    CustomerID      string         `json:"customer_id"`
    SiteID          string         `json:"site_id"`
    GroupID         string         `json:"group_id,omitempty"`
    ModelID         string         `json:"model_id,omitempty"`
    SerialNo        string         `json:"serial_no"`
    Name            string         `json:"name"`
    Type            string         `json:"type"`
    Protocol        string         `json:"protocol"`
    Status          string         `json:"status"`
    MQTTClientID    string         `json:"mqtt_client_id,omitempty"`
    FirmwareVersion string         `json:"firmware_version,omitempty"`
    Capabilities    []CapabilityDTO `json:"capabilities,omitempty"`
    Tags            []string       `json:"tags,omitempty"`
    Metadata        map[string]any `json:"metadata,omitempty"`
    LastSeenAt      *time.Time     `json:"last_seen_at,omitempty"`
    ProvisionedAt   *time.Time     `json:"provisioned_at,omitempty"`
    CreatedAt       time.Time      `json:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at"`
}

type ProvisionResponse struct {
    DeviceID     string `json:"device_id"`
    MQTTClientID string `json:"mqtt_client_id"`
    DeviceToken  string `json:"device_token"`  // return ครั้งเดียว
    MQTTBroker   string `json:"mqtt_broker"`
    TopicPrefix  string `json:"topic_prefix"`
}

type CapabilityDTO struct {
    Type     string   `json:"type"`
    Metric   string   `json:"metric,omitempty"`
    Command  string   `json:"command,omitempty"`
    MinValue *float64 `json:"min_value,omitempty"`
    MaxValue *float64 `json:"max_value,omitempty"`
    Unit     string   `json:"unit,omitempty"`
}

type ListDevicesInput struct {
    TenantID   string
    CustomerID string
    SiteID     string
    GroupID    string
    ModelID    string
    Statuses   []string
    Types      []string
    Protocol   string
    Search     string
    Tags       []string
    Page       int
    PageSize   int
    SortBy     string
    SortOrder  string
}

type ListDevicesResponse struct {
    Items    []DeviceResponse `json:"items"`
    Total    int64            `json:"total"`
    Page     int              `json:"page"`
    PageSize int              `json:"page_size"`
    Pages    int              `json:"pages"`
}

// ============================================================
// TELEMETRY
// ============================================================

type IngestTelemetryInput struct {
    TenantID   string         `json:"tenant_id" binding:"required"`
    DeviceID   string         `json:"device_id" binding:"required"`
    SiteID     string         `json:"site_id"`
    Metrics    []MetricDTO    `json:"metrics" binding:"required,min=1"`
    Timestamp  time.Time      `json:"timestamp"`
    Source     string         `json:"source,omitempty"`
    AuthToken  string         `json:"-"` // device token (header)
}

type MetricDTO struct {
    Metric  string  `json:"metric" binding:"required"`
    Value   float64 `json:"value"`
    Quality string  `json:"quality,omitempty"`
}

type QueryTelemetryInput struct {
    TenantID  string
    DeviceID  string
    Metrics   []string
    From      time.Time
    To        time.Time
    Aggregate string // raw | 1m | 5m | 1h | 1d
    Limit     int
}

type TelemetrySeriesResponse struct {
    DeviceID   string              `json:"device_id"`
    Series     map[string][]DataPoint `json:"series"`
    From       time.Time           `json:"from"`
    To         time.Time           `json:"to"`
    Aggregated bool                `json:"aggregated"`
}

type DataPoint struct {
    Timestamp time.Time `json:"t"`
    Value     float64   `json:"v"`
}

// ============================================================
// COMMAND
// ============================================================

type SendCommandInput struct {
    TenantID  string         `json:"-"`
    DeviceID  string         `json:"-"`
    Command   string         `json:"command" binding:"required"`
    Payload   map[string]any `json:"payload,omitempty"`
    Priority  int            `json:"priority,omitempty"`
    ActorID   string         `json:"-"`
}

type AckCommandInput struct {
    TenantID  string         `json:"-"`
    CommandID string         `json:"-"`
    Success   bool           `json:"success"`
    Result    map[string]any `json:"result,omitempty"`
    ErrorMsg  string         `json:"error_msg,omitempty"`
}

type CommandResponse struct {
    ID         string         `json:"id"`
    DeviceID   string         `json:"device_id"`
    Command    string         `json:"command"`
    Payload    map[string]any `json:"payload,omitempty"`
    Status     string         `json:"status"`
    IssuedBy   string         `json:"issued_by"`
    IssuedAt   time.Time      `json:"issued_at"`
    SentAt     *time.Time     `json:"sent_at,omitempty"`
    AckedAt    *time.Time     `json:"acked_at,omitempty"`
    ExpiresAt  time.Time      `json:"expires_at"`
    ErrorMsg   string         `json:"error_msg,omitempty"`
    Result     map[string]any `json:"result,omitempty"`
    RetryCount int            `json:"retry_count"`
}

// ============================================================
// SHADOW
// ============================================================

type UpdateShadowDesiredInput struct {
    TenantID string         `json:"-"`
    DeviceID string         `json:"-"`
    Desired  map[string]any `json:"desired" binding:"required"`
    Merge    bool           `json:"merge,omitempty"`
    ActorID  string         `json:"-"`
}

type UpdateShadowReportedInput struct {
    TenantID string         `json:"-"`
    DeviceID string         `json:"-"`
    Reported map[string]any `json:"reported" binding:"required"`
}

type ShadowResponse struct {
    DeviceID  string         `json:"device_id"`
    Desired   map[string]any `json:"desired"`
    Reported  map[string]any `json:"reported"`
    Delta     map[string]any `json:"delta"`
    InSync    bool           `json:"in_sync"`
    Version   int64          `json:"version"`
    UpdatedAt time.Time      `json:"updated_at"`
}

// ============================================================
// ALERT RULE
// ============================================================

type CreateAlertRuleInput struct {
    TenantID    string         `json:"-"`
    Name        string         `json:"name" binding:"required"`
    DeviceID    string         `json:"device_id,omitempty"`
    GroupID     string         `json:"group_id,omitempty"`
    SiteID      string         `json:"site_id,omitempty"`
    Metric      string         `json:"metric" binding:"required"`
    Operator    string         `json:"operator" binding:"required"`
    Threshold   float64        `json:"threshold"`
    DurationSec int            `json:"duration_sec,omitempty"`
    Severity    string         `json:"severity" binding:"required"`
    Actions     []AlertActionDTO `json:"actions,omitempty"`
    ActorID     string         `json:"-"`
}

type AlertActionDTO struct {
    Type    string         `json:"type"`
    Channel string         `json:"channel,omitempty"`
    Target  string         `json:"target,omitempty"`
    Payload map[string]any `json:"payload,omitempty"`
}

type AlertRuleResponse struct {
    ID          string           `json:"id"`
    Name        string           `json:"name"`
    DeviceID    string           `json:"device_id,omitempty"`
    Metric      string           `json:"metric"`
    Operator    string           `json:"operator"`
    Threshold   float64          `json:"threshold"`
    Severity    string           `json:"severity"`
    Actions     []AlertActionDTO `json:"actions"`
    IsActive    bool             `json:"is_active"`
    LastTriggeredAt *time.Time   `json:"last_triggered_at,omitempty"`
    CreatedAt   time.Time        `json:"created_at"`
}

// ============================================================
// AUTOMATION
// ============================================================

type CreateAutomationInput struct {
    TenantID    string         `json:"-"`
    Name        string         `json:"name" binding:"required"`
    Description string         `json:"description,omitempty"`
    Trigger     TriggerDTO     `json:"trigger" binding:"required"`
    Conditions  []ConditionDTO `json:"conditions,omitempty"`
    Actions     []ActionDTO    `json:"actions" binding:"required,min=1"`
    Priority    int            `json:"priority,omitempty"`
    ActorID     string         `json:"-"`
}

type TriggerDTO struct {
    Type     string `json:"type" binding:"required"`
    DeviceID string `json:"device_id,omitempty"`
    GroupID  string `json:"group_id,omitempty"`
    Metric   string `json:"metric,omitempty"`
    CronExpr string `json:"cron_expr,omitempty"`
}

type ConditionDTO struct {
    Field    string `json:"field" binding:"required"`
    Operator string `json:"operator" binding:"required"`
    Value    any    `json:"value"`
}

type ActionDTO struct {
    Type     string         `json:"type" binding:"required"`
    DeviceID string         `json:"device_id,omitempty"`
    Command  string         `json:"command,omitempty"`
    Payload  map[string]any `json:"payload,omitempty"`
    Channel  string         `json:"channel,omitempty"`
    Target   string         `json:"target,omitempty"`
    DelayMs  int            `json:"delay_ms,omitempty"`
}

type AutomationResponse struct {
    ID          string         `json:"id"`
    Name        string         `json:"name"`
    Description string         `json:"description,omitempty"`
    Trigger     TriggerDTO     `json:"trigger"`
    Conditions  []ConditionDTO `json:"conditions"`
    Actions     []ActionDTO    `json:"actions"`
    IsActive    bool           `json:"is_active"`
    Priority    int            `json:"priority"`
    RunCount    int64          `json:"run_count"`
    LastRunAt   *time.Time     `json:"last_run_at,omitempty"`
    LastError   string         `json:"last_error,omitempty"`
    CreatedAt   time.Time      `json:"created_at"`
}
```

### `application/mappers.go`
```go
package application

import (
    "icmongolang/internal/modules/device/domain/entity"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

func toDeviceResponse(d *entity.Device) *DeviceResponse {
    r := &DeviceResponse{
        ID:              d.ID.String(),
        TenantID:        d.TenantID.String(),
        CustomerID:      d.CustomerID.String(),
        SiteID:          d.SiteID.String(),
        SerialNo:        d.SerialNo.String(),
        Name:            d.Name,
        Type:            string(d.Type),
        Protocol:        string(d.Protocol),
        Status:          string(d.Status),
        MQTTClientID:    d.MQTTClientID,
        FirmwareVersion: d.FirmwareVersion,
        Tags:            d.Tags,
        Metadata:        d.Metadata,
        LastSeenAt:      d.LastSeenAt,
        ProvisionedAt:   d.ProvisionedAt,
        CreatedAt:       d.CreatedAt,
        UpdatedAt:       d.UpdatedAt,
    }
    if d.GroupID != nil { r.GroupID = d.GroupID.String() }
    if d.ModelID != nil { r.ModelID = d.ModelID.String() }
    if len(d.Capabilities) > 0 {
        r.Capabilities = make([]CapabilityDTO, 0, len(d.Capabilities))
        for _, c := range d.Capabilities {
            r.Capabilities = append(r.Capabilities, CapabilityDTO{
                Type: c.Type, Metric: string(c.Metric), Command: c.Command,
                MinValue: c.MinValue, MaxValue: c.MaxValue, Unit: c.Unit,
            })
        }
    }
    return r
}

func toCommandResponse(c *entity.Command) *CommandResponse {
    return &CommandResponse{
        ID:         c.ID.String(),
        DeviceID:   c.DeviceID.String(),
        Command:    c.Command,
        Payload:    c.Payload,
        Status:     string(c.Status),
        IssuedBy:   c.IssuedBy.String(),
        IssuedAt:   c.IssuedAt,
        SentAt:     c.SentAt,
        AckedAt:    c.AckedAt,
        ExpiresAt:  c.ExpiresAt,
        ErrorMsg:   c.ErrorMsg,
        Result:     c.ResultPayload,
        RetryCount: c.RetryCount,
    }
}

func toShadowResponse(s *entity.DeviceShadow) *ShadowResponse {
    return &ShadowResponse{
        DeviceID:  s.DeviceID.String(),
        Desired:   s.State.Desired,
        Reported:  s.State.Reported,
        Delta:     s.State.Delta,
        InSync:    s.IsInSync(),
        Version:   s.Version(),
        UpdatedAt: s.UpdatedAt,
    }
}

func toAlertRuleResponse(r *entity.AlertRule) *AlertRuleResponse {
    res := &AlertRuleResponse{
        ID:              r.ID.String(),
        Name:            r.Name,
        Metric:          string(r.Metric),
        Operator:        r.Operator,
        Threshold:       r.Threshold,
        Severity:        string(r.Severity),
        IsActive:        r.IsActive,
        LastTriggeredAt: r.LastTriggeredAt,
        CreatedAt:       r.CreatedAt,
    }
    if r.DeviceID != nil { res.DeviceID = r.DeviceID.String() }
    for _, a := range r.Actions {
        res.Actions = append(res.Actions, AlertActionDTO{
            Type: a.Type, Channel: a.Channel, Target: a.Target, Payload: a.Payload,
        })
    }
    return res
}

func toAutomationResponse(a *entity.Automation) *AutomationResponse {
    res := &AutomationResponse{
        ID:          a.ID.String(),
        Name:        a.Name,
        Description: a.Description,
        IsActive:    a.IsActive,
        Priority:    a.Priority,
        RunCount:    a.RunCount,
        LastRunAt:   a.LastRunAt,
        LastError:   a.LastError,
        CreatedAt:   a.CreatedAt,
    }
    res.Trigger = TriggerDTO{
        Type: string(a.Trigger.Type),
    }
    if a.Trigger.DeviceID != nil { res.Trigger.DeviceID = *a.Trigger.DeviceID }
    if a.Trigger.GroupID != nil  { res.Trigger.GroupID = *a.Trigger.GroupID }
    if a.Trigger.Metric != nil   { res.Trigger.Metric = string(*a.Trigger.Metric) }
    res.Trigger.CronExpr = a.Trigger.CronExpr

    for _, c := range a.Conditions {
        res.Conditions = append(res.Conditions, ConditionDTO{
            Field: c.Field, Operator: c.Operator, Value: c.Value,
        })
    }
    for _, act := range a.Actions {
        dto := ActionDTO{
            Type: act.Type.String(), Command: act.Command, Payload: act.Payload,
            Channel: act.Channel, Target: act.Target, DelayMs: act.DelayMs,
        }
        if act.DeviceID != nil { dto.DeviceID = *act.DeviceID }
        res.Actions = append(res.Actions, dto)
    }
    return res
}

var _ = valueobject.MetricTemperature
```

### `application/register_device.go`
```go
package application

import (
    "context"
    "time"

    "github.com/google/uuid"

    "icmongolang/internal/modules/device/domain/entity"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
    "icmongolang/internal/modules/device/domain/event"
    "icmongolang/internal/modules/device/domain/repository"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
)

type RegisterDeviceUseCase struct {
    deviceRepo     repository.DeviceRepository
    modelRepo      repository.DeviceModelRepository
    entitlementChk EntitlementChecker // outbound port
    producer       kafka.Producer
    log            logger.Logger
}

// EntitlementChecker – port ตรวจ quota ก่อน create
type EntitlementChecker interface {
    CheckDeviceQuota(ctx context.Context, tenantID, customerID uuid.UUID) error
}

func NewRegisterDeviceUseCase(
    deviceRepo repository.DeviceRepository,
    modelRepo repository.DeviceModelRepository,
    entitlementChk EntitlementChecker,
    producer kafka.Producer,
    log logger.Logger,
) *RegisterDeviceUseCase {
    return &RegisterDeviceUseCase{
        deviceRepo: deviceRepo, modelRepo: modelRepo,
        entitlementChk: entitlementChk, producer: producer, log: log,
    }
}

func (uc *RegisterDeviceUseCase) Execute(ctx context.Context, in RegisterDeviceInput) (*DeviceResponse, error) {
    tenantID, err := uuid.Parse(in.TenantID)
    if err != nil { return nil, domainerrors.ErrInvalidTenantID }
    customerID, _ := uuid.Parse(in.CustomerID)
    siteID, _ := uuid.Parse(in.SiteID)
    actorID, _ := uuid.Parse(in.ActorID)

    // 1. Quota check
    if err := uc.entitlementChk.CheckDeviceQuota(ctx, tenantID, customerID); err != nil {
        return nil, err
    }

    // 2. Parse value objects
    serial, err := valueobject.NewDeviceSerial(in.SerialNo)
    if err != nil { return nil, err }

    typ := valueobject.DeviceType(in.Type)
    protocol := valueobject.Protocol(in.Protocol)

    // 3. Create aggregate
    d, err := entity.NewDevice(tenantID, customerID, siteID, actorID, serial, typ, protocol, in.Name)
    if err != nil { return nil, err }

    // 4. Optional: group, model
    if in.GroupID != "" {
        gid, _ := uuid.Parse(in.GroupID)
        d.GroupID = &gid
    }
    if in.ModelID != "" {
        mid, _ := uuid.Parse(in.ModelID)
        d.ModelID = &mid
        // copy capabilities
        if m, err := uc.modelRepo.FindByID(ctx, mid); err == nil && m != nil {
            d.Capabilities = m.Capabilities
        }
    }
    if len(in.Tags) > 0 { d.Tags = in.Tags }
    if in.Metadata != nil { d.Metadata = in.Metadata }

    // 5. Persist
    if err := uc.deviceRepo.Save(ctx, d); err != nil {
        uc.log.Error("save device failed", "err", err)
        return nil, domainerrors.ErrPersistenceFailure
    }

    // 6. Publish event
    _ = uc.producer.Publish(ctx, event.TopicDeviceRegistered, d.ID.String(), event.DeviceRegistered{
        EventID: uuid.New(), DeviceID: d.ID, TenantID: tenantID,
        CustomerID: customerID, SiteID: siteID,
        SerialNo: d.SerialNo.String(), Protocol: string(d.Protocol),
        OccurredAt: time.Now(),
    })
    return toDeviceResponse(d), nil
}
```

### `application/provision_device.go`
```go
package application

import (
    "context"
    "time"

    "github.com/google/uuid"

    domainerrors "icmongolang/internal/modules/device/domain/errors"
    "icmongolang/internal/modules/device/domain/event"
    "icmongolang/internal/modules/device/domain/repository"
    "icmongolang/internal/modules/device/domain/service"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
)

type ProvisionDeviceUseCase struct {
    deviceRepo repository.DeviceRepository
    provSvc    *service.ProvisioningService
    producer   kafka.Producer
    log        logger.Logger
}

func NewProvisionDeviceUseCase(
    deviceRepo repository.DeviceRepository,
    provSvc *service.ProvisioningService,
    producer kafka.Producer,
    log logger.Logger,
) *ProvisionDeviceUseCase {
    return &ProvisionDeviceUseCase{
        deviceRepo: deviceRepo, provSvc: provSvc,
        producer: producer, log: log,
    }
}

func (uc *ProvisionDeviceUseCase) Execute(ctx context.Context, in ProvisionDeviceInput) (*ProvisionResponse, error) {
    tenantID, _ := uuid.Parse(in.TenantID)
    deviceID, _ := uuid.Parse(in.DeviceID)

    d, err := uc.deviceRepo.FindByID(ctx, tenantID, deviceID)
    if err != nil { return nil, err }

    if d.IsProvisioned() {
        return nil, domainerrors.ErrDeviceAlreadyProvisioned
    }

    // Delegate to domain service
    result, err := uc.provSvc.Provision(ctx, d)
    if err != nil { return nil, err }

    if err := uc.deviceRepo.Save(ctx, d); err != nil {
        return nil, domainerrors.ErrPersistenceFailure
    }

    _ = uc.producer.Publish(ctx, event.TopicDeviceProvisioned, d.ID.String(), event.DeviceProvisioned{
        EventID: uuid.New(), DeviceID: d.ID, TenantID: tenantID,
        MQTTClientID: result.MQTTClientID, OccurredAt: time.Now(),
    })

    return &ProvisionResponse{
        DeviceID:     d.ID.String(),
        MQTTClientID: result.MQTTClientID,
        DeviceToken:  result.Token.Plaintext, // ⚠️ return ครั้งเดียว
        MQTTBroker:   in.MQTTBroker,
        TopicPrefix:  result.TopicPrefix,
    }, nil
}
```

### `application/send_command.go`
```go
package application

import (
    "context"
    "time"

    "github.com/google/uuid"

    "icmongolang/internal/modules/device/domain/entity"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
    "icmongolang/internal/modules/device/domain/event"
    "icmongolang/internal/modules/device/domain/repository"
    "icmongolang/internal/modules/device/domain/service"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
    "icmongolang/pkg/transaction"
)

type SendCommandUseCase struct {
    deviceRepo  repository.DeviceRepository
    commandRepo repository.CommandRepository
    dispatcher  *service.CommandDispatcher
    producer    kafka.Producer
    tx          *transaction.Manager
    log         logger.Logger
}

func NewSendCommandUseCase(
    deviceRepo repository.DeviceRepository,
    commandRepo repository.CommandRepository,
    dispatcher *service.CommandDispatcher,
    producer kafka.Producer,
    tx *transaction.Manager,
    log logger.Logger,
) *SendCommandUseCase {
    return &SendCommandUseCase{
        deviceRepo: deviceRepo, commandRepo: commandRepo,
        dispatcher: dispatcher, producer: producer, tx: tx, log: log,
    }
}

func (uc *SendCommandUseCase) Execute(ctx context.Context, in SendCommandInput) (*CommandResponse, error) {
    tenantID, _ := uuid.Parse(in.TenantID)
    deviceID, _ := uuid.Parse(in.DeviceID)
    actorID, _ := uuid.Parse(in.ActorID)

    d, err := uc.deviceRepo.FindByID(ctx, tenantID, deviceID)
    if err != nil { return nil, err }

    if d.IsDecommissioned() {
        return nil, domainerrors.ErrDeviceDecommissioned
    }
    if !d.HasCommandCapability(in.Command) && len(d.Capabilities) > 0 {
        return nil, domainerrors.ErrInvalidCommand
    }

    cmd, err := entity.NewCommand(tenantID, deviceID, actorID, in.Command, in.Payload)
    if err != nil { return nil, err }
    if in.Priority > 0 { cmd.Priority = in.Priority }

    // Decide: send now or queue
    offline := !d.CanReceiveCommand()
    if offline {
        _ = cmd.Queue("device_offline")
        if err := uc.commandRepo.Save(ctx, cmd); err != nil {
            return nil, domainerrors.ErrPersistenceFailure
        }
    } else {
        if err := uc.commandRepo.Save(ctx, cmd); err != nil {
            return nil, domainerrors.ErrPersistenceFailure
        }
        // Dispatch ผ่าน MQTT
        if err := uc.dispatcher.Dispatch(ctx, d, cmd); err != nil {
            _ = cmd.Nack("mqtt_publish_failed")
            _ = uc.commandRepo.Save(ctx, cmd)
            return nil, domainerrors.ErrMQTTPublishFailure
        }
        _ = cmd.MarkSent(time.Now())
        _ = uc.commandRepo.Save(ctx, cmd)
    }

    _ = uc.producer.Publish(ctx, event.TopicCommandIssued, cmd.ID.String(), event.CommandIssued{
        EventID: uuid.New(), CommandID: cmd.ID, DeviceID: d.ID,
        TenantID: tenantID, Command: cmd.Command, IssuedBy: actorID,
        OccurredAt: time.Now(),
    })
    return toCommandResponse(cmd), nil
}

var _ = valueobject.CommandStatusSent
```

### `application/ack_command.go`
```go
package application

import (
    "context"
    "time"

    "github.com/google/uuid"

    domainerrors "icmongolang/internal/modules/device/domain/errors"
    "icmongolang/internal/modules/device/domain/event"
    "icmongolang/internal/modules/device/domain/repository"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
)

type AckCommandUseCase struct {
    commandRepo repository.CommandRepository
    producer    kafka.Producer
    log         logger.Logger
}

func NewAckCommandUseCase(
    commandRepo repository.CommandRepository,
    producer kafka.Producer,
    log logger.Logger,
) *AckCommandUseCase {
    return &AckCommandUseCase{commandRepo: commandRepo, producer: producer, log: log}
}

// Execute – เรียกจาก Kafka consumer (device.command.ack)
func (uc *AckCommandUseCase) Execute(ctx context.Context, in AckCommandInput) error {
    tenantID, _ := uuid.Parse(in.TenantID)
    cmdID, _ := uuid.Parse(in.CommandID)

    cmd, err := uc.commandRepo.FindByID(ctx, tenantID, cmdID)
    if err != nil { return err }

    if cmd.IsTerminal() {
        return nil // idempotent
    }

    if in.Success {
        if err := cmd.Ack(in.Result); err != nil {
            return err
        }
        _ = uc.producer.Publish(ctx, event.TopicCommandAcked, cmd.ID.String(), event.CommandAcked{
            EventID: uuid.New(), CommandID: cmd.ID, DeviceID: cmd.DeviceID,
            TenantID: tenantID, Result: in.Result, OccurredAt: time.Now(),
        })
    } else {
        if err := cmd.Nack(in.ErrorMsg); err != nil {
            return err
        }
        _ = uc.producer.Publish(ctx, event.TopicCommandFailed, cmd.ID.String(), event.CommandFailed{
            EventID: uuid.New(), CommandID: cmd.ID, DeviceID: cmd.DeviceID,
            TenantID: tenantID, Reason: in.ErrorMsg, OccurredAt: time.Now(),
        })
    }

    return uc.commandRepo.Save(ctx, cmd)
}

var _ = domainerrors.ErrCommandAlreadyTerminal
```

### `application/ingest_telemetry.go` ⭐ (Hot path)
```go
package application

import (
    "context"
    "sync"
    "time"

    "github.com/google/uuid"

    "icmongolang/internal/modules/device/domain/entity"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
    "icmongolang/internal/modules/device/domain/event"
    "icmongolang/internal/modules/device/domain/repository"
    "icmongolang/internal/modules/device/domain/service"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
)

// DeviceStateCache – Redis cache สำหรับ device state (fast)
type DeviceStateCache interface {
    SetLastSeen(ctx context.Context, deviceID uuid.UUID, ts time.Time) error
    GetLastSeen(ctx context.Context, deviceID uuid.UUID) (time.Time, error)
    SetStatus(ctx context.Context, deviceID uuid.UUID, status string) error
}

// IngestTelemetryUseCase – HOT PATH
// ลำดับความสำคัญ:
//   1. Fast Redis write (LastSeen + online)
//   2. InfluxDB batch write (async)
//   3. Alert evaluation (async, ใน worker)
//   4. Automation engine (async)
//   5. Kafka aggregated event (async)
type IngestTelemetryUseCase struct {
    telemetryRepo repository.TelemetryRepository
    deviceRepo    repository.DeviceRepository
    stateCache    DeviceStateCache
    producer      kafka.Producer
    log           logger.Logger

    // batch buffer
    mu    sync.Mutex
    batch *entity.TelemetryBatch
    flushInterval time.Duration
}

func NewIngestTelemetryUseCase(
    telemetryRepo repository.TelemetryRepository,
    deviceRepo repository.DeviceRepository,
    stateCache DeviceStateCache,
    producer kafka.Producer,
    log logger.Logger,
) *IngestTelemetryUseCase {
    return &IngestTelemetryUseCase{
        telemetryRepo: telemetryRepo, deviceRepo: deviceRepo,
        stateCache: stateCache, producer: producer, log: log,
        batch:         &entity.TelemetryBatch{},
        flushInterval: 100 * time.Millisecond,
    }
}

func (uc *IngestTelemetryUseCase) Execute(ctx context.Context, in IngestTelemetryInput) error {
    tenantID, err := uuid.Parse(in.TenantID)
    if err != nil { return domainerrors.ErrInvalidTenantID }
    deviceID, err := uuid.Parse(in.DeviceID)
    if err != nil { return domainerrors.ErrDeviceNotFound }

    // 1. Build metrics
    ts := in.Timestamp
    if ts.IsZero() { ts = time.Now() }

    metrics := make([]valueobject.MetricValue, 0, len(in.Metrics))
    for _, m := range in.Metrics {
        mv, err := valueobject.NewMetricValue(valueobject.Metric(m.Metric), m.Value, ts)
        if err != nil {
            uc.log.Warn("invalid metric", "device_id", deviceID, "metric", m.Metric, "err", err)
            continue
        }
        metrics = append(metrics, mv)
    }
    if len(metrics) == 0 {
        return domainerrors.ErrInvalidMetric
    }

    // 2. Fast state update (Redis) — non-blocking
    _ = uc.stateCache.SetLastSeen(ctx, deviceID, ts)
    _ = uc.stateCache.SetStatus(ctx, deviceID, "ONLINE")

    // 3. Build telemetry event
    t := &entity.Telemetry{
        TenantID: tenantID, DeviceID: deviceID,
        Metrics: metrics, Timestamp: ts, Source: defaultSource(in.Source),
    }
    if in.SiteID != "" {
        sid, _ := uuid.Parse(in.SiteID)
        t.SiteID = sid
    }

    // 4. Write to InfluxDB (batch)
    if err := uc.telemetryRepo.Write(ctx, t); err != nil {
        uc.log.Error("influx write failed", "err", err)
        return domainerrors.ErrInfluxWriteFailure
    }

    // 5. Publish to Kafka (raw) — consumer จะทำ alert + automation + aggregate
    _ = uc.producer.Publish(ctx, event.TopicTelemetryRaw, deviceID.String(), event.TelemetryRaw{
        EventID: uuid.New(), TenantID: tenantID, DeviceID: deviceID,
        SiteID: t.SiteID, Metrics: toMetricsDTO(metrics),
        Timestamp: ts, Source: t.Source, OccurredAt: time.Now(),
    })
    return nil
}

func toMetricsDTO(metrics []valueobject.MetricValue) []event.TelemetryMetric {
    out := make([]event.TelemetryMetric, 0, len(metrics))
    for _, mv := range metrics {
        out = append(out, event.TelemetryMetric{
            Metric: string(mv.Metric), Value: mv.Value,
            Unit: mv.Metric.Unit(), Quality: mv.Quality,
        })
    }
    return out
}

func defaultSource(s string) string {
    if s == "" { return "mqtt" }
    return s
}
```

### `application/query_telemetry.go`
```go
package application

import (
    "context"

    "github.com/google/uuid"

    "icmongolang/internal/modules/device/domain/repository"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

type QueryTelemetryUseCase struct {
    telemetryRepo repository.TelemetryRepository
}

func NewQueryTelemetryUseCase(telemetryRepo repository.TelemetryRepository) *QueryTelemetryUseCase {
    return &QueryTelemetryUseCase{telemetryRepo: telemetryRepo}
}

func (uc *QueryTelemetryUseCase) Execute(ctx context.Context, in QueryTelemetryInput) (*TelemetrySeriesResponse, error) {
    tenantID, _ := uuid.Parse(in.TenantID)
    deviceID, _ := uuid.Parse(in.DeviceID)

    metrics := make([]valueobject.Metric, len(in.Metrics))
    for i, m := range in.Metrics {
        metrics[i] = valueobject.Metric(m)
    }

    q := repository.TelemetryQuery{
        TenantID: tenantID, DeviceID: deviceID,
        Metrics: metrics, From: in.From, To: in.To,
        Aggregate: defaultAgg(in.Aggregate),
        Limit: defaultLimit(in.Limit),
    }

    res := &TelemetrySeriesResponse{
        DeviceID: deviceID.String(),
        Series:   map[string][]DataPoint{},
        From:     in.From, To: in.To,
        Aggregated: q.Aggregate != "raw",
    }

    if q.Aggregate == "raw" {
        items, err := uc.telemetryRepo.Query(ctx, q)
        if err != nil { return nil, err }
        for _, t := range items {
            for _, mv := range t.Metrics {
                key := string(mv.Metric)
                res.Series[key] = append(res.Series[key], DataPoint{
                    Timestamp: mv.Timestamp, Value: mv.Value,
                })
            }
        }
        return res, nil
    }

    aggs, err := uc.telemetryRepo.QueryAggregated(ctx, q)
    if err != nil { return nil, err }
    for _, a := range aggs {
        key := string(a.Metric)
        res.Series[key] = append(res.Series[key], DataPoint{
            Timestamp: a.Timestamp, Value: a.Avg,
        })
    }
    return res, nil
}

func defaultAgg(a string) string {
    if a == "" { return "raw" }
    return a
}
func defaultLimit(l int) int {
    if l <= 0 || l > 10000 { return 1000 }
    return l
}
```

### `application/update_shadow.go`
```go
package application

import (
    "context"
    "time"

    "github.com/google/uuid"

    domainerrors "icmongolang/internal/modules/device/domain/errors"
    "icmongolang/internal/modules/device/domain/repository"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
)

type UpdateShadowUseCase struct {
    shadowRepo repository.ShadowRepository
    producer   kafka.Producer
    log        logger.Logger
}

func NewUpdateShadowUseCase(shadowRepo repository.ShadowRepository, producer kafka.Producer, log logger.Logger) *UpdateShadowUseCase {
    return &UpdateShadowUseCase{shadowRepo: shadowRepo, producer: producer, log: log}
}

// UpdateDesired – user ตั้ง desired
func (uc *UpdateShadowUseCase) UpdateDesired(ctx context.Context, in UpdateShadowDesiredInput) (*ShadowResponse, error) {
    tenantID, _ := uuid.Parse(in.TenantID)
    deviceID, _ := uuid.Parse(in.DeviceID)

    s, err := uc.shadowRepo.FindByDeviceID(ctx, tenantID, deviceID)
    if err != nil { return nil, err }

    if in.Merge {
        if err := s.MergeDesired(in.Desired); err != nil { return nil, err }
    } else {
        if err := s.SetDesired(in.Desired); err != nil { return nil, err }
    }

    if err := uc.shadowRepo.Save(ctx, s); err != nil {
        return nil, domainerrors.ErrPersistenceFailure
    }

    _ = uc.producer.Publish(ctx, "device.shadow.updated", deviceID.String(), map[string]any{
        "device_id": deviceID.String(),
        "delta":     s.GetDelta(),
        "version":   s.Version(),
    })
    return toShadowResponse(s), nil
}

// UpdateReported – device ส่ง reported กลับ
func (uc *UpdateShadowUseCase) UpdateReported(ctx context.Context, in UpdateShadowReportedInput) (*ShadowResponse, error) {
    tenantID, _ := uuid.Parse(in.TenantID)
    deviceID, _ := uuid.Parse(in.DeviceID)

    s, err := uc.shadowRepo.FindByDeviceID(ctx, tenantID, deviceID)
    if err != nil { return nil, err }

    if err := s.UpdateReported(in.Reported); err != nil { return nil, err }

    if err := uc.shadowRepo.Save(ctx, s); err != nil {
        return nil, domainerrors.ErrPersistenceFailure
    }
    return toShadowResponse(s), nil
}

// Get – อ่าน shadow
func (uc *UpdateShadowUseCase) Get(ctx context.Context, tenantID, deviceID uuid.UUID) (*ShadowResponse, error) {
    s, err := uc.shadowRepo.FindByDeviceID(ctx, tenantID, deviceID)
    if err != nil { return nil, err }
    return toShadowResponse(s), nil
}

var _ = time.Now
```

### `application/create_alert_rule.go`
```go
package application

import (
    "context"

    "github.com/google/uuid"

    "icmongolang/internal/modules/device/domain/entity"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
    "icmongolang/internal/modules/device/domain/repository"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
    "icmongolang/pkg/logger"
)

type CreateAlertRuleUseCase struct {
    ruleRepo repository.AlertRuleRepository
    log      logger.Logger
}

func NewCreateAlertRuleUseCase(ruleRepo repository.AlertRuleRepository, log logger.Logger) *CreateAlertRuleUseCase {
    return &CreateAlertRuleUseCase{ruleRepo: ruleRepo, log: log}
}

func (uc *CreateAlertRuleUseCase) Execute(ctx context.Context, in CreateAlertRuleInput) (*AlertRuleResponse, error) {
    tenantID, err := uuid.Parse(in.TenantID)
    if err != nil { return nil, domainerrors.ErrInvalidTenantID }
    actorID, _ := uuid.Parse(in.ActorID)

    r, err := entity.NewAlertRule(
        tenantID, in.Name,
        valueobject.Metric(in.Metric), in.Operator, in.Threshold,
        valueobject.AlertSeverity(in.Severity),
    )
    if err != nil { return nil, err }

    r.CreatedBy = actorID
    r.DurationSec = in.DurationSec

    if in.DeviceID != "" {
        did, _ := uuid.Parse(in.DeviceID)
        r.BindToDevice(did)
    }
    if in.GroupID != "" {
        gid, _ := uuid.Parse(in.GroupID)
        r.BindToGroup(gid)
    }
    if in.SiteID != "" {
        sid, _ := uuid.Parse(in.SiteID)
        r.SiteID = &sid
    }

    for _, a := range in.Actions {
        r.Actions = append(r.Actions, entity.AlertAction{
            Type: a.Type, Channel: a.Channel, Target: a.Target, Payload: a.Payload,
        })
    }

    if err := uc.ruleRepo.Save(ctx, r); err != nil {
        return nil, domainerrors.ErrPersistenceFailure
    }
    return toAlertRuleResponse(r), nil
}
```

### `application/create_automation.go`
```go
package application

import (
    "context"

    "github.com/google/uuid"

    "icmongolang/internal/modules/device/domain/entity"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
    "icmongolang/internal/modules/device/domain/repository"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
    "icmongolang/pkg/logger"
)

type CreateAutomationUseCase struct {
    autoRepo repository.AutomationRepository
    log      logger.Logger
}

func NewCreateAutomationUseCase(autoRepo repository.AutomationRepository, log logger.Logger) *CreateAutomationUseCase {
    return &CreateAutomationUseCase{autoRepo: autoRepo, log: log}
}

func (uc *CreateAutomationUseCase) Execute(ctx context.Context, in CreateAutomationInput) (*AutomationResponse, error) {
    tenantID, err := uuid.Parse(in.TenantID)
    if err != nil { return nil, domainerrors.ErrInvalidTenantID }
    actorID, _ := uuid.Parse(in.ActorID)

    trigger := valueobject.AutomationTrigger{
        Type:     valueobject.AutomationTriggerType(in.Trigger.Type),
        CronExpr: in.Trigger.CronExpr,
    }
    if in.Trigger.DeviceID != "" { trigger.DeviceID = &in.Trigger.DeviceID }
    if in.Trigger.GroupID != ""  { trigger.GroupID = &in.Trigger.GroupID }
    if in.Trigger.Metric != "" {
        m := valueobject.Metric(in.Trigger.Metric)
        trigger.Metric = &m
    }

    a, err := entity.NewAutomation(tenantID, actorID, in.Name, trigger)
    if err != nil { return nil, err }
    a.Description = in.Description
    if in.Priority > 0 { a.Priority = in.Priority }

    for _, c := range in.Conditions {
        if err := a.AddCondition(valueobject.Condition{
            Field: c.Field, Operator: c.Operator, Value: c.Value,
        }); err != nil {
            return nil, err
        }
    }
    for _, act := range in.Actions {
        action := valueobject.Action{
            Type: valueobject.ActionType(act.Type),
            Command: act.Command, Payload: act.Payload,
            Channel: act.Channel, Target: act.Target, DelayMs: act.DelayMs,
        }
        if act.DeviceID != "" { action.DeviceID = &act.DeviceID }
        if err := a.AddAction(action); err != nil {
            return nil, err
        }
    }

    if err := uc.autoRepo.Save(ctx, a); err != nil {
        return nil, domainerrors.ErrPersistenceFailure
    }
    return toAutomationResponse(a), nil
}
```

> **หมายเหตุ**: Use case อื่นๆ (ListDevices, GetDevice, UpdateDevice, DecommissionDevice, OTAUpdateFirmware, RunAIInference, ListAlerts, AcknowledgeAlert) ตามรูปแบบเดียวกัน — อยู่ใน PART 3B เต็ม

---

# 🅲 PART 3C — INFRASTRUCTURE LAYER

## C.1 Persistence – Postgres Models

### `infrastructure/persistence/postgres/models.go`
```go
package postgres

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/datatypes"
)

// ============================================================
// DEVICE
// ============================================================

type DeviceModel struct {
    ID               uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID         uuid.UUID      `gorm:"type:uuid;not null;index:idx_dev_tenant_status,priority:1"`
    CustomerID       uuid.UUID      `gorm:"type:uuid;not null;index"`
    SiteID           uuid.UUID      `gorm:"type:uuid;not null;index"`
    GroupID          *uuid.UUID     `gorm:"type:uuid;index"`
    ModelID          *uuid.UUID     `gorm:"type:uuid;index"`
    SerialNo         string         `gorm:"size:100;not null;index:idx_dev_tenant_serial,unique"`
    Name             string         `gorm:"size:255;not null"`
    Type             string         `gorm:"size:30;not null;index"`
    Protocol         string         `gorm:"size:30;not null;index"`
    Status           string         `gorm:"size:20;not null;default:'REGISTERED';index:idx_dev_tenant_status,priority:2"`
    MQTTClientID     string         `gorm:"size:100;index"`
    DeviceTokenHash  string         `gorm:"size:255"`
    FirmwareVersion  string         `gorm:"size:50"`
    FirmwareID       *uuid.UUID     `gorm:"type:uuid"`
    LastSeenAt       *time.Time     `gorm:"index"`
    InstalledAt      *time.Time
    ProvisionedAt    *time.Time
    DecommissionedAt *time.Time
    Capabilities     datatypes.JSON `gorm:"type:jsonb"`
    Tags             datatypes.JSON `gorm:"type:jsonb"`
    Metadata         datatypes.JSON `gorm:"type:jsonb"`
    CreatedBy        uuid.UUID      `gorm:"type:uuid"`
    CreatedAt        time.Time      `gorm:"index"`
    UpdatedAt        time.Time
}

func (DeviceModel) TableName() string { return "device_devices" }

// ============================================================
// SHADOW
// ============================================================

type DeviceShadowModel struct {
    DeviceID  uuid.UUID      `gorm:"type:uuid;primaryKey"`
    TenantID  uuid.UUID      `gorm:"type:uuid;not null;index"`
    State     datatypes.JSON `gorm:"type:jsonb;not null"`
    Version   int64          `gorm:"not null;default:1"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (DeviceShadowModel) TableName() string { return "device_shadows" }

// ============================================================
// COMMAND
// ============================================================

type CommandModel struct {
    ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID      uuid.UUID      `gorm:"type:uuid;not null;index:idx_cmd_tenant_status,priority:1"`
    DeviceID      uuid.UUID      `gorm:"type:uuid;not null;index"`
    Command       string         `gorm:"size:100;not null"`
    Payload       datatypes.JSON `gorm:"type:jsonb"`
    Status        string         `gorm:"size:20;not null;index:idx_cmd_tenant_status,priority:2"`
    Priority      int            `gorm:"default:3"`
    IssuedBy      uuid.UUID      `gorm:"type:uuid"`
    IssuedAt      time.Time      `gorm:"index"`
    SentAt        *time.Time
    AckedAt       *time.Time
    ExpiresAt     time.Time      `gorm:"index"`
    ErrorMsg      string         `gorm:"type:text"`
    ResultPayload datatypes.JSON `gorm:"type:jsonb"`
    RetryCount    int            `gorm:"default:0"`
    MaxRetries    int            `gorm:"default:2"`
}

func (CommandModel) TableName() string { return "device_commands" }

// ============================================================
// ALERT RULE + EVENT
// ============================================================

type AlertRuleModel struct {
    ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID        uuid.UUID      `gorm:"type:uuid;not null;index:idx_alert_tenant_active,priority:1"`
    Name            string         `gorm:"size:255;not null"`
    DeviceID        *uuid.UUID     `gorm:"type:uuid;index"`
    GroupID         *uuid.UUID     `gorm:"type:uuid;index"`
    SiteID          *uuid.UUID     `gorm:"type:uuid;index"`
    Metric          string         `gorm:"size:50;not null;index"`
    Operator        string         `gorm:"size:10;not null"`
    Threshold       float64        `gorm:"type:numeric(15,4)"`
    DurationSec     int            `gorm:"default:0"`
    Severity        string         `gorm:"size:20;not null;index"`
    Actions         datatypes.JSON `gorm:"type:jsonb"`
    IsActive        bool           `gorm:"default:true;index:idx_alert_tenant_active,priority:2"`
    CooldownSec     int            `gorm:"default:300"`
    LastTriggeredAt *time.Time
    CreatedBy       uuid.UUID      `gorm:"type:uuid"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

func (AlertRuleModel) TableName() string { return "device_alert_rules" }

type AlertEventModel struct {
    ID           int64          `gorm:"primaryKey;autoIncrement"`
    TenantID     uuid.UUID      `gorm:"type:uuid;not null;index:idx_alert_evt_tenant,priority:1"`
    RuleID       uuid.UUID      `gorm:"type:uuid;not null;index"`
    DeviceID     uuid.UUID      `gorm:"type:uuid;not null;index"`
    Metric       string         `gorm:"size:50;not null"`
    Value        float64        `gorm:"type:numeric(15,4)"`
    Threshold    float64        `gorm:"type:numeric(15,4)"`
    Operator     string         `gorm:"size:10"`
    Severity     string         `gorm:"size:20;not null;index"`
    Message      string         `gorm:"type:text"`
    Acknowledged bool           `gorm:"default:false;index"`
    AckedBy      *uuid.UUID     `gorm:"type:uuid"`
    AckedAt      *time.Time
    ResolvedAt   *time.Time
    TriggeredAt  time.Time      `gorm:"index:idx_alert_evt_tenant,priority:2"`
}

func (AlertEventModel) TableName() string { return "device_alert_events" }

// ============================================================
// AUTOMATION
// ============================================================

type AutomationModel struct {
    ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID    uuid.UUID      `gorm:"type:uuid;not null;index:idx_auto_tenant_active,priority:1"`
    Name        string         `gorm:"size:255;not null"`
    Description string         `gorm:"type:text"`
    Trigger     datatypes.JSON `gorm:"type:jsonb;not null"`
    Conditions  datatypes.JSON `gorm:"type:jsonb"`
    Actions     datatypes.JSON `gorm:"type:jsonb;not null"`
    IsActive    bool           `gorm:"default:true;index:idx_auto_tenant_active,priority:2"`
    Priority    int            `gorm:"default:5"`
    RunCount    int64          `gorm:"default:0"`
    LastRunAt   *time.Time
    LastError   string         `gorm:"type:text"`
    CreatedBy   uuid.UUID      `gorm:"type:uuid"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func (AutomationModel) TableName() string { return "device_automations" }

// ============================================================
// GROUP / MODEL / FIRMWARE / AI
// ============================================================

type DeviceGroupModel struct {
    ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID  uuid.UUID      `gorm:"type:uuid;not null;index"`
    SiteID    uuid.UUID      `gorm:"type:uuid;not null;index"`
    ParentID  *uuid.UUID     `gorm:"type:uuid;index"`
    Name      string         `gorm:"size:255;not null"`
    GroupType string         `gorm:"size:30"`
    Metadata  datatypes.JSON `gorm:"type:jsonb"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (DeviceGroupModel) TableName() string { return "device_groups" }

type DeviceMasterModel struct {
    ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Vendor          string         `gorm:"size:100"`
    ModelNo         string         `gorm:"size:100;uniqueIndex"`
    Name            string         `gorm:"size:255"`
    DeviceType      string         `gorm:"size:30"`
    Protocol        string         `gorm:"size:30"`
    Capabilities    datatypes.JSON `gorm:"type:jsonb"`
    DefaultConfig   datatypes.JSON `gorm:"type:jsonb"`
    FirmwareChannel string         `gorm:"size:20;default:'stable'"`
    CreatedAt       time.Time
}

func (DeviceMasterModel) TableName() string { return "device_models" }

type FirmwareModel struct {
    ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    ModelID       uuid.UUID `gorm:"type:uuid;not null;index"`
    Version       string    `gorm:"size:50;not null"`
    Channel       string    `gorm:"size:20;not null;index"`
    FileURL       string    `gorm:"type:text"`
    FileSize      int64
    Checksum      string    `gorm:"size:64"`
    ReleaseNotes  string    `gorm:"type:text"`
    MinVersion    string    `gorm:"size:50"`
    IsActive      bool      `gorm:"default:true"`
    ReleasedAt    time.Time
    CreatedAt     time.Time
}

func (FirmwareModel) TableName() string { return "device_firmwares" }

type AIModelModel struct {
    ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID  *uuid.UUID     `gorm:"type:uuid;index"`
    Name      string         `gorm:"size:100;not null"`
    Type      string         `gorm:"size:30;not null"`
    Provider  string         `gorm:"size:30;not null"`
    Endpoint  string         `gorm:"size:255"`
    ModelKey  string         `gorm:"size:100"`
    Config    datatypes.JSON `gorm:"type:jsonb"`
    IsActive  bool           `gorm:"default:true"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (AIModelModel) TableName() string { return "device_ai_models" }
```

### `infrastructure/persistence/postgres/mappers.go`
```go
package postgres

import (
    "encoding/json"

    "gorm.io/datatypes"

    "icmongolang/internal/modules/device/domain/entity"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

func toDeviceModel(d *entity.Device) *DeviceModel {
    m := &DeviceModel{
        ID: d.ID, TenantID: d.TenantID, CustomerID: d.CustomerID, SiteID: d.SiteID,
        GroupID: d.GroupID, ModelID: d.ModelID,
        SerialNo: d.SerialNo.String(), Name: d.Name,
        Type: string(d.Type), Protocol: string(d.Protocol), Status: string(d.Status),
        MQTTClientID: d.MQTTClientID, DeviceTokenHash: d.DeviceTokenHash,
        FirmwareVersion: d.FirmwareVersion, FirmwareID: d.FirmwareID,
        LastSeenAt: d.LastSeenAt, InstalledAt: d.InstalledAt,
        ProvisionedAt: d.ProvisionedAt, DecommissionedAt: d.DecommissionedAt,
        Capabilities: marshalJSON(d.Capabilities),
        Tags:         marshalJSON(d.Tags),
        Metadata:     marshalJSON(d.Metadata),
        CreatedBy:    d.CreatedBy, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
    }
    return m
}

func toDeviceEntity(m *DeviceModel) *entity.Device {
    d := &entity.Device{
        ID: m.ID, TenantID: m.TenantID, CustomerID: m.CustomerID, SiteID: m.SiteID,
        GroupID: m.GroupID, ModelID: m.ModelID,
        SerialNo: valueobject.DeviceSerial(m.SerialNo), Name: m.Name,
        Type: valueobject.DeviceType(m.Type), Protocol: valueobject.Protocol(m.Protocol),
        Status: valueobject.DeviceStatus(m.Status),
        MQTTClientID: m.MQTTClientID, DeviceTokenHash: m.DeviceTokenHash,
        FirmwareVersion: m.FirmwareVersion, FirmwareID: m.FirmwareID,
        LastSeenAt: m.LastSeenAt, InstalledAt: m.InstalledAt,
        ProvisionedAt: m.ProvisionedAt, DecommissionedAt: m.DecommissionedAt,
        CreatedBy: m.CreatedBy, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
    }
    if len(m.Capabilities) > 0 { _ = json.Unmarshal(m.Capabilities, &d.Capabilities) }
    if len(m.Tags) > 0         { _ = json.Unmarshal(m.Tags, &d.Tags) }
    if len(m.Metadata) > 0     { _ = json.Unmarshal(m.Metadata, &d.Metadata) }
    if d.Tags == nil { d.Tags = []string{} }
    if d.Metadata == nil { d.Metadata = map[string]any{} }
    return d
}

func toShadowModel(s *entity.DeviceShadow) *DeviceShadowModel {
    return &DeviceShadowModel{
        DeviceID: s.DeviceID, TenantID: s.TenantID,
        State:     marshalJSON(s.State),
        Version:   s.Version(),
        CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
    }
}

func toShadowEntity(m *DeviceShadowModel) *entity.DeviceShadow {
    s := &entity.DeviceShadow{
        DeviceID: m.DeviceID, TenantID: m.TenantID,
        CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
    }
    if len(m.State) > 0 {
        _ = json.Unmarshal(m.State, &s.State)
    }
    return s
}

func toCommandModel(c *entity.Command) *CommandModel {
    return &CommandModel{
        ID: c.ID, TenantID: c.TenantID, DeviceID: c.DeviceID,
        Command: c.Command, Payload: marshalJSON(c.Payload),
        Status: string(c.Status), Priority: c.Priority,
        IssuedBy: c.IssuedBy, IssuedAt: c.IssuedAt,
        SentAt: c.SentAt, AckedAt: c.AckedAt, ExpiresAt: c.ExpiresAt,
        ErrorMsg: c.ErrorMsg, ResultPayload: marshalJSON(c.ResultPayload),
        RetryCount: c.RetryCount, MaxRetries: c.MaxRetries,
    }
}

func toCommandEntity(m *CommandModel) *entity.Command {
    c := &entity.Command{
        ID: m.ID, TenantID: m.TenantID, DeviceID: m.DeviceID,
        Command: m.Command, Status: valueobject.CommandStatus(m.Status),
        Priority: m.Priority, IssuedBy: m.IssuedBy, IssuedAt: m.IssuedAt,
        SentAt: m.SentAt, AckedAt: m.AckedAt, ExpiresAt: m.ExpiresAt,
        ErrorMsg: m.ErrorMsg, RetryCount: m.RetryCount, MaxRetries: m.MaxRetries,
    }
    if len(m.Payload) > 0 { _ = json.Unmarshal(m.Payload, &c.Payload) }
    if len(m.ResultPayload) > 0 { _ = json.Unmarshal(m.ResultPayload, &c.ResultPayload) }
    return c
}

func marshalJSON(v any) datatypes.JSON {
    if v == nil { return datatypes.JSON([]byte("{}")) }
    b, err := json.Marshal(v)
    if err != nil { return datatypes.JSON([]byte("{}")) }
    return datatypes.JSON(b)
}
```

### `infrastructure/persistence/postgres/device_repository.go`
```go
package postgres

import (
    "context"
    "errors"
    "fmt"
    "strings"
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"

    "icmongolang/internal/modules/device/domain/entity"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
    "icmongolang/internal/modules/device/domain/repository"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

type deviceRepository struct{ db *gorm.DB }

func NewDeviceRepository(db *gorm.DB) repository.DeviceRepository {
    return &deviceRepository{db: db}
}

func (r *deviceRepository) Save(ctx context.Context, d *entity.Device) error {
    m := toDeviceModel(d)
    return r.db.WithContext(ctx).Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "id"}},
        UpdateAll: true,
    }).Create(m).Error
}

func (r *deviceRepository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Device, error) {
    var m DeviceModel
    err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND id = ?", tenantID, id).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrDeviceNotFound
    }
    if err != nil { return nil, err }
    return toDeviceEntity(&m), nil
}

func (r *deviceRepository) FindBySerial(ctx context.Context, tenantID uuid.UUID, serial valueobject.DeviceSerial) (*entity.Device, error) {
    var m DeviceModel
    err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND serial_no = ?", tenantID, serial.String()).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrDeviceNotFound
    }
    if err != nil { return nil, err }
    return toDeviceEntity(&m), nil
}

func (r *deviceRepository) FindByMQTTClientID(ctx context.Context, clientID string) (*entity.Device, error) {
    var m DeviceModel
    err := r.db.WithContext(ctx).
        Where("mqtt_client_id = ?", clientID).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrDeviceNotFound
    }
    if err != nil { return nil, err }
    return toDeviceEntity(&m), nil
}

func (r *deviceRepository) List(ctx context.Context, f repository.DeviceFilter) ([]*entity.Device, int64, error) {
    q := r.db.WithContext(ctx).Model(&DeviceModel{}).Where("tenant_id = ?", f.TenantID)
    if f.CustomerID != nil { q = q.Where("customer_id = ?", *f.CustomerID) }
    if f.SiteID != nil { q = q.Where("site_id = ?", *f.SiteID) }
    if f.GroupID != nil { q = q.Where("group_id = ?", *f.GroupID) }
    if f.ModelID != nil { q = q.Where("model_id = ?", *f.ModelID) }
    if f.Protocol != nil { q = q.Where("protocol = ?", string(*f.Protocol)) }
    if len(f.Statuses) > 0 {
        ss := make([]string, len(f.Statuses))
        for i, s := range f.Statuses { ss[i] = s.String() }
        q = q.Where("status IN ?", ss)
    }
    if len(f.Types) > 0 {
        ts := make([]string, len(f.Types))
        for i, t := range f.Types { ts[i] = string(t) }
        q = q.Where("type IN ?", ts)
    }
    if f.Search != "" {
        s := "%" + strings.ToLower(f.Search) + "%"
        q = q.Where("(LOWER(name) LIKE ? OR LOWER(serial_no) LIKE ?)", s, s)
    }
    if len(f.Tags) > 0 {
        // JSONB @> tags
        for _, tag := range f.Tags {
            q = q.Where("tags @> ?", fmt.Sprintf(`["%s"]`, tag))
        }
    }
    var total int64
    if err := q.Count(&total).Error; err != nil { return nil, 0, err }

    sortBy := sanitizeDeviceSortBy(f.SortBy)
    order := "ASC"
    if strings.EqualFold(f.SortOrder, "desc") { order = "DESC" }
    q = q.Order(fmt.Sprintf("%s %s", sortBy, order))

    if f.PageSize > 0 {
        q = q.Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize)
    }
    var models []DeviceModel
    if err := q.Find(&models).Error; err != nil { return nil, 0, err }
    out := make([]*entity.Device, 0, len(models))
    for i := range models {
        out = append(out, toDeviceEntity(&models[i]))
    }
    return out, total, nil
}

func (r *deviceRepository) ListBySite(ctx context.Context, tenantID, siteID uuid.UUID) ([]*entity.Device, error) {
    var models []DeviceModel
    if err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND site_id = ?", tenantID, siteID).
        Find(&models).Error; err != nil {
        return nil, err
    }
    out := make([]*entity.Device, 0, len(models))
    for i := range models {
        out = append(out, toDeviceEntity(&models[i]))
    }
    return out, nil
}

func (r *deviceRepository) ListByGroup(ctx context.Context, tenantID, groupID uuid.UUID) ([]*entity.Device, error) {
    var models []DeviceModel
    if err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND group_id = ?", tenantID, groupID).
        Find(&models).Error; err != nil {
        return nil, err
    }
    out := make([]*entity.Device, 0, len(models))
    for i := range models {
        out = append(out, toDeviceEntity(&models[i]))
    }
    return out, nil
}

func (r *deviceRepository) ListStale(ctx context.Context, threshold time.Time) ([]*entity.Device, error) {
    var models []DeviceModel
    err := r.db.WithContext(ctx).
        Where("status IN ? AND (last_seen_at IS NULL OR last_seen_at < ?)",
            []string{"ONLINE", "PROVISIONED"}, threshold).
        Limit(500).
        Find(&models).Error
    if err != nil { return nil, err }
    out := make([]*entity.Device, 0, len(models))
    for i := range models {
        out = append(out, toDeviceEntity(&models[i]))
    }
    return out, nil
}

func (r *deviceRepository) CountByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) (int64, error) {
    var n int64
    err := r.db.WithContext(ctx).Model(&DeviceModel{}).
        Where("tenant_id = ? AND customer_id = ? AND status != ?", tenantID, customerID, "DECOMMISSIONED").
        Count(&n).Error
    return n, err
}

func (r *deviceRepository) CountByStatus(ctx context.Context, tenantID uuid.UUID) (map[valueobject.DeviceStatus]int64, error) {
    type row struct {
        Status string
        N      int64
    }
    var rows []row
    err := r.db.WithContext(ctx).Model(&DeviceModel{}).
        Select("status, COUNT(*) as n").
        Where("tenant_id = ?", tenantID).
        Group("status").Scan(&rows).Error
    if err != nil { return nil, err }
    out := map[valueobject.DeviceStatus]int64{}
    for _, r := range rows {
        out[valueobject.DeviceStatus(r.Status)] = r.N
    }
    return out, nil
}

func (r *deviceRepository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
    res := r.db.WithContext(ctx).
        Where("tenant_id = ? AND id = ?", tenantID, id).
        Delete(&DeviceModel{})
    if res.Error != nil { return res.Error }
    if res.RowsAffected == 0 { return domainerrors.ErrDeviceNotFound }
    return nil
}

func sanitizeDeviceSortBy(s string) string {
    allowed := map[string]bool{
        "created_at": true, "updated_at": true, "name": true,
        "serial_no": true, "status": true, "last_seen_at": true,
    }
    if allowed[s] { return s }
    return "created_at"
}
```

### `infrastructure/persistence/postgres/shadow_repository.go`
```go
package postgres

import (
    "context"
    "errors"

    "github.com/google/uuid"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"

    domainerrors "icmongolang/internal/modules/device/domain/errors"
    "icmongolang/internal/modules/device/domain/entity"
    "icmongolang/internal/modules/device/domain/repository"
)

type shadowRepository struct{ db *gorm.DB }

func NewShadowRepository(db *gorm.DB) repository.ShadowRepository {
    return &shadowRepository{db: db}
}

func (r *shadowRepository) Save(ctx context.Context, s *entity.DeviceShadow) error {
    m := toShadowModel(s)
    return r.db.WithContext(ctx).Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "device_id"}},
        UpdateAll: true,
    }).Create(m).Error
}

// SaveWithVersion – optimistic locking
func (r *shadowRepository) SaveWithVersion(ctx context.Context, s *entity.DeviceShadow, expected int64) error {
    m := toShadowModel(s)
    res := r.db.WithContext(ctx).Model(&DeviceShadowModel{}).
        Where("device_id = ? AND version = ?", s.DeviceID, expected).
        Updates(map[string]any{
            "state":      m.State,
            "version":    m.Version,
            "updated_at": m.UpdatedAt,
        })
    if res.Error != nil { return res.Error }
    if res.RowsAffected == 0 {
        return domainerrors.ErrShadowVersionMismatch
    }
    return nil
}

func (r *shadowRepository) FindByDeviceID(ctx context.Context, tenantID, deviceID uuid.UUID) (*entity.DeviceShadow, error) {
    var m DeviceShadowModel
    err := r.db.WithContext(ctx).
        Where("device_id = ? AND tenant_id = ?", deviceID, tenantID).
        First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrShadowNotFound
    }
    if err != nil { return nil, err }
    return toShadowEntity(&m), nil
}

func (r *shadowRepository) Delete(ctx context.Context, tenantID, deviceID uuid.UUID) error {
    return r.db.WithContext(ctx).
        Where("device_id = ? AND tenant_id = ?", deviceID, tenantID).
        Delete(&DeviceShadowModel{}).Error
}
```

### `infrastructure/persistence/postgres/command_repository.go`
```go
package postgres

import (
    "context"
    "errors"
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"

    "icmongolang/internal/modules/device/domain/entity"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
    "icmongolang/internal/modules/device/domain/repository"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

type commandRepository struct{ db *gorm.DB }

func NewCommandRepository(db *gorm.DB) repository.CommandRepository {
    return &commandRepository{db: db}
}

func (r *commandRepository) Save(ctx context.Context, c *entity.Command) error {
    m := toCommandModel(c)
    return r.db.WithContext(ctx).Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "id"}},
        UpdateAll: true,
    }).Create(m).Error
}

func (r *commandRepository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Command, error) {
    var m CommandModel
    err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND id = ?", tenantID, id).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrInvalidCommand
    }
    if err != nil { return nil, err }
    return toCommandEntity(&m), nil
}

func (r *commandRepository) List(ctx context.Context, f repository.CommandFilter) ([]*entity.Command, int64, error) {
    q := r.db.WithContext(ctx).Model(&CommandModel{}).Where("tenant_id = ?", f.TenantID)
    if f.DeviceID != nil { q = q.Where("device_id = ?", *f.DeviceID) }
    if len(f.Status) > 0 {
        ss := make([]string, len(f.Status))
        for i, s := range f.Status { ss[i] = s.String() }
        q = q.Where("status IN ?", ss)
    }
    if f.From != nil { q = q.Where("issued_at >= ?", f.From) }
    if f.To != nil   { q = q.Where("issued_at <= ?", f.To) }

    var total int64
    if err := q.Count(&total).Error; err != nil { return nil, 0, err }

    if f.PageSize > 0 {
        q = q.Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize)
    }
    var models []CommandModel
    if err := q.Order("issued_at DESC").Find(&models).Error; err != nil {
        return nil, 0, err
    }
    out := make([]*entity.Command, 0, len(models))
    for i := range models {
        out = append(out, toCommandEntity(&models[i]))
    }
    return out, total, nil
}

func (r *commandRepository) ListPending(ctx context.Context, deviceID uuid.UUID) ([]*entity.Command, error) {
    return r.listByStatuses(ctx, deviceID,
        []string{"PENDING", "QUEUED", "SENT"})
}

func (r *commandRepository) ListQueuedByDevice(ctx context.Context, deviceID uuid.UUID) ([]*entity.Command, error) {
    return r.listByStatuses(ctx, deviceID, []string{"QUEUED"})
}

func (r *commandRepository) listByStatuses(ctx context.Context, deviceID uuid.UUID, statuses []string) ([]*entity.Command, error) {
    var models []CommandModel
    if err := r.db.WithContext(ctx).
        Where("device_id = ? AND status IN ?", deviceID, statuses).
        Order("priority DESC, issued_at ASC").
        Find(&models).Error; err != nil {
        return nil, err
    }
    out := make([]*entity.Command, 0, len(models))
    for i := range models {
        out = append(out, toCommandEntity(&models[i]))
    }
    return out, nil
}

func (r *commandRepository) ListExpired(ctx context.Context, asOf time.Time, limit int) ([]*entity.Command, error) {
    if limit <= 0 { limit = 100 }
    var models []CommandModel
    err := r.db.WithContext(ctx).
        Where("status IN ? AND expires_at < ?", []string{"PENDING", "QUEUED", "SENT"}, asOf).
        Limit(limit).Find(&models).Error
    if err != nil { return nil, err }
    out := make([]*entity.Command, 0, len(models))
    for i := range models {
        out = append(out, toCommandEntity(&models[i]))
    }
    return out, nil
}

var _ = valueobject.CommandStatusPending
```

### `infrastructure/persistence/postgres/alert_repository.go`
```go
package postgres

import (
    "context"
    "errors"

    "github.com/google/uuid"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"

    "icmongolang/internal/modules/device/domain/entity"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
    "icmongolang/internal/modules/device/domain/repository"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

type alertRuleRepository struct{ db *gorm.DB }

func NewAlertRuleRepository(db *gorm.DB) repository.AlertRuleRepository {
    return &alertRuleRepository{db: db}
}

func (r *alertRuleRepository) Save(ctx context.Context, rule *entity.AlertRule) error {
    m := toAlertRuleModel(rule)
    return r.db.WithContext(ctx).Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "id"}},
        UpdateAll: true,
    }).Create(m).Error
}

func (r *alertRuleRepository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.AlertRule, error) {
    var m AlertRuleModel
    err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND id = ?", tenantID, id).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrAlertRuleNotFound
    }
    if err != nil { return nil, err }
    return toAlertRuleEntity(&m), nil
}

func (r *alertRuleRepository) ListActive(ctx context.Context, tenantID uuid.UUID) ([]*entity.AlertRule, error) {
    var models []AlertRuleModel
    if err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND is_active = true", tenantID).
        Find(&models).Error; err != nil {
        return nil, err
    }
    return toAlertRuleEntities(models), nil
}

func (r *alertRuleRepository) ListActiveByDevice(ctx context.Context, tenantID, deviceID uuid.UUID) ([]*entity.AlertRule, error) {
    var models []AlertRuleModel
    if err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND is_active = true AND (device_id IS NULL OR device_id = ?)", tenantID, deviceID).
        Find(&models).Error; err != nil {
        return nil, err
    }
    return toAlertRuleEntities(models), nil
}

func (r *alertRuleRepository) ListActiveByGroup(ctx context.Context, tenantID, groupID uuid.UUID) ([]*entity.AlertRule, error) {
    var models []AlertRuleModel
    if err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND is_active = true AND (group_id = ? OR group_id IS NULL)", tenantID, groupID).
        Find(&models).Error; err != nil {
        return nil, err
    }
    return toAlertRuleEntities(models), nil
}

func (r *alertRuleRepository) ListByMetric(ctx context.Context, tenantID uuid.UUID, metric valueobject.Metric) ([]*entity.AlertRule, error) {
    var models []AlertRuleModel
    if err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND is_active = true AND metric = ?", tenantID, string(metric)).
        Find(&models).Error; err != nil {
        return nil, err
    }
    return toAlertRuleEntities(models), nil
}

func (r *alertRuleRepository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
    res := r.db.WithContext(ctx).
        Where("tenant_id = ? AND id = ?", tenantID, id).
        Delete(&AlertRuleModel{})
    if res.Error != nil { return res.Error }
    if res.RowsAffected == 0 { return domainerrors.ErrAlertRuleNotFound }
    return nil
}

// ============================================================
// Alert Event Repo
// ============================================================

type alertEventRepository struct{ db *gorm.DB }

func NewAlertEventRepository(db *gorm.DB) repository.AlertEventRepository {
    return &alertEventRepository{db: db}
}

func (r *alertEventRepository) Save(ctx context.Context, e *entity.AlertEvent) error {
    m := &AlertEventModel{
        TenantID: e.TenantID, RuleID: e.RuleID, DeviceID: e.DeviceID,
        Metric: string(e.Metric), Value: e.Value, Threshold: e.Threshold,
        Operator: e.Operator, Severity: string(e.Severity),
        Message: e.Message, TriggeredAt: e.TriggeredAt,
    }
    if err := r.db.WithContext(ctx).Create(m).Error; err != nil { return err }
    e.ID = m.ID
    return nil
}

func (r *alertEventRepository) FindByID(ctx context.Context, id int64) (*entity.AlertEvent, error) {
    var m AlertEventModel
    err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrAlertEventNotFound
    }
    if err != nil { return nil, err }
    return toAlertEventEntity(&m), nil
}

func (r *alertEventRepository) List(ctx context.Context, f repository.AlertEventFilter) ([]*entity.AlertEvent, int64, error) {
    q := r.db.WithContext(ctx).Model(&AlertEventModel{}).Where("tenant_id = ?", f.TenantID)
    if f.DeviceID != nil { q = q.Where("device_id = ?", *f.DeviceID) }
    if f.RuleID != nil { q = q.Where("rule_id = ?", *f.RuleID) }
    if f.Acknowledged != nil { q = q.Where("acknowledged = ?", *f.Acknowledged) }
    if len(f.Severity) > 0 {
        ss := make([]string, len(f.Severity))
        for i, s := range f.Severity { ss[i] = string(s) }
        q = q.Where("severity IN ?", ss)
    }
    if f.From != nil { q = q.Where("triggered_at >= ?", f.From) }
    if f.To != nil   { q = q.Where("triggered_at <= ?", f.To) }

    var total int64
    if err := q.Count(&total).Error; err != nil { return nil, 0, err }

    if f.PageSize > 0 { q = q.Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize) }
    var models []AlertEventModel
    if err := q.Order("triggered_at DESC").Find(&models).Error; err != nil {
        return nil, 0, err
    }
    out := make([]*entity.AlertEvent, 0, len(models))
    for i := range models {
        out = append(out, toAlertEventEntity(&models[i]))
    }
    return out, total, nil
}

func (r *alertEventRepository) CountUnacknowledged(ctx context.Context, tenantID uuid.UUID) (int64, error) {
    var n int64
    err := r.db.WithContext(ctx).Model(&AlertEventModel{}).
        Where("tenant_id = ? AND acknowledged = false", tenantID).Count(&n).Error
    return n, err
}

// ============================================================
// Mappers (alert)
// ============================================================

func toAlertRuleModel(r *entity.AlertRule) *AlertRuleModel {
    return &AlertRuleModel{
        ID: r.ID, TenantID: r.TenantID, Name: r.Name,
        DeviceID: r.DeviceID, GroupID: r.GroupID, SiteID: r.SiteID,
        Metric: string(r.Metric), Operator: r.Operator, Threshold: r.Threshold,
        DurationSec: r.DurationSec, Severity: string(r.Severity),
        Actions: marshalJSON(r.Actions),
        IsActive: r.IsActive, CooldownSec: r.CooldownSec,
        LastTriggeredAt: r.LastTriggeredAt,
        CreatedBy: r.CreatedBy, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
    }
}

func toAlertRuleEntity(m *AlertRuleModel) *entity.AlertRule {
    r := &entity.AlertRule{
        ID: m.ID, TenantID: m.TenantID, Name: m.Name,
        DeviceID: m.DeviceID, GroupID: m.GroupID, SiteID: m.SiteID,
        Metric: valueobject.Metric(m.Metric), Operator: m.Operator,
        Threshold: m.Threshold, DurationSec: m.DurationSec,
        Severity: valueobject.AlertSeverity(m.Severity),
        IsActive: m.IsActive, CooldownSec: m.CooldownSec,
        LastTriggeredAt: m.LastTriggeredAt,
        CreatedBy: m.CreatedBy, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
    }
    if len(m.Actions) > 0 {
        var acts []entity.AlertAction
        _ = jsonUnmarshal(m.Actions, &acts)
        r.Actions = acts
    }
    return r
}

func toAlertRuleEntities(models []AlertRuleModel) []*entity.AlertRule {
    out := make([]*entity.AlertRule, 0, len(models))
    for i := range models {
        out = append(out, toAlertRuleEntity(&models[i]))
    }
    return out
}

func toAlertEventEntity(m *AlertEventModel) *entity.AlertEvent {
    return &entity.AlertEvent{
        ID: m.ID, TenantID: m.TenantID, RuleID: m.RuleID, DeviceID: m.DeviceID,
        Metric: valueobject.Metric(m.Metric), Value: m.Value, Threshold: m.Threshold,
        Operator: m.Operator, Severity: valueobject.AlertSeverity(m.Severity),
        Message: m.Message, Acknowledged: m.Acknowledged,
        AckedBy: m.AckedBy, AckedAt: m.AckedAt, ResolvedAt: m.ResolvedAt,
        TriggeredAt: m.TriggeredAt,
    }
}
```

*(ต้องการ helper `jsonUnmarshal`)*

### `infrastructure/persistence/postgres/automation_repository.go`
```go
package postgres

import (
    "context"
    "errors"

    "github.com/google/uuid"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"

    "icmongolang/internal/modules/device/domain/entity"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
    "icmongolang/internal/modules/device/domain/repository"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

type automationRepository struct{ db *gorm.DB }

func NewAutomationRepository(db *gorm.DB) repository.AutomationRepository {
    return &automationRepository{db: db}
}

func (r *automationRepository) Save(ctx context.Context, a *entity.Automation) error {
    m := toAutomationModel(a)
    return r.db.WithContext(ctx).Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "id"}},
        UpdateAll: true,
    }).Create(m).Error
}

func (r *automationRepository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Automation, error) {
    var m AutomationModel
    err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND id = ?", tenantID, id).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrAutomationNotFound
    }
    if err != nil { return nil, err }
    return toAutomationEntity(&m), nil
}

func (r *automationRepository) ListActiveByTrigger(ctx context.Context, tenantID uuid.UUID, triggerType valueobject.AutomationTriggerType) ([]*entity.Automation, error) {
    var models []AutomationModel
    // Filter by trigger.type in JSONB
    err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND is_active = true", tenantID).
        Where("trigger->>'type' = ?", string(triggerType)).
        Order("priority DESC, created_at ASC").
        Find(&models).Error
    if err != nil { return nil, err }
    return toAutomationEntities(models), nil
}

func (r *automationRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.Automation, error) {
    var models []AutomationModel
    if err := r.db.WithContext(ctx).
        Where("tenant_id = ?", tenantID).
        Order("priority DESC, created_at DESC").
        Find(&models).Error; err != nil {
        return nil, err
    }
    return toAutomationEntities(models), nil
}

func (r *automationRepository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
    res := r.db.WithContext(ctx).
        Where("tenant_id = ? AND id = ?", tenantID, id).
        Delete(&AutomationModel{})
    if res.Error != nil { return res.Error }
    if res.RowsAffected == 0 { return domainerrors.ErrAutomationNotFound }
    return nil
}

func toAutomationModel(a *entity.Automation) *AutomationModel {
    return &AutomationModel{
        ID: a.ID, TenantID: a.TenantID, Name: a.Name, Description: a.Description,
        Trigger: marshalJSON(a.Trigger),
        Conditions: marshalJSON(a.Conditions),
        Actions: marshalJSON(a.Actions),
        IsActive: a.IsActive, Priority: a.Priority,
        RunCount: a.RunCount, LastRunAt: a.LastRunAt, LastError: a.LastError,
        CreatedBy: a.CreatedBy, CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt,
    }
}

func toAutomationEntity(m *AutomationModel) *entity.Automation {
    a := &entity.Automation{
        ID: m.ID, TenantID: m.TenantID, Name: m.Name, Description: m.Description,
        IsActive: m.IsActive, Priority: m.Priority,
        RunCount: m.RunCount, LastRunAt: m.LastRunAt, LastError: m.LastError,
        CreatedBy: m.CreatedBy, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
    }
    if len(m.Trigger) > 0 { _ = jsonUnmarshal(m.Trigger, &a.Trigger) }
    if len(m.Conditions) > 0 { _ = jsonUnmarshal(m.Conditions, &a.Conditions) }
    if len(m.Actions) > 0 { _ = jsonUnmarshal(m.Actions, &a.Actions) }
    return a
}

func toAutomationEntities(models []AutomationModel) []*entity.Automation {
    out := make([]*entity.Automation, 0, len(models))
    for i := range models {
        out = append(out, toAutomationEntity(&models[i]))
    }
    return out
}
```

### `infrastructure/persistence/postgres/helpers.go`
```go
package postgres

import "encoding/json"

func jsonUnmarshal(data []byte, v any) error {
    return json.Unmarshal(data, v)
}
```

---

## C.2 Persistence – InfluxDB (Telemetry)

### `infrastructure/persistence/influxdb/telemetry_repository.go`
```go
package influxdb

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"
    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"

    "icmongolang/internal/modules/device/domain/entity"
    "icmongolang/internal/modules/device/domain/repository"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
)

const (
    measurementTelemetry = "device_telemetry"
)

type telemetryRepo struct {
    client     influxdb2.Client
    writeAPI   api.WriteAPI
    queryAPI   api.QueryAPI
    org        string
    bucket     string
}

func NewTelemetryRepository(client influxdb2.Client, org, bucket string) repository.TelemetryRepository {
    return &telemetryRepo{
        client:   client,
        writeAPI: client.WriteAPI(org, bucket),
        queryAPI: client.QueryAPI(org),
        org:      org,
        bucket:   bucket,
    }
}

func (r *telemetryRepo) Write(ctx context.Context, items ...*entity.Telemetry) error {
    for _, t := range items {
        for _, mv := range t.Metrics {
            p := influxdb2.NewPointWithMeasurement(measurementTelemetry).
                AddTag("tenant_id", t.TenantID.String()).
                AddTag("device_id", t.DeviceID.String()).
                AddTag("site_id", t.SiteID.String()).
                AddTag("metric", string(mv.Metric)).
                AddTag("unit", mv.Metric.Unit()).
                AddTag("quality", mv.Quality).
                AddField("value", mv.Value).
                SetTime(mv.Timestamp)
            r.writeAPI.WritePoint(p)
        }
    }
    // async — ไม่ block
    r.writeAPI.Flush()
    return nil
}

func (r *telemetryRepo) WriteBatch(ctx context.Context, batch *entity.TelemetryBatch) error {
    return r.Write(ctx, batch.Items...)
}

func (r *telemetryRepo) Query(ctx context.Context, q repository.TelemetryQuery) ([]*entity.Telemetry, error) {
    flux := r.buildRawQuery(q)
    result, err := r.queryAPI.Query(ctx, flux)
    if err != nil { return nil, err }

    grouped := map[string]*entity.Telemetry{}
    for result.Next() {
        rec := result.Record()
        deviceID, _ := uuid.Parse(rec.ValueByKey("device_id").(string))
        tenantID, _ := uuid.Parse(rec.ValueByKey("tenant_id").(string))
        metricStr := rec.ValueByKey("metric").(string)
        val, _ := rec.Value().(float64)

        key := fmt.Sprintf("%s-%d", deviceID.String(), rec.Time().UnixNano())
        t, ok := grouped[key]
        if !ok {
            t = &entity.Telemetry{
                TenantID: tenantID, DeviceID: deviceID, Timestamp: rec.Time(),
            }
            grouped[key] = t
        }
        t.Metrics = append(t.Metrics, valueobject.MetricValue{
            Metric: valueobject.Metric(metricStr), Value: val,
            Quality: "GOOD", Timestamp: rec.Time(),
        })
    }
    if err := result.Err(); err != nil { return nil, err }

    out := make([]*entity.Telemetry, 0, len(grouped))
    for _, t := range grouped { out = append(out, t) }
    return out, nil
}

func (r *telemetryRepo) QueryAggregated(ctx context.Context, q repository.TelemetryQuery) ([]*repository.TelemetryAggregate, error) {
    flux := r.buildAggQuery(q)
    result, err := r.queryAPI.Query(ctx, flux)
    if err != nil { return nil, err }

    var out []*repository.TelemetryAggregate
    for result.Next() {
        rec := result.Record()
        metricStr, _ := rec.ValueByKey("metric").(string)
        out = append(out, &repository.TelemetryAggregate{
            Metric: valueobject.Metric(metricStr),
            Timestamp: rec.Time(),
            Min: toF(rec.ValueByKey("min")),
            Max: toF(rec.ValueByKey("max")),
            Avg: toF(rec.ValueByKey("mean")),
            Sum: toF(rec.ValueByKey("sum")),
            Count: toF(rec.ValueByKey("count")),
        })
    }
    return out, result.Err()
}

func (r *telemetryRepo) LatestByDevice(ctx context.Context, tenantID, deviceID uuid.UUID) (map[valueobject.Metric]float64, error) {
    flux := fmt.Sprintf(`
from(bucket: "%s")
  |> range(start: -1h)
  |> filter(fn: (r) => r._measurement == "%s")
  |> filter(fn: (r) => r.tenant_id == "%s")
  |> filter(fn: (r) => r.device_id == "%s")
  |> filter(fn: (r) => r._field == "value")
  |> last()
`, r.bucket, measurementTelemetry, tenantID.String(), deviceID.String())

    result, err := r.queryAPI.Query(ctx, flux)
    if err != nil { return nil, err }

    out := map[valueobject.Metric]float64{}
    for result.Next() {
        metricStr, _ := result.Record().ValueByKey("metric").(string)
        val, _ := result.Record().Value().(float64)
        out[valueobject.Metric(metricStr)] = val
    }
    return out, result.Err()
}

func (r *telemetryRepo) DeleteOlderThan(ctx context.Context, asOf time.Time) error {
    // InfluxDB retention policy จัดการ หรือใช้ delete API
    deleteAPI := r.client.DeleteAPI()
    return deleteAPI.DeleteWithName(ctx, r.org, r.bucket, asOf, time.Now(), "")
}

func (r *telemetryRepo) buildRawQuery(q repository.TelemetryQuery) string {
    filter := fmt.Sprintf(`
  |> filter(fn: (r) => r._measurement == "%s")
  |> filter(fn: (r) => r.tenant_id == "%s")
  |> filter(fn: (r) => r.device_id == "%s")
  |> filter(fn: (r) => r._field == "value")
`, measurementTelemetry, q.TenantID.String(), q.DeviceID.String())

    if len(q.Metrics) > 0 {
        filter += "  |> filter(fn: (r) => "
        for i, m := range q.Metrics {
            if i > 0 { filter += " or " }
            filter += fmt.Sprintf(`r.metric == "%s"`, string(m))
        }
        filter += ")\n"
    }

    limit := q.Limit
    if limit <= 0 { limit = 1000 }

    return fmt.Sprintf(`
from(bucket: "%s")
  |> range(start: %s, stop: %s)
%s
  |> limit(n: %d)
`, r.bucket, q.From.Format(time.RFC3339), q.To.Format(time.RFC3339), filter, limit)
}

func (r *telemetryRepo) buildAggQuery(q repository.TelemetryQuery) string {
    window := aggToWindow(q.Aggregate)
    return fmt.Sprintf(`
from(bucket: "%s")
  |> range(start: %s, stop: %s)
  |> filter(fn: (r) => r._measurement == "%s")
  |> filter(fn: (r) => r.tenant_id == "%s")
  |> filter(fn: (r) => r.device_id == "%s")
  |> filter(fn: (r) => r._field == "value")
  |> aggregateWindow(every: %s, fn: mean, createEmpty: false)
`, r.bucket,
        q.From.Format(time.RFC3339), q.To.Format(time.RFC3339),
        measurementTelemetry, q.TenantID.String(), q.DeviceID.String(),
        window)
}

func aggToWindow(agg string) string {
    switch agg {
    case "1m":  return "1m"
    case "5m":  return "5m"
    case "15m": return "15m"
    case "1h":  return "1h"
    case "1d":  return "1d"
    }
    return "5m"
}

func toF(v any) float64 {
    if f, ok := v.(float64); ok { return f }
    return 0
}
```

---

## C.3 MQTT Integration

### `infrastructure/messaging/mqtt/broker.go`
```go
package mqtt

import (
    "context"
    "fmt"
    "strings"
    "time"

    mqttlib "github.com/eclipse/paho.mqtt.golang"

    "icmongolang/pkg/logger"
)

type MessageHandler func(ctx context.Context, topic string, payload []byte) error

type Broker struct {
    client     mqttlib.Client
    handlers   map[string]MessageHandler
    log        logger.Logger
    brokerURL  string
}

type Config struct {
    BrokerURL      string
    ClientID       string
    Username       string
    Password       string
    ConnectTimeout time.Duration
    KeepAlive      time.Duration
    CleanSession   bool
}

func NewBroker(cfg Config, log logger.Logger) (*Broker, error) {
    opts := mqttlib.NewClientOptions().
        AddBroker(cfg.BrokerURL).
        SetClientID(cfg.ClientID).
        SetUsername(cfg.Username).
        SetPassword(cfg.Password).
        SetAutoReconnect(true).
        SetConnectRetry(true).
        SetConnectTimeout(orDefault(cfg.ConnectTimeout, 10*time.Second)).
        SetKeepAlive(orDefault(cfg.KeepAlive, 30*time.Second)).
        SetCleanSession(cfg.CleanSession)

    b := &Broker{
        handlers:  map[string]MessageHandler{},
        log:       log,
        brokerURL: cfg.BrokerURL,
    }

    opts.SetOnConnectHandler(func(c mqttlib.Client) {
        log.Info("mqtt connected")
        b.resubscribeAll(c)
    })
    opts.SetConnectionLostHandler(func(c mqttlib.Client, err error) {
        log.Warn("mqtt connection lost", "err", err)
    })

    client := mqttlib.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        return nil, token.Error()
    }
    b.client = client
    return b, nil
}

func (b *Broker) Register(topicPattern string, handler MessageHandler) {
    b.handlers[topicPattern] = handler
    if b.client.IsConnected() {
        b.subscribe(topicPattern)
    }
}

func (b *Broker) resubscribeAll(c mqttlib.Client) {
    for pattern := range b.handlers {
        b.subscribe(pattern)
    }
}

func (b *Broker) subscribe(topicPattern string) {
    handler := b.handlers[topicPattern]
    b.client.Subscribe(topicPattern, 1, func(_ mqttlib.Client, msg mqttlib.Message) {
        go func() {
            ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
            defer cancel()
            if err := handler(ctx, msg.Topic(), msg.Payload()); err != nil {
                b.log.Error("mqtt handler error", "topic", msg.Topic(), "err", err)
            }
        }()
    })
}

func (b *Broker) Publish(ctx context.Context, topic, key string, payload any) error {
    bts, err := jsonMarshal(payload)
    if err != nil { return err }
    token := b.client.Publish(topic, 1, false, bts)
    if !token.WaitTimeout(5*time.Second) {
        return fmt.Errorf("mqtt publish timeout: %s", topic)
    }
    return token.Error()
}

func (b *Broker) Disconnect() {
    b.client.Disconnect(250)
}

func orDefault[T comparable](v, def T) T {
    var zero T
    if v == zero { return def }
    return v
}

func jsonMarshal(v any) ([]byte, error) {
    if b, ok := v.([]byte); ok { return b, nil }
    return jsonMarshalAny(v)
}

// Topic helpers
type TopicParts struct {
    TenantID string
    DeviceID string
    Action   string // telemetry | status | cmd | ack | shadow | ota
    Sub      string
}

// ParseTopic – iot/{tenant}/{device}/{action}[/{sub}]
func ParseTopic(topic string) (*TopicParts, error) {
    parts := strings.Split(topic, "/")
    if len(parts) < 4 || parts[0] != "iot" {
        return nil, fmt.Errorf("invalid topic: %s", topic)
    }
    tp := &TopicParts{
        TenantID: parts[1],
        DeviceID: parts[2],
        Action:   parts[3],
    }
    if len(parts) > 4 { tp.Sub = parts[4] }
    return tp, nil
}

func BuildCommandTopic(tenantID, deviceID string) string {
    return fmt.Sprintf("iot/%s/%s/cmd", tenantID, deviceID)
}

func BuildOTACommandTopic(tenantID, deviceID string) string {
    return fmt.Sprintf("iot/%s/%s/ota", tenantID, deviceID)
}
```

### `infrastructure/messaging/mqtt/topic_router.go`
```go
package mqtt

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/google/uuid"

    "icmongolang/internal/modules/device/application"
    "icmongolang/internal/modules/device/domain/repository"
    "icmongolang/pkg/logger"
)

// TopicRouter – route MQTT → Use Case
type TopicRouter struct {
    broker        *Broker
    deviceRepo    repository.DeviceRepository
    ingestUC      *application.IngestTelemetryUseCase
    ackUC         *application.AckCommandUseCase
    shadowUC      *application.UpdateShadowUseCase
    log           logger.Logger
}

func NewTopicRouter(
    broker *Broker,
    deviceRepo repository.DeviceRepository,
    ingestUC *application.IngestTelemetryUseCase,
    ackUC *application.AckCommandUseCase,
    shadowUC *application.UpdateShadowUseCase,
    log logger.Logger,
) *TopicRouter {
    r := &TopicRouter{
        broker: broker, deviceRepo: deviceRepo,
        ingestUC: ingestUC, ackUC: ackUC, shadowUC: shadowUC, log: log,
    }
    r.register()
    return r
}

func (r *TopicRouter) register() {
    r.broker.Register("iot/+/+/telemetry", r.onTelemetry)
    r.broker.Register("iot/+/+/status", r.onStatus)
    r.broker.Register("iot/+/+/cmd/ack", r.onCommandAck)
    r.broker.Register("iot/+/+/shadow/reported", r.onShadowReported)
    r.broker.Register("iot/+/+/fault", r.onFault)
}

// onTelemetry – hot path
func (r *TopicRouter) onTelemetry(ctx context.Context, topic string, payload []byte) error {
    tp, err := ParseTopic(topic)
    if err != nil { return err }

    var msg struct {
        Timestamp int64   `json:"ts"`
        Metrics   []struct {
            Metric  string  `json:"metric"`
            Value   float64 `json:"value"`
            Quality string  `json:"quality,omitempty"`
        } `json:"metrics"`
    }
    if err := json.Unmarshal(payload, &msg); err != nil {
        r.log.Warn("telemetry unmarshal failed", "topic", topic, "err", err)
        return err
    }

    metrics := make([]application.MetricDTO, len(msg.Metrics))
    for i, m := range msg.Metrics {
        metrics[i] = application.MetricDTO{
            Metric: m.Metric, Value: m.Value, Quality: m.Quality,
        }
    }

    return r.ingestUC.Execute(ctx, application.IngestTelemetryInput{
        TenantID: tp.TenantID,
        DeviceID: tp.DeviceID,
        Metrics:  metrics,
        Source:   "mqtt",
    })
}

// onStatus – LWT (Last Will and Testament) — device online/offline
func (r *TopicRouter) onStatus(ctx context.Context, topic string, payload []byte) error {
    tp, err := ParseTopic(topic)
    if err != nil { return err }

    var msg struct {
        Status string `json:"status"` // ONLINE | OFFLINE
    }
    if err := json.Unmarshal(payload, &msg); err != nil { return err }

    tid, _ := uuid.Parse(tp.TenantID)
    did, _ := uuid.Parse(tp.DeviceID)
    d, err := r.deviceRepo.FindByID(ctx, tid, did)
    if err != nil { return err }

    switch msg.Status {
    case "ONLINE":
        if err := d.MarkOnline(now()); err != nil { return err }
    case "OFFLINE":
        if err := d.MarkOffline("lwt"); err != nil { return err }
    }
    return r.deviceRepo.Save(ctx, d)
}

// onCommandAck – device ตอบรับ command
func (r *TopicRouter) onCommandAck(ctx context.Context, topic string, payload []byte) error {
    tp, err := ParseTopic(topic)
    if err != nil { return err }

    var msg struct {
        CommandID string         `json:"command_id"`
        Success   bool           `json:"success"`
        Result    map[string]any `json:"result,omitempty"`
        Error     string         `json:"error,omitempty"`
    }
    if err := json.Unmarshal(payload, &msg); err != nil { return err }

    return r.ackUC.Execute(ctx, application.AckCommandInput{
        TenantID:  tp.TenantID,
        CommandID: msg.CommandID,
        Success:   msg.Success,
        Result:    msg.Result,
        ErrorMsg:  msg.Error,
    })
}

// onShadowReported – device ส่ง reported state
func (r *TopicRouter) onShadowReported(ctx context.Context, topic string, payload []byte) error {
    tp, err := ParseTopic(topic)
    if err != nil { return err }

    var reported map[string]any
    if err := json.Unmarshal(payload, &reported); err != nil { return err }

    _, err = r.shadowUC.UpdateReported(ctx, application.UpdateShadowReportedInput{
        TenantID: tp.TenantID,
        DeviceID: tp.DeviceID,
        Reported: reported,
    })
    return err
}

// onFault – device แจ้ง fault
func (r *TopicRouter) onFault(ctx context.Context, topic string, payload []byte) error {
    tp, err := ParseTopic(topic)
    if err != nil { return err }

    var msg struct {
        Code    string `json:"code"`
        Message string `json:"message"`
    }
    if err := json.Unmarshal(payload, &msg); err != nil { return err }

    tid, _ := uuid.Parse(tp.TenantID)
    did, _ := uuid.Parse(tp.DeviceID)
    d, err := r.deviceRepo.FindByID(ctx, tid, did)
    if err != nil { return err }

    if err := d.ReportFault(msg.Code, msg.Message); err != nil { return err }
    return r.deviceRepo.Save(ctx, d)
}

func now() time.Time { return time.Now() }

var _ = fmt.Sprintf
```

*(ต้อง import `time`)*

### `infrastructure/messaging/mqtt/publisher_adapter.go`
```go
package mqtt

import (
    "context"

    "icmongolang/internal/modules/device/domain/service/port"
)

// PublisherAdapter – implements port.MQTTPublisherPort
type PublisherAdapter struct{ broker *Broker }

func NewPublisherAdapter(b *Broker) port.MQTTPublisherPort {
    return &PublisherAdapter{broker: b}
}

func (p *PublisherAdapter) Publish(ctx context.Context, topic, key string, payload any) error {
    return p.broker.Publish(ctx, topic, key, payload)
}
```

---

## C.4 Kafka Consumers (async processing)

### `infrastructure/messaging/kafka/consumers/telemetry_processor_consumer.go`
```go
package consumers

import (
    "context"
    "encoding/json"
    "log"
    "sync"
    "time"

    "github.com/IBM/sarama"
    "github.com/google/uuid"

    "icmongolang/internal/modules/device/domain/entity"
    "icmongolang/internal/modules/device/domain/event"
    "icmongolang/internal/modules/device/domain/service"
    valueobject "icmongolang/internal/modules/device/domain/value_object"
    "icmongolang/pkg/kafka"
)

// TelemetryProcessorConsumer – consume iot.telemetry.raw → alert + automation + aggregate
type TelemetryProcessorConsumer struct {
    alertEval  *service.AlertEvaluator
    automation *service.AutomationEngine
    producer   kafka.Producer

    // For aggregated output
    mu       sync.Mutex
    windows  map[string]*aggWindow
    flushInt time.Duration
}

type aggWindow struct {
    tenantID   uuid.UUID
    deviceID   uuid.UUID
    customerID uuid.UUID
    metric     string
    count      int
    min, max, sum float64
    start      time.Time
}

func NewTelemetryProcessorConsumer(
    alertEval *service.AlertEvaluator,
    automation *service.AutomationEngine,
    producer kafka.Producer,
) *TelemetryProcessorConsumer {
    c := &TelemetryProcessorConsumer{
        alertEval: alertEval, automation: automation, producer: producer,
        windows:  map[string]*aggWindow{},
        flushInt: 5 * time.Second,
    }
    go c.flusher()
    return c
}

func (c *TelemetryProcessorConsumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (c *TelemetryProcessorConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *TelemetryProcessorConsumer) ConsumeClaim(
    session sarama.ConsumerGroupSession,
    claim sarama.ConsumerGroupClaim,
) error {
    for msg := range claim.Messages() {
        var evt event.TelemetryRaw
        if err := json.Unmarshal(msg.Value, &evt); err != nil {
            log.Printf("[telemetry.processor] unmarshal: %v", err)
            session.MarkMessage(msg, "")
            continue
        }

        ctx := context.Background()
        t := toEntity(&evt)

        // 1. Alert evaluation
        if alerts, err := c.alertEval.Evaluate(ctx, t); err == nil {
            for _, a := range alerts {
                _ = c.producer.Publish(ctx, event.TopicAlertTriggered, a.Event.DeviceID.String(), event.AlertTriggered{
                    EventID:    uuid.New(),
                    AlertID:    a.Event.ID,
                    RuleID:     a.Rule.ID,
                    DeviceID:   a.Event.DeviceID,
                    TenantID:   a.Event.TenantID,
                    Metric:     string(a.Event.Metric),
                    Value:      a.Event.Value,
                    Threshold:  a.Event.Threshold,
                    Severity:   string(a.Event.Severity),
                    Message:    a.Event.Message,
                    OccurredAt: time.Now(),
                })
            }
        }

        // 2. Automation engine
        _ = c.automation.OnTelemetry(ctx, t)

        // 3. Aggregate (in-memory 5s windows)
        c.addToWindow(&evt)

        session.MarkMessage(msg, "")
    }
    return nil
}

func (c *TelemetryProcessorConsumer) addToWindow(evt *event.TelemetryRaw) {
    c.mu.Lock()
    defer c.mu.Unlock()
    for _, m := range evt.Metrics {
        key := evt.DeviceID.String() + ":" + m.Metric
        w, ok := c.windows[key]
        if !ok {
            w = &aggWindow{
                tenantID: evt.TenantID, deviceID: evt.DeviceID,
                metric: m.Metric, start: time.Now(),
                min: m.Value, max: m.Value,
            }
            c.windows[key] = w
        }
        w.count++
        w.sum += m.Value
        if m.Value < w.min { w.min = m.Value }
        if m.Value > w.max { w.max = m.Value }
    }
}

func (c *TelemetryProcessorConsumer) flusher() {
    ticker := time.NewTicker(c.flushInt)
    defer ticker.Stop()
    for range ticker.C {
        c.flush()
    }
}

func (c *TelemetryProcessorConsumer) flush() {
    c.mu.Lock()
    windows := c.windows
    c.windows = map[string]*aggWindow{}
    c.mu.Unlock()

    for _, w := range windows {
        avg := 0.0
        if w.count > 0 { avg = w.sum / float64(w.count) }
        _ = c.producer.Publish(context.Background(), event.TopicTelemetryAggregated, w.deviceID.String(), event.TelemetryAggregated{
            EventID: uuid.New(), TenantID: w.tenantID, DeviceID: w.deviceID,
            Metric: w.metric, Count: w.count,
            Min: w.min, Max: w.max, Avg: avg,
            Window: w.start, OccurredAt: time.Now(),
        })
    }
}

func toEntity(evt *event.TelemetryRaw) *entity.Telemetry {
    metrics := make([]valueobject.MetricValue, 0, len(evt.Metrics))
    for _, m := range evt.Metrics {
        metrics = append(metrics, valueobject.MetricValue{
            Metric: valueobject.Metric(m.Metric), Value: m.Value,
            Quality: m.Quality, Timestamp: evt.Timestamp,
        })
    }
    return &entity.Telemetry{
        TenantID: evt.TenantID, DeviceID: evt.DeviceID, SiteID: evt.SiteID,
        Metrics: metrics, Timestamp: evt.Timestamp, Source: evt.Source,
    }
}
```

### `infrastructure/messaging/kafka/consumers/command_timeout_consumer.go`
```go
package consumers

import (
    "context"
    "log"
    "time"

    "icmongolang/internal/modules/device/domain/repository"
    "icmongolang/pkg/kafka"
)

// CommandTimeoutSweeper – ไม่ใช่ Kafka consumer แต่เป็น scheduler
// หา commands ที่หมดอายุ → mark timeout
type CommandTimeoutSweeper struct {
    commandRepo repository.CommandRepository
    producer    kafka.Producer
}

func NewCommandTimeoutSweeper(repo repository.CommandRepository, producer kafka.Producer) *CommandTimeoutSweeper {
    return &CommandTimeoutSweeper{commandRepo: repo, producer: producer}
}

func (s *CommandTimeoutSweeper) Run(ctx context.Context) error {
    expired, err := s.commandRepo.ListExpired(ctx, time.Now(), 200)
    if err != nil { return err }

    for _, cmd := range expired {
        if err := cmd.Timeout(); err != nil { continue }
        if err := s.commandRepo.Save(ctx, cmd); err != nil {
            log.Printf("[cmd.timeout] save %s: %v", cmd.ID, err)
        }
    }
    return nil
}
```

---

## C.5 Redis Cache (Device State + Shadow)

### `infrastructure/persistence/redis/device_state_cache.go`
```go
package redis

import (
    "context"
    "fmt"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/google/uuid"

    "icmongolang/internal/modules/device/application"
)

type deviceStateCache struct{ client *redis.Client }

func NewDeviceStateCache(client *redis.Client) application.DeviceStateCache {
    return &deviceStateCache{client: client}
}

func (c *deviceStateCache) lastSeenKey(deviceID uuid.UUID) string {
    return fmt.Sprintf("device:last_seen:%s", deviceID)
}

func (c *deviceStateCache) statusKey(deviceID uuid.UUID) string {
    return fmt.Sprintf("device:status:%s", deviceID)
}

func (c *deviceStateCache) SetLastSeen(ctx context.Context, deviceID uuid.UUID, ts time.Time) error {
    return c.client.Set(ctx, c.lastSeenKey(deviceID), ts.Unix(), 1*time.Hour).Err()
}

func (c *deviceStateCache) GetLastSeen(ctx context.Context, deviceID uuid.UUID) (time.Time, error) {
    v, err := c.client.Get(ctx, c.lastSeenKey(deviceID)).Int64()
    if err != nil { return time.Time{}, err }
    return time.Unix(v, 0), nil
}

func (c *deviceStateCache) SetStatus(ctx context.Context, deviceID uuid.UUID, status string) error {
    return c.client.Set(ctx, c.statusKey(deviceID), status, 1*time.Hour).Err()
}
```

### `infrastructure/persistence/redis/shadow_cache.go`
```go
package redis

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/google/uuid"

    "icmongolang/internal/modules/device/domain/entity"
)

// ShadowCache – cache shadow state (hot read path)
type ShadowCache interface {
    Set(ctx context.Context, s *entity.DeviceShadow, ttl time.Duration) error
    Get(ctx context.Context, tenantID, deviceID uuid.UUID) (*entity.DeviceShadow, error)
    Invalidate(ctx context.Context, tenantID, deviceID uuid.UUID) error
}

type shadowCache struct{ client *redis.Client }

func NewShadowCache(client *redis.Client) ShadowCache {
    return &shadowCache{client: client}
}

func (c *shadowCache) key(tenantID, deviceID uuid.UUID) string {
    return fmt.Sprintf("device:shadow:%s:%s", tenantID, deviceID)
}

func (c *shadowCache) Set(ctx context.Context, s *entity.DeviceShadow, ttl time.Duration) error {
    b, err := json.Marshal(s)
    if err != nil { return err }
    return c.client.Set(ctx, c.key(s.TenantID, s.DeviceID), b, ttl).Err()
}

func (c *shadowCache) Get(ctx context.Context, tenantID, deviceID uuid.UUID) (*entity.DeviceShadow, error) {
    b, err := c.client.Get(ctx, c.key(tenantID, deviceID)).Bytes()
    if err != nil { return nil, err }
    var s entity.DeviceShadow
    if err := json.Unmarshal(b, &s); err != nil { return nil, err }
    return &s, nil
}

func (c *shadowCache) Invalidate(ctx context.Context, tenantID, deviceID uuid.UUID) error {
    return c.client.Del(ctx, c.key(tenantID, deviceID)).Err()
}
```

---

## C.6 Elasticsearch Indexer

### `infrastructure/search/elasticsearch/device_indexer.go`
```go
package elasticsearch

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"

    "github.com/elastic/go-elasticsearch/v8"

    "icmongolang/internal/modules/device/application"
)

type DeviceIndexer interface {
    Index(ctx context.Context, tenantID string, dev *application.DeviceResponse) error
    Delete(ctx context.Context, tenantID, id string) error
}

type deviceIndexer struct{ client *elasticsearch.Client }

func NewDeviceIndexer(client *elasticsearch.Client) DeviceIndexer {
    return &deviceIndexer{client: client}
}

func (i *deviceIndexer) indexName(tenantID string) string {
    return fmt.Sprintf("device_devices_%s", tenantID)
}

func (i *deviceIndexer) Index(ctx context.Context, tenantID string, dev *application.DeviceResponse) error {
    b, err := json.Marshal(dev)
    if err != nil { return err }
    _, err = i.client.Index(
        i.indexName(tenantID), bytes.NewReader(b),
        i.client.Index.WithDocumentID(dev.ID),
        i.client.Index.WithContext(ctx),
    )
    return err
}

func (i *deviceIndexer) Delete(ctx context.Context, tenantID, id string) error {
    _, err := i.client.Delete(i.indexName(tenantID), id, i.client.Delete.WithContext(ctx))
    return err
}
```

---

## C.7 AI Inference Service

### `infrastructure/services/ai/ollama_client.go`
```go
package ai

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "icmongolang/internal/modules/device/domain/service/port"
)

type ollamaClient struct {
    baseURL string
    http    *http.Client
}

func NewOllamaClient(baseURL string) port.AIInferencePort {
    return &ollamaClient{
        baseURL: baseURL,
        http:    &http.Client{Timeout: 60 * time.Second},
    }
}

type generateRequest struct {
    Model  string `json:"model"`
    Prompt string `json:"prompt"`
    Stream bool   `json:"stream"`
    Format string `json:"format,omitempty"`
}

type generateResponse struct {
    Response string `json:"response"`
    Done     bool   `json:"done"`
}

func (c *ollamaClient) generate(ctx context.Context, prompt string, format string) (string, error) {
    body := generateRequest{
        Model: "llama3", Prompt: prompt, Stream: false, Format: format,
    }
    b, _ := json.Marshal(body)
    req, _ := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/generate", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.http.Do(req)
    if err != nil { return "", err }
    defer resp.Body.Close()

    var out generateResponse
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return "", err }
    return out.Response, nil
}

func (c *ollamaClient) AnomalyDetect(ctx context.Context, model string, values map[string]float64) (*port.AnomalyResult, error) {
    prompt := fmt.Sprintf(`Analyze these IoT sensor values for anomalies. Reply in JSON: {"is_anomaly": bool, "score": 0-1, "reason": "..."}.
Values: %v`, values)

    raw, err := c.generate(ctx, prompt, "json")
    if err != nil { return nil, err }

    var result struct {
        IsAnomaly bool    `json:"is_anomaly"`
        Score     float64 `json:"score"`
        Reason    string  `json:"reason"`
    }
    if err := json.Unmarshal([]byte(raw), &result); err != nil {
        return nil, err
    }
    return &port.AnomalyResult{
        IsAnomaly: result.IsAnomaly,
        Score:     result.Score,
        Reason:    result.Reason,
    }, nil
}

func (c *ollamaClient) Forecast(ctx context.Context, model string, metric string, horizonHours int, history []float64) (*port.ForecastResult, error) {
    prompt := fmt.Sprintf(`Forecast %s for next %d hours based on history: %v. Reply JSON: {"values": [...], "confidence": 0-1}`,
        metric, horizonHours, history)

    raw, err := c.generate(ctx, prompt, "json")
    if err != nil { return nil, err }

    var result struct {
        Values     []float64 `json:"values"`
        Confidence float64   `json:"confidence"`
    }
    if err := json.Unmarshal([]byte(raw), &result); err != nil {
        return nil, err
    }
    return &port.ForecastResult{Values: result.Values, Confidence: result.Confidence}, nil
}

func (c *ollamaClient) Classify(ctx context.Context, model, text string) (string, float64, error) {
    prompt := fmt.Sprintf(`Classify: "%s". Reply with just the category.`, text)
    raw, err := c.generate(ctx, prompt, "")
    if err != nil { return "", 0, err }
    return raw, 1.0, nil
}
```

---

## C.8 Scheduler Jobs

### `infrastructure/scheduler/offline_detector_job.go`
```go
package scheduler

import (
    "context"
    "log"
    "time"

    "icmongolang/internal/modules/device/domain/repository"
    "icmongolang/internal/modules/device/domain/service"
    "icmongolang/pkg/kafka"

    "icmongolang/internal/modules/device/domain/event"
    "github.com/google/uuid"
)

type OfflineDetectorJob struct {
    deviceRepo repository.DeviceRepository
    producer   kafka.Producer
    threshold  time.Duration
    health     *service.HealthMonitorService
}

func NewOfflineDetectorJob(
    deviceRepo repository.DeviceRepository,
    producer kafka.Producer,
    offlineThreshold time.Duration,
) *OfflineDetectorJob {
    return &OfflineDetectorJob{
        deviceRepo: deviceRepo, producer: producer,
        threshold: offlineThreshold,
        health:    service.NewHealthMonitorService(offlineThreshold, 2*offlineThreshold),
    }
}

func (j *OfflineDetectorJob) Run() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()

    stale, err := j.deviceRepo.ListStale(ctx, time.Now().Add(-j.threshold))
    if err != nil {
        log.Printf("[offline.detector] list stale: %v", err)
        return
    }

    log.Printf("[offline.detector] found %d stale devices", len(stale))
    for _, d := range stale {
        h := j.health.Evaluate(d, time.Now())
        if !h.ShouldMark {
            continue
        }
        if err := d.MarkOffline(h.Reason); err != nil {
            log.Printf("[offline.detector] mark offline %s: %v", d.ID, err)
            continue
        }
        if err := j.deviceRepo.Save(ctx, d); err != nil {
            log.Printf("[offline.detector] save %s: %v", d.ID, err)
            continue
        }
        _ = j.producer.Publish(ctx, event.TopicDeviceOffline, d.ID.String(), event.DeviceOffline{
            EventID: uuid.New(), DeviceID: d.ID, TenantID: d.TenantID,
            Reason: h.Reason, OccurredAt: time.Now(),
        })
    }
}
```

---

# 🅳 PART 3D — INTERFACE + WIRING + MIGRATION + TESTS

## D.1 HTTP Handlers

### `interfaces/http/device_handler.go`
```go
package http

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    "icmongolang/internal/modules/device/application"
)

type DeviceHandler struct {
    registerUC   *application.RegisterDeviceUseCase
    provisionUC  *application.ProvisionDeviceUseCase
    getUC        *application.GetDeviceUseCase
    listUC       *application.ListDevicesUseCase
    updateUC     *application.UpdateDeviceUseCase
    decommUC     *application.DecommissionDeviceUseCase
}

func NewDeviceHandler(
    registerUC *application.RegisterDeviceUseCase,
    provisionUC *application.ProvisionDeviceUseCase,
    getUC *application.GetDeviceUseCase,
    listUC *application.ListDevicesUseCase,
    updateUC *application.UpdateDeviceUseCase,
    decommUC *application.DecommissionDeviceUseCase,
) *DeviceHandler {
    return &DeviceHandler{
        registerUC: registerUC, provisionUC: provisionUC,
        getUC: getUC, listUC: listUC,
        updateUC: updateUC, decommUC: decommUC,
    }
}

func (h *DeviceHandler) Register(c *gin.Context) {
    tid, _ := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    var in application.RegisterDeviceInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
    }
    in.TenantID = tid.String()
    in.ActorID = uid.String()

    res, err := h.registerUC.Execute(c.Request.Context(), in)
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusCreated, res)
}

func (h *DeviceHandler) Provision(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    res, err := h.provisionUC.Execute(c.Request.Context(), application.ProvisionDeviceInput{
        TenantID: tid.String(), DeviceID: id.String(), ActorID: uid.String(),
    })
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusOK, res)
}

func (h *DeviceHandler) Get(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))
    res, err := h.getUC.Execute(c.Request.Context(), tid, id)
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusOK, res)
}

func (h *DeviceHandler) List(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    in := application.ListDevicesInput{
        TenantID:   tid.String(),
        CustomerID: c.Query("customer_id"),
        SiteID:     c.Query("site_id"),
        GroupID:    c.Query("group_id"),
        Statuses:   c.QueryArray("status"),
        Types:      c.QueryArray("type"),
        Protocol:   c.Query("protocol"),
        Search:     c.Query("q"),
        Tags:       c.QueryArray("tag"),
        SortBy:     c.Query("sort_by"),
        SortOrder:  c.DefaultQuery("sort_order", "desc"),
        Page:       parseInt(c.DefaultQuery("page", "1")),
        PageSize:   parseInt(c.DefaultQuery("page_size", "20")),
    }
    res, err := h.listUC.Execute(c.Request.Context(), in)
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusOK, res)
}

func (h *DeviceHandler) Update(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    var in application.UpdateDeviceInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
    }
    in.TenantID = tid.String()
    in.DeviceID = id.String()
    in.ActorID = uid.String()

    res, err := h.updateUC.Execute(c.Request.Context(), in)
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusOK, res)
}

func (h *DeviceHandler) Decommission(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    var req struct {
        Reason string `json:"reason" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
    }
    err := h.decommUC.Execute(c.Request.Context(), application.DecommissionDeviceInput{
        TenantID: tid.String(), DeviceID: id.String(),
        Reason: req.Reason, ActorID: uid.String(),
    })
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusOK, gin.H{"message": "decommissioned"})
}
```

### `interfaces/http/telemetry_handler.go`
```go
package http

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    "icmongolang/internal/modules/device/application"
)

type TelemetryHandler struct {
    ingestUC *application.IngestTelemetryUseCase
    queryUC  *application.QueryTelemetryUseCase
}

func NewTelemetryHandler(ingestUC *application.IngestTelemetryUseCase, queryUC *application.QueryTelemetryUseCase) *TelemetryHandler {
    return &TelemetryHandler{ingestUC: ingestUC, queryUC: queryUC}
}

// Ingest – HTTP fallback (device ที่ไม่ใช้ MQTT)
func (h *TelemetryHandler) Ingest(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    var in application.IngestTelemetryInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
    }
    in.TenantID = tid.String()
    in.DeviceID = id.String()
    in.Source = "http"

    if err := h.ingestUC.Execute(c.Request.Context(), in); err != nil {
        writeError(c, err); return
    }
    c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
}

// Query – GET /devices/:id/telemetry
func (h *TelemetryHandler) Query(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    from := parseTimeOrDefault(c.Query("from"), time.Now().Add(-1*time.Hour))
    to := parseTimeOrDefault(c.Query("to"), time.Now())

    in := application.QueryTelemetryInput{
        TenantID: tid.String(), DeviceID: id.String(),
        Metrics: c.QueryArray("metric"),
        From: from, To: to,
        Aggregate: c.DefaultQuery("agg", "raw"),
        Limit: parseInt(c.DefaultQuery("limit", "1000")),
    }
    res, err := h.queryUC.Execute(c.Request.Context(), in)
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusOK, res)
}

func parseTimeOrDefault(s string, def time.Time) time.Time {
    if s == "" { return def }
    t, err := time.Parse(time.RFC3339, s)
    if err != nil { return def }
    return t
}
```

### `interfaces/http/command_handler.go`
```go
package http

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    "icmongolang/internal/modules/device/application"
)

type CommandHandler struct {
    sendUC *application.SendCommandUseCase
}

func NewCommandHandler(sendUC *application.SendCommandUseCase) *CommandHandler {
    return &CommandHandler{sendUC: sendUC}
}

func (h *CommandHandler) Send(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    var req struct {
        Command  string         `json:"command" binding:"required"`
        Payload  map[string]any `json:"payload"`
        Priority int            `json:"priority,omitempty"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
    }
    res, err := h.sendUC.Execute(c.Request.Context(), application.SendCommandInput{
        TenantID: tid.String(), DeviceID: id.String(),
        Command: req.Command, Payload: req.Payload,
        Priority: req.Priority, ActorID: uid.String(),
    })
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusAccepted, res)
}
```

### `interfaces/http/shadow_handler.go`
```go
package http

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    "icmongolang/internal/modules/device/application"
)

type ShadowHandler struct {
    shadowUC *application.UpdateShadowUseCase
}

func NewShadowHandler(shadowUC *application.UpdateShadowUseCase) *ShadowHandler {
    return &ShadowHandler{shadowUC: shadowUC}
}

func (h *ShadowHandler) Get(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))
    res, err := h.shadowUC.Get(c.Request.Context(), tid, id)
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusOK, res)
}

func (h *ShadowHandler) UpdateDesired(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    var req struct {
        Desired map[string]any `json:"desired" binding:"required"`
        Merge   bool           `json:"merge,omitempty"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
    }
    res, err := h.shadowUC.UpdateDesired(c.Request.Context(), application.UpdateShadowDesiredInput{
        TenantID: tid.String(), DeviceID: id.String(),
        Desired: req.Desired, Merge: req.Merge, ActorID: uid.String(),
    })
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusOK, res)
}
```

### `interfaces/http/alert_handler.go`
```go
package http

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    "icmongolang/internal/modules/device/application"
)

type AlertHandler struct {
    createUC *application.CreateAlertRuleUseCase
    listUC   *application.ListAlertRulesUseCase
    eventsUC *application.ListAlertEventsUseCase
    ackUC    *application.AcknowledgeAlertUseCase
}

func NewAlertHandler(
    createUC *application.CreateAlertRuleUseCase,
    listUC *application.ListAlertRulesUseCase,
    eventsUC *application.ListAlertEventsUseCase,
    ackUC *application.AcknowledgeAlertUseCase,
) *AlertHandler {
    return &AlertHandler{createUC: createUC, listUC: listUC, eventsUC: eventsUC, ackUC: ackUC}
}

func (h *AlertHandler) CreateRule(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    var in application.CreateAlertRuleInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
    }
    in.TenantID = tid.String()
    in.ActorID = uid.String()
    res, err := h.createUC.Execute(c.Request.Context(), in)
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusCreated, res)
}

func (h *AlertHandler) ListEvents(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    in := application.ListAlertEventsInput{
        TenantID: tid.String(),
        DeviceID: c.Query("device_id"),
        Severity: c.QueryArray("severity"),
        Page:     parseInt(c.DefaultQuery("page", "1")),
        PageSize: parseInt(c.DefaultQuery("page_size", "50")),
    }
    res, err := h.eventsUC.Execute(c.Request.Context(), in)
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusOK, res)
}

func (h *AlertHandler) Acknowledge(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    id := parseInt64(c.Param("id"))
    if err := h.ackUC.Execute(c.Request.Context(), application.AckAlertInput{
        TenantID: tid.String(), AlertID: id, ActorID: uid.String(),
    }); err != nil {
        writeError(c, err); return
    }
    c.JSON(http.StatusOK, gin.H{"message": "acknowledged"})
}
```

### `interfaces/http/automation_handler.go`
```go
package http

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    "icmongolang/internal/modules/device/application"
)

type AutomationHandler struct {
    createUC *application.CreateAutomationUseCase
    listUC   *application.ListAutomationsUseCase
    toggleUC *application.ToggleAutomationUseCase
}

func NewAutomationHandler(
    createUC *application.CreateAutomationUseCase,
    listUC *application.ListAutomationsUseCase,
    toggleUC *application.ToggleAutomationUseCase,
) *AutomationHandler {
    return &AutomationHandler{createUC: createUC, listUC: listUC, toggleUC: toggleUC}
}

func (h *AutomationHandler) Create(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    var in application.CreateAutomationInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
    }
    in.TenantID = tid.String()
    in.ActorID = uid.String()
    res, err := h.createUC.Execute(c.Request.Context(), in)
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusCreated, res)
}

func (h *AutomationHandler) List(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    res, err := h.listUC.Execute(c.Request.Context(), tid)
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusOK, res)
}

func (h *AutomationHandler) Toggle(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))
    var req struct {
        Enabled bool `json:"enabled"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
    }
    if err := h.toggleUC.Execute(c.Request.Context(), tid, id, req.Enabled); err != nil {
        writeError(c, err); return
    }
    c.JSON(http.StatusOK, gin.H{"enabled": req.Enabled})
}
```

### `interfaces/http/errors.go`
```go
package http

import (
    "errors"
    "net/http"

    "github.com/gin-gonic/gin"

    domainerrors "icmongolang/internal/modules/device/domain/errors"
)

func writeError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, domainerrors.ErrDeviceNotFound),
        errors.Is(err, domainerrors.ErrAlertRuleNotFound),
        errors.Is(err, domainerrors.ErrAlertEventNotFound),
        errors.Is(err, domainerrors.ErrAutomationNotFound),
        errors.Is(err, domainerrors.ErrShadowNotFound),
        errors.Is(err, domainerrors.ErrAIModelNotFound):
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

    case errors.Is(err, domainerrors.ErrDeviceAlreadyProvisioned),
        errors.Is(err, domainerrors.ErrDeviceAlreadyDecommissioned),
        errors.Is(err, domainerrors.ErrShadowVersionMismatch),
        errors.Is(err, domainerrors.ErrConcurrentModification):
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

    case errors.Is(err, domainerrors.ErrDeviceDecommissioned),
        errors.Is(err, domainerrors.ErrDeviceCannotReceiveCommand):
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})

    case errors.Is(err, domainerrors.ErrInvalidTenantID),
        errors.Is(err, domainerrors.ErrInvalidCustomerID),
        errors.Is(err, domainerrors.ErrInvalidSiteID),
        errors.Is(err, domainerrors.ErrInvalidDeviceName),
        errors.Is(err, domainerrors.ErrInvalidDeviceType),
        errors.Is(err, domainerrors.ErrInvalidSerial),
        errors.Is(err, domainerrors.ErrInvalidProtocol),
        errors.Is(err, domainerrors.ErrInvalidMetric),
        errors.Is(err, domainerrors.ErrInvalidMetricValue),
        errors.Is(err, domainerrors.ErrInvalidCapability),
        errors.Is(err, domainerrors.ErrEmptyCapabilities),
        errors.Is(err, domainerrors.ErrInvalidCommand),
        errors.Is(err, domainerrors.ErrCommandTooLong),
        errors.Is(err, domainerrors.ErrCommandAlreadyTerminal),
        errors.Is(err, domainerrors.ErrInvalidCommandState),
        errors.Is(err, domainerrors.ErrEmptyPatch),
        errors.Is(err, domainerrors.ErrInvalidAlertRuleName),
        errors.Is(err, domainerrors.ErrInvalidOperator),
        errors.Is(err, domainerrors.ErrInvalidSeverity),
        errors.Is(err, domainerrors.ErrInvalidAutomationName),
        errors.Is(err, domainerrors.ErrInvalidTrigger),
        errors.Is(err, domainerrors.ErrTriggerMetricRequired),
        errors.Is(err, domainerrors.ErrTriggerCronRequired),
        errors.Is(err, domainerrors.ErrInvalidCondition),
        errors.Is(err, domainerrors.ErrInvalidAction),
        errors.Is(err, domainerrors.ErrTooManyActions),
        errors.Is(err, domainerrors.ErrTooManyConditions),
        errors.Is(err, domainerrors.ErrInvalidGroupName),
        errors.Is(err, domainerrors.ErrInvalidAIModelName):
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

    case errors.Is(err, domainerrors.ErrMQTTPublishFailure),
        errors.Is(err, domainerrors.ErrInfluxWriteFailure),
        errors.Is(err, domainerrors.ErrPersistenceFailure),
        errors.Is(err, domainerrors.ErrInferenceFailed):
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})

    default:
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
    }
}

func parseInt(s string) int {
    var n int
    for _, r := range s {
        if r < '0' || r > '9' { return 0 }
        n = n*10 + int(r-'0')
    }
    return n
}

func parseInt64(s string) int64 { return int64(parseInt(s)) }
```

### `interfaces/http/routes.go`
```go
package http

import "github.com/gin-gonic/gin"

type Handlers struct {
    Device     *DeviceHandler
    Telemetry  *TelemetryHandler
    Command    *CommandHandler
    Shadow     *ShadowHandler
    Alert      *AlertHandler
    Automation *AutomationHandler
}

func RegisterRoutes(r *gin.RouterGroup, h *Handlers, auth, tenant gin.HandlerFunc) {
    // Devices
    g := r.Group("/devices")
    g.Use(auth, tenant)
    g.POST   ("",                       h.Device.Register)
    g.GET    ("",                       h.Device.List)
    g.GET    ("/:id",                   h.Device.Get)
    g.PUT    ("/:id",                   h.Device.Update)
    g.POST   ("/:id/provision",         h.Device.Provision)
    g.DELETE ("/:id",                   h.Device.Decommission)

    // Telemetry
    g.POST   ("/:id/telemetry",         h.Telemetry.Ingest)
    g.GET    ("/:id/telemetry",         h.Telemetry.Query)

    // Commands
    g.POST   ("/:id/commands",          h.Command.Send)

    // Shadow
    g.GET    ("/:id/shadow",            h.Shadow.Get)
    g.PATCH  ("/:id/shadow/desired",    h.Shadow.UpdateDesired)

    // Alert Rules
    ar := r.Group("/alert-rules")
    ar.Use(auth, tenant)
    ar.POST("",          h.Alert.CreateRule)

    // Alert Events
    ae := r.Group("/alert-events")
    ae.Use(auth, tenant)
    ae.GET   ("",        h.Alert.ListEvents)
    ae.POST  ("/:id/ack", h.Alert.Acknowledge)

    // Automations
    au := r.Group("/automations")
    au.Use(auth, tenant)
    au.POST  ("",         h.Automation.Create)
    au.GET   ("",         h.Automation.List)
    au.PATCH ("/:id/toggle", h.Automation.Toggle)
}
```

## D.2 WebSocket Live Telemetry

### `interfaces/websocket/live_hub.go`
```go
package websocket

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

// LiveHub – per-tenant broadcast
type LiveHub struct {
    mu      sync.RWMutex
    tenants map[string]map[*websocket.Conn]bool
}

func NewLiveHub() *LiveHub {
    return &LiveHub{tenants: map[string]map[*websocket.Conn]bool{}}
}

func (h *LiveHub) Broadcast(tenantID string, event any) {
    b, _ := json.Marshal(event)
    h.mu.RLock()
    conns := h.tenants[tenantID]
    h.mu.RUnlock()
    for c := range conns {
        if err := c.WriteMessage(websocket.TextMessage, b); err != nil {
            log.Printf("ws write: %v", err)
        }
    }
}

func (h *LiveHub) Handle(c *gin.Context, tenantID string) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("ws upgrade: %v", err)
        return
    }
    h.mu.Lock()
    if h.tenants[tenantID] == nil {
        h.tenants[tenantID] = map[*websocket.Conn]bool{}
    }
    h.tenants[tenantID][conn] = true
    h.mu.Unlock()

    defer func() {
        h.mu.Lock()
        delete(h.tenants[tenantID], conn)
        h.mu.Unlock()
        conn.Close()
    }()

    for {
        if _, _, err := conn.ReadMessage(); err != nil { return }
    }
}
```

## D.3 Composition Root

### `module.go`
```go
package device

import (
    "time"

    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
    "gorm.io/gorm"

    "icmongolang/internal/modules/device/application"
    "icmongolang/internal/modules/device/domain/service"
    "icmongolang/internal/modules/device/domain/service/port"
    "icmongolang/internal/modules/device/infrastructure/messaging/kafka/consumers"
    "icmongolang/internal/modules/device/infrastructure/messaging/mqtt"
    "icmongolang/internal/modules/device/infrastructure/persistence/influxdb"
    "icmongolang/internal/modules/device/infrastructure/persistence/postgres"
    redisrepo "icmongolang/internal/modules/device/infrastructure/persistence/redis"
    httpiface "icmongolang/internal/modules/device/interfaces/http"
    wsiface "icmongolang/internal/modules/device/interfaces/websocket"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
    "icmongolang/pkg/transaction"
)

type Dependencies struct {
    DB              *gorm.DB
    Redis           *redis.Client
    Influx          influxdb2.Client
    InfluxOrg       string
    InfluxBucket    string
    Producer        kafka.Producer
    MQTTConfig      mqtt.Config
    AIModel         port.AIInferencePort
    Logger          logger.Logger
    OfflineThreshold time.Duration
}

type Wiring struct {
    LiveHub    *wsiface.LiveHub
    MQTTBroker *mqtt.Broker
}

func Init(router *gin.RouterGroup, deps Dependencies, auth, tenant gin.HandlerFunc) (*Wiring, error) {
    // Repos
    deviceRepo := postgres.NewDeviceRepository(deps.DB)
    shadowRepo := postgres.NewShadowRepository(deps.DB)
    commandRepo := postgres.NewCommandRepository(deps.DB)
    alertRuleRepo := postgres.NewAlertRuleRepository(deps.DB)
    alertEventRepo := postgres.NewAlertEventRepository(deps.DB)
    autoRepo := postgres.NewAutomationRepository(deps.DB)
    groupRepo := postgres.NewDeviceGroupRepository(deps.DB)
    modelRepo := postgres.NewDeviceModelRepository(deps.DB)
    firmwareRepo := postgres.NewFirmwareRepository(deps.DB)

    // InfluxDB
    telemetryRepo := influxdb.NewTelemetryRepository(deps.Influx, deps.InfluxOrg, deps.InfluxBucket)

    // Redis
    stateCache := redisrepo.NewDeviceStateCache(deps.Redis)
    shadowCache := redisrepo.NewShadowCache(deps.Redis)

    // MQTT
    broker, err := mqtt.NewBroker(deps.MQTTConfig, deps.Logger)
    if err != nil { return nil, err }

    // Domain services
    provSvc := service.NewProvisioningService(modelRepo, shadowRepo)
    alertEval := service.NewAlertEvaluator(alertRuleRepo, alertEventRepo)
    actionExec := newActionExecutor(broker, deps.Logger)
    autoEngine := service.NewAutomationEngine(autoRepo, actionExec)
    mqttPublisher := mqtt.NewPublisherAdapter(broker)
    cmdDispatcher := service.NewCommandDispatcher(mqttPublisher)

    // Tx
    txMgr := transaction.NewManager(deps.DB)

    // Entitlement checker (stub — wire to package module in real deployment)
    entChk := newEntitlementCheckerStub()

    // Use cases
    registerUC := application.NewRegisterDeviceUseCase(deviceRepo, modelRepo, entChk, deps.Producer, deps.Logger)
    provisionUC := application.NewProvisionDeviceUseCase(deviceRepo, provSvc, deps.Producer, deps.Logger)
    sendCmdUC := application.NewSendCommandUseCase(deviceRepo, commandRepo, cmdDispatcher, deps.Producer, txMgr, deps.Logger)
    ackCmdUC := application.NewAckCommandUseCase(commandRepo, deps.Producer, deps.Logger)
    ingestUC := application.NewIngestTelemetryUseCase(telemetryRepo, deviceRepo, stateCache, deps.Producer, deps.Logger)
    queryUC := application.NewQueryTelemetryUseCase(telemetryRepo)
    shadowUC := application.NewUpdateShadowUseCase(shadowRepo, deps.Producer, deps.Logger)
    createAlertUC := application.NewCreateAlertRuleUseCase(alertRuleRepo, deps.Logger)
    createAutoUC := application.NewCreateAutomationUseCase(autoRepo, deps.Logger)

    // Get/List/Update — assumed defined in same pattern
    getUC := application.NewGetDeviceUseCase(deviceRepo)
    listUC := application.NewListDevicesUseCase(deviceRepo)
    updateUC := application.NewUpdateDeviceUseCase(deviceRepo, deps.Logger)
    decommUC := application.NewDecommissionDeviceUseCase(deviceRepo, deps.Producer, deps.Logger)

    // MQTT Router
    _ = mqtt.NewTopicRouter(broker, deviceRepo, ingestUC, ackCmdUC, shadowUC, deps.Logger)

    // WebSocket
    liveHub := wsiface.NewLiveHub()

    // Kafka consumer (run in separate process — see cmd/workers)
    _ = consumers.NewTelemetryProcessorConsumer(alertEval, autoEngine, deps.Producer)

    // Handler
    devH := httpiface.NewDeviceHandler(registerUC, provisionUC, getUC, listUC, updateUC, decommUC)
    telH := httpiface.NewTelemetryHandler(ingestUC, queryUC)
    cmdH := httpiface.NewCommandHandler(sendCmdUC)
    shadowH := httpiface.NewShadowHandler(shadowUC)
    alertH := httpiface.NewAlertHandler(createAlertUC,
        application.NewListAlertRulesUseCase(alertRuleRepo),
        application.NewListAlertEventsUseCase(alertEventRepo),
        application.NewAcknowledgeAlertUseCase(alertEventRepo))
    autoH := httpiface.NewAutomationHandler(createAutoUC,
        application.NewListAutomationsUseCase(autoRepo),
        application.NewToggleAutomationUseCase(autoRepo))

    httpiface.RegisterRoutes(router, &httpiface.Handlers{
        Device: devH, Telemetry: telH, Command: cmdH,
        Shadow: shadowH, Alert: alertH, Automation: autoH,
    }, auth, tenant)

    // WS route
    ws := router.Group("/ws")
    ws.Use(auth, tenant)
    ws.GET("/live", func(c *gin.Context) {
        tid := c.MustGet("tenant_id").(string)
        liveHub.Handle(c, tid)
    })

    return &Wiring{LiveHub: liveHub, MQTTBroker: broker}, nil
}

// --- stubs ---

type actionExecutor struct {
    broker *mqtt.Broker
    log    logger.Logger
}

func newActionExecutor(b *mqtt.Broker, log logger.Logger) service.ActionExecutor {
    return &actionExecutor{broker: b, log: log}
}

func (a *actionExecutor) Execute(ctx context.Context, tenantID string, act valueobject.Action) error {
    switch act.Type {
    case valueobject.ActionTypeSendCommand:
        if act.DeviceID == nil { return nil }
        topic := mqtt.BuildCommandTopic(tenantID, *act.DeviceID)
        return a.broker.Publish(ctx, topic, "", map[string]any{
            "command": act.Command, "payload": act.Payload,
        })
    case valueobject.ActionTypeBroadcastWS:
        // broadcast ผ่าน hub (ต้อง inject)
        return nil
    case valueobject.ActionTypeNotify:
        // forward to notifier module ผ่าน Kafka
        return nil
    }
    return nil
}

type entitlementCheckerStub struct{}
func newEntitlementCheckerStub() application.EntitlementChecker { return &entitlementCheckerStub{} }
func (e *entitlementCheckerStub) CheckDeviceQuota(ctx context.Context, tenantID, customerID uuid.UUID) error {
    return nil // TODO: wire to package module
}

var _ = port.MQTTPublisherPort(nil)
```

*(ต้องเพิ่ม imports: `context`, `uuid`, `valueobject`)*

---

## D.4 Worker Entry Points

### `cmd/workers/telemetry/main.go`
```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/IBM/sarama"
    "github.com/joho/godotenv"

    deviceconsumers "icmongolang/internal/modules/device/infrastructure/messaging/kafka/consumers"
)

func main() {
    _ = godotenv.Load()

    // ... init db, kafka producer, repos, services (ดู module.go)
    // แล้วสร้าง consumer

    brokers := []string{os.Getenv("KAFKA_BROKERS")}
    cfg := sarama.NewConfig()
    cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
    cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
    cfg.Consumer.Fetch.Default = 10 * 1024 * 1024 // 10MB — hot path

    group, err := sarama.NewConsumerGroup(brokers, "device-telemetry-processor", cfg)
    if err != nil { log.Fatal(err) }
    defer group.Close()

    // handler ที่ได้จาก module.go wiring
    var handler sarama.ConsumerGroupHandler = newTelemetryHandler()

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() {
        for {
            if err := group.Consume(ctx, []string{"iot.telemetry.raw"}, handler); err != nil {
                log.Printf("consume error: %v", err)
            }
        }
    }()

    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    <-sig
    log.Println("shutting down telemetry worker")
}

func newTelemetryHandler() sarama.ConsumerGroupHandler {
    return deviceconsumers.NewTelemetryProcessorConsumer(nil, nil, nil)
}
```

### `cmd/scheduler/main.go`
```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/joho/godotenv"
    "github.com/robfig/cron/v3"

    devicescheduler "icmongolang/internal/modules/device/infrastructure/scheduler"
)

func main() {
    _ = godotenv.Load()

    // init wiring ...

    c := cron.New(cron.WithSeconds())

    // Offline detection ทุก 1 นาที
    offlineJob := devicescheduler.NewOfflineDetectorJob(nil, nil, 5*time.Minute)
    c.AddFunc("0 * * * * *", offlineJob.Run)

    // Command timeout sweeper ทุก 30 วินาที
    cmdSweeper := devicescheduler.NewCommandTimeoutSweeper(nil, nil)
    c.AddFunc("*/30 * * * * *", func() {
        _ = cmdSweeper.Run(context.Background())
    })

    c.Start()
    log.Println("scheduler started")

    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    <-sig
    c.Stop()
}
```

---

## D.5 Migration SQL

### `migrations/20260105_device_init.sql`
```sql
-- ============================================================
-- device module — initial schema
-- Prefix: device_
-- ============================================================

-- Device master
CREATE TABLE IF NOT EXISTS device_models (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vendor            VARCHAR(100),
    model_no          VARCHAR(100) UNIQUE,
    name              VARCHAR(255),
    device_type       VARCHAR(30),
    protocol          VARCHAR(30),
    capabilities      JSONB DEFAULT '[]'::jsonb,
    default_config    JSONB DEFAULT '{}'::jsonb,
    firmware_channel  VARCHAR(20) DEFAULT 'stable',
    created_at        TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS device_groups (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL,
    site_id     UUID NOT NULL,
    parent_id   UUID REFERENCES device_groups(id),
    name        VARCHAR(255) NOT NULL,
    group_type  VARCHAR(30) DEFAULT 'ZONE',
    metadata    JSONB DEFAULT '{}'::jsonb,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_device_groups_tenant_site ON device_groups(tenant_id, site_id);

-- Devices
CREATE TABLE IF NOT EXISTS device_devices (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL,
    customer_id       UUID NOT NULL,
    site_id           UUID NOT NULL,
    group_id          UUID REFERENCES device_groups(id),
    model_id          UUID REFERENCES device_models(id),
    serial_no         VARCHAR(100) NOT NULL,
    name              VARCHAR(255) NOT NULL,
    type              VARCHAR(30) NOT NULL,
    protocol          VARCHAR(30) NOT NULL,
    status            VARCHAR(20) NOT NULL DEFAULT 'REGISTERED',
    mqtt_client_id    VARCHAR(100),
    device_token_hash VARCHAR(255),
    firmware_version  VARCHAR(50),
    firmware_id       UUID,
    last_seen_at      TIMESTAMP,
    installed_at      TIMESTAMP,
    provisioned_at    TIMESTAMP,
    decommissioned_at TIMESTAMP,
    capabilities      JSONB DEFAULT '[]'::jsonb,
    tags              JSONB DEFAULT '[]'::jsonb,
    metadata          JSONB DEFAULT '{}'::jsonb,
    created_by        UUID,
    created_at        TIMESTAMP DEFAULT NOW(),
    updated_at        TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_device_serial UNIQUE (tenant_id, serial_no)
);
CREATE INDEX idx_dev_tenant_status    ON device_devices(tenant_id, status);
CREATE INDEX idx_dev_tenant_customer  ON device_devices(tenant_id, customer_id);
CREATE INDEX idx_dev_site             ON device_devices(site_id);
CREATE INDEX idx_dev_group            ON device_devices(group_id);
CREATE INDEX idx_dev_mqtt_client      ON device_devices(mqtt_client_id) WHERE mqtt_client_id IS NOT NULL;
CREATE INDEX idx_dev_last_seen        ON device_devices(last_seen_at) WHERE status IN ('ONLINE', 'PROVISIONED');
CREATE INDEX idx_dev_tags_gin         ON device_devices USING GIN (tags);

-- Shadow (Digital Twin)
CREATE TABLE IF NOT EXISTS device_shadows (
    device_id  UUID PRIMARY KEY REFERENCES device_devices(id) ON DELETE CASCADE,
    tenant_id  UUID NOT NULL,
    state      JSONB NOT NULL,
    version    BIGINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_shadows_tenant ON device_shadows(tenant_id);

-- Commands
CREATE TABLE IF NOT EXISTS device_commands (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID NOT NULL,
    device_id      UUID NOT NULL REFERENCES device_devices(id) ON DELETE CASCADE,
    command        VARCHAR(100) NOT NULL,
    payload        JSONB DEFAULT '{}'::jsonb,
    status         VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    priority       INT DEFAULT 3,
    issued_by      UUID,
    issued_at      TIMESTAMP DEFAULT NOW(),
    sent_at        TIMESTAMP,
    acked_at       TIMESTAMP,
    expires_at     TIMESTAMP NOT NULL,
    error_msg      TEXT,
    result_payload JSONB,
    retry_count    INT DEFAULT 0,
    max_retries    INT DEFAULT 2
);
CREATE INDEX idx_cmd_tenant_status    ON device_commands(tenant_id, status);
CREATE INDEX idx_cmd_device           ON device_commands(device_id, issued_at DESC);
CREATE INDEX idx_cmd_expires          ON device_commands(expires_at) WHERE status IN ('PENDING', 'QUEUED', 'SENT');

-- Alert Rules
CREATE TABLE IF NOT EXISTS device_alert_rules (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL,
    name              VARCHAR(255) NOT NULL,
    device_id         UUID REFERENCES device_devices(id) ON DELETE CASCADE,
    group_id          UUID REFERENCES device_groups(id) ON DELETE CASCADE,
    site_id           UUID,
    metric            VARCHAR(50) NOT NULL,
    operator          VARCHAR(10) NOT NULL,
    threshold         NUMERIC(15,4),
    duration_sec      INT DEFAULT 0,
    severity          VARCHAR(20) NOT NULL,
    actions           JSONB DEFAULT '[]'::jsonb,
    is_active         BOOLEAN DEFAULT TRUE,
    cooldown_sec      INT DEFAULT 300,
    last_triggered_at TIMESTAMP,
    created_by        UUID,
    created_at        TIMESTAMP DEFAULT NOW(),
    updated_at        TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_alert_tenant_active ON device_alert_rules(tenant_id, is_active);
CREATE INDEX idx_alert_device        ON device_alert_rules(device_id) WHERE device_id IS NOT NULL;
CREATE INDEX idx_alert_metric        ON device_alert_rules(tenant_id, metric);

-- Alert Events
CREATE TABLE IF NOT EXISTS device_alert_events (
    id            BIGSERIAL PRIMARY KEY,
    tenant_id     UUID NOT NULL,
    rule_id       UUID NOT NULL REFERENCES device_alert_rules(id) ON DELETE CASCADE,
    device_id     UUID NOT NULL,
    metric        VARCHAR(50) NOT NULL,
    value         NUMERIC(15,4),
    threshold     NUMERIC(15,4),
    operator      VARCHAR(10),
    severity      VARCHAR(20) NOT NULL,
    message       TEXT,
    acknowledged  BOOLEAN DEFAULT FALSE,
    acked_by      UUID,
    acked_at      TIMESTAMP,
    resolved_at   TIMESTAMP,
    triggered_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_alert_evt_tenant ON device_alert_events(tenant_id, triggered_at DESC);
CREATE INDEX idx_alert_evt_device ON device_alert_events(device_id, triggered_at DESC);
CREATE INDEX idx_alert_evt_unack  ON device_alert_events(tenant_id, acknowledged) WHERE acknowledged = false;

-- Automations
CREATE TABLE IF NOT EXISTS device_automations (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL,
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    trigger      JSONB NOT NULL,
    conditions   JSONB DEFAULT '[]'::jsonb,
    actions      JSONB NOT NULL,
    is_active    BOOLEAN DEFAULT TRUE,
    priority     INT DEFAULT 5,
    run_count    BIGINT DEFAULT 0,
    last_run_at  TIMESTAMP,
    last_error   TEXT,
    created_by   UUID,
    created_at   TIMESTAMP DEFAULT NOW(),
    updated_at   TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_auto_tenant_active ON device_automations(tenant_id, is_active);
CREATE INDEX idx_auto_trigger_type  ON device_automations((trigger->>'type')) WHERE is_active = true;

-- Firmware
CREATE TABLE IF NOT EXISTS device_firmwares (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    model_id      UUID NOT NULL REFERENCES device_models(id) ON DELETE CASCADE,
    version       VARCHAR(50) NOT NULL,
    channel       VARCHAR(20) NOT NULL,
    file_url      TEXT,
    file_size     BIGINT,
    checksum      VARCHAR(64),
    release_notes TEXT,
    min_version   VARCHAR(50),
    is_active     BOOLEAN DEFAULT TRUE,
    released_at   TIMESTAMP DEFAULT NOW(),
    created_at    TIMESTAMP DEFAULT NOW(),
    UNIQUE(model_id, version, channel)
);

-- AI Models
CREATE TABLE IF NOT EXISTS device_ai_models (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID,
    name       VARCHAR(100) NOT NULL,
    type       VARCHAR(30) NOT NULL,
    provider   VARCHAR(30) NOT NULL,
    endpoint   VARCHAR(255),
    model_key  VARCHAR(100),
    config     JSONB DEFAULT '{}'::jsonb,
    is_active  BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ============================================================
-- InfluxDB setup (รันผ่าน CLI หรือ init script)
-- ============================================================
-- CREATE BUCKET iot_telemetry WITH RETENTION 90d
-- CREATE BUCKET iot_telemetry_5y WITH RETENTION 1825d
-- measurement: device_telemetry
-- tags: tenant_id, device_id, site_id, metric, unit, quality
-- fields: value
```

---

## D.6 .env

```env
# device module
DEVICE_OFFLINE_THRESHOLD_SECONDS=300
DEVICE_COMMAND_TIMEOUT_SECONDS=30
DEVICE_COMMAND_MAX_RETRIES=2
DEVICE_CACHE_TTL_SECONDS=3600
DEVICE_SHADOW_CACHE_TTL_SECONDS=300

# MQTT
MQTT_BROKER=tcp://mosquitto:1883
MQTT_CLIENT_ID=device-platform
MQTT_USERNAME=iot-platform
MQTT_PASSWORD=***
MQTT_KEEPALIVE_SECONDS=30
MQTT_CONNECT_TIMEOUT_SECONDS=10
MQTT_CLEAN_SESSION=false

# InfluxDB
INFLUX_URL=http://influxdb:8086
INFLUX_TOKEN=***
INFLUX_ORG=icmon
INFLUX_BUCKET=iot_telemetry
INFLUX_RETENTION_DAYS=90

# Kafka
KAFKA_BROKERS=kafka:9092
KAFKA_TELEMETRY_GROUP=device-telemetry-processor

# AI
AI_PROVIDER=ollama
AI_ENDPOINT=http://ollama:11434
AI_MODEL=llama3
AI_ENABLED=true
```

## D.7 Run Instructions

```bash
# 1. Migration
psql "$DB_DSN" -f migrations/20260105_device_init.sql

# 2. InfluxDB setup
influx bucket create -n iot_telemetry -o icmon -r 90d
influx bucket create -n iot_telemetry_5y -o icmon -r 1825d

# 3. Build
go build ./internal/modules/device/... ./cmd/...

# 4. Test
go test -v ./internal/modules/device/...

# 5. Run API (with MQTT)
go run ./cmd/api

# 6. Run telemetry worker
go run ./cmd/workers/telemetry

# 7. Run scheduler
go run ./cmd/scheduler

# 8. Smoke test — register + provision
curl -X POST http://localhost:8080/api/v1/devices \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id":"...",
    "site_id":"...",
    "serial_no":"SN-FARM-0001",
    "name":"Soil Sensor 1",
    "type":"SENSOR",
    "protocol":"MQTT"
  }'

# 9. Provision → get token
curl -X POST http://localhost:8080/api/v1/devices/$DEVICE_ID/provision \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: $TENANT_ID"
# → { "device_token": "xxx", "mqtt_client_id": "dev-...", "mqtt_broker": "..." }

# 10. Device publishes telemetry via MQTT
mosquitto_pub -h localhost -p 1883 \
  -u iot-platform -P *** \
  -t "iot/$TENANT_ID/$DEVICE_ID/telemetry" \
  -m '{"ts": 1735000000, "metrics":[{"metric":"temperature","value":28.5},{"metric":"humidity","value":65}]}'

# 11. Query telemetry
curl "http://localhost:8080/api/v1/devices/$DEVICE_ID/telemetry?from=2025-12-01T00:00:00Z&metric=temperature" \
  -H "Authorization: Bearer $TOKEN" -H "X-Tenant-ID: $TENANT_ID"

# 12. Send command
curl -X POST http://localhost:8080/api/v1/devices/$DEVICE_ID/commands \
  -H "Authorization: Bearer $TOKEN" -H "X-Tenant-ID: $TENANT_ID" \
  -d '{"command":"pump_on","payload":{"duration_sec":60}}'
```

---

## D.8 DDD Validation Checklist — device Module

- [x] 6 Aggregate Roots: `Device`, `DeviceShadow`, `Command`, `AlertRule`, `Automation`, `DeviceGroup`
- [x] Shadow ใช้ **optimistic locking** (version)
- [x] Value Objects: 12 ตัว (`DeviceSerial`, `DeviceToken`, `ShadowState`, ...)
- [x] State machine ครบ (`DeviceStatus.CanTransitionTo`, `CommandStatus.IsTerminal`)
- [x] Hot path telemetry → Redis (fast) + InfluxDB (batch) + Kafka (async)
- [x] MQTT topic pattern: `iot/{tenant}/{device}/{action}[/{sub}]`
- [x] InfluxDB สำหรับ time-series (ไม่ใช้ PG)
- [x] AI เป็น **outbound port** (`port.AIInferencePort`)
- [x] Automation engine (ECA pattern) → outbound `ActionExecutor`
- [x] Cross-module ผ่าน Kafka (ไม่ import module อื่นตรง)
- [x] Multi-tenant: `tenant_id` ทุกตาราง, Redis key, MQTT topic, InfluxDB tag
- [x] Domain errors: 50+ sentinel errors
- [x] Idempotency: command ACK, telemetry ingest
- [x] Unit tests: entity (device lifecycle, shadow delta), value objects
- [x] **Device Shadow**: desired/reported/delta pattern ครบ
- [x] **Health monitor**: offline detector เป็น scheduler job
- [x] **Command lifecycle**: PENDING → SENT → ACKED/FAILED/TIMEOUT
- [x] **Retry**: command max retries + sweeper
- [x] Import whitelist ✅

---

# ✅ PART 3 (device) — เสร็จสมบูรณ์

**สถิติ:**
- ไฟล์ทั้งหมด: **~65 ไฟล์**
- Domain: 30 ไฟล์ (12 VO, 11 entities, 10 repos, 5 services+ports, 5 events, 1 error)
- Application: 15 ไฟล์ (12 use cases + DTO + mappers)
- Infrastructure: 15 ไฟล์ (6 repos PG, InfluxDB, MQTT broker+router+publisher, Kafka consumers, Redis x2, ES, AI, Scheduler)
- Interface: 7 ไฟล์ (6 handlers + routes + errors)
- Worker: 2 entry points
- Migration: 9 tables
- Kafka topics: 15
- MQTT topics: 5 patterns

**Pattern พิเศษที่ใช้:**
1. ✅ **Device Shadow (Digital Twin)** — desired/reported/delta + optimistic lock
2. ✅ **MQTT Hot Path** — Redis fast write + InfluxDB batch + Kafka async
3. ✅ **InfluxDB Time-Series** — แยกจาก PostgreSQL
4. ✅ **Command & Control** — state machine + ACK + timeout sweeper
5. ✅ **Automation Engine** — ECA pattern
6. ✅ **AI Inference** — outbound port (Ollama)
7. ✅ **Multi-protocol** — MQTT/LoRaWAN/Modbus/HTTP

---

**PART ถัดไป**: `PART 5 / 7 — Module: iotlogistics` (IoT Logistics)

ต้องการให้ทำ PART 5 ต่อเลยไหม หรือต้องการให้ลงลึกส่วนไหนของ device เพิ่ม (เช่น full code ของ use cases ที่ยังเหลือ, integration tests, docker-compose config)?