package config

import (
	"log"

	"github.com/segmentio/kafka-go"
)

var (
	KafkaWriter *kafka.Writer
	KafkaReader *kafka.Reader
)

func InitializeKafka() {
	KafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(Env("KAFKA_BROKER", "localhost:9092")),
		Topic:    "chat-messages",
		Balancer: &kafka.LeastBytes{},
	}

	KafkaReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{Env("KAFKA_BROKER", "localhost:9092")},
		Topic:    "chat-messages",
		GroupID:  "websocket-service",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	log.Println("Kafka initialized")
}
