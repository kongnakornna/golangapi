
---

# 🎯 PART 6 — APPLICATION LAYER DEEP DIVE

> **ขนาด**: ใหญ่ — แยก 9 ตอนย่อย
> **Part 6A**: Application Layer Architecture
> **Part 6B**: Use Case Patterns (Command / Query)
> **Part 6C**: All Use Cases by Module
> **Part 6D**: DTOs & Mappers
> **Part 6E**: Query Handlers (CQRS read side)
> **Part 6F**: Transaction & Idempotency
> **Part 6G**: Validation & Orchestration
> **Part 6H**: Event Publishing & Side Effects
> **Part 6I**: Application Errors

---

## 🅰️ PART 6A — APPLICATION LAYER ARCHITECTURE

### A.1 Position & Responsibilities

```
┌─────────────────────────────────────────────────────────────┐
│                    Interface Layer                          │
│  HTTP · WebSocket · MQTT · CLI · Scheduler                  │
├─────────────────────────────────────────────────────────────┤
│  ★★★ APPLICATION LAYER ★★★  ← PART 6 นี้                   │
│  Use Cases · DTOs · Orchestration · Idempotency             │
│  ────────────────────────────────────────────────────────   │
│  ✓ Orchestrate domain objects                              │
│  ✓ Transaction boundary                                     │
│  ✓ Idempotency handling                                     │
│  ✓ Event publishing (via port)                              │
│  ✗ ไม่มี business logic (อยู่ที่ domain)                    │
│  ✗ ไม่มี SQL / HTTP / MQTT โดยตรง                          │
├─────────────────────────────────────────────────────────────┤
│                      Domain Layer                           │
│  Entities · VOs · Repos (interface) · Services              │
└─────────────────────────────────────────────────────────────┘
```

### A.2 Use Case Rules

| # | กฎ | เหตุผล |
|:-:|:---|:---|
| 1 | 1 use case / 1 file | `{{verb}}_{{aggregate}}.go` |
| 2 | `Execute(ctx, input) (output, error)` | Interface ชัดเจน |
| 3 | Input/Output เป็น struct (DTO) | ไม่ใช้ primitive |
| 4 | ✅ Transaction boundary ที่นี่ | 1 use case = 1 transaction |
| 5 | ✅ Idempotency key ที่นี่ | ป้องกัน duplicate |
| 6 | ❌ ไม่มี SQL / HTTP / MQTT | แยก infrastructure ออก |
| 7 | ❌ ไม่มี business logic | ย้ายไป domain |
| 8 | Side effects ที่ไม่ critical → log warn | ไม่ fail transaction |
| 9 | Constructor injection | Testable |
| 10 | Pull events → publish หลัง commit | Consistency |

### A.3 Folder Structure

```
internal/modules/{{module}}/application/
├── {{verb}}_{{aggregate}}.go         # 1 use case / file
├── {{verb}}_{{aggregate}}_test.go
├── query/
│   ├── {{name}}_query.go             # Query handler
│   └── {{name}}_query_test.go
├── dto.go                            # Request/Response DTOs
├── mappers.go                        # Entity ↔ DTO
├── ports.go                          # Outbound ports (EventProducer, etc.)
└── errors.go                         # Application errors
```

### A.4 Use Case Template

```go
// internal/modules/{{module}}/application/{{verb}}_{{aggregate}}.go
package application

import (
    "context"
    "errors"
    "fmt"

    "github.com/google/uuid"
    "icmongolang/internal/modules/{{module}}/domain/entity"
    domainerrors "icmongolang/internal/modules/{{module}}/domain/errors"
    "icmongolang/internal/modules/{{module}}/domain/repository"
)

// ============================================================
// UseCase
// ============================================================
type {{Verb}}{{Aggregate}}UseCase struct {
    repo      repository.{{Aggregate}}Repository
    txManager TransactionManager
    producer  EventProducer
    logger    Logger
    metrics   MetricsRecorder
}

func New{{Verb}}{{Aggregate}}UseCase(
    repo repository.{{Aggregate}}Repository,
    txManager TransactionManager,
    producer EventProducer,
    logger Logger,
    metrics MetricsRecorder,
) *{{Verb}}{{Aggregate}}UseCase {
    return &{{Verb}}{{Aggregate}}UseCase{
        repo:      repo,
        txManager: txManager,
        producer:  producer,
        logger:    logger,
        metrics:   metrics,
    }
}

// ============================================================
// Input / Output DTOs
// ============================================================
type {{Verb}}{{Aggregate}}Input struct {
    TenantID     uuid.UUID
    ActorID      uuid.UUID
    IdempotencyKey string
    // ... other fields
}

type {{Verb}}{{Aggregate}}Output struct {
    ID        uuid.UUID
    // ... other fields
}

// ============================================================
// Execute
// ============================================================
func (uc *{{Verb}}{{Aggregate}}UseCase) Execute(
    ctx context.Context,
    in {{Verb}}{{Aggregate}}Input,
) (*{{Verb}}{{Aggregate}}Output, error) {
    // 1. Validate input
    if err := uc.validate(in); err != nil {
        return nil, fmt.Errorf("validate: %w", err)
    }

    // 2. Idempotency check
    if in.IdempotencyKey != "" {
        cached, err := uc.checkIdempotency(ctx, in)
        if err != nil {
            return nil, err
        }
        if cached != nil {
            return cached, nil
        }
    }

    // 3. Execute in transaction
    var out *{{Verb}}{{Aggregate}}Output
    err := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
        var err error
        out, err = uc.execute(txCtx, in)
        return err
    })
    if err != nil {
        uc.metrics.IncError("{{verb}}_{{aggregate}}")
        return nil, err
    }

    // 4. Post-commit: publish events
    uc.publishEvents(ctx, out)

    // 5. Save idempotency result
    if in.IdempotencyKey != "" {
        _ = uc.saveIdempotency(ctx, in, out)
    }

    uc.metrics.IncSuccess("{{verb}}_{{aggregate}}")
    return out, nil
}

// ============================================================
// Internal
// ============================================================
func (uc *{{Verb}}{{Aggregate}}UseCase) validate(in {{Verb}}{{Aggregate}}Input) error {
    if in.TenantID == uuid.Nil {
        return domainerrors.ErrInvalidTenant
    }
    if in.ActorID == uuid.Nil {
        return domainerrors.ErrUnauthorized
    }
    return nil
}

func (uc *{{Verb}}{{Aggregate}}UseCase) execute(ctx context.Context, in {{Verb}}{{Aggregate}}Input) (*{{Verb}}{{Aggregate}}Output, error) {
    // Business logic here (delegate to domain)
    // ...
    return &{{Verb}}{{Aggregate}}Output{}, nil
}

func (uc *{{Verb}}{{Aggregate}}UseCase) publishEvents(ctx context.Context, out *{{Verb}}{{Aggregate}}Output) {
    // ...
}
```

---

## 🅱️ PART 6B — USE CASE PATTERNS

### B.1 Command Use Case (Write)

