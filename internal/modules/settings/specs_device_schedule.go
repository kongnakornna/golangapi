package settings

// ---- Device + Schedule specs (ported from NestJS device_list_paginate / schedule queries) ----

var deviceSelects = []string{
	"d.device_id AS device_id", "d.mqtt_id AS mqtt_id", "d.setting_id AS setting_id",
	"d.type_id AS type_id", "d.device_name AS device_name", "d.sn AS sn",
	"d.hardware_id AS hardware_id", "d.status_warning AS status_warning",
	"d.recovery_warning AS recovery_warning", "d.status_alert AS status_alert",
	"d.recovery_alert AS recovery_alert", "d.time_life AS time_life", "d.period AS period",
	"d.work_status AS work_status", "d.layout AS layout", "d.menu AS menu",
	`d.max AS "max"`, `d.min AS "min"`, "d.oid AS oid", "d.mqtt_data_value AS mqtt_data_value",
	"d.mqtt_data_control AS mqtt_data_control", "d.model AS model", "d.vendor AS vendor",
	"d.comparevalue AS comparevalue", "d.createddate AS createddate",
	"d.updateddate AS updateddate", "d.status AS status", "d.unit AS unit",
	"d.action_id AS action_id", "d.status_alert_id AS status_alert_id",
	"d.measurement AS measurement",
	"d.mqtt_control_on AS mqtt_control_on", "d.mqtt_control_off AS mqtt_control_off",
	"d.org AS device_org", "d.bucket AS device_bucket",
	"t.type_name AS type_name", "l.location_name AS location_name", "l.configdata AS configdata",
	"mq.mqtt_name AS mqtt_name", "mq.org AS mqtt_org", "mq.bucket AS mqtt_bucket",
	"mq.envavorment AS mqtt_envavorment", "mq.host AS mqtt_host", "mq.port AS mqtt_port",
	"d.mqtt_device_name AS mqtt_device_name", "d.mqtt_status_over_name AS mqtt_status_over_name",
	"d.mqtt_status_data_name AS mqtt_status_data_name", "d.mqtt_act_relay_name AS mqtt_act_relay_name",
	"d.mqtt_control_relay_name AS mqtt_control_relay_name",
	"h.host_name AS host_name", "h.port AS port", "h.host_id AS host_id",
	`CASE WHEN d.hardware_id = 1 THEN 'Sensor' WHEN d.hardware_id = 2 THEN 'IO Sensor'
	      WHEN d.hardware_id = 3 THEN 'IO Control' WHEN d.hardware_id = 4 THEN 'Critical Sensor'
	      ELSE 'Unknown' END AS hardware_type_name`,
	`CASE WHEN d.layout = 1 THEN 'Right Menu' WHEN d.layout = 2 THEN 'Card'
	      WHEN d.layout = 3 THEN 'Left Menu' WHEN d.layout = 4 THEN 'Footer Menu'
	      ELSE 'Unknown' END AS layoutapp`,
	"d.calibration_add AS calibration_add", "d.calibration_subtract AS calibration_subtract",
	"d.calibration_type AS calibration_type",
	`CASE WHEN d.calibration_type = 1 THEN 'Calibration Add' WHEN d.calibration_type = 2 THEN 'Calibration Subtract'
	      ELSE 'Non calibration' END AS calibrationtype`,
}

