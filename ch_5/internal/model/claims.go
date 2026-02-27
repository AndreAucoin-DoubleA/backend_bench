package model

import j "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	j.RegisteredClaims
}
