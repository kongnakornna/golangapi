# Settings API — `/api/settings`

> เจ้าของ: icmongolang · Framework: Go + Chi · Auth: **JWT Bearer (ทุก endpoint)**
> Middleware: `Verifier(true)` → `Authenticator()` → `CurrentUser()` → `ActiveUser()`
> ไม่มี endpoint public ใน module นี้

---

## Response Envelope

ทุก endpoint ใช้ `pkg/responses` envelope เดียวกัน:

```jsonc
// Success (HTTP 200)
{ "data": <payload>, "is_success": true }

// Error
{
  "data": null,
  "error": { "err": "<string>", "status": <int>, "status_text": "<string>", "msg": "<string>" },
  "is_success": false
}
```

### รูปแบบ paginated list
```jsonc
{
  "data": {
    "page": 1, "currentPage": 1, "pageSize": 10,
    "totalPages": 1, "total": 0, "filter": {}, "data": []
  },
  "is_success": true
}
```

### หลักการของ CRUD (generic handler)
- **Create** (`createHandler`): ต้องมี field อย่างน้อย 1 ตัว → ถ้าไม่มี คืน `400 "request body has no valid fields"`; บางตัวมี `dupCheck` → `422` ถ้าซ้ำ
- **Update** (`updateBodyHandler`): body ต้องมี id column (เช่น `device_id`, `schedule_id`) → ถ้าไม่มี `400 "<idcol> is required"`; เพิ่ม `updateddate` อัตโนมัติ
- **Delete via GET** (`deleteGetHandler`): ใช้ query param `id`; ถ้าไม่พบ → `404 "Data not found."`
- พารามิเตอร์ที่ส่งนอก **whitelist** จะถูกละเลย

### Query params มาตรฐานของ list/page
- `page` (int, default 1)
- `pageSize` (int, default 1000, max 5000)
- `sort` — `field-ASC` / `field-DESC`
- `keyword` — filter LIKE
- field เฉพาะ endpoint

---

## 1. Master: Setting

| # | Method | Path | Handler |
|---|--------|------|---------|
| 1 | GET | `/settings/listsetting` | `ListSetting` |
| 2 | GET | `/settings/settingall` | `SettingAll` |
| 3 | POST | `/settings/createsetting` | `CreateSetting` |
| 4 | POST | `/settings/updatesetting` | `UpdateSetting` |
| 5 | GET | `/settings/deletesetting` | `DeleteSettingViaGet` |
| 6 | DELETE | `/settings/deletesetting` | `DeleteSetting` |

## 2. Master: Location

| # | Method | Path | Handler |
|---|--------|------|---------|
| 7 | GET | `/settings/listlocation` | `ListLocation` |
| 8 | GET | `/settings/locationall` | `LocationAll` |
| 9 | POST | `/settings/createlocation` | `CreateLocation` |
| 10 | POST | `/settings/updatelocation` | `UpdateLocation` |
| 11 | GET | `/settings/deletelocation` | `DeleteLocation` |

### POST `/settings/createlocation`
- Body whitelist (`locationCols`): `location_name`, `ipaddress`, `location_detail`, `configdata`, `status` (default `status:1`)
- **dupCheck**: `ipaddress` ซ้ำ → 422
- Response: แถวที่ insert

## 3. Master: Type

| # | Method | Path | Handler |
|---|--------|------|---------|
| 12 | GET | `/settings/listtype` | `ListType` |
| 13 | GET | `/settings/typeall` | `TypeAll` |
| 14 | POST | `/settings/createtype` | `CreateType` |
| 15 | POST | `/settings/updatetype` | `UpdateType` |
| 16 | GET | `/settings/deletetype` | `DeleteType` |

## 4. Master: DeviceType

| # | Method | Path | Handler |
|---|--------|------|---------|
| 17 | GET | `/settings/listdevicetype` | `ListDeviceTypePage` |
| 18 | GET | `/settings/devicetypeall` | `DeviceTypeAll` |
| 19 | GET | `/settings/devicetypeallcontrol` | `DeviceTypeAllControl` |
| 20 | POST | `/settings/createdevicetype` | `CreateDeviceType` |
| 21 | POST | `/settings/updatedevicetype` | `UpdateDeviceType` |
| 22 | GET | `/settings/deletedevicetype` | `DeleteDeviceType` |

