package auth

import (
	"context"
	"time"

	"icmongolang/internal/models"

	"github.com/google/uuid"
)

// AuthPgRepository defines the database operations for authentication data.
// AuthPgRepository กำหนดการดำเนินการกับฐานข้อมูลสำหรับข้อมูลการยืนยันตัวตน
type AuthPgRepository interface {
	StoreRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error
	FindRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error
}
