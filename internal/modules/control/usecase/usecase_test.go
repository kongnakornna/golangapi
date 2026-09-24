package usecase_test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"icmongolang/internal/modules/control/presenter"
	"icmongolang/internal/modules/control/repository"
	"icmongolang/internal/modules/control/usecase"
	iotmodels "icmongolang/internal/modules/iot/models"
	"icmongolang/pkg/logger"
)

type stubRepo struct {
	devices     map[int]*iotmodels.Device
	schedules   []iotmodels.Schedule
	scheduleDev map[int][]int
	logs        int
}

func (s *stubRepo) GetDeviceByID(ctx context.Context, deviceID int) (*iotmodels.Device, error) {
	if s.devices == nil {
		s.devices = map[int]*iotmodels.Device{}
	}
	if s.devices[deviceID] == nil {
		return nil, errors.New("not found")
	}
	return s.devices[deviceID], nil
}

func (s *stubRepo) GetActiveSchedules(ctx context.Context) ([]iotmodels.Schedule, error) {
	return s.schedules, nil
}

func (s *stubRepo) GetScheduleByID(ctx context.Context, scheduleID int) (*iotmodels.Schedule, error) {
	for i := range s.schedules {
		if s.schedules[i].ScheduleID == scheduleID {
			return &s.schedules[i], nil
		}
	}
	return nil, errors.New("not found")
}

func (s *stubRepo) GetScheduleDevices(ctx context.Context, scheduleID int) ([]int, error) {
	return s.scheduleDev[scheduleID], nil
}

func (s *stubRepo) CreateScheduleLog(ctx context.Context, log *iotmodels.ScheduleProcessLog) error {
	s.logs++
	return nil
}

type stubMQTT struct {
	published []string
	connected bool
}

func (s *stubMQTT) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	s.published = append(s.published, topic+"="+toString(payload))
	return nil
}

func (s *stubMQTT) IsConnected() bool { return s.connected }

func toString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	}
	return strconv.Itoa(0)
}

func newUC(r repository.Repository, m usecase.Publisher) usecase.ControlUseCase {
	return usecase.NewControlUseCase(r, m, logger.NewLogger("control-test"))
}

func TestExecute_InvalidRequest(t *testing.T) {
	uc := newUC(&stubRepo{}, &stubMQTT{connected: true})
	if _, err := uc.ExecuteControl(context.Background(), nil); err == nil {
		t.Fatal("expected error for nil request")
	}
	if _, err := uc.ExecuteControl(context.Background(), &presenter.ExecuteRequest{}); err == nil {
		t.Fatal("expected error for missing device id")
	}
}

