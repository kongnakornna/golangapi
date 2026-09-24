\## การนำเทคนิค SQL ขั้นสูงมาใช้กับ GORM (Go)



จากคำถามเรื่องการ filter ผู้ใช้ตามฟิลด์ต่าง ๆ และการใช้ `CASE WHEN`, `IF-ELSE`, `dynamic WHERE` ใน SQL เราจะมาแปลงเป็น \*\*GORM\*\* ซึ่งเป็น ORM ที่ใช้ในโปรเจกต์ของคุณอยู่แล้ว



\---



\### 1. การใช้ `CASE WHEN` ใน GORM (คล้าย SQL)



GORM รองรับ `CASE` ผ่าน `gorm.Expr` หรือ raw SQL



\#### ตัวอย่าง: ดึงผู้ใช้พร้อมสถานะเป็นข้อความ (Active/Inactive)



```go

type UserWithStatusText struct {

&#x20;   ID         uuid.UUID

&#x20;   Email      string

&#x20;   Status     int16

&#x20;   StatusText string

}



var results \[]UserWithStatusText

db.Model(\&models.SdUser{}).

&#x20;   Select("id, email, status, " +

&#x20;       "CASE WHEN status = 1 THEN 'Active' ELSE 'Inactive' END AS status\_text").

&#x20;   Scan(\&results)

```



\#### ตัวอย่าง: ใช้ `CASE` ใน `ORDER BY` (เรียงลำดับตามเงื่อนไข)



```go

// เรียงลำดับ: verified users ก่อน, แล้วตามด้วย created\_date

db.Model(\&models.SdUser{}).

&#x20;   Order(gorm.Expr("CASE WHEN verified = true THEN 0 ELSE 1 END")).

&#x20;   Order("created\_date DESC").

&#x20;   Find(\&users)

```



\#### ตัวอย่าง: ใช้ `CASE` ใน `UPDATE` (อัปเดตตามเงื่อนไข)



```go

// อัปเดต role\_id: ถ้า is\_superuser = true ให้ role\_id = 1, ถ้าไม่ให้คงเดิม

db.Model(\&models.SdUser{}).

&#x20;   Where("is\_superuser = true").

&#x20;   Update("role\_id", gorm.Expr("CASE WHEN is\_superuser = true THEN 1 ELSE role\_id END"))

```



\---



\### 2. Dynamic Filter (WHERE clause แบบมีเงื่อนไข) – วิธีที่ดีที่สุดใน GORM



ใช้ \*\*method chaining\*\* และ \*\*conditional where\*\* ตาม pattern ที่แนะนำ



\#### ตัวอย่าง: ฟังก์ชัน FilterUsers รองรับพารามิเตอร์หลายตัว (email, status, role\_id, verified, pagination)



```go

func (r \*UserPgRepo) FilterUsers(ctx context.Context, req \*FilterRequest) (\[]\*models.SdUser, int64, error) {

&#x20;   var users \[]\*models.SdUser

&#x20;   var total int64



&#x20;   query := r.DB.WithContext(ctx).Model(\&models.SdUser{})



&#x20;   // Dynamic where clauses

&#x20;   if req.Email != "" {

&#x20;       query = query.Where("email ILIKE ?", "%"+req.Email+"%")

&#x20;   }

&#x20;   if req.Status != nil {

&#x20;       query = query.Where("status = ?", \*req.Status)

&#x20;   }

&#x20;   if req.RoleID != nil {

&#x20;       query = query.Where("role\_id = ?", \*req.RoleID)

&#x20;   }

&#x20;   if req.Verified != nil {

&#x20;       query = query.Where("verified = ?", \*req.Verified)

&#x20;   }

&#x20;   if req.IsSuperUser != nil {

&#x20;       query = query.Where("is\_superuser = ?", \*req.IsSuperUser)

&#x20;   }

&#x20;   if req.LocationID != nil {

&#x20;       query = query.Where("location\_id = ?", \*req.LocationID)

&#x20;   }



&#x20;   // Count total before pagination

&#x20;   if err := query.Count(\&total).Error; err != nil {

&#x20;       return nil, 0, err

&#x20;   }



&#x20;   // Sorting (default by created\_date desc)

&#x20;   sortField := req.SortField

&#x20;   if sortField == "" {

&#x20;       sortField = "created\_date"

&#x20;   }

&#x20;   sortOrder := req.SortOrder

&#x20;   if sortOrder == "" {

&#x20;       sortOrder = "DESC"

&#x20;   }

&#x20;   query = query.Order(sortField + " " + sortOrder)



&#x20;   // Pagination

&#x20;   if req.Limit > 0 {

&#x20;       query = query.Limit(req.Limit).Offset(req.Offset)

&#x20;   }



&#x20;   err := query.Find(\&users).Error

&#x20;   return users, total, err

}



// FilterRequest struct

type FilterRequest struct {

&#x20;   Email       string

&#x20;   Status      \*int16

&#x20;   RoleID      \*int

&#x20;   Verified    \*bool

&#x20;   IsSuperUser \*bool

&#x20;   LocationID  \*string

&#x20;   SortField   string

&#x20;   SortOrder   string

&#x20;   Limit       int

&#x20;   Offset      int

}

```



