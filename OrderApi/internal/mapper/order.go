package mapper

import (
	db "ITK_Code/m/v2/internal/storage"

	pb "github.com/Samurosa/exchange-contract/generated"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProto(
	order db.Order,
) *pb.Order {
	return &pb.Order{
		Id:     order.ID,
		UserId: order.UserId,
		SpotId: order.SpotId,

		OrderType:   ToOrderStorageType(order.Type),
		OrderStatus: ToProtoStatus(order.Status),

		Price:    order.Price,
		Quantity: order.Quantity,

		CreatedAt: timestamppb.New(order.CreatedAt),
		UpdatedAt: timestamppb.New(order.UpdatedAt),
	}
}

func ToOrderStorageType(
	orderType db.OrderTypeStorage,
) pb.OrderType {
	switch orderType {
	case db.Buy:
		return pb.OrderType_BUY
	case db.Sell:
		return pb.OrderType_SELL
	default:
		return pb.OrderType_ORDER_TYPE_UNSPECIFIED
	}

}

func ToStorageOrderType(
	t pb.OrderType,
) db.OrderTypeStorage {

	switch t {

	case pb.OrderType_BUY:
		return db.Buy

	case pb.OrderType_SELL:
		return db.Sell

	default:
		return db.OrderTypeUnspecified
	}
}

func ToProtoStatus(
	orderStatus db.OrderStatusStorage,
) pb.OrderStatus {
	switch orderStatus {
	case db.Created:
		return pb.OrderStatus_CREATED
	case db.Canceled:
		return pb.OrderStatus_CANCELED
	case db.Rejected:
		return pb.OrderStatus_REJECTED
	case db.Completed:
		return pb.OrderStatus_COMPLETED
	default:
		return pb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}
