package settings

// ---- Integration resource specs ----
// Aligned with gistdaapi settings.service.ts (columns / filters / default sorts).

var MqttSpec = ListSpec{
	Table: "sd_iot_mqtt m",
	Selects: []string{
		"m.mqtt_id AS mqtt_id", "l.location_id AS location_id", "h.idhost AS idhost",
		"m.sort AS sort", "m.mqtt_type_id AS mqtt_type_id", "m.mqtt_name AS mqtt_name",
		"m.host AS host", "m.port AS port", "m.username AS username", "m.password AS password",
		"m.secret AS secret", "m.expire_in AS expire_in", "m.token_value AS token_value",
		"m.org AS org", "m.bucket AS bucket", "m.envavorment AS envavorment",
		"m.updateddate AS updateddate", "m.status AS status",
		"m.latitude AS latitude", "m.longitude AS longitude", "m.zoom AS zoom",
		"m.mqtt_main_id AS mqtt_main_id", "m.configuration AS configuration",
		"t.type_name AS type_name", "l.location_name AS location_name",
		"l.ipaddress AS ipaddress", "l.location_detail AS location_detail",
	},
	Joins: []string{
		"LEFT JOIN sd_iot_type t ON t.type_id = m.mqtt_type_id",
		"LEFT JOIN sd_iot_host h ON m.mqtt_main_id = h.idhost",
		"LEFT JOIN sd_iot_location l ON l.location_id = m.location_id",
	},
	Filters: []Filter{
		{Column: "m.mqtt_name", Key: "keyword", Op: "LIKE"},
		{Column: "m.mqtt_id", Key: "mqtt_id"},
		{Column: "m.mqtt_main_id", Key: "idhost"},
		{Column: "m.location_id", Key: "location_id"},
		{Column: "m.mqtt_type_id", Key: "mqtt_type_id"},
		{Column: "m.createddate", Key: "createddate"},
		{Column: "m.secret", Key: "secret"},
		{Column: "m.expire_in", Key: "expire_in"},
		{Column: "m.token_value", Key: "token_value"},
		{Column: "m.org", Key: "org"},
		{Column: "m.bucket", Key: "bucket"},
		{Column: "m.envavorment", Key: "envavorment"},
		{Column: "m.updateddate", Key: "updateddate"},
		{Column: "m.status", Key: "status"},
	},
	SortCols: map[string]string{
		"mqtt_id":     "m.mqtt_id",
		"sort":        "m.sort",
		"mqtt_name":   "m.mqtt_name",
		"createddate": "m.createddate",
		"updateddate": "m.updateddate",
		"status":      "m.status",
	},
	DefaultSort: "m.sort ASC, m.mqtt_id ASC",
}

var MqttHostSpec = ListSpec{
	Table: "sd_mqtt_host mh",
	Selects: []string{
		"mh.id AS id", "mh.hostname AS hostname", "mh.host AS host", "mh.port AS port",
		"mh.username AS username", "mh.password AS password", "mh.updateddate AS updateddate",
		"mh.status AS status", "mh.idhost AS idhost",
	},
	Filters: []Filter{
		{Column: "mh.hostname", Key: "keyword", Op: "LIKE"},
		{Column: "mh.id::text", Key: "id"},
		{Column: "mh.createddate", Key: "createddate"},
		{Column: "mh.updateddate", Key: "updateddate"},
		{Column: "mh.status", Key: "status"},
	},
	SortCols: map[string]string{
		"hostname":    "mh.hostname",
		"createddate": "mh.createddate",
		"updateddate": "mh.updateddate",
	},
	DefaultSort: "mh.updateddate DESC",
}

var ApiSpec = ListSpec{
	Table: "sd_iot_api a",
	Selects: []string{
		"a.api_id AS api_id", "a.api_name AS api_name", "a.host AS host", "a.port AS port",
		"a.token_value AS token_value", "a.password AS password",
		"a.updateddate AS updateddate", "a.status AS status",
	},
	Filters: []Filter{
		{Column: "a.api_name", Key: "keyword", Op: "LIKE"},
		{Column: "a.api_id", Key: "api_id"},
		{Column: "a.createddate", Key: "createddate"},
		{Column: "a.updateddate", Key: "updateddate"},
		{Column: "a.status", Key: "status"},
	},
	SortCols: map[string]string{
		"api_id":      "a.api_id",
		"api_name":    "a.api_name",
		"createddate": "a.createddate",
		"updateddate": "a.updateddate",
	},
	DefaultSort: "a.api_id ASC",
}

