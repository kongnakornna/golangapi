# 🏭 PART 4 — MODULE: `erp` (Enterprise Resource Planning)

> **ขนาด**: ใหญ่พิเศษ — แยก 4 ตอนย่อย
> **Part 4A**: Domain Layer (Entities + VOs + Services + Events + Errors)
> **Part 4B**: Application Layer (Use Cases + DTO + Mappers)
> **Part 4C**: Infrastructure (Postgres + Kafka + Scheduler + Cache)
> **Part 4D**: Interface + Wiring + Migration + Tests

> **Pattern เฉพาะของ module นี้**:
> 1. **Document Numbering** — Running numbers per tenant (SO-2026-000001, INV-2026-000001)
> 2. **Double-entry Accounting** — Journal entries + Chart of Accounts
> 3. **Stock Movements Ledger** — Event-sourced inventory (ไม่ mutate qty ตรงๆ)
> 4. **Multi-warehouse** — Transfer, allocation, reservation
> 5. **Sales Order Lifecycle** — DRAFT → CONFIRMED → PARTIALLY_SHIPPED → SHIPPED → INVOICED → PAID → CLOSED
> 6. **Purchase Order Lifecycle** — DRAFT → APPROVED → SENT → PARTIALLY_RECEIVED → RECEIVED → CLOSED
> 7. **Tax Engine** — VAT 7% (TH), WHT, Reverse Charge (B2B cross-border)
> 8. **Price Lists** — Customer-specific / currency / date-effective
> 9. **Payment Allocation** — Payment ↔ Invoice matching + partial payment
> 10. **Credit / Debit Notes** — Returns & adjustments
> 11. **Reorder Point Automation** — Auto-generate PO when stock low
> 12. **Multi-currency** — FX rate snapshot at transaction time

---

## 🅰️ PART 4A — DOMAIN LAYER

### A.1 โครงสร้าง Domain

```
internal/modules/erp/domain/
├── entity/
│   ├── product.go                    # Aggregate Root #1
│   ├── price_list.go                 # Aggregate Root #2
│   ├── warehouse.go                  # Aggregate Root #3
│   ├── inventory_item.go             # Aggregate Root #4 (stock at warehouse)
│   ├── stock_movement.go             # Entity (ledger)
│   ├── stock_reservation.go          # Entity
│   ├── order.go                      # Aggregate Root #5 (SO + PO)
│   ├── order_line.go                 # Entity ย่อย
│   ├── invoice.go                    # Aggregate Root #6 (AR + AP)
│   ├── invoice_line.go               # Entity ย่อย
│   ├── payment.go                    # Aggregate Root #7
│   ├── payment_allocation.go         # Entity ย่อย
│   ├── credit_note.go                # Aggregate Root #8
│   ├── supplier.go                   # Aggregate Root #9
│   ├── journal_entry.go              # Aggregate Root #10
│   └── tax_rate.go                   # Entity
├── value_object/
│   ├── sku.go
│   ├── document_number.go
│   ├── order_status.go
│   ├── order_type.go
│   ├── invoice_status.go
│   ├── invoice_type.go
│   ├── payment_method.go
│   ├── payment_status.go
│   ├── movement_type.go
│   ├── stock_status.go
│   ├── money.go                      # Re-export จาก pkg
│   ├── quantity.go
│   ├── tax_type.go
│   ├── address.go
│   ├── currency.go
│   └── account_code.go
├── repository/
│   ├── product_repository.go
│   ├── price_list_repository.go
│   ├── warehouse_repository.go
│   ├── inventory_repository.go
│   ├── order_repository.go
│   ├── invoice_repository.go
│   ├── payment_repository.go
│   ├── credit_note_repository.go
│   ├── supplier_repository.go
│   ├── journal_repository.go
│   └── audit_repository.go
├── service/
│   ├── document_number_service.go
│   ├── tax_calculator.go
│   ├── pricing_service.go
│   ├── inventory_service.go
│   ├── allocation_service.go
│   ├── journal_service.go
│   ├── reorder_service.go
│   └── port/
│       ├── payment_port.go
│       ├── customer_port.go
│       ├── notifier_port.go
│       ├── pdf_renderer_port.go
│       └── code_generator_port.go
├── event/
│   ├── product_events.go
│   ├── order_events.go
│   ├── invoice_events.go
│   ├── inventory_events.go
│   └── payment_events.go
└── errors/
    └── errors.go
```

### A.2 Value Objects

#### `domain/value_object/sku.go`
```go
package valueobject

import (
    "regexp"
    "strings"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
)

// SKU – Stock Keeping Unit: "SENSOR-TEMP-001"
type SKU string

var skuPattern = regexp.MustCompile(`^[A-Z0-9][A-Z0-9\-]{2,49}$`)

func NewSKU(s string) (SKU, error) {
    s = strings.ToUpper(strings.TrimSpace(s))
    if !skuPattern.MatchString(s) {
        return "", domainerrors.ErrInvalidSKU
    }
    return SKU(s), nil
}

func (s SKU) String() string { return string(s) }
```

#### `domain/value_object/document_number.go`
```go
package valueobject

import (
    "fmt"
    "regexp"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
)

// DocType – ประเภทเอกสาร
type DocType string

const (
    DocTypeSalesOrder    DocType = "SO"
    DocTypePurchaseOrder DocType = "PO"
    DocTypeInvoice       DocType = "INV"
    DocTypeARInvoice     DocType = "INV"  // ลูกค้า
    DocTypeAPInvoice     DocType = "BILL" // supplier
    DocTypePayment       DocType = "PAY"
    DocTypeReceipt       DocType = "RCP"
    DocTypeCreditNote    DocType = "CN"
    DocTypeDebitNote     DocType = "DN"
    DocTypeStockTransfer DocType = "TRF"
    DocTypeStockAdjust   DocType = "ADJ"
    DocTypeJournalEntry  DocType = "JE"
)

func (t DocType) IsValid() bool {
    switch t {
    case DocTypeSalesOrder, DocTypePurchaseOrder, DocTypeInvoice,
        DocTypeAPInvoice, DocTypePayment, DocTypeReceipt,
        DocTypeCreditNote, DocTypeDebitNote, DocTypeStockTransfer,
        DocTypeStockAdjust, DocTypeJournalEntry:
        return true
    }
    return false
}

// DocumentNumber – รูปแบบ {TYPE}-{YYYY}-{NNNNNN}
type DocumentNumber string

var docNoPattern = regexp.MustCompile(`^[A-Z]{2,4}-\d{4}-\d{6}$`)

func NewDocumentNumber(docType DocType, year, seq int) (DocumentNumber, error) {
    if !docType.IsValid() { return "", domainerrors.ErrInvalidDocType }
    if year < 2000 || year > 9999 || seq < 1 {
        return "", domainerrors.ErrInvalidDocNumber
    }
    return DocumentNumber(fmt.Sprintf("%s-%04d-%06d", docType, year, seq)), nil
}

func ParseDocumentNumber(s string) (DocumentNumber, error) {
    if !docNoPattern.MatchString(s) {
        return "", domainerrors.ErrInvalidDocNumber
    }
    return DocumentNumber(s), nil
}

func (d DocumentNumber) String() string { return string(d) }
```

#### `domain/value_object/order_status.go`
```go
package valueobject

type OrderStatus string

const (
    OrderStatusDraft            OrderStatus = "DRAFT"
    OrderStatusPendingApproval  OrderStatus = "PENDING_APPROVAL"
    OrderStatusConfirmed        OrderStatus = "CONFIRMED"
    OrderStatusPartiallyShipped OrderStatus = "PARTIALLY_SHIPPED"
    OrderStatusShipped          OrderStatus = "SHIPPED"
    OrderStatusPartiallyReceived OrderStatus = "PARTIALLY_RECEIVED"
    OrderStatusReceived         OrderStatus = "RECEIVED"
    OrderStatusInvoiced         OrderStatus = "INVOICED"
    OrderStatusPaid             OrderStatus = "PAID"
    OrderStatusClosed           OrderStatus = "CLOSED"
    OrderStatusCancelled        OrderStatus = "CANCELLED"
    OrderStatusOnHold           OrderStatus = "ON_HOLD"
)

func (s OrderStatus) IsValid() bool {
    switch s {
    case OrderStatusDraft, OrderStatusPendingApproval, OrderStatusConfirmed,
        OrderStatusPartiallyShipped, OrderStatusShipped,
        OrderStatusPartiallyReceived, OrderStatusReceived,
        OrderStatusInvoiced, OrderStatusPaid, OrderStatusClosed,
        OrderStatusCancelled, OrderStatusOnHold:
        return true
    }
    return false
}

// CanTransitionTo – state machine
func (s OrderStatus) CanTransitionTo(next OrderStatus) bool {
    t := map[OrderStatus][]OrderStatus{
        OrderStatusDraft:             {OrderStatusPendingApproval, OrderStatusConfirmed, OrderStatusCancelled},
        OrderStatusPendingApproval:   {OrderStatusConfirmed, OrderStatusCancelled, OrderStatusOnHold},
        OrderStatusConfirmed:         {OrderStatusPartiallyShipped, OrderStatusShipped, OrderStatusPartiallyReceived, OrderStatusReceived, OrderStatusCancelled, OrderStatusOnHold},
        OrderStatusPartiallyShipped:  {OrderStatusShipped, OrderStatusCancelled},
        OrderStatusShipped:           {OrderStatusInvoiced, OrderStatusClosed},
        OrderStatusPartiallyReceived: {OrderStatusReceived, OrderStatusCancelled},
        OrderStatusReceived:          {OrderStatusInvoiced, OrderStatusClosed},
        OrderStatusInvoiced:          {OrderStatusPaid, OrderStatusClosed},
        OrderStatusPaid:              {OrderStatusClosed},
        OrderStatusOnHold:            {OrderStatusConfirmed, OrderStatusCancelled},
        OrderStatusClosed:            {}, // terminal
        OrderStatusCancelled:         {}, // terminal
    }
    for _, allowed := range t[s] {
        if allowed == next { return true }
    }
    return false
}

func (s OrderStatus) IsTerminal() bool {
    return s == OrderStatusClosed || s == OrderStatusCancelled
}

func (s OrderStatus) IsEditable() bool {
    return s == OrderStatusDraft
}

func (s OrderStatus) CanShip() bool {
    return s == OrderStatusConfirmed || s == OrderStatusPartiallyShipped
}

func (s OrderStatus) CanReceive() bool {
    return s == OrderStatusConfirmed || s == OrderStatusPartiallyReceived
}

func (s OrderStatus) String() string { return string(s) }
```

#### `domain/value_object/order_type.go`
```go
package valueobject

type OrderType string

const (
    OrderTypeSales    OrderType = "SALES"
    OrderTypePurchase OrderType = "PURCHASE"
)

func (t OrderType) IsValid() bool {
    return t == OrderTypeSales || t == OrderTypePurchase
}

func (t OrderType) IsSales() bool    { return t == OrderTypeSales }
func (t OrderType) IsPurchase() bool { return t == OrderTypePurchase }

// DocTypeFor – แปลงเป็น DocType
func (t OrderType) DocTypeFor() DocType {
    if t == OrderTypeSales { return DocTypeSalesOrder }
    return DocTypePurchaseOrder
}
```

#### `domain/value_object/invoice_status.go`
```go
package valueobject

type InvoiceStatus string

const (
    InvoiceStatusDraft    InvoiceStatus = "DRAFT"
    InvoiceStatusIssued   InvoiceStatus = "ISSUED"
    InvoiceStatusPartial  InvoiceStatus = "PARTIALLY_PAID"
    InvoiceStatusPaid     InvoiceStatus = "PAID"
    InvoiceStatusOverdue  InvoiceStatus = "OVERDUE"
    InvoiceStatusVoid     InvoiceStatus = "VOID"
    InvoiceStatusWritten  InvoiceStatus = "WRITTEN_OFF"
)

func (s InvoiceStatus) IsValid() bool {
    switch s {
    case InvoiceStatusDraft, InvoiceStatusIssued, InvoiceStatusPartial,
        InvoiceStatusPaid, InvoiceStatusOverdue, InvoiceStatusVoid,
        InvoiceStatusWritten:
        return true
    }
    return false
}

func (s InvoiceStatus) CanTransitionTo(next InvoiceStatus) bool {
    t := map[InvoiceStatus][]InvoiceStatus{
        InvoiceStatusDraft:   {InvoiceStatusIssued, InvoiceStatusVoid},
        InvoiceStatusIssued:  {InvoiceStatusPartial, InvoiceStatusPaid, InvoiceStatusOverdue, InvoiceStatusVoid},
        InvoiceStatusPartial: {InvoiceStatusPaid, InvoiceStatusOverdue, InvoiceStatusVoid},
        InvoiceStatusOverdue: {InvoiceStatusPartial, InvoiceStatusPaid, InvoiceStatusWritten, InvoiceStatusVoid},
        InvoiceStatusPaid:    {}, // terminal
        InvoiceStatusVoid:    {}, // terminal
        InvoiceStatusWritten: {}, // terminal
    }
    for _, allowed := range t[s] {
        if allowed == next { return true }
    }
    return false
}

func (s InvoiceStatus) IsTerminal() bool {
    return s == InvoiceStatusPaid || s == InvoiceStatusVoid || s == InvoiceStatusWritten
}

func (s InvoiceStatus) IsPayable() bool {
    return s == InvoiceStatusIssued || s == InvoiceStatusPartial || s == InvoiceStatusOverdue
}

func (s InvoiceStatus) String() string { return string(s) }
```

#### `domain/value_object/invoice_type.go`
```go
package valueobject

type InvoiceType string

const (
    InvoiceTypeAR InvoiceType = "AR" // Accounts Receivable (จากลูกค้า)
    InvoiceTypeAP InvoiceType = "AP" // Accounts Payable (ถึง supplier)
)

func (t InvoiceType) IsValid() bool {
    return t == InvoiceTypeAR || t == InvoiceTypeAP
}

func (t InvoiceType) DocTypeFor() DocType {
    if t == InvoiceTypeAR { return DocTypeInvoice }
    return DocTypeAPInvoice
}
```

#### `domain/value_object/movement_type.go`
```go
package valueobject

type MovementType string

const (
    // Inbound
    MovementTypeReceipt       MovementType = "RECEIPT"        // รับจาก PO
    MovementTypeReturnIn      MovementType = "RETURN_IN"      // ลูกค้าคืน
    MovementTypeTransferIn    MovementType = "TRANSFER_IN"    // ย้ายเข้าคลัง
    MovementTypeAdjustIn      MovementType = "ADJUST_IN"      // ปรับเพิ่ม
    MovementTypeProductionIn  MovementType = "PRODUCTION_IN"  // ผลิตเสร็จ

    // Outbound
    MovementTypeIssue         MovementType = "ISSUE"          // ตัดขาย
    MovementTypeReturnOut     MovementType = "RETURN_OUT"     // คืน supplier
    MovementTypeTransferOut   MovementType = "TRANSFER_OUT"   // ย้ายออก
    MovementTypeAdjustOut     MovementType = "ADJUST_OUT"     // ปรับลด
    MovementTypeScrap         MovementType = "SCRAP"          // ของเสีย

    // Reservation (ไม่กระทบ qty_on_hand)
    MovementTypeReserve       MovementType = "RESERVE"
    MovementTypeRelease       MovementType = "RELEASE"
)

func (m MovementType) IsValid() bool {
    switch m {
    case MovementTypeReceipt, MovementTypeReturnIn, MovementTypeTransferIn,
        MovementTypeAdjustIn, MovementTypeProductionIn,
        MovementTypeIssue, MovementTypeReturnOut, MovementTypeTransferOut,
        MovementTypeAdjustOut, MovementTypeScrap,
        MovementTypeReserve, MovementTypeRelease:
        return true
    }
    return false
}

// Direction – +1 (เข้า) / -1 (ออก) / 0 (reservation)
func (m MovementType) Direction() int {
    switch m {
    case MovementTypeReceipt, MovementTypeReturnIn,
        MovementTypeTransferIn, MovementTypeAdjustIn, MovementTypeProductionIn:
        return 1
    case MovementTypeIssue, MovementTypeReturnOut,
        MovementTypeTransferOut, MovementTypeAdjustOut, MovementTypeScrap:
        return -1
    }
    return 0
}

// AffectsStock – กระทบ qty_on_hand หรือไม่
func (m MovementType) AffectsStock() bool {
    return m.Direction() != 0
}

// AffectsReserved – กระทบ qty_reserved
func (m MovementType) AffectsReserved() bool {
    return m == MovementTypeReserve || m == MovementTypeRelease
}

func (m MovementType) String() string { return string(m) }
```

#### `domain/value_object/quantity.go`
```go
package valueobject

import (
    "fmt"
    "math"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
)

// Quantity – ปริมาณ (4 ตำแหน่ง)
type Quantity struct {
    Amount float64
    UOM    string // Unit of Measure
}

func NewQuantity(amount float64, uom string) (Quantity, error) {
    if amount < 0 {
        return Quantity{}, domainerrors.ErrNegativeQuantity
    }
    if uom == "" {
        return Quantity{}, domainerrors.ErrUOMRequired
    }
    return Quantity{Amount: round4(amount), UOM: uom}, nil
}

func (q Quantity) Add(other Quantity) (Quantity, error) {
    if q.UOM != other.UOM {
        return Quantity{}, domainerrors.ErrUOMMismatch
    }
    return Quantity{Amount: round4(q.Amount + other.Amount), UOM: q.UOM}, nil
}

func (q Quantity) Sub(other Quantity) (Quantity, error) {
    if q.UOM != other.UOM {
        return Quantity{}, domainerrors.ErrUOMMismatch
    }
    return Quantity{Amount: round4(q.Amount - other.Amount), UOM: q.UOM}, nil
}

func (q Quantity) Mul(factor float64) Quantity {
    return Quantity{Amount: round4(q.Amount * factor), UOM: q.UOM}
}

func (q Quantity) IsZero() bool     { return q.Amount == 0 }
func (q Quantity) IsPositive() bool { return q.Amount > 0 }

func (q Quantity) String() string {
    return fmt.Sprintf("%.4f %s", q.Amount, q.UOM)
}

func round4(v float64) float64 {
    return math.Round(v*10000) / 10000
}
```

#### `domain/value_object/tax_type.go`
```go
package valueobject

type TaxType string

const (
    TaxTypeVAT       TaxType = "VAT"        // ภาษีมูลค่าเพิ่ม 7% (TH)
    TaxTypeZeroRated TaxType = "ZERO_RATED" // 0% (ส่งออก)
    TaxTypeExempt    TaxType = "EXEMPT"     // ยกเว้น
    TaxTypeWHT       TaxType = "WHT"        // หัก ณ ที่จ่าย
    TaxTypeReverse   TaxType = "REVERSE_CHARGE" // B2B cross-border
)

func (t TaxType) IsValid() bool {
    switch t {
    case TaxTypeVAT, TaxTypeZeroRated, TaxTypeExempt, TaxTypeWHT, TaxTypeReverse:
        return true
    }
    return false
}

// DefaultRate – อัตราภาษีเริ่มต้น (Thailand)
func (t TaxType) DefaultRate() float64 {
    switch t {
    case TaxTypeVAT:       return 0.07
    case TaxTypeZeroRated: return 0.0
    case TaxTypeExempt:    return 0.0
    case TaxTypeWHT:       return 0.03 // default 3%
    case TaxTypeReverse:   return 0.0
    }
    return 0
}
```

#### `domain/value_object/payment_method.go`
```go
package valueobject

type PaymentMethod string

const (
    PaymentMethodCash         PaymentMethod = "CASH"
    PaymentMethodBankTransfer PaymentMethod = "BANK_TRANSFER"
    PaymentMethodCreditCard   PaymentMethod = "CREDIT_CARD"
    PaymentMethodCheque       PaymentMethod = "CHEQUE"
    PaymentMethodPromptPay    PaymentMethod = "PROMPTPAY"
    PaymentMethodStripe       PaymentMethod = "STRIPE"
    PaymentMethodOmise        PaymentMethod = "OMISE"
    PaymentMethodCredit       PaymentMethod = "CREDIT" // ใช้ credit limit
)

func (m PaymentMethod) IsValid() bool {
    switch m {
    case PaymentMethodCash, PaymentMethodBankTransfer, PaymentMethodCreditCard,
        PaymentMethodCheque, PaymentMethodPromptPay, PaymentMethodStripe,
        PaymentMethodOmise, PaymentMethodCredit:
        return true
    }
    return false
}

func (m PaymentMethod) IsOnline() bool {
    return m == PaymentMethodStripe || m == PaymentMethodOmise
}
```

#### `domain/value_object/payment_status.go`
```go
package valueobject

type PaymentStatus string

const (
    PaymentStatusPending   PaymentStatus = "PENDING"
    PaymentStatusConfirmed PaymentStatus = "CONFIRMED"
    PaymentStatusFailed    PaymentStatus = "FAILED"
    PaymentStatusReversed  PaymentStatus = "REVERSED"
    PaymentStatusRefunded  PaymentStatus = "REFUNDED"
)

func (s PaymentStatus) IsValid() bool {
    switch s {
    case PaymentStatusPending, PaymentStatusConfirmed,
        PaymentStatusFailed, PaymentStatusReversed, PaymentStatusRefunded:
        return true
    }
    return false
}

func (s PaymentStatus) IsTerminal() bool {
    return s == PaymentStatusReversed || s == PaymentStatusRefunded
}

func (s PaymentStatus) String() string { return string(s) }
```

#### `domain/value_object/account_code.go`
```go
package valueobject

import (
    "regexp"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
)

// AccountCode – Chart of Accounts: "1000", "1100", "4000", "5000"
type AccountCode string

var accountPattern = regexp.MustCompile(`^[1-9][0-9]{3,5}$`)

func NewAccountCode(s string) (AccountCode, error) {
    if !accountPattern.MatchString(s) {
        return "", domainerrors.ErrInvalidAccountCode
    }
    return AccountCode(s), nil
}

// Category – หมวดบัญชี
func (c AccountCode) Category() string {
    if len(c) == 0 { return "" }
    switch c[0] {
    case '1': return "ASSET"
    case '2': return "LIABILITY"
    case '3': return "EQUITY"
    case '4': return "REVENUE"
    case '5': return "EXPENSE"
    }
    return "OTHER"
}

func (c AccountCode) String() string { return string(c) }
```

#### `domain/value_object/address.go`
```go
package valueobject

type Address struct {
    Line1    string `json:"line1"`
    Line2    string `json:"line2,omitempty"`
    District string `json:"district,omitempty"`
    Amphoe   string `json:"amphoe,omitempty"`
    Province string `json:"province,omitempty"`
    Postcode string `json:"postcode,omitempty"`
    Country  string `json:"country"`
    TaxID    string `json:"tax_id,omitempty"`
}

func (a Address) IsEmpty() bool {
    return a.Line1 == "" && a.Country == ""
}

func (a Address) Full() string {
    parts := []string{a.Line1, a.Line2, a.District, a.Amphoe, a.Province, a.Postcode, a.Country}
    out := ""
    for _, p := range parts {
        if p == "" { continue }
        if out != "" { out += ", " }
        out += p
    }
    return out
}
```

### A.3 Domain Entities

#### `domain/entity/product.go` ⭐ (Aggregate Root #1)
```go
package entity

import (
    "strings"
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type ProductType string

const (
    ProductTypeGoods     ProductType = "GOODS"     // สินค้าจับต้องได้
    ProductTypeService   ProductType = "SERVICE"   // บริการ
    ProductTypeBundle    ProductType = "BUNDLE"    // ชุดสินค้า
    ProductTypeRaw       ProductType = "RAW"       // วัตถุดิบ
)

// Product – Aggregate Root
//
// Invariants:
//  1. SKU unique per tenant
//  2. cost_price >= 0, sell_price >= 0
//  3. reorder_point >= 0
//  4. เป็น GOODS/RAW → ต้อง track inventory
type Product struct {
    ID          uuid.UUID
    TenantID    uuid.UUID
    SKU         valueobject.SKU
    Name        string
    Description string
    Type        ProductType
    Category    string
    Brand       string
    UOM         string  // unit of measure: PCS, KG, L, BOX

    CostPrice   valueobject.Money // ราคาทุน
    SellPrice   valueobject.Money // ราคาขายเริ่มต้น
    Currency    string

    TaxType     valueobject.TaxType
    TaxRate     float64

    TrackInventory bool
    ReorderPoint   float64
    ReorderQty     float64
    LeadTimeDays   int
    SafetyStock    float64

    // Physical
    Weight       float64 // kg
    Volume       float64 // m³
    Barcode      string

    // Media
    ImageURLs    []string

    IsActive     bool
    IsSellable   bool
    IsPurchasable bool

    Metadata     map[string]any
    Tags         []string

    CreatedBy    uuid.UUID
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

func NewProduct(
    tenantID uuid.UUID,
    sku valueobject.SKU,
    name string,
    ptype ProductType,
    uom string,
    costPrice, sellPrice valueobject.Money,
    actor uuid.UUID,
) (*Product, error) {
    if tenantID == uuid.Nil {
        return nil, domainerrors.ErrInvalidTenantID
    }
    name = strings.TrimSpace(name)
    if name == "" || len(name) > 255 {
        return nil, domainerrors.ErrInvalidProductName
    }
    if uom == "" {
        return nil, domainerrors.ErrUOMRequired
    }
    if costPrice.Currency != sellPrice.Currency {
        return nil, domainerrors.ErrCurrencyMismatch
    }

    now := time.Now()
    return &Product{
        ID: uuid.New(), TenantID: tenantID, SKU: sku,
        Name: name, Type: ptype, UOM: uom,
        CostPrice: costPrice, SellPrice: sellPrice,
        Currency: costPrice.Currency,
        TaxType:  valueobject.TaxTypeVAT,
        TaxRate:  0.07,
        TrackInventory: ptype == ProductTypeGoods || ptype == ProductTypeRaw,
        ReorderPoint:   0, ReorderQty: 0, SafetyStock: 0,
        IsActive: true, IsSellable: true, IsPurchasable: true,
        Metadata: map[string]any{}, Tags: []string{},
        CreatedBy: actor, CreatedAt: now, UpdatedAt: now,
    }, nil
}

func (p *Product) Rename(name string) error {
    name = strings.TrimSpace(name)
    if name == "" || len(name) > 255 {
        return domainerrors.ErrInvalidProductName
    }
    p.Name = name
    p.touch()
    return nil
}

func (p *Product) UpdatePrices(cost, sell valueobject.Money) error {
    if cost.Currency != p.Currency || sell.Currency != p.Currency {
        return domainerrors.ErrCurrencyMismatch
    }
    p.CostPrice = cost
    p.SellPrice = sell
    p.touch()
    return nil
}

func (p *Product) SetReorderPolicy(point, qty, safetyStock float64, leadTime int) error {
    if point < 0 || qty < 0 || safetyStock < 0 || leadTime < 0 {
        return domainerrors.ErrInvalidReorderPolicy
    }
    p.ReorderPoint = point
    p.ReorderQty = qty
    p.SafetyStock = safetyStock
    p.LeadTimeDays = leadTime
    p.touch()
    return nil
}

func (p *Product) SetTax(tt valueobject.TaxType, rate float64) error {
    if !tt.IsValid() { return domainerrors.ErrInvalidTaxType }
    if rate < 0 || rate > 1 { return domainerrors.ErrInvalidTaxRate }
    p.TaxType = tt
    p.TaxRate = rate
    p.touch()
    return nil
}

func (p *Product) Activate()   { p.IsActive = true; p.touch() }
func (p *Product) Deactivate() { p.IsActive = false; p.touch() }

// --- Query ---

func (p *Product) IsTangible() bool  { return p.Type == ProductTypeGoods || p.Type == ProductTypeRaw }
func (p *Product) NeedsReorder(stock float64) bool {
    return p.TrackInventory && p.ReorderPoint > 0 && stock <= p.ReorderPoint
}
func (p *Product) GrossMarginPercent() float64 {
    if p.SellPrice.Amount == 0 { return 0 }
    return (p.SellPrice.Amount - p.CostPrice.Amount) / p.SellPrice.Amount * 100
}

func (p *Product) touch() { p.UpdatedAt = time.Now() }
```

#### `domain/entity/warehouse.go` ⭐ (Aggregate Root #3)
```go
package entity

import (
    "strings"
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type WarehouseType string

const (
    WarehouseTypeMain     WarehouseType = "MAIN"
    WarehouseTypeBranch   WarehouseType = "BRANCH"
    WarehouseTypeTransit  WarehouseType = "TRANSIT"  // ระหว่างขนส่ง
    WarehouseTypeReturn   WarehouseType = "RETURN"   // คลังรับคืน
    WarehouseTypeVirtual  WarehouseType = "VIRTUAL"  // virtual (สำหรับ adjustment)
)

// Warehouse – Aggregate Root
type Warehouse struct {
    ID          uuid.UUID
    TenantID    uuid.UUID
    Code        string
    Name        string
    Type        WarehouseType
    Address     valueobject.Address
    ManagerID   *uuid.UUID
    Phone       string
    Email       string

    IsActive    bool
    IsDefault   bool

    Metadata    map[string]any
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func NewWarehouse(tenantID uuid.UUID, code, name string, wtype WarehouseType) (*Warehouse, error) {
    code = strings.TrimSpace(strings.ToUpper(code))
    name = strings.TrimSpace(name)
    if code == "" || name == "" {
        return nil, domainerrors.ErrInvalidWarehouse
    }
    now := time.Now()
    return &Warehouse{
        ID: uuid.New(), TenantID: tenantID,
        Code: code, Name: name, Type: wtype,
        IsActive: true, Metadata: map[string]any{},
        CreatedAt: now, UpdatedAt: now,
    }, nil
}

func (w *Warehouse) SetDefault(isDefault bool) {
    w.IsDefault = isDefault
    w.touch()
}

func (w *Warehouse) Activate()   { w.IsActive = true; w.touch() }
func (w *Warehouse) Deactivate() { w.IsActive = false; w.touch() }

func (w *Warehouse) touch() { w.UpdatedAt = time.Now() }
```

#### `domain/entity/inventory_item.go` ⭐ (Aggregate Root #4)
```go
package entity

import (
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
)

// InventoryItem – stock ของ product ที่ warehouse หนึ่ง
//
// Invariants:
//  1. (tenant_id, product_id, warehouse_id) unique
//  2. qty_on_hand >= 0 (ห้ามติดลบ)
//  3. qty_reserved >= 0 และ <= qty_on_hand
//  4. qty_available = qty_on_hand - qty_reserved
type InventoryItem struct {
    ID          uuid.UUID
    TenantID    uuid.UUID
    ProductID   uuid.UUID
    WarehouseID uuid.UUID

    QtyOnHand   float64
    QtyReserved float64

    // Cost tracking
    AvgCost     float64
    LastCost    float64

    // Bin/Location (optional)
    BinLocation string

    UpdatedAt   time.Time
    CreatedAt   time.Time
}

// NewInventoryItem – สร้างรายการ inventory ใหม่
func NewInventoryItem(tenantID, productID, warehouseID uuid.UUID) *InventoryItem {
    now := time.Now()
    return &InventoryItem{
        ID: uuid.New(), TenantID: tenantID,
        ProductID: productID, WarehouseID: warehouseID,
        UpdatedAt: now, CreatedAt: now,
    }
}

// QtyAvailable – ปริมาณที่ใช้ได้
func (i *InventoryItem) QtyAvailable() float64 {
    v := i.QtyOnHand - i.QtyReserved
    if v < 0 { return 0 }
    return v
}

// ApplyMovement – ใช้ movement บันทึกการเปลี่ยนแปลง
//   direction: +1 (in) / -1 (out)
func (i *InventoryItem) ApplyMovement(direction int, qty float64, unitCost float64) error {
    if qty < 0 {
        return domainerrors.ErrNegativeQuantity
    }
    switch direction {
    case 1:
        i.addStock(qty, unitCost)
    case -1:
        if err := i.removeStock(qty); err != nil {
            return err
        }
    default:
        return domainerrors.ErrInvalidDirection
    }
    i.UpdatedAt = time.Now()
    return nil
}

func (i *InventoryItem) addStock(qty, unitCost float64) {
    // weighted average cost
    if qty > 0 {
        totalValue := i.QtyOnHand*i.AvgCost + qty*unitCost
        totalQty := i.QtyOnHand + qty
        if totalQty > 0 {
            i.AvgCost = round4(totalValue / totalQty)
        }
    }
    i.QtyOnHand += qty
    if unitCost > 0 {
        i.LastCost = unitCost
    }
}

func (i *InventoryItem) removeStock(qty float64) error {
    if qty > i.QtyOnHand {
        return domainerrors.ErrInsufficientStock
    }
    i.QtyOnHand -= qty
    return nil
}

// Reserve – จองสินค้า (เพิ่ม qty_reserved)
func (i *InventoryItem) Reserve(qty float64) error {
    if qty <= 0 { return domainerrors.ErrInvalidQuantity }
    if i.QtyAvailable() < qty {
        return domainerrors.ErrInsufficientAvailable
    }
    i.QtyReserved += qty
    i.UpdatedAt = time.Now()
    return nil
}

// Release – คืนสินค้าที่จองไว้
func (i *InventoryItem) Release(qty float64) error {
    if qty <= 0 { return domainerrors.ErrInvalidQuantity }
    if i.QtyReserved < qty {
        return domainerrors.ErrInvalidReleaseAmount
    }
    i.QtyReserved -= qty
    i.UpdatedAt = time.Now()
    return nil
}

// ConsumeReservation – ใช้ reserved ตัดสต็อกจริง (ตอน ship)
func (i *InventoryItem) ConsumeReservation(qty float64) error {
    if qty <= 0 { return domainerrors.ErrInvalidQuantity }
    if i.QtyReserved < qty {
        return domainerrors.ErrInvalidReleaseAmount
    }
    if i.QtyOnHand < qty {
        return domainerrors.ErrInsufficientStock
    }
    i.QtyReserved -= qty
    i.QtyOnHand -= qty
    i.UpdatedAt = time.Now()
    return nil
}

// SetBinLocation – ตำแหน่งชั้นวาง
func (i *InventoryItem) SetBinLocation(bin string) {
    i.BinLocation = bin
    i.UpdatedAt = time.Now()
}

// Value – มูลค่าสต็อก (at cost)
func (i *InventoryItem) Value() float64 {
    return round4(i.QtyOnHand * i.AvgCost)
}

func round4(v float64) float64 {
    return float64(int(v*10000+0.5)) / 10000
}
```

