package models

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseOrderHeader struct {
	ID                  uuid.UUID  `gorm:"column:id;type:uuid;default:gen_random_uuid();primary_key"`
	PONo                string     `gorm:"column:po_no;type:varchar(20);unique;not null"`
	QuotationID         *uuid.UUID `gorm:"column:quotation_id;type:uuid"`
	JobID               *uuid.UUID `gorm:"column:job_id;type:uuid"`
	SupplierID          uuid.UUID  `gorm:"column:supplier_id;type:uuid;not null"`
	PODate              time.Time  `gorm:"column:po_date;not null;default:now()"`
	ExpectedDeliveryDate *time.Time `gorm:"column:expected_delivery_date"`
	ActualDeliveryDate  *time.Time `gorm:"column:actual_delivery_date"`
	Status              string     `gorm:"column:status;type:varchar(20);not null;default:DRAFT"`
	Subtotal            float64    `gorm:"column:subtotal;type:decimal(15,2);not null;default:0"`
	TaxRate             float64    `gorm:"column:tax_rate;type:decimal(5,2);default:7.00"`
	TaxAmount           float64    `gorm:"column:tax_amount;type:decimal(15,2);default:0"`
	DiscountType        *string    `gorm:"column:discount_type;type:varchar(20)"`
	DiscountValue       float64    `gorm:"column:discount_value;type:decimal(15,2);default:0"`
	Total               float64    `gorm:"column:total;type:decimal(15,2);not null;default:0"`
	Currency            string     `gorm:"column:currency;type:varchar(10);default:THB"`
	ExchangeRate        float64    `gorm:"column:exchange_rate;type:decimal(10,4);default:1.0000"`
	ShippingCost        float64    `gorm:"column:shipping_cost;type:decimal(15,2);default:0"`
	PaymentTerms        *string    `gorm:"column:payment_terms;type:text"`
	DeliveryAddress     *string    `gorm:"column:delivery_address;type:text"`
	Notes               *string    `gorm:"column:notes;type:text"`
	TermsAndConditions  *string    `gorm:"column:terms_and_conditions;type:text"`
	SentAt              *time.Time `gorm:"column:sent_at"`
	ConfirmedAt         *time.Time `gorm:"column:confirmed_at"`
	ReceivedBy          *uuid.UUID `gorm:"column:received_by;type:uuid"`
	CreatedAt           time.Time  `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt           *time.Time `gorm:"column:updated_at"`
	DeletedAt           *time.Time `gorm:"column:deleted_at"`
	Deleted             bool       `gorm:"column:deleted;default:false"`
	UserID              uuid.UUID  `gorm:"column:user_id;type:uuid;not null"`
	WhitelabelID        uuid.UUID  `gorm:"column:whitelabel_id;type:uuid;not null"`

	Details []PurchaseOrderDetail `gorm:"foreignKey:PoHeaderID"`
}

func (PurchaseOrderHeader) TableName() string { return "t_purchase_order_header" }

type PurchaseOrderDetail struct {
	ID               uuid.UUID  `gorm:"column:id;type:uuid;default:gen_random_uuid();primary_key"`
	PoHeaderID       uuid.UUID  `gorm:"column:po_header_id;type:uuid;not null"`
	PartID           uuid.UUID  `gorm:"column:part_id;type:uuid;not null"`
	QuantityOrdered  int        `gorm:"column:quantity_ordered;not null;default:1"`
	QuantityReceived int        `gorm:"column:quantity_received;default:0"`
	UnitPrice        float64    `gorm:"column:unit_price;type:decimal(15,2);not null"`
	TotalPrice       float64    `gorm:"column:total_price;type:decimal(15,2)"`
	Discount         float64    `gorm:"column:discount;type:decimal(15,2);default:0"`
	NetPrice         float64    `gorm:"column:net_price;type:decimal(15,2)"`
	Note             *string    `gorm:"column:note;type:text"`
	CreatedAt        time.Time  `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt        *time.Time `gorm:"column:updated_at"`
	UserID           uuid.UUID  `gorm:"column:user_id;type:uuid;not null"`
	WhitelabelID     uuid.UUID  `gorm:"column:whitelabel_id;type:uuid;not null"`
}

func (PurchaseOrderDetail) TableName() string { return "t_purchase_order_detail" }

type PurchaseOrderStatusHistory struct {
	ID          uuid.UUID  `gorm:"column:id;type:uuid;default:gen_random_uuid();primary_key"`
	PoHeaderID  uuid.UUID  `gorm:"column:po_header_id;type:uuid;not null"`
	FromStatus  *string    `gorm:"column:from_status;type:varchar(20)"`
	ToStatus    string     `gorm:"column:to_status;type:varchar(20);not null"`
	ChangedBy   uuid.UUID  `gorm:"column:changed_by;type:uuid;not null"`
	ChangedAt   time.Time  `gorm:"column:changed_at;not null;default:now()"`
	Reason      *string    `gorm:"column:reason;type:text"`
	WhitelabelID uuid.UUID `gorm:"column:whitelabel_id;type:uuid;not null"`
}

func (PurchaseOrderStatusHistory) TableName() string { return "t_purchase_order_status_history" }
