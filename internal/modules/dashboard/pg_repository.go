package dashboard

import (
	"context"

	"icmongolang/internal/modules/dashboard/presenter"
)

// DashboardPgRepository defines data access methods for dashboard aggregate queries.
// ดึงข้อมูลสถิติรวมสำหรับแดชบอร์ด
type DashboardPgRepository interface {
	GetDashboardStats(ctx context.Context) (*presenter.DashboardResponse, error)
	GetRevenueChart(ctx context.Context, period string) ([]*presenter.RevenueData, error)
	GetTopParts(ctx context.Context, limit int) ([]*presenter.TopPartData, error)
	GetJobStatusSummary(ctx context.Context) ([]*presenter.JobStatusSummary, error)
}
