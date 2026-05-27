package handler

import (
	"ITK_Code/m/v2/internal/mapper"
	db "ITK_Code/m/v2/internal/storage"
	"context"
	"time"

	pb "github.com/Truncklin/exchange-contract/generated"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer

	data *db.OrderData

	spotClient pb.SpotInstrumentServiceClient

	userClient pb.UserServiceClient
}

func NewOrderService(
	spotClient pb.SpotInstrumentServiceClient,
	userClient pb.UserServiceClient,
) *OrderService {
	return &OrderService{
		data:       db.NewOrderStorage(),
		spotClient: spotClient,
		userClient: userClient,
	}
}

func (h *OrderService) CreateOrder(
	ctx context.Context,
	req *pb.CreateOrderRequest,
) (*pb.CreateOrderResponse, error) {

	if req.UserId == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"user_id required",
		)
	}

	if req.SpotId == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"spot_id required",
		)
	}

	spotResp, err := h.spotClient.GetSpot(ctx, &pb.GetSpotRequest{
		Id: req.SpotId,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "spot not found")
	}

	userResp, err := h.userClient.Get(ctx, &pb.GetRequest{
		Id: req.UserId,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	id := uuid.New().String()
	now := time.Now()

	order := db.Order{
		ID:     id,
		UserId: userResp.Id,
		SpotId: spotResp.Spot.Id,

		Type:   mapper.ToStorageOrderType(req.OrderType),
		Status: db.Created,

		Price:    req.Price,
		Quantity: req.Quantity,

		CreatedAt: now,
		UpdatedAt: now,
	}

	h.data.Set(&order)

	return &pb.CreateOrderResponse{
		Order: mapper.ToProto(order),
	}, nil
}

func (h *OrderService) GetOrderStatus(
	ctx context.Context,
	req *pb.GetOrderStatusRequest,
) (*pb.GetOrderStatusResponse, error) {

	order, ok := h.data.Get(req.OrderId)

	if !ok {
		return nil, status.Error(
			codes.NotFound,
			"order not found",
		)
	}

	if order.UserId != req.UserId {

		return nil, status.Error(
			codes.PermissionDenied,
			"access denied",
		)
	}

	return &pb.GetOrderStatusResponse{
		OrderStatus: mapper.ToProtoStatus(order.Status),
	}, nil
}
