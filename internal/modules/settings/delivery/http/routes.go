package http

import (
	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/settings"

	"github.com/go-chi/chi/v5"
)

// MapSettingRoute registers every settings endpoint under /api/settings.
// Paths + verbs mirror the NestJS controller 1:1 (including GET deletes).
func MapSettingRoute(router chi.Router, h settings.Handlers, mw *middleware.MiddlewareManager) {
	router.Route("/settings", func(r chi.Router) {
		r.Use(mw.Verifier(true))
		r.Use(mw.Authenticator())
		r.Use(mw.CurrentUser())
		r.Use(mw.ActiveUser())

		// ---- Master: setting ----
		r.Get("/listsetting", h.ListSetting())
		r.Get("/settingall", h.SettingAll())
		r.Post("/createsetting", h.CreateSetting())
		r.Post("/updatesetting", h.UpdateSetting())
		r.Get("/deletesetting", h.DeleteSettingViaGet())
		r.Delete("/deletesetting", h.DeleteSetting())

		// ---- Master: location ----
		r.Get("/listlocation", h.ListLocation())
		r.Get("/locationall", h.LocationAll())
		r.Post("/createlocation", h.CreateLocation())
		r.Post("/updatelocation", h.UpdateLocation())
		r.Get("/deletelocation", h.DeleteLocation())

		// ---- Master: type ----
		r.Get("/listtype", h.ListType())
		r.Get("/typeall", h.TypeAll())
		r.Post("/createtype", h.CreateType())
		r.Post("/updatetype", h.UpdateType())
		r.Get("/deletetype", h.DeleteType())

		// ---- Master: devicetype ----
		r.Get("/listdevicetype", h.ListDeviceTypePage())
		r.Get("/devicetypeall", h.DeviceTypeAll()) // declared twice in NestJS; registered once
		r.Get("/devicetypeallcontrol", h.DeviceTypeAllControl())
		r.Post("/createdevicetype", h.CreateDeviceType())
		r.Post("/updatedevicetype", h.UpdateDeviceType())
		r.Get("/deletedevicetype", h.DeleteDeviceType())

		// ---- Master: group ----
		r.Get("/lisgroup", h.ListGroup())
		r.Get("/lisgroupall", h.GroupAll())
		r.Get("/listgrouppage", h.ListGroupPage())
		r.Post("/creategroup", h.CreateGroup())
		r.Post("/updategroup", h.UpdateGroup())
		r.Get("/deletegroup", h.DeleteGroup())

		// ---- Master: sensor ----
		r.Get("/listsensor", h.ListSensor())
		r.Get("/sensorall", h.SensorAll())
		r.Post("/createsensor", h.CreateSensor())
		r.Post("/updatesensor", h.UpdateSensor())
		r.Get("/deletesensor", h.DeleteSensor())

		// ---- Utility ----
		r.Get("/sensortype", h.SensorType())
		r.Get("/testgemail", h.TestGmailConnection())
		r.Get("/sendemail", h.SendEmail())
		r.Get("/mqttdata", h.MqttData())

		// ---- Device ----
		r.Get("/deviceall", h.ListDeviceAll())
		r.Get("/listdevicepage", h.ListDevicePage())
		r.Get("/listdevicepagess", h.ListDevicePagess())
		r.Get("/listdevicepageactive1", h.ListDevicePageActive1())
		r.Get("/listdevicepageactive", h.ListDevicePageActive())
		r.Get("/listdevicepageall", h.ListDevicePageAll())
		r.Get("/listdevicepagesensor", h.ListDevicePageSensor())
		r.Get("/listdevicepageallactive", h.ListDevicePageAllActive())
		r.Get("/listdevicepageallactiveschedule", h.ListDevicePageAllActiveSchedule())
		r.Get("/deviceeditget", h.DeviceEditGet())
		r.Get("/devicedetail", h.DeviceDetail())
		r.Get("/devicedelete", h.DeviceDeleteCheck())
		r.Get("/deletedevice", h.DeleteDevice())
		r.Post("/createdevice", h.CreateDevice())
		r.Post("/updatedevice", h.UpdateDevice())
		r.Post("/updatestatusdeviceid", h.UpdateStatusDeviceId())
		r.Post("/deviceactionuser", h.DeviceActionUser())

		// ---- Schedule ----
		r.Get("/listscheduledevice", h.ListScheduleDevice())
		r.Get("/findscheduledevicechk", h.FindScheduleDeviceChk())
		r.Get("/schedulelist", h.ScheduleList())
		r.Get("/scheduleall", h.ScheduleAll())
		r.Get("/listschedulepage", h.ListSchedulePage())
		r.Get("/scheduledevicepage", h.ScheduleDevicePage())
		r.Get("/listdevicescheduledata", h.ListDeviceScheduleData())
		r.Get("/createscheduledevice", h.CreateScheduleDeviceViaGet())
		r.Get("/deletescheduledevice", h.DeleteScheduleDevices())
		r.Get("/deletedeviceschedule", h.DeleteDeviceSchedule())
		r.Get("/deletedeviceandschedule", h.DeleteDeviceAndSchedule())
		r.Post("/createschedule", h.CreateSchedule())
		r.Post("/createscheduledevice", h.CreateScheduleDevice())
		r.Post("/updateschedule", h.UpdateSchedule())
		r.Post("/updateschedulestatus", h.UpdateScheduleStatus())
		r.Post("/updatescheduledaystatus", h.UpdateScheduleDayStatus())
		r.Get("/deleteschedule", h.DeleteSchedule())

		// ---- Integration: mqtt ----
		r.Get("/lismqtt", h.ListMqtt())
		r.Get("/listmqtt", h.ListMqttAlt())
		r.Get("/lismqttall", h.MqttAll())
		r.Get("/getmqttdetail", h.GetMqttDetail())
		r.Get("/listmqttpaginate", h.ListMqttPaginate())
		r.Get("/listmqttpaginateactive", h.ListMqttPaginateActive())
		r.Get("/listmqttdevicepaginate", h.ListMqttDevicePaginate())
		r.Get("/mqttdelete", h.DeleteMqtt())
		r.Get("/deletemqtt", h.DeleteMqttAlt())
		r.Post("/createmqtt", h.CreateMqtt())
		r.Post("/updatemqtt", h.UpdateMqtt())
		r.Post("/updatemqttstatus", h.UpdateMqttStatus())
		r.Post("/mqtttsort", h.UpdateMqtttSort())

		// ---- Integration: mqtthost ----
		r.Get("/listmqtthost", h.ListMqttHost())
		r.Get("/mqtthostall", h.MqttHostAll())
		r.Post("/createmqtthost", h.CreateMqttHost())
		r.Post("/updatemqtthost", h.UpdateMqttHost())
		r.Post("/updatemqtthoststatus", h.UpdateMqttHostStatus())
		r.Get("/deletemqtthost", h.DeleteMqttHost())

		// ---- Integration: api ----
		r.Get("/apiall", h.ApiAll())
		r.Get("/listapipage", h.ListApiPage())
		r.Post("/createapi", h.CreateApi())
		r.Post("/updateapi", h.UpdateApi())
		r.Get("/deleteapi", h.DeleteApi())

		// ---- Integration: email ----
		r.Get("/listemail", h.ListEmail())
		r.Get("/emailall", h.EmailAll())
		r.Post("/createemail", h.CreateEmail())
		r.Post("/updateemail", h.UpdateEmail())
		r.Post("/updateemailstatus", h.UpdateEmailStatus())
		r.Get("/deleteemail", h.DeleteEmail())

		// ---- Integration: host ----
		r.Get("/hostall", h.HostAll())
		r.Get("/listhostpage", h.ListHostPage())
		r.Post("/createhost", h.CreateHost())
		r.Post("/updatehost", h.UpdateHost())
		r.Get("/deletehost", h.DeleteHost())

		// ---- Integration: influxdb ----
		r.Get("/influxdball", h.InfluxdbAll())
		r.Get("/listinfluxdbpage", h.ListInfluxdbPage())
		r.Post("/createinfluxdb", h.CreateInfluxdb())
		r.Post("/updateinfluxdb", h.UpdateInfluxdb())
		r.Post("/updateinfluxdbstatus", h.UpdateInfluxdbStatus())
		r.Get("/deleteinfluxdb", h.DeleteInfluxdb())

		// ---- Integration: line ----
		r.Get("/lineall", h.LineAll())
		r.Get("/listlinepage", h.ListLinePage())
		r.Post("/createline", h.CreateLine())
		r.Post("/updateline", h.UpdateLine())
		r.Post("/updatelinestatus", h.UpdateLineStatus())
		r.Get("/deleteline", h.DeleteLine())

		// ---- Integration: nodered ----
		r.Get("/noderedall", h.NoderedAll())
		r.Get("/listnoderedpaginate", h.ListNoderedPaginate())
		r.Post("/createnodered", h.CreateNodered())
		r.Post("/updatenodered", h.UpdateNodered())
		r.Post("/updatenoderedstatus", h.UpdateNoderedStatus())
		r.Get("/deletenodered", h.DeleteNodered())

		// ---- Integration: sms ----
		r.Get("/smsall", h.SmsAll())
		r.Get("/listsmspage", h.ListSmsPage())
		r.Post("/createsms", h.CreateSms())
		r.Post("/updatesms", h.UpdateSms())
		r.Post("/updatesmsstatus", h.UpdateSmsStatus())
		r.Get("/deletesms", h.DeleteSms())

		// ---- Integration: token ----
		r.Get("/tokenall", h.TokenAll())
		r.Get("/tokensmspage", h.ListTokenPage())
		r.Post("/createtoken", h.CreateToken())
		r.Post("/updatetoken", h.UpdateToken())
		r.Get("/deletetoken", h.DeleteToken())

		// ---- Integration: telegram ----
		r.Post("/createtelegram", h.CreateTelegram())
		r.Post("/updatetelegram", h.UpdateTelegram())
		r.Get("/deletetelegram", h.DeleteTelegram())

		// ---- Dashboard config (static before {id}) ----
		r.Post("/dashboardconfig", h.CreateDashboardConfig())
		r.Get("/dashboardconfig_1", h.DashboardConfigByLocation())
		r.Get("/dashboardconfig", h.ListDashboardConfig())
		r.Get("/dashboardconfig/search", h.FindOrCreateDashboardConfig())
		r.Get("/dashboardconfig/{id}", h.GetDashboardConfig())
		r.Patch("/dashboardconfig/{id}", h.UpdateDashboardConfig())
		r.Delete("/dashboardconfig/{id}", h.RemoveDashboardConfig())

		// ---- Alarm device/event links ----
		r.Get("/listalarmdevicepage", h.ListAlarmDevicePage())
		r.Get("/listalarmdeviceactivepage", h.ListAlarmDeviceActivePage())
		r.Get("/listalarmeventdevicepage", h.ListAlarmEventDevicePage())
		r.Get("/listalarmeventdevicecontrolpage", h.ListAlarmEventDeviceControlPage())
		r.Get("/activealarmdevicepage", h.ActiveAlarmDevicePage())
		r.Get("/activealarmeventdeviceeventpage", h.ActiveAlarmEventDeviceEventPage())
		r.Get("/alarmdevice", h.AlarmDevice())
		r.Get("/alarmdevicestatus", h.AlarmDeviceStatus())
		r.Get("/deviceactivemqttalarm", h.DeviceActiveMqttAlarm())
		r.Get("/devicealarm", h.DeviceAlarm())
		r.Get("/createalarmdevice", h.CreateAlarmDeviceViaGet())
		r.Get("/deletealarmdevice", h.DeleteAlarmDevices())
		r.Get("/createalarmeventdevice", h.CreateAlarmEventDeviceViaGet())
		r.Get("/deletealarmeventdevice", h.DeleteAlarmEventDevices())
		r.Get("/deletearmdevice", h.DeleteArmDevice())
		r.Get("/deletearmdevicev2", h.DeleteArmDeviceV2())
		r.Post("/createalarmDevice", h.CreateAlarmDevice())
		r.Post("/createalarmdevicepaginate", h.CreateAlarmDevicePaginate())
		r.Post("/updatealarmdevice", h.UpdateAlarmDevice())
		r.Post("/updatealarmstatus", h.UpdateAlarmStatus())
		r.Post("/createdevicealarmaction", h.CreateDeviceAlarmAction())

		// ---- Device alarm views / monitors ----
		r.Get("/listdevicealarm", h.ListDeviceAlarm())
		r.Get("/listdevicealarmairV1", h.ListDeviceAlarmAirV1())
		r.Get("/listdevicealarmair", h.ListDeviceAlarmAir())
		r.Get("/listdevicealarmall", h.ListDeviceAlarmAll())
		r.Get("/_listdevicealarmfan", h.UnderscoreListDeviceAlarmFan())
		r.Get("/listdevicealarmfan", h.ListDeviceAlarmFan())
		r.Get("/listdevicealarmlimit", h.ListDeviceAlarmLimit())
		r.Get("/_devicemonitor", h.UnderscoreDeviceMonitor())
		r.Get("/devicemonitor", h.DeviceMonitor())
		r.Get("/devicemonitors", h.DeviceMonitors())

		// ---- Process logs ----
		r.Get("/scheduleproces", h.ScheduleProces())
		r.Get("/scheduleprocesslog", h.ScheduleProcessLog())
		r.Get("/scheduleprocesslogpaginate", h.ScheduleProcessLogPaginate())
		r.Get("/mqtterrorlogpaginate", h.MqttErrorLogPaginate())
		r.Get("/alarmlogpaginate", h.AlarmLogPaginate())
		r.Get("/alarmlogpaginateemail", h.AlarmLogPaginateEmail())
		r.Get("/alarmlogpaginateline", h.AlarmLogPaginateLine())
		r.Get("/alarmlogpaginatesms", h.AlarmLogPaginateSms())
		r.Get("/alarmlogpaginatetelegram", h.AlarmLogPaginateTelegram())
		r.Get("/alarmlogpaginatecontrols", h.AlarmLogPaginateControls())
		r.Get("/alarmlogpaginatecontrol", h.AlarmLogPaginateControl())
	})
}
