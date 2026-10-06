// Package presenter จัดเตรียม HTTP request/response DTOs ของ keys module
package presenter

// ============================================================
// Request bodies
// ============================================================

// CreateAPIKeyRequest is the HTTP body for creating an API key.
type CreateAPIKeyRequest struct {
	Name        string  `json:"name" example:"mobile-app" validate:"required,min=3,max=100"`
	Description *string `json:"description,omitempty" example:"API key for mobile application" validate:"omitempty,max=500"`
	ExpiresIn   *int    `json:"expires_in,omitempty" example:"2592000" validate:"omitempty,min=1"`
	// ExpiresIn: อายุ key เป็น "นาที" — 0 = ไม่หมดอายุ
}

// RotateAPIKeyRequest is the HTTP body for rotating an API key.
type RotateAPIKeyRequest struct {
	// Empty body - rotation generates new secret automatically
}

// ============================================================
// Responses
// ============================================================

// APIKeyResponse is the HTTP response for an API key (without secret).
type APIKeyResponse struct {
	ID          string  `json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	Name        string  `json:"name" example:"mobile-app"`
	Description *string `json:"description,omitempty" example:"API key for mobile application"`
	Prefix      string  `json:"prefix" example:"icmon_"`
	LastFour    string  `json:"last_four" example:"a1b2"`
	ExpiresAt   *string `json:"expires_at,omitempty" example:"2026-12-31T23:59:59Z"`
	LastUsedAt  *string `json:"last_used_at,omitempty" example:"2026-10-04T12:34:56Z"`
	IsActive    bool    `json:"is_active" example:"true"`
	CreatedAt   string  `json:"created_at" example:"2026-10-04T00:00:00Z"`
	UpdatedAt   string  `json:"updated_at" example:"2026-10-04T00:00:00Z"`
}

// APIKeyWithSecretResponse is returned ONLY on create/rotate (raw secret once).
type APIKeyWithSecretResponse struct {
	APIKeyResponse
	// APIKey = prefix.base64secret — returned only once, not stored
	APIKey string `json:"api_key" example:"icmon_abc123def456..."`
}

// ListAPIKeysResponse is the paginated list response.
type ListAPIKeysResponse struct {
	Items  []APIKeyResponse `json:"items"`
	Total  int64            `json:"total" example:"42"`
	Limit  int              `json:"limit" example:"20"`
	Offset int              `json:"offset" example:"0"`
}
