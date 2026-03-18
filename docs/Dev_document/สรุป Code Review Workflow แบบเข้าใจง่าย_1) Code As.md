
# สรุป Code Review Workflow แบบเข้าใจง่าย

1) Code Assistant — เขียนโค้ดด้วยเครื่องมือช่วย
• ใช้ AI/Code Assistant ช่วยเขียนโค้ดให้เร็วขึ้น
• ตรวจสไตล์โค้ดเบื้องต้นก่อนส่งขึ้นระบบ
เริ่มต้นให้โค้ดพร้อมตรวจ
⸻
2) Pull Request — ส่งโค้ดให้ทีมตรวจ
• เปิด Pull Request พร้อมคำอธิบาย
• แนบรายการเปลี่ยนแปลง, Issue ที่เกี่ยวข้อง
• Reviewer สามารถเข้ามาตรวจได้ทันที
เป็นขั้นตอนขอให้ทีมช่วยตรวจ
⸻
3) CI Pipeline — ตรวจอัตโนมัติ
• ระบบรัน Test / Lint / Build อัตโนมัติ
• ถ้าไม่ผ่าน จะต้องแก้ไขก่อนเข้าสู่ Code Review
ป้องกันโค้ดเสียตั้งแต่ต้นทาง
⸻
4) Code Review — ทีมตรวจโค้ด
• Reviewer ตรวจคุณภาพโค้ด ความถูกต้อง ความเหมาะสม
• อาจมีข้อเสนอแนะหรือให้แก้ไข
ให้แน่ใจว่าคุณภาพได้มาตรฐานทีม
⸻
5) Code Update — แก้ไขตามรีวิว
• ปรับตามคำแนะนำ
• Push โค้ดขึ้นมาใหม่ให้ระบบ rerun
รอบแก้ไขเพื่อให้โค้ดสมบูรณ์ขึ้น
⸻
6) Deploy — ปล่อยสู่ระบบ
• เมื่อรีวิวผ่านและ CI ผ่าน
• โค้ดถูก merge และนำไป deploy
ก้าวสู่ production อย่างมั่นใจ
⸻
เหมาะกับใคร?
• Developer ทุกระดับ
• ทีมที่ต้องการ workflow ที่ชัดเจน
• ผู้เริ่มต้นที่อยากเรียนรู้ขั้นตอน CI/CD + Code Review
• Tech Lead ที่ต้องการอธิบาย Flow ให้ทีมใหม่เข้าใจเร็ว

นี่คือสรุป **Code Review Workflow แบบเข้าใจง่าย** สำหรับทีมที่ต้องการกระบวนการที่มีคุณภาพและรวดเร็ว ตั้งแต่เริ่มเขียนจนถึงขึ้นระบบ โดยเน้นความชัดเจนในแต่ละขั้นตอน:

### 1. Code Assistant — เขียนโค้ดด้วยตัวช่วย

* **หน้าที่:** Developer
* **การกระทำ:**
    * ใช้ AI หรือ Code Assistant ช่วยร่างโค้ดและตรวจสอบ Syntax เบื้องต้น
    * รัน Unit Test (local) ด้วยตัวเองให้ผ่านทุกข้อ
    * จัด Format โค้ดให้สวยงามตามมาตรฐานทีม (Linting)
* **เป้าหมาย:** ส่งโค้ดที่ "สะอาด" และ "ทำงานได้" เข้าสู่ระบบ ลดภาระคนตรวจ


### 2. Pull Request (PR) — ส่งโค้ดให้ทีมตรวจ

* **หน้าที่:** Developer
* **การกระทำ:**
    * สร้าง PR/MR เข้า Branch หลัก (เช่น `develop`)
    * **สำคัญ:** เขียนคำอธิบาย PR ให้ชัดเจน (ทำอะไร? เพื่อแก้ Ticket ไหน? มีผลกระทบอะไร?)
    * แนบรูปภาพหรือผลเทสประกอบถ้ามี
* **เป้าหมาย:** แจ้งทีมว่า "งานเสร็จแล้ว ช่วยมาดูหน่อย"


### 3. CI Pipeline — ตรวจอัตโนมัติ (ด่านหน้า)

* **หน้าที่:** ระบบอัตโนมัติ (System)
* **การกระทำ:**
    * ทันทีที่เปิด PR ระบบจะรัน Test, Lint, และ Build Docker Image
    * เช็ค Code Coverage (ต้องผ่านเกณฑ์ที่ตั้งไว้)
    * *ถ้าไม่ผ่าน:* ระบบจะ Block ไม่ให้ Merge และแจ้งเตือน Developer ให้ไปแก้ก่อน
