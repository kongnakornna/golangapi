package repository

import (
	"context"
	"time"

	"icmongolang/internal/models"
	"icmongolang/internal/modules/auth"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// authPgRepo implements AuthPgRepository for PostgreSQL using GORM.
// authPgRepo implement อินเทอร์เฟซ AuthPgRepository สำหรับ PostgreSQL ด้วย GORM
type authPgRepo struct {
	db *gorm.DB
}

// CreateAuthPgRepository creates a new AuthPgRepository instance.
// CreateAuthPgRepository สร้างอินสแตนซ์ AuthPgRepository ใหม่
func CreateAuthPgRepository(db *gorm.DB) auth.AuthPgRepository {
	return &authPgRepo{db: db}
}

// StoreRefreshToken saves a new refresh token in the database.
// StoreRefreshToken บันทึก refresh token ใหม่ลงในฐานข้อมูล
func (r *authPgRepo) StoreRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	rt := &models.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	return r.db.WithContext(ctx).Create(rt).Error
}

// FindRefreshToken retrieves a refresh token by its value.
// FindRefreshToken ค้นหา refresh token ตามค่าของ token
func (r *authPgRepo) FindRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&rt).Error; err != nil {
		return nil, err
	}
	return &rt, nil
}

// DeleteRefreshToken removes a refresh token by its value.
// DeleteRefreshToken ลบ refresh token ตามค่าของ token
func (r *authPgRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}

// DeleteAllUserRefreshTokens removes all refresh tokens for a specific user.
// DeleteAllUserRefreshTokens ลบ refresh token ทั้งหมดของผู้ใช้ที่ระบุ
func (r *authPgRepo) DeleteAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}
