package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Payload struct {
	Id        string
	Username  string
	IssuedAt  time.Time
	ExpiresAt time.Time
	jwt.StandardClaims
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid for jwt")
	}

	return &Payload{
		Id:        id.String(),
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}, nil
}
