package handler

import (
	mapper "ITK_Code/m/v2/internal/mapper"
	db "ITK_Code/m/v2/internal/storage"
	"context"
	"time"

	pb "github.com/Samurosa/exchange-contract/generated"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MarketHandler struct {
	pb.UnimplementedSpotInstrumentServiceServer

	data *db.Market
}

func NewSpotService() *MarketHandler {
	return &MarketHandler{
		data: db.NewMarketStorage(),
	}
}

func (h *MarketHandler) CreateSpot(
	ctx context.Context,
	req *pb.SetSpotRequest,
) (
	*pb.SetSpotResponse,
	error,
) {

	id := uuid.New().String()
	now := time.Now()
	spot := db.Spot{
		ID:     id,
		Symbol: req.Symbol,

		Base_asset:  req.BaseAsset,
		Quote_asset: req.QuoteAsset,

		Price_precision:    req.PricePrecision,
		Quantity_precision: req.QuantityPrecision,
		Min_order_size:     req.MinOrderSize,
		Max_order_size:     req.MaxOrderSize,

		Status:       db.StatusActive,
		AllowedRoles: req.AllowedRoles,

		Name:        req.Name,
		Description: req.Description,

		Created_at: now,
		Updated_at: now,
	}

	h.data.Create(&spot)

	protoSpot := mapper.ToProto(spot)

	return &pb.SetSpotResponse{
		Spot: protoSpot,
	}, nil
}

func (h *MarketHandler) GetSpot(
	ctx context.Context,
	req *pb.GetSpotRequest,
) (
	*pb.GetSpotResponse,
	error,
) {

	spot, ok := h.data.Get(req.Id)
	if !ok {
		return nil, status.Error(
			codes.NotFound,
			"spot not found",
		)
	}

	protoSpot := mapper.ToProto(*spot)

	return &pb.GetSpotResponse{
		Spot: protoSpot,
	}, nil
}

func (h *MarketHandler) DeleteSpot(
	ctx context.Context,
	req *pb.DeleteSpotRequest,
) (
	*pb.DeleteSpotResponse,
	error,
) {

	success, deleted_at := h.data.Delete(req.Id)
	if !success {
		return nil, status.Error(
			codes.NotFound,
			"spot not found",
		)
	}

	return &pb.DeleteSpotResponse{
		Success:   success,
		DeletedAt: timestamppb.New(deleted_at),
	}, nil
}

func (h *MarketHandler) ViewMarkets(
	ctx context.Context,
	req *pb.ViewMarketsRequest,
) (
	*pb.ViewMarketsResponse,
	error,
) {

	spots := h.data.List()

	protoSpots := make([]*pb.Spot, 0, len(spots))

	for _, spot := range spots {

		if spot.Status != db.StatusActive {
			continue
		}

		if !spot.Deleted_at.IsZero() {
			continue
		}

		if !HasAccess(req.UserRoles, spot.AllowedRoles) {
			continue
		}

		protoSpots = append(protoSpots, mapper.ToProto(*spot))
	}

	return &pb.ViewMarketsResponse{
		Spots: protoSpots,
	}, nil
}

func HasAccess(
	userRoles []string,
	spotRoles []string,
) bool {

	roleSet := make(map[string]struct{}, len(spotRoles))

	for _, r := range spotRoles {
		roleSet[r] = struct{}{}
	}

	for _, ur := range userRoles {
		if _, ok := roleSet[ur]; ok {
			return true
		}
	}

	return false
}
