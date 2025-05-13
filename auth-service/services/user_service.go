package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mohammad-ammad/auth-service/config"
	"github.com/mohammad-ammad/auth-service/dto"
	"github.com/mohammad-ammad/auth-service/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(req dto.RegisterRequest) error {
	var existing models.User

	if err := config.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return errors.New("user already exists with this email")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	return config.DB.Create(&user).Error
}

func AuthenticateUser(req dto.LoginRequest) (string, error) {
	var user models.User

	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.Env("JWT_SECRET", "secret")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserByID(userID uint) (models.User, error) {
	var user models.User

	if err := config.DB.First(&user, userID).Error; err != nil {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}
