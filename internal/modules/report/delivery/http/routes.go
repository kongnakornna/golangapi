package http

import (
	"icmongolang/config"
	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/report"

	"github.com/go-chi/chi/v5"
)

func MapReportRoute(router chi.Router, cfg *config.Config, h report.Handlers, mw *middleware.MiddlewareManager) {
	router.Route("/reports", func(r chi.Router) {
		// PDF rendering: full access token via the Authorization header, OR a
		// short-lived ?token= download token so a browser can navigate straight
		// to the PDF (NIT-3). The download token is scoped to one report type
		// (+ optional source) and is validated per-handler via downloadScopeOK.
		r.Group(func(r chi.Router) {
			r.Use(mw.VerifierForReport())
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())
			r.Get("/daily-sales/pdf", h.DailySalesPDF())
			r.Get("/inventory-summary/pdf", h.InventorySummaryPDF())
			r.Get("/customer-list/pdf", h.CustomerListPDF())
			r.Get("/invoice/pdf", h.InvoicePDF())
			r.Get("/credit-note/pdf", h.CreditNotePDF())
			r.Get("/debit-note/pdf", h.DebitNotePDF())
		})
		// Download-token minting requires a full access token only — never a
		// download token (that would let a leaked download token mint others).
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())
			r.Post("/token", h.CreateDownloadToken())
		})
	})
}
