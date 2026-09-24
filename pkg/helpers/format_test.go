package helpers

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewLocalTimeNormalizesToBangkok(t *testing.T) {
	utc := time.Date(2026, 8, 30, 10, 0, 0, 0, time.UTC)

	lt := NewLocalTime(utc)
	if lt.Time.Format("-07:00") != "+07:00" {
		t.Fatalf("NewLocalTime zone offset = %s, want +07:00", lt.Time.Format("-07:00"))
	}
	if lt.String() != "2026-08-30 17:00:00" {
		t.Fatalf("String() = %q, want 2026-08-30 17:00:00 (UTC+7)", lt.String())
	}

	// Already Thai time → unchanged wall-clock.
	bangkok := time.Date(2026, 8, 30, 17, 0, 0, 0, GetTimeLocation())
	if got := NewLocalTime(bangkok).String(); got != "2026-08-30 17:00:00" {
		t.Fatalf("NewLocalTime(Bangkok) = %q, want 2026-08-30 17:00:00", got)
	}
}

func TestLocalTimeMarshalJSON(t *testing.T) {
	tm := time.Date(2026, 8, 30, 10, 5, 7, 0, time.UTC)

	b, err := json.Marshal(NewLocalTime(tm))
	if err != nil {
		t.Fatalf("Marshal error = %v", err)
	}
	want := `"2026-08-30 17:05:07"` // UTC+7 (Asia/Bangkok) + "2006-01-02 15:04:05"
	if string(b) != want {
		t.Fatalf("MarshalJSON() = %s, want %s", b, want)
	}
}

func TestLocalTimeMarshalJSONZero(t *testing.T) {
	b, err := json.Marshal(NewLocalTime(time.Time{}))
	if err != nil {
		t.Fatalf("Marshal error = %v", err)
	}
	if string(b) != `""` {
		t.Fatalf("MarshalJSON(zero) = %s, want \"\"", b)
	}
}

func TestLocalTimeUnmarshalJSON(t *testing.T) {
	var got LocalTime
	if err := json.Unmarshal([]byte(`"2026-08-30 10:05:07"`), &got); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}
	if got.String() != "2026-08-30 10:05:07" {
		t.Fatalf("got %q, want 2026-08-30 10:05:07", got.String())
	}

	var empty LocalTime
	if err := json.Unmarshal([]byte(`""`), &empty); err != nil {
		t.Fatalf("Unmarshal(empty) error = %v", err)
	}
	if !empty.IsZero() {
		t.Fatalf("empty should stay zero, got %v", empty.Time)
	}
}

func TestLocalTimeString(t *testing.T) {
	tm := time.Date(2026, 8, 30, 3, 0, 0, 0, time.UTC)
	if got := NewLocalTime(tm).String(); got != "2026-08-30 10:00:00" {
		t.Fatalf("String() = %q, want 2026-08-30 10:00:00", got)
	}
	if got := NewLocalTime(time.Time{}).String(); got != "" {
		t.Fatalf("String(zero) = %q, want empty", got)
	}
}
