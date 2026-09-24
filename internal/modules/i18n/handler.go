package i18n

import "net/http"

// Handlers defines HTTP handler methods for translations.
// ตัวจัดการ HTTP สำหรับการแปลภาษา
type Handlers interface {
	GetTranslations() func(w http.ResponseWriter, r *http.Request)
	GetTranslationByKey() func(w http.ResponseWriter, r *http.Request)
	CreateTranslation() func(w http.ResponseWriter, r *http.Request)
	UpdateTranslation() func(w http.ResponseWriter, r *http.Request)
	DeleteTranslation() func(w http.ResponseWriter, r *http.Request)
}
