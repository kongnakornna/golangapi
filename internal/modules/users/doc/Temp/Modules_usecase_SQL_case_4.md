# สารบัญ

1. [หลักการทำงาน SQL แบบ Advance](#หลักการทำงาน-sql-แบบ-advance)
2. [ส่วนที่ 1: SQL Advance บน GORM (GoLang)](#ส่วนที่-1-sql-advance-บน-gorm-golang)
3. [เอกสารประกอบแต่ละ Component](#เอกสารประกอบแต่ละ-component)

---

# หลักการทำงาน SQL แบบ Advance

## คืออะไร?

SQL ขั้นสูง (Advanced SQL) คือเทคนิคการเขียนคำสั่ง SQL ที่ซับซ้อนกว่า CRUD พื้นฐาน เพื่อจัดการกับ:
- **Dynamic filtering** - การสร้างเงื่อนไข WHERE แบบพลวัต
- **Conditional logic** - การใช้ CASE, COALESCE, NULLIF
- **Recursive queries** - การ query ข้อมูลแบบลำดับชั้น
- **Window functions** - การคำนวณข้ามแถวข้อมูล
- **JSON/JSONB operations** - การจัดการข้อมูลแบบ JSON
- **CTE (Common Table Expressions)** - การสร้าง temporary result sets

## มีกี่แบบ?

| # | ประเภท | คำอธิบาย | ตัวอย่างการใช้งาน |
|---|--------|----------|------------------|
| 1 | **Conditional Expressions** | CASE, COALESCE, NULLIF | แปลงค่า, จัดการ NULL |
| 2 | **Dynamic WHERE** | Conditional filtering | ระบบค้นหาพร้อม filter หลายตัว |
| 3 | **Window Functions** | ROW_NUMBER, RANK, LAG, LEAD | Ranking, การเปรียบเทียบแถวก่อนหน้า |
| 4 | **CTE (WITH clause)** | WITH ... AS (...) | ทำให้ query ซับซ้อนอ่านง่าย |
| 5 | **Recursive CTE** | WITH RECURSIVE |  query ข้อมูล tree/hierarchy |
| 6 | **JSON Operations** | ->, ->>, jsonb_agg | เก็บ/ query ข้อมูล semi-structured |
| 7 | **Aggregate Functions** | STRING_AGG, ARRAY_AGG | รวมหลายแถวเป็นค่าเดียว |
| 8 | **Subquery** | EXISTS, IN, ANY, ALL | query ซ้อน query |
| 9 | **Full-Text Search** | tsvector, tsquery | ค้นหาข้อความประสิทธิภาพสูง |
| 10 | **Partitioning** | PARTITION BY | แบ่งตารางใหญ่เป็นส่วนย่อย |

## ข้อห้ามสำคัญ

```
⚠️ CRITICAL RULES:

1. ห้ามใช้ loop ใน SQL ถ้าใช้ set-based operation ได้ (performance ตกมาก)
2. ห้าม SELECT * ใน production (ใช้เฉพาะ column ที่จำเป็น)
3. ห้ามทำ dynamic SQL โดยไม่ sanitize input (SQL Injection)
4. ห้ามใช้ recursive query กับข้อมูลมากๆ (stack overflow)
5. ห้าม join เกิน 5-7 ตารางใน query เดียว
6. ห้ามใช้ functions ใน WHERE clause ถ้า column ถูก index
7. ห้ามทำ aggregation บนข้อมูลที่ไม่ได้ filter ก่อน
```

---

# ส่วนที่ 1: SQL Advance บน GORM (GoLang)

## โครงสร้างโปรเจกต์

```
advanced-sql-gorm/
├── README.md
├── go.mod
├── go.sum
├── main.go
├── config/
│   └── database.go
├── models/
│   └── user.go
├── repositories/
│   ├── user_repo.go
│   ├── filter_repo.go
│   ├── window_repo.go
│   ├── recursive_repo.go
│   └── json_repo.go
├── services/
│   ├── user_service.go
│   └── filter_service.go
├── controllers/
│   └── user_controller.go
├── middleware/
│   └── logger.go
├── utils/
│   └── query_builder.go
├── tests/
│   ├── user_repo_test.go
│   └── integration_test.go
├── migrations/
│   ├── 001_create_users_table.sql
│   ├── 002_create_employee_tree.sql
│   └── 003_add_json_metadata.sql
└── scripts/
    ├── seed_data.go
    └── benchmark.go
```

## คำอธิบายแต่ละโฟลเดอร์/ไฟล์

| โฟลเดอร์/ไฟล์ | คำอธิบาย (ไทย) | Description (English) |
|---------------|----------------|----------------------|
| `main.go` | จุดเริ่มต้นของโปรแกรม, เรียกใช้ service และ route | Application entry point, initializes services and routes |
| `config/` | ตั้งค่าการเชื่อมต่อฐานข้อมูล PostgreSQL | Database connection configuration for PostgreSQL |
| `models/` | โครงสร้างข้อมูลสอดคล้องกับตารางใน DB | Data structures mapping to database tables |
| `repositories/` | รวม logic การ query ฐานข้อมูลทั้งหมด | Contains all database query logic |
| `services/` | ประมวลผล business logic ก่อนเรียก repository | Processes business logic before calling repositories |
| `controllers/` | รับ request และส่ง response กลับ client | Handles HTTP requests and responses |
| `middleware/` | ฟังก์ชันทำงานก่อน/หลัง request (logging, auth) | Functions executed before/after requests |
| `utils/` | ฟังก์ชันช่วยเหลือที่ใช้ร่วมกันได้ | Reusable helper functions |
| `tests/` | ชุดทดสอบสำหรับ repository และ integration | Test suites for repositories and integration |
| `migrations/` | SQL scripts สำหรับสร้างและปรับปรุง schema | SQL scripts for schema creation and updates |
| `scripts/` | สคริปต์เสริมสำหรับ seed data และ benchmark | Utility scripts for data seeding and benchmarking |

---

# เอกสารประกอบ: repositories/user_repo.go

## Concept
### คืออะไร?
Repository layer ที่รวบรวมการทำงานกับฐานข้อมูลทั้งหมด ใช้ GORM เป็น ORM

### มีกี่แบบ?
- Basic CRUD (Create, Read, Update, Delete)
- Complex query with filters
- Transaction operations
- Raw SQL execution

## Comment CODE (ไทย/อังกฤษ คนละบรรทัด)

```go
// CreateUser - สร้างผู้ใช้ใหม่ในระบบ
// CreateUser - Creates a new user in the system
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
    // ใช้ WithContext เพื่อ support context cancellation และ timeout
    // Use WithContext to support context cancellation and timeout
    return r.DB.WithContext(ctx).Create(user).Error
}

// GetUserByID - ค้นหาผู้ใช้ด้วย ID พร้อม preload ความสัมพันธ์
// GetUserByID - Finds user by ID with preloaded relationships
func (r *UserRepository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
    var user models.User
    
    // Preload ช่วยลด N+1 query problem
    // Preload helps reduce N+1 query problem
    err := r.DB.WithContext(ctx).
        Preload("Roles").           // โหลดข้อมูล roles ที่เกี่ยวข้อง
        Preload("Profile").         // โหลดข้อมูล profile ที่เกี่ยวข้อง
        Where("id = ? AND deleted_at IS NULL", id).
        First(&user).Error
    
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

### ใช้อย่างไร / นำไปใช้กรณีไหน
```go
// ตัวอย่างการเรียกใช้ใน service layer
// Example usage in service layer
func (s *UserService) GetUser(ctx context.Context, id uint) (*models.User, error) {
    // validate input ก่อนเรียก repository
    if id == 0 {
        return nil, errors.New("invalid user id")
    }
    return s.userRepo.GetUserByID(ctx, id)
}
```

### ประโยชน์ที่ได้รับ
- แยกการเข้าถึงฐานข้อมูลออกจาก business logic
- สามารถ mock repository ใน unit test ได้
- ลด code duplication

### ข้อควรระวัง
- ระวังการเรียก Preload มากเกินไป (over-fetching)
- ต้องจัดการ context timeout ให้เหมาะสม
- ระวัง N+1 query เมื่อใช้ loop ร่วมกับ database call

### ข้อดี
- ควบคุม database operation ได้จากที่เดียว
- ง่ายต่อการทำ caching
- สามารถเปลี่ยน ORM ได้โดยไม่กระทบ service layer

### ข้อเสีย
- มี boilerplate code เยอะ
- ต้องเขียน repository ทุก entity
- อาจซับซ้อนเกินไปสำหรับโปรเจกต์เล็ก

### ข้อห้าม
- ห้ามใส่ business logic ใน repository
- ห้ามเรียก repository ซ้อน repository โดยตรง
- ห้าม return *gorm.DB ออกจาก repository

---

# เอกสารประกอบ: repositories/filter_repo.go

## Concept
### คืออะไร?
Dynamic filter system ที่สร้าง WHERE clause ตามเงื่อนไขที่ได้รับจาก client

### มีกี่แบบ?
1. Map-based filtering
2. Struct-based filtering with omitempty
3. Raw SQL with conditional builder
4. Scopes pattern (GORM specific)

## Comment CODE

```go
// FilterUsers - กรองผู้ใช้ตามเงื่อนไขที่ส่งมาแบบ dynamic
// FilterUsers - Dynamically filters users based on provided conditions
func (r *UserRepository) FilterUsers(ctx context.Context, filters map[string]interface{}) ([]models.User, int64, error) {
    var users []models.User
    var total int64
    
    // เริ่มต้น query ด้วย model
    // Start query with base model
    query := r.DB.WithContext(ctx).Model(&models.User{})
    
    // ========== DYNAMIC WHERE CLAUSE ==========
    // ตรวจสอบและเพิ่มเงื่อนไขทีละฟิลด์
    // Check and add conditions field by field
    
    // กรณีค้นหาด้วยชื่อแบบ partial match (case-insensitive)
    // Case-insensitive partial match for name search
    if name, ok := filters["name"].(string); ok && name != "" {
        // ILIKE ใน PostgreSQL, GORM จะแปลงเป็น LIKE สำหรับ DB อื่น
        // ILIKE in PostgreSQL, GORM converts for other databases
        query = query.Where("name ILIKE ?", "%"+name+"%")
    }
    
    // กรณีกําหนด email แบบ exact match
    // Exact match for email
    if email, ok := filters["email"].(string); ok && email != "" {
        query = query.Where("email = ?", email)
    }
    
    // กรณีกําหนด age range (min-max)
    // Age range filtering
    if minAge, ok := filters["min_age"].(int); ok && minAge > 0 {
        query = query.Where("age >= ?", minAge)
    }
    if maxAge, ok := filters["max_age"].(int); ok && maxAge > 0 {
        query = query.Where("age <= ?", maxAge)
    }
    
    // กรณีกําหนด status เป็น array (IN clause)
    // IN clause for multiple statuses
    if statuses, ok := filters["statuses"].([]int); ok && len(statuses) > 0 {
        query = query.Where("status IN ?", statuses)
    }
    
    // กรณี filter วันที่ (ช่วงเวลา)
    // Date range filtering
    if startDate, ok := filters["start_date"].(time.Time); ok && !startDate.IsZero() {
        query = query.Where("created_at >= ?", startDate)
    }
    if endDate, ok := filters["end_date"].(time.Time); ok && !endDate.IsZero() {
        query = query.Where("created_at <= ?", endDate)
    }
    
    // ========== SORTING ==========
    // จัดเรียงตามฟิลด์ที่กำหนด
    // Sort by specified field
    if sortBy, ok := filters["sort_by"].(string); ok && sortBy != "" {
        sortOrder := "ASC" // default
        if order, ok := filters["sort_order"].(string); ok && order == "desc" {
            sortOrder = "DESC"
        }
        query = query.Order(sortBy + " " + sortOrder)
    }
    
    // ========== PAGINATION ==========
    // แบ่งหน้าเพื่อลดภาระ database
    // Pagination to reduce database load
    var limit, offset int
    if l, ok := filters["limit"].(int); ok && l > 0 {
        limit = l
    } else {
        limit = 20 // default limit
    }
    
    if o, ok := filters["offset"].(int); ok && o >= 0 {
        offset = o
    }
    
    // ========== COUNT TOTAL BEFORE PAGINATION ==========
    // นับจํานวนทั้งหมดก่อนแบ่งหน้า (สําหรับ frontend pagination)
    // Count total before pagination (for frontend pagination)
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // ========== EXECUTE QUERY WITH PAGINATION ==========
    // Execute query with pagination applied
    err := query.Limit(limit).Offset(offset).Find(&users).Error
    
    return users, total, err
}
```

## GORM Scopes Pattern (Advanced)

```go
// UserFilterScopes - สร้าง reusable scopes สําหรับ GORM
// UserFilterScopes - Creates reusable scopes for GORM
type UserFilterScopes struct {
    Name   string
    Email  string
    Status int
    MinAge int
    MaxAge int
}

// ToScopes - แปลง filter struct เป็น GORM scopes
// ToScopes - Converts filter struct to GORM scopes
func (f UserFilterScopes) ToScopes() []func(*gorm.DB) *gorm.DB {
    var scopes []func(*gorm.DB) *gorm.DB
    
    if f.Name != "" {
        scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
            return db.Where("name ILIKE ?", "%"+f.Name+"%")
        })
    }
    
    if f.Email != "" {
        scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
            return db.Where("email = ?", f.Email)
        })
    }
    
    if f.Status > 0 {
        scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
            return db.Where("status = ?", f.Status)
        })
    }
    
    if f.MinAge > 0 {
        scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
            return db.Where("age >= ?", f.MinAge)
        })
    }
    
    if f.MaxAge > 0 {
        scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
            return db.Where("age <= ?", f.MaxAge)
        })
    }
    
    return scopes
}

