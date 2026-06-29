package user

import (
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/user/models"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProtoBalances(
	balance map[string]models.Balance,
) []*pb.Balance {
	result := make(
		[]*pb.Balance,
		0,
		len(balance),
	)

	for assets, _ := range balance {
		result = append(result,
			&pb.Balance{
				Asset:     assets,
				Available: "", //заглушка
				Locked:    "", //заглушка
			},
		)
	}
	return result
}

func ToProtoBalance(balance models.Balance) *pb.Balance {
	return &pb.Balance{
		Asset:     balance.Asset,
		Available: "", //заглушка
		Locked:    "", //заглушка
	}
}

func ToProtoRole(role models.Role) pb.Role {
	switch role {
	case models.UserRole:
		return pb.Role_ROLE_USER
	case models.GuestRole:
		return pb.Role_ROLE_GUEST
	case models.PremiumRole:
		return pb.Role_ROLE_PREMIUM
	case models.AdminRole:
		return pb.Role_ROLE_ADMIN
	default:
		return pb.Role_ROLE_UNSPECIFIED
	}
}

func ToProtoTokens(tokens auth.TokensModel) *pb.TokenPairResponse {
	return &pb.TokenPairResponse{
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		AccessExpiresAt:  timestamppb.New(tokens.AccessExpiresAt),
		RefreshExpiresAt: timestamppb.New(tokens.RefreshExpiresAt),
	}
}
