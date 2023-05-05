package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jerrykhh/job-queue/server"
)

func main() {
	config, err := server.LoadConfig(".env")
	if err != nil {
		fmt.Println("load config file failed")
		log.Fatalln(err)
	}
	fmt.Println("run")
	err = runGrpcServer(config)
	if err != nil {
		error.Error(err)
	}
}

func runGrpcServer(config server.Config) error {
	serv, err := server.NewServer(config)
	if err != nil {
		return err
	}

	// gprcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer, _ := serv.RunGrpcServer()
	//defer closeFunc()
	// reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		return err
	}

	fmt.Printf("start gRPC server at %s\n", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}