// ตัวอย่างการใช้งาน scopes
// Example usage of scopes
func (r *UserRepository) FindWithScopes(ctx context.Context, filters UserFilterScopes) ([]models.User, error) {
    var users []models.User
    err := r.DB.WithContext(ctx).Scopes(filters.ToScopes()...).Find(&users).Error
    return users, err
}
```

### ใช้อย่างไร / นำไปใช้กรณีไหน
```go
// ตัวอย่าง: API endpoint สําหรับค้นหาผู้ใช้
// Example: API endpoint for user search
func (c *UserController) SearchUsers(ginctx *gin.Context) {
    // สร้าง filters จาก query parameters
    // Build filters from query parameters
    filters := make(map[string]interface{})
    
    if name := ginctx.Query("name"); name != "" {
        filters["name"] = name
    }
    if status := ginctx.Query("status"); status != "" {
        statusInt, _ := strconv.Atoi(status)
        filters["status"] = statusInt
    }
    if limit := ginctx.Query("limit"); limit != "" {
        limitInt, _ := strconv.Atoi(limit)
        filters["limit"] = limitInt
    }
    
    users, total, err := c.userService.FilterUsers(ginctx.Request.Context(), filters)
    // ... handle response
}
```

### ประโยชน์ที่ได้รับ
- ลดจํานวน function ที่ต้องเขียน (function เดียวใช้ได้ทุก filter)
- Frontend ส่ง filter อะไรมาก็ได้ ไม่ต้องแก้ไข backend
- รองรับการเพิ่ม filter ใหม่โดยไม่ต้องเขียน query ใหม่

### ข้อควรระวัง
- ต้อง validate ทุก input ที่มาจาก client
- ระวัง performance เมื่อมี filter เยอะๆ
- map[string]interface{} ทําให้ type safety ลดลง
- SQL Injection risk ถ้าไม่ใช้ parameterized query

### ข้อดี
- ยืดหยุ่นสูง รองรับการเปลี่ยนแปลงได้ดี
- ลด boilerplate code
- รวม logic การ filter ไว้ที่เดียว

### ข้อเสีย
- สูญเสีย type safety (ใช้ interface{})
- debug ยากกว่า static query
- อาจเกิด performance issue ถ้า filter เยอะและไม่ optimize

### ข้อห้าม
- ห้ามรับ filter จาก client แล้วเอาไปใส่ใน WHERE โดยตรง
- ห้ามใช้ reflection เพื่อสร้าง dynamic query โดยไม่จําเป็น
- ห้าม filter บน column ที่ไม่มี index

---

# เอกสารประกอบ: repositories/window_repo.go

## Concept
### คืออะไร?
การใช้ Window Functions ของ PostgreSQL ผ่าน GORM เพื่อคํานวณข้ามแถวข้อมูลโดยไม่ยุบรวมแถว

### มีกี่แบบ?
1. ROW_NUMBER() - เรียงลําดับแถว
2. RANK() / DENSE_RANK() - เรียงลําดับแบบมีอันดับซ้ํา
3. LAG() / LEAD() - เข้าถึงแถวก่อนหน้า/ถัดไป
4. SUM() OVER() - สะสมยอดรวม
5. AVG() OVER() - ค่าเฉลี่ยแบบเลื่อน

## Comment CODE

```go
// GetUsersWithRanking - หาผู้ใช้พร้อมอันดับตามคะแนน
// GetUsersWithRanking - Get users with ranking by score
func (r *UserRepository) GetUsersWithRanking(ctx context.Context) ([]UserWithRank, error) {
    var results []UserWithRank
    
    // raw SQL พร้อม window function
    // Raw SQL with window function
    sql := `
        SELECT 
            id,
            name,
            email,
            score,
            ROW_NUMBER() OVER (ORDER BY score DESC) as row_num,
            RANK() OVER (ORDER BY score DESC) as rank_num,
            DENSE_RANK() OVER (ORDER BY score DESC) as dense_rank_num
        FROM users
        WHERE deleted_at IS NULL
    `
    
    err := r.DB.WithContext(ctx).Raw(sql).Scan(&results).Error
    return results, err
}

