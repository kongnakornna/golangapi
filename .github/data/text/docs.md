# โครงสร้าง Folder และ File

## ภาพรวมทั้งโปรเจกต์

```mermaid
flowchart TB
    subgraph Entry["cmd/ - Entry Points"]
        CMD_API["cmd/api/main.go"]
        CMD_KAFKA["cmd/kafka/main.go"]
        CMD_WS["cmd/websocket/main.go"]
        CMD_SERVE["cmd/serve.go"]
        CMD_WORKER["cmd/worker.go"]
        CMD_MIGRATE["cmd/migrate.go"]
        CMD_INIT["cmd/initdata.go"]
    end

    subgraph Config["config/"]
        CFG["config.go"]
        CFG_YML["config.dev.yml"]
    end

    subgraph Internal["internal/ - Core Logic (DDD)"]
        subgraph Middleware["middleware/"]
            MW_CORS["cors.go"]
            MW_JWT["jwtauth.go"]
            MW_RATE["rate_limit.go"]
            MW_SEC["security.go"]
            MW_MON["monitoring.go"]
            MW_LOG["logging.go"]
        end

        subgraph Domains["Domain Modules"]
            direction TB
            IOT["iot/"]
            USERS["users/"]
            ALARM["alarm/"]
            AUTH["auth/"]
            ITEMS["items/"]
            KAFKA["kafka/"]
            MQTT["mqtt/"]
            INFLUX["influxdb/"]
            WS["websocket/"]
        end

        subgraph Shared["Shared"]
            MODELS["models/"]
            REPO["repository/"]
            UC["usecase/"]
            SERVER["server/"]
            DELIVERY["delivery/"]
            DIST["distributor/"]
            PROC["processor/"]
            QUEUE["queue/"]
            WORKER["worker/"]
        end
    end

    subgraph Pkg["pkg/ - Reusable Packages"]
        PG["db/postgres/"]
        RD["db/redis/"]
        KFK["kafka/"]
        MQ["mqtt/"]
        WS_PKG["websocket/"]
        JWT["jwt/"]
        LOG["logger/"]
        HTTP_ERR["httpErrors/"]
        RESP["responses/"]
        INFLUX_PKG["influxdb/"]
        EMAIL["sendEmail/"]
        CRYPT["cryptpass/"]
    end

    subgraph Docs["docs/"]
        SWAGGER["swagger.json/yaml"]
    end

    subgraph DB_Files["db/, migrations/"]
        SQL["icmon.sql"]
        MIG["migration files"]
    end

    Entry --> Config
    Entry --> Internal
    Entry --> Pkg

    Internal --> Middleware
    Internal --> Domains
    Internal --> Shared

    Domains --> Pkg
    Domains --> Shared

    Pkg --> PG & RD & KFK & MQ & WS_PKG

    subgraph DDD_Pattern["DDD Layer per Domain"]
        direction LR
        DEL["delivery/http/\n(handler + routes)"]
        UC_LAYER["usecase/\n(business logic)"]
        REPO_LAYER["repository/\n(data access)"]
        PRES["presenter/\n(response format)"]
        PROC_LAYER["processor/\n(async processing)"]
        DIST_LAYER["distributor/\n(event distribution)"]
    end

    DEL --> UC_LAYER --> REPO_LAYER
    UC_LAYER --> PRES
    UC_LAYER --> PROC_LAYER --> DIST_LAYER
```

---

## แยกแต่ละ Module

### 1. Module: iot/ — main IoT domain

