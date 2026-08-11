package mapper

import (
	"ITK_Code/m/v2/internal/core/dto"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/spot"
	userPB "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
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
		AllowedRoles:      FromProtoRoles(req.GetAllowedRoles()),
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
		AllowedRoles:      ToProtoRoles(spot.AllowedRoles),
		Name:              spot.Name,
		Description:       spot.Description,
		Status:            ToProtoStatus(spot.Status),
		CreatedAt:         timestamppb.New(spot.CreatedAt),
		UpdatedAt:         timestamppb.New(spot.UpdatedAt),
		DisableAt:         timestamppb.New(spot.DisabledAt),
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

func ToProtoRoles(roles []dto.Role) []userPB.Role {
	result := make([]userPB.Role, len(roles))
	for _, role := range roles {
		switch role {
		case dto.UnspecifiedRole:
			result = append(result, userPB.Role_ROLE_UNSPECIFIED)
		case dto.UserRole:
			result = append(result, userPB.Role_ROLE_USER)
		case dto.GuestRole:
			result = append(result, userPB.Role_ROLE_GUEST)
		case dto.PremiumRole:
			result = append(result, userPB.Role_ROLE_PREMIUM)
		case dto.AdminRole:
			result = append(result, userPB.Role_ROLE_ADMIN)

		default:
			result = append(result, userPB.Role_ROLE_UNSPECIFIED)
		}
	}
	return result
}

func FromProtoRoles(roles []userPB.Role) []dto.Role {
	result := make([]dto.Role, len(roles))
	for _, role := range roles {
		switch role {
		case userPB.Role_ROLE_UNSPECIFIED:
			result = append(result, dto.UnspecifiedRole)
		case userPB.Role_ROLE_USER:
			result = append(result, dto.UserRole)
		case userPB.Role_ROLE_GUEST:
			result = append(result, dto.GuestRole)
		case userPB.Role_ROLE_PREMIUM:
			result = append(result, dto.PremiumRole)
		case userPB.Role_ROLE_ADMIN:
			result = append(result, dto.AdminRole)

		default:
			result = append(result, dto.UnspecifiedRole)
		}
	}
	return result
}
