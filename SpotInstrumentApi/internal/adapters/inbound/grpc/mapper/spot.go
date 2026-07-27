package mapper

import (
	"ITK_Code/m/v2/internal/core/dto"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/spot"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromProtoCreateSpot(req *pb.CreateSpotRequest) dto.CreateSpot {
	return dto.CreateSpot{
		BaseAsset:         req.GetBaseAsset(),
		QuoteAsset:        req.GetQuoteAsset(),
		PricePrecision:    req.GetPricePrecision(),
		QuantityPrecision: req.GetQuantityPrecision(),
		MinOrderSize:      req.GetMinOrderSize(),
		MaxOrderSize:      req.GetMaxOrderSize(),
		AllowedRoles:      req.GetAllowedRoles(),
		Name:              req.GetName(),
		Description:       req.GetDescription(),
	}
}

func ToProtoSpot(spot dto.Spot) *pb.GetSpotResponse {
	return &pb.GetSpotResponse{
		Id:                spot.ID,
		Symbol:            spot.Symbol,
		BaseAsset:         spot.BaseAsset,
		QuoteAsset:        spot.QuoteAsset,
		PricePrecision:    spot.PricePrecision,
		QuantityPrecision: spot.QuantityPrecision,
		MinOrderSize:      spot.MinOrderSize,
		MaxOrderSize:      spot.MaxOrderSize,
		AllowedRoles:      spot.AllowedRoles,
		Name:              spot.Name,
		Description:       spot.Description,
		Status:            ToProtoStatus(spot.Status),
		CreatedAt:         timestamppb.New(spot.CreatedAt),
		UpdatedAt:         timestamppb.New(spot.UpdatedAt),
		DisableAt:         timestamppb.New(spot.DeletedAt),
	}
}

func ToProtoMarkets(markets []dto.Market) []*pb.Market {
	protoMarkets := make([]*pb.Market, len(markets))
	for _, market := range markets {
		protoMarkets = append(protoMarkets, &pb.Market{
			SpotId:                 market.SpotID,
			Symbol:                 market.Symbol,
			BaseAsset:              market.BaseAsset,
			QuoteAsset:             market.QuoteAsset,
			Status:                 ToProtoStatus(market.Status),
			LastPrice:              market.LastPrice,
			PriceChange_24H:        market.PriceChange24h,
			PriceChangePercent_24H: market.PriceChangePercent24h,
			UpdatedAt:              timestamppb.New(market.UpdatedAt),
		})
	}

	return protoMarkets
}

func ToProtoDescriptionMarket(market dto.DescriptionMarket) *pb.DescribeMarketResponse {
	return &pb.DescribeMarketResponse{
		BaseAsset:   market.BaseAsset,
		QuoteAsset:  market.QuoteAsset,
		Name:        market.Name,
		Description: market.Description,
	}
}

func ToProtoStatus(status dto.SpotStatus) pb.SpotStatus {
	switch status {
	case dto.UnspecifiedStatus:
		return pb.SpotStatus_SPOT_STATUS_UNSPECIFIED
	case dto.ActiveStatus:
		return pb.SpotStatus_SPOT_STATUS_ACTIVE
	case dto.DisabledStatus:
		return pb.SpotStatus_SPOT_STATUS_DISABLED

	default:
		return pb.SpotStatus_SPOT_STATUS_UNSPECIFIED
	}
}
