package server

import (
	"ITK_Code/m/v2/internal/application"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/spot"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedSpotInstrumentServiceServer
	log  *zap.Logger
	spot application.Spot
}

func RegisterSpotService(grpc *grpc.Server, spot application.Spot, log *zap.Logger) {
	pb.RegisterSpotInstrumentServiceServer(grpc, &Server{log: log, spot: spot})
}
