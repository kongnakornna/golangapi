package http

import (
	"net/http"

	"icmongolang/internal/modules/notifier/presenter"
	"icmongolang/internal/modules/notifier/usecase"
	"icmongolang/pkg/logger"

	"github.com/go-chi/render"
)

// NotifierHandler is the HTTP handler for notification dispatch.
type NotifierHandler struct {
	uc     usecase.NotifierUseCase
	logger logger.Logger
}

// CreateNotifierHandler creates a new notifier HTTP handler.
func CreateNotifierHandler(uc usecase.NotifierUseCase, log logger.Logger) *NotifierHandler {
	return &NotifierHandler{uc: uc, logger: log}
}

// Dispatch godoc
// @Summary      Dispatch a notification through a channel
// @Description  Sends an alarm notification via the requested channel
// @Tags         notifier
// @Accept       json
// @Produce      json
// @Param        request body presenter.DispatchRequest true "Dispatch payload"
// @Success      200 {object} presenter.DispatchResponse
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /notifier/dispatch [post]
// @Security     BearerAuth
func (h *NotifierHandler) Dispatch(w http.ResponseWriter, r *http.Request) {
	var req presenter.DispatchRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	resp, err := h.uc.Dispatch(r.Context(), &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

// Health godoc
// @Summary      Notifier health check
// @Tags         notifier
// @Produce      json
// @Success      200 {object} presenter.HealthResponse
// @Router       /notifier/health [get]
// @Security     BearerAuth
func (h *NotifierHandler) Health(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, h.uc.Health(r.Context()))
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
