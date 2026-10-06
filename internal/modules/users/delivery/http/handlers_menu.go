package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"icmongolang/internal/modules/users"
	usermodels "icmongolang/internal/modules/users/models"
	"icmongolang/internal/modules/users/presenter"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/responses"
	"icmongolang/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// CreateMenu – POST /user/menus
// @Security     BearerAuth
// @Summary      Create menu
// @Tags         menus
// @Accept       json
// @Produce      json
// @Param        body body presenter.MenuCreate true "Menu data"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=usermodels.SdMenu}
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/menus [post]
func (h *userHandler) CreateMenu() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.MenuCreate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))); return
		}
		m := &usermodels.SdMenu{
			MenuParentID: &req.MenuParentID,
			ModulesID:    req.ModulesID,
			Name:         req.Name,
			Path:         ptrStr(req.Path),
			Icon:         ptrStr(req.Icon),
			Sort:         req.Sort,
			Status:       req.Status,
		}
		out, err := h.menuUC.CreateMenu(r.Context(), m)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(out))
	}
}

// GetMenu – GET /user/menus/{id}
// @Security     BearerAuth
// @Summary      Get menu by ID
// @Tags         menus
// @Produce      json
// @Param        id path int true "Menu ID"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=usermodels.SdMenu}
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      404  {object}  responses.ErrorResponse
// @Router       /user/menus/{id} [get]
func (h *userHandler) GetMenu() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconvAtoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))); return
		}
		out, err := h.menuUC.GetMenu(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(out))
	}
}

// UpdateMenu – PUT /user/menus/{id}
// @Security     BearerAuth
// @Summary      Update menu
// @Tags         menus
// @Accept       json
// @Produce      json
// @Param        id path int true "Menu ID"
// @Param        body body presenter.MenuUpdate true "Fields to update"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=usermodels.SdMenu}
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/menus/{id} [put]
func (h *userHandler) UpdateMenu() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconvAtoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))); return
		}
		req := new(presenter.MenuUpdate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		values := map[string]interface{}{}
		if req.MenuParentID != nil { values["menu_parent_id"] = *req.MenuParentID }
		if req.ModulesID != nil    { values["modules_id"] = *req.ModulesID }
		if req.Name != nil         { values["name"] = *req.Name }
		if req.Path != nil         { values["path"] = *req.Path }
		if req.Icon != nil         { values["icon"] = *req.Icon }
		if req.Status != nil       { values["status"] = *req.Status }
		out, err := h.menuUC.UpdateMenu(r.Context(), id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(out))
	}
}

// DeleteMenu – DELETE /user/menus/{id}
// @Security     BearerAuth
// @Summary      Delete menu
// @Tags         menus
// @Param        id path int true "Menu ID"
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/menus/{id} [delete]
func (h *userHandler) DeleteMenu() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconvAtoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))); return
		}
		if err := h.menuUC.DeleteMenu(r.Context(), id); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(struct{}{}))
	}
}

// ListMenus – GET /user/menus
// @Security     BearerAuth
// @Summary      List menus
// @Tags         menus
// @Produce      json
// @Param        page           query int    false "Page" default(1)
// @Param        per_page       query int    false "Items per page" default(20)
// @Param        modules_id     query int    false "Filter by modules_id"
// @Param        menu_parent_id query int    false "Filter by menu_parent_id"
// @Param        keyword        query string false "Search keyword"
// @Param        status         query string false "Status csv"
// @Param        sort_by        query string false "Sort by field"
// @Param        sort_order     query string false "asc|desc"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=presenter.MenuListResponse}
// @Router       /user/menus [get]
func (h *userHandler) ListMenus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		page, perPage := parsePageParams(q)
		f := users.MenuFilter{
			Limit: perPage, Offset: (page - 1) * perPage,
			Keyword:   q.Get("keyword"),
			SortBy:    q.Get("sort_by"),
			SortOrder: q.Get("sort_order"),
		}
		if v := q.Get("modules_id"); v != "" {
			if n, err := strconvAtoi(v); err == nil { f.ModulesID = &n }
		}
		if v := q.Get("menu_parent_id"); v != "" {
			if n, err := strconvAtoi(v); err == nil { f.MenuParentID = &n }
		}
		if v := q.Get("status"); v != "" {
			for _, s := range strings.Split(v, ",") {
				if n, err := strconvAtoi16(strings.TrimSpace(s)); err == nil {
					f.Statuses = append(f.Statuses, n)
				}
			}
		}
		items, total, err := h.menuUC.ListMenus(r.Context(), f)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		list := make([]*presenter.MenuResponse, 0, len(items))
		for _, m := range items {
			list = append(list, &presenter.MenuResponse{
				MenuID:       m.MenuID,
				MenuParentID: derefInt(m.MenuParentID),
				ModulesID:    m.ModulesID,
				Name:         m.Name,
				Path:         derefStr(m.Path),
				Icon:         derefStr(m.Icon),
				Sort:         m.Sort,
				CreateDate:   m.CreateDate,
				UpdateDate:   m.UpdateDate,
				Status:       m.Status,
			})
		}
		render.Respond(w, r, responses.CreateSuccessResponse(
			presenter.NewMenuListResponse(list, page, perPage, total),
		))
	}
}

// BulkUpdateMenuSort – PATCH /user/menus/sort
// @Security     BearerAuth
// @Summary      Bulk update menu sort (drag & drop)
// @Tags         menus
// @Accept       json
// @Produce      json
// @Param        body body presenter.SortUpdateRequest true "Items to sort"
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/menus/sort [patch]
func (h *userHandler) BulkUpdateMenuSort() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.SortUpdateRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))); return
		}
		items := make([]users.SortItem, 0, len(req.Items))
		for _, it := range req.Items {
			items = append(items, users.SortItem{ID: it.ID, Sort: it.Sort, ParentID: it.ParentID})
		}
		if err := h.menuUC.BulkUpdateMenuSort(r.Context(), items); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(map[string]int{"updated": len(items)}))
	}
}