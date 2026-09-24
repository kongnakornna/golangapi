package influxdb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"icmongolang/config"
	"icmongolang/pkg/logger"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxClient struct {
	client   influxdb2.Client
	writeApi api.WriteAPI
	queryApi api.QueryAPI
	cfg      *config.Config
	logger   logger.Logger
	org      string
	bucket   string
	username string
	password string
}

type QueryParams struct {
	Start           string
	Stop            string
	Bucket          string
	Measurement     string
	Field           string
	Limit           int
	Offset          int
	WindowPeriod    string
	Mean            string
	TzString        string
	FailSilently    bool
	FallbackOnError bool
	Percentile      float64
}

type StatisticalResult struct {
	Type       string      `json:"type"`
	Value      interface{} `json:"value"`
	Time       string      `json:"time,omitempty"`
	DataPoints int         `json:"dataPoints,omitempty"`
}

type MeanCalculationResult struct {
	Success  bool                `json:"success"`
	Data     []StatisticalResult `json:"data"`
	Summary  *SummaryStats       `json:"summary,omitempty"`
	Error    string              `json:"error,omitempty"`
	Metadata *QueryMetadata      `json:"metadata,omitempty"`
}

type SummaryStats struct {
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
	Avg      float64 `json:"avg"`
	Count    int     `json:"count"`
	StdDev   float64 `json:"stdDev"`
	Variance float64 `json:"variance"`
	Median   float64 `json:"median"`
	P95      float64 `json:"p95"`
	P99      float64 `json:"p99"`
}

type QueryMetadata struct {
	QueryTime string `json:"queryTime"`
	Duration  int64  `json:"duration"`
	Method    string `json:"method"`
}

type CountResult struct {
	Total  int64  `json:"total"`
	Method string `json:"method,omitempty"`
	Error  string `json:"error,omitempty"`
}

type Client interface {
	WritePoint(measurement string, tags map[string]string, fields map[string]interface{}, t time.Time) error
	Close()
}

type influxClient struct {
	client   influxdb2.Client
	writeAPI api.WriteAPI
}

func NewInfluxClient(cfg *config.Config, log logger.Logger) (*InfluxClient, error) {
	url := cfg.InfluxDB.URL
	token := cfg.InfluxDB.Token
	org := cfg.InfluxDB.Org
	bucket := cfg.InfluxDB.Bucket
	username := cfg.InfluxDB.Username
	password := cfg.InfluxDB.Password

	if url == "" || token == "" {
		return nil, fmt.Errorf("influxdb config missing (url and token are required; set INFLUXDB_URL / INFLUXDB_TOKEN)")
	}

	timeout := cfg.InfluxDB.Timeout
	if timeout <= 0 {
		timeout = 30
	}

	client := influxdb2.NewClientWithOptions(url, token,
		influxdb2.DefaultOptions().SetHTTPRequestTimeout(uint(timeout)),
	)
	writeApi := client.WriteAPI(org, bucket)
	queryApi := client.QueryAPI(org)

	return &InfluxClient{
		client:   client,
		writeApi: writeApi,
		queryApi: queryApi,
		cfg:      cfg,
		logger:   log,
		org:      org,
		bucket:   bucket,
		username: username,
		password: password,
	}, nil
}

// DashboardURL returns the InfluxDB UI URL the user can open to log in with
// the configured username/password (InfluxDB 2.x API calls use the token instead).
func (i *InfluxClient) DashboardURL() string {
	if i.cfg == nil {
		return ""
	}
	return i.cfg.InfluxDB.URL
}

// Credentials returns the configured InfluxDB UI credentials.
func (i *InfluxClient) Credentials() (username, password string) {
	return i.username, i.password
}

func (i *InfluxClient) Close() {
	i.client.Close()
}

func (i *InfluxClient) WriteData(measurement string, fields map[string]interface{}, tags map[string]string) error {
	point := influxdb2.NewPoint(measurement, tags, fields, time.Now())
	i.writeApi.WritePoint(point)
	i.writeApi.Flush()
	return nil
}

// 🔧 WritePoint ทำให้ InfluxClient implement Client interface
func (i *InfluxClient) WritePoint(measurement string, tags map[string]string, fields map[string]interface{}, t time.Time) error {
	point := influxdb2.NewPoint(measurement, tags, fields, t)
	i.writeApi.WritePoint(point)
	i.writeApi.Flush()
	return nil
}