```go
// internal/modules/device/application/register_device.go
package application

type RegisterDeviceUseCase struct {
    deviceRepo  repository.DeviceRepository
    siteRepo    repository.SiteRepository
    quotaSvc    *service.QuotaService
    subRepo     repository.SubscriptionRepository
    txManager   TransactionManager
    producer    EventProducer
    logger      Logger
}

type RegisterDeviceInput struct {
    TenantID       uuid.UUID
    CustomerID     uuid.UUID
    SiteID         uuid.UUID
    ZoneID         uuid.UUID
    SerialNo       string
    Name           string
    Type           string
    Protocol       string
    ModelID        uuid.UUID
    Tags           []string
    IdempotencyKey string
    ActorID        uuid.UUID
}

type RegisterDeviceOutput struct {
    ID       uuid.UUID
    SerialNo string
}

func (uc *RegisterDeviceUseCase) Execute(ctx context.Context, in RegisterDeviceInput) (*RegisterDeviceOutput, error) {
    // 1. Validate
    if in.TenantID == uuid.Nil {
        return nil, domainerrors.ErrInvalidTenant
    }

    // 2. Parse VOs
    serial, err := vo.NewDeviceSerial(in.SerialNo)
    if err != nil {
        return nil, fmt.Errorf("%w: %v", domainerrors.ErrInvalidSerial, err)
    }
    dtype := vo.DeviceType(in.Type)
    if !dtype.IsValid() {
        return nil, domainerrors.ErrInvalidDeviceType
    }
    proto := vo.Protocol(in.Protocol)
    if !proto.IsValid() {
        return nil, domainerrors.ErrInvalidProtocol
    }

    // 3. Quota check
    sub, err := uc.subRepo.FindActiveByCustomer(ctx, in.TenantID, in.CustomerID)
    if err != nil {
        return nil, fmt.Errorf("find subscription: %w", err)
    }
    pkg, err := uc.subRepo.FindPackage(ctx, sub.PackageID)
    if err != nil {
        return nil, fmt.Errorf("find package: %w", err)
    }
    count, err := uc.deviceRepo.CountByCustomer(ctx, in.TenantID, in.CustomerID)
    if err != nil {
        return nil, err
    }
    quotaCheck := uc.quotaSvc.CheckDeviceQuota(*pkg, count)
    if !quotaCheck.Allowed {
        return nil, fmt.Errorf("%w: %s", domainerrors.ErrQuotaExceeded, quotaCheck.Message)
    }

    // 4. Duplicate check
    exists, err := uc.deviceRepo.ExistsBySerial(ctx, in.TenantID, serial)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, fmt.Errorf("%w: serial=%s", domainerrors.ErrSerialAlreadyExists, serial)
    }

    // 5. Create aggregate
    device, err := entity.NewDevice(
        in.TenantID, in.CustomerID, in.SiteID, in.ZoneID,
        serial, dtype, proto, in.Name, in.ActorID,
    )
    if err != nil {
        return nil, err
    }
    if len(in.Tags) > 0 {
        device.Tags = in.Tags
    }
    if in.ModelID != uuid.Nil {
        device.ModelID = in.ModelID
    }

    // 6. Persist + publish in transaction
    var out *RegisterDeviceOutput
    err = uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
        if err := uc.deviceRepo.Save(txCtx, device); err != nil {
            return fmt.Errorf("save device: %w", err)
        }
        // Pull domain events
        events := device.PullEvents()
        if err := uc.producer.PublishBatch(txCtx, events); err != nil {
            return fmt.Errorf("publish events: %w", err)
        }
        out = &RegisterDeviceOutput{
            ID:       device.ID,
            SerialNo: device.SerialNo.String(),
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    return out, nil
}
```

### B.2 Query Use Case (Read)

```go
// internal/modules/device/application/query/get_device.go
package query

type GetDeviceQuery struct {
    readModel repository.DeviceReadModelRepository
    cache     Cache
}

type GetDeviceInput struct {
    TenantID uuid.UUID
    DeviceID uuid.UUID
}

type DeviceDTO struct {
    ID           string      `json:"id"`
    TenantID     string      `json:"tenant_id"`
    SerialNo     string      `json:"serial_no"`
    Name         string      `json:"name"`
    Type         string      `json:"type"`
    Protocol     string      `json:"protocol"`
    Status       string      `json:"status"`
    SiteID       string      `json:"site_id"`
    SiteName     string      `json:"site_name"`
    Firmware     string      `json:"firmware,omitempty"`
    LastSeenAt   *time.Time  `json:"last_seen_at,omitempty"`
    CreatedAt    time.Time   `json:"created_at"`
    UpdatedAt    time.Time   `json:"updated_at"`
}

func NewGetDeviceQuery(rm repository.DeviceReadModelRepository, cache Cache) *GetDeviceQuery {
    return &GetDeviceQuery{readModel: rm, cache: cache}
}

func (q *GetDeviceQuery) Execute(ctx context.Context, in GetDeviceInput) (*DeviceDTO, error) {
    // 1. Try cache
    cacheKey := fmt.Sprintf("device:%s:%s", in.TenantID, in.DeviceID)
    if cached, err := q.cache.Get(ctx, cacheKey); err == nil && cached != nil {
        var dto DeviceDTO
        if err := json.Unmarshal(cached, &dto); err == nil {
            return &dto, nil
        }
    }

    // 2. Query read model
    rm, err := q.readModel.Get(ctx, in.TenantID, in.DeviceID)
    if err != nil {
        return nil, err
    }

    // 3. Map to DTO
    dto := &DeviceDTO{
        ID:         rm.ID.String(),
        TenantID:   in.TenantID.String(),
        SerialNo:   rm.SerialNo,
        Name:       rm.Name,
        Type:       rm.Type,
        Status:     rm.Status,
        SiteID:     rm.SiteID.String(),
        SiteName:   rm.SiteName,
        LastSeenAt: rm.LastSeenAt,
        CreatedAt:  rm.CreatedAt,
    }

    // 4. Cache (TTL 5 นาที)
    if data, err := json.Marshal(dto); err == nil {
        _ = q.cache.Set(ctx, cacheKey, data, 5*time.Minute)
    }

    return dto, nil
}
```

### B.3 List Query (Paginated)

```go
// internal/modules/device/application/query/list_devices.go
package query

type ListDevicesQuery struct {
    readModel repository.DeviceReadModelRepository
}

type ListDevicesInput struct {
    TenantID   uuid.UUID
    SiteID     *uuid.UUID
    Status     *string
    Type       *string
    Search     string
    Page       int
    PageSize   int
    SortBy     string
    SortOrder  string
}

type ListDevicesOutput struct {
    Items      []*DeviceDTO `json:"items"`
    Pagination PaginationDTO `json:"pagination"`
}

type PaginationDTO struct {
    Page       int   `json:"page"`
    PageSize   int   `json:"page_size"`
    Total      int64 `json:"total"`
    TotalPages int   `json:"total_pages"`
}

func (q *ListDevicesQuery) Execute(ctx context.Context, in ListDevicesInput) (*ListDevicesOutput, error) {
    // Normalize
    if in.Page < 1 {
        in.Page = 1
    }
    if in.PageSize < 1 || in.PageSize > 100 {
        in.PageSize = 20
    }
    if in.SortBy == "" {
        in.SortBy = "created_at"
    }
    if in.SortOrder != "asc" && in.SortOrder != "desc" {
        in.SortOrder = "desc"
    }

    // Query
    var statusPtr *vo.DeviceStatus
    if in.Status != nil {
        s := vo.DeviceStatus(*in.Status)
        statusPtr = &s
    }
    var typePtr *vo.DeviceType
    if in.Type != nil {
        t := vo.DeviceType(*in.Type)
        typePtr = &t
    }

    result, err := q.readModel.List(ctx, repository.DeviceListQuery{
        TenantID:  in.TenantID,
        SiteID:    in.SiteID,
        Status:    statusPtr,
        Type:      typePtr,
        Search:    in.Search,
        Page:      in.Page,
        PageSize:  in.PageSize,
        SortBy:    in.SortBy,
        SortOrder: in.SortOrder,
    })
    if err != nil {
        return nil, err
    }

    // Map
    items := make([]*DeviceDTO, 0, len(result.Items))
    for _, rm := range result.Items {
        items = append(items, &DeviceDTO{
            ID:         rm.ID.String(),
            SerialNo:   rm.SerialNo,
            Name:       rm.Name,
            Type:       rm.Type,
            Status:     rm.Status,
            SiteID:     rm.SiteID.String(),
            SiteName:   rm.SiteName,
            LastSeenAt: rm.LastSeenAt,
            CreatedAt:  rm.CreatedAt,
        })
    }

    totalPages := int(result.Pagination.Total) / in.PageSize
    if int(result.Pagination.Total)%in.PageSize != 0 {
        totalPages++
    }

    return &ListDevicesOutput{
        Items: items,
        Pagination: PaginationDTO{
            Page:       in.Page,
            PageSize:   in.PageSize,
            Total:      result.Pagination.Total,
            TotalPages: totalPages,
        },
    }, nil
}
```

---

## 🅲 PART 6C — ALL USE CASES BY MODULE

### C.1 Use Case Inventory

| Module | Commands | Queries | Total |
|:---|:-:|:-:|:-:|
| **Auth** | 5 | 1 | 6 |
| **Customer** | 8 | 4 | 12 |
| **Package** | 6 | 4 | 10 |
| **Device** | 15 | 8 | 23 |
| **ERP** | 20 | 15 | 35 |
| **CRM** | 12 | 8 | 20 |
| **Logistics** | 15 | 10 | 25 |
| **Report** | 6 | 8 | 14 |
| **Alarm** | 5 | 4 | 9 |
| **Notifier** | 3 | 2 | 5 |
| **รวม** | **95** | **64** | **159** |

### C.2 Customer Module — Use Cases

