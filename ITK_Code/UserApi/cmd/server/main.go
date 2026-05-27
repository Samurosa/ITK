package main

import (
	"ITK_Code/m/v2/internal/handler"
	"log"
	"net"

	commonInterceptor "github.com/Truncklin/exchange-common/interceptor"
	metrics "github.com/Truncklin/exchange-common/metrics"
	pb "github.com/Truncklin/exchange-contract/generated"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("error listen port 50051: ", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				commonInterceptor.RequestID(),
				commonInterceptor.NewLogger(&zap.Logger{}).Unary(),
				commonInterceptor.Recovery(),
				metrics.Prometheus(),
			),
		),
	)
	userService := handler.NewUserService()
	pb.RegisterUserServiceServer(
		grpcServer,
		userService,
	)
	log.Println("server starting")
	grpcServer.Serve(lis)
}