// WritePointToBucket writes a point to an explicit bucket (falls back to the configured bucket when empty).
func (i *InfluxClient) WritePointToBucket(bucket, measurement string, tags map[string]string, fields map[string]interface{}, t time.Time) error {
	if bucket == "" {
		bucket = i.bucket
	}
	wa := i.client.WriteAPI(i.org, bucket)
	point := influxdb2.NewPoint(measurement, tags, fields, t)
	wa.WritePoint(point)
	wa.Flush()
	return nil
}

// QueryFilterData ดึงข้อมูลตาม measurement, field, ช่วงเวลา และ limit/offset
func (i *InfluxClient) QueryFilterData(params QueryParams) ([]map[string]interface{}, error) {
	bucket := params.Bucket
	if bucket == "" {
		bucket = i.bucket
	}
	start := params.Start
	if start == "" {
		start = "-1h"
	}
	stop := params.Stop
	if stop == "" {
		stop = "now()"
	}
	limit := params.Limit
	if limit <= 0 {
		limit = 1000
	}
	offset := params.Offset
	if offset < 0 {
		offset = 0
	}

	if !validRangeExpr(start) || !validRangeExpr(stop) {
		return nil, fmt.Errorf("invalid range expression")
	}

	fluxQuery := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r["_measurement"] == "%s")
			|> filter(fn: (r) => r["_field"] == "%s")
			|> limit(n: %d, offset: %d)
			|> yield(name: "filtered_data")`,
		fluxEscape(bucket), start, stop, fluxEscape(params.Measurement), fluxEscape(params.Field), limit, offset)

	return i.executeQuery(fluxQuery, params.TzString)
}

// Querydevicechart ดึงข้อมูลตาม measurement, field, ช่วงเวลา และ limit/offset
func (i *InfluxClient) Querydevicechart(params QueryParams) ([]map[string]interface{}, error) {
	bucket := params.Bucket
	if bucket == "" {
		bucket = i.bucket
	}
	start := params.Start
	if start == "" {
		start = "-1h"
	}
	stop := params.Stop
	if stop == "" {
		stop = "now()"
	}
	limit := params.Limit
	if limit <= 0 {
		limit = 1000
	}
	offset := params.Offset
	if offset < 0 {
		offset = 0
	}

	if !validRangeExpr(start) || !validRangeExpr(stop) {
		return nil, fmt.Errorf("invalid range expression")
	}

	fluxQuery := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r["_measurement"] == "%s")
			|> filter(fn: (r) => r["_field"] == "%s")
			|> limit(n: %d, offset: %d)
			|> yield(name: "filtered_data")`,
		fluxEscape(bucket), start, stop, fluxEscape(params.Measurement), fluxEscape(params.Field), limit, offset)

	return i.executeQuery(fluxQuery, params.TzString)
}

// QueryFilterDataRs ดึงข้อมูลแบบเรียงลำดับตามเวลา (ascending)
func (i *InfluxClient) QueryFilterDataRs(params QueryParams) ([]map[string]interface{}, error) {
	bucket := params.Bucket
	if bucket == "" {
		bucket = i.bucket
	}
	start := params.Start
	if start == "" {
		start = "-1h"
	}
	stop := params.Stop
	if stop == "" {
		stop = "now()"
	}
	limit := params.Limit
	if limit <= 0 {
		limit = 100
	}
	offset := params.Offset
	if offset < 0 {
		offset = 0
	}

	if !validRangeExpr(start) || !validRangeExpr(stop) {
		return nil, fmt.Errorf("invalid range expression")
	}

	fluxQuery := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r["_measurement"] == "%s")
			|> filter(fn: (r) => r["_field"] == "%s")
			|> sort(columns: ["_time"], desc: false)
			|> limit(n: %d, offset: %d)
			|> yield(name: "sorted_data")`,
		fluxEscape(bucket), start, stop, fluxEscape(params.Measurement), fluxEscape(params.Field), limit, offset)

	return i.executeQuery(fluxQuery, params.TzString)
}

func (i *InfluxClient) CountRows(params QueryParams) (CountResult, error) {
	bucket := params.Bucket
	if bucket == "" {
		bucket = i.bucket
	}
	start := params.Start
	if start == "" {
		start = "-30d"
	}
	stop := params.Stop
	if stop == "" {
		stop = "now()"
	}

	if !validRangeExpr(start) || !validRangeExpr(stop) {
		return CountResult{Error: "invalid range expression"}, fmt.Errorf("invalid range expression")
	}

	fluxQuery := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r["_measurement"] == "%s")
			|> filter(fn: (r) => r["_field"] == "%s")
			|> count()
			|> yield(name: "count")`,
		fluxEscape(bucket), start, stop, fluxEscape(params.Measurement), fluxEscape(params.Field))

	result, err := i.executeQuery(fluxQuery, params.TzString)
	if err != nil {
		return CountResult{Error: err.Error()}, err
	}
	if len(result) > 0 {
		if val, ok := result[0]["_value"].(float64); ok {
			return CountResult{Total: int64(val), Method: "direct_count"}, nil
		}
	}
	return CountResult{Total: 0, Method: "no_data"}, nil
}

