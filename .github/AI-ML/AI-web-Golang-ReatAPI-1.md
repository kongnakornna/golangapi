# เอกสารประกอบการอบรม Go Bootcamp: สถาปัตยกรรมระบบแบบ Clean Architecture และการประมวลผลแบบ Real-Time
## Go Bootcamp: With gRPC and Protocol Buffers (HTTP/S, HTTP2)
## Architecture

In this project use 3 layer architecture

- Models
- Repository
- Usecase
- Delivery

## Features

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

```
# ตัวอย่าง โครงสร้าง Folder

```bash
  api/
    ├── cmd/
    │   ├── apiser/                 # REST API หลัก (มีอยู่แล้ว)
    │   ├── websocket/              # *** WebSocket server
    │   │   └── main.go
    │   ├── initdata.go
    │   ├── root.go
    │   ├── serve.go
    │   └── worker.go
    ├── internal/
    │   ├── websocket/              # **** โค้ดเฉพาะของ WebSocket
    │   │   ├── delivery/
    │   │   │   └── ws/
    │   │   │       ├── hub.go
    │   │   │       ├── client.go
    │   │   │       └── handler.go     # HTTP endpoint สำหรับ upgrade
    │   │   ├── usecase/
    │   │   │   └── ws_usecase.go      # business logic (save message, auth)
    │   │   ├── repository/
    │   │   │   └── ws_repo.go         # interface สำหรับ DB
    │   │   └── models/
    │   │       └── ws_models.go       # entity ของ message, session
    │   ├── pkg/                    # shared packages (มีอยู่แล้ว)
    │   │   ├── websocket/          # **ปรับปรุง** ใช้ร่วมกันได้
    │   │   │   ├── hub.go          # core Hub logic
    │   │   │   ├── client.go
    │   │   │   └── message.go      # struct ของ message
    │   │   └── ...
    ├── migrations/                 # **เพิ่ม** SQL schema สำหรับ websocket
    │   └── 20250619_websocket_tables.sql
    └── 
```
```bash
    icmongolang/
    ├── pkg/
    │   ├── helpers/
    │   │   ├── iot.go          # Alarm logic (สมบูรณ์)
    │   │   └── format.go       # ฟังก์ชันช่วยเหลือ (time, string, random)
    │   ├── mqtt/
    │   │   └── client.go       # MQTT client พร้อม GetDataFromTopic
    │   ├── influxdb/
    │   │   └── client.go       # InfluxDB client
    │   └── redis/
    │       └── redis_conn.go   # Redis client + Cache interface
    ├── internal/
    │   ├── mqtt/
    │   │   ├── delivery/http/
    │   │   │   ├── handler.go
    │   │   │   └── routes.go
    │   │   ├── presenter/
    │   │   │   └── presenter.go
    │   │   └── usecase/
    │   │       └── usecase.go
    │   ├── influxdb/
    │   │   ├── delivery/http/
    │   │   │   ├── handler.go
    │   │   │   └── routes.go
    │   │   ├── presenter/
    │   │   │   └── presenter.go
    │   │   └── usecase/
    │   │       └── usecase.go
    │   ├── alarm/
    │   │   ├── delivery/http/
    │   │   │   ├── handler.go
    │   │   │   └── routes.go
    │   │   ├── repository/
    │   │   │   └── alarm_log_repo.go
    │   │   └── usecase/
    │   │       └── usecase.go
    │   └── server/
    │       ├── handlers.go
    │       └── server.go
    └── cmd/api/main.go (สมมติตามเดิม)
