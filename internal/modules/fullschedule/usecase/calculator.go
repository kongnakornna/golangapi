package usecase

import (
	"errors"
	"fmt"
	"time"

	"icmongolang/internal/modules/fullschedule"

	"github.com/robfig/cron/v3"
)

const (
	// maxLookaheadDays bounds the search window for weekly/calendar modes.
	maxLookaheadDays = 366
)

// ErrNoNextRun is returned when no future occurrence exists in the window.
var ErrNoNextRun = errors.New("fullschedule: no next run found within lookahead window")

// parseTimeStart parses "HH:MM" into hour/minute. Returns an error for invalid
// values or times outside 00:00-23:59.
func parseTimeStart(s string) (hour, minute int, err error) {
	if s == "" {
		return 0, 0, nil
	}
	t, cerr := time.Parse("15:04", s)
	if cerr != nil {
		return 0, 0, fmt.Errorf("invalid time_start %q: %w", s, cerr)
	}
	return t.Hour(), t.Minute(), nil
}

// calcNextRun computes the next occurrence of a schedule after now.
// คำนวณเวลาทำงานครั้งถัดไปตามโหมด (normal / full / batch)
func calcNextRun(s *fullschedule.Schedule, now time.Time) (*time.Time, error) {
	switch s.Mode {
	case fullschedule.ModeBatch:
		return nextBatchRun(s, now)
	case fullschedule.ModeFull:
		return nextCalendarRun(s, now)
	default:
		return nextWeeklyRun(s, now)
	}
}

// nextWeeklyRun searches the next day whose weekday flag (sunday..saturday) is 1.
// หาวันถัดไปที่ตรงกับวันในสัปดาห์
func nextWeeklyRun(s *fullschedule.Schedule, now time.Time) (*time.Time, error) {
	hour, minute, err := parseTimeStart(s.TimeStart)
	if err != nil {
		return nil, err
	}
	loc := now.Location()

	for i := 1; i <= maxLookaheadDays; i++ {
		candidate := now.AddDate(0, 0, i)
		wd := candidate.Weekday()
		// sunday..saturday: Sunday=1 ... Saturday=7 (flag 1=on)
		if !weekdayFlag(s, wd) {
			continue
		}
		t := time.Date(candidate.Year(), candidate.Month(), candidate.Day(), hour, minute, 0, 0, loc)
		if t.Before(now) {
			continue
		}
		return &t, nil
	}
	return nil, ErrNoNextRun
}

// weekdayFlag reports whether the given weekday is enabled (sd_iot weekday int8 = 1).
// sunday..saturday: time.Weekday Sunday=0 .. Saturday=6, map to Sunday=1 .. Saturday=7.
func weekdayFlag(s *fullschedule.Schedule, wd time.Weekday) bool {
	switch wd {
	case time.Sunday:
		return s.Sunday == 1
	case time.Monday:
		return s.Monday == 1
	case time.Tuesday:
		return s.Tuesday == 1
	case time.Wednesday:
		return s.Wednesday == 1
	case time.Thursday:
		return s.Thursday == 1
	case time.Friday:
		return s.Friday == 1
	case time.Saturday:
		return s.Saturday == 1
	}
	return false
}

// nextCalendarRun searches the next day that matches month ∈ months and
// day ∈ dates, skipping impossible dates (e.g. 31 Feb).
// หาวันถัดไปที่ตรงกับเดือน/วันที่ในปฏิทิน
func nextCalendarRun(s *fullschedule.Schedule, now time.Time) (*time.Time, error) {
	hour, minute, err := parseTimeStart(s.TimeStart)
	if err != nil {
		return nil, err
	}
	monthSet := toSet([]int(s.Months))
	dateSet := toSet([]int(s.Dates))
	loc := now.Location()

	for i := 1; i <= maxLookaheadDays; i++ {
		candidate := now.AddDate(0, 0, i)
		if !monthSet[int(candidate.Month())] {
			continue
		}
		if !dateSet[candidate.Day()] {
			continue
		}
		// Guard impossible date (time.Date normalizes; check day unchanged).
		probe := time.Date(candidate.Year(), candidate.Month(), candidate.Day(), 0, 0, 0, 0, loc)
		if probe.Day() != candidate.Day() {
			continue
		}
		t := time.Date(candidate.Year(), candidate.Month(), candidate.Day(), hour, minute, 0, 0, loc)
		if t.Before(now) {
			continue
		}
		return &t, nil
	}
	return nil, ErrNoNextRun
}

// nextBatchRun uses robfig/cron v3 to compute the next run.
// คำนวณครั้งถัดไปจาก Cron Expression
func nextBatchRun(s *fullschedule.Schedule, now time.Time) (*time.Time, error) {
	if s.CronExpr == nil || *s.CronExpr == "" {
		return nil, errors.New("fullschedule: cron_expr is required for batch mode")
	}
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	cs, err := parser.Parse(*s.CronExpr)
	if err != nil {
		return nil, fmt.Errorf("fullschedule: invalid cron_expr: %w", err)
	}
	next := cs.Next(now)
	return &next, nil
}

// toSet converts a []int into a map[int]bool helper.
func toSet(vals []int) map[int]bool {
	out := make(map[int]bool, len(vals))
	for _, v := range vals {
		out[v] = true
	}
	return out
}
