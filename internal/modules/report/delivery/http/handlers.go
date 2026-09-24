package http

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"icmongolang/config"
	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/customer"
	reportmodule "icmongolang/internal/modules/report"
	"icmongolang/internal/modules/report/usecase"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/jwt"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/report"
	"icmongolang/pkg/responses"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type reportHandler struct {
	cfg        *config.Config
	logger     logger.Logger
	customerUC customer.CustomerUseCaseI
	reportUC   usecase.ReportUseCaseI
}

func CreateReportHandler(cfg *config.Config, logger logger.Logger, customerUC customer.CustomerUseCaseI, reportUC usecase.ReportUseCaseI) reportmodule.Handlers {
	return &reportHandler{cfg: cfg, logger: logger, customerUC: customerUC, reportUC: reportUC}
}

func pdfError(w http.ResponseWriter, r *http.Request, log logger.Logger, msg string, err error) {
	log.Errorf("%s: %v", msg, err)
	render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrInternalServer(err)))
}

// downloadScopeOK verifies that a request authenticated via a short-lived
// download token is allowed to render the requested report. Requests that came
// in with a full access token (no download token in context) are always allowed
// here — they were already vetted by the JWT middleware. A download token is
// bound to one report type (and, for the invoice report, to one source id), so
// a token minted for one PDF cannot fetch a different one.
func downloadScopeOK(r *http.Request, reportType, source string) error {
	claims, ok := r.Context().Value(middleware.DownloadTokenCtxKey).(*jwt.DownloadClaims)
	if !ok || claims == nil {
		return nil
	}
	if claims.ReportType != reportType {
		return fmt.Errorf("download token scoped to %q cannot render %q", claims.ReportType, reportType)
	}
	if source != "" && claims.Source != source {
		return fmt.Errorf("download token scoped to source %q cannot render source %q", claims.Source, source)
	}
	return nil
}

func isValidReportType(t string) bool {
	switch t {
	case "daily_sales", "inventory_summary", "customer_list", "invoice", "credit_note", "debit_note":
		return true
	}
	return false
}

// CreateDownloadToken godoc
// @Summary Mint a short-lived download token for a report PDF
// @Description Returns a 5-minute, report-scoped token that can be appended to a /reports/*/pdf URL (?token=) so a browser can render the PDF without the long-lived access token. Requires a full access token.
// @Tags Reports
// @Accept json
// @Produce json
// @Param request body downloadTokenRequest true "Report type + optional source (quotation id for invoice)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /reports/token [post]
func (h *reportHandler) CreateDownloadToken() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := middleware.GetUserFromCtx(r.Context())
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		var req downloadTokenRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(errors.New("invalid request body")))
			return
		}
		if !isValidReportType(req.Type) {
			render.Render(w, r, responses.CreateErrorResponse(errors.New("invalid report type")))
			return
		}
		if req.Type == "invoice" && req.Source == "" {
			render.Render(w, r, responses.CreateErrorResponse(errors.New("source (quotation id) is required for invoice")))
			return
		}

		ttl := 5 * time.Minute
		token, err := jwt.CreateDownloadTokenHS256(user.ID.String(), req.Type, req.Source, h.cfg.Jwt.DownloadTokenSecret, ttl, h.cfg.Jwt.Issuer)
		if err != nil {
			pdfError(w, r, h.logger, "Failed to create download token", err)
			return
		}

		render.JSON(w, r, map[string]interface{}{
			"token":      token,
			"token_type": "download",
			"expires_in": int(ttl.Seconds()),
			"type":       req.Type,
		})
	}
}

type downloadTokenRequest struct {
	Type   string `json:"type"`
	Source string `json:"source,omitempty"`
}

// DailySalesPDF godoc
// @Summary Generate Daily Sales Report PDF
// @Description Generate daily sales report PDF.
// @Tags Reports
// @Produce application/pdf
// @Param date query string false "Date (YYYY-MM-DD)"
// @Success 200 {file} byte "PDF file"
// @Failure 500 {object} responses.ErrorResponse
// @Router /reports/daily-sales/pdf [get]
func (h *reportHandler) DailySalesPDF() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := downloadScopeOK(r, "daily_sales", ""); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrUnauthorized(err)))
			return
		}
		dateStr := r.URL.Query().Get("date")
		reportDate := time.Now()
		if dateStr != "" {
			parsed, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(errors.New("invalid date, expected YYYY-MM-DD")))
				return
			}
			reportDate = parsed
		}

		data := report.DailySalesData{
			Company:     reportmodule.CompanyInfo(),
			Date:        reportDate,
			Sales:       []report.DailySaleRow{},
			TotalRev:    0,
			TotalCost:   0,
			TotalProfit: 0,
			Summary:     report.DailySalesSummary{},
		}

		pdf, err := report.GeneratePDF(r.Context(), report.TplDailySales, data)
		if err != nil {
			pdfError(w, r, h.logger, "Failed to generate daily sales PDF", err)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "inline; filename=daily_sales_"+reportDate.Format("20060102")+".pdf")
		w.Write(pdf)
	}
}

// InventorySummaryPDF godoc
// @Summary Generate Inventory Summary Report PDF
// @Description Generate inventory summary report PDF.
// @Tags Reports
// @Produce application/pdf
// @Success 200 {file} byte "PDF file"
// @Failure 500 {object} responses.ErrorResponse
// @Router /reports/inventory-summary/pdf [get]
func (h *reportHandler) InventorySummaryPDF() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := downloadScopeOK(r, "inventory_summary", ""); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrUnauthorized(err)))
			return
		}
		data := report.InventorySummaryData{
			Company:    reportmodule.CompanyInfo(),
			Date:       time.Now(),
			Items:      []report.InventoryItem{},
			TotalValue: 0,
			TotalQty:   0,
		}

		pdf, err := report.GeneratePDF(r.Context(), report.TplInventorySum, data)
		if err != nil {
			pdfError(w, r, h.logger, "Failed to generate inventory summary PDF", err)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "inline; filename=inventory_summary.pdf")
		w.Write(pdf)
	}
}

