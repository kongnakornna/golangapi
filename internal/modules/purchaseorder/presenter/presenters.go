package presenter

import (
	"time"

	"icmongolang/pkg/helpers"

	"github.com/google/uuid"
)

type PurchaseOrderCreateRequest struct {
	QuotationID          *uuid.UUID                   `json:"quotation_id" example:"00000000-0000-0000-0000-000000000000"`
	JobID                *uuid.UUID                   `json:"job_id" example:"00000000-0000-0000-0000-000000000000"`
	SupplierID           uuid.UUID                    `json:"supplier_id" validate:"required" example:"00000000-0000-0000-0000-000000000000"`
	ExpectedDeliveryDate *string                      `json:"expected_delivery_date" example:"2026-07-20T00:00:00Z"`
	Currency             string                       `json:"currency" example:"THB"`
	ExchangeRate         float64                      `json:"exchange_rate" example:"1.0000"`
	ShippingCost         float64                      `json:"shipping_cost" example:"0"`
	PaymentTerms         string                       `json:"payment_terms" example:"Net 30"`
	DeliveryAddress      string                       `json:"delivery_address"`
	Notes                string                       `json:"notes"`
	TermsAndConditions   string                       `json:"terms_and_conditions"`
	TaxRate              float64                      `json:"tax_rate" example:"7.00"`
	DiscountType         string                       `json:"discount_type" example:"PERCENTAGE"`
	DiscountValue        float64                      `json:"discount_value" example:"0"`
	Items                []PurchaseOrderDetailRequest `json:"items" validate:"required,min=1"`
}

type PurchaseOrderDetailRequest struct {
	PartID          uuid.UUID `json:"part_id" validate:"required" example:"00000000-0000-0000-0000-000000000000"`
	QuantityOrdered int       `json:"quantity_ordered" validate:"required,min=1" example:"10"`
	UnitPrice       float64   `json:"unit_price" validate:"required" example:"150.00"`
	Discount        float64   `json:"discount" example:"0"`
	Note            string    `json:"note"`
}

type PurchaseOrderUpdateRequest struct {
	ExpectedDeliveryDate *string                      `json:"expected_delivery_date"`
	Currency             string                       `json:"currency"`
	ExchangeRate         *float64                     `json:"exchange_rate"`
	ShippingCost         *float64                     `json:"shipping_cost"`
	PaymentTerms         string                       `json:"payment_terms"`
	DeliveryAddress      string                       `json:"delivery_address"`
	Notes                string                       `json:"notes"`
	TermsAndConditions   string                       `json:"terms_and_conditions"`
	TaxRate              *float64                     `json:"tax_rate"`
	DiscountType         string                       `json:"discount_type"`
	DiscountValue        *float64                     `json:"discount_value"`
	Items                []PurchaseOrderDetailRequest `json:"items"`
}

type PurchaseOrderStatusRequest struct {
	Reason string `json:"reason" example:"ขอเปลี่ยนสถานะใบสั่งซื้อ"`
}

type PurchaseOrderReceiveRequest struct {
	Items []ReceiveItemRequest `json:"items" validate:"required,min=1"`
}

type ReceiveItemRequest struct {
	DetailID         uuid.UUID `json:"detail_id" validate:"required" example:"00000000-0000-0000-0000-000000000000"`
	ReceivedQuantity int       `json:"received_quantity" validate:"required,min=1" example:"5"`
}

type PurchaseOrderDetailRequestDTO struct {
	ID               uuid.UUID `json:"id,omitempty"`
	PartID           uuid.UUID `json:"part_id"`
	QuantityOrdered  int       `json:"quantity_ordered"`
	QuantityReceived int       `json:"quantity_received"`
	UnitPrice        float64   `json:"unit_price"`
	TotalPrice       float64   `json:"total_price"`
	Discount         float64   `json:"discount"`
	NetPrice         float64   `json:"net_price"`
	Note             string    `json:"note"`
}

