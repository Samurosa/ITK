package inmemory

import (
	"ITK_Code/m/v2/internal/core/wallet"
	"context"
	"errors"
	"sync"

	"github.com/shopspring/decimal"
)

type BalanceKey struct {
	UserID string
	Asset  string
}
type BalanceRepository struct {
	mu sync.RWMutex

	balances map[BalanceKey]*wallet.Balance
}

func NewBalanceStorage() *BalanceRepository {
	return &BalanceRepository{
		balances: make(map[BalanceKey]*wallet.Balance),
	}
}

func (b *BalanceRepository) Create(ctx context.Context, userID string, currency string) (wallet.Balance, error) {
	if userID == "" || currency == "" {
		return wallet.Balance{}, errors.New("invalid params")
	}

	balance := wallet.Balance{
		UserID:    userID,
		Asset:     currency,
		Available: decimal.Zero,
		Locked:    decimal.Zero,
	}
	return balance, nil
}

func (b *BalanceRepository) Get(ctx context.Context, userID string, currency string) (wallet.Balance, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	gotBalance, ok := b.balances[BalanceKey{userID, currency}]
	if !ok {
		return wallet.Balance{}, wallet.ErrBalanceNotFound
	}
	return *gotBalance, nil
}

func (b *BalanceRepository) GetOrCreate(ctx context.Context, userID string, currency string) (wallet.Balance, error) {

	gotBalance, ok := b.balances[BalanceKey{userID, currency}]
	if !ok {
		return b.Create(ctx, userID, currency)
	}
	return *gotBalance, nil
}

func (b *BalanceRepository) Save(ctx context.Context, balance wallet.Balance) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.balances[BalanceKey{balance.UserID, balance.Asset}] = &balance
	return nil
}

func (b *BalanceRepository) GetAll(ctx context.Context, userID string) ([]wallet.Balance, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	balancesByID := make([]wallet.Balance, 0)
	for _, balance := range b.balances {
		if balance.UserID != userID {
			continue
		}
		balancesByID = append(balancesByID, *balance)
	}
	if len(balancesByID) == 0 {
		return balancesByID, wallet.ErrBalanceNotFound
	}

	return balancesByID, nil
}
