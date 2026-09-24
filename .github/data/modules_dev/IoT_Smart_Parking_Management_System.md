### ระบบบริหารจัดการอาคารจอดรถอัจฉริยะ(Smart Parking Management System) Module

> **เป้าหมาย:**
> ที่สมบูรณ์ด้วย Clean Architecture คือ แนวทางการออกแบบซอฟต์แวร์ (Software Design Approach) ที่คิดค้นโดย Robert C. Martin (หรือ Uncle Bob) ซึ่งจัดระเบียบโครงสร้างโค้ดด้วยการแบ่งออกเป็นชั้น ๆ (Layers) เพื่อแยกส่วนตรรกะทางธุรกิจ (Business Rules) ออกจากเทคโนโลยีภายนอก เช่น ฐานข้อมูล (Database) หรือ UI  
> DDD  Domain-Driven Design คือ แนวทางและปรัชญาในการออกแบบซอฟต์แวร์ที่เน้นจำลองโครงสร้างและตรรกะของโค้ดให้สอดคล้องกับ "โดเมนธุรกิจ" (Business Domain) หรือปัญหาที่ซับซ้อนของกิจการนั้น ๆ อย่างแท้จริง  
> **เหมาะสำหรับ:** นักพัฒนาที่ต้องการระบบ Clean Architecture + DDD Domain-Driven Design
> **รายละเอียด**
> - Kafka Consumer Group (หลัก)
> - Elasticsearch Bulk Indexer
> - WebSocket Broadcaster
> - LLM Consumer (เรียก LLM แบบ Async)
> - Embedding Consumer (สร้างเวกเตอร์)
> - Blockchain Consumer (บันทึกข้อมูลลง Blockchain)
> - QR Code & Payment System
> - การติดตั้งและใช้งาน
> - Business Model
> - Mobile App
> - Prompt สำหรับการขยายระบบ

