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
