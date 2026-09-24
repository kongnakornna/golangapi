// internal/modules/websocket/repository/postgres/ws_repo_pg.go
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"icmongolang/internal/modules/websocket/models"
	"icmongolang/internal/modules/websocket/repository"
)

type pgRepo struct {
	db *sql.DB
}

func NewPGRepository(db *sql.DB) repository.WSRepository {
	return &pgRepo{db: db}
}

func (r *pgRepo) SaveMessage(ctx context.Context, msg *models.WSMessage) error {
	query := `INSERT INTO ws_messages (id, topic, payload, sender_id, sent_at)
              VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, msg.ID, msg.Topic, msg.Payload, msg.SenderID, msg.SentAt)
	return err
}

func (r *pgRepo) GetMessagesByTopic(ctx context.Context, topic string, limit int) ([]*models.WSMessage, error) {
	query := `SELECT id, topic, payload, sender_id, sent_at
              FROM ws_messages WHERE topic = $1 ORDER BY sent_at DESC LIMIT $2`
	rows, err := r.db.QueryContext(ctx, query, topic, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []*models.WSMessage
	for rows.Next() {
		var m models.WSMessage
		var payloadJSON []byte
		if err := rows.Scan(&m.ID, &m.Topic, &payloadJSON, &m.SenderID, &m.SentAt); err != nil {
			return nil, err
		}
		m.Payload = json.RawMessage(payloadJSON)
		msgs = append(msgs, &m)
	}
	return msgs, rows.Err()
}

func (r *pgRepo) ValidateSession(ctx context.Context, token string) (*models.Session, error) {
	query := `SELECT id, user_id, token, created_at, expires_at
              FROM ws_sessions WHERE token = $1 AND expires_at > NOW()`
	var s models.Session
	err := r.db.QueryRowContext(ctx, query, token).Scan(&s.ID, &s.UserID, &s.Token, &s.CreatedAt, &s.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
