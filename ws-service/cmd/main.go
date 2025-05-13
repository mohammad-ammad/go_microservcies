package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mohammad-ammad/ws-service/config"

	"github.com/mohammad-ammad/ws-service/service"
	"github.com/mohammad-ammad/ws-service/ws"
)

func main() {
	config.LoadEnv()

	config.InitializeKafka()
	go service.ConsumeMessages()

	r := gin.Default()
	r.GET("/ws", ws.HandleWebSocket)

	r.Run(":" + config.Env("PORT", "8080"))
}
