package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/repository"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
	base "icmongolang/internal/repository"
)

type DSARRepoImpl struct {
	base.PgRepo[entity.DSARRequest]
	db *gorm.DB
}

func NewDSARPgRepository(db *gorm.DB) repository.DSARRepository {
	return &DSARRepoImpl{
		PgRepo: base.CreatePgRepo[entity.DSARRequest](db),
		db:     db,
	}
}

func (r *DSARRepoImpl) Save(ctx context.Context, dsar *entity.DSARRequest) error {
	return r.db.WithContext(ctx).Create(dsar).Error
}

func (r *DSARRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.DSARRequest, error) {
	var d entity.DSARRequest
	err := r.db.WithContext(ctx).First(&d, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrDSARNotFound
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DSARRepoImpl) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.DSARRequest, error) {
	var out []entity.DSARRequest
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("requested_at DESC").Find(&out).Error
	return out, err
}

func (r *DSARRepoImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status valueobject.DSARStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.DSARRequest{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *DSARRepoImpl) Update(ctx context.Context, d *entity.DSARRequest) error {
	return r.db.WithContext(ctx).Save(d).Error
}

func (r *DSARRepoImpl) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.DSARRequest{}, "user_id = ?", userID).Error
}
