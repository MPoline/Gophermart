package models

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	TokenTTL  = time.Hour * 24 * 7 // 1 неделя
	secretKey = "my_very_secret_key"
)

type Claims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

func GenerateAuthToken(login string) (string, error) {
	claims := &Claims{
		Login: login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenTTL).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ValidateAuthToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Login, nil
	}

	return "", errors.New("invalid token")
}
