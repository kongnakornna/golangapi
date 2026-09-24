# 📁 โครงสร้างภายใต้ `internal/modules/`

ตามที่คุณต้องการให้ **ทุก Module ใช้โครงสร้างแบบเดียวกัน** นี่คือโครงสร้างที่สมบูรณ์ภายใต้ `internal/modules/`:

---

```
internal/modules/
│
├── auth/                                    # 📦 MODULE 1: Authentication & User Management
│   │
│   ├── domain/                              # 🏛️ DOMAIN LAYER
│   │   ├── entity/
│   │   │   ├── user.go                      # User Aggregate Root
│   │   │   ├── session.go                   # Session Entity
│   │   │   ├── verification_token.go        # VerificationToken Entity
│   │   │   └── permission.go                # Permission Entity
│   │   │
│   │   ├── value_object/
│   │   │   ├── email.go                     # Email Value Object
│   │   │   ├── password.go                  # Password Value Object
│   │   │   ├── role.go                      # Role Value Object (admin, user, moderator)
│   │   │   ├── status.go                    # Status Value Object (active, inactive, suspended)
│   │   │   ├── token_type.go                # TokenType Value Object
│   │   │   └── user_id.go                   # UserID Value Object
│   │   │
│   │   ├── repository/
│   │   │   ├── user_repository.go           # Interface
│   │   │   ├── session_repository.go        # Interface
│   │   │   └── verification_token_repository.go # Interface
│   │   │
│   │   ├── service/
│   │   │   ├── auth_service.go              # Domain Service
│   │   │   ├── password_hasher.go           # Interface (bcrypt)
│   │   │   └── token_maker.go               # Interface (JWT)
│   │   │
│   │   └── errors/
│   │       └── errors.go                    # Auth Domain Errors
│   │
│   ├── application/                         # 🎯 APPLICATION LAYER
│   │   ├── register.go                      # Register UseCase
│   │   ├── login.go                         # Login UseCase
│   │   ├── refresh_token.go                 # RefreshToken UseCase
│   │   ├── logout.go                        # Logout UseCase
│   │   ├── verify_email.go                  # VerifyEmail UseCase
│   │   ├── change_password.go               # ChangePassword UseCase
│   │   ├── get_profile.go                   # GetProfile UseCase
│   │   ├── update_profile.go                # UpdateProfile UseCase
│   │   ├── list_users.go                    # ListUsers UseCase
│   │   ├── delete_user.go                   # DeleteUser UseCase
│   │   ├── forgot_password.go               # ForgotPassword UseCase
│   │   ├── reset_password.go                # ResetPassword UseCase
│   │   └── dto.go                           # Request/Response DTOs
│   │
│   ├── infrastructure/                      # 🔧 INFRASTRUCTURE LAYER
│   │   ├── persistence/
│   │   │   ├── postgres/
│   │   │   │   ├── user_repo_impl.go        # PostgreSQL Implementation
│   │   │   │   ├── session_repo_impl.go
│   │   │   │   ├── verification_token_repo_impl.go
│   │   │   │   └── models.go                # GORM Models
│   │   │   └── redis/
│   │   │       ├── session_repo_impl.go     # Redis Implementation
│   │   │       └── cache_repo_impl.go
│   │   │
│   │   └── security/
│   │       ├── jwt_maker.go                 # JWT Implementation
│   │       ├── bcrypt_hasher.go             # Bcrypt Implementation
│   │       └── auth_middleware.go           # Auth Middleware
│   │
│   └── interfaces/                          # 🌐 INTERFACE LAYER
│       ├── http/
│       │   ├── auth_handler.go              # HTTP Handlers
│       │   ├── user_handler.go
│       │   ├── routes.go                    # Route Registration
│       │   └── dto.go                       # HTTP DTOs
│       └── middleware/
│           ├── auth.go                      # Auth Middleware
│           ├── cors.go                      # CORS Middleware
│           ├── logging.go                   # Logging Middleware
│           └── rate_limit.go                # Rate Limit Middleware
│
├── iot/                                     # 📦 MODULE 2: IoT Device Management
│   │
│   ├── domain/                              # 🏛️ DOMAIN LAYER
│   │   ├── entity/
│   │   │   ├── device.go                    # Device Aggregate Root
│   │   │   ├── telemetry.go                 # Telemetry Entity
│   │   │   ├── alert.go                     # Alert Entity
│   │   │   ├── alert_rule.go                # AlertRule Entity
│   │   │   ├── device_group.go              # DeviceGroup Entity
│   │   │   ├── device_group_member.go       # DeviceGroupMember Entity
│   │   │   ├── schedule.go                  # Schedule Entity
│   │   │   ├── device_schedule.go           # DeviceSchedule Entity
│   │   │   ├── notification.go              # Notification Entity
│   │   │   ├── notification_channel.go      # NotificationChannel Entity
│   │   │   ├── notification_config.go       # NotificationConfig Entity
│   │   │   ├── location.go                  # Location Entity
│   │   │   ├── device_config.go             # DeviceConfig Entity
│   │   │   ├── device_status.go             # DeviceStatus Entity
│   │   │   ├── device_type.go               # DeviceType Entity
│   │   │   ├── device_category.go           # DeviceCategory Entity
│   │   │   ├── mqtt_host.go                 # MQTTHost Entity
│   │   │   ├── mqtt_log.go                  # MQTTLog Entity
│   │   │   ├── command_log.go               # CommandLog Entity
│   │   │   ├── activity_log.go              # ActivityLog Entity
│   │   │   ├── audit_log.go                 # AuditLog Entity
│   │   │   ├── api_key.go                   # APIKey Entity
│   │   │   ├── report_data.go               # ReportData Entity
│   │   │   ├── sensor_data.go               # SensorData Entity
│   │   │   ├── iot_data.go                  # IoTData Entity
│   │   │   └── system_setting.go            # SystemSetting Entity
│   │   │
│   │   ├── value_object/
│   │   │   ├── device_id.go                 # DeviceID Value Object
│   │   │   ├── device_status.go             # DeviceStatus Value Object (online, offline, error)
│   │   │   ├── device_type.go               # DeviceType Value Object
│   │   │   ├── metric.go                    # Metric Value Object
│   │   │   ├── severity.go                  # Severity Value Object (info, warning, critical)
│   │   │   ├── alert_status.go              # AlertStatus Value Object
│   │   │   ├── condition.go                 # Condition Value Object (gt, lt, eq, between)
│   │   │   ├── coordinates.go               # Coordinates Value Object
│   │   │   ├── time_range.go                # TimeRange Value Object
│   │   │   ├── notification_type.go         # NotificationType Value Object
│   │   │   ├── channel_type.go              # ChannelType Value Object
│   │   │   ├── schedule_type.go             # ScheduleType Value Object
│   │   │   ├── command_status.go            # CommandStatus Value Object
│   │   │   └── report_type.go               # ReportType Value Object
│   │   │
│   │   ├── repository/
│   │   │   ├── device_repository.go         # Interface
│   │   │   ├── telemetry_repository.go      # Interface
│   │   │   ├── alert_repository.go          # Interface
│   │   │   ├── alert_rule_repository.go     # Interface
│   │   │   ├── device_group_repository.go   # Interface
│   │   │   ├── schedule_repository.go       # Interface
│   │   │   ├── notification_repository.go   # Interface
│   │   │   ├── location_repository.go       # Interface
│   │   │   ├── mqtt_repository.go           # Interface
│   │   │   └── audit_repository.go          # Interface
│   │   │
│   │   ├── service/
│   │   │   ├── device_manager.go            # Domain Service
│   │   │   ├── alert_engine.go              # Domain Service
│   │   │   ├── telemetry_processor.go       # Domain Service
│   │   │   ├── notification_service.go      # Domain Service
│   │   │   ├── schedule_executor.go         # Domain Service
│   │   │   ├── mqtt_manager.go              # Domain Service
│   │   │   ├── location_service.go          # Domain Service
│   │   │   └── data_aggregator.go           # Domain Service
│   │   │
│   │   └── errors/
│   │       └── errors.go                    # IoT Domain Errors
│   │
│   ├── application/                         # 🎯 APPLICATION LAYER
│   │   ├── device/
│   │   │   ├── register_device.go
│   │   │   ├── update_device.go
│   │   │   ├── delete_device.go
│   │   │   ├── get_device.go
│   │   │   ├── list_devices.go
│   │   │   ├── get_device_telemetry.go
│   │   │   ├── send_device_command.go
│   │   │   ├── update_device_status.go
│   │   │   ├── get_device_status_history.go
│   │   │   └── dto.go
│   │   │
│   │   ├── telemetry/
│   │   │   ├── process_telemetry.go
│   │   │   ├── process_batch_telemetry.go
│   │   │   ├── get_telemetry_stats.go
│   │   │   ├── query_telemetry.go
│   │   │   ├── get_telemetry_aggregation.go
│   │   │   └── dto.go
│   │   │
│   │   ├── alert/
│   │   │   ├── create_alert_rule.go
│   │   │   ├── update_alert_rule.go
│   │   │   ├── delete_alert_rule.go
│   │   │   ├── get_alert_rule.go
│   │   │   ├── list_alert_rules.go
│   │   │   ├── list_alerts.go
│   │   │   ├── acknowledge_alert.go
│   │   │   ├── resolve_alert.go
│   │   │   ├── get_alert_stats.go
│   │   │   └── dto.go
│   │   │
│   │   ├── device_group/
│   │   │   ├── create_group.go
│   │   │   ├── update_group.go
│   │   │   ├── delete_group.go
│   │   │   ├── get_group.go
│   │   │   ├── list_groups.go
│   │   │   ├── add_device_to_group.go
│   │   │   ├── remove_device_from_group.go
│   │   │   └── dto.go
│   │   │
│   │   ├── schedule/
│   │   │   ├── create_schedule.go
│   │   │   ├── update_schedule.go
│   │   │   ├── delete_schedule.go
│   │   │   ├── get_schedule.go
│   │   │   ├── list_schedules.go
│   │   │   ├── execute_schedule.go
│   │   │   └── dto.go
│   │   │
│   │   ├── notification/
│   │   │   ├── create_notification.go
│   │   │   ├── list_notifications.go
│   │   │   ├── mark_as_read.go
│   │   │   ├── create_channel.go
│   │   │   ├── update_channel.go
│   │   │   └── dto.go
│   │   │
│   │   └── location/
│   │       ├── create_location.go
│   │       ├── update_location.go
│   │       ├── delete_location.go
│   │       ├── get_location.go
│   │       ├── list_locations.go
│   │       └── dto.go
│   │
│   ├── infrastructure/                      # 🔧 INFRASTRUCTURE LAYER
│   │   ├── persistence/
│   │   │   ├── postgres/
│   │   │   │   ├── device_repo_impl.go
│   │   │   │   ├── telemetry_repo_impl.go
│   │   │   │   ├── alert_repo_impl.go
│   │   │   │   ├── alert_rule_repo_impl.go
│   │   │   │   ├── device_group_repo_impl.go
│   │   │   │   ├── schedule_repo_impl.go
│   │   │   │   ├── notification_repo_impl.go
│   │   │   │   ├── location_repo_impl.go
│   │   │   │   ├── mqtt_repo_impl.go
│   │   │   │   ├── audit_repo_impl.go
│   │   │   │   └── models.go               # GORM Models
│   │   │   └── influxdb/
│   │   │       └── telemetry_repo_impl.go  # InfluxDB Implementation
│   │   │
│   │   ├── mqtt/
│   │   │   ├── client.go                   # MQTT Client
│   │   │   ├── manager.go                  # Connection Manager
│   │   │   ├── handler.go                  # Message Handler
│   │   │   └── options.go                  # Configuration
│   │   │
│   │   ├── websocket/
│   │   │   ├── hub.go                      # WebSocket Hub
│   │   │   ├── client.go                   # WebSocket Client
│   │   │   └── message.go                  # Message Types
│   │   │
│   │   └── external/
│   │       ├── influx_client.go            # InfluxDB Client
│   │       └── redis_client.go             # Redis Client
│   │
│   └── interfaces/                          # 🌐 INTERFACE LAYER
│       ├── http/
│       │   ├── device_handler.go
│       │   ├── telemetry_handler.go
│       │   ├── alert_handler.go
│       │   ├── device_group_handler.go
│       │   ├── schedule_handler.go
│       │   ├── notification_handler.go
│       │   ├── location_handler.go
│       │   ├── websocket_handler.go
│       │   ├── routes.go
│       │   └── dto.go
│       │
│       └── worker/
│           ├── mqtt_worker.go              # MQTT Consumer
│           ├── alert_worker.go             # Alert Processor
│           ├── schedule_worker.go          # Schedule Executor
│           ├── notification_worker.go      # Notification Sender
│           └── report_worker.go            # Report Generator
│
├── dashboard/                               # 📦 MODULE 3: Dashboard & Visualization
│   │
│   ├── domain/                              # 🏛️ DOMAIN LAYER
│   │   ├── entity/
│   │   │   ├── dashboard.go                 # Dashboard Aggregate Root
│   │   │   ├── widget.go                    # Widget Entity
│   │   │   ├── report.go                    # Report Entity
│   │   │   ├── chart.go                     # Chart Entity
│   │   │   ├── data_source.go               # DataSource Entity
│   │   │   └── user_preference.go           # UserPreference Entity
│   │   │
│   │   ├── value_object/
│   │   │   ├── widget_type.go               # WidgetType Value Object
│   │   │   ├── chart_type.go                # ChartType Value Object
│   │   │   ├── position.go                  # Position Value Object
│   │   │   ├── size.go                      # Size Value Object
│   │   │   ├── report_format.go             # ReportFormat Value Object
│   │   │   ├── time_range.go                # TimeRange Value Object
│   │   │   ├── aggregation_type.go          # AggregationType Value Object
│   │   │   └── color_scheme.go              # ColorScheme Value Object
│   │   │
│   │   ├── repository/
│   │   │   ├── dashboard_repository.go      # Interface
│   │   │   ├── widget_repository.go         # Interface
│   │   │   ├── report_repository.go         # Interface
│   │   │   └── data_source_repository.go    # Interface
│   │   │
│   │   ├── service/
│   │   │   ├── dashboard_service.go         # Domain Service
│   │   │   ├── chart_service.go             # Domain Service
│   │   │   ├── report_service.go            # Domain Service
│   │   │   ├── data_aggregator.go           # Domain Service
│   │   │   └── export_service.go            # Domain Service
│   │   │
│   │   └── errors/
│   │       └── errors.go                    # Dashboard Domain Errors
│   │
│   ├── application/                         # 🎯 APPLICATION LAYER
│   │   ├── create_dashboard.go
│   │   ├── update_dashboard.go
│   │   ├── delete_dashboard.go
│   │   ├── get_dashboard.go
│   │   ├── list_dashboards.go
│   │   ├── clone_dashboard.go
│   │   ├── add_widget.go
│   │   ├── update_widget.go
│   │   ├── delete_widget.go
│   │   ├── rearrange_widgets.go
│   │   ├── get_chart_data.go
│   │   ├── generate_report.go
│   │   ├── get_report.go
│   │   ├── list_reports.go
│   │   ├── delete_report.go
│   │   ├── export_report.go
│   │   ├── get_dashboard_stats.go
│   │   └── dto.go
│   │
│   ├── infrastructure/                      # 🔧 INFRASTRUCTURE LAYER
│   │   ├── persistence/
│   │   │   ├── postgres/
│   │   │   │   ├── dashboard_repo_impl.go
│   │   │   │   ├── widget_repo_impl.go
│   │   │   │   ├── report_repo_impl.go
│   │   │   │   └── models.go
│   │   │   └── redis/
│   │   │       └── cache_repo_impl.go
│   │   │
│   │   └── export/
│   │       ├── pdf_exporter.go
│   │       ├── excel_exporter.go
│   │       └── csv_exporter.go
│   │
│   └── interfaces/                          # 🌐 INTERFACE LAYER
│       ├── http/
│       │   ├── dashboard_handler.go
│       │   ├── widget_handler.go
│       │   ├── report_handler.go
│       │   ├── routes.go
│       │   └── dto.go
│       │
│       └── worker/
│           └── report_generator.go          # Background Report Generation
│
├── shared/                                   # 🔄 SHARED BETWEEN MODULES
│   ├── domain/
│   │   ├── errors/
│   │   │   └── errors.go                    # Shared Domain Errors
│   │   └── value_object/
│   │       ├── pagination.go                # Pagination Value Object
│   │       ├── time_range.go                # TimeRange Value Object
│   │       └── sort_order.go                # SortOrder Value Object
│   │
│   ├── application/
│   │   └── dto/
│   │       ├── common_response.go           # Common Response DTO
│   │       └── pagination_request.go        # Pagination Request DTO
│   │
│   └── infrastructure/
│       ├── database/
│       │   ├── postgres.go                  # PostgreSQL Connection
│       │   └── redis.go                     # Redis Connection
│       ├── logger/
│       │   └── logger.go                    # Shared Logger
│       ├── config/
│       │   └── config.go                    # Shared Config
│       └── utils/
│           ├── validator.go                 # Shared Validator
│           ├── id_generator.go              # ID Generator
│           └── time_utils.go                # Time Utilities
│
├── wire/                                     # 🔌 DEPENDENCY INJECTION (Wire)
│   ├── auth_wire.go
│   ├── iot_wire.go
│   ├── dashboard_wire.go
│   └── wire.go
│
└── migrations/                               # 🗄️ DATABASE MIGRATIONS
    ├── auth/
    │   ├── 001_create_users_table.up.sql
    │   ├── 001_create_users_table.down.sql
    │   ├── 002_create_sessions_table.up.sql
    │   └── 003_create_verification_tokens_table.up.sql
    │
    ├── iot/
    │   ├── 001_create_devices_table.up.sql
    │   ├── 002_create_telemetry_table.up.sql
    │   ├── 003_create_alerts_table.up.sql
    │   ├── 004_create_alert_rules_table.up.sql
    │   ├── 005_create_device_groups_table.up.sql
    │   ├── 006_create_schedules_table.up.sql
    │   ├── 007_create_notifications_table.up.sql
    │   ├── 008_create_locations_table.up.sql
    │   ├── 009_create_audit_logs_table.up.sql
    │   └── 010_create_api_keys_table.up.sql
    │
    └── dashboard/
        ├── 001_create_dashboards_table.up.sql
        ├── 002_create_widgets_table.up.sql
        └── 003_create_reports_table.up.sql
```

---

## 📊 สรุปโครงสร้างของ 1 Module (Template)

```yaml
{module_name}/:
  domain/:
    entity/:
      - aggregate_root.go
      - entities.go
    value_object/:
      - value_objects.go
    repository/:
      - interfaces.go
    service/:
      - domain_services.go
    errors/:
      - errors.go
  
  application/:
    {feature}/:
      - usecase.go
      - dto.go
  
  infrastructure/:
    persistence/:
      postgres/:
        - repo_impl.go
        - models.go
      redis/:
        - repo_impl.go
    external/:
      - clients.go
  
  interfaces/:
    http/:
      - handler.go
      - routes.go
      - dto.go
    worker/:
      - workers.go
```

### กฎการพึ่งพา (Dependency Rule)

```
interfaces/  →  application/  →  domain/
     ↑              ↑               ↑
infrastructure/  →  domain/    ←  repository/
     ↑
(implements)
```