
# คู่มือ SQL สำหรับ PostgreSQL   (Practical Guide)

## 1. วัตถุประสงค์
เพื่อให้ผู้เรียนสามารถเขียน SQL บน PostgreSQL ได้อย่างถูกต้อง มีประสิทธิภาพ เข้าใจการออกแบบฐานข้อมูล ฟังก์ชัน โพรซีเยอร์ และสามารถป้องกันช่องโหว่ SQL Injection ได้ด้วยตนเอง

## 2. กลุ่มเป้าหมาย
- นักพัฒนาโปรแกรมเมอร์ (Back-end, Full-stack)
- Data Analyst / Data Engineer
- ผู้ดูแลฐานข้อมูล (DBA) มือใหม่ถึงกลาง
- นักศึกษาที่ต้องการเรียนรู้ PostgreSQL เชิงปฏิบัติ

## 3. ความรู้พื้นฐาน
- พื้นฐานการเขียนโปรแกรม (ตัวแปร, เงื่อนไข, ลูป)
- ความเข้าใจเรื่องฐานข้อมูลเชิงสัมพันธ์ (Tables, Rows, Columns)
- รู้จัก PostgreSQL environment (pgAdmin, psql, หรือ DBeaver)

---

## บทนำ
PostgreSQL เป็นระบบฐานข้อมูลเชิงวัตถุ-เชิงสัมพันธ์ (ORDBMS) ที่มีประสิทธิภาพสูง รองรับมาตรฐาน SQL ได้ดีที่สุดตัวหนึ่ง คู่มือนี้จะไม่เพียงแต่สอนไวยากรณ์ แต่จะเน้น **การออกแบบที่นำไปใช้จริง** เช่น การสร้างฟังก์ชันที่ reuse ได้, การเขียน Query แบบมีเงื่อนไข, การป้องกันภัยคุกคาม และการจัดโครงสร้างโค้ดในฐานข้อมูล

---

## บทนิยาม
| ศัพท์ | ความหมาย |
|------|------------|
| **Query** | คำสั่งที่ใช้สอบถามหรือจัดการข้อมูล |
| **Function** | โค้ดที่封装 ไว้ใน DB รับพารามิเตอร์และคืนค่า |
| **Stored Procedure** | คล้าย Function แต่ไม่จำเป็นต้องคืนค่า และทำ Transaction ได้ |
| **Subquery** | Query ที่ซ้อนอยู่ในอีก Query หนึ่ง |
| **SQL Injection** | การแทรกคำสั่ง SQL โดยใช้ input ที่ไม่ถูก sanitize |

---

## บทหัวข้อ (สารบัญคู่มือ)

### ส่วนที่ 1: พื้นฐาน SQL (ทุกหัวข้อที่คุณระบุ)
เราจะสรุปแต่ละหัวข้อพร้อมตัวอย่างสั้น ๆ ที่ใช้ได้จริง

#### SQL Syntax
```sql
-- คำสั่งพื้นฐานลงท้ายด้วย ;
SELECT * FROM customers;
```

#### SQL SELECT
```sql
SELECT first_name, last_name FROM employees;
```

#### SELECT DISTINCT
```sql
SELECT DISTINCT city FROM customers;
```

#### WHERE
```sql
SELECT * FROM products WHERE price > 100;
```

#### ORDER BY
```sql
SELECT name, price FROM products ORDER BY price DESC;
```

#### AND, OR, NOT
```sql
SELECT * FROM orders WHERE status = 'paid' AND total >= 1000;
SELECT * FROM users WHERE city = 'Bangkok' OR city = 'Chiang Mai';
SELECT * FROM products WHERE NOT discontinued = true;
```

#### INSERT INTO
```sql
INSERT INTO customers (name, email) VALUES ('สมชาย', 'somchai@email.com');
```

#### NULL Values
```sql
SELECT * FROM employees WHERE phone IS NULL;
UPDATE employees SET phone = '000-000-0000' WHERE phone IS NULL;
```

#### UPDATE, DELETE
```sql
UPDATE products SET price = price * 1.10 WHERE category = 'electronic';
DELETE FROM logs WHERE created_at < NOW() - INTERVAL '30 days';
```

#### SELECT TOP (ใช้ LIMIT ใน PostgreSQL)
```sql
SELECT * FROM sales ORDER BY amount DESC LIMIT 10;
```

#### Aggregate Functions
```sql
SELECT MIN(price), MAX(price), COUNT(*), SUM(quantity), AVG(price) FROM orders;
```

#### LIKE & Wildcards
```sql
SELECT * FROM customers WHERE name LIKE 'สม%';  -- ขึ้นต้นด้วย สม
SELECT * FROM products WHERE code LIKE '_A%';    -- ตัวที่สองคือ A
```

#### IN, BETWEEN
```sql
SELECT * FROM employees WHERE dept_id IN (1,2,3);
SELECT * FROM products WHERE price BETWEEN 500 AND 1500;
```

#### Aliases
```sql
SELECT first_name AS "ชื่อ" FROM users;
```

#### JOINs (สำคัญ)
```sql
-- INNER JOIN
SELECT orders.id, customers.name
FROM orders
INNER JOIN customers ON orders.cust_id = customers.id;

-- LEFT JOIN
SELECT customers.name, orders.id
FROM customers
LEFT JOIN orders ON customers.id = orders.cust_id;

-- RIGHT JOIN, FULL JOIN, SELF JOIN
```

#### UNION / UNION ALL
```sql
SELECT name FROM staff_old
UNION
SELECT name FROM staff_new;  -- ตัดซ้ำ
```

#### GROUP BY + HAVING
```sql
SELECT dept_id, COUNT(*) 
FROM employees
GROUP BY dept_id
HAVING COUNT(*) > 5;
```

#### EXISTS, ANY, ALL
```sql
SELECT * FROM products p
WHERE EXISTS (SELECT 1 FROM orders o WHERE o.product_id = p.id);

SELECT * FROM products WHERE price > ANY (SELECT price from products where category='premium');
```

#### CASE (เงื่อนไขใน SQL)
```sql
SELECT name,
  CASE 
    WHEN score >= 80 THEN 'A'
    WHEN score >= 70 THEN 'B'
    ELSE 'C'
  END AS grade
FROM students;
```

#### Stored Procedures (PostgreSQL ใช้ CREATE PROCEDURE)
```sql
CREATE PROCEDURE transfer_money(from_acc int, to_acc int, amount dec)
LANGUAGE plpgsql
AS $$
BEGIN
  UPDATE accounts SET balance = balance - amount WHERE id = from_acc;
  UPDATE accounts SET balance = balance + amount WHERE id = to_acc;
  COMMIT;
END; $$;
```

#### Views
```sql
CREATE VIEW active_users AS
SELECT * FROM users WHERE status = 'active';
```

---

## การออกแบบ SQL Functions (ละเอียด)

### รูปแบบของ Function ใน PostgreSQL
1. **Scalar function** – คืนค่าธรรมดา (int, text, date)
2. **Set-returning function** – คืนค่าเป็นชุดข้อมูล (SETOF หรือ TABLE)
3. **Polymorphic function** – รองรับหลายชนิดข้อมูล

### ตัวอย่างจริง: คำนวณภาษีขาย
```sql
CREATE FUNCTION calculate_vat(price numeric, vat_rate numeric DEFAULT 0.07)
RETURNS numeric
LANGUAGE plpgsql
IMMUTABLE
AS $$
BEGIN
  RETURN price * vat_rate;
END; $$;

-- ใช้งาน
SELECT calculate_vat(1000, 0.07); -- 70
```

### ข้อควรระวัง
- ใช้ `IMMUTABLE` ถ้าผลลัพธ์ขึ้นกับ input เท่านั้น (ช่วยให้ index ทำงานเร็ว)
- ไม่ควรใช้ Function ที่มี `INSERT/UPDATE` โดยไม่ระบุ `VOLATILE`

---

## การออกแบบ SQL Subquery Functions

Subquery คือการเอา Query ไปแทรกใน `SELECT`, `FROM`, `WHERE`, `HAVING`

### ประเภท
- **Scalar Subquery** – คืนค่าเดี่ยว
- **Row Subquery** – คืนค่าหนึ่งแถว
- **Table Subquery** – คืนค่าหลายแถว

### ตัวอย่างการใช้งานจริง
```sql
-- หาพนักงานที่มีเงินเดือนสูงกว่าแผนกเฉลี่ย
SELECT name, salary, dept_id
FROM employees e
WHERE salary > (SELECT AVG(salary) FROM employees WHERE dept_id = e.dept_id);
```

### Subquery ใน FROM (Derived Table)
```sql
SELECT dept_id, avg_salary
FROM (
  SELECT dept_id, AVG(salary) as avg_salary
  FROM employees
  GROUP BY dept_id
) AS dept_avg
WHERE avg_salary > 50000;
```

### ข้อดี
- อ่านง่าย แก้ปัญหา complex query
- ไม่ต้องสร้าง view ชั่วคราว

### ข้อเสีย
- อาจช้า ถ้า subquery ถูกเรียกทุกแถว (ใช้ `EXPLAIN` ตรวจสอบ)

---

## การออกแบบ SQL Query ด้วย CASE WHEN, IF-ELSE, SWITCH

### CASE WHEN (มาตรฐาน SQL)
```sql
SELECT 
  order_id,
  total,
  CASE 
    WHEN total >= 1000 THEN 'แพง'
    WHEN total >= 500 THEN 'กลาง'
    ELSE 'ถูก'
  END AS price_level
FROM orders;
```

### IF-ELSE ใน PostgreSQL (เฉพาะใน PL/pgSQL)
```sql
DO $$
DECLARE
  score int := 85;
BEGIN
  IF score >= 80 THEN
    RAISE NOTICE 'Grade A';
  ELSIF score >= 70 THEN
    RAISE NOTICE 'Grade B';
  ELSE
    RAISE NOTICE 'Grade C';
  END IF;
END; $$;
```

### SWITCH (ใช้ CASE แทนได้)
ไม่มี SWITCH โดยตรง แต่ใช้ CASE expression ได้เหมือนกัน

### CONTINUE ใน Loop
```sql
FOR i IN 1..10 LOOP
  IF i % 2 = 0 THEN
    CONTINUE;  -- ข้ามเลขคู่
  END IF;
  RAISE NOTICE 'Odd number: %', i;
END LOOP;
```

---

## การออกแบบ SQL Procedures (Stored Procedures)

### แตกต่างจาก Function อย่างไร?
| Feature | Function | Procedure |
|--------|-----------|-------------|
| คืนค่า | ต้องคืนค่า (RETURNS) | ไม่จำเป็น |
| ใช้ใน SELECT | ได้ | ไม่ได้ |
| Transaction ควบคุม | ไม่ได้ (อยู่ใน atomic context) | ได้ (COMMIT/ROLLBACK) |
| OUT Parameters | ได้ | ได้ |

### ตัวอย่าง Procedure โอนเงิน พร้อม Transaction
```sql
CREATE OR REPLACE PROCEDURE transfer_funds(
  from_account INT,
  to_account INT,
  amount DECIMAL
)
LANGUAGE plpgsql
AS $$
BEGIN
  -- หักเงิน
  UPDATE accounts SET balance = balance - amount WHERE id = from_account;
  -- เพิ่มเงิน
  UPDATE accounts SET balance = balance + amount WHERE id = to_account;
  
  -- บันทึกประวัติ
  INSERT INTO transfer_logs (from_acc, to_acc, amount, trans_date)
  VALUES (from_account, to_account, amount, NOW());
  
  COMMIT;
EXCEPTION
  WHEN OTHERS THEN
    ROLLBACK;
    RAISE;
END; $$;

-- เรียกใช้
CALL transfer_funds(101, 202, 5000);
```

