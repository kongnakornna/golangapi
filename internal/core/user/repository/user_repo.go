package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/kongnakornna/golangapi/internal/core/user/model"
	apperrors "github.com/kongnakornna/golangapi/pkg/errors"
)

// UserRepository กำหนดอินเทอร์เฟซของคลังข้อมูลผู้ใช้
type UserRepository interface {
	Create(ctx context.Context, tx *gorm.DB, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	//GetByUsername(ctx context.Context, Username string) (*model.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Update(ctx context.Context, tx *gorm.DB, user *model.User) error
	Delete(ctx context.Context, tx *gorm.DB, id uint) error
	List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository สร้างอินสแตนซ์ใหม่ของ UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create สร้างผู้ใช้
func (r *userRepository) Create(ctx context.Context, tx *gorm.DB, user *model.User) error {
	result := tx.WithContext(ctx).Create(user)
	if result.Error != nil {
		return apperrors.InternalError("สร้างผู้ใช้ล้มเหลว", result.Error)
	}
	return nil
}

// GetByID ดึงข้อมูลผู้ใช้ตาม ID
func (r *userRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, apperrors.NotFoundError("ผู้ใช้", result.Error)
		}
		return nil, apperrors.InternalError("ดึงข้อมูลผู้ใช้ล้มเหลว", result.Error)
	}
	return &user, nil
}

// GetByEmail ดึงข้อมูลผู้ใช้ตามอีเมล
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, apperrors.NotFoundError("ผู้ใช้", result.Error)
		}
		return nil, apperrors.InternalError("ดึงข้อมูลผู้ใช้ล้มเหลว", result.Error)
	}
	return &user, nil
}

// ExistsByEmail ตรวจสอบว่ามีอีเมลนี้อยู่ในระบบหรือไม่
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, apperrors.InternalError("ตรวจสอบอีเมลล้มเหลว", result.Error)
	}
	return count > 0, nil
}

// Update อัปเดตข้อมูลผู้ใช้
func (r *userRepository) Update(ctx context.Context, tx *gorm.DB, user *model.User) error {
	result := tx.WithContext(ctx).Save(user)
	if result.Error != nil {
		return apperrors.InternalError("อัปเดตผู้ใช้ล้มเหลว", result.Error)
	}
	return nil
}

// Delete ลบผู้ใช้
func (r *userRepository) Delete(ctx context.Context, tx *gorm.DB, id uint) error {
	result := tx.WithContext(ctx).Delete(&model.User{}, id)
	if result.Error != nil {
		return apperrors.InternalError("ลบผู้ใช้ล้มเหลว", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.NotFoundError("ผู้ใช้", nil)
	}
	return nil
}

// List ดึงรายการผู้ใช้
func (r *userRepository) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var users []*model.User
	result := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&users)
	if result.Error != nil {
		return nil, 0, apperrors.InternalError("ดึงรายการผู้ใช้ล้มเหลว", result.Error)
	}

	var total int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, apperrors.InternalError("ดึงจำนวนผู้ใช้ทั้งหมดล้มเหลว", err)
	}

	return users, total, nil
}