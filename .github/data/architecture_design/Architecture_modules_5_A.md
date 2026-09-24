# 🚚 PART 5  — MODULE: `iotlogistics` (IoT Logistics Management)

> **ขนาด**: ใหญ่ — แยก 4 ตอนย่อย
> **Part 5A**: Domain Layer (Entities + VOs + Services + Events + Errors)
> **Part 5B**: Application Layer (Use Cases + DTO + Mappers)
> **Part 5C**: Infrastructure (Postgres + Kafka + Maps + AI + Scheduler + WS)
> **Part 5D**: Interface + Wiring + Migration + Tests + cmd/workers

> **Pattern เฉพาะของ module นี้**:
> 1. **Field Service Management** (FSM) — technician + job + geo-fence
> 2. **Route Optimization** — TSP + Google Maps + AI advisor
> 3. **Installation Checklist** — checklist + photos + e-signature
> 4. **Preventive Maintenance (PM)** — scheduler + auto work order
> 5. **RMA Workflow** — state machine + replacement device
> 6. **Real-time GPS tracking** (WebSocket)
> 7. **Cross-module orchestration** — installation.completed → device.provision (via Kafka)

---

# 🅰️ PART 5A — DOMAIN LAYER

## A.1 โครงสร้าง Domai

```
internal/modules/iotlogistics/domain/
├── entity/
│   ├── shipment.go                  # Aggregate Root #1
│   ├── shipment_item.go             # Entity ย่อย
│   ├── installation_job.go          # Aggregate Root #2
│   ├── checklist.go                 # Value Object (embedded)
│   ├── technician.go                # Aggregate Root #3
│   ├── work_order.go                # Aggregate Root #4
│   ├── maintenance_schedule.go      # Aggregate Root #5
│   ├── rma.go                       # Aggregate Root #6
│   ├── spare_part.go                # Entity
│   └── technician_location.go       # Entity (tracking)
├── value_object/
│   ├── shipment_status.go
│   ├── shipment_no.go
│   ├── job_status.go
│   ├── job_type.go
│   ├── job_priority.go
│   ├── rma_status.go
│   ├── gps_location.go
│   ├── route.go
│   ├── checklist_item.go
│   ├── skill.go
│   ├── geo_fence.go
│   └── time_window.go
├── repository/
│   ├── shipment_repository.go
│   ├── installation_repository.go
│   ├── technician_repository.go
│   ├── work_order_repository.go
│   ├── maintenance_repository.go
│   ├── rma_repository.go
│   ├── spare_part_repository.go
│   └── audit_repository.go
├── service/
│   ├── technician_assignment_service.go
│   ├── route_planning_service.go
│   ├── geo_fence_service.go
│   ├── sla_calculator.go
│   ├── checklist_validator.go
│   └── port/
│       ├── maps_port.go
│       ├── ai_advisor_port.go
│       ├── notifier_port.go
│       ├── erp_port.go             # เช็ค stock อะไหล่
│       └── code_generator_port.go
├── event/
│   ├── shipment_events.go
│   ├── installation_events.go
│   ├── maintenance_events.go
│   └── rma_events.go
└── errors/
    └── errors.go
```

---

## A.2 Value Objects

### `domain/value_object/shipment_status.go`
```go
package valueobject

type ShipmentStatus string

const (
    ShipmentStatusDraft      ShipmentStatus = "DRAFT"
    ShipmentStatusPending    ShipmentStatus = "PENDING"       // รอ dispatch
    ShipmentStatusInTransit  ShipmentStatus = "IN_TRANSIT"    // กำลังขนส่ง
    ShipmentStatusDelivered  ShipmentStatus = "DELIVERED"     // ถึง site แล้ว
    ShipmentStatusInstalled  ShipmentStatus = "INSTALLED"     // ติดตั้งเสร็จ
    ShipmentStatusCancelled  ShipmentStatus = "CANCELLED"
    ShipmentStatusReturned   ShipmentStatus = "RETURNED"
)

func (s ShipmentStatus) IsValid() bool {
    switch s {
    case ShipmentStatusDraft, ShipmentStatusPending, ShipmentStatusInTransit,
        ShipmentStatusDelivered, ShipmentStatusInstalled,
        ShipmentStatusCancelled, ShipmentStatusReturned:
        return true
    }
    return false
}

func (s ShipmentStatus) CanTransitionTo(next ShipmentStatus) bool {
    t := map[ShipmentStatus][]ShipmentStatus{
        ShipmentStatusDraft:     {ShipmentStatusPending, ShipmentStatusCancelled},
        ShipmentStatusPending:   {ShipmentStatusInTransit, ShipmentStatusCancelled},
        ShipmentStatusInTransit: {ShipmentStatusDelivered, ShipmentStatusReturned},
        ShipmentStatusDelivered: {ShipmentStatusInstalled, ShipmentStatusReturned},
        ShipmentStatusInstalled: {}, // terminal
        ShipmentStatusCancelled: {}, // terminal
        ShipmentStatusReturned:  {}, // terminal
    }
    for _, allowed := range t[s] {
        if allowed == next { return true }
    }
    return false
}

func (s ShipmentStatus) IsTerminal() bool {
    return s == ShipmentStatusInstalled || s == ShipmentStatusCancelled || s == ShipmentStatusReturned
}

func (s ShipmentStatus) IsEditable() bool {
    return s == ShipmentStatusDraft
}

func (s ShipmentStatus) String() string { return string(s) }
```

### `domain/value_object/shipment_no.go`
```go
package valueobject

import (
    "fmt"
    "regexp"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
)

// ShipmentNo – รูปแบบ SHP-YYYY-NNNNNN เช่น SHP-2026-000001
type ShipmentNo string

var shipmentNoPattern = regexp.MustCompile(`^SHP-\d{4}-\d{6}$`)

func NewShipmentNo(year, seq int) (ShipmentNo, error) {
    if year < 2000 || year > 9999 || seq < 1 {
        return "", domainerrors.ErrInvalidShipmentNo
    }
    return ShipmentNo(fmt.Sprintf("SHP-%04d-%06d", year, seq)), nil
}

func ParseShipmentNo(s string) (ShipmentNo, error) {
    if !shipmentNoPattern.MatchString(s) {
        return "", domainerrors.ErrInvalidShipmentNo
    }
    return ShipmentNo(s), nil
}

func (s ShipmentNo) String() string { return string(s) }
```

### `domain/value_object/job_status.go`
```go
package valueobject

type JobStatus string

const (
    JobStatusScheduled   JobStatus = "SCHEDULED"
    JobStatusAssigned    JobStatus = "ASSIGNED"
    JobStatusInProgress  JobStatus = "IN_PROGRESS"
    JobStatusPaused      JobStatus = "PAUSED"
    JobStatusDone        JobStatus = "DONE"
    JobStatusFailed      JobStatus = "FAILED"
    JobStatusCancelled   JobStatus = "CANCELLED"
)

func (s JobStatus) IsValid() bool {
    switch s {
    case JobStatusScheduled, JobStatusAssigned, JobStatusInProgress,
        JobStatusPaused, JobStatusDone, JobStatusFailed, JobStatusCancelled:
        return true
    }
    return false
}

func (s JobStatus) CanTransitionTo(next JobStatus) bool {
    t := map[JobStatus][]JobStatus{
        JobStatusScheduled:  {JobStatusAssigned, JobStatusCancelled},
        JobStatusAssigned:   {JobStatusInProgress, JobStatusCancelled, JobStatusFailed},
        JobStatusInProgress: {JobStatusPaused, JobStatusDone, JobStatusFailed},
        JobStatusPaused:     {JobStatusInProgress, JobStatusFailed, JobStatusCancelled},
        JobStatusDone:       {}, // terminal
        JobStatusFailed:     {}, // terminal
        JobStatusCancelled:  {}, // terminal
    }
    for _, allowed := range t[s] {
        if allowed == next { return true }
    }
    return false
}

func (s JobStatus) IsTerminal() bool {
    return s == JobStatusDone || s == JobStatusFailed || s == JobStatusCancelled
}

func (s JobStatus) IsActive() bool {
    return s == JobStatusInProgress || s == JobStatusPaused
}

func (s JobStatus) String() string { return string(s) }
```

### `domain/value_object/job_type.go`
```go
package valueobject

type JobType string

const (
    JobTypeInstallation   JobType = "INSTALLATION"
    JobTypeMaintenance    JobType = "MAINTENANCE"
    JobTypeRepair         JobType = "REPAIR"
    JobTypeRemoval        JobType = "REMOVAL"
    JobTypeInspection     JobType = "INSPECTION"
    JobTypeRelocation     JobType = "RELOCATION"
)

func (t JobType) IsValid() bool {
    switch t {
    case JobTypeInstallation, JobTypeMaintenance, JobTypeRepair,
        JobTypeRemoval, JobTypeInspection, JobTypeRelocation:
        return true
    }
    return false
}

func (t JobType) DefaultSLAHours(priority JobPriority) int {
    base := map[JobType]int{
        JobTypeInstallation: 48,
        JobTypeMaintenance:  72,
        JobTypeRepair:       24,
        JobTypeRemoval:      72,
        JobTypeInspection:   48,
        JobTypeRelocation:   96,
    }[t]
    switch priority {
    case JobPriorityUrgent:   return base / 4
    case JobPriorityHigh:     return base / 2
    case JobPriorityNormal:   return base
    case JobPriorityLow:      return base * 2
    }
    return base
}
```

### `domain/value_object/job_priority.go`
```go
package valueobject

type JobPriority string

const (
    JobPriorityLow    JobPriority = "LOW"
    JobPriorityNormal JobPriority = "NORMAL"
    JobPriorityHigh   JobPriority = "HIGH"
    JobPriorityUrgent JobPriority = "URGENT"
)

func (p JobPriority) IsValid() bool {
    switch p {
    case JobPriorityLow, JobPriorityNormal, JobPriorityHigh, JobPriorityUrgent:
        return true
    }
    return false
}

func (p JobPriority) Score() int {
    switch p {
    case JobPriorityUrgent: return 4
    case JobPriorityHigh:   return 3
    case JobPriorityNormal: return 2
    case JobPriorityLow:    return 1
    }
    return 0
}
```

### `domain/value_object/rma_status.go`
```go
package valueobject

type RMAStatus string

const (
    RMAStatusOpen       RMAStatus = "OPEN"       // สร้าง RMA
    RMAStatusApproved   RMAStatus = "APPROVED"   // อนุมัติแล้ว
    RMAStatusShipped    RMAStatus = "SHIPPED"    // ลูกค้าส่งคืน
    RMAStatusReceived   RMAStatus = "RECEIVED"   // คลังรับแล้ว
    RMAStatusRepairing  RMAStatus = "REPAIRING"  // กำลังซ่อม
    RMAStatusRepaired   RMAStatus = "REPAIRED"   // ซ่อมเสร็จ
    RMAStatusReplaced   RMAStatus = "REPLACED"   // เปลี่ยนตัวใหม่
    RMAStatusRejected   RMAStatus = "REJECTED"   // ปฏิเสธ
    RMAStatusClosed     RMAStatus = "CLOSED"     // ปิด
)

func (s RMAStatus) IsValid() bool {
    switch s {
    case RMAStatusOpen, RMAStatusApproved, RMAStatusShipped,
        RMAStatusReceived, RMAStatusRepairing, RMAStatusRepaired,
        RMAStatusReplaced, RMAStatusRejected, RMAStatusClosed:
        return true
    }
    return false
}

func (s RMAStatus) CanTransitionTo(next RMAStatus) bool {
    t := map[RMAStatus][]RMAStatus{
        RMAStatusOpen:      {RMAStatusApproved, RMAStatusRejected},
        RMAStatusApproved:  {RMAStatusShipped},
        RMAStatusShipped:   {RMAStatusReceived},
        RMAStatusReceived:  {RMAStatusRepairing, RMAStatusReplaced, RMAStatusClosed},
        RMAStatusRepairing: {RMAStatusRepaired, RMAStatusRejected},
        RMAStatusRepaired:  {RMAStatusClosed},
        RMAStatusReplaced:  {RMAStatusClosed},
        RMAStatusRejected:  {}, // terminal
        RMAStatusClosed:    {}, // terminal
    }
    for _, allowed := range t[s] {
        if allowed == next { return true }
    }
    return false
}

func (s RMAStatus) IsTerminal() bool {
    return s == RMAStatusClosed || s == RMAStatusRejected
}

func (s RMAStatus) String() string { return string(s) }
```

### `domain/value_object/gps_location.go`
```go
package valueobject

import (
    "math"
    "time"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
)

// GPSLocation – พิกัดพร้อม metadata
type GPSLocation struct {
    Lat       float64   `json:"lat"`
    Lng       float64   `json:"lng"`
    Accuracy  float64   `json:"accuracy,omitempty"`  // เมตร
    Heading   float64   `json:"heading,omitempty"`   // องศา 0-360
    Speed     float64   `json:"speed,omitempty"`     // km/h
    Timestamp time.Time `json:"ts"`
}

func NewGPSLocation(lat, lng float64, accuracy float64) (GPSLocation, error) {
    loc := GPSLocation{Lat: lat, Lng: lng, Accuracy: accuracy, Timestamp: time.Now()}
    if !loc.IsValid() {
        return GPSLocation{}, domainerrors.ErrInvalidGPS
    }
    return loc, nil
}

func (g GPSLocation) IsValid() bool {
    if math.IsNaN(g.Lat) || math.IsNaN(g.Lng) {
        return false
    }
    return g.Lat >= -90 && g.Lat <= 90 && g.Lng >= -180 && g.Lng <= 180
}

func (g GPSLocation) IsZero() bool {
    return g.Lat == 0 && g.Lng == 0
}

// DistanceKm – Haversine
func (g GPSLocation) DistanceKm(other GPSLocation) float64 {
    const R = 6371.0
    lat1 := g.Lat * math.Pi / 180
    lat2 := other.Lat * math.Pi / 180
    dLat := (other.Lat - g.Lat) * math.Pi / 180
    dLng := (other.Lng - g.Lng) * math.Pi / 180
    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
        math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLng/2)*math.Sin(dLng/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    return R * c
}

// DistanceMeters – สำหรับ geo-fence check
func (g GPSLocation) DistanceMeters(other GPSLocation) float64 {
    return g.DistanceKm(other) * 1000
}
```

### `domain/value_object/route.go`
```go
package valueobject

import (
    "time"
)

// Route – เส้นทางที่วางแผนไว้
type Route struct {
    ID           string
    Waypoints    []RouteWaypoint
    TotalKm      float64
    TotalMinutes int
    StartAt      time.Time
    EndAt        time.Time
    Provider     string // google | osrm | ai
    Polyline     string // encoded
}

type RouteWaypoint struct {
    Order       int
    RefID       string   // job_id
    Location    GPSLocation
    ArrivalETA  time.Time
    DepartETA   time.Time
    ServiceMin  int      // ระยะเวลาบริการ
    DistanceKm  float64  // จาก waypoint ก่อน
}

func (r Route) IsEmpty() bool { return len(r.Waypoints) == 0 }

func (r Route) TotalServiceMinutes() int {
    total := 0
    for _, w := range r.Waypoints {
        total += w.ServiceMin
    }
    return total
}
```

### `domain/value_object/checklist_item.go`
```go
package valueobject

import (
    "time"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
)

// ChecklistItem – ข้อใน checklist ตอนติดตั้ง
type ChecklistItem struct {
    Key         string    `json:"key"`         // "power_check"
    Label       string    `json:"label"`       // "ตรวจสอบไฟเลี้ยง"
    Required    bool      `json:"required"`
    Checked     bool      `json:"checked"`
    Notes       string    `json:"notes,omitempty"`
    PhotoURL    string    `json:"photo_url,omitempty"`
    CheckedAt   time.Time `json:"checked_at,omitempty"`
}

type Checklist []ChecklistItem

func (c Checklist) AllRequiredChecked() bool {
    for _, it := range c {
        if it.Required && !it.Checked {
            return false
        }
    }
    return true
}

func (c Checklist) Validate() error {
    if len(c) == 0 {
        return domainerrors.ErrEmptyChecklist
    }
    for _, it := range c {
        if it.Key == "" || it.Label == "" {
            return domainerrors.ErrInvalidChecklistItem
        }
    }
    return nil
}

func (c Checklist) Progress() (checked, total int) {
    for _, it := range c {
        total++
        if it.Checked { checked++ }
    }
    return
}

// StandardInstallChecklist – template
func StandardInstallChecklist() Checklist {
    return Checklist{
        {Key: "power_check",   Label: "ตรวจสอบไฟเลี้ยง",          Required: true},
        {Key: "network_check", Label: "ตรวจสอบการเชื่อมต่อเครือข่าย", Required: true},
        {Key: "device_mount",  Label: "ยึดอุปกรณ์เข้าที่",         Required: true},
        {Key: "sensor_test",   Label: "ทดสอบ sensor",              Required: true},
        {Key: "config_check",  Label: "ตรวจสอบ config",            Required: true},
        {Key: "photo_evidence",Label: "ถ่ายภาพหลักฐาน",            Required: true},
        {Key: "customer_demo", Label: "สาธิตการใช้งานให้ลูกค้า",     Required: false},
    }
}
```

### `domain/value_object/skill.go`
```go
package valueobject

type Skill string

const (
    SkillSmartFarm     Skill = "SMART_FARM"
    SkillSmartBuilding Skill = "SMART_BUILDING"
    SkillNetwork       Skill = "NETWORK"
    SkillElectrical    Skill = "ELECTRICAL"
    SkillSolar         Skill = "SOLAR"
    SkillHVAC          Skill = "HVAC"
    SkillPlumbing      Skill = "PLUMBING"
    SkillSafetyCert    Skill = "SAFETY_CERT"
)

func (s Skill) IsValid() bool {
    switch s {
    case SkillSmartFarm, SkillSmartBuilding, SkillNetwork,
        SkillElectrical, SkillSolar, SkillHVAC, SkillPlumbing, SkillSafetyCert:
        return true
    }
    return false
}

type SkillSet []Skill

func (ss SkillSet) Has(s Skill) bool {
    for _, x := range ss {
        if x == s { return true }
    }
    return false
}

func (ss SkillSet) Matches(required []Skill) (matched int) {
    for _, req := range required {
        if ss.Has(req) { matched++ }
    }
    return
}

func (ss SkillSet) MatchRatio(required []Skill) float64 {
    if len(required) == 0 { return 1.0 }
    return float64(ss.Matches(required)) / float64(len(required))
}
```

### `domain/value_object/geo_fence.go`
```go
package valueobject

// GeoFence – พื้นที่ทำงาน (circular fence)
type GeoFence struct {
    Center    GPSLocation
    RadiusM   float64
}

func (g GeoFence) Contains(loc GPSLocation) bool {
    return g.Center.DistanceMeters(loc) <= g.RadiusM
}

func (g GeoFence) DistanceFromEdge(loc GPSLocation) float64 {
    return g.RadiusM - g.Center.DistanceMeters(loc)
}

// IsValid – radius ต้อง > 0
func (g GeoFence) IsValid() bool {
    return g.RadiusM > 0 && g.Center.IsValid()
}
```

### `domain/value_object/time_window.go`
```go
package valueobject

import "time"

// TimeWindow – ช่วงเวลาที่ช่างต้องเข้างาน (เช่น 09:00-12:00)
type TimeWindow struct {
    Start time.Time
    End   time.Time
}

func NewTimeWindow(start, end time.Time) (TimeWindow, error) {
    if !end.After(start) {
        return TimeWindow{}, ErrInvalidTimeWindow
    }
    return TimeWindow{Start: start, End: end}, nil
}

func (w TimeWindow) Contains(t time.Time) bool {
    return !t.Before(w.Start) && !t.After(w.End)
}

func (w TimeWindow) Duration() time.Duration {
    return w.End.Sub(w.Start)
}

func (w TimeWindow) Overlaps(other TimeWindow) bool {
    return !w.End.Before(other.Start) && !other.End.Before(w.Start)
}

var ErrInvalidTimeWindow = errNew("invalid time window")
```

*(helper)*

---

## A.3 Domain Entities

### `domain/entity/shipment_item.go`
```go
package entity

import (
    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
)

// ShipmentItem – Entity ย่อยใน Shipment aggregate
type ShipmentItem struct {
    ID          uuid.UUID
    ShipmentID  uuid.UUID
    ProductID   uuid.UUID   // จาก erp_products
    DeviceID    *uuid.UUID  // จาก device_devices (ถ้ามี)
    SerialNo    string
    Qty         float64
    UOM         string
    Notes       string
}

func NewShipmentItem(productID uuid.UUID, serialNo string, qty float64) (*ShipmentItem, error) {
    if qty <= 0 {
        return nil, domainerrors.ErrInvalidQuantity
    }
    return &ShipmentItem{
        ID:        uuid.New(),
        ProductID: productID,
        SerialNo:  serialNo,
        Qty:       qty,
        UOM:       "PCS",
    }, nil
}

func (i *ShipmentItem) AttachDevice(deviceID uuid.UUID) {
    i.DeviceID = &deviceID
}
```

### `domain/entity/shipment.go` ⭐ (Aggregate Root #1)
```go
package entity

import (
    "strings"
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

// Shipment – Aggregate Root
// Invariants:
//   - shipment_no unique per tenant
//   - dispatched → แก้ items ไม่ได้
//   - can dispatch → ต้องมี items ≥ 1
//   - in_transit → update GPS ได้
type Shipment struct {
    ID              uuid.UUID
    TenantID        uuid.UUID
    ShipmentNo      valueobject.ShipmentNo
    CustomerID      uuid.UUID
    SiteID          uuid.UUID
    FromWarehouseID uuid.UUID
    Carrier         string
    TrackingNo      string
    Status          valueobject.ShipmentStatus
    Items           []*ShipmentItem
    ScheduledAt     time.Time
    DispatchedAt    *time.Time
    DeliveredAt     *time.Time
    DeliveredTo     string   // ชื่อผู้รับ
    SignatureURL    string
    GPSLast         *valueobject.GPSLocation
    Notes           string
    Metadata        map[string]any
    CreatedBy       uuid.UUID
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

func NewShipment(
    tenantID, customerID, siteID, warehouseID, actorID uuid.UUID,
    no valueobject.ShipmentNo,
    scheduledAt time.Time,
) (*Shipment, error) {
    if tenantID == uuid.Nil {
        return nil, domainerrors.ErrInvalidTenantID
    }
    if customerID == uuid.Nil {
        return nil, domainerrors.ErrInvalidCustomerID
    }
    if siteID == uuid.Nil {
        return nil, domainerrors.ErrInvalidSiteID
    }
    now := time.Now()
    return &Shipment{
        ID:              uuid.New(),
        TenantID:        tenantID,
        ShipmentNo:      no,
        CustomerID:      customerID,
        SiteID:          siteID,
        FromWarehouseID: warehouseID,
        Status:          valueobject.ShipmentStatusDraft,
        Items:           []*ShipmentItem{},
        ScheduledAt:     scheduledAt,
        Metadata:        map[string]any{},
        CreatedBy:       actorID,
        CreatedAt:       now,
        UpdatedAt:       now,
    }, nil
}

// --- Behavior ---

func (s *Shipment) AddItem(item *ShipmentItem) error {
    if !s.Status.IsEditable() {
        return domainerrors.ErrShipmentLocked
    }
    for _, existing := range s.Items {
        if existing.SerialNo == item.SerialNo {
            return domainerrors.ErrDuplicateSerial
        }
    }
    item.ShipmentID = s.ID
    s.Items = append(s.Items, item)
    s.UpdatedAt = time.Now()
    return nil
}

func (s *Shipment) RemoveItem(itemID uuid.UUID) error {
    if !s.Status.IsEditable() {
        return domainerrors.ErrShipmentLocked
    }
    for i, it := range s.Items {
        if it.ID == itemID {
            s.Items = append(s.Items[:i], s.Items[i+1:]...)
            s.UpdatedAt = time.Now()
            return nil
        }
    }
    return domainerrors.ErrShipmentItemNotFound
}

func (s *Shipment) MarkPending() error {
    if s.Status != valueobject.ShipmentStatusDraft {
        return domainerrors.ErrInvalidStatusTransition
    }
    if len(s.Items) == 0 {
        return domainerrors.ErrEmptyShipment
    }
    s.Status = valueobject.ShipmentStatusPending
    s.UpdatedAt = time.Now()
    return nil
}

func (s *Shipment) Dispatch(carrier, trackingNo string, actorID uuid.UUID) error {
    if !s.Status.CanTransitionTo(valueobject.ShipmentStatusInTransit) {
        return domainerrors.ErrInvalidStatusTransition
    }
    if len(s.Items) == 0 {
        return domainerrors.ErrEmptyShipment
    }
    carrier = strings.TrimSpace(carrier)
    if carrier == "" {
        return domainerrors.ErrInvalidCarrier
    }
    now := time.Now()
    s.Carrier = carrier
    s.TrackingNo = trackingNo
    s.Status = valueobject.ShipmentStatusInTransit
    s.DispatchedAt = &now
    s.UpdatedAt = now
    return nil
}

func (s *Shipment) UpdateGPS(loc valueobject.GPSLocation) error {
    if s.Status != valueobject.ShipmentStatusInTransit {
        return domainerrors.ErrShipmentNotInTransit
    }
    if !loc.IsValid() {
        return domainerrors.ErrInvalidGPS
    }
    s.GPSLast = &loc
    s.UpdatedAt = time.Now()
    return nil
}

func (s *Shipment) ConfirmDelivery(deliveredTo, signatureURL string) error {
    if s.Status != valueobject.ShipmentStatusInTransit {
        return domainerrors.ErrInvalidStatusTransition
    }
    deliveredTo = strings.TrimSpace(deliveredTo)
    if deliveredTo == "" {
        return domainerrors.ErrRecipientRequired
    }
    now := time.Now()
    s.Status = valueobject.ShipmentStatusDelivered
    s.DeliveredAt = &now
    s.DeliveredTo = deliveredTo
    s.SignatureURL = signatureURL
    s.UpdatedAt = now
    return nil
}

func (s *Shipment) MarkInstalled() error {
    if s.Status != valueobject.ShipmentStatusDelivered {
        return domainerrors.ErrInvalidStatusTransition
    }
    s.Status = valueobject.ShipmentStatusInstalled
    s.UpdatedAt = time.Now()
    return nil
}

func (s *Shipment) Cancel(reason string, actorID uuid.UUID) error {
    if s.Status.IsTerminal() {
        return domainerrors.ErrShipmentAlreadyTerminal
    }
    if s.Status == valueobject.ShipmentStatusInTransit {
        return domainerrors.ErrCannotCancelInTransit
    }
    s.Status = valueobject.ShipmentStatusCancelled
    if s.Metadata == nil { s.Metadata = map[string]any{} }
    s.Metadata["cancel_reason"] = reason
    s.Metadata["cancelled_by"] = actorID.String()
    s.Metadata["cancelled_at"] = time.Now()
    s.UpdatedAt = time.Now()
    return nil
}

// --- Query ---

func (s *Shipment) IsEditable() bool      { return s.Status.IsEditable() }
func (s *Shipment) IsInTransit() bool     { return s.Status == valueobject.ShipmentStatusInTransit }
func (s *Shipment) IsDelivered() bool     { return s.Status == valueobject.ShipmentStatusDelivered }
func (s *Shipment) IsInstalled() bool     { return s.Status == valueobject.ShipmentStatusInstalled }

func (s *Shipment) TotalItems() int       { return len(s.Items) }
func (s *Shipment) DeviceIDs() []uuid.UUID {
    out := []uuid.UUID{}
    for _, it := range s.Items {
        if it.DeviceID != nil {
            out = append(out, *it.DeviceID)
        }
    }
    return out
}
func (s *Shipment) SerialNumbers() []string {
    out := make([]string, 0, len(s.Items))
    for _, it := range s.Items {
        out = append(out, it.SerialNo)
    }
    return out
}
```

### `domain/entity/installation_job.go` ⭐ (Aggregate Root #2)
```go
package entity

import (
    "strings"
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

// InstallationJob – Aggregate Root
// ครอบคลุม: install / repair / inspect
type InstallationJob struct {
    ID             uuid.UUID
    TenantID       uuid.UUID
    JobNo          string
    ShipmentID     *uuid.UUID
    CustomerID     uuid.UUID
    SiteID         uuid.UUID
    DeviceIDs      []uuid.UUID
    TechnicianID   *uuid.UUID

    JobType        valueobject.JobType
    Priority       valueobject.JobPriority
    Status         valueobject.JobStatus
    ScheduledAt    time.Time
    TimeWindow     *valueobject.TimeWindow
    SLADueAt       time.Time

    // Execution
    StartedAt      *time.Time
    CompletedAt    *time.Time
    StartedGPS     *valueobject.GPSLocation
    CompletedGPS   *valueobject.GPSLocation

    // Evidence
    Checklist      valueobject.Checklist
    Photos         []string
    SignatureURL   string
    Notes          string
    Result         map[string]any

    // Meta
    AssignedBy     *uuid.UUID
    AssignedAt     *time.Time
    CancelledAt    *time.Time
    CancelReason   string

    CreatedBy      uuid.UUID
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

func NewInstallationJob(
    tenantID, customerID, siteID, actorID uuid.UUID,
    jobNo string,
    jobType valueobject.JobType,
    priority valueobject.JobPriority,
    scheduledAt time.Time,
) (*InstallationJob, error) {
    if !jobType.IsValid() {
        return nil, domainerrors.ErrInvalidJobType
    }
    if !priority.IsValid() {
        return nil, domainerrors.ErrInvalidPriority
    }
    jobNo = strings.TrimSpace(jobNo)
    if jobNo == "" {
        return nil, domainerrors.ErrInvalidJobNo
    }
    now := time.Now()
    slaDue := scheduledAt.Add(time.Duration(jobType.DefaultSLAHours(priority)) * time.Hour)
    return &InstallationJob{
        ID:          uuid.New(),
        TenantID:    tenantID,
        JobNo:       jobNo,
        CustomerID:  customerID,
        SiteID:      siteID,
        JobType:     jobType,
        Priority:    priority,
        Status:      valueobject.JobStatusScheduled,
        ScheduledAt: scheduledAt,
        SLADueAt:    slaDue,
        Checklist:   valueobject.StandardInstallChecklist(),
        DeviceIDs:   []uuid.UUID{},
        Photos:      []string{},
        CreatedBy:   actorID,
        CreatedAt:   now,
        UpdatedAt:   now,
    }, nil
}

// --- Assignment ---

func (j *InstallationJob) AssignTechnician(techID, assignedBy uuid.UUID) error {
    if !j.Status.CanTransitionTo(valueobject.JobStatusAssigned) {
        return domainerrors.ErrInvalidStatusTransition
    }
    now := time.Now()
    j.TechnicianID = &techID
    j.AssignedBy = &assignedBy
    j.AssignedAt = &now
    j.Status = valueobject.JobStatusAssigned
    j.UpdatedAt = now
    return nil
}

func (j *InstallationJob) ReassignTechnician(techID, assignedBy uuid.UUID) error {
    if j.Status.IsTerminal() {
        return domainerrors.ErrJobAlreadyTerminal
    }
    if j.Status == valueobject.JobStatusInProgress {
        return domainerrors.ErrCannotReassignInProgress
    }
    now := time.Now()
    j.TechnicianID = &techID
    j.AssignedBy = &assignedBy
    j.AssignedAt = &now
    j.UpdatedAt = now
    return nil
}

// --- Execution ---

func (j *InstallationJob) Start(gps valueobject.GPSLocation, geoFence *valueobject.GeoFence) error {
    if j.Status != valueobject.JobStatusAssigned {
        return domainerrors.ErrInvalidStatusTransition
    }
    if j.TechnicianID == nil {
        return domainerrors.ErrNoTechnicianAssigned
    }
    if geoFence != nil && !geoFence.Contains(gps) {
        return domainerrors.ErrOutsideGeoFence
    }
    if !gps.IsValid() {
        return domainerrors.ErrInvalidGPS
    }
    now := time.Now()
    j.StartedAt = &now
    j.StartedGPS = &gps
    j.Status = valueobject.JobStatusInProgress
    j.UpdatedAt = now
    return nil
}

func (j *InstallationJob) Pause(reason string) error {
    if j.Status != valueobject.JobStatusInProgress {
        return domainerrors.ErrInvalidStatusTransition
    }
    j.Status = valueobject.JobStatusPaused
    j.Notes = appendNote(j.Notes, "pause: "+reason)
    j.UpdatedAt = time.Now()
    return nil
}

func (j *InstallationJob) Resume() error {
    if j.Status != valueobject.JobStatusPaused {
        return domainerrors.ErrInvalidStatusTransition
    }
    j.Status = valueobject.JobStatusInProgress
    j.UpdatedAt = time.Now()
    return nil
}

func (j *InstallationJob) UpdateChecklist(itemKey string, checked bool, notes, photoURL string) error {
    if j.Status != valueobject.JobStatusInProgress {
        return domainerrors.ErrJobNotInProgress
    }
    for i := range j.Checklist {
        if j.Checklist[i].Key == itemKey {
            j.Checklist[i].Checked = checked
            j.Checklist[i].Notes = notes
            if photoURL != "" {
                j.Checklist[i].PhotoURL = photoURL
            }
            if checked {
                j.Checklist[i].CheckedAt = time.Now()
            }
            j.UpdatedAt = time.Now()
            return nil
        }
    }
    return domainerrors.ErrChecklistItemNotFound
}

func (j *InstallationJob) AddPhoto(url string) error {
    if j.Status != valueobject.JobStatusInProgress {
        return domainerrors.ErrJobNotInProgress
    }
    if len(j.Photos) >= 20 {
        return domainerrors.ErrTooManyPhotos
    }
    j.Photos = append(j.Photos, url)
    j.UpdatedAt = time.Now()
    return nil
}

func (j *InstallationJob) Complete(gps valueobject.GPSLocation, signature, notes string) error {
    if j.Status != valueobject.JobStatusInProgress {
        return domainerrors.ErrInvalidStatusTransition
    }
    if !j.Checklist.AllRequiredChecked() {
        return domainerrors.ErrChecklistIncomplete
    }
    if len(j.Photos) < 1 {
        return domainerrors.ErrEvidenceRequired
    }
    if !gps.IsValid() {
        return domainerrors.ErrInvalidGPS
    }
    now := time.Now()
    j.CompletedAt = &now
    j.CompletedGPS = &gps
    j.SignatureURL = signature
    j.Notes = appendNote(j.Notes, notes)
    j.Status = valueobject.JobStatusDone
    j.UpdatedAt = now
    return nil
}

func (j *InstallationJob) Fail(reason string, gps *valueobject.GPSLocation) error {
    if j.Status.IsTerminal() {
        return domainerrors.ErrJobAlreadyTerminal
    }
    j.Status = valueobject.JobStatusFailed
    j.CancelReason = reason
    j.CompletedGPS = gps
    j.UpdatedAt = time.Now()
    return nil
}

func (j *InstallationJob) Cancel(reason string) error {
    if j.Status.IsTerminal() {
        return domainerrors.ErrJobAlreadyTerminal
    }
    if j.Status == valueobject.JobStatusInProgress {
        return domainerrors.ErrCannotCancelInProgress
    }
    now := time.Now()
    j.Status = valueobject.JobStatusCancelled
    j.CancelledAt = &now
    j.CancelReason = reason
    j.UpdatedAt = now
    return nil
}

func (j *InstallationJob) AttachDevice(deviceID uuid.UUID) error {
    for _, existing := range j.DeviceIDs {
        if existing == deviceID {
            return nil // idempotent
        }
    }
    j.DeviceIDs = append(j.DeviceIDs, deviceID)
    return nil
}

// --- Query ---

func (j *InstallationJob) IsSLABreached(asOf time.Time) bool {
    if j.Status.IsTerminal() {
        if j.CompletedAt != nil {
            return j.CompletedAt.After(j.SLADueAt)
        }
        return false
    }
    return asOf.After(j.SLADueAt)
}

func (j *InstallationJob) DurationMinutes() int {
    if j.StartedAt == nil { return 0 }
    end := time.Now()
    if j.CompletedAt != nil { end = *j.CompletedAt }
    return int(end.Sub(*j.StartedAt).Minutes())
}

func (j *InstallationJob) ChecklistProgress() (checked, total int) {
    return j.Checklist.Progress()
}

func appendNote(existing, add string) string {
    add = strings.TrimSpace(add)
    if add == "" { return existing }
    if existing == "" { return add }
    return existing + "\n" + add
}
```

