# IoT-API-Sevice
ออกแบบ  ภาษาไทย พร้อมภาษาอังกฤษกำกับในบางส่วน
ออกแบบ ระบบ API  Golang นำไปใช้งานจริง

โครงสร้างหลัก ทั้ง โครงการ
### โฟลเดอร์หลัก  gobackend
```
gobackend/
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

# Architecture

In this project use 3 layer architecture

- Models
- Repository
- Usecase
- Delivery

## 🌟 Features

1.สร้างบทนำ
2.สร้างบทนิยาม
3.สร้างบทหัวข้อ
4.ออกแบบคู่มือ
5.TASK LIST CHECKLIST Template
6.Time line project Template





ข้อกำหนด:
1. แต่ละบทไม่จำกัดความยาว เน้นความสมบูนณ์ของเนื้อหา
- โครงสร้างการทำงาน
- ออกแบบ workflow
  - วาดรูป dataflow สร้าง รูปแบบ dataflow เหมือนจริง ลักษณะ flowchart   เพื่ออธิบายกระบวนการ ทำความเข้าใจ
  - พร้อมอธิบาย แบบ ละเอียด 
  - คอมเม้น code ภาษาไทย และ ภาษาอังถถษ อธิบาย การทำงาน แต่ละจุด
  - ยกตัวอย่างการใช้งานจริง หรือ กรณีศึกษา แนวทางแก้ไขปัญหา ที่อาจจะเกิดขึ้น
  - เทมเพลต และ ตัวอย่างโค้ด พร้อมนำไป run ได้ทันที  มีคำอธิบายการใช้งานแต่ละจุด การคอมเม้น  
- สรุป
   -ประโยชน์ที่ได้รับ
   -ข้อควรระวัง
   -ข้อดี
   -ข้อเสีย
   -ข้อห้าม ถ้ามี
2. ทุกบทต้องประกอบด้วย:
   - คำอธิบายแนวคิด (Concept Explanation)
   - ตัวอย่างโค้ดที่รันได้จริง (Runnable Code Example)
   - ตารางสรุป (ถ้ามีการเปรียบเทียบ)
   - แบบฝึกหัดท้ายบท 3–5 ข้อ
   - ส่วน "แหล่งอ้างอิง" ท้ายบท (References)
3. บทที่มีการออกแบบ Workflow, Task List, Checklist, Dataflow Diagram ให้:
   - แสดงเทมเพลตเป็น Markdown Table หรือลิงก์ดาวน์โหลด
   - อธิบายวิธีการใช้งานแต่ละจุด (step-by-step)
   - แทรกรูปภาพโดยระบุเป็น "รูปที่ X: คำอธิบาย"
4. สำหรับบทที่เกี่ยวข้องกับ Draw.io: ให้อธิบายวิธีการวาด Flowchart แบบ Top-to-Bottom (TB) พร้อมแสดงตัวอย่างโค้ด Mermaid หรือ ASCII flowchart
5. ใช้ภาษาไทยที่เป็นทางการ แต่เข้าใจง่ายและมีภาษอังถถษจุดสำคัญเสริม ไม่ใช้ศัพท์เทคนิคที่ซับซ้อนเกินไปโดยไม่มีการอธิบาย
5.หากใช้ศัพท์เทคนิค ต้องอธิบายความหมาย หลัการทำงาน วิธีการสำไปประยุตใช้   
ไม่จำกัดความยาว เน้นความสมบูนณ์ของเนื้อหา  มี สรุปสั้น ก่อน เนื้อหา แต่ละส่วน  มีหัวหนหัวข้อสำคัญ
คืออะไร
มีกี่แบบ
ใช้อย่างไร นำในกรณีไหน ทำไม่ต้องใช้ ประโยชน์ที่ได้รับ 
   -ประโยชน์ที่ได้รับ
   -ข้อควรระวัง
   -ข้อดี
   -ข้อเสีย
   -ข้อห้าม ถ้ามี 

  สร้าง : เทมเพลต และ ตัวอย่างโค้ด พร้อมนำไป run ได้ทันที  มีคำอธิบายการใช้งานแต่ละจุด การคอมเม้น ทุกส่วน 2ภาษา ไทย อังกถษ อออแบบระบบ พร้อมนำไป งานจริง  

  ทำให้ครบสุกส่วน แยก เล่ม เช่น เล่ม 1 ภาคทฤษฎี,  เล่ม 2 สถาบัตยกรรมโครงสร้างระบบ, เล่ม 3 การพัศนาระบบ การลงมือเขียน code นำไปใช้งานจิง , เล่ม 4 คู่การใช้งาน การบำรุงรักษา การขยายสวน การ ตั้งค่าระบบ ,    เพื่อสะดวกการใช้งาน



### 🚀 Core Features
- **🏭 Clean Architecture** - Three-layer architecture (Repository/Service/Handler) with comprehensive dependency injection
- **🔒 JWT Authentication** - Complete authentication system with access/refresh tokens and token blacklisting
- **👥 User Management** - Full CRUD operations with basic role-based protection (admin-only create)
- **📝 Structured Logging** - Unified logging via `pkg/logger` (built on Go's slog) with trace/request context
- **🚫 Rate Limiting** - IP-based request throttling with automatic cleanup
- **📊 Health Monitoring** - Comprehensive health checks with dependency monitoring
- **🌐 Redis Cache** - Production-ready caching layer with TTL management and object serialization
- **📦 Message Queue** - Redis-based pub/sub messaging with worker pools and dead letter queue support,kafka support,Raqbit mq,
- **👥 IOT Management**   mqtt ,socket io, snmp
- **👥 Token Management**  ลงทะเบียน token key สำหรับให้ผู้อื่น มาใชโดย กำหนดแต่ละ token ได้บ้าง manange role กำหนดวันหมดอายุได้เอง แต่ละ Pagekege token key มีสิทธิ์ต่างกัน, ระบบ Payment Token Key สำหรับ Page การใช้งาน
- **💼 Transaction Management** - GORM transaction manager with nested transaction support
- **🛡️ Security** - Multiple security layers including CORS, security headers, and input validation
### 🚀 Core database  PostgreSQL,postgresql,influxdb,TimescaleDB (TigerData),MongoDB 

### 🛠️ Middleware Stack
- **Request Context** - Trace IDs, request IDs, and user context propagation
- **Security Headers** - CSP, HSTS, X-Frame-Options, XSS Protection
- **CORS Handling** - Configurable cross-origin resource sharing
- **Panic Recovery** - Application-level panic handling with graceful error responses
- **Request Logging** - Structured request/response logging with performance metrics
- **Authentication** - JWT middleware with role-based route protection
- **Input Validation** - Comprehensive request validation using go-playground/validator

### 📈 Health & Monitoring
- **Health Endpoints** - Basic, detailed, readiness, and liveness probes
- **Dependency Checks** - Database and Redis connection status via `/health/detailed`
- **System Metrics (Optional)** - Runtime/memory snapshot handler available for wiring
- **Performance Tracking (Optional)** - Monitoring middleware + metrics handler available
- **Kubernetes Ready** - `/ready` and `/live` probes

- CRUD
- Jwt, refresh token saved in redis
- Cached user in redis
- Email verification
- Forget/reset password, send email

## Technical

- `chi`: router and middleware
- `viper`: configuration
- `cobra`: CLI features
- `gorm`: orm
- `validator`: data validation
- `jwt`: jwt authentication
- `zap`: logger
- `gomail`: email
- `hermes`: generate email body
- `air`: hot-reload

## Start Application

### Generate the Private and Public Keys

- Generate the private and public keys: [travistidwell.com/jsencrypt/demo/](https://travistidwell.com/jsencrypt/demo/)
- Copy the generated private key and visit this Base64 encoding website to convert it to base64
- Copy the base64 encoded key and add it to the `config/config-local.yml` file as `jwt`
- Similar for public key

### Stmp mail config

- Create [mailtrap](https://mailtrap.io/) account
- Create new inboxes
- Update smtp config `config/config-local.yml` file as `smtpEmail`



เริ่มเขียนได้ทันที โดยเริ่มจากบทที่ ล่าสุด ไปจนครบ
ข้อกำหนด:
1. แต่ละบทไม่จำกัดความยาว เน้นความสมบูนณ์ของเนื้อหา
- โครงสร้างการทำงาน
- ออกแบบ workflow
  - วาดรูป dataflow สร้าง รูปแบบ dataflow เหมือนจริง ลักษณะ flowchart   เพื่ออธิบายกระบวนการ ทำความเข้าใจ
  - พร้อมอธิบาย แบบ ละเอียด 
  - คอมเม้น code ภาษาไทย และ ภาษาอังถถษ อธิบาย การทำงาน แต่ละจุด
  - ยกตัวอย่างการใช้งานจริง หรือ กรณีศึกษา แนวทางแก้ไขปัญหา ที่อาจจะเกิดขึ้น
  - เทมเพลต และ ตัวอย่างโค้ด พร้อมนำไป run ได้ทันที  มีคำอธิบายการใช้งานแต่ละจุด การคอมเม้น  
- สรุป
   -ประโยชน์ที่ได้รับ
   -ข้อควรระวัง
   -ข้อดี
   -ข้อเสีย
   -ข้อห้าม ถ้ามี
2. ทุกบทต้องประกอบด้วย:
   - คำอธิบายแนวคิด (Concept Explanation)
   - ตัวอย่างโค้ดที่รันได้จริง (Runnable Code Example)
   - ตารางสรุป (ถ้ามีการเปรียบเทียบ)
   - แบบฝึกหัดท้ายบท 3–5 ข้อ
   - ส่วน "แหล่งอ้างอิง" ท้ายบท (References)
3. บทที่มีการออกแบบ การออกแบบ workflow หรือการสร้าง diagram), Task List, Checklist,  จนครบทุกบท ให้:
   - แสดงเทมเพลตเป็น Markdown Table หรือลิงก์ดาวน์โหลด
   - อธิบายวิธีการใช้งานแต่ละจุด (step-by-step)
   - แทรกรูปภาพโดยระบุเป็น "รูปที่ X: คำอธิบาย"
4. สำหรับบทที่เกี่ยวข้องกับ Draw.io: ให้อธิบายวิธีการวาด Flowchart แบบ Top-to-Bottom (TB) พร้อมแสดงตัวอย่างโค้ด Mermaid หรือ ASCII flowchart
5. ใช้ภาษาไทยที่เป็นทางการ แต่เข้าใจง่ายและมีภาษอังถถษจุดสำคัญเสริม ไม่ใช้ศัพท์เทคนิคที่ซับซ้อนเกินไปโดยไม่มีการอธิบาย
5.หากใช้ศัพท์เทคนิค ต้องอธิบายความหมาย หลัการทำงาน วิธีการสำไปประยุตใช้   
ไม่จำกัดความยาว เน้นความสมบูนณ์ของเนื้อหา  มี สรุปสั้น ก่อน เนื้อหา แต่ละส่วน  มีหัวหนหัวข้อสำคัญ
คืออะไร
มีกี่แบบ
ใช้อย่างไร นำในกรณีไหน ทำไม่ต้องใช้ ประโยชน์ที่ได้รับ 
   -ประโยชน์ที่ได้รับ
   -ข้อควรระวัง
   -ข้อดี
   -ข้อเสีย
   -ข้อห้าม ถ้ามี
CMON IoT Solution
1.	Monitoring or auto control and management and notification system 
(ระบบติดตามเฝ้าระวังภัย และ ควบคุม บริการ จัดการอุปกรณ์ อัตโนมัติและการแจ้งเตือน ภัย)
 
. ภาพรวมระบบ (System Overview)
ระบบนี้เป็นแพลตฟอร์มกลางสำหรับการตรวจสอบสภาพแวดล้อมและอุปกรณ์ภายในห้อง Data Center และสถานที่สำคัญแบบ Real-time โดยอัตโนมัติ ระบบสามารถรับข้อมูลจากเซนเซอร์และอุปกรณ์ได้หลายประเภท
 วิเคราะห์ข้อมูลตามเงื่อนไขที่กำหนด ส่งการแจ้งเตือนผ่านหลายช่องทาง และสั่งงานควบคุมอุปกรณ์อื่นๆ ได้โดยอัตโนมัติ พร้อมทั้งมีระบบรายงานสำหรับการวิเคราะห์เชิงลึก
เป้าหมายหลัก:
•	ลด Downtime โดยการตรวจจับความผิดปกติตั้งแต่เริ่มต้น
•	เพิ่มประสิทธิภาพ การทำงานของเจ้าหน้าที่ผ่านระบบอัตโนมัติ มีรายงานประวัติการทำงานของระบบ มี dashboard แสดงตาม ความต้องการ
•	รักษาความปลอดภัย ของโครงสร้างพื้นฐานทางไอทีอย่างต่อเนื่อง
•	เป็นศูนย์กลางการจัดการ (Single Pane of Glass) สำหรับการตรวจสอบทั้งหมด
•	Schedule ตั้งค่า กำหนดการสั่งงานอุปกรณ์ ตามตารางเวลา ตามที่ได้มีการตั้งค่าไว้
•	แจ้งเตียน และสั่งงานอุปกรณ์ สำรอง แบบ อัตโนมัติ ตามที่ได้มีการตั้งค่าไว้
2.	สถาปัตยกรรมระบบ (System Architecture) แบบ Modular ที่สามารถขยายได้
คุณสมบัติหลัก (Core Features)
3.1 ระบบแจ้งเตือนและเฝ้าระวังภัยแบบ Real-time
•	ความเร็ว: ข้อมูลอัพเดทแบบ Real-time ด้วยความคลาดเคลื่อนไม่เกิน 45 วินาที โดยขึ้นอยู่กับสภาพเครือข่าย
•	การแจ้งเตือนสองรูปแบบ:
o	รูปแบบที่ 1: แจ้งเตือนผ่านระบบ
	Web Dashboard: แสดงการแจ้งเตือนแบบ Pop-up, Notification Bell ด้วยสีตามระดับความรุนแรง
	อีเมล: ส่งอีเมลอัตโนมัติพร้อมรายละเอียดเหตุการณ์
	Mobile Application/Push Notification: แจ้งเตือนไปยังแอปบนมือถือ
	Line Application / SMS: ส่งข้อความสั้นๆ พร้อมลิงก์สำหรับดูรายละเอียดเพิ่มเติม
o	รูปแบบที่ 2: แจ้งเตือนผ่านอุปกรณ์
	สัญญาณไฟ: ไฟสีเขียว(ปกติ), สีเหลือง(Warning), สีแดง(Alarm)
	สัญญาณเสียง: Siren หรือเสียงเตือนที่ปรับระดับความดังได้
	การเชื่อมต่อกับระบบภายในอาคาร (NO/NC): ใช้รีเลย์ Dry Contact เพื่อเปิด/ปิดวงจรในระบบอื่นๆ ของอาคารได้ตามความต้องการ
3.2 ระบบเฝ้าระวังและแจ้งเตือนอัตโนมัติ
•	ระดับความรุนแรงของการแจ้งเตือน:
o	Normal (ปกติ): สภาพแวดล้อมและอุปกรณ์ทำงานปกติ
o	Warning (ระดับกลาง): ค่าใกล้เคียงขีดจำกัด หรือพบความผิดปกติเล็กน้อย (เช่น อุณหภูมิเริ่มสูง) แจ้งเตือนเพื่อเตรียมการ
o	Alarm (ระดับร้ายแรง): ค่าเกินขีดจำกัดที่กำหนด หรือพบความผิดปกติที่ส่งผลกระทบทันที (เช่น ตรวจพบน้ำรั่ว, ไฟดับ) แจ้งเตือนทันทีและอาจสั่งงานอัตโนมัติ
•	การตรวจสอบสภาพแวดล้อม:
o	อุณหภูมิและความชื้น (Temperature/Humidity): ตั้งค่า Threshold สำหรับการแจ้งเตือน
o	การรั่วซึมของน้ำ (Water Leak): ตรวจจับการรั่วไหลของน้ำ
o	ควันและไฟ (Smoke/Fire): ตรวจจับควันหรือสัญญาณจากระบบดับเพลิง
o	ระบบไฟฟ้าและ UPS: ตรวจสอบสถานะ แบตเตอรี่ โหลด ของ UPS ผ่าน SNMP
3.3 ระบบจัดการและควบคุมอัตโนมัติ
•	การตั้งค่าเงื่อนไข (Rule Engine): ผู้ใช้สามารถกำหนด "If-This-Then-That" ได้
o	ตัวอย่าง: IF อุณหภูมิ > 35°C THEN ส่ง Email แจ้งเตือนระดับ Alarm AND เปิดสัญญาณไฟแดง AND ส่งคำสั่ง MQTT ไปเพิ่มความเย็นของเครื่องปรับอากาศ
•	การตั้งเวลางาน (Scheduler):
o	ตั้งเวลาส่งรายงานอัตโนมัติ (รายวัน, สัปดาห์, เดือน, ปี)
o	ตั้งเวลาเปิด/ปิดอุปกรณ์ที่ไม่สำคัญตามเวลาทำงาน
3.4 ระบบรายงานและแดชบอร์ด
•	แดชบอร์ดแบบกราฟิก (Graphic Interface):
o	แสดงสถานะอุปกรณ์และเซนเซอร์แบบ Real-time ด้วยไอคอนและสี
o	แสดงข้อมูลเป็นกราฟเส้น (Line Chart) สำหรับแนวโน้ม และกราฟแท่ง (Bar Chart) สำหรับเปรียบเทียบ
•	ระบบรายงาน (Reporting):
o	รูปแบบ: ตาราง, กราฟเส้น, กราฟแท่ง
o	ประเภท: รายงานประจำวัน, สัปดาห์, เดือน, ปี
o	เนื้อหา: สรุปเหตุการณ์แจ้งเตือน, สถิติการทำงานของอุปกรณ์, แนวโน้มอุณหภูมิ/ความชื้น
o	การส่งออก: ดูบนเว็บ, ดาวน์โหลดเป็น PDF/Excel, ส่งอีเมลอัตโนมัติตามตารางเวลาที่ตั้งไว้
4. การรองรับการเชื่อมต่ออุปกรณ์
ระบบรองรับการเชื่อมต่อกับอุปกรณ์ที่หลากหลายผ่านโปรโตคอลมาตรฐาน:
1.	SNMP (Simple Network Management Protocol):
o	สำหรับอุปกรณ์เครือข่ายและไฟฟ้า: เครื่องปรับอากาศในห้อง server, UPS, PDU, Switch
2.	MQTT (The Standard for IoT Messaging):
o	สำหรับเซนเซอร์ IoT: เซนเซอร์วัดอุณหภูมิ/ความชื้น, ความเข้มแสง, ตรวจจับน้ำรั่ว, วัดกระแสไฟฟ้า
3.	อุปกรณ์และเซนเซอร์แบบกำหนดเอง (Custom Sensors):
o	เช่น: Sensor RFID, TEME Sensor (วัดความชื้น, อุณหภูมิ, ความเข้มแสง)
o	การเชื่อมต่อ: ผ่าน WiFi หรือ LAN โดยตรงไปยังระบบ
4.	การเชื่อมต่อแบบดิจิทัล (Digital I/O):
o	Digital Input (DI): รับสัญญาณจากสวิตช์หรือเซนเซอร์แบบง่าย (เช่น Door Contact Sensor)
o	Digital Output (DO): ใช้ส่งสัญญาณควบคุม (เปิด/ปิด) อุปกรณ์ภายนอกผ่าน Contact (NO/NC) เช่น ควบคุมไฟ, เปิดปั๊มน้ำ
________________________________________
5. ประโยชน์ของระบบ
•	รู้ทันปัญหา: รับรู้ความผิดปกติของอุปกรณ์และสภาพแวดล้อมได้ทันท่วงที
•	ป้องกันความเสียหายร้ายแรง: สามารถแก้ไขปัญหาได้ก่อนจะลุกลาม
•	ลดภาระงาน manual: ระบบอัตโนมัติและการแจ้งเตือนที่ชาญฉลาดช่วยให้เจ้าหน้าที่ทำงานมีประสิทธิภาพมากขึ้น
•	ตัดสินใจได้ดีขึ้น: ข้อมูลและรายงานแบบกราฟิกที่เข้าใจง่ายช่วยในการวิเคราะห์และวางแผน
•	บริหารจัดการจากที่ใดก็ได้: ผ่าน Web Dashboard และ Mobile Application
•	เป็นระบบเปิด (Open System): รองรับการเชื่อมต่อกับอุปกรณ์และระบบใหม่ๆ ในอนาคตได้ง่าย
________________________________________
6. ตัวอย่างการทำงาน (Use Case Scenario)
สถานการณ์: อุณหภูมิในห้อง Server สูงเกินกำหนด
1.	ตรวจจับ: เซนเซอร์วัดอุณหภูมิ (เชื่อม via MQTT) ส่งค่าอุณหภูมิ 38°C เข้าสู่ระบบ
2.	วิเคราะห์: Rule Engine ตรวจพบว่าค่าเกิน Threshold ที่ตั้งไว้ (35°C) และประเมินเป็นระดับ "Alarm"
3.	ดำเนินการ:
o	แจ้งเตือน: ระบบส่งการแจ้งเตือนแบบ Real-time ไปยัง:
	Dashboard: เปลี่ยนสีไอคอนเป็นแดงและแสดง Pop-up
	Email & Line: ส่งข้อความ "ALARM: อุณหภูมิใน Rack A1 สูงถึง 38°C"
	อุปกรณ์: เปิดสัญญาณไฟแดงและเสียงไซเรนในห้อง
o	ควบคุมอัตโนมัติ (หากตั้งค่า): ระบบส่งคำสั่ง Digital Output (DO) ไปเปิดเครื่องปรับอากาศเครื่องสำรองทันที
4.	ติดตามผล: เจ้าหน้าที่รับการแจ้งเตือนและเข้าตรวจสอบได้ทันที ในขณะที่ระบบพยายามแก้ไขเบื้องต้นโดยอัตโนมัติแล้ว
5.	รายงาน: ระบบบันทึกเหตุการณ์ทั้งหมดไว้สำหรับการสร้างรายงานสรุปประจำสัปดาห์




-คุณสมบัติ
-ระบบแจ้งเตือนและเฝ้าระวังภัย ทำงานของระบบ แบบ Realtime  
ความคลาดเคลื่อนไม่เกิน 45 วินาที ตามภาวะ การ ของ ระบบ Network
-ระบบการแจ้งเตือนจะแยกออกเป็น 2 รูปแบบ
-รูปแบบที่ 1 คือแจ้งเตือนผ่านระบบ 
-เช่นสามารถแจ้งเตือนผ่านหน้า dashboard,สามารถแจ้งเตือนผ่านอีเมล,สามารถแจ้งเตือนผ่าน application Email เป็นต้น
-รูปแบบที่ 2 คือการแจ้งเตือนผ่านอุปกรณ์
-เช่นสามารถแจ้งเตือนจากสัญญาณไฟ,สามารถแจ้งเตือนผ่านสัญญาณเสียง,สามารถแจ้งเตือนเข้ากับระบบภายในอาคารโดยผ่าน NO/NC  
-NO (Normally Open) NC (Normally Closed) Sensor ตามความต้องการของลูกค้า เป็นต้น
- ระบบเฝ้าระวังและแจ้งเตือนอัตโนมัติ
ยกตัวอย่าง รูปแบบการแจ้งเตือน เช่น ระดับความรุนแรง ของอุปกรณ์ ที่ระบบตรวจพบความผิดปกติในการทำงาน    1. Warning เตือนภัยระดับกลาง     2. Alarm เตือนภัยระดับร้ายแรง เป็นต้น 
Temperature/Humidity Monitoring
โซลูชั่นสำหรับวัดอุณหภูมิและความชื้น หากอุณภูมิหรือความชื้นเกินหรือต่ำกว่าค่าที่เรากำหนดสามารถให้ส่งการแจ้งเตือนผ่าน Email, Line, SMS ได้ โดยบางยี่ห้อจะมีระบบ Cloud เพื่อจัดการกับอุปกรณ์, Report, และแจ้งเตือนในอีกช่องทาง
ระบบเฝ้าระวังและแจ้งเตือนอัตโนมัติ สามารถออกรายงานและส่งอีเมลอัตโนมัติ ได้ทั้งรายวัน รายสัปดาห์ รายเดือนและรายปี  ตามที่มีการตั้งค่าไว้
ประโยชน์ของการใช้ระบบ 
เพื่อให้เจ้าหน้าที่ที่รับผิดชอบภายในห้อง Data Center ทราบสถานะความผิดปกติของอุปกรณ์ และสามารถแก้ไขปัญหาที่เกิดขึ้นก่อนที่ปัญหาจะลุกลามเกินกว่าจะแก้ไขเพื่อง่ายต่อการแสดงผลในการตรวจสอบเนื่องจากแสดงผลในรูปแบบ Graphic ระบบ มี Sensor ในการตรวจสอบระบบที่หลากหลาย ระบบตรวจสอบสภาพแวดล้อมภายในห้อง Data Center  เป็นระบบที่มีความสำคัญสำหรับหน่วยงานที่มีห้อง Server  และห้อง Data Center   โดยระบบ  จะช่วยตรวจสอบและเฝ้าระวังการทำงานของระบบต่างๆภายให้ห้อง Server, Data Center เช่น ระบบตรวจจับอุณหภูมิ และ ความชื้นภายในห้อง,  ระบบตรวจจับน้ำรั่วซึม, ระบบตรวจจับควันไฟ, ระบบไฟฟ้าและระบบสำรองไฟฟ้า (UPS)  โดยเซนเซอร์ต่างๆ ที่ติดตั้งในห้อง Server, ห้อง Data Center จะส่งสัญญาณแจ้งเตือนไปยังระบบตรวจสอบสภาพสิ่งแวดล้อม และกำหนดให้มีการทำงานโดยการส่ง Email, SMS, หรือส่งข้อความแจ้งเตือนไปยัง Line Application บนโทรศัพท์มือถือได้อีกด้วย
           โดยระบบ  นั้นมีความสามารถในการเชื่อมต่อเช้ากับระบบต่างๆผ่าน Digital Input และระบบอื่น  ซึ่งสามารถทำการสั่งให้ระบบ Digital Output ไปควบคุมให้ระบบต่างๆทำงาน และสามารถวัดค่าผ่านระบบ Sensor หรือผ่านระบบ MQTT และแสดงผลเป็น Graphic Interface
รองรับการ เชื่อต่อกับอุปกรณ์ ต่างๆ 
1.อุปกรณ์ ที่รองระบ สัญญาณ SNMP - Simple Network Management Protocol
-เช่น เครื่องปรับอากาศ(air conditioner)   เครื่องสำรองไฟฟฟ้า (UPS) เป็นต้น
2.อุปกรณ์  IOT  ที่รองระบ สัญญาณ  MQTT - The Standard for IoT Messaging
--- MQTT Protocol | Messaging & Data Exchange for the IoT
เช่น  sensor วัดความชื่น อุณหภูมิ   ความเข้มแสง  ความถี่ของคลื่นเสียง (sound wave frequency) 
กระแสไฟฟ้า  เป็นต้น  
3.Sensor ที่ทางทีมงานพัฒนาขึ้น เช่น Sensor RFID, TEME Sensor วัดความชื่น อุณหภูมิ   ความเข้มแสง  
-การรองรับ กับอุปกรณ์
     ไร้สาย  WIFI  / LAN
-การทำงาน
 --แบบ อัตโนมัติ ตั้งค่าการ สั่งงานตาม เวลาในปฏิทิน หรือส่ง รายงาน ได้ ตามวันและเวลา 
ในแต่ละช่วงเวลา จากการตั้งค่าในระบบ  1.กำหนด วัน เวลา สั่งงาน 2.เลือกอุปกรณ์ที่ต้องการสั่งให้ทำงาน
- แสดงรายงาน แบบ ตาราง แบบการแท่ง
- แสดงรายงาน แบบ ตาราง แบบกราฟเส้น
Benefit  (ประโยชน์ที่ได้)
- สามารถติดตามการเปลี่ยนแปลงของอุณหภูมิและความชื้น และการทำงานของอุปกรณ์ IoT ได้ตลอด 24 ชั่วโมง
- สำรองข้อมูลและสามารถเรียกดูรีพอร์ตย้อนหลังได้
- ลดภาระการทำงานของเจ้าหน้าที่ ในการจดบันทึกข้อมูล
- มีข้อมูลทั้งอุณหภูมิและความชื้น (2in1) ที่น่าเชื่อถือ
ระบบ จะ แสดงการแจ้งเตียนแบบเรียลไทม์(Real Time Alert)ผ่านทางหน้า จอ และส่งการแจ้งเตียน ไปทาง email 
และสั่งอุปกรณ์สำรองทำงานทดแทนอุปกรณ์ ที่เสียหาย ทันที ที่เกิด เหตุการณ์  
The system will display real-time alerts on the screen and send notifications via email. It will also immediately dispatch backup equipment to replace damaged equipment upon an incident.
สามารถตั้งค่าการ แจ้งเตียน ทาง Email, SMS และ สั่งอุปกรณ์สำรองทำงานทดแทนอุปกรณ์ ที่เสียหาย ที่ตั้งค่าไว้ได้
You can set up notifications via Email, SMS and order backup devices to replace damaged devices that have been set up.
สั่งอุปกรณ์สำรองทำงานทดแทนอุปกรณ์ ที่เสียหาย ที่ตั้งค่าไว้ได้
Order backup devices to work as replacements for damaged devices that have been set up.
 
ระบบ ออกแบบให้ใช้งาน  ง่ายโดย ลดขั้นตอนที่ซ้ำซ้อนออกไป เพื่อให้เกิดความสะดวก  ที่สุด
The system is designed to be easy to use by reducing redundant steps for maximum convenience.
 
ระบบ จะมีหน้า แสดงผลการทำงาน ดูย้องหลังได้ ออกแบบให้ใช้งายง่ายแสดง ข้อมูลเท่าที่จำเป็น
The system will have a page that shows the results of work and allows for retrospective viewing. It is designed to be easy to use and shows only the necessary information.
 
หน้ารายงานสามารถ กรอง ข้อมูลใน ช่วงเวลา หรือค้นหาหัวข้อที่สำคัญได้ ใน มากี่คลิก 
The report page can filter data by time period or search for important topics in just a few clicks.
ระบบ จะทำการ ควบคุมอุปกรณ์  แบบอัตโนมัติ ทันที่ที่มีเหตุการณ์ เช่น เปิด หรือ ปิด พัดลม เมื่ออุณหภูมิ ห้อง สูงเกิกำนด
และสั่ง ปิด /เปิด เมื่อ อุณหภูมิกลับสู่สภาวะปกติ  และส่งแจ้งเตียนไป ทาง  Email / SMS  
The system automatically controls devices immediately upon an event, such as turning the fan on or off when the room temperature exceeds a certain level.
It also automatically turns the unit off or on when the temperature returns to normal. and send notifications via Email / SMS
Schedule work on  by time period (กำหนดตารางการทำงานตามช่วงเวลา)
 
เมื่อ ถึงกำนด วันเวลา ทำงานตามช่วงเวลา ระบบสั่งอุปกรณ์ทำงาน ส่งรายงานการทำงานไป ทาง  Email / SMS
When the time is set to work according to the time period, the system will order the equipment to work and send the work report via Email / SMS.
-สามารถสั่งอุปกรณ์ให้ทำงาน ได้ ผ่านหน้าจอควบคุมอุปกรณ์ ทันที เช่นสั่งเปิด ปิด พัดลม ระบายอากาศ
-Can command devices to work immediately through the device control screen, such as commanding to turn on/off ventilation fans.
-สามารถสั่งอุปกรณ์ให้ทำงาน ได้ ผ่านหน้าจอควบคุมอุปกรณ์ ทันที  เช่นสั่งเปิด ปิด แอร์
- Can command devices to work immediately through the device control screen, such as turning the air conditioner on and off.
การเพิ่มอุปกรณ์หรือเปลี่ยนอุกรกรณ์ สามารถตั้งค่า บนหน้าจอจะทำได้ง่าย ลดความซับซ้อน เพื่อปังกันความผิดพลาด
Adding or replacing devices can be easily configured on the screen, reducing complexity to prevent errors.
3.Smart building management (การจัดการอาคารอัจฉริยะ)
rprise Application + Database + Cache + Message Queue" 
 