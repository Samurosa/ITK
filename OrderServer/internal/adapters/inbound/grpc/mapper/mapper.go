package mapper

import (
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/order/models"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/order"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProtoOrder(dtoOrder dto.Order) *pb.Order {
	return &pb.Order{
		OrderId:     dtoOrder.OrderID,
		UserId:      dtoOrder.UserID,
		SpotId:      dtoOrder.SpotID,
		OrderSide:   ToProtoSide(dtoOrder.OrderSide),
		Price:       dtoOrder.Price.String(),
		Quantity:    dtoOrder.Quantity,
		OrderStatus: ToProtoStatus(dtoOrder.OrderStatus),
		CreatedAt:   timestamppb.New(dtoOrder.CreatedAt),
		UpdatedAt:   timestamppb.New(dtoOrder.UpdatedAt),
	}
}

func ToProtoOrderList(listOrders []dto.Order) []*pb.Order {
	if len(listOrders) == 0 {
		return nil
	}

	result := make([]*pb.Order, 0, len(listOrders))

	for _, spot := range listOrders {
		result = append(result, ToProtoOrder(spot))
	}

	return result
}

func ToProtoStatus(status dto.OrderStatus) pb.OrderStatus {
	switch status {
	case dto.StatusUnspecified:
		return pb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	case dto.StatusNew:
		return pb.OrderStatus_ORDER_STATUS_NEW
	case dto.StatusOpen:
		return pb.OrderStatus_ORDER_STATUS_OPEN
	case dto.StatusFilled:
		return pb.OrderStatus_ORDER_STATUS_FILLED
	case dto.StatusCanceled:
		return pb.OrderStatus_ORDER_STATUS_CANCELED
	case dto.StatusRejected:
		return pb.OrderStatus_ORDER_STATUS_REJECTED

	default:
		return pb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}

func ToProtoSide(side dto.OrderSide) pb.OrderSide {
	switch side {
	case dto.SideUnspecified:
		return pb.OrderSide_ORDER_SIDE_UNSPECIFIED
	case dto.SideBuy:
		return pb.OrderSide_ORDER_SIDE_BUY
	case dto.SideSell:
		return pb.OrderSide_ORDER_SIDE_SELL

	default:
		return pb.OrderSide_ORDER_SIDE_UNSPECIFIED
	}
}

func FromProtoCreateOrder(req *pb.CreateOrderRequest) models.CreateOrder {
	return models.CreateOrder{
		SpotId:         req.SpotId,
		OrderSide:      FromProtoOrderSide(req.OrderSide),
		IdempotencyKey: req.IdempotencyKey,
		Price:          ConvertToDecimal(req.Price),
		Quantity:       req.Quantity,
	}
}

func FromProtoListOrdersRequest(req *pb.ListOrdersRequest) models.ListOrdersRequest {
	result := models.ListOrdersRequest{
		PageSize: req.GetPageSize(),
		Cursor:   req.GetCursor(),
		SpotId:   req.GetSpotId(),
	}

	if req.Status != nil {
		result.Status = FromProtoStatus(req.GetStatus())
	}
	if req.Side != nil {
		result.Side = FromProtoOrderSide(req.GetSide())
	}

	return result
}

func FromProtoOrderSide(orderProto pb.OrderSide) dto.OrderSide {
	switch orderProto {
	case pb.OrderSide_ORDER_SIDE_UNSPECIFIED:
		return dto.SideUnspecified
	case pb.OrderSide_ORDER_SIDE_BUY:
		return dto.SideBuy
	case pb.OrderSide_ORDER_SIDE_SELL:
		return dto.SideSell
	default:
		return dto.SideUnspecified
	}
}

func FromProtoStatus(orderProto pb.OrderStatus) dto.OrderStatus {
	switch orderProto {
	case pb.OrderStatus_ORDER_STATUS_UNSPECIFIED:
		return dto.StatusUnspecified
	case pb.OrderStatus_ORDER_STATUS_NEW:
		return dto.StatusNew
	case pb.OrderStatus_ORDER_STATUS_OPEN:
		return dto.StatusOpen
	case pb.OrderStatus_ORDER_STATUS_FILLED:
		return dto.StatusFilled
	case pb.OrderStatus_ORDER_STATUS_CANCELED:
		return dto.StatusCanceled
	case pb.OrderStatus_ORDER_STATUS_REJECTED:
		return dto.StatusRejected

	default:
		return dto.StatusUnspecified
	}
}

func ConvertToDecimal(amount string) decimal.Decimal {
	if amount == "" {
		return decimal.Zero
	}
	return decimal.RequireFromString(amount)
}
