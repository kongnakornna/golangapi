นี่คือ **คู่มือ SQL สำหรับ PostgreSQL** ที่ครอบคลุมทุกหัวข้อที่คุณระบุ พร้อมตัวอย่างการใช้งานจริง การออกแบบฟังก์ชัน, Subquery, CASE, Stored Procedures, การป้องกัน SQL Injection และส่วนของการออกแบบ Workflow, Task List, Checklist ตามที่ต้องการ

---

# คู่มือ SQL สำหรับ PostgreSQL ฉบับสมบูรณ์ (Practical Guide)

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

หากต้องการให้ขยายหัวข้อใดเป็นบทความแยก หรือเขียนตัวอย่างเพิ่ม (เช่น Trigger, Window Functions, CTE) แจ้งได้เลยครับ