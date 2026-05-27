package main

import (
	"ITK_Code/m/v2/internal/handler"
	"log"
	"net"

	pb "github.com/Truncklin/exchange-contract/generated"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal("error listen port 50052: ", err)
	}
	grpcServer := grpc.NewServer()
	spotService := handler.NewSpotService()
	pb.RegisterSpotInstrumentServiceServer(
		grpcServer,
		spotService,
	)
	log.Println("server starting")
	grpcServer.Serve(lis)
}
