# 🚀 SETUP + SESSION G04 + P04 — Parallel Start

> ทำทุกอย่างที่ขอในชุดเดียว: Setup → CONTEXT × 2 → Prompts × 2 → เริ่ม Session ทั้ง 2 สายพร้อมกัน

---

# PART A — SETUP `.ai/`

## A.1 Bash Commands (รันครั้งเดียว)

```bash
# === สร้างโครง ===
mkdir -p .ai/go/{modules,sessions,checkpoints,prompts}
mkdir -p .ai/py/{modules,sessions,checkpoints,prompts}
mkdir -p .ai/specs
mkdir -p .ai/checkpoints

# === ไฟล์ Go ===
touch .ai/go/CONTEXT.md
touch .ai/go/PROGRESS.md
touch .ai/go/CONVENTIONS.md
touch .ai/go/checkpoints/LAST_SESSION.md
touch .ai/go/checkpoints/NEXT_SESSION.md
touch .ai/go/prompts/generate_module.md
touch .ai/go/prompts/fix_error.md

# === ไฟล์ Python ===
touch .ai/py/CONTEXT.md
touch .ai/py/PROGRESS.md
touch .ai/py/CONVENTIONS.md
touch .ai/py/checkpoints/LAST_SESSION.md
touch .ai/py/checkpoints/NEXT_SESSION.md
touch .ai/py/prompts/generate_module.md
touch .ai/py/prompts/fix_error.md

# === Shared ===
touch .ai/MASTER_INDEX.md
touch .ai/TOKEN_LOG.md
touch .ai/specs/04_erp.yaml
```

## A.2 โครงสุดท้าย

```
.ai/
├── MASTER_INDEX.md
├── TOKEN_LOG.md
├── specs/
│   ├── 01_customer.yaml
│   ├── 02_package.yaml
│   ├── 03_device.yaml
│   └── 04_erp.yaml                 ← Day 1 target
├── go/                             🟦
│   ├── CONTEXT.md
│   ├── CONVENTIONS.md
│   ├── PROGRESS.md
│   ├── modules/
│   ├── sessions/
│   ├── checkpoints/
│   └── prompts/
└── py/                             🐍
    ├── CONTEXT.md
    ├── CONVENTIONS.md
    ├── PROGRESS.md
    ├── modules/
    ├── sessions/
    ├── checkpoints/
    └── prompts/
```

---

# PART B — CONTEXT.md

## B.1 `.ai/go/CONTEXT.md` (Go)

```markdown
# Go Context — icmongolang

## Stack
Go 1.21+ | Gin | GORM (PostgreSQL) | IBM Sarama (Kafka) | go-redis/v8 |
go-elasticsearch/v8 | influxdb-client-go/v2 | paho.mqtt.golang | google/uuid

## Structure
internal/modules/{name}/
├── domain/{entity,value_object,repository,service,event,errors}/
├── application/{ports.go,dto.go,{verb}_{entity}.go}
├── infrastructure/{persistence/postgres,persistence/redis,persistence/influxdb,
│                   messaging,iot,search/elasticsearch,scheduler}/
├── interfaces/{http,websocket}/
└── module.go

## Naming
- Entity: PascalCase
- File: snake_case
- Table: {module}_{plural}
- Migration: migrations/YYYYMMDD_{module}_init.sql
- Kafka: {context}.{aggregate}.{event}
- Redis: {module}:{entity}:{id}
- MQTT: iot/{serial}/{suffix}

## Layer Rules
| Layer | Import | Forbidden |
|-------|--------|-----------|
| domain | stdlib + uuid | gorm gin sarama redis es |
| application | domain | infrastructure interface |
| infrastructure | domain + app | interface |
| interface | all | – |

## Patterns

### Entity
```go
type Customer struct { ID uuid.UUID; ... }
func NewCustomer(tenantID uuid.UUID, code, name string) (*Customer, error)
func (c *Customer) Activate() error       // behavior
func (c *Customer) IsActive() bool        // query
```

### Use Case
```go
type CreateCustomerUseCase struct { repo repository.CustomerRepository }
type CreateCustomerInput struct { ... }
type CreateCustomerOutput struct { ... }
func (uc *CreateCustomerUseCase) Execute(ctx context.Context, in CreateCustomerInput) (*CreateCustomerOutput, error)
```

### Repository
- Interface: domain/repository/{entity}_repository.go
- Impl: infrastructure/persistence/postgres/{entity}_repo_impl.go
- Mapper: toEntity(), toModel()

### Error
```go
var (
    ErrCustomerNotFound = errors.New("customer not found")
    ErrInvalidCode      = errors.New("invalid customer code")
)
```

### Port
```go
type EventProducer interface {
    PublishCustomerCreated(ctx context.Context, evt event.CustomerCreated) error
}
```

## pkg/ Whitelist
cryptpass db elasticsearch emailTemplates helpers http-swagger
httpErrors influxdb jwt kafka llm logger mqtt report responses
secureRandom sendEmail transaction utils vectordb websocket

## Build/Test
go build ./... && go vet ./... && go test ./...

## Reference
Structure: internal/modules/device/
Patterns:  internal/modules/customer/
```

## B.2 `.ai/py/CONTEXT.md` (Python)

```markdown
# Python Context — icmongolang

## Stack
Python 3.11+ | FastAPI | Pydantic v2 | SQLAlchemy 2.0 async + asyncpg |
aiokafka | redis.asyncio | elasticsearch-py async | influxdb-client async |
aiomqtt | pytest + pytest-asyncio

## Structure
app/modules/{name}/
├── domain/{entities,value_objects,repositories,services,events,errors}/
├── application/{use_cases,ports}/
├── infrastructure/{persistence/postgres,persistence/redis,persistence/influxdb,
│                   messaging,iot,search/elasticsearch,scheduler}/
├── interfaces/{http,middleware}/
└── module.py

## Naming
- Entity: PascalCase
- File: snake_case
- Table: {module}_{plural}
- Migration: alembic/versions/YYYYMMDD_{module}_init.py
- Kafka: {context}.{aggregate}.{event}
- Redis: {module}:{entity}:{id}

## Layer Rules
| Layer | Import | Forbidden |
|-------|--------|-----------|
| domain | stdlib + pydantic | sqlalchemy fastapi aiokafka |
| application | domain | infrastructure interface |
| infrastructure | domain + app | interface |
| interface | all | – |

## Patterns

### Entity
```python
class Customer(BaseModel):
    id: UUID = Field(default_factory=uuid4)
    tenant_id: UUID
    code: str
    name: str
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    @classmethod
    def create(cls, tenant_id: UUID, code: str, name: str) -> "Customer": ...

    def activate(self) -> None: ...      # behavior
    def is_active(self) -> bool: ...     # query
