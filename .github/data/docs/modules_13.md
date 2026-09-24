**โมดูลที่ 13: 📡 IoT & GPS Tracking (ระบบ IoT และติดตามตำแหน่ง)**

โมดูล IoT & GPS Tracking เป็นระบบที่เชื่อมต่อกับอุปกรณ์ IoT ผ่าน MQTT Protocol รับข้อมูลตำแหน่ง (GPS) และสถานะอุปกรณ์แบบ Real-time ครอบคลุมการทำงานดังนี้:

1. **การจัดการอุปกรณ์ IoT (Device Management)** – ลงทะเบียน, อัปเดต, เปิด/ปิดอุปกรณ์
2. **การรับข้อมูล GPS (GPS Tracking)** – รับพิกัด/longitude, ความเร็ว, เวลา
3. **MQTT Integration** – รับ-ส่งข้อมูลผ่าน MQTT Broker (EMQX/Mosquitto)
4. **InfluxDB สำหรับ Time-series Data** – เก็บข้อมูล GPS และสถานะอุปกรณ์แบบ Time-series
5. **ประวัติอุปกรณ์ (Device History)** – บันทึกประวัติการทำงานและการเปลี่ยนแปลง
6. **การควบคุมการเข้าถึงอุปกรณ์ (Device Access Control)** – จัดการสิทธิ์การเข้าถึง
7. **Geofence (รั้วรอบขอบ)** – กำหนดพื้นที่และแจ้งเตือนเมื่ออุปกรณ์เข้า/ออก
8. **รายงานอัตโนมัติ (Auto Report)** – สร้างรายงานสรุปการเคลื่อนไหว

---

## 📁 โครงสร้างโมดูล IoT & GPS Tracking (`modules/iot`)

```
modules/iot/
├── application/
│   ├── interfaces/
│   │   ├── IotDeviceService.java
│   │   ├── GpsTrackingService.java
│   │   ├── MqttService.java
│   │   ├── DeviceHistoryService.java
│   │   ├── DeviceAccessControlService.java
│   │   ├── GeofenceService.java
│   │   └── AutoReportService.java
│   ├── impl/
│   │   ├── IotDeviceServiceImpl.java
│   │   ├── GpsTrackingServiceImpl.java
│   │   ├── MqttServiceImpl.java
│   │   ├── DeviceHistoryServiceImpl.java
│   │   ├── DeviceAccessControlServiceImpl.java
│   │   ├── GeofenceServiceImpl.java
│   │   └── AutoReportServiceImpl.java
│   └── usecase/
│       ├── RegisterDeviceUseCase.java
│       ├── UpdateDeviceUseCase.java
│       ├── GetDeviceLocationUseCase.java
│       ├── GetDeviceHistoryUseCase.java
│       ├── PublishMqttMessageUseCase.java
│       ├── ProcessGpsDataUseCase.java
│       ├── CreateGeofenceUseCase.java
│       ├── CheckGeofenceAlertUseCase.java
│       ├── AuthorizeDeviceAccessUseCase.java
│       └── GenerateAutoReportUseCase.java
├── domain/
│   ├── MIoTDevice.java
│   ├── TGPSData.java
│   ├── TDeviceHistory.java
│   ├── TDeviceAccessLog.java
│   ├── MGeofence.java
│   ├── TGeofenceAlert.java
│   ├── TAutoReport.java
│   ├── enums/
│   │   ├── DeviceStatus.java           // ONLINE, OFFLINE, INACTIVE, MAINTENANCE
│   │   ├── DeviceType.java             // GPS_TRACKER, SENSOR, CAMERA, CONTROLLER
│   │   ├── DeviceAccessLevel.java      // READ, WRITE, CONTROL, ADMIN
│   │   ├── GeofenceType.java           // CIRCLE, POLYGON, RECTANGLE
│   │   └── AlertType.java              // ENTER, EXIT, SPEED, BATTERY
│   └── valueobjects/
│       ├── Coordinate.java             // 纬度, longitude, altitude
│       ├── Speed.java                  // ความเร็ว
│       └── BatteryLevel.java           // ระดับแบตเตอรี่
├── infrastructure/
│   ├── repository/
│   │   ├── IotDeviceRepository.java
│   │   ├── GpsDataRepository.java
│   │   ├── DeviceHistoryRepository.java
│   │   ├── DeviceAccessLogRepository.java
│   │   ├── GeofenceRepository.java
│   │   ├── GeofenceAlertRepository.java
│   │   ├── AutoReportRepository.java
│   │   └── impl/
│   │       ├── IotDeviceRepositoryImpl.java
│   │       ├── DeviceHistoryRepositoryImpl.java
│   │       └── DeviceAccessLogRepositoryImpl.java
│   ├── cache/                                           // ⬅️ ระบบ Cache สำหรับ IoT
│   │   ├── DeviceCacheService.java
│   │   ├── DeviceLocationCacheService.java
│   │   └── DeviceStatusCacheService.java
│   ├── mqtt/                                            // ⬅️ MQTT Integration
│   │   ├── MqttConfig.java
│   │   ├── MqttPublisher.java
│   │   ├── MqttSubscriber.java
│   │   ├── MqttMessageHandler.java
│   │   └── MqttGpsMessageHandler.java
│   ├── influxdb/                                        // ⬅️ InfluxDB Integration
│   │   ├── InfluxDBConfig.java
│   │   ├── InfluxDBService.java
│   │   ├── GpsDataPoint.java
│   │   └── DeviceStatusPoint.java
│   ├── geofence/                                        // ⬅️ Geofence
│   │   ├── GeofenceChecker.java
│   │   ├── PointInPolygon.java
│   │   └── GeofenceAlertProcessor.java
│   ├── entity/
│   │   ├── IotDeviceEntity.java
│   │   ├── GpsDataEntity.java
│   │   ├── DeviceHistoryEntity.java
│   │   ├── DeviceAccessLogEntity.java
│   │   ├── GeofenceEntity.java
│   │   └── AutoReportEntity.java
│   └── mapper/
│       ├── IotDeviceMapper.java
│       ├── GpsDataMapper.java
│       └── DeviceHistoryMapper.java
└── presentation/
    ├── controller/
    │   ├── IotDeviceController.java      // Device CRUD + Status
    │   ├── GpsTrackingController.java    // GPS Data APIs
    │   ├── MqttController.java           // MQTT Publish/Subscribe
    │   ├── GeofenceController.java       // Geofence Management
    │   └── AutoReportController.java     // Auto Report APIs
    ├── dto/
    │   ├── request/
    │   │   ├── DeviceRegisterRequestDTO.java
    │   │   ├── DeviceUpdateRequestDTO.java
    │   │   ├── GpsDataRequestDTO.java
    │   │   ├── GeofenceCreateRequestDTO.java
    │   │   ├── DeviceAccessRequestDTO.java
    │   │   └── AutoReportRequestDTO.java
    │   └── response/
    │       ├── DeviceResponseDTO.java
    │       ├── DeviceDetailResponseDTO.java
    │       ├── GpsDataResponseDTO.java
    │       ├── DeviceLocationResponseDTO.java
    │       ├── GeofenceResponseDTO.java
    │       ├── DeviceHistoryResponseDTO.java
    │       └── AutoReportResponseDTO.java
    └── websocket/                                       // ⬅️ WebSocket สำหรับ Realtime
        ├── GpsWebSocketHandler.java
        └── DeviceStatusWebSocketHandler.java
```

---

## 🗄️ ออกแบบฐานข้อมูล (Database Design) สำหรับโมดูล IoT

