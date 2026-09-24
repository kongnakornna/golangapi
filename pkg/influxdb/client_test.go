package influxdb

import "testing"

func TestFluxEscape(t *testing.T) {
	cases := map[string]string{
		"value":                  "value",
		`value"`:                 `value\"`,
		`a"b\c`:                  `a\"b\\c`,
		"line\nbreak":            `line\nbreak`,
		"tab\there":              `tab\there`,
		`payload" ; |> drop()`:   `payload\" ; |> drop()`,
		`100% >= 90%`:            `100% >= 90%`,
	}
	for in, want := range cases {
		if got := fluxEscape(in); got != want {
			t.Errorf("fluxEscape(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestValidRangeExpr(t *testing.T) {
	valid := []string{"-1h", "now()", "-30d", "2023-01-01T00:00:00Z", "2023-01-01T00:00:00.000Z"}
	invalid := []string{"", "-1h; drop()", `"`, "1h>0", "`"}
	for _, s := range valid {
		if !validRangeExpr(s) {
			t.Errorf("validRangeExpr(%q) = false, want true", s)
		}
	}
	for _, s := range invalid {
		if validRangeExpr(s) {
			t.Errorf("validRangeExpr(%q) = true, want false", s)
		}
	}
}
