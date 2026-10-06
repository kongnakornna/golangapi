package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/repository"
	base "icmongolang/internal/repository"
)

type PolicyRepoImpl struct {
	base.PgRepo[entity.PdpaPolicy]
	db *gorm.DB
}

func CreatePolicyPgRepository(db *gorm.DB) repository.PolicyRepository {
	return &PolicyRepoImpl{
		PgRepo: base.CreatePgRepo[entity.PdpaPolicy](db),
		db:     db,
	}
}

func (r *PolicyRepoImpl) Save(ctx context.Context, p *entity.PdpaPolicy) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *PolicyRepoImpl) Update(ctx context.Context, p *entity.PdpaPolicy) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *PolicyRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.PdpaPolicy, error) {
	var p entity.PdpaPolicy
	err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &p, err
}

func (r *PolicyRepoImpl) FindByVersion(ctx context.Context, policyKey, version string) (*entity.PdpaPolicy, error) {
	var p entity.PdpaPolicy
	err := r.db.WithContext(ctx).
		First(&p, "policy_key = ? AND version = ?", policyKey, version).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &p, err
}

func (r *PolicyRepoImpl) FindActive(ctx context.Context) ([]entity.PdpaPolicy, error) {
	var rows []entity.PdpaPolicy
	err := r.db.WithContext(ctx).
		Where("status = ?", 1).
		Order("effective_at DESC, created_at DESC").
		Find(&rows).Error
	return rows, err
}

func (r *PolicyRepoImpl) FindAll(ctx context.Context) ([]entity.PdpaPolicy, error) {
	var rows []entity.PdpaPolicy
	err := r.db.WithContext(ctx).
		Order("created_at DESC").Find(&rows).Error
	return rows, err
}

func (r *PolicyRepoImpl) DeactivateAll(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Model(&entity.PdpaPolicy{}).
		Where("1=1").
		Update("status", 0).Error
}