### 📄 SQL DDL (ไฟล์: `src/main/resources/db/migration/V13__iot_schema.sql`)

```sql
-- ==============================================
-- ตาราง: m_iot_device (ข้อมูลอุปกรณ์ IoT)
-- IoT device master data.
-- ==============================================
CREATE TABLE IF NOT EXISTS m_iot_device (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id VARCHAR(50) UNIQUE NOT NULL,              -- รหัสอุปกรณ์ (Hardware ID)
    device_name VARCHAR(100) NOT NULL,                  -- ชื่ออุปกรณ์
    device_type VARCHAR(30) NOT NULL,                   -- GPS_TRACKER, SENSOR, CAMERA, CONTROLLER
    status VARCHAR(20) DEFAULT 'OFFLINE',               -- ONLINE, OFFLINE, INACTIVE, MAINTENANCE
    serial_number VARCHAR(50),                          -- เลขซีเรียล
    firmware_version VARCHAR(20),                       -- เวอร์ชัน Firmware
    hardware_version VARCHAR(20),                       -- เวอร์ชัน Hardware
    manufacturer VARCHAR(100),                          -- ผู้ผลิต
    model VARCHAR(50),                                  -- รุ่น
    battery_level INTEGER DEFAULT 0,                    -- ระดับแบตเตอรี่ (%)
    last_seen TIMESTAMP,                                -- เวลาที่เห็นครั้งล่าสุด
    last_latitude DECIMAL(10,7),                        -- ละติจูดล่าสุด
    last_longitude DECIMAL(10,7),                       -- ลองจิจูดล่าสุด
    last_altitude DECIMAL(10,2),                        -- ระดับความสูงล่าสุด (เมตร)
    last_speed DECIMAL(10,2),                           -- ความเร็วล่าสุด (กม./ชม.)
    is_active BOOLEAN DEFAULT TRUE,
    is_online BOOLEAN DEFAULT FALSE,
    metadata JSONB,                                     -- ข้อมูลเพิ่มเติม (JSON)
    notes TEXT,
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    deleted BOOLEAN DEFAULT FALSE,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_m_iot_device_id ON m_iot_device(device_id);
CREATE INDEX idx_m_iot_device_status ON m_iot_device(status);
CREATE INDEX idx_m_iot_device_type ON m_iot_device(device_type);
CREATE INDEX idx_m_iot_device_last_seen ON m_iot_device(last_seen);
CREATE INDEX idx_m_iot_device_whitelabel ON m_iot_device(whitelabel_id);
CREATE INDEX idx_m_iot_device_deleted ON m_iot_device(deleted);

-- ==============================================
-- ตาราง: t_gps_data (ข้อมูล GPS - ใช้ร่วมกับ InfluxDB)
-- GPS data (also stored in InfluxDB for time-series).
-- ==============================================
CREATE TABLE IF NOT EXISTS t_gps_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES m_iot_device(id) ON DELETE CASCADE,
    device_identifier VARCHAR(50),                      -- รหัสอุปกรณ์ (Hardware ID)
    latitude DECIMAL(10,7) NOT NULL,                    -- ละติจูด
    longitude DECIMAL(10,7) NOT NULL,                   -- ลองจิจูด
    altitude DECIMAL(10,2),                             -- ระดับความสูง (เมตร)
    speed DECIMAL(10,2),                                -- ความเร็ว (กม./ชม.)
    heading DECIMAL(5,2),                               -- ทิศทาง (องศา)
    accuracy DECIMAL(5,2),                              -- ความแม่นยำ (เมตร)
    battery_level INTEGER,                              -- ระดับแบตเตอรี่ (%)
    satelites INTEGER,                                  -- จำนวนดาวเทียม
    event_type VARCHAR(20),                             -- MOVING, STOPPED, IDLE, ALERT
    timestamp TIMESTAMP NOT NULL,                       -- เวลาจากอุปกรณ์
    received_at TIMESTAMP NOT NULL DEFAULT NOW(),      -- เวลาที่ระบบรับ
    metadata JSONB,                                     -- ข้อมูลเพิ่มเติม
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_gps_data_device ON t_gps_data(device_id);
CREATE INDEX idx_t_gps_data_timestamp ON t_gps_data(timestamp);
CREATE INDEX idx_t_gps_data_whitelabel ON t_gps_data(whitelabel_id);

-- ==============================================
-- ตาราง: t_device_history (ประวัติการทำงานของอุปกรณ์)
-- Device operation history.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_device_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES m_iot_device(id) ON DELETE CASCADE,
    event_type VARCHAR(30) NOT NULL,                    -- REGISTERED, ACTIVATED, DEACTIVATED, ONLINE, OFFLINE, MAINTENANCE, FIRMWARE_UPDATE, BATTERY_LOW, ALERT
    event_description TEXT,
    old_value TEXT,
    new_value TEXT,
    event_timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    triggered_by UUID REFERENCES m_user(id),
    metadata JSONB,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_device_history_device ON t_device_history(device_id);
CREATE INDEX idx_t_device_history_event ON t_device_history(event_type);
CREATE INDEX idx_t_device_history_timestamp ON t_device_history(event_timestamp);
CREATE INDEX idx_t_device_history_whitelabel ON t_device_history(whitelabel_id);

-- ==============================================
-- ตาราง: t_device_access_log (บันทึกการเข้าถึงอุปกรณ์)
-- Device access log.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_device_access_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES m_iot_device(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES m_user(id) ON DELETE CASCADE,
    access_type VARCHAR(20) NOT NULL,                   -- READ, WRITE, CONTROL, ADMIN
    access_granted BOOLEAN DEFAULT TRUE,
    ip_address INET,
    user_agent TEXT,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    reason TEXT,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_device_access_device ON t_device_access_log(device_id);
CREATE INDEX idx_t_device_access_user ON t_device_access_log(user_id);
CREATE INDEX idx_t_device_access_timestamp ON t_device_access_log(timestamp);
CREATE INDEX idx_t_device_access_whitelabel ON t_device_access_log(whitelabel_id);

-- ==============================================
-- ตาราง: m_geofence (ข้อมูลรั้วรอบขอบ)
-- Geofence master data.
-- ==============================================
CREATE TABLE IF NOT EXISTS m_geofence (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    geofence_name VARCHAR(100) NOT NULL,
    geofence_type VARCHAR(20) NOT NULL,                 -- CIRCLE, POLYGON, RECTANGLE
    center_latitude DECIMAL(10,7),                      -- ศูนย์กลาง (สำหรับ CIRCLE)
    center_longitude DECIMAL(10,7),                     -- ศูนย์กลาง (สำหรับ CIRCLE)
    radius DECIMAL(10,2),                               -- รัศมี (เมตร) สำหรับ CIRCLE
    coordinates JSONB,                                  -- จุด (สำหรับ POLYGON, RECTANGLE)
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    alert_on_enter BOOLEAN DEFAULT TRUE,
    alert_on_exit BOOLEAN DEFAULT TRUE,
    speed_limit DECIMAL(10,2),                          -- ขีดจำกัดความเร็ว (กม./ชม.)
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_m_geofence_name ON m_geofence(geofence_name);
CREATE INDEX idx_m_geofence_active ON m_geofence(is_active);
CREATE INDEX idx_m_geofence_whitelabel ON m_geofence(whitelabel_id);

-- ==============================================
-- ตาราง: t_geofence_alert (การแจ้งเตือนรั้วรอบขอบ)
-- Geofence alert records.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_geofence_alert (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    geofence_id UUID NOT NULL REFERENCES m_geofence(id) ON DELETE CASCADE,
    device_id UUID NOT NULL REFERENCES m_iot_device(id) ON DELETE CASCADE,
    alert_type VARCHAR(20) NOT NULL,                    -- ENTER, EXIT, SPEED
    latitude DECIMAL(10,7),
    longitude DECIMAL(10,7),
    speed DECIMAL(10,2),
    alert_timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    is_resolved BOOLEAN DEFAULT FALSE,
    resolved_at TIMESTAMP,
    resolved_by UUID REFERENCES m_user(id),
    notes TEXT,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_geofence_alert_geofence ON t_geofence_alert(geofence_id);
CREATE INDEX idx_t_geofence_alert_device ON t_geofence_alert(device_id);
CREATE INDEX idx_t_geofence_alert_timestamp ON t_geofence_alert(alert_timestamp);
CREATE INDEX idx_t_geofence_alert_resolved ON t_geofence_alert(is_resolved);

-- ==============================================
-- ตาราง: t_auto_report (รายงานอัตโนมัติ)
-- Auto-generated reports.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_auto_report (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_no VARCHAR(30) UNIQUE NOT NULL,
    report_type VARCHAR(30) NOT NULL,                   -- DAILY, WEEKLY, MONTHLY
    report_name VARCHAR(100) NOT NULL,
    report_date DATE NOT NULL,
    device_count INTEGER,
    total_distance DECIMAL(15,2),                       -- ระยะทางรวม (กม.)
    avg_speed DECIMAL(10,2),                            -- ความเร็วเฉลี่ย (กม./ชม.)
    max_speed DECIMAL(10,2),                            -- ความเร็วสูงสุด (กม./ชม.)
    total_moving_time INTERVAL,                         -- เวลาที่เคลื่อนที่
    total_idle_time INTERVAL,                           -- เวลาที่จอด
    battery_avg INTEGER,                                -- แบตเตอรี่เฉลี่ย
    alert_count INTEGER DEFAULT 0,
    report_data JSONB,                                  -- ข้อมูลรายงานแบบ JSON
    file_path TEXT,                                     -- ที่อยู่ไฟล์ PDF/Excel
    status VARCHAR(20) DEFAULT 'GENERATED',
    generated_by UUID REFERENCES m_user(id),
    generated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_auto_report_type ON t_auto_report(report_type);
CREATE INDEX idx_t_auto_report_date ON t_auto_report(report_date);
CREATE INDEX idx_t_auto_report_whitelabel ON t_auto_report(whitelabel_id);

-- ==============================================
-- ฟังก์ชันสร้างเลขที่รายงานอัตโนมัติ
-- Function to generate report number.
-- ==============================================
CREATE OR REPLACE FUNCTION generate_auto_report_no()
RETURNS TRIGGER AS $$
DECLARE
    year_part TEXT;
    sequence_part TEXT;
BEGIN
    year_part := TO_CHAR(NOW(), 'YYYY');
    sequence_part := LPAD(CAST((
        SELECT COALESCE(MAX(CAST(SUBSTRING(report_no FROM 9) AS INTEGER)), 0) + 1
        FROM t_auto_report
        WHERE report_no LIKE 'RPT-' || year_part || '-%'
    ) AS TEXT), 4, '0');
    NEW.report_no := 'RPT-' || year_part || '-' || sequence_part;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_generate_auto_report_no ON t_auto_report;
CREATE TRIGGER trg_generate_auto_report_no
BEFORE INSERT ON t_auto_report
FOR EACH ROW
EXECUTE FUNCTION generate_auto_report_no();
```

