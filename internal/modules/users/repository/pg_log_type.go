package repository

import (
	"context"

	"icmongolang/internal/modules/users"
	usermodels "icmongolang/internal/modules/users/models"

	"gorm.io/gorm"
)

type logTypeRepo struct{ db *gorm.DB }

func CreateLogTypePgRepository(db *gorm.DB) users.LogTypeRepository {
	return &logTypeRepo{db: db}
}

func (r *logTypeRepo) Create(ctx context.Context, lt *usermodels.SdUserLogType) (*usermodels.SdUserLogType, error) {
	if err := r.db.WithContext(ctx).Create(lt).Error; err != nil { return nil, err }
	return lt, nil
}

func (r *logTypeRepo) Get(ctx context.Context, id int) (*usermodels.SdUserLogType, error) {
	var lt usermodels.SdUserLogType
	if err := r.db.WithContext(ctx).First(&lt, "log_type_id = ?", id).Error; err != nil { return nil, err }
	return &lt, nil
}

func (r *logTypeRepo) Update(ctx context.Context, lt *usermodels.SdUserLogType, values map[string]interface{}) (*usermodels.SdUserLogType, error) {
	if err := r.db.WithContext(ctx).Model(lt).Updates(values).Error; err != nil { return nil, err }
	return lt, nil
}

func (r *logTypeRepo) Delete(ctx context.Context, id int) error {
	res := r.db.WithContext(ctx).Delete(&usermodels.SdUserLogType{}, "log_type_id = ?", id)
	if res.Error != nil { return res.Error }
	if res.RowsAffected == 0 { return gorm.ErrRecordNotFound }
	return nil
}

func (r *logTypeRepo) List(ctx context.Context, limit, offset int) ([]*usermodels.SdUserLogType, int64, error) {
	q := r.db.WithContext(ctx).Model(&usermodels.SdUserLogType{})
	var total int64
	if err := q.Count(&total).Error; err != nil { return nil, 0, err }
	if limit <= 0 { limit = 20 }
	if offset < 0 { offset = 0 }
	var out []*usermodels.SdUserLogType
	if err := q.Order("log_type_id ASC").Limit(limit).Offset(offset).Find(&out).Error; err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

var _ users.LogTypeRepository = (*logTypeRepo)(nil)