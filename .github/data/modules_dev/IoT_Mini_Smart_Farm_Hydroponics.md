# ระบบ Mini Smart Farm สำหรับ Hydroponics

## สารบัญ

1. [ภาพรวมระบบ](#1-ภาพรวมระบบ)
2. [โครงสร้างระบบ](#2-โครงสร้างระบบ)
3. [Hardware Components](#3-hardware-components)
4. [Software Architecture](#4-software-architecture)
5. [Domain Layer](#5-domain-layer)
6. [Application Layer](#6-application-layer)
7. [Infrastructure Layer](#7-infrastructure-layer)
8. [API Endpoints](#8-api-endpoints)
9. [Mobile App Integration](#9-mobile-app-integration)
10. [Dashboard UI](#10-dashboard-ui)
11. [การติดตั้งและใช้งาน](#11-การติดตั้งและใช้งาน)
12. [Business Model](#12-business-model)

---

## 1. ภาพรวมระบบ

### 1.1 ระบบ Mini Smart Farm Hydroponics

ระบบ Mini Smart Farm เป็นโซลูชันสำหรับการปลูกพืชแบบไฮโดรโพนิกส์อัจฉริยะ ขนาดเล็กถึงกลาง เหมาะสำหรับ:

- **ร้านอาหาร/โรงแรม** - ปลูกผักสดใช้ในครัว
- **โรงเรียน/มหาวิทยาลัย** - สื่อการเรียนรู้เกษตรสมัยใหม่
- **คอนโด/บ้านพักอาศัย** - ปลูกผักทานเอง
- **สตาร์ทอัพเกษตร** - เริ่มต้นธุรกิจผักไฮโดรโพนิกส์
- **ชุมชน/สหกรณ์** - กลุ่มเกษตรกรรุ่นใหม่

### 1.2 ฟีเจอร์หลัก

| ฟีเจอร์ | รายละเอียด |
|---------|------------|
| **ควบคุมอัตโนมัติ** | ควบคุมปั๊มน้ำ, แสง, พัดลม, วาล์วอัตโนมัติ |
| **ตรวจวัดสภาพแวดล้อม** | อุณหภูมิ, ความชื้น, pH, EC, น้ำ, แสง, CO2 |
| **ตั้งเวลาทำงาน** | ตั้งเวลาเปิด-ปิดอุปกรณ์ตามรอบการปลูก |
| **แจ้งเตือน** | แจ้งเตือนเมื่อค่าน้ำ/สารอาหารผิดปกติ |
| **ดูข้อมูลย้อนหลัง** | เก็บข้อมูลและแสดงกราฟย้อนหลัง |
| **AI ช่วยแนะนำ** | แนะนำการปรับค่า pH/EC และการปลูก |
| **ควบคุมระยะไกล** | ควบคุมผ่านแอปพลิเคชันมือถือ |
| **Multi-farm** | รองรับหลายแปลงปลูกพร้อมกัน |

### 1.3 ระบบย่อย

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      Mini Smart Farm System                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                        │
│  ┌─────────────────────────────────────────────────────────────────┐  │
│  │                       User Layer                                │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │  │
│  │  │  Mobile  │  │   Web    │  │  Tablet  │  │  Smart   │    │  │
│  │  │    App   │  │ Dashboard│  │   App    │  │  Watch   │    │  │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │  │
│  └─────────────────────────────────────────────────────────────────┘  │
│                                    │                                    │
│  ┌─────────────────────────────────▼─────────────────────────────────┐ │
│  │                      Cloud Platform                               │ │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │  │
│  │  │   API    │  │  MQTT    │  │  WebSocket│  │  AI/ML   │    │  │
│  │  │ Gateway  │  │ Broker   │  │  Server  │  │  Engine  │    │  │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │  │
│  └─────────────────────────────────────────────────────────────────┘  │
│                                    │                                    │
│  ┌─────────────────────────────────▼─────────────────────────────────┐ │
│  │                      Edge Gateway                                 │ │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │  │
│  │  │ ESP32/   │  │  Wi-Fi/  │  │  Sensor  │  │  Relay   │    │  │
│  │  │ Arduino  │  │  LoRa    │  │  Hub     │  │  Control │    │  │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │  │
│  └─────────────────────────────────────────────────────────────────┘  │
│                                    │                                    │
│  ┌─────────────────────────────────▼─────────────────────────────────┐ │
│  │                      Hardware Layer                               │ │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │  │
│  │  │  Sensors │  │  Pumps   │  │  Lights  │  │  Fans    │    │  │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │  │
│  └─────────────────────────────────────────────────────────────────┘  │
│                                                                        │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 2. โครงสร้างระบบ

### 2.1 สถาปัตยกรรมโมดูล

```
internal/modules/hydroponics/
├── domain/
│   ├── entity/
│   │   ├── farm.go           # ข้อมูลฟาร์ม
│   │   ├── sensor.go         # ข้อมูลเซ็นเซอร์
│   │   ├── actuator.go       # ข้อมูลอุปกรณ์ควบคุม
│   │   ├── schedule.go       # ตารางเวลา
│   │   ├── alert.go          # การแจ้งเตือน
│   │   ├── crop.go           # ข้อมูลพืช
│   │   └── nutrient.go       # สูตรสารอาหาร
│   ├── value_object/
│   │   ├── sensor_type.go
│   │   ├── actuator_type.go
│   │   ├── crop_status.go
│   │   └── alert_severity.go
│   ├── repository/
│   │   ├── farm_repository.go
│   │   ├── sensor_repository.go
│   │   └── alert_repository.go
│   └── service/
│       ├── farm_service.go
│       ├── automation_service.go
│       └── recommendation_service.go
├── application/
│   ├── create_farm.go
│   ├── get_sensor_data.go
│   ├── control_actuator.go
│   ├── create_schedule.go
│   ├── get_alerts.go
│   └── get_recommendation.go
├── infrastructure/
│   ├── mqtt/
│   │   └── mqtt_client.go
│   ├── persistence/
│   │   ├── postgres/
│   │   │   ├── models.go
│   │   │   └── farm_repo_impl.go
│   │   └── influxdb/
│   │       └── sensor_writer.go
│   ├── ai/
│   │   ├── recommendation.go
│   │   └── anomaly_detection.go
│   └── notification/
│       ├── line_notify.go
│       └── email_notify.go
└── interfaces/
    ├── http/
    │   ├── farm_handler.go
    │   └── routes.go
    ├── websocket/
    │   └── realtime_handler.go
    └── mqtt/
        └── mqtt_handler.go
```

---

## 3. Hardware Components

### 3.1 รายการอุปกรณ์

| Component | Model | จำนวน | ราคา (บาท) |
|-----------|-------|-------|------------|
| **ESP32** | ESP32-WROOM | 1 | 350 |
| **Sensor pH** | Atlas Scientific EZO-pH | 1 | 2,500 |
| **Sensor EC** | Atlas Scientific EZO-EC | 1 | 2,500 |
| **Sensor Temperature** | DS18B20 | 2 | 200 |
| **Sensor Humidity** | DHT22 | 1 | 150 |
| **TDS Sensor** | Gravity TDS | 1 | 350 |
| **Water Level** | Ultrasonic HC-SR04 | 1 | 100 |
| **Light Sensor** | BH1750 | 1 | 150 |
| **Relay Module** | 8-Channel Relay | 1 | 250 |
| **Water Pump** | DC 12V | 2 | 300 |
| **LED Grow Light** | 100W Full Spectrum | 1 | 800 |
| **Fan** | DC 12V | 2 | 200 |
| **Water Valve** | Solenoid 12V | 2 | 250 |
| **Power Supply** | 12V 5A | 1 | 300 |
| **LCD Display** | 16x2 I2C | 1 | 150 |
| **รวม** | | | **~8,650** |

### 3.2 การเชื่อมต่อเซ็นเซอร์

```
ESP32 Pin Map:
├── GPIO 4  → DHT22 (Temperature + Humidity)
├── GPIO 5  → DS18B20 (Water Temperature)
├── GPIO 18 → pH Sensor (UART)
├── GPIO 19 → EC Sensor (UART)
├── GPIO 21 → I2C SDA (LCD, BH1750)
├── GPIO 22 → I2C SCL (LCD, BH1750)
├── GPIO 23 → TDS Sensor (Analog)
├── GPIO 32 → HC-SR04 (Trig)
├── GPIO 33 → HC-SR04 (Echo)
├── GPIO 26 → Relay 1 (Pump 1)
├── GPIO 27 → Relay 2 (Pump 2)
├── GPIO 14 → Relay 3 (Light)
├── GPIO 12 → Relay 4 (Fan 1)
├── GPIO 13 → Relay 5 (Fan 2)
├── GPIO 25 → Relay 6 (Valve 1)
└── GPIO 26 → Relay 7 (Valve 2)
```

### 3.3 Firmware (ESP32 - Arduino)

```cpp
// hydroponics_firmware.ino
#include <WiFi.h>
#include <PubSubClient.h>
#include <DHT.h>
#include <OneWire.h>
#include <DallasTemperature.h>
#include <Wire.h>
#include <LiquidCrystal_I2C.h>
#include <BH1750.h>

// WiFi Configuration
const char* ssid = "YOUR_WIFI_SSID";
const char* password = "YOUR_WIFI_PASSWORD";
const char* mqtt_server = "mqtt.yourserver.com";
const char* mqtt_topic = "hydroponics/farm01";

// Pin Definitions
#define DHTPIN 4
#define DHTTYPE DHT22
#define ONE_WIRE_BUS 5
#define TDS_PIN 23
#define TRIG_PIN 32
#define ECHO_PIN 33

// Relay Pins
#define PUMP1_PIN 26
#define PUMP2_PIN 27
#define LIGHT_PIN 14
#define FAN1_PIN 12
#define FAN2_PIN 13
#define VALVE1_PIN 25
#define VALVE2_PIN 26

// Objects
DHT dht(DHTPIN, DHTTYPE);
OneWire oneWire(ONE_WIRE_BUS);
DallasTemperature waterTemp(&oneWire);
LiquidCrystal_I2C lcd(0x27, 16, 2);
BH1750 lightMeter;

WiFiClient espClient;
PubSubClient client(espClient);

// Variables
float pH = 0.0;
float ec = 0.0;
float tds = 0.0;
float temperature = 0.0;
float humidity = 0.0;
float water_temp = 0.0;
float light = 0.0;
float water_level = 0.0;

unsigned long lastSensorRead = 0;
const unsigned long sensorInterval = 5000;

// Function prototypes
void setup_wifi();
void reconnect_mqtt();
void read_sensors();
void control_actuators();
void publish_data();

void setup() {
    Serial.begin(115200);
    
    // Init sensors
    dht.begin();
    waterTemp.begin();
    Wire.begin();
    lcd.init();
    lcd.backlight();
    lightMeter.begin();
    
    // Init pins
    pinMode(PUMP1_PIN, OUTPUT);
    pinMode(PUMP2_PIN, OUTPUT);
    pinMode(LIGHT_PIN, OUTPUT);
    pinMode(FAN1_PIN, OUTPUT);
    pinMode(FAN2_PIN, OUTPUT);
    pinMode(VALVE1_PIN, OUTPUT);
    pinMode(VALVE2_PIN, OUTPUT);
    pinMode(TRIG_PIN, OUTPUT);
    pinMode(ECHO_PIN, INPUT);
    
    // Turn off all
    digitalWrite(PUMP1_PIN, LOW);
    digitalWrite(PUMP2_PIN, LOW);
    digitalWrite(LIGHT_PIN, LOW);
    digitalWrite(FAN1_PIN, LOW);
    digitalWrite(FAN2_PIN, LOW);
    digitalWrite(VALVE1_PIN, LOW);
    digitalWrite(VALVE2_PIN, LOW);
    
    setup_wifi();
    client.setServer(mqtt_server, 1883);
    client.setCallback(mqtt_callback);
    
    lcd.setCursor(0, 0);
    lcd.print("Hydroponics");
    lcd.setCursor(0, 1);
    lcd.print("Ready...");
    
    delay(2000);
}

void setup_wifi() {
    delay(10);
    Serial.println();
    Serial.print("Connecting to ");
    Serial.println(ssid);
    
    WiFi.begin(ssid, password);
    
    while (WiFi.status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }
    
    Serial.println();
    Serial.println("WiFi connected");
    Serial.print("IP address: ");
    Serial.println(WiFi.localIP());
}

void reconnect_mqtt() {
    while (!client.connected()) {
        Serial.print("Attempting MQTT connection...");
        if (client.connect("ESP32Client01")) {
            Serial.println("connected");
            client.subscribe("hydroponics/farm01/control");
        } else {
            Serial.print("failed, rc=");
            Serial.print(client.state());
            Serial.println(" try again in 5 seconds");
            delay(5000);
        }
    }
}

void read_sensors() {
    // Read DHT22
    temperature = dht.readTemperature();
    humidity = dht.readHumidity();
    
    // Read water temperature
    waterTemp.requestTemperatures();
    water_temp = waterTemp.getTempCByIndex(0);
    
    // Read pH (simulated - actual uses UART)
    pH = 6.5 + (random(0, 100) / 100.0);
    
    // Read EC (simulated)
    ec = 1.2 + (random(0, 50) / 100.0);
    
    // Read TDS
    tds = ec * 500.0;
    
    // Read light
    light = lightMeter.readLightLevel();
    
    // Read water level
    digitalWrite(TRIG_PIN, LOW);
    delayMicroseconds(2);
    digitalWrite(TRIG_PIN, HIGH);
    delayMicroseconds(10);
    digitalWrite(TRIG_PIN, LOW);
    long duration = pulseIn(ECHO_PIN, HIGH);
    water_level = duration * 0.034 / 2;
}

void control_actuators() {
    // Auto control based on sensor values
    // pH control - if pH too high, add acid
    if (pH > 7.0) {
        digitalWrite(VALVE1_PIN, HIGH); // Acid valve
    } else {
        digitalWrite(VALVE1_PIN, LOW);
    }
    
    // EC control - if EC too low, add nutrient
    if (ec < 1.0) {
        digitalWrite(VALVE2_PIN, HIGH); // Nutrient valve
    } else {
        digitalWrite(VALVE2_PIN, LOW);
    }
    
    // Temperature control
    if (temperature > 30.0) {
        digitalWrite(FAN1_PIN, HIGH);
        digitalWrite(FAN2_PIN, HIGH);
    } else {
        digitalWrite(FAN1_PIN, LOW);
        digitalWrite(FAN2_PIN, LOW);
    }
    
    // Light control - based on time or ambient light
    if (light < 100 && (millis() % 86400000) > 21600000) {
        digitalWrite(LIGHT_PIN, HIGH);
    } else {
        digitalWrite(LIGHT_PIN, LOW);
    }
}

void publish_data() {
    char buffer[256];
    
    snprintf(buffer, sizeof(buffer),
        "{\"device\":\"farm01\",\"data\":{"
        "\"temperature\":%.2f,"
        "\"humidity\":%.2f,"
        "\"water_temp\":%.2f,"
        "\"ph\":%.2f,"
        "\"ec\":%.2f,"
        "\"tds\":%.2f,"
        "\"light\":%.2f,"
        "\"water_level\":%.2f,"
        "\"timestamp\":%lu"
        "}}",
        temperature, humidity, water_temp,
        pH, ec, tds, light, water_level,
        millis() / 1000
    );
    
    client.publish(mqtt_topic, buffer);
    
    // Update LCD
    lcd.clear();
    lcd.setCursor(0, 0);
    lcd.printf("pH:%.2f EC:%.2f", pH, ec);
    lcd.setCursor(0, 1);
    lcd.printf("T:%.1fC H:%.0f%%", temperature, humidity);
}

void mqtt_callback(char* topic, byte* payload, unsigned int length) {
    String message;
    for (int i = 0; i < length; i++) {
        message += (char)payload[i];
    }
    
    // Parse control commands
    if (message == "pump1_on") {
        digitalWrite(PUMP1_PIN, HIGH);
    } else if (message == "pump1_off") {
        digitalWrite(PUMP1_PIN, LOW);
    } else if (message == "pump2_on") {
        digitalWrite(PUMP2_PIN, HIGH);
    } else if (message == "pump2_off") {
        digitalWrite(PUMP2_PIN, LOW);
    } else if (message == "light_on") {
        digitalWrite(LIGHT_PIN, HIGH);
    } else if (message == "light_off") {
        digitalWrite(LIGHT_PIN, LOW);
    }
}

void loop() {
    if (!client.connected()) {
        reconnect_mqtt();
    }
    client.loop();
    
    unsigned long currentMillis = millis();
    if (currentMillis - lastSensorRead >= sensorInterval) {
        lastSensorRead = currentMillis;
        
        read_sensors();
        control_actuators();
        publish_data();
    }
}
```

---

## 4. Software Architecture

### 4.1 Go Backend Structure

```go
// cmd/hydroponics/main.go
package main

import (
    "context"
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/yourproject/internal/modules/hydroponics/infrastructure/mqtt"
    "github.com/yourproject/internal/modules/hydroponics/infrastructure/persistence/postgres"
    "github.com/yourproject/internal/modules/hydroponics/infrastructure/persistence/influxdb"
    "github.com/yourproject/internal/modules/hydroponics/infrastructure/ai"
    "github.com/yourproject/internal/modules/hydroponics/interfaces/http"
    "github.com/yourproject/internal/modules/hydroponics/interfaces/websocket"
)

func main() {
    ctx := context.Background()

    // 1. PostgreSQL
    db := postgres.NewDB(postgres.Config{
        Host:     os.Getenv("DB_HOST"),
        Port:     os.Getenv("DB_PORT"),
        User:     os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASSWORD"),
        DBName:   os.Getenv("DB_NAME"),
    })

    // 2. InfluxDB
    influx := influxdb.NewClient(influxdb.Config{
        URL:    os.Getenv("INFLUX_URL"),
        Token:  os.Getenv("INFLUX_TOKEN"),
        Org:    os.Getenv("INFLUX_ORG"),
        Bucket: os.Getenv("INFLUX_BUCKET"),
    })

    // 3. MQTT
    mqttClient := mqtt.NewClient(mqtt.Config{
        Broker:   os.Getenv("MQTT_BROKER"),
        ClientID: "hydroponics_backend",
    })

    // 4. AI Service
    aiService := ai.NewRecommendationService()

    // 5. WebSocket
    wsServer := websocket.NewServer()

    // 6. Repositories
    farmRepo := postgres.NewFarmRepository(db)
    sensorRepo := influxdb.NewSensorRepository(influx)
    alertRepo := postgres.NewAlertRepository(db)

    // 7. Services
    farmService := service.NewFarmService(farmRepo)
    automationService := service.NewAutomationService(farmRepo, sensorRepo, alertRepo)
    recommendationService := service.NewRecommendationService(aiService)

    // 8. HTTP Handlers
    farmHandler := http.NewFarmHandler(farmService, automationService, recommendationService)
    sensorHandler := http.NewSensorHandler(sensorRepo)
    alertHandler := http.NewAlertHandler(alertRepo)

    // 9. MQTT Handler
    mqttHandler := mqtt.NewHandler(sensorRepo, alertRepo, automationService)
    mqttClient.SetHandler(mqttHandler)
    go mqttClient.Connect(ctx)

    // 10. Setup Routes
    r := gin.Default()
    http.SetupRoutes(r, farmHandler, sensorHandler, alertHandler)
    
    // 11. WebSocket
    go wsServer.Run()

    // 12. Start Server
    r.Run(":8080")
}
```

---

## 5. Domain Layer

### 5.1 Farm Entity

```go
// domain/entity/farm.go
package entity

import (
    "time"
)

type Farm struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Location    string    `json:"location"`
    UserID      int64     `json:"user_id"`
    Status      string    `json:"status"` // active, inactive, maintenance
    Config      FarmConfig `json:"config"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type FarmConfig struct {
    WaterTankCapacity float64 `json:"water_tank_capacity"` // liters
    NutrientType      string  `json:"nutrient_type"`
    LightSchedule     Schedule `json:"light_schedule"`
    PumpSchedule      Schedule `json:"pump_schedule"`
    PHMin             float64 `json:"ph_min"`
    PHMax             float64 `json:"ph_max"`
    ECMin             float64 `json:"ec_min"`
    ECMax             float64 `json:"ec_max"`
    TemperatureMin    float64 `json:"temperature_min"`
    TemperatureMax    float64 `json:"temperature_max"`
}

type Schedule struct {
    StartTime string `json:"start_time"` // "06:00"
    EndTime   string `json:"end_time"`   // "18:00"
    Interval  int    `json:"interval"`   // minutes
    Duration  int    `json:"duration"`   // seconds
}

type Crop struct {
    ID          int64     `json:"id"`
    FarmID      int64     `json:"farm_id"`
    Name        string    `json:"name"`
    Type        string    `json:"type"` // lettuce, basil, mint, etc.
    PlantDate   time.Time `json:"plant_date"`
    HarvestDate time.Time `json:"harvest_date"`
    Status      string    `json:"status"` // growing, harvested, failed
    Quantity    int       `json:"quantity"`
    Notes       string    `json:"notes"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### 5.2 Sensor Entity

```go
// domain/entity/sensor.go
package entity

import "time"

type SensorData struct {
    ID          string    `json:"id"`
    FarmID      int64     `json:"farm_id"`
    DeviceID    string    `json:"device_id"`
    Temperature float64   `json:"temperature"`    // °C
    Humidity    float64   `json:"humidity"`       // %
    WaterTemp   float64   `json:"water_temp"`     // °C
    PH          float64   `json:"ph"`
    EC          float64   `json:"ec"`             // mS/cm
    TDS         float64   `json:"tds"`            // ppm
    Light       float64   `json:"light"`          // lux
    WaterLevel  float64   `json:"water_level"`    // cm
    CO2         float64   `json:"co2"`            // ppm
    Timestamp   time.Time `json:"timestamp"`
}

type SensorThreshold struct {
    ID          int64     `json:"id"`
    FarmID      int64     `json:"farm_id"`
    Parameter   string    `json:"parameter"` // temperature, ph, ec, etc.
    MinValue    float64   `json:"min_value"`
    MaxValue    float64   `json:"max_value"`
    AlertLevel  string    `json:"alert_level"` // info, warning, critical
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### 5.3 Actuator Entity

```go
// domain/entity/actuator.go
package entity

import "time"

type Actuator struct {
    ID          int64     `json:"id"`
    FarmID      int64     `json:"farm_id"`
    Name        string    `json:"name"`
    Type        string    `json:"type"` // pump, light, fan, valve
    Pin         int       `json:"pin"`
    Status      string    `json:"status"` // on, off, error
    LastCommand string    `json:"last_command"`
    LastUpdate  time.Time `json:"last_update"`
    CreatedAt   time.Time `json:"created_at"`
}

type ActuatorCommand struct {
    ActuatorID int64     `json:"actuator_id"`
    Command    string    `json:"command"` // on, off, toggle
    Duration   int       `json:"duration"` // seconds, 0 = indefinite
    Timestamp  time.Time `json:"timestamp"`
}

type ActuatorLog struct {
    ID          int64     `json:"id"`
    ActuatorID  int64     `json:"actuator_id"`
    Command     string    `json:"command"`
    Result      string    `json:"result"` // success, failed
    ErrorMsg    string    `json:"error_msg"`
    Timestamp   time.Time `json:"timestamp"`
}
```

### 5.4 Alert Entity

```go
// domain/entity/alert.go
package entity

import "time"

type Alert struct {
    ID          int64     `json:"id"`
    FarmID      int64     `json:"farm_id"`
    Type        string    `json:"type"` // sensor, actuator, system
    Severity    string    `json:"severity"` // info, warning, critical
    Title       string    `json:"title"`
    Message     string    `json:"message"`
    Parameter   string    `json:"parameter"`
    Value       float64   `json:"value"`
    Threshold   float64   `json:"threshold"`
    Status      string    `json:"status"` // active, acknowledged, resolved
    ResolvedAt  *time.Time `json:"resolved_at"`
    CreatedAt   time.Time `json:"created_at"`
}

type AlertRule struct {
    ID          int64     `json:"id"`
    FarmID      int64     `json:"farm_id"`
    Parameter   string    `json:"parameter"`
    Operator    string    `json:"operator"` // gt, lt, eq, between
    Value1      float64   `json:"value1"`
    Value2      float64   `json:"value2"` // for between
    Severity    string    `json:"severity"`
    Message     string    `json:"message"`
    IsActive    bool      `json:"is_active"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

---

## 6. Application Layer

### 6.1 Farm Use Cases

```go
// application/create_farm.go
package application

import (
    "context"
    "time"

    "github.com/yourproject/internal/modules/hydroponics/domain/entity"
    "github.com/yourproject/internal/modules/hydroponics/domain/repository"
)

type CreateFarmUseCase struct {
    repo repository.FarmRepository
}

func NewCreateFarmUseCase(repo repository.FarmRepository) *CreateFarmUseCase {
    return &CreateFarmUseCase{repo: repo}
}

type CreateFarmRequest struct {
    Name        string               `json:"name"`
    Description string               `json:"description"`
    Location    string               `json:"location"`
    UserID      int64                `json:"user_id"`
    Config      entity.FarmConfig    `json:"config"`
}

type CreateFarmResponse struct {
    ID          int64                `json:"id"`
    Name        string               `json:"name"`
    Status      string               `json:"status"`
    CreatedAt   time.Time            `json:"created_at"`
}

func (uc *CreateFarmUseCase) Execute(ctx context.Context, req CreateFarmRequest) (*CreateFarmResponse, error) {
    farm := &entity.Farm{
        Name:        req.Name,
        Description: req.Description,
        Location:    req.Location,
        UserID:      req.UserID,
        Status:      "active",
        Config:      req.Config,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    if err := uc.repo.Create(ctx, farm); err != nil {
        return nil, err
    }

    return &CreateFarmResponse{
        ID:        farm.ID,
        Name:      farm.Name,
        Status:    farm.Status,
        CreatedAt: farm.CreatedAt,
    }, nil
}
```

### 6.2 Sensor Use Cases

```go
// application/get_sensor_data.go
package application

import (
    "context"
    "time"

    "github.com/yourproject/internal/modules/hydroponics/domain/entity"
    "github.com/yourproject/internal/modules/hydroponics/domain/repository"
)

type GetSensorDataUseCase struct {
    repo repository.SensorRepository
}

func NewGetSensorDataUseCase(repo repository.SensorRepository) *GetSensorDataUseCase {
    return &GetSensorDataUseCase{repo: repo}
}

type GetSensorDataRequest struct {
    FarmID    int64     `json:"farm_id"`
    StartTime time.Time `json:"start_time"`
    EndTime   time.Time `json:"end_time"`
    Limit     int       `json:"limit"`
}

type GetSensorDataResponse struct {
    Data       []entity.SensorData `json:"data"`
    Total      int                 `json:"total"`
    Statistics SensorStatistics    `json:"statistics"`
}

type SensorStatistics struct {
    AvgTemp     float64 `json:"avg_temp"`
    MaxTemp     float64 `json:"max_temp"`
    MinTemp     float64 `json:"min_temp"`
    AvgPH       float64 `json:"avg_ph"`
    MaxPH       float64 `json:"max_ph"`
    MinPH       float64 `json:"min_ph"`
    AvgEC       float64 `json:"avg_ec"`
    MaxEC       float64 `json:"max_ec"`
    MinEC       float64 `json:"min_ec"`
    AvgHumidity float64 `json:"avg_humidity"`
    MaxHumidity float64 `json:"max_humidity"`
    MinHumidity float64 `json:"min_humidity"`
}

func (uc *GetSensorDataUseCase) Execute(ctx context.Context, req GetSensorDataRequest) (*GetSensorDataResponse, error) {
    data, err := uc.repo.GetByFarmAndTimeRange(ctx, req.FarmID, req.StartTime, req.EndTime, req.Limit)
    if err != nil {
        return nil, err
    }

    stats := calculateStatistics(data)

    return &GetSensorDataResponse{
        Data:       data,
        Total:      len(data),
        Statistics: stats,
    }, nil
}

func calculateStatistics(data []entity.SensorData) SensorStatistics {
    stats := SensorStatistics{}
    if len(data) == 0 {
        return stats
    }

    stats.MinTemp = data[0].Temperature
    stats.MaxTemp = data[0].Temperature
    stats.MinPH = data[0].PH
    stats.MaxPH = data[0].PH
    stats.MinEC = data[0].EC
    stats.MaxEC = data[0].EC
    stats.MinHumidity = data[0].Humidity
    stats.MaxHumidity = data[0].Humidity

    var sumTemp, sumPH, sumEC, sumHumidity float64
    for _, d := range data {
        sumTemp += d.Temperature
        sumPH += d.PH
        sumEC += d.EC
        sumHumidity += d.Humidity

        if d.Temperature < stats.MinTemp {
            stats.MinTemp = d.Temperature
        }
        if d.Temperature > stats.MaxTemp {
            stats.MaxTemp = d.Temperature
        }
        if d.PH < stats.MinPH {
            stats.MinPH = d.PH
        }
        if d.PH > stats.MaxPH {
            stats.MaxPH = d.PH
        }
        if d.EC < stats.MinEC {
            stats.MinEC = d.EC
        }
        if d.EC > stats.MaxEC {
            stats.MaxEC = d.EC
        }
        if d.Humidity < stats.MinHumidity {
            stats.MinHumidity = d.Humidity
        }
        if d.Humidity > stats.MaxHumidity {
            stats.MaxHumidity = d.Humidity
        }
    }

    n := float64(len(data))
    stats.AvgTemp = sumTemp / n
    stats.AvgPH = sumPH / n
    stats.AvgEC = sumEC / n
    stats.AvgHumidity = sumHumidity / n

    return stats
}
```

### 6.3 Automation Service

```go
// application/automation_service.go
package application

import (
    "context"
    "time"

    "github.com/yourproject/internal/modules/hydroponics/domain/entity"
    "github.com/yourproject/internal/modules/hydroponics/domain/repository"
)

type AutomationService struct {
    farmRepo   repository.FarmRepository
    sensorRepo repository.SensorRepository
    actuatorRepo repository.ActuatorRepository
    alertRepo  repository.AlertRepository
}

func NewAutomationService(
    farmRepo repository.FarmRepository,
    sensorRepo repository.SensorRepository,
    actuatorRepo repository.ActuatorRepository,
    alertRepo repository.AlertRepository,
) *AutomationService {
    return &AutomationService{
        farmRepo:   farmRepo,
        sensorRepo: sensorRepo,
        actuatorRepo: actuatorRepo,
        alertRepo:  alertRepo,
    }
}

func (s *AutomationService) ProcessAutomation(ctx context.Context, farmID int64) error {
    // Get latest sensor data
    latest, err := s.sensorRepo.GetLatest(ctx, farmID)
    if err != nil {
        return err
    }

    // Get farm config
    farm, err := s.farmRepo.GetByID(ctx, farmID)
    if err != nil {
        return err
    }

    // Check thresholds
    alerts := s.checkThresholds(ctx, farmID, latest, farm.Config)

    // Control actuators
    actions := s.controlActuators(ctx, farmID, latest, farm.Config)

    // Process schedules
    s.processSchedules(ctx, farmID)

    // Save alerts
    for _, alert := range alerts {
        s.alertRepo.Create(ctx, alert)
    }

    // Send commands for actions
    for _, action := range actions {
        s.actuatorRepo.SendCommand(ctx, action)
    }

    return nil
}

func (s *AutomationService) checkThresholds(ctx context.Context, farmID int64, data *entity.SensorData, config entity.FarmConfig) []*entity.Alert {
    var alerts []*entity.Alert

    // Check temperature
    if data.Temperature < config.TemperatureMin {
        alerts = append(alerts, &entity.Alert{
            FarmID:    farmID,
            Type:      "sensor",
            Severity:  "warning",
            Title:     "อุณหภูมิต่ำเกินไป",
            Message:   "อุณหภูมิ %.1f°C ต่ำกว่าค่าต่ำสุด %.1f°C",
            Parameter: "temperature",
            Value:     data.Temperature,
            Threshold: config.TemperatureMin,
            Status:    "active",
            CreatedAt: time.Now(),
        })
    } else if data.Temperature > config.TemperatureMax {
        alerts = append(alerts, &entity.Alert{
            FarmID:    farmID,
            Type:      "sensor",
            Severity:  "critical",
            Title:     "อุณหภูมิสูงเกินไป",
            Message:   "อุณหภูมิ %.1f°C สูงกว่าค่าสูงสุด %.1f°C",
            Parameter: "temperature",
            Value:     data.Temperature,
            Threshold: config.TemperatureMax,
            Status:    "active",
            CreatedAt: time.Now(),
        })
    }

    // Check pH
    if data.PH < config.PHMin {
        alerts = append(alerts, &entity.Alert{
            FarmID:    farmID,
            Type:      "sensor",
            Severity:  "warning",
            Title:     "pH ต่ำเกินไป",
            Message:   "pH %.2f ต่ำกว่าค่าต่ำสุด %.2f",
            Parameter: "ph",
            Value:     data.PH,
            Threshold: config.PHMin,
            Status:    "active",
            CreatedAt: time.Now(),
        })
    } else if data.PH > config.PHMax {
        alerts = append(alerts, &entity.Alert{
            FarmID:    farmID,
            Type:      "sensor",
            Severity:  "warning",
            Title:     "pH สูงเกินไป",
            Message:   "pH %.2f สูงกว่าค่าสูงสุด %.2f",
            Parameter: "ph",
            Value:     data.PH,
            Threshold: config.PHMax,
            Status:    "active",
            CreatedAt: time.Now(),
        })
    }

    // Check EC
    if data.EC < config.ECMin {
        alerts = append(alerts, &entity.Alert{
            FarmID:    farmID,
            Type:      "sensor",
            Severity:  "warning",
            Title:     "EC ต่ำเกินไป",
            Message:   "EC %.2f mS/cm ต่ำกว่าค่าต่ำสุด %.2f",
            Parameter: "ec",
            Value:     data.EC,
            Threshold: config.ECMin,
            Status:    "active",
            CreatedAt: time.Now(),
        })
    } else if data.EC > config.ECMax {
        alerts = append(alerts, &entity.Alert{
            FarmID:    farmID,
            Type:      "sensor",
            Severity:  "critical",
            Title:     "EC สูงเกินไป",
            Message:   "EC %.2f mS/cm สูงกว่าค่าสูงสุด %.2f",
            Parameter: "ec",
            Value:     data.EC,
            Threshold: config.ECMax,
            Status:    "active",
            CreatedAt: time.Now(),
        })
    }

    return alerts
}

func (s *AutomationService) controlActuators(ctx context.Context, farmID int64, data *entity.SensorData, config entity.FarmConfig) []*entity.ActuatorCommand {
    var commands []*entity.ActuatorCommand

    // pH control
    if data.PH < config.PHMin {
        // Add pH Up (valve 1)
        commands = append(commands, &entity.ActuatorCommand{
            ActuatorID: 1, // pH Up valve
            Command:    "on",
            Duration:   30,
            Timestamp:  time.Now(),
        })
    } else if data.PH > config.PHMax {
        // Add pH Down (valve 2)
        commands = append(commands, &entity.ActuatorCommand{
            ActuatorID: 2, // pH Down valve
            Command:    "on",
            Duration:   30,
            Timestamp:  time.Now(),
        })
    }

    // EC control
    if data.EC < config.ECMin {
        // Add nutrient (valve 3)
        commands = append(commands, &entity.ActuatorCommand{
            ActuatorID: 3, // Nutrient valve
            Command:    "on",
            Duration:   60,
            Timestamp:  time.Now(),
        })
    }

    // Temperature control
    if data.Temperature > config.TemperatureMax {
        commands = append(commands, &entity.ActuatorCommand{
            ActuatorID: 4, // Fan
            Command:    "on",
            Duration:   300,
            Timestamp:  time.Now(),
        })
    }

    // Water level control
    if data.WaterLevel < 5.0 {
        commands = append(commands, &entity.ActuatorCommand{
            ActuatorID: 5, // Water pump
            Command:    "on",
            Duration:   120,
            Timestamp:  time.Now(),
        })
    }

    return commands
}

func (s *AutomationService) processSchedules(ctx context.Context, farmID int64) {
    now := time.Now()
    currentTime := now.Format("15:04")

    // Get farm for schedule
    farm, err := s.farmRepo.GetByID(ctx, farmID)
    if err != nil {
        return
    }

    // Check light schedule
    if currentTime == farm.Config.LightSchedule.StartTime {
        s.actuatorRepo.SendCommand(ctx, &entity.ActuatorCommand{
            ActuatorID: 6, // Light
            Command:    "on",
            Duration:   0,
            Timestamp:  time.Now(),
        })
    } else if currentTime == farm.Config.LightSchedule.EndTime {
        s.actuatorRepo.SendCommand(ctx, &entity.ActuatorCommand{
            ActuatorID: 6, // Light
            Command:    "off",
            Duration:   0,
            Timestamp:  time.Now(),
        })
    }

    // Check pump schedule
    if now.Hour()%farm.Config.PumpSchedule.Interval == 0 && now.Minute() == 0 {
        s.actuatorRepo.SendCommand(ctx, &entity.ActuatorCommand{
            ActuatorID: 5, // Water pump
            Command:    "on",
            Duration:   farm.Config.PumpSchedule.Duration,
            Timestamp:  time.Now(),
        })
    }
}
```

---

## 7. Infrastructure Layer

### 7.1 MQTT Handler

```go
// infrastructure/mqtt/handler.go
package mqtt

import (
    "context"
    "encoding/json"
    "time"

    "github.com/yourproject/internal/modules/hydroponics/domain/entity"
    "github.com/yourproject/internal/modules/hydroponics/domain/repository"
)

type MQTTMessage struct {
    Device    string                 `json:"device"`
    Data      map[string]interface{} `json:"data"`
    Timestamp int64                  `json:"timestamp"`
}

type MQTTHandler struct {
    sensorRepo   repository.SensorRepository
    alertRepo    repository.AlertRepository
    automationSvc *AutomationService
}

func NewMQTTHandler(
    sensorRepo repository.SensorRepository,
    alertRepo repository.AlertRepository,
    automationSvc *AutomationService,
) *MQTTHandler {
    return &MQTTHandler{
        sensorRepo:   sensorRepo,
        alertRepo:    alertRepo,
        automationSvc: automationSvc,
    }
}

func (h *MQTTHandler) HandleMessage(topic string, payload []byte) {
    var msg MQTTMessage
    if err := json.Unmarshal(payload, &msg); err != nil {
        return
    }

    // Parse farm ID from topic
    farmID := extractFarmID(topic)

    // Convert to sensor data
    sensorData := &entity.SensorData{
        FarmID:      farmID,
        DeviceID:    msg.Device,
        Temperature: getFloat(msg.Data, "temperature"),
        Humidity:    getFloat(msg.Data, "humidity"),
        WaterTemp:   getFloat(msg.Data, "water_temp"),
        PH:          getFloat(msg.Data, "ph"),
        EC:          getFloat(msg.Data, "ec"),
        TDS:         getFloat(msg.Data, "tds"),
        Light:       getFloat(msg.Data, "light"),
        WaterLevel:  getFloat(msg.Data, "water_level"),
        CO2:         getFloat(msg.Data, "co2"),
        Timestamp:   time.Unix(msg.Timestamp, 0),
    }

    // Save to database
    ctx := context.Background()
    if err := h.sensorRepo.Save(ctx, sensorData); err != nil {
        return
    }

    // Process automation
    h.automationSvc.ProcessAutomation(ctx, farmID)
}

func getFloat(data map[string]interface{}, key string) float64 {
    if val, ok := data[key]; ok {
        if f, ok := val.(float64); ok {
            return f
        }
        if f, ok := val.(int); ok {
            return float64(f)
        }
    }
    return 0.0
}

func extractFarmID(topic string) int64 {
    // topic format: hydroponics/farm01/data
    // Extract "farm01" and convert to ID
    return 1 // Simplified
}
```

### 7.2 Postgres Repository

```go
// infrastructure/persistence/postgres/farm_repo.go
package postgres

import (
    "context"
    "encoding/json"
    "time"

    "gorm.io/gorm"
    "github.com/yourproject/internal/modules/hydroponics/domain/entity"
    "github.com/yourproject/internal/modules/hydroponics/domain/repository"
)

type FarmModel struct {
    ID          int64     `gorm:"primaryKey"`
    Name        string    `gorm:"type:varchar(255);not null"`
    Description string    `gorm:"type:text"`
    Location    string    `gorm:"type:varchar(255)"`
    UserID      int64     `gorm:"not null;index"`
    Status      string    `gorm:"type:varchar(50);default:'active'"`
    Config      string    `gorm:"type:jsonb"`
    CreatedAt   time.Time `gorm:"autoCreateTime"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (FarmModel) TableName() string {
    return "farms"
}

type FarmRepository struct {
    db *gorm.DB
}

func NewFarmRepository(db *gorm.DB) repository.FarmRepository {
    return &FarmRepository{db: db}
}

func (r *FarmRepository) Create(ctx context.Context, farm *entity.Farm) error {
    configJSON, err := json.Marshal(farm.Config)
    if err != nil {
        return err
    }

    model := &FarmModel{
        Name:        farm.Name,
        Description: farm.Description,
        Location:    farm.Location,
        UserID:      farm.UserID,
        Status:      farm.Status,
        Config:      string(configJSON),
    }

    if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
        return err
    }

    farm.ID = model.ID
    farm.CreatedAt = model.CreatedAt
    farm.UpdatedAt = model.UpdatedAt

    return nil
}

func (r *FarmRepository) GetByID(ctx context.Context, id int64) (*entity.Farm, error) {
    var model FarmModel
    if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
        return nil, err
    }

    return r.toEntity(&model), nil
}

func (r *FarmRepository) GetByUserID(ctx context.Context, userID int64) ([]*entity.Farm, error) {
    var models []FarmModel
    if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models).Error; err != nil {
        return nil, err
    }

    farms := make([]*entity.Farm, len(models))
    for i, model := range models {
        farms[i] = r.toEntity(&model)
    }

    return farms, nil
}

func (r *FarmRepository) Update(ctx context.Context, farm *entity.Farm) error {
    configJSON, err := json.Marshal(farm.Config)
    if err != nil {
        return err
    }

    model := &FarmModel{
        ID:          farm.ID,
        Name:        farm.Name,
        Description: farm.Description,
        Location:    farm.Location,
        UserID:      farm.UserID,
        Status:      farm.Status,
        Config:      string(configJSON),
        UpdatedAt:   time.Now(),
    }

    return r.db.WithContext(ctx).Save(model).Error
}

func (r *FarmRepository) Delete(ctx context.Context, id int64) error {
    return r.db.WithContext(ctx).Delete(&FarmModel{}, id).Error
}

func (r *FarmRepository) toEntity(model *FarmModel) *entity.Farm {
    var config entity.FarmConfig
    json.Unmarshal([]byte(model.Config), &config)

    return &entity.Farm{
        ID:          model.ID,
        Name:        model.Name,
        Description: model.Description,
        Location:    model.Location,
        UserID:      model.UserID,
        Status:      model.Status,
        Config:      config,
        CreatedAt:   model.CreatedAt,
        UpdatedAt:   model.UpdatedAt,
    }
}
```

### 7.3 InfluxDB Repository

```go
// infrastructure/persistence/influxdb/sensor_repo.go
package influxdb

import (
    "context"
    "fmt"
    "time"

    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"
    "github.com/yourproject/internal/modules/hydroponics/domain/entity"
    "github.com/yourproject/internal/modules/hydroponics/domain/repository"
)

type SensorRepository struct {
    client   influxdb2.Client
    bucket   string
    org      string
    writeAPI api.WriteAPIBlocking
    queryAPI api.QueryAPI
}

func NewSensorRepository(config Config) *SensorRepository {
    client := influxdb2.NewClient(config.URL, config.Token)

    return &SensorRepository{
        client:   client,
        bucket:   config.Bucket,
        org:      config.Org,
        writeAPI: client.WriteAPIBlocking(config.Org, config.Bucket),
        queryAPI: client.QueryAPI(config.Org),
    }
}

func (r *SensorRepository) Save(ctx context.Context, data *entity.SensorData) error {
    point := influxdb2.NewPoint(
        "sensor_data",
        map[string]string{
            "farm_id":   fmt.Sprintf("%d", data.FarmID),
            "device_id": data.DeviceID,
        },
        map[string]interface{}{
            "temperature": data.Temperature,
            "humidity":    data.Humidity,
            "water_temp":  data.WaterTemp,
            "ph":          data.PH,
            "ec":          data.EC,
            "tds":         data.TDS,
            "light":       data.Light,
            "water_level": data.WaterLevel,
            "co2":         data.CO2,
        },
        data.Timestamp,
    )

    return r.writeAPI.WritePoint(ctx, point)
}

func (r *SensorRepository) GetLatest(ctx context.Context, farmID int64) (*entity.SensorData, error) {
    query := fmt.Sprintf(`
        from(bucket: "%s")
        |> range(start: -1h)
        |> filter(fn: (r) => r._measurement == "sensor_data" and r.farm_id == "%d")
        |> last()
    `, r.bucket, farmID)

    result, err := r.queryAPI.Query(ctx, query)
    if err != nil {
        return nil, err
    }

    data := &entity.SensorData{}
    for result.Next() {
        record := result.Record()
        data.Timestamp = record.Time()
        data.FarmID = farmID

        switch record.Field() {
        case "temperature":
            data.Temperature = record.Value().(float64)
        case "humidity":
            data.Humidity = record.Value().(float64)
        case "water_temp":
            data.WaterTemp = record.Value().(float64)
        case "ph":
            data.PH = record.Value().(float64)
        case "ec":
            data.EC = record.Value().(float64)
        case "tds":
            data.TDS = record.Value().(float64)
        case "light":
            data.Light = record.Value().(float64)
        case "water_level":
            data.WaterLevel = record.Value().(float64)
        case "co2":
            data.CO2 = record.Value().(float64)
        }
    }

    return data, nil
}

func (r *SensorRepository) GetByFarmAndTimeRange(
    ctx context.Context,
    farmID int64,
    start, end time.Time,
    limit int,
) ([]entity.SensorData, error) {
    query := fmt.Sprintf(`
        from(bucket: "%s")
        |> range(start: %d, stop: %d)
        |> filter(fn: (r) => r._measurement == "sensor_data" and r.farm_id == "%d")
        |> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")
        |> limit(n: %d)
    `, r.bucket, start.Unix(), end.Unix(), farmID, limit)

    result, err := r.queryAPI.Query(ctx, query)
    if err != nil {
        return nil, err
    }

    var data []entity.SensorData
    for result.Next() {
        record := result.Record()
        data = append(data, entity.SensorData{
            Timestamp:   record.Time(),
            FarmID:      farmID,
            Temperature: r.getFloat(record, "temperature"),
            Humidity:    r.getFloat(record, "humidity"),
            WaterTemp:   r.getFloat(record, "water_temp"),
            PH:          r.getFloat(record, "ph"),
            EC:          r.getFloat(record, "ec"),
            TDS:         r.getFloat(record, "tds"),
            Light:       r.getFloat(record, "light"),
            WaterLevel:  r.getFloat(record, "water_level"),
            CO2:         r.getFloat(record, "co2"),
        })
    }

    return data, nil
}

func (r *SensorRepository) getFloat(record *api.Record, field string) float64 {
    if val, ok := record.ValueByKey(field); ok {
        if f, ok := val.(float64); ok {
            return f
        }
    }
    return 0.0
}
```

---

## 8. API Endpoints

### 8.1 Routes

```go
// interfaces/http/routes.go
package http

import (
    "github.com/gin-gonic/gin"
    "github.com/yourproject/internal/modules/hydroponics/interfaces/middleware"
)

func SetupRoutes(
    r *gin.Engine,
    farmHandler *FarmHandler,
    sensorHandler *SensorHandler,
    alertHandler *AlertHandler,
    auth *middleware.AuthMiddleware,
) {
    api := r.Group("/api/v1/hydroponics")
    api.Use(auth.Authenticate())

    // Farm routes
    farms := api.Group("/farms")
    {
        farms.POST("/", farmHandler.CreateFarm)
        farms.GET("/", farmHandler.GetFarms)
        farms.GET("/:id", farmHandler.GetFarm)
        farms.PUT("/:id", farmHandler.UpdateFarm)
        farms.DELETE("/:id", farmHandler.DeleteFarm)
        farms.POST("/:id/activate", farmHandler.ActivateFarm)
        farms.POST("/:id/deactivate", farmHandler.DeactivateFarm)
    }

    // Sensor routes
    sensors := api.Group("/sensors")
    {
        sensors.GET("/:farm_id/latest", sensorHandler.GetLatest)
        sensors.GET("/:farm_id/history", sensorHandler.GetHistory)
        sensors.GET("/:farm_id/statistics", sensorHandler.GetStatistics)
        sensors.POST("/thresholds", sensorHandler.CreateThreshold)
        sensors.GET("/thresholds/:farm_id", sensorHandler.GetThresholds)
        sensors.PUT("/thresholds/:id", sensorHandler.UpdateThreshold)
    }

    // Actuator routes
    actuators := api.Group("/actuators")
    {
        actuators.GET("/:farm_id", actuatorHandler.GetActuators)
        actuators.POST("/:id/command", actuatorHandler.SendCommand)
        actuators.GET("/:id/logs", actuatorHandler.GetLogs)
    }

    // Schedule routes
    schedules := api.Group("/schedules")
    {
        schedules.POST("/", scheduleHandler.CreateSchedule)
        schedules.GET("/:farm_id", scheduleHandler.GetSchedules)
        schedules.PUT("/:id", scheduleHandler.UpdateSchedule)
        schedules.DELETE("/:id", scheduleHandler.DeleteSchedule)
    }

    // Alert routes
    alerts := api.Group("/alerts")
    {
        alerts.GET("/:farm_id", alertHandler.GetAlerts)
        alerts.PUT("/:id/acknowledge", alertHandler.AcknowledgeAlert)
        alerts.PUT("/:id/resolve", alertHandler.ResolveAlert)
        alerts.POST("/rules", alertHandler.CreateRule)
        alerts.GET("/rules/:farm_id", alertHandler.GetRules)
        alerts.PUT("/rules/:id", alertHandler.UpdateRule)
        alerts.DELETE("/rules/:id", alertHandler.DeleteRule)
    }

    // Crop routes
    crops := api.Group("/crops")
    {
        crops.POST("/", cropHandler.CreateCrop)
        crops.GET("/:farm_id", cropHandler.GetCrops)
        crops.GET("/:id", cropHandler.GetCrop)
        crops.PUT("/:id", cropHandler.UpdateCrop)
        crops.DELETE("/:id", cropHandler.DeleteCrop)
        crops.POST("/:id/harvest", cropHandler.HarvestCrop)
    }

    // Recommendation routes
    recs := api.Group("/recommendations")
    {
        recs.GET("/:farm_id", recHandler.GetRecommendations)
        recs.GET("/:farm_id/crop/:crop_id", recHandler.GetCropRecommendations)
    }

    // WebSocket
    r.GET("/ws", websocket.Handler)
}
```

### 8.2 API Documentation

| Method | Endpoint | Description |
|--------|----------|-------------|
| **Farm** | | |
| POST | `/api/v1/hydroponics/farms` | สร้างฟาร์มใหม่ |
| GET | `/api/v1/hydroponics/farms` | ดูรายการฟาร์ม |
| GET | `/api/v1/hydroponics/farms/:id` | ดูข้อมูลฟาร์ม |
| PUT | `/api/v1/hydroponics/farms/:id` | แก้ไขฟาร์ม |
| DELETE | `/api/v1/hydroponics/farms/:id` | ลบฟาร์ม |
| POST | `/api/v1/hydroponics/farms/:id/activate` | เปิดใช้งานฟาร์ม |
| POST | `/api/v1/hydroponics/farms/:id/deactivate` | ปิดใช้งานฟาร์ม |
| **Sensor** | | |
| GET | `/api/v1/hydroponics/sensors/:farm_id/latest` | ดูข้อมูลล่าสุด |
| GET | `/api/v1/hydroponics/sensors/:farm_id/history` | ดูข้อมูลย้อนหลัง |
| GET | `/api/v1/hydroponics/sensors/:farm_id/statistics` | ดูค่าสถิติ |
| POST | `/api/v1/hydroponics/sensors/thresholds` | ตั้งค่าเกณฑ์ |
| GET | `/api/v1/hydroponics/sensors/thresholds/:farm_id` | ดูเกณฑ์ที่ตั้ง |
| **Actuator** | | |
| GET | `/api/v1/hydroponics/actuators/:farm_id` | ดูอุปกรณ์ทั้งหมด |
| POST | `/api/v1/hydroponics/actuators/:id/command` | สั่งงานอุปกรณ์ |
| **Schedule** | | |
| POST | `/api/v1/hydroponics/schedules` | สร้างตารางเวลา |
| GET | `/api/v1/hydroponics/schedules/:farm_id` | ดูตารางเวลา |
| **Alert** | | |
| GET | `/api/v1/hydroponics/alerts/:farm_id` | ดูการแจ้งเตือน |
| PUT | `/api/v1/hydroponics/alerts/:id/acknowledge` | รับทราบการแจ้งเตือน |
| POST | `/api/v1/hydroponics/alerts/rules` | สร้างกฎการแจ้งเตือน |
| **Crop** | | |
| POST | `/api/v1/hydroponics/crops` | ปลูกพืชใหม่ |
| GET | `/api/v1/hydroponics/crops/:farm_id` | ดูพืชทั้งหมด |
| POST | `/api/v1/hydroponics/crops/:id/harvest` | เก็บเกี่ยว |
| **Recommendation** | | |
| GET | `/api/v1/hydroponics/recommendations/:farm_id` | ดูคำแนะนำ |

---

## 9. Mobile App Integration

### 9.1 Flutter App Structure

```dart
// lib/main.dart
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:hydroponics_app/providers/farm_provider.dart';
import 'package:hydroponics_app/providers/sensor_provider.dart';
import 'package:hydroponics_app/screens/dashboard_screen.dart';
import 'package:hydroponics_app/screens/farm_detail_screen.dart';
import 'package:hydroponics_app/screens/settings_screen.dart';
import 'package:hydroponics_app/services/api_service.dart';
import 'package:hydroponics_app/services/websocket_service.dart';

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => FarmProvider()),
        ChangeNotifierProvider(create: (_) => SensorProvider()),
        Provider(create: (_) => ApiService()),
        Provider(create: (_) => WebSocketService()),
      ],
      child: MaterialApp(
        title: 'Hydroponics Smart Farm',
        theme: ThemeData(
          primaryColor: Color(0xFF2E7D32),
          accentColor: Color(0xFF4CAF50),
          fontFamily: 'Sarabun',
        ),
        home: DashboardScreen(),
        routes: {
          '/farm': (context) => FarmDetailScreen(),
          '/settings': (context) => SettingsScreen(),
        },
      ),
    );
  }
}
```

### 9.2 Dashboard Screen

```dart
// lib/screens/dashboard_screen.dart
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:hydroponics_app/providers/farm_provider.dart';
import 'package:hydroponics_app/providers/sensor_provider.dart';
import 'package:hydroponics_app/widgets/sensor_card.dart';
import 'package:hydroponics_app/widgets/alert_card.dart';
import 'package:hydroponics_app/widgets/farm_card.dart';

class DashboardScreen extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    final farmProvider = Provider.of<FarmProvider>(context);
    final sensorProvider = Provider.of<SensorProvider>(context);

    return Scaffold(
      appBar: AppBar(
        title: Text('Smart Hydroponics'),
        actions: [
          IconButton(
            icon: Icon(Icons.notifications),
            onPressed: () {
              // Navigate to alerts
            },
          ),
          IconButton(
            icon: Icon(Icons.settings),
            onPressed: () {
              Navigator.pushNamed(context, '/settings');
            },
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          await farmProvider.loadFarms();
          await sensorProvider.loadLatestData();
        },
        child: ListView(
          padding: EdgeInsets.all(16),
          children: [
            // Farm selector
            Container(
              height: 120,
              child: ListView.builder(
                scrollDirection: Axis.horizontal,
                itemCount: farmProvider.farms.length,
                itemBuilder: (context, index) {
                  return FarmCard(farm: farmProvider.farms[index]);
                },
              ),
            ),
            SizedBox(height: 16),

            // Latest sensor data
            if (sensorProvider.latestData != null)
              SensorCard(data: sensorProvider.latestData!),

            SizedBox(height: 16),

            // Alerts
            AlertCard(alerts: sensorProvider.alerts),

            SizedBox(height: 16),

            // Quick actions
            GridView.count(
              shrinkWrap: true,
              physics: NeverScrollableScrollPhysics(),
              crossAxisCount: 4,
              childAspectRatio: 1.0,
              children: [
                _buildQuickAction(
                  Icons.lightbulb,
                  'ไฟ',
                  () => _controlActuator(context, 'light'),
                ),
                _buildQuickAction(
                  Icons.water_drop,
                  'ปั๊มน้ำ',
                  () => _controlActuator(context, 'pump'),
                ),
                _buildQuickAction(
                  Icons.toys,
                  'พัดลม',
                  () => _controlActuator(context, 'fan'),
                ),
                _buildQuickAction(
                  Icons.grass,
                  'สารอาหาร',
                  () => _controlActuator(context, 'nutrient'),
                ),
              ],
            ),
          ],
        ),
      ),
      bottomNavigationBar: BottomNavigationBar(
        items: [
          BottomNavigationBarItem(
            icon: Icon(Icons.dashboard),
            label: 'หน้าหลัก',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.show_chart),
            label: 'กราฟ',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.history),
            label: 'ประวัติ',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.person),
            label: 'โปรไฟล์',
          ),
        ],
      ),
    );
  }

  Widget _buildQuickAction(IconData icon, String label, VoidCallback onTap) {
    return GestureDetector(
      onTap: onTap,
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            padding: EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: Colors.green[50],
              borderRadius: BorderRadius.circular(12),
            ),
            child: Icon(icon, color: Colors.green[700], size: 30),
          ),
          SizedBox(height: 4),
          Text(
            label,
            style: TextStyle(fontSize: 12),
          ),
        ],
      ),
    );
  }

  void _controlActuator(BuildContext context, String actuator) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('ควบคุม $actuator'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ElevatedButton(
              onPressed: () {
                // Send command to actuator
                Navigator.pop(context);
              },
              child: Text('เปิด'),
            ),
            SizedBox(height: 8),
            ElevatedButton(
              onPressed: () {
                // Send command to actuator
                Navigator.pop(context);
              },
              child: Text('ปิด'),
              style: ElevatedButton.styleFrom(
                primary: Colors.red,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
```

### 9.3 Sensor Card Widget

```dart
// lib/widgets/sensor_card.dart
import 'package:flutter/material.dart';
import 'package:hydroponics_app/models/sensor_data.dart';

class SensorCard extends StatelessWidget {
  final SensorData data;

  const SensorCard({Key? key, required this.data}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 2,
      child: Padding(
        padding: EdgeInsets.all(16),
        child: Column(
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                _buildSensorItem('🌡️', 'อุณหภูมิ', data.temperature, '°C'),
                _buildSensorItem('💧', 'ความชื้น', data.humidity, '%'),
                _buildSensorItem('📊', 'pH', data.ph, ''),
              ],
            ),
            SizedBox(height: 12),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                _buildSensorItem('⚡', 'EC', data.ec, 'mS/cm'),
                _buildSensorItem('💦', 'น้ำ', data.waterLevel, 'cm'),
                _buildSensorItem('☀️', 'แสง', data.light, 'lux'),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSensorItem(String icon, String label, double value, String unit) {
    return Column(
      children: [
        Text(icon, style: TextStyle(fontSize: 24)),
        Text(
          value.toStringAsFixed(1),
          style: TextStyle(
            fontSize: 20,
            fontWeight: FontWeight.bold,
            color: _getColor(label, value),
          ),
        ),
        Text(
          '$label $unit',
          style: TextStyle(fontSize: 12, color: Colors.grey),
        ),
      ],
    );
  }

  Color _getColor(String label, double value) {
    switch (label) {
      case 'อุณหภูมิ':
        if (value > 30) return Colors.red;
        if (value < 20) return Colors.blue;
        return Colors.green;
      case 'pH':
        if (value > 7.0 || value < 5.5) return Colors.orange;
        return Colors.green;
      case 'EC':
        if (value > 2.5 || value < 0.8) return Colors.orange;
        return Colors.green;
      default:
        return Colors.black;
    }
  }
}
```

---

## 10. Dashboard UI

### 10.1 Web Dashboard (React)

```tsx
// src/pages/Dashboard.tsx
import React, { useEffect, useState } from 'react';
import {
  Grid,
  Paper,
  Typography,
  Box,
  Card,
  CardContent,
  IconButton,
  LinearProgress,
} from '@mui/material';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import {
  WbSunny,
  WaterDrop,
  Thermostat,
  Opacity,
  Grass,
  Notifications,
  Settings,
} from '@mui/icons-material';
import { useWebSocket } from '../hooks/useWebSocket';
import { SensorData, FarmData } from '../types';

const Dashboard: React.FC = () => {
  const [farmData, setFarmData] = useState<FarmData | null>(null);
  const [sensorHistory, setSensorHistory] = useState<SensorData[]>([]);
  const [alerts, setAlerts] = useState<any[]>([]);
  const { lastMessage, sendMessage } = useWebSocket('/ws');

  useEffect(() => {
    // Fetch initial data
    fetchFarmData();
    fetchSensorHistory();
    fetchAlerts();

    // Subscribe to real-time updates
    sendMessage(JSON.stringify({
      type: 'subscribe',
      farmId: 'farm01',
    }));
  }, []);

  useEffect(() => {
    if (lastMessage) {
      const data = JSON.parse(lastMessage);
      if (data.type === 'sensor_update') {
        updateSensorData(data.payload);
      } else if (data.type === 'alert') {
        setAlerts((prev) => [data.payload, ...prev].slice(0, 10));
      }
    }
  }, [lastMessage]);

  const updateSensorData = (data: SensorData) => {
    setFarmData((prev) => ({
      ...prev!,
      latestData: data,
    }));
    setSensorHistory((prev) => [...prev, data].slice(-60));
  };

  const fetchFarmData = async () => {
    const response = await fetch('/api/v1/hydroponics/farms/1');
    const data = await response.json();
    setFarmData(data);
  };

  const fetchSensorHistory = async () => {
    const response = await fetch(
      '/api/v1/hydroponics/sensors/1/history?limit=60'
    );
    const data = await response.json();
    setSensorHistory(data.data);
  };

  const fetchAlerts = async () => {
    const response = await fetch('/api/v1/hydroponics/alerts/1');
    const data = await response.json();
    setAlerts(data);
  };

  if (!farmData) {
    return <LinearProgress />;
  }

  const latest = farmData.latestData;

  return (
    <Box sx={{ flexGrow: 1, p: 3 }}>
      {/* Header */}
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">
          🌱 {farmData.name}
        </Typography>
        <Box>
          <IconButton color="primary">
            <Notifications />
          </IconButton>
          <IconButton>
            <Settings />
          </IconButton>
        </Box>
      </Box>

      {/* Sensor Cards */}
      <Grid container spacing={3} mb={3}>
        <Grid item xs={12} sm={6} md={4} lg={2}>
          <SensorCard
            icon={<Thermostat />}
            title="อุณหภูมิ"
            value={latest?.temperature || 0}
            unit="°C"
            color={getTempColor(latest?.temperature)}
          />
        </Grid>
        <Grid item xs={12} sm={6} md={4} lg={2}>
          <SensorCard
            icon={<Opacity />}
            title="ความชื้น"
            value={latest?.humidity || 0}
            unit="%"
            color={getHumidityColor(latest?.humidity)}
          />
        </Grid>
        <Grid item xs={12} sm={6} md={4} lg={2}>
          <SensorCard
            icon={<WaterDrop />}
            title="pH"
            value={latest?.ph || 0}
            unit=""
            color={getPHColor(latest?.ph)}
          />
        </Grid>
        <Grid item xs={12} sm={6} md={4} lg={2}>
          <SensorCard
            icon={<Grass />}
            title="EC"
            value={latest?.ec || 0}
            unit="mS/cm"
            color={getECColor(latest?.ec)}
          />
        </Grid>
        <Grid item xs={12} sm={6} md={4} lg={2}>
          <SensorCard
            icon={<WbSunny />}
            title="แสง"
            value={latest?.light || 0}
            unit="lux"
            color="#FFA000"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={4} lg={2}>
          <SensorCard
            icon={<WaterDrop />}
            title="ระดับน้ำ"
            value={latest?.waterLevel || 0}
            unit="cm"
            color="#1976D2"
          />
        </Grid>
      </Grid>

      {/* Charts */}
      <Grid container spacing={3} mb={3}>
        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 2 }}>
            <Typography variant="h6" mb={2}>
              อุณหภูมิและความชื้น
            </Typography>
            <ResponsiveContainer width="100%" height={250}>
              <LineChart data={sensorHistory}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="timestamp" />
                <YAxis yAxisId="left" />
                <YAxis yAxisId="right" orientation="right" />
                <Tooltip />
                <Legend />
                <Line
                  yAxisId="left"
                  type="monotone"
                  dataKey="temperature"
                  stroke="#FF6B6B"
                  name="อุณหภูมิ °C"
                />
                <Line
                  yAxisId="right"
                  type="monotone"
                  dataKey="humidity"
                  stroke="#4ECDC4"
                  name="ความชื้น %"
                />
              </LineChart>
            </ResponsiveContainer>
          </Paper>
        </Grid>
        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 2 }}>
            <Typography variant="h6" mb={2}>
              pH และ EC
            </Typography>
            <ResponsiveContainer width="100%" height={250}>
              <LineChart data={sensorHistory}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="timestamp" />
                <YAxis yAxisId="left" />
                <YAxis yAxisId="right" orientation="right" />
                <Tooltip />
                <Legend />
                <Line
                  yAxisId="left"
                  type="monotone"
                  dataKey="ph"
                  stroke="#FFD93D"
                  name="pH"
                />
                <Line
                  yAxisId="right"
                  type="monotone"
                  dataKey="ec"
                  stroke="#6BCB77"
                  name="EC mS/cm"
                />
              </LineChart>
            </ResponsiveContainer>
          </Paper>
        </Grid>
      </Grid>

      {/* Alerts */}
      {alerts.length > 0 && (
        <Paper sx={{ p: 2 }}>
          <Typography variant="h6" mb={2}>
            ⚠️ การแจ้งเตือนล่าสุด
          </Typography>
          {alerts.map((alert, index) => (
            <AlertItem key={index} alert={alert} />
          ))}
        </Paper>
      )}
    </Box>
  );
};