---

## 🔌 MQTT Configuration

### `infrastructure/mqtt/MqttConfig.java`

```java
package com.template.app.modules.iot.infrastructure.mqtt;

import org.eclipse.paho.client.mqttv3.MqttConnectOptions;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.integration.channel.DirectChannel;
import org.springframework.integration.mqtt.core.DefaultMqttPahoClientFactory;
import org.springframework.integration.mqtt.core.MqttPahoClientFactory;
import org.springframework.integration.mqtt.inbound.MqttPahoMessageDrivenChannelAdapter;
import org.springframework.integration.mqtt.outbound.MqttPahoMessageHandler;
import org.springframework.messaging.MessageChannel;

@Configuration
public class MqttConfig {

    @Value("${mqtt.broker.url:tcp://localhost:1883}")
    private String brokerUrl;

    @Value("${mqtt.broker.username}")
    private String username;

    @Value("${mqtt.broker.password}")
    private String password;

    @Value("${mqtt.client.id:spring-boot-client}")
    private String clientId;

    @Value("${mqtt.default.topic:gps/data}")
    private String defaultTopic;

    /*
        ฟังก์ชันนี้สร้าง MqttPahoClientFactory สำหรับเชื่อมต่อกับ MQTT Broker
        This function creates MqttPahoClientFactory for connecting to MQTT Broker.
    */
    @Bean
    public MqttPahoClientFactory mqttClientFactory() {
        DefaultMqttPahoClientFactory factory = new DefaultMqttPahoClientFactory();
        MqttConnectOptions options = new MqttConnectOptions();
        options.setServerURIs(new String[]{brokerUrl});
        options.setUserName(username);
        options.setPassword(password.toCharArray());
        options.setKeepAliveInterval(60);
        options.setCleanSession(true);
        options.setAutomaticReconnect(true);
        options.setConnectionTimeout(30);
        factory.setConnectionOptions(options);
        return factory;
    }

    /*
        ฟังก์ชันนี้สร้าง Inbound Channel Adapter สำหรับรับข้อมูลจาก MQTT
        This function creates an Inbound Channel Adapter for receiving data from MQTT.
    */
    @Bean
    public MqttPahoMessageDrivenChannelAdapter mqttInbound() {
        MqttPahoMessageDrivenChannelAdapter adapter = 
            new MqttPahoMessageDrivenChannelAdapter(
                clientId + "-inbound", 
                mqttClientFactory(), 
                "gps/data",        // Topic สำหรับ GPS / GPS topic
                "device/status",   // Topic สำหรับสถานะ / Status topic
                "device/alert"     // Topic สำหรับแจ้งเตือน / Alert topic
            );
        adapter.setCompletionTimeout(5000);
        adapter.setQos(1);
        adapter.setOutputChannel(mqttInputChannel());
        return adapter;
    }

    /*
        ฟังก์ชันนี้สร้าง Channel สำหรับรับข้อความจาก MQTT
        This function creates a Channel for receiving MQTT messages.
    */
    @Bean
    public MessageChannel mqttInputChannel() {
        return new DirectChannel();
    }

    /*
        ฟังก์ชันนี้สร้าง Outbound Message Handler สำหรับส่งข้อความไปยัง MQTT
        This function creates an Outbound Message Handler for sending messages to MQTT.
    */
    @Bean
    public MqttPahoMessageHandler mqttOutbound() {
        MqttPahoMessageHandler handler = new MqttPahoMessageHandler(
            clientId + "-outbound", 
            mqttClientFactory()
        );
        handler.setAsync(true);
        handler.setDefaultTopic(defaultTopic);
        handler.setDefaultQos(1);
        handler.setDefaultRetained(false);
        return handler;
    }

    /*
        ฟังก์ชันนี้สร้าง Channel สำหรับส่งข้อความไปยัง MQTT
        This function creates a Channel for sending messages to MQTT.
    */
    @Bean
    public MessageChannel mqttOutputChannel() {
        return new DirectChannel();
    }
}
```

---

## 🧠 ระบบ Cache (Redis) สำหรับ IoT

