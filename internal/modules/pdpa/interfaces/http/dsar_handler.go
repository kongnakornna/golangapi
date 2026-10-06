package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/application/command"
	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/application/query"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"
)

type dsarHandler struct {
	submitCmd  *command.SubmitDSARRequestCommand
	processCmd *command.ProcessDSARCommand
	listQuery  *query.ListDSARByUserQuery
	getQuery   *query.GetDSARStatusQuery
	log        logger.Logger
}

func NewDSARHandler(
	submit *command.SubmitDSARRequestCommand,
	process *command.ProcessDSARCommand,
	list *query.ListDSARByUserQuery,
	get *query.GetDSARStatusQuery,
	log logger.Logger,
) DSARHandler {
	return &dsarHandler{submitCmd: submit, processCmd: process, listQuery: list, getQuery: get, log: log}
}

// Submit godoc
//
// @Summary      ยื่นคำขอใช้สิทธิ (DSAR Submit)
// @Description  ผู้ใช้ยื่นคำขอ Data Subject Access Request — จะสร้าง record พร้อม OTP 15 นาที
// @Tags         pdpa-dsar
// @Accept       json
// @Produce      json
// @Param        request body dto.SubmitDSARBody true "DSAR payload"
// @Success      202 {object} responses.SwaggerSuccessResponse "Request accepted"
// @Failure      400 {object} responses.ErrorResponse "Invalid request body"
// @Failure      401 {object} responses.ErrorResponse "Unauthorized"
// @Failure      500 {object} responses.ErrorResponse "Internal error"
// @Security     BearerAuth
// @Router       /v1/pdpa/dsar [post]
func (h *dsarHandler) Submit(w http.ResponseWriter, r *http.Request) {
	uid, ok := currentUserID(r)
	if !ok {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusUnauthorized, "unauthorized")))
		return
	}

	var body dto.SubmitDSARBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
		return
	}

	d, err := h.submitCmd.Execute(r.Context(), command.SubmitDSARInput{
		UserID:      uid,
		RequestType: valueobject.DSARType(body.RequestType),
		IP:          r.RemoteAddr,
		UA:          r.UserAgent(),
	})
	if err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusInternalServerError, err.Error())))
		return
	}

	render.Status(r, http.StatusAccepted)
	render.Render(w, r, responses.CreateSuccessResponse(d))
}

// List godoc
//
// @Summary      รายการ DSAR ของผู้ใช้
// @Description  ดึงรายการ DSAR ทั้งหมดของผู้ใช้ปัจจุบัน เรียงจากใหม่ไปเก่า
// @Tags         pdpa-dsar
// @Produce      json
// @Success      200 {object} responses.SwaggerSuccessResponse "DSAR list"
// @Failure      401 {object} responses.ErrorResponse "Unauthorized"
// @Failure      500 {object} responses.ErrorResponse "Internal error"
// @Security     BearerAuth
// @Router       /v1/pdpa/dsar [get]
func (h *dsarHandler) List(w http.ResponseWriter, r *http.Request) {
	uid, ok := currentUserID(r)
	if !ok {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusUnauthorized, "unauthorized")))
		return
	}

	rows, err := h.listQuery.Execute(r.Context(), uid)
	if err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusInternalServerError, err.Error())))
		return
	}

	render.Render(w, r, responses.CreateSuccessResponse(rows))
}

// Get godoc
//
// @Summary      รายละเอียด DSAR ตาม ID
// @Description  ดึงรายละเอียด DSAR 1 รายการ (เฉพาะของผู้ใช้ปัจจุบัน)
// @Tags         pdpa-dsar
// @Produce      json
// @Param        id path string true "DSAR ID (UUID)" format(uuid)
// @Success      200 {object} responses.SwaggerSuccessResponse "DSAR detail"
// @Failure      400 {object} responses.ErrorResponse "Invalid ID"
// @Failure      404 {object} responses.ErrorResponse "DSAR not found"
// @Security     BearerAuth
// @Router       /v1/pdpa/dsar/{id} [get]
func (h *dsarHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusBadRequest, "invalid id")))
		return
	}

	d, err := h.getQuery.Execute(r.Context(), id)
	if err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusNotFound, err.Error())))
		return
	}

	render.Render(w, r, responses.CreateSuccessResponse(d))
}

// Process godoc
//
// @Summary      ประมวลผล DSAR (admin only)
// @Description  เปลี่ยนสถานะ DSAR: 3=PROCESSING, 4=COMPLETED, 5=REJECTED — ต้องเป็น admin
// @Tags         pdpa-dsar
// @Accept       json
// @Produce      json
// @Param        id path string true "DSAR ID (UUID)" format(uuid)
// @Param        request body dto.ProcessDSARBody true "Process payload"
// @Success      200 {object} responses.SwaggerSuccessResponse "Updated status"
// @Failure      400 {object} responses.ErrorResponse "Invalid ID or body"
// @Failure      403 {object} responses.ErrorResponse "Admin only"
// @Failure      404 {object} responses.ErrorResponse "DSAR not found"
// @Failure      500 {object} responses.ErrorResponse "Internal error"
// @Security     BearerAuth
// @Router       /v1/pdpa/dsar/{id}/process [post]
func (h *dsarHandler) Process(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusBadRequest, "invalid id")))
		return
	}

	// (optional) decode body ถ้าต้องการใช้ ProcessDSARBody
	var body dto.ProcessDSARBody
	_ = json.NewDecoder(r.Body).Decode(&body)

	if err := h.processCmd.Execute(r.Context(), id); err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusInternalServerError, err.Error())))
		return
	}

	render.Render(w, r, responses.CreateSuccessResponse(map[string]string{"status": "processing"}))
}
