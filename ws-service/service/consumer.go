package service

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
	"github.com/mohammad-ammad/ws-service/config"
	"github.com/mohammad-ammad/ws-service/ws"
)

func ConsumeMessages() {
	for {
		log.Println("Waiting for messages from Kafka...")

		msg, err := config.KafkaReader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading from Kafka: %v", err)
			continue
		}

		for client := range ws.Clients {
			err := client.WriteMessage(websocket.TextMessage, msg.Value)
			log.Printf("Broadcasting message to client: %s", msg.Value)
			if err != nil {
				log.Printf("Error sending to client: %v", err)
				client.Close()
				delete(ws.Clients, client)
			}
		}
	}
}
