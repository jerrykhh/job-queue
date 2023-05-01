package server

import (
	"context"

	pb "github.com/jerrykhh/job-queue/grpc/pb"
	"github.com/jerrykhh/job-queue/server/utils/pwd"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) Login(ctx context.Context, req *pb.User) (*pb.LoginResponse, error) {

	err := server.CompareUsername(req.GetUsername())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	err = pwd.ComparePwd(server.rootHashPwd, req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "incorrect password")
	}

	accessToken, accessPayload, err := server.jwtCreator.CreateToken(req.GetUsername(), server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create JWT access token ")
	}

	refreshToken, refeshPayload, err := server.jwtCreator.CreateToken(req.GetUsername(), server.config.RefreshTokenDueation)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create JWT refresh token")
	}

	return &pb.LoginResponse{
		Username:          req.GetUsername(),
		AccessToken:       accessToken,
		RefreshToken:      refreshToken,
		AccessTokenExpAt:  timestamppb.New(accessPayload.ExpiresAt),
		RefreshTokenExpAt: timestamppb.New(refeshPayload.ExpiresAt),
	}, nil
}