### ข้อห้าม
- ห้ามใช้ Procedure ในการ query ข้อมูลเพื่อแสดงผลในแอพ (ควรใช้ Function หรือ View)
- อย่าใช้ Procedure ถ้าคุณต้องการค่าคืนใน SELECT

---

## การหา SQL Injection ในโค้ดแบบเป็นขั้นตอน

### ขั้นตอนที่ 1: ตรวจสอบการต่อ string สร้าง SQL
```python
# เสี่ยงมาก
user_input = request.GET['username']
query = f"SELECT * FROM users WHERE name = '{user_input}'"
```
**วิธีตรวจจับ:** ค้นหา pattern `+`, `f-string`, `format` ที่เอา input มาแปะใน SQL

### ขั้นตอนที่ 2: ทดสอบด้วย payload
ลองป้อน input:
```
' OR '1'='1
'; DROP TABLE users; --
```
ถ้าเห็นข้อมูลทั้งหมด หรือ error แปลก แสดงว่ามีช่องโหว่

### ขั้นตอนที่ 3: ตรวจสอบการใช้ Dynamic Identifier (table name, column name)
```sql
-- เสี่ยง
EXECUTE 'SELECT * FROM ' || table_name;
```
ต้องใช้ `quote_ident()` ใน PostgreSQL:
```sql
EXECUTE 'SELECT * FROM ' || quote_ident(table_name);
```

### ขั้นตอนที่ 4: ใช้ Parameterized Query (Prepared Statement) เสมอ
```python
# ปลอดภัย
cursor.execute("SELECT * FROM users WHERE name = %s", (user_input,))
```

### ขั้นตอนที่ 5: ใช้ ORM (SQLAlchemy, Prisma, TypeORM) จะช่วยลดความเสี่ยง

### ขั้นตอนที่ 6: ใช้เครื่องมือสแกนอัตโนมัติ (sqlmap, Burp Suite, SonarQube)

---

## CLI Commands สำหรับ Admin Tasks (migrate, seed)

### CLI commands มีกี่แบบ?
สำหรับ PostgreSQL + โปรเจคทั่วไป มี 3 แบบ:
1. **psql native commands** – `\l`, `\dt`, `\d table`
2. **Migration tools** (Alembic, Flyway, golang-migrate)
3. **Seed commands** – custom scripts หรือ `COPY`

### ใช้อย่างไร / นำไปใช้กรณีไหน
| คำสั่ง | วิธีใช้ | กรณี |
|-------|--------|------|
| `psql -d mydb -f migrate.sql` | รันไฟล์ SQL | Deploy ฐานข้อมูล |
| `alembic upgrade head` | อัปเกรด schema | ระหว่างพัฒนา |
| `psql -c "COPY users FROM 'data.csv' CSV HEADER"` | นำเข้าข้อมูล | ตั้งค่าข้อมูลเริ่มต้น |

### ทำไมต้องใช้
- เพื่อให้การเปลี่ยนแปลงฐานข้อมูลเป็นระบบ (version control)
- ลดความผิดพลาดจากการทำ manual
- ทำซ้ำได้ (repeatable)

### ประโยชน์ที่ได้รับ
- rollback ได้เมื่อ schema ผิดพลาด
- seed ข้อมูลทดสอบซ้ำ ๆ
- ทีมงานใช้ migration เดียวกัน

### ข้อควรระวัง
- อย่าแก้ migration ที่ถูก apply แล้วโดยตรง (ให้สร้าง migration ใหม่)
- ทดสอบ migration ใน environment ทดสอบก่อน

### ข้อดี
- ปลอดภัยสูง
- ทำงานอัตโนมัติใน CI/CD

### ข้อเสีย
- ต้องเรียนรู้เครื่องมือเพิ่ม
- บางที debug ยาก

### ข้อห้าม
- ห้ามใช้ migration บน production โดยไม่ backup
- ห้าม seed ข้อมูลจริง (real PII) ในไฟล์ seed

---

## ออกแบบคู่มือ (เอกสารเล่มนี้)

คู่มือถูกออกแบบเป็น 3 ส่วน:
- **Part A:** Basic SQL (ทุกคำสั่งที่คุณระบุ)
- **Part B:** Advanced Design (Function, Procedure, Subquery, CASE)
- **Part C:** Security & Administration (Injection, CLI, Migration)

แต่ละส่วนมี:
- บทนำ
- ตารางสรุปไวยากรณ์
- ตัวอย่างจริง
- ข้อควรระวัง + ข้อห้าม
- แบบฝึกหัดสั้น ๆ

---

## ออกแบบ Workflow (สำหรับการพัฒนา SQL ในโปรเจค)

```
[ต้องการรายงาน/ฟีเจอร์ใหม่] 
  ↓
[เขียน Query ใน SQL Editor (DBeaver/pgAdmin)]
  ↓
[ใช้ EXPLAIN ANALYZE ดูประสิทธิภาพ]
  ↓
[สร้าง Function หรือ Procedure ถ้าจำเป็น]
  ↓
[เขียน Migration (.sql) + Rollback]
  ↓
[ทดสอบด้วยข้อมูลตัวอย่าง]
  ↓
[Commit + Deploy ผ่าน CI (รัน migration อัตโนมัติ)]
  ↓
[Monitor ด้วย pg_stat_statements]
```

---

## TASK LIST Template (สำหรับงาน SQL)

| Task ID | รายการ | สถานะ | หมายเหตุ |
|---------|--------|--------|----------|
| T01 | สร้างตาราง `users` พร้อม PK | ✅ Done | |
| T02 | เขียน Function `get_user_orders(user_id)` | 🔄 In Progress | |
| T03 | ป้องกัน SQL Injection ทุก Query | ⬜ Pending | |
| T04 | สร้าง Index บน `created_at` | ✅ Done | |
| T05 | Seed ข้อมูลทดสอบ 1000 แถว | ⬜ Pending | ใช้ generate_series |

---

## CHECKLIST Template (ก่อนขึ้น Production)

- [ ] ทุก Query ใช้ Parameterized หรือ ORM
- [ ] มีการทดสอบ SQL Injection ด้วย payload `' OR '1'='1`
- [ ] มี Index ที่เหมาะสม (ดูจาก EXPLAIN)
- [ ] ไม่มี `SELECT *` บนตารางใหญ่
- [ ] ฟังก์ชันระบุ `IMMUTABLE`/`STABLE`/`VOLATILE` ถูกต้อง
- [ ] Procedure มีการจัดการ Exception และ ROLLBACK
- [ ] Migration มีทั้ง up และ down
- [ ] Backup database ก่อน deploy

---

## สรุป

คู่มือนี้ครอบคลุม:
- **ทุกหัวข้อ SQL พื้นฐาน** ที่คุณระบุ พร้อมตัวอย่างรวบรัดแต่ใช้ได้จริง
- **การออกแบบ Functions** ทั้งแบบ Scalar และ Set-returning
- **Subquery** ชนิดต่าง ๆ และข้อควรระวังเรื่องประสิทธิภาพ
- **CASE WHEN, IF-ELSE, CONTINUE** ใน PostgreSQL
- **Stored Procedures** เทียบกับ Functions พร้อมตัวอย่างโอนเงิน
- **การหา SQL Injection** แบบเป็นขั้นตอน (6 ขั้นตอน)
- **CLI commands** สำหรับ migrate/seed พร้อมข้อดี/ข้อเสีย/ข้อห้าม
- **Workflow, Task List, Checklist** สำหรับการทำงานจริง

> **คำแนะนำสุดท้าย:** อย่าเรียนแค่ไวยากรณ์ – ให้ลองสร้างฐานข้อมูลจำลอง (เช่น ร้านค้า, ระบบจองโรงแรม) แล้วลองเขียน Query ทุกแบบด้วยตัวเอง เสร็จแล้วใช้ `EXPLAIN ANALYZE` ปรับปรุง performance เสมอ

# ส่วนเพิ่มเติม: Trigger, Window Functions, CTE สำหรับ PostgreSQL

คู่มือนี้ต่อเนื่องจากคู่มือหลัก โดยเจาะลึก 3 หัวข้อสำคัญที่ช่วยให้เขียน SQL เชิงวิเคราะห์และบริหารจัดการข้อมูลได้อย่างมืออาชีพ

---

## 1. Trigger (ทริกเกอร์)

### ความหมาย
Trigger เป็นกลไกที่ทำให้ฟังก์ชันทำงานโดยอัตโนมัติ เมื่อมีเหตุการณ์ `INSERT`, `UPDATE`, `DELETE` เกิดขึ้นบนตารางที่กำหนด

### ประเภทของ Trigger ใน PostgreSQL
| ประเภท | คำอธิบาย |
|--------|------------|
| `BEFORE` | ทำงานก่อนการเปลี่ยนแปลงข้อมูล (ใช้ตรวจสอบหรือแก้ไขข้อมูลก่อนบันทึก) |
| `AFTER`  | ทำงานหลังการเปลี่ยนแปลงข้อมูลสำเร็จ (ใช้บันทึกประวัติ หรือแจ้งเตือน) |
| `INSTEAD OF` | ใช้กับ View แทนการทำงานปกติ (rare) |
| `FOR EACH ROW` | ทำงานทุกแถวที่ถูกเปลี่ยนแปลง |
| `FOR EACH STATEMENT` | ทำงานครั้งเดียวต่อคำสั่ง SQL (แม้เปลี่ยนหลายแถว) |

### โครงสร้างการสร้าง Trigger
```sql
-- 1. สร้างฟังก์ชันที่ return type = TRIGGER
CREATE OR REPLACE FUNCTION function_name()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
   -- ใช้ NEW, OLD เข้าถึงข้อมูล
   -- NEW คือแถวใหม่ (สำหรับ INSERT, UPDATE)
   -- OLD คือแถวเก่า (สำหรับ UPDATE, DELETE)
   RETURN NEW|OLD;
END; $$;

-- 2. สร้าง Trigger ผูกกับตาราง
CREATE TRIGGER trigger_name
{BEFORE | AFTER} {INSERT | UPDATE | DELETE} ON table_name
FOR EACH ROW
EXECUTE FUNCTION function_name();
```

### ตัวอย่างจริง 1: อัปเดต `updated_at` อัตโนมัติ (ใช้บ่อยที่สุด)
```sql
-- สร้างฟังก์ชัน
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- สร้าง trigger
CREATE TRIGGER trigger_update_users_modified
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();
```

### ตัวอย่างจริง 2: เก็บบันทึกประวัติ (Audit Log) เมื่อมีการลบ
```sql
-- ตารางเก็บ log
CREATE TABLE user_delete_log (
    id SERIAL PRIMARY KEY,
    deleted_user_id INT,
    deleted_data JSONB,
    deleted_at TIMESTAMP DEFAULT NOW()
);

-- ฟังก์ชัน trigger
CREATE OR REPLACE FUNCTION log_user_deletion()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO user_delete_log (deleted_user_id, deleted_data)
    VALUES (OLD.id, row_to_json(OLD));
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Trigger AFTER DELETE
CREATE TRIGGER trigger_log_user_delete
AFTER DELETE ON users
FOR EACH ROW
EXECUTE FUNCTION log_user_deletion();
```

### ข้อดี
- ช่วยรักษาความถูกต้องของข้อมูล (data integrity)
- ลดโค้ดในแอพพลิเคชัน (centralize business logic)
- ทำงานอัตโนมัติ ไม่ต้องเรียกเอง

### ข้อเสีย
- ซ่อน logic ทำให้ debug ยาก
- อาจทำให้ performance ช้าลง ถ้า trigger ทำงานหนัก
- ไม่สามารถเรียก trigger ด้วยตนเองได้ (เฉพาะอัตโนมัติ)

