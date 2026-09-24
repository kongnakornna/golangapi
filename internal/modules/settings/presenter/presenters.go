package presenter

// Row is a free-form database row returned by list/all endpoints.
type Row map[string]interface{}

// ListResult is the standard paginated payload of the settings module.
// Key names mirror the gistdaapi NestJS envelope payload.
type ListResult struct {
	Page        int                      `json:"page"`
	CurrentPage int                      `json:"currentPage"`
	PageSize    int                      `json:"pageSize"`
	TotalPages  int64                    `json:"totalPages"`
	Total       int64                    `json:"total"`
	Filter      map[string]interface{}   `json:"filter"`
	Data        []map[string]interface{} `json:"data"`
}

func NewListResult(items []map[string]interface{}, total int64, page, pageSize int, filter map[string]interface{}) *ListResult {
	if items == nil {
		items = []map[string]interface{}{}
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	tp := total / int64(pageSize)
	if total%int64(pageSize) > 0 {
		tp++
	}
	if tp < 1 {
		tp = 1
	}
	if filter == nil {
		filter = map[string]interface{}{}
	}
	return &ListResult{
		Page:        page,
		CurrentPage: page,
		PageSize:    pageSize,
		TotalPages:  tp,
		Total:       total,
		Filter:      filter,
		Data:        items,
	}
}

// SendEmailResult is the payload of /sendemail and /testgemail.
type SendEmailResult struct {
	Success   bool   `json:"success"`
	Code      int    `json:"code"`
	Error     string `json:"error,omitempty"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	Content   string `json:"content"`
	Message   string `json:"message"`
	MessageTh string `json:"message_th"`
}

// MqttDataResult is the payload of /mqttdata.
type MqttDataResult struct {
	GetDataFrom string                   `json:"getdataFrom"`
	Payload     []map[string]interface{} `json:"payload"`
}

type Affected struct {
	Affected int64 `json:"affected"`
}
