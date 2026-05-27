package mapper

import (
	db "ITK_Code/m/v2/internal/storage"

	pb "github.com/Truncklin/exchange-contract/generated"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProto(
	user db.User,
) *pb.User {
	return &pb.User{
		Id:    user.ID,
		Name:  user.Name,
		Login: user.Login,

		Balances: ToProtoBalances(user.Balances),

		Role: ToProtoRole(user.Role),

		CreatedAt: timestamppb.New(
			user.CreateTime,
		),
		UpdatedAt: timestamppb.New(
			user.UpdateTime,
		),
	}
}

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

func ToProtoRole(role db.Role) pb.Role {
	switch role {
	case db.AdminRole:
		return pb.Role_ADMIN
	default:
		return pb.Role_USER
	}

}
