# จะตั้งเกณฑ์อะไรให้ reviewer ตรวจโค้ดอย่างมีประสิทธิภาพ

การกำหนด **"เกณฑ์ (Criteria)"** ให้ Reviewer ใช้ตรวจงาน เป็นสิ่งสำคัญมาก เพื่อไม่ให้การตรวจกลายเป็นแค่เรื่อง "รสนิยมส่วนตัว" นี่คือ Checklist ที่แบ่งตามความสำคัญ (Must/Should/Could) ที่คุณสามารถนำไปแปะไว้ใน Template ของ Pull Request ได้เลยครับ

### 1. Functionality (ใช้งานได้จริงไหม?) — **[MUST]**

*ด่านแรกที่สำคัญที่สุด ถ้าไม่ผ่านข้อนี้ ห้าม Approve เด็ดขาด*

* [ ] **Requirements Met:** โค้ดทำงานตรงตาม Ticket หรือ User Story หรือไม่?
* [ ] **Edge Cases:** รองรับกรณีแปลกๆ หรือยัง? (เช่น ค่า NULL, ข้อมูลว่าง, เน็ตหลุด, User กรอกภาษาต่างดาว)
* [ ] **No Regression:** การแก้นี้ไปทำลายฟีเจอร์เก่าที่เคยทำงานได้หรือไม่?


### 2. Unit Test Quality (เทสมีคุณภาพไหม?) — **[MUST]**

*อ้างอิงจากบทความ EPT: ใช้หลักการ 5 ข้อมาจับ*

* [ ] **Arrange-Act-Assert:** เขียนเทสชัดเจนไหม? (เตรียมของ -> ทำ -> ตรวจผล)
* [ ] **Isolation:** เทสนี้ "แยกขาด" จริงไหม? (ต้องไม่มีการต่อ Database จริง, ไม่เรียก API ภายนอกจริง ต้องใช้ Mock เท่านั้น)
* [ ] **Coverage:** ครอบคลุมทั้งเคส "ปกติ" (Happy Path) และเคส "Error" (Unhappy Path) หรือยัง?
* [ ] **Readability:** ชื่อ Test Function อ่านแล้วรู้เรื่องทันทีไหมว่าเทสอะไร? (เช่น `should_throw_error_when_password_too_short` ไม่ใช่ `test_fail_1`)


### 3. Code Quality \& Readability (อ่านรู้เรื่องไหม?) — **[SHOULD]**

*เน้นความยั่งยืน (Maintainability) เพื่อให้คนอื่นมาแก้ต่อได้*

* [ ] **Naming:** ชื่อตัวแปร/ฟังก์ชัน สื่อความหมายชัดเจน ไม่ใช้ชื่อย่อที่รู้กันเอง (เช่น `x`, `data`, `temp` ❌ -> `userList`, `totalPrice` ✅)
* [ ] **Complexity:** ฟังก์ชันยาวเกินไปไหม? (ถ้าเกิน 20-30 บรรทัด ควรแยกฟังก์ชัน)
* [ ] **DRY (Don't Repeat Yourself):** มีการ Copy-Paste โค้ดเดิมซ้ำๆ ไหม? (ถ้ามีควรยุบเป็นฟังก์ชันกลาง)
* [ ] **Comments:** คอมเมนต์อธิบาย "ทำไม" (Why) ไม่ใช่อธิบายว่า "ทำอะไร" (What) (โค้ดที่ดีควรอธิบายตัวเองได้อยู่แล้ว)


### 4. Security \& Performance (ปลอดภัยและเร็วไหม?) — **[SHOULD]**

* [ ] **SQL Injection:** มีการต่อ String ใน SQL ตรงๆ ไหม? (ต้องใช้ ORM หรือ Parameterized Query)
* [ ] **Sensitive Data:** เผลอ Hardcode รหัสผ่านหรือ API Key ลงไปในโค้ดหรือเปล่า? ☠️
* [ ] **N+1 Problem:** มีการ Loop เรียก Database ทีละ row ไหม? (ถ้ามีควรแก้เป็น Query ทีเดียว)


### 5. Style \& Housekeeping (เรื่องจุกจิก) — **[COULD]**

*เรื่องพวกนี้ควรให้ "ระบบอัตโนมัติ (Linter/Prettier)" จัดการแทนคน*

* [ ] **Formatting:** เว้นวรรค, ปีกกา, ย่อหน้า ตรงตามมาตรฐานทีมไหม?
* [ ] **Unused Code:** มีตัวแปรที่ประกาศไว้แต่ไม่ได้ใช้ หรือ `console.log` ที่ลืมลบไหม?

***

### ตัวอย่าง Comment ที่ดี vs ไม่ดี (สำหรับ Reviewer)

| ❌ ไม่ดี (Vague/Rude) | ✅ ดี (Specific/Constructive) |
| :-- | :-- |
| "โค้ดแย่มาก ไปแก้มาใหม่" | "ฟังก์ชันนี้ดูซับซ้อนไปนิด ลองแยก Logic ส่วนคำนวณภาษีออกมาเป็นอีกฟังก์ชันดีไหมครับ? จะได้เทสง่ายขึ้น" |
| "ชื่อตัวแปรงง" | "ตัวแปร `d` สื่อความหมายไม่ชัดเจน แนะนำให้เปลี่ยนเป็น `transactionDate` เพื่อให้อ่านง่ายขึ้นครับ" |
| "ทำไมไม่เขียนเทส?" | "รบกวนเพิ่ม Unit Test สำหรับเคสที่ 'User ไม่ได้ Login' ด้วยครับ เพื่อให้ Coverage ครอบคลุมตามเกณฑ์" |

### คำแนะนำสำหรับการนำไปใช้

ให้สร้าง **Pull Request Template** ใน GitLab/GitHub แล้วใส่ Checklist นี้ลงไป เพื่อให้ Developer ต้องติ๊กยืนยันก่อนกดส่ง PR ครับ

```markdown
## Reviewer Checklist
- [ ] Functionality: ทำงานถูกต้องตาม Requirement และรองรับ Edge Cases
- [ ] Testing: มี Unit Test ครอบคลุม (Arrange, Act, Assert, Isolation)
- [ ] Readability: ตั้งชื่อตัวแปรชัดเจน ไม่ซับซ้อนเกินไป
- [ ] Security: ไม่มี Sensitive Data หรือช่องโหว่
```
