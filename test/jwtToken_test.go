package test

import (
	"testing"
	"time"

	"github.com/jerrykhh/job-queue/server/utils/jwt"
	"github.com/stretchr/testify/assert"
)

var secretKey string

func init() {
	secretKey = "306247005340623062470053406230624700534062a"
}

func TestNewTokenCreatorSecretKeyErr(t *testing.T) {
	secretKey := "123"
	_, err := jwt.NewJWTCreator(secretKey)
	assert.NotNil(t, err)
}

func TestNewTokenCreator(t *testing.T) {
	creator, err := jwt.NewJWTCreator(secretKey)
	assert.Nil(t, err)
	assert.NotNil(t, creator)
}

func TestCreateToken(t *testing.T) {
	creator, _ := jwt.NewJWTCreator(secretKey)
	token, payload, err := creator.CreateToken("test", time.Duration(10))
	assert.Nil(t, err)
	assert.Equal(t, payload.Username, "test")
	assert.Nil(t, payload.Valid())
	assert.Equal(t, len(token) > 0, true)
}

func TestValidToken(t *testing.T) {
	creator, _ := jwt.NewJWTCreator(secretKey)
	token, _, _ := creator.CreateToken("test", time.Duration(100))
	claims, err := creator.VerifyToken(token)
	assert.Nil(t, err)
	assert.NoError(t, claims.Valid())
}

func TestInvalidToken(t *testing.T) {
	creator, _ := jwt.NewJWTCreator(secretKey)
	_, err := creator.VerifyToken("sdafasdfasdf")
	assert.NotNil(t, err)
}