\---



\### 3. การใช้ `FOR` / `WHILE` loop (ทำใน Go แทน SQL Loop)



ใน GORM เราไม่ควรทำ loop ใน SQL แต่ให้ทำใน Go แทน เช่น การอัปเดตทีละแถว หรือ generate token



\#### ตัวอย่าง: Generate unique verification code (loop จนไม่ซ้ำ)



```go

func (r \*UserPgRepo) GenerateUniqueVerificationCode(ctx context.Context) (string, error) {

&#x20;   for {

&#x20;       code := secureRandom.RandomHex(16) // ใช้ฟังก์ชันที่มีอยู่แล้ว

&#x20;       var count int64

&#x20;       err := r.DB.WithContext(ctx).Model(\&models.SdUser{}).

&#x20;           Where("verification\_code = ?", code).

&#x20;           Count(\&count).Error

&#x20;       if err != nil {

&#x20;           return "", err

&#x20;       }

&#x20;       if count == 0 {

&#x20;           return code, nil

&#x20;       }

&#x20;   }

}

```



\#### ตัวอย่าง: Bulk update ทีละแถว (แต่ควรใช้ `UPDATE ... WHERE` แทนถ้าเป็นไปได้)



```go

// ไม่แนะนำ: loop ทีละ record

for \_, user := range users {

&#x20;   db.Model(\&user).Update("status", 1)

}



// แนะนำ: update เป็น batch

db.Model(\&models.SdUser{}).Where("verified = ?", true).Update("status", 1)

```



\---



\### 4. Recursive Query (WITH RECURSIVE) ใน GORM



ถ้าต้องการ hierarchical query (เช่น user มี manager\_id) ต้องใช้ raw SQL เพราะ GORM ไม่มี native recursive



```go

type UserTree struct {

&#x20;   ID    uuid.UUID

&#x20;   Email string

&#x20;   Level int

}



var tree \[]UserTree

rawSQL := `

&#x20;   WITH RECURSIVE user\_tree AS (

&#x20;       SELECT id, email, 1 AS level

&#x20;       FROM sd\_user

&#x20;       WHERE id = ?

&#x20;       UNION ALL

&#x20;       SELECT u.id, u.email, ut.level + 1

&#x20;       FROM sd\_user u

&#x20;       JOIN user\_tree ut ON u.manager\_id = ut.id

&#x20;   )

&#x20;   SELECT id, email, level FROM user\_tree

`

db.Raw(rawSQL, userID).Scan(\&tree)

```



\---



\### 5. JSONB Query ใน GORM (ถ้ามี metadata field)



สมมติคุณเพิ่มฟิลด์ `metadata` type `JSONB`



```go

// ค้นหาผู้ใช้ที่มี metadata->>'city' = 'Bangkok'

var users \[]models.SdUser

db.Where("metadata->>? = ?", "city", "Bangkok").Find(\&users)



// ค้นหาผู้ใช้ที่มีคะแนน > 100 (ต้อง cast)

db.Where("(metadata->>?)::int > ?", "score", 100).Find(\&users)

```



\---



\### 6. Window Functions (ROW\_NUMBER, RANK) ใน GORM



ใช้ `ROW\_NUMBER()` เพื่อทำ complex pagination หรือ ranking



