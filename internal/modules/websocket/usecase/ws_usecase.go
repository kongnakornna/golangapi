package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"icmongolang/internal/modules/queue"
	"icmongolang/internal/modules/websocket/models"
	"icmongolang/internal/modules/websocket/repository"
	"icmongolang/pkg/websocket"

	"github.com/google/uuid"
)

type WSUsecase interface {
	Authenticate(ctx context.Context, token string) (string, error)
	SaveMessage(ctx context.Context, topic string, payload []byte, senderID string) (string, error)
	GetTopicHistory(ctx context.Context, topic string, limit int) ([]*models.WSMessage, error)
	// New pagination methods
	GetMessages(ctx context.Context, params models.MessageQueryParams) (*models.MessagePage, error)
	GetMessagesByRoom(ctx context.Context, room string, params models.MessageQueryParams) (*models.MessagePage, error)
	GetMessagesBySender(ctx context.Context, senderID string, params models.MessageQueryParams) (*models.MessagePage, error)
	CountMessages(ctx context.Context, params models.MessageQueryParams) (int64, error)
	// Search methods
	SearchMessages(ctx context.Context, params models.SearchQueryParams) (*models.SearchPage, error)
	CountSearchMessages(ctx context.Context, params models.SearchQueryParams) (int64, error)
	HandleIncomingMessage(ctx context.Context, topic, room string, payload []byte, senderID string) error

	// Retention Policy methods
	CreateRetentionPolicy(ctx context.Context, req *models.RetentionPolicyRequest) (*models.RetentionPolicyResponse, error)
	GetRetentionPolicy(ctx context.Context, id string) (*models.RetentionPolicyResponse, error)
	ListRetentionPolicies(ctx context.Context, limit, offset int) (*models.RetentionPolicyListResponse, error)
	UpdateRetentionPolicy(ctx context.Context, id string, req *models.RetentionPolicyRequest) (*models.RetentionPolicyResponse, error)
	DeleteRetentionPolicy(ctx context.Context, id string) error
	GetRetentionStats(ctx context.Context) (*models.RetentionStats, error)
	RunRetentionCleanup(ctx context.Context) (*models.RetentionStats, error)
}

type wsUsecase struct {
	repo  repository.WSRepository
	queue queue.Queue
	hub   *websocket.Hub // สำหรับ broadcast
}

func NewWSUsecase(repo repository.WSRepository, q queue.Queue, hub *websocket.Hub) WSUsecase {
	return &wsUsecase{
		repo:  repo,
		queue: q,
		hub:   hub,
	}
}

func (u *wsUsecase) Authenticate(ctx context.Context, token string) (string, error) {
	session, err := u.repo.ValidateSession(ctx, token)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}
	return session.UserID, nil
}

func (u *wsUsecase) SaveMessage(ctx context.Context, topic string, payload []byte, senderID string) (string, error) {
	msg := &models.WSMessage{
		ID:       uuid.New().String(),
		Topic:    topic,
		Payload:  json.RawMessage(payload),
		SenderID: senderID,
		SentAt:   time.Now(),
	}
	if err := u.repo.SaveMessage(ctx, msg); err != nil {
		return "", err
	}
	return msg.ID, nil
}

func (u *wsUsecase) GetTopicHistory(ctx context.Context, topic string, limit int) ([]*models.WSMessage, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return u.repo.GetMessagesByTopic(ctx, topic, limit)
}

// GetMessages retrieves messages with full pagination support.
func (u *wsUsecase) GetMessages(ctx context.Context, params models.MessageQueryParams) (*models.MessagePage, error) {
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
	return u.repo.GetMessages(ctx, params)
}

// GetMessagesByRoom retrieves messages for a room with pagination.
func (u *wsUsecase) GetMessagesByRoom(ctx context.Context, room string, params models.MessageQueryParams) (*models.MessagePage, error) {
	params.Room = room
	return u.GetMessages(ctx, params)
}

// GetMessagesBySender retrieves messages sent by a specific user with pagination.
func (u *wsUsecase) GetMessagesBySender(ctx context.Context, senderID string, params models.MessageQueryParams) (*models.MessagePage, error) {
	params.SenderID = senderID
	return u.GetMessages(ctx, params)
}

// CountMessages returns the total count of messages matching the query.
func (u *wsUsecase) CountMessages(ctx context.Context, params models.MessageQueryParams) (int64, error) {
	return u.repo.CountMessages(ctx, params)
}

// SearchMessages performs full-text search on messages.
func (u *wsUsecase) SearchMessages(ctx context.Context, params models.SearchQueryParams) (*models.SearchPage, error) {
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
	return u.repo.SearchMessages(ctx, params)
}

// CountSearchMessages returns the total count of messages matching the search query.
func (u *wsUsecase) CountSearchMessages(ctx context.Context, params models.SearchQueryParams) (int64, error) {
	return u.repo.CountSearchMessages(ctx, params)
}

