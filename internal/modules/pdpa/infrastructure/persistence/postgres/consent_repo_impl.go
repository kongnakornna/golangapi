package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/repository"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
	base "icmongolang/internal/repository"
)

type ConsentRepoImpl struct {
	base.PgRepo[entity.ConsentLog]
	db *gorm.DB
}

func NewConsentPgRepository(db *gorm.DB) repository.ConsentRepository {
	return &ConsentRepoImpl{
		PgRepo: base.CreatePgRepo[entity.ConsentLog](db),
		db:     db,
	}
}

func (r *ConsentRepoImpl) Save(ctx context.Context, consent *entity.ConsentLog) error {
	return r.db.WithContext(ctx).Create(consent).Error
}

func (r *ConsentRepoImpl) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.ConsentLog, error) {
	var out []entity.ConsentLog
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("consented_at DESC").Find(&out).Error
	return out, err
}

func (r *ConsentRepoImpl) FindLatestByUserAndPurpose(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) (*entity.ConsentLog, error) {
	var c entity.ConsentLog
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND purpose = ?", userID, purpose.String()).
		Order("consented_at DESC").First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrConsentNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ConsentRepoImpl) Revoke(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) error {
	now := time.Now()
	res := r.db.WithContext(ctx).
		Model(&entity.ConsentLog{}).
		Where("user_id = ? AND purpose = ? AND status = ?",
			userID, purpose.String(), valueobject.ConsentGranted).
		Updates(map[string]interface{}{
			"status":     valueobject.ConsentRevoked,
			"revoked_at": now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrConsentNotFound
	}
	return nil
}

func (r *ConsentRepoImpl) FindExpiredRevoked(ctx context.Context, cutoff time.Time) ([]entity.ConsentLog, error) {
	var out []entity.ConsentLog
	err := r.db.WithContext(ctx).
		Where("status = ? AND revoked_at <= ?", valueobject.ConsentRevoked, cutoff).
		Find(&out).Error
	return out, err
}

func (r *ConsentRepoImpl) DeletePermanently(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.ConsentLog{}, "id = ?", id).Error
}

func (r *ConsentRepoImpl) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.ConsentLog{}, "user_id = ?", userID).Error
}

func (r *ConsentRepoImpl) IsConsentActive(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.ConsentLog{}).
		Where("user_id = ? AND purpose = ? AND status = ? AND expires_at > NOW()",
			userID, purpose.String(), valueobject.ConsentGranted).
		Count(&count).Error
	return count > 0, err
}
