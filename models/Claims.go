package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Email  string `json:"email"`
	UserID uint   `json:"user_id"`
	jwt.StandardClaims
}