// Sensor Card Component
const SensorCard: React.FC<{
  icon: React.ReactNode;
  title: string;
  value: number;
  unit: string;
  color: string;
}> = ({ icon, title, value, unit, color }) => (
  <Card>
    <CardContent>
      <Box display="flex" alignItems="center" gap={1}>
        <Box color={color}>{icon}</Box>
        <Typography variant="body2" color="textSecondary">
          {title}
        </Typography>
      </Box>
      <Typography variant="h4" color={color}>
        {value.toFixed(1)}
      </Typography>
      <Typography variant="caption" color="textSecondary">
        {unit}
      </Typography>
    </CardContent>
  </Card>
);

// Alert Item
const AlertItem: React.FC<{ alert: any }> = ({ alert }) => (
  <Box
    sx={{
      p: 1,
      mb: 1,
      borderLeft: `4px solid ${getSeverityColor(alert.severity)}`,
      bgcolor: getSeverityBg(alert.severity),
      borderRadius: 1,
    }}
  >
    <Typography variant="body1">
      {alert.title}
    </Typography>
    <Typography variant="body2" color="textSecondary">
      {alert.message}
    </Typography>
    <Typography variant="caption" color="textSecondary">
      {new Date(alert.createdAt).toLocaleString()}
    </Typography>
  </Box>
);