#### `domain/entity/stock_movement.go`
```go
package entity

import (
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

// StockMovement – ledger entry (immutable, event-sourced)
type StockMovement struct {
    ID             uuid.UUID
    TenantID       uuid.UUID
    ProductID      uuid.UUID
    WarehouseID    uuid.UUID

    Type           valueobject.MovementType
    Direction      int     // +1, -1, 0 (snapshot)
    Quantity       float64
    UOM            string

    UnitCost       float64
    TotalCost      float64

    // Balance after movement (snapshot)
    QtyOnHandAfter    float64
    QtyReservedAfter  float64

    // Reference
    RefType    string  // "ORDER" | "INVOICE" | "TRANSFER" | "ADJUSTMENT"
    RefID      *uuid.UUID
    RefNumber  string

    Reason     string
    Notes      string

    // For transfer
    CounterWarehouseID *uuid.UUID

    OccurredAt time.Time
    CreatedBy  uuid.UUID
}

func NewStockMovement(
    tenantID, productID, warehouseID uuid.UUID,
    mvType valueobject.MovementType,
    qty float64,
    unitCost float64,
    refType string,
    refID *uuid.UUID,
    actor uuid.UUID,
) (*StockMovement, error) {
    if !mvType.IsValid() {
        return nil, domainerrors.ErrInvalidMovementType
    }
    if qty <= 0 {
        return nil, domainerrors.ErrInvalidQuantity
    }
    return &StockMovement{
        ID: uuid.New(), TenantID: tenantID,
        ProductID: productID, WarehouseID: warehouseID,
        Type: mvType, Direction: mvType.Direction(),
        Quantity: qty, UnitCost: unitCost,
        TotalCost: qty * unitCost,
        RefType: refType, RefID: refID,
        OccurredAt: time.Now(), CreatedBy: actor,
    }, nil
}

// SnapshotBalance – บันทึกยอดคงเหลือหลัง movement
func (m *StockMovement) SnapshotBalance(onHand, reserved float64) {
    m.QtyOnHandAfter = onHand
    m.QtyReservedAfter = reserved
}

func (m *StockMovement) IsInbound() bool  { return m.Direction > 0 }
func (m *StockMovement) IsOutbound() bool { return m.Direction < 0 }
func (m *StockMovement) IsReservation() bool {
    return m.Type == valueobject.MovementTypeReserve ||
        m.Type == valueobject.MovementTypeRelease
}
```

#### `domain/entity/stock_reservation.go`
```go
package entity

import (
    "time"

    "github.com/google/uuid"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type ReservationStatus string

const (
    ReservationStatusActive    ReservationStatus = "ACTIVE"
    ReservationStatusConsumed  ReservationStatus = "CONSUMED"
    ReservationStatusReleased  ReservationStatus = "RELEASED"
    ReservationStatusExpired   ReservationStatus = "EXPIRED"
)

// StockReservation – การจองสต็อก (ผูกกับ order line)
type StockReservation struct {
    ID          uuid.UUID
    TenantID    uuid.UUID
    ProductID   uuid.UUID
    WarehouseID uuid.UUID
    OrderID     uuid.UUID
    OrderLineID uuid.UUID

    Quantity    float64
    Status      ReservationStatus
    ExpiresAt   *time.Time
    ConsumedAt  *time.Time
    ReleasedAt  *time.Time

    CreatedAt   time.Time
}

func (r *StockReservation) Consume(actor uuid.UUID) error {
    if r.Status != ReservationStatusActive {
        return ErrReservationNotActive
    }
    now := time.Now()
    r.Status = ReservationStatusConsumed
    r.ConsumedAt = &now
    return nil
}

func (r *StockReservation) Release() {
    if r.Status != ReservationStatusActive { return }
    now := time.Now()
    r.Status = ReservationStatusReleased
    r.ReleasedAt = &now
}

func (r *StockReservation) Expire() {
    if r.Status != ReservationStatusActive { return }
    r.Status = ReservationStatusExpired
}

var _ = valueobject.MovementTypeReserve
```

#### `domain/entity/order_line.go`
```go
package entity

import (
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

// OrderLine – รายการใน order
type OrderLine struct {
    ID          uuid.UUID
    OrderID     uuid.UUID
    LineNo      int

    ProductID   uuid.UUID
    SKU         string
    Name        string
    Description string

    Quantity    float64
    UOM         string

    UnitPrice   float64
    Discount    float64     // ส่วนลดต่อหน่วย
    DiscountPct float64     // หรือ %
    TaxType     valueobject.TaxType
    TaxRate     float64

    // Computed
    Subtotal    float64 // (unit_price - discount) * qty
    TaxAmount   float64
    LineTotal   float64 // subtotal + tax

    // Fulfillment tracking
    QtyShipped  float64
    QtyReceived float64
    QtyInvoiced float64

    // Warehouse for this line (อาจต่างกันต่อ line)
    WarehouseID *uuid.UUID

    Notes       string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func NewOrderLine(
    orderID uuid.UUID,
    lineNo int,
    productID uuid.UUID,
    sku, name, uom string,
    qty, unitPrice float64,
    taxType valueobject.TaxType,
    taxRate float64,
) (*OrderLine, error) {
    if qty <= 0 {
        return nil, domainerrors.ErrInvalidQuantity
    }
    if unitPrice < 0 {
        return nil, domainerrors.ErrInvalidPrice
    }
    now := time.Now()
    line := &OrderLine{
        ID: uuid.New(), OrderID: orderID, LineNo: lineNo,
        ProductID: productID, SKU: sku, Name: name,
        Quantity: qty, UOM: uom,
        UnitPrice: unitPrice,
        TaxType: taxType, TaxRate: taxRate,
        CreatedAt: now, UpdatedAt: now,
    }
    line.recalculate()
    return line, nil
}

// SetDiscount – ตั้งส่วนลด (fixed amount per unit)
func (l *OrderLine) SetDiscount(amount float64) error {
    if amount < 0 { return domainerrors.ErrInvalidDiscount }
    if amount > l.UnitPrice { return domainerrors.ErrDiscountExceedsPrice }
    l.Discount = amount
    l.DiscountPct = 0
    l.recalculate()
    return nil
}

// SetDiscountPercent – ตั้งส่วนลดเป็น %
func (l *OrderLine) SetDiscountPercent(pct float64) error {
    if pct < 0 || pct > 100 { return domainerrors.ErrInvalidDiscount }
    l.DiscountPct = pct
    l.Discount = 0
    l.recalculate()
    return nil
}

// SetQuantity – เปลี่ยนจำนวน (เฉพาะ DRAFT)
func (l *OrderLine) SetQuantity(qty float64) error {
    if qty <= 0 { return domainerrors.ErrInvalidQuantity }
    l.Quantity = qty
    l.recalculate()
    return nil
}

// Ship – ตัดสต็อก
func (l *OrderLine) Ship(qty float64) error {
    if qty <= 0 { return domainerrors.ErrInvalidQuantity }
    if l.QtyShipped + qty > l.Quantity {
        return domainerrors.ErrOverShip
    }
    l.QtyShipped += qty
    l.UpdatedAt = time.Now()
    return nil
}

// Receive – รับสินค้า (PO)
func (l *OrderLine) Receive(qty float64) error {
    if qty <= 0 { return domainerrors.ErrInvalidQuantity }
    if l.QtyReceived + qty > l.Quantity {
        return domainerrors.ErrOverReceive
    }
    l.QtyReceived += qty
    l.UpdatedAt = time.Now()
    return nil
}

// Invoice – ออกใบแจ้งหนี้
func (l *OrderLine) Invoice(qty float64) error {
    if qty <= 0 { return domainerrors.ErrInvalidQuantity }
    if l.QtyInvoiced + qty > l.Quantity {
        return domainerrors.ErrOverInvoice
    }
    l.QtyInvoiced += qty
    l.UpdatedAt = time.Now()
    return nil
}

// RemainingToShip – ยังไม่ส่งเท่าไหร่
func (l *OrderLine) RemainingToShip() float64 {
    return l.Quantity - l.QtyShipped
}

func (l *OrderLine) RemainingToReceive() float64 {
    return l.Quantity - l.QtyReceived
}

func (l *OrderLine) RemainingToInvoice() float64 {
    return l.Quantity - l.QtyInvoiced
}

// Recalculate – คำนวณ subtotal, tax, total
func (l *OrderLine) Recalculate() {
    l.recalculate()
}

func (l *OrderLine) recalculate() {
    // Effective unit price after discount
    effPrice := l.UnitPrice
    if l.Discount > 0 {
        effPrice -= l.Discount
    }
    if l.DiscountPct > 0 {
        effPrice = effPrice * (1 - l.DiscountPct/100)
    }
    if effPrice < 0 { effPrice = 0 }

    l.Subtotal = round2(effPrice * l.Quantity)

    // Tax
    if l.TaxType == valueobject.TaxTypeVAT {
        l.TaxAmount = round2(l.Subtotal * l.TaxRate)
    } else {
        l.TaxAmount = 0
    }
    l.LineTotal = round2(l.Subtotal + l.TaxAmount)
}

func round2(v float64) float64 {
    return float64(int(v*100+0.5)) / 100
}
```

#### `domain/entity/order.go` ⭐ (Aggregate Root #5)
```go
package entity

import (
    "strings"
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

// Order – Aggregate Root (Sales Order + Purchase Order)
//
// Invariants:
//  1. order_no unique per tenant
//  2. ต้องมีอย่างน้อย 1 line
//  3. total = sum(lines.total) - order_discount + tax
//  4. ยอด shipped/received ≤ ordered
//  5. ยอด invoiced ≤ shipped/received
//  6. CONFIRMED ขึ้นไป → แก้ lines ไม่ได้
type Order struct {
    ID       uuid.UUID
    TenantID uuid.UUID
    OrderNo  valueobject.DocumentNumber
    Type     valueobject.OrderType

    // Sales = customer_id, Purchase = supplier_id
    CustomerID *uuid.UUID
    SupplierID *uuid.UUID

    Status valueobject.OrderStatus

    // Dates
    OrderDate   time.Time
    ExpectedAt  *time.Time
    ConfirmedAt *time.Time
    ShippedAt   *time.Time
    ReceivedAt  *time.Time
    ClosedAt    *time.Time
    CancelledAt *time.Time

    // Shipping
    ShipToAddress   valueobject.Address
    BillToAddress   valueobject.Address
    ShippingMethod  string
    TrackingNumber  string

    // Lines
    Lines []*OrderLine

    // Totals
    Subtotal       float64
    DiscountAmount float64
    TaxAmount      float64
    ShippingFee    float64
    TotalAmount    float64
    Currency       string

    // Currency / FX
    FXRate         float64 // rate ณ วันที่สร้าง
    BaseCurrency   string
    BaseTotal      float64

    // Reference
    POReference    string // ลูกค้า PO number
    QuotationID    *uuid.UUID
    ContractID     *uuid.UUID

    Notes          string
    InternalNotes  string
    Tags           []string
    Metadata       map[string]any

    CreatedBy      uuid.UUID
    ApprovedBy     *uuid.UUID
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

// NewOrder – Factory
func NewOrder(
    tenantID uuid.UUID,
    orderNo valueobject.DocumentNumber,
    otype valueobject.OrderType,
    orderDate time.Time,
    currency string,
    actor uuid.UUID,
) (*Order, error) {
    if tenantID == uuid.Nil { return nil, domainerrors.ErrInvalidTenantID }
    if !otype.IsValid() { return nil, domainerrors.ErrInvalidOrderType }

    now := time.Now()
    return &Order{
        ID: uuid.New(), TenantID: tenantID,
        OrderNo: orderNo, Type: otype,
        Status:    valueobject.OrderStatusDraft,
        OrderDate: orderDate,
        Lines:     []*OrderLine{},
        Currency:  currency,
        FXRate:    1.0, BaseCurrency: "THB",
        Tags:      []string{},
        Metadata:  map[string]any{},
        CreatedBy: actor,
        CreatedAt: now, UpdatedAt: now,
    }, nil
}

// --- Behavior: Lines ---

func (o *Order) AddLine(line *OrderLine) error {
    if !o.Status.IsEditable() {
        return domainerrors.ErrOrderNotEditable
    }
    for _, l := range o.Lines {
        if l.ProductID == line.ProductID && l.WarehouseID == line.WarehouseID {
            // merge qty
            return l.SetQuantity(l.Quantity + line.Quantity)
        }
    }
    line.OrderID = o.ID
    line.LineNo = len(o.Lines) + 1
    o.Lines = append(o.Lines, line)
    o.recalculate()
    o.touch()
    return nil
}

func (o *Order) RemoveLine(lineID uuid.UUID) error {
    if !o.Status.IsEditable() {
        return domainerrors.ErrOrderNotEditable
    }
    for i, l := range o.Lines {
        if l.ID == lineID {
            o.Lines = append(o.Lines[:i], o.Lines[i+1:]...)
            // renumber
            for j, l2 := range o.Lines {
                l2.LineNo = j + 1
            }
            o.recalculate()
            o.touch()
            return nil
        }
    }
    return domainerrors.ErrOrderLineNotFound
}

func (o *Order) SetOrderDiscount(amount float64) error {
    if !o.Status.IsEditable() {
        return domainerrors.ErrOrderNotEditable
    }
    if amount < 0 { return domainerrors.ErrInvalidDiscount }
    o.DiscountAmount = amount
    o.recalculate()
    o.touch()
    return nil
}

func (o *Order) SetShippingFee(fee float64) error {
    if fee < 0 { return domainerrors.ErrInvalidPrice }
    o.ShippingFee = fee
    o.recalculate()
    o.touch()
    return nil
}

// --- Behavior: Lifecycle ---

func (o *Order) SetCustomer(customerID uuid.UUID, addr valueobject.Address) error {
    if !o.Type.IsSales() {
        return domainerrors.ErrOrderTypeMismatch
    }
    o.CustomerID = &customerID
    o.ShipToAddress = addr
    if o.BillToAddress.IsEmpty() { o.BillToAddress = addr }
    o.touch()
    return nil
}

func (o *Order) SetSupplier(supplierID uuid.UUID, addr valueobject.Address) error {
    if !o.Type.IsPurchase() {
        return domainerrors.ErrOrderTypeMismatch
    }
    o.SupplierID = &supplierID
    o.ShipToAddress = addr
    o.touch()
    return nil
}

func (o *Order) SubmitForApproval() error {
    if o.Status != valueobject.OrderStatusDraft {
        return domainerrors.ErrInvalidStatusTransition
    }
    if len(o.Lines) == 0 {
        return domainerrors.ErrEmptyOrder
    }
    o.Status = valueobject.OrderStatusPendingApproval
    o.touch()
    return nil
}

func (o *Order) Confirm(approverID uuid.UUID) error {
    if !o.Status.CanTransitionTo(valueobject.OrderStatusConfirmed) {
        return domainerrors.ErrInvalidStatusTransition
    }
    if len(o.Lines) == 0 {
        return domainerrors.ErrEmptyOrder
    }
    if o.Type.IsSales() && o.CustomerID == nil {
        return domainerrors.ErrCustomerRequired
    }
    if o.Type.IsPurchase() && o.SupplierID == nil {
        return domainerrors.ErrSupplierRequired
    }
    now := time.Now()
    o.Status = valueobject.OrderStatusConfirmed
    o.ConfirmedAt = &now
    o.ApprovedBy = &approverID
    o.touch()
    return nil
}

// Ship – ตัดสต็อก (Sales Order)
// ส่ง lineQuantities เป็น map[lineID]qty
func (o *Order) Ship(lineQuantities map[uuid.UUID]float64) error {
    if !o.Type.IsSales() {
        return domainerrors.ErrNotSalesOrder
    }
    if !o.Status.CanShip() {
        return domainerrors.ErrInvalidStatusTransition
    }

    totalShipped := 0.0
    for _, l := range o.Lines {
        qty := lineQuantities[l.ID]
        if qty == 0 { continue }
        if err := l.Ship(qty); err != nil { return err }
        totalShipped += qty
    }
    if totalShipped == 0 {
        return domainerrors.ErrNoLinesToShip
    }

    // Determine new status
    allShipped := true
    for _, l := range o.Lines {
        if l.QtyShipped < l.Quantity {
            allShipped = false
            break
        }
    }

    now := time.Now()
    if allShipped {
        o.Status = valueobject.OrderStatusShipped
        o.ShippedAt = &now
    } else {
        o.Status = valueobject.OrderStatusPartiallyShipped
    }
    o.touch()
    return nil
}

// Receive – รับสินค้า (Purchase Order)
func (o *Order) Receive(lineQuantities map[uuid.UUID]float64) error {
    if !o.Type.IsPurchase() {
        return domainerrors.ErrNotPurchaseOrder
    }
    if !o.Status.CanReceive() {
        return domainerrors.ErrInvalidStatusTransition
    }

    totalReceived := 0.0
    for _, l := range o.Lines {
        qty := lineQuantities[l.ID]
        if qty == 0 { continue }
        if err := l.Receive(qty); err != nil { return err }
        totalReceived += qty
    }
    if totalReceived == 0 {
        return domainerrors.ErrNoLinesToReceive
    }

    allReceived := true
    for _, l := range o.Lines {
        if l.QtyReceived < l.Quantity {
            allReceived = false
            break
        }
    }

    now := time.Now()
    if allReceived {
        o.Status = valueobject.OrderStatusReceived
        o.ReceivedAt = &now
    } else {
        o.Status = valueobject.OrderStatusPartiallyReceived
    }
    o.touch()
    return nil
}

// Invoice – mark ว่า line นี้ถูกออกใบแจ้งหนี้แล้ว
func (o *Order) Invoice(lineQuantities map[uuid.UUID]float64) error {
    if o.Status != valueobject.OrderStatusShipped &&
        o.Status != valueobject.OrderStatusReceived &&
        o.Status != valueobject.OrderStatusPartiallyShipped &&
        o.Status != valueobject.OrderStatusPartiallyReceived {
        return domainerrors.ErrInvalidStatusTransition
    }
    for _, l := range o.Lines {
        qty := lineQuantities[l.ID]
        if qty == 0 { continue }
        if err := l.Invoice(qty); err != nil { return err }
    }
    o.Status = valueobject.OrderStatusInvoiced
    o.touch()
    return nil
}

func (o *Order) MarkPaid() error {
    if o.Status != valueobject.OrderStatusInvoiced {
        return domainerrors.ErrInvalidStatusTransition
    }
    o.Status = valueobject.OrderStatusPaid
    o.touch()
    return nil
}

func (o *Order) Close() error {
    if o.Status.IsTerminal() {
        return domainerrors.ErrOrderAlreadyTerminal
    }
    now := time.Now()
    o.Status = valueobject.OrderStatusClosed
    o.ClosedAt = &now
    o.touch()
    return nil
}

func (o *Order) Cancel(reason string, actor uuid.UUID) error {
    if !o.Status.CanTransitionTo(valueobject.OrderStatusCancelled) {
        return domainerrors.ErrInvalidStatusTransition
    }
    if o.Status == valueobject.OrderStatusShipped ||
        o.Status == valueobject.OrderStatusReceived {
        return domainerrors.ErrCannotCancelShippedOrder
    }
    now := time.Now()
    o.Status = valueobject.OrderStatusCancelled
    o.CancelledAt = &now
    o.Metadata["cancel_reason"] = reason
    o.Metadata["cancelled_by"] = actor.String()
    o.touch()
    return nil
}

func (o *Order) PutOnHold(reason string) error {
    if !o.Status.CanTransitionTo(valueobject.OrderStatusOnHold) {
        return domainerrors.ErrInvalidStatusTransition
    }
    o.Status = valueobject.OrderStatusOnHold
    o.Metadata["hold_reason"] = reason
    o.touch()
    return nil
}

func (o *Order) Resume() error {
    if o.Status != valueobject.OrderStatusOnHold {
        return domainerrors.ErrInvalidStatusTransition
    }
    o.Status = valueobject.OrderStatusConfirmed
    o.touch()
    return nil
}

// --- Query ---

func (o *Order) IsEditable() bool       { return o.Status.IsEditable() }
func (o *Order) IsTerminal() bool       { return o.Status.IsTerminal() }
func (o *Order) IsFullyShipped() bool {
    for _, l := range o.Lines {
        if l.QtyShipped < l.Quantity { return false }
    }
    return true
}
func (o *Order) IsFullyReceived() bool {
    for _, l := range o.Lines {
        if l.QtyReceived < l.Quantity { return false }
    }
    return true
}
func (o *Order) TotalItems() int { return len(o.Lines) }

// --- Private ---

func (o *Order) recalculate() {
    var subtotal, taxTotal float64
    for _, l := range o.Lines {
        subtotal += l.Subtotal
        taxTotal += l.TaxAmount
    }
    o.Subtotal = round2(subtotal)
    o.TaxAmount = round2(taxTotal)
    o.TotalAmount = round2(subtotal - o.DiscountAmount + taxTotal + o.ShippingFee)

    if o.FXRate > 0 {
        o.BaseTotal = round2(o.TotalAmount * o.FXRate)
    }
}

func (o *Order) touch() { o.UpdatedAt = time.Now() }

var _ = strings.TrimSpace
```

#### `domain/entity/invoice_line.go`
```go
package entity

import (
    "time"

    "github.com/google/uuid"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

// InvoiceLine – รายการในใบแจ้งหนี้
type InvoiceLine struct {
    ID          uuid.UUID
    InvoiceID   uuid.UUID
    LineNo      int
    OrderLineID *uuid.UUID // อ้างอิง order line (optional)

    ProductID   uuid.UUID
    SKU         string
    Name        string
    Description string

    Quantity    float64
    UOM         string
    UnitPrice   float64
    Discount    float64
    TaxType     valueobject.TaxType
    TaxRate     float64

    Subtotal    float64
    TaxAmount   float64
    LineTotal   float64

    CreatedAt   time.Time
}

func (l *InvoiceLine) Recalculate() {
    effPrice := l.UnitPrice - l.Discount
    if effPrice < 0 { effPrice = 0 }
    l.Subtotal = round2(effPrice * l.Quantity)
    if l.TaxType == valueobject.TaxTypeVAT {
        l.TaxAmount = round2(l.Subtotal * l.TaxRate)
    }
    l.LineTotal = round2(l.Subtotal + l.TaxAmount)
}
```

#### `domain/entity/invoice.go` ⭐ (Aggregate Root #6)
```go
package entity

import (
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

// Invoice – Aggregate Root (AR + AP)
//
// Invariants:
//  1. invoice_no unique per tenant
//  2. due_date >= issue_date
//  3. amount_paid ≤ total
//  4. VOID → ไม่แก้
type Invoice struct {
    ID        uuid.UUID
    TenantID  uuid.UUID
    InvoiceNo valueobject.DocumentNumber
    Type      valueobject.InvoiceType

    CustomerID *uuid.UUID
    SupplierID *uuid.UUID
    OrderID    *uuid.UUID

    Status valueobject.InvoiceStatus

    IssueDate time.Time
    DueDate   time.Time
    PaidAt    *time.Time
    VoidedAt  *time.Time
    VoidReason string

    Lines []*InvoiceLine

    // Totals
    Subtotal      float64
    DiscountAmount float64
    TaxAmount     float64
    WHTAmount     float64 // withholding tax
    TotalAmount   float64
    AmountPaid    float64
    AmountDue     float64
    Currency      string

    // FX
    FXRate       float64
    BaseCurrency string
    BaseTotal    float64

    // Payment terms
    PaymentTerms string  // "NET_30", "NET_15", "DUE_ON_RECEIPT"
    Notes        string

    BillingAddress valueobject.Address

    // References
    POReference string
    RefNumber   string

    Metadata    map[string]any
    Tags        []string

    CreatedBy   uuid.UUID
    IssuedBy    *uuid.UUID
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// NewInvoice
func NewInvoice(
    tenantID uuid.UUID,
    invoiceNo valueobject.DocumentNumber,
    itype valueobject.InvoiceType,
    issueDate, dueDate time.Time,
    currency string,
    actor uuid.UUID,
) (*Invoice, error) {
    if tenantID == uuid.Nil { return nil, domainerrors.ErrInvalidTenantID }
    if !itype.IsValid() { return nil, domainerrors.ErrInvalidInvoiceType }
    if dueDate.Before(issueDate) {
        return nil, domainerrors.ErrInvalidDueDate
    }
    now := time.Now()
    return &Invoice{
        ID: uuid.New(), TenantID: tenantID,
        InvoiceNo: invoiceNo, Type: itype,
        Status:    valueobject.InvoiceStatusDraft,
        IssueDate: issueDate, DueDate: dueDate,
        Lines:     []*InvoiceLine{},
        Currency:  currency,
        FXRate:    1.0, BaseCurrency: "THB",
        Metadata:  map[string]any{}, Tags: []string{},
        CreatedBy: actor,
        CreatedAt: now, UpdatedAt: now,
    }, nil
}

// --- Behavior ---

func (i *Invoice) AddLine(line *InvoiceLine) error {
    if i.Status != valueobject.InvoiceStatusDraft {
        return domainerrors.ErrInvoiceNotEditable
    }
    line.InvoiceID = i.ID
    line.LineNo = len(i.Lines) + 1
    line.Recalculate()
    i.Lines = append(i.Lines, line)
    i.recalculate()
    i.touch()
    return nil
}

func (i *Invoice) RemoveLine(lineID uuid.UUID) error {
    if i.Status != valueobject.InvoiceStatusDraft {
        return domainerrors.ErrInvoiceNotEditable
    }
    for idx, l := range i.Lines {
        if l.ID == lineID {
            i.Lines = append(i.Lines[:idx], i.Lines[idx+1:]...)
            for j, l2 := range i.Lines { l2.LineNo = j + 1 }
            i.recalculate()
            i.touch()
            return nil
        }
    }
    return domainerrors.ErrInvoiceLineNotFound
}

func (i *Invoice) Issue(actor uuid.UUID) error {
    if i.Status != valueobject.InvoiceStatusDraft {
        return domainerrors.ErrInvalidStatusTransition
    }
    if len(i.Lines) == 0 {
        return domainerrors.ErrEmptyInvoice
    }
    if i.Type == valueobject.InvoiceTypeAR && i.CustomerID == nil {
        return domainerrors.ErrCustomerRequired
    }
    if i.Type == valueobject.InvoiceTypeAP && i.SupplierID == nil {
        return domainerrors.ErrSupplierRequired
    }
    i.Status = valueobject.InvoiceStatusIssued
    i.IssuedBy = &actor
    i.touch()
    return nil
}

// ApplyPayment – ใช้ payment ตัด invoice
func (i *Invoice) ApplyPayment(amount float64) error {
    if !i.Status.IsPayable() {
        return domainerrors.ErrInvoiceNotPayable
    }
    if amount <= 0 {
        return domainerrors.ErrInvalidPaymentAmount
    }
    if amount > i.AmountDue + 0.01 { // tolerance 1 สตางค์
        return domainerrors.ErrOverpayment
    }
    i.AmountPaid = round2(i.AmountPaid + amount)
    i.AmountDue = round2(i.TotalAmount - i.AmountPaid)

    now := time.Now()
    if i.AmountDue <= 0.01 {
        i.Status = valueobject.InvoiceStatusPaid
        i.PaidAt = &now
    } else {
        i.Status = valueobject.InvoiceStatusPartial
    }
    i.touch()
    return nil
}

func (i *Invoice) MarkOverdue() error {
    if i.Status != valueobject.InvoiceStatusIssued &&
        i.Status != valueobject.InvoiceStatusPartial {
        return domainerrors.ErrInvalidStatusTransition
    }
    i.Status = valueobject.InvoiceStatusOverdue
    i.touch()
    return nil
}

func (i *Invoice) Void(reason string, actor uuid.UUID) error {
    if i.Status == valueobject.InvoiceStatusPaid {
        return domainerrors.ErrCannotVoidPaidInvoice
    }
    if i.Status == valueobject.InvoiceStatusVoid {
        return domainerrors.ErrInvoiceAlreadyVoid
    }
    now := time.Now()
    i.Status = valueobject.InvoiceStatusVoid
    i.VoidedAt = &now
    i.VoidReason = reason
    i.touch()
    return nil
}

func (i *Invoice) WriteOff(reason string) error {
    if i.Status != valueobject.InvoiceStatusOverdue {
        return domainerrors.ErrInvalidStatusTransition
    }
    i.Status = valueobject.InvoiceStatusWritten
    i.Metadata["writeoff_reason"] = reason
    i.Metadata["writeoff_at"] = time.Now()
    i.touch()
    return nil
}

// --- Query ---

func (i *Invoice) IsPaid() bool      { return i.Status == valueobject.InvoiceStatusPaid }
func (i *Invoice) IsOverdue() bool {
    return i.DueDate.Before(time.Now()) && i.Status.IsPayable()
}
func (i *Invoice) DaysOverdue() int {
    if !i.IsOverdue() { return 0 }
    return int(time.Since(i.DueDate).Hours() / 24)
}
func (i *Invoice) AgingBucket() string {
    d := i.DaysOverdue()
    switch {
    case d == 0:  return "CURRENT"
    case d <= 30: return "1_30"
    case d <= 60: return "31_60"
    case d <= 90: return "61_90"
    default:      return "90_PLUS"
    }
}

// --- Private ---

func (i *Invoice) recalculate() {
    var subtotal, tax, wht float64
    for _, l := range i.Lines {
        subtotal += l.Subtotal
        tax += l.TaxAmount
    }
    i.Subtotal = round2(subtotal)
    i.TaxAmount = round2(tax)
    i.TotalAmount = round2(subtotal - i.DiscountAmount + tax - wht)
    i.AmountDue = round2(i.TotalAmount - i.AmountPaid)
    if i.FXRate > 0 {
        i.BaseTotal = round2(i.TotalAmount * i.FXRate)
    }
}

func (i *Invoice) touch() { i.UpdatedAt = time.Now() }
```

#### `domain/entity/payment.go` ⭐ (Aggregate Root #7)
```go
package entity

import (
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

// Payment – Aggregate Root
//
// Invariants:
//  1. payment_no unique
//  2. total_allocated ≤ amount
//  3. payment_date <= now
type Payment struct {
    ID        uuid.UUID
    TenantID  uuid.UUID
    PaymentNo valueobject.DocumentNumber
    Direction string  // IN (จากลูกค้า) / OUT (จ่าย supplier)

    CustomerID *uuid.UUID
    SupplierID *uuid.UUID

    Method valueobject.PaymentMethod
    Status valueobject.PaymentStatus

    Amount     float64
    Allocated  float64
    Unallocated float64
    Currency   string

    // FX
    FXRate     float64
    BaseAmount float64

    PaymentDate time.Time
    ValueDate   *time.Time
    Reference   string  // bank ref, txn id
    BankAccount string

    Allocations []*PaymentAllocation

    Notes    string
    Metadata map[string]any

    CreatedBy uuid.UUID
    CreatedAt time.Time
    UpdatedAt time.Time
    ConfirmedAt *time.Time
    ReversedAt  *time.Time
    RefundedAt  *time.Time
}

func NewPayment(
    tenantID uuid.UUID,
    paymentNo valueobject.DocumentNumber,
    direction string,
    method valueobject.PaymentMethod,
    amount float64,
    currency string,
    paymentDate time.Time,
    actor uuid.UUID,
) (*Payment, error) {
    if tenantID == uuid.Nil { return nil, domainerrors.ErrInvalidTenantID }
    if amount <= 0 { return nil, domainerrors.ErrInvalidPaymentAmount }
    if direction != "IN" && direction != "OUT" {
        return nil, domainerrors.ErrInvalidDirection
    }
    if !method.IsValid() { return nil, domainerrors.ErrInvalidPaymentMethod }
    now := time.Now()
    return &Payment{
        ID: uuid.New(), TenantID: tenantID,
        PaymentNo: paymentNo, Direction: direction,
        Method: method, Status: valueobject.PaymentStatusPending,
        Amount: amount, Unallocated: amount,
        Currency: currency, FXRate: 1.0, BaseAmount: amount,
        PaymentDate: paymentDate,
        Allocations: []*PaymentAllocation{},
        Metadata:    map[string]any{},
        CreatedBy:   actor,
        CreatedAt:   now, UpdatedAt: now,
    }, nil
}

func (p *Payment) SetCustomer(customerID uuid.UUID) error {
    if p.Direction != "IN" { return domainerrors.ErrDirectionMismatch }
    p.CustomerID = &customerID
    p.touch()
    return nil
}

func (p *Payment) SetSupplier(supplierID uuid.UUID) error {
    if p.Direction != "OUT" { return domainerrors.ErrDirectionMismatch }
    p.SupplierID = &supplierID
    p.touch()
    return nil
}

// Confirm – ยืนยันการรับ/จ่ายเงิน
func (p *Payment) Confirm() error {
    if p.Status != valueobject.PaymentStatusPending {
        return domainerrors.ErrInvalidStatusTransition
    }
    now := time.Now()
    p.Status = valueobject.PaymentStatusConfirmed
    p.ConfirmedAt = &now
    p.touch()
    return nil
}

// Allocate – จับคู่ payment กับ invoice
func (p *Payment) Allocate(invoiceID uuid.UUID, amount float64) error {
    if p.Status != valueobject.PaymentStatusConfirmed {
        return domainerrors.ErrPaymentNotConfirmed
    }
    if amount <= 0 { return domainerrors.ErrInvalidPaymentAmount }
    if amount > p.Unallocated+0.01 {
        return domainerrors.ErrOverAllocation
    }

    // Check existing allocation to same invoice → update
    for _, a := range p.Allocations {
        if a.InvoiceID == invoiceID {
            a.Amount += amount
            p.Allocated += amount
            p.Unallocated -= amount
            p.touch()
            return nil
        }
    }

    alloc := &PaymentAllocation{
        ID:        uuid.New(),
        PaymentID: p.ID,
        InvoiceID: invoiceID,
        Amount:    amount,
        Currency:  p.Currency,
        CreatedAt: time.Now(),
    }
    p.Allocations = append(p.Allocations, alloc)
    p.Allocated += amount
    p.Unallocated -= amount
    p.touch()
    return nil
}

func (p *Payment) Deallocate(invoiceID uuid.UUID, amount float64) error {
    if p.Status == valueobject.PaymentStatusReversed {
        return domainerrors.ErrPaymentReversed
    }
    for i, a := range p.Allocations {
        if a.InvoiceID == invoiceID {
            if amount > a.Amount {
                return domainerrors.ErrOverDeallocation
            }
            a.Amount -= amount
            if a.Amount <= 0.01 {
                p.Allocations = append(p.Allocations[:i], p.Allocations[i+1:]...)
            }
            p.Allocated -= amount
            p.Unallocated += amount
            p.touch()
            return nil
        }
    }
    return domainerrors.ErrAllocationNotFound
}

func (p *Payment) Reverse(reason string) error {
    if p.Status == valueobject.PaymentStatusReversed {
        return domainerrors.ErrPaymentAlreadyReversed
    }
    now := time.Now()
    p.Status = valueobject.PaymentStatusReversed
    p.ReversedAt = &now
    p.Metadata["reverse_reason"] = reason
    p.touch()
    return nil
}

func (p *Payment) Refund(amount float64) error {
    if p.Status != valueobject.PaymentStatusConfirmed {
        return domainerrors.ErrPaymentNotConfirmed
    }
    if amount > p.Unallocated+0.01 {
        return domainerrors.ErrOverRefund
    }
    now := time.Now()
    p.Status = valueobject.PaymentStatusRefunded
    p.RefundedAt = &now
    p.Unallocated -= amount
    p.touch()
    return nil
}

func (p *Payment) IsFullyAllocated() bool { return p.Unallocated <= 0.01 }
func (p *Payment) touch()                 { p.UpdatedAt = time.Now() }
```

