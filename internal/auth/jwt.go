package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secretKey []byte
}

func (j *JWT) ValidateToken(tokenString string) error {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {

		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func NewJWT(secret string) *JWT {

	return &JWT{
		secretKey: []byte(secret),
	}
}
