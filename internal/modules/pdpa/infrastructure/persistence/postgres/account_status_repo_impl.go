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

type AccountStatusRepoImpl struct {
	base.PgRepo[entity.UserAccountStatus]
	db *gorm.DB
}

func NewAccountStatusPgRepository(db *gorm.DB) repository.UserAccountStatusRepository {
	return &AccountStatusRepoImpl{
		PgRepo: base.CreatePgRepo[entity.UserAccountStatus](db),
		db:     db,
	}
}

func (r *AccountStatusRepoImpl) Save(ctx context.Context, s *entity.UserAccountStatus) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *AccountStatusRepoImpl) FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserAccountStatus, error) {
	var s entity.UserAccountStatus
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *AccountStatusRepoImpl) FindSuspendedWithoutDeletion(ctx context.Context, cutoffDate time.Time) ([]entity.UserAccountStatus, error) {
	var out []entity.UserAccountStatus
	err := r.db.WithContext(ctx).
		Where("status = ? AND deletion_confirmed_at IS NULL AND retention_deadline <= ?",
			valueobject.AccountSuspended, cutoffDate).
		Find(&out).Error
	return out, err
}

func (r *AccountStatusRepoImpl) UpdateStatus(ctx context.Context, userID uuid.UUID, status valueobject.AccountStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.UserAccountStatus{}).
		Where("user_id = ?", userID).
		Update("status", status).Error
}

func (r *AccountStatusRepoImpl) UpdateDeletionConfirmed(ctx context.Context, userID uuid.UUID, confirmedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.UserAccountStatus{}).
		Where("user_id = ?", userID).
		Update("deletion_confirmed_at", confirmedAt).Error
}

func (r *AccountStatusRepoImpl) UpdateAutoDeleted(ctx context.Context, userID uuid.UUID, deletedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.UserAccountStatus{}).
		Where("user_id = ?", userID).
		Update("auto_deleted_at", deletedAt).Error
}
