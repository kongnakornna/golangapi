package email

import "net/http"

// Handlers defines HTTP handler methods for email operations.
// ตัวจัดการ HTTP สำหรับอีเมล
type Handlers interface {
	SendEmail() func(w http.ResponseWriter, r *http.Request)
	GetEmailLog() func(w http.ResponseWriter, r *http.Request)
	ListEmailLogs() func(w http.ResponseWriter, r *http.Request)
	GetEmailConfig() func(w http.ResponseWriter, r *http.Request)
	UpdateEmailConfig() func(w http.ResponseWriter, r *http.Request)
}
