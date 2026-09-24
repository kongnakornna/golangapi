package repository

import (
	"context"
	"os"
	"testing"

	"icmongolang/internal/modules/fullschedule"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestListPage_ScanMonthsDates is an integration test that verifies the reported
// bug (scanning PostgreSQL integer[] months/dates into []int) is fixed against a
// real DB. It is skipped unless FULLSCHEDULE_TEST_DB_DSN is set, so it never
// runs in a normal unit-test/CI environment.
func TestListPage_ScanMonthsDates(t *testing.T) {
	dsn := os.Getenv("FULLSCHEDULE_TEST_DB_DSN")
	if dsn == "" {
		t.Skip("FULLSCHEDULE_TEST_DB_DSN not set; skipping live DB integration test")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	repo := CreateFullSchedulePgRepository(db)
	ctx := context.Background()

	rows, total, err := repo.ListPage(ctx, fullschedule.ScheduleFilter{}, 1, 10)
	if err != nil {
		t.Fatalf("ListPage returned scan error: %v", err)
	}
	if total == 0 {
		t.Fatal("expected seeded fs_schedule rows, got 0")
	}
	for _, r := range rows {
		// If scanning failed, Months/Dates would be nil; ensure non-nil slices present.
		if r.Months == nil || r.Dates == nil {
			t.Fatalf("schedule %s months/dates not scanned (nil)", r.ID)
		}
	}
}
