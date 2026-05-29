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
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal("error listen port 50051: ", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	spotConn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer spotConn.Close()

	userConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTransportCredentials(commonInterceptor.RequestID()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer userConn.Close()

	spotClient := pb.NewSpotInstrumentServiceClient(spotConn)
	userClient := pb.NewUserServiceClient(userConn)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				commonInterceptor.RequestID(),
				commonInterceptor.NewLogger(logger).Unary(),
				commonInterceptor.Recovery(logger),
				metrics.Prometheus(),
			),
		),
	)

	orderService := handler.NewOrderService(spotClient, userClient)
	pb.RegisterOrderServiceServer(
		grpcServer,
		orderService,
	)

	logger.Info(
		"grpc server started",
		zap.String("port", ":50053"),
	)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("grpc server failed",
			zap.Error(err),
		)
	}
}