var DeviceSpec = ListSpec{
	Table:   "sd_iot_device d",
	Selects: deviceSelects,
	Joins: []string{
		"LEFT JOIN sd_iot_device_type t ON t.type_id = d.type_id",
		"LEFT JOIN sd_iot_mqtt mq ON mq.mqtt_id = d.mqtt_id",
		"LEFT JOIN sd_iot_location l ON l.location_id = d.location_id",
		"LEFT JOIN sd_iot_host h ON h.idhost = mq.mqtt_main_id",
	},
	Filters: []Filter{
		{Column: "d.device_name", Key: "keyword", Op: "LIKE"},
		{Column: "d.device_id", Key: "device_id"},
		{Column: "d.mqtt_id", Key: "mqtt_id"},
		{Column: "d.setting_id", Key: "setting_id"},
		{Column: "d.type_id", Key: "type_id"},
		{Column: "d.location_id", Key: "location_id"},
		{Column: "d.org", Key: "org"},
		{Column: "d.bucket", Key: "bucket"},
		{Column: "d.sn", Key: "sn"},
		{Column: "d.createddate", Key: "createddate"},
		{Column: "d.updateddate", Key: "updateddate"},
		{Column: "d.status_warning", Key: "status_warning"},
		{Column: "d.recovery_warning", Key: "recovery_warning"},
		{Column: "d.status_alert", Key: "status_alert"},
		{Column: "d.recovery_alert", Key: "recovery_alert"},
		{Column: "d.time_life", Key: "time_life"},
		{Column: "d.period", Key: "period"},
		{Column: "d.max::text", Key: "max"},
		{Column: "d.min::text", Key: "min"},
		{Column: "d.hardware_id", Key: "hardware_id"},
		{Column: "d.model", Key: "model"},
		{Column: "d.vendor", Key: "vendor"},
		{Column: "d.comparevalue", Key: "comparevalue"},
		{Column: "d.oid", Key: "oid"},
		{Column: "d.action_id", Key: "action_id"},
		{Column: "d.mqtt_data_value", Key: "mqtt_data_value"},
		{Column: "d.mqtt_data_control", Key: "mqtt_data_control"},
		{Column: "d.work_status", Key: "work_status"},
		{Column: "d.menu", Key: "menu"},
	},
	SortCols: map[string]string{
		"device_id":   "d.device_id",
		"device_name": "d.device_name",
		"sort":        "mq.sort",
		"status":      "d.status",
		"work_status": "d.work_status",
		"createddate": "d.createddate",
		"updateddate": "d.updateddate",
	},
	DefaultSort: "mq.sort ASC, d.device_id ASC",
}

// deviceCommonFilters is the WHERE set shared by the device_list_paginate*
// variants (keyword + the d.* equality filters; model is skipped because the
// NestJS condition "d.model=::model" is malformed and always throws).
var deviceCommonFilters = []Filter{
	{Column: "d.device_name", Key: "keyword", Op: "LIKE"},
	{Column: "d.device_id", Key: "device_id"},
	{Column: "d.mqtt_id", Key: "mqtt_id"},
	{Column: "d.type_id", Key: "type_id"},
	{Column: "d.location_id", Key: "location_id"},
	{Column: "d.org", Key: "org"},
	{Column: "d.bucket", Key: "bucket"},
	{Column: "d.sn", Key: "sn"},
	{Column: "d.createddate", Key: "createddate"},
	{Column: "d.updateddate", Key: "updateddate"},
	{Column: "d.status_warning", Key: "status_warning"},
	{Column: "d.recovery_warning", Key: "recovery_warning"},
	{Column: "d.status_alert", Key: "status_alert"},
	{Column: "d.recovery_alert", Key: "recovery_alert"},
	{Column: "d.time_life", Key: "time_life"},
	{Column: "d.period", Key: "period"},
	{Column: "d.max::text", Key: "max"},
	{Column: "d.min::text", Key: "min"},
	{Column: "d.hardware_id", Key: "hardware_id"},
	{Column: "d.vendor", Key: "vendor"},
	{Column: "d.comparevalue", Key: "comparevalue"},
	{Column: "d.oid", Key: "oid"},
	{Column: "d.action_id", Key: "action_id"},
	{Column: "d.mqtt_data_value", Key: "mqtt_data_value"},
	{Column: "d.mqtt_data_control", Key: "mqtt_data_control"},
}

var deviceVariantSorts = map[string]string{
	"device_id":   "d.device_id",
	"device_name": "d.device_name",
	"sort":        "mq.sort",
	"status":      "d.status",
	"work_status": "d.work_status",
	"createddate": "d.createddate",
	"updateddate": "d.updateddate",
}