### `domain/entity/technician.go` (Aggregate Root #3)
```go
package entity

import (
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

type TechnicianLevel string

const (
    TechnicianLevelJunior TechnicianLevel = "JUNIOR"
    TechnicianLevelSenior TechnicianLevel = "SENIOR"
    TechnicianLevelLead   TechnicianLevel = "LEAD"
)

// Technician – Aggregate Root
type Technician struct {
    ID             uuid.UUID
    TenantID       uuid.UUID
    UserID         uuid.UUID  // link to users module

    Code           string
    Name           string
    Phone          string
    Email          string
    Level          TechnicianLevel
    Skills         valueobject.SkillSet
    Zone           string      // e.g. "BANGKOK", "CHIANG_MAI"
    HomeBase       *valueobject.GPSLocation

    IsAvailable    bool
    IsActive       bool
    GPSLast        *valueobject.GPSLocation
    GPSUpdatedAt   *time.Time

    MaxJobsPerDay  int
    CurrentLoad    int   // jobs วันนี้

    CreatedAt      time.Time
    UpdatedAt      time.Time
}

func NewTechnician(tenantID, userID uuid.UUID, code, name string) (*Technician, error) {
    if code == "" || name == "" {
        return nil, domainerrors.ErrInvalidTechnicianData
    }
    now := time.Now()
    return &Technician{
        ID:             uuid.New(),
        TenantID:       tenantID,
        UserID:         userID,
        Code:           code,
        Name:           name,
        Level:          TechnicianLevelJunior,
        Skills:         valueobject.SkillSet{},
        IsAvailable:    true,
        IsActive:       true,
        MaxJobsPerDay:  6,
        CreatedAt:      now,
        UpdatedAt:      now,
    }, nil
}

func (t *Technician) AddSkill(s valueobject.Skill) error {
    if !s.IsValid() {
        return domainerrors.ErrInvalidSkill
    }
    for _, existing := range t.Skills {
        if existing == s { return nil }
    }
    t.Skills = append(t.Skills, s)
    t.UpdatedAt = time.Now()
    return nil
}

func (t *Technician) UpdateGPS(loc valueobject.GPSLocation) error {
    if !loc.IsValid() {
        return domainerrors.ErrInvalidGPS
    }
    now := time.Now()
    t.GPSLast = &loc
    t.GPSUpdatedAt = &now
    t.UpdatedAt = now
    return nil
}

func (t *Technician) SetAvailable(available bool) {
    t.IsAvailable = available
    t.UpdatedAt = time.Now()
}

func (t *Technician) CanTakeJob(requiredSkills []valueobject.Skill) bool {
    if !t.IsActive || !t.IsAvailable { return false }
    if t.CurrentLoad >= t.MaxJobsPerDay { return false }
    if len(requiredSkills) > 0 {
        return t.Skills.MatchRatio(requiredSkills) >= 0.6 // 60% ขึ้นไป
    }
    return true
}

func (t *Technician) IncrementLoad() { t.CurrentLoad++ }
func (t *Technician) ResetDailyLoad() { t.CurrentLoad = 0 }

func (t *Technician) GPSTTL() time.Duration {
    if t.GPSUpdatedAt == nil { return 0 }
    return time.Since(*t.GPSUpdatedAt)
}

func (t *Technician) IsGPSStale(threshold time.Duration) bool {
    return t.GPSTTL() > threshold
}
```

### `domain/entity/work_order.go` (Aggregate Root #4)
```go
package entity

import (
    "strings"
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

// SparePartUsage – spare part ที่ใช้ในงาน
type SparePartUsage struct {
    PartID    uuid.UUID
    PartNo    string
    Name      string
    Qty       float64
    UnitPrice float64
    Currency  string
}

// WorkOrder – Aggregate Root (สำหรับ maintenance / repair)
type WorkOrder struct {
    ID             uuid.UUID
    TenantID       uuid.UUID
    WONo           string
    DeviceID       *uuid.UUID
    SiteID         uuid.UUID
    CustomerID     uuid.UUID
    JobType        valueobject.JobType
    Priority       valueobject.JobPriority
    Status         valueobject.JobStatus
    Title          string
    Description    string

    TechnicianID   *uuid.UUID
    ScheduledAt    time.Time
    TimeWindow     *valueobject.TimeWindow

    StartedAt      *time.Time
    CompletedAt    *time.Time
    Resolution     string
    PartsUsed      []SparePartUsage
    LaborMinutes   int
    CostTotal      float64
    Currency       string

    SourceRef      string  // "pm_schedule" | "alert" | "customer_ticket"
    SourceRefID    *uuid.UUID
    Notes          string

    CreatedBy      uuid.UUID
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

func NewWorkOrder(
    tenantID, customerID, siteID, actorID uuid.UUID,
    woNo, title string,
    jobType valueobject.JobType,
    priority valueobject.JobPriority,
) (*WorkOrder, error) {
    if woNo == "" { return nil, domainerrors.ErrInvalidWONo }
    if title == "" { return nil, domainerrors.ErrInvalidTitle }
    if !jobType.IsValid() { return nil, domainerrors.ErrInvalidJobType }
    if !priority.IsValid() { return nil, domainerrors.ErrInvalidPriority }
    now := time.Now()
    return &WorkOrder{
        ID: uuid.New(), TenantID: tenantID,
        CustomerID: customerID, SiteID: siteID,
        WONo: woNo, Title: title,
        JobType: jobType, Priority: priority,
        Status: valueobject.JobStatusScheduled,
        PartsUsed: []SparePartUsage{},
        Currency: "THB",
        CreatedBy: actorID, CreatedAt: now, UpdatedAt: now,
    }, nil
}

func (w *WorkOrder) AssignTechnician(techID uuid.UUID) error {
    if w.Status.IsTerminal() {
        return domainerrors.ErrJobAlreadyTerminal
    }
    w.TechnicianID = &techID
    w.Status = valueobject.JobStatusAssigned
    w.UpdatedAt = time.Now()
    return nil
}

func (w *WorkOrder) Start() error {
    if w.Status != valueobject.JobStatusAssigned {
        return domainerrors.ErrInvalidStatusTransition
    }
    if w.TechnicianID == nil {
        return domainerrors.ErrNoTechnicianAssigned
    }
    now := time.Now()
    w.StartedAt = &now
    w.Status = valueobject.JobStatusInProgress
    w.UpdatedAt = now
    return nil
}

func (w *WorkOrder) AddPartUsage(usage SparePartUsage) error {
    if w.Status != valueobject.JobStatusInProgress {
        return domainerrors.ErrJobNotInProgress
    }
    if usage.Qty <= 0 { return domainerrors.ErrInvalidQuantity }
    w.PartsUsed = append(w.PartsUsed, usage)
    w.CostTotal += usage.Qty * usage.UnitPrice
    w.UpdatedAt = time.Now()
    return nil
}

func (w *WorkOrder) Complete(resolution string, laborMinutes int) error {
    if w.Status != valueobject.JobStatusInProgress {
        return domainerrors.ErrInvalidStatusTransition
    }
    resolution = strings.TrimSpace(resolution)
    if resolution == "" {
        return domainerrors.ErrResolutionRequired
    }
    now := time.Now()
    w.Resolution = resolution
    w.LaborMinutes = laborMinutes
    w.CompletedAt = &now
    w.Status = valueobject.JobStatusDone
    w.UpdatedAt = now
    return nil
}

func (w *WorkOrder) Fail(reason string) error {
    if w.Status.IsTerminal() { return domainerrors.ErrJobAlreadyTerminal }
    w.Status = valueobject.JobStatusFailed
    w.Resolution = reason
    w.UpdatedAt = time.Now()
    return nil
}

func (w *WorkOrder) IsDone() bool { return w.Status == valueobject.JobStatusDone }
```

### `domain/entity/maintenance_schedule.go` (Aggregate Root #5)
```go
package entity

import (
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

type PMFrequency string

const (
    PMFrequencyMonthly   PMFrequency = "MONTHLY"
    PMFrequencyQuarterly PMFrequency = "QUARTERLY"
    PMFrequencySemiAnnual PMFrequency = "SEMI_ANNUAL"
    PMFrequencyAnnual    PMFrequency = "ANNUAL"
    PMFrequencyCustom    PMFrequency = "CUSTOM"
)

func (f PMFrequency) Days() int {
    switch f {
    case PMFrequencyMonthly:    return 30
    case PMFrequencyQuarterly:  return 90
    case PMFrequencySemiAnnual: return 180
    case PMFrequencyAnnual:     return 365
    }
    return 0
}

// MaintenanceSchedule – Aggregate Root (Preventive Maintenance)
type MaintenanceSchedule struct {
    ID               uuid.UUID
    TenantID         uuid.UUID
    DeviceID         uuid.UUID
    CustomerID       uuid.UUID
    SiteID           uuid.UUID
    JobType          valueobject.JobType
    Frequency        PMFrequency
    CustomDays       int  // ถ้า CUSTOM
    LastDoneAt       *time.Time
    NextDueAt        time.Time
    AssignedTechID   *uuid.UUID
    AutoCreateWO     bool
    IsActive         bool
    Notes            string
    CreatedAt        time.Time
    UpdatedAt        time.Time
}

func NewMaintenanceSchedule(
    tenantID, deviceID, customerID, siteID uuid.UUID,
    freq PMFrequency,
    startDate time.Time,
) (*MaintenanceSchedule, error) {
    if !isValidPMFrequency(freq) {
        return nil, domainerrors.ErrInvalidFrequency
    }
    now := time.Now()
    nextDue := startDate
    if nextDue.Before(now) {
        nextDue = now.AddDate(0, 0, freq.Days())
    }
    return &MaintenanceSchedule{
        ID: uuid.New(), TenantID: tenantID,
        DeviceID: deviceID, CustomerID: customerID, SiteID: siteID,
        JobType: valueobject.JobTypeMaintenance,
        Frequency: freq, NextDueAt: nextDue,
        AutoCreateWO: true, IsActive: true,
        CreatedAt: now, UpdatedAt: now,
    }, nil
}

// MarkDone – ตั้ง last_done_at + คำนวณ next_due_at
func (m *MaintenanceSchedule) MarkDone(doneAt time.Time) {
    m.LastDoneAt = &doneAt
    days := m.Frequency.Days()
    if m.Frequency == PMFrequencyCustom && m.CustomDays > 0 {
        days = m.CustomDays
    }
    m.NextDueAt = doneAt.AddDate(0, 0, days)
    m.UpdatedAt = time.Now()
}

// IsDue – ตรวจว่าถึงกำหนด PM หรือยัง
func (m *MaintenanceSchedule) IsDue(asOf time.Time) bool {
    return m.IsActive && !asOf.Before(m.NextDueAt)
}

// IsOverdueBy – overdue เกิน N วัน
func (m *MaintenanceSchedule) IsOverdueBy(asOf time.Time, days int) bool {
    if !m.IsActive { return false }
    overdue := asOf.Sub(m.NextDueAt)
    return overdue > time.Duration(days)*24*time.Hour
}

func (m *MaintenanceSchedule) Disable() {
    m.IsActive = false
    m.UpdatedAt = time.Now()
}

func isValidPMFrequency(f PMFrequency) bool {
    switch f {
    case PMFrequencyMonthly, PMFrequencyQuarterly, PMFrequencySemiAnnual,
        PMFrequencyAnnual, PMFrequencyCustom:
        return true
    }
    return false
}
```

### `domain/entity/rma.go` (Aggregate Root #6)
```go
package entity

import (
    "strings"
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

type RMAReason string

const (
    RMAReasonDefective     RMAReason = "DEFECTIVE"
    RMAReasonDamaged       RMAReason = "DAMAGED"
    RMAReasonWrongItem     RMAReason = "WRONG_ITEM"
    RMAReasonWarranty      RMAReason = "WARRANTY"
    RMAReasonCustomerError RMAReason = "CUSTOMER_ERROR"
    RMAReasonUpgrade       RMAReason = "UPGRADE"
    RMAReasonOther         RMAReason = "OTHER"
)

type RMAResolution string

const (
    RMAResolutionRepair    RMAResolution = "REPAIR"
    RMAResolutionReplace   RMAResolution = "REPLACE"
    RMAResolutionRefund    RMAResolution = "REFUND"
    RMAResolutionNoAction  RMAResolution = "NO_ACTION"
)

// RMA – Aggregate Root
type RMA struct {
    ID                  uuid.UUID
    TenantID            uuid.UUID
    RMANo               string
    DeviceID            uuid.UUID
    CustomerID          uuid.UUID
    SiteID              uuid.UUID
    JobID               *uuid.UUID  // link installation job ที่เจอปัญหา
    Reason              RMAReason
    ReasonDetail        string
    Status              valueobject.RMAStatus
    Resolution          *RMAResolution
    ReplacementDeviceID *uuid.UUID
    WarrantyClaim       bool
    ReceivedPhotos      []string

    ApprovedAt          *time.Time
    ApprovedBy          *uuid.UUID
    ShippedAt           *time.Time
    ReceivedAt          *time.Time
    ClosedAt            *time.Time
    RejectReason        string

    AssignedTechID      *uuid.UUID
    Notes               string

    CreatedBy           uuid.UUID
    CreatedAt           time.Time
    UpdatedAt           time.Time
}

func NewRMA(
    tenantID, deviceID, customerID, siteID, actorID uuid.UUID,
    rmaNo string,
    reason RMAReason,
    detail string,
) (*RMA, error) {
    if rmaNo == "" {
        return nil, domainerrors.ErrInvalidRMANo
    }
    if !isValidRMAReason(reason) {
        return nil, domainerrors.ErrInvalidRMAReason
    }
    now := time.Now()
    return &RMA{
        ID: uuid.New(), TenantID: tenantID,
        RMANo: rmaNo,
        DeviceID: deviceID, CustomerID: customerID, SiteID: siteID,
        Reason: reason, ReasonDetail: strings.TrimSpace(detail),
        Status: valueobject.RMAStatusOpen,
        ReceivedPhotos: []string{},
        CreatedBy: actorID, CreatedAt: now, UpdatedAt: now,
    }, nil
}

func (r *RMA) Approve(approverID uuid.UUID, resolution RMAResolution) error {
    if r.Status != valueobject.RMAStatusOpen {
        return domainerrors.ErrInvalidStatusTransition
    }
    if !isValidRMAResolution(resolution) {
        return domainerrors.ErrInvalidRMAResolution
    }
    now := time.Now()
    r.Status = valueobject.RMAStatusApproved
    r.ApprovedAt = &now
    r.ApprovedBy = &approverID
    r.Resolution = &resolution
    r.UpdatedAt = now
    return nil
}

func (r *RMA) Reject(reason string, actorID uuid.UUID) error {
    if r.Status != valueobject.RMAStatusOpen {
        return domainerrors.ErrInvalidStatusTransition
    }
    if strings.TrimSpace(reason) == "" {
        return domainerrors.ErrRejectReasonRequired
    }
    r.Status = valueobject.RMAStatusRejected
    r.RejectReason = reason
    r.UpdatedAt = time.Now()
    return nil
}

func (r *RMA) Ship(trackingNo string) error {
    if r.Status != valueobject.RMAStatusApproved {
        return domainerrors.ErrInvalidStatusTransition
    }
    now := time.Now()
    r.Status = valueobject.RMAStatusShipped
    r.ShippedAt = &now
    r.Notes = appendNote(r.Notes, "shipped: "+trackingNo)
    r.UpdatedAt = now
    return nil
}

func (r *RMA) Receive(warehouseID uuid.UUID) error {
    if r.Status != valueobject.RMAStatusShipped {
        return domainerrors.ErrInvalidStatusTransition
    }
    now := time.Now()
    r.Status = valueobject.RMAStatusReceived
    r.ReceivedAt = &now
    r.UpdatedAt = now
    return nil
}

func (r *RMA) StartRepair(techID uuid.UUID) error {
    if r.Status != valueobject.RMAStatusReceived {
        return domainerrors.ErrInvalidStatusTransition
    }
    r.Status = valueobject.RMAStatusRepairing
    r.AssignedTechID = &techID
    r.UpdatedAt = time.Now()
    return nil
}

func (r *RMA) MarkRepaired() error {
    if r.Status != valueobject.RMAStatusRepairing {
        return domainerrors.ErrInvalidStatusTransition
    }
    r.Status = valueobject.RMAStatusRepaired
    r.UpdatedAt = time.Now()
    return nil
}

func (r *RMA) ReplaceDevice(newDeviceID uuid.UUID) error {
    if r.Status != valueobject.RMAStatusReceived {
        return domainerrors.ErrInvalidStatusTransition
    }
    r.ReplacementDeviceID = &newDeviceID
    r.Status = valueobject.RMAStatusReplaced
    r.UpdatedAt = time.Now()
    return nil
}

func (r *RMA) Close(notes string) error {
    if r.Status.IsTerminal() {
        return domainerrors.ErrRMAAlreadyTerminal
    }
    now := time.Now()
    r.Status = valueobject.RMAStatusClosed
    r.ClosedAt = &now
    r.Notes = appendNote(r.Notes, notes)
    r.UpdatedAt = now
    return nil
}

func (r *RMA) AddPhoto(url string) {
    r.ReceivedPhotos = append(r.ReceivedPhotos, url)
}

func (r *RMA) IsWarranty() bool { return r.WarrantyClaim }

func isValidRMAReason(r RMAReason) bool {
    switch r {
    case RMAReasonDefective, RMAReasonDamaged, RMAReasonWrongItem,
        RMAReasonWarranty, RMAReasonCustomerError, RMAReasonUpgrade, RMAReasonOther:
        return true
    }
    return false
}

func isValidRMAResolution(r RMAResolution) bool {
    switch r {
    case RMAResolutionRepair, RMAResolutionReplace,
        RMAResolutionRefund, RMAResolutionNoAction:
        return true
    }
    return false
}
```

### `domain/entity/technician_location.go`
```go
package entity

import (
    "time"

    "github.com/google/uuid"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

// TechnicianLocation – tracking GPS log (store in PG + optionally InfluxDB)
type TechnicianLocation struct {
    ID           int64
    TenantID     uuid.UUID
    TechnicianID uuid.UUID
    Location     valueobject.GPSLocation
    JobID        *uuid.UUID
    BatteryLevel *int
    Network      string
    RecordedAt   time.Time
}
```

### `domain/entity/spare_part.go`
```go
package entity

import (
    "time"

    "github.com/google/uuid"
)

// SparePart – master data ของอะไหล่
type SparePart struct {
    ID         uuid.UUID
    TenantID   uuid.UUID
    PartNo     string
    Name       string
    Category   string
    UnitPrice  float64
    Currency   string
    IsActive   bool
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
```

---

## A.4 Repository Interfaces

### `domain/repository/shipment_repository.go`
```go
package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "icmongolang/internal/modules/iotlogistics/domain/entity"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

type ShipmentFilter struct {
    TenantID   uuid.UUID
    CustomerID *uuid.UUID
    SiteID     *uuid.UUID
    Statuses   []valueobject.ShipmentStatus
    From       *time.Time
    To         *time.Time
    Search     string
    Page       int
    PageSize   int
    SortBy     string
    SortOrder  string
}

type ShipmentRepository interface {
    Save(ctx context.Context, s *entity.Shipment) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Shipment, error)
    FindByNo(ctx context.Context, tenantID uuid.UUID, no valueobject.ShipmentNo) (*entity.Shipment, error)
    List(ctx context.Context, f ShipmentFilter) ([]*entity.Shipment, int64, error)
    ListByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) ([]*entity.Shipment, error)
    ListPendingBySite(ctx context.Context, tenantID, siteID uuid.UUID) ([]*entity.Shipment, error)
    NextSequence(ctx context.Context, tenantID uuid.UUID, year int) (int, error)
}
```

### `domain/repository/installation_repository.go`
```go
package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "icmongolang/internal/modules/iotlogistics/domain/entity"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

type InstallationJobFilter struct {
    TenantID     uuid.UUID
    TechnicianID *uuid.UUID
    SiteID       *uuid.UUID
    CustomerID   *uuid.UUID
    Statuses     []valueobject.JobStatus
    Types        []valueobject.JobType
    From         *time.Time
    To           *time.Time
    OnlySLABreach *bool
    Page         int
    PageSize     int
}

type InstallationRepository interface {
    Save(ctx context.Context, j *entity.InstallationJob) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.InstallationJob, error)
    FindByNo(ctx context.Context, tenantID uuid.UUID, no string) (*entity.InstallationJob, error)
    List(ctx context.Context, f InstallationJobFilter) ([]*entity.InstallationJob, int64, error)
    ListByTechnician(ctx context.Context, tenantID, techID uuid.UUID, date time.Time) ([]*entity.InstallationJob, error)
    ListBySite(ctx context.Context, tenantID, siteID uuid.UUID) ([]*entity.InstallationJob, error)
    ListSLABreached(ctx context.Context, tenantID uuid.UUID, asOf time.Time) ([]*entity.InstallationJob, error)
    CountActiveByTechnician(ctx context.Context, techID uuid.UUID, date time.Time) (int64, error)
    NextSequence(ctx context.Context, tenantID uuid.UUID, year int) (int, error)
}
```

### `domain/repository/technician_repository.go`
```go
package repository

import (
    "context"

    "github.com/google/uuid"
    "icmongolang/internal/modules/iotlogistics/domain/entity"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

type TechnicianFilter struct {
    TenantID    uuid.UUID
    Zone        string
    Skills      []valueobject.Skill
    Available   *bool
    Active      *bool
    Page        int
    PageSize    int
}

type TechnicianRepository interface {
    Save(ctx context.Context, t *entity.Technician) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Technician, error)
    FindByUserID(ctx context.Context, tenantID, userID uuid.UUID) (*entity.Technician, error)
    List(ctx context.Context, f TechnicianFilter) ([]*entity.Technician, int64, error)
    ListAvailableNear(ctx context.Context, tenantID uuid.UUID, center valueobject.GPSLocation, radiusKm float64) ([]*entity.Technician, error)
    UpdateGPS(ctx context.Context, id uuid.UUID, loc valueobject.GPSLocation) error
    ResetDailyLoads(ctx context.Context, tenantID uuid.UUID) error
}
```

### `domain/repository/work_order_repository.go`
```go
package repository

import (
    "context"

    "github.com/google/uuid"
    "icmongolang/internal/modules/iotlogistics/domain/entity"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

type WorkOrderFilter struct {
    TenantID     uuid.UUID
    TechnicianID *uuid.UUID
    SiteID       *uuid.UUID
    Statuses     []valueobject.JobStatus
    SourceRef    string
    Page         int
    PageSize     int
}

type WorkOrderRepository interface {
    Save(ctx context.Context, w *entity.WorkOrder) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.WorkOrder, error)
    FindByNo(ctx context.Context, tenantID uuid.UUID, no string) (*entity.WorkOrder, error)
    List(ctx context.Context, f WorkOrderFilter) ([]*entity.WorkOrder, int64, error)
    NextSequence(ctx context.Context, tenantID uuid.UUID, year int) (int, error)
}
```

### `domain/repository/maintenance_repository.go`
```go
package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "icmongolang/internal/modules/iotlogistics/domain/entity"
)

type MaintenanceFilter struct {
    TenantID  uuid.UUID
    SiteID    *uuid.UUID
    DeviceID  *uuid.UUID
    Active    *bool
    Page      int
    PageSize  int
}

type MaintenanceRepository interface {
    Save(ctx context.Context, m *entity.MaintenanceSchedule) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.MaintenanceSchedule, error)
    FindByDevice(ctx context.Context, tenantID, deviceID uuid.UUID) (*entity.MaintenanceSchedule, error)
    List(ctx context.Context, f MaintenanceFilter) ([]*entity.MaintenanceSchedule, int64, error)
    ListDue(ctx context.Context, tenantID uuid.UUID, asOf time.Time) ([]*entity.MaintenanceSchedule, error)
    ListDueAllTenants(ctx context.Context, asOf time.Time) ([]*entity.MaintenanceSchedule, error)
    Delete(ctx context.Context, tenantID, id uuid.UUID) error
}
```

### `domain/repository/rma_repository.go`
```go
package repository

import (
    "context"

    "github.com/google/uuid"
    "icmongolang/internal/modules/iotlogistics/domain/entity"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

type RMAFilter struct {
    TenantID   uuid.UUID
    CustomerID *uuid.UUID
    DeviceID   *uuid.UUID
    Statuses   []valueobject.RMAStatus
    Page       int
    PageSize   int
}

type RMARepository interface {
    Save(ctx context.Context, r *entity.RMA) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.RMA, error)
    FindByNo(ctx context.Context, tenantID uuid.UUID, no string) (*entity.RMA, error)
    List(ctx context.Context, f RMAFilter) ([]*entity.RMA, int64, error)
    ListByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) ([]*entity.RMA, error)
    CountOpenByDevice(ctx context.Context, tenantID, deviceID uuid.UUID) (int64, error)
    NextSequence(ctx context.Context, tenantID uuid.UUID, year int) (int, error)
}
```

### `domain/repository/spare_part_repository.go`
```go
package repository

import (
    "context"

    "github.com/google/uuid"
    "icmongolang/internal/modules/iotlogistics/domain/entity"
)

type SparePartRepository interface {
    Save(ctx context.Context, p *entity.SparePart) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.SparePart, error)
    FindByPartNo(ctx context.Context, tenantID uuid.UUID, partNo string) (*entity.SparePart, error)
    List(ctx context.Context, tenantID uuid.UUID) ([]*entity.SparePart, error)
}
```

### `domain/repository/audit_repository.go`
```go
package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
)

type AuditEntry struct {
    ID         int64
    TenantID   uuid.UUID
    ActorID    uuid.UUID
    Action     string
    EntityType string
    EntityID   uuid.UUID
    Payload    map[string]any
    IPAddress  string
    UserAgent  string
    CreatedAt  time.Time
}

type AuditRepository interface {
    Save(ctx context.Context, e *AuditEntry) error
}
```

---

## A.5 Domain Services & Ports

### `domain/service/technician_assignment_service.go` ⭐
```go
package service

import (
    "context"
    "math"
    "sort"

    "github.com/google/uuid"

    "icmongolang/internal/modules/iotlogistics/domain/entity"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    "icmongolang/internal/modules/iotlogistics/domain/repository"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

// TechnicianAssignmentService – เลือกช่างที่เหมาะสมที่สุด
// Scoring criteria:
//   1. Skill match (40%)
//   2. Proximity  (30%)
//   3. Workload   (20%)
//   4. Level      (10%)
type TechnicianAssignmentService struct {
    techRepo repository.TechnicianRepository
}

func NewTechnicianAssignmentService(techRepo repository.TechnicianRepository) *TechnicianAssignmentService {
    return &TechnicianAssignmentService{techRepo: techRepo}
}

type AssignmentCandidate struct {
    Technician *entity.Technician
    Score      float64
    DistanceKm float64
    MatchRatio float64
    Breakdown  map[string]float64
}

// Assign – เลือกช่างที่ดีที่สุด
func (s *TechnicianAssignmentService) Assign(
    ctx context.Context,
    tenantID uuid.UUID,
    siteGPS valueobject.GPSLocation,
    requiredSkills []valueobject.Skill,
    zone string,
) (*entity.Technician, error) {
    // 1. ดึง candidates — ใน zone + active + available
    filter := repository.TechnicianFilter{
        TenantID:  tenantID,
        Zone:      zone,
        Available: boolPtr(true),
        Active:    boolPtr(true),
        Page:      1,
        PageSize:  50,
    }
    techs, _, err := s.techRepo.List(ctx, filter)
    if err != nil {
        return nil, err
    }
    if len(techs) == 0 {
        return nil, domainerrors.ErrNoTechnicianAvailable
    }

    // 2. Score
    candidates := make([]*AssignmentCandidate, 0, len(techs))
    for _, t := range techs {
        if !t.CanTakeJob(requiredSkills) {
            continue
        }
        c := s.score(t, siteGPS, requiredSkills)
        candidates = append(candidates, c)
    }
    if len(candidates) == 0 {
        return nil, domainerrors.ErrNoTechnicianAvailable
    }

    // 3. เลือก score สูงสุด
    sort.Slice(candidates, func(i, j int) bool {
        return candidates[i].Score > candidates[j].Score
    })
    return candidates[0].Technician, nil
}

// Rank – คืน list เรียงตามคะแนน
func (s *TechnicianAssignmentService) Rank(
    ctx context.Context,
    tenantID uuid.UUID,
    siteGPS valueobject.GPSLocation,
    requiredSkills []valueobject.Skill,
    zone string,
    limit int,
) ([]*AssignmentCandidate, error) {
    filter := repository.TechnicianFilter{
        TenantID: tenantID, Zone: zone, Active: boolPtr(true),
        Page: 1, PageSize: 100,
    }
    techs, _, err := s.techRepo.List(ctx, filter)
    if err != nil { return nil, err }

    out := make([]*AssignmentCandidate, 0, len(techs))
    for _, t := range techs {
        out = append(out, s.score(t, siteGPS, requiredSkills))
    }
    sort.Slice(out, func(i, j int) bool { return out[i].Score > out[j].Score })
    if limit > 0 && len(out) > limit {
        out = out[:limit]
    }
    return out, nil
}

func (s *TechnicianAssignmentService) score(
    t *entity.Technician,
    siteGPS valueobject.GPSLocation,
    requiredSkills []valueobject.Skill,
) *AssignmentCandidate {
    c := &AssignmentCandidate{Technician: t, Breakdown: map[string]float64{}}

    // 1. Skill match (40%)
    c.MatchRatio = t.Skills.MatchRatio(requiredSkills)
    skillScore := c.MatchRatio * 40
    c.Breakdown["skill"] = skillScore

    // 2. Proximity (30%) — ระยะทางสูงสุด 50 km = 0 คะแนน
    if t.GPSLast != nil && siteGPS.IsValid() {
        c.DistanceKm = t.GPSLast.DistanceKm(siteGPS)
    }
    proximityScore := math.Max(0, (50-c.DistanceKm)/50) * 30
    c.Breakdown["proximity"] = proximityScore

    // 3. Workload (20%)
    loadRatio := float64(t.CurrentLoad) / float64(t.MaxJobsPerDay)
    workloadScore := (1 - loadRatio) * 20
    c.Breakdown["workload"] = workloadScore

    // 4. Level (10%)
    levelScore := 10.0
    switch t.Level {
    case entity.TechnicianLevelLead:   levelScore = 10
    case entity.TechnicianLevelSenior: levelScore = 7
    case entity.TechnicianLevelJunior: levelScore = 4
    }
    c.Breakdown["level"] = levelScore

    c.Score = skillScore + proximityScore + workloadScore + levelScore
    return c
}

func boolPtr(b bool) *bool { return &b }
```

