package main

import (
	"ITK_Code/m/v1/internal/handler"
	pb "ITK_Code/m/v1/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}

	grpcserver := grpc.NewServer()
	userservice := handler.NewChatApi()
	pb.RegisterChatAPIServer(
		grpcserver,
		userservice,
	)
	log.Println("start server")
	grpcserver.Serve(lis)
}
