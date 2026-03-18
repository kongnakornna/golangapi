# คู่มือการพัฒนา

## การเตรียมสภาพแวดล้อม
- Go 1.25 ขึ้นไป
- PostgreSQL 12 ขึ้นไป
- Redis 6 ขึ้นไป (ไม่จำเป็น)

เริ่มต้นการกำหนดค่า:

```bash
cp configs/config.example.yaml configs/config.yaml
```

โดยค่าเริ่มต้น จะอ่าน `configs/config.yaml` ซึ่งสามารถระบุได้ผ่าน `CONFIG_PATH`
สามารถแทนที่การกำหนดค่าได้ผ่านตัวแปรสภาพแวดล้อม (ขึ้นต้นด้วย `APP_` เช่น `APP_DB_HOST`)


## คำสั่งที่ใช้บ่อย
- การพัฒนาและการรัน: `./scripts/dev.sh`
- รันโดยตรง: `go run cmd/app/main.go`
- สร้าง Swagger: `./scripts/swagger.sh`
- รันการทดสอบ: `go test ./...`

## กระบวนการพัฒนาโมดูล

ยกตัวอย่างการเพิ่มโมดูล `order`:

1. สร้างไดเร็กทอรี: `internal/core/order/{handler,service,repository,dto,model}`
2. สร้างคลาส handler/service/repository
3. ลงทะเบียนเส้นทางใน `internal/core/order/routes.go`
4. อัปเดตการฉีด: `internal/apps/app/bootstrap/injection/*`
5. อัปเดตการกำหนดค่าเส้นทาง: `internal/apps/app/router` (หากต้องการรวมการกำหนดเส้นทาง v1)

## การเพิ่ม API Endpoints
- เพิ่มเมธอด ในโมดูล `handler`
- ลงทะเบียนในโมดูล `routes.go`
- กรอกข้อมูล Swagger annotations ให้ครบถ้วนและรัน `./scripts/swagger.sh`

## ส่วนขยายสำหรับหลายแอปพลิเคชัน

- เพิ่ม `cmd/<app>/main.go`
- เพิ่ม `internal/apps/<app>/bootstrap` และ `router`
- นำโมดูล `internal/core` กลับมาใช้ใหม่

## มาตรฐานการเขียนโค้ด

- `gofmt` (การเยื้องด้วยแท็บ)
- ชื่อไฟล์: snake_case
- ลำดับชั้นของโมดูลที่ชัดเจน: handler/service/repository