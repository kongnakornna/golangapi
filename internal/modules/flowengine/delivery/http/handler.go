package http

import (
	"net/http"

	"icmongolang/internal/modules/flowengine/presenter"
	"icmongolang/internal/modules/flowengine/usecase"
	"icmongolang/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// FlowEngineHandler is the HTTP handler for workflow diagrams (Node-RED-like).
type FlowEngineHandler struct {
	uc     usecase.FlowEngineUseCase
	logger logger.Logger
}

// CreateFlowEngineHandler creates a new flow engine HTTP handler.
func CreateFlowEngineHandler(uc usecase.FlowEngineUseCase, log logger.Logger) *FlowEngineHandler {
	return &FlowEngineHandler{uc: uc, logger: log}
}

// Create godoc
// @Summary      Create or update a workflow flow
// @Description  Saves a Node-RED-like flow (nodes + edges)
// @Tags         flowengine
// @Accept       json
// @Produce      json
// @Param        request body presenter.CreateRequest true "Flow definition"
// @Success      200 {object} presenter.Flow
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /flowengine/flow [post]
// @Security     BearerAuth
func (h *FlowEngineHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req presenter.CreateRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	flow, err := h.uc.Create(r.Context(), &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, flow)
}

// Get godoc
// @Summary      Get a workflow flow
// @Tags         flowengine
// @Param        id path string true "Flow ID"
// @Success      200 {object} presenter.Flow
// @Failure      500 {object} errResponse
// @Router       /flowengine/flow/{id} [get]
// @Security     BearerAuth
func (h *FlowEngineHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	flow, err := h.uc.Get(r.Context(), id)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, flow)
}

// Validate godoc
// @Summary      Validate a workflow graph
// @Tags         flowengine
// @Accept       json
// @Produce      json
// @Param        request body presenter.CreateRequest true "Flow definition"
// @Success      200 {object} presenter.ValidateResponse
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /flowengine/validate [post]
// @Security     BearerAuth
func (h *FlowEngineHandler) Validate(w http.ResponseWriter, r *http.Request) {
	var req presenter.CreateRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	resp, err := h.uc.Validate(r.Context(), &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

// Trigger godoc
// @Summary      Trigger a workflow execution
// @Tags         flowengine
// @Accept       json
// @Produce      json
// @Param        request body presenter.TriggerRequest true "Trigger payload"
// @Success      200 {object} presenter.TriggerResponse
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /flowengine/trigger [post]
// @Security     BearerAuth
func (h *FlowEngineHandler) Trigger(w http.ResponseWriter, r *http.Request) {
	var req presenter.TriggerRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	resp, err := h.uc.Trigger(r.Context(), &req)
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
