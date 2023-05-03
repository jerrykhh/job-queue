package test

import (
	"testing"
	"time"

	"github.com/jerrykhh/job-queue/server/utils/jwt"
	"github.com/stretchr/testify/assert"
)

func TestPayloadInValid(t *testing.T) {
	payload, err := jwt.NewPayload("test", time.Duration(1)*time.Second)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)
	assert.Error(t, payload.Valid())
}

func TestPayloadValid(t *testing.T) {
	payload, err := jwt.NewPayload("test", time.Duration(1)*time.Second)
	assert.Nil(t, err)
	assert.Nil(t, payload.Valid())
}