var EmailSpec = ListSpec{
	Table: "sd_iot_email e",
	Selects: []string{
		"e.email_id AS email_id", "e.email_name AS email_name", "e.host AS host",
		"e.port AS port", "e.username AS username", "e.password AS password",
		"e.updateddate AS updateddate", "e.status AS status",
	},
	Filters: []Filter{
		{Column: "e.email_name", Key: "keyword", Op: "LIKE"},
		{Column: "e.email_id::text", Key: "email_id"},
		{Column: "e.createddate", Key: "createddate"},
		{Column: "e.updateddate", Key: "updateddate"},
		{Column: "e.status", Key: "status"},
	},
	SortCols: map[string]string{
		"email_id":    "e.email_id",
		"email_name":  "e.email_name",
		"createddate": "e.createddate",
		"updateddate": "e.updateddate",
	},
	DefaultSort: "e.email_id ASC",
}

var HostSpec = ListSpec{
	Table: "sd_iot_host h",
	Selects: []string{
		"h.host_id AS host_id", "h.host_name AS host_name", "h.port AS port",
		"h.username AS username", "h.password AS password",
		"h.createddate AS createddate", "h.updateddate AS updateddate",
		"h.status AS status", "h.idhost AS idhost",
	},
	Filters: []Filter{
		{Column: "h.host_name", Key: "keyword", Op: "LIKE"},
		{Column: "h.host_id::text", Key: "host_id"},
		{Column: "h.createddate", Key: "createddate"},
		{Column: "h.updateddate", Key: "updateddate"},
		{Column: "h.status", Key: "status"},
	},
	SortCols: map[string]string{
		"host_id":     "h.host_id",
		"host_name":   "h.host_name",
		"createddate": "h.createddate",
		"updateddate": "h.updateddate",
	},
	DefaultSort: "h.host_id ASC",
}

var InfluxdbSpec = ListSpec{
	Table: "sd_iot_influxdb i",
	Selects: []string{
		"i.influxdb_id AS influxdb_id", "i.influxdb_name AS influxdb_name", "i.host AS host",
		"i.port AS port", "i.username AS username", "i.password AS password",
		"i.buckets AS buckets", "i.token_value AS token_value",
		"i.updateddate AS updateddate", "i.status AS status",
	},
	Filters: []Filter{
		{Column: "i.influxdb_name", Key: "keyword", Op: "LIKE"},
		{Column: "i.influxdb_id::text", Key: "influxdb_id"},
		{Column: "i.username", Key: "username"},
		{Column: "i.password", Key: "password"},
		{Column: "i.buckets", Key: "buckets"},
		{Column: "i.token_value", Key: "token_value"},
		{Column: "i.status", Key: "status"},
	},
	SortCols: map[string]string{
		"influxdb_id":   "i.influxdb_id",
		"influxdb_name": "i.influxdb_name",
		"createddate":   "i.createddate",
		"updateddate":   "i.updateddate",
	},
	DefaultSort: "i.createddate DESC",
}

var LineSpec = ListSpec{
	Table:   "sd_iot_line l",
	Selects: []string{"l.*"},
	Filters: []Filter{
		{Column: "l.line_name", Key: "keyword", Op: "LIKE"},
		{Column: "l.line_id::text", Key: "line_id"},
		{Column: "l.createddate", Key: "createddate"},
		{Column: "l.updateddate", Key: "updateddate"},
		{Column: "l.status", Key: "status"},
	},
	SortCols: map[string]string{
		"line_id":     "l.line_id",
		"line_name":   "l.line_name",
		"createddate": "l.createddate",
		"updateddate": "l.updateddate",
	},
	DefaultSort: "l.line_id ASC",
}

var NoderedSpec = ListSpec{
	Table: "sd_iot_nodered n",
	Selects: []string{
		"n.nodered_id AS nodered_id", "n.nodered_name AS nodered_name", "n.host AS host",
		"n.port AS port", "n.routing AS routing", "n.client_id AS client_id",
		"n.grant_type AS grant_type", "n.scope AS scope", "n.username AS username",
		"n.password AS password", "n.createddate AS createddate",
		"n.updateddate AS updateddate", "n.status AS status",
	},
	Filters: []Filter{
		{Column: "n.nodered_name", Key: "keyword", Op: "LIKE"},
		{Column: "n.nodered_id::text", Key: "nodered_id"},
		{Column: "n.createddate", Key: "createddate"},
		{Column: "n.updateddate", Key: "updateddate"},
		{Column: "n.status", Key: "status"},
	},
	SortCols: map[string]string{
		"nodered_id":   "n.nodered_id",
		"nodered_name": "n.nodered_name",
		"createddate":  "n.createddate",
		"updateddate":  "n.updateddate",
	},
	DefaultSort: "n.createddate ASC",
}

