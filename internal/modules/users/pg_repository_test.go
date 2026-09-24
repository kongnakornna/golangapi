package users

import "testing"

func TestResolveUserSortColumn(t *testing.T) {
	tests := []struct {
		name  string
		field string
		want  string
	}{
		{name: "whitelist username", field: "username", want: "username"},
		{name: "whitelist email", field: "email", want: "email"},
		{name: "whitelist createddate", field: "createddate", want: "createddate"},
		{name: "alias created_at", field: "created_at", want: "createddate"},
		{name: "whitelist status", field: "status", want: "status"},
		{name: "whitelist active_status", field: "active_status", want: "active_status"},
		{name: "case-insensitive", field: "USERNAME", want: "username"},
		{name: "trim spaces", field: "  email  ", want: "email"},
		{name: "unknown falls back to createddate", field: "password; DROP TABLE sd_user", want: "createddate"},
		{name: "empty falls back to createddate", field: "", want: "createddate"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ResolveUserSortColumn(tt.field); got != tt.want {
				t.Fatalf("ResolveUserSortColumn(%q) = %q, want %q", tt.field, got, tt.want)
			}
		})
	}
}

func TestUserFilterSortOrderEffective(t *testing.T) {
	tests := []struct {
		name  string
		order string
		want  string
	}{
		{name: "asc lower", order: "asc", want: "ASC"},
		{name: "ASC upper", order: "ASC", want: "ASC"},
		{name: "asc with spaces", order: "  asc ", want: "ASC"},
		{name: "desc upper", order: "DESC", want: "DESC"},
		{name: "empty defaults DESC", order: "", want: "DESC"},
		{name: "garbage defaults DESC", order: "sideways", want: "DESC"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := UserFilter{SortOrder: tt.order}
			if got := f.SortOrderEffective(); got != tt.want {
				t.Fatalf("SortOrderEffective(%q) = %q, want %q", tt.order, got, tt.want)
			}
		})
	}
}
