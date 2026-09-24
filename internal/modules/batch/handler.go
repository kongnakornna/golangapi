package batch

import "net/http"

// Handlers defines HTTP handler methods for batch jobs.
// ตัวจัดการ HTTP สำหรับงานแบตช์
type Handlers interface {
	CreateJob() func(w http.ResponseWriter, r *http.Request)
	GetJob() func(w http.ResponseWriter, r *http.Request)
	ListJobs() func(w http.ResponseWriter, r *http.Request)
	UpdateJob() func(w http.ResponseWriter, r *http.Request)
	DeleteJob() func(w http.ResponseWriter, r *http.Request)
	RunJobNow() func(w http.ResponseWriter, r *http.Request)
	GetJobLogs() func(w http.ResponseWriter, r *http.Request)
}
