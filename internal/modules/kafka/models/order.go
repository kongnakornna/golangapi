package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProductID string    `gorm:"not null"`
	Quantity  int       `gorm:"not null"`
	Status    string    `gorm:"default:PENDING"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateOrderRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}