// CalculateStatistics คำนวณสถิติ (mean, median, percentile, etc.)
func (i *InfluxClient) CalculateStatistics(params QueryParams) MeanCalculationResult {
	startTime := time.Now()

	bucket := params.Bucket
	if bucket == "" {
		bucket = i.bucket
	}
	start := params.Start
	if start == "" {
		start = "-15s"
	}
	stop := params.Stop
	if stop == "" {
		stop = "now()"
	}
	meanType := params.Mean
	if meanType == "" {
		meanType = "last"
	}
	windowPeriod := params.WindowPeriod
	if windowPeriod == "" {
		windowPeriod = "15s"
	}

	var fluxQuery string
	var method string

	if !validRangeExpr(start) || !validRangeExpr(stop) || !validRangeExpr(windowPeriod) {
		return MeanCalculationResult{
			Success: false,
			Error:   "invalid range or window period expression",
			Metadata: &QueryMetadata{
				QueryTime: time.Now().Format(time.RFC3339),
				Duration:  time.Since(startTime).Milliseconds(),
				Method:    meanType,
			},
		}
	}

	switch meanType {
	case "mean", "average":
		method = "mean"
		fluxQuery = fmt.Sprintf(`
			from(bucket: "%s")
				|> range(start: %s, stop: %s)
				|> filter(fn: (r) => r["_measurement"] == "%s")
				|> filter(fn: (r) => r["_field"] == "%s")
				|> aggregateWindow(every: %s, fn: mean, createEmpty: false)
				|> yield(name: "mean")`,
			fluxEscape(bucket), start, stop, fluxEscape(params.Measurement), fluxEscape(params.Field), windowPeriod)
	case "median":
		method = "median"
		fluxQuery = fmt.Sprintf(`
			from(bucket: "%s")
				|> range(start: %s, stop: %s)
				|> filter(fn: (r) => r["_measurement"] == "%s")
				|> filter(fn: (r) => r["_field"] == "%s")
				|> aggregateWindow(every: %s, fn: median, createEmpty: false)
				|> yield(name: "median")`,
			fluxEscape(bucket), start, stop, fluxEscape(params.Measurement), fluxEscape(params.Field), windowPeriod)
	case "mode":
		method = "mode"
		fluxQuery = fmt.Sprintf(`
			from(bucket: "%s")
				|> range(start: %s, stop: %s)
				|> filter(fn: (r) => r["_measurement"] == "%s")
				|> filter(fn: (r) => r["_field"] == "%s")
				|> group(columns: ["_value"])
				|> count()
				|> group(columns: ["_measurement"])
				|> top(n: 1, columns: ["_value"])
				|> yield(name: "mode")`,
			fluxEscape(bucket), start, stop, fluxEscape(params.Measurement), fluxEscape(params.Field))
	case "last":
		method = "last"
		fluxQuery = fmt.Sprintf(`
			from(bucket: "%s")
				|> range(start: %s, stop: %s)
				|> filter(fn: (r) => r["_measurement"] == "%s")
				|> filter(fn: (r) => r["_field"] == "%s")
				|> last()
				|> yield(name: "last")`,
			fluxEscape(bucket), start, stop, fluxEscape(params.Measurement), fluxEscape(params.Field))
	case "first":
		method = "first"
		fluxQuery = fmt.Sprintf(`
			from(bucket: "%s")
				|> range(start: %s, stop: %s)
				|> filter(fn: (r) => r["_measurement"] == "%s")
				|> filter(fn: (r) => r["_field"] == "%s")
				|> first()
				|> yield(name: "first")`,
			fluxEscape(bucket), start, stop, fluxEscape(params.Measurement), fluxEscape(params.Field))
	case "stddev", "standarddeviation":
		method = "stddev"
		fluxQuery = fmt.Sprintf(`
			from(bucket: "%s")
				|> range(start: %s, stop: %s)
				|> filter(fn: (r) => r["_measurement"] == "%s")
				|> filter(fn: (r) => r["_field"] == "%s")
				|> stddev()
				|> yield(name: "stddev")`,
			fluxEscape(bucket), start, stop, fluxEscape(params.Measurement), fluxEscape(params.Field))
	case "variance":
		method = "variance"
		fluxQuery = fmt.Sprintf(`
			from(bucket: "%s")
				|> range(start: %s, stop: %s)
				|> filter(fn: (r) => r["_measurement"] == "%s")
				|> filter(fn: (r) => r["_field"] == "%s")
				|> variance()
				|> yield(name: "variance")`,
			fluxEscape(bucket), start, stop, fluxEscape(params.Measurement), fluxEscape(params.Field))
	case "percentile":
		method = "percentile"
		percentile := params.Percentile
		if percentile <= 0 {
			percentile = 0.95
		}
		fluxQuery = fmt.Sprintf(`
			from(bucket: "%s")
				|> range(start: %s, stop: %s)
				|> filter(fn: (r) => r["_measurement"] == "%s")
				|> filter(fn: (r) => r["_field"] == "%s")
				|> percentile(percentile: %f)
				|> yield(name: "percentile_%f")`,
			fluxEscape(bucket), start, stop, fluxEscape(params.Measurement), fluxEscape(params.Field), percentile, percentile)
	default:
		method = "last"
		fluxQuery = fmt.Sprintf(`
			from(bucket: "%s")
				|> range(start: %s, stop: %s)
				|> filter(fn: (r) => r["_measurement"] == "%s")
				|> filter(fn: (r) => r["_field"] == "%s")
				|> last()
				|> yield(name: "last")`,
			fluxEscape(bucket), start, stop, fluxEscape(params.Measurement), fluxEscape(params.Field))
	}

	results, err := i.executeQuery(fluxQuery, params.TzString)
	if err != nil {
		return MeanCalculationResult{
			Success: false,
			Error:   err.Error(),
			Metadata: &QueryMetadata{
				QueryTime: time.Now().Format(time.RFC3339),
				Duration:  time.Since(startTime).Milliseconds(),
				Method:    method,
			},
		}
	}

	var stats []StatisticalResult
	for _, r := range results {
		if val, ok := r["_value"]; ok {
			stat := StatisticalResult{
				Type:  method,
				Value: val,
			}
			if t, ok := r["_time"]; ok {
				stat.Time = fmt.Sprintf("%v", t)
			}
			stats = append(stats, stat)
		}
	}

	var summary *SummaryStats
	if len(results) > 1 && (method == "mean" || method == "median") {
		summary, _ = i.calculateSummary(params)
	}

	return MeanCalculationResult{
		Success: true,
		Data:    stats,
		Summary: summary,
		Metadata: &QueryMetadata{
			QueryTime: time.Now().Format(time.RFC3339),
			Duration:  time.Since(startTime).Milliseconds(),
			Method:    method,
		},
	}
}

