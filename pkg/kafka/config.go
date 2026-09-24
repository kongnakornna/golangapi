package kafka

import (
	"os"
)

type Config struct {
	Brokers []string
	Topic   string
	GroupID string
}

func LoadConfig() Config {
	return Config{
		Brokers: []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
		Topic:   getEnv("KAFKA_TOPIC", "icmon-events"),
		GroupID: getEnv("KAFKA_GROUP_ID", "icmon-consumer-group"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
