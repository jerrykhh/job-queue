package server

import (
	"github.com/jerrykhh/job-queue/grpc/pb"
	"github.com/jerrykhh/job-queue/server/utils/jwt"
)

type Server struct {
	pb.UnimplementedJobQueueServiceServer
	config     Config
	jwtCreator *jwt.JWTCreator
}

func NewServer(config Config) (*Server, error) {

	jwtTokenCreator, err := jwt.NewJWTCreator(config.TokenHashKey)

	if err != nil {
		return nil, err
	}

	return &Server{
		config:     config,
		jwtCreator: jwtTokenCreator,
	}, nil

}