## 5. Master: Group

| # | Method | Path | Handler |
|---|--------|------|---------|
| 23 | GET | `/settings/lisgroup` | `ListGroup` |
| 24 | GET | `/settings/lisgroupall` | `GroupAll` |
| 25 | GET | `/settings/listgrouppage` | `ListGroupPage` |
| 26 | POST | `/settings/creategroup` | `CreateGroup` |
| 27 | POST | `/settings/updategroup` | `UpdateGroup` |
| 28 | GET | `/settings/deletegroup` | `DeleteGroup` |

## 6. Master: Sensor

| # | Method | Path | Handler |
|---|--------|------|---------|
| 29 | GET | `/settings/listsensor` | `ListSensor` |
| 30 | GET | `/settings/sensorall` | `SensorAll` |
| 31 | POST | `/settings/createsensor` | `CreateSensor` |
| 32 | POST | `/settings/updatesensor` | `UpdateSensor` |
| 33 | GET | `/settings/deletesensor` | `DeleteSensor` |

## 7. Utility

| # | Method | Path | Handler |
|---|--------|------|---------|
| 34 | GET | `/settings/sensortype` | `SensorType` |
| 35 | GET | `/settings/testgemail` | `TestGmailConnection` |
| 36 | GET | `/settings/sendemail` | `SendEmail` |
| 37 | GET | `/settings/mqttdata` | `MqttData` |

### GET `/settings/sensortype`
- Query: `lang` (en default / th)

### GET `/settings/testgemail`
- ไม่มี param; Response `presenter.SendEmailResult`:
```jsonc
{"success":bool,"code":int,"to":"","subject":"","content":"","message":"","message_th":""}
```

### GET `/settings/sendemail`
- Query: `to`, `subject`, `content` (optional) — stub

### GET `/settings/mqttdata`
- Query: `mqttdata` (topic) — Redis cache → request-and-wait (cache 30s)
- Response `presenter.MqttDataResult`: `{"getdataFrom":"<cache|mqtt>","payload":[]}`

## 8. Device

| # | Method | Path | Handler |
|---|--------|------|---------|
| 38 | GET | `/settings/deviceall` | `ListDeviceAll` |
| 39 | GET | `/settings/listdevicepage` | `ListDevicePage` |
| 40 | GET | `/settings/listdevicepagess` | `ListDevicePagess` |
| 41 | GET | `/settings/listdevicepageactive1` | `ListDevicePageActive1` |
| 42 | GET | `/settings/listdevicepageactive` | `ListDevicePageActive` |
| 43 | GET | `/settings/listdevicepageall` | `ListDevicePageAll` |
| 44 | GET | `/settings/listdevicepagesensor` | `ListDevicePageSensor` |
| 45 | GET | `/settings/listdevicepageallactive` | `ListDevicePageAllActive` |
| 46 | GET | `/settings/listdevicepageallactiveschedule` | `ListDevicePageAllActiveSchedule` |
| 47 | GET | `/settings/deviceeditget` | `DeviceEditGet` |
| 48 | GET | `/settings/devicedetail` | `DeviceDetail` |
| 49 | GET | `/settings/devicedelete` | `DeviceDeleteCheck` |
| 50 | GET | `/settings/deletedevice` | `DeleteDevice` |
| 51 | POST | `/settings/createdevice` | `CreateDevice` |
| 52 | POST | `/settings/updatedevice` | `UpdateDevice` |
| 53 | POST | `/settings/updatestatusdeviceid` | `UpdateStatusDeviceId` |
| 54 | POST | `/settings/deviceactionuser` | `DeviceActionUser` |

### GET `/settings/listdevicepage`
Query params ฟิลเตอร์: `keyword` (LIKE `d.device_name`), `device_id`, `mqtt_id`, `setting_id`, `type_id`, `location_id`, `org`, `bucket`, `sn`, `createddate`, `updateddate`, `status_warning`, `recovery_warning`, `status_alert`, `recovery_alert`, `time_life`, `period`, `max`, `min`, `hardware_id`, `model`, `vendor`, `comparevalue`, `oid`, `action_id`, `mqtt_data_value`, `mqtt_data_control`, `work_status`, `menu`

