package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"

	"icmongolang/internal/modules/pdpa/application/command"
	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/application/query"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"
)

type consentHandler struct {
	grantCmd  *command.GrantConsentCommand
	revokeCmd *command.RevokeConsentCommand
	listQuery *query.ListConsentLogByUserQuery
	log       logger.Logger
}

func NewConsentHandler(
	grant *command.GrantConsentCommand,
	revoke *command.RevokeConsentCommand,
	list *query.ListConsentLogByUserQuery,
	log logger.Logger,
) ConsentHandler {
	return &consentHandler{grantCmd: grant, revokeCmd: revoke, listQuery: list, log: log}
}

// RecordConsent godoc
//
// @Summary      บันทึกความยินยอม (Consent)
// @Description  บันทึก consent ของผู้ใช้ตาม purpose ที่กำหนด พร้อม publish event consent.granted
// @Tags         pdpa-consent
// @Accept       json
// @Produce      json
// @Param        request body dto.RecordConsentBody true "Consent payload"
// @Success      201 {object} responses.SwaggerSuccessResponse "Consent created"
// @Failure      400 {object} responses.ErrorResponse "Invalid request body"
// @Failure      401 {object} responses.ErrorResponse "Unauthorized"
// @Failure      500 {object} responses.ErrorResponse "Internal error"
// @Security     BearerAuth
// @Router       /v1/pdpa/consent [post]
func (h *consentHandler) RecordConsent(w http.ResponseWriter, r *http.Request) {
	uid, ok := currentUserID(r)
	if !ok {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusUnauthorized, "unauthorized")))
		return
	}

	var body dto.RecordConsentBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
		return
	}

	result, err := h.grantCmd.Execute(r.Context(), command.GrantConsentInput{
		UserID:  uid,
		Purpose: valueobject.ConsentPurpose(body.Purpose),
		Source:  body.Source,
		IP:      r.RemoteAddr,
		UA:      r.UserAgent(),
	})
	if err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusInternalServerError, err.Error())))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, responses.CreateSuccessResponse(result))
}

// RevokeConsent godoc
//
// @Summary      ถอนความยินยอม (Revoke Consent)
// @Description  ถอน consent ของผู้ใช้ตาม purpose ที่กำหนด พร้อม publish event consent.revoked
// @Tags         pdpa-consent
// @Accept       json
// @Produce      json
// @Param        request body dto.RevokeConsentBody true "Revoke payload"
// @Success      200 {object} responses.SwaggerSuccessResponse "Revoked successfully"
// @Failure      400 {object} responses.ErrorResponse "Invalid request body"
// @Failure      401 {object} responses.ErrorResponse "Unauthorized"
// @Failure      404 {object} responses.ErrorResponse "Consent not found"
// @Failure      500 {object} responses.ErrorResponse "Internal error"
// @Security     BearerAuth
// @Router       /v1/pdpa/consent [delete]
func (h *consentHandler) RevokeConsent(w http.ResponseWriter, r *http.Request) {
	uid, ok := currentUserID(r)
	if !ok {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusUnauthorized, "unauthorized")))
		return
	}

	var body dto.RevokeConsentBody
	_ = json.NewDecoder(r.Body).Decode(&body)

	if err := h.revokeCmd.Execute(r.Context(), command.RevokeConsentInput{
		UserID:  uid,
		Purpose: valueobject.ConsentPurpose(body.Purpose),
	}); err != nil {
		render.Render(w, r, responses.CreateErrorResponse(httpErrors.NewError(http.StatusInternalServerError, err.Error())))
		return
	}

	render.Render(w, r, responses.CreateSuccessResponse(map[string]string{"status": "revoked"}))
}

// GetHistory godoc
//
// @Summary      ประวัติความยินยอม
// @Description  ดึงประวัติ consent ทั้งหมดของผู้ใช้ปัจจุบัน เรียงจากใหม่ไปเก่า
// @Tags         pdpa-consent
// @Produce      json
// @Success      200 {object} responses.SwaggerSuccessResponse "Consent history"
// @Failure      401 {object} responses.ErrorResponse "Unauthorized"
// @Failure      500 {object} responses.ErrorResponse "Internal error"
// @Security     BearerAuth
// @Router       /v1/pdpa/consent/history [get]
func (h *consentHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
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