// deviceActiveSelects mirrors device_list_paginate_active (listdevicepageactive,
// listdevicepageactive1): the shared columns WITHOUT the sd_iot_host fields —
// that query has no host join. Fixed WHERE: d.status=1 AND mq.status=1.
var deviceActiveSelects = []string{
	"d.device_id AS device_id", "d.mqtt_id AS mqtt_id", "d.setting_id AS setting_id",
	"d.type_id AS type_id", "d.device_name AS device_name", "d.sn AS sn",
	"d.hardware_id AS hardware_id", "d.status_warning AS status_warning",
	"d.recovery_warning AS recovery_warning", "d.status_alert AS status_alert",
	"d.recovery_alert AS recovery_alert", "d.time_life AS time_life", "d.period AS period",
	"d.work_status AS work_status", "d.layout AS layout", "d.menu AS menu",
	`d.max AS "max"`, `d.min AS "min"`, "d.oid AS oid", "d.mqtt_data_value AS mqtt_data_value",
	"d.mqtt_data_control AS mqtt_data_control", "d.model AS model", "d.vendor AS vendor",
	"d.comparevalue AS comparevalue", "d.createddate AS createddate",
	"d.updateddate AS updateddate", "d.status AS status", "d.unit AS unit",
	"d.action_id AS action_id", "d.status_alert_id AS status_alert_id",
	"d.measurement AS measurement",
	"d.mqtt_control_on AS mqtt_control_on", "d.mqtt_control_off AS mqtt_control_off",
	"d.org AS device_org", "d.bucket AS device_bucket",
	"t.type_name AS type_name", "l.location_name AS location_name", "l.configdata AS configdata",
	"mq.mqtt_name AS mqtt_name", "mq.org AS mqtt_org", "mq.bucket AS mqtt_bucket",
	"mq.envavorment AS mqtt_envavorment", "mq.host AS mqtt_host", "mq.port AS mqtt_port",
	"d.mqtt_device_name AS mqtt_device_name", "d.mqtt_status_over_name AS mqtt_status_over_name",
	"d.mqtt_status_data_name AS mqtt_status_data_name", "d.mqtt_act_relay_name AS mqtt_act_relay_name",
	"d.mqtt_control_relay_name AS mqtt_control_relay_name",
	`CASE WHEN d.hardware_id = 1 THEN 'Sensor' WHEN d.hardware_id = 2 THEN 'IO Sensor'
	      WHEN d.hardware_id = 3 THEN 'IO Control' WHEN d.hardware_id = 4 THEN 'Critical Sensor'
	      ELSE 'Unknown' END AS hardware_type_name`,
	`CASE WHEN d.layout = 1 THEN 'Right Menu' WHEN d.layout = 2 THEN 'Card'
	      WHEN d.layout = 3 THEN 'Left Menu' WHEN d.layout = 4 THEN 'Footer Menu'
	      ELSE 'Unknown' END AS layoutapp`,
	"d.calibration_add AS calibration_add", "d.calibration_subtract AS calibration_subtract",
	"d.calibration_type AS calibration_type",
	`CASE WHEN d.calibration_type = 1 THEN 'Calibration Add' WHEN d.calibration_type = 2 THEN 'Calibration Subtract'
	      ELSE 'Non calibration' END AS calibrationtype`,
}

// DeviceActiveSpec ports device_list_paginate_active.
var DeviceActiveSpec = ListSpec{
	Table:   "sd_iot_device d",
	Selects: deviceActiveSelects,
	Joins: []string{
		"LEFT JOIN sd_iot_setting st ON st.setting_id = d.setting_id",
		"LEFT JOIN sd_iot_device_type t ON t.type_id = d.type_id",
		"LEFT JOIN sd_iot_mqtt mq ON mq.mqtt_id = d.mqtt_id",
		"LEFT JOIN sd_iot_location l ON l.location_id = d.location_id",
	},
	Filters:      deviceCommonFilters,
	FixedFilters: []FixedFilter{{Column: "d.status", Value: "1"}, {Column: "mq.status", Value: "1"}},
	SortCols:     deviceVariantSorts,
	DefaultSort:  "mq.sort ASC, d.device_id ASC",
}

