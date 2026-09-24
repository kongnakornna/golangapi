package http

import (
	"net/http"

	"icmongolang/internal/modules/apimanager/presenter"
	"icmongolang/internal/modules/apimanager/usecase"
	"icmongolang/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// APIManagerHandler is the HTTP handler for API management.
type APIManagerHandler struct {
	uc     usecase.APIManagerUseCase
	logger logger.Logger
}

// CreateAPIManagerHandler creates a new API manager HTTP handler.
func CreateAPIManagerHandler(uc usecase.APIManagerUseCase, log logger.Logger) *APIManagerHandler {
	return &APIManagerHandler{uc: uc, logger: log}
}

// CreateKey godoc
// @Summary      Create an API key
// @Tags         apimanager
// @Accept       json
// @Produce      json
// @Param        request body presenter.APIKeyRequest true "API key request"
// @Success      200 {object} presenter.APIKeyResponse
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /apimanager/key [post]
// @Security     BearerAuth
func (h *APIManagerHandler) CreateKey(w http.ResponseWriter, r *http.Request) {
	var req presenter.APIKeyRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	resp, err := h.uc.CreateKey(r.Context(), &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

// RevokeKey godoc
// @Summary      Revoke an API key
// @Tags         apimanager
// @Param        key_id path string true "Key ID"
// @Success      200 {object} map[string]string
// @Failure      500 {object} errResponse
// @Router       /apimanager/key/{key_id} [delete]
// @Security     BearerAuth
func (h *APIManagerHandler) RevokeKey(w http.ResponseWriter, r *http.Request) {
	keyID := chi.URLParam(r, "key_id")
	if err := h.uc.RevokeKey(r.Context(), keyID); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, map[string]string{"status": "ok"})
}

// CheckRateLimit godoc
// @Summary      Check current rate-limit state for a key
// @Tags         apimanager
// @Param        key_id path string true "Key ID"
// @Success      200 {object} presenter.RateLimitResponse
// @Failure      500 {object} errResponse
// @Router       /apimanager/ratelimit/{key_id} [get]
// @Security     BearerAuth
func (h *APIManagerHandler) CheckRateLimit(w http.ResponseWriter, r *http.Request) {
	keyID := chi.URLParam(r, "key_id")
	resp, err := h.uc.CheckRateLimit(r.Context(), keyID)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

// Usage godoc
// @Summary      Get API usage for a key
// @Tags         apimanager
// @Param        key_id path string true "Key ID"
// @Success      200 {array} presenter.UsageEntry
// @Failure      500 {object} errResponse
// @Router       /apimanager/usage/{key_id} [get]
// @Security     BearerAuth
func (h *APIManagerHandler) Usage(w http.ResponseWriter, r *http.Request) {
	keyID := chi.URLParam(r, "key_id")
	entries, err := h.uc.Usage(r.Context(), keyID)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, entries)
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
