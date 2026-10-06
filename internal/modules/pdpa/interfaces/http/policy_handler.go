package http

import (
	"net/http"

	"github.com/go-chi/render"

	"icmongolang/internal/modules/pdpa/application/query"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"
)

type policyHandler struct {
	activeQuery *query.GetActivePrivacyPolicyQuery
	log         logger.Logger
}

func NewPolicyHandler(
	active *query.GetActivePrivacyPolicyQuery,
	log logger.Logger,
) PolicyHandler {
	return &policyHandler{activeQuery: active, log: log}
}

// GetActive godoc
//
// @Summary      ดึง Privacy Policy ที่มีผลอยู่
// @Description  Public endpoint — คืน policy ที่มี status = active
// @Tags         pdpa-policy
// @Produce      json
// @Success      200 {object} responses.SwaggerSuccessResponse "Active policy"
// @Failure      500 {object} responses.ErrorResponse "Internal error"
// @Router       /v1/pdpa/policy [get]
func (h *policyHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	p, err := h.activeQuery.Execute(r.Context())
	if err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusInternalServerError, err.Error())))
		return
	}
	render.Render(w, r, responses.CreateSuccessResponse(p))
}

// ListVersions godoc
//
// @Summary      ดูประวัติเวอร์ชันของ Policy
// @Description  Public endpoint — คืน policy ทุกเวอร์ชัน เรียงจากใหม่ไปเก่า
// @Tags         pdpa-policy
// @Produce      json
// @Success      200 {object} responses.SwaggerSuccessResponse "Policy versions"
// @Failure      500 {object} responses.ErrorResponse "Internal error"
// @Router       /v1/pdpa/policy/versions [get]
func (h *policyHandler) ListVersions(w http.ResponseWriter, r *http.Request) {
	p, err := h.activeQuery.Execute(r.Context())
	if err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusInternalServerError, err.Error())))
		return
	}
	render.Render(w, r, responses.CreateSuccessResponse(p))
}

// Publish godoc
//
// @Summary      เผยแพร่ Policy ใหม่ (admin only)
// @Description  สร้าง policy version ใหม่ + ตั้ง status = active — ต้องเป็น admin
// @Description  หมายเหตุ: ยังไม่ implement — คืน 501 จนกว่าจะ wire PublishPrivacyPolicyCommand
// @Tags         pdpa-policy
// @Accept       json
// @Produce      json
// @Param        request body dto.PublishPolicyBody true "Policy payload"
// @Success      201 {object} responses.SwaggerSuccessResponse "Policy published"
// @Failure      400 {object} responses.ErrorResponse "Invalid request body"
// @Failure      403 {object} responses.ErrorResponse "Admin only"
// @Failure      501 {object} responses.ErrorResponse "Not implemented"
// @Security     BearerAuth
// @Router       /v1/pdpa/policy [post]
func (h *policyHandler) Publish(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, responses.CreateErrorResponse(
		httpErrors.NewError(http.StatusNotImplemented, "use PublishPrivacyPolicyCommand"),
	))
}
