package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"icmongolang/internal/modules/users"
	"icmongolang/internal/modules/users/presenter"
	usermodels "icmongolang/internal/modules/users/models"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/responses"
	"icmongolang/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// CreateModule – POST /user/modules
// @Security     BearerAuth
// @Summary      Create module
// @Tags         modules
// @Accept       json
// @Produce      json
// @Param        body body presenter.ModuleCreate true "Module data"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=usermodels.SdModules}
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/modules [post]
func (h *userHandler) CreateModule() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.ModuleCreate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		m := &usermodels.SdModules{
			ModuleParentID: &req.ModuleParentID,
			Name:           req.Name,
			Path:           ptrStr(req.Path),
			Icon:           ptrStr(req.Icon),
			Sort:           req.Sort,
			Status:         req.Status,
		}
		out, err := h.moduleUC.CreateModule(r.Context(), m)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(out))
	}
}

// GetModule – GET /user/modules/{id}
// @Security     BearerAuth
// @Summary      Get module by ID
// @Tags         modules
// @Produce      json
// @Param        id path int true "Module ID"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=usermodels.SdModules}
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      404  {object}  responses.ErrorResponse
// @Router       /user/modules/{id} [get]
func (h *userHandler) GetModule() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconvAtoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		out, err := h.moduleUC.GetModule(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(out))
	}
}

// UpdateModule – PUT /user/modules/{id}
// @Security     BearerAuth
// @Summary      Update module
// @Tags         modules
// @Accept       json
// @Produce      json
// @Param        id path int true "Module ID"
// @Param        body body presenter.ModuleUpdate true "Fields to update"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=usermodels.SdModules}
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/modules/{id} [put]
func (h *userHandler) UpdateModule() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconvAtoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		req := new(presenter.ModuleUpdate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		values := map[string]interface{}{}
		if req.ModuleParentID != nil {
			values["module_parent_id"] = *req.ModuleParentID
		}
		if req.Name != nil {
			values["name"] = *req.Name
		}
		if req.Path != nil {
			values["path"] = *req.Path
		}
		if req.Icon != nil {
			values["icon"] = *req.Icon
		}
		if req.Status != nil {
			values["status"] = *req.Status
		}
		out, err := h.moduleUC.UpdateModule(r.Context(), id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(out))
	}
}

// DeleteModule – DELETE /user/modules/{id}
// @Security     BearerAuth
// @Summary      Delete module
// @Tags         modules
// @Param        id path int true "Module ID"
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/modules/{id} [delete]
func (h *userHandler) DeleteModule() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconvAtoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := h.moduleUC.DeleteModule(r.Context(), id); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(struct{}{}))
	}
}

// ListModules – GET /user/modules
// @Security     BearerAuth
// @Summary      List modules
// @Tags         modules
// @Produce      json
// @Param        page             query int    false "Page" default(1)
// @Param        per_page         query int    false "Items per page" default(20)
// @Param        module_parent_id query int    false "Filter by module_parent_id"
// @Param        keyword          query string false "Search keyword"
// @Param        status           query string false "Status csv"
// @Param        sort_by          query string false "Sort by field"
// @Param        sort_order       query string false "asc|desc"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=presenter.ModuleListResponse}
// @Router       /user/modules [get]
func (h *userHandler) ListModules() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		page, perPage := parsePageParams(q)
		f := users.ModuleFilter{
			Limit:     perPage,
			Offset:    (page - 1) * perPage,
			Keyword:   q.Get("keyword"),
			SortBy:    q.Get("sort_by"),
			SortOrder: q.Get("sort_order"),
		}
		if v := q.Get("module_parent_id"); v != "" {
			if n, err := strconvAtoi(v); err == nil {
				f.ModuleParentID = &n
			}
		}
		if v := q.Get("status"); v != "" {
			for _, s := range strings.Split(v, ",") {
				if n, err := strconvAtoi16(strings.TrimSpace(s)); err == nil {
					f.Statuses = append(f.Statuses, n)
				}
			}
		}
		items, total, err := h.moduleUC.ListModules(r.Context(), f)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		list := make([]*presenter.ModuleResponse, 0, len(items))
		for _, m := range items {
			list = append(list, &presenter.ModuleResponse{
				ModuleID:       m.ModuleID,
				ModuleParentID: derefInt(m.ModuleParentID),
				Name:           m.Name,
				Path:           derefStr(m.Path),
				Icon:           derefStr(m.Icon),
				Sort:           m.Sort,
				CreateDate:     m.CreateDate,
				UpdateDate:     m.UpdateDate,
				Status:         m.Status,
			})
		}
		render.Respond(w, r, responses.CreateSuccessResponse(
			presenter.NewModuleListResponse(list, page, perPage, total),
		))
	}
}

// BulkUpdateModuleSort – PATCH /user/modules/sort
// @Security     BearerAuth
// @Summary      Bulk update module sort (drag & drop)
// @Tags         modules
// @Accept       json
// @Produce      json
// @Param        body body presenter.SortUpdateRequest true "Items to sort"
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/modules/sort [patch]
func (h *userHandler) BulkUpdateModuleSort() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.SortUpdateRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		items := make([]users.SortItem, 0, len(req.Items))
		for _, it := range req.Items {
			items = append(items, users.SortItem{ID: it.ID, Sort: it.Sort, ParentID: it.ParentID})
		}
		if err := h.moduleUC.BulkUpdateModuleSort(r.Context(), items); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(map[string]int{"updated": len(items)}))
	}
}