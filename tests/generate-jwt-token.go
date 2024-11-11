package tests

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type jwtCustomClaims struct {
	CustomerId string `json:"customerId"`
	jwt.StandardClaims
}

func GenerateJwtAccessToken(customerId string, authAccessToken string) (string, error) {
	claims := &jwtCustomClaims{
		CustomerId: customerId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(authAccessToken))
}
