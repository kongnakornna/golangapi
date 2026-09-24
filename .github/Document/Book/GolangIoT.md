# 📡 Go Programming for IoT: 10 Real-World Case Studies with Open Source Maps

> **สารบัญ (Table of Contents)**
> - [บทนำ](#บทนำ)
> - [สถาปัตยกรรม IoT กับ Go](#สถาปัตยกรรม-iot-กับ-go)
> - [10 กรณีศึกษา (Case Studies)](#10-กรณีศึกษา-case-studies)
>   - [Case 1: Live GPS Tracker - OsmAnd + Leaflet](#case-1-live-gps-tracker---osmand--leaflet)
>   - [Case 2: Edge IoT Gateway - MQTT + Kafka + Raspberry Pi](#case-2-edge-iot-gateway---mqtt--kafka--raspberry-pi)
>   - [Case 3: EcoTracker Map - OpenStreetMap + React + Go](#case-3-ecotracker-map---openstreetmap--react--go)
>   - [Case 4: Device Hub - Local IoT Data Infrastructure](#case-4-device-hub---local-iot-data-infrastructure)
>   - [Case 5: Multi-Source Fusion Positioning - GPS/WiFi/UWB/IMU](#case-5-multi-source-fusion-positioning---gpswifiuwbimu)
>   - [Case 6: Smart Container Tracking - LoRaWAN + ESP32](#case-6-smart-container-tracking---lorawan--esp32)
>   - [Case 7: Cloud Microservices IoT Platform](#case-7-cloud-microservices-iot-platform)
>   - [Case 8: 3D OSM Visualization with Go Backend](#case-8-3d-osm-visualization-with-go-backend)
>   - [Case 9: Garden Monitor Dataflow](#case-9-garden-monitor-dataflow)
>   - [Case 10: Fleet Management System](#case-10-fleet-management-system)
> - [เทมเพลตและโค้ดที่รันได้จริง](#เทมเพลตและโค้ดที่รันได้จริง)
> - [แบบฝึกหัดท้ายบท](#แบบฝึกหัดท้ายบท)
> - [สรุป](#สรุป)

---

## บทนำ (Introduction)

> **สรุปสั้น:** บทนี้รวบรวม 10 กรณีศึกษา (Case Studies) เกี่ยวกับการพัฒนา IoT (Internet of Things) ด้วยภาษา Go ร่วมกับ Open Source Maps (OSM, Leaflet, MapLibre) สำหรับงาน Tracking, Monitoring และ Visualization แบบเรียลไทม์

### วัตถุประสงค์ (Objectives)

| วัตถุประสงค์ | คำอธิบาย |
|-------------|----------|
| **เรียนรู้สถาปัตยกรรม IoT** | เข้าใจโครงสร้างระบบ IoT ด้วย Go |
| **ประยุกต์ใช้ Open Source Maps** | ใช้ Leaflet, MapLibre, OpenStreetMap ฟรี |
| **สร้างระบบ Tracking** | GPS Tracking, Asset Tracking, Real-time Monitoring |
| **จัดการข้อมูล Sensor** | MQTT, Kafka, WebSocket, Data Aggregation |
| **ปรับใช้บน Edge Device** | Raspberry Pi, ESP32, Docker |

### กลุ่มเป้าหมาย (Target Audience)

| กลุ่ม | ความเหมาะสม | เหตุผล |
|------|-------------|---------|
| **นักพัฒนา IoT** | ⭐⭐⭐⭐⭐ | เรียนรู้การประยุกต์ใช้ Go ในงาน IoT |
| **DevOps/System Engineer** | ⭐⭐⭐⭐ | เข้าใจ Edge Computing, Data Pipeline |
| **นักศึกษา/นักวิจัย** | ⭐⭐⭐⭐ | มีกรณีศึกษาและโค้ดพร้อมใช้ |
| **ผู้สนใจ Open Source** | ⭐⭐⭐⭐⭐ | ใช้ OSM, Leaflet, MQTT ฟรี |

---

## สถาปัตยกรรม IoT กับ Go (IoT Architecture with Go)

### ภาพรวมสถาปัตยกรรม (Architecture Overview)

```
รูปที่ 1: สถาปัตยกรรมระบบ IoT Tracking แบบสมบูรณ์

┌─────────────────────────────────────────────────────────────────────────────┐
│                           IoT TRACKING ARCHITECTURE                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐                  │
│  │   DEVICE     │    │   GATEWAY    │    │   CLOUD      │                  │
│  │   LAYER      │───▶│   LAYER      │───▶│   LAYER      │                  │
│  │              │    │              │    │              │                  │
│  │ • GPS Module │    │ • MQTT       │    │ • Kafka      │                  │
│  │ • ESP32      │    │ • Go Service │    │ • PostgreSQL │                  │
│  │ • LoRaWAN    │    │ • SQLite     │    │ • TimescaleDB│                  │
│  │ • BLE        │    │ • Buffer     │    │ • Redis      │                  │
│  └──────────────┘    └──────────────┘    └──────────────┘                  │
│                              │                    │                          │
│                              ▼                    ▼                          │
│                    ┌──────────────┐    ┌──────────────┐                    │
│                    │   MAP UI     │    │   WEB UI     │                    │
│                    │   LAYER      │    │   LAYER      │                    │
│                    │              │    │              │                    │
│                    │ • Leaflet    │◀───│ • React      │                    │
│                    │ • MapLibre   │    │ • WebSocket  │                    │
│                    │ • OpenStreet │    │ • REST API   │                    │
│                    │   Map        │    │              │                    │
│                    └──────────────┘    └──────────────┘                    │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Open Source Maps ที่ใช้ฟรี (Free Open Source Maps)

| แผนที่ (Map) | ไลบรารี (Library) | การใช้งาน (Usage) | ข้อดี (Advantage) |
|--------------|-------------------|-------------------|-------------------|
| **OpenStreetMap** | Leaflet.js | แผนที่พื้นฐาน, เลเยอร์ | ฟรี, ข้อมูลสมบูรณ์ |
| **MapLibre GL** | MapLibre | Render 3D, Vector Tiles | เร็ว, สวยงาม |
| **OpenLayers** | OpenLayers | GIS ขั้นสูง | รองรับหลายรูปแบบ |

---

## 10 กรณีศึกษา (10 Case Studies)

### Case 1: Live GPS Tracker - OsmAnd + Leaflet

#### โครงสร้างการทำงาน (Architecture)

```
รูปที่ 2: Live GPS Tracker Dataflow

    [OsmAnd App]          [Go Backend]           [Web Browser]
         │                     │                       │
         │  HTTP POST          │                       │
         │  /track?token=xxx   │                       │
         │  &lat={0}&lon={1}   │                       │
         │────────────────────▶│                       │
         │                     │                       │
         │                     │  Store to SQLite      │
         │                     │──────────────────────▶│
         │                     │                       │
         │                     │  WebSocket Push       │
         │                     │──────────────────────▶│
         │                     │                       │
         │                     │                       │  [Leaflet Map]
         │                     │                       │  Real-time Update
```

#### โค้ด Backend (Go)

```go
// File: main.go - Live GPS Tracker Server
// ไฟล์: main.go - เซิร์ฟเวอร์ติดตามตำแหน่งแบบเรียลไทม์

package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"
    "sync"
    "time"
    
    "github.com/gorilla/websocket"
    _ "github.com/mattn/go-sqlite3"
)

// Location โครงสร้างข้อมูลตำแหน่ง GPS
// Location structure for GPS data
type Location struct {
    Lat       float64   `json:"lat"`        // ละติจูด / Latitude
    Lon       float64   `json:"lon"`        // ลองจิจูด / Longitude
    Timestamp time.Time `json:"timestamp"`  // เวลาที่บันทึก / Record time
    Speed     float64   `json:"speed"`      // ความเร็ว (km/h) / Speed
    Altitude  float64   `json:"altitude"`   // ความสูง (m) / Altitude
}

// WebSocket Hub สำหรับจัดการ connection ทั้งหมด
// WebSocket Hub for managing all connections
type Hub struct {
    clients    map[*websocket.Conn]bool
    broadcast  chan Location
    register   chan *websocket.Conn
    unregister chan *websocket.Conn
    mu         sync.RWMutex
}

// NewHub สร้าง Hub ใหม่
// NewHub creates a new Hub
func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*websocket.Conn]bool),
        broadcast:  make(chan Location),
        register:   make(chan *websocket.Conn),
        unregister: make(chan *websocket.Conn),
    }
}

// Run เริ่มการทำงานของ Hub
// Run starts the Hub
func (h *Hub) Run() {
    for {
        select {
        case conn := <-h.register:
            h.mu.Lock()
            h.clients[conn] = true
            h.mu.Unlock()
            log.Println("New WebSocket client connected")
            
        case conn := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[conn]; ok {
                delete(h.clients, conn)
                conn.Close()
            }
            h.mu.Unlock()
            log.Println("WebSocket client disconnected")
            
        case location := <-h.broadcast:
            h.mu.RLock()
            for conn := range h.clients {
                err := conn.WriteJSON(location)
                if err != nil {
                    conn.Close()
                    delete(h.clients, conn)
                }
            }
            h.mu.RUnlock()
        }
    }
}

var (
    hub     = NewHub()
    upgrader = websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool { return true },
    }
    db *sql.DB
)

func initDB() {
    var err error
    // เปิดเชื่อมต่อ SQLite database
    // Open SQLite database connection
    db, err = sql.Open("sqlite3", "./tracker.db")
    if err != nil {
        log.Fatal(err)
    }
    
    // สร้างตารางเก็บตำแหน่ง
    // Create location table
    createTableSQL := `
    CREATE TABLE IF NOT EXISTS locations (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        lat REAL NOT NULL,
        lon REAL NOT NULL,
        timestamp DATETIME NOT NULL,
        speed REAL,
        altitude REAL
    );`
    
    _, err = db.Exec(createTableSQL)
    if err != nil {
        log.Fatal(err)
    }
    
    // สร้าง index เพื่อเพิ่มความเร็วในการค้นหา
    // Create index for faster queries
    _, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_timestamp ON locations(timestamp)")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Database initialized successfully")
}

// saveLocation บันทึกตำแหน่งลงฐานข้อมูล
// saveLocation saves location to database
func saveLocation(loc Location) error {
    _, err := db.Exec(
        "INSERT INTO locations (lat, lon, timestamp, speed, altitude) VALUES (?, ?, ?, ?, ?)",
        loc.Lat, loc.Lon, loc.Timestamp, loc.Speed, loc.Altitude,
    )
    return err
}

// getRecentLocations ดึงตำแหน่งล่าสุด (3 ชั่วโมง)
// getRecentLocations gets recent locations (last 3 hours)
func getRecentLocations() ([]Location, error) {
    threeHoursAgo := time.Now().Add(-3 * time.Hour)
    
    rows, err := db.Query(
        "SELECT lat, lon, timestamp, speed, altitude FROM locations WHERE timestamp > ? ORDER BY timestamp DESC",
        threeHoursAgo,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var locations []Location
    for rows.Next() {
        var loc Location
        err := rows.Scan(&loc.Lat, &loc.Lon, &loc.Timestamp, &loc.Speed, &loc.Altitude)
        if err != nil {
            return nil, err
        }
        locations = append(locations, loc)
    }
    
    return locations, nil
}

// trackHandler รับข้อมูลจาก OsmAnd
// trackHandler receives data from OsmAnd
func trackHandler(w http.ResponseWriter, r *http.Request) {
    // ตรวจสอบ API Token
    // Verify API Token
    token := r.URL.Query().Get("token")
    if token != "yourtoken" {  // ควรเปลี่ยนเป็นค่าจริง / Should change to actual value
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    // ดึงค่าจาก query parameters
    // Extract values from query parameters
    latStr := r.URL.Query().Get("lat")
    lonStr := r.URL.Query().Get("lon")
    timestampStr := r.URL.Query().Get("timestamp")
    speedStr := r.URL.Query().Get("speed")
    altitudeStr := r.URL.Query().Get("altitude")
    
    if latStr == "" || lonStr == "" {
        http.Error(w, "Missing lat or lon", http.StatusBadRequest)
        return
    }
    
    lat, _ := strconv.ParseFloat(latStr, 64)
    lon, _ := strconv.ParseFloat(lonStr, 64)
    speed, _ := strconv.ParseFloat(speedStr, 64)
    altitude, _ := strconv.ParseFloat(altitudeStr, 64)
    
    // แปลง timestamp (Unix timestamp)
    // Parse timestamp
    var timestamp time.Time
    if timestampStr != "" {
        ts, _ := strconv.ParseInt(timestampStr, 10, 64)
        timestamp = time.Unix(ts, 0)
    } else {
        timestamp = time.Now()
    }
    
    loc := Location{
        Lat:       lat,
        Lon:       lon,
        Timestamp: timestamp,
        Speed:     speed,
        Altitude:  altitude,
    }
    
    // บันทึกและ broadcast
    // Save and broadcast
    if err := saveLocation(loc); err != nil {
        log.Printf("Error saving location: %v", err)
    }
    
    hub.broadcast <- loc
    
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

// wsHandler จัดการ WebSocket connections
// wsHandler handles WebSocket connections
func wsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket upgrade error: %v", err)
        return
    }
    
    hub.register <- conn
    
    // ส่งประวัติล่าสุดเมื่อเชื่อมต่อ
    // Send recent history on connection
    recent, err := getRecentLocations()
    if err == nil && len(recent) > 0 {
        for _, loc := range recent {
            conn.WriteJSON(loc)
        }
    }
    
    // รับ message จาก client (ping/pong)
    // Handle client messages
    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            hub.unregister <- conn
            break
        }
    }
}

// mapHandler ส่งหน้า HTML
// mapHandler serves HTML page
func mapHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte(mapHTML))
}

const mapHTML = `<!DOCTYPE html>
<html>
<head>
    <title>Live GPS Tracker</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" />
    <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
    <style>
        #map { height: 100vh; margin: 0; padding: 0; }
        .info-panel {
            position: absolute;
            top: 10px;
            right: 10px;
            background: white;
            padding: 10px;
            border-radius: 5px;
            z-index: 1000;
            box-shadow: 0 2px 5px rgba(0,0,0,0.2);
        }
        .live-badge {
            color: red;
            animation: blink 1s infinite;
        }
        @keyframes blink {
            0% { opacity: 1; }
            50% { opacity: 0.3; }
            100% { opacity: 1; }
        }
    </style>
</head>
<body>
    <div id="map"></div>
    <div class="info-panel">
        <strong>📍 Live GPS Tracking</strong>
        <span class="live-badge">● LIVE</span>
        <div id="coords">Waiting for GPS...</div>
        <div id="speed">Speed: -- km/h</div>
    </div>

    <script>
        // เริ่มต้นแผนที่
        // Initialize map
        var map = L.map('map').setView([13.736717, 100.523186], 13);
        
        // เพิ่ม OpenStreetMap layer
        // Add OpenStreetMap layer
        L.tileLayer('https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png', {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OSM</a>'
        }).addTo(map);
        
        // ตัวแปรเก็บ marker และ polyline
        // Variables for marker and polyline
        var currentMarker = null;
        var pathPoints = [];
        var pathLine = null;
        
        // เชื่อมต่อ WebSocket
        // Connect WebSocket
        var ws = new WebSocket('ws://' + window.location.host + '/ws');
        
        ws.onmessage = function(event) {
            var data = JSON.parse(event.data);
            updateMap(data);
        };
        
        function updateMap(loc) {
            var lat = loc.lat;
            var lon = loc.lon;
            var speed = loc.speed || 0;
            
            // อัปเดตข้อมูล
            // Update info panel
            document.getElementById('coords').innerHTML = 
                '📍 ' + lat.toFixed(6) + ', ' + lon.toFixed(6);
            document.getElementById('speed').innerHTML = 
                '⚡ Speed: ' + speed.toFixed(1) + ' km/h';
            
            // อัปเดต marker
            // Update marker
            if (currentMarker) {
                currentMarker.setLatLng([lat, lon]);
            } else {
                currentMarker = L.marker([lat, lon]).addTo(map);
                currentMarker.bindPopup('Current Location').openPopup();
            }
            
            // เพิ่มจุดในเส้นทาง
            // Add point to path
            pathPoints.push([lat, lon]);
            
            // จำกัดจำนวนจุด (200 จุดล่าสุด)
            // Limit points (last 200 points)
            if (pathPoints.length > 200) {
                pathPoints.shift();
            }
            
            // วาดเส้นทาง
            // Draw path
            if (pathLine) {
                pathLine.setLatLngs(pathPoints);
            } else {
                pathLine = L.polyline(pathPoints, {
                    color: '#3388ff',
                    weight: 3,
                    opacity: 0.7
                }).addTo(map);
            }
            
            // เลื่อนแผนที่ตามตำแหน่ง
            // Center map on location
            map.setView([lat, lon], map.getZoom());
        }
        
        // ฟังก์ชันซูมไปยังตำแหน่งปัจจุบัน
        // Function to zoom to current location
        function zoomToCurrent() {
            if (currentMarker) {
                map.setView(currentMarker.getLatLng(), 15);
            }
        }
        
        console.log('Live Tracker Ready!');
    </script>
</body>
</html>`

func main() {
    // เริ่มต้นฐานข้อมูล
    // Initialize database
    initDB()
    defer db.Close()
    
    // เริ่ม Hub
    // Start Hub
    go hub.Run()
    
    // ตั้งค่า Routes
    // Setup routes
    http.HandleFunc("/track", trackHandler)
    http.HandleFunc("/ws", wsHandler)
    http.HandleFunc("/", mapHandler)
    
    log.Println("Server starting on :8080")
    log.Println("OsmAnd URL: http://YOUR_IP:8080/track?token=yourtoken&lat={0}&lon={1}&timestamp={2}&speed={3}&altitude={4}")
    
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
```

#### การใช้งาน (Usage)

```bash
# 1. ติดตั้ง Dependencies
go mod init livetracker
go get github.com/gorilla/websocket
go get github.com/mattn/go-sqlite3

# 2. รันเซิร์ฟเวอร์
go run main.go

# 3. ตั้งค่า OsmAnd App
#   - เปิด OsmAnd → Trip Recording → Online Tracking
#   - เพิ่ม URL: http://YOUR_IP:8080/track?token=yourtoken&lat={0}&lon={1}&timestamp={2}&speed={3}&altitude={4}
#   - เริ่มบันทึกเส้นทาง

# 4. เปิด Browser: http://localhost:8080
```

---

### Case 2: Edge IoT Gateway - MQTT + Kafka + Raspberry Pi

#### โครงสร้างการทำงาน (Architecture)

```
รูปที่ 3: Edge IoT Gateway Dataflow

    [IoT Sensors]     [MQTT Broker]      [Edge Gateway]       [Kafka]        [Analytics]
         │                 │                   │                  │               │
         │  MQTT Publish   │                   │                  │               │
         │────────────────▶│                   │                  │               │
         │                 │  MQTT Subscribe    │                  │               │
         │                 │──────────────────▶│                  │               │
         │                 │                   │                  │               │
         │                 │                   │  Filter &        │               │
         │                 │                   │  Aggregate       │               │
         │                 │                   │  (SQLite Buffer) │               │
         │                 │                   │                  │               │
         │                 │                   │  Kafka Produce   │               │
         │                 │                   │─────────────────▶│               │
         │                 │                   │                  │  Consume      │
         │                 │                   │                  │──────────────▶│
```

#### โค้ด Edge Gateway (Go)

```go
// File: main.go - Edge IoT Gateway
// ไฟล์: main.go - เกตเวย์ IoT บริเวณขอบเครือข่าย

package main

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    mqtt "github.com/eclipse/paho.mqtt.golang"
    "github.com/confluentinc/confluent-kafka-go/kafka"
    _ "github.com/mattn/go-sqlite3"
    "gopkg.in/yaml.v3"
)

// Config โครงสร้างการตั้งค่า
// Config structure
type Config struct {
    MQTT struct {
        Broker   string `yaml:"broker"`
        ClientID string `yaml:"client_id"`
        Topic    string `yaml:"topic"`
        QoS      int    `yaml:"qos"`
    } `yaml:"mqtt"`
    
    Kafka struct {
        Brokers  []string `yaml:"brokers"`
        Topic    string   `yaml:"topic"`
        ClientID string   `yaml:"client_id"`
    } `yaml:"kafka"`
    
    Processing struct {
        AggregationWindow int `yaml:"aggregation_window_seconds"`
        Rules            []Rule `yaml:"rules"`
    } `yaml:"processing"`
    
    Buffer struct {
        Path          string `yaml:"path"`
        MaxSizeMB     int    `yaml:"max_size_mb"`
        FlushInterval int    `yaml:"flush_interval_seconds"`
    } `yaml:"buffer"`
}

// Rule กฎการประมวลผล
// Processing rule
type Rule struct {
    Type     string `yaml:"type"`      // filter, aggregate, enrich
    Field    string `yaml:"field"`
    Operator string `yaml:"operator"`  // <, >, =, !=
    Value    interface{} `yaml:"value"`
    Action   string `yaml:"action"`    // drop, keep, transform
    Function string `yaml:"function"`  // average, sum, min, max
    Output   string `yaml:"output_name"`
}

// SensorData โครงสร้างข้อมูลจากเซนเซอร์
// Sensor data structure
type SensorData struct {
    DeviceID    string    `json:"device_id"`
    Temperature float64   `json:"temperature"`
    Humidity    float64   `json:"humidity"`
    Pressure    float64   `json:"pressure"`
    Timestamp   time.Time `json:"timestamp"`
    Quality     float64   `json:"quality_score"`
}

// AggregatedData โครงสร้างข้อมูลหลังการรวม
// Aggregated data structure
type AggregatedData struct {
    DeviceID      string    `json:"device_id"`
    AvgTemp       float64   `json:"avg_temperature"`
    AvgHumidity   float64   `json:"avg_humidity"`
    MinTemp       float64   `json:"min_temperature"`
    MaxTemp       float64   `json:"max_temperature"`
    SampleCount   int       `json:"sample_count"`
    WindowStart   time.Time `json:"window_start"`
    WindowEnd     time.Time `json:"window_end"`
    GatewayID     string    `json:"gateway_id"`
}

type EdgeGateway struct {
    config      *Config
    mqttClient  mqtt.Client
    kafkaProd   *kafka.Producer
    db          *sql.DB
    dataBuffer  []SensorData
    windowData  map[string][]SensorData  // key: device_id
    msgChan     chan SensorData
    stopChan    chan struct{}
}

// NewEdgeGateway สร้าง Edge Gateway ใหม่
// NewEdgeGateway creates a new Edge Gateway
func NewEdgeGateway(configPath string) (*EdgeGateway, error) {
    // โหลด configuration
    // Load configuration
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("read config error: %v", err)
    }
    
    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("parse config error: %v", err)
    }
    
    // เปิด SQLite buffer
    // Open SQLite buffer
    db, err := sql.Open("sqlite3", config.Buffer.Path)
    if err != nil {
        return nil, fmt.Errorf("open database error: %v", err)
    }
    
    // สร้างตาราง buffer
    // Create buffer table
    createTableSQL := `
    CREATE TABLE IF NOT EXISTS sensor_buffer (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        device_id TEXT NOT NULL,
        data JSON NOT NULL,
        timestamp DATETIME NOT NULL,
        retry_count INTEGER DEFAULT 0
    );`
    
    if _, err := db.Exec(createTableSQL); err != nil {
        return nil, err
    }
    
    return &EdgeGateway{
        config:     &config,
        db:         db,
        dataBuffer: make([]SensorData, 0),
        windowData: make(map[string][]SensorData),
        msgChan:    make(chan SensorData, 10000),
        stopChan:   make(chan struct{}),
    }, nil
}

// setupMQTT ตั้งค่า MQTT client
// setupMQTT configures MQTT client
func (g *EdgeGateway) setupMQTT() error {
    opts := mqtt.NewClientOptions()
    opts.AddBroker(g.config.MQTT.Broker)
    opts.SetClientID(g.config.MQTT.ClientID)
    opts.SetAutoReconnect(true)
    opts.SetConnectRetry(true)
    opts.SetOnConnectHandler(g.onMQTTConnect)
    opts.SetConnectionLostHandler(g.onMQTTLost)
    
    g.mqttClient = mqtt.NewClient(opts)
    
    if token := g.mqttClient.Connect(); token.Wait() && token.Error() != nil {
        return token.Error()
    }
    
    return nil
}

// onMQTTConnect เรียกเมื่อเชื่อมต่อ MQTT สำเร็จ
// onMQTTConnect called when MQTT connected
func (g *EdgeGateway) onMQTTConnect(client mqtt.Client) {
    log.Println("MQTT Connected, subscribing to topics...")
    
    token := client.Subscribe(g.config.MQTT.Topic, byte(g.config.MQTT.QoS), g.onMessage)
    if token.Wait() && token.Error() != nil {
        log.Printf("Subscribe error: %v", token.Error())
    }
}

// onMQTTLost เรียกเมื่อสูญเสียการเชื่อมต่อ MQTT
// onMQTTLost called when MQTT connection lost
func (g *EdgeGateway) onMQTTLost(client mqtt.Client, err error) {
    log.Printf("MQTT Connection lost: %v", err)
}

// onMessage รับ和处理ข้อความ MQTT
// onMessage receives and processes MQTT messages
func (g *EdgeGateway) onMessage(client mqtt.Client, msg mqtt.Message) {
    log.Printf("Received message on topic: %s", msg.Topic())
    
    var data SensorData
    if err := json.Unmarshal(msg.Payload(), &data); err != nil {
        log.Printf("JSON parse error: %v", err)
        return
    }
    
    data.Timestamp = time.Now()
    
    // ประมวลผลและส่งเข้า channel
    // Process and send to channel
    if g.applyFilterRules(&data) {
        g.msgChan <- data
    }
}

// applyFilterRules ใช้กฎการกรอง
// applyFilterRules applies filter rules
func (g *EdgeGateway) applyFilterRules(data *SensorData) bool {
    for _, rule := range g.config.Processing.Rules {
        if rule.Type != "filter" {
            continue
        }
        
        var value float64
        switch rule.Field {
        case "temperature":
            value = data.Temperature
        case "humidity":
            value = data.Humidity
        case "pressure":
            value = data.Pressure
        case "quality_score":
            value = data.Quality
        default:
            continue
        }
        
        ruleValue, ok := rule.Value.(float64)
        if !ok {
            continue
        }
        
        // ตรวจสอบเงื่อนไข
        // Check condition
        switch rule.Operator {
        case "<":
            if value < ruleValue {
                return rule.Action != "drop"
            }
        case ">":
            if value > ruleValue {
                return rule.Action != "drop"
            }
        case "=":
            if value == ruleValue {
                return rule.Action != "drop"
            }
        }
        
        if rule.Action == "drop" {
            log.Printf("Dropped data: %s=%f violates rule %s", rule.Field, value, rule.Operator)
            return false
        }
    }
    return true
}

// setupKafka ตั้งค่า Kafka producer
// setupKafka configures Kafka producer
func (g *EdgeGateway) setupKafka() error {
    conf := &kafka.ConfigMap{
        "bootstrap.servers":     strings.Join(g.config.Kafka.Brokers, ","),
        "client.id":            g.config.Kafka.ClientID,
        "acks":                 "all",
        "retries":              5,
        "retry.backoff.ms":     100,
        "enable.idempotence":   true,
    }
    
    prod, err := kafka.NewProducer(conf)
    if err != nil {
        return err
    }
    
    g.kafkaProd = prod
    
    // เริ่ม goroutine สำหรับ delivery reports
    // Start goroutine for delivery reports
    go func() {
        for e := range prod.Events() {
            switch ev := e.(type) {
            case *kafka.Message:
                if ev.TopicPartition.Error != nil {
                    log.Printf("Kafka delivery failed: %v", ev.TopicPartition.Error)
                    // บันทึกไว้ใน buffer เพื่อ retry
                    // Save to buffer for retry
                    g.saveToBuffer(ev.Value)
                } else {
                    log.Printf("Kafka delivered to %v", ev.TopicPartition)
                }
            }
        }
    }()
    
    return nil
}

// saveToBuffer บันทึกข้อมูลลง SQLite buffer
// saveToBuffer saves data to SQLite buffer
func (g *EdgeGateway) saveToBuffer(data []byte) error {
    _, err := g.db.Exec(
        "INSERT INTO sensor_buffer (device_id, data, timestamp) VALUES (?, ?, ?)",
        "unknown", string(data), time.Now(),
    )
    return err
}

// aggregateData รวมข้อมูลตามช่วงเวลา
// aggregateData aggregates data by time window
func (g *EdgeGateway) aggregateData(deviceID string, samples []SensorData) *AggregatedData {
    if len(samples) == 0 {
        return nil
    }
    
    var sumTemp, sumHumidity float64
    var minTemp, maxTemp float64 = 1e9, -1e9
    
    for _, s := range samples {
        sumTemp += s.Temperature
        sumHumidity += s.Humidity
        
        if s.Temperature < minTemp {
            minTemp = s.Temperature
        }
        if s.Temperature > maxTemp {
            maxTemp = s.Temperature
        }
    }
    
    return &AggregatedData{
        DeviceID:    deviceID,
        AvgTemp:     sumTemp / float64(len(samples)),
        AvgHumidity: sumHumidity / float64(len(samples)),
        MinTemp:     minTemp,
        MaxTemp:     maxTemp,
        SampleCount: len(samples),
        WindowStart: samples[0].Timestamp,
        WindowEnd:   samples[len(samples)-1].Timestamp,
        GatewayID:   g.config.MQTT.ClientID,
    }
}

// processWorker ประมวลผลข้อมูลจาก channel
// processWorker processes data from channel
func (g *EdgeGateway) processWorker() {
    ticker := time.NewTicker(time.Duration(g.config.Processing.AggregationWindow) * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case data := <-g.msgChan:
            // เก็บข้อมูลใน window
            // Store data in window
            g.windowData[data.DeviceID] = append(g.windowData[data.DeviceID], data)
            
        case <-ticker.C:
            // ส่งข้อมูลที่รวมแล้วไป Kafka
            // Send aggregated data to Kafka
            for deviceID, samples := range g.windowData {
                if len(samples) == 0 {
                    continue
                }
                
                aggregated := g.aggregateData(deviceID, samples)
                if aggregated == nil {
                    continue
                }
                
                // ส่งไป Kafka
                // Send to Kafka
                jsonData, err := json.Marshal(aggregated)
                if err != nil {
                    log.Printf("JSON marshal error: %v", err)
                    continue
                }
                
                msg := &kafka.Message{
                    TopicPartition: kafka.TopicPartition{
                        Topic:     &g.config.Kafka.Topic,
                        Partition: kafka.PartitionAny,
                    },
                    Value: jsonData,
                    Timestamp: time.Now(),
                }
                
                if err := g.kafkaProd.Produce(msg, nil); err != nil {
                    log.Printf("Kafka produce error: %v", err)
                    g.saveToBuffer(jsonData)
                } else {
                    log.Printf("Sent aggregated data for %s: avgTemp=%.2f, samples=%d",
                        deviceID, aggregated.AvgTemp, aggregated.SampleCount)
                }
                
                // เคลียร์ window
                // Clear window
                g.windowData[deviceID] = make([]SensorData, 0)
            }
        }
    }
}

// retryFailedMessages ลองส่งข้อความที่ล้มเหลวอีกครั้ง
// retryFailedMessages retries failed messages
func (g *EdgeGateway) retryFailedMessages() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        rows, err := g.db.Query(
            "SELECT id, data FROM sensor_buffer WHERE retry_count < 5",
        )
        if err != nil {
            log.Printf("Query buffer error: %v", err)
            continue
        }
        
        var ids []int
        var dataList [][]byte
        
        for rows.Next() {
            var id int
            var dataStr string
            if err := rows.Scan(&id, &dataStr); err != nil {
                continue
            }
            ids = append(ids, id)
            dataList = append(dataList, []byte(dataStr))
        }
        rows.Close()
        
        for i, data := range dataList {
            msg := &kafka.Message{
                TopicPartition: kafka.TopicPartition{
                    Topic: &g.config.Kafka.Topic,
                },
                Value: data,
            }
            
            if err := g.kafkaProd.Produce(msg, nil); err == nil {
                // ลบจาก buffer
                // Delete from buffer
                g.db.Exec("DELETE FROM sensor_buffer WHERE id = ?", ids[i])
                log.Printf("Retried message %d sent successfully", ids[i])
            } else {
                // อัปเดต retry count
                // Update retry count
                g.db.Exec("UPDATE sensor_buffer SET retry_count = retry_count + 1 WHERE id = ?", ids[i])
                log.Printf("Retry failed for message %d: %v", ids[i], err)
            }
        }
    }
}

// Run เริ่มทำงาน Edge Gateway
// Run starts the Edge Gateway
func (g *EdgeGateway) Run() error {
    // ตั้งค่า MQTT
    if err := g.setupMQTT(); err != nil {
        return fmt.Errorf("MQTT setup error: %v", err)
    }
    
    // ตั้งค่า Kafka
    if err := g.setupKafka(); err != nil {
        return fmt.Errorf("Kafka setup error: %v", err)
    }
    
    // เริ่ม workers
    go g.processWorker()
    go g.retryFailedMessages()
    
    log.Println("Edge Gateway started successfully")
    log.Printf("Listening on MQTT topic: %s", g.config.MQTT.Topic)
    log.Printf("Forwarding to Kafka: %v", g.config.Kafka.Brokers)
    
    // รอ signal สำหรับ graceful shutdown
    // Wait for graceful shutdown signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
    
    log.Println("Shutting down gracefully...")
    g.mqttClient.Disconnect(250)
    g.kafkaProd.Flush(5000)
    g.kafkaProd.Close()
    g.db.Close()
    
    return nil
}

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Usage: go run main.go <config.yaml>")
    }
    
    gateway, err := NewEdgeGateway(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }
    
    if err := gateway.Run(); err != nil {
        log.Fatal(err)
    }
}
```

#### ไฟล์ Configuration (config.yaml)

```yaml
# config.yaml - การตั้งค่า Edge Gateway
# config.yaml - Edge Gateway configuration

