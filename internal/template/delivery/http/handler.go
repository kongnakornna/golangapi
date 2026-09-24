package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"icmongolang/config"
	"icmongolang/internal/template"
	"icmongolang/internal/template/models"
	"icmongolang/internal/template/presenter"
	"icmongolang/pkg/helpers"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type templateHandler struct {
	cfg    *config.Config
	uc     template.UseCaseI
	logger logger.Logger
}

func NewTemplateHandler(uc template.UseCaseI, cfg *config.Config, log logger.Logger) template.Handlers {
	return &templateHandler{cfg: cfg, uc: uc, logger: log}
}

func (h *templateHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req presenter.CreateTemplateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrBadRequest(err)))
			return
		}

		model := &models.Template{Name: req.Name}
		created, err := h.uc.Create(r.Context(), model)
		if err != nil {
			h.logger.Errorf("template.Create failed: %v", err)
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		resp := toResponse(created)
		render.Render(w, r, responses.CreateSuccessResponse(resp))
	}
}

func (h *templateHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrBadRequest(err)))
			return
		}

		obj, err := h.uc.Get(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Render(w, r, responses.CreateSuccessResponse(toResponse(obj)))
	}
}

func (h *templateHandler) GetMulti() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		if limit <= 0 {
			limit = 50
		}

		objs, err := h.uc.GetMulti(r.Context(), limit, offset)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		items := make([]presenter.TemplateResponse, len(objs))
		for i, obj := range objs {
			items[i] = toResponse(obj)
		}
		render.Render(w, r, responses.CreateSuccessResponse(items))
	}
}

func (h *templateHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrBadRequest(err)))
			return
		}

		var req presenter.UpdateTemplateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrBadRequest(err)))
			return
		}

		values := make(map[string]interface{})
		if req.Name != "" {
			values["name"] = req.Name
		}
		if req.Status != nil {
			values["status"] = *req.Status
		}

		updated, err := h.uc.Update(r.Context(), id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Render(w, r, responses.CreateSuccessResponse(toResponse(updated)))
	}
}

func (h *templateHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrBadRequest(err)))
			return
		}

		obj, err := h.uc.Delete(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Render(w, r, responses.CreateSuccessResponse(toResponse(obj)))
	}
}

func toResponse(m *models.Template) presenter.TemplateResponse {
	return presenter.TemplateResponse{
		ID:        m.ID.String(),
		Name:      m.Name,
		Status:    m.Status,
		CreatedAt: helpers.NewLocalTime(m.CreatedAt),
		UpdatedAt: helpers.NewLocalTime(m.UpdatedAt),
	}
}
