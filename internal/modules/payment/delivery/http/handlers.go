package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"icmongolang/config"
	"icmongolang/internal/middleware"
	"icmongolang/internal/models"
	"icmongolang/internal/modules/payment"
	"icmongolang/internal/modules/payment/presenter"
	"icmongolang/pkg/helpers"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/report"
	"icmongolang/pkg/responses"
	"icmongolang/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type paymentHandler struct {
	cfg       *config.Config
	paymentUC payment.PaymentUseCaseI
	logger    logger.Logger
}

func CreatePaymentHandler(uc payment.PaymentUseCaseI, cfg *config.Config, logger logger.Logger) payment.Handlers {
	return &paymentHandler{cfg: cfg, paymentUC: uc, logger: logger}
}

// Create godoc
// @Summary Record payment
// @Description Record a payment from a customer, link to Invoice, and auto-generate a receipt.
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body presenter.PaymentRecordRequest true "Payment record"
// @Success 200 {object} responses.SuccessResponse[presenter.PaymentResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /payments [post]
func (h *paymentHandler) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := new(presenter.PaymentRecordRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		if err := utils.ValidateStruct(ctx, req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		// ✅ สร้าง model Payment
		paymentModel := mapPaymentRequest(req)
		paymentModel.ReceivedBy = user.ID
		paymentModel.UserID = user.ID

		createdPayment, err := h.paymentUC.Create(ctx, paymentModel)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapPaymentResponse(createdPayment)))
	}
}

// Get godoc
// @Summary Get payment by ID
// @Description Retrieve payment details by ID.
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} responses.SuccessResponse[presenter.PaymentResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /payments/{id} [get]
func (h *paymentHandler) Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		payment, err := h.paymentUC.Get(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapPaymentResponse(payment)))
	}
}

// List godoc
// @Summary List payments
// @Description List payments with pagination.
// @Tags payments
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} responses.SuccessResponse[presenter.PaginatedPaymentResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /payments [get]
func (h *paymentHandler) List() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		page := 1
		perPage := 10

		if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
			page = p
		}
		if pp, err := strconv.Atoi(r.URL.Query().Get("per_page")); err == nil && pp > 0 {
			perPage = pp
		}
		const maxPerPage = 100
		if perPage > maxPerPage {
			perPage = maxPerPage
		}

		limit := perPage
		offset := (page - 1) * perPage

		payments, err := h.paymentUC.GetMulti(ctx, limit, offset)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		total, err := h.paymentUC.Count(ctx)
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

		paginatedRes := &presenter.PaginatedPaymentResponse{
			Payments:   mapPaymentsResponse(payments),
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages,
		}
		render.Respond(w, r, responses.CreateSuccessResponse(paginatedRes))
	}
}

// GetByInvoice godoc
// @Summary Get payment by Invoice ID
// @Description Retrieve payment details by Invoice ID.
// @Tags payments
// @Accept json
// @Produce json
// @Param invoiceId path string true "Invoice ID"
// @Success 200 {object} responses.SuccessResponse[presenter.PaymentResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /payments/invoice/{invoiceId} [get]
func (h *paymentHandler) GetByInvoice() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		invoiceId, err := uuid.Parse(chi.URLParam(r, "invoiceId"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		payment, err := h.paymentUC.GetByInvoiceId(ctx, invoiceId)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapPaymentResponse(payment)))
	}
}

// Search godoc
// @Summary Search payments
// @Description Search payments with filters: customer, date range, status.
// @Tags payments
// @Accept json
// @Produce json
// @Param search body presenter.PaymentSearchRequest true "Search filters"
// @Success 200 {object} responses.SuccessResponse[presenter.PaginatedPaymentResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /payments/search [post]
func (h *paymentHandler) Search() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := new(presenter.PaymentSearchRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		filter := make(map[string]interface{})
		if req.CustomerID != nil {
			filter["customer_id"] = req.CustomerID.String()
		}
		if req.InvoiceID != nil {
			filter["invoice_id"] = req.InvoiceID.String()
		}
		if req.Status != nil {
			filter["status"] = *req.Status
		}
		if req.PaymentMethodID != nil {
			filter["payment_method_id"] = req.PaymentMethodID.String()
		}
		if req.DateFrom != nil {
			filter["date_from"] = *req.DateFrom
		}
		if req.DateTo != nil {
			filter["date_to"] = *req.DateTo
		}

		page := req.Page
		if page < 1 {
			page = 1
		}
		perPage := req.PerPage
		if perPage < 1 {
			perPage = 10
		}
		const maxPerPage = 100
		if perPage > maxPerPage {
			perPage = maxPerPage
		}
		limit := perPage
		offset := (page - 1) * perPage

		payments, err := h.paymentUC.Search(ctx, filter, limit, offset)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		total, err := h.paymentUC.Count(ctx)
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

		paginatedRes := &presenter.PaginatedPaymentResponse{
			Payments:   mapPaymentsResponse(payments),
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages,
		}
		render.Respond(w, r, responses.CreateSuccessResponse(paginatedRes))
	}
}

