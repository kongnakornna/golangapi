## การสร้างฟังก์ชัน Insert แบบแยกตัวแปร (Parameterized Insert)

ใน repository layer เราสามารถสร้าง method ที่รับค่าฟิลด์ต่าง ๆ แบบแยกกัน (ไม่ต้องส่งทั้ง struct) เพื่อเพิ่มความยืดหยุ่นและลดการสร้าง object ที่ไม่จำเป็น ตัวอย่างด้านล่างจะแสดง:

1. **InsertUserWithFields** – รับค่าฟิลด์สำคัญแบบแยกตัวแปร สร้าง `models.SdUser` แล้วบันทึกลงฐานข้อมูล
2. **InsertUserWithMap** – รับค่าเป็น `map[string]interface{}` สำหรับกรณีที่ต้องการ dynamic fields
3. **BatchInsertUsers** – แทรกผู้ใช้หลายคนในครั้งเดียว (batch insert)

### เพิ่ม method ต่อไปนี้ใน `internal/modules/users/repository/pg_repository.go` (ต่อจากฟังก์ชันเดิม)

```go
// InsertUserWithFields creates a new user using separate parameters.
// InsertUserWithFields สร้างผู้ใช้ใหม่โดยรับค่าฟิลด์แบบแยกตัวแปร
func (r *UserPgRepo) InsertUserWithFields(
	ctx context.Context,
	email string,
	passwordHash string,
	roleID int,
	firstname, lastname, fullname *string,
	status int16,
	isSuperUser bool,
	verified bool,
) (*models.SdUser, error) {
	user := &models.SdUser{
		Email:       email,
		Password:    passwordHash,
		RoleID:      roleID,
		Firstname:   firstname,
		Lastname:    lastname,
		Fullname:    fullname,
		Status:      status,
		IsSuperUser: isSuperUser,
		Verified:    verified,
	}
	// กำหนด Username ถ้ายังไม่มี
	if user.Username == "" {
		user.Username = email
	}
	// บันทึกลงฐานข้อมูล
	if err := r.DB.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// InsertUserWithMap creates a new user using a map of field names to values.
// InsertUserWithMap สร้างผู้ใช้ใหม่โดยใช้ map ชื่อฟิลด์ -> ค่า
func (r *UserPgRepo) InsertUserWithMap(ctx context.Context, fields map[string]interface{}) (*models.SdUser, error) {
	user := &models.SdUser{}
	// เติมค่าเริ่มต้นที่จำเป็น
	if email, ok := fields["email"]; ok {
		user.Email = email.(string)
		user.Username = email.(string)
	}
	if password, ok := fields["password"]; ok {
		user.Password = password.(string)
	}
	if roleID, ok := fields["role_id"]; ok {
		user.RoleID = roleID.(int)
	}
	if firstname, ok := fields["firstname"]; ok {
		user.Firstname = firstname.(*string)
	}
	if lastname, ok := fields["lastname"]; ok {
		user.Lastname = lastname.(*string)
	}
	if fullname, ok := fields["fullname"]; ok {
		user.Fullname = fullname.(*string)
	}
	if status, ok := fields["status"]; ok {
		user.Status = status.(int16)
	} else {
		user.Status = 1 // default active
	}
	if isSuper, ok := fields["is_superuser"]; ok {
		user.IsSuperUser = isSuper.(bool)
	}
	if verified, ok := fields["verified"]; ok {
		user.Verified = verified.(bool)
	}
	// ใช้ GORM Create จาก struct
	if err := r.DB.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// BatchInsertUsers inserts multiple users in a single database call.
// BatchInsertUsers แทรกผู้ใช้หลายคนในคำสั่งเดียว
func (r *UserPgRepo) BatchInsertUsers(ctx context.Context, users []*models.SdUser, batchSize int) error {
	if len(users) == 0 {
		return nil
	}
	// GORM CreateInBatches จะแบ่งเป็นชุดย่อยตาม batchSize
	return r.DB.WithContext(ctx).CreateInBatches(users, batchSize).Error
}
```

### ตัวอย่างการเรียกใช้ใน UseCase หรือ Service

```go
// ใน usecase (internal/modules/users/usecase/usecase.go)
func (u *userUseCase) CreateUserWithParams(ctx context.Context, email, rawPassword string, roleID int) (*models.SdUser, error) {
	// hash password
	hashed, err := cryptpass.HashPassword(rawPassword)
	if err != nil {
		return nil, err
	}
	// เรียก repository แบบแยกตัวแปร
	user, err := u.pgRepo.InsertUserWithFields(
		ctx,
		email,
		hashed,
		roleID,
		nil, nil, nil, // firstname, lastname, fullname
		1,   // status active
		false, // not superuser
		false, // not verified yet
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
```

### ข้อดีของการแยกตัวแปร

- **ชัดเจน** – รู้ว่าต้องส่งค่าอะไรบ้าง โดยไม่ต้องสร้าง struct เปล่าแล้ว assign ทีละ field
- **ยืดหยุ่น** – สามารถปรับค่าเริ่มต้นภายใน method ได้ (เช่น status, verified)
- **ลด dependency** – caller ไม่จำเป็นต้องรู้จัก `models.SdUser` ครบทุก field
- **เหมาะกับ API ที่รับพารามิเตอร์ไม่กี่ตัว** – ลดโค้ดที่ต้อง map ข้อมูล

### ข้อควรระวัง

- ถ้ามีหลายฟิลด์ optional การใช้ `map[string]interface{}` จะสะดวกกว่า แต่จะสูญเสีย type safety
- แนะนำให้ใช้ `InsertUserWithFields` สำหรับกรณีที่ฟิลด์สำคัญมีจำนวนจำกัด และใช้ `InsertUserWithMap` เมื่อต้องการ dynamic fields จริง ๆ

### การใช้งานร่วมกับ Generic `PgRepo.Create`

คุณสามารถเรียก `r.PgRepo.Create(ctx, user)` แทน `r.DB.Create(user)` ได้เช่นกัน เพราะ `PgRepo` มี method `Create` อยู่แล้ว แต่วิธีข้างต้นใช้ GORM โดยตรงเพื่อความเรียบง่ายและเห็นตัวอย่างการ insert แบบแยกตัวแปร

---

**หมายเหตุ:** อย่าลืม import `"gorm.io/gorm"` (ถ้ายังไม่มี) และตรวจสอบว่า `models.SdUser` มี struct tags ที่ถูกต้องสำหรับ GORM