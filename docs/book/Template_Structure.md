# โครงสร้างโปรเจกต์ icmongolang (Structure Map)

รวบรวมจากการตรวจสอบ live filesystem ณ root ` icmongolang`
คำศัพท์ 3 กลุ่ม: **ระบบ** = Global infrastructure / root / อื่น ๆ ที่ไม่ใช่ตัว module, **รวม** = shared packages / ร่วมกันทุก module, **แยกรราย** = แต่ละ module แยกตามโฟลเดอร์ของตัวเอง

ไฟล์บางรายการถูกตัดออกจากแผนที่: `vendor/`, `.git/`, `app-postgres-data/`, `tmp/`

---

## ระบบ (โครงสร้างระดับ Global)

```
icmongolang
├─ main.go    # entrypoint หลัก
├─ go.mod / go.sum
├─ Makefile / build.ps1 / run.ps1
├─ air.cmd / air.ps1 / .air.toml / .air.toml.linux / .air_1.toml / .air_2.toml
├─ Dockerfile (dev / prod / ollama-exporter)
├─ docker-compose.yml (base / .dev / .prod)
├─ .env / .env.example / .gitattributes / .gitignore
├─ icmongolang.exe / ollama-exporter.exe / package-lock.json / queryex
├─ Note.md / NoteAI.md / README.md / README_IO.md
├─ .github/            # CI/CD workflows
├─ .gitlab-ci/         # GitLab CI
├─ .vscode/
├─ cmd/
│  ├─ api/
│  ├─ kafka/
│  ├─ ollama-exporter/
│  ├─ scheduler/
│  ├─ websocket/
│  ├─ workers/
│  ├─ initdata.go
│  ├─ migrate.go
│  ├─ root.go
│  ├─ serve.go
│  ├─ vectordata.go
│  └─ worker.go
├─ config/
│  ├─ config.default.yml
│  ├─ config.dev.yml
│  ├─ config.go
│  └─ verify_env_test.go
├─ db/
│  ├─ icmon.sql
│  ├─ public.sql
│  ├─ public_DB.sql
│  ├─ public_icmongolang.sql
│  ├─ public_icmongolang_v1.sql
│  └─ sd_iot_device.sql
├─ migrations/
│  ├─ 20250619_kafka_tables.sql
│  ├─ 20250619_websocket_tables.sql
│  ├─ 20260712_new_modules_schema.sql
│  ├─ 20260712_seed_data.sql
│  ├─ 20260827_slow_sql_device_indexes.sql
│  ├─ 20260829_fullschedule.sql
│  ├─ 20260830_fullschedule_scope.sql
│  ├─ 20260831_fullschedule_legacy_backfill.sql
│  ├─ 20260832_fs_schedule_sdiot_columns.sql
│  ├─ 20260833_fullschedule_baseline_audit.sql
│  ├─ 20260834_fullschedule_sort_order.sql
│  ├─ 20260835_fullschedule_baseline_gap.sql
│  ├─ 20260901_fs_schedule_demo_data.sql
│  ├─ 20260901_fs_schedule_history_settings_seed.sql
│  ├─ db.sql
│  ├─ icmon.sql
│  ├─ icmon_database.sql
│  ├─ public.sql
│  ├─ public_db.sql
│  └─ sd_user_access_menu.sql
├─ monitoring/
│  ├─ prometheus.yml
│  ├─ kafka-dashboard.json
│  └─ grafana/
│     ├─ dashboards/
│     │  ├─ elasticsearch/ *-dashboard.json
│     │  ├─ golang/ *-dashboard.json
│     │  ├─ health/ *-dashboard.json
│     │  ├─ influxdb/ *-dashboard.json
│     │  ├─ kafka/ *-dashboard.json
│     │  ├─ mqtt/ *-dashboard.json
│     │  ├─ nodered/ *-dashboard.json
│     │  ├─ ollama/ *-dashboard.json
│     │  ├─ postgres/ *-dashboard.json
│     │  ├─ redis/ *-dashboard.json
│     │  └─ websocket/ *-dashboard.json
│     └─ provisioning/
│        ├─ dashboards/
│        │  ├─ elasticsearch.yml / golang.yml / health.yml / influxdb.yml
│        │  ├─ kafka.yml / mqtt.yml / nodered.yml / ollama.yml
│        │  └─ postgres.yml / redis.yml / websocket.yml
│        └─ datasources/
│           └─ prometheus.yml
├─ elk/                # Elasticsearch / Logstash / Kibana
├─ influxdb/
├─ mqtt/
├─ linux/
├─ scripts/
├─ test/
├─ postman/
├─ docs/
│  ├─ structure-map.md            # ไฟล์นี้
│  ├─ docs.md
│  ├─ docs.go (swagger embed)
│  ├─ swagger.json / swagger.yaml
│  └─ modules_dev/
│     ├─ Template_Modules.md
│     ├─ Template_Modules_PDPA.md
│     ├─ Template_prompt_ai.md
│     ├─ PDPA_Modules_V4.md
│     └─ PDPA_Modules_V5.md
├─ ai/
├─ internal/            # Go source (ดู รวม และ แยกรราย ด้านล่าง)
└─ pkg/                 # shared packages (ดู รวม)
```