```

ส่วน แผนการสอน 
	1.วัตุประสงค์
	2.กลุ่มเป้าหมาย
	3.ความรู้พื้นฐาน
	4.เนื้อหา โดยย่อ กระชับ เน้น วัตถุประสงค์  ประโยชน์ของการใช้
ส่วน เอกสาร
	1.สร้างบทนำ
	2.สร้างบทนิยาม
	3.สร้างบทหัวข้อ
	5.ออกแบบคู่มือ
	6.ออกแบบ workflow
	7.TASK LIST Template
	8.CHECKLIST Template
	9.สรุป
1.บทนำ
2.บทนิยาม ศัพท์ ทุกส่วน
 เช่น 
	Queue Processor คืออะไร
	Queue Processor มีกี่แบบ
	Queue Processor ใช้อย่างไร นำในกรณีไหน ทำไม่ต้องใช้ ประโยชน์ที่ได้รับ
3.ออกแบบ workflow
  - วาดรูป dataflow สร้างรูปแบบ draw.io เหมือนจริง ลักษณะ flowchart TB   เพื่ออธิบายกระบวนการ ทำความเข้าใจ
  - พร้อมอธิบาย แบบ ละเอียด 
  - ยกตัวอย่างการใช้งานจริง หรือ กรณีศึกษา 
  - เทมเพลตและตัวอย่างโค้ด พร้อมนำไป run ได้ทันที  มีคำอธิบายการใช้งานแต่ละจุด การคอมเม้น ภาษาไทย และ ภาษาอังกถษ
4.ทำบทสรุปท้ายบท
   -ประโยชน์ที่ได้รับ
   -ข้อควรระวัง
   -ข้อดี
   -ข้อเสีย
   -ข้อห้าม ถ้ามี
   -แหล่ง อ้างอิ่ง ที่มา 
5.คู่มือ
   -TASK LIST Template
    -กระบวนการทำงานแต่ละขั้นตอน
	-กระบวนการตรวจสอบ การทำงานแต่ละขั้นตอน
	-กระบวนการ อธิบายปัญหา สรุป 
	   - Root Cause Analysis (RCA) คือแนวทางหรือกระบวนการวิเคราะห์ปัญหาอย่างเป็นระบบ เพื่อค้นหา "สาเหตุที่แท้จริง" (Root Cause) ของปัญหาที่เกิดขึ้น 
			แทนที่จะแก้ไขเพียงอาการ (Symptoms) ที่ปรากฏ ช่วยป้องกันปัญหาไม่ให้เกิดซ้ำอย่างยั่งยืน 
			โดยนิยมใช้เครื่องมืออย่างเทคนิค 5 Why หรือแผนภูมิก้างปลามาช่วยในการวิเคราะห์ 
			แนวทางและขั้นตอนการวิเคราะห์ Root Cause
			กำหนดปัญหาให้ชัดเจน (Define the Problem): ระบุสิ่งที่เกิดขึ้นจริง ผลกระทบ และขอบเขตของปัญหาอย่างชัดเจนและวัดผลได้
			รวบรวมข้อมูล (Gather Data): หาหลักฐาน บันทึกข้อมูล หรือสัมภาษณ์ผู้เกี่ยวข้อง เพื่อทำความเข้าใจสถานการณ์
			วิเคราะห์สาเหตุ (Analyze the Root Cause): ใช้เทคนิค เช่น "5 Why" (ถาม "ทำไม" เพื่อเจาะลึก 5 ครั้ง) หรือ แผนภูมิก้างปลา (Fishbone Diagram) 
			เพื่อแยกแยะสาเหตุที่เป็นไปได้จนถึงสาเหตุรากฐาน
			หาวิธีแก้ไข (Implement Solutions): กำหนดแนวทางแก้ไขที่ตรงจุดเพื่อกำจัดต้นตอของปัญหา
			ติดตามผล (Monitor): ติดตามและประเมินผลลัพธ์เพื่อยืนยันว่าการแก้ไขได้ผลและปัญหาไม่กลับมาเกิดขึ้นซ้ำ 
		- ประโยชน์ของการใช้ RCA
			แก้ไขปัญหาได้อย่างยั่งยืน: ป้องกันไม่ให้ปัญหาเดิมเกิดขึ้นซ้ำอีก
			ลดต้นทุนและความสูญเสีย: ไม่ต้องเสียเวลาและค่าใช้จ่ายในการแก้ไขปัญหาเดิมๆ ซ้ำซาก
			ปรับปรุงกระบวนการ: ช่วยให้เห็นจุดอ่อนในกระบวนการทำงานและแก้ไขให้ดี
   -CHECKLIST Template เช่่น
        -กระบวนการทำงานแต่ละขั้นตอน
		-กระบวนการตรวจสอบ การทำงานแต่ละขั้นตอน
		
### โฟลเดอร์หลัก (Modules)
- REST API
- Golang   
- GORM 
- Entities

โปรเจกต์ใช้ **Clean Architecture** 3-layer + Delivery:

1. MQTT Request-Response Pattern** – รองรับการส่งคำขอและรอรับ response ตาม pattern ที่มีใน TypeScript (subscribe ชั่วคราว, timeout, unsubscribe อัตโนมัติ)
2. Cache Layer** – สำหรับลดการเรียก MQTT ซ้ำ (Redis cache)
3. InfluxDB Integration** – บันทึกข้อมูล MQTT และ alarm logs ลง InfluxDB (time-series)
4. HTTP API** – เพิ่ม endpoint สำหรับดึงข้อมูล MQTT แบบ request-response (คล้าย `/v1/mqtt2/topic`) และจัดการ cache
5. Alarm Processing** – ปรับใช้ logic การประเมิน alarm จาก `iot.helper.ts` มาเป็น Go (ใช้ struct และฟังก์ชัน)
6. Websockets  ,socket IO  Real-Time Real-time communication with WebSockets & Kafka 
7. Kafka  queue process


- **ข้อมูล (Data)** : ตัวเลข, ข้อความ, รายการต่างๆ
- **การประมวลผล (Processing)** : การดำเนินการกับข้อมูล เช่น การคำนวณ การเปรียบเทียบ
- **การควบคุมการทำงาน (Control Flow)** : การตัดสินใจ (if-else), การวนซ้ำ (loop)
- **การจัดเก็บ (Storage)** : หน่วยความจำ, ไฟล์, ฐานข้อมูล
- **อินพุต/เอาท์พุต (I/O)** : การรับข้อมูลจากผู้ใช้ หรือแสดงผล

### โฟลเดอร์หลัก (Modules)
โปรเจกต์ใช้ **Clean Architecture** 3-layer + Delivery:

| Layer | ตำแหน่ง | หน้าที่ |
|-------|---------|--------|
| **Model** | `internal/models/` | Entity (GORM) – `User`, `Session`, `VerificationToken` |
| **Repository** | `internal/repository/` | อ่าน/เขียน DB และ Redis ผ่าน interface |
| **Usecase** | `internal/usecase/` | Business logic: hash, JWT, email queue, validation |
| **Delivery** | `internal/delivery/rest/` | HTTP handlers, middleware, DTO, router |
| **Worker** | `internal/delivery/worker/` | Background job สำหรับส่งอีเมล |
- -------------
- What you'll learn
- Comprehensive Examples of Basic Concepts.
- Detailed Explanation and Practice of Intermediate level Concepts in Go.
- Highly Extensive Section on Advanced Concepts in Golang.
- Detailed Explanation of GoRoutines: Complete Coverage with many examples to master the concept.
- Comprehensive Explanation and Extensive Practice on Protocol Buffers and gRPC.
- We will make REST API in Go.
- We will make a gRPC API in Go.
- How concurrency works in Go?
- Quizzes and Slides with downloadable PDF material.
- Git and Github.
- Pointers in Go.
- Detailed Explanation and Practical Examples of Struct, Maps, Slices in Go.
- Importance and Various Use Cases of CHANNELS in Go.
- Real Use Case Based API Examples with SQL and NoSQL Usage.
- API Folder Structure.
- Learn How to Plan Before Making an API. *** Important for beginners***
- Learn How to Make Professional, Industry Standard APIs
- MongoDB and MariaDB(Drop in replacement for MySQL)
- Advanced API Benchmarking Tools like wrk, h2load, ghz etc.
- Make HTTP2, HTTPS API.
- Learn How to Implement TLS/SSL in API.
- Learn How to Code Your Own Middleware from scratch
- Learn how to read Go Source Code and Find Solutions to Any Problem
- Learn to use Algorithms in Real World Cases
- Interview Preparation: Question Bank with 350+ Questions and Answers
- How Does Go runtime Work? Why is it important to understand it?
- Become an expert in using Reflect Package. Comprehensive use of Reflect in gRPC & REST API projects in this course.

-  *************************

# Go Programming Course - Complete Curriculum
  
## Section 1: Introduction and Setup
- Greetings and Welcome!
- Some tips while using this course
- Course Content
- About Go Language
- Why choose Go?
- Go Playground
- Installing Go on Linux
- Installing Go on Windows
- Installing Go on Mac
- IDE/Code Editor
- Installing VS Code on Linux
- Installing VS Code on Windows
- Installing VS Code on Mac
- Setting Up Development environment: Extensions
- Resources

## Section 2: Version Control with Git and GitHub
- What is Git? What is VCS?
- Installing Git on Linux
- Installing Git on Windows
- Installing Git on Mac
- Github
- Github and Git : SSH
- git init
- git add
- git commit
- git remote
- git push
- Course Setup
- Important Note

## Section 3: Go Fundamentals
- Hello World!
- Go Run
- Go Compiler
- The Standard Library
- Import statement
- Data Types
- Variables
- Naming Conventions
- Constants
- Arithmetic Operations
- Loop: For (break continue)
- Loop: For (using as while)
- Operators
- Conditions: If else
- Conditions: Switch
- Arrays and Blank Identifier
- Slices
- Maps
- Range
- Functions
- Multiple Return Values
- Variadic functions
- Defer
- Panic
  - Learn how the panic function halts execution, unwinds the stack, and runs deferred cleanup, signaling an unrecoverable error in Go, with any type via interface.
- Recover
- Exit
- Init function
- Basics Quiz
- Section Summary and Motivation

## Section 4: Intermediate Go Concepts
- Closures
- Recursion
- Pointers
- Strings and Runes
- Formatting Verbs
- Fmt package
- Structs
- Methods
- Interfaces
- Struct Embedding
- Generics
- Intermediate Quiz 1

### Error Handling and String Operations
- Errors
- Custom Errors
- String Functions
- String Formatting
- Text Templates
- Regular Expressions
- Time
- Epoch
- Time Formatting / Parsing
- Random Numbers
- Number Parsing
- Intermediate Quiz 2

### File I/O and System Operations
- URL Parsing
- bufio package
- Base64 Coding
- SHA 256/512 Hashes / Hashing / Cryptography / Crypto Package
- Writing Files
- Reading Files
- Line Filters
- File Paths
- Directories
- Temporary Files and Directories
- Embed Directive
- Intermediate Quiz 3

### Command Line and Configuration
- Command Line Arguments/Flags
- Command Line Sub Commands
- Environment Variables
- Logging
- JSON
- Struct Tags
- XML
- Go Extension
- Type Conversions
- IO package
- Math package
- Math package (Code Examples)
- Intermediate Quiz 4
- Section Summary and Motivation

## Section 5: Concurrency in Go
- Goroutines
- Channels - Introduction
- Unbuffered Channels and Runtime Mechanism
- Buffered Channels
- Channel Synchronization
- Advanced Quiz 1

### Advanced Channel Operations
- Channel Directions
- Multiplexing using Select
- Non blocking channel operations
- Closing Channels
- Advanced Quiz 2

### Concurrency Patterns
- Context
- Timers
- Tickers
- Worker Pools
- Wait Groups
- Advanced Quiz 3

### Synchronization and Advanced Topics
- Mutexes
- Atomic Counters
- Rate Limiting
- Rate Limiting - Token Bucket Algorithm
- Rate Limiting - Fixed Window Counter
- Rate Limiting - Leaky Bucket Algorithm
- Stateful Goroutines
- Sorting
- Advanced Quiz 4

### Testing and Reflection
- Testing / Benchmarking
- Executing Processes / OS Processes / Other Processes
- Signals
- Reflect
- Advanced Quiz 5
- Section Summary and Congratulations

### Advanced Concurrency Patterns
- Concurrency vs Parallelism
- Race Conditions
- Deadlocks
- RWMutex
- sync.NewCond
- sync.Once
- sync.Pool
- for select statement
- Advanced Concurrency Quiz

## Section 6: Web Development and API Building
### Internet Fundamentals
- URL/URI
- Request Response Cycle
- What is Frontend Dev/ Client Side
- What is Backend Dev/ API / Server Side
- HTTP 1/2/3, HTTPS
- Internet Quiz
- OS Choice for Development
- What is REST API
- Endpoints
- HTTP Client
- HTTP Server
- Ports

### API Testing and Tools
- Postman for API Testing
- Install wrk (Benchmarking Tool)
- Install Htop
- Benchmarking an API
- Modules - go mod init

### Building the API
- Let's begin making the API/Server
- Downloading Third Party/External Packages - go get <package link>
- Let's add HTTP2 and HTTPS to our API
- https certificates - SSL/TLS
- Postman for TLS + HTTP2 Requests
- Using Curl to make http2 request
- HTTP2/HTTPS/HTTP Connections, TLS Handshake
- mTLS and Postman Settings
- Benchmarking HTTP1 vs HTTP2 -H2Load BM Tool
- Serialization/Deserialization - Marshal/Unmarshal - Encode/Decode

### API Structure and Routing
- API Folder Structure
- API Planning Stage
- Basic Routing-CRUD-HTTP Methods
- Processing Requests
- Path Params
- Query Params
- .gitignore file
- Multiplexer (mux)
- Middlewares
- Middlewares - Security Headers
- Middlewares - CORS
- Middlewares - Response Time
- Middlewares - Compression
- Middlewares - Rate Limiter
- Middlewares - HPP
- Middlewares - Ordering
- Efficient Middleware Chaining
- Older Routing Technique (Pre Go 1.22)

### CRUD Operations - GET and POST
- Getting All/Filtered/One Entry(ies) - GET
- Adding Single Entry/Multiple Entries - POST
- Handlers Refactoring

### Database Integration - MariaDB/MySQL
- MariaDB/MySQL - Introduction
- MariaDB Installation
- MariaDB GUI Tool - DBeaver Installation
- SQL Primer - CRUD - Command Line
- SQL Primer - CRUD - DBeaver
- Connect API to SQL
- Environment Variables (.env file)
- Creating our SQL Database
- Updating POST methods to post in Database
- Updating GET method to Fetch One Entry from Database
- Updating GET method to Fetch Multiple Entries from Database
- WHERE 1=1. WHY???
- Advanced Filtering Technique - GET - Getting Entries Based on Multiple Criteria
- Advanced Sort Order Technique - GET - Get Entries Based on Multiple Criteria

### CRUD Operations - PUT, PATCH, DELETE
- Updating a 'Complete Entry' - PUT
- Modifying An Entry - PATCH
- Improving our PATCH function - Reflect Package
- Deleting An Entry - DELETE
- Modernizing Routes - Older Routing Technique and its Limitations
- Refactoring Mux
- Using Path Params for Specific Entry
- Modifying Multiple Entries - PATCH
- Deleting Multiple Entries - DELETE

### Data Modeling and Validation
- Modelling Data
- Refactoring Database Operations
- Error Handling
- Struct Tags
- Data Validation

### Student and Teacher Routes
- Students Database Creation
- CRUD for Students Route
- Students Routes and Testing
- New Subroutes
- Getting Student List for a Specific Teacher
- Getting Student Count for a Specific Teacher
- Router Refactoring

### Execs Routes
- Execs Router
- Execs Model and Database Table
- CRUD for Execs Route

### Authentication and Authorization
- Passwords - Hashing
- Authorization and Authentication
- Cookies, Sessions and JWT
- Login Route Part 1 - Data Validation
- Login Route Part 2 - Password Hashing - Argon2
- Login Route Part 3 - JWT, Cookie
- Login Route Refactoring
- Logout
- Authentication Middleware - JWT
- Skipping Routes With Middlewares - PreLogin
- Update Password
- Sending Emails - MailHog
- Forgot Password
- Reset Password
- CSRF

### Advanced API Features
- Adding Pagination
- Data Sanitization - XSS Middleware
- Authorization
- Middleware Sequence Revisited
- Code Obfuscation
- Adjustments Before Final Binary
- API Binary
- Extensive Benchmarking - Source Code v/s Go Binary v/s Obfuscated
- Section Summary and Motivation

## Section 7: Protocol Buffers and gRPC
### Protocol Buffers Fundamentals
- What are Protocol Buffers?
- Syntax and Structure of .proto Files
- Packages in Protocol Buffers
- Messages in Protocol Buffers
- Fields in Protocol Buffers
- Field Types and Data Types
- Field Numbers
- Serialization and Deserialization
- RPC (Remote Procedure Call) in Protocol Buffers
- Versioning and Backward Compatibility
- Best Practices for .proto Files
- Installing Protoc Compiler to Generate Code from .proto Files
- Protocol Buffers in Practice
- Protocol Buffers Quiz

### gRPC Basics
- What is gRPC?
- Stubs
- What is Service?
- REST vs gRPC
- Creating Simple gRPC Server
- Creating a Simple gRPC Client
- gRPC + TLS
- Deep Dive - Proto Buf Packages + RPC

### gRPC Streaming
- gRPC Streaming
- Server Side Stream
- Client Side Stream
- BiDirectional Stream
- Advanced gRPC Features

### gRPC Tools and Testing
- Metadata, Headers and Trailers
- Postman for gRPC
- gRPCurl for gRPC
- Protoc Gen Validate Plugin
- Combo API (gRPC + REST functionality in One API)
- Benchmarking Combo API - GHZ BM Tool
- gRPC Quiz

## Section 8: MongoDB Integration
### MongoDB Fundamentals
- Intro
- MongoDB and NoSQL - Introduction
  - Compare SQL and NoSQL databases, and show how MongoDB stores data in BSON documents with a flexible schema for scalable, high-performance data management.
- MongoDB - Installation
  - Install MongoDB community edition on Debian or other platforms, following step by step instructions to start mongod, enable it, and verify with the mongo-shell.
- MongoDB Compass - GUI for MongoDB
- MongoDB Primer - CRUD

### Building gRPC API with MongoDB
- gRPC API Folder Structure and Project Requirements
- Creating Proto files based on our REST API routes
- Creating Project's gRPC Server
- Downloading Known Required Dependencies
- Connect API to MongoDB
- Error Handling
- Adding New Teacher(s)
- Refactoring
- Getting Teacher(s) - Filter
- Getting Teacher(s) - Sorting
- Getting Teacher(s) - Finalizing
- Interfaces - Common Filter for all Get RPCs
- Decode Function
- Generics - Common Decode for all Get Functions
- Modifying Teacher(s)
- Generics - Mapping Helpers - Refactored
- Deleting Teacher(s)

### Student and Exec Operations with MongoDB
- Adding New Student(s) and Exec(s)
- Getting Student(s) and Exec(s)
- Modifying Student(s) and Exec(s)
- Deleting Student(s) and Exec(s)

### Advanced MongoDB Operations
- Relationships in NoSQL (MongoDB)
- Getting Students By Teacher - RPC
- Getting Student Count By Teacher - RPC

### Authentication and Authorization with MongoDB
- Login RPC
- Update Password RPC
- Deactivate User RPC
- Forgot Password RPC
- Reset Password RPC

### gRPC Middleware and Security
- Response Time Interceptor
- Rate Limiting Interceptor
- Authentication Interceptor
- Logout RPC
- Authorization
- gRPC Advantage - No need for HPP, Sanitize, Compression, HTTP Headers, CORS
- Interceptor Sequence
- Data Validation using Protoc Gen Validate
- TLS/SSL + gRPC

### Final Deployment
- Code Obfuscation and API Binary
- Benchmarking

---

## Course Summary
- **Total Duration**: 96 hours 31 minutes
- **Total Lectures**: 328
- **Sections**: 8 comprehensive sections covering everything from Go basics to advanced microservices with gRPC and MongoDB
- **Projects**: REST API with MySQL/MariaDB and gRPC API with MongoDB

---

## สารบัญ

1. บทนำ
2. บทนิยามศัพท์
3. การออกแบบ Workflow
4. บทสรุป
5. คู่มือการปฏิบัติงาน
   - TASK LIST Template
   - CHECKLIST Template
   - การวิเคราะห์สาเหตุที่แท้จริง (RCA)
6. ภาคผนวก: เทมเพลตและตัวอย่างโค้ด

---

## 1. บทนำ

เอกสารนี้จัดทำขึ้นเพื่อเป็นคู่มือประกอบการพัฒนาแอปพลิเคชัน Backend ด้วยภาษา **Go** โดยใช้สถาปัตยกรรมแบบ **Clean Architecture (3‑Layer + Delivery)** ซึ่งประกอบด้วยชั้น Model, Repository, Usecase และ Delivery ภายในโปรเจกต์มีการบูรณาการระบบย่อยหลายตัว ได้แก่:

- **REST API** – ให้บริการ HTTP endpoints สำหรับการดำเนินการ CRUD, การยืนยันตัวตนด้วย JWT, การจัดการรหัสผ่าน และการแจ้งเตือนทางอีเมล
- **WebSocket** – รองรับการสื่อสารแบบ Real‑Time ระหว่าง client และ server ผ่าน WebSocket protocol
- **MQTT Request‑Response Pattern** – รองรับการส่งคำขอและรอรับ response ผ่าน MQTT พร้อมกลไกการ subscribe ชั่วคราวและ timeout
- **Redis Cache** – ใช้สำหรับเก็บ session JWT, cache ข้อมูลผู้ใช้ และลดการเรียก MQTT ซ้ำ
- **InfluxDB** – ใช้สำหรับบันทึกข้อมูลเชิงเวลา (time‑series) เช่น ข้อมูลจากเซ็นเซอร์ MQTT และประวัติแจ้งเตือน (alarm logs)
- **Alarm Processing** – ตรรกะการประเมินค่าเพื่อสร้างการแจ้งเตือนจากข้อมูลที่ได้รับ
- **Kafka** – ใช้เป็นระบบคิว (Queue Processor) สำหรับประมวลผลงานแบบอะซิงโครนัส เช่น การส่งอีเมล หรือการประมวลผลข้อมูลจำนวนมาก
- ** gRPC **  (gRPC Remote Procedure Call) คือ Open-Source Framework ประสิทธิภาพสูงจาก Google ที่ใช้สำหรับให้แอปพลิเคชันหรือเซอร์วิสต่างๆ สื่อสารกันผ่านเครือข่าย นิยมนำมาใช้ในระบบ Microservices เพราะรับ-ส่งข้อมูลได้รวดเร็ว ขนาดเล็ก และสร้างโค้ดอัตโนมัติจุดเด่นสำคัญของ gRPC มีดังนี้:ใช้ Protocol Buffers (Protobuf): แทนที่จะส่งข้อมูลแบบ JSON ทั่วไป gRPC จะแปลงข้อมูลให้เป็นรูปแบบ Binary ทำให้ข้อมูลมีขนาดเล็ก ส่งได้เร็วขึ้น และใช้ทรัพยากรน้อยกว่าใช้โปรโตคอล HTTP/2: รองรับการเชื่อมต่อแบบ Multiplexing และการสื่อสารแบบ Streaming (ส่งข้อมูลต่อเนื่องทั้งทางเดียวและสองทางพร้อมกัน)สร้างโค้ดได้อัตโนมัติ (Code Generation): สามารถเขียนไฟล์ .proto เพื่อกำหนดโครงสร้างข้อมูลและฟังก์ชันที่ต้องการเรียกใช้งาน แล้วคอมไพเลอร์จะสร้างโค้ดสำหรับ Client และ Server   Go, 

เอกสารนี้มีเป้าหมายเพื่อให้ผู้อ่านเข้าใจภาพรวมของระบบ รู้จักส่วนประกอบต่าง ๆ และสามารถนำไปปรับใช้ในโครงการจริงได้อย่างมีประสิทธิภาพ พร้อมทั้งมีเทมเพลตและ checklist สำหรับการวางแผนและตรวจสอบการทำงาน

---

## 2. บทนิยามศัพท์

### 2.1 Clean Architecture (สถาปัตยกรรมแบบสะอาด)

Clean Architecture เป็นรูปแบบการจัดโครงสร้างซอฟต์แวร์ที่แยกความรับผิดชอบออกเป็นชั้นต่าง ๆ เพื่อให้โค้ดมีความเป็นระเบียบ ทดสอบได้ง่าย และบำรุงรักษาได้ระยะยาว ในโปรเจกต์นี้แบ่งเป็น 4 ชั้นหลัก:

- **Model** – กำหนดโครงสร้างข้อมูล (entity) ที่ใช้ในแอปพลิเคชัน เช่น `User`, `Session`, `VerificationToken` โดยใช้ GORM เป็น ORM
- **Repository** – ทำหน้าที่ติดต่อกับแหล่งข้อมูล (Database, Redis, InfluxDB, MQTT) ผ่าน interface ทำให้ชั้นที่สูงกว่าไม่ต้องรู้รายละเอียดการจัดเก็บ
- **Usecase** – บรรจุตรรกะทางธุรกิจ (business logic) เช่น การแฮชรหัสผ่าน, การสร้าง JWT, การตรวจสอบข้อมูล, การจัดการคิวอีเมล
- **Delivery** – ชั้นที่ติดต่อกับโลกภายนอก ประกอบด้วย HTTP handlers (REST API), WebSocket handlers, และ background workers

### 2.2 MQTT Request‑Response Pattern

MQTT (Message Queuing Telemetry Transport) เป็นโปรโตคอลแบบ publish/subscribe ที่มีน้ำหนักเบา เหมาะกับ IoT และระบบที่มีข้อจำกัดด้านแบนด์วิดท์

**Request‑Response Pattern** คือรูปแบบที่ client ส่งคำขอ (publish) ไปยัง topic หนึ่ง และรอ response จาก server ผ่าน topic อีกตัวหนึ่ง โดยจะมีการ subscribe ชั่วคราวเพื่อรอ response และ unsubscribe อัตโนมัติเมื่อได้รับหรือ timeout

**ประโยชน์ที่ได้รับ**:
- รองรับการสื่อสารแบบสองทาง (bidirectional) บน MQTT
- ช่วยให้ระบบที่ใช้ MQTT สามารถทำงานแบบ request‑response ได้คล้าย HTTP
- ลดภาระการเก็บสถานะ (stateless) บนเซิร์ฟเวอร์

### 2.3 WebSocket

WebSocket เป็นโปรโตคอลที่ให้การสื่อสารแบบ full‑duplex ผ่านการเชื่อมต่อ TCP เดียว ช่วยให้ server สามารถส่งข้อมูลไปยัง client ได้ทันทีโดยไม่ต้องรอคำขอ (real‑time)

**Hub** – เป็นโครงสร้างที่รวบรวมและจัดการ client connections ทั้งหมด รวมถึงการ broadcast ข้อความไปยัง client ต่างๆ  
**Client** – แทนการเชื่อมต่อ WebSocket แต่ละตัว มี channel สำหรับรับและส่งข้อความ

### 2.4 Kafka Queue Processor

Kafka เป็นแพลตฟอร์มกระจายสำหรับการสตรีมข้อมูล (distributed streaming platform) ทำหน้าที่เป็นระบบคิวข้อความ (message queue) ที่มีความทนทานสูงและรองรับปริมาณงานมหาศาล

**Queue Processor** คือกระบวนการที่อ่านข้อความจาก Kafka topic และประมวลผลตามลำดับ เช่น การส่งอีเมล, การบันทึก log, การวิเคราะห์ข้อมูล

**ประโยชน์**:
- แยกส่วนการผลิตและการบริโภคข้อมูล ทำให้ระบบมีความยืดหยุ่น
- รองรับการทำงานแบบอะซิงโครนัส ลดการรอคอยใน request cycle
- สามารถขยายจำนวน consumer ได้ตามโหลด

### 2.5 Redis Cache

Redis เป็นฐานข้อมูลแบบ in‑memory ที่ใช้สำหรับแคชข้อมูลและเก็บ session

ในโปรเจกต์นี้ Redis ถูกใช้เพื่อ:
- เก็บ Refresh Token ของผู้ใช้
- แคชข้อมูลผู้ใช้เพื่อลดการ query ฐานข้อมูล
- แคชผลลัพธ์จากการเรียก MQTT เพื่อลดการเรียกซ้ำในระยะเวลาสั้น

### 2.6 InfluxDB

InfluxDB เป็นฐานข้อมูลเชิงเวลา (time‑series database) ที่ออกแบบมาเพื่อจัดเก็บและสอบถามข้อมูลที่มีป้ายกำกับเวลา เช่น ค่าจากเซ็นเซอร์, logs

ในระบบนี้ InfluxDB ใช้บันทึก:
- ข้อมูล MQTT ที่ได้รับจากอุปกรณ์
- ประวัติการแจ้งเตือน (alarm logs) ที่เกิดจากตรรกะการประเมินค่า

### 2.7 Alarm Processing

Alarm Processing เป็นกระบวนการประเมินข้อมูลที่ได้รับ (เช่น ค่าอุณหภูมิ, ความชื้น) เทียบกับเกณฑ์ที่กำหนด เพื่อตัดสินใจว่าควรแจ้งเตือนหรือไม่

**ตรรกะ** มักประกอบด้วย:
- การเปรียบเทียบค่ากับขีดจำกัดสูง/ต่ำ
- การตรวจจับการเปลี่ยนแปลงอย่างรวดเร็ว (rate of change)
- การหน่วงเวลา (debounce) เพื่อป้องกันการแจ้งเตือนซ้ำซ้อน

### 2.8 gRPC และ Protocol Buffers

**Protocol Buffers (Protobuf)** เป็นภาษาสำหรับกำหนดโครงสร้างข้อมูล (interface definition language) ที่ใช้ในการส่งข้อมูลแบบไบนารี มีประสิทธิภาพสูงกว่า JSON/XML

**gRPC** เป็นเฟรมเวิร์ก RPC ที่ใช้ Protobuf เป็นข้อมูลพื้นฐาน รองรับการสตรีมแบบต่างๆ และทำงานบน HTTP/2 ให้ประสิทธิภาพสูง

---

## 3. การออกแบบ Workflow

ในส่วนนี้จะอธิบายกระแสข้อมูล (data flow) ของระบบหลัก ๆ พร้อมตัวอย่างการทำงานจริง

### 3.1 Data Flow Diagram (Draw.io Style)

```
┌─────────────────────────────────────────────────────────────────────┐
│                          ผู้ใช้ / Client                            │
└───────────────┬─────────────────┬─────────────────┬───────────────┘
                │                 │                 │
                │ HTTP/REST       │ WebSocket       │ MQTT
                │ (chi)           │ (gorilla/ws)    │ (paho)
                ▼                 ▼                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        Delivery Layer                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐             │
