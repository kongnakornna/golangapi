package usecase

import (
	"context"

	"icmongolang/config"
	"icmongolang/internal/modules/dashboard"
	"icmongolang/internal/modules/dashboard/presenter"
	"icmongolang/pkg/logger"
)

type dashboardUseCase struct {
	pgRepo dashboard.DashboardPgRepository
	cfg    *config.Config
	logger logger.Logger
}

// CreateDashboardUseCaseI creates a new dashboard use case instance.
// สร้างอินสแตนซ์สำหรับธุรกิจแดชบอร์ด
func CreateDashboardUseCaseI(
	pgRepo dashboard.DashboardPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) dashboard.DashboardUseCaseI {
	return &dashboardUseCase{
		pgRepo: pgRepo,
		cfg:    cfg,
		logger: logger,
	}
}

func (u *dashboardUseCase) GetDashboardStats(ctx context.Context) (*presenter.DashboardResponse, error) {
	return u.pgRepo.GetDashboardStats(ctx)
}

func (u *dashboardUseCase) GetRevenueChart(ctx context.Context, period string) ([]*presenter.RevenueData, error) {
	return u.pgRepo.GetRevenueChart(ctx, period)
}

func (u *dashboardUseCase) GetTopParts(ctx context.Context, limit int) ([]*presenter.TopPartData, error) {
	return u.pgRepo.GetTopParts(ctx, limit)
}

func (u *dashboardUseCase) GetJobStatusSummary(ctx context.Context) ([]*presenter.JobStatusSummary, error) {
	return u.pgRepo.GetJobStatusSummary(ctx)
}