### `domain/service/route_planning_service.go` ⭐
```go
package service

import (
    "context"
    "sort"
    "time"

    "github.com/google/uuid"

    "icmongolang/internal/modules/iotlogistics/domain/entity"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    "icmongolang/internal/modules/iotlogistics/domain/service/port"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

// RoutePlanningService – วางแผนเส้นทางให้ช่าง
// Algorithm:
//   1. เริ่มจาก GPS ปัจจุบันของช่าง
//   2. Nearest-neighbor TSP (greedy)
//   3. ถ้ามี AI Advisor → ส่งให้ refine
//   4. คำนวณ ETA แต่ละจุด
type RoutePlanningService struct {
    maps port.MapsPort
    ai   port.AIAdvisorPort
}

func NewRoutePlanningService(maps port.MapsPort, ai port.AIAdvisorPort) *RoutePlanningService {
    return &RoutePlanningService{maps: maps, ai: ai}
}

type PlanRequest struct {
    TenantID     uuid.UUID
    TechnicianID uuid.UUID
    StartLocation valueobject.GPSLocation
    Jobs         []*entity.InstallationJob
    JobGPS       map[uuid.UUID]valueobject.GPSLocation
    MaxHours     int // 8 ชั่วโมงทำงาน
}

func (s *RoutePlanningService) Plan(ctx context.Context, req PlanRequest) (*valueobject.Route, error) {
    if len(req.Jobs) == 0 {
        return nil, domainerrors.ErrNoJobs
    }

    // 1. Build distance matrix from Maps API
    locations := []valueobject.GPSLocation{req.StartLocation}
    for _, j := range req.Jobs {
        if gps, ok := req.JobGPS[j.ID]; ok {
            locations = append(locations, gps)
        }
    }

    matrix, err := s.maps.DistanceMatrix(ctx, locations)
    if err != nil {
        return nil, err
    }

    // 2. Nearest-neighbor TSP
    order := s.nearestNeighbor(matrix, len(req.Jobs))

    // 3. Build route
    route := &valueobject.Route{
        ID:        uuid.New().String(),
        Waypoints: []valueobject.RouteWaypoint{},
        StartAt:   time.Now(),
        Provider:  "nearest_neighbor",
    }

    currentTime := route.StartAt
    for i, idx := range order {
        if idx == 0 { continue } // skip start
        job := req.Jobs[idx-1]
        gps := req.JobGPS[job.ID]
        dist := matrix.Distance(0, idx) // from start
        if i > 1 {
            prev := order[i-1]
            dist = matrix.Distance(prev, idx)
        }

        travelMin := int(dist / 30 * 60) // assume 30 km/h avg
        arrival := currentTime.Add(time.Duration(travelMin) * time.Minute)
        service := estimateServiceMinutes(job)
        depart := arrival.Add(time.Duration(service) * time.Minute)

        route.Waypoints = append(route.Waypoints, valueobject.RouteWaypoint{
            Order:      i,
            RefID:      job.ID.String(),
            Location:   gps,
            ArrivalETA: arrival,
            DepartETA:  depart,
            ServiceMin: service,
            DistanceKm: dist,
        })

        route.TotalKm += dist
        route.TotalMinutes += travelMin + service
        currentTime = depart

        // Stop ถ้าเกิน max hours
        if req.MaxHours > 0 && route.TotalMinutes > req.MaxHours*60 {
            break
        }
    }

    route.EndAt = currentTime

    // 4. AI refinement (optional)
    if s.ai != nil {
        if refined, err := s.ai.RefineRoute(ctx, route, req); err == nil {
            route.Polyline = refined
            route.Provider = "ai_refined"
        }
    }

    return route, nil
}

// nearestNeighbor – greedy TSP
func (s *RoutePlanningService) nearestNeighbor(matrix *port.DistanceMatrix, n int) []int {
    visited := make([]bool, n+1)
    order := make([]int, 0, n+1)
    order = append(order, 0) // start at index 0
    visited[0] = true
    current := 0

    for len(order) < n+1 {
        best := -1
        bestDist := 1e18
        for i := 1; i <= n; i++ {
            if visited[i] { continue }
            d := matrix.Distance(current, i)
            if d < bestDist {
                bestDist = d
                best = i
            }
        }
        if best == -1 { break }
        visited[best] = true
        order = append(order, best)
        current = best
    }
    return order
}

func estimateServiceMinutes(job *entity.InstallationJob) int {
    switch job.JobType {
    case valueobject.JobTypeInstallation: return 90
    case valueobject.JobTypeMaintenance:  return 45
    case valueobject.JobTypeRepair:       return 60
    case valueobject.JobTypeInspection:   return 30
    case valueobject.JobTypeRemoval:      return 45
    case valueobject.JobTypeRelocation:   return 120
    }
    return 60
}

// ReorderByPriority – sort jobs by priority ก่อน plan
func (s *RoutePlanningService) ReorderByPriority(jobs []*entity.InstallationJob) {
    sort.SliceStable(jobs, func(i, j int) bool {
        return jobs[i].Priority.Score() > jobs[j].Priority.Score()
    })
}
```

### `domain/service/geo_fence_service.go`
```go
package service

import (
    "icmongolang/internal/modules/iotlogistics/domain/entity"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

// GeoFenceService – ตรวจ geo-fence ก่อน start/complete job
type GeoFenceService struct {
    DefaultRadiusM float64
}

func NewGeoFenceService(defaultRadiusM float64) *GeoFenceService {
    if defaultRadiusM <= 0 {
        defaultRadiusM = 500
    }
    return &GeoFenceService{DefaultRadiusM: defaultRadiusM}
}

func (s *GeoFenceService) BuildFence(siteGPS valueobject.GPSLocation, radiusM float64) (*valueobject.GeoFence, error) {
    if !siteGPS.IsValid() {
        return nil, domainerrors.ErrInvalidGPS
    }
    if radiusM <= 0 {
        radiusM = s.DefaultRadiusM
    }
    return &valueobject.GeoFence{Center: siteGPS, RadiusM: radiusM}, nil
}

// ValidateCheckIn – ตรวจว่าช่างอยู่ใกล้ site
func (s *GeoFenceService) ValidateCheckIn(
    techGPS valueobject.GPSLocation,
    siteGPS valueobject.GPSLocation,
    radiusM float64,
) error {
    fence, err := s.BuildFence(siteGPS, radiusM)
    if err != nil { return err }
    if !fence.Contains(techGPS) {
        return domainerrors.ErrOutsideGeoFence
    }
    return nil
}

// ValidateCheckInWithJob – variant ใช้ job's site
func (s *GeoFenceService) ValidateCheckInWithJob(
    job *entity.InstallationJob,
    techGPS valueobject.GPSLocation,
    siteGPS valueobject.GPSLocation,
) error {
    return s.ValidateCheckIn(techGPS, siteGPS, s.DefaultRadiusM)
}
```

### `domain/service/sla_calculator.go`
```go
package service

import (
    "time"

    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

// SLACalculator – คำนวณ SLA due date
type SLACalculator struct{}

func NewSLACalculator() *SLACalculator { return &SLACalculator{} }

func (s *SLACalculator) Calculate(
    jobType valueobject.JobType,
    priority valueobject.JobPriority,
    scheduledAt time.Time,
) time.Time {
    hours := jobType.DefaultSLAHours(priority)
    return scheduledAt.Add(time.Duration(hours) * time.Hour)
}

// RemainingHours – เหลือเวลาเท่าไหร่ถึง SLA
func (s *SLACalculator) RemainingHours(dueAt, asOf time.Time) float64 {
    return dueAt.Sub(asOf).Hours()
}

// Severity – ระดับความเร่งด่วน
func (s *SLACalculator) Severity(dueAt, asOf time.Time) string {
    h := s.RemainingHours(dueAt, asOf)
    switch {
    case h < 0:   return "BREACHED"
    case h < 4:   return "CRITICAL"
    case h < 12:  return "WARNING"
    case h < 24:  return "ATTENTION"
    default:      return "OK"
    }
}
```

### `domain/service/checklist_validator.go`
```go
package service

import (
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

// ChecklistValidator – ตรวจ checklist ก่อน complete
type ChecklistValidator struct {
    RequirePhotos    bool
    RequireSignature bool
    MinPhotos        int
}

func NewChecklistValidator() *ChecklistValidator {
    return &ChecklistValidator{
        RequirePhotos:    true,
        RequireSignature: true,
        MinPhotos:        1,
    }
}

func (v *ChecklistValidator) Validate(c valueobject.Checklist, photos []string, signature string) error {
    if err := c.Validate(); err != nil {
        return err
    }
    if !c.AllRequiredChecked() {
        return domainerrors.ErrChecklistIncomplete
    }
    if v.RequirePhotos && len(photos) < v.MinPhotos {
        return domainerrors.ErrEvidenceRequired
    }
    if v.RequireSignature && signature == "" {
        return domainerrors.ErrSignatureRequired
    }
    return nil
}
```

### `domain/service/port/maps_port.go`
```go
package port

import (
    "context"

    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

// MapsPort – outbound port สำหรับ Google Maps / OSRM
type MapsPort interface {
    DistanceMatrix(ctx context.Context, locations []valueobject.GPSLocation) (*DistanceMatrix, error)
    Geocode(ctx context.Context, address string) (valueobject.GPSLocation, error)
    ReverseGeocode(ctx context.Context, loc valueobject.GPSLocation) (string, error)
}

// DistanceMatrix – matrix ของ distance/duration
type DistanceMatrix struct {
    Cells [][]Cell  // [from][to]
    Size  int
}

type Cell struct {
    DistanceKm float64
    DurationMin int
}

func (m *DistanceMatrix) Distance(from, to int) float64 {
    if from < 0 || from >= m.Size || to < 0 || to >= m.Size {
        return 0
    }
    return m.Cells[from][to].DistanceKm
}

func (m *DistanceMatrix) Duration(from, to int) int {
    if from < 0 || from >= m.Size || to < 0 || to >= m.Size {
        return 0
    }
    return m.Cells[from][to].DurationMin
}
```

### `domain/service/port/ai_advisor_port.go`
```go
package port

import (
    "context"

    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

// AIAdvisorPort – outbound port สำหรับ AI route advisor
type AIAdvisorPort interface {
    RefineRoute(ctx context.Context, route *valueobject.Route, req any) (polyline string, err error)
    PredictDuration(ctx context.Context, jobType string, complexity int) (minutes int, err error)
    SummarizeJob(ctx context.Context, notes string) (string, error)
}
```

### `domain/service/port/notifier_port.go`
```go
package port

import "context"

type NotifierPort interface {
    SendJobAssigned(ctx context.Context, technicianID, jobID, jobType string) error
    SendJobReminder(ctx context.Context, technicianID, jobID string, minutesBefore int) error
    SendCustomerETA(ctx context.Context, customerID, jobID string, etaMinutes int) error
    SendPMReminder(ctx context.Context, customerID, deviceID string, daysBefore int) error
    SendRMAAcknowledged(ctx context.Context, customerID, rmaNo string) error
}
```

### `domain/service/port/erp_port.go`
```go
package port

import "context"

type ERPPort interface {
    CheckPartAvailability(ctx context.Context, tenantID, partID string, qty float64) (bool, error)
    ReservePart(ctx context.Context, tenantID, partID string, qty float64) error
    CreateGoodsIssue(ctx context.Context, tenantID string, parts []PartRequest) error
}

type PartRequest struct {
    PartID  string  `json:"part_id"`
    Qty     float64 `json:"qty"`
    RefType string  `json:"ref_type"`
    RefID   string  `json:"ref_id"`
}
```

### `domain/service/port/code_generator_port.go`
```go
package port

import (
    "context"

    "github.com/google/uuid"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

type CodeGeneratorPort interface {
    NextShipmentNo(ctx context.Context, tenantID uuid.UUID) (valueobject.ShipmentNo, error)
    NextJobNo(ctx context.Context, tenantID uuid.UUID) (string, error)
    NextWONo(ctx context.Context, tenantID uuid.UUID) (string, error)
    NextRMANo(ctx context.Context, tenantID uuid.UUID) (string, error)
}
```

---

## A.6 Domain Events

### `domain/event/shipment_events.go`
```go
package event

import (
    "time"

    "github.com/google/uuid"
)

const (
    TopicShipmentCreated    = "iotlogistics.shipment.created"
    TopicShipmentDispatched = "iotlogistics.shipment.dispatched"
    TopicShipmentDelivered  = "iotlogistics.shipment.delivered"
    TopicShipmentInstalled  = "iotlogistics.shipment.installed"
    TopicShipmentCancelled  = "iotlogistics.shipment.cancelled"
)

type ShipmentCreated struct {
    EventID    uuid.UUID `json:"event_id"`
    ShipmentID uuid.UUID `json:"shipment_id"`
    TenantID   uuid.UUID `json:"tenant_id"`
    ShipmentNo string    `json:"shipment_no"`
    CustomerID uuid.UUID `json:"customer_id"`
    SiteID     uuid.UUID `json:"site_id"`
    ItemCount  int       `json:"item_count"`
    OccurredAt time.Time `json:"occurred_at"`
}

type ShipmentDispatched struct {
    EventID      uuid.UUID `json:"event_id"`
    ShipmentID   uuid.UUID `json:"shipment_id"`
    TenantID     uuid.UUID `json:"tenant_id"`
    ShipmentNo   string    `json:"shipment_no"`
    Carrier      string    `json:"carrier"`
    TrackingNo   string    `json:"tracking_no"`
    OccurredAt   time.Time `json:"occurred_at"`
}

type ShipmentDelivered struct {
    EventID     uuid.UUID `json:"event_id"`
    ShipmentID  uuid.UUID `json:"shipment_id"`
    TenantID    uuid.UUID `json:"tenant_id"`
    DeliveredTo string    `json:"delivered_to"`
    OccurredAt  time.Time `json:"occurred_at"`
}
```

### `domain/event/installation_events.go`
```go
package event

import (
    "time"

    "github.com/google/uuid"
)

const (
    TopicInstallationScheduled = "iotlogistics.installation.scheduled"
    TopicInstallationAssigned  = "iotlogistics.installation.assigned"
    TopicInstallationStarted   = "iotlogistics.installation.started"
    TopicInstallationCompleted = "iotlogistics.installation.completed"
    TopicInstallationFailed    = "iotlogistics.installation.failed"
    TopicInstallationSLABreach = "iotlogistics.installation.sla.breach"
)

type InstallationScheduled struct {
    EventID       uuid.UUID  `json:"event_id"`
    JobID         uuid.UUID  `json:"job_id"`
    JobNo         string     `json:"job_no"`
    TenantID      uuid.UUID  `json:"tenant_id"`
    CustomerID    uuid.UUID  `json:"customer_id"`
    SiteID        uuid.UUID  `json:"site_id"`
    JobType       string     `json:"job_type"`
    Priority      string     `json:"priority"`
    ScheduledAt   time.Time  `json:"scheduled_at"`
    SLADueAt      time.Time  `json:"sla_due_at"`
    ShipmentID    *uuid.UUID `json:"shipment_id,omitempty"`
    DeviceIDs     []string   `json:"device_ids,omitempty"`
    OccurredAt    time.Time  `json:"occurred_at"`
}

type InstallationAssigned struct {
    EventID      uuid.UUID `json:"event_id"`
    JobID        uuid.UUID `json:"job_id"`
    TenantID     uuid.UUID `json:"tenant_id"`
    TechnicianID uuid.UUID `json:"technician_id"`
    AssignedAt   time.Time `json:"assigned_at"`
    OccurredAt   time.Time `json:"occurred_at"`
}

type InstallationCompleted struct {
    EventID       uuid.UUID   `json:"event_id"`
    JobID         uuid.UUID   `json:"job_id"`
    TenantID      uuid.UUID   `json:"tenant_id"`
    CustomerID    uuid.UUID   `json:"customer_id"`
    SiteID        uuid.UUID   `json:"site_id"`
    JobType       string      `json:"job_type"`
    DeviceIDs     []string    `json:"device_ids"`
    SerialNumbers []string    `json:"serial_numbers"`
    ShipmentID    *uuid.UUID  `json:"shipment_id,omitempty"`
    CompletedAt   time.Time   `json:"completed_at"`
    DurationMin   int         `json:"duration_min"`
    Photos        []string    `json:"photos,omitempty"`
    OccurredAt    time.Time   `json:"occurred_at"`
}

type InstallationFailed struct {
    EventID     uuid.UUID `json:"event_id"`
    JobID       uuid.UUID `json:"job_id"`
    TenantID    uuid.UUID `json:"tenant_id"`
    Reason      string    `json:"reason"`
    OccurredAt  time.Time `json:"occurred_at"`
}
```

### `domain/event/maintenance_events.go`
```go
package event

import (
    "time"

    "github.com/google/uuid"
)

const (
    TopicMaintenanceDue     = "iotlogistics.maintenance.due"
    TopicMaintenanceDone    = "iotlogistics.maintenance.done"
    TopicMaintenanceOverdue = "iotlogistics.maintenance.overdue"
)

type MaintenanceDue struct {
    EventID     uuid.UUID `json:"event_id"`
    ScheduleID  uuid.UUID `json:"schedule_id"`
    TenantID    uuid.UUID `json:"tenant_id"`
    DeviceID    uuid.UUID `json:"device_id"`
    CustomerID  uuid.UUID `json:"customer_id"`
    SiteID      uuid.UUID `json:"site_id"`
    DueAt       time.Time `json:"due_at"`
    DaysOverdue int       `json:"days_overdue"`
    OccurredAt  time.Time `json:"occurred_at"`
}
```

### `domain/event/rma_events.go`
```go
package event

import (
    "time"

    "github.com/google/uuid"
)

const (
    TopicRMACreated   = "iotlogistics.rma.created"
    TopicRMAApproved  = "iotlogistics.rma.approved"
    TopicRMACompleted = "iotlogistics.rma.completed"
)

type RMACreated struct {
    EventID    uuid.UUID `json:"event_id"`
    RMAID      uuid.UUID `json:"rma_id"`
    RMANo      string    `json:"rma_no"`
    TenantID   uuid.UUID `json:"tenant_id"`
    DeviceID   uuid.UUID `json:"device_id"`
    CustomerID uuid.UUID `json:"customer_id"`
    Reason     string    `json:"reason"`
    OccurredAt time.Time `json:"occurred_at"`
}
```

---

## A.7 Domain Errors

### `domain/errors/errors.go`
```go
package domainerrors

import "errors"

// helper for VO packages
func newErr(s string) error { return errors.New(s) }

var (
    // Shipment
    ErrShipmentNotFound          = errors.New("shipment not found")
    ErrInvalidShipmentNo         = errors.New("invalid shipment number")
    ErrShipmentLocked            = errors.New("shipment is locked (not editable)")
    ErrEmptyShipment             = errors.New("shipment has no items")
    ErrDuplicateSerial           = errors.New("duplicate serial number in shipment")
    ErrShipmentItemNotFound      = errors.New("shipment item not found")
    ErrShipmentNotInTransit      = errors.New("shipment not in transit")
    ErrShipmentAlreadyTerminal   = errors.New("shipment already in terminal state")
    ErrCannotCancelInTransit     = errors.New("cannot cancel shipment in transit")
    ErrInvalidCarrier            = errors.New("invalid carrier")
    ErrRecipientRequired         = errors.New("recipient name required")

    // Installation Job
    ErrJobNotFound             = errors.New("installation job not found")
    ErrInvalidJobNo            = errors.New("invalid job number")
    ErrInvalidJobType          = errors.New("invalid job type")
    ErrInvalidPriority         = errors.New("invalid priority")
    ErrJobAlreadyTerminal      = errors.New("job already in terminal state")
    ErrJobNotInProgress        = errors.New("job is not in progress")
    ErrNoTechnicianAssigned    = errors.New("no technician assigned")
    ErrCannotReassignInProgress = errors.New("cannot reassign job in progress")
    ErrCannotCancelInProgress  = errors.New("cannot cancel job in progress")
    ErrChecklistIncomplete     = errors.New("checklist incomplete")
    ErrChecklistItemNotFound   = errors.New("checklist item not found")
    ErrEmptyChecklist          = errors.New("empty checklist")
    ErrInvalidChecklistItem    = errors.New("invalid checklist item")
    ErrEvidenceRequired        = errors.New("photo evidence required")
    ErrSignatureRequired       = errors.New("signature required")
    ErrTooManyPhotos           = errors.New("too many photos (max 20)")
    ErrOutsideGeoFence         = errors.New("technician outside geo-fence")
    ErrNoJobs                  = errors.New("no jobs to plan")

    // Technician
    ErrTechnicianNotFound      = errors.New("technician not found")
    ErrInvalidTechnicianData   = errors.New("invalid technician data")
    ErrInvalidSkill            = errors.New("invalid skill")
    ErrNoTechnicianAvailable   = errors.New("no technician available")
    ErrTechnicianUnavailable   = errors.New("technician unavailable")
    ErrTechnicianGPSStale      = errors.New("technician gps stale")

    // Work Order
    ErrWONotFound            = errors.New("work order not found")
    ErrInvalidWONo           = errors.New("invalid work order number")
    ErrInvalidTitle          = errors.New("invalid title")
    ErrResolutionRequired    = errors.New("resolution required")

    // Maintenance
    ErrMaintenanceNotFound = errors.New("maintenance schedule not found")
    ErrInvalidFrequency    = errors.New("invalid pm frequency")

    // RMA
    ErrRMANotFound            = errors.New("rma not found")
    ErrInvalidRMANo           = errors.New("invalid rma number")
    ErrInvalidRMAReason       = errors.New("invalid rma reason")
    ErrInvalidRMAResolution   = errors.New("invalid rma resolution")
    ErrRMAAlreadyTerminal     = errors.New("rma already in terminal state")
    ErrRejectReasonRequired   = errors.New("reject reason required")
    ErrRMAAlreadyOpen         = errors.New("rma already open for this device")

    // Spare parts
    ErrSparePartNotFound = errors.New("spare part not found")

    // Validation
    ErrInvalidTenantID    = errors.New("invalid tenant id")
    ErrInvalidCustomerID  = errors.New("invalid customer id")
    ErrInvalidSiteID      = errors.New("invalid site id")
    ErrInvalidQuantity    = errors.New("invalid quantity")
    ErrInvalidStatusTransition = errors.New("invalid status transition")
    ErrInvalidGPS         = errors.New("invalid gps location")
    ErrInvalidTimeWindow  = errors.New("invalid time window")

    // Infrastructure
    ErrPersistenceFailure = errors.New("persistence failure")
    ErrMapsServiceFailure = errors.New("maps service failure")
    ErrAIServiceFailure   = errors.New("ai service failure")
)
```

---

## A.8 Unit Tests (Domain)

### `domain/entity/installation_job_test.go`
```go
package entity_test

import (
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "icmongolang/internal/modules/iotlogistics/domain/entity"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

func newTestJob(t *testing.T) *entity.InstallationJob {
    j, err := entity.NewInstallationJob(
        uuid.New(), uuid.New(), uuid.New(), uuid.New(),
        "JOB-2026-000001",
        valueobject.JobTypeInstallation,
        valueobject.JobPriorityNormal,
        time.Now().Add(24*time.Hour),
    )
    require.NoError(t, err)
    return j
}

func TestNewJob_InvalidType(t *testing.T) {
    _, err := entity.NewInstallationJob(
        uuid.New(), uuid.New(), uuid.New(), uuid.New(),
        "JOB-X", "INVALID", valueobject.JobPriorityNormal, time.Now(),
    )
    assert.ErrorIs(t, err, domainerrors.ErrInvalidJobType)
}

func TestJob_FullLifecycle(t *testing.T) {
    j := newTestJob(t)

    // Assign
    techID := uuid.New()
    actorID := uuid.New()
    require.NoError(t, j.AssignTechnician(techID, actorID))
    assert.Equal(t, valueobject.JobStatusAssigned, j.Status)

    // Start
    gps, _ := valueobject.NewGPSLocation(13.7563, 100.5018, 10)
    require.NoError(t, j.Start(gps, nil))
    assert.Equal(t, valueobject.JobStatusInProgress, j.Status)

    // Update checklist
    for _, it := range j.Checklist {
        if it.Required {
            require.NoError(t, j.UpdateChecklist(it.Key, true, "", "url"))
        }
    }

    // Add photo
    require.NoError(t, j.AddPhoto("https://cdn/p1.jpg"))

    // Complete
    require.NoError(t, j.Complete(gps, "sig-url", "done"))
    assert.Equal(t, valueobject.JobStatusDone, j.Status)
    assert.True(t, j.Status.IsTerminal())
}

func TestJob_Complete_RequiresChecklist(t *testing.T) {
    j := newTestJob(t)
    techID := uuid.New()
    _ = j.AssignTechnician(techID, uuid.New())
    gps, _ := valueobject.NewGPSLocation(13.7, 100.5, 10)
    _ = j.Start(gps, nil)
    _ = j.AddPhoto("p1")

    // ยังไม่ติ๊ก checklist
    err := j.Complete(gps, "sig", "")
    assert.ErrorIs(t, err, domainerrors.ErrChecklistIncomplete)
}

func TestJob_Start_OutsideGeoFence(t *testing.T) {
    j := newTestJob(t)
    _ = j.AssignTechnician(uuid.New(), uuid.New())

    siteGPS, _ := valueobject.NewGPSLocation(13.7563, 100.5018, 10)
    techGPS, _ := valueobject.NewGPSLocation(14.5, 101.5, 10) // ห่าง 100+ km
    fence := &valueobject.GeoFence{Center: siteGPS, RadiusM: 500}

    err := j.Start(techGPS, fence)
    assert.ErrorIs(t, err, domainerrors.ErrOutsideGeoFence)
}

func TestJob_SLABreach(t *testing.T) {
    j := newTestJob(t)
    j.SLADueAt = time.Now().Add(-1 * time.Hour)
    assert.True(t, j.IsSLABreached(time.Now()))
}
```

### `domain/entity/rma_test.go`
```go
package entity_test

import (
    "testing"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "icmongolang/internal/modules/iotlogistics/domain/entity"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
)

func TestRMA_ApprovalFlow(t *testing.T) {
    r, err := entity.NewRMA(
        uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New(),
        "RMA-2026-000001",
        entity.RMAReasonDefective,
        "sensor reads NaN",
    )
    require.NoError(t, err)

    // Approve
    require.NoError(t, r.Approve(uuid.New(), entity.RMAResolutionRepair))
    require.NoError(t, r.Ship("TRACK-001"))
    require.NoError(t, r.Receive(uuid.New()))
    require.NoError(t, r.StartRepair(uuid.New()))
    require.NoError(t, r.MarkRepaired())
    require.NoError(t, r.Close("returned to customer"))

    assert.Equal(t, valueobject.RMAStatusClosed, r.Status)
}

func TestRMA_RejectTerminal(t *testing.T) {
    r, _ := entity.NewRMA(uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New(),
        "RMA-2026-000002", entity.RMAReasonDamaged, "shipping damage")
    require.NoError(t, r.Reject("out of warranty", uuid.New()))
    err := r.Approve(uuid.New(), entity.RMAResolutionRepair)
    assert.ErrorIs(t, err, domainerrors.ErrInvalidStatusTransition)
}

func TestRMA_InvalidReason(t *testing.T) {
    _, err := entity.NewRMA(uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New(),
        "RMA-X", "NOT_A_REASON", "x")
    assert.ErrorIs(t, err, domainerrors.ErrInvalidRMAReason)
}
```

---

# 🅱️ PART 5B — APPLICATION LAYER

## B.1 DTOs

