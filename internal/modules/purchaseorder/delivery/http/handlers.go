package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"icmongolang/config"
	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/purchaseorder"
	pomodels "icmongolang/internal/modules/purchaseorder/models"
	"icmongolang/internal/modules/purchaseorder/presenter"
	"icmongolang/pkg/helpers"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"
	"icmongolang/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type purchaseOrderHandler struct {
	cfg    *config.Config
	poUC   purchaseorder.PurchaseOrderUseCaseI
	logger logger.Logger
}

func CreatePurchaseOrderHandler(uc purchaseorder.PurchaseOrderUseCaseI, cfg *config.Config, logger logger.Logger) purchaseorder.Handlers {
	return &purchaseOrderHandler{cfg: cfg, poUC: uc, logger: logger}
}

// Create godoc
// @Summary Create purchase order
// @Description Create a new purchase order with items.
// @Tags purchase-orders
// @Accept json
// @Produce json
// @Param body body presenter.PurchaseOrderCreateRequest true "Purchase Order"
// @Success 200 {object} responses.SuccessResponse[presenter.PurchaseOrderResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 422 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /purchase-orders [post]
func (h *purchaseOrderHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := new(presenter.PurchaseOrderCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
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

		header := mapCreateHeader(req)
		header.UserID = user.ID

		details := mapCreateDetails(req.Items)

		created, err := h.poUC.CreateWithDetails(ctx, header, details)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapPOResponse(created)))
	}
}

// GetByID godoc
// @Summary Get purchase order by ID
// @Description Get purchase order with items by ID.
// @Tags purchase-orders
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Success 200 {object} responses.SuccessResponse[presenter.PurchaseOrderResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /purchase-orders/{id} [get]
func (h *purchaseOrderHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		po, err := h.poUC.GetByIDWithDetails(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapPOResponse(po)))
	}
}

// List godoc
// @Summary List purchase orders
// @Description List purchase orders with pagination and filters.
// @Tags purchase-orders
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Param supplier_id query string false "Supplier ID"
// @Param status query string false "Status filter (DRAFT,SENT,CONFIRMED,SHIPPED,RECEIVED,CANCELLED)"
// @Param start_date query string false "Start date filter"
// @Param end_date query string false "End date filter"
// @Success 200 {object} responses.SuccessResponse[presenter.PaginatedPurchaseOrdersResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /purchase-orders [get]
func (h *purchaseOrderHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		q := r.URL.Query()

		limit := 10
		offset := 0
		page := 1
		perPage := 10

		pageStr := q.Get("page")
		perPageStr := q.Get("per_page")

		if pageStr != "" || perPageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
			if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 {
				perPage = pp
			}
			const maxPerPage = 100
			if perPage > maxPerPage {
				perPage = maxPerPage
			}
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
			if limit > 0 {
				perPage = limit
				page = (offset / limit) + 1
			}
		}

		var supplierID *uuid.UUID
		if sid := q.Get("supplier_id"); sid != "" {
			parsed, err := uuid.Parse(sid)
			if err == nil {
				supplierID = &parsed
			}
		}

		status := q.Get("status")
		startDate := q.Get("start_date")
		endDate := q.Get("end_date")

		items, total, err := h.poUC.List(ctx, limit, offset, supplierID, status, startDate, endDate)
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

		paginatedRes := &presenter.PaginatedPurchaseOrdersResponse{
			Items:      mapPOResponses(items),
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages,
		}
		render.Respond(w, r, responses.CreateSuccessResponse(paginatedRes))
	}
}