mqtt:
  broker: "tcp://localhost:1883"      # MQTT broker address
  client_id: "edge-gateway-01"        # Gateway identifier
  topic: "sensors/+/data"             # Topic pattern
  qos: 1                               # Quality of Service

kafka:
  brokers:
    - "localhost:9092"                 # Kafka broker address
  topic: "iot-sensor-aggregated"       # Output topic
  client_id: "edge-producer"

processing:
  aggregation_window_seconds: 60       # 60-second aggregation window
  rules:
    - type: "filter"                   # Filter out invalid temperatures
      field: "temperature"
      operator: "<"
      value: -50
      action: "drop"
    - type: "filter"
      field: "temperature"
      operator: ">"
      value: 100
      action: "drop"
    - type: "filter"                   # Filter low quality data
      field: "quality_score"
      operator: "<"
      value: 0.5
      action: "drop"

buffer:
  path: "./data/buffer.db"             # SQLite buffer path
  max_size_mb: 100                     # Maximum buffer size
  flush_interval_seconds: 30           # Buffer flush interval
```

#### การติดตั้งบน Raspberry Pi

```bash
# 1. SSH เข้า Raspberry Pi
ssh pi@raspberrypi.local

# 2. ติดตั้ง Go (ถ้ายังไม่มี)
wget https://go.dev/dl/go1.22.0.linux-arm64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-arm64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# 3. ติดตั้ง Dependencies
sudo apt update
sudo apt install -y mosquitto mosquitto-clients sqlite3

