package presenter

// TranslationRequest represents a translation creation/update payload.
// คำขอสร้าง/แก้ไขคำแปล
type TranslationRequest struct {
	Locale string `json:"locale" validate:"required,min=2,max=10"`
	Key    string `json:"key" validate:"required,min=1,max=255"`
	Value  string `json:"value" validate:"required"`
}

// TranslationResponse represents a translation record returned to the client.
// ข้อมูลคำแปลที่ส่งกลับไปยังผู้ใช้
type TranslationResponse struct {
	ID        int    `json:"id"`
	Locale    string `json:"locale"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
