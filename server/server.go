package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jerrykhh/job-queue/grpc/pb"
	server_queue "github.com/jerrykhh/job-queue/server/queue"
	"github.com/jerrykhh/job-queue/server/utils/jwt"
	"github.com/jerrykhh/job-queue/server/utils/pwd"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type Server struct {
	pb.JobQueueServiceServer
	pb.UserServiceServer

	config      Config
	jwtCreator  *jwt.JWTCreator
	rootHashPwd string
	queues      map[string]*server_queue.JobQueue
	redis       *redis.Client

	// server
	Lis *bufconn.Listener
}

func NewServer(config Config) (*Server, error) {

	jwtTokenCreator, err := jwt.NewJWTCreator(config.TokenHashKey)

	if err != nil {
		return nil, err
	}

	serv := &Server{
		config:     config,
		jwtCreator: jwtTokenCreator,
		Lis:        bufconn.Listen(1024 * 1024),
	}

	hashPwd, err := pwd.HashPassword(config.RootPwd)

	if err != nil {
		return nil, err
	}

	serv.rootHashPwd = hashPwd
	serv.queues = make(map[string]*server_queue.JobQueue)
	serv.redis = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.RedisAddress, config.RedisPort),
	})

	err = serv.LoadQueueFromRedis()

	if err != nil {
		return nil, err
	}

	return serv, nil

}

func (server *Server) RunGrpcServer() (*grpc.Server, func()) {
	grpcServer := grpc.NewServer()
	pb.RegisterJobQueueServiceServer(grpcServer, server)
	pb.RegisterUserServiceServer(grpcServer, server)
	go func() {
		if err := grpcServer.Serve(server.Lis); err != nil {
			fmt.Printf("error servinf server: %v", err)
		}
	}()

	closeFunc := func() {
		err := server.Lis.Close()
		if err != nil {
			fmt.Printf("error closing listener: %v\n", err)
		}
		grpcServer.Stop()
	}

	return grpcServer, closeFunc
}

func (server *Server) LoadQueueFromRedis() error {
	result, err := server.ListJobQueue()
	if err != nil {
		return err
	}

	for _, jobQueue := range result {
		server.queues[jobQueue.Id] = jobQueue
	}

	return nil

}

func (server *Server) CompareUsername(name string) error {
	if name != server.config.RootUsername {
		return fmt.Errorf("incorrect username")
	}
	return nil
}

func (server *Server) NewJobQueue(name string, runEverySec, seed, dequeueCount int) (*server_queue.JobQueue, error) {
	ctx := context.Background()

	newQueue, err := server_queue.NewJobQueue(name, runEverySec, seed, dequeueCount, server.redis)
	if err != nil {
		return nil, err
	}
	server.queues[newQueue.Id] = newQueue

	queueJson, err := json.Marshal(newQueue)
	if err != nil {
		return nil, err
	}

	server.redis.RPush(ctx, "job-queue", queueJson).Result()

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

	queueJson, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}

	server.redis.LRem(context.Background(), "job-queue", 1, queueJson).Result()
	q.Pause = true
	delete(server.queues, queueId)

	return q, nil
}

func (server *Server) ListJobQueue() ([]*server_queue.JobQueue, error) {
	ctx := context.Background()
	result, err := server.redis.LRange(ctx, "job-queue", 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to load queue data from redis")
	}

	queues := make([]*server_queue.JobQueue, len(result))

	for i, data := range result {
		var jobQueue server_queue.JobQueue
		err := json.Unmarshal([]byte(data), &jobQueue)
		if err != nil {
			log.Println(err)
		}

		queues[i] = &jobQueue
	}
	return queues, nil
}

func (server *Server) TriggerBgSave() error {
	_, err := server.redis.BgSave(context.Background()).Result()
	return err
}