### ข้อควรระวัง
- Trigger ที่มี recursive (UPDATE แล้ว trigger อีก UPDATE) ต้องควบคุมให้ดี
- ระวังการใช้ `FOR EACH STATEMENT` กับตารางใหญ่ เพราะอาจทำงานไม่ตรงตามที่คิด
- ต้องมีสิทธิ์สร้างฟังก์ชันและ trigger

### ข้อห้าม
- ห้ามใช้ trigger สำหรับงานที่ควรทำในแอพ เช่น การส่ง email, call API (เพราะทำใน DB ไม่เหมาะ)
- ห้ามใช้ trigger ที่ซับซ้อนมากจนทำให้ transaction ยาวเกินไป

---

## 2. Window Functions (ฟังก์ชันหน้าต่าง)

### ความหมาย
Window functions คำนวณค่าตามกลุ่มของแถว (window) โดยไม่ยุบรวมแถวเหมือน `GROUP BY` – แต่ละแถวยังคงอยู่และได้ค่าเพิ่มเติมจากการคำนวณข้ามแถว

### ไวยากรณ์พื้นฐาน
```sql
function_name(...) OVER (
    [PARTITION BY column1, column2...]
    [ORDER BY column1 [ASC|DESC]]
    [ROWS | RANGE BETWEEN frame_start AND frame_end]
)
```

### ประเภทฟังก์ชันยอดนิยม
| ฟังก์ชัน | การทำงาน |
|----------|------------|
| `ROW_NUMBER()` | ลำดับที่ (ไม่ซ้ำ) |
| `RANK()` | ลำดับที่มีอันดับเว้นช่องว่าง |
| `DENSE_RANK()` | ลำดับไม่เว้นช่องว่าง |
| `LEAD(column, offset)` | ค่าของแถวถัดไป |
| `LAG(column, offset)` | ค่าของแถวก่อนหน้า |
| `SUM() OVER()` | ยอดรวมสะสม |
| `AVG() OVER()` | ค่าเฉลี่ยเคลื่อนที่ |

### ตัวอย่างจริง 1: จัดอันดับพนักงานตามเงินเดือนในแต่ละแผนก
```sql
SELECT 
    name,
    dept_id,
    salary,
    ROW_NUMBER() OVER (PARTITION BY dept_id ORDER BY salary DESC) AS rank_in_dept
FROM employees;
```

### ตัวอย่างจริง 2: ยอดขายสะสมของแต่ละวัน (cumulative sum)
```sql
SELECT 
    sale_date,
    amount,
    SUM(amount) OVER (ORDER BY sale_date ROWS UNBOUNDED PRECEDING) AS running_total
FROM sales;
```

### ตัวอย่างจริง 3: เปรียบเทียบเงินเดือนกับพนักงานที่อยู่ก่อนหน้า (LAG)
```sql
SELECT 
    name,
    hire_date,
    salary,
    LAG(salary, 1) OVER (ORDER BY hire_date) AS prev_salary,
    salary - LAG(salary, 1) OVER (ORDER BY hire_date) AS diff_from_prev
FROM employees;
```

### ข้อดี
- เขียน简洁 比 self-join หรือ subquery อ่านง่ายกว่า
- ประสิทธิภาพดีกว่าการ join ตัวเองหลายรอบ
- เหมาะกับการวิเคราะห์ (reporting, dashboard)

### ข้อเสีย
- ใช้ทรัพยากร memory สูง ถ้า partition ใหญ่
- ไม่สามารถใช้ใน `WHERE` ได้โดยตรง (ต้องทำ subquery หรือ CTE)
- ผู้เริ่มต้นอาจงงกับ syntax

### ข้อควรระวัง
- `ORDER BY` ใน `OVER` มีผลต่อการคำนวณ frame เริ่มต้น (default คือ `RANGE UNBOUNDED PRECEDING`)
- ถ้าไม่ใส่ `ORDER BY` ฟังก์ชันอันดับจะไม่มีความหมาย

### ข้อห้าม
- ห้ามใช้ window function ซ้อนกัน (nest) โดยตรง – ต้องใช้ CTE หรือ subquery

---

## 3. CTE – Common Table Expression (WITH clause)

### ความหมาย
CTE เป็นตารางชั่วคราวที่อยู่ในคำสั่ง SQL เดียวกัน ช่วยให้ query ซับซ้อนอ่านง่าย และสามารถเรียกซ้ำได้ (Recursive CTE)

### ประเภท
- **Non-recursive CTE** – ใช้แทน subquery หรือสร้าง temporary view
- **Recursive CTE** – เรียกตัวเอง เหมาะกับโครงสร้างต้นไม้ (tree/graph)

### ไวยากรณ์พื้นฐาน
```sql
WITH cte_name AS (
    SELECT query
)
SELECT * FROM cte_name;
```

### ตัวอย่างจริง 1: CTE แทน subquery ทำให้อ่านง่าย
```sql
WITH high_value_orders AS (
    SELECT customer_id, SUM(amount) AS total_spent
    FROM orders
    GROUP BY customer_id
    HAVING SUM(amount) > 10000
)
SELECT c.name, h.total_spent
FROM customers c
JOIN high_value_orders h ON c.id = h.customer_id;
```

### ตัวอย่างจริง 2: Recursive CTE สำหรับจัดลำดับผู้บังคับบัญชา (org chart)
```sql
WITH RECURSIVE org_tree AS (
    -- Anchor: หัวหน้าสูงสุด
    SELECT id, name, manager_id, 1 AS level
    FROM employees
    WHERE manager_id IS NULL
    
    UNION ALL
    
    -- Recursive: ไล่ลูกน้อง
    SELECT e.id, e.name, e.manager_id, ot.level + 1
    FROM employees e
    INNER JOIN org_tree ot ON e.manager_id = ot.id
)
SELECT * FROM org_tree ORDER BY level, name;
```

### ตัวอย่างจริง 3: หาเส้นทางในกราฟ (เช่น ระบบแนะนำเพื่อน)
```sql
WITH RECURSIVE friend_path AS (
    SELECT user_id, friend_id, 1 AS depth, ARRAY[user_id] AS path
    FROM friendships
    WHERE user_id = 1
    UNION ALL
    SELECT fp.user_id, f.friend_id, fp.depth + 1, fp.path || f.user_id
    FROM friend_path fp
    JOIN friendships f ON fp.friend_id = f.user_id
    WHERE f.friend_id <> ALL(fp.path) AND fp.depth < 5
)
SELECT * FROM friend_path;
```

### ข้อดี
- อ่านง่ายกว่า subquery ซับซ้อน
- recursive ทำให้ query ต้นไม้ทำได้ง่าย
- สามารถอ้างอิง CTE หลายตัวใน query เดียว

### ข้อเสีย
- Recursive CTE อาจช้ามากถ้าข้อมูลลึกหรือไม่มี index
- ไม่สามารถสร้าง index บน CTE ได้
- บาง optimizer จัดการ CTE ไม่ดีเท่า subquery (PostgreSQL 12+ ดีขึ้น)

### ข้อควรระวัง
- Recursive CTE ต้องมีเงื่อนไขสิ้นสุด (stop condition) มิฉะนั้น loop อนันต์
- PostgreSQL จำกัดจำนวน recursion ด้วย `max_recursion_depth`

### ข้อห้าม
- ห้ามใช้ CTE สำหรับชุดข้อมูลขนาดใหญ่ที่ต้องใช้หลายครั้ง (เพราะคำนวณใหม่ทุกครั้ง ยกเว้นใช้ `MATERIALIZED` hint)
- ห้ามใช้ recursive CTE บนตารางที่ไม่มี index บนคอลัมน์ที่ใช้ join

---

## การออกแบบ Workflow สำหรับใช้ Trigger, Window Functions, CTE ในโปรเจค

```
[เริ่ม需求]
   ↓
[เลือกเครื่องมือ] 
   - Trigger: เมื่อต้องการความอัตโนมัติ / audit / data integrity
   - Window Functions: เมื่อต้องการ ranking, running total, lag/lead
   - CTE: เมื่อ query ซับซ้อนหรือต้อง recursive
   ↓
[เขียนต้นแบบ (prototype)]
   ↓
[ทดสอบกับข้อมูลจริงปริมาณน้อย]
   ↓
[ตรวจสอบประสิทธิภาพด้วย EXPLAIN ANALYZE]
   ↓
[ปรับ index / ปรับ query]
   ↓
[สร้าง migration]
   ↓
[บันทึกเอกสารการใช้งาน]
```

---

## TASK LIST Template (สำหรับงานที่ใช้ Trigger / Window / CTE)

| Task ID | รายการ | เทคนิค | สถานะ |
|---------|--------|--------|--------|
| T06 | สร้าง Trigger อัปเดต updated_at ทุกตาราง | BEFORE UPDATE | ✅ |
| T07 | สร้าง Audit Log สำหรับ DELETE บน orders | AFTER DELETE + JSONB | 🔄 |
| T08 | เขียนรายงานจัดอันดับสินค้าขายดีตามหมวดหมู่ | RANK() OVER(PARTITION BY category) | ⬜ |
| T09 | แสดงยอดขายสะสมรายเดือน | SUM() OVER(ORDER BY month) | ⬜ |
| T10 | Query โครงสร้างองค์กร (CEO -> พนักงาน) | Recursive CTE | ✅ |

---

## CHECKLIST Template (เพิ่มเติม)

ก่อนใช้ Trigger / Window / CTE ใน production:

### Trigger
- [ ] ฟังก์ชัน trigger ถูกเขียนเป็น `RETURNS TRIGGER` และมี `LANGUAGE plpgsql`
- [ ] ทดสอบกรณี `NEW` / `OLD` สำหรับ INSERT, UPDATE, DELETE ครบ
- [ ] ตรวจสอบว่าไม่มี recursive trigger loop
- [ ] มีการจัดการ exception ใน trigger (ถ้าจำเป็น)

### Window Functions
- [ ] ใช้ `PARTITION BY` เฉพาะคอลัมน์ที่จำเป็น เพื่อลดขนาด window
- [ ] ระบุ `ROWS` หรือ `RANGE` ให้ชัดเจน ถ้าต้องการ cumulative แบบเจาะจง
- [ ] ไม่มี window function ใน `WHERE` – ใช้ subquery หรือ CTE แทน

### CTE
- [ ] Recursive CTE มี `UNION ALL` และมี anchor member + recursive member
- [ ] มีเงื่อนไขสิ้นสุด recursive (เช่น `WHERE depth < 10`)
- [ ] ทดสอบ query ด้วย `EXPLAIN` เพื่อดูว่า CTE ถูก materialized หรือไม่
- [ ] ถ้า CTE ถูกเรียกใช้ครั้งเดียว ไม่ต้องกังวล ถ้าถูกเรียกหลายครั้ง ให้พิจารณาใช้ `MATERIALIZED` หรือ temporary table

---

## สรุปส่วนเพิ่มเติม

| หัวข้อ | ใช้เมื่อใด | ข้อควรจำ |
|--------|------------|-------------|
| **Trigger** | ต้องการให้ DB ทำงานอัตโนมัติเมื่อมี DML | ระวัง performance และ recursion |
| **Window Functions** | ต้องการ ranking, cumulative, หรือเข้าถึงแถวข้างเคียงโดยไม่ GROUP BY | ไม่ใช้ใน WHERE ต้องใช้ subquery/CTE |
| **CTE** | query ซับซ้อน อ่านยาก หรือต้อง recursive | recursive ต้องมีเงื่อนไขสิ้นสุด ระวัง performance |

