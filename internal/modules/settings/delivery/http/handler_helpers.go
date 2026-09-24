package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"icmongolang/internal/modules/settings"
	"icmongolang/pkg/helpers"

	"github.com/go-chi/render"
)

// ---- generic handler factories shared by all domain files ----

func (h *settingsHandler) listHandler(spec settings.ListSpec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := h.uc.ListPaginate(r.Context(), spec, listQueryFrom(r))
		if err != nil {
			fail(w, r, err)
			return
		}
		ok(w, r, res)
	}
}

func (h *settingsHandler) allHandler(spec settings.ListSpec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := listQueryFrom(r)
		q.PageSize = 0 // no LIMIT/OFFSET
		items, err := h.uc.RowsAll(r.Context(), spec, q)
		if err != nil {
			fail(w, r, err)
			return
		}
		ok(w, r, items)
	}
}

// dupCheck mirrors the NestJS create guards ("The <field> duplicate this
// data cannot create."): when any row already holds the submitted value the
// insert is rejected with a 422.
type dupCheck struct {
	BodyKey string // JSON body key holding the submitted value
	Column  string // table column to compare against
	Field   string // field label used in the error message
}

func (h *settingsHandler) createHandler(table string, cols []string, defaults map[string]interface{}, dupChecks ...dupCheck) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := decodeBody(r)
		if err != nil {
			badRequest(w, r, err.Error())
			return
		}
		for _, dc := range dupChecks {
			val := strings.TrimSpace(toStr(body[dc.BodyKey]))
			if val == "" {
				continue
			}
			spec := settings.ListSpec{
				Table:        table,
				Selects:      []string{dc.Column + " AS dup_hit"},
				FixedFilters: []settings.FixedFilter{{Column: dc.Column, Value: val}},
			}
			items, err := h.uc.RowsAll(r.Context(), spec, settings.NewListQuery(1, 1, "", nil))
			if err != nil {
				fail(w, r, err)
				return
			}
			if len(items) > 0 {
				render422(w, r, "The "+dc.Field+" "+val+" duplicate this data cannot create.")
				return
			}
		}
		row := filterKeys(body, cols)
		if len(row) == 0 {
			badRequest(w, r, "request body has no valid fields")
			return
		}
		for k, v := range defaults {
			if _, exists := row[k]; !exists {
				row[k] = v
			}
		}
		if err = h.uc.CreateRow(r.Context(), table, row); err != nil {
			fail(w, r, err)
			return
		}
		ok(w, r, row)
	}
}

func (h *settingsHandler) updateBodyHandler(table, idCol string, cols []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := decodeBody(r)
		if err != nil {
			badRequest(w, r, err.Error())
			return
		}
		rawID, _ := body[idCol].(string)
		if rawID == "" {
			if v, ok2 := body[idCol]; ok2 && v != nil {
				rawID = strings.TrimSpace(toStr(v))
			}
		}
		if rawID == "" {
			badRequest(w, r, idCol+" is required")
			return
		}
		fields := filterKeys(body, cols)
		fields["updateddate"] = nowPtr()
		n, err := h.uc.UpdateFields(r.Context(), table, idCol, coerceID(rawID), fields)
		if err != nil {
			fail(w, r, err)
			return
		}
		ok(w, r, map[string]interface{}{"affected": n})
	}
}

func toStr(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case int64:
		return strconv.FormatInt(t, 10)
	case int:
		return strconv.Itoa(t)
	case float64:
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10)
		}
		return strconv.FormatFloat(t, 'f', -1, 64)
	default:
		return ""
	}
}

func (h *settingsHandler) deleteGetHandler(table, idCol, param string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.URL.Query().Get(param))
		if id == "" {
			badRequest(w, r, param+" is required")
			return
		}
		n, err := h.uc.DeleteRow(r.Context(), table, idCol+" = ?", coerceID(id))
		if err != nil {
			fail(w, r, err)
			return
		}
		if n == 0 {
			notFound(w, r, "Data not found.")
			return
		}
		affected(w, r, n)
	}
}