func (i *InfluxClient) calculateSummary(params QueryParams) (*SummaryStats, error) {
	bucket := params.Bucket
	if bucket == "" {
		bucket = i.bucket
	}
	start := params.Start
	if start == "" {
		start = "-1h"
	}
	stop := params.Stop
	if stop == "" {
		stop = "now()"
	}

	if !validRangeExpr(start) || !validRangeExpr(stop) {
		return nil, fmt.Errorf("invalid range expression")
	}

	fluxQuery := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r["_measurement"] == "%s")
			|> filter(fn: (r) => r["_field"] == "%s")
			|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
			|> yield(name: "summary")`,
		fluxEscape(bucket), start, stop, fluxEscape(params.Measurement), fluxEscape(params.Field))

	results, err := i.executeQuery(fluxQuery, params.TzString)
	if err != nil || len(results) == 0 {
		return nil, err
	}

	var values []float64
	for _, r := range results {
		if val, ok := r["_value"].(float64); ok {
			values = append(values, val)
		}
	}
	if len(values) == 0 {
		return nil, nil
	}

	// Calculate statistics
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	// Sort for percentiles
	sorted := make([]float64, len(values))
	copy(sorted, values)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	min := sorted[0]
	max := sorted[len(sorted)-1]
	median := sorted[len(sorted)/2]

	// Standard deviation
	variance := 0.0
	for _, v := range values {
		variance += (v - mean) * (v - mean)
	}
	variance /= float64(len(values))
	stdDev := 0.0
	if variance > 0 {
		stdDev = sqrt(variance)
	}

	// Percentiles
	p95Idx := int(float64(len(sorted)) * 0.95)
	p99Idx := int(float64(len(sorted)) * 0.99)
	p95 := sorted[p95Idx]
	p99 := sorted[p99Idx]

	return &SummaryStats{
		Min:      min,
		Max:      max,
		Avg:      mean,
		Count:    len(values),
		StdDev:   stdDev,
		Variance: variance,
		Median:   median,
		P95:      p95,
		P99:      p99,
	}, nil
}

func (i *InfluxClient) queryTimeout() time.Duration {
	t := i.cfg.InfluxDB.Timeout
	if t <= 0 {
		t = 30
	}
	return time.Duration(t) * time.Second
}

// validRangeExpr reports whether s is a safe Flux range start/stop expression
// (duration like "-1h", "now()", or an RFC3339 timestamp).
func validRangeExpr(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		case r == '-' || r == ':' || r == '.' || r == ',' || r == '+' || r == '(' || r == ')' || r == ' ' || r == '_':
		default:
			return false
		}
	}
	return true
}

// fluxEscape escapes a value for embedding inside a Flux double-quoted string
// literal, so user-supplied strings cannot terminate the literal early (breaking
// the query) or inject arbitrary Flux expressions.
func fluxEscape(s string) string {
	r := strings.NewReplacer(
		`\`, `\\`,
		`"`, `\"`,
		"\n", `\n`,
		"\r", `\r`,
		"\t", `\t`,
	)
	return r.Replace(s)
}