type PurchaseOrderResponse struct {
	ID                   uuid.UUID                     `json:"id"`
	PONo                 string                        `json:"po_no"`
	QuotationID          *uuid.UUID                    `json:"quotation_id,omitempty"`
	JobID                *uuid.UUID                    `json:"job_id,omitempty"`
	SupplierID           uuid.UUID                     `json:"supplier_id"`
	PODate               time.Time                     `json:"po_date"`
	ExpectedDeliveryDate *time.Time                    `json:"expected_delivery_date,omitempty"`
	ActualDeliveryDate   *time.Time                    `json:"actual_delivery_date,omitempty"`
	Status               string                        `json:"status"`
	Subtotal             float64                       `json:"subtotal"`
	TaxRate              float64                       `json:"tax_rate"`
	TaxAmount            float64                       `json:"tax_amount"`
	DiscountType         string                        `json:"discount_type,omitempty"`
	DiscountValue        float64                       `json:"discount_value"`
	Total                float64                       `json:"total"`
	Currency             string                        `json:"currency"`
	ExchangeRate         float64                       `json:"exchange_rate"`
	ShippingCost         float64                       `json:"shipping_cost"`
	PaymentTerms         string                        `json:"payment_terms,omitempty"`
	DeliveryAddress      string                        `json:"delivery_address,omitempty"`
	Notes                string                        `json:"notes,omitempty"`
	TermsAndConditions   string                        `json:"terms_and_conditions,omitempty"`
	SentAt               *time.Time                    `json:"sent_at,omitempty"`
	ConfirmedAt          *time.Time                    `json:"confirmed_at,omitempty"`
	ReceivedBy           *uuid.UUID                    `json:"received_by,omitempty"`
	CreatedAt            helpers.LocalTime             `json:"created_at"`
	UpdatedAt            *helpers.LocalTime            `json:"updated_at,omitempty"`
	UserID               uuid.UUID                     `json:"user_id"`
	WhitelabelID         uuid.UUID                     `json:"whitelabel_id"`
	Items                []PurchaseOrderDetailResponse `json:"items,omitempty"`
}

type PurchaseOrderDetailResponse struct {
	ID               uuid.UUID `json:"id"`
	PartID           uuid.UUID `json:"part_id"`
	QuantityOrdered  int       `json:"quantity_ordered"`
	QuantityReceived int       `json:"quantity_received"`
	UnitPrice        float64   `json:"unit_price"`
	TotalPrice       float64   `json:"total_price"`
	Discount         float64   `json:"discount"`
	NetPrice         float64   `json:"net_price"`
	Note             string    `json:"note,omitempty"`
}

type PurchaseOrderDetailResponseDTO struct {
	ID               uuid.UUID `json:"id"`
	PoHeaderID       uuid.UUID `json:"po_header_id"`
	PartID           uuid.UUID `json:"part_id"`
	QuantityOrdered  int       `json:"quantity_ordered"`
	QuantityReceived int       `json:"quantity_received"`
	UnitPrice        float64   `json:"unit_price"`
	TotalPrice       float64   `json:"total_price"`
	Discount         float64   `json:"discount"`
	NetPrice         float64   `json:"net_price"`
	Note             string    `json:"note,omitempty"`
}

type PurchaseOrderSuggestionDTO struct {
	PartID        uuid.UUID `json:"part_id"`
	PartName      string    `json:"part_name"`
	PartCode      string    `json:"part_code"`
	SuggestedQty  int       `json:"suggested_qty"`
	CurrentStock  int       `json:"current_stock"`
	UnitPrice     float64   `json:"unit_price"`
	FromQuotation bool      `json:"from_quotation"`
}

type PurchaseOrderStatusHistoryResponse struct {
	ID         uuid.UUID `json:"id"`
	PoHeaderID uuid.UUID `json:"po_header_id"`
	FromStatus string    `json:"from_status,omitempty"`
	ToStatus   string    `json:"to_status"`
	ChangedBy  uuid.UUID `json:"changed_by"`
	ChangedAt  time.Time `json:"changed_at"`
	Reason     string    `json:"reason,omitempty"`
}

type PaginatedPurchaseOrdersResponse struct {
	Items      []*PurchaseOrderResponse `json:"items"`
	Total      int64                    `json:"total"`
	Page       int                      `json:"page"`
	PerPage    int                      `json:"per_page"`
	TotalPages int                      `json:"total_pages"`
}
