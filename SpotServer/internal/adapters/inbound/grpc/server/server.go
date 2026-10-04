package server

import (
	"ITK_Code/m/v2/internal/application"
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/spot/models"
	"context"
	"time"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/spot"
	"google.golang.org/grpc"
)

type Spot interface {
	CreateSpot(context.Context, models.CreateSpot) (string, time.Time, error)
	GetSpot(context.Context, string) (dto.Spot, error)
	EnableSpot(context.Context, string) error
	DisableSpot(context.Context, string) error
	ListSpots(context.Context, models.ListSpotsRequest) ([]dto.SpotListItem, string, bool, error)
}

type Server struct {
	pb.UnimplementedSpotInstrumentServiceServer
	spot *application.Spot
}

func RegisterSpotService(grpc *grpc.Server, spot *application.Spot) {
	pb.RegisterSpotInstrumentServiceServer(grpc, &Server{spot: spot})
}