// GetTopNPerGroup - หา top N ผู้ใช้ในแต่ละกลุ่ม (department)
// GetTopNPerGroup - Get top N users per group (department)
func (r *UserRepository) GetTopNPerGroup(ctx context.Context, n int) ([]UserWithGroupRank, error) {
    var results []UserWithGroupRank
    
    // PARTITION BY แบ่งกลุ่มตาม department_id
    // PARTITION BY groups by department_id
    sql := `
        SELECT 
            id,
            name,
            department_id,
            score,
            ROW_NUMBER() OVER (
                PARTITION BY department_id 
                ORDER BY score DESC
            ) as rank_in_dept
        FROM users
        WHERE status = 'active'
    `
    
    // subquery เพื่อ filter เฉพาะ top N
    // subquery to filter only top N
    finalSQL := `
        SELECT * FROM (
            ` + sql + `
        ) ranked
        WHERE rank_in_dept <= ?
    `
    
    err := r.DB.WithContext(ctx).Raw(finalSQL, n).Scan(&results).Error
    return results, err
}

// GetUserScoreComparison - เปรียบเทียบคะแนนกับผู้ใช้ก่อนหน้าและถัดไป
// GetUserScoreComparison - Compare score with previous and next users
func (r *UserRepository) GetUserScoreComparison(ctx context.Context) ([]UserComparison, error) {
    var results []UserComparison
    
    sql := `
        SELECT 
            id,
            name,
            score,
            LAG(score, 1, 0) OVER (ORDER BY score) as previous_score,
            LAG(name, 1, '') OVER (ORDER BY score) as previous_name,
            LEAD(score, 1, 0) OVER (ORDER BY score) as next_score,
            LEAD(name, 1, '') OVER (ORDER BY score) as next_name,
            score - LAG(score, 1, 0) OVER (ORDER BY score) as score_diff_from_previous
        FROM users
        WHERE status = 'active'
        ORDER BY score DESC
    `
    
    err := r.DB.WithContext(ctx).Raw(sql).Scan(&results).Error
    return results, err
}

