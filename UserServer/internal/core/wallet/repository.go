package wallet

import (
	"ITK_Code/m/v2/internal/core/dto"
	"context"
)

type Repository interface {
	Create(ctx context.Context, userID string, currency string) (dto.Balance, error)

	Get(ctx context.Context, userID string, currency string) (dto.Balance, error)

	GetOrCreate(ctx context.Context, userID string, currency string) (dto.Balance, error)

	Save(ctx context.Context, balance dto.Balance) error

	GetAll(ctx context.Context, userID string) ([]dto.Balance, error)
}
