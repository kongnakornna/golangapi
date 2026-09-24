package presenter

// APIKeyRequest is used to create or revoke an API key.
type APIKeyRequest struct {
	Name   string `json:"name"`
	Scopes string `json:"scopes,omitempty"`
	Revoke bool   `json:"revoke,omitempty"`
}

// APIKeyResponse is the result of key creation.
type APIKeyResponse struct {
	Key       string `json:"key,omitempty"`
	KeyID     string `json:"key_id"`
	HashedKey string `json:"hashed_key"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at,omitempty"`
	Revoked   bool   `json:"revoked"`
}

// RateLimitResponse reports the current rate-limit window state.
type RateLimitResponse struct {
	Allowed   bool  `json:"allowed"`
	Limit     int64 `json:"limit"`
	Remaining int64 `json:"remaining"`
	ResetIn   int64 `json:"reset_in_seconds"`
}

// UsageEntry is a single API usage record.
type UsageEntry struct {
	KeyID     string `json:"key_id"`
	Route     string `json:"route"`
	Status    int    `json:"status"`
	LatencyMs int64  `json:"latency_ms"`
	At        string `json:"at"`
}