// GetRunningTotal - คํานวณยอดสะสม (running total) ตามลําดับเวลา
// GetRunningTotal - Calculate running total by time order
func (r *UserRepository) GetRunningTotal(ctx context.Context) ([]UserRunningTotal, error) {
    var results []UserRunningTotal
    
    sql := `
        SELECT 
            id,
            name,
            created_at,
            score,
            SUM(score) OVER (
                ORDER BY created_at 
                ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
            ) as running_total,
            AVG(score) OVER (
                ORDER BY created_at 
                ROWS BETWEEN 2 PRECEDING AND CURRENT ROW
            ) as moving_avg_3
        FROM users
        ORDER BY created_at
    `
    
    err := r.DB.WithContext(ctx).Raw(sql).Scan(&results).Error
    return results, err
}
```

### ใช้อย่างไร / นําไปใช้กรณีไหน
```go
// ตัวอย่าง: Leaderboard API
// Example: Leaderboard API
func (s *UserService) GetLeaderboard(ctx context.Context) ([]UserWithRank, error) {
    return s.userRepo.GetUsersWithRanking(ctx)
}

// ตัวอย่าง: หาคนทํายอดขายดีสุดในแต่ละ region
// Example: Find top salesperson per region
func (s *SalesService) GetTopPerRegion(ctx context.Context) ([]SalesPersonRank, error) {
    return s.salesRepo.GetTopNPerGroup(ctx, 3)
}
```

### ประโยชน์ที่ได้รับ
- คํานวณ ranking โดยไม่ต้องใช้ subquery ซับซ้อน
- performance ดีกว่า self-join หลายเท่า
- readable code มากขึ้น

### ข้อควรระวัง
- Window functions อาจใช้ memory สูง
- บาง database (MySQL < 8.0) ไม่ support
- ต้องเข้าใจ ORDER BY ใน window function

### ข้อดี
- ลด complexity ของ query
- performance ดีมากสําหรับการคํานวณข้ามแถว
- support ใน database สมัยใหม่ทุกตัว

### ข้อเสีย
- syntax ซับซ้อนกว่า query ปกติ
- debug ยากกว่า
- ไม่ support ใน database ทุก version

### ข้อห้าม
- ห้ามใช้ window function ใน WHERE clause โดยตรง
- ห้ามใช้กับตารางที่ไม่มี index บน column ที่ใช้ ORDER BY
- ห้ามใช้ ROWS UNBOUNDED PRECEDING กับตารางขนาดใหญ่มาก

---

# เอกสารประกอบ: repositories/recursive_repo.go

## Concept
### คืออะไร?
Recursive Query (WITH RECURSIVE) สําหรับ query ข้อมูลแบบ tree structure (组织结构, comments, categories)

### มีกี่แบบ?
1. Hierarchy traversal (ต้นไม้) - หา descendants/ancestors ทั้งหมด
2. Path enumeration - หาเส้นทางจาก root ถึง node
3. Level calculation - คํานวณระดับความลึก
4. Cycle detection - ตรวจจับวงจรใน tree

## Comment CODE

```go
// EmployeeNode - โครงสร้างข้อมูลพนักงานแบบ tree
// EmployeeNode - Tree structure for employee data
type EmployeeNode struct {
    ID        uint   `json:"id"`
    Name      string `json:"name"`
    ManagerID *uint  `json:"manager_id"`
    Level     int    `json:"level"`
    Path      string `json:"path"`
}