# 4. โคลนและ build โปรเจค
git clone https://github.com/your-repo/edge-gateway.git
cd edge-gateway
go mod download
go build -o edge-gateway

# 5. สร้าง systemd service
sudo cat > /etc/systemd/system/edge-gateway.service << EOF
[Unit]
Description=Edge IoT Gateway
After=network.target mosquitto.service

[Service]
Type=simple
User=pi
WorkingDirectory=/home/pi/edge-gateway
ExecStart=/home/pi/edge-gateway/edge-gateway config.yaml
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# 6. เริ่ม service
sudo systemctl daemon-reload
sudo systemctl enable edge-gateway
sudo systemctl start edge-gateway

# 7. ตรวจสอบ logs
sudo journalctl -u edge-gateway -f
```

---

### Case 3: EcoTracker Map - OpenStreetMap + React + Go

#### โครงสร้างการทำงาน (Architecture)

```
รูปที่ 4: EcoTracker System Architecture

    [React Frontend]         [Go Backend]           [OpenStreetMap APIs]
         │                        │                          │
         │  HTTP Request          │                          │
         │  /api/points           │                          │
         │───────────────────────▶│                          │
         │                        │                          │
         │                        │  Nominatim API           │
         │                        │  (Geocoding)            │
         │                        │─────────────────────────▶│
         │                        │                          │
         │                        │  Overpass API            │
         │                        │  (POI Data)             │
         │                        │─────────────────────────▶│
         │                        │                          │
         │  JSON Response         │                          │
         │◀───────────────────────│                          │
         │                        │                          │
         │  MapLibre GL           │                          │
         │  Render Map            │                          │
```

#### โค้ด Backend (Go)

```go
// File: backend/main.go - EcoTracker Backend
// ไฟล์: backend/main.go - Backend ของ EcoTracker

package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "math"
    "math/rand"
    "net/http"
    "os"
    "os/signal"
    "strconv"
    "syscall"
    "time"
    
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

// EcoPoint จุดรวบรวมขยะรีไซเคิล
// EcoPoint represents recycling collection point
type EcoPoint struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Lat         float64 `json:"lat"`
    Lon         float64 `json:"lon"`
    Address     string  `json:"address"`
    Type        string  `json:"type"`      // plastic, glass, paper, electronic
    Capacity    int     `json:"capacity"`  // กำลังการผลิต (kg/day)
    CurrentLoad int     `json:"current_load"` // ปริมาณปัจจุบัน
    Status      string  `json:"status"`    // active, full, maintenance
    LastUpdate  string  `json:"last_update"`
}

// RouteResponse เส้นทางระหว่างจุด
// RouteResponse represents route between points
type RouteResponse struct {
    StartPoint  EcoPoint    `json:"start_point"`
    EndPoint    EcoPoint    `json:"end_point"`
    Distance    float64     `json:"distance_km"`   // ระยะทาง (km)
    Duration    int         `json:"duration_min"`  // ระยะเวลา (minutes)
    Path        [][2]float64 `json:"path"`         // เส้นทาง [lat, lon][]
}

// NearbyRequest คำขอค้นหาจุดใกล้เคียง
// NearbyRequest for finding nearby points
type NearbyRequest struct {
    Lat      float64 `form:"lat" binding:"required"`
    Lon      float64 `form:"lon" binding:"required"`
    Radius   float64 `form:"radius" default:"5"`      // กิโลเมตร
    MaxItems int     `form:"max_items" default:"10"`
}

var (
    ecoPoints []EcoPoint
    port      = getEnv("PORT", "8080")
)

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

// initData สร้างข้อมูลตัวอย่าง
// initData creates sample data
func initData() {
    // จุดรวบรวมในเขตกรุงเทพฯ
    // Collection points in Bangkok
    ecoPoints = []EcoPoint{
        {
            ID:          1,
            Name:        "ศูนย์รีไซเคิลสุขุมวิท (Sukhumvit Recycling Center)",
            Lat:         13.736717,
            Lon:         100.523186,
            Address:     "สุขุมวิท 23, คลองเตยเหนือ, วัฒนา, กรุงเทพฯ",
            Type:        "plastic",
            Capacity:    500,
            CurrentLoad: 320,
            Status:      "active",
            LastUpdate:  time.Now().Format(time.RFC3339),
        },
        {
            ID:          2,
            Name:        "ธนบุรีรีไซเคิล (Thonburi Recycling)",
            Lat:         13.726417,
            Lon:         100.512186,
            Address:     "ถนนราชดำเนิน, พระนคร, กรุงเทพฯ",
            Type:        "glass",
            Capacity:    300,
            CurrentLoad: 180,
            Status:      "active",
            LastUpdate:  time.Now().Format(time.RFC3339),
        },
        {
            ID:          3,
            Name:        "ลาดพร้าวรีไซเคิลฮับ (Lat Phrao Recycling Hub)",
            Lat:         13.803817,
            Lon:         100.609022,
            Address:     "ลาดพร้าว 101, จอมพล, จตุจักร, กรุงเทพฯ",
            Type:        "paper",
            Capacity:    400,
            CurrentLoad: 250,
            Status:      "active",
            LastUpdate:  time.Now().Format(time.RFC3339),
        },
        {
            ID:          4,
            Name:        "บางนาอีโว (Bangna E-Waste)",
            Lat:         13.669586,
            Lon:         100.624971,
            Address:     "บางนา-ตราด กม.3, บางนา, กรุงเทพฯ",
            Type:        "electronic",
            Capacity:    200,
            CurrentLoad: 195,
            Status:      "full",
            LastUpdate:  time.Now().Format(time.RFC3339),
        },
        {
            ID:          5,
            Name:        "รังสิตพลาสติก (Rangsit Plastic Center)",
            Lat:         13.987918,
            Lon:         100.620053,
            Address:     "รังสิต-นครนายก คลอง 5, ธัญบุรี, ปทุมธานี",
            Type:        "plastic",
            Capacity:    600,
            CurrentLoad: 420,
            Status:      "active",
            LastUpdate:  time.Now().Format(time.RFC3339),
        },
        {
            ID:          6,
            Name:        "ศูนย์รีไซเคิลพระราม 9 (Rama 9 Recycling)",
            Lat:         13.761852,
            Lon:         100.568247,
            Address:     "พระราม 9 ซอย 41, สวนหลวง, กรุงเทพฯ",
            Type:        "glass",
            Capacity:    250,
            CurrentLoad: 120,
            Status:      "active",
            LastUpdate:  time.Now().Format(time.RFC3339),
        },
        {
            ID:          7,
            Name:        "อ่อนนุชรีไซเคิล (On Nut Recycling)",
            Lat:         13.711898,
            Lon:         100.616921,
            Address:     "อ่อนนุช 46, ประเวศ, กรุงเทพฯ",
            Type:        "paper",
            Capacity:    350,
            CurrentLoad: 280,
            Status:      "maintenance",
            LastUpdate:  time.Now().Format(time.RFC3339),
        },
        {
            ID:          8,
            Name:        "บางแครีไซเคิล (Bang Khae Recycling)",
            Lat:         13.695318,
            Lon:         100.409342,
            Address:     "เพชรเกษม 81, บางแค, กรุงเทพฯ",
            Type:        "plastic",
            Capacity:    450,
            CurrentLoad: 310,
            Status:      "active",
            LastUpdate:  time.Now().Format(time.RFC3339),
        },
    }
}

