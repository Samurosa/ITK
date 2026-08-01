package mapper

import (
	"ITK_Code/m/v2/internal/core/dto"
	"errors"

	"github.com/Samurosa/exchange-contract/protobuf/gen/go/shared"
	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"github.com/shopspring/decimal"
)

func WithProtoMoney(protoMoney *shared.Money) (dto.Money, error) {
	if protoMoney == nil {
		return dto.Money{}, errors.New("money is nil")
	}

	amount, err := decimal.NewFromString(protoMoney.Amount)
	if err != nil {
		return dto.Money{}, err
	}

	return dto.Money{
		Currency: protoMoney.Currency,
		Amount:   amount,
	}, nil
}

func ToProtoBalance(balance dto.Balance) *pb.Balance {
	return &pb.Balance{
		Asset:     balance.Asset,
		Available: balance.Available.String(),
		Locked:    balance.Locked.String(),
	}
}

func ToProtoBalances(
	balance []dto.Balance,
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
