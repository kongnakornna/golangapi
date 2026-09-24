package presenter

import (
	"time"

	"icmongolang/pkg/helpers"

	"github.com/google/uuid"
)

// PaymentRecordRequest – ข้อมูลที่ใช้บันทึกการชำระเงิน / Payment record request payload.
type PaymentRecordRequest struct {
	InvoiceID       uuid.UUID `json:"invoice_id" validate:"required" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	PaymentMethodID uuid.UUID `json:"payment_method_id" validate:"required" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	Amount          float64   `json:"amount" validate:"required,gt=0" example:"1500.00"`
	AmountReceived  float64   `json:"amount_received" validate:"required,gt=0" example:"2000.00"`
	ChangeAmount    float64   `json:"change_amount" example:"500.00"`
	Currency        string    `json:"currency" example:"THB"`
	ExchangeRate    float64   `json:"exchange_rate" example:"1.0000"`
	ReferenceNumber *string   `json:"reference_number,omitempty" example:"TRF-2026-0001"`
	BankName        *string   `json:"bank_name,omitempty" example:"ธนาคารกรุงเทพ"`
	ChequeNumber    *string   `json:"cheque_number,omitempty" example:"CHQ-001"`
	ChequeBank      *string   `json:"cheque_bank,omitempty" example:"ธนาคารกสิกรไทย"`
	ChequeDate      *string   `json:"cheque_date,omitempty" example:"2026-01-15"`
	Notes           *string   `json:"notes,omitempty" example:"ชำระค่าซ่อมบำรุง"`
}

// PaymentSearchRequest – ตัวกรองสำหรับค้นหาการชำระเงิน / Payment search filter.
type PaymentSearchRequest struct {
	CustomerID      *uuid.UUID `json:"customer_id,omitempty"`
	InvoiceID       *uuid.UUID `json:"invoice_id,omitempty"`
	Status          *string    `json:"status,omitempty" example:"COMPLETED"`
	PaymentMethodID *uuid.UUID `json:"payment_method_id,omitempty"`
	DateFrom        *string    `json:"date_from,omitempty" example:"2026-01-01"`
	DateTo          *string    `json:"date_to,omitempty" example:"2026-12-31"`
	Page            int        `json:"page" example:"1"`
	PerPage         int        `json:"per_page" example:"10"`
}

// RefundRequest – ข้อมูลคำขอคืนเงิน / Refund request payload.
type RefundRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0" example:"500.00"`
	Reason string  `json:"reason" validate:"required" example:"ลูกค้าขอยกเลิกบริการ"`
}

// PaymentResponse – ข้อมูลการชำระเงินที่ส่งกลับ / Payment response payload.
type PaymentResponse struct {
	ID              uuid.UUID         `json:"id,omitempty"`
	PaymentNo       string            `json:"payment_no,omitempty"`
	InvoiceID       uuid.UUID         `json:"invoice_id,omitempty"`
	JobID           *uuid.UUID        `json:"job_id,omitempty"`
	CustomerID      uuid.UUID         `json:"customer_id,omitempty"`
	PaymentDate     time.Time         `json:"payment_date,omitempty"`
	PaymentMethodID uuid.UUID         `json:"payment_method_id,omitempty"`
	Amount          float64           `json:"amount,omitempty"`
	AmountReceived  float64           `json:"amount_received,omitempty"`
	ChangeAmount    float64           `json:"change_amount,omitempty"`
	Currency        string            `json:"currency,omitempty"`
	ExchangeRate    float64           `json:"exchange_rate,omitempty"`
	Status          string            `json:"status,omitempty"`
	ReferenceNumber *string           `json:"reference_number,omitempty"`
	BankName        *string           `json:"bank_name,omitempty"`
	ChequeNumber    *string           `json:"cheque_number,omitempty"`
	ChequeBank      *string           `json:"cheque_bank,omitempty"`
	ChequeDate      *string           `json:"cheque_date,omitempty"`
	Notes           *string           `json:"notes,omitempty"`
	ReceivedBy      uuid.UUID         `json:"received_by,omitempty"`
	ApprovedBy      *uuid.UUID        `json:"approved_by,omitempty"`
	ApprovedAt      *time.Time        `json:"approved_at,omitempty"`
	RefundedAmount  float64           `json:"refunded_amount,omitempty"`
	RefundedAt      *time.Time        `json:"refunded_at,omitempty"`
	CreatedAt       helpers.LocalTime `json:"created_at,omitempty"`
}

// ReceiptResponse – ข้อมูลใบเสร็จที่ส่งกลับ / Receipt response payload.
type ReceiptResponse struct {
	ID              uuid.UUID         `json:"id,omitempty"`
	ReceiptNo       string            `json:"receipt_no,omitempty"`
	PaymentID       uuid.UUID         `json:"payment_id,omitempty"`
	InvoiceID       uuid.UUID         `json:"invoice_id,omitempty"`
	CustomerID      uuid.UUID         `json:"customer_id,omitempty"`
	ReceiptDate     time.Time         `json:"receipt_date,omitempty"`
	ReceiptType     string            `json:"receipt_type,omitempty"`
	Amount          float64           `json:"amount,omitempty"`
	AmountInWordsTh *string           `json:"amount_in_words_th,omitempty"`
	AmountInWordsEn *string           `json:"amount_in_words_en,omitempty"`
	Currency        string            `json:"currency,omitempty"`
	Status          string            `json:"status,omitempty"`
	Notes           *string           `json:"notes,omitempty"`
	IssuedBy        uuid.UUID         `json:"issued_by,omitempty"`
	CreatedAt       helpers.LocalTime `json:"created_at,omitempty"`
}

// PaymentHistoryResponse – ประวัติการชำระเงิน / Payment history response.
type PaymentHistoryResponse struct {
	ID         uuid.UUID `json:"id,omitempty"`
	PaymentID  uuid.UUID `json:"payment_id,omitempty"`
	FromStatus *string   `json:"from_status,omitempty"`
	ToStatus   string    `json:"to_status,omitempty"`
	ChangedBy  uuid.UUID `json:"changed_by,omitempty"`
	ChangedAt  time.Time `json:"changed_at,omitempty"`
	Reason     *string   `json:"reason,omitempty"`
}

// OutstandingBalanceResponse – ยอดคงค้างของลูกค้า / Outstanding balance response.
type OutstandingBalanceResponse struct {
	InvoiceID         uuid.UUID  `json:"invoice_id,omitempty"`
	InvoiceTotal      float64    `json:"invoice_total,omitempty"`
	AmountPaid        float64    `json:"amount_paid,omitempty"`
	OutstandingAmount float64    `json:"outstanding_amount,omitempty"`
	LastPaymentDate   *time.Time `json:"last_payment_date,omitempty"`
	Status            string     `json:"status,omitempty"`
}

// PaginatedPaymentResponse – รายการชำระเงินแบบแบ่งหน้า / Paginated payment list.
type PaginatedPaymentResponse struct {
	Payments   []*PaymentResponse `json:"payments"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PerPage    int                `json:"per_page"`
	TotalPages int                `json:"total_pages"`
}