func TestExecute_DeviceOn_UsesDeviceConfig(t *testing.T) {
	m := &stubMQTT{connected: true}
	uc := newUC(&stubRepo{devices: map[int]*iotmodels.Device{
		7: {DeviceID: 7, DeviceName: "pump", MqttDataControl: "dev/7/CONTROL", MqttControlOn: "1", MqttControlOff: "0"},
	}}, m)
	resp, err := uc.ExecuteControl(context.Background(), &presenter.ExecuteRequest{DeviceID: 7, Event: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.OK {
		t.Fatalf("expected ok, got %s", resp.Message)
	}
	if resp.Topic != "dev/7/CONTROL" || resp.Payload != "1" {
		t.Fatalf("unexpected target: topic=%s payload=%s", resp.Topic, resp.Payload)
	}
	if len(m.published) != 1 || m.published[0] != "dev/7/CONTROL=1" {
		t.Fatalf("unexpected publish: %v", m.published)
	}
}

func TestExecute_DeviceOff_WithOverride(t *testing.T) {
	m := &stubMQTT{connected: true}
	uc := newUC(&stubRepo{devices: map[int]*iotmodels.Device{
		3: {DeviceID: 3, MqttDataControl: "dev/3/CONTROL", MqttControlOn: "1", MqttControlOff: "0"},
	}}, m)
	resp, err := uc.ExecuteControl(context.Background(), &presenter.ExecuteRequest{
		DeviceID: 3, Event: 0, Topic: "custom/CTRL", MqttControlOff: "OFF",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Topic != "custom/CTRL" || resp.Payload != "OFF" {
		t.Fatalf("unexpected target: topic=%s payload=%s", resp.Topic, resp.Payload)
	}
}

func TestExecute_DeviceNotFound(t *testing.T) {
	uc := newUC(&stubRepo{}, &stubMQTT{connected: true})
	_, err := uc.ExecuteControl(context.Background(), &presenter.ExecuteRequest{DeviceID: 999, Event: 1})
	if err == nil {
		t.Fatal("expected device-not-found error")
	}
}

func TestExecute_NilMQTT_Graceful(t *testing.T) {
	uc := newUC(&stubRepo{devices: map[int]*iotmodels.Device{
		1: {DeviceID: 1, MqttDataControl: "x"},
	}}, nil)
	resp, err := uc.ExecuteControl(context.Background(), &presenter.ExecuteRequest{DeviceID: 1, Event: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.OK {
		t.Fatalf("expected not-ok with nil mqtt")
	}
	if uc.Enabled() {
		t.Fatalf("expected disabled when mqtt is nil")
	}
}

func TestExecute_NotConnected_Graceful(t *testing.T) {
	uc := newUC(&stubRepo{devices: map[int]*iotmodels.Device{
		1: {DeviceID: 1, MqttDataControl: "x"},
	}}, &stubMQTT{connected: false})
	resp, _ := uc.ExecuteControl(context.Background(), &presenter.ExecuteRequest{DeviceID: 1, Event: 1})
	if resp.OK {
		t.Fatalf("expected not-ok when disconnected")
	}
}

func TestRunSchedule_Force_ExecutesAll(t *testing.T) {
	m := &stubMQTT{connected: true}
	repo := &stubRepo{
		devices: map[int]*iotmodels.Device{
			10: {DeviceID: 10, MqttDataControl: "dev/10/CTRL", MqttControlOn: "1", MqttControlOff: "0"},
			20: {DeviceID: 20, MqttDataControl: "dev/20/CTRL", MqttControlOn: "1", MqttControlOff: "0"},
		},
		schedules: []iotmodels.Schedule{
			{ScheduleID: 1, DeviceID: 10, Event: 1, Start: "08:00"},
			{ScheduleID: 2, DeviceID: 20, Event: 0, Start: "18:00"},
		},
		scheduleDev: map[int][]int{1: {10}, 2: {20}},
	}
	uc := newUC(repo, m)
	resp, err := uc.RunSchedule(context.Background(), &presenter.ScheduleRunRequest{Force: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Total != 2 || resp.Succeeded != 2 || resp.Failed != 0 {
		t.Fatalf("unexpected totals: %+v", resp)
	}
	if len(m.published) != 2 {
		t.Fatalf("expected 2 publishes, got %v", m.published)
	}
	if repo.logs != 2 {
		t.Fatalf("expected 2 schedule logs, got %d", repo.logs)
	}
}

func TestRunSchedule_SpecificSchedule(t *testing.T) {
	m := &stubMQTT{connected: true}
	repo := &stubRepo{
		devices: map[int]*iotmodels.Device{
			10: {DeviceID: 10, MqttDataControl: "dev/10/CTRL", MqttControlOn: "1", MqttControlOff: "0"},
			20: {DeviceID: 20, MqttDataControl: "dev/20/CTRL", MqttControlOn: "1", MqttControlOff: "0"},
		},
		schedules: []iotmodels.Schedule{
			{ScheduleID: 1, DeviceID: 10, Event: 1, Start: "08:00"},
			{ScheduleID: 2, DeviceID: 20, Event: 0, Start: "18:00"},
		},
		scheduleDev: map[int][]int{1: {10}, 2: {20}},
	}
	uc := newUC(repo, m)
	resp, err := uc.RunSchedule(context.Background(), &presenter.ScheduleRunRequest{ScheduleID: 2, Force: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Total != 1 || resp.Items[0].DeviceID != 20 {
		t.Fatalf("unexpected result: %+v", resp.Items)
	}
}

func TestRunSchedule_DayMismatch_Skips(t *testing.T) {
	m := &stubMQTT{connected: true}
	repo := &stubRepo{
		schedules: []iotmodels.Schedule{
			// Start obviously not matching "now"; lists only Monday-enabled but force=false
			{ScheduleID: 1, DeviceID: 10, Event: 1, Start: "23:59", Monday: 1},
		},
		scheduleDev: map[int][]int{1: {10}},
	}
	uc := newUC(repo, m)
	resp, _ := uc.RunSchedule(context.Background(), &presenter.ScheduleRunRequest{})
	// Unless now is exactly 23:59 on Monday (unlikely), nothing executes.
	if resp.Total != 0 {
		t.Fatalf("expected skip, got %+v", resp)
	}
}

func TestRunSchedule_NilRepo_Graceful(t *testing.T) {
	uc := newUC(nil, &stubMQTT{connected: true})
	resp, err := uc.RunSchedule(context.Background(), &presenter.ScheduleRunRequest{Force: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Total != 0 {
		t.Fatalf("expected empty result, got %+v", resp)
	}
}

func TestHealth_Enabled(t *testing.T) {
	uc := newUC(&stubRepo{}, &stubMQTT{connected: true})
	if !uc.Health(context.Background()).Enabled {
		t.Fatal("expected enabled=true")
	}
}
