package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	userdto "github.com/kongnakornna/golangapi/internal/core/user/dto"
	"github.com/kongnakornna/golangapi/internal/core/user/model"
	apperrors "github.com/kongnakornna/golangapi/pkg/errors"
)

// MockUserRepository คือการจำลองการทำงานของ UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, tx *gorm.DB, user *model.User) error {
	args := m.Called(ctx, tx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, tx *gorm.DB, user *model.User) error {
	args := m.Called(ctx, tx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, tx *gorm.DB, id uint) error {
	args := m.Called(ctx, tx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

// MockDB คือการจำลองการทำงานของ gorm.DB
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Transaction(fc func(tx *gorm.DB) error) error {
	args := m.Called(fc)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	// เรียกใช้ฟังก์ชัน callback
	return fc(nil)
}

// MockCache คือการจำลองการทำงานของ Cache
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(ctx context.Context, key string) ([]byte, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCache) SetObject(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCache) GetObject(ctx context.Context, key string, dest interface{}) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func (m *MockCache) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockCache) Clear(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	return db
}

func TestUserService_CreateUser(t *testing.T) {
	// ตั้งค่าข้อมูลทดสอบ
	mockRepo := new(MockUserRepository)
	mockCache := new(MockCache)
	validator := validator.New()
	db := newTestDB(t)

	service := NewUserService(mockRepo, validator, db, mockCache)

	ctx := context.Background()
	input := userdto.CreateUserInput{
		Name:     "ผู้ใช้ทดสอบ",
		Email:    "test@example.com",
		Password: "password123",
	}

	// ทดสอบการสร้างผู้ใช้สำเร็จ
	t.Run("สำเร็จ", func(t *testing.T) {
		// ตั้งค่าความคาดหวัง
		mockRepo.On("ExistsByEmail", ctx, input.Email).Return(false, nil)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*gorm.DB"), mock.AnythingOfType("*model.User")).Return(nil)
		mockCache.On("Delete", ctx, userListCacheKey).Return(nil)

		// ดำเนินการทดสอบ
		user, err := service.CreateUser(ctx, input)

		// ตรวจสอบผลลัพธ์
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, input.Name, user.Name)
		assert.Equal(t, input.Email, user.Email)
		assert.Equal(t, "user", user.Role)

		// ตรวจสอบว่ารหัสผ่านถูกเข้ารหัสอย่างถูกต้อง
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		assert.NoError(t, err)

		// ตรวจสอบการเรียกใช้ mock
		mockRepo.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	// ทดสอบกรณีอีเมลมีอยู่แล้ว
	t.Run("อีเมลมีอยู่แล้ว", func(t *testing.T) {
		mockRepo2 := new(MockUserRepository)
		service2 := NewUserService(mockRepo2, validator, db, mockCache)

		// ตั้งค่าความคาดหวัง
		mockRepo2.On("ExistsByEmail", ctx, input.Email).Return(true, nil)

		// ดำเนินการทดสอบ
		user, err := service2.CreateUser(ctx, input)

		// ตรวจสอบผลลัพธ์
		assert.Error(t, err)
		assert.Nil(t, user)

		appErr, ok := err.(*apperrors.Error)
		assert.True(t, ok)
		assert.Equal(t, apperrors.ErrorTypeConflict, appErr.Type)

		// ตรวจสอบการเรียกใช้ mock
		mockRepo2.AssertExpectations(t)
	})

	// ทดสอบกรณีข้อมูลไม่ผ่านการตรวจสอบ
	t.Run("ข้อผิดพลาดการตรวจสอบ", func(t *testing.T) {
		mockRepo3 := new(MockUserRepository)
		service3 := NewUserService(mockRepo3, validator, db, mockCache)

		invalidInput := userdto.CreateUserInput{
			Name:     "", // ชื่อว่างควรจะไม่ผ่าน
			Email:    "อีเมลไม่ถูกต้อง",
			Password: "123", // รหัสผ่านสั้นเกินไป
		}

		// ดำเนินการทดสอบ
		user, err := service3.CreateUser(ctx, invalidInput)

		// ตรวจสอบผลลัพธ์
		assert.Error(t, err)
		assert.Nil(t, user)

		appErr, ok := err.(*apperrors.Error)
		assert.True(t, ok)
		assert.Equal(t, apperrors.ErrorTypeValidation, appErr.Type)
	})
}

func TestUserService_GetByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCache := new(MockCache)
	validator := validator.New()
	db := newTestDB(t)
	service := NewUserService(mockRepo, validator, db, mockCache)

	ctx := context.Background()
	userID := "1"
	expectedUser := &model.User{
		Name:  "ผู้ใช้ทดสอบ",
		Email: "test@example.com",
		Role:  "user",
	}
	expectedUser.ID = 1

	// ทดสอบกรณีเจอข้อมูลในแคช
	t.Run("เจอแคช", func(t *testing.T) {
		cacheKey := getUserCacheKey(userID)
		mockCache.On("GetObject", ctx, cacheKey, mock.AnythingOfType("*model.User")).Return(nil).Run(func(args mock.Arguments) {
			user := args[2].(*model.User)
			*user = *expectedUser
		})

		// ดำเนินการทดสอบ
		user, err := service.GetByID(ctx, userID)

		// ตรวจสอบผลลัพธ์
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.Name, user.Name)
		assert.Equal(t, expectedUser.Email, user.Email)

		// ตรวจสอบการเรียกใช้ mock
		mockCache.AssertExpectations(t)
	})

	// ทดสอบกรณีไม่เจอแคช แต่ดึงจากฐานข้อมูลสำเร็จ
	t.Run("ไม่เจอแคช ดึงจากฐานข้อมูลสำเร็จ", func(t *testing.T) {
		mockRepo2 := new(MockUserRepository)
		mockCache2 := new(MockCache)
		service2 := NewUserService(mockRepo2, validator, db, mockCache2)

		cacheKey := getUserCacheKey(userID)

		// ตั้งค่าความคาดหวัง
		mockCache2.On("GetObject", ctx, cacheKey, mock.AnythingOfType("*model.User")).Return(errors.New("ไม่เจอแคช"))
		mockRepo2.On("GetByID", ctx, userID).Return(expectedUser, nil)
		mockCache2.On("SetObject", ctx, cacheKey, expectedUser, userCacheTTL).Return(nil)

		// ดำเนินการทดสอบ
		user, err := service2.GetByID(ctx, userID)

		// ตรวจสอบผลลัพธ์
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.Name, user.Name)
		assert.Equal(t, expectedUser.Email, user.Email)

		// ตรวจสอบการเรียกใช้ mock
		mockRepo2.AssertExpectations(t)
		mockCache2.AssertExpectations(t)
	})

	// ทดสอบกรณีไม่พบผู้ใช้
	t.Run("ไม่พบผู้ใช้", func(t *testing.T) {
		mockRepo3 := new(MockUserRepository)
		mockCache3 := new(MockCache)
		service3 := NewUserService(mockRepo3, validator, db, mockCache3)

		cacheKey := getUserCacheKey(userID)

		// ตั้งค่าความคาดหวัง
		mockCache3.On("GetObject", ctx, cacheKey, mock.AnythingOfType("*model.User")).Return(errors.New("ไม่เจอแคช"))
		mockRepo3.On("GetByID", ctx, userID).Return(nil, apperrors.NotFoundError("ผู้ใช้", nil))

		// ดำเนินการทดสอบ
		user, err := service3.GetByID(ctx, userID)

		// ตรวจสอบผลลัพธ์
		assert.Error(t, err)
		assert.Nil(t, user)

		appErr, ok := err.(*apperrors.Error)
		assert.True(t, ok)
		assert.Equal(t, apperrors.ErrorTypeNotFound, appErr.Type)

		// ตรวจสอบการเรียกใช้ mock
		mockRepo3.AssertExpectations(t)
		mockCache3.AssertExpectations(t)
	})
}