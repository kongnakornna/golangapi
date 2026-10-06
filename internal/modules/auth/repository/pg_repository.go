package repository

import (
	"context"
	"errors"
	"time"

	"icmongolang/internal/models"
	"icmongolang/internal/modules/auth"
	domainErr "icmongolang/internal/modules/auth/domain/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type authPgRepo struct {
	db *gorm.DB
}

func CreateAuthPgRepository(db *gorm.DB) auth.AuthPgRepository {
	return &authPgRepo{db: db}
}

func (r *authPgRepo) StoreRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	rt := &models.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	return r.db.WithContext(ctx).Create(rt).Error
}

func (r *authPgRepo) FindRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErr.ErrInvalidRefreshToken
		}
		return nil, err
	}
	return &rt, nil
}

func (r *authPgRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}

func (r *authPgRepo) DeleteAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}