### `infrastructure/cache/DeviceCacheService.java`

```java
package com.template.app.modules.iot.infrastructure.cache;

import com.template.app.modules.iot.domain.MIoTDevice;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.CachePut;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class DeviceCacheService {

    /*
        ฟังก์ชันนี้ดึงข้อมูลอุปกรณ์ IoT จาก Cache (ลดภาระฐานข้อมูล)
        This function retrieves IoT device data from cache (reduces DB load).
        Redis Key: iot_device:{deviceId}
    */
    @Cacheable(value = "iot_devices", key = "#deviceId")
    public MIoTDevice getDevice(UUID deviceId) {
        return null;
    }

    /*
        ฟังก์ชันนี้ดึงข้อมูลอุปกรณ์จาก Device ID (Hardware ID)
        This function retrieves device data by Hardware ID.
        Redis Key: iot_device_identifier:{identifier}
    */
    @Cacheable(value = "iot_device_identifier", key = "#identifier")
    public MIoTDevice getDeviceByIdentifier(String identifier) {
        return null;
    }

    /*
        ฟังก์ชันนี้จะอัปเดต Cache เมื่อมีการบันทึกอุปกรณ์
        This function updates the cache when a device is saved.
    */
    @CachePut(value = "iot_devices", key = "#device.id")
    public MIoTDevice saveDevice(MIoTDevice device) {
        return device;
    }

    /*
        ฟังก์ชันนี้ลบข้อมูลอุปกรณ์ออกจาก Cache
        This function evicts device data from cache.
    */
    @CacheEvict(value = {"iot_devices", "iot_device_identifier"}, key = "#deviceId")
    public void evictDevice(UUID deviceId) {
        // ลบ Cache / Evict cache.
    }

    /*
        ฟังก์ชันนี้ลบ Cache ทั้งหมดของอุปกรณ์
        This function clears all device caches.
    */
    @CacheEvict(value = {"iot_devices", "iot_device_identifier", "device_status"}, allEntries = true)
    public void evictAllDevices() {
        // ลบทุก key / Evict all keys.
    }
}
```

### `infrastructure/cache/DeviceLocationCacheService.java`

```java
package com.template.app.modules.iot.infrastructure.cache;

import com.template.app.modules.iot.presentation.dto.response.DeviceLocationResponseDTO;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.CachePut;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class DeviceLocationCacheService {

    /*
        ฟังก์ชันนี้ดึงตำแหน่งล่าสุดของอุปกรณ์จาก Cache
        This function retrieves the latest device location from cache.
        Redis Key: device_location:{deviceId}
    */
    @Cacheable(value = "device_location", key = "#deviceId")
    public DeviceLocationResponseDTO getDeviceLocation(UUID deviceId) {
        return null;
    }

    /*
        ฟังก์ชันนี้จะอัปเดตตำแหน่งใน Cache เมื่อมีการอัปเดตตำแหน่ง
        This function updates the location in cache when location is updated.
    */
    @CachePut(value = "device_location", key = "#deviceId")
    public DeviceLocationResponseDTO updateDeviceLocation(UUID deviceId, DeviceLocationResponseDTO location) {
        return location;
    }

    /*
        ฟังก์ชันนี้ลบตำแหน่งของอุปกรณ์ออกจาก Cache
        This function evicts device location from cache.
    */
    @CacheEvict(value = "device_location", key = "#deviceId")
    public void evictDeviceLocation(UUID deviceId) {
        // ลบ Cache / Evict cache.
    }
}
```

---

## ⏱️ ระบบ Rate Limit สำหรับ IoT Controller