```go
// ============================================================
// COMMANDS
// ============================================================
type CreateCustomerUseCase struct{...}      // POST /customers
type UpdateCustomerUseCase struct{...}      // PUT /customers/:id
type DeleteCustomerUseCase struct{...}      // DELETE /customers/:id
type QualifyLeadUseCase struct{...}         // POST /customers/:id/qualify
type ActivateCustomerUseCase struct{...}    // POST /customers/:id/activate
type SuspendCustomerUseCase struct{...}     // POST /customers/:id/suspend
type OnboardCustomerUseCase struct{...}     // POST /customers/:id/onboard
type AddContactUseCase struct{...}          // POST /customers/:id/contacts

// ============================================================
// QUERIES
// ============================================================
type GetCustomerQuery struct{...}           // GET /customers/:id
type ListCustomersQuery struct{...}         // GET /customers
type ListContactsQuery struct{...}          // GET /customers/:id/contacts
type CustomerStatsQuery struct{...}         // GET /customers/stats
```

```go
// internal/modules/customer/application/create_customer.go
package application

type CreateCustomerUseCase struct {
    repo     repository.CustomerRepository
    codeGen  CodeGenerator
    txMgr    TransactionManager
    producer EventProducer
}

type CreateCustomerInput struct {
    TenantID       uuid.UUID
    ActorID        uuid.UUID
    Type           string
    Name           string
    LegalName      string
    TaxID          string
    Email          string
    Phone          string
    Segment        string
    Address        AddressDTO
    CreditLimit    float64
    Currency       string
    IdempotencyKey string
}

type CreateCustomerOutput struct {
    ID   uuid.UUID `json:"id"`
    Code string    `json:"code"`
}

func (uc *CreateCustomerUseCase) Execute(ctx context.Context, in CreateCustomerInput) (*CreateCustomerOutput, error) {
    // 1. Validate
    if in.TenantID == uuid.Nil {
        return nil, domainerrors.ErrInvalidTenant
    }
    if err := validateAddress(in.Address); err != nil {
        return nil, err
    }

    // 2. Parse VOs
    ctype := vo.CustomerType(in.Type)
    if !ctype.IsValid() {
        return nil, domainerrors.ErrInvalidCustomer
    }
    seg := vo.Segment(in.Segment)
    if !seg.IsValid() {
        return nil, domainerrors.ErrInvalidCustomer
    }

    // 3. Generate code
    code, err := uc.codeGen.NextCustomerCode(ctx, in.TenantID)
    if err != nil {
        return nil, fmt.Errorf("generate code: %w", err)
    }

    // 4. Create aggregate
    customer, err := entity.NewCustomer(in.TenantID, code, ctype, in.Name, in.ActorID)
    if err != nil {
        return nil, err
    }
    customer.LegalName = in.LegalName
    customer.Email = in.Email
    customer.Phone = in.Phone
    customer.Segment = seg
    customer.Address = mapAddress(in.Address)

    if in.TaxID != "" {
        taxID, err := vo.NewTaxID(in.TaxID)
        if err != nil {
            return nil, fmt.Errorf("%w: %v", domainerrors.ErrInvalidTaxID, err)
        }
        if err := customer.SetTaxID(taxID); err != nil {
            return nil, err
        }
    }
    if in.CreditLimit > 0 {
        money, _ := vo.NewMoney(in.CreditLimit, in.Currency)
        if err := customer.UpdateCreditLimit(money); err != nil {
            return nil, err
        }
    }

    // 5. Persist + publish
    err = uc.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
        if err := uc.repo.Save(txCtx, customer); err != nil {
            return fmt.Errorf("save customer: %w", err)
        }
        return uc.producer.PublishBatch(txCtx, customer.PullEvents())
    })
    if err != nil {
        return nil, err
    }

    return &CreateCustomerOutput{ID: customer.ID, Code: customer.Code.String()}, nil
}
```

### C.3 Device Module — Use Cases (23)

```go
// ---- Commands ----
RegisterDeviceUseCase           // POST /devices
UpdateDeviceUseCase             // PUT /devices/:id
DeleteDeviceUseCase             // DELETE /devices/:id
ProvisionDeviceUseCase          // POST /devices/:id/provision
RotateDeviceTokenUseCase        // POST /devices/:id/rotate-token
SendCommandUseCase              // POST /devices/:id/commands
IngestTelemetryUseCase          // POST /devices/:id/telemetry (or MQTT)
UpdateFirmwareUseCase           // POST /devices/:id/firmware
UpdateShadowDesiredUseCase      // PATCH /devices/:id/shadow/desired
ClearShadowDesiredUseCase       // DELETE /devices/:id/shadow/desired
CreateAlertRuleUseCase          // POST /devices/alert-rules
UpdateAlertRuleUseCase          // PUT /devices/alert-rules/:id
DeleteAlertRuleUseCase          // DELETE /devices/alert-rules/:id
AcknowledgeAlertUseCase         // POST /devices/alert-events/:id/ack
BulkProvisionUseCase            // POST /devices/bulk-provision

// ---- Queries ----
GetDeviceQuery                  // GET /devices/:id
ListDevicesQuery                // GET /devices
GetDeviceShadowQuery            // GET /devices/:id/shadow
QueryTelemetryQuery             // GET /devices/:id/telemetry
LatestTelemetryQuery            // GET /devices/:id/telemetry/latest
ListCommandsQuery               // GET /devices/:id/commands
ListAlertRulesQuery             // GET /devices/alert-rules
ListAlertEventsQuery            // GET /devices/alert-events
```

### C.4 ERP Module — Use Cases (35)

```go
// ---- Products ----
CreateProductUseCase            // POST /erp/products
UpdateProductUseCase            // PUT /erp/products/:id
DeleteProductUseCase            // DELETE /erp/products/:id
ListProductsQuery               // GET /erp/products
GetProductQuery                 // GET /erp/products/:id

// ---- Warehouses ----
CreateWarehouseUseCase          // POST /erp/warehouses
ListWarehousesQuery             // GET /erp/warehouses

// ---- Inventory ----
AdjustStockUseCase              // POST /erp/inventory/adjust
TransferStockUseCase            // POST /erp/inventory/transfer
ReserveStockUseCase             // POST /erp/inventory/reserve
ReleaseReservationUseCase       // POST /erp/inventory/release
ListInventoryQuery              // GET /erp/inventory
GetStockQuery                   // GET /erp/inventory/:productID/:warehouseID
ListStockMovementsQuery         // GET /erp/stock-movements
LowStockQuery                   // GET /erp/inventory/low-stock

// ---- Orders ----
CreateOrderUseCase              // POST /erp/orders
UpdateOrderUseCase              // PUT /erp/orders/:id
AddOrderLineUseCase             // POST /erp/orders/:id/lines
RemoveOrderLineUseCase          // DELETE /erp/orders/:id/lines/:lineID
SubmitOrderUseCase              // POST /erp/orders/:id/submit
ConfirmOrderUseCase             // POST /erp/orders/:id/confirm
ShipOrderUseCase                // POST /erp/orders/:id/ship
ReceiveOrderUseCase             // POST /erp/orders/:id/receive
CloseOrderUseCase               // POST /erp/orders/:id/close
CancelOrderUseCase              // POST /erp/orders/:id/cancel
ListOrdersQuery                 // GET /erp/orders
GetOrderQuery                   // GET /erp/orders/:id

// ---- Invoices ----
CreateInvoiceUseCase            // POST /erp/invoices
IssueInvoiceUseCase             // POST /erp/invoices/:id/issue
VoidInvoiceUseCase              // POST /erp/invoices/:id/void
ListInvoicesQuery               // GET /erp/invoices
GetInvoiceQuery                 // GET /erp/invoices/:id
AgingReportQuery                // GET /erp/invoices/aging

// ---- Payments ----
CreatePaymentUseCase            // POST /erp/payments
AllocatePaymentUseCase          // POST /erp/payments/:id/allocate
ListPaymentsQuery               // GET /erp/payments
GetPaymentQuery                 // GET /erp/payments/:id

// ---- Reports ----
TrialBalanceQuery               // GET /erp/reports/trial-balance
ProfitLossQuery                 // GET /erp/reports/profit-loss
BalanceSheetQuery               // GET /erp/reports/balance-sheet
JournalEntriesQuery             // GET /erp/journal-entries
```

### C.5 Use Case Cross-Reference (ตัวอย่าง ShipOrder)

