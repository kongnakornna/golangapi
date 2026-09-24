# Module xxx: xxx/xxx 

## สำหรับโฟลเดอร์ `xxx/xxx/`

ไฟล์ที่เกี่ยวข้อง:
 
---
# การคอมเม้น ใช้ 2 ภาษา อังกถษ และ ภาษาไทย คนละบรรทัด
## หลักการ (Concept)

###  คืออะไร?
 
### มีกี่แบบ?  

**ข้อห้ามสำคัญ:**  

### ใช้อย่างไร / นำไปใช้กรณีไหน

 
### ประโยชน์ที่ได้รับ
 
### ข้อควรระวัง

 
### ข้อดี
 
### ข้อเสีย
 
### ข้อห้าม
 

## การออกแบบ Workflow และ Dataflow

---

## ตัวอย่างโค้ดที่รันได้จริง
 
---

## วิธีใช้งาน module นี้
 
---
## การติดตั้ง
 
---
## การตังค่า configuration

---
## การรวมกับ GROM

---

# design file  table sql ที่เกียวข้อง

# retrun เป็น RestAPI


## การใช้งานจริง
 
---

## ตารางสรุป   Components
 
---

## แบบฝึกหัดท้าย module (5 ข้อ)
 
---

## แหล่งอ้างอิง
 
---

**หมายเหตุ:** module นี้ครบถ้วนสำหรับ `xxx` สำหรับระบบ ccc  หากต้องการ module เพิ่มเติม (เช่น cc`)  



### โฟลเดอร์หลัก  icmongolang
```
icmongolang/
├── .vscode/
│   ├── launch.json
│   └── settings.json
├── cmd/
│   ├── api/
│   │   └── main.go
│   ├── initdata.go
│   ├── migrate.go
│   ├── root.go
│   ├── serve.go
│   └── worker.go
├── config/
│   ├── config-local.yml
│   ├── config-prod.yml
│   └── config.go
├── docdev/
├── docs/
├── internal/
│   ├── models/
│   │   ├── base.go
│   │   ├── session.go
│   │   ├── user.go
│   │   └── verification.go
│   ├── repository/
│   │   ├── pg_repository.go
│   │   ├── redis_repo.go
│   │   ├── session_repo.go
│   │   └── user_repo.go
│   ├── usecase/
│   │   ├── auth_usecase.go
│   │   ├── cache_usecase.go
│   │   └── user_usecase.go
│   ├── delivery/
│   │   ├── rest/
│   │   │   ├── handler/
│   │   │   │   ├── auth_handler.go
│   │   │   │   ├── health_handler.go
│   │   │   │   └── user_handler.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go
│   │   │   │   ├── cors.go
│   │   │   │   ├── logger.go
│   │   │   │   ├── monitoring.go
│   │   │   │   ├── rate_limit.go
│   │   │   │   └── security.go
│   │   │   ├── dto/
│   │   │   │   ├── auth_dto.go
│   │   │   │   ├── error_dto.go
│   │   │   │   └── user_dto.go
│   │   │   └── router.go
│   │   └── worker/
│   │       └── email_worker.go
│   └── pkg/
│       ├── email/
│       │   ├── gomail_sender.go
│       │   ├── sender.go
│       │   └── templates/
│       │       ├── reset_password.html
│       │       └── verification.html
│       ├── hash/
│       │   └── bcrypt.go
│       ├── jwt/
│       │   ├── maker.go
│       │   ├── payload.go
│       │   └── rsa_maker.go
│       ├── logger/
│       │   └── zap_logger.go
│       ├── redis/
│       │   ├── cache.go
│       │   ├── client.go
│       │   └── refresh_store.go
│       ├── utils/
│       │   ├── random.go
│       │   └── time.go
│       └── validator/
│           └── custom_validator.go
├── migrations/
│   ├── 000001_create_users_table.down.sql
│   └── 000001_create_users_table.up.sql
├── pkg/
│   └── utils/
├── scripts/
│   ├── build.sh
│   └── deploy.sh
├── vendor/
├── .air.toml
├── .dockerignore
├── .env.dev
├── .env.prod
├── .gitignore
├── docker-compose.dev.yml
├── docker-compose.prod.yml
├── Dockerfile.dev
├── Dockerfile.prod
├── go.mod
├── go.sum
├── LICENSE
├── README.md
└── BookGolang.md
```


ใช้ Template ด้านบน   golang  

ดำเนินการต่อไปโดยอัตโนมัติ