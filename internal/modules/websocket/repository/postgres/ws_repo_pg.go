// internal/modules/websocket/repository/postgres/ws_repo_pg.go
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

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
	query := `INSERT INTO ws_messages (id, topic, room, payload, sender_id, sender_role, sent_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, msg.ID, msg.Topic, msg.Room, msg.Payload, msg.SenderID, msg.SenderRole, msg.SentAt)
	return err
}

func (r *pgRepo) GetMessagesByTopic(ctx context.Context, topic string, limit int) ([]*models.WSMessage, error) {
	query := `SELECT id, topic, room, payload, sender_id, sender_role, sent_at
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
		if err := rows.Scan(&m.ID, &m.Topic, &m.Room, &payloadJSON, &m.SenderID, &m.SenderRole, &m.SentAt); err != nil {
			return nil, err
		}
		m.Payload = json.RawMessage(payloadJSON)
		msgs = append(msgs, &m)
	}
	return msgs, rows.Err()
}

func (r *pgRepo) GetMessages(ctx context.Context, params models.MessageQueryParams) (*models.MessagePage, error) {
	return r.queryMessages(ctx, params)
}

func (r *pgRepo) GetMessagesByRoom(ctx context.Context, room string, params models.MessageQueryParams) (*models.MessagePage, error) {
	params.Room = room
	return r.queryMessages(ctx, params)
}

func (r *pgRepo) GetMessagesBySender(ctx context.Context, senderID string, params models.MessageQueryParams) (*models.MessagePage, error) {
	params.SenderID = senderID
	return r.queryMessages(ctx, params)
}

