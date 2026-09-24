package settings

import (
	"context"
)

// SettingsRepositoryI is the storage boundary of the settings module.
// All row payloads are map[string]interface{} keyed by the SELECT column aliases.
type SettingsRepositoryI interface {
	// ListPaginate runs a COUNT + paginated SELECT for the given spec.
	ListPaginate(ctx context.Context, spec ListSpec, q ListQuery) ([]map[string]interface{}, int64, error)
	// RowsAll runs the spec SELECT without LIMIT/OFFSET (PageSize <= 0 disables paging).
	RowsAll(ctx context.Context, spec ListSpec, q ListQuery) ([]map[string]interface{}, error)

	// CreateMap inserts one row built from column->value pairs.
	CreateMap(ctx context.Context, table string, m map[string]interface{}) error
	// GetWhere loads the first matching row into dest (map[string]interface{} or struct pointer).
	GetWhere(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	// ExistsWhere reports whether any row in table matches.
	ExistsWhere(ctx context.Context, table string, query string, args ...interface{}) (bool, error)
	// CountWhere returns the number of rows in table matching query.
	CountWhere(ctx context.Context, table string, query string, args ...interface{}) (int64, error)
	// UpdateFieldsMap applies a partial update (column->value) to matching rows.
	UpdateFieldsMap(ctx context.Context, table string, idCol string, idVal interface{}, fields map[string]interface{}) (int64, error)
	// UpdateAllMap applies a partial update to EVERY row of the table (no WHERE).
	UpdateAllMap(ctx context.Context, table string, fields map[string]interface{}) (int64, error)
	// DeleteWhere removes matching rows and returns the affected count.
	DeleteWhere(ctx context.Context, table string, query string, args ...interface{}) (int64, error)
}