// GetOutstanding godoc
// @Summary Get customer outstanding balance
// @Description Show outstanding balance for a customer (unpaid invoices).
// @Tags payments
// @Accept json
// @Produce json
// @Param customerId path string true "Customer ID"
// @Success 200 {object} responses.SuccessResponse[[]presenter.OutstandingBalanceResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /payments/outstanding/{customerId} [get]
func (h *paymentHandler) GetOutstanding() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		customerId, err := uuid.Parse(chi.URLParam(r, "customerId"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		balances, err := h.paymentUC.GetOutstandingByCustomerId(ctx, customerId)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapOutstandingResponse(balances)))
	}
}

// GetHistory godoc
// @Summary Get customer payment history
// @Description Show payment history of a customer.
// @Tags payments
// @Accept json
// @Produce json
// @Param customerId path string true "Customer ID"
// @Success 200 {object} responses.SuccessResponse[[]presenter.PaymentHistoryResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /payments/history/{customerId} [get]
func (h *paymentHandler) GetHistory() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		customerId, err := uuid.Parse(chi.URLParam(r, "customerId"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		history, err := h.paymentUC.GetPaymentHistory(ctx, customerId)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapHistoryResponse(history)))
	}
}

// Refund godoc
// @Summary Process refund
// @Description Process a refund to a customer for a payment.
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Param refund body presenter.RefundRequest true "Refund request"
// @Success 200 {object} responses.SuccessResponse[presenter.PaymentResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /payments/{id}/refund [post]
func (h *paymentHandler) Refund() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		req := new(presenter.RefundRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		if err := utils.ValidateStruct(ctx, req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		updatedPayment, err := h.paymentUC.ProcessRefund(ctx, id, req.Amount, req.Reason)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapPaymentResponse(updatedPayment)))
	}
}

// Cancel godoc
// @Summary Cancel payment
// @Description Cancel a payment (only if not yet confirmed).
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Param reason query string true "Cancellation reason"
// @Success 200 {object} responses.SuccessResponse[presenter.PaymentResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /payments/{id}/cancel [put]
func (h *paymentHandler) Cancel() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		reason := r.URL.Query().Get("reason")
		if reason == "" {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(nil)))
			return
		}

		err = h.paymentUC.CancelPayment(ctx, id, reason)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse[interface{}](nil))
	}
}

// GetReceipt godoc
// @Summary Get receipt by ID
// @Description Retrieve receipt details by ID.
// @Tags receipts
// @Accept json
// @Produce json
// @Param id path string true "Receipt ID"
// @Success 200 {object} responses.SuccessResponse[presenter.ReceiptResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /receipts/{id} [get]
func (h *paymentHandler) GetReceipt() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		receipt, err := h.paymentUC.GetReceipt(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapReceiptResponse(receipt)))
	}
}

// GetReceiptByPayment godoc
// @Summary Get receipt by Payment ID
// @Description Retrieve receipt by Payment ID.
// @Tags receipts
// @Accept json
// @Produce json
// @Param paymentId path string true "Payment ID"
// @Success 200 {object} responses.SuccessResponse[presenter.ReceiptResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /receipts/payment/{paymentId} [get]
func (h *paymentHandler) GetReceiptByPayment() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		paymentId, err := uuid.Parse(chi.URLParam(r, "paymentId"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		receipt, err := h.paymentUC.GetReceiptByPaymentId(ctx, paymentId)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapReceiptResponse(receipt)))
	}
}

// GetReceiptPDF godoc
// @Summary Generate receipt PDF
// @Description Generate a PDF file of the receipt.
// @Tags receipts
// @Accept json
// @Produce application/pdf
// @Param id path string true "Receipt ID"
// @Success 200 {file} byte[]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /receipts/{id}/pdf [get]
func (h *paymentHandler) GetReceiptPDF() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		receipt, err := h.paymentUC.GetReceipt(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		data := report.ReceiptData{
			Company: report.CompanyInfo{
				Name:    "ICMON Auto Repair",
				Address: "123/4 ถนนสุขุมวิท แขวงคลองเตย เขตคลองเตย กรุงเทพฯ 10110",
				Phone:   "02-123-4567",
				TaxID:   "0123456789012",
			},
			ReceiptNo:     receipt.ReceiptNo,
			Date:          receipt.ReceiptDate,
			CustomerName:  receipt.CustomerID.String(),
			PaymentMethod: receipt.ReceiptType,
			Items: []report.ReceiptItem{
				{Description: "ค่าบริการซ่อมบำรุง", Total: receipt.Amount},
			},
			Amount:      receipt.Amount,
			AmountWords: report.BahtThai(receipt.Amount),
		}

		pdf, err := report.GeneratePDF(ctx, report.TplReceipt, data)
		if err != nil {
			h.logger.Errorf("Failed to generate receipt PDF: %v", err)
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrInternalServer(err)))
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "inline; filename=receipt_"+id.String()+".pdf")
		w.Write(pdf)
	}
}