// GetAllSubordinates - หาพนักงานใต้บังคับบัญชาทั้งหมด (递归)
// GetAllSubordinates - Get all subordinates recursively
func (r *EmployeeRepository) GetAllSubordinates(ctx context.Context, managerID uint) ([]EmployeeNode, error) {
    var results []EmployeeNode
    
    sql := `
        WITH RECURSIVE employee_tree AS (
            -- Anchor member: เริ่มจาก manager ที่ต้องการ
            -- Anchor member: Start from the specified manager
            SELECT 
                id, 
                name, 
                manager_id, 
                1 as level,
                name as path
            FROM employees
            WHERE id = ?
            
            UNION ALL
            
            -- Recursive member: หาลูกน้องทั้งหมด
            -- Recursive member: Find all subordinates
            SELECT 
                e.id, 
                e.name, 
                e.manager_id, 
                et.level + 1,
                et.path || ' > ' || e.name
            FROM employees e
            INNER JOIN employee_tree et ON e.manager_id = et.id
        )
        SELECT * FROM employee_tree
        WHERE id != ?  -- ไม่รวมตัว manager เอง
        ORDER BY level, name
    `
    
    err := r.DB.WithContext(ctx).Raw(sql, managerID, managerID).Scan(&results).Error
    return results, err
}

// GetOrganizationTree - สร้าง tree structure ทั้งองค์กร
// GetOrganizationTree - Build entire organization tree
func (r *EmployeeRepository) GetOrganizationTree(ctx context.Context) ([]EmployeeNode, error) {
    var results []EmployeeNode
    
    sql := `
        WITH RECURSIVE org_tree AS (
            -- เริ่มจาก CEO (manager_id IS NULL)
            -- Start from CEO (manager_id IS NULL)
            SELECT 
                id,
                name,
                manager_id,
                1 as level,
                ARRAY[id] as path_array,
                name as path_string
            FROM employees
            WHERE manager_id IS NULL
            
            UNION ALL
            
            SELECT 
                e.id,
                e.name,
                e.manager_id,
                ot.level + 1,
                ot.path_array || e.id,
                ot.path_string || ' > ' || e.name
            FROM employees e
            INNER JOIN org_tree ot ON e.manager_id = ot.id
        )
        SELECT 
            id,
            name,
            manager_id,
            level,
            path_string as path
        FROM org_tree
        ORDER BY path_array
    `
    
    err := r.DB.WithContext(ctx).Raw(sql).Scan(&results).Error
    return results, err
}

// GetManagerHierarchy - หาสายการบังคับบัญชาข้างบน (向上递归)
// GetManagerHierarchy - Get management chain upward
func (r *EmployeeRepository) GetManagerHierarchy(ctx context.Context, employeeID uint) ([]EmployeeNode, error) {
    var results []EmployeeNode
    
    sql := `
        WITH RECURSIVE manager_chain AS (
            -- เริ่มจากพนักงานที่ต้องการ
            -- Start from the specified employee
            SELECT 
                id,
                name,
                manager_id,
                1 as level,
                name as path
            FROM employees
            WHERE id = ?
            
            UNION ALL
            
            -- หัวหน้าของหัวหน้า (向上)
            -- Manager of manager (upward)
            SELECT 
                e.id,
                e.name,
                e.manager_id,
                mc.level + 1,
                e.name || ' > ' || mc.path
            FROM employees e
            INNER JOIN manager_chain mc ON e.id = mc.manager_id
        )
        SELECT * FROM manager_chain
        ORDER BY level
    `
    
    err := r.DB.WithContext(ctx).Raw(sql, employeeID).Scan(&results).Error
    return results, err
}

// GetDepartmentBudget - คํานวณงบประมาณรวมทั้ง department (รวม sub-department)
// GetDepartmentBudget - Calculate total budget including sub-departments
func (r *EmployeeRepository) GetDepartmentBudget(ctx context.Context, deptID uint) (float64, error) {
    var total float64
    
    sql := `
        WITH RECURSIVE dept_tree AS (
            SELECT id, budget, parent_id
            FROM departments
            WHERE id = ?
            
            UNION ALL
            
            SELECT d.id, d.budget, d.parent_id
            FROM departments d
            INNER JOIN dept_tree dt ON d.parent_id = dt.id
        )
        SELECT COALESCE(SUM(budget), 0) as total_budget
        FROM dept_tree
    `
    
    err := r.DB.WithContext(ctx).Raw(sql, deptID).Scan(&total).Error
    return total, err
}
```

### ใช้อย่างไร / นําไปใช้กรณีไหน
```go
// ตัวอย่าง: แสดง组织结构图
// Example: Display organization chart
func (s *OrganizationService) GetTeam(ctx context.Context, managerID uint) ([]EmployeeNode, error) {
    return s.employeeRepo.GetAllSubordinates(ctx, managerID)
}

// ตัวอย่าง: ตรวจสอบเส้นทางการอนุมัติ
// Example: Check approval chain
func (s *ApprovalService) GetApprovalChain(ctx context.Context, requesterID uint) ([]EmployeeNode, error) {
    return s.employeeRepo.GetManagerHierarchy(ctx, requesterID)
}
```

### ประโยชน์ที่ได้รับ
- query ข้อมูล tree structure ด้วย query เดียว
- performance ดีกว่าการ query แบบ recursive ใน application
- รองรับ deep hierarchy

### ข้อควรระวัง
- ต้องมี termination condition (UNION ALL จะจบเมื่อไม่มีแถวเพิ่ม)
- ระวัง infinite loop ถ้า data มี cycle
- Performance อาจลดลงถ้า tree ลึกมากๆ
- PostgreSQL มี recursion depth limit (default 100)

### ข้อดี
- elegant solution สําหรับ hierarchical data
- ได้ข้อมูลทั้ง tree ใน query เดียว
- รองรับ几乎所有 relational database

### ข้อเสีย
- syntax ซับซ้อน เรียนรู้ยาก
- debug ยาก
- performance อาจไม่ดีถ้า tree ใหญ่มาก (>10000 nodes)

### ข้อห้าม
- ห้ามใช้ recursive query ถ้า depth ไม่เกิน 3-4 (ใช้ join แทน)
- ห้ามใช้กับ table ที่มี cycle (ต้อง detect cycle ก่อน)
- ห้าม recursive โดยไม่มี index บน foreign key

---

# เอกสารประกอบ: repositories/json_repo.go

## Concept
### คืออะไร?
การใช้ JSON/JSONB features ของ PostgreSQL ผ่าน GORM สําหรับเก็บและ query ข้อมูล semi-structured

### มีกี่แบบ?
1. JSONB operators (->, ->>, @>, ?)
2. JSONB functions (jsonb_agg, jsonb_build_object)
3. Partial update on JSONB
4. Index on JSONB fields

## Comment CODE

```go
// UserMetadata - metadata รูปแบบ JSON
// UserMetadata - JSON formatted metadata
type UserMetadata struct {
    Preferences map[string]interface{} `json:"preferences"`
    Tags        []string               `json:"tags"`
    Address     Address                `json:"address"`
}

