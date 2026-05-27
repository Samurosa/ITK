package main

import (
	"ITK_Code/m/v2/internal/handler"
	"log"
	"net"

	pb "github.com/Truncklin/exchange-contract/generated"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal("error listen port 50051: ", err)
	}

	spotConn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	spotClient := pb.NewSpotInstrumentServiceClient(spotConn)

	grpcServer := grpc.NewServer()
	orderService := handler.NewOrderService(spotClient)
	pb.RegisterOrderServiceServer(
		grpcServer,
		orderService,
	)
	log.Println("server starting")
	grpcServer.Serve(lis)
}
