package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"icmongolang/internal/modules/kafka/models"
	"icmongolang/internal/modules/kafka/repository"
	"icmongolang/pkg/helpers"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/websocket"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, req *models.CreateOrderRequest) (*models.Order, error)
	ProcessOrderMessage(msg *kafka.OrderMessage) error
	HandleWebSocketMessage(client *websocket.Client, message []byte) error
}

type orderUsecase struct {
	repo     repository.OrderRepository
	producer *kafka.Producer
	topic    string
	hub      *websocket.Hub
}

func NewOrderUsecase(repo repository.OrderRepository, producer *kafka.Producer, topic string, hub *websocket.Hub) OrderUsecase {
	return &orderUsecase{
		repo:     repo,
		producer: producer,
		topic:    topic,
		hub:      hub,
	}
}

// CreateOrder – ใช้สำหรับ REST API
func (u *orderUsecase) CreateOrder(ctx context.Context, req *models.CreateOrderRequest) (*models.Order, error) {
	order := &models.Order{
		ID:        uuid.New(),
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Status:    "PENDING",
		CreatedAt: time.Now().In(helpers.GetTimeLocation()),
		UpdatedAt: time.Now().In(helpers.GetTimeLocation()),
	}

	msg := &kafka.OrderMessage{
		OrderID:   order.ID.String(),
		ProductID: order.ProductID,
		Quantity:  order.Quantity,
		CreatedAt: order.CreatedAt,
	}

	if err := u.producer.PublishMessage(order.ID.String(), msg); err != nil {
		return nil, fmt.Errorf("publish failed: %w", err)
	}

	if err := u.repo.Create(order); err != nil {
		return nil, fmt.Errorf("save order failed: %w", err)
	}

	return order, nil
}

// ProcessOrderMessage – เรียกโดย Consumer
func (u *orderUsecase) ProcessOrderMessage(msg *kafka.OrderMessage) error {
	order, err := u.repo.FindByID(msg.OrderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("order not found: %s", msg.OrderID)
		}
		return err
	}

	if order.Status != "PENDING" {
		return nil // already processed
	}

	order.Status = "PROCESSED"
	order.UpdatedAt = time.Now().In(helpers.GetTimeLocation())
	if err := u.repo.Update(order); err != nil {
		return err
	}

	u.BroadcastOrderStatus(order.ID.String(), order.Status)
	return nil
}

// HandleWebSocketMessage – รับข้อความจาก WebSocket แล้ว publish
func (u *orderUsecase) HandleWebSocketMessage(client *websocket.Client, message []byte) error {
	var req models.CreateOrderRequest
	if err := json.Unmarshal(message, &req); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	order := &models.Order{
		ID:        uuid.New(),
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Status:    "PENDING",
		CreatedAt: time.Now().In(helpers.GetTimeLocation()),
		UpdatedAt: time.Now().In(helpers.GetTimeLocation()),
	}

	msg := &kafka.OrderMessage{
		OrderID:   order.ID.String(),
		ProductID: order.ProductID,
		Quantity:  order.Quantity,
		CreatedAt: order.CreatedAt,
	}

	if err := u.producer.PublishMessage(order.ID.String(), msg); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	if err := u.repo.Create(order); err != nil {
		return fmt.Errorf("save failed: %w", err)
	}

	// ✅ ส่งตอบกลับ client โดยใช้เมธอด SendMessage
	response := map[string]interface{}{
		"event":    "order_created",
		"order_id": order.ID.String(),
		"status":   order.Status,
	}
	respData, _ := json.Marshal(response)
	client.SendMessage(respData)

	return nil
}

// BroadcastOrderStatus – ส่งสถานะไปยังทุก WebSocket client โดยใช้ BroadcastMessage
func (u *orderUsecase) BroadcastOrderStatus(orderID, status string) {
	u.hub.BroadcastMessage("order_updated", map[string]interface{}{
		"order_id":  orderID,
		"status":    status,
		"timestamp": time.Now().In(helpers.GetTimeLocation()),
	})
}
