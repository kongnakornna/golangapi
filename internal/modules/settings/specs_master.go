package settings

// ---- Master resource specs (ported from NestJS *_list_paginate QueryBuilders) ----

var SettingSpec = ListSpec{
	Table: "sd_iot_setting s",
	Selects: []string{
		"s.setting_id AS setting_id", "s.location_id AS location_id",
		"s.setting_type_id AS setting_type_id", "s.setting_name AS setting_name", "s.sn AS sn",
		"l.location_name AS location_name", "t.type_name AS type_name",
		"s.status AS status", "s.createddate AS createddate", "s.updateddate AS updateddate",
	},
	Joins: []string{
		"LEFT JOIN sd_iot_location l ON l.location_id = s.location_id",
		"LEFT JOIN sd_iot_type t ON t.type_id = s.setting_type_id",
	},
	Filters: []Filter{
		{Column: "s.setting_name", Key: "keyword", Op: "LIKE"},
		{Column: "s.setting_id", Key: "setting_id"},
		{Column: "s.location_id", Key: "location_id"},
		{Column: "s.setting_type_id", Key: "setting_type_id"},
		{Column: "s.sn", Key: "sn"},
		// NestJS applies the status filter to BOTH tables when provided.
		{Column: "l.status", Key: "status"},
		{Column: "s.status", Key: "status"},
	},
	SortCols: map[string]string{
		"setting_id":   "s.setting_id",
		"setting_name": "s.setting_name",
		"createddate":  "s.createddate",
	},
	DefaultSort: "s.createddate ASC",
}

var LocationSpec = ListSpec{
	Table: "sd_iot_location l",
	Selects: []string{
		"l.location_id AS location_id", "l.location_name AS location_name",
		"l.ipaddress AS ipaddress", "l.location_detail AS location_detail",
		"l.configdata AS configdata", "l.status AS status",
		"l.createddate AS createddate", "l.updateddate AS updateddate",
	},
	Filters: []Filter{
		{Column: "l.location_name", Key: "keyword", Op: "LIKE"},
		{Column: "l.location_id", Key: "location_id"},
		{Column: "l.ipaddress", Key: "ipaddress"},
		{Column: "l.createddate", Key: "createddate"},
		{Column: "l.updateddate", Key: "updateddate"},
		{Column: "l.status", Key: "status"},
	},
	SortCols: map[string]string{
		"location_id":   "l.location_id",
		"location_name": "l.location_name",
		"createddate":   "l.createddate",
	},
	DefaultSort: "l.location_id ASC",
}

var TypeSpec = ListSpec{
	Table: "sd_iot_type t",
	Selects: []string{
		"t.type_id AS type_id", "t.type_name AS type_name", "t.group_id AS group_id",
		"g.group_name AS group_name", "t.status AS status",
		"t.createddate AS createddate", "t.updateddate AS updateddate",
	},
	Joins: []string{"LEFT JOIN sd_iot_group g ON g.group_id = t.group_id"},
	Filters: []Filter{
		{Column: "t.type_name", Key: "keyword", Op: "LIKE"},
		{Column: "t.type_id", Key: "type_id"},
		{Column: "t.group_id", Key: "group_id"},
		{Column: "t.createddate", Key: "createddate"},
		{Column: "t.updateddate", Key: "updateddate"},
		{Column: "t.status", Key: "status"},
	},
	SortCols: map[string]string{
		"type_id":     "t.type_id",
		"type_name":   "t.type_name",
		"createddate": "t.createddate",
	},
	DefaultSort: "t.type_id ASC",
}

var DeviceTypeSpec = ListSpec{
	Table: "sd_iot_device_type d",
	Selects: []string{
		"d.type_id AS type_id", "d.type_name AS type_name", "d.status AS status",
		"d.createddate AS createddate", "d.updateddate AS updateddate",
	},
	Filters: []Filter{
		{Column: "d.type_name", Key: "keyword", Op: "LIKE"},
		{Column: "d.type_id", Key: "type_id"},
		{Column: "d.status", Key: "status"},
	},
	SortCols: map[string]string{
		"type_id":     "d.type_id",
		"type_name":   "d.type_name",
		"createddate": "d.createddate",
	},
	DefaultSort: "d.type_id ASC",
}

var GroupSpec = ListSpec{
	Table: "sd_iot_group g",
	Selects: []string{
		"g.group_id AS group_id", "g.group_name AS group_name", "g.status AS status",
		"g.createddate AS createddate", "g.updateddate AS updateddate",
	},
	Filters: []Filter{
		{Column: "g.group_name", Key: "keyword", Op: "LIKE"},
		{Column: "g.group_id", Key: "group_id"},
		{Column: "g.createddate", Key: "createddate"},
		{Column: "g.updateddate", Key: "updateddate"},
		{Column: "g.status", Key: "status"},
	},
	SortCols: map[string]string{
		"group_id":    "g.group_id",
		"group_name":  "g.group_name",
		"createddate": "g.createddate",
	},
	DefaultSort: "g.group_id ASC",
}

// SensorSpec ports sensor_list_paginate (sensor query semantics — see plan deviations).
var SensorSpec = ListSpec{
	Table: "sd_iot_sensor s",
	Selects: []string{
		"s.sensor_id AS sensor_id", "s.setting_id AS setting_id", "s.setting_type_id AS setting_type_id",
		"s.sensor_name AS sensor_name", "s.sn AS sn", "s.max AS max", "s.min AS min",
		"s.hardware_id AS hardware_id", "s.status_high AS status_high", "s.status_warning AS status_warning",
		"s.status_alert AS status_alert", "s.model AS model", "s.vendor AS vendor",
		"s.comparevalue AS comparevalue", "s.unit AS unit", "s.mqtt_id AS mqtt_id", "s.oid AS oid",
		"s.action_id AS action_id", "s.status_alert_id AS status_alert_id",
		"s.mqtt_data_value AS mqtt_data_value", "s.mqtt_data_control AS mqtt_data_control",
		"s.status AS status", "s.createddate AS createddate", "s.updateddate AS updateddate",
		"st.setting_name AS setting_name", "t.type_name AS type_name", "l.location_name AS location_name",
	},
	Joins: []string{
		"LEFT JOIN sd_iot_setting st ON st.setting_id = s.setting_id",
		"LEFT JOIN sd_iot_type t ON t.type_id = s.setting_type_id",
		"LEFT JOIN sd_iot_location l ON l.location_id = st.location_id",
	},
	Filters: []Filter{
		{Column: "s.sensor_name", Key: "keyword", Op: "LIKE"},
		{Column: "s.sensor_id", Key: "sensor_id"},
		{Column: "s.setting_id", Key: "setting_id"},
		{Column: "s.setting_type_id", Key: "setting_type_id"},
		{Column: "s.sn", Key: "sn"},
		{Column: "s.hardware_id", Key: "hardware_id"},
		{Column: "s.mqtt_id", Key: "mqtt_id"},
		{Column: "s.model", Key: "model"},
		{Column: "s.vendor", Key: "vendor"},
		{Column: "s.status", Key: "status"},
	},
	SortCols: map[string]string{
		"sensor_id":   "s.sensor_id",
		"sensor_name": "s.sensor_name",
		"createddate": "s.createddate",
	},
	DefaultSort: "s.sensor_id ASC",
}
