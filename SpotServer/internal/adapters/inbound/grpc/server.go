package grpc

import (
	"ITK_Code/m/v2/internal/core/market"
	"ITK_Code/m/v2/internal/core/spot"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/spot"
	"go.uber.org/zap"
)

type Server struct {
	pb.UnimplementedSpotInstrumentServiceServer
	log    *zap.Logger
	spot   spot.Service
	market market.Service
}