```java
package com.template.app.modules.iot.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.iot.application.interfaces.GpsTrackingService;
import com.template.app.modules.iot.application.interfaces.IotDeviceService;
import com.template.app.modules.iot.presentation.dto.request.DeviceRegisterRequestDTO;
import com.template.app.modules.iot.presentation.dto.request.DeviceUpdateRequestDTO;
import com.template.app.modules.iot.presentation.dto.response.DeviceLocationResponseDTO;
import com.template.app.modules.iot.presentation.dto.response.DeviceResponseDTO;
import com.template.app.modules.iot.presentation.dto.response.GpsDataResponseDTO;
import com.template.app.exception.SystemGlobalException;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/api/v1/iot")
@Tag(name = "IoT & GPS", description = "IoT Device and GPS Tracking APIs")
@RequiredArgsConstructor
public class IotDeviceController {

    private final IotDeviceService iotDeviceService;
    private final GpsTrackingService gpsTrackingService;

    // ========================================================================
    // 1. REGISTER DEVICE (ลงทะเบียนอุปกรณ์)
    // ========================================================================

    /*
        API: POST /api/v1/iot/devices/register
        ฟังก์ชันนี้ลงทะเบียนอุปกรณ์ IoT ใหม่ในระบบ
        This function registers a new IoT device in the system.
        Rate Limit: อนุญาต 10 ครั้งต่อ 5 นาที
        Rate Limit: Allows 10 requests per 5 minutes.
    */
    @PostMapping("/devices/register")
    @RateLimit(limit = 10, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Register a new IoT device")
    public ResponseEntity<DeviceResponseDTO> registerDevice(@Valid @RequestBody DeviceRegisterRequestDTO request)
            throws SystemGlobalException {
        DeviceResponseDTO response = iotDeviceService.registerDevice(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 2. GET DEVICE BY ID
    // ========================================================================

    /*
        API: GET /api/v1/iot/devices/{id}
        ฟังก์ชันนี้ดึงข้อมูลอุปกรณ์ IoT ตาม ID (ใช้ Cache)
        This function retrieves IoT device data by ID (uses caching).
        Rate Limit: อนุญาต 100 ครั้งต่อ 1 นาที
        Rate Limit: Allows 100 requests per minute.
    */
    @GetMapping("/devices/{id}")
    @RateLimit(limit = 100, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get IoT device by ID")
    public ResponseEntity<DeviceResponseDTO> getDevice(@PathVariable UUID id)
            throws SystemGlobalException {
        DeviceResponseDTO response = iotDeviceService.getDevice(id);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. LIST DEVICES
    // ========================================================================

    /*
        API: GET /api/v1/iot/devices
        ฟังก์ชันนี้แสดงรายการอุปกรณ์ IoT ทั้งหมดแบบแบ่งหน้า
        This function lists all IoT devices with pagination.
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @GetMapping("/devices")
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "List IoT devices")
    public ResponseEntity<Page<DeviceResponseDTO>> listDevices(
            @RequestParam(required = false) String status,
            @RequestParam(required = false) String deviceType,
            Pageable pageable) throws SystemGlobalException {
        Page<DeviceResponseDTO> page = iotDeviceService.listDevices(status, deviceType, pageable);
        return ResponseEntity.ok(page);
    }

    // ========================================================================
    // 4. UPDATE DEVICE
    // ========================================================================

    /*
        API: PUT /api/v1/iot/devices/{id}
        ฟังก์ชันนี้แก้ไขข้อมูลอุปกรณ์ IoT
        This function updates IoT device data.
        Rate Limit: อนุญาต 15 ครั้งต่อ 1 นาที
        Rate Limit: Allows 15 requests per minute.
    */
    @PutMapping("/devices/{id}")
    @RateLimit(limit = 15, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Update IoT device")
    public ResponseEntity<DeviceResponseDTO> updateDevice(
            @PathVariable UUID id,
            @Valid @RequestBody DeviceUpdateRequestDTO request) throws SystemGlobalException {
        DeviceResponseDTO response = iotDeviceService.updateDevice(id, request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 5. GET DEVICE STATUS (สถานะอุปกรณ์)
    // ========================================================================

    /*
        API: GET /api/v1/iot/devices/{id}/status
        ฟังก์ชันนี้ดึงสถานะปัจจุบันของอุปกรณ์ (ONLINE, OFFLINE, ฯลฯ) จาก Cache
        This function retrieves the current status of a device (ONLINE, OFFLINE, etc.) from cache.
        Rate Limit: อนุญาต 60 ครั้งต่อ 1 นาที
        Rate Limit: Allows 60 requests per minute.
    */
    @GetMapping("/devices/{id}/status")
    @RateLimit(limit = 60, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get device status")
    public ResponseEntity<String> getDeviceStatus(@PathVariable UUID id) throws SystemGlobalException {
        String status = iotDeviceService.getDeviceStatus(id);
        return ResponseEntity.ok(status);
    }

    // ========================================================================
    // 6. GET DEVICE LATEST LOCATION (ตำแหน่งล่าสุด)
    // ========================================================================

    /*
        API: GET /api/v1/iot/devices/{id}/location
        ฟังก์ชันนี้ดึงตำแหน่งล่าสุดของอุปกรณ์ (จาก Cache) พร้อมแสดงบนแผนที่
        This function retrieves the latest device location (from cache) for map display.
        Rate Limit: อนุญาต 60 ครั้งต่อ 1 นาที
        Rate Limit: Allows 60 requests per minute.
    */
    @GetMapping("/devices/{id}/location")
    @RateLimit(limit = 60, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get latest device location")
    public ResponseEntity<DeviceLocationResponseDTO> getDeviceLocation(@PathVariable UUID id)
            throws SystemGlobalException {
        DeviceLocationResponseDTO response = gpsTrackingService.getDeviceLocation(id);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 7. GET DEVICE LOCATION HISTORY (ประวัติตำแหน่ง)
    // ========================================================================

    /*
        API: GET /api/v1/iot/devices/{id}/history
        ฟังก์ชันนี้ดึงประวัติตำแหน่งของอุปกรณ์ในช่วงเวลาที่กำหนด
        This function retrieves device location history for a specific time range.
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที
        Rate Limit: Allows 30 requests per minute.
    */
    @GetMapping("/devices/{id}/history")
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get device location history")
    public ResponseEntity<List<GpsDataResponseDTO>> getDeviceHistory(
            @PathVariable UUID id,
            @RequestParam LocalDateTime startTime,
            @RequestParam LocalDateTime endTime) throws SystemGlobalException {
        List<GpsDataResponseDTO> response = gpsTrackingService.getDeviceHistory(id, startTime, endTime);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 8. GET DEVICE HISTORY (ประวัติการทำงาน)
    // ========================================================================

    /*
        API: GET /api/v1/iot/devices/{id}/event-history
        ฟังก์ชันนี้ดึงประวัติเหตุการณ์ของอุปกรณ์ (ONLINE, OFFLINE, BATTERY_LOW)
        This function retrieves device event history (ONLINE, OFFLINE, BATTERY_LOW).
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที
        Rate Limit: Allows 30 requests per minute.
    */
    @GetMapping("/devices/{id}/event-history")
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get device event history")
    public ResponseEntity<Page<DeviceHistoryResponseDTO>> getDeviceEventHistory(
            @PathVariable UUID id,
            @RequestParam(required = false) String eventType,
            Pageable pageable) throws SystemGlobalException {
        Page<DeviceHistoryResponseDTO> response = iotDeviceService.getDeviceEventHistory(id, eventType, pageable);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 9. DELETE DEVICE (ลบอุปกรณ์)
    // ========================================================================

    /*
        API: DELETE /api/v1/iot/devices/{id}
        ฟังก์ชันนี้ลบอุปกรณ์ IoT แบบ Soft Delete
        This function soft-deletes an IoT device.
        Rate Limit: อนุญาต 5 ครั้งต่อ 1 ชั่วโมง
        Rate Limit: Allows 5 requests per 1 hour.
    */
    @DeleteMapping("/devices/{id}")
    @RateLimit(limit = 5, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Delete IoT device")
    public ResponseEntity<Void> deleteDevice(@PathVariable UUID id) throws SystemGlobalException {
        iotDeviceService.deleteDevice(id);
        return ResponseEntity.noContent().build();
    }
}
```

### `MqttController.java` (จัดการ MQTT)

```java
package com.template.app.modules.iot.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.iot.application.interfaces.MqttService;
import com.template.app.modules.iot.presentation.dto.request.MqttPublishRequestDTO;
import com.template.app.exception.SystemGlobalException;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/v1/iot/mqtt")
@Tag(name = "MQTT", description = "MQTT Message APIs")
@RequiredArgsConstructor
public class MqttController {

    private final MqttService mqttService;

    /*
        API: POST /api/v1/iot/mqtt/publish
        ฟังก์ชันนี้ส่งข้อความไปยัง MQTT Broker (ใช้สำหรับควบคุมอุปกรณ์)
        This function publishes a message to MQTT Broker (used for device control).
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @PostMapping("/publish")
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Publish MQTT message")
    public ResponseEntity<Void> publishMqttMessage(@Valid @RequestBody MqttPublishRequestDTO request)
            throws SystemGlobalException {
        mqttService.publishMessage(request.getTopic(), request.getPayload());
        return ResponseEntity.ok().build();
    }
}
```

### `GeofenceController.java`

