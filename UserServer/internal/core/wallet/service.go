package wallet

import (
	"ITK_Code/m/v2/internal/core/dto"
	"context"
)

type Service interface {
	Deposit(ctx context.Context,
		id string,
		asset string,
		amount dto.Money,
	) (
		success bool,
		balance dto.Balance,
		err error,
	)

	GetBalances(ctx context.Context,
	) (
		[]dto.Balance,
		error,
	)
}
