package presenter

import (
	"time"

	"icmongolang/pkg/helpers"

	"github.com/google/uuid"
)

type QuotationCreate struct {
	QuotationNo        string    `json:"quotation_no" validate:"required" example:"QTN-2024-0001"`
	JobID              uuid.UUID `json:"job_id" validate:"required"`
	CustomerID         uuid.UUID `json:"customer_id" validate:"required"`
	QuotationDate      string    `json:"quotation_date" example:"2024-01-15T10:00:00Z"`
	ExpiryDate         string    `json:"expiry_date" validate:"required" example:"2024-02-15T10:00:00Z"`
	Subtotal           float64   `json:"subtotal"`
	TaxRate            float64   `json:"tax_rate" example:"7.00"`
	TaxAmount          float64   `json:"tax_amount"`
	DiscountType       *string   `json:"discount_type" example:"PERCENTAGE"`
	DiscountValue      float64   `json:"discount_value"`
	Total              float64   `json:"total"`
	AmountInWordsTh    *string   `json:"amount_in_words_th"`
	AmountInWordsEn    *string   `json:"amount_in_words_en"`
	Currency           *string   `json:"currency" example:"THB"`
	ExchangeRate       float64   `json:"exchange_rate" example:"1.0000"`
	Notes              *string   `json:"notes"`
	TermsAndConditions *string   `json:"terms_and_conditions"`
}

type QuotationResponse struct {
	ID                 uuid.UUID          `json:"id,omitempty"`
	QuotationNo        string             `json:"quotation_no,omitempty"`
	JobID              uuid.UUID          `json:"job_id,omitempty"`
	CustomerID         uuid.UUID          `json:"customer_id,omitempty"`
	QuotationDate      time.Time          `json:"quotation_date,omitempty"`
	ExpiryDate         time.Time          `json:"expiry_date,omitempty"`
	Status             string             `json:"status,omitempty"`
	Subtotal           float64            `json:"subtotal,omitempty"`
	TaxRate            float64            `json:"tax_rate,omitempty"`
	TaxAmount          float64            `json:"tax_amount,omitempty"`
	DiscountType       *string            `json:"discount_type,omitempty"`
	DiscountValue      float64            `json:"discount_value,omitempty"`
	Total              float64            `json:"total,omitempty"`
	AmountInWordsTh    *string            `json:"amount_in_words_th,omitempty"`
	AmountInWordsEn    *string            `json:"amount_in_words_en,omitempty"`
	Currency           *string            `json:"currency,omitempty"`
	ExchangeRate       float64            `json:"exchange_rate,omitempty"`
	Notes              *string            `json:"notes,omitempty"`
	TermsAndConditions *string            `json:"terms_and_conditions,omitempty"`
	ApprovedBy         *uuid.UUID         `json:"approved_by,omitempty"`
	ApprovedAt         *time.Time         `json:"approved_at,omitempty"`
	RejectedReason     *string            `json:"rejected_reason,omitempty"`
	ConvertedToPo      bool               `json:"converted_to_po,omitempty"`
	UserID             uuid.UUID          `json:"user_id,omitempty"`
	WhitelabelID       uuid.UUID          `json:"whitelabel_id,omitempty"`
	CreatedAt          helpers.LocalTime  `json:"created_at,omitempty"`
	UpdatedAt          *helpers.LocalTime `json:"updated_at,omitempty"`
}

type QuotationUpdate struct {
	QuotationDate      *string  `json:"quotation_date" example:"2024-01-15T10:00:00Z"`
	ExpiryDate         *string  `json:"expiry_date" example:"2024-02-15T10:00:00Z"`
	Subtotal           *float64 `json:"subtotal"`
	TaxRate            *float64 `json:"tax_rate" example:"7.00"`
	TaxAmount          *float64 `json:"tax_amount"`
	DiscountType       *string  `json:"discount_type" example:"PERCENTAGE"`
	DiscountValue      *float64 `json:"discount_value"`
	Total              *float64 `json:"total"`
	AmountInWordsTh    *string  `json:"amount_in_words_th"`
	AmountInWordsEn    *string  `json:"amount_in_words_en"`
	Currency           *string  `json:"currency" example:"THB"`
	ExchangeRate       *float64 `json:"exchange_rate" example:"1.0000"`
	Notes              *string  `json:"notes"`
	TermsAndConditions *string  `json:"terms_and_conditions"`
}

type PaginatedQuotationResponse struct {
	Items      []*QuotationResponse `json:"items"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PerPage    int                  `json:"per_page"`
	TotalPages int                  `json:"total_pages"`
}
