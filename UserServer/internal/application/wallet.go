package application

import (
	"ITK_Code/m/v2/internal/core/wallet"
	"context"

	"go.uber.org/zap"
)

func (w *Wallet) Deposit(ctx context.Context,
	id string,
	asset string,
	amount wallet.Money,
	idempotentKey string,
) (
	wallet.Balance,
	error,
) {
	log := w.log.Named("Deposit")

	newBalance, err := w.balanceRepository.Deposit(ctx, id, asset, amount, idempotentKey)
	if err != nil {
		log.Error("failed to save balance", zap.Error(err))
		return wallet.Balance{}, wallet.ErrSaveBalance
	}
	log.Info("deposit successful", zap.String("balance id", newBalance.ID))

	return newBalance, nil
}

func (w *Wallet) GetBalances(ctx context.Context,
	id string,
) (
	[]wallet.Balance,
	error,
) {
	log := w.log.Named("GetBalances")

	gotBalances, err := w.balanceRepository.GetAll(ctx, id)
	if err != nil {
		log.Error("failed to get balances", zap.String("id", id), zap.Error(err))
		return nil, wallet.ErrBalanceNotFound
	}
	log.Info("balances retrieved", zap.Int("count elements in wallet", len(gotBalances)))

	return gotBalances, nil
}
