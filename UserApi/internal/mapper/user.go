package mapper

import (
	db "ITK_Code/m/v2/internal/storage"

	pb "github.com/Truncklin/exchange-contract/protobuf/gen/go/user"
)

func ToProtoBalances(
	balance map[string]*db.Balance,
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

func ToProtoBalance(balance *db.Balance) *pb.Balance {
	return &pb.Balance{
		Asset:     balance.Asset,
		Available: balance.Available,
		Locked:    balance.Locked,
	}
}

func ToProtoRole(role db.Role) pb.Role {
	switch role {
	case db.UserRole:
		return pb.Role_ROLE_USER
	case db.AdminRole:
		return pb.Role_ROLE_ADMIN
	default:
		return pb.Role_ROLE_UNSPECIFIED
	}

}
