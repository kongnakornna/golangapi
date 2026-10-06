package repository

import (
	"context"
	"icmongolang/internal/modules/websocket/models"
)

type WSRepository interface {
	// SaveMessage stores a message in the database.
	SaveMessage(ctx context.Context, msg *models.WSMessage) error
	// GetMessagesByTopic retrieves recent messages for a topic (pagination optional).
	GetMessagesByTopic(ctx context.Context, topic string, limit int) ([]*models.WSMessage, error)
	// GetMessages retrieves messages with full pagination support.
	GetMessages(ctx context.Context, params models.MessageQueryParams) (*models.MessagePage, error)
	// GetMessagesByRoom retrieves messages for a room with pagination.
	GetMessagesByRoom(ctx context.Context, room string, params models.MessageQueryParams) (*models.MessagePage, error)
	// GetMessagesBySender retrieves messages sent by a specific user with pagination.
	GetMessagesBySender(ctx context.Context, senderID string, params models.MessageQueryParams) (*models.MessagePage, error)
	// CountMessages returns the total count of messages matching the query.
	CountMessages(ctx context.Context, params models.MessageQueryParams) (int64, error)
	// SearchMessages performs full-text search on messages.
	SearchMessages(ctx context.Context, params models.SearchQueryParams) (*models.SearchPage, error)
	// CountSearchMessages returns the total count of messages matching the search query.
	CountSearchMessages(ctx context.Context, params models.SearchQueryParams) (int64, error)
	// ValidateSession checks if a token is valid and returns the session.
	ValidateSession(ctx context.Context, token string) (*models.Session, error)

	// Retention Policy methods
	// CreateRetentionPolicy creates a new retention policy.
	CreateRetentionPolicy(ctx context.Context, policy *models.RetentionPolicy) error
	// GetRetentionPolicy retrieves a retention policy by ID.
	GetRetentionPolicy(ctx context.Context, id string) (*models.RetentionPolicy, error)
	// ListRetentionPolicies lists all retention policies with pagination.
	ListRetentionPolicies(ctx context.Context, limit, offset int) ([]*models.RetentionPolicy, int64, error)
	// UpdateRetentionPolicy updates an existing retention policy.
	UpdateRetentionPolicy(ctx context.Context, policy *models.RetentionPolicy) error
	// DeleteRetentionPolicy deletes a retention policy by ID.
	DeleteRetentionPolicy(ctx context.Context, id string) error
	// GetActiveRetentionPolicies retrieves all active retention policies.
	GetActiveRetentionPolicies(ctx context.Context) ([]*models.RetentionPolicy, error)
	// GetRetentionStats retrieves statistics about message retention.
	GetRetentionStats(ctx context.Context) (*models.RetentionStats, error)
	// CleanupExpiredMessages deletes messages older than the retention policy allows.
	CleanupExpiredMessages(ctx context.Context, policy *models.RetentionPolicy) (int64, error)
}