```

### Use Case
```python
@dataclass
class CreateCustomerInput:
    tenant_id: UUID
    code: str

@dataclass
class CreateCustomerOutput:
    customer_id: UUID

class CreateCustomerUseCase:
    def __init__(self, repo: CustomerRepository, producer: EventProducer): ...
    async def execute(self, inp: CreateCustomerInput) -> CreateCustomerOutput: ...
```

### Repository (Protocol)
```python
class CustomerRepository(Protocol):
    async def save(self, c: Customer) -> None: ...
    async def find_by_id(self, id: UUID) -> Customer | None: ...
```

### Error
```python
class DomainError(Exception):
    code: str = "DOMAIN_ERROR"

class CustomerNotFoundError(DomainError):
    code = "CUSTOMER_NOT_FOUND"
```

### Port (Protocol)
```python
class EventProducer(Protocol):
    async def publish_customer_created(self, evt: CustomerCreatedEvent) -> None: ...
```

### FastAPI Router
```python
@router.post("", response_model=CustomerResponse, status_code=201)
async def create_customer(
    req: CreateCustomerRequest,
    uc: CreateCustomerUseCase = Depends(get_create_uc),
    user_id: UUID = Depends(get_current_user),
):
    try:
        return await uc.execute(CreateCustomerInput(...))
    except DomainError as e:
        raise HTTPException(400, detail={"code": e.code, "message": str(e)})
```

## Build/Test
python -m compileall app/ && pytest tests/

## Reference
Structure: app/modules/device/
Patterns:  app/modules/customer/
```

---

# PART C — PROMPTS

## C.1 `.ai/go/prompts/generate_module.md`

```markdown
# Generate Go Module

## Context
{{paste .ai/go/CONTEXT.md}}

## Spec
{{paste .ai/specs/{NN}_{name}.yaml}}

## Task
Generate complete Go module per spec.

## Output Protocol
Reply **"ready"** → wait for **"start"**.

Then generate in 5 steps. After each step, reply "OK {step}. Say continue."
I'll say "continue" between steps.

### Step 1 — Domain Layer
Files:
- domain/errors/errors.go
- domain/value_object/*.go
- domain/entity/*.go  (entities + behavior methods)
- domain/event/events.go
- domain/repository/*.go  (interfaces only)
- domain/service/*.go  (if needed)

### Step 2 — Application Layer
Files:
- application/ports.go
- application/dto.go
- application/{verb}_{entity}.go  (1 file per use case)

### Step 3 — Infrastructure Layer
Files:
- infrastructure/persistence/postgres/models.go
- infrastructure/persistence/postgres/*_repo_impl.go
- infrastructure/messaging/producer.go
- infrastructure/{others per spec}

### Step 4 — Interface Layer + module.go
Files:
- interfaces/http/*_handler.go
- interfaces/http/routes.go
- interfaces/http/dto.go
- module.go

### Step 5 — Migration + Tests
Files:
- migrations/YYYYMMDD_{name}_init.sql
- domain/entity/{entity}_test.go

## Rules
1. Domain: NO gorm/gin/sarama imports
2. Entities: constructor + behavior methods (no setters)
3. Repository: interface in domain, impl in infrastructure
4. Use case: Execute(ctx, input) (output, error)
5. Errors: sentinel errors
6. Complete code, do NOT truncate
7. Reference structure: internal/modules/device/
8. Reference patterns: internal/modules/customer/

## After generation
Run `go build ./...` and report errors.
```

## C.2 `.ai/py/prompts/generate_module.md`

```markdown
# Generate Python Module

## Context
{{paste .ai/py/CONTEXT.md}}

## Spec
{{paste .ai/specs/{NN}_{name}.yaml}}

## Task
Generate complete Python FastAPI module per spec.

## Output Protocol
Reply **"ready"** → wait for **"start"**.

Then generate in 5 steps. After each step, reply "OK {step}. Say continue."

### Step 1 — Domain Layer
Files:
- domain/errors/errors.py
- domain/value_objects/*.py
- domain/entities/*.py
- domain/events/events.py
- domain/repositories/*.py  (Protocol)
- domain/services/*.py  (if needed)

### Step 2 — Application Layer
Files:
- application/ports/*.py
- application/use_cases/*.py

### Step 3 — Infrastructure Layer
Files:
- infrastructure/persistence/postgres/models.py
- infrastructure/persistence/postgres/*_repository_impl.py
- infrastructure/messaging/kafka_producer.py
- infrastructure/{others per spec}

### Step 4 — Interface Layer + module.py
Files:
- interfaces/http/schemas.py
- interfaces/http/*_router.py
- interfaces/http/dependencies.py
- module.py

### Step 5 — Migration + Tests
Files:
- alembic/versions/YYYYMMDD_{name}_init.py
- tests/modules/{name}/test_{entity}_entity.py

## Rules
1. Domain: NO sqlalchemy/fastapi/aiokafka imports
2. Entities: Pydantic BaseModel + behavior methods
3. Repository: Protocol in domain, impl in infrastructure
4. Use case: async execute(input) -> output
5. Errors: exception classes with `code` attribute
6. Complete code, do NOT truncate
7. Type hints everywhere
8. Reference structure: app/modules/device/
9. Reference patterns: app/modules/customer/

## After generation
Run `python -m compileall app/` and `pytest tests/`.
```