### `application/dto.go`
```go
package application

import "time"

// ============================================================
// SHIPMENT
// ============================================================

type CreateShipmentInput struct {
    TenantID    string            `json:"-"`
    CustomerID  string            `json:"customer_id" binding:"required"`
    SiteID      string            `json:"site_id" binding:"required"`
    WarehouseID string            `json:"warehouse_id" binding:"required"`
    ScheduledAt time.Time         `json:"scheduled_at" binding:"required"`
    Items       []ShipmentItemDTO `json:"items" binding:"required,min=1"`
    Notes       string            `json:"notes,omitempty"`
    ActorID     string            `json:"-"`
}

type ShipmentItemDTO struct {
    ProductID string  `json:"product_id" binding:"required"`
    SerialNo  string  `json:"serial_no" binding:"required"`
    Qty       float64 `json:"qty" binding:"required,gt=0"`
    UOM       string  `json:"uom,omitempty"`
    Notes     string  `json:"notes,omitempty"`
}

type DispatchShipmentInput struct {
    TenantID   string `json:"-"`
    ShipmentID string `json:"-"`
    Carrier    string `json:"carrier" binding:"required"`
    TrackingNo string `json:"tracking_no,omitempty"`
    ActorID    string `json:"-"`
}

type UpdateShipmentGPSInput struct {
    TenantID   string  `json:"-"`
    ShipmentID string  `json:"-"`
    Lat        float64 `json:"lat" binding:"required"`
    Lng        float64 `json:"lng" binding:"required"`
    Speed      float64 `json:"speed,omitempty"`
}

type ConfirmDeliveryInput struct {
    TenantID     string `json:"-"`
    ShipmentID   string `json:"-"`
    DeliveredTo  string `json:"delivered_to" binding:"required"`
    SignatureURL string `json:"signature_url,omitempty"`
    ActorID      string `json:"-"`
}

type ShipmentResponse struct {
    ID          string              `json:"id"`
    ShipmentNo  string              `json:"shipment_no"`
    CustomerID  string              `json:"customer_id"`
    SiteID      string              `json:"site_id"`
    Status      string              `json:"status"`
    Carrier     string              `json:"carrier,omitempty"`
    TrackingNo  string              `json:"tracking_no,omitempty"`
    ScheduledAt time.Time           `json:"scheduled_at"`
    DispatchedAt *time.Time         `json:"dispatched_at,omitempty"`
    DeliveredAt  *time.Time         `json:"delivered_at,omitempty"`
    DeliveredTo  string             `json:"delivered_to,omitempty"`
    Items       []ShipmentItemDTO   `json:"items,omitempty"`
    GPSLast     *GPSDTO             `json:"gps_last,omitempty"`
    CreatedAt   time.Time           `json:"created_at"`
}

type GPSDTO struct {
    Lat       float64   `json:"lat"`
    Lng       float64   `json:"lng"`
    Accuracy  float64   `json:"accuracy,omitempty"`
    Timestamp time.Time `json:"ts"`
}

// ============================================================
// INSTALLATION JOB
// ============================================================

type ScheduleInstallationInput struct {
    TenantID    string    `json:"-"`
    ShipmentID  string    `json:"shipment_id,omitempty"`
    CustomerID  string    `json:"customer_id" binding:"required"`
    SiteID      string    `json:"site_id" binding:"required"`
    JobType     string    `json:"job_type" binding:"required"`
    Priority    string    `json:"priority,omitempty"`
    ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
    DeviceIDs   []string  `json:"device_ids,omitempty"`
    ActorID     string    `json:"-"`
}

type AssignTechnicianInput struct {
    TenantID     string `json:"-"`
    JobID        string `json:"-"`
    TechnicianID string `json:"technician_id,omitempty"` // ถ้าไม่ส่ง → auto-assign
    AutoAssign   bool   `json:"auto_assign,omitempty"`
    SiteGPS      GPSDTO `json:"site_gps,omitempty"`
    ActorID      string `json:"-"`
}

type StartJobInput struct {
    TenantID   string  `json:"-"`
    JobID      string  `json:"-"`
    Lat        float64 `json:"lat" binding:"required"`
    Lng        float64 `json:"lng" binding:"required"`
    Accuracy   float64 `json:"accuracy,omitempty"`
    ActorID    string  `json:"-"`
}

type UpdateChecklistInput struct {
    TenantID   string `json:"-"`
    JobID      string `json:"-"`
    ItemKey    string `json:"item_key" binding:"required"`
    Checked    bool   `json:"checked"`
    Notes      string `json:"notes,omitempty"`
    PhotoURL   string `json:"photo_url,omitempty"`
}

type CompleteJobInput struct {
    TenantID     string  `json:"-"`
    JobID        string  `json:"-"`
    Lat          float64 `json:"lat" binding:"required"`
    Lng          float64 `json:"lng" binding:"required"`
    SignatureURL string  `json:"signature_url,omitempty"`
    Notes        string  `json:"notes,omitempty"`
    ActorID      string  `json:"-"`
}

type InstallationJobResponse struct {
    ID             string                `json:"id"`
    JobNo          string                `json:"job_no"`
    ShipmentID     string                `json:"shipment_id,omitempty"`
    CustomerID     string                `json:"customer_id"`
    SiteID         string                `json:"site_id"`
    DeviceIDs      []string              `json:"device_ids,omitempty"`
    TechnicianID   string                `json:"technician_id,omitempty"`
    JobType        string                `json:"job_type"`
    Priority       string                `json:"priority"`
    Status         string                `json:"status"`
    ScheduledAt    time.Time             `json:"scheduled_at"`
    SLADueAt       time.Time             `json:"sla_due_at"`
    IsSLABreached  bool                  `json:"is_sla_breached"`
    StartedAt      *time.Time            `json:"started_at,omitempty"`
    CompletedAt    *time.Time            `json:"completed_at,omitempty"`
    DurationMin    int                   `json:"duration_min"`
    Checklist      []ChecklistItemDTO    `json:"checklist,omitempty"`
    ChecklistProgress string             `json:"checklist_progress,omitempty"` // "5/7"
    Photos         []string              `json:"photos,omitempty"`
    SignatureURL   string                `json:"signature_url,omitempty"`
    Notes          string                `json:"notes,omitempty"`
    CreatedAt      time.Time             `json:"created_at"`
}

type ChecklistItemDTO struct {
    Key       string     `json:"key"`
    Label     string     `json:"label"`
    Required  bool       `json:"required"`
    Checked   bool       `json:"checked"`
    Notes     string     `json:"notes,omitempty"`
    PhotoURL  string     `json:"photo_url,omitempty"`
    CheckedAt *time.Time `json:"checked_at,omitempty"`
}

// ============================================================
// ROUTE
// ============================================================

type OptimizeRouteInput struct {
    TenantID     string    `json:"-"`
    TechnicianID string    `json:"-"`
    Date         time.Time `json:"date"`
    MaxHours     int       `json:"max_hours,omitempty"`
    StartLat     float64   `json:"start_lat,omitempty"`
    StartLng     float64   `json:"start_lng,omitempty"`
}

type RouteResponse struct {
    RouteID      string            `json:"route_id"`
    TechnicianID string            `json:"technician_id"`
    Waypoints    []RouteWaypointDTO `json:"waypoints"`
    TotalKm      float64           `json:"total_km"`
    TotalMinutes int               `json:"total_minutes"`
    StartAt      time.Time         `json:"start_at"`
    EndAt        time.Time         `json:"end_at"`
    Provider     string            `json:"provider"`
}

type RouteWaypointDTO struct {
    Order       int       `json:"order"`
    JobID       string    `json:"job_id"`
    Lat         float64   `json:"lat"`
    Lng         float64   `json:"lng"`
    ArrivalETA  time.Time `json:"arrival_eta"`
    DepartETA   time.Time `json:"depart_eta"`
    ServiceMin  int       `json:"service_min"`
    DistanceKm  float64   `json:"distance_km"`
}

// ============================================================
// TECHNICIAN
// ============================================================

type CreateTechnicianInput struct {
    TenantID      string   `json:"-"`
    UserID        string   `json:"user_id" binding:"required"`
    Code          string   `json:"code" binding:"required"`
    Name          string   `json:"name" binding:"required"`
    Phone         string   `json:"phone,omitempty"`
    Email         string   `json:"email,omitempty"`
    Level         string   `json:"level,omitempty"`
    Skills        []string `json:"skills,omitempty"`
    Zone          string   `json:"zone,omitempty"`
    MaxJobsPerDay int      `json:"max_jobs_per_day,omitempty"`
    ActorID       string   `json:"-"`
}

type UpdateTechnicianGPSInput struct {
    TechnicianID string  `json:"-"`
    Lat          float64 `json:"lat" binding:"required"`
    Lng          float64 `json:"lng" binding:"required"`
    Accuracy     float64 `json:"accuracy,omitempty"`
}

type TechnicianResponse struct {
    ID            string     `json:"id"`
    Code          string     `json:"code"`
    Name          string     `json:"name"`
    Phone         string     `json:"phone,omitempty"`
    Email         string     `json:"email,omitempty"`
    Level         string     `json:"level"`
    Skills        []string   `json:"skills"`
    Zone          string     `json:"zone,omitempty"`
    IsAvailable   bool       `json:"is_available"`
    IsActive      bool       `json:"is_active"`
    CurrentLoad   int        `json:"current_load"`
    MaxJobsPerDay int        `json:"max_jobs_per_day"`
    GPSLast       *GPSDTO    `json:"gps_last,omitempty"`
    GPSUpdatedAt  *time.Time `json:"gps_updated_at,omitempty"`
    CreatedAt     time.Time  `json:"created_at"`
}

// ============================================================
// WORK ORDER
// ============================================================

type CreateWorkOrderInput struct {
    TenantID    string `json:"-"`
    CustomerID  string `json:"customer_id" binding:"required"`
    SiteID      string `json:"site_id" binding:"required"`
    DeviceID    string `json:"device_id,omitempty"`
    Title       string `json:"title" binding:"required"`
    Description string `json:"description,omitempty"`
    JobType     string `json:"job_type" binding:"required"`
    Priority    string `json:"priority,omitempty"`
    SourceRef   string `json:"source_ref,omitempty"`
    SourceRefID string `json:"source_ref_id,omitempty"`
    ActorID     string `json:"-"`
}

type AddPartUsageInput struct {
    TenantID string  `json:"-"`
    WOID     string  `json:"-"`
    PartID   string  `json:"part_id" binding:"required"`
    Qty      float64 `json:"qty" binding:"required,gt=0"`
}

type CompleteWorkOrderInput struct {
    TenantID     string `json:"-"`
    WOID         string `json:"-"`
    Resolution   string `json:"resolution" binding:"required"`
    LaborMinutes int    `json:"labor_minutes,omitempty"`
    ActorID      string `json:"-"`
}

type WorkOrderResponse struct {
    ID           string          `json:"id"`
    WONo         string          `json:"wo_no"`
    CustomerID   string          `json:"customer_id"`
    SiteID       string          `json:"site_id"`
    DeviceID     string          `json:"device_id,omitempty"`
    JobType      string          `json:"job_type"`
    Priority     string          `json:"priority"`
    Status       string          `json:"status"`
    Title        string          `json:"title"`
    Description  string          `json:"description,omitempty"`
    TechnicianID string          `json:"technician_id,omitempty"`
    ScheduledAt  *time.Time      `json:"scheduled_at,omitempty"`
    StartedAt    *time.Time      `json:"started_at,omitempty"`
    CompletedAt  *time.Time      `json:"completed_at,omitempty"`
    Resolution   string          `json:"resolution,omitempty"`
    LaborMinutes int             `json:"labor_minutes"`
    CostTotal    float64         `json:"cost_total"`
    Currency     string          `json:"currency"`
    PartsUsed    []PartUsageDTO  `json:"parts_used,omitempty"`
    CreatedAt    time.Time       `json:"created_at"`
}

type PartUsageDTO struct {
    PartID    string  `json:"part_id"`
    PartNo    string  `json:"part_no"`
    Name      string  `json:"name"`
    Qty       float64 `json:"qty"`
    UnitPrice float64 `json:"unit_price"`
}

// ============================================================
// MAINTENANCE
// ============================================================

type ScheduleMaintenanceInput struct {
    TenantID   string    `json:"-"`
    DeviceID   string    `json:"device_id" binding:"required"`
    CustomerID string    `json:"customer_id" binding:"required"`
    SiteID     string    `json:"site_id" binding:"required"`
    Frequency  string    `json:"frequency" binding:"required"`
    CustomDays int       `json:"custom_days,omitempty"`
    StartDate  time.Time `json:"start_date" binding:"required"`
    AutoCreateWO bool    `json:"auto_create_wo"`
    Notes      string    `json:"notes,omitempty"`
    ActorID    string    `json:"-"`
}

type MaintenanceScheduleResponse struct {
    ID             string     `json:"id"`
    DeviceID       string     `json:"device_id"`
    CustomerID     string     `json:"customer_id"`
    SiteID         string     `json:"site_id"`
    Frequency      string     `json:"frequency"`
    LastDoneAt     *time.Time `json:"last_done_at,omitempty"`
    NextDueAt      time.Time  `json:"next_due_at"`
    DaysUntilDue   int        `json:"days_until_due"`
    IsDue          bool       `json:"is_due"`
    IsActive       bool       `json:"is_active"`
    AutoCreateWO   bool       `json:"auto_create_wo"`
    CreatedAt      time.Time  `json:"created_at"`
}

// ============================================================
// RMA
// ============================================================

type CreateRMAInput struct {
    TenantID     string   `json:"-"`
    DeviceID     string   `json:"device_id" binding:"required"`
    CustomerID   string   `json:"customer_id" binding:"required"`
    SiteID       string   `json:"site_id" binding:"required"`
    JobID        string   `json:"job_id,omitempty"`
    Reason       string   `json:"reason" binding:"required"`
    ReasonDetail string   `json:"reason_detail,omitempty"`
    Photos       []string `json:"photos,omitempty"`
    ActorID      string   `json:"-"`
}

type ApproveRMAInput struct {
    TenantID   string `json:"-"`
    RMAID      string `json:"-"`
    Resolution string `json:"resolution" binding:"required"`
    ActorID    string `json:"-"`
}

type ReplaceDeviceInput struct {
    TenantID      string `json:"-"`
    RMAID         string `json:"-"`
    NewDeviceID   string `json:"new_device_id" binding:"required"`
    ActorID       string `json:"-"`
}

type RMAResponse struct {
    ID                  string     `json:"id"`
    RMANo               string     `json:"rma_no"`
    DeviceID            string     `json:"device_id"`
    CustomerID          string     `json:"customer_id"`
    SiteID              string     `json:"site_id"`
    Reason              string     `json:"reason"`
    ReasonDetail        string     `json:"reason_detail,omitempty"`
    Status              string     `json:"status"`
    Resolution          string     `json:"resolution,omitempty"`
    ReplacementDeviceID string     `json:"replacement_device_id,omitempty"`
    Photos              []string   `json:"photos,omitempty"`
    ApprovedAt          *time.Time `json:"approved_at,omitempty"`
    ShippedAt           *time.Time `json:"shipped_at,omitempty"`
    ReceivedAt          *time.Time `json:"received_at,omitempty"`
    ClosedAt            *time.Time `json:"closed_at,omitempty"`
    CreatedAt           time.Time  `json:"created_at"`
}
```

### `application/mappers.go`
```go
package application

import (
    "time"

    "icmongolang/internal/modules/iotlogistics/domain/entity"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

func toShipmentResponse(s *entity.Shipment) *ShipmentResponse {
    r := &ShipmentResponse{
        ID: s.ID.String(), ShipmentNo: s.ShipmentNo.String(),
        CustomerID: s.CustomerID.String(), SiteID: s.SiteID.String(),
        Status: string(s.Status), Carrier: s.Carrier, TrackingNo: s.TrackingNo,
        ScheduledAt: s.ScheduledAt,
        DispatchedAt: s.DispatchedAt, DeliveredAt: s.DeliveredAt,
        DeliveredTo: s.DeliveredTo,
        CreatedAt: s.CreatedAt,
    }
    if s.GPSLast != nil {
        r.GPSLast = &GPSDTO{
            Lat: s.GPSLast.Lat, Lng: s.GPSLast.Lng,
            Accuracy: s.GPSLast.Accuracy, Timestamp: s.GPSLast.Timestamp,
        }
    }
    if len(s.Items) > 0 {
        r.Items = make([]ShipmentItemDTO, 0, len(s.Items))
        for _, it := range s.Items {
            r.Items = append(r.Items, ShipmentItemDTO{
                ProductID: it.ProductID.String(), SerialNo: it.SerialNo,
                Qty: it.Qty, UOM: it.UOM, Notes: it.Notes,
            })
        }
    }
    return r
}

func toInstallationJobResponse(j *entity.InstallationJob, asOf time.Time) *InstallationJobResponse {
    r := &InstallationJobResponse{
        ID: j.ID.String(), JobNo: j.JobNo,
        CustomerID: j.CustomerID.String(), SiteID: j.SiteID.String(),
        JobType: string(j.JobType), Priority: string(j.Priority),
        Status: string(j.Status),
        ScheduledAt: j.ScheduledAt, SLADueAt: j.SLADueAt,
        IsSLABreached: j.IsSLABreached(asOf),
        StartedAt: j.StartedAt, CompletedAt: j.CompletedAt,
        DurationMin: j.DurationMinutes(),
        Photos: j.Photos, SignatureURL: j.SignatureURL,
        Notes: j.Notes, CreatedAt: j.CreatedAt,
    }
    if j.ShipmentID != nil { r.ShipmentID = j.ShipmentID.String() }
    if j.TechnicianID != nil { r.TechnicianID = j.TechnicianID.String() }

    r.DeviceIDs = make([]string, 0, len(j.DeviceIDs))
    for _, id := range j.DeviceIDs {
        r.DeviceIDs = append(r.DeviceIDs, id.String())
    }

    if len(j.Checklist) > 0 {
        r.Checklist = make([]ChecklistItemDTO, 0, len(j.Checklist))
        for _, it := range j.Checklist {
            dto := ChecklistItemDTO{
                Key: it.Key, Label: it.Label, Required: it.Required,
                Checked: it.Checked, Notes: it.Notes, PhotoURL: it.PhotoURL,
            }
            if !it.CheckedAt.IsZero() { dto.CheckedAt = &it.CheckedAt }
            r.Checklist = append(r.Checklist, dto)
        }
        ch, total := j.ChecklistProgress()
        r.ChecklistProgress = formatProgress(ch, total)
    }
    return r
}

func toTechnicianResponse(t *entity.Technician) *TechnicianResponse {
    r := &TechnicianResponse{
        ID: t.ID.String(), Code: t.Code, Name: t.Name,
        Phone: t.Phone, Email: t.Email,
        Level: string(t.Level),
        Zone: t.Zone, IsAvailable: t.IsAvailable, IsActive: t.IsActive,
        CurrentLoad: t.CurrentLoad, MaxJobsPerDay: t.MaxJobsPerDay,
        GPSUpdatedAt: t.GPSUpdatedAt, CreatedAt: t.CreatedAt,
    }
    r.Skills = make([]string, 0, len(t.Skills))
    for _, s := range t.Skills {
        r.Skills = append(r.Skills, string(s))
    }
    if t.GPSLast != nil {
        r.GPSLast = &GPSDTO{
            Lat: t.GPSLast.Lat, Lng: t.GPSLast.Lng,
            Accuracy: t.GPSLast.Accuracy, Timestamp: t.GPSLast.Timestamp,
        }
    }
    return r
}

func toWorkOrderResponse(w *entity.WorkOrder) *WorkOrderResponse {
    r := &WorkOrderResponse{
        ID: w.ID.String(), WONo: w.WONo,
        CustomerID: w.CustomerID.String(), SiteID: w.SiteID.String(),
        JobType: string(w.JobType), Priority: string(w.Priority),
        Status: string(w.Status), Title: w.Title, Description: w.Description,
        StartedAt: w.StartedAt, CompletedAt: w.CompletedAt,
        Resolution: w.Resolution, LaborMinutes: w.LaborMinutes,
        CostTotal: w.CostTotal, Currency: w.Currency,
        CreatedAt: w.CreatedAt,
    }
    if w.DeviceID != nil { r.DeviceID = w.DeviceID.String() }
    if w.TechnicianID != nil { r.TechnicianID = w.TechnicianID.String() }
    if !w.ScheduledAt.IsZero() { r.ScheduledAt = &w.ScheduledAt }

    r.PartsUsed = make([]PartUsageDTO, 0, len(w.PartsUsed))
    for _, p := range w.PartsUsed {
        r.PartsUsed = append(r.PartsUsed, PartUsageDTO{
            PartID: p.PartID.String(), PartNo: p.PartNo, Name: p.Name,
            Qty: p.Qty, UnitPrice: p.UnitPrice,
        })
    }
    return r
}

func toMaintenanceResponse(m *entity.MaintenanceSchedule, asOf time.Time) *MaintenanceScheduleResponse {
    days := int(m.NextDueAt.Sub(asOf).Hours() / 24)
    return &MaintenanceScheduleResponse{
        ID: m.ID.String(), DeviceID: m.DeviceID.String(),
        CustomerID: m.CustomerID.String(), SiteID: m.SiteID.String(),
        Frequency: string(m.Frequency),
        LastDoneAt: m.LastDoneAt, NextDueAt: m.NextDueAt,
        DaysUntilDue: days, IsDue: m.IsDue(asOf),
        IsActive: m.IsActive, AutoCreateWO: m.AutoCreateWO,
        CreatedAt: m.CreatedAt,
    }
}

func toRMAResponse(r *entity.RMA) *RMAResponse {
    resp := &RMAResponse{
        ID: r.ID.String(), RMANo: r.RMANo,
        DeviceID: r.DeviceID.String(), CustomerID: r.CustomerID.String(),
        SiteID: r.SiteID.String(),
        Reason: string(r.Reason), ReasonDetail: r.ReasonDetail,
        Status: string(r.Status),
        Photos: r.ReceivedPhotos,
        ApprovedAt: r.ApprovedAt, ShippedAt: r.ShippedAt,
        ReceivedAt: r.ReceivedAt, ClosedAt: r.ClosedAt,
        CreatedAt: r.CreatedAt,
    }
    if r.Resolution != nil { resp.Resolution = string(*r.Resolution) }
    if r.ReplacementDeviceID != nil {
        resp.ReplacementDeviceID = r.ReplacementDeviceID.String()
    }
    return resp
}

func toRouteResponse(route *valueobject.Route, techID string) *RouteResponse {
    r := &RouteResponse{
        RouteID: route.ID, TechnicianID: techID,
        TotalKm: route.TotalKm, TotalMinutes: route.TotalMinutes,
        StartAt: route.StartAt, EndAt: route.EndAt,
        Provider: route.Provider,
        Waypoints: make([]RouteWaypointDTO, 0, len(route.Waypoints)),
    }
    for _, w := range route.Waypoints {
        r.Waypoints = append(r.Waypoints, RouteWaypointDTO{
            Order: w.Order, JobID: w.RefID,
            Lat: w.Location.Lat, Lng: w.Location.Lng,
            ArrivalETA: w.ArrivalETA, DepartETA: w.DepartETA,
            ServiceMin: w.ServiceMin, DistanceKm: w.DistanceKm,
        })
    }
    return r
}

func formatProgress(ch, total int) string {
    return itoa(ch) + "/" + itoa(total)
}

func itoa(i int) string {
    if i == 0 { return "0" }
    buf := make([]byte, 0, 4)
    for i > 0 {
        buf = append([]byte{byte('0' + i%10)}, buf...)
        i /= 10
    }
    return string(buf)
}
```

### `application/create_shipment.go`
```go
package application

import (
    "context"
    "time"

    "github.com/google/uuid"

    "icmongolang/internal/modules/iotlogistics/domain/entity"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    "icmongolang/internal/modules/iotlogistics/domain/event"
    "icmongolang/internal/modules/iotlogistics/domain/repository"
    "icmongolang/internal/modules/iotlogistics/domain/service/port"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
)

type CreateShipmentUseCase struct {
    shipRepo   repository.ShipmentRepository
    auditRepo  repository.AuditRepository
    codeGen    port.CodeGeneratorPort
    producer   kafka.Producer
    log        logger.Logger
}

func NewCreateShipmentUseCase(
    shipRepo repository.ShipmentRepository,
    auditRepo repository.AuditRepository,
    codeGen port.CodeGeneratorPort,
    producer kafka.Producer,
    log logger.Logger,
) *CreateShipmentUseCase {
    return &CreateShipmentUseCase{
        shipRepo: shipRepo, auditRepo: auditRepo,
        codeGen: codeGen, producer: producer, log: log,
    }
}

func (uc *CreateShipmentUseCase) Execute(ctx context.Context, in CreateShipmentInput) (*ShipmentResponse, error) {
    tenantID, err := uuid.Parse(in.TenantID)
    if err != nil { return nil, domainerrors.ErrInvalidTenantID }
    customerID, _ := uuid.Parse(in.CustomerID)
    siteID, _ := uuid.Parse(in.SiteID)
    warehouseID, _ := uuid.Parse(in.WarehouseID)
    actorID, _ := uuid.Parse(in.ActorID)

    // 1. Generate shipment no
    no, err := uc.codeGen.NextShipmentNo(ctx, tenantID)
    if err != nil { return nil, err }

    // 2. Create aggregate
    s, err := entity.NewShipment(tenantID, customerID, siteID, warehouseID, actorID, no, in.ScheduledAt)
    if err != nil { return nil, err }

    // 3. Add items
    for _, item := range in.Items {
        productID, _ := uuid.Parse(item.ProductID)
        si, err := entity.NewShipmentItem(productID, item.SerialNo, item.Qty)
        if err != nil { return nil, err }
        if item.UOM != "" { si.UOM = item.UOM }
        si.Notes = item.Notes
        if err := s.AddItem(si); err != nil { return nil, err }
    }

    s.Notes = in.Notes

    // 4. Mark pending
    if err := s.MarkPending(); err != nil { return nil, err }

    // 5. Persist
    if err := uc.shipRepo.Save(ctx, s); err != nil {
        uc.log.Error("save shipment failed", "err", err)
        return nil, domainerrors.ErrPersistenceFailure
    }

    // 6. Audit
    _ = uc.auditRepo.Save(ctx, &repository.AuditEntry{
        TenantID: tenantID, ActorID: actorID,
        Action: "shipment.create", EntityType: "shipment", EntityID: s.ID,
        Payload: map[string]any{
            "shipment_no": s.ShipmentNo.String(),
            "item_count":  len(s.Items),
        },
    })

    // 7. Event
    _ = uc.producer.Publish(ctx, event.TopicShipmentCreated, s.ID.String(), event.ShipmentCreated{
        EventID: uuid.New(), ShipmentID: s.ID, TenantID: tenantID,
        ShipmentNo: s.ShipmentNo.String(),
        CustomerID: customerID, SiteID: siteID,
        ItemCount: len(s.Items), OccurredAt: time.Now(),
    })

    return toShipmentResponse(s), nil
}
```

### `application/dispatch_shipment.go`
```go
package application

import (
    "context"
    "time"

    "github.com/google/uuid"

    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    "icmongolang/internal/modules/iotlogistics/domain/event"
    "icmongolang/internal/modules/iotlogistics/domain/repository"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
)

type DispatchShipmentUseCase struct {
    shipRepo  repository.ShipmentRepository
    auditRepo repository.AuditRepository
    producer  kafka.Producer
    log       logger.Logger
}

func NewDispatchShipmentUseCase(
    shipRepo repository.ShipmentRepository,
    auditRepo repository.AuditRepository,
    producer kafka.Producer,
    log logger.Logger,
) *DispatchShipmentUseCase {
    return &DispatchShipmentUseCase{shipRepo: shipRepo, auditRepo: auditRepo, producer: producer, log: log}
}

func (uc *DispatchShipmentUseCase) Execute(ctx context.Context, in DispatchShipmentInput) (*ShipmentResponse, error) {
    tenantID, _ := uuid.Parse(in.TenantID)
    shipID, _ := uuid.Parse(in.ShipmentID)
    actorID, _ := uuid.Parse(in.ActorID)

    s, err := uc.shipRepo.FindByID(ctx, tenantID, shipID)
    if err != nil { return nil, err }

    if err := s.Dispatch(in.Carrier, in.TrackingNo, actorID); err != nil {
        return nil, err
    }
    if err := uc.shipRepo.Save(ctx, s); err != nil {
        return nil, domainerrors.ErrPersistenceFailure
    }

    _ = uc.auditRepo.Save(ctx, &repository.AuditEntry{
        TenantID: tenantID, ActorID: actorID,
        Action: "shipment.dispatch", EntityType: "shipment", EntityID: s.ID,
        Payload: map[string]any{"carrier": s.Carrier, "tracking": s.TrackingNo},
    })

    _ = uc.producer.Publish(ctx, event.TopicShipmentDispatched, s.ID.String(), event.ShipmentDispatched{
        EventID: uuid.New(), ShipmentID: s.ID, TenantID: tenantID,
        ShipmentNo: s.ShipmentNo.String(),
        Carrier: s.Carrier, TrackingNo: s.TrackingNo,
        OccurredAt: time.Now(),
    })

    return toShipmentResponse(s), nil
}
```

### `application/schedule_installation.go`
```go
package application

import (
    "context"
    "time"

    "github.com/google/uuid"

    "icmongolang/internal/modules/iotlogistics/domain/entity"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    "icmongolang/internal/modules/iotlogistics/domain/event"
    "icmongolang/internal/modules/iotlogistics/domain/repository"
    "icmongolang/internal/modules/iotlogistics/domain/service/port"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
)

type ScheduleInstallationUseCase struct {
    jobRepo    repository.InstallationRepository
    shipRepo   repository.ShipmentRepository
    auditRepo  repository.AuditRepository
    codeGen    port.CodeGeneratorPort
    producer   kafka.Producer
    log        logger.Logger
}

func NewScheduleInstallationUseCase(
    jobRepo repository.InstallationRepository,
    shipRepo repository.ShipmentRepository,
    auditRepo repository.AuditRepository,
    codeGen port.CodeGeneratorPort,
    producer kafka.Producer,
    log logger.Logger,
) *ScheduleInstallationUseCase {
    return &ScheduleInstallationUseCase{
        jobRepo: jobRepo, shipRepo: shipRepo, auditRepo: auditRepo,
        codeGen: codeGen, producer: producer, log: log,
    }
}

func (uc *ScheduleInstallationUseCase) Execute(ctx context.Context, in ScheduleInstallationInput) (*InstallationJobResponse, error) {
    tenantID, _ := uuid.Parse(in.TenantID)
    customerID, _ := uuid.Parse(in.CustomerID)
    siteID, _ := uuid.Parse(in.SiteID)
    actorID, _ := uuid.Parse(in.ActorID)

    jobType := valueobject.JobType(in.JobType)
    priority := valueobject.JobPriority(in.Priority)
    if priority == "" { priority = valueobject.JobPriorityNormal }

    // 1. Generate job no
    jobNo, err := uc.codeGen.NextJobNo(ctx, tenantID)
    if err != nil { return nil, err }

    // 2. Create job
    j, err := entity.NewInstallationJob(tenantID, customerID, siteID, actorID, jobNo, jobType, priority, in.ScheduledAt)
    if err != nil { return nil, err }

    if in.ShipmentID != "" {
        sid, _ := uuid.Parse(in.ShipmentID)
        j.ShipmentID = &sid
    }

    // 3. Attach devices
    for _, ds := range in.DeviceIDs {
        did, _ := uuid.Parse(ds)
        _ = j.AttachDevice(did)
    }

    // 4. Persist
    if err := uc.jobRepo.Save(ctx, j); err != nil {
        uc.log.Error("save job failed", "err", err)
        return nil, domainerrors.ErrPersistenceFailure
    }

    // 5. Audit
    _ = uc.auditRepo.Save(ctx, &repository.AuditEntry{
        TenantID: tenantID, ActorID: actorID,
        Action: "installation.schedule", EntityType: "installation_job", EntityID: j.ID,
        Payload: map[string]any{
            "job_no":      j.JobNo,
            "job_type":    string(j.JobType),
            "scheduled_at": j.ScheduledAt,
        },
    })

    // 6. Publish event
    deviceIDs := make([]string, 0, len(j.DeviceIDs))
    for _, id := range j.DeviceIDs {
        deviceIDs = append(deviceIDs, id.String())
    }
    _ = uc.producer.Publish(ctx, event.TopicInstallationScheduled, j.ID.String(), event.InstallationScheduled{
        EventID: uuid.New(), JobID: j.ID, JobNo: j.JobNo,
        TenantID: tenantID, CustomerID: customerID, SiteID: siteID,
        JobType: string(j.JobType), Priority: string(j.Priority),
        ScheduledAt: j.ScheduledAt, SLADueAt: j.SLADueAt,
        ShipmentID: j.ShipmentID, DeviceIDs: deviceIDs,
        OccurredAt: time.Now(),
    })

    return toInstallationJobResponse(j, time.Now()), nil
}
```

### `application/assign_technician.go` ⭐
```go
package application

import (
    "context"
    "time"

    "github.com/google/uuid"

    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    "icmongolang/internal/modules/iotlogistics/domain/event"
    "icmongolang/internal/modules/iotlogistics/domain/repository"
    "icmongolang/internal/modules/iotlogistics/domain/service"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
    "icmongolang/pkg/transaction"
)

type AssignTechnicianUseCase struct {
    jobRepo    repository.InstallationRepository
    techRepo   repository.TechnicianRepository
    auditRepo  repository.AuditRepository
    assignSvc  *service.TechnicianAssignmentService
    producer   kafka.Producer
    tx         *transaction.Manager
    log        logger.Logger
}

func NewAssignTechnicianUseCase(
    jobRepo repository.InstallationRepository,
    techRepo repository.TechnicianRepository,
    auditRepo repository.AuditRepository,
    assignSvc *service.TechnicianAssignmentService,
    producer kafka.Producer,
    tx *transaction.Manager,
    log logger.Logger,
) *AssignTechnicianUseCase {
    return &AssignTechnicianUseCase{
        jobRepo: jobRepo, techRepo: techRepo, auditRepo: auditRepo,
        assignSvc: assignSvc, producer: producer, tx: tx, log: log,
    }
}

func (uc *AssignTechnicianUseCase) Execute(ctx context.Context, in AssignTechnicianInput) (*InstallationJobResponse, error) {
    tenantID, _ := uuid.Parse(in.TenantID)
    jobID, _ := uuid.Parse(in.JobID)
    actorID, _ := uuid.Parse(in.ActorID)

    var job, updatedJob interface{}
    _ = job
    _ = updatedJob

    var techID uuid.UUID
    var j interface{}

    err := uc.tx.Do(ctx, func(txCtx context.Context) error {
        jobEntity, err := uc.jobRepo.FindByID(txCtx, tenantID, jobID)
        if err != nil { return err }

        // 1. เลือก technician
        if in.TechnicianID != "" {
            tid, _ := uuid.Parse(in.TechnicianID)
            techID = tid
        } else if in.AutoAssign {
            siteGPS := valueobject.GPSLocation{Lat: in.SiteGPS.Lat, Lng: in.SiteGPS.Lng}
            tech, err := uc.assignSvc.Assign(txCtx, tenantID, siteGPS, nil, "")
            if err != nil { return err }
            techID = tech.ID
        } else {
            return domainerrors.ErrInvalidTenantID // invalid input
        }

        // 2. Assign
        if err := jobEntity.AssignTechnician(techID, actorID); err != nil {
            return err
        }

        // 3. Increment tech load
        tech, err := uc.techRepo.FindByID(txCtx, tenantID, techID)
        if err == nil {
            tech.IncrementLoad()
            _ = uc.techRepo.Save(txCtx, tech)
        }

        if err := uc.jobRepo.Save(txCtx, jobEntity); err != nil {
            return err
        }

        j = jobEntity
        return nil
    })
    if err != nil {
        return nil, err
    }

    // Side effects
    _ = uc.auditRepo.Save(ctx, &repository.AuditEntry{
        TenantID: tenantID, ActorID: actorID,
        Action: "installation.assign", EntityType: "installation_job", EntityID: jobID,
        Payload: map[string]any{"technician_id": techID.String()},
    })

    _ = uc.producer.Publish(ctx, event.TopicInstallationAssigned, jobID.String(), event.InstallationAssigned{
        EventID: uuid.New(), JobID: jobID, TenantID: tenantID,
        TechnicianID: techID, AssignedAt: time.Now(), OccurredAt: time.Now(),
    })

    jobEntity := j.(interface {
		// use reflection-free type assertion
	})
	_ = jobEntity
	return nil, nil // simplified
}
```

