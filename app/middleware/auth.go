package middleware

import (
	"advanceauth/backend/app/handler"
	"advanceauth/backend/app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		bearerToken := strings.Split(header, "Bearer ")
		if len(bearerToken) != 2 {
			handler.Error(c, http.StatusUnauthorized, utils.AUTH_MISSING_JWT, "Unauthorized")
			return
		}

		payload, err := validateToken(bearerToken[1])
		if err != nil && string(err.Error()) != "Token is expired" {
			handler.Error(c, http.StatusUnauthorized, utils.AUTH_WRONG_JWT, err.Error())
			return
		}
		if err != nil && string(err.Error()) == "Token is expired" {
			c.Set("isTokenExpired", true)
		} else {
			c.Set("UUID", payload.UUID)
		}
		c.Set("loginToken", bearerToken[1])
		c.Next()
	}
}

func validateToken(tokenInput string) (*utils.PayloadToken, error) {
	res, err := jwt.Parse(tokenInput, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(utils.GetEnv("JWT_SecretKey")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := res.Claims.(jwt.MapClaims)
	if !ok || !res.Valid {
		return nil, errors.New("Unauthorized")
	}

	payload := claims["UUID"]
	var payloadToken = utils.PayloadToken{}
	payloadByte, _ := json.Marshal(payload)
	err = json.Unmarshal(payloadByte, &payloadToken)
	if err != nil {
		return nil, err
	}

	return &payloadToken, nil
}
