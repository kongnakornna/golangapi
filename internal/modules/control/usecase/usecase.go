package usecase

import (
	"context"
	"errors"
	"time"

	iotmodels "icmongolang/internal/modules/iot/models"
	"icmongolang/internal/modules/control/presenter"
	"icmongolang/internal/modules/control/repository"
	"icmongolang/pkg/logger"
)

// Publisher is the narrow MQTT surface the control module needs. The real
// pkg/mqtt.Client satisfies it; tests provide a lightweight stub.
type Publisher interface {
	Publish(topic string, qos byte, retained bool, payload interface{}) error
	IsConnected() bool
}

// ErrInvalidRequest is returned for a malformed control request.
var ErrInvalidRequest = errors.New("control: invalid execute request")

// ErrDeviceNotFound is returned when the requested device does not exist.
var ErrDeviceNotFound = errors.New("control: device not found")

// ControlUseCase defines the auto-control contract.
type ControlUseCase interface {
	ExecuteControl(ctx context.Context, req *presenter.ExecuteRequest) (*presenter.ExecuteResult, error)
	RunSchedule(ctx context.Context, req *presenter.ScheduleRunRequest) (*presenter.ScheduleRunResult, error)
	Health(ctx context.Context) *presenter.HealthResponse
	// Enabled reports whether MQTT publish backend is wired.
	Enabled() bool
}

type controlUseCase struct {
	repo       repository.Repository
	publisher  Publisher
	logger     logger.Logger
	enabled    bool
}

// NewControlUseCase builds the control use case. repo/publisher may be nil
// (allows graceful no-op) — existing flows never break.
func NewControlUseCase(repo repository.Repository, publisher Publisher, log logger.Logger) ControlUseCase {
	return &controlUseCase{repo: repo, publisher: publisher, logger: log, enabled: publisher != nil}
}

func (u *controlUseCase) Enabled() bool { return u.enabled }

func (u *controlUseCase) Health(ctx context.Context) *presenter.HealthResponse {
	return &presenter.HealthResponse{Enabled: u.enabled}
}

// ExecuteControl publishes an MQTT control message (ON/OFF) for a device.
func (u *controlUseCase) ExecuteControl(ctx context.Context, req *presenter.ExecuteRequest) (*presenter.ExecuteResult, error) {
	if req == nil || req.DeviceID <= 0 {
		return nil, ErrInvalidRequest
	}
	if u.publisher == nil {
		return &presenter.ExecuteResult{DeviceID: req.DeviceID, OK: false, Message: "mqtt backend not wired"}, nil
	}
	device, err := u.repo.GetDeviceByID(ctx, req.DeviceID)
	if err != nil {
		return nil, ErrDeviceNotFound
	}
	topic, payload := u.resolveTarget(req, device)
	if topic == "" {
		return &presenter.ExecuteResult{DeviceID: device.DeviceID, DeviceName: device.DeviceName, OK: false, Message: "no control topic configured"}, nil
	}
	if !u.publisher.IsConnected() {
		return &presenter.ExecuteResult{DeviceID: device.DeviceID, DeviceName: device.DeviceName, Topic: topic, Payload: payload, OK: false, Message: "mqtt not connected"}, nil
	}
	u.logger.Infof("control: device=%d event=%d topic=%s payload=%s", device.DeviceID, req.Event, topic, payload)
	if err := u.publisher.Publish(topic, 1, false, payload); err != nil {
		u.logger.Errorf("control: publish device=%d error: %v", device.DeviceID, err)
		return &presenter.ExecuteResult{DeviceID: device.DeviceID, DeviceName: device.DeviceName, Topic: topic, Payload: payload, OK: false, Message: err.Error()}, nil
	}
	return &presenter.ExecuteResult{DeviceID: device.DeviceID, DeviceName: device.DeviceName, Topic: topic, Payload: payload, OK: true, Message: "published"}, nil
}

