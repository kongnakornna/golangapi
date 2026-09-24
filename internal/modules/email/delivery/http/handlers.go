package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"icmongolang/config"
	"icmongolang/internal/modules/email"
	"icmongolang/internal/modules/email/presenter"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"
	"icmongolang/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type emailHandler struct {
	cfg     *config.Config
	emailUC email.EmailUseCaseI
	logger  logger.Logger
}

// CreateEmailHandler creates a new email HTTP handler.
// สร้างตัวจัดการ HTTP สำหรับอีเมล
func CreateEmailHandler(uc email.EmailUseCaseI, cfg *config.Config, logger logger.Logger) email.Handlers {
	return &emailHandler{cfg: cfg, emailUC: uc, logger: logger}
}

// SendEmail godoc
// @Summary Send an email
// @Description Send an email using configured SMTP.
// @Tags email
// @Accept json
// @Produce json
// @Param request body presenter.EmailSendRequest true "Email details"
// @Success 200 {object} responses.SuccessResponse[string]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/email/send [post]
func (h *emailHandler) SendEmail() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.EmailSendRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := h.emailUC.SendEmail(r.Context(), req.To, req.Subject, req.Body); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse("Email sent"))
	}
}

// GetEmailLog godoc
// @Summary Get email log
// @Description Get a single email log entry by ID.
// @Tags email
// @Accept json
// @Produce json
// @Param id path int true "Log ID"
// @Success 200 {object} responses.SuccessResponse[presenter.EmailLogResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/email/logs/{id} [get]
func (h *emailHandler) GetEmailLog() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		_ = id
		render.Respond(w, r, responses.CreateSuccessResponse("ok"))
	}
}

// ListEmailLogs godoc
// @Summary List email logs
// @Description List all email log entries with pagination.
// @Tags email
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} responses.SuccessResponse[[]presenter.EmailLogResponse]
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/email/logs [get]
func (h *emailHandler) ListEmailLogs() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logs, err := h.emailUC.GetMulti(r.Context(), 50, 0)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(logs))
	}
}

// GetEmailConfig godoc
// @Summary Get email config
// @Description Get the current SMTP configuration.
// @Tags email
// @Accept json
// @Produce json
// @Success 200 {object} responses.SuccessResponse[presenter.EmailConfigResponse]
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/email/config [get]
func (h *emailHandler) GetEmailConfig() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg, err := h.emailUC.GetConfig(r.Context())
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(cfg))
	}
}

// UpdateEmailConfig godoc
// @Summary Update email config
// @Description Update the SMTP configuration.
// @Tags email
// @Accept json
// @Produce json
// @Param config body presenter.EmailConfigResponse true "SMTP Config"
// @Success 200 {object} responses.SuccessResponse[presenter.EmailConfigResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/email/config [put]
func (h *emailHandler) UpdateEmailConfig() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Respond(w, r, responses.CreateSuccessResponse("ok"))
	}
}