```go
// internal/modules/erp/application/ship_order.go
package application

type ShipOrderUseCase struct {
    orderRepo     repository.OrderRepository
    inventoryRepo repository.InventoryRepository
    movementRepo  repository.StockMovementRepository
    journalRepo   repository.JournalRepository
    txMgr         TransactionManager
    producer      EventProducer
    clock         Clock
}

type ShipOrderInput struct {
    TenantID    uuid.UUID
    OrderID     uuid.UUID
    WarehouseID uuid.UUID
    Lines       []ShipLineInput
    ActorID     uuid.UUID
}

type ShipLineInput struct {
    LineID   uuid.UUID
    Quantity decimal.Decimal
}

type ShipOrderOutput struct {
    OrderID     uuid.UUID
    Status      string
    Movements   []MovementDTO
}

func (uc *ShipOrderUseCase) Execute(ctx context.Context, in ShipOrderInput) (*ShipOrderOutput, error) {
    var out *ShipOrderOutput

    err := uc.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
        // 1. Load order
        order, err := uc.orderRepo.FindByID(txCtx, in.TenantID, in.OrderID)
        if err != nil {
            return err
        }

        // 2. Validate state
        if order.Status != vo.OrderStatusConfirmed {
            return domainerrors.ErrMustBeConfirmed
        }

        // 3. Check stock + reserve
        for _, line := range in.Lines {
            inv, err := uc.inventoryRepo.FindByKey(txCtx, in.TenantID, line.ProductID(order, line.LineID), in.WarehouseID)
            if err != nil {
                return err
            }
            available := inv.QtyOnHand.Sub(inv.QtyReserved)
            if available.LessThan(line.Quantity) {
                return fmt.Errorf("%w: product=%s available=%s need=%s",
                    domainerrors.ErrInsufficientStock,
                    inv.ProductID, available, line.Quantity)
            }
            // Reserve
            if err := inv.Reserve(line.Quantity); err != nil {
                return err
            }
            if err := uc.inventoryRepo.Update(txCtx, inv); err != nil {
                return err
            }
        }

        // 4. Ship (change order state)
        if err := order.Ship(in.ActorID); err != nil {
            return err
        }
        if err := uc.orderRepo.Update(txCtx, order); err != nil {
            return err
        }

        // 5. Create stock movements
        movements := make([]*entity.StockMovement, 0)
        for _, line := range in.Lines {
            inv, _ := uc.inventoryRepo.FindByKey(txCtx, in.TenantID, line.ProductID, in.WarehouseID)
            if err := inv.Issue(line.Quantity); err != nil {
                return err
            }
            if err := uc.inventoryRepo.Update(txCtx, inv); err != nil {
                return err
            }
            mv := entity.NewStockMovement(
                in.TenantID, inv.ProductID, in.WarehouseID,
                entity.MovementTypeIssue, -1, line.Quantity,
                vo.Money{}, inv.QtyOnHand, "ORDER", order.OrderNo.String(),
                in.ActorID,
            )
            if err := uc.movementRepo.Save(txCtx, mv); err != nil {
                return err
            }
            movements = append(movements, mv)
        }

        // 6. Post journal entry (COGS)
        if err := uc.postCOGS(txCtx, order, movements); err != nil {
            return err
        }

        // 7. Publish events
        if err := uc.producer.PublishBatch(txCtx, order.PullEvents()); err != nil {
            return err
        }

        out = &ShipOrderOutput{
            OrderID: order.ID,
            Status:  string(order.Status),
            Movements: toMovementDTOs(movements),
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    return out, nil
}
```

---

## 🅳 PART 6D — DTOs & MAPPERS

### D.1 DTO Convention

```go
// internal/modules/device/application/dto.go
package application

// ============================================================
// REQUEST DTOs (input from HTTP)
// ============================================================
type RegisterDeviceRequest struct {
    CustomerID string   `json:"customer_id" binding:"required,uuid4"`
    SiteID     string   `json:"site_id"     binding:"required,uuid4"`
    ZoneID     string   `json:"zone_id"     binding:"omitempty,uuid4"`
    SerialNo   string   `json:"serial_no"   binding:"required,min=3,max=50"`
    Name       string   `json:"name"        binding:"required,min=1,max=255"`
    Type       string   `json:"type"        binding:"required,oneof=SENSOR ACTUATOR GATEWAY CAMERA METER TRACKER"`
    Protocol   string   `json:"protocol"    binding:"required,oneof=MQTT MODBUS LORAWAN ZIGBEE HTTP COAP"`
    ModelID    string   `json:"model_id"    binding:"omitempty,uuid4"`
    Tags       []string `json:"tags"        binding:"omitempty,max=10,dive,min=1,max=30"`
}

type UpdateDeviceRequest struct {
    Name string `json:"name" binding:"omitempty,min=1,max=255"`
    Tags []string `json:"tags" binding:"omitempty,max=10,dive,min=1,max=30"`
}

type SendCommandRequest struct {
    Command      string                 `json:"command" binding:"required,min=1,max=50"`
    Payload      map[string]interface{} `json:"payload"`
    Priority     int                    `json:"priority" binding:"omitempty,min=1,max=5"`
    ExpiresInSec int                    `json:"expires_in_sec" binding:"omitempty,min=1,max=3600"`
}

type UpdateShadowDesiredRequest struct {
    Desired map[string]interface{} `json:"desired" binding:"required"`
}

// ============================================================
// RESPONSE DTOs (output to HTTP)
// ============================================================
type DeviceResponse struct {
    ID           string                 `json:"id"`
    TenantID     string                 `json:"tenant_id"`
    CustomerID   string                 `json:"customer_id"`
    SiteID       string                 `json:"site_id"`
    ZoneID       string                 `json:"zone_id,omitempty"`
    SerialNo     string                 `json:"serial_no"`
    Name         string                 `json:"name"`
    Type         string                 `json:"type"`
    Protocol     string                 `json:"protocol"`
    Status       string                 `json:"status"`
    Firmware     string                 `json:"firmware,omitempty"`
    LastSeenAt   *time.Time             `json:"last_seen_at,omitempty"`
    ProvisionedAt *time.Time            `json:"provisioned_at,omitempty"`
    Capabilities []CapabilityResponse   `json:"capabilities,omitempty"`
    Tags         []string               `json:"tags,omitempty"`
    Config       map[string]interface{} `json:"config,omitempty"`
    CreatedAt    time.Time              `json:"created_at"`
    UpdatedAt    time.Time              `json:"updated_at"`
}

type CapabilityResponse struct {
    Type    string `json:"type"`
    Metric  string `json:"metric,omitempty"`
    Command string `json:"command,omitempty"`
    Unit    string `json:"unit,omitempty"`
}

type CommandResponse struct {
    ID         string                 `json:"id"`
    DeviceID   string                 `json:"device_id"`
    Type       string                 `json:"type"`
    Payload    map[string]interface{} `json:"payload,omitempty"`
    Status     string                 `json:"status"`
    Priority   int                    `json:"priority"`
    IssuedBy   string                 `json:"issued_by"`
    IssuedAt   time.Time              `json:"issued_at"`
    SentAt     *time.Time             `json:"sent_at,omitempty"`
    AckedAt    *time.Time             `json:"acked_at,omitempty"`
    ExpiresAt  time.Time              `json:"expires_at"`
    Result     map[string]interface{} `json:"result,omitempty"`
}

type ShadowResponse struct {
    DeviceID  string                 `json:"device_id"`
    Desired   map[string]interface{} `json:"desired"`
    Reported  map[string]interface{} `json:"reported"`
    Delta     map[string]interface{} `json:"delta"`
    Version   int                    `json:"version"`
    UpdatedAt time.Time              `json:"updated_at"`
}

type TelemetryResponse struct {
    DeviceID string                       `json:"device_id"`
    From     time.Time                    `json:"from"`
    To       time.Time                    `json:"to"`
    Series   map[string][]TelemetryPointDTO `json:"series"`
}

type TelemetryPointDTO struct {
    Timestamp time.Time `json:"ts"`
    Value     float64   `json:"value"`
    Unit      string    `json:"unit,omitempty"`
    Quality   string    `json:"quality,omitempty"`
}
```

### D.2 Mapper Pattern

