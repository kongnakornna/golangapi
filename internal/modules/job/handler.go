package job

import "net/http"

// Handlers defines HTTP handler methods for the Job module.
// ตัวจัดการ HTTP สำหรับโมดูลใบรับงานซ่อม
type Handlers interface {
	Create() func(w http.ResponseWriter, r *http.Request)
	GetByID() func(w http.ResponseWriter, r *http.Request)
	List() func(w http.ResponseWriter, r *http.Request)
	Update() func(w http.ResponseWriter, r *http.Request)
	Delete() func(w http.ResponseWriter, r *http.Request)
	ChangeStatus() func(w http.ResponseWriter, r *http.Request)
	AddService() func(w http.ResponseWriter, r *http.Request)
	AddPart() func(w http.ResponseWriter, r *http.Request)
	GetReport() func(w http.ResponseWriter, r *http.Request)
	GetStatusHistory() func(w http.ResponseWriter, r *http.Request)
	GetServices() func(w http.ResponseWriter, r *http.Request)
	GetParts() func(w http.ResponseWriter, r *http.Request)
	GetPDF() func(w http.ResponseWriter, r *http.Request)
	GetPickingPDF() func(w http.ResponseWriter, r *http.Request)
	GetDeliveryPDF() func(w http.ResponseWriter, r *http.Request)
}
