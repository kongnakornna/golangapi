package purchaseorder

import "net/http"

type Handlers interface {
	Create() http.HandlerFunc
	GetByID() http.HandlerFunc
	List() http.HandlerFunc
	Update() http.HandlerFunc
	Delete() http.HandlerFunc
	Send() http.HandlerFunc
	Confirm() http.HandlerFunc
	Receive() http.HandlerFunc
	Cancel() http.HandlerFunc
	GetPDF() http.HandlerFunc
	GetSuggestions() http.HandlerFunc
	GetStatusHistory() http.HandlerFunc
	CreateFromQuotation() http.HandlerFunc
}