// haversineDistance คำนวณระยะทางระหว่างจุด 2 จุด (km)
// haversineDistance calculates distance between two points (km)
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
    const R = 6371 // Earth radius in km
    
    dLat := (lat2 - lat1) * math.Pi / 180
    dLon := (lon2 - lon1) * math.Pi / 180
    
    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
        math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
            math.Sin(dLon/2)*math.Sin(dLon/2)
    
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    
    return R * c
}

// calculateRoute คำนวณเส้นทางระหว่างจุด (จำลอง)
// calculateRoute calculates route between points (simulation)
func calculateRoute(start, end EcoPoint) RouteResponse {
    distance := haversineDistance(start.Lat, start.Lon, end.Lat, end.Lon)
    
    // สร้างเส้นทางด้วยจุดเชื่อมต่อ (interpolation)
    // Create path with interpolated points
    steps := 10
    path := make([][2]float64, steps+1)
    
    for i := 0; i <= steps; i++ {
        t := float64(i) / float64(steps)
        lat := start.Lat + (end.Lat-start.Lat)*t
        lon := start.Lon + (end.Lon-start.Lon)*t
        path[i] = [2]float64{lat, lon}
    }
    
    return RouteResponse{
        StartPoint: start,
        EndPoint:   end,
        Distance:   distance,
        Duration:   int(distance * 2), // สมมติ 2 นาทีต่อกิโลเมตร
        Path:       path,
    }
}

// getPointsHandler ส่งคืนรายการจุดทั้งหมด
// getPointsHandler returns all eco points
func getPointsHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "count":  len(ecoPoints),
        "points": ecoPoints,
    })
}

// getPointByIDHandler ส่งคืนจุดตาม ID
// getPointByIDHandler returns point by ID
func getPointByIDHandler(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }
    
    for _, point := range ecoPoints {
        if point.ID == id {
            c.JSON(http.StatusOK, point)
            return
        }
    }
    
    c.JSON(http.StatusNotFound, gin.H{"error": "Point not found"})
}

// nearbyPointsHandler ค้นหาจุดใกล้เคียง
// nearbyPointsHandler finds nearby points
func nearbyPointsHandler(c *gin.Context) {
    var req NearbyRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if req.Radius == 0 {
        req.Radius = 5
    }
    if req.MaxItems == 0 {
        req.MaxItems = 10
    }
    
    // คำนวณระยะทางและเรียงลำดับ
    // Calculate distances and sort
    type PointWithDistance struct {
        Point    EcoPoint
        Distance float64
    }
    
    var nearby []PointWithDistance
    for _, point := range ecoPoints {
        dist := haversineDistance(req.Lat, req.Lon, point.Lat, point.Lon)
        if dist <= req.Radius {
            nearby = append(nearby, PointWithDistance{
                Point:    point,
                Distance: dist,
            })
        }
    }
    
    // เรียงตามระยะทาง
    // Sort by distance
    for i := 0; i < len(nearby)-1; i++ {
        for j := i + 1; j < len(nearby); j++ {
            if nearby[i].Distance > nearby[j].Distance {
                nearby[i], nearby[j] = nearby[j], nearby[i]
            }
        }
    }
    
    // จำกัดจำนวน
    // Limit count
    if len(nearby) > req.MaxItems {
        nearby = nearby[:req.MaxItems]
    }
    
    result := make([]gin.H, len(nearby))
    for i, np := range nearby {
        result[i] = gin.H{
            "point":    np.Point,
            "distance": np.Distance,
        }
    }
    
    c.JSON(http.StatusOK, gin.H{
        "count":     len(result),
        "nearby":    result,
        "center_lat": req.Lat,
        "center_lon": req.Lon,
        "radius_km": req.Radius,
    })
}

// routeToNearestHandler สร้างเส้นทางไปยังจุดที่ใกล้ที่สุด
// routeToNearestHandler creates route to nearest point
func routeToNearestHandler(c *gin.Context) {
    var req NearbyRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // หาจุดที่ใกล้ที่สุด
    // Find nearest point
    var nearest *EcoPoint
    var minDistance float64 = 1e9
    
    for i := range ecoPoints {
        dist := haversineDistance(req.Lat, req.Lon, ecoPoints[i].Lat, ecoPoints[i].Lon)
        if dist < minDistance && ecoPoints[i].Status == "active" {
            minDistance = dist
            nearest = &ecoPoints[i]
        }
    }
    
    if nearest == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "No active points found"})
        return
    }
    
    // สร้างจุดเริ่มต้นจำลอง
    // Create simulated start point
    startPoint := EcoPoint{
        ID:   0,
        Name: "Your Location",
        Lat:  req.Lat,
        Lon:  req.Lon,
    }
    
    route := calculateRoute(startPoint, *nearest)
    c.JSON(http.StatusOK, route)
}

// randomPointsHandler ส่งคืนจุดสุ่ม
// randomPointsHandler returns random points
func randomPointsHandler(c *gin.Context) {
    nStr := c.DefaultQuery("n", "10")
    n, err := strconv.Atoi(nStr)
    if err != nil || n > len(ecoPoints) {
        n = len(ecoPoints)
    }
    
    // สุ่มจุด
    // Shuffle points
    shuffled := make([]EcoPoint, len(ecoPoints))
    copy(shuffled, ecoPoints)
    
    rand.Seed(time.Now().UnixNano())
    for i := len(shuffled) - 1; i > 0; i-- {
        j := rand.Intn(i + 1)
        shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
    }
    
    c.JSON(http.StatusOK, gin.H{
        "count":  n,
        "points": shuffled[:n],
    })
}

// healthHandler ตรวจสอบสถานะเซิร์ฟเวอร์
// healthHandler checks server health
func healthHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":  "ok",
        "version": "1.0.0",
        "time":    time.Now().Unix(),
        "points_count": len(ecoPoints),
    })
}

// statsHandler ส่งคืนสถิติ
// statsHandler returns statistics
func statsHandler(c *gin.Context) {
    stats := make(map[string]interface{})
    
    // สถิติแยกตามประเภท
    // Statistics by type
    typeStats := make(map[string]int)
    statusStats := make(map[string]int)
    totalCapacity := 0
    totalLoad := 0
    
    for _, p := range ecoPoints {
        typeStats[p.Type]++
        statusStats[p.Status]++
        totalCapacity += p.Capacity
        totalLoad += p.CurrentLoad
    }
    
    stats["by_type"] = typeStats
    stats["by_status"] = statusStats
    stats["total_points"] = len(ecoPoints)
    stats["total_capacity_kg"] = totalCapacity
    stats["total_current_load_kg"] = totalLoad
    stats["utilization_percent"] = float64(totalLoad) / float64(totalCapacity) * 100
    
    c.JSON(http.StatusOK, stats)
}