---

# PART D — SPEC `.ai/specs/04_erp.yaml`

```yaml
module: erp
context: Enterprise Resource Planning
priority: P1
status: queued

aggregates:
  Order:
    fields:
      - id: UUID
      - tenant_id: UUID
      - customer_id: UUID
      - order_number: string(50)
      - type: OrderType
      - items: []OrderItem
      - subtotal: Money
      - tax_amount: Money
      - discount: Money
      - total: Money
      - status: OrderStatus
      - notes: text
      - shipping_addr: Address
      - billing_addr: Address
      - ordered_at: timestamp
      - confirmed_at: *timestamp
      - shipped_at: *timestamp
      - delivered_at: *timestamp
      - cancelled_at: *timestamp
      - cancel_reason: text
      - created_by: UUID
    behavior: [AddItem, RemoveItem, Recalculate, Confirm, Ship, Deliver, Cancel]
    invariants:
      - total = subtotal - discount + tax
      - must have items before confirm
      - cannot cancel delivered

  OrderItem:
    parent: Order
    fields: [id, order_id, product_id, sku, name, quantity, unit_price, discount, tax_rate, subtotal]

  Invoice:
    fields:
      - id, tenant_id, order_id?, customer_id
      - number: string(50)
      - type: InvoiceType
      - subtotal, tax_amount, total: Money
      - status: InvoiceStatus
      - issued_at, due_date, paid_at: timestamp
      - paid_amount: Money
      - payment_method: string
    behavior: [MarkPaid, MarkOverdue, Cancel]
    invariants:
      - total = subtotal + tax_amount
      - cannot pay if already paid
      - due_date > issued_at

  Warehouse:
    fields: [id, tenant_id, code, name, address, type, capacity, manager_id, is_active]

  StockItem:
    fields: [id, tenant_id, product_id, warehouse_id, sku, quantity, reserved, reorder_point, unit_cost]
    behavior: [Reserve, Release, Issue, Restock, NeedsReorder]
    invariants:
      - reserved <= quantity
      - available = quantity - reserved

value_objects:
  Money: {amount: float, currency: string}
  OrderType: [SALES, PURCHASE]
  OrderStatus: [DRAFT, CONFIRMED, SHIPPED, DELIVERED, CANCELLED]
  InvoiceType: [TAX, RECEIPT, DEBIT_NOTE, CREDIT_NOTE]
  InvoiceStatus: [DRAFT, ISSUED, PAID, OVERDUE, CANCELLED]
  WarehouseType: [MAIN, BRANCH, TRANSIT, COLD]
  Address: {line1, line2, district, amphoe, province, postcode, country}

use_cases:
  - CreateOrder
  - GetOrder
  - ListOrders
  - AddOrderItem
  - RemoveOrderItem
  - ConfirmOrder
  - ShipOrder
  - DeliverOrder
  - CancelOrder
  - CreateInvoice
  - IssueInvoice
  - RecordPayment
  - GetInvoice
  - ListInvoices
  - CreateWarehouse
  - ListWarehouses
  - RestockItem
  - CheckStock
  - ReorderCheck       # scheduled
  - InvoiceOverdueCheck # scheduled

tables:
  - erp_orders
  - erp_order_items
  - erp_invoices
  - erp_invoice_items
  - erp_warehouses
  - erp_stock_items

kafka_publish:
  - erp.order.created
  - erp.order.confirmed
  - erp.order.shipped
  - erp.order.delivered
  - erp.order.cancelled
  - erp.invoice.issued
  - erp.invoice.paid
  - erp.stock.low
  - erp.stock.reordered

kafka_consume:
  - payment.payment.succeeded    # mark invoice paid
  - logistics.shipment.delivered # mark order delivered
  - customer.customer.suspended  # block new orders

sync_apis_out:
  - GET customer /internal/customers/{id}
  - GET items /internal/items/{id}
  - POST payment /internal/payments

ports_needed:
  - CustomerClient
  - ItemClient
  - PaymentClient
  - EventProducer
  - AuditRepository
  - WSHub

scheduler_jobs:
  - InvoiceOverdueJob   # daily 3am
  - ReorderCheckJob     # daily 6am

migration: 20260104_erp_init.sql
references:
  full_doc: docs/modules/04_erp.md
```

---

# PART E — SESSION G04 (Go ERP) — Step 1 Domain

> **Instruction:** ใน chat ใหม่ paste CONTEXT.md + spec + prompt แล้ว AI จะ generate ตามนี้

## E.1 Directory Tree

```
internal/modules/erp/
├── domain/
│   ├── entity/
│   │   ├── order.go
│   │   ├── order_item.go
│   │   ├── invoice.go
│   │   ├── warehouse.go
│   │   └── stock_item.go
│   ├── value_object/
│   │   ├── order_type.go
│   │   ├── order_status.go
│   │   ├── invoice_type.go
│   │   ├── invoice_status.go
│   │   ├── warehouse_type.go
│   │   ├── money.go
│   │   └── address.go
│   ├── repository/
│   │   ├── order_repository.go
│   │   ├── invoice_repository.go
│   │   ├── warehouse_repository.go
│   │   └── stock_repository.go
│   ├── event/
│   │   └── events.go
│   └── errors/
│       └── errors.go
```

## E.2 `domain/value_object/money.go`

```go
package valueobject

import "errors"

type Money struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

func NewMoney(amount float64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("money amount cannot be negative")
	}
	if currency == "" {
		return Money{}, errors.New("currency required")
	}
	return Money{Amount: amount, Currency: currency}, nil
}

func (m Money) Add(o Money) Money {
	if m.Currency != o.Currency {
		return m
	}
	return Money{Amount: m.Amount + o.Amount, Currency: m.Currency}
}

func (m Money) Sub(o Money) Money {
	if m.Currency != o.Currency {
		return m
	}
	return Money{Amount: m.Amount - o.Amount, Currency: m.Currency}
}

func (m Money) Mul(f float64) Money {
	return Money{Amount: m.Amount * f, Currency: m.Currency}
}

func (m Money) IsZero() bool  { return m.Amount == 0 }
func (m Money) IsNegative() bool { return m.Amount < 0 }
```

