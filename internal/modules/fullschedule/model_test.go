package fullschedule_test

import (
	"testing"

	"icmongolang/internal/modules/fullschedule"
)

func TestScheduleStatusScan(t *testing.T) {
	tests := []struct {
		name string
		src  interface{}
		want fullschedule.ScheduleStatus
	}{
		{name: "int active", src: int64(1), want: fullschedule.StatusActive},
		{name: "int inactive", src: int64(0), want: fullschedule.StatusInactive},
		{name: "byte 1", src: []byte("1"), want: fullschedule.StatusActive},
		{name: "byte active", src: []byte("active"), want: fullschedule.StatusActive},
		{name: "byte 0", src: []byte("0"), want: fullschedule.StatusInactive},
		{name: "string 1", src: "1", want: fullschedule.StatusActive},
		{name: "string active", src: "active", want: fullschedule.StatusActive},
		{name: "string inactive", src: "inactive", want: fullschedule.StatusInactive},
		{name: "nil", src: nil, want: fullschedule.StatusInactive},
		{name: "unknown type", src: 3.14, want: fullschedule.StatusInactive},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var s fullschedule.ScheduleStatus
			if err := s.Scan(tc.src); err != nil {
				t.Fatalf("Scan() error = %v", err)
			}
			if s != tc.want {
				t.Fatalf("Scan() = %q, want %q", s, tc.want)
			}
		})
	}
}

func TestScheduleStatusValue(t *testing.T) {
	if v, err := fullschedule.StatusActive.Value(); err != nil || v != int64(1) {
		t.Fatalf("StatusActive.Value() = %v, %v; want int64(1)", v, err)
	}
	if v, err := fullschedule.StatusInactive.Value(); err != nil || v != int64(0) {
		t.Fatalf("StatusInactive.Value() = %v, %v; want int64(0)", v, err)
	}
}

func TestScheduleStatusValueScan(t *testing.T) {
	tests := []struct {
		name string
		src  interface{}
		want fullschedule.ScheduleStatusValue
	}{
		{name: "int64 draft", src: int64(3), want: fullschedule.StatusValueDraft},
		{name: "int64 active", src: int64(1), want: fullschedule.StatusValueActive},
		{name: "byte draft", src: []byte("3"), want: fullschedule.StatusValueDraft},
		{name: "byte active", src: []byte("active"), want: fullschedule.StatusValueActive},
		{name: "byte draft word", src: []byte("draft"), want: fullschedule.StatusValueDraft},
		{name: "string inactive", src: "inactive", want: fullschedule.StatusValueInactive},
		{name: "nil", src: nil, want: fullschedule.StatusValueInactive},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var s fullschedule.ScheduleStatusValue
			if err := s.Scan(tc.src); err != nil {
				t.Fatalf("Scan() error = %v", err)
			}
			if s != tc.want {
				t.Fatalf("Scan() = %d, want %d", s, tc.want)
			}
		})
	}
}

func TestScheduleStatusValueValue(t *testing.T) {
	for st, want := range map[fullschedule.ScheduleStatusValue]int64{
		fullschedule.StatusValueInactive: 0,
		fullschedule.StatusValueActive:   1,
		fullschedule.StatusValueDraft:    3,
	} {
		v, err := st.Value()
		if err != nil || v != want {
			t.Fatalf("ScheduleStatusValue(%d).Value() = %v, %v; want int64(%d)", st, v, err, want)
		}
	}
}

func TestIntArrayScan(t *testing.T) {
	tests := []struct {
		name string
		src  interface{}
		want fullschedule.IntArray
	}{
		{name: "nil", src: nil, want: nil},
		{name: "string ok", src: "{1,15}", want: fullschedule.IntArray{1, 15}},
		{name: "string empty", src: "{}", want: fullschedule.IntArray{}},
		{name: "string spaces", src: " {2, 3, 4} ", want: fullschedule.IntArray{2, 3, 4}},
		{name: "bytes ok", src: []byte("{7,8}"), want: fullschedule.IntArray{7, 8}},
		{name: "bytes empty", src: []byte("{}"), want: fullschedule.IntArray{}},
		{name: "int64 slice", src: []int64{5, 6, 7}, want: fullschedule.IntArray{5, 6, 7}},
		{name: "int slice", src: []int{9, 10}, want: fullschedule.IntArray{9, 10}},
		{name: "interface slice", src: []interface{}{int64(1), int64(2)}, want: fullschedule.IntArray{1, 2}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var a fullschedule.IntArray
			if err := a.Scan(tc.src); err != nil {
				t.Fatalf("Scan() error = %v", err)
			}
			if len(a) != len(tc.want) {
				t.Fatalf("Scan() = %v, want %v", a, tc.want)
			}
			for i := range a {
				if a[i] != tc.want[i] {
					t.Fatalf("Scan() = %v, want %v", a, tc.want)
				}
			}
		})
	}
}

func TestIntArrayScanError(t *testing.T) {
	bad := []interface{}{
		"not-an-array",
		"{1,x}",
		`{1,"unterminated`,
		42,
		"{NULL}",
	}
	for _, src := range bad {
		var a fullschedule.IntArray
		if err := a.Scan(src); err == nil {
			t.Fatalf("Scan(%v) expected error, got nil (%v)", src, a)
		}
	}
}

func TestIntArrayValue(t *testing.T) {
	tests := []struct {
		name string
		in   fullschedule.IntArray
		want string
	}{
		{name: "nil", in: nil, want: "{}"},
		{name: "single", in: fullschedule.IntArray{1}, want: "{1}"},
		{name: "multi", in: fullschedule.IntArray{1, 15, 31}, want: "{1,15,31}"},
		{name: "empty", in: fullschedule.IntArray{}, want: "{}"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v, err := tc.in.Value()
			if err != nil {
				t.Fatalf("Value() error = %v", err)
			}
			if v != tc.want {
				t.Fatalf("Value() = %v, want %q", v, tc.want)
			}
		})
	}
}
