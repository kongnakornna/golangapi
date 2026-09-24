package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"icmongolang/config"
	"icmongolang/internal/middleware"
	"icmongolang/internal/models"
	"icmongolang/internal/modules/quotation"
	"icmongolang/internal/modules/quotation/presenter"
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

type quotationHandler struct {
	cfg         *config.Config
	quotationUC quotation.QuotationUseCaseI
	logger      logger.Logger
}

func CreateQuotationHandler(uc quotation.QuotationUseCaseI, cfg *config.Config, logger logger.Logger) quotation.Handlers {
	return &quotationHandler{cfg: cfg, quotationUC: uc, logger: logger}
}

// Create godoc
// @Summary Create Quotation
// @Description Create new quotation.
// @Tags quotation
// @Accept json
// @Produce json
// @Param quotation body presenter.QuotationCreate true "Add quotation"
// @Success 200 {object} responses.SuccessResponse[presenter.QuotationResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /quotation [post]
func (h *quotationHandler) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := new(presenter.QuotationCreate)

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		err = utils.ValidateStruct(ctx, req)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		newQuotation := mapCreateToModel(req, user.ID, uuid.Nil)
		createdQuotation, err := h.quotationUC.Create(ctx, newQuotation)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelToResponse(createdQuotation)))
	}
}

// Get godoc
// @Summary Read quotation
// @Description Get quotation by ID.
// @Tags quotation
// @Accept json
// @Produce json
// @Param id path string true "Quotation Id"
// @Success 200 {object} responses.SuccessResponse[presenter.QuotationResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /quotation/{id} [get]
func (h *quotationHandler) Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		quotation, err := h.quotationUC.Get(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelToResponse(quotation)))
	}
}

// GetMulti godoc
// @Summary Read quotations
// @Description Retrieve quotations with pagination.
// @Tags quotation
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)" Format(page)
// @Param per_page query int false "Items per page (default 10)" Format(per_page)
// @Success 200 {object} responses.SuccessResponse[presenter.PaginatedQuotationResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /quotation [get]
func (h *quotationHandler) GetMulti() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		q := r.URL.Query()

		limit := 10
		offset := 0
		page := 1
		perPage := 10

		pageStr := q.Get("page")
		if pageStr == "" {
			pageStr = q.Get("offset")
		}
		perPageStr := q.Get("per_page")
		if perPageStr == "" {
			perPageStr = q.Get("limit")
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

		items, err := h.quotationUC.GetMulti(ctx, limit, offset)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		total, err := h.quotationUC.Count(ctx)
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

		paginatedRes := &presenter.PaginatedQuotationResponse{
			Items:      mapModelsToResponse(items),
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages,
		}
		render.Respond(w, r, responses.CreateSuccessResponse(paginatedRes))
	}
}

// Delete godoc
// @Summary Delete quotation
// @Description Delete a quotation by ID.
// @Tags quotation
// @Accept json
// @Produce json
// @Param id path string true "Quotation Id"
// @Success 200 {object} responses.SuccessResponse[presenter.QuotationResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /quotation/{id} [delete]
func (h *quotationHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		item, err := h.quotationUC.Delete(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelToResponse(item)))
	}
}

// Update godoc
// @Summary Update quotation
// @Description Update a quotation by ID.
// @Tags quotation
// @Accept json
// @Produce json
// @Param id path string true "Quotation Id"
// @Param quotation body presenter.QuotationUpdate true "Update quotation"
// @Success 200 {object} responses.SuccessResponse[presenter.QuotationResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /quotation/{id} [put]
func (h *quotationHandler) Update() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		req := new(presenter.QuotationUpdate)

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		err = utils.ValidateStruct(ctx, req)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		values := make(map[string]interface{})
		if req.QuotationDate != nil {
			values["quotation_date"] = *req.QuotationDate
		}
		if req.ExpiryDate != nil {
			values["expiry_date"] = *req.ExpiryDate
		}
		if req.Subtotal != nil {
			values["subtotal"] = *req.Subtotal
		}
		if req.TaxRate != nil {
			values["tax_rate"] = *req.TaxRate
		}
		if req.TaxAmount != nil {
			values["tax_amount"] = *req.TaxAmount
		}
		if req.DiscountType != nil {
			values["discount_type"] = *req.DiscountType
		}
		if req.DiscountValue != nil {
			values["discount_value"] = *req.DiscountValue
		}
		if req.Total != nil {
			values["total"] = *req.Total
		}
		if req.AmountInWordsTh != nil {
			values["amount_in_words_th"] = *req.AmountInWordsTh
		}
		if req.AmountInWordsEn != nil {
			values["amount_in_words_en"] = *req.AmountInWordsEn
		}
		if req.Currency != nil {
			values["currency"] = *req.Currency
		}
		if req.ExchangeRate != nil {
			values["exchange_rate"] = *req.ExchangeRate
		}
		if req.Notes != nil {
			values["notes"] = *req.Notes
		}
		if req.TermsAndConditions != nil {
			values["terms_and_conditions"] = *req.TermsAndConditions
		}

		updatedItem, err := h.quotationUC.Update(ctx, id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelToResponse(updatedItem)))
	}
}

