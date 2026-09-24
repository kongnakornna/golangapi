package usecase

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"icmongolang/internal/modules/queue"
	"icmongolang/internal/modules/websocket/models"
	"icmongolang/pkg/websocket"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository
type mockWSRepo struct {
	mock.Mock
}

func (m *mockWSRepo) SaveMessage(ctx context.Context, msg *models.WSMessage) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}
func (m *mockWSRepo) GetMessagesByTopic(ctx context.Context, topic string, limit int) ([]*models.WSMessage, error) {
	args := m.Called(ctx, topic, limit)
	return args.Get(0).([]*models.WSMessage), args.Error(1)
}
func (m *mockWSRepo) ValidateSession(ctx context.Context, token string) (*models.Session, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*models.Session), args.Error(1)
}

// mock queue
type mockQueue struct{ mock.Mock }

func (m *mockQueue) Publish(ctx context.Context, topic string, payload interface{}) error {
	args := m.Called(ctx, topic, payload)
	return args.Error(0)
}
func (m *mockQueue) Subscribe(ctx context.Context, topic string, handler queue.Handler) error {
	args := m.Called(ctx, topic, handler)
	return args.Error(0)
}
func (m *mockQueue) PublishDelayed(ctx context.Context, topic string, payload interface{}, delay time.Duration) error {
	args := m.Called(ctx, topic, payload, delay)
	return args.Error(0)
}
func (m *mockQueue) Close() error { return nil }

func TestWSUsecase_SaveMessage(t *testing.T) {
	repo := new(mockWSRepo)
	q := new(mockQueue)
	hub := &websocket.Hub{}
	uc := NewWSUsecase(repo, q, hub)

	ctx := context.Background()
	payload := []byte(`{"test":"ok"}`)
	repo.On("SaveMessage", ctx, mock.AnythingOfType("*models.WSMessage")).Return(nil).Run(func(args mock.Arguments) {
		msg := args.Get(1).(*models.WSMessage)
		assert.Equal(t, "test", msg.Topic)
		assert.Equal(t, json.RawMessage(payload), msg.Payload)
		assert.Equal(t, "user1", msg.SenderID)
	})

	id, err := uc.SaveMessage(ctx, "test", payload, "user1")
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
	repo.AssertExpectations(t)
}

func TestWSUsecase_HandleIncomingMessage(t *testing.T) {
	repo := new(mockWSRepo)
	q := new(mockQueue)
	hub := &websocket.Hub{}
	uc := NewWSUsecase(repo, q, hub)

	ctx := context.Background()
	payload := []byte(`{"msg":"hello"}`)
	repo.On("SaveMessage", ctx, mock.AnythingOfType("*models.WSMessage")).Return(nil)
	q.On("Publish", ctx, "topic1", mock.Anything).Return(nil)

	err := uc.HandleIncomingMessage(ctx, "topic1", "", payload, "user1")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
	q.AssertExpectations(t)
}