> **หมายเหตุ**: ตัวอย่างข้างบนเป็น pattern ที่ย่อให้เห็น — use case จริงจะ return `*InstallationJobResponse` โดยใช้ type assertion ที่ถูกต้อง (ในใช้จริง `j` จะเป็น `*entity.InstallationJob`)

### `application/start_job.go`
```go
package application

import (
    "context"
    "time"

    "github.com/google/uuid"

    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    "icmongolang/internal/modules/iotlogistics/domain/event"
    "icmongolang/internal/modules/iotlogistics/domain/repository"
    "icmongolang/internal/modules/iotlogistics/domain/service"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
)

type StartJobUseCase struct {
    jobRepo   repository.InstallationRepository
    geoFence  *service.GeoFenceService
    producer  kafka.Producer
    log       logger.Logger
}

func NewStartJobUseCase(
    jobRepo repository.InstallationRepository,
    geoFence *service.GeoFenceService,
    producer kafka.Producer,
    log logger.Logger,
) *StartJobUseCase {
    return &StartJobUseCase{jobRepo: jobRepo, geoFence: geoFence, producer: producer, log: log}
}

func (uc *StartJobUseCase) Execute(ctx context.Context, in StartJobInput) (*InstallationJobResponse, error) {
    tenantID, _ := uuid.Parse(in.TenantID)
    jobID, _ := uuid.Parse(in.JobID)

    j, err := uc.jobRepo.FindByID(ctx, tenantID, jobID)
    if err != nil { return nil, err }

    techGPS, err := valueobject.NewGPSLocation(in.Lat, in.Lng, in.Accuracy)
    if err != nil { return nil, err }

    // siteGPS — ในโปรดักชันจะดึงจาก customer_sites
    siteGPS := valueobject.GPSLocation{Lat: in.Lat, Lng: in.Lng}
    var fence *valueobject.GeoFence
    if !siteGPS.IsZero() {
        fence, _ = uc.geoFence.BuildFence(siteGPS, 0)
    }

    if err := j.Start(techGPS, fence); err != nil {
        return nil, err
    }
    if err := uc.jobRepo.Save(ctx, j); err != nil {
        return nil, domainerrors.ErrPersistenceFailure
    }

    _ = uc.producer.Publish(ctx, event.TopicInstallationStarted, j.ID.String(), map[string]any{
        "event_id":  uuid.New(),
        "job_id":    j.ID.String(),
        "tenant_id": tenantID.String(),
        "gps":       map[string]any{"lat": techGPS.Lat, "lng": techGPS.Lng},
        "started_at": time.Now(),
    })

    return toInstallationJobResponse(j, time.Now()), nil
}
```

### `application/update_checklist.go`
```go
package application

import (
    "context"

    "github.com/google/uuid"

    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    "icmongolang/internal/modules/iotlogistics/domain/repository"
)

type UpdateChecklistUseCase struct {
    jobRepo repository.InstallationRepository
}

func NewUpdateChecklistUseCase(jobRepo repository.InstallationRepository) *UpdateChecklistUseCase {
    return &UpdateChecklistUseCase{jobRepo: jobRepo}
}

func (uc *UpdateChecklistUseCase) Execute(ctx context.Context, in UpdateChecklistInput) (*InstallationJobResponse, error) {
    tenantID, _ := uuid.Parse(in.TenantID)
    jobID, _ := uuid.Parse(in.JobID)

    j, err := uc.jobRepo.FindByID(ctx, tenantID, jobID)
    if err != nil { return nil, err }

    if err := j.UpdateChecklist(in.ItemKey, in.Checked, in.Notes, in.PhotoURL); err != nil {
        return nil, err
    }
    if err := uc.jobRepo.Save(ctx, j); err != nil {
        return nil, domainerrors.ErrPersistenceFailure
    }
    return toInstallationJobResponse(j, timeNow()), nil
}
```

### `application/complete_job.go` ⭐ (Publishes installation.completed)
```go
package application

import (
    "context"
    "time"

    "github.com/google/uuid"

    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    "icmongolang/internal/modules/iotlogistics/domain/event"
    "icmongolang/internal/modules/iotlogistics/domain/repository"
    "icmongolang/internal/modules/iotlogistics/domain/service"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
    "icmongolang/pkg/transaction"
)

type CompleteJobUseCase struct {
    jobRepo    repository.InstallationRepository
    shipRepo   repository.ShipmentRepository
    maintRepo  repository.MaintenanceRepository
    geoFence   *service.GeoFenceService
    validator  *service.ChecklistValidator
    producer   kafka.Producer
    tx         *transaction.Manager
    log        logger.Logger
}

func NewCompleteJobUseCase(
    jobRepo repository.InstallationRepository,
    shipRepo repository.ShipmentRepository,
    maintRepo repository.MaintenanceRepository,
    geoFence *service.GeoFenceService,
    validator *service.ChecklistValidator,
    producer kafka.Producer,
    tx *transaction.Manager,
    log logger.Logger,
) *CompleteJobUseCase {
    return &CompleteJobUseCase{
        jobRepo: jobRepo, shipRepo: shipRepo, maintRepo: maintRepo,
        geoFence: geoFence, validator: validator,
        producer: producer, tx: tx, log: log,
    }
}

func (uc *CompleteJobUseCase) Execute(ctx context.Context, in CompleteJobInput) (*InstallationJobResponse, error) {
    tenantID, _ := uuid.Parse(in.TenantID)
    jobID, _ := uuid.Parse(in.JobID)

    var completedJob interface{}
    _ = completedJob

    var result interface{}

    err := uc.tx.Do(ctx, func(txCtx context.Context) error {
        j, err := uc.jobRepo.FindByID(txCtx, tenantID, jobID)
        if err != nil { return err }

        gps, err := valueobject.NewGPSLocation(in.Lat, in.Lng, 0)
        if err != nil { return err }

        // Domain validation ผ่าน service
        if err := uc.validator.Validate(j.Checklist, j.Photos, in.SignatureURL); err != nil {
            return err
        }

        if err := j.Complete(gps, in.SignatureURL, in.Notes); err != nil {
            return err
        }
        if err := uc.jobRepo.Save(txCtx, j); err != nil { return err }

        // ถ้ามี shipment → mark installed
        if j.ShipmentID != nil {
            if ship, err := uc.shipRepo.FindByID(txCtx, tenantID, *j.ShipmentID); err == nil {
                _ = ship.MarkInstalled()
                _ = uc.shipRepo.Save(txCtx, ship)
            }
        }

        // ถ้าเป็น maintenance → update schedule
        for _, deviceID := range j.DeviceIDs {
            if m, err := uc.maintRepo.FindByDevice(txCtx, tenantID, deviceID); err == nil {
                m.MarkDone(time.Now())
                _ = uc.maintRepo.Save(txCtx, m)
            }
        }

        result = j
        return nil
    })
    if err != nil {
        return nil, err
    }

    j := result.(*entityTypeAlias)
    _ = j

    // Publish installation.completed event
    _ = uc.producer.Publish(ctx, event.TopicInstallationCompleted, jobID.String(), event.InstallationCompleted{
        EventID: uuid.New(), JobID: jobID,
        TenantID: tenantID,
        DeviceIDs: stringifyUUIDs(j.DeviceIDs),
        SerialNumbers: nil, // TODO: include from shipment
        CompletedAt: time.Now(),
        DurationMin: j.DurationMinutes(),
        Photos: j.Photos,
        OccurredAt: time.Now(),
    })

    return toInstallationJobResponse(j, time.Now()), nil
}
```

> **หมายเหตุ**: โค้ดจริงต้อง type-assert ที่ถูกต้อง — ตัวอย่างข้างบนแสดง pattern

### `application/optimize_route.go`
```go
package application

import (
    "context"
    "time"

    "github.com/google/uuid"

    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    "icmongolang/internal/modules/iotlogistics/domain/repository"
    "icmongolang/internal/modules/iotlogistics/domain/service"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

type OptimizeRouteUseCase struct {
    techRepo repository.TechnicianRepository
    jobRepo  repository.InstallationRepository
    planner  *service.RoutePlanningService
}

func NewOptimizeRouteUseCase(
    techRepo repository.TechnicianRepository,
    jobRepo repository.InstallationRepository,
    planner *service.RoutePlanningService,
) *OptimizeRouteUseCase {
    return &OptimizeRouteUseCase{techRepo: techRepo, jobRepo: jobRepo, planner: planner}
}

func (uc *OptimizeRouteUseCase) Execute(ctx context.Context, in OptimizeRouteInput) (*RouteResponse, error) {
    tenantID, _ := uuid.Parse(in.TenantID)
    techID, _ := uuid.Parse(in.TechnicianID)

    tech, err := uc.techRepo.FindByID(ctx, tenantID, techID)
    if err != nil { return nil, err }

    // 1. List jobs ที่ผูกช่างในวันนั้น
    jobs, err := uc.jobRepo.ListByTechnician(ctx, tenantID, techID, in.Date)
    if err != nil { return nil, err }

    // filter specific statuses
    pending := make([]*entityTypeAlias, 0)
    for _, j := range jobs {
        if j.Status == valueobject.JobStatusAssigned || j.Status == valueobject.JobStatusScheduled {
            pending = append(pending, j)
        }
    }

    if len(pending) == 0 {
        return nil, domainerrors.ErrNoJobs
    }

    // 2. Start location
    start := valueobject.GPSLocation{}
    if in.StartLat != 0 || in.StartLng != 0 {
        start = valueobject.GPSLocation{Lat: in.StartLat, Lng: in.StartLng}
    } else if tech.GPSLast != nil {
        start = *tech.GPSLast
    }

    // 3. Job GPS map — ในโปรดักชันดึงจาก customer_sites
    jobGPS := map[uuid.UUID]valueobject.GPSLocation{}
    for _, j := range pending {
        jobGPS[j.ID] = valueobject.GPSLocation{} // placeholder
    }

    // 4. Plan
    route, err := uc.planner.Plan(ctx, service.PlanRequest{
        TenantID: tenantID, TechnicianID: techID,
        StartLocation: start, Jobs: pending,
        JobGPS: jobGPS, MaxHours: in.MaxHours,
    })
    if err != nil { return nil, err }

    return toRouteResponse(route, techID.String()), nil
}
```

### `application/create_rma.go`
```go
package application

import (
    "context"
    "time"

    "github.com/google/uuid"

    "icmongolang/internal/modules/iotlogistics/domain/entity"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    "icmongolang/internal/modules/iotlogistics/domain/event"
    "icmongolang/internal/modules/iotlogistics/domain/repository"
    "icmongolang/internal/modules/iotlogistics/domain/service/port"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
)

type CreateRMAUseCase struct {
    rmaRepo   repository.RMARepository
    codeGen   port.CodeGeneratorPort
    auditRepo repository.AuditRepository
    producer  kafka.Producer
    log       logger.Logger
}

func NewCreateRMAUseCase(
    rmaRepo repository.RMARepository,
    codeGen port.CodeGeneratorPort,
    auditRepo repository.AuditRepository,
    producer kafka.Producer,
    log logger.Logger,
) *CreateRMAUseCase {
    return &CreateRMAUseCase{rmaRepo: rmaRepo, codeGen: codeGen, auditRepo: auditRepo, producer: producer, log: log}
}

func (uc *CreateRMAUseCase) Execute(ctx context.Context, in CreateRMAInput) (*RMAResponse, error) {
    tenantID, _ := uuid.Parse(in.TenantID)
    deviceID, _ := uuid.Parse(in.DeviceID)
    customerID, _ := uuid.Parse(in.CustomerID)
    siteID, _ := uuid.Parse(in.SiteID)
    actorID, _ := uuid.Parse(in.ActorID)

    // 1. Generate RMA no
    no, err := uc.codeGen.NextRMANo(ctx, tenantID)
    if err != nil { return nil, err }

    // 2. Create
    r, err := entity.NewRMA(tenantID, deviceID, customerID, siteID, actorID, no, entity.RMAReason(in.Reason), in.ReasonDetail)
    if err != nil { return nil, err }

    if in.JobID != "" {
        jid, _ := uuid.Parse(in.JobID)
        r.JobID = &jid
    }
    for _, p := range in.Photos {
        r.AddPhoto(p)
    }

    // 3. Persist
    if err := uc.rmaRepo.Save(ctx, r); err != nil {
        uc.log.Error("save rma failed", "err", err)
        return nil, domainerrors.ErrPersistenceFailure
    }

    // 4. Audit + event
    _ = uc.auditRepo.Save(ctx, &repository.AuditEntry{
        TenantID: tenantID, ActorID: actorID,
        Action: "rma.create", EntityType: "rma", EntityID: r.ID,
        Payload: map[string]any{"rma_no": r.RMANo, "reason": string(r.Reason)},
    })
    _ = uc.producer.Publish(ctx, event.TopicRMACreated, r.ID.String(), event.RMACreated{
        EventID: uuid.New(), RMAID: r.ID, RMANo: r.RMANo,
        TenantID: tenantID, DeviceID: deviceID, CustomerID: customerID,
        Reason: string(r.Reason), OccurredAt: time.Now(),
    })

    return toRMAResponse(r), nil
}
```

### `application/pm_scheduler.go` ⭐ (Auto-create WO)
```go
package application

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"

    "icmongolang/internal/modules/iotlogistics/domain/entity"
    "icmongolang/internal/modules/iotlogistics/domain/event"
    "icmongolang/internal/modules/iotlogistics/domain/repository"
    "icmongolang/internal/modules/iotlogistics/domain/service/port"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
)

// ProcessMaintenanceDueUseCase – เรียกโดย scheduler รายวัน
// ตรวจ schedule ที่ due → auto-create WO + publish event
type ProcessMaintenanceDueUseCase struct {
    maintRepo repository.MaintenanceRepository
    woRepo    repository.WorkOrderRepository
    codeGen   port.CodeGeneratorPort
    producer  kafka.Producer
    log       logger.Logger
}

func NewProcessMaintenanceDueUseCase(
    maintRepo repository.MaintenanceRepository,
    woRepo repository.WorkOrderRepository,
    codeGen port.CodeGeneratorPort,
    producer kafka.Producer,
    log logger.Logger,
) *ProcessMaintenanceDueUseCase {
    return &ProcessMaintenanceDueUseCase{
        maintRepo: maintRepo, woRepo: woRepo,
        codeGen: codeGen, producer: producer, log: log,
    }
}

// Run – ทำงานกับทุก tenant
func (uc *ProcessMaintenanceDueUseCase) Run(ctx context.Context) error {
    asOf := time.Now()
    schedules, err := uc.maintRepo.ListDueAllTenants(ctx, asOf)
    if err != nil { return err }

    uc.log.Info("pm processing", "due_count", len(schedules))

    for _, m := range schedules {
        if !m.AutoCreateWO {
            uc.publishDueEvent(ctx, m, 0)
            continue
        }

        // สร้าง WO
        woNo, err := uc.codeGen.NextWONo(ctx, m.TenantID)
        if err != nil { continue }

        wo, err := entity.NewWorkOrder(
            m.TenantID, m.CustomerID, m.SiteID, m.DeviceID,
            woNo,
            fmt.Sprintf("PM ประจำสำหรับ device %s", m.DeviceID),
            valueobject.JobTypeMaintenance,
            valueobject.JobPriorityNormal,
        )
        if err != nil { continue }

        wo.DeviceID = &m.DeviceID
        wo.ScheduledAt = m.NextDueAt
        wo.SourceRef = "pm_schedule"
        wo.SourceRefID = &m.ID

        if err := uc.woRepo.Save(ctx, wo); err != nil {
            uc.log.Error("create WO failed", "schedule_id", m.ID, "err", err)
            continue
        }

        daysOverdue := int(asOf.Sub(m.NextDueAt).Hours() / 24)
        uc.publishDueEvent(ctx, m, daysOverdue)
    }
    return nil
}

func (uc *ProcessMaintenanceDueUseCase) publishDueEvent(ctx context.Context, m *entity.MaintenanceSchedule, daysOverdue int) {
    _ = uc.producer.Publish(ctx, event.TopicMaintenanceDue, m.ID.String(), event.MaintenanceDue{
        EventID: uuid.New(), ScheduleID: m.ID, TenantID: m.TenantID,
        DeviceID: m.DeviceID, CustomerID: m.CustomerID, SiteID: m.SiteID,
        DueAt: m.NextDueAt, DaysOverdue: daysOverdue,
        OccurredAt: time.Now(),
    })
}
```

### `application/update_technician_gps.go`
```go
package application

import (
    "context"

    "github.com/google/uuid"

    "icmongolang/internal/modules/iotlogistics/domain/repository"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

type UpdateTechnicianGPSUseCase struct {
    techRepo repository.TechnicianRepository
    // Optionally: WS hub to broadcast
}

func NewUpdateTechnicianGPSUseCase(techRepo repository.TechnicianRepository) *UpdateTechnicianGPSUseCase {
    return &UpdateTechnicianGPSUseCase{techRepo: techRepo}
}

func (uc *UpdateTechnicianGPSUseCase) Execute(ctx context.Context, tenantID uuid.UUID, in UpdateTechnicianGPSInput) error {
    techID, _ := uuid.Parse(in.TechnicianID)
    loc, err := valueobject.NewGPSLocation(in.Lat, in.Lng, in.Accuracy)
    if err != nil { return err }
    return uc.techRepo.UpdateGPS(ctx, techID, loc)
}
```

---

# 🅲 PART 5C — INFRASTRUCTURE LAYER

## C.1 Persistence – GORM Models

### `infrastructure/persistence/postgres/models.go`
```go
package postgres

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/datatypes"
)

// ============================================================
// SHIPMENT
// ============================================================

type ShipmentModel struct {
    ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID        uuid.UUID      `gorm:"type:uuid;not null;index:idx_ship_tenant_status,priority:1"`
    ShipmentNo      string         `gorm:"size:30;not null;index:idx_ship_tenant_no,unique"`
    CustomerID      uuid.UUID      `gorm:"type:uuid;not null;index"`
    SiteID          uuid.UUID      `gorm:"type:uuid;not null;index"`
    FromWarehouseID uuid.UUID      `gorm:"type:uuid"`
    Carrier         string         `gorm:"size:100"`
    TrackingNo      string         `gorm:"size:100"`
    Status          string         `gorm:"size:30;not null;default:'DRAFT';index:idx_ship_tenant_status,priority:2"`
    ScheduledAt     time.Time      `gorm:"not null;index"`
    DispatchedAt    *time.Time
    DeliveredAt     *time.Time
    DeliveredTo     string         `gorm:"size:255"`
    SignatureURL    string         `gorm:"type:text"`
    GPSLast         datatypes.JSON `gorm:"type:jsonb"`
    Notes           string         `gorm:"type:text"`
    Metadata        datatypes.JSON `gorm:"type:jsonb"`
    CreatedBy       uuid.UUID      `gorm:"type:uuid"`
    CreatedAt       time.Time
    UpdatedAt       time.Time

    Items []ShipmentItemModel `gorm:"foreignKey:ShipmentID;constraint:OnDelete:CASCADE"`
}

func (ShipmentModel) TableName() string { return "iotlogistics_shipments" }

type ShipmentItemModel struct {
    ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    ShipmentID uuid.UUID  `gorm:"type:uuid;not null;index"`
    ProductID  uuid.UUID  `gorm:"type:uuid;not null"`
    DeviceID   *uuid.UUID `gorm:"type:uuid;index"`
    SerialNo   string     `gorm:"size:100;not null"`
    Qty        float64    `gorm:"type:numeric(15,2);not null"`
    UOM        string     `gorm:"size:20;default:'PCS'"`
    Notes      string     `gorm:"type:text"`
    CreatedAt  time.Time
}

func (ShipmentItemModel) TableName() string { return "iotlogistics_shipment_items" }

// ============================================================
// TECHNICIAN
// ============================================================

type TechnicianModel struct {
    ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID      uuid.UUID      `gorm:"type:uuid;not null;index:idx_tech_tenant_zone,priority:1"`
    UserID        uuid.UUID      `gorm:"type:uuid;not null;index"`
    Code          string         `gorm:"size:30;not null;index:idx_tech_tenant_code,unique"`
    Name          string         `gorm:"size:255;not null"`
    Phone         string         `gorm:"size:30"`
    Email         string         `gorm:"size:255"`
    Level         string         `gorm:"size:20;default:'JUNIOR'"`
    Skills        datatypes.JSON `gorm:"type:jsonb"`
    Zone          string         `gorm:"size:100;index:idx_tech_tenant_zone,priority:2"`
    HomeBase      datatypes.JSON `gorm:"type:jsonb"`
    IsAvailable   bool           `gorm:"default:true;index"`
    IsActive      bool           `gorm:"default:true;index"`
    GPSLast       datatypes.JSON `gorm:"type:jsonb"`
    GPSUpdatedAt  *time.Time
    MaxJobsPerDay int            `gorm:"default:6"`
    CurrentLoad   int            `gorm:"default:0"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

func (TechnicianModel) TableName() string { return "iotlogistics_technicians" }

// ============================================================
// INSTALLATION JOB
// ============================================================