func main() {
    // ตั้งค่า Gin
    // Setup Gin
    gin.SetMode(gin.ReleaseMode)
    r := gin.Default()
    
    // ตั้งค่า CORS (allow all for development)
    // Setup CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))
    
    // 初始化ข้อมูล
    // Initialize data
    initData()
    
    // Routes
    r.GET("/health", healthHandler)
    r.GET("/api/points", getPointsHandler)
    r.GET("/api/points/:id", getPointByIDHandler)
    r.GET("/api/points/random", randomPointsHandler)
    r.GET("/api/nearby", nearbyPointsHandler)
    r.GET("/api/route/nearest", routeToNearestHandler)
    r.GET("/api/stats", statsHandler)
    
    // Start server with graceful shutdown
    srv := &http.Server{
        Addr:    ":" + port,
        Handler: r,
    }
    
    go func() {
        log.Printf("EcoTracker Backend starting on port %s", port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()
    
    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }
    
    log.Println("Server exited")
}
```

#### โค้ด Frontend (React + TypeScript)

```tsx
// File: frontend/src/App.tsx - EcoTracker Frontend
// ไฟล์: frontend/src/App.tsx - หน้าแสดงแผนที่ EcoTracker

import React, { useState, useEffect, useRef } from 'react';
import Map, { Marker, Popup, NavigationControl, Source, Layer } from 'react-map-gl/maplibre';
import 'maplibre-gl/dist/maplibre-gl.css';
import axios from 'axios';

// Types
interface EcoPoint {
  id: number;
  name: string;
  lat: number;
  lon: number;
  address: string;
  type: string;
  capacity: number;
  current_load: number;
  status: string;
  last_update: string;
}

interface Route {
  start_point: EcoPoint;
  end_point: EcoPoint;
  distance_km: number;
  duration_min: number;
  path: [number, number][];
}

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

// EcoTracker Component
// คอมโพเนนต์หลักของ EcoTracker
const EcoTracker: React.FC = () => {
  const [points, setPoints] = useState<EcoPoint[]>([]);
  const [selectedPoint, setSelectedPoint] = useState<EcoPoint | null>(null);
  const [userLocation, setUserLocation] = useState<[number, number] | null>(null);
  const [route, setRoute] = useState<Route | null>(null);
  const [loading, setLoading] = useState(false);
  const [stats, setStats] = useState<any>(null);
  const mapRef = useRef<any>(null);

  // โหลดข้อมูลจุดรวบรวม
  // Load eco points data
  useEffect(() => {
    fetchPoints();
    fetchStats();
    
    // ขออนุญาตเข้าถึงตำแหน่ง
    // Request geolocation permission
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          setUserLocation([position.coords.longitude, position.coords.latitude]);
        },
        (error) => {
          console.error('Geolocation error:', error);
          // ตั้งค่าตำแหน่งเริ่มต้นเป็นกรุงเทพฯ
          // Set default location to Bangkok
          setUserLocation([100.523186, 13.736717]);
        }
      );
    }
  }, []);

  // ดึงข้อมูลจุดรวบรวม
  // Fetch eco points
  const fetchPoints = async () => {
    try {
      const response = await axios.get(`${API_URL}/api/points`);
      setPoints(response.data.points);
    } catch (error) {
      console.error('Error fetching points:', error);
    }
  };

  // ดึงข้อมูลสถิติ
  // Fetch statistics
  const fetchStats = async () => {
    try {
      const response = await axios.get(`${API_URL}/api/stats`);
      setStats(response.data);
    } catch (error) {
      console.error('Error fetching stats:', error);
    }
  };

  // ค้นหาจุดใกล้เคียง
  // Find nearby points
  const findNearbyPoints = async () => {
    if (!userLocation) return;
    
    setLoading(true);
    try {
      const response = await axios.get(`${API_URL}/api/nearby`, {
        params: {
          lat: userLocation[1],
          lon: userLocation[0],
          radius: 10,
          max_items: 10
        }
      });
      
      // แสดงผลใน console
      console.log('Nearby points:', response.data);
      
      // ซูมไปยังตำแหน่งผู้ใช้
      if (mapRef.current) {
        mapRef.current.flyTo({
          center: userLocation,
          zoom: 12,
          duration: 2000
        });
      }
    } catch (error) {
      console.error('Error finding nearby points:', error);
    } finally {
      setLoading(false);
    }
  };

  // สร้างเส้นทางไปยังจุดที่ใกล้ที่สุด
  // Create route to nearest point
  const findRouteToNearest = async () => {
    if (!userLocation) {
      alert('Please enable location access');
      return;
    }
    
    setLoading(true);
    try {
      const response = await axios.get(`${API_URL}/api/route/nearest`, {
        params: {
          lat: userLocation[1],
          lon: userLocation[0]
        }
      });
      setRoute(response.data);
      
      // ซูมเพื่อแสดงเส้นทางทั้งหมด
      if (mapRef.current && response.data.path.length > 0) {
        const bounds = response.data.path.reduce(
          (bounds: [[number, number], [number, number]], coord: [number, number]) => {
            return [
              [Math.min(bounds[0][0], coord[0]), Math.min(bounds[0][1], coord[1])],
              [Math.max(bounds[1][0], coord[0]), Math.max(bounds[1][1], coord[1])]
            ];
          },
          [[180, 90], [-180, -90]]
        );
        mapRef.current.fitBounds(bounds, { padding: 50, duration: 2000 });
      }
    } catch (error) {
      console.error('Error finding route:', error);
    } finally {
      setLoading(false);
    }
  };

  // รับสีตามประเภทขยะ
  // Get color by waste type
  const getMarkerColor = (type: string): string => {
    const colors: Record<string, string> = {
      plastic: '#4CAF50',
      glass: '#2196F3',
      paper: '#FF9800',
      electronic: '#9C27B0'
    };
    return colors[type] || '#757575';
  };

  // รับสถานะ
  // Get status badge
  const getStatusBadge = (status: string): JSX.Element => {
    const styles: Record<string, { color: string; bg: string; text: string }> = {
      active: { color: '#4CAF50', bg: '#E8F5E9', text: '🟢 เปิดให้บริการ' },
      full: { color: '#FF9800', bg: '#FFF3E0', text: '🟠 เต็มแล้ว' },
      maintenance: { color: '#F44336', bg: '#FFEBEE', text: '🔴 ปิดซ่อมบำรุง' }
    };
    const style = styles[status] || styles.active;
    return (
      <span style={{ background: style.bg, color: style.color, padding: '2px 8px', borderRadius: '12px', fontSize: '12px' }}>
        {style.text}
      </span>
    );
  };

  // แปลงเส้นทางเป็น GeoJSON
  // Convert route to GeoJSON
  const routeGeoJSON = route ? {
    type: 'Feature' as const,
    geometry: {
      type: 'LineString' as const,
      coordinates: route.path
    },
    properties: {}
  } : null;

  return (
    <div style={{ position: 'relative', height: '100vh', width: '100%' }}>
      {/* Control Panel */}
      {/* แผงควบคุม */}
      <div style={{
        position: 'absolute',
        top: 20,
        right: 20,
        zIndex: 1000,
        background: 'white',
        borderRadius: 8,
        boxShadow: '0 2px 10px rgba(0,0,0,0.1)',
        padding: 15,
        minWidth: 250,
        maxWidth: 350,
        maxHeight: '80vh',
        overflowY: 'auto'
      }}>
        <h3 style={{ margin: '0 0 10px 0' }}>🗺️ EcoTracker</h3>
        
        {/* Statistics */}
        {/* สถิติ */}
        {stats && (
          <div style={{ marginBottom: 15, padding: 10, background: '#f5f5f5', borderRadius: 5 }}>
            <div><strong>📊 สถิติรวม</strong></div>
            <div>📍 จุดรวบรวม: {stats.total_points} แห่ง</div>
            <div>♻️ กำลังการผลิตรวม: {stats.total_capacity_kg?.toLocaleString()} kg/วัน</div>
            <div>📦 ปริมาณปัจจุบัน: {stats.total_current_load_kg?.toLocaleString()} kg</div>
            <div>📈 อัตราการใช้: {stats.utilization_percent?.toFixed(1)}%</div>
          </div>
        )}
        
        {/* Action Buttons */}
        {/* ปุ่มดำเนินการ */}
        <div style={{ display: 'flex', gap: 10, marginBottom: 15 }}>
          <button
            onClick={findNearbyPoints}
            disabled={loading}
            style={{
              flex: 1,
              padding: '8px 12px',
              background: '#4CAF50',
              color: 'white',
              border: 'none',
              borderRadius: 5,
              cursor: loading ? 'not-allowed' : 'pointer'
            }}
          >
            {loading ? '⏳ กำลังค้นหา...' : '🔍 จุดใกล้ฉัน'}
          </button>
          <button
            onClick={findRouteToNearest}
            disabled={loading || !userLocation}
            style={{
              flex: 1,
              padding: '8px 12px',
              background: '#2196F3',
              color: 'white',
              border: 'none',
              borderRadius: 5,
              cursor: (loading || !userLocation) ? 'not-allowed' : 'pointer'
            }}
          >
            {loading ? '⏳ กำลังสร้าง...' : '🎯 เส้นทางไปจุดใกล้สุด'}
          </button>
        </div>
        
        {/* Route Info */}
        {/* ข้อมูลเส้นทาง */}
        {route && (
          <div style={{ marginBottom: 15, padding: 10, background: '#E3F2FD', borderRadius: 5 }}>
            <div><strong>🚗 ข้อมูลเส้นทาง</strong></div>
            <div>📍 ไปยัง: {route.end_point.name}</div>
            <div>📏 ระยะทาง: {route.distance_km.toFixed(2)} km</div>
            <div>⏱️ เวลาโดยประมาณ: {route.duration_min} นาที</div>
          </div>
        )}
        
        {/* Points List */}
        {/* รายการจุดรวบรวม */}
        <div>
          <strong>📍 จุดรวบรวมทั้งหมด ({points.length})</strong>
          <div style={{ marginTop: 10, maxHeight: 300, overflowY: 'auto' }}>
            {points.map(point => (
              <div
                key={point.id}
                onClick={() => setSelectedPoint(point)}
                style={{
                  padding: '8px',
                  marginBottom: 5,
                  background: selectedPoint?.id === point.id ? '#E3F2FD' : '#f9f9f9',
                  borderRadius: 5,
                  cursor: 'pointer',
                  border: '1px solid #eee'
                }}
              >
                <div><strong>{point.name}</strong></div>
                <div style={{ fontSize: '12px', color: '#666' }}>{point.address}</div>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: 5 }}>
                  {getStatusBadge(point.status)}
                  <span>📦 {point.current_load}/{point.capacity} kg</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
      
      {/* Map */}
      {/* แผนที่ */}
      <Map
        ref={mapRef}
        initialViewState={{
          longitude: 100.523186,
          latitude: 13.736717,
          zoom: 11
        }}
        style={{ width: '100%', height: '100%' }}
        mapStyle="https://basemaps.cartocdn.com/gl/positron-gl-style/style.json"
      >
        <NavigationControl position="top-left" />
        
        {/* User Location Marker */}
        {/* มาร์กเกอร์ตำแหน่งผู้ใช้ */}
        {userLocation && (
          <Marker longitude={userLocation[0]} latitude={userLocation[1]}>
            <div style={{
              width: 20,
              height: 20,
              background: '#E91E63',
              border: '2px solid white',
              borderRadius: '50%',
              boxShadow: '0 0 0 3px rgba(233,30,99,0.3)'
            }} />
          </Marker>
        )}
        
        {/* Eco Points Markers */}
        {/* มาร์กเกอร์จุดรวบรวม */}
        {points.map(point => (
          <Marker
            key={point.id}
            longitude={point.lon}
            latitude={point.lat}
            onClick={() => setSelectedPoint(point)}
          >
            <div
              style={{
                width: 32,
                height: 32,
                background: getMarkerColor(point.type),
                border: '2px solid white',
                borderRadius: '50%',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                cursor: 'pointer',
                boxShadow: '0 2px 5px rgba(0,0,0,0.2)'
              }}
            >
              ♻️
            </div>
          </Marker>
        ))}
        
        {/* Route Line */}
        {/* เส้นทาง */}
        {routeGeoJSON && (
          <Source id="route" type="geojson" data={routeGeoJSON}>
            <Layer
              id="route-line"
              type="line"
              paint={{
                'line-color': '#2196F3',
                'line-width': 4,
                'line-dasharray': [5, 5]
              }}
            />
          </Source>
        )}
        
        {/* Selected Point Popup */}
        {/* ป๊อปอัปจุดที่เลือก */}
        {selectedPoint && (
          <Popup
            longitude={selectedPoint.lon}
            latitude={selectedPoint.lat}
            onClose={() => setSelectedPoint(null)}
            closeButton={true}
            closeOnClick={false}
          >
            <div style={{ minWidth: 200 }}>
              <h4 style={{ margin: '0 0 5px 0' }}>{selectedPoint.name}</h4>
              <p style={{ margin: '5px 0', fontSize: '12px', color: '#666' }}>{selectedPoint.address}</p>
              <div><strong>ประเภท:</strong> {selectedPoint.type}</div>
              <div><strong>กำลังการผลิต:</strong> {selectedPoint.capacity} kg/วัน</div>
              <div><strong>ปริมาณปัจจุบัน:</strong> {selectedPoint.current_load} kg</div>
              <div>{getStatusBadge(selectedPoint.status)}</div>
              <div style={{ fontSize: '11px', color: '#999', marginTop: 5 }}>
                อัปเดตล่าสุด: {new Date(selectedPoint.last_update).toLocaleString()}
              </div>
            </div>
          </Popup>
        )}
      </Map>
    </div>
  );
};

export default EcoTracker;
```

#### การติดตั้งและรัน (Installation & Run)

```bash
# 1. Backend Setup
cd backend
go mod init ecotracker
go get github.com/gin-gonic/gin
go get github.com/gin-contrib/cors
go run main.go

# 2. Frontend Setup (อีก terminal)
cd frontend
npm create vite@latest . -- --template react-ts
npm install
npm install react-map-gl maplibre-gl axios
npm run dev

# 3. เปิด Browser: http://localhost:5173
```

---

### Case 4: Device Hub - Local IoT Data Infrastructure

#### โครงสร้างการทำงาน (Architecture)

```
รูปที่ 5: Device Hub Architecture

    [IoT Devices]         [Device Hub]           [Storage]
         │                     │                      │
         │  HTTP API           │                      │
         │  Register Device    │                      │
         │────────────────────▶│                      │
         │                     │                      │
         │  mDNS Discovery     │                      │
         │◀───────────────────▶│                      │
         │                     │                      │
         │  Telemetry Data     │                      │
         │────────────────────▶│  Store to InfluxDB   │
         │                     │─────────────────────▶│
         │                     │                      │
         │                     │  Monitor Status      │
         │                     │─────────────────────▶│
```

#### โค้ด Device Hub (Go)

```go
// File: main.go - Device Hub for IoT Data Collection
// ไฟล์: main.go - ศูนย์กลางรวบรวมข้อมูล IoT

package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net"
    "net/http"
    "sync"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/gorilla/websocket"
    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
    "github.com/miekg/dns"
)

// Device ข้อมูลอุปกรณ์ IoT
// IoT Device information
type Device struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Type        string            `json:"type"`        // sensor, actuator, gateway
    IP          string            `json:"ip"`
    Port        int               `json:"port"`
    Status      string            `json:"status"`      // online, offline, inactive
    LastSeen    time.Time         `json:"last_seen"`
    Metadata    map[string]string `json:"metadata"`
    Telemetry   map[string]interface{} `json:"telemetry"`
}

// TelemetryData ข้อมูล telemetry
// Telemetry data structure
type TelemetryData struct {
    DeviceID  string                 `json:"device_id"`
    Timestamp time.Time              `json:"timestamp"`
    Metrics   map[string]interface{} `json:"metrics"`
    Tags      map[string]string      `json:"tags"`
}

type DeviceHub struct {
    devices      map[string]*Device
    mu           sync.RWMutex
    influxClient influxdb2.Client
    httpServer   *http.Server
    upgrader     websocket.Upgrader
    wsClients    map[*websocket.Conn]bool
    wsMu         sync.RWMutex
}

// NewDeviceHub สร้าง Device Hub ใหม่
// NewDeviceHub creates a new Device Hub
func NewDeviceHub(influxURL, influxToken, influxOrg, influxBucket string) *DeviceHub {
    dh := &DeviceHub{
        devices:   make(map[string]*Device),
        wsClients: make(map[*websocket.Conn]bool),
        upgrader: websocket.Upgrader{
            CheckOrigin: func(r *http.Request) bool { return true },
        },
    }
    
    // ตั้งค่า InfluxDB client
    // Setup InfluxDB client
    if influxURL != "" {
        dh.influxClient = influxdb2.NewClient(influxURL, influxToken)
    }
    
    return dh
}

// registerDeviceHandler ลงทะเบียนอุปกรณ์ใหม่
// registerDeviceHandler registers a new device
func (dh *DeviceHub) registerDeviceHandler(w http.ResponseWriter, r *http.Request) {
    var device Device
    if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    device.LastSeen = time.Now()
    device.Status = "online"
    
    dh.mu.Lock()
    dh.devices[device.ID] = &device
    dh.mu.Unlock()
    
    log.Printf("Device registered: %s (%s)", device.Name, device.ID)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(device)
}

// telemetryHandler รับข้อมูล telemetry
// telemetryHandler receives telemetry data
func (dh *DeviceHub) telemetryHandler(w http.ResponseWriter, r *http.Request) {
    var telemetry TelemetryData
    if err := json.NewDecoder(r.Body).Decode(&telemetry); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    telemetry.Timestamp = time.Now()
    
    // อัปเดตสถานะอุปกรณ์
    // Update device status
    dh.mu.Lock()
    if device, exists := dh.devices[telemetry.DeviceID]; exists {
        device.LastSeen = time.Now()
        device.Status = "online"
        device.Telemetry = telemetry.Metrics
    }
    dh.mu.Unlock()
    
    // บันทึกลง InfluxDB
    // Save to InfluxDB
    if dh.influxClient != nil {
        go dh.saveToInfluxDB(telemetry)
    }
    
    // Broadcast ผ่าน WebSocket
    // Broadcast via WebSocket
    dh.broadcastTelemetry(telemetry)
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// saveToInfluxDB บันทึก telemetry ลง InfluxDB
// saveToInfluxDB saves telemetry to InfluxDB
func (dh *DeviceHub) saveToInfluxDB(telemetry TelemetryData) {
    writeAPI := dh.influxClient.WriteAPI("my-org", "my-bucket")
    
    point := influxdb2.NewPoint(
        "telemetry",
        map[string]string{
            "device_id": telemetry.DeviceID,
        },
        telemetry.Metrics,
        telemetry.Timestamp,
    )
    
    writeAPI.WritePoint(point)
    writeAPI.Flush()
}

// broadcastTelemetry ส่งข้อมูล telemetry ผ่าน WebSocket
// broadcastTelemetry sends telemetry via WebSocket
func (dh *DeviceHub) broadcastTelemetry(telemetry TelemetryData) {
    dh.wsMu.RLock()
    defer dh.wsMu.RUnlock()
    
    data, _ := json.Marshal(telemetry)
    for client := range dh.wsClients {
        if err := client.WriteMessage(websocket.TextMessage, data); err != nil {
            client.Close()
            delete(dh.wsClients, client)
        }
    }
}

// wsHandler จัดการ WebSocket connections
// wsHandler handles WebSocket connections
func (dh *DeviceHub) wsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := dh.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket upgrade error: %v", err)
        return
    }
    
    dh.wsMu.Lock()
    dh.wsClients[conn] = true
    dh.wsMu.Unlock()
    
    // ส่งรายการอุปกรณ์ทั้งหมดเมื่อเชื่อมต่อ
    // Send all devices list on connection
    dh.mu.RLock()
    devicesList := make([]*Device, 0, len(dh.devices))
    for _, device := range dh.devices {
        devicesList = append(devicesList, device)
    }
    dh.mu.RUnlock()
    
    conn.WriteJSON(map[string]interface{}{
        "type":    "devices",
        "devices": devicesList,
    })
}

