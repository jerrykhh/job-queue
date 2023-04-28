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

func (jwtCreator *JWTCreator) CreateToken(username string, duration time.Duration) (string, jwt.MapClaims, error) {

	payload := jwt.MapClaims{
		"sub": username,
		"iat": time.Now(),
		"exp": time.Now().Add(duration),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	token, err := t.SignedString([]byte(jwtCreator.secretKey))
	return token, payload, err
}

func (jwtCreator *JWTCreator) VerifyToken(token string) (jwt.Claims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtCreator.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, parsedToken.Claims.Valid()
}