// DeviceActive1Spec shares the exact query of listdevicepageactive1
// (NestJS routes active and active1 both call device_list_paginate_active).
var DeviceActive1Spec = DeviceActiveSpec

// DeviceAllSpec ports device_list_paginate_all (listdevicepageall,
// listdevicepagesensor): full join set incl. sd_iot_host, extra mq.sort /
// timestamp selects, a layout filter, and NO fixed status conditions.
var DeviceAllSpec = ListSpec{
	Table:   "sd_iot_device d",
	Selects: append(append([]string{}, deviceSelects...), "mq.sort AS sort", "d.updateddate AS timestamp"),
	Joins: []string{
		"LEFT JOIN sd_iot_setting st ON st.setting_id = d.setting_id",
		"LEFT JOIN sd_iot_device_type t ON t.type_id = d.type_id",
		"LEFT JOIN sd_iot_mqtt mq ON mq.mqtt_id = d.mqtt_id",
		"LEFT JOIN sd_iot_location l ON l.location_id = d.location_id",
		"LEFT JOIN sd_iot_host h ON h.idhost = mq.mqtt_main_id",
	},
	Filters:     append(append([]Filter{}, deviceCommonFilters...), Filter{Column: "d.layout", Key: "layout"}),
	SortCols:    deviceVariantSorts,
	DefaultSort: "mq.sort ASC, d.device_id ASC",
}

// deviceAllActiveSelects mirrors device_list_paginate_all_active: like the
// active variant plus host fields and the "timestamp" alias, but without the
// raw hardware_id / layout columns. Joined with INNER JOINs in NestJS.
var deviceAllActiveSelects = []string{
	"d.device_id AS device_id", "d.mqtt_id AS mqtt_id", "d.setting_id AS setting_id",
	"d.type_id AS type_id", "d.device_name AS device_name", "d.sn AS sn",
	"d.status_warning AS status_warning",
	"d.recovery_warning AS recovery_warning", "d.status_alert AS status_alert",
	"d.recovery_alert AS recovery_alert", "d.time_life AS time_life", "d.period AS period",
	"d.work_status AS work_status", "d.menu AS menu",
	`d.max AS "max"`, `d.min AS "min"`, "d.oid AS oid", "d.mqtt_data_value AS mqtt_data_value",
	"d.mqtt_data_control AS mqtt_data_control", "d.model AS model", "d.vendor AS vendor",
	"d.comparevalue AS comparevalue", "d.createddate AS createddate",
	"d.updateddate AS updateddate", "d.status AS status", "d.unit AS unit",
	"d.action_id AS action_id", "d.status_alert_id AS status_alert_id",
	"d.measurement AS measurement",
	"d.mqtt_control_on AS mqtt_control_on", "d.mqtt_control_off AS mqtt_control_off",
	"d.org AS device_org", "d.bucket AS device_bucket", "d.updateddate AS timestamp",
	"d.mqtt_device_name AS mqtt_device_name", "d.mqtt_status_over_name AS mqtt_status_over_name",
	"d.mqtt_status_data_name AS mqtt_status_data_name", "d.mqtt_act_relay_name AS mqtt_act_relay_name",
	"d.mqtt_control_relay_name AS mqtt_control_relay_name",
	"t.type_name AS type_name", "l.location_name AS location_name", "l.configdata AS configdata",
	"mq.mqtt_name AS mqtt_name", "mq.org AS mqtt_org", "mq.bucket AS mqtt_bucket",
	"mq.envavorment AS mqtt_envavorment", "mq.host AS mqtt_host", "mq.port AS mqtt_port",
	"h.host_name AS host_name", "h.port AS port", "h.host_id AS host_id",
	`CASE WHEN d.hardware_id = 1 THEN 'Sensor' WHEN d.hardware_id = 2 THEN 'IO Sensor'
	      WHEN d.hardware_id = 3 THEN 'IO Control' WHEN d.hardware_id = 4 THEN 'Critical Sensor'
	      ELSE 'Unknown' END AS hardware_type_name`,
	`CASE WHEN d.layout = 1 THEN 'Right Menu' WHEN d.layout = 2 THEN 'Card'
	      WHEN d.layout = 3 THEN 'Left Menu' WHEN d.layout = 4 THEN 'Footer Menu'
	      ELSE 'Unknown' END AS layoutapp`,
	"d.calibration_add AS calibration_add", "d.calibration_subtract AS calibration_subtract",
	"d.calibration_type AS calibration_type",
	`CASE WHEN d.calibration_type = 1 THEN 'Calibration Add' WHEN d.calibration_type = 2 THEN 'Calibration Subtract'
	      ELSE 'Non calibration' END AS calibrationtype`,
}