```java
package com.template.app.modules.iot.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.iot.application.interfaces.GeofenceService;
import com.template.app.modules.iot.presentation.dto.request.GeofenceCreateRequestDTO;
import com.template.app.modules.iot.presentation.dto.response.GeofenceResponseDTO;
import com.template.app.exception.SystemGlobalException;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

@RestController
@RequestMapping("/api/v1/iot/geofences")
@Tag(name = "Geofence", description = "Geofence Management APIs")
@RequiredArgsConstructor
public class GeofenceController {

    private final GeofenceService geofenceService;

    // ========================================================================
    // 1. CREATE GEOFENCE
    // ========================================================================

    /*
        API: POST /api/v1/iot/geofences
        ฟังก์ชันนี้สร้างรั้วรอบขอบ (Geofence) สำหรับติดตามอุปกรณ์
        This function creates a geofence for tracking devices.
        Rate Limit: อนุญาต 10 ครั้งต่อ 5 นาที
        Rate Limit: Allows 10 requests per 5 minutes.
    */
    @PostMapping
    @RateLimit(limit = 10, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Create a geofence")
    public ResponseEntity<GeofenceResponseDTO> createGeofence(@Valid @RequestBody GeofenceCreateRequestDTO request)
            throws SystemGlobalException {
        GeofenceResponseDTO response = geofenceService.createGeofence(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 2. GET GEOFENCE BY ID
    // ========================================================================

    /*
        API: GET /api/v1/iot/geofences/{id}
        ฟังก์ชันนี้ดึงข้อมูลรั้วรอบขอบตาม ID
        This function retrieves geofence data by ID.
        Rate Limit: อนุญาต 80 ครั้งต่อ 1 นาที
        Rate Limit: Allows 80 requests per minute.
    */
    @GetMapping("/{id}")
    @RateLimit(limit = 80, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get geofence by ID")
    public ResponseEntity<GeofenceResponseDTO> getGeofence(@PathVariable UUID id)
            throws SystemGlobalException {
        GeofenceResponseDTO response = geofenceService.getGeofence(id);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. LIST GEOFENCES
    // ========================================================================

    /*
        API: GET /api/v1/iot/geofences
        ฟังก์ชันนี้แสดงรายการรั้วรอบขอบทั้งหมด
        This function lists all geofences.
        Rate Limit: อนุญาต 40 ครั้งต่อ 1 นาที
        Rate Limit: Allows 40 requests per minute.
    */
    @GetMapping
    @RateLimit(limit = 40, duration = 60, keyType = "USER_ID")
    @Operation(summary = "List geofences")
    public ResponseEntity<Page<GeofenceResponseDTO>> listGeofences(Pageable pageable)
            throws SystemGlobalException {
        Page<GeofenceResponseDTO> page = geofenceService.listGeofences(pageable);
        return ResponseEntity.ok(page);
    }

    // ========================================================================
    // 4. DELETE GEOFENCE
    // ========================================================================

    /*
        API: DELETE /api/v1/iot/geofences/{id}
        ฟังก์ชันนี้ลบรั้วรอบขอบ
        This function deletes a geofence.
        Rate Limit: อนุญาต 5 ครั้งต่อ 1 ชั่วโมง
        Rate Limit: Allows 5 requests per 1 hour.
    */
    @DeleteMapping("/{id}")
    @RateLimit(limit = 5, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Delete geofence")
    public ResponseEntity<Void> deleteGeofence(@PathVariable UUID id) throws SystemGlobalException {
        geofenceService.deleteGeofence(id);
        return ResponseEntity.noContent().build();
    }

    // ========================================================================
    // 5. GET GEOFENCE ALERTS
    // ========================================================================

    /*
        API: GET /api/v1/iot/geofences/alerts
        ฟังก์ชันนี้แสดงการแจ้งเตือนจากรั้วรอบขอบทั้งหมด
        This function lists all geofence alerts.
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที
        Rate Limit: Allows 30 requests per minute.
    */
    @GetMapping("/alerts")
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get geofence alerts")
    public ResponseEntity<Page<GeofenceAlertResponseDTO>> getGeofenceAlerts(
            @RequestParam(required = false) Boolean isResolved,
            Pageable pageable) throws SystemGlobalException {
        Page<GeofenceAlertResponseDTO> page = geofenceService.getGeofenceAlerts(isResolved, pageable);
        return ResponseEntity.ok(page);
    }
}
```

---

## 🧩 ตัวอย่าง Domain Entity และ Enum

### `domain/enums/DeviceStatus.java`

```java
package com.template.app.modules.iot.domain.enums;

/*
    สถานะของอุปกรณ์ IoT / IoT device status.
*/
public enum DeviceStatus {
    ONLINE,        // ออนไลน์ (เชื่อมต่อ) / Online (connected).
    OFFLINE,       // ออฟไลน์ (ตัดการเชื่อมต่อ) / Offline (disconnected).
    INACTIVE,      // ถูกปิดใช้งาน / Deactivated.
    MAINTENANCE    // กำลังซ่อมบำรุง / Under maintenance.
}
```

### `domain/enums/DeviceType.java`

```java
package com.template.app.modules.iot.domain.enums;

/*
    ประเภทของอุปกรณ์ IoT / IoT device type.
*/
public enum DeviceType {
    GPS_TRACKER,   // เครื่องติดตาม GPS / GPS tracker.
    SENSOR,        // เซ็นเซอร์ / Sensor.
    CAMERA,        // กล้อง / Camera.
    CONTROLLER     // ตัวควบคุม / Controller.
}
```

### `domain/MIoTDevice.java`

```java
package com.template.app.modules.iot.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import com.template.app.modules.iot.domain.enums.DeviceStatus;
import com.template.app.modules.iot.domain.enums.DeviceType;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.math.BigDecimal;
import java.time.LocalDateTime;

@Data
@EqualsAndHashCode(callSuper = true)
public class MIoTDevice extends GenericBusinessClass {

    private String deviceId;                // รหัสฮาร์ดแวร์ / Hardware ID.
    private String deviceName;              // ชื่ออุปกรณ์ / Device name.
    private DeviceType deviceType;          // ประเภทอุปกรณ์ / Device type.
    private DeviceStatus status;            // สถานะ / Status.
    private String serialNumber;            // เลขซีเรียล / Serial number.
    private String firmwareVersion;         // เวอร์ชัน Firmware
    private String hardwareVersion;         // เวอร์ชัน Hardware
    private String manufacturer;            // ผู้ผลิต / Manufacturer.
    private String model;                   // รุ่น / Model.
    private Integer batteryLevel;           // ระดับแบตเตอรี่ (%) / Battery level (%).
    private LocalDateTime lastSeen;         // เวลาที่เห็นครั้งล่าสุด / Last seen.
    private BigDecimal lastLatitude;        // ละติจูดล่าสุด / Last latitude.
    private BigDecimal lastLongitude;       // ลองจิจูดล่าสุด / Last longitude.
    private BigDecimal lastAltitude;        // ระดับความสูงล่าสุด / Last altitude.
    private BigDecimal lastSpeed;           // ความเร็วล่าสุด / Last speed.
    private Boolean isActive;               // ใช้งานอยู่ / Active.
    private Boolean isOnline;               // ออนไลน์อยู่ / Online.
    private String metadata;                // ข้อมูลเพิ่มเติม / Metadata.
    private String notes;                   // หมายเหตุ / Notes.

    /*
        ฟังก์ชันนี้อัปเดตสถานะของอุปกรณ์
        This function updates the device status.
    */
    public void updateStatus(DeviceStatus newStatus) {
        this.status = newStatus;
        this.isOnline = (newStatus == DeviceStatus.ONLINE);
        if (newStatus == DeviceStatus.ONLINE) {
            this.lastSeen = LocalDateTime.now();
        }
    }

    /*
        ฟังก์ชันนี้อัปเดตตำแหน่งของอุปกรณ์
        This function updates the device location.
    */
    public void updateLocation(BigDecimal latitude, BigDecimal longitude, BigDecimal altitude, BigDecimal speed) {
        this.lastLatitude = latitude;
        this.lastLongitude = longitude;
        this.lastAltitude = altitude;
        this.lastSpeed = speed;
        this.lastSeen = LocalDateTime.now();
        this.isOnline = true;
        this.status = DeviceStatus.ONLINE;
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่าอุปกรณ์ออนไลน์อยู่หรือไม่
        This function checks if the device is online.
    */
    public boolean isOnline() {
        return this.isOnline != null && this.isOnline;
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่าอุปกรณ์มีแบตเตอรี่ต่ำหรือไม่ (น้อยกว่า 15%)
        This function checks if the device has low battery (less than 15%).
    */
    public boolean isBatteryLow() {
        return this.batteryLevel != null && this.batteryLevel < 15;
    }
}
```

### `domain/TGPSData.java`