**คำแนะนำปฏิบัติ:** 
- เริ่มจากเขียน query แบบธรรมดา แล้วค่อยเพิ่ม window/CTE
- ใช้ `EXPLAIN (ANALYZE, BUFFERS)` เปรียบเทียบก่อน-หลัง
- เอกสารทุก trigger และ recursive CTE เพราะคนอื่นอาจงง

 # คู่มืออ้างอิงแบบสอบถาม SQL (SQL Query Reference Guide)

ด้านล่างนี้คือการสรุป **ทุกหัวข้อที่คุณระบุ** ในรูปแบบคำถาม-คำตอบ พร้อมตัวอย่างการใช้งานจริง ข้อดี ข้อเสีย และข้อห้ามสำหรับ PostgreSQL

---

## 1. SQL Query / Advanced SQL Query

### คืออะไร?
- **SQL Query** คือคำสั่งที่ใช้สอบถามหรือจัดการข้อมูลในฐานข้อมูล relational
- **Advanced SQL Query** คือเทคนิคการเขียน Query ขั้นสูง เช่น subquery, CTE, window functions, recursive query, และการ optimize ประสิทธิภาพ

### มีกี่แบบ?
แบ่งตามวัตถุประสงค์:
1. **Data Query Language (DQL)** – `SELECT`
2. **Data Manipulation Language (DML)** – `INSERT`, `UPDATE`, `DELETE`
3. **Data Definition Language (DDL)** – `CREATE`, `ALTER`, `DROP`
4. **Data Control Language (DCL)** – `GRANT`, `REVOKE`
5. **Transaction Control** – `BEGIN`, `COMMIT`, `ROLLBACK`

Advanced มีหลายเทคนิค: CTE, Window Functions, Recursive Query, Pivot, Lateral Join, etc.

### ข้อห้ามสำคัญ
- ห้ามใช้ `SELECT *` บนตารางขนาดใหญ่ใน Production
- ห้ามเขียน Query ที่ไม่มี `WHERE` บน `UPDATE`/`DELETE`
- ห้ามต่อ string สร้าง SQL แบบ dynamic โดยไม่ใช้ parameterized query (เสี่ยง SQL Injection)

### ใช้อย่างไร / นำไปใช้กรณีไหน
- Query ทั่วไป: ดึงข้อมูลตามเงื่อนไข, รายงาน, dashboard
- Advanced: วิเคราะห์แนวโน้ม, จัดอันดับ, tree structure (org chart), การ migrate ข้อมูลซับซ้อน

### ประโยชน์
- ดึงข้อมูลได้ตรงตามต้องการ
- ลดภาระการประมวลผลที่แอพพลิเคชัน
- ใช้ทรัพยากร database engine ให้เต็มประสิทธิภาพ

### ข้อควรระวัง
- Query ที่ซับซ้อนอาจทำให้ database ตายได้ (full table scan, missing index)
- ควรใช้ `EXPLAIN` วิเคราะห์ก่อนขึ้น Production

### ข้อดี
- มาตรฐานสากล เรียนครั้งเดียวใช้ได้หลาย DB
- ทรงพลัง สามารถทำ data transformation ได้ใน DB

### ข้อเสีย
- ไวยากรณ์บางอย่างต่างกันระหว่าง DB (Oracle, SQL Server, PostgreSQL)
- Debug ยากถ้า query ยาวมาก

### ข้อห้าม
- ห้ามเขียน nested subquery หลายชั้นโดยไม่จำเป็น (ให้ใช้ CTE แทน)
- ห้ามใช้ `SELECT DISTINCT` เป็นทางลัดโดยไม่เข้าใจข้อมูล

---

## 2. SQL Syntax

### คืออะไร?
ชุดของกฎและโครงสร้างที่ใช้เขียนคำสั่ง SQL ให้ถูกต้องและทำงานได้

### มีกี่แบบ?
- คำสั่งหลัก: `SELECT`, `INSERT`, `UPDATE`, `DELETE`, `CREATE`, `DROP`
- ส่วนประกอบ: `FROM`, `WHERE`, `GROUP BY`, `HAVING`, `ORDER BY`, `LIMIT`
- ตัวดำเนินการ: `=`, `<>`, `>`, `<`, `AND`, `OR`, `NOT`, `IN`, `BETWEEN`, `LIKE`

### ข้อห้ามสำคัญ
- ห้ามละเว้นเครื่องหมาย semicolon (`;`) ตอนมีหลายคำสั่ง
- ห้ามใช้ reserved word เป็นชื่อตาราง/คอลัมน์ โดยไม่ใส่ double quote (`"select"`)

### ใช้อย่างไร
```sql
SELECT column1, column2
FROM table_name
WHERE condition
ORDER BY column1;
```

### ประโยชน์
- เขียนได้มาตรฐาน เรียนง่าย
- ระบบสามารถ parse และ execute ได้

### ข้อควรระวัง
- PostgreSQL เป็น case-sensitive สำหรับชื่อที่ใส่ double quote
- String ต้องใช้ single quote (`'text'`)

### ข้อดี
- อ่านง่าย (clarity)
- มีเครื่องมือช่วย auto-complete มากมาย

### ข้อเสีย
- ความแตกต่างเล็กน้อยระหว่างรุ่น DB

### ข้อห้าม
- ห้ามเขียน SQL โดยไม่จัด indent (โค้ดรก อ่านยาก)

---

## 3. SQL SELECT

### คืออะไร?
คำสั่งที่ใช้ดึงข้อมูลจากฐานข้อมูล

### มีกี่แบบ?
- `SELECT *` – ดึงทุกคอลัมน์
- `SELECT DISTINCT` – ดึงเฉพาะค่าที่ไม่ซ้ำ
- `SELECT column1, column2` – เลือกบางคอลัมน์
- `SELECT INTO` – ดึงแล้วสร้างตารางใหม่

### ข้อห้ามสำคัญ
ห้ามใช้ `SELECT *` ใน production ถ้าตารางมีคอลัมน์ BLOB/JSON ขนาดใหญ่ หรือมีคอลัมน์ที่ไม่ได้ใช้

### ใช้อย่างไร
```sql
SELECT first_name, last_name FROM employees WHERE department = 'IT';
```

### ประโยชน์
เป็นพื้นฐานของทุกการอ่านข้อมูล

### ข้อควรระวัง
ถ้าลืม `WHERE` จะดึงข้อมูลทั้งหมด ทำให้ overload ได้

### ข้อดี
ยืดหยุ่น ใช้ร่วมกับ JOIN, WHERE, GROUP BY ได้

### ข้อเสีย
บางทีก็ซับซ้อนเกินไปสำหรับผู้เริ่มต้น

### ข้อห้าม
ห้าม `SELECT *` แล้วนำไป `INSERT INTO` อีกตารางโดยไม่ระบุคอลัมน์ (เสี่ยงตอน schema เปลี่ยน)

---

## 4. SQL SELECT DISTINCT

### คืออะไร?
เลือกเฉพาะแถวที่มีค่าแตกต่างกัน (unique combination)

### มีกี่แบบ?
- `SELECT DISTINCT col1` – unique values ของ col1
- `SELECT DISTINCT col1, col2` – unique pairs

### ข้อห้ามสำคัญ
ห้ามใช้ DISTINCT เพื่อ掩盖การ JOIN ที่ผิดพลาด (เช่น cartesian product)

### ใช้อย่างไร
```sql
SELECT DISTINCT city FROM customers;
```

### ประโยชน์
ลบข้อมูลซ้ำออกจากผลลัพธ์

### ข้อควรระวัง
DISTINCT ใช้ resources สูง ถ้ามีหลายคอลัมน์และข้อมูล

### ข้อดี
ได้ unique list ง่ายๆ

### ข้อเสีย
ไม่สามารถบอกได้ว่าเอาตัวไหนออก (random)

### ข้อห้าม
ห้ามใช้ DISTINCT ทุก query โดยไม่คิด เพราะอาจช้า

---

## 5. SQL WHERE

### คืออะไร?
กรองแถวตามเงื่อนไขที่กำหนด

### มีกี่แบบ?
- เงื่อนไข comparison: `=`, `>`, `<`, `>=`, `<=`, `<>`
- เงื่อนไข logical: `AND`, `OR`, `NOT`
- เงื่อนไข set: `IN`, `BETWEEN`, `LIKE`, `IS NULL`

### ข้อห้ามสำคัญ
ห้ามใช้ `WHERE` กับคอลัมน์ที่ถูกฟังก์ชันครอบ (เช่น `WHERE LOWER(name) = 'john'`) เพราะจะไม่ใช้ index

### ใช้อย่างไร
```sql
SELECT * FROM products WHERE price > 100 AND category = 'Electronics';
```

### ประโยชน์
ลดปริมาณข้อมูลที่ต้องประมวลผล

### ข้อควรระวัง
`NULL` เปรียบเทียบไม่ได้ ต้องใช้ `IS NULL` หรือ `IS NOT NULL`

### ข้อดี
ทำให้ query เร็วขึ้นถ้ามี index

### ข้อเสีย
เงื่อนไขซับซ้อนอาจทำให้ optimizer เลือกแผนไม่ดี

### ข้อห้าม
ห้ามเขียน `WHERE column = NULL` (ผิดตลอด)

---

## 6. SQL ORDER BY

### คืออะไร?
เรียงลำดับผลลัพธ์จากมากไปน้อย หรือน้อยไปมาก

### มีกี่แบบ?
- `ORDER BY col ASC` (default)
- `ORDER BY col DESC`
- เรียงหลายคอลัมน์: `ORDER BY col1 ASC, col2 DESC`

### ข้อห้ามสำคัญ
ห้ามใช้ `ORDER BY` กับคอลัมน์ที่ไม่มี index และมีข้อมูลหลายล้านแถว ถ้าไม่จำเป็น เพราะต้อง sort ทั้งชุด

### ใช้อย่างไร
```sql
SELECT name, salary FROM employees ORDER BY salary DESC LIMIT 10;
```

### ประโยชน์
แสดงข้อมูลเป็นลำดับ เหมาะกับการทำ ranking, report

### ข้อควรระวัง
`ORDER BY` จะทำงานหลังจาก `WHERE`, `GROUP BY`, `HAVING` แต่ก่อน `LIMIT`

### ข้อดี
ควบคุมลำดับได้เต็มที่

### ข้อเสีย
ใช้ memory และ CPU มาก ถ้าข้อมูล

### ข้อห้าม
ห้าม `ORDER BY RANDOM()` บนตารางใหญ่ (ช้ามาก)

---

## 7. SQL AND, OR, NOT

### คืออะไร?
ตัวดำเนินการตรรกะสำหรับรวมเงื่อนไขใน `WHERE` หรือ `HAVING`

### มีกี่แบบ?
- `AND` – ทุกเงื่อนไขเป็นจริง
- `OR` – อย่างน้อยหนึ่งเงื่อนไขเป็นจริง
- `NOT` – กลับค่าความจริง

### ข้อห้ามสำคัญ
ระวัง precedence: `NOT` > `AND` > `OR` – ควรใช้วงเล็บ `()` เสมอเมื่อมีหลายเงื่อนไข

### ใช้อย่างไร
```sql
SELECT * FROM users WHERE (status = 'active' OR role = 'admin') AND created_at > '2024-01-01';
```

### ประโยชน์
สร้างเงื่อนไขซับซ้อนได้

### ข้อควรระวัง
`OR` อาจทำให้ optimizer เลือกไม่ใช้ index (เปลี่ยนเป็น `UNION` บางครั้งดีกว่า)

### ข้อดี
flexible

### ข้อเสีย
ซับซ้อนเกินไปอาจอ่านยาก