---

## รวม (โครงสร้างระดับ shared / ร่วมกัน)

```
internal/                       # application-level shared + modules
├─ server/
│  ├─ api_cache_adapter.go
│  ├─ handlers.go
│  └─ server.go
├─ delivery/
│  └─ rest/
│     ├─ router.go
│     └─ middleware/monitoring.go
├─ middleware/
│  ├─ cors.go / jwtauth.go / jwtauth_test.go / logging.go
│  ├─ middleware.go / monitoring.go / prometheus.go
│  ├─ rate_limit.go / security.go
├─ models/
│  ├─ air_control.go / base.go / device.go / models.go / payment.go
│  ├─ refresh_token.go / sd_user.go / sd_user_role.go
│  ├─ sd_user_roles_access.go / sd_user_roles_permission.go
│  ├─ user.go / ws_models.go
├─ repository/
│  ├─ pg.go / redis.go
│  └─ ReadMe_Modules.md
├─ distributor/
│  └─ distributor.go
├─ processor/
│  └─ processor.go
├─ usecase/
│  └─ usecase.go
├─ worker/
│  └─ worker.go
└─ template/                    # ต้นแบบสำหรับสร้าง module ใหม่
   ├─ handler.go / migration.sql / pg_repository.go / README.md
   ├─ redis_repository.go / usecase.go / worker.go
   ├─ delivery/http/handler.go / routes.go
   ├─ distributor/distributor.go
   ├─ models/model.go
   ├─ presenter/presenter.go
   ├─ processor/processor.go
   ├─ repository/pg_repository.go / redis_repository.go
   └─ usecase/usecase.go

pkg/                            # 21 shared packages
├─ cryptpass
├─ db
├─ elasticsearch
├─ emailTemplates
├─ helpers
├─ http-swagger
├─ httpErrors
├─ influxdb
├─ jwt
├─ kafka
├─ llm
├─ logger
├─ mqtt
├─ report
├─ responses
├─ secureRandom
├─ sendEmail
├─ transaction
├─ utils
├─ vectordb
└─ websocket
```

> หมายเหตุ: `internal/modules/*` เป็นกลุ่ม **แยกรราย** (ดูด้านล่าง)

---

## แยกรราย (โครงสร้างระดับ module — 33 modules)

### แบบที่ 1 — ลำดับชั้นเต็ม 6 ชั้น (alarm)
```
internal/modules/alarm/
├─ delivery/ + presenter/ + processor/ + models/
├─ repository/ + usecase/
└─ (handler / pg_repository / usecase อยู่ใต้ sub-folder ตามโครงสร้าง)
```

### แบบที่ 2 — 4 ชั้นแบบบาง (apimanager / control / flowengine)
```
internal/modules/<name>/
├─ delivery/ + presenter/
└─ repository/ + usecase/
```

### แบบที่ 3 — 5 ชั้น + provider (auditlog / notifier)
```
internal/modules/<name>/
├─ delivery/ + presenter/ + provider/
└─ repository/ + usecase/
```

### แบบที่ 4 — main file ที่ root (auth/batch/customer/dashboard/document/email/i18n/items/job/payment/settings/users/wos + ฯลฯ)
```
internal/modules/<name>/
├─ delivery/ + presenter/
├─ repository/ + usecase/
├─ handler.go
├─ pg_repository.go
└─ usecase.go
```

### แบบที่ 5 — จัดกลุ่มตาม service (elasticsearch / influxdb / mqtt / vectordata / realtime / websocket)
```
internal/modules/<name>/
├─ delivery/ + presenter/
└─ usecase/
```