```mermaid
flowchart TB
    subgraph IOT["internal/modules/iot/"]
        direction TB

        subgraph DELIVERY["delivery/http/"]
            H["handler.go"]
            R["routes.go"]
        end

        subgraph MODELS["models/ (44 files)"]
            DEVICE["device.go"]
            SENSOR["sensor_data.go"]
            DEVCFG["device_config.go"]
            DEVSTATUS["device_status.go"]
            ALARM_IOT["alarm.go"]
            MQTT_IOT["mqtt.go"]
            SCHEDULE["schedule.go"]
            OTHERS["activity_log.go,\nair_control.go,\nlocation.go,\n..."]
        end

        subgraph REPOSITORY["repository/"]
            DEVREPO["device_repo.go"]
            DEVCFGREPO["device_config_repo.go"]
            DEVSTATREPO["device_status_repo.go"]
            IOTDATA["iot_data_repo.go"]
            SCHEDREPO["schedule_repo.go"]
            ACTLOG["activity_log_repo.go"]
            ALARMREPO["alarm_log_repo.go"]
            DEVALERT["device_alert_repo.go"]
        end

        subgraph USECASE["usecase/"]
            UC_IOT["usecase.go"]
        end

        subgraph PRESENTER["presenter/"]
            PRES_IOT["presenter.go"]
        end

        subgraph IOTHELPER["iothelper/"]
            ALARM_HELPER["alarm.go"]
        end
    end

    subgraph PKG_IOT["pkg/ dependencies"]
        PG_IOT["db/postgres/"]
        RD_IOT["db/redis/"]
        MQ_IOT["mqtt/"]
        KFK_IOT["kafka/"]
        INFLUX_IOT["influxdb/"]
        LOG_IOT["logger/"]
    end

    subgraph EXTERNAL["External"]
        API_REQ["HTTP Request"]
        MQTT_BROKER["MQTT Broker"]
        KAFKA_BROKER["Kafka"]
        WS_CLIENT["WebSocket Client"]
    end

    API_REQ --> DELIVERY
    MQTT_BROKER --> PKG_IOT
    KAFKA_BROKER --> PKG_IOT

    DELIVERY --> UC_IOT
    UC_IOT --> MODELS
    UC_IOT --> REPOSITORY
    UC_IOT --> PRES_IOT
    UC_IOT --> IOTHELPER
    UC_IOT --> PKG_IOT

    REPOSITORY --> PG_IOT & RD_IOT
    IOTHELPER --> MODELS
```

---

### 2. Module: users/ — User management

```mermaid
flowchart TB
    subgraph USERS["internal/modules/users/"]
        direction TB

        subgraph DEL_USERS["delivery/http/"]
            H_USERS["handlers.go"]
            R_USERS["routes.go"]
        end

        subgraph UC_USERS["usecase/"]
            UC_U["usecase.go"]
        end

        subgraph REP_USERS["repository/"]
            PG_U["pg_repository.go"]
            RD_U["redis_repository.go"]
        end

        subgraph PRES_USERS["presenter/"]
            PRES_U["presenters.go"]
        end

        subgraph PROC_USERS["processor/"]
            PROC_U["processor.go"]
        end

        subgraph DIST_USERS["distributor/"]
            DIST_U["distributor.go"]
        end

        subgraph ROOT_USERS["root-level legacy"]
            H_LEGACY["handler.go"]
            PG_LEGACY["pg_repository.go"]
            RD_LEGACY["redis_repository.go"]
            UC_LEGACY["usecase.go"]
            WORKER_U["worker.go"]
        end
    end

    API_USERS["HTTP Request"] --> DEL_USERS
    DEL_USERS --> UC_U
    UC_U --> REP_USERS
    UC_U --> PRES_U
    UC_U --> PROC_U
    PROC_U --> DIST_U
    UC_U --> ROOT_USERS

    REP_USERS --> PG_U & RD_U
```

---

### 3. Module: alarm/ — Alarm & notification

```mermaid
flowchart TB
    subgraph ALARM["internal/modules/alarm/"]
        direction TB

        DEL_ALARM["delivery/http/\nhandler.go + routes.go"]
        UC_ALARM["usecase/usecase.go"]
        REP_ALARM["repository/alarm_log_repo.go"]
        PRES_ALARM["presenter/presenter.go"]
        PROC_ALARM["processor/processor.go"]
        MOD_ALARM["models/alarm_log.go"]
    end

    API_ALARM["HTTP Request"] --> DEL_ALARM
    DEL_ALARM --> UC_ALARM
    UC_ALARM --> REP_ALARM
    UC_ALARM --> PRES_ALARM
    UC_ALARM --> PROC_ALARM
    REP_ALARM --> MOD_ALARM
    REP_ALARM --> PG_ALARM["PostgreSQL"]
```

---

### 4. Module: auth/ — Authentication