// GetUsersByJSONCondition - ค้นหาผู้ใช้จาก JSON field
// GetUsersByJSONCondition - Find users by JSON field condition
func (r *UserRepository) GetUsersByJSONCondition(ctx context.Context, key, value string) ([]models.User, error) {
    var users []models.User
    
    // ใช้ ->> เพื่อ extract value เป็น text
    // Use ->> to extract value as text
    err := r.DB.WithContext(ctx).
        Where("metadata->>? = ?", key, value).
        Find(&users).Error
    
    return users, err
}

// GetUsersWithTag - ค้นหาผู้ใช้ที่มี tag เฉพาะ (JSON array)
// GetUsersWithTag - Find users with specific tag (JSON array)
func (r *UserRepository) GetUsersWithTag(ctx context.Context, tag string) ([]models.User, error) {
    var users []models.User
    
    // @> operator ใช้เช็ค JSON containment
    // @> operator checks JSON containment
    err := r.DB.WithContext(ctx).
        Where("metadata->'tags' @> ?", fmt.Sprintf(`["%s"]`, tag)).
        Find(&users).Error
    
    return users, err
}

// UpdateJSONField - อัปเดตเฉพาะบาง field ใน JSONB (ไม่ต้อง update ทั้ง record)
// UpdateJSONField - Update specific fields in JSONB (no full record update)
func (r *UserRepository) UpdateJSONField(ctx context.Context, userID uint, path string, value interface{}) error {
    // jsonb_set ใช้อัปเดต nested field
    // jsonb_set updates nested field
    sql := `
        UPDATE users 
        SET metadata = jsonb_set(
            COALESCE(metadata, '{}'::jsonb),
            ?,
            ?,
            true
        )
        WHERE id = ?
    `
    
    // path ต้องเป็น array ของ text เช่น '{preferences,theme}'
    // path must be text array e.g., '{preferences,theme}'
    pathArray := fmt.Sprintf("{%s}", path)
    
    valueJSON, _ := json.Marshal(value)
    
    return r.DB.WithContext(ctx).Exec(sql, pathArray, valueJSON, userID).Error
}

// AggregateJSONData - รวม JSON data จากหลายแถว
// AggregateJSONData - Aggregate JSON data from multiple rows
func (r *UserRepository) AggregateJSONData(ctx context.Context, roleID int) (map[string]interface{}, error) {
    var result struct {
        UsersJSON string `json:"users_json"`
    }
    
    sql := `
        SELECT 
            jsonb_agg(
                jsonb_build_object(
                    'id', id,
                    'name', name,
                    'email', email,
                    'metadata', metadata
                )
            ) as users_json
        FROM users
        WHERE role_id = ?
    `
    
    err := r.DB.WithContext(ctx).Raw(sql, roleID).Scan(&result).Error
    if err != nil {
        return nil, err
    }
    
    var usersData map[string]interface{}
    json.Unmarshal([]byte(result.UsersJSON), &usersData)
    return usersData, nil
}

// SearchInJSON - ค้นหาข้อความใน JSON field ทุก nested level
// SearchInJSON - Search text in JSON field at all nested levels
func (r *UserRepository) SearchInJSON(ctx context.Context, searchText string) ([]models.User, error) {
    var users []models.User
    
    // jsonb_path_exists ใช้ค้นหาแบบ pattern
    // jsonb_path_exists searches with pattern
    sql := `
        SELECT * FROM users
        WHERE jsonb_path_exists(
            metadata,
            '$.* ? (@.type() == "string" && @ like_regex $pattern)',
            jsonb_build_object('pattern', ?)
        )
    `
    
    err := r.DB.WithContext(ctx).Raw(sql, searchText).Scan(&users).Error
    return users, err
}
```

### ใช้อย่างไร / นําไปใช้กรณีไหน
```go
// ตัวอย่าง: เก็บ user preferences
// Example: Store user preferences
func (s *UserService) UpdateTheme(ctx context.Context, userID uint, theme string) error {
    return s.userRepo.UpdateJSONField(ctx, userID, "preferences,theme", theme)
}

// ตัวอย่าง: ค้นหาผู้ใช้ตาม city ใน address JSON
// Example: Find users by city in address JSON
func (s *UserService) FindByCity(ctx context.Context, city string) ([]models.User, error) {
    return s.userRepo.GetUsersByJSONCondition(ctx, "address->>city", city)
}
```

### ประโยชน์ที่ได้รับ
- schema flexibility - เปลี่ยน structure ได้โดยไม่ต้อง migrate
- performance ดี (JSONB มี binary format และ index)
- เหมาะกับข้อมูลที่ไม่รู้ schema ล่วงหน้า

### ข้อควรระวัง
- JSONB มี overhead มากกว่า normal column
- query ซับซ้อนกว่า relational data
- ไม่มี foreign key constraint ใน JSON

### ข้อดี
- เก็บ structured data ใน single column
- query ได้เร็ว (JSONB มี GIN index)
- รองรับ partial update

### ข้อเสีย
- ไม่มี type safety
- debug ยากกว่า relational
- migration ยากเมื่อ structure เปลี่ยน

### ข้อห้าม
- ห้ามเก็บข้อมูลที่มี relation กับตารางอื่นใน JSON
- ห้ามใช้ JSONB แทน normalization
- ห้าม query JSONB field บ่อยๆ โดยไม่มี index

---

# การออกแบบ Workflow และ Dataflow

## Workflow การทํางานของ Dynamic Filter System

```
[Client Request]
      |
      v