#### `domain/entity/payment_allocation.go`
```go
package entity

import (
    "time"

    "github.com/google/uuid"
)

// PaymentAllocation – การจับคู่ payment กับ invoice
type PaymentAllocation struct {
    ID         uuid.UUID
    PaymentID  uuid.UUID
    InvoiceID  uuid.UUID
    Amount     float64
    Currency   string
    FXRate     float64
    BaseAmount float64
    CreatedAt  time.Time
}
```

#### `domain/entity/credit_note.go` ⭐ (Aggregate Root #8)
```go
package entity

import (
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type CreditNoteReason string

const (
    CreditNoteReasonReturn       CreditNoteReason = "RETURN"
    CreditNoteReasonDiscount     CreditNoteReason = "DISCOUNT"
    CreditNoteReasonPriceAdj     CreditNoteReason = "PRICE_ADJUSTMENT"
    CreditNoteReasonDamage       CreditNoteReason = "DAMAGE"
    CreditNoteReasonWarranty     CreditNoteReason = "WARRANTY"
    CreditNoteReasonOvercharge   CreditNoteReason = "OVERCHARGE"
    CreditNoteReasonOther        CreditNoteReason = "OTHER"
)

// CreditNote – Aggregate Root (CN/DN)
type CreditNote struct {
    ID         uuid.UUID
    TenantID   uuid.UUID
    NoteNo     valueobject.DocumentNumber
    NoteType   string // "CREDIT" | "DEBIT"
    InvoiceID  *uuid.UUID
    CustomerID *uuid.UUID
    SupplierID *uuid.UUID

    Reason     CreditNoteReason
    ReasonDetail string

    IssueDate  time.Time
    Currency   string

    Subtotal   float64
    TaxAmount  float64
    TotalAmount float64

    // ถ้า CN → เครดิตให้ลูกค้า (ลด AR)
    // ถ้า DN → เรียกเก็บเพิ่ม
    Status     string // "DRAFT" | "ISSUED" | "APPLIED" | "VOID"

    AppliedToInvoice bool
    AppliedAt        *time.Time
    Refunded         bool
    RefundedAt       *time.Time

    Notes    string
    Metadata map[string]any

    CreatedBy uuid.UUID
    CreatedAt time.Time
    UpdatedAt time.Time
}

func NewCreditNote(
    tenantID uuid.UUID,
    noteNo valueobject.DocumentNumber,
    noteType string,
    issueDate time.Time,
    currency string,
    reason CreditNoteReason,
    actor uuid.UUID,
) (*CreditNote, error) {
    if tenantID == uuid.Nil { return nil, domainerrors.ErrInvalidTenantID }
    if noteType != "CREDIT" && noteType != "DEBIT" {
        return nil, domainerrors.ErrInvalidNoteType
    }
    now := time.Now()
    return &CreditNote{
        ID: uuid.New(), TenantID: tenantID,
        NoteNo: noteNo, NoteType: noteType,
        Reason: reason, IssueDate: issueDate,
        Currency: currency, Status: "DRAFT",
        Metadata: map[string]any{},
        CreatedBy: actor, CreatedAt: now, UpdatedAt: now,
    }, nil
}

func (c *CreditNote) SetAmounts(subtotal, tax float64) error {
    if subtotal < 0 || tax < 0 {
        return domainerrors.ErrInvalidAmount
    }
    c.Subtotal = subtotal
    c.TaxAmount = tax
    c.TotalAmount = round2(subtotal + tax)
    c.touch()
    return nil
}

func (c *CreditNote) Issue() error {
    if c.Status != "DRAFT" { return domainerrors.ErrInvalidStatusTransition }
    if c.TotalAmount <= 0 { return domainerrors.ErrInvalidAmount }
    c.Status = "ISSUED"
    c.touch()
    return nil
}

func (c *CreditNote) ApplyToInvoice() error {
    if c.Status != "ISSUED" { return domainerrors.ErrInvalidStatusTransition }
    now := time.Now()
    c.Status = "APPLIED"
    c.AppliedToInvoice = true
    c.AppliedAt = &now
    c.touch()
    return nil
}

func (c *CreditNote) Refund() error {
    if c.Status != "ISSUED" && c.Status != "APPLIED" {
        return domainerrors.ErrInvalidStatusTransition
    }
    now := time.Now()
    c.Refunded = true
    c.RefundedAt = &now
    c.touch()
    return nil
}

func (c *CreditNote) Void(reason string) error {
    if c.Status == "VOID" { return domainerrors.ErrAlreadyVoid }
    c.Status = "VOID"
    c.Metadata["void_reason"] = reason
    c.touch()
    return nil
}

func (c *CreditNote) touch() { c.UpdatedAt = time.Now() }
```

#### `domain/entity/supplier.go` ⭐ (Aggregate Root #9)
```go
package entity

import (
    "strings"
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type SupplierStatus string

const (
    SupplierStatusActive    SupplierStatus = "ACTIVE"
    SupplierStatusSuspended SupplierStatus = "SUSPENDED"
    SupplierStatusBlacklist SupplierStatus = "BLACKLIST"
)

// Supplier – Aggregate Root
type Supplier struct {
    ID       uuid.UUID
    TenantID uuid.UUID
    Code     string
    Name     string
    LegalName string
    TaxID    string
    Branch   string

    Email     string
    Phone     string
    ContactName string

    Address     valueobject.Address
    BillingAddr valueobject.Address
    ShippingAddr valueobject.Address

    Status SupplierStatus

    // Payment terms
    PaymentTerms   string  // NET_30
    CreditLimit    float64
    Currency       string

    // Performance
    OnTimeDeliveryPct  float64
    QualityScore       float64
    AvgLeadTimeDays    int

    Notes     string
    Metadata  map[string]any
    Tags      []string

    IsActive  bool
    CreatedBy uuid.UUID
    CreatedAt time.Time
    UpdatedAt time.Time
}

func NewSupplier(
    tenantID uuid.UUID,
    code, name string,
    actor uuid.UUID,
) (*Supplier, error) {
    code = strings.ToUpper(strings.TrimSpace(code))
    name = strings.TrimSpace(name)
    if code == "" || name == "" {
        return nil, domainerrors.ErrInvalidSupplier
    }
    now := time.Now()
    return &Supplier{
        ID: uuid.New(), TenantID: tenantID,
        Code: code, Name: name,
        Status: SupplierStatusActive,
        IsActive: true,
        Currency: "THB",
        PaymentTerms: "NET_30",
        Tags: []string{}, Metadata: map[string]any{},
        CreatedBy: actor, CreatedAt: now, UpdatedAt: now,
    }, nil
}

func (s *Supplier) Suspend(reason string) {
    s.Status = SupplierStatusSuspended
    s.Metadata["suspend_reason"] = reason
    s.touch()
}

func (s *Supplier) Blacklist(reason string) {
    s.Status = SupplierStatusBlacklist
    s.IsActive = false
    s.Metadata["blacklist_reason"] = reason
    s.touch()
}

func (s *Supplier) Activate() {
    s.Status = SupplierStatusActive
    s.IsActive = true
    s.touch()
}

func (s *Supplier) touch() { s.UpdatedAt = time.Now() }
```

#### `domain/entity/journal_entry.go` ⭐ (Aggregate Root #10)
```go
package entity

import (
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

// JournalLine – debit/credit line
type JournalLine struct {
    ID          uuid.UUID
    JournalID   uuid.UUID
    LineNo      int
    AccountCode valueobject.AccountCode
    AccountName string
    Debit       float64
    Credit      float64
    Description string
    RefType     string
    RefID       *uuid.UUID
}

// JournalEntry – Aggregate Root (Double-entry accounting)
//
// Invariants:
//  1. sum(debit) == sum(credit) — double-entry
//  2. ต้องมีอย่างน้อย 2 lines
//  3. POSTED → ไม่แก้
type JournalEntry struct {
    ID          uuid.UUID
    TenantID    uuid.UUID
    EntryNo     valueobject.DocumentNumber
    EntryDate   time.Time
    Description string

    Lines []*JournalLine

    TotalDebit  float64
    TotalCredit float64
    Currency    string

    // Source
    SourceType string  // "INVOICE" | "PAYMENT" | "CREDIT_NOTE" | "MANUAL"
    SourceID   *uuid.UUID

    Status string // "DRAFT" | "POSTED" | "VOID"

    PostedAt  *time.Time
    PostedBy  *uuid.UUID
    VoidedAt  *time.Time
    VoidReason string

    CreatedBy uuid.UUID
    CreatedAt time.Time
    UpdatedAt time.Time
}

func NewJournalEntry(
    tenantID uuid.UUID,
    entryNo valueobject.DocumentNumber,
    entryDate time.Time,
    description, currency string,
    actor uuid.UUID,
) (*JournalEntry, error) {
    if tenantID == uuid.Nil { return nil, domainerrors.ErrInvalidTenantID }
    if description == "" { return nil, domainerrors.ErrDescriptionRequired }
    now := time.Now()
    return &JournalEntry{
        ID: uuid.New(), TenantID: tenantID,
        EntryNo: entryNo, EntryDate: entryDate,
        Description: description, Currency: currency,
        Lines: []*JournalLine{}, Status: "DRAFT",
        CreatedBy: actor, CreatedAt: now, UpdatedAt: now,
    }, nil
}

func (j *JournalEntry) AddLine(accountCode valueobject.AccountCode, accountName string, debit, credit float64, desc string) error {
    if j.Status != "DRAFT" { return domainerrors.ErrJournalNotEditable }
    if debit < 0 || credit < 0 {
        return domainerrors.ErrInvalidAmount
    }
    if debit == 0 && credit == 0 {
        return domainerrors.ErrInvalidAmount
    }
    if debit > 0 && credit > 0 {
        return domainerrors.ErrCannotHaveBothDebitCredit
    }
    line := &JournalLine{
        ID: uuid.New(), JournalID: j.ID,
        LineNo: len(j.Lines) + 1,
        AccountCode: accountCode, AccountName: accountName,
        Debit: debit, Credit: credit, Description: desc,
    }
    j.Lines = append(j.Lines, line)
    j.recalculate()
    j.touch()
    return nil
}

func (j *JournalEntry) Post(actor uuid.UUID) error {
    if j.Status != "DRAFT" { return domainerrors.ErrInvalidStatusTransition }
    if len(j.Lines) < 2 { return domainerrors.ErrJournalNeedsTwoLines }
    if !j.IsBalanced() { return domainerrors.ErrJournalNotBalanced }
    now := time.Now()
    j.Status = "POSTED"
    j.PostedAt = &now
    j.PostedBy = &actor
    j.touch()
    return nil
}

func (j *JournalEntry) Void(reason string, actor uuid.UUID) error {
    if j.Status == "VOID" { return domainerrors.ErrAlreadyVoid }
    now := time.Now()
    j.Status = "VOID"
    j.VoidedAt = &now
    j.VoidReason = reason
    j.touch()
    return nil
}

func (j *JournalEntry) IsBalanced() bool {
    return round2(j.TotalDebit) == round2(j.TotalCredit)
}

func (j *JournalEntry) recalculate() {
    var d, c float64
    for _, l := range j.Lines {
        d += l.Debit
        c += l.Credit
    }
    j.TotalDebit = round2(d)
    j.TotalCredit = round2(c)
}

func (j *JournalEntry) touch() { j.UpdatedAt = time.Now() }
```

#### `domain/entity/price_list.go` ⭐ (Aggregate Root #2)
```go
package entity

import (
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
)

type PriceListEntry struct {
    ProductID   uuid.UUID
    SKU         string
    UnitPrice   float64
    MinQty      float64
    MaxQty      *float64
    DiscountPct float64
}

// PriceList – Aggregate Root (customer-specific pricing)
type PriceList struct {
    ID          uuid.UUID
    TenantID    uuid.UUID
    Code        string
    Name        string
    Description string

    CustomerID  *uuid.UUID // NULL = public
    Currency    string

    ValidFrom   time.Time
    ValidTo     *time.Time

    Entries     []PriceListEntry

    IsActive    bool
    Priority    int  // สูงกว่า = ใช้ก่อน

    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func NewPriceList(tenantID uuid.UUID, code, name, currency string, validFrom time.Time) (*PriceList, error) {
    if code == "" || name == "" {
        return nil, domainerrors.ErrInvalidPriceList
    }
    now := time.Now()
    return &PriceList{
        ID: uuid.New(), TenantID: tenantID,
        Code: code, Name: name, Currency: currency,
        ValidFrom: validFrom,
        Entries:   []PriceListEntry{},
        IsActive:  true, Priority: 5,
        CreatedAt: now, UpdatedAt: now,
    }, nil
}

func (p *PriceList) AddEntry(e PriceListEntry) error {
    for _, existing := range p.Entries {
        if existing.ProductID == e.ProductID {
            return domainerrors.ErrPriceEntryDuplicate
        }
    }
    p.Entries = append(p.Entries, e)
    p.touch()
    return nil
}

func (p *PriceList) FindPrice(productID uuid.UUID, qty float64) (*PriceListEntry, bool) {
    for i := range p.Entries {
        e := &p.Entries[i]
        if e.ProductID != productID { continue }
        if qty < e.MinQty { continue }
        if e.MaxQty != nil && qty > *e.MaxQty { continue }
        return e, true
    }
    return nil, false
}

func (p *PriceList) IsApplicable(at time.Time) bool {
    if !p.IsActive { return false }
    if at.Before(p.ValidFrom) { return false }
    if p.ValidTo != nil && at.After(*p.ValidTo) { return false }
    return true
}

func (p *PriceList) touch() { p.UpdatedAt = time.Now() }
```

#### `domain/entity/tax_rate.go`
```go
package entity

import (
    "time"

    "github.com/google/uuid"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

// TaxRate – master data ของอัตราภาษี
type TaxRate struct {
    ID         uuid.UUID
    TenantID   *uuid.UUID
    Code       string
    Name       string
    Type       valueobject.TaxType
    Rate       float64
    Country    string
    IsCompound bool
    IsActive   bool
    EffectiveFrom time.Time
    EffectiveTo   *time.Time
    CreatedAt  time.Time
}
```

### A.4 Repository Interfaces

#### `domain/repository/product_repository.go`
```go
package repository

import (
    "context"

    "github.com/google/uuid"
    "icmongolang/internal/modules/erp/domain/entity"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type ProductFilter struct {
    TenantID  uuid.UUID
    Type      *entity.ProductType
    Category  string
    Search    string
    IsActive  *bool
    Tags      []string
    Page      int
    PageSize  int
}

type ProductRepository interface {
    Save(ctx context.Context, p *entity.Product) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Product, error)
    FindBySKU(ctx context.Context, tenantID uuid.UUID, sku valueobject.SKU) (*entity.Product, error)
    List(ctx context.Context, f ProductFilter) ([]*entity.Product, int64, error)
    Delete(ctx context.Context, tenantID, id uuid.UUID) error
}
```

#### `domain/repository/order_repository.go`
```go
package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "icmongolang/internal/modules/erp/domain/entity"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type OrderFilter struct {
    TenantID   uuid.UUID
    Type       *valueobject.OrderType
    CustomerID *uuid.UUID
    SupplierID *uuid.UUID
    Statuses   []valueobject.OrderStatus
    From       *time.Time
    To         *time.Time
    Search     string
    Page       int
    PageSize   int
}

type OrderRepository interface {
    Save(ctx context.Context, o *entity.Order) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Order, error)
    FindByNo(ctx context.Context, tenantID uuid.UUID, no valueobject.DocumentNumber) (*entity.Order, error)
    List(ctx context.Context, f OrderFilter) ([]*entity.Order, int64, error)
    NextSequence(ctx context.Context, tenantID uuid.UUID, docType valueobject.DocType, year int) (int, error)
}
```

#### `domain/repository/inventory_repository.go`
```go
package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "icmongolang/internal/modules/erp/domain/entity"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type InventoryFilter struct {
    TenantID    uuid.UUID
    ProductID   *uuid.UUID
    WarehouseID *uuid.UUID
    LowStock    bool // ต่ำกว่า reorder point
    Page        int
    PageSize    int
}

type InventoryRepository interface {
    // Item
    Save(ctx context.Context, i *entity.InventoryItem) error
    FindByKey(ctx context.Context, tenantID, productID, warehouseID uuid.UUID) (*entity.InventoryItem, error)
    List(ctx context.Context, f InventoryFilter) ([]*entity.InventoryItem, int64, error)
    ListByProduct(ctx context.Context, tenantID, productID uuid.UUID) ([]*entity.InventoryItem, error)

    // Movement
    SaveMovement(ctx context.Context, m *entity.StockMovement) error
    ListMovements(ctx context.Context, tenantID, productID uuid.UUID, from, to time.Time) ([]*entity.StockMovement, error)

    // Reservation
    SaveReservation(ctx context.Context, r *entity.StockReservation) error
    FindReservation(ctx context.Context, id uuid.UUID) (*entity.StockReservation, error)
    ListActiveReservations(ctx context.Context, tenantID, productID, warehouseID uuid.UUID) ([]*entity.StockReservation, error)
    ListExpiredReservations(ctx context.Context, asOf time.Time) ([]*entity.StockReservation, error)
}

var _ = valueobject.MovementTypeReceipt
```

#### `domain/repository/invoice_repository.go`
```go
package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "icmongolang/internal/modules/erp/domain/entity"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type InvoiceFilter struct {
    TenantID   uuid.UUID
    Type       *valueobject.InvoiceType
    CustomerID *uuid.UUID
    SupplierID *uuid.UUID
    Statuses   []valueobject.InvoiceStatus
    From       *time.Time
    To         *time.Time
    Overdue    *bool
    Page       int
    PageSize   int
}

type InvoiceRepository interface {
    Save(ctx context.Context, i *entity.Invoice) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Invoice, error)
    FindByNo(ctx context.Context, tenantID uuid.UUID, no valueobject.DocumentNumber) (*entity.Invoice, error)
    List(ctx context.Context, f InvoiceFilter) ([]*entity.Invoice, int64, error)
    ListOverdue(ctx context.Context, asOf time.Time, limit int) ([]*entity.Invoice, error)
    NextSequence(ctx context.Context, tenantID uuid.UUID, docType valueobject.DocType, year int) (int, error)
}
```

#### `domain/repository/payment_repository.go`
```go
package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "icmongolang/internal/modules/erp/domain/entity"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type PaymentFilter struct {
    TenantID   uuid.UUID
    Direction  *string
    CustomerID *uuid.UUID
    SupplierID *uuid.UUID
    Statuses   []valueobject.PaymentStatus
    From       *time.Time
    To         *time.Time
    Page       int
    PageSize   int
}

type PaymentRepository interface {
    Save(ctx context.Context, p *entity.Payment) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Payment, error)
    FindByNo(ctx context.Context, tenantID uuid.UUID, no valueobject.DocumentNumber) (*entity.Payment, error)
    List(ctx context.Context, f PaymentFilter) ([]*entity.Payment, int64, error)
    NextSequence(ctx context.Context, tenantID uuid.UUID, docType valueobject.DocType, year int) (int, error)
}
```

#### `domain/repository/warehouse_repository.go`
```go
package repository

import (
    "context"

    "github.com/google/uuid"
    "icmongolang/internal/modules/erp/domain/entity"
)

type WarehouseRepository interface {
    Save(ctx context.Context, w *entity.Warehouse) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Warehouse, error)
    FindByCode(ctx context.Context, tenantID uuid.UUID, code string) (*entity.Warehouse, error)
    List(ctx context.Context, tenantID uuid.UUID) ([]*entity.Warehouse, error)
    FindDefault(ctx context.Context, tenantID uuid.UUID) (*entity.Warehouse, error)
}
```

#### `domain/repository/journal_repository.go`
```go
package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "icmongolang/internal/modules/erp/domain/entity"
)

type JournalRepository interface {
    Save(ctx context.Context, j *entity.JournalEntry) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.JournalEntry, error)
    List(ctx context.Context, tenantID uuid.UUID, from, to time.Time, page, pageSize int) ([]*entity.JournalEntry, int64, error)
    NextSequence(ctx context.Context, tenantID uuid.UUID, year int) (int, error)
}
```

#### `domain/repository/supplier_repository.go`
```go
package repository

import (
    "context"

    "github.com/google/uuid"
    "icmongolang/internal/modules/erp/domain/entity"
)

type SupplierRepository interface {
    Save(ctx context.Context, s *entity.Supplier) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Supplier, error)
    FindByCode(ctx context.Context, tenantID uuid.UUID, code string) (*entity.Supplier, error)
    List(ctx context.Context, tenantID uuid.UUID) ([]*entity.Supplier, error)
}
```

#### `domain/repository/credit_note_repository.go`
```go
package repository

import (
    "context"

    "github.com/google/uuid"
    "icmongolang/internal/modules/erp/domain/entity"
)

type CreditNoteRepository interface {
    Save(ctx context.Context, c *entity.CreditNote) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.CreditNote, error)
    ListByInvoice(ctx context.Context, tenantID, invoiceID uuid.UUID) ([]*entity.CreditNote, error)
    NextSequence(ctx context.Context, tenantID uuid.UUID, noteType string, year int) (int, error)
}
```

#### `domain/repository/price_list_repository.go`
```go
package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "icmongolang/internal/modules/erp/domain/entity"
)

type PriceListRepository interface {
    Save(ctx context.Context, p *entity.PriceList) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.PriceList, error)
    ListApplicable(ctx context.Context, tenantID, customerID uuid.UUID, at time.Time) ([]*entity.PriceList, error)
    List(ctx context.Context, tenantID uuid.UUID) ([]*entity.PriceList, error)
}
```

#### `domain/repository/audit_repository.go`
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

### A.5 Domain Services & Ports

#### `domain/service/tax_calculator.go` ⭐
```go
package service

import (
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

// TaxCalculator – คำนวณภาษีสำหรับ order/invoice
type TaxCalculator struct{}

func NewTaxCalculator() *TaxCalculator { return &TaxCalculator{} }

type TaxRequest struct {
    Subtotal  float64
    TaxType   valueobject.TaxType
    TaxRate   float64
    Inclusive bool // ราคารวมภาษีแล้วหรือยัง
}

type TaxResult struct {
    BaseAmount  float64
    TaxAmount   float64
    TotalAmount float64
    WHTAmount   float64
    NetPayable  float64
}

// Calculate – คำนวณภาษี
func (c *TaxCalculator) Calculate(req TaxRequest) TaxResult {
    res := TaxResult{}
    switch req.TaxType {
    case valueobject.TaxTypeVAT:
        if req.Inclusive {
            // ราคารวม VAT แล้ว
            res.TotalAmount = req.Subtotal
            res.BaseAmount = round2(req.Subtotal / (1 + req.TaxRate))
            res.TaxAmount = round2(res.TotalAmount - res.BaseAmount)
        } else {
            res.BaseAmount = req.Subtotal
            res.TaxAmount = round2(req.Subtotal * req.TaxRate)
            res.TotalAmount = round2(req.Subtotal + res.TaxAmount)
        }
        res.NetPayable = res.TotalAmount
    case valueobject.TaxTypeZeroRated, valueobject.TaxTypeExempt:
        res.BaseAmount = req.Subtotal
        res.TotalAmount = req.Subtotal
        res.NetPayable = req.Subtotal
    case valueobject.TaxTypeWHT:
        res.BaseAmount = req.Subtotal
        res.TotalAmount = req.Subtotal
        res.WHTAmount = round2(req.Subtotal * req.TaxRate)
        res.NetPayable = round2(req.Subtotal - res.WHTAmount)
    case valueobject.TaxTypeReverse:
        res.BaseAmount = req.Subtotal
        res.TaxAmount = round2(req.Subtotal * 0.07)
        res.TotalAmount = req.Subtotal
        res.NetPayable = req.Subtotal
    }
    return res
}

// CalculateLine – สำหรับ order/invoice line
func (c *TaxCalculator) CalculateLine(unitPrice, discount float64, qty float64, taxType valueobject.TaxType, taxRate float64) TaxResult {
    eff := unitPrice - discount
    if eff < 0 { eff = 0 }
    subtotal := round2(eff * qty)
    return c.Calculate(TaxRequest{
        Subtotal: subtotal, TaxType: taxType, TaxRate: taxRate,
    })
}

func round2(v float64) float64 {
    return float64(int(v*100+0.5)) / 100
}
```

#### `domain/service/document_number_service.go`
```go
package service

import (
    "context"
    "time"

    "github.com/google/uuid"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

// DocumentNumberService – generate running numbers
type DocumentNumberService struct {
    generator DocNumberGenerator
}

type DocNumberGenerator interface {
    NextSeq(ctx context.Context, tenantID uuid.UUID, docType valueobject.DocType, year int) (int, error)
}

func NewDocumentNumberService(g DocNumberGenerator) *DocumentNumberService {
    return &DocumentNumberService{generator: g}
}

func (s *DocumentNumberService) Next(ctx context.Context, tenantID uuid.UUID, docType valueobject.DocType) (valueobject.DocumentNumber, error) {
    year := time.Now().Year()
    seq, err := s.generator.NextSeq(ctx, tenantID, docType, year)
    if err != nil { return "", err }
    return valueobject.NewDocumentNumber(docType, year, seq)
}
```

#### `domain/service/pricing_service.go`
```go
package service

import (
    "time"

    "github.com/google/uuid"

    "icmongolang/internal/modules/erp/domain/entity"
    "icmongolang/internal/modules/erp/domain/repository"
)

// PricingService – หาราคาที่ดีที่สุดสำหรับ product
type PricingService struct {
    plRepo  repository.PriceListRepository
    prRepo  repository.ProductRepository
}

func NewPricingService(plRepo repository.PriceListRepository, prRepo repository.ProductRepository) *PricingService {
    return &PricingService{plRepo: plRepo, prRepo: prRepo}
}

type PricingResult struct {
    UnitPrice   float64
    Currency    string
    Source      string  // "price_list" | "default"
    PriceListID *uuid.UUID
    DiscountPct float64
}

// ResolvePrice – หาราคาที่เหมาะสม
func (s *PricingService) ResolvePrice(ctx contextContext, tenantID, customerID, productID uuid.UUID, qty float64) (*PricingResult, error) {
    // 1. Try customer-specific price list
    lists, err := s.plRepo.ListApplicable(ctx, tenantID, customerID, time.Now())
    if err == nil && len(lists) > 0 {
        // Highest priority first
        for _, pl := range lists {
            if entry, ok := pl.FindPrice(productID, qty); ok {
                plID := pl.ID
                return &PricingResult{
                    UnitPrice: entry.UnitPrice,
                    Currency:  pl.Currency,
                    Source:    "price_list",
                    PriceListID: &plID,
                    DiscountPct: entry.DiscountPct,
                }, nil
            }
        }
    }

    // 2. Fallback to product default
    p, err := s.prRepo.FindByID(ctx, tenantID, productID)
    if err != nil { return nil, err }
    return &PricingResult{
        UnitPrice: p.SellPrice.Amount,
        Currency:  p.SellPrice.Currency,
        Source:    "default",
    }, nil
}

type contextContext = interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key any) any
}
```

#### `domain/service/inventory_service.go` ⭐
```go
package service

import (
    "context"
    "fmt"

    "github.com/google/uuid"

    "icmongolang/internal/modules/erp/domain/entity"
    "icmongolang/internal/modules/erp/domain/repository"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

// InventoryService – orchestrate stock movements
type InventoryService struct {
    invRepo repository.InventoryRepository
}

func NewInventoryService(invRepo repository.InventoryRepository) *InventoryService {
    return &InventoryService{invRepo: invRepo}
}

// ApplyMovement – apply + record movement (atomic)
func (s *InventoryService) ApplyMovement(
    ctx context.Context,
    tenantID, productID, warehouseID uuid.UUID,
    mvType valueobject.MovementType,
    qty, unitCost float64,
    refType string, refID *uuid.UUID,
    actor uuid.UUID,
) (*entity.InventoryItem, *entity.StockMovement, error) {
    // 1. Load or create inventory item
    item, err := s.invRepo.FindByKey(ctx, tenantID, productID, warehouseID)
    if err != nil {
        item = entity.NewInventoryItem(tenantID, productID, warehouseID)
    }

    // 2. Apply movement
    if err := item.ApplyMovement(mvType.Direction(), qty, unitCost); err != nil {
        return nil, nil, err
    }

    // 3. Create movement record
    mv, err := entity.NewStockMovement(
        tenantID, productID, warehouseID,
        mvType, qty, unitCost, refType, refID, actor,
    )
    if err != nil { return nil, nil, err }

    mv.SnapshotBalance(item.QtyOnHand, item.QtyReserved)

    // 4. Save both
    if err := s.invRepo.Save(ctx, item); err != nil {
        return nil, nil, err
    }
    if err := s.invRepo.SaveMovement(ctx, mv); err != nil {
        return nil, nil, err
    }

    return item, mv, nil
}

// Reserve – จองสต็อก
func (s *InventoryService) Reserve(
    ctx context.Context,
    tenantID, productID, warehouseID uuid.UUID,
    orderID, orderLineID uuid.UUID,
    qty float64,
    actor uuid.UUID,
) error {
    item, err := s.invRepo.FindByKey(ctx, tenantID, productID, warehouseID)
    if err != nil {
        return fmt.Errorf("inventory not found: %w", err)
    }
    if err := item.Reserve(qty); err != nil {
        return err
    }

    // Record RESERVE movement
    mv, _ := entity.NewStockMovement(
        tenantID, productID, warehouseID,
        valueobject.MovementTypeReserve, qty, item.AvgCost,
        "ORDER", &orderID, actor,
    )
    mv.SnapshotBalance(item.QtyOnHand, item.QtyReserved)

    if err := s.invRepo.Save(ctx, item); err != nil { return err }
    _ = s.invRepo.SaveMovement(ctx, mv)

    // Save reservation record
    r := &entity.StockReservation{
        ID: uuid.New(), TenantID: tenantID,
        ProductID: productID, WarehouseID: warehouseID,
        OrderID: orderID, OrderLineID: orderLineID,
        Quantity: qty, Status: entity.ReservationStatusActive,
    }
    return s.invRepo.SaveReservation(ctx, r)
}

// ConsumeReservation – ตัดสต็อกจริงตอน ship
func (s *InventoryService) ConsumeReservation(
    ctx context.Context,
    reservation *entity.StockReservation,
    actor uuid.UUID,
) error {
    item, err := s.invRepo.FindByKey(ctx, reservation.TenantID, reservation.ProductID, reservation.WarehouseID)
    if err != nil { return err }

    if err := item.ConsumeReservation(reservation.Quantity); err != nil {
        return err
    }

    // Record ISSUE movement
    mv, _ := entity.NewStockMovement(
        reservation.TenantID, reservation.ProductID, reservation.WarehouseID,
        valueobject.MovementTypeIssue, reservation.Quantity, item.AvgCost,
        "ORDER", &reservation.OrderID, actor,
    )
    mv.SnapshotBalance(item.QtyOnHand, item.QtyReserved)

    if err := s.invRepo.Save(ctx, item); err != nil { return err }
    _ = s.invRepo.SaveMovement(ctx, mv)

    if err := reservation.Consume(actor); err != nil { return err }
    return s.invRepo.SaveReservation(ctx, reservation)
}

// CheckAvailability – ตรวจว่ามีสต็อกพอไหม
func (s *InventoryService) CheckAvailability(
    ctx context.Context,
    tenantID, productID, warehouseID uuid.UUID,
    qty float64,
) (bool, float64, error) {
    item, err := s.invRepo.FindByKey(ctx, tenantID, productID, warehouseID)
    if err != nil { return false, 0, nil }
    return item.QtyAvailable() >= qty, item.QtyAvailable(), nil
}
```