// Update godoc
// @Summary Update purchase order
// @Description Update purchase order details (only if DRAFT status).
// @Tags purchase-orders
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Param body body presenter.PurchaseOrderUpdateRequest true "Update data"
// @Success 200 {object} responses.SuccessResponse[presenter.PurchaseOrderResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /purchase-orders/{id} [put]
func (h *purchaseOrderHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		req := new(presenter.PurchaseOrderUpdateRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		values := make(map[string]interface{})
		if req.ExpectedDeliveryDate != nil {
			values["expected_delivery_date"] = *req.ExpectedDeliveryDate
		}
		if req.Currency != "" {
			values["currency"] = req.Currency
		}
		if req.ExchangeRate != nil {
			values["exchange_rate"] = *req.ExchangeRate
		}
		if req.ShippingCost != nil {
			values["shipping_cost"] = *req.ShippingCost
		}
		if req.PaymentTerms != "" {
			values["payment_terms"] = req.PaymentTerms
		}
		if req.DeliveryAddress != "" {
			values["delivery_address"] = req.DeliveryAddress
		}
		if req.Notes != "" {
			values["notes"] = req.Notes
		}
		if req.TermsAndConditions != "" {
			values["terms_and_conditions"] = req.TermsAndConditions
		}
		if req.TaxRate != nil {
			values["tax_rate"] = *req.TaxRate
		}
		if req.DiscountType != "" {
			values["discount_type"] = req.DiscountType
		}
		if req.DiscountValue != nil {
			values["discount_value"] = *req.DiscountValue
		}

		updated, err := h.poUC.Update(ctx, id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapPOResponse(updated)))
	}
}

// Delete godoc
// @Summary Delete purchase order
// @Description Soft delete a purchase order (only if DRAFT status).
// @Tags purchase-orders
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Success 200 {object} responses.SuccessResponse[presenter.PurchaseOrderResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /purchase-orders/{id} [delete]
func (h *purchaseOrderHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		deleted, err := h.poUC.Delete(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapPOResponse(deleted)))
	}
}

// Send godoc
// @Summary Send purchase order to supplier
// @Description Send purchase order to supplier via email with PDF.
// @Tags purchase-orders
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Success 200 {object} responses.SuccessResponse[presenter.PurchaseOrderResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /purchase-orders/{id}/send [post]
func (h *purchaseOrderHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		po, err := h.poUC.Send(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapPOResponse(po)))
	}
}

// Confirm godoc
// @Summary Confirm purchase order
// @Description Confirm purchase order when supplier acknowledges the order.
// @Tags purchase-orders
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Success 200 {object} responses.SuccessResponse[presenter.PurchaseOrderResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /purchase-orders/{id}/confirm [put]
func (h *purchaseOrderHandler) Confirm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		po, err := h.poUC.Confirm(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapPOResponse(po)))
	}
}

// Receive godoc
// @Summary Receive goods
// @Description Receive goods into inventory and update PO.
// @Tags purchase-orders
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Param body body presenter.PurchaseOrderReceiveRequest true "Receive items"
// @Success 200 {object} responses.SuccessResponse[presenter.PurchaseOrderResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /purchase-orders/{id}/receive [post]
func (h *purchaseOrderHandler) Receive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		req := new(presenter.PurchaseOrderReceiveRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		if err := utils.ValidateStruct(ctx, req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		receiveReq := &purchaseorder.ReceiveRequest{
			Items: make([]purchaseorder.ReceiveItem, len(req.Items)),
		}
		for i, item := range req.Items {
			receiveReq.Items[i] = purchaseorder.ReceiveItem{
				DetailID:         item.DetailID,
				ReceivedQuantity: item.ReceivedQuantity,
			}
		}

		po, err := h.poUC.Receive(ctx, id, receiveReq)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapPOResponse(po)))
	}
}

// Cancel godoc
// @Summary Cancel purchase order
// @Description Cancel a purchase order with reason.
// @Tags purchase-orders
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Param reason body presenter.PurchaseOrderStatusRequest true "Cancel reason"
// @Success 200 {object} responses.SuccessResponse[presenter.PurchaseOrderResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /purchase-orders/{id}/cancel [put]
func (h *purchaseOrderHandler) Cancel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		req := new(presenter.PurchaseOrderStatusRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		po, err := h.poUC.Cancel(ctx, id, req.Reason, user.ID)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapPOResponse(po)))
	}
}

// GetPDF godoc
// @Summary Generate purchase order PDF
// @Description Generate PDF document for purchase order.
// @Tags purchase-orders
// @Produce application/pdf
// @Param id path string true "Purchase Order ID"
// @Success 200 {file} byte "PDF file"
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /purchase-orders/{id}/pdf [get]
func (h *purchaseOrderHandler) GetPDF() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		pdf, err := h.poUC.GetPDF(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "inline; filename=po_"+id.String()+".pdf")
		w.Write(pdf)
	}
}

