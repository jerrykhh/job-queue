package server

import (
	"fmt"

	"github.com/jerrykhh/job-queue/grpc/pb"
	server_queue "github.com/jerrykhh/job-queue/server/queue"
	"github.com/jerrykhh/job-queue/server/utils/jwt"
	"github.com/jerrykhh/job-queue/server/utils/pwd"
)

type Server struct {
	pb.JobQueueServiceServer
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

func (server *Server) CompareUsername(name string) error {
	if name != server.config.RootUsername {
		return fmt.Errorf("incorrect username")
	}
	return nil
}

func (server *Server) NewJobQueue(name string, runEverySec, seed, dequeueCount int) (*server_queue.JobQueue, error) {
	newQueue, err := server_queue.NewJobQueue(name, runEverySec, seed, dequeueCount, server.config.RedisAddress, server.config.RedisPort)
	if err != nil {
		return nil, err
	}
	server.queues[newQueue.Id] = newQueue
	return newQueue, nil
}

func (server *Server) GetJobQueue(queueId string) (*server_queue.JobQueue, error) {
	if q, ok := server.queues[queueId]; ok {
		return q, nil
	} else {
		return nil, fmt.Errorf("queue id not found")
	}
}

func (server *Server) RemoveJobQueue(queueId string) (*server_queue.JobQueue, error) {
	q, err := server.GetJobQueue(queueId)
	if err != nil {
		return nil, err
	}
	delete(server.queues, queueId)

	return q, nil
}