#### `domain/service/allocation_service.go`
```go
package service

import (
    "github.com/google/uuid"

    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

// AllocationService – จับคู่ payment กับ invoice (oldest first = FIFO)
type AllocationService struct{}

func NewAllocationService() *AllocationService { return &AllocationService{} }

type OpenInvoice struct {
    InvoiceID  uuid.UUID
    InvoiceNo  string
    DueDate    time.Time
    AmountDue  float64
    Currency   string
    Priority   int
}

type AllocationResult struct {
    InvoiceID   uuid.UUID
    InvoiceNo   string
    Amount      float64
}

// Allocate – allocate payment ให้ invoice (FIFO by due date)
func (s *AllocationService) Allocate(
    paymentAmount float64,
    invoices []OpenInvoice,
) ([]AllocationResult, float64) {
    remaining := paymentAmount
    results := []AllocationResult{}

    // Sort by due date ascending
    sort.Slice(invoices, func(i, j int) bool {
        return invoices[i].DueDate.Before(invoices[j].DueDate)
    })

    for _, inv := range invoices {
        if remaining <= 0.01 { break }
        alloc := min(remaining, inv.AmountDue)
        results = append(results, AllocationResult{
            InvoiceID: inv.InvoiceID,
            InvoiceNo: inv.InvoiceNo,
            Amount:    round2(alloc),
        })
        remaining = round2(remaining - alloc)
    }
    return results, remaining
}

func min(a, b float64) float64 {
    if a < b { return a }
    return b
}

var _ = valueobject.InvoiceStatusIssued
```

#### `domain/service/journal_service.go` ⭐
```go
package service

import (
    "context"
    "time"

    "github.com/google/uuid"

    "icmongolang/internal/modules/erp/domain/entity"
    "icmongolang/internal/modules/erp/domain/repository"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

// JournalService – สร้าง journal entry อัตโนมัติจาก business events
type JournalService struct {
    journalRepo repository.JournalRepository
    seqGen      DocNumberGenerator
}

func NewJournalService(journalRepo repository.JournalRepository, seqGen DocNumberGenerator) *JournalService {
    return &JournalService{journalRepo: journalRepo, seqGen: seqGen}
}

// PostARInvoice – บันทึก AR invoice
//   Dr. Accounts Receivable    100
//     Cr. Revenue                      100
func (s *JournalService) PostARInvoice(
    ctx context.Context,
    tenantID uuid.UUID,
    invoiceID uuid.UUID,
    invoiceNo string,
    subtotal, tax, total float64,
    currency string,
    actor uuid.UUID,
) error {
    no, err := s.nextNo(ctx, tenantID)
    if err != nil { return err }

    je, err := entity.NewJournalEntry(tenantID, no, time.Now(),
        "AR Invoice "+invoiceNo, currency, actor)
    if err != nil { return err }

    // Dr. A/R (1100) = total
    arCode, _ := valueobject.NewAccountCode("1100")
    _ = je.AddLine(arCode, "Accounts Receivable", total, 0, invoiceNo)

    // Cr. Revenue (4000) = subtotal
    revCode, _ := valueobject.NewAccountCode("4000")
    _ = je.AddLine(revCode, "Revenue", 0, subtotal, invoiceNo)

    // Cr. VAT Payable (2200) = tax
    if tax > 0 {
        vatCode, _ := valueobject.NewAccountCode("2200")
        _ = je.AddLine(vatCode, "VAT Payable", 0, tax, invoiceNo)
    }

    if err := je.Post(actor); err != nil { return err }
    return s.journalRepo.Save(ctx, je)
}

// PostPaymentIn – รับชำระ
//   Dr. Cash/Bank       100
//     Cr. Accounts Receivable     100
func (s *JournalService) PostPaymentIn(
    ctx context.Context,
    tenantID uuid.UUID,
    paymentID uuid.UUID,
    paymentNo string,
    amount float64,
    currency string,
    actor uuid.UUID,
) error {
    no, err := s.nextNo(ctx, tenantID)
    if err != nil { return err }

    je, _ := entity.NewJournalEntry(tenantID, no, time.Now(),
        "Payment Received "+paymentNo, currency, actor)

    cashCode, _ := valueobject.NewAccountCode("1010") // Cash
    _ = je.AddLine(cashCode, "Cash/Bank", amount, 0, paymentNo)

    arCode, _ := valueobject.NewAccountCode("1100")
    _ = je.AddLine(arCode, "Accounts Receivable", 0, amount, paymentNo)

    if err := je.Post(actor); err != nil { return err }
    return s.journalRepo.Save(ctx, je)
}

// PostCOGS – ตัดต้นทุนขาย
//   Dr. Cost of Goods Sold    = qty * avg_cost
//     Cr. Inventory                    = qty * avg_cost
func (s *JournalService) PostCOGS(
    ctx context.Context,
    tenantID uuid.UUID,
    refNo string,
    cost float64,
    currency string,
    actor uuid.UUID,
) error {
    no, err := s.nextNo(ctx, tenantID)
    if err != nil { return err }

    je, _ := entity.NewJournalEntry(tenantID, no, time.Now(),
        "COGS "+refNo, currency, actor)

    cogsCode, _ := valueobject.NewAccountCode("5000")
    _ = je.AddLine(cogsCode, "Cost of Goods Sold", cost, 0, refNo)

    invCode, _ := valueobject.NewAccountCode("1200")
    _ = je.AddLine(invCode, "Inventory", 0, cost, refNo)

    if err := je.Post(actor); err != nil { return err }
    return s.journalRepo.Save(ctx, je)
}

func (s *JournalService) nextNo(ctx context.Context, tenantID uuid.UUID) (valueobject.DocumentNumber, error) {
    seq, err := s.seqGen.NextSeq(ctx, tenantID, valueobject.DocTypeJournalEntry, time.Now().Year())
    if err != nil { return "", err }
    return valueobject.NewDocumentNumber(valueobject.DocTypeJournalEntry, time.Now().Year(), seq)
}
```

#### `domain/service/reorder_service.go`
```go
package service

import (
    "context"

    "github.com/google/uuid"

    "icmongolang/internal/modules/erp/domain/entity"
    "icmongolang/internal/modules/erp/domain/repository"
)

// ReorderService – ตรวจสอบ reorder point + สร้าง draft PO
type ReorderService struct {
    invRepo     repository.InventoryRepository
    prodRepo    repository.ProductRepository
    orderRepo   repository.OrderRepository
    seqSvc      *DocumentNumberService
}

func NewReorderService(
    invRepo repository.InventoryRepository,
    prodRepo repository.ProductRepository,
    orderRepo repository.OrderRepository,
    seqSvc *DocumentNumberService,
) *ReorderService {
    return &ReorderService{invRepo: invRepo, prodRepo: prodRepo, orderRepo: orderRepo, seqSvc: seqSvc}
}

type ReorderSuggestion struct {
    ProductID    uuid.UUID
    SKU          string
    Name         string
    WarehouseID  uuid.UUID
    CurrentStock float64
    ReorderPoint float64
    SuggestedQty float64
    LeadTimeDays int
}

// ScanLowStock – หา products ที่ต่ำกว่า reorder point
func (s *ReorderService) ScanLowStock(ctx context.Context, tenantID uuid.UUID) ([]ReorderSuggestion, error) {
    items, _, err := s.invRepo.List(ctx, repository.InventoryFilter{
        TenantID: tenantID, LowStock: true,
    })
    if err != nil { return nil, err }

    var out []ReorderSuggestion
    for _, item := range items {
        p, err := s.prodRepo.FindByID(ctx, tenantID, item.ProductID)
        if err != nil { continue }
        if !p.TrackInventory || p.ReorderPoint <= 0 { continue }
        if item.QtyAvailable() > p.ReorderPoint { continue }
        out = append(out, ReorderSuggestion{
            ProductID: p.ID, SKU: p.SKU.String(), Name: p.Name,
            WarehouseID: item.WarehouseID,
            CurrentStock: item.QtyAvailable(),
            ReorderPoint: p.ReorderPoint,
            SuggestedQty: p.ReorderQty,
            LeadTimeDays: p.LeadTimeDays,
        })
    }
    return out, nil
}

// CreateDraftPO – สร้าง draft PO จาก suggestions (group by supplier — simplified)
func (s *ReorderService) CreateDraftPO(ctx context.Context, tenantID, supplierID uuid.UUID, suggestions []ReorderSuggestion) (*entity.Order, error) {
    no, err := s.seqSvc.Next(ctx, tenantID, valueobject.DocTypePurchaseOrder)
    if err != nil { return nil, err }

    order, err := entity.NewOrder(tenantID, no, valueobject.OrderTypePurchase, time.Now(), "THB", uuid.Nil)
    if err != nil { return nil, err }

    order.SetSupplier(supplierID, valueobject.Address{})

    for _, sug := range suggestions {
        p, _ := s.prodRepo.FindByID(ctx, tenantID, sug.ProductID)
        line, _ := entity.NewOrderLine(
            order.ID, 0, sug.ProductID,
            sug.SKU, sug.Name, p.UOM,
            sug.SuggestedQty, p.CostPrice.Amount,
            p.TaxType, p.TaxRate,
        )
        whID := sug.WarehouseID
        line.WarehouseID = &whID
        _ = order.AddLine(line)
    }

    if err := s.orderRepo.Save(ctx, order); err != nil { return nil, err }
    return order, nil
}

type valueobjectAlias = interface{}
var _ = valueobjectAlias(nil)
```

#### `domain/service/port/payment_port.go`
```go
package port

import (
    "context"

    "github.com/google/uuid"
)

type PaymentChargeRequest struct {
    TenantID      uuid.UUID
    CustomerID    uuid.UUID
    Amount        float64
    Currency      string
    Description   string
    ReferenceType string
    ReferenceID   uuid.UUID
}

type PaymentChargeResponse struct {
    ChargeID    string
    Status      string
    CheckoutURL string
}

type PaymentPort interface {
    CreateCharge(ctx context.Context, req PaymentChargeRequest) (*PaymentChargeResponse, error)
}
```

#### `domain/service/port/customer_port.go`
```go
package port

import (
    "context"

    "github.com/google/uuid"
)

type CustomerInfo struct {
    ID       uuid.UUID
    TenantID uuid.UUID
    Code     string
    Name     string
    Email    string
    Phone    string
    TaxID    string
    IsActive bool
}

type CustomerPort interface {
    Get(ctx context.Context, tenantID, customerID uuid.UUID) (*CustomerInfo, error)
    Exists(ctx context.Context, tenantID, customerID uuid.UUID) (bool, error)
}
```

#### `domain/service/port/notifier_port.go`
```go
package port

import "context"

type NotifyMessage struct {
    Channel  string
    Target   string
    Subject  string
    Body     string
    Metadata map[string]any
}

type NotifierPort interface {
    Send(ctx context.Context, tenantID string, msg NotifyMessage) error
}
```

#### `domain/service/port/pdf_renderer_port.go`
```go
package port

import "context"

type PDFRenderRequest struct {
    Template string
    Data     map[string]any
    Filename string
}

type PDFRendererPort interface {
    Render(ctx context.Context, req PDFRenderRequest) ([]byte, error)
}
```

#### `domain/service/port/code_generator_port.go`
```go
package port

import (
    "context"

    "github.com/google/uuid"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type CodeGeneratorPort interface {
    NextDocumentNumber(ctx context.Context, tenantID uuid.UUID, docType valueobject.DocType) (valueobject.DocumentNumber, error)
}
```

### A.6 Domain Events

#### `domain/event/product_events.go`
```go
package event

import (
    "time"

    "github.com/google/uuid"
)

const (
    TopicProductCreated  = "erp.product.created"
    TopicProductUpdated  = "erp.product.updated"
    TopicProductArchived = "erp.product.archived"
)

type ProductCreated struct {
    EventID    uuid.UUID `json:"event_id"`
    ProductID  uuid.UUID `json:"product_id"`
    TenantID   uuid.UUID `json:"tenant_id"`
    SKU        string    `json:"sku"`
    Name       string    `json:"name"`
    Type       string    `json:"type"`
    SellPrice  float64   `json:"sell_price"`
    Currency   string    `json:"currency"`
    OccurredAt time.Time `json:"occurred_at"`
}
```

#### `domain/event/order_events.go`
```go
package event

import (
    "time"

    "github.com/google/uuid"
)

const (
    TopicOrderCreated    = "erp.order.created"
    TopicOrderConfirmed  = "erp.order.confirmed"
    TopicOrderShipped    = "erp.order.shipped"
    TopicOrderReceived   = "erp.order.received"
    TopicOrderInvoiced   = "erp.order.invoiced"
    TopicOrderCancelled  = "erp.order.cancelled"
    TopicOrderClosed     = "erp.order.closed"
)

type OrderCreated struct {
    EventID    uuid.UUID  `json:"event_id"`
    OrderID    uuid.UUID  `json:"order_id"`
    TenantID   uuid.UUID  `json:"tenant_id"`
    OrderNo    string     `json:"order_no"`
    OrderType  string     `json:"order_type"`
    CustomerID *uuid.UUID `json:"customer_id,omitempty"`
    SupplierID *uuid.UUID `json:"supplier_id,omitempty"`
    TotalAmount float64   `json:"total_amount"`
    Currency   string     `json:"currency"`
    OccurredAt time.Time  `json:"occurred_at"`
}

type OrderConfirmed struct {
    EventID     uuid.UUID `json:"event_id"`
    OrderID     uuid.UUID `json:"order_id"`
    TenantID    uuid.UUID `json:"tenant_id"`
    OrderNo     string    `json:"order_no"`
    OrderType   string    `json:"order_type"`
    ConfirmedBy uuid.UUID `json:"confirmed_by"`
    OccurredAt  time.Time `json:"occurred_at"`
}

type OrderShipped struct {
    EventID     uuid.UUID `json:"event_id"`
    OrderID     uuid.UUID `json:"order_id"`
    TenantID    uuid.UUID `json:"tenant_id"`
    OrderNo     string    `json:"order_no"`
    CustomerID  uuid.UUID `json:"customer_id"`
    ShippedLines []ShippedLine `json:"shipped_lines"`
    OccurredAt  time.Time `json:"occurred_at"`
}

type ShippedLine struct {
    LineID      uuid.UUID `json:"line_id"`
    ProductID   uuid.UUID `json:"product_id"`
    WarehouseID uuid.UUID `json:"warehouse_id"`
    Quantity    float64   `json:"quantity"`
}
```

#### `domain/event/invoice_events.go`
```go
package event

import (
    "time"

    "github.com/google/uuid"
)

const (
    TopicInvoiceDraft    = "erp.invoice.draft"
    TopicInvoiceIssued   = "erp.invoice.issued"
    TopicInvoicePaid     = "erp.invoice.paid"
    TopicInvoiceOverdue  = "erp.invoice.overdue"
    TopicInvoiceVoided   = "erp.invoice.voided"
)

type InvoiceIssued struct {
    EventID     uuid.UUID  `json:"event_id"`
    InvoiceID   uuid.UUID  `json:"invoice_id"`
    TenantID    uuid.UUID  `json:"tenant_id"`
    InvoiceNo   string     `json:"invoice_no"`
    InvoiceType string     `json:"invoice_type"`
    CustomerID  *uuid.UUID `json:"customer_id,omitempty"`
    SupplierID  *uuid.UUID `json:"supplier_id,omitempty"`
    OrderID     *uuid.UUID `json:"order_id,omitempty"`
    Subtotal    float64    `json:"subtotal"`
    TaxAmount   float64    `json:"tax_amount"`
    TotalAmount float64    `json:"total_amount"`
    Currency    string     `json:"currency"`
    IssueDate   time.Time  `json:"issue_date"`
    DueDate     time.Time  `json:"due_date"`
    OccurredAt  time.Time  `json:"occurred_at"`
}

type InvoicePaid struct {
    EventID     uuid.UUID `json:"event_id"`
    InvoiceID   uuid.UUID `json:"invoice_id"`
    TenantID    uuid.UUID `json:"tenant_id"`
    InvoiceNo   string    `json:"invoice_no"`
    AmountPaid  float64   `json:"amount_paid"`
    Currency    string    `json:"currency"`
    OccurredAt  time.Time `json:"occurred_at"`
}

type InvoiceOverdue struct {
    EventID     uuid.UUID `json:"event_id"`
    InvoiceID   uuid.UUID `json:"invoice_id"`
    TenantID    uuid.UUID `json:"tenant_id"`
    InvoiceNo   string    `json:"invoice_no"`
    AmountDue   float64   `json:"amount_due"`
    Currency    string    `json:"currency"`
    DaysOverdue int       `json:"days_overdue"`
    OccurredAt  time.Time `json:"occurred_at"`
}
```

#### `domain/event/inventory_events.go`
```go
package event

import (
    "time"

    "github.com/google/uuid"
)

const (
    TopicStockLow       = "erp.inventory.low_stock"
    TopicStockOut       = "erp.inventory.out_of_stock"
    TopicStockAdjusted  = "erp.inventory.adjusted"
    TopicStockTransferred = "erp.inventory.transferred"
)

type StockLow struct {
    EventID      uuid.UUID `json:"event_id"`
    TenantID     uuid.UUID `json:"tenant_id"`
    ProductID    uuid.UUID `json:"product_id"`
    SKU          string    `json:"sku"`
    WarehouseID  uuid.UUID `json:"warehouse_id"`
    CurrentStock float64   `json:"current_stock"`
    ReorderPoint float64   `json:"reorder_point"`
    OccurredAt   time.Time `json:"occurred_at"`
}

type StockOut struct {
    EventID      uuid.UUID `json:"event_id"`
    TenantID     uuid.UUID `json:"tenant_id"`
    ProductID    uuid.UUID `json:"product_id"`
    SKU          string    `json:"sku"`
    WarehouseID  uuid.UUID `json:"warehouse_id"`
    OccurredAt   time.Time `json:"occurred_at"`
}

type StockAdjusted struct {
    EventID     uuid.UUID `json:"event_id"`
    TenantID    uuid.UUID `json:"tenant_id"`
    ProductID   uuid.UUID `json:"product_id"`
    WarehouseID uuid.UUID `json:"warehouse_id"`
    OldQty      float64   `json:"old_qty"`
    NewQty      float64   `json:"new_qty"`
    Delta       float64   `json:"delta"`
    Reason      string    `json:"reason"`
    ActorID     uuid.UUID `json:"actor_id"`
    OccurredAt  time.Time `json:"occurred_at"`
}
```

#### `domain/event/payment_events.go`
```go
package event

import (
    "time"

    "github.com/google/uuid"
)

const (
    TopicPaymentCreated  = "erp.payment.created"
    TopicPaymentConfirmed = "erp.payment.confirmed"
    TopicPaymentApplied  = "erp.payment.applied"
    TopicPaymentReversed = "erp.payment.reversed"
)

type PaymentConfirmed struct {
    EventID     uuid.UUID `json:"event_id"`
    PaymentID   uuid.UUID `json:"payment_id"`
    TenantID    uuid.UUID `json:"tenant_id"`
    PaymentNo   string    `json:"payment_no"`
    Direction   string    `json:"direction"`
    Amount      float64   `json:"amount"`
    Currency    string    `json:"currency"`
    OccurredAt  time.Time `json:"occurred_at"`
}
```

### A.7 Domain Errors

#### `domain/errors/errors.go`
```go
package domainerrors

import "errors"

var (
    // Product
    ErrProductNotFound    = errors.New("product not found")
    ErrInvalidSKU         = errors.New("invalid sku")
    ErrInvalidProductName = errors.New("invalid product name")
    ErrSKUDuplicate       = errors.New("sku already exists")
    ErrInvalidReorderPolicy = errors.New("invalid reorder policy")

    // Warehouse
    ErrWarehouseNotFound = errors.New("warehouse not found")
    ErrInvalidWarehouse  = errors.New("invalid warehouse")
    ErrWarehouseDuplicate = errors.New("warehouse code duplicate")

    // Inventory
    ErrInventoryNotFound       = errors.New("inventory not found")
    ErrInsufficientStock       = errors.New("insufficient stock")
    ErrInsufficientAvailable   = errors.New("insufficient available stock")
    ErrInvalidDirection        = errors.New("invalid direction")
    ErrInvalidMovementType     = errors.New("invalid movement type")
    ErrInvalidQuantity         = errors.New("invalid quantity")
    ErrInvalidReleaseAmount    = errors.New("invalid release amount")
    ErrReservationNotActive    = errors.New("reservation not active")

    // Order
    ErrOrderNotFound          = errors.New("order not found")
    ErrOrderNotEditable       = errors.New("order not editable")
    ErrOrderAlreadyTerminal   = errors.New("order already terminal")
    ErrEmptyOrder             = errors.New("order has no lines")
    ErrOrderLineNotFound      = errors.New("order line not found")
    ErrInvalidOrderType       = errors.New("invalid order type")
    ErrOrderTypeMismatch      = errors.New("order type mismatch")
    ErrNotSalesOrder          = errors.New("not a sales order")
    ErrNotPurchaseOrder       = errors.New("not a purchase order")
    ErrCustomerRequired       = errors.New("customer required")
    ErrSupplierRequired       = errors.New("supplier required")
    ErrCannotCancelShippedOrder = errors.New("cannot cancel shipped order")
    ErrOverShip               = errors.New("cannot ship more than ordered")
    ErrOverReceive            = errors.New("cannot receive more than ordered")
    ErrOverInvoice            = errors.New("cannot invoice more than ordered")
    ErrNoLinesToShip          = errors.New("no lines to ship")
    ErrNoLinesToReceive       = errors.New("no lines to receive")

    // Invoice
    ErrInvoiceNotFound         = errors.New("invoice not found")
    ErrInvoiceNotEditable      = errors.New("invoice not editable")
    ErrInvoiceNotPayable       = errors.New("invoice not payable")
    ErrInvoiceLineNotFound     = errors.New("invoice line not found")
    ErrInvalidInvoiceType      = errors.New("invalid invoice type")
    ErrInvalidDueDate          = errors.New("invalid due date")
    ErrEmptyInvoice            = errors.New("invoice has no lines")
    ErrOverpayment             = errors.New("payment exceeds amount due")
    ErrCannotVoidPaidInvoice   = errors.New("cannot void paid invoice")
    ErrInvoiceAlreadyVoid      = errors.New("invoice already void")

    // Payment
    ErrPaymentNotFound        = errors.New("payment not found")
    ErrPaymentNotConfirmed    = errors.New("payment not confirmed")
    ErrPaymentReversed        = errors.New("payment reversed")
    ErrPaymentAlreadyReversed = errors.New("payment already reversed")
    ErrInvalidPaymentMethod   = errors.New("invalid payment method")
    ErrInvalidPaymentAmount   = errors.New("invalid payment amount")
    ErrDirectionMismatch      = errors.New("direction mismatch")
    ErrOverAllocation         = errors.New("over allocation")
    ErrOverDeallocation       = errors.New("over deallocation")
    ErrOverRefund             = errors.New("over refund")
    ErrAllocationNotFound     = errors.New("allocation not found")

    // Credit Note
    ErrCreditNoteNotFound  = errors.New("credit note not found")
    ErrInvalidNoteType     = errors.New("invalid note type")
    ErrAlreadyVoid         = errors.New("already void")

    // Journal
    ErrJournalNotFound        = errors.New("journal entry not found")
    ErrJournalNotEditable     = errors.New("journal not editable")
    ErrJournalNeedsTwoLines   = errors.New("journal needs at least 2 lines")
    ErrJournalNotBalanced     = errors.New("journal not balanced")
    ErrCannotHaveBothDebitCredit = errors.New("line cannot have both debit and credit")
    ErrDescriptionRequired    = errors.New("description required")

    // Supplier
    ErrSupplierNotFound = errors.New("supplier not found")
    ErrInvalidSupplier  = errors.New("invalid supplier")

    // Price List
    ErrPriceListNotFound   = errors.New("price list not found")
    ErrInvalidPriceList    = errors.New("invalid price list")
    ErrPriceEntryDuplicate = errors.New("price entry duplicate")

    // Validation
    ErrInvalidTenantID   = errors.New("invalid tenant id")
    ErrInvalidCustomerID = errors.New("invalid customer id")
    ErrInvalidDocNumber  = errors.New("invalid document number")
    ErrInvalidDocType    = errors.New("invalid document type")
    ErrInvalidAccountCode = errors.New("invalid account code")
    ErrInvalidTaxType    = errors.New("invalid tax type")
    ErrInvalidTaxRate    = errors.New("invalid tax rate")
    ErrInvalidAmount     = errors.New("invalid amount")
    ErrInvalidPrice      = errors.New("invalid price")
    ErrInvalidDiscount   = errors.New("invalid discount")
    ErrDiscountExceedsPrice = errors.New("discount exceeds price")
    ErrNegativeQuantity  = errors.New("negative quantity")
    ErrNegativeMoney     = errors.New("negative money")
    ErrUOMRequired       = errors.New("uom required")
    ErrUOMMismatch       = errors.New("uom mismatch")
    ErrCurrencyMismatch  = errors.New("currency mismatch")
    ErrInvalidCurrency   = errors.New("invalid currency")
    ErrInvalidStatusTransition = errors.New("invalid status transition")
    ErrInvalidDateRange  = errors.New("invalid date range")
    ErrReasonRequired    = errors.New("reason required")

    // Infrastructure
    ErrPersistenceFailure = errors.New("persistence failure")
    ErrPDFRenderFailure   = errors.New("pdf render failure")
    ErrNotifierFailure    = errors.New("notification failure")
    ErrPaymentGateway     = errors.New("payment gateway error")
)
```

### A.8 Unit Tests (Domain)

#### `domain/entity/order_test.go`
```go
package entity_test

import (
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "icmongolang/internal/modules/erp/domain/entity"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

func newSalesOrder(t *testing.T) *entity.Order {
    no, _ := valueobject.NewDocumentNumber(valueobject.DocTypeSalesOrder, 2026, 1)
    o, err := entity.NewOrder(uuid.New(), no, valueobject.OrderTypeSales, time.Now(), "THB", uuid.New())
    require.NoError(t, err)
    return o
}

func TestOrder_AddLine_CalculatesTotals(t *testing.T) {
    o := newSalesOrder(t)

    line, _ := entity.NewOrderLine(
        o.ID, 0, uuid.New(), "SKU-001", "Product A", "PCS",
        10, 100, valueobject.TaxTypeVAT, 0.07,
    )
    require.NoError(t, o.AddLine(line))

    // subtotal 1000, tax 70, total 1070
    assert.InDelta(t, 1000.0, o.Subtotal, 0.01)
    assert.InDelta(t, 70.0, o.TaxAmount, 0.01)
    assert.InDelta(t, 1070.0, o.TotalAmount, 0.01)
}

func TestOrder_AddLine_DuplicateProduct_Merges(t *testing.T) {
    o := newSalesOrder(t)
    pid := uuid.New()

    l1, _ := entity.NewOrderLine(o.ID, 0, pid, "SKU-1", "A", "PCS", 5, 100, valueobject.TaxTypeVAT, 0.07)
    l2, _ := entity.NewOrderLine(o.ID, 0, pid, "SKU-1", "A", "PCS", 3, 100, valueobject.TaxTypeVAT, 0.07)

    require.NoError(t, o.AddLine(l1))
    require.NoError(t, o.AddLine(l2))

    assert.Equal(t, 1, len(o.Lines))
    assert.InDelta(t, 8.0, o.Lines[0].Quantity, 0.01)
}

func TestOrder_Lifecycle_ConfirmShipInvoice(t *testing.T) {
    o := newSalesOrder(t)
    custID := uuid.New()
    require.NoError(t, o.SetCustomer(custID, valueobject.Address{Country: "TH"}))

    line, _ := entity.NewOrderLine(o.ID, 0, uuid.New(), "SKU-1", "A", "PCS", 10, 100, valueobject.TaxTypeVAT, 0.07)
    _ = o.AddLine(line)

    // Submit
    require.NoError(t, o.SubmitForApproval())
    assert.Equal(t, valueobject.OrderStatusPendingApproval, o.Status)

    // Confirm
    require.NoError(t, o.Confirm(uuid.New()))
    assert.Equal(t, valueobject.OrderStatusConfirmed, o.Status)

    // Ship full
    qty := map[uuid.UUID]float64{line.ID: 10}
    require.NoError(t, o.Ship(qty))
    assert.Equal(t, valueobject.OrderStatusShipped, o.Status)

    // Invoice
    require.NoError(t, o.Invoice(qty))
    assert.Equal(t, valueobject.OrderStatusInvoiced, o.Status)

    // Paid
    require.NoError(t, o.MarkPaid())
    assert.Equal(t, valueobject.OrderStatusPaid, o.Status)

    // Close
    require.NoError(t, o.Close())
    assert.Equal(t, valueobject.OrderStatusClosed, o.Status)
    assert.True(t, o.IsTerminal())
}

func TestOrder_Ship_Partial(t *testing.T) {
    o := newSalesOrder(t)
    _ = o.SetCustomer(uuid.New(), valueobject.Address{})
    line, _ := entity.NewOrderLine(o.ID, 0, uuid.New(), "SKU-1", "A", "PCS", 10, 100, valueobject.TaxTypeVAT, 0.07)
    _ = o.AddLine(line)
    _ = o.SubmitForApproval()
    _ = o.Confirm(uuid.New())

    // Ship 6 of 10
    require.NoError(t, o.Ship(map[uuid.UUID]float64{line.ID: 6}))
    assert.Equal(t, valueobject.OrderStatusPartiallyShipped, o.Status)
    assert.False(t, o.IsFullyShipped())

    // Ship remaining
    require.NoError(t, o.Ship(map[uuid.UUID]float64{line.ID: 4}))
    assert.Equal(t, valueobject.OrderStatusShipped, o.Status)
    assert.True(t, o.IsFullyShipped())
}

func TestOrder_Ship_OverQty_Fails(t *testing.T) {
    o := newSalesOrder(t)
    _ = o.SetCustomer(uuid.New(), valueobject.Address{})
    line, _ := entity.NewOrderLine(o.ID, 0, uuid.New(), "SKU-1", "A", "PCS", 10, 100, valueobject.TaxTypeVAT, 0.07)
    _ = o.AddLine(line)
    _ = o.SubmitForApproval()
    _ = o.Confirm(uuid.New())

    err := o.Ship(map[uuid.UUID]float64{line.ID: 11})
    assert.ErrorIs(t, err, domainerrors.ErrOverShip)
}

func TestOrder_Cancel_AfterShip_Fails(t *testing.T) {
    o := newSalesOrder(t)
    _ = o.SetCustomer(uuid.New(), valueobject.Address{})
    line, _ := entity.NewOrderLine(o.ID, 0, uuid.New(), "SKU-1", "A", "PCS", 10, 100, valueobject.TaxTypeVAT, 0.07)
    _ = o.AddLine(line)
    _ = o.SubmitForApproval()
    _ = o.Confirm(uuid.New())
    _ = o.Ship(map[uuid.UUID]float64{line.ID: 10})

    err := o.Cancel("test", uuid.New())
    assert.ErrorIs(t, err, domainerrors.ErrCannotCancelShippedOrder)
}
```

#### `domain/entity/inventory_test.go`
```go
package entity_test

import (
    "testing"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"

    "icmongolang/internal/modules/erp/domain/entity"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
)

func TestInventory_Movement_In(t *testing.T) {
    item := entity.NewInventoryItem(uuid.New(), uuid.New(), uuid.New())
    require := assert.New(t)

    // รับเข้า 100 @ cost 50
    require.NoError(item.ApplyMovement(1, 100, 50))
    require.Equal(100.0, item.QtyOnHand)
    require.Equal(50.0, item.AvgCost)

    // รับเข้า 100 @ cost 70 → avg = 60
    require.NoError(item.ApplyMovement(1, 100, 70))
    require.Equal(200.0, item.QtyOnHand)
    require.InDelta(60.0, item.AvgCost, 0.01)
}

func TestInventory_Movement_Out_Insufficient(t *testing.T) {
    item := entity.NewInventoryItem(uuid.New(), uuid.New(), uuid.New())
    item.QtyOnHand = 10

    err := item.ApplyMovement(-1, 20, 0)
    assert.ErrorIs(t, err, domainerrors.ErrInsufficientStock)
}

func TestInventory_Reserve_Release(t *testing.T) {
    item := entity.NewInventoryItem(uuid.New(), uuid.New(), uuid.New())
    item.QtyOnHand = 100

    require := assert.New(t)
    require.NoError(item.Reserve(30))
    require.Equal(30.0, item.QtyReserved)
    require.Equal(70.0, item.QtyAvailable())

    require.NoError(item.Release(10))
    require.Equal(20.0, item.QtyReserved)
    require.Equal(80.0, item.QtyAvailable())
}

func TestInventory_ConsumeReservation(t *testing.T) {
    item := entity.NewInventoryItem(uuid.New(), uuid.New(), uuid.New())
    item.QtyOnHand = 100
    _ = item.Reserve(30)

    require := assert.New(t)
    require.NoError(item.ConsumeReservation(30))
    require.Equal(70.0, item.QtyOnHand)   // stock ลด
    require.Equal(0.0, item.QtyReserved)  // reserved ลด
    require.Equal(70.0, item.QtyAvailable())
}
```

