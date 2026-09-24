package dashboard

import (
	"context"

	"icmongolang/internal/modules/dashboard/presenter"
)

// DashboardUseCaseI defines business logic methods for dashboard.
// อินเทอร์เฟซธุรกิจสำหรับแดชบอร์ด
type DashboardUseCaseI interface {
	GetDashboardStats(ctx context.Context) (*presenter.DashboardResponse, error)
	GetRevenueChart(ctx context.Context, period string) ([]*presenter.RevenueData, error)
	GetTopParts(ctx context.Context, limit int) ([]*presenter.TopPartData, error)
	GetJobStatusSummary(ctx context.Context) ([]*presenter.JobStatusSummary, error)
}