```go
// internal/modules/device/application/mappers.go
package application

import (
    "icmongolang/internal/modules/device/domain/entity"
)

func ToDeviceResponse(d *entity.Device) *DeviceResponse {
    if d == nil {
        return nil
    }
    resp := &DeviceResponse{
        ID:            d.ID.String(),
        TenantID:      d.TenantID.String(),
        CustomerID:    d.CustomerID.String(),
        SiteID:        d.SiteID.String(),
        SerialNo:      d.SerialNo.String(),
        Name:          d.Name,
        Type:          string(d.Type),
        Protocol:      string(d.Protocol),
        Status:        string(d.Status),
        Firmware:      d.FirmwareVersion,
        LastSeenAt:    d.LastSeenAt,
        ProvisionedAt: d.ProvisionedAt,
        Tags:          d.Tags,
        Config:        d.Config,
        CreatedAt:     d.CreatedAt,
        UpdatedAt:     d.UpdatedAt,
    }
    if d.ZoneID != uuid.Nil {
        resp.ZoneID = d.ZoneID.String()
    }
    for _, c := range d.Capabilities {
        resp.Capabilities = append(resp.Capabilities, CapabilityResponse{
            Type:    c.Type,
            Metric:  c.Metric,
            Command: c.Command,
            Unit:    c.Unit,
        })
    }
    return resp
}

func ToDeviceResponses(devices []*entity.Device) []*DeviceResponse {
    out := make([]*DeviceResponse, 0, len(devices))
    for _, d := range devices {
        out = append(out, ToDeviceResponse(d))
    }
    return out
}

func ToCommandResponse(c *entity.Command) *CommandResponse {
    if c == nil {
        return nil
    }
    return &CommandResponse{
        ID:        c.ID.String(),
        DeviceID:  c.DeviceID.String(),
        Type:      c.Type,
        Payload:   c.Payload,
        Status:    string(c.Status),
        Priority:  c.Priority,
        IssuedBy:  c.IssuedBy.String(),
        IssuedAt:  c.IssuedAt,
        SentAt:    c.SentAt,
        AckedAt:   c.AckedAt,
        ExpiresAt: c.ExpiresAt,
        Result:    c.Result,
    }
}

func ToShadowResponse(s *entity.DeviceShadow) *ShadowResponse {
    if s == nil {
        return nil
    }
    return &ShadowResponse{
        DeviceID:  s.DeviceID.String(),
        Desired:   s.Desired,
        Reported:  s.Reported,
        Delta:     s.ComputeDelta(),
        Version:   s.Version,
        UpdatedAt: s.UpdatedAt,
    }
}
```

### D.3 Request → Input Mapper

```go
// internal/modules/device/application/mappers.go (continued)
func (uc *RegisterDeviceUseCase) fromRequest(req RegisterDeviceRequest, tenantID, actorID uuid.UUID) (*RegisterDeviceInput, error) {
    customerID, err := uuid.Parse(req.CustomerID)
    if err != nil {
        return nil, domainerrors.ErrInvalidInput
    }
    siteID, err := uuid.Parse(req.SiteID)
    if err != nil {
        return nil, domainerrors.ErrInvalidInput
    }
    var zoneID uuid.UUID
    if req.ZoneID != "" {
        zoneID, err = uuid.Parse(req.ZoneID)
        if err != nil {
            return nil, domainerrors.ErrInvalidInput
        }
    }
    var modelID uuid.UUID
    if req.ModelID != "" {
        modelID, err = uuid.Parse(req.ModelID)
        if err != nil {
            return nil, domainerrors.ErrInvalidInput
        }
    }
    return &RegisterDeviceInput{
        TenantID:   tenantID,
        CustomerID: customerID,
        SiteID:     siteID,
        ZoneID:     zoneID,
        SerialNo:   req.SerialNo,
        Name:       req.Name,
        Type:       req.Type,
        Protocol:   req.Protocol,
        ModelID:    modelID,
        Tags:       req.Tags,
        ActorID:    actorID,
    }, nil
}
```

---

## 🅴 PART 6E — QUERY HANDLERS (CQRS READ SIDE)

### E.1 Read Model Repository (Interface)

```go
// internal/modules/device/domain/repository/read_model_repository.go
package repository

type DeviceReadModelRepository interface {
    List(ctx context.Context, q DeviceListQuery) (*DeviceListResult, error)
    Get(ctx context.Context, tenantID, id uuid.UUID) (*DeviceReadModel, error)
    StatsBySite(ctx context.Context, tenantID uuid.UUID) ([]SiteStats, error)
    StatsByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) (CustomerStats, error)
    Timeline(ctx context.Context, tenantID, deviceID uuid.UUID, from, to time.Time) ([]TimelineEvent, error)
}
```

### E.2 Read Model Implementation (Infrastructure)

```go
// internal/modules/device/infrastructure/persistence/postgres/read_model_repo_impl.go
package postgres

type DeviceReadModelRepoImpl struct {
    db *gorm.DB
}

func (r *DeviceReadModelRepoImpl) List(ctx context.Context, q repository.DeviceListQuery) (*repository.DeviceListResult, error) {
    // Build query
    tx := r.db.WithContext(ctx).
        Model(&DeviceModel{}).
        Where("tenant_id = ?", q.TenantID)

    if q.SiteID != nil {
        tx = tx.Where("site_id = ?", *q.SiteID)
    }
    if q.CustomerID != nil {
        tx = tx.Where("customer_id = ?", *q.CustomerID)
    }
    if q.Status != nil {
        tx = tx.Where("status = ?", string(*q.Status))
    }
    if q.Type != nil {
        tx = tx.Where("type = ?", string(*q.Type))
    }
    if q.Search != "" {
        pattern := "%" + strings.ToLower(q.Search) + "%"
        tx = tx.Where("LOWER(name) LIKE ? OR LOWER(serial_no) LIKE ?", pattern, pattern)
    }

    // Count
    var total int64
    if err := tx.Count(&total).Error; err != nil {
        return nil, err
    }

    // Sort
    sortBy := "created_at"
    switch q.SortBy {
    case "name", "status", "created_at", "updated_at":
        sortBy = q.SortBy
    }
    sortOrder := "DESC"
    if q.SortOrder == "asc" {
        sortOrder = "ASC"
    }
    tx = tx.Order(sortBy + " " + sortOrder)

    // Paginate
    offset := (q.Page - 1) * q.PageSize
    tx = tx.Offset(offset).Limit(q.PageSize)

    var models []DeviceModel
    if err := tx.Find(&models).Error; err != nil {
        return nil, err
    }

    // Map
    items := make([]*repository.DeviceReadModel, 0, len(models))
    for _, m := range models {
        items = append(items, &repository.DeviceReadModel{
            ID:         m.ID,
            SerialNo:   m.SerialNumber,
            Name:       m.Name,
            Type:       m.Type,
            Status:     m.Status,
            SiteID:     m.SiteID,
            LastSeenAt: m.LastSeenAt,
            CreatedAt:  m.CreatedAt,
        })
    }
    return &repository.DeviceListResult{
        Items: items,
        Pagination: repository.Pagination{
            Page:     q.Page,
            PageSize: q.PageSize,
            Total:    total,
        },
    }, nil
}
```

### E.3 Query Handler (Application)

```go
// internal/modules/device/application/query/list_devices.go
package query

type ListDevicesHandler struct {
    readModel repository.DeviceReadModelRepository
}

func NewListDevicesHandler(rm repository.DeviceReadModelRepository) *ListDevicesHandler {
    return &ListDevicesHandler{readModel: rm}
}

type ListDevicesInput struct {
    TenantID uuid.UUID
    SiteID   *uuid.UUID
    Status   *string
    Search   string
    Page     int
    PageSize int
}

type ListDevicesOutput struct {
    Items      []*DeviceDTO   `json:"items"`
    Pagination PaginationDTO  `json:"pagination"`
}

func (h *ListDevicesHandler) Handle(ctx context.Context, in ListDevicesInput) (*ListDevicesOutput, error) {
    result, err := h.readModel.List(ctx, repository.DeviceListQuery{
        TenantID: in.TenantID,
        SiteID:   in.SiteID,
        Status:   in.Status,
        Search:   in.Search,
        Page:     in.Page,
        PageSize: in.PageSize,
    })
    if err != nil {
        return nil, err
    }
    // Map to output
    out := &ListDevicesOutput{
        Items: make([]*DeviceDTO, 0, len(result.Items)),
        Pagination: PaginationDTO{
            Page:     in.Page,
            PageSize: in.PageSize,
            Total:    result.Pagination.Total,
        },
    }
    for _, rm := range result.Items {
        out.Items = append(out.Items, &DeviceDTO{
            ID:       rm.ID.String(),
            SerialNo: rm.SerialNo,
            Name:     rm.Name,
            Type:     rm.Type,
            Status:   rm.Status,
            SiteID:   rm.SiteID.String(),
        })
    }
    return out, nil
}
```

### E.4 CQRS Read Model Sync (Event Handler)

```go
// internal/modules/device/application/query/sync_handlers.go
package query

// SyncDeviceReadModel — consume domain events → update read model
type SyncDeviceReadModelHandler struct {
    readModel repository.DeviceReadModelRepository
}

func (h *SyncDeviceReadModelHandler) HandleDeviceCreated(ctx context.Context, event event.DeviceCreated) error {
    return h.readModel.Insert(ctx, &repository.DeviceReadModel{
        ID:        event.AggregateID(),
        SerialNo:  event.SerialNo,
        Type:      event.Type,
        SiteID:    event.SiteID,
        Status:    "OFFLINE",
        CreatedAt: event.OccurredAt(),
    })
}

func (h *SyncDeviceReadModelHandler) HandleDeviceOnline(ctx context.Context, event event.DeviceOnline) error {
    return h.readModel.UpdateStatus(ctx, event.AggregateID(), "ONLINE", &event.LastSeenAt)
}

func (h *SyncDeviceReadModelHandler) HandleDeviceOffline(ctx context.Context, event event.DeviceOffline) error {
    return h.readModel.UpdateStatus(ctx, event.AggregateID(), "OFFLINE", nil)
}
```