// CustomerListPDF godoc
// @Summary Generate Customer List Report PDF
// @Description Generate customer list report PDF from real customer data.
// @Tags Reports
// @Produce application/pdf
// @Success 200 {file} byte "PDF file"
// @Failure 500 {object} responses.ErrorResponse
// @Router /reports/customer-list/pdf [get]
func (h *reportHandler) CustomerListPDF() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := downloadScopeOK(r, "customer_list", ""); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrUnauthorized(err)))
			return
		}
		ctx := r.Context()
		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		customers, err := h.customerUC.GetMultiByUserID(ctx, user.ID, 1000, 0)
		if err != nil {
			pdfError(w, r, h.logger, "Failed to load customers", err)
			return
		}

		items := make([]report.CustomerListItem, 0, len(customers))
		for i, c := range customers {
			name := c.FullName
			if c.DisplayName != nil && *c.DisplayName != "" {
				name = *c.DisplayName
			}
			email := ""
			if c.Email != nil {
				email = *c.Email
			}
			lastVisit := ""
			if c.LastVisitDate != nil {
				lastVisit = c.LastVisitDate.Format("2006-01-02")
			}
			items = append(items, report.CustomerListItem{
				No:         i + 1,
				Code:       c.CustomerCode,
				Name:       name,
				Phone:      c.PhoneNumber,
				Email:      email,
				CarCount:   0,
				LastVisit:  lastVisit,
				TotalSpent: c.TotalSpent,
				Status:     c.Status,
			})
		}

		data := report.CustomerListData{
			Company:    reportmodule.CompanyInfo(),
			Date:       time.Now(),
			Customers:  items,
			TotalCount: len(items),
		}

		pdf, err := report.GeneratePDF(r.Context(), report.TplCustomerList, data)
		if err != nil {
			pdfError(w, r, h.logger, "Failed to generate customer list PDF", err)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "inline; filename=customer_list.pdf")
		w.Write(pdf)
	}
}

// InvoicePDF godoc
// @Summary Generate Invoice PDF from quotation
// @Description Generate invoice PDF from a quotation (source = quotation id).
// @Tags Reports
// @Produce application/pdf
// @Param source query string true "Quotation ID (uuid)"
// @Success 200 {file} byte "PDF file"
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /reports/invoice/pdf [get]
func (h *reportHandler) InvoicePDF() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		source := r.URL.Query().Get("source")
		if source == "" {
			render.Render(w, r, responses.CreateErrorResponse(errors.New("source (quotation id) is required")))
			return
		}
		if err := downloadScopeOK(r, "invoice", source); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrUnauthorized(err)))
			return
		}
		qid, err := uuid.Parse(source)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(errors.New("invalid source (quotation id)")))
			return
		}

		data, err := h.reportUC.BuildInvoice(r.Context(), qid)
		if err != nil {
			pdfError(w, r, h.logger, "Failed to build invoice", err)
			return
		}

		pdf, err := report.GeneratePDF(r.Context(), report.TplInvoice, data)
		if err != nil {
			pdfError(w, r, h.logger, "Failed to generate invoice PDF", err)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "inline; filename=invoice_"+data.InvoiceNo+".pdf")
		w.Write(pdf)
	}
}

// CreditNotePDF godoc
// @Summary Generate Credit Note PDF
// @Description Generate credit note PDF.
// @Tags Reports
// @Produce application/pdf
// @Success 200 {file} byte "PDF file"
// @Failure 500 {object} responses.ErrorResponse
// @Router /reports/credit-note/pdf [get]
func (h *reportHandler) CreditNotePDF() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := downloadScopeOK(r, "credit_note", ""); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrUnauthorized(err)))
			return
		}
		data := report.CreditNoteData{
			Company:      reportmodule.CompanyInfo(),
			CreditNoteNo: "CN-2026-001",
			Date:         time.Now(),
			CustomerName: "ลูกค้า",
			InvoiceNo:    "INV-2026-001",
			Items:        []report.CreditNoteItem{},
		}

		pdf, err := report.GeneratePDF(r.Context(), report.TplCreditNote, data)
		if err != nil {
			pdfError(w, r, h.logger, "Failed to generate credit note PDF", err)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "inline; filename=credit_note.pdf")
		w.Write(pdf)
	}
}

// DebitNotePDF godoc
// @Summary Generate Debit Note PDF
// @Description Generate debit note PDF.
// @Tags Reports
// @Produce application/pdf
// @Success 200 {file} byte "PDF file"
// @Failure 500 {object} responses.ErrorResponse
// @Router /reports/debit-note/pdf [get]
func (h *reportHandler) DebitNotePDF() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := downloadScopeOK(r, "debit_note", ""); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrUnauthorized(err)))
			return
		}
		data := report.DebitNoteData{
			Company:      reportmodule.CompanyInfo(),
			DebitNoteNo:  "DN-2026-001",
			Date:         time.Now(),
			CustomerName: "ลูกค้า",
			InvoiceNo:    "INV-2026-001",
			Items:        []report.DebitNoteItem{},
		}

		pdf, err := report.GeneratePDF(r.Context(), report.TplDebitNote, data)
		if err != nil {
			pdfError(w, r, h.logger, "Failed to generate debit note PDF", err)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "inline; filename=debit_note.pdf")
		w.Write(pdf)
	}
}
