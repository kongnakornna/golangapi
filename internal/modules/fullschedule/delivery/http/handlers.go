package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"icmongolang/config"
	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/fullschedule"
	"icmongolang/internal/modules/fullschedule/presenter"
	iotmodels "icmongolang/internal/modules/iot/models"
	"icmongolang/pkg/helpers"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"
	"icmongolang/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type fullscheduleHandler struct {
	cfg    *config.Config
	uc     fullschedule.FullScheduleUseCaseI
	logger logger.Logger
}

// CreateFullscheduleHandler creates the HTTP handler.
// สร้างฮาร์ดเลอร์สำหรับโมดูล Full Schedule
func CreateFullscheduleHandler(uc fullschedule.FullScheduleUseCaseI, cfg *config.Config, logger logger.Logger) fullschedule.Handlers {
	return &fullscheduleHandler{cfg: cfg, uc: uc, logger: logger}
}

// Create godoc
// @Summary Create a fullschedule
// @Description Create a new schedule (normal/full/batch).
// @Tags fullschedule
// @Accept json
// @Produce json
// @Param schedule body presenter.ScheduleCreate true "Schedule details"
// @Success 200 {object} responses.SuccessResponse[presenter.ScheduleResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule [post]
func (h *fullscheduleHandler) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := new(presenter.ScheduleCreate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := utils.ValidateStruct(ctx, req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		model := mapCreateRequest(req, &user.ID)
		created, err := h.uc.Create(ctx, model)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapScheduleResponse(created)))
	}
}

// GetByID godoc
// @Summary Get a fullschedule
// @Description Get a schedule by ID.
// @Tags fullschedule
// @Accept json
// @Produce json
// @Param id path string true "Schedule ID"
// @Success 200 {object} responses.SuccessResponse[presenter.ScheduleResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/{id} [get]
func (h *fullscheduleHandler) GetByID() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		s, err := h.uc.Get(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapScheduleResponse(s)))
	}
}

// List godoc
// @Summary List fullschedules
// @Description List schedules with pagination.
// @Tags fullschedule
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} responses.SuccessResponse[presenter.PaginatedScheduleResponse]
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule [get]
func (h *fullscheduleHandler) List() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		limit, offset, page, perPage := parsePagination(r)
		schedules, err := h.uc.GetMulti(ctx, limit, offset)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		total, err := h.uc.Count(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		totalPages := int(total) / perPage
		if int(total)%perPage != 0 {
			totalPages++
		}
		if totalPages < 1 {
			totalPages = 1
		}
		res := &presenter.PaginatedScheduleResponse{
			Schedules:  mapSchedulesResponse(schedules),
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages,
		}
		render.Respond(w, r, responses.CreateSuccessResponse(res))
	}
}

// ListSchedulePage godoc
// @Summary List schedules (settings-style page)
// @Description Settings-style paginated list: ?page=&pageSize=&start=&keyword=&mode=&status=&event_type=&event_action=&group_id=&zone_id=&area_id=&sort=. Mirrors /api/settings/listschedulepage.
// @Tags fullschedule
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Items per page"
// @Param start query string false "Filter by start time"
// @Param keyword query string false "Search by schedule name"
// @Param mode query string false "Schedule mode" Enums(normal,full,batch)
// @Param status query string false "Filter by status" Enums(active,inactive,draft)
// @Param event_type query string false "Filter by event type"
// @Param event_action query string false "Filter by event action"
// @Param group_id query string false "Filter by group id (uuid)"
// @Param zone_id query string false "Filter by zone id (uuid)"
// @Param area_id query string false "Filter by area id (uuid)"
// @Param sort query string false "Sort e.g. start-DESC"
// @Success 200 {object} responses.SuccessResponse[presenter.ListResult]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/listschedulepage [get]
func (h *fullscheduleHandler) ListSchedulePage() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		q := r.URL.Query()
		filter, err := parseScheduleFilter(q)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		page, pageSize := parseListPaging(q)
		list, total, err := h.uc.ListPage(ctx, filter, page, pageSize)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		res := presenter.NewScheduleListResult(mapSchedulesResponse(list), total, page, pageSize, filter.FilterMap())
		render.Respond(w, r, responses.CreateSuccessResponse(res))
	}
}

