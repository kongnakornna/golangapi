package models

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken stores JWT refresh tokens for secure token rotation.
// RefreshToken เก็บ JWT refresh token เพื่อการหมุนเวียน token อย่างปลอดภัย
type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Token     string    `gorm:"type:text;not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