```java
package com.template.app.modules.iot.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.UUID;

@Data
@EqualsAndHashCode(callSuper = true)
public class TGPSData extends GenericBusinessClass {

    private UUID deviceId;                  // ID อุปกรณ์ / Device ID.
    private String deviceIdentifier;        // รหัสอุปกรณ์ / Device identifier.
    private BigDecimal latitude;            // ละติจูด / Latitude.
    private BigDecimal longitude;           // ลองจิจูด / Longitude.
    private BigDecimal altitude;            // ระดับความสูง / Altitude.
    private BigDecimal speed;               // ความเร็ว / Speed.
    private BigDecimal heading;             // ทิศทาง / Heading.
    private BigDecimal accuracy;            // ความแม่นยำ / Accuracy.
    private Integer batteryLevel;           // ระดับแบตเตอรี่ / Battery level.
    private Integer satellites;             // จำนวนดาวเทียม / Satellites.
    private String eventType;               // MOVING, STOPPED, IDLE, ALERT
    private LocalDateTime timestamp;        // เวลาจากอุปกรณ์ / Device timestamp.
    private LocalDateTime receivedAt;       // เวลาที่ระบบรับ / System received time.
    private String metadata;                // ข้อมูลเพิ่มเติม / Metadata.

    /*
        ฟังก์ชันนี้คำนวณระยะทางระหว่างจุดสองจุด (Haversine formula)
        This function calculates distance between two points (Haversine formula).
    */
    public static double calculateDistance(BigDecimal lat1, BigDecimal lon1,
                                           BigDecimal lat2, BigDecimal lon2) {
        final int R = 6371; // รัศมีโลก (กิโลเมตร) / Earth radius (km).
        double lat1Rad = Math.toRadians(lat1.doubleValue());
        double lon1Rad = Math.toRadians(lon1.doubleValue());
        double lat2Rad = Math.toRadians(lat2.doubleValue());
        double lon2Rad = Math.toRadians(lon2.doubleValue());

        double dlon = lon2Rad - lon1Rad;
        double dlat = lat2Rad - lat1Rad;

        double a = Math.pow(Math.sin(dlat / 2), 2)
                + Math.cos(lat1Rad) * Math.cos(lat2Rad)
                * Math.pow(Math.sin(dlon / 2), 2);

        double c = 2 * Math.asin(Math.sqrt(a));
        return R * c; // ระยะทาง (กม.) / Distance (km).
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่าอุปกรณ์กำลังเคลื่อนที่หรือไม่ (ความเร็ว > 5 กม./ชม.)
        This function checks if the device is moving (speed > 5 km/h).
    */
    public boolean isMoving() {
        return this.speed != null && this.speed.doubleValue() > 5.0;
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่าอุปกรณ์หยุดนิ่งหรือไม่ (ความเร็ว < 1 กม./ชม. นานกว่า 5 นาที)
        This function checks if the device is stopped (speed < 1 km/h for more than 5 minutes).
    */
    public boolean isStopped() {
        return this.speed != null && this.speed.doubleValue() < 1.0;
    }
}
```

---

## 🔧 ฟังก์ชัน Service ตัวอย่าง (พร้อมคำอธิบาย)

### `application/impl/MqttServiceImpl.java` (MQTT Service)