// ListScheduleAll godoc
// @Summary List all schedules (no pagination)
// @Description Settings-style full list (like /api/settings/scheduleall) with the same filters as listschedulepage.
// @Tags fullschedule
// @Accept json
// @Produce json
// @Param start query string false "Filter by start time"
// @Param keyword query string false "Search by schedule name"
// @Param mode query string false "Schedule mode" Enums(normal,full,batch)
// @Param status query string false "Filter by status" Enums(active,inactive,draft)
// @Param event_type query string false "Filter by event type"
// @Param event_action query string false "Filter by event action"
// @Param sort query string false "Sort e.g. start-DESC"
// @Success 200 {object} responses.SuccessResponse[[]presenter.ScheduleResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/scheduleall [get]
func (h *fullscheduleHandler) ListScheduleAll() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		filter, err := parseScheduleFilter(r.URL.Query())
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		list, err := h.uc.ListAll(ctx, filter)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapSchedulesResponse(list)))
	}
}

// ListScheduleDevicePage godoc
// @Summary List schedule-device mappings (settings-style page)
// @Description Settings-style paginated list of schedule-device mappings (like /api/settings/listscheduledevice).
// @Tags fullschedule
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Items per page"
// @Param keyword query string false "Search by schedule name or device name"
// @Param schedule_id query string false "Filter by schedule id (uuid)"
// @Param device_id query int false "Filter by device id"
// @Param status query string false "Filter by schedule status"
// @Param sort query string false "Sort e.g. start-DESC"
// @Success 200 {object} responses.SuccessResponse[presenter.ListResult]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/scheduledevicepage [get]
func (h *fullscheduleHandler) ListScheduleDevicePage() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		q := r.URL.Query()
		filter, err := parseScheduleDeviceFilter(q)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		page, pageSize := parseListPaging(q)
		rows, total, err := h.uc.ListScheduleDevicePage(ctx, filter, page, pageSize)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		data := make([]map[string]interface{}, 0, len(rows))
		for _, row := range rows {
			item := map[string]interface{}{
				"schedule_id":   row.ScheduleID,
				"schedule_name": row.ScheduleName,
				"mode":          row.Mode,
				"status":        row.Status,
				"start":         row.Start,
				"event_action":  row.EventAction,
				"device_id":     row.DeviceID,
			}
			if row.DeviceName != nil {
				item["device_name"] = *row.DeviceName
			}
			if row.SN != nil {
				item["sn"] = *row.SN
			}
			data = append(data, item)
		}
		res := presenter.NewScheduleListResult(data, total, page, pageSize, filter.FilterMap())
		render.Respond(w, r, responses.CreateSuccessResponse(res))
	}
}