## E.3 `domain/value_object/order_type.go`

```go
package valueobject

type OrderType string

const (
	OrderTypeSales    OrderType = "SALES"
	OrderTypePurchase OrderType = "PURCHASE"
)

func (t OrderType) IsValid() bool {
	switch t {
	case OrderTypeSales, OrderTypePurchase:
		return true
	}
	return false
}
```

## E.4 `domain/value_object/order_status.go`

```go
package valueobject

type OrderStatus string

const (
	OrderStatusDraft     OrderStatus = "DRAFT"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusShipped   OrderStatus = "SHIPPED"
	OrderStatusDelivered OrderStatus = "DELIVERED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusDraft, OrderStatusConfirmed, OrderStatusShipped,
		OrderStatusDelivered, OrderStatusCancelled:
		return true
	}
	return false
}

func (s OrderStatus) IsTerminal() bool {
	return s == OrderStatusDelivered || s == OrderStatusCancelled
}

func (s OrderStatus) CanTransitionTo(next OrderStatus) bool {
	transitions := map[OrderStatus]map[OrderStatus]bool{
		OrderStatusDraft:     {OrderStatusConfirmed: true, OrderStatusCancelled: true},
		OrderStatusConfirmed: {OrderStatusShipped: true, OrderStatusCancelled: true},
		OrderStatusShipped:   {OrderStatusDelivered: true, OrderStatusCancelled: true},
		OrderStatusDelivered: {},
		OrderStatusCancelled: {},
	}
	return transitions[s][next]
}
```

## E.5 `domain/value_object/invoice_type.go`

```go
package valueobject

type InvoiceType string

const (
	InvoiceTypeTax        InvoiceType = "TAX"
	InvoiceTypeReceipt    InvoiceType = "RECEIPT"
	InvoiceTypeDebitNote  InvoiceType = "DEBIT_NOTE"
	InvoiceTypeCreditNote InvoiceType = "CREDIT_NOTE"
)

func (t InvoiceType) IsValid() bool {
	switch t {
	case InvoiceTypeTax, InvoiceTypeReceipt, InvoiceTypeDebitNote, InvoiceTypeCreditNote:
		return true
	}
	return false
}
```

## E.6 `domain/value_object/invoice_status.go`

```go
package valueobject

type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "DRAFT"
	InvoiceStatusIssued    InvoiceStatus = "ISSUED"
	InvoiceStatusPaid      InvoiceStatus = "PAID"
	InvoiceStatusOverdue   InvoiceStatus = "OVERDUE"
	InvoiceStatusCancelled InvoiceStatus = "CANCELLED"
)

func (s InvoiceStatus) IsValid() bool {
	switch s {
	case InvoiceStatusDraft, InvoiceStatusIssued, InvoiceStatusPaid,
		InvoiceStatusOverdue, InvoiceStatusCancelled:
		return true
	}
	return false
}

func (s InvoiceStatus) IsTerminal() bool {
	return s == InvoiceStatusPaid || s == InvoiceStatusCancelled
}

func (s InvoiceStatus) CanTransitionTo(next InvoiceStatus) bool {
	transitions := map[InvoiceStatus]map[InvoiceStatus]bool{
		InvoiceStatusDraft:   {InvoiceStatusIssued: true, InvoiceStatusCancelled: true},
		InvoiceStatusIssued:  {InvoiceStatusPaid: true, InvoiceStatusOverdue: true, InvoiceStatusCancelled: true},
		InvoiceStatusOverdue: {InvoiceStatusPaid: true, InvoiceStatusCancelled: true},
		InvoiceStatusPaid:    {},
		InvoiceStatusCancelled: {},
	}
	return transitions[s][next]
}
```

## E.7 `domain/value_object/warehouse_type.go`

```go
package valueobject

type WarehouseType string

const (
	WarehouseTypeMain    WarehouseType = "MAIN"
	WarehouseTypeBranch  WarehouseType = "BRANCH"
	WarehouseTypeTransit WarehouseType = "TRANSIT"
	WarehouseTypeCold    WarehouseType = "COLD"
)

func (t WarehouseType) IsValid() bool {
	switch t {
	case WarehouseTypeMain, WarehouseTypeBranch, WarehouseTypeTransit, WarehouseTypeCold:
		return true
	}
	return false
}
```

## E.8 `domain/errors/errors.go`

```go
package domainerrors

import "errors"

var (
	// Validation
	ErrInvalidTenant    = errors.New("invalid tenant")
	ErrInvalidCustomer  = errors.New("invalid customer")
	ErrInvalidOrderType = errors.New("invalid order type")
	ErrInvalidQuantity  = errors.New("invalid quantity (must be > 0)")
	ErrInvalidPrice     = errors.New("invalid price (must be >= 0)")
	ErrInvalidOrderNum  = errors.New("invalid order number")
	ErrInvalidInvoiceType = errors.New("invalid invoice type")
	ErrInvalidWarehouseType = errors.New("invalid warehouse type")
	ErrInvalidDueDate   = errors.New("due date must be after issue date")

	// Not found
	ErrOrderNotFound     = errors.New("order not found")
	ErrOrderItemNotFound = errors.New("order item not found")
	ErrInvoiceNotFound   = errors.New("invoice not found")
	ErrWarehouseNotFound = errors.New("warehouse not found")
	ErrStockNotFound     = errors.New("stock item not found")

	// Conflict
	ErrOrderNumberDuplicate   = errors.New("order number already exists")
	ErrInvoiceNumberDuplicate = errors.New("invoice number already exists")
	ErrWarehouseCodeDuplicate = errors.New("warehouse code already exists")

	// State machine
	ErrOrderNotDraft     = errors.New("order is not in draft state")
	ErrOrderNotConfirmed = errors.New("order is not confirmed")
	ErrOrderNotShipped   = errors.New("order is not shipped")
	ErrOrderEmpty        = errors.New("order has no items")
	ErrCannotCancelDelivered = errors.New("cannot cancel delivered order")
	ErrInvalidOrderTransition = errors.New("invalid order status transition")
	ErrInvoiceAlreadyPaid = errors.New("invoice already paid")
	ErrCannotCancelPaid   = errors.New("cannot cancel paid invoice")
	ErrInvalidInvoiceTransition = errors.New("invalid invoice status transition")

	// Business rules
	ErrInsufficientStock   = errors.New("insufficient stock")
	ErrReservationMismatch = errors.New("reservation mismatch")
	ErrInsufficientPayment = errors.New("insufficient payment")
	ErrCurrencyMismatch    = errors.New("currency mismatch")
)
```