#### `domain/service/proration_service_test.go` (simplified)
```go
package service_test

import (
    "testing"

    "github.com/stretchr/testify/assert"

    "icmongolang/internal/modules/erp/domain/service"
)

func TestTaxCalculator_VAT_Exclusive(t *testing.T) {
    c := service.NewTaxCalculator()
    r := c.Calculate(service.TaxRequest{
        Subtotal: 1000,
        TaxType:  "VAT",
        TaxRate:  0.07,
    })
    assert.InDelta(t, 1000.0, r.BaseAmount, 0.01)
    assert.InDelta(t, 70.0, r.TaxAmount, 0.01)
    assert.InDelta(t, 1070.0, r.TotalAmount, 0.01)
}

func TestTaxCalculator_VAT_Inclusive(t *testing.T) {
    c := service.NewTaxCalculator()
    r := c.Calculate(service.TaxRequest{
        Subtotal: 1070,
        TaxType:  "VAT",
        TaxRate:  0.07,
        Inclusive: true,
    })
    assert.InDelta(t, 1000.0, r.BaseAmount, 0.01)
    assert.InDelta(t, 70.0, r.TaxAmount, 0.01)
    assert.InDelta(t, 1070.0, r.TotalAmount, 0.01)
}

func TestTaxCalculator_WHT(t *testing.T) {
    c := service.NewTaxCalculator()
    r := c.Calculate(service.TaxRequest{
        Subtotal: 10000,
        TaxType:  "WHT",
        TaxRate:  0.03,
    })
    assert.InDelta(t, 300.0, r.WHTAmount, 0.01)
    assert.InDelta(t, 9700.0, r.NetPayable, 0.01)
}
```

---

## 🅱️ PART 4B — APPLICATION LAYER

### B.1 DTOs

```go
package application

import "time"

// ============================================================
// PRODUCT
// ============================================================

type CreateProductInput struct {
    TenantID     string   `json:"-"`
    SKU          string   `json:"sku" binding:"required"`
    Name         string   `json:"name" binding:"required"`
    Description  string   `json:"description,omitempty"`
    Type         string   `json:"type" binding:"required"`
    Category     string   `json:"category,omitempty"`
    Brand        string   `json:"brand,omitempty"`
    UOM          string   `json:"uom" binding:"required"`
    CostPrice    float64  `json:"cost_price" binding:"gte=0"`
    SellPrice    float64  `json:"sell_price" binding:"gte=0"`
    Currency     string   `json:"currency,omitempty"`
    TaxType      string   `json:"tax_type,omitempty"`
    TaxRate      float64  `json:"tax_rate,omitempty"`
    TrackInventory bool   `json:"track_inventory"`
    ReorderPoint float64  `json:"reorder_point,omitempty"`
    ReorderQty   float64  `json:"reorder_qty,omitempty"`
    LeadTimeDays int      `json:"lead_time_days,omitempty"`
    SafetyStock  float64  `json:"safety_stock,omitempty"`
    Tags         []string `json:"tags,omitempty"`
    ActorID      string   `json:"-"`
}

type ProductResponse struct {
    ID           string    `json:"id"`
    TenantID     string    `json:"tenant_id"`
    SKU          string    `json:"sku"`
    Name         string    `json:"name"`
    Description  string    `json:"description,omitempty"`
    Type         string    `json:"type"`
    Category     string    `json:"category,omitempty"`
    Brand        string    `json:"brand,omitempty"`
    UOM          string    `json:"uom"`
    CostPrice    float64   `json:"cost_price"`
    SellPrice    float64   `json:"sell_price"`
    Currency     string    `json:"currency"`
    TaxType      string    `json:"tax_type"`
    TaxRate      float64   `json:"tax_rate"`
    TrackInventory bool    `json:"track_inventory"`
    ReorderPoint float64   `json:"reorder_point"`
    ReorderQty   float64   `json:"reorder_qty"`
    LeadTimeDays int       `json:"lead_time_days"`
    IsActive     bool      `json:"is_active"`
    IsSellable   bool      `json:"is_sellable"`
    Tags         []string  `json:"tags,omitempty"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

type ListProductsInput struct {
    TenantID string
    Type     string
    Category string
    Search   string
    Active   *bool
    Tags     []string
    Page     int
    PageSize int
}

// ============================================================
// WAREHOUSE
// ============================================================

type CreateWarehouseInput struct {
    TenantID  string `json:"-"`
    Code      string `json:"code" binding:"required"`
    Name      string `json:"name" binding:"required"`
    Type      string `json:"type,omitempty"`
    Address   AddressDTO `json:"address"`
    ActorID   string `json:"-"`
}

type AddressDTO struct {
    Line1    string `json:"line1,omitempty"`
    Line2    string `json:"line2,omitempty"`
    District string `json:"district,omitempty"`
    Amphoe   string `json:"amphoe,omitempty"`
    Province string `json:"province,omitempty"`
    Postcode string `json:"postcode,omitempty"`
    Country  string `json:"country" binding:"required"`
    TaxID    string `json:"tax_id,omitempty"`
}

type WarehouseResponse struct {
    ID        string     `json:"id"`
    Code      string     `json:"code"`
    Name      string     `json:"name"`
    Type      string     `json:"type"`
    Address   AddressDTO `json:"address"`
    IsActive  bool       `json:"is_active"`
    IsDefault bool       `json:"is_default"`
    CreatedAt time.Time  `json:"created_at"`
}

// ============================================================
// ORDER
// ============================================================

type CreateOrderInput struct {
    TenantID    string       `json:"-"`
    OrderType   string       `json:"order_type" binding:"required"` // SALES | PURCHASE
    CustomerID  string       `json:"customer_id,omitempty"`
    SupplierID  string       `json:"supplier_id,omitempty"`
    OrderDate   time.Time    `json:"order_date"`
    ExpectedAt  *time.Time   `json:"expected_at,omitempty"`
    Currency    string       `json:"currency,omitempty"`
    POReference string       `json:"po_reference,omitempty"`
    ShipTo      *AddressDTO  `json:"ship_to,omitempty"`
    BillTo      *AddressDTO  `json:"bill_to,omitempty"`
    Lines       []OrderLineDTO `json:"lines" binding:"required,min=1"`
    Notes       string       `json:"notes,omitempty"`
    ActorID     string       `json:"-"`
}

type OrderLineDTO struct {
    ProductID   string  `json:"product_id" binding:"required"`
    Quantity    float64 `json:"quantity" binding:"required,gt=0"`
    UnitPrice   float64 `json:"unit_price,omitempty"` // ถ้าไม่ส่ง → ใช้ default
    Discount    float64 `json:"discount,omitempty"`
    DiscountPct float64 `json:"discount_pct,omitempty"`
    WarehouseID string  `json:"warehouse_id,omitempty"`
    Notes       string  `json:"notes,omitempty"`
}

type OrderResponse struct {
    ID           string             `json:"id"`
    TenantID     string             `json:"tenant_id"`
    OrderNo      string             `json:"order_no"`
    Type         string             `json:"type"`
    Status       string             `json:"status"`
    CustomerID   string             `json:"customer_id,omitempty"`
    SupplierID   string             `json:"supplier_id,omitempty"`
    OrderDate    time.Time          `json:"order_date"`
    ExpectedAt   *time.Time         `json:"expected_at,omitempty"`
    ConfirmedAt  *time.Time         `json:"confirmed_at,omitempty"`
    ShippedAt    *time.Time         `json:"shipped_at,omitempty"`
    ReceivedAt   *time.Time         `json:"received_at,omitempty"`
    Lines        []OrderLineResponse `json:"lines"`
    Subtotal     float64            `json:"subtotal"`
    DiscountAmount float64          `json:"discount_amount"`
    TaxAmount    float64            `json:"tax_amount"`
    ShippingFee  float64            `json:"shipping_fee"`
    TotalAmount  float64            `json:"total_amount"`
    Currency     string             `json:"currency"`
    POReference  string             `json:"po_reference,omitempty"`
    Notes        string             `json:"notes,omitempty"`
    CreatedAt    time.Time          `json:"created_at"`
}

type OrderLineResponse struct {
    ID          string  `json:"id"`
    LineNo      int     `json:"line_no"`
    ProductID   string  `json:"product_id"`
    SKU         string  `json:"sku"`
    Name        string  `json:"name"`
    Quantity    float64 `json:"quantity"`
    UOM         string  `json:"uom"`
    UnitPrice   float64 `json:"unit_price"`
    Discount    float64 `json:"discount"`
    Subtotal    float64 `json:"subtotal"`
    TaxAmount   float64 `json:"tax_amount"`
    LineTotal   float64 `json:"line_total"`
    QtyShipped  float64 `json:"qty_shipped"`
    QtyReceived float64 `json:"qty_received"`
    QtyInvoiced float64 `json:"qty_invoiced"`
}

type ShipOrderInput struct {
    TenantID string                 `json:"-"`
    OrderID  string                 `json:"-"`
    Lines    map[string]float64     `json:"lines"` // lineID → qty
    Tracking string                 `json:"tracking,omitempty"`
    ActorID  string                 `json:"-"`
}

type ReceiveOrderInput struct {
    TenantID string                 `json:"-"`
    OrderID  string                 `json:"-"`
    Lines    map[string]float64     `json:"lines"`
    ActorID  string                 `json:"-"`
}

// ============================================================
// INVOICE
// ============================================================

type CreateInvoiceInput struct {
    TenantID      string          `json:"-"`
    InvoiceType   string          `json:"invoice_type" binding:"required"` // AR | AP
    CustomerID    string          `json:"customer_id,omitempty"`
    SupplierID    string          `json:"supplier_id,omitempty"`
    OrderID       string          `json:"order_id,omitempty"`
    IssueDate     time.Time       `json:"issue_date"`
    DueDate       time.Time       `json:"due_date"`
    Currency      string          `json:"currency,omitempty"`
    PaymentTerms  string          `json:"payment_terms,omitempty"`
    BillTo        *AddressDTO     `json:"bill_to,omitempty"`
    Lines         []InvoiceLineDTO `json:"lines" binding:"required,min=1"`
    Notes         string          `json:"notes,omitempty"`
    ActorID       string          `json:"-"`
}

type InvoiceLineDTO struct {
    ProductID   string  `json:"product_id" binding:"required"`
    Quantity    float64 `json:"quantity" binding:"required,gt=0"`
    UnitPrice   float64 `json:"unit_price" binding:"required,gte=0"`
    Discount    float64 `json:"discount,omitempty"`
    OrderLineID string  `json:"order_line_id,omitempty"`
}

type InvoiceResponse struct {
    ID          string                `json:"id"`
    InvoiceNo   string                `json:"invoice_no"`
    Type        string                `json:"type"`
    Status      string                `json:"status"`
    CustomerID  string                `json:"customer_id,omitempty"`
    SupplierID  string                `json:"supplier_id,omitempty"`
    OrderID     string                `json:"order_id,omitempty"`
    IssueDate   time.Time             `json:"issue_date"`
    DueDate     time.Time             `json:"due_date"`
    PaidAt      *time.Time            `json:"paid_at,omitempty"`
    Lines       []InvoiceLineResponse `json:"lines"`
    Subtotal    float64               `json:"subtotal"`
    TaxAmount   float64               `json:"tax_amount"`
    TotalAmount float64               `json:"total_amount"`
    AmountPaid  float64               `json:"amount_paid"`
    AmountDue   float64               `json:"amount_due"`
    Currency    string                `json:"currency"`
    DaysOverdue int                   `json:"days_overdue"`
    AgingBucket string                `json:"aging_bucket"`
    IsOverdue   bool                  `json:"is_overdue"`
    CreatedAt   time.Time             `json:"created_at"`
}

type InvoiceLineResponse struct {
    ID        string  `json:"id"`
    LineNo    int     `json:"line_no"`
    ProductID string  `json:"product_id"`
    SKU       string  `json:"sku"`
    Name      string  `json:"name"`
    Quantity  float64 `json:"quantity"`
    UOM       string  `json:"uom"`
    UnitPrice float64 `json:"unit_price"`
    Subtotal  float64 `json:"subtotal"`
    TaxAmount float64 `json:"tax_amount"`
    LineTotal float64 `json:"line_total"`
}

// ============================================================
// PAYMENT
// ============================================================

type CreatePaymentInput struct {
    TenantID    string    `json:"-"`
    Direction   string    `json:"direction" binding:"required,oneof=IN OUT"`
    CustomerID  string    `json:"customer_id,omitempty"`
    SupplierID  string    `json:"supplier_id,omitempty"`
    Method      string    `json:"method" binding:"required"`
    Amount      float64   `json:"amount" binding:"required,gt=0"`
    Currency    string    `json:"currency,omitempty"`
    PaymentDate time.Time `json:"payment_date"`
    Reference   string    `json:"reference,omitempty"`
    BankAccount string    `json:"bank_account,omitempty"`
    Notes       string    `json:"notes,omitempty"`
    Allocations []AllocationDTO `json:"allocations,omitempty"`
    ActorID     string    `json:"-"`
}

type AllocationDTO struct {
    InvoiceID string  `json:"invoice_id"`
    Amount    float64 `json:"amount" binding:"gt=0"`
}

type PaymentResponse struct {
    ID          string               `json:"id"`
    PaymentNo   string               `json:"payment_no"`
    Direction   string               `json:"direction"`
    Method      string               `json:"method"`
    Status      string               `json:"status"`
    CustomerID  string               `json:"customer_id,omitempty"`
    SupplierID  string               `json:"supplier_id,omitempty"`
    Amount      float64              `json:"amount"`
    Allocated   float64              `json:"allocated"`
    Unallocated float64              `json:"unallocated"`
    Currency    string               `json:"currency"`
    PaymentDate time.Time            `json:"payment_date"`
    Allocations []PaymentAllocationResponse `json:"allocations"`
    CreatedAt   time.Time            `json:"created_at"`
}

type PaymentAllocationResponse struct {
    InvoiceID uuid.UUID `json:"invoice_id"`
    Amount    float64   `json:"amount"`
}

// ============================================================
// INVENTORY
// ============================================================