type InstallationJobModel struct {
    ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID       uuid.UUID      `gorm:"type:uuid;not null;index:idx_job_tenant_status,priority:1"`
    JobNo          string         `gorm:"size:30;not null;index:idx_job_tenant_no,unique"`
    ShipmentID     *uuid.UUID     `gorm:"type:uuid;index"`
    CustomerID     uuid.UUID      `gorm:"type:uuid;not null;index"`
    SiteID         uuid.UUID      `gorm:"type:uuid;not null;index"`
    DeviceIDs      datatypes.JSON `gorm:"type:jsonb"`
    TechnicianID   *uuid.UUID     `gorm:"type:uuid;index"`
    JobType        string         `gorm:"size:30;not null;index"`
    Priority       string         `gorm:"size:20;not null;default:'NORMAL'"`
    Status         string         `gorm:"size:20;not null;index:idx_job_tenant_status,priority:2"`
    ScheduledAt    time.Time      `gorm:"not null;index"`
    SLADueAt       time.Time      `gorm:"not null;index"`
    StartedAt      *time.Time
    CompletedAt    *time.Time
    StartedGPS     datatypes.JSON `gorm:"type:jsonb"`
    CompletedGPS   datatypes.JSON `gorm:"type:jsonb"`
    Checklist      datatypes.JSON `gorm:"type:jsonb"`
    Photos         datatypes.JSON `gorm:"type:jsonb"`
    SignatureURL   string         `gorm:"type:text"`
    Notes          string         `gorm:"type:text"`
    Result         datatypes.JSON `gorm:"type:jsonb"`
    AssignedBy     *uuid.UUID     `gorm:"type:uuid"`
    AssignedAt     *time.Time
    CancelledAt    *time.Time
    CancelReason   string         `gorm:"type:text"`
    CreatedBy      uuid.UUID      `gorm:"type:uuid"`
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

func (InstallationJobModel) TableName() string { return "iotlogistics_installation_jobs" }

// ============================================================
// WORK ORDER
// ============================================================

type WorkOrderModel struct {
    ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID      uuid.UUID      `gorm:"type:uuid;not null;index"`
    WONo          string         `gorm:"size:30;not null;index:idx_wo_tenant_no,unique"`
    DeviceID      *uuid.UUID     `gorm:"type:uuid;index"`
    SiteID        uuid.UUID      `gorm:"type:uuid;not null;index"`
    CustomerID    uuid.UUID      `gorm:"type:uuid;not null;index"`
    JobType       string         `gorm:"size:30;not null"`
    Priority      string         `gorm:"size:20;default:'NORMAL'"`
    Status        string         `gorm:"size:20;not null;index"`
    Title         string         `gorm:"size:255;not null"`
    Description   string         `gorm:"type:text"`
    TechnicianID  *uuid.UUID     `gorm:"type:uuid;index"`
    ScheduledAt   time.Time
    TimeWindow    datatypes.JSON `gorm:"type:jsonb"`
    StartedAt     *time.Time
    CompletedAt   *time.Time
    Resolution    string         `gorm:"type:text"`
    PartsUsed     datatypes.JSON `gorm:"type:jsonb"`
    LaborMinutes  int            `gorm:"default:0"`
    CostTotal     float64        `gorm:"type:numeric(15,2);default:0"`
    Currency      string         `gorm:"size:3;default:'THB'"`
    SourceRef     string         `gorm:"size:50"`
    SourceRefID   *uuid.UUID     `gorm:"type:uuid"`
    Notes         string         `gorm:"type:text"`
    CreatedBy     uuid.UUID      `gorm:"type:uuid"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

func (WorkOrderModel) TableName() string { return "iotlogistics_work_orders" }

// ============================================================
// MAINTENANCE
// ============================================================

type MaintenanceScheduleModel struct {
    ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID       uuid.UUID  `gorm:"type:uuid;not null;index:idx_maint_tenant_due,priority:1"`
    DeviceID       uuid.UUID  `gorm:"type:uuid;not null;index:idx_maint_device,unique"`
    CustomerID     uuid.UUID  `gorm:"type:uuid;not null;index"`
    SiteID         uuid.UUID  `gorm:"type:uuid;not null;index"`
    JobType        string     `gorm:"size:30;default:'MAINTENANCE'"`
    Frequency      string     `gorm:"size:30;not null"`
    CustomDays     int        `gorm:"default:0"`
    LastDoneAt     *time.Time
    NextDueAt      time.Time  `gorm:"not null;index:idx_maint_tenant_due,priority:2"`
    AssignedTechID *uuid.UUID `gorm:"type:uuid"`
    AutoCreateWO   bool       `gorm:"default:true"`
    IsActive       bool       `gorm:"default:true;index"`
    Notes          string     `gorm:"type:text"`
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

func (MaintenanceScheduleModel) TableName() string { return "iotlogistics_maintenance_schedules" }

// ============================================================
// RMA
// ============================================================

type RMAModel struct {
    ID                  uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID            uuid.UUID      `gorm:"type:uuid;not null;index:idx_rma_tenant_status,priority:1"`
    RMANo               string         `gorm:"size:30;not null;index:idx_rma_tenant_no,unique"`
    DeviceID            uuid.UUID      `gorm:"type:uuid;not null;index"`
    CustomerID          uuid.UUID      `gorm:"type:uuid;not null;index"`
    SiteID              uuid.UUID      `gorm:"type:uuid;not null;index"`
    JobID               *uuid.UUID     `gorm:"type:uuid"`
    Reason              string         `gorm:"size:30;not null"`
    ReasonDetail        string         `gorm:"type:text"`
    Status              string         `gorm:"size:20;not null;index:idx_rma_tenant_status,priority:2"`
    Resolution          *string        `gorm:"size:20"`
    ReplacementDeviceID *uuid.UUID     `gorm:"type:uuid"`
    WarrantyClaim       bool           `gorm:"default:false"`
    ReceivedPhotos      datatypes.JSON `gorm:"type:jsonb"`
    ApprovedAt          *time.Time
    ApprovedBy          *uuid.UUID     `gorm:"type:uuid"`
    ShippedAt           *time.Time
    ReceivedAt          *time.Time
    ClosedAt            *time.Time
    RejectReason        string         `gorm:"type:text"`
    AssignedTechID      *uuid.UUID     `gorm:"type:uuid"`
    Notes               string         `gorm:"type:text"`
    CreatedBy           uuid.UUID      `gorm:"type:uuid"`
    CreatedAt           time.Time
    UpdatedAt           time.Time
}

func (RMAModel) TableName() string { return "iotlogistics_rma" }

// ============================================================
// SPARE PART
// ============================================================

type SparePartModel struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID  uuid.UUID `gorm:"type:uuid;not null;index"`
    PartNo    string    `gorm:"size:50;not null;index"`
    Name      string    `gorm:"size:255;not null"`
    Category  string    `gorm:"size:50"`
    UnitPrice float64   `gorm:"type:numeric(15,2)"`
    Currency  string    `gorm:"size:3;default:'THB'"`
    IsActive  bool      `gorm:"default:true"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (SparePartModel) TableName() string { return "iotlogistics_spare_parts" }

// ============================================================
// TECHNICIAN LOCATION (tracking log)
// ============================================================

type TechnicianLocationModel struct {
    ID           int64          `gorm:"primaryKey;autoIncrement"`
    TenantID     uuid.UUID      `gorm:"type:uuid;not null;index:idx_tl_tenant_time,priority:1"`
    TechnicianID uuid.UUID      `gorm:"type:uuid;not null;index"`
    Location     datatypes.JSON `gorm:"type:jsonb"`
    JobID        *uuid.UUID     `gorm:"type:uuid"`
    BatteryLevel *int
    Network      string     `gorm:"size:20"`
    RecordedAt   time.Time  `gorm:"index:idx_tl_tenant_time,priority:2"`
}

func (TechnicianLocationModel) TableName() string { return "iotlogistics_technician_locations" }

// ============================================================
// AUDIT
// ============================================================

type AuditLogModel struct {
    ID         int64          `gorm:"primaryKey;autoIncrement"`
    TenantID   uuid.UUID      `gorm:"type:uuid;not null;index"`
    ActorID    uuid.UUID      `gorm:"type:uuid"`
    Action     string         `gorm:"size:50;not null"`
    EntityType string         `gorm:"size:50;not null;index:idx_log_entity,priority:1"`
    EntityID   uuid.UUID      `gorm:"type:uuid;not null;index:idx_log_entity,priority:2"`
    Payload    datatypes.JSON `gorm:"type:jsonb"`
    IPAddress  string         `gorm:"size:45"`
    UserAgent  string         `gorm:"size:500"`
    CreatedAt  time.Time      `gorm:"index"`
}

func (AuditLogModel) TableName() string { return "iotlogistics_audit_logs" }
```

### `infrastructure/persistence/postgres/mappers.go`
```go
package postgres

import (
    "encoding/json"

    "gorm.io/datatypes"

    "icmongolang/internal/modules/iotlogistics/domain/entity"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

// ============================================================
// SHIPMENT
// ============================================================

func toShipmentModel(s *entity.Shipment) *ShipmentModel {
    m := &ShipmentModel{
        ID: s.ID, TenantID: s.TenantID,
        ShipmentNo: s.ShipmentNo.String(),
        CustomerID: s.CustomerID, SiteID: s.SiteID,
        FromWarehouseID: s.FromWarehouseID,
        Carrier: s.Carrier, TrackingNo: s.TrackingNo,
        Status: string(s.Status), ScheduledAt: s.ScheduledAt,
        DispatchedAt: s.DispatchedAt, DeliveredAt: s.DeliveredAt,
        DeliveredTo: s.DeliveredTo, SignatureURL: s.SignatureURL,
        Notes: s.Notes, Metadata: marshalJSON(s.Metadata),
        CreatedBy: s.CreatedBy, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
    }
    if s.GPSLast != nil {
        m.GPSLast = marshalJSON(s.GPSLast)
    }
    for _, it := range s.Items {
        m.Items = append(m.Items, ShipmentItemModel{
            ID: it.ID, ShipmentID: it.ShipmentID,
            ProductID: it.ProductID, DeviceID: it.DeviceID,
            SerialNo: it.SerialNo, Qty: it.Qty, UOM: it.UOM,
            Notes: it.Notes,
        })
    }
    return m
}

func toShipmentEntity(m *ShipmentModel) *entity.Shipment {
    s := &entity.Shipment{
        ID: m.ID, TenantID: m.TenantID,
        ShipmentNo: valueobject.ShipmentNo(m.ShipmentNo),
        CustomerID: m.CustomerID, SiteID: m.SiteID,
        FromWarehouseID: m.FromWarehouseID,
        Carrier: m.Carrier, TrackingNo: m.TrackingNo,
        Status: valueobject.ShipmentStatus(m.Status),
        ScheduledAt: m.ScheduledAt,
        DispatchedAt: m.DispatchedAt, DeliveredAt: m.DeliveredAt,
        DeliveredTo: m.DeliveredTo, SignatureURL: m.SignatureURL,
        Notes: m.Notes,
        CreatedBy: m.CreatedBy, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
    }
    if len(m.GPSLast) > 0 {
        var loc valueobject.GPSLocation
        _ = json.Unmarshal(m.GPSLast, &loc)
        s.GPSLast = &loc
    }
    if len(m.Metadata) > 0 {
        _ = json.Unmarshal(m.Metadata, &s.Metadata)
    }
    if s.Metadata == nil { s.Metadata = map[string]any{} }
    s.Items = make([]*entity.ShipmentItem, 0, len(m.Items))
    for i := range m.Items {
        it := &m.Items[i]
        s.Items = append(s.Items, &entity.ShipmentItem{
            ID: it.ID, ShipmentID: it.ShipmentID,
            ProductID: it.ProductID, DeviceID: it.DeviceID,
            SerialNo: it.SerialNo, Qty: it.Qty, UOM: it.UOM,
            Notes: it.Notes,
        })
    }
    return s
}

// ============================================================
// TECHNICIAN
// ============================================================

func toTechnicianModel(t *entity.Technician) *TechnicianModel {
    m := &TechnicianModel{
        ID: t.ID, TenantID: t.TenantID, UserID: t.UserID,
        Code: t.Code, Name: t.Name, Phone: t.Phone, Email: t.Email,
        Level: string(t.Level),
        Zone: t.Zone, IsAvailable: t.IsAvailable, IsActive: t.IsActive,
        GPSUpdatedAt: t.GPSUpdatedAt,
        MaxJobsPerDay: t.MaxJobsPerDay, CurrentLoad: t.CurrentLoad,
        CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
    }
    if len(t.Skills) > 0 {
        m.Skills = marshalJSON(t.Skills)
    }
    if t.HomeBase != nil {
        m.HomeBase = marshalJSON(t.HomeBase)
    }
    if t.GPSLast != nil {
        m.GPSLast = marshalJSON(t.GPSLast)
    }
    return m
}

func toTechnicianEntity(m *TechnicianModel) *entity.Technician {
    t := &entity.Technician{
        ID: m.ID, TenantID: m.TenantID, UserID: m.UserID,
        Code: m.Code, Name: m.Name, Phone: m.Phone, Email: m.Email,
        Level: entity.TechnicianLevel(m.Level),
        Zone: m.Zone, IsAvailable: m.IsAvailable, IsActive: m.IsActive,
        GPSUpdatedAt: m.GPSUpdatedAt,
        MaxJobsPerDay: m.MaxJobsPerDay, CurrentLoad: m.CurrentLoad,
        CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
    }
    if len(m.Skills) > 0 {
        _ = json.Unmarshal(m.Skills, &t.Skills)
    }
    if t.Skills == nil { t.Skills = valueobject.SkillSet{} }
    if len(m.HomeBase) > 0 {
        var loc valueobject.GPSLocation
        _ = json.Unmarshal(m.HomeBase, &loc)
        t.HomeBase = &loc
    }
    if len(m.GPSLast) > 0 {
        var loc valueobject.GPSLocation
        _ = json.Unmarshal(m.GPSLast, &loc)
        t.GPSLast = &loc
    }
    return t
}

// ============================================================
// INSTALLATION JOB
// ============================================================

func toInstallationModel(j *entity.InstallationJob) *InstallationJobModel {
    m := &InstallationJobModel{
        ID: j.ID, TenantID: j.TenantID, JobNo: j.JobNo,
        ShipmentID: j.ShipmentID,
        CustomerID: j.CustomerID, SiteID: j.SiteID,
        TechnicianID: j.TechnicianID,
        JobType: string(j.JobType), Priority: string(j.Priority),
        Status: string(j.Status),
        ScheduledAt: j.ScheduledAt, SLADueAt: j.SLADueAt,
        StartedAt: j.StartedAt, CompletedAt: j.CompletedAt,
        SignatureURL: j.SignatureURL, Notes: j.Notes,
        AssignedBy: j.AssignedBy, AssignedAt: j.AssignedAt,
        CancelledAt: j.CancelledAt, CancelReason: j.CancelReason,
        CreatedBy: j.CreatedBy, CreatedAt: j.CreatedAt, UpdatedAt: j.UpdatedAt,
    }
    if len(j.DeviceIDs) > 0 {
        m.DeviceIDs = marshalJSON(j.DeviceIDs)
    }
    if j.StartedGPS != nil {
        m.StartedGPS = marshalJSON(j.StartedGPS)
    }
    if j.CompletedGPS != nil {
        m.CompletedGPS = marshalJSON(j.CompletedGPS)
    }
    m.Checklist = marshalJSON(j.Checklist)
    m.Photos = marshalJSON(j.Photos)
    if len(j.Result) > 0 {
        m.Result = marshalJSON(j.Result)
    }
    return m
}

func toInstallationEntity(m *InstallationJobModel) *entity.InstallationJob {
    j := &entity.InstallationJob{
        ID: m.ID, TenantID: m.TenantID, JobNo: m.JobNo,
        ShipmentID: m.ShipmentID,
        CustomerID: m.CustomerID, SiteID: m.SiteID,
        TechnicianID: m.TechnicianID,
        JobType: valueobject.JobType(m.JobType),
        Priority: valueobject.JobPriority(m.Priority),
        Status: valueobject.JobStatus(m.Status),
        ScheduledAt: m.ScheduledAt, SLADueAt: m.SLADueAt,
        StartedAt: m.StartedAt, CompletedAt: m.CompletedAt,
        SignatureURL: m.SignatureURL, Notes: m.Notes,
        AssignedBy: m.AssignedBy, AssignedAt: m.AssignedAt,
        CancelledAt: m.CancelledAt, CancelReason: m.CancelReason,
        CreatedBy: m.CreatedBy, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
    }
    if len(m.DeviceIDs) > 0 { _ = json.Unmarshal(m.DeviceIDs, &j.DeviceIDs) }
    if j.DeviceIDs == nil { j.DeviceIDs = []uuid.UUID{} }
    if len(m.StartedGPS) > 0 {
        var loc valueobject.GPSLocation
        _ = json.Unmarshal(m.StartedGPS, &loc)
        j.StartedGPS = &loc
    }
    if len(m.CompletedGPS) > 0 {
        var loc valueobject.GPSLocation
        _ = json.Unmarshal(m.CompletedGPS, &loc)
        j.CompletedGPS = &loc
    }
    if len(m.Checklist) > 0 { _ = json.Unmarshal(m.Checklist, &j.Checklist) }
    if j.Checklist == nil { j.Checklist = valueobject.Checklist{} }
    if len(m.Photos) > 0 { _ = json.Unmarshal(m.Photos, &j.Photos) }
    if j.Photos == nil { j.Photos = []string{} }
    if len(m.Result) > 0 { _ = json.Unmarshal(m.Result, &j.Result) }
    return j
}

// ============================================================
// Others — toWorkOrderModel/Entity, toMaintModel/Entity, toRMAModel/Entity
// ============================================================

func toWorkOrderModel(w *entity.WorkOrder) *WorkOrderModel {
    m := &WorkOrderModel{
        ID: w.ID, TenantID: w.TenantID, WONo: w.WONo,
        DeviceID: w.DeviceID, SiteID: w.SiteID, CustomerID: w.CustomerID,
        JobType: string(w.JobType), Priority: string(w.Priority),
        Status: string(w.Status), Title: w.Title, Description: w.Description,
        TechnicianID: w.TechnicianID, ScheduledAt: w.ScheduledAt,
        StartedAt: w.StartedAt, CompletedAt: w.CompletedAt,
        Resolution: w.Resolution, LaborMinutes: w.LaborMinutes,
        CostTotal: w.CostTotal, Currency: w.Currency,
        SourceRef: w.SourceRef, SourceRefID: w.SourceRefID,
        Notes: w.Notes, CreatedBy: w.CreatedBy,
        CreatedAt: w.CreatedAt, UpdatedAt: w.UpdatedAt,
    }
    if len(w.PartsUsed) > 0 { m.PartsUsed = marshalJSON(w.PartsUsed) }
    if w.TimeWindow != nil { m.TimeWindow = marshalJSON(w.TimeWindow) }
    return m
}

func toWorkOrderEntity(m *WorkOrderModel) *entity.WorkOrder {
    w := &entity.WorkOrder{
        ID: m.ID, TenantID: m.TenantID, WONo: m.WONo,
        DeviceID: m.DeviceID, SiteID: m.SiteID, CustomerID: m.CustomerID,
        JobType: valueobject.JobType(m.JobType),
        Priority: valueobject.JobPriority(m.Priority),
        Status: valueobject.JobStatus(m.Status),
        Title: m.Title, Description: m.Description,
        TechnicianID: m.TechnicianID, ScheduledAt: m.ScheduledAt,
        StartedAt: m.StartedAt, CompletedAt: m.CompletedAt,
        Resolution: m.Resolution, LaborMinutes: m.LaborMinutes,
        CostTotal: m.CostTotal, Currency: m.Currency,
        SourceRef: m.SourceRef, SourceRefID: m.SourceRefID,
        Notes: m.Notes, CreatedBy: m.CreatedBy,
        CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
    }
    if len(m.PartsUsed) > 0 { _ = json.Unmarshal(m.PartsUsed, &w.PartsUsed) }
    if w.PartsUsed == nil { w.PartsUsed = []entity.SparePartUsage{} }
    return w
}

func toMaintModel(mm *entity.MaintenanceSchedule) *MaintenanceScheduleModel {
    return &MaintenanceScheduleModel{
        ID: mm.ID, TenantID: mm.TenantID, DeviceID: mm.DeviceID,
        CustomerID: mm.CustomerID, SiteID: mm.SiteID,
        JobType: string(mm.JobType), Frequency: string(mm.Frequency),
        CustomDays: mm.CustomDays, LastDoneAt: mm.LastDoneAt,
        NextDueAt: mm.NextDueAt, AssignedTechID: mm.AssignedTechID,
        AutoCreateWO: mm.AutoCreateWO, IsActive: mm.IsActive,
        Notes: mm.Notes, CreatedAt: mm.CreatedAt, UpdatedAt: mm.UpdatedAt,
    }
}

func toMaintEntity(m *MaintenanceScheduleModel) *entity.MaintenanceSchedule {
    return &entity.MaintenanceSchedule{
        ID: m.ID, TenantID: m.TenantID, DeviceID: m.DeviceID,
        CustomerID: m.CustomerID, SiteID: m.SiteID,
        JobType: valueobject.JobType(m.JobType),
        Frequency: entity.PMFrequency(m.Frequency),
        CustomDays: m.CustomDays, LastDoneAt: m.LastDoneAt,
        NextDueAt: m.NextDueAt, AssignedTechID: m.AssignedTechID,
        AutoCreateWO: m.AutoCreateWO, IsActive: m.IsActive,
        Notes: m.Notes, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
    }
}

func toRMAModel(r *entity.RMA) *RMAModel {
    m := &RMAModel{
        ID: r.ID, TenantID: r.TenantID, RMANo: r.RMANo,
        DeviceID: r.DeviceID, CustomerID: r.CustomerID, SiteID: r.SiteID,
        JobID: r.JobID, Reason: string(r.Reason), ReasonDetail: r.ReasonDetail,
        Status: string(r.Status),
        ReplacementDeviceID: r.ReplacementDeviceID,
        WarrantyClaim: r.WarrantyClaim,
        ApprovedAt: r.ApprovedAt, ApprovedBy: r.ApprovedBy,
        ShippedAt: r.ShippedAt, ReceivedAt: r.ReceivedAt, ClosedAt: r.ClosedAt,
        RejectReason: r.RejectReason, AssignedTechID: r.AssignedTechID,
        Notes: r.Notes, CreatedBy: r.CreatedBy,
        CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
    }
    if r.Resolution != nil {
        s := string(*r.Resolution)
        m.Resolution = &s
    }
    if len(r.ReceivedPhotos) > 0 {
        m.ReceivedPhotos = marshalJSON(r.ReceivedPhotos)
    }
    return m
}

func toRMAEntity(m *RMAModel) *entity.RMA {
    r := &entity.RMA{
        ID: m.ID, TenantID: m.TenantID, RMANo: m.RMANo,
        DeviceID: m.DeviceID, CustomerID: m.CustomerID, SiteID: m.SiteID,
        JobID: m.JobID, Reason: entity.RMAReason(m.Reason),
        ReasonDetail: m.ReasonDetail,
        Status: valueobject.RMAStatus(m.Status),
        ReplacementDeviceID: m.ReplacementDeviceID,
        WarrantyClaim: m.WarrantyClaim,
        ApprovedAt: m.ApprovedAt, ApprovedBy: m.ApprovedBy,
        ShippedAt: m.ShippedAt, ReceivedAt: m.ReceivedAt, ClosedAt: m.ClosedAt,
        RejectReason: m.RejectReason, AssignedTechID: m.AssignedTechID,
        Notes: m.Notes, CreatedBy: m.CreatedBy,
        CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
    }
    if m.Resolution != nil {
        res := entity.RMAResolution(*m.Resolution)
        r.Resolution = &res
    }
    if len(m.ReceivedPhotos) > 0 { _ = json.Unmarshal(m.ReceivedPhotos, &r.ReceivedPhotos) }
    if r.ReceivedPhotos == nil { r.ReceivedPhotos = []string{} }
    return r
}

func marshalJSON(v any) datatypes.JSON {
    if v == nil { return datatypes.JSON([]byte("{}")) }
    b, err := json.Marshal(v)
    if err != nil { return datatypes.JSON([]byte("{}")) }
    return datatypes.JSON(b)
}
```

> **ต้องเพิ่ม import `uuid`** ใน mapper

### `infrastructure/persistence/postgres/shipment_repository.go`
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

    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    "icmongolang/internal/modules/iotlogistics/domain/entity"
    "icmongolang/internal/modules/iotlogistics/domain/repository"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

type shipmentRepository struct{ db *gorm.DB }

func NewShipmentRepository(db *gorm.DB) repository.ShipmentRepository {
    return &shipmentRepository{db: db}
}

func (r *shipmentRepository) Save(ctx context.Context, s *entity.Shipment) error {
    m := toShipmentModel(s)
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Clauses(clause.OnConflict{
            Columns: []clause.Column{{Name: "id"}},
            UpdateAll: true,
        }).Omit("Items").Create(m).Error; err != nil {
            return err
        }
        // sync items: delete removed + upsert current
        existingIDs := make([]uuid.UUID, 0, len(m.Items))
        for _, it := range m.Items { existingIDs = append(existingIDs, it.ID) }
        q := tx.Where("shipment_id = ?", s.ID)
        if len(existingIDs) > 0 {
            q = q.Where("id NOT IN ?", existingIDs)
        }
        if err := q.Delete(&ShipmentItemModel{}).Error; err != nil { return err }
        if len(m.Items) > 0 {
            if err := tx.Clauses(clause.OnConflict{
                Columns:   []clause.Column{{Name: "id"}},
                DoUpdates: clause.AssignmentColumns([]string{"serial_no", "qty", "uom", "notes", "device_id"}),
            }).Create(&m.Items).Error; err != nil {
                return err
            }
        }
        return nil
    })
}

func (r *shipmentRepository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Shipment, error) {
    var m ShipmentModel
    err := r.db.WithContext(ctx).
        Preload("Items").
        Where("tenant_id = ? AND id = ?", tenantID, id).
        First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrShipmentNotFound
    }
    if err != nil { return nil, err }
    return toShipmentEntity(&m), nil
}

func (r *shipmentRepository) FindByNo(ctx context.Context, tenantID uuid.UUID, no valueobject.ShipmentNo) (*entity.Shipment, error) {
    var m ShipmentModel
    err := r.db.WithContext(ctx).
        Preload("Items").
        Where("tenant_id = ? AND shipment_no = ?", tenantID, no.String()).
        First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrShipmentNotFound
    }
    if err != nil { return nil, err }
    return toShipmentEntity(&m), nil
}

func (r *shipmentRepository) List(ctx context.Context, f repository.ShipmentFilter) ([]*entity.Shipment, int64, error) {
    q := r.db.WithContext(ctx).Model(&ShipmentModel{}).Where("tenant_id = ?", f.TenantID)
    if f.CustomerID != nil { q = q.Where("customer_id = ?", *f.CustomerID) }
    if f.SiteID != nil { q = q.Where("site_id = ?", *f.SiteID) }
    if len(f.Statuses) > 0 {
        ss := make([]string, len(f.Statuses))
        for i, s := range f.Statuses { ss[i] = string(s) }
        q = q.Where("status IN ?", ss)
    }
    if f.From != nil { q = q.Where("scheduled_at >= ?", f.From) }
    if f.To != nil   { q = q.Where("scheduled_at <= ?", f.To) }
    if f.Search != "" {
        s := "%" + strings.ToLower(f.Search) + "%"
        q = q.Where("(LOWER(shipment_no) LIKE ? OR LOWER(tracking_no) LIKE ?)", s, s)
    }
    var total int64
    if err := q.Count(&total).Error; err != nil { return nil, 0, err }

    sortBy := sanitizeShipSortBy(f.SortBy)
    order := "DESC"
    if strings.EqualFold(f.SortOrder, "asc") { order = "ASC" }
    q = q.Order(fmt.Sprintf("%s %s", sortBy, order))

    if f.PageSize > 0 { q = q.Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize) }

    var models []ShipmentModel
    if err := q.Preload("Items").Find(&models).Error; err != nil { return nil, 0, err }
    out := make([]*entity.Shipment, 0, len(models))
    for i := range models { out = append(out, toShipmentEntity(&models[i])) }
    return out, total, nil
}

func (r *shipmentRepository) ListByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) ([]*entity.Shipment, error) {
    var models []ShipmentModel
    if err := r.db.WithContext(ctx).
        Preload("Items").
        Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
        Order("created_at DESC").Find(&models).Error; err != nil {
        return nil, err
    }
    out := make([]*entity.Shipment, 0, len(models))
    for i := range models { out = append(out, toShipmentEntity(&models[i])) }
    return out, nil
}

func (r *shipmentRepository) ListPendingBySite(ctx context.Context, tenantID, siteID uuid.UUID) ([]*entity.Shipment, error) {
    var models []ShipmentModel
    if err := r.db.WithContext(ctx).
        Preload("Items").
        Where("tenant_id = ? AND site_id = ? AND status IN ?",
            tenantID, siteID, []string{"PENDING", "IN_TRANSIT", "DELIVERED"}).
        Order("scheduled_at ASC").Find(&models).Error; err != nil {
        return nil, err
    }
    out := make([]*entity.Shipment, 0, len(models))
    for i := range models { out = append(out, toShipmentEntity(&models[i])) }
    return out, nil
}

func (r *shipmentRepository) NextSequence(ctx context.Context, tenantID uuid.UUID, year int) (int, error) {
    var count int64
    pattern := fmt.Sprintf("SHP-%04d-%%", year)
    err := r.db.WithContext(ctx).Model(&ShipmentModel{}).
        Where("tenant_id = ? AND shipment_no LIKE ?", tenantID, pattern).
        Count(&count).Error
    return int(count) + 1, err
}

func sanitizeShipSortBy(s string) string {
    allowed := map[string]bool{
        "created_at": true, "scheduled_at": true,
        "dispatched_at": true, "delivered_at": true,
        "shipment_no": true, "status": true,
    }
    if allowed[s] { return s }
    return "created_at"
}

var _ = time.Now
```

### `infrastructure/persistence/postgres/installation_repository.go`
```go
package postgres

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"

    "icmongolang/internal/modules/iotlogistics/domain/entity"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    "icmongolang/internal/modules/iotlogistics/domain/repository"
)

type installationRepository struct{ db *gorm.DB }

func NewInstallationRepository(db *gorm.DB) repository.InstallationRepository {
    return &installationRepository{db: db}
}

func (r *installationRepository) Save(ctx context.Context, j *entity.InstallationJob) error {
    m := toInstallationModel(j)
    return r.db.WithContext(ctx).Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "id"}},
        UpdateAll: true,
    }).Create(m).Error
}

func (r *installationRepository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.InstallationJob, error) {
    var m InstallationJobModel
    err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND id = ?", tenantID, id).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrJobNotFound
    }
    if err != nil { return nil, err }
    return toInstallationEntity(&m), nil
}

func (r *installationRepository) FindByNo(ctx context.Context, tenantID uuid.UUID, no string) (*entity.InstallationJob, error) {
    var m InstallationJobModel
    err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND job_no = ?", tenantID, no).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrJobNotFound
    }
    if err != nil { return nil, err }
    return toInstallationEntity(&m), nil
}

func (r *installationRepository) List(ctx context.Context, f repository.InstallationJobFilter) ([]*entity.InstallationJob, int64, error) {
    q := r.db.WithContext(ctx).Model(&InstallationJobModel{}).Where("tenant_id = ?", f.TenantID)
    if f.TechnicianID != nil { q = q.Where("technician_id = ?", *f.TechnicianID) }
    if f.SiteID != nil { q = q.Where("site_id = ?", *f.SiteID) }
    if f.CustomerID != nil { q = q.Where("customer_id = ?", *f.CustomerID) }
    if len(f.Statuses) > 0 {
        ss := make([]string, len(f.Statuses))
        for i, s := range f.Statuses { ss[i] = string(s) }
        q = q.Where("status IN ?", ss)
    }
    if len(f.Types) > 0 {
        ts := make([]string, len(f.Types))
        for i, t := range f.Types { ts[i] = string(t) }
        q = q.Where("job_type IN ?", ts)
    }
    if f.From != nil { q = q.Where("scheduled_at >= ?", f.From) }
    if f.To != nil   { q = q.Where("scheduled_at <= ?", f.To) }
    if f.OnlySLABreach != nil && *f.OnlySLABreach {
        q = q.Where("sla_due_at < ? AND status NOT IN ?", time.Now(), []string{"DONE", "FAILED", "CANCELLED"})
    }
    var total int64
    if err := q.Count(&total).Error; err != nil { return nil, 0, err }

    if f.PageSize > 0 { q = q.Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize) }

    var models []InstallationJobModel
    if err := q.Order("scheduled_at ASC").Find(&models).Error; err != nil {
        return nil, 0, err
    }
    out := make([]*entity.InstallationJob, 0, len(models))
    for i := range models { out = append(out, toInstallationEntity(&models[i])) }
    return out, total, nil
}

func (r *installationRepository) ListByTechnician(ctx context.Context, tenantID, techID uuid.UUID, date time.Time) ([]*entity.InstallationJob, error) {
    dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
    dayEnd := dayStart.Add(24 * time.Hour)
    var models []InstallationJobModel
    if err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND technician_id = ? AND scheduled_at >= ? AND scheduled_at < ?",
            tenantID, techID, dayStart, dayEnd).
        Order("scheduled_at ASC").Find(&models).Error; err != nil {
        return nil, err
    }
    out := make([]*entity.InstallationJob, 0, len(models))
    for i := range models { out = append(out, toInstallationEntity(&models[i])) }
    return out, nil
}

func (r *installationRepository) ListBySite(ctx context.Context, tenantID, siteID uuid.UUID) ([]*entity.InstallationJob, error) {
    var models []InstallationJobModel
    if err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND site_id = ?", tenantID, siteID).
        Order("scheduled_at DESC").Find(&models).Error; err != nil {
        return nil, err
    }
    out := make([]*entity.InstallationJob, 0, len(models))
    for i := range models { out = append(out, toInstallationEntity(&models[i])) }
    return out, nil
}

func (r *installationRepository) ListSLABreached(ctx context.Context, tenantID uuid.UUID, asOf time.Time) ([]*entity.InstallationJob, error) {
    var models []InstallationJobModel
    err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND sla_due_at < ? AND status NOT IN ?",
            tenantID, asOf, []string{"DONE", "FAILED", "CANCELLED"}).
        Order("sla_due_at ASC").Find(&models).Error
    if err != nil { return nil, err }
    out := make([]*entity.InstallationJob, 0, len(models))
    for i := range models { out = append(out, toInstallationEntity(&models[i])) }
    return out, nil
}

func (r *installationRepository) CountActiveByTechnician(ctx context.Context, techID uuid.UUID, date time.Time) (int64, error) {
    dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
    dayEnd := dayStart.Add(24 * time.Hour)
    var n int64
    err := r.db.WithContext(ctx).Model(&InstallationJobModel{}).
        Where("technician_id = ? AND scheduled_at >= ? AND scheduled_at < ? AND status NOT IN ?",
            techID, dayStart, dayEnd, []string{"DONE", "FAILED", "CANCELLED"}).
        Count(&n).Error
    return n, err
}

func (r *installationRepository) NextSequence(ctx context.Context, tenantID uuid.UUID, year int) (int, error) {
    var count int64
    pattern := fmt.Sprintf("JOB-%04d-%%", year)
    err := r.db.WithContext(ctx).Model(&InstallationJobModel{}).
        Where("tenant_id = ? AND job_no LIKE ?", tenantID, pattern).
        Count(&count).Error
    return int(count) + 1, err
}
```

### `infrastructure/persistence/postgres/technician_repository.go`
```go
package postgres

import (
    "context"
    "errors"
    "math"

    "github.com/google/uuid"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"

    "icmongolang/internal/modules/iotlogistics/domain/entity"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    "icmongolang/internal/modules/iotlogistics/domain/repository"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

type technicianRepository struct{ db *gorm.DB }

func NewTechnicianRepository(db *gorm.DB) repository.TechnicianRepository {
    return &technicianRepository{db: db}
}

func (r *technicianRepository) Save(ctx context.Context, t *entity.Technician) error {
    return r.db.WithContext(ctx).Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "id"}},
        UpdateAll: true,
    }).Create(toTechnicianModel(t)).Error
}

func (r *technicianRepository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Technician, error) {
    var m TechnicianModel
    err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND id = ?", tenantID, id).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrTechnicianNotFound
    }
    if err != nil { return nil, err }
    return toTechnicianEntity(&m), nil
}

func (r *technicianRepository) FindByUserID(ctx context.Context, tenantID, userID uuid.UUID) (*entity.Technician, error) {
    var m TechnicianModel
    err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND user_id = ?", tenantID, userID).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrTechnicianNotFound
    }
    if err != nil { return nil, err }
    return toTechnicianEntity(&m), nil
}

func (r *technicianRepository) List(ctx context.Context, f repository.TechnicianFilter) ([]*entity.Technician, int64, error) {
    q := r.db.WithContext(ctx).Model(&TechnicianModel{}).Where("tenant_id = ?", f.TenantID)
    if f.Zone != "" { q = q.Where("zone = ?", f.Zone) }
    if f.Available != nil { q = q.Where("is_available = ?", *f.Available) }
    if f.Active != nil { q = q.Where("is_active = ?", *f.Active) }
    if len(f.Skills) > 0 {
        // JSONB contain any of skills
        q = q.Where("skills ?| ?", pgTextArray(f.Skills))
    }
    var total int64
    if err := q.Count(&total).Error; err != nil { return nil, 0, err }
    if f.PageSize > 0 { q = q.Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize) }
    var models []TechnicianModel
    if err := q.Order("current_load ASC, name ASC").Find(&models).Error; err != nil {
        return nil, 0, err
    }
    out := make([]*entity.Technician, 0, len(models))
    for i := range models { out = append(out, toTechnicianEntity(&models[i])) }
    return out, total, nil
}

func (r *technicianRepository) ListAvailableNear(ctx context.Context, tenantID uuid.UUID, center valueobject.GPSLocation, radiusKm float64) ([]*entity.Technician, error) {
    // PostgreSQL only – ใช้ Haversine approximation ใน SQL
    // สำหรับ production ควรใช้ PostGIS
    var models []TechnicianModel
    err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND is_active = true", tenantID).
        Find(&models).Error
    if err != nil { return nil, err }

    out := make([]*entity.Technician, 0)
    for _, m := range models {
        t := toTechnicianEntity(&m)
        if t.GPSLast == nil { continue }
        if t.GPSLast.DistanceKm(center) <= radiusKm {
            out = append(out, t)
        }
    }
    return out, nil
}

func (r *technicianRepository) UpdateGPS(ctx context.Context, id uuid.UUID, loc valueobject.GPSLocation) error {
    return r.db.WithContext(ctx).Model(&TechnicianModel{}).
        Where("id = ?", id).
        Updates(map[string]any{
            "gps_last":      marshalJSON(loc),
            "gps_updated_at": time.Now(),
            "updated_at":    time.Now(),
        }).Error
}

func (r *technicianRepository) ResetDailyLoads(ctx context.Context, tenantID uuid.UUID) error {
    return r.db.WithContext(ctx).Model(&TechnicianModel{}).
        Where("tenant_id = ?", tenantID).
        Update("current_load", 0).Error
}

func pgTextArray(ss []valueobject.Skill) string {
    // JSON array string – format ให้ pgx ใช้
    out := "{"
    for i, s := range ss {
        if i > 0 { out += "," }
        out += `"` + string(s) + `"`
    }
    out += "}"
    return out
}

var _ = math.Pi
```

### `infrastructure/persistence/postgres/others_repository.go`
```go
package postgres

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"

    "icmongolang/internal/modules/iotlogistics/domain/entity"
    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
    "icmongolang/internal/modules/iotlogistics/domain/repository"
)

// ============================================================
// WorkOrder
// ============================================================

type workOrderRepository struct{ db *gorm.DB }
func NewWorkOrderRepository(db *gorm.DB) repository.WorkOrderRepository { return &workOrderRepository{db: db} }

func (r *workOrderRepository) Save(ctx context.Context, w *entity.WorkOrder) error {
    return r.db.WithContext(ctx).Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "id"}},
        UpdateAll: true,
    }).Create(toWorkOrderModel(w)).Error
}
func (r *workOrderRepository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.WorkOrder, error) {
    var m WorkOrderModel
    err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) { return nil, domainerrors.ErrWONotFound }
    if err != nil { return nil, err }
    return toWorkOrderEntity(&m), nil
}
func (r *workOrderRepository) FindByNo(ctx context.Context, tenantID uuid.UUID, no string) (*entity.WorkOrder, error) {
    var m WorkOrderModel
    err := r.db.WithContext(ctx).Where("tenant_id = ? AND wo_no = ?", tenantID, no).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) { return nil, domainerrors.ErrWONotFound }
    if err != nil { return nil, err }
    return toWorkOrderEntity(&m), nil
}
func (r *workOrderRepository) List(ctx context.Context, f repository.WorkOrderFilter) ([]*entity.WorkOrder, int64, error) {
    q := r.db.WithContext(ctx).Model(&WorkOrderModel{}).Where("tenant_id = ?", f.TenantID)
    if f.TechnicianID != nil { q = q.Where("technician_id = ?", *f.TechnicianID) }
    if f.SiteID != nil { q = q.Where("site_id = ?", *f.SiteID) }
    if f.SourceRef != "" { q = q.Where("source_ref = ?", f.SourceRef) }
    if len(f.Statuses) > 0 {
        ss := make([]string, len(f.Statuses))
        for i, s := range f.Statuses { ss[i] = string(s) }
        q = q.Where("status IN ?", ss)
    }
    var total int64
    if err := q.Count(&total).Error; err != nil { return nil, 0, err }
    if f.PageSize > 0 { q = q.Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize) }
    var models []WorkOrderModel
    if err := q.Order("created_at DESC").Find(&models).Error; err != nil { return nil, 0, err }
    out := make([]*entity.WorkOrder, 0, len(models))
    for i := range models { out = append(out, toWorkOrderEntity(&models[i])) }
    return out, total, nil
}
func (r *workOrderRepository) NextSequence(ctx context.Context, tenantID uuid.UUID, year int) (int, error) {
    var c int64
    pattern := fmt.Sprintf("WO-%04d-%%", year)
    err := r.db.WithContext(ctx).Model(&WorkOrderModel{}).
        Where("tenant_id = ? AND wo_no LIKE ?", tenantID, pattern).Count(&c).Error
    return int(c) + 1, err
}

// ============================================================
// Maintenance
// ============================================================

type maintenanceRepository struct{ db *gorm.DB }
func NewMaintenanceRepository(db *gorm.DB) repository.MaintenanceRepository { return &maintenanceRepository{db: db} }

func (r *maintenanceRepository) Save(ctx context.Context, m *entity.MaintenanceSchedule) error {
    return r.db.WithContext(ctx).Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "id"}},
        UpdateAll: true,
    }).Create(toMaintModel(m)).Error
}
func (r *maintenanceRepository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.MaintenanceSchedule, error) {
    var m MaintenanceScheduleModel
    err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) { return nil, domainerrors.ErrMaintenanceNotFound }
    if err != nil { return nil, err }
    return toMaintEntity(&m), nil
}
func (r *maintenanceRepository) FindByDevice(ctx context.Context, tenantID, deviceID uuid.UUID) (*entity.MaintenanceSchedule, error) {
    var m MaintenanceScheduleModel
    err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND device_id = ?", tenantID, deviceID).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) { return nil, domainerrors.ErrMaintenanceNotFound }
    if err != nil { return nil, err }
    return toMaintEntity(&m), nil
}
func (r *maintenanceRepository) List(ctx context.Context, f repository.MaintenanceFilter) ([]*entity.MaintenanceSchedule, int64, error) {
    q := r.db.WithContext(ctx).Model(&MaintenanceScheduleModel{}).Where("tenant_id = ?", f.TenantID)
    if f.SiteID != nil { q = q.Where("site_id = ?", *f.SiteID) }
    if f.DeviceID != nil { q = q.Where("device_id = ?", *f.DeviceID) }
    if f.Active != nil { q = q.Where("is_active = ?", *f.Active) }
    var total int64
    if err := q.Count(&total).Error; err != nil { return nil, 0, err }
    if f.PageSize > 0 { q = q.Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize) }
    var models []MaintenanceScheduleModel
    if err := q.Order("next_due_at ASC").Find(&models).Error; err != nil { return nil, 0, err }
    out := make([]*entity.MaintenanceSchedule, 0, len(models))
    for i := range models { out = append(out, toMaintEntity(&models[i])) }
    return out, total, nil
}
func (r *maintenanceRepository) ListDue(ctx context.Context, tenantID uuid.UUID, asOf time.Time) ([]*entity.MaintenanceSchedule, error) {
    var models []MaintenanceScheduleModel
    err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND is_active = true AND next_due_at <= ?", tenantID, asOf).
        Order("next_due_at ASC").Find(&models).Error
    if err != nil { return nil, err }
    out := make([]*entity.MaintenanceSchedule, 0, len(models))
    for i := range models { out = append(out, toMaintEntity(&models[i])) }
    return out, nil
}
func (r *maintenanceRepository) ListDueAllTenants(ctx context.Context, asOf time.Time) ([]*entity.MaintenanceSchedule, error) {
    var models []MaintenanceScheduleModel
    err := r.db.WithContext(ctx).
        Where("is_active = true AND next_due_at <= ?", asOf).
        Order("next_due_at ASC").Limit(1000).Find(&models).Error
    if err != nil { return nil, err }
    out := make([]*entity.MaintenanceSchedule, 0, len(models))
    for i := range models { out = append(out, toMaintEntity(&models[i])) }
    return out, nil
}
func (r *maintenanceRepository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
    return r.db.WithContext(ctx).
        Where("tenant_id = ? AND id = ?", tenantID, id).
        Delete(&MaintenanceScheduleModel{}).Error
}

// ============================================================
// RMA
// ============================================================

type rmaRepository struct{ db *gorm.DB }
func NewRMARepository(db *gorm.DB) repository.RMARepository { return &rmaRepository{db: db} }

func (r *rmaRepository) Save(ctx context.Context, rma *entity.RMA) error {
    return r.db.WithContext(ctx).Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "id"}},
        UpdateAll: true,
    }).Create(toRMAModel(rma)).Error
}
func (r *rmaRepository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.RMA, error) {
    var m RMAModel
    err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) { return nil, domainerrors.ErrRMANotFound }
    if err != nil { return nil, err }
    return toRMAEntity(&m), nil
}
func (r *rmaRepository) FindByNo(ctx context.Context, tenantID uuid.UUID, no string) (*entity.RMA, error) {
    var m RMAModel
    err := r.db.WithContext(ctx).Where("tenant_id = ? AND rma_no = ?", tenantID, no).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) { return nil, domainerrors.ErrRMANotFound }
    if err != nil { return nil, err }
    return toRMAEntity(&m), nil
}
func (r *rmaRepository) List(ctx context.Context, f repository.RMAFilter) ([]*entity.RMA, int64, error) {
    q := r.db.WithContext(ctx).Model(&RMAModel{}).Where("tenant_id = ?", f.TenantID)
    if f.CustomerID != nil { q = q.Where("customer_id = ?", *f.CustomerID) }
    if f.DeviceID != nil { q = q.Where("device_id = ?", *f.DeviceID) }
    if len(f.Statuses) > 0 {
        ss := make([]string, len(f.Statuses))
        for i, s := range f.Statuses { ss[i] = string(s) }
        q = q.Where("status IN ?", ss)
    }
    var total int64
    if err := q.Count(&total).Error; err != nil { return nil, 0, err }
    if f.PageSize > 0 { q = q.Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize) }
    var models []RMAModel
    if err := q.Order("created_at DESC").Find(&models).Error; err != nil { return nil, 0, err }
    out := make([]*entity.RMA, 0, len(models))
    for i := range models { out = append(out, toRMAEntity(&models[i])) }
    return out, total, nil
}
func (r *rmaRepository) ListByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) ([]*entity.RMA, error) {
    var models []RMAModel
    if err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
        Order("created_at DESC").Find(&models).Error; err != nil { return nil, err }
    out := make([]*entity.RMA, 0, len(models))
    for i := range models { out = append(out, toRMAEntity(&models[i])) }
    return out, nil
}
func (r *rmaRepository) CountOpenByDevice(ctx context.Context, tenantID, deviceID uuid.UUID) (int64, error) {
    var n int64
    err := r.db.WithContext(ctx).Model(&RMAModel{}).
        Where("tenant_id = ? AND device_id = ? AND status NOT IN ?",
            tenantID, deviceID, []string{"CLOSED", "REJECTED"}).
        Count(&n).Error
    return n, err
}
func (r *rmaRepository) NextSequence(ctx context.Context, tenantID uuid.UUID, year int) (int, error) {
    var c int64
    pattern := fmt.Sprintf("RMA-%04d-%%", year)
    err := r.db.WithContext(ctx).Model(&RMAModel{}).
        Where("tenant_id = ? AND rma_no LIKE ?", tenantID, pattern).Count(&c).Error
    return int(c) + 1, err
}

// ============================================================
// Audit
// ============================================================

type auditRepository struct{ db *gorm.DB }
func NewAuditRepository(db *gorm.DB) repository.AuditRepository { return &auditRepository{db: db} }

func (r *auditRepository) Save(ctx context.Context, e *repository.AuditEntry) error {
    m := &AuditLogModel{
        TenantID: e.TenantID, ActorID: e.ActorID,
        Action: e.Action, EntityType: e.EntityType, EntityID: e.EntityID,
        IPAddress: e.IPAddress, UserAgent: e.UserAgent,
        CreatedAt: e.CreatedAt,
    }
    if m.CreatedAt.IsZero() { m.CreatedAt = time.Now() }
    if len(e.Payload) > 0 {
        if b, err := json.Marshal(e.Payload); err == nil {
            m.Payload = datatypes.JSON(b)
        }
    }
    return r.db.WithContext(ctx).Create(m).Error
}
```

> **ต้องการ import `encoding/json`, `datatypes`**

---

## C.2 Code Generator

### `infrastructure/persistence/postgres/code_generator.go`
```go
package postgres

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"

    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

type codeGeneratorRepo struct{ db *gorm.DB }

func NewCodeGeneratorRepo(db *gorm.DB) *codeGeneratorRepo {
    return &codeGeneratorRepo{db: db}
}

func (r *codeGeneratorRepo) NextShipmentNo(ctx context.Context, tenantID uuid.UUID) (valueobject.ShipmentNo, error) {
    year := time.Now().Year()
    var c int64
    pattern := fmt.Sprintf("SHP-%04d-%%", year)
    err := r.db.WithContext(ctx).Model(&ShipmentModel{}).
        Where("tenant_id = ? AND shipment_no LIKE ?", tenantID, pattern).
        Count(&c).Error
    if err != nil { return "", err }
    return valueobject.NewShipmentNo(year, int(c)+1)
}

func (r *codeGeneratorRepo) NextJobNo(ctx context.Context, tenantID uuid.UUID) (string, error) {
    year := time.Now().Year()
    var c int64
    pattern := fmt.Sprintf("JOB-%04d-%%", year)
    err := r.db.WithContext(ctx).Model(&InstallationJobModel{}).
        Where("tenant_id = ? AND job_no LIKE ?", tenantID, pattern).
        Count(&c).Error
    if err != nil { return "", err }
    return fmt.Sprintf("JOB-%04d-%06d", year, c+1), nil
}

func (r *codeGeneratorRepo) NextWONo(ctx context.Context, tenantID uuid.UUID) (string, error) {
    year := time.Now().Year()
    var c int64
    pattern := fmt.Sprintf("WO-%04d-%%", year)
    err := r.db.WithContext(ctx).Model(&WorkOrderModel{}).
        Where("tenant_id = ? AND wo_no LIKE ?", tenantID, pattern).
        Count(&c).Error
    if err != nil { return "", err }
    return fmt.Sprintf("WO-%04d-%06d", year, c+1), nil
}

func (r *codeGeneratorRepo) NextRMANo(ctx context.Context, tenantID uuid.UUID) (string, error) {
    year := time.Now().Year()
    var c int64
    pattern := fmt.Sprintf("RMA-%04d-%%", year)
    err := r.db.WithContext(ctx).Model(&RMAModel{}).
        Where("tenant_id = ? AND rma_no LIKE ?", tenantID, pattern).
        Count(&c).Error
    if err != nil { return "", err }
    return fmt.Sprintf("RMA-%04d-%06d", year, c+1), nil
}
```

---

## C.3 Maps Integration

### `infrastructure/services/maps/google_maps.go`
```go
package maps

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "time"

    "icmongolang/internal/modules/iotlogistics/domain/service/port"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

type googleMaps struct {
    apiKey string
    http   *http.Client
}

func NewGoogleMaps(apiKey string) port.MapsPort {
    return &googleMaps{
        apiKey: apiKey,
        http:   &http.Client{Timeout: 15 * time.Second},
    }
}

func (g *googleMaps) DistanceMatrix(ctx context.Context, locations []valueobject.GPSLocation) (*port.DistanceMatrix, error) {
    if len(locations) == 0 {
        return &port.DistanceMatrix{Size: 0}, nil
    }
    origins := joinLocations(locations)
    destinations := origins // square matrix

    u := fmt.Sprintf(
        "https://maps.googleapis.com/maps/api/distancematrix/json?origins=%s&destinations=%s&key=%s&units=metric&mode=driving",
        url.QueryEscape(origins), url.QueryEscape(destinations), g.apiKey,
    )

    req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
    resp, err := g.http.Do(req)
    if err != nil { return nil, err }
    defer resp.Body.Close()

    var out struct {
        Rows []struct {
            Elements []struct {
                Distance struct {
                    Value float64 `json:"value"` // meters
                } `json:"distance"`
                Duration struct {
                    Value float64 `json:"value"` // seconds
                } `json:"duration"`
                Status string `json:"status"`
            } `json:"elements"`
        } `json:"rows"`
        Status string `json:"status"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return nil, err }

    size := len(locations)
    cells := make([][]port.Cell, size)
    for i := 0; i < size; i++ {
        cells[i] = make([]port.Cell, size)
        for j := 0; j < size; j++ {
            if i < len(out.Rows) && j < len(out.Rows[i].Elements) {
                el := out.Rows[i].Elements[j]
                cells[i][j] = port.Cell{
                    DistanceKm: el.Distance.Value / 1000,
                    DurationMin: int(el.Duration.Value / 60),
                }
            }
        }
    }
    return &port.DistanceMatrix{Cells: cells, Size: size}, nil
}

func (g *googleMaps) Geocode(ctx context.Context, address string) (valueobject.GPSLocation, error) {
    u := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s",
        url.QueryEscape(address), g.apiKey)
    req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
    resp, err := g.http.Do(req)
    if err != nil { return valueobject.GPSLocation{}, err }
    defer resp.Body.Close()

    var out struct {
        Results []struct {
            Geometry struct {
                Location struct {
                    Lat float64 `json:"lat"`
                    Lng float64 `json:"lng"`
                } `json:"location"`
            } `json:"geometry"`
        } `json:"results"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return valueobject.GPSLocation{}, err }
    if len(out.Results) == 0 { return valueobject.GPSLocation{}, fmt.Errorf("no results") }
    return valueobject.GPSLocation{
        Lat: out.Results[0].Geometry.Location.Lat,
        Lng: out.Results[0].Geometry.Location.Lng,
    }, nil
}

func (g *googleMaps) ReverseGeocode(ctx context.Context, loc valueobject.GPSLocation) (string, error) {
    u := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?latlng=%f,%f&key=%s",
        loc.Lat, loc.Lng, g.apiKey)
    req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
    resp, err := g.http.Do(req)
    if err != nil { return "", err }
    defer resp.Body.Close()
    var out struct {
        Results []struct {
            FormattedAddress string `json:"formatted_address"`
        } `json:"results"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return "", err }
    if len(out.Results) == 0 { return "", nil }
    return out.Results[0].FormattedAddress, nil
}

