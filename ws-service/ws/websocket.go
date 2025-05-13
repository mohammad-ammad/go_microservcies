package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mohammad-ammad/ws-service/config"
	"github.com/mohammad-ammad/ws-service/middleware"
	"github.com/mohammad-ammad/ws-service/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var Clients = make(map[*websocket.Conn]bool)

func HandleWebSocket(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authentication token"})
		return
	}

	parsedToken, err := middleware.ValidateToken(token)
	if err != nil || !parsedToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not parse token claims"})
		return
	}

	log.Printf("Claims: %v", claims)

	userID := uint(claims["user_id"].(float64))
	username := claims["username"].(string)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}

	Clients[conn] = true
	log.Printf("Client connected: User %s (ID: %d)", username, userID)

	HandleConnection(c.Request.Context(), conn, userID, username)
}

func HandleConnection(ctx context.Context, conn *websocket.Conn, userID uint, username string) {
	defer func() {
		conn.Close()
		delete(Clients, conn)
		log.Printf("Client disconnected: User %s (ID: %d)", username, userID)
	}()

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		msg := model.Message{
			UserID:    userID,
			Username:  username,
			Content:   string(msgBytes),
			Timestamp: time.Now(),
		}

		msgJSON, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshaling message: %v", err)
			continue
		}

		log.Printf("msgJSON: %s", msgJSON)

		err = config.KafkaWriter.WriteMessages(ctx, kafka.Message{
			Value: msgJSON,
		})

		if err != nil {
			log.Printf("Error publishing to Kafka: %v", err)
		}
	}
}
