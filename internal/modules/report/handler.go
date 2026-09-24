package report

import "net/http"

type Handlers interface {
	DailySalesPDF() func(w http.ResponseWriter, r *http.Request)
	InventorySummaryPDF() func(w http.ResponseWriter, r *http.Request)
	CustomerListPDF() func(w http.ResponseWriter, r *http.Request)
	InvoicePDF() func(w http.ResponseWriter, r *http.Request)
	CreditNotePDF() func(w http.ResponseWriter, r *http.Request)
	DebitNotePDF() func(w http.ResponseWriter, r *http.Request)
	CreateDownloadToken() func(w http.ResponseWriter, r *http.Request)
}
