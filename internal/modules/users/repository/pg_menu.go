package repository

import (
	"context"
	"strings"

	"icmongolang/internal/modules/users"
	usermodels "icmongolang/internal/modules/users/models"

	"gorm.io/gorm"
)

type menuRepo struct{ db *gorm.DB }

func CreateMenuPgRepository(db *gorm.DB) users.MenuRepository {
	return &menuRepo{db: db}
}

func (r *menuRepo) Create(ctx context.Context, m *usermodels.SdMenu) (*usermodels.SdMenu, error) {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil { return nil, err }
	return m, nil
}

func (r *menuRepo) Get(ctx context.Context, id int) (*usermodels.SdMenu, error) {
	var m usermodels.SdMenu
	if err := r.db.WithContext(ctx).First(&m, "menu_id = ?", id).Error; err != nil { return nil, err }
	return &m, nil
}

func (r *menuRepo) Update(ctx context.Context, m *usermodels.SdMenu, values map[string]interface{}) (*usermodels.SdMenu, error) {
	values["update_date"] = gorm.Expr("now()")
	if err := r.db.WithContext(ctx).Model(m).Updates(values).Error; err != nil { return nil, err }
	return m, nil
}

func (r *menuRepo) Delete(ctx context.Context, id int) error {
	res := r.db.WithContext(ctx).Delete(&usermodels.SdMenu{}, "menu_id = ?", id)
	if res.Error != nil { return res.Error }
	if res.RowsAffected == 0 { return gorm.ErrRecordNotFound }
	return nil
}

func (r *menuRepo) List(ctx context.Context, f users.MenuFilter) ([]*usermodels.SdMenu, int64, error) {
	q := r.db.WithContext(ctx).Model(&usermodels.SdMenu{})
	if f.ModulesID != nil    { q = q.Where("modules_id = ?", *f.ModulesID) }
	if f.MenuParentID != nil { q = q.Where("menu_parent_id = ?", *f.MenuParentID) }
	if len(f.Statuses) > 0   { q = q.Where("status IN ?", f.Statuses) }
	if kw := strings.TrimSpace(f.Keyword); kw != "" {
		like := "%" + kw + "%"
		q = q.Where("name ILIKE ? OR path ILIKE ?", like, like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil { return nil, 0, err }

	sortBy := "sort"
	switch strings.ToLower(f.SortBy) {
	case "name": sortBy = "name"
	case "create_date": sortBy = "create_date"
	case "update_date": sortBy = "update_date"
	}
	sortOrder := "ASC"
	if strings.EqualFold(f.SortOrder, "desc") { sortOrder = "DESC" }

	if f.Limit <= 0 { f.Limit = 20 }
	if f.Offset < 0 { f.Offset = 0 }
	var out []*usermodels.SdMenu
	if err := q.Order(sortBy + " " + sortOrder).Limit(f.Limit).Offset(f.Offset).Find(&out).Error; err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func (r *menuRepo) BulkUpdateSort(ctx context.Context, items []users.SortItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, it := range items {
			values := map[string]interface{}{
				"sort":        it.Sort,
				"update_date": gorm.Expr("now()"),
			}
			if it.ParentID != nil { values["menu_parent_id"] = *it.ParentID }
			if err := tx.Model(&usermodels.SdMenu{}).Where("menu_id = ?", it.ID).Updates(values).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

var _ users.MenuRepository = (*menuRepo)(nil)