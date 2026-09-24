package repository

import (
	"context"

	"icmongolang/internal/repository"
	"icmongolang/internal/template"
	"icmongolang/internal/template/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type templatePgRepo struct {
	repository.PgRepo[models.Template]
	db *gorm.DB
}

func NewTemplatePgRepository(db *gorm.DB) template.PgRepository {
	return &templatePgRepo{
		PgRepo: repository.CreatePgRepo[models.Template](db),
		db:     db,
	}
}

func (r *templatePgRepo) Create(ctx context.Context, exp *models.Template) (*models.Template, error) {
	if err := r.db.WithContext(ctx).Create(exp).Error; err != nil {
		return nil, err
	}
	return exp, nil
}

func (r *templatePgRepo) Get(ctx context.Context, id uuid.UUID) (*models.Template, error) {
	var obj models.Template
	if err := r.db.WithContext(ctx).First(&obj, "id = ?", id.String()).Error; err != nil {
		return nil, err
	}
	return &obj, nil
}

func (r *templatePgRepo) GetMulti(ctx context.Context, limit, offset int) ([]*models.Template, error) {
	var objs []*models.Template
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at DESC").Find(&objs).Error; err != nil {
		return nil, err
	}
	return objs, nil
}

func (r *templatePgRepo) Update(ctx context.Context, exp *models.Template, values map[string]interface{}) (*models.Template, error) {
	if err := r.db.WithContext(ctx).Model(&exp).Updates(values).Error; err != nil {
		return nil, err
	}
	return exp, nil
}

func (r *templatePgRepo) Delete(ctx context.Context, id uuid.UUID) (*models.Template, error) {
	var obj models.Template
	if err := r.db.WithContext(ctx).First(&obj, "id = ?", id.String()).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Delete(&obj).Error; err != nil {
		return nil, err
	}
	return &obj, nil
}

func (r *templatePgRepo) GetByName(ctx context.Context, name string) (*models.Template, error) {
	var obj models.Template
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&obj).Error; err != nil {
		return nil, err
	}
	return &obj, nil
}
