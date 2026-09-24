package usecase

import (
	"context"

	"icmongolang/config"
	"icmongolang/internal/models"
	"icmongolang/internal/modules/document"
	"icmongolang/internal/usecase"
	"icmongolang/pkg/logger"

)

type documentUseCase struct {
	usecase.UseCase[models.Document]
	pgRepo document.DocumentPgRepository
}

// CreateDocumentUseCaseI creates a new document use case instance.
// สร้างอินสแตนซ์สำหรับธุรกิจเอกสาร
func CreateDocumentUseCaseI(
	pgRepo document.DocumentPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) document.DocumentUseCaseI {
	return &documentUseCase{
		UseCase: usecase.CreateUseCase[models.Document](pgRepo, cfg, logger),
		pgRepo:  pgRepo,
	}
}

func (u *documentUseCase) GetByFilename(ctx context.Context, filename string) (*models.Document, error) {
	return u.pgRepo.GetByFilename(ctx, filename)
}
