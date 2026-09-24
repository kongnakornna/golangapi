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
	HandleIncomingMessage(ctx context.Context, topic, room string, payload []byte, senderID string) error
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
