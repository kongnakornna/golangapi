package repository

import (
	"context"
	"strings"

	"icmongolang/internal/modules/settings"

	"gorm.io/gorm"
)

type SettingsPgRepo struct {
	DB *gorm.DB
}

func CreateSettingsPgRepository(db *gorm.DB) *SettingsPgRepo {
	return &SettingsPgRepo{DB: db}
}

func (r *SettingsPgRepo) listBase(ctx context.Context, spec settings.ListSpec, q settings.ListQuery) *gorm.DB {
	db := r.DB.WithContext(ctx).Table(spec.Table)
	for _, j := range spec.Joins {
		db = db.Joins(j)
	}
	for _, f := range spec.Filters {
		v := q.Value(f.Key)
		if v == "" {
			continue
		}
		if f.Op == "LIKE" {
			db = db.Where(f.Column+" ILIKE ?", "%"+settings.EscapeLike(v)+"%")
		} else if f.Op == ">=" || f.Op == "<=" {
			db = db.Where(f.Column+" "+f.Op+" ?", v)
		} else {
			db = db.Where(f.Column+" = ?", v)
		}
	}
	for _, pc := range spec.ParamConds {
		v := q.Value(pc.Key)
		if v == "" {
			continue
		}
		if cond, ok := pc.Conds[v]; ok {
			db = db.Where(cond)
		}
	}
	for _, f := range spec.FixedFilters {
		if f.Value == "" {
			db = db.Where(f.Column + " IS NULL")
		} else {
			db = db.Where(f.Column+" = ?", f.Value)
		}
	}
	return db
}

func (r *SettingsPgRepo) orderOf(spec settings.ListSpec, q settings.ListQuery) (string, error) {
	if q.Sort != "" {
		o, ok := settings.ParseSort(q.Sort, spec.SortCols, spec.DefaultSort)
		if !ok {
			return "", settings.ErrInvalidSort
		}
		return o, nil
	}
	return spec.DefaultSort, nil
}

func (r *SettingsPgRepo) ListPaginate(ctx context.Context, spec settings.ListSpec, q settings.ListQuery) ([]map[string]interface{}, int64, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 5000 {
		q.PageSize = 1000
	}
	order, err := r.orderOf(spec, q)
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if err := r.listBase(ctx, spec, q).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	items := []map[string]interface{}{}
	err = r.listBase(ctx, spec, q).
		Select(strings.Join(spec.Selects, ", ")).
		Order(order).
		Limit(q.PageSize).
		Offset((q.Page - 1) * q.PageSize).
		Find(&items).Error
	return items, total, err
}

func (r *SettingsPgRepo) RowsAll(ctx context.Context, spec settings.ListSpec, q settings.ListQuery) ([]map[string]interface{}, error) {
	order, err := r.orderOf(spec, q)
	if err != nil {
		return nil, err
	}
	db := r.listBase(ctx, spec, q).
		Select(strings.Join(spec.Selects, ", ")).
		Order(order)
	if q.PageSize > 0 {
		page := q.Page
		if page < 1 {
			page = 1
		}
		db = db.Limit(q.PageSize).Offset((page - 1) * q.PageSize)
	}
	items := []map[string]interface{}{}
	err = db.Find(&items).Error
	return items, err
}

func (r *SettingsPgRepo) CreateMap(ctx context.Context, table string, m map[string]interface{}) error {
	return r.DB.WithContext(ctx).Table(table).Create(m).Error
}

func (r *SettingsPgRepo) GetWhere(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return r.DB.WithContext(ctx).Where(query, args...).First(dest).Error
}

func (r *SettingsPgRepo) ExistsWhere(ctx context.Context, table string, query string, args ...interface{}) (bool, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Table(table).Where(query, args...).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *SettingsPgRepo) CountWhere(ctx context.Context, table string, query string, args ...interface{}) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Table(table).Where(query, args...).Count(&count).Error
	return count, err
}

func (r *SettingsPgRepo) UpdateFieldsMap(ctx context.Context, table string, idCol string, idVal interface{}, fields map[string]interface{}) (int64, error) {
	res := r.DB.WithContext(ctx).Table(table).Where(idCol+" = ?", idVal).Updates(fields)
	return res.RowsAffected, res.Error
}

func (r *SettingsPgRepo) UpdateAllMap(ctx context.Context, table string, fields map[string]interface{}) (int64, error) {
	res := r.DB.WithContext(ctx).Table(table).Session(&gorm.Session{AllowGlobalUpdate: true}).Updates(fields)
	return res.RowsAffected, res.Error
}

func (r *SettingsPgRepo) DeleteWhere(ctx context.Context, table string, query string, args ...interface{}) (int64, error) {
	res := r.DB.WithContext(ctx).Table(table).Where(query, args...).Delete(nil)
	return res.RowsAffected, res.Error
}
