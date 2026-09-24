package document

import "net/http"

// Handlers defines HTTP handler methods for document management.
// ตัวจัดการ HTTP สำหรับเอกสาร
type Handlers interface {
	Upload() func(w http.ResponseWriter, r *http.Request)
	Download() func(w http.ResponseWriter, r *http.Request)
	List() func(w http.ResponseWriter, r *http.Request)
	Delete() func(w http.ResponseWriter, r *http.Request)
}
