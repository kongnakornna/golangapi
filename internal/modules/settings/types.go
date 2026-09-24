package settings

import (
	"errors"
	"strings"
)

var ErrInvalidSort = errors.New("invalid sort option")

// Filter describes one optional WHERE clause driven by a query param.
type Filter struct {
	Column string // e.g. "s.setting_name"
	Key    string // query param name: "keyword", "status", ...
	Op     string // "" (=), "LIKE", ">=" or "<="
} // ParamCond applies a verbatim SQL condition when a query param carries one
// of the mapped values (ports the NestJS type_id_log switch: 1 -> al.email_alarm = 1,
// 2 -> al.line_alarm = 1, 3 -> telegram, 4 -> sms, 5 -> nonc).
type ParamCond struct {
	Key   string
	Conds map[string]string // param value -> SQL fragment
}

// ListSpec describes how to paginate one resource.
type ListSpec struct {
	Table        string
	Selects      []string
	Joins        []string
	Filters      []Filter
	FixedFilters []FixedFilter // applied unconditionally
	ParamConds   []ParamCond   // conditional SQL fragments keyed by a param value
	SortCols     map[string]string
	DefaultSort  string
}

// ListQuery carries parsed pagination/sort/filter params.
type ListQuery struct {
	Page     int
	PageSize int
	Sort     string
	values   map[string]string
}

func NewListQuery(page, pageSize int, sort string, values map[string]string) ListQuery {
	return ListQuery{Page: page, PageSize: pageSize, Sort: sort, values: values}
}

func (q ListQuery) Value(key string) string {
	if q.values == nil {
		return ""
	}
	return strings.TrimSpace(q.values[key])
}

// filterExcludedKeys are paging/control params that never appear
// in the echoed "filter" object of a paginated response.
var filterExcludedKeys = map[string]bool{
	"page": true, "pageSize": true, "sort": true,
	"isCount": true, "deletecache": true,
}

// FilterMap returns the request filters echoed back to the client,
// matching the NestJS payload.filter behaviour.
func (q ListQuery) FilterMap() map[string]interface{} {
	out := map[string]interface{}{}
	for k, v := range q.values {
		if filterExcludedKeys[k] {
			continue
		}
		out[k] = v
	}
	return out
}

func EscapeLike(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return r.Replace(s)
}

// ParseSort converts "createddate-DESC" -> "s.createddate DESC" using the whitelist.
func ParseSort(sort string, allowed map[string]string, fallback string) (string, bool) {
	if strings.TrimSpace(sort) == "" {
		return fallback, true
	}
	parts := strings.SplitN(sort, "-", 2)
	if len(parts) != 2 {
		return "", false
	}
	col, ok := allowed[parts[0]]
	if !ok {
		return "", false
	}
	dir := strings.ToUpper(parts[1])
	if dir != "ASC" && dir != "DESC" {
		return "", false
	}
	return col + " " + dir, true
}
