package settings

// ---- Alarm link + process log specs ----
// Ported from the NestJS alarmlogpaginate* QueryBuilders: every variant joins
// device / device_type / mqtt / location / alarm-action and re-aliases the
// daa.* columns (warning, recoverywarning, alert, ..., event).

var alarmLogSelects = []string{
	"al.id AS id", "al.alarm_action_id AS alarm_action_id", "al.event AS event_log",
	"al.alarm_type AS alarm_type", "al.status_warning AS status_warning",
	"al.recovery_warning AS recovery_warning", "al.status_alert AS status_alert",
	"al.recovery_alert AS recovery_alert", "al.email_alarm AS email_alarm",
	"al.line_alarm AS line_alarm", "al.telegram_alarm AS telegram_alarm",
	"al.sms_alarm AS sms_alarm", "al.nonc_alarm AS nonc_alarm",
	"al.status AS status", "al.date AS date", "al.time AS time",
	"al.data AS data", "al.data_alarm AS data_alarm", "al.alarm_status AS alarm_status",
	"al.updateddate AS updateddate", "al.subject AS subject", "al.content AS content",
	"al.createddate AS createddate",
	"d.device_id AS device_id", "d.bucket AS bucket", "d.mqtt_id AS mqtt_id",
	"d.setting_id AS setting_id", "d.type_id AS type_id", "d.device_name AS device_name",
	"d.mqtt_device_name AS mqtt_device_name", "t.type_name AS type_name",
	"l.location_name AS location_name",
	"mq.mqtt_name AS mqtt_name", "mq.org AS mqtt_org", "mq.bucket AS mqtt_bucket",
	"mq.envavorment AS mqtt_envavorment",
	"daa.action_name AS action_name", "daa.status_warning AS warning",
	"daa.recovery_warning AS recoverywarning", "daa.status_alert AS alert",
	"daa.recovery_alert AS recoveryalert", "daa.email_alarm AS emailalarm",
	"daa.line_alarm AS linealarm", "daa.telegram_alarm AS telegramalarm",
	"daa.sms_alarm AS smsalarm", "daa.nonc_alarm AS noncalarm",
	"daa.time_life AS timelife", "daa.event AS event",
}

// typeIdLogConds ports the NestJS switch over dto.type_id_log.
var typeIdLogConds = map[string]string{
	"1": "al.email_alarm = 1",
	"2": "al.line_alarm = 1",
	"3": "al.telegram_alarm = 1",
	"4": "al.sms_alarm = 1",
	"5": "al.nonc_alarm = 1",
}

func alarmLogFilters(dateCol string) []Filter {
	return []Filter{
		{Column: "daa.action_name", Key: "keyword", Op: "LIKE"},
		{Column: "daa.event", Key: "event"},
		{Column: "al.alarm_action_id", Key: "alarm_action_id"},
		{Column: "d.device_id", Key: "device_id"},
		{Column: "daa.status", Key: "status"},
		{Column: "d.type_id", Key: "type_id"},
		{Column: "l.location_id", Key: "location_id"},
		{Column: "d.bucket", Key: "bucket"},
		{Column: dateCol, Key: "start", Op: ">="},
		{Column: dateCol, Key: "end", Op: "<="},
	}
}

func alarmLogJoins(inner bool) []string {
	kw := "LEFT"
	if inner {
		kw = "INNER"
	}
	return []string{
		kw + " JOIN sd_iot_device_alarm_action daa ON daa.alarm_action_id = al.alarm_action_id",
		kw + " JOIN sd_iot_device d ON d.device_id = al.device_id",
		kw + " JOIN sd_iot_device_type t ON t.type_id = d.type_id",
		kw + " JOIN sd_iot_mqtt mq ON mq.mqtt_id = d.mqtt_id",
		kw + " JOIN sd_iot_location l ON l.location_id = mq.location_id",
	}
}

var alarmLogSortCols = map[string]string{
	"id":              "al.id",
	"alarm_action_id": "al.alarm_action_id",
	"date":            "al.date",
	"time":            "al.time",
	"status":          "al.status",
	"alarm_status":    "al.alarm_status",
	"createddate":     "al.createddate",
	"updateddate":     "al.updateddate",
}

// AlarmProcessLogSpec ports alarmlogpaginate (table sd_alarm_process_log_temp):
// LEFT JOINs, start/end range on al.date, fixed order al.date DESC, al.time DESC.
var AlarmProcessLogSpec = ListSpec{
	Table:   "sd_alarm_process_log_temp al",
	Selects: alarmLogSelects,
	Joins:   alarmLogJoins(false),
	Filters: alarmLogFilters("al.date"),
	ParamConds: []ParamCond{
		{Key: "type_id_log", Conds: typeIdLogConds},
	},
	SortCols:    alarmLogSortCols,
	DefaultSort: "al.date DESC, al.time DESC",
}

// alarmLogVariantSpec ports the email / line / sms / telegram / control
// variants: same shape but INNER JOINs, start/end range on al.createddate and
// default order createddate DESC, updateddate DESC (sortable in NestJS).
func alarmLogVariantSpec(table string) ListSpec {
	return ListSpec{
		Table:   table + " al",
		Selects: alarmLogSelects,
		Joins:   alarmLogJoins(true),
		Filters: alarmLogFilters("al.createddate"),
		ParamConds: []ParamCond{
			{Key: "type_id_log", Conds: typeIdLogConds},
		},
		SortCols:    alarmLogSortCols,
		DefaultSort: "al.createddate DESC, al.updateddate DESC",
	}
}

