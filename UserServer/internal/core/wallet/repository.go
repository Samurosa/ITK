package wallet

import (
	"ITK_Code/m/v2/internal/core/dto"
	"context"
)

type Repository interface {
	Deposit(ctx context.Context, userID string, asset string, amount dto.Money) (dto.Balance, error)

	GetAll(ctx context.Context, userID string) ([]dto.Balance, error)
}
