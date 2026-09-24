 
# วิธีหาจำนวนคอลัมน์ก่อนใช้ UNION SELECT

เทคนิคการหาจำนวนคอลัมน์ (Column Enumeration) เป็น **ขั้นตอนแรกที่สำคัญที่สุด** ในการทำ SQL Injection แบบ Union-Based เพราะกฎเหล็กของ `UNION` คือ **"จำนวนคอลัมน์ของทั้งสองตารางต้องเท่ากัน"** ถ้าไม่เท่ากัน Database จะแจ้ง Error ทันที

มี 2 วิธีหลักที่ Reviewer สามารถใช้ทดสอบได้:

### 1. วิธีใช้ `ORDER BY` (วิธีที่นิยมที่สุด)

วิธีนี้ง่ายและเร็วกว่า โดยใช้การสั่งเรียงลำดับข้อมูลตามเลขคอลัมน์ ถ้าเราสั่งเรียงตามคอลัมน์ที่ "ไม่มีอยู่จริง" Database จะ Error

**สมมติ Query หลังบ้านเป็น:**
`SELECT id, title, content FROM news WHERE id = $id` (มี 3 คอลัมน์)

**การทดสอบ:** ใส่เลขไปเรื่อยๆ ใน URL หรือ Input

* `id=1 ORDER BY 1 --` ✅ (ผ่าน: มีคอลัมน์ที่ 1)
* `id=1 ORDER BY 2 --` ✅ (ผ่าน: มีคอลัมน์ที่ 2)
* `id=1 ORDER BY 3 --` ✅ (ผ่าน: มีคอลัมน์ที่ 3)
* `id=1 ORDER BY 4 --` ❌ **(Error!)**
    * *สรุป:* ตารางนี้มี **3 คอลัมน์**

***

### 2. วิธีใช้ `UNION SELECT NULL` (วิธีที่แม่นยำกว่า)

ถ้าวิธีแรกใช้ไม่ได้ผล (เช่น ระบบปิด Error Message ไว้) ให้ลองใช้ `UNION SELECT` แล้วเพิ่มจำนวน `NULL` ไปเรื่อยๆ จนกว่าหน้าเว็บจะแสดงผลปกติ

**การทดสอบ:**

* `id=1 UNION SELECT NULL --` ❌ (Error: คอลัมน์ไม่เท่า)
* `id=1 UNION SELECT NULL, NULL --` ❌ (Error: คอลัมน์ไม่เท่า)
* `id=1 UNION SELECT NULL, NULL, NULL --` ✅ **(ผ่าน!)** หน้าเว็บโหลดข้อมูลปกติ
    * *สรุป:* ตารางนี้มี **3 คอลัมน์**

**ทำไมต้องใช้ NULL?**
เพราะ `NULL` สามารถเข้าได้กับทุก Data Type (ไม่ว่าคอลัมน์นั้นจะเป็น Int, String, หรือ Date) ทำให้เราไม่ต้องเดาชนิดข้อมูลในขั้นตอนนี้ครับ

***

### สรุปตารางเปรียบเทียบ

| วิธีการ | ข้อดี | ข้อเสีย |
| :-- | :-- | :-- |
| **ORDER BY** | เร็ว, พิมพ์น้อย | อาจใช้ไม่ได้กับบาง Database ที่เคร่งครัดมาก |
| **UNION SELECT NULL** | แม่นยำ, ใช้ได้เกือบทุก DB | พิมพ์เยอะ, ต้องลองผิดลองถูกหลายรอบ |

### 🛡️ แนวทางแก้ไข (สำหรับ Developer)

ถ้าเจอช่องโหว่นี้ในโค้ด ไม่ใช่ไปแก้ที่ SQL แต่ต้องแก้ที่ **วิธีส่งค่า**

* **ห้าม:** เอา Input ไปต่อท้าย `ORDER BY` ตรงๆ
* **ต้องทำ:** ตรวจสอบค่า Input ก่อนเสมอ (Allowlist)

```typescript
// ✅ ตัวอย่างการป้องกัน (Node.js)
const allowedSorts = ['id', 'title', 'date'];
const sortBy = allowedSorts.includes(req.query.sort) ? req.query.sort : 'id'; // ถ้าส่งค่าแปลกๆ มา ให้ใช้ค่า default

const sql = `SELECT * FROM news ORDER BY ${sortBy}`; // ปลอดภัยเพราะค่า sortBy มาจาก list ที่เราคุมเอง
```


<span style="display:none">[^1][^2][^3][^4][^5][^6][^7][^8][^9]</span>

<div align="center">⁂</div>

[^1]: https://saixiii.com/database-sql-union-clause/

[^2]: https://www.reddit.com/r/SQL/comments/1fmbj1d/how_to_union_all_two_select_statements_where_for/

[^3]: https://www.reddit.com/r/SQL/comments/14fbypt/can_somebody_explain_this_rule_for_union_the/

[^4]: https://www.mindphp.com/บทเรียนออนไลน์/สอน-sql/3928-postgresql-sql-union.html

[^5]: https://www.codebee.co.th/labs/วิธีใช้งาน-union-ใน-mysql/

[^6]: https://www.9experttraining.com/articles/dax-functions-a-to-z

[^7]: https://staff.informatics.buu.ac.th/~komate/886301/DB-Chpater-8.pdf

[^8]: https://www.guru99.com/th/unions.html

[^9]: https://www.youtube.com/watch?v=oIqMpafmYbc

