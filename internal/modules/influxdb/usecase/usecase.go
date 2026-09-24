package usecase

import (
	"context"
	"fmt"
	"icmongolang/internal/modules/influxdb/presenter"
	"icmongolang/pkg/helpers"
	"icmongolang/pkg/influxdb"
	"icmongolang/pkg/logger"
	"time"
)

type InfluxUseCase interface {
	WriteData(ctx context.Context, req *presenter.WriteDataRequest) error
	QueryFilter(ctx context.Context, req *presenter.QueryFilterRequest) ([]presenter.DataPointResponse, error)
	Querydevicechart(ctx context.Context, req *presenter.Querydevicechart) (*presenter.DeviceChartResponse, error)
	QueryStatistics(ctx context.Context, req *presenter.StatisticsRequest) (*presenter.StatisticsResponse, error)
}

type influxUseCase struct {
	client *influxdb.InfluxClient
	logger logger.Logger
}

func NewInfluxUseCase(client *influxdb.InfluxClient, log logger.Logger) InfluxUseCase {
	return &influxUseCase{client: client, logger: log}
}

func (u *influxUseCase) WriteData(ctx context.Context, req *presenter.WriteDataRequest) error {
	return u.client.WriteData(req.Measurement, req.Fields, req.Tags)
}

func (u *influxUseCase) QueryFilter(ctx context.Context, req *presenter.QueryFilterRequest) ([]presenter.DataPointResponse, error) {
	params := influxdb.QueryParams{
		Measurement: req.Measurement,
		Field:       req.Field,
		Bucket:      req.Bucket,
		Start:       req.Start,
		Stop:        req.Stop,
		Limit:       req.Limit,
		Offset:      req.Offset,
	}
	results, err := u.client.QueryFilterData(params)
	if err != nil {
		return nil, err
	}
	// โหลด Bangkok timezone (UTC+7)
	TimeLoc := helpers.GetTimeLocation()
	var resp []presenter.DataPointResponse
	for _, r := range results {
		timeStr := ""
		if t, ok := r["_time"]; ok {
			timeStr = helpers.TimeConvertermas(t.(time.Time).In(TimeLoc))
		}
		resp = append(resp, presenter.DataPointResponse{
			Time:  timeStr,
			Value: r["_value"],
		})
	}
	return resp, nil
}

func (u *influxUseCase) Querydevicechart(ctx context.Context, req *presenter.Querydevicechart) (*presenter.DeviceChartResponse, error) {
	params := influxdb.QueryParams{
		Measurement: req.Measurement,
		Field:       req.Field,
		Bucket:      req.Bucket,
		Start:       req.Start,
		Stop:        req.Stop,
		Limit:       req.Limit,
		Offset:      req.Offset,
	}
	results, err := u.client.Querydevicechart(params)
	if err != nil {
		return nil, err
	}

	// Load Bangkok timezone
	TimeLoc := helpers.GetTimeLocation()

	// Prepare slices
	var dataValues []int64
	var dateStrings []string
	var lastValue int64
	var lastTimeStr string

	for _, r := range results {
		// Extract time
		var timeStr string
		if t, ok := r["_time"]; ok {
			tt := t.(time.Time).In(TimeLoc)
			timeStr = helpers.TimeConvertermas(tt) // match sample format
			dateStrings = append(dateStrings, timeStr)
			lastTimeStr = timeStr
		}

		// Extract value (assuming numeric)
		var val int64
		if v, ok := r["_value"]; ok {
			switch vv := v.(type) {
			case float64:
				val = int64(vv)
			case int64:
				val = vv
			case int:
				val = int64(vv)
			}
			dataValues = append(dataValues, val)
			lastValue = val
		}
	}

	// Build response
	resp := &presenter.DeviceChartResponse{
		Bucket: req.Bucket,
		Field:  req.Field,
		Info: presenter.InfoData{
			Bucket:      req.Bucket,
			Measurement: req.Measurement,
			Result:      "last", // or dynamic from query
			Table:       0,
			Field:       req.Field,
			Start:       req.Start,
			Stop:        req.Stop,
			Time:        lastTimeStr,
			Value:       lastValue,
		},
		Data:  dataValues,
		Date:  dateStrings,
		Name:  req.Field,
		Cache: "cache", // or dynamic
	}

	return resp, nil
}

func (u *influxUseCase) QueryStatistics(ctx context.Context, req *presenter.StatisticsRequest) (*presenter.StatisticsResponse, error) {
	params := influxdb.QueryParams{
		Measurement:  req.Measurement,
		Field:        req.Field,
		Bucket:       req.Bucket,
		Start:        req.Start,
		Stop:         req.Stop,
		Mean:         req.Aggregate,
		WindowPeriod: req.WindowPeriod,
		Percentile:   req.Percentile,
	}
	result := u.client.CalculateStatistics(params)
	if !result.Success {
		return nil, fmt.Errorf("%s", result.Error)
	}
	// Build response
	resp := &presenter.StatisticsResponse{
		Aggregate: result.Metadata.Method,
	}
	if len(result.Data) == 1 {
		resp.Value = result.Data[0].Value
	} else {
		var dataPoints []presenter.DataPointResponse
		for _, d := range result.Data {
			dataPoints = append(dataPoints, presenter.DataPointResponse{
				Time:  d.Time,
				Value: d.Value,
			})
		}
		resp.Data = dataPoints
	}
	return resp, nil
}