```go

type RankedUser struct {

&#x20;   models.SdUser

&#x20;   RowNum int64

}



var ranked \[]RankedUser

db.Table("sd\_user").

&#x20;   Select("\*, ROW\_NUMBER() OVER (ORDER BY created\_date DESC) as row\_num").

&#x20;   Scan(\&ranked)

```



\---



\### 7. สร้าง View แล้วใช้ GORM query ปกติ



ถ้า filter ซับซ้อนและใช้บ่อย ให้สร้าง View ใน PostgreSQL แล้ว map เป็น model ใน GORM



```sql

CREATE VIEW active\_user\_summary AS

SELECT id, email, fullname, role\_id, created\_date

FROM sd\_user

WHERE status = 1 AND verified = true;

```



ใน Go:



```go

type ActiveUserSummary struct {

&#x20;   ID          uuid.UUID

&#x20;   Email       string

&#x20;   Fullname    \*string

&#x20;   RoleID      int

&#x20;   CreatedDate time.Time

}



func (ActiveUserSummary) TableName() string {

&#x20;   return "active\_user\_summary"

}



// แล้วใช้ query ปกติ

var summaries \[]ActiveUserSummary

db.Find(\&summaries)

```



\---



\### 8. ตัวอย่างการรวม dynamic filter + pagination + sorting ครบวงจร (ใช้ใน Handler)



```go

// ใน repository method

func (r \*UserPgRepo) GetFilteredUsers(ctx context.Context, filter map\[string]interface{}, page, pageSize int, sortBy, sortOrder string) (\[]\*models.SdUser, int64, error) {

&#x20;   var users \[]\*models.SdUser

&#x20;   var total int64



&#x20;   query := r.DB.WithContext(ctx).Model(\&models.SdUser{})



&#x20;   // Dynamic filters

&#x20;   if email, ok := filter\["email"]; ok \&\& email != "" {

&#x20;       query = query.Where("email ILIKE ?", "%"+email.(string)+"%")

&#x20;   }

&#x20;   if status, ok := filter\["status"]; ok {

&#x20;       query = query.Where("status = ?", status)

&#x20;   }

&#x20;   if roleID, ok := filter\["role\_id"]; ok {

&#x20;       query = query.Where("role\_id = ?", roleID)

&#x20;   }

&#x20;   // ... เพิ่ม filter ตามต้องการ



&#x20;   // Count total

&#x20;   if err := query.Count(\&total).Error; err != nil {

&#x20;       return nil, 0, err

&#x20;   }



&#x20;   // Pagination

&#x20;   offset := (page - 1) \* pageSize

&#x20;   query = query.Limit(pageSize).Offset(offset)



&#x20;   // Sorting (ป้องกัน SQL injection โดย mapping field)

&#x20;   allowedSortFields := map\[string]bool{"created\_date": true, "email": true, "status": true}

&#x20;   if allowedSortFields\[sortBy] {

&#x20;       query = query.Order(sortBy + " " + sortOrder)

&#x20;   } else {

&#x20;       query = query.Order("created\_date DESC")

&#x20;   }



&#x20;   err := query.Find(\&users).Error

&#x20;   return users, total, err

}

```



\---



\### สรุปเทคนิค GORM ที่ควรใช้



| SQL เทคนิค | วิธีทำใน GORM |

|-------------|----------------|

| `CASE WHEN` | `gorm.Expr` หรือ `Select` พร้อม raw SQL |

| Dynamic `WHERE` | Method chaining + conditional |

| Loop / Recursive | ทำใน Go หรือใช้ raw SQL |

| Window functions | `Select` + raw SQL fragment |

| JSONB | ใช้ `Where("field->>key = ?", value)` |

| View | สร้าง View ใน DB แล้ว map เป็น struct |

| Bulk update | `Model().Where().Update()` |



การนำเทคนิคเหล่านี้ไปใช้จะช่วยให้ระบบ filter ผู้ใช้ของคุณยืดหยุ่นและมีประสิทธิภาพมากขึ้น โดยยังคงใช้ GORM เป็นหลัก



##### ถ้าต้องการตัวอย่างเพิ่มเติมสำหรับ use case เฉพาะ (เช่น filter ตามช่วงวันที่, การเรียงลำดับหลายเงื่อนไข, หรือการ join กับตารางอื่น) บอกได้เลยครับ

