package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"icmongolang/config"
	"icmongolang/internal/modules/notifier/repository"
	"icmongolang/pkg/mqtt"
	"icmongolang/pkg/sendEmail"
)

// Sender sends a notification through a specific channel.
type Sender interface {
	Name() string
	// Send performs the dispatch; returns (sent, error). sent=false + nil error
	// means the channel is not configured/disabled (handled gracefully).
	Send(ctx context.Context, cfg repository.ChannelConfig, subject, content string, extra map[string]string) (bool, error)
}

// NewChannelSenders builds all channel senders that have the required deps.
// emailSender is nil-safe: if the SMTP email sender is nil, email returns
// (false, nil) so existing flows never break.
func NewChannelSenders(appCfg *config.Config, mqttClient mqtt.Client) []Sender {
	senders := []Sender{
		newDiscordSender(),
		newLineSender(),
		newSmsSender(),
		newIOSender(mqttClient),
	}
	if appCfg != nil {
		senders = append(senders, &emailSender{es: sendEmail.NewEmailSender(appCfg), cfg: appCfg})
	}
	return senders
}

// --- email ---

type emailSender struct {
	es  sendEmail.EmailSender
	cfg *config.Config
}

func (e *emailSender) Name() string { return "email" }

func (e *emailSender) Send(ctx context.Context, cfg repository.ChannelConfig, subject, content string, _ map[string]string) (bool, error) {
	if e.es == nil || e.cfg == nil {
		return false, nil
	}
	from := e.cfg.SmtpEmail.User
	if from == "" {
		from = cfg.Username
	}
	to := cfg.To
	if to == "" {
		to = e.cfg.Email.From
	}
	if from == "" || to == "" {
		return false, nil
	}
	if err := e.es.SendEmail(ctx, from, to, subject, content, content); err != nil {
		return true, fmt.Errorf("email send: %w", err)
	}
	return true, nil
}

// --- sms (generic HTTP gateway) ---

type smsSender struct{}

func newSmsSender() Sender { return &smsSender{} }

func (s *smsSender) Name() string { return "sms" }

func (s *smsSender) Send(ctx context.Context, cfg repository.ChannelConfig, msg, _ string, extra map[string]string) (bool, error) {
	if !cfg.Enabled || cfg.Host == "" {
		return false, nil
	}
	target := extra["phone"]
	if target == "" {
		target = cfg.Username
	}
	if target == "" {
		return false, nil
	}
	endpoint := cfg.Host
	form := url.Values{}
	form.Set("to", target)
	form.Set("message", msg)
	if cfg.APIKey != "" {
		form.Set("apikey", cfg.APIKey)
	}
	if cfg.Originator != "" {
		form.Set("originator", cfg.Originator)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return true, fmt.Errorf("sms request build: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return true, fmt.Errorf("sms send: %w", err)
	}
	defer resp.Body.Close()
	return true, nil
}

// --- line (LINE Notify) ---

type lineSender struct{}

func newLineSender() Sender { return &lineSender{} }

func (l *lineSender) Name() string { return "line" }

func (l *lineSender) Send(ctx context.Context, cfg repository.ChannelConfig, msg, _ string, _ map[string]string) (bool, error) {
	token := cfg.AccessToken
	if token == "" {
		token = cfg.Password
	}
	if !cfg.Enabled || token == "" {
		return false, nil
	}
	form := url.Values{}
	form.Set("message", msg)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://notify-api.line.me/api/notify", strings.NewReader(form.Encode()))
	if err != nil {
		return true, fmt.Errorf("line request build: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return true, fmt.Errorf("line send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return true, fmt.Errorf("line notify status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	return true, nil
}

// --- discord (webhook) ---

type discordSender struct{}

func newDiscordSender() Sender { return &discordSender{} }

func (d *discordSender) Name() string { return "discord" }

func (d *discordSender) Send(ctx context.Context, cfg repository.ChannelConfig, msg, _ string, extra map[string]string) (bool, error) {
	webhookURL := extra["webhook_url"]
	if webhookURL == "" {
		webhookURL = cfg.WebhookURL
	}
	if webhookURL == "" {
		return false, nil
	}
	body, _ := json.Marshal(map[string]string{"content": msg})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, strings.NewReader(string(body)))
	if err != nil {
		return true, fmt.Errorf("discord request build: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return true, fmt.Errorf("discord send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return true, fmt.Errorf("discord webhook status=%d", resp.StatusCode)
	}
	return true, nil
}

// --- io (MQTT device control) ---

type ioSender struct {
	mqttClient mqtt.Client
}

func newIOSender(mqttClient mqtt.Client) Sender {
	return &ioSender{mqttClient: mqttClient}
}

func (i *ioSender) Name() string { return "io" }

func (i *ioSender) Send(ctx context.Context, cfg repository.ChannelConfig, msg, _ string, extra map[string]string) (bool, error) {
	if i.mqttClient == nil || !i.mqttClient.IsConnected() {
		return false, nil
	}
	topic := extra["control_topic"]
	if topic == "" {
		topic = cfg.ControlTopic
	}
	if topic == "" {
		return false, nil
	}
	payload := extra["control_payload"]
	if payload == "" {
		payload = msg
	}
	if err := i.mqttClient.Publish(topic, 1, false, payload); err != nil {
		return true, fmt.Errorf("io mqtt publish: %w", err)
	}
	return true, nil
}
