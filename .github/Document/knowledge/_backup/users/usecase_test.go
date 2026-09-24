package usecase

import (
	"context"
	"icmongolang/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	// ... imports อื่น ๆ ที่จำเป็น
)

// ฟังก์ชันทดสอบของคุณต่อจากนี้
func TestCreateUser(t *testing.T) {
    // Mock repository
    mockPgRepo := new(MockUserPgRepository)
    mockRedisRepo := new(MockUserRedisRepository)
    mockDistributor := new(MockUserRedisTaskDistributor)
    
    uc := CreateUserUseCaseI(mockPgRepo, mockRedisRepo, mockDistributor, cfg, logger)
    
    user := &models.SdUser{
        Email:    "test@example.com",
        Password: "password123",
        Fullname: stringPtr("Test User"),
    }
    
    mockPgRepo.On("Create", mock.Anything, mock.Anything).Return(user, nil)
    
    result, err := uc.CreateUser(context.Background(), user, "password123")
    
    assert.NoError(t, err)
    assert.Equal(t, user.Email, result.Email)
}

// go test -v ./internal/users/...
// rm internal/users/usecase_test.go
// # หรือใช้คำสั่ง del ใน Windows:
// del internal\users\usecase_test.go