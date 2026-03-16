package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	userdto "github.com/vadxq/go-rest-starter/internal/core/user/dto"
	"github.com/vadxq/go-rest-starter/internal/core/user/model"
	userrepo "github.com/vadxq/go-rest-starter/internal/core/user/repository"
	"github.com/vadxq/go-rest-starter/pkg/cache"
	apperrors "github.com/vadxq/go-rest-starter/pkg/errors"
)

const (
	// คำนำหน้าคีย์แคชของผู้ใช้
	userCachePrefix = "user:"

	// คีย์แคชสำหรับรายการผู้ใช้
	userListCacheKey = "user:list"

	// อายุของแคชผู้ใช้
	userCacheTTL = 30 * time.Minute
)

// UserService อินเทอร์เฟซสำหรับบริการจัดการผู้ใช้
type UserService interface {
	CreateUser(ctx context.Context, input userdto.CreateUserInput) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, id string, input userdto.UpdateUserInput) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
}

// userService โครงสร้างสำหรับimplementบริการจัดการผู้ใช้
type userService struct {
	userRepo  userrepo.UserRepository
	validator *validator.Validate
	db        *gorm.DB
	cache     cache.Cache
}

// NewUserService สร้างinstanceใหม่ของUserService
func NewUserService(ur userrepo.UserRepository, v *validator.Validate, db *gorm.DB, c cache.Cache) UserService {
	return &userService{
		userRepo:  ur,
		validator: v,
		db:        db,
		cache:     c,
	}
}

// สร้างคีย์แคชสำหรับผู้ใช้ตามID
func getUserCacheKey(id string) string {
	return fmt.Sprintf("%s%s", userCachePrefix, id)
}

// CreateUser สร้างผู้ใช้ใหม่
func (s *userService) CreateUser(ctx context.Context, input userdto.CreateUserInput) (*model.User, error) {
	// ตรวจสอบความถูกต้องของข้อมูลนำเข้า
	if err := s.validator.Struct(input); err != nil {
		return nil, apperrors.ValidationError("ข้อมูลนำเข้าไม่ถูกต้อง", err)
	}

	// ตรวจสอบว่าอีเมลมีอยู่แล้วหรือไม่
	exists, err := s.userRepo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return nil, err // ส่งต่อข้อผิดพลาดจากrepository
	}

	if exists {
		return nil, apperrors.ConflictError("อีเมลนี้ถูกใช้งานแล้ว", nil)
	}

	// เข้ารหัสรหัสผ่าน
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.InternalError("ไม่สามารถเข้ารหัสรหัสผ่านได้", err)
	}

	user := &model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "user", // บทบาทเริ่มต้น
	}

	// เริ่มtransaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.userRepo.Create(ctx, tx, user); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err // ส่งต่อข้อผิดพลาดจากrepository
	}

	// ลบแคชรายการผู้ใช้
	_ = s.cache.Delete(ctx, userListCacheKey)

	return user, nil
}

// GetByID ดึงข้อมูลผู้ใช้ตามID
func (s *userService) GetByID(ctx context.Context, id string) (*model.User, error) {
	// พยายามดึงจากแคช
	cacheKey := getUserCacheKey(id)
	var user model.User

	err := s.cache.GetObject(ctx, cacheKey, &user)
	if err == nil {
		return &user, nil
	}

	// ไม่พบในแคช ดึงจากฐานข้อมูล
	user2, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err // ส่งต่อข้อผิดพลาดจากrepository
	}

	// เก็บลงแคช
	_ = s.cache.SetObject(ctx, cacheKey, user2, userCacheTTL)

	return user2, nil
}

