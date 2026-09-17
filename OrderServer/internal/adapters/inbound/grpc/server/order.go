package server

import (
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/mapper"
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/validate"
	"context"

	"github.com/Samurosa/exchange-common/shared/auth/sharedContext"
	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/order"
	"go.uber.org/zap"
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

	requestOrder := mapper.FromProtoCreateOrder(req)

	tokenContext, err := sharedContext.GetRequestContext(ctx)
	if err != nil {
		return nil, mapper.ToGRPC(err)
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

	order, err := o.order.Get(ctx, req.OrderId)
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
				OrderId:     update.OrderID,
				OrderStatus: mapper.ToProtoStatus(update.OrderStatus),
				Quantity:    update.Quantity,
				UpdatedAt:   timestamppb.New(update.UpdatedAt),
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
	log := o.log.Named("ListOrders")
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	tokenContext, err := sharedContext.GetRequestContext(ctx)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	if validate.AccessCheck(tokenContext.Principal.Role) {
		log.Info("attempt to gain access with a ", zap.String("role: ", tokenContext.Principal.Role))
		return nil, status.Error(codes.PermissionDenied, "insufficient privileges")
	}

	request := mapper.FromProtoListOrdersRequest(req)

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