│  │ REST Handler │  │ WS Handler   │  │ MQTT Handler │             │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘             │
└─────────┼─────────────────┼─────────────────┼─────────────────────┘
          │                 │                 │
          ▼                 ▼                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        Usecase Layer                               │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  Business Logic: JWT, Password Hash, Validation, Alarm      │   │
│  └─────────────────────────────────────────────────────────────┘   │
└─────────┬───────────────────────┬───────────────────────┬─────────┘
          │                       │                       │
          ▼                       ▼                       ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Repository Layer                              │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌──────────┐  ┌────────┐ │
│  │  GORM   │  │  Redis  │  │InfluxDB │  │  MQTT   │  │ Kafka  │ │
│  │ (MySQL) │  │ (Cache) │  │(Time-   │  │(Request-│  │(Queue) │ │
│  └─────────┘  └─────────┘  │ Series) │  │Response)│  └────────┘ │
│                             └─────────┘  └──────────┘            │
└─────────────────────────────────────────────────────────────────────┘
```

### 3.2 คำอธิบายขั้นตอนการทำงานหลัก

#### 3.2.1 MQTT Request‑Response Workflow

1. **Client** ส่งคำขอ (request) โดย publish ไปยัง topic เช่น `request/device/123` พร้อม payload ที่มี `correlationId` และ `replyTopic`
2. **MQTT Handler** ใน Delivery Layer รับข้อความและส่งต่อไปยัง Usecase
3. **Usecase** ประมวลผลคำขอ เช่น อ่านค่าจากเซ็นเซอร์ หรือเขียนข้อมูลลง InfluxDB
4. **Usecase** สร้าง response และ publish ไปยัง `replyTopic` ที่ client ระบุ พร้อม `correlationId` เดิม
5. **Client** ซึ่ง subscribe ไว้ที่ `replyTopic` จะได้รับ response และจับคู่กับคำขอผ่าน `correlationId`

**กลไกการจัดการ timeout**: Client จะตั้ง timer ไว้ หากครบเวลาที่กำหนดยังไม่ได้รับ response จะถือว่าคำขอหมดอายุและดำเนินการตามที่กำหนด (เช่น retry หรือแจ้ง error)

#### 3.2.2 Alarm Processing Workflow

1. **MQTT Client** รับข้อมูลจากเซ็นเซอร์ (เช่น อุณหภูมิ) ผ่าน MQTT subscription
2. ข้อมูลถูกส่งไปยัง **Alarm Usecase** ซึ่งจะประเมินค่าตามเกณฑ์ที่กำหนด (เช่น ถ้าอุณหภูมิ > 40°C เป็นเวลา 5 วินาที ให้แจ้งเตือน)
3. หากเข้าเงื่อนไข **Alarm Usecase** จะ:
   - บันทึกเหตุการณ์ลง InfluxDB (alarm log)
   - ส่งข้อความแจ้งเตือนผ่าน WebSocket ไปยัง client ที่เกี่ยวข้อง (real‑time)
   - อาจส่งอีเมลหรือ LINE notify ผ่าน Kafka queue
4. **WebSocket Hub** จะ broadcast ข้อความแจ้งเตือนไปยังทุก client ที่ subscribe ไว้

#### 3.2.3 WebSocket Real‑Time Communication

1. **Client** ทำ HTTP request เพื่อ upgrade เป็น WebSocket ที่ endpoint เช่น `/ws`
2. **WebSocket Handler** รับ connection และสร้าง `Client` object พร้อมลงทะเบียนใน `Hub`
3. **Hub** จะมี goroutine หลัก 2 ตัว:
   - `run()`: รับคำสั่ง register/unregister และ broadcast ข้อความไปยัง clients
   - `broadcast()`: ส่งข้อความที่ได้รับจาก Kafka หรือ internal event ไปยัง client ทั้งหมดหรือเฉพาะกลุ่ม
4. Client สามารถส่งข้อความไปยัง server ได้ ซึ่ง server จะประมวลผลและอาจตอบกลับ หรือส่งต่อไปยังระบบอื่น (เช่น บันทึก InfluxDB)

#### 3.2.4 Kafka Queue Processing

1. **Producer** (เช่น REST API เมื่อมีการส่งอีเมล) ส่งข้อความไปยัง Kafka topic `email-queue`
2. **Worker** (background process) จะ consume ข้อความจาก topic ดังกล่าว
3. **Worker** ประมวลผล (เช่น สร้างเนื้อหาอีเมลด้วย Hermes และส่งผ่าน Gomail)
4. เมื่อส่งสำเร็จ Worker อาจบันทึก log หรืออัปเดตสถานะในฐานข้อมูล

### 3.3 ตัวอย่างการใช้งานจริง (Use Case)

**กรณีศึกษา: ระบบตรวจสอบอุณหภูมิในคลังสินค้า**

- เซ็นเซอร์ IoT ส่งค่าอุณหภูมิทุก 5 วินาทีผ่าน MQTT ไปยัง topic `sensors/temperature`
- ระบบจะรับข้อมูลและประเมินแจ้งเตือน หากอุณหภูมิสูงกว่า 35°C ต่อเนื่อง 3 ครั้ง จะเกิด alarm
- Alarm จะถูกส่งไปยัง WebSocket client ที่เป็นหน้า dashboard และบันทึกใน InfluxDB เพื่อนำไปวิเคราะห์แนวโน้ม
- ผู้ดูแลสามารถเรียกดูข้อมูลย้อนหลังผ่าน REST API หรือ gRPC API
- หากต้องการแจ้งเตือนทางอีเมลเมื่อเกิด alarm รุนแรง ระบบจะใส่ message ลง Kafka queue และ worker จะทำการส่งอีเมลให้ผู้รับผิดชอบ

### 3.4 เทมเพลตและตัวอย่างโค้ด

#### 3.4.1 การตั้งค่า MQTT Client (พร้อม Cache)

```go
// pkg/mqtt/client.go
package mqtt

