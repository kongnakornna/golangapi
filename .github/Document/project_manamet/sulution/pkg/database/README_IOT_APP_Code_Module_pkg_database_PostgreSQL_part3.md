A subquery is a powerful technique where one query (`SELECT`, `INSERT`, `UPDATE`, or `DELETE`) is nested inside another[reference:0], allowing you to create complex, multi-step logic. However, it's important to know when to use a subquery versus a `JOIN` to optimize your application's performance.

### 📝 ประเภทและโครงสร้างของ Subquery
*   **Basic Subquery**: Subqueries are enclosed in parentheses and can be placed in many parts of an SQL statement[reference:1], typically used with operators like `IN`, `EXISTS`, `ANY`, or `ALL`[reference:2].
*   **Scalar Subquery**: Must return exactly one row and one column, often used in `SELECT` or `WHERE` clauses to compare with a single value. If a scalar subquery returns no rows, the result is `NULL`[reference:3].
*   **Correlated Subquery**: References a column from the outer query, executed repeatedly for each row processed by the outer query[reference:4]. It's powerful but must be used carefully due to performance implications[reference:5].

### 📈 Subquery กับ JOIN
*   **Performance**: Uncorrelated subqueries are generally not a problem. However, **correlated scalar subqueries** often cause performance issues because PostgreSQL can only execute them as a nested loop, running the subquery once for each row in the outer query[reference:6].
*   **Solution**: Rewriting a correlated scalar subquery as a `JOIN` is often a better solution[reference:7], especially when one row from the main table matches many rows in the subquery.
*   **Example**: To find users with a profile from a specific domain, a `JOIN` query is often more elegant and efficient than a subquery[reference:8].

### 💡 GORM in Action: `Raw()` vs. Chainable API
While GORM's chainable API is great, `Raw()` is often the simpler, more readable, and less error-prone choice for non-trivial queries.

*   **When to use `Raw()`**: For complex `FROM` clauses, when performance is critical, or when the Chainable API becomes cumbersome.
*   **When to use Chainable API**: For simple subqueries, dynamic queries built at runtime, and in teams that prioritize ORM consistency.

### 📋 GORM Subquery Implementations

| การใช้งาน | Chainable API | Raw SQL (`db.Raw()`) |
| :--- | :--- | :--- |
| **WHERE IN** | `db.Where("amount > (?)", db.Table("orders").Select("AVG(amount)")).Find(&orders)` | `db.Raw("SELECT * FROM orders WHERE amount > (SELECT AVG(amount) FROM orders)").Scan(&orders)` |
| **FROM Clause** | `subQuery := db.Table("users").Select("name").Where("age > ?", 18).SubQuery()`<br>`db.Table(subQuery).Find(&results)` | `db.Raw("SELECT * FROM (?) AS t WHERE ...", subQuerySQL)` |
| **UPDATE** | `subQuery := db.Table("table_b").Select("col_c").Where("col_d = ?", value)`<br>`db.Table("table_a").Where("col_c = (?)", subQuery).Updates(data)` | `db.Exec("UPDATE table_a SET col_a = ? WHERE col_c = (SELECT col_c FROM table_b WHERE col_d = ?)", valueA, valueD)` |

#### ✍️ แบบฝึกหัด
1.  **Subquery in `WHERE`**: Write a query using a subquery to find all products whose price is greater than the average price of all products in the same category.
2.  **Subquery with `EXISTS`**: Write a query to find all customers who have placed at least one order in the last 30 days.
3.  **GORM Chainable Subquery**: Write a GORM query to find all users whose age is greater than the average age of all users.
4.  **GORM `Raw` Subquery**: Write a GORM `Raw` query to find the top 5 most expensive products by joining with a subquery that calculates the average price per category.
5.  **Refactoring**: Refactor a correlated scalar subquery that finds the latest order for each user into a more efficient `JOIN` query.

---

Do you have a specific use case or scenario where you're trying to decide between a subquery and a `JOIN`? If you can share the Go structs and the logic you have in mind, I can help you write the most suitable and efficient code.

## การจัดการ Query ที่ซับซ้อนใน Golang (30 ตัวอย่าง)

ในโลกจริง เรามักเจอ SQL ที่ซับซ้อน เช่น การรวมหลายตาราง, window functions, recursive CTE, pivot, หรือการ query ข้อมูล JSON ใน PostgreSQL การใช้ GORM แบบ Chainable API อาจไม่ครอบคลุมทุกกรณี ดังนั้นวิธีที่ดีที่สุดคือใช้ **Raw SQL** ผสมกับ GORM หรือใช้ `sqlx` สำหรับงานที่ต้องการประสิทธิภาพสูงสุด

