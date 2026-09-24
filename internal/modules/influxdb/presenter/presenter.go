package presenter

type WriteDataRequest struct {
	Measurement string                 `json:"measurement" validate:"required"`
	Fields      map[string]interface{} `json:"fields" validate:"required"`
	Tags        map[string]string      `json:"tags,omitempty"`
}

type QueryFilterRequest struct {
	Measurement string `json:"measurement" validate:"required" example:"temperature"`
	Field       string `json:"field" validate:"required" example:"value"`
	Bucket      string `json:"bucket,omitempty" example:"AIRCOM1"`
	Start       string `json:"start,omitempty" example:"-1h"`
	Stop        string `json:"stop,omitempty" example:"now()"`
	Limit       int    `json:"limit,omitempty" example:"10000"`
	Offset      int    `json:"offset,omitempty" example:"1"`
}
type Querydevicechart struct {
	Measurement string `json:"measurement" validate:"required" example:"temperature"`
	Field       string `json:"field" validate:"required" example:"value"`
	Bucket      string `json:"bucket,omitempty" example:"AIRCOM1"`
	Start       string `json:"start,omitempty" example:"-1h"`
	Stop        string `json:"stop,omitempty" example:"now()"`
	Limit       int    `json:"limit,omitempty" example:"10000"`
	Offset      int    `json:"offset,omitempty" example:"1"`
}
type StatisticsRequest struct {
	Measurement  string  `json:"measurement" validate:"required" example:"temperature"`
	Bucket       string  `json:"bucket" validate:"required" example:"AIRCOM1"`
	Field        string  `json:"field" validate:"required" example:"value"`
	Start        string  `json:"start,omitempty" example:"-15s"`
	Stop         string  `json:"stop,omitempty" example:"now()"`
	Aggregate    string  `json:"aggregate" validate:"required" example:"value"`
	WindowPeriod string  `json:"windowPeriod,omitempty" example:"15s"`
	Percentile   float64 `json:"percentile,omitempty" example:"95"`
}

type DataPointResponse struct {
	Time  string      `json:"time"`
	Value interface{} `json:"value"`
}

type StatisticsResponse struct {
	Aggregate string              `json:"aggregate"`
	Value     interface{}         `json:"value"`
	Data      []DataPointResponse `json:"data,omitempty"`
}

// DeviceChartResponse represents the structure for /influx/devicechart endpoint
type DeviceChartResponse struct {
	Bucket string   `json:"bucket"`
	Field  string   `json:"field"`
	Info   InfoData `json:"info"`
	Data   []int64  `json:"data"`
	Date   []string `json:"date"`
	Name   string   `json:"name"`
	Cache  string   `json:"cache"`
}

// InfoData contains metadata for the chart
type InfoData struct {
	Bucket      string `json:"bucket"`
	Measurement string `json:"measurement"`
	Result      string `json:"result"`
	Table       int    `json:"table"`
	Field       string `json:"field"`
	Start       string `json:"start"`
	Stop        string `json:"stop"`
	Time        string `json:"time"`
	Value       int64  `json:"value"`
}
