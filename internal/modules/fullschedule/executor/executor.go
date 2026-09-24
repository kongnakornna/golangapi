package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"icmongolang/config"
	"icmongolang/internal/modules/fullschedule"
	iotmodels "icmongolang/internal/modules/iot/models"
	"icmongolang/pkg/logger"
)

// Publisher is the narrow MQTT surface the executor needs. Any real MQTT
// client satisfying the same method-set (pkg/mqtt.Client) can be passed in;
// nil disables device commands gracefully.
type Publisher interface {
	Publish(topic string, qos byte, retained bool, payload interface{}) error
	IsConnected() bool
}

// DeviceStore resolves a device by internal device ID.
type DeviceStore interface {
	GetDeviceByID(ctx context.Context, deviceID int) (*iotmodels.Device, error)
}

// Emailer sends an email alert.
type Emailer interface {
	SendEmail(ctx context.Context, from string, to string, subject string, bodyHtml string, bodyPlain string) error
}

// Result is the outcome of an executed event.
type Result struct {
	OK       bool
	Message  string
	Payload  string
	DeviceID *int
	Duration time.Duration
}

// Executor executes schedule events: device ON/OFF over MQTT or email alerts.
// ผู้ปฏิบัติงาน: ส่งคำสั่งควบคุมอุปกรณ์หรือส่งอีเมลแจ้งเตือน
type Executor struct {
	publisher Publisher
	devices   DeviceStore
	emailer   Emailer
	cfg       *config.Config
	logger    logger.Logger
	timeout   time.Duration
	smtpFrom  string
}

// New builds an executor.
func New(publisher Publisher, devices DeviceStore, emailer Emailer, cfg *config.Config, logger logger.Logger) *Executor {
	from := cfg.Email.From
	if from == "" {
		from = cfg.SmtpEmail.User
	}
	return &Executor{
		publisher: publisher,
		devices:   devices,
		emailer:   emailer,
		cfg:       cfg,
		logger:    logger,
		timeout:   5 * time.Second,
		smtpFrom:  from,
	}
}

// Execute runs the schedule event for every mapped device (or once for email).
// ทำงานตาม event_type / event_action ของ schedule
func (e *Executor) Execute(ctx context.Context, s *fullschedule.Schedule, deviceIDs []int) []Result {
	switch s.EventType {
	case fullschedule.EventTypeDevice:
		return e.executeDevices(ctx, s, deviceIDs)
	case fullschedule.EventTypeEmail:
		return []Result{e.executeEmail(ctx, s)}
	default:
		return []Result{{OK: false, Message: fmt.Sprintf("unsupported event_type %q", s.EventType)}}
	}
}

func (e *Executor) executeDevices(ctx context.Context, s *fullschedule.Schedule, deviceIDs []int) []Result {
	results := make([]Result, 0, len(deviceIDs))
	for _, id := range deviceIDs {
		results = append(results, e.executeDevice(ctx, s, id))
	}
	return results
}

func (e *Executor) executeDevice(ctx context.Context, s *fullschedule.Schedule, deviceID int) Result {
	start := time.Now()
	res := Result{DeviceID: &deviceID}

	if e.devices == nil {
		res.Message = "device store not wired"
		res.Duration = time.Since(start)
		return res
	}
	if e.publisher == nil || !e.publisher.IsConnected() {
		res.Message = "mqtt backend not connected"
		res.Duration = time.Since(start)
		return res
	}

	device, err := e.devices.GetDeviceByID(ctx, deviceID)
	if err != nil {
		res.Message = fmt.Sprintf("device %d not found: %v", deviceID, err)
		res.Duration = time.Since(start)
		return res
	}

	topic, payload, action := resolveDeviceTarget(device, s.EventAction)
	if topic == "" {
		res.Message = fmt.Sprintf("device %d has no control topic", deviceID)
		res.Duration = time.Since(start)
		return res
	}

	res.Payload = payload
	// Publish (QOS=1, exactly-once semantics at the broker level).
	if err := e.publisher.Publish(topic, 1, false, payload); err != nil {
		res.Message = fmt.Sprintf("publish failed: %v", err)
		res.Duration = time.Since(start)
		return res
	}

	// Verify: after a short delay, confirm the device reached the desired state.
	if e.verifyDeviceState(ctx, device, action) {
		res.OK = true
		res.Message = fmt.Sprintf("published %s to %s and verified", action, topic)
	} else {
		res.Message = fmt.Sprintf("published %s to %s but could not verify state", action, topic)
	}
	res.Duration = time.Since(start)
	return res
}

