package server

import (
	"ITK_Code/m/v2/internal/application"
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/order/models"
	"context"
	"time"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/order"
	"google.golang.org/grpc"
)

type Order interface {
	Create(context.Context, models.CreateOrder, string) (string, dto.OrderStatus, time.Time, error)
	Get(context.Context, string, string) (dto.Order, error)
	SubscribeOrderUpdates(context.Context, string, string) (<-chan models.OrderUpdateResult, error)
	ListOrders(context.Context, models.ListOrdersRequest) ([]dto.Order, string, bool, error)
}

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	order *application.OrderService
}

func NewOrderServer(grpc *grpc.Server, order *application.OrderService) {
	pb.RegisterOrderServiceServer(grpc, &OrderServer{order: order})
}