type InventoryResponse struct {
    ID            string    `json:"id"`
    ProductID     string    `json:"product_id"`
    SKU           string    `json:"sku"`
    ProductName   string    `json:"product_name"`
    WarehouseID   string    `json:"warehouse_id"`
    WarehouseName string    `json:"warehouse_name"`
    QtyOnHand     float64   `json:"qty_on_hand"`
    QtyReserved   float64   `json:"qty_reserved"`
    QtyAvailable  float64   `json:"qty_available"`
    AvgCost       float64   `json:"avg_cost"`
    Value         float64   `json:"value"`
    BinLocation   string    `json:"bin_location,omitempty"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type AdjustInventoryInput struct {
    TenantID    string  `json:"-"`
    ProductID   string  `json:"product_id" binding:"required"`
    WarehouseID string  `json:"warehouse_id" binding:"required"`
    NewQty      float64 `json:"new_qty" binding:"gte=0"`
    Reason      string  `json:"reason" binding:"required"`
    ActorID     string  `json:"-"`
}

type TransferStockInput struct {
    TenantID          string  `json:"-"`
    ProductID         string  `json:"product_id" binding:"required"`
    FromWarehouseID   string  `json:"from_warehouse_id" binding:"required"`
    ToWarehouseID     string  `json:"to_warehouse_id" binding:"required"`
    Quantity          float64 `json:"quantity" binding:"required,gt=0"`
    Reference         string  `json:"reference,omitempty"`
    ActorID           string  `json:"-"`
}

type StockMovementResponse struct {
    ID          int64     `json:"id"`
    ProductID   string    `json:"product_id"`
    WarehouseID string    `json:"warehouse_id"`
    Type        string    `json:"type"`
    Direction   int       `json:"direction"`
    Quantity    float64   `json:"quantity"`
    UnitCost    float64   `json:"unit_cost"`
    TotalCost   float64   `json:"total_cost"`
    RefType     string    `json:"ref_type,omitempty"`
    RefNumber   string    `json:"ref_number,omitempty"`
    QtyOnHandAfter float64 `json:"qty_on_hand_after"`
    OccurredAt  time.Time `json:"occurred_at"`
}
```

### B.2 Mappers

```go
package application

import (
    "icmongolang/internal/modules/erp/domain/entity"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

func toProductResponse(p *entity.Product) *ProductResponse {
    if p == nil { return nil }
    return &ProductResponse{
        ID: p.ID.String(), TenantID: p.TenantID.String(),
        SKU: p.SKU.String(), Name: p.Name, Description: p.Description,
        Type: string(p.Type), Category: p.Category, Brand: p.Brand,
        UOM: p.UOM,
        CostPrice: p.CostPrice.Amount, SellPrice: p.SellPrice.Amount,
        Currency: p.Currency,
        TaxType: string(p.TaxType), TaxRate: p.TaxRate,
        TrackInventory: p.TrackInventory,
        ReorderPoint: p.ReorderPoint, ReorderQty: p.ReorderQty,
        LeadTimeDays: p.LeadTimeDays,
        IsActive: p.IsActive, IsSellable: p.IsSellable,
        Tags: p.Tags,
        CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
    }
}

func toOrderResponse(o *entity.Order, productNames map[string]string) *OrderResponse {
    if o == nil { return nil }
    r := &OrderResponse{
        ID: o.ID.String(), TenantID: o.TenantID.String(),
        OrderNo: o.OrderNo.String(), Type: string(o.Type),
        Status: string(o.Status),
        OrderDate: o.OrderDate, ExpectedAt: o.ExpectedAt,
        ConfirmedAt: o.ConfirmedAt, ShippedAt: o.ShippedAt,
        ReceivedAt: o.ReceivedAt,
        Subtotal: o.Subtotal, DiscountAmount: o.DiscountAmount,
        TaxAmount: o.TaxAmount, ShippingFee: o.ShippingFee,
        TotalAmount: o.TotalAmount, Currency: o.Currency,
        POReference: o.POReference, Notes: o.Notes,
        CreatedAt: o.CreatedAt,
    }
    if o.CustomerID != nil { r.CustomerID = o.CustomerID.String() }
    if o.SupplierID != nil { r.SupplierID = o.SupplierID.String() }
    for _, l := range o.Lines {
        r.Lines = append(r.Lines, OrderLineResponse{
            ID: l.ID.String(), LineNo: l.LineNo,
            ProductID: l.ProductID.String(),
            SKU: l.SKU, Name: l.Name,
            Quantity: l.Quantity, UOM: l.UOM,
            UnitPrice: l.UnitPrice, Discount: l.Discount,
            Subtotal: l.Subtotal, TaxAmount: l.TaxAmount, LineTotal: l.LineTotal,
            QtyShipped: l.QtyShipped, QtyReceived: l.QtyReceived, QtyInvoiced: l.QtyInvoiced,
        })
    }
    return r
}

func toInvoiceResponse(i *entity.Invoice) *InvoiceResponse {
    if i == nil { return nil }
    r := &InvoiceResponse{
        ID: i.ID.String(), InvoiceNo: i.InvoiceNo.String(),
        Type: string(i.Type), Status: string(i.Status),
        IssueDate: i.IssueDate, DueDate: i.DueDate, PaidAt: i.PaidAt,
        Subtotal: i.Subtotal, TaxAmount: i.TaxAmount,
        TotalAmount: i.TotalAmount, AmountPaid: i.AmountPaid,
        AmountDue: i.AmountDue, Currency: i.Currency,
        DaysOverdue: i.DaysOverdue(), AgingBucket: i.AgingBucket(),
        IsOverdue: i.IsOverdue(),
        CreatedAt: i.CreatedAt,
    }
    if i.CustomerID != nil { r.CustomerID = i.CustomerID.String() }
    if i.SupplierID != nil { r.SupplierID = i.SupplierID.String() }
    if i.OrderID != nil { r.OrderID = i.OrderID.String() }
    for _, l := range i.Lines {
        r.Lines = append(r.Lines, InvoiceLineResponse{
            ID: l.ID.String(), LineNo: l.LineNo,
            ProductID: l.ProductID.String(), SKU: l.SKU, Name: l.Name,
            Quantity: l.Quantity, UOM: l.UOM, UnitPrice: l.UnitPrice,
            Subtotal: l.Subtotal, TaxAmount: l.TaxAmount, LineTotal: l.LineTotal,
        })
    }
    return r
}

func toPaymentResponse(p *entity.Payment) *PaymentResponse {
    if p == nil { return nil }
    r := &PaymentResponse{
        ID: p.ID.String(), PaymentNo: p.PaymentNo.String(),
        Direction: p.Direction, Method: string(p.Method), Status: string(p.Status),
        Amount: p.Amount, Allocated: p.Allocated, Unallocated: p.Unallocated,
        Currency: p.Currency, PaymentDate: p.PaymentDate,
        CreatedAt: p.CreatedAt,
    }
    if p.CustomerID != nil { r.CustomerID = p.CustomerID.String() }
    if p.SupplierID != nil { r.SupplierID = p.SupplierID.String() }
    for _, a := range p.Allocations {
        r.Allocations = append(r.Allocations, PaymentAllocationResponse{
            InvoiceID: a.InvoiceID, Amount: a.Amount,
        })
    }
    return r
}

func dtoToAddress(d AddressDTO) valueobject.Address {
    return valueobject.Address{
        Line1: d.Line1, Line2: d.Line2, District: d.District,
        Amphoe: d.Amphoe, Province: d.Province, Postcode: d.Postcode,
        Country: d.Country, TaxID: d.TaxID,
    }
}

func addressToDTO(a valueobject.Address) AddressDTO {
    return AddressDTO{
        Line1: a.Line1, Line2: a.Line2, District: a.District,
        Amphoe: a.Amphoe, Province: a.Province, Postcode: a.Postcode,
        Country: a.Country, TaxID: a.TaxID,
    }
}
```

### B.3 Use Cases

#### `application/create_order.go` ⭐
```go
package application

import (
    "context"
    "time"

    "github.com/google/uuid"

    "icmongolang/internal/modules/erp/domain/entity"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
    "icmongolang/internal/modules/erp/domain/event"
    "icmongolang/internal/modules/erp/domain/repository"
    "icmongolang/internal/modules/erp/domain/service"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
)

type CreateOrderUseCase struct {
    orderRepo    repository.OrderRepository
    productRepo  repository.ProductRepository
    warehouseRepo repository.WarehouseRepository
    inventorySvc *service.InventoryService
    pricingSvc   *service.PricingService
    numSvc       *service.DocumentNumberService
    auditRepo    repository.AuditRepository
    producer     kafka.Producer
    log          logger.Logger
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, in CreateOrderInput) (*OrderResponse, error) {
    tenantID, err := uuid.Parse(in.TenantID)
    if err != nil { return nil, domainerrors.ErrInvalidTenantID }
    actorID, _ := uuid.Parse(in.ActorID)

    otype := valueobject.OrderType(in.OrderType)
    if !otype.IsValid() { return nil, domainerrors.ErrInvalidOrderType }

    // 1. Generate order number
    orderNo, err := uc.numSvc.Next(ctx, tenantID, otype.DocTypeFor())
    if err != nil { return nil, err }

    // 2. Create order
    orderDate := in.OrderDate
    if orderDate.IsZero() { orderDate = time.Now() }
    currency := in.Currency
    if currency == "" { currency = "THB" }

    order, err := entity.NewOrder(tenantID, orderNo, otype, orderDate, currency, actorID)
    if err != nil { return nil, err }

    // 3. Set customer/supplier
    if otype.IsSales() {
        if in.CustomerID == "" { return nil, domainerrors.ErrCustomerRequired }
        custID, _ := uuid.Parse(in.CustomerID)
        var addr valueobject.Address
        if in.ShipTo != nil { addr = dtoToAddress(*in.ShipTo) }
        _ = order.SetCustomer(custID, addr)
    } else {
        if in.SupplierID == "" { return nil, domainerrors.ErrSupplierRequired }
        supID, _ := uuid.Parse(in.SupplierID)
        var addr valueobject.Address
        if in.ShipTo != nil { addr = dtoToAddress(*in.ShipTo) }
        _ = order.SetSupplier(supID, addr)
    }

    order.ExpectedAt = in.ExpectedAt
    order.POReference = in.POReference
    order.Notes = in.Notes

    if in.BillTo != nil {
        order.BillToAddress = dtoToAddress(*in.BillTo)
    }

    // 4. Add lines
    for _, ld := range in.Lines {
        pid, _ := uuid.Parse(ld.ProductID)
        p, err := uc.productRepo.FindByID(ctx, tenantID, pid)
        if err != nil { return nil, err }

        // Resolve price
        unitPrice := ld.UnitPrice
        if unitPrice == 0 {
            unitPrice = p.SellPrice.Amount
            if otype.IsPurchase() {
                unitPrice = p.CostPrice.Amount
            }
        }

        line, err := entity.NewOrderLine(
            order.ID, 0, pid, p.SKU.String(), p.Name, p.UOM,
            ld.Quantity, unitPrice, p.TaxType, p.TaxRate,
        )
        if err != nil { return nil, err }

        if ld.Discount > 0 { _ = line.SetDiscount(ld.Discount) }
        if ld.DiscountPct > 0 { _ = line.SetDiscountPercent(ld.DiscountPct) }

        if ld.WarehouseID != "" {
            whID, _ := uuid.Parse(ld.WarehouseID)
            line.WarehouseID = &whID
        }

        if err := order.AddLine(line); err != nil { return nil, err }
    }

    // 5. Persist
    if err := uc.orderRepo.Save(ctx, order); err != nil {
        uc.log.Error("save order failed", "err", err)
        return nil, domainerrors.ErrPersistenceFailure
    }

    // 6. Audit
    _ = uc.auditRepo.Save(ctx, &repository.AuditEntry{
        TenantID: tenantID, ActorID: actorID,
        Action: "order.create", EntityType: "order", EntityID: order.ID,
        Payload: map[string]any{
            "order_no": order.OrderNo.String(),
            "type":     string(otype),
            "total":    order.TotalAmount,
        },
    })

    // 7. Publish event
    _ = uc.producer.Publish(ctx, event.TopicOrderCreated, order.ID.String(), event.OrderCreated{
        EventID: uuid.New(), OrderID: order.ID, TenantID: tenantID,
        OrderNo: order.OrderNo.String(), OrderType: string(otype),
        CustomerID: order.CustomerID, SupplierID: order.SupplierID,
        TotalAmount: order.TotalAmount, Currency: order.Currency,
        OccurredAt: time.Now(),
    })

    return toOrderResponse(order, nil), nil
}

func NewCreateOrderUseCase(
    orderRepo repository.OrderRepository,
    productRepo repository.ProductRepository,
    warehouseRepo repository.WarehouseRepository,
    inventorySvc *service.InventoryService,
    pricingSvc *service.PricingService,
    numSvc *service.DocumentNumberService,
    auditRepo repository.AuditRepository,
    producer kafka.Producer,
    log logger.Logger,
) *CreateOrderUseCase {
    return &CreateOrderUseCase{
        orderRepo: orderRepo, productRepo: productRepo,
        warehouseRepo: warehouseRepo, inventorySvc: inventorySvc,
        pricingSvc: pricingSvc, numSvc: numSvc,
        auditRepo: auditRepo, producer: producer, log: log,
    }
}
```

#### `application/ship_order.go` ⭐
```go
package application

import (
    "context"
    "time"

    "github.com/google/uuid"

    domainerrors "icmongolang/internal/modules/erp/domain/errors"
    "icmongolang/internal/modules/erp/domain/event"
    "icmongolang/internal/modules/erp/domain/repository"
    "icmongolang/internal/modules/erp/domain/service"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
    "icmongolang/pkg/transaction"
)

// ShipOrderUseCase – ตัดสต็อก + mark shipped + publish event
type ShipOrderUseCase struct {
    orderRepo    repository.OrderRepository
    inventorySvc *service.InventoryService
    journalSvc   *service.JournalService
    auditRepo    repository.AuditRepository
    producer     kafka.Producer
    tx           *transaction.Manager
    log          logger.Logger
}

func (uc *ShipOrderUseCase) Execute(ctx context.Context, in ShipOrderInput) (*OrderResponse, error) {
    tenantID, _ := uuid.Parse(in.TenantID)
    orderID, _ := uuid.Parse(in.OrderID)
    actorID, _ := uuid.Parse(in.ActorID)

    var result *OrderResponse
    var shippedLines []event.ShippedLine

    err := uc.tx.Do(ctx, func(txCtx context.Context) error {
        order, err := uc.orderRepo.FindByID(txCtx, tenantID, orderID)
        if err != nil { return err }

        // Convert lineID-string map to uuid map
        lineQtys := map[uuid.UUID]float64{}
        for k, v := range in.Lines {
            lid, _ := uuid.Parse(k)
            lineQtys[lid] = v
        }

        // 1. Ship (validates qty)
        if err := order.Ship(lineQtys); err != nil { return err }

        // 2. Apply stock movements per line
        totalCOGS := 0.0
        for _, l := range order.Lines {
            qty := lineQtys[l.ID]
            if qty <= 0 { continue }
            whID := uuid.Nil
            if l.WarehouseID != nil { whID = *l.WarehouseID }

            item, mv, err := uc.inventorySvc.ApplyMovement(
                txCtx, tenantID, l.ProductID, whID,
                valueobjectMovementTypeIssue(),
                qty, 0, // unit cost — will pull avg_cost
                "ORDER", &order.ID, actorID,
            )
            if err != nil { return err }

            // COGS
            cogs := mv.Quantity * item.AvgCost
            totalCOGS += cogs

            shippedLines = append(shippedLines, event.ShippedLine{
                LineID: l.ID, ProductID: l.ProductID,
                WarehouseID: whID, Quantity: qty,
            })
        }

        // 3. Post COGS journal
        if totalCOGS > 0 {
            _ = uc.journalSvc.PostCOGS(txCtx, tenantID, order.OrderNo.String(),
                totalCOGS, order.Currency, actorID)
        }

        // 4. Persist
        if err := uc.orderRepo.Save(txCtx, order); err != nil { return err }

        result = toOrderResponse(order, nil)
        return nil
    })
    if err != nil { return nil, err }

    // Side effects
    _ = uc.auditRepo.Save(ctx, &repository.AuditEntry{
        TenantID: tenantID, ActorID: actorID,
        Action: "order.ship", EntityType: "order", EntityID: orderID,
    })

    _ = uc.producer.Publish(ctx, event.TopicOrderShipped, orderID.String(), event.OrderShipped{
        EventID: uuid.New(), OrderID: orderID, TenantID: tenantID,
        OrderNo: result.OrderNo,
        CustomerID: uuid.MustParse(result.CustomerID),
        ShippedLines: shippedLines,
        OccurredAt: time.Now(),
    })

    return result, nil
}

func valueobjectMovementTypeIssue() valueobject.MovementType { return valueobject.MovementTypeIssue }

func NewShipOrderUseCase(
    orderRepo repository.OrderRepository,
    inventorySvc *service.InventoryService,
    journalSvc *service.JournalService,
    auditRepo repository.AuditRepository,
    producer kafka.Producer,
    tx *transaction.Manager,
    log logger.Logger,
) *ShipOrderUseCase {
    return &ShipOrderUseCase{
        orderRepo: orderRepo, inventorySvc: inventorySvc,
        journalSvc: journalSvc, auditRepo: auditRepo,
        producer: producer, tx: tx, log: log,
    }
}

var _ = domainerrors.ErrInsufficientStock
```

#### `application/create_invoice.go`
```go
package application

import (
    "context"
    "time"

    "github.com/google/uuid"

    "icmongolang/internal/modules/erp/domain/entity"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
    "icmongolang/internal/modules/erp/domain/event"
    "icmongolang/internal/modules/erp/domain/repository"
    "icmongolang/internal/modules/erp/domain/service"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
)

type CreateInvoiceUseCase struct {
    invoiceRepo repository.InvoiceRepository
    orderRepo   repository.OrderRepository
    productRepo repository.ProductRepository
    journalSvc  *service.JournalService
    numSvc      *service.DocumentNumberService
    auditRepo   repository.AuditRepository
    producer    kafka.Producer
    log         logger.Logger
}

func (uc *CreateInvoiceUseCase) Execute(ctx context.Context, in CreateInvoiceInput) (*InvoiceResponse, error) {
    tenantID, _ := uuid.Parse(in.TenantID)
    actorID, _ := uuid.Parse(in.ActorID)

    itype := valueobject.InvoiceType(in.InvoiceType)
    if !itype.IsValid() { return nil, domainerrors.ErrInvalidInvoiceType }

    no, err := uc.numSvc.Next(ctx, tenantID, itype.DocTypeFor())
    if err != nil { return nil, err }

    issueDate := in.IssueDate
    if issueDate.IsZero() { issueDate = time.Now() }
    dueDate := in.DueDate
    if dueDate.IsZero() { dueDate = issueDate.AddDate(0, 0, 30) }

    currency := in.Currency
    if currency == "" { currency = "THB" }

    inv, err := entity.NewInvoice(tenantID, no, itype, issueDate, dueDate, currency, actorID)
    if err != nil { return nil, err }

    // Set customer/supplier
    if itype == valueobject.InvoiceTypeAR && in.CustomerID != "" {
        cid, _ := uuid.Parse(in.CustomerID)
        inv.CustomerID = &cid
    }
    if itype == valueobject.InvoiceTypeAP && in.SupplierID != "" {
        sid, _ := uuid.Parse(in.SupplierID)
        inv.SupplierID = &sid
    }
    if in.OrderID != "" {
        oid, _ := uuid.Parse(in.OrderID)
        inv.OrderID = &oid
    }

    inv.PaymentTerms = in.PaymentTerms
    inv.Notes = in.Notes
    if in.BillTo != nil {
        inv.BillingAddress = dtoToAddress(*in.BillTo)
    }

    // Add lines
    for _, ld := range in.Lines {
        pid, _ := uuid.Parse(ld.ProductID)
        p, err := uc.productRepo.FindByID(ctx, tenantID, pid)
        if err != nil { return nil, err }

        line := &entity.InvoiceLine{
            ID: uuid.New(),
            ProductID: pid, SKU: p.SKU.String(), Name: p.Name,
            Quantity: ld.Quantity, UOM: p.UOM,
            UnitPrice: ld.UnitPrice, Discount: ld.Discount,
            TaxType: p.TaxType, TaxRate: p.TaxRate,
        }
        if ld.OrderLineID != "" {
            olid, _ := uuid.Parse(ld.OrderLineID)
            line.OrderLineID = &olid
        }
        if err := inv.AddLine(line); err != nil { return nil, err }
    }

    // Issue immediately
    if err := inv.Issue(actorID); err != nil { return nil, err }

    // Persist
    if err := uc.invoiceRepo.Save(ctx, inv); err != nil {
        uc.log.Error("save invoice failed", "err", err)
        return nil, domainerrors.ErrPersistenceFailure
    }

    // Post journal
    _ = uc.journalSvc.PostARInvoice(ctx, tenantID, inv.ID, inv.InvoiceNo.String(),
        inv.Subtotal, inv.TaxAmount, inv.TotalAmount, inv.Currency, actorID)

    // Update order status
    if inv.OrderID != nil {
        if order, err := uc.orderRepo.FindByID(ctx, tenantID, *inv.OrderID); err == nil {
            lineQtys := map[uuid.UUID]float64{}
            for _, l := range inv.Lines {
                if l.OrderLineID != nil {
                    lineQtys[*l.OrderLineID] = l.Quantity
                }
            }
            _ = order.Invoice(lineQtys)
            _ = uc.orderRepo.Save(ctx, order)
        }
    }

    // Audit + event
    _ = uc.auditRepo.Save(ctx, &repository.AuditEntry{
        TenantID: tenantID, ActorID: actorID,
        Action: "invoice.create", EntityType: "invoice", EntityID: inv.ID,
    })

    _ = uc.producer.Publish(ctx, event.TopicInvoiceIssued, inv.ID.String(), event.InvoiceIssued{
        EventID: uuid.New(), InvoiceID: inv.ID, TenantID: tenantID,
        InvoiceNo: inv.InvoiceNo.String(), InvoiceType: string(itype),
        CustomerID: inv.CustomerID, SupplierID: inv.SupplierID, OrderID: inv.OrderID,
        Subtotal: inv.Subtotal, TaxAmount: inv.TaxAmount, TotalAmount: inv.TotalAmount,
        Currency: inv.Currency,
        IssueDate: inv.IssueDate, DueDate: inv.DueDate,
        OccurredAt: time.Now(),
    })

    return toInvoiceResponse(inv), nil
}

func NewCreateInvoiceUseCase(
    invoiceRepo repository.InvoiceRepository,
    orderRepo repository.OrderRepository,
    productRepo repository.ProductRepository,
    journalSvc *service.JournalService,
    numSvc *service.DocumentNumberService,
    auditRepo repository.AuditRepository,
    producer kafka.Producer,
    log logger.Logger,
) *CreateInvoiceUseCase {
    return &CreateInvoiceUseCase{
        invoiceRepo: invoiceRepo, orderRepo: orderRepo,
        productRepo: productRepo, journalSvc: journalSvc,
        numSvc: numSvc, auditRepo: auditRepo,
        producer: producer, log: log,
    }
}
```

#### `application/record_payment.go` ⭐
```go
package application

import (
    "context"
    "time"

    "github.com/google/uuid"

    "icmongolang/internal/modules/erp/domain/entity"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
    "icmongolang/internal/modules/erp/domain/event"
    "icmongolang/internal/modules/erp/domain/repository"
    "icmongolang/internal/modules/erp/domain/service"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
    "icmongolang/pkg/transaction"
)

type RecordPaymentUseCase struct {
    paymentRepo repository.PaymentRepository
    invoiceRepo repository.InvoiceRepository
    journalSvc  *service.JournalService
    numSvc      *service.DocumentNumberService
    auditRepo   repository.AuditRepository
    producer    kafka.Producer
    tx          *transaction.Manager
    log         logger.Logger
}

func (uc *RecordPaymentUseCase) Execute(ctx context.Context, in CreatePaymentInput) (*PaymentResponse, error) {
    tenantID, _ := uuid.Parse(in.TenantID)
    actorID, _ := uuid.Parse(in.ActorID)

    // 1. Generate number
    no, err := uc.numSvc.Next(ctx, tenantID, valueobject.DocTypePayment)
    if err != nil { return nil, err }

    // 2. Create payment
    payDate := in.PaymentDate
    if payDate.IsZero() { payDate = time.Now() }
    currency := in.Currency
    if currency == "" { currency = "THB" }

    payment, err := entity.NewPayment(
        tenantID, no, in.Direction,
        valueobject.PaymentMethod(in.Method),
        in.Amount, currency, payDate, actorID,
    )
    if err != nil { return nil, err }

    if in.CustomerID != "" {
        cid, _ := uuid.Parse(in.CustomerID)
        _ = payment.SetCustomer(cid)
    }
    if in.SupplierID != "" {
        sid, _ := uuid.Parse(in.SupplierID)
        _ = payment.SetSupplier(sid)
    }
    payment.Reference = in.Reference
    payment.BankAccount = in.BankAccount
    payment.Notes = in.Notes

    // 3. Confirm
    if err := payment.Confirm(); err != nil { return nil, err }

    // 4. Persist + allocate in transaction
    var updatedInvoices []*entity.Invoice
    err = uc.tx.Do(ctx, func(txCtx context.Context) error {
        if err := uc.paymentRepo.Save(txCtx, payment); err != nil { return err }

        // 5. Apply allocations
        for _, a := range in.Allocations {
            invID, _ := uuid.Parse(a.InvoiceID)

            inv, err := uc.invoiceRepo.FindByID(txCtx, tenantID, invID)
            if err != nil { return err }
            if !inv.Status.IsPayable() {
                return domainerrors.ErrInvoiceNotPayable
            }

            if err := payment.Allocate(invID, a.Amount); err != nil {
                return err
            }
            if err := inv.ApplyPayment(a.Amount); err != nil {
                return err
            }

            if err := uc.invoiceRepo.Save(txCtx, inv); err != nil {
                return err
            }
            updatedInvoices = append(updatedInvoices, inv)
        }

        return uc.paymentRepo.Save(txCtx, payment)
    })
    if err != nil { return nil, err }

    // 6. Post journal
    _ = uc.journalSvc.PostPaymentIn(ctx, tenantID, payment.ID,
        payment.PaymentNo.String(), payment.Amount, payment.Currency, actorID)

    // 7. Events
    _ = uc.producer.Publish(ctx, event.TopicPaymentConfirmed, payment.ID.String(), event.PaymentConfirmed{
        EventID: uuid.New(), PaymentID: payment.ID, TenantID: tenantID,
        PaymentNo: payment.PaymentNo.String(), Direction: payment.Direction,
        Amount: payment.Amount, Currency: payment.Currency,
        OccurredAt: time.Now(),
    })

    for _, inv := range updatedInvoices {
        if inv.IsPaid() {
            _ = uc.producer.Publish(ctx, event.TopicInvoicePaid, inv.ID.String(), event.InvoicePaid{
                EventID: uuid.New(), InvoiceID: inv.ID, TenantID: tenantID,
                InvoiceNo: inv.InvoiceNo.String(), AmountPaid: inv.AmountPaid,
                Currency: inv.Currency, OccurredAt: time.Now(),
            })
        }
    }

    return toPaymentResponse(payment), nil
}

func NewRecordPaymentUseCase(
    paymentRepo repository.PaymentRepository,
    invoiceRepo repository.InvoiceRepository,
    journalSvc *service.JournalService,
    numSvc *service.DocumentNumberService,
    auditRepo repository.AuditRepository,
    producer kafka.Producer,
    tx *transaction.Manager,
    log logger.Logger,
) *RecordPaymentUseCase {
    return &RecordPaymentUseCase{
        paymentRepo: paymentRepo, invoiceRepo: invoiceRepo,
        journalSvc: journalSvc, numSvc: numSvc,
        auditRepo: auditRepo, producer: producer, tx: tx, log: log,
    }
}
```

#### `application/adjust_inventory.go`
```go
package application

import (
    "context"
    "time"

    "github.com/google/uuid"

    domainerrors "icmongolang/internal/modules/erp/domain/errors"
    "icmongolang/internal/modules/erp/domain/event"
    "icmongolang/internal/modules/erp/domain/repository"
    "icmongolang/internal/modules/erp/domain/service"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
)

type AdjustInventoryUseCase struct {
    invRepo    repository.InventoryRepository
    prodRepo   repository.ProductRepository
    invSvc     *service.InventoryService
    auditRepo  repository.AuditRepository
    producer   kafka.Producer
    log        logger.Logger
}

func (uc *AdjustInventoryUseCase) Execute(ctx context.Context, in AdjustInventoryInput) (*InventoryResponse, error) {
    tenantID, _ := uuid.Parse(in.TenantID)
    productID, _ := uuid.Parse(in.ProductID)
    warehouseID, _ := uuid.Parse(in.WarehouseID)
    actorID, _ := uuid.Parse(in.ActorID)

    // 1. Load current
    item, err := uc.invRepo.FindByKey(ctx, tenantID, productID, warehouseID)
    if err != nil {
        // create new
        item = entity.NewInventoryItem(tenantID, productID, warehouseID)
    }
    oldQty := item.QtyOnHand

    // 2. Calculate delta
    delta := in.NewQty - oldQty
    if delta == 0 {
        return toInventoryResponse(item, nil, nil), nil
    }

    mvType := valueobject.MovementTypeAdjustIn
    if delta < 0 { mvType = valueobject.MovementTypeAdjustOut }

    // 3. Apply
    _, mv, err := uc.invSvc.ApplyMovement(
        ctx, tenantID, productID, warehouseID,
        mvType, abs(delta), item.AvgCost,
        "ADJUSTMENT", nil, actorID,
    )
    if err != nil { return nil, err }
    mv.Reason = in.Reason

    // 4. Publish event
    _ = uc.producer.Publish(ctx, event.TopicStockAdjusted, productID.String(), event.StockAdjusted{
        EventID: uuid.New(), TenantID: tenantID,
        ProductID: productID, WarehouseID: warehouseID,
        OldQty: oldQty, NewQty: in.NewQty, Delta: delta,
        Reason: in.Reason, ActorID: actorID,
        OccurredAt: time.Now(),
    })

    // Reload
    item2, _ := uc.invRepo.FindByKey(ctx, tenantID, productID, warehouseID)
    return toInventoryResponse(item2, nil, nil), nil
}

func abs(v float64) float64 { if v < 0 { return -v }; return v }

func NewAdjustInventoryUseCase(
    invRepo repository.InventoryRepository,
    prodRepo repository.ProductRepository,
    invSvc *service.InventoryService,
    auditRepo repository.AuditRepository,
    producer kafka.Producer,
    log logger.Logger,
) *AdjustInventoryUseCase {
    return &AdjustInventoryUseCase{
        invRepo: invRepo, prodRepo: prodRepo,
        invSvc: invSvc, auditRepo: auditRepo,
        producer: producer, log: log,
    }
}

var _ = domainerrors.ErrInventoryNotFound
```

#### `application/transfer_stock.go`
```go
package application

import (
    "context"
    "time"

    "github.com/google/uuid"

    "icmongolang/internal/modules/erp/domain/entity"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
    "icmongolang/internal/modules/erp/domain/repository"
    "icmongolang/internal/modules/erp/domain/service"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
    "icmongolang/pkg/transaction"
)

// TransferStockUseCase – ย้ายสินค้าระหว่าง warehouse (2 movements)
type TransferStockUseCase struct {
    invRepo    repository.InventoryRepository
    invSvc     *service.InventoryService
    producer   kafka.Producer
    tx         *transaction.Manager
    log        logger.Logger
}

func (uc *TransferStockUseCase) Execute(ctx context.Context, in TransferStockInput) error {
    tenantID, _ := uuid.Parse(in.TenantID)
    productID, _ := uuid.Parse(in.ProductID)
    fromWH, _ := uuid.Parse(in.FromWarehouseID)
    toWH, _ := uuid.Parse(in.ToWarehouseID)
    actorID, _ := uuid.Parse(in.ActorID)

    if fromWH == toWH {
        return domainerrors.ErrInvalidWarehouse
    }

    return uc.tx.Do(ctx, func(txCtx context.Context) error {
        // 1. Load source item
        item, err := uc.invRepo.FindByKey(txCtx, tenantID, productID, fromWH)
        if err != nil { return err }
        if item.QtyAvailable() < in.Quantity {
            return domainerrors.ErrInsufficientStock
        }

        unitCost := item.AvgCost

        // 2. OUT from source
        _, _, err = uc.invSvc.ApplyMovement(
            txCtx, tenantID, productID, fromWH,
            valueobject.MovementTypeTransferOut,
            in.Quantity, unitCost,
            "TRANSFER", nil, actorID,
        )
        if err != nil { return err }

        // 3. IN to destination
        _, _, err = uc.invSvc.ApplyMovement(
            txCtx, tenantID, productID, toWH,
            valueobject.MovementTypeTransferIn,
            in.Quantity, unitCost,
            "TRANSFER", nil, actorID,
        )
        return err
    })
}

func NewTransferStockUseCase(
    invRepo repository.InventoryRepository,
    invSvc *service.InventoryService,
    producer kafka.Producer,
    tx *transaction.Manager,
    log logger.Logger,
) *TransferStockUseCase {
    return &TransferStockUseCase{
        invRepo: invRepo, invSvc: invSvc,
        producer: producer, tx: tx, log: log,
    }
}

var _ = entity.NewInventoryItem
var _ = time.Now
```

> **หมายเหตุ**: use cases อื่น (`UpdateProduct`, `ListProducts`, `GetOrder`, `CancelOrder`, `IssueInvoice`, `VoidInvoice`, `CreateCreditNote`, `GetInventory`, `ListStockMovements`, `ReserveStock`, `ListInvoices`, `GetAgingReport`) ทำตามรูปแบบเดียวกัน — pattern เดียวกัน

---

## 🅲 PART 4C — INFRASTRUCTURE LAYER

### C.1 GORM Models

```go
package postgres

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/datatypes"
)

// Product
type ProductModel struct {
    ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID        uuid.UUID      `gorm:"type:uuid;not null;index:idx_erp_prod_tenant,priority:1"`
    SKU             string         `gorm:"size:50;not null;index:idx_erp_prod_tenant_sku,unique"`
    Name            string         `gorm:"size:255;not null"`
    Description     string         `gorm:"type:text"`
    Type            string         `gorm:"size:20;not null;index"`
    Category        string         `gorm:"size:100;index"`
    Brand           string         `gorm:"size:100"`
    UOM             string         `gorm:"size:20;not null"`
    CostPrice       float64        `gorm:"type:numeric(15,2);not null;default:0"`
    SellPrice       float64        `gorm:"type:numeric(15,2);not null;default:0"`
    Currency        string         `gorm:"size:3;not null;default:'THB'"`
    TaxType         string         `gorm:"size:20;default:'VAT'"`
    TaxRate         float64        `gorm:"type:numeric(5,4);default:0.07"`
    TrackInventory  bool           `gorm:"not null;default:false"`
    ReorderPoint    float64        `gorm:"type:numeric(15,4);default:0"`
    ReorderQty      float64        `gorm:"type:numeric(15,4);default:0"`
    LeadTimeDays    int            `gorm:"default:0"`
    SafetyStock     float64        `gorm:"type:numeric(15,4);default:0"`
    Weight          float64        `gorm:"type:numeric(10,3);default:0"`
    Volume          float64        `gorm:"type:numeric(10,4);default:0"`
    Barcode         string         `gorm:"size:100;index"`
    ImageURLs       datatypes.JSON `gorm:"type:jsonb"`
    IsActive        bool           `gorm:"not null;default:true;index:idx_erp_prod_tenant,priority:2"`
    IsSellable      bool           `gorm:"default:true"`
    IsPurchasable   bool           `gorm:"default:true"`
    Metadata        datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
    Tags            datatypes.JSON `gorm:"type:jsonb;default:'[]'"`
    CreatedBy       uuid.UUID      `gorm:"type:uuid"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

func (ProductModel) TableName() string { return "erp_products" }

// Warehouse
type WarehouseModel struct {
    ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID   uuid.UUID      `gorm:"type:uuid;not null;index:idx_erp_wh_tenant_code,unique"`
    Code       string         `gorm:"size:30;not null"`
    Name       string         `gorm:"size:255;not null"`
    Type       string         `gorm:"size:20;not null"`
    Address    datatypes.JSON `gorm:"type:jsonb"`
    ManagerID  *uuid.UUID     `gorm:"type:uuid"`
    Phone      string         `gorm:"size:30"`
    Email      string         `gorm:"size:255"`
    IsActive   bool           `gorm:"not null;default:true"`
    IsDefault  bool           `gorm:"default:false;index"`
    Metadata   datatypes.JSON `gorm:"type:jsonb"`
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

func (WarehouseModel) TableName() string { return "erp_warehouses" }

// InventoryItem
type InventoryItemModel struct {
    ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID      uuid.UUID `gorm:"type:uuid;not null;index:idx_erp_inv_tenant_product,unique:idx_erp_inv_unique,priority:1"`
    ProductID     uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_erp_inv_unique,priority:2"`
    WarehouseID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_erp_inv_unique,priority:3"`
    QtyOnHand     float64   `gorm:"type:numeric(15,4);not null;default:0"`
    QtyReserved   float64   `gorm:"type:numeric(15,4);not null;default:0"`
    AvgCost       float64   `gorm:"type:numeric(15,4);not null;default:0"`
    LastCost      float64   `gorm:"type:numeric(15,4);not null;default:0"`
    BinLocation   string    `gorm:"size:50"`
    UpdatedAt     time.Time
    CreatedAt     time.Time
}

func (InventoryItemModel) TableName() string { return "erp_inventory_items" }

// StockMovement
type StockMovementModel struct {
    ID                int64          `gorm:"primaryKey;autoIncrement"`
    TenantID          uuid.UUID      `gorm:"type:uuid;not null;index:idx_erp_mv_tenant_product,priority:1"`
    ProductID         uuid.UUID      `gorm:"type:uuid;not null;index:idx_erp_mv_tenant_product,priority:2"`
    WarehouseID       uuid.UUID      `gorm:"type:uuid;not null;index"`
    Type              string         `gorm:"size:20;not null;index"`
    Direction         int            `gorm:"not null"`
    Quantity          float64        `gorm:"type:numeric(15,4);not null"`
    UOM               string         `gorm:"size:20"`
    UnitCost          float64        `gorm:"type:numeric(15,4)"`
    TotalCost         float64        `gorm:"type:numeric(15,2)"`
    QtyOnHandAfter    float64        `gorm:"type:numeric(15,4)"`
    QtyReservedAfter  float64        `gorm:"type:numeric(15,4)"`
    RefType           string         `gorm:"size:30;index"`
    RefID             *uuid.UUID     `gorm:"type:uuid"`
    RefNumber         string         `gorm:"size:50"`
    Reason            string         `gorm:"type:text"`
    Notes             string         `gorm:"type:text"`
    CounterWarehouseID *uuid.UUID    `gorm:"type:uuid"`
    OccurredAt        time.Time      `gorm:"index:idx_erp_mv_tenant_product,priority:3,sort:desc"`
    CreatedBy         uuid.UUID      `gorm:"type:uuid"`
}

func (StockMovementModel) TableName() string { return "erp_stock_movements" }

// StockReservation
type StockReservationModel struct {
    ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID    uuid.UUID `gorm:"type:uuid;not null;index"`
    ProductID   uuid.UUID `gorm:"type:uuid;not null;index"`
    WarehouseID uuid.UUID `gorm:"type:uuid;not null;index"`
    OrderID     uuid.UUID `gorm:"type:uuid;not null;index"`
    OrderLineID uuid.UUID `gorm:"type:uuid;not null;index"`
    Quantity    float64   `gorm:"type:numeric(15,4);not null"`
    Status      string    `gorm:"size:20;not null;index"`
    ExpiresAt   *time.Time
    ConsumedAt  *time.Time
    ReleasedAt  *time.Time
    CreatedAt   time.Time
}

func (StockReservationModel) TableName() string { return "erp_stock_reservations" }

// Order
type OrderModel struct {
    ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID        uuid.UUID      `gorm:"type:uuid;not null;index:idx_erp_order_tenant_status,priority:1"`
    OrderNo         string         `gorm:"size:30;not null;index:idx_erp_order_tenant_no,unique"`
    Type            string         `gorm:"size:20;not null;index"`
    CustomerID      *uuid.UUID     `gorm:"type:uuid;index"`
    SupplierID      *uuid.UUID     `gorm:"type:uuid;index"`
    Status          string         `gorm:"size:30;not null;default:'DRAFT';index:idx_erp_order_tenant_status,priority:2"`
    OrderDate       time.Time      `gorm:"not null;index"`
    ExpectedAt      *time.Time
    ConfirmedAt     *time.Time
    ShippedAt       *time.Time
    ReceivedAt      *time.Time
    ClosedAt        *time.Time
    CancelledAt     *time.Time
    ShipToAddress   datatypes.JSON `gorm:"type:jsonb"`
    BillToAddress   datatypes.JSON `gorm:"type:jsonb"`
    ShippingMethod  string         `gorm:"size:50"`
    TrackingNumber  string         `gorm:"size:100"`
    Subtotal        float64        `gorm:"type:numeric(15,2);default:0"`
    DiscountAmount  float64        `gorm:"type:numeric(15,2);default:0"`
    TaxAmount       float64        `gorm:"type:numeric(15,2);default:0"`
    ShippingFee     float64        `gorm:"type:numeric(15,2);default:0"`
    TotalAmount     float64        `gorm:"type:numeric(15,2);default:0"`
    Currency        string         `gorm:"size:3;not null;default:'THB'"`
    FXRate          float64        `gorm:"type:numeric(15,6);default:1"`
    BaseCurrency    string         `gorm:"size:3;default:'THB'"`
    BaseTotal       float64        `gorm:"type:numeric(15,2);default:0"`
    POReference     string         `gorm:"size:100"`
    QuotationID     *uuid.UUID     `gorm:"type:uuid"`
    ContractID      *uuid.UUID     `gorm:"type:uuid"`
    Notes           string         `gorm:"type:text"`
    InternalNotes   string         `gorm:"type:text"`
    Tags            datatypes.JSON `gorm:"type:jsonb"`
    Metadata        datatypes.JSON `gorm:"type:jsonb"`
    CreatedBy       uuid.UUID      `gorm:"type:uuid"`
    ApprovedBy      *uuid.UUID     `gorm:"type:uuid"`
    CreatedAt       time.Time
    UpdatedAt       time.Time

    Lines []OrderLineModel `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

func (OrderModel) TableName() string { return "erp_orders" }

type OrderLineModel struct {
    ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    OrderID     uuid.UUID  `gorm:"type:uuid;not null;index"`
    LineNo      int        `gorm:"not null"`
    ProductID   uuid.UUID  `gorm:"type:uuid;not null;index"`
    SKU         string     `gorm:"size:50;not null"`
    Name        string     `gorm:"size:255;not null"`
    Description string     `gorm:"type:text"`
    Quantity    float64    `gorm:"type:numeric(15,4);not null"`
    UOM         string     `gorm:"size:20;not null"`
    UnitPrice   float64    `gorm:"type:numeric(15,2);not null"`
    Discount    float64    `gorm:"type:numeric(15,2);default:0"`
    DiscountPct float64    `gorm:"type:numeric(5,2);default:0"`
    TaxType     string     `gorm:"size:20"`
    TaxRate     float64    `gorm:"type:numeric(5,4)"`
    Subtotal    float64    `gorm:"type:numeric(15,2);not null"`
    TaxAmount   float64    `gorm:"type:numeric(15,2);default:0"`
    LineTotal   float64    `gorm:"type:numeric(15,2);not null"`
    QtyShipped  float64    `gorm:"type:numeric(15,4);default:0"`
    QtyReceived float64    `gorm:"type:numeric(15,4);default:0"`
    QtyInvoiced float64    `gorm:"type:numeric(15,4);default:0"`
    WarehouseID *uuid.UUID `gorm:"type:uuid"`
    Notes       string     `gorm:"type:text"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func (OrderLineModel) TableName() string { return "erp_order_lines" }

// Invoice
type InvoiceModel struct {
    ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID        uuid.UUID      `gorm:"type:uuid;not null;index:idx_erp_inv_tenant_status,priority:1"`
    InvoiceNo       string         `gorm:"size:30;not null;index:idx_erp_inv_tenant_no,unique"`
    Type            string         `gorm:"size:5;not null;index"`
    CustomerID      *uuid.UUID     `gorm:"type:uuid;index"`
    SupplierID      *uuid.UUID     `gorm:"type:uuid;index"`
    OrderID         *uuid.UUID     `gorm:"type:uuid;index"`
    Status          string         `gorm:"size:20;not null;default:'DRAFT';index:idx_erp_inv_tenant_status,priority:2"`
    IssueDate       time.Time      `gorm:"not null;index"`
    DueDate         time.Time      `gorm:"not null;index"`
    PaidAt          *time.Time
    VoidedAt        *time.Time
    VoidReason      string         `gorm:"type:text"`
    Subtotal        float64        `gorm:"type:numeric(15,2);not null;default:0"`
    DiscountAmount  float64        `gorm:"type:numeric(15,2);default:0"`
    TaxAmount       float64        `gorm:"type:numeric(15,2);default:0"`
    WHTAmount       float64        `gorm:"type:numeric(15,2);default:0"`
    TotalAmount     float64        `gorm:"type:numeric(15,2);not null;default:0"`
    AmountPaid      float64        `gorm:"type:numeric(15,2);default:0"`
    AmountDue       float64        `gorm:"type:numeric(15,2);default:0"`
    Currency        string         `gorm:"size:3;not null;default:'THB'"`
    FXRate          float64        `gorm:"type:numeric(15,6);default:1"`
    BaseCurrency    string         `gorm:"size:3;default:'THB'"`
    BaseTotal       float64        `gorm:"type:numeric(15,2);default:0"`
    PaymentTerms    string         `gorm:"size:50"`
    Notes           string         `gorm:"type:text"`
    BillingAddress  datatypes.JSON `gorm:"type:jsonb"`
    POReference     string         `gorm:"size:100"`
    RefNumber       string         `gorm:"size:100"`
    Metadata        datatypes.JSON `gorm:"type:jsonb"`
    Tags            datatypes.JSON `gorm:"type:jsonb"`
    CreatedBy       uuid.UUID      `gorm:"type:uuid"`
    IssuedBy        *uuid.UUID     `gorm:"type:uuid"`
    CreatedAt       time.Time
    UpdatedAt       time.Time

    Lines []InvoiceLineModel `gorm:"foreignKey:InvoiceID;constraint:OnDelete:CASCADE"`
}

func (InvoiceModel) TableName() string { return "erp_invoices" }

type InvoiceLineModel struct {
    ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    InvoiceID   uuid.UUID  `gorm:"type:uuid;not null;index"`
    LineNo      int        `gorm:"not null"`
    OrderLineID *uuid.UUID `gorm:"type:uuid;index"`
    ProductID   uuid.UUID  `gorm:"type:uuid;not null;index"`
    SKU         string     `gorm:"size:50;not null"`
    Name        string     `gorm:"size:255;not null"`
    Description string     `gorm:"type:text"`
    Quantity    float64    `gorm:"type:numeric(15,4);not null"`
    UOM         string     `gorm:"size:20;not null"`
    UnitPrice   float64    `gorm:"type:numeric(15,2);not null"`
    Discount    float64    `gorm:"type:numeric(15,2);default:0"`
    TaxType     string     `gorm:"size:20"`
    TaxRate     float64    `gorm:"type:numeric(5,4)"`
    Subtotal    float64    `gorm:"type:numeric(15,2);not null"`
    TaxAmount   float64    `gorm:"type:numeric(15,2);default:0"`
    LineTotal   float64    `gorm:"type:numeric(15,2);not null"`
    CreatedAt   time.Time
}

func (InvoiceLineModel) TableName() string { return "erp_invoice_lines" }

// Payment
type PaymentModel struct {
    ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID    uuid.UUID      `gorm:"type:uuid;not null;index"`
    PaymentNo   string         `gorm:"size:30;not null;index:idx_erp_pay_tenant_no,unique"`
    Direction   string         `gorm:"size:5;not null;index"`
    CustomerID  *uuid.UUID     `gorm:"type:uuid;index"`
    SupplierID  *uuid.UUID     `gorm:"type:uuid;index"`
    Method      string         `gorm:"size:30;not null"`
    Status      string         `gorm:"size:20;not null;default:'PENDING';index"`
    Amount      float64        `gorm:"type:numeric(15,2);not null"`
    Allocated   float64        `gorm:"type:numeric(15,2);default:0"`
    Unallocated float64        `gorm:"type:numeric(15,2);not null"`
    Currency    string         `gorm:"size:3;not null;default:'THB'"`
    FXRate      float64        `gorm:"type:numeric(15,6);default:1"`
    BaseAmount  float64        `gorm:"type:numeric(15,2);default:0"`
    PaymentDate time.Time      `gorm:"not null;index"`
    ValueDate   *time.Time
    Reference   string         `gorm:"size:100"`
    BankAccount string         `gorm:"size:50"`
    Notes       string         `gorm:"type:text"`
    Metadata    datatypes.JSON `gorm:"type:jsonb"`
    CreatedBy   uuid.UUID      `gorm:"type:uuid"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    ConfirmedAt *time.Time
    ReversedAt  *time.Time
    RefundedAt  *time.Time

    Allocations []PaymentAllocationModel `gorm:"foreignKey:PaymentID;constraint:OnDelete:CASCADE"`
}

func (PaymentModel) TableName() string { return "erp_payments" }

type PaymentAllocationModel struct {
    ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    PaymentID  uuid.UUID `gorm:"type:uuid;not null;index"`
    InvoiceID  uuid.UUID `gorm:"type:uuid;not null;index"`
    Amount     float64   `gorm:"type:numeric(15,2);not null"`
    Currency   string    `gorm:"size:3;not null"`
    FXRate     float64   `gorm:"type:numeric(15,6);default:1"`
    BaseAmount float64   `gorm:"type:numeric(15,2);default:0"`
    CreatedAt  time.Time
}

func (PaymentAllocationModel) TableName() string { return "erp_payment_allocations" }

// CreditNote
type CreditNoteModel struct {
    ID               uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID         uuid.UUID      `gorm:"type:uuid;not null;index"`
    NoteNo           string         `gorm:"size:30;not null;index:idx_erp_cn_tenant_no,unique"`
    NoteType         string         `gorm:"size:10;not null;index"`
    InvoiceID        *uuid.UUID     `gorm:"type:uuid;index"`
    CustomerID       *uuid.UUID     `gorm:"type:uuid;index"`
    SupplierID       *uuid.UUID     `gorm:"type:uuid;index"`
    Reason           string         `gorm:"size:30;not null"`
    ReasonDetail     string         `gorm:"type:text"`
    IssueDate        time.Time      `gorm:"not null"`
    Currency         string         `gorm:"size:3;not null;default:'THB'"`
    Subtotal         float64        `gorm:"type:numeric(15,2);not null"`
    TaxAmount        float64        `gorm:"type:numeric(15,2);default:0"`
    TotalAmount      float64        `gorm:"type:numeric(15,2);not null"`
    Status           string         `gorm:"size:20;not null;default:'DRAFT'"`
    AppliedToInvoice bool           `gorm:"default:false"`
    AppliedAt        *time.Time
    Refunded         bool           `gorm:"default:false"`
    RefundedAt       *time.Time
    Notes            string         `gorm:"type:text"`
    Metadata         datatypes.JSON `gorm:"type:jsonb"`
    CreatedBy        uuid.UUID      `gorm:"type:uuid"`
    CreatedAt        time.Time
    UpdatedAt        time.Time
}

func (CreditNoteModel) TableName() string { return "erp_credit_notes" }

// Supplier
type SupplierModel struct {
    ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID    uuid.UUID      `gorm:"type:uuid;not null;index:idx_erp_sup_tenant_code,unique"`
    Code        string         `gorm:"size:30;not null"`
    Name        string         `gorm:"size:255;not null"`
    LegalName   string         `gorm:"size:255"`
    TaxID       string         `gorm:"size:20;index"`
    Branch      string         `gorm:"size:20"`
    Email       string         `gorm:"size:255"`
    Phone       string         `gorm:"size:30"`
    ContactName string         `gorm:"size:255"`
    Address     datatypes.JSON `gorm:"type:jsonb"`
    BillingAddr datatypes.JSON `gorm:"type:jsonb"`
    ShippingAddr datatypes.JSON `gorm:"type:jsonb"`
    Status      string         `gorm:"size:20;not null;default:'ACTIVE'"`
    PaymentTerms string        `gorm:"size:50"`
    CreditLimit float64        `gorm:"type:numeric(15,2);default:0"`
    Currency    string         `gorm:"size:3;default:'THB'"`
    OnTimeDeliveryPct float64  `gorm:"type:numeric(5,2);default:0"`
    QualityScore float64       `gorm:"type:numeric(5,2);default:0"`
    AvgLeadTimeDays int
    Notes       string         `gorm:"type:text"`
    Metadata    datatypes.JSON `gorm:"type:jsonb"`
    Tags        datatypes.JSON `gorm:"type:jsonb"`
    IsActive    bool           `gorm:"default:true"`
    CreatedBy   uuid.UUID      `gorm:"type:uuid"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func (SupplierModel) TableName() string { return "erp_suppliers" }

// JournalEntry
type JournalEntryModel struct {
    ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID     uuid.UUID      `gorm:"type:uuid;not null;index"`
    EntryNo      string         `gorm:"size:30;not null;index:idx_erp_je_tenant_no,unique"`
    EntryDate    time.Time      `gorm:"not null;index"`
    Description  string         `gorm:"type:text"`
    TotalDebit   float64        `gorm:"type:numeric(15,2);default:0"`
    TotalCredit  float64        `gorm:"type:numeric(15,2);default:0"`
    Currency     string         `gorm:"size:3;not null;default:'THB'"`
    SourceType   string         `gorm:"size:30;index"`
    SourceID     *uuid.UUID     `gorm:"type:uuid"`
    Status       string         `gorm:"size:20;not null;default:'DRAFT'"`
    PostedAt     *time.Time
    PostedBy     *uuid.UUID     `gorm:"type:uuid"`
    VoidedAt     *time.Time
    VoidReason   string         `gorm:"type:text"`
    CreatedBy    uuid.UUID      `gorm:"type:uuid"`
    CreatedAt    time.Time
    UpdatedAt    time.Time

    Lines []JournalLineModel `gorm:"foreignKey:JournalID;constraint:OnDelete:CASCADE"`
}

func (JournalEntryModel) TableName() string { return "erp_journal_entries" }

type JournalLineModel struct {
    ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    JournalID   uuid.UUID  `gorm:"type:uuid;not null;index"`
    LineNo      int        `gorm:"not null"`
    AccountCode string     `gorm:"size:10;not null;index"`
    AccountName string     `gorm:"size:255"`
    Debit       float64    `gorm:"type:numeric(15,2);default:0"`
    Credit      float64    `gorm:"type:numeric(15,2);default:0"`
    Description string     `gorm:"type:text"`
    RefType     string     `gorm:"size:30"`
    RefID       *uuid.UUID `gorm:"type:uuid"`
}

func (JournalLineModel) TableName() string { return "erp_journal_lines" }

// AuditLog
type AuditLogModel struct {
    ID         int64          `gorm:"primaryKey;autoIncrement"`
    TenantID   uuid.UUID      `gorm:"type:uuid;not null;index"`
    ActorID    uuid.UUID      `gorm:"type:uuid"`
    Action     string         `gorm:"size:50;not null"`
    EntityType string         `gorm:"size:50;not null;index:idx_erp_audit_entity,priority:1"`
    EntityID   uuid.UUID      `gorm:"type:uuid;not null;index:idx_erp_audit_entity,priority:2"`
    Payload    datatypes.JSON `gorm:"type:jsonb"`
    IPAddress  string         `gorm:"size:45"`
    UserAgent  string         `gorm:"size:500"`
    CreatedAt  time.Time      `gorm:"index"`
}

func (AuditLogModel) TableName() string { return "erp_audit_logs" }

// DocumentSequence – running number
type DocumentSequenceModel struct {
    TenantID  uuid.UUID `gorm:"type:uuid;primaryKey"`
    DocType   string    `gorm:"size:10;primaryKey"`
    Year      int       `gorm:"primaryKey"`
    LastSeq   int       `gorm:"not null;default:0"`
    UpdatedAt time.Time
}

func (DocumentSequenceModel) TableName() string { return "erp_document_sequences" }
```

### C.2 Document Number Generator (Postgres)

```go
package postgres

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"

    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type docNumberGen struct{ db *gorm.DB }

func NewDocumentNumberGenerator(db *gorm.DB) *docNumberGen {
    return &docNumberGen{db: db}
}

// NextSeq – atomic increment
func (g *docNumberGen) NextSeq(ctx context.Context, tenantID uuid.UUID, docType valueobject.DocType, year int) (int, error) {
    var seq DocumentSequenceModel
    err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // Upsert + increment
        if err := tx.Clauses(clause.OnConflict{
            Columns: []clause.Column{
                {Name: "tenant_id"}, {Name: "doc_type"}, {Name: "year"},
            },
            DoUpdates: clause.Assignments(map[string]any{
                "last_seq": gorm.Expr("erp_document_sequences.last_seq + 1"),
            }),
        }).Create(&DocumentSequenceModel{
            TenantID: tenantID, DocType: string(docType), Year: year, LastSeq: 1,
        }).Error; err != nil {
            return err
        }

        // Read back
        if err := tx.Where("tenant_id = ? AND doc_type = ? AND year = ?",
            tenantID, docType, year).First(&seq).Error; err != nil {
            return err
        }
        return nil
    })
    if err != nil { return 0, err }
    return seq.LastSeq, nil
}

var _ = fmt.Sprintf
```

### C.3 Repositories (Representative — Product + Order + Inventory)

```go
package postgres

import (
    "context"
    "encoding/json"
    "errors"
    "strings"

    "github.com/google/uuid"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"

    "icmongolang/internal/modules/erp/domain/entity"
    domainerrors "icmongolang/internal/modules/erp/domain/errors"
    "icmongolang/internal/modules/erp/domain/repository"
    valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

// ============================================================
// Product
// ============================================================

type productRepository struct{ db *gorm.DB }

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
    return &productRepository{db: db}
}

func (r *productRepository) Save(ctx context.Context, p *entity.Product) error {
    return r.db.WithContext(ctx).Clauses(clause.OnConflict{
        Columns: []clause.Column{{Name: "id"}},
        UpdateAll: true,
    }).Create(toProductModel(p)).Error
}

func (r *productRepository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Product, error) {
    var m ProductModel
    err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND id = ?", tenantID, id).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrProductNotFound
    }
    if err != nil { return nil, err }
    return toProductEntity(&m), nil
}

func (r *productRepository) FindBySKU(ctx context.Context, tenantID uuid.UUID, sku valueobject.SKU) (*entity.Product, error) {
    var m ProductModel
    err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND sku = ?", tenantID, sku.String()).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrProductNotFound
    }
    if err != nil { return nil, err }
    return toProductEntity(&m), nil
}

