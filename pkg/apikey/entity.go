package apikey

import (
	"time"

	"github.com/google/uuid"
)

// Key is the aggregate root for an API key. The raw secret is transient and
// never persisted — only the HMAC fingerprint is stored.
type Key struct {
	ID          uuid.UUID  `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name        string     `gorm:"column:name;not null" json:"name"`
	Description *string    `gorm:"column:description" json:"description,omitempty"`
	Prefix      string     `gorm:"column:prefix;not null;index" json:"prefix"`
	LastFour    string     `gorm:"column:last_four;not null" json:"last_four"`
	HashedKey   string     `gorm:"column:hashed_key;not null;uniqueIndex" json:"-"`
	ExpiresAt   *time.Time `gorm:"column:expires_at" json:"expires_at,omitempty"`
	LastUsedAt  *time.Time `gorm:"column:last_used_at" json:"last_used_at,omitempty"`
	CreatedBy   uuid.UUID  `gorm:"column:created_by;not null;index" json:"created_by"`
	UpdatedBy   uuid.UUID  `gorm:"column:updated_by;not null;index" json:"updated_by"`
	IsActive    bool       `gorm:"column:is_active;not null;default:true;index" json:"is_active"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`

	// Transient — populated only on create/rotate, never persisted.
	PlainKey string `gorm:"-" json:"api_key,omitempty"`
}

// TableName pins the entity.
func (Key) TableName() string { return "api_keys" }