## E.9 `domain/entity/order_item.go`

```go
package entity

import (
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/erp/domain/errors"
	valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type OrderItem struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	ProductID uuid.UUID
	SKU       string
	Name      string
	Quantity  int
	UnitPrice valueobject.Money
	Discount  valueobject.Money
	TaxRate   float64
	Subtotal  valueobject.Money
	CreatedAt time.Time
}

func NewOrderItem(
	productID uuid.UUID,
	sku, name string,
	qty int,
	unitPrice valueobject.Money,
) (*OrderItem, error) {
	if productID == uuid.Nil {
		return nil, domainerrors.ErrInvalidQuantity
	}
	if qty <= 0 {
		return nil, domainerrors.ErrInvalidQuantity
	}
	if unitPrice.Amount < 0 {
		return nil, domainerrors.ErrInvalidPrice
	}
	subtotal := unitPrice.Mul(float64(qty))
	return &OrderItem{
		ID:        uuid.New(),
		ProductID: productID,
		SKU:       sku,
		Name:      name,
		Quantity:  qty,
		UnitPrice: unitPrice,
		Discount:  valueobject.Money{Currency: unitPrice.Currency},
		Subtotal:  subtotal,
		CreatedAt: time.Now(),
	}, nil
}

func (i *OrderItem) UpdateQuantity(qty int) error {
	if qty <= 0 {
		return domainerrors.ErrInvalidQuantity
	}
	i.Quantity = qty
	i.recalcSubtotal()
	return nil
}

func (i *OrderItem) UpdateUnitPrice(price valueobject.Money) error {
	if price.Amount < 0 {
		return domainerrors.ErrInvalidPrice
	}
	if price.Currency != i.UnitPrice.Currency {
		return domainerrors.ErrCurrencyMismatch
	}
	i.UnitPrice = price
	i.recalcSubtotal()
	return nil
}

func (i *OrderItem) ApplyDiscount(discount valueobject.Money) error {
	if discount.Amount < 0 {
		return domainerrors.ErrInvalidPrice
	}
	if discount.Currency != i.Subtotal.Currency {
		return domainerrors.ErrCurrencyMismatch
	}
	i.Discount = discount
	return nil
}

func (i *OrderItem) recalcSubtotal() {
	i.Subtotal = i.UnitPrice.Mul(float64(i.Quantity))
}
```

## E.10 `domain/entity/order.go`

```go
package entity

import (
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/erp/domain/errors"
	valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type Order struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	CustomerID   uuid.UUID
	OrderNumber  string
	Type         valueobject.OrderType
	Items        []OrderItem
	Subtotal     valueobject.Money
	TaxAmount    valueobject.Money
	Discount     valueobject.Money
	Total        valueobject.Money
	Status       valueobject.OrderStatus
	Notes        string
	ShippingAddr valueobject.Address
	BillingAddr  valueobject.Address
	OrderedAt    time.Time
	ConfirmedAt  *time.Time
	ShippedAt    *time.Time
	DeliveredAt  *time.Time
	CancelledAt  *time.Time
	CancelReason string
	CreatedBy    uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewOrder(
	tenantID, customerID uuid.UUID,
	otype valueobject.OrderType,
	createdBy uuid.UUID,
) (*Order, error) {
	if tenantID == uuid.Nil {
		return nil, domainerrors.ErrInvalidTenant
	}
	if customerID == uuid.Nil {
		return nil, domainerrors.ErrInvalidCustomer
	}
	if !otype.IsValid() {
		return nil, domainerrors.ErrInvalidOrderType
	}
	now := time.Now()
	return &Order{
		ID:         uuid.New(),
		TenantID:   tenantID,
		CustomerID: customerID,
		Type:       otype,
		Items:      []OrderItem{},
		Status:     valueobject.OrderStatusDraft,
		OrderedAt:  now,
		CreatedBy:  createdBy,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// ---------- Behavior ----------

func (o *Order) SetNumber(number string) error {
	if number == "" || len(number) > 50 {
		return domainerrors.ErrInvalidOrderNum
	}
	o.OrderNumber = number
	o.touch()
	return nil
}

func (o *Order) AddItem(item OrderItem) error {
	if o.Status != valueobject.OrderStatusDraft {
		return domainerrors.ErrOrderNotDraft
	}
	item.OrderID = o.ID
	o.Items = append(o.Items, item)
	o.recalculate()
	o.touch()
	return nil
}

func (o *Order) RemoveItem(itemID uuid.UUID) error {
	if o.Status != valueobject.OrderStatusDraft {
		return domainerrors.ErrOrderNotDraft
	}
	for i, item := range o.Items {
		if item.ID == itemID {
			o.Items = append(o.Items[:i], o.Items[i+1:]...)
			o.recalculate()
			o.touch()
			return nil
		}
	}
	return domainerrors.ErrOrderItemNotFound
}

func (o *Order) Confirm() error {
	if !o.Status.CanTransitionTo(valueobject.OrderStatusConfirmed) {
		return domainerrors.ErrInvalidOrderTransition
	}
	if len(o.Items) == 0 {
		return domainerrors.ErrOrderEmpty
	}
	now := time.Now()
	o.Status = valueobject.OrderStatusConfirmed
	o.ConfirmedAt = &now
	o.touch()
	return nil
}

func (o *Order) Ship() error {
	if !o.Status.CanTransitionTo(valueobject.OrderStatusShipped) {
		return domainerrors.ErrInvalidOrderTransition
	}
	now := time.Now()
	o.Status = valueobject.OrderStatusShipped
	o.ShippedAt = &now
	o.touch()
	return nil
}

func (o *Order) Deliver() error {
	if !o.Status.CanTransitionTo(valueobject.OrderStatusDelivered) {
		return domainerrors.ErrInvalidOrderTransition
	}
	now := time.Now()
	o.Status = valueobject.OrderStatusDelivered
	o.DeliveredAt = &now
	o.touch()
	return nil
}

func (o *Order) Cancel(reason string) error {
	if !o.Status.CanTransitionTo(valueobject.OrderStatusCancelled) {
		return domainerrors.ErrCannotCancelDelivered
	}
	now := time.Now()
	o.Status = valueobject.OrderStatusCancelled
	o.CancelledAt = &now
	o.CancelReason = reason
	o.touch()
	return nil
}

func (o *Order) SetShippingAddress(addr valueobject.Address) error {
	if err := addr.Validate(); err != nil {
		return err
	}
	o.ShippingAddr = addr
	o.touch()
	return nil
}

func (o *Order) SetBillingAddress(addr valueobject.Address) error {
	if err := addr.Validate(); err != nil {
		return err
	}
	o.BillingAddr = addr
	o.touch()
	return nil
}

// ---------- Query ----------

func (o *Order) IsDraft() bool     { return o.Status == valueobject.OrderStatusDraft }
func (o *Order) IsConfirmed() bool { return o.Status == valueobject.OrderStatusConfirmed }
func (o *Order) IsDelivered() bool { return o.Status == valueobject.OrderStatusDelivered }
func (o *Order) IsCancelled() bool { return o.Status == valueobject.OrderStatusCancelled }
func (o *Order) IsSales() bool     { return o.Type == valueobject.OrderTypeSales }
func (o *Order) ItemCount() int    { return len(o.Items) }

// ---------- Private ----------

func (o *Order) recalculate() {
	var subtotal, tax, discount float64
	currency := "THB"
	if len(o.Items) > 0 {
		currency = o.Items[0].Subtotal.Currency
	}
	for _, item := range o.Items {
		subtotal += item.Subtotal.Amount
		discount += item.Discount.Amount
		tax += item.Subtotal.Amount * item.TaxRate / 100
	}
	o.Subtotal = valueobject.Money{Amount: subtotal, Currency: currency}
	o.Discount = valueobject.Money{Amount: discount, Currency: currency}
	o.TaxAmount = valueobject.Money{Amount: tax, Currency: currency}
	o.Total = valueobject.Money{Amount: subtotal - discount + tax, Currency: currency}
}

func (o *Order) touch() { o.UpdatedAt = time.Now() }
```