// devicesHandler ส่งคืนรายการอุปกรณ์ทั้งหมด
// devicesHandler returns all devices
func (dh *DeviceHub) devicesHandler(w http.ResponseWriter, r *http.Request) {
    dh.mu.RLock()
    defer dh.mu.RUnlock()
    
    devices := make([]*Device, 0, len(dh.devices))
    for _, device := range dh.devices {
        devices = append(devices, device)
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(devices)
}

// mDNSServer เริ่ม mDNS server สำหรับการค้นพบอุปกรณ์
// mDNSServer starts mDNS server for device discovery
func (dh *DeviceHub) mDNSServer() {
    // mDNS query handler
    handler := func(w dns.ResponseWriter, r *dns.Msg) {
        m := new(dns.Msg)
        m.SetReply(r)
        
        for _, question := range r.Question {
            if question.Qtype == dns.TypeA {
                // ตอบกลับด้วย IP ของ hub
                // Respond with hub IP
                rr, _ := dns.NewRR(fmt.Sprintf("%s A %s", question.Name, getLocalIP()))
                m.Answer = append(m.Answer, rr)
            }
        }
        
        w.WriteMsg(m)
    }
    
    server := &dns.Server{Addr: ":5353", Net: "udp", Handler: dns.HandlerFunc(handler)}
    
    go func() {
        if err := server.ListenAndServe(); err != nil {
            log.Printf("mDNS server error: %v", err)
        }
    }()
    
    log.Println("mDNS server started on port 5353")
}

// getLocalIP ดึง IP ในเครื่อง
// getLocalIP gets local IP address
func getLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return "127.0.0.1"
    }
    
    for _, addr := range addrs {
        if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    
    return "127.0.0.1"
}

// monitorDevices ตรวจสอบสถานะอุปกรณ์เป็นระยะ
// monitorDevices periodically checks device status
func (dh *DeviceHub) monitorDevices() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        dh.mu.Lock()
        now := time.Now()
        for id, device := range dh.devices {
            // ถ้าไม่มีการอัปเดตเกิน 2 นาที ให้เปลี่ยนสถานะเป็น offline
            // If no update for 2 minutes, mark as offline
            if now.Sub(device.LastSeen) > 2*time.Minute {
                device.Status = "offline"
                log.Printf("Device %s is offline", id)
            }
        }
        dh.mu.Unlock()
    }
}

// Run เริ่มทำงาน Device Hub
// Run starts the Device Hub
func (dh *DeviceHub) Run(port string) {
    // เริ่ม mDNS server
    go dh.mDNSServer()
    
    // เริ่ม device monitor
    go dh.monitorDevices()
    
    // ตั้งค่า routes
    r := mux.NewRouter()
    r.HandleFunc("/api/devices/register", dh.registerDeviceHandler).Methods("POST")
    r.HandleFunc("/api/devices", dh.devicesHandler).Methods("GET")
    r.HandleFunc("/api/telemetry", dh.telemetryHandler).Methods("POST")
    r.HandleFunc("/ws", dh.wsHandler)
    
    // Static files
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
    
    dh.httpServer = &http.Server{
        Addr:    ":" + port,
        Handler: r,
    }
    
    log.Printf("Device Hub starting on port %s", port)
    log.Printf("mDNS service: _device-hub._tcp.local")
    
    if err := dh.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatal(err)
    }
}

func main() {
    hub := NewDeviceHub(
        "http://localhost:9087",
        "my-token",
        "my-org",
        "my-bucket",
    )
    
    hub.Run("8080")
}
```

---

### Case 5: Multi-Source Fusion Positioning - GPS/WiFi/UWB/IMU

#### โครงสร้างการทำงาน (Architecture)

```
รูปที่ 6: Multi-Source Fusion Positioning System

    [GPS]     [WiFi]     [UWB]      [IMU]
      │          │          │          │
      │  1Hz     │  2Hz     │ 100Hz    │ 200Hz
      │          │          │          │
      ▼          ▼          ▼          ▼
  ┌─────────────────────────────────────────┐
  │         Time Synchronization            │
  │         (UTC Nanosecond)                │
  └─────────────────────────────────────────┘
                      │
                      ▼
  ┌─────────────────────────────────────────┐
  │         Data Fusion Engine              │
  │  • Kalman Filter (GPS/WiFi)            │
  │  • Particle Filter (UWB/IMU)           │
  │  • Coordinate Transformation            │
  └─────────────────────────────────────────┘
                      │
                      ▼
  ┌─────────────────────────────────────────┐
  │         Position Output (50Hz)          │
  │  • x, y, z coordinates                  │
  │  • Confidence score                     │
  │  • Velocity & heading                   │
  └─────────────────────────────────────────┘
```

#### โค้ด Multi-Source Fusion (Go)

```go
// File: fusion.go - Multi-Source Positioning Engine
// ไฟล์: fusion.go - เอ็นจิ้นคำนวณตำแหน่งจากหลายแหล่งข้อมูล

package main

import (
    "fmt"
    "log"
    "math"
    "sync"
    "time"
)

// SensorEvent อินเทอร์เฟซสำหรับข้อมูลจากเซนเซอร์
// SensorEvent interface for sensor data
type SensorEvent interface {
    Timestamp() time.Time
    Position() (x, y, z float64)
    Confidence() float64
    Source() string
}

// GPSData ข้อมูลจาก GPS
// GPS data structure
type GPSData struct {
    Time     time.Time
    Lat      float64   // ละติจูด / Latitude
    Lon      float64   // ลองจิจูด / Longitude
    Alt      float64   // ความสูง / Altitude
    Speed    float64   // ความเร็ว / Speed
    Heading  float64   // ทิศทาง / Heading
    HDOP     float64   // Horizontal Dilution of Precision
}

func (g GPSData) Timestamp() time.Time { return g.Time }
func (g GPSData) Position() (x, y, z float64) {
    // แปลง GPS (lat, lon) เป็น Cartesian (x, y)
    // Convert GPS to Cartesian coordinates
    const R = 6371000 // Earth radius in meters
    x = R * math.Cos(g.Lat*math.Pi/180) * math.Cos(g.Lon*math.Pi/180)
    y = R * math.Cos(g.Lat*math.Pi/180) * math.Sin(g.Lon*math.Pi/180)
    z = g.Alt
    return
}
func (g GPSData) Confidence() float64 {
    // ความมั่นใจขึ้นอยู่กับ HDOP (ค่า越小ยิ่งดี)
    // Confidence based on HDOP (lower is better)
    if g.HDOP <= 1 {
        return 0.95
    } else if g.HDOP <= 2 {
        return 0.85
    } else if g.HDOP <= 5 {
        return 0.70
    }
    return 0.50
}
func (g GPSData) Source() string { return "gps" }

// WiFiData ข้อมูลจาก WiFi fingerprint
// WiFi data structure
type WiFiData struct {
    Time      time.Time
    BSSID     string
    RSSI      int
    Location  struct{ X, Y, Z float64 }
}

func (w WiFiData) Timestamp() time.Time { return w.Time }
func (w WiFiData) Position() (x, y, z float64) { return w.Location.X, w.Location.Y, w.Location.Z }
func (w WiFiData) Confidence() float64 {
    // ความมั่นใจขึ้นอยู่กับ RSSI (ค่าใกล้ 0 ยิ่งดี)
    // Confidence based on RSSI (closer to 0 is better)
    rssiNorm := float64(w.RSSI+100) / 100
    if rssiNorm > 0.8 {
        return 0.70
    }
    return 0.50
}
func (w WiFiData) Source() string { return "wifi" }

// UWBData ข้อมูลจาก UWB (Ultra-Wideband)
// UWB data structure
type UWBData struct {
    Time      time.Time
    AnchorID  string
    Distance  float64   // ระยะทางจาก anchor / Distance from anchor
    AnchorPos struct{ X, Y, Z float64 }
}

func (u UWBData) Timestamp() time.Time { return u.Time }
func (u UWBData) Position() (x, y, z float64) {
    // trilateration จะทำใน fusion engine
    // Trilateration will be done in fusion engine
    return u.AnchorPos.X, u.AnchorPos.Y, u.AnchorPos.Z
}
func (u UWBData) Confidence() float64 { return 0.95 }
func (u UWBData) Source() string { return "uwb" }

// IMUData ข้อมูลจาก IMU (Inertial Measurement Unit)
// IMU data structure
type IMUData struct {
    Time      time.Time
    AccX, AccY, AccZ float64   // Acceleration (m/s²)
    GyroX, GyroY, GyroZ float64 // Angular velocity (rad/s)
}

func (i IMUData) Timestamp() time.Time { return i.Time }
func (i IMUData) Position() (x, y, z float64) { return 0, 0, 0 } // จะคำนวณจากการอินทิเกรต
func (i IMUData) Confidence() float64 { return 0.60 }
func (i IMUData) Source() string { return "imu" }

// KalmanFilter ฟิลเตอร์คาลมานสำหรับ GPS/WiFi
// Kalman filter for GPS/WiFi
type KalmanFilter struct {
    // State vector [x, y, vx, vy]
    x [4]float64  // State
    P [4][4]float64 // Covariance matrix
    F [4][4]float64 // State transition
    H [2][4]float64 // Observation matrix
    R [2][2]float64 // Measurement noise
    Q [4][4]float64 // Process noise
}

// NewKalmanFilter สร้าง Kalman filter ใหม่
// NewKalmanFilter creates a new Kalman filter
func NewKalmanFilter(dt float64) *KalmanFilter {
    kf := &KalmanFilter{}
    
    // State transition matrix (constant velocity model)
    kf.F = [4][4]float64{
        {1, 0, dt, 0},
        {0, 1, 0, dt},
        {0, 0, 1, 0},
        {0, 0, 0, 1},
    }
    
    // Observation matrix (we observe position)
    kf.H = [2][4]float64{
        {1, 0, 0, 0},
        {0, 1, 0, 0},
    }
    
    // Measurement noise (GPS/WiFi error)
    kf.R = [2][2]float64{
        {25, 0},
        {0, 25},
    }
    
    // Process noise (model uncertainty)
    kf.Q = [4][4]float64{
        {0.1, 0, 0, 0},
        {0, 0.1, 0, 0},
        {0, 0, 0.5, 0},
        {0, 0, 0, 0.5},
    }
    
    // Initial covariance
    kf.P = [4][4]float64{
        {100, 0, 0, 0},
        {0, 100, 0, 0},
        {0, 0, 10, 0},
        {0, 0, 0, 10},
    }
    
    return kf
}

// Predict ทำนายขั้นตอนถัดไป
// Predict next state
func (kf *KalmanFilter) Predict() {
    // x = F * x
    var newX [4]float64
    for i := 0; i < 4; i++ {
        for j := 0; j < 4; j++ {
            newX[i] += kf.F[i][j] * kf.x[j]
        }
    }
    kf.x = newX
    
    // P = F * P * F^T + Q
    var newP [4][4]float64
    for i := 0; i < 4; i++ {
        for j := 0; j < 4; j++ {
            for k := 0; k < 4; k++ {
                for l := 0; l < 4; l++ {
                    newP[i][j] += kf.F[i][k] * kf.P[k][l] * kf.F[j][l]
                }
            }
            newP[i][j] += kf.Q[i][j]
        }
    }
    kf.P = newP
}

// Update ปรับปรุงค่าด้วยการวัด
// Update with measurement
func (kf *KalmanFilter) Update(zx, zy float64) {
    // y = z - H*x
    yx := zx - (kf.H[0][0]*kf.x[0] + kf.H[0][1]*kf.x[1])
    yy := zy - (kf.H[1][0]*kf.x[0] + kf.H[1][1]*kf.x[1])
    
    // S = H*P*H^T + R
    S := [2][2]float64{
        {kf.P[0][0] + kf.R[0][0], kf.P[0][1]},
        {kf.P[1][0], kf.P[1][1] + kf.R[1][1]},
    }
    
    // Kalman gain K = P*H^T * S^-1
    det := S[0][0]*S[1][1] - S[0][1]*S[1][0]
    if det == 0 {
        return
    }
    
    invS00 := S[1][1] / det
    invS01 := -S[0][1] / det
    invS10 := -S[1][0] / det
    invS11 := S[0][0] / det
    
    K := [4][2]float64{
        {(kf.P[0][0]*invS00 + kf.P[0][1]*invS10),
         (kf.P[0][0]*invS01 + kf.P[0][1]*invS11)},
        {(kf.P[1][0]*invS00 + kf.P[1][1]*invS10),
         (kf.P[1][0]*invS01 + kf.P[1][1]*invS11)},
        {(kf.P[2][0]*invS00 + kf.P[2][1]*invS10),
         (kf.P[2][0]*invS01 + kf.P[2][1]*invS11)},
        {(kf.P[3][0]*invS00 + kf.P[3][1]*invS10),
         (kf.P[3][0]*invS01 + kf.P[3][1]*invS11)},
    }
    
    // x = x + K*y
    kf.x[0] += K[0][0]*yx + K[0][1]*yy
    kf.x[1] += K[1][0]*yx + K[1][1]*yy
    kf.x[2] += K[2][0]*yx + K[2][1]*yy
    kf.x[3] += K[3][0]*yx + K[3][1]*yy
    
    // P = (I - K*H) * P
    for i := 0; i < 4; i++ {
        for j := 0; j < 4; j++ {
            kh := K[i][0]*kf.H[0][j] + K[i][1]*kf.H[1][j]
            kf.P[i][j] -= kh
        }
    }
}

// GetPosition ส่งคืนตำแหน่งปัจจุบัน
// GetPosition returns current position
func (kf *KalmanFilter) GetPosition() (x, y float64) {
    return kf.x[0], kf.x[1]
}

// FusionEngine เอ็นจิ้นรวมข้อมูล
// Fusion engine
type FusionEngine struct {
    kalman    *KalmanFilter
    particles []Particle
    lastTime  time.Time
    position  struct{ X, Y, Z float64 }
    mu        sync.RWMutex
    
    // Input channels
    gpsChan   chan GPSData
    wifiChan  chan WiFiData
    uwbChan   chan UWBData
    imuChan   chan IMUData
    
    // Output
    outputChan chan PositionEstimate
}

// Particle สำหรับ particle filter (UWB/IMU)
// Particle for particle filter
type Particle struct {
    X, Y, Z float64
    Weight  float64
}

