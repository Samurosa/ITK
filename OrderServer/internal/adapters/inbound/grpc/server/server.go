package server

import (
	"ITK_Code/m/v2/internal/core/order/service"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/order"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	order service.Order
	log   *zap.Logger
}

func NewOrderServer(grpc *grpc.Server, order service.Order, log *zap.Logger) {
	pb.RegisterOrderServiceServer(grpc, &OrderServer{order: order, log: log})
}
