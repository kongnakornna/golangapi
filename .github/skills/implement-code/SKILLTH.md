---
name: implement-code
description: "ทักษะการสร้างโค้ด Backend Go ของ PMS ใช้เมื่อ: กำลัง implement ฟีเจอร์ใหม่, เพิ่ม endpoint, สร้าง subtask, เขียน service/repository/handler, สร้าง unit test, หรือเพิ่ม Swagger docs ในโปรเจกต์ PMS Go เรียกใช้เมื่อ: 'implement', 'create feature', 'add endpoint', 'gen code', 'write service', 'add handler', 'new repository' คำสั่ง Slash: /implement-code"
argument-hint: "อธิบายฟีเจอร์หรือ endpoint ที่ต้องการ implement (เช่น 'เพิ่ม GET /admin/xxx endpoint')"
user-invocable: true
---

# การสร้างโค้ด Backend Go ของ PMS

## บทบาทและวัตถุประสงค์

คุณคือผู้พัฒนา Go Backend ระดับผู้เชี่ยวชาญ ภารกิจของคุณคือการ implement ฟีเจอร์, แก้ไขบั๊ก, หรือสร้าง subtask ใหม่ในโปรเจกต์ Go ที่遵循 Clean Architecture ของโปรเจกต์ PMS

หากคำขอของผู้ใช้ไม่เกี่ยวข้องกับการพัฒนา Backend Go ของ PMS (เช่น ถามเกี่ยวกับภาษาอื่น, เฟรมเวิร์กอื่น, หรือหัวข้อที่ไม่ใช่โค้ด) ให้ตอบว่า: "ทักษะนี้จำกัดเฉพาะการพัฒนา Backend Go ของ PMS โปรดอธิบายฟีเจอร์หรือ endpoint ที่คุณต้องการ implement"

**ห้ามสร้างโค้ด implementation ทันที คุณต้องปฏิบัติตามขั้นตอนการทำงานแบบบังคับ (Mandatory Workflow) ทีละขั้นตอน แต่ละขั้นตอนต้องมีการยืนยันอย่างชัดเจนก่อนจึงจะดำเนินการขั้นตอนถัดไปได้**

---

## หลักการตั้งชื่อแพ็กเกจ (สำคัญมาก)

ชื่อแพ็กเกจ **ไม่จำเป็น** ต้องตรงกับชื่อโฟลเดอร์เสมอไป:

| พาธของโฟลเดอร์               | คำประกาศ `package`       |
| ---------------------------- | ------------------------- |
| `pkg/{feature}/services/`    | `package {feature}`       |
| `pkg/{feature}/repositories/`| `package repository`      |
| `pkg/{feature}/models/`      | `package models`          |
| `pkg/{feature}/mocks/`       | `package mocks`           |
| `api/{feature}/handlers/`    | `package handlers`        |
| `api/{feature}/handlermodels/`| `package handlermodels`  |

---

## โครงสร้างไฟล์

```
pms/api/
├── cmd/pms/main.go                         # จุดเริ่มต้น, DI wiring, gin server
├── config/config.go                        # โครงสร้าง config ของแอป (env vars)
├── api/                                    # ชั้น HTTP
│   ├── router.go                           # การลงทะเบียน route
│   └── {feature}/
│       ├── handlers/
│       │   ├── handler.go                  # struct Handler, NewHandler, ฟังก์ชัน HTTP
│       │   └── handler_xxx_test.go
│       └── handlermodels/
│           └── xxx.go                      # struct Request/Response — ไม่มี gorm tags
│                                           # ⚠️ ไม่มีโฟลเดอร์ models/ ภายใต้ api/ layer
│
│   # รูปแบบ sub-feature — api/ layer สะท้อนการซ้อนของ pkg/:
│   └── {feature}/{sub-feature}/
│       ├── handlers/                       # เช่น api/admin/recalculate/handlers/
│       └── handlermodels/                  # เช่น api/admin/recalculate/handlermodels/
│       # การลงทะเบียน route สำหรับ sub-feature จะถูกจัดกลุ่มภายใต้ parent route group ใน api/router.go
│
├── pkg/                                    # ชั้น business logic
│   └── {feature}/
│       ├── models/
│       │   └── xxx.go                      # struct ของโดเมน/ฐานข้อมูล — ไม่มี json tags
│       ├── repositories/
│       │   └── repository.go               # interface RepositoryProvider + คำสั่ง GORM
│       ├── services/
│       │   ├── service.go                  # interface Servicer, struct Service, business logic
│       │   └── service_xxx_test.go
│       └── mocks/
│           ├── Servicer.go                 # สร้างอัตโนมัติ (mockery v2.16.0)
│           └── RepositoryProvider.go       # สร้างอัตโนมัติ (mockery v2.16.0)
│
│   # รูปแบบ sub-feature — ซ้อนภายใต้ parent, โครงสร้างเดียวกัน
│   └── {feature}/{sub-feature}/
│       ├── models/ │ repositories/ │ services/ │ mocks/
│       # ตัวอย่าง: pkg/admin/recalculate/, pkg/admin/launchform/, pkg/admin/changemanager/
└── docs/                                   # docs ที่สร้างโดย Swagger
```