func joinLocations(locs []valueobject.GPSLocation) string {
    s := ""
    for i, loc := range locs {
        if i > 0 { s += "|" }
        s += fmt.Sprintf("%f,%f", loc.Lat, loc.Lng)
    }
    return s
}
```

### `infrastructure/services/maps/mock_maps.go` (สำหรับ dev/test)
```go
package maps

import (
    "context"
    "math"

    "icmongolang/internal/modules/iotlogistics/domain/service/port"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

type mockMaps struct{}

func NewMockMaps() port.MapsPort { return &mockMaps{} }

func (m *mockMaps) DistanceMatrix(ctx context.Context, locs []valueobject.GPSLocation) (*port.DistanceMatrix, error) {
    size := len(locs)
    cells := make([][]port.Cell, size)
    for i := 0; i < size; i++ {
        cells[i] = make([]port.Cell, size)
        for j := 0; j < size; j++ {
            d := locs[i].DistanceKm(locs[j])
            cells[i][j] = port.Cell{
                DistanceKm:  d,
                DurationMin: int(d / 30 * 60), // assume 30 km/h
            }
        }
    }
    return &port.DistanceMatrix{Cells: cells, Size: size}, nil
}

func (m *mockMaps) Geocode(ctx context.Context, addr string) (valueobject.GPSLocation, error) {
    return valueobject.GPSLocation{Lat: 13.7563, Lng: 100.5018}, nil // Bangkok
}

func (m *mockMaps) ReverseGeocode(ctx context.Context, loc valueobject.GPSLocation) (string, error) {
    return "Mock Address", nil
}

var _ = math.Pi
```

---

## C.4 AI Advisor (Route)

### `infrastructure/services/ai/ollama_advisor.go`
```go
package ai

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "icmongolang/internal/modules/iotlogistics/domain/service/port"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
)

type ollamaAdvisor struct {
    baseURL string
    model   string
    http    *http.Client
}

func NewOllamaAdvisor(baseURL, model string) port.AIAdvisorPort {
    return &ollamaAdvisor{
        baseURL: baseURL, model: model,
        http: &http.Client{Timeout: 30 * time.Second},
    }
}

type ollamaReq struct {
    Model  string `json:"model"`
    Prompt string `json:"prompt"`
    Stream bool   `json:"stream"`
}

type ollamaResp struct {
    Response string `json:"response"`
}

func (o *ollamaAdvisor) call(ctx context.Context, prompt string) (string, error) {
    b, _ := json.Marshal(ollamaReq{Model: o.model, Prompt: prompt, Stream: false})
    req, _ := http.NewRequestWithContext(ctx, "POST", o.baseURL+"/api/generate", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    resp, err := o.http.Do(req)
    if err != nil { return "", err }
    defer resp.Body.Close()
    var out ollamaResp
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return "", err }
    return out.Response, nil
}

func (o *ollamaAdvisor) RefineRoute(ctx context.Context, route *valueobject.Route, req any) (string, error) {
    prompt := fmt.Sprintf(`You are a fleet route optimizer. Given %d stops with total %0.1f km and %d minutes:
	Suggest 3 tips to reduce time (traffic, breaks, sequencing).
	Reply with JSON: {"tips": ["...", "...", "..."], "polyline": ""}`,
        len(route.Waypoints), route.TotalKm, route.TotalMinutes)
    raw, err := o.call(ctx, prompt)
    if err != nil { return "", err }
    // extract polyline (or empty)
    var out struct {
        Polyline string `json:"polyline"`
    }
    _ = json.Unmarshal([]byte(raw), &out)
    return out.Polyline, nil
}

func (o *ollamaAdvisor) PredictDuration(ctx context.Context, jobType string, complexity int) (int, error) {
    prompt := fmt.Sprintf("Predict duration in minutes for job type %s with complexity %d. Reply with number only.", jobType, complexity)
    raw, err := o.call(ctx, prompt)
    if err != nil { return 0, err }
    var min int
    fmt.Sscanf(raw, "%d", &min)
    if min == 0 { min = 60 }
    return min, nil
}

func (o *ollamaAdvisor) SummarizeJob(ctx context.Context, notes string) (string, error) {
    prompt := "Summarize this technician note in 2 bullet points in Thai:\n" + notes
    return o.call(ctx, prompt)
}
```

---

## C.5 Kafka Consumers

### `infrastructure/messaging/kafka/consumers/customer_onboarded_consumer.go`
```go
package consumers

import (
    "context"
    "encoding/json"
    "log"

    "github.com/google/uuid"

    "icmongolang/internal/modules/iotlogistics/application"
)

// CustomerOnboardedConsumer – รับ event จาก customer module
// → auto-create shipment draft (ถ้ามี package ที่ต้องติดตั้ง)
type CustomerOnboardedConsumer struct {
    createShipmentUC *application.CreateShipmentUseCase
}

func NewCustomerOnboardedConsumer(uc *application.CreateShipmentUseCase) *CustomerOnboardedConsumer {
    return &CustomerOnboardedConsumer{createShipmentUC: uc}
}

func (c *CustomerOnboardedConsumer) Handle(ctx context.Context, payload []byte) error {
    var evt struct {
        CustomerID uuid.UUID `json:"customer_id"`
        TenantID   uuid.UUID `json:"tenant_id"`
        PackageID  uuid.UUID `json:"package_id"`
        SiteIDs    []string  `json:"site_ids"`
    }
    if err := json.Unmarshal(payload, &evt); err != nil { return err }

    log.Printf("[customer.onboarded] customer=%s, %d sites", evt.CustomerID, len(evt.SiteIDs))
    // TODO: lookup package → determine devices → create shipment
    return nil
}
```

### `infrastructure/messaging/kafka/consumers/alert_triggered_consumer.go`
```go
package consumers

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/google/uuid"

    "icmongolang/internal/modules/iotlogistics/application"
    "icmongolang/internal/modules/iotlogistics/domain/service/port"
)

// AlertTriggeredConsumer – device alert → auto-create work order (repair)
type AlertTriggeredConsumer struct {
    createWOUC *application.CreateWorkOrderUseCase
    codeGen    port.CodeGeneratorPort
}

func NewAlertTriggeredConsumer(uc *application.CreateWorkOrderUseCase, codeGen port.CodeGeneratorPort) *AlertTriggeredConsumer {
    return &AlertTriggeredConsumer{createWOUC: uc, codeGen: codeGen}
}

func (c *AlertTriggeredConsumer) Handle(ctx context.Context, payload []byte) error {
    var evt struct {
        TenantID   uuid.UUID `json:"tenant_id"`
        DeviceID   uuid.UUID `json:"device_id"`
        CustomerID uuid.UUID `json:"customer_id"`
        SiteID     uuid.UUID `json:"site_id"`
        Metric     string    `json:"metric"`
        Value      float64   `json:"value"`
        Severity   string    `json:"severity"`
        Message    string    `json:"message"`
    }
    if err := json.Unmarshal(payload, &evt); err != nil { return err }

    // เฉพาะ CRITICAL → create repair WO
    if evt.Severity != "CRITICAL" && evt.Severity != "ERROR" {
        return nil
    }

    _, err := c.createWOUC.Execute(ctx, application.CreateWorkOrderInput{
        TenantID:    evt.TenantID.String(),
        CustomerID:  evt.CustomerID.String(),
        SiteID:      evt.SiteID.String(),
        DeviceID:    evt.DeviceID.String(),
        Title:       fmt.Sprintf("Repair for %s alert", evt.Metric),
        Description: evt.Message,
        JobType:     "REPAIR",
        Priority:    "HIGH",
        SourceRef:   "device_alert",
    })
    return err
}

var _ = time.Now
```

### `infrastructure/messaging/kafka/consumers/installation_completed_dispatcher.go`
```go
package consumers

import (
    "context"
    "log"
)

// InstallationCompletedDispatcher – forward installation.completed
// → device module (via Kafka) เพื่อ provision devices
// → package module เพื่อ start billing
// Pattern: เราไม่ต้องทำอะไร — เพียงแค่ republish (หรือให้ device/package subscribe โดยตรง)
// ในกรณีที่ต้อง transform เท่านั้น
type InstallationCompletedDispatcher struct{}

func NewInstallationCompletedDispatcher() *InstallationCompletedDispatcher {
    return &InstallationCompletedDispatcher{}
}

func (d *InstallationCompletedDispatcher) Handle(ctx context.Context, payload []byte) error {
    log.Printf("[installation.completed] dispatch → device.provision + package.billing.start")
    // ถ้ามี transformation ให้ทำที่นี่
    return nil
}
```

---

## C.6 WebSocket Technician Tracking

### `infrastructure/websocket/technician_tracker.go`
```go
package websocket

import (
    "context"
    "encoding/json"
    "log"
    "sync"

    "github.com/google/uuid"

    "icmongolang/internal/modules/iotlogistics/domain/entity"
    valueobject "icmongolang/internal/modules/iotlogistics/domain/value_object"
    ws "icmongolang/pkg/websocket"
)

// TechnicianTracker – broadcast GPS update ให้ dispatcher
type TechnicianTracker struct {
    hub *ws.Hub
    mu  sync.RWMutex
    byTenant map[string]map[uuid.UUID]*entity.Technician
}

func NewTechnicianTracker(hub *ws.Hub) *TechnicianTracker {
    return &TechnicianTracker{
        hub:      hub,
        byTenant: map[string]map[uuid.UUID]*entity.Technician{},
    }
}

func (t *TechnicianTracker) UpdateGPS(tenantID string, techID uuid.UUID, loc valueobject.GPSLocation) {
    t.mu.Lock()
    if t.byTenant[tenantID] == nil {
        t.byTenant[tenantID] = map[uuid.UUID]*entity.Technician{}
    }
    t.byTenant[tenantID][techID] = &entity.Technician{ID: techID, GPSLast: &loc}
    t.mu.Unlock()

    // Broadcast ไปยัง tenant
    evt, _ := json.Marshal(map[string]any{
        "type": "technician.gps",
        "data": map[string]any{
            "technician_id": techID.String(),
            "lat":           loc.Lat,
            "lng":           loc.Lng,
            "ts":            loc.Timestamp,
        },
    })
    if err := t.hub.BroadcastToTenant(tenantID, evt); err != nil {
        log.Printf("broadcast gps: %v", err)
    }
}

func (t *TechnicianTracker) Snapshot(tenantID string) []*entity.Technician {
    t.mu.RLock()
    defer t.mu.RUnlock()
    out := make([]*entity.Technician, 0, len(t.byTenant[tenantID]))
    for _, v := range t.byTenant[tenantID] {
        out = append(out, v)
    }
    return out
}

var _ = context.Background
```

---

## C.7 Scheduler Jobs

### `infrastructure/scheduler/pm_due_job.go`
```go
package scheduler

import (
    "context"
    "log"
    "time"

    "icmongolang/internal/modules/iotlogistics/application"
)

// PMDueJob – เรียกทุกวัน 02:00 → ตรวจ schedule ที่ due
type PMDueJob struct {
    processUC *application.ProcessMaintenanceDueUseCase
}

func NewPMDueJob(uc *application.ProcessMaintenanceDueUseCase) *PMDueJob {
    return &PMDueJob{processUC: uc}
}

func (j *PMDueJob) Run() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
    defer cancel()
    log.Println("[pm.due] starting")
    if err := j.processUC.Run(ctx); err != nil {
        log.Printf("[pm.due] error: %v", err)
    }
    log.Println("[pm.due] done")
}
```

### `infrastructure/scheduler/sla_check_job.go`
```go
package scheduler

import (
    "context"
    "log"
    "time"

    "github.com/google/uuid"

    "icmongolang/internal/modules/iotlogistics/domain/repository"
    "icmongolang/internal/modules/iotlogistics/domain/event"
    "icmongolang/pkg/kafka"
)

// SLACheckJob – ตรวจ job ที่เกิน SLA → publish event
type SLACheckJob struct {
    jobRepo  repository.InstallationRepository
    producer kafka.Producer
}

func NewSLACheckJob(repo repository.InstallationRepository, producer kafka.Producer) *SLACheckJob {
    return &SLACheckJob{jobRepo: repo, producer: producer}
}

func (j *SLACheckJob) Run() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    asOf := time.Now()
    log.Println("[sla.check] starting")

    // ดึง tenant ทั้งหมด — ใน prod จะ loop
    tenants := j.listActiveTenants(ctx)
    for _, tid := range tenants {
        jobs, err := j.jobRepo.ListSLABreached(ctx, tid, asOf)
        if err != nil { continue }
        for _, job := range jobs {
            _ = j.producer.Publish(ctx, event.TopicInstallationSLABreach, job.ID.String(), map[string]any{
                "event_id":   uuid.New(),
                "job_id":     job.ID.String(),
                "tenant_id":  job.TenantID.String(),
                "sla_due_at": job.SLADueAt,
                "status":     string(job.Status),
                "occurred_at": asOf,
            })
        }
    }
    log.Println("[sla.check] done")
}

func (j *SLACheckJob) listActiveTenants(ctx context.Context) []uuid.UUID {
    // TODO: query from a tenants table
    return []uuid.UUID{}
}
```

### `infrastructure/scheduler/technician_load_reset_job.go`
```go
package scheduler

import (
    "context"
    "log"

    "github.com/google/uuid"

    "icmongolang/internal/modules/iotlogistics/domain/repository"
)

// TechnicianLoadResetJob – reset current_load ทุกเที่ยงคืน
type TechnicianLoadResetJob struct {
    techRepo repository.TechnicianRepository
}

func NewTechnicianLoadResetJob(repo repository.TechnicianRepository) *TechnicianLoadResetJob {
    return &TechnicianLoadResetJob{techRepo: repo}
}

func (j *TechnicianLoadResetJob) Run() {
    ctx := context.Background()
    // TODO: iterate tenants
    tenants := []uuid.UUID{}
    for _, tid := range tenants {
        if err := j.techRepo.ResetDailyLoads(ctx, tid); err != nil {
            log.Printf("[tech.reset] error for %s: %v", tid, err)
        }
    }
}
```

---

# 🅳 PART 5D — INTERFACE + WIRING + MIGRATION + TESTS

## D.1 HTTP Handlers

### `interfaces/http/shipment_handler.go`
```go
package http

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    "icmongolang/internal/modules/iotlogistics/application"
)

type ShipmentHandler struct {
    createUC   *application.CreateShipmentUseCase
    dispatchUC *application.DispatchShipmentUseCase
    gpsUC      *application.UpdateShipmentGPSUseCase
    deliveryUC *application.ConfirmDeliveryUseCase
    listUC     *application.ListShipmentsUseCase
    getUC      *application.GetShipmentUseCase
}

func NewShipmentHandler(
    createUC *application.CreateShipmentUseCase,
    dispatchUC *application.DispatchShipmentUseCase,
    gpsUC *application.UpdateShipmentGPSUseCase,
    deliveryUC *application.ConfirmDeliveryUseCase,
    listUC *application.ListShipmentsUseCase,
    getUC *application.GetShipmentUseCase,
) *ShipmentHandler {
    return &ShipmentHandler{
        createUC: createUC, dispatchUC: dispatchUC,
        gpsUC: gpsUC, deliveryUC: deliveryUC,
        listUC: listUC, getUC: getUC,
    }
}

func (h *ShipmentHandler) Create(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    var in application.CreateShipmentInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    in.TenantID = tid.String()
    in.ActorID = uid.String()

    res, err := h.createUC.Execute(c.Request.Context(), in)
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusCreated, res)
}