sort whitelist: `device_id`, `device_name`, `sort`, `status`, `work_status`, `createddate`, `updateddate` (default `mq.sort ASC, d.device_id ASC`)

Response แต่ละแถว: `device_id`, `mqtt_id`, `setting_id`, `type_id`, `device_name`, `sn`, `hardware_id`, `status_warning`, `recovery_warning`, `status_alert`, `recovery_alert`, `time_life`, `period`, `work_status`, `layout`, `menu`, `max`, `min`, `oid`, `mqtt_data_value`, `mqtt_data_control`, `model`, `vendor`, `comparevalue`, `createddate`, `updateddate`, `status`, `unit`, `action_id`, `status_alert_id`, `measurement`, `mqtt_control_on`, `mqtt_control_off`, `device_org`, `device_bucket`, `type_name`, `location_name`, `configdata`, `mqtt_name`, `mqtt_org`, `mqtt_bucket`, `mqtt_envavorment`, `mqtt_host`, `mqtt_port`, `mqtt_device_name`, `mqtt_status_over_name`, `mqtt_status_data_name`, `mqtt_act_relay_name`, `mqtt_control_relay_name`, `host_name`, `port`, `host_id`, `hardware_type_name` (1=Sensor,2=IO Sensor,3=IO Control,4=Critical Sensor), `layoutapp` (1=Right Menu,2=Card,3=Left Menu,4=Footer Menu), `calibration_add`, `calibration_subtract`, `calibration_type`, `calibrationtype`

### POST `/settings/createdevice`
- Body whitelist (`deviceCols`): `setting_id`, `type_id`, `location_id`, `device_name`, `sn`, `hardware_id`, `status_warning`, `recovery_warning`, `status_alert`, `recovery_alert`, `time_life`, `period`, `work_status`, `model`, `vendor`, `comparevalue`, `unit`, `mqtt_id`, `oid`, `action_id`, `status_alert_id`, `measurement`, `mqtt_data_value`, `mqtt_data_control`, `mqtt_control_on`, `mqtt_control_off`, `org`, `bucket`, `status`, `mqtt_device_name`, `mqtt_status_over_name`, `mqtt_status_data_name`, `mqtt_act_relay_name`, `mqtt_control_relay_name`, `mqtt_config`, `max`, `min`, `layout`, `alert_set`, `menu`, `calibration_add`, `calibration_subtract`, `calibration_type`
- defaults: `status:1`, `work_status:1`
- **dupCheck**: `sn` ซ้ำ → 422 `"The SN <val> duplicate this data cannot create."`
- Response: แถวที่ insert

### GET `/settings/listdevicepageallactiveschedule`
- Query: `schedule_id` (required → 400 `"schedule_id is null."`); `hardware_id` (optional, default `3` = IO Control)
- Response: device list + ต่อแถว `schedule_id`, `schedule_status` (0/1), `count_schedule_device`, `schedule_name`, `schedule_start`, `schedule_title`

### GET `/settings/deviceall`
- ไม่มี filter, ไม่ paginate — คืนทุกแถว `sd_iot_device` sort `device_id ASC`

## 9. Schedule

| # | Method | Path | Handler |
|---|--------|------|---------|
| 55 | GET | `/settings/listscheduledevice` | `ListScheduleDevice` |
| 56 | GET | `/settings/findscheduledevicechk` | `FindScheduleDeviceChk` |
| 57 | GET | `/settings/schedulelist` | `ScheduleList` |
| 58 | GET | `/settings/scheduleall` | `ScheduleAll` |
| 59 | GET | `/settings/listschedulepage` | `ListSchedulePage` |
| 60 | GET | `/settings/scheduledevicepage` | `ScheduleDevicePage` |
| 61 | GET | `/settings/listdevicescheduledata` | `ListDeviceScheduleData` |
| 62 | GET | `/settings/createscheduledevice` | `CreateScheduleDeviceViaGet` |
| 63 | GET | `/settings/deletescheduledevice` | `DeleteScheduleDevices` |
| 64 | GET | `/settings/deletedeviceschedule` | `DeleteDeviceSchedule` |
| 65 | GET | `/settings/deletedeviceandschedule` | `DeleteDeviceAndSchedule` |
| 66 | POST | `/settings/createschedule` | `CreateSchedule` |
| 67 | POST | `/settings/createscheduledevice` | `CreateScheduleDevice` |
| 68 | POST | `/settings/updateschedule` | `UpdateSchedule` |
| 69 | POST | `/settings/updateschedulestatus` | `UpdateScheduleStatus` |
| 70 | POST | `/settings/updatescheduledaystatus` | `UpdateScheduleDayStatus` |
| 71 | GET | `/settings/deleteschedule` | `DeleteSchedule` |

