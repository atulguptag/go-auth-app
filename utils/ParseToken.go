package utils

import (
	"errors"
	"go-auth-app/models"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

func ParseToken(tokenString string) (*models.Claims, error) {
	if tokenString == "" {
		return nil, errors.New("empty token")
	}
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.Claims)

	if !ok {
		return nil, err
	}

	return claims, nil
}