// UpdateUser อัปเดตข้อมูลผู้ใช้
func (s *userService) UpdateUser(ctx context.Context, id string, input userdto.UpdateUserInput) (*model.User, error) {
	// ตรวจสอบความถูกต้องของข้อมูลนำเข้า
	if err := s.validator.Struct(input); err != nil {
		return nil, apperrors.ValidationError("ข้อมูลนำเข้าไม่ถูกต้อง", err)
	}

	// ดึงข้อมูลผู้ใช้
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err // ส่งต่อข้อผิดพลาดจากrepository
	}

	// อัปเดตฟิลด์ต่างๆ
	if input.Name != "" {
		user.Name = input.Name
	}

	if input.Email != "" && input.Email != user.Email {
		// ตรวจสอบอีเมลใหม่ว่ามีอยู่แล้วหรือไม่
		exists, err := s.userRepo.ExistsByEmail(ctx, input.Email)
		if err != nil {
			return nil, err // ส่งต่อข้อผิดพลาดจากrepository
		}

		if exists {
			return nil, apperrors.ConflictError("อีเมลนี้ถูกใช้งานแล้ว", nil)
		}

		user.Email = input.Email
	}

	if input.Password != "" {
		// เข้ารหัสรหัสผ่านใหม่
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, apperrors.InternalError("ไม่สามารถเข้ารหัสรหัสผ่านได้", err)
		}

		user.Password = string(hashedPassword)
	}

	// เริ่มtransaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.userRepo.Update(ctx, tx, user); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err // ส่งต่อข้อผิดพลาดจากrepository
	}

	// อัปเดตแคช
	cacheKey := getUserCacheKey(id)
	_ = s.cache.SetObject(ctx, cacheKey, user, userCacheTTL)

	// ลบแคชรายการผู้ใช้
	_ = s.cache.Delete(ctx, userListCacheKey)

	return user, nil
}

// DeleteUser ลบผู้ใช้
func (s *userService) DeleteUser(ctx context.Context, id string) error {
	// ดึงข้อมูลผู้ใช้
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err // ส่งต่อข้อผิดพลาดจากrepository
	}

	// เริ่มtransaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.userRepo.Delete(ctx, tx, user.ID); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err // ส่งต่อข้อผิดพลาดจากrepository
	}

	// ลบแคช
	cacheKey := getUserCacheKey(id)
	_ = s.cache.Delete(ctx, cacheKey)

	// ลบแคชรายการผู้ใช้
	_ = s.cache.Delete(ctx, userListCacheKey)

	return nil
}

// ListUsers ดึงรายการผู้ใช้แบบแบ่งหน้า
func (s *userService) ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	// สร้างคีย์แคชตามข้อมูลการแบ่งหน้า
	cacheKey := fmt.Sprintf("%s:%d:%d", userListCacheKey, page, pageSize)

	// พยายามดึงจากแคช
	var cachedResult struct {
		Users []*model.User `json:"users"`
		Total int64         `json:"total"`
	}

	err := s.cache.GetObject(ctx, cacheKey, &cachedResult)
	if err == nil {
		return cachedResult.Users, cachedResult.Total, nil
	}

	// ไม่พบในแคช ดึงจากฐานข้อมูล
	users, total, err := s.userRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err // ส่งต่อข้อผิดพลาดจากrepository
	}

	// เก็บลงแคช
	cachedResult = struct {
		Users []*model.User `json:"users"`
		Total int64         `json:"total"`
	}{
		Users: users,
		Total: total,
	}

	_ = s.cache.SetObject(ctx, cacheKey, cachedResult, userCacheTTL)

	return users, total, nil
}
/********
*
นี่โค้ด `userService` จากภาษาอังกฤษเป็นภาษาไทย พร้อมอธิบายการทำงานของแต่ละส่วน:
### คำอธิบายเพิ่มเติม:
- **Service นี้ทำหน้าที่**: จัดการหลักของผู้ใช้ (CRUD) พร้อมกับการตรวจสอบข้อมูล การเข้ารหัสรหัสผ่าน การใช้ Transaction และการจัดการ Cache
- **การจัดการข้อผิดพลาด**: ใช้ `apperrors` ในการสร้างข้อผิดพลาดที่มีความหมายเฉพาะ (Validation, Conflict, Internal)
- **การใช้งาน Cache**: ใช้ Redis หรือ cache อื่นๆ ผ่าน interface `cache.Cache` เพื่อลดภาระฐานข้อมูล โดยมี TTL 30 นาที
- **ความปลอดภัย**: รหัสผ่านถูกเข้ารหัสด้วย bcrypt ก่อนบันทึก
*/