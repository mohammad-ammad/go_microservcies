package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mohammad-ammad/auth-service/controllers"
	"github.com/mohammad-ammad/auth-service/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.GET("/me", middleware.AuthMiddleware(), controllers.Me)
	}
}