func (i *InfluxClient) executeQuery(fluxQuery string, tzString string) ([]map[string]interface{}, error) {
	i.logger.Infof("Executing InfluxDB query: %s", fluxQuery)

	ctx, cancel := context.WithTimeout(context.Background(), i.queryTimeout())
	defer cancel()

	var results []map[string]interface{}
	result, err := i.queryApi.Query(ctx, fluxQuery)
	if err != nil {
		return nil, fmt.Errorf("🔴 query failed: %w", err)
	}
	defer result.Close()
	if result.Err() != nil {
		return nil, fmt.Errorf("🔴 query error: %w", result.Err())
	}
	for result.Next() {
		record := result.Record()
		row := map[string]interface{}{
			"_time":        record.Time(),
			"_measurement": record.Measurement(),
			"_field":       record.Field(),
			"_value":       record.Value(),
		}
		if tzString != "" {
			// Timezone conversion can be added here if needed
		}
		results = append(results, row)
	}

	if result.Err() != nil {
		return nil, fmt.Errorf("🔴 query error: %w", result.Err())
	}

	return results, nil
}

// sqrt is a simple implementation for float64
func sqrt(x float64) float64 {
	z := 1.0
	for i := 0; i < 100; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}

// Bucket returns the bucket name
func (i *InfluxClient) Bucket() string {
	return i.bucket
}

// New สำหรับใช้งาน Client interface (optional)
func New(cfg *config.InfluxDBConfig) (Client, error) {
	client := influxdb2.NewClient(cfg.URL, cfg.Token)
	_, err := client.Health(context.Background())
	if err != nil {
		return nil, err
	}
	writeAPI := client.WriteAPI(cfg.Org, cfg.Bucket)
	return &influxClient{
		client:   client,
		writeAPI: writeAPI,
	}, nil
}

func (ic *influxClient) WritePoint(measurement string, tags map[string]string, fields map[string]interface{}, t time.Time) error {
	p := influxdb2.NewPoint(measurement, tags, fields, t)
	ic.writeAPI.WritePoint(p)
	return nil
}

func (ic *influxClient) Close() {
	ic.writeAPI.Flush()
	ic.client.Close()
}