// DeviceAllActiveSpec ports device_list_paginate_all_active
// (listdevicepageallactive): INNER JOINs across all five tables and the fixed
// d.status=1 AND mq.status=1 conditions.
var DeviceAllActiveSpec = ListSpec{
	Table:   "sd_iot_device d",
	Selects: deviceAllActiveSelects,
	Joins: []string{
		"INNER JOIN sd_iot_mqtt mq ON mq.mqtt_id = d.mqtt_id",
		"INNER JOIN sd_iot_setting st ON st.setting_id = d.setting_id",
		"INNER JOIN sd_iot_device_type t ON t.type_id = d.type_id",
		"INNER JOIN sd_iot_location l ON l.location_id = d.location_id",
		"INNER JOIN sd_iot_host h ON h.idhost = mq.mqtt_main_id",
	},
	Filters:      deviceCommonFilters,
	FixedFilters: []FixedFilter{{Column: "d.status", Value: "1"}, {Column: "mq.status", Value: "1"}},
	SortCols:     deviceVariantSorts,
	DefaultSort:  "mq.sort ASC, d.device_id ASC",
}

// DeviceAllActiveSchedSpec is the base query of
// listdevicepageallactiveschedule; the handler adds the schedule merge logic.
var DeviceAllActiveSchedSpec = DeviceAllActiveSpec

var ScheduleSpec = ListSpec{
	Table: "sd_iot_schedule sc",
	Selects: []string{
		"sc.schedule_id AS schedule_id", "sc.schedule_name AS schedule_name",
		"sc.device_id AS device_id", "sc.start AS start", "sc.event AS event",
		"sc.sunday AS sunday", "sc.monday AS monday", "sc.tuesday AS tuesday",
		"sc.wednesday AS wednesday", "sc.thursday AS thursday", "sc.friday AS friday",
		"sc.saturday AS saturday", "sc.status AS status",
		"sc.createddate AS createddate", "sc.updateddate AS updateddate",
		"d.device_name AS device_name",
		"(SELECT COUNT(DISTINCT device_id) FROM sd_iot_schedule_device WHERE schedule_id = sc.schedule_id) AS \"countDevice\"",
	},
	Joins: []string{"LEFT JOIN sd_iot_device d ON d.device_id = sc.device_id"},
	Filters: []Filter{
		{Column: "sc.schedule_name", Key: "keyword", Op: "LIKE"},
		{Column: "sc.schedule_id", Key: "schedule_id"},
		{Column: "sc.device_id", Key: "device_id"},
		{Column: "sc.event", Key: "event"},
		{Column: "sc.monday", Key: "monday"},
		{Column: "sc.tuesday", Key: "tuesday"},
		{Column: "sc.wednesday", Key: "wednesday"},
		{Column: "sc.thursday", Key: "thursday"},
		{Column: "sc.friday", Key: "friday"},
		{Column: "sc.saturday", Key: "saturday"},
		{Column: "sc.sunday", Key: "sunday"},
		{Column: "sc.start", Key: "start"},
		{Column: "sc.createddate", Key: "createddate"},
		{Column: "sc.updateddate", Key: "updateddate"},
		{Column: "sc.status", Key: "status"},
	},
	SortCols: map[string]string{
		"schedule_id":   "sc.schedule_id",
		"schedule_name": "sc.schedule_name",
		"start":         "sc.start",
		"createddate":   "sc.createddate",
		"updateddate":   "sc.updateddate",
	},
	DefaultSort: "sc.start ASC",
}

