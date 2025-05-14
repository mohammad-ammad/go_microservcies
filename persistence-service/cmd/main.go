package main

import (
	"github.com/mohammad-ammad/persistence-service/config"
	"github.com/mohammad-ammad/persistence-service/models"
	"github.com/mohammad-ammad/persistence-service/services"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	config.DB.AutoMigrate(&models.User{}, &models.Message{})

	config.InitializeKafka()

	services.ProcessMessages()

}
