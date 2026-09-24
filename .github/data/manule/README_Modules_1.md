# Go Restful API 

**หมายเหตุ เนื้อหาในหนังสือ:**  
เนื้อหาในหนังสือ "ใช้ AI ช่วยเขียน  เพื่อทดสอบ AI Model ผู้เขียนเป็นผู้ออกแบบ ใช้ AI ช่วยจัดเรียง ซึ่งมีค่าใช้จ่ายพอสมควร ให้ใช้ฟรีก่อน ต้องการสนับสนุนเพื่อทำเนื้อหาแนวนี้ต่อ สามารถให้การสนับสนุนได้ครับ ตามกำลังศรัทธา 
📞 โทรศัพท์ / พร้อมเพย์: **0955088091**  
- Code ใช้งานได้จริง
---

An API dev written in Golang with chi-route and Gorm. Write restful API with fast development and developer friendly.

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

### Run
- `docker-compose up`
- OR  go run ./main.go serve  on loca Windows OS
- Swagger: [localhost:5000/swagger/](http://localhost:5000/swagger/)
- http://localhost:5000/swagger/index.html#/

```bash
  Email: root@gmail.com
  Password: root_password
```
## TODO

- Traefik
- Config using .env
- Linter
- Jaeger
- Production docker file version
- Mock database using gomock

## Acknowledgements

- [github.com/dhax/go-base](https://github.com/dhax/go-base)
- [github.com/akmamun/go-fication](https://github.com/akmamun/go-fication)
- [github.com/wpcodevo/golang-fiber-jwt](https://github.com/wpcodevo/golang-fiber-jwt)
- [github.com/wpcodevo/golang-fiber](https://github.com/wpcodevo/golang-fiber)
- [github.com/kienmatu/togo](https://github.com/kienmatu/togo)
- [github.com/AleksK1NG/Go-Clean-Architecture-REST-API](https://github.com/AleksK1NG/Go-Clean-Architecture-REST-API)
- [github.com/bxcodec/go-clean-arch](https://github.com/bxcodec/go-clean-arch)
- [codevoweb.com/golang-and-gorm-user-registration-email-verification/](https://codevoweb.com/golang-and-gorm-user-registration-email-verification/)
- [codevoweb.com/golang-gorm-postgresql-user-registration-with-refresh-tokens/](https://codevoweb.com/golang-gorm-postgresql-user-registration-with-refresh-tokens/)
- [codevoweb.com/how-to-implement-google-oauth2-in-golang/](https://codevoweb.com/how-to-implement-google-oauth2-in-golang/)
- [codevoweb.com/how-to-upload-single-and-multiple-files-in-golang/](https://codevoweb.com/how-to-upload-single-and-multiple-files-in-golang/)
- [codevoweb.com/forgot-reset-passwords-in-golang-with-html-email/](https://codevoweb.com/forgot-reset-passwords-in-golang-with-html-email/)
- [techmaster.vn/posts/34577/kien-truc-sach-voi-golang](https://techmaster.vn/posts/34577/kien-truc-sach-voi-golang)



### Installation


```bash
- ตรวจสอบว่า go.mod มี replace directive หรือใช้ local module หรือไม่
- ถ้า icmongolang เป็น local module ให้ใช้ replace icmongolang => icmongolang
- หรือถ้าเป็น private repo ให้ตั้ง GOPRIVATE และใช้ access token
- Perfect! You're setting up an existing Go project (icmongolang). Here's how to properly set it up and run it:
```bash
## Complete Setup Steps for Your icmongolang Project
## 📘 การจัดการ `go.mod` และ dependencies สำหรับโปรเจกต์ `icmongolang`  
## 📘 Managing `go.mod` and dependencies for `icmongolang` project

> คำแนะนำแบบทีละขั้นตอน (ไทย / อังกฤษ)  
> Step-by-step guide (Thai / English)

---

### 🧱 1. โคลนโปรเจกต์และเข้าไปในโฟลเดอร์  
### 1. Clone and enter the project

```bash
# ไทย: โคลน repository จาก GitHub และเปลี่ยนไปยังไดเรกทอรีโปรเจกต์
# EN: Clone the repository from GitHub and change into the project directory
git clone github.com/kongnakornna/icmongolang.git
cd icmongolang
```
```bash

go mod tidy
go mod download
go mod verify

# ล้างฐานข้อมูลเก่า (ระวังข้อมูล)
# go run cmd/api/main.go migrate:reset (ถ้ามีคำสั่ง)

# รัน migrate ใหม่
go run  main.go migrate


go run ./main.go serve

OR ใช้ go run โดยตรง (ไม่ต้อง build exe)

go run ./cmd/api/main.go serve


 
go mod tidy
go mod download
go mod verify
go run cmd/api/main.go migrate
go mod vendor
air




 ```
# Auto Run 

air

```bash


- for windows 10
 
 - .air.toml
 
    root = "."
    tmp_dir = "tmp"
    env_files = [".env.dev"]   # โหลด env โดยอัตโนมัติ

    [build]
    # ใช้ array: [binary, argument1, argument2, ...]
    entrypoint = ["./tmp/main.exe", "serve"]
    cmd = "go build -o ./tmp/main.exe ./cmd/api"
    env = ["GOOS=windows", "GOARCH=amd64"]
    clean_on_exit = true

    [log]
    time = true

    [misc]
    clean_on_exit = true

 ```


![Icmon5](https://github.com/user-attachments/assets/fa802d05-f4f7-4f60-bab0-897a95cea541)

# 🚀 โครงสร้างและ Workflow ของโปรเจกต์ `icmongolang` (Go Backend Clean Architecture)

เอกสารนี้  ประกอบด้วย  
- โครงสร้างการทำงานแบบละเอียด  
- Dataflow Diagram (Flowchart TB สำหรับ Draw.io) พร้อมคำอธิบาย  
- ตัวอย่างโค้ดพร้อมคอมเมนต์ไทย/อังกฤษ ที่รันได้จริง  
- กรณีศึกษา  
- สรุป (ประโยชน์, ข้อควรระวัง, ข้อดี/เสีย, ข้อห้าม, แหล่งอ้างอิง)
- มี Boilerplate พร้อมนำไปใช้
### Boilerplate คือ โค้ดหรือข้อความรูปแบบมาตรฐานที่สามารถนำกลับมาใช้ใหม่ได้หลายครั้งโดยมีการเปลี่ยนแปลงแก้ไขน้อยมากหรือไม่มีเลย 
- มีวัตถุประสงค์หลักเพื่อลดเวลาในการทำงานซ้ำซ้อน เพิ่มมาตรฐานให้กับชิ้นงาน และช่วยให้โครงสร้างไฟล์เริ่มต้นเป็นระเบียบ 
- จุดเด่นและประโยชน์ของ Boilerplate:
- ความรวดเร็ว: ไม่ต้องเสียเวลาเขียนโค้ดตั้งต้นใหม่ทุกครั้ง
- มาตรฐาน: สร้างความสม่ำเสมอในโค้ดหรือเอกสาร
- ลดข้อผิดพลาด: เนื่องจากใช้โค้ดที่ผ่านการตรวจสอบมาแล้ว


  
![Icmon8](https://github.com/user-attachments/assets/214c3fbe-dd7f-458c-9ea7-b1c13d6fea67)

---

## 1. โครงสร้างการทำงานของโปรเจกต์ (Architecture Overview)

โปรเจกต์ใช้ **Clean Architecture** 3-layer + Delivery:

| Layer | ตำแหน่ง | หน้าที่ |
|-------|---------|--------|
| **Model** | `internal/models/` | Entity (GORM) – `User`, `Session`, `VerificationToken` |
| **Repository** | `internal/repository/` | อ่าน/เขียน DB และ Redis ผ่าน interface |
| **Usecase** | `internal/usecase/` | Business logic: hash, JWT, email queue, validation |
| **Delivery** | `internal/delivery/rest/` | HTTP handlers, middleware, DTO, router |
| **Worker** | `internal/delivery/worker/` | Background job สำหรับส่งอีเมล |
 
### Model คืออะไร?
Models คือโครงสร้างข้อมูล (struct) ที่แทน entity ในฐานข้อมูล หรือข้อมูลที่ใช้ในการสื่อสารระหว่าง layers (DTO) โดยปกติจะสอดคล้องกับตารางใน PostgreSQL และใช้ GORM tags สำหรับ mapping

### มีกี่แบบ?
1. **Entity Model** – สอดคล้องกับตาราง DB โดยตรง (user, session)
2. **DTO (Data Transfer Object)** – ใช้รับ/ส่งข้อมูลระหว่าง API (Request/Response)
3. **Embedded Model** – struct ที่ถูกแทรกใน model อื่น (เช่น BaseModel)
4. **Enum-like Model** – ใช้ iota สำหรับ status constants

### ใช้อย่างไร / นำไปใช้กรณีไหน
- ใช้ GORM annotations (`gorm:"column:name;type:..."`) เพื่อกำหนด schema
- ใช้ JSON tags (`json:"field_name"`) สำหรับ serialization
- ใช้ Validator tags (`validate:"required,email"`) สำหรับ input validation

### ทำไมต้องใช้
- จัดระเบียบโครงสร้างข้อมูลให้เป็นหนึ่งเดียว
- ช่วยให้ GORM สร้างตารางอัตโนมัติ (AutoMigrate)
- แยก business entity ออกจาก database details

### ประโยชน์ที่ได้รับ
- Type safety ใน Go (ไม่ต้องใช้ map[string]interface{})
- ลด boilerplate code สำหรับ CRUD
- รองรับความสัมพันธ์ระหว่างตาราง (Relationships: BelongsTo, HasMany)

### Boilerplate คือ โค้ดหรือข้อความรูปแบบมาตรฐานที่สามารถนำกลับมาใช้ใหม่ได้หลายครั้งโดยมีการเปลี่ยนแปลงแก้ไขน้อยมากหรือไม่มีเลย 
- มีวัตถุประสงค์หลักเพื่อลดเวลาในการทำงานซ้ำซ้อน เพิ่มมาตรฐานให้กับชิ้นงาน และช่วยให้โครงสร้างไฟล์เริ่มต้นเป็นระเบียบ เช่น โครงสร้างพื้นฐานของ HTML หรือการตั้งค่าเริ่มต้นในโปรเจกต์ซอฟต์แวร์ Amazon Web Services
- จุดเด่นและประโยชน์ของ Boilerplate:
- ความรวดเร็ว: ไม่ต้องเสียเวลาเขียนโค้ดตั้งต้นใหม่ทุกครั้ง
- มาตรฐาน: สร้างความสม่ำเสมอในโค้ดหรือเอกสาร
- ลดข้อผิดพลาด: เนื่องจากใช้โค้ดที่ผ่านการตรวจสอบมาแล้ว

### ข้อควรระวัง
- ห้ามเก็บ password plain text (ต้อง hashed)
- ใช้ pointer type สำหรับ nullable fields (`*time.Time` แทน `time.Time`)
- ระวัง zero values (0, "", false) vs null

### ข้อดี
- ชัดเจน, ตรวจสอบได้ตอน compile
- รองรับ GORM hooks (BeforeCreate, AfterUpdate)

### ข้อเสีย
- ต้องเปลี่ยนแปลง struct เมื่อ schema เปลี่ยน
- อาจมีหลาย struct ที่คล้ายกัน (entity vs response DTO)

### ข้อห้าม
- ห้ามใช้ model สำหรับ business logic (ควรอยู่ใน usecase)
- ห้าม serialize model ที่มี password ไปเป็น JSON


### Repository คืออะไร?
Repository Pattern คือตัวกลาง (abstraction) ระหว่าง business logic (usecase) และแหล่งข้อมูล (database, cache, external API) โดยกำหนด interface สำหรับการเข้าถึงข้อมูล และมี implementation ที่เป็นรูปธรรม (PostgreSQL, Redis) แยกออกจากกัน

### มีกี่แบบ?
1. **Specific Repository** – แต่ละ entity มี interface ของตัวเอง (UserRepository, SessionRepository) – ใช้ในโปรเจกต์นี้
2. **Generic Repository** – interface เดียวที่ใช้กับ entity ใดก็ได้ (ใช้ reflection หรือ interface{})
3. **Transaction Repository** – repository ที่มี method สำหรับ transaction (Begin, Commit, Rollback)
4. **Cached Repository** – decorator ที่เพิ่ม cache layer ให้กับ repository หลัก

### ใช้อย่างไร / นำไปใช้กรณีไหน
- ใช้ interface เพื่อกำหนด method ที่ usecase จะเรียก
- implementation จริงใช้ GORM สำหรับ PostgreSQL และ go-redis สำหรับ Redis
- usecase ไม่รู้ว่าข้อมูลมาจาก DB หรือ Cache
- เหมาะกับระบบที่ต้องเปลี่ยนแหล่งข้อมูลบ่อย หรือต้องการ mock สำหรับ unit test

### ทำไมต้องใช้
- แยก business logic ออกจาก data access code
- ทดสอบ usecase ได้ง่ายโดยใช้ mock repository
- สามารถเปลี่ยนจาก PostgreSQL เป็น MongoDB ได้โดยไม่ต้องแก้ usecase

### ประโยชน์ที่ได้รับ
- ลด dependency coupling
- โค้ดสะอาดขึ้น (Clean Architecture)
- รองรับ caching, logging, monitoring ได้โดยไม่แก้ usecase

### ข้อควรระวัง
- repository ควรคืนค่าเป็น model ของ domain (internal/models) ไม่ใช่ DTO
- repository ไม่ควรมี business logic (เช่น การตรวจสอบว่า email ซ้ำ ควรอยู่ใน usecase)
- ระวังเรื่อง transaction: ถ้าต้องการ atomic operation ควรส่ง transaction object (`*gorm.DB`) เข้าไปใน method

### ข้อดี
- ทดสอบง่าย, เปลี่ยนแหล่งข้อมูลได้, แยกความรับผิดชอบชัดเจน

### ข้อเสีย
- มีโค้ดเพิ่มขึ้น (interface + implementation หลายตัว)
- อาจเพิ่มความซับซ้อนในโปรเจกต์เล็ก

### ข้อห้าม
- ห้ามเรียก repository โดยตรงจาก handler (ต้องผ่าน usecase)
- ห้ามใช้ repository ใน repository อื่น (ควรใช้ service หรือ usecase แทน)
- ห้ามใส่ context ลงใน struct repository (ควรส่งผ่าน method argument)

### Usecase คืออะไร?
Usecase (หรือ Service layer) คือชั้นที่บรรจุ **business logic** ของแอปพลิเคชัน ทำหน้าที่ประสานงานระหว่าง repository ต่างๆ ตรวจสอบกฎทางธุรกิจ และแปลงข้อมูลจากรูปแบบของ repository ให้เป็นรูปแบบที่ delivery (handler) ต้องการ โดย usecase **ไม่รู้** ว่า repository ใช้ PostgreSQL หรือ Redis หรือ external API

### มีกี่แบบ?
1. **Specific Usecase** – แต่ละฟีเจอร์มี usecase ของตัวเอง (AuthUsecase, UserUsecase) – ใช้ในโปรเจกต์นี้
2. **Generic Usecase** – ใช้ interface เดียวกันกับหลาย entity (ไม่ค่อยพบใน Go)
3. **Command/Query Segregation** – แยก Usecase สำหรับการแก้ไขข้อมูล (Command) และการอ่านข้อมูล (Query)
4. **Domain Service** – เมื่อ logic ซับซ้อนและเกี่ยวข้องกับหลาย entity

### ใช้อย่างไร / นำไปใช้กรณีไหน
- ใช้ใน handler: `authUsecase.Login(ctx, email, password)` 
- usecase จะเรียก repository method ต่างๆ และอาจเรียกใช้ transaction manager
- คืนค่า business result หรือ error (ไม่คืน HTTP status code)

### ทำไมต้องใช้
- ป้องกัน business logic กระจายอยู่ใน handler หรือ repository
- ทำให้ทดสอบ business logic ได้โดยไม่ต้องมี HTTP request หรือ database จริง (ใช้ mock repository)
- สอดคล้องกับ Clean Architecture

### ประโยชน์ที่ได้รับ
- เปลี่ยน business logic ได้โดยไม่กระทบ delivery (handler) และ repository
- รองรับการ reuse logic (handler เดียวกันใช้ usecase เดียว)
- ง่ายต่อการเพิ่ม logging, tracing, metrics ในชั้นเดียว

### ข้อควรระวัง
- usecase **ห้าม** import package `net/http` หรือ `gin/chi` เพราะจะทำให้ coupling กับ delivery
  ***
  ## Coupling คืออะไร? (ในบริบทการออกแบบซอฟต์แวร์)

**Coupling (การผูกพัน)** คือระดับที่ **โมดูล / คลาส / คอมโพเนนต์** หนึ่งต้องพึ่งพาอีกโมดูลหนึ่ง **มากน้อยแค่ไหน**

- **High coupling (ผูกพันสูง)** – เปลี่ยนอะไรที่ A แล้ว B พังไปหมด  
- **Low coupling (ผูกพันต่ำ)** – แต่ละส่วนเป็นอิสระ เปลี่ยนแปลงได้โดยไม่กระทบกันมาก

---

### ประเภทของ Coupling (เรียงจากแย่ที่สุดไปดีที่สุด)

| ประเภท | คำอธิบาย | ตัวอย่าง |
|--------|----------|----------|
| **Content coupling** | โมดูลเข้าถึงข้อมูลภายในของอีกโมดูลโดยตรง | `otherModule.internalVar = 5` |
| **Common coupling** | ใช้ global variable หรือ shared state ร่วมกัน | `var db *sql.DB` ทั่วทั้งโปรแกรม |
| **Control coupling** | ส่ง flag หรือ parameter เพื่อควบคุมลำดับการทำงานของอีกโมดูล | `ProcessData(shouldSave bool)` |
| **Stamp coupling** | ส่งโครงสร้างข้อมูลที่ใหญ่เกินจำเป็น (ทั้ง struct ทั้งที่ใช้แค่ฟิลด์เดียว) | `func Save(user User)` แต่ใช้แค่ `user.ID` |
| **Data coupling** | ส่งเฉพาะข้อมูลที่จำเป็นผ่านพารามิเตอร์ | `func Save(userID int)` ✅ |
| **No coupling** | ไม่มีการพึ่งพากันเลย (หายาก) | |

---

### ทำไมต้องสนใจ Coupling?

- **Low coupling + High cohesion (การเกาะกลุ่มกันภายใน)** = โค้ดที่บำรุงรักษาง่าย  
- ถ้า coupling สูง →  
  - แก้ไขที่หนึ่งแล้วกระทบหลายที่  
  - ทดสอบยาก (ต้อง mock เยอะ)  
  - reuse โมดุลยาก  
  - เข้าใจระบบยาก

---

### Coupling กับ Worker ใน Go (เชื่อมกับคำถามก่อนหน้า)

ใน pattern **Worker Pool** ถ้าออกแบบไม่ดีจะเกิด coupling สูง เช่น

```go
// High coupling: worker รู้จัก database โดยตรง
func worker(jobs <-chan Job) {
    for job := range jobs {
        db.Exec("INSERT ...", job.Data) // coupling กับ DB driver
    }
}
```

วิธีลด coupling:

- ส่ง **interface** ให้ worker แทนคอนกรีต type  
  ```go
  type Saver interface { Save(data []byte) error }
  func worker(jobs <-chan Job, saver Saver) { ... }
  ```
- ใช้ **dependency injection**  
- ใช้ **message queue** (RabbitMQ, Kafka) เป็นตัวกลาง – workers ไม่รู้จักกันและกัน

---

### ข้อควรปฏิบัติ

✅ **ชอบ data coupling / message coupling** (ผ่าน channel หรือ queue)  
✅ **ใช้ interface เพื่อลด coupling**  
✅ **แยก business logic ออกจาก infrastructure** (DB, HTTP, file)  
❌ **ห้ามใช้ global state ร่วมกันระหว่าง workers**  
❌ **ห้ามให้ worker เรียก method อีก worker โดยตรง** (ควรผ่าน channel)

---

### สรุป

> **Coupling** = ระดับการพึ่งพากันระหว่างโมดูล  
> **Low coupling** = ดี – เปลี่ยนง่าย, ทดสอบง่าย, reuse ได้  
> **High coupling** = ร้าย – โค้ดเปราะบาง, แก้ไขลำบาก  

ใน Go การใช้ **channel, interface, dependency injection** ช่วยให้ workers มี coupling ต่ำและยืดหยุ่นมากขึ้น
  ***
- usecase **ห้าม** ส่งออก HTTP status code หรือ JSON
- ควรใช้ interface สำหรับ usecase เพื่อให้ handler มองเห็นแค่ method ที่จำเป็น

### ข้อดี
- แยก business logic ชัดเจน
- ทดสอบ unit ได้ง่าย (ใช้ mock)
- ปรับเปลี่ยน flow ได้โดยไม่แก้ handler

### ข้อเสีย
- เพิ่ม layer ทำให้มีไฟล์มากขึ้น
- มือใหม่อาจเข้าใจยากว่าควรใส่ logic ตรงไหน (repository หรือ usecase)

### ข้อห้าม
- ห้ามเรียก handler โดยตรงจาก usecase
- ห้ามใช้ `*gorm.DB` ใน usecase (ใช้ repository interface แทน)
- ห้ามใช้ context เพื่อส่งค่าที่ไม่เกี่ยวกับ request (ใช้ argument ปกติ)

---
### Delivery คืออะไร?
Delivery layer เป็นชั้นที่อยู่ด้านนอกสุดของ Clean Architecture ทำหน้าที่รับ request จากผู้ใช้ (HTTP, gRPC, CLI) แปลงข้อมูล, เรียกใช้ usecase, และส่ง response กลับ โดยไม่มีการประมวลผลทางธุรกิจใดๆ

### มีกี่แบบ?
1. **HTTP/REST handlers** – รับ JSON, เรียก usecase, ส่ง JSON response
2. **Middleware** – ทำงานก่อน/หลัง handlers (authentication, logging, CORS, rate limiting)
3. **WebSocket handlers** – จัดการ real-time connections
4. **gRPC services** – สำหรับ internal microservices
5. **CLI commands** – สำหรับ admin tasks (migrate, seed)

### ใช้อย่างไร / นำไปใช้กรณีไหน
- Handler: แปลง HTTP request → usecase input, usecase output → HTTP response
- Middleware: ตรวจสอบ token, log request, จำกัด rate, เพิ่ม security headers
- DTO: กำหนดโครงสร้าง JSON request/response (แยกจาก entity model)
- Router: จับคู่ path กับ handler และ middleware

### ทำไมต้องใช้
- แยกการรับ/ส่งข้อมูลออกจาก business logic
- เปลี่ยนจาก REST เป็น GraphQL ได้โดยไม่ต้องแก้ usecase
- จัดการ cross-cutting concerns (logging, auth) เป็น centralized

### ประโยชน์ที่ได้รับ
- เปลี่ยนรูปแบบ API (REST → gRPC) โดยไม่กระทบ usecase
- ทดสอบ handler แบบ integration ได้ง่าย
- middleware reuse

### ข้อควรระวัง
- handler ควรสั้น (แค่ binding, validation, call usecase, response)
- อย่าใส่ business logic ใน handler
- DTO ควรแยกจาก entity model เพื่อป้องกันข้อมูล泄露 (password hash)

### ข้อดี
- แยกชัดเจน, ยืดหยุ่น, middleware จัดการ统一

### ข้อเสีย
- มีไฟล์จำนวนมาก (handler, dto, middleware แต่ละตัว)
- อาจมีการ mapping ซ้ำซ้อน (entity → dto)

### ข้อห้าม
- ห้ามเรียก repository โดยตรงจาก handler
- ห้ามทำ business logic (if-else ที่เกี่ยวกับธุรกิจ) ใน handler
- ห้ามใช้ entity model เป็น request DTO ถ้ามี field ที่ไม่ต้องการให้ client ส่งมา



## Golang: Worker คืออะไร?

Worker ในภาษา Go คือ **กระบวนการทำงานเบื้องหลัง (background process)** ที่ทำงานในรูปแบบ concurrent โดยใช้ **goroutine** และรับงานผ่าน **channel** หรือระบบ queue ต่าง ๆ เพื่อประมวลผลแบบไม่รอ (non-blocking) ช่วยให้โปรแกรมหลักทำงานต่อไปได้โดยไม่ต้องรอผลลัพธ์จากงานหนักหรืองานที่ใช้เวลานาน

---

## มีกี่แบบ?

แบ่งตามรูปแบบการทำงานและได้ดังนี้

1. **Single Worker**  
   ใช้ goroutine ตัวเดียวรับงานจาก channel และประมวลผลทีละงาน – เหมาะกับงานที่เรียงลำดับ

2. **Worker Pool**  
   มี goroutine หลายตัว (จำนวนคงที่) แชร์ channel เดียวกัน รับงานมาแล้วกระจายไปยัง workers – เพิ่ม throughput

3. **Scheduled / Cron Worker**  
   ทำงานตามเวลาที่กำหนด เช่น ทุกเที่ยงคืน – ใช้ `time.Ticker` หรือ library `cron`

4. **CLI Command Worker**  
   ทำงานแบบครั้งเดียวจบ (one-off) สำหรับงานระบบ เช่น  
   - `go run cmd/migrate/main.go` – migrate database  
   - `go run cmd/seed/main.go` – seed ข้อมูลเริ่มต้น  
   ทำงานแยกจาก main application ไม่รันตลอดเวลา

5. **Message Queue Consumer**  
   ฟัง messages จาก RabbitMQ, Kafka, NATS – workers จะรอและประมวลผลแบบ long-running

---

## ใช้อย่างไร / นำไปใช้กรณีไหน?

- **API Server ที่ต้องส่งอีเมล / อัปโหลดไฟล์** – ส่งงานเข้า channel ให้ worker จัดการ async  
- **ประมวลผลรูปภาพหรือข้อมูลจำนวนมาก** – ใช้ worker pool แบ่งงานกันทำ  
- **งานประจำ (batch jobs)** – สรุปยอดขายทุกเที่ยงคืน, ล้าง cache  
- **CLI admin tasks** – `migrate`, `seed`, `backup` ใช้เป็น disposable worker  
- **ระบบ Event-driven** – worker คอย consume events จาก Kafka แล้วอัปเดตฐานข้อมูล

ตัวอย่าง Worker Pool ง่าย ๆ ใน Go:

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        results <- job * 2 // simulate work
    }
}

func main() {
    const numWorkers = 5
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    for w := 1; w <= numWorkers; w++ {
        go worker(w, jobs, results)
    }

    // ส่งงาน
    for j := 1; j <= 20; j++ {
        jobs <- j
    }
    close(jobs)
}
```

---

## ทำไมต้องใช้?

- **ไม่ block main goroutine** – โดยเฉพาะใน web server ที่ต้องตอบสนอง client เร็ว  
- **ใช้ทรัพยากรอย่างคุ้มค่า** – Go goroutine มี overhead ต่ำ (เริ่มต้น ~2KB stack)  
- **แยกความรับผิดชอบ (separation of concerns)** – โค้ดส่วน worker ไม่ปนกับ business logic  
- **รองรับการ scaling** – เพิ่มจำนวน worker ได้ง่ายเมื่อโหลดสูง  
- **ทำงานขนาน (parallelism)** – ถ้ามี CPU หลาย core ก็ทำงานพร้อมกันจริง

---

## ประโยชน์ที่ได้รับ

- **Response time ดีขึ้น** – งาน async ไม่ทำให้ client รอนาน  
- **ใช้ goroutine ได้ตรงตามหลัก Go** – “Do not communicate by sharing memory; instead, share memory by communicating”  
- **ทนทานต่อ traffic สูง** – worker pool ควบคุมจำนวน goroutine ไม่ให้หลุดมือ  
- **โครงสร้างโค้ดสะอาด** – แยก worker logic ออกเป็น function / package ได้  
- **รองรับ graceful shutdown** – ใช้ `context` สั่งให้ worker หยุดรับงานใหม่และจบงานปัจจุบัน

---

## ข้อควรระวัง

- **Goroutine leak** – ลืมปิด channel หรือไม่รับ signal ให้ worker ออก → memory leak  
- **Race condition** – ถ้า workers แชร์โครงสร้างข้อมูลเดียวกัน ต้องใช้ `sync.Mutex` หรือ channel  
- **Error handling** – ถ้า worker panic แล้วไม่ recover ทั้งโปรแกรมอาจ crash  
- **Buffer overflow** – channel ที่ไม่มี buffer หรือ buffer น้อยเกินไปอาจ deadlock หรือ block  
- **ไม่กำหนด timeout** – worker อาจทำงานค้างตลอดไป → ใช้ `context.WithTimeout`  
- **Graceful shutdown ไม่สมบูรณ์** – ทำให้งานค้างหายหรือข้อมูลเสีย

---

## ข้อดี

- **น้ำหนักเบา** – สร้าง goroutine นับหมื่นได้โดยไม่กิน memory มาก  
- **ง่ายด้วย channel** – สื่อสารระหว่าง worker และ main ได้เป็นธรรมชาติ  
- **concurrency สูง** – ได้ parallelism โดยไม่ต้องยุ่งยากเรื่อง thread  
- **เหมาะกับ I/O-bound และ CPU-bound** – ปรับจำนวน worker ตามลักษณะงาน  
- **มี pipeline pattern** – ต่อ worker หลายชั้น (fan-out, fan-in) ได้สวยงาม

---

## ข้อเสีย

- **ดีบักยากกว่า synchronous code** – race condition, deadlock หายาก  
- **ต้องออกแบบ concurrency ให้ดี** – มิฉะนั้น performance แย่กว่า sequential  
- **test ยากขึ้น** – ต้องจำลอง channel, time, context  
- **overhead เล็กน้อย** – channel และ scheduler มี cost (แต่ต่ำมาก)  
- **ควบคุมลำดับงานยาก** – ถ้าต้องการ strict order อาจต้องมี worker เดียว

---

## ข้อห้าม (สิ่งที่ควรหลีกเลี่ยง)

- **อย่าสร้าง worker แบบไม่จำกัด** – ใช้ worker pool แทน goroutine ระเบิด  
- **อย่าใช้ global state ร่วมกันโดยไม่ sync** – ใช้ mutex หรือ channel แทน  
- **อย่าละเลย context cancellation** – worker ควรเช็ค `ctx.Done()` เสมอ  
- **อย่าปิด channel ขณะมี worker กำลังใช้อยู่** – จะ panic  
- **อย่าให้ worker ทำ blocking I/O แบบ infinite** – ควรมี timeout หรือ circuit breaker  
- **อย่าใช้ worker สำหรับงานที่ต้อง transaction ACID โดยไม่ careful** – rollback ยาก  
- **อย่าเพิกเฉยต่อ worker panic** – ควรมี `recover` และ log

---

## สรุป

Worker ใน Go เป็น **concurrency pattern ที่ใช้ goroutine + channel** เพื่อทำงานเบื้องหลังแบบ async, parallel หรือ scheduled  
มีหลายรูปแบบ เช่น single worker, worker pool, cron worker, CLI command worker, message queue consumer  
**ประโยชน์หลัก** – ช่วยเพิ่ม throughput, ลด response time, ใช้ทรัพยากรน้อย, รองรับการ scaling  
**ข้อควรระวัง** – goroutine leak, race, graceful shutdown, error handling  
**ข้อห้าม** – ห้ามสร้าง worker แบบไร้ขีดจำกัด, ห้ามใช้ shared state โดยไม่ป้องกัน, ห้ามลืม context  
โดยเฉพาะ **CLI command** (เช่น `migrate`, `seed`) ก็ถูกมองเป็น worker แบบ one-off ที่ใช้สำหรับ admin tasks โดยเฉพาะ – ทำงานแยกจาก main app ช่วยรักษาความสะอาดของระบบและความปลอดภัย

---

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

- NodeJS type script  convert to golang   GORM entities
 
### โฟลเดอร์หลัก (Modules)
โปรเจกต์ใช้ **Clean Architecture** 3-layer + Delivery:
| Layer | ตำแหน่ง | หน้าที่ |
|-------|---------|--------|
| **Model** | `internal/models/` | Entity (GORM) – `User`, `Session`, `VerificationToken` |
| **Repository** | `internal/repository/` | อ่าน/เขียน DB และ Redis ผ่าน interface |
| **Usecase** | `internal/usecase/` | Business logic: hash, JWT, email queue, validation |
| **Delivery** | `internal/delivery/rest/` | HTTP handlers, middleware, DTO, router |
| **Worker** | `internal/delivery/worker/` | Background job สำหรับส่งอีเมล |
 
```bash
  api/
    ├── cmd/                     # Cobra CLI (serve, migrate, initdata, worker)
    │   ├──apir/
    │   │  └── main.go
    │   ├── initdata.go             
    │   ├── root.go
    │   ├── serve.go
    │   └── worker.go        
    ├── config/                  # Viper config (YAML + env)
    │   ├── config.default.yml
    │   ├── config.dev.yml
    │   └── config.go  
    ├── internal/                # โค้ดส่วนตัว (ไม่ถูก import จากภายนอก)
    │   ├── iot/                 # shared packages (jwt, redis, email, logger, hash, utils)
    │   │   ├── delivery/        #  HTTP handlers, middleware, dto, router
    │   │   │     └── http/
    │   │   │          ├── handler.go
    │   │   │          └── routes.go
    │   │   ├── worker/        
    │   │   │     └──worker.go
    │   │   ├── helper/          #helper
    │   │   │     └──alarm.go
    │   │   ├── models/          # GORM entities
    │   │   │     ├── alarm.go
    │   │   │     ├── common.go
    │   │   │     └── device_type.go
    │   │   ├── presenter/          
    │   │   │     └── presenter.go
    │   │   ├── repository/         
    │   │   │     ├── alarm_log_repo.go
    │   │   │     ├── device_repo.go
    │   │   │     └── schedule_repo.go
    │   │   ├── usecase/          # GORM entities
    │   │   │     └── usecase.go
    │   │   └── iot.go
    │   ├── repository/          # interfaces + impl (postgres, redis)
    │   │       ├── pg.go
    │   │       └── redis.go
    │   ├── server/           
    │   │       ├── handlers.go
    │   │       └── server.go
    │   ├── usecase/             # business logic
    │   │       └── usecase.go
    │   ├── pg_repository.go   
    │   ├── redis_repository.go   
    │   ├── usecase.go   
    │   ├── pkg/                 #  shared packages (jwt, redis, email, logger, hash, utils)
    │   │   ├── db/              #  DB postgres,redis
    │   │   │   ├── postgres/ 
    │   │   │   │     └── db_conn.go
    │   │   │   └── redis/ 
    │   │   │          └── redis_conn.go
    │   │   ├── helpers/
    │   │   │   ├── iot.go          # Alarm logic (สมบูรณ์)
    │   │   │   └── format.go       # ฟังก์ชันช่วยเหลือ (time, string, random)
    │   │   ├── websocket/
    │   │   │   └── websocket.go   
    │   │   ├── mqtt/
    │   │   │   └── client.go       # MQTT client พร้อม GetDataFromTopic
    │   │   ├── influxdb/
    │   │   │   └── client.go       # InfluxDB client
    │   │   └── httpErrors/
    │   │       └── httpErrors.go   # httpErrors
    ├── migrations/              # raw SQL (optional)
    ├── docker-compose.dev.yml   # Postgres + Redis + MailHog
    ├── Dockerfile.dev / .air.toml
    └── go.mod
```
# WebSocket server

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

---
# ----
go mod tidy
go mod download
go mod verify
go run cmd/api/main.go migrate
swag init -g cmd/api/main.go 
mockery --all   
go clean -cache        
go clean -modcache
go mod vendor 
go test ./...     
air


  


# ----

go clean -cache        
go clean -modcache    
mockery --all     
go test ./...  
go build ./...   
go test -v ./internal/modules/mqtt/delivery/http -run TestMQTTHandler_Publish
go test -cover ./internal/modules/mqtt/...        
go test ./internal/modules/mqtt/...        


----  

go clean -cache        
go clean -modcache    
go mod tidy
go mod download
go mod verify
go run cmd/api/main.go migrate
swag init -g cmd/api/main.go 
go mod vendor 
go test ./...     
air 


----  

- http://localhost:5000/api/iot/monitordevicegroup?bucket=BAACTW01&location_id=&hardware_id=&lang=th&delcache=0


go clean -cache        
go clean -modcache    
go mod tidy  
swag init -g cmd/api/main.go  
air 

lsnrctl start
net start OracleServiceXE

### วิธีที่ 1: รันทีละคำสั่ง (ง่ายที่สุด)
- คัดลอกและวางทีละบรรทัดลงใน PowerShell:
# 1. ล้าง cache
go clean -cache
go clean -modcache

# 2. จัดการ dependencies
go mod tidy
go mod download
go mod verify

# 3. Migrate database
go run cmd/api/main.go migrate

# 4. สร้าง Swagger docs (ต้องติดตั้ง swag ก่อน)
swag init -g cmd/api/main.go

# 5. Vendor (ถ้าต้องการ)
go mod vendor

# 6. ทดสอบ
go test ./...

# 7. รัน server (ใช้ air หรือ go run)
air
# หรือถ้าไม่มี air:
go run cmd/api/main.go serve