var (
	AlarmProcessLogEmailSpec    = alarmLogVariantSpec("sd_alarm_process_log_email")
	AlarmProcessLogLineSpec     = alarmLogVariantSpec("sd_alarm_process_log_line")
	AlarmProcessLogSmsSpec      = alarmLogVariantSpec("sd_alarm_process_log_sms")
	AlarmProcessLogTelegramSpec = alarmLogVariantSpec("sd_alarm_process_log_telegram")
	// AlarmProcessLogControlSpec ports alarmlogpaginateecontrol
	// (/alarmlogpaginatecontrol[s]): sd_alarm_process_log with INNER JOINs.
	AlarmProcessLogControlSpec = alarmLogVariantSpec("sd_alarm_process_log")
)

// ScheduleProcessLogSpec ports scheduleprocesslog_paginate
// (sd_schedule_process_log sl): the sl.* columns plus the joined device block.
// The NestJS sort stage references a nonexistent "alarm" alias, so ordering is
// broken there; Go keeps a deterministic default instead. Several NestJS
// filters reference a missing "st" alias or l.schedule_id/l.device_id (broken)
// — their intent is preserved via the working sl.* equality filters below.
var ScheduleProcessLogSpec = ListSpec{
	Table: "sd_schedule_process_log sl",
	Selects: append([]string{
		"sl.id AS id", "sl.schedule_id AS schedule_id",
		"sl.schedule_event_start AS schedule_event_start", "sl.day AS day", "sl.doday AS doday",
		"sl.dotime AS dotime", "sl.schedule_event AS schedule_event", "sl.device_status AS device_status",
	}, deviceActiveSelects...),
	Joins: []string{
		"LEFT JOIN sd_iot_device d ON d.device_id = sl.device_id",
		"LEFT JOIN sd_iot_device_type t ON t.type_id = d.type_id",
		"LEFT JOIN sd_iot_mqtt mq ON mq.mqtt_id = d.mqtt_id",
		"LEFT JOIN sd_iot_location l ON l.location_id = mq.location_id",
	},
	Filters: []Filter{
		{Column: "d.device_name", Key: "keyword", Op: "LIKE"},
		{Column: "d.org", Key: "org"},
		{Column: "d.bucket", Key: "bucket"},
		{Column: "d.type_id", Key: "type_id"},
		{Column: "d.status_warning", Key: "status_warning"},
		{Column: "d.recovery_warning", Key: "recovery_warning"},
		{Column: "d.status_alert", Key: "status_alert"},
		{Column: "d.recovery_alert", Key: "recovery_alert"},
		{Column: "d.time_life", Key: "time_life"},
		{Column: "d.period", Key: "period"},
		// intent-parity for the broken-in-NestJS l.* filters:
		{Column: "sl.schedule_id", Key: "schedule_id"},
		{Column: "sl.device_id", Key: "device_id"},
		{Column: "sl.day", Key: "day"},
		{Column: "sl.doday", Key: "doday"},
		{Column: "sl.dotime", Key: "dotime"},
		{Column: "sl.schedule_event", Key: "schedule_event"},
		{Column: "sl.device_status", Key: "device_status"},
		{Column: "sl.status", Key: "status"},
	},
	SortCols: map[string]string{
		"id":          "sl.id",
		"schedule_id": "sl.schedule_id",
		"day":         "sl.day",
		"doday":       "sl.doday",
		"device_name": "d.device_name",
		"createddate": "d.createddate",
	},
	DefaultSort: "sl.id ASC",
}

// MqttErrorLogSpec ports mqttlogpaginate (/mqtterrorlogpaginate): the MQTT
// alarm log table with the same INNER JOIN shape as the other log variants.
var MqttErrorLogSpec = alarmLogVariantSpec("sd_alarm_process_log_mqtt")

// alarmLinkSelects is shared by sd_iot_alarm_device / sd_iot_alarm_device_event paginates.
var alarmLinkSelects = []string{
	"a.id AS id", "a.alarm_action_id AS alarm_action_id", "a.device_id AS device_id",
	"aa.action_name AS action_name", "aa.status_warning AS status_warning",
	"aa.recovery_warning AS recovery_warning", "aa.status_alert AS status_alert",
	"aa.recovery_alert AS recovery_alert", "aa.email_alarm AS email_alarm",
	"aa.line_alarm AS line_alarm", "aa.telegram_alarm AS telegram_alarm",
	"aa.sms_alarm AS sms_alarm", "aa.nonc_alarm AS nonc_alarm",
	"aa.time_life AS time_life", "aa.event AS event", "aa.status AS action_status",
	"d.device_name AS device_name", "d.sn AS sn", "d.mqtt_data_value AS mqtt_data_value",
}

func alarmLinkSpec(linkTable string, joinAlias string) ListSpec {
	return ListSpec{
		Table:   linkTable + " a",
		Selects: alarmLinkSelects,
		Joins: []string{
			"LEFT JOIN sd_iot_device_alarm_action aa ON aa.alarm_action_id = a.alarm_action_id",
			"LEFT JOIN sd_iot_device d ON d.device_id = a.device_id",
		},
		Filters: []Filter{
			{Column: "a.device_id", Key: "device_id"},
			{Column: "a.alarm_action_id", Key: "alarm_action_id"},
			{Column: "d.device_name", Key: "keyword", Op: "LIKE"},
			{Column: joinAlias + ".status", Key: "action_status"},
		},
		SortCols: map[string]string{
			"device_id":   "a.device_id",
			"device_name": "d.device_name",
		},
		DefaultSort: "a.device_id ASC",
	}
}

var (
	AlarmDeviceJoinSpec      = alarmLinkSpec("sd_iot_alarm_device", "aa")
	AlarmEventDeviceJoinSpec = alarmLinkSpec("sd_iot_alarm_device_event", "aa")
)
