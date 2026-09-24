package presenter

import "icmongolang/pkg/helpers"

// AuditEntry is the canonical audit log record.
type AuditEntry struct {
	ID         string            `json:"id"`
	UserID     string            `json:"user_id,omitempty"`
	Action     string            `json:"action"` // e.g. settings.update, device.control
	Resource   string            `json:"resource,omitempty"`
	ResourceID string            `json:"resource_id,omitempty"`
	RequestID  string            `json:"request_id,omitempty"`
	Before     string            `json:"before,omitempty"`
	After      string            `json:"after,omitempty"`
	Hash       string            `json:"hash"`
	PrevHash   string            `json:"prev_hash,omitempty"`
	CreatedAt  helpers.LocalTime `json:"created_at"`
	Anchored   bool              `json:"anchored"`
	TxHash     string            `json:"tx_hash,omitempty"`
}

// CreateRequest is the input to record a new audit entry.
type CreateRequest struct {
	Action     string `json:"action"`
	Resource   string `json:"resource,omitempty"`
	ResourceID string `json:"resource_id,omitempty"`
	Before     string `json:"before,omitempty"`
	After      string `json:"after,omitempty"`
}

// VerifyResponse is the result of checking a hash-chain entry.
type VerifyResponse struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason,omitempty"`
}

// ChainResponse returns the current chain tail (latest entry).
type ChainResponse struct {
	Length   int    `json:"length"`
	HeadID   string `json:"head_id"`
	HeadHash string `json:"head_hash"`
}
