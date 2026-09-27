package server

import (
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/mapper"
	"context"

	"github.com/Samurosa/exchange-common/shared/auth/sharedContext"
	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/order"
	"github.com/Samurosa/exchange-contract/protobuf/gen/go/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (o *OrderServer) CreateOrder(ctx context.Context,
	req *pb.CreateOrderRequest,
) (
	*pb.CreateOrderResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	requestOrder, err := mapper.FromProtoCreateOrder(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order amount")
	}

	tokenContext, ok := sharedContext.GetRequestContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "request context not provided")
	}

	requestOrder.UserId = tokenContext.Principal.UserID

	orderID, orderStatus, createdTo, err := o.order.Create(ctx, requestOrder, tokenContext.Principal.Role)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &pb.CreateOrderResponse{
		OrderId:     orderID,
		OrderStatus: mapper.ToProtoStatus(orderStatus),
		CreatedAt:   timestamppb.New(createdTo),
	}, nil
}

func (o *OrderServer) GetOrder(ctx context.Context,
	req *pb.GetOrderRequest) (
	*pb.GetOrderResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	tokenContext, ok := sharedContext.GetRequestContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "request context not provided")
	}

	order, err := o.order.Get(ctx, req.OrderId, tokenContext.Principal.UserID)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &pb.GetOrderResponse{
		Order: mapper.ToProtoOrder(order),
	}, nil
}

func (o *OrderServer) StreamOrderUpdate(
	req *pb.StreamOrderUpdateRequest,
	stream pb.OrderService_StreamOrderUpdateServer,
) error {
	if err := req.Validate(); err != nil {
		return status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	ctx := stream.Context()
	tokenContext, ok := sharedContext.GetRequestContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "request context not provided")
	}
	if _, err := o.order.Get(ctx, req.OrderId, tokenContext.Principal.UserID); err != nil {
		return mapper.ToGRPC(err)
	}

	updates, err := o.order.SubscribeOrderUpdates(ctx, req.OrderId)
	if err != nil {
		return mapper.ToGRPC(err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case update, ok := <-updates:
			if !ok {
				return nil
			}

			response := &pb.StreamOrderUpdateResponse{
				OrderId:        update.OrderID,
				OrderStatus:    mapper.ToProtoStatus(update.OrderStatus),
				FilledQuantity: &shared.Money{Currency: update.QuantityCurrency, Amount: update.Quantity},
				UpdatedAt:      timestamppb.New(update.UpdatedAt),
			}

			if err := stream.Send(response); err != nil {
				return err
			}
		}
	}
}

func (o *OrderServer) ListOrders(ctx context.Context,
	req *pb.ListOrdersRequest,
) (
	*pb.ListOrdersResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	tokenContext, ok := sharedContext.GetRequestContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "request context not provided")
	}

	request := mapper.FromProtoListOrdersRequest(req)
	request.UserID = tokenContext.Principal.UserID

	listOrders, cursor, hasMore, err := o.order.ListOrders(ctx, request)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}
	return &pb.ListOrdersResponse{
		Orders:     mapper.ToProtoOrderList(listOrders),
		NextCursor: cursor,
		HasMore:    hasMore,
	}, nil
}
