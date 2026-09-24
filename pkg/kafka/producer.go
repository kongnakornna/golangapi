package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	client sarama.SyncProducer
	topic  string
}

func NewProducer(brokers []string, topic string, timeoutSec int) (*Producer, error) {
	if timeoutSec <= 0 {
		timeoutSec = 30
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Net.DialTimeout = time.Duration(timeoutSec) * time.Second
	config.Net.ReadTimeout = time.Duration(timeoutSec) * time.Second
	config.Net.WriteTimeout = time.Duration(timeoutSec) * time.Second

	client, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{
		client: client,
		topic:  topic,
	}, nil
}

func (p *Producer) PublishMessage(key string, msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := p.client.SendMessage(kafkaMsg)
	if err != nil {
		return fmt.Errorf("send error: %w", err)
	}

	log.Printf("Message sent to topic %s, partition %d, offset %d", p.topic, partition, offset)
	return nil
}

func (p *Producer) Close() error {
	return p.client.Close()
}
