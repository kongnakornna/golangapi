package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/application/command"
	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"
)

type accountHandler struct {
	suspendCmd   *command.SuspendAccountCommand
	terminateCmd *command.TerminateAccountCommand
	confirmCmd   *command.ConfirmDeletionCommand
	log          logger.Logger
}

func NewAccountHandler(
	suspend *command.SuspendAccountCommand,
	terminate *command.TerminateAccountCommand,
	confirm *command.ConfirmDeletionCommand,
	log logger.Logger,
) AccountHandler {
	return &accountHandler{suspendCmd: suspend, terminateCmd: terminate, confirmCmd: confirm, log: log}
}

// Suspend godoc
//
// @Summary      ระงับบัญชีผู้ใช้ (admin only)
// @Description  ระงับบัญชี + เริ่มนับ retention period — publish event account.suspended
// @Tags         pdpa-account
// @Accept       json
// @Produce      json
// @Param        request body dto.SuspendAccountBody true "Suspend payload"
// @Success      200 {object} responses.SwaggerSuccessResponse "Account suspended"
// @Failure      400 {object} responses.ErrorResponse "Invalid user_id"
// @Failure      403 {object} responses.ErrorResponse "Admin only"
// @Failure      500 {object} responses.ErrorResponse "Internal error"
// @Security     BearerAuth
// @Router       /v1/pdpa/account/suspend [post]
func (h *accountHandler) Suspend(w http.ResponseWriter, r *http.Request) {
	var body dto.SuspendAccountBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
		return
	}

	uid, err := uuid.Parse(body.UserID)
	if err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusBadRequest, "invalid user_id")))
		return
	}

	if err := h.suspendCmd.Execute(r.Context(), command.SuspendAccountInput{
		UserID:         uid,
		RetentionYears: body.RetentionYears,
		Reason:         body.Reason,
	}); err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusInternalServerError, err.Error())))
		return
	}

	render.Render(w, r, responses.CreateSuccessResponse(map[string]string{"status": "suspended"}))
}

// Terminate godoc
//
// @Summary      ยุติบัญชีผู้ใช้ (admin only)
// @Description  ยุติบัญชี + เข้าสู่กระบวนการลบแบบรอเวลาตาม retention — publish event account.terminated
// @Tags         pdpa-account
// @Accept       json
// @Produce      json
// @Param        request body dto.TerminateAccountBody true "Terminate payload"
// @Success      200 {object} responses.SwaggerSuccessResponse "Account terminated"
// @Failure      400 {object} responses.ErrorResponse "Invalid user_id"
// @Failure      403 {object} responses.ErrorResponse "Admin only"
// @Failure      500 {object} responses.ErrorResponse "Internal error"
// @Security     BearerAuth
// @Router       /v1/pdpa/account/terminate [post]
func (h *accountHandler) Terminate(w http.ResponseWriter, r *http.Request) {
	var body dto.TerminateAccountBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
		return
	}

	uid, err := uuid.Parse(body.UserID)
	if err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusBadRequest, "invalid user_id")))
		return
	}

	if err := h.terminateCmd.Execute(r.Context(), uid); err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusInternalServerError, err.Error())))
		return
	}

	render.Render(w, r, responses.CreateSuccessResponse(map[string]string{"status": "terminated"}))
}

// ConfirmDeletion godoc
//
// @Summary      ยืนยันการลบข้อมูลส่วนบุคคล (admin only)
// @Description  ยืนยันว่าผู้ใช้ประสงค์จะลบข้อมูลส่วนบุคคลจริง — publish event data.deletion_requested
// @Tags         pdpa-account
// @Accept       json
// @Produce      json
// @Param        request body dto.ConfirmDeletionBody true "Confirm payload"
// @Success      200 {object} responses.SwaggerSuccessResponse "Deletion confirmed"
// @Failure      400 {object} responses.ErrorResponse "Invalid user_id"
// @Failure      403 {object} responses.ErrorResponse "Admin only"
// @Failure      500 {object} responses.ErrorResponse "Internal error"
// @Security     BearerAuth
// @Router       /v1/pdpa/account/confirm-deletion [post]
func (h *accountHandler) ConfirmDeletion(w http.ResponseWriter, r *http.Request) {
	var body dto.ConfirmDeletionBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
		return
	}

	uid, err := uuid.Parse(body.UserID)
	if err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusBadRequest, "invalid user_id")))
		return
	}

	if err := h.confirmCmd.Execute(r.Context(), uid); err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusInternalServerError, err.Error())))
		return
	}

	render.Render(w, r, responses.CreateSuccessResponse(map[string]string{"status": "confirmed"}))
}