// Helper functions
const getTempColor = (temp?: number) => {
  if (!temp) return '#000';
  if (temp > 30) return '#D32F2F';
  if (temp < 20) return '#1976D2';
  return '#388E3C';
};

const getHumidityColor = (humidity?: number) => {
  if (!humidity) return '#000';
  if (humidity > 80) return '#D32F2F';
  if (humidity < 40) return '#F57C00';
  return '#388E3C';
};

const getPHColor = (ph?: number) => {
  if (!ph) return '#000';
  if (ph > 7.0 || ph < 5.5) return '#F57C00';
  return '#388E3C';
};

const getECColor = (ec?: number) => {
  if (!ec) return '#000';
  if (ec > 2.5 || ec < 0.8) return '#F57C00';
  return '#388E3C';
};

const getSeverityColor = (severity: string) => {
  switch (severity) {
    case 'critical': return '#D32F2F';
    case 'warning': return '#F57C00';
    case 'info': return '#1976D2';
    default: return '#757575';
  }
};

const getSeverityBg = (severity: string) => {
  switch (severity) {
    case 'critical': return '#FFEBEE';
    case 'warning': return '#FFF3E0';
    case 'info': return '#E3F2FD';
    default: return '#F5F5F5';
  }
};

export default Dashboard;
```

---

## 11. การติดตั้งและใช้งาน

### 11.1 System Requirements

| Component | Requirement |
|-----------|-------------|
| **Backend** | Go 1.21+, PostgreSQL 15+, InfluxDB 2.x |
| **Frontend** | Node.js 18+, React 18 / Angular 17 |
| **Mobile** | Flutter 3.16+ |
| **Hardware** | ESP32, Sensors, Actuators |
| **Cloud** | Ubuntu 20.04+, Docker 24+, Kubernetes 1.28+ |

### 11.2 Quick Start

```bash
# 1. Clone repository
git clone https://github.com/yourproject/hydroponics-smart-farm
cd hydroponics-smart-farm

