package presenter

import (
	"github.com/google/uuid"
)

// DocumentUploadRequest represents a document upload payload.
// คำขออัปโหลดเอกสาร
type DocumentUploadRequest struct {
	Filename string `json:"filename" validate:"required"`
	MimeType string `json:"mimeType" validate:"required"`
	Size     int64  `json:"size" validate:"required"`
}

// DocumentResponse represents a document record returned to the client.
// ข้อมูลเอกสารที่ส่งกลับไปยังผู้ใช้
type DocumentResponse struct {
	ID           uuid.UUID `json:"id"`
	Filename     string    `json:"filename"`
	OriginalName string    `json:"originalName"`
	MimeType     string    `json:"mimeType"`
	Size         int64     `json:"size"`
	CreatedAt    string    `json:"createdAt"`
}
