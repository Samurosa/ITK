package mapper

import (
	"ITK_Code/m/v2/internal/core/wallet"
	"errors"

	"github.com/Samurosa/exchange-contract/protobuf/gen/go/shared"
	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"github.com/shopspring/decimal"
)

func WithProtoMoney(protoMoney *shared.Money) (wallet.Money, error) {
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

func ToProtoBalance(balance wallet.Balance) *pb.Balance {
	return &pb.Balance{
		Asset:     balance.Asset,
		Available: balance.Available.String(),
		Locked:    balance.Locked.String(),
	}
}

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
