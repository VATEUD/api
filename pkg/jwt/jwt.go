package jwt

import (
	"api/utils"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

type Token struct {
	*jwt.Token
	MapClaims jwt.MapClaims
}

func New(tokenString string) (*Token, error) {
	secret := utils.Getenv("JWT_SECRET", "")

	if len(secret) < 1 {
		return nil, errors.New("JWT Secret not available")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("Unexpected signing method: %s", token.Header["alg"]))
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("Token claims could not be converted")
	}

	return &Token{token, claims}, nil
}
