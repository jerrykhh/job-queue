package server

import (
	"github.com/jerrykhh/job-queue/grpc/pb"
	server_queue "github.com/jerrykhh/job-queue/server/queue"
	"github.com/jerrykhh/job-queue/server/utils/jwt"
	"github.com/jerrykhh/job-queue/server/utils/pwd"
)

type Server struct {
	pb.UnimplementedJobQueueServiceServer
	pb.UserServiceServer

	config      Config
	jwtCreator  *jwt.JWTCreator
	rootHashPwd string
	queues      map[string]*server_queue.JobQueue
}

func NewServer(config Config) (*Server, error) {

	jwtTokenCreator, err := jwt.NewJWTCreator(config.TokenHashKey)

	if err != nil {
		return nil, err
	}

	serv := &Server{
		config:     config,
		jwtCreator: jwtTokenCreator,
	}

	hashPwd, err := pwd.HashPassword(config.RootPwd)

	if err != nil {
		return nil, err
	}

	serv.rootHashPwd = hashPwd
	serv.queues = make(map[string]*server_queue.JobQueue)
	return serv, nil

}

func (server *Server) NewJobQueue(name string, runEverySec, seed, dequeueCount int) (*server_queue.JobQueue, error) {
	newQueue, err := server_queue.NewJobQueue(name, runEverySec, seed, dequeueCount)
	if err != nil {
		return nil, err
	}
	server.queues[newQueue.Id] = newQueue
	return newQueue, nil
}
