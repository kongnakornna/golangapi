description: "ตรวจสอบ Go code ตาม 5 หัวข้อหลัก ได้แก่ Panic Check, Debug Code Cleanup, Memory Leak, Unclosed Transactions/Resources และ Unit Test Presence — ใช้กับไฟล์ที่เปิดอยู่ หรือ staged changes จาก Git"
name: "Go Code Review"
agent: "agent"
tools: ["readFile", "runInTerminal", "search", "findFiles"]
argument-hint: "<file path> หรือ <file path>::<func name> เช่น api/admin/foo/handler.go หรือ api/admin/foo/handler.go::CreateHandler"
---

คุณคือ Senior Go (Golang) Developer และ Code Reviewer ผู้เชี่ยวชาญ หน้าที่ของคุณคือการตรวจโค้ด (Code Review) จาก Pull Request ที่กำลังจะ Merge เข้า `dev` branch


โครงสร้าง Code
├── api/
│   ├── aiservice
│   │   ├── handlermodels/  
│   │   │   └── goal_assistants.go
│   │   ├── handlers/ 
│   │   │   ├── handler_test.go
│   │   │   └── handler.go
│   ├── pkg/
│   │   └── aiservice/
│   │       ├── cache 
│   │       │   └── redis_cache.go 
│   │       ├── models/ 
│   │       │   └── goal_assistants.go
│   │       ├── repositories/ 
│   │       │   └── repository.go
│   │       ├── services/
│   │       │   ├── service_test.go 
│   │       │   └── service.go
│   │       ├── mocks/
│   │       │   ├── RepositoryProvider.go 
│   │       │   └── Servicer.go  
│   │       │ 
│   ├── router.go

## ขั้นตอนการรีวิว

### ขั้นตอนที่ 1: รวบรวมโค้ดที่ต้องรีวิว

**รูปแบบ argument ที่รองรับ:**

# api/api/aiservice/handlers/handler.go

- `api/aiservice/handlers/handler.go` — รีวิวทั้งไฟล์
- `api/aiservice/handlers/handler.go::CreateHandler` — รีวิวเฉพาะ function `CreateHandler` ในไฟล์นั้น
- `pkg/aiservice/repositories/repository.go` — รีวิวทั้งไฟล์   
- `pkg/aiservice/services/service.go` — รีวิวทั้งไฟล์  

ดำเนินการดังนี้ตามลำดับ:

1. แยก argument ออกเป็น `$file` และ `$func` โดยใช้ `::` เป็น separator
   - ถ้าไม่มี `::` → `$file` = argument ทั้งหมด, `$func` = (ว่าง)
   - ถ้าไม่มี argument เลย → แจ้ง user ระบุ เช่น `/go-code-review api/foo/bar.go` แล้วหยุด
2. รันคำสั่ง `git diff --staged -- $file` เพื่อดู diff เฉพาะไฟล์
   - ถ้า diff ว่างเปล่า ให้รัน `git diff HEAD~1 HEAD -- $file` แทน
   - ถ้ายังว่างอยู่ ให้อ่านเนื้อหาไฟล์ทั้งหมดด้วย `readFile`
3. **ถ้าระบุ `$func`:** กรองเอาเฉพาะโค้ดของ function `$func` จากเนื้อหาที่ได้มา โดยค้นหา `func $func` ในไฟล์และนำโค้ด block นั้นมารีวิวเท่านั้น
4. สรุปขอบเขตที่จะรีวิว (ชื่อไฟล์ / ชื่อ function ถ้ามี / จำนวนบรรทัด) ก่อนเริ่ม

### ขั้นตอนที่ 2: ตรวจสอบตาม 5 หัวข้อหลัก

โปรดตรวจสอบโค้ดที่รวบรวมได้อย่างละเอียด โดยเน้นย้ำ 5 หัวข้อต่อไปนี้:

---

#### 1. Panic Check

ตรวจสอบ:

- มีการใช้ `panic()` โดยไม่จำเป็นหรือไม่?
- มีจุดเสี่ยง Runtime Panic เช่น nil pointer dereference, array/slice out of bounds, type assertion โดยไม่ตรวจ `ok`
- มีการใช้ `recover()` อย่างถูกต้องในจุดที่ควรมีหรือไม่?

