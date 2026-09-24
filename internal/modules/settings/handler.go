package settings

import "net/http"

// Handlers is the HTTP boundary of the settings module.
// One method per endpoint, named after the route path (PascalCase).
type Handlers interface {
	// ---- Master: setting ----
	ListSetting() http.HandlerFunc
	SettingAll() http.HandlerFunc
	CreateSetting() http.HandlerFunc
	UpdateSetting() http.HandlerFunc
	DeleteSettingViaGet() http.HandlerFunc // GET /deletesetting
	DeleteSetting() http.HandlerFunc       // DELETE /deletesetting

	// ---- Master: location ----
	ListLocation() http.HandlerFunc
	LocationAll() http.HandlerFunc
	CreateLocation() http.HandlerFunc
	UpdateLocation() http.HandlerFunc
	DeleteLocation() http.HandlerFunc

	// ---- Master: type ----
	ListType() http.HandlerFunc
	TypeAll() http.HandlerFunc
	CreateType() http.HandlerFunc
	UpdateType() http.HandlerFunc
	DeleteType() http.HandlerFunc

	// ---- Master: devicetype ----
	ListDeviceTypePage() http.HandlerFunc
	DeviceTypeAll() http.HandlerFunc // GET /devicetypeall (registered once)
	DeviceTypeAllControl() http.HandlerFunc
	CreateDeviceType() http.HandlerFunc
	UpdateDeviceType() http.HandlerFunc
	DeleteDeviceType() http.HandlerFunc

	// ---- Master: group ----
	ListGroup() http.HandlerFunc
	GroupAll() http.HandlerFunc
	ListGroupPage() http.HandlerFunc
	CreateGroup() http.HandlerFunc
	UpdateGroup() http.HandlerFunc
	DeleteGroup() http.HandlerFunc

	// ---- Master: sensor ----
	ListSensor() http.HandlerFunc
	SensorAll() http.HandlerFunc
	CreateSensor() http.HandlerFunc
	UpdateSensor() http.HandlerFunc
	DeleteSensor() http.HandlerFunc

	// ---- Utility ----
	SensorType() http.HandlerFunc
	TestGmailConnection() http.HandlerFunc
	SendEmail() http.HandlerFunc
	MqttData() http.HandlerFunc

	// ---- Device ----
	ListDeviceAll() http.HandlerFunc
	ListDevicePage() http.HandlerFunc
	ListDevicePagess() http.HandlerFunc
	ListDevicePageActive1() http.HandlerFunc
	ListDevicePageActive() http.HandlerFunc
	ListDevicePageAll() http.HandlerFunc
	ListDevicePageSensor() http.HandlerFunc
	ListDevicePageAllActive() http.HandlerFunc
	ListDevicePageAllActiveSchedule() http.HandlerFunc
	DeviceEditGet() http.HandlerFunc
	DeviceDetail() http.HandlerFunc
	DeviceDeleteCheck() http.HandlerFunc
	DeleteDevice() http.HandlerFunc
	CreateDevice() http.HandlerFunc
	UpdateDevice() http.HandlerFunc
	UpdateStatusDeviceId() http.HandlerFunc
	DeviceActionUser() http.HandlerFunc

	// ---- Schedule ----
	ListScheduleDevice() http.HandlerFunc
	FindScheduleDeviceChk() http.HandlerFunc
	ScheduleList() http.HandlerFunc
	ScheduleAll() http.HandlerFunc
	ListSchedulePage() http.HandlerFunc
	ScheduleDevicePage() http.HandlerFunc
	ListDeviceScheduleData() http.HandlerFunc
	CreateScheduleDeviceViaGet() http.HandlerFunc
	DeleteScheduleDevices() http.HandlerFunc
	DeleteDeviceSchedule() http.HandlerFunc
	DeleteDeviceAndSchedule() http.HandlerFunc
	CreateSchedule() http.HandlerFunc
	CreateScheduleDevice() http.HandlerFunc
	UpdateSchedule() http.HandlerFunc
	UpdateScheduleStatus() http.HandlerFunc
	UpdateScheduleDayStatus() http.HandlerFunc
	DeleteSchedule() http.HandlerFunc

	// ---- Integration: mqtt ----
	ListMqtt() http.HandlerFunc
	ListMqttAlt() http.HandlerFunc
	MqttAll() http.HandlerFunc
	GetMqttDetail() http.HandlerFunc
	ListMqttPaginate() http.HandlerFunc
	ListMqttPaginateActive() http.HandlerFunc
	ListMqttDevicePaginate() http.HandlerFunc
	DeleteMqtt() http.HandlerFunc    // GET /mqttdelete
	DeleteMqttAlt() http.HandlerFunc // GET /deletemqtt (existence-checked)
	CreateMqtt() http.HandlerFunc
	UpdateMqtt() http.HandlerFunc
	UpdateMqttStatus() http.HandlerFunc
	UpdateMqtttSort() http.HandlerFunc

	// ---- Integration: mqtthost ----
	ListMqttHost() http.HandlerFunc
	MqttHostAll() http.HandlerFunc
	CreateMqttHost() http.HandlerFunc
	UpdateMqttHost() http.HandlerFunc
	UpdateMqttHostStatus() http.HandlerFunc
	DeleteMqttHost() http.HandlerFunc

	// ---- Integration: api ----
	ApiAll() http.HandlerFunc
	ListApiPage() http.HandlerFunc
	CreateApi() http.HandlerFunc
	UpdateApi() http.HandlerFunc
	DeleteApi() http.HandlerFunc

	// ---- Integration: email ----
	ListEmail() http.HandlerFunc
	EmailAll() http.HandlerFunc
	CreateEmail() http.HandlerFunc
	UpdateEmail() http.HandlerFunc
	UpdateEmailStatus() http.HandlerFunc
	DeleteEmail() http.HandlerFunc

	// ---- Integration: host ----
	HostAll() http.HandlerFunc
	ListHostPage() http.HandlerFunc
	CreateHost() http.HandlerFunc
	UpdateHost() http.HandlerFunc
	DeleteHost() http.HandlerFunc

	// ---- Integration: influxdb ----
	InfluxdbAll() http.HandlerFunc
	ListInfluxdbPage() http.HandlerFunc
	CreateInfluxdb() http.HandlerFunc
	UpdateInfluxdb() http.HandlerFunc
	UpdateInfluxdbStatus() http.HandlerFunc
	DeleteInfluxdb() http.HandlerFunc

	// ---- Integration: line ----
	LineAll() http.HandlerFunc
	ListLinePage() http.HandlerFunc
	CreateLine() http.HandlerFunc
	UpdateLine() http.HandlerFunc
	UpdateLineStatus() http.HandlerFunc
	DeleteLine() http.HandlerFunc

	// ---- Integration: nodered ----
	NoderedAll() http.HandlerFunc
	ListNoderedPaginate() http.HandlerFunc
	CreateNodered() http.HandlerFunc
	UpdateNodered() http.HandlerFunc
	UpdateNoderedStatus() http.HandlerFunc
	DeleteNodered() http.HandlerFunc

	// ---- Integration: sms ----
	SmsAll() http.HandlerFunc
	ListSmsPage() http.HandlerFunc
	CreateSms() http.HandlerFunc
	UpdateSms() http.HandlerFunc
	UpdateSmsStatus() http.HandlerFunc
	DeleteSms() http.HandlerFunc

	// ---- Integration: token ----
	TokenAll() http.HandlerFunc
	ListTokenPage() http.HandlerFunc
	CreateToken() http.HandlerFunc
	UpdateToken() http.HandlerFunc
	DeleteToken() http.HandlerFunc

	// ---- Integration: telegram ----
	CreateTelegram() http.HandlerFunc
	UpdateTelegram() http.HandlerFunc
	DeleteTelegram() http.HandlerFunc

	// ---- Dashboard config ----
	CreateDashboardConfig() http.HandlerFunc
	DashboardConfigByLocation() http.HandlerFunc
	ListDashboardConfig() http.HandlerFunc
	FindOrCreateDashboardConfig() http.HandlerFunc
	GetDashboardConfig() http.HandlerFunc
	UpdateDashboardConfig() http.HandlerFunc
	RemoveDashboardConfig() http.HandlerFunc

	// ---- Alarm device/event links ----
	ListAlarmDevicePage() http.HandlerFunc
	ListAlarmDeviceActivePage() http.HandlerFunc
	ListAlarmEventDevicePage() http.HandlerFunc
	ListAlarmEventDeviceControlPage() http.HandlerFunc
	ActiveAlarmDevicePage() http.HandlerFunc
	ActiveAlarmEventDeviceEventPage() http.HandlerFunc
	AlarmDevice() http.HandlerFunc
	AlarmDeviceStatus() http.HandlerFunc
	DeviceActiveMqttAlarm() http.HandlerFunc
	DeviceAlarm() http.HandlerFunc
	CreateAlarmDeviceViaGet() http.HandlerFunc
	DeleteAlarmDevices() http.HandlerFunc
	CreateAlarmEventDeviceViaGet() http.HandlerFunc
	DeleteAlarmEventDevices() http.HandlerFunc
	DeleteArmDevice() http.HandlerFunc
	DeleteArmDeviceV2() http.HandlerFunc
	CreateAlarmDevice() http.HandlerFunc
	CreateAlarmDevicePaginate() http.HandlerFunc // POST /createalarmdevicepaginate (SN dup check)
	UpdateAlarmDevice() http.HandlerFunc
	UpdateAlarmStatus() http.HandlerFunc
	CreateDeviceAlarmAction() http.HandlerFunc

	// ---- Device alarm views / monitors ----
	ListDeviceAlarm() http.HandlerFunc
	ListDeviceAlarmAirV1() http.HandlerFunc
	ListDeviceAlarmAir() http.HandlerFunc
	ListDeviceAlarmAll() http.HandlerFunc
	UnderscoreListDeviceAlarmFan() http.HandlerFunc
	ListDeviceAlarmFan() http.HandlerFunc
	ListDeviceAlarmLimit() http.HandlerFunc
	UnderscoreDeviceMonitor() http.HandlerFunc
	DeviceMonitor() http.HandlerFunc
	DeviceMonitors() http.HandlerFunc

	// ---- Process logs ----
	ScheduleProces() http.HandlerFunc
	ScheduleProcessLog() http.HandlerFunc
	ScheduleProcessLogPaginate() http.HandlerFunc
	MqttErrorLogPaginate() http.HandlerFunc
	AlarmLogPaginate() http.HandlerFunc
	AlarmLogPaginateEmail() http.HandlerFunc
	AlarmLogPaginateLine() http.HandlerFunc
	AlarmLogPaginateSms() http.HandlerFunc
	AlarmLogPaginateTelegram() http.HandlerFunc
	AlarmLogPaginateControls() http.HandlerFunc
	AlarmLogPaginateControl() http.HandlerFunc
}