### module พิเศษแบบละเอียด
```
internal/modules/fullschedule/
├─ delivery/ + executor/ + presenter/
├─ repository/ + scheduler/ + usecase/
└─ filter.go / handler.go / master.go / model.go / model_test.go
   pg_repository.go / usecase.go

internal/modules/iot/
├─ delivery/ + iothelper/ + models/ + presenter/ + repository/ + usecase/

internal/modules/kafka/
├─ delivery/ + models/ + repository/ + usecase/ + ws/

internal/modules/pdpa/
├─ application/ + domain/ + infrastructure/ + interfaces/ + repository/

internal/modules/queue/
├─ manager.go / manager_test.go / noop_queue.go / queue_test.go

internal/modules/report/
├─ delivery/ + usecase/ + company.go + handler.go

internal/modules/purchaseorder/  และ quotation/ และ wos/
├─ delivery/ + presenter/ + repository/ + usecase/
└─ handler.go / pg_repository.go / usecase.go

internal/modules/settings/
├─ delivery/ + presenter/ + repository/ + usecase/
└─ handler.go / repository.go / specs_alarm_logs.go / specs_device_schedule.go
   specs_integration.go / specs_master.go / spec_with.go / types.go / types_test.go / usecase.go

internal/modules/users/        # module อเนกประสงค์สุด
├─ delivery/ + distributor/ + doc/ + presenter/ + processor/
├─ repository/ + usecase/
└─ handler.go / pg_repository.go / pg_repository_test.go
   redis_repository.go / usecase.go / worker.go

internal/modules/websocket/
├─ delivery/ + models/ + presenter/ + repository/ + usecase/
└─ interface.go
```

### รายการ 33 modules (แยกตามโฟลเดอร์ `internal/modules/`)
```
alarm        apimanager   auditlog     auth         batch
control      customer     dashboard    document     elasticsearch
email        flowengine   fullschedule i18n         influxdb
iot          items        job          kafka        mqtt
notifier     payment      pdpa         purchaseorder queue
quotation    realtime     report       settings     users
vectordata   websocket    wos
```

---

## Template Prompt สำหรับสร้าง / แก้ไข Module (ใช้งานซ้ำได้)

เมื่อต้อง สร้าง module ใหม่ หรือ แก้ไข module ที่มีอยู่ ให้ใช้ prompt ต่อไปนี้

```
สร้าง/แก้ไขโมดูลใหม่ในโปรเจกต์ icmongolang

1. ตั้งชื่อโมดูล: <module_name>
2. ดูโครงสร้างที่ได้จาก docs\Template_Module.md 
3. สร้างโฟลเดอร์ที่จำเป็นภายใต้ internal\modules\<module_name>:
   - delivery\http\handler.go, routes.go
   - models\model.go
   - repository\pg_repository.go (และ redis_repository.go ถ้าต้องการใช้ Redis)
   - businesslogic\  (ถ้ามี logic แยก)
   - optional: processor\engine\session\template\ + es.go / influx.go / vectordb.go
4. ใช้ internal\modules\users และ internal\template เป็นตัวอย่างอ้างอิง
   รวมถึง internal\repository\ReadMe_Modules.md
5. ตรวจสอบว่า import package ใช้ชื่อจริงจาก pkg\ เท่านั้น:
   pkg\cryptpass, pkg\db, pkg\elasticsearch, pkg\emailTemplates, pkg\helpers,
   pkg\http-swagger, pkg\httpErrors, pkg\influxdb, pkg\jwt, pkg\kafka, pkg\llm,
   pkg\logger, pkg\mqtt, pkg\report, pkg\responses, pkg\secureRandom,
   pkg\sendEmail, pkg\transaction, pkg\utils, pkg\vectordb, pkg\websocket
   (ห้ามใช้ package ที่ไม่อยู่ในรายการนี้ - ตรวจสอบ live listing ด้วย)
6. ถ้าเป็น module ใหม่:
   - เพิ่ม schema/migration ไฟล์ใน migrations\ (ชื่อ YYYYMMDD_<module>_*.sql)
   - เส้นทาง route ลงทะเบียนใน internal\delivery\rest\router.go
7. ถ้าเป็น module เดิม: แก้ไขไฟล์เฉพาะในโครงสร้างเดิมของโมดูลนั้นเท่านั้น
   (ดู 'แยกรราย' ใน structure-map.md เพื่อหาแบบที่ตรงกับโมดูล)
8. สร้าง  Unit Test ตัวอย่าง
    internal\modules\<module_name>\myapp
            ├── go.mod
            ├── calculator.go
            └── calculator_test.go
9. รัน go build ./... และ/หรือ go test ./... ให้ผ่านก่อนจบ
```

---

> เอกสารนี้สร้างจาก live filesystem — หากโครงสร้างเปลี่ยน ให้ capture `Get-ChildItem -Recurse` ใหม่แล้วอัปเดตไฟล์นี้