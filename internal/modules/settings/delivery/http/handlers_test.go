package http

import (
	"net/http"

	_ "icmongolang/internal/modules/settings/presenter"
	_ "icmongolang/pkg/responses"
	"net/http/httptest"
	"strings"
	"testing"

	"icmongolang/config"
	"icmongolang/internal/modules/settings"
	"icmongolang/pkg/logger"

	"github.com/go-chi/chi/v5"
)

type nopLogger struct{}

func (nopLogger) InitLogger()                        {}
func (nopLogger) Debug(args ...interface{})          {}
func (nopLogger) Debugf(t string, a ...interface{})  {}
func (nopLogger) Info(args ...interface{})           {}
func (nopLogger) Infof(t string, a ...interface{})   {}
func (nopLogger) Warn(args ...interface{})           {}
func (nopLogger) Warnf(t string, a ...interface{})   {}
func (nopLogger) Error(args ...interface{})          {}
func (nopLogger) Errorf(t string, a ...interface{})  {}
func (nopLogger) DPanic(args ...interface{})         {}
func (nopLogger) DPanicf(t string, a ...interface{}) {}
func (nopLogger) Fatal(args ...interface{})          {}
func (nopLogger) Fatalf(t string, a ...interface{})  {}
func (nopLogger) Sync() error                        { return nil }

var _ = logger.Logger(nopLogger{})

func newTestHandler() settings.Handlers {
	return CreateSettingsHandler(nil, &config.Config{}, nopLogger{})
}

func TestDecodeBodyRejectsGarbage(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/x", strings.NewReader("{not json"))
	if _, err := decodeBody(req); err == nil {
		t.Fatal("expected decode error for malformed JSON")
	}
}

func TestDecodeBodyAcceptsObject(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/x", strings.NewReader(`{"setting_name":"s1"}`))
	body, err := decodeBody(req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if body["setting_name"] != "s1" {
		t.Fatalf("body not decoded: %+v", body)
	}
}

func TestFilterKeysWhitelists(t *testing.T) {
	body := map[string]interface{}{
		"setting_name": "s1",
		"sn":           "SN1",
		"evil":         "drop me",
		"nil_field":    nil,
	}
	got := filterKeys(body, settingCols)
	if _, exists := got["evil"]; exists {
		t.Fatal("non-whitelisted key must be dropped")
	}
	if _, exists := got["nil_field"]; exists {
		t.Fatal("nil values must be dropped")
	}
	if got["setting_name"] != "s1" || got["sn"] != "SN1" {
		t.Fatalf("whitelisted values missing: %+v", got)
	}
}

func TestCoerceID(t *testing.T) {
	if v := coerceID("42"); v != 42 {
		t.Fatalf("numeric id should coerce to int, got %T", v)
	}
	if v := coerceID("abc-uuid"); v != "abc-uuid" {
		t.Fatalf("non-numeric id should stay string, got %T", v)
	}
}

func TestToStr(t *testing.T) {
	cases := []struct {
		in   interface{}
		want string
	}{
		{"s", "s"},
		{float64(7), "7"},
		{float64(7.5), "7.5"},
		{nil, ""},
	}
	for _, c := range cases {
		if got := toStr(c.in); got != c.want {
			t.Errorf("toStr(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestRoutesRegisterWithoutPanic verifies MapSettingRoute registers the full
// inventory and routes requests without panicking.
func TestRoutesRegisterWithoutPanic(t *testing.T) {
	h := newTestHandler()
	r := chi.NewRouter()
	MapSettingRoute(r, h, nil)

	paths := 0
	err := chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		paths++
		return nil
	})
	if err != nil {
		t.Fatalf("walk failed: %v", err)
	}
	if paths < 150 {
		t.Fatalf("expected >=150 registered routes, got %d", paths)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/settings/listsetting", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code == 0 {
		t.Fatal("no response written")
	}
}