func mapCreateToModel(req *presenter.QuotationCreate, userID uuid.UUID, whitelabelID uuid.UUID) *models.Quotation {
	taxAmount := req.Subtotal * (req.TaxRate / 100.0)
	discount := 0.0
	if req.DiscountType != nil && *req.DiscountType == "PERCENTAGE" {
		discount = req.Subtotal * (req.DiscountValue / 100.0)
	} else {
		discount = req.DiscountValue
	}
	total := req.Subtotal - discount + taxAmount

	m := &models.Quotation{
		QuotationNo:        req.QuotationNo,
		JobID:              req.JobID,
		CustomerID:         req.CustomerID,
		Subtotal:           req.Subtotal,
		TaxRate:            req.TaxRate,
		TaxAmount:          taxAmount,
		DiscountType:       req.DiscountType,
		DiscountValue:      req.DiscountValue,
		Total:              total,
		AmountInWordsTh:    req.AmountInWordsTh,
		AmountInWordsEn:    req.AmountInWordsEn,
		Currency:           req.Currency,
		ExchangeRate:       req.ExchangeRate,
		Notes:              req.Notes,
		TermsAndConditions: req.TermsAndConditions,
		UserID:             userID,
		WhitelabelID:       whitelabelID,
	}
	if req.QuotationDate != "" {
		if t, err := time.Parse(time.RFC3339, req.QuotationDate); err == nil {
			m.QuotationDate = t
		}
	}
	if req.ExpiryDate != "" {
		if t, err := time.Parse(time.RFC3339, req.ExpiryDate); err == nil {
			m.ExpiryDate = t
		}
	}
	return m
}

func mapModelToResponse(m *models.Quotation) *presenter.QuotationResponse {
	if m == nil {
		return nil
	}
	return &presenter.QuotationResponse{
		ID:                 m.ID,
		QuotationNo:        m.QuotationNo,
		JobID:              m.JobID,
		CustomerID:         m.CustomerID,
		QuotationDate:      m.QuotationDate,
		ExpiryDate:         m.ExpiryDate,
		Status:             m.Status,
		Subtotal:           m.Subtotal,
		TaxRate:            m.TaxRate,
		TaxAmount:          m.TaxAmount,
		DiscountType:       m.DiscountType,
		DiscountValue:      m.DiscountValue,
		Total:              m.Total,
		AmountInWordsTh:    m.AmountInWordsTh,
		AmountInWordsEn:    m.AmountInWordsEn,
		Currency:           m.Currency,
		ExchangeRate:       m.ExchangeRate,
		Notes:              m.Notes,
		TermsAndConditions: m.TermsAndConditions,
		ApprovedBy:         m.ApprovedBy,
		ApprovedAt:         m.ApprovedAt,
		RejectedReason:     m.RejectedReason,
		ConvertedToPo:      m.ConvertedToPo,
		UserID:             m.UserID,
		WhitelabelID:       m.WhitelabelID,
		CreatedAt:          helpers.NewLocalTime(m.CreatedAt),
		UpdatedAt:          helpers.NewLocalTimePtr(m.UpdatedAt),
	}
}

// GetPDF godoc
// @Summary Generate quotation PDF
// @Description Generate PDF document for a quotation.
// @Tags Quotation
// @Produce application/pdf
// @Param id path string true "Quotation ID"
// @Success 200 {file} byte "PDF file"
// @Failure 400 {object} responses.ErrorResponse
// @Router /quotations/{id}/pdf [get]
func (h *quotationHandler) GetPDF() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		ctx := r.Context()
		q, err := h.quotationUC.Get(ctx, id)
		if err != nil {
			h.logger.Errorf("Failed to get quotation for PDF: %v", err)
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotFound(err)))
			return
		}

		data := report.QuotationData{
			Company: report.CompanyInfo{
				Name:    "ICMON Auto Repair",
				Address: "123/4 ถนนสุขุมวิท แขวงคลองเตย เขตคลองเตย กรุงเทพฯ 10110",
				Phone:   "02-123-4567",
				TaxID:   "0123456789012",
			},
			QuotationNo:  q.QuotationNo,
			Date:         q.QuotationDate,
			ExpiryDate:   q.ExpiryDate,
			CustomerName: q.CustomerID.String(),
			Subtotal:     q.Subtotal,
			TaxRate:      q.TaxRate,
			TaxAmount:    q.TaxAmount,
			Discount:     q.DiscountValue,
			GrandTotal:   q.Total,
			AmountWords:  report.BahtThai(q.Total),
			CreatedBy:    q.UserID.String(),
		}

		pdf, err := report.GeneratePDF(ctx, report.TplQuotation, data)
		if err != nil {
			h.logger.Errorf("Failed to generate quotation PDF: %v", err)
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrInternalServer(err)))
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "inline; filename=quotation_"+idStr+".pdf")
		w.Write(pdf)
	}
}

func mapModelsToResponse(items []*models.Quotation) []*presenter.QuotationResponse {
	out := make([]*presenter.QuotationResponse, len(items))
	for i, item := range items {
		out[i] = mapModelToResponse(item)
	}
	return out
}
