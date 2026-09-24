
  influxdb:
    image: influxdb:2.7
    ports:
      - 9087:8086
    environment:
      - DOCKER_INFLUXDB_INIT_MODE=setup
      - DOCKER_INFLUXDB_INIT_USERNAME=admin
      - DOCKER_INFLUXDB_INIT_PASSWORD=${INFLUXDB_PASSWORD:-admin}
      - DOCKER_INFLUXDB_INIT_ORG=cmon_org
      - DOCKER_INFLUXDB_INIT_BUCKET=AIRCOM1
      - DOCKER_INFLUXDB_INIT_ADMIN_TOKEN=${INFLUXDB_TOKEN:-icmongolang-super-secret-token}
    volumes:
      - app-influxdb-data:/var/lib/influxdb2

BUCKET


AIRCOM1
AIRCOM2
AIRCOM3
AIRCOM4
AIRCOM5
AIRCOM6
AIRCOM7
AIRCOM8
AIRCOM9
AIRCOM10
AIRCOM11
AIRCOM12
AIRCOM13
AIRCOM14
AIRCOM15
AIRCOM16
AIRCOM17
AIRCOM18
AIRCOM19
AIRCOM20
BAACTW01
BAACTW02
BAACTW03
BAACTW04
BAACTW05
BAACTW06
BAACTW07
BAACTW08
BAACTW09
BAACTW10
BAACTW11
BAACTW12
BAACTW13
BAACTW14
BAACTW15
BAACTW16
BAACTW17
BAACTW18
BAACTW19
BAACTW20
CMONBUGKET
CMONBUGKET1
CMONBUGKET2
CMONBUGKET3
CMONBUGKET4
CMONBUGKET5
CMONBUGKET6
CMONBUGKET7
CMONBUGKET8
CMONBUGKET9
CMONBUGKET10
SORACELL1
SORACELL2
SORACELL3
SORACELL4
SORACELL5
SORACELL6
SORACELL7
SORACELL8
SORACELL9
SORACELL10


		
อ่าน
 1.C:\github\gistdaapi\src\modules\mqtt2
 2.C:\github\gistdaapi\src\modules\mqtt2\mqtt2.service.ts
 3.C:\github\gistdaapi\src\modules\mqtt2\mqtt2.controller.ts
 
มาปรับใช้กับ
 C:\github\icmongolang\internal\modules\iot
 
 UPDATE  ต้องการให้แสกงผลเหมือน  C:\github\gistdaapi\src\modules\mqtt2\mqtt2.controller.ts


{{baseUrl}}/api/iot/monitordevicegroup?bucket=SORACELL1


