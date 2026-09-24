package settings

import (
	"testing"
)

func TestParseSort(t *testing.T) {
	allowed := map[string]string{
		"setting_id":   "s.setting_id",
		"setting_name": "s.setting_name",
		"createddate":  "s.createddate",
	}
	fallback := "s.setting_id ASC"

	tests := []struct {
		name   string
		sort   string
		want   string
		wantOK bool
	}{
		{"empty falls back", "", fallback, true},
		{"valid asc", "createddate-ASC", "s.createddate ASC", true},
		{"valid desc", "setting_name-desc", "s.setting_name DESC", true},
		{"missing direction", "createddate", "", false},
		{"unknown field", "hacker-ASC", "", false},
		{"bad direction", "setting_id-SIDEWAYS", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ParseSort(tt.sort, allowed, fallback)
			if ok != tt.wantOK || got != tt.want {
				t.Fatalf("ParseSort(%q) = (%q, %v), want (%q, %v)", tt.sort, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestEscapeLike(t *testing.T) {
	tests := []struct{ in, want string }{
		{`plain`, `plain`},
		{`50%`, `50\%`},
		{`a_b`, `a\_b`},
		{`back\slash`, `back\\slash`},
		{`%_\`, `\%\_\\`},
	}
	for _, tt := range tests {
		if got := EscapeLike(tt.in); got != tt.want {
			t.Errorf("EscapeLike(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestListQueryValue(t *testing.T) {
	q := NewListQuery(2, 50, "createddate-DESC", map[string]string{"status": "1"})
	if q.Value("status") != "1" {
		t.Fatalf("Value(status) = %q", q.Value("status"))
	}
	if q.Value("missing") != "" {
		t.Fatalf("Value(missing) should be empty")
	}
	var nilQ ListQuery
	if nilQ.Value("x") != "" {
		t.Fatal("nil values map should be safe")
	}
}

func TestSpecWithAppendsFixedFilters(t *testing.T) {
	spec := SettingSpec.With(FixedFilter{Column: "s.status", Value: "1"})
	if len(spec.FixedFilters) != 1 {
		t.Fatalf("expected 1 fixed filter, got %d", len(spec.FixedFilters))
	}
	if len(SettingSpec.FixedFilters) != 0 {
		t.Fatal("base spec must not be mutated")
	}
}