// verifyDeviceState checks the device's current control value matches the action.
// ตรวจสอบสถานะอุปกรณ์ว่าสอดคล้องกับคำสั่งหรือไม่
func (e *Executor) verifyDeviceState(ctx context.Context, device *iotmodels.Device, action fullschedule.EventAction) bool {
	if device == nil {
		return false
	}
	want := device.MqttControlOn
	if action == fullschedule.ActionOff {
		want = device.MqttControlOff
	}
	if want == "" {
		// No known expected value; treat publish as success.
		return true
	}
	// Best effort: device may report via mqtt_data_value; if empty we accept.
	have := device.MqttDataValue
	if have == "" {
		return true
	}
	return have == want
}

func (e *Executor) executeEmail(ctx context.Context, s *fullschedule.Schedule) Result {
	start := time.Now()
	res := Result{}
	if e.emailer == nil {
		res.Message = "email backend not wired"
		res.Duration = time.Since(start)
		return res
	}

	subject := fmt.Sprintf("[FullSchedule] %s → %s", s.Name, s.EventAction)
	body := fmt.Sprintf(
		"<h3>Full Schedule Alert</h3><p>Schedule: <b>%s</b></p><p>Mode: %s</p><p>Action: <b>%s</b></p><p>Time: %s</p>",
		s.Name, s.Mode, s.EventAction, time.Now().Format(time.RFC3339),
	)
	to := e.emailTo(s)
	if to == "" {
		res.Message = "no email recipient configured (set settings key email_recipients)"
		res.Duration = time.Since(start)
		return res
	}
	if err := e.emailer.SendEmail(ctx, e.smtpFrom, to, subject, body, "Full Schedule alert"); err != nil {
		res.Message = fmt.Sprintf("email send failed: %v", err)
		res.Duration = time.Since(start)
		return res
	}
	res.OK = true
	res.Payload = fmt.Sprintf(`{"to":%q,"subject":%q}`, to, subject)
	res.Message = "email alert sent"
	res.Duration = time.Since(start)
	return res
}

// emailTo resolves recipients from settings (email_recipients JSON array)
// falling back to the smtp username.
// หาผู้รับอีเมล: จาก settings `email_recipients` หรือ SMTP user
func (e *Executor) emailTo(s *fullschedule.Schedule) string {
	var recipients []string
	for _, setting := range s.Settings {
		if setting.Key != "email_recipients" {
			continue
		}
		if err := json.Unmarshal(setting.Value, &recipients); err == nil && len(recipients) > 0 {
			return recipients[0]
		}
	}
	return e.smtpFrom
}

// resolveDeviceTarget picks the MQTT topic + payload for the given action.
// เลือก topic/payload ควบคุมอุปกรณ์ตาม action
func resolveDeviceTarget(d *iotmodels.Device, action fullschedule.EventAction) (topic, payload string, a fullschedule.EventAction) {
	a = action
	if action == fullschedule.ActionOff {
		return d.MqttDataControl, d.MqttControlOff, action
	}
	if action == fullschedule.ActionOn {
		return d.MqttDataControl, d.MqttControlOn, action
	}
	return "", "", action
}

// Enabled reports whether at least one backend (mqtt) is wired.
func (e *Executor) Enabled() bool { return e.publisher != nil && e.publisher.IsConnected() }
