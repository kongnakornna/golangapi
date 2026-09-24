package http

import (
	"net/http"

	"icmongolang/internal/modules/auditlog/presenter"
	"icmongolang/internal/modules/auditlog/usecase"
	"icmongolang/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// AuditLogHandler is the HTTP handler for audit logs + blockchain anchoring.
type AuditLogHandler struct {
	uc     usecase.AuditLogUseCase
	logger logger.Logger
}

// CreateAuditLogHandler creates a new audit log HTTP handler.
func CreateAuditLogHandler(uc usecase.AuditLogUseCase, log logger.Logger) *AuditLogHandler {
	return &AuditLogHandler{uc: uc, logger: log}
}

// Create godoc
// @Summary      Record an audit log entry
// @Description  Creates a new audit entry (hash-chained for tamper-evidence)
// @Tags         auditlog
// @Accept       json
// @Produce      json
// @Param        request body presenter.CreateRequest true "Audit entry payload"
// @Success      200 {object} presenter.AuditEntry
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /auditlog [post]
// @Security     BearerAuth
func (h *AuditLogHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req presenter.CreateRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	entry, err := h.uc.Create(r.Context(), &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, entry)
}

// Verify godoc
// @Summary      Verify an audit log entry hash-chain
// @Description  Checks that an entry's hash/prev-hash chain is intact
// @Tags         auditlog
// @Produce      json
// @Param        id path string true "Audit entry ID"
// @Success      200 {object} presenter.VerifyResponse
// @Failure      500 {object} errResponse
// @Router       /auditlog/verify/{id} [get]
// @Security     BearerAuth
func (h *AuditLogHandler) Verify(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	resp, err := h.uc.Verify(r.Context(), id)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

// Chain godoc
// @Summary      Get current audit chain tail
// @Tags         auditlog
// @Produce      json
// @Success      200 {object} presenter.ChainResponse
// @Failure      500 {object} errResponse
// @Router       /auditlog/chain [get]
// @Security     BearerAuth
func (h *AuditLogHandler) Chain(w http.ResponseWriter, r *http.Request) {
	resp, err := h.uc.Chain(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

func ErrInvalidRequest(err error) render.Renderer {
	return &errResponse{HTTPStatusCode: http.StatusBadRequest, ErrorText: err.Error()}
}

func ErrInternal(err error) render.Renderer {
	return &errResponse{HTTPStatusCode: http.StatusInternalServerError, ErrorText: err.Error()}
}

type errResponse struct {
	HTTPStatusCode int    `json:"-"`
	ErrorText      string `json:"error"`
}

func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
