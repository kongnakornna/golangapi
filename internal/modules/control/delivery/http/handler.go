package http

import (
	"net/http"

	"icmongolang/internal/modules/control/presenter"
	"icmongolang/internal/modules/control/usecase"
	"icmongolang/pkg/logger"

	"github.com/go-chi/render"
)

// ControlHandler is the HTTP handler for auto-control actions.
type ControlHandler struct {
	uc     usecase.ControlUseCase
	logger logger.Logger
}

// CreateControlHandler creates a new control HTTP handler.
func CreateControlHandler(uc usecase.ControlUseCase, log logger.Logger) *ControlHandler {
	return &ControlHandler{uc: uc, logger: log}
}

// ExecuteControl godoc
// @Summary      Execute a device control (ON/OFF)
// @Description  Publishes an MQTT control message for a device immediately
// @Tags         control
// @Accept       json
// @Produce      json
// @Param        request body presenter.ExecuteRequest true "Control payload"
// @Success      200 {object} presenter.ExecuteResult
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /control/execute [post]
// @Security     BearerAuth
func (h *ControlHandler) ExecuteControl(w http.ResponseWriter, r *http.Request) {
	var req presenter.ExecuteRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	resp, err := h.uc.ExecuteControl(r.Context(), &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

// RunSchedule godoc
// @Summary      Run active control schedules
// @Description  Executes schedules matched for today/time (or forced)
// @Tags         control
// @Accept       json
// @Produce      json
// @Param        request body presenter.ScheduleRunRequest true "Schedule run scoping"
// @Success      200 {object} presenter.ScheduleRunResult
// @Failure      500 {object} errResponse
// @Router       /control/schedule/run [post]
// @Security     BearerAuth
func (h *ControlHandler) RunSchedule(w http.ResponseWriter, r *http.Request) {
	var req presenter.ScheduleRunRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	resp, err := h.uc.RunSchedule(r.Context(), &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

// Health godoc
// @Summary      Control health check
// @Tags         control
// @Produce      json
// @Success      200 {object} presenter.HealthResponse
// @Router       /control/health [get]
// @Security     BearerAuth
func (h *ControlHandler) Health(w http.ResponseWriter, r *http.Request) {
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
