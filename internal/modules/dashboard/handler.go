package dashboard

import "net/http"

// Handlers defines HTTP handler methods for dashboard endpoints.
// ตัวจัดการ HTTP สำหรับแดชบอร์ด
type Handlers interface {
	GetDashboardStats() func(w http.ResponseWriter, r *http.Request)
	GetRevenueChart() func(w http.ResponseWriter, r *http.Request)
	GetTopParts() func(w http.ResponseWriter, r *http.Request)
	GetJobStatusSummary() func(w http.ResponseWriter, r *http.Request)
}