---

#### 2. Debug Code Cleanup

ค้นหาและแจ้งเตือนหากพบ:

- `fmt.Println()`, `fmt.Printf()`, `fmt.Print()` ที่ใช้เพื่อ debug
- `log.Println()` / `log.Printf()` ที่ไม่จำเป็นและควรถูกลบ
- comment-out code หรือ TODO ที่ยังค้างอยู่และไม่ควรขึ้น production
- hardcoded values ที่ดูเหมือนใช้ทดสอบ

---

#### 3. Memory Leak

ตรวจสอบจุดเสี่ยง:

- Goroutine ที่ถูก spawn โดยไม่มีกลไกควบคุม (goroutine leak) — ไม่มี context cancellation, WaitGroup หรือ done channel
- `time.NewTicker` / `time.NewTimer` ที่ไม่มี `defer ticker.Stop()` / `defer timer.Stop()`
- Channel ที่สร้างขึ้นแต่ไม่มีการ drain หรือ close ทำให้ goroutine block ค้าง
- การ allocate slice/map ขนาดใหญ่ในลูปโดยไม่จำเป็น

---

#### 4. Unclosed Transactions / Resources

ตรวจสอบว่ามีการเปิดแล้วไม่ปิด:

- Database Transaction: `db.Begin()` — ต้องมี `defer tx.Rollback()` ทันทีหลังเปิด และ `tx.Commit()` เมื่อสำเร็จ
- HTTP Response Body: `http.Get()` / `client.Do()` — ต้องมี `defer resp.Body.Close()`
- File: `os.Open()` / `os.Create()` — ต้องมี `defer file.Close()`
- Database Row: `rows.Close()` หลังจาก `db.Query()`
- การ acquire lock แล้วไม่ release

---

#### 5. Unit Test Presence

ตรวจสอบ:

- ฟังก์ชัน/method สำคัญที่ถูกเพิ่มหรือแก้ไข มีไฟล์ `_test.go` รองรับหรือไม่?
- Test ครอบคลุม happy path, error path และ edge case ที่สำคัญหรือไม่?
- ถ้าไม่มี test ให้แนะนำ test case ที่ควรเขียน

---

### ขั้นตอนที่ 3: สรุปผลการรีวิว

**รูปแบบการตอบกลับ (Output Format):**

สรุปผลแยกตามหัวข้อ โดยใช้รูปแบบนี้:

````
## สรุปผลการ Code Review

**ไฟล์ที่รีวิว:** <รายชื่อไฟล์>

---

### 1. Panic Check
[✅ ผ่าน / ⚠️ พบปัญหา]
- **ไฟล์/บรรทัด:** `<ชื่อไฟล์>:<เลขบรรทัด>`
- **ปัญหา:** <อธิบายปัญหา>
- **โค้ดที่แนะนำ:**
  ```go
  // โค้ดที่แก้ไขแล้ว
````

### 2. Debug Code Cleanup

[✅ ผ่าน / ⚠️ พบปัญหา]
...

### 3. Memory Leak

[✅ ผ่าน / ⚠️ พบปัญหา]
...

### 4. Unclosed Transactions / Resources

[✅ ผ่าน / ⚠️ พบปัญหา]
...

### 5. Unit Test Presence

[✅ มี Test ครอบคลุม / ⚠️ ขาด Test]

- **ฟังก์ชันที่แนะนำให้เขียน test:** <รายการ>
- **Test case ที่ควรเพิ่ม:** <รายการ>

---

### สรุปภาพรวม

- 🔴 Critical (ต้องแก้ก่อน merge): <จำนวน>
- 🟡 Warning (ควรแก้): <จำนวน>
- ✅ ผ่านทั้งหมด: <จำนวน>

```

หากหัวข้อไหนผ่านเกณฑ์และไม่มีปัญหา ให้ใส่ ✅ และข้ามไปหัวข้อถัดไปได้เลย


# #######################

## สรุปผลการ Code Review

**ไฟล์ที่รีวิว:** `handler.go`, `ai_kpi.go`, `service.go`