[Controller Layer]
รับ query parameters: ?name=john&status=1&limit=10
      |
      v
[Validation Layer]
ตรวจสอบ input: 
- name: string, max 100 chars
- status: int, 0-2
- limit: int, 1-100
      |
      v
[Service Layer]
แปลง parameters เป็น map[string]interface{}
เรียก repository.FilterUsers(filters)
      |
      v
[Repository Layer]
สร้าง GORM query
เพิ่ม WHERE clauses แบบ conditional
เพิ่ม ORDER BY, LIMIT, OFFSET
Execute query
      |
      v
[Database]
PostgreSQL execute query
ใช้ indexes (ถ้ามี) เพื่อ optimize
Return result set
      |
      v
[Response]
JSON response กลับ client
```

## Dataflow สําหรับ Recursive Query

```
[Client] -> [API: GET /api/org/:id/subordinates]
                         |
                         v
              [Controller: GetSubordinates]
                         |
                         v
              [Service: GetTeamTree]
                         |
                         v
              [Repository: GetAllSubordinates]
                         |
                         v
              [PostgreSQL: WITH RECURSIVE]
                         |
    +--------------------+--------------------+
    |                    |                    |
[Anchor]            [Recursive]          [Terminate]
SELECT *           UNION ALL            WHEN no more
WHERE id = X       JOIN tree             rows found
    |                    |
    +--------+-----------+
             |
             v
    [Result Set: id, name, level, path]
             |
             v
    [Build Tree Structure]
             |
             v
    [Return JSON to Client]
```

---

# คู่มือการทดสอบ

## 1. Unit Test สําหรับ Repository

```go
// tests/user_repo_test.go
package tests

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type UserRepoTestSuite struct {
    suite.Suite
    db       *gorm.DB
    repo     *repositories.UserRepository
    testUser *models.User
}

func (s *UserRepoTestSuite) SetupTest() {
    // สร้าง test database
    s.db = setupTestDB()
    s.repo = repositories.NewUserRepository(s.db)
    
    // สร้าง test data
    s.testUser = &models.User{
        Name:  "Test User",
        Email: "test@example.com",
        Age:   25,
    }
    s.db.Create(s.testUser)
}

func (s *UserRepoTestSuite) TearDownTest() {
    // ลบ test data
    s.db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
}

// ทดสอบ dynamic filter
func (s *UserRepoTestSuite) TestFilterUsers() {
    tests := []struct {
        name    string
        filters map[string]interface{}
        want    int
        wantErr bool
    }{
        {
            name: "filter by name partial match",
            filters: map[string]interface{}{
                "name": "Test",
            },
            want: 1,
        },
        {
            name: "filter by age range",
            filters: map[string]interface{}{
                "min_age": 20,
                "max_age": 30,
            },
            want: 1,
        },
        {
            name: "filter with pagination",
            filters: map[string]interface{}{
                "limit":  10,
                "offset": 0,
            },
            want: 1,
        },
        {
            name: "empty filters should return all",
            filters: map[string]interface{}{},
            want:    1,
        },
    }
    
    for _, tt := range tests {
        s.Run(tt.name, func() {
            users, total, err := s.repo.FilterUsers(context.Background(), tt.filters)
            
            if tt.wantErr {
                s.Error(err)
            } else {
                s.NoError(err)
                s.Equal(int64(tt.want), total)
                s.NotNil(users)
            }
        })
    }
}

// ทดสอบ recursive query
func (s *UserRepoTestSuite) TestRecursiveQuery() {
    // สร้าง hierarchy data
    ceo := createEmployee(s.db, "CEO", nil)
    manager := createEmployee(s.db, "Manager", &ceo.ID)
    staff := createEmployee(s.db, "Staff", &manager.ID)
    
    subordinates, err := s.employeeRepo.GetAllSubordinates(context.Background(), ceo.ID)
    
    s.NoError(err)
    s.Equal(2, len(subordinates)) // manager + staff
}

// ทดสอบ window function
func (s *UserRepoTestSuite) TestWindowFunction() {
    // สร้าง users with different scores
    users := []models.User{
        {Name: "User1", Score: 100},
        {Name: "User2", Score: 90},
        {Name: "User3", Score: 80},
    }
    s.db.Create(&users)
    
    rankings, err := s.repo.GetUsersWithRanking(context.Background())
    
    s.NoError(err)
    s.Equal(3, len(rankings))
    s.Equal(1, rankings[0].RankNum) // highest score rank 1
}

// รัน test suite
func TestUserRepoTestSuite(t *testing.T) {
    suite.Run(t, new(UserRepoTestSuite))
}
```

## 2. Integration Test

```go
// tests/integration_test.go
func TestIntegration_FilterAPI(t *testing.T) {
    // สร้าง test server
    router := setupRouter()
    
    tests := []struct {
        name       string
        query      string
        statusCode int
    }{
        {
            name:       "search by name",
            query:      "/users?name=john",
            statusCode: 200,
        },
        {
            name:       "search with pagination",
            query:      "/users?limit=10&page=1",
            statusCode: 200,
        },
        {
            name:       "invalid limit",
            query:      "/users?limit=1000",
            statusCode: 400, // limit too high
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req, _ := http.NewRequest("GET", tt.query, nil)
            resp := performRequest(router, req)
            assert.Equal(t, tt.statusCode, resp.Code)
        })
    }
}

// Benchmark test
func BenchmarkFilterUsers(b *testing.B) {
    db := setupTestDB()
    repo := repositories.NewUserRepository(db)
    
    // seed 10000 users
    seedUsers(db, 10000)
    
    filters := map[string]interface{}{
        "status": 1,
        "min_age": 18,
        "max_age": 30,
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        repo.FilterUsers(context.Background(), filters)
    }
}
```

---

# คู่มือการการใช้งาน

## การติดตั้ง

```bash
# 1. Clone project
git clone https://github.com/yourrepo/advanced-sql-gorm.git
cd advanced-sql-gorm