---

## 🅵 PART 6F — TRANSACTION & IDEMPOTENCY

### F.1 Transaction Manager Port

```go
// internal/shared/application/transaction.go
package application

import "context"

type TransactionManager interface {
    WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
    WithIsolation(ctx context.Context, level IsolationLevel, fn func(ctx context.Context) error) error
}

type IsolationLevel int

const (
    IsolationDefault IsolationLevel = iota
    IsolationReadCommitted
    IsolationRepeatableRead
    IsolationSerializable
)
```

### F.2 Transaction Manager Impl

```go
// internal/shared/infrastructure/postgres/tx_manager.go
package postgres

type TxManager struct {
    db *gorm.DB
}

func NewTxManager(db *gorm.DB) *TxManager {
    return &TxManager{db: db}
}

func (m *TxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
    return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        txCtx := context.WithValue(ctx, txKey{}, tx)
        return fn(txCtx)
    })
}

func (m *TxManager) WithIsolation(ctx context.Context, level application.IsolationLevel, fn func(ctx context.Context) error) error {
    var sqlLevel string
    switch level {
    case application.IsolationReadCommitted:
        sqlLevel = "READ COMMITTED"
    case application.IsolationRepeatableRead:
        sqlLevel = "REPEATABLE READ"
    case application.IsolationSerializable:
        sqlLevel = "SERIALIZABLE"
    default:
        return m.WithTransaction(ctx, fn)
    }
    return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Exec("SET TRANSACTION ISOLATION LEVEL " + sqlLevel).Error; err != nil {
            return err
        }
        txCtx := context.WithValue(ctx, txKey{}, tx)
        return fn(txCtx)
    })
}

type txKey struct{}

func TxFromContext(ctx context.Context) *gorm.DB {
    if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
        return tx
    }
    return nil
}
```

### F.3 Repository ใช้ Tx

```go
// internal/modules/device/infrastructure/persistence/postgres/device_repo_impl.go
func (r *deviceRepoImpl) Save(ctx context.Context, d *entity.Device) error {
    db := postgres.TxFromContext(ctx)
    if db == nil {
        db = r.db // fallback (ควรใช้ tx เสมอ)
    }
    m := toDeviceModel(d)
    return db.WithContext(ctx).Save(m).Error
}

func (r *deviceRepoImpl) Update(ctx context.Context, d *entity.Device) error {
    db := postgres.TxFromContext(ctx)
    if db == nil {
        db = r.db
    }
    m := toDeviceModel(d)
    return db.WithContext(ctx).
        Model(&DeviceModel{}).
        Where("id = ? AND tenant_id = ?", d.ID, d.TenantID).
        Updates(map[string]interface{}{
            "name":           m.Name,
            "status":         m.Status,
            "firmware":       m.Firmware,
            "last_seen_at":   m.LastSeenAt,
            "config":         m.Config,
            "updated_at":     time.Now().UTC(),
        }).Error
}
```

### F.4 Idempotency Pattern

```go
// internal/shared/application/idempotency.go
package application

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "time"
)

type IdempotencyStore interface {
    Get(ctx context.Context, key string) (*IdempotencyRecord, error)
    Save(ctx context.Context, key string, rec *IdempotencyRecord, ttl time.Duration) error
}

type IdempotencyRecord struct {
    Key          string          `json:"key"`
    RequestHash  string          `json:"request_hash"`
    ResponseJSON json.RawMessage `json:"response"`
    StatusCode   int             `json:"status_code"`
    CreatedAt    time.Time       `json:"created_at"`
}

type IdempotencyService struct {
    store IdempotencyStore
}

func NewIdempotencyService(store IdempotencyStore) *IdempotencyService {
    return &IdempotencyService{store: store}
}

// CheckOrExecute — ถ้ามี key เดิม + request ตรง → return cached; ถ้า key เดิม + request ต่าง → error
func (s *IdempotencyService) CheckOrExecute(
    ctx context.Context,
    key string,
    request interface{},
    ttl time.Duration,
    fn func(ctx context.Context) (interface{}, int, error),
) (interface{}, int, error) {
    if key == "" {
        return fn(ctx)
    }

    // 1. Hash request
    reqBytes, _ := json.Marshal(request)
    reqHash := sha256.Sum256(reqBytes)
    hashHex := hex.EncodeToString(reqHash[:])

    // 2. Check existing
    existing, err := s.store.Get(ctx, key)
    if err != nil {
        return nil, 0, err
    }
    if existing != nil {
        if existing.RequestHash != hashHex {
            return nil, 0, fmt.Errorf("idempotency key reused with different request")
        }
        var resp interface{}
        _ = json.Unmarshal(existing.ResponseJSON, &resp)
        return resp, existing.StatusCode, nil
    }

    // 3. Execute
    result, statusCode, err := fn(ctx)
    if err != nil {
        return nil, statusCode, err
    }

    // 4. Save
    respBytes, _ := json.Marshal(result)
    _ = s.store.Save(ctx, key, &IdempotencyRecord{
        Key:          key,
        RequestHash:  hashHex,
        ResponseJSON: respBytes,
        StatusCode:   statusCode,
        CreatedAt:    time.Now().UTC(),
    }, ttl)

    return result, statusCode, nil
}
```

### F.5 Use Case ใช้ Idempotency

```go
func (uc *RegisterDeviceUseCase) Execute(ctx context.Context, in RegisterDeviceInput) (*RegisterDeviceOutput, error) {
    if in.IdempotencyKey != "" {
        result, _, err := uc.idempotency.CheckOrExecute(
            ctx,
            in.IdempotencyKey,
            in,
            24*time.Hour,
            func(ctx context.Context) (interface{}, int, error) {
                out, err := uc.execute(ctx, in)
                return out, 201, err
            },
        )
        if err != nil {
            return nil, err
        }
        return result.(*RegisterDeviceOutput), nil
    }
    return uc.execute(ctx, in)
}
```

---

## 🅶 PART 6G — VALIDATION & ORCHESTRATION

### G.1 Two-Level Validation

```
┌──────────────────────────────────────────────────────────┐
│  Level 1: Handler Validation (defense in depth)          │
│  - binding tags (required, min, max, oneof, uuid4)       │
│  - format check                                          │
│  - return 400 ก่อนใช้ use case                            │
├──────────────────────────────────────────────────────────┤
│  Level 2: Domain Validation (in constructor / behavior)  │
│  - invariant (business rules)                            │
│  - return sentinel error                                 │
│  - use case map → 400/409/422                            │
└──────────────────────────────────────────────────────────┘
```

### G.2 Handler-Level Validation

```go
// internal/modules/device/interfaces/http/device_handler.go
func (h *DeviceHandler) Register(c *gin.Context) {
    var req application.RegisterDeviceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": err.Error(),
            },
        })
        return
    }

    userID := c.MustGet("user_id").(uuid.UUID)
    tenantID := c.MustGet("tenant_id").(uuid.UUID)
    idemKey := c.GetHeader("Idempotency-Key")

    out, err := h.registerUC.Execute(c.Request.Context(), application.RegisterDeviceInput{
        TenantID:       tenantID,
        ActorID:        userID,
        CustomerID:     uuid.MustParse(req.CustomerID),
        SiteID:         uuid.MustParse(req.SiteID),
        SerialNo:       req.SerialNo,
        Name:           req.Name,
        Type:           req.Type,
        Protocol:       req.Protocol,
        IdempotencyKey: idemKey,
    })
    if err != nil {
        status := apperrors.HTTPStatus(err)
        c.JSON(status, gin.H{"error": apperrors.ToResponse(err)})
        return
    }

    c.Header("Location", "/api/v1/devices/"+out.ID.String())
    c.JSON(http.StatusCreated, gin.H{"data": out})
}
```

### G.3 Orchestration Pattern (Saga)