```mermaid
flowchart TB
    subgraph AUTH["internal/modules/auth/"]
        direction TB

        DEL_AUTH["delivery/http/\nhandlers.go + routes.go"]
        H_AUTH["handler.go (root)"]
    end

    API_AUTH["HTTP Request (login, register,\nrefresh token)"] --> DEL_AUTH
    DEL_AUTH --> H_AUTH
    H_AUTH --> PKG_AUTH["pkg/jwt/"]
    H_AUTH --> PKG_CRYPT["pkg/cryptpass/"]
    H_AUTH --> USERS_MOD["internal/modules/users/"]
```

---

### 5. Module: items/ — CRUD items

```mermaid
flowchart TB
    subgraph ITEMS["internal/modules/items/"]
        direction TB

        DEL_ITEM["delivery/http/\nhandlers.go + routes.go"]
        UC_ITEM["usecase/usecase.go"]
        REP_ITEM["repository/pg_repository.go"]
        PRES_ITEM["presenter/presenter.go"]

        ROOT_H["handler.go"]
        ROOT_PG["pg_repository.go"]
        ROOT_UC["usecase.go"]
    end

    API_ITEM["HTTP Request"] --> DEL_ITEM
    DEL_ITEM --> UC_ITEM
    UC_ITEM --> REP_ITEM
    UC_ITEM --> PRES_ITEM
    UC_ITEM --> ROOT_UC
    ROOT_UC --> ROOT_PG
    ROOT_H --> DEL_ITEM
    REP_ITEM --> PG_ITEM["PostgreSQL"]
```

---

### 6. Module: kafka/ — Kafka message handling

```mermaid
flowchart TB
    subgraph KAFKA["internal/modules/kafka/"]
        direction TB

        DEL_KAFKA_HTTP["delivery/http/handler.go"]
        DEL_KAFKA_WS["delivery/ws/handler.go"]
        WS_KAFKA["ws/handler.go"]
        UC_KAFKA["usecase/order_usecase.go"]
        REP_KAFKA["repository/order_repo.go"]
        MOD_KAFKA["models/order.go"]
    end

    subgraph KAFKA_PKG["pkg/kafka/"]
        PROD["producer.go"]
        CONS["consumer.go"]
        CFG_K["config.go"]
        MSG["message.go"]
    end

    KAFKA_BROKER["Kafka Broker"] <--> KAFKA_PKG
    KAFKA_PKG --> UC_KAFKA
    UC_KAFKA --> REP_KAFKA
    UC_KAFKA --> MOD_KAFKA
    DEL_KAFKA_HTTP --> UC_KAFKA
    DEL_KAFKA_WS --> UC_KAFKA
    WS_KAFKA --> UC_KAFKA
    REP_KAFKA --> PG_KAFKA["PostgreSQL"]
```

---

### 7. Module: mqtt/ — MQTT IoT gateway

```mermaid
flowchart TB
    subgraph MQTT["internal/modules/mqtt/"]
        direction TB

        DEL_MQTT["delivery/http/\nhandler.go + routes.go"]
        UC_MQTT["usecase/usecase.go"]
        PRES_MQTT["presenter/presenter.go"]
    end

    subgraph MQTT_PKG["pkg/mqtt/"]
        CLIENT["client.go"]
    end

    API_MQTT["HTTP Request"] --> DEL_MQTT
    DEL_MQTT --> UC_MQTT
    UC_MQTT --> PRES_MQTT
    UC_MQTT --> MQTT_PKG
    MQTT_PKG <--> MQTT_BROKER["MQTT Broker (HiveMQ)"]
    MQTT_BROKER <--> IOT_DEVICES["IoT Devices"]
    UC_MQTT --> INFLUX_MQTT["InfluxDB"]
    UC_MQTT --> PG_MQTT["PostgreSQL"]
```

---

### 8. Module: influxdb/ — Time-series data query

```mermaid
flowchart TB
    subgraph INFLUX["internal/modules/influxdb/"]
        direction TB

        DEL_INFLUX["delivery/http/\nhandler.go + routes.go"]
        UC_INFLUX["usecase/usecase.go"]
        PRES_INFLUX["presenter/presenter.go"]
    end

    subgraph INFLUX_PKG["pkg/influxdb/"]
        CLIENT_INFLUX["client.go"]
    end

    API_INFLUX["HTTP Request"] --> DEL_INFLUX
    DEL_INFLUX --> UC_INFLUX
    UC_INFLUX --> PRES_INFLUX
    UC_INFLUX --> INFLUX_PKG
    INFLUX_PKG <--> INFLUXDB["InfluxDB (Time-series)"]
```