* **เป้าหมาย:** คัดกรอง Error พื้นฐานออกไป ไม่ให้เสียเวลาคนตรวจ


### 4. Code Review — ทีมตรวจคุณภาพ

* **หน้าที่:** Reviewer (Senior/Lead/Peer)
* **การกระทำ:**
    * อ่าน Logic ว่าถูกต้องและปลอดภัยไหม
    * เช็คความอ่านง่าย (Readability) และการตั้งชื่อตัวแปร
    * ดูว่า Unit Test ครอบคลุมและแยกส่วน (Isolation) จริงหรือไม่
    * ให้ Comment แนะนำจุดที่ควรปรับ (เน้นสร้างสรรค์ ไม่ใช่จับผิด)
* **เป้าหมาย:** รักษามาตรฐานทีม และแชร์ความรู้ร่วมกัน


### 5. Code Update — แก้ไขตามรีวิว

* **หน้าที่:** Developer
* **การกระทำ:**
    * อ่าน Comment และแก้ไขโค้ดตามคำแนะนำ
    * Push โค้ดใหม่ขึ้นไป (ระบบ CI จะรันซ้ำอัตโนมัติ)
    * Reply ใน PR ว่า "แก้แล้วครับ" เพื่อเรียกคนตรวจกลับมาดู
* **เป้าหมาย:** ปรับปรุงงานให้สมบูรณ์ที่สุดก่อนรวมร่าง


### 6. Deploy — ปล่อยสู่ระบบ

* **หน้าที่:** Team Lead / System
* **การกระทำ:**
    * เมื่อ Reviewer กด **Approve** และ CI เป็นสีเขียว ✅
    * กดปุ่ม **Merge** โค้ดเข้าสู่ Branch หลัก
    * ระบบ CD (Continuous Deployment) จะนำโค้ดใหม่ไป Deploy ขึ้น Server (Staging/Production) ทันที
* **เป้าหมาย:** ส่งมอบคุณค่า (Value) ให้ผู้ใช้งานอย่างรวดเร็วและมั่นใจ

***

### แผนภาพสรุป (Workflow Diagram)

```mermaid
graph TD
    A[1. Code Assistant<br/>(Dev เขียนโค้ด)] -->|Push Code| B[3. CI Pipeline<br/>(Auto Test/Lint)]
    B -->|❌ ไม่ผ่าน| A
    B -->|✅ ผ่าน| C[2. Pull Request<br/>(ขอตรวจงาน)]
    C --> D[4. Code Review<br/>(ทีมช่วยดู)]
    D -->|❌ ขอแก้| E[5. Code Update<br/>(Dev แก้ตามสั่ง)]
    E -->|Push ใหม่| B
    D -->|✅ Approve| F[6. Deploy<br/>(Merge & Run)]
```

**เหมาะกับใคร?**

* **Developer ทุกระดับ:** เพื่อให้รู้หน้าที่ตัวเองในแต่ละขั้นตอน
* **ทีมที่ต้องการมาตรฐาน:** ลดการ "หลุด" ของบั๊กง่ายๆ
* **Tech Lead:** ใช้สอนน้องใหม่ในทีมให้เข้าใจ Flow การทำงานร่วมกันได้อย่างรวดเร็ว
<span style="display:none">[^1][^2][^3][^4][^5][^6][^7][^8]</span>

<div align="center">⁂</div>

[^1]: https://sennalabs.com/blog/code-review

[^2]: https://www.facebook.com/groups/programmerthailand/posts/4108233416091932/

[^3]: https://www.somkiat.cc/keys-for-code-review/

[^4]: https://www.reddit.com/r/ExperiencedDevs/comments/1k53zol/how_to_give_code_reviews_without_offending_other/

[^5]: https://www.techstarthailand.com/blog/detail/7-steps-to-better-code-reviews/634

[^6]: https://expert-programming-tutor.com/tutorial/article/KE000561_What_is_the_most_simple_code_reviews_8_years_old_children_understand_how_to_write_programs_With_the_easiest_example.php

[^7]: https://www.saladpuk.com/basic/agile-methodology/code-review

[^8]: https://www.youtube.com/watch?v=9Vc7zueVbRo

