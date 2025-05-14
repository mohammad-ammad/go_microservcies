package config

import (
	"log"

	"github.com/segmentio/kafka-go"
)

var KafkaReader *kafka.Reader

func InitializeKafka() {
	KafkaReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{Env("KAFKA_BROKER", "localhost:9092")},
		Topic:          "chat-messages",
		GroupID:        "persistence-service",
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: 0,
	})

	log.Println("Kafka consumer initialized")
}
