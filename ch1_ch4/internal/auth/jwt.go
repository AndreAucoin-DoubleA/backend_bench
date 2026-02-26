package auth

import (
	"backend_bench/internal/model"
	"errors"
	"time"

	j "github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string, email string, secretKey string) (string, error) {
	var jwtKey = []byte(secretKey)
	expiration := time.Now().Add(24 * time.Hour)
	claims := &model.Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: j.RegisteredClaims{
			ExpiresAt: j.NewNumericDate(expiration),
		},
	}
	token := j.NewWithClaims(j.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func VerifyJWT(tokenStr string, secretKey string) (*model.Claims, error) {
	var jwtKey = []byte(secretKey)
	token, err := j.ParseWithClaims(tokenStr, &model.Claims{}, func(token *j.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return token.Claims.(*model.Claims), nil
}