# 2. Setup database
docker-compose up -d postgres influxdb

# 3. Run migrations
make migrate-up

# 4. Start backend
go run cmd/hydroponics/main.go

# 5. Start frontend
cd frontend
npm install
npm start

# 6. Start mobile app
cd mobile
flutter pub get
flutter run
```

### 11.3 Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: hydroponics
      POSTGRES_PASSWORD: password
      POSTGRES_DB: hydroponics
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  influxdb:
    image: influxdb:2.7-alpine
    environment:
      INFLUXDB_INIT_MODE: setup
      INFLUXDB_INIT_USERNAME: admin
      INFLUXDB_INIT_PASSWORD: password
      INFLUXDB_INIT_ORG: hydroponics
      INFLUXDB_INIT_BUCKET: sensor_data
      INFLUXDB_INIT_RETENTION: 30d
    ports:
      - "8086:8086"
    volumes:
      - influxdb_data:/var/lib/influxdb2

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  mosquitto:
    image: eclipse-mosquitto:2.0
    ports:
      - "1883:1883"
      - "9001:9001"
    volumes:
      - ./mosquitto/config:/mosquitto/config

  backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: hydroponics
      DB_PASSWORD: password
      DB_NAME: hydroponics
      INFLUX_URL: http://influxdb:8086
      INFLUX_TOKEN: admin:password
      INFLUX_ORG: hydroponics
      INFLUX_BUCKET: sensor_data
      MQTT_BROKER: mosquitto:1883
      REDIS_HOST: redis:6379
    depends_on:
      - postgres
      - influxdb
      - redis
      - mosquitto

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    depends_on:
      - backend

volumes:
  postgres_data:
  influxdb_data:
```

