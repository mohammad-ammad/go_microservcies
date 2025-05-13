package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/mohammad-ammad/ws-service/config"
)

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(config.Env("JWT_SECRET", "secret")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
