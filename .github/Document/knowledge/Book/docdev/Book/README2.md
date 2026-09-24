# The Complete Go Programming Guide: From Basics to Production-Ready Applications

**โดย คงนคร จันทะคุณ**  
kongnakornjantakun@gmail.com  
*Version 4.0 – April 2026*

---

## 📖 บทนำ

ในยุคที่ซอฟต์แวร์มีความซับซ้อนมากขึ้นเรื่อย ๆ ภาษาโปรแกรมมิ่งที่เรียบง่าย มีประสิทธิภาพสูง และสามารถจัดการกับการทำงานพร้อมกันได้ดี กลายเป็นสิ่งที่นักพัฒนาต้องการอย่างยิ่ง ภาษา Go (หรือ Golang) ถือกำเนิดขึ้นจากความต้องการของ engineers ที่ Google ซึ่งเผชิญกับความท้าทายในการพัฒนาและบำรุงรักษาระบบขนาดใหญ่ที่มีการทำงานพร้อมกันสูง พวกเขาต้องการภาษาใหม่ที่ผสมผสานความรวดเร็วในการทำงานของภาษา C ความง่ายของภาษา Python และความสามารถในการจัดการ concurrency ที่ดีขึ้น

คู่มือเล่มนี้เกิดจากความตั้งใจที่จะรวบรวมองค์ความรู้เกี่ยวกับภาษา Go ตั้งแต่ระดับพื้นฐานจนถึงระดับมืออาชีพ ครอบคลุมทั้งไวยากรณ์พื้นฐาน การจัดการโปรเจกต์ด้วย Go Modules การทดสอบหน่วย การทำงานพร้อมกัน (concurrency) ไปจนถึงการออกแบบสถาปัตยกรรมระดับ Production และการประยุกต์ใช้ Domain-Driven Design (DDD) ร่วมกับ Go รวมถึงการเชื่อมต่อกับระบบภายนอกที่ใช้ในโลกแห่งความจริง เช่น Redis, RabbitMQ, MQTT, InfluxDB, WebSocket, SMS, LINE Notify และ Discord

คู่มือนี้ถูกออกแบบให้เป็น **ทั้งตำราเรียนและคู่มืออ้างอิง** โดยเน้นให้ผู้อ่านสามารถนำไปประยุกต์ใช้ได้ทันที ตั้งแต่การติดตั้ง การเขียนโปรแกรมพื้นฐาน ไปจนถึงการออกแบบสถาปัตยกรรมแบบ Clean Architecture และ Domain‑Driven Design (DDD) รวมถึงการเชื่อมต่อกับระบบภายนอกที่พบได้บ่อยในโลกแห่งความจริง

### วัตถุประสงค์
- ให้ผู้อ่านเข้าใจภาษา Go อย่างลึกซึ้ง ตั้งแต่ไวยากรณ์จนถึง concurrency
- เสนอแนวทางการจัดโครงสร้างโปรเจกต์สำหรับการผลิตจริง
- นำเสนอเทคนิคการทดสอบหน่วย (Unit Test) และการวัดประสิทธิภาพ
- แนะนำรูปแบบสถาปัตยกรรม Clean Architecture + DDD + CQRS
- สอนการผสาน Redis, RabbitMQ, MQTT, InfluxDB, WebSocket, SMS, LINE Notify, Discord
- จัดเตรียมเทมเพลตและ checklist ที่ช่วยให้ทีมทำงานเป็นระบบ
- ประกอบด้วยแผนภาพ (Mermaid) สำหรับอธิบายโครงสร้างและกระบวนการทำงาน
- ให้โค้ดตัวอย่างที่สามารถนำไปรันทดสอบได้จริง

### กลุ่มเป้าหมาย
- **ผู้เริ่มต้น** ที่ต้องการเรียนรู้ภาษา Go ตั้งแต่ศูนย์
- **นักพัฒนาที่เปลี่ยนภาษา** จากภาษาอื่นมาสู่ Go
- **นักพัฒนาที่ต้องการยกระดับ** สู่การเป็น Go Developer มืออาชีพ
- **สถาปนิกซอฟต์แวร์** ที่สนใจการออกแบบระบบด้วย Go

### วิธีการอ่าน
- หากยังไม่เคยเขียน Go มาก่อน ให้เริ่มจาก **ภาคที่ 1–3** เพื่อทำความเข้าใจพื้นฐาน
- หากต้องการออกแบบแอปพลิเคชันทันที ให้ข้ามไป **ภาคที่ 7–8** เพื่อศึกษา Clean Architecture และ DDD
- หากต้องการเชื่อมต่อกับระบบอื่น (ฐานข้อมูล time‑series, message queue, IoT) ให้ดู **ภาคที่ 9**

---

## 🧭 สารบัญ

### ภาคที่ 1: ปฐมบทกับการเขียนโปรแกรม
- บทที่ 1: ความรู้เบื้องต้นเกี่ยวกับการเขียนโปรแกรมคอมพิวเตอร์
- บทที่ 2: รู้จักกับภาษา Go
- บทที่ 3: พื้นฐานการใช้งาน Terminal
- บทที่ 4: เตรียมสภาพแวดล้อมสำหรับพัฒนา
- บทที่ 5: สร้างแอปพลิเคชันแรกของคุณ

### ภาคที่ 2: พื้นฐานภาษาและโครงสร้างข้อมูล
- บทที่ 6: ระบบเลขฐานสองและฐานสิบ
- บทที่ 7: เลขฐานสิบหก, ฐานแปด, ASCII, UTF8, Unicode และ Runes
- บทที่ 8: ตัวแปร, ค่าคงที่ และชนิดข้อมูลพื้นฐาน
- บทที่ 9: คำสั่งควบคุมการทำงาน
- บทที่ 10: ฟังก์ชัน
- บทที่ 11: แพคเกจและการนำเข้า
- บทที่ 12: การเริ่มต้นทำงานของแพคเกจ
- บทที่ 13: การสร้างชนิดข้อมูลใหม่ (Types)
- บทที่ 14: เมธอด (Methods)
- บทที่ 15: พอยน์เตอร์ (Pointer)
- บทที่ 16: อินเทอร์เฟซ (Interfaces)

### ภาคที่ 3: การจัดการโปรเจกต์และโครงสร้างข้อมูลขั้นสูง
- บทที่ 17: Go Modules - การจัดการโปรเจกต์สมัยใหม่
- บทที่ 18: Go Module Proxies
- บทที่ 19: การทดสอบหน่วย (Unit Tests)
- บทที่ 20: อาเรย์ (Arrays)
- บทที่ 21: สไลซ์ (Slices)
- บทที่ 22: แมพ (Maps)
- บทที่ 23: การจัดการข้อผิดพลาด (Errors)

### ภาคที่ 4: การพัฒนาแอปพลิเคชันเชิงปฏิบัติ
- บทที่ 24: ฟังก์ชันนิรนาม (Anonymous functions) และ Closure
- บทที่ 25: การจัดการข้อมูล JSON และ XML
- บทที่ 26: พื้นฐานการสร้าง HTTP Server
- บทที่ 27: Enum, Iota และ Bitmask
- บทที่ 28: วันที่และเวลา
- บทที่ 29: การจัดเก็บข้อมูล: ไฟล์และฐานข้อมูล
- บทที่ 30: การทำงานพร้อมกัน (Concurrency)
- บทที่ 31: การบันทึกเหตุการณ์ (Logging)
- บทที่ 32: เทมเพลต (Templates)
- บทที่ 33: การจัดการค่า Configuration

### ภาคที่ 5: สู่การเป็นนักพัฒนา Go มืออาชีพ
- บทที่ 34: การวัดประสิทธิภาพ (Benchmarks)
- บทที่ 35: สร้าง HTTP Client
- บทที่ 36: การวิเคราะห์โปรไฟล์ (Program Profiling)
- บทที่ 37: การจัดการ Context
- บทที่ 38: Generics - การเขียนโค้ดแบบยืดหยุ่น
- บทที่ 39: Go กับกระบวนทัศน์ OOP?
- บทที่ 40: การอัปเกรดหรือดาวน์เกรดเวอร์ชัน Go
- บทที่ 41: คำแนะนำในการออกแบบโค้ดที่ดี
- บทที่ 42: ชีทสรุป (Cheatsheet)

### ภาคที่ 6: เครื่องมือและไลบรารียอดนิยม
- บทที่ 43: chi, viper, cobra, zap และเครื่องมือสำคัญ
- บทที่ 44: GORM – ORM ทรงพลังสำหรับ Go
- บทที่ 45: การส่งอีเมลด้วย gomail และ hermes

### ภาคที่ 7: การออกแบบสถาปัตยกรรมและ Workflow
- บทที่ 46: Clean Architecture และโครงสร้างโปรเจกต์
- บทที่ 47: Blueprint สำหรับโปรเจกต์ Go ระดับ Production
- บทที่ 48: การออกแบบ Workflow และ Task Management

### ภาคที่ 8: Domain-Driven Design (DDD) กับ Go
- บทที่ 49: หลักการ DDD และการนำไปใช้ใน Go
- บทที่ 50: Aggregates, Event Storming และ CQRS
- บทที่ 51: การออกแบบบริการด้วย Go-DDD

### ภาคที่ 9: การผสานระบบภายนอกและคุณลักษณะเสริม
- บทที่ 52: Redis สำหรับ Cache และ Message Queue
- บทที่ 53: RabbitMQ – Message Broker มาตรฐานองค์กร
- บทที่ 54: MQTT สำหรับ IoT และระบบเรียลไทม์
- บทที่ 55: InfluxDB – Time‑Series Database
- บทที่ 56: WebSocket และ Socket.IO
- บทที่ 57: การส่ง SMS และ LINE Notify
- บทที่ 58: Discord Webhook สำหรับแจ้งเตือน

### ภาคที่ 10: เทมเพลต กระบวนการพัฒนา และตัวอย่างโค้ด
- บทที่ 59: ตัวอย่างโค้ดครบวงจร (Full‑stack Example)
- บทที่ 60: Task List Template
- บทที่ 61: Checklist Template
- บทที่ 62: แผนภาพการทำงาน (Workflow Diagram)
- บทที่ 63: mop Config – การจัดการ Configuration

---

## ภาคที่ 1: ปฐมบทกับการเขียนโปรแกรม

### บทที่ 1: ความรู้เบื้องต้นเกี่ยวกับการเขียนโปรแกรมคอมพิวเตอร์

#### 1.1 การเขียนโปรแกรมคืออะไร?
การเขียนโปรแกรม (Programming) คือกระบวนการสร้างชุดคำสั่งที่ใช้ควบคุมการทำงานของคอมพิวเตอร์ให้ทำงานตามที่เราต้องการ โดยใช้ภาษาเฉพาะที่คอมพิวเตอร์สามารถเข้าใจได้ ภาษาที่มนุษย์ใช้เขียนเรียกว่า "ภาษาคอมพิวเตอร์ระดับสูง" (High-Level Language) เช่น Go, Python, Java ซึ่งจากนั้นจะถูกแปลงเป็นภาษาเครื่อง (Machine Language) ที่เป็นเลขฐานสอง (0 และ 1) ที่ซีพียูสามารถประมวลผลได้

#### 1.2 โครงสร้างพื้นฐานของโปรแกรม
โปรแกรมคอมพิวเตอร์โดยทั่วไปประกอบด้วย:
- **ข้อมูล (Data)** : ตัวเลข, ข้อความ, รายการต่างๆ
- **การประมวลผล (Processing)** : การดำเนินการกับข้อมูล เช่น การคำนวณ การเปรียบเทียบ
- **การควบคุมการทำงาน (Control Flow)** : การตัดสินใจ (if-else), การวนซ้ำ (loop)
- **การจัดเก็บ (Storage)** : หน่วยความจำ, ไฟล์, ฐานข้อมูล
- **อินพุต/เอาท์พุต (I/O)** : การรับข้อมูลจากผู้ใช้ หรือแสดงผล

#### 1.3 ตัวแปลภาษาและคอมไพเลอร์
- **คอมไพเลอร์ (Compiler)** : แปลงซอร์สโค้ดทั้งหมดเป็นไฟล์ได้ก่อนรัน (เช่น Go, C)
- **อินเทอร์พรีเตอร์ (Interpreter)** : แปลงและรันทีละคำสั่ง (เช่น Python, JavaScript)

Go เป็นภาษาแบบคอมไพล์ (compiled) ซึ่งมีข้อดีคือทำงานเร็วและสร้างไฟล์ binary ที่รันได้ทันทีโดยไม่ต้องพึ่งพาสิ่งแวดล้อมอื่น (ยกเว้นระบบปฏิบัติการ)

#### 1.4 กระบวนทัศน์การเขียนโปรแกรม
Go รองรับการเขียนโปรแกรมแบบ:
- **Procedural** : ใช้ฟังก์ชันและลำดับขั้นตอน
- **Concurrent** : ทำงานพร้อมกันด้วย goroutine
- **Functional** (บางส่วน) : ฟังก์ชันเป็น first-class citizen
- **ไม่ใช่ OOP แบบคลาสสิก** : ใช้ struct และ interface แทน inheritance

#### 1.5 ขั้นตอนการพัฒนาโปรแกรม
1. เขียนซอร์สโค้ด (.go)
2. คอมไพล์ (go build)
3. ทดสอบรัน (go run)
4. แก้ไขข้อผิดพลาด (debug)
5. จัดการแพคเกจ (go mod)
6. ทดสอบหน่วย (go test)

---

### บทที่ 2: รู้จักกับภาษา Go

#### 2.1 ประวัติความเป็นมา
ภาษา Go (หรือ Golang) ถูกพัฒนาโดย Google เริ่มต้นในปี 2007 โดย Robert Griesemer, Rob Pike, และ Ken Thompson เปิดตัวเป็นโอเพนซอร์สในปี 2009 จุดประสงค์เพื่อแก้ปัญหาที่เกิดขึ้นในภาษา C++ และ Java ในระบบขนาดใหญ่ของ Google เช่น การคอมไพล์ที่ช้า, การจัดการการทำงานพร้อมกันที่ซับซ้อน, และความยุ่งยากในการบำรุงรักษา

#### 2.2 จุดเด่นของภาษา Go
- **เรียบง่ายและอ่านง่าย** : ไวยากรณ์กระชับ ไม่มีฟีเจอร์ที่ซับซ้อนเกินจำเป็น
- **คอมไพล์เร็ว** : สามารถคอมไพล์โปรเจกต์ขนาดใหญ่ได้ในไม่กี่วินาที
- **การจัดการหน่วยความจำอัตโนมัติ** : มี garbage collector ที่มีประสิทธิภาพ
- **Concurrency ระดับภาษา** : goroutine และ channel ทำให้เขียนโปรแกรม concurrent ได้ง่าย
- **Static typing** : ตรวจสอบชนิดข้อมูลตั้งแต่ตอนคอมไพล์ ป้องกันข้อผิดพลาด
- **เครื่องมือที่ครบครัน** : go fmt, go test, go mod, go vet, go doc
- **สามารถคอมไพล์ข้ามแพลตฟอร์ม** (cross-compile) ไปยัง Windows, Linux, macOS, ARM เป็นต้น

#### 2.3 โครงสร้างภาษา Go
- ไม่มี class แต่ใช้ struct และ method
- ไม่มี inheritance แต่ใช้ composition และ interface
- ไม่มี exception handling แต่ใช้ error return value
- มี garbage collection
- มี pointer แต่ไม่มี pointer arithmetic

#### 2.4 ตัวอย่าง Hello World ใน Go
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

#### 2.5 ใครใช้ Go บ้าง?
- **Google** : ระบบ backend, Kubernetes, Docker
- **Uber** : ระบบการจับคู่การเดินทาง
- **Dropbox** : เปลี่ยนจาก Python มาเป็น Go สำหรับระบบจัดเก็บไฟล์
- **Netflix** : ส่วนของ proxy และ caching
- **Cloudflare** : เครื่องมือโครงสร้างพื้นฐาน

---

### บทที่ 3: พื้นฐานการใช้งาน Terminal

#### 3.1 Terminal คืออะไร?
Terminal (หรือ command line, console) เป็นเครื่องมือที่ให้เราสั่งงานคอมพิวเตอร์ผ่านข้อความ แทนการใช้ GUI การพัฒนา Go มักใช้ terminal ในการรันคำสั่งต่างๆ เช่น go build, go run, go test, git เป็นต้น

#### 3.2 คำสั่งพื้นฐาน (Unix/Linux/macOS)
- `pwd` : แสดงไดเรกทอรีปัจจุบัน
- `ls` : แสดงรายการไฟล์ (ใช้ `ls -la` แสดงรายละเอียด)
- `cd <path>` : เปลี่ยนไดเรกทอรี
- `mkdir <name>` : สร้างโฟลเดอร์
- `touch <file>` : สร้างไฟล์
- `rm <file>` : ลบไฟล์ (ใช้ `rm -rf` ลบโฟลเดอร์)
- `cp <source> <dest>` : คัดลอกไฟล์
- `mv <source> <dest>` : ย้ายหรือเปลี่ยนชื่อ
- `cat <file>` : แสดงเนื้อหาไฟล์
- `echo <text>` : แสดงข้อความ
- `grep <pattern> <file>` : ค้นหาข้อความ

#### 3.3 คำสั่งพื้นฐานสำหรับ Windows (Command Prompt หรือ PowerShell)
- `cd` : เปลี่ยนไดเรกทอรี
- `dir` : แสดงรายการไฟล์
- `mkdir` : สร้างโฟลเดอร์
- `del` : ลบไฟล์
- `copy` : คัดลอก
- `move` : ย้าย
- `type` : แสดงเนื้อหาไฟล์

#### 3.4 การตั้งค่า PATH
หลังจากติดตั้ง Go แล้ว เราต้องให้ terminal สามารถเรียกใช้คำสั่ง `go` ได้ โดยเพิ่มไดเรกทอรีของ Go (เช่น `/usr/local/go/bin` หรือ `C:\Go\bin`) ลงในตัวแปรสภาพแวดล้อม PATH

#### 3.5 การใช้ go command
- `go version` : ตรวจสอบเวอร์ชัน Go
- `go env` : แสดงตัวแปรสภาพแวดล้อมของ Go
- `go run <file>` : คอมไพล์และรันโปรแกรม
- `go build` : คอมไพล์เป็น binary
- `go mod init` : เริ่มต้น go module
- `go get` : ดาวน์โหลดแพคเกจ (deprecated ใช้ `go install` หรือ `go mod download`)
- `go test` : รัน test

---

### บทที่ 4: เตรียมสภาพแวดล้อมสำหรับพัฒนา

#### 4.1 การติดตั้ง Go
1. ดาวน์โหลดจาก https://go.dev/dl/
2. เลือกเวอร์ชันที่เหมาะสมกับ OS ของคุณ
3. ติดตั้งตามขั้นตอน
4. ตรวจสอบการติดตั้ง: เปิด terminal แล้วพิมพ์ `go version`

