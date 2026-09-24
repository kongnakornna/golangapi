package fullschedule

import "net/http"

// Handlers defines HTTP handler methods for the Full Schedule module.
// ตัวจัดการ HTTP สำหรับโมดูล Full Schedule
type Handlers interface {
	Create() func(w http.ResponseWriter, r *http.Request)
	GetByID() func(w http.ResponseWriter, r *http.Request)
	List() func(w http.ResponseWriter, r *http.Request)
	Update() func(w http.ResponseWriter, r *http.Request)
	Delete() func(w http.ResponseWriter, r *http.Request)
	SetStatus() func(w http.ResponseWriter, r *http.Request)
	UpdateDayStatus() func(w http.ResponseWriter, r *http.Request)
	Trigger() func(w http.ResponseWriter, r *http.Request)
	SetSettings() func(w http.ResponseWriter, r *http.Request)
	GetSettings() func(w http.ResponseWriter, r *http.Request)
	GetHistoryBySchedule() func(w http.ResponseWriter, r *http.Request)
	GetHistory() func(w http.ResponseWriter, r *http.Request)
	GetReport() func(w http.ResponseWriter, r *http.Request)
	ListDevices() func(w http.ResponseWriter, r *http.Request)
	ListSchedulePage() func(w http.ResponseWriter, r *http.Request)
	ListScheduleAll() func(w http.ResponseWriter, r *http.Request)
	ListScheduleDevicePage() func(w http.ResponseWriter, r *http.Request)
}
