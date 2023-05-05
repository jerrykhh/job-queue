package test

import (
	"fmt"

	"github.com/jerrykhh/job-queue/grpc/pb"
	"github.com/jerrykhh/job-queue/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
}

func ConnServer() *grpc.ClientConn {
	config, err := server.LoadConfig("../.env")
	if err != nil {
		fmt.Println("read .env failed")
	}
	// serv, err := server.NewServer(config)
	// if err != nil {
	// 	fmt.Println("failed to New Server")
	// }
	// _, closeFunc := serv.RunGrpcServer()
	fmt.Println(config.GRPCServerAddress)
	fmt.Println("ConnServer")
	// dialAddr := config.GRPCServerAddress

	conn, err := grpc.Dial(config.GRPCServerAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("error connecting to server: %v\n", err)
	}

	return conn
}

func ConnJobqueueServer() pb.JobQueueServiceClient {
	conn := ConnServer()
	// defer conn.Close()
	client := pb.NewJobQueueServiceClient(conn)
	return client
}

func ConnUserServer() pb.UserServiceClient {
	conn := ConnServer()
	// defer conn.Close()
	client := pb.NewUserServiceClient(conn)
	return client
}

func int32Ptr(v int) *int32 {
	p := int32(v)
	return &p
}