#### 4.2 ตัวแปรสภาพแวดล้อมที่สำคัญ
- `GOROOT` : ตำแหน่งที่ติดตั้ง Go (มักตั้งค่าให้อัตโนมัติ)
- `GOPATH` : ตำแหน่ง workspace ของ Go (เลิกใช้แล้ว ถ้าใช้ modules)
- `GOBIN` : ตำแหน่งที่ติดตั้ง binaries ที่ `go install`
- `GOOS`/`GOARCH` : ใช้ cross-compile

#### 4.3 การเลือก Editor/IDE
- **Visual Studio Code** : มี extension "Go" โดย官方 ให้ linting, autocomplete, debugging
- **GoLand** : IDE เฉพาะของ JetBrains
- **Vim/Neovim** : ใช้ plugin vim-go
- **Sublime Text** : มี GoSublime

#### 4.4 การตั้งค่า Go Modules
Go Modules เป็นระบบจัดการ dependencies เริ่มใช้ตั้งแต่ Go 1.11 โดยค่าเริ่มต้นใช้งานได้ตั้งแต่ Go 1.16 เป็นต้นไป
```bash
mkdir myproject
cd myproject
go mod init example.com/myproject
```
ไฟล์ `go.mod` จะถูกสร้างขึ้น

#### 4.5 workspace structure (แบบเดิมกับแบบใหม่)
- **แบบ GOPATH** (ก่อน Go 1.11): ต้องวางโค้ดใน `$GOPATH/src`
- **แบบ Modules** (ปัจจุบัน): สามารถสร้างโปรเจกต์ได้ทุกที่ โดยมี go.mod กำกับ

#### 4.6 การติดตั้งเครื่องมือเสริม
- **gopls** : language server (ติดตั้งโดยอัตโนมัติใน VS Code)
- **dlv** : delve debugger (`go install github.com/go-delve/delve/cmd/dlv@latest`)
- **golangci-lint** : linter (`go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`)

---

### บทที่ 5: สร้างแอปพลิเคชันแรกของคุณ

#### 5.1 สร้างโปรเจกต์
```bash
mkdir hello
cd hello
go mod init hello
```

#### 5.2 เขียนโค้ด
สร้างไฟล์ `main.go`:
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```

#### 5.3 รันโปรแกรม
```bash
go run main.go
```
หรือคอมไพล์แล้วรัน:
```bash
go build -o hello
./hello   # (Linux/macOS) หรือ hello.exe (Windows)
```

#### 5.4 เพิ่มฟังก์ชัน
```go
package main

import "fmt"

func greet(name string) string {
    return "Hello, " + name
}

func main() {
    message := greet("Gopher")
    fmt.Println(message)
}
```
ทดสอบ `go run main.go` ควรได้ "Hello, Gopher"

#### 5.5 การอ่านค่าจาก command line
```go
package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Please provide a name")
        return
    }
    name := os.Args[1]
    fmt.Printf("Hello, %s!\n", name)
}
```
รัน: `go run main.go Somchai`

#### 5.6 การจัดรูปแบบโค้ด
ใช้คำสั่ง `go fmt` เพื่อจัดรูปแบบโค้ดให้เป็นมาตรฐาน:
```bash
go fmt ./...
```

#### 5.7 ข้อผิดพลาดที่พบบ่อย
- ใช้ตัวแปรแต่ไม่ได้ใช้ -> Go จะคอมไพล์ไม่ผ่าน (ยกเว้นตัวแปร _)
- import แพคเกจแต่ไม่ได้ใช้ -> ต้องลบออกหรือใช้ _
- ใส่เครื่องหมาย `,` หรือ `;` ไม่ถูกต้อง

---

## ภาคที่ 2: พื้นฐานภาษาและโครงสร้างข้อมูล

*(หมายเหตุ: เนื้อหาในภาคที่ 2–5 ได้ถูกนำมาจาก BOOK.md ซึ่งครอบคลุมพื้นฐานภาษา Go อย่างละเอียด เนื่องจากพื้นที่จำกัด จึงขอนำเสนอเฉพาะสารบัญและเนื้อหาสำคัญบางส่วนเพื่อให้เห็นภาพรวม)*

- บทที่ 6: ระบบเลขฐานสองและฐานสิบ
- บทที่ 7: เลขฐานสิบหก, ฐานแปด, ASCII, UTF8, Unicode และ Runes
- บทที่ 8: ตัวแปร, ค่าคงที่ และชนิดข้อมูลพื้นฐาน
- บทที่ 9: คำสั่งควบคุมการทำงาน
- บทที่ 10: ฟังก์ชัน
- บทที่ 11: แพคเกจและการนำเข้า
- บทที่ 12: การเริ่มต้นทำงานของแพคเกจ
- บทที่ 13: การสร้างชนิดข้อมูลใหม่ (Types)
- บทที่ 14: เมธอด (Methods)
- บทที่ 15: พอยน์เตอร์ (Pointer)
- บทที่ 16: อินเทอร์เฟซ (Interfaces)

---

## ภาคที่ 3: การจัดการโปรเจกต์และโครงสร้างข้อมูลขั้นสูง

- บทที่ 17: Go Modules - การจัดการโปรเจกต์สมัยใหม่
- บทที่ 18: Go Module Proxies
- บทที่ 19: การทดสอบหน่วย (Unit Tests)
- บทที่ 20: อาเรย์ (Arrays)
- บทที่ 21: สไลซ์ (Slices)
- บทที่ 22: แมพ (Maps)
- บทที่ 23: การจัดการข้อผิดพลาด (Errors)

---

## ภาคที่ 4: การพัฒนาแอปพลิเคชันเชิงปฏิบัติ

- บทที่ 24: ฟังก์ชันนิรนาม (Anonymous functions) และ Closure
- บทที่ 25: การจัดการข้อมูล JSON และ XML
- บทที่ 26: พื้นฐานการสร้าง HTTP Server
- บทที่ 27: Enum, Iota และ Bitmask
- บทที่ 28: วันที่และเวลา
- บทที่ 29: การจัดเก็บข้อมูล: ไฟล์และฐานข้อมูล
- บทที่ 30: การทำงานพร้อมกัน (Concurrency)
- บทที่ 31: การบันทึกเหตุการณ์ (Logging)
- บทที่ 32: เทมเพลต (Templates)
- บทที่ 33: การจัดการค่า Configuration

---

## ภาคที่ 5: สู่การเป็นนักพัฒนา Go มืออาชีพ

- บทที่ 34: การวัดประสิทธิภาพ (Benchmarks)
- บทที่ 35: สร้าง HTTP Client
- บทที่ 36: การวิเคราะห์โปรไฟล์ (Program Profiling)
- บทที่ 37: การจัดการ Context
- บทที่ 38: Generics - การเขียนโค้ดแบบยืดหยุ่น
- บทที่ 39: Go กับกระบวนทัศน์ OOP?
- บทที่ 40: การอัปเกรดหรือดาวน์เกรดเวอร์ชัน Go
- บทที่ 41: คำแนะนำในการออกแบบโค้ดที่ดี
- บทที่ 42: ชีทสรุป (Cheatsheet)

---

## ภาคที่ 6: เครื่องมือและไลบรารียอดนิยม

### บทที่ 43: chi, viper, cobra, zap และเครื่องมือสำคัญ

#### 43.1 chi – เราเตอร์และมิดเดิลแวร์

[chi](https://github.com/go-chi/chi) เป็น lightweight router ที่มีประสิทธิภาพสูง รองรับ middleware และเข้ากันได้กับ net/http มาตรฐาน

**การติดตั้ง**
```bash
go get github.com/go-chi/chi/v5
```

**ตัวอย่างพื้นฐาน**
```go
package main

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    r := chi.NewRouter()
    
    // ใช้ middleware พื้นฐาน
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    
    // กำหนด routes
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello World"))
    })
    
    // route พร้อมพารามิเตอร์
    r.Get("/users/{id}", getUser)
    
    // group routes
    r.Route("/api", func(r chi.Router) {
        r.Get("/users", listUsers)
        r.Post("/users", createUser)
        
        // sub-router พร้อม middleware เฉพาะ
        r.Route("/admin", func(r chi.Router) {
            r.Use(adminOnly)
            r.Get("/dashboard", adminDashboard)
        })
    })
    
    http.ListenAndServe(":8080", r)
}

func getUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    w.Write([]byte("User ID: " + id))
}
```

**มิดเดิลแวร์ที่ใช้บ่อย**
- `middleware.Logger` – บันทึก request
- `middleware.Recoverer` – จับ panic
- `middleware.Timeout` – กำหนด timeout
- `middleware.Compress` – บีบอัด response
- `middleware.RealIP` – ดึง IP จริงจาก proxy

#### 43.2 viper – การจัดการ configuration

[viper](https://github.com/spf13/viper) รองรับหลายรูปแบบ (JSON, YAML, ENV, flags) และสามารถโหลดจากไฟล์, environment variables, หรือ remote system

**การติดตั้ง**
```bash
go get github.com/spf13/viper
```

**ตัวอย่างการใช้งาน**
```go
package config

import (
    "log"
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    JWT      JWTConfig
}

type ServerConfig struct {
    Port int
    Mode string
}

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    Name     string
}

func LoadConfig() (*Config, error) {
    viper.SetConfigName("config")      // ชื่อไฟล์ (ไม่รวมนามสกุล)
    viper.SetConfigType("yaml")        // yaml, json, toml, etc.
    viper.AddConfigPath(".")           // path ที่ค้นหา
    viper.AddConfigPath("/etc/app/")
    viper.AutomaticEnv()               // อ่านจาก environment variables
    
    // กำหนดค่า default
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("server.mode", "debug")
    
    // อ่านไฟล์
    if err := viper.ReadInConfig(); err != nil {
        log.Printf("Config file not found: %v", err)
    }
    
    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }
    
    return &cfg, nil
}
```

**ไฟล์ config.yaml ตัวอย่าง**
```yaml
server:
  port: 8080
  mode: release

database:
  host: localhost
  port: 3306
  user: root
  password: secret
  name: mydb

redis:
  addr: localhost:6379
  password: ""
  db: 0

jwt:
  secret: "your-secret-key"
  access_expiry: 15m
  refresh_expiry: 7d
```

#### 43.3 cobra – การสร้าง CLI

[cobra](https://github.com/spf13/cobra) เป็นไลบรารีสำหรับสร้าง command-line application รองรับ commands, flags, และ subcommands

**การติดตั้ง**
```bash
go get -u github.com/spf13/cobra/cobra
```

**การเริ่มต้นโปรเจกต์**
```bash
cobra init --pkg-name mycli
```

**ตัวอย่างการเพิ่ม command**
```go
// cmd/root.go
package cmd

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "mycli",
    Short: "My CLI application",
    Long:  "A sample CLI built with cobra",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Hello from CLI")
    },
}

// cmd/serve.go
var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start the server",
    Run: func(cmd *cobra.Command, args []string) {
        port, _ := cmd.Flags().GetInt("port")
        fmt.Printf("Starting server on port %d\n", port)
    },
}

