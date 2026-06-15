package mapper

import (
	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"

	models "ITK_Code/m/v2/internal/domain/models"
)

func ToProtoBalances(
	balance map[string]*models.Balance,
) []*pb.Balance {
	result := make(
		[]*pb.Balance,
		0,
		len(balance),
	)

	for assets, b := range balance {
		result = append(result,
			&pb.Balance{
				Asset:     assets,
				Available: b.Available,
				Locked:    b.Locked,
			},
		)
	}
	return result
}

func ToProtoBalance(balance *models.Balance) *pb.Balance {
	return &pb.Balance{
		Asset:     balance.Asset,
		Available: balance.Available,
		Locked:    balance.Locked,
	}
}

func ToProtoRole(role models.Role) pb.Role {
	switch role {
	case models.UserRole:
		return pb.Role_ROLE_USER
	case models.AdminRole:
		return pb.Role_ROLE_ADMIN
	default:
		return pb.Role_ROLE_UNSPECIFIED
	}

}