// FullScheduleSpec ports the fs_schedule table into the same display shape as
// ScheduleSpec so /settings/listschedulepage can merge both together.
// ตัวสเปก fs_schedule — จัดรูปรายการให้เหมือนกับ sd_iot_schedule
var FullScheduleSpec = ListSpec{
	Table: "fs_schedule",
	Selects: []string{
		"id::text AS schedule_id", "name AS schedule_name",
		"time_start AS start", "event_action AS event",
		"sunday AS sunday", "monday AS monday", "tuesday AS tuesday",
		"wednesday AS wednesday", "thursday AS thursday", "friday AS friday",
		"saturday AS saturday",
		"status AS status",
		"created_at AS createddate", "updated_at AS updateddate",
		"(SELECT COUNT(*) FROM fs_schedule_device fsd WHERE fsd.schedule_id = fs_schedule.id) AS \"countDevice\"",
	},
	Filters: []Filter{
		{Column: "name", Key: "keyword", Op: "LIKE"},
		{Column: "id::text", Key: "schedule_id"},
		{Column: "time_start", Key: "start"},
		{Column: "event_action", Key: "event"},
		{Column: "created_at", Key: "createddate"},
		{Column: "updated_at", Key: "updateddate"},
	},
	// status in fs_schedule is int8 (1=active, 0=inactive, 3=draft); the settings
	// page expects 1/0, so map the query values to the int8 column.
	ParamConds: []ParamCond{
		{
			Key: "status",
			Conds: map[string]string{
				"1": "status = 1",
				"0": "status = 0",
				"3": "status = 3",
			},
		},
	},
	// Value == "" signals "deleted_at IS NULL" (empty Value means NULL-check,
	// see SettingsPgRepo.listBase).
	FixedFilters: []FixedFilter{{Column: "deleted_at"}},
	SortCols: map[string]string{
		"schedule_id":   "id",
		"schedule_name": "name",
		"start":         "time_start",
		"createddate":   "created_at",
		"updateddate":   "updated_at",
	},
	DefaultSort: "time_start ASC",
}

var ScheduleDeviceJoinSpec = ListSpec{
	Table: "sd_iot_schedule_device s",
	Selects: []string{
		"s.id AS id", "s.schedule_id AS schedule_id", "s.device_id AS device_id",
		"sc.schedule_name AS schedule_name", "sc.start AS start", "sc.event AS event",
		"sc.sunday AS sunday", "sc.monday AS monday", "sc.tuesday AS tuesday",
		"sc.wednesday AS wednesday", "sc.thursday AS thursday", "sc.friday AS friday",
		"sc.saturday AS saturday", "sc.status AS schedule_status",
		"d.device_name AS device_name", "d.sn AS sn", "d.mqtt_data_value AS mqtt_data_value",
		"d.mqtt_control_on AS mqtt_control_on", "d.mqtt_control_off AS mqtt_control_off",
	},
	Joins: []string{
		"LEFT JOIN sd_iot_schedule sc ON sc.schedule_id = s.schedule_id",
		"LEFT JOIN sd_iot_device d ON d.device_id = s.device_id",
	},
	Filters: []Filter{
		{Column: "s.schedule_id", Key: "schedule_id"},
		{Column: "s.device_id", Key: "device_id"},
		{Column: "sc.monday", Key: "monday"},
		{Column: "sc.tuesday", Key: "tuesday"},
		{Column: "sc.wednesday", Key: "wednesday"},
		{Column: "sc.thursday", Key: "thursday"},
		{Column: "sc.friday", Key: "friday"},
		{Column: "sc.saturday", Key: "saturday"},
		{Column: "sc.sunday", Key: "sunday"},
		{Column: "sc.start", Key: "start"},
		{Column: "sc.event", Key: "event"},
	},
	SortCols:    map[string]string{"schedule_id": "s.schedule_id", "device_id": "s.device_id"},
	DefaultSort: "s.schedule_id ASC",
}
