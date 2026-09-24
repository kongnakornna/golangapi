package presenter

import (
	"encoding/json"

	"icmongolang/pkg/helpers"

	"github.com/google/uuid"
)

// WosOrderRequest represents a customer order creation payload.
// คำขอสร้างคำสั่งซื้อ
type WosOrderRequest struct {
	CustomerName  string          `json:"customerName" validate:"required"`
	CustomerEmail string          `json:"customerEmail" validate:"required,email"`
	CustomerPhone string          `json:"customerPhone"`
	Items         json.RawMessage `json:"items" validate:"required" swaggertype:"object"`
	TotalAmount   float64         `json:"totalAmount" validate:"required"`
	Notes         string          `json:"notes"`
}

// WosOrderResponse represents an order record returned to the client.
// ข้อมูลคำสั่งซื้อที่ส่งกลับไปยังผู้ใช้
type WosOrderResponse struct {
	ID            uuid.UUID         `json:"id"`
	OrderNumber   string            `json:"orderNumber"`
	CustomerName  string            `json:"customerName"`
	CustomerEmail string            `json:"customerEmail"`
	CustomerPhone *string           `json:"customerPhone,omitempty"`
	Items         json.RawMessage   `json:"items" swaggertype:"object"`
	TotalAmount   float64           `json:"totalAmount"`
	Status        string            `json:"status"`
	Notes         *string           `json:"notes,omitempty"`
	CreatedAt     helpers.LocalTime `json:"createdAt"`
	UpdatedAt     helpers.LocalTime `json:"updatedAt"`
}
