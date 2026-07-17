package application

import (
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/wallet"
	"context"

	"go.uber.org/zap"
)

func (w *Wallet) Deposit(ctx context.Context,
	id string,
	asset string,
	amount wallet.Money,
) (
	bool,
	wallet.Balance,
	error,
) {
	log := w.log.Named("Deposit")
	balance, err := w.balanceRepository.GetOrCreate(ctx, id, asset)
	if err != nil {
		log.Error("balance not found, error creating new balance", zap.String("id", id), zap.Error(err))
		return false, wallet.Balance{}, wallet.ErrCreateNewBalance
	}
	log.Info("balance created", zap.String("id", id))

	newBalance := balance
	newBalance.Available = balance.Available.Add(amount.Amount)

	err = w.balanceRepository.Save(ctx, newBalance)
	if err != nil {
		log.Error("failed to save balance", zap.Error(err))
		return false, wallet.Balance{}, wallet.ErrSaveBalance
	}
	log.Info("deposit successful", zap.String("id", id))

	return true, newBalance, nil
}

func (w *Wallet) GetBalances(ctx context.Context,
) (
	[]wallet.Balance,
	error,
) {
	log := w.log.Named("GetBalances")

	id, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		log.Error("context is not valid", zap.Error(err))
		return []wallet.Balance{}, auth.ErrInvalidContext
	}
	log.Info("user id from context", zap.String("id", id))

	gotBalances, err := w.balanceRepository.GetAll(ctx, id)
	if err != nil {
		log.Error("failed to get balances", zap.String("id", id), zap.Error(err))
		return nil, wallet.ErrBalanceNotFound
	}
	log.Info("balances retrieved", zap.Int("count elements in wallet", len(gotBalances)))

	return gotBalances, nil
}