ด้านล่างนี้คือ **30 ตัวอย่าง** ของ complex queries ที่ใช้บ่อย พร้อมโค้ด Go ที่รันได้จริง (ใช้ GORM + PostgreSQL)

---

### ข้อควรรู้ก่อนเริ่ม
- ใช้ `db.Raw()` เมื่อ query ซับซ้อนเกินความสามารถของ GORM chain
- ใช้ `db.Exec()` สำหรับ DML ที่ซับซ้อน (UPDATE/DELETE พร้อม subquery)
- ใช้ `Scan()` หรือ `Find()` เพื่อ map ผลลัพธ์ไปยัง struct
- หาก query ส่งคืนโครงสร้างที่ไม่ตรงกับ model ให้สร้าง struct ชั่วคราว

---

## 30 ตัวอย่าง Query ซับซ้อน

### 1. **Window Function: ROW_NUMBER() เพื่อจัดอันดับภายในกลุ่ม**
```go
type RankedProduct struct {
    ID        uint
    Name      string
    Category  string
    Price     float64
    Rank      int
}

query := `
    SELECT id, name, category, price,
           ROW_NUMBER() OVER (PARTITION BY category ORDER BY price DESC) AS rank
    FROM products
`
var results []RankedProduct
db.Raw(query).Scan(&results)
```

### 2. **Running Total (สะสมยอดขาย)**
```go
type SalesRunning struct {
    Date      time.Time
    DailySale float64
    Running   float64
}

db.Raw(`
    SELECT date, amount AS daily_sale,
           SUM(amount) OVER (ORDER BY date) AS running
    FROM sales
`).Scan(&results)
```

### 3. **Recursive CTE (จัดการกราฟ / tree structure)**
```go
// หาลูกหลานทั้งหมดของ department_id = 1
var deptIDs []uint
db.Raw(`
    WITH RECURSIVE dept_tree AS (
        SELECT id FROM departments WHERE parent_id = 1
        UNION ALL
        SELECT d.id FROM departments d
        JOIN dept_tree dt ON dt.id = d.parent_id
    )
    SELECT id FROM dept_tree
`).Scan(&deptIDs)
```

### 4. **LATERAL JOIN (เรียก subquery แบบ row-by-row)**
```go
type UserWithLatestOrder struct {
    UserID    uint
    UserName  string
    OrderID   uint
    OrderDate time.Time
}

db.Raw(`
    SELECT u.id AS user_id, u.name, o.id AS order_id, o.created_at
    FROM users u
    LEFT JOIN LATERAL (
        SELECT id, created_at FROM orders
        WHERE user_id = u.id
        ORDER BY created_at DESC LIMIT 1
    ) o ON true
`).Scan(&results)
```

### 5. **Pivot Table (แปลงแถวเป็นคอลัมน์)**
```go
type MonthlySales struct {
    Year  int
    Jan   float64
    Feb   float64
    Mar   float64
    -- ...
}

db.Raw(`
    SELECT year,
           SUM(CASE WHEN month = 1 THEN amount END) AS Jan,
           SUM(CASE WHEN month = 2 THEN amount END) AS Feb,
           SUM(CASE WHEN month = 3 THEN amount END) AS Mar
    FROM sales
    GROUP BY year
`).Scan(&results)
```

### 6. **JSON Aggregation (รวมแถวเป็น JSON array)**
```go
type CategoryWithProducts struct {
    Category   string
    ProductList string // JSON string
}

db.Raw(`
    SELECT category,
           json_agg(json_build_object('id', id, 'name', name)) AS product_list
    FROM products
    GROUP BY category
`).Scan(&results)
```

### 7. **ค้นหาในคอลัมน์ JSONB**
```go
var users []User
db.Raw(`
    SELECT * FROM users
    WHERE profile->>'nickname' = 'john_doe'
      AND profile->'address'->>'city' = 'Bangkok'
`).Scan(&users)
```

### 8. **ใช้ @> (contains) กับ JSONB array**
```go
db.Raw(`
    SELECT * FROM products
    WHERE tags @> '["electronic", "sale"]'
`).Scan(&products)
```

### 9. **Array Operators (ANY / ALL)**
```go
var ids []uint
db.Raw(`
    SELECT id FROM users
    WHERE age = ANY(ARRAY[25,30,35])
`).Scan(&ids)
```

### 10. **Range Types (overlap, contain)**
```go
db.Raw(`
    SELECT * FROM reservations
    WHERE period && '[2025-01-01, 2025-01-07]'::daterange
`).Scan(&reservations)
```

