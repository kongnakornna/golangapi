package repository

import (
	"context"

	"icmongolang/internal/modules/users"
	usermodels "icmongolang/internal/modules/users/models"

	"gorm.io/gorm"
)

type fileRepo struct{ db *gorm.DB }

func CreateFilePgRepository(db *gorm.DB) users.FileRepository {
	return &fileRepo{db: db}
}

func (r *fileRepo) Create(ctx context.Context, f *usermodels.SdUserFile) (*usermodels.SdUserFile, error) {
	if err := r.db.WithContext(ctx).Create(f).Error; err != nil {
		return nil, err
	}
	return f, nil
}

func (r *fileRepo) Get(ctx context.Context, id int) (*usermodels.SdUserFile, error) {
	var f usermodels.SdUserFile
	if err := r.db.WithContext(ctx).First(&f, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *fileRepo) Update(ctx context.Context, f *usermodels.SdUserFile, values map[string]interface{}) (*usermodels.SdUserFile, error) {
	if err := r.db.WithContext(ctx).Model(f).Updates(values).Error; err != nil {
		return nil, err
	}
	return f, nil
}

func (r *fileRepo) Delete(ctx context.Context, id int) error {
	res := r.db.WithContext(ctx).Delete(&usermodels.SdUserFile{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *fileRepo) List(ctx context.Context, f users.FileFilter) ([]*usermodels.SdUserFile, int64, error) {
	q := r.db.WithContext(ctx).Model(&usermodels.SdUserFile{})
	if f.UID != nil {
		q = q.Where("uid = ?", *f.UID)
	}
	if f.FileTypeID != nil {
		q = q.Where("file_type_id = ?", *f.FileTypeID)
	}
	if len(f.Statuses) > 0 {
		q = q.Where("status IN ?", f.Statuses)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if f.Limit <= 0 {
		f.Limit = 20
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	var out []*usermodels.SdUserFile
	if err := q.Order("id DESC").Limit(f.Limit).Offset(f.Offset).Find(&out).Error; err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

var _ users.FileRepository = (*fileRepo)(nil)