package http

import (
	"net/http"

	"icmongolang/internal/modules/realtime/presenter"
	"icmongolang/internal/modules/realtime/usecase"
	"icmongolang/pkg/logger"

	"github.com/go-chi/render"
)

// RealtimeHandler is the HTTP handler for real-time dashboard push.
type RealtimeHandler struct {
	uc     usecase.RealtimeUseCase
	logger logger.Logger
}

// CreateRealtimeHandler creates a new realtime HTTP handler.
func CreateRealtimeHandler(uc usecase.RealtimeUseCase, log logger.Logger) *RealtimeHandler {
	return &RealtimeHandler{uc: uc, logger: log}
}

// Publish godoc
// @Summary      Broadcast a real-time event to the dashboard
// @Description  Pushes device/sensor/alarm events to websocket rooms for live monitoring
// @Tags         realtime
// @Accept       json
// @Produce      json
// @Param        request body presenter.PublishRequest true "Realtime event"
// @Success      200 {object} presenter.PublishResponse
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /realtime/publish [post]
// @Security     BearerAuth
func (h *RealtimeHandler) Publish(w http.ResponseWriter, r *http.Request) {
	var req presenter.PublishRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	resp, err := h.uc.Publish(r.Context(), &req)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.JSON(w, r, resp)
}

// Health godoc
// @Summary      Realtime service health (broadcaster wired?)
// @Tags         realtime
// @Produce      json
// @Success      200 {object} presenter.HealthResponse
// @Router       /realtime/health [get]
// @Security     BearerAuth
func (h *RealtimeHandler) Health(w http.ResponseWriter, r *http.Request) {
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
