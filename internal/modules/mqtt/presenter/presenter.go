package presenter

import "time"

type PublishRequest struct {
	Topic    string      `json:"topic" validate:"required"`
	QoS      byte        `json:"qos" validate:"omitempty,min=0,max=2"`
	Retained bool        `json:"retained"`
	Payload  interface{} `json:"payload" validate:"required"`
}

type SubscribeRequest struct {
	Topic string `json:"topic" validate:"required"`
	QoS   byte   `json:"qos" validate:"omitempty,min=0,max=2"`
}

type UnsubscribeRequest struct {
	Topic string `json:"topic" validate:"required"`
}

type UnsubscribeResponse struct {
	Success bool   `json:"success"`
	Topic   string `json:"topic"`
	Message string `json:"message,omitempty"`
}

type SubscriptionsResponse struct {
	Topics []string `json:"topics"`
	Count  int      `json:"count"`
}

type PublishResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type SubscribeResponse struct {
	Topic        string    `json:"topic"`
	QoS          byte      `json:"qos"`
	SubscribedAt time.Time `json:"subscribed_at"`
	Message      string    `json:"message,omitempty"`
}

// MQTTTopicDataResponse สำหรับ GET /mqtt/gettopicdata
type MQTTTopicDataResponse struct {
	StatusCode      int         `json:"statuscode" example:"200"`
	Code            int         `json:"code" example:"200"`
	Topic           string      `json:"topic" example:"AIRCOM1/DATA"`
	Payload         interface{} `json:"payload"` // ✅ remove example:"25.5,60,1013"
	From            string      `json:"from" example:"mqtt"`
	Timestamp       string      `json:"timestamp" example:"2026-06-08T12:36:58+07:00"`
	Message         string      `json:"message" example:"live data"`
	MqttConnected   bool        `json:"mqtt_connected,omitempty" example:"true"`
	CacheEnabled    bool        `json:"cache_enabled,omitempty" example:"true"`
	CacheHit        bool        `json:"cache_hit,omitempty" example:"false"`
	DataLength      int         `json:"data_length,omitempty" example:"128"`
	FetchDurationMs int64       `json:"fetch_duration_ms,omitempty" example:"145"`
	ErrorDetail     string      `json:"error_detail,omitempty" example:"timeout waiting for message"`
}

type DeviceControlRequest struct {
	Topic   string `json:"topic"`   // MQTT topic to publish control message
	Message string `json:"message"` // Control message (e.g., "ON", "OFF", "1", "0")
}

type DeviceControlResponse struct {
	StatusCode      int         `json:"statuscode"`
	Code            int         `json:"code"`
	TopicControl    string      `json:"topic_control"`
	TopicData       string      `json:"topic_data"`
	MessageSent     string      `json:"message_sent"`
	Payload         interface{} `json:"payload"`
	Data            interface{} `json:"data"`
	Status          int         `json:"status"`
	StatusMsg       string      `json:"status_msg"`
	Timestamp       string      `json:"timestamp"`
	Message         string      `json:"message"`
	MessageTh       string      `json:"message_th"`
	From            string      `json:"from"`
	CacheLift       interface{} `json:"cachelift,omitempty"`
	FetchDurationMs int64       `json:"fetch_duration_ms,omitempty"`
	ErrorDetail     string      `json:"error_detail,omitempty"`
}
