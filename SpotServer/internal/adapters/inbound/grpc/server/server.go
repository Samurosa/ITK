package server

import (
	"ITK_Code/m/v2/internal/application"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/spot"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedSpotInstrumentServiceServer
	spot *application.Spot
}

func RegisterSpotService(grpc *grpc.Server, spot *application.Spot) {
	pb.RegisterSpotInstrumentServiceServer(grpc, &Server{spot: spot})
}