## สารบัญ
1. [ภาพรวมระบบ](#1-ภาพรวมระบบ)
2. [โครงสร้าง Module](#2-โครงสร้าง-module)
3. [Domain Layer](#3-domain-layer)
   - 3.1 Entities
   - 3.2 Value Objects
   - 3.3 Repository Interfaces
   - 3.4 Domain Services
   - 3.5 Domain Errors
4. [Application Layer](#4-application-layer)
   - 4.1 Use Cases
   - 4.2 DTOs
5. [Infrastructure Layer](#5-infrastructure-layer)
   - 5.1 Repository Implementations
   - 5.2 Hardware Controllers
   - 5.3 Payment Gateway Integration
   - 5.4 MQTT/Kafka Implementation
   - 5.5 QR Code Generation
6. [Interface Layer](#6-interface-layer)
   - 6.1 HTTP Handlers
   - 6.2 Routes
   - 6.3 Middleware
   - 6.4 WebSocket Handlers
7. [Database Migrations](#7-database-migrations)
8. [Workflow Diagram](#8-workflow-diagram)
9. [System Flow](#9-system-flow)
10. [การติดตั้งและใช้งาน](#10-การติดตั้งและใช้งาน)
11. [Business Model](#11-business-model)
12. [Prompt สำหรับการขยายระบบ](#12-prompt-สำหรับการขยายระบบ)
13. [ภาคผนวก](#13-ภาคผนวก)

---

## 1. ภาพรวมระบบ

### 1.1 ระบบ Smart Parking Management

ระบบบริหารจัดการอาคารจอดรถอัจฉริยะ ครบวงจร รองรับ:

- **อาคารจอดรถหลายชั้น** - รองรับหลายอาคาร หลายชั้น หลายโซน
- **ระบบลิฟต์จอดรถ** - ควบคุมลิฟต์ขนย้ายรถอัตโนมัติ
- **ระบบชาร์จรถไฟฟ้า (EV Charging)** - รองรับการชาร์จ EV
- **RFID & กล้อง CCTV AI** - ตรวจจับและจดจำป้ายทะเบียน
- **ระบบสแกนป้ายทะเบียน (LPR)** - อ่านป้ายทะเบียนอัตโนมัติ
- **QR Code Scan & Payment** - ชำระเงินผ่าน QR Code
- **ระบบนำทาง (Parking Guidance)** - แสดงที่จอดว่าง/ไม่ว่าง
- **ระบบไฟแสดงสถานะ** - สีเขียว (ว่าง) / สีแดง (ไม่ว่าง)
- **ระบบวิเคราะห์สถิติ** - รายชั่วโมง/วัน/สัปดาห์/เดือน/ปี/ไตรมาส
- **ระบบพิมพ์ใบเสร็จ/บัตรจอดรถ/ใบแจ้งหนี้** - รองรับหลายรูปแบบ
- **ระบบจัดการส่วนลด** - สแกนส่วนลดผ่าน QR Code
- **ระบบแจ้งเตือนเหตุฉุกเฉิน** - ไฟแสดงสถานะและแจ้งเตือน
- **ระบบบำรุงรักษาอุปกรณ์** - ตารางบำรุงรักษาและแจ้งเตือน
- **ระบบป้องกันการกระแทก** - ระบบกันกระแทกอัตโนมัติ
- **ระบบสลับโหมด** - Auto / Manual Mode

### 1.2 โครงสร้างระบบ

```
┌──────────────────────────────────────────────────────────────────────────────────────┐
│                         Smart Parking Management System                              │
├──────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                      │
│  ┌──────────────────────────────────────────────────────────────────────────────────┐ │
│  │                           Building A (5 Floors)                                  │ │
│  │  ┌───────────────────────────────────────────────────────────────────────────┐ │ │
│  │  │  Floor 1 - Zone A1 (ลานจอดทั่วไป)                                       │ │ │
│  │  │  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐      │ │ │
│  │  │  │ 🟢 A1-01│ │ 🔴 A1-02│ │ 🟢 A1-03│ │ 🔴 A1-04│ │ 🟢 A1-05│ │ 🔴 A1-06│ │ 🟢 A1-07│ │ │ │
│  │  │  └──────┘ └──────┘ └──────┘ └──────┘ └──────┘ └──────┘ └──────┘      │ │ │
│  │  └───────────────────────────────────────────────────────────────────────────┘ │ │
│  │                                                                                 │ │
│  │  ┌───────────────────────────────────────────────────────────────────────────┐ │ │
│  │  │  Floor 2 - Zone A2 (EV Charging + EV Parking)                           │ │ │
│  │  │  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐      │ │ │
│  │  │  │ ⚡🔴  │ │ ⚡🟢  │ │ ⚡🟢  │ │ ⚡🔴  │ │ ⚡🟢  │ │ 🔴    │ │ 🟢    │      │ │ │
│  │  │  └──────┘ └──────┘ └──────┘ └──────┘ └──────┘ └──────┘ └──────┘      │ │ │
│  │  └───────────────────────────────────────────────────────────────────────────┘ │ │
│  │                                                                                 │ │
│  │  ┌───────────────────────────────────────────────────────────────────────────┐ │ │
│  │  │  Floor 3 - Zone A3 (VIP / Monthly Rental)                               │ │ │
│  │  │  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐      │ │ │
│  │  │  │ 🔴    │ │ 🟢    │ │ 🟢    │ │ 🔴    │ │ 🟢    │ │ 🟢    │ │ 🔴    │      │ │ │
│  │  │  └──────┘ └──────┘ └──────┘ └──────┘ └──────┘ └──────┘ └──────┘      │ │ │
│  │  └───────────────────────────────────────────────────────────────────────────┘ │ │
│  │                                                                                 │ │
│  │  ┌───────────────────────────────────────────────────────────────────────────┐ │ │
│  │  │  Floor 4-5 - Zone A4 (Automated Parking Lift System)                    │ │ │
│  │  │                    ┌────────────┐                                       │ │ │
│  │  │                    │  Elevator  │                                       │ │ │
│  │  │                    │  Control   │                                       │ │ │
│  │  │                    │  System    │                                       │ │ │
│  │  │                    └────────────┘                                       │ │ │
│  │  └───────────────────────────────────────────────────────────────────────────┘ │ │
│  └──────────────────────────────────────────────────────────────────────────────────┘ │
│                                                                                      │
│  ┌──────────────────────────────────────────────────────────────────────────────────┐ │
│  │                        Control Center / Lobby                                      │ │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐              │ │
│  │  │  CCTV      │  │  LPR       │  │  RFID      │  │  BOOM      │              │ │
│  │  │  Monitor   │  │  System    │  │  Reader    │  │  Barrier   │              │ │
│  │  └────────────┘  └────────────┘  └────────────┘  └────────────┘              │ │
│  │                                                                                 │ │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐              │ │
│  │  │  Payment   │  │  QR Code   │  │  Ticket    │  │  EV       │              │ │
│  │  │  Kiosk     │  │  Scanner   │  │  Printer   │  │  Charger  │              │ │
│  │  └────────────┘  └────────────┘  └────────────┘  └────────────┘              │ │
│  │                                                                                 │ │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐              │ │
│  │  │  Invoice   │  │  Receipt   │  │  Parking   │  │  Emergency │              │ │
│  │  │  Printer   │  │  Printer   │  │  Ticket    │  │  Light     │              │ │
│  │  │  (ใบแจ้งหนี้)│  │  (ใบเสร็จ)  │  │  Printer   │  │  System    │              │ │
│  │  └────────────┘  └────────────┘  └────────────┘  └────────────┘              │ │
│  └──────────────────────────────────────────────────────────────────────────────────┘ │
│                                                                                      │
└──────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. โครงสร้าง Module

```
internal/modules/parking/
│
├── domain/                                    # 🏛️ DOMAIN LAYER
│   ├── entity/
│   │   ├── building.go                        # Building Aggregate Root
│   │   ├── floor.go                           # Floor Entity
│   │   ├── parking_space.go                   # ParkingSpace Entity
│   │   ├── vehicle.go                         # Vehicle Aggregate Root
│   │   ├── parking_ticket.go                  # ParkingTicket Aggregate Root
│   │   ├── parking_rate.go                    # ParkingRate Entity
│   │   ├── ev_charging.go                     # EVCharger & EVChargingSession
│   │   ├── lift.go                            # Lift & LiftOperation
│   │   ├── surveillance.go                    # Camera & LPREvent
│   │   ├── payment.go                         # Payment Entity
│   │   ├── discount.go                        # Discount & DiscountUsage
│   │   ├── printing.go                        # PrintDocument & Receipts
│   │   ├── maintenance.go                     # Equipment & MaintenanceLog
│   │   ├── emergency.go                       # EmergencyEvent & SystemMode
│   │   ├── statistics.go                      # ParkingStatistics
│   │   └── member.go                          # Member & Loyalty
│   │
│   ├── value_object/
│   │   ├── plate_no.go                        # PlateNo Value Object
│   │   ├── ticket_no.go                       # TicketNo Value Object
│   │   ├── space_status.go                    # SpaceStatus Value Object
│   │   ├── money.go                           # Money Value Object
│   │   ├── duration.go                        # Duration Value Object
│   │   ├── qr_code.go                         # QRCode Value Object
│   │   └── location.go                        # Location Value Object
│   │
│   ├── repository/
│   │   ├── building_repository.go             # Interface
│   │   ├── floor_repository.go                # Interface
│   │   ├── parking_space_repository.go        # Interface
│   │   ├── vehicle_repository.go              # Interface
│   │   ├── parking_ticket_repository.go       # Interface
│   │   ├── parking_rate_repository.go         # Interface
│   │   ├── ev_charging_repository.go          # Interface
│   │   ├── lift_repository.go                 # Interface
│   │   ├── payment_repository.go              # Interface
│   │   ├── discount_repository.go             # Interface
│   │   ├── printing_repository.go             # Interface
│   │   ├── maintenance_repository.go          # Interface
│   │   ├── statistics_repository.go           # Interface
│   │   └── member_repository.go               # Interface
│   │
│   ├── service/
│   │   ├── parking_service.go                 # Domain Service
│   │   ├── pricing_service.go                 # Pricing Domain Service
│   │   ├── space_allocation_service.go        # Space Allocation Service
│   │   ├── validation_service.go              # Validation Service
│   │   └── notification_service.go            # Notification Interface
│   │
│   └── errors/
│       └── errors.go                          # Domain Errors
│
├── application/                               # 🎯 APPLICATION LAYER
│   ├── entry_vehicle.go                       # EntryVehicle UseCase
│   ├── exit_vehicle.go                        # ExitVehicle UseCase
│   ├── calculate_fee.go                       # CalculateFee UseCase
│   ├── process_payment.go                     # ProcessPayment UseCase
│   ├── generate_qr_payment.go                 # GenerateQRPayment UseCase
│   ├── print_ticket.go                        # PrintTicket UseCase
│   ├── print_receipt.go                       # PrintReceipt UseCase
│   ├── print_invoice.go                       # PrintInvoice UseCase
│   ├── apply_discount.go                      # ApplyDiscount UseCase
│   ├── manage_maintenance.go                  # ManageMaintenance UseCase
│   ├── emergency_mode.go                      # EmergencyMode UseCase
│   ├── generate_report.go                     # GenerateReport UseCase
│   ├── get_available_spaces.go                # GetAvailableSpaces UseCase
│   ├── get_statistics.go                      # GetStatistics UseCase
│   ├── member_register.go                     # MemberRegister UseCase
│   ├── redeem_points.go                       # RedeemPoints UseCase
│   └── dto.go                                 # Request/Response DTOs
│
├── infrastructure/                            # 🔧 INFRASTRUCTURE LAYER
│   ├── persistence/
│   │   ├── postgres/
│   │   │   ├── building_repo_impl.go          # PostgreSQL Implementation
│   │   │   ├── floor_repo_impl.go
│   │   │   ├── parking_space_repo_impl.go
│   │   │   ├── vehicle_repo_impl.go
│   │   │   ├── parking_ticket_repo_impl.go
│   │   │   ├── parking_rate_repo_impl.go
│   │   │   ├── ev_charging_repo_impl.go
│   │   │   ├── lift_repo_impl.go
│   │   │   ├── payment_repo_impl.go
│   │   │   ├── discount_repo_impl.go
│   │   │   ├── printing_repo_impl.go
│   │   │   ├── maintenance_repo_impl.go
│   │   │   ├── statistics_repo_impl.go
│   │   │   ├── member_repo_impl.go
│   │   │   └── models.go                      # GORM Models
│   │   │
│   │   └── redis/
│   │       ├── cache_repo_impl.go
│   │       ├── session_repo_impl.go
│   │       └── realtime_status_repo.go
│   │
│   ├── hardware/
│   │   ├── led_controller.go                  # LED Control
│   │   ├── boom_controller.go                 # Boom Barrier Control
│   │   ├── lift_controller.go                 # Lift Control
│   │   ├── anti_crash.go                      # Anti-Crash System
│   │   ├── mode_switch.go                     # Auto/Manual Mode
│   │   ├── ev_charger_controller.go           # EV Charger Control
│   │   ├── rfid_reader.go                     # RFID Reader
│   │   └── lpr_service.go                     # LPR Service
│   │
│   ├── payment/
│   │   ├── qr_payment.go                      # QR Payment Generation
│   │   ├── cash_payment.go                    # Cash Payment Processing
│   │   ├── card_payment.go                    # Card Payment Processing
│   │   └── payment_gateway.go                 # Payment Gateway Integration
│   │
│   ├── printing/
│   │   ├── escpos_client.go                   # ESC/POS Protocol
│   │   ├── receipt_printer.go                 # Receipt Printing
│   │   ├── invoice_printer.go                 # Invoice Printing
│   │   ├── ticket_printer.go                  # Ticket Printing
│   │   └── template_engine.go                 # Print Template Engine
│   │
│   ├── messaging/
│   │   ├── kafka/
│   │   │   ├── producer.go                    # Kafka Producer
│   │   │   ├── consumer_group.go              # Kafka Consumer Group
│   │   │   ├── event_types.go                 # Event Definitions
│   │   │   └── handlers.go                    # Event Handlers
│   │   │
│   │   └── mqtt/
│   │       └── mqtt_client.go                 # MQTT Client
│   │
│   ├── ai/
│   │   ├── llm_consumer.go                    # LLM Consumer (Async)
│   │   ├── embedding_consumer.go              # Embedding Consumer
│   │   ├── anomaly_detection.go               # Anomaly Detection
│   │   └── plate_recognition.go               # Plate Recognition
│   │
│   ├── blockchain/
│   │   ├── blockchain_consumer.go             # Blockchain Consumer
│   │   ├── smart_contract.go                  # Smart Contract Integration
│   │   └── transaction_logger.go              # Transaction Logger
│   │
│   ├── elasticsearch/
│   │   ├── bulk_indexer.go                    # Bulk Indexer
│   │   ├── search_service.go                  # Search Service
│   │   └── mapping.go                         # Index Mapping
│   │
│   ├── websocket/
│   │   ├── broadcaster.go                     # WebSocket Broadcaster
│   │   ├── hub.go                             # WebSocket Hub
│   │   └── client.go                          # WebSocket Client
│   │
│   ├── notification/
│   │   ├── line_notify.go                     # LINE Notify
│   │   ├── email_notify.go                    # Email Notify
│   │   └── sms_notify.go                      # SMS Notify
│   │
│   └── security/
│       ├── jwt_maker.go                       # JWT Implementation
│       ├── bcrypt_hasher.go                   # Bcrypt Implementation
│       └── rate_limit.go                      # Rate Limit Middleware
│
└── interfaces/                                # 🌐 INTERFACE LAYER
    ├── http/
    │   ├── parking_handler.go                 # Parking HTTP Handlers
    │   ├── payment_handler.go                 # Payment HTTP Handlers
    │   ├── printing_handler.go                # Printing HTTP Handlers
    │   ├── admin_handler.go                   # Admin HTTP Handlers
    │   ├── member_handler.go                  # Member HTTP Handlers
    │   ├── dashboard_handler.go               # Dashboard HTTP Handlers
    │   ├── routes.go                          # Route Registration
    │   └── dto.go                             # HTTP DTOs
    │
    ├── websocket/
    │   ├── realtime_handler.go                # Real-time WebSocket
    │   └── dashboard_handler.go               # Dashboard WebSocket
    │
    └── middleware/
        ├── auth.go                            # Auth Middleware
        ├── cors.go                            # CORS Middleware
        ├── logging.go                         # Logging Middleware
        ├── rate_limit.go                      # Rate Limit Middleware
        └── role.go                            # Role-based Middleware
```

---

## 3. Domain Layer

### 3.1 Entities

#### 3.1.1 Building - Aggregate Root

```go
// domain/entity/building.go
package entity

import (
    "errors"
    "time"
    "github.com/yourproject/internal/modules/parking/domain/value_object"
)

type Building struct {
    ID          int64              `json:"id"`
    Code        value_object.Code  `json:"code"`
    Name        string             `json:"name"`
    Address     string             `json:"address"`
    TotalFloors int                `json:"total_floors"`
    TotalSpaces int                `json:"total_spaces"`
    HasEVCharger bool              `json:"has_ev_charger"`
    HasLift     bool               `json:"has_lift"`
    Status      BuildingStatus     `json:"status"`
    Floors      []*Floor           `json:"floors,omitempty"`
    CreatedAt   time.Time          `json:"created_at"`
    UpdatedAt   time.Time          `json:"updated_at"`
}

type BuildingStatus string

const (
    BuildingStatusActive     BuildingStatus = "active"
    BuildingStatusInactive   BuildingStatus = "inactive"
    BuildingStatusMaintenance BuildingStatus = "maintenance"
)

func NewBuilding(code, name, address string, totalFloors int, hasEV, hasLift bool) (*Building, error) {
    if code == "" {
        return nil, errors.New("building code is required")
    }
    if name == "" {
        return nil, errors.New("building name is required")
    }
    if totalFloors < 1 {
        return nil, errors.New("building must have at least 1 floor")
    }

    return &Building{
        Code:        value_object.Code{Value: code},
        Name:        name,
        Address:     address,
        TotalFloors: totalFloors,
        HasEVCharger: hasEV,
        HasLift:     hasLift,
        Status:      BuildingStatusActive,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }, nil
}

func (b *Building) AddFloor(floor *Floor) error {
    if floor == nil {
        return errors.New("floor cannot be nil")
    }
    if len(b.Floors) >= b.TotalFloors {
        return errors.New("building has reached maximum floors")
    }
    b.Floors = append(b.Floors, floor)
    b.UpdatedAt = time.Now()
    return nil
}

func (b *Building) UpdateStatus(status BuildingStatus) {
    b.Status = status
    b.UpdatedAt = time.Now()
}
```

#### 3.1.2 Floor - Entity

```go
// domain/entity/floor.go
package entity

import (
    "errors"
    "time"
)

type Floor struct {
    ID          int64     `json:"id"`
    BuildingID  int64     `json:"building_id"`
    FloorNo     int       `json:"floor_no"`
    Name        string    `json:"name"`
    ZoneCode    string    `json:"zone_code"`
    ZoneType    ZoneType  `json:"zone_type"`
    TotalSpaces int       `json:"total_spaces"`
    HasEVCharger bool     `json:"has_ev_charger"`
    HasLift     bool      `json:"has_lift"`
    Spaces      []*ParkingSpace `json:"spaces,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type ZoneType string

const (
    ZoneTypeGeneral    ZoneType = "general"
    ZoneTypeEV         ZoneType = "ev_charging"
    ZoneTypeVIP        ZoneType = "vip"
    ZoneTypeMonthly    ZoneType = "monthly"
    ZoneTypeHandicap   ZoneType = "handicap"
    ZoneTypeAutomated  ZoneType = "automated"
)

func NewFloor(buildingID int64, floorNo int, zoneCode string, zoneType ZoneType) (*Floor, error) {
    if buildingID <= 0 {
        return nil, errors.New("building_id is required")
    }
    if floorNo < 1 {
        return nil, errors.New("floor number must be positive")
    }
    if zoneCode == "" {
        return nil, errors.New("zone code is required")
    }

    return &Floor{
        BuildingID: buildingID,
        FloorNo:    floorNo,
        Name:       string(zoneType),
        ZoneCode:   zoneCode,
        ZoneType:   zoneType,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }, nil
}
```

#### 3.1.3 ParkingSpace - Entity

```go
// domain/entity/parking_space.go
package entity

import (
    "errors"
    "time"
    "github.com/yourproject/internal/modules/parking/domain/value_object"
)

type ParkingSpace struct {
    ID           int64                  `json:"id"`
    FloorID      int64                  `json:"floor_id"`
    SpaceNo      string                 `json:"space_no"`
    SpaceType    SpaceType              `json:"space_type"`
    Status       value_object.SpaceStatus `json:"status"`
    IsEVCharger  bool                   `json:"is_ev_charger"`
    IsHandicap   bool                   `json:"is_handicap"`
    Width        float64                `json:"width"`
    Length       float64                `json:"length"`
    Height       float64                `json:"height"`
    ZoneCode     string                 `json:"zone_code"`
    FloorNo      int                    `json:"floor_no"`
    LEDColor     string                 `json:"led_color"`
    OccupiedBy   string                 `json:"occupied_by"`
    EntryTime    *time.Time             `json:"entry_time"`
    ExitTime     *time.Time             `json:"exit_time"`
    CreatedAt    time.Time              `json:"created_at"`
    UpdatedAt    time.Time              `json:"updated_at"`
}

type SpaceType string

const (
    SpaceTypeCar        SpaceType = "car"
    SpaceTypeSUV        SpaceType = "suv"
    SpaceTypeTruck      SpaceType = "truck"
    SpaceTypeVan        SpaceType = "van"
    SpaceTypeMotorcycle SpaceType = "motorcycle"
    SpaceTypeEV         SpaceType = "ev"
)

func NewParkingSpace(floorID int64, spaceNo string, spaceType SpaceType) (*ParkingSpace, error) {
    if floorID <= 0 {
        return nil, errors.New("floor_id is required")
    }
    if spaceNo == "" {
        return nil, errors.New("space number is required")
    }

    return &ParkingSpace{
        FloorID:    floorID,
        SpaceNo:    spaceNo,
        SpaceType:  spaceType,
        Status:     value_object.SpaceStatusAvailable,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }, nil
}

func (s *ParkingSpace) Occupied(vehicleID string) error {
    if s.Status == value_object.SpaceStatusOccupied {
        return errors.New("space is already occupied")
    }
    if s.Status == value_object.SpaceStatusMaintenance {
        return errors.New("space is under maintenance")
    }

    now := time.Now()
    s.Status = value_object.SpaceStatusOccupied
    s.OccupiedBy = vehicleID
    s.EntryTime = &now
    s.LEDColor = "red"
    s.UpdatedAt = now
    return nil
}

func (s *ParkingSpace) Vacate() error {
    if s.Status == value_object.SpaceStatusAvailable {
        return errors.New("space is already available")
    }

    now := time.Now()
    s.Status = value_object.SpaceStatusAvailable
    s.OccupiedBy = ""
    s.ExitTime = &now
    s.LEDColor = "green"
    s.UpdatedAt = now
    return nil
}
```

#### 3.1.4 ParkingTicket - Aggregate Root

```go
// domain/entity/parking_ticket.go
package entity

import (
    "errors"
    "time"
    "github.com/yourproject/internal/modules/parking/domain/value_object"
)

type ParkingTicket struct {
    ID             int64                  `json:"id"`
    TicketNo       value_object.TicketNo  `json:"ticket_no"`
    BuildingID     int64                  `json:"building_id"`
    FloorID        int64                  `json:"floor_id"`
    SpaceID        int64                  `json:"space_id"`
    VehicleID      int64                  `json:"vehicle_id"`
    PlateNo        value_object.PlateNo   `json:"plate_no"`
    VehicleType    string                 `json:"vehicle_type"`
    EntryTime      time.Time              `json:"entry_time"`
    ExitTime       *time.Time             `json:"exit_time"`
    Duration       int                    `json:"duration"`       // in minutes
    RateType       string                 `json:"rate_type"`
    RateAmount     float64                `json:"rate_amount"`
    TotalAmount    float64                `json:"total_amount"`
    Discount       float64                `json:"discount"`
    NetAmount      float64                `json:"net_amount"`
    PaymentStatus  PaymentStatus          `json:"payment_status"`
    PaymentMethod  PaymentMethod          `json:"payment_method"`
    QRCode         string                 `json:"qr_code"`
    Status         TicketStatus           `json:"status"`
    CreatedAt      time.Time              `json:"created_at"`
    UpdatedAt      time.Time              `json:"updated_at"`
}

type PaymentStatus string

const (
    PaymentStatusPending   PaymentStatus = "pending"
    PaymentStatusPaid      PaymentStatus = "paid"
    PaymentStatusFailed    PaymentStatus = "failed"
    PaymentStatusRefunded  PaymentStatus = "refunded"
)

type TicketStatus string

const (
    TicketStatusActive    TicketStatus = "active"
    TicketStatusCompleted TicketStatus = "completed"
    TicketStatusCancelled TicketStatus = "cancelled"
    TicketStatusExpired   TicketStatus = "expired"
)

func NewParkingTicket(
    ticketNo string,
    buildingID, floorID, spaceID, vehicleID int64,
    plateNo string,
    vehicleType string,
) (*ParkingTicket, error) {
    if ticketNo == "" {
        return nil, errors.New("ticket number is required")
    }
    if buildingID <= 0 || floorID <= 0 || spaceID <= 0 || vehicleID <= 0 {
        return nil, errors.New("invalid reference IDs")
    }
    if plateNo == "" {
        return nil, errors.New("plate number is required")
    }

    ticketNoVal, err := value_object.NewTicketNo(ticketNo)
    if err != nil {
        return nil, err
    }

    plateNoVal, err := value_object.NewPlateNo(plateNo)
    if err != nil {
        return nil, err
    }

    return &ParkingTicket{
        TicketNo:      *ticketNoVal,
        BuildingID:    buildingID,
        FloorID:       floorID,
        SpaceID:       spaceID,
        VehicleID:     vehicleID,
        PlateNo:       *plateNoVal,
        VehicleType:   vehicleType,
        EntryTime:     time.Now(),
        PaymentStatus: PaymentStatusPending,
        Status:        TicketStatusActive,
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }, nil
}

func (t *ParkingTicket) Exit() error {
    if t.Status != TicketStatusActive {
        return errors.New("ticket is not active")
    }
    if t.ExitTime != nil {
        return errors.New("ticket already exited")
    }

    now := time.Now()
    t.ExitTime = &now
    t.Duration = int(now.Sub(t.EntryTime).Minutes())
    t.UpdatedAt = now
    return nil
}

func (t *ParkingTicket) CalculateFee(rates []ParkingRate) (float64, error) {
    if t.ExitTime == nil {
        return 0, errors.New("ticket must have exit time to calculate fee")
    }

    duration := t.ExitTime.Sub(t.EntryTime)
    hours := int(duration.Hours())
    minutes := int(duration.Minutes()) % 60

    var rate ParkingRate
    for _, r := range rates {
        if r.VehicleType == t.VehicleType {
            rate = r
            break
        }
    }

    if rate.FirstHour == 0 {
        return 0, errors.New("no rate found for vehicle type")
    }

    fee := rate.FirstHour
    hours--

    if hours > 0 {
        fee += float64(hours) * rate.SubsequentHour
    }

    if minutes > 0 {
        fee += rate.SubsequentHour
    }

    if fee > rate.MaxCharge {
        fee = rate.MaxCharge
    }

    t.TotalAmount = fee
    return fee, nil
}

func (t *ParkingTicket) ApplyDiscount(discountValue float64, discountType string) error {
    if t.PaymentStatus == PaymentStatusPaid {
        return errors.New("ticket already paid")
    }

    var discount float64
    if discountType == "percentage" {
        discount = t.TotalAmount * (discountValue / 100)
    } else {
        discount = discountValue
    }

    if discount > t.TotalAmount {
        discount = t.TotalAmount
    }

    t.Discount = discount
    t.NetAmount = t.TotalAmount - discount
    t.UpdatedAt = time.Now()
    return nil
}
```

#### 3.1.5 Vehicle - Aggregate Root

```go
// domain/entity/vehicle.go
package entity

import (
    "errors"
    "time"
    "github.com/yourproject/internal/modules/parking/domain/value_object"
)

type Vehicle struct {
    ID          int64                  `json:"id"`
    PlateNo     value_object.PlateNo   `json:"plate_no"`
    VehicleType VehicleType            `json:"vehicle_type"`
    Brand       string                 `json:"brand"`
    Model       string                 `json:"model"`
    Color       string                 `json:"color"`
    IsEV        bool                   `json:"is_ev"`
    OwnerName   string                 `json:"owner_name"`
    OwnerPhone  string                 `json:"owner_phone"`
    OwnerEmail  string                 `json:"owner_email"`
    RFIDTag     string                 `json:"rfid_tag"`
    IsVIP       bool                   `json:"is_vip"`
    MonthlyPlan bool                   `json:"monthly_plan"`
    CreatedAt   time.Time              `json:"created_at"`
    UpdatedAt   time.Time              `json:"updated_at"`
}

type VehicleType string

const (
    VehicleTypeCar        VehicleType = "car"
    VehicleTypeSUV        VehicleType = "suv"
    VehicleTypeTruck      VehicleType = "truck"
    VehicleTypeVan        VehicleType = "van"
    VehicleTypeMotorcycle VehicleType = "motorcycle"
    VehicleTypeEV         VehicleType = "ev"
)

func NewVehicle(plateNo, ownerName, ownerPhone string, vehicleType VehicleType) (*Vehicle, error) {
    if plateNo == "" {
        return nil, errors.New("plate number is required")
    }
    if ownerName == "" {
        return nil, errors.New("owner name is required")
    }
    if ownerPhone == "" {
        return nil, errors.New("owner phone is required")
    }

    plateNoVal, err := value_object.NewPlateNo(plateNo)
    if err != nil {
        return nil, err
    }

    return &Vehicle{
        PlateNo:     *plateNoVal,
        VehicleType: vehicleType,
        OwnerName:   ownerName,
        OwnerPhone:  ownerPhone,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }, nil
}
```

#### 3.1.6 Payment - Entity

```go
// domain/entity/payment.go
package entity

import (
    "errors"
    "time"
)

type Payment struct {
    ID          int64        `json:"id"`
    TicketID    int64        `json:"ticket_id"`
    Amount      float64      `json:"amount"`
    Discount    float64      `json:"discount"`
    NetAmount   float64      `json:"net_amount"`
    Method      PaymentMethod `json:"method"`
    Status      PaymentStatus `json:"status"`
    QRCode      string       `json:"qr_code"`
    PaymentDate time.Time    `json:"payment_date"`
    CreatedAt   time.Time    `json:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at"`
}

type PaymentMethod string

const (
    PaymentMethodQR     PaymentMethod = "qr"
    PaymentMethodCash   PaymentMethod = "cash"
    PaymentMethodCard   PaymentMethod = "card"
    PaymentMethodRFID   PaymentMethod = "rfid"
    PaymentMethodPromptPay PaymentMethod = "promptpay"
)

func NewPayment(ticketID int64, amount float64, method PaymentMethod) (*Payment, error) {
    if ticketID <= 0 {
        return nil, errors.New("ticket_id is required")
    }
    if amount <= 0 {
        return nil, errors.New("amount must be positive")
    }

    return &Payment{
        TicketID:    ticketID,
        Amount:      amount,
        Method:      method,
        Status:      PaymentStatusPending,
        PaymentDate: time.Now(),
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }, nil
}

func (p *Payment) Confirm() {
    p.Status = PaymentStatusPaid
    p.UpdatedAt = time.Now()
}

func (p *Payment) Fail() {
    p.Status = PaymentStatusFailed
    p.UpdatedAt = time.Now()
}
```

### 3.2 Value Objects

#### 3.2.1 PlateNo - Value Object

```go
// domain/value_object/plate_no.go
package value_object

import (
    "errors"
    "regexp"
)

type PlateNo struct {
    Value string
}

var plateNoRegex = regexp.MustCompile(`^[ก-ฮ]{2}\s?\d{4}[ก-ฮ]?$`)

func NewPlateNo(value string) (*PlateNo, error) {
    if value == "" {
        return nil, errors.New("plate number cannot be empty")
    }

    // Basic validation for Thai license plate
    if !plateNoRegex.MatchString(value) {
        // Allow any format for flexibility
        // return nil, errors.New("invalid plate number format")
    }

    return &PlateNo{Value: value}, nil
}

func (p PlateNo) String() string {
    return p.Value
}

func (p PlateNo) Equals(other PlateNo) bool {
    return p.Value == other.Value
}
```

#### 3.2.2 TicketNo - Value Object

```go
// domain/value_object/ticket_no.go
package value_object

import (
    "errors"
    "fmt"
    "time"
)

type TicketNo struct {
    Value string
}

func NewTicketNo(value string) (*TicketNo, error) {
    if value == "" {
        return nil, errors.New("ticket number cannot be empty")
    }
    return &TicketNo{Value: value}, nil
}

func GenerateTicketNo() TicketNo {
    return TicketNo{
        Value: fmt.Sprintf("TK-%s-%04d", time.Now().Format("20060102"), 
            time.Now().UnixNano()%10000),
    }
}

func (t TicketNo) String() string {
    return t.Value
}
```

#### 3.2.3 Money - Value Object

```go
// domain/value_object/money.go
package value_object

import (
    "errors"
    "fmt"
)

type Money struct {
    Amount   float64
    Currency string
}

func NewMoney(amount float64, currency string) (*Money, error) {
    if amount < 0 {
        return nil, errors.New("amount cannot be negative")
    }
    if currency == "" {
        currency = "THB"
    }
    return &Money{
        Amount:   amount,
        Currency: currency,
    }, nil
}

func (m Money) Add(other Money) (Money, error) {
    if m.Currency != other.Currency {
        return Money{}, errors.New("currencies must match")
    }
    return Money{
        Amount:   m.Amount + other.Amount,
        Currency: m.Currency,
    }, nil
}

func (m Money) Subtract(other Money) (Money, error) {
    if m.Currency != other.Currency {
        return Money{}, errors.New("currencies must match")
    }
    return Money{
        Amount:   m.Amount - other.Amount,
        Currency: m.Currency,
    }, nil
}

func (m Money) String() string {
    return fmt.Sprintf("%.2f %s", m.Amount, m.Currency)
}
```

#### 3.2.4 SpaceStatus - Value Object

```go
// domain/value_object/space_status.go
package value_object

type SpaceStatus string

const (
    SpaceStatusAvailable  SpaceStatus = "available"
    SpaceStatusOccupied   SpaceStatus = "occupied"
    SpaceStatusReserved   SpaceStatus = "reserved"
    SpaceStatusMaintenance SpaceStatus = "maintenance"
    SpaceStatusEVCharging SpaceStatus = "ev_charging"
)

func (s SpaceStatus) IsAvailable() bool {
    return s == SpaceStatusAvailable
}

func (s SpaceStatus) IsOccupied() bool {
    return s == SpaceStatusOccupied
}
```

### 3.3 Repository Interfaces

#### 3.3.1 ParkingTicketRepository

```go
// domain/repository/parking_ticket_repository.go
package repository

import (
    "context"
    "time"
    "github.com/yourproject/internal/modules/parking/domain/entity"
    "github.com/yourproject/internal/modules/parking/domain/value_object"
)

type ParkingTicketRepository interface {
    Create(ctx context.Context, ticket *entity.ParkingTicket) error
    Update(ctx context.Context, ticket *entity.ParkingTicket) error
    GetByID(ctx context.Context, id int64) (*entity.ParkingTicket, error)
    GetByNo(ctx context.Context, ticketNo string) (*entity.ParkingTicket, error)
    FindByVehicle(ctx context.Context, vehicleID int64) ([]*entity.ParkingTicket, error)
    FindByPlate(ctx context.Context, plateNo string) ([]*entity.ParkingTicket, error)
    FindActive(ctx context.Context) ([]*entity.ParkingTicket, error)
    FindActiveByBuilding(ctx context.Context, buildingID int64) ([]*entity.ParkingTicket, error)
    FindByDateRange(ctx context.Context, start, end time.Time) ([]*entity.ParkingTicket, error)
    FindByStatus(ctx context.Context, status entity.TicketStatus) ([]*entity.ParkingTicket, error)
    CountActive(ctx context.Context) (int64, error)
    CountByBuilding(ctx context.Context, buildingID int64) (int64, error)
    CountByDateRange(ctx context.Context, start, end time.Time) (int64, error)
    Delete(ctx context.Context, id int64) error
}
```

#### 3.3.2 ParkingSpaceRepository

```go
// domain/repository/parking_space_repository.go
package repository

import (
    "context"
    "github.com/yourproject/internal/modules/parking/domain/entity"
    "github.com/yourproject/internal/modules/parking/domain/value_object"
)

type ParkingSpaceRepository interface {
    Create(ctx context.Context, space *entity.ParkingSpace) error
    Update(ctx context.Context, space *entity.ParkingSpace) error
    GetByID(ctx context.Context, id int64) (*entity.ParkingSpace, error)
    GetByNo(ctx context.Context, floorID int64, spaceNo string) (*entity.ParkingSpace, error)
    FindByFloor(ctx context.Context, floorID int64) ([]*entity.ParkingSpace, error)
    FindByBuilding(ctx context.Context, buildingID int64) ([]*entity.ParkingSpace, error)
    FindAvailable(ctx context.Context, buildingID int64, vehicleType string, isEV bool) (*entity.ParkingSpace, error)
    FindAvailableByFloor(ctx context.Context, floorID int64) ([]*entity.ParkingSpace, error)
    FindByStatus(ctx context.Context, status value_object.SpaceStatus) ([]*entity.ParkingSpace, error)
    FindEVSpaces(ctx context.Context, buildingID int64) ([]*entity.ParkingSpace, error)
    CountByStatus(ctx context.Context, status value_object.SpaceStatus) (int64, error)
    CountByBuilding(ctx context.Context, buildingID int64) (int64, error)
    CountAvailableByBuilding(ctx context.Context, buildingID int64) (int64, error)
    Delete(ctx context.Context, id int64) error
}
```

#### 3.3.3 PaymentRepository

```go
// domain/repository/payment_repository.go
package repository

import (
    "context"
    "time"
    "github.com/yourproject/internal/modules/parking/domain/entity"
)

type PaymentRepository interface {
    Create(ctx context.Context, payment *entity.Payment) error
    Update(ctx context.Context, payment *entity.Payment) error
    GetByID(ctx context.Context, id int64) (*entity.Payment, error)
    GetByTicketID(ctx context.Context, ticketID int64) (*entity.Payment, error)
    FindByStatus(ctx context.Context, status entity.PaymentStatus) ([]*entity.Payment, error)
    FindByDateRange(ctx context.Context, start, end time.Time) ([]*entity.Payment, error)
    FindByMethod(ctx context.Context, method entity.PaymentMethod) ([]*entity.Payment, error)
    GetTotalRevenue(ctx context.Context, start, end time.Time) (float64, error)
    GetRevenueByBuilding(ctx context.Context, buildingID int64, start, end time.Time) (float64, error)
    CountByDateRange(ctx context.Context, start, end time.Time) (int64, error)
    Delete(ctx context.Context, id int64) error
}
```

### 3.4 Domain Services

#### 3.4.1 PricingService

```go
// domain/service/pricing_service.go
package service

import (
    "context"
    "errors"
    "time"
    "github.com/yourproject/internal/modules/parking/domain/entity"
    "github.com/yourproject/internal/modules/parking/domain/repository"
)

type PricingService struct {
    rateRepo repository.ParkingRateRepository
}

func NewPricingService(rateRepo repository.ParkingRateRepository) *PricingService {
    return &PricingService{rateRepo: rateRepo}
}

func (s *PricingService) CalculateFee(
    ctx context.Context,
    ticket *entity.ParkingTicket,
    rates []entity.ParkingRate,
) (float64, error) {
    if ticket.ExitTime == nil {
        return 0, errors.New("ticket must have exit time to calculate fee")
    }

    duration := ticket.ExitTime.Sub(ticket.EntryTime)
    hours := int(duration.Hours())
    minutes := int(duration.Minutes()) % 60

    var rate *entity.ParkingRate
    for _, r := range rates {
        if r.VehicleType == ticket.VehicleType {
            rate = &r
            break
        }
    }

    if rate == nil {
        return 0, errors.New("no rate found for vehicle type")
    }

    fee := rate.FirstHour
    hours--

    if hours > 0 {
        fee += float64(hours) * rate.SubsequentHour
    }

    if minutes > 0 {
        fee += rate.SubsequentHour
    }

    if fee > rate.MaxCharge {
        fee = rate.MaxCharge
    }

    return fee, nil
}

func (s *PricingService) GetRatesForVehicleType(
    ctx context.Context,
    buildingID int64,
    vehicleType string,
) ([]entity.ParkingRate, error) {
    return s.rateRepo.FindByBuildingAndType(ctx, buildingID, vehicleType)
}
```

#### 3.4.2 SpaceAllocationService

```go
// domain/service/space_allocation_service.go
package service

import (
    "context"
    "errors"
    "github.com/yourproject/internal/modules/parking/domain/entity"
    "github.com/yourproject/internal/modules/parking/domain/repository"
    "github.com/yourproject/internal/modules/parking/domain/value_object"
)

type SpaceAllocationService struct {
    spaceRepo repository.ParkingSpaceRepository
}

func NewSpaceAllocationService(spaceRepo repository.ParkingSpaceRepository) *SpaceAllocationService {
    return &SpaceAllocationService{spaceRepo: spaceRepo}
}

func (s *SpaceAllocationService) FindAvailableSpace(
    ctx context.Context,
    buildingID int64,
    vehicleType string,
    isEV bool,
) (*entity.ParkingSpace, error) {
    // Find space based on vehicle type and requirements
    space, err := s.spaceRepo.FindAvailable(ctx, buildingID, vehicleType, isEV)
    if err != nil {
        return nil, err
    }
    if space == nil {
        return nil, errors.New("no available parking space")
    }
    return space, nil
}

func (s *SpaceAllocationService) AllocateSpace(
    ctx context.Context,
    spaceID int64,
    vehicleID string,
) error {
    space, err := s.spaceRepo.GetByID(ctx, spaceID)
    if err != nil {
        return err
    }

    if space.Status != value_object.SpaceStatusAvailable {
        return errors.New("space is not available")
    }

    return space.Occupied(vehicleID)
}

func (s *SpaceAllocationService) ReleaseSpace(
    ctx context.Context,
    spaceID int64,
) error {
    space, err := s.spaceRepo.GetByID(ctx, spaceID)
    if err != nil {
        return err
    }

    if space.Status != value_object.SpaceStatusOccupied {
        return errors.New("space is not occupied")
    }

    return space.Vacate()
}
```

### 3.5 Domain Errors

```go
// domain/errors/errors.go
package errors

import "errors"

var (
    // Building errors
    ErrBuildingNotFound      = errors.New("building not found")
    ErrBuildingInactive      = errors.New("building is inactive")
    ErrBuildingFull          = errors.New("building has reached maximum capacity")
    
    // Space errors
    ErrSpaceNotFound         = errors.New("parking space not found")
    ErrSpaceNotAvailable     = errors.New("parking space is not available")
    ErrSpaceOccupied         = errors.New("parking space is already occupied")
    ErrSpaceMaintenance      = errors.New("parking space is under maintenance")
    ErrNoAvailableSpace      = errors.New("no available parking space found")
    
    // Ticket errors
    ErrTicketNotFound        = errors.New("parking ticket not found")
    ErrTicketNotActive       = errors.New("parking ticket is not active")
    ErrTicketAlreadyExited   = errors.New("ticket already exited")
    ErrTicketAlreadyPaid     = errors.New("ticket already paid")
    ErrTicketExpired         = errors.New("ticket has expired")
    
    // Vehicle errors
    ErrVehicleNotFound       = errors.New("vehicle not found")
    ErrVehicleAlreadyExists  = errors.New("vehicle already exists")
    ErrVehicleTypeInvalid    = errors.New("invalid vehicle type")
    
    // Payment errors
    ErrPaymentNotFound       = errors.New("payment not found")
    ErrPaymentAlreadyProcessed = errors.New("payment already processed")
    ErrPaymentFailed         = errors.New("payment failed")
    ErrInsufficientAmount    = errors.New("insufficient amount")
    
    // Discount errors
    ErrDiscountNotFound      = errors.New("discount not found")
    ErrDiscountExpired       = errors.New("discount has expired")
    ErrDiscountInactive      = errors.New("discount is not active")
    ErrDiscountLimitReached  = errors.New("discount usage limit reached")
    
    // Rate errors
    ErrRateNotFound          = errors.New("parking rate not found")
    
    // Validation errors
    ErrInvalidPlateNumber    = errors.New("invalid plate number")
    ErrInvalidTicketNumber   = errors.New("invalid ticket number")
    ErrInvalidAmount         = errors.New("invalid amount")
    
    // Authorization errors
    ErrUnauthorized          = errors.New("unauthorized access")
    ErrForbidden             = errors.New("forbidden access")
    
    // System errors
    ErrSystemModeInvalid     = errors.New("invalid system mode")
    ErrSystemInEmergency     = errors.New("system is in emergency mode")
    ErrSystemInMaintenance   = errors.New("system is in maintenance mode")
    
    // Hardware errors
    ErrHardwareNotConnected  = errors.New("hardware device not connected")
    ErrHardwareError         = errors.New("hardware error occurred")
    
    // Printing errors
    ErrPrinterNotAvailable   = errors.New("printer not available")
    ErrPrintFailed           = errors.New("print job failed")
    
    // EV Charging errors
    ErrEVChargerNotFound     = errors.New("EV charger not found")
    ErrEVChargerNotAvailable = errors.New("EV charger not available")
    ErrEVChargerInUse        = errors.New("EV charger is in use")
)
```

---

## 4. Application Layer

### 4.1 Use Cases

#### 4.1.1 EntryVehicle - Use Case

```go
// application/entry_vehicle.go
package application

import (
    "context"
    "fmt"
    "time"

    "github.com/yourproject/internal/modules/parking/domain/entity"
    "github.com/yourproject/internal/modules/parking/domain/repository"
    "github.com/yourproject/internal/modules/parking/domain/service"
    "github.com/yourproject/internal/modules/parking/domain/value_object"
    "github.com/yourproject/internal/modules/parking/infrastructure/hardware"
    "github.com/yourproject/internal/modules/parking/infrastructure/printing"
    "github.com/yourproject/internal/modules/parking/infrastructure/messaging/kafka"
)

type EntryVehicleUseCase struct {
    ticketRepo          repository.ParkingTicketRepository
    spaceRepo           repository.ParkingSpaceRepository
    vehicleRepo         repository.VehicleRepository
    rateRepo            repository.ParkingRateRepository
    allocationService   *service.SpaceAllocationService
    printer             *printing.TicketPrinter
    ledCtrl             *hardware.LEDController
    boomCtrl            *hardware.BoomController
    lprService          *hardware.LPRService
    kafkaProducer       *kafka.Producer
}

func NewEntryVehicleUseCase(
    ticketRepo repository.ParkingTicketRepository,
    spaceRepo repository.ParkingSpaceRepository,
    vehicleRepo repository.VehicleRepository,
    rateRepo repository.ParkingRateRepository,
    allocationService *service.SpaceAllocationService,
    printer *printing.TicketPrinter,
    ledCtrl *hardware.LEDController,
    boomCtrl *hardware.BoomController,
    lprService *hardware.LPRService,
    kafkaProducer *kafka.Producer,
) *EntryVehicleUseCase {
    return &EntryVehicleUseCase{
        ticketRepo:        ticketRepo,
        spaceRepo:         spaceRepo,
        vehicleRepo:       vehicleRepo,
        rateRepo:          rateRepo,
        allocationService: allocationService,
        printer:           printer,
        ledCtrl:           ledCtrl,
        boomCtrl:          boomCtrl,
        lprService:        lprService,
        kafkaProducer:     kafkaProducer,
    }
}

type EntryVehicleRequest struct {
    BuildingID   int64  `json:"building_id"`
    PlateNo      string `json:"plate_no"`
    VehicleType  string `json:"vehicle_type"`
    IsEV         bool   `json:"is_ev"`
    EntryMethod  string `json:"entry_method"` // lpr, rfid, manual
}

type EntryVehicleResponse struct {
    TicketNo   string    `json:"ticket_no"`
    SpaceNo    string    `json:"space_no"`
    ZoneCode   string    `json:"zone_code"`
    FloorNo    int       `json:"floor_no"`
    EntryTime  time.Time `json:"entry_time"`
    QRCode     string    `json:"qr_code"`
    IsEV       bool      `json:"is_ev"`
    IsVIP      bool      `json:"is_vip"`
}

func (uc *EntryVehicleUseCase) Execute(ctx context.Context, req EntryVehicleRequest) (*EntryVehicleResponse, error) {
    // 1. Validate request
    if req.PlateNo == "" {
        return nil, domainErr.ErrInvalidPlateNumber
    }

    // 2. Find available space
    space, err := uc.allocationService.FindAvailableSpace(ctx, req.BuildingID, req.VehicleType, req.IsEV)
    if err != nil {
        return nil, err
    }

    // 3. Find or create vehicle
    vehicle, err := uc.vehicleRepo.FindByPlate(ctx, req.PlateNo)
    if err != nil {
        vehicle, err = entity.NewVehicle(req.PlateNo, "", "", entity.VehicleType(req.VehicleType))
        if err != nil {
            return nil, err
        }
        vehicle.IsEV = req.IsEV
        if err := uc.vehicleRepo.Create(ctx, vehicle); err != nil {
            return nil, err
        }
    }

    // 4. Check if vehicle is already parked
    activeTickets, err := uc.ticketRepo.FindByVehicle(ctx, vehicle.ID)
    if err == nil {
        for _, ticket := range activeTickets {
            if ticket.Status == entity.TicketStatusActive {
                return nil, fmt.Errorf("vehicle already parked with ticket %s", ticket.TicketNo)
            }
        }
    }

    // 5. Create parking ticket
    ticketNo := value_object.GenerateTicketNo()
    ticket, err := entity.NewParkingTicket(
        ticketNo.String(),
        req.BuildingID,
        space.FloorID,
        space.ID,
        vehicle.ID,
        req.PlateNo,
        req.VehicleType,
    )
    if err != nil {
        return nil, err
    }

    // 6. Save ticket
    if err := uc.ticketRepo.Create(ctx, ticket); err != nil {
        return nil, err
    }

    // 7. Allocate space
    if err := uc.allocationService.AllocateSpace(ctx, space.ID, fmt.Sprintf("%d", vehicle.ID)); err != nil {
        return nil, err
    }

    // 8. Update LED status
    if err := uc.ledCtrl.SetOccupied(space.ID); err != nil {
        // Log error but continue
    }

    // 9. Open boom barrier
    if err := uc.boomCtrl.Open(req.BuildingID); err != nil {
        // Log error but continue
    }

    // 10. Print ticket
    ticketPrint := &entity.ParkingTicketPrint{
        TicketNo:    ticket.TicketNo.String(),
        EntryTime:   ticket.EntryTime,
        PlateNo:     req.PlateNo,
        VehicleType: req.VehicleType,
        ZoneCode:    space.ZoneCode,
        SpaceNo:     space.SpaceNo,
        QRCode:      generateQRCode(ticket.TicketNo.String()),
        Barcode:     ticket.TicketNo.String(),
    }
    if err := uc.printer.PrintTicket(ticketPrint); err != nil {
        // Log error but continue
    }

    // 11. Publish event to Kafka
    event := kafka.EntryEvent{
        TicketNo:   ticket.TicketNo.String(),
        PlateNo:    req.PlateNo,
        SpaceNo:    space.SpaceNo,
        BuildingID: req.BuildingID,
        Timestamp:  time.Now(),
    }
    if err := uc.kafkaProducer.PublishEntryEvent(ctx, event); err != nil {
        // Log error but continue
    }

    // 12. Build response
    return &EntryVehicleResponse{
        TicketNo:   ticket.TicketNo.String(),
        SpaceNo:    space.SpaceNo,
        ZoneCode:   space.ZoneCode,
        FloorNo:    space.FloorNo,
        EntryTime:  ticket.EntryTime,
        QRCode:     ticketPrint.QRCode,
        IsEV:       req.IsEV,
        IsVIP:      vehicle.IsVIP,
    }, nil
}

func generateQRCode(ticketNo string) string {
    return fmt.Sprintf("https://parking.com/ticket/%s", ticketNo)
}
```

#### 4.1.2 ExitVehicle - Use Case

```go
// application/exit_vehicle.go
package application

import (
    "context"
    "fmt"
    "time"

    "github.com/yourproject/internal/modules/parking/domain/entity"
    "github.com/yourproject/internal/modules/parking/domain/repository"
    "github.com/yourproject/internal/modules/parking/domain/service"
    domainErr "github.com/yourproject/internal/modules/parking/domain/errors"
    "github.com/yourproject/internal/modules/parking/infrastructure/hardware"
    "github.com/yourproject/internal/modules/parking/infrastructure/printing"
    "github.com/yourproject/internal/modules/parking/infrastructure/messaging/kafka"
)

type ExitVehicleUseCase struct {
    ticketRepo          repository.ParkingTicketRepository
    spaceRepo           repository.ParkingSpaceRepository
    paymentRepo         repository.PaymentRepository
    rateRepo            repository.ParkingRateRepository
    discountRepo        repository.DiscountRepository
    pricingService      *service.PricingService
    printer             *printing.ReceiptPrinter
    ledCtrl             *hardware.LEDController
    boomCtrl            *hardware.BoomController
    kafkaProducer       *kafka.Producer
}

func NewExitVehicleUseCase(
    ticketRepo repository.ParkingTicketRepository,
    spaceRepo repository.ParkingSpaceRepository,
    paymentRepo repository.PaymentRepository,
    rateRepo repository.ParkingRateRepository,
    discountRepo repository.DiscountRepository,
    pricingService *service.PricingService,
    printer *printing.ReceiptPrinter,
    ledCtrl *hardware.LEDController,
    boomCtrl *hardware.BoomController,
    kafkaProducer *kafka.Producer,
) *ExitVehicleUseCase {
    return &ExitVehicleUseCase{
        ticketRepo:     ticketRepo,
        spaceRepo:      spaceRepo,
        paymentRepo:    paymentRepo,
        rateRepo:       rateRepo,
        discountRepo:   discountRepo,
        pricingService: pricingService,
        printer:        printer,
        ledCtrl:        ledCtrl,
        boomCtrl:       boomCtrl,
        kafkaProducer:  kafkaProducer,
    }
}

type ExitVehicleRequest struct {
    TicketNo     string `json:"ticket_no"`
    PlateNo      string `json:"plate_no"`
    DiscountCode string `json:"discount_code"`
    PaymentMethod string `json:"payment_method"`
}

type ExitVehicleResponse struct {
    TicketNo    string    `json:"ticket_no"`
    PlateNo     string    `json:"plate_no"`
    EntryTime   time.Time `json:"entry_time"`
    ExitTime    time.Time `json:"exit_time"`
    Duration    string    `json:"duration"`
    TotalAmount float64   `json:"total_amount"`
    Discount    float64   `json:"discount"`
    NetAmount   float64   `json:"net_amount"`
    PaymentMethod string  `json:"payment_method"`
    QRCode      string    `json:"qr_code"`
}

func (uc *ExitVehicleUseCase) Execute(ctx context.Context, req ExitVehicleRequest) (*ExitVehicleResponse, error) {
    // 1. Find ticket
    ticket, err := uc.ticketRepo.GetByNo(ctx, req.TicketNo)
    if err != nil {
        return nil, domainErr.ErrTicketNotFound
    }

    // 2. Validate ticket
    if ticket.Status != entity.TicketStatusActive {
        return nil, domainErr.ErrTicketNotActive
    }
    if ticket.ExitTime != nil {
        return nil, domainErr.ErrTicketAlreadyExited
    }

    // 3. Verify plate number
    if ticket.PlateNo.String() != req.PlateNo {
        return nil, fmt.Errorf("plate number does not match ticket")
    }

    // 4. Set exit time and calculate fee
    if err := ticket.Exit(); err != nil {
        return nil, err
    }

    // 5. Get rates and calculate fee
    rates, err := uc.rateRepo.FindByBuildingAndType(ctx, ticket.BuildingID, ticket.VehicleType)
    if err != nil {
        return nil, err
    }

    totalAmount, err := uc.pricingService.CalculateFee(ctx, ticket, rates)
    if err != nil {
        return nil, err
    }

    // 6. Apply discount if provided
    discount := 0.0
    if req.DiscountCode != "" {
        discount, err = uc.applyDiscount(ctx, req.DiscountCode, ticket.ID)
        if err != nil {
            // Log error but continue without discount
        }
    }

    // 7. Create payment
    netAmount := totalAmount - discount
    payment, err := entity.NewPayment(ticket.ID, netAmount, entity.PaymentMethod(req.PaymentMethod))
    if err != nil {
        return nil, err
    }

    // 8. Process payment (simplified)
    payment.Confirm()

    // 9. Update ticket with payment info
    ticket.TotalAmount = totalAmount
    ticket.Discount = discount
    ticket.NetAmount = netAmount
    ticket.PaymentStatus = entity.PaymentStatusPaid
    ticket.PaymentMethod = entity.PaymentMethod(req.PaymentMethod)
    ticket.UpdatedAt = time.Now()

    if err := uc.ticketRepo.Update(ctx, ticket); err != nil {
        return nil, err
    }

    // 10. Save payment
    if err := uc.paymentRepo.Create(ctx, payment); err != nil {
        return nil, err
    }

    // 11. Release parking space
    space, err := uc.spaceRepo.GetByID(ctx, ticket.SpaceID)
    if err != nil {
        return nil, err
    }

    if err := space.Vacate(); err != nil {
        return nil, err
    }

    if err := uc.spaceRepo.Update(ctx, space); err != nil {
        return nil, err
    }

    // 12. Update LED status
    if err := uc.ledCtrl.SetAvailable(space.ID); err != nil {
        // Log error but continue
    }

    // 13. Open boom barrier
    if err := uc.boomCtrl.Open(ticket.BuildingID); err != nil {
        // Log error but continue
    }

    // 14. Print receipt
    receipt := &entity.ReceiptPrint{
        ReceiptNo:    generateReceiptNo(),
        TicketNo:     ticket.TicketNo.String(),
        PlateNo:      ticket.PlateNo.String(),
        EntryTime:    ticket.EntryTime,
        ExitTime:     *ticket.ExitTime,
        Duration:     formatDuration(*ticket.ExitTime, ticket.EntryTime),
        RateType:     ticket.RateType,
        RateAmount:   ticket.RateAmount,
        TotalAmount:  totalAmount,
        Discount:     discount,
        NetAmount:    netAmount,
        PaymentMethod: req.PaymentMethod,
        PaymentDate:  time.Now(),
        QRCode:       generateQRCode(ticket.TicketNo.String()),
    }
    if err := uc.printer.PrintReceipt(receipt); err != nil {
        // Log error but continue
    }

    // 15. Publish exit event to Kafka
    event := kafka.ExitEvent{
        TicketNo:   ticket.TicketNo.String(),
        PlateNo:    ticket.PlateNo.String(),
        Amount:     netAmount,
        BuildingID: ticket.BuildingID,
        Timestamp:  time.Now(),
    }
    if err := uc.kafkaProducer.PublishExitEvent(ctx, event); err != nil {
        // Log error but continue
    }

    // 16. Build response
    return &ExitVehicleResponse{
        TicketNo:     ticket.TicketNo.String(),
        PlateNo:      ticket.PlateNo.String(),
        EntryTime:    ticket.EntryTime,
        ExitTime:     *ticket.ExitTime,
        Duration:     formatDuration(*ticket.ExitTime, ticket.EntryTime),
        TotalAmount:  totalAmount,
        Discount:     discount,
        NetAmount:    netAmount,
        PaymentMethod: req.PaymentMethod,
        QRCode:       receipt.QRCode,
    }, nil
}

func (uc *ExitVehicleUseCase) applyDiscount(ctx context.Context, code string, ticketID int64) (float64, error) {
    discount, err := uc.discountRepo.FindByCode(ctx, code)
    if err != nil {
        return 0, domainErr.ErrDiscountNotFound
    }

    // Check if discount is valid
    if !discount.IsActive {
        return 0, domainErr.ErrDiscountInactive
    }

    if time.Now().After(discount.EndDate) || time.Now().Before(discount.StartDate) {
        return 0, domainErr.ErrDiscountExpired
    }

    if discount.UsedCount >= discount.MaxUses {
        return 0, domainErr.ErrDiscountLimitReached
    }

    // Get ticket to calculate discount amount
    ticket, err := uc.ticketRepo.GetByID(ctx, ticketID)
    if err != nil {
        return 0, err
    }

    var discountAmount float64
    if discount.Type == "percentage" {
        discountAmount = ticket.TotalAmount * (discount.Value / 100)
    } else {
        discountAmount = discount.Value
    }

    if discountAmount > ticket.TotalAmount {
        discountAmount = ticket.TotalAmount
    }

    // Record discount usage
    usage := &entity.DiscountUsage{
        DiscountID: discount.ID,
        TicketID:   ticketID,
        UsedAt:     time.Now(),
    }
    if err := uc.discountRepo.CreateUsage(ctx, usage); err != nil {
        return 0, err
    }

    // Update discount usage count
    discount.UsedCount++
    if err := uc.discountRepo.Update(ctx, discount); err != nil {
        return 0, err
    }

    return discountAmount, nil
}

func formatDuration(exitTime, entryTime time.Time) string {
    d := exitTime.Sub(entryTime)
    hours := int(d.Hours())
    minutes := int(d.Minutes()) % 60
    return fmt.Sprintf("%d ชั่วโมง %d นาที", hours, minutes)
}

func generateReceiptNo() string {
    return fmt.Sprintf("RCP-%s-%04d", time.Now().Format("20060102"), 
        time.Now().UnixNano()%10000)
}
```

### 4.2 DTOs

```go
// application/dto.go
package application

import "time"

// Entry/Exit DTOs
type EntryRequest struct {
    BuildingID   int64  `json:"building_id" binding:"required"`
    PlateNo      string `json:"plate_no" binding:"required"`
    VehicleType  string `json:"vehicle_type" binding:"required"`
    IsEV         bool   `json:"is_ev"`
    EntryMethod  string `json:"entry_method"`
}

type EntryResponse struct {
    TicketNo   string    `json:"ticket_no"`
    SpaceNo    string    `json:"space_no"`
    ZoneCode   string    `json:"zone_code"`
    FloorNo    int       `json:"floor_no"`
    EntryTime  time.Time `json:"entry_time"`
    QRCode     string    `json:"qr_code"`
    IsEV       bool      `json:"is_ev"`
    IsVIP      bool      `json:"is_vip"`
}

type ExitRequest struct {
    TicketNo     string `json:"ticket_no" binding:"required"`
    PlateNo      string `json:"plate_no" binding:"required"`
    DiscountCode string `json:"discount_code"`
    PaymentMethod string `json:"payment_method" binding:"required"`
}

type ExitResponse struct {
    TicketNo     string    `json:"ticket_no"`
    PlateNo      string    `json:"plate_no"`
    EntryTime    time.Time `json:"entry_time"`
    ExitTime     time.Time `json:"exit_time"`
    Duration     string    `json:"duration"`
    TotalAmount  float64   `json:"total_amount"`
    Discount     float64   `json:"discount"`
    NetAmount    float64   `json:"net_amount"`
    PaymentMethod string   `json:"payment_method"`
    QRCode       string    `json:"qr_code"`
}

// Payment DTOs
type PaymentRequest struct {
    TicketNo string  `json:"ticket_no" binding:"required"`
    Amount   float64 `json:"amount" binding:"required"`
    Method   string  `json:"method" binding:"required"`
}

type PaymentResponse struct {
    PaymentID   int64     `json:"payment_id"`
    TicketNo    string    `json:"ticket_no"`
    Amount      float64   `json:"amount"`
    Status      string    `json:"status"`
    QRCode      string    `json:"qr_code"`
    PaymentDate time.Time `json:"payment_date"`
}

type QRPaymentRequest struct {
    TicketNo string  `json:"ticket_no" binding:"required"`
    Amount   float64 `json:"amount" binding:"required"`
    Discount float64 `json:"discount"`
}

type QRPaymentResponse struct {
    QRCode    string `json:"qr_code"`
    PaymentID int64  `json:"payment_id"`
    Amount    float64 `json:"amount"`
    ExpiresIn int    `json:"expires_in"`
}

// Statistics DTOs
type StatisticsRequest struct {
    BuildingID int64  `json:"building_id"`
    Period     string `json:"period" binding:"required"` // daily, weekly, monthly, quarterly, yearly
    StartDate  string `json:"start_date"`
    EndDate    string `json:"end_date"`
}

type StatisticsResponse struct {
    Period         string    `json:"period"`
    Date           time.Time `json:"date"`
    TotalEntries   int       `json:"total_entries"`
    TotalExits     int       `json:"total_exits"`
    TotalSpaces    int       `json:"total_spaces"`
    AvgOccupancy   float64   `json:"avg_occupancy"`
    MaxOccupancy   float64   `json:"max_occupancy"`
    MinOccupancy   float64   `json:"min_occupancy"`
    TotalRevenue   float64   `json:"total_revenue"`
    AvgRevenue     float64   `json:"avg_revenue"`
    EVRevenue      float64   `json:"ev_revenue"`
}

// Maintenance DTOs
type MaintenanceRequest struct {
    EquipmentID int64  `json:"equipment_id" binding:"required"`
    Type        string `json:"type" binding:"required"`
    Description string `json:"description" binding:"required"`
    PerformedBy string `json:"performed_by" binding:"required"`
    Duration    int    `json:"duration"`
    Cost        float64 `json:"cost"`
    Notes       string `json:"notes"`
}

type MaintenanceResponse struct {
    LogID       int64     `json:"log_id"`
    EquipmentID int64     `json:"equipment_id"`
    Type        string    `json:"type"`
    Description string    `json:"description"`
    PerformedBy string    `json:"performed_by"`
    Duration    int       `json:"duration"`
    Cost        float64   `json:"cost"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
}

// Emergency DTOs
type EmergencyRequest struct {
    BuildingID  int64  `json:"building_id" binding:"required"`
    Type        string `json:"type" binding:"required"`
    Severity    string `json:"severity" binding:"required"`
    Description string `json:"description" binding:"required"`
    Location    string `json:"location"`
}

type EmergencyResponse struct {
    EventID     int64     `json:"event_id"`
    Type        string    `json:"type"`
    Severity    string    `json:"severity"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
}

// Mode DTOs
type ModeRequest struct {
    Mode   string `json:"mode" binding:"required"` // auto, manual, emergency, maintenance
    Reason string `json:"reason"`
}

type ModeResponse struct {
    Mode      string    `json:"mode"`
    ChangedBy string    `json:"changed_by"`
    StartedAt time.Time `json:"started_at"`
    IsActive  bool      `json:"is_active"`
}
```

---

## 5. Infrastructure Layer

### 5.1 Repository Implementations

#### 5.1.1 ParkingTicketRepository - PostgreSQL Implementation

```go
// infrastructure/persistence/postgres/parking_ticket_repo_impl.go
package postgres

import (
    "context"
    "time"
    "gorm.io/gorm"
    "github.com/yourproject/internal/modules/parking/domain/entity"
    "github.com/yourproject/internal/modules/parking/domain/repository"
)

type ParkingTicketRepositoryImpl struct {
    db *gorm.DB
}

func NewParkingTicketRepository(db *gorm.DB) repository.ParkingTicketRepository {
    return &ParkingTicketRepositoryImpl{db: db}
}

func (r *ParkingTicketRepositoryImpl) Create(ctx context.Context, ticket *entity.ParkingTicket) error {
    return r.db.WithContext(ctx).Create(ticket).Error
}

func (r *ParkingTicketRepositoryImpl) Update(ctx context.Context, ticket *entity.ParkingTicket) error {
    return r.db.WithContext(ctx).Save(ticket).Error
}

func (r *ParkingTicketRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.ParkingTicket, error) {
    var ticket entity.ParkingTicket
    err := r.db.WithContext(ctx).First(&ticket, id).Error
    if err != nil {
        return nil, err
    }
    return &ticket, nil
}

func (r *ParkingTicketRepositoryImpl) GetByNo(ctx context.Context, ticketNo string) (*entity.ParkingTicket, error) {
    var ticket entity.ParkingTicket
    err := r.db.WithContext(ctx).Where("ticket_no = ?", ticketNo).First(&ticket).Error
    if err != nil {
        return nil, err
    }
    return &ticket, nil
}

func (r *ParkingTicketRepositoryImpl) FindByVehicle(ctx context.Context, vehicleID int64) ([]*entity.ParkingTicket, error) {
    var tickets []*entity.ParkingTicket
    err := r.db.WithContext(ctx).Where("vehicle_id = ?", vehicleID).Find(&tickets).Error
    return tickets, err
}

func (r *ParkingTicketRepositoryImpl) FindByPlate(ctx context.Context, plateNo string) ([]*entity.ParkingTicket, error) {
    var tickets []*entity.ParkingTicket
    err := r.db.WithContext(ctx).Where("plate_no = ?", plateNo).Find(&tickets).Error
    return tickets, err
}

func (r *ParkingTicketRepositoryImpl) FindActive(ctx context.Context) ([]*entity.ParkingTicket, error) {
    var tickets []*entity.ParkingTicket
    err := r.db.WithContext(ctx).Where("status = ?", entity.TicketStatusActive).Find(&tickets).Error
    return tickets, err
}

func (r *ParkingTicketRepositoryImpl) FindActiveByBuilding(ctx context.Context, buildingID int64) ([]*entity.ParkingTicket, error) {
    var tickets []*entity.ParkingTicket
    err := r.db.WithContext(ctx).Where("building_id = ? AND status = ?", buildingID, entity.TicketStatusActive).Find(&tickets).Error
    return tickets, err
}

func (r *ParkingTicketRepositoryImpl) FindByDateRange(ctx context.Context, start, end time.Time) ([]*entity.ParkingTicket, error) {
    var tickets []*entity.ParkingTicket
    err := r.db.WithContext(ctx).Where("created_at BETWEEN ? AND ?", start, end).Find(&tickets).Error
    return tickets, err
}

func (r *ParkingTicketRepositoryImpl) FindByStatus(ctx context.Context, status entity.TicketStatus) ([]*entity.ParkingTicket, error) {
    var tickets []*entity.ParkingTicket
    err := r.db.WithContext(ctx).Where("status = ?", status).Find(&tickets).Error
    return tickets, err
}

func (r *ParkingTicketRepositoryImpl) CountActive(ctx context.Context) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).Model(&entity.ParkingTicket{}).Where("status = ?", entity.TicketStatusActive).Count(&count).Error
    return count, err
}

func (r *ParkingTicketRepositoryImpl) CountByBuilding(ctx context.Context, buildingID int64) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).Model(&entity.ParkingTicket{}).Where("building_id = ?", buildingID).Count(&count).Error
    return count, err
}

func (r *ParkingTicketRepositoryImpl) CountByDateRange(ctx context.Context, start, end time.Time) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).Model(&entity.ParkingTicket{}).Where("created_at BETWEEN ? AND ?", start, end).Count(&count).Error
    return count, err
}

func (r *ParkingTicketRepositoryImpl) Delete(ctx context.Context, id int64) error {
    return r.db.WithContext(ctx).Delete(&entity.ParkingTicket{}, id).Error
}
```

### 5.2 Hardware Controllers

#### 5.2.1 LED Controller

```go
// infrastructure/hardware/led_controller.go
package hardware

import (
    "fmt"
    "sync"
    "time"
)

type LEDController struct {
    spaces map[int64]LEDStatus
    mu     sync.RWMutex
    conn   *SerialConnection // or TCP connection
}

type LEDStatus struct {
    SpaceID int64
    Color   string // green, red, blue, yellow, flashing_red
    Blink   bool
    Pattern string // solid, blink, flash
}

func NewLEDController(conn *SerialConnection) *LEDController {
    return &LEDController{
        spaces: make(map[int64]LEDStatus),
        conn:   conn,
    }
}

func (c *LEDController) UpdateLED(spaceID int64, color string, pattern string) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.spaces[spaceID] = LEDStatus{
        SpaceID: spaceID,
        Color:   color,
        Pattern: pattern,
        Blink:   pattern == "blink" || pattern == "flash",
    }

    return c.sendCommand(spaceID, color, pattern)
}

func (c *LEDController) sendCommand(spaceID int64, color string, pattern string) error {
    command := fmt.Sprintf("LED:%d:%s:%s\n", spaceID, color, pattern)
    _, err := c.conn.Write([]byte(command))
    return err
}

func (c *LEDController) SetAvailable(spaceID int64) error {
    return c.UpdateLED(spaceID, "green", "solid")
}

func (c *LEDController) SetOccupied(spaceID int64) error {
    return c.UpdateLED(spaceID, "red", "solid")
}

func (c *LEDController) SetReserved(spaceID int64) error {
    return c.UpdateLED(spaceID, "blue", "solid")
}

func (c *LEDController) SetMaintenance(spaceID int64) error {
    return c.UpdateLED(spaceID, "yellow", "blink")
}

func (c *LEDController) SetEVCharging(spaceID int64) error {
    return c.UpdateLED(spaceID, "blue", "blink")
}

func (c *LEDController) SetEmergency(spaceID int64) error {
    return c.UpdateLED(spaceID, "red", "flash")
}
```

#### 5.2.2 Boom Barrier Controller

```go
// infrastructure/hardware/boom_controller.go
package hardware

import (
    "fmt"
    "sync"
    "time"
)

type BoomController struct {
    barriers map[int64]BoomStatus
    mu       sync.RWMutex
    conn     *SerialConnection
}

type BoomStatus struct {
    BarrierID int64
    Position  string // up, down, opening, closing
    Status    string // normal, error, blocked
    LastUpdate time.Time
}

func NewBoomController(conn *SerialConnection) *BoomController {
    return &BoomController{
        barriers: make(map[int64]BoomStatus),
        conn:     conn,
    }
}

func (c *BoomController) Open(buildingID int64) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    command := fmt.Sprintf("BOOM:%d:OPEN\n", buildingID)
    _, err := c.conn.Write([]byte(command))
    if err != nil {
        return err
    }

    c.barriers[buildingID] = BoomStatus{
        BarrierID:  buildingID,
        Position:   "opening",
        Status:     "normal",
        LastUpdate: time.Now(),
    }
    return nil
}

func (c *BoomController) Close(buildingID int64) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    command := fmt.Sprintf("BOOM:%d:CLOSE\n", buildingID)
    _, err := c.conn.Write([]byte(command))
    if err != nil {
        return err
    }

    c.barriers[buildingID] = BoomStatus{
        BarrierID:  buildingID,
        Position:   "closing",
        Status:     "normal",
        LastUpdate: time.Now(),
    }
    return nil
}

func (c *BoomController) GetStatus(buildingID int64) (*BoomStatus, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    status, ok := c.barriers[buildingID]
    if !ok {
        return nil, fmt.Errorf("barrier %d not found", buildingID)
    }
    return &status, nil
}
```

#### 5.2.3 LPR Service

```go
// infrastructure/hardware/lpr_service.go
package hardware

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type LPRService struct {
    cameraURL string
    client    *http.Client
}

type LPRResponse struct {
    PlateNo    string  `json:"plate_no"`
    Confidence float64 `json:"confidence"`
    ImageURL   string  `json:"image_url"`
    Timestamp  string  `json:"timestamp"`
}

func NewLPRService(cameraURL string) *LPRService {
    return &LPRService{
        cameraURL: cameraURL,
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (s *LPRService) ScanPlate(ctx context.Context, cameraID int64) (*LPRResponse, error) {
    url := fmt.Sprintf("%s/api/lpr/scan/%d", s.cameraURL, cameraID)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    resp, err := s.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("LPR service returned status %d", resp.StatusCode)
    }

    var result LPRResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &result, nil
}
```

### 5.3 Payment Gateway Integration

```go
// infrastructure/payment/payment_gateway.go
package payment

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type PaymentGateway struct {
    apiURL    string
    apiKey    string
    secretKey string
    client    *http.Client
}

type PaymentRequest struct {
    Amount    float64 `json:"amount"`
    Currency  string  `json:"currency"`
    Reference string  `json:"reference"`
    Method    string  `json:"method"`
    Timestamp int64   `json:"timestamp"`
}

type PaymentResponse struct {
    TransactionID string `json:"transaction_id"`
    Status        string `json:"status"`
    Message       string `json:"message"`
    QRCode        string `json:"qr_code,omitempty"`
}

func NewPaymentGateway(apiURL, apiKey, secretKey string) *PaymentGateway {
    return &PaymentGateway{
        apiURL:    apiURL,
        apiKey:    apiKey,
        secretKey: secretKey,
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (g *PaymentGateway) ProcessPayment(ctx context.Context, req PaymentRequest) (*PaymentResponse, error) {
    req.Timestamp = time.Now().Unix()
    
    body, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }

    httpReq, err := http.NewRequestWithContext(ctx, "POST", g.apiURL+"/api/payment", bytes.NewBuffer(body))
    if err != nil {
        return nil, err
    }

    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("X-API-Key", g.apiKey)
    httpReq.Header.Set("X-Secret-Key", g.secretKey)

    resp, err := g.client.Do(httpReq)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        return nil, fmt.Errorf("payment gateway returned status %d", resp.StatusCode)
    }

    var result PaymentResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &result, nil
}

func (g *PaymentGateway) GenerateQRCode(ctx context.Context, amount float64, reference string) (string, error) {
    req := PaymentRequest{
        Amount:    amount,
        Currency:  "THB",
        Reference: reference,
        Method:    "qr",
    }

    resp, err := g.ProcessPayment(ctx, req)
    if err != nil {
        return "", err
    }

    return resp.QRCode, nil
}
```

### 5.4 MQTT/Kafka Implementation

#### 5.4.1 Kafka Producer

```go
// infrastructure/messaging/kafka/producer.go
package kafka

import (
    "context"
    "encoding/json"
    "time"

    "github.com/IBM/sarama"
)

type Producer struct {
    producer sarama.SyncProducer
    topic    string
}

type EntryEvent struct {
    TicketNo   string    `json:"ticket_no"`
    PlateNo    string    `json:"plate_no"`
    SpaceNo    string    `json:"space_no"`
    BuildingID int64     `json:"building_id"`
    Timestamp  time.Time `json:"timestamp"`
}

type ExitEvent struct {
    TicketNo   string    `json:"ticket_no"`
    PlateNo    string    `json:"plate_no"`
    Amount     float64   `json:"amount"`
    BuildingID int64     `json:"building_id"`
    Timestamp  time.Time `json:"timestamp"`
}

type PaymentEvent struct {
    PaymentID   int64     `json:"payment_id"`
    TicketNo    string    `json:"ticket_no"`
    Amount      float64   `json:"amount"`
    Method      string    `json:"method"`
    Timestamp   time.Time `json:"timestamp"`
}

type SpaceEvent struct {
    SpaceID     int64     `json:"space_id"`
    SpaceNo     string    `json:"space_no"`
    Status      string    `json:"status"`
    BuildingID  int64     `json:"building_id"`
    Timestamp   time.Time `json:"timestamp"`
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Retry.Max = 5
    config.Producer.Return.Successes = true

    producer, err := sarama.NewSyncProducer(brokers, config)
    if err != nil {
        return nil, err
    }

    return &Producer{
        producer: producer,
        topic:    topic,
    }, nil
}

func (p *Producer) PublishEntryEvent(ctx context.Context, event EntryEvent) error {
    return p.publish(ctx, "entry", event)
}

func (p *Producer) PublishExitEvent(ctx context.Context, event ExitEvent) error {
    return p.publish(ctx, "exit", event)
}

func (p *Producer) PublishPaymentEvent(ctx context.Context, event PaymentEvent) error {
    return p.publish(ctx, "payment", event)
}

func (p *Producer) PublishSpaceEvent(ctx context.Context, event SpaceEvent) error {
    return p.publish(ctx, "space", event)
}

func (p *Producer) publish(ctx context.Context, eventType string, data interface{}) error {
    msg, err := json.Marshal(data)
    if err != nil {
        return err
    }

    partition, offset, err := p.producer.SendMessage(&sarama.ProducerMessage{
        Topic: p.topic,
        Key:   sarama.StringEncoder(eventType),
        Value: sarama.ByteEncoder(msg),
        Headers: []sarama.RecordHeader{
            {Key: []byte("event_type"), Value: []byte(eventType)},
            {Key: []byte("timestamp"), Value: []byte(time.Now().Format(time.RFC3339))},
        },
    })
    if err != nil {
        return err
    }

    // Log successful publish
    _ = partition
    _ = offset
    return nil
}

func (p *Producer) Close() error {
    return p.producer.Close()
}
```

#### 5.4.2 Kafka Consumer Group

```go
// infrastructure/messaging/kafka/consumer_group.go
package kafka

import (
    "context"
    "encoding/json"
    "log"
    "sync"
    "time"

    "github.com/IBM/sarama"
)

type ConsumerGroup struct {
    consumer sarama.ConsumerGroup
    handlers map[string]EventHandler
    mu       sync.RWMutex
}

type EventHandler interface {
    Handle(ctx context.Context, eventType string, data []byte) error
}

type ConsumerGroupHandler struct {
    consumer *ConsumerGroup
    ready    chan bool
}

func NewConsumerGroup(brokers []string, groupID string, topics []string) (*ConsumerGroup, error) {
    config := sarama.NewConfig()
    config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
    config.Consumer.Offsets.Initial = sarama.OffsetNewest
    config.Consumer.Return.Errors = true

    consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
    if err != nil {
        return nil, err
    }

    cg := &ConsumerGroup{
        consumer: consumer,
        handlers: make(map[string]EventHandler),
    }

    // Start consuming
    go cg.consume(topics)

    return cg, nil
}

func (cg *ConsumerGroup) RegisterHandler(eventType string, handler EventHandler) {
    cg.mu.Lock()
    defer cg.mu.Unlock()
    cg.handlers[eventType] = handler
}

func (cg *ConsumerGroup) consume(topics []string) {
    handler := &ConsumerGroupHandler{
        consumer: cg,
        ready:    make(chan bool),
    }

    ctx := context.Background()
    for {
        if err := cg.consumer.Consume(ctx, topics, handler); err != nil {
            log.Printf("Error from consumer: %v", err)
        }
        if ctx.Err() != nil {
            return
        }
        handler.ready = make(chan bool)
    }
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
    close(h.ready)
    return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
    return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        eventType := string(msg.Headers[0].Value)
        
        h.consumer.mu.RLock()
        handler, ok := h.consumer.handlers[eventType]
        h.consumer.mu.RUnlock()

        if !ok {
            log.Printf("No handler for event type: %s", eventType)
            session.MarkMessage(msg, "")
            continue
        }

        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        if err := handler.Handle(ctx, eventType, msg.Value); err != nil {
            log.Printf("Error handling event: %v", err)
            cancel()
            session.MarkMessage(msg, "")
            continue
        }
        cancel()

        session.MarkMessage(msg, "")
    }
    return nil
}

func (cg *ConsumerGroup) Close() error {
    return cg.consumer.Close()
}
```

### 5.5 QR Code Generation

```go
// infrastructure/payment/qr_code.go
package payment

import (
    "encoding/base64"
    "fmt"
    "time"

    "github.com/skip2/go-qrcode"
)

type QRCodeGenerator struct {
    size      int
    expirySec int
}

func NewQRCodeGenerator(size int, expirySec int) *QRCodeGenerator {
    return &QRCodeGenerator{
        size:      size,
        expirySec: expirySec,
    }
}

func (g *QRCodeGenerator) GeneratePaymentQR(ticketNo string, amount float64, discount float64) (string, error) {
    // Create payment data
    netAmount := amount - discount
    data := fmt.Sprintf(
        "PTT|%s|%.2f|%.2f|%d",
        ticketNo,
        netAmount,
        discount,
        time.Now().Unix()+int64(g.expirySec),
    )

    // Generate QR code
    qrCode, err := qrcode.Encode(data, qrcode.Medium, g.size)
    if err != nil {
        return "", err
    }

    // Encode as base64
    return base64.StdEncoding.EncodeToString(qrCode), nil
}

func (g *QRCodeGenerator) GenerateTicketQR(ticketNo string) (string, error) {
    data := fmt.Sprintf("https://parking.com/ticket/%s", ticketNo)
    
    qrCode, err := qrcode.Encode(data, qrcode.Medium, g.size)
    if err != nil {
        return "", err
    }

    return base64.StdEncoding.EncodeToString(qrCode), nil
}
```

---

## 6. Interface Layer

### 6.1 HTTP Handlers

#### 6.1.1 ParkingHandler

```go
// interfaces/http/parking_handler.go
package http

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/yourproject/internal/modules/parking/application"
)

type ParkingHandler struct {
    entryUseCase *application.EntryVehicleUseCase
    exitUseCase  *application.ExitVehicleUseCase
}

func NewParkingHandler(
    entryUseCase *application.EntryVehicleUseCase,
    exitUseCase *application.ExitVehicleUseCase,
) *ParkingHandler {
    return &ParkingHandler{
        entryUseCase: entryUseCase,
        exitUseCase:  exitUseCase,
    }
}

func (h *ParkingHandler) EntryVehicle(c *gin.Context) {
    var req application.EntryVehicleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := h.entryUseCase.Execute(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *ParkingHandler) ExitVehicle(c *gin.Context) {
    var req application.ExitVehicleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := h.exitUseCase.Execute(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *ParkingHandler) GetAvailableSpaces(c *gin.Context) {
    buildingID, err := strconv.ParseInt(c.Query("building_id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid building_id"})
        return
    }

    // Implementation
    c.JSON(http.StatusOK, gin.H{"message": "available spaces"})
}

func (h *ParkingHandler) GetTicketInfo(c *gin.Context) {
    ticketNo := c.Param("ticket_no")
    if ticketNo == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_no is required"})
        return
    }

    // Implementation
    c.JSON(http.StatusOK, gin.H{"ticket_no": ticketNo})
}
```

### 6.2 Routes

```go
// interfaces/http/routes.go
package http

import (
    "github.com/gin-gonic/gin"
    "github.com/yourproject/internal/modules/parking/interfaces/middleware"
)

func SetupRoutes(
    r *gin.Engine,
    parkingHandler *ParkingHandler,
    paymentHandler *PaymentHandler,
    printingHandler *PrintingHandler,
    adminHandler *AdminHandler,
    memberHandler *MemberHandler,
) {
    api := r.Group("/api/v1/parking")

    // Public endpoints (no auth required for entry/exit)
    api.POST("/entry", parkingHandler.EntryVehicle)
    api.POST("/exit", parkingHandler.ExitVehicle)
    api.GET("/space/available", parkingHandler.GetAvailableSpaces)
    api.GET("/qr/:ticket_no", parkingHandler.GetTicketInfo)

    // Payment endpoints
    payment := api.Group("/payment")
    {
        payment.POST("/qr", paymentHandler.GenerateQRPayment)
        payment.POST("/confirm", paymentHandler.ConfirmPayment)
        payment.POST("/discount", paymentHandler.ApplyDiscount)
        payment.GET("/rates", paymentHandler.GetRates)
    }

    // Printing endpoints
    printing := api.Group("/print")
    {
        printing.POST("/ticket/:ticket_id", printingHandler.PrintTicket)
        printing.POST("/receipt/:ticket_id", printingHandler.PrintReceipt)
        printing.POST("/invoice/:ticket_id", printingHandler.PrintInvoice)
        printing.GET("/status/:printer_id", printingHandler.GetPrinterStatus)
    }

    // Member endpoints
    member := api.Group("/member")
    {
        member.POST("/register", memberHandler.Register)
        member.POST("/login", memberHandler.Login)
        member.GET("/profile", middleware.AuthMiddleware(), memberHandler.GetProfile)
        member.GET("/points", middleware.AuthMiddleware(), memberHandler.GetPoints)
        member.POST("/redeem", middleware.AuthMiddleware(), memberHandler.RedeemPoints)
    }

    // Admin endpoints (require auth)
    admin := api.Group("/admin")
    admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
    {
        // Building management
        admin.GET("/buildings", adminHandler.GetBuildings)
        admin.POST("/buildings", adminHandler.CreateBuilding)
        admin.PUT("/buildings/:id", adminHandler.UpdateBuilding)
        admin.DELETE("/buildings/:id", adminHandler.DeleteBuilding)

        // Space management
        admin.GET("/spaces", adminHandler.GetSpaces)
        admin.POST("/spaces", adminHandler.CreateSpace)
        admin.PUT("/spaces/:id", adminHandler.UpdateSpace)
        admin.DELETE("/spaces/:id", adminHandler.DeleteSpace)

        // Rate management
        admin.GET("/rates", adminHandler.GetRates)
        admin.POST("/rates", adminHandler.CreateRate)
        admin.PUT("/rates/:id", adminHandler.UpdateRate)

        // Discount management
        admin.GET("/discounts", adminHandler.GetDiscounts)
        admin.POST("/discounts", adminHandler.CreateDiscount)
        admin.PUT("/discounts/:id", adminHandler.UpdateDiscount)
        admin.DELETE("/discounts/:id", adminHandler.DeleteDiscount)

        // Statistics
        admin.GET("/statistics/daily", adminHandler.GetDailyStatistics)
        admin.GET("/statistics/weekly", adminHandler.GetWeeklyStatistics)
        admin.GET("/statistics/monthly", adminHandler.GetMonthlyStatistics)
        admin.GET("/statistics/quarterly", adminHandler.GetQuarterlyStatistics)
        admin.GET("/statistics/yearly", adminHandler.GetYearlyStatistics)

        // Maintenance
        admin.GET("/equipment", adminHandler.GetEquipment)
        admin.POST("/equipment", adminHandler.CreateEquipment)
        admin.PUT("/equipment/:id", adminHandler.UpdateEquipment)
        admin.GET("/maintenance/schedule", adminHandler.GetMaintenanceSchedule)
        admin.POST("/maintenance/log", adminHandler.CreateMaintenanceLog)

        // Emergency
        admin.GET("/emergency/events", adminHandler.GetEmergencyEvents)
        admin.POST("/emergency/trigger", adminHandler.TriggerEmergency)
        admin.POST("/emergency/resolve", adminHandler.ResolveEmergency)

        // System mode
        admin.GET("/mode", adminHandler.GetSystemMode)
        admin.POST("/mode", adminHandler.SetSystemMode)

        // Printer management
        admin.GET("/printers", adminHandler.GetPrinters)
        admin.POST("/printers", adminHandler.CreatePrinter)
        admin.PUT("/printers/:id", adminHandler.UpdatePrinter)
        admin.POST("/printers/:id/test", adminHandler.TestPrinter)
    }

    // WebSocket endpoints
    r.GET("/ws/parking", WebSocketHandler)
}
```

### 6.3 Middleware

```go
// interfaces/middleware/auth.go
package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/yourproject/internal/modules/parking/infrastructure/security"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
            c.Abort()
            return
        }

        token := parts[1]
        claims, err := security.ValidateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("role", claims.Role)
        c.Next()
    }
}

