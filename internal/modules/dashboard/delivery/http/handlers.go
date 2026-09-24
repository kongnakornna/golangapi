package http

import (
	"net/http"

	"icmongolang/config"
	"icmongolang/internal/modules/dashboard"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"

	"github.com/go-chi/render"
)

type dashboardHandler struct {
	cfg    *config.Config
	dashUC dashboard.DashboardUseCaseI
	logger logger.Logger
}

// CreateDashboardHandler creates a new dashboard HTTP handler.
// สร้างตัวจัดการ HTTP สำหรับแดชบอร์ด
func CreateDashboardHandler(uc dashboard.DashboardUseCaseI, cfg *config.Config, logger logger.Logger) dashboard.Handlers {
	return &dashboardHandler{cfg: cfg, dashUC: uc, logger: logger}
}

// GetDashboardStats godoc
// @Summary Get dashboard statistics
// @Description Get overall dashboard statistics including device counts, alerts, and commands.
// @Tags dashboard
// @Accept json
// @Produce json
// @Success 200 {object} responses.SuccessResponse[presenter.DashboardResponse]
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/dashboard/stats [get]
func (h *dashboardHandler) GetDashboardStats() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := h.dashUC.GetDashboardStats(r.Context())
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(data))
	}
}

// GetRevenueChart godoc
// @Summary Get revenue chart data
// @Description Get revenue data grouped by period (daily, weekly, monthly).
// @Tags dashboard
// @Accept json
// @Produce json
// @Param period query string false "Period: daily, weekly, monthly" default(monthly)
// @Success 200 {object} responses.SuccessResponse[[]presenter.RevenueData]
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/dashboard/revenue [get]
func (h *dashboardHandler) GetRevenueChart() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		period := r.URL.Query().Get("period")
		if period == "" {
			period = "monthly"
		}
		data, err := h.dashUC.GetRevenueChart(r.Context(), period)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(data))
	}
}

// GetTopParts godoc
// @Summary Get top parts
// @Description Get top most used/sold parts.
// @Tags dashboard
// @Accept json
// @Produce json
// @Param limit query int false "Number of top parts" default(10)
// @Success 200 {object} responses.SuccessResponse[[]presenter.TopPartData]
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/dashboard/top-parts [get]
func (h *dashboardHandler) GetTopParts() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 10
		data, err := h.dashUC.GetTopParts(r.Context(), limit)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(data))
	}
}

// GetJobStatusSummary godoc
// @Summary Get job status summary
// @Description Get summary of job statuses (pending, running, completed, failed).
// @Tags dashboard
// @Accept json
// @Produce json
// @Success 200 {object} responses.SuccessResponse[[]presenter.JobStatusSummary]
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/dashboard/job-status [get]
func (h *dashboardHandler) GetJobStatusSummary() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := h.dashUC.GetJobStatusSummary(r.Context())
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(data))
	}
}