{
    "data": {
        "bucket": "SORACELL1",
        "cache_used": true,
        "data": [
            {
                "count": 9,
                "devices": [
                    {
                        "alarm_status": 5,
                        "alarm_subject": "Normal",
                        "alarm_title": "Normal",
                        "cache_used": true,
                        "control": [],
                        "device_id": 84,
                        "device_name": "AMP",
                        "devicedata": "250.00 A",
                        "graph": "http://localhost:5000/api/iot/monitordevicechart?bucket=SORACELL1&measurement=amp1&field=value&start=-5m&stop=now()&limit=120&lang=en",
                        "hardware_id": 1,
                        "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"icon icon-tabler icons-tabler-outline icon-tabler-circuit-ammeter\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M5 12a7 7 0 1 0 14 0a7 7 0 1 0 -14 0\" /><path d=\"M5 12h-3\" /><path d=\"M19 12h3\" /><path d=\"M10 14v-3c0 -1.036 .895 -2 2 -2s2 .964 2 2v3\" /><path d=\"M14 12h-4\" /></svg>",
                        "icon_access": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"icon icon-tabler icons-tabler-outline icon-tabler-circuit-ammeter\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M5 12a7 7 0 1 0 14 0a7 7 0 1 0 -14 0\" /><path d=\"M5 12h-3\" /><path d=\"M19 12h3\" /><path d=\"M10 14v-3c0 -1.036 .895 -2 2 -2s2 .964 2 2v3\" /><path d=\"M14 12h-4\" /></svg>",
                        "icon_off": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_on": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"fan-spin windmill-spin me-2 text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "layout": 2,
                        "location_name": "Location Solar",
                        "measurement": "amp1",
                        "menu": 1,
                        "mqtt_connected": true,
                        "mqtt_control_off": "0",
                        "mqtt_control_on": "1",
                        "mqtt_data_control": "SORACELL1/CONTROL",
                        "mqtt_data_value": "SORACELL1/DATA",
                        "status": 1,
                        "timestamp": "2026-08-27 21:58:07",
                        "type_id": 8,
                        "type_name": "ZONE 1",
                        "unit": "A",
                        "value_data": "250.00"
                    },
                    {
                        "alarm_status": 5,
                        "alarm_subject": "Normal",
                        "alarm_title": "Normal",
                        "cache_used": true,
                        "control": [],
                        "device_id": 85,
                        "device_name": "Volt",
                        "devicedata": "300.00 V",
                        "graph": "http://localhost:5000/api/iot/monitordevicechart?bucket=SORACELL1&measurement=volt1&field=value&start=-5m&stop=now()&limit=120&lang=en",
                        "hardware_id": 1,
                        "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"icon icon-tabler icons-tabler-outline icon-tabler-circuit-voltmeter\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M5 12a7 7 0 1 0 14 0a7 7 0 1 0 -14 0\" /><path d=\"M5 12h-3\" /><path d=\"M19 12h3\" /><path d=\"M10 10l2 4l2 -4\" /></svg>",
                        "icon_access": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"icon icon-tabler icons-tabler-outline icon-tabler-circuit-voltmeter\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M5 12a7 7 0 1 0 14 0a7 7 0 1 0 -14 0\" /><path d=\"M5 12h-3\" /><path d=\"M19 12h3\" /><path d=\"M10 10l2 4l2 -4\" /></svg>",
                        "icon_off": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_on": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"fan-spin windmill-spin me-2 text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "layout": 2,
                        "location_name": "Location Solar",
                        "measurement": "volt1",
                        "menu": 1,
                        "mqtt_connected": true,
                        "mqtt_control_off": "0",
                        "mqtt_control_on": "1",
                        "mqtt_data_control": "SORACELL1/CONTROL",
                        "mqtt_data_value": "SORACELL1/DATA",
                        "status": 1,
                        "timestamp": "2026-08-27 21:58:07",
                        "type_id": 8,
                        "type_name": "ZONE 1",
                        "unit": "V",
                        "value_data": "300.00"
                    },
                    {
                        "alarm_status": 5,
                        "alarm_subject": "Normal",
                        "alarm_title": "Normal",
                        "cache_used": true,
                        "control": [],
                        "device_id": 86,
                        "device_name": "Watt",
                        "devicedata": "1.50 Kw",
                        "graph": "http://localhost:5000/api/iot/monitordevicechart?bucket=SORACELL1&measurement=watt1&field=value&start=-5m&stop=now()&limit=120&lang=en",
                        "hardware_id": 1,
                        "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"icon icon-tabler icons-tabler-outline icon-tabler-circuit-resistor\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M2 12h2l2 -5l3 10l3 -10l3 10l3 -10l1.5 5h2.5\" /></svg>",
                        "icon_access": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"icon icon-tabler icons-tabler-outline icon-tabler-circuit-resistor\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M2 12h2l2 -5l3 10l3 -10l3 10l3 -10l1.5 5h2.5\" /></svg>",
                        "icon_off": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_on": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"fan-spin windmill-spin me-2 text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "layout": 2,
                        "location_name": "Location Solar",
                        "measurement": "watt1",
                        "menu": 1,
                        "mqtt_connected": true,
                        "mqtt_control_off": "0",
                        "mqtt_control_on": "1",
                        "mqtt_data_control": "SORACELL1/CONTROL",
                        "mqtt_data_value": "SORACELL1/DATA",
                        "status": 1,
                        "timestamp": "2026-08-27 21:58:07",
                        "type_id": 8,
                        "type_name": "ZONE 1",
                        "unit": "Kw",
                        "value_data": "1.50"
                    },
                    {
                        "alarm_status": 5,
                        "alarm_subject": "Normal",
                        "alarm_title": "Normal",
                        "cache_used": true,
                        "control": [],
                        "device_id": 87,
                        "device_name": "EV",
                        "devicedata": "70.00 %",
                        "graph": "http://localhost:5000/api/iot/monitordevicechart?bucket=SORACELL1&measurement=ev1&field=value&start=-5m&stop=now()&limit=120&lang=en",
                        "hardware_id": 1,
                        "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"currentColor\" class=\"icon icon-tabler icons-tabler-filled icon-tabler-charging-pile\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M12 3a3 3 0 0 1 3 3v4a3 3 0 0 1 3 3v3a.5 .5 0 1 0 1 0v-6.585l-1 -1l-.293 .292a1 1 0 0 1 -1.414 -1.414l.292 -.293l-.292 -.293a1 1 0 0 1 -.083 -1.32l.083 -.094a1 1 0 0 1 1.414 0l3 3a1 1 0 0 1 .293 .707v7a2.5 2.5 0 1 1 -5 0v-3a1 1 0 0 0 -1 -1v7a1 1 0 0 1 0 2h-12a1 1 0 0 1 0 -2v-13a3 3 0 0 1 3 -3zm-2.486 7.643a1 1 0 0 0 -1.371 .343l-1.5 2.5l-.054 .1a1 1 0 0 0 .911 1.414h1.233l-.59 .986a1 1 0 0 0 1.714 1.028l1.5 -2.5l.054 -.1a1 1 0 0 0 -.911 -1.414h-1.235l.592 -.986a1 1 0 0 0 -.343 -1.371m2.486 -5.643h-6a1 1 0 0 0 -1 1v1h8v-1a1 1 0 0 0 -1 -1\" /></svg>",
                        "icon_access": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"currentColor\" class=\"icon icon-tabler icons-tabler-filled icon-tabler-charging-pile\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M12 3a3 3 0 0 1 3 3v4a3 3 0 0 1 3 3v3a.5 .5 0 1 0 1 0v-6.585l-1 -1l-.293 .292a1 1 0 0 1 -1.414 -1.414l.292 -.293l-.292 -.293a1 1 0 0 1 -.083 -1.32l.083 -.094a1 1 0 0 1 1.414 0l3 3a1 1 0 0 1 .293 .707v7a2.5 2.5 0 1 1 -5 0v-3a1 1 0 0 0 -1 -1v7a1 1 0 0 1 0 2h-12a1 1 0 0 1 0 -2v-13a3 3 0 0 1 3 -3zm-2.486 7.643a1 1 0 0 0 -1.371 .343l-1.5 2.5l-.054 .1a1 1 0 0 0 .911 1.414h1.233l-.59 .986a1 1 0 0 0 1.714 1.028l1.5 -2.5l.054 -.1a1 1 0 0 0 -.911 -1.414h-1.235l.592 -.986a1 1 0 0 0 -.343 -1.371m2.486 -5.643h-6a1 1 0 0 0 -1 1v1h8v-1a1 1 0 0 0 -1 -1\" /></svg>",
                        "icon_off": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_on": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"fan-spin windmill-spin me-2 text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "layout": 2,
                        "location_name": "Location Solar",
                        "measurement": "ev1",
                        "menu": 1,
                        "mqtt_connected": true,
                        "mqtt_control_off": "0",
                        "mqtt_control_on": "1",
                        "mqtt_data_control": "SORACELL1/CONTROL",
                        "mqtt_data_value": "SORACELL1/DATA",
                        "status": 1,
                        "timestamp": "2026-08-27 21:58:07",
                        "type_id": 9,
                        "type_name": "ZONE 2",
                        "unit": "%",
                        "value_data": "70.00"
                    },
                    {
                        "alarm_status": 5,
                        "alarm_subject": "Normal",
                        "alarm_title": "Normal",
                        "cache_used": true,
                        "control": [],
                        "device_id": 88,
                        "device_name": "Amp",
                        "devicedata": "75.00 A",
                        "graph": "http://localhost:5000/api/iot/monitordevicechart?bucket=SORACELL1&measurement=ev1amp&field=value&start=-5m&stop=now()&limit=120&lang=en",
                        "hardware_id": 1,
                        "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"icon icon-tabler icons-tabler-outline icon-tabler-circuit-ammeter\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M5 12a7 7 0 1 0 14 0a7 7 0 1 0 -14 0\" /><path d=\"M5 12h-3\" /><path d=\"M19 12h3\" /><path d=\"M10 14v-3c0 -1.036 .895 -2 2 -2s2 .964 2 2v3\" /><path d=\"M14 12h-4\" /></svg>",
                        "icon_access": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"icon icon-tabler icons-tabler-outline icon-tabler-circuit-ammeter\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M5 12a7 7 0 1 0 14 0a7 7 0 1 0 -14 0\" /><path d=\"M5 12h-3\" /><path d=\"M19 12h3\" /><path d=\"M10 14v-3c0 -1.036 .895 -2 2 -2s2 .964 2 2v3\" /><path d=\"M14 12h-4\" /></svg>",
                        "icon_off": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_on": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"fan-spin windmill-spin me-2 text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "layout": 2,
                        "location_name": "Location Solar",
                        "measurement": "ev1amp",
                        "menu": 1,
                        "mqtt_connected": true,
                        "mqtt_control_off": "0",
                        "mqtt_control_on": "1",
                        "mqtt_data_control": "SORACELL1/CONTROL",
                        "mqtt_data_value": "SORACELL1/DATA",
                        "status": 1,
                        "timestamp": "2026-08-27 21:58:07",
                        "type_id": 9,
                        "type_name": "ZONE 2",
                        "unit": "A",
                        "value_data": "75.00"
                    },
                    {
                        "alarm_status": 5,
                        "alarm_subject": "Normal",
                        "alarm_title": "Normal",
                        "cache_used": true,
                        "control": [],
                        "device_id": 89,
                        "device_name": "Volt",
                        "devicedata": "220.00 V",
                        "graph": "http://localhost:5000/api/iot/monitordevicechart?bucket=SORACELL1&measurement=ev1volt&field=value&start=-5m&stop=now()&limit=120&lang=en",
                        "hardware_id": 1,
                        "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"icon icon-tabler icons-tabler-outline icon-tabler-circuit-voltmeter\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M5 12a7 7 0 1 0 14 0a7 7 0 1 0 -14 0\" /><path d=\"M5 12h-3\" /><path d=\"M19 12h3\" /><path d=\"M10 10l2 4l2 -4\" /></svg>",
                        "icon_access": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"icon icon-tabler icons-tabler-outline icon-tabler-circuit-voltmeter\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M5 12a7 7 0 1 0 14 0a7 7 0 1 0 -14 0\" /><path d=\"M5 12h-3\" /><path d=\"M19 12h3\" /><path d=\"M10 10l2 4l2 -4\" /></svg>",
                        "icon_off": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_on": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"fan-spin windmill-spin me-2 text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "layout": 2,
                        "location_name": "Location Solar",
                        "measurement": "ev1volt",
                        "menu": 1,
                        "mqtt_connected": true,
                        "mqtt_control_off": "0",
                        "mqtt_control_on": "1",
                        "mqtt_data_control": "SORACELL1/CONTROL",
                        "mqtt_data_value": "SORACELL1/DATA",
                        "status": 1,
                        "timestamp": "2026-08-27 21:58:07",
                        "type_id": 9,
                        "type_name": "ZONE 2",
                        "unit": "V",
                        "value_data": "220.00"
                    },
                    {
                        "alarm_status": 5,
                        "alarm_subject": "Normal",
                        "alarm_title": "Normal",
                        "cache_used": true,
                        "control": [],
                        "device_id": 90,
                        "device_name": "Watt",
                        "devicedata": "1.50 Kw",
                        "graph": "http://localhost:5000/api/iot/monitordevicechart?bucket=SORACELL1&measurement=ev1watt&field=value&start=-5m&stop=now()&limit=120&lang=en",
                        "hardware_id": 1,
                        "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"icon icon-tabler icons-tabler-outline icon-tabler-circuit-resistor\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M2 12h2l2 -5l3 10l3 -10l3 10l3 -10l1.5 5h2.5\" /></svg>",
                        "icon_access": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"icon icon-tabler icons-tabler-outline icon-tabler-circuit-resistor\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M2 12h2l2 -5l3 10l3 -10l3 10l3 -10l1.5 5h2.5\" /></svg>",
                        "icon_off": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_on": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"fan-spin windmill-spin me-2 text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "layout": 2,
                        "location_name": "Location Solar",
                        "measurement": "ev1watt",
                        "menu": 1,
                        "mqtt_connected": true,
                        "mqtt_control_off": "0",
                        "mqtt_control_on": "1",
                        "mqtt_data_control": "SORACELL1/CONTROL",
                        "mqtt_data_value": "SORACELL1/DATA",
                        "status": 1,
                        "timestamp": "2026-08-27 21:58:07",
                        "type_id": 9,
                        "type_name": "ZONE 2",
                        "unit": "Kw",
                        "value_data": "1.50"
                    },
                    {
                        "alarm_status": 5,
                        "alarm_subject": "Normal",
                        "alarm_title": "Normal",
                        "cache_used": true,
                        "control": [],
                        "device_id": 91,
                        "device_name": "temperature",
                        "devicedata": "25.50 ",
                        "graph": "http://localhost:5000/api/iot/monitordevicechart?bucket=SORACELL1&measurement=temperature&field=value&start=-5m&stop=now()&limit=120&lang=en",
                        "hardware_id": 1,
                        "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"icon icon-tabler icons-tabler-outline icon-tabler-temperature\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M10 13.5a4 4 0 1 0 4 0v-8.5a2 2 0 0 0 -4 0v8.5\" /><path d=\"M10 9l4 0\" /></svg>",
                        "icon_access": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"icon icon-tabler icons-tabler-outline icon-tabler-temperature\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M10 13.5a4 4 0 1 0 4 0v-8.5a2 2 0 0 0 -4 0v8.5\" /><path d=\"M10 9l4 0\" /></svg>",
                        "icon_off": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_on": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"fan-spin windmill-spin me-2 text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "layout": 2,
                        "location_name": "Location Solar",
                        "measurement": "temperature",
                        "menu": 1,
                        "mqtt_connected": true,
                        "mqtt_control_off": "0",
                        "mqtt_control_on": "1",
                        "mqtt_data_control": "SORACELL1/CONTROL",
                        "mqtt_data_value": "SORACELL1/DATA",
                        "status": 1,
                        "timestamp": "2026-08-27 21:58:07",
                        "type_id": 10,
                        "type_name": "ZONE 3",
                        "unit": "",
                        "value_data": "25.50"
                    },
                    {
                        "alarm_status": 5,
                        "alarm_subject": "Normal",
                        "alarm_title": "Normal",
                        "cache_used": true,
                        "control": [],
                        "device_id": 92,
                        "device_name": "humidity",
                        "devicedata": "65.50 %",
                        "graph": "http://localhost:5000/api/iot/monitordevicechart?bucket=SORACELL1&measurement=humidity&field=value&start=-5m&stop=now()&limit=120&lang=en",
                        "hardware_id": 1,
                        "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"icon icon-tabler icons-tabler-outline icon-tabler-droplet\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M7.502 19.423c2.602 2.105 6.395 2.105 8.996 0c2.602 -2.105 3.262 -5.708 1.566 -8.546l-4.89 -7.26c-.42 -.625 -1.287 -.803 -1.936 -.397a1.376 1.376 0 0 0 -.41 .397l-4.893 7.26c-1.695 2.838 -1.035 6.441 1.567 8.546\" /></svg>",
                        "icon_access": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"icon icon-tabler icons-tabler-outline icon-tabler-droplet\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M7.502 19.423c2.602 2.105 6.395 2.105 8.996 0c2.602 -2.105 3.262 -5.708 1.566 -8.546l-4.89 -7.26c-.42 -.625 -1.287 -.803 -1.936 -.397a1.376 1.376 0 0 0 -.41 .397l-4.893 7.26c-1.695 2.838 -1.035 6.441 1.567 8.546\" /></svg>",
                        "icon_off": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_on": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"fan-spin windmill-spin me-2 text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "layout": 2,
                        "location_name": "Location Solar",
                        "measurement": "humidity",
                        "menu": 1,
                        "mqtt_connected": true,
                        "mqtt_control_off": "0",
                        "mqtt_control_on": "1",
                        "mqtt_data_control": "SORACELL1/CONTROL",
                        "mqtt_data_value": "SORACELL1/DATA",
                        "status": 1,
                        "timestamp": "2026-08-27 21:58:07",
                        "type_id": 10,
                        "type_name": "ZONE 3",
                        "unit": "%",
                        "value_data": "65.50"
                    }
                ],
                "group_id": 1,
                "group_name": "Sensor"
            },
            {
                "count": 4,
                "devices": [
                    {
                        "alarm_status": 5,
                        "alarm_subject": "Normal",
                        "alarm_title": "Normal",
                        "cache_used": true,
                        "control": "http://localhost:5000/api/iot/controls?topic=SORACELL1/CONTROL&message=0",
                        "device_id": 93,
                        "device_name": "fan1",
                        "devicedata": "OFF",
                        "graph": "http://localhost:5000/api/iot/monitordevicechart?bucket=SORACELL1&measurement=fan1&field=value&start=-5m&stop=now()&limit=120&lang=en",
                        "hardware_id": 2,
                        "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"fan-spin windmill-spin me-2 text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_access": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_off": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_on": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"fan-spin windmill-spin me-2 text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "layout": 2,
                        "location_name": "Location Solar",
                        "measurement": "fan1",
                        "menu": 1,
                        "mqtt_connected": true,
                        "mqtt_control_off": "0",
                        "mqtt_control_on": "1",
                        "mqtt_data_control": "SORACELL1/CONTROL",
                        "mqtt_data_value": "SORACELL1/DATA",
                        "status": 1,
                        "timestamp": "2026-08-27 21:58:07",
                        "type_id": 11,
                        "type_name": "ZONE 4",
                        "unit": "",
                        "value_data": "1"
                    },
                    {
                        "alarm_status": 5,
                        "alarm_subject": "Normal",
                        "alarm_title": "Normal",
                        "cache_used": true,
                        "control": "http://localhost:5000/api/iot/controls?topic=SORACELL1/CONTROL&message=0",
                        "device_id": 94,
                        "device_name": "fan2",
                        "devicedata": "OFF",
                        "graph": "http://localhost:5000/api/iot/monitordevicechart?bucket=SORACELL1&measurement=fan2&field=value&start=-5m&stop=now()&limit=120&lang=en",
                        "hardware_id": 2,
                        "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"fan-spin windmill-spin me-2 text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_access": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_off": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_on": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"fan-spin windmill-spin me-2 text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "layout": 2,
                        "location_name": "Location Solar",
                        "measurement": "fan2",
                        "menu": 1,
                        "mqtt_connected": true,
                        "mqtt_control_off": "0",
                        "mqtt_control_on": "1",
                        "mqtt_data_control": "SORACELL1/CONTROL",
                        "mqtt_data_value": "SORACELL1/DATA",
                        "status": 1,
                        "timestamp": "2026-08-27 21:58:07",
                        "type_id": 12,
                        "type_name": "ZONE 5",
                        "unit": "",
                        "value_data": "1"
                    },
                    {
                        "alarm_status": 5,
                        "alarm_subject": "Normal",
                        "alarm_title": "Normal",
                        "cache_used": true,
                        "control": "http://localhost:5000/api/iot/controls?topic=SORACELL1/CONTROL&message=0",
                        "device_id": 95,
                        "device_name": "fan1",
                        "devicedata": "OFF",
                        "graph": "http://localhost:5000/api/iot/monitordevicechart?bucket=SORACELL1&measurement=fan1&field=value&start=-5m&stop=now()&limit=120&lang=en",
                        "hardware_id": 2,
                        "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"fan-spin windmill-spin me-2 text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_access": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_off": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_on": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"fan-spin windmill-spin me-2 text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "layout": 1,
                        "location_name": "Location Solar",
                        "measurement": "fan1",
                        "menu": 1,
                        "mqtt_connected": true,
                        "mqtt_control_off": "0",
                        "mqtt_control_on": "1",
                        "mqtt_data_control": "SORACELL1/CONTROL",
                        "mqtt_data_value": "SORACELL1/DATA",
                        "status": 1,
                        "timestamp": "2026-08-27 21:58:07",
                        "type_id": 2,
                        "type_name": "ZONE Fan1",
                        "unit": "",
                        "value_data": "1"
                    },
                    {
                        "alarm_status": 5,
                        "alarm_subject": "Normal",
                        "alarm_title": "Normal",
                        "cache_used": true,
                        "control": "http://localhost:5000/api/iot/controls?topic=SORACELL1/CONTROL&message=0",
                        "device_id": 96,
                        "device_name": "fan2",
                        "devicedata": "OFF",
                        "graph": "http://localhost:5000/api/iot/monitordevicechart?bucket=SORACELL1&measurement=fan2&field=value&start=-5m&stop=now()&limit=120&lang=en",
                        "hardware_id": 2,
                        "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"fan-spin windmill-spin me-2 text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_access": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_off": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "icon_on": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" \n            fill=\"currentColor\" \n            stroke-width=\"2\" stroke-linecap=\"round\" \n            stroke-linejoin=\"round\" \n            class=\"fan-spin windmill-spin me-2 text-warning\"\n            style=\"vertical-align: middle;\">\n            <path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"></path>\n            <path d=\"M12.363 2.068l4.912 1.914a2.7 2.7 0 0 1 .68 4.646l-3.045 2.371l6.09 .001a1 1 0 0 1 .932 1.363l-1.914 4.912a2.7 2.7 0 0 1 -4.646 .68l-2.372 -3.047v6.092a1 1 0 0 1 -1.363 .932l-4.912 -1.914a2.7 2.7 0 0 1 -.68 -4.646l3.045 -2.372h-6.09a1 1 0 0 1 -.932 -1.363l1.914 -4.912a2.7 2.7 0 0 1 4.646 -.68l2.371 3.044l.001 -6.089a1 1 0 0 1 1.363 -.932\"></path>\n        </svg>",
                        "layout": 1,
                        "location_name": "Location Solar",
                        "measurement": "fan2",
                        "menu": 1,
                        "mqtt_connected": true,
                        "mqtt_control_off": "0",
                        "mqtt_control_on": "1",
                        "mqtt_data_control": "SORACELL1/CONTROL",
                        "mqtt_data_value": "SORACELL1/DATA",
                        "status": 1,
                        "timestamp": "2026-08-27 21:58:07",
                        "type_id": 2,
                        "type_name": "ZONE Fan1",
                        "unit": "",
                        "value_data": "1"
                    }
                ],
                "group_id": 2,
                "group_name": "IO Sensor"
            }
        ],
        "device_count": 13,
        "device_type": "",
        "group_name": "",
        "layout": 2,
        "layout_name": "Card",
        "mqtt_connected": true,
        "mqtt_raw_payload": "250,300,1.5,70,75,220,1.5,25.5,65.5,1,1",
        "timestamp": "2026-08-27 21:58:07"
    },
    "is_success": true
}

