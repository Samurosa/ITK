package grpc

import (
	"ITK_Code/m/v2/internal/core/spot"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/spot"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedSpotInstrumentServiceServer
	log  *zap.Logger
	spot spot.Service
}

func RegisterUserService(grpc *grpc.Server, spot spot.Service, log *zap.Logger) {
	pb.RegisterSpotInstrumentServiceServer(grpc, &Server{log: log, spot: spot})
}
