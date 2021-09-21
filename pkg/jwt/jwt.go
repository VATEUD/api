package jwt

import (
	"auth/utils"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	*jwt.Token
}

func New(tokenString string) (*Token, error) {
	secret := utils.Getenv("JWT_SECRET", "")

	if len(secret) < 1 {
		return nil, errors.New("JWT Secret not available")
	}

	token, err := jwt.Parse(tokenString, func (token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("Unexpected signing method: %s", token.Header["alg"]))
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	return &Token{token}, nil
}
