package document

import (
	"context"

	"icmongolang/internal"
	"icmongolang/internal/models"
)

// DocumentPgRepository defines data access methods for documents.
// ดึงข้อมูลเอกสารจากฐานข้อมูล
type DocumentPgRepository interface {
	internal.PgRepository[models.Document]
	GetByFilename(ctx context.Context, filename string) (*models.Document, error)
}