---

## 12. Business Model

### 12.1 Packaging & Pricing

| Package | Price/เดือน | ฟาร์ม | Users | Features |
|---------|-------------|-------|-------|----------|
| **Starter** | ฿990 | 1 | 1 | พื้นฐาน, ควบคุมระยะไกล |
| **Pro** | ฿2,990 | 5 | 5 | Starter + AI แนะนำ, วิเคราะห์ข้อมูล |
| **Business** | ฿9,900 | 20 | 20 | Pro + Multi-farm, API, Report |
| **Enterprise** | ฿29,900 | Unlimited | Unlimited | Business + On-premise, Support |

### 12.2 Hardware Packages

| Package | Price | Components |
|---------|-------|------------|
| **Basic Kit** | ฿8,650 | ESP32, Sensors, Relay, Pumps |
| **Standard Kit** | ฿15,900 | Basic + LED Grow Light, Fans |
| **Premium Kit** | ฿25,900 | Standard + pH/EC Pro, LCD, Valve |
| **Commercial Kit** | ฿49,900 | Premium x 2 + Backup System |

### 12.3 Revenue Model

```
┌─────────────────────────────────────────────────────────────────┐
│                    Revenue Model                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │  Hardware Sales (One-time)                               │ │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐       │ │
│  │  │ Basic Kit  │  │ Standard   │  │ Premium    │       │ │
│  │  │ ฿8,650     │  │ ฿15,900   │  │ ฿25,900   │       │ │
│  │  └────────────┘  └────────────┘  └────────────┘       │ │
│  └───────────────────────────────────────────────────────────┘ │
│                              +                                  │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │  Subscription (Recurring)                                │ │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐       │ │
│  │  │  Starter   │  │    Pro     │  │  Business  │       │ │
│  │  │  ฿990/เดือน│  │  ฿2,990/เดือน│  │  ฿9,900/เดือน│   │ │
│  │  └────────────┘  └────────────┘  └────────────┘       │ │
│  └───────────────────────────────────────────────────────────┘ │
│                              +                                  │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │  Professional Services                                   │ │
│  │  ┌────────────────────────────────────────────────────┐ │ │
│  │  │  Installation    : ฿5,000-15,000                  │ │ │
│  │  │  Training        : ฿3,000-10,000                  │ │ │
│  │  │  Customization   : ฿10,000-50,000                 │ │ │
│  │  │  Consulting      : ฿1,500/ชั่วโมง                 │ │ │
│  │  └────────────────────────────────────────────────────┘ │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 12.4 Target Market

| กลุ่ม | ขนาด | ราคาที่เหมาะสม | ช่องทางขาย |
|------|------|---------------|-----------|
| **ร้านอาหาร** | 50,000+ | 15,000-50,000 บาท | Direct, Distributor |
| **โรงแรม** | 10,000+ | 50,000-200,000 บาท | Direct, Event |
| **โรงเรียน** | 30,000+ | 20,000-100,000 บาท | Direct, Government |
| **คอนโด/บ้าน** | 1,000,000+ | 10,000-30,000 บาท | Online, Retail |
| **สตาร์ทอัพเกษตร** | 500+ | 50,000-200,000 บาท | Direct, Event |
| **ชุมชน/สหกรณ์** | 10,000+ | 100,000-500,000 บาท | Government, NGO |

### 12.5 Go-to-Market Strategy

| ระยะ | กิจกรรม | เป้าหมาย |
|------|---------|----------|
| **Q1** | เปิดตัวผลิตภัณฑ์, งานมหกรรมเกษตร | 50 ฟาร์ม |
| **Q2** | โรดโชว์ 5 จังหวัด, ร้านอาหารพรีเมียม | 200 ฟาร์ม |
| **Q3** | โปรโมชั่น Back to School, หน่วยงานรัฐ | 500 ฟาร์ม |
| **Q4** | Campaign สิ้นปี, Export ตลาด CLMV | 1,000 ฟาร์ม |

### 12.6 Channel Partners

| Partner Type | Margin | Role |
|--------------|--------|------|
| **Distributor** | 20-30% | ขายปลีก, ติดตั้ง |
| **System Integrator** | 15-25% | โซลูชันครบวงจร |
| **Reseller** | 10-20% | ขายต่อ |
| **Affiliate** | 5-10% | แนะนำลูกค้า |
| **Trainer** | 5% | อบรมการใช้งาน |

---

## สรุป

ระบบ Mini Smart Farm Hydroponics เป็นโซลูชันที่ครบวงจรสำหรับการทำเกษตรสมัยใหม่ โดยมี:

1. **Hardware** - อุปกรณ์ IoT ราคาประหยัด ใช้งานง่าย
2. **Software** - ระบบควบคุมและตรวจสอบอัจฉริยะ
3. **AI** - คำแนะนำและวิเคราะห์ข้อมูลอัตโนมัติ
4. **Mobile App** - ควบคุมจากที่ไหนก็ได้
5. **Business Model** - ชัดเจน ขายได้จริง

เหมาะสำหรับการขายให้กับ:
- ร้านอาหาร/โรงแรมที่ต้องการผักสดคุณภาพ
- โรงเรียน/มหาวิทยาลัยเพื่อการเรียนการสอน
- บ้านพักอาศัยยุคใหม่
- ผู้ประกอบการเกษตรรุ่นใหม่