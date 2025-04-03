package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	SecretKey []byte
}

func ConfigureJWT(key string) *JWTService {
	return &JWTService{SecretKey: []byte(key)}
}

func (j *JWTService) GenerateJWT(username, id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"id":       id,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})
	println("key", j.SecretKey)
	return token.SignedString(j.SecretKey)
}

func (j *JWTService) ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неверный метод подписи")
		}
		return (j.SecretKey), nil
	})
}

func (j *JWTService) GetExpTime(token *jwt.Token) float64 {
	return token.Claims.(jwt.MapClaims)["exp"].(float64)
}

func (j *JWTService) GetUsername(token *jwt.Token) string {
	return token.Claims.(jwt.MapClaims)["username"].(string)
}
