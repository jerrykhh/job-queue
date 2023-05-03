package test

import (
	"context"
	"testing"

	"github.com/jerrykhh/job-queue/grpc/pb"
	"github.com/stretchr/testify/assert"
)

var client pb.UserServiceClient

func init() {
	client = ConnUserServer()
}

func TestLogin(t *testing.T) {
	res, err := client.Login(context.Background(), &pb.User{
		Username: "root",
		Password: "test",
	})
	assert.Nil(t, err)
	assert.Equal(t, res.Username, "root")
}

func TestLoginIncorrectPwd(t *testing.T) {
	_, err := client.Login(context.Background(), &pb.User{
		Username: "root",
		Password: "",
	})
	assert.NotNil(t, err)
}

func TestLoginIncorrectUsername(t *testing.T) {
	_, err := client.Login(context.Background(), &pb.User{
		Username: "a",
		Password: "test",
	})
	assert.NotNil(t, err)
}