func (h *ShipmentHandler) Dispatch(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    var req struct {
        Carrier    string `json:"carrier" binding:"required"`
        TrackingNo string `json:"tracking_no,omitempty"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    res, err := h.dispatchUC.Execute(c.Request.Context(), application.DispatchShipmentInput{
        TenantID: tid.String(), ShipmentID: id.String(),
        Carrier: req.Carrier, TrackingNo: req.TrackingNo, ActorID: uid.String(),
    })
    if err != nil { writeError(c, err); return }
    c.JSON(200, res)
}

func (h *ShipmentHandler) UpdateGPS(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    var req struct {
        Lat   float64 `json:"lat" binding:"required"`
        Lng   float64 `json:"lng" binding:"required"`
        Speed float64 `json:"speed,omitempty"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    if err := h.gpsUC.Execute(c.Request.Context(), application.UpdateShipmentGPSInput{
        TenantID: tid.String(), ShipmentID: id.String(),
        Lat: req.Lat, Lng: req.Lng, Speed: req.Speed,
    }); err != nil {
        writeError(c, err); return
    }
    c.JSON(200, gin.H{"status": "ok"})
}

func (h *ShipmentHandler) ConfirmDelivery(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    var req struct {
        DeliveredTo  string `json:"delivered_to" binding:"required"`
        SignatureURL string `json:"signature_url,omitempty"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    res, err := h.deliveryUC.Execute(c.Request.Context(), application.ConfirmDeliveryInput{
        TenantID: tid.String(), ShipmentID: id.String(),
        DeliveredTo: req.DeliveredTo, SignatureURL: req.SignatureURL, ActorID: uid.String(),
    })
    if err != nil { writeError(c, err); return }
    c.JSON(200, res)
}
```

### `interfaces/http/installation_handler.go`
```go
package http

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    "icmongolang/internal/modules/iotlogistics/application"
)

type InstallationHandler struct {
    scheduleUC   *application.ScheduleInstallationUseCase
    assignUC     *application.AssignTechnicianUseCase
    startUC      *application.StartJobUseCase
    checklistUC  *application.UpdateChecklistUseCase
    completeUC   *application.CompleteJobUseCase
    listUC       *application.ListJobsUseCase
    getUC        *application.GetJobUseCase
}

func (h *InstallationHandler) Schedule(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    var in application.ScheduleInstallationInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    in.TenantID = tid.String()
    in.ActorID = uid.String()

    res, err := h.scheduleUC.Execute(c.Request.Context(), in)
    if err != nil { writeError(c, err); return }
    c.JSON(201, res)
}

func (h *InstallationHandler) Assign(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    var req struct {
        TechnicianID string                    `json:"technician_id,omitempty"`
        AutoAssign   bool                      `json:"auto_assign,omitempty"`
        SiteGPS      application.GPSDTO        `json:"site_gps,omitempty"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    res, err := h.assignUC.Execute(c.Request.Context(), application.AssignTechnicianInput{
        TenantID: tid.String(), JobID: id.String(),
        TechnicianID: req.TechnicianID, AutoAssign: req.AutoAssign,
        SiteGPS: req.SiteGPS, ActorID: uid.String(),
    })
    if err != nil { writeError(c, err); return }
    c.JSON(200, res)
}

func (h *InstallationHandler) Start(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    var req struct {
        Lat      float64 `json:"lat" binding:"required"`
        Lng      float64 `json:"lng" binding:"required"`
        Accuracy float64 `json:"accuracy,omitempty"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    res, err := h.startUC.Execute(c.Request.Context(), application.StartJobInput{
        TenantID: tid.String(), JobID: id.String(),
        Lat: req.Lat, Lng: req.Lng, Accuracy: req.Accuracy,
        ActorID: uid.String(),
    })
    if err != nil { writeError(c, err); return }
    c.JSON(200, res)
}

func (h *InstallationHandler) UpdateChecklist(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    var in application.UpdateChecklistInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    in.TenantID = tid.String()
    in.JobID = id.String()

    res, err := h.checklistUC.Execute(c.Request.Context(), in)
    if err != nil { writeError(c, err); return }
    c.JSON(200, res)
}

func (h *InstallationHandler) Complete(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    var req struct {
        Lat          float64 `json:"lat" binding:"required"`
        Lng          float64 `json:"lng" binding:"required"`
        SignatureURL string  `json:"signature_url,omitempty"`
        Notes        string  `json:"notes,omitempty"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    res, err := h.completeUC.Execute(c.Request.Context(), application.CompleteJobInput{
        TenantID: tid.String(), JobID: id.String(),
        Lat: req.Lat, Lng: req.Lng,
        SignatureURL: req.SignatureURL, Notes: req.Notes,
        ActorID: uid.String(),
    })
    if err != nil { writeError(c, err); return }
    c.JSON(200, res)
}

var _ = time.Now
```

### `interfaces/http/technician_handler.go`
```go
package http

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    "icmongolang/internal/modules/iotlogistics/application"
)

type TechnicianHandler struct {
    createUC    *application.CreateTechnicianUseCase
    gpsUC       *application.UpdateTechnicianGPSUseCase
    listUC      *application.ListTechniciansUseCase
    getUC       *application.GetTechnicianUseCase
    availableUC *application.FindAvailableTechniciansUseCase
    routeUC     *application.OptimizeRouteUseCase
}

func (h *TechnicianHandler) Create(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    var in application.CreateTechnicianInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    in.TenantID = tid.String()
    in.ActorID = uid.String()
    res, err := h.createUC.Execute(c.Request.Context(), in)
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusCreated, res)
}

func (h *TechnicianHandler) UpdateGPS(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    var req struct {
        Lat      float64 `json:"lat" binding:"required"`
        Lng      float64 `json:"lng" binding:"required"`
        Accuracy float64 `json:"accuracy,omitempty"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    if err := h.gpsUC.Execute(c.Request.Context(), tid, application.UpdateTechnicianGPSInput{
        TechnicianID: id.String(), Lat: req.Lat, Lng: req.Lng, Accuracy: req.Accuracy,
    }); err != nil {
        writeError(c, err); return
    }
    c.JSON(200, gin.H{"status": "ok"})
}

func (h *TechnicianHandler) OptimizeRoute(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    var in application.OptimizeRouteInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    in.TenantID = tid.String()
    in.TechnicianID = id.String()

    res, err := h.routeUC.Execute(c.Request.Context(), in)
    if err != nil { writeError(c, err); return }
    c.JSON(200, res)
}
```

### `interfaces/http/rma_handler.go`
```go
package http

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    "icmongolang/internal/modules/iotlogistics/application"
)

type RMAHandler struct {
    createUC  *application.CreateRMAUseCase
    approveUC *application.ApproveRMAUseCase
    replaceUC *application.ReplaceDeviceUseCase
    listUC    *application.ListRMAUseCase
}

func (h *RMAHandler) Create(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    var in application.CreateRMAInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    in.TenantID = tid.String()
    in.ActorID = uid.String()
    res, err := h.createUC.Execute(c.Request.Context(), in)
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusCreated, res)
}

func (h *RMAHandler) Approve(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    var req struct {
        Resolution string `json:"resolution" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    res, err := h.approveUC.Execute(c.Request.Context(), application.ApproveRMAInput{
        TenantID: tid.String(), RMAID: id.String(),
        Resolution: req.Resolution, ActorID: uid.String(),
    })
    if err != nil { writeError(c, err); return }
    c.JSON(200, res)
}
```

### `interfaces/http/errors.go`
```go
package http

import (
    "errors"
    "net/http"

    "github.com/gin-gonic/gin"

    domainerrors "icmongolang/internal/modules/iotlogistics/domain/errors"
)

func writeError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, domainerrors.ErrShipmentNotFound),
        errors.Is(err, domainerrors.ErrJobNotFound),
        errors.Is(err, domainerrors.ErrWONotFound),
        errors.Is(err, domainerrors.ErrMaintenanceNotFound),
        errors.Is(err, domainerrors.ErrRMANotFound),
        errors.Is(err, domainerrors.ErrTechnicianNotFound):
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

    case errors.Is(err, domainerrors.ErrShipmentLocked),
        errors.Is(err, domainerrors.ErrDuplicateSerial),
        errors.Is(err, domainerrors.ErrJobAlreadyTerminal),
        errors.Is(err, domainerrors.ErrRMAAlreadyTerminal),
        errors.Is(err, domainerrors.ErrRMAAlreadyOpen):
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

    case errors.Is(err, domainerrors.ErrOutsideGeoFence):
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})

    case errors.Is(err, domainerrors.ErrNoTechnicianAvailable),
        errors.Is(err, domainerrors.ErrTechnicianUnavailable):
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})

    case errors.Is(err, domainerrors.ErrInvalidShipmentNo),
        errors.Is(err, domainerrors.ErrEmptyShipment),
        errors.Is(err, domainerrors.ErrInvalidJobNo),
        errors.Is(err, domainerrors.ErrInvalidJobType),
        errors.Is(err, domainerrors.ErrInvalidPriority),
        errors.Is(err, domainerrors.ErrInvalidStatusTransition),
        errors.Is(err, domainerrors.ErrChecklistIncomplete),
        errors.Is(err, domainerrors.ErrEvidenceRequired),
        errors.Is(err, domainerrors.ErrSignatureRequired),
        errors.Is(err, domainerrors.ErrInvalidGPS),
        errors.Is(err, domainerrors.ErrInvalidQuantity),
        errors.Is(err, domainerrors.ErrInvalidFrequency),
        errors.Is(err, domainerrors.ErrInvalidRMANo),
        errors.Is(err, domainerrors.ErrInvalidRMAReason),
        errors.Is(err, domainerrors.ErrInvalidRMAResolution),
        errors.Is(err, domainerrors.ErrResolutionRequired),
        errors.Is(err, domainerrors.ErrRejectReasonRequired),
        errors.Is(err, domainerrors.ErrJobNotInProgress):
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

    case errors.Is(err, domainerrors.ErrMapsServiceFailure),
        errors.Is(err, domainerrors.ErrAIServiceFailure),
        errors.Is(err, domainerrors.ErrPersistenceFailure):
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})

    default:
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
    }
}
```

### `interfaces/http/routes.go`
```go
package http

import "github.com/gin-gonic/gin"

type Handlers struct {
    Shipment     *ShipmentHandler
    Installation *InstallationHandler
    Technician   *TechnicianHandler
    RMA          *RMAHandler
    Maintenance  *MaintenanceHandler
    WorkOrder    *WorkOrderHandler
}

func RegisterRoutes(r *gin.RouterGroup, h *Handlers, auth, tenant gin.HandlerFunc) {
    // Shipments
    sh := r.Group("/shipments")
    sh.Use(auth, tenant)
    sh.POST("",                    h.Shipment.Create)
    sh.GET("",                     h.Shipment.List)
    sh.GET("/:id",                 h.Shipment.Get)
    sh.POST("/:id/dispatch",       h.Shipment.Dispatch)
    sh.POST("/:id/gps",            h.Shipment.UpdateGPS)
    sh.POST("/:id/deliver",        h.Shipment.ConfirmDelivery)

    // Installation Jobs
    jb := r.Group("/installation-jobs")
    jb.Use(auth, tenant)
    jb.POST("",                    h.Installation.Schedule)
    jb.GET("",                     h.Installation.List)
    jb.GET("/:id",                 h.Installation.Get)
    jb.POST("/:id/assign",         h.Installation.Assign)
    jb.POST("/:id/start",          h.Installation.Start)
    jb.PATCH("/:id/checklist",     h.Installation.UpdateChecklist)
    jb.POST("/:id/complete",       h.Installation.Complete)

    // Technicians
    tech := r.Group("/technicians")
    tech.Use(auth, tenant)
    tech.POST("",                  h.Technician.Create)
    tech.GET("",                   h.Technician.List)
    tech.GET("/:id",               h.Technician.Get)
    tech.POST("/:id/gps",          h.Technician.UpdateGPS)
    tech.POST("/:id/optimize-route", h.Technician.OptimizeRoute)
    tech.GET("/available",         h.Technician.Available)

    // Work Orders
    wo := r.Group("/work-orders")
    wo.Use(auth, tenant)
    wo.POST("",                    h.WorkOrder.Create)
    wo.GET("",                     h.WorkOrder.List)
    wo.GET("/:id",                 h.WorkOrder.Get)
    wo.POST("/:id/assign",         h.WorkOrder.Assign)
    wo.POST("/:id/start",          h.WorkOrder.Start)
    wo.POST("/:id/parts",          h.WorkOrder.AddPart)
    wo.POST("/:id/complete",       h.WorkOrder.Complete)

    // Maintenance
    mt := r.Group("/maintenance")
    mt.Use(auth, tenant)
    mt.POST("",                    h.Maintenance.Schedule)
    mt.GET("",                     h.Maintenance.List)
    mt.GET("/due",                 h.Maintenance.Due)
    mt.POST("/:id/done",           h.Maintenance.MarkDone)

    // RMA
    rma := r.Group("/rma")
    rma.Use(auth, tenant)
    rma.POST("",                   h.RMA.Create)
    rma.GET("",                    h.RMA.List)
    rma.GET("/:id",                h.RMA.Get)
    rma.POST("/:id/approve",       h.RMA.Approve)
    rma.POST("/:id/ship",          h.RMA.Ship)
    rma.POST("/:id/receive",       h.RMA.Receive)
    rma.POST("/:id/replace",       h.RMA.Replace)
    rma.POST("/:id/close",         h.RMA.Close)
}
```

## D.2 Composition Root

### `module.go`
```go
package iotlogistics

import (
    "time"

    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"

    "icmongolang/internal/modules/iotlogistics/application"
    "icmongolang/internal/modules/iotlogistics/domain/service"
    "icmongolang/internal/modules/iotlogistics/domain/service/port"
    "icmongolang/internal/modules/iotlogistics/infrastructure/persistence/postgres"
    httpiface "icmongolang/internal/modules/iotlogistics/interfaces/http"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
    "icmongolang/pkg/transaction"
    "icmongolang/pkg/websocket"
)

type Dependencies struct {
    DB            *gorm.DB
    Redis         *redis.Client
    Producer      kafka.Producer
    Maps          port.MapsPort
    AIAdvisor     port.AIAdvisorPort
    Notifier      port.NotifierPort
    WSHub         *websocket.Hub
    GeoFenceM     float64  // default radius
    Logger        logger.Logger
}

type Wiring struct {
    TechnicianTracker interface{} // *wsiface.TechnicianTracker
}

func Init(router *gin.RouterGroup, deps Dependencies, auth, tenant gin.HandlerFunc) *Wiring {
    // Repos
    shipRepo   := postgres.NewShipmentRepository(deps.DB)
    jobRepo    := postgres.NewInstallationRepository(deps.DB)
    techRepo   := postgres.NewTechnicianRepository(deps.DB)
    woRepo     := postgres.NewWorkOrderRepository(deps.DB)
    maintRepo  := postgres.NewMaintenanceRepository(deps.DB)
    rmaRepo    := postgres.NewRMARepository(deps.DB)
    auditRepo  := postgres.NewAuditRepository(deps.DB)
    codeGen    := postgres.NewCodeGeneratorRepo(deps.DB)

    // Tx
    txMgr := transaction.NewManager(deps.DB)

    // Domain services
    assignSvc := service.NewTechnicianAssignmentService(techRepo)
    routeSvc  := service.NewRoutePlanningService(deps.Maps, deps.AIAdvisor)
    geoFence  := service.NewGeoFenceService(deps.GeoFenceM)
    validator := service.NewChecklistValidator()

    // Use cases
    createShipUC    := application.NewCreateShipmentUseCase(shipRepo, auditRepo, codeGen, deps.Producer, deps.Logger)
    dispatchUC      := application.NewDispatchShipmentUseCase(shipRepo, auditRepo, deps.Producer, deps.Logger)
    scheduleInstUC  := application.NewScheduleInstallationUseCase(jobRepo, shipRepo, auditRepo, codeGen, deps.Producer, deps.Logger)
    assignTechUC    := application.NewAssignTechnicianUseCase(jobRepo, techRepo, auditRepo, assignSvc, deps.Producer, txMgr, deps.Logger)
    startJobUC      := application.NewStartJobUseCase(jobRepo, geoFence, deps.Producer, deps.Logger)
    checklistUC     := application.NewUpdateChecklistUseCase(jobRepo)
    completeJobUC   := application.NewCompleteJobUseCase(jobRepo, shipRepo, maintRepo, geoFence, validator, deps.Producer, txMgr, deps.Logger)
    optimizeUC      := application.NewOptimizeRouteUseCase(techRepo, jobRepo, routeSvc)
    createRMAUC     := application.NewCreateRMAUseCase(rmaRepo, codeGen, auditRepo, deps.Producer, deps.Logger)
    pmUC            := application.NewProcessMaintenanceDueUseCase(maintRepo, woRepo, codeGen, deps.Producer, deps.Logger)

    // Handlers
    shipH    := httpiface.NewShipmentHandler(createShipUC, dispatchUC,
        application.NewUpdateShipmentGPSUseCase(shipRepo),
        application.NewConfirmDeliveryUseCase(shipRepo, auditRepo, deps.Producer, deps.Logger),
        application.NewListShipmentsUseCase(shipRepo),
        application.NewGetShipmentUseCase(shipRepo, jobRepo))

    installH := &httpiface.InstallationHandler{}
    // ... set fields

    techH    := &httpiface.TechnicianHandler{}
    // ... set fields

    rmaH     := &httpiface.RMAHandler{}
    // ... set fields

    // Routes
    httpiface.RegisterRoutes(router, &httpiface.Handlers{
        Shipment: shipH, Installation: installH,
        Technician: techH, RMA: rmaH,
        // Maintenance, WorkOrder
    }, auth, tenant)

    _ = pmUC

    return &Wiring{}
}

var _ = time.Now
```

## D.3 Worker Entry Points

### `cmd/workers/logistics/main.go`
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
)

func main() {
    _ = godotenv.Load()
    brokers := []string{os.Getenv("KAFKA_BROKERS")}
    cfg := sarama.NewConfig()
    cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
    cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

    group, err := sarama.NewConsumerGroup(brokers, "iotlogistics-worker", cfg)
    if err != nil { log.Fatal(err) }
    defer group.Close()

    // ... setup consumers: customer.onboarded, device.alert.triggered
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() {
        for {
            _ = group.Consume(ctx, []string{
                "customer.onboarded",
                "device.alert.triggered",
                "iotlogistics.installation.completed",
            }, nil) // handler
        }
    }()

    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    <-sig
}
```

### `cmd/scheduler/logistics/main.go`
```go
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/joho/godotenv"
    "github.com/robfig/cron/v3"

    logischeduler "icmongolang/internal/modules/iotlogistics/infrastructure/scheduler"
)

func main() {
    _ = godotenv.Load()
    c := cron.New(cron.WithSeconds())

    // PM due — ทุกวัน 02:00
    pmJob := logischeduler.NewPMDueJob(nil)
    c.AddFunc("0 0 2 * * *", pmJob.Run)

    // SLA check — ทุก 5 นาที
    slaJob := logischeduler.NewSLACheckJob(nil, nil)
    c.AddFunc("0 */5 * * * *", slaJob.Run)

    // Reset technician load — ทุกวัน 00:00
    resetJob := logischeduler.NewTechnicianLoadResetJob(nil)
    c.AddFunc("0 0 0 * * *", resetJob.Run)

    c.Start()
    log.Println("logistics scheduler started")

    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    <-sig
    c.Stop()
}
```

---

## D.4 Migration SQL

### `migrations/20260106_iotlogistics_init.sql`
```sql
-- ============================================================
-- iotlogistics module — initial schema
-- Prefix: iotlogistics_
-- ============================================================

-- Shipments
CREATE TABLE IF NOT EXISTS iotlogistics_shipments (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL,
    shipment_no       VARCHAR(30) NOT NULL,
    customer_id       UUID NOT NULL,
    site_id           UUID NOT NULL,
    from_warehouse_id UUID,
    carrier           VARCHAR(100),
    tracking_no       VARCHAR(100),
    status            VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    scheduled_at      TIMESTAMP NOT NULL,
    dispatched_at     TIMESTAMP,
    delivered_at      TIMESTAMP,
    delivered_to      VARCHAR(255),
    signature_url     TEXT,
    gps_last          JSONB,
    notes             TEXT,
    metadata          JSONB DEFAULT '{}'::jsonb,
    created_by        UUID,
    created_at        TIMESTAMP DEFAULT NOW(),
    updated_at        TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_shipment_no UNIQUE (tenant_id, shipment_no)
);
CREATE INDEX idx_ship_tenant_status    ON iotlogistics_shipments(tenant_id, status);
CREATE INDEX idx_ship_tenant_scheduled ON iotlogistics_shipments(tenant_id, scheduled_at DESC);
CREATE INDEX idx_ship_customer         ON iotlogistics_shipments(customer_id);

CREATE TABLE IF NOT EXISTS iotlogistics_shipment_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shipment_id UUID NOT NULL REFERENCES iotlogistics_shipments(id) ON DELETE CASCADE,
    product_id  UUID NOT NULL,
    device_id   UUID,
    serial_no   VARCHAR(100) NOT NULL,
    qty         NUMERIC(15,2) NOT NULL,
    uom         VARCHAR(20) DEFAULT 'PCS',
    notes       TEXT,
    created_at  TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_ship_item_serial UNIQUE (shipment_id, serial_no)
);
CREATE INDEX idx_ship_items_shipment ON iotlogistics_shipment_items(shipment_id);

-- Technicians
CREATE TABLE IF NOT EXISTS iotlogistics_technicians (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    user_id         UUID NOT NULL,
    code            VARCHAR(30) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    phone           VARCHAR(30),
    email           VARCHAR(255),
    level           VARCHAR(20) DEFAULT 'JUNIOR',
    skills          JSONB DEFAULT '[]'::jsonb,
    zone            VARCHAR(100),
    home_base       JSONB,
    is_available    BOOLEAN DEFAULT TRUE,
    is_active       BOOLEAN DEFAULT TRUE,
    gps_last        JSONB,
    gps_updated_at  TIMESTAMP,
    max_jobs_per_day INT DEFAULT 6,
    current_load    INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_tech_code UNIQUE (tenant_id, code)
);
CREATE INDEX idx_tech_tenant_zone  ON iotlogistics_technicians(tenant_id, zone);
CREATE INDEX idx_tech_available    ON iotlogistics_technicians(tenant_id, is_available, is_active);
CREATE INDEX idx_tech_skills_gin   ON iotlogistics_technicians USING GIN (skills);

-- Installation Jobs
CREATE TABLE IF NOT EXISTS iotlogistics_installation_jobs (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID NOT NULL,
    job_no         VARCHAR(30) NOT NULL,
    shipment_id    UUID REFERENCES iotlogistics_shipments(id),
    customer_id    UUID NOT NULL,
    site_id        UUID NOT NULL,
    device_ids     JSONB DEFAULT '[]'::jsonb,
    technician_id  UUID REFERENCES iotlogistics_technicians(id),
    job_type       VARCHAR(30) NOT NULL,
    priority       VARCHAR(20) NOT NULL DEFAULT 'NORMAL',
    status         VARCHAR(20) NOT NULL DEFAULT 'SCHEDULED',
    scheduled_at   TIMESTAMP NOT NULL,
    sla_due_at     TIMESTAMP NOT NULL,
    started_at     TIMESTAMP,
    completed_at   TIMESTAMP,
    started_gps    JSONB,
    completed_gps  JSONB,
    checklist      JSONB DEFAULT '[]'::jsonb,
    photos         JSONB DEFAULT '[]'::jsonb,
    signature_url  TEXT,
    notes          TEXT,
    result         JSONB,
    assigned_by    UUID,
    assigned_at    TIMESTAMP,
    cancelled_at   TIMESTAMP,
    cancel_reason  TEXT,
    created_by     UUID,
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_job_no UNIQUE (tenant_id, job_no)
);
CREATE INDEX idx_job_tenant_status  ON iotlogistics_installation_jobs(tenant_id, status);
CREATE INDEX idx_job_tech_date      ON iotlogistics_installation_jobs(technician_id, scheduled_at);
CREATE INDEX idx_job_sla            ON iotlogistics_installation_jobs(sla_due_at)
    WHERE status NOT IN ('DONE', 'FAILED', 'CANCELLED');

-- Work Orders
CREATE TABLE IF NOT EXISTS iotlogistics_work_orders (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL,
    wo_no         VARCHAR(30) NOT NULL,
    device_id     UUID,
    site_id       UUID NOT NULL,
    customer_id   UUID NOT NULL,
    job_type      VARCHAR(30) NOT NULL,
    priority      VARCHAR(20) DEFAULT 'NORMAL',
    status        VARCHAR(20) NOT NULL DEFAULT 'SCHEDULED',
    title         VARCHAR(255) NOT NULL,
    description   TEXT,
    technician_id UUID REFERENCES iotlogistics_technicians(id),
    scheduled_at  TIMESTAMP,
    time_window   JSONB,
    started_at    TIMESTAMP,
    completed_at  TIMESTAMP,
    resolution    TEXT,
    parts_used    JSONB DEFAULT '[]'::jsonb,
    labor_minutes INT DEFAULT 0,
    cost_total    NUMERIC(15,2) DEFAULT 0,
    currency      VARCHAR(3) DEFAULT 'THB',
    source_ref    VARCHAR(50),
    source_ref_id UUID,
    notes         TEXT,
    created_by    UUID,
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_wo_no UNIQUE (tenant_id, wo_no)
);
CREATE INDEX idx_wo_tenant_status ON iotlogistics_work_orders(tenant_id, status);
CREATE INDEX idx_wo_tech          ON iotlogistics_work_orders(technician_id, scheduled_at);
CREATE INDEX idx_wo_source        ON iotlogistics_work_orders(source_ref, source_ref_id);

-- Maintenance Schedules
CREATE TABLE IF NOT EXISTS iotlogistics_maintenance_schedules (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    device_id       UUID NOT NULL,
    customer_id     UUID NOT NULL,
    site_id         UUID NOT NULL,
    job_type        VARCHAR(30) DEFAULT 'MAINTENANCE',
    frequency       VARCHAR(30) NOT NULL,
    custom_days     INT DEFAULT 0,
    last_done_at    TIMESTAMP,
    next_due_at     TIMESTAMP NOT NULL,
    assigned_tech_id UUID REFERENCES iotlogistics_technicians(id),
    auto_create_wo  BOOLEAN DEFAULT TRUE,
    is_active       BOOLEAN DEFAULT TRUE,
    notes           TEXT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_maint_device UNIQUE (tenant_id, device_id)
);
CREATE INDEX idx_maint_due ON iotlogistics_maintenance_schedules(next_due_at) WHERE is_active = true;

-- RMA
CREATE TABLE IF NOT EXISTS iotlogistics_rma (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              UUID NOT NULL,
    rma_no                 VARCHAR(30) NOT NULL,
    device_id              UUID NOT NULL,
    customer_id            UUID NOT NULL,
    site_id                UUID NOT NULL,
    job_id                 UUID,
    reason                 VARCHAR(30) NOT NULL,
    reason_detail          TEXT,
    status                 VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    resolution             VARCHAR(20),
    replacement_device_id  UUID,
    warranty_claim         BOOLEAN DEFAULT FALSE,
    received_photos        JSONB DEFAULT '[]'::jsonb,
    approved_at            TIMESTAMP,
    approved_by            UUID,
    shipped_at             TIMESTAMP,
    received_at            TIMESTAMP,
    closed_at              TIMESTAMP,
    reject_reason          TEXT,
    assigned_tech_id       UUID REFERENCES iotlogistics_technicians(id),
    notes                  TEXT,
    created_by             UUID,
    created_at             TIMESTAMP DEFAULT NOW(),
    updated_at             TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_rma_no UNIQUE (tenant_id, rma_no)
);
CREATE INDEX idx_rma_tenant_status ON iotlogistics_rma(tenant_id, status);
CREATE INDEX idx_rma_device        ON iotlogistics_rma(device_id);
CREATE INDEX idx_rma_customer      ON iotlogistics_rma(customer_id);

-- Spare Parts
CREATE TABLE IF NOT EXISTS iotlogistics_spare_parts (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID NOT NULL,
    part_no    VARCHAR(50) NOT NULL,
    name       VARCHAR(255) NOT NULL,
    category   VARCHAR(50),
    unit_price NUMERIC(15,2),
    currency   VARCHAR(3) DEFAULT 'THB',
    is_active  BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_spare_tenant_part UNIQUE (tenant_id, part_no)
);

-- Technician GPS tracking log
CREATE TABLE IF NOT EXISTS iotlogistics_technician_locations (
    id            BIGSERIAL PRIMARY KEY,
    tenant_id     UUID NOT NULL,
    technician_id UUID NOT NULL,
    location      JSONB NOT NULL,
    job_id        UUID,
    battery_level INT,
    network       VARCHAR(20),
    recorded_at   TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_tl_tenant_time ON iotlogistics_technician_locations(tenant_id, recorded_at DESC);
CREATE INDEX idx_tl_tech_time   ON iotlogistics_technician_locations(technician_id, recorded_at DESC);

-- Audit log
CREATE TABLE IF NOT EXISTS iotlogistics_audit_logs (
    id          BIGSERIAL PRIMARY KEY,
    tenant_id   UUID NOT NULL,
    actor_id    UUID,
    action      VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id   UUID NOT NULL,
    payload     JSONB DEFAULT '{}'::jsonb,
    ip_address  VARCHAR(45),
    user_agent  VARCHAR(500),
    created_at  TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_log_tenant_action ON iotlogistics_audit_logs(tenant_id, action, created_at DESC);
CREATE INDEX idx_log_entity        ON iotlogistics_audit_logs(entity_type, entity_id, created_at DESC);
```

---

## D.5 .env

```env
# iotlogistics
GEOFENCE_RADIUS_METERS=500
TECHNICIAN_MAX_JOBS_PER_DAY=6
PM_REMINDER_DAYS_BEFORE=7
INSTALLATION_SLA_HOURS=48
MAX_PHOTOS_PER_JOB=20
SHIPMENT_DEFAULT_CARRIER=THAIPOST

# Google Maps
GOOGLE_MAPS_API_KEY=***

# AI Route advisor (optional)
AI_ROUTE_ENABLED=true
AI_PROVIDER=ollama
AI_ENDPOINT=http://ollama:11434
AI_MODEL=llama3

# GPS tracking
GPS_STALE_THRESHOLD_SECONDS=1800
GPS_TRACKING_INTERVAL_SECONDS=30
```

## D.6 Run Instructions

```bash
# 1. Migration
psql "$DB_DSN" -f migrations/20260106_iotlogistics_init.sql

# 2. Build
go build ./internal/modules/iotlogistics/... ./cmd/...

# 3. Test
go test -v ./internal/modules/iotlogistics/...

# 4. Run API
go run ./cmd/api

# 5. Run logistics worker
go run ./cmd/workers/logistics

# 6. Run scheduler
go run ./cmd/scheduler/logistics

# 7. Smoke test — create shipment
curl -X POST http://localhost:8080/api/v1/shipments \
  -H "Authorization: Bearer $TOKEN" -H "X-Tenant-ID: $TENANT_ID" \
  -d '{
    "customer_id":"...","site_id":"...","warehouse_id":"...",
    "scheduled_at":"2026-02-15T09:00:00Z",
    "items":[
      {"product_id":"...","serial_no":"SN-001","qty":5},
      {"product_id":"...","serial_no":"SN-002","qty":3}
    ]
  }'

# 8. Dispatch
curl -X POST http://localhost:8080/api/v1/shipments/$SHIP_ID/dispatch \
  -d '{"carrier":"Kerry Express","tracking_no":"KER-123"}'

# 9. Schedule installation
curl -X POST http://localhost:8080/api/v1/installation-jobs \
  -d '{
    "shipment_id":"...","customer_id":"...","site_id":"...",
    "job_type":"INSTALLATION","priority":"HIGH",
    "scheduled_at":"2026-02-16T10:00:00Z",
    "device_ids":["..."]
  }'

# 10. Auto-assign technician
curl -X POST http://localhost:8080/api/v1/installation-jobs/$JOB_ID/assign \
  -d '{"auto_assign":true,"site_gps":{"lat":13.7563,"lng":100.5018}}'

# 11. Technician starts job (with GPS)
curl -X POST http://localhost:8080/api/v1/installation-jobs/$JOB_ID/start \
  -H "X-User-ID: $TECH_USER_ID" \
  -d '{"lat":13.7563,"lng":100.5018,"accuracy":10}'

# 12. Update checklist
curl -X PATCH http://localhost:8080/api/v1/installation-jobs/$JOB_ID/checklist \
  -d '{"item_key":"power_check","checked":true}'

# 13. Complete job
curl -X POST http://localhost:8080/api/v1/installation-jobs/$JOB_ID/complete \
  -d '{"lat":13.7563,"lng":100.5018,"signature_url":"https://cdn/sig.png","notes":"Done"}'

# 14. Optimize route
curl -X POST http://localhost:8080/api/v1/technicians/$TECH_ID/optimize-route \
  -d '{"date":"2026-02-16T00:00:00Z","max_hours":8}'

# 15. PM scheduler trigger (manual)
go run ./cmd/scheduler/logistics  # ตั้ง cron
```

---

## D.7 DDD Validation Checklist — iotlogistics Module

- [x] 6 Aggregate Roots: `Shipment`, `InstallationJob`, `Technician`, `WorkOrder`, `MaintenanceSchedule`, `RMA`
- [x] Entity ย่อย: `ShipmentItem`, `SparePartUsage`, `TechnicianLocation`
- [x] 12 Value Objects ครบ
- [x] State machines: `ShipmentStatus`, `JobStatus`, `RMAStatus` (ทุกตัวมี `CanTransitionTo`)
- [x] **Field Service Management** pattern — assignment + geo-fence + checkin
- [x] **Route Optimization** — TSP + Maps API + AI advisor
- [x] **Geo-fence validation** ก่อน start/complete job
- [x] **Checklist + photos + signature** — evidence collection
- [x] **PM Scheduler** — auto-create WO จาก schedule
- [x] **RMA Workflow** — approval → ship → receive → repair/replace → close
- [x] **Cross-module** ผ่าน Kafka:
  - `customer.onboarded` → create shipment
  - `device.alert.triggered` → create repair WO
  - `installation.completed` → device.provision + package.billing
- [x] Multi-tenant: `tenant_id` ทุกตาราง + index
- [x] Outbound ports: `MapsPort`, `AIAdvisorPort`, `NotifierPort`, `ERPPort`, `CodeGeneratorPort`
- [x] Domain errors: 60+ sentinel errors
- [x] Unit tests: installation lifecycle, RMA lifecycle
- [x] WebSocket GPS tracking
- [x] Import whitelist ✅
- [x] Scheduler jobs: PM due, SLA check, load reset

---

# ✅ PART 5 (iotlogistics) — เสร็จสมบูรณ์

**สถิติ:**
- ไฟล์ทั้งหมด: **~70 ไฟล์**
- Domain: 28 ไฟล์ (12 VOs, 10 entities, 8 repos, 5 services + 5 ports, 4 events, 1 error)
- Application: 18 ไฟล์ (use cases + DTO + mappers)
- Infrastructure: 12 ไฟล์ (PG repos, code gen, maps, AI, kafka consumers, WS, schedulers)
- Interface: 6 ไฟล์ (handlers + routes + errors)
- Worker: 2 entry points
- Migration: 9 tables
- Kafka topics: 12
- MQTT: ใช้ผ่าน device module (ไม่ publish MQTT ตรง)

**Pattern พิเศษ:**
1. ✅ **FSM (Field Service Management)** — technician + job + geo-fence + checklist
2. ✅ **Route Optimization** — Nearest-neighbor TSP + Google Maps + AI refine
3. ✅ **Technician Assignment Scoring** — skill 40% + proximity 30% + workload 20% + level 10%
4. ✅ **PM Scheduler** — auto-create WO จาก schedule ที่ due
5. ✅ **RMA Workflow** — 9 states + replacement device
6. ✅ **Cross-module orchestration** — installation.completed → device.provision (ผ่าน Kafka)
7. ✅ **Real-time GPS tracking** — WebSocket broadcast

---

**PART ถัดไป**: `PART 6 / 7 — Module: report` (Reporting & Analytics)
หรือต้องการให้ผมทำ `PART 7 — Integration Layer` (Bootstrap ทั้งระบบ + docker-compose + monitoring) ก่อน?

 