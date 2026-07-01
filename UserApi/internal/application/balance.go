package application

import (
	"ITK_Code/m/v2/internal/core/wallet"
	"context"

	"go.uber.org/zap"
)

func (u *User) Deposit(ctx context.Context,
	id string,
	asset string,
	amount wallet.Money,
) (
	bool,
	wallet.Balance,
	error,
) {
	log := u.log.Named("Deposit")

	balance, err := u.balanceRepository.GetOrCreate(ctx, id, asset)
	if err != nil {
		log.Error("balance not found, error creating new balance", zap.String("id", id), zap.Error(wallet.ErrBalanceNotFound))
		return false, wallet.Balance{}, wallet.ErrCreateNewBalance
	}

	newBalance := balance
	newBalance.Available = balance.Available.Add(amount.Amount)

	err = u.balanceRepository.Save(ctx, newBalance)
	if err != nil {
		log.Error("failed to save balance", zap.Error(err))
		return false, wallet.Balance{}, wallet.ErrSaveBalance
	}

	log.Info("deposit successful", zap.String("id", id))
	return true, newBalance, nil
}

func (u *User) GetBalances(ctx context.Context,
	id string,
) (
	[]wallet.Balance,
	error,
) {
	gotBalances, err := u.balanceRepository.GetAll(ctx, id)
	if err != nil {
		return nil, wallet.ErrBalanceNotFound
	}

	return gotBalances, nil
}