import (
    "context"
    "time"
    mqtt "github.com/eclipse/paho.mqtt.golang"
    "github.com/redis/go-redis/v9"
)

type Client struct {
    conn   mqtt.Client
    redis  *redis.Client
    cacheTTL time.Duration
}

// GetDataFromTopic ดึงข้อมูลจาก MQTT โดยใช้ Redis cache เพื่อลดการเรียกซ้ำ
// หากมีใน cache จะคืนค่าทันที มิฉะนั้นจะ subscribe แบบ request-response
func (c *Client) GetDataFromTopic(ctx context.Context, topic string) ([]byte, error) {
    // 1. ตรวจสอบ cache
    cached, err := c.redis.Get(ctx, "mqtt:"+topic).Bytes()
    if err == nil {
        return cached, nil // cache hit
    }

    // 2. สร้าง correlationId และ replyTopic
    corrID := uuid.New().String()
    replyTopic := "reply/" + corrID

    // 3. Subscribe ชั่วคราวเพื่อรอ response
    responseChan := make(chan []byte, 1)
    token := c.conn.Subscribe(replyTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
        responseChan <- msg.Payload()
    })
    if !token.WaitTimeout(5 * time.Second) {
        return nil, errors.New("subscribe timeout")
    }

    // 4. Publish request พร้อม correlationId และ replyTopic
    payload := map[string]interface{}{
        "correlationId": corrID,
        "replyTopic":    replyTopic,
        "data":          "your request payload",
    }
    data, _ := json.Marshal(payload)
    token = c.conn.Publish(topic, 1, false, data)
    if !token.WaitTimeout(5 * time.Second) {
        return nil, errors.New("publish timeout")
    }

    // 5. รอ response (timeout 10 วินาที)
    select {
    case resp := <-responseChan:
        // บันทึก cache
        c.redis.Set(ctx, "mqtt:"+topic, resp, c.cacheTTL)
        return resp, nil
    case <-time.After(10 * time.Second):
        return nil, errors.New("response timeout")
    }
}
```

**คำอธิบาย**:
- ฟังก์ชันนี้จะพยายามดึงข้อมูลจาก Redis ก่อน ถ้าไม่มีจะทำ MQTT request-response
- มีการตั้ง timeout ทั้งตอน subscribe, publish และรอ response เพื่อป้องกันการค้าง
- เมื่อได้ response จะเก็บไว้ใน Redis เป็นเวลา `cacheTTL` เพื่อลดการเรียกซ้ำ

#### 3.4.2 WebSocket Hub และ Client

```go
// internal/pkg/websocket/hub.go
package websocket

import "sync"

// Hub จัดการ client connections ทั้งหมด
type Hub struct {
    clients    map[*Client]bool          // ลงทะเบียน client
    broadcast  chan []byte               // channel สำหรับ broadcast ข้อความ
    register   chan *Client              // รับ client ใหม่
    unregister chan *Client              // รับ client ที่断开
    mu         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

// Run เป็น goroutine หลักสำหรับจัดการ events
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()
        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mu.Unlock()
        case message := <-h.broadcast:
            h.mu.RLock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mu.RUnlock()
        }
    }
}
```

```go
// internal/pkg/websocket/client.go
package websocket

import (
    "github.com/gorilla/websocket"
)

type Client struct {
    hub  *Hub
    conn *websocket.Conn
    send chan []byte       // buffered channel สำหรับส่งข้อความ
}

// readPump อ่านข้อความจาก WebSocket และส่งไปยัง Hub
func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }
        // ประมวลผล message (เช่น ส่งไปยัง Kafka หรือ MQTT)
        // ...
        // ตัวอย่าง: broadcast กลับไปยังทุกคน
        c.hub.broadcast <- message
    }
}

// writePump ส่งข้อความจาก channel send ไปยัง WebSocket
func (c *Client) writePump() {
    defer c.conn.Close()
    for message := range c.send {
        err := c.conn.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            break
        }
    }
}
```

#### 3.4.3 Kafka Producer – การส่งอีเมลผ่าน Queue

```go
// internal/delivery/worker/email_worker.go
package worker

import (
    "context"
    "encoding/json"
    "github.com/IBM/sarama"
    "gomail"
)

type EmailTask struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

// EmailWorker consume messages from Kafka and send emails
func EmailWorker(brokers []string, topic string) {
    consumer, err := sarama.NewConsumer(brokers, nil)
    if err != nil {
        panic(err)
    }
    defer consumer.Close()

    partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
    if err != nil {
        panic(err)
    }
    defer partitionConsumer.Close()

    for msg := range partitionConsumer.Messages() {
        var task EmailTask
        if err := json.Unmarshal(msg.Value, &task); err != nil {
            // log error
            continue
        }
        // ส่งอีเมล
        m := gomail.NewMessage()
        m.SetHeader("From", "no-reply@example.com")
        m.SetHeader("To", task.To)
        m.SetHeader("Subject", task.Subject)
        m.SetBody("text/html", task.Body)

        d := gomail.NewDialer("smtp.example.com", 587, "user", "pass")
        if err := d.DialAndSend(m); err != nil {
            // log error, อาจ retry หรือส่งไป dead letter queue
        }
    }
}
```

#### 3.4.4 การบันทึกข้อมูลลง InfluxDB

```go
// internal/modules/influxdb/usecase/usecase.go
package usecase

import (
    "context"
    "time"
    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxUsecase struct {
    client influxdb2.Client
    org    string
    bucket string
}

func (u *InfluxUsecase) WriteSensorData(ctx context.Context, deviceID string, temperature float64, humidity float64) error {
    writeAPI := u.client.WriteAPI(u.org, u.bucket)
    // สร้าง point
    p := influxdb2.NewPoint(
        "sensor",
        map[string]string{"device": deviceID},
        map[string]interface{}{
            "temperature": temperature,
            "humidity":    humidity,
        },
        time.Now(),
    )
    writeAPI.WritePoint(p)
    writeAPI.Flush()
    return nil
}

func (u *InfluxUsecase) QueryAlarmLogs(ctx context.Context, deviceID string, start, stop time.Time) ([]AlarmLog, error) {
    queryAPI := u.client.QueryAPI(u.org)
    query := `from(bucket:"` + u.bucket + `")
              |> range(start: ` + start.Format(time.RFC3339) + `, stop: ` + stop.Format(time.RFC3339) + `)
              |> filter(fn: (r) => r._measurement == "alarm" and r.device == "` + deviceID + `")
              |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")`
    result, err := queryAPI.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    // แปลงผลลัพธ์เป็น struct
    var logs []AlarmLog
    for result.Next() {
        // ... mapping
    }
    return logs, nil
}
```

---

## 4. บทสรุป

### 4.1 ประโยชน์ที่ได้รับ

- **ความยืดหยุ่นสูง** – การแยกชั้นตาม Clean Architecture ทำให้สามารถเปลี่ยนเทคโนโลยีในแต่ละชั้นได้โดยไม่กระทบต่อชั้นอื่น (เช่น เปลี่ยน GORM เป็น Ent หรือเปลี่ยนฐานข้อมูล)
- **รองรับ Real‑Time** – WebSocket และ MQTT ช่วยให้ระบบสามารถตอบสนองแบบทันทีทันใด เหมาะกับ IoT และแอปพลิเคชันแบบ monitoring
- **ประสิทธิภาพและความทนทาน** – การใช้ Redis Cache ช่วยลดภาระระบบ, InfluxDB ช่วยจัดการข้อมูลเวลา, Kafka ช่วยแยกงานแบบอะซิงโครนัส ทำให้ระบบขยายตัวได้ดี
- **ความปลอดภัย** – มีระบบ JWT, การแฮชรหัสผ่าน, และ middleware ต่างๆ ช่วยป้องกันภัยคุกคามทั่วไป
- **ทดสอบง่าย** – การใช้ interface และ dependency injection ทำให้สามารถเขียน unit test ได้สะดวก

### 4.2 ข้อควรระวัง

- **การจัดการ Goroutine** – ต้องระวัง goroutine leak โดยเฉพาะใน WebSocket และ Kafka consumer ควรมีกลไก graceful shutdown
- **Timeout และ Retry** – ในการสื่อสารผ่าน MQTT และ Kafka ควรตั้งค่า timeout และ retry อย่างเหมาะสม เพื่อป้องกันการค้างของระบบ
- **การซิงโครไนซ์ข้อมูล** – การใช้ cache (Redis) อาจทำให้ข้อมูลไม่สดใหม่ ควรกำหนด TTL ให้เหมาะสม หรือมีกลไก invalidate cache เมื่อข้อมูลเปลี่ยนแปลง
- **ความซับซ้อน** – การมีหลาย component (MQTT, WebSocket, Kafka, InfluxDB) ทำให้ระบบมีความซับซ้อนในการดูแลและดีบัก ควรมีระบบ logging และ monitoring ที่ดี

### 4.3 ข้อดี

- รองรับการทำงานแบบเรียลไทม์และ asynchronous อย่างมีประสิทธิภาพ
- สถาปัตยกรรมที่ยืดหยุ่นและเป็นมาตรฐาน ทำให้ทีมพัฒนาเข้าใจและทำงานร่วมกันได้ง่าย
- ใช้เทคโนโลยีที่ได้รับความนิยมและมี community ขนาดใหญ่ (Go, gRPC, Redis, Kafka, InfluxDB)

### 4.4 ข้อเสีย

- การติดตั้งและกำหนดค่าเริ่มต้นค่อนข้างซับซ้อน ต้องจัดการหลาย service พร้อมกัน
- ต้องการความรู้พื้นฐานด้าน distributed systems และ concurrency เพื่อให้ใช้งานได้อย่างถูกต้อง
- การ debug ในระบบที่มีหลาย component อาจทำได้ยากกว่า monolithic แบบเดิม

### 4.5 ข้อห้าม (ถ้ามี)

- **ห้ามใช้ JWT แบบไม่ตั้งค่าหมดอายุ** – ควรตั้งค่า expiration ให้สั้น และใช้ refresh token หมุนเวียน
- **ห้ามเก็บข้อมูลสำคัญ (เช่น password) ใน cache โดยไม่เข้ารหัส** – Redis อาจถูกโจมตีได้ ควรเข้ารหัสข้อมูลที่ sensitive
- **ห้ามใช้ MQTT QoS 0 ในกรณีที่ต้องการความเชื่อถือได้** – ควรเลือก QoS 1 หรือ 2 สำหรับข้อมูลสำคัญ
- **ห้ามข้ามการ validate input** – ควรใช้ validator และ sanitize ทุกครั้ง เพื่อป้องกัน injection และ XSS

### 4.6 แหล่งอ้างอิง

1. เอกสารทางการของ Go: [https://go.dev/doc/](https://go.dev/doc/)
2. Clean Architecture โดย Robert C. Martin: [https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
3. MQTT Specification: [https://mqtt.org/mqtt-specification/](https://mqtt.org/mqtt-specification/)
4. gRPC Documentation: [https://grpc.io/docs/](https://grpc.io/docs/)
5. Apache Kafka Documentation: [https://kafka.apache.org/documentation/](https://kafka.apache.org/documentation/)
6. InfluxDB Documentation: [https://docs.influxdata.com/influxdb/](https://docs.influxdata.com/influxdb/)
7. Redis Documentation: [https://redis.io/documentation](https://redis.io/documentation)

---

## 5. คู่มือการปฏิบัติงาน

### 5.1 TASK LIST Template

ใช้สำหรับวางแผนและติดตามขั้นตอนการพัฒนาในแต่ละ sprint หรือ feature

| Task ID | ชื่องาน | รายละเอียด | ผู้รับผิดชอบ | สถานะ | หมายเหตุ |
|---------|--------|------------|-------------|--------|----------|
| T001 | ตั้งค่าโครงสร้างโปรเจกต์ | สร้างโฟลเดอร์ตาม Clean Architecture, ตั้งค่า go.mod | Dev A | Done | ใช้ Go 1.21+ |
| T002 | ติดตั้งและตั้งค่า Database | ติดตั้ง MySQL/PostgreSQL, สร้าง schema ตาม models | Dev B | Done | ใช้ GORM AutoMigrate |
| T003 | Implement JWT Authentication | สร้าง login/logout, middleware, Redis สำหรับ refresh token | Dev A | In Progress | ใช้ golang-jwt |
| T004 | สร้าง MQTT Client + Request-Response | implement pkg/mqtt พร้อม cache | Dev C | To Do | |
| T005 | Integrate InfluxDB | ตั้งค่า client และเขียน usecase สำหรับบันทึกข้อมูล | Dev C | To Do | |
| T006 | พัฒนา Alarm Logic | แปลงจาก iot.helper.ts มาเป็น Go struct และฟังก์ชัน | Dev A | To Do | |
| T007 | สร้าง WebSocket Hub และ Client | ใช้ gorilla/websocket, จัดการ connection | Dev B | To Do | |
| T008 | ตั้งค่า Kafka Producer/Consumer | กำหนด topic, implement worker สำหรับ email | Dev D | To Do | |
| T009 | สร้าง REST API endpoints | CRUD สำหรับ user, device, alarm config | Dev A | To Do | |
| T010 | ทดสอบและ Benchmark | ใช้ wrk, h2load, ghz, เขียน unit test | Dev D | To Do | |
| T011 | จัดทำเอกสาร | เขียน README, API docs, คู่มือนี้ | Dev D | In Progress | |

### 5.2 CHECKLIST Template

ใช้สำหรับตรวจสอบความสมบูรณ์ก่อนการ deploy หรือส่งมอบงาน

#### 5.2.1 การตั้งค่า Infrastructure

- [ ] ติดตั้งและตั้งค่า MySQL/PostgreSQL พร้อมกำหนด user, password, database
- [ ] ติดตั้ง Redis และกำหนดค่า maxmemory, policy
- [ ] ติดตั้ง InfluxDB และสร้าง bucket, organization
- [ ] ติดตั้ง Kafka และสร้าง topics ที่จำเป็น (email-queue, alarm-queue, log-queue)
- [ ] ตรวจสอบการเชื่อมต่อระหว่าง service ทั้งหมด (ping, telnet)

#### 5.2.2 การพัฒนา Code

- [ ] โค้ดผ่านการ go fmt, go vet, golint
- [ ] มี unit test ครอบคลุม usecase หลัก (อย่างน้อย 70%)
- [ ] มี integration test สำหรับ API endpoints
- [ ] การจัดการ error ครบถ้วน (ไม่ panic ใน production)
- [ ] มี logging ที่เพียงพอ (ใช้ zap)
- [ ] ตั้งค่า environment variables ทั้งหมดใน .env และมี .env.example

#### 5.2.3 ความปลอดภัย

- [ ] ใช้ HTTPS สำหรับ REST และ gRPC (TLS)
- [ ] JWT มี expiration และ refresh token mechanism
- [ ] Password ถูก hash ด้วย Argon2 หรือ bcrypt
- [ ] มี middleware สำหรับ CORS, Rate Limiting, XSS Prevention
- [ ] ตรวจสอบ input ทั้งหมดด้วย validator

#### 5.2.4 Performance

- [ ] มีการตั้งค่า connection pool สำหรับ database และ Redis
- [ ] ใช้ Redis cache สำหรับข้อมูลที่เรียกบ่อย
- [ ] มีการ benchmark และปรับแต่งให้เหมาะสม (ใช้ wrk, ghz)
- [ ] ตรวจสอบ goroutine leak (ใช้ pprof)

#### 5.2.5 Deployment

- [ ] สร้าง binary สำหรับ production (GOOS=linux GOARCH=amd64 go build)
- [ ] มี Dockerfile และ docker-compose สำหรับ development
- [ ] กำหนด health check endpoint
- [ ] มี graceful shutdown (handle SIGTERM)

### 5.3 การวิเคราะห์สาเหตุที่แท้จริง (RCA – Root Cause Analysis)

RCA เป็นกระบวนการค้นหาสาเหตุพื้นฐานของปัญหา เพื่อป้องกันไม่ให้เกิดซ้ำ ใช้เครื่องมือเช่น **5 Why** และ **Fishbone Diagram**

#### 5.3.1 แนวทางและขั้นตอนการวิเคราะห์

1. **กำหนดปัญหาให้ชัดเจน (Define the Problem)**  
   - ระบุสิ่งที่เกิดขึ้นจริง, ผลกระทบ, ขอบเขต เช่น “ระบบไม่สามารถส่งอีเมลยืนยันได้เป็นเวลา 10 นาที ส่งผลให้ผู้ใช้สมัครสมาชิกไม่สำเร็จ”

2. **รวบรวมข้อมูล (Gather Data)**  
   - ดู logs, metric, สอบถามผู้เกี่ยวข้อง, ตรวจสอบการตั้งค่า

3. **วิเคราะห์สาเหตุ (Analyze the Root Cause)**  
   - ใช้เทคนิค **5 Why** ถาม “ทำไม” ซ้ำ ๆ เพื่อเจาะลึก  
   - ใช้ **Fishbone Diagram** แยกสาเหตุเป็นด้านต่าง ๆ (คน, กระบวนการ, เทคโนโลยี, สภาพแวดล้อม)

4. **กำหนดแนวทางแก้ไข (Implement Solutions)**  
   - แก้ที่สาเหตุ ไม่ใช่แค่อาการ เช่น ถ้าสาเหตุคือ SMTP server timeout ให้เพิ่ม retry และปรับ timeout

5. **ติดตามผล (Monitor)**  
   - ตรวจสอบว่าปัญหาไม่กลับมา และประเมินประสิทธิภาพของแนวทางแก้ไข

#### 5.3.2 ตัวอย่างการวิเคราะห์ด้วย 5 Why

**ปัญหา**: ผู้ใช้ไม่ได้รับอีเมลยืนยันการลงทะเบียน

- **Why 1**: ทำไมอีเมลไม่ถูกส่ง?  
  → เพราะ worker ที่ส่งอีเมลไม่สามารถเชื่อมต่อกับ SMTP server ได้  
- **Why 2**: ทำไม worker เชื่อมต่อ SMTP ไม่ได้?  
  → เพราะ SMTP server response timeout หลังจาก 5 วินาที  
- **Why 3**: ทำไม SMTP server ถึง timeout?  
  → เพราะปริมาณอีเมลในช่วงนั้นสูงกว่า normal 5 เท่า (เกิด spike)  
- **Why 4**: ทำไมถึงมี spike?  
  → เพราะมีการส่งอีเมลแจ้งเตือนจำนวนมากจาก alarm system พร้อมกัน  
- **Why 5**: ทำไม alarm system ส่งพร้อมกัน?  
  → เพราะไม่มี rate limiting ในการส่งอีเมลแจ้งเตือน และไม่มี queue buffer

**Root Cause**: ขาด rate limiting และ buffer mechanism สำหรับการส่งอีเมลแจ้งเตือน

**แนวทางแก้ไข**:
- เพิ่ม rate limiter ใน worker (จำกัดจำนวนอีเมลต่อนาที)
- ใช้ Kafka เป็น buffer ก่อนส่งจริง
- ปรับปรุง retry logic ให้มีการ backoff

#### 5.3.3 ประโยชน์ของการใช้ RCA

- **แก้ไขปัญหาอย่างยั่งยืน** – ป้องกันไม่ให้เกิดซ้ำอีก
- **ลดต้นทุนและความสูญเสีย** – ไม่ต้องเสียเวลาแก้ปัญหาเดิมซ้ำ ๆ
- **ปรับปรุงกระบวนการ** – ช่วยให้เห็นจุดอ่อนในระบบและพัฒนาต่อไป

---

## ภาคผนวก: เทมเพลตและตัวอย่างโค้ดเพิ่มเติม

### A. การตั้งค่า Viper และ Cobra

```go
// cmd/root.go
package cmd

import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
    Use:   "myapp",
    Short: "MyApp is a Go backend with Clean Architecture",
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    cobra.OnInitialize(initConfig)
}

