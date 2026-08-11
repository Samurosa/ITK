package wallet

import (
	"context"
)

type Repository interface {
	Deposit(ctx context.Context, userID string, asset string, amount Money, idempotentKey string) (Balance, error)

	GetAll(ctx context.Context, userID string) ([]Balance, error)
}