// GetSuggestions godoc
// @Summary Get PO suggestions from job
// @Description Suggest parts to order from a job based on quotation and stock.
// @Tags purchase-orders
// @Accept json
// @Produce json
// @Param jobId path string true "Job ID"
// @Success 200 {object} responses.SuccessResponse[[]presenter.PurchaseOrderSuggestionDTO]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /purchase-orders/suggestions/{jobId} [get]
func (h *purchaseOrderHandler) GetSuggestions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		jobID, err := uuid.Parse(chi.URLParam(r, "jobId"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		suggestions, err := h.poUC.GetSuggestions(ctx, jobID)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		dtos := make([]presenter.PurchaseOrderSuggestionDTO, len(suggestions))
		for i, s := range suggestions {
			dtos[i] = presenter.PurchaseOrderSuggestionDTO{
				PartID:        s.PartID,
				PartName:      s.PartName,
				PartCode:      s.PartCode,
				SuggestedQty:  s.SuggestedQty,
				CurrentStock:  s.CurrentStock,
				UnitPrice:     s.UnitPrice,
				FromQuotation: s.FromQuotation,
			}
		}

		render.Respond(w, r, responses.CreateSuccessResponse(dtos))
	}
}

// GetStatusHistory godoc
// @Summary Get purchase order status history
// @Description Get status change history of a purchase order.
// @Tags purchase-orders
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Success 200 {object} responses.SuccessResponse[[]presenter.PurchaseOrderStatusHistoryResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /purchase-orders/{id}/history [get]
func (h *purchaseOrderHandler) GetStatusHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		history, err := h.poUC.GetStatusHistory(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		dtos := make([]presenter.PurchaseOrderStatusHistoryResponse, len(history))
		for i, h := range history {
			var fromStatus string
			if h.FromStatus != nil {
				fromStatus = *h.FromStatus
			}
			var reason string
			if h.Reason != nil {
				reason = *h.Reason
			}
			dtos[i] = presenter.PurchaseOrderStatusHistoryResponse{
				ID:         h.ID,
				PoHeaderID: h.PoHeaderID,
				FromStatus: fromStatus,
				ToStatus:   h.ToStatus,
				ChangedBy:  h.ChangedBy,
				ChangedAt:  h.ChangedAt,
				Reason:     reason,
			}
		}

		render.Respond(w, r, responses.CreateSuccessResponse(dtos))
	}
}

// CreateFromQuotation godoc
// @Summary Create purchase order from quotation
// @Description Create a purchase order from an approved quotation.
// @Tags purchase-orders
// @Accept json
// @Produce json
// @Param quotationId path string true "Quotation ID"
// @Success 200 {object} responses.SuccessResponse[presenter.PurchaseOrderResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /purchase-orders/from-quotation/{quotationId} [post]
func (h *purchaseOrderHandler) CreateFromQuotation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		quotationID, err := uuid.Parse(chi.URLParam(r, "quotationId"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		po, err := h.poUC.CreateFromQuotation(ctx, quotationID, user.ID, uuid.Nil)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapPOResponse(po)))
	}
}

func mapCreateHeader(req *presenter.PurchaseOrderCreateRequest) *pomodels.PurchaseOrderHeader {
	header := &pomodels.PurchaseOrderHeader{
		SupplierID:    req.SupplierID,
		QuotationID:   req.QuotationID,
		JobID:         req.JobID,
		Currency:      req.Currency,
		ExchangeRate:  req.ExchangeRate,
		ShippingCost:  req.ShippingCost,
		TaxRate:       req.TaxRate,
		DiscountType:  stringPtr(req.DiscountType),
		DiscountValue: req.DiscountValue,
		Status:        "DRAFT",
		Subtotal:      0,
		Total:         0,
	}
	if req.PaymentTerms != "" {
		header.PaymentTerms = &req.PaymentTerms
	}
	if req.DeliveryAddress != "" {
		header.DeliveryAddress = &req.DeliveryAddress
	}
	if req.Notes != "" {
		header.Notes = &req.Notes
	}
	if req.TermsAndConditions != "" {
		header.TermsAndConditions = &req.TermsAndConditions
	}
	if req.ExpectedDeliveryDate != nil {
		t, err := timeParse(*req.ExpectedDeliveryDate)
		if err == nil {
			header.ExpectedDeliveryDate = t
		}
	}
	return header
}

