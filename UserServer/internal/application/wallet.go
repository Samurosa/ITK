package application

import (
	"ITK_Code/m/v2/internal/adapters/outbound/requestContext"
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/errors"
	"context"

	"go.uber.org/zap"
)

func (w *Wallet) Deposit(ctx context.Context,
	id string,
	asset string,
	amount dto.Money,
) (
	bool,
	dto.Balance,
	error,
) {
	log := w.log.Named("Deposit")

	if _, err := w.userProvider.Get(ctx, id); err != nil {
		log.Error("failed to get user", zap.String("id", id), zap.Error(err))
		return false, dto.Balance{}, errors.ErrUserNotFound
	}
	log.Info("health check user successful")

	newBalance, err := w.balanceRepository.Deposit(ctx, id, asset, amount)
	if err != nil {
		log.Error("failed to save balance", zap.Error(err))
		return false, dto.Balance{}, errors.ErrSaveBalance
	}
	log.Info("deposit successful", zap.String("id", id))

	return true, newBalance, nil
}

func (w *Wallet) GetBalances(ctx context.Context,
) (
	[]dto.Balance,
	error,
) {
	log := w.log.Named("GetBalances")

	id, err := requestContext.GetUserIDFromContext(ctx)
	if err != nil {
		log.Error("context is not valid", zap.Error(err))
		return []dto.Balance{}, errors.ErrInvalidContext
	}
	log.Info("user id from context", zap.String("id", id))

	gotBalances, err := w.balanceRepository.GetAll(ctx, id)
	if err != nil {
		log.Error("failed to get balances", zap.String("id", id), zap.Error(err))
		return nil, errors.ErrBalanceNotFound
	}
	log.Info("balances retrieved", zap.Int("count elements in wallet", len(gotBalances)))

	return gotBalances, nil
}
