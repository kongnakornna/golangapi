บันทึกความคืบหน้าลงในเอกสารแผนแล้ว (Section 15 — Progress Log)
สรุปที่บันทึกไว้:
- ✅ เสร็จ: Phase 0 (skeleton 4 module + wiring), Phase 2 (realtime dashboard push + tests), Phase 3 (notifier จริง email/sms/line/discord/io + tests)
- ⏹️ ข้าม: Phase 3b (auto-hook alarm engine) — ตามที่เลือก option 1 ไม่แตะของเดิม
- ⏳ เหลือ: Phase 4–8
- ไฟล์ที่แตะ: internal/modules/{notifier,realtime,auditlog,apimanager,flowengine}/, internal/server/handlers.go, docs/plan/2026-08-28-applied-iot-platform-plan.md
- หมายเหตุ: restore docs/docs.go จาก git (เสียมาก่อน ไม่ใช่จากงานนี้); config test ล้มเพราะ env MQTT_BROKER ใน shell (ไม่เกี่ยวกับโค้ด)
วันหลังเปิด docs/plan/2026-08-28-applied-iot-platform-plan.md → Section 15 มีวิธีทำต่อครบ พร้อมเบาะแสว่า Phase 7 ต้อง wire notifier + realtime เข้า flow engine




1.C:\github\icmongolang\docs\InfluxDBServiceModuleV2.md
2.C:\github\icmongolang\docs\InfluxDBServiceModule.md
นำมา ปรับใช้กับ ของเดิม  โดยต้องไม่กระทบการทำงาน เดิม
ให้ ทำเอกสาร แผนงานออกมาก่อน
ต้องการ คือ
1.ระบบ settings
2.ระบบแสดงข้อมูล Read time dashboard
3.


time.Now().In(GetTimeLocation()).Format


TimeLoc := helpers.GetTimeLocation()
"time":        helpers.TimeConvertermas(time.Now().In(TimeLoc)),
timePoints = append(timePoints, helpers.TimeConvertermas(t.In(TimeLoc)))

tt := t.(time.Time).In(TimeLoc)
timeStr = helpers.TimeConvertermas(tt) // match sample format
dateStrings = append(dateStrings, timeStr)

NewLocalTime

CreatedAt
UpdatedAt,
เพิ่ม
time.Now().In(GetTimeLocation()).Format


TimeLoc := helpers.GetTimeLocation()
"time":        helpers.TimeConvertermas(time.Now().In(TimeLoc)),
timePoints = append(timePoints, helpers.TimeConvertermas(t.In(TimeLoc)))

tt := t.(time.Time).In(TimeLoc)
timeStr = helpers.TimeConvertermas(tt) // match sample format
dateStrings = append(dateStrings, timeStr)



go mod tidy
go mod download
go mod verify
go run cmd/api/main.go migrate
swag init -g cmd/api/main.go 
go clean -cache        
go clean -modcache
go mod vendor 
go test ./...  
air
./air.cmd

-----------------------------------

.\run.ps1 -All

-----------------------------------

อ่าน  C:\github\sd_iot_schedule.sql
อ่าน C:\github\fs_schedule.sql

update for fs_schedule
 ใช้ 
 "event" int8,
  "sunday" int8,
  "monday" int8,
  "tuesday" int8,
  "wednesday" int8,
  "thursday" int8,
  "friday" int8,
  "saturday" int8,
  "status" int8
  
อ่านจาก  C:\github\sd_iot_schedule.sql

