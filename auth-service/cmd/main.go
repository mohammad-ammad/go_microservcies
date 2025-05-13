package main

import (
	"github.com/mohammad-ammad/auth-service/config"
	"github.com/mohammad-ammad/auth-service/models"
	"github.com/mohammad-ammad/auth-service/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{}) 

	r := gin.Default()
	routes.RegisterRoutes(r)

	r.Run(":" + config.Env("PORT", "8080"))
}