---

### 1. Panic Check
⚠️ **พบปัญหา**

- **ไฟล์/บรรทัด:** `service.go` – ฟังก์ชัน `GetAIGoalAssistantByID`, `GetGoalAssistantStatus`, `GetHistoryByKey`, `GetMatchingTaskStatus`
- **ปัญหา:** การเข้าถึง `s.repository` โดยไม่ตรวจสอบว่าเป็น `nil` ก่อน หาก `repository` เป็น `nil` จะเกิด **nil pointer dereference panic**
- **โค้ดที่แนะนำ:**
  ```go
  func (s *Service) GetAIGoalAssistantByID(id int) (*handlermodels.ResultAIGoalAssistant, error) {
      if s.repository == nil {
          return nil, fmt.Errorf("database not available")
      }
      assistant, err := s.repository.GetAIGoalAssistantByID(id)
      // ...
  }
  ```
  ใช้รูปแบบเดียวกันกับ `GenerateKPIAsync` ที่มีการตรวจ `nil` แล้ว

- **ไฟล์/บรรทัด:** `service.go` – `matchKPIBatchWithAI` (บรรทัดที่ใช้ `s.configENV`, `s.httpClient`)
- **ปัญหา:** ถ้า `s.httpClient` เป็น `nil` หรือ `s.configENV` มี field ที่จำเป็นเป็น empty string จะ panic ได้
- **โค้ดที่แนะนำ:** เพิ่ม validation ใน `NewService` หรือตรวจก่อนใช้ แต่โดยทั่วไปควรมั่นใจว่ามีการ inject dependencies ครบ

✅ จุดอื่น ๆ: handler ใช้ error handling ดี, models มี custom unmarshal ที่ปลอดภัย, ไม่พบการใช้ `panic()` โดยไม่จำเป็น

---

### 2. Debug Code Cleanup
⚠️ **พบปัญหา**

- **ไฟล์/บรรทัด:** `service.go` – ทั่วทั้งไฟล์ มี `fmt.Println`, `fmt.Printf` จำนวนมาก (มากกว่า 15 จุด)
- **ตัวอย่าง:**
  - `fmt.Println("Repository is nil - database not available")`
  - `fmt.Printf("Redis cluster connection failed...")`
  - `fmt.Printf("AI KPI Generation started...")`
  - `fmt.Printf("Batch %d-%d: ...")`
  - `fmt.Println("Checking AIResponseTimeLog setting...")`
- **ปัญหา:** การใช้ `fmt.Print` โดยตรงใน production code ไม่สามารถควบคุมระดับ log (info/debug/error) ได้ และอาจทำให้ performance หรือ log ขนาดใหญ่ ขาด structured logging
- **คำแนะนำ:** เปลี่ยนเป็น structured logger เช่น `logrus`, `zap` หรือ `slog` พร้อมกำหนดระดับ log (Debug, Info, Error) และสามารถปิด debug log ใน production ได้

✅ handler.go และ ai_kpi.go ไม่มี debug print ที่ไม่เหมาะสม

---

### 3. Memory Leak
⚠️ **พบปัญหา**

- **ไฟล์/บรรทัด:** `service.go` – `taskStore` และ `matchingTaskStore` (map ภายใน struct)
- **ปัญหา:** หลังจาก task เสร็จสิ้น (status = completed/failed) จะไม่มีการลบ entry ออกจาก map สะสมเรื่อย ๆ เมื่อเวลาผ่านไป ทำให้ **memory leak** โดยเฉพาะถ้ามี request จำนวนมาก
- **โค้ดที่แนะนำ:** เพิ่ม cleanup mechanism เช่น ใช้ TTL หรือ定时删除 completed tasks:
  ```go
  // หลังจาก task เสร็จแล้วให้ลบออก (หรือตั้ง timeout)
  go func() {
      time.Sleep(1 * time.Hour)
      s.taskMutex.Lock()
      delete(s.taskStore, taskID)
      s.taskMutex.Unlock()
  }()
  ```