## E.11 `domain/entity/invoice.go`

```go
package entity

import (
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/erp/domain/errors"
	valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type Invoice struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	OrderID       *uuid.UUID
	CustomerID    uuid.UUID
	Number        string
	Type          valueobject.InvoiceType
	Subtotal      valueobject.Money
	TaxAmount     valueobject.Money
	Total         valueobject.Money
	Status        valueobject.InvoiceStatus
	IssuedAt      time.Time
	DueDate       time.Time
	PaidAt        *time.Time
	PaidAmount    valueobject.Money
	PaymentMethod string
	Notes         string
	CreatedBy     uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewInvoice(
	tenantID, customerID uuid.UUID,
	itype valueobject.InvoiceType,
	subtotal, tax, total valueobject.Money,
	dueDate time.Time,
	createdBy uuid.UUID,
) (*Invoice, error) {
	if tenantID == uuid.Nil {
		return nil, domainerrors.ErrInvalidTenant
	}
	if customerID == uuid.Nil {
		return nil, domainerrors.ErrInvalidCustomer
	}
	if !itype.IsValid() {
		return nil, domainerrors.ErrInvalidInvoiceType
	}
	now := time.Now()
	if !dueDate.After(now) {
		return nil, domainerrors.ErrInvalidDueDate
	}
	return &Invoice{
		ID:         uuid.New(),
		TenantID:   tenantID,
		CustomerID: customerID,
		Type:       itype,
		Subtotal:   subtotal,
		TaxAmount:  tax,
		Total:      total,
		Status:     valueobject.InvoiceStatusDraft,
		IssuedAt:   now,
		DueDate:    dueDate,
		CreatedBy:  createdBy,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// ---------- Behavior ----------

func (i *Invoice) SetNumber(number string) error {
	if number == "" || len(number) > 50 {
		return domainerrors.ErrInvalidOrderNum
	}
	i.Number = number
	i.touch()
	return nil
}

func (i *Invoice) LinkOrder(orderID uuid.UUID) {
	i.OrderID = &orderID
	i.touch()
}

func (i *Invoice) Issue() error {
	if !i.Status.CanTransitionTo(valueobject.InvoiceStatusIssued) {
		return domainerrors.ErrInvalidInvoiceTransition
	}
	i.Status = valueobject.InvoiceStatusIssued
	i.IssuedAt = time.Now()
	i.touch()
	return nil
}

func (i *Invoice) MarkPaid(amount valueobject.Money, method string) error {
	if i.Status == valueobject.InvoiceStatusPaid {
		return domainerrors.ErrInvoiceAlreadyPaid
	}
	if amount.Amount < i.Total.Amount {
		return domainerrors.ErrInsufficientPayment
	}
	if amount.Currency != i.Total.Currency {
		return domainerrors.ErrCurrencyMismatch
	}
	now := time.Now()
	i.Status = valueobject.InvoiceStatusPaid
	i.PaidAt = &now
	i.PaidAmount = amount
	i.PaymentMethod = method
	i.touch()
	return nil
}

func (i *Invoice) MarkOverdue() {
	if i.Status == valueobject.InvoiceStatusIssued && time.Now().After(i.DueDate) {
		i.Status = valueobject.InvoiceStatusOverdue
		i.touch()
	}
}

func (i *Invoice) Cancel() error {
	if i.Status == valueobject.InvoiceStatusPaid {
		return domainerrors.ErrCannotCancelPaid
	}
	if !i.Status.CanTransitionTo(valueobject.InvoiceStatusCancelled) {
		return domainerrors.ErrInvalidInvoiceTransition
	}
	i.Status = valueobject.InvoiceStatusCancelled
	i.touch()
	return nil
}

// ---------- Query ----------

func (i *Invoice) IsPaid() bool    { return i.Status == valueobject.InvoiceStatusPaid }
func (i *Invoice) IsIssued() bool  { return i.Status == valueobject.InvoiceStatusIssued }
func (i *Invoice) IsOverdue() bool { return i.Status == valueobject.InvoiceStatusOverdue }

func (i *Invoice) DaysUntilDue() int {
	return int(time.Until(i.DueDate).Hours() / 24)
}

func (i *Invoice) touch() { i.UpdatedAt = time.Now() }
```

