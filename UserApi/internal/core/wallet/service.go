package wallet

import "context"

type Service interface {
	Deposit(ctx context.Context,
		id string,
		asset string,
		amount Money,
	) (
		success bool,
		balance Balance,
		err error,
	)

	GetBalances(ctx context.Context,
		id string,
	) (
		[]Balance,
		error,
	)
}
