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
	// ValidateSession checks if a token is valid and returns the session.
	ValidateSession(ctx context.Context, token string) (*models.Session, error)
}
