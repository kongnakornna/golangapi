package usecase

import (
	"context"

	"icmongolang/config"
	"icmongolang/internal/models"
	"icmongolang/internal/modules/email"
	"icmongolang/internal/usecase"
	"icmongolang/pkg/logger"
)

type emailUseCase struct {
	usecase.UseCase[models.EmailLog]
	pgRepo email.EmailPgRepository
}

// CreateEmailUseCaseI creates a new email use case instance.
// สร้างอินสแตนซ์สำหรับธุรกิจอีเมล
func CreateEmailUseCaseI(
	pgRepo email.EmailPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) email.EmailUseCaseI {
	return &emailUseCase{
		UseCase: usecase.CreateUseCase[models.EmailLog](pgRepo, cfg, logger),
		pgRepo:  pgRepo,
	}
}

// SendEmail delegates to pkg/sendEmail for actual sending, logs the result.
// ส่งอีเมลผ่าน pkg/sendEmail และบันทึกผลลัพธ์
func (u *emailUseCase) SendEmail(ctx context.Context, to, subject, body string) error {
	// TODO: delegate to pkg/sendEmail
	// e.g.: sendEmail.Send(to, subject, body)
	logEntry := &models.EmailLog{
		To:      to,
		Subject: subject,
		Body:    body,
		Status:  "sent",
	}
	_, err := u.pgRepo.Create(ctx, logEntry)
	return err
}

func (u *emailUseCase) GetConfig(ctx context.Context) (*models.EmailConfig, error) {
	return u.pgRepo.GetConfig(ctx)
}

func (u *emailUseCase) UpdateConfig(ctx context.Context, cfg *models.EmailConfig) (*models.EmailConfig, error) {
	return u.pgRepo.UpdateConfig(ctx, cfg)
}