```go
// internal/modules/package/application/subscribe.go
package application

// SubscribeUseCase — orchestrate: package + subscription + payment + quota
type SubscribeUseCase struct {
    pkgRepo     repository.PackageRepository
    subRepo     repository.SubscriptionRepository
    paymentSvc  PaymentService
    quotaSvc    *service.QuotaService
    txMgr       TransactionManager
    producer    EventProducer
    logger      Logger
}

type SubscribeInput struct {
    TenantID    uuid.UUID
    CustomerID  uuid.UUID
    PackageID   uuid.UUID
    Cycle       string
    PaymentRef  string
    ActorID     uuid.UUID
    IdempotencyKey string
}

type SubscribeOutput struct {
    SubscriptionID uuid.UUID `json:"subscription_id"`
    PaymentURL     string    `json:"payment_url,omitempty"`
    Status         string    `json:"status"`
}

func (uc *SubscribeUseCase) Execute(ctx context.Context, in SubscribeInput) (*SubscribeOutput, error) {
    // 1. Load package
    pkg, err := uc.pkgRepo.FindByID(ctx, in.TenantID, in.PackageID)
    if err != nil {
        return nil, fmt.Errorf("find package: %w", err)
    }
    if err := pkg.CanBeSubscribed(); err != nil {
        return nil, err
    }

    // 2. Check existing subscription
    existing, err := uc.subRepo.FindActiveByCustomer(ctx, in.TenantID, in.CustomerID)
    if err != nil && !errors.Is(err, domainerrors.ErrSubscriptionNotFound) {
        return nil, err
    }
    if existing != nil {
        return nil, domainerrors.ErrAlreadySubscribed
    }

    // 3. Create subscription (aggregate)
    cycle := vo.BillingCycle(in.Cycle)
    if !cycle.IsValid() {
        return nil, domainerrors.ErrInvalidCycle
    }
    sub, err := entity.NewSubscription(in.TenantID, in.CustomerID, pkg.ID, cycle)
    if err != nil {
        return nil, err
    }

    // 4. Create payment (external — ไม่ rollback)
    payment, err := uc.paymentSvc.CreateCheckout(ctx, PaymentRequest{
        TenantID:  in.TenantID,
        CustomerID: in.CustomerID,
        Amount:    pkg.Price,
        Reference: "SUB-" + sub.ID.String(),
        ReturnURL: in.PaymentRef,
    })
    if err != nil {
        return nil, fmt.Errorf("create payment: %w", err)
    }
    sub.PaymentID = payment.ID

    // 5. Persist (transaction)
    err = uc.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
        if err := uc.subRepo.Save(txCtx, sub); err != nil {
            return err
        }
        return uc.producer.PublishBatch(txCtx, sub.PullEvents())
    })
    if err != nil {
        // Compensate: cancel payment
        _ = uc.paymentSvc.CancelCheckout(ctx, payment.ID)
        return nil, err
    }

    return &SubscribeOutput{
        SubscriptionID: sub.ID,
        PaymentURL:     payment.CheckoutURL,
        Status:         string(sub.Status),
    }, nil
}
```

### G.4 Parallel Orchestration

```go
// Use case ที่ต้องดึงข้อมูลหลายที่พร้อมกัน
func (uc *GetDeviceDashboardUseCase) Execute(ctx context.Context, in GetDashboardInput) (*DashboardOutput, error) {
    var wg sync.WaitGroup
    var device *entity.Device
    var telemetry *entity.Telemetry
    var alerts []entity.AlertEvent
    var commands []entity.Command

    errCh := make(chan error, 4)
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    wg.Add(4)
    go func() {
        defer wg.Done()
        var err error
        device, err = uc.deviceRepo.FindByID(ctx, in.TenantID, in.DeviceID)
        if err != nil { errCh <- err }
    }()
    go func() {
        defer wg.Done()
        var err error
        telemetry, err = uc.telemetryRepo.Latest(ctx, in.TenantID, in.DeviceID, "temperature")
        if err != nil { errCh <- err }
    }()
    go func() {
        defer wg.Done()
        var err error
        alerts, err = uc.alertRepo.FindRecentByDevice(ctx, in.TenantID, in.DeviceID, 10)
        if err != nil { errCh <- err }
    }()
    go func() {
        defer wg.Done()
        var err error
        commands, err = uc.commandRepo.FindByDevice(ctx, in.TenantID, in.DeviceID, 10)
        if err != nil { errCh <- err }
    }()

    wg.Wait()
    close(errCh)
    for err := range errCh {
        if err != nil {
            return nil, err
        }
    }

    return &DashboardOutput{
        Device:    ToDeviceResponse(device),
        Telemetry: ToTelemetryPointDTO(telemetry),
        Alerts:    alerts,
        Commands:  commands,
    }, nil
}
```

---

## 🅷 PART 6H — EVENT PUBLISHING & SIDE EFFECTS

### H.1 Event Producer Port

```go
// internal/shared/application/ports.go
package application

import (
    "context"
    "icmongolang/internal/shared/domain/event"
)

type EventProducer interface {
    Publish(ctx context.Context, evt event.DomainEvent) error
    PublishBatch(ctx context.Context, events []event.DomainEvent) error
}

// EventConsumer port
type EventConsumer interface {
    Subscribe(topic string, handler EventHandler) error
    Start(ctx context.Context) error
    Close() error
}

type EventHandler func(ctx context.Context, payload []byte) error

// Cache port
type Cache interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
}

// Logger port
type Logger interface {
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
}

type Field struct {
    Key   string
    Value interface{}
}

// Metrics port
type MetricsRecorder interface {
    IncCounter(name string, labels map[string]string)
    ObserveHistogram(name string, value float64, labels map[string]string)
    IncSuccess(useCase string)
    IncError(useCase string)
}

// Clock port (สำหรับ testable time)
type Clock interface {
    Now() time.Time
}
```

### H.2 Event Producer Impl (Kafka)

```go
// internal/shared/infrastructure/kafka/producer.go
package kafka

type Producer struct {
    writer sarama.SyncProducer
    topicPrefix string
    logger application.Logger
}

func NewProducer(brokers []string, prefix string, logger application.Logger) (*Producer, error) {
    cfg := sarama.NewConfig()
    cfg.Producer.RequiredAcks = sarama.WaitForAll
    cfg.Producer.Retry.Max = 3
    cfg.Producer.Return.Successes = true
    cfg.Version = sarama.V3_0_0_0

    writer, err := sarama.NewSyncProducer(brokers, cfg)
    if err != nil {
        return nil, err
    }
    return &Producer{writer: writer, topicPrefix: prefix, logger: logger}, nil
}

func (p *Producer) Publish(ctx context.Context, evt event.DomainEvent) error {
    topic := p.topicFor(evt.EventName())
    payload, err := json.Marshal(evt)
    if err != nil {
        return fmt.Errorf("marshal event: %w", err)
    }

    msg := &sarama.ProducerMessage{
        Topic: topic,
        Key:   sarama.StringEncoder(evt.AggregateID().String()),
        Value: sarama.ByteEncoder(payload),
        Headers: []sarama.RecordHeader{
            {Key: []byte("event_name"), Value: []byte(evt.EventName())},
            {Key: []byte("event_id"), Value: []byte(evt.EventID().String())},
            {Key: []byte("occurred_at"), Value: []byte(evt.OccurredAt().Format(time.RFC3339Nano))},
            {Key: []byte("tenant_id"), Value: []byte(evt.TenantID().String())},
        },
    }

    partition, offset, err := p.writer.SendMessage(msg)
    if err != nil {
        p.logger.Error("kafka publish failed",
            application.Field{Key: "topic", Value: topic},
            application.Field{Key: "error", Value: err},
        )
        return err
    }
    p.logger.Info("event published",
        application.Field{Key: "topic", Value: topic},
        application.Field{Key: "partition", Value: partition},
        application.Field{Key: "offset", Value: offset},
    )
    return nil
}

func (p *Producer) PublishBatch(ctx context.Context, events []event.DomainEvent) error {
    if len(events) == 0 {
        return nil
    }
    msgs := make([]*sarama.ProducerMessage, 0, len(events))
    for _, evt := range events {
        payload, err := json.Marshal(evt)
        if err != nil {
            return err
        }
        msgs = append(msgs, &sarama.ProducerMessage{
            Topic: p.topicFor(evt.EventName()),
            Key:   sarama.StringEncoder(evt.AggregateID().String()),
            Value: sarama.ByteEncoder(payload),
        })
    }
    return p.writer.SendMessages(msgs)
}

// device.device.created → icmon.device.device.created
func (p *Producer) topicFor(eventName string) string {
    return p.topicPrefix + "." + eventName
}
```

### H.3 Side Effects in Use Case

