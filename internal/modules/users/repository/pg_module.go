package repository

import (
	"context"
	"strings"

	"icmongolang/internal/modules/users"
	usermodels "icmongolang/internal/modules/users/models"

	"gorm.io/gorm"
)

type moduleRepo struct{ db *gorm.DB }

func CreateModulePgRepository(db *gorm.DB) users.ModuleRepository {
	return &moduleRepo{db: db}
}

func (r *moduleRepo) Create(ctx context.Context, m *usermodels.SdModules) (*usermodels.SdModules, error) {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil { return nil, err }
	return m, nil
}

func (r *moduleRepo) Get(ctx context.Context, id int) (*usermodels.SdModules, error) {
	var m usermodels.SdModules
	if err := r.db.WithContext(ctx).First(&m, "module_id = ?", id).Error; err != nil { return nil, err }
	return &m, nil
}

func (r *moduleRepo) Update(ctx context.Context, m *usermodels.SdModules, values map[string]interface{}) (*usermodels.SdModules, error) {
	values["update_date"] = gorm.Expr("now()")
	if err := r.db.WithContext(ctx).Model(m).Updates(values).Error; err != nil { return nil, err }
	return m, nil
}

func (r *moduleRepo) Delete(ctx context.Context, id int) error {
	res := r.db.WithContext(ctx).Delete(&usermodels.SdModules{}, "module_id = ?", id)
	if res.Error != nil { return res.Error }
	if res.RowsAffected == 0 { return gorm.ErrRecordNotFound }
	return nil
}

func (r *moduleRepo) List(ctx context.Context, f users.ModuleFilter) ([]*usermodels.SdModules, int64, error) {
	q := r.db.WithContext(ctx).Model(&usermodels.SdModules{})
	if f.ModuleParentID != nil { q = q.Where("module_parent_id = ?", *f.ModuleParentID) }
	if len(f.Statuses) > 0     { q = q.Where("status IN ?", f.Statuses) }
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
	var out []*usermodels.SdModules
	if err := q.Order(sortBy + " " + sortOrder).Limit(f.Limit).Offset(f.Offset).Find(&out).Error; err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func (r *moduleRepo) BulkUpdateSort(ctx context.Context, items []users.SortItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, it := range items {
			values := map[string]interface{}{
				"sort":        it.Sort,
				"update_date": gorm.Expr("now()"),
			}
			if it.ParentID != nil { values["module_parent_id"] = *it.ParentID }
			if err := tx.Model(&usermodels.SdModules{}).Where("module_id = ?", it.ID).Updates(values).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

var _ users.ModuleRepository = (*moduleRepo)(nil)