# 2. Install dependencies
go mod tidy

# 3. Setup PostgreSQL
docker run -d \
  --name postgres-sql \
  -e POSTGRES_USER=myuser \
  -e POSTGRES_PASSWORD=mypassword \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  postgres:15

# 4. Run migrations
go run scripts/migrate.go

# 5. Seed test data
go run scripts/seed_data.go

# 6. Run application
go run main.go
```

## API Endpoints

```bash
# 1. Dynamic Filter Users
GET /api/users?name=john&status=1&min_age=18&limit=20&offset=0

# Response
{
  "data": [...],
  "total": 150,
  "limit": 20,
  "offset": 0
}

# 2. Get User Ranking
GET /api/users/ranking

# Response
[
  {"id":1, "name":"John", "score":100, "rank":1},
  {"id":2, "name":"Jane", "score":95, "rank":2}
]

# 3. Get Organization Tree
GET /api/org/:id/subordinates

# Response
[
  {"id":2, "name":"Manager", "level":2, "path":"CEO > Manager"},
  {"id":3, "name":"Staff", "level":3, "path":"CEO > Manager > Staff"}
]

# 4. Search by JSON Metadata
GET /api/users/metadata?key=preferences.theme&value=dark

# 5. Get Top N per Department
GET /api/departments/top?n=3
```

## ตัวอย่างการใช้งาน Client

```javascript
// Frontend: React example
const searchUsers = async (filters) => {
  const params = new URLSearchParams(filters);
  const response = await fetch(`/api/users?${params}`);
  return response.json();
};

// ใช้งาน
searchUsers({
  name: 'john',
  status: 1,
  limit: 20
}).then(data => console.log(data));
```

```python
# Backend: Python client
import requests

def filter_users(name=None, status=None, min_age=None):
    params = {}
    if name:
        params['name'] = name
    if status:
        params['status'] = status
    if min_age:
        params['min_age'] = min_age
    
    response = requests.get('http://localhost:8080/api/users', params=params)
    return response.json()
```

---

# คู่มือการบำรุงรักษา

## การ Monitor และ Logging

```go
// middleware/logger.go
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        // บันทึก query ก่อน execute
        c.Next()
        
        // บันทึกหลังจาก execute
        duration := time.Since(start)
        
        // log slow queries (>100ms)
        if duration > 100*time.Millisecond {
            log.Printf("[SLOW QUERY] %s %s took %v", 
                c.Request.Method, 
                c.Request.URL.Path,
                duration)
        }
    }
}
```

## การ Optimize Performance

```sql
-- 1. สร้าง indexes ที่จําเป็น
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_age ON users(age);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_status_age ON users(status, age);

-- 2. Partial index for active users only
CREATE INDEX CONCURRENTLY idx_active_users_name ON users(name) 
WHERE status = 1 AND deleted_at IS NULL;

-- 3. GIN index for JSONB
CREATE INDEX idx_users_metadata ON users USING gin(metadata);

-- 4. Index for ILIKE search (pg_trgm)
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_users_name_trgm ON users USING gin(name gin_trgm_ops);
```

## การ Backup และ Recovery

```bash
# Backup database
pg_dump -U myuser -h localhost mydb > backup_$(date +%Y%m%d).sql

# Backup specific tables
pg_dump -U myuser -h localhost -t users -t employees mydb > backup_users.sql

# Restore
psql -U myuser -h localhost mydb < backup_20240101.sql
```

## การ Migration Guide

```go
// migrations/004_add_jsonb_index.sql
-- +migrate Up
CREATE INDEX CONCURRENTLY idx_users_metadata_preferences 
ON users USING gin((metadata->'preferences'));

-- +migrate Down
DROP INDEX CONCURRENTLY idx_users_metadata_preferences;
```

## Troubleshooting Guide

| ปัญหา | สาเหตุ | วิธีแก้ไข |
|-------|--------|----------|
| Query ช้า | Missing index | ใช้ EXPLAIN ANALYZE ตรวจสอบ |
| Recursive query infinite loop | Cycle in data | เพิ่ม cycle detection |
| JSON query ช้า | No GIN index | สร้าง GIN index บน JSONB |
| Memory leak | Too many open connections | ปรับ connection pool |
| Deadlock | Bad transaction order | ใช้ advisory lock |

## Health Check Script

```go
// scripts/health_check.go
func healthCheck(db *gorm.DB) {
    // 1. Check database connection
    sqlDB, _ := db.DB()
    if err := sqlDB.Ping(); err != nil {
        log.Fatal("Database connection failed:", err)
    }
    
    // 2. Check slow queries
    var slowQueries int
    db.Raw(`
        SELECT count(*) FROM pg_stat_statements 
        WHERE mean_time > 100
    `).Scan(&slowQueries)
    
    if slowQueries > 10 {
        log.Warn("High number of slow queries:", slowQueries)
    }
    
    // 3. Check index usage
    var unusedIndexes []string
    db.Raw(`
        SELECT indexname FROM pg_stat_user_indexes 
        WHERE idx_scan = 0
    `).Scan(&unusedIndexes)
    
    if len(unusedIndexes) > 0 {
        log.Info("Unused indexes found:", unusedIndexes)
    }
}
```

---

## สรุป Best Practices

1. **Dynamic Filter**: ใช้ map-based filter กับ conditional WHERE clauses
2. **Window Functions**: ใช้แทน subquery สําหรับ ranking และ running totals
3. **Recursive Query**: ใช้กับ hierarchical data แต่ต้องมี index และ cycle detection
4. **JSONB**: ใช้เมื่อ schema เปลี่ยนแปลงบ่อย แต่ต้องมี GIN index
5. **Index**: สร้าง indexes บน columns ที่ใช้ใน WHERE, ORDER BY, JOIN
6. **Monitoring**: log slow queries และ unused indexes
7. **Testing**: ทดสอบทั้ง unit และ integration รวมถึง benchmark

---

**เอกสารนี้สร้างขึ้นสําหรับเรียนรู้ SQL Advance บน GORM (GoLang)**