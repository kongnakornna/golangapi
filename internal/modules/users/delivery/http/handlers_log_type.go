package http

import (
	"encoding/json"
	"net/http"

	usermodels "icmongolang/internal/modules/users/models"
	"icmongolang/internal/modules/users/presenter"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/responses"
	"icmongolang/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// CreateLogType – POST /user/log-types
// @Security     BearerAuth
// @Summary      Create log type
// @Tags         log-types
// @Accept       json
// @Produce      json
// @Param        body body presenter.LogTypeCreate true "Log type data"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=usermodels.SdUserLogType}
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/log-types [post]
func (h *userHandler) CreateLogType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.LogTypeCreate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))); return
		}
		lt := &usermodels.SdUserLogType{
			TypeName:   req.TypeName,
			TypeDetail: req.TypeDetail,
			Status:     ptrIntFromInt16(req.Status),
		}
		out, err := h.logTypeUC.CreateLogType(r.Context(), lt)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(out))
	}
}

// GetLogType – GET /user/log-types/{id}
// @Security     BearerAuth
// @Summary      Get log type by ID
// @Tags         log-types
// @Produce      json
// @Param        id path int true "Log type ID"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=usermodels.SdUserLogType}
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      404  {object}  responses.ErrorResponse
// @Router       /user/log-types/{id} [get]
func (h *userHandler) GetLogType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconvAtoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))); return
		}
		out, err := h.logTypeUC.GetLogType(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(out))
	}
}

// UpdateLogType – PUT /user/log-types/{id}
// @Security     BearerAuth
// @Summary      Update log type
// @Tags         log-types
// @Accept       json
// @Produce      json
// @Param        id path int true "Log type ID"
// @Param        body body presenter.LogTypeUpdate true "Fields to update"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=usermodels.SdUserLogType}
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/log-types/{id} [put]
func (h *userHandler) UpdateLogType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconvAtoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))); return
		}
		req := new(presenter.LogTypeUpdate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		values := map[string]interface{}{}
		if req.TypeName != nil   { values["type_name"] = *req.TypeName }
		if req.TypeDetail != nil { values["type_detail"] = *req.TypeDetail }
		if req.Status != nil     { values["status"] = *req.Status }
		out, err := h.logTypeUC.UpdateLogType(r.Context(), id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(out))
	}
}

// DeleteLogType – DELETE /user/log-types/{id}
// @Security     BearerAuth
// @Summary      Delete log type
// @Tags         log-types
// @Param        id path int true "Log type ID"
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/log-types/{id} [delete]
func (h *userHandler) DeleteLogType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconvAtoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))); return
		}
		if err := h.logTypeUC.DeleteLogType(r.Context(), id); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(struct{}{}))
	}
}

// ListLogTypes – GET /user/log-types
// @Security     BearerAuth
// @Summary      List log types
// @Tags         log-types
// @Produce      json
// @Param        page     query int false "Page" default(1)
// @Param        per_page query int false "Items per page" default(20)
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Router       /user/log-types [get]
func (h *userHandler) ListLogTypes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		page, perPage := parsePageParams(q)
		items, total, err := h.logTypeUC.ListLogTypes(r.Context(), perPage, (page-1)*perPage)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		resp := map[string]interface{}{
			"data":        items,
			"page":        page,
			"per_page":    perPage,
			"total":       total,
			"total_pages": int((total + int64(perPage) - 1) / int64(perPage)),
		}
		render.Respond(w, r, responses.CreateSuccessResponse(resp))
	}
}