func initConfig() {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AutomaticEnv()
    if err := viper.ReadInConfig(); err != nil {
        // ignore
    }
}
```

### B. Middleware ตัวอย่าง (Logger + Recovery)

```go
// internal/delivery/rest/middleware/logger.go
package middleware

import (
    "net/http"
    "time"
    "go.uber.org/zap"
)

func Logger(logger *zap.Logger) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            next.ServeHTTP(w, r)
            logger.Info("request",
                zap.String("method", r.Method),
                zap.String("path", r.URL.Path),
                zap.Duration("latency", time.Since(start)),
                zap.String("ip", r.RemoteAddr),
            )
        })
    }
}
```

### C. การใช้ Air สำหรับ Hot‑Reload

ติดตั้ง:
```bash
go install github.com/cosmtrek/air@latest
```

สร้างไฟล์ `.air.toml`:
```toml
root = "."
tmp_dir = "tmp"
[build]
  cmd = "go build -o ./tmp/main ./cmd/apiser"
  bin = "./tmp/main"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  include_dir = []
  include_file = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_error = false
[color]
  main = "magenta"
  watcher = "cyan"
  build = "yellow"
  runner = "green"
[log]
  time = false
[proxy]
  enabled = false
  proxy_port = 0
  app_port = 0
```

จากนั้นรัน `air` เพื่อเริ่มพัฒนาด้วย hot‑reload

---

 
# ตัวอย่าง โครงสร้าง Folder

```bash
 icmongolang/
  api/
    ├── cmd/
    │   ├── apiser/                 # REST API หลัก (มีอยู่แล้ว)
    │   ├── websocket/              # *** WebSocket server
    │   │   └── main.go
    │   ├── initdata.go
    │   ├── root.go
    │   ├── serve.go
    │   └── worker.go
    ├── internal/
    │   ├── websocket/              # **** โค้ดเฉพาะของ WebSocket
    │   │   ├── delivery/
    │   │   │   └── ws/
    │   │   │       ├── hub.go
    │   │   │       ├── client.go
    │   │   │       └── handler.go     # HTTP endpoint สำหรับ upgrade
    │   │   ├── usecase/
    │   │   │   └── ws_usecase.go      # business logic (save message, auth)
    │   │   ├── repository/
    │   │   │   └── ws_repo.go         # interface สำหรับ DB
    │   │   └── models/
    │   │       └── ws_models.go       # entity ของ message, session
    │   ├── pkg/                    # shared packages (มีอยู่แล้ว)
    │   │   ├── websocket/          # **ปรับปรุง** ใช้ร่วมกันได้
    │   │   │   ├── hub.go          # core Hub logic
    │   │   │   ├── client.go
    │   │   │   └── message.go      # struct ของ message
    │   │   └── ...
    ├── migrations/                 # **เพิ่ม** SQL schema สำหรับ websocket
    │   └── 20250619_websocket_tables.sql
    └── 
```
```bash
    icmongolang/
    ├── pkg/
    │   ├── helpers/
    │   │   ├── iot.go          # Alarm logic (สมบูรณ์)
    │   │   └── format.go       # ฟังก์ชันช่วยเหลือ (time, string, random)
    │   ├── mqtt/
    │   │   └── client.go       # MQTT client พร้อม GetDataFromTopic
    │   ├── influxdb/
    │   │   └── client.go       # InfluxDB client
    │   └── redis/
    │       └── redis_conn.go   # Redis client + Cache interface
    ├── internal/
    │   ├── mqtt/
    │   │   ├── delivery/http/
    │   │   │   ├── handler.go
    │   │   │   └── routes.go
    │   │   ├── presenter/
    │   │   │   └── presenter.go
    │   │   └── usecase/
    │   │       └── usecase.go
    │   ├── influxdb/
    │   │   ├── delivery/http/
    │   │   │   ├── handler.go
    │   │   │   └── routes.go
    │   │   ├── presenter/
    │   │   │   └── presenter.go
    │   │   └── usecase/
    │   │       └── usecase.go
    │   ├── alarm/
    │   │   ├── delivery/http/
    │   │   │   ├── handler.go
    │   │   │   └── routes.go
    │   │   ├── repository/
    │   │   │   └── alarm_log_repo.go
    │   │   └── usecase/
    │   │       └── usecase.go
    │   └── server/
    │       ├── handlers.go
    │       └── server.go
    └── cmd/api/main.go (สมมติตามเดิม)