### ข้อห้าม
ห้ามเขียน `WHERE col = 1 OR 2` (ต้อง `col = 1 OR col = 2`)

---

## 8. SQL INSERT INTO

### คืออะไร?
เพิ่มแถวใหม่ลงในตาราง

### มีกี่แบบ?
- `INSERT INTO table (col1, col2) VALUES (val1, val2)`
- `INSERT INTO table VALUES (val1, val2)` (ต้องตรงทุกคอลัมน์)
- `INSERT INTO table SELECT ...` (insert จาก query)

### ข้อห้ามสำคัญ
ห้าม insert โดยไม่ระบุรายชื่อคอลัมน์ ถ้าตารางมีการเปลี่ยนแปลงในอนาคต

### ใช้อย่างไร
```sql
INSERT INTO products (name, price) VALUES ('Laptop', 25000);
INSERT INTO logs SELECT * from temp_logs WHERE created_at > NOW() - interval '1 day';
```

### ประโยชน์
เพิ่มข้อมูล

### ข้อควรระวัง
Constraint (`NOT NULL`, `UNIQUE`, `FOREIGN KEY`) อาจทำให้ insert ล้มเหลว

### ข้อดี
สามารถ insert ทีละหลายแถว: `VALUES (1,'a'), (2,'b')`

### ข้อเสีย
ถ้าไม่มี `RETURNING` จะไม่รู้ id ที่ auto increment

### ข้อห้าม
ห้าม insert ค่าซ้ำใน primary key

---

## 9. SQL NULL Values

### คืออะไร?
NULL หมายถึง "ไม่มีค่า" หรือ "ไม่ทราบค่า" ไม่ใช่ 0 หรือ empty string

### มีกี่แบบ?
- `IS NULL`
- `IS NOT NULL`
- ฟังก์ชันจัดการ NULL: `COALESCE`, `NULLIF`

### ข้อห้ามสำคัญ
ห้ามใช้ `= NULL` หรือ `!= NULL` เพราะผลลัพธ์เป็น NULL เสมอ (ไม่เป็น true/false)

### ใช้อย่างไร
```sql
SELECT * FROM employees WHERE phone_number IS NULL;
UPDATE users SET middle_name = COALESCE(middle_name, 'N/A');
```

### ประโยชน์
แทนค่าที่ไม่มีข้อมูลได้อย่างถูกต้อง

### ข้อควรระวัง
Aggregate functions (`SUM`, `AVG`) จะ ignore NULL

### ข้อดี
แยกความแตกต่างระหว่าง "0" กับ "ไม่มีข้อมูล"

### ข้อเสีย
ทำให้ query ซับซ้อนขึ้น ต้องจัดการ NULL เสมอ

### ข้อห้าม
ห้ามใช้ NULL ในคอลัมน์ที่เป็น primary key (PK ต้อง NOT NULL)

---

## 10. SQL UPDATE

### คืออะไร?
แก้ไขข้อมูลที่มีอยู่แล้วในตาราง

### มีกี่แบบ?
- `UPDATE table SET col = value WHERE condition`
- `UPDATE ... FROM` (join กับตารางอื่น)

### ข้อห้ามสำคัญ
**ห้ามลืม `WHERE`** – ถ้าไม่มีจะอัปเดตทั้งตาราง!

### ใช้อย่างไร
```sql
UPDATE products SET price = price * 1.10 WHERE category = 'Electronics';
UPDATE orders SET status = 'shipped' FROM shipments WHERE orders.id = shipments.order_id;
```

### ประโยชน์
ปรับปรุงข้อมูล

### ข้อควรระวัง
ควร `BEGIN;` แล้ว `COMMIT;` หรือทดสอบด้วย `SELECT` ก่อน update

### ข้อดี
สามารถใช้ `RETURNING` เพื่อดูค่าที่เปลี่ยน
```sql
UPDATE users SET last_login = NOW() WHERE id = 1 RETURNING *;
```

### ข้อเสีย
ถ้าไม่มี index บน `WHERE` จะช้ามาก

### ข้อห้าม
ห้าม update คอลัมน์ที่ใช้ใน `WHERE` โดยไม่ระวัง (อาจทำให้แถวเลื่อน)

---

## 11. SQL DELETE

### คืออะไร?
ลบแถวออกจากตาราง

### มีกี่แบบ?
- `DELETE FROM table WHERE condition`
- `DELETE USING` (join)

### ข้อห้ามสำคัญ
**ห้ามลืม `WHERE`** – มิฉะนั้นข้อมูลทั้งหมดหาย!

### ใช้อย่างไร
```sql
DELETE FROM logs WHERE created_at < NOW() - INTERVAL '30 days';
```

### ประโยชน์
ล้างข้อมูลเก่า

### ข้อควรระวัง
ถ้ามี foreign key อาจถูก reject หรือ cascade ได้

### ข้อดี
`RETURNING` ช่วยให้เห็นข้อมูลที่ถูกลบ

### ข้อเสีย
การลบทีละแถวไม่คืนพื้นที่ว่าง (ต้อง `VACUUM` ใน PostgreSQL)

### ข้อห้าม
ห้าม `DELETE` โดยไม่มี `WHERE` ใน Production เด็ดขาด (ใช้ `TRUNCATE` ถ้าต้องการลบทั้งหมด)

---

## 12. SQL SELECT TOP

### คืออะไร?
จำกัดจำนวนแถวที่คืนกลับ ใน PostgreSQL ใช้ `LIMIT` แทน

### มีกี่แบบ?
- `LIMIT n`
- `LIMIT n OFFSET m`
- `FETCH FIRST n ROWS ONLY` (SQL standard)

### ข้อห้ามสำคัญ
`LIMIT` โดยไม่มี `ORDER BY` จะให้ผลลัพธ์ที่ไม่แน่นอน (random)

### ใช้อย่างไร
```sql
SELECT * FROM sales ORDER BY amount DESC LIMIT 10; -- ขายดี 10 อันดับ
```

### ประโยชน์
ลด traffic, ใช้ทำ pagination

### ข้อควรระวัง
`OFFSET` บนตารางใหญ่ช้ามาก (ต้อง scan ผ่านแถวก่อนหน้า) ให้ใช้ keyset pagination แทน

### ข้อดี
ง่าย

### ข้อเสีย
`OFFSET` มี performance แย่เมื่อ offset มาก

### ข้อห้าม
ห้ามใช้ `LIMIT 1` โดยไม่ระบุ `ORDER BY` ถ้าต้องการ "แถวแรก" ที่แน่นอน

---

## 13. SQL Aggregate Functions

### คืออะไร?
ฟังก์ชันที่คำนวณจากหลายแถวแล้วให้ผลลัพธ์ค่าเดียว

### มีกี่แบบ?
`MIN()`, `MAX()`, `COUNT()`, `SUM()`, `AVG()`, และ `STDDEV()`, `VARIANCE()` เป็นต้น

### ข้อห้ามสำคัญ
Aggregate functions ignore NULL (ยกเว้น `COUNT(*)`)

### ใช้อย่างไร
```sql
SELECT department, AVG(salary), COUNT(*) FROM employees GROUP BY department;
```

### ประโยชน์
สรุปข้อมูล (reporting, KPI)

### ข้อควรระวัง
ถ้าใช้ aggregate โดยไม่มี `GROUP BY` จะรวมทั้งตารางเป็นแถวเดียว

### ข้อดี
ทำงานเร็วถ้ามี index

### ข้อเสีย
ไม่สามารถผสม aggregate กับ non-aggregate โดยไม่ `GROUP BY`

### ข้อห้าม
ห้ามใช้ `COUNT(column)` ถ้าคอลัมน์มี NULL เพราะจะไม่นับ (ใช้ `COUNT(*)` แทน)

---

## 14. SQL MIN(), MAX()

### คืออะไร?
หาค่าต่ำสุดและสูงสุด

### ใช้อย่างไร
```sql
SELECT MIN(price) AS cheapest, MAX(price) AS most_expensive FROM products;
```

### ประโยชน์
หาช่วงของข้อมูล

### ข้อควรระวัง
กับ text, min/max ใช้ lexical order

### ข้อดี
ใช้ index ได้ดี (B-tree)

### ข้อเสีย
ไม่บอกว่าแถวไหนมีค่านั้น (ต้อง subquery)

### ข้อห้าม
ห้ามใช้กับคอลัมน์ที่ไม่มี index บนตารางขนาดใหญ่ ถ้าต้องการความเร็ว

---

## 15. SQL COUNT()

### คืออะไร?
นับจำนวนแถวหรือ non-null values

### มีกี่แบบ?
- `COUNT(*)` – นับทุกแถวรวม NULL
- `COUNT(column)` – นับเฉพาะ non-NULL
- `COUNT(DISTINCT column)`

### ใช้อย่างไร
```sql
SELECT COUNT(DISTINCT customer_id) FROM orders;
```

### ประโยชน์
หาจำนวน record

### ข้อควรระวัง
`COUNT(*)` บนตารางใหญ่ PostgreSQL จะ scan ทั้งตาราง (ใช้ estimate จาก statistic แทน)

### ข้อดี
ง่าย

### ข้อเสีย
ช้ามากถ้าไม่มี index

### ข้อห้าม
ห้ามใช้ `COUNT(*)` ทุกครั้งที่ refresh dashboard (ใช้ materialized view หรือ cache)

---

## 16. SQL SUM(), AVG()

### คืออะไร?
รวมค่า และค่าเฉลี่ย

### ใช้อย่างไร
```sql
SELECT SUM(quantity) AS total_sold, AVG(price) AS avg_price FROM order_items;
```

### ประโยชน์
คํานวณยอดรวม, ค่าเฉลี่ย

### ข้อควรระวัง
ถ้าทุกแถวเป็น NULL, SUM ได้ NULL, AVG ได้ NULL (ใช้ COALESCE)

### ข้อดี
ใช้ index ในการคำนวณ range scan

### ข้อเสีย
ต้องแปลงประเภทถ้า integer แล้วต้องการ decimal

### ข้อห้าม
ห้ามใช้ SUM กับ string

---

## 17. SQL LIKE & Wildcards

### คืออะไร?
ค้นหารูปแบบข้อความ (pattern matching)

### Wildcards
- `%` – แทนอักขระ 0 หรือมากกว่า
- `_` – แทนอักขระเดียว
- ใน PostgreSQL ยังมี `ILIKE` (case-insensitive)

### ข้อห้ามสำคัญ
`LIKE '%text'` หรือ `LIKE '%text%'` ไม่สามารถใช้ B-tree index ได้ (ใช้ trigram index แทน)

### ใช้อย่างไร
```sql
SELECT * FROM users WHERE email LIKE '%@gmail.com';
SELECT * FROM products WHERE name ILIKE '%laptop%';
```

### ประโยชน์
ค้นหาแบบ flexible

### ข้อควรระวัง
Wildcard ขึ้นต้น `%` จะช้ามาก

### ข้อดี
ง่าย เหมาะกับ search ธรรมดา

### ข้อเสีย
ไม่ support regex เต็มรูปแบบ (ต้องใช้ `~` ใน PostgreSQL)

### ข้อห้าม
ห้ามใช้ `LIKE` กับคอลัมน์ที่ไม่มี index และข้อมูลเป็นล้านแถว

---

## 18. SQL IN

### คืออะไร?
ตรวจสอบว่าค่าอยู่ใน list หรือ subquery

### ใช้อย่างไร
```sql
SELECT * FROM employees WHERE dept_id IN (1, 2, 3);
SELECT * FROM orders WHERE customer_id IN (SELECT id FROM customers WHERE vip = true);
```

### ประโยชน์
แทน `OR` ซ้ำๆ

