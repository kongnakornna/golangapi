package email

import (
	"context"

	"icmongolang/internal"
	"icmongolang/internal/models"
)

// EmailPgRepository defines data access methods for email logs and config.
// ดึงข้อมูลบันทึกอีเมลและการตั้งค่าจากฐานข้อมูล
type EmailPgRepository interface {
	internal.PgRepository[models.EmailLog]
	GetConfig(ctx context.Context) (*models.EmailConfig, error)
	UpdateConfig(ctx context.Context, cfg *models.EmailConfig) (*models.EmailConfig, error)
}
