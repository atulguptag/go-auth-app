package utils

import (
	"go-auth-app/models"

	"github.com/dgrijalva/jwt-go"
)

func ParseToken(tokenString string) (claims *models.Claims, err error) {
	access_token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(access_token *jwt.Token) (interface{}, error) {
		return []byte("my_secret_key"), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := access_token.Claims.(*models.Claims)

	if !ok {
		return nil, err
	}

	return claims, nil
}
