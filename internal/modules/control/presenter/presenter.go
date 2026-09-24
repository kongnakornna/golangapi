package presenter

// ExecuteRequest defines a manual/immediate device control action.
type ExecuteRequest struct {
	DeviceID int `json:"device_id"` // required
	Event    int `json:"event"`     // 1=ON, 0=OFF
	// Optional overrides. When empty, the device's own MQTT config is used.
	Topic          string `json:"topic,omitempty"`
	MqttControlOn  string `json:"mqtt_control_on,omitempty"`
	MqttControlOff string `json:"mqtt_control_off,omitempty"`
}

// ExecuteResult reports the outcome of a single device control action.
type ExecuteResult struct {
	DeviceID   int    `json:"device_id"`
	DeviceName string `json:"device_name,omitempty"`
	Topic      string `json:"topic,omitempty"`
	Payload    string `json:"payload,omitempty"`
	OK         bool   `json:"ok"`
	Message    string `json:"message,omitempty"`
}

// ScheduleRunRequest optionally scopes which schedules to process.
type ScheduleRunRequest struct {
	ScheduleID int  `json:"schedule_id,omitempty"` // 0 = run all active
	Force      bool `json:"force,omitempty"`       // ignore weekday/time match
}

// ScheduleRunResult aggregates a schedule run.
type ScheduleRunResult struct {
	Total     int               `json:"total"`
	Succeeded int               `json:"succeeded"`
	Failed    int               `json:"failed"`
	Items     []ScheduleRunItem `json:"items,omitempty"`
}

// ScheduleRunItem reports one device execution inside a schedule run.
type ScheduleRunItem struct {
	ScheduleID int    `json:"schedule_id"`
	DeviceID   int    `json:"device_id"`
	Topic      string `json:"topic,omitempty"`
	Payload    string `json:"payload,omitempty"`
	OK         bool   `json:"ok"`
	Message    string `json:"message,omitempty"`
}

// HealthResponse reports whether the control service is wired.
type HealthResponse struct {
	Enabled bool `json:"enabled"`
}
