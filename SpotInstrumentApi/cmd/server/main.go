package main

import (
	"ITK_Code/m/v2/internal/handler"
	"log"
	"net"

	commonInterceptor "github.com/Samurosa/exchange-common/interceptor"
	metrics "github.com/Samurosa/exchange-common/metrics"
	pb "github.com/Samurosa/exchange-contract/generated"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal("error listen port 50052: ", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

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

	spotService := handler.NewSpotService()

	pb.RegisterSpotInstrumentServiceServer(
		grpcServer,
		spotService,
	)

	logger.Info(
		"grpc server started",
		zap.String("port", ":50052"),
	)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("grpc server failed",
			zap.Error(err),
		)
	}
}
