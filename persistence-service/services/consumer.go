package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mohammad-ammad/persistence-service/config"
	"github.com/mohammad-ammad/persistence-service/models"
)

func ProcessMessages() {
	for {

		msg, err := config.KafkaReader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading from Kafka: %v", err)
			continue
		}

		var chatMsg models.Message
		err = json.Unmarshal(msg.Value, &chatMsg)
		if err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		result := config.DB.Create(&chatMsg)
		if result.Error != nil {
			log.Printf("Error storing message in DB: %v", result.Error)
			continue
		}

		log.Printf("Stored message from user %d: %s", chatMsg.UserID, chatMsg.Content)

		if err := config.KafkaReader.CommitMessages(context.Background(), msg); err != nil {
			log.Printf("Error committing message: %v", err)
		}
	}
}
