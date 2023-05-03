package test

import (
	"testing"

	"github.com/jerrykhh/job-queue/server/utils/pwd"
	"github.com/stretchr/testify/assert"
)

func TestPwdHash(t *testing.T) {

	password := "testPassword"

	hashed, err := pwd.HashPassword(password)
	assert.Nil(t, err)
	assert.Nil(t, pwd.ComparePwd(hashed, password))
	assert.NotNil(t, pwd.ComparePwd("testPassword1", hashed))
}
