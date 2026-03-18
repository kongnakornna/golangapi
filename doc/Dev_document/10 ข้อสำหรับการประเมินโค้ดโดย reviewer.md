# คู่มือ  Code Review Criteria
### เกณฑ์สำคัญ 10 ข้อสำหรับการประเมินโค้ดโดย reviewer

**เกณฑ์สำคัญ 10 ข้อสำหรับการประเมินโค้ด (Code Review Criteria)**
เพื่อช่วยให้ Reviewer ตรวจงานได้อย่างมีทิศทาง ลดความขัดแย้ง และยกระดับคุณภาพซอฟต์แวร์ นี่คือ Checklist 10 ข้อที่ครอบคลุมทั้ง Functionality, Quality และ Testing ครับ:

***

### หมวดที่ 1: ความถูกต้องและการใช้งาน (Does it work?)

**1. ตรงตาม Requirement (Correctness):**

* โค้ดทำงานถูกต้องตาม Ticket/User Story หรือไม่?
* Logic การคำนวณหรือ Flow การทำงานถูกต้องตาม Business Rule ไหม?
* *คำถาม:* "ถ้า User ทำตาม Step 1-2-3 ผลลัพธ์ออกมาถูกเป๊ะไหม?"

**2. รองรับ Edge Cases (Robustness):**

* ทดสอบกรณี "ข้อมูลแปลกๆ" หรือยัง? (เช่น ข้อมูลเป็น Null, Array ว่าง, User กรอก Emoji, เน็ตหลุดกลางทาง)
* มีการจัดการ Error (Error Handling) ที่เหมาะสมหรือไม่? ไม่ใช่แค่ `try-catch` ทิ้งไว้เฉยๆ

**3. ความปลอดภัย (Security):**

* มีการตรวจสอบข้อมูลนำเข้า (Input Validation) ไหม?
* มีช่องโหว่พื้นฐานหรือไม่? (เช่น SQL Injection, XSS, หรือเผลอ Hardcode Password/API Key ลงไปในโค้ด)

***

### หมวดที่ 2: คุณภาพโค้ด (Is it clean?)

**4. อ่านง่ายและสื่อความหมาย (Readability \& Naming):**

* ชื่อตัวแปร ฟังก์ชัน และคลาส สื่อความหมายชัดเจนหรือไม่? (เช่น `d` ❌ vs `daysSinceLastLogin` ✅)
* โครงสร้างโค้ดซับซ้อนเกินไปไหม? (Cyclomatic Complexity) ถ้าอ่านแล้วต้องขมวดคิ้วเกิน 3 วิ แสดงว่าควรแก้

**5. ไม่ทำงานซ้ำซ้อน (DRY - Don't Repeat Yourself):**

* มีการ Copy-Paste Logic เดิมไปแปะหลายที่ไหม?
* ถ้ามี ควรยุบรวมเป็น Function กลาง หรือ Component ที่ใช้ร่วมกันได้

**6. ประสิทธิภาพ (Performance):**

* มีการ Loop ซ้อน Loop (O(n^2)) โดยไม่จำเป็นไหม?
* มีการ Query Database ใน Loop (N+1 Problem) หรือไม่?
* มีการโหลดข้อมูลมาเยอะเกินความจำเป็นไหม? (เช่น `SELECT *` แต่ใช้แค่ 2 fields)

**7. สไตล์และมาตรฐาน (Coding Standard):**

* การจัด Format (เว้นวรรค, ย่อหน้า) ตรงตามมาตรฐานทีมหรือ Linter ไหม?
* โครงสร้าง Folder/File ถูกต้องตาม Architecture ของโปรเจกต์ไหม?

***

### หมวดที่ 3: การทดสอบและการดูแลรักษา (Can we maintain it?)

**8. การทดสอบ (Test Coverage \& Quality):**

* มี Unit Test ครอบคลุม Logic ใหม่หรือไม่?
* Test เขียนตามหลัก Arrange-Act-Assert และ Isolation (Mock dependency) หรือไม่?
* Test เคส Unhappy Path (กรณี Error) ด้วยหรือเปล่า?

**9. ไม่กระทบของเดิม (No Regression):**

* การแก้นี้ไปทำให้ฟีเจอร์เก่าที่เคยดีอยู่...พังไหม?
* ควรเช็คว่ามีการแก้ไขไฟล์ที่ไม่เกี่ยวข้องโดยไม่ตั้งใจหรือไม่

**10. เอกสารประกอบ (Documentation):**

* ถ้ามีการแก้ API มีการอัปเดต Swagger/Postman หรือยัง?
* ถ้า Logic ซับซ้อนมาก มี Comment อธิบาย "Why" (ทำไมถึงเขียนแบบนี้) ไว้ไหม?

***

### 💡 เคล็ดลับสำหรับ Reviewer

* **Be Constructive:** วิจารณ์ที่ "โค้ด" ไม่ใช่ "คน" (เช่น "ตรงนี้อาจทำให้ช้า" แทน "ทำไมเขียนแบบนี้")
* **Nitpicks:** เรื่องเล็กน้อย (เช่น ลืมลบ console.log) ให้ระบุว่าเป็น "Nitpick" (แก้ก็ดี ไม่แก้ก็ได้) เพื่อไม่ให้ผู้ถูกตรวจรู้สึกกดดัน
* **Approve with Comments:** ถ้ามีแก้เล็กน้อย ให้ Approve ไปเลยแต่ฝากแก้ด้วย เพื่อไม่ให้งานสะดุด
<span style="display:none">[^1][^2][^3][^4][^5][^6][^7][^8]</span>

<div align="center">⁂</div>

[^1]: https://www.reddit.com/r/ExperiencedDevs/comments/1k53zol/how_to_give_code_reviews_without_offending_other/

[^2]: https://www.softnova.co/post/การสื่อสารอย่างมีประสิทธิภาพกับทีมเทคโนโลยี-เพื่อความสำเร็จของโปรเจกต์

[^3]: https://sennalabs.com/blog/code-review

[^4]: https://www.facebook.com/groups/ThaiPGAssociateSociety/posts/2903426459868604/

[^5]: https://developers.ascendcorp.com/review-code-เมื่อโค้ดจบแต่คนยังไม่จบ-บทบาทสำคัญในการเป็น-software-developer-aafbbc456d80

[^6]: https://www.techstarthailand.com/blog/detail/How-to-properly-give-a-code-review/933

[^7]: https://anyflip.com/ipbly/jekh/basic

[^8]: https://www.catc.or.th/internalaudit/file/หลักเกณฑ์ประเมินCBE_%E0%B8%A3%E0%B8%AA._(%E0%B8%89%E0%B8%9A%E0%B8%B1%E0%B8%9A%E0%B8%9B%E0%B8%A3%E0%B8%B1%E0%B8%9A%E0%B8%9B%E0%B8%A3%E0%B8%B8%E0%B8%87%20%E0%B8%9B%E0%B8%B5%202567).pdf

