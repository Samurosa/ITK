package server

import (
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/mapper"
	"context"
	"errors"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"github.com/Samurosa/exchange-common/shared/auth/sharedContext"
	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/order"
	"github.com/Samurosa/exchange-contract/protobuf/gen/go/shared"
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
		logging.FromContext(ctx).Debug("create order request validation failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	requestOrder, err := mapper.FromProtoCreateOrder(req)
	if err != nil {
		logging.FromContext(ctx).Debug("create order amount parsing failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid order amount")
	}

	tokenContext, ok := sharedContext.GetRequestContext(ctx)
	if !ok {
		logging.FromContext(ctx).Error("authenticated request context missing")
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
		logging.FromContext(ctx).Debug("get order request validation failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	tokenContext, ok := sharedContext.GetRequestContext(ctx)
	if !ok {
		logging.FromContext(ctx).Error("authenticated request context missing")
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
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	log := logging.FromContext(ctx).With(zap.String("order_id", req.GetOrderId()))
	if err := req.Validate(); err != nil {
		log.Debug("stream order request validation failed", zap.Error(err))
		return status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	tokenContext, ok := sharedContext.GetRequestContext(ctx)
	if !ok {
		log.Error("authenticated request context missing")
		return status.Error(codes.Unauthenticated, "request context not provided")
	}
	updates, err := o.order.SubscribeOrderUpdates(ctx, req.OrderId, tokenContext.Principal.UserID)
	if err != nil {
		return mapper.ToGRPC(err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case result, ok := <-updates:
			if !ok {
				if err := ctx.Err(); err != nil {
					return mapper.ToGRPC(err)
				}
				log.Error("order update channel closed unexpectedly")
				return status.Error(codes.Internal, "order subscription stopped unexpectedly")
			}
			if result.Err != nil {
				return mapper.ToGRPC(result.Err)
			}

			update := result.Update
			response := &pb.StreamOrderUpdateResponse{
				OrderId:        update.OrderID,
				OrderStatus:    mapper.ToProtoStatus(update.OrderStatus),
				FilledQuantity: &shared.Money{Currency: update.QuantityCurrency, Amount: update.Quantity},
				UpdatedAt:      timestamppb.New(update.UpdatedAt),
			}

			if err := stream.Send(response); err != nil {
				if ctx.Err() != nil || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) ||
					status.Code(err) == codes.Canceled || status.Code(err) == codes.DeadlineExceeded {
					log.Debug("order update send canceled", zap.Error(err))
				} else {
					log.Warn("order update send failed", zap.Error(err))
				}
				return err
			}
			log.Debug("order update sent", zap.String("order_status", string(update.OrderStatus)))
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
		logging.FromContext(ctx).Debug("list orders request validation failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	tokenContext, ok := sharedContext.GetRequestContext(ctx)
	if !ok {
		logging.FromContext(ctx).Error("authenticated request context missing")
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
