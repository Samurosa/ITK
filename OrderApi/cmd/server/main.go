package main

import (
	config "ITK_Code/m/v2/config"
	"ITK_Code/m/v2/internal/handler"
	"fmt"
	"log"
	"net"

	commonInterceptor "github.com/Samurosa/exchange-common/interceptor"
	metrics "github.com/Samurosa/exchange-common/metrics"
	pb "github.com/Samurosa/exchange-contract/generated"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	cfg, err := config.Load("ITK_Code/m/v2/internal/config/local.yaml")
	if err != nil {
		logger.Fatal(
			"config not load",
			zap.Error(err),
		)
	}

	orderAddr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)

	lis, err := net.Listen("tcp", orderAddr)
	if err != nil {
		log.Fatal("error listen port 50051: ", err)
	}

	spotConn, err := grpc.NewClient(
		cfg.Clients.User.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer spotConn.Close()

	userConn, err := grpc.NewClient(
		cfg.Clients.Spot.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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
