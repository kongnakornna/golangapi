package repository

import (
	"context"

	"icmongolang/internal/modules/users"
	usermodels "icmongolang/internal/modules/users/models"

	"gorm.io/gorm"
)

type accessMenuRepo struct{ db *gorm.DB }

func CreateAccessMenuPgRepository(db *gorm.DB) users.AccessMenuRepository {
	return &accessMenuRepo{db: db}
}

func (r *accessMenuRepo) Create(ctx context.Context, m *usermodels.SdUserAccessMenu) (*usermodels.SdUserAccessMenu, error) {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil { return nil, err }
	return m, nil
}

func (r *accessMenuRepo) Get(ctx context.Context, id int) (*usermodels.SdUserAccessMenu, error) {
	var m usermodels.SdUserAccessMenu
	if err := r.db.WithContext(ctx).First(&m, "user_access_id = ?", id).Error; err != nil { return nil, err }
	return &m, nil
}

func (r *accessMenuRepo) Update(ctx context.Context, m *usermodels.SdUserAccessMenu, values map[string]interface{}) (*usermodels.SdUserAccessMenu, error) {
	if err := r.db.WithContext(ctx).Model(m).Updates(values).Error; err != nil { return nil, err }
	return m, nil
}

func (r *accessMenuRepo) Delete(ctx context.Context, id int) error {
	res := r.db.WithContext(ctx).Delete(&usermodels.SdUserAccessMenu{}, "user_access_id = ?", id)
	if res.Error != nil { return res.Error }
	if res.RowsAffected == 0 { return gorm.ErrRecordNotFound }
	return nil
}

func (r *accessMenuRepo) List(ctx context.Context, f users.AccessMenuFilter) ([]*usermodels.SdUserAccessMenu, int64, error) {
	q := r.db.WithContext(ctx).Model(&usermodels.SdUserAccessMenu{})
	if f.UserTypeID != nil { q = q.Where("user_type_id = ?", *f.UserTypeID) }
	if f.RoleID != nil     { q = q.Where("role_id = ?", *f.RoleID) }
	if f.MenuID != nil     { q = q.Where("menu_id = ?", *f.MenuID) }
	if f.ParentID != nil   { q = q.Where("parent_id = ?", *f.ParentID) }

	var total int64
	if err := q.Count(&total).Error; err != nil { return nil, 0, err }
	if f.Limit <= 0 { f.Limit = 20 }
	if f.Offset < 0 { f.Offset = 0 }
	var out []*usermodels.SdUserAccessMenu
	if err := q.Order("user_access_id DESC").Limit(f.Limit).Offset(f.Offset).Find(&out).Error; err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func (r *accessMenuRepo) GetMenusByRole(ctx context.Context, roleID int) ([]*usermodels.SdUserAccessMenu, error) {
	var out []*usermodels.SdUserAccessMenu
	if err := r.db.WithContext(ctx).Where("role_id = ?", roleID).Order("user_access_id ASC").Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

var _ users.AccessMenuRepository = (*accessMenuRepo)(nil)