### 11. **Full-Text Search (tsvector / tsquery)**
```go
var articles []Article
db.Raw(`
    SELECT * FROM articles
    WHERE to_tsvector('english', title || ' ' || content) @@ to_tsquery('english', 'database & performance')
`).Scan(&articles)
```

### 12. **Grouping Sets / ROLLUP (สรุปหลายมิติ)**
```go
type SalesSummary struct {
    Region    string
    Product   string
    Total     float64
}

db.Raw(`
    SELECT region, product, SUM(amount)
    FROM sales
    GROUP BY ROLLUP(region, product)
`).Scan(&results)
```

### 13. **Filtered Aggregates (เฉพาะบางแถว)**
```go
db.Raw(`
    SELECT
        COUNT(*) FILTER (WHERE status = 'active') AS active_count,
        COUNT(*) FILTER (WHERE status = 'inactive') AS inactive_count
    FROM users
`).Scan(&stats)
```

### 14. **DISTINCT ON (เลือกแถวแรกในกลุ่ม)**
```go
// ออเดอร์ล่าสุดของแต่ละ user
db.Raw(`
    SELECT DISTINCT ON (user_id) *
    FROM orders
    ORDER BY user_id, created_at DESC
`).Scan(&latestOrders)
```

### 15. **Common Table Expression (CTE) พร้อมการอัปเดต**
```go
db.Exec(`
    WITH moved AS (
        DELETE FROM logs WHERE created_at < NOW() - INTERVAL '90 days'
        RETURNING *
    )
    INSERT INTO logs_archive SELECT * FROM moved
`)
```

### 16. **UPDATE จากผลลัพธ์ของ JOIN**
```go
db.Exec(`
    UPDATE products p
    SET price = p.price * 0.9
    FROM categories c
    WHERE p.category_id = c.id AND c.name = 'clearance'
`)
```

### 17. **DELETE โดยใช้ USING (PostgreSQL specific)**
```go
db.Exec(`
    DELETE FROM users u
    USING inactive_users i
    WHERE u.id = i.user_id
`)
```

### 18. **MERGE (Upsert) แบบมีเงื่อนไข**
```go
db.Exec(`
    INSERT INTO inventory (product_id, quantity)
    VALUES (1, 100)
    ON CONFLICT (product_id) DO UPDATE
    SET quantity = inventory.quantity + EXCLUDED.quantity
    WHERE inventory.quantity < 500
`)
```

### 19. **Recursive Query หาเส้นทางในกราฟ (BFS)**
```go
type Path struct {
    StartID uint
    EndID   uint
    Depth   int
}

db.Raw(`
    WITH RECURSIVE path AS (
        SELECT start_id, end_id, 1 AS depth FROM connections WHERE start_id = 1
        UNION
        SELECT p.start_id, c.end_id, p.depth + 1
        FROM path p
        JOIN connections c ON p.end_id = c.start_id
        WHERE p.depth < 5
    )
    SELECT * FROM path
`).Scan(&paths)
```

### 20. **Window Function: LAG / LEAD (เข้าถึงแถวก่อน/หลัง)**
```go
db.Raw(`
    SELECT date, amount,
           LAG(amount, 1) OVER (ORDER BY date) AS prev_day,
           amount - LAG(amount, 1) OVER (ORDER BY date) AS diff
    FROM daily_sales
`).Scan(&results)
```

### 21. **Percent Rank และ Cume Dist**
```go
db.Raw(`
    SELECT name, salary,
           PERCENT_RANK() OVER (ORDER BY salary) AS percentile
    FROM employees
`).Scan(&results)
```

### 22. **NTILE (แบ่งเป็น N buckets)**
```go
db.Raw(`
    SELECT name, score,
           NTILE(4) OVER (ORDER BY score) AS quartile
    FROM exam_results
`).Scan(&results)
```

### 23. **UNPIVOT (แปลงคอลัมน์เป็นแถว) ด้วย VALUES**
```go
db.Raw(`
    SELECT id, unnest(ARRAY['q1','q2','q3']) AS quarter,
           unnest(ARRAY[q1_sales, q2_sales, q3_sales]) AS sales
    FROM yearly_sales
`).Scan(&results)
```

### 24. **Pivot แบบ Dynamic (สร้าง SQL แบบข้อความก่อน)**
```go
// สมมติมี list of months
months := []string{"Jan", "Feb", "Mar"}
pivotCols := strings.Join(months, ",")
sql := fmt.Sprintf(`
    SELECT year, %s
    FROM crosstab(
        'SELECT year, month, amount FROM sales ORDER BY 1,2',
        'SELECT unnest(ARRAY[''%s''])'
    ) AS ct(year int, %s numeric)