### ข้อควรระวัง
`IN (subquery)` อาจช้าถ้า subquery คืนค่ามาก (ใช้ `EXISTS` แทน)

### ข้อดี
อ่านง่าย

### ข้อเสีย
list ยาวเกินไป (1000+) อาจมี performance issues

### ข้อห้าม
ห้ามใช้ `IN` กับ NULL ใน list (NULL ไม่เท่ากับอะไร)

---

## 19. SQL BETWEEN

### คืออะไร?
ตรวจสอบว่าค่าอยู่ในช่วง [start, end] (รวม endpoint)

### ใช้อย่างไร
```sql
SELECT * FROM products WHERE price BETWEEN 500 AND 1500;
SELECT * FROM orders WHERE order_date BETWEEN '2024-01-01' AND '2024-12-31';
```

### ประโยชน์
อ่านง่ายกว่า `>= AND <=`

### ข้อควรระวัง
สำหรับวันที่ ต้องระวังเวลา (timestamp) – `BETWEEN '2024-01-01' AND '2024-01-31'` จะไม่รวมวันที่ 31 ถ้ามีเวลา 00:00:00? จริงๆ รวมถ้าใช้ date เท่านั้น

### ข้อดี
ใช้ range index scan ได้ดี

### ข้อเสีย
กับ text อาจไม่เป็นไปตามที่คิด

### ข้อห้าม
ห้ามใช้ BETWEEN กับ timestamp โดยไม่แปลงเป็น date ถ้าต้องการ inclusive ขอบบน

---

## 20. SQL Aliases

### คืออะไร?
ตั้งชื่อชั่วคราวให้คอลัมน์หรือตาราง

### มีกี่แบบ?
- คอลัมน์ alias: `SELECT col AS alias`
- ตาราง alias: `FROM table AS t`

### ใช้อย่างไร
```sql
SELECT c.name AS customer_name, o.total
FROM customers AS c
JOIN orders AS o ON c.id = o.cust_id;
```

### ประโยชน์
ทำให้ query สั้นลง อ่านง่าย

### ข้อควรระวัง
alias ของคอลัมน์ใช้ใน `ORDER BY` ได้ แต่ใช้ใน `WHERE` ไม่ได้ (เพราะประมวลผลทีหลัง)

### ข้อดี
ช่วยให้ self-join และ subquery ใช้งานง่าย

### ข้อเสีย
ถ้าใช้มากเกินไปอาจทำให้งง

### ข้อห้าม
ห้ามใช้ alias เดียวกันใน query เดียว (ambiguous)

---

## 21. SQL Joins (รวมทุกชนิด)

### คืออะไร?
รวมข้อมูลจาก 2 ตารางขึ้นไปโดยใช้คอลัมน์ที่เกี่ยวข้อง

### มีกี่แบบ?
- `INNER JOIN` – เอาเฉพาะแถวที่ match ทั้งสองฝั่ง
- `LEFT JOIN` – เอาแถวจากซ้ายทั้งหมด + match ขวา (null ถ้าไม่มี)
- `RIGHT JOIN` – กลับกัน
- `FULL JOIN` – ทั้งสองฝั่ง (union)
- `CROSS JOIN` – cartesian product
- `SELF JOIN` – join ตารางตัวเอง

### ข้อห้ามสำคัญ
ห้าม `JOIN` โดยไม่มี `ON` (กลายเป็น cross join)

### ใช้อย่างไร
```sql
SELECT orders.id, customers.name
FROM orders
INNER JOIN customers ON orders.customer_id = customers.id;
```

### ประโยชน์
Normalized database จำเป็นต้อง join

### ข้อควรระวัง
Join หลายตารางอาจทำให้ query ช้า – ตรวจสอบ execution plan

### ข้อดี
ยืดหยุ่น

### ข้อเสีย
ซับซ้อน อ่านยากถ้าเยอะ

### ข้อห้าม
ห้าม `SELECT *` เวลา join หลายตาราง (คอลัมน์ซ้ำ)

---

## 22. SQL Self Join

### คืออะไร?
Join ตารางกับตัวเอง โดยใช้ alias ต่างกัน

### ใช้อย่างไร
```sql
-- หาพนักงานที่มีหัวหน้าชื่อเดียวกัน
SELECT e1.name AS employee, e2.name AS manager
FROM employees e1
LEFT JOIN employees e2 ON e1.manager_id = e2.id;
```

### ประโยชน์
จัดการ hierarchical data (org chart, category tree)

### ข้อควรระวัง
ต้องใช้ alias เสมอ

### ข้อดี
ไม่ต้องสร้างตารางเพิ่ม

### ข้อเสีย
อาจทำให้เกิด infinite recursion ถ้าไม่ระวัง (ใช้ recursive CTE ดีกว่า)

### ข้อห้าม
ห้าม self join โดยไม่มีเงื่อนไข join ที่ถูกต้อง (จะได้ cross join)

---

## 23. SQL UNION / UNION ALL

### คืออะไร?
รวมผลลัพธ์ของสอง query เข้าด้วยกันเป็นผลลัพธ์ชุดเดียว

### ความแตกต่าง
- `UNION` – ตัดแถวซ้ำ (expensive)
- `UNION ALL` – รวมทั้งหมด (เร็วกว่า)

### ข้อห้ามสำคัญ
คอลัมน์ที่ต้อง match กันทั้งจำนวนและชนิดข้อมูล

### ใช้อย่างไร
```sql
SELECT name FROM active_customers
UNION ALL
SELECT name FROM inactive_customers;
```

### ประโยชน์
รวมข้อมูลจากหลายตารางหรือหลาย query

### ข้อควรระวัง
`ORDER BY` ต้องอยู่ท้ายสุดของคำสั่งรวมเท่านั้น

### ข้อดี
`UNION ALL` เร็ว

### ข้อเสีย
`UNION` ต้อง sort เพื่อตัดซ้ำ

### ข้อห้าม
ห้ามใช้ `UNION` ถ้าไม่จำเป็นต้องตัดซ้ำ (ใช้ `UNION ALL` แทน)

---

## 24. SQL GROUP BY

### คืออะไร?
จัดกลุ่มแถวที่มีค่าเดียวกันในคอลัมน์ที่ระบุ เพื่อใช้ aggregate

### ใช้อย่างไร
```sql
SELECT department, COUNT(*), AVG(salary)
FROM employees
GROUP BY department;
```

### ประโยชน์
สร้างรายงานสรุป

### ข้อควรระวัง
ทุกคอลัมน์ใน `SELECT` ที่ไม่ใช่ aggregate ต้องอยู่ใน `GROUP BY`

### ข้อดี
ลดข้อมูลเหลือกลุ่ม

### ข้อเสีย
อาจไม่สามารถใช้ index ได้ดีถ้า group by หลายคอลัมน์

### ข้อห้าม
ห้าม `GROUP BY` โดยไม่เข้าใจว่าจะเหลือแถวละกลุ่ม

---

## 25. SQL HAVING

### คืออะไร?
กรองผลลัพธ์หลังจาก `GROUP BY` (เหมือน `WHERE` ของกลุ่ม)

### ใช้อย่างไร
```sql
SELECT department, AVG(salary)
FROM employees
GROUP BY department
HAVING AVG(salary) > 50000;
```

### ประโยชน์
กรอง aggregate

### ข้อควรระวัง
`HAVING` ไม่สามารถใช้กับคอลัมน์ที่ไม่ได้ aggregate โดยไม่มี `GROUP BY`

### ข้อดี
ทำให้ report มีเงื่อนไข

### ข้อเสีย
บางครั้งใช้ subquery แทนได้ดีกว่า

### ข้อห้าม
ห้ามใช้ `HAVING` แทน `WHERE` (WHERE เร็วกว่าเพราะ filter ก่อน group)

---

## 26. SQL EXISTS

### คืออะไร?
ตรวจสอบว่า subquery คืนค่าแถวอย่างน้อยหนึ่งแถวหรือไม่

### ใช้อย่างไร
```sql
SELECT * FROM customers c
WHERE EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id);
```

### ประโยชน์
มีประสิทธิภาพดีกว่า `IN` เมื่อ subquery ใหญ่

### ข้อควรระวัง
`SELECT 1` หรือ `SELECT *` ไม่มีผลต่อ EXISTS (แค่ตรวจสอบ)

### ข้อดี
หยุดทำงานทันทีที่พบแถวแรก

### ข้อเสีย
อ่านยากกว่าสำหรับผู้เริ่มต้น

### ข้อห้าม
ห้ามใช้ `EXISTS` กับ subquery ที่ไม่มี correlation (ทำได้ แต่เสียประโยชน์)

---

## 27. SQL ANY / ALL

### คืออะไร?
เปรียบเทียบค่ากับทุกค่าที่คืนจาก subquery
- `ANY` – เป็นจริงถ้ามีค่าใดค่าหนึ่งใน subquery ที่เป็นจริง
- `ALL` – เป็นจริงถ้าทุกค่าใน subquery เป็นจริง

### ใช้อย่างไร
```sql
SELECT name, salary FROM employees
WHERE salary > ANY (SELECT salary FROM employees WHERE department = 'Sales');

SELECT name, salary FROM employees
WHERE salary > ALL (SELECT salary FROM employees WHERE department = 'Intern');
```

### ประโยชน์
ใช้กับ operator `=`, `>`, `<`, `<>`

### ข้อควรระวัง
`= ANY` เหมือน `IN` แต่ `<> ALL` เหมือน `NOT IN`

### ข้อดี
สะดวก

### ข้อเสีย
อาจสับสนระหว่าง `ANY` กับ `SOME`

### ข้อห้าม
ห้ามใช้ `ALL` กับ subquery ที่อาจคืนค่า NULL (จะได้ผลลัพธ์เป็น NULL)

---

## 28. SQL SELECT INTO

### คืออะไร?
สร้างตารางใหม่จากผลลัพธ์ของ query (ใน PostgreSQL ใช้ `CREATE TABLE ... AS` แทน)

### ใช้อย่างไร
```sql
CREATE TABLE high_value_orders AS
SELECT * FROM orders WHERE total > 10000;
```

### ประโยชน์
backup, snapshot, หรือสร้าง staging table

### ข้อควรระวัง
ตารางใหม่จะไม่มี index, constraint, หรือ default

### ข้อดี
เร็วกว่า insert ทีละแถว

### ข้อเสีย
ไม่ copy โครงสร้างทั้งหมด

### ข้อห้าม
ห้ามใช้ใน production ถ้าต้องการ index หรือ fk (ต้องสร้างทีหลัง)

---

## 29. SQL INSERT INTO SELECT

### คืออะไร?
Insert ผลลัพธ์จาก query ลงในตารางที่มีอยู่แล้ว

### ใช้อย่างไร
```sql
INSERT INTO archive_orders (id, date, total)
SELECT id, order_date, total FROM orders WHERE order_date < '2023-01-01';
```

### ประโยชน์
ย้ายข้อมูล, สร้าง report table

### ข้อควรระวัง
จำนวนและชนิดคอลัมน์ต้อง match

### ข้อดี
ทำงานเป็นชุด (set-based)

### ข้อเสีย
ถ้าตารางปลายทางมี trigger หรือ constraint อาจช้า

### ข้อห้าม
ห้าม insert select โดยไม่มี `WHERE` ถ้าไม่ได้ต้องการทั้งหมด

---

## 30. SQL CASE

### คืออะไร?
สร้างเงื่อนไขแบบ if-then-else ใน SQL

### มีกี่แบบ?
- Simple CASE: `CASE col WHEN val1 THEN ...`
- Searched CASE: `CASE WHEN condition THEN ...`

