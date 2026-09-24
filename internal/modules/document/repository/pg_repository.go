package repository

import (
	"context"

	"icmongolang/internal/models"
	"icmongolang/internal/modules/document"
	"icmongolang/internal/repository"

	"gorm.io/gorm"
)

// DocumentPgRepo implements document.DocumentPgRepository.
// รีโพสิทอรีสำหรับเอกสาร
type DocumentPgRepo struct {
	repository.PgRepo[models.Document]
}

// CreateDocumentPgRepository creates a new document repository.
// สร้างรีโพสิทอรีสำหรับเอกสาร
func CreateDocumentPgRepository(db *gorm.DB) document.DocumentPgRepository {
	return &DocumentPgRepo{
		PgRepo: repository.CreatePgRepo[models.Document](db),
	}
}

// GetByFilename finds a document by its filename.
// ค้นหาเอกสารตามชื่อไฟล์
func (r *DocumentPgRepo) GetByFilename(ctx context.Context, filename string) (*models.Document, error) {
	var doc models.Document
	if result := r.DB.WithContext(ctx).Where("filename = ?", filename).First(&doc); result.Error != nil {
		return nil, result.Error
	}
	return &doc, nil
}
