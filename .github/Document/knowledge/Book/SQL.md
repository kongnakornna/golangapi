\# สารบัญ



1\. \[หลักการทำงาน SQL แบบ Advance](#หลักการทำงาน-sql-แบบ-advance)

2\. \[ส่วนที่ 1: SQL Advance บน GORM (GoLang)](#ส่วนที่-1-sql-advance-บน-gorm-golang)

3\. \[เอกสารประกอบแต่ละ Component](#เอกสารประกอบแต่ละ-component)



\---



\# หลักการทำงาน SQL แบบ Advance



\## คืออะไร?



SQL ขั้นสูง (Advanced SQL) คือเทคนิคการเขียนคำสั่ง SQL ที่ซับซ้อนกว่า CRUD พื้นฐาน เพื่อจัดการกับ:

\- \*\*Dynamic filtering\*\* - การสร้างเงื่อนไข WHERE แบบพลวัต

\- \*\*Conditional logic\*\* - การใช้ CASE, COALESCE, NULLIF

\- \*\*Recursive queries\*\* - การ query ข้อมูลแบบลำดับชั้น

\- \*\*Window functions\*\* - การคำนวณข้ามแถวข้อมูล

\- \*\*JSON/JSONB operations\*\* - การจัดการข้อมูลแบบ JSON

\- \*\*CTE (Common Table Expressions)\*\* - การสร้าง temporary result sets



\## มีกี่แบบ?



| # | ประเภท | คำอธิบาย | ตัวอย่างการใช้งาน |

|---|--------|----------|------------------|

| 1 | \*\*Conditional Expressions\*\* | CASE, COALESCE, NULLIF | แปลงค่า, จัดการ NULL |

| 2 | \*\*Dynamic WHERE\*\* | Conditional filtering | ระบบค้นหาพร้อม filter หลายตัว |

| 3 | \*\*Window Functions\*\* | ROW\_NUMBER, RANK, LAG, LEAD | Ranking, การเปรียบเทียบแถวก่อนหน้า |

| 4 | \*\*CTE (WITH clause)\*\* | WITH ... AS (...) | ทำให้ query ซับซ้อนอ่านง่าย |

| 5 | \*\*Recursive CTE\*\* | WITH RECURSIVE |  query ข้อมูล tree/hierarchy |

| 6 | \*\*JSON Operations\*\* | ->, ->>, jsonb\_agg | เก็บ/ query ข้อมูล semi-structured |

| 7 | \*\*Aggregate Functions\*\* | STRING\_AGG, ARRAY\_AGG | รวมหลายแถวเป็นค่าเดียว |

| 8 | \*\*Subquery\*\* | EXISTS, IN, ANY, ALL | query ซ้อน query |

| 9 | \*\*Full-Text Search\*\* | tsvector, tsquery | ค้นหาข้อความประสิทธิภาพสูง |

| 10 | \*\*Partitioning\*\* | PARTITION BY | แบ่งตารางใหญ่เป็นส่วนย่อย |



\## ข้อห้ามสำคัญ



```

⚠️ CRITICAL RULES:



1\. ห้ามใช้ loop ใน SQL ถ้าใช้ set-based operation ได้ (performance ตกมาก)

2\. ห้าม SELECT \* ใน production (ใช้เฉพาะ column ที่จำเป็น)

3\. ห้ามทำ dynamic SQL โดยไม่ sanitize input (SQL Injection)

4\. ห้ามใช้ recursive query กับข้อมูลมากๆ (stack overflow)

5\. ห้าม join เกิน 5-7 ตารางใน query เดียว

6\. ห้ามใช้ functions ใน WHERE clause ถ้า column ถูก index

7\. ห้ามทำ aggregation บนข้อมูลที่ไม่ได้ filter ก่อน

```



\---



\# ส่วนที่ 1: SQL Advance บน GORM (GoLang)



\## โครงสร้างโปรเจกต์



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

│   ├── user\_repo.go

│   ├── filter\_repo.go

│   ├── window\_repo.go

│   ├── recursive\_repo.go

│   └── json\_repo.go

├── services/

│   ├── user\_service.go

│   └── filter\_service.go

├── controllers/

│   └── user\_controller.go

├── middleware/

│   └── logger.go

├── utils/

│   └── query\_builder.go

├── tests/

│   ├── user\_repo\_test.go

│   └── integration\_test.go

├── migrations/

│   ├── 001\_create\_users\_table.sql

│   ├── 002\_create\_employee\_tree.sql

│   └── 003\_add\_json\_metadata.sql

└── scripts/

&#x20;   ├── seed\_data.go

&#x20;   └── benchmark.go

```



\## คำอธิบายแต่ละโฟลเดอร์/ไฟล์



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



\---



\# เอกสารประกอบ: repositories/user\_repo.go



\## Concept

\### คืออะไร?

Repository layer ที่รวบรวมการทำงานกับฐานข้อมูลทั้งหมด ใช้ GORM เป็น ORM



\### มีกี่แบบ?

\- Basic CRUD (Create, Read, Update, Delete)

\- Complex query with filters

\- Transaction operations

\- Raw SQL execution



\## Comment CODE (ไทย/อังกฤษ คนละบรรทัด)



```go

// CreateUser - สร้างผู้ใช้ใหม่ในระบบ

// CreateUser - Creates a new user in the system

func (r \*UserRepository) CreateUser(ctx context.Context, user \*models.User) error {

&#x20;   // ใช้ WithContext เพื่อ support context cancellation และ timeout

&#x20;   // Use WithContext to support context cancellation and timeout

&#x20;   return r.DB.WithContext(ctx).Create(user).Error

}



// GetUserByID - ค้นหาผู้ใช้ด้วย ID พร้อม preload ความสัมพันธ์

// GetUserByID - Finds user by ID with preloaded relationships

func (r \*UserRepository) GetUserByID(ctx context.Context, id uint) (\*models.User, error) {

&#x20;   var user models.User

&#x20;   

&#x20;   // Preload ช่วยลด N+1 query problem

&#x20;   // Preload helps reduce N+1 query problem

&#x20;   err := r.DB.WithContext(ctx).

&#x20;       Preload("Roles").           // โหลดข้อมูล roles ที่เกี่ยวข้อง

&#x20;       Preload("Profile").         // โหลดข้อมูล profile ที่เกี่ยวข้อง

&#x20;       Where("id = ? AND deleted\_at IS NULL", id).

&#x20;       First(\&user).Error

&#x20;   

&#x20;   if err != nil {

&#x20;       return nil, err

&#x20;   }

&#x20;   return \&user, nil

}

```



\### ใช้อย่างไร / นำไปใช้กรณีไหน

```go

// ตัวอย่างการเรียกใช้ใน service layer

// Example usage in service layer

func (s \*UserService) GetUser(ctx context.Context, id uint) (\*models.User, error) {

&#x20;   // validate input ก่อนเรียก repository

&#x20;   if id == 0 {

&#x20;       return nil, errors.New("invalid user id")

&#x20;   }

&#x20;   return s.userRepo.GetUserByID(ctx, id)

}

```



\### ประโยชน์ที่ได้รับ

\- แยกการเข้าถึงฐานข้อมูลออกจาก business logic

\- สามารถ mock repository ใน unit test ได้

\- ลด code duplication



\### ข้อควรระวัง

\- ระวังการเรียก Preload มากเกินไป (over-fetching)

\- ต้องจัดการ context timeout ให้เหมาะสม

\- ระวัง N+1 query เมื่อใช้ loop ร่วมกับ database call



\### ข้อดี

\- ควบคุม database operation ได้จากที่เดียว

\- ง่ายต่อการทำ caching

\- สามารถเปลี่ยน ORM ได้โดยไม่กระทบ service layer



\### ข้อเสีย

\- มี boilerplate code เยอะ

\- ต้องเขียน repository ทุก entity

\- อาจซับซ้อนเกินไปสำหรับโปรเจกต์เล็ก



\### ข้อห้าม

\- ห้ามใส่ business logic ใน repository

\- ห้ามเรียก repository ซ้อน repository โดยตรง

\- ห้าม return \*gorm.DB ออกจาก repository



\---



\# เอกสารประกอบ: repositories/filter\_repo.go



\## Concept

\### คืออะไร?

Dynamic filter system ที่สร้าง WHERE clause ตามเงื่อนไขที่ได้รับจาก client



\### มีกี่แบบ?

1\. Map-based filtering

2\. Struct-based filtering with omitempty

3\. Raw SQL with conditional builder

4\. Scopes pattern (GORM specific)



\## Comment CODE



```go

// FilterUsers - กรองผู้ใช้ตามเงื่อนไขที่ส่งมาแบบ dynamic

// FilterUsers - Dynamically filters users based on provided conditions

func (r \*UserRepository) FilterUsers(ctx context.Context, filters map\[string]interface{}) (\[]models.User, int64, error) {

&#x20;   var users \[]models.User

&#x20;   var total int64

&#x20;   

&#x20;   // เริ่มต้น query ด้วย model

&#x20;   // Start query with base model

&#x20;   query := r.DB.WithContext(ctx).Model(\&models.User{})

&#x20;   

&#x20;   // ========== DYNAMIC WHERE CLAUSE ==========

&#x20;   // ตรวจสอบและเพิ่มเงื่อนไขทีละฟิลด์

&#x20;   // Check and add conditions field by field

&#x20;   

&#x20;   // กรณีค้นหาด้วยชื่อแบบ partial match (case-insensitive)

&#x20;   // Case-insensitive partial match for name search

&#x20;   if name, ok := filters\["name"].(string); ok \&\& name != "" {

&#x20;       // ILIKE ใน PostgreSQL, GORM จะแปลงเป็น LIKE สำหรับ DB อื่น

&#x20;       // ILIKE in PostgreSQL, GORM converts for other databases

&#x20;       query = query.Where("name ILIKE ?", "%"+name+"%")

&#x20;   }

&#x20;   

&#x20;   // กรณีกําหนด email แบบ exact match

&#x20;   // Exact match for email

&#x20;   if email, ok := filters\["email"].(string); ok \&\& email != "" {

&#x20;       query = query.Where("email = ?", email)

&#x20;   }

&#x20;   

&#x20;   // กรณีกําหนด age range (min-max)

&#x20;   // Age range filtering

&#x20;   if minAge, ok := filters\["min\_age"].(int); ok \&\& minAge > 0 {

&#x20;       query = query.Where("age >= ?", minAge)

&#x20;   }

&#x20;   if maxAge, ok := filters\["max\_age"].(int); ok \&\& maxAge > 0 {

&#x20;       query = query.Where("age <= ?", maxAge)

&#x20;   }

&#x20;   

&#x20;   // กรณีกําหนด status เป็น array (IN clause)

&#x20;   // IN clause for multiple statuses

&#x20;   if statuses, ok := filters\["statuses"].(\[]int); ok \&\& len(statuses) > 0 {

&#x20;       query = query.Where("status IN ?", statuses)

&#x20;   }

&#x20;   

&#x20;   // กรณี filter วันที่ (ช่วงเวลา)

&#x20;   // Date range filtering

&#x20;   if startDate, ok := filters\["start\_date"].(time.Time); ok \&\& !startDate.IsZero() {

&#x20;       query = query.Where("created\_at >= ?", startDate)

&#x20;   }

&#x20;   if endDate, ok := filters\["end\_date"].(time.Time); ok \&\& !endDate.IsZero() {

&#x20;       query = query.Where("created\_at <= ?", endDate)

&#x20;   }

&#x20;   

&#x20;   // ========== SORTING ==========

&#x20;   // จัดเรียงตามฟิลด์ที่กำหนด

&#x20;   // Sort by specified field

&#x20;   if sortBy, ok := filters\["sort\_by"].(string); ok \&\& sortBy != "" {

&#x20;       sortOrder := "ASC" // default

&#x20;       if order, ok := filters\["sort\_order"].(string); ok \&\& order == "desc" {

&#x20;           sortOrder = "DESC"

&#x20;       }

&#x20;       query = query.Order(sortBy + " " + sortOrder)

&#x20;   }

&#x20;   

&#x20;   // ========== PAGINATION ==========

&#x20;   // แบ่งหน้าเพื่อลดภาระ database

&#x20;   // Pagination to reduce database load

&#x20;   var limit, offset int

&#x20;   if l, ok := filters\["limit"].(int); ok \&\& l > 0 {

&#x20;       limit = l

&#x20;   } else {

&#x20;       limit = 20 // default limit

&#x20;   }

&#x20;   

&#x20;   if o, ok := filters\["offset"].(int); ok \&\& o >= 0 {

&#x20;       offset = o

&#x20;   }

&#x20;   

&#x20;   // ========== COUNT TOTAL BEFORE PAGINATION ==========

&#x20;   // นับจํานวนทั้งหมดก่อนแบ่งหน้า (สําหรับ frontend pagination)

&#x20;   // Count total before pagination (for frontend pagination)

&#x20;   if err := query.Count(\&total).Error; err != nil {

&#x20;       return nil, 0, err

&#x20;   }

&#x20;   

&#x20;   // ========== EXECUTE QUERY WITH PAGINATION ==========

&#x20;   // Execute query with pagination applied

&#x20;   err := query.Limit(limit).Offset(offset).Find(\&users).Error

&#x20;   

&#x20;   return users, total, err

}

```



\## GORM Scopes Pattern (Advanced)



```go

// UserFilterScopes - สร้าง reusable scopes สําหรับ GORM

// UserFilterScopes - Creates reusable scopes for GORM

type UserFilterScopes struct {

&#x20;   Name   string

&#x20;   Email  string

&#x20;   Status int

&#x20;   MinAge int

&#x20;   MaxAge int

}



// ToScopes - แปลง filter struct เป็น GORM scopes

// ToScopes - Converts filter struct to GORM scopes

func (f UserFilterScopes) ToScopes() \[]func(\*gorm.DB) \*gorm.DB {

&#x20;   var scopes \[]func(\*gorm.DB) \*gorm.DB

&#x20;   

&#x20;   if f.Name != "" {

&#x20;       scopes = append(scopes, func(db \*gorm.DB) \*gorm.DB {

&#x20;           return db.Where("name ILIKE ?", "%"+f.Name+"%")

&#x20;       })

&#x20;   }

&#x20;   

&#x20;   if f.Email != "" {

&#x20;       scopes = append(scopes, func(db \*gorm.DB) \*gorm.DB {

&#x20;           return db.Where("email = ?", f.Email)

&#x20;       })

&#x20;   }

&#x20;   

&#x20;   if f.Status > 0 {

&#x20;       scopes = append(scopes, func(db \*gorm.DB) \*gorm.DB {

&#x20;           return db.Where("status = ?", f.Status)

&#x20;       })

&#x20;   }

&#x20;   

&#x20;   if f.MinAge > 0 {

&#x20;       scopes = append(scopes, func(db \*gorm.DB) \*gorm.DB {

&#x20;           return db.Where("age >= ?", f.MinAge)

&#x20;       })

&#x20;   }

&#x20;   

&#x20;   if f.MaxAge > 0 {

&#x20;       scopes = append(scopes, func(db \*gorm.DB) \*gorm.DB {

&#x20;           return db.Where("age <= ?", f.MaxAge)

&#x20;       })

&#x20;   }

&#x20;   

&#x20;   return scopes

}



// ตัวอย่างการใช้งาน scopes

// Example usage of scopes

func (r \*UserRepository) FindWithScopes(ctx context.Context, filters UserFilterScopes) (\[]models.User, error) {

&#x20;   var users \[]models.User

&#x20;   err := r.DB.WithContext(ctx).Scopes(filters.ToScopes()...).Find(\&users).Error

&#x20;   return users, err

}

```



\### ใช้อย่างไร / นำไปใช้กรณีไหน

```go

// ตัวอย่าง: API endpoint สําหรับค้นหาผู้ใช้

// Example: API endpoint for user search

func (c \*UserController) SearchUsers(ginctx \*gin.Context) {

&#x20;   // สร้าง filters จาก query parameters

&#x20;   // Build filters from query parameters

&#x20;   filters := make(map\[string]interface{})

&#x20;   

&#x20;   if name := ginctx.Query("name"); name != "" {

&#x20;       filters\["name"] = name

&#x20;   }

&#x20;   if status := ginctx.Query("status"); status != "" {

&#x20;       statusInt, \_ := strconv.Atoi(status)

&#x20;       filters\["status"] = statusInt

&#x20;   }

&#x20;   if limit := ginctx.Query("limit"); limit != "" {

&#x20;       limitInt, \_ := strconv.Atoi(limit)

&#x20;       filters\["limit"] = limitInt

&#x20;   }

&#x20;   

&#x20;   users, total, err := c.userService.FilterUsers(ginctx.Request.Context(), filters)

&#x20;   // ... handle response

}

```



\### ประโยชน์ที่ได้รับ

\- ลดจํานวน function ที่ต้องเขียน (function เดียวใช้ได้ทุก filter)

\- Frontend ส่ง filter อะไรมาก็ได้ ไม่ต้องแก้ไข backend

\- รองรับการเพิ่ม filter ใหม่โดยไม่ต้องเขียน query ใหม่



\### ข้อควรระวัง

\- ต้อง validate ทุก input ที่มาจาก client

\- ระวัง performance เมื่อมี filter เยอะๆ

\- map\[string]interface{} ทําให้ type safety ลดลง

\- SQL Injection risk ถ้าไม่ใช้ parameterized query



\### ข้อดี

\- ยืดหยุ่นสูง รองรับการเปลี่ยนแปลงได้ดี

\- ลด boilerplate code

\- รวม logic การ filter ไว้ที่เดียว



\### ข้อเสีย

\- สูญเสีย type safety (ใช้ interface{})

\- debug ยากกว่า static query

\- อาจเกิด performance issue ถ้า filter เยอะและไม่ optimize



\### ข้อห้าม

\- ห้ามรับ filter จาก client แล้วเอาไปใส่ใน WHERE โดยตรง

\- ห้ามใช้ reflection เพื่อสร้าง dynamic query โดยไม่จําเป็น

\- ห้าม filter บน column ที่ไม่มี index



\---



\# เอกสารประกอบ: repositories/window\_repo.go



\## Concept

\### คืออะไร?

การใช้ Window Functions ของ PostgreSQL ผ่าน GORM เพื่อคํานวณข้ามแถวข้อมูลโดยไม่ยุบรวมแถว



\### มีกี่แบบ?

1\. ROW\_NUMBER() - เรียงลําดับแถว

2\. RANK() / DENSE\_RANK() - เรียงลําดับแบบมีอันดับซ้ํา

3\. LAG() / LEAD() - เข้าถึงแถวก่อนหน้า/ถัดไป

4\. SUM() OVER() - สะสมยอดรวม

5\. AVG() OVER() - ค่าเฉลี่ยแบบเลื่อน



\## Comment CODE



```go

// GetUsersWithRanking - หาผู้ใช้พร้อมอันดับตามคะแนน

// GetUsersWithRanking - Get users with ranking by score

func (r \*UserRepository) GetUsersWithRanking(ctx context.Context) (\[]UserWithRank, error) {

&#x20;   var results \[]UserWithRank

&#x20;   

&#x20;   // raw SQL พร้อม window function

&#x20;   // Raw SQL with window function

&#x20;   sql := `

&#x20;       SELECT 

&#x20;           id,

&#x20;           name,

&#x20;           email,

&#x20;           score,

&#x20;           ROW\_NUMBER() OVER (ORDER BY score DESC) as row\_num,

&#x20;           RANK() OVER (ORDER BY score DESC) as rank\_num,

&#x20;           DENSE\_RANK() OVER (ORDER BY score DESC) as dense\_rank\_num

&#x20;       FROM users

&#x20;       WHERE deleted\_at IS NULL

&#x20;   `

&#x20;   

&#x20;   err := r.DB.WithContext(ctx).Raw(sql).Scan(\&results).Error

&#x20;   return results, err

}



// GetTopNPerGroup - หา top N ผู้ใช้ในแต่ละกลุ่ม (department)

// GetTopNPerGroup - Get top N users per group (department)

func (r \*UserRepository) GetTopNPerGroup(ctx context.Context, n int) (\[]UserWithGroupRank, error) {

&#x20;   var results \[]UserWithGroupRank

&#x20;   

&#x20;   // PARTITION BY แบ่งกลุ่มตาม department\_id

&#x20;   // PARTITION BY groups by department\_id

&#x20;   sql := `

&#x20;       SELECT 

&#x20;           id,

&#x20;           name,

&#x20;           department\_id,

&#x20;           score,

&#x20;           ROW\_NUMBER() OVER (

&#x20;               PARTITION BY department\_id 

&#x20;               ORDER BY score DESC

&#x20;           ) as rank\_in\_dept

&#x20;       FROM users

&#x20;       WHERE status = 'active'

&#x20;   `

&#x20;   

&#x20;   // subquery เพื่อ filter เฉพาะ top N

&#x20;   // subquery to filter only top N

&#x20;   finalSQL := `

&#x20;       SELECT \* FROM (

&#x20;           ` + sql + `

&#x20;       ) ranked

&#x20;       WHERE rank\_in\_dept <= ?

&#x20;   `

&#x20;   

&#x20;   err := r.DB.WithContext(ctx).Raw(finalSQL, n).Scan(\&results).Error

&#x20;   return results, err

}



// GetUserScoreComparison - เปรียบเทียบคะแนนกับผู้ใช้ก่อนหน้าและถัดไป

// GetUserScoreComparison - Compare score with previous and next users

func (r \*UserRepository) GetUserScoreComparison(ctx context.Context) (\[]UserComparison, error) {

&#x20;   var results \[]UserComparison

&#x20;   

&#x20;   sql := `

&#x20;       SELECT 

&#x20;           id,

&#x20;           name,

&#x20;           score,

&#x20;           LAG(score, 1, 0) OVER (ORDER BY score) as previous\_score,

&#x20;           LAG(name, 1, '') OVER (ORDER BY score) as previous\_name,

&#x20;           LEAD(score, 1, 0) OVER (ORDER BY score) as next\_score,

&#x20;           LEAD(name, 1, '') OVER (ORDER BY score) as next\_name,

&#x20;           score - LAG(score, 1, 0) OVER (ORDER BY score) as score\_diff\_from\_previous

&#x20;       FROM users

&#x20;       WHERE status = 'active'

&#x20;       ORDER BY score DESC

&#x20;   `

&#x20;   

&#x20;   err := r.DB.WithContext(ctx).Raw(sql).Scan(\&results).Error

&#x20;   return results, err

}



// GetRunningTotal - คํานวณยอดสะสม (running total) ตามลําดับเวลา

// GetRunningTotal - Calculate running total by time order

func (r \*UserRepository) GetRunningTotal(ctx context.Context) (\[]UserRunningTotal, error) {

&#x20;   var results \[]UserRunningTotal

&#x20;   

&#x20;   sql := `

&#x20;       SELECT 

&#x20;           id,

&#x20;           name,

&#x20;           created\_at,

&#x20;           score,

&#x20;           SUM(score) OVER (

&#x20;               ORDER BY created\_at 

&#x20;               ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW

&#x20;           ) as running\_total,

&#x20;           AVG(score) OVER (

&#x20;               ORDER BY created\_at 

&#x20;               ROWS BETWEEN 2 PRECEDING AND CURRENT ROW

&#x20;           ) as moving\_avg\_3

&#x20;       FROM users

&#x20;       ORDER BY created\_at

&#x20;   `

&#x20;   

&#x20;   err := r.DB.WithContext(ctx).Raw(sql).Scan(\&results).Error

&#x20;   return results, err

}

```



\### ใช้อย่างไร / นําไปใช้กรณีไหน

```go

// ตัวอย่าง: Leaderboard API

// Example: Leaderboard API

func (s \*UserService) GetLeaderboard(ctx context.Context) (\[]UserWithRank, error) {

&#x20;   return s.userRepo.GetUsersWithRanking(ctx)

}



// ตัวอย่าง: หาคนทํายอดขายดีสุดในแต่ละ region

// Example: Find top salesperson per region

func (s \*SalesService) GetTopPerRegion(ctx context.Context) (\[]SalesPersonRank, error) {

&#x20;   return s.salesRepo.GetTopNPerGroup(ctx, 3)

}

```



\### ประโยชน์ที่ได้รับ

\- คํานวณ ranking โดยไม่ต้องใช้ subquery ซับซ้อน

\- performance ดีกว่า self-join หลายเท่า

\- readable code มากขึ้น



\### ข้อควรระวัง

\- Window functions อาจใช้ memory สูง

\- บาง database (MySQL < 8.0) ไม่ support

\- ต้องเข้าใจ ORDER BY ใน window function



\### ข้อดี

\- ลด complexity ของ query

\- performance ดีมากสําหรับการคํานวณข้ามแถว

\- support ใน database สมัยใหม่ทุกตัว



\### ข้อเสีย

\- syntax ซับซ้อนกว่า query ปกติ

\- debug ยากกว่า

\- ไม่ support ใน database ทุก version



\### ข้อห้าม

\- ห้ามใช้ window function ใน WHERE clause โดยตรง

\- ห้ามใช้กับตารางที่ไม่มี index บน column ที่ใช้ ORDER BY

\- ห้ามใช้ ROWS UNBOUNDED PRECEDING กับตารางขนาดใหญ่มาก



\---



\# เอกสารประกอบ: repositories/recursive\_repo.go



\## Concept

\### คืออะไร?

Recursive Query (WITH RECURSIVE) สําหรับ query ข้อมูลแบบ tree structure (组织结构, comments, categories)



\### มีกี่แบบ?

1\. Hierarchy traversal (ต้นไม้) - หา descendants/ancestors ทั้งหมด

2\. Path enumeration - หาเส้นทางจาก root ถึง node

3\. Level calculation - คํานวณระดับความลึก

4\. Cycle detection - ตรวจจับวงจรใน tree



\## Comment CODE



```go

// EmployeeNode - โครงสร้างข้อมูลพนักงานแบบ tree

// EmployeeNode - Tree structure for employee data

type EmployeeNode struct {

&#x20;   ID        uint   `json:"id"`

&#x20;   Name      string `json:"name"`

&#x20;   ManagerID \*uint  `json:"manager\_id"`

&#x20;   Level     int    `json:"level"`

&#x20;   Path      string `json:"path"`

}



// GetAllSubordinates - หาพนักงานใต้บังคับบัญชาทั้งหมด (递归)

// GetAllSubordinates - Get all subordinates recursively

func (r \*EmployeeRepository) GetAllSubordinates(ctx context.Context, managerID uint) (\[]EmployeeNode, error) {

&#x20;   var results \[]EmployeeNode

&#x20;   

&#x20;   sql := `

&#x20;       WITH RECURSIVE employee\_tree AS (

&#x20;           -- Anchor member: เริ่มจาก manager ที่ต้องการ

&#x20;           -- Anchor member: Start from the specified manager

&#x20;           SELECT 

&#x20;               id, 

&#x20;               name, 

&#x20;               manager\_id, 

&#x20;               1 as level,

&#x20;               name as path

&#x20;           FROM employees

&#x20;           WHERE id = ?

&#x20;           

&#x20;           UNION ALL

&#x20;           

&#x20;           -- Recursive member: หาลูกน้องทั้งหมด

&#x20;           -- Recursive member: Find all subordinates

&#x20;           SELECT 

&#x20;               e.id, 

&#x20;               e.name, 

&#x20;               e.manager\_id, 

&#x20;               et.level + 1,

&#x20;               et.path || ' > ' || e.name

&#x20;           FROM employees e

&#x20;           INNER JOIN employee\_tree et ON e.manager\_id = et.id

&#x20;       )

&#x20;       SELECT \* FROM employee\_tree

&#x20;       WHERE id != ?  -- ไม่รวมตัว manager เอง

&#x20;       ORDER BY level, name

&#x20;   `

&#x20;   

&#x20;   err := r.DB.WithContext(ctx).Raw(sql, managerID, managerID).Scan(\&results).Error

&#x20;   return results, err

}



// GetOrganizationTree - สร้าง tree structure ทั้งองค์กร

// GetOrganizationTree - Build entire organization tree

func (r \*EmployeeRepository) GetOrganizationTree(ctx context.Context) (\[]EmployeeNode, error) {

&#x20;   var results \[]EmployeeNode

&#x20;   

&#x20;   sql := `

&#x20;       WITH RECURSIVE org\_tree AS (

&#x20;           -- เริ่มจาก CEO (manager\_id IS NULL)

&#x20;           -- Start from CEO (manager\_id IS NULL)

&#x20;           SELECT 

&#x20;               id,

&#x20;               name,

&#x20;               manager\_id,

&#x20;               1 as level,

&#x20;               ARRAY\[id] as path\_array,

&#x20;               name as path\_string

&#x20;           FROM employees

&#x20;           WHERE manager\_id IS NULL

&#x20;           

&#x20;           UNION ALL

&#x20;           

&#x20;           SELECT 

&#x20;               e.id,

&#x20;               e.name,

&#x20;               e.manager\_id,

&#x20;               ot.level + 1,

&#x20;               ot.path\_array || e.id,

&#x20;               ot.path\_string || ' > ' || e.name

&#x20;           FROM employees e

&#x20;           INNER JOIN org\_tree ot ON e.manager\_id = ot.id

&#x20;       )

&#x20;       SELECT 

&#x20;           id,

&#x20;           name,

&#x20;           manager\_id,

&#x20;           level,

&#x20;           path\_string as path

&#x20;       FROM org\_tree

&#x20;       ORDER BY path\_array

&#x20;   `

&#x20;   

&#x20;   err := r.DB.WithContext(ctx).Raw(sql).Scan(\&results).Error

&#x20;   return results, err

}



// GetManagerHierarchy - หาสายการบังคับบัญชาข้างบน (向上递归)

// GetManagerHierarchy - Get management chain upward

func (r \*EmployeeRepository) GetManagerHierarchy(ctx context.Context, employeeID uint) (\[]EmployeeNode, error) {

&#x20;   var results \[]EmployeeNode

&#x20;   

&#x20;   sql := `

&#x20;       WITH RECURSIVE manager\_chain AS (

&#x20;           -- เริ่มจากพนักงานที่ต้องการ

&#x20;           -- Start from the specified employee

&#x20;           SELECT 

&#x20;               id,

&#x20;               name,

&#x20;               manager\_id,

&#x20;               1 as level,

&#x20;               name as path

&#x20;           FROM employees

&#x20;           WHERE id = ?

&#x20;           

&#x20;           UNION ALL

&#x20;           

&#x20;           -- หัวหน้าของหัวหน้า (向上)

&#x20;           -- Manager of manager (upward)

&#x20;           SELECT 

&#x20;               e.id,

&#x20;               e.name,

&#x20;               e.manager\_id,

&#x20;               mc.level + 1,

&#x20;               e.name || ' > ' || mc.path

&#x20;           FROM employees e

&#x20;           INNER JOIN manager\_chain mc ON e.id = mc.manager\_id

&#x20;       )

&#x20;       SELECT \* FROM manager\_chain

&#x20;       ORDER BY level

&#x20;   `

&#x20;   

&#x20;   err := r.DB.WithContext(ctx).Raw(sql, employeeID).Scan(\&results).Error

&#x20;   return results, err

}



// GetDepartmentBudget - คํานวณงบประมาณรวมทั้ง department (รวม sub-department)

// GetDepartmentBudget - Calculate total budget including sub-departments

func (r \*EmployeeRepository) GetDepartmentBudget(ctx context.Context, deptID uint) (float64, error) {

&#x20;   var total float64

&#x20;   

&#x20;   sql := `

&#x20;       WITH RECURSIVE dept\_tree AS (

&#x20;           SELECT id, budget, parent\_id

&#x20;           FROM departments

&#x20;           WHERE id = ?

&#x20;           

&#x20;           UNION ALL

&#x20;           

&#x20;           SELECT d.id, d.budget, d.parent\_id

&#x20;           FROM departments d

&#x20;           INNER JOIN dept\_tree dt ON d.parent\_id = dt.id

&#x20;       )

&#x20;       SELECT COALESCE(SUM(budget), 0) as total\_budget

&#x20;       FROM dept\_tree

&#x20;   `

&#x20;   

&#x20;   err := r.DB.WithContext(ctx).Raw(sql, deptID).Scan(\&total).Error

&#x20;   return total, err

}

```



\### ใช้อย่างไร / นําไปใช้กรณีไหน

```go

// ตัวอย่าง: แสดง组织结构图

// Example: Display organization chart

func (s \*OrganizationService) GetTeam(ctx context.Context, managerID uint) (\[]EmployeeNode, error) {

&#x20;   return s.employeeRepo.GetAllSubordinates(ctx, managerID)

}



// ตัวอย่าง: ตรวจสอบเส้นทางการอนุมัติ

// Example: Check approval chain

func (s \*ApprovalService) GetApprovalChain(ctx context.Context, requesterID uint) (\[]EmployeeNode, error) {

&#x20;   return s.employeeRepo.GetManagerHierarchy(ctx, requesterID)

}

```



\### ประโยชน์ที่ได้รับ

\- query ข้อมูล tree structure ด้วย query เดียว

\- performance ดีกว่าการ query แบบ recursive ใน application

\- รองรับ deep hierarchy



\### ข้อควรระวัง

\- ต้องมี termination condition (UNION ALL จะจบเมื่อไม่มีแถวเพิ่ม)

\- ระวัง infinite loop ถ้า data มี cycle

\- Performance อาจลดลงถ้า tree ลึกมากๆ

\- PostgreSQL มี recursion depth limit (default 100)



\### ข้อดี

\- elegant solution สําหรับ hierarchical data

\- ได้ข้อมูลทั้ง tree ใน query เดียว

\- รองรับ几乎所有 relational database



\### ข้อเสีย

\- syntax ซับซ้อน เรียนรู้ยาก

\- debug ยาก

\- performance อาจไม่ดีถ้า tree ใหญ่มาก (>10000 nodes)



\### ข้อห้าม

\- ห้ามใช้ recursive query ถ้า depth ไม่เกิน 3-4 (ใช้ join แทน)

\- ห้ามใช้กับ table ที่มี cycle (ต้อง detect cycle ก่อน)

\- ห้าม recursive โดยไม่มี index บน foreign key



\---



\# เอกสารประกอบ: repositories/json\_repo.go



\## Concept

\### คืออะไร?

การใช้ JSON/JSONB features ของ PostgreSQL ผ่าน GORM สําหรับเก็บและ query ข้อมูล semi-structured



\### มีกี่แบบ?

1\. JSONB operators (->, ->>, @>, ?)

2\. JSONB functions (jsonb\_agg, jsonb\_build\_object)

3\. Partial update on JSONB

4\. Index on JSONB fields



\## Comment CODE



```go

// UserMetadata - metadata รูปแบบ JSON

// UserMetadata - JSON formatted metadata

type UserMetadata struct {

&#x20;   Preferences map\[string]interface{} `json:"preferences"`

&#x20;   Tags        \[]string               `json:"tags"`

&#x20;   Address     Address                `json:"address"`

}



// GetUsersByJSONCondition - ค้นหาผู้ใช้จาก JSON field

// GetUsersByJSONCondition - Find users by JSON field condition

func (r \*UserRepository) GetUsersByJSONCondition(ctx context.Context, key, value string) (\[]models.User, error) {

&#x20;   var users \[]models.User

&#x20;   

&#x20;   // ใช้ ->> เพื่อ extract value เป็น text

&#x20;   // Use ->> to extract value as text

&#x20;   err := r.DB.WithContext(ctx).

&#x20;       Where("metadata->>? = ?", key, value).

&#x20;       Find(\&users).Error

&#x20;   

&#x20;   return users, err

}



// GetUsersWithTag - ค้นหาผู้ใช้ที่มี tag เฉพาะ (JSON array)

// GetUsersWithTag - Find users with specific tag (JSON array)

func (r \*UserRepository) GetUsersWithTag(ctx context.Context, tag string) (\[]models.User, error) {

&#x20;   var users \[]models.User

&#x20;   

&#x20;   // @> operator ใช้เช็ค JSON containment

&#x20;   // @> operator checks JSON containment

&#x20;   err := r.DB.WithContext(ctx).

&#x20;       Where("metadata->'tags' @> ?", fmt.Sprintf(`\["%s"]`, tag)).

&#x20;       Find(\&users).Error

&#x20;   

&#x20;   return users, err

}



// UpdateJSONField - อัปเดตเฉพาะบาง field ใน JSONB (ไม่ต้อง update ทั้ง record)

// UpdateJSONField - Update specific fields in JSONB (no full record update)

func (r \*UserRepository) UpdateJSONField(ctx context.Context, userID uint, path string, value interface{}) error {

&#x20;   // jsonb\_set ใช้อัปเดต nested field

&#x20;   // jsonb\_set updates nested field

&#x20;   sql := `

&#x20;       UPDATE users 

&#x20;       SET metadata = jsonb\_set(

&#x20;           COALESCE(metadata, '{}'::jsonb),

&#x20;           ?,

&#x20;           ?,

&#x20;           true

&#x20;       )

&#x20;       WHERE id = ?

&#x20;   `

&#x20;   

&#x20;   // path ต้องเป็น array ของ text เช่น '{preferences,theme}'

&#x20;   // path must be text array e.g., '{preferences,theme}'

&#x20;   pathArray := fmt.Sprintf("{%s}", path)

&#x20;   

&#x20;   valueJSON, \_ := json.Marshal(value)

&#x20;   

&#x20;   return r.DB.WithContext(ctx).Exec(sql, pathArray, valueJSON, userID).Error

}



// AggregateJSONData - รวม JSON data จากหลายแถว

// AggregateJSONData - Aggregate JSON data from multiple rows

func (r \*UserRepository) AggregateJSONData(ctx context.Context, roleID int) (map\[string]interface{}, error) {

&#x20;   var result struct {

&#x20;       UsersJSON string `json:"users\_json"`

&#x20;   }

&#x20;   

&#x20;   sql := `

&#x20;       SELECT 

&#x20;           jsonb\_agg(

&#x20;               jsonb\_build\_object(

&#x20;                   'id', id,

&#x20;                   'name', name,

&#x20;                   'email', email,

&#x20;                   'metadata', metadata

&#x20;               )

&#x20;           ) as users\_json

&#x20;       FROM users

&#x20;       WHERE role\_id = ?

&#x20;   `

&#x20;   

&#x20;   err := r.DB.WithContext(ctx).Raw(sql, roleID).Scan(\&result).Error

&#x20;   if err != nil {

&#x20;       return nil, err

&#x20;   }

&#x20;   

&#x20;   var usersData map\[string]interface{}

&#x20;   json.Unmarshal(\[]byte(result.UsersJSON), \&usersData)

&#x20;   return usersData, nil

}



// SearchInJSON - ค้นหาข้อความใน JSON field ทุก nested level

// SearchInJSON - Search text in JSON field at all nested levels

func (r \*UserRepository) SearchInJSON(ctx context.Context, searchText string) (\[]models.User, error) {

&#x20;   var users \[]models.User

&#x20;   

&#x20;   // jsonb\_path\_exists ใช้ค้นหาแบบ pattern

&#x20;   // jsonb\_path\_exists searches with pattern

&#x20;   sql := `

&#x20;       SELECT \* FROM users

&#x20;       WHERE jsonb\_path\_exists(

&#x20;           metadata,

&#x20;           '$.\* ? (@.type() == "string" \&\& @ like\_regex $pattern)',

&#x20;           jsonb\_build\_object('pattern', ?)

&#x20;       )

&#x20;   `

&#x20;   

&#x20;   err := r.DB.WithContext(ctx).Raw(sql, searchText).Scan(\&users).Error

&#x20;   return users, err

}

```



\### ใช้อย่างไร / นําไปใช้กรณีไหน

```go

// ตัวอย่าง: เก็บ user preferences

// Example: Store user preferences

func (s \*UserService) UpdateTheme(ctx context.Context, userID uint, theme string) error {

&#x20;   return s.userRepo.UpdateJSONField(ctx, userID, "preferences,theme", theme)

}



// ตัวอย่าง: ค้นหาผู้ใช้ตาม city ใน address JSON

// Example: Find users by city in address JSON

func (s \*UserService) FindByCity(ctx context.Context, city string) (\[]models.User, error) {

&#x20;   return s.userRepo.GetUsersByJSONCondition(ctx, "address->>city", city)

}

```



\### ประโยชน์ที่ได้รับ

\- schema flexibility - เปลี่ยน structure ได้โดยไม่ต้อง migrate

\- performance ดี (JSONB มี binary format และ index)

\- เหมาะกับข้อมูลที่ไม่รู้ schema ล่วงหน้า



\### ข้อควรระวัง

\- JSONB มี overhead มากกว่า normal column

\- query ซับซ้อนกว่า relational data

\- ไม่มี foreign key constraint ใน JSON



\### ข้อดี

\- เก็บ structured data ใน single column

\- query ได้เร็ว (JSONB มี GIN index)

\- รองรับ partial update



\### ข้อเสีย

\- ไม่มี type safety

\- debug ยากกว่า relational

\- migration ยากเมื่อ structure เปลี่ยน



\### ข้อห้าม

\- ห้ามเก็บข้อมูลที่มี relation กับตารางอื่นใน JSON

\- ห้ามใช้ JSONB แทน normalization

\- ห้าม query JSONB field บ่อยๆ โดยไม่มี index



\---



\# การออกแบบ Workflow และ Dataflow



\## Workflow การทํางานของ Dynamic Filter System



```

\[Client Request]

&#x20;     |

&#x20;     v

\[Controller Layer]

รับ query parameters: ?name=john\&status=1\&limit=10

&#x20;     |

&#x20;     v

\[Validation Layer]

ตรวจสอบ input: 

\- name: string, max 100 chars

\- status: int, 0-2

\- limit: int, 1-100

&#x20;     |

&#x20;     v

\[Service Layer]

แปลง parameters เป็น map\[string]interface{}

เรียก repository.FilterUsers(filters)

&#x20;     |

&#x20;     v

\[Repository Layer]

สร้าง GORM query

เพิ่ม WHERE clauses แบบ conditional

เพิ่ม ORDER BY, LIMIT, OFFSET

Execute query

&#x20;     |

&#x20;     v

\[Database]

PostgreSQL execute query

ใช้ indexes (ถ้ามี) เพื่อ optimize

Return result set

&#x20;     |

&#x20;     v

\[Response]

JSON response กลับ client

```



\## Dataflow สําหรับ Recursive Query



```

\[Client] -> \[API: GET /api/org/:id/subordinates]

&#x20;                        |

&#x20;                        v

&#x20;             \[Controller: GetSubordinates]

&#x20;                        |

&#x20;                        v

&#x20;             \[Service: GetTeamTree]

&#x20;                        |

&#x20;                        v

&#x20;             \[Repository: GetAllSubordinates]

&#x20;                        |

&#x20;                        v

&#x20;             \[PostgreSQL: WITH RECURSIVE]

&#x20;                        |

&#x20;   +--------------------+--------------------+

&#x20;   |                    |                    |

\[Anchor]            \[Recursive]          \[Terminate]

SELECT \*           UNION ALL            WHEN no more

WHERE id = X       JOIN tree             rows found

&#x20;   |                    |

&#x20;   +--------+-----------+

&#x20;            |

&#x20;            v

&#x20;   \[Result Set: id, name, level, path]

&#x20;            |

&#x20;            v

&#x20;   \[Build Tree Structure]

&#x20;            |

&#x20;            v

&#x20;   \[Return JSON to Client]

```



\---



\# คู่มือการทดสอบ



\## 1. Unit Test สําหรับ Repository



```go

// tests/user\_repo\_test.go

package tests



import (

&#x20;   "context"

&#x20;   "testing"

&#x20;   "github.com/stretchr/testify/assert"

&#x20;   "github.com/stretchr/testify/suite"

)



type UserRepoTestSuite struct {

&#x20;   suite.Suite

&#x20;   db       \*gorm.DB

&#x20;   repo     \*repositories.UserRepository

&#x20;   testUser \*models.User

}



func (s \*UserRepoTestSuite) SetupTest() {

&#x20;   // สร้าง test database

&#x20;   s.db = setupTestDB()

&#x20;   s.repo = repositories.NewUserRepository(s.db)

&#x20;   

&#x20;   // สร้าง test data

&#x20;   s.testUser = \&models.User{

&#x20;       Name:  "Test User",

&#x20;       Email: "test@example.com",

&#x20;       Age:   25,

&#x20;   }

&#x20;   s.db.Create(s.testUser)

}



func (s \*UserRepoTestSuite) TearDownTest() {

&#x20;   // ลบ test data

&#x20;   s.db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")

}



// ทดสอบ dynamic filter

func (s \*UserRepoTestSuite) TestFilterUsers() {

&#x20;   tests := \[]struct {

&#x20;       name    string

&#x20;       filters map\[string]interface{}

&#x20;       want    int

&#x20;       wantErr bool

&#x20;   }{

&#x20;       {

&#x20;           name: "filter by name partial match",

&#x20;           filters: map\[string]interface{}{

&#x20;               "name": "Test",

&#x20;           },

&#x20;           want: 1,

&#x20;       },

&#x20;       {

&#x20;           name: "filter by age range",

&#x20;           filters: map\[string]interface{}{

&#x20;               "min\_age": 20,

&#x20;               "max\_age": 30,

&#x20;           },

&#x20;           want: 1,

&#x20;       },

&#x20;       {

&#x20;           name: "filter with pagination",

&#x20;           filters: map\[string]interface{}{

&#x20;               "limit":  10,

&#x20;               "offset": 0,

&#x20;           },

&#x20;           want: 1,

&#x20;       },

&#x20;       {

&#x20;           name: "empty filters should return all",

&#x20;           filters: map\[string]interface{}{},

&#x20;           want:    1,

&#x20;       },

&#x20;   }

&#x20;   

&#x20;   for \_, tt := range tests {

&#x20;       s.Run(tt.name, func() {

&#x20;           users, total, err := s.repo.FilterUsers(context.Background(), tt.filters)

&#x20;           

&#x20;           if tt.wantErr {

&#x20;               s.Error(err)

&#x20;           } else {

&#x20;               s.NoError(err)

&#x20;               s.Equal(int64(tt.want), total)

&#x20;               s.NotNil(users)

&#x20;           }

&#x20;       })

&#x20;   }

}



// ทดสอบ recursive query

func (s \*UserRepoTestSuite) TestRecursiveQuery() {

&#x20;   // สร้าง hierarchy data

&#x20;   ceo := createEmployee(s.db, "CEO", nil)

&#x20;   manager := createEmployee(s.db, "Manager", \&ceo.ID)

&#x20;   staff := createEmployee(s.db, "Staff", \&manager.ID)

&#x20;   

&#x20;   subordinates, err := s.employeeRepo.GetAllSubordinates(context.Background(), ceo.ID)

&#x20;   

&#x20;   s.NoError(err)

&#x20;   s.Equal(2, len(subordinates)) // manager + staff

}



// ทดสอบ window function

func (s \*UserRepoTestSuite) TestWindowFunction() {

&#x20;   // สร้าง users with different scores

&#x20;   users := \[]models.User{

&#x20;       {Name: "User1", Score: 100},

&#x20;       {Name: "User2", Score: 90},

&#x20;       {Name: "User3", Score: 80},

&#x20;   }

&#x20;   s.db.Create(\&users)

&#x20;   

&#x20;   rankings, err := s.repo.GetUsersWithRanking(context.Background())

&#x20;   

&#x20;   s.NoError(err)

&#x20;   s.Equal(3, len(rankings))

&#x20;   s.Equal(1, rankings\[0].RankNum) // highest score rank 1

}



// รัน test suite

func TestUserRepoTestSuite(t \*testing.T) {

&#x20;   suite.Run(t, new(UserRepoTestSuite))

}

```



\## 2. Integration Test



```go

// tests/integration\_test.go

func TestIntegration\_FilterAPI(t \*testing.T) {

&#x20;   // สร้าง test server

&#x20;   router := setupRouter()

&#x20;   

&#x20;   tests := \[]struct {

&#x20;       name       string

&#x20;       query      string

&#x20;       statusCode int

&#x20;   }{

&#x20;       {

&#x20;           name:       "search by name",

&#x20;           query:      "/users?name=john",

&#x20;           statusCode: 200,

&#x20;       },

&#x20;       {

&#x20;           name:       "search with pagination",

&#x20;           query:      "/users?limit=10\&page=1",

&#x20;           statusCode: 200,

&#x20;       },

&#x20;       {

&#x20;           name:       "invalid limit",

&#x20;           query:      "/users?limit=1000",

&#x20;           statusCode: 400, // limit too high

&#x20;       },

&#x20;   }

&#x20;   

&#x20;   for \_, tt := range tests {

&#x20;       t.Run(tt.name, func(t \*testing.T) {

&#x20;           req, \_ := http.NewRequest("GET", tt.query, nil)

&#x20;           resp := performRequest(router, req)

&#x20;           assert.Equal(t, tt.statusCode, resp.Code)

&#x20;       })

&#x20;   }

}



// Benchmark test

func BenchmarkFilterUsers(b \*testing.B) {

&#x20;   db := setupTestDB()

&#x20;   repo := repositories.NewUserRepository(db)

&#x20;   

&#x20;   // seed 10000 users

&#x20;   seedUsers(db, 10000)

&#x20;   

&#x20;   filters := map\[string]interface{}{

&#x20;       "status": 1,

&#x20;       "min\_age": 18,

&#x20;       "max\_age": 30,

&#x20;   }

&#x20;   

&#x20;   b.ResetTimer()

&#x20;   for i := 0; i < b.N; i++ {

&#x20;       repo.FilterUsers(context.Background(), filters)

&#x20;   }

}

```



\---



\# คู่มือการการใช้งาน



\## การติดตั้ง



```bash

\# 1. Clone project

git clone https://github.com/yourrepo/advanced-sql-gorm.git

cd advanced-sql-gorm



\# 2. Install dependencies

go mod tidy



\# 3. Setup PostgreSQL

docker run -d \\

&#x20; --name postgres-sql \\

&#x20; -e POSTGRES\_USER=myuser \\

&#x20; -e POSTGRES\_PASSWORD=mypassword \\

&#x20; -e POSTGRES\_DB=mydb \\

&#x20; -p 5432:5432 \\

&#x20; postgres:15



\# 4. Run migrations

go run scripts/migrate.go



\# 5. Seed test data

go run scripts/seed\_data.go



\# 6. Run application

go run main.go

```



\## API Endpoints



```bash

\# 1. Dynamic Filter Users

GET /api/users?name=john\&status=1\&min\_age=18\&limit=20\&offset=0



\# Response

{

&#x20; "data": \[...],

&#x20; "total": 150,

&#x20; "limit": 20,

&#x20; "offset": 0

}



\# 2. Get User Ranking

GET /api/users/ranking



\# Response

\[

&#x20; {"id":1, "name":"John", "score":100, "rank":1},

&#x20; {"id":2, "name":"Jane", "score":95, "rank":2}

]



\# 3. Get Organization Tree

GET /api/org/:id/subordinates



\# Response

\[

&#x20; {"id":2, "name":"Manager", "level":2, "path":"CEO > Manager"},

&#x20; {"id":3, "name":"Staff", "level":3, "path":"CEO > Manager > Staff"}

]



\# 4. Search by JSON Metadata

GET /api/users/metadata?key=preferences.theme\&value=dark



\# 5. Get Top N per Department

GET /api/departments/top?n=3

```



\## ตัวอย่างการใช้งาน Client



```javascript

// Frontend: React example

const searchUsers = async (filters) => {

&#x20; const params = new URLSearchParams(filters);

&#x20; const response = await fetch(`/api/users?${params}`);

&#x20; return response.json();

};



// ใช้งาน

searchUsers({

&#x20; name: 'john',

&#x20; status: 1,

&#x20; limit: 20

}).then(data => console.log(data));

```



```python

\# Backend: Python client

import requests



def filter\_users(name=None, status=None, min\_age=None):

&#x20;   params = {}

&#x20;   if name:

&#x20;       params\['name'] = name

&#x20;   if status:

&#x20;       params\['status'] = status

&#x20;   if min\_age:

&#x20;       params\['min\_age'] = min\_age

&#x20;   

&#x20;   response = requests.get('http://localhost:8080/api/users', params=params)

&#x20;   return response.json()

```



\---



\# คู่มือการบำรุงรักษา



\## การ Monitor และ Logging



```go

// middleware/logger.go

func LoggerMiddleware() gin.HandlerFunc {

&#x20;   return func(c \*gin.Context) {

&#x20;       start := time.Now()

&#x20;       

&#x20;       // บันทึก query ก่อน execute

&#x20;       c.Next()

&#x20;       

&#x20;       // บันทึกหลังจาก execute

&#x20;       duration := time.Since(start)

&#x20;       

&#x20;       // log slow queries (>100ms)

&#x20;       if duration > 100\*time.Millisecond {

&#x20;           log.Printf("\[SLOW QUERY] %s %s took %v", 

&#x20;               c.Request.Method, 

&#x20;               c.Request.URL.Path,

&#x20;               duration)

&#x20;       }

&#x20;   }

}

```



\## การ Optimize Performance



```sql

\-- 1. สร้าง indexes ที่จําเป็น

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx\_users\_status ON users(status);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx\_users\_age ON users(age);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx\_users\_status\_age ON users(status, age);



\-- 2. Partial index for active users only

CREATE INDEX CONCURRENTLY idx\_active\_users\_name ON users(name) 

WHERE status = 1 AND deleted\_at IS NULL;



\-- 3. GIN index for JSONB

CREATE INDEX idx\_users\_metadata ON users USING gin(metadata);



\-- 4. Index for ILIKE search (pg\_trgm)

CREATE EXTENSION IF NOT EXISTS pg\_trgm;

CREATE INDEX idx\_users\_name\_trgm ON users USING gin(name gin\_trgm\_ops);

```



\## การ Backup และ Recovery



```bash

\# Backup database

pg\_dump -U myuser -h localhost mydb > backup\_$(date +%Y%m%d).sql



\# Backup specific tables

pg\_dump -U myuser -h localhost -t users -t employees mydb > backup\_users.sql



\# Restore

psql -U myuser -h localhost mydb < backup\_20240101.sql

```



\## การ Migration Guide



```go

// migrations/004\_add\_jsonb\_index.sql

\-- +migrate Up

CREATE INDEX CONCURRENTLY idx\_users\_metadata\_preferences 

ON users USING gin((metadata->'preferences'));



\-- +migrate Down

DROP INDEX CONCURRENTLY idx\_users\_metadata\_preferences;

```



\## Troubleshooting Guide



| ปัญหา | สาเหตุ | วิธีแก้ไข |

|-------|--------|----------|

| Query ช้า | Missing index | ใช้ EXPLAIN ANALYZE ตรวจสอบ |

| Recursive query infinite loop | Cycle in data | เพิ่ม cycle detection |

| JSON query ช้า | No GIN index | สร้าง GIN index บน JSONB |

| Memory leak | Too many open connections | ปรับ connection pool |

| Deadlock | Bad transaction order | ใช้ advisory lock |



\## Health Check Script



```go

// scripts/health\_check.go

func healthCheck(db \*gorm.DB) {

&#x20;   // 1. Check database connection

&#x20;   sqlDB, \_ := db.DB()

&#x20;   if err := sqlDB.Ping(); err != nil {

&#x20;       log.Fatal("Database connection failed:", err)

&#x20;   }

&#x20;   

&#x20;   // 2. Check slow queries

&#x20;   var slowQueries int

&#x20;   db.Raw(`

&#x20;       SELECT count(\*) FROM pg\_stat\_statements 

&#x20;       WHERE mean\_time > 100

&#x20;   `).Scan(\&slowQueries)

&#x20;   

&#x20;   if slowQueries > 10 {

&#x20;       log.Warn("High number of slow queries:", slowQueries)

&#x20;   }

&#x20;   

&#x20;   // 3. Check index usage

&#x20;   var unusedIndexes \[]string

&#x20;   db.Raw(`

&#x20;       SELECT indexname FROM pg\_stat\_user\_indexes 

&#x20;       WHERE idx\_scan = 0

&#x20;   `).Scan(\&unusedIndexes)

&#x20;   

&#x20;   if len(unusedIndexes) > 0 {

&#x20;       log.Info("Unused indexes found:", unusedIndexes)

&#x20;   }

}

```



\---



\## สรุป Best Practices



1\. \*\*Dynamic Filter\*\*: ใช้ map-based filter กับ conditional WHERE clauses

2\. \*\*Window Functions\*\*: ใช้แทน subquery สําหรับ ranking และ running totals

3\. \*\*Recursive Query\*\*: ใช้กับ hierarchical data แต่ต้องมี index และ cycle detection

4\. \*\*JSONB\*\*: ใช้เมื่อ schema เปลี่ยนแปลงบ่อย แต่ต้องมี GIN index

5\. \*\*Index\*\*: สร้าง indexes บน columns ที่ใช้ใน WHERE, ORDER BY, JOIN

6\. \*\*Monitoring\*\*: log slow queries และ unused indexes

7\. \*\*Testing\*\*: ทดสอบทั้ง unit และ integration รวมถึง benchmark



\---



\*\*เอกสารนี้สร้างขึ้นสําหรับเรียนรู้ SQL Advance บน GORM (GoLang)\*\*

