package document

import (
	"context"

	"icmongolang/internal"
	"icmongolang/internal/models"
)

// DocumentUseCaseI defines business logic methods for documents.
// อินเทอร์เฟซธุรกิจสำหรับเอกสาร
type DocumentUseCaseI interface {
	internal.UseCaseI[models.Document]
	GetByFilename(ctx context.Context, filename string) (*models.Document, error)
}
