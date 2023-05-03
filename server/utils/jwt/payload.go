package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Payload struct {
	jwt.StandardClaims
	Username string
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid for jwt")
	}

	return &Payload{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			Id:        id.String(),
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}, nil
}

func (p *Payload) Valid() error {
	return p.StandardClaims.Valid()
}