`, pivotCols, strings.Join(months, "','"), pivotCols)

db.Raw(sql).Scan(&result)
```

### 25. **Hypothetical-set function (rank with tie)**
```go
var rank int
db.Raw(`
    SELECT rank(70) WITHIN GROUP (ORDER BY score DESC)
    FROM students
`).Scan(&rank)
```

### 26. **GROUP BY ตามช่วงเวลา (time bucket)**
```go
db.Raw(`
    SELECT date_trunc('hour', created_at) AS hour,
           COUNT(*)
    FROM events
    GROUP BY hour
    ORDER BY hour
`).Scan(&hourlyStats)
```

### 27. **ค้นหาผู้ใช้ที่ทำรายการซ้ำ (self-join)**
```go
var dupUsers []User
db.Raw(`
    SELECT DISTINCT u1.*
    FROM users u1
    JOIN users u2 ON u1.email = u2.email AND u1.id != u2.id
`).Scan(&dupUsers)
```

### 28. **Fuzzy string matching (pg_trgm)**
```go
db.Raw(`
    SELECT *, similarity(name, 'Jonn') AS sml
    FROM users
    WHERE name % 'Jonn'
    ORDER BY sml DESC
`).Scan(&users)
```

### 29. **Partitioned Table Query (เฉพาะ partition)**
```go
db.Raw(`
    SELECT * FROM orders_2025_01
    WHERE order_date BETWEEN '2025-01-01' AND '2025-01-31'
`)
```

### 30. **EXPLAIN ANALYZE เพื่อ debug performance**
```go
var plan string
db.Raw("EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) " + complexSQL).Scan(&plan)
fmt.Println(plan) // log เอาไปวิเคราะห์
```

---

## วิธีการใช้ GORM ร่วมกับ complex query

**1. ใช้ `Raw()` + `Scan()`** (เหมาะกับ query ที่ไม่เปลี่ยนโครงสร้าง)
```go
type TempResult struct {
    Col1 string
    Col2 int
}
db.Raw("SELECT ...").Scan(&temp)
```

**2. ใช้ `Table()` + `Select()` + `Joins()`** (ถ้าเป็น window function หรือ CTE ที่ไม่ซับซ้อนเกิน)
```go
subQuery := db.Table("orders").Select("user_id, SUM(amount) as total").Group("user_id")
db.Table("users u").
   Select("u.*, sq.total").
   Joins("JOIN (?) sq ON u.id = sq.user_id", subQuery).
   Find(&users)
```

**3. สร้าง `Scope` สำหรับ complex condition ที่ reuse บ่อย**
```go
func FilterByFullText(search string) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("to_tsvector(title) @@ plainto_tsquery(?)", search)
    }
}
```

**4. ใช้ `db.Exec()` สำหรับ DML ที่ซับซ้อน**
```go
db.Exec(`
    WITH updated AS (
        UPDATE products SET price = price * 1.1
        WHERE category_id IN (SELECT id FROM categories WHERE on_sale = true)
        RETURNING id
    )
    INSERT INTO price_logs (product_id, new_price, changed_at)
    SELECT id, price, NOW() FROM products WHERE id IN (SELECT id FROM updated)
`)
```

---

## ข้อควรระวังและข้อห้าม

- **ห้ามใช้ Bucket Pattern ร่วมกับ Time Series Collections** (ตามเอกสารอ้างอิงก่อนหน้า) เพราะจะลดประสิทธิภาพการ query
- หลีกเลี่ยง correlated subquery ใน `SELECT` หากสามารถ改用 `JOIN` หรือ `LATERAL` ได้
- ระวัง SQL Injection: **ห้ามใช้ `fmt.Sprintf` ใส่ค่าที่รับจาก user โดยตรง** ให้ใช้ `db.Raw(sql, param1, param2)` แทน
- หาก query ส่งคืนแถวเป็นล้าน ให้ใช้ `Rows()` + `Scan()` แบบ stream แทน `Scan(&slice)`

---

## สรุป

การจัดการ complex query ใน Golang ด้วย PostgreSQL:
- **30 ตัวอย่าง** ข้างต้นครอบคลุมกรณีการใช้งานจริงส่วนใหญ่
- เลือกใช้ **Raw SQL** เมื่อ query มี window function, CTE, pivot, หรือ JSON operators
- ใช้ **GORM chain API** เฉพาะส่วนที่ dynamic และไม่ซับซ้อนเกิน
- เสริมประสิทธิภาพด้วย `EXPLAIN` และปรับ index ให้เหมาะสม