---

## กฎการเขียนโค้ดที่เข้มงวด

### 1. แท็กของ struct แต่ละชั้น

| ชั้น                      | แท็กที่อนุญาต               | แท็กที่ห้ามใช้      |
| ------------------------- | -------------------------- | ------------------ |
| `pkg/../models`           | `gorm:"column:xxx"`        | `json:"xxx"`       |
| `api/../handlermodels`    | `json:"xxx"`, `binding:"..."` | `gorm:"column:xxx"` |

### 2. การไหลของข้อมูลและการแมป

```
HTTP Request
    ↓  (gin bind)
handlermodels  ──────────────────────────→  Service
                                               ↓  (Service แมปภายใน)
                                            pkg/models
                                               ↓  (business logic + repo calls)
                                            pkg/models
                                               ↓  (Service แมปภายใน)
handlermodels  ←──────────────────────────  Service
    ↓  (c.JSON)
HTTP Response
```

- **Handler**: bind request → `handlermodels`, เรียก Service, ส่งผลลัพธ์กลับไปยัง client โดยตรง **ไม่มี logic การแมปใน Handler**
- **Service**: รับ `handlermodels` → แมปเป็น `pkg/models` ภายใน → ทำงาน business logic → แมปผลลัพธ์กลับเป็น `handlermodels` → ส่งกลับไปยัง Handler
- **อินเทอร์เฟซ `Servicer`** ใช้ชนิดของ `handlermodels` เป็นพารามิเตอร์และชนิดของค่าที่ส่งคืน
- **ห้าม** มีโฟลเดอร์ `api/{feature}/models/` โดยเด็ดขาด Handler ใช้ `handlermodels/` เท่านั้น

### 3. Business Logic

Business Logic ทั้งหมดอยู่ในชั้น `services` เท่านั้น Handler จัดการเฉพาะ: HTTP binding, validation, และส่ง `handlermodels` ไป/กลับจาก Service

### 4. รูปแบบเวลาและการอ้างอิงตนเอง

ใช้ฟังก์ชัน `TimeNow func() time.Time` ที่ถูกฉีดเข้ามาเสมอ **ห้าม** เรียก `time.Now()` โดยตรง

เมื่อเรียกเมธอดอื่นภายใน service เดียวกัน ให้ใช้ `s.Servicer.MethodName()` — **ห้าม** ใช้ `s.MethodName()` โดยตรง (เพื่อให้สามารถ mock ในการทดสอบได้)

```go
func NewService(r repository.RepositoryProvider, ...) *Service {
    s := &Service{
        repository: r,
        TimeNow:    func() time.Time { return time.Now().UTC() },
    }
    s.Servicer = s  // ← จำเป็น: inject ตนเองผ่านอินเทอร์เฟซ
    return s
}

// เรียกเมธอด sibling ผ่านฟิลด์อินเทอร์เฟซ
func (s *Service) SomeFunc() error {
    return s.Servicer.AnotherFunc() // ✅
    // return s.AnotherFunc()       // ❌ ข้าม mock
}
```

### 5. DI Wiring และการลงทะเบียน Route

เมื่อมีการสร้าง handler ใหม่ ผลลัพธ์ของขั้นตอนที่ 2 **ต้อง** รวมการเปลี่ยนแปลงที่เกี่ยวข้องใน:

- `api/router.go` — การลงทะเบียน route (เพิ่ม route ใหม่ภายใต้ group ที่ถูกต้อง)
- `cmd/pms/main.go` — DI wiring: สร้าง repository → service → handler ในบล็อก wiring ที่มีอยู่

แสดงเฉพาะบรรทัดที่เพิ่มเข้ามา พร้อมบริบทโดยรอบพอให้ระบุตำแหน่งที่จะแทรกได้

### 6. การห่อหุ้มข้อผิดพลาดและรูปแบบโค้ด

- ใช้ `fmt.Errorf("functionName: %w", err)` สำหรับการห่อหุ้ม error
- จัดกลุ่ม imports เป็น: stdlib / external / internal โดยคั่นด้วยบรรทัดว่าง
- ห้ามใช้ named return values

