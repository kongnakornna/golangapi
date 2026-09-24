package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	entity "icmongolang/internal/modules/pdpa/domain/entity"
	repository "icmongolang/internal/modules/pdpa/domain/repository"
	base "icmongolang/internal/repository"
)

type pdpaPolicyRow struct {
	ID          uuid.UUID  `gorm:"column:id"`
	PolicyKey   string     `gorm:"column:policy_key"`
	Version     string     `gorm:"column:version"`
	TitleEn     string     `gorm:"column:title_en"`
	TitleTh     string     `gorm:"column:title_th"`
	ContentEn   *string    `gorm:"column:content_en"`
	ContentTh   *string    `gorm:"column:content_th"`
	Status      int64      `gorm:"column:status"`
	EffectiveAt *time.Time `gorm:"column:effective_at"`
	PublishedAt *time.Time `gorm:"column:published_at"`
	CreatedBy   *uuid.UUID `gorm:"column:created_by"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

func (pdpaPolicyRow) TableName() string {
	return "pdpa_privacy_policies"
}

func (r pdpaPolicyRow) toPolicy() entity.PdpaPolicy {
	p := entity.PdpaPolicy{
		ID:          r.ID,
		PolicyKey:   r.PolicyKey,
		Version:     r.Version,
		TitleEn:     r.TitleEn,
		TitleTh:     r.TitleTh,
		ContentEn:   r.ContentEn,
		ContentTh:   r.ContentTh,
		Status:      r.Status,
		EffectiveAt: r.EffectiveAt,
		PublishedAt: r.PublishedAt,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		DeletedAt:   r.DeletedAt,
	}
	if r.CreatedBy != nil {
		p.CreatedBy = stringPtr(r.CreatedBy.String())
	}
	return p
}

func toPolicyRow(p entity.PdpaPolicy) pdpaPolicyRow {
	r := pdpaPolicyRow{
		ID:          p.ID,
		PolicyKey:   p.PolicyKey,
		Version:     p.Version,
		TitleEn:     p.TitleEn,
		TitleTh:     p.TitleTh,
		ContentEn:   p.ContentEn,
		ContentTh:   p.ContentTh,
		Status:      p.Status,
		EffectiveAt: p.EffectiveAt,
		PublishedAt: p.PublishedAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		DeletedAt:   p.DeletedAt,
	}
	if p.CreatedBy != nil {
		id, err := uuid.Parse(*p.CreatedBy)
		if err == nil {
			r.CreatedBy = &id
		}
	}
	return r
}

type policyPgRepository struct {
	base.PgRepo[entity.PdpaPolicy]
	db *gorm.DB
}

func CreatePolicyPgRepository(db *gorm.DB) repository.PolicyRepository {
	return &policyPgRepository{
		PgRepo: base.CreatePgRepo[entity.PdpaPolicy](db),
		db:     db,
	}
}

func (r *policyPgRepository) Save(ctx context.Context, policy *entity.PdpaPolicy) error {
	row := toPolicyRow(*policy)
	return r.db.WithContext(ctx).Create(&row).Error
}

func (r *policyPgRepository) Update(ctx context.Context, policy *entity.PdpaPolicy) error {
	row := toPolicyRow(*policy)
	return r.db.WithContext(ctx).Model(&row).Updates(map[string]interface{}{
		"policy_key":   row.PolicyKey,
		"version":      row.Version,
		"title_en":     row.TitleEn,
		"title_th":     row.TitleTh,
		"content_en":   row.ContentEn,
		"content_th":   row.ContentTh,
		"status":       row.Status,
		"effective_at": row.EffectiveAt,
		"published_at": row.PublishedAt,
		"created_by":   row.CreatedBy,
		"updated_at":   time.Now(),
	}).Where("id = ?", policy.ID.String()).Error
}

func (r *policyPgRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.PdpaPolicy, error) {
	var row pdpaPolicyRow
	if err := r.db.WithContext(ctx).First(&row, "id = ? AND deleted_at IS NULL", id.String()).Error; err != nil {
		return nil, err
	}
	p := row.toPolicy()
	return &p, nil
}

func (r *policyPgRepository) FindActive(ctx context.Context) ([]entity.PdpaPolicy, error) {
	var rows []pdpaPolicyRow
	if err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Where("status = ?", 1).
		Order("effective_at DESC, created_at DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	policies := make([]entity.PdpaPolicy, 0, len(rows))
	for _, row := range rows {
		policies = append(policies, row.toPolicy())
	}
	return policies, nil
}

func (r *policyPgRepository) FindAll(ctx context.Context) ([]entity.PdpaPolicy, error) {
	var rows []pdpaPolicyRow
	if err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	policies := make([]entity.PdpaPolicy, 0, len(rows))
	for _, row := range rows {
		policies = append(policies, row.toPolicy())
	}
	return policies, nil
}

func (r *policyPgRepository) FindByVersion(ctx context.Context, policyKey, version string) (*entity.PdpaPolicy, error) {
	var row pdpaPolicyRow
	if err := r.db.WithContext(ctx).
		First(&row, "policy_key = ? AND version = ? AND deleted_at IS NULL", policyKey, version).Error; err != nil {
		return nil, err
	}
	p := row.toPolicy()
	return &p, nil
}

func (r *policyPgRepository) DeactivateAll(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Model(&pdpaPolicyRow{}).
		Where("deleted_at IS NULL").
		Update("status", 0).Error
}