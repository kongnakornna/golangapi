package http

import (
	"encoding/json"
	"net/http"

	"errors"

	"icmongolang/config"
	"icmongolang/internal/modules/wos"
	"icmongolang/internal/modules/wos/presenter"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"
	"icmongolang/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type wosHandler struct {
	cfg    *config.Config
	wosUC  wos.WosUseCaseI
	logger logger.Logger
}

// CreateWosHandler creates a new WOS HTTP handler.
// สร้างตัวจัดการ HTTP สำหรับระบบสั่งซื้อออนไลน์
func CreateWosHandler(uc wos.WosUseCaseI, cfg *config.Config, logger logger.Logger) wos.Handlers {
	return &wosHandler{cfg: cfg, wosUC: uc, logger: logger}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new customer order in the Web Order System.
// @Tags wos
// @Accept json
// @Produce json
// @Param order body presenter.WosOrderRequest true "Order details"
// @Success 200 {object} responses.SuccessResponse[presenter.WosOrderResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/wos/orders [post]
func (h *wosHandler) CreateOrder() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.WosOrderRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(req))
	}
}

// GetOrder godoc
// @Summary Get an order
// @Description Get an order by ID.
// @Tags wos
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} responses.SuccessResponse[presenter.WosOrderResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/wos/orders/{id} [get]
func (h *wosHandler) GetOrder() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		order, err := h.wosUC.Get(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(order))
	}
}

// ListOrders godoc
// @Summary List orders
// @Description List all orders with pagination.
// @Tags wos
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} responses.SuccessResponse[[]presenter.WosOrderResponse]
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/wos/orders [get]
func (h *wosHandler) ListOrders() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		orders, err := h.wosUC.GetMulti(r.Context(), 50, 0)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(orders))
	}
}

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update the status of an existing order.
// @Tags wos
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param status query string true "New status (pending, confirmed, shipped, delivered, cancelled)"
// @Success 200 {object} responses.SuccessResponse[string]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/wos/orders/{id}/status [put]
func (h *wosHandler) UpdateOrderStatus() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		status := r.URL.Query().Get("status")
		if status == "" {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrBadRequest(errors.New("status is required"))))
			return
		}
		if err := h.wosUC.UpdateStatus(r.Context(), id, status); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse("Status updated"))
	}
}