func (u *wsUsecase) HandleIncomingMessage(ctx context.Context, topic, room string, payload []byte, senderID string) error {
	// บันทึก history
	msgID, err := u.SaveMessage(ctx, topic, payload, senderID)
	if err != nil {
		return err
	}
	_ = msgID

	// Publish ไปยัง queue (topic)
	if topic != "" {
		if err := u.queue.Publish(ctx, topic, json.RawMessage(payload)); err != nil {
			return err
		}
	}

	// ถ้ามี room ก็ broadcast ไปห้องนั้น
	if room != "" {
		u.hub.BroadcastToRoom(room, "message", json.RawMessage(payload))
	}

	return nil
}

// CreateRetentionPolicy creates a new retention policy.
func (u *wsUsecase) CreateRetentionPolicy(ctx context.Context, req *models.RetentionPolicyRequest) (*models.RetentionPolicyResponse, error) {
	now := time.Now()
	policy := &models.RetentionPolicy{
		ID:        uuid.New().String(),
		Topic:     req.Topic,
		Room:      req.Room,
		MaxAge:    req.MaxAge,
		MaxCount:  req.MaxCount,
		Enabled:   req.Enabled,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := u.repo.CreateRetentionPolicy(ctx, policy); err != nil {
		return nil, err
	}
	return u.policyToResponse(policy), nil
}

// GetRetentionPolicy retrieves a retention policy by ID.
func (u *wsUsecase) GetRetentionPolicy(ctx context.Context, id string) (*models.RetentionPolicyResponse, error) {
	policy, err := u.repo.GetRetentionPolicy(ctx, id)
	if err != nil {
		return nil, err
	}
	return u.policyToResponse(policy), nil
}

// ListRetentionPolicies lists all retention policies with pagination.
func (u *wsUsecase) ListRetentionPolicies(ctx context.Context, limit, offset int) (*models.RetentionPolicyListResponse, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 1000 {
		limit = 1000
	}
	if offset < 0 {
		offset = 0
	}

	policies, total, err := u.repo.ListRetentionPolicies(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	response := &models.RetentionPolicyListResponse{
		Policies: make([]models.RetentionPolicyResponse, len(policies)),
		Total:    total,
	}
	for i, p := range policies {
		response.Policies[i] = *u.policyToResponse(p)
	}
	return response, nil
}

// UpdateRetentionPolicy updates an existing retention policy.
func (u *wsUsecase) UpdateRetentionPolicy(ctx context.Context, id string, req *models.RetentionPolicyRequest) (*models.RetentionPolicyResponse, error) {
	policy, err := u.repo.GetRetentionPolicy(ctx, id)
	if err != nil {
		return nil, err
	}

	policy.Topic = req.Topic
	policy.Room = req.Room
	policy.MaxAge = req.MaxAge
	policy.MaxCount = req.MaxCount
	policy.Enabled = req.Enabled
	policy.UpdatedAt = time.Now()

	if err := u.repo.UpdateRetentionPolicy(ctx, policy); err != nil {
		return nil, err
	}
	return u.policyToResponse(policy), nil
}

// DeleteRetentionPolicy deletes a retention policy by ID.
func (u *wsUsecase) DeleteRetentionPolicy(ctx context.Context, id string) error {
	return u.repo.DeleteRetentionPolicy(ctx, id)
}

// GetRetentionStats retrieves statistics about message retention.
func (u *wsUsecase) GetRetentionStats(ctx context.Context) (*models.RetentionStats, error) {
	return u.repo.GetRetentionStats(ctx)
}

// RunRetentionCleanup runs the retention cleanup for all active policies.
func (u *wsUsecase) RunRetentionCleanup(ctx context.Context) (*models.RetentionStats, error) {
	policies, err := u.repo.GetActiveRetentionPolicies(ctx)
	if err != nil {
		return nil, err
	}

	var totalDeleted int64
	for _, policy := range policies {
		deleted, err := u.repo.CleanupExpiredMessages(ctx, policy)
		if err != nil {
			return nil, err
		}
		totalDeleted += deleted
	}

	stats, err := u.repo.GetRetentionStats(ctx)
	if err != nil {
		return nil, err
	}
	stats.MessagesDeleted = totalDeleted
	stats.LastCleanupAt = time.Now()

	return stats, nil
}

// policyToResponse converts a RetentionPolicy to RetentionPolicyResponse.
func (u *wsUsecase) policyToResponse(p *models.RetentionPolicy) *models.RetentionPolicyResponse {
	return &models.RetentionPolicyResponse{
		ID:        p.ID,
		Topic:     p.Topic,
		Room:      p.Room,
		MaxAge:    p.MaxAge,
		MaxCount:  p.MaxCount,
		Enabled:   p.Enabled,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
