
 
swag init -g cmd/api/main.go 

/api/settings/listschedulepage?page=1&pageSize=20&event=1&keyword=TASK&start=05:10
http://localhost:3003/v1/settings/listschedulepage?page=1&event=&pageSize=10&keyword=



อ่าน
C:\github\gistdaapi\src\modules\settings
C:\github\gistdaapi\src\modules\settings\settings.module.ts
C:\github\gistdaapi\src\modules\settings\settings.controller.ts
C:\github\gistdaapi\src\modules\settings\settings.service.ts
C:\github\gistdaapi\src\modules\settings\entities
C:\github\gistdaapi\src\modules\settings\dto

Update
C:\github\icmongolang\internal\modules\settings

ให้ มี ข้อมูลแสดงผลและ business เหมือนกับ  C:\github\gistdaapi\src\modules\settings

โดยใช้โครงสร้าง  C:\github\icmongolang 

ทำการ  update code C:\github\icmongolang\internal\modules\settings  เอกสารได้เลย ไม่ต้องถาม 








อ่าน
C:\github\icmongolang\internal\modules\users
## 1. ปัญหาเดิม

ในโปรเจคใช้ `time.Now()` ตรงๆ สำหรับ datetime fields เช่น `updateddate`, `createddate`, `CreatedAt`, `UpdatedAt` ซึ่งจะใช้ **UTC timezone ของ server** ไม่ใช่ Asia/Bangkok ทำให้:

- เวลาใน DB ไม่ตรงกับเวลาจริงของผู้ใช้
- API response ส่งเวลา UTC กลับไป
- Log และ audit trail ไม่สอดคล้องกัน

---

## 2. รูปแบบที่ใช้แก้ไข

### 2.1 สำหรับ `time.Time` struct fields

```go
// BEFORE
CreatedAt: time.Now(),
UpdatedAt: time.Now(),

// AFTER
CreatedAt: time.Now().In(helpers.GetTimeLocation()),
UpdatedAt: time.Now().In(helpers.GetTimeLocation()),
```

### 2.2 สำหรับ map updates (GORM)

```go
// BEFORE
"updateddate": time.Now(),

// AFTER
"updateddate": time.Now().In(helpers.GetTimeLocation()),
```

### 2.3 สำหรับ string formatting

```go
// BEFORE
time.Now().Format("2006-01-02 15:04:05")

// AFTER
helpers.GetCurrentFullDatenow()
```

---

## 3. รายละเอียดการแก้ไข

### 3.1 `internal/modules/settings/usecase/usecase.go`

| บรรทัด | โค้ดเดิม | โค้ดใหม่ |
|--------|----------|----------|
| 95 | `"updateddate": time.Now()` | `"updateddate": time.Now().In(helpers.GetTimeLocation())` |
| 102 | `"updateddate": time.Now()` | `"updateddate": time.Now().In(helpers.GetTimeLocation())` |

**สิ่งที่เพิ่ม:** import `"icmongolang/pkg/helpers"`


skill
1. สั่ง copilot  
   - copilot
   - /tb-brainstorming
   - แล้ว เอาโจทย์มา วาง เพื่อให้  AI วิเคาะห์ 
   - /tb-writing-plans
   -  ถาม และ ให้ ทางเลือก แสดงแนวทางตัดสินใจ ไปจน กว่าจะ สมบุรณ์
   - /tb-executing-plans
   - เริ่ม ทำงาน ตาม plans


/tb-scrutinize 
- Affected กับส่วนไหนบ้าง
- panic
- รายงานสรุปผลการดำเนิดการ อะไรไปบ้าง และบอกกระบวนการทดสอบผล



C:\github\icmongolang\internal\modules\iot\usecase\usecase.go



