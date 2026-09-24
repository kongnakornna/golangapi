package usecase

import (
	"context"
	"errors"
	"strings"

	"icmongolang/config"
	"icmongolang/internal/modules/notifier/presenter"
	"icmongolang/internal/modules/notifier/provider"
	"icmongolang/internal/modules/notifier/repository"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/mqtt"
)

// ErrInvalidRequest is returned for a malformed dispatch request.
var ErrInvalidRequest = errors.New("notifier: invalid dispatch request")

// ErrUnknownChannel is returned for an unsupported channel name.
var ErrUnknownChannel = errors.New("notifier: unknown channel")

// NotifierUseCase defines the notification dispatch contract.
type NotifierUseCase interface {
	Dispatch(ctx context.Context, req *presenter.DispatchRequest) (*presenter.DispatchResponse, error)
	Health(ctx context.Context) *presenter.HealthResponse
	// Enabled reports whether this notifier has any dispatch backend wired.
	Enabled() bool
}

type notifierUseCase struct {
	repo    repository.Repository
	senders map[string]provider.Sender
	logger  logger.Logger
	enabled bool
}

// NewNotifierUseCase builds the notifier use case. repo may be nil (dispatch
// then still proceeds via request-provided targets); appCfg/mqttClient may be
// nil (email/io channels disabled) — existing flows never break.
func NewNotifierUseCase(repo repository.Repository, appCfg *config.Config, mqttClient mqtt.Client, log logger.Logger) NotifierUseCase {
	uc := &notifierUseCase{repo: repo, logger: log, senders: make(map[string]provider.Sender)}
	for _, s := range provider.NewChannelSenders(appCfg, mqttClient) {
		uc.senders[s.Name()] = s
		uc.enabled = true
	}
	return uc
}

func (u *notifierUseCase) Enabled() bool { return u.enabled }

func (u *notifierUseCase) Dispatch(ctx context.Context, req *presenter.DispatchRequest) (*presenter.DispatchResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	channel := strings.ToLower(strings.TrimSpace(req.Channel))
	if channel == "" {
		return nil, ErrInvalidRequest
	}
	sender, ok := u.senders[channel]
	if !ok {
		return nil, ErrUnknownChannel
	}

	// Load channel config (nil-safe repo).
	var cfg repository.ChannelConfig
	if u.repo != nil {
		cfg, _ = u.repo.GetChannelConfig(ctx, channel)
	}

	extra := map[string]string{
		"control_topic":   req.ControlTopic,
		"control_payload": req.MqttControlOn,
	}
	if extra["control_payload"] == "" {
		extra["control_payload"] = req.MqttControlOff
	}
	extra["webhook_url"] = req.WebhookURL
	if req.DeviceID != nil {
		extra["device_id"] = itoa(*req.DeviceID)
	}

	sent, err := sender.Send(ctx, cfg, req.Content, req.Subject, extra)
	if err != nil {
		u.logger.Errorf("notifier: channel=%s send error: %v", channel, err)
		return dispatchedResponse(req, false, err.Error()), nil
	}
	u.logger.Infof("notifier: channel=%s sent=%v status=%d", channel, sent, req.Status)
	return dispatchedResponse(req, sent, ""), nil
}

func (u *notifierUseCase) Health(ctx context.Context) *presenter.HealthResponse {
	return &presenter.HealthResponse{Enabled: u.enabled}
}

func dispatchedResponse(req *presenter.DispatchRequest, ok bool, msg string) *presenter.DispatchResponse {
	if !ok {
		return &presenter.DispatchResponse{Channel: req.Channel, OK: false, Message: msg}
	}
	return &presenter.DispatchResponse{Channel: req.Channel, OK: true, Message: "sent"}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	b := []byte{}
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	if neg {
		b = append([]byte{'-'}, b...)
	}
	return string(b)
}
