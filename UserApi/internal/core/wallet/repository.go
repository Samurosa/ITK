package wallet

import "context"

type Repository interface {
	Create(ctx context.Context, userID string, currency string) (Balance, error)
	Get(ctx context.Context, userID string, currency string) (Balance, error)
	GetOrCreate(ctx context.Context, userID string, currency string) (Balance, error)
	Save(ctx context.Context, balance Balance) error
	GetAll(ctx context.Context, userID string) ([]Balance, error)
}
