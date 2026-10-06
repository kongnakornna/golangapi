package http

import (
	"net/url"
	"testing"
)

func TestParsePageParams(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		wantPage    int
		wantPerPage int
	}{
		{"defaults (empty)", "", 1, 10},
		{"explicit page & per_page", "page=2&per_page=20", 2, 20},
		{"page_size fallback", "page=3&page_size=15", 3, 15},
		{"limit fallback", "page=1&limit=5", 1, 5},
		{"per_page clamps to 100", "per_page=999", 1, 100},
		{"page < 1 becomes 1", "page=0", 1, 10},
		{"negative page becomes 1", "page=-5", 1, 10},
		{"per_page preferred over page_size", "per_page=25&page_size=15", 1, 25},
		{"per_page preferred over limit", "per_page=25&limit=5", 1, 25},
		{"page_size preferred over limit", "page_size=30&limit=5", 1, 30},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, _ := url.ParseQuery(tt.query)
			page, perPage := parsePageParams(q)
			if page != tt.wantPage || perPage != tt.wantPerPage {
				t.Fatalf("got page=%d per_page=%d, want page=%d per_page=%d",
					page, perPage, tt.wantPage, tt.wantPerPage)
			}
		})
	}
}

func TestParseSortParam(t *testing.T) {
	tests := []struct {
		in        string
		wantField string
		wantOrder string
	}{
		{"", "", ""},
		{"email-ASC", "email", "ASC"},
		{"createddate-DESC", "createddate", "DESC"},
		{"username", "username", ""},
		{"username-asc", "username", "asc"},
		{"with-dash-in-name-ASC", "with-dash-in-name", "ASC"},
	}
	for _, tt := range tests {
		field, order := parseSortParam(tt.in)
		if field != tt.wantField || order != tt.wantOrder {
			t.Fatalf("parseSortParam(%q) = (%q, %q), want (%q, %q)",
				tt.in, field, order, tt.wantField, tt.wantOrder)
		}
	}
}