// GetAlarmDeviceStatus
func (u *mqtt3UseCase) GetAlarmDeviceStatus(ctx context.Context, req map[string]interface{}) (interface{}, error) {
	// โหลด Bangkok timezone (UTC+7)
	TimeLoc := helpers.GetTimeLocation()
	// 1. Extract bucket (required)
	bucket := ""
	if v, ok := req["bucket"]; ok {
		bucket = fmt.Sprintf("%v", v)
	}
	if bucket == "" {
		return nil, fmt.Errorf("bucket is required")
	}

	page := 1
	if v, ok := req["page"]; ok {
		if p, err := strconv.Atoi(fmt.Sprintf("%v", v)); err == nil && p > 0 {
			page = p
		}
	}
	if page == 0 {
		return nil, fmt.Errorf("page is required")
	}

	pageSize := 1000
	if v, ok := req["pageSize"]; ok {
		if ps, err := strconv.Atoi(fmt.Sprintf("%v", v)); err == nil && ps > 0 {
			pageSize = ps
		}
	}
	if pageSize == 0 {
		return nil, fmt.Errorf("pageSize is required")
	}
	lang := "en"
	if v, ok := req["lang"]; ok {
		lang = fmt.Sprintf("%v", v)
	}
	measurement := ""
	if v, ok := req["measurement"]; ok {
		measurement = fmt.Sprintf("%v", v)
	}
	if measurement == "" {
		//return nil, fmt.Errorf("measurement is required")
		measurement = "temperature" // default
	}

	// 2. Prepare device list request with sensible defaults
	listReq := &repository.DeviceListAlarmRequest{
		Page:     page,
		PageSize: pageSize,
		Status:   1,
		Bucket:   bucket,
	}
	if page, ok := req["page"]; ok {
		if p, err := strconv.Atoi(fmt.Sprintf("%v", page)); err == nil && p > 0 {
			listReq.Page = p
		}
	}
	if pageSize, ok := req["pageSize"]; ok {
		if ps, err := strconv.Atoi(fmt.Sprintf("%v", pageSize)); err == nil && ps > 0 {
			listReq.PageSize = ps
		}
	}
	// Optionally accept other filters
	if v, ok := req["device_id"]; ok {
		listReq.DeviceID = fmt.Sprintf("%v", v)
	}
	if v, ok := req["type_id"]; ok {
		listReq.TypeID = toInt(v)
	}
	if v, ok := req["hardware_id"]; ok {
		listReq.HardwareID = toInt(v)
	}
	if v, ok := req["keyword"]; ok {
		listReq.Keyword = fmt.Sprintf("%v", v)
	}

	// 3. Fetch devices from repository
	//
	deviceResult, err := u.deviceRepo.ListDevicesWithAlarm(ctx, listReq)
	if err != nil {
		return nil, err
	}

	// 4. ดึงข้อมูล MQTT ล่าสุดแบบยืดหยุ่น (ย้ายมาอยู่ก่อน grouping เพื่อให้
	//    deviceioinfo / devicesensor แสดง "ค่าจริง" จาก MQTT ไม่ใช่ topic)
	mqttConnected := u.mqttClient.IsConnected()
	var mqttRawPayload string
	var fromCache bool
	var cacheTime int64
	cacheEnabled := u.IsCacheEnabled()
	var mqttDataMap map[string]interface{} // map ชื่อฟิลด์ -> ค่า (string หรือ number)

	if mqttConnected && len(deviceResult.Items) > 0 {
		firstDevice := deviceResult.Items[0]
		topic := firstDevice.MqttDataValue
		if topic == "" {
			topic = bucket + "/DATA"
		}

		mqttCacheKey := "mqtt_payload:" + bucket // ใช้ key เดียวกันกับ GetMonitorDeviceGroup

		// 1) ลองอ่านจาก Redis cache
		if cacheEnabled {
			cached, _ := u.redisClient.Get(ctx, mqttCacheKey).Bytes()
			if len(cached) > 0 {
				mqttRawPayload = string(cached)
				fromCache = true
				cacheTime = 15 // โดยประมาณ
				u.logger.Debug("MQTT payload from cache")
			}
		}

		// 2) ถ้าไม่เจอใน cache หรือ cache ว่าง ให้ดึงจาก MQTT
		if mqttRawPayload == "" {
			mqttCtx, mqttCancel := context.WithTimeout(ctx, 5*time.Second)
			defer mqttCancel()

			payload, err := u.mqttClient.GetDataFromTopic(mqttCtx, topic, 5*time.Second)
			if err == nil && len(payload) > 0 {
				mqttRawPayload = string(payload)
				fromCache = false
				cacheTime = 0
				if cacheEnabled {
					_ = u.redisClient.Set(ctx, mqttCacheKey, payload, 60*time.Second).Err()
				}
			} else {
				u.logger.Warnf("Failed to get data from MQTT topic %s: %v", topic, err)
				// ถ้าดึงไม่ได้ แต่มี cache เก่า (ซึ่งอาจจะยังอยู่) เราก็จะใช้ cache ที่ได้ไว้แล้ว
				// แต่ถ้า cache ก็ไม่มี mqttRawPayload จะเป็น "" และจะแสดง "No data"
			}
		}

		// 3) ถ้ามี payload ให้แยกส่วนและสร้าง mqttDataMap
		if mqttRawPayload != "" {
			parts := strings.Split(mqttRawPayload, ",")
			var configMap map[string]string
			if firstDevice.MqttStatusDataName != "" {
				if err := json.Unmarshal([]byte(firstDevice.MqttStatusDataName), &configMap); err != nil {
					u.logger.Warnf("Failed to parse MqttStatusDataName: %v", err)
					configMap = nil
				}
			}
			mqttDataMap = make(map[string]interface{}, len(parts))
			for i, val := range parts {
				key := fmt.Sprintf("%d", i)
				if configMap != nil {
					if mapped, ok := configMap[key]; ok {
						key = mapped
					}
				}
				// ตัดช่องว่าง และพยายามแปลงเป็นตัวเลขถ้าเป็นไปได้ (เพื่อให้ response มีข้อมูลที่ใช้งานง่าย)
				trimmed := strings.TrimSpace(val)
				if f, err := strconv.ParseFloat(trimmed, 64); err == nil {
					mqttDataMap[key] = f
				} else {
					mqttDataMap[key] = trimmed
				}
			}
		}
	}

	// 5. Group devices by hardware_id
	var deviceSensors   []interface{} // hardware_id == 1
	var deviceIO        []interface{} // hardware_id == 2
	var deviceControl   []interface{} // hardware_id == 3
	var deviceCritical  []interface{} // hardware_id == 4
	var deviceIOInfo    []interface{} // simplified version for "deviceioinfo"
	var deviceArray     []map[string]interface{} // rss_data-style array (mqtt2/devicemqtt)

	// baseUrl for building control/sensercharts URLs
	baseUrl := ""
	if u.cfg != nil {
		baseUrl = strings.TrimSuffix(u.cfg.Server.BaseUrl, "/")
	}

	for _, dev := range deviceResult.Items {
		devMap := structToMap(dev)
		switch dev.HardwareID {
		case 1:
			deviceSensors = append(deviceSensors, devMap)
		case 2:
			deviceIO = append(deviceIO, devMap)
		case 3:
			deviceControl = append(deviceControl, devMap)
		case 4:
			deviceCritical = append(deviceCritical, devMap)
		default:
			// Unknown hardware_id -> also expose as sensor group to avoid data loss
			deviceSensors = append(deviceSensors, devMap)
		}
		// deviceioinfo: simpler representation
		// แสดง "ค่าจริง" จาก MQTT payload (lookup ด้วย measurement -> mqttDataMap)
		// แทน topic เดิม ช่วยให้ dataAlarm/eventControl สะท้อนสถานะจริงของ device
		realValue := "0"
		if mqttDataMap != nil {
			if v, ok := mqttDataMap[dev.Measurement]; ok {
				realValue = fmt.Sprintf("%v", v)
			} else if v, ok := mqttDataMap[dev.MqttDeviceName]; ok {
				realValue = fmt.Sprintf("%v", v)
			}
		}
		dataAlarm := 0
		if f, err := strconv.ParseFloat(realValue, 64); err == nil {
			if f >= 1 {
				dataAlarm = 1
			}
		}
		ioInfo := map[string]interface{}{
			"device_id":      dev.DeviceID,
			"type_id":        dev.TypeID,
			"status":         dev.Status,
			"device_name":    dev.DeviceName,
			"timestamp":      helpers.GetCurrentFullDatenow(),
			"subject":        dev.StatusWarning,
			"value_data":     realValue,
			"dataAlarm":      dataAlarm,
			"eventControl":   1,
			"value_data_msg": realValue,
		}
		deviceIOInfo = append(deviceIOInfo, ioInfo)

		// rss_data-style enriched item (mqtt2/devicemqtt output)
		deviceArray = append(deviceArray, u.buildDeviceMqttItem(dev, mqttDataMap, fromCache, lang, baseUrl, 60))
	}

	// 5. MQTT connection status
	checkConnectionMqtt := map[string]interface{}{
		"isConnected": mqttConnected,
		"connected":   mqttConnected,
		"status":      1,
		"msg":         "MQTT Connection Status: Connected",
	}
	if !mqttConnected {
		checkConnectionMqtt["status"] = 0
		checkConnectionMqtt["msg"] = "MQTT Connection Status: Disconnected"
	}

	// ====== 6. สร้าง mqttrs และ mqttData สำหรับ response ======
	mqttrs := map[string]interface{}{
		"case":        0,
		"status":      0,
		"msg":         "No data available",
		"fromCache":   false,
		"time":        0,
		"timestamp":   helpers.GetCurrentFullDatenow(),
		"isConnected": mqttConnected,
	}
	mqttData := make(map[string]interface{})

	if mqttRawPayload != "" && mqttDataMap != nil {
		mqttrs = map[string]interface{}{
			"case":        1,
			"status":      1,
			"msg":         mqttRawPayload, // ข้อมูลดิบ
			"fromCache":   fromCache,
			"time":        cacheTime,
			"timestamp":   helpers.GetCurrentFullDatenow(),
			"isConnected": mqttConnected,
		}
		// คัดลอกข้อมูลที่แยกแล้ว (อาจมีทั้งตัวเลขและข้อความ)
		for k, v := range mqttDataMap {
			mqttData[k] = v
		}
	}

	// 7. Chart data from InfluxDB
	chartData := map[string]interface{}{
		"bucket": bucket,
		"field":  "value",
		"info":   map[string]interface{}{},
		"data":   []float64{},
		"date":   []string{},
		"name":   "value",
		"cache":  "no cache",
	}
	if u.influxClient != nil {
		now := time.Now()
		start := now.Add(-15 * time.Minute).Format(time.RFC3339)
		stop := now.Format(time.RFC3339)
		params := influxdb.QueryParams{
			Measurement: measurement, // could be made configurable per device
			Field:       "value",
			Bucket:      bucket,
			Start:       start,
			Stop:        stop,
			Limit:       150,
		}
		results, err := u.influxClient.QueryFilterData(params)
		if err == nil {
			var dataPoints []float64
			var timePoints []string
			for _, r := range results {
				if val, ok := r["_value"].(float64); ok {
					dataPoints = append(dataPoints, val)
				}
				if t, ok := r["_time"].(time.Time); ok {
					timePoints = append(timePoints, helpers.TimeConvertermas(t.In(TimeLoc)))
				}
			}
			chartData["data"] = dataPoints
			chartData["date"] = timePoints
			chartData["info"] = map[string]interface{}{
				"bucket":      bucket,helpers.TimeConvertermas(t.In(
				"measurement": measurement,
				"result":      "last",
				"table":       0,
				"field":       "value",
				"start":       start,
				"stop":        stop,
				"time":        helpers.TimeConvertermas(time.Now().In(TimeLoc)),
				"value": func() interface{} {
					if len(dataPoints) > 0 {
						return dataPoints[len(dataPoints)-1]
					}
					return nil
				}(),
			}
		} else {
			u.logger.Warnf("Failed to query InfluxDB: %v", err)
		}
	}

	// 8. Build final response
	response := map[string]interface{}{
		"statuscode": 200,
		"status":     "success",
		"Mqttstatus": 1,
		"payload": map[string]interface{}{
			"checkConnectionMqtt": checkConnectionMqtt,
			"mqttrs":              mqttrs,
			"mqttname":            getMqttNameFromDevices(deviceResult.Items),
			"bucket":              bucket,
			"time":                helpers.GetCurrentFullDatenow(),
			"mqttdata":            mqttData,
			"deviceioinfo":        deviceIOInfo,
			"devicesensor":        deviceSensors,
			"deviceio":            deviceIO,
			"devicecritical":      deviceCritical, 
			"cache": "cache",
			"chart": chartData,
			// mqtt2/devicemqtt-compatible enriched output below
			"lang":         lang,
			"page":         page,
			"currentPage":  page,
			"pageSize":     pageSize,
			"total":        int64(len(deviceResult.Items)),
			"device_count": len(deviceArray),
			"device":       deviceArray,
		},
		"message":    "check Connection Status Mqtt",
		"message_th": "check Connection Status Mqtt",
	}
	return response, nil
}



แก้   SLOW SQL


2026/08/30 00:27:21 C:/github/icmongolang/cmd/migrate.go:281 SLOW SQL >= 200ms
[291.003ms] [rows:0] ALTER TABLE "sd_iot_device" ALTER COLUMN "device_name" TYPE varchar(255) USING "device_name"::varchar(255)

2026/08/30 00:27:21 C:/github/icmongolang/cmd/migrate.go:281 SLOW SQL >= 200ms
[273.330ms] [rows:0] ALTER TABLE "sd_iot_device" ALTER COLUMN "sn" TYPE varchar(255) USING "sn"::varchar(255)

2026/08/30 00:27:21 C:/github/icmongolang/cmd/migrate.go:281 SLOW SQL >= 200ms
[239.711ms] [rows:0] ALTER TABLE "sd_iot_device" ALTER COLUMN "status_warning" TYPE varchar(150) USING "status_warning"::varchar(150)

2026/08/30 00:27:21 C:/github/icmongolang/cmd/migrate.go:281 SLOW SQL >= 200ms
[218.787ms] [rows:0] ALTER TABLE "sd_iot_device" ALTER COLUMN "recovery_warning" TYPE varchar(150) USING "recovery_warning"::varchar(150)

2026/08/30 00:27:22 C:/github/icmongolang/cmd/migrate.go:281 SLOW SQL >= 200ms
[217.449ms] [rows:0] ALTER TABLE "sd_iot_device" ALTER COLUMN "status_alert" TYPE varchar(150) USING "status_alert"::varchar(150)

2026/08/30 00:27:22 C:/github/icmongolang/cmd/migrate.go:281 SLOW SQL >= 200ms
[214.530ms] [rows:0] ALTER TABLE "sd_iot_device" ALTER COLUMN "recovery_alert" TYPE varchar(150) USING "recovery_alert"::varchar(150)

2026/08/30 00:27:22 C:/github/icmongolang/cmd/migrate.go:281 SLOW SQL >= 200ms
[207.864ms] [rows:0] ALTER TABLE "sd_iot_device" ALTER COLUMN "period" TYPE varchar(150) USING "period"::varchar(150)

2026/08/30 00:27:23 C:/github/icmongolang/cmd/migrate.go:281 SLOW SQL >= 200ms
[217.400ms] [rows:0] ALTER TABLE "sd_iot_device" ALTER COLUMN "min" TYPE varchar(255) USING "min"::varchar(255)

2026/08/30 00:27:24 C:/github/icmongolang/cmd/migrate.go:281 SLOW SQL >= 200ms
[206.652ms] [rows:0] ALTER TABLE "sd_iot_device" ALTER COLUMN "unit" TYPE varchar(255) USING "unit"::varchar(255)

2026/08/30 00:27:24 C:/github/icmongolang/cmd/migrate.go:281 SLOW SQL >= 200ms
[210.640ms] [rows:0] ALTER TABLE "sd_iot_device" ALTER COLUMN "mqtt_data_value" TYPE varchar(255) USING "mqtt_data_value"::varchar(255)

2026/08/30 00:27:25 C:/github/icmongolang/cmd/migrate.go:281 SLOW SQL >= 200ms
[221.204ms] [rows:0] ALTER TABLE "sd_iot_device" ALTER COLUMN "mqtt_data_control" TYPE varchar(255) USING "mqtt_data_control"::varchar(255)

2026/08/30 00:27:25 C:/github/icmongolang/cmd/migrate.go:281 SLOW SQL >= 200ms
[205.428ms] [rows:0] ALTER TABLE "sd_iot_device" ALTER COLUMN "measurement" TYPE varchar(255) USING "measurement"::varchar(255)

2026/08/30 00:27:25 C:/github/icmongolang/cmd/migrate.go:281 SLOW SQL >= 200ms
[232.407ms] [rows:0] ALTER TABLE "sd_iot_device" ALTER COLUMN "org" TYPE varchar(255) USING "org"::varchar(255)

2026/08/30 00:27:25 C:/github/icmongolang/cmd/migrate.go:281 SLOW SQL >= 200ms
[202.185ms] [rows:0] ALTER TABLE "sd_iot_device" ALTER COLUMN "bucket" TYPE varchar(255) USING "bucket"::varchar(255)

2026/08/30 00:27:28  ERROR: constraint "uni_device_status_device_id" of relation "device_status" does not exist (SQLSTATE 42704)
[1.509ms] [rows:0] ALTER TABLE "device_status" DROP CONSTRAINT "uni_device_status_device_id"
2026-08-30T00:27:28.027+0700    WARN    cmd/migrate.go:287      Skipping non-fatal migrate error on icmongolang/internal/modules/iot/models.DeviceStatus (constraint/index not found): ERROR: constraint "uni_device_status_device_id" of relation "device_status" does not exist (SQLSTATE 42704)

2026/08/30 00:27:28  ERROR: constraint "uni_device_config_device_id" of relation "device_config" does not exist (SQLSTATE 42704)

สร้าง demodata  ลง databse  fs_schedule 10 rows    map data จาก sd_iot_device 

