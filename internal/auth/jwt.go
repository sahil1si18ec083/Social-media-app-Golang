package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secretKey []byte
	expiry    time.Duration
}

func (j *JWT) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	subject, err := mapClaims.GetSubject()
	if err != nil {
		return nil, err
	}

	username, ok := mapClaims["username"].(string)
	if !ok || username == "" {
		return nil, fmt.Errorf("invalid username claim")
	}

	return &Claims{
		Subject:  subject,
		Username: username,
	}, nil
}

func NewJWT(secret string, exp time.Duration) *JWT {

	return &JWT{
		secretKey: []byte(secret),
		expiry:    exp,
	}
}

func (j *JWT) GenerateToken(userID int64, username string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":      strconv.FormatInt(userID, 10),
		"username": username,
		"iat":      now.Unix(),
		"exp":      now.Add(j.expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil

}