### ใช้อย่างไร
```sql
SELECT name,
  CASE 
    WHEN score >= 80 THEN 'A'
    WHEN score >= 70 THEN 'B'
    ELSE 'C'
  END AS grade
FROM students;
```

### ประโยชน์
ทำ transformation ข้อมูลใน query

### ข้อควรระวัง
CASE จะประเมินตามลำดับ – หยุดเมื่อเจอ true แรก

### ข้อดี
ใช้ใน `SELECT`, `WHERE`, `ORDER BY` ได้

### ข้อเสีย
ซับซ้อนเกินไปอาจอ่านยาก

### ข้อห้าม
ห้ามใช้ CASE เมื่อสามารถใช้ `COALESCE` หรือ `NULLIF` ได้ง่ายกว่า

---

## 31. SQL NULL Functions

### คืออะไร?
ฟังก์ชันจัดการค่า NULL

### ใน PostgreSQL มี
- `COALESCE(val1, val2, ...)` – คืนค่า non-null ตัวแรก
- `NULLIF(a, b)` – คืน NULL ถ้า a = b
- `GREATEST`, `LEAST` – ignore NULL? (ไม่ ถ้ามี NULL คืน NULL)

### ใช้อย่างไร
```sql
SELECT COALESCE(phone, 'ไม่มีเบอร์') FROM contacts;
SELECT NULLIF(amount, 0) FROM payments; -- เปลี่ยน 0 เป็น NULL
```

### ประโยชน์
ป้องกัน error จาก NULL

### ข้อควรระวัง
`COALESCE` มี argument ได้หลายตัว

### ข้อดี
เขียน简洁

### ข้อเสีย
อาจซ่อนข้อมูล null โดยไม่ตั้งใจ

### ข้อห้าม
ห้ามใช้ `COALESCE` บนคอลัมน์ใน `WHERE` เพราะจะไม่ใช้ index (ควร redesign)

---

## 32. SQL Stored Procedures

### คืออะไร?
ชุดคำสั่ง SQL ที่เก็บไว้ใน database และเรียกใช้ได้

### แตกต่างจาก Function?
Procedure ไม่คืนค่า, ควบคุม transaction ได้, ใช้ `CALL`

### ใช้อย่างไร
```sql
CREATE PROCEDURE update_salary(emp_id INT, new_salary DECIMAL)
LANGUAGE plpgsql AS $$
BEGIN
  UPDATE employees SET salary = new_salary WHERE id = emp_id;
  COMMIT;
END; $$;

CALL update_salary(10, 75000);
```

### ประโยชน์
 encapsulate business logic, ลด network round-trip

### ข้อควรระวัง
PostgreSQL 14+ support procedure, ก่อนหน้านี้ใช้ function

### ข้อดี
reuse, ปลอดภัย (กำหนดสิทธิ์)

### ข้อเสีย
debug ยาก, migration ยาก

### ข้อห้าม
ห้ามใช้ procedure สำหรับ logic ที่ควรอยู่ในแอพ (เช่น send email)

---

## 33. SQL Comments

### คืออะไร?
ข้อความที่ database ignore ใช้สำหรับอธิบายโค้ด

### มีกี่แบบ?
- single line: `-- comment`
- multi-line: `/* comment */`

### ใช้อย่างไร
```sql
-- นับจำนวนพนักงาน
SELECT COUNT(*) FROM employees;
```

### ประโยชน์
ทำให้โค้ดเข้าใจง่าย

### ข้อควรระวัง
ไม่มี performance impact

### ข้อดี
document ใน script

### ข้อเสีย
ลืมอัปเดต comment เมื่อโค้ดเปลี่ยน

### ข้อห้าม
ห้าม comment รหัสผ่านหรือ secret

---

## 34. SQL Operators

### คืออะไร?
สัญลักษณ์ที่ใช้ในการดำเนินการ

### มีกี่ประเภท?
- Arithmetic: `+`, `-`, `*`, `/`, `%`
- Comparison: `=`, `!=`, `<`, `>`, `<=`, `>=`
- Logical: `AND`, `OR`, `NOT`
- String: `||` (concatenation)
- Bitwise: `&`, `|`, `#`, `~`

### ใช้อย่างไร
```sql
SELECT (price * 0.07) AS vat FROM products;
SELECT 'Hello ' || 'World';
```

### ประโยชน์
คำนวณ, เปรียบเทียบ, รวม string

### ข้อควรระวัง
`||` กับ NULL ให้ NULL (ใช้ `CONCAT` แทน)

### ข้อดี
เยอะมาก

### ข้อเสีย
บาง operator ต่างกันระหว่าง DB

### ข้อห้าม
ห้ามใช้ `+` กับ string ใน PostgreSQL (ใช้ `||`)

---

## 35. SQL Database (CREATE, DROP, BACKUP)

### CREATE DATABASE
```sql
CREATE DATABASE mydb;
```

### DROP DATABASE
```sql
DROP DATABASE mydb; -- ระวังข้อมูลหาย
```

### BACKUP (ใช้ pg_dump)
```bash
pg_dump mydb > backup.sql
```

### ข้อห้าม
ห้าม DROP DATABASE บน production โดยไม่ backup

---

## 36. SQL CREATE TABLE, DROP TABLE, ALTER TABLE

### CREATE TABLE
```sql
CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT);
```

### DROP TABLE
```sql
DROP TABLE users; -- หรือ CASCADE
```

### ALTER TABLE
```sql
ALTER TABLE users ADD COLUMN email TEXT;
ALTER TABLE users DROP COLUMN email;
ALTER TABLE users RENAME TO customers;
```

### ข้อห้าม
ห้าม DROP TABLE โดยไม่ตรวจสอบ dependencies (foreign keys)

---

## 37. SQL Constraints

### คืออะไร?
กฎที่บังคับข้อมูลในตารางให้ถูกต้อง

### มีกี่แบบ?
- `NOT NULL`
- `UNIQUE`
- `PRIMARY KEY`
- `FOREIGN KEY`
- `CHECK`
- `DEFAULT`
- `EXCLUSION` (PostgreSQL)

### ข้อห้าม
ห้ามสร้าง foreign key โดยไม่มี index บนคอลัมน์อ้างอิง (จะช้า)

---

## 38. SQL Index

### คืออะไร?
โครงสร้างข้อมูลที่ช่วยเพิ่มความเร็วในการค้นหา

### มีกี่แบบ?
- B-tree (default)
- Hash
- GIN (สำหรับ JSON, array)
- GiST, SP-GiST, BRIN

### ใช้อย่างไร
```sql
CREATE INDEX idx_users_email ON users(email);
```

### ข้อห้าม
ห้ามสร้าง index บนคอลัมน์ที่มีการ update บ่อย (overhead)

---

## 39. SQL Auto Increment

### ใน PostgreSQL ใช้ SERIAL หรือ IDENTITY
```sql
CREATE TABLE t (id SERIAL PRIMARY KEY);
-- หรือ
CREATE TABLE t (id INT GENERATED ALWAYS AS IDENTITY);
```

### ข้อห้าม
ห้าม insert ค่าซ้ำใน auto increment column

---

## 40. SQL Dates

### ชนิดข้อมูล
- `DATE`, `TIME`, `TIMESTAMP`, `TIMESTAMPTZ`

### ฟังก์ชัน
```sql
NOW(), CURRENT_DATE, EXTRACT(YEAR FROM date), AGE(date1, date2)
```

### ข้อห้าม
ห้ามเก็บ timestamp เป็น string หรือ integer

---

## 41. SQL Views

### คืออะไร?
ตารางเสมือนที่เก็บ query

### มีกี่แบบ?
- Ordinary view
- Materialized view (เก็บข้อมูลจริง)

### ใช้อย่างไร
```sql
CREATE VIEW active_users AS SELECT * FROM users WHERE status = 'active';
```

### ข้อห้าม
ห้ามใช้ view ซ้อน view หลายชั้น (performance)

---

## 42. SQL Injection

### คืออะไร?
การแทรกคำสั่ง SQL ผ่าน input ที่ไม่ถูก sanitize

### การป้องกัน
- Parameterized query (prepared statement)
- ใช้ ORM
- ใช้ `quote_ident`, `quote_literal` ใน dynamic SQL

### ข้อห้าม
ห้ามต่อ string สร้าง SQL โดยเด็ดขาด

---

## 43. SQL Parameters / Prepared Statements

### คืออะไร?
การแยก SQL structure ออกจากค่าพารามิเตอร์

### ใช้อย่างไร (ในแอพ)
```python
cursor.execute("SELECT * FROM users WHERE id = %s", (user_id,))
```

### ประโยชน์
ป้องกัน SQL injection, เพิ่มประสิทธิภาพ (parse ครั้งเดียว)

### ข้อห้าม
ห้ามใช้ parameter สำหรับชื่อตารางหรือคอลัมน์ (ใช้ identifier quoting แทน)

---

## 44. SQL Hosting

### คืออะไร?
การนำฐานข้อมูล PostgreSQL ไปไว้บน server

### มีกี่แบบ?
- On-premise
- Cloud: AWS RDS, Google Cloud SQL, Azure Database, Supabase, Neon

### ข้อควรระวัง
การตั้งค่า connection pool, SSL, backup, monitoring

### ข้อห้าม
ห้าม expose PostgreSQL port 5432 ตรงสู่ internet โดยไม่มี firewall

---

## สรุปหัวข้อทั้งหมด (Checklist)

✅ ครอบคลุมทุกหัวข้อที่คุณระบุ ตั้งแต่ SQL Query พื้นฐาน, Advanced, จนถึง SQL Hosting
✅ แต่ละหัวข้อมี: ความหมาย, ประเภท, ข้อห้ามสำคัญ, วิธีใช้, ประโยชน์, ข้อดี/เสีย, ข้อห้าม
✅ ยกตัวอย่างจริงใน PostgreSQL
 
 # เฉลยละเอียดเพิ่มเติม + โจทย์ใหม่: MERGE, Table Partitioning, JSON Functions

ตามที่คุณขอ ผมจะจัดทำเป็น 2 ส่วน:
1. **เฉลยละเอียด** สำหรับแบบฝึกหัดบางข้อใน workshop ก่อนหน้า (โดยเฉพาะที่ซับซ้อน)
2. **โจทย์ใหม่** พร้อมเฉลย สำหรับหัวข้อ `MERGE` (UPSERT), `Table Partitioning`, และ `JSON Functions` ใน PostgreSQL

---

## ส่วนที่ 1: เฉลยละเอียดเพิ่มเติม (จาก Workshop เดิม)

### เฉลยข้อ 2.4 – Subquery แบบ Scalar ใน SELECT (พร้อมอธิบาย)

**โจทย์เดิม:** แสดงสินค้าแต่ละตัว พร้อมราคา และราคาเฉลี่ยของสินค้าทั้งหมวด

```sql
SELECT 
    p1.name, 
    p1.price, 
    (SELECT AVG(p2.price) 
     FROM products p2 
     WHERE p2.category = p1.category) AS avg_category_price
FROM products p1;
```

**คำอธิบายเพิ่มเติม:**
- Subquery นี้เรียกว่า **correlated subquery** เพราะอ้างอิงถึง `p1.category` จาก outer query
- ทำงานโดยการรัน subquery **หนึ่งครั้งต่อแต่ละแถว** ของ products
- ถ้าตารางใหญ่ อาจช้า – ควร `LEFT JOIN` กับ aggregated subquery หรือใช้ Window function `AVG() OVER(PARTITION BY category)` ซึ่งมีประสิทธิภาพกว่า

**ตัวอย่างปรับปรุงด้วย Window Function (แนะนำ):**
```sql
SELECT 
    name, 
    price, 
    AVG(price) OVER (PARTITION BY category) AS avg_category_price
FROM products;
```

