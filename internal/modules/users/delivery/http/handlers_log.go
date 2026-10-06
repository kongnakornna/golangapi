package http

import (
	"encoding/json"
	"net/http"
	"time"

	"icmongolang/internal/modules/users"
	usermodels "icmongolang/internal/modules/users/models"
	"icmongolang/internal/modules/users/presenter"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/responses"
	"icmongolang/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// CreateLog – POST /user/logs
// @Security     BearerAuth
// @Summary      Create user log (manual)
// @Tags         user-logs
// @Accept       json
// @Produce      json
// @Param        body body presenter.LogCreate true "Log data"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=usermodels.SdUserLog}
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/logs [post]
func (h *userHandler) CreateLog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.LogCreate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))); return
		}
		uidParsed, err := uuid.Parse(req.UID)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))); return
		}
		ip, device, browser := extractClientMeta(r)
		if req.IPAddress == nil    { req.IPAddress = &ip }
		if req.DeviceAccess == nil { req.DeviceAccess = &device }
		if req.BrowserAgent == nil { req.BrowserAgent = &browser }

		l := &usermodels.SdUserLog{
			LogTypeID:    req.LogTypeID,
			UID:          uidParsed,
			Name:         req.Name,
			Detail:       req.Detail,
			IPAddress:    req.IPAddress,
			DeviceAccess: req.DeviceAccess,
			BrowserAgent: req.BrowserAgent,
			Status:       ptrIntFromInt16(req.Status),
			Lang:         ptrStr(req.Lang),
		}
		out, err := h.logUC.CreateLog(r.Context(), l)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(out))
	}
}

// GetLog – GET /user/logs/{log_type_id}/{id}
// @Security     BearerAuth
// @Summary      Get log by composite ID
// @Tags         user-logs
// @Produce      json
// @Param        log_type_id path int true "Log type ID"
// @Param        id          path int true "Log ID"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=usermodels.SdUserLog}
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      404  {object}  responses.ErrorResponse
// @Router       /user/logs/{log_type_id}/{id} [get]
func (h *userHandler) GetLog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconvAtoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))); return
		}
		logTypeID, err := strconvAtoi(chi.URLParam(r, "log_type_id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))); return
		}
		out, err := h.logUC.GetLog(r.Context(), id, logTypeID)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(out))
	}
}

// ListLogs – GET /user/logs
// @Security     BearerAuth
// @Summary      List user logs
// @Tags         user-logs
// @Produce      json
// @Param        page        query int    false "Page" default(1)
// @Param        per_page    query int    false "Items per page" default(20)
// @Param        uid         query string false "Filter by uid"
// @Param        log_type_id query int    false "Filter by log_type_id"
// @Param        ipaddress   query string false "Filter by IP"
// @Param        from        query string false "From (RFC3339)"
// @Param        to          query string false "To (RFC3339)"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=presenter.LogListResponse}
// @Router       /user/logs [get]
func (h *userHandler) ListLogs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		page, perPage := parsePageParams(q)
		f := users.LogFilter{Limit: perPage, Offset: (page - 1) * perPage}
		if v := q.Get("uid"); v != "" {
			if uid, err := uuid.Parse(v); err == nil { f.UID = &uid }
		}
		if v := q.Get("log_type_id"); v != "" {
			if n, err := strconvAtoi(v); err == nil { f.LogTypeID = &n }
		}
		if v := q.Get("ipaddress"); v != "" { f.IPAddress = &v }
		if v := q.Get("from"); v != "" {
			if t, err := time.Parse(time.RFC3339, v); err == nil { f.From = &t }
		}
		if v := q.Get("to"); v != "" {
			if t, err := time.Parse(time.RFC3339, v); err == nil { f.To = &t }
		}
		items, total, err := h.logUC.ListLogs(r.Context(), f)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		list := make([]*presenter.LogResponse, 0, len(items))
		for _, it := range items {
			list = append(list, &presenter.LogResponse{
				ID:           it.ID,
				LogTypeID:    it.LogTypeID,
				UID:          it.UID.String(),
				Name:         it.Name,
				Detail:       it.Detail,
				IPAddress:    it.IPAddress,
				DeviceAccess: it.DeviceAccess,
				BrowserAgent: it.BrowserAgent,
				Status:       derefInt16FromIntPtr(it.Status),
				Create:       it.Create,
				Update:       it.Update,
				Lang:         derefStr(it.Lang),
			})
		}
		render.Respond(w, r, responses.CreateSuccessResponse(
			presenter.NewLogListResponse(list, page, perPage, total),
		))
	}
}

// DeleteLog – DELETE /user/logs/{log_type_id}/{id}
// @Security     BearerAuth
// @Summary      Delete log
// @Tags         user-logs
// @Param        log_type_id path int true "Log type ID"
// @Param        id          path int true "Log ID"
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/logs/{log_type_id}/{id} [delete]
func (h *userHandler) DeleteLog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconvAtoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))); return
		}
		logTypeID, err := strconvAtoi(chi.URLParam(r, "log_type_id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))); return
		}
		if err := h.logUC.DeleteLog(r.Context(), id, logTypeID); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(struct{}{}))
	}
}