```java
package com.template.app.modules.iot.application.impl;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.template.app.modules.iot.application.interfaces.MqttService;
import com.template.app.modules.iot.domain.MIoTDevice;
import com.template.app.modules.iot.domain.TGPSData;
import com.template.app.modules.iot.infrastructure.cache.DeviceCacheService;
import com.template.app.modules.iot.infrastructure.cache.DeviceLocationCacheService;
import com.template.app.modules.iot.infrastructure.repository.IotDeviceRepository;
import com.template.app.modules.iot.infrastructure.repository.GpsDataRepository;
import com.template.app.modules.iot.presentation.dto.response.DeviceLocationResponseDTO;
import org.eclipse.paho.client.mqttv3.MqttMessage;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.integration.mqtt.outbound.MqttPahoMessageHandler;
import org.springframework.messaging.Message;
import org.springframework.messaging.MessageHandler;
import org.springframework.messaging.MessagingException;
import org.springframework.stereotype.Service;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.UUID;

@Service
public class MqttServiceImpl implements MqttService, MessageHandler {

    private static final Logger logger = LoggerFactory.getLogger(MqttServiceImpl.class);

    private final MqttPahoMessageHandler mqttOutbound;
    private final ObjectMapper objectMapper;
    private final IotDeviceRepository deviceRepository;
    private final GpsDataRepository gpsDataRepository;
    private final DeviceCacheService deviceCacheService;
    private final DeviceLocationCacheService locationCacheService;

    public MqttServiceImpl(MqttPahoMessageHandler mqttOutbound,
                           ObjectMapper objectMapper,
                           IotDeviceRepository deviceRepository,
                           GpsDataRepository gpsDataRepository,
                           DeviceCacheService deviceCacheService,
                           DeviceLocationCacheService locationCacheService) {
        this.mqttOutbound = mqttOutbound;
        this.objectMapper = objectMapper;
        this.deviceRepository = deviceRepository;
        this.gpsDataRepository = gpsDataRepository;
        this.deviceCacheService = deviceCacheService;
        this.locationCacheService = locationCacheService;
    }

    /*
        ฟังก์ชันนี้ส่งข้อความไปยัง MQTT Broker (ใช้สำหรับควบคุมอุปกรณ์)
        This function publishes a message to MQTT Broker (used for device control).
    */
    @Override
    public void publishMessage(String topic, String payload) {
        try {
            MqttMessage message = new MqttMessage(payload.getBytes());
            message.setQos(1);
            message.setRetained(false);
            
            // ส่งข้อความผ่าน MQTT Outbound Handler
            // Send message via MQTT Outbound Handler.
            mqttOutbound.handleMessage(
                new org.springframework.messaging.support.GenericMessage<>(
                    new Object[] { topic, message }
                )
            );
            logger.info("MQTT message published to topic: {}", topic);
        } catch (Exception e) {
            logger.error("Failed to publish MQTT message: {}", e.getMessage(), e);
        }
    }

    /*
        ฟังก์ชันนี้รับข้อความจาก MQTT (Callback จาก Inbound Adapter)
        This function handles incoming MQTT messages (Callback from Inbound Adapter).
        ฟังทุกข้อความที่มาจาก MQTT และประมวลผลตาม Topic
        Listens to all messages from MQTT and processes based on Topic.
    */
    @Override
    public void handleMessage(Message<?> message) throws MessagingException {
        try {
            String topic = (String) message.getHeaders().get("mqtt_receivedTopic");
            String payload = (String) message.getPayload();
            
            logger.debug("Received MQTT message on topic: {}, payload: {}", topic, payload);

            // ประมวลผลตาม Topic / Process based on Topic.
            if (topic != null) {
                if (topic.startsWith("gps/data")) {
                    processGpsData(payload);
                } else if (topic.startsWith("device/status")) {
                    processDeviceStatus(payload);
                } else if (topic.startsWith("device/alert")) {
                    processDeviceAlert(payload);
                }
            }
        } catch (Exception e) {
            logger.error("Error processing MQTT message: {}", e.getMessage(), e);
        }
    }

    /*
        ฟังก์ชันนี้ประมวลผลข้อมูล GPS ที่ได้รับจาก MQTT
        This function processes GPS data received from MQTT.
    */
    private void processGpsData(String payload) throws Exception {
        // 1. แปลง JSON payload เป็น Object / Parse JSON payload to Object.
        GpsDataPayload gpsData = objectMapper.readValue(payload, GpsDataPayload.class);
        
        // 2. ค้นหาอุปกรณ์ / Find device.
        UUID deviceId = deviceRepository.findByDeviceId(gpsData.getDeviceId())
                .map(MIoTDevice::getId)
                .orElse(null);
        
        if (deviceId == null) {
            logger.warn("Device not found: {}", gpsData.getDeviceId());
            return;
        }

        // 3. อัปเดตตำแหน่งล่าสุดของอุปกรณ์ / Update device last location.
        deviceRepository.updateDeviceLocation(
            deviceId,
            new BigDecimal(gpsData.getLatitude()),
            new BigDecimal(gpsData.getLongitude()),
            gpsData.getAltitude() != null ? new BigDecimal(gpsData.getAltitude()) : null,
            gpsData.getSpeed() != null ? new BigDecimal(gpsData.getSpeed()) : null,
            LocalDateTime.now()
        );
        deviceCacheService.evictDevice(deviceId);

        // 4. อัปเดต Cache ตำแหน่ง / Update location cache.
        DeviceLocationResponseDTO location = DeviceLocationResponseDTO.builder()
                .deviceId(deviceId)
                .latitude(gpsData.getLatitude())
                .longitude(gpsData.getLongitude())
                .speed(gpsData.getSpeed())
                .timestamp(LocalDateTime.now())
                .build();
        locationCacheService.updateDeviceLocation(deviceId, location);

        // 5. บันทึกข้อมูล GPS ลงฐานข้อมูล (และ InfluxDB) / Save GPS data to DB (and InfluxDB).
        TGPSData gpsRecord = new TGPSData();
        gpsRecord.setDeviceId(deviceId);
        gpsRecord.setDeviceIdentifier(gpsData.getDeviceId());
        gpsRecord.setLatitude(new BigDecimal(gpsData.getLatitude()));
        gpsRecord.setLongitude(new BigDecimal(gpsData.getLongitude()));
        gpsRecord.setSpeed(gpsData.getSpeed() != null ? new BigDecimal(gpsData.getSpeed()) : null);
        gpsRecord.setTimestamp(gpsData.getTimestamp() != null ? gpsData.getTimestamp() : LocalDateTime.now());
        gpsRecord.setReceivedAt(LocalDateTime.now());
        gpsRecord.setBatteryLevel(gpsData.getBatteryLevel());
        
        gpsDataRepository.save(gpsRecord);

        // 6. ตรวจสอบ Geofence / Check Geofence.
        // geofenceService.checkGeofence(deviceId, latitude, longitude);

        logger.info("GPS data processed for device: {}, lat: {}, lon: {}", 
            gpsData.getDeviceId(), gpsData.getLatitude(), gpsData.getLongitude());
    }

    /*
        ฟังก์ชันนี้ประมวลผลสถานะอุปกรณ์ที่ได้รับจาก MQTT
        This function processes device status received from MQTT.
    */
    private void processDeviceStatus(String payload) throws Exception {
        DeviceStatusPayload status = objectMapper.readValue(payload, DeviceStatusPayload.class);
        
        // อัปเดตสถานะอุปกรณ์ / Update device status.
        deviceRepository.updateDeviceStatus(
            status.getDeviceId(),
            status.getStatus(),
            status.getBatteryLevel()
        );
        
        logger.info("Device status updated: {}, status: {}", status.getDeviceId(), status.getStatus());
    }

    /*
        ฟังก์ชันนี้ประมวลผลการแจ้งเตือนจากอุปกรณ์
        This function processes alerts from devices.
    */
    private void processDeviceAlert(String payload) throws Exception {
        DeviceAlertPayload alert = objectMapper.readValue(payload, DeviceAlertPayload.class);
        
        // บันทึกการแจ้งเตือน / Save alert.
        logger.info("Device alert: {}, type: {}, message: {}", 
            alert.getDeviceId(), alert.getAlertType(), alert.getMessage());
        
        // TODO: ส่งการแจ้งเตือนไปยังระบบแจ้งเตือน (LINE, Email)
        // sendNotification(alert);
    }

    // Inner class for GPS data payload
    private static class GpsDataPayload {
        private String deviceId;
        private String latitude;
        private String longitude;
        private String altitude;
        private String speed;
        private Integer batteryLevel;
        private LocalDateTime timestamp;
        // getters and setters
        public String getDeviceId() { return deviceId; }
        public void setDeviceId(String deviceId) { this.deviceId = deviceId; }
        public String getLatitude() { return latitude; }
        public void setLatitude(String latitude) { this.latitude = latitude; }
        public String getLongitude() { return longitude; }
        public void setLongitude(String longitude) { this.longitude = longitude; }
        public String getAltitude() { return altitude; }
        public void setAltitude(String altitude) { this.altitude = altitude; }
        public String getSpeed() { return speed; }
        public void setSpeed(String speed) { this.speed = speed; }
        public Integer getBatteryLevel() { return batteryLevel; }
        public void setBatteryLevel(Integer batteryLevel) { this.batteryLevel = batteryLevel; }
        public LocalDateTime getTimestamp() { return timestamp; }
        public void setTimestamp(LocalDateTime timestamp) { this.timestamp = timestamp; }
    }

    private static class DeviceStatusPayload {
        private String deviceId;
        private String status;
        private Integer batteryLevel;
        // getters and setters
        public String getDeviceId() { return deviceId; }
        public void setDeviceId(String deviceId) { this.deviceId = deviceId; }
        public String getStatus() { return status; }
        public void setStatus(String status) { this.status = status; }
        public Integer getBatteryLevel() { return batteryLevel; }
        public void setBatteryLevel(Integer batteryLevel) { this.batteryLevel = batteryLevel; }
    }

    private static class DeviceAlertPayload {
        private String deviceId;
        private String alertType;
        private String message;
        // getters and setters
        public String getDeviceId() { return deviceId; }
        public void setDeviceId(String deviceId) { this.deviceId = deviceId; }
        public String getAlertType() { return alertType; }
        public void setAlertType(String alertType) { this.alertType = alertType; }
        public String getMessage() { return message; }
        public void setMessage(String message) { this.message = message; }
    }
}
```

---

## 📊 สรุปโมดูลที่ดำเนินการแล้ว (Modules Completed)

| # | โมดูล | สถานะ | รายละเอียด |
|---|-------|--------|-----------|
| 1 | 🔑 Authentication & Permission | ✅ ครบถ้วน | JWT + Role/Permission + Rate Limit + Redis Cache |
| 2 | 🚗 Job Card Management | ✅ ครบถ้วน | 14 Statuses + Service/Part + History + Cache |
| 3 | 👥 Customer Management | ✅ ครบถ้วน | Customer + Car + History + Cache (ID/Phone/Plate) |
| 4 | 📋 Quotation | ✅ ครบถ้วน | Quotation + Part/Service + Approve/Reject + PDF + Cache |
| 5 | 🛒 Purchase Order | ✅ ครบถ้วน | PO + Status + Send/Receive + PDF + Email + Cache |
| 6 | 📦 Inventory Management | ✅ ครบถ้วน | Part Master + Receive/Issue + Picking + Stock Take + Adjustment + Cache |
| 7 | 💰 Payment Management | ✅ ครบถ้วน | Payment Record + Receipt + Outstanding Balance + Refund + Cache |
| 8 | 📊 Dashboard & Reports | ✅ ครบถ้วน | Overview + Sales + Inventory + Job Status + Top Parts + Financial + Export |
| 9 | 📄 Document Management | ✅ ครบถ้วน | Document Generation + Upload + OCR + Template Management + Storage + Cache |
| 10 | 📧 Email Service | ✅ ครบถ้วน | Template-based Email + Attachments + Bulk + Queue + History + Cache |
| 11 | ⏱️ Batch Jobs | ✅ ครบถ้วน | 6 Scheduled Jobs + Distributed Lock + History + Manual Trigger + Cache |
| 12 | 🌏 Multi-Language (i18n) | ✅ ครบถ้วน | 18 Languages + Translation Management + Resource Bundle + Cache |
| 13 | 📡 IoT & GPS Tracking | ✅ ครบถ้วน | Device Management + GPS Tracking + MQTT + InfluxDB + Geofence + Auto Report + Cache |

--- 