```
*****************************************************************************
# เอกสารประกอบการอบรม Go Bootcamp  
## สถาปัตยกรรมระบบแบบ Clean Architecture และการประมวลผลแบบ Real‑Time

---

## สารบัญ

1. แผนการสอน
2. บทนำ
3. บทนิยามศัพท์
4. การออกแบบ Workflow
5. บทสรุป
6. คู่มือการปฏิบัติงาน
   - TASK LIST Template
   - CHECKLIST Template
   - การวิเคราะห์สาเหตุที่แท้จริง (RCA)
7. ภาคผนวก: เทมเพลตและตัวอย่างโค้ด

---

## 1. แผนการสอน

### 1.1 วัตถุประสงค์
- เพื่อให้ผู้เรียนเข้าใจแนวคิดและประโยชน์ของ **Clean Architecture** ในการพัฒนา Backend ด้วยภาษา Go
- เพื่อให้ผู้เรียนสามารถออกแบบและพัฒนา REST API, WebSocket, MQTT Request‑Response, และระบบ Real‑Time ด้วย Go
- เพื่อให้ผู้เรียนสามารถบูรณาการระบบย่อยต่าง ๆ ได้แก่ Redis Cache, InfluxDB, Kafka Queue, และ gRPC
- เพื่อให้ผู้เรียนสามารถเขียนโค้ดที่มีประสิทธิภาพ ปลอดภัย และบำรุงรักษาได้ง่าย
- เพื่อฝึกปฏิบัติการใช้เครื่องมือและไลบรารียอดนิยม เช่น `chi`, `viper`, `cobra`, `gorm`, `zap`, `gomail`, `hermes`

### 1.2 กลุ่มเป้าหมาย
- นักพัฒนาซอฟต์แวร์ระดับกลางถึงสูงที่ต้องการพัฒนา Backend ด้วย Go
- วิศวกรระบบที่ต้องการออกแบบระบบที่มีการสื่อสารแบบ Real‑Time และประมวลผลข้อมูลปริมาณมาก
- ผู้ที่สนใจสถาปัตยกรรม Clean Architecture และการใช้งานร่วมกับ IoT, Microservices

### 1.3 ความรู้พื้นฐาน
- มีประสบการณ์การเขียนโปรแกรมภาษา Go ขั้นพื้นฐาน (ตัวแปร, ฟังก์ชัน, struct, interface)
- เข้าใจแนวคิดของการเขียน REST API และ HTTP protocol
- มีความรู้เกี่ยวกับฐานข้อมูลเชิงสัมพันธ์ (SQL) พอสมควร
- (แนะนำ) เคยใช้งาน Docker และเครื่องมือ DevOps มาก่อน

### 1.4 เนื้อหาโดยย่อ (เน้นวัตถุประสงค์และประโยชน์)
| หัวข้อ | รายละเอียด | ประโยชน์ |
|--------|------------|----------|
| Clean Architecture | แบ่งชั้น Model, Repository, Usecase, Delivery | ทำให้โค้ดเป็นระเบียบ ทดสอบง่าย เปลี่ยนเทคโนโลยีได้ยืดหยุ่น |
| REST API with JWT | สร้าง CRUD, Authentication, Email Verification, Forgot Password | ระบบยืนยันตัวตนที่ปลอดภัย พร้อมใช้งานจริง |
| WebSocket | สื่อสารแบบ Real‑Time ด้วย Hub/Client pattern | รองรับการอัปเดตข้อมูลทันที เช่น Dashboard, Chat |
| MQTT Request‑Response | ส่งคำขอและรอรับ response ผ่าน MQTT พร้อม Cache | เชื่อมต่อกับอุปกรณ์ IoT ได้อย่างมีประสิทธิภาพ ลดการเรียกซ้ำ |
| Redis Cache | เก็บ Session, Cache ข้อมูลผู้ใช้ และผลลัพธ์ MQTT | เพิ่มความเร็ว ลดภาระ Database และ MQTT |
| InfluxDB | บันทึกข้อมูล Time‑Series (Sensor Data, Alarm Logs) | วิเคราะห์แนวโน้มและประวัติข้อมูลเชิงเวลาได้ง่าย |
| Kafka Queue | ประมวลผลงานแบบอะซิงโครนัส (ส่งอีเมล, บันทึก Log) | แยกงานหนักออกจาก Request Cycle, ขยายระบบได้ |
| gRPC & Protobuf | สื่อสารระหว่างบริการด้วยประสิทธิภาพสูง, รองรับ Streaming | เหมาะกับ Microservices, ใช้ Bandwidth น้อยกว่า REST |

---

## 2. บทนำ

เอกสารนี้จัดทำขึ้นเพื่อเป็นคู่มือประกอบการพัฒนาแอปพลิเคชัน Backend ด้วยภาษา **Go** โดยใช้สถาปัตยกรรมแบบ **Clean Architecture (3‑Layer + Delivery)** ซึ่งประกอบด้วยชั้น Model, Repository, Usecase และ Delivery ภายในโปรเจกต์มีการบูรณาการระบบย่อยหลายตัว ได้แก่:

- **REST API** – ให้บริการ HTTP endpoints สำหรับการดำเนินการ CRUD, การยืนยันตัวตนด้วย JWT, การจัดการรหัสผ่าน และการแจ้งเตือนทางอีเมล
- **WebSocket** – รองรับการสื่อสารแบบ Real‑Time ระหว่าง client และ server ผ่าน WebSocket protocol
- **MQTT Request‑Response Pattern** – รองรับการส่งคำขอและรอรับ response ผ่าน MQTT พร้อมกลไกการ subscribe ชั่วคราวและ timeout
- **Redis Cache** – ใช้สำหรับเก็บ session JWT, cache ข้อมูลผู้ใช้ และลดการเรียก MQTT ซ้ำ
- **InfluxDB** – ใช้สำหรับบันทึกข้อมูลเชิงเวลา (time‑series) เช่น ข้อมูลจากเซ็นเซอร์ MQTT และประวัติแจ้งเตือน (alarm logs)
- **Alarm Processing** – ตรรกะการประเมินค่าเพื่อสร้างการแจ้งเตือนจากข้อมูลที่ได้รับ
- **Kafka** – ใช้เป็นระบบคิว (Queue Processor) สำหรับประมวลผลงานแบบอะซิงโครนัส เช่น การส่งอีเมล หรือการประมวลผลข้อมูลจำนวนมาก
- **gRPC** – Open‑Source Framework ประสิทธิภาพสูงจาก Google ที่ใช้สำหรับให้แอปพลิเคชันหรือเซอร์วิสต่างๆ สื่อสารกันผ่านเครือข่าย รองรับการสตรีม และใช้ Protocol Buffers เป็นรูปแบบข้อมูลแบบไบนารี

เอกสารนี้มีเป้าหมายเพื่อให้ผู้อ่านเข้าใจภาพรวมของระบบ รู้จักส่วนประกอบต่าง ๆ และสามารถนำไปปรับใช้ในโครงการจริงได้อย่างมีประสิทธิภาพ พร้อมทั้งมีเทมเพลตและ checklist สำหรับการวางแผนและตรวจสอบการทำงาน

---

## 3. บทนิยามศัพท์

### 3.1 Clean Architecture (สถาปัตยกรรมแบบสะอาด)
Clean Architecture เป็นรูปแบบการจัดโครงสร้างซอฟต์แวร์ที่แยกความรับผิดชอบออกเป็นชั้นต่าง ๆ เพื่อให้โค้ดมีความเป็นระเบียบ ทดสอบได้ง่าย และบำรุงรักษาได้ระยะยาว ในโปรเจกต์นี้แบ่งเป็น 4 ชั้นหลัก:

- **Model** – กำหนดโครงสร้างข้อมูล (entity) ที่ใช้ในแอปพลิเคชัน เช่น `User`, `Session`, `VerificationToken` โดยใช้ GORM เป็น ORM
- **Repository** – ทำหน้าที่ติดต่อกับแหล่งข้อมูล (Database, Redis, InfluxDB, MQTT) ผ่าน interface ทำให้ชั้นที่สูงกว่าไม่ต้องรู้รายละเอียดการจัดเก็บ
- **Usecase** – บรรจุตรรกะทางธุรกิจ (business logic) เช่น การแฮชรหัสผ่าน, การสร้าง JWT, การตรวจสอบข้อมูล, การจัดการคิวอีเมล
- **Delivery** – ชั้นที่ติดต่อกับโลกภายนอก ประกอบด้วย HTTP handlers (REST API), WebSocket handlers, และ background workers

### 3.2 MQTT Request‑Response Pattern
MQTT (Message Queuing Telemetry Transport) เป็นโปรโตคอลแบบ publish/subscribe ที่มีน้ำหนักเบา เหมาะกับ IoT และระบบที่มีข้อจำกัดด้านแบนด์วิดท์

**Request‑Response Pattern** คือรูปแบบที่ client ส่งคำขอ (publish) ไปยัง topic หนึ่ง และรอ response จาก server ผ่าน topic อีกตัวหนึ่ง โดยจะมีการ subscribe ชั่วคราวเพื่อรอ response และ unsubscribe อัตโนมัติเมื่อได้รับหรือ timeout

**ประโยชน์ที่ได้รับ**:
- รองรับการสื่อสารแบบสองทาง (bidirectional) บน MQTT
- ช่วยให้ระบบที่ใช้ MQTT สามารถทำงานแบบ request‑response ได้คล้าย HTTP
- ลดภาระการเก็บสถานะ (stateless) บนเซิร์ฟเวอร์

### 3.3 WebSocket
WebSocket เป็นโปรโตคอลที่ให้การสื่อสารแบบ full‑duplex ผ่านการเชื่อมต่อ TCP เดียว ช่วยให้ server สามารถส่งข้อมูลไปยัง client ได้ทันทีโดยไม่ต้องรอคำขอ (real‑time)

- **Hub** – เป็นโครงสร้างที่รวบรวมและจัดการ client connections ทั้งหมด รวมถึงการ broadcast ข้อความไปยัง client ต่างๆ
- **Client** – แทนการเชื่อมต่อ WebSocket แต่ละตัว มี channel สำหรับรับและส่งข้อความ

### 3.4 Kafka Queue Processor
Kafka เป็นแพลตฟอร์มกระจายสำหรับการสตรีมข้อมูล (distributed streaming platform) ทำหน้าที่เป็นระบบคิวข้อความ (message queue) ที่มีความทนทานสูงและรองรับปริมาณงานมหาศาล

**Queue Processor** คือกระบวนการที่อ่านข้อความจาก Kafka topic และประมวลผลตามลำดับ เช่น การส่งอีเมล, การบันทึก log, การวิเคราะห์ข้อมูล

**ประโยชน์**:
- แยกส่วนการผลิตและการบริโภคข้อมูล ทำให้ระบบมีความยืดหยุ่น
- รองรับการทำงานแบบอะซิงโครนัส ลดการรอคอยใน request cycle
- สามารถขยายจำนวน consumer ได้ตามโหลด

### 3.5 Redis Cache
Redis เป็นฐานข้อมูลแบบ in‑memory ที่ใช้สำหรับแคชข้อมูลและเก็บ session ในโปรเจกต์นี้ Redis ถูกใช้เพื่อ:
- เก็บ Refresh Token ของผู้ใช้
- แคชข้อมูลผู้ใช้เพื่อลดการ query ฐานข้อมูล
- แคชผลลัพธ์จากการเรียก MQTT เพื่อลดการเรียกซ้ำในระยะเวลาสั้น

### 3.6 InfluxDB
InfluxDB เป็นฐานข้อมูลเชิงเวลา (time‑series database) ที่ออกแบบมาเพื่อจัดเก็บและสอบถามข้อมูลที่มีป้ายกำกับเวลา เช่น ค่าจากเซ็นเซอร์, logs ในระบบนี้ InfluxDB ใช้บันทึก:
- ข้อมูล MQTT ที่ได้รับจากอุปกรณ์
- ประวัติการแจ้งเตือน (alarm logs) ที่เกิดจากตรรกะการประเมินค่า

### 3.7 Alarm Processing
Alarm Processing เป็นกระบวนการประเมินข้อมูลที่ได้รับ (เช่น ค่าอุณหภูมิ, ความชื้น) เทียบกับเกณฑ์ที่กำหนด เพื่อตัดสินใจว่าควรแจ้งเตือนหรือไม่ ตรรกะมักประกอบด้วย:
- การเปรียบเทียบค่ากับขีดจำกัดสูง/ต่ำ
- การตรวจจับการเปลี่ยนแปลงอย่างรวดเร็ว (rate of change)
- การหน่วงเวลา (debounce) เพื่อป้องกันการแจ้งเตือนซ้ำซ้อน

### 3.8 gRPC และ Protocol Buffers
- **Protocol Buffers (Protobuf)** เป็นภาษาสำหรับกำหนดโครงสร้างข้อมูล (interface definition language) ที่ใช้ในการส่งข้อมูลแบบไบนารี มีประสิทธิภาพสูงกว่า JSON/XML
- **gRPC** เป็นเฟรมเวิร์ก RPC ที่ใช้ Protobuf เป็นข้อมูลพื้นฐาน รองรับการสตรีมแบบต่างๆ และทำงานบน HTTP/2 ให้ประสิทธิภาพสูง

---

## 4. การออกแบบ Workflow

ในส่วนนี้จะอธิบายกระแสข้อมูล (data flow) ของระบบหลัก ๆ พร้อมตัวอย่างการทำงานจริง

### 4.1 Data Flow Diagram (Draw.io Style)

```
┌─────────────────────────────────────────────────────────────────────┐
│                          ผู้ใช้ / Client                            │
└───────────────┬─────────────────┬─────────────────┬───────────────┘
                │                 │                 │
                │ HTTP/REST       │ WebSocket       │ MQTT
                │ (chi)           │ (gorilla/ws)    │ (paho)
                ▼                 ▼                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        Delivery Layer                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐             │
│  │ REST Handler │  │ WS Handler   │  │ MQTT Handler │             │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘             │
└─────────┼─────────────────┼─────────────────┼─────────────────────┘
          │                 │                 │
          ▼                 ▼                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        Usecase Layer                               │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  Business Logic: JWT, Password Hash, Validation, Alarm      │   │
│  └─────────────────────────────────────────────────────────────┘   │
└─────────┬───────────────────────┬───────────────────────┬─────────┘
          │                       │                       │
          ▼                       ▼                       ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Repository Layer                              │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌──────────┐  ┌────────┐ │