func (r *pgRepo) CountMessages(ctx context.Context, params models.MessageQueryParams) (int64, error) {
	whereClause, args := r.buildWhereClause(params)
	query := fmt.Sprintf("SELECT COUNT(*) FROM ws_messages %s", whereClause)
	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
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

// queryMessages is the internal method that handles all message queries with pagination.
func (r *pgRepo) queryMessages(ctx context.Context, params models.MessageQueryParams) (*models.MessagePage, error) {
	// Set defaults
	if params.Limit <= 0 {
		params.Limit = 50
	}
	if params.Limit > 1000 {
		params.Limit = 1000
	}
	if params.Order == "" {
		params.Order = "desc"
	}
	if params.Order != "asc" && params.Order != "desc" {
		params.Order = "desc"
	}

	whereClause, args := r.buildWhereClause(params)

	// Build ORDER BY clause
	orderBy := "sent_at " + strings.ToUpper(params.Order)

	// Build LIMIT/OFFSET
	limit := params.Limit
	offset := params.Offset

	// For cursor-based pagination, we use the Before/After timestamps
	if params.Before != nil {
		if params.Order == "desc" {
			// For desc order, "before" means older messages (sent_at < before)
			whereClause += " AND sent_at < $" + fmt.Sprintf("%d", len(args)+1)
			args = append(args, *params.Before)
		} else {
			// For asc order, "before" means newer messages (sent_at > before)
			whereClause += " AND sent_at > $" + fmt.Sprintf("%d", len(args)+1)
			args = append(args, *params.Before)
		}
	}
	if params.After != nil {
		if params.Order == "desc" {
			// For desc order, "after" means newer messages (sent_at > after)
			whereClause += " AND sent_at > $" + fmt.Sprintf("%d", len(args)+1)
			args = append(args, *params.After)
		} else {
			// For asc order, "after" means older messages (sent_at < after)
			whereClause += " AND sent_at < $" + fmt.Sprintf("%d", len(args)+1)
			args = append(args, *params.After)
		}
	}

	// Get total count
	totalCount, err := r.CountMessages(context.Background(), models.MessageQueryParams{
		Topic:    params.Topic,
		Room:     params.Room,
		SenderID: params.SenderID,
		Before:   params.Before,
		After:    params.After,
	})
	if err != nil {
		return nil, err
	}

	// Fetch messages with pagination
	query := fmt.Sprintf(`
		SELECT id, topic, room, payload, sender_id, sender_role, sent_at
		FROM ws_messages %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, len(args)+1, len(args)+2)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []*models.WSMessage
	for rows.Next() {
		var m models.WSMessage
		if err := rows.Scan(&m.ID, &m.Topic, &m.Room, &m.Payload, &m.SenderID, &m.SenderRole, &m.SentAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, &m)
	}

	// Determine cursors for next/prev pages
	var nextCursor, prevCursor *time.Time
	if len(msgs) > 0 {
		lastMsg := msgs[len(msgs)-1]
		firstMsg := msgs[0]

		if params.Order == "desc" {
			// For desc order, next page is older messages (before last message)
			nextCursor = &lastMsg.SentAt
			// Previous page is newer messages (after first message)
			prevCursor = &firstMsg.SentAt
		} else {
			// For asc order, next page is newer messages (after last message)
			nextCursor = &lastMsg.SentAt
			// Previous page is older messages (before first message)
			prevCursor = &firstMsg.SentAt
		}
	}

	hasMore := int64(len(msgs)) == int64(params.Limit)

	return &models.MessagePage{
		Messages:   msgs,
		TotalCount: totalCount,
		HasMore:    hasMore,
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
	}, nil
}

// buildWhereClause constructs the WHERE clause and arguments for message queries.
func (r *pgRepo) buildWhereClause(params models.MessageQueryParams) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.Topic != "" {
		conditions = append(conditions, fmt.Sprintf("topic = $%d", argIndex))
		args = append(args, params.Topic)
		argIndex++
	}
	if params.Room != "" {
		conditions = append(conditions, fmt.Sprintf("room = $%d", argIndex))
		args = append(args, params.Room)
		argIndex++
	}
	if params.SenderID != "" {
		conditions = append(conditions, fmt.Sprintf("sender_id = $%d", argIndex))
		args = append(args, params.SenderID)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}
	return whereClause, args
}

// buildSearchWhereClause constructs the WHERE clause for full-text search.
func (r *pgRepo) buildSearchWhereClause(params models.SearchQueryParams) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Full-text search using PostgreSQL's tsvector/tsquery
	if params.Query != "" {
		conditions = append(conditions, fmt.Sprintf("to_tsvector('english', payload::text) @@ plainto_tsquery('english', $%d)", argIndex))
		args = append(args, params.Query)
		argIndex++
	}

	if params.Topic != "" {
		conditions = append(conditions, fmt.Sprintf("topic = $%d", argIndex))
		args = append(args, params.Topic)
		argIndex++
	}
	if params.Room != "" {
		conditions = append(conditions, fmt.Sprintf("room = $%d", argIndex))
		args = append(args, params.Room)
		argIndex++
	}
	if params.SenderID != "" {
		conditions = append(conditions, fmt.Sprintf("sender_id = $%d", argIndex))
		args = append(args, params.SenderID)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}
	return whereClause, args
}

// SearchMessages performs full-text search on messages.
func (r *pgRepo) SearchMessages(ctx context.Context, params models.SearchQueryParams) (*models.SearchPage, error) {
	// Set defaults
	if params.Limit <= 0 {
		params.Limit = 50
	}
	if params.Limit > 1000 {
		params.Limit = 1000
	}
	if params.Order == "" {
		params.Order = "desc"
	}
	if params.Order != "asc" && params.Order != "desc" {
		params.Order = "desc"
	}

	whereClause, args := r.buildSearchWhereClause(params)

	// Build ORDER BY clause - order by rank for search results
	orderBy := "ts_rank_cd(to_tsvector('english', payload::text), plainto_tsquery('english', $1)) " + strings.ToUpper(params.Order)
	// Add the query as first argument for ranking
	rankArgs := []interface{}{params.Query}
	rankArgs = append(rankArgs, args...)

	// Build LIMIT/OFFSET
	limit := params.Limit
	offset := params.Offset

	// For cursor-based pagination
	if params.Before != nil {
		if params.Order == "desc" {
			whereClause += " AND sent_at < $" + fmt.Sprintf("%d", len(args)+1)
			args = append(args, *params.Before)
		} else {
			whereClause += " AND sent_at > $" + fmt.Sprintf("%d", len(args)+1)
			args = append(args, *params.Before)
		}
	}
	if params.After != nil {
		if params.Order == "desc" {
			whereClause += " AND sent_at > $" + fmt.Sprintf("%d", len(args)+1)
			args = append(args, *params.After)
		} else {
			whereClause += " AND sent_at < $" + fmt.Sprintf("%d", len(args)+1)
			args = append(args, *params.After)
		}
	}

	// Get total count
	totalCount, err := r.CountSearchMessages(ctx, params)
	if err != nil {
		return nil, err
	}

	// Build SELECT clause with optional highlights
	selectClause := `
		SELECT id, topic, room, payload, sender_id, sender_role, sent_at,
		       ts_rank_cd(to_tsvector('english', payload::text), plainto_tsquery('english', $1)) as rank
	`
	if params.Highlight {
		selectClause = `
			SELECT id, topic, room, payload, sender_id, sender_role, sent_at,
			       ts_rank_cd(to_tsvector('english', payload::text), plainto_tsquery('english', $1)) as rank,
			       ts_headline('english', payload::text, plainto_tsquery('english', $1), 'StartSel=<mark>, StopSel=</mark>') as highlight
		`
	}

	// Fetch messages with pagination
	query := fmt.Sprintf(`
		%s
		FROM ws_messages %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, selectClause, whereClause, orderBy, len(args)+1, len(args)+2)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.SearchResult
	for rows.Next() {
		var result models.SearchResult
		var payloadJSON []byte
		var rank *float64
		var highlight *string

		scanArgs := []interface{}{&result.Message.ID, &result.Message.Topic, &result.Message.Room, &payloadJSON, &result.Message.SenderID, &result.Message.SenderRole, &result.Message.SentAt, &rank}
		if params.Highlight {
			scanArgs = append(scanArgs, &highlight)
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		result.Message.Payload = json.RawMessage(payloadJSON)
		if rank != nil {
			result.Rank = *rank
		}
		if highlight != nil {
			result.Highlights = map[string][]string{"payload": {*highlight}}
		}

		results = append(results, &result)
	}

	// Determine cursors for next/prev pages
	var nextCursor, prevCursor *time.Time
	if len(results) > 0 {
		lastMsg := results[len(results)-1].Message
		firstMsg := results[0].Message

		if params.Order == "desc" {
			nextCursor = &lastMsg.SentAt
			prevCursor = &firstMsg.SentAt
		} else {
			nextCursor = &lastMsg.SentAt
			prevCursor = &firstMsg.SentAt
		}
	}

	hasMore := int64(len(results)) == int64(params.Limit)

	return &models.SearchPage{
		Results:    results,
		TotalCount: totalCount,
		HasMore:    hasMore,
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
	}, nil
}

// CountSearchMessages returns the total count of messages matching the search query.
func (r *pgRepo) CountSearchMessages(ctx context.Context, params models.SearchQueryParams) (int64, error) {
	whereClause, args := r.buildSearchWhereClause(params)
	query := fmt.Sprintf("SELECT COUNT(*) FROM ws_messages %s", whereClause)
	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

// CreateRetentionPolicy creates a new retention policy.
func (r *pgRepo) CreateRetentionPolicy(ctx context.Context, policy *models.RetentionPolicy) error {
	query := `INSERT INTO ws_retention_policies (id, topic, room, max_age, max_count, enabled, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, policy.ID, policy.Topic, policy.Room, policy.MaxAge, policy.MaxCount, policy.Enabled, policy.CreatedAt, policy.UpdatedAt)
	return err
}

// GetRetentionPolicy retrieves a retention policy by ID.
func (r *pgRepo) GetRetentionPolicy(ctx context.Context, id string) (*models.RetentionPolicy, error) {
	query := `SELECT id, topic, room, max_age, max_count, enabled, created_at, updated_at
              FROM ws_retention_policies WHERE id = $1`
	var p models.RetentionPolicy
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Topic, &p.Room, &p.MaxAge, &p.MaxCount, &p.Enabled, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListRetentionPolicies lists all retention policies with pagination.
func (r *pgRepo) ListRetentionPolicies(ctx context.Context, limit, offset int) ([]*models.RetentionPolicy, int64, error) {
	// Get total count
	var total int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM ws_retention_policies").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, topic, room, max_age, max_count, enabled, created_at, updated_at
              FROM ws_retention_policies ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var policies []*models.RetentionPolicy
	for rows.Next() {
		var p models.RetentionPolicy
		if err := rows.Scan(&p.ID, &p.Topic, &p.Room, &p.MaxAge, &p.MaxCount, &p.Enabled, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		policies = append(policies, &p)
	}
	return policies, total, rows.Err()
}

// UpdateRetentionPolicy updates an existing retention policy.
func (r *pgRepo) UpdateRetentionPolicy(ctx context.Context, policy *models.RetentionPolicy) error {
	query := `UPDATE ws_retention_policies SET topic = $2, room = $3, max_age = $4, max_count = $5, enabled = $6, updated_at = $7
              WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, policy.ID, policy.Topic, policy.Room, policy.MaxAge, policy.MaxCount, policy.Enabled, policy.UpdatedAt)
	return err
}

// DeleteRetentionPolicy deletes a retention policy by ID.
func (r *pgRepo) DeleteRetentionPolicy(ctx context.Context, id string) error {
	query := `DELETE FROM ws_retention_policies WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetActiveRetentionPolicies retrieves all active retention policies.
func (r *pgRepo) GetActiveRetentionPolicies(ctx context.Context) ([]*models.RetentionPolicy, error) {
	query := `SELECT id, topic, room, max_age, max_count, enabled, created_at, updated_at
              FROM ws_retention_policies WHERE enabled = true`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var policies []*models.RetentionPolicy
	for rows.Next() {
		var p models.RetentionPolicy
		if err := rows.Scan(&p.ID, &p.Topic, &p.Room, &p.MaxAge, &p.MaxCount, &p.Enabled, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		policies = append(policies, &p)
	}
	return policies, rows.Err()
}

// GetRetentionStats retrieves statistics about message retention.
func (r *pgRepo) GetRetentionStats(ctx context.Context) (*models.RetentionStats, error) {
	var stats models.RetentionStats

	// Total messages
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM ws_messages").Scan(&stats.TotalMessages)
	if err != nil {
		return nil, err
	}

	// Active policies count
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM ws_retention_policies WHERE enabled = true").Scan(&stats.ActivePolicies)
	if err != nil {
		return nil, err
	}

	// Total policies count
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM ws_retention_policies").Scan(&stats.PoliciesCount)
	if err != nil {
		return nil, err
	}

	// Last cleanup time (from a metadata table or we can track it)
	// For now, we'll use a simple approach - check the oldest message
	var lastCleanup sql.NullTime
	err = r.db.QueryRowContext(ctx, "SELECT MIN(sent_at) FROM ws_messages").Scan(&lastCleanup)
	if err == nil && lastCleanup.Valid {
		stats.LastCleanupAt = lastCleanup.Time
	}

	// Next cleanup - estimate based on oldest message + min retention
	var oldestMsg sql.NullTime
	err = r.db.QueryRowContext(ctx, "SELECT MIN(sent_at) FROM ws_messages").Scan(&oldestMsg)
	if err == nil && oldestMsg.Valid {
		// Find minimum max_age among active policies
		var minAge sql.NullString
		err = r.db.QueryRowContext(ctx, "SELECT MIN(max_age) FROM ws_retention_policies WHERE enabled = true").Scan(&minAge)
		if err == nil && minAge.Valid {
			// Parse duration and add to oldest message
			if d, err := time.ParseDuration(minAge.String); err == nil {
				stats.NextCleanupAt = oldestMsg.Time.Add(d)
			}
		}
	}

	return &stats, nil
}

// CleanupExpiredMessages deletes messages older than the retention policy allows.
func (r *pgRepo) CleanupExpiredMessages(ctx context.Context, policy *models.RetentionPolicy) (int64, error) {
	// Build WHERE clause for topic/room matching
	var conditions []string
	args := []interface{}{policy.MaxAge}
	argIndex := 2

	if policy.Topic != "" {
		conditions = append(conditions, fmt.Sprintf("topic LIKE $%d", argIndex))
		args = append(args, policy.Topic)
		argIndex++
	}
	if policy.Room != "" {
		conditions = append(conditions, fmt.Sprintf("room LIKE $%d", argIndex))
		args = append(args, policy.Room)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Delete by age
	ageQuery := fmt.Sprintf(`DELETE FROM ws_messages %s AND sent_at < NOW() - $1::interval`, whereClause)
	result, err := r.db.ExecContext(ctx, ageQuery, args...)
	if err != nil {
		return 0, err
	}
	deletedByAge, _ := result.RowsAffected()

	// Delete by count (keep only the most recent N messages per topic/room)
	if policy.MaxCount > 0 {
		// This is more complex - we need to delete old messages beyond max_count per topic/room
		// Using a CTE with ROW_NUMBER()
		countQuery := fmt.Sprintf(`
			DELETE FROM ws_messages
			WHERE id IN (
				SELECT id FROM (
					SELECT id,
					       ROW_NUMBER() OVER (PARTITION BY topic, room ORDER BY sent_at DESC) as rn
					FROM ws_messages
					%s
				) t
				WHERE t.rn > $%d
			)
		`, whereClause, argIndex)
		args = append(args, policy.MaxCount)
		result, err := r.db.ExecContext(ctx, countQuery, args...)
		if err != nil {
			return deletedByAge, err
		}
		deletedByCount, _ := result.RowsAffected()
		return deletedByAge + deletedByCount, nil
	}

	return deletedByAge, nil
}
