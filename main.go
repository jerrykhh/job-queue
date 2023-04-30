package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/jerrykhh/job-queue/grpc/pb"
	"github.com/jerrykhh/job-queue/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := server.LoadConfig(".")
	if err != nil {
		fmt.Println("load config file failed")
		log.Fatalln(err)
	}
	runGrpcServer(config)
}

func runGrpcServer(config server.Config) {
	serv, err := server.NewServer(config)
	if err != nil {
		fmt.Println(err)
	}

	// gprcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, serv)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("start gRPC server at %s\n", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		fmt.Println(err)
	}
}
