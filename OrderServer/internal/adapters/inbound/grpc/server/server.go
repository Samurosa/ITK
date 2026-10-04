package server

import (
	"ITK_Code/m/v2/internal/application"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/order"
	"google.golang.org/grpc"
)

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	order *application.OrderService
}

func NewOrderServer(grpc *grpc.Server, order *application.OrderService) {
	pb.RegisterOrderServiceServer(grpc, &OrderServer{order: order})
}