// parseListPaging reads page/pageSize (accepts per_page as alias, max 100).
// อ่านค่า page/pageSize จาก query (รองรับ per_page ด้วย สูงสุด 100)
func parseListPaging(q map[string][]string) (page, pageSize int) {
	page, pageSize = 1, 10
	if raw := firstParam(q, "page"); raw != "" {
		if p, err := strconv.Atoi(raw); err == nil && p > 0 {
			page = p
		}
	}
	if raw := firstParam(q, "pageSize"); raw != "" {
		if ps, err := strconv.Atoi(raw); err == nil && ps > 0 {
			pageSize = ps
		}
	} else if raw := firstParam(q, "per_page"); raw != "" {
		if ps, err := strconv.Atoi(raw); err == nil && ps > 0 {
			pageSize = ps
		}
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

// parseStatusValue parses a status query value ("1"/"active"->1, "3"/"draft"->3, else 0).
// แปลงค่า status จาก query (1/0/3 หรือ active/inactive/draft) เป็น int8
func parseStatusValue(raw string) int8 {
	switch raw {
	case "1", "active":
		return 1
	case "3", "draft":
		return 3
	default:
		return 0
	}
}

// parseScheduleFilter builds a ScheduleFilter from query params.
// สร้าง ScheduleFilter จาก query params
func parseScheduleFilter(q map[string][]string) (fullschedule.ScheduleFilter, error) {
	var f fullschedule.ScheduleFilter
	f.Keyword = firstParam(q, "keyword")
	f.Start = firstParam(q, "start")
	f.Sort = firstParam(q, "sort")
	if f.Sort != "" {
		if _, ok := fullschedule.ParseScheduleSort(f.Sort); !ok {
			return f, fullschedule.ErrInvalidSort
		}
	}
	if raw := firstParam(q, "mode"); raw != "" {
		m := fullschedule.ScheduleMode(raw)
		f.Mode = &m
	}
	if raw := firstParam(q, "status"); raw != "" {
		v := fullschedule.ScheduleStatusValue(parseStatusValue(raw))
		f.Status = &v
	}
	if raw := firstParam(q, "event_type"); raw != "" {
		e := fullschedule.EventType(raw)
		f.EventType = &e
	}
	if raw := firstParam(q, "event_action"); raw != "" {
		a := fullschedule.EventAction(raw)
		f.EventAction = &a
	}
	var err error
	f.Scope, err = parseScopeQuery(q)
	if err != nil {
		return f, err
	}
	return f, nil
}

// parseScheduleDeviceFilter builds a ScheduleDeviceFilter from query params.
// สร้าง ScheduleDeviceFilter จาก query params
func parseScheduleDeviceFilter(q map[string][]string) (fullschedule.ScheduleDeviceFilter, error) {
	var f fullschedule.ScheduleDeviceFilter
	f.Keyword = firstParam(q, "keyword")
	f.Sort = firstParam(q, "sort")
	if f.Sort != "" {
		if _, ok := fullschedule.ParseScheduleDeviceSort(f.Sort); !ok {
			return f, fullschedule.ErrInvalidSort
		}
	}
	if raw := firstParam(q, "schedule_id"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			return f, err
		}
		f.ScheduleID = &id
	}
	if raw := firstParam(q, "device_id"); raw != "" {
		d, err := strconv.Atoi(raw)
		if err != nil {
			return f, err
		}
		f.DeviceID = &d
	}
	if raw := firstParam(q, "status"); raw != "" {
		v := fullschedule.ScheduleStatusValue(parseStatusValue(raw))
		f.Status = &v
	}
	return f, nil
}

// Update godoc
// @Summary Update a fullschedule
// @Description Update a schedule by ID.
// @Tags fullschedule
// @Accept json
// @Produce json
// @Param id path string true "Schedule ID"
// @Param schedule body presenter.ScheduleUpdate true "Schedule update fields"
// @Success 200 {object} responses.SuccessResponse[presenter.ScheduleResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/{id} [put]
func (h *fullscheduleHandler) Update() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		req := new(presenter.ScheduleUpdate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		values := make(map[string]interface{})
		if req.Name != nil {
			values["name"] = *req.Name
		}
		if req.Mode != nil {
			values["mode"] = *req.Mode
		}
		if req.Status != nil {
			values["status"] = *req.Status
		}
		if req.TimeStart != nil {
			values["time_start"] = *req.TimeStart
		}
		if req.EventType != nil {
			values["event_type"] = *req.EventType
		}
		if req.EventAction != nil {
			values["event_action"] = *req.EventAction
		}
		if req.Event != nil {
			values["event"] = *req.Event
		}
		if req.Sunday != nil {
			values["sunday"] = *req.Sunday
		}
		if req.Monday != nil {
			values["monday"] = *req.Monday
		}
		if req.Tuesday != nil {
			values["tuesday"] = *req.Tuesday
		}
		if req.Wednesday != nil {
			values["wednesday"] = *req.Wednesday
		}
		if req.Thursday != nil {
			values["thursday"] = *req.Thursday
		}
		if req.Friday != nil {
			values["friday"] = *req.Friday
		}
		if req.Saturday != nil {
			values["saturday"] = *req.Saturday
		}
		if req.Months != nil {
			values["months"] = req.Months
		}
		if req.Dates != nil {
			values["dates"] = req.Dates
		}
		if req.CronExpr != nil {
			values["cron_expr"] = req.CronExpr
		}
		if req.DeviceIDs != nil {
			values["_device_ids"] = req.DeviceIDs
		}
		if req.GroupID != nil {
			values["group_id"] = *req.GroupID
		}
		if req.ZoneID != nil {
			values["zone_id"] = *req.ZoneID
		}
		if req.AreaID != nil {
			values["area_id"] = *req.AreaID
		}

		updated, err := h.uc.Update(ctx, id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapScheduleResponse(updated)))
	}
}

