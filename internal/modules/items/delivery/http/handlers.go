package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"icmongolang/config"
	"icmongolang/internal/modules/items"
	"icmongolang/internal/modules/items/presenter"
	"icmongolang/internal/middleware"
	"icmongolang/internal/models"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"
	"icmongolang/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type itemHandler struct {
	cfg     *config.Config
	itemsUC items.ItemUseCaseI
	logger  logger.Logger
}

func CreateItemHandler(uc items.ItemUseCaseI, cfg *config.Config, logger logger.Logger) items.Handlers {
	return &itemHandler{cfg: cfg, itemsUC: uc, logger: logger}
}

// Create godoc
// @Summary Create Item
// @Description Create new item.
// @Tags items
// @Accept json
// @Produce json
// @Param item body presenter.ItemCreate true "Add item"
// @Success 200 {object} responses.SuccessResponse[presenter.ItemResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /item [post]
func (h *itemHandler) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		itemReq := new(presenter.ItemCreate)

		err := json.NewDecoder(r.Body).Decode(&itemReq)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		err = utils.ValidateStruct(ctx, itemReq)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		// ✅ สร้าง model และกำหนด OwnerId
		newItem := mapModel(itemReq)
		newItem.OwnerId = user.ID // กำหนดเจ้าของจาก user ที่ login

		createdItem, err := h.itemsUC.CreateWithOwner(ctx, user.ID, newItem)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		itemResponse := *mapModelResponse(createdItem)
		render.Respond(w, r, responses.CreateSuccessResponse(itemResponse))
	}
}

// Get godoc
// @Summary Read item
// @Description Get item by ID.
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item Id"
// @Success 200 {object} responses.SuccessResponse[presenter.ItemResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /item/{id} [get]
func (h *itemHandler) Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		item, err := h.itemsUC.Get(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		// เปลี่ยน user.Id เป็น user.ID และ item.OwnerId (คงไว้ตาม model item)
		if !user.IsSuperUser && item.OwnerId != user.ID {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotEnoughPrivileges(err)))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(item)))
	}
}

// GetMulti godoc
// @Summary Read Items
// @Description Retrieve items with pagination (page & per_page) or offset/limit.
// @Tags items
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)" Format(page)
// @Param per_page query int false "Items per page (default 10)" Format(per_page)
// @Success 200 {object} responses.SuccessResponse[presenter.PaginatedItemsResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /item [get]
func (h *itemHandler) GetMulti() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		q := r.URL.Query()

		// Default values
		limit := 10
		offset := 0
		page := 1
		perPage := 10

		// Check for page/per_page first (preferred)
		pageStr := q.Get("page")
		if pageStr == "" {
			pageStr = q.Get("offset") // treat limit as page only when page is missing
		}
		perPageStr := q.Get("per_page")
		if perPageStr == "" {
			perPageStr = q.Get("limit") // treat offset as per_page only when per_page missing
		}

		if pageStr != "" || perPageStr != "" {
			if pageStr != "" {
				if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
					page = p
				}
			}
			if perPageStr != "" {
				if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 {
					perPage = pp
				}
			}
			// Cap max per page
			const maxPerPage = 100
			if perPage > maxPerPage {
				perPage = maxPerPage
			}
			limit = perPage
			offset = (page - 1) * perPage
		} else {
			// Fallback to legacy offset/limit
			if l := q.Get("limit"); l != "" {
				if lim, err := strconv.Atoi(l); err == nil && lim > 0 {
					limit = lim
				}
			}
			if o := q.Get("offset"); o != "" {
				if off, err := strconv.Atoi(o); err == nil && off >= 0 {
					offset = off
				}
			}
			// Recalculate page/per_page from offset/limit for response consistency
			if limit > 0 {
				perPage = limit
				page = (offset / limit) + 1
			}
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		var items []*models.Item
		var total int64

		if user.IsSuperUser {
			items, err = h.itemsUC.GetMulti(ctx, limit, offset)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
			total, err = h.itemsUC.Count(ctx)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
		} else {
			items, err = h.itemsUC.GetMultiByOwnerId(ctx, user.ID, limit, offset)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
			total, err = h.itemsUC.CountByOwnerId(ctx, user.ID)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
		}

		// Calculate total pages
		totalPages := int(total) / perPage
		if int(total)%perPage != 0 {
			totalPages++
		}
		if totalPages < 1 {
			totalPages = 1
		}

		paginatedRes := &presenter.PaginatedItemsResponse{
			Items:      mapModelsResponse(items),
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages,
		}
		render.Respond(w, r, responses.CreateSuccessResponse(paginatedRes))
	}
}

// Delete godoc
// @Summary Delete item
// @Description Delete an item by ID.
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item Id"
// @Success 200 {object} responses.SuccessResponse[presenter.ItemResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /item/{id} [delete]
func (h *itemHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		item, err := h.itemsUC.Get(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		// เปลี่ยน user.Id เป็น user.ID
		if !user.IsSuperUser && item.OwnerId != user.ID {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotEnoughPrivileges(err)))
			return
		}

		err = h.itemsUC.DeleteWithoutGet(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(item)))
	}
}

// Update godoc
// @Summary Update item
// @Description Update an item by ID.
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item Id"
// @Param item body presenter.ItemUpdate true "Update item"
// @Success 200 {object} responses.SuccessResponse[presenter.ItemResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /item/{id} [put]
func (h *itemHandler) Update() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		item := new(presenter.ItemUpdate)

		err = json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		err = utils.ValidateStruct(r.Context(), item)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		dbItem, err := h.itemsUC.Get(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		// เปลี่ยน user.Id เป็น user.ID
		if !user.IsSuperUser && dbItem.OwnerId != user.ID {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotEnoughPrivileges(err)))
			return
		}

		values := make(map[string]interface{})
		if item.Title != "" {
			values["title"] = item.Title
		}
		if item.Description != "" {
			values["description"] = item.Description
		}

		updatedItem, err := h.itemsUC.Update(r.Context(), id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedItem)))
	}
}

func mapModel(exp *presenter.ItemCreate) *models.Item {
	return &models.Item{
		Title:       exp.Title,
		Description: exp.Description,
	}
}

func mapModelResponse(exp *models.Item) *presenter.ItemResponse {
	return &presenter.ItemResponse{
		Id:          exp.Id,
		Title:       exp.Title,
		Description: exp.Description,
		OwnerId:     exp.OwnerId,
	}
}

func mapModelsResponse(exp []*models.Item) []*presenter.ItemResponse {
	out := make([]*presenter.ItemResponse, len(exp))
	for i, user := range exp {
		out[i] = mapModelResponse(user)
	}
	return out
}
