package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type PayloadToken struct {
	UUID string
}

func GenerateToken(UUID *PayloadToken) (string, error) {
	claims := jwt.MapClaims{
		"UUID": UUID,
		"exp":  time.Now().Add(time.Hour * 24 * 14).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(GetEnv("JWT_SecretKey")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
