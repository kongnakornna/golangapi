// Package models จัดเตรียม GORM models + request/response DTOs ของโมดูล kafka
package models

import (
	"time"

	"github.com/google/uuid"
)

// ============================================================
// GORM model
// ============================================================

// Order — GORM model สำหรับตาราง orders
type Order struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"         example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	ProductID string    `gorm:"not null"             json:"product_id" example:"PROD-001"`
	Quantity  int       `gorm:"not null"             json:"quantity"   example:"3"`
	Status    string    `gorm:"default:PENDING"      json:"status"     example:"PENDING" enums:"PENDING,PROCESSING,COMPLETED,FAILED,CANCELLED"`
	CreatedAt time.Time `json:"created_at" example:"2026-10-04T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2026-10-04T12:00:00Z"`
}

func (Order) TableName() string { return "orders" }

// ============================================================
// Request DTOs
// ============================================================

// CreateOrderRequest — body ของ POST /orders
type CreateOrderRequest struct {
	ProductID string `json:"product_id" validate:"required"       example:"PROD-001"`
	Quantity  int    `json:"quantity"   validate:"required,min=1" example:"3" minimum:"1"`
}

// ============================================================
// Response DTOs
// ============================================================

// OrderResponse — representation ของ Order
type OrderResponse struct {
	ID        string    `json:"id"         example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	ProductID string    `json:"product_id" example:"PROD-001"`
	Quantity  int       `json:"quantity"   example:"3"`
	Status    string    `json:"status"     example:"PENDING"`
	CreatedAt time.Time `json:"created_at" example:"2026-10-04T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2026-10-04T12:00:00Z"`
}

// CreateOrderAcceptedResponse — response 202 Accepted ของ POST /orders
type CreateOrderAcceptedResponse struct {
	OrderID string `json:"order_id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	Status  string `json:"status"   example:"PENDING"`
	Message string `json:"message"  example:"Order accepted for processing"`
}

// ErrorResponse — generic error response
type ErrorResponse struct {
	Error string `json:"error" example:"product_id and quantity (>=1) are required"`
}
