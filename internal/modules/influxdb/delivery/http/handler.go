package http

import (
	"fmt"
	"net/http"
	"strconv"

	"icmongolang/internal/modules/influxdb/presenter"
	"icmongolang/internal/modules/influxdb/usecase"
	"icmongolang/pkg/logger"

	"github.com/go-chi/render"
)

// InfluxHandler is the HTTP handler for InfluxDB operations
type InfluxHandler struct {
	uc     usecase.InfluxUseCase
	logger logger.Logger
}

// CreateInfluxHandler creates a new InfluxDB HTTP handler
func CreateInfluxHandler(uc usecase.InfluxUseCase, log logger.Logger) *InfluxHandler {
	return &InfluxHandler{uc: uc, logger: log}
}

// WriteData godoc
// @Summary      Write a data point to InfluxDB
// @Description  Writes a single measurement with fields and tags
// @Tags         influxdb
// @Accept       json
// @Produce      json
// @Param        request body presenter.WriteDataRequest true "Data point"
// @Success      200 {object} map[string]string
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /influx/write [post]
// @Security     BearerAuth
func (h *InfluxHandler) WriteData(w http.ResponseWriter, r *http.Request) {
	var req presenter.WriteDataRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	if err := h.uc.WriteData(r.Context(), &req); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, map[string]string{"status": "ok"})
}

// QueryGetFilter godoc
// @Summary      Query raw time-series data via query parameters
// @Description  Retrieves filtered data points from InfluxDB. Use query parameters for filtering.
// @Tags         influxdb
// @Accept       json
// @Produce      json
// @Param        measurement query string true "Measurement name"
// @Param        field query string true "Field name"
// @Param        start query string false "Start time (e.g. -1h, 2023-01-01T00:00:00Z)"
// @Param        stop query string false "Stop time (default: now())"
// @Param        limit query int false "Limit number of records (default: 1000)"
// @Param        offset query int false "Offset for pagination (default: 0)"
// @Success      200 {array} presenter.DataPointResponse
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /influx/query [get]
// @Security     BearerAuth
func (h *InfluxHandler) QueryGetFilter(w http.ResponseWriter, r *http.Request) {
	// อ่านค่าจาก query parameters
	query := r.URL.Query()

	req := presenter.QueryFilterRequest{
		Measurement: query.Get("measurement"),
		Field:       query.Get("field"),
		Start:       query.Get("start"),
		Stop:        query.Get("stop"),
	}

	// Parse limit (optional, default 0 = use usecase default)
	if limitStr := query.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}

	// Parse offset (optional)
	if offsetStr := query.Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			req.Offset = offset
		}
	}

	// ตรวจสอบค่าที่จำเป็น
	if req.Measurement == "" || req.Field == "" {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("measurement and field are required")))
		return
	}

	// เรียก usecase
	data, err := h.uc.QueryFilter(r.Context(), &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, data)
}

// Querydevicechart godoc
// @Summary      Query raw time-series data
// @Description  Retrieves filtered data points from InfluxDB
// @Tags         influxdb
// @Accept       json
// @Produce      json
// @Param        request body presenter.Querydevicechart true "Filter parameters"
// @Success      200 {array} presenter.DataPointResponse
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /influx/devicechart [post]
// @Security     BearerAuth
func (h *InfluxHandler) Querydevicechart(w http.ResponseWriter, r *http.Request) {
	var req presenter.Querydevicechart
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	data, err := h.uc.Querydevicechart(r.Context(), &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, data) // data is now *presenter.DeviceChartResponse
}

// QueryFilters godoc
// @Summary      Query raw time-series data
// @Description  Retrieves filtered data points from InfluxDB
// @Tags         influxdb
// @Accept       json
// @Produce      json
// @Param        request body presenter.QueryFilterRequest true "Filter parameters"
// @Success      200 {array} presenter.DataPointResponse
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /influx/filters [post]
// @Security     BearerAuth
func (h *InfluxHandler) QueryFilters(w http.ResponseWriter, r *http.Request) {
	var req presenter.QueryFilterRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	data, err := h.uc.QueryFilter(r.Context(), &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, data)
}

// QueryStatistics godoc
// @Summary      Compute statistics over time-series data
// @Description  Aggregates data using mean, median, percentile, etc.
// @Tags         influxdb
// @Accept       json
// @Produce      json
// @Param        request body presenter.StatisticsRequest true "Statistics parameters"
// @Success      200 {object} presenter.StatisticsResponse
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /influx/statistics [post]
// @Security     BearerAuth
func (h *InfluxHandler) QueryStatistics(w http.ResponseWriter, r *http.Request) {
	var req presenter.StatisticsRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	stats, err := h.uc.QueryStatistics(r.Context(), &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, stats)
}

// Error helpers (same as MQTT)
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
