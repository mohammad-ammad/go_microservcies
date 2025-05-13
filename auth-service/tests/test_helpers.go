package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/mohammad-ammad/auth-service/config"
	"github.com/mohammad-ammad/auth-service/controllers"
	"github.com/mohammad-ammad/auth-service/models"
	"github.com/mohammad-ammad/auth-service/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"golang.org/x/crypto/bcrypt"
)

var DB *gorm.DB

func SetupTestDB(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	config.DB = db

	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("Failed to migrate models: %v", err)
	}

	DB = db
}

func CleanupTestDB() {
	// Optional: Cleanup tasks, like dropping tables
	if config.DB != nil {
		config.DB.Migrator().DropTable(&models.User{})
	}
}

func setupTestRouter() *gin.Engine {
	SetupTestDB(nil)
	r := gin.Default()
	r.POST("/auth/register", controllers.Register)
	r.POST("/auth/login", controllers.Login)

	authGroup := r.Group("/auth")
	authGroup.Use(middleware.AuthMiddleware())
	authGroup.GET("/me", controllers.Me)
	return r
}


func createTestUser(email, password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := models.User{
		Username: "testuser",
		Email:    email,
		Password: string(hashedPassword),
	}

	DB.Create(&user)
}