func (r *productRepository) List(ctx context.Context, f repository.ProductFilter) ([]*entity.Product, int64, error) {
    q := r.db.WithContext(ctx).Model(&ProductModel{}).Where("tenant_id = ?", f.TenantID)
    if f.Type != nil { q = q.Where("type = ?", string(*f.Type)) }
    if f.Category != "" { q = q.Where("category = ?", f.Category) }
    if f.IsActive != nil { q = q.Where("is_active = ?", *f.IsActive) }
    if f.Search != "" {
        s := "%" + strings.ToLower(f.Search) + "%"
        q = q.Where("(LOWER(name) LIKE ? OR LOWER(sku) LIKE ? OR LOWER(barcode) LIKE ?)", s, s, s)
    }
    if len(f.Tags) > 0 {
        for _, t := range f.Tags {
            q = q.Where("tags @> ?", `["`+t+`"]`)
        }
    }

    var total int64
    if err := q.Count(&total).Error; err != nil { return nil, 0, err }

    if f.PageSize > 0 { q = q.Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize) }

    var models []ProductModel
    if err := q.Order("name ASC").Find(&models).Error; err != nil { return nil, 0, err }
    out := make([]*entity.Product, 0, len(models))
    for i := range models { out = append(out, toProductEntity(&models[i])) }
    return out, total, nil
}

func (r *productRepository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
    res := r.db.WithContext(ctx).
        Where("tenant_id = ? AND id = ?", tenantID, id).
        Delete(&ProductModel{})
    if res.RowsAffected == 0 { return domainerrors.ErrProductNotFound }
    return res.Error
}

// Mappers
func toProductModel(p *entity.Product) *ProductModel {
    metaJSON, _ := json.Marshal(p.Metadata)
    tagsJSON, _ := json.Marshal(p.Tags)
    urlsJSON, _ := json.Marshal(p.ImageURLs)
    return &ProductModel{
        ID: p.ID, TenantID: p.TenantID,
        SKU: p.SKU.String(), Name: p.Name, Description: p.Description,
        Type: string(p.Type), Category: p.Category, Brand: p.Brand,
        UOM: p.UOM,
        CostPrice: p.CostPrice.Amount, SellPrice: p.SellPrice.Amount,
        Currency: p.Currency,
        TaxType: string(p.TaxType), TaxRate: p.TaxRate,
        TrackInventory: p.TrackInventory,
        ReorderPoint: p.ReorderPoint, ReorderQty: p.ReorderQty,
        LeadTimeDays: p.LeadTimeDays, SafetyStock: p.SafetyStock,
        Weight: p.Weight, Volume: p.Volume, Barcode: p.Barcode,
        ImageURLs: datatypes.JSON(urlsJSON),
        IsActive: p.IsActive, IsSellable: p.IsSellable, IsPurchasable: p.IsPurchasable,
        Metadata: datatypes.JSON(metaJSON), Tags: datatypes.JSON(tagsJSON),
        CreatedBy: p.CreatedBy, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
    }
}

func toProductEntity(m *ProductModel) *entity.Product {
    p := &entity.Product{
        ID: m.ID, TenantID: m.TenantID,
        SKU: valueobject.SKU(m.SKU),
        Name: m.Name, Description: m.Description,
        Type: entity.ProductType(m.Type),
        Category: m.Category, Brand: m.Brand, UOM: m.UOM,
        CostPrice: valueobject.Money{Amount: m.CostPrice, Currency: m.Currency},
        SellPrice: valueobject.Money{Amount: m.SellPrice, Currency: m.Currency},
        Currency: m.Currency,
        TaxType: valueobject.TaxType(m.TaxType), TaxRate: m.TaxRate,
        TrackInventory: m.TrackInventory,
        ReorderPoint: m.ReorderPoint, ReorderQty: m.ReorderQty,
        LeadTimeDays: m.LeadTimeDays, SafetyStock: m.SafetyStock,
        Weight: m.Weight, Volume: m.Volume, Barcode: m.Barcode,
        IsActive: m.IsActive, IsSellable: m.IsSellable, IsPurchasable: m.IsPurchasable,
        CreatedBy: m.CreatedBy, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
    }
    if len(m.Metadata) > 0 { _ = json.Unmarshal(m.Metadata, &p.Metadata) }
    if len(m.Tags) > 0     { _ = json.Unmarshal(m.Tags, &p.Tags) }
    if len(m.ImageURLs) > 0 { _ = json.Unmarshal(m.ImageURLs, &p.ImageURLs) }
    if p.Metadata == nil { p.Metadata = map[string]any{} }
    if p.Tags == nil { p.Tags = []string{} }
    return p
}

// ============================================================
// Order (with eager loading of lines)
// ============================================================

type orderRepository struct{ db *gorm.DB }

func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
    return &orderRepository{db: db}
}

func (r *orderRepository) Save(ctx context.Context, o *entity.Order) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        m := toOrderModel(o)
        if err := tx.Clauses(clause.OnConflict{
            Columns: []clause.Column{{Name: "id"}},
            UpdateAll: true,
        }).Omit("Lines").Create(m).Error; err != nil {
            return err
        }

        // Sync lines: delete removed
        ids := make([]uuid.UUID, 0, len(m.Lines))
        for _, l := range m.Lines { ids = append(ids, l.ID) }
        q := tx.Where("order_id = ?", o.ID)
        if len(ids) > 0 { q = q.Where("id NOT IN ?", ids) }
        if err := q.Delete(&OrderLineModel{}).Error; err != nil { return err }

        if len(m.Lines) > 0 {
            if err := tx.Clauses(clause.OnConflict{
                Columns:   []clause.Column{{Name: "id"}},
                DoUpdates: clause.AssignmentColumns([]string{
                    "quantity", "unit_price", "discount", "subtotal", "tax_amount",
                    "line_total", "qty_shipped", "qty_received", "qty_invoiced",
                    "warehouse_id", "notes", "updated_at",
                }),
            }).Create(&m.Lines).Error; err != nil {
                return err
            }
        }
        return nil
    })
}

func (r *orderRepository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Order, error) {
    var m OrderModel
    err := r.db.WithContext(ctx).
        Preload("Lines").Where("tenant_id = ? AND id = ?", tenantID, id).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrOrderNotFound
    }
    if err != nil { return nil, err }
    return toOrderEntity(&m), nil
}

func (r *orderRepository) FindByNo(ctx context.Context, tenantID uuid.UUID, no valueobject.DocumentNumber) (*entity.Order, error) {
    var m OrderModel
    err := r.db.WithContext(ctx).
        Preload("Lines").Where("tenant_id = ? AND order_no = ?", tenantID, no.String()).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrOrderNotFound
    }
    if err != nil { return nil, err }
    return toOrderEntity(&m), nil
}

func (r *orderRepository) List(ctx context.Context, f repository.OrderFilter) ([]*entity.Order, int64, error) {
    q := r.db.WithContext(ctx).Model(&OrderModel{}).Where("tenant_id = ?", f.TenantID)
    if f.Type != nil { q = q.Where("type = ?", string(*f.Type)) }
    if f.CustomerID != nil { q = q.Where("customer_id = ?", *f.CustomerID) }
    if f.SupplierID != nil { q = q.Where("supplier_id = ?", *f.SupplierID) }
    if len(f.Statuses) > 0 {
        ss := make([]string, len(f.Statuses))
        for i, s := range f.Statuses { ss[i] = string(s) }
        q = q.Where("status IN ?", ss)
    }
    if f.From != nil { q = q.Where("order_date >= ?", *f.From) }
    if f.To != nil   { q = q.Where("order_date <= ?", *f.To) }
    if f.Search != "" {
        s := "%" + strings.ToLower(f.Search) + "%"
        q = q.Where("(LOWER(order_no) LIKE ? OR LOWER(po_reference) LIKE ?)", s, s)
    }

    var total int64
    if err := q.Count(&total).Error; err != nil { return nil, 0, err }
    if f.PageSize > 0 { q = q.Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize) }

    var models []OrderModel
    if err := q.Preload("Lines").Order("order_date DESC").Find(&models).Error; err != nil {
        return nil, 0, err
    }
    out := make([]*entity.Order, 0, len(models))
    for i := range models { out = append(out, toOrderEntity(&models[i])) }
    return out, total, nil
}

func (r *orderRepository) NextSequence(ctx context.Context, tenantID uuid.UUID, docType valueobject.DocType, year int) (int, error) {
    gen := NewDocumentNumberGenerator(r.db)
    return gen.NextSeq(ctx, tenantID, docType, year)
}

// Mappers
func toOrderModel(o *entity.Order) *OrderModel {
    shipJSON, _ := json.Marshal(o.ShipToAddress)
    billJSON, _ := json.Marshal(o.BillToAddress)
    tagsJSON, _ := json.Marshal(o.Tags)
    metaJSON, _ := json.Marshal(o.Metadata)

    m := &OrderModel{
        ID: o.ID, TenantID: o.TenantID,
        OrderNo: o.OrderNo.String(), Type: string(o.Type),
        CustomerID: o.CustomerID, SupplierID: o.SupplierID,
        Status:    string(o.Status),
        OrderDate: o.OrderDate, ExpectedAt: o.ExpectedAt,
        ConfirmedAt: o.ConfirmedAt, ShippedAt: o.ShippedAt,
        ReceivedAt: o.ReceivedAt, ClosedAt: o.ClosedAt, CancelledAt: o.CancelledAt,
        ShipToAddress: datatypes.JSON(shipJSON),
        BillToAddress: datatypes.JSON(billJSON),
        ShippingMethod: o.ShippingMethod, TrackingNumber: o.TrackingNumber,
        Subtotal: o.Subtotal, DiscountAmount: o.DiscountAmount,
        TaxAmount: o.TaxAmount, ShippingFee: o.ShippingFee,
        TotalAmount: o.TotalAmount, Currency: o.Currency,
        FXRate: o.FXRate, BaseCurrency: o.BaseCurrency, BaseTotal: o.BaseTotal,
        POReference: o.POReference,
        QuotationID: o.QuotationID, ContractID: o.ContractID,
        Notes: o.Notes, InternalNotes: o.InternalNotes,
        Tags: datatypes.JSON(tagsJSON), Metadata: datatypes.JSON(metaJSON),
        CreatedBy: o.CreatedBy, ApprovedBy: o.ApprovedBy,
        CreatedAt: o.CreatedAt, UpdatedAt: o.UpdatedAt,
    }
    for _, l := range o.Lines {
        m.Lines = append(m.Lines, OrderLineModel{
            ID: l.ID, OrderID: o.ID, LineNo: l.LineNo,
            ProductID: l.ProductID, SKU: l.SKU, Name: l.Name,
            Description: l.Description,
            Quantity: l.Quantity, UOM: l.UOM,
            UnitPrice: l.UnitPrice, Discount: l.Discount, DiscountPct: l.DiscountPct,
            TaxType: string(l.TaxType), TaxRate: l.TaxRate,
            Subtotal: l.Subtotal, TaxAmount: l.TaxAmount, LineTotal: l.LineTotal,
            QtyShipped: l.QtyShipped, QtyReceived: l.QtyReceived, QtyInvoiced: l.QtyInvoiced,
            WarehouseID: l.WarehouseID, Notes: l.Notes,
            CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
        })
    }
    return m
}

func toOrderEntity(m *OrderModel) *entity.Order {
    o := &entity.Order{
        ID: m.ID, TenantID: m.TenantID,
        OrderNo: valueobject.DocumentNumber(m.OrderNo),
        Type:    valueobject.OrderType(m.Type),
        CustomerID: m.CustomerID, SupplierID: m.SupplierID,
        Status:    valueobject.OrderStatus(m.Status),
        OrderDate: m.OrderDate, ExpectedAt: m.ExpectedAt,
        ConfirmedAt: m.ConfirmedAt, ShippedAt: m.ShippedAt,
        ReceivedAt: m.ReceivedAt, ClosedAt: m.ClosedAt, CancelledAt: m.CancelledAt,
        ShippingMethod: m.ShippingMethod, TrackingNumber: m.TrackingNumber,
        Subtotal: m.Subtotal, DiscountAmount: m.DiscountAmount,
        TaxAmount: m.TaxAmount, ShippingFee: m.ShippingFee,
        TotalAmount: m.TotalAmount, Currency: m.Currency,
        FXRate: m.FXRate, BaseCurrency: m.BaseCurrency, BaseTotal: m.BaseTotal,
        POReference: m.POReference,
        QuotationID: m.QuotationID, ContractID: m.ContractID,
        Notes: m.Notes, InternalNotes: m.InternalNotes,
        CreatedBy: m.CreatedBy, ApprovedBy: m.ApprovedBy,
        CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
        Lines: []*entity.OrderLine{},
    }
    if len(m.ShipToAddress) > 0 { _ = json.Unmarshal(m.ShipToAddress, &o.ShipToAddress) }
    if len(m.BillToAddress) > 0 { _ = json.Unmarshal(m.BillToAddress, &o.BillToAddress) }
    if len(m.Tags) > 0     { _ = json.Unmarshal(m.Tags, &o.Tags) }
    if len(m.Metadata) > 0 { _ = json.Unmarshal(m.Metadata, &o.Metadata) }
    if o.Tags == nil { o.Tags = []string{} }
    if o.Metadata == nil { o.Metadata = map[string]any{} }

    for _, lm := range m.Lines {
        line := &entity.OrderLine{
            ID: lm.ID, OrderID: lm.OrderID, LineNo: lm.LineNo,
            ProductID: lm.ProductID, SKU: lm.SKU, Name: lm.Name,
            Description: lm.Description,
            Quantity: lm.Quantity, UOM: lm.UOM,
            UnitPrice: lm.UnitPrice, Discount: lm.Discount, DiscountPct: lm.DiscountPct,
            TaxType: valueobject.TaxType(lm.TaxType), TaxRate: lm.TaxRate,
            Subtotal: lm.Subtotal, TaxAmount: lm.TaxAmount, LineTotal: lm.LineTotal,
            QtyShipped: lm.QtyShipped, QtyReceived: lm.QtyReceived, QtyInvoiced: lm.QtyInvoiced,
            WarehouseID: lm.WarehouseID, Notes: lm.Notes,
            CreatedAt: lm.CreatedAt, UpdatedAt: lm.UpdatedAt,
        }
        o.Lines = append(o.Lines, line)
    }
    return o
}
```

> **หมายเหตุ**: repositories อื่นๆ (`Invoice`, `Payment`, `Inventory`, `Warehouse`, `Journal`, `Supplier`, `CreditNote`, `PriceList`) ทำตาม pattern เดียวกัน

### C.4 Kafka Producer / Consumer

```go
package messaging

import (
    "context"
    "encoding/json"
    "time"

    "github.com/IBM/sarama"
)

type kafkaProducer struct{ p sarama.SyncProducer }

func NewKafkaProducer(brokers []string, clientID string) (*kafkaProducer, error) {
    cfg := sarama.NewConfig()
    cfg.ClientID = clientID
    cfg.Producer.Return.Successes = true
    cfg.Producer.RequiredAcks = sarama.WaitForAll
    cfg.Producer.Retry.Max = 5
    cfg.Producer.Idempotent = true
    cfg.Net.MaxOpenRequests = 1

    p, err := sarama.NewSyncProducer(brokers, cfg)
    if err != nil { return nil, err }
    return &kafkaProducer{p: p}, nil
}

func (k *kafkaProducer) Publish(ctx context.Context, topic, key string, payload any) error {
    b, err := json.Marshal(payload)
    if err != nil { return err }
    _, _, err = k.p.SendMessage(&sarama.ProducerMessage{
        Topic: topic, Key: sarama.StringEncoder(key),
        Value: sarama.ByteEncoder(b), Timestamp: time.Now(),
    })
    return err
}

func (k *kafkaProducer) Close() error { return k.p.Close() }
```

#### `infrastructure/messaging/consumers/device_alert_to_order_consumer.go`
```go
package consumers

import (
    "context"
    "encoding/json"

    "github.com/google/uuid"

    "icmongolang/internal/modules/erp/application"
)

// DeviceAlertToOrderConsumer – CRITICAL device alert → สร้าง PR (Purchase Requisition)
type DeviceAlertToOrderConsumer struct {
    createOrderUC *application.CreateOrderUseCase
}

func (c *DeviceAlertToOrderConsumer) Handle(ctx context.Context, payload []byte) error {
    var evt struct {
        TenantID   uuid.UUID `json:"tenant_id"`
        DeviceID   uuid.UUID `json:"device_id"`
        CustomerID uuid.UUID `json:"customer_id"`
        Severity   string    `json:"severity"`
        Message    string    `json:"message"`
    }
    _ = json.Unmarshal(payload, &evt)

    if evt.Severity != "CRITICAL" { return nil }
    // ... create PR
    return nil
}
```

### C.5 Scheduler Jobs

#### `infrastructure/scheduler/invoice_overdue_job.go`
```go
package scheduler

import (
    "context"
    "log"
    "time"

    "github.com/google/uuid"

    "icmongolang/internal/modules/erp/domain/event"
    "icmongolang/internal/modules/erp/domain/repository"
    "icmongolang/pkg/kafka"
)

type InvoiceOverdueJob struct {
    invoiceRepo repository.InvoiceRepository
    producer    kafka.Producer
}

func NewInvoiceOverdueJob(invoiceRepo repository.InvoiceRepository, producer kafka.Producer) *InvoiceOverdueJob {
    return &InvoiceOverdueJob{invoiceRepo: invoiceRepo, producer: producer}
}

func (j *InvoiceOverdueJob) Run() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    asOf := time.Now()
    overdue, err := j.invoiceRepo.ListOverdue(ctx, asOf, 500)
    if err != nil {
        log.Printf("[invoice.overdue] list: %v", err)
        return
    }

    for _, inv := range overdue {
        // Mark if not already
        if inv.Status == "ISSUED" || inv.Status == "PARTIAL" {
            _ = inv.MarkOverdue()
            _ = j.invoiceRepo.Save(ctx, inv)
        }

        _ = j.producer.Publish(ctx, event.TopicInvoiceOverdue, inv.ID.String(), event.InvoiceOverdue{
            EventID: uuid.New(), InvoiceID: inv.ID, TenantID: inv.TenantID,
            InvoiceNo: inv.InvoiceNo.String(),
            AmountDue: inv.AmountDue, Currency: inv.Currency,
            DaysOverdue: inv.DaysOverdue(),
            OccurredAt:  asOf,
        })
    }
    log.Printf("[invoice.overdue] processed %d invoices", len(overdue))
}
```

#### `infrastructure/scheduler/low_stock_alert_job.go`
```go
package scheduler

import (
    "context"
    "log"
    "time"

    "github.com/google/uuid"

    "icmongolang/internal/modules/erp/domain/event"
    "icmongolang/internal/modules/erp/domain/repository"
    "icmongolang/internal/modules/erp/domain/service"
    "icmongolang/pkg/kafka"
)

type LowStockAlertJob struct {
    reorderSvc *service.ReorderService
    producer   kafka.Producer
}

func (j *LowStockAlertJob) Run() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    // In prod: iterate over tenants
    tenantIDs := []uuid.UUID{}

    for _, tid := range tenantIDs {
        suggestions, err := j.reorderSvc.ScanLowStock(ctx, tid)
        if err != nil {
            log.Printf("[low.stock] scan %s: %v", tid, err)
            continue
        }
        for _, s := range suggestions {
            _ = j.producer.Publish(ctx, event.TopicStockLow, s.ProductID.String(), event.StockLow{
                EventID: uuid.New(), TenantID: tid,
                ProductID: s.ProductID, SKU: s.SKU,
                WarehouseID: s.WarehouseID,
                CurrentStock: s.CurrentStock, ReorderPoint: s.ReorderPoint,
                OccurredAt: time.Now(),
            })
        }
    }
}

var _ = repository.ProductFilter{}
```

### C.6 Redis Cache (Aging report)

```go
package redis

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    goredis "github.com/go-redis/redis/v8"
    "github.com/google/uuid"
)

type AgingCache interface {
    Set(ctx context.Context, tenantID uuid.UUID, key string, v any, ttl time.Duration) error
    Get(ctx context.Context, tenantID uuid.UUID, key string) ([]byte, error)
    Invalidate(ctx context.Context, tenantID uuid.UUID) error
}

type agingCache struct{ client *goredis.Client }

func NewAgingCache(c *goredis.Client) AgingCache { return &agingCache{client: c} }

func (a *agingCache) key(tenantID uuid.UUID, k string) string {
    return fmt.Sprintf("erp:aging:%s:%s", tenantID, k)
}

func (a *agingCache) Set(ctx context.Context, tenantID uuid.UUID, k string, v any, ttl time.Duration) error {
    b, _ := json.Marshal(v)
    return a.client.Set(ctx, a.key(tenantID, k), b, ttl).Err()
}

func (a *agingCache) Get(ctx context.Context, tenantID uuid.UUID, k string) ([]byte, error) {
    return a.client.Get(ctx, a.key(tenantID, k)).Bytes()
}

func (a *agingCache) Invalidate(ctx context.Context, tenantID uuid.UUID) error {
    iter := a.client.Scan(ctx, 0, fmt.Sprintf("erp:aging:%s:*", tenantID), 100).Iterator()
    for iter.Next(ctx) {
        _ = a.client.Del(ctx, iter.Val()).Err()
    }
    return iter.Err()
}
```

---

## 🅳 PART 4D — INTERFACE + WIRING + MIGRATION + TESTS

### D.1 HTTP Handlers

```go
package http

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    "icmongolang/internal/modules/erp/application"
)

type ProductHandler struct {
    createUC *application.CreateProductUseCase
    listUC   *application.ListProductsUseCase
    getUC    *application.GetProductUseCase
    updateUC *application.UpdateProductUseCase
}

func (h *ProductHandler) Create(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    var in application.CreateProductInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    in.TenantID = tid.String()
    in.ActorID = uid.String()
    res, err := h.createUC.Execute(c.Request.Context(), in)
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusCreated, res)
}

func (h *ProductHandler) List(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    var active *bool
    if v := c.Query("active"); v == "true" { t := true; active = &t } else if v == "false" { f := false; active = &f }
    res, err := h.listUC.Execute(c.Request.Context(), application.ListProductsInput{
        TenantID: tid.String(),
        Type:     c.Query("type"),
        Category: c.Query("category"),
        Search:   c.Query("q"),
        Active:   active,
        Tags:     c.QueryArray("tag"),
        Page:     parseInt(c.DefaultQuery("page", "1")),
        PageSize: parseInt(c.DefaultQuery("page_size", "20")),
    })
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusOK, res)
}

type OrderHandler struct {
    createUC *application.CreateOrderUseCase
    shipUC   *application.ShipOrderUseCase
    recvUC   *application.ReceiveOrderUseCase
    cancelUC *application.CancelOrderUseCase
    getUC    *application.GetOrderUseCase
    listUC   *application.ListOrdersUseCase
}

func (h *OrderHandler) Create(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    var in application.CreateOrderInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    in.TenantID = tid.String()
    in.ActorID = uid.String()
    res, err := h.createUC.Execute(c.Request.Context(), in)
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusCreated, res)
}

func (h *OrderHandler) Ship(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))
    var in application.ShipOrderInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    in.TenantID = tid.String()
    in.OrderID = id.String()
    in.ActorID = uid.String()
    res, err := h.shipUC.Execute(c.Request.Context(), in)
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusOK, res)
}

type InvoiceHandler struct {
    createUC *application.CreateInvoiceUseCase
    getUC    *application.GetInvoiceUseCase
    listUC   *application.ListInvoicesUseCase
    voidUC   *application.VoidInvoiceUseCase
    agingUC  *application.AgingReportUseCase
}

func (h *InvoiceHandler) Create(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    var in application.CreateInvoiceInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    in.TenantID = tid.String()
    in.ActorID = uid.String()
    res, err := h.createUC.Execute(c.Request.Context(), in)
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusCreated, res)
}

type PaymentHandler struct {
    createUC *application.RecordPaymentUseCase
    listUC   *application.ListPaymentsUseCase
    getUC    *application.GetPaymentUseCase
}

func (h *PaymentHandler) Create(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    var in application.CreatePaymentInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    in.TenantID = tid.String()
    in.ActorID = uid.String()
    res, err := h.createUC.Execute(c.Request.Context(), in)
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusCreated, res)
}

type InventoryHandler struct {
    adjustUC   *application.AdjustInventoryUseCase
    transferUC *application.TransferStockUseCase
    listUC     *application.ListInventoryUseCase
    movementUC *application.ListMovementsUseCase
}

func (h *InventoryHandler) Adjust(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    var in application.AdjustInventoryInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    in.TenantID = tid.String()
    in.ActorID = uid.String()
    res, err := h.adjustUC.Execute(c.Request.Context(), in)
    if err != nil { writeError(c, err); return }
    c.JSON(http.StatusOK, res)
}

func (h *InventoryHandler) Transfer(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid, _ := c.MustGet("user_id").(uuid.UUID)
    var in application.TransferStockInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    in.TenantID = tid.String()
    in.ActorID = uid.String()
    if err := h.transferUC.Execute(c.Request.Context(), in); err != nil {
        writeError(c, err); return
    }
    c.JSON(http.StatusOK, gin.H{"status": "transferred"})
}

type WarehouseHandler struct {
    createUC *application.CreateWarehouseUseCase
    listUC   *application.ListWarehousesUseCase
}
```

### D.2 Routes

```go
func RegisterRoutes(r *gin.RouterGroup, h *Handlers, auth, tenant gin.HandlerFunc) {
    // Products
    p := r.Group("/erp/products")
    p.Use(auth, tenant)
    p.POST("",      h.Product.Create)
    p.GET ("",      h.Product.List)
    p.GET ("/:id",  h.Product.Get)
    p.PUT ("/:id",  h.Product.Update)
    p.DELETE("/:id", h.Product.Delete)

    // Warehouses
    w := r.Group("/erp/warehouses")
    w.Use(auth, tenant)
    w.POST("",     h.Warehouse.Create)
    w.GET ("",     h.Warehouse.List)

    // Inventory
    i := r.Group("/erp/inventory")
    i.Use(auth, tenant)
    i.GET ("",              h.Inventory.List)
    i.GET ("/:productID",   h.Inventory.GetByProduct)
    i.POST("/adjust",       h.Inventory.Adjust)
    i.POST("/transfer",     h.Inventory.Transfer)
    i.GET ("/movements",    h.Inventory.ListMovements)

    // Orders
    o := r.Group("/erp/orders")
    o.Use(auth, tenant)
    o.POST("",                    h.Order.Create)
    o.GET ("",                    h.Order.List)
    o.GET ("/:id",                h.Order.Get)
    o.POST("/:id/confirm",        h.Order.Confirm)
    o.POST("/:id/ship",           h.Order.Ship)
    o.POST("/:id/receive",        h.Order.Receive)
    o.POST("/:id/cancel",         h.Order.Cancel)

    // Invoices
    inv := r.Group("/erp/invoices")
    inv.Use(auth, tenant)
    inv.POST("",                h.Invoice.Create)
    inv.GET ("",                h.Invoice.List)
    inv.GET ("/:id",            h.Invoice.Get)
    inv.POST("/:id/void",       h.Invoice.Void)
    inv.GET ("/aging",          h.Invoice.Aging)

    // Payments
    pay := r.Group("/erp/payments")
    pay.Use(auth, tenant)
    pay.POST("",       h.Payment.Create)
    pay.GET ("",       h.Payment.List)
    pay.GET ("/:id",   h.Payment.Get)

    // Suppliers
    s := r.Group("/erp/suppliers")
    s.Use(auth, tenant)
    s.POST("",     h.Supplier.Create)
    s.GET ("",     h.Supplier.List)

    // Reports
    rep := r.Group("/erp/reports")
    rep.Use(auth, tenant)
    rep.GET("/aging",      h.Invoice.Aging)
    rep.GET("/stock-value", h.Inventory.StockValue)
}
```

### D.3 Errors

```go
func writeError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, domainerrors.ErrProductNotFound),
        errors.Is(err, domainerrors.ErrOrderNotFound),
        errors.Is(err, domainerrors.ErrInvoiceNotFound),
        errors.Is(err, domainerrors.ErrPaymentNotFound),
        errors.Is(err, domainerrors.ErrWarehouseNotFound),
        errors.Is(err, domainerrors.ErrInventoryNotFound),
        errors.Is(err, domainerrors.ErrSupplierNotFound),
        errors.Is(err, domainerrors.ErrCreditNoteNotFound),
        errors.Is(err, domainerrors.ErrPriceListNotFound):
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

    case errors.Is(err, domainerrors.ErrSKUDuplicate),
        errors.Is(err, domainerrors.ErrWarehouseDuplicate),
        errors.Is(err, domainerrors.ErrPriceEntryDuplicate):
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

    case errors.Is(err, domainerrors.ErrInsufficientStock),
        errors.Is(err, domainerrors.ErrInsufficientAvailable),
        errors.Is(err, domainerrors.ErrOrderNotEditable),
        errors.Is(err, domainerrors.ErrInvoiceNotEditable),
        errors.Is(err, domainerrors.ErrInvoiceNotPayable),
        errors.Is(err, domainerrors.ErrJournalNotBalanced),
        errors.Is(err, domainerrors.ErrOverShip),
        errors.Is(err, domainerrors.ErrOverReceive),
        errors.Is(err, domainerrors.ErrOverAllocation):
        c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})

    case errors.Is(err, domainerrors.ErrPersistenceFailure),
        errors.Is(err, domainerrors.ErrPDFRenderFailure),
        errors.Is(err, domainerrors.ErrNotifierFailure),
        errors.Is(err, domainerrors.ErrPaymentGateway):
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})

    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
