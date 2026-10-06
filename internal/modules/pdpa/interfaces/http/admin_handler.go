package http

import (
	"net/http"
	"time"

	"github.com/go-chi/render"

	"icmongolang/internal/modules/pdpa/application/query"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"
)

type adminHandler struct {
	reportQuery *query.GetAdminReportQuery
	log         logger.Logger
}

func NewAdminHandler(report *query.GetAdminReportQuery, log logger.Logger) AdminHandler {
	return &adminHandler{reportQuery: report, log: log}
}

// GetReport godoc
//
// @Summary      รายงาน PDPA compliance (admin only)
// @Description  สรุปจำนวน audit trail ใน 1 เดือนล่าสุด จัดกลุ่มตาม action
// @Tags         pdpa-admin
// @Produce      json
// @Success      200 {object} responses.SwaggerSuccessResponse "Compliance report"
// @Failure      403 {object} responses.ErrorResponse "Admin only"
// @Failure      500 {object} responses.ErrorResponse "Internal error"
// @Security     BearerAuth
// @Router       /v1/pdpa/admin/reports [get]
func (h *adminHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	to := time.Now()
	from := to.AddDate(0, -1, 0)

	rep, err := h.reportQuery.Execute(r.Context(), from, to)
	if err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusInternalServerError, err.Error())))
		return
	}

	render.Render(w, r, responses.CreateSuccessResponse(rep))
}