// CancelReceipt godoc
// @Summary Cancel receipt
// @Description Cancel a receipt.
// @Tags receipts
// @Accept json
// @Produce json
// @Param id path string true "Receipt ID"
// @Param reason query string true "Cancellation reason"
// @Success 200 {object} responses.SwaggerSuccessResponse
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /receipts/{id}/cancel [put]
func (h *paymentHandler) CancelReceipt() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		reason := r.URL.Query().Get("reason")
		if reason == "" {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(nil)))
			return
		}

		err = h.paymentUC.CancelReceipt(ctx, id, reason)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse[interface{}](nil))
	}
}

// ============================================================================
// Mapper functions (ฟังก์ชันแปลงข้อมูล)
// ============================================================================

func mapPaymentRequest(req *presenter.PaymentRecordRequest) *models.Payment {
	return &models.Payment{
		InvoiceID:       req.InvoiceID,
		PaymentMethodID: req.PaymentMethodID,
		Amount:          req.Amount,
		AmountReceived:  req.AmountReceived,
		ChangeAmount:    req.ChangeAmount,
		Currency:        req.Currency,
		ExchangeRate:    req.ExchangeRate,
		ReferenceNumber: req.ReferenceNumber,
		BankName:        req.BankName,
		ChequeNumber:    req.ChequeNumber,
		ChequeBank:      req.ChequeBank,
		Notes:           req.Notes,
		Status:          "PENDING",
	}
}

func mapPaymentResponse(p *models.Payment) *presenter.PaymentResponse {
	if p == nil {
		return nil
	}
	resp := &presenter.PaymentResponse{
		ID:              p.ID,
		PaymentNo:       p.PaymentNo,
		InvoiceID:       p.InvoiceID,
		CustomerID:      p.CustomerID,
		PaymentDate:     p.PaymentDate,
		PaymentMethodID: p.PaymentMethodID,
		Amount:          p.Amount,
		AmountReceived:  p.AmountReceived,
		ChangeAmount:    p.ChangeAmount,
		Currency:        p.Currency,
		ExchangeRate:    p.ExchangeRate,
		Status:          p.Status,
		ReferenceNumber: p.ReferenceNumber,
		BankName:        p.BankName,
		ChequeNumber:    p.ChequeNumber,
		ChequeBank:      p.ChequeBank,
		Notes:           p.Notes,
		ReceivedBy:      p.ReceivedBy,
		ApprovedBy:      p.ApprovedBy,
		ApprovedAt:      p.ApprovedAt,
		RefundedAmount:  p.RefundedAmount,
		RefundedAt:      p.RefundedAt,
		CreatedAt:       helpers.NewLocalTime(p.CreatedAt),
	}
	if p.JobID != nil {
		resp.JobID = p.JobID
	}
	if p.ChequeDate != nil {
		dateStr := p.ChequeDate.Format("2006-01-02")
		resp.ChequeDate = &dateStr
	}
	return resp
}

func mapPaymentsResponse(payments []*models.Payment) []*presenter.PaymentResponse {
	out := make([]*presenter.PaymentResponse, len(payments))
	for i, p := range payments {
		out[i] = mapPaymentResponse(p)
	}
	return out
}

func mapReceiptResponse(r *models.Receipt) *presenter.ReceiptResponse {
	if r == nil {
		return nil
	}
	return &presenter.ReceiptResponse{
		ID:              r.ID,
		ReceiptNo:       r.ReceiptNo,
		PaymentID:       r.PaymentID,
		InvoiceID:       r.InvoiceID,
		CustomerID:      r.CustomerID,
		ReceiptDate:     r.ReceiptDate,
		ReceiptType:     r.ReceiptType,
		Amount:          r.Amount,
		AmountInWordsTh: r.AmountInWordsTh,
		AmountInWordsEn: r.AmountInWordsEn,
		Currency:        r.Currency,
		Status:          r.Status,
		Notes:           r.Notes,
		IssuedBy:        r.IssuedBy,
		CreatedAt:       helpers.NewLocalTime(r.CreatedAt),
	}
}

func mapOutstandingResponse(balances []*models.OutstandingBalance) []*presenter.OutstandingBalanceResponse {
	out := make([]*presenter.OutstandingBalanceResponse, len(balances))
	for i, b := range balances {
		out[i] = &presenter.OutstandingBalanceResponse{
			InvoiceID:         b.InvoiceID,
			InvoiceTotal:      b.InvoiceTotal,
			AmountPaid:        b.AmountPaid,
			OutstandingAmount: b.OutstandingAmount,
			LastPaymentDate:   b.LastPaymentDate,
			Status:            b.Status,
		}
	}
	return out
}

func mapHistoryResponse(history []*models.PaymentHistory) []*presenter.PaymentHistoryResponse {
	out := make([]*presenter.PaymentHistoryResponse, len(history))
	for i, h := range history {
		out[i] = &presenter.PaymentHistoryResponse{
			ID:         h.ID,
			PaymentID:  h.PaymentID,
			FromStatus: h.FromStatus,
			ToStatus:   h.ToStatus,
			ChangedBy:  h.ChangedBy,
			ChangedAt:  h.ChangedAt,
			Reason:     h.Reason,
		}
	}
	return out
}
