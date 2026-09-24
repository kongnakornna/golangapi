## การทำ Transaction (Rollback/Commit) ใน GORM

Transaction ช่วยให้การดำเนินงานหลาย ๆ อย่างสำเร็จทั้งหมดหรือไม่สำเร็จเลย (all-or-nothing) เช่น การสร้างผู้ใช้ + ส่งอีเมล + อัปเดต log

### 1. วิธีที่ 1: ใช้ `db.Transaction` (recommended)

```go
// CreateUserWithTransaction สร้างผู้ใช้และบันทึก log แบบ transaction
// CreateUserWithTransaction creates user and audit log in a transaction
func (r *UserPgRepo) CreateUserWithTransaction(ctx context.Context, user *models.SdUser, logMessage string) error {
	// เริ่ม transaction
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// ขั้นตอนที่ 1: สร้างผู้ใช้
		if err := tx.Create(user).Error; err != nil {
			return err // ถ้า error จะ rollback อัตโนมัติ
		}

		// ขั้นตอนที่ 2: บันทึก audit log (สมมติมีตาราง user_logs)
		logEntry := UserLog{
			UserID:    user.ID,
			Message:   logMessage,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&logEntry).Error; err != nil {
			return err // rollback ทั้งสอง operation
		}

		// commit อัตโนมัติเมื่อ return nil
		return nil
	})
}
```

### 2. วิธีที่ 2: จัดการ transaction ด้วยตนเอง (manual commit/rollback)

```go
// UpdateUserWithManualTransaction อัปเดตผู้ใช้หลายคนแบบ manual transaction
// UpdateUserWithManualTransaction updates multiple users with manual transaction control
func (r *UserPgRepo) UpdateUserWithManualTransaction(ctx context.Context, userIDs []string, newStatus int16) error {
	// เริ่ม transaction
	tx := r.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// ทำการอัปเดต
	if err := tx.Model(&models.SdUser{}).Where("id IN ?", userIDs).Update("status", newStatus).Error; err != nil {
		tx.Rollback()
		return err
	}

	// ตรวจสอบเงื่อนไขเพิ่มเติม (ตัวอย่าง: ต้องมีผู้ใช้อย่างน้อย 1 คน)
	var count int64
	tx.Model(&models.SdUser{}).Where("status = ?", newStatus).Count(&count)
	if count == 0 {
		tx.Rollback()
		return errors.New("no users updated")
	}

	// commit transaction
	return tx.Commit().Error
}
```

### 3. ตัวอย่าง method ใน `UserPgRepo` สำหรับ transaction ที่ใช้บ่อย

```go
// TransferUserRole เปลี่ยน role ของผู้ใช้และบันทึกประวัติ (transaction)
// TransferUserRole changes user role and logs history in a transaction
func (r *UserPgRepo) TransferUserRole(ctx context.Context, userID uuid.UUID, newRoleID int, performedBy uuid.UUID) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. อัปเดต role ของ user
		if err := tx.Model(&models.SdUser{}).Where("id = ?", userID).Update("role_id", newRoleID).Error; err != nil {
			return fmt.Errorf("update role failed: %w", err)
		}

		// 2. บันทึกประวัติ (สมมติตาราง role_history)
		history := map[string]interface{}{
			"user_id":       userID,
			"old_role_id":   0, // ควร query ก่อน หรือส่งมา
			"new_role_id":   newRoleID,
			"performed_by":  performedBy,
			"changed_at":    time.Now(),
		}
		if err := tx.Table("role_histories").Create(history).Error; err != nil {
			return fmt.Errorf("create history failed: %w", err)
		}
		return nil
	})
}
```

### 4. การใช้ transaction ร่วมกับ repository อื่น ๆ (cross-repo)

ใน usecase layer สามารถเรียกหลาย repo ภายใน transaction เดียวกันได้ โดยส่ง `*gorm.DB` ผ่าน context หรือใช้ `db.Transaction` ที่ระดับ usecase

```go
// ใน usecase
func (u *userUseCase) CreateUserWithOrg(ctx context.Context, user *models.SdUser, orgName string) error {
	return u.pgRepo.DB.Transaction(func(tx *gorm.DB) error {
		// สร้าง user โดยใช้ tx
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		// สร้าง org โดยใช้ tx (ผ่าน orgRepo ที่รับ tx)
		org := &models.Organization{Name: orgName, OwnerID: user.ID}
		if err := u.orgRepo.CreateWithTx(ctx, tx, org); err != nil {
			return err
		}
		return nil
	})
}
```

### 5. ตัวอย่างการ rollback แบบมีเงื่อนไข

```go
// ConditionalTransaction แสดงการ rollback ตามเงื่อนไขทางธุรกิจ
// ConditionalTransaction demonstrates rollback based on business rules
func (r *UserPgRepo) ConditionalTransaction(ctx context.Context, userID uuid.UUID, points int) error {
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// ดึงคะแนนปัจจุบัน (สมมติฟิลด์ loyalty_points)
	var currentPoints int
	if err := tx.Model(&models.SdUser{}).Where("id = ?", userID).Select("loyalty_points").Scan(&currentPoints).Error; err != nil {
		tx.Rollback()
		return err
	}

	// ตรวจสอบเงื่อนไข: ถ้าคะแนนใหม่ติดลบ ให้ rollback
	if currentPoints+points < 0 {
		tx.Rollback()
		return errors.New("insufficient points")
	}

	// อัปเดตคะแนน
	if err := tx.Model(&models.SdUser{}).Where("id = ?", userID).Update("loyalty_points", currentPoints+points).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
```

### สรุป

| การทำงาน | Method ที่ใช้ |
|-----------|---------------|
| Auto rollback/commit | `db.Transaction(func(tx *gorm.DB) error { ... })` |
| Manual control | `tx := db.Begin()`, `tx.Commit()`, `tx.Rollback()` |
| ส่ง tx ไปยัง repo อื่น | ออกแบบ method ให้รับ `*gorm.DB` เป็น parameter |
| Rollback เมื่อ panic | ใช้ `defer` + `recover()` ร่วมกับ `tx.Rollback()` |

**คำแนะนำ:** ใช้ `db.Transaction` เป็นหลัก เพราะจัดการ rollback/commit ให้อัตโนมัติ และป้องกันการลืม commit/rollback

ต้องการให้เพิ่ม method ที่ใช้ transaction ใน `UserPgRepo` ของคุณโดยเฉพาะไหมครับ? บอก use case มาได้เลย