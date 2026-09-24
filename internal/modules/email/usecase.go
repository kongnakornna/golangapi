package email

import (
	"context"

	"icmongolang/internal"
	"icmongolang/internal/models"
)

// EmailUseCaseI defines business logic methods for email.
// อินเทอร์เฟซธุรกิจสำหรับอีเมล
type EmailUseCaseI interface {
	internal.UseCaseI[models.EmailLog]
	SendEmail(ctx context.Context, to, subject, body string) error
	GetConfig(ctx context.Context) (*models.EmailConfig, error)
	UpdateConfig(ctx context.Context, cfg *models.EmailConfig) (*models.EmailConfig, error)
}
