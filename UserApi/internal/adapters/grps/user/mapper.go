package user

import (
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/user"
	"ITK_Code/m/v2/internal/core/wallet"
	"errors"

	"github.com/Samurosa/exchange-contract/protobuf/gen/go/shared"
	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProtoBalances(
	balance []wallet.Balance,
) []*pb.Balance {
	result := make(
		[]*pb.Balance,
		0,
		len(balance),
	)

	for _, b := range balance {
		result = append(result,
			&pb.Balance{
				Asset:     b.Asset,
				Available: b.Available.String(),
				Locked:    b.Locked.String(),
			},
		)
	}
	return result
}

func ToProtoBalance(balance wallet.Balance) *pb.Balance {
	return &pb.Balance{
		Asset:     balance.Asset,
		Available: balance.Available.String(),
		Locked:    balance.Locked.String(),
	}
}

func ToProtoRole(role user.Role) pb.Role {
	switch role {
	case user.UserRole:
		return pb.Role_ROLE_USER
	case user.GuestRole:
		return pb.Role_ROLE_GUEST
	case user.PremiumRole:
		return pb.Role_ROLE_PREMIUM
	case user.AdminRole:
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

func ToProtoMoney(protoMoney *shared.Money) (wallet.Money, error) {
	if protoMoney == nil {
		return wallet.Money{}, errors.New("money is nil")
	}

	amount, err := decimal.NewFromString(protoMoney.Amount)
	if err != nil {
		return wallet.Money{}, err
	}

	return wallet.Money{
		Currency: protoMoney.Currency,
		Amount:   amount,
	}, nil
}