// Delete godoc
// @Summary Delete a fullschedule
// @Description Soft-delete a schedule by ID.
// @Tags fullschedule
// @Accept json
// @Produce json
// @Param id path string true "Schedule ID"
// @Success 200 {object} responses.SuccessResponse[presenter.ScheduleResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/{id} [delete]
func (h *fullscheduleHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		s, err := h.uc.Delete(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapScheduleResponse(s)))
	}
}

// SetStatus godoc
// @Summary Change fullschedule status
// @Description Activate or deactivate a schedule.
// @Tags fullschedule
// @Accept json
// @Produce json
// @Param id path string true "Schedule ID"
// @Param status body presenter.ScheduleStatusChange true "New status"
// @Success 200 {object} responses.SuccessResponse[presenter.ScheduleResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/{id}/status [put]
func (h *fullscheduleHandler) SetStatus() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		req := new(presenter.ScheduleStatusChange)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := utils.ValidateStruct(ctx, req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		updated, err := h.uc.SetStatus(ctx, id, fullschedule.ScheduleStatusValue(req.Status))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapScheduleResponse(updated)))
	}
}

// UpdateDayStatus godoc
// @Summary Toggle a single weekday flag on a fullschedule
// @Description Activate or deactivate one weekday (sunday..saturday) of a schedule.
// @Tags fullschedule
// @Accept json
// @Produce json
// @Param id path string true "Schedule ID"
// @Param body body presenter.ScheduleDayStatusChange true "Day + value"
// @Success 200 {object} responses.SuccessResponse[presenter.ScheduleResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/{id}/day-status [put]
func (h *fullscheduleHandler) UpdateDayStatus() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		req := new(presenter.ScheduleDayStatusChange)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := utils.ValidateStruct(ctx, req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		updated, err := h.uc.Update(ctx, id, map[string]interface{}{req.Day: req.Value})
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapScheduleResponse(updated)))
	}
}

// Trigger godoc
// @Summary Manually trigger a fullschedule
// @Description Force execution now, bypassing day/time/month conditions. Only applies to normal/full modes.
// @Tags fullschedule
// @Accept json
// @Produce json
// @Param id path string true "Schedule ID"
// @Success 200 {object} responses.SuccessResponse[presenter.TriggerResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/{id}/trigger [post]
func (h *fullscheduleHandler) Trigger() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		h.logger.Infof("fullschedule: manual trigger %s", id)
		history, err := h.uc.Trigger(ctx, id, fullschedule.TriggeredManual)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		res := &presenter.TriggerResponse{
			Message:   "manual trigger executed",
			HistoryID: history.ID,
		}
		render.Respond(w, r, responses.CreateSuccessResponse(res))
	}
}

// SetSettings godoc
// @Summary Upsert a fullschedule setting
// @Description Set a per-schedule setting (e.g. email_recipients).
// @Tags fullschedule
// @Accept json
// @Produce json
// @Param id path string true "Schedule ID"
// @Param setting body presenter.ScheduleSettingRequest true "Setting"
// @Success 200 {object} responses.SuccessResponse[presenter.SettingResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/{id}/settings [post]
func (h *fullscheduleHandler) SetSettings() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		req := new(presenter.ScheduleSettingRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := utils.ValidateStruct(ctx, req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		setting, err := h.uc.SetSetting(ctx, id, req.Key, []byte(req.Value))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapSettingResponse(setting)))
	}
}

// GetSettings godoc
// @Summary Get fullschedule settings
// @Description Get all settings of a schedule.
// @Tags fullschedule
// @Accept json
// @Produce json
// @Param id path string true "Schedule ID"
// @Success 200 {object} responses.SuccessResponse[[]presenter.SettingResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/{id}/settings [get]
func (h *fullscheduleHandler) GetSettings() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		settings, err := h.uc.GetSettings(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapSettingsResponse(settings)))
	}
}

