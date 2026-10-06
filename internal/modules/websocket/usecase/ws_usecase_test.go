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
func (m *mockWSRepo) GetMessages(ctx context.Context, params models.MessageQueryParams) (*models.MessagePage, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*models.MessagePage), args.Error(1)
}
func (m *mockWSRepo) GetMessagesByRoom(ctx context.Context, room string, params models.MessageQueryParams) (*models.MessagePage, error) {
	args := m.Called(ctx, room, params)
	return args.Get(0).(*models.MessagePage), args.Error(1)
}
func (m *mockWSRepo) GetMessagesBySender(ctx context.Context, senderID string, params models.MessageQueryParams) (*models.MessagePage, error) {
	args := m.Called(ctx, senderID, params)
	return args.Get(0).(*models.MessagePage), args.Error(1)
}
func (m *mockWSRepo) CountMessages(ctx context.Context, params models.MessageQueryParams) (int64, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockWSRepo) SearchMessages(ctx context.Context, params models.SearchQueryParams) (*models.SearchPage, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*models.SearchPage), args.Error(1)
}
func (m *mockWSRepo) CountSearchMessages(ctx context.Context, params models.SearchQueryParams) (int64, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockWSRepo) CreateRetentionPolicy(ctx context.Context, policy *models.RetentionPolicy) error {
	args := m.Called(ctx, policy)
	return args.Error(0)
}
func (m *mockWSRepo) GetRetentionPolicy(ctx context.Context, id string) (*models.RetentionPolicy, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.RetentionPolicy), args.Error(1)
}
func (m *mockWSRepo) ListRetentionPolicies(ctx context.Context, limit, offset int) ([]*models.RetentionPolicy, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*models.RetentionPolicy), args.Get(1).(int64), args.Error(2)
}
func (m *mockWSRepo) UpdateRetentionPolicy(ctx context.Context, policy *models.RetentionPolicy) error {
	args := m.Called(ctx, policy)
	return args.Error(0)
}
func (m *mockWSRepo) DeleteRetentionPolicy(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockWSRepo) GetActiveRetentionPolicies(ctx context.Context) ([]*models.RetentionPolicy, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.RetentionPolicy), args.Error(1)
}
func (m *mockWSRepo) GetRetentionStats(ctx context.Context) (*models.RetentionStats, error) {
	args := m.Called(ctx)
	return args.Get(0).(*models.RetentionStats), args.Error(1)
}
func (m *mockWSRepo) CleanupExpiredMessages(ctx context.Context, policy *models.RetentionPolicy) (int64, error) {
	args := m.Called(ctx, policy)
	return args.Get(0).(int64), args.Error(1)
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