func (h *settingsHandler) deleteParamHandler(table, idCol string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := urlID(r)
		if id == "" {
			badRequest(w, r, "id is required")
			return
		}
		n, err := h.uc.DeleteRow(r.Context(), table, idCol+" = ?", coerceID(id))
		if err != nil {
			fail(w, r, err)
			return
		}
		if n == 0 {
			notFound(w, r, "Data not found.")
			return
		}
		affected(w, r, n)
	}
}

// statusBodyHandler ports the NestJS update_*_status endpoints: body must
// carry <idCol> and status; when resetAll is true and status==1 every other
// row of the table is deactivated first (mirrors update_email_status etc.).
// A missing target row yields 404 like the NestJS NotFoundException.
func (h *settingsHandler) statusBodyHandler(table, idCol string, resetAll bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := decodeBody(r)
		if err != nil {
			badRequest(w, r, err.Error())
			return
		}
		rawID := toStr(body[idCol])
		if rawID == "" {
			badRequest(w, r, idCol+" is required")
			return
		}
		status := 0
		switch v := body["status"].(type) {
		case float64:
			status = int(v)
		case string:
			status, _ = strconv.Atoi(strings.TrimSpace(v))
		}
		n, err := h.uc.SetStatusExclusive(r.Context(), table, idCol, coerceID(rawID), status, resetAll)
		if err != nil {
			fail(w, r, err)
			return
		}
		if n == 0 {
			notFound(w, r, "No devices found with bucket '"+rawID+"'")
			return
		}
		affected(w, r, n)
	}
}

// getOneByParam loads the first row of spec filtered by a fixed query param.
func (h *settingsHandler) getOneByParam(spec settings.ListSpec, param, fixedCol string, required bool, missingMsg string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.URL.Query().Get(param))
		if id == "" {
			id = urlID(r)
		}
		if id == "" && required {
			render.Render(w, r, response422(missingMsg))
			return
		}
		s := spec
		if id != "" {
			s = s.With(settings.FixedFilter{Column: fixedCol, Value: id})
		}
		items, err := h.uc.RowsAll(r.Context(), s, settings.NewListQuery(1, 1, "", nil))
		if err != nil {
			fail(w, r, err)
			return
		}
		if len(items) == 0 {
			notFound(w, r, missingMsg)
			return
		}
		ok(w, r, items[0])
	}
}

func (h *settingsHandler) insertLinkViaGet(table string, params ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		row := map[string]interface{}{}
		for _, p := range params {
			v := strings.TrimSpace(r.URL.Query().Get(p))
			if v == "" {
				badRequest(w, r, p+" is required")
				return
			}
			row[p] = coerceID(v)
		}
		if err := h.uc.CreateRow(r.Context(), table, row); err != nil {
			fail(w, r, err)
			return
		}
		ok(w, r, row)
	}
}

// convertDatesToBangkok converts createddate and updateddate fields
// from their raw database format to "2006-01-02 15:04:05" in Bangkok timezone.
func convertDatesToBangkok(items []map[string]interface{}) {
	loc := helpers.GetTimeLocation()
	for _, item := range items {
		for _, key := range []string{"createddate", "updateddate"} {
			if v, ok := item[key]; ok && v != nil {
				item[key] = toBangkokString(v, loc)
			}
		}
	}
}

func toBangkokString(v interface{}, loc *time.Location) string {
	var t time.Time
	switch val := v.(type) {
	case time.Time:
		t = val
	case string:
		for _, layout := range []string{
			"2006-01-02T15:04:05.000000Z",
			time.RFC3339,
			"2006-01-02 15:04:05",
			"2006-01-02",
		} {
			if parsed, err := time.Parse(layout, val); err == nil {
				t = parsed
				break
			}
		}
	}
	if t.IsZero() {
		return fmt.Sprintf("%v", v)
	}
	return helpers.TimeConvertermas(t.In(loc))
}