func init() {
    serveCmd.Flags().IntP("port", "p", 8080, "port to listen on")
    rootCmd.AddCommand(serveCmd)
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

**การใช้งาน**
```bash
go run main.go serve --port 3000
```

#### 43.4 gorm – ORM

[gorm](https://gorm.io) เป็น ORM ที่มีฟีเจอร์ครบถ้วน: associations, hooks, preloading, transactions, etc. (รายละเอียดในบทที่ 44)

#### 43.5 validator – การตรวจสอบข้อมูล

[validator](https://github.com/go-playground/validator) ใช้ struct tags ในการกำหนด validation rules

**การติดตั้ง**
```bash
go get github.com/go-playground/validator/v10
```

**ตัวอย่าง**
```go
package main

import (
    "fmt"
    "github.com/go-playground/validator/v10"
)

type RegisterRequest struct {
    Name     string `validate:"required,min=3,max=50"`
    Email    string `validate:"required,email"`
    Password string `validate:"required,min=8"`
    Age      int    `validate:"gte=18,lte=99"`
}

func main() {
    validate := validator.New()
    
    req := RegisterRequest{
        Name:     "Jo",
        Email:    "invalid",
        Password: "short",
        Age:      16,
    }
    
    err := validate.Struct(req)
    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            fmt.Printf("Field %s failed on tag %s\n", err.Field(), err.Tag())
        }
    }
}
```

**Tags ที่ใช้บ่อย**
- `required` – ต้องมีค่า
- `email` – รูปแบบ email
- `min`, `max` – ขนาดต่ำสุด/สูงสุด
- `gte`, `lte` – มากกว่าหรือเท่ากับ, น้อยกว่าหรือเท่ากับ
- `oneof=red blue` – ต้องเป็นหนึ่งในค่าที่กำหนด
- `uuid` – ต้องเป็น UUID
- `url` – ต้องเป็น URL

#### 43.6 jwt – การจัดการ JWT

[jwt-go](https://github.com/golang-jwt/jwt) เป็นไลบรารีมาตรฐานสำหรับ JWT

**การติดตั้ง**
```bash
go get github.com/golang-jwt/jwt/v5
```

**ตัวอย่างการสร้างและตรวจสอบ token**
```go
package jwtutil

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID uint `json:"user_id"`
    jwt.RegisteredClaims
}

var secretKey = []byte("your-secret-key")

func GenerateAccessToken(userID uint) (string, error) {
    claims := Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, jwt.ErrInvalidKey
}
```

#### 43.7 zap – structured logging

[zap](https://github.com/uber-go/zap) เป็น logger ที่มีความเร็วสูงและรองรับ structured logging

**การติดตั้ง**
```bash
go get go.uber.org/zap
```

**ตัวอย่าง**
```go
package logger

import (
    "go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger(mode string) {
    var err error
    if mode == "production" {
        Log, err = zap.NewProduction()
    } else {
        Log, err = zap.NewDevelopment()
    }
    if err != nil {
        panic(err)
    }
    defer Log.Sync()
}

func main() {
    InitLogger("development")
    defer Log.Sync()
    
    Log.Info("User logged in",
        zap.String("user", "john"),
        zap.Int("id", 123),
        zap.Duration("duration", time.Second*2),
    )
    
    Log.Error("Database error",
        zap.Error(err),
        zap.String("query", "SELECT * FROM users"),
    )
}
```

#### 43.8 gomail – การส่งอีเมล

[gomail](https://github.com/go-gomail/gomail) เป็นไลบรารีที่ใช้งานง่ายสำหรับ SMTP

**การติดตั้ง**
```bash
go get gopkg.in/gomail.v2
```

**ตัวอย่าง**
```go
package mail

import (
    "gopkg.in/gomail.v2"
)

type Mailer struct {
    dialer *gomail.Dialer
    from   string
}

func NewMailer(host string, port int, user, pass, from string) *Mailer {
    return &Mailer{
        dialer: gomail.NewDialer(host, port, user, pass),
        from:   from,
    }
}

func (m *Mailer) Send(to, subject, body string) error {
    msg := gomail.NewMessage()
    msg.SetHeader("From", m.from)
    msg.SetHeader("To", to)
    msg.SetHeader("Subject", subject)
    msg.SetBody("text/html", body)
    
    return m.dialer.DialAndSend(msg)
}
```

#### 43.9 hermes – สร้าง HTML email ที่สวยงาม

[hermes](https://github.com/matcornic/hermes) ใช้สร้าง email template แบบ responsive

**การติดตั้ง**
```bash
go get github.com/matcornic/hermes/v2
```

**ตัวอย่าง**
```go
package email

import (
    "github.com/matcornic/hermes/v2"
)

type EmailGenerator struct {
    h hermes.Hermes
}

func NewEmailGenerator(appURL, appName string) *EmailGenerator {
    h := hermes.Hermes{
        Product: hermes.Product{
            Name: appName,
            Link: appURL,
            Logo: appURL + "/logo.png",
        },
    }
    return &EmailGenerator{h: h}
}

func (g *EmailGenerator) WelcomeEmail(name, verifyURL string) (string, error) {
    email := hermes.Email{
        Body: hermes.Body{
            Name: name,
            Intros: []string{
                "Welcome to our platform!",
            },
            Actions: []hermes.Action{
                {
                    Instructions: "Please click below to verify your email address:",
                    Button: hermes.Button{
                        Color: "#22BC66",
                        Text:  "Verify Email",
                        Link:  verifyURL,
                    },
                },
            },
            Outros: []string{
                "If you didn't sign up, you can ignore this email.",
            },
        },
    }
    
    return g.h.GenerateHTML(email)
}
```

#### 43.10 air – hot-reload

[air](https://github.com/cosmtrek/air) ใช้สำหรับ reload อัตโนมัติเมื่อไฟล์เปลี่ยนแปลง เหมาะสำหรับการพัฒนา

**การติดตั้ง**
```bash
go install github.com/cosmtrek/air@latest
```

**การใช้งาน**
สร้างไฟล์ `.air.toml` ใน root ของโปรเจกต์ (หรือใช้ default) แล้วรัน:
```bash
air
```

**ตัวอย่าง .air.toml**
```toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/main ."
  bin = "./tmp/main"
  include_ext = ["go", "tpl", "tmpl", "html"]
  exclude_dir = ["assets", "tmp", "vendor"]
  delay = 1000
  stop_on_error = true
  send_interrupt = false
  kill_delay = 500
```

---

### บทที่ 44: GORM – ORM ทรงพลังสำหรับ Go

#### 44.1 บทนำ

GORM (Go Object Relational Mapping) เป็น ORM ที่ได้รับความนิยมสูงสุดในภาษา Go ช่วยลดความซับซ้อนในการจัดการฐานข้อมูล ด้วยการแมป struct กับตารางโดยอัตโนมัติ พร้อมฟังก์ชัน CRUD ที่ใช้งานง่าย การจัดการ transaction แบบมีระบบ และฟีเจอร์ขั้นสูงอย่าง query cache, queue processor

#### 44.2 การติดตั้งและการตั้งค่าพื้นฐาน

```go
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres // หรือ driver อื่นตามที่ใช้
```

**ตัวอย่างการเชื่อมต่อ PostgreSQL**

```go
package main

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
)

func main() {
    dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Bangkok"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect database")
    }

    // ใช้ db ต่อไป
}
```

#### 44.3 กำหนด Model (Entity)

```go
type User struct {
    ID        uint           `gorm:"primaryKey"`
    Name      string         `gorm:"size:100;not null"`
    Email     string         `gorm:"uniqueIndex;size:100;not null"`
    Age       int            `gorm:"default:0"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

#### 44.4 CRUD Operations

**Create**

```go
user := User{Name: "สมชาย", Email: "somchai@example.com", Age: 30}
result := db.Create(&user) // สร้าง record ใหม่
fmt.Println(user.ID)       // คืนค่า ID ที่ถูกสร้าง
fmt.Println(result.Error)  // error ถ้ามี
```

**Read**

```go
// ดึง record แรกที่ตรงเงื่อนไข
var user User
db.First(&user, 1)                 // by primary key
db.First(&user, "email = ?", "somchai@example.com")

// ดึงทั้งหมด
var users []User
db.Find(&users)

// พร้อมเงื่อนไข
db.Where("age > ?", 20).Find(&users)
db.Where(&User{Name: "สมชาย"}).Find(&users)
```

**Update**

```go
// อัปเดต single column
db.Model(&user).Update("Name", "สมชาย ใหม่")

// อัปเดตหลาย columns ด้วย struct (ไม่สนใจ zero values)
db.Model(&user).Updates(User{Name: "สมชาย ใหม่", Age: 31})

// อัปเดตหลาย columns ด้วย map
db.Model(&user).Updates(map[string]interface{}{"name": "สมชาย ใหม่", "age": 31})
```

**Delete**

```go
// soft delete (ถ้ามี gorm.DeletedAt)
db.Delete(&user, 1)

// hard delete
db.Unscoped().Delete(&user, 1)
```

#### 44.5 SessionFactory

Session factory เป็นรูปแบบการสร้าง `*gorm.DB` ที่มี configuration คงที่ (เช่น logging, connection pool) และสามารถสร้าง session ใหม่สำหรับแต่ละ request หรือ transaction

**ตัวอย่าง session factory**

```go
package db

import (
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "log"
    "time"
)

type SessionFactory struct {
    db *gorm.DB
}

func NewSessionFactory(dsn string) (*SessionFactory, error) {
    // ตั้งค่า logger และ connection pool
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold: time.Second,
            LogLevel:      logger.Info,
            Colorful:      true,
        },
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: newLogger,
        NowFunc: func() time.Time { return time.Now().UTC() },
    })
    if err != nil {
        return nil, err
    }

    // ตั้งค่า connection pool
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    return &SessionFactory{db: db}, nil
}

// NewSession คืนค่า session ใหม่สำหรับการทำ transaction หรือ query แบบแยก context
func (sf *SessionFactory) NewSession() *gorm.DB {
    return sf.db.Session(&gorm.Session{})
}
```

**การใช้งาน**

```go
factory, _ := db.NewSessionFactory(dsn)

// สร้าง session ใหม่สำหรับ request นี้
session := factory.NewSession()
var user User
session.First(&user, 1)

// เมื่อต้องการ transaction
session.Transaction(func(tx *gorm.DB) error {
    // ใช้ tx แทน session
    return nil
})
```

#### 44.6 การใช้ GORM Transaction เพื่อ Rollback

GORM มีฟังก์ชัน `db.Transaction` ที่ช่วยให้เราสามารถรวมหลายคำสั่ง SQL ไว้ใน transaction เดียวกันได้ โดยหากฟังก์ชันที่ส่งเข้าไปคืนค่า `error` GORM จะทำการ rollback โดยอัตโนมัติ ถ้าคืน `nil` จะ commit

**ตัวอย่างการใช้งาน Transaction**

```go
func PlaceOrder(db *gorm.DB, userID uint, items []CartItem) error {
    // เริ่ม transaction
    return db.Transaction(func(tx *gorm.DB) error {
        // 1. คำนวณราคารวม และตรวจสอบสต็อกพร้อม lock
        var total float64
        for _, item := range items {
            var stock Stock
            // Lock แถว stock เพื่อป้องกัน race condition
            if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
                Where("product_id = ?", item.ProductID).
                First(&stock).Error; err != nil {
                return err // สินค้าไม่มีในระบบ
            }
            if stock.Quantity < item.Quantity {
                return errors.New("สินค้าไม่พอ")
            }
            // หักสต็อก (จะบันทึกภายหลัง)
            stock.Quantity -= item.Quantity
            if err := tx.Save(&stock).Error; err != nil {
                return err
            }
            // คำนวณราคารวม
            total += getPrice(item.ProductID) * float64(item.Quantity)
        }

        // 2. สร้าง order
        order := Order{UserID: userID, Total: total}
        if err := tx.Create(&order).Error; err != nil {
            return err
        }

        // 3. สร้าง receipt
        receipt := Receipt{OrderID: order.ID, Amount: total}
        if err := tx.Create(&receipt).Error; err != nil {
            return err
        }

        // ทุกอย่างสำเร็จ -> commit อัตโนมัติ
        return nil
    })
}
```

#### 44.7 Query Cache

GORM เองไม่มี query cache ในตัว แต่สามารถใช้ plugin หรือจัดการเองผ่าน Redis หรือ memory cache

**ตัวอย่างการใช้ Redis cache กับ GORM**

```go
import (
    "context"
    "encoding/json"
    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"
    "time"
)

type CachedDB struct {
    db    *gorm.DB
    cache *redis.Client
}

func (c *CachedDB) FirstWithCache(dest interface{}, conds ...interface{}) error {
    // สร้าง cache key จากเงื่อนไข
    key := fmt.Sprintf("query:%v", conds)

    // พยายามอ่านจาก cache
    val, err := c.cache.Get(context.Background(), key).Result()
    if err == nil {
        // พบใน cache
        return json.Unmarshal([]byte(val), dest)
    }

    // ไม่พบใน cache, query ฐานข้อมูล
    if err := c.db.First(dest, conds...).Error; err != nil {
        return err
    }

    // บันทึกผลลัพธ์ลง cache (serialize)
    data, _ := json.Marshal(dest)
    c.cache.Set(context.Background(), key, data, 5*time.Minute)
    return nil
}
```

---

### บทที่ 45: การส่งอีเมลด้วย gomail และ hermes

*เนื้อหาดูได้ในบทที่ 43 (gomail) และ 43.9 (hermes)*

---

## ภาคที่ 7: การออกแบบสถาปัตยกรรมและ Workflow

### บทที่ 46: Clean Architecture และโครงสร้างโปรเจกต์

#### 46.1 โครงสร้างโปรเจกต์แบบ Clean Architecture

โครงสร้างที่แบ่งเป็น 3 ชั้นหลัก:
- **Delivery** – รับและส่งข้อมูล (HTTP handlers, gRPC, CLI)
- **Usecase** – business logic (interactors)
- **Repository** – การเข้าถึงข้อมูล (database, external API)

นอกจากนี้ยังมี **Models** (Entities) ที่ใช้ร่วมกันทุกชั้น

**โครงสร้างโฟลเดอร์ตัวอย่าง**
```
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/            # การตั้งค่า
│   ├── delivery/
│   │   ├── http/          # HTTP handlers
│   │   │   ├── handler.go
│   │   │   ├── middleware.go
│   │   │   └── routes.go
│   │   └── cli/           # (optional) CLI commands
│   ├── models/            # entities / DTOs
│   ├── repository/        # implementations
│   │   ├── user_repo.go
│   │   ├── user_repo_mock.go (สำหรับ test)
│   │   └── redis/         # redis implementation
│   └── usecase/           # business logic
│       ├── user_usecase.go
│       └── auth_usecase.go
├── pkg/                   # reusable packages
│   ├── jwt/
│   ├── mail/
│   └── redis/
├── go.mod
└── config.yaml
```

#### 46.2 Models (Entities)

**models/user.go**
```go
package models

import "time"

type User struct {
    ID          uint      `json:"id"`
    Name        string    `json:"name"`
    Email       string    `json:"email"`
    Password    string    `json:"-"`          // ไม่ส่งกลับใน JSON
    IsVerified  bool      `json:"is_verified"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// สำหรับการ register
type RegisterRequest struct {
    Name     string `json:"name" validate:"required,min=3"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

// สำหรับ login
type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}
```

#### 46.3 Repository Interface

**internal/repository/user_repo.go**
```go
package repository

import (
    "context"
    "your-project/internal/models"
)

type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id uint) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    Update(ctx context.Context, user *models.User) error
    Delete(ctx context.Context, id uint) error
}
```

**Implementation ด้วย GORM (internal/repository/user_repo_gorm.go)**
```go
package repository

import (
    "context"
    "gorm.io/gorm"
    "your-project/internal/models"
)

type userRepositoryGorm struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepositoryGorm{db: db}
}

func (r *userRepositoryGorm) Create(ctx context.Context, user *models.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepositoryGorm) GetByID(ctx context.Context, id uint) (*models.User, error) {
    var user models.User
    err := r.db.WithContext(ctx).First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}
// ... methods อื่นๆ
```

#### 46.4 Usecase

**internal/usecase/user_usecase.go**
```go
package usecase

import (
    "context"
    "errors"
    "golang.org/x/crypto/bcrypt"
    "your-project/internal/models"
    "your-project/internal/repository"
)

type UserUsecase interface {
    Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error)
    GetUserByID(ctx context.Context, id uint) (*models.User, error)
}

type userUsecase struct {
    userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
    return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
    // ตรวจสอบว่ามี email ซ้ำหรือไม่
    existing, _ := u.userRepo.GetByEmail(ctx, req.Email)
    if existing != nil {
        return nil, errors.New("email already registered")
    }
    
    // Hash password
    hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    user := &models.User{
        Name:       req.Name,
        Email:      req.Email,
        Password:   string(hashed),
        IsVerified: false,
    }
    
    if err := u.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    return user, nil
}
```

#### 46.5 Delivery (HTTP)

**internal/delivery/http/handler.go**
```go
package http

import (
    "encoding/json"
    "net/http"
    "your-project/internal/usecase"
    "your-project/internal/models"
    "github.com/go-chi/chi/v5"
    "github.com/go-playground/validator/v10"
)

type UserHandler struct {
    userUsecase usecase.UserUsecase
    validate    *validator.Validate
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
    return &UserHandler{
        userUsecase: userUsecase,
        validate:    validator.New(),
    }
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req models.RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    if err := h.validate.Struct(req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    user, err := h.userUsecase.Register(r.Context(), &req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
```

**internal/delivery/http/routes.go**
```go
package http

import (
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(userHandler *UserHandler, authHandler *AuthHandler) *chi.Mux {
    r := chi.NewRouter()
    
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    
    r.Post("/api/register", userHandler.Register)
    r.Post("/api/login", authHandler.Login)
    
    // Protected routes
    r.Group(func(r chi.Router) {
        r.Use(authHandler.AuthMiddleware)
        r.Get("/api/users/{id}", userHandler.GetUser)
        // ...
    })
    
    return r
}
```

---

### บทที่ 47: Blueprint สำหรับโปรเจกต์ Go ระดับ Production

#### 47.1 โครงสร้างโฟลเดอร์หลัก

```
project-root/
├── cmd/                 # จุดเริ่มต้นของโปรแกรม (executable)
├── internal/            # โค้ดเฉพาะของแอปพลิเคชัน (ไม่ถูก import จากภายนอก)
├── pkg/                 # โค้ดที่นำกลับไปใช้ได้ (reusable) ในโปรเจกต์อื่น
├── api/                 # ไฟล์ที่เกี่ยวข้องกับ API (เช่น Swagger)
├── configs/             # ไฟล์ configuration
├── deploy/              # Docker, Kubernetes files
├── migrations/          # SQL migration files
├── scripts/             # สคริปต์ช่วยพัฒนา
└── test/                # integration tests
```

#### 47.2 ชั้นการทำงานภายใน `internal/`

**apps/** – จุดรวม dependencies และกำหนด routing

- `bootstrap/injection/` – ใช้สำหรับ wire dependencies (repository, service, handler) ด้วย dependency injection แบบ manual
- `router/v1/` – กำหนด routes สำหรับ API version 1 แยก public/protected

**core/** – ชั้นโดเมนของแอปพลิเคชัน (แบ่งตาม domain module)

แต่ละ module (เช่น `auth`, `user`, `iot`) มีโครงสร้างย่อย:

- `entity/` – entity หลัก (struct พร้อม behavior) และ value objects
- `repository/` – interface สำหรับ repository
- `service/` – interface สำหรับ service และการ implement business logic
- `dto/` – data transfer objects (request/response)
- `model/` – model ที่ใช้กับฐานข้อมูล (optional)
- `handler/` – HTTP handlers (รับ request, เรียก service, ส่ง response)
- `routes.go` – ลงทะเบียน routes ของ module นี้

**platform/** – ชั้น infrastructure

- `config/` – โหลด configuration (viper)
- `db/` – การเชื่อมต่อ PostgreSQL, Redis
- `cache/` – ฟังก์ชันช่วยสำหรับ Redis cache
- `queue/` – message queue (Redis pub/sub) + worker pool + dead letter queue
- `logger/` – structured logger (slog) + middleware

**transport/** – ชั้น delivery

- `middleware/` – middleware ต่างๆ (CORS, rate limit, auth, logging, recovery, security headers)
- `httpx/` – utilities สำหรับ response, request, validation
- `utils/` – helper functions (context, pagination)

#### 47.3 การเพิ่มโมดูลใหม่ (Feature)

สมมติต้องการเพิ่มโมดูล `product`:

1. สร้างโครงสร้างใน `internal/core/product/`:
   - `entity/product.go` – entity และ behavior
   - `repository/repository.go` – interface สำหรับ repository
   - `service/service.go` – interface สำหรับ service + implementation
   - `dto/product_dto.go` – request/response structs
   - `handler/product_handler.go` – HTTP handlers
   - `routes.go` – routes สำหรับ product

2. เพิ่ม repository implementation ใน `internal/platform/db/postgres/`
3. เพิ่ม service implementation ใน `internal/core/product/service/service_impl.go`
4. เพิ่ม handler ใน `internal/core/product/handler/`
5. ลงทะเบียน dependencies ใน `internal/apps/app/bootstrap/injection/`
6. เพิ่ม routes ใน `internal/apps/app/router/v1/protected_routes.go`
7. อย่าลืมเพิ่ม migrations ถ้ามีการเปลี่ยนแปลง schema

---

### บทที่ 48: การออกแบบ Workflow และ Task Management

#### 48.1 Workflow การพัฒนา Feature ใหม่

1. **Analyze** – ทำความเข้าใจ requirement, ระบุ domain models, use cases
2. **Design** – ออกแบบ entities, value objects, repository interface, service interface
3. **Implement Domain** – เขียน entity, repository interface, service interface ใน `internal/core/<module>`
4. **Implement Infrastructure** – เขียน repository implementation (GORM), cache, queue (ถ้าจำเป็น)
5. **Implement Service** – เขียน business logic ใน service implementation
6. **Implement Handler** – เขียน HTTP handlers, ตรวจสอบ input validation, ใช้ service
7. **Register Routes** – ลงทะเบียน routes ใน router
8. **Test** – เขียน unit tests (domain, service), integration tests (handler)
9. **Document** – อัปเดต API docs (Swagger) ถ้ามี

#### 48.2 Task List Template

**Phase 1: Domain Design**
- [ ] ระบุ domain model (entity, value objects)
- [ ] กำหนด invariants (business rules)
- [ ] ระบุ use cases (methods in service)
- [ ] กำหนด events (ถ้ามี)
- [ ] ออกแบบ repository interface (methods)
- [ ] ออกแบบ DTOs (request/response)

**Phase 2: Implementation**
- [ ] สร้าง entity struct และ behavior methods
- [ ] สร้าง repository interface
- [ ] สร้าง service interface
- [ ] สร้าง DTO structs
- [ ] เขียน unit tests สำหรับ entity
- [ ] เขียน unit tests สำหรับ service (mock repository)

**Phase 3: Infrastructure**
- [ ] สร้าง repository implementation (GORM)
- [ ] สร้าง migration file (ถ้ามี)
- [ ] ตั้งค่า Redis cache (ถ้าจำเป็น)
- [ ] ตั้งค่า message queue (ถ้าจำเป็น)
- [ ] ทดสอบ repository ด้วย integration test

**Phase 4: Delivery**
- [ ] สร้าง HTTP handlers
- [ ] เพิ่ม input validation (go-playground/validator)
- [ ] สร้าง routes
- [ ] ลงทะเบียน dependencies ใน injection
- [ ] ทดสอบ handler ด้วย httptest

**Phase 5: Integration & Documentation**
- [ ] ทดสอบ end-to-end ด้วย curl/Postman
- [ ] อัปเดต Swagger docs (ถ้ามี)
- [ ] อัปเดต README (ถ้าจำเป็น)
- [ ] รัน linter (`golangci-lint run`) และแก้ไข warnings

**Phase 6: Review & Deploy**
- [ ] Code review
- [ ] ตรวจสอบ performance (ถ้ามีการ query มาก)
- [ ] รัน test coverage (`go test -cover`), ควร > 80%
- [ ] รัน race detector (`go test -race`)
- [ ] Deploy to staging
- [ ] ทดสอบใน staging
- [ ] Deploy to production

#### 48.3 Checklist Template

**Code Quality Checklist**
- [ ] All exported functions have comments
- [ ] No unused imports or variables (go vet)
- [ ] Code formatted with go fmt
- [ ] Error handling is explicit (no ignored errors)
- [ ] No use of panic in library code (only in main/init)
- [ ] Context is passed as first parameter
- [ ] Interfaces are small and focused
- [ ] No global state except configuration

**Security Checklist**
- [ ] Passwords hashed with bcrypt
- [ ] JWT secret loaded from environment, not hardcoded
- [ ] Refresh tokens stored in Redis, not in DB
- [ ] Access token short-lived (≤15min)
- [ ] CORS configured properly (allow only trusted origins)
- [ ] Input validation on all endpoints
- [ ] SQL injection prevented by using parameterized queries (GORM)
- [ ] No sensitive data in logs
- [ ] HTTPS enforced in production
- [ ] Rate limiting on auth endpoints

**Performance Checklist**
- [ ] Database indexes created on frequently queried columns
- [ ] User data cached in Redis
- [ ] Connection pools configured for DB and Redis
- [ ] Use of goroutines for non-blocking tasks (email sending)
- [ ] Avoid N+1 queries (use Preload in GORM)
- [ ] Benchmarks for critical paths

**Testing Checklist**
- [ ] Unit tests cover business logic (usecase)
- [ ] Repository tests with testcontainers or in-memory DB
- [ ] HTTP handler tests with httptest
- [ ] Mock external dependencies (Redis, Mailer)
- [ ] Race condition tests with `-race` flag
- [ ] Test coverage >80%

**Deployment Checklist**
- [ ] Configurable via environment variables
- [ ] Graceful shutdown (wait for existing requests)
- [ ] Health check endpoint
- [ ] Logging to stdout (for container)
- [ ] Docker image built with non-root user
- [ ] Secrets not baked into image
- [ ] Database migration runs automatically on startup (or separate step)
- [ ] Readiness and liveness probes configured
- [ ] Monitoring (Prometheus metrics) exposed

---

## ภาคที่ 8: Domain-Driven Design (DDD) กับ Go

### บทที่ 49: หลักการ DDD และการนำไปใช้ใน Go

#### 49.1 Domain-Driven Design (DDD) คืออะไร?

Domain-Driven Design เป็นแนวทางการออกแบบซอฟต์แวร์ที่เน้นการสร้างโมเดลที่สะท้อนความรู้ความเข้าใจ (domain knowledge) อย่างแท้จริง โดยมีหลักการสำคัญคือการทำให้ซอฟต์แวร์สอดคล้องกับความต้องการทางธุรกิจผ่านการร่วมมือกันระหว่างนักพัฒนาและผู้เชี่ยวชาญในโดเมน (domain experts)

#### 49.2 หลักการสำคัญของ DDD

1. **Ubiquitous Language (ภาษาร่วม)**
   - สร้างภาษากลางที่ใช้ร่วมกันระหว่างนักพัฒนาและผู้เชี่ยวชาญโดเมน
   - ใช้ศัพท์เดียวกันในโค้ด, การสนทนา, และเอกสาร
   - ลดความเข้าใจผิดและเพิ่มความสอดคล้อง

2. **Bounded Context (บริบทที่จำกัด)**
   - แบ่งโดเมนขนาดใหญ่ออกเป็นบริทย่อยที่มีขอบเขตชัดเจน
   - แต่ละ Bounded Context มีโมเดลของตัวเองและภาษาร่วมของตัวเอง
   - ลดความซับซ้อนและความขัดแย้งของโมเดล

3. **Entities และ Value Objects**
   - **Entity**: วัตถุที่มีเอกลักษณ์ (identity) และสามารถเปลี่ยนแปลงได้ (mutable) เช่น `User`, `Order`
   - **Value Object**: วัตถุที่ไม่มีเอกลักษณ์ในตัวเอง ถูกกำหนดโดยคุณสมบัติ (immutable) เช่น `Address`, `Money`

4. **Aggregates**
   - กลุ่มของ Entities และ Value Objects ที่ถูกจัดการเป็นหน่วยเดียวกัน
   - มี Aggregate Root (รูท) เป็นตัวควบคุมความสอดคล้องของข้อมูล
   - เช่น `Order` (aggregate root) ประกอบด้วย `OrderItem` (entity) และ `Address` (value object)

5. **Domain Events**
   - เหตุการณ์สำคัญในโดเมนที่เกิดขึ้น เช่น `OrderPlaced`, `UserRegistered`
   - ใช้ในการสื่อสารระหว่าง aggregates หรือระหว่าง bounded contexts

6. **Repositories**
   - ให้ abstraction ในการเข้าถึง aggregate roots
   - ซ่อนรายละเอียดของแหล่งข้อมูล (database, cache)

7. **Domain Services**
   - ใช้สำหรับ logic ที่ไม่เหมาะจะอยู่ใน entity หรือ value object
   - เช่น `TransferService` ที่โอนเงินระหว่างบัญชี

#### 49.3 การนำ DDD ไปใช้กับ Go: ข้อแนะนำเฉพาะภาษา

1. **จัดโครงสร้างโปรเจกต์ตามโมดูล**
   ```
   /cmd
     /myapp
       main.go
   /internal
     /domain
       /order
         order.go (entity, value objects)
         repository.go (interface)
         events.go
     /application
       order_service.go
     /infrastructure
       /persistence
         order_repo_mysql.go
       /bus
         event_bus_kafka.go
     /presentation
       /http
         order_handler.go
   /pkg
     /shared
       errors.go, uuid.go, etc.
   ```

2. **ใช้ interface เพื่อ Dependency Inversion**
   - Application layer รับ domain interface
   - Infrastructure ถูก inject ผ่าน constructor
   - ใช้ `wire` (Google Wire) หรือ manual DI สำหรับการประกอบ dependencies

3. **จัดการ Transaction**
   - นิยมใช้ `Unit of Work` pattern: application service เริ่ม transaction ผ่าน interface

4. **Domain Events**
   - ใช้ channel หรือ event bus ภายใน memory ก่อน แล้วค่อยเพิ่ม infrastructure

5. **Value Objects กับ immutability**
   - ใช้ struct พร้อม private fields และ constructor functions
   - เปรียบเทียบด้วย `==` หรือ implement `Equals` method

---

### บทที่ 50: Aggregates, Event Storming และ CQRS

#### 50.1 Aggregate (กลุ่มวัตถุที่มีความสอดคล้อง)

**Aggregate** คือกลุ่มของวัตถุ (Entities + Value Objects) ที่ถูกยึดเข้าด้วยกันโดย **Aggregate Root** (รูท) ซึ่งเป็นตัวเดียวที่อนุญาตให้เข้าถึงหรือแก้ไขวัตถุอื่นภายในกลุ่มจากภายนอก การออกแบบ Aggregate ช่วยรักษาความถูกต้องของข้อมูล (invariants) และลดความซับซ้อนในการจัดการธุรกรรม

**หลักการสำคัญ:**
- **Aggregate Root** มี identity และเป็นจุดเดียวที่ภายนอกเข้าถึงได้
- การเปลี่ยนแปลงใด ๆ ภายใน Aggregate ต้องทำผ่าน Root เท่านั้น
- ภายใน Aggregate เดียวกันต้องมีความสอดคล้องกันในทันที (consistency boundary)
- ระหว่าง Aggregates ควรใช้ **eventual consistency** ผ่าน Domain Events

**ตัวอย่างใน Go:**

```go
// aggregate root: Order
type Order struct {
    id      OrderID
    status  OrderStatus
    items   []OrderItem   // value object
    total   Money
}

func (o *Order) AddItem(product Product, quantity int) error {
    if o.status != Draft {
        return errors.New("cannot add item to non-draft order")
    }
    // invariant: total must not exceed limit
    newTotal := o.total.Add(product.Price.Mul(quantity))
    if newTotal.GreaterThan(MaxOrderTotal) {
        return errors.New("order total exceeds limit")
    }
    o.items = append(o.items, NewOrderItem(product, quantity))
    o.total = newTotal
    o.addDomainEvent(OrderItemAdded{OrderID: o.id, ProductID: product.ID})
    return nil
}

// repository interface รับเฉพาะ aggregate root
type OrderRepository interface {
    Save(order *Order) error
    FindByID(id OrderID) (*Order, error)
}
```

#### 50.2 Event Storming (เทคนิคค้นพบโดเมน)

**Event Storming** เป็นเวิร์กช็อปแบบมีส่วนร่วมที่ช่วยให้ทีม (นักพัฒนา, นักธุรกิจ, ผู้เชี่ยวชาญโดเมน) ระบุและเข้าใจโดเมนผ่านเหตุการณ์สำคัญที่เกิดขึ้นในระบบ โดยใช้โน้ตสีต่าง ๆ บนกระดาน

**สัญลักษณ์ทั่วไป:**
- **สีส้ม** – Domain Events (สิ่งที่เกิดขึ้นแล้ว) เช่น `OrderPlaced`, `PaymentReceived`
- **สีน้ำเงิน** – Commands (คำสั่งที่ทำให้เกิดเหตุการณ์) เช่น `PlaceOrder`, `RefundPayment`
- **สีเหลือง** – Aggregates (กลุ่มข้อมูลที่ถูกคำสั่งเรียก) เช่น `Order`, `Customer`
- **สีม่วง** – External Systems / Policies (ระบบภายนอกหรือกฎ)
- **สีเขียว** – Read Models / Views (ข้อมูลที่แสดงผล)

**ขั้นตอนคร่าว ๆ:**
1. ระบุ Domain Events (อดีต) โดยเรียงตามลำดับเวลา
2. ระบุ Commands ที่ทำให้เกิดเหตุการณ์นั้น
3. จับคู่ Command กับ Aggregate (ผู้รับผิดชอบ)
4. เพิ่ม Policies / Rules และ External Systems
5. ระบุ Read Models ที่จำเป็นต่อการแสดงผล

Event Storming นำไปสู่การกำหนด **Bounded Context** และ **Aggregates** ที่ชัดเจน ก่อนเริ่มเขียนโค้ด

#### 50.3 CQRS (Command Query Responsibility Segregation)

CQRS แยกโมเดลการ **เขียน** (Command) และ **อ่าน** (Query) ออกจากกัน ทำให้สามารถปรับแต่งให้เหมาะสมกับแต่ละฝั่งได้ เช่น ใช้ฐานข้อมูลแบบ Event Sourcing สำหรับ Command และฐานข้อมูลแบบ Materialized View สำหรับ Query

**ใน Go สามารถจัดโครงสร้างได้ดังนี้:**

##### แยก Command และ Query Models
```go
// Command models (เขียน)
type PlaceOrderCommand struct {
    OrderID string
    Items   []OrderItemDTO
}

// Query models (อ่าน)
type OrderView struct {
    OrderID    string
    Status     string
    TotalPrice float64
    Items      []OrderItemView
}
```

##### Command Handlers
```go
type PlaceOrderHandler struct {
    repo      domain.OrderRepository
    eventBus  domain.EventBus
}

func (h *PlaceOrderHandler) Handle(ctx context.Context, cmd PlaceOrderCommand) error {
    order, err := domain.NewOrder(cmd.OrderID)
    if err != nil {
        return err
    }
    for _, item := range cmd.Items {
        order.AddItem(item.ProductID, item.Quantity)
    }
    if err := h.repo.Save(order); err != nil {
        return err
    }
    // Publish events for projection
    for _, event := range order.Events() {
        h.eventBus.Publish(event)
    }
    return nil
}
```

##### Query Handlers (อ่านจาก Read Database)
```go
type OrderQueryHandler struct {
    db *sql.DB // หรือ ORM
}

func (h *OrderQueryHandler) GetOrder(ctx context.Context, orderID string) (*OrderView, error) {
    var view OrderView
    err := h.db.QueryRowContext(ctx, "SELECT ... FROM order_views WHERE id = ?", orderID).Scan(&view)
    return &view, err
}
```

##### Projection (สร้าง Read Model จาก Events)
```go
type OrderProjection struct {
    db *sql.DB
}

func (p *OrderProjection) HandleEvent(event domain.DomainEvent) {
    switch e := event.(type) {
    case OrderPlaced:
        p.db.Exec("INSERT INTO order_views (id, status, total) VALUES (?, ?, ?)", e.OrderID, "Placed", e.Total)
    case OrderItemAdded:
        p.db.Exec("INSERT INTO order_items_view ...")
    }
}
```

#### 50.4 การนำ CQRS ไปใช้ใน Go อย่างมีประสิทธิภาพ

- **ใช้ Interface ในการแยก**: Command handlers, Query handlers, Projection handlers แต่ละตัวเป็น struct ที่ implements interface ต่างกัน ทำให้ test และ replace ได้ง่าย
- **Event Store**: ใน Go สามารถใช้ database เช่น PostgreSQL พร้อม JSONB เก็บ events
- **Read Database**: อาจใช้ฐานข้อมูลแยก (SQL, NoSQL) หรือ cache เช่น Redis สำหรับการ query ที่รวดเร็ว
- **Concurrency**: Goroutine + channel ใช้ในการประมวลผล projection แบบ asynchronous

**ข้อควรระวัง:**
- CQRS เพิ่มความซับซ้อน เหมาะสำหรับระบบที่ต้องการความยืดหยุ่นสูง (microservices, high scalability)
- ไม่จำเป็นต้องใช้ Event Sourcing เสมอไป CQRS สามารถแยกโมเดลอ่าน-เขียนโดยใช้ฐานข้อมูลเดียวกันได้ (แต่แยกตารางหรือ schema)
- การจัดการ eventual consistency ต้องออกแบบ UX ให้เหมาะสม

---

### บทที่ 51: การออกแบบบริการด้วย Go-DDD

#### 51.1 โครงสร้างโปรเจกต์แบบ DDD เต็มรูปแบบ

```
project/
├── cmd/
│   └── api/
│       └── main.go                 # entry point
├── internal/
│   ├── domain/                     # ชั้นโดเมน (core business logic)
│   │   ├── user/
│   │   │   ├── entity.go           # User entity
│   │   │   ├── value_objects.go    # Email, Password, etc.
│   │   │   ├── repository.go       # interface
│   │   │   ├── service.go          # domain services
│   │   │   └── events.go           # domain events
│   │   ├── order/
│   │   │   └── ...
│   │   └── shared/                 # shared value objects
│   ├── application/                # ชั้นแอปพลิเคชัน (use cases)
│   │   ├── user/
│   │   │   ├── register.go         # use case
│   │   │   ├── login.go
│   │   │   └── dto.go              # input/output DTOs
│   │   └── order/
│   │       └── ...
│   ├── infrastructure/             # ชั้นโครงสร้างพื้นฐาน
│   │   ├── persistence/
│   │   │   ├── gorm/
│   │   │   │   ├── user_repo.go    # implementation
│   │   │   │   └── models.go       # GORM models
│   │   │   └── redis/
│   │   ├── mail/
│   │   └── bus/                    # event bus
│   └── interfaces/                 # ชั้นติดต่อกับภายนอก
│       ├── http/
│       │   ├── handlers/
│       │   ├── middleware/
│       │   └── routes.go
│       └── cli/
├── pkg/                            # reusable packages
└── go.mod
```

#### 51.2 การสร้าง Domain Layer

**1. Entities (domain/user/entity.go)**
```go
package user

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    id          uuid.UUID
    email       Email          // value object
    password    Password       // value object
    name        string
    isVerified  bool
    createdAt   time.Time
    updatedAt   time.Time
    events      []DomainEvent
}

// Constructor
func NewUser(email, password, name string) (*User, error) {
    emailVO, err := NewEmail(email)
    if err != nil {
        return nil, err
    }
    passwordVO, err := NewPassword(password)
    if err != nil {
        return nil, err
    }
    
    user := &User{
        id:         uuid.New(),
        email:      *emailVO,
        password:   *passwordVO,
        name:       name,
        isVerified: false,
        createdAt:  time.Now(),
        updatedAt:  time.Now(),
        events:     []DomainEvent{},
    }
    
    user.addDomainEvent(NewUserRegisteredEvent(user.id, user.email.String()))
    
    return user, nil
}

// Getters (exported)
func (u *User) ID() uuid.UUID      { return u.id }
func (u *User) Email() Email       { return u.email }
func (u *User) Name() string       { return u.name }
func (u *User) IsVerified() bool   { return u.isVerified }
func (u *User) Events() []DomainEvent { return u.events }
func (u *User) ClearEvents()          { u.events = nil }

// Behavior methods
func (u *User) Verify() {
    u.isVerified = true
    u.updatedAt = time.Now()
    u.addDomainEvent(NewUserVerifiedEvent(u.id))
}

func (u *User) ChangePassword(old, new string) error {
    if err := u.password.Compare(old); err != nil {
        return ErrInvalidPassword
    }
    newPassword, err := NewPassword(new)
    if err != nil {
        return err
    }
    u.password = *newPassword
    u.updatedAt = time.Now()
    return nil
}

func (u *User) addDomainEvent(event DomainEvent) {
    u.events = append(u.events, event)
}
```

**2. Value Objects (domain/user/value_objects.go)**
```go
package user

import (
    "regexp"
    "errors"
    "golang.org/x/crypto/bcrypt"
)

type Email struct {
    value string
}

func NewEmail(email string) (*Email, error) {
    // validate email format
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
    if !emailRegex.MatchString(email) {
        return nil, errors.New("invalid email format")
    }
    return &Email{value: email}, nil
}

func (e Email) String() string { return e.value }

type Password struct {
    hash string
}

func NewPassword(plain string) (*Password, error) {
    if len(plain) < 8 {
        return nil, errors.New("password must be at least 8 characters")
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    return &Password{hash: string(hash)}, nil
}

func (p *Password) Compare(plain string) error {
    return bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plain))
}
```

**3. Domain Events (domain/user/events.go)**
```go
package user

import (
    "time"
    "github.com/google/uuid"
)

type DomainEvent interface {
    OccurredAt() time.Time
}

type UserRegisteredEvent struct {
    UserID    uuid.UUID `json:"user_id"`
    Email     string    `json:"email"`
    Occurred  time.Time `json:"occurred_at"`
}

func (e UserRegisteredEvent) OccurredAt() time.Time { return e.Occurred }

func NewUserRegisteredEvent(userID uuid.UUID, email string) UserRegisteredEvent {
    return UserRegisteredEvent{
        UserID:   userID,
        Email:    email,
        Occurred: time.Now(),
    }
}

type UserVerifiedEvent struct {
    UserID   uuid.UUID `json:"user_id"`
    Occurred time.Time `json:"occurred_at"`
}

func NewUserVerifiedEvent(userID uuid.UUID) UserVerifiedEvent {
    return UserVerifiedEvent{
        UserID:   userID,
        Occurred: time.Now(),
    }
}
```

#### 51.3 การสร้าง Application Layer

**Use Case (application/user/register.go)**
```go
package user

import (
    "context"
    "your-project/internal/domain/user"
    "your-project/internal/infrastructure/bus"
)

type RegisterUseCase struct {
    userRepo user.Repository
    eventBus *bus.EventBus
}

type RegisterInput struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required"`
}

type RegisterOutput struct {
    ID    string `json:"id"`
    Email string `json:"email"`
    Name  string `json:"name"`
}

func NewRegisterUseCase(repo user.Repository, eventBus *bus.EventBus) *RegisterUseCase {
    return &RegisterUseCase{
        userRepo: repo,
        eventBus: eventBus,
    }
}

func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
    // 1. ตรวจสอบว่า email ซ้ำหรือไม่
    emailVO, _ := user.NewEmail(input.Email)
    existing, _ := uc.userRepo.FindByEmail(ctx, *emailVO)
    if existing != nil {
        return nil, ErrEmailAlreadyExists
    }
    
    // 2. สร้าง User entity
    newUser, err := user.NewUser(input.Email, input.Password, input.Name)
    if err != nil {
        return nil, err
    }
    
    // 3. บันทึกผ่าน repository
    if err := uc.userRepo.Save(ctx, newUser); err != nil {
        return nil, err
    }
    
    // 4. Dispatch domain events
    for _, event := range newUser.Events() {
        uc.eventBus.Publish(event)
    }
    newUser.ClearEvents()
    
    // 5. ส่ง output
    return &RegisterOutput{
        ID:    newUser.ID().String(),
        Email: newUser.Email().String(),
        Name:  newUser.Name(),
    }, nil
}
```

#### 51.4 การสร้าง Infrastructure Layer

**Event Bus (infrastructure/bus/event_bus.go)**
```go
package bus

import (
    "context"
    "sync"
    "your-project/internal/domain/user"
)

type EventHandler func(context.Context, user.DomainEvent) error

type EventBus struct {
    handlers map[string][]EventHandler
    mu       sync.RWMutex
}

func NewEventBus() *EventBus {
    return &EventBus{
        handlers: make(map[string][]EventHandler),
    }
}

func (b *EventBus) Subscribe(eventName string, handler EventHandler) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.handlers[eventName] = append(b.handlers[eventName], handler)
}

func (b *EventBus) Publish(event user.DomainEvent) {
    b.mu.RLock()
    handlers := b.handlers[eventName(event)]
    b.mu.RUnlock()
    
    for _, h := range handlers {
        go h(context.Background(), event) // async
    }
}

func eventName(event user.DomainEvent) string {
    switch event.(type) {
    case user.UserRegisteredEvent:
        return "UserRegistered"
    case user.UserVerifiedEvent:
        return "UserVerified"
    default:
        return "Unknown"
    }
}
```

#### 51.5 สรุปประโยชน์ของการใช้ DDD ร่วมกับ Go

- **ความชัดเจนของโดเมน**: โค้ดสะท้อนภาษาธุรกิจ (Ubiquitous Language)
- **การแยกหน้าที่**: แต่ละ layer มีความรับผิดชอบชัดเจน ลด coupling
- **ทดสอบง่าย**: domain layer ไม่ขึ้นกับ infrastructure สามารถ unit test ด้วย mock
- **ปรับเปลี่ยน infrastructure ได้**: เปลี่ยนฐานข้อมูลหรือ event bus โดยไม่กระทบโดเมน
- **Go เหมาะสม**: struct, interface, package system ช่วยให้จัดระเบียบตาม bounded context ได้ดี และการทำงาน concurrency ผ่าน goroutine ช่วยให้จัดการ domain events ได้มีประสิทธิภาพ

---

## ภาคที่ 9: การผสานระบบภายนอกและคุณลักษณะเสริม

### บทที่ 52: Redis สำหรับ Cache และ Message Queue

#### 52.1 บทนำ

Redis (Remote Dictionary Server) เป็นฐานข้อมูลแบบ in-memory ที่รวดเร็วสูง นิยมใช้สำหรับ caching, session storage, message queue, และ real-time analytics ในบทนี้เราจะสำรวจการใช้งาน Redis กับ Go ผ่านไลบรารี `go-redis/redis` ในรูปแบบ cache และ message queue พร้อมตัวอย่างการนำไปใช้จริง

#### 52.2 การติดตั้ง Redis และ go-redis

**ติดตั้ง Redis** (ผ่าน Docker เพื่อความสะดวก)
```bash
docker run -d --name redis -p 6379:6379 redis:alpine
```

**ติดตั้ง go-redis**
```bash
go get github.com/redis/go-redis/v9
```

#### 52.3 การเชื่อมต่อและตั้งค่าเบื้องต้น

```go
package main

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
)

func main() {
    // สร้าง Redis client
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    ctx := context.Background()

    // ทดสอบการเชื่อมต่อ
    pong, err := rdb.Ping(ctx).Result()
    if err != nil {
        fmt.Println("Error connecting to Redis:", err)
        return
    }
    fmt.Println("Connected to Redis:", pong)
}
```

#### 52.4 Redis Cache – การใช้งานพื้นฐาน

**Set/Get**
```go
// เก็บค่า
err := rdb.Set(ctx, "key", "value", 10*time.Minute).Err()

// ดึงค่า
val, err := rdb.Get(ctx, "key").Result()
if err == redis.Nil {
    fmt.Println("key does not exist")
} else if err != nil {
    panic(err)
} else {
    fmt.Println("key:", val)
}
```

**Hash – สำหรับเก็บ object**
```go
// เก็บ hash
err := rdb.HSet(ctx, "user:123", map[string]interface{}{
    "name":  "John",
    "email": "john@example.com",
}).Err()

// ดึง hash
user, err := rdb.HGetAll(ctx, "user:123").Result()
```

**Set – สำหรับ unique items**
```go
// เพิ่มสมาชิก
rdb.SAdd(ctx, "tags", "golang", "redis", "cache")

// ดึงสมาชิกทั้งหมด
members, err := rdb.SMembers(ctx, "tags").Result()
```

**Expiration และ TTL**
```go
// ตั้งเวลา expire
rdb.Expire(ctx, "key", 1*time.Hour)

// ดู TTL ที่เหลือ
ttl, err := rdb.TTL(ctx, "key").Result()
```

#### 52.5 Redis Cache – การใช้งานระดับ Production

**Cache-aside pattern (Lazy loading)**
```go
func GetUser(ctx context.Context, rdb *redis.Client, db *gorm.DB, id uint) (*User, error) {
    key := fmt.Sprintf("user:%d", id)
    
    // 1. พยายามอ่านจาก cache
    val, err := rdb.Get(ctx, key).Result()
    if err == nil {
        var user User
        if err := json.Unmarshal([]byte(val), &user); err == nil {
            return &user, nil
        }
    }
    
    // 2. Cache miss: query database
    var user User
    if err := db.First(&user, id).Error; err != nil {
        return nil, err
    }
    
    // 3. เก็บลง cache (async)
    go func() {
        data, _ := json.Marshal(user)
        rdb.Set(context.Background(), key, data, 10*time.Minute)
    }()
    
    return &user, nil
}
```

**Cache invalidation**
```go
func UpdateUser(ctx context.Context, rdb *redis.Client, db *gorm.DB, user *User) error {
    // 1. อัปเดตฐานข้อมูล
    if err := db.Save(user).Error; err != nil {
        return err
    }
    
    // 2. ลบ cache
    key := fmt.Sprintf("user:%d", user.ID)
    if err := rdb.Del(ctx, key).Err(); err != nil {
        // log error but not fail
    }
    
    return nil
}
```

**Cache warming (preload cache)**
```go
func WarmupUserCache(ctx context.Context, rdb *redis.Client, db *gorm.DB) {
    var users []User
    db.Find(&users)
    for _, u := range users {
        key := fmt.Sprintf("user:%d", u.ID)
        data, _ := json.Marshal(u)
        rdb.Set(ctx, key, data, 10*time.Minute)
    }
}
```

#### 52.6 Redis Message Queue (Pub/Sub)

Redis มีระบบ publish/subscribe ที่ใช้สำหรับส่งข้อความระหว่าง services แบบ asynchronous

**Publisher**
```go
// ส่งข้อความไปยัง channel
err := rdb.Publish(ctx, "order_events", `{"event":"order_created","order_id":123}`).Err()
```

**Subscriber**
```go
// สมัครรับข้อความ
pubsub := rdb.Subscribe(ctx, "order_events")
defer pubsub.Close()

// รอรับข้อความ
ch := pubsub.Channel()
for msg := range ch {
    fmt.Println("Received:", msg.Payload)
    // ประมวลผลข้อความ
}
```

#### 52.7 Redis Queue with List (Producer-Consumer)

**Producer (LPush)**
```go
// ส่งงานเข้ารายการ
job := map[string]interface{}{"task": "send_email", "user_id": 123}
data, _ := json.Marshal(job)
rdb.LPush(ctx, "email_queue", data)
```

**Consumer (BRPop) – blocking pop**
```go
// รับงานจาก queue (block จนกว่าจะมีงาน)
result, err := rdb.BRPop(ctx, 0, "email_queue").Result()
if err != nil {
    // handle error
}
data := result[1] // ค่าที่ได้เป็น string

var job EmailJob
json.Unmarshal([]byte(data), &job)
// ประมวลผล job
```

#### 52.8 Worker Pool with Redis Queue

```go
type Worker struct {
    id      int
    rdb     *redis.Client
    queue   string
    handler func(string) error
    stop    chan struct{}
}

func (w *Worker) Start(ctx context.Context) {
    for {
        select {
        case <-w.stop:
            return
        default:
            // Blocking pop
            result, err := w.rdb.BRPop(ctx, 0, w.queue).Result()
            if err != nil {
                continue
            }
            data := result[1]
            if err := w.handler(data); err != nil {
                // ถ้าผิดพลาด อาจ push กลับไป dead letter queue
                w.rdb.LPush(ctx, "dead_letter_queue", data)
            }
        }
    }
}

func (w *Worker) Stop() {
    close(w.stop)
}
```

#### 52.9 Dead Letter Queue

Dead Letter Queue (DLQ) ใช้สำหรับเก็บข้อความที่ไม่สามารถประมวลผลได้หลังจากลองหลายครั้ง

```go
type JobWithRetry struct {
    Data      string
    RetryCount int
}

func processWithRetry(rdb *redis.Client, data string, maxRetries int) error {
    err := process(data) // ลองประมวลผล
    if err == nil {
        return nil
    }
    
    // ถ้า error ให้ดู retry count
    var job JobWithRetry
    json.Unmarshal([]byte(data), &job)
    if job.RetryCount < maxRetries {
        job.RetryCount++
        newData, _ := json.Marshal(job)
        // ส่งกลับไปยัง queue เพื่อลองใหม่
        rdb.LPush(context.Background(), "main_queue", newData)
    } else {
        // ส่งไป DLQ
        rdb.LPush(context.Background(), "dead_letter_queue", data)
        // แจ้งเตือนผู้ดูแล
    }
    return err
}
```

---

### บทที่ 53: RabbitMQ – Message Broker มาตรฐานองค์กร

#### 53.1 บทนำ

RabbitMQ เป็น message broker ที่ใช้โปรโตคอล AMQP (Advanced Message Queuing Protocol) ถูกออกแบบมาเพื่อการสื่อสารระหว่าง services ในระบบขนาดใหญ่ มีความน่าเชื่อถือสูง รองรับ routing ที่ซับซ้อน และสามารถจัดการกับปริมาณข้อความสูงได้ดี ในบทนี้เราจะแนะนำการใช้งาน RabbitMQ ร่วมกับ Go ผ่านไลบรารี `rabbitmq/amqp091-go`

#### 53.2 การติดตั้ง RabbitMQ และไลบรารี

**ติดตั้ง RabbitMQ ผ่าน Docker**
```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management
```
- Port 5672: สำหรับ AMQP protocol
- Port 15672: สำหรับ management UI (http://localhost:15672, user: guest, pass: guest)

**ติดตั้ง Go client**
```bash
go get github.com/rabbitmq/amqp091-go
```

#### 53.3 การเชื่อมต่อเบื้องต้น

```go
package main

import (
    "log"
    "github.com/rabbitmq/amqp091-go"
)

func main() {
    // เชื่อมต่อ RabbitMQ
    conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatal("Failed to connect to RabbitMQ:", err)
    }
    defer conn.Close()

    // สร้าง channel
    ch, err := conn.Channel()
    if err != nil {
        log.Fatal("Failed to open channel:", err)
    }
    defer ch.Close()

    // ประกาศ queue
    q, err := ch.QueueDeclare(
        "hello", // name
        false,   // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )
    if err != nil {
        log.Fatal("Failed to declare queue:", err)
    }

    log.Println("Connected and queue declared")
}
```

#### 53.4 Producer – ส่งข้อความ

```go
// ส่งข้อความไปยัง queue
body := "Hello World!"
err = ch.Publish(
    "",     // exchange
    q.Name, // routing key (queue name)
    false,  // mandatory
    false,  // immediate
    amqp091.Publishing{
        ContentType: "text/plain",
        Body:        []byte(body),
    })
if err != nil {
    log.Fatal("Failed to publish message:", err)
}
log.Printf("Sent %s", body)
```

#### 53.5 Consumer – รับข้อความ

```go
// สร้าง consumer
msgs, err := ch.Consume(
    q.Name, // queue
    "",     // consumer
    true,   // auto-ack
    false,  // exclusive
    false,  // no-local
    false,  // no-wait
    nil,    // args
)
if err != nil {
    log.Fatal("Failed to register consumer:", err)
}

// รับข้อความแบบ continuous
go func() {
    for d := range msgs {
        log.Printf("Received a message: %s", d.Body)
        // ประมวลผลข้อความ
    }
}()

// ปิดโปรแกรมด้วย ctrl+c
select {}
```

#### 53.6 Work Queue – การกระจายงาน

ในรูปแบบ work queue งานจะถูกกระจายไปยัง consumer หลายตัวแบบ round-robin

**Producer (ส่งงาน)**
```go
// ประกาศ durable queue (รอด survive restart)
q, err := ch.QueueDeclare(
    "task_queue",
    true,  // durable
    false, // delete when unused
    false, // exclusive
    false, // no-wait
    nil,
)

// ส่งงานโดยตั้งค่า persistent
err = ch.Publish(
    "",
    q.Name,
    false,
    false,
    amqp091.Publishing{
        DeliveryMode: amqp091.Persistent, // ทำให้ message อยู่รอดแม้ RabbitMQ restart
        ContentType:  "text/plain",
        Body:         []byte(task),
    })
```

**Consumer (รับงาน)**
```go
// ตั้งค่า prefetch count เพื่อไม่ให้ consumer รับงานเยอะเกินไป
err = ch.Qos(
    1,     // prefetch count
    0,     // prefetch size
    false, // global
)

msgs, err := ch.Consume(
    q.Name,
    "",
    false, // auto-ack (ต้อง ack เอง)
    false,
    false,
    false,
    nil,
)

for d := range msgs {
    log.Printf("Received a message: %s", d.Body)
    // ทำงาน
    time.Sleep(time.Second) // simulate work
    
    // ack เพื่อบอก RabbitMQ ว่าทำงานสำเร็จ
    d.Ack(false)
}
```

#### 53.7 Publish/Subscribe (Fanout Exchange)

ใช้เมื่อต้องการส่งข้อความไปยังหลาย consumer ผ่าน exchange แบบ fanout

**Producer (ส่งข้อความไปยัง exchange)**
```go
// ประกาศ exchange
err = ch.ExchangeDeclare(
    "logs",   // name
    "fanout", // type
    true,     // durable
    false,    // auto-deleted
    false,    // internal
    false,    // no-wait
    nil,      // arguments
)

// ส่งข้อความไปยัง exchange (ไม่ต้องระบุ routing key)
err = ch.Publish(
    "logs", // exchange
    "",     // routing key
    false,
    false,
    amqp091.Publishing{
        ContentType: "text/plain",
        Body:        []byte("Hello World!"),
    })
```

**Consumer (รับข้อความจาก exchange)**
```go
// ประกาศ exchange (ต้องตรงกัน)
ch.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)

// สร้าง queue ชั่วคราว (ไม่คงทน)
q, err := ch.QueueDeclare(
    "",    // name (auto-generated)
    false, // durable
    false, // delete when unused
    true,  // exclusive (เฉพาะ connection นี้)
    false, // no-wait
    nil,
)

// Bind queue กับ exchange
err = ch.QueueBind(
    q.Name, // queue name
    "",     // routing key (ignored for fanout)
    "logs", // exchange
    false,
    nil,
)

// consume
msgs, _ := ch.Consume(q.Name, "", true, false, false, false, nil)
```

#### 53.8 Routing (Direct Exchange)

ใช้ direct exchange เพื่อส่งข้อความไปยัง queue ที่มี routing key ตรงกัน

**Producer**
```go
// ประกาศ direct exchange
ch.ExchangeDeclare("direct_logs", "direct", true, false, false, false, nil)

severity := "error" // routing key
err = ch.Publish(
    "direct_logs",
    severity,
    false,
    false,
    amqp091.Publishing{
        ContentType: "text/plain",
        Body:        []byte("Error message"),
    })
```

**Consumer (รับเฉพาะ routing key ที่สนใจ)**
```go
// ประกาศ exchange
ch.ExchangeDeclare("direct_logs", "direct", true, false, false, false, nil)

// สร้าง queue ชั่วคราว
q, _ := ch.QueueDeclare("", false, false, true, false, nil)

// Bind queue กับ routing keys ที่ต้องการ
for _, s := range []string{"error", "warning"} {
    ch.QueueBind(q.Name, s, "direct_logs", false, nil)
}

// consume
msgs, _ := ch.Consume(q.Name, "", true, false, false, false, nil)
```

#### 53.9 Topic Exchange – การจับคู่แบบ wildcard

Topic exchange ใช้ routing key ที่มี pattern เช่น `*.error` หรือ `order.*`

**Producer**
```go
// ประกาศ topic exchange
ch.ExchangeDeclare("topic_logs", "topic", true, false, false, false, nil)

// ส่งข้อความพร้อม routing key
ch.Publish("topic_logs", "order.created", false, false, ...)
ch.Publish("topic_logs", "order.paid", false, false, ...)
ch.Publish("topic_logs", "order.shipped", false, false, ...)
```

**Consumer (รับ pattern)**
```go
// Bind queue กับ routing pattern
ch.QueueBind(q.Name, "order.*", "topic_logs", false, nil) // รับ order.created, order.paid, order.shipped
ch.QueueBind(q.Name, "*.error", "topic_logs", false, nil) // รับทุกเหตุการณ์ที่มี .error
```

#### 53.10 RPC (Remote Procedure Call)

RabbitMQ สามารถใช้ทำ RPC โดยส่ง request ไปยัง queue และรอ response กลับผ่าน reply queue

**Client**
```go
// สร้าง reply queue ชั่วคราว
replyQueue, _ := ch.QueueDeclare("", false, false, true, false, nil)

// ตั้งค่า correlationId และ replyTo
corrId := uuid.New().String()
err = ch.Publish(
    "",
    "rpc_queue",
    false,
    false,
    amqp091.Publishing{
        ContentType:   "text/plain",
        CorrelationId: corrId,
        ReplyTo:       replyQueue.Name,
        Body:          []byte("ping"),
    })

// รอรับ response
msgs, _ := ch.Consume(replyQueue.Name, "", true, false, false, false, nil)
for d := range msgs {
    if d.CorrelationId == corrId {
        fmt.Println("Response:", string(d.Body))
        break
    }
}
```

**Server**
```go
// ประกาศ queue สำหรับ RPC
q, _ := ch.QueueDeclare("rpc_queue", false, false, false, false, nil)

msgs, _ := ch.Consume(q.Name, "", false, false, false, false, nil)

for d := range msgs {
    // ประมวลผลคำขอ
    response := "pong"
    
    // ส่ง response กลับไปยัง reply queue
    ch.Publish(
        "",
        d.ReplyTo,
        false,
        false,
        amqp091.Publishing{
            ContentType:   "text/plain",
            CorrelationId: d.CorrelationId,
            Body:          []byte(response),
        })
    d.Ack(false)
}
```

---

### บทที่ 54: MQTT สำหรับ IoT และระบบเรียลไทม์

#### 54.1 บทนำ

MQTT (Message Queuing Telemetry Transport) เป็นโปรโตคอล lightweight สำหรับการสื่อสารแบบ publish/subscribe เหมาะสำหรับ IoT, mobile apps, และระบบที่ต้องการการส่งข้อมูลแบบ real-time ด้วยแบนด์วิดท์ต่ำ ใน Go มีไลบรารียอดนิยมเช่น `eclipse/paho.mqtt.golang`

#### 54.2 การติดตั้งและเชื่อมต่อ

**ติดตั้งไลบรารี**
```bash
go get github.com/eclipse/paho.mqtt.golang
```

**เชื่อมต่อกับ MQTT broker** (เช่น Mosquitto)
```go
package main

import (
    "fmt"
    "time"
    mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
    opts := mqtt.NewClientOptions()
    opts.AddBroker("tcp://localhost:1883")
    opts.SetClientID("go_mqtt_client")
    opts.SetUsername("user") // ถ้าต้องการ authentication
    opts.SetPassword("pass")
    opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
        fmt.Printf("Received: %s from topic: %s\n", msg.Payload(), msg.Topic())
    })

    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }
    defer client.Disconnect(250)

    fmt.Println("Connected to MQTT broker")
    // ... do something
}
```

#### 54.3 Publish และ Subscribe

**Publish**
```go
// publish ข้อความ
token := client.Publish("sensors/temperature", 0, false, "25.6")
token.Wait()
```

**Subscribe**
```go
// subscribe topic
token := client.Subscribe("sensors/#", 0, nil) // wildcard # รับทุก subtopic
token.Wait()
```

**Subscribe ด้วย handler**
```go
// กำหนด handler เฉพาะ topic
handler := func(client mqtt.Client, msg mqtt.Message) {
    fmt.Printf("Temperature: %s\n", msg.Payload())
}
token := client.Subscribe("sensors/temperature", 0, handler)
token.Wait()
```

#### 54.4 QoS (Quality of Service) Levels

- **QoS 0** – At most once (fire and forget)
- **QoS 1** – At least once (รับประกันถึง broker แต่ซ้ำได้)
- **QoS 2** – Exactly once (รับประกันไม่ซ้ำ)

```go
// ตัวอย่างใช้ QoS 1
token := client.Publish("important/topic", 1, false, "important message")
token.Wait()
```

#### 54.5 Retained Messages

Retained message คือข้อความที่ broker จะเก็บไว้และส่งให้ subscriber ใหม่ทันทีเมื่อ subscribe

```go
// publish with retained flag
token := client.Publish("device/status", 0, true, "online")
token.Wait()
```

#### 54.6 Last Will and Testament (LWT)

LWT คือข้อความที่ broker จะส่งเมื่อ client disconnect อย่างไม่ปกติ

```go
opts.SetWill("device/status", "offline", 0, true) // กำหนด LWT
```

#### 54.7 ตัวอย่าง: IoT Sensor Data Collector

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "time"
    mqtt "github.com/eclipse/paho.mqtt.golang"
)

type SensorData struct {
    DeviceID  string    `json:"device_id"`
    Temp      float64   `json:"temperature"`
    Humidity  float64   `json:"humidity"`
    Timestamp time.Time `json:"timestamp"`
}

func main() {
    opts := mqtt.NewClientOptions()
    opts.AddBroker("tcp://localhost:1883")
    opts.SetClientID("data_collector")
    opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
        var data SensorData
        if err := json.Unmarshal(msg.Payload(), &data); err != nil {
            log.Printf("Invalid JSON: %v", err)
            return
        }
        fmt.Printf("Received from %s: temp=%.2f, humidity=%.2f\n", 
            data.DeviceID, data.Temp, data.Humidity)
        // บันทึกลง InfluxDB หรือ database อื่น
    })

    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        log.Fatal(token.Error())
    }

    // Subscribe topic ทั้งหมดของ sensor
    if token := client.Subscribe("sensors/+/data", 0, nil); token.Wait() && token.Error() != nil {
        log.Fatal(token.Error())
    }

    // ปิดโปรแกรมด้วย ctrl+c
    select {}
}
```

---

### บทที่ 55: InfluxDB – Time‑Series Database

#### 55.1 บทนำ

InfluxDB เป็นฐานข้อมูล time-series ที่ออกแบบมาเพื่อจัดเก็บและสืบค้นข้อมูลที่มี timestamp เช่น ข้อมูลเซ็นเซอร์, metrics, logs มีประสิทธิภาพสูงและมี query language เฉพาะ (Flux หรือ InfluxQL) ใน Go มี client library อย่าง `influxdata/influxdb-client-go`

#### 55.2 การติดตั้งและเชื่อมต่อ

**ติดตั้ง InfluxDB (Docker)**
```bash
docker run -d --name influxdb -p 9087:8086 influxdb:latest
```

**ติดตั้ง Go client**
```bash
go get github.com/influxdata/influxdb-client-go/v2
```

**เชื่อมต่อ**
```go
package main

import (
    "fmt"
    "time"
    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
    // สร้าง client
    client := influxdb2.NewClient("http://localhost:9087", "my-token")
    defer client.Close()

    // ใช้ write API
    writeAPI := client.WriteAPI("my-org", "my-bucket")

    // สร้าง point
    p := influxdb2.NewPoint("temperature",
        map[string]string{"device": "sensor1"},
        map[string]interface{}{"value": 25.6},
        time.Now())
    writeAPI.WritePoint(p)

    // flush
    writeAPI.Flush()

    fmt.Println("Data written")
}
```

#### 55.3 การเขียนข้อมูล (Write)

**เขียนข้อมูลทีละหลายจุด**
```go
writeAPI := client.WriteAPI("my-org", "my-bucket")
for i := 0; i < 10; i++ {
    p := influxdb2.NewPoint("cpu",
        map[string]string{"host": "server1"},
        map[string]interface{}{"usage": 50 + i},
        time.Now())
    writeAPI.WritePoint(p)
    time.Sleep(1 * time.Second)
}
writeAPI.Flush()
```

**ใช้แบบ batch ด้วย WriteAPIBlocking**
```go
writeAPI := client.WriteAPIBlocking("my-org", "my-bucket")
p := influxdb2.NewPoint("temperature",
    map[string]string{"device": "sensor2"},
    map[string]interface{}{"value": 30.2},
    time.Now())
err := writeAPI.WritePoint(context.Background(), p)
```

#### 55.4 การอ่านข้อมูล (Query) ด้วย Flux

```go
// สร้าง query API
queryAPI := client.QueryAPI("my-org")

// Flux query
flux := `from(bucket:"my-bucket")
  |> range(start: -1h)
  |> filter(fn: (r) => r._measurement == "temperature")
  |> filter(fn: (r) => r.device == "sensor1")
  |> aggregateWindow(every: 1m, fn: mean)`

// Execute query
result, err := queryAPI.Query(context.Background(), flux)
if err != nil {
    panic(err)
}

// Iterate over results
for result.Next() {
    record := result.Record()
    fmt.Printf("Time: %s, Value: %v\n", record.Time(), record.Value())
}
if result.Err() != nil {
    panic(result.Err())
}
```

#### 55.5 การใช้ InfluxQL (ภาษา SQL-like)

```go
query := `SELECT mean(value) FROM temperature WHERE time > now() - 1h GROUP BY time(5m)`
result, err := queryAPI.Query(context.Background(), query)
// ... iterate
```

#### 55.6 การจัดการ Tags และ Fields

- **Tags**: ใช้สำหรับ indexing (เช่น device ID, location) – ต้องเป็น string
- **Fields**: ใช้สำหรับค่าที่วัด (เช่น temperature, humidity) – 可以是数值

```go
p := influxdb2.NewPoint("sensor_data",
    map[string]string{
        "device_id": "sensor1",
        "location":  "room_101",
    },
    map[string]interface{}{
        "temperature": 25.6,
        "humidity":    60,
        "battery":     4.2,
    },
    time.Now())
```

---

### บทที่ 56: WebSocket และ Socket.IO

#### 56.1 WebSocket พื้นฐาน

WebSocket เป็นโปรโตคอลที่ช่วยให้ client-server สื่อสารแบบ full-duplex เหนือ TCP เหมาะสำหรับ real-time applications เช่น chat, live notifications, gaming ใน Go เราใช้ไลบรารี `gorilla/websocket`

**การติดตั้ง**
```bash
go get github.com/gorilla/websocket
```

**Server ตัวอย่าง**
```go
package main

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade failed:", err)
        return
    }
    defer conn.Close()

    for {
        // อ่านข้อความจาก client
        msgType, msg, err := conn.ReadMessage()
        if err != nil {
            break
        }
        log.Printf("Received: %s", msg)

        // ส่งข้อความกลับ
        err = conn.WriteMessage(msgType, []byte("Echo: "+string(msg)))
        if err != nil {
            break
        }
    }
}

func main() {
    http.HandleFunc("/ws", wsHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

**Client (JavaScript)**
```javascript
const socket = new WebSocket('ws://localhost:8080/ws');
socket.onmessage = (event) => console.log('Received:', event.data);
socket.send('Hello');
```

#### 56.2 WebSocket Hub – กระจายข้อความ

ใช้ pattern hub เพื่อจัดการ connections หลายตัวและ broadcast ข้อความ

```go
type Client struct {
    conn *websocket.Conn
    send chan []byte
    hub  *Hub
}

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
        }
    }
}

func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    for {
        _, msg, err := c.conn.ReadMessage()
        if err != nil {
            break
        }
        c.hub.broadcast <- msg
    }
}

func (c *Client) writePump() {
    defer c.conn.Close()
    for {
        select {
        case msg, ok := <-c.send:
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
                return
            }
        }
    }
}
```

#### 56.3 Socket.IO – การใช้งานที่ง่ายขึ้น

Socket.IO เป็น wrapper เหนือ WebSocket ที่ให้คุณสมบัติเพิ่มเติม เช่น fallback, rooms, auto-reconnect ใน Go มีไลบรารี `googollee/go-socket.io`

**การติดตั้ง**
```bash
go get github.com/googollee/go-socket.io
```

**Server ตัวอย่าง**
```go
package main

import (
    "log"
    socketio "github.com/googollee/go-socket.io"
    "net/http"
)

func main() {
    server := socketio.NewServer(nil)

    server.OnConnect("/", func(s socketio.Conn) error {
        s.SetContext("")
        log.Println("Connected:", s.ID())
        return nil
    })

    server.OnEvent("/", "chat message", func(s socketio.Conn, msg string) {
        log.Println("Message:", msg)
        s.Emit("chat message", "Echo: "+msg)
    })

    server.OnDisconnect("/", func(s socketio.Conn, reason string) {
        log.Println("Disconnected:", reason)
    })

    go server.Serve()
    defer server.Close()

    http.Handle("/socket.io/", server)
    http.Handle("/", http.FileServer(http.Dir("./static")))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

**Client (JavaScript)**
```html
<script src="/socket.io/socket.io.js"></script>
<script>
  const socket = io();
  socket.on('chat message', (msg) => console.log(msg));
  socket.emit('chat message', 'Hello from client');
</script>
```

---

### บทที่ 57: การส่ง SMS และ LINE Notify

#### 57.1 การส่ง SMS ผ่าน Twilio

Twilio เป็นบริการ SMS API ยอดนิยม

**การติดตั้ง**
```bash
go get github.com/twilio/twilio-go
```

**ตัวอย่างการส่ง SMS**
```go
package main

import (
    "fmt"
    "github.com/twilio/twilio-go"
    twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func main() {
    client := twilio.NewRestClientWithParams(twilio.ClientParams{
        Username: "ACxxxxxxxxxxxxxx",
        Password: "yyyyyyyyyyyyyyyy",
    })

    params := &twilioApi.CreateMessageParams{}
    params.SetTo("+66812345678")
    params.SetFrom("+1234567890")
    params.SetBody("Hello from Go!")

    resp, err := client.Api.CreateMessage(params)
    if err != nil {
        fmt.Println("Error sending SMS:", err)
    } else {
        fmt.Println("SMS SID:", *resp.Sid)
    }
}
```

#### 57.2 การส่ง SMS ผ่านบริการไทย (SMS360)

บริการในไทย เช่น SMS360, ThaiBulkSMS ฯลฯ มักใช้ HTTP API

```go
func sendSMS360(phone, message string) error {
    url := "https://api.sms360.com/send"
    data := url.Values{}
    data.Set("api_key", "your-api-key")
    data.Set("phone", phone)
    data.Set("message", message)

    resp, err := http.PostForm(url, data)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to send SMS: %s", resp.Status)
    }
    return nil
}
```

#### 57.3 LINE Notify – การแจ้งเตือนผ่าน LINE

LINE Notify ให้บริการแจ้งเตือนฟรี โดยต้องสร้าง token ที่ https://notify-bot.line.me/

**ส่งข้อความผ่าน LINE Notify API**
```go
func sendLineNotify(message string) error {
    url := "https://notify-api.line.me/api/notify"
    data := url.Values{}
    data.Set("message", message)

    req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
    if err != nil {
        return err
    }
    req.Header.Set("Authorization", "Bearer YOUR_TOKEN")
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("LINE Notify failed: %s", resp.Status)
    }
    return nil
}
```

**ส่งรูปภาพ (image)**
```go
// สามารถส่งภาพโดยใช้ multipart form
func sendLineNotifyWithImage(message, imagePath string) error {
    var b bytes.Buffer
    w := multipart.NewWriter(&b)
    w.WriteField("message", message)
    file, _ := os.Open(imagePath)
    defer file.Close()
    part, _ := w.CreateFormFile("imageFile", imagePath)
    io.Copy(part, file)
    w.Close()

    req, _ := http.NewRequest("POST", "https://notify-api.line.me/api/notify", &b)
    req.Header.Set("Authorization", "Bearer YOUR_TOKEN")
    req.Header.Set("Content-Type", w.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

---

### บทที่ 58: Discord Webhook สำหรับแจ้งเตือน

#### 58.1 Discord Webhook คืออะไร

Discord Webhook เป็นวิธีส่งข้อความอัตโนมัติไปยัง channel ใน Discord โดยไม่ต้องใช้ bot token เหมาะสำหรับการแจ้งเตือนจากระบบ (เช่น deployment status, error logs)

#### 58.2 การสร้าง Webhook ใน Discord

1. เปิด Discord, ไปที่ server ของคุณ
2. คลิกขวาที่ channel -> Edit Channel -> Integrations -> Webhooks -> Create Webhook
3. คัดลอก URL (https://discord.com/api/webhooks/...)

#### 58.3 ส่งข้อความผ่าน Go

**โครงสร้าง payload**
```go
type DiscordWebhook struct {
    Content   string         `json:"content,omitempty"`
    Username  string         `json:"username,omitempty"`
    AvatarURL string         `json:"avatar_url,omitempty"`
    Embeds    []DiscordEmbed `json:"embeds,omitempty"`
}

type DiscordEmbed struct {
    Title       string `json:"title,omitempty"`
    Description string `json:"description,omitempty"`
    Color       int    `json:"color,omitempty"` // hexadecimal color (e.g., 0x00ff00)
    Fields      []struct {
        Name   string `json:"name"`
        Value  string `json:"value"`
        Inline bool   `json:"inline,omitempty"`
    } `json:"fields,omitempty"`
    Timestamp string `json:"timestamp,omitempty"`
}
```

**ฟังก์ชันส่งข้อความ**
```go
func sendDiscordWebhook(webhookURL string, content string) error {
    payload := DiscordWebhook{Content: content}
    data, _ := json.Marshal(payload)

    resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(data))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
        return fmt.Errorf("discord webhook returned %s", resp.Status)
    }
    return nil
}
```

**ส่งแบบ embed (สวยงาม)**
```go
embed := DiscordEmbed{
    Title:       "🚀 Deployment",
    Description: "New version v1.2.3 deployed successfully",
    Color:       0x00ff00,
    Fields: []struct{
        Name   string `json:"name"`
        Value  string `json:"value"`
        Inline bool   `json:"inline,omitempty"`
    }{
        {Name: "Environment", Value: "Production", Inline: true},
        {Name: "Version", Value: "v1.2.3", Inline: true},
        {Name: "Deployed By", Value: "CI/CD", Inline: false},
    },
    Timestamp: time.Now().Format(time.RFC3339),
}
payload := DiscordWebhook{Embeds: []DiscordEmbed{embed}}
```

---

## ภาคที่ 10: เทมเพลต กระบวนการพัฒนา และตัวอย่างโค้ด

### บทที่ 59: ตัวอย่างโค้ดครบวงจร (Full‑stack Example)

*(หมายเหตุ: ตัวอย่างโค้ดเต็มรูปแบบจะถูกแทรกในฉบับสมบูรณ์ แต่เนื่องจากข้อจำกัดของพื้นที่ ขอสรุปเป็นโครงสร้างหลัก)*

- โครงสร้างโปรเจกต์ตาม Clean Architecture
- การลงทะเบียนผู้ใช้ด้วย JWT และ refresh token ใน Redis
- การส่งอีเมลยืนยันด้วย gomail + hermes
- การ query ข้อมูลพร้อม cache ใน Redis
- การใช้ GORM transaction สำหรับการสร้าง order
- การใช้ RabbitMQ สำหรับ async email sending
- การแจ้งเตือนผ่าน Discord Webhook เมื่อเกิด error

### บทที่ 60: Task List Template (แสดงไว้ในบทที่ 48)

### บทที่ 61: Checklist Template (แสดงไว้ในบทที่ 48)

### บทที่ 62: แผนภาพการทำงาน (Workflow Diagram)

แผนภาพ Mermaid ที่แทรกในบทต่างๆ สามารถนำไปใช้ในการอธิบายระบบ

**ตัวอย่าง: Clean Architecture Workflow**
```mermaid
graph TD
    A[HTTP Request] --> B[Handler]
    B --> C[Usecase]
    C --> D[Repository Interface]
    D --> E[Database Implementation]
    C --> F[Domain Events]
    F --> G[Event Bus]
    G --> H[Other Services]
    C --> I[Response]
    I --> B
```

**ตัวอย่าง: JWT Authentication Flow**
```mermaid
sequenceDiagram
    participant Client
    participant API
    participant DB
    participant Redis

    Client->>API: POST /login (email, password)
    API->>DB: Validate credentials
    DB-->>API: User data
    API->>API: Generate access & refresh tokens
    API->>Redis: Store refresh token
    API-->>Client: Tokens
```

### บทที่ 63: mop Config – การจัดการ Configuration

**ไฟล์ config/config.yaml**
```yaml
server:
  port: 8080
  mode: release
  read_timeout: 15s
  write_timeout: 15s

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  name: userdb
  sslmode: disable
  max_open_conns: 25
  max_idle_conns: 25
  conn_max_lifetime: 5m

redis:
  addr: localhost:6379
  password: ""
  db: 0
  pool_size: 10
  ttl: 10m

jwt:
  secret: "your-secret-key-change-in-production"
  access_expiry: 15m
  refresh_expiry: 168h  # 7 days

smtp:
  host: smtp.gmail.com
  port: 587
  username: your-email@gmail.com
  password: your-app-password
  from: your-email@gmail.com

log:
  level: info
  format: json
  output: stdout

rabbitmq:
  url: amqp://guest:guest@localhost:5672/
  exchange: events

mqtt:
  broker: tcp://localhost:1883
  client_id: go_mqtt_client

influxdb:
  url: http://localhost:9087
  token: my-token
  org: my-org
  bucket: my-bucket
```

**การโหลด config ด้วย viper (Go)**
```go
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    JWT      JWTConfig
    SMTP     SMTPConfig
    Log      LogConfig
    RabbitMQ RabbitMQConfig
    MQTT     MQTTConfig
    InfluxDB InfluxDBConfig
}

func LoadConfig() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
```

---
เราจะเพิ่มคำอธิบายใต้แผนภาพภาษาไทยสำหรับแต่ละส่วน เพื่อให้ผู้อ่านเข้าใจการไหลของข้อมูลได้ง่ายขึ้น โดยใช้แผนภาพ Mermaid ที่นำเสนอในคำตอบก่อนหน้านี้ พร้อมคำอธิบายภาษาไทยใต้ภาพ

---

## ภาพรวมการไหลของข้อมูล (Overview Data Flow)

```mermaid
graph TD
    subgraph Part1[Part 1: Fundamentals]
        A1[Source Code] --> A2[Compiler]
        A2 --> A3[Executable]
        A3 --> A4[Run & Debug]
    end

    subgraph Part2[Part 2: Basic Language & Data Structures]
        B1[Variables & Types] --> B2[Control Flow]
        B2 --> B3[Functions]
        B3 --> B4[Structs & Interfaces]
    end

    subgraph Part3[Part 3: Project Management & Advanced Data Structures]
        C1[Go Modules] --> C2[Tests]
        C2 --> C3[Arrays/Slices/Maps]
        C3 --> C4[Error Handling]
    end

    subgraph Part4[Part 4: Practical Application Development]
        D1[HTTP Server] --> D2[JSON/XML]
        D2 --> D3[Concurrency]
        D3 --> D4[Logging & Config]
    end

    subgraph Part5[Part 5: Professional Go Development]
        E1[Benchmarks] --> E2[Profiling]
        E2 --> E3[Context]
        E3 --> E4[Generics]
    end

    subgraph Part6[Part 6: Popular Tools & Libraries]
        F1[chi, viper, cobra, zap] --> F2[GORM]
        F2 --> F3[gomail & hermes]
    end

    subgraph Part7[Part 7: Architecture Design & Workflow]
        G1[Clean Architecture] --> G2[Blueprint]
        G2 --> G3[Workflow & Tasks]
    end

    subgraph Part8[Part 8: Domain-Driven Design]
        H1[DDD Principles] --> H2[Aggregates & Events]
        H2 --> H3[CQRS & Services]
    end

    subgraph Part9[Part 9: External Systems Integration]
        I1[Redis] --> I2[RabbitMQ]
        I2 --> I3[MQTT]
        I3 --> I4[InfluxDB]
        I4 --> I5[WebSocket]
        I5 --> I6[SMS/LINE/Discord]
    end

    subgraph Part10[Part 10: Templates & Examples]
        J1[Full‑stack Example] --> J2[Task List]
        J2 --> J3[Checklist]
        J3 --> J4[Workflow Diagram]
        J4 --> J5[mop Config]
    end

    A4 --> B1
    C4 --> D1
    D4 --> E1
    E4 --> F1
    F3 --> G1
    G3 --> H1
    H3 --> I1
    I6 --> J1
```

**คำอธิบาย:**  
แผนภาพนี้แสดงความสัมพันธ์ของเนื้อหาทั้ง 10 ภาค โดยเริ่มจากภาคที่ 1 (พื้นฐานการเขียนโปรแกรม) ไหลไปสู่ภาคที่ 2 (โครงสร้างข้อมูลพื้นฐาน) จากนั้นต่อเนื่องไปจนถึงภาคที่ 10 (เทมเพลตและตัวอย่างโค้ด) แต่ละภาคส่งต่อข้อมูลและแนวคิดไปยังภาคถัดไป ทำให้เห็นภาพรวมว่าความรู้แต่ละส่วนเชื่อมโยงกันอย่างไร

---

## ภาคที่ 1: ปฐมบทกับการเขียนโปรแกรม (บทที่ 1–5)

**แผนภาพ: จากแนวคิดสู่โปรแกรมแรก**

```mermaid
graph LR
    subgraph User[ผู้ใช้]
        U1[เขียนโค้ด]
    end

    subgraph Terminal[เทอร์มินัล]
        T1[go run] --> T2[คอมไพเลอร์]
        T2 --> T3[ไบนารี]
        T3 --> T4[ผลลัพธ์]
    end

    U1 -->|main.go| T1
    T4 -->|แสดงผล| U1
```

**คำอธิบาย:**  
ผู้ใช้เขียนโค้ด (ไฟล์ `.go`) แล้วสั่ง `go run` ซึ่งจะส่งให้คอมไพเลอร์แปลงเป็นไบนารีและทำงานทันที ผลลัพธ์แสดงบนเทอร์มินัล ผู้ใช้เห็นผลและสามารถปรับปรุงโค้ดต่อไป วงจรนี้แสดงการพัฒนาโปรแกรมแรกด้วย Go

---

## ภาคที่ 2: พื้นฐานภาษาและโครงสร้างข้อมูล (บทที่ 6–16)

**แผนภาพ: การไหลของข้อมูลผ่านโครงสร้างภาษา**

```mermaid
graph TD
    A[การประกาศตัวแปร] --> B[โครงสร้างควบคุม<br/>(if, for, switch)]
    B --> C[การเรียกฟังก์ชัน]
    C --> D[สตรัคและเมธอด]
    D --> E[พอยน์เตอร์]
    E --> F[อินเทอร์เฟซ]

    subgraph Data[ชนิดข้อมูล]
        A1[ชนิดพื้นฐาน<br/>int, string, bool]
        A2[ชนิดประกอบ<br/>array, slice, map]
    end
    A1 --> A
    A2 --> A
```

**คำอธิบาย:**  
ข้อมูลเริ่มจากตัวแปรที่มีชนิดข้อมูลพื้นฐานหรือประกอบ ผ่านโครงสร้างควบคุม (if, for) แล้วถูกส่งไปยังฟังก์ชันเพื่อประมวลผล จากนั้นนำผลลัพธ์มาจัดเก็บในสตรัคและเมธอด สามารถใช้พอยน์เตอร์เพื่อส่งต่ออ้างอิง และใช้ interface ในการสร้างความยืดหยุ่น

---

## ภาคที่ 3: การจัดการโปรเจกต์และโครงสร้างข้อมูลขั้นสูง (บทที่ 17–23)

**แผนภาพ: การจัดการโปรเจกต์ → ทดสอบ → คอลเลกชัน → ข้อผิดพลาด**

```mermaid
graph TD
    subgraph Modules[โมดูล]
        M1[go mod init] --> M2[go.mod]
        M2 --> M3[go get / tidy]
        M3 --> M4[go.sum]
    end

    subgraph Testing[การทดสอบ]
        T1[_test.go] --> T2[go test]
        T2 --> T3[รายงาน coverage]
    end

    subgraph Collections[คอลเลกชัน]
        C1[อาเรย์] --> C2[สไลซ์]
        C2 --> C3[แมพ]
    end

    subgraph Errors[ข้อผิดพลาด]
        E1[errors.New] --> E2[fmt.Errorf]
        E2 --> E3[errors.Is / As]
    end

    M4 --> T1
    T3 --> C1
    C3 --> E1
```

**คำอธิบาย:**  
เริ่มจากสร้างโมดูลด้วย `go mod init` ได้ไฟล์ `go.mod` และ `go.sum` จากนั้นเขียน test (`_test.go`) แล้วรัน `go test` เพื่อวัด coverage ข้อมูลที่ถูกต้องผ่านการทดสอบจะถูกนำไปใช้กับคอลเลกชัน (array, slice, map) และในที่สุดการจัดการข้อผิดพลาดจะเกิดขึ้นเมื่อพบ error

---

## ภาคที่ 4: การพัฒนาแอปพลิเคชันเชิงปฏิบัติ (บทที่ 24–33)

**แผนภาพ: คำขอ HTTP → การประมวลผล → การตอบกลับ**

```mermaid
graph LR
    subgraph Client[ไคลเอนต์]
        C1[เบราว์เซอร์ / API Client]
    end

    subgraph Server[เซิร์ฟเวอร์]
        S1[HTTP Handler] --> S2[แปลง JSON/XML]
        S2 --> S3[ตรรกะธุรกิจ<br/>(ฟังก์ชัน, การทำงานพร้อมกัน)]
        S3 --> S4[แปลง JSON/XML]
        S4 --> S5[ส่งคำตอบ]
    end

    subgraph Storage[ที่เก็บข้อมูล]
        ST1[ไฟล์]
        ST2[ฐานข้อมูล]
    end

    C1 -->|คำขอ| S1
    S3 <--> ST1
    S3 <--> ST2
    S5 -->|คำตอบ| C1
```

**คำอธิบาย:**  
ไคลเอนต์ส่งคำขอ HTTP มาที่ Handler ฝั่งเซิร์ฟเวอร์ Handler แปลงข้อมูล JSON/XML แล้วส่งต่อไปยังชั้นตรรกะธุรกิจ ซึ่งอาจมีการอ่าน/เขียนไฟล์หรือฐานข้อมูล หลังจากประมวลผลเสร็จ ข้อมูลจะถูกแปลงกลับเป็น JSON/XML แล้วส่งกลับไปยังไคลเอนต์

---

## ภาคที่ 5: สู่การเป็นนักพัฒนา Go มืออาชีพ (บทที่ 34–42)

**แผนภาพ: วัดประสิทธิภาพ → วิเคราะห์โปรไฟล์ → ปรับปรุง**

```mermaid
graph TD
    A[เขียนโค้ด] --> B[Benchmark]
    B --> C[Profiling<br/>(CPU, หน่วยความจำ)]
    C --> D[ระบุจุดที่ช้า]
    D --> E[ปรับปรุง<br/>(Generics, context)]
    E --> A

    subgraph Tools[เครื่องมือ]
        T1[go test -bench]
        T2[pprof]
        T3[go vet / golangci-lint]
    end
    B --> T1
    C --> T2
    D --> T3
```

**คำอธิบาย:**  
หลังจากเขียนโค้ดแล้ว ให้นำไป benchmark ด้วย `go test -bench` เพื่อวัดประสิทธิภาพ จากนั้นใช้ pprof วิเคราะห์โปรไฟล์ CPU และหน่วยความจำ เพื่อหาโค้ดที่ช้า (bottleneck) แล้วปรับปรุงด้วยเทคนิคต่างๆ เช่น generics หรือ context หลังจากปรับปรุงแล้วจะวนกลับไปวัดประสิทธิภาพอีกครั้ง

---

## ภาคที่ 6: เครื่องมือและไลบรารียอดนิยม (บทที่ 43–45)

**แผนภาพ: Router → Config → CLI → Logger → ORM → Email**

```mermaid
graph LR
    subgraph Router[เราเตอร์]
        R1[chi] --> R2[มิดเดิลแวร์]
    end

    subgraph Config[การตั้งค่า]
        V1[viper] --> V2[config.yaml]
        V2 --> R1
    end

    subgraph CLI[บรรทัดคำสั่ง]
        C1[cobra] --> C2[คำสั่ง]
        C2 --> V1
    end

    subgraph Logger[การบันทึก]
        Z1[zap] --> Z2[บันทึกแบบมีโครงสร้าง]
    end

    subgraph ORM[ORM]
        G1[GORM] --> G2[การจัดการฐานข้อมูล]
    end

    subgraph Email[อีเมล]
        M1[gomail] --> M2[hermes]
        M2 --> M3[HTML Email]
    end

    R2 --> Z1
    R2 --> G1
    G1 --> M1
```

**คำอธิบาย:**  
การกำหนดค่า (viper) จะถูกโหลดจากไฟล์ `config.yaml` และใช้ใน router (chi) ส่วน cobra ใช้สร้าง CLI ที่อาจเรียกใช้ viper หรือ chi ได้ ระหว่างการทำงาน router จะเรียก logger (zap) และ ORM (GORM) ซึ่ง GORM สามารถเชื่อมต่อกับ gomail เพื่อส่งอีเมลผ่าน hermes ที่สร้าง HTML template

---

## ภาคที่ 7: การออกแบบสถาปัตยกรรมและ Workflow (บทที่ 46–48)

**แผนภาพ: ชั้น Clean Architecture**

```mermaid
graph TD
    subgraph Delivery[ชั้นส่งข้อมูล]
        D1[HTTP Handlers] --> D2[DTOs]
    end

    subgraph Usecase[ชั้นกรณีการใช้งาน]
        U1[ตรรกะธุรกิจ] --> U2[Repository Interface]
    end

    subgraph Repository[ชั้นเก็บข้อมูล]
        R1[Repository Interface] --> R2[Implementation<br/>(GORM, Redis)]
    end

    subgraph External[ระบบภายนอก]
        E1[ฐานข้อมูล] --> E2[แคช]
    end

    D2 --> U1
    U2 --> R1
    R2 --> E1
    R2 --> E2
```

**คำอธิบาย:**  
คำขอ HTTP เข้าสู่ Delivery (handlers) ซึ่งแปลงเป็น DTO แล้วส่งไปยัง Usecase เพื่อประมวลผลตรรกะธุรกิจ Usecase เรียก Repository Interface โดยไม่ต้องรู้รายละเอียดของฐานข้อมูล Repository Implementation จะติดต่อกับฐานข้อมูลหรือแคชจริง ๆ แล้วคืนผลลัพธ์กลับไปตามลำดับ

---

## ภาคที่ 8: Domain-Driven Design (DDD) (บทที่ 49–51)

**แผนภาพ: Domain → Application → Infrastructure**

```mermaid
graph LR
    subgraph Domain[ชั้นโดเมน]
        A[Entities & Value Objects] --> B[Aggregates]
        B --> C[Domain Events]
        C --> D[Repository Interface]
    end

    subgraph Application[ชั้นแอปพลิเคชัน]
        E[Use Cases] --> F[DTOs]
        F --> G[Service Orchestration]
    end

    subgraph Infrastructure[ชั้นโครงสร้างพื้นฐาน]
        H[Repository Implementation] --> I[Event Bus]
        I --> J[Message Broker]
    end

    A --> E
    D --> H
    C --> I
```

**คำอธิบาย:**  
ชั้นโดเมน (Entities, Value Objects, Aggregates) กำหนดพฤติกรรมทางธุรกิจและสร้าง Domain Events ชั้นแอปพลิเคชันใช้ Use Cases ในการจัดลำดับการทำงาน โดยรับ DTO และเรียกใช้ Service Orchestration ชั้นโครงสร้างพื้นฐาน Implement Repository และ Event Bus เพื่อติดต่อกับฐานข้อมูลหรือ Message Broker

---


---

## ภาคที่ 10: เทมเพลต กระบวนการพัฒนา และตัวอย่างโค้ด (บทที่ 59–63)

**แผนภาพ: Example Application → Tasks → Checklist → Diagrams → Config**

```mermaid
graph LR
    subgraph Example[ตัวอย่าง]
        E1[Full‑stack Example] --> E2[โครงสร้างโค้ด]
    end

    subgraph Processes[กระบวนการ]
        P1[Task List] --> P2[Checklist]
        P2 --> P3[Workflow Diagram]
    end

    subgraph Config[การตั้งค่า]
        C1[mop Config] --> C2[YAML/ENV]
        C2 --> E1
    end

    E2 --> P1
    P3 --> C1
```

**คำอธิบาย:**  
จากตัวอย่างโปรเจกต์ครบวงจร (Full‑stack Example) จะได้โครงสร้างโค้ดที่สามารถนำไปใช้เป็นต้นแบบ โครงสร้างนี้ช่วยให้ทีมสร้าง Task List และ Checklist เพื่อติดตามงาน พร้อมทั้ง Workflow Diagram ที่อธิบายขั้นตอนการทำงาน สุดท้าย mop Config จัดการค่า configuration (YAML/ENV) เพื่อให้แอปพลิเคชันทำงานได้ในหลายสภาพแวดล้อม

---
เราจะนำเสนอแผนภาพพร้อมคำอธิบายภาษาไทยสำหรับแต่ละภาค โดยสามารถนำ Mermaid code ไปใช้ใน draw.io หรือโปรแกรมที่รองรับ Mermaid ได้เลย

---

## 📌 ภาพรวมการไหลของข้อมูลในคู่มือ

```mermaid
graph TD
    subgraph Part1[Part 1: Fundamentals]
        A1[Source Code] --> A2[Compiler]
        A2 --> A3[Executable]
        A3 --> A4[Run & Debug]
    end

    subgraph Part2[Part 2: Basic Language & Data Structures]
        B1[Variables & Types] --> B2[Control Flow]
        B2 --> B3[Functions]
        B3 --> B4[Structs & Interfaces]
    end

    subgraph Part3[Part 3: Project Management & Advanced Data Structures]
        C1[Go Modules] --> C2[Tests]
        C2 --> C3[Arrays/Slices/Maps]
        C3 --> C4[Error Handling]
    end

    subgraph Part4[Part 4: Practical Application Development]
        D1[HTTP Server] --> D2[JSON/XML]
        D2 --> D3[Concurrency]
        D3 --> D4[Logging & Config]
    end

    subgraph Part5[Part 5: Professional Go Development]
        E1[Benchmarks] --> E2[Profiling]
        E2 --> E3[Context]
        E3 --> E4[Generics]
    end

    subgraph Part6[Part 6: Popular Tools & Libraries]
        F1[chi, viper, cobra, zap] --> F2[GORM]
        F2 --> F3[gomail & hermes]
    end

    subgraph Part7[Part 7: Architecture Design & Workflow]
        G1[Clean Architecture] --> G2[Blueprint]
        G2 --> G3[Workflow & Tasks]
    end

    subgraph Part8[Part 8: Domain-Driven Design]
        H1[DDD Principles] --> H2[Aggregates & Events]
        H2 --> H3[CQRS & Services]
    end

    subgraph Part9[Part 9: External Systems Integration]
        I1[Redis] --> I2[RabbitMQ]
        I2 --> I3[MQTT]
        I3 --> I4[InfluxDB]
        I4 --> I5[WebSocket]
        I5 --> I6[SMS/LINE/Discord]
    end

    subgraph Part10[Part 10: Templates & Examples]
        J1[Full‑stack Example] --> J2[Task List]
        J2 --> J3[Checklist]
        J3 --> J4[Workflow Diagram]
        J4 --> J5[mop Config]
    end

    A4 --> B1
    C4 --> D1
    D4 --> E1
    E4 --> F1
    F3 --> G1
    G3 --> H1
    H3 --> I1
    I6 --> J1
```

**คำอธิบายภาษาไทย:**  
แผนภาพนี้แสดงภาพรวมการไหลของเนื้อหาในคู่มือ แต่ละภาคจะนำไปสู่ภาคถัดไป โดยเริ่มจากพื้นฐาน (คอมไพล์, รัน) ไปสู่โครงสร้างข้อมูล, การจัดการโปรเจกต์, การพัฒนาเชิงปฏิบัติ, เครื่องมือ, สถาปัตยกรรม, DDD, การเชื่อมต่อระบบภายนอก และสุดท้ายคือเทมเพลตและตัวอย่างโค้ด ซึ่งเป็นการเรียนรู้แบบลำดับและต่อยอดกัน

---

## ภาคที่ 1: ปฐมบทกับการเขียนโปรแกรม (บทที่ 1–5)

### แผนภาพ: จากแนวคิดสู่โปรแกรมแรก

```mermaid
graph LR
    subgraph User[ผู้ใช้]
        U1[เขียนโค้ด]
    end

    subgraph Terminal
        T1[go run] --> T2[คอมไพเลอร์]
        T2 --> T3[ไฟล์ binary]
        T3 --> T4[แสดงผล]
    end

    U1 -->|main.go| T1
    T4 -->|คอนโซล| U1
```

**คำอธิบายภาษาไทย:**  
ผู้ใช้เขียนซอร์สโค้ด (`main.go`) จากนั้นใช้คำสั่ง `go run` ซึ่งจะเรียกคอมไพเลอร์ (compiler) แปลงเป็นไฟล์ binary แล้วรันโปรแกรม ผลลัพธ์แสดงออกทางคอนโซล แผนภาพนี้แสดงกระบวนการพื้นฐานของการพัฒนาโปรแกรมด้วย Go ตั้งแต่การเขียนโค้ดไปจนถึงการรันและเห็นผลลัพธ์

---

## ภาคที่ 2: พื้นฐานภาษาและโครงสร้างข้อมูล (บทที่ 6–16)

### แผนภาพ: การไหลของข้อมูลในภาษา Go

```mermaid
graph TD
    A[ประกาศตัวแปร] --> B[โครงสร้างควบคุม<br/>(if, for, switch)]
    B --> C[เรียกใช้ฟังก์ชัน]
    C --> D[Struct และ Method]
    D --> E[พอยน์เตอร์]
    E --> F[อินเทอร์เฟซ]

    subgraph Data[ชนิดข้อมูล]
        A1[ชนิดพื้นฐาน<br/>int, string, bool]
        A2[ชนิดประกอบ<br/>array, slice, map]
    end
    A1 --> A
    A2 --> A
```

**คำอธิบายภาษาไทย:**  
การเขียนโปรแกรมใน Go เริ่มจากการประกาศตัวแปร (ใช้ชนิดข้อมูลพื้นฐานหรือชนิดประกอบ) จากนั้นใช้โครงสร้างควบคุมเพื่อกำหนดทิศทางการทำงาน เรียกใช้ฟังก์ชัน ซึ่งอาจเป็น method ที่ผูกกับ struct หรือรับพอยน์เตอร์ และสุดท้ายนำไปใช้งานผ่านอินเทอร์เฟซ แผนภาพนี้แสดงลำดับการใช้ส่วนประกอบภาษาเพื่อสร้างโปรแกรมที่ซับซ้อนขึ้น

---

## ภาคที่ 3: การจัดการโปรเจกต์และโครงสร้างข้อมูลขั้นสูง (บทที่ 17–23)

### แผนภาพ: โมดูล → ทดสอบ → คอลเลกชัน → การจัดการข้อผิดพลาด

```mermaid
graph TD
    subgraph Modules[Go Modules]
        M1[go mod init] --> M2[go.mod]
        M2 --> M3[go get / tidy]
        M3 --> M4[go.sum]
    end

    subgraph Testing[การทดสอบ]
        T1[_test.go] --> T2[go test]
        T2 --> T3[รายงาน coverage]
    end

    subgraph Collections[โครงสร้างข้อมูล]
        C1[อาเรย์] --> C2[สไลซ์]
        C2 --> C3[แมพ]
    end

    subgraph Errors[การจัดการข้อผิดพลาด]
        E1[errors.New] --> E2[fmt.Errorf]
        E2 --> E3[errors.Is / As]
    end

    M4 --> T1
    T3 --> C1
    C3 --> E1
```

**คำอธิบายภาษาไทย:**  
เริ่มต้นจัดการโปรเจกต์ด้วย Go Modules สร้างไฟล์ `go.mod` และ `go.sum` จากนั้นเขียนการทดสอบในไฟล์ `_test.go` แล้วรันด้วย `go test` ข้อมูลจากการทดสอบจะนำไปปรับปรุงโครงสร้างข้อมูลอย่างอาเรย์, สไลซ์, แมพ ส่วนการจัดการข้อผิดพลาด (error) จะถูกนำไปใช้ในโค้ดจริง เพื่อให้โปรแกรมมีความน่าเชื่อถือ

---

## ภาคที่ 4: การพัฒนาแอปพลิเคชันเชิงปฏิบัติ (บทที่ 24–33)

### แผนภาพ: HTTP Request → ประมวลผล → Response

```mermaid
graph LR
    subgraph Client[ไคลเอนต์]
        C1[เบราว์เซอร์ / API Client]
    end

    subgraph Server[เซิร์ฟเวอร์]
        S1[HTTP Handler] --> S2[แปลง JSON/XML]
        S2 --> S3[ตรรกะธุรกิจ<br/>(ฟังก์ชัน, การทำงานพร้อมกัน)]
        S3 --> S4[แปลงเป็น JSON/XML]
        S4 --> S5[ส่ง Response]
    end

    subgraph Storage[ที่เก็บข้อมูล]
        ST1[ไฟล์]
        ST2[ฐานข้อมูล]
    end

    C1 -->|Request| S1
    S3 <--> ST1
    S3 <--> ST2
    S5 -->|Response| C1
```

**คำอธิบายภาษาไทย:**  
ไคลเอนต์ส่งคำขอ (HTTP request) มาที่เซิร์ฟเวอร์ ผ่าน Handler ซึ่งจะแปลงข้อมูล JSON/XML จาก body ให้เป็น struct ของ Go จากนั้นตรรกะธุรกิจจะประมวลผล อาจอ่านหรือเขียนข้อมูลจากไฟล์หรือฐานข้อมูล ผลลัพธ์ถูกแปลงกลับเป็น JSON/XML แล้วส่งกลับไปยังไคลเอนต์ แผนภาพนี้แสดงการทำงานพื้นฐานของเว็บแอปพลิเคชันด้วย Go

---

## ภาคที่ 5: สู่การเป็นนักพัฒนา Go มืออาชีพ (บทที่ 34–42)

### แผนภาพ: Benchmark → Profile → Optimize

```mermaid
graph TD
    A[เขียนโค้ด] --> B[เขียน Benchmark]
    B --> C[Profiling<br/>(CPU, memory)]
    C --> D[ระบุจุดช้า]
    D --> E[ปรับปรุงประสิทธิภาพ<br/>(Generics, context)]
    E --> A

    subgraph Tools[เครื่องมือ]
        T1[go test -bench]
        T2[pprof]
        T3[go vet / golangci-lint]
    end
    B --> T1
    C --> T2
    D --> T3
```

**คำอธิบายภาษาไทย:**  
หลังจากเขียนโค้ดแล้ว ให้เขียน benchmark ในไฟล์ `_test.go` รันด้วย `go test -bench` จากนั้นใช้ profiler (pprof) วิเคราะห์การใช้ CPU และหน่วยความจำ เพื่อหาจุดที่ช้า ปรับปรุงโดยใช้ Generics, context หรือเทคนิคอื่น ๆ แล้วนำโค้ดที่ปรับแล้วกลับมาวัดซ้ำ เป็นวงจรการเพิ่มประสิทธิภาพอย่างต่อเนื่อง

---

## ภาคที่ 6: เครื่องมือและไลบรารียอดนิยม (บทที่ 43–45)

### แผนภาพ: Router → Config → CLI → Logger → ORM → Email

```mermaid
graph LR
    subgraph Router
        R1[chi] --> R2[Middleware]
    end

    subgraph Config
        V1[viper] --> V2[config.yaml]
        V2 --> R1
    end

    subgraph CLI
        C1[cobra] --> C2[Commands]
        C2 --> V1
    end

    subgraph Logger
        Z1[zap] --> Z2[Structured Logs]
    end

    subgraph ORM
        G1[GORM] --> G2[Database Operations]
    end

    subgraph Email
        M1[gomail] --> M2[hermes]
        M2 --> M3[HTML Email]
    end

    R2 --> Z1
    R2 --> G1
    G1 --> M1
```

**คำอธิบายภาษาไทย:**  
การพัฒนาแอปพลิเคชันยุคใหม่มักใช้เครื่องมือหลายตัวร่วมกัน เช่น chi สำหรับ routing และ middleware, viper สำหรับอ่านค่าคอนฟิกจากไฟล์, cobra สำหรับสร้าง CLI, zap สำหรับบันทึก log แบบมีโครงสร้าง, GORM สำหรับจัดการฐานข้อมูล, gomail และ hermes สำหรับส่งอีเมลแบบ HTML สวยงาม แผนภาพนี้แสดงการเชื่อมโยงเครื่องมือเหล่านี้เข้าด้วยกัน

---

## ภาคที่ 7: การออกแบบสถาปัตยกรรมและ Workflow (บทที่ 46–48)

### แผนภาพ: Clean Architecture Layers

```mermaid
graph TD
    subgraph Delivery
        D1[HTTP Handlers] --> D2[DTOs]
    end

    subgraph Usecase
        U1[Business Logic] --> U2[Repository Interface]
    end

    subgraph Repository
        R1[Repository Interface] --> R2[Implementation<br/>(GORM, Redis)]
    end

    subgraph External
        E1[Database] --> E2[Cache]
    end

    D2 --> U1
    U2 --> R1
    R2 --> E1
    R2 --> E2
```

**คำอธิบายภาษาไทย:**  
สถาปัตยกรรม Clean Architecture แบ่งเป็นสามชั้นหลัก:
- **Delivery**: รับ Request และส่ง Response (เช่น HTTP handlers) ใช้ DTO เพื่อสื่อสารกับชั้นถัดไป
- **Usecase**: ตรรกะธุรกิจ (business logic) เรียกใช้ repository ผ่าน interface
- **Repository**: นำ interface ไป implement ด้วย GORM หรือ Redis เพื่อติดต่อฐานข้อมูลหรือ cache

การแยกชั้นนี้ทำให้โค้ดยืดหยุ่น ทดสอบง่าย และบำรุงรักษาง่าย

---

## ภาคที่ 8: Domain-Driven Design (DDD) (บทที่ 49–51)

### แผนภาพ: Domain → Application → Infrastructure

```mermaid
graph LR
    subgraph Domain
        A[Entities & Value Objects] --> B[Aggregates]
        B --> C[Domain Events]
        C --> D[Repository Interface]
    end

    subgraph Application
        E[Use Cases] --> F[DTOs]
        F --> G[Service Orchestration]
    end

    subgraph Infrastructure
        H[Repository Implementation] --> I[Event Bus]
        I --> J[Message Broker]
    end

    A --> E
    D --> H
    C --> I
```

**คำอธิบายภาษาไทย:**  
ใน DDD ชั้น **Domain** จะเก็บโมเดลธุรกิจหลัก (Entity, Value Object, Aggregate) และกำหนด interface สำหรับ repository รวมถึง domain events  
ชั้น **Application** จะนำ domain มาใช้ใน use cases ผ่าน DTO และจัดลำดับการทำงาน  
ชั้น **Infrastructure** รับผิดชอบการ implement repository และส่ง domain events ผ่าน event bus หรือ message broker ไปยังระบบอื่น  

การแยกแบบนี้ทำให้โมเดลธุรกิจไม่ขึ้นอยู่กับเทคโนโลยีภายนอก

---

## ภาคที่ 9: การผสานระบบภายนอกและคุณลักษณะเสริม (บทที่ 52–58)

### แผนภาพ: Redis, RabbitMQ, MQTT, InfluxDB, WebSocket, Notifications
 
```mermaid
flowchart TD
    subgraph CacheQueue[Cache & Queue]
        R1[Redis] -->|Cache| R2[In-memory Data]
        R1 -->|Queue| R3[Pub/Sub]
    end

    subgraph MessageBroker[Message Broker]
        M1[RabbitMQ] --> M2[Exchanges & Queues]
        M2 --> M3[Workers]
    end

    subgraph IoT[IoT Devices]
        I1[MQTT Broker] --> I2[Sensors]
        I2 --> I3[Raw Data]
    end

    subgraph TimeSeries[Time‑Series Database]
        T1[InfluxDB] --> T2[Metrics & Analytics]
    end

    subgraph Realtime[Real‑time Communication]
        W1[WebSocket Server] --> W2[Live Updates to Clients]
    end

    subgraph Notifications[Notification Channels]
        N1[SMS Gateway] --> N2[LINE Notify]
        N2 --> N3[Discord Webhook]
    end

    R2 -->|Cache hit| T1
    R3 -->|Message| M1
    M3 -->|Processed data| T1
    I3 -->|Sensor telemetry| T1
    T2 -->|Metrics & alerts| W1
    W2 -->|Trigger| N1
```

**คำอธิบายภาษาไทย:**

แผนภาพนี้แสดงการทำงานร่วมกันของระบบภายนอกที่ใช้ในแอปพลิเคชัน Go ระดับ production โดยมีลำดับการไหลของข้อมูลดังนี้:

1. **Cache & Queue (Redis)**  
   - ใช้ Redis เป็น **cache** เก็บข้อมูลในหน่วยความจำ (In-memory Data) เพื่อลดภาระฐานข้อมูล  
   - ใช้ Redis เป็น **queue** (Pub/Sub) สำหรับส่งข้อความระหว่าง services

2. **Message Broker (RabbitMQ)**  
   - RabbitMQ ทำหน้าที่เป็นตัวกลางรับข้อความจาก Redis Pub/Sub ผ่าน exchanges และ queues  
   - Workers จะดึงงานจาก queue และประมวลผล จากนั้นส่งข้อมูลที่ประมวลผลแล้วไปยัง InfluxDB

3. **IoT Devices (MQTT)**  
   - อุปกรณ์ IoT ส่งข้อมูลผ่าน MQTT Broker ไปยังเซ็นเซอร์ (Sensors)  
   - ข้อมูลดิบ (Raw Data) จะถูกส่งต่อไปยัง InfluxDB เพื่อจัดเก็บเป็นอนุกรมเวลา

4. **Time‑Series Database (InfluxDB)**  
   - รับข้อมูลจาก Redis cache, RabbitMQ workers, และ MQTT  
   - เก็บเป็น metrics และ analytics พร้อมให้บริการข้อมูลแก่ระบบอื่น

5. **Real‑time Communication (WebSocket)**  
   - WebSocket Server ดึงข้อมูลจาก InfluxDB เพื่อส่งการอัปเดตแบบเรียลไทม์ไปยังไคลเอนต์

6. **Notifications (SMS, LINE, Discord)**  
   - เมื่อเกิดเหตุการณ์สำคัญ (เช่น ข้อมูลเกินค่าที่กำหนด) WebSocket จะส่งสัญญาณไปยัง SMS Gateway  
   - จากนั้นส่งต่อไปยัง LINE Notify และ Discord Webhook เพื่อแจ้งเตือนผู้ใช้ผ่านช่องทางที่กำหนด

**สรุป:** ระบบทั้งหมดทำงานร่วมกันเป็นวงจร เริ่มจากการรับข้อมูลจาก IoT และแคช/คิวใน Redis, ประมวลผลผ่าน RabbitMQ, เก็บข้อมูลใน InfluxDB, แสดงผลแบบเรียลไทม์ผ่าน WebSocket และแจ้งเตือนผ่านช่องทางต่างๆ ทำให้แอปพลิเคชันมีความเสถียร รองรับปริมาณข้อมูลสูง และตอบสนองได้ทันที

**คำอธิบายภาษาไทย:**  
ระบบสมัยใหม่มักเชื่อมต่อกับหลายบริการ:
- **Redis**: ใช้เป็น cache และ message queue (pub/sub)
- **RabbitMQ**: เป็น message broker สำหรับกระจายงานระหว่าง services
- **MQTT**: รับข้อมูลจากอุปกรณ์ IoT (เซ็นเซอร์)
- **InfluxDB**: เก็บบันทึกข้อมูลอนุกรมเวลา (time‑series) สำหรับ metrics
- **WebSocket**: สื่อสารแบบ real‑time กับไคลเอนต์
- **SMS / LINE / Discord**: แจ้งเตือนผู้ใช้ผ่านช่องทางต่าง ๆ

ข้อมูลจากทุกแหล่งจะถูกนำไปจัดเก็บและแสดงผลแบบเรียลไทม์

---

## ภาคที่ 10: เทมเพลต กระบวนการพัฒนา และตัวอย่างโค้ด (บทที่ 59–63)

### แผนภาพ: ตัวอย่างแอป → Task List → Checklist → Diagram → Config

```mermaid
graph LR
    subgraph Example
        E1[Full‑stack Example] --> E2[Code Structure]
    end

    subgraph Processes
        P1[Task List] --> P2[Checklist]
        P2 --> P3[Workflow Diagram]
    end

    subgraph Config
        C1[mop Config] --> C2[YAML/ENV]
        C2 --> E1
    end

    E2 --> P1
    P3 --> C1
```

**คำอธิบายภาษาไทย:**  
ในภาคสุดท้ายเราจะนำเสนอตัวอย่างแอปพลิเคชันครบวงจร (Full‑stack Example) ซึ่งแสดงโครงสร้างโค้ดตาม Clean Architecture และ DDD จากนั้นใช้ Task List และ Checklist เป็นแนวทางในการพัฒนาและการตรวจสอบคุณภาพ แผนภาพการทำงาน (Workflow Diagram) ช่วยให้ทีมเห็นภาพรวมของกระบวนการ และ mop Config ช่วยจัดการการตั้งค่าผ่านไฟล์ YAML หรือ environment variables เพื่อให้แอปพลิเคชันปรับค่าได้ง่ายเมื่อนำไปใช้งานจริง

---

## วิธีใช้แผนภาพใน draw.io

1. คัดลอกโค้ด Mermaid ที่ต้องการ  
2. ใน draw.io ไปที่ **Arrange → Insert → Advanced → Mermaid**  
3. วางโค้ดแล้วกด **Insert**  
4. แผนภาพจะถูกสร้างขึ้นและสามารถแก้ไขเพิ่มเติมได้

หรือใช้โปรแกรมที่รองรับ Mermaid เช่น GitHub, GitLab, VS Code (ติดตั้ง extension Mermaid)

แผนภาพเหล่านี้ช่วยให้เห็นภาพการทำงานของระบบและแนวคิดในคู่มืออย่างชัดเจน เหมาะสำหรับใช้สอนหรืออธิบายสถาปัตยกรรมให้ทีมฟัง
 

## 📚 บทสรุป

คู่มือนี้ครอบคลุมเนื้อหาตั้งแต่พื้นฐานภาษา Go ไปจนถึงการออกแบบสถาปัตยกรรมระดับองค์กรและการผสานระบบภายนอกที่ใช้ในโลกแห่งความจริง โดยมุ่งเน้นให้ผู้อ่านสามารถนำไปประยุกต์ใช้ได้ทันที

**หวังว่าคู่มือนี้จะเป็นประโยชน์สำหรับผู้ที่ต้องการเริ่มต้นและพัฒนาทักษะการเขียนโปรแกรมด้วย Go อย่างจริงจัง ขอให้สนุกกับการเขียนโปรแกรม!**

---

**ผู้เขียน:** คงนคร จันทะคุณ  
**อีเมล:** kongnakornjantakun@gmail.com  
**วันที่:** เมษายน 2026