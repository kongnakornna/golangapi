# คู่มือการใช้งาน Repository (Repo)

## โครงสร้างโปรเจกต์และการจัดระเบียบโมดูล
จุดเริ่มต้นของ service อยู่ที่ `cmd/app/main.go`  
การประกอบแอปพลิเคชันอยู่ที่ `internal/apps/app/bootstrap/`  
การประกอบเส้นทาง (routing) อยู่ที่ `internal/apps/app/router/`  
โมดูลหลัก (domain) อยู่ที่ `internal/core/` (เช่น `auth`, `user`, `health`)  
แต่ละโมดูลจะมีโครงสร้างย่อย เช่น `handler`, `service`, `repository`, `dto`, `model`, `routes`  
โครงสร้างพื้นฐาน (infrastructure) อยู่ที่ `internal/platform/` (เช่น `config`, `db`)  
เครื่องมือสำหรับ transport layer อยู่ที่ `internal/transport/` (เช่น `httpx`, `middleware`)  
ส่วนประกอบที่สามารถนำมาใช้ซ้ำได้ (reusable components) อยู่ที่ `pkg/` (cache, jwt, logger, queue, transaction, utils)  
เอกสาร API ที่ถูกสร้างจะอยู่ที่ `api/app/`  
การตั้งค่าอยู่ที่ `configs/`  
SQL migrations อยู่ที่ `migrations/`  
ทรัพยากรสำหรับ deploy อยู่ที่ `deploy/`  
สคริปต์ต่างๆ อยู่ที่ `scripts/`

## คำสั่ง build, test และพัฒนา
- `./scripts/dev.sh` ใช้สร้าง Swagger และใช้ `air` สำหรับ hot reload (ดู `.air.toml`)
- `go run cmd/app/main.go` สำหรับรัน service โดยตรง
- `./scripts/swagger.sh` ใช้สร้าง Swagger ใหม่ไปยัง `api/app/`
- `./scripts/build.sh` ใช้ build binary สำหรับหลาย ๆ แพลตฟอร์มไปยัง `build/`
- `go test ./...` ใช้รัน tests ทั้งหมด; สามารถเพิ่ม `-race` หรือ `-coverprofile=coverage.out`
- `docker compose -f deploy/docker/docker-compose.yaml up -d` ใช้เริ่ม PostgreSQL และ Redis

## รูปแบบการเขียนโค้ดและการตั้งชื่อ
ใช้ `gofmt` (ย่อหน้า (indent) ด้วย tab)  
ชื่อ package เป็นตัวพิมพ์เล็ก, ชื่อไฟล์ใช้ snake_case (เช่น `user_service.go`, `auth_handler.go`)  
identifier ที่ export ใช้ PascalCase, ที่ไม่ export ใช้ camelCase  
เมื่อเพิ่ม type ใหม่ ควรใส่ไว้ในโมดูลหลักที่เกี่ยวข้องก่อน และรักษาการแบ่งชั้น (layer) handler/service/repository ให้ชัดเจน

## ข้อกำหนดการทดสอบ
tests จะอยู่ไดเรกทอรีเดียวกับโค้ด, ชื่อไฟล์ `*_test.go`, ชื่อฟังก์ชัน `TestXxx`  
โปรเจกต์ใช้ `testify` (ดูตัวอย่างได้ที่ `internal/core/user/service/user_service_test.go`)  
เมื่อแก้ไข logic ของ service สามารถรัน: `go test ./internal/core/user/service/`

## ข้อกำหนดการ commit และ PR
ข้อความ commit ให้ใช้รูปแบบ prefix ตามที่เคยใช้ (เช่น `feat: ...`, `fix: ...`)  
PR ควรมีคำอธิบายสั้น ๆ, คำสั่งสำหรับทดสอบ, และระบุว่ามีการอัปเดตการตั้งค่า/migration/Swagger หรือไม่ (หากมีการเปลี่ยนแปลงที่ `api/app/` กรุณาระบุด้วย)

## การตั้งค่าและความปลอดภัย
สำหรับการพัฒนาในเครื่อง ให้คัดลอก `configs/config.example.yaml` เป็น `configs/config.yaml` และสามารถใช้ environment variable ที่ขึ้นต้นด้วย `APP_` เพื่อแทนที่ค่าได้  
อย่า commit คีย์สำคัญ (secret); สำหรับ production ให้ใช้ `configs/config.production.yaml` หรือ environment variable ในการแทนที่  
หากมีการแก้ไข route หรือโครงสร้าง request/response กรุณารัน `./scripts/swagger.sh` เพื่ออัปเดตเอกสาร

## การตอบกลับ
- ใช้ภาษาไทยในการสื่อสารและตอบกลับ
- เอกสารใช้ภาษาไทยในรูปแบบ Markdown