// Package presenter จัดเตรียม HTTP request/response DTOs ของโมดูล mqtt
package presenter

import "time"

// ============================================================
// Publish
// ============================================================

// PublishRequest — body ของ POST /mqtt/publish
type PublishRequest struct {
	Topic    string      `json:"topic"    validate:"required"               example:"BAACTW05/CONTROL"`
	QoS      byte        `json:"qos"      validate:"omitempty,min=0,max=2"  example:"1" minimum:"0" maximum:"2"`
	Retained bool        `json:"retained"                                   example:"false"`
	Payload  interface{} `json:"payload"  validate:"required" swaggertype:"string" example:"ON"`
}

// PublishResponse — response ของ POST /mqtt/publish
type PublishResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message,omitempty" example:"published successfully"`
}

// ============================================================
// Subscribe / Unsubscribe
// ============================================================

// SubscribeRequest — body ของ POST /mqtt/subscribe
type SubscribeRequest struct {
	Topic string `json:"topic" validate:"required"              example:"BAACTW05/DATA"`
	QoS   byte   `json:"qos"   validate:"omitempty,min=0,max=2" example:"1" minimum:"0" maximum:"2"`
}

// SubscribeResponse — response ของ POST /mqtt/subscribe
type SubscribeResponse struct {
	Topic        string    `json:"topic"         example:"BAACTW05/DATA"`
	QoS          byte      `json:"qos"           example:"1"`
	SubscribedAt time.Time `json:"subscribed_at" example:"2026-10-04T12:00:00Z"`
	Message      string    `json:"message,omitempty" example:"subscribed successfully"`
}

// UnsubscribeRequest — body ของ POST /mqtt/unsubscribe
type UnsubscribeRequest struct {
	Topic string `json:"topic" validate:"required" example:"BAACTW05/DATA"`
}

// UnsubscribeResponse — response ของ POST /mqtt/unsubscribe
type UnsubscribeResponse struct {
	Success bool   `json:"success" example:"true"`
	Topic   string `json:"topic"   example:"BAACTW05/DATA"`
	Message string `json:"message,omitempty" example:"unsubscribed successfully"`
}

// SubscriptionsResponse — response ของ GET /mqtt/subscriptions
type SubscriptionsResponse struct {
	Topics []string `json:"topics" example:"BAACTW05/DATA,BAACTW05/ALARM"`
	Count  int      `json:"count"  example:"2"`
}

// ============================================================
// Status
// ============================================================

// StatusResponse — response ของ GET /mqtt/status
//
// ใช้แทน map[string]interface{} เพื่อให้ Swagger generate schema ได้
type StatusResponse struct {
	Connected bool   `json:"connected" example:"true"`
	Status    string `json:"status"    example:"connected" enums:"connected,disconnected"`
	Timestamp string `json:"timestamp" example:"2026-10-04T12:00:00+07:00"`
}

// ============================================================
// Topic data
// ============================================================

// MQTTTopicDataResponse — response ของ GET /mqtt/gettopicdata
type MQTTTopicDataResponse struct {
	StatusCode      int         `json:"statuscode"                    example:"200"`
	Code            int         `json:"code"                          example:"200"`
	Topic           string      `json:"topic"                         example:"AIRCOM1/DATA"`
	Payload         interface{} `json:"payload"                       swaggertype:"string" example:"25.5,60,1013"`
	From            string      `json:"from"                          example:"mqtt" enums:"mqtt,cache,fallback,error"`
	Timestamp       string      `json:"timestamp"                     example:"2026-06-08T12:36:58+07:00"`
	Message         string      `json:"message"                       example:"live data"`
	MqttConnected   bool        `json:"mqtt_connected,omitempty"      example:"true"`
	CacheEnabled    bool        `json:"cache_enabled,omitempty"       example:"true"`
	CacheHit        bool        `json:"cache_hit,omitempty"           example:"false"`
	DataLength      int         `json:"data_length,omitempty"         example:"128"`
	FetchDurationMs int64       `json:"fetch_duration_ms,omitempty"   example:"145"`
	ErrorDetail     string      `json:"error_detail,omitempty"        example:"timeout waiting for message"`
}

// ============================================================
// Device control
// ============================================================

// DeviceControlRequest — internal payload (ไม่ใช่ HTTP body แต่ใช้ใน usecase)
type DeviceControlRequest struct {
	Topic   string `json:"topic"   example:"BAACTW02/CONTROL"`
	Message string `json:"message" example:"ON"`
}

// DeviceControlResponse — response ของ GET /mqtt/devicecontrol
type DeviceControlResponse struct {
	StatusCode      int         `json:"statuscode"                    example:"200"`
	Code            int         `json:"code"                          example:"200"`
	TopicControl    string      `json:"topic_control"                 example:"BAACTW02/CONTROL"`
	TopicData       string      `json:"topic_data"                    example:"BAACTW02/DATA"`
	MessageSent     string      `json:"message_sent"                  example:"ON"`
	Payload         interface{} `json:"payload"                       swaggertype:"string" example:"25.5,60,1013"`
	Data            interface{} `json:"data"                          swaggertype:"array,string" example:"25.5"`
	Status          int         `json:"status"                        example:"1" enums:"0,1"`
	StatusMsg       string      `json:"status_msg"                    example:"ON" enums:"ON,OFF"`
	Timestamp       string      `json:"timestamp"                     example:"2026-06-08T12:36:58+07:00"`
	Message         string      `json:"message"                       example:"Control sent, response received"`
	MessageTh       string      `json:"message_th"                    example:"ส่งคำสั่งและได้รับข้อมูลตอบกลับ"`
	From            string      `json:"from"                          example:"mqtt" enums:"mqtt,cache,fallback"`
	CacheLift       interface{} `json:"cachelift,omitempty"`
	FetchDurationMs int64       `json:"fetch_duration_ms,omitempty"   example:"245"`
	ErrorDetail     string      `json:"error_detail,omitempty"`
}

// ============================================================
// Error (public)
// ============================================================

// ErrResponse — error response มาตรฐานของ MQTT module (ใช้ใน @Failure annotation)
type ErrResponse struct {
	Error string `json:"error" example:"topic is required"`
}
