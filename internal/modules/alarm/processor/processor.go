package alarm

import (
	"fmt"
	"strings"
)

type AlarmDetailDto struct {
	HardwareID      int         `json:"hardware_id"`
	ValueData       interface{} `json:"value_data"`
	ValueAlarm      interface{} `json:"value_alarm"`
	Max             *float64    `json:"max"`
	Min             *float64    `json:"min"`
	StatusAlert     float64     `json:"status_alert"`
	StatusWarning   float64     `json:"status_warning"`
	RecoveryWarning float64     `json:"recovery_warning"`
	RecoveryAlert   float64     `json:"recovery_alert"`
	CountAlarm      float64     `json:"count_alarm"`
	Event           int         `json:"event"`
	Unit            string      `json:"unit"`
	DeviceName      string      `json:"device_name"`
	ActionName      string      `json:"action_name"`
	MqttControlOn   string      `json:"mqtt_control_on"`
	MqttControlOff  string      `json:"mqtt_control_off"`
}

type AlarmResult struct {
	Status             int     `json:"status"` // 1 warning,2 critical,3 recovery warning,4 recovery critical,5 normal
	Title              string  `json:"title"`
	Subject            string  `json:"subject"`
	Content            string  `json:"content"`
	DataAlarm          float64 `json:"data_alarm"`
	EventControl       int     `json:"event_control"`
	MessageMqttControl string  `json:"message_mqtt_control"`
}

func AlarmDetailValidate(dto AlarmDetailDto, lang string) AlarmResult {
	// แปลงภาษา (lang == "th" หรือ "en")
	messages := getMessages(lang)

	// parse numeric values
	sensorValue := normalizeSensorValue(dto.ValueData)
	maxVal := dto.Max
	minVal := dto.Min
	statusAlert := dto.StatusAlert
	statusWarning := dto.StatusWarning
	recoveryWarning := dto.RecoveryWarning
	recoveryAlert := dto.RecoveryAlert
	countAlarm := dto.CountAlarm
	event := dto.Event

	// กำหนดค่าเริ่มต้น
	// var alarmStatusSet int
	var title, subject, content string
	var dataAlarm float64
	eventControl := event
	messageMqttControl := ""
	status := 5 // default normal

	// ใช้ switch ตาม hardware_id และเงื่อนไขอื่นๆ (ย่อมาจาก TypeScript)
	switch dto.HardwareID {
	case 3: // IO Control
		if sensorValue == 1 || sensorValue == "ON" {
			status = 5
			title = messages["normal"]
			subject = messages["normal"]
			content = fmt.Sprintf("%s %v %s", messages["normal"], sensorValue, dto.Unit)
		}
	case 4: // Critical Sensor
		if sensorValue != 1 {
			status = 2
			title = messages["critical"]
			subject = fmt.Sprintf("%s %s: %v %s", dto.DeviceName, messages["critical"], sensorValue, dto.Unit)
			content = subject
			dataAlarm = statusWarning
		} else {
			status = 5
			title = messages["normal"]
		}
	default: // Sensor & IO Sensor
		if maxVal != nil && toFloat(sensorValue) >= *maxVal {
			status = 2
			title = messages["critical_max"]
			subject = fmt.Sprintf("%s %s: %v %s", dto.DeviceName, messages["critical_max"], sensorValue, dto.Unit)
			content = subject
			dataAlarm = statusWarning
		} else if minVal != nil && toFloat(sensorValue) <= *minVal {
			status = 1
			title = messages["critical_min"]
			subject = fmt.Sprintf("%s %s: %v %s", dto.DeviceName, messages["critical_min"], sensorValue, dto.Unit)
			content = subject
			dataAlarm = statusWarning
		} else if dto.HardwareID == 1 && statusWarning > 0 && toFloat(sensorValue) >= statusWarning && toFloat(sensorValue) < statusAlert {
			status = 1
			title = messages["warning"]
			subject = fmt.Sprintf("%s %s: %v %s", dto.DeviceName, messages["warning"], sensorValue, dto.Unit)
			content = subject
			dataAlarm = statusWarning
		} else if dto.HardwareID == 1 && statusAlert > 0 && toFloat(sensorValue) >= statusAlert {
			status = 2
			title = messages["critical"]
			subject = fmt.Sprintf("%s %s: %v %s", dto.DeviceName, messages["critical"], sensorValue, dto.Unit)
			content = subject
			dataAlarm = statusAlert
		} else if countAlarm >= 1 && recoveryWarning > 0 && toFloat(sensorValue) <= recoveryWarning {
			status = 3
			title = messages["recovery_warning"]
			subject = fmt.Sprintf("%s %s: %v %s", dto.DeviceName, messages["recovery_warning"], sensorValue, dto.Unit)
			content = subject
			dataAlarm = recoveryWarning
			eventControl = 1 - event
			messageMqttControl = map[bool]string{true: dto.MqttControlOff, false: dto.MqttControlOn}[event == 1]
		} else if countAlarm >= 1 && recoveryAlert > 0 && toFloat(sensorValue) <= recoveryAlert {
			status = 4
			title = messages["recovery_critical"]
			subject = fmt.Sprintf("%s %s: %v %s", dto.DeviceName, messages["recovery_critical"], sensorValue, dto.Unit)
			content = subject
			dataAlarm = recoveryAlert
			eventControl = 1 - event
			messageMqttControl = map[bool]string{true: dto.MqttControlOff, false: dto.MqttControlOn}[event == 1]
		} else {
			status = 5
			title = messages["normal"]
			subject = messages["normal"]
			content = messages["normal"]
			dataAlarm = 0
		}
	}

	return AlarmResult{
		Status:             status,
		Title:              title,
		Subject:            subject,
		Content:            content,
		DataAlarm:          dataAlarm,
		EventControl:       eventControl,
		MessageMqttControl: messageMqttControl,
	}
}

// helper functions
func getMessages(lang string) map[string]string {
	if lang == "th" {
		return map[string]string{
			"warning":           "คำเตือน มีความผิดปกติ",
			"critical":          "ภาวะวิกฤตต้องแก้ไขทันที",
			"recovery_warning":  "คืนสู่ภาวะปกติ (คำเตือน)",
			"recovery_critical": "คืนสู่ภาวะปกติ (วิกฤต)",
			"normal":            "ปกติ",
			"critical_max":      "วิกฤต มีค่าสูงเกินกำหนด",
			"critical_min":      "วิกฤต มีค่าต่ำกว่ากำหนด",
		}
	}
	return map[string]string{
		"warning":           "Warning",
		"critical":          "Critical",
		"recovery_warning":  "Recovery Warning",
		"recovery_critical": "Recovery Critical",
		"normal":            "Normal",
		"critical_max":      "Critical! Maximum limit.",
		"critical_min":      "Critical! Minimum limit",
	}
}

func normalizeSensorValue(v interface{}) interface{} {
	switch val := v.(type) {
	case string:
		if strings.ToUpper(val) == "ON" {
			return 1
		}
		if strings.ToUpper(val) == "OFF" {
			return 0
		}
		// try parse float
		var f float64
		if _, err := fmt.Sscanf(val, "%f", &f); err == nil {
			return f
		}
		return val
	default:
		return val
	}
}

func toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		var f float64
		fmt.Sscanf(val, "%f", &f)
		return f
	default:
		return 0
	}
}
