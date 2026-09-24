package payment

import "net/http"

type Handlers interface {
	Create() func(w http.ResponseWriter, r *http.Request)
	Get() func(w http.ResponseWriter, r *http.Request)
	List() func(w http.ResponseWriter, r *http.Request)
	GetByInvoice() func(w http.ResponseWriter, r *http.Request)
	Search() func(w http.ResponseWriter, r *http.Request)
	GetOutstanding() func(w http.ResponseWriter, r *http.Request)
	GetHistory() func(w http.ResponseWriter, r *http.Request)
	Refund() func(w http.ResponseWriter, r *http.Request)
	Cancel() func(w http.ResponseWriter, r *http.Request)
	// Receipt
	GetReceipt() func(w http.ResponseWriter, r *http.Request)
	GetReceiptByPayment() func(w http.ResponseWriter, r *http.Request)
	GetReceiptPDF() func(w http.ResponseWriter, r *http.Request)
	CancelReceipt() func(w http.ResponseWriter, r *http.Request)
}
