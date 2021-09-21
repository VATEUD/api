package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

func GenerateNewJWT(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := Getenv("JWT_SECRET", "")

	if len(secret) < 1 {
		return "", errors.New("JWT Secret not provided")
	}

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
