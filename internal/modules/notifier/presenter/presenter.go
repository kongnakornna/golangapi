package presenter

// DispatchRequest defines the payload for dispatching an alarm notification
// through configured channels (email/sms/line/discord/io-device).
type DispatchRequest struct {
	Channel        string      `json:"channel"` // email|sms|line|discord|io
	Title          string      `json:"title"`
	Subject        string      `json:"subject"`
	Content        string      `json:"content"`
	Status         int         `json:"status"` // 1..5 (see alarm processor)
	DeviceID       *int        `json:"device_id,omitempty"`
	DeviceName     string      `json:"device_name,omitempty"`
	MqttControlOn  string      `json:"mqtt_control_on,omitempty"`
	MqttControlOff string      `json:"mqtt_control_off,omitempty"`
	ControlTopic   string      `json:"control_topic,omitempty"`
	WebhookURL     string      `json:"webhook_url,omitempty"`
	EventControl   int         `json:"event_control,omitempty"`
	Extra          interface{} `json:"extra,omitempty"`
}

// DispatchResponse reports a single channel dispatch result.
type DispatchResponse struct {
	Channel  string `json:"channel"`
	OK       bool   `json:"ok"`
	Message  string `json:"message,omitempty"`
	LoggedID string `json:"logged_id,omitempty"`
}

// HealthResponse reports whether the notifier service is wired.
type HealthResponse struct {
	Enabled bool `json:"enabled"`
}