---

### เฉลยข้อ 3.5 – Recursive CTE (โครงสร้างองค์กร) แบบละเอียด

**โจทย์:** แสดงพนักงานทั้งหมดภายใต้หัวหน้า Alice (id=1) พร้อมระดับชั้น

```sql
WITH RECURSIVE org AS (
    -- Anchor member: หัวหน้าระดับสูงสุด
    SELECT id, name, manager_id, 1 AS level
    FROM employees
    WHERE name = 'Alice'
    
    UNION ALL
    
    -- Recursive member: เอาลูกน้องของแถวก่อนหน้า
    SELECT e.id, e.name, e.manager_id, org.level + 1
    FROM employees e
    JOIN org ON e.manager_id = org.id
)
SELECT * FROM org ORDER BY level, id;
```

**คำอธิบายทีละขั้นตอน:**
1. **Anchor** ดึงแถวของ Alice (level=1)
2. **Recursive** ใช้ผลลัพธ์จาก org (ที่มี level=1) ไป join กับ employees เพื่อหาคนที่มี manager_id = org.id -> ได้ Bob (level=2), Charlie (level=2)
3. ทำซ้ำ: จาก Bob (id=2) หาคนที่มี manager_id=2 -> ได้ David (level=3), จาก Charlie ได้ Eve (level=3)
4. หยุดเมื่อไม่พบแถวใหม่
5. `UNION ALL` รวมผลลัพธ์ทั้งหมด

**ข้อควรระวัง:** ถ้าข้อมูลมีวงจร (เช่น A บังคับ B, B บังคับ A) จะเกิด infinite loop ต้องจำกัด depth หรือใช้ `CYCLE` clause ใน PostgreSQL 14+

---

### เฉลยข้อ 4.3 – Trigger อัปเดต updated_at (พร้อมการทดสอบ)

**สร้างตารางและ trigger อย่างสมบูรณ์:**
```sql
-- เพิ่มคอลัมน์ updated_at (ถ้ายังไม่มี)
ALTER TABLE products ADD COLUMN updated_at TIMESTAMP DEFAULT NOW();

-- สร้างฟังก์ชัน
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- สร้าง trigger
CREATE TRIGGER trigger_products_updated
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- ทดสอบ
UPDATE products SET stock = 20 WHERE id = 1;
SELECT id, name, stock, updated_at FROM products WHERE id = 1;
-- ควรเห็น updated_at เปลี่ยนเป็นเวลาปัจจุบัน
```

**หมายเหตุ:** ถ้าใช้ `FOR EACH STATEMENT` จะไม่สามารถเข้าถึง `NEW`/`OLD` ได้ เพราะอาจมีหลายแถว

---

## ส่วนที่ 2: โจทย์ใหม่ + เฉลย (MERGE, Partition, JSON)

### หัวข้อที่ 1: MERGE (Upsert) – ใช้ใน PostgreSQL 15+

**MERGE** (หรือที่เรียกว่า UPSERT) ช่วยให้สามารถ `INSERT`, `UPDATE`, `DELETE` ในคำสั่งเดียว โดยอิงจากเงื่อนไขการ match

#### โจทย์ที่ 1.1: สร้างหรืออัปเดตข้อมูลสินค้าคงคลัง

สมมติเรามีตาราง `inventory` ที่เก็บ stock ของสินค้าในแต่ละคลัง:

```sql
CREATE TABLE inventory (
    product_id INT,
    warehouse_id INT,
    stock INT,
    last_updated TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (product_id, warehouse_id)
);

-- ข้อมูลตัวอย่าง
INSERT INTO inventory VALUES (1, 101, 10, NOW()), (2, 101, 5, NOW());
```

**โจทย์:** ต้องการอัปเดต stock ถ้ามีคู่ (product_id, warehouse_id) อยู่แล้ว ถ้าไม่มีให้ insert ใหม่ โดยใช้ข้อมูลจากตาราง `new_stock`:

```sql
CREATE TEMP TABLE new_stock AS
SELECT 1 AS product_id, 101 AS warehouse_id, 15 AS stock  -- มีอยู่แล้ว
UNION ALL
SELECT 3, 101, 20;  -- ยังไม่มี
```

**เขียนคำสั่ง MERGE** เพื่อให้ได้ผลลัพธ์: product_id=1 ได้ stock=15, product_id=3 ได้แถวใหม่ stock=20

<details>
<summary>เฉลย</summary>

```sql
MERGE INTO inventory AS target
USING new_stock AS source
ON target.product_id = source.product_id AND target.warehouse_id = source.warehouse_id
WHEN MATCHED THEN
    UPDATE SET stock = source.stock, last_updated = NOW()
WHEN NOT MATCHED THEN
    INSERT (product_id, warehouse_id, stock, last_updated)
    VALUES (source.product_id, source.warehouse_id, source.stock, NOW());
```
</details>

#### โจทย์ที่ 1.2: ใช้ MERGE แบบมี DELETE

**โจทย์:** ถ้า stock ใน source เป็น 0 ให้ลบแถวนั้นออกจาก target (แทนที่จะอัปเดต)

<details>
<summary>เฉลย</summary>

```sql
MERGE INTO inventory AS target
USING new_stock AS source
ON target.product_id = source.product_id AND target.warehouse_id = source.warehouse_id
WHEN MATCHED AND source.stock = 0 THEN
    DELETE
WHEN MATCHED THEN
    UPDATE SET stock = source.stock, last_updated = NOW()
WHEN NOT MATCHED THEN
    INSERT (product_id, warehouse_id, stock, last_updated)
    VALUES (source.product_id, source.warehouse_id, source.stock, NOW());
```
</details>

---

### หัวข้อที่ 2: Table Partitioning (การแบ่งพาร์ติชัน)

**ประโยชน์:** เพิ่มประสิทธิภาพการข้อมูล ตามช่วงเวลา หรือตามคีย์เฉพาะ

#### โจทย์ที่ 2.1: สร้างตาราง orders ที่แบ่งพาร์ติชันตามปี (range partitioning)

**โจทย์:** สร้างตาราง `orders_partitioned` โดยมี partition รายปี (2023, 2024, 2025) และให้แทรกข้อมูลตัวอย่าง จากนั้น query เฉพาะ partition ที่เกี่ยวข้อง

<details>
<summary>เฉลย</summary>

```sql
-- 1. สร้างตารางแม่ (parent table) พร้อมระบุ partition key
CREATE TABLE orders_partitioned (
    id SERIAL,
    order_date DATE NOT NULL,
    customer_id INT,
    total DECIMAL(10,2)
) PARTITION BY RANGE (order_date);

-- 2. สร้าง partitions ย่อย
CREATE TABLE orders_2023 PARTITION OF orders_partitioned
    FOR VALUES FROM ('2023-01-01') TO ('2024-01-01');

CREATE TABLE orders_2024 PARTITION OF orders_partitioned
    FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');

CREATE TABLE orders_2025 PARTITION OF orders_partitioned
    FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');

-- 3. แทรกข้อมูล (จะถูก route ไปยัง partition อัตโนมัติ)
INSERT INTO orders_partitioned (order_date, customer_id, total)
VALUES 
    ('2023-05-10', 1, 1000),
    ('2024-02-20', 2, 2000),
    ('2025-01-15', 3, 1500);

-- 4. Query – จะอ่านเฉพาะ partition ที่เกี่ยวข้อง
EXPLAIN SELECT * FROM orders_partitioned WHERE order_date = '2024-02-20';
-- ควรเห็น "Seq Scan on orders_2024" เท่านั้น

-- 5. ดู partition ที่มีข้อมูล
SELECT tableoid::regclass, count(*) FROM orders_partitioned GROUP BY tableoid;
```
</details>

#### โจทย์ที่ 2.2: เพิ่ม partition ใหม่สำหรับปี 2026

<details>
<summary>เฉลย</summary>

```sql
CREATE TABLE orders_2026 PARTITION OF orders_partitioned
    FOR VALUES FROM ('2026-01-01') TO ('2027-01-01');
```
</details>

---

### หัวข้อที่ 3: JSON Functions ใน PostgreSQL

PostgreSQL มีฟังก์ชัน JSON มากมาย เช่น `jsonb_build_object`, `jsonb_agg`, `->`, `->>`, `@>`, `jsonb_set`

#### โจทย์ที่ 3.1: สร้างตารางที่มีคอลัมน์ JSONB และ query หาข้อมูล

**โจทย์:** สร้างตาราง `user_profiles` ที่มีคอลัมน์ `metadata` เป็น JSONB เก็บข้อมูลเพิ่มเติม (เช่น ที่อยู่, เบอร์โทร) แล้วเขียน query เพื่อหาผู้ใช้ที่อาศัยอยู่ใน "กรุงเทพ"

<details>
<summary>เฉลย</summary>

```sql
CREATE TABLE user_profiles (
    id SERIAL PRIMARY KEY,
    name TEXT,
    metadata JSONB
);

INSERT INTO user_profiles (name, metadata) VALUES
('สมชาย', '{"city": "กรุงเทพ", "phone": "0812345678"}'),
('สมหญิง', '{"city": "เชียงใหม่", "phone": "0898765432"}');

-- Query หาคนที่ city = กรุงเทพ
SELECT name, metadata->>'phone' AS phone
FROM user_profiles
WHERE metadata->>'city' = 'กรุงเทพ';

-- หรือใช้ @> operator (ถ้าเป็น JSONB)
SELECT * FROM user_profiles WHERE metadata @> '{"city": "กรุงเทพ"}';
```
</details>

#### โจทย์ที่ 3.2: อัปเดตค่าใน JSONB (เพิ่มฟิลด์ใหม่)

**โจทย์:** เพิ่มฟิลด์ `"age": 30` ใน metadata ของ user ที่ชื่อ "สมชาย"

<details>
<summary>เฉลย</summary>

```sql
UPDATE user_profiles
SET metadata = jsonb_set(metadata, '{age}', '30')
WHERE name = 'สมชาย';

-- หรือใช้ || operator (concat)
UPDATE user_profiles
SET metadata = metadata || '{"age": 30}'
WHERE name = 'สมชาย';
```
</details>

#### โจทย์ที่ 3.3: รวม JSONB จากหลายแถวเป็น array

**โจทย์:** สร้าง query ที่คืนค่าเป็น JSON array ของชื่อและเมืองทั้งหมด

<details>
<summary>เฉลย</summary>

```sql
SELECT jsonb_agg(jsonb_build_object('name', name, 'city', metadata->>'city'))
FROM user_profiles;
```
</details>

#### โจทย์ที่ 3.4 (ขั้นสูง): JSONB พร้อม index GIN เพื่อความเร็ว

**โจทย์:** สร้าง index บน metadata แล้วทดสอบ performance

<details>
<summary>เฉลย</summary>

```sql
CREATE INDEX idx_user_metadata ON user_profiles USING GIN (metadata);

-- จากนั้น query ที่ใช้ @> หรือ ? จะใช้ index
EXPLAIN SELECT * FROM user_profiles WHERE metadata @> '{"city": "กรุงเทพ"}';
```
</details>

---

## สรุปคำขอที่ทำครบถ้วน

✅ **เฉลยละเอียดเพิ่มเติม** 3 ข้อ (Scalar subquery, Recursive CTE, Trigger)  
✅ **โจทย์ใหม่ MERGE** 2 ข้อ พร้อมเฉลย  
✅ **โจทย์ใหม่ Table Partitioning** 2 ข้อ พร้อมเฉลย  
✅ **โจทย์ใหม่ JSON Functions** 4 ข้อ พร้อมเฉลย