│  │  GORM   │  │  Redis  │  │InfluxDB │  │  MQTT   │  │ Kafka  │ │
│  │ (MySQL) │  │ (Cache) │  │(Time-   │  │(Request-│  │(Queue) │ │
│  └─────────┘  └─────────┘  │ Series) │  │Response)│  └────────┘ │
│                             └─────────┘  └──────────┘            │
└─────────────────────────────────────────────────────────────────────┘
```

### 4.2 คำอธิบายขั้นตอนการทำงานหลัก

#### 4.2.1 MQTT Request‑Response Workflow
1. **Client** ส่งคำขอ (request) โดย publish ไปยัง topic เช่น `request/device/123` พร้อม payload ที่มี `correlationId` และ `replyTopic`
2. **MQTT Handler** ใน Delivery Layer รับข้อความและส่งต่อไปยัง Usecase
3. **Usecase** ประมวลผลคำขอ เช่น อ่านค่าจากเซ็นเซอร์ หรือเขียนข้อมูลลง InfluxDB
4. **Usecase** สร้าง response และ publish ไปยัง `replyTopic` ที่ client ระบุ พร้อม `correlationId` เดิม
5. **Client** ซึ่ง subscribe ไว้ที่ `replyTopic` จะได้รับ response และจับคู่กับคำขอผ่าน `correlationId`

**กลไกการจัดการ timeout**: Client จะตั้ง timer ไว้ หากครบเวลาที่กำหนดยังไม่ได้รับ response จะถือว่าคำขอหมดอายุและดำเนินการตามที่กำหนด (เช่น retry หรือแจ้ง error)

#### 4.2.2 Alarm Processing Workflow
1. **MQTT Client** รับข้อมูลจากเซ็นเซอร์ (เช่น อุณหภูมิ) ผ่าน MQTT subscription
2. ข้อมูลถูกส่งไปยัง **Alarm Usecase** ซึ่งจะประเมินค่าตามเกณฑ์ที่กำหนด (เช่น ถ้าอุณหภูมิ > 40°C เป็นเวลา 5 วินาที ให้แจ้งเตือน)
3. หากเข้าเงื่อนไข **Alarm Usecase** จะ:
   - บันทึกเหตุการณ์ลง InfluxDB (alarm log)
   - ส่งข้อความแจ้งเตือนผ่าน WebSocket ไปยัง client ที่เกี่ยวข้อง (real‑time)
   - อาจส่งอีเมลหรือ LINE notify ผ่าน Kafka queue
4. **WebSocket Hub** จะ broadcast ข้อความแจ้งเตือนไปยังทุก client ที่ subscribe ไว้

#### 4.2.3 WebSocket Real‑Time Communication
1. **Client** ทำ HTTP request เพื่อ upgrade เป็น WebSocket ที่ endpoint เช่น `/ws`
2. **WebSocket Handler** รับ connection และสร้าง `Client` object พร้อมลงทะเบียนใน `Hub`
3. **Hub** จะมี goroutine หลัก 2 ตัว:
   - `run()`: รับคำสั่ง register/unregister และ broadcast ข้อความไปยัง clients
   - `broadcast()`: ส่งข้อความที่ได้รับจาก Kafka หรือ internal event ไปยัง client ทั้งหมดหรือเฉพาะกลุ่ม
4. Client สามารถส่งข้อความไปยัง server ได้ ซึ่ง server จะประมวลผลและอาจตอบกลับ หรือส่งต่อไปยังระบบอื่น (เช่น บันทึก InfluxDB)

#### 4.2.4 Kafka Queue Processing
1. **Producer** (เช่น REST API เมื่อมีการส่งอีเมล) ส่งข้อความไปยัง Kafka topic `email-queue`
2. **Worker** (background process) จะ consume ข้อความจาก topic ดังกล่าว
3. **Worker** ประมวลผล (เช่น สร้างเนื้อหาอีเมลด้วย Hermes และส่งผ่าน Gomail)
4. เมื่อส่งสำเร็จ Worker อาจบันทึก log หรืออัปเดตสถานะในฐานข้อมูล

### 4.3 ตัวอย่างการใช้งานจริง (Use Case)

**กรณีศึกษา: ระบบตรวจสอบอุณหภูมิในคลังสินค้า**

- เซ็นเซอร์ IoT ส่งค่าอุณหภูมิทุก 5 วินาทีผ่าน MQTT ไปยัง topic `sensors/temperature`
- ระบบจะรับข้อมูลและประเมินแจ้งเตือน หากอุณหภูมิสูงกว่า 35°C ต่อเนื่อง 3 ครั้ง จะเกิด alarm
- Alarm จะถูกส่งไปยัง WebSocket client ที่เป็นหน้า dashboard และบันทึกใน InfluxDB เพื่อนำไปวิเคราะห์แนวโน้ม
- ผู้ดูแลสามารถเรียกดูข้อมูลย้อนหลังผ่าน REST API หรือ gRPC API
- หากต้องการแจ้งเตือนทางอีเมลเมื่อเกิด alarm รุนแรง ระบบจะใส่ message ลง Kafka queue และ worker จะทำการส่งอีเมลให้ผู้รับผิดชอบ

---

## 5. บทสรุป

### 5.1 ประโยชน์ที่ได้รับ
- **ความยืดหยุ่นสูง** – การแยกชั้นตาม Clean Architecture ทำให้สามารถเปลี่ยนเทคโนโลยีในแต่ละชั้นได้โดยไม่กระทบต่อชั้นอื่น (เช่น เปลี่ยน GORM เป็น Ent หรือเปลี่ยนฐานข้อมูล)
- **รองรับ Real‑Time** – WebSocket และ MQTT ช่วยให้ระบบสามารถตอบสนองแบบทันทีทันใด เหมาะกับ IoT และแอปพลิเคชันแบบ monitoring
- **ประสิทธิภาพและความทนทาน** – การใช้ Redis Cache ช่วยลดภาระระบบ, InfluxDB ช่วยจัดการข้อมูลเวลา, Kafka ช่วยแยกงานแบบอะซิงโครนัส ทำให้ระบบขยายตัวได้ดี
- **ความปลอดภัย** – มีระบบ JWT, การแฮชรหัสผ่าน, และ middleware ต่างๆ ช่วยป้องกันภัยคุกคามทั่วไป
- **ทดสอบง่าย** – การใช้ interface และ dependency injection ทำให้สามารถเขียน unit test ได้สะดวก

### 5.2 ข้อควรระวัง
- **การจัดการ Goroutine** – ต้องระวัง goroutine leak โดยเฉพาะใน WebSocket และ Kafka consumer ควรมีกลไก graceful shutdown
- **Timeout และ Retry** – ในการสื่อสารผ่าน MQTT และ Kafka ควรตั้งค่า timeout และ retry อย่างเหมาะสม เพื่อป้องกันการค้างของระบบ
- **การซิงโครไนซ์ข้อมูล** – การใช้ cache (Redis) อาจทำให้ข้อมูลไม่สดใหม่ ควรกำหนด TTL ให้เหมาะสม หรือมีกลไก invalidate cache เมื่อข้อมูลเปลี่ยนแปลง
- **ความซับซ้อน** – การมีหลาย component (MQTT, WebSocket, Kafka, InfluxDB) ทำให้ระบบมีความซับซ้อนในการดูแลและดีบัก ควรมีระบบ logging และ monitoring ที่ดี

### 5.3 ข้อดี
- รองรับการทำงานแบบเรียลไทม์และ asynchronous อย่างมีประสิทธิภาพ
- สถาปัตยกรรมที่ยืดหยุ่นและเป็นมาตรฐาน ทำให้ทีมพัฒนาเข้าใจและทำงานร่วมกันได้ง่าย
- ใช้เทคโนโลยีที่ได้รับความนิยมและมี community ขนาดใหญ่ (Go, gRPC, Redis, Kafka, InfluxDB)

### 5.4 ข้อเสีย
- การติดตั้งและกำหนดค่าเริ่มต้นค่อนข้างซับซ้อน ต้องจัดการหลาย service พร้อมกัน
- ต้องการความรู้พื้นฐานด้าน distributed systems และ concurrency เพื่อให้ใช้งานได้อย่างถูกต้อง
- การ debug ในระบบที่มีหลาย component อาจทำได้ยากกว่า monolithic แบบเดิม

### 5.5 ข้อห้าม (ถ้ามี)
- **ห้ามใช้ JWT แบบไม่ตั้งค่าหมดอายุ** – ควรตั้งค่า expiration ให้สั้น และใช้ refresh token หมุนเวียน
- **ห้ามเก็บข้อมูลสำคัญ (เช่น password) ใน cache โดยไม่เข้ารหัส** – Redis อาจถูกโจมตีได้ ควรเข้ารหัสข้อมูลที่ sensitive
- **ห้ามใช้ MQTT QoS 0 ในกรณีที่ต้องการความเชื่อถือได้** – ควรเลือก QoS 1 หรือ 2 สำหรับข้อมูลสำคัญ
- **ห้ามข้ามการ validate input** – ควรใช้ validator และ sanitize ทุกครั้ง เพื่อป้องกัน injection และ XSS

### 5.6 แหล่งอ้างอิง
1. เอกสารทางการของ Go: [https://go.dev/doc/](https://go.dev/doc/)
2. Clean Architecture โดย Robert C. Martin: [https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
3. MQTT Specification: [https://mqtt.org/mqtt-specification/](https://mqtt.org/mqtt-specification/)
4. gRPC Documentation: [https://grpc.io/docs/](https://grpc.io/docs/)
5. Apache Kafka Documentation: [https://kafka.apache.org/documentation/](https://kafka.apache.org/documentation/)
6. InfluxDB Documentation: [https://docs.influxdata.com/influxdb/](https://docs.influxdata.com/influxdb/)
7. Redis Documentation: [https://redis.io/documentation](https://redis.io/documentation)

---

## 6. คู่มือการปฏิบัติงาน

### 6.1 TASK LIST Template
ใช้สำหรับวางแผนและติดตามขั้นตอนการพัฒนาในแต่ละ sprint หรือ feature

| Task ID | ชื่องาน | รายละเอียด | ผู้รับผิดชอบ | สถานะ | หมายเหตุ |
|---------|--------|------------|-------------|--------|----------|
| T001 | ตั้งค่าโครงสร้างโปรเจกต์ | สร้างโฟลเดอร์ตาม Clean Architecture, ตั้งค่า go.mod | Dev A | Done | ใช้ Go 1.21+ |
| T002 | ติดตั้งและตั้งค่า Database | ติดตั้ง MySQL/PostgreSQL, สร้าง schema ตาม models | Dev B | Done | ใช้ GORM AutoMigrate |
| T003 | Implement JWT Authentication | สร้าง login/logout, middleware, Redis สำหรับ refresh token | Dev A | In Progress | ใช้ golang-jwt |
| T004 | สร้าง MQTT Client + Request-Response | implement pkg/mqtt พร้อม cache | Dev C | To Do | |
| T005 | Integrate InfluxDB | ตั้งค่า client และเขียน usecase สำหรับบันทึกข้อมูล | Dev C | To Do | |
| T006 | พัฒนา Alarm Logic | แปลงจาก iot.helper.ts มาเป็น Go struct และฟังก์ชัน | Dev A | To Do | |
| T007 | สร้าง WebSocket Hub และ Client | ใช้ gorilla/websocket, จัดการ connection | Dev B | To Do | |
| T008 | ตั้งค่า Kafka Producer/Consumer | กำหนด topic, implement worker สำหรับ email | Dev D | To Do | |
| T009 | สร้าง REST API endpoints | CRUD สำหรับ user, device, alarm config | Dev A | To Do | |
| T010 | ทดสอบและ Benchmark | ใช้ wrk, h2load, ghz, เขียน unit test | Dev D | To Do | |
| T011 | จัดทำเอกสาร | เขียน README, API docs, คู่มือนี้ | Dev D | In Progress | |

### 6.2 CHECKLIST Template
ใช้สำหรับตรวจสอบความสมบูรณ์ก่อนการ deploy หรือส่งมอบงาน

#### 6.2.1 การตั้งค่า Infrastructure
- [ ] ติดตั้งและตั้งค่า MySQL/PostgreSQL พร้อมกำหนด user, password, database
- [ ] ติดตั้ง Redis และกำหนดค่า maxmemory, policy
- [ ] ติดตั้ง InfluxDB และสร้าง bucket, organization
- [ ] ติดตั้ง Kafka และสร้าง topics ที่จำเป็น (email-queue, alarm-queue, log-queue)
- [ ] ตรวจสอบการเชื่อมต่อระหว่าง service ทั้งหมด (ping, telnet)

#### 6.2.2 การพัฒนา Code
- [ ] โค้ดผ่านการ go fmt, go vet, golint
- [ ] มี unit test ครอบคลุม usecase หลัก (อย่างน้อย 70%)
- [ ] มี integration test สำหรับ API endpoints
- [ ] การจัดการ error ครบถ้วน (ไม่ panic ใน production)
- [ ] มี logging ที่เพียงพอ (ใช้ zap)
- [ ] ตั้งค่า environment variables ทั้งหมดใน .env และมี .env.example

#### 6.2.3 ความปลอดภัย
- [ ] ใช้ HTTPS สำหรับ REST และ gRPC (TLS)
- [ ] JWT มี expiration และ refresh token mechanism
- [ ] Password ถูก hash ด้วย Argon2 หรือ bcrypt
- [ ] มี middleware สำหรับ CORS, Rate Limiting, XSS Prevention
- [ ] ตรวจสอบ input ทั้งหมดด้วย validator

#### 6.2.4 Performance
- [ ] มีการตั้งค่า connection pool สำหรับ database และ Redis
- [ ] ใช้ Redis cache สำหรับข้อมูลที่เรียกบ่อย
- [ ] มีการ benchmark และปรับแต่งให้เหมาะสม (ใช้ wrk, ghz)
- [ ] ตรวจสอบ goroutine leak (ใช้ pprof)

#### 6.2.5 Deployment
- [ ] สร้าง binary สำหรับ production (GOOS=linux GOARCH=amd64 go build)
- [ ] มี Dockerfile และ docker-compose สำหรับ development
- [ ] กำหนด health check endpoint
- [ ] มี graceful shutdown (handle SIGTERM)

### 6.3 การวิเคราะห์สาเหตุที่แท้จริง (RCA – Root Cause Analysis)

RCA เป็นกระบวนการค้นหาสาเหตุพื้นฐานของปัญหา เพื่อป้องกันไม่ให้เกิดซ้ำ ใช้เครื่องมือเช่น **5 Why** และ **Fishbone Diagram**

#### 6.3.1 แนวทางและขั้นตอนการวิเคราะห์
1. **กำหนดปัญหาให้ชัดเจน (Define the Problem)**  
   - ระบุสิ่งที่เกิดขึ้นจริง, ผลกระทบ, ขอบเขต เช่น “ระบบไม่สามารถส่งอีเมลยืนยันได้เป็นเวลา 10 นาที ส่งผลให้ผู้ใช้สมัครสมาชิกไม่สำเร็จ”
2. **รวบรวมข้อมูล (Gather Data)**  
   - ดู logs, metric, สอบถามผู้เกี่ยวข้อง, ตรวจสอบการตั้งค่า
3. **วิเคราะห์สาเหตุ (Analyze the Root Cause)**  
   - ใช้เทคนิค **5 Why** ถาม “ทำไม” ซ้ำ ๆ เพื่อเจาะลึก  
   - ใช้ **Fishbone Diagram** แยกสาเหตุเป็นด้านต่าง ๆ (คน, กระบวนการ, เทคโนโลยี, สภาพแวดล้อม)
4. **กำหนดแนวทางแก้ไข (Implement Solutions)**  
   - แก้ที่สาเหตุ ไม่ใช่แค่อาการ เช่น ถ้าสาเหตุคือ SMTP server timeout ให้เพิ่ม retry และปรับ timeout
5. **ติดตามผล (Monitor)**  
   - ตรวจสอบว่าปัญหาไม่กลับมา และประเมินประสิทธิภาพของแนวทางแก้ไข

#### 6.3.2 ตัวอย่างการวิเคราะห์ด้วย 5 Why
**ปัญหา**: ผู้ใช้ไม่ได้รับอีเมลยืนยันการลงทะเบียน

- **Why 1**: ทำไมอีเมลไม่ถูกส่ง?  
  → เพราะ worker ที่ส่งอีเมลไม่สามารถเชื่อมต่อกับ SMTP server ได้
- **Why 2**: ทำไม worker เชื่อมต่อ SMTP ไม่ได้?  
  → เพราะ SMTP server response timeout หลังจาก 5 วินาที
- **Why 3**: ทำไม SMTP server ถึง timeout?  
  → เพราะปริมาณอีเมลในช่วงนั้นสูงกว่า normal 5 เท่า (เกิด spike)
- **Why 4**: ทำไมถึงมี spike?  
  → เพราะมีการส่งอีเมลแจ้งเตือนจำนวนมากจาก alarm system พร้อมกัน
- **Why 5**: ทำไม alarm system ส่งพร้อมกัน?  
  → เพราะไม่มี rate limiting ในการส่งอีเมลแจ้งเตือน และไม่มี queue buffer

**Root Cause**: ขาด rate limiting และ buffer mechanism สำหรับการส่งอีเมลแจ้งเตือน

**แนวทางแก้ไข**:
- เพิ่ม rate limiter ใน worker (จำกัดจำนวนอีเมลต่อนาที)
- ใช้ Kafka เป็น buffer ก่อนส่งจริง
- ปรับปรุง retry logic ให้มีการ backoff

#### 6.3.3 ประโยชน์ของการใช้ RCA
- **แก้ไขปัญหาอย่างยั่งยืน** – ป้องกันไม่ให้เกิดซ้ำอีก
- **ลดต้นทุนและความสูญเสีย** – ไม่ต้องเสียเวลาแก้ปัญหาเดิมซ้ำ ๆ
- **ปรับปรุงกระบวนการ** – ช่วยให้เห็นจุดอ่อนในระบบและพัฒนาต่อไป

---

## 7. ภาคผนวก: เทมเพลตและตัวอย่างโค้ด

### 7.1 การตั้งค่า MQTT Client (พร้อม Cache)

```go
// pkg/mqtt/client.go
package mqtt

