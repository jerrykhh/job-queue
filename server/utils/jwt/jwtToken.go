package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTCreator struct {
	secretKey string
}

const minSecretKeySize = 32

func NewJWTCreator(secretKey string) (*JWTCreator, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("secretKey size must be at least %d", minSecretKeySize)
	}

	return &JWTCreator{secretKey}, nil
}

func (jwtCreator *JWTCreator) CreateToken(username string, duration time.Duration) (string, *Payload, error) {

	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	token, err := t.SignedString([]byte(jwtCreator.secretKey))
	return token, payload, err
}

func (jwtCreator *JWTCreator) VerifyToken(token string) (jwt.Claims, error) {
	claims := &Payload{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtCreator.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if ok := parsedToken.Valid; ok {
		return claims, nil
	}

	return nil, fmt.Errorf("parsed token is invalid")
}
