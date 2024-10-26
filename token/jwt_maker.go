package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const MIN_SECRET_KEY_SIZE = 32

type JWTMaker struct {
	secretKey string
}

func (j *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(j.secretKey))
}

func (j *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken

		}
		return []byte(j.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	// token is invalid or has expired
	if err != nil {
		return nil, err
	}

	if !jwtToken.Valid {
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < MIN_SECRET_KEY_SIZE {
		return nil, fmt.Errorf("invalid key size, must be %d characters", MIN_SECRET_KEY_SIZE)
	}

	return &JWTMaker{secretKey}, nil
}