// PositionEstimate ผลลัพธ์ตำแหน่ง
// Position estimate result
type PositionEstimate struct {
    Timestamp  time.Time
    X, Y, Z    float64
    Confidence float64
    Source     string
    Velocity   float64
    Heading    float64
}

// NewFusionEngine สร้าง fusion engine ใหม่
// NewFusionEngine creates a new fusion engine
func NewFusionEngine() *FusionEngine {
    fe := &FusionEngine{
        kalman:     NewKalmanFilter(0.02), // 50Hz, dt=0.02s
        gpsChan:    make(chan GPSData, 100),
        wifiChan:   make(chan WiFiData, 100),
        uwbChan:    make(chan UWBData, 100),
        imuChan:    make(chan IMUData, 100),
        outputChan: make(chan PositionEstimate, 100),
        particles:  make([]Particle, 1000),
    }
    
    // Initialize particles
    for i := range fe.particles {
        fe.particles[i] = Particle{
            X:      0,
            Y:      0,
            Z:      0,
            Weight: 1.0 / float64(len(fe.particles)),
        }
    }
    
    return fe
}

// ProcessGPS ประมวลผลข้อมูล GPS
// Process GPS data
func (fe *FusionEngine) ProcessGPS(data GPSData) {
    fe.mu.Lock()
    defer fe.mu.Unlock()
    
    // คำนวณตำแหน่งจาก GPS
    x, y, z := data.Position()
    
    // Predict แล้ว Update
    fe.kalman.Predict()
    fe.kalman.Update(x, y)
    
    // อัปเดตตำแหน่ง
    kx, ky := fe.kalman.GetPosition()
    fe.position.X = kx
    fe.position.Y = ky
    fe.position.Z = z
    
    // ส่ง output
    fe.outputChan <- PositionEstimate{
        Timestamp:  data.Time,
        X:          kx,
        Y:          ky,
        Z:          z,
        Confidence: data.Confidence(),
        Source:     "gps+kf",
        Velocity:   data.Speed,
        Heading:    data.Heading,
    }
}

// ProcessIMU ประมวลผลข้อมูล IMU
// Process IMU data
func (fe *FusionEngine) ProcessIMU(data IMUData) {
    fe.mu.Lock()
    defer fe.mu.Unlock()
    
    if fe.lastTime.IsZero() {
        fe.lastTime = data.Time
        return
    }
    
    dt := data.Time.Sub(fe.lastTime).Seconds()
    if dt <= 0 {
        return
    }
    
    // Dead reckoning with IMU
    // อัปเดตตำแหน่งด้วยการอินทิเกรตความเร่ง
    vx := fe.kalman.x[2]
    vy := fe.kalman.x[3]
    
    // อัปเดตความเร็ว
    fe.kalman.x[2] += data.AccX * dt
    fe.kalman.x[3] += data.AccY * dt
    
    // Dead reckoning position
    fe.position.X += vx * dt + 0.5*data.AccX*dt*dt
    fe.position.Y += vy * dt + 0.5*data.AccY*dt*dt
    
    // Dead reckoning for particles
    for i := range fe.particles {
        fe.particles[i].X += vx*dt + 0.5*data.AccX*dt*dt
        fe.particles[i].Y += vy*dt + 0.5*data.AccY*dt*dt
    }
    
    fe.lastTime = data.Time
}

// ProcessUWB ประมวลผลข้อมูล UWB ด้วย trilateration
// Process UWB data with trilateration
func (fe *FusionEngine) ProcessUWB(data UWBData) {
    fe.mu.Lock()
    defer fe.mu.Unlock()
    
    // trilateration จะทำที่นี่
    // Trilateration will be done here
    // (ตัวอย่าง: คำนวณตำแหน่งจากระยะทางไปยัง anchor 3 จุด)
    
    // ปรับปรุงน้ำหนัก particles
    // Update particle weights
    totalWeight := 0.0
    for i := range fe.particles {
        // คำนวณระยะทางจาก particle ไปยัง anchor
        dx := fe.particles[i].X - data.AnchorPos.X
        dy := fe.particles[i].Y - data.AnchorPos.Y
        dz := fe.particles[i].Z - data.AnchorPos.Z
        dist := math.Sqrt(dx*dx + dy*dy + dz*dz)
        
        // อัปเดตน้ำหนักตามความแตกต่างของระยะทาง
        diff := math.Abs(dist - data.Distance)
        fe.particles[i].Weight *= math.Exp(-diff * diff / (2 * 0.1))
        totalWeight += fe.particles[i].Weight
    }
    
    // Normalize weights
    if totalWeight > 0 {
        for i := range fe.particles {
            fe.particles[i].Weight /= totalWeight
        }
    }
    
    // Resampling
    fe.resampleParticles()
}

// resampleParticles สุ่มเลือก particles ใหม่ตามน้ำหนัก
// resampleParticles resamples particles based on weights
func (fe *FusionEngine) resampleParticles() {
    newParticles := make([]Particle, len(fe.particles))
    
    // Cumulative sum of weights
    cumWeights := make([]float64, len(fe.particles))
    cumWeights[0] = fe.particles[0].Weight
    for i := 1; i < len(fe.particles); i++ {
        cumWeights[i] = cumWeights[i-1] + fe.particles[i].Weight
    }
    
    // Systematic resampling
    step := 1.0 / float64(len(fe.particles))
    u := math.Float64frombits(uint64(rand.Int63())) * step
    
    j := 0
    for i := 0; i < len(fe.particles); i++ {
        for u > cumWeights[j] {
            j++
        }
        newParticles[i] = fe.particles[j]
        u += step
    }
    
    // Reset weights
    for i := range newParticles {
        newParticles[i].Weight = 1.0 / float64(len(newParticles))
    }
    
    fe.particles = newParticles
}

// GetPositionEstimate ส่งคืนค่าประมาณตำแหน่งล่าสุด
// GetPositionEstimate returns latest position estimate
func (fe *FusionEngine) GetPositionEstimate() PositionEstimate {
    fe.mu.RLock()
    defer fe.mu.RUnlock()
    
    // คำนวณค่าเฉลี่ยของ particles
    var avgX, avgY, avgZ float64
    for _, p := range fe.particles {
        avgX += p.X * p.Weight
        avgY += p.Y * p.Weight
        avgZ += p.Z * p.Weight
    }
    
    return PositionEstimate{
        Timestamp:  time.Now(),
        X:          fe.position.X,
        Y:          fe.position.Y,
        Z:          fe.position.Z,
        Confidence: 0.85,
        Source:     "fusion",
    }
}

// Run เริ่ม fusion engine
// Run starts fusion engine
func (fe *FusionEngine) Run() {
    ticker := time.NewTicker(20 * time.Millisecond) // 50Hz output
    defer ticker.Stop()
    
    for {
        select {
        case gps := <-fe.gpsChan:
            fe.ProcessGPS(gps)
        case wifi := <-fe.wifiChan:
            // Process WiFi (similar to GPS)
            _ = wifi
        case uwb := <-fe.uwbChan:
            fe.ProcessUWB(uwb)
        case imu := <-fe.imuChan:
            fe.ProcessIMU(imu)
        case <-ticker.C:
            // Predict step for IMU dead reckoning
            estimate := fe.GetPositionEstimate()
            log.Printf("Position: (%.2f, %.2f, %.2f) conf=%.2f",
                estimate.X, estimate.Y, estimate.Z, estimate.Confidence)
        }
    }
}

func main() {
    engine := NewFusionEngine()
    go engine.Run()
    
    // Simulate data
    go func() {
        for {
            // Simulate GPS at 1Hz
            time.Sleep(1 * time.Second)
            engine.gpsChan <- GPSData{
                Time:  time.Now(),
                Lat:   13.736717,
                Lon:   100.523186,
                Alt:   10.5,
                Speed: 5.2,
                HDOP:  1.2,
            }
        }
    }()
    
    select {}
}
```

---

## เทมเพลตและตัวอย่างโค้ดที่รันได้จริง (Templates & Runnable Code)

### เทมเพลตที่ 1: โครงสร้างโปรเจค Go สำหรับ IoT (Project Structure)

```
my-iot-project/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point
├── internal/
│   ├── device/                     # Device management
│   │   ├── registry.go
│   │   └── types.go
│   ├── ingestion/                  # Data ingestion
│   │   ├── mqtt.go
│   │   ├── http.go
│   │   └── websocket.go
│   ├── processor/                  # Data processing
│   │   ├── filter.go
│   │   ├── aggregator.go
│   │   └── fusion.go
│   ├── storage/                    # Data storage
│   │   ├── influxdb.go
│   │   ├── timescaledb.go
│   │   └── redis.go
│   └── web/                        # Web interface
│       ├── handlers.go
│       ├── middleware.go
│       └── static/
├── pkg/
│   ├── logger/                     # Logging
│   ├── metrics/                    # Metrics collection
│   └── utils/                      # Utilities
├── configs/
│   └── config.yaml                 # Configuration
├── deployments/
│   ├── docker/
│   │   └── Dockerfile
│   └── kubernetes/
│       └── deployment.yaml
├── go.mod
├── go.sum
└── Makefile
```

### เทมเพลตที่ 2: Docker Compose สำหรับ IoT Stack

```yaml
# docker-compose.yml
# Docker Compose configuration for IoT stack

version: '3.8'

services:
  # MQTT Broker
  mosquitto:
    image: eclipse-mosquitto:latest
    container_name: mosquitto
    ports:
      - "1883:1883"
      - "9001:9001"
    volumes:
      - ./mosquitto/config:/mosquitto/config
      - ./mosquitto/data:/mosquitto/data
      - ./mosquitto/log:/mosquitto/log
    networks:
      - iot-network

  # Time Series Database (InfluxDB)
  influxdb:
    image: influxdb:2.7
    container_name: influxdb
    ports:
      - "9087:8086"
    environment:
      - DOCKER_INFLUXDB_INIT_MODE=setup
      - DOCKER_INFLUXDB_INIT_USERNAME=admin
      - DOCKER_INFLUXDB_INIT_PASSWORD=admin123
      - DOCKER_INFLUXDB_INIT_ORG=my-org
      - DOCKER_INFLUXDB_INIT_BUCKET=iot-data
    volumes:
      - ./influxdb/data:/var/lib/influxdb2
    networks:
      - iot-network

  # PostgreSQL with TimescaleDB
  timescaledb:
    image: timescale/timescaledb:latest-pg14
    container_name: timescaledb
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=iotuser
      - POSTGRES_PASSWORD=iotpass
      - POSTGRES_DB=iotdb
    volumes:
      - ./timescaledb/data:/var/lib/postgresql/data
    networks:
      - iot-network

  # Redis for caching
  redis:
    image: redis:7-alpine
    container_name: redis
    ports:
      - "6379:6379"
    volumes:
      - ./redis/data:/data
    networks:
      - iot-network

  # Grafana for visualization
  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_USER=admin
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - ./grafana/data:/var/lib/grafana
      - ./grafana/provisioning:/etc/grafana/provisioning
    networks:
      - iot-network
    depends_on:
      - influxdb
      - timescaledb

  # Go IoT Gateway
  iot-gateway:
    build: .
    container_name: iot-gateway
    ports:
      - "8080:8080"
      - "8081:8081"
    environment:
      - MQTT_BROKER=tcp://mosquitto:1883
      - INFLUXDB_URL=http://influxdb:9087
      - INFLUXDB_TOKEN=my-token
      - REDIS_ADDR=redis:6379
    volumes:
      - ./config.yaml:/app/config.yaml
      - ./data:/app/data
    networks:
      - iot-network
    depends_on:
      - mosquitto
      - influxdb
      - redis
    restart: unless-stopped

networks:
  iot-network:
    driver: bridge
```

---

## แบบฝึกหัดท้ายบท (Exercises)

### แบบฝึกหัดที่ 1: สร้าง MQTT Publisher สำหรับจำลองข้อมูล Sensor

**โจทย์:** เขียนโปรแกรม Go ที่ส่งข้อมูล Temperature และ Humidity ไปยัง MQTT broker ทุก 5 วินาที

<details>
<summary>เฉลย (Solution)</summary>

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "time"
    
    mqtt "github.com/eclipse/paho.mqtt.golang"
)

type SensorData struct {
    DeviceID    string    `json:"device_id"`
    Temperature float64   `json:"temperature"`
    Humidity    float64   `json:"humidity"`
    Timestamp   time.Time `json:"timestamp"`
}

func main() {
    opts := mqtt.NewClientOptions()
    opts.AddBroker("tcp://localhost:1883")
    opts.SetClientID("sensor-simulator")
    
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        log.Fatal(token.Error())
    }
    
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        data := SensorData{
            DeviceID:    "sensor-001",
            Temperature: 20 + rand.Float64()*15,
            Humidity:    40 + rand.Float64()*40,
            Timestamp:   time.Now(),
        }
        
        jsonData, _ := json.Marshal(data)
        token := client.Publish("sensors/data", 1, false, jsonData)
        token.Wait()
        
        fmt.Printf("Published: %.2f°C, %.2f%%\n", data.Temperature, data.Humidity)
    }
}
```

</details>

### แบบฝึกหัดที่ 2: สร้าง WebSocket Dashboard สำหรับแสดงข้อมูลแบบ Real-time

**โจทย์:** สร้าง WebSocket server ที่ broadcast ข้อมูลตำแหน่งไปยัง client และแสดงบนแผนที่

<details>
<summary>เฉลย (Solution)</summary>