// GetHistoryBySchedule godoc
// @Summary Get fullschedule history
// @Description Get execution history of one schedule.
// @Tags fullschedule
// @Accept json
// @Produce json
// @Param id path string true "Schedule ID"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} responses.SuccessResponse[presenter.PaginatedHistoryResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/{id}/history [get]
func (h *fullscheduleHandler) GetHistoryBySchedule() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		limit, offset, _, _ := parsePagination(r)
		rows, err := h.uc.GetHistoryBySchedule(r.Context(), id, limit, offset)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		res := &presenter.PaginatedHistoryResponse{
			History: mapHistoryResponse(rows),
			Total:   int64(len(rows)),
			Page:    1,
			PerPage: limit,
		}
		render.Respond(w, r, responses.CreateSuccessResponse(res))
	}
}

// GetHistory godoc
// @Summary Get fullschedule history (all)
// @Description Get execution history across all schedules.
// @Tags fullschedule
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Param status query string false "Filter by status" Enums(success,failed,skipped,processing)
// @Success 200 {object} responses.SuccessResponse[presenter.PaginatedHistoryResponse]
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/history [get]
func (h *fullscheduleHandler) GetHistory() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset, _, _ := parsePagination(r)
		q := r.URL.Query()
		var status *fullschedule.HistoryStatus
		if s := q.Get("status"); s != "" {
			st := fullschedule.HistoryStatus(s)
			status = &st
		}
		rows, err := h.uc.GetHistory(r.Context(), limit, offset, status)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		res := &presenter.PaginatedHistoryResponse{
			History: mapHistoryResponse(rows),
			Total:   int64(len(rows)),
			Page:    1,
			PerPage: limit,
		}
		render.Respond(w, r, responses.CreateSuccessResponse(res))
	}
}

// GetReport godoc
// @Summary Get fullschedule report
// @Description Aggregate success/failure report by schedule.
// @Tags fullschedule
// @Accept json
// @Produce json
// @Param from query string false "From date (RFC3339)"
// @Param to query string false "To date (RFC3339)"
// @Param schedule_id query string false "Filter by schedule id"
// @Success 200 {object} responses.SuccessResponse[[]presenter.ReportResponse]
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/report [get]
func (h *fullscheduleHandler) GetReport() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		q := r.URL.Query()
		from, to, err := parseRange(q)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		var scheduleID *uuid.UUID
		if raw := q.Get("schedule_id"); raw != "" {
			parsed, perr := uuid.Parse(raw)
			if perr != nil {
				render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(perr)))
				return
			}
			scheduleID = &parsed
		}
		scope, err := parseScopeQuery(q)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		rows, err := h.uc.GetReport(ctx, from, to, scheduleID, scope)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapReportResponse(rows)))
	}
}

// ListDevices godoc
// @Summary List IoT devices
// @Description List available IoT devices for mapping (device picker).
// @Tags fullschedule
// @Accept json
// @Produce json
// @Param keyword query string false "Search by device name or SN"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} responses.SuccessResponse[[]presenter.IoTDeviceBrief]
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/devices [get]
func (h *fullscheduleHandler) ListDevices() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		limit, offset, _, _ := parsePagination(r)
		page := (offset / limit) + 1
		filter := map[string]interface{}{}
		if k := r.URL.Query().Get("keyword"); k != "" {
			filter["keyword"] = k
		}
		provider, ok := h.uc.(fullschedule.DeviceListProvider)
		if !ok {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrBadRequest(fullschedule.ErrDeviceStoreNotWired)))
			return
		}
		devices, _, err := provider.ListDevices(ctx, filter, page, limit)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapDevicesBrief(devices)))
	}
}

// parseSortOrder converts the sort query param into a SortOrder. Invalid or
// empty values fall back to ascending (default).
func parseSortOrder(raw string) fullschedule.SortOrder {
	if strings.EqualFold(strings.TrimSpace(raw), "desc") {
		return fullschedule.SortOrderDESC
	}
	return fullschedule.SortOrderASC
}

