package main

import (
	"ITK_Code/m/v2/internal/handler"
	pb "ITK_Code/m/v2/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("error listen port 50051: ", err)
	}
	grpcServer := grpc.NewServer()
	userService := handler.NewUserService()
	pb.RegisterUserServiceServer(
		grpcServer,
		userService,
	)
	log.Println("server starting")
	grpcServer.Serve(lis)
}