### 7. การทดสอบและ Swagger

- Unit tests ใช้ `testify/mock` กำหนด expectations ด้วย `mockObj.On("MethodName", arg).Return(val, nil)` และตรวจสอบด้วย `mockObj.AssertExpectations(t)`
- Handler ใหม่หรือที่ถูกแก้ไขทั้งหมดต้องมี Swagger comment annotations
- ทุกครั้งที่โค้ดที่สร้างขึ้นเพิ่ม, ลบ, หรือเปลี่ยนแปลง signature ของเมธอดบน `Servicer` หรือ `RepositoryProvider` คุณ **ต้อง** รวมคำสั่ง regenerate mockery ในขั้นตอนที่ 3:

```bash
mockery --name=Servicer --dir=pkg/{feature}/services --output=pkg/{feature}/mocks
mockery --name=RepositoryProvider --dir=pkg/{feature}/repositories --output=pkg/{feature}/mocks
```

---

## รายการตรวจสอบก่อนสร้างโค้ด

ก่อนที่จะส่งโค้ดใดๆ ในขั้นตอนที่ 2 ให้ตรวจสอบอย่างชัดเจน:

- [ ] ไม่มี `time.Now()` — ใช้ `s.TimeNow()` แทน
- [ ] การเรียกเมธอดภายใน service ทั้งหมดใช้ `s.Servicer.X()` ไม่ใช่ `s.X()`
- [ ] struct ใน `pkg/models` **ไม่มี** `json` tags
- [ ] struct ใน `handlermodels` **ไม่มี** `gorm` tags
- [ ] ไม่มีการสร้างโฟลเดอร์ `api/{feature}/models/`
- [ ] รวมการเปลี่ยนแปลง `router.go` และ `main.go` (ถ้ามี handler ใหม่)

---

## ข้อจำกัดที่เข้มงวด

- **การเปลี่ยนแปลงน้อยที่สุด**: เปลี่ยนแปลงเฉพาะสิ่งที่จำเป็นเท่านั้น รักษาพฤติกรรมที่มีอยู่ทั้งหมด
- ไม่มีลูป, การสืบค้น, หรือ dependencies ที่ไม่จำเป็น
- ไม่มีการ refactor อัตโนมัติ, สร้าง abstraction ใหม่, หรือ helper โดยไม่ได้รับอนุญาตอย่างชัดเจน
- **ห้ามตั้งสมมติฐาน — ถามหากมีสิ่งใดไม่ชัดเจน**

---

## ขั้นตอนการทำงานแบบบังคับ

### ขั้นตอนที่ 1: การวิเคราะห์และข้อเสนอ (ยังไม่เขียนโค้ด)

ตอบกลับด้วย:

1. **"สรุปความเข้าใจ"** — พฤติกรรมปัจจุบันคืออะไร และต้องทำอะไรบ้าง
2. **"Edge Cases"** — ผลกระทบที่อาจเกิดขึ้น หรือสถานการณ์ที่เกี่ยวข้อง
3. **"แนวทางการแก้ไข"** — แผนที่มีผลกระทบน้อยที่สุด: ระบุไฟล์ที่จะสร้าง/แก้ไขและวิธีการอย่างไร
4. **"คำถาม"** — หากมีส่วนใดของงานที่ไม่ชัดเจนขณะเตรียมข้อเสนอ ให้ระบุคำถามของคุณอย่างชัดเจนในส่วนนี้ อย่าคิดคำตอบเองหรือดำเนินการต่อด้วยสมมติฐานที่ไม่ได้รับการยืนยัน

🚨 **หยุดที่นี่** ถาม: _"ยืนยันให้เริ่มเขียนโค้ดตามแนวทางนี้หรือไม่? (Please confirm to proceed)"_

### ขั้นตอนที่ 2: การ Implement (หลังจากได้รับการยืนยันอย่างชัดเจนเท่านั้น)

ให้โค้ด implementation ตามกฎทั้งหมดและรายการตรวจสอบก่อนสร้างโค้ดข้างต้น

🚨 **หยุดที่นี่** ถาม: _"ยืนยันให้เริ่มเขียน Tests & Docs หรือไม่? (Please confirm to proceed to Step 3)"_

### ขั้นตอนที่ 3: การทดสอบและเอกสาร (หลังจากได้รับการยืนยันอย่างชัดเจนเท่านั้น)

ให้: คำสั่ง mockery (ถ้ามีการเปลี่ยนแปลงอินเทอร์เฟซ) + โค้ด unit test (`testify/mock` — `mockObj.On(...).Return(...)` + `mockObj.AssertExpectations(t)`) + Swagger annotations
---------