func parsePagination(r *http.Request) (limit, offset, page, perPage int) {
	limit, offset, page, perPage = 10, 0, 1, 10
	q := r.URL.Query()
	if pageStr := q.Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if perPageStr := q.Get("per_page"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 {
			perPage = pp
		}
	}
	const maxPerPage = 100
	if perPage > maxPerPage {
		perPage = maxPerPage
	}
	if perPageStr := q.Get("per_page"); perPageStr != "" {
		limit = perPage
		offset = (page - 1) * perPage
	} else {
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
		perPage = limit
		page = 1
		if limit > 0 {
			page = (offset / limit) + 1
		}
	}
	return limit, offset, page, perPage
}

func parseRange(q map[string][]string) (*time.Time, *time.Time, error) {
	fromStr := firstParam(q, "from")
	toStr := firstParam(q, "to")
	now := time.Now().UTC()
	var from, to *time.Time

	if fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			return nil, nil, err
		}
		from = &t
	}
	if toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			return nil, nil, err
		}
		to = &t
	}
	if from == nil {
		f := now.AddDate(0, -1, 0)
		from = &f
	}
	if to == nil {
		to = &now
	}
	return from, to, nil
}

func firstParam(q map[string][]string, key string) string {
	if vals, ok := q[key]; ok && len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// parseScopeQuery parses optional group_id/zone_id/area_id query params.
func parseScopeQuery(q map[string][]string) (fullschedule.ScopeFilter, error) {
	var scope fullschedule.ScopeFilter
	for key, dst := range map[string]**uuid.UUID{
		"group_id": &scope.GroupID,
		"zone_id":  &scope.ZoneID,
		"area_id":  &scope.AreaID,
	} {
		raw := firstParam(q, key)
		if raw == "" {
			continue
		}
		parsed, err := uuid.Parse(raw)
		if err != nil {
			return scope, err
		}
		*dst = &parsed
	}
	return scope, nil
}

func mapCreateRequest(req *presenter.ScheduleCreate, createdBy *uuid.UUID) *fullschedule.Schedule {
	status := fullschedule.StatusValueActive
	if req.Status != nil {
		status = fullschedule.ScheduleStatusValue(*req.Status)
	}
	cronExpr := req.CronExpr
	s := &fullschedule.Schedule{
		Name:        req.Name,
		Mode:        fullschedule.ScheduleMode(req.Mode),
		Status:      status,
		TimeStart:   req.TimeStart,
		EventType:   fullschedule.EventType(req.EventType),
		EventAction: fullschedule.EventAction(req.EventAction),
		Months:      req.Months,
		Dates:       req.Dates,
		CronExpr:    cronExpr,
		GroupID:     req.GroupID,
		ZoneID:      req.ZoneID,
		AreaID:      req.AreaID,
		CreatedBy:   createdBy,
	}
	if req.Event != nil {
		s.Event = *req.Event
	}
	if req.Sunday != nil {
		s.Sunday = *req.Sunday
	}
	if req.Monday != nil {
		s.Monday = *req.Monday
	}
	if req.Tuesday != nil {
		s.Tuesday = *req.Tuesday
	}
	if req.Wednesday != nil {
		s.Wednesday = *req.Wednesday
	}
	if req.Thursday != nil {
		s.Thursday = *req.Thursday
	}
	if req.Friday != nil {
		s.Friday = *req.Friday
	}
	if req.Saturday != nil {
		s.Saturday = *req.Saturday
	}
	return s
}

func mapScheduleResponse(s *fullschedule.Schedule) *presenter.ScheduleResponse {
	res := &presenter.ScheduleResponse{
		ID:            s.ID,
		Name:          s.Name,
		Mode:          string(s.Mode),
		Status:        int8(s.Status),
		TimeStart:     s.TimeStart,
		EventType:     string(s.EventType),
		EventAction:   string(s.EventAction),
		Event:         s.Event,
		Sunday:        s.Sunday,
		Monday:        s.Monday,
		Tuesday:       s.Tuesday,
		Wednesday:     s.Wednesday,
		Thursday:      s.Thursday,
		Friday:        s.Friday,
		Saturday:      s.Saturday,
		Months:        s.Months,
		Dates:         s.Dates,
		CronExpr:      s.CronExpr,
		ManualTrigger: s.ManualTrigger,
		LastRunAt:     s.LastRunAt,
		NextRunAt:     s.NextRunAt,
		RunCount:      s.RunCount,
		SuccessCount:  s.SuccessCount,
		FailedCount:   s.FailedCount,
		DeviceCount:   s.DeviceCount,
		GroupID:       s.GroupID,
		ZoneID:        s.ZoneID,
		AreaID:        s.AreaID,
		CreatedBy:     s.CreatedBy,
		CreatedAt:     helpers.NewLocalTime(s.CreatedAt),
		UpdatedBy:     s.UpdatedBy,
		UpdatedAt:     helpers.NewLocalTime(s.UpdatedAt),
		Version:       s.Version,
	}
	for _, d := range s.Devices {
		item := presenter.DeviceResponse{DeviceID: d.DeviceID}
		if d.DeviceSN != nil {
			item.DeviceSN = *d.DeviceSN
		}
		res.Devices = append(res.Devices, item)
	}
	return res
}

func mapSchedulesResponse(list []*fullschedule.Schedule) []*presenter.ScheduleResponse {
	out := make([]*presenter.ScheduleResponse, 0, len(list))
	for _, s := range list {
		out = append(out, mapScheduleResponse(s))
	}
	return out
}

func mapSettingResponse(s *fullschedule.ScheduleSetting) *presenter.SettingResponse {
	return &presenter.SettingResponse{
		ID:         s.ID,
		ScheduleID: s.ScheduleID,
		Key:        s.Key,
		Value:      string(s.Value),
		UpdatedAt:  helpers.NewLocalTime(s.UpdatedAt),
	}
}

func mapSettingsResponse(settings []*fullschedule.ScheduleSetting) []*presenter.SettingResponse {
	out := make([]*presenter.SettingResponse, 0, len(settings))
	for _, s := range settings {
		out = append(out, mapSettingResponse(s))
	}
	return out
}

func mapHistoryResponse(rows []*fullschedule.ScheduleHistory) []*presenter.HistoryResponse {
	out := make([]*presenter.HistoryResponse, 0, len(rows))
	for _, h := range rows {
		out = append(out, &presenter.HistoryResponse{
			ID:            h.ID,
			ScheduleID:    h.ScheduleID,
			DeviceID:      h.DeviceID,
			TriggeredBy:   string(h.TriggeredBy),
			TriggerSource: string(h.TriggerSource),
			Status:        string(h.Status),
			EventAction:   string(h.EventAction),
			Payload:       h.Payload,
			Message:       h.Message,
			DurationMs:    h.DurationMs,
			RetryCount:    h.RetryCount,
			ExecutedAt:    h.ExecutedAt,
			Timezone:      h.Timezone,
			Date:          h.Date,
			Time:          h.Time,
			CreatedAt:     helpers.NewLocalTime(h.CreatedAt),
		})
	}
	return out
}

func mapReportResponse(rows []*fullschedule.ScheduleReport) []*presenter.ReportResponse {
	out := make([]*presenter.ReportResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, &presenter.ReportResponse{
			ScheduleID:   r.ScheduleID,
			Name:         r.Name,
			Mode:         r.Mode,
			TotalRuns:    r.TotalRuns,
			SuccessCount: r.SuccessCount,
			FailedCount:  r.FailedCount,
			SkippedCount: r.SkippedCount,
			Processing:   r.Processing,
			SuccessRate:  r.SuccessRate,
			AvgDuration:  r.AvgDuration,
			LastRunAt:    r.LastRunAt,
		})
	}
	return out
}

func mapDevicesBrief(devices []*iotmodels.Device) []*presenter.IoTDeviceBrief {
	out := make([]*presenter.IoTDeviceBrief, 0, len(devices))
	for _, d := range devices {
		out = append(out, &presenter.IoTDeviceBrief{
			DeviceID:   d.DeviceID,
			DeviceName: d.DeviceName,
			SN:         d.SN,
			Status:     d.Status,
		})
	}
	return out
}
