package http

import (
	"net/http"

	"icmongolang/internal/modules/alarm/usecase"
	"icmongolang/pkg/helpers"
	"icmongolang/pkg/logger"

	"github.com/go-chi/render"
)

type AlarmHandler struct {
	uc     usecase.AlarmUseCase
	logger logger.Logger
}

func NewAlarmHandler(uc usecase.AlarmUseCase, log logger.Logger) *AlarmHandler {
	return &AlarmHandler{uc: uc, logger: log}
}

// ValidateAlarm godoc
// @Summary      Validate alarm status (default Thai)
// @Description  Calculates alarm level, title, subject, content and control messages based on sensor data and thresholds. Uses Thai language messages.
// @Tags         alarm
// @Accept       json
// @Produce      json
// @Param        request body helpers.AlarmDetailDto true "Alarm input data"
// @Success      200 {object} helpers.AlarmDetailResult
// @Failure      400 {object} errResponse "invalid JSON body"
// @Router       /alarm/validate [post]
// @Security     BearerAuth
func (h *AlarmHandler) ValidateAlarm(w http.ResponseWriter, r *http.Request) {
	var req helpers.AlarmDetailDto
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	result := h.uc.ValidateAlarm(req)
	render.JSON(w, r, result)
}

// ValidateAlarmEn godoc
// @Summary      Validate alarm status (English)
// @Description  Same as /validate but uses English messages.
// @Tags         alarm
// @Accept       json
// @Produce      json
// @Param        request body helpers.AlarmDetailDto true "Alarm input data"
// @Success      200 {object} helpers.AlarmDetailResult
// @Failure      400 {object} errResponse
// @Router       /alarm/validate/en [post]
// @Security     BearerAuth
func (h *AlarmHandler) ValidateAlarmEn(w http.ResponseWriter, r *http.Request) {
	var req helpers.AlarmDetailDto
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	result := h.uc.ValidateAlarmEn(req)
	render.JSON(w, r, result)
}

// ValidateAlarmTh godoc
// @Summary      Validate alarm status (Thai)
// @Description  Same as /validate but explicitly uses Thai messages.
// @Tags         alarm
// @Accept       json
// @Produce      json
// @Param        request body helpers.AlarmDetailDto true "Alarm input data"
// @Success      200 {object} helpers.AlarmDetailResult
// @Failure      400 {object} errResponse
// @Router       /alarm/validate/th [post]
// @Security     BearerAuth
func (h *AlarmHandler) ValidateAlarmTh(w http.ResponseWriter, r *http.Request) {
	var req helpers.AlarmDetailDto
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	result := h.uc.ValidateAlarmTh(req)
	render.JSON(w, r, result)
}

func ErrInvalidRequest(err error) render.Renderer {
	return &errResponse{HTTPStatusCode: http.StatusBadRequest, ErrorText: err.Error()}
}

type errResponse struct {
	HTTPStatusCode int    `json:"-"`
	ErrorText      string `json:"error"`
}

func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