### POST `/settings/createschedule`
- Body whitelist (`scheduleCols`): `schedule_name`, `device_id`, `start`, `event`, `sunday`, `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, `status` (default `status:1`)
- **dupCheck**: `schedule_name` ซ้ำ → 422
- Response: แถวที่ insert

### GET `/settings/listschedulepage`
- Query ฟิลเตอร์: `keyword` (LIKE `schedule_name`), `schedule_id`, `device_id`, `event`, `monday`..`sunday`, `start`, `createddate`, `updateddate`, `status`
- Response แต่ละแถว: `schedule_id`, `schedule_name`, `device_id`, `start`, `event`, `sunday`..`saturday`, `status`, `createddate`, `updateddate`, `device_name`, `countDevice`

### POST `/settings/updateschedulestatus`
- Body: ต้องมี `schedule_id`; อัปเดตเฉพาะ `status`
- Response: `{"affected": n}`

### POST `/settings/updatescheduledaystatus`
- Body: ต้องมี `schedule_id` + อย่างน้อย 1 ใน `sunday`, `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday` (ไม่มี day field → 400 `"no day field provided"`)
- Response: `{"affected": n}`

### POST `/settings/updateschedule`
- Body: ต้องมี `schedule_id`; อัปเดตจาก `scheduleCols` + `updateddate` อัตโนมัติ

## 10. Integration: MQTT

| # | Method | Path | Handler |
|---|--------|------|---------|
| 72 | GET | `/settings/lismqtt` | `ListMqtt` |
| 73 | GET | `/settings/listmqtt` | `ListMqttAlt` |
| 74 | GET | `/settings/lismqttall` | `MqttAll` |
| 75 | GET | `/settings/getmqttdetail` | `GetMqttDetail` |
| 76 | GET | `/settings/listmqttpaginate` | `ListMqttPaginate` |
| 77 | GET | `/settings/listmqttpaginateactive` | `ListMqttPaginateActive` |
| 78 | GET | `/settings/listmqttdevicepaginate` | `ListMqttDevicePaginate` |
| 79 | GET | `/settings/mqttdelete` | `DeleteMqtt` |
| 80 | GET | `/settings/deletemqtt` | `DeleteMqttAlt` |
| 81 | POST | `/settings/createmqtt` | `CreateMqtt` |
| 82 | POST | `/settings/updatemqtt` | `UpdateMqtt` |
| 83 | POST | `/settings/updatemqttstatus` | `UpdateMqttStatus` |
| 84 | POST | `/settings/mqtttsort` | `UpdateMqtttSort` |

### POST `/settings/createmqtt`
- Body whitelist (`mqttCols`): `mqtt_type_id`, `sort`, `mqtt_name`, `host`, `port`, `username`, `password`, `secret`, `expire_in`, `token_value`, `org`, `bucket`, `envavorment`, `location_id`, `latitude`, `longitude`, `mqtt_main_id`, `configuration`, `zoom`, `status`
- defaults: `status:1`, `sort:1`
- **dupCheck**: `mqtt_name` และ `bucket` ซ้ำ → 422

### POST `/settings/updatemqtt`
- Body ต้องมี `mqtt_id`; อัปเดตจาก `mqttCols`
- Response: `{"data":{"affected":<n>},"is_success":true}`

### POST `/settings/updatemqttstatus`
- Body: `mqtt_id`, `status` (int)
- ถ้าไม่เจอแถว → 404 `"No devices found with bucket '<id>'"`
- Response: `{"affected": n}`

### GET `/settings/getmqttdetail`
- Query: `mqtt_id` (required); ถ้าไม่เจอ → 404 `"mqtt not found"`
- Response fields: `mqtt_id`, `location_id`, `idhost`, `sort`, `mqtt_type_id`, `mqtt_name`, `host`, `port`, `username`, `password`, `secret`, `expire_in`, `token_value`, `org`, `bucket`, `envavorment`, `updateddate`, `status`, `latitude`, `longitude`, `zoom`, `mqtt_main_id`, `configuration`, `type_name`, `location_name`, `ipaddress`, `location_detail`

### GET `/settings/listmqttpaginateactive`
- fixed filter `m.status = 1`

## 11. Integration: MQTTHost

| # | Method | Path | Handler |
|---|--------|------|---------|
| 85 | GET | `/settings/listmqtthost` | `ListMqttHost` |
| 86 | GET | `/settings/mqtthostall` | `MqttHostAll` |
| 87 | POST | `/settings/createmqtthost` | `CreateMqttHost` |
| 88 | POST | `/settings/updatemqtthost` | `UpdateMqttHost` |
| 89 | POST | `/settings/updatemqtthoststatus` | `UpdateMqttHostStatus` |
| 90 | GET | `/settings/deletemqtthost` | `DeleteMqttHost` |

## 12. Integration: API

| # | Method | Path | Handler |
|---|--------|------|---------|
| 91 | GET | `/settings/apiall` | `ApiAll` |
| 92 | GET | `/settings/listapipage` | `ListApiPage` |
| 93 | POST | `/settings/createapi` | `CreateApi` |
| 94 | POST | `/settings/updateapi` | `UpdateApi` |
| 95 | GET | `/settings/deleteapi` | `DeleteApi` |

## 13. Integration: Email

| # | Method | Path | Handler |
|---|--------|------|---------|
| 96 | GET | `/settings/listemail` | `ListEmail` |
| 97 | GET | `/settings/emailall` | `EmailAll` |
| 98 | POST | `/settings/createemail` | `CreateEmail` |
| 99 | POST | `/settings/updateemail` | `UpdateEmail` |
| 100 | POST | `/settings/updateemailstatus` | `UpdateEmailStatus` |
| 101 | GET | `/settings/deleteemail` | `DeleteEmail` |

### POST `/settings/updateemailstatus`
- Body: `email_id`, `status` — `resetAll=true`: ถ้า status=1 จะ set status=0 ให้ทุกแถวอื่นก่อน

## 14. Integration: Host

| # | Method | Path | Handler |
|---|--------|------|---------|
| 102 | GET | `/settings/hostall` | `HostAll` |
| 103 | GET | `/settings/listhostpage` | `ListHostPage` |
| 104 | POST | `/settings/createhost` | `CreateHost` |
| 105 | POST | `/settings/updatehost` | `UpdateHost` |
| 106 | GET | `/settings/deletehost` | `DeleteHost` |

## 15. Integration: InfluxDB

| # | Method | Path | Handler |
|---|--------|------|---------|
| 107 | GET | `/settings/influxdball` | `InfluxdbAll` |
| 108 | GET | `/settings/listinfluxdbpage` | `ListInfluxdbPage` |
| 109 | POST | `/settings/createinfluxdb` | `CreateInfluxdb` |
| 110 | POST | `/settings/updateinfluxdb` | `UpdateInfluxdb` |
| 111 | POST | `/settings/updateinfluxdbstatus` | `UpdateInfluxdbStatus` |
| 112 | GET | `/settings/deleteinfluxdb` | `DeleteInfluxdb` |

## 16. Integration: Line

| # | Method | Path | Handler |
|---|--------|------|---------|
| 113 | GET | `/settings/lineall` | `LineAll` |
| 114 | GET | `/settings/listlinepage` | `ListLinePage` |
| 115 | POST | `/settings/createline` | `CreateLine` |
| 116 | POST | `/settings/updateline` | `UpdateLine` |
| 117 | POST | `/settings/updatelinestatus` | `UpdateLineStatus` |
| 118 | GET | `/settings/deleteline` | `DeleteLine` |

## 17. Integration: NodeRED

| # | Method | Path | Handler |
|---|--------|------|---------|
| 119 | GET | `/settings/noderedall` | `NoderedAll` |
| 120 | GET | `/settings/listnoderedpaginate` | `ListNoderedPaginate` |
| 121 | POST | `/settings/createnodered` | `CreateNodered` |
| 122 | POST | `/settings/updatenodered` | `UpdateNodered` |
| 123 | POST | `/settings/updatenoderedstatus` | `UpdateNoderedStatus` |
| 124 | GET | `/settings/deletenodered` | `DeleteNodered` |

## 18. Integration: SMS

| # | Method | Path | Handler |
|---|--------|------|---------|
| 125 | GET | `/settings/smsall` | `SmsAll` |
| 126 | GET | `/settings/listsmspage` | `ListSmsPage` |
| 127 | POST | `/settings/createsms` | `CreateSms` |
| 128 | POST | `/settings/updatesms` | `UpdateSms` |
| 129 | POST | `/settings/updatesmsstatus` | `UpdateSmsStatus` |
| 130 | GET | `/settings/deletesms` | `DeleteSms` |

## 19. Integration: Token

| # | Method | Path | Handler |
|---|--------|------|---------|
| 131 | GET | `/settings/tokenall` | `TokenAll` |
| 132 | GET | `/settings/tokensmspage` | `ListTokenPage` |
| 133 | POST | `/settings/createtoken` | `CreateToken` |
| 134 | POST | `/settings/updatetoken` | `UpdateToken` |
| 135 | GET | `/settings/deletetoken` | `DeleteToken` |

## 20. Integration: Telegram

| # | Method | Path | Handler |
|---|--------|------|---------|
| 136 | POST | `/settings/createtelegram` | `CreateTelegram` |
| 137 | POST | `/settings/updatetelegram` | `UpdateTelegram` |
| 138 | GET | `/settings/deletetelegram` | `DeleteTelegram` |

## 21. Dashboard Config

| # | Method | Path | Handler |
|---|--------|------|---------|
| 139 | POST | `/settings/dashboardconfig` | `CreateDashboardConfig` |
| 140 | GET | `/settings/dashboardconfig_1` | `DashboardConfigByLocation` |
| 141 | GET | `/settings/dashboardconfig` | `ListDashboardConfig` |
| 142 | GET | `/settings/dashboardconfig/search` | `FindOrCreateDashboardConfig` |
| 143 | GET | `/settings/dashboardconfig/{id}` | `GetDashboardConfig` |
| 144 | PATCH | `/settings/dashboardconfig/{id}` | `UpdateDashboardConfig` |
| 145 | DELETE | `/settings/dashboardconfig/{id}` | `RemoveDashboardConfig` |

### POST `/settings/dashboardconfig`
- Body: ต้องมี `name`, `location_id`, `config` (ขาด → 400 `"name, location_id and config are required"`) — เป็น **upsert** keyed on `(name, location_id)`
- Response: แถว dashboard config

### GET `/settings/dashboardconfig/search`
- Query: `name` (required), `location_id` (required); ขาด → 422
- ถ้าไม่มีแถว → insert ใหม่ด้วย `config_data:{}`, `status:1`

## 22. Alarm Device / Event links

| # | Method | Path | Handler |
|---|--------|------|---------|
| 146 | GET | `/settings/listalarmdevicepage` | `ListAlarmDevicePage` |
| 147 | GET | `/settings/listalarmdeviceactivepage` | `ListAlarmDeviceActivePage` |
| 148 | GET | `/settings/listalarmeventdevicepage` | `ListAlarmEventDevicePage` |
| 149 | GET | `/settings/listalarmeventdevicecontrolpage` | `ListAlarmEventDeviceControlPage` |
| 150 | GET | `/settings/activealarmdevicepage` | `ActiveAlarmDevicePage` |
| 151 | GET | `/settings/activealarmeventdeviceeventpage` | `ActiveAlarmEventDeviceEventPage` |
| 152 | GET | `/settings/alarmdevice` | `AlarmDevice` |
| 153 | GET | `/settings/alarmdevicestatus` | `AlarmDeviceStatus` |
| 154 | GET | `/settings/deviceactivemqttalarm` | `DeviceActiveMqttAlarm` |
| 155 | GET | `/settings/devicealarm` | `DeviceAlarm` |
| 156 | GET | `/settings/createalarmdevice` | `CreateAlarmDeviceViaGet` |
| 157 | GET | `/settings/deletealarmdevice` | `DeleteAlarmDevices` |
| 158 | GET | `/settings/createalarmeventdevice` | `CreateAlarmEventDeviceViaGet` |
| 159 | GET | `/settings/deletealarmeventdevice` | `DeleteAlarmEventDevices` |
| 160 | GET | `/settings/deletearmdevice` | `DeleteArmDevice` |
| 161 | GET | `/settings/deletearmdevicev2` | `DeleteArmDeviceV2` |
| 162 | POST | `/settings/createalarmDevice` | `CreateAlarmDevice` |
| 163 | POST | `/settings/createalarmdevicepaginate` | `CreateAlarmDevicePaginate` |
| 164 | POST | `/settings/updatealarmdevice` | `UpdateAlarmDevice` |
| 165 | POST | `/settings/updatealarmstatus` | `UpdateAlarmStatus` |
| 166 | POST | `/settings/createdevicealarmaction` | `CreateDeviceAlarmAction` |

### GET `/settings/alarmdevicestatus`
- Query: `device_id` (required → 400 `"device_id is required"`; ไม่พบ → 404 `"device not found"`)
- Response: แถวเดียว `DeviceSpec`

### POST `/settings/createalarmdevicepaginate`
- Body: ต้องมี `sn` (ขาด → 400 `"sn is required"`); sn ซ้ำ → 422
- Body whitelist (`settingCols`): `location_id`, `setting_type_id`, `setting_name`, `sn`, `status` (default `status:1`)

## 23. Device Alarm Monitors

| # | Method | Path | Handler |
|---|--------|------|---------|
| 167 | GET | `/settings/listdevicealarm` | `ListDeviceAlarm` |
| 168 | GET | `/settings/listdevicealarmairV1` | `ListDeviceAlarmAirV1` |
| 169 | GET | `/settings/listdevicealarmair` | `ListDeviceAlarmAir` |
| 170 | GET | `/settings/listdevicealarmall` | `ListDeviceAlarmAll` |
| 171 | GET | `/settings/_listdevicealarmfan` | `UnderscoreListDeviceAlarmFan` |
| 172 | GET | `/settings/listdevicealarmfan` | `ListDeviceAlarmFan` |
| 173 | GET | `/settings/listdevicealarmlimit` | `ListDeviceAlarmLimit` |
| 174 | GET | `/settings/_devicemonitor` | `UnderscoreDeviceMonitor` |
| 175 | GET | `/settings/devicemonitor` | `DeviceMonitor` |
| 176 | GET | `/settings/devicemonitors` | `DeviceMonitors` |

## 24. Process Logs

| # | Method | Path | Handler |
|---|--------|------|---------|
| 177 | GET | `/settings/scheduleproces` | `ScheduleProces` |
| 178 | GET | `/settings/scheduleprocesslog` | `ScheduleProcessLog` |
| 179 | GET | `/settings/scheduleprocesslogpaginate` | `ScheduleProcessLogPaginate` |
| 180 | GET | `/settings/mqtterrorlogpaginate` | `MqttErrorLogPaginate` |
| 181 | GET | `/settings/alarmlogpaginate` | `AlarmLogPaginate` |
| 182 | GET | `/settings/alarmlogpaginateemail` | `AlarmLogPaginateEmail` |
| 183 | GET | `/settings/alarmlogpaginateline` | `AlarmLogPaginateLine` |
| 184 | GET | `/settings/alarmlogpaginatesms` | `AlarmLogPaginateSms` |
| 185 | GET | `/settings/alarmlogpaginatetelegram` | `AlarmLogPaginateTelegram` |
| 186 | GET | `/settings/alarmlogpaginatecontrols` | `AlarmLogPaginateControls` |
| 187 | GET | `/settings/alarmlogpaginatecontrol` | `AlarmLogPaginateControl` |

---

## หมายเหตุ
- **Delete** ส่วนใหญ่เป็น `GET` (mirror NestJS) ยกเว้น `deletesetting`, `dashboardconfig/{id}`
- id ใน query/path ใช้ `coerceID` — ถ้าเป็นตัวเลข cast เป็น int
- path หลัง `/api/settings`; path-sensitive (เช่น `createalarmDevice` ตัวพิมพ์ใหญ่ D)
- มี path พร้อม underscore นำหน้า: `_listdevicealarmfan`, `_devicemonitor`
