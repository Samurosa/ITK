package mapper

import (
	db "ITK_Code/m/v2/internal/storage"

	pb "github.com/Truncklin/exchange-contract/generated"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProto(
	spot db.Spot,
) *pb.Spot {
	return &pb.Spot{
		Id:     spot.ID,
		Symbol: spot.Symbol,

		BaseAsset:  spot.Base_asset,
		QuoteAsset: spot.Quote_asset,

		PricePrecision:    spot.Price_precision,
		QuantityPrecision: spot.Quantity_precision,
		MinOrderSize:      spot.Min_order_size,
		MaxOrderSize:      spot.Max_order_size,

		Status:       ToProtoStatus(spot.Status),
		AllowedRoles: spot.AllowedRoles,

		Name:        spot.Name,
		Description: spot.Description,

		CreatedAt: timestamppb.New(spot.Created_at),
		UpdatedAt: timestamppb.New(spot.Updated_at),
		DeletedAt: timestamppb.New(spot.Deleted_at),
	}
}

func ToProtoStatus(
	spotStatus db.SpotStatus,
) pb.SpotStatus {

	switch spotStatus {
	case db.StatusActive:
		return pb.SpotStatus_SPOT_STATUS_ACTIVE
	case db.StatusDisabled:
		return pb.SpotStatus_SPOT_STATUS_DISABLED
	default:
		return pb.SpotStatus_SPOT_STATUS_UNSPECIFIED
	}

}