import (
    "context"
    "encoding/json"
    "errors"
    "time"
    "github.com/eclipse/paho.mqtt.golang"
    "github.com/google/uuid"
    "github.com/redis/go-redis/v9"
)

type Client struct {
    conn     mqtt.Client
    redis    *redis.Client
    cacheTTL time.Duration
}

// GetDataFromTopic ดึงข้อมูลจาก MQTT โดยใช้ Redis cache เพื่อลดการเรียกซ้ำ
// หากมีใน cache จะคืนค่าทันที มิฉะนั้นจะ subscribe แบบ request-response
func (c *Client) GetDataFromTopic(ctx context.Context, topic string) ([]byte, error) {
    // 1. ตรวจสอบ cache
    cached, err := c.redis.Get(ctx, "mqtt:"+topic).Bytes()
    if err == nil {
        return cached, nil // cache hit
    }

    // 2. สร้าง correlationId และ replyTopic
    corrID := uuid.New().String()
    replyTopic := "reply/" + corrID

    // 3. Subscribe ชั่วคราวเพื่อรอ response
    responseChan := make(chan []byte, 1)
    token := c.conn.Subscribe(replyTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
        responseChan <- msg.Payload()
    })
    if !token.WaitTimeout(5 * time.Second) {
        return nil, errors.New("subscribe timeout")
    }

    // 4. Publish request พร้อม correlationId และ replyTopic
    payload := map[string]interface{}{
        "correlationId": corrID,
        "replyTopic":    replyTopic,
        "data":          "your request payload",
    }
    data, _ := json.Marshal(payload)
    token = c.conn.Publish(topic, 1, false, data)
    if !token.WaitTimeout(5 * time.Second) {
        return nil, errors.New("publish timeout")
    }

    // 5. รอ response (timeout 10 วินาที)
    select {
    case resp := <-responseChan:
        // บันทึก cache
        c.redis.Set(ctx, "mqtt:"+topic, resp, c.cacheTTL)
        return resp, nil
    case <-time.After(10 * time.Second):
        return nil, errors.New("response timeout")
    }
}
```

**คำอธิบาย**:  
- ฟังก์ชันนี้จะพยายามดึงข้อมูลจาก Redis ก่อน ถ้าไม่มีจะทำ MQTT request-response  
- มีการตั้ง timeout ทั้งตอน subscribe, publish และรอ response เพื่อป้องกันการค้าง  
- เมื่อได้ response จะเก็บไว้ใน Redis เป็นเวลา `cacheTTL` เพื่อลดการเรียกซ้ำ

---

### 7.2 WebSocket Hub และ Client

```go
// internal/pkg/websocket/hub.go
package websocket

import "sync"

// Hub จัดการ client connections ทั้งหมด
type Hub struct {
    clients    map[*Client]bool          // ลงทะเบียน client
    broadcast  chan []byte               // channel สำหรับ broadcast ข้อความ
    register   chan *Client              // รับ client ใหม่
    unregister chan *Client              // รับ client ที่断开
    mu         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

// Run เป็น goroutine หลักสำหรับจัดการ events
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()
        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mu.Unlock()
        case message := <-h.broadcast:
            h.mu.RLock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mu.RUnlock()
        }
    }
}
```

```go
// internal/pkg/websocket/client.go
package websocket

import (
    "github.com/gorilla/websocket"
)

type Client struct {
    hub  *Hub
    conn *websocket.Conn
    send chan []byte       // buffered channel สำหรับส่งข้อความ
}

// readPump อ่านข้อความจาก WebSocket และส่งไปยัง Hub
func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }
        // ประมวลผล message (เช่น ส่งไปยัง Kafka หรือ MQTT)
        // ...
        // ตัวอย่าง: broadcast กลับไปยังทุกคน
        c.hub.broadcast <- message
    }
}

// writePump ส่งข้อความจาก channel send ไปยัง WebSocket
func (c *Client) writePump() {
    defer c.conn.Close()
    for message := range c.send {
        err := c.conn.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            break
        }
    }
}
```

---

### 7.3 Kafka Producer – การส่งอีเมลผ่าน Queue

```go
// internal/delivery/worker/email_worker.go
package worker

import (
    "context"
    "encoding/json"
    "github.com/IBM/sarama"
    "gomail"
)

type EmailTask struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

// EmailWorker consume messages from Kafka and send emails
func EmailWorker(brokers []string, topic string) {
    consumer, err := sarama.NewConsumer(brokers, nil)
    if err != nil {
        panic(err)
    }
    defer consumer.Close()

    partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
    if err != nil {
        panic(err)
    }
    defer partitionConsumer.Close()

    for msg := range partitionConsumer.Messages() {
        var task EmailTask
        if err := json.Unmarshal(msg.Value, &task); err != nil {
            // log error
            continue
        }
        // ส่งอีเมล
        m := gomail.NewMessage()
        m.SetHeader("From", "no-reply@example.com")
        m.SetHeader("To", task.To)
        m.SetHeader("Subject", task.Subject)
        m.SetBody("text/html", task.Body)

        d := gomail.NewDialer("smtp.example.com", 587, "user", "pass")
        if err := d.DialAndSend(m); err != nil {
            // log error, อาจ retry หรือส่งไป dead letter queue
        }
    }
}
```

---

### 7.4 การบันทึกข้อมูลลง InfluxDB

```go
// internal/modules/influxdb/usecase/usecase.go
package usecase

import (
    "context"
    "time"
    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InfluxUsecase struct {
    client influxdb2.Client
    org    string
    bucket string
}

func (u *InfluxUsecase) WriteSensorData(ctx context.Context, deviceID string, temperature float64, humidity float64) error {
    writeAPI := u.client.WriteAPI(u.org, u.bucket)
    p := influxdb2.NewPoint(
        "sensor",
        map[string]string{"device": deviceID},
        map[string]interface{}{
            "temperature": temperature,
            "humidity":    humidity,
        },
        time.Now(),
    )
    writeAPI.WritePoint(p)
    writeAPI.Flush()
    return nil
}
```

---

### 7.5 Middleware ตัวอย่าง (Logger + Recovery)

```go
// internal/delivery/rest/middleware/logger.go
package middleware

import (
    "net/http"
    "time"
    "go.uber.org/zap"
)

func Logger(logger *zap.Logger) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            next.ServeHTTP(w, r)
            logger.Info("request",
                zap.String("method", r.Method),
                zap.String("path", r.URL.Path),
                zap.Duration("latency", time.Since(start)),
                zap.String("ip", r.RemoteAddr),
            )
        })
    }
}
```

---

### 7.6 การใช้ Air สำหรับ Hot‑Reload

ติดตั้ง:
```bash
go install github.com/cosmtrek/air@latest
```

สร้างไฟล์ `.air.toml`:
```toml
root = "."
tmp_dir = "tmp"
[build]
  cmd = "go build -o ./tmp/main ./cmd/apiser"
  bin = "./tmp/main"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  include_dir = []
  include_file = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_error = false
[color]
  main = "magenta"
  watcher = "cyan"
  build = "yellow"
  runner = "green"
[log]
  time = false
[proxy]
  enabled = false
  proxy_port = 0
  app_port = 0
```

จากนั้นรัน `air` เพื่อเริ่มพัฒนาด้วย hot‑reload

---

## โครงสร้างโฟลเดอร์อ้างอิง

```bash
icmongolang/
├── api/
│   ├── cmd/
│   │   ├── apiser/                 # REST API หลัก
│   │   ├── websocket/              # WebSocket server
│   │   │   └── main.go
│   │   ├── initdata.go
│   │   ├── root.go
│   │   ├── serve.go
│   │   └── worker.go
│   ├── internal/
│   │   ├── websocket/
│   │   │   ├── delivery/ws/
│   │   │   │   ├── hub.go
│   │   │   │   ├── client.go
│   │   │   │   └── handler.go
│   │   │   ├── usecase/
│   │   │   │   └── ws_usecase.go
│   │   │   ├── repository/
│   │   │   │   └── ws_repo.go
│   │   │   └── models/
│   │   │       └── ws_models.go
│   │   └── pkg/
│   │       └── websocket/
│   │           ├── hub.go
│   │           ├── client.go
│   │           └── message.go
│   └── migrations/
│       └── 20250619_websocket_tables.sql
├── pkg/
│   ├── helpers/
│   │   ├── iot.go
│   │   └── format.go
│   ├── mqtt/
│   │   └── client.go
│   ├── influxdb/
│   │   └── client.go
│   └── redis/
│       └── redis_conn.go
└── internal/
    ├── mqtt/
    │   ├── delivery/http/
    │   │   ├── handler.go
    │   │   └── routes.go
    │   ├── presenter/
    │   │   └── presenter.go
    │   └── usecase/
    │       └── usecase.go
    ├── influxdb/
    │   ├── delivery/http/
    │   │   ├── handler.go
    │   │   └── routes.go
    │   ├── presenter/
    │   │   └── presenter.go
    │   └── usecase/
    │       └── usecase.go
    ├── alarm/
    │   ├── delivery/http/
    │   │   ├── handler.go
    │   │   └── routes.go
    │   ├── repository/
    │   │   └── alarm_log_repo.go
    │   └── usecase/
    │       └── usecase.go
    └── server/
        ├── handlers.go
        └── server.go
```

---
 

*****************************************************************************