func RoleMiddleware(role string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("role")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{"error": "role not found"})
            c.Abort()
            return
        }

        if userRole != role {
            c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### 6.4 WebSocket Handlers

```go
// interfaces/websocket/realtime_handler.go
package websocket

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type WebSocketHub struct {
    clients    map[*websocket.Conn]bool
    broadcast  chan []byte
    register   chan *websocket.Conn
    unregister chan *websocket.Conn
    mu         sync.RWMutex
}

type WebSocketMessage struct {
    Type    string      `json:"type"`
    Payload interface{} `json:"payload"`
}

var hub *WebSocketHub
var once sync.Once

func GetHub() *WebSocketHub {
    once.Do(func() {
        hub = &WebSocketHub{
            clients:    make(map[*websocket.Conn]bool),
            broadcast:  make(chan []byte),
            register:   make(chan *websocket.Conn),
            unregister: make(chan *websocket.Conn),
        }
        go hub.run()
    })
    return hub
}

func (h *WebSocketHub) run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()
            log.Printf("New client connected. Total clients: %d", len(h.clients))

        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                client.Close()
            }
            h.mu.Unlock()
            log.Printf("Client disconnected. Total clients: %d", len(h.clients))

        case message := <-h.broadcast:
            h.mu.RLock()
            for client := range h.clients {
                if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
                    log.Printf("Error broadcasting to client: %v", err)
                    client.Close()
                    h.mu.RUnlock()
                    h.mu.Lock()
                    delete(h.clients, client)
                    h.mu.Unlock()
                    h.mu.RLock()
                }
            }
            h.mu.RUnlock()
        }
    }
}

func (h *WebSocketHub) Broadcast(eventType string, payload interface{}) {
    msg := WebSocketMessage{
        Type:    eventType,
        Payload: payload,
    }
    data, err := json.Marshal(msg)
    if err != nil {
        log.Printf("Error marshaling message: %v", err)
        return
    }
    h.broadcast <- data
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Error upgrading connection: %v", err)
        return
    }

    hub := GetHub()
    hub.register <- conn

    defer func() {
        hub.unregister <- conn
    }()

    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("WebSocket error: %v", err)
            }
            break
        }
    }
}
```

---

## 7. Database Migrations

```sql
-- migrations/001_create_buildings_table.sql
CREATE TABLE IF NOT EXISTS buildings (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    total_floors INTEGER NOT NULL DEFAULT 1,
    total_spaces INTEGER NOT NULL DEFAULT 0,
    has_ev_charger BOOLEAN DEFAULT FALSE,
    has_lift BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- migrations/002_create_floors_table.sql
CREATE TABLE IF NOT EXISTS floors (
    id BIGSERIAL PRIMARY KEY,
    building_id BIGINT NOT NULL REFERENCES buildings(id) ON DELETE CASCADE,
    floor_no INTEGER NOT NULL,
    name VARCHAR(100),
    zone_code VARCHAR(50) NOT NULL,
    zone_type VARCHAR(50) DEFAULT 'general',
    total_spaces INTEGER NOT NULL DEFAULT 0,
    has_ev_charger BOOLEAN DEFAULT FALSE,
    has_lift BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(building_id, floor_no)
);

-- migrations/003_create_parking_spaces_table.sql
CREATE TABLE IF NOT EXISTS parking_spaces (
    id BIGSERIAL PRIMARY KEY,
    floor_id BIGINT NOT NULL REFERENCES floors(id) ON DELETE CASCADE,
    space_no VARCHAR(20) NOT NULL,
    space_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) DEFAULT 'available',
    is_ev_charger BOOLEAN DEFAULT FALSE,
    is_handicap BOOLEAN DEFAULT FALSE,
    width DECIMAL(5,2),
    length DECIMAL(5,2),
    height DECIMAL(5,2),
    zone_code VARCHAR(50),
    floor_no INTEGER,
    led_color VARCHAR(20) DEFAULT 'green',
    occupied_by VARCHAR(50),
    entry_time TIMESTAMP WITH TIME ZONE,
    exit_time TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(floor_id, space_no)
);

-- migrations/004_create_vehicles_table.sql
CREATE TABLE IF NOT EXISTS vehicles (
    id BIGSERIAL PRIMARY KEY,
    plate_no VARCHAR(20) NOT NULL UNIQUE,
    vehicle_type VARCHAR(50) NOT NULL,
    brand VARCHAR(100),
    model VARCHAR(100),
    color VARCHAR(50),
    is_ev BOOLEAN DEFAULT FALSE,
    owner_name VARCHAR(255),
    owner_phone VARCHAR(20),
    owner_email VARCHAR(255),
    rfid_tag VARCHAR(50),
    is_vip BOOLEAN DEFAULT FALSE,
    monthly_plan BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- migrations/005_create_parking_tickets_table.sql
CREATE TABLE IF NOT EXISTS parking_tickets (
    id BIGSERIAL PRIMARY KEY,
    ticket_no VARCHAR(50) NOT NULL UNIQUE,
    building_id BIGINT NOT NULL REFERENCES buildings(id),
    floor_id BIGINT NOT NULL REFERENCES floors(id),
    space_id BIGINT NOT NULL REFERENCES parking_spaces(id),
    vehicle_id BIGINT NOT NULL REFERENCES vehicles(id),
    plate_no VARCHAR(20) NOT NULL,
    vehicle_type VARCHAR(50) NOT NULL,
    entry_time TIMESTAMP WITH TIME ZONE NOT NULL,
    exit_time TIMESTAMP WITH TIME ZONE,
    duration INTEGER DEFAULT 0,
    rate_type VARCHAR(50),
    rate_amount DECIMAL(10,2) DEFAULT 0,
    total_amount DECIMAL(10,2) DEFAULT 0,
    discount DECIMAL(10,2) DEFAULT 0,
    net_amount DECIMAL(10,2) DEFAULT 0,
    payment_status VARCHAR(20) DEFAULT 'pending',
    payment_method VARCHAR(50),
    qr_code TEXT,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- migrations/006_create_parking_rates_table.sql
CREATE TABLE IF NOT EXISTS parking_rates (
    id BIGSERIAL PRIMARY KEY,
    building_id BIGINT NOT NULL REFERENCES buildings(id),
    vehicle_type VARCHAR(50) NOT NULL,
    rate_type VARCHAR(50) DEFAULT 'standard',
    first_hour DECIMAL(10,2) NOT NULL,
    subsequent_hour DECIMAL(10,2) NOT NULL,
    daily_rate DECIMAL(10,2),
    monthly_rate DECIMAL(10,2),
    event_rate DECIMAL(10,2),
    overstay_fee DECIMAL(10,2),
    max_charge DECIMAL(10,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(building_id, vehicle_type, rate_type)
);

-- migrations/007_create_payments_table.sql
CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,
    ticket_id BIGINT NOT NULL REFERENCES parking_tickets(id),
    amount DECIMAL(10,2) NOT NULL,
    discount DECIMAL(10,2) DEFAULT 0,
    net_amount DECIMAL(10,2) NOT NULL,
    method VARCHAR(50) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    qr_code TEXT,
    payment_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- migrations/008_create_discounts_table.sql
CREATE TABLE IF NOT EXISTS discounts (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL,
    value DECIMAL(10,2) NOT NULL,
    description TEXT,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    max_uses INTEGER DEFAULT 1,
    used_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- migrations/009_create_discount_usages_table.sql
CREATE TABLE IF NOT EXISTS discount_usages (
    id BIGSERIAL PRIMARY KEY,
    discount_id BIGINT NOT NULL REFERENCES discounts(id),
    ticket_id BIGINT NOT NULL REFERENCES parking_tickets(id),
    used_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- migrations/010_create_print_documents_table.sql
CREATE TABLE IF NOT EXISTS print_documents (
    id BIGSERIAL PRIMARY KEY,
    ticket_id BIGINT NOT NULL REFERENCES parking_tickets(id),
    document_type VARCHAR(50) NOT NULL,
    document_no VARCHAR(50) NOT NULL,
    content TEXT,
    qr_code_data TEXT,
    printer_name VARCHAR(100),
    status VARCHAR(20) DEFAULT 'pending',
    printed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- migrations/011_create_equipment_table.sql
CREATE TABLE IF NOT EXISTS equipment (
    id BIGSERIAL PRIMARY KEY,
    building_id BIGINT NOT NULL REFERENCES buildings(id),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    model VARCHAR(100),
    serial_no VARCHAR(100),
    install_date TIMESTAMP WITH TIME ZONE,
    warranty_end TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'active',
    last_maintenance TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- migrations/012_create_maintenance_logs_table.sql
CREATE TABLE IF NOT EXISTS maintenance_logs (
    id BIGSERIAL PRIMARY KEY,
    equipment_id BIGINT NOT NULL REFERENCES equipment(id),
    type VARCHAR(50) NOT NULL,
    description TEXT,
    performed_by VARCHAR(255),
    duration INTEGER DEFAULT 0,
    cost DECIMAL(10,2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'completed',
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- migrations/013_create_emergency_events_table.sql
CREATE TABLE IF NOT EXISTS emergency_events (
    id BIGSERIAL PRIMARY KEY,
    building_id BIGINT NOT NULL REFERENCES buildings(id),
    type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    description TEXT,
    location VARCHAR(255),
    status VARCHAR(20) DEFAULT 'active',
    reported_by VARCHAR(255),
    resolved_by VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP WITH TIME ZONE
);

-- migrations/014_create_system_modes_table.sql
CREATE TABLE IF NOT EXISTS system_modes (
    id BIGSERIAL PRIMARY KEY,
    building_id BIGINT NOT NULL REFERENCES buildings(id),
    mode VARCHAR(50) NOT NULL,
    changed_by VARCHAR(255),
    reason TEXT,
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ended_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- migrations/015_create_members_table.sql
CREATE TABLE IF NOT EXISTS members (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    tier VARCHAR(50) DEFAULT 'silver',
    points INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- migrations/016_create_statistics_table.sql
CREATE TABLE IF NOT EXISTS parking_statistics (
    id BIGSERIAL PRIMARY KEY,
    building_id BIGINT NOT NULL REFERENCES buildings(id),
    period VARCHAR(50) NOT NULL,
    date DATE NOT NULL,
    total_entries INTEGER DEFAULT 0,
    total_exits INTEGER DEFAULT 0,
    total_spaces INTEGER DEFAULT 0,
    avg_occupancy DECIMAL(5,2) DEFAULT 0,
    max_occupancy DECIMAL(5,2) DEFAULT 0,
    min_occupancy DECIMAL(5,2) DEFAULT 0,
    total_revenue DECIMAL(10,2) DEFAULT 0,
    avg_revenue DECIMAL(10,2) DEFAULT 0,
    ev_charging_revenue DECIMAL(10,2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(building_id, period, date)
);

-- Indexes
CREATE INDEX idx_tickets_ticket_no ON parking_tickets(ticket_no);
CREATE INDEX idx_tickets_plate_no ON parking_tickets(plate_no);
CREATE INDEX idx_tickets_status ON parking_tickets(status);
CREATE INDEX idx_tickets_building_id ON parking_tickets(building_id);
CREATE INDEX idx_tickets_created_at ON parking_tickets(created_at);
CREATE INDEX idx_payments_ticket_id ON payments(ticket_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_spaces_status ON parking_spaces(status);
CREATE INDEX idx_spaces_floor_id ON parking_spaces(floor_id);
CREATE INDEX idx_emergency_building_id ON emergency_events(building_id);
CREATE INDEX idx_emergency_status ON emergency_events(status);
```

---

## 8. Workflow Diagram

### 8.1 Vehicle Entry Flow

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                           VEHICLE ENTRY FLOW                                         │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  ┌──────────────┐                                                                   │
│  │  Vehicle     │                                                                   │
│  │  Arrives     │                                                                   │
│  └──────┬───────┘                                                                   │
│         │                                                                           │
│         ▼                                                                           │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐                        │
│  │  LPR Camera  │────▶│  Plate No    │────▶│  Validate    │                        │
│  │  Scan        │     │  Detection   │     │  Plate No    │                        │
│  └──────────────┘     └──────────────┘     └──────┬───────┘                        │
│                                                    │                                │
│                                                    ▼                                │
│                                     ┌──────────────────────────┐                    │
│                                     │  Check Vehicle Status     │                    │
│                                     │  - Is vehicle already?    │                    │
│                                     │  - Is parking allowed?    │                    │
│                                     └──────────┬───────────────┘                    │
│                                                │                                    │
│                                                ▼                                    │
│                                     ┌──────────────────────────┐                    │
│                                     │  Find Available Space    │                    │
│                                     │  - Check building        │                    │
│                                     │  - Check zone            │                    │
│                                     │  - Check EV status       │                    │
│                                     └──────────┬───────────────┘                    │
│                                                │                                    │
│                                                ▼                                    │
│                                     ┌──────────────────────────┐                    │
│                                     │  Create Parking Ticket   │                    │
│                                     │  - Generate Ticket No    │                    │
│                                     │  - Record Entry Time     │                    │
│                                     │  - Assign Space          │                    │
│                                     └──────────┬───────────────┘                    │
│                                                │                                    │
│                                                ▼                                    │
│                              ┌──────────────────────────────────┐                  │
│                              │  Update System State              │                  │
│                              │  - Update Space Status           │                  │
│                              │  - Update LED Color (Red)        │                  │
│                              │  - Open Boom Barrier             │                  │
│                              │  - Print Ticket                  │                  │
│                              └──────────┬───────────────────────┘                  │
│                                         │                                          │
│                                         ▼                                          │
│                              ┌──────────────────────────────────┐                  │
│                              │  Publish Event to Kafka          │                  │
│                              │  - Entry Event                   │                  │
│                              │  - Space Event                   │                  │
│                              └──────────┬───────────────────────┘                  │
│                                         │                                          │
│                                         ▼                                          │
│                              ┌──────────────────────────────────┐                  │
│                              │  Send WebSocket Update           │                  │
│                              │  - Update Dashboard              │                  │
│                              │  - Notify Mobile App             │                  │
│                              └──────────┬───────────────────────┘                  │
│                                         │                                          │
│                                         ▼                                          │
│                              ┌──────────────────────────────────┐                  │
│                              │  Vehicle Enters Parking          │                  │
│                              └──────────────────────────────────┘                  │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 8.2 Vehicle Exit Flow

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                           VEHICLE EXIT FLOW                                          │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  ┌──────────────┐                                                                   │
│  │  Vehicle     │                                                                   │
│  │  Prepares    │                                                                   │
│  │  to Exit     │                                                                   │
│  └──────┬───────┘                                                                   │
│         │                                                                           │
│         ▼                                                                           │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐                        │
│  │  Scan Ticket │────▶│  Validate    │────▶│  Verify      │                        │
│  │  or Plate    │     │  Ticket      │     │  Plate No    │                        │
│  └──────────────┘     └──────────────┘     └──────┬───────┘                        │
│                                                    │                                │
│                                                    ▼                                │
│                                     ┌──────────────────────────┐                    │
│                                     │  Calculate Fee           │                    │
│                                     │  - Calculate Duration    │                    │
│                                     │  - Apply Rate            │                    │
│                                     │  - Apply Discount        │                    │
│                                     └──────────┬───────────────┘                    │
│                                                │                                    │
│                                                ▼                                    │
│                                     ┌──────────────────────────┐                    │
│                                     │  Process Payment         │                    │
│                                     │  - Generate QR Code      │                    │
│                                     │  - Process Payment       │                    │
│                                     │  - Confirm Payment       │                    │
│                                     └──────────┬───────────────┘                    │
│                                                │                                    │
│                                                ▼                                    │
│                              ┌──────────────────────────────────┐                  │
│                              │  Update System State              │                  │
│                              │  - Update Ticket Status          │                  │
│                              │  - Release Space                 │                  │
│                              │  - Update LED Color (Green)      │                  │
│                              │  - Open Boom Barrier             │                  │
│                              │  - Print Receipt                 │                  │
│                              └──────────┬───────────────────────┘                  │
│                                         │                                          │
│                                         ▼                                          │
│                              ┌──────────────────────────────────┐                  │
│                              │  Publish Event to Kafka          │                  │
│                              │  - Exit Event                   │                  │
│                              │  - Payment Event                │                  │
│                              │  - Space Event                  │                  │
│                              └──────────┬───────────────────────┘                  │
│                                         │                                          │
│                                         ▼                                          │
│                              ┌──────────────────────────────────┐                  │
│                              │  Send WebSocket Update           │                  │
│                              │  - Update Dashboard              │                  │
│                              │  - Notify Mobile App             │                  │
│                              └──────────┬───────────────────────┘                  │
│                                         │                                          │
│                                         ▼                                          │
│                              ┌──────────────────────────────────┐                  │
│                              │  Vehicle Exits Parking           │                  │
│                              └──────────────────────────────────┘                  │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 9. System Flow

### 9.1 Complete System Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                         COMPLETE SYSTEM FLOW                                        │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │                          FRONTEND LAYER                                      │   │
│  │  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐          │   │
│  │  │  React Dashboard │  │  Mobile App      │  │  Admin Panel     │          │   │
│  │  │  - Real-time     │  │  - QR Scanner    │  │  - Management    │          │   │
│  │  │  - Charts        │  │  - Payment       │  │  - Reports       │          │   │
│  │  │  - Controls      │  │  - History       │  │  - Settings      │          │   │
│  │  └────────┬─────────┘  └────────┬─────────┘  └────────┬─────────┘          │   │
│  └───────────┼─────────────────────┼─────────────────────┼─────────────────────┘   │
│              │                     │                     │                          │
│              ▼                     ▼                     ▼                          │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │                          API GATEWAY / LOAD BALANCER                        │   │
│  │                    (HTTP/REST, WebSocket, MQTT)                             │   │
│  └──────────────────────────────────┬──────────────────────────────────────────┘   │
│                                     │                                              │
│                                     ▼                                              │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │                          INTERFACE LAYER                                     │   │
│  │  ┌────────────────────┐  ┌────────────────────┐  ┌────────────────────┐     │   │
│  │  │  HTTP Handlers     │  │  WebSocket         │  │  Middleware        │     │   │
│  │  │  - Parking         │  │  - Real-time       │  │  - Auth            │     │   │
│  │  │  - Payment         │  │  - Dashboard       │  │  - Rate Limit      │     │   │
│  │  │  - Admin           │  │  - Notifications   │  │  - Logging         │     │   │
│  │  └────────────────────┘  └────────────────────┘  └────────────────────┘     │   │
│  └──────────────────────────────────┬──────────────────────────────────────────┘   │
│                                     │                                              │
│                                     ▼                                              │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │                          APPLICATION LAYER                                   │   │
│  │  ┌────────────────────┐  ┌────────────────────┐  ┌────────────────────┐     │   │
│  │  │  Use Cases         │  │  Use Cases         │  │  Use Cases         │     │   │
│  │  │  - Entry Vehicle   │  │  - Exit Vehicle    │  │  - Process Payment │     │   │
│  │  │  - Print Ticket    │  │  - Print Receipt   │  │  - Generate Report │     │   │
│  │  └────────────────────┘  └────────────────────┘  └────────────────────┘     │   │
│  └──────────────────────────────────┬──────────────────────────────────────────┘   │
│                                     │                                              │
│                                     ▼                                              │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │                          DOMAIN LAYER                                       │   │
│  │  ┌────────────────────┐  ┌────────────────────┐  ┌────────────────────┐     │   │
│  │  │  Entities          │  │  Value Objects     │  │  Domain Services   │     │   │
│  │  │  - Building        │  │  - Plate No        │  │  - Pricing         │     │   │
│  │  │  - Parking Space   │  │  - Ticket No       │  │  - Allocation      │     │   │
│  │  │  - Ticket          │  │  - Money           │  │  - Validation      │     │   │
│  │  └────────────────────┘  └────────────────────┘  └────────────────────┘     │   │
│  └──────────────────────────────────┬──────────────────────────────────────────┘   │
│                                     │                                              │
│                                     ▼                                              │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │                          INFRASTRUCTURE LAYER                                │   │
│  │                                                                             │   │
│  │  ┌──────────────────────────────────────────────────────────────────────┐   │   │
│  │  │                     DATA PERSISTENCE                                 │   │   │
│  │  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐    │   │   │
│  │  │  │ PostgreSQL │  │  Redis     │  │  InfluxDB  │  │ Elasticsearch│    │   │   │
│  │  │  │ - Tickets  │  │ - Cache    │  │ - Stats    │  │ - Search   │    │   │   │
│  │  │  │ - Spaces   │  │ - Sessions │  │ - Metrics  │  │ - Logs     │    │   │   │
│  │  │  │ - Payments │  │ - Real-time│  │ - Analytics│  │ - Audit    │    │   │   │
│  │  │  └────────────┘  └────────────┘  └────────────┘  └────────────┘    │   │   │
│  │  └──────────────────────────────────────────────────────────────────────┘   │   │
│  │                                                                             │   │
│  │  ┌──────────────────────────────────────────────────────────────────────┐   │   │
│  │  │                     MESSAGING & STREAMING                           │   │   │
│  │  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐    │   │   │
│  │  │  │   Kafka    │  │   MQTT     │  │  WebSocket │  │   gRPC     │    │   │   │
│  │  │  │ - Events   │  │ - Hardware │  │ - Real-time│  │ - Services │    │   │   │
│  │  │  │ - Streams  │  │ - Sensors  │  │ - Updates  │  │ - Internal │    │   │   │
│  │  │  └────────────┘  └────────────┘  └────────────┘  └────────────┘    │   │   │
│  │  └──────────────────────────────────────────────────────────────────────┘   │   │
│  │                                                                             │   │
│  │  ┌──────────────────────────────────────────────────────────────────────┐   │   │
│  │  │                     HARDWARE & DEVICES                               │   │   │
│  │  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐    │   │   │
│  │  │  │ LPR Camera │  │ RFID Reader│  │ Boom Barrier│  │ LED System │    │   │   │
│  │  │  │ EV Charger │  │ Lift System│  │ Printers   │  │ Sensors    │    │   │   │
│  │  │  └────────────┘  └────────────┘  └────────────┘  └────────────┘    │   │   │
│  │  └──────────────────────────────────────────────────────────────────────┘   │   │
│  │                                                                             │   │
│  │  ┌──────────────────────────────────────────────────────────────────────┐   │   │
│  │  │                     EXTERNAL SERVICES                                 │   │   │
│  │  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐    │   │   │
│  │  │  │ Payment    │  │ LLM/AI     │  │ Blockchain │  │ Notification│   │   │   │
│  │  │  │ Gateway    │  │ Services   │  │ Services   │  │ Services   │    │   │   │
│  │  │  └────────────┘  └────────────┘  └────────────┘  └────────────┘    │   │   │
│  │  └──────────────────────────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 10. การติดตั้งและใช้งาน

### 10.1 Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: parking
      POSTGRES_PASSWORD: password
      POSTGRES_DB: parking
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U parking"]
      interval: 10s
      timeout: 5s
      retries: 5

  influxdb:
    image: influxdb:2.7-alpine
    environment:
      INFLUXDB_INIT_MODE: setup
      INFLUXDB_INIT_USERNAME: admin
      INFLUXDB_INIT_PASSWORD: password
      INFLUXDB_INIT_ORG: parking
      INFLUXDB_INIT_BUCKET: parking_stats
      INFLUXDB_INIT_RETENTION: 365d
    ports:
      - "8086:8086"
    volumes:
      - influxdb_data:/var/lib/influxdb2

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.10.0
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    ports:
      - "9200:9200"
    volumes:
      - elasticsearch_data:/usr/share/elasticsearch/data

  kibana:
    image: docker.elastic.co/kibana/kibana:8.10.0
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
    ports:
      - "5601:5601"
    depends_on:
      - elasticsearch

  zookeeper:
    image: confluentinc/cp-zookeeper:latest
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    ports:
      - "2181:2181"

  kafka:
    image: confluentinc/cp-kafka:latest
    depends_on:
      - zookeeper
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    ports:
      - "9092:9092"

  mosquitto:
    image: eclipse-mosquitto:2.0
    ports:
      - "1883:1883"
      - "9001:9001"
    volumes:
      - ./mosquitto/config:/mosquitto/config

  backend:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: parking
      DB_PASSWORD: password
      DB_NAME: parking
      INFLUX_URL: http://influxdb:8086
      INFLUX_TOKEN: admin:password
      INFLUX_ORG: parking
      INFLUX_BUCKET: parking_stats
      MQTT_BROKER: mosquitto:1883
      REDIS_HOST: redis:6379
      KAFKA_BROKERS: kafka:9092
      ELASTICSEARCH_URL: http://elasticsearch:9200
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      kafka:
        condition: service_started
      elasticsearch:
        condition: service_started
      mosquitto:
        condition: service_started
    volumes:
      - ./uploads:/app/uploads
      - /dev/usb:/dev/usb  # For printers

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    depends_on:
      - backend

volumes:
  postgres_data:
  influxdb_data:
  elasticsearch_data:
```

### 10.2 Environment Variables

```env
# .env
DB_HOST=localhost
DB_PORT=5432
DB_USER=parking
DB_PASSWORD=password
DB_NAME=parking

INFLUX_URL=http://localhost:8086
INFLUX_TOKEN=admin:password
INFLUX_ORG=parking
INFLUX_BUCKET=parking_stats

MQTT_BROKER=localhost:1883
REDIS_HOST=localhost:6379

KAFKA_BROKERS=localhost:9092
ELASTICSEARCH_URL=http://localhost:9200

# Printer Settings
PRINTER_RECEIPT_HOST=localhost
PRINTER_RECEIPT_PORT=9100
PRINTER_INVOICE_HOST=localhost
PRINTER_INVOICE_PORT=9101
PRINTER_TICKET_HOST=localhost
PRINTER_TICKET_PORT=9102

# QR Code Settings
QR_CODE_SIZE=256
QR_CODE_EXPIRY=600

# Payment Settings
PAYMENT_GATEWAY_URL=https://payment.example.com
PAYMENT_API_KEY=your_api_key
PAYMENT_SECRET_KEY=your_secret_key

# JWT Settings
JWT_SECRET=your_jwt_secret
JWT_EXPIRY=86400

# Rate Limit
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60

# LLM Settings
LLM_API_URL=https://api.llm.example.com
LLM_API_KEY=your_llm_api_key

# Blockchain Settings
BLOCKCHAIN_URL=https://blockchain.example.com
BLOCKCHAIN_API_KEY=your_blockchain_api_key

# Notification Settings
LINE_NOTIFY_TOKEN=your_line_token
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_email_password
```

### 10.3 Installation Steps

```bash
# 1. Clone repository
git clone https://github.com/yourproject/smart-parking.git
cd smart-parking

# 2. Copy environment file
cp .env.example .env

# 3. Build and start services
docker-compose up -d

# 4. Run database migrations
docker-compose exec backend go run cmd/migrate/main.go

# 5. Create admin user
docker-compose exec backend go run cmd/seed/main.go

# 6. Access dashboard
# Frontend: http://localhost:3000
# API: http://localhost:8080
# Kibana: http://localhost:5601
```

---

## 11. Business Model

### 11.1 Packaging & Pricing

| Package | Price | Spaces | Features |
|---------|-------|--------|----------|
| **Starter** | ฿49,900 | 50 | 1 Building, LPR, QR Payment, Basic Report |
| **Professional** | ฿149,900 | 200 | 3 Buildings, EV Charging, RFID, Statistics |
| **Enterprise** | ฿499,900 | 500+ | Unlimited Buildings, Lift System, Custom Reports |
| **Ultimate** | ฿999,900 | 1000+ | Full Features, On-premise, 24/7 Support |

### 11.2 Hardware Packages

| Package | Price | Components |
|---------|-------|------------|
| **Basic Kit** | ฿250,000 | LPR Camera, Boom Barrier, LED, Ticket Printer |
| **Standard Kit** | ฿500,000 | Basic + RFID, EV Charger, Payment Kiosk |
| **Premium Kit** | ฿1,000,000 | Standard + Lift System, Full CCTV, Emergency System |

### 11.3 Revenue Model

```
┌─────────────────────────────────────────────────────────────────┐
│                    Revenue Model                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Hardware Sales (One-time)                                     │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │  Basic Kit: ฿250,000                                     │ │
│  │  Standard Kit: ฿500,000                                  │ │
│  │  Premium Kit: ฿1,000,000                                 │ │
│  └───────────────────────────────────────────────────────────┘ │
│                              +                                  │
│  Subscription (Annual)                                         │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │  Starter: ฿15,000/ปี                                     │ │
│  │  Professional: ฿45,000/ปี                                │ │
│  │  Enterprise: ฿150,000/ปี                                 │ │
│  │  Ultimate: ฿300,000/ปี                                   │ │
│  └───────────────────────────────────────────────────────────┘ │
│                              +                                  │
│  Transaction Fee                                               │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │  QR Payment: 1-2% per transaction                        │ │
│  │  EV Charging: 5-10% per session                          │ │
│  └───────────────────────────────────────────────────────────┘ │
│                              +                                  │
│  Professional Services                                         │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │  Installation: ฿50,000-200,000                           │ │
│  │  Training: ฿20,000-50,000                                │ │
│  │  Custom Integration: ฿100,000-500,000                   │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 12. Prompt สำหรับการขยายระบบ

### Prompt 1: เพิ่มระบบจองที่จอดล่วงหน้า (Pre-booking)

```
ต้องการเพิ่มระบบจองที่จอดรถล่วงหน้าผ่าน Mobile App และ Website

รายละเอียด:
1. ผู้ใช้สามารถเลือกโซน/ชั้น/วันที่/เวลาที่ต้องการจอง
2. ระบบแสดงที่จอดว่างและราคา
3. ชำระเงินล่วงหน้า (Deposit) ผ่าน QR Code
4. ระบบกันที่จอดไว้ 15 นาทีก่อนถึงเวลาจอง
5. แจ้งเตือนผ่าน Line/Email เมื่อถึงเวลาจอง
6. ยกเลิกการจองและคืนเงิน (Refund) ตามนโยบาย

สิ่งที่ต้องเพิ่ม:
- Booking Entity และ CRUD
- Pre-payment System
- Space Reservation Logic
- Notification Service
- Booking Dashboard

ไฟล์ที่เกี่ยวข้อง:
- internal/modules/parking/domain/entity/booking.go
- internal/modules/parking/application/booking/
- internal/modules/parking/interfaces/http/booking_handler.go
- frontend/src/pages/Booking/
```

### Prompt 2: เพิ่มระบบวิเคราะห์ Predictive Analytics

```
ต้องการเพิ่มระบบวิเคราะห์และพยากรณ์อัตราการจอดรถ

รายละเอียด:
1. วิเคราะห์ข้อมูลย้อนหลัง 1 ปี
2. พยากรณ์อัตราการจอดในวัน/เดือนถัดไป
3. แนะนำการปรับราคาตามความต้องการ (Dynamic Pricing)
4. แจ้งเตือนเมื่อคาดการณ์ว่าจะเต็ม
5. แสดง Heatmap ของการจอดแต่ละโซน

สิ่งที่ต้องเพิ่ม:
- ML Model (Python/TensorFlow)
- Predictive Service
- Dynamic Pricing Engine
- Heatmap Visualization

ไฟล์ที่เกี่ยวข้อง:
- internal/modules/parking/infrastructure/analytics/
- python/models/parking_forecast.py
- frontend/src/components/Heatmap/
```

### Prompt 3: เพิ่มระบบ Integration กับ ERP

```
ต้องการเชื่อมต่อระบบ Parking กับ ERP ที่มีอยู่

รายละเอียด:
1. ส่งข้อมูลรายได้ไปยังระบบบัญชีอัตโนมัติ
2. ซิงค์ข้อมูลพนักงานและบัตร RFID
3. สร้าง Invoice และส่งไปยังระบบ ERP
4. เชื่อมต่อกับระบบ CRM
5. API Gateway สำหรับการ Integration

สิ่งที่ต้องเพิ่ม:
- ERP Adapter
- Invoice Export Service
- Data Synchronization
- Webhook System

ไฟล์ที่เกี่ยวข้อง:
- internal/modules/parking/infrastructure/integration/
- internal/modules/parking/application/sync/
- internal/modules/parking/interfaces/webhook/
```

### Prompt 4: เพิ่มระบบกล้อง AI และความปลอดภัย

```
ต้องการเพิ่มระบบกล้อง AI สำหรับความปลอดภัย

รายละเอียด:
1. ตรวจจับใบหน้าและจดจำ (Face Recognition)
2. ตรวจจับพฤติกรรมผิดปกติ (Anomaly Detection)
3. นับจำนวนรถอัตโนมัติ (Vehicle Counting)
4. ตรวจจับป้ายทะเบียนที่ถูกขโมย
5. ระบบแจ้งเตือนเหตุการณ์ผิดปกติ
6. บันทึกวิดีโอและค้นหาย้อนหลัง

สิ่งที่ต้องเพิ่ม:
- AI Camera Service
- Face Recognition
- Anomaly Detection
- Video Recording & Search

ไฟล์ที่เกี่ยวข้อง:
- internal/modules/parking/infrastructure/ai/
- internal/modules/parking/infrastructure/video/
- frontend/src/components/CameraView/
```

### Prompt 5: เพิ่มระบบ Loyalty & Membership

```
ต้องการเพิ่มระบบสมาชิกและสิทธิประโยชน์

รายละเอียด:
1. สมัครสมาชิก (Member Registration)
2. สะสมแต้ม (Points Accumulation)
3. แลกส่วนลด (Discount Redemption)
4. ระดับสมาชิก (Tier Levels: Silver, Gold, Platinum)
5. สิทธิพิเศษตามระดับ (Free parking, EV charging)
6. แจ้งเตือนโปรโมชั่น

สิ่งที่ต้องเพิ่ม:
- Member Entity และ CRUD
- Points System
- Tier Management
- Promotion Engine
- Member Dashboard

ไฟล์ที่เกี่ยวข้อง:
- internal/modules/parking/domain/entity/member.go
- internal/modules/parking/application/membership/
- internal/modules/parking/application/promotion/
- frontend/src/pages/Membership/
```

### Prompt 6: เพิ่มระบบ Report & Dashboard ขั้นสูง

```
ต้องการเพิ่มระบบรายงานและ Dashboard แบบ Advance

รายละเอียด:
1. Dashboard แบบ Real-time (ทุก 5 วินาที)
2. รายงานแบบ Interactive (กรองข้อมูลเอง)
3. Export Report เป็น Excel, PDF, CSV
4. Scheduled Report (ส่งอีเมลอัตโนมัติ)
5. เปรียบเทียบข้อมูลข้ามปี/ไตรมาส
6. KPI Dashboard (Occupancy, Revenue, EV Usage)

สิ่งที่ต้องเพิ่ม:
- Advanced Report Engine
- Interactive Dashboard
- Export Service
- Scheduled Report
- KPI Calculator

ไฟล์ที่เกี่ยวข้อง:
- internal/modules/parking/infrastructure/report/
- internal/modules/parking/application/dashboard/
- frontend/src/pages/Reports/
- frontend/src/components/Charts/
```

---

## 13. ภาคผนวก

### 13.1 Hardware Specifications

| Component | Model | จำนวน | ราคา (บาท) |
|-----------|-------|-------|------------|
| **LPR Camera** | Hikvision DS-2CD7A26G0 | 4 | 35,000 |
| **CCTV Camera** | Hikvision DS-2CD2347G2-LU | 20 | 15,000 |
| **RFID Reader** | Zebra FX9600 | 2 | 25,000 |
| **RFID Tag** | UHF Passive | 1,000 | 5,000 |
| **Boom Barrier** | Nice PLATINUM | 2 | 45,000 |
| **LED Indicator** | Green/Red LED | 200 | 2,000 |
| **Ultrasonic Sensor** | HC-SR04 | 200 | 100 |
| **EV Charger** | EVBox Troniq | 10 | 80,000 |
| **LCD Display** | 32" Digital Signage | 4 | 15,000 |
| **Payment Kiosk** | Self-service Kiosk | 2 | 50,000 |
| **Ticket Printer** | Zebra ZD420 | 2 | 20,000 |
| **QR Scanner** | Zebra DS4608 | 4 | 8,000 |
| **Receipt Printer** | Epson TM-T88VII | 2 | 12,000 |
| **Invoice Printer** | Epson LQ-590II | 1 | 8,000 |
| **Industrial PC** | Advantech IPC | 4 | 25,000 |
| **Network Switch** | PoE Switch 24-port | 6 | 12,000 |
| **UPS** | APC SMT1500 | 4 | 15,000 |
| **Lift Control** | PLC Controller | 1 | 100,000 |
| **Emergency Light** | LED Emergency | 20 | 2,500 |
| **รวม** | | | **~1,000,000** |

### 13.2 API Reference

| Endpoint | Method | Description | Auth |
|----------|--------|-------------|------|
| `/api/v1/parking/entry` | POST | Vehicle entry | No |
| `/api/v1/parking/exit` | POST | Vehicle exit | No |
| `/api/v1/parking/space/available` | GET | Get available spaces | No |
| `/api/v1/parking/qr/:ticket_no` | GET | Get ticket info | No |
| `/api/v1/parking/payment/qr` | POST | Generate QR payment | No |
| `/api/v1/parking/payment/confirm` | POST | Confirm payment | No |
| `/api/v1/parking/payment/discount` | POST | Apply discount | No |
| `/api/v1/parking/payment/rates` | GET | Get parking rates | No |
| `/api/v1/parking/print/ticket/:id` | POST | Print ticket | No |
| `/api/v1/parking/print/receipt/:id` | POST | Print receipt | No |
| `/api/v1/parking/print/invoice/:id` | POST | Print invoice | No |
| `/api/v1/parking/member/register` | POST | Member registration | No |
| `/api/v1/parking/member/login` | POST | Member login | No |
| `/api/v1/parking/member/profile` | GET | Get member profile | Yes |
| `/api/v1/parking/admin/buildings` | GET | Get buildings | Yes |
| `/api/v1/parking/admin/buildings` | POST | Create building | Yes |
| `/api/v1/parking/admin/spaces` | GET | Get spaces | Yes |
| `/api/v1/parking/admin/rates` | GET | Get rates | Yes |
| `/api/v1/parking/admin/discounts` | GET | Get discounts | Yes |
| `/api/v1/parking/admin/statistics/daily` | GET | Daily statistics | Yes |
| `/api/v1/parking/admin/equipment` | GET | Get equipment | Yes |
| `/api/v1/parking/admin/mode` | GET | Get system mode | Yes |
| `/api/v1/parking/admin/emergency/events` | GET | Get emergency events | Yes |

### 13.3 Glossary

| Term | Description |
|------|-------------|
| **LPR** | License Plate Recognition - ระบบจดจำป้ายทะเบียน |
| **RFID** | Radio Frequency Identification - ระบบระบุคลื่นความถี่วิทยุ |
| **EV** | Electric Vehicle - รถยนต์ไฟฟ้า |
| **QR Code** | Quick Response Code - รหัสตอบสนองเร็ว |
| **IoT** | Internet of Things - อินเทอร์เน็ตของสรรพสิ่ง |
| **M2M** | Machine to Machine - การสื่อสารระหว่างเครื่องจักร |
| **PLC** | Programmable Logic Controller - ตัวควบคุมตรรกะที่ตั้งโปรแกรมได้ |
| **ESC/POS** | Epson Standard Code for Point of Service - รหัสมาตรฐานสำหรับจุดบริการ |
| **DDD** | Domain-Driven Design - การออกแบบที่ขับเคลื่อนด้วยโดเมน |
| **Kafka** | Apache Kafka - แพลตฟอร์มสตรีมมิ่งแบบกระจาย |
| **Elasticsearch** | Elasticsearch - เครื่องมือค้นหาและวิเคราะห์แบบกระจาย |
| **InfluxDB** | InfluxDB - ฐานข้อมูลอนุกรมเวลาสำหรับการวิเคราะห์ |
| **WebSocket** | WebSocket - โปรโตคอลการสื่อสารแบบสองทิศทาง |
| **REST** | Representational State Transfer - สถาปัตยกรรม API |
| **JWT** | JSON Web Token - โทเค็นสำหรับการรับรองความถูกต้อง |

---

## สรุป

ระบบ Smart Parking Management System เป็นโซลูชันที่ครบวงจรสำหรับการบริหารจัดการอาคารจอดรถอัจฉริยะ โดยมี:

1. **Multi-Building & Multi-Zone** - รองรับหลายอาคาร หลายชั้น หลายโซน
2. **Smart Entry/Exit** - LPR, RFID, QR Code, Boom Barrier
3. **EV Charging** - รองรับการชาร์จรถไฟฟ้า
4. **Parking Guidance** - LED แสดงสถานะ, แผนที่นำทาง
5. **Payment System** - QR Code, Cash, Card, RFID
6. **Printing System** - บัตรจอดรถ, ใบเสร็จ, ใบแจ้งหนี้
7. **Discount System** - QR Code ส่วนลด, โปรโมชั่น
8. **Statistics & Reports** - รายชั่วโมง/วัน/สัปดาห์/เดือน/ปี/ไตรมาส
9. **Emergency System** - ไฟฉุกเฉิน, แจ้งเตือน
10. **Maintenance System** - ตารางบำรุงรักษา, แจ้งเตือน
11. **Anti-Crash System** - ระบบป้องกันการกระแทก
12. **Auto/Manual Mode** - สลับโหมดการทำงาน

เหมาะสำหรับการขายให้กับ:
- ห้างสรรพสินค้าและศูนย์การค้า
- อาคารสำนักงานและคอนโดมิเนียม
- โรงพยาบาลและมหาวิทยาลัย
- สนามบินและสถานีขนส่ง
- โรงแรมและรีสอร์ท