```

### D.4 Composition Root

```go
package erp

import (
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "github.com/minio/minio-go/v7"
    "gorm.io/gorm"

    "icmongolang/internal/modules/erp/application"
    "icmongolang/internal/modules/erp/domain/service"
    "icmongolang/internal/modules/erp/domain/service/port"
    "icmongolang/internal/modules/erp/infrastructure/persistence/postgres"
    httpiface "icmongolang/internal/modules/erp/interfaces/http"
    "icmongolang/pkg/kafka"
    "icmongolang/pkg/logger"
    "icmongolang/pkg/transaction"
)

type Deps struct {
    DB         *gorm.DB
    Redis      *redis.Client
    Producer   kafka.Producer
    MinIO      *minio.Client
    Bucket     string
    CustomerCli port.CustomerPort
    PaymentCli  port.PaymentPort
    Notifier    port.NotifierPort
    Logger     logger.Logger
}

func Init(router *gin.RouterGroup, deps Deps, auth, tenant gin.HandlerFunc) {
    // Repos
    prodRepo := postgres.NewProductRepository(deps.DB)
    whRepo   := postgres.NewWarehouseRepository(deps.DB)
    invRepo  := postgres.NewInventoryRepository(deps.DB)
    orderRepo := postgres.NewOrderRepository(deps.DB)
    invRepo2 := postgres.NewInvoiceRepository(deps.DB)
    payRepo  := postgres.NewPaymentRepository(deps.DB)
    supRepo  := postgres.NewSupplierRepository(deps.DB)
    jRepo    := postgres.NewJournalRepository(deps.DB)
    cnRepo   := postgres.NewCreditNoteRepository(deps.DB)
    plRepo   := postgres.NewPriceListRepository(deps.DB)
    auditRepo := postgres.NewAuditRepository(deps.DB)

    seqGen := postgres.NewDocumentNumberGenerator(deps.DB)

    // Domain services
    numSvc  := service.NewDocumentNumberService(seqGen)
    taxSvc  := service.NewTaxCalculator()
    invSvc  := service.NewInventoryService(invRepo)
    allocSvc := service.NewAllocationService()
    jSvc    := service.NewJournalService(jRepo, seqGen)
    pricingSvc := service.NewPricingService(plRepo, prodRepo)
    reorderSvc := service.NewReorderService(invRepo, prodRepo, orderRepo, numSvc)

    _ = taxSvc
    _ = allocSvc

    txMgr := transaction.NewManager(deps.DB)

    // Use cases
    createProdUC := application.NewCreateProductUseCase(prodRepo, auditRepo, deps.Producer, deps.Logger)
    listProdUC := application.NewListProductsUseCase(prodRepo)
    getProdUC := application.NewGetProductUseCase(prodRepo)
    updateProdUC := application.NewUpdateProductUseCase(prodRepo, auditRepo, deps.Logger)

    createOrderUC := application.NewCreateOrderUseCase(
        orderRepo, prodRepo, whRepo, invSvc, pricingSvc,
        numSvc, auditRepo, deps.Producer, deps.Logger,
    )
    shipOrderUC := application.NewShipOrderUseCase(
        orderRepo, invSvc, jSvc, auditRepo,
        deps.Producer, txMgr, deps.Logger,
    )
    getOrderUC := application.NewGetOrderUseCase(orderRepo)
    listOrderUC := application.NewListOrdersUseCase(orderRepo)
    cancelOrderUC := application.NewCancelOrderUseCase(orderRepo, auditRepo, deps.Producer, deps.Logger)

    createInvUC := application.NewCreateInvoiceUseCase(
        invRepo2, orderRepo, prodRepo, jSvc,
        numSvc, auditRepo, deps.Producer, deps.Logger,
    )
    getInvUC := application.NewGetInvoiceUseCase(invRepo2)
    listInvUC := application.NewListInvoicesUseCase(invRepo2)
    agingUC := application.NewAgingReportUseCase(invRepo2, deps.Redis)

    recordPayUC := application.NewRecordPaymentUseCase(
        payRepo, invRepo2, jSvc, numSvc,
        auditRepo, deps.Producer, txMgr, deps.Logger,
    )

    adjustInvUC := application.NewAdjustInventoryUseCase(
        invRepo, prodRepo, invSvc, auditRepo, deps.Producer, deps.Logger,
    )
    transferUC := application.NewTransferStockUseCase(
        invRepo, invSvc, deps.Producer, txMgr, deps.Logger,
    )

    // Handlers
    prodH := &httpiface.ProductHandler{CreateUC: createProdUC, ListUC: listProdUC, GetUC: getProdUC, UpdateUC: updateProdUC}
    orderH := &httpiface.OrderHandler{CreateUC: createOrderUC, ShipUC: shipOrderUC, GetUC: getOrderUC, ListUC: listOrderUC, CancelUC: cancelOrderUC}
    invH := &httpiface.InvoiceHandler{CreateUC: createInvUC, GetUC: getInvUC, ListUC: listInvUC, AgingUC: agingUC}
    payH := &httpiface.PaymentHandler{CreateUC: recordPayUC}
    invH2 := &httpiface.InventoryHandler{AdjustUC: adjustInvUC, TransferUC: transferUC}

    httpiface.RegisterRoutes(router, &httpiface.Handlers{
        Product: prodH, Order: orderH, Invoice: invH,
        Payment: payH, Inventory: invH2,
    }, auth, tenant)

    _ = reorderSvc
    _ = cnRepo
    _ = supRepo
}

var _ = application.NewCreateWarehouseUseCase
```

### D.5 Migration SQL

```sql
-- ============================================================
-- erp module — initial schema
-- Prefix: erp_
-- ============================================================

-- ============================================================
-- PRODUCTS
-- ============================================================
CREATE TABLE IF NOT EXISTS erp_products (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL,
    sku               VARCHAR(50) NOT NULL,
    name              VARCHAR(255) NOT NULL,
    description       TEXT,
    type              VARCHAR(20) NOT NULL,
    category          VARCHAR(100),
    brand             VARCHAR(100),
    uom               VARCHAR(20) NOT NULL,
    cost_price        NUMERIC(15,2) NOT NULL DEFAULT 0,
    sell_price        NUMERIC(15,2) NOT NULL DEFAULT 0,
    currency          VARCHAR(3) NOT NULL DEFAULT 'THB',
    tax_type          VARCHAR(20) DEFAULT 'VAT',
    tax_rate          NUMERIC(5,4) DEFAULT 0.07,
    track_inventory   BOOLEAN NOT NULL DEFAULT FALSE,
    reorder_point     NUMERIC(15,4) DEFAULT 0,
    reorder_qty       NUMERIC(15,4) DEFAULT 0,
    lead_time_days    INT DEFAULT 0,
    safety_stock      NUMERIC(15,4) DEFAULT 0,
    weight            NUMERIC(10,3) DEFAULT 0,
    volume            NUMERIC(10,4) DEFAULT 0,
    barcode           VARCHAR(100),
    image_urls        JSONB DEFAULT '[]'::jsonb,
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,
    is_sellable       BOOLEAN DEFAULT TRUE,
    is_purchasable    BOOLEAN DEFAULT TRUE,
    metadata          JSONB DEFAULT '{}'::jsonb,
    tags              JSONB DEFAULT '[]'::jsonb,
    created_by        UUID,
    created_at        TIMESTAMP DEFAULT NOW(),
    updated_at        TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_erp_prod_tenant_sku UNIQUE (tenant_id, sku)
);
CREATE INDEX idx_erp_prod_tenant_active ON erp_products(tenant_id, is_active);
CREATE INDEX idx_erp_prod_type          ON erp_products(type);
CREATE INDEX idx_erp_prod_category      ON erp_products(category);
CREATE INDEX idx_erp_prod_barcode       ON erp_products(barcode) WHERE barcode IS NOT NULL;
CREATE INDEX idx_erp_prod_tags_gin      ON erp_products USING GIN (tags);

-- ============================================================
-- WAREHOUSES
-- ============================================================
CREATE TABLE IF NOT EXISTS erp_warehouses (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL,
    code        VARCHAR(30) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    type        VARCHAR(20) NOT NULL,
    address     JSONB,
    manager_id  UUID,
    phone       VARCHAR(30),
    email       VARCHAR(255),
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    is_default  BOOLEAN DEFAULT FALSE,
    metadata    JSONB,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_erp_wh_tenant_code UNIQUE (tenant_id, code)
);
CREATE INDEX idx_erp_wh_default ON erp_warehouses(tenant_id, is_default) WHERE is_default = true;

-- ============================================================
-- INVENTORY ITEMS
-- ============================================================
CREATE TABLE IF NOT EXISTS erp_inventory_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    product_id      UUID NOT NULL REFERENCES erp_products(id),
    warehouse_id    UUID NOT NULL REFERENCES erp_warehouses(id),
    qty_on_hand     NUMERIC(15,4) NOT NULL DEFAULT 0,
    qty_reserved    NUMERIC(15,4) NOT NULL DEFAULT 0,
    avg_cost        NUMERIC(15,4) NOT NULL DEFAULT 0,
    last_cost       NUMERIC(15,4) NOT NULL DEFAULT 0,
    bin_location    VARCHAR(50),
    updated_at      TIMESTAMP DEFAULT NOW(),
    created_at      TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_erp_inv_unique UNIQUE (tenant_id, product_id, warehouse_id),
    CONSTRAINT ck_qty_onhand_nonneg CHECK (qty_on_hand >= 0),
    CONSTRAINT ck_qty_reserved_nonneg CHECK (qty_reserved >= 0),
    CONSTRAINT ck_qty_reserved_le_onhand CHECK (qty_reserved <= qty_on_hand)
);
CREATE INDEX idx_erp_inv_tenant_product ON erp_inventory_items(tenant_id, product_id);
CREATE INDEX idx_erp_inv_warehouse      ON erp_inventory_items(warehouse_id);

-- ============================================================
-- STOCK MOVEMENTS (ledger, immutable)
-- ============================================================
CREATE TABLE IF NOT EXISTS erp_stock_movements (
    id                    BIGSERIAL PRIMARY KEY,
    tenant_id             UUID NOT NULL,
    product_id            UUID NOT NULL,
    warehouse_id          UUID NOT NULL,
    type                  VARCHAR(20) NOT NULL,
    direction             INT NOT NULL,
    quantity              NUMERIC(15,4) NOT NULL,
    uom                   VARCHAR(20),
    unit_cost             NUMERIC(15,4),
    total_cost            NUMERIC(15,2),
    qty_on_hand_after     NUMERIC(15,4),
    qty_reserved_after    NUMERIC(15,4),
    ref_type              VARCHAR(30),
    ref_id                UUID,
    ref_number            VARCHAR(50),
    reason                TEXT,
    notes                 TEXT,
    counter_warehouse_id  UUID,
    occurred_at           TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by            UUID
);
CREATE INDEX idx_erp_mv_tenant_product ON erp_stock_movements(tenant_id, product_id, occurred_at DESC);
CREATE INDEX idx_erp_mv_type           ON erp_stock_movements(type);
CREATE INDEX idx_erp_mv_ref            ON erp_stock_movements(ref_type, ref_id);

-- ============================================================
-- STOCK RESERVATIONS
-- ============================================================
CREATE TABLE IF NOT EXISTS erp_stock_reservations (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID NOT NULL,
    product_id     UUID NOT NULL,
    warehouse_id   UUID NOT NULL,
    order_id       UUID NOT NULL,
    order_line_id  UUID NOT NULL,
    quantity       NUMERIC(15,4) NOT NULL,
    status         VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    expires_at     TIMESTAMP,
    consumed_at    TIMESTAMP,
    released_at    TIMESTAMP,
    created_at     TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_erp_res_tenant_status ON erp_stock_reservations(tenant_id, status);
CREATE INDEX idx_erp_res_order         ON erp_stock_reservations(order_id);

-- ============================================================
-- ORDERS
-- ============================================================
CREATE TABLE IF NOT EXISTS erp_orders (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL,
    order_no          VARCHAR(30) NOT NULL,
    type              VARCHAR(20) NOT NULL,
    customer_id       UUID,
    supplier_id       UUID,
    status            VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    order_date        TIMESTAMP NOT NULL,
    expected_at       TIMESTAMP,
    confirmed_at      TIMESTAMP,
    shipped_at        TIMESTAMP,
    received_at       TIMESTAMP,
    closed_at         TIMESTAMP,
    cancelled_at      TIMESTAMP,
    ship_to_address   JSONB,
    bill_to_address   JSONB,
    shipping_method   VARCHAR(50),
    tracking_number   VARCHAR(100),
    subtotal          NUMERIC(15,2) DEFAULT 0,
    discount_amount   NUMERIC(15,2) DEFAULT 0,
    tax_amount        NUMERIC(15,2) DEFAULT 0,
    shipping_fee      NUMERIC(15,2) DEFAULT 0,
    total_amount      NUMERIC(15,2) DEFAULT 0,
    currency          VARCHAR(3) NOT NULL DEFAULT 'THB',
    fx_rate           NUMERIC(15,6) DEFAULT 1,
    base_currency     VARCHAR(3) DEFAULT 'THB',
    base_total        NUMERIC(15,2) DEFAULT 0,
    po_reference      VARCHAR(100),
    quotation_id      UUID,
    contract_id       UUID,
    notes             TEXT,
    internal_notes    TEXT,
    tags              JSONB DEFAULT '[]'::jsonb,
    metadata          JSONB DEFAULT '{}'::jsonb,
    created_by        UUID,
    approved_by       UUID,
    created_at        TIMESTAMP DEFAULT NOW(),
    updated_at        TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_erp_order_tenant_no UNIQUE (tenant_id, order_no)
);
CREATE INDEX idx_erp_order_tenant_status ON erp_orders(tenant_id, status);
CREATE INDEX idx_erp_order_customer      ON erp_orders(customer_id) WHERE customer_id IS NOT NULL;
CREATE INDEX idx_erp_order_supplier      ON erp_orders(supplier_id) WHERE supplier_id IS NOT NULL;
CREATE INDEX idx_erp_order_date          ON erp_orders(tenant_id, order_date DESC);

CREATE TABLE IF NOT EXISTS erp_order_lines (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id      UUID NOT NULL REFERENCES erp_orders(id) ON DELETE CASCADE,
    line_no       INT NOT NULL,
    product_id    UUID NOT NULL,
    sku           VARCHAR(50) NOT NULL,
    name          VARCHAR(255) NOT NULL,
    description   TEXT,
    quantity      NUMERIC(15,4) NOT NULL,
    uom           VARCHAR(20) NOT NULL,
    unit_price    NUMERIC(15,2) NOT NULL,
    discount      NUMERIC(15,2) DEFAULT 0,
    discount_pct  NUMERIC(5,2) DEFAULT 0,
    tax_type      VARCHAR(20),
    tax_rate      NUMERIC(5,4),
    subtotal      NUMERIC(15,2) NOT NULL,
    tax_amount    NUMERIC(15,2) DEFAULT 0,
    line_total    NUMERIC(15,2) NOT NULL,
    qty_shipped   NUMERIC(15,4) DEFAULT 0,
    qty_received  NUMERIC(15,4) DEFAULT 0,
    qty_invoiced  NUMERIC(15,4) DEFAULT 0,
    warehouse_id  UUID,
    notes         TEXT,
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_erp_order_lines_order ON erp_order_lines(order_id, line_no);

-- ============================================================
-- INVOICES
-- ============================================================
CREATE TABLE IF NOT EXISTS erp_invoices (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL,
    invoice_no        VARCHAR(30) NOT NULL,
    type              VARCHAR(5) NOT NULL,
    customer_id       UUID,
    supplier_id       UUID,
    order_id          UUID,
    status            VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    issue_date        TIMESTAMP NOT NULL,
    due_date          TIMESTAMP NOT NULL,
    paid_at           TIMESTAMP,
    voided_at         TIMESTAMP,
    void_reason       TEXT,
    subtotal          NUMERIC(15,2) NOT NULL DEFAULT 0,
    discount_amount   NUMERIC(15,2) DEFAULT 0,
    tax_amount        NUMERIC(15,2) DEFAULT 0,
    wht_amount        NUMERIC(15,2) DEFAULT 0,
    total_amount      NUMERIC(15,2) NOT NULL DEFAULT 0,
    amount_paid       NUMERIC(15,2) DEFAULT 0,
    amount_due        NUMERIC(15,2) DEFAULT 0,
    currency          VARCHAR(3) NOT NULL DEFAULT 'THB',
    fx_rate           NUMERIC(15,6) DEFAULT 1,
    base_currency     VARCHAR(3) DEFAULT 'THB',
    base_total        NUMERIC(15,2) DEFAULT 0,
    payment_terms     VARCHAR(50),
    notes             TEXT,
    billing_address   JSONB,
    po_reference      VARCHAR(100),
    ref_number        VARCHAR(100),
    metadata          JSONB DEFAULT '{}'::jsonb,
    tags              JSONB DEFAULT '[]'::jsonb,
    created_by        UUID,
    issued_by         UUID,
    created_at        TIMESTAMP DEFAULT NOW(),
    updated_at        TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_erp_inv_tenant_no UNIQUE (tenant_id, invoice_no)
);
CREATE INDEX idx_erp_inv_tenant_status ON erp_invoices(tenant_id, status);
CREATE INDEX idx_erp_inv_customer      ON erp_invoices(customer_id) WHERE customer_id IS NOT NULL;
CREATE INDEX idx_erp_inv_supplier      ON erp_invoices(supplier_id) WHERE supplier_id IS NOT NULL;
CREATE INDEX idx_erp_inv_due_date      ON erp_invoices(due_date) WHERE status IN ('ISSUED', 'PARTIAL', 'OVERDUE');

CREATE TABLE IF NOT EXISTS erp_invoice_lines (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id     UUID NOT NULL REFERENCES erp_invoices(id) ON DELETE CASCADE,
    line_no        INT NOT NULL,
    order_line_id  UUID,
    product_id     UUID NOT NULL,
    sku            VARCHAR(50) NOT NULL,
    name           VARCHAR(255) NOT NULL,
    description    TEXT,
    quantity       NUMERIC(15,4) NOT NULL,
    uom            VARCHAR(20) NOT NULL,
    unit_price     NUMERIC(15,2) NOT NULL,
    discount       NUMERIC(15,2) DEFAULT 0,
    tax_type       VARCHAR(20),
    tax_rate       NUMERIC(5,4),
    subtotal       NUMERIC(15,2) NOT NULL,
    tax_amount     NUMERIC(15,2) DEFAULT 0,
    line_total     NUMERIC(15,2) NOT NULL,
    created_at     TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_erp_inv_lines_inv ON erp_invoice_lines(invoice_id, line_no);

-- ============================================================
-- PAYMENTS
-- ============================================================
CREATE TABLE IF NOT EXISTS erp_payments (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL,
    payment_no    VARCHAR(30) NOT NULL,
    direction     VARCHAR(5) NOT NULL,
    customer_id   UUID,
    supplier_id   UUID,
    method        VARCHAR(30) NOT NULL,
    status        VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    amount        NUMERIC(15,2) NOT NULL,
    allocated     NUMERIC(15,2) DEFAULT 0,
    unallocated   NUMERIC(15,2) NOT NULL,
    currency      VARCHAR(3) NOT NULL DEFAULT 'THB',
    fx_rate       NUMERIC(15,6) DEFAULT 1,
    base_amount   NUMERIC(15,2) DEFAULT 0,
    payment_date  TIMESTAMP NOT NULL,
    value_date    TIMESTAMP,
    reference     VARCHAR(100),
    bank_account  VARCHAR(50),
    notes         TEXT,
    metadata      JSONB DEFAULT '{}'::jsonb,
    created_by    UUID,
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW(),
    confirmed_at  TIMESTAMP,
    reversed_at   TIMESTAMP,
    refunded_at   TIMESTAMP,
    CONSTRAINT uq_erp_pay_tenant_no UNIQUE (tenant_id, payment_no)
);
CREATE INDEX idx_erp_pay_tenant_dir    ON erp_payments(tenant_id, direction);
CREATE INDEX idx_erp_pay_status        ON erp_payments(status);
CREATE INDEX idx_erp_pay_customer      ON erp_payments(customer_id) WHERE customer_id IS NOT NULL;
CREATE INDEX idx_erp_pay_supplier      ON erp_payments(supplier_id) WHERE supplier_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS erp_payment_allocations (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id   UUID NOT NULL REFERENCES erp_payments(id) ON DELETE CASCADE,
    invoice_id   UUID NOT NULL REFERENCES erp_invoices(id),
    amount       NUMERIC(15,2) NOT NULL,
    currency     VARCHAR(3) NOT NULL,
    fx_rate      NUMERIC(15,6) DEFAULT 1,
    base_amount  NUMERIC(15,2) DEFAULT 0,
    created_at   TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_erp_alloc_payment ON erp_payment_allocations(payment_id);
CREATE INDEX idx_erp_alloc_invoice ON erp_payment_allocations(invoice_id);

-- ============================================================
-- CREDIT NOTES
-- ============================================================
CREATE TABLE IF NOT EXISTS erp_credit_notes (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL,
    note_no             VARCHAR(30) NOT NULL,
    note_type           VARCHAR(10) NOT NULL,
    invoice_id          UUID,
    customer_id         UUID,
    supplier_id         UUID,
    reason              VARCHAR(30) NOT NULL,
    reason_detail       TEXT,
    issue_date          TIMESTAMP NOT NULL,
    currency            VARCHAR(3) NOT NULL DEFAULT 'THB',
    subtotal            NUMERIC(15,2) NOT NULL,
    tax_amount          NUMERIC(15,2) DEFAULT 0,
    total_amount        NUMERIC(15,2) NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    applied_to_invoice  BOOLEAN DEFAULT FALSE,
    applied_at          TIMESTAMP,
    refunded            BOOLEAN DEFAULT FALSE,
    refunded_at         TIMESTAMP,
    notes               TEXT,
    metadata            JSONB DEFAULT '{}'::jsonb,
    created_by          UUID,
    created_at          TIMESTAMP DEFAULT NOW(),
    updated_at          TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_erp_cn_tenant_no UNIQUE (tenant_id, note_no)
);
CREATE INDEX idx_erp_cn_invoice ON erp_credit_notes(invoice_id) WHERE invoice_id IS NOT NULL;

-- ============================================================
-- SUPPLIERS
-- ============================================================
CREATE TABLE IF NOT EXISTS erp_suppliers (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id             UUID NOT NULL,
    code                  VARCHAR(30) NOT NULL,
    name                  VARCHAR(255) NOT NULL,
    legal_name            VARCHAR(255),
    tax_id                VARCHAR(20),
    branch                VARCHAR(20),
    email                 VARCHAR(255),
    phone                 VARCHAR(30),
    contact_name          VARCHAR(255),
    address               JSONB,
    billing_addr          JSONB,
    shipping_addr         JSONB,
    status                VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    payment_terms         VARCHAR(50),
    credit_limit          NUMERIC(15,2) DEFAULT 0,
    currency              VARCHAR(3) DEFAULT 'THB',
    on_time_delivery_pct  NUMERIC(5,2) DEFAULT 0,
    quality_score         NUMERIC(5,2) DEFAULT 0,
    avg_lead_time_days    INT DEFAULT 0,
    notes                 TEXT,
    metadata              JSONB DEFAULT '{}'::jsonb,
    tags                  JSONB DEFAULT '[]'::jsonb,
    is_active             BOOLEAN DEFAULT TRUE,
    created_by            UUID,
    created_at            TIMESTAMP DEFAULT NOW(),
    updated_at            TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_erp_sup_tenant_code UNIQUE (tenant_id, code)
);
CREATE INDEX idx_erp_sup_status ON erp_suppliers(status);

-- ============================================================
-- JOURNAL ENTRIES (Double-entry)
-- ============================================================
CREATE TABLE IF NOT EXISTS erp_journal_entries (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL,
    entry_no      VARCHAR(30) NOT NULL,
    entry_date    TIMESTAMP NOT NULL,
    description   TEXT,
    total_debit   NUMERIC(15,2) DEFAULT 0,
    total_credit  NUMERIC(15,2) DEFAULT 0,
    currency      VARCHAR(3) NOT NULL DEFAULT 'THB',
    source_type   VARCHAR(30),
    source_id     UUID,
    status        VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    posted_at     TIMESTAMP,
    posted_by     UUID,
    voided_at     TIMESTAMP,
    void_reason   TEXT,
    created_by    UUID,
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_erp_je_tenant_no UNIQUE (tenant_id, entry_no)
);
CREATE INDEX idx_erp_je_tenant_date ON erp_journal_entries(tenant_id, entry_date DESC);
CREATE INDEX idx_erp_je_source      ON erp_journal_entries(source_type, source_id);

CREATE TABLE IF NOT EXISTS erp_journal_lines (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    journal_id    UUID NOT NULL REFERENCES erp_journal_entries(id) ON DELETE CASCADE,
    line_no       INT NOT NULL,
    account_code  VARCHAR(10) NOT NULL,
    account_name  VARCHAR(255),
    debit         NUMERIC(15,2) DEFAULT 0,
    credit        NUMERIC(15,2) DEFAULT 0,
    description   TEXT,
    ref_type      VARCHAR(30),
    ref_id        UUID,
    CONSTRAINT ck_debit_or_credit CHECK (
        (debit > 0 AND credit = 0) OR (credit > 0 AND debit = 0)
    )
);
CREATE INDEX idx_erp_je_lines_journal ON erp_journal_lines(journal_id, line_no);
CREATE INDEX idx_erp_je_lines_account ON erp_journal_lines(account_code);

-- ============================================================
-- DOCUMENT SEQUENCES
-- ============================================================
CREATE TABLE IF NOT EXISTS erp_document_sequences (
    tenant_id   UUID NOT NULL,
    doc_type    VARCHAR(10) NOT NULL,
    year        INT NOT NULL,
    last_seq    INT NOT NULL DEFAULT 0,
    updated_at  TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (tenant_id, doc_type, year)
);

-- ============================================================
-- AUDIT LOGS
-- ============================================================
CREATE TABLE IF NOT EXISTS erp_audit_logs (
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
CREATE INDEX idx_erp_audit_tenant_action ON erp_audit_logs(tenant_id, action, created_at DESC);
CREATE INDEX idx_erp_audit_entity        ON erp_audit_logs(entity_type, entity_id, created_at DESC);

-- ============================================================
-- DEFAULT WAREHOUSE (per tenant — created by app)
-- ============================================================
-- App creates via migration or seed per tenant
```

### D.6 .env

```env
# erp module
ERP_DEFAULT_CURRENCY=THB
ERP_DEFAULT_TAX_RATE=0.07
ERP_DEFAULT_PAYMENT_TERMS=NET_30
ERP_OVERDUE_JOB_CRON="0 0 9 * * *"
ERP_LOW_STOCK_JOB_CRON="0 0 */4 * * *"
ERP_REORDER_AUTO_ENABLED=true
ERP_AGING_CACHE_TTL_SEC=300
ERP_INVOICE_FILE_TTL_DAYS=90
ERP_MAX_LINES_PER_ORDER=500
```

### D.7 DDD Validation Checklist

- [x] **10 Aggregate Roots**: `Product`, `PriceList`, `Warehouse`, `InventoryItem`, `Order`, `Invoice`, `Payment`, `CreditNote`, `Supplier`, `JournalEntry`
- [x] **Entities**: `OrderLine`, `InvoiceLine`, `PaymentAllocation`, `StockMovement`, `StockReservation`, `JournalLine`, `TaxRate`, `PriceListEntry`
- [x] **16 Value Objects** ครบ
- [x] **State machines**: `OrderStatus`, `InvoiceStatus`, `PaymentStatus` (ทุกตัวมี `CanTransitionTo`)
- [x] **Document Numbering** — atomic sequence per (tenant, doc_type, year)
- [x] **Double-entry Accounting** — journal entry + chart of accounts + balance check
- [x] **Stock Movements Ledger** — immutable + snapshot balance
- [x] **Multi-warehouse** — transfer + reservation + consume
- [x] **Multi-currency** — FX rate snapshot
- [x] **Tax Engine** — VAT inclusive/exclusive + WHT + reverse charge
- [x] **Partial fulfillment** — ship/receive/invoice per line
- [x] **Payment Allocation** — FIFO + partial payment
- [x] **Aging Report** — bucket by 30/60/90 days
- [x] **Reorder Point Automation** — ScanLowStock + CreateDraftPO
- [x] **Idempotency** — stock movement ledger prevents double-apply
- [x] **Multi-tenant** — `tenant_id` ทุกตาราง
- [x] **Domain errors**: 60+ sentinel errors
- [x] **Unit tests**: order lifecycle, inventory movement, tax calculation
- [x] **Scheduler**: overdue invoice + low stock
- [x] **Kafka**: 20+ topics (product, order, invoice, payment, inventory)
- [x] **Outbound ports**: `PaymentPort`, `CustomerPort`, `NotifierPort`, `PDFRendererPort`, `CodeGeneratorPort`
- [x] **Import whitelist**: `pkg/kafka, pkg/logger, pkg/transaction` ✅

---

# ✅ PART 4 (erp) — เสร็จสมบูรณ์

**สถิติ:**
- ไฟล์: **~60 ไฟล์**
- Domain: 30 (16 VOs, 15 entities, 11 repos, 7 services + 5 ports, 5 events, 1 error)
- Application: 15 (12 use cases + DTO + mappers)
- Infrastructure: 10 (PG models + repo + seq gen, Kafka, 2 schedulers, cache)
- Interface: 5 (handlers + routes + errors)
- Migration: 13 tables

**Pattern พิเศษ:**
1. ✅ **Document Numbering** — atomic sequence
2. ✅ **Double-entry Accounting** — automatic journal from business events
3. ✅ **Stock Ledger** — event-sourced + snapshot balance
4. ✅ **Multi-warehouse** — transfer + reservation
5. ✅ **Tax Engine** — VAT/WHT/Reverse Charge
6. ✅ **Partial Fulfillment** — per-line tracking
7. ✅ **Payment Allocation** — FIFO + partial
8. ✅ **Aging Report** — 30/60/90 buckets
9. ✅ **Reorder Automation** — auto-PO
10. ✅ **Multi-currency** — FX snapshot
11. ✅ **Weighted Average Cost** — real-time avg cost
12. ✅ **Supplier Performance** — on-time %, quality score

---

# 🎯 ความคืบหน้า (ลำดับ B)

| # | งาน | สถานะ |
|:-:|---|:-:|
| 1 | `packagecatalog` deep dive | ✅ เสร็จ |
| 2 | `erp` deep dive | ✅ เสร็จ (response นี้) |
| 3 | UML / Sequence Diagrams | ⏳ ถัดไป |
| 4 | Sample Data / Fixtures | ⏳ |
| 5 | Postman Collection | ⏳ |
| 6 | Executive Summary | ⏳ |

---

# 📋 Response ถัดไป: **PART 8 — UML / Sequence Diagrams**

จะมี:
- **System Context Diagram** — modules + external systems
- **Sequence Diagrams**:
  1. Customer Onboarding (customer → package → shipment → install → device)
  2. Telemetry Hot Path (MQTT → Kafka → InfluxDB + WS)
  3. Order-to-Cash (ERP): SO → Ship → Invoice → Payment
  4. Procure-to-Pay (ERP): PO → Receive → Bill → Pay
  5. Report Generation (schedule → query → render → deliver)
  6. Device Provisioning + Shadow Sync
- **State Diagrams**:
  - Order lifecycle
  - Subscription lifecycle
  - Device lifecycle
  - Invoice lifecycle
  - Shipment lifecycle
- **ERD** (Entity Relationship Diagram) ของ 47 tables
- **Kafka Topic Map** — producers → consumers matrix
- **Deployment Diagram** — services + infra

 