// RunSchedule executes active schedules. When req.ScheduleID is set, only that
// schedule runs; otherwise all active schedules. When req.Force is set, the
// weekday/time match is ignored.
func (u *controlUseCase) RunSchedule(ctx context.Context, req *presenter.ScheduleRunRequest) (*presenter.ScheduleRunResult, error) {
	if u.repo == nil {
		return &presenter.ScheduleRunResult{}, nil
	}
	force := req != nil && req.Force
	now := time.Now()

	schedules, err := u.loadSchedules(ctx, req)
	if err != nil {
		return nil, err
	}
	result := &presenter.ScheduleRunResult{}
	for _, sched := range schedules {
		if !force && !schedMatches(sched, now) {
			continue
		}
		devices, err := u.repo.GetScheduleDevices(ctx, sched.ScheduleID)
		if err != nil {
			continue
		}
		if len(devices) == 0 {
			devices = []int{sched.DeviceID}
		}
		result.Total += len(devices)
		for _, deviceID := range devices {
			item := u.runScheduleDevice(ctx, sched, deviceID, now)
			result.Items = append(result.Items, item)
			if item.OK {
				result.Succeeded++
			} else {
				result.Failed++
			}
		}
	}
	return result, nil
}

func (u *controlUseCase) loadSchedules(ctx context.Context, req *presenter.ScheduleRunRequest) ([]iotmodels.Schedule, error) {
	if req != nil && req.ScheduleID > 0 {
		s, err := u.repo.GetScheduleByID(ctx, req.ScheduleID)
		if err != nil {
			return nil, err
		}
		return []iotmodels.Schedule{*s}, nil
	}
	return u.repo.GetActiveSchedules(ctx)
}

func (u *controlUseCase) runScheduleDevice(ctx context.Context, sched iotmodels.Schedule, deviceID int, now time.Time) presenter.ScheduleRunItem {
	execReq := &presenter.ExecuteRequest{DeviceID: deviceID, Event: sched.Event}
	execResult, err := u.ExecuteControl(ctx, execReq)
	item := presenter.ScheduleRunItem{ScheduleID: sched.ScheduleID, DeviceID: deviceID}
	if err != nil {
		item.Message = err.Error()
		item.OK = false
	} else {
		item.Topic = execResult.Topic
		item.Payload = execResult.Payload
		item.OK = execResult.OK
		item.Message = execResult.Message
	}
	// Record schedule process log (nil-safe repo path handled by ExecuteControl).
	u.logSchedule(ctx, sched, deviceID, item, now)
	return item
}

func (u *controlUseCase) logSchedule(ctx context.Context, sched iotmodels.Schedule, deviceID int, item presenter.ScheduleRunItem, now time.Time) {
	if u.repo == nil {
		return
	}
	status := 0
	if item.OK {
		status = 1
	}
	event := "OFF"
	if sched.Event == 1 {
		event = "ON"
	}
	log := &iotmodels.ScheduleProcessLog{
		ScheduleID:         sched.ScheduleID,
		DeviceID:           deviceID,
		ScheduleEventStart: sched.Start,
		Day:                now.Weekday().String(),
		Doday:              now.Format("2006-01-02"),
		DoTime:             now.Format("15:04:05"),
		ScheduleEvent:      event,
		DeviceStatus:       deviceStatus(item.OK),
		Status:             status,
		Date:               now.Format("2006-01-02"),
		Time:               now.Format("15:04:05"),
	}
	if err := u.repo.CreateScheduleLog(ctx, log); err != nil {
		u.logger.Errorf("control: schedule log schedule=%d device=%d error: %v", sched.ScheduleID, deviceID, err)
	}
}

func deviceStatus(ok bool) string {
	if ok {
		return "success"
	}
	return "failed"
}

func (u *controlUseCase) resolveTarget(req *presenter.ExecuteRequest, device *iotmodels.Device) (string, string) {
	topic := req.Topic
	if topic == "" {
		topic = device.MqttDataControl
	}
	payload := ""
	if req.Event == 1 {
		payload = req.MqttControlOn
		if payload == "" {
			payload = device.MqttControlOn
		}
	} else {
		payload = req.MqttControlOff
		if payload == "" {
			payload = device.MqttControlOff
		}
	}
	return topic, payload
}

// schedMatches reports whether the schedule should run today at the given time.
func schedMatches(s iotmodels.Schedule, now time.Time) bool {
	if !dayEnabled(s, now.Weekday()) {
		return false
	}
	return s.Start == now.Format("15:04")
}

func dayEnabled(s iotmodels.Schedule, weekday time.Weekday) bool {
	switch weekday {
	case time.Sunday:
		return s.Sunday == 1
	case time.Monday:
		return s.Monday == 1
	case time.Tuesday:
		return s.Tuesday == 1
	case time.Wednesday:
		return s.Wednesday == 1
	case time.Thursday:
		return s.Thursday == 1
	case time.Friday:
		return s.Friday == 1
	case time.Saturday:
		return s.Saturday == 1
	}
	return false
}