---

### 9. Module: websocket/ — Real-time WebSocket

```mermaid
flowchart TB
    subgraph WS["internal/modules/websocket/"]
        direction TB

        DEL_WS_HTTP["delivery/http/handler.go"]
        DEL_WS_WS["delivery/ws/handler.go"]
        UC_WS["usecase/ws_usecase.go"]
        REP_WS["repository/ws_repo.go"]
        REP_WS_PG["repository/postgres/ws_repo_pg.go"]
        MOD_WS["models/ws_models.go"]
        IFACE_WS["interface.go"]
    end

    subgraph WS_PKG["pkg/websocket/"]
        HUB["hub.go"]
        CLIENT_WS["client.go"]
        IFACE_PKG["interface.go"]
        MSG_WS["message.go"]
    end

    WS_CLIENTS["WebSocket Clients"] <--> WS_PKG
    WS_PKG --> UC_WS
    UC_WS --> REP_WS
    REP_WS --> REP_WS_PG
    UC_WS --> MOD_WS
    DEL_WS_WS --> UC_WS
    DEL_WS_HTTP --> UC_WS
    REP_WS_PG --> PG_WS["PostgreSQL"]
```

---

### 10. Shared Modules

```mermaid
flowchart LR
    subgraph SHARED["internal/ (shared)"]
        MODELS_SH["models/\n(base.go, device.go,\nuser.go, ...)"]
        REPO_SH["repository/\npg.go, redis.go"]
        UC_SH["usecase/usecase.go"]
        SERVER_SH["server/\nserver.go, handlers.go"]
        DEL_SH["delivery/rest/\nrouter.go"]
        DIST_SH["distributor/distributor.go"]
        PROC_SH["processor/processor.go"]
        QUEUE_SH["queue/\nmanager.go"]
        WORKER_SH["worker/worker.go"]
    end

    DOMAINS["Domain Modules\n(iot, users, alarm, ...)"] --> MODELS_SH
    DOMAINS --> REPO_SH
    DOMAINS --> UC_SH
    DOMAINS --> DIST_SH
    DOMAINS --> PROC_SH

    SERVER_SH --> DEL_SH
    DEL_SH --> MIDDLEWARE["middleware/\n(cors, jwt, rate_limit,\nmonitoring, security)"]
```

---

## Flow คำอธิบาย

1. **cmd/** — entry points เริ่มต้น (API, Kafka consumer, WebSocket, Worker)
2. **internal/** — core logic แบบ DDD แต่ละ domain มี layer: `delivery/http` → `usecase` → `repository` + `presenter` + `processor` + `distributor`
3. **pkg/** — shared packages (DB, Kafka, MQTT, JWT, Logger, ฯลฯ)
4. **middleware/** — intercept request (CORS, JWT auth, rate limit, monitoring)
5. **config/** → inject config ให้ทุก module
6. **db/ + migrations/** → schema ฐานข้อมูล


------------------------------
skill
1. สั่ง copilot  
   - copilot
   - /tb-brainstorming
   - แล้ว เอาโจทย์มา วาง เพื่อให้  AI วิเคาะห์ 
   - /tb-writing-plans
   -  ถาม และ ให้ ทางเลือก แสดงแนวทางตัดสินใจ ไปจน กว่าจะ สมบุรณ์
   - /tb-executing-plans
   - เริ่ม ทำงาน ตาม plans


/tb-scrutinize 
- Affected กับส่วนไหนบ้าง
- panic
- code ใน branches feature/rtrpms-2916 มันไปเปลียน  พถติกรรม ใน branches  dev  หรือไหม
- สิ่งที่ต้องการ code ใน  branches feature/rtrmps-2916   ต้องไม่เปลียนแปลง  พถติกรรม ใน branches  dev  code ใน branches  dev ต้องใช้งานได้ โดยไม่มีผลกระทบ 
- แค่นำ code นี้ไป เพิ่มเข้าไป โดยไม่มีผลกระทบ code ใน BRANCHES dev  
- รายงานสรุปผลการดำเนิดการ อะไรไปบ้าง และบอกกระบวนการทดสอบผล

