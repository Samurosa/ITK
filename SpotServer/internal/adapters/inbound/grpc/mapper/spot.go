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
	response := &pb.GetSpotResponse{
		Id:                spot.ID,
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
	}

	if spot.DisabledAt != nil {
		response.DisableAt = timestamppb.New(*spot.DisabledAt)
	}

	return response
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
	result := make([]userPB.Role, 0, len(roles))

	for _, role := range roles {
		result = append(result, ToProtoRole(role))
	}

	return result
}

func ToProtoRole(role dto.Role) userPB.Role {
	switch role {
	case dto.UnspecifiedRole:
		return userPB.Role_ROLE_UNSPECIFIED
	case dto.UserRole:
		return userPB.Role_ROLE_USER
	case dto.GuestRole:
		return userPB.Role_ROLE_GUEST
	case dto.PremiumRole:
		return userPB.Role_ROLE_PREMIUM
	case dto.AdminRole:
		return userPB.Role_ROLE_ADMIN

	default:
		return userPB.Role_ROLE_UNSPECIFIED
	}
}

func FromProtoRoles(roles []userPB.Role) []dto.Role {
	result := make([]dto.Role, 0, len(roles))

	for _, role := range roles {
		result = append(result, FromProtoRole(role))
	}

	return result
}

func FromProtoRole(role userPB.Role) dto.Role {
	switch role {
	case userPB.Role_ROLE_UNSPECIFIED:
		return dto.UnspecifiedRole
	case userPB.Role_ROLE_USER:
		return dto.UserRole
	case userPB.Role_ROLE_GUEST:
		return dto.GuestRole
	case userPB.Role_ROLE_PREMIUM:
		return dto.PremiumRole
	case userPB.Role_ROLE_ADMIN:
		return dto.AdminRole

	default:
		return dto.UnspecifiedRole
	}
}
