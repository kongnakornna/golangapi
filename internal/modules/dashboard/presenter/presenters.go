package presenter

// DashboardResponse represents overall dashboard statistics.
// สถิติภาพรวมของแดชบอร์ด
type DashboardResponse struct {
	TotalDevices  int64 `json:"totalDevices"`
	OnlineDevices int64 `json:"onlineDevices"`
	ActiveAlerts  int64 `json:"activeAlerts"`
	TodayCommands int64 `json:"todayCommands"`
}

// RevenueData represents revenue for a specific period.
// ข้อมูลรายได้ตามช่วงเวลา
type RevenueData struct {
	Period string  `json:"period"`
	Amount float64 `json:"amount"`
}

// TopPartData represents a top-selling part.
// อะไหล่ที่ขายดีที่สุด
type TopPartData struct {
	PartName string `json:"partName"`
	Count    int64  `json:"count"`
}

// JobStatusSummary represents job status counts.
// จำนวนงานตามสถานะ
type JobStatusSummary struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}
