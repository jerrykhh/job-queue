package test

import (
	"context"
	"testing"

	"github.com/jerrykhh/job-queue/grpc/pb"
	"github.com/stretchr/testify/assert"
)

var userClient pb.UserServiceClient

func init() {
	userClient = ConnUserServer()
}

func TestLogin(t *testing.T) {
	res, err := userClient.Login(context.Background(), &pb.User{
		Username: "root",
		Password: "test",
	})
	assert.Nil(t, err)
	assert.Equal(t, res.Username, "root")
}

func TestLoginIncorrectPwd(t *testing.T) {
	_, err := userClient.Login(context.Background(), &pb.User{
		Username: "root",
		Password: "",
	})
	assert.NotNil(t, err)
}

func TestLoginIncorrectUsername(t *testing.T) {
	_, err := userClient.Login(context.Background(), &pb.User{
		Username: "a",
		Password: "test",
	})
	assert.NotNil(t, err)
}