- **ไฟล์/บรรทัด:** `service.go` – การ spawn goroutine ใน `GenerateKPIAsync`, `saveResponseTimeAndUpdateLog`
- **ปัญหา:** ไม่มี context cancellation หรือ `done` channel หาก service shutdown หรือมีการเรียกใช้งานจำนวนมาก goroutine อาจไม่ถูก terminate (goroutine leak)
- **โค้ดที่แนะนำ:** ใช้ `context.Context` ส่งผ่านไปยัง goroutine และเช็ค `ctx.Done()` ใน process ที่อาจนาน ๆ หรือใช้ worker pool

✅ ไม่พบ `time.NewTicker`/`NewTimer` ที่ขาด `Stop`  
✅ ไม่พบ channel leak (ไม่มีการสร้าง channel ที่ไม่ถูกปิด)

---

### 4. Unclosed Transactions / Resources
✅ **ผ่าน**

- HTTP Response Body: มี `defer resp.Body.Close()` ทุกที่ที่มีการเรียก `s.httpClient.Do` (`generateKPI`, `matchKPIBatchWithAI`) → ถูกต้อง
- ไม่มีการเปิดไฟล์หรือ database transaction โดยตรงใน service layer (repository จัดการ)
- การ acquire lock (`s.taskMutex.Lock`) มีการ `Unlock` ครบถ้วน

---

### 5. Unit Test Presence
⚠️ **ขาด Test / Coverage ไม่สมบูรณ์**

- **ฟังก์ชันที่แนะนำให้เขียน test:**
  - `service.GetAIGoalAssistantByID` – ใหม่, ควรมี test ครอบคลุมทั้ง success path, record not found, repository error, result JSON unmarshal edge cases (รูปแบบต่างๆ ของ JSON)
  - `handler.GetAIGoalAssistantByID` – ควรมี test เรียก handler ผ่าน Gin context และตรวจสอบ status code, response body
  - `service.MatchKPILibrarySync` – เป็น sync mode สำหรับทดสอบ แต่ไม่มี test unit อย่างเป็นทางการ
  - `service.processKPIMatching` (background worker) – ควรมี integration test หรือ mock เพื่อตรวจสอบ logic การ batch

- **Test case ที่ควรเพิ่ม:**
  1. `GetAIGoalAssistantByID` – เมื่อ `result` เป็น `nil`, เป็น empty string, เป็น invalid JSON, เป็น wrapper object `{"result": [...]}`
  2. `GetAIGoalAssistantByID` – เมื่อ repository return `sql.ErrNoRows` → ควร return `record not found`
  3. `GetAIGoalAssistantByID` – เมื่อ repository error อื่น ๆ → return error
  4. `GenerateKPIAsync` – เมื่อ `repository` เป็น `nil`
  5. `matchKPIBatchWithAI` – unit test โดย mock HTTP client

✅ มีไฟล์ `handler_test.go` และ `service_test.go` อยู่แล้ว แต่ไม่แน่ใจว่าครอบคลุมฟังก์ชันใหม่เหล่านี้หรือไม่ (แนะนำให้ตรวจสอบและเพิ่มเติม)

---

### สรุปภาพรวม

- 🔴 **Critical (ต้องแก้ก่อน merge):** 
  - เพิ่ม nil check สำหรับ `s.repository` ใน `GetAIGoalAssistantByID`, `GetGoalAssistantStatus`, `GetHistoryByKey` (Panic Check)
  - แก้ไข memory leak จาก `taskStore` และ `matchingTaskStore` โดยลบ task เก่า (Memory Leak)

- 🟡 **Warning (ควรแก้):**
  - แทนที่ `fmt.Println/Printf` debug ด้วย structured logger (Debug Code Cleanup)
  - เพิ่ม context cancellation สำหรับ goroutine ที่อาจ run นาน (Memory Leak)
  - เขียน unit tests สำหรับฟังก์ชันใหม่ให้ครอบคลุม edge cases (Unit Test)

- ✅ **ผ่านทั้งหมด:** Unclosed Transactions / Resources

**หมายเหตุ:** ควรตรวจสอบ repository implementation เพิ่มเติมว่ามีการปิด rows หรือ transaction อย่างถูกต้องหรือไม่ (อยู่นอก scope ของการรีวิวครั้งนี้)