```go
package main

import (
    "log"
    "net/http"
    "sync"
    "time"
    
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
    clients    map[*websocket.Conn]bool
    broadcast  chan []byte
    mu         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:   make(map[*websocket.Conn]bool),
        broadcast: make(chan []byte),
    }
}

func (h *Hub) Run() {
    for msg := range h.broadcast {
        h.mu.RLock()
        for client := range h.clients {
            err := client.WriteMessage(websocket.TextMessage, msg)
            if err != nil {
                client.Close()
                delete(h.clients, client)
            }
        }
        h.mu.RUnlock()
    }
}

func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Upgrade error: %v", err)
        return
    }
    
    h.mu.Lock()
    h.clients[conn] = true
    h.mu.Unlock()
    
    // Keep connection alive
    for {
        if _, _, err := conn.ReadMessage(); err != nil {
            h.mu.Lock()
            delete(h.clients, conn)
            h.mu.Unlock()
            break
        }
    }
}

func main() {
    hub := NewHub()
    go hub.Run()
    
    http.HandleFunc("/ws", hub.HandleWebSocket)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`<html><body><div id="map"></div><script>
            const ws = new WebSocket('ws://' + location.host + '/ws');
            ws.onmessage = (e) => {
                const data = JSON.parse(e.data);
                console.log('Location:', data);
            };
        </script></body></html>`))
    })
    
    // Simulate location updates
    go func() {
        for {
            time.Sleep(1 * time.Second)
            hub.broadcast <- []byte(`{"lat":13.736717,"lon":100.523186}`)
        }
    }()
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

</details>

### แบบฝึกหัดที่ 3: สร้าง Edge Gateway ที่กรองข้อมูล

**โจทย์:** สร้าง edge gateway ที่รับข้อมูล MQTT กรองค่า temperature ที่เกินช่วง 0-50°C แล้วส่งต่อไปยัง Kafka

<details>
<summary>เฉลย (Solution)</summary>

```go
package main

import (
    "encoding/json"
    "log"
    "os"
    "os/signal"
    
    mqtt "github.com/eclipse/paho.mqtt.golang"
    "github.com/confluentinc/confluent-kafka-go/kafka"
)

type SensorReading struct {
    DeviceID    string  `json:"device_id"`
    Temperature float64 `json:"temperature"`
}

func main() {
    // MQTT subscriber
    mqttOpts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
    mqttClient := mqtt.NewClient(mqttOpts)
    mqttClient.Connect()
    
    // Kafka producer
    kafkaProd, _ := kafka.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers": "localhost:9092",
    })
    
    // Subscribe and filter
    mqttClient.Subscribe("sensors/+/data", 1, func(client mqtt.Client, msg mqtt.Message) {
        var reading SensorReading
        if err := json.Unmarshal(msg.Payload(), &reading); err != nil {
            return
        }
        
        // Filter: only pass temperature between 0-50°C
        if reading.Temperature >= 0 && reading.Temperature <= 50 {
            kafkaProd.Produce(&kafka.Message{
                TopicPartition: kafka.TopicPartition{Topic: stringPtr("filtered-data")},
                Value:          msg.Payload(),
            }, nil)
            log.Printf("Forwarded: %s - %.2f°C", reading.DeviceID, reading.Temperature)
        } else {
            log.Printf("Filtered out: %s - %.2f°C", reading.DeviceID, reading.Temperature)
        }
    })
    
    // Wait for interrupt
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, os.Interrupt)
    <-sig
}

func stringPtr(s string) *string { return &s }
```

</details>

### แบบฝึกหัดที่ 4: สร้าง Time Series Aggregation

**โจทย์:** เขียนฟังก์ชันที่รวมข้อมูล sensor ทุก 1 นาที (average, min, max)

<details>
<summary>เฉลย (Solution)</summary>

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Aggregator struct {
    windows map[string]*TimeWindow
    mu      sync.RWMutex
}

type TimeWindow struct {
    DeviceID    string
    Readings    []float64
    StartTime   time.Time
    EndTime     time.Time
    mu          sync.Mutex
}

type AggregatedResult struct {
    DeviceID  string
    Avg       float64
    Min       float64
    Max       float64
    Count     int
    StartTime time.Time
    EndTime   time.Time
}

func NewAggregator() *Aggregator {
    a := &Aggregator{windows: make(map[string]*TimeWindow)}
    go a.cleanup()
    return a
}

func (a *Aggregator) AddReading(deviceID string, value float64) {
    a.mu.Lock()
    defer a.mu.Unlock()
    
    now := time.Now()
    window, exists := a.windows[deviceID]
    
    if !exists || now.After(window.EndTime) {
        // Create new window
        window = &TimeWindow{
            DeviceID:  deviceID,
            Readings:  []float64{},
            StartTime: now.Truncate(time.Minute),
            EndTime:   now.Truncate(time.Minute).Add(time.Minute),
        }
        a.windows[deviceID] = window
    }
    
    window.mu.Lock()
    window.Readings = append(window.Readings, value)
    window.mu.Unlock()
}

func (a *Aggregator) cleanup() {
    ticker := time.NewTicker(30 * time.Second)
    for range ticker.C {
        a.mu.Lock()
        now := time.Now()
        for id, window := range a.windows {
            if now.After(window.EndTime.Add(time.Minute)) {
                // Process and remove
                result := a.processWindow(window)
                if result != nil {
                    fmt.Printf("Aggregated %s: avg=%.2f, min=%.2f, max=%.2f, count=%d\n",
                        result.DeviceID, result.Avg, result.Min, result.Max, result.Count)
                }
                delete(a.windows, id)
            }
        }
        a.mu.Unlock()
    }
}

func (a *Aggregator) processWindow(window *TimeWindow) *AggregatedResult {
    window.mu.Lock()
    defer window.mu.Unlock()
    
    if len(window.Readings) == 0 {
        return nil
    }
    
    sum := 0.0
    min := window.Readings[0]
    max := window.Readings[0]
    
    for _, v := range window.Readings {
        sum += v
        if v < min {
            min = v
        }
        if v > max {
            max = v
        }
    }
    
    return &AggregatedResult{
        DeviceID:  window.DeviceID,
        Avg:       sum / float64(len(window.Readings)),
        Min:       min,
        Max:       max,
        Count:     len(window.Readings),
        StartTime: window.StartTime,
        EndTime:   window.EndTime,
    }
}

func main() {
    agg := NewAggregator()
    
    // Simulate data
    go func() {
        for i := 0; i < 100; i++ {
            agg.AddReading("sensor-001", 20+float64(i%30))
            time.Sleep(500 * time.Millisecond)
        }
    }()
    
    select {}
}
```

</details>

### แบบฝึกหัดที่ 5: สร้าง REST API สำหรับจัดการอุปกรณ์

**โจทย์:** สร้าง REST API ด้วย Gin ที่รองรับ CRUD operations สำหรับอุปกรณ์ IoT

<details>
<summary>เฉลย (Solution)</summary>

```go
package main

import (
    "net/http"
    "sync"
    
    "github.com/gin-gonic/gin"
)

type Device struct {
    ID      string            `json:"id"`
    Name    string            `json:"name"`
    Type    string            `json:"type"`
    Status  string            `json:"status"`
    Metadata map[string]string `json:"metadata"`
}

type DeviceManager struct {
    devices map[string]Device
    mu      sync.RWMutex
}

func NewDeviceManager() *DeviceManager {
    return &DeviceManager{
        devices: make(map[string]Device),
    }
}

func (dm *DeviceManager) Create(device Device) {
    dm.mu.Lock()
    defer dm.mu.Unlock()
    dm.devices[device.ID] = device
}

func (dm *DeviceManager) Get(id string) (Device, bool) {
    dm.mu.RLock()
    defer dm.mu.RUnlock()
    d, ok := dm.devices[id]
    return d, ok
}

func (dm *DeviceManager) Update(id string, device Device) bool {
    dm.mu.Lock()
    defer dm.mu.Unlock()
    if _, ok := dm.devices[id]; ok {
        dm.devices[id] = device
        return true
    }
    return false
}

func (dm *DeviceManager) Delete(id string) bool {
    dm.mu.Lock()
    defer dm.mu.Unlock()
    if _, ok := dm.devices[id]; ok {
        delete(dm.devices, id)
        return true
    }
    return false
}

func (dm *DeviceManager) List() []Device {
    dm.mu.RLock()
    defer dm.mu.RUnlock()
    devices := make([]Device, 0, len(dm.devices))
    for _, d := range dm.devices {
        devices = append(devices, d)
    }
    return devices
}

func main() {
    r := gin.Default()
    dm := NewDeviceManager()
    
    // Create device
    r.POST("/devices", func(c *gin.Context) {
        var device Device
        if err := c.ShouldBindJSON(&device); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        dm.Create(device)
        c.JSON(http.StatusCreated, device)
    })
    
    // Get device
    r.GET("/devices/:id", func(c *gin.Context) {
        id := c.Param("id")
        device, ok := dm.Get(id)
        if !ok {
            c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
            return
        }
        c.JSON(http.StatusOK, device)
    })
    
    // Update device
    r.PUT("/devices/:id", func(c *gin.Context) {
        id := c.Param("id")
        var device Device
        if err := c.ShouldBindJSON(&device); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        if !dm.Update(id, device) {
            c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
            return
        }
        c.JSON(http.StatusOK, device)
    })
    
    // Delete device
    r.DELETE("/devices/:id", func(c *gin.Context) {
        id := c.Param("id")
        if !dm.Delete(id) {
            c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
            return
        }
        c.JSON(http.StatusNoContent, nil)
    })
    
    // List devices
    r.GET("/devices", func(c *gin.Context) {
        c.JSON(http.StatusOK, dm.List())
    })
    
    r.Run(":8080")
}
```

</details>

---

## สรุป (Summary)

### ประโยชน์ที่ได้รับ (Benefits)

| ข้อ | ประโยชน์ | คำอธิบาย |
|-----|----------|----------|
| 1 | **Performance สูง** | Go จัดการ concurrent requests ได้ดี |
| 2 | **Deploy เดียว** | Binary ไฟล์เดียว ทำงานบน Raspberry Pi ได้ |
| 3 | **Open Source ฟรี** | OSM, Leaflet, MQTT, InfluxDB ใช้ฟรี |
| 4 | **Real-time Tracking** | WebSocket + Leaflet อัปเดตแบบเรียลไทม์ |
| 5 | **Edge Computing** | ประมวลผลข้อมูลใกล้แหล่ง ลด latency |
| 6 | **ต้นทุนต่ำ** | ใช้ฮาร์ดแวร์ราคาถูก (RPi, ESP32) |

### ข้อควรระวัง (Cautions)

| ข้อ | ข้อควรระวัง | คำอธิบาย |
|-----|-------------|----------|
| 1 | **ความปลอดภัย** | ควรใช้ HTTPS/WSS และ authentication |
| 2 | **การจัดการ Buffer** | ต้องออกแบบ buffer เผื่อ network ดับ |
| 3 | **GPS Signal** | GPS อาจไม่มีสัญญาณในอาคาร |
| 4 | **Battery Life** | อุปกรณ์พกพาต้องจัดการ power |
| 5 | **Data Privacy** | ข้อมูลตำแหน่งเป็นข้อมูล sensitive |

### ข้อดี (Advantages)

```
✅ Go มี performance สูง เหมาะกับ edge device
✅ Concurrent ดี รองรับอุปกรณ์หลายตัวพร้อมกัน
✅ Cross-compile รองรับ ARM (Raspberry Pi)
✅ Ecosystem ครบ (MQTT, Kafka, InfluxDB)
✅ Open Source Maps ฟรี ไม่มีค่าใช้จ่าย
✅ WebSocket รองรับ real-time update
✅ Deploy ด้วย Docker สะดวก
```

### ข้อเสีย (Disadvantages)

```
❌ GPS อาจไม่แม่นในเมืองหรือในอาคาร
❌ ต้องมี internet สำหรับ map tiles (OSM)
❌ การจัดการ power สำหรับอุปกรณ์พกพาทำได้ยาก
❌ การ sync time ระหว่างอุปกรณ์ต้องแม่นยำ
❌ Big data ต้องใช้ time-series database เพิ่ม
```

### ข้อห้าม (Prohibitions)

| ข้อห้าม | เหตุผล | วิธีแก้ไข |
|---------|--------|----------|
| **ห้ามเก็บข้อมูลตำแหน่งโดยไม่ได้รับ consent** | ละเมิด PDPA/GDPR | ขออนุญาตก่อนเก็บ |
| **ห้ามใช้ OSM tiles เกิน rate limit** | OSM มีนโยบายจำกัด | ใช้ cache หรือ self-host |
| **ห้ามส่งข้อมูล GPS ทุก millisecond** | กิน bandwidth และ battery | จำกัดความถี่ 1-5 Hz |
| **ห้าม hardcode credentials** | ไม่ปลอดภัย | ใช้ environment variables |
| **ห้าม ignore error** | อาจทำให้ระบบพัง | ตรวจสอบ err ทุกครั้ง |

---

## แหล่งอ้างอิง (References)

| แหล่งข้อมูล | URL | คำอธิบาย |
|------------|-----|----------|
| **LiveTracker** | https://github.com/jlelse/livetracker | GPS tracker ด้วย Go + OsmAnd |
| **Device Hub** | https://github.com/tendry-lab/device-hub | IoT data infrastructure |
| **IoT Edge Gateway** | https://github.com/VDarius-IT/IoT_Data_Aggregation_and_Processing_Gateway | Edge computing with Go |
| **EcoTracker** | https://github.com/xlurr/eco-tracker | OSM + React + Go |
| **OpenStreetMap** | https://www.openstreetmap.org | ฟรีแผนที่โลก |
| **Leaflet.js** | https://leafletjs.com | ไลบรารีแผนที่ JavaScript |
| **MapLibre GL** | https://maplibre.org | ไลบรารีแผนที่ 3D |
| **MQTT.org** | https://mqtt.org | MQTT protocol |
| **InfluxDB** | https://www.influxdata.com | Time-series database |
| **TimescaleDB** | https://www.timescale.com | PostgreSQL for time-series |

---

**จบบทที่: IoT with Go and Open Source Maps** ✅

ในบทนี้คุณได้เรียนรู้:
- 10 กรณีศึกษา IoT Tracking ด้วย Go
- การใช้ Open Source Maps (OSM, Leaflet, MapLibre)
- การสร้าง Real-time GPS Tracking System
- Edge Gateway ด้วย MQTT + Kafka
- Multi-source Positioning (GPS/WiFi/UWB/IMU)
- Docker Compose สำหรับ IoT Stack
- แบบฝึกหัด 5 ข้อพร้อมเฉลย