func mapCreateDetails(items []presenter.PurchaseOrderDetailRequest) []*pomodels.PurchaseOrderDetail {
	details := make([]*pomodels.PurchaseOrderDetail, len(items))
	for i, item := range items {
		netPrice := (float64(item.QuantityOrdered) * item.UnitPrice) - item.Discount
		details[i] = &pomodels.PurchaseOrderDetail{
			PartID:          item.PartID,
			QuantityOrdered: item.QuantityOrdered,
			UnitPrice:       item.UnitPrice,
			TotalPrice:      float64(item.QuantityOrdered) * item.UnitPrice,
			Discount:        item.Discount,
			NetPrice:        netPrice,
		}
		if item.Note != "" {
			details[i].Note = &item.Note
		}
	}
	return details
}

func mapPOResponse(po *pomodels.PurchaseOrderHeader) *presenter.PurchaseOrderResponse {
	resp := &presenter.PurchaseOrderResponse{
		ID:                   po.ID,
		PONo:                 po.PONo,
		QuotationID:          po.QuotationID,
		JobID:                po.JobID,
		SupplierID:           po.SupplierID,
		PODate:               po.PODate,
		ExpectedDeliveryDate: po.ExpectedDeliveryDate,
		ActualDeliveryDate:   po.ActualDeliveryDate,
		Status:               po.Status,
		Subtotal:             po.Subtotal,
		TaxRate:              po.TaxRate,
		TaxAmount:            po.TaxAmount,
		DiscountValue:        po.DiscountValue,
		Total:                po.Total,
		Currency:             po.Currency,
		ExchangeRate:         po.ExchangeRate,
		ShippingCost:         po.ShippingCost,
		SentAt:               po.SentAt,
		ConfirmedAt:          po.ConfirmedAt,
		ReceivedBy:           po.ReceivedBy,
		CreatedAt:            helpers.NewLocalTime(po.CreatedAt),
		UpdatedAt:            helpers.NewLocalTimePtr(po.UpdatedAt),
		UserID:               po.UserID,
		WhitelabelID:         po.WhitelabelID,
	}
	if po.DiscountType != nil {
		resp.DiscountType = *po.DiscountType
	}
	if po.PaymentTerms != nil {
		resp.PaymentTerms = *po.PaymentTerms
	}
	if po.DeliveryAddress != nil {
		resp.DeliveryAddress = *po.DeliveryAddress
	}
	if po.Notes != nil {
		resp.Notes = *po.Notes
	}
	if po.TermsAndConditions != nil {
		resp.TermsAndConditions = *po.TermsAndConditions
	}
	if len(po.Details) > 0 {
		resp.Items = make([]presenter.PurchaseOrderDetailResponse, len(po.Details))
		for i, d := range po.Details {
			resp.Items[i] = presenter.PurchaseOrderDetailResponse{
				ID:               d.ID,
				PartID:           d.PartID,
				QuantityOrdered:  d.QuantityOrdered,
				QuantityReceived: d.QuantityReceived,
				UnitPrice:        d.UnitPrice,
				TotalPrice:       d.TotalPrice,
				Discount:         d.Discount,
				NetPrice:         d.NetPrice,
			}
			if d.Note != nil {
				resp.Items[i].Note = *d.Note
			}
		}
	}
	return resp
}

func mapPOResponses(pos []*pomodels.PurchaseOrderHeader) []*presenter.PurchaseOrderResponse {
	out := make([]*presenter.PurchaseOrderResponse, len(pos))
	for i, po := range pos {
		out[i] = mapPOResponse(po)
	}
	return out
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func timeParse(s string) (*time.Time, error) {
	// โปรดระวัง รูปแบบเวลาที่ใช้
	// Parse time in RFC3339 format
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
