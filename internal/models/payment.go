package models

import (
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// Payment Models (โมเดลการชำระเงิน)
// ============================================================================

// PaymentMethod stores payment method master data.
// วิธีการชำระเงิน (เงินสด, โอนเงิน, บัตรเครดิต, เช็ค, พร้อมเพย์)
type PaymentMethod struct {
	ID               uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	MethodCode       string     `gorm:"column:method_code;type:varchar(20);unique;not null"`
	MethodName       string     `gorm:"column:method_name;type:varchar(100);not null"`
	MethodNameEn     *string    `gorm:"column:method_name_en;type:varchar(100)"`
	IsActive         bool       `gorm:"column:is_active;default:true"`
	RequiresApproval bool       `gorm:"column:requires_approval;default:false"`
	FeePercentage    *float64   `gorm:"column:fee_percentage;type:decimal(5,2)"`
	FeeFixed         *float64   `gorm:"column:fee_fixed;type:decimal(15,2)"`
	Description      *string    `gorm:"column:description;type:text"`
	CreatedAt        time.Time  `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt        *time.Time `gorm:"column:updated_at"`
	UserID           uuid.UUID  `gorm:"column:user_id;type:uuid;not null"`
	WhitelabelID     uuid.UUID  `gorm:"column:whitelabel_id;type:uuid;not null"`
}

func (PaymentMethod) TableName() string { return "m_payment_method" }

// Payment records payment transactions from customers.
// การชำระเงินจากลูกค้า เชื่อมโยงกับ Invoice
type Payment struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	PaymentNo       string     `gorm:"column:payment_no;type:varchar(20);unique;not null"`
	InvoiceID       uuid.UUID  `gorm:"column:invoice_id;type:uuid;not null"`
	JobID           *uuid.UUID `gorm:"column:job_id;type:uuid"`
	CustomerID      uuid.UUID  `gorm:"column:customer_id;type:uuid;not null"`
	PaymentDate     time.Time  `gorm:"column:payment_date;not null;default:now()"`
	PaymentMethodID uuid.UUID  `gorm:"column:payment_method_id;type:uuid;not null"`
	Amount          float64    `gorm:"column:amount;type:decimal(15,2);not null"`
	AmountReceived  float64    `gorm:"column:amount_received;type:decimal(15,2);not null"`
	ChangeAmount    float64    `gorm:"column:change_amount;type:decimal(15,2);default:0"`
	Currency        string     `gorm:"column:currency;type:varchar(10);default:'THB'"`
	ExchangeRate    float64    `gorm:"column:exchange_rate;type:decimal(10,4);default:1.0000"`
	Status          string     `gorm:"column:status;type:varchar(20);not null;default:'PENDING'"`
	ReferenceNumber *string    `gorm:"column:reference_number;type:varchar(50)"`
	BankName        *string    `gorm:"column:bank_name;type:varchar(100)"`
	ChequeNumber    *string    `gorm:"column:cheque_number;type:varchar(50)"`
	ChequeBank      *string    `gorm:"column:cheque_bank;type:varchar(100)"`
	ChequeDate      *time.Time `gorm:"column:cheque_date;type:date"`
	Notes           *string    `gorm:"column:notes;type:text"`
	ReceivedBy      uuid.UUID  `gorm:"column:received_by;type:uuid;not null"`
	ApprovedBy      *uuid.UUID `gorm:"column:approved_by;type:uuid"`
	ApprovedAt      *time.Time `gorm:"column:approved_at"`
	RefundedAmount  float64    `gorm:"column:refunded_amount;type:decimal(15,2);default:0"`
	RefundedAt      *time.Time `gorm:"column:refunded_at"`
	CreatedAt       time.Time  `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt       *time.Time `gorm:"column:updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at"`
	Deleted         bool       `gorm:"column:deleted;default:false"`
	UserID          uuid.UUID  `gorm:"column:user_id;type:uuid;not null"`
	WhitelabelID    uuid.UUID  `gorm:"column:whitelabel_id;type:uuid;not null"`
}

func (Payment) TableName() string { return "t_payment" }

// Receipt stores receipt documents issued for payments.
// ใบเสร็จรับเงินที่ออกให้ลูกค้า
type Receipt struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	ReceiptNo       string     `gorm:"column:receipt_no;type:varchar(20);unique;not null"`
	PaymentID       uuid.UUID  `gorm:"column:payment_id;type:uuid;not null"`
	InvoiceID       uuid.UUID  `gorm:"column:invoice_id;type:uuid;not null"`
	CustomerID      uuid.UUID  `gorm:"column:customer_id;type:uuid;not null"`
	ReceiptDate     time.Time  `gorm:"column:receipt_date;not null;default:now()"`
	ReceiptType     string     `gorm:"column:receipt_type;type:varchar(20);default:'FULL'"`
	Amount          float64    `gorm:"column:amount;type:decimal(15,2);not null"`
	AmountInWordsTh *string    `gorm:"column:amount_in_words_th;type:text"`
	AmountInWordsEn *string    `gorm:"column:amount_in_words_en;type:text"`
	Currency        string     `gorm:"column:currency;type:varchar(10);default:'THB'"`
	Status          string     `gorm:"column:status;type:varchar(20);default:'ISSUED'"`
	Notes           *string    `gorm:"column:notes;type:text"`
	IssuedBy        uuid.UUID  `gorm:"column:issued_by;type:uuid;not null"`
	CreatedAt       time.Time  `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt       *time.Time `gorm:"column:updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at"`
	Deleted         bool       `gorm:"column:deleted;default:false"`
	UserID          uuid.UUID  `gorm:"column:user_id;type:uuid;not null"`
	WhitelabelID    uuid.UUID  `gorm:"column:whitelabel_id;type:uuid;not null"`
}

func (Receipt) TableName() string { return "t_receipt" }

// PaymentHistory tracks payment status changes over time.
// ประวัติการเปลี่ยนสถานะการชำระเงิน
type PaymentHistory struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	PaymentID    uuid.UUID  `gorm:"column:payment_id;type:uuid;not null"`
	FromStatus   *string    `gorm:"column:from_status;type:varchar(20)"`
	ToStatus     string     `gorm:"column:to_status;type:varchar(20);not null"`
	ChangedBy    uuid.UUID  `gorm:"column:changed_by;type:uuid;not null"`
	ChangedAt    time.Time  `gorm:"column:changed_at;not null;default:now()"`
	Reason       *string    `gorm:"column:reason;type:text"`
	WhitelabelID uuid.UUID  `gorm:"column:whitelabel_id;type:uuid;not null"`
}

func (PaymentHistory) TableName() string { return "t_payment_history" }

// OutstandingBalance stores invoice outstanding balances (summary/dashboard).
// ยอดคงค้างของ Invoice สำหรับ Dashboard
type OutstandingBalance struct {
	ID                uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	InvoiceID         uuid.UUID  `gorm:"column:invoice_id;type:uuid;unique;not null"`
	CustomerID        uuid.UUID  `gorm:"column:customer_id;type:uuid;not null"`
	InvoiceTotal      float64    `gorm:"column:invoice_total;type:decimal(15,2);not null"`
	AmountPaid        float64    `gorm:"column:amount_paid;type:decimal(15,2);default:0"`
	OutstandingAmount float64    `gorm:"column:outstanding_amount;type:decimal(15,2);-:migration;->"` // GENERATED ALWAYS AS
	LastPaymentDate   *time.Time `gorm:"column:last_payment_date"`
	Status            string     `gorm:"column:status;type:varchar(20);default:'OUTSTANDING'"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;not null;default:now()"`
	WhitelabelID      uuid.UUID  `gorm:"column:whitelabel_id;type:uuid;not null"`
}

func (OutstandingBalance) TableName() string { return "t_outstanding_balance" }