## E.12 `domain/entity/warehouse.go`

```go
package entity

import (
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/erp/domain/errors"
	valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type Warehouse struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	Code      string
	Name      string
	Address   valueobject.Address
	Type      valueobject.WarehouseType
	Capacity  int
	ManagerID *uuid.UUID
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewWarehouse(
	tenantID uuid.UUID,
	code, name string,
	wtype valueobject.WarehouseType,
) (*Warehouse, error) {
	if tenantID == uuid.Nil {
		return nil, domainerrors.ErrInvalidTenant
	}
	if code == "" || len(code) > 50 {
		return nil, domainerrors.ErrInvalidOrderNum
	}
	if name == "" || len(name) > 255 {
		return nil, domainerrors.ErrInvalidOrderNum
	}
	if !wtype.IsValid() {
		return nil, domainerrors.ErrInvalidWarehouseType
	}
	now := time.Now()
	return &Warehouse{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Code:      code,
		Name:      name,
		Type:      wtype,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (w *Warehouse) Rename(name string) error {
	if name == "" || len(name) > 255 {
		return domainerrors.ErrInvalidOrderNum
	}
	w.Name = name
	w.touch()
	return nil
}

func (w *Warehouse) SetAddress(addr valueobject.Address) error {
	if err := addr.Validate(); err != nil {
		return err
	}
	w.Address = addr
	w.touch()
	return nil
}

func (w *Warehouse) SetCapacity(capacity int) error {
	if capacity < 0 {
		return domainerrors.ErrInvalidQuantity
	}
	w.Capacity = capacity
	w.touch()
	return nil
}

func (w *Warehouse) AssignManager(managerID uuid.UUID) {
	w.ManagerID = &managerID
	w.touch()
}

func (w *Warehouse) Activate()   { w.IsActive = true; w.touch() }
func (w *Warehouse) Deactivate() { w.IsActive = false; w.touch() }

func (w *Warehouse) touch() { w.UpdatedAt = time.Now() }
```

## E.13 `domain/entity/stock_item.go`

```go
package entity

import (
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/erp/domain/errors"
	valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type StockItem struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	ProductID     uuid.UUID
	WarehouseID   uuid.UUID
	SKU           string
	Quantity      int
	Reserved      int
	ReorderPoint  int
	UnitCost      valueobject.Money
	LastRestocked *time.Time
	UpdatedAt     time.Time
}

func NewStockItem(
	tenantID, productID, warehouseID uuid.UUID,
	sku string,
	unitCost valueobject.Money,
) (*StockItem, error) {
	if tenantID == uuid.Nil {
		return nil, domainerrors.ErrInvalidTenant
	}
	if productID == uuid.Nil || warehouseID == uuid.Nil {
		return nil, domainerrors.ErrInvalidQuantity
	}
	if sku == "" {
		return nil, domainerrors.ErrInvalidOrderNum
	}
	return &StockItem{
		ID:          uuid.New(),
		TenantID:    tenantID,
		ProductID:   productID,
		WarehouseID: warehouseID,
		SKU:         sku,
		UnitCost:    unitCost,
		UpdatedAt:   time.Now(),
	}, nil
}

// ---------- Behavior ----------

func (s *StockItem) Reserve(qty int) error {
	if qty <= 0 {
		return domainerrors.ErrInvalidQuantity
	}
	if s.Available() < qty {
		return domainerrors.ErrInsufficientStock
	}
	s.Reserved += qty
	s.touch()
	return nil
}

func (s *StockItem) Release(qty int) {
	s.Reserved -= qty
	if s.Reserved < 0 {
		s.Reserved = 0
	}
	s.touch()
}

func (s *StockItem) Issue(qty int) error {
	if qty <= 0 {
		return domainerrors.ErrInvalidQuantity
	}
	if s.Reserved < qty {
		return domainerrors.ErrReservationMismatch
	}
	s.Reserved -= qty
	s.Quantity -= qty
	s.touch()
	return nil
}

func (s *StockItem) Restock(qty int, unitCost valueobject.Money) error {
	if qty <= 0 {
		return domainerrors.ErrInvalidQuantity
	}
	if unitCost.Currency != s.UnitCost.Currency {
		return domainerrors.ErrCurrencyMismatch
	}
	s.Quantity += qty
	s.UnitCost = unitCost
	now := time.Now()
	s.LastRestocked = &now
	s.touch()
	return nil
}

func (s *StockItem) SetReorderPoint(point int) error {
	if point < 0 {
		return domainerrors.ErrInvalidQuantity
	}
	s.ReorderPoint = point
	s.touch()
	return nil
}

// ---------- Query ----------

func (s *StockItem) Available() int { return s.Quantity - s.Reserved }
func (s *StockItem) NeedsReorder() bool { return s.Available() <= s.ReorderPoint }
func (s *StockItem) IsOutOfStock() bool { return s.Available() <= 0 }

func (s *StockItem) touch() { s.UpdatedAt = time.Now() }
```

## E.14 `domain/repository/order_repository.go`