```go
func (uc *RegisterDeviceUseCase) Execute(ctx context.Context, in RegisterDeviceInput) (*RegisterDeviceOutput, error) {
    // ... (business logic)

    err := uc.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
        // 1. Persist (must succeed)
        if err := uc.deviceRepo.Save(txCtx, device); err != nil {
            return fmt.Errorf("save device: %w", err)
        }

        // 2. Publish events (must succeed — same tx via outbox)
        events := device.PullEvents()
        if err := uc.producer.PublishBatch(txCtx, events); err != nil {
            return fmt.Errorf("publish events: %w", err)
        }

        return nil
    })
    if err != nil {
        return nil, err
    }

    // 3. Non-critical side effects (log warn, ไม่ fail)
    if err := uc.cache.Delete(ctx, "devices:list:"+in.TenantID.String()); err != nil {
        uc.logger.Warn("cache invalidation failed",
            application.Field{Key: "error", Value: err},
        )
    }

    if err := uc.auditLog.Record(ctx, AuditEntry{
        TenantID:  in.TenantID,
        ActorID:   in.ActorID,
        Action:    "DEVICE_REGISTERED",
        EntityID:  device.ID,
        EntityType: "device",
    }); err != nil {
        uc.logger.Warn("audit log failed",
            application.Field{Key: "error", Value: err},
        )
    }

    // 4. Metrics (non-critical)
    uc.metrics.IncSuccess("register_device")

    return &RegisterDeviceOutput{ID: device.ID, SerialNo: device.SerialNo.String()}, nil
}
```

### H.4 Outbox Pattern (สำหรับ consistency)

```go
// internal/shared/infrastructure/outbox/outbox.go
package outbox

type OutboxEntry struct {
    ID          uuid.UUID       `gorm:"type:uuid;primaryKey"`
    TenantID    uuid.UUID       `gorm:"type:uuid;index"`
    Topic       string          `gorm:"type:varchar(255);not null"`
    EventName   string          `gorm:"type:varchar(255);not null;index"`
    AggregateID uuid.UUID       `gorm:"type:uuid;index"`
    Payload     json.RawMessage `gorm:"type:jsonb;not null"`
    Status      string          `gorm:"type:varchar(20);index;default:'PENDING'"`
    RetryCount  int             `gorm:"default:0"`
    LastError   string          `gorm:"type:text"`
    CreatedAt   time.Time       `gorm:"index"`
    ProcessedAt *time.Time
}

func (OutboxEntry) TableName() string { return "outbox_entries" }

// OutboxProducer — write to same DB transaction
type OutboxProducer struct {
    db *gorm.DB
}

func (p *OutboxProducer) PublishBatch(ctx context.Context, events []event.DomainEvent) error {
    db := postgres.TxFromContext(ctx)
    if db == nil {
        db = p.db
    }
    entries := make([]*OutboxEntry, 0, len(events))
    for _, evt := range events {
        payload, err := json.Marshal(evt)
        if err != nil {
            return err
        }
        entries = append(entries, &OutboxEntry{
            ID:          uuid.New(),
            TenantID:    evt.TenantID(),
            Topic:       topicFor(evt.EventName()),
            EventName:   evt.EventName(),
            AggregateID: evt.AggregateID(),
            Payload:     payload,
            Status:      "PENDING",
            CreatedAt:   time.Now().UTC(),
        })
    }
    return db.WithContext(ctx).Create(&entries).Error
}

// OutboxRelay — background worker อ่านจาก outbox → publish Kafka → mark SENT
type OutboxRelay struct {
    db       *gorm.DB
    producer application.EventProducer
    interval time.Duration
    batchSize int
}

func (r *OutboxRelay) Run(ctx context.Context) {
    ticker := time.NewTicker(r.interval)
    defer ticker.Stop()
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            _ = r.processBatch(ctx)
        }
    }
}

func (r *OutboxRelay) processBatch(ctx context.Context) error {
    var entries []OutboxEntry
    if err := r.db.WithContext(ctx).
        Where("status = ?", "PENDING").
        Order("created_at ASC").
        Limit(r.batchSize).
        Find(&entries).Error; err != nil {
        return err
    }
    for _, entry := range entries {
        // Publish (decode payload กลับเป็น event)
        // ...
        // Mark SENT
        now := time.Now().UTC()
        r.db.WithContext(ctx).Model(&OutboxEntry{}).
            Where("id = ?", entry.ID).
            Updates(map[string]interface{}{
                "status":       "SENT",
                "processed_at": now,
            })
    }
    return nil
}
```

---

## 🅸 PART 6I — APPLICATION ERRORS

### I.1 Error Response Format

```go
// internal/shared/application/errors.go
package application

type ErrorResponse struct {
    Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
    Code      string                 `json:"code"`
    Message   string                 `json:"message"`
    Details   map[string]interface{} `json:"details,omitempty"`
    RequestID string                 `json:"request_id,omitempty"`
    Timestamp time.Time              `json:"timestamp"`
}
```

### I.2 Error Mapping Helper

```go
// internal/shared/application/errors.go (continued)
package application

import (
    "errors"
    "net/http"
    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
)

func ToResponse(err error) ErrorResponse {
    code := "INTERNAL_ERROR"
    message := "internal server error"
    details := map[string]interface{}{}

    switch {
    case errors.Is(err, domainerrors.ErrNotFound),
        errors.Is(err, domainerrors.ErrDeviceNotFound),
        errors.Is(err, domainerrors.ErrCustomerNotFound),
        errors.Is(err, domainerrors.ErrPackageNotFound),
        errors.Is(err, domainerrors.ErrOrderNotFound),
        errors.Is(err, domainerrors.ErrInvoiceNotFound):
        code = "NOT_FOUND"
        message = err.Error()

    case errors.Is(err, domainerrors.ErrAlreadyExists),
        errors.Is(err, domainerrors.ErrSerialAlreadyExists),
        errors.Is(err, domainerrors.ErrCustomerAlreadyExists):
        code = "ALREADY_EXISTS"
        message = err.Error()

    case errors.Is(err, domainerrors.ErrInvalidInput),
        errors.Is(err, domainerrors.ErrInvalidSerial),
        errors.Is(err, domainerrors.ErrInvalidDeviceType):
        code = "VALIDATION_ERROR"
        message = err.Error()

    case errors.Is(err, domainerrors.ErrUnauthorized):
        code = "UNAUTHORIZED"
        message = "authentication required"

    case errors.Is(err, domainerrors.ErrForbidden):
        code = "FORBIDDEN"
        message = "permission denied"

    case errors.Is(err, domainerrors.ErrQuotaExceeded):
        code = "QUOTA_EXCEEDED"
        message = err.Error()

    case errors.Is(err, domainerrors.ErrDeviceOffline):
        code = "DEVICE_OFFLINE"
        message = "device is offline"

    case errors.Is(err, domainerrors.ErrRateLimitExceeded):
        code = "RATE_LIMIT_EXCEEDED"
        message = "too many requests"
    }

    return ErrorResponse{
        Error: ErrorDetail{
            Code:      code,
            Message:   message,
            Details:   details,
            Timestamp: time.Now().UTC(),
        },
    }
}

func HTTPStatus(err error) int {
    switch {
    case errors.Is(err, domainerrors.ErrNotFound),
        errors.Is(err, domainerrors.ErrDeviceNotFound):
        return http.StatusNotFound
    case errors.Is(err, domainerrors.ErrAlreadyExists),
        errors.Is(err, domainerrors.ErrSerialAlreadyExists):
        return http.StatusConflict
    case errors.Is(err, domainerrors.ErrUnauthorized):
        return http.StatusUnauthorized
    case errors.Is(err, domainerrors.ErrForbidden):
        return http.StatusForbidden
    case errors.Is(err, domainerrors.ErrRateLimitExceeded):
        return http.StatusTooManyRequests
    case errors.Is(err, domainerrors.ErrQuotaExceeded):
        return http.StatusPaymentRequired
    case errors.Is(err, domainerrors.ErrDeviceOffline):
        return http.StatusServiceUnavailable
    default:
        return http.StatusBadRequest
    }
}
```

---

## 🎯 PART 6 — SUMMARY

### Application Layer Deliverables

| Category | Count | Location |
|---|:-:|---|
| **Command Use Cases** | 95 | `application/{{verb}}_{{aggregate}}.go` |
| **Query Handlers** | 64 | `application/query/` |
| **Request DTOs** | 120 | `application/dto.go` |
| **Response DTOs** | 120 | `application/dto.go` |
| **Mappers** | 40 | `application/mappers.go` |
| **Ports** | 8 | `application/ports.go` |
| **Error Types** | 60+ | `application/errors.go` |

### Layer Dependency Check

```
✅ Application Layer imports:
   - domain (entities, repos, events, errors)
   - stdlib, uuid, decimal
   - shared/application (ports)

❌ Application Layer NEVER imports:
   - gorm.io/gorm
   - github.com/gin-gonic/gin
   - github.com/IBM/sarama
   - github.com/redis/go-redis/v9
   - github.com/eclipse/paho.mqtt.golang
   - github.com/influxdata/influxdb-client-go/v2
```

---