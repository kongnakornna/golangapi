package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type MessageHandler func(msg *OrderMessage) error

type Consumer struct {
	client  sarama.ConsumerGroup
	topic   string
	groupID string
	handler MessageHandler
	wg      sync.WaitGroup
}

func NewConsumer(brokers []string, groupID, topic string, handler MessageHandler, timeoutSec int) (*Consumer, error) {
	if timeoutSec <= 0 {
		timeoutSec = 30
	}

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true
	config.Net.DialTimeout = time.Duration(timeoutSec) * time.Second
	config.Net.ReadTimeout = time.Duration(timeoutSec) * time.Second
	config.Net.WriteTimeout = time.Duration(timeoutSec) * time.Second

	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &Consumer{
		client:  client,
		topic:   topic,
		groupID: groupID,
		handler: handler,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := c.client.Consume(ctx, []string{c.topic}, &consumerGroupHandler{handler: c.handler}); err != nil {
					log.Printf("Consumer error: %v", err)
				}
			}
		}
	}()
	return nil
}

func (c *Consumer) Wait() {
	c.wg.Wait()
}

func (c *Consumer) Close() error {
	return c.client.Close()
}

type consumerGroupHandler struct {
	handler MessageHandler
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var orderMsg OrderMessage
		if err := json.Unmarshal(msg.Value, &orderMsg); err != nil {
			log.Printf("Unmarshal error: %v", err)
			sess.MarkMessage(msg, "")
			continue
		}

		if err := h.handler(&orderMsg); err != nil {
			log.Printf("Handler error: %v", err)
			// ไม่ commit เพื่อให้ retry (อาจมี retry logic เพิ่มเติม)
			continue
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