```go
package repository

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/erp/domain/entity"
	valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type OrderRepository interface {
	Save(ctx context.Context, o *entity.Order) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Order, error)
	FindByNumber(ctx context.Context, tenantID uuid.UUID, number string) (*entity.Order, error)
	FindByCustomer(ctx context.Context, customerID uuid.UUID) ([]entity.Order, error)
	FindByStatus(ctx context.Context, tenantID uuid.UUID, status valueobject.OrderStatus) ([]entity.Order, error)
	FindByTenant(ctx context.Context, tenantID uuid.UUID) ([]entity.Order, error)
	NextOrderNumber(ctx context.Context, tenantID uuid.UUID, otype valueobject.OrderType) (string, error)
	ExistsByNumber(ctx context.Context, tenantID uuid.UUID, number string) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
```

## E.15 `domain/repository/invoice_repository.go`

```go
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/erp/domain/entity"
	valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type InvoiceRepository interface {
	Save(ctx context.Context, i *entity.Invoice) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Invoice, error)
	FindByNumber(ctx context.Context, tenantID uuid.UUID, number string) (*entity.Invoice, error)
	FindByCustomer(ctx context.Context, customerID uuid.UUID) ([]entity.Invoice, error)
	FindByOrder(ctx context.Context, orderID uuid.UUID) ([]entity.Invoice, error)
	FindByStatus(ctx context.Context, tenantID uuid.UUID, status valueobject.InvoiceStatus) ([]entity.Invoice, error)
	FindOverdue(ctx context.Context, before time.Time) ([]entity.Invoice, error)
	NextInvoiceNumber(ctx context.Context, tenantID uuid.UUID) (string, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
```

## E.16 `domain/repository/warehouse_repository.go`

```go
package repository

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/erp/domain/entity"
)

type WarehouseRepository interface {
	Save(ctx context.Context, w *entity.Warehouse) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Warehouse, error)
	FindByCode(ctx context.Context, tenantID uuid.UUID, code string) (*entity.Warehouse, error)
	FindByTenant(ctx context.Context, tenantID uuid.UUID) ([]entity.Warehouse, error)
	ExistsByCode(ctx context.Context, tenantID uuid.UUID, code string) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
```

## E.17 `domain/repository/stock_repository.go`

```go
package repository

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/erp/domain/entity"
)

type StockRepository interface {
	Save(ctx context.Context, s *entity.StockItem) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.StockItem, error)
	FindByProductWarehouse(ctx context.Context, productID, warehouseID uuid.UUID) (*entity.StockItem, error)
	FindByProduct(ctx context.Context, productID uuid.UUID) ([]entity.StockItem, error)
	FindByWarehouse(ctx context.Context, warehouseID uuid.UUID) ([]entity.StockItem, error)
	FindNeedsReorder(ctx context.Context, tenantID uuid.UUID) ([]entity.StockItem, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
```

## E.18 `domain/event/events.go`

```go
package event

import (
	"time"

	"github.com/google/uuid"
)

type OrderCreated struct {
	OrderID     uuid.UUID
	TenantID    uuid.UUID
	CustomerID  uuid.UUID
	OrderNumber string
	Type        string
	Total       float64
	Currency    string
	ItemCount   int
	OccurredAt  time.Time
}

type OrderConfirmed struct {
	OrderID    uuid.UUID
	TenantID   uuid.UUID
	OccurredAt time.Time
}

type OrderShipped struct {
	OrderID    uuid.UUID
	TenantID   uuid.UUID
	OccurredAt time.Time
}

type OrderDelivered struct {
	OrderID    uuid.UUID
	TenantID   uuid.UUID
	OccurredAt time.Time
}

type OrderCancelled struct {
	OrderID    uuid.UUID
	TenantID   uuid.UUID
	Reason     string
	OccurredAt time.Time
}

type InvoiceIssued struct {
	InvoiceID  uuid.UUID
	TenantID   uuid.UUID
	CustomerID uuid.UUID
	Total      float64
	Currency   string
	DueDate    time.Time
	OccurredAt time.Time
}

type InvoicePaid struct {
	InvoiceID  uuid.UUID
	TenantID   uuid.UUID
	Amount     float64
	Currency   string
	OccurredAt time.Time
}

type StockLow struct {
	TenantID   uuid.UUID
	ProductID  uuid.UUID
	SKU        string
	WarehouseID uuid.UUID
	Available  int
	ReorderPoint int
	OccurredAt time.Time
}

type StockReordered struct {
	TenantID   uuid.UUID
	ProductID  uuid.UUID
	Quantity   int
	OccurredAt time.Time
}
```

## E.19 `domain/value_object/address.go`

```go
package valueobject

import "regexp"

var postcodePattern = regexp.MustCompile(`^\d{5}$`)

type Address struct {
	Line1    string `json:"line1"`
	Line2    string `json:"line2,omitempty"`
	District string `json:"district"`
	Amphoe   string `json:"amphoe"`
	Province string `json:"province"`
	Postcode string `json:"postcode"`
	Country  string `json:"country"`
}

func (a Address) Validate() error {
	if len(a.Line1) > 255 {
		return errTooLong
	}
	if a.Postcode != "" && !postcodePattern.MatchString(a.Postcode) {
		return errInvalidPostcode
	}
	if a.Country == "" {
		return errCountryRequired
	}
	return nil
}

var (
	errTooLong          = errorString("address line too long")
	errInvalidPostcode  = errorString("invalid postcode")
	errCountryRequired  = errorString("country required")
)

type errorString string

func (e errorString) Error() string { return string(e) }
```

## E.20 Status After Step 1

```
✅ Step 1 Domain complete (19 files)
   - 5 entities (Order, OrderItem, Invoice, Warehouse, StockItem)
   - 7 value objects
   - 4 repository interfaces
   - 1 errors file
   - 1 events file

Reply: "OK domain. Say continue for application."

Build check:
go build ./internal/modules/erp/domain/...
Expected: PASS (no external deps)
```

------------------------------
------------------------------
------------------------------