var SmsSpec = ListSpec{
	Table: "sd_iot_sms s",
	Selects: []string{
		"s.sms_id AS sms_id", "s.sms_name AS sms_name", "s.host AS host", "s.port AS port",
		"s.username AS username", "s.password AS password", "s.apikey AS apikey",
		"s.originator AS originator", "s.updateddate AS updateddate", "s.status AS status",
	},
	Filters: []Filter{
		{Column: "s.sms_name", Key: "keyword", Op: "LIKE"},
		{Column: "s.sms_id::text", Key: "sms_id"},
		{Column: "s.host", Key: "host"},
		{Column: "s.port", Key: "port"},
		{Column: "s.username", Key: "username"},
		{Column: "s.password", Key: "password"},
		{Column: "s.apikey", Key: "apikey"},
		{Column: "s.originator", Key: "originator"},
		{Column: "s.updateddate", Key: "updateddate"},
		{Column: "s.status", Key: "status"},
	},
	SortCols: map[string]string{
		"sms_id":      "s.sms_id",
		"sms_name":    "s.sms_name",
		"createddate": "s.createddate",
		"updateddate": "s.updateddate",
	},
	DefaultSort: "s.createddate DESC",
}

var TokenSpec = ListSpec{
	Table: "sd_iot_token t",
	Selects: []string{
		"t.token_id AS token_id", "t.token_name AS token_name", "t.host AS host",
		"t.port AS port", "t.username AS username", "t.password AS password",
		"t.updateddate AS updateddate", "t.status AS status",
	},
	Filters: []Filter{
		{Column: "t.token_name", Key: "keyword", Op: "LIKE"},
		{Column: "t.token_id", Key: "token_id"},
		{Column: "t.host", Key: "host"},
		{Column: "t.port", Key: "port"},
		{Column: "t.username", Key: "username"},
		{Column: "t.password", Key: "password"},
		{Column: "t.updateddate", Key: "updateddate"},
		{Column: "t.status", Key: "status"},
	},
	SortCols: map[string]string{
		"token_id":    "t.token_id",
		"token_name":  "t.token_name",
		"createddate": "t.createddate",
		"updateddate": "t.updateddate",
	},
	DefaultSort: "t.token_id ASC",
}

var TelegramSpec = ListSpec{
	Table: "sd_iot_telegram tg",
	Selects: []string{
		"tg.telegram_id AS telegram_id", "tg.telegram_name AS telegram_name",
		"tg.host AS host", "tg.port AS port", "tg.username AS username",
		"tg.password AS password", "tg.updateddate AS updateddate", "tg.status AS status",
	},
	Filters: []Filter{
		{Column: "tg.telegram_name", Key: "keyword", Op: "LIKE"},
		{Column: "tg.telegram_id::text", Key: "telegram_id"},
		{Column: "tg.host", Key: "host"},
		{Column: "tg.port", Key: "port"},
		{Column: "tg.username", Key: "username"},
		{Column: "tg.password", Key: "password"},
		{Column: "tg.updateddate", Key: "updateddate"},
		{Column: "tg.status", Key: "status"},
	},
	SortCols: map[string]string{
		"telegram_id":   "tg.telegram_id",
		"telegram_name": "tg.telegram_name",
		"createddate":   "tg.createddate",
		"updateddate":   "tg.updateddate",
	},
	DefaultSort: "tg.telegram_id ASC",
}

var DashboardConfigSpec = ListSpec{
	Table: "sd_dashboard_config dc",
	Selects: []string{
		"dc.id AS id", "dc.location_id AS location_id", "dc.name AS name",
		"dc.config_data AS config_data", "dc.status AS status",
		"dc.created_date AS created_date", "dc.updated_date AS updated_date",
	},
	Filters: []Filter{
		{Column: "dc.name", Key: "keyword", Op: "LIKE"},
		{Column: "dc.id", Key: "id"},
		{Column: "dc.location_id", Key: "location_id"},
		{Column: "dc.status", Key: "status"},
	},
	SortCols:    map[string]string{"id": "dc.id", "name": "dc.name", "created_date": "dc.created_date"},
	DefaultSort: "dc.id ASC",
}
