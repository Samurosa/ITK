package wallet

import (
	"github.com/shopspring/decimal"
)

type Balance struct {
	UserID string

	Asset string

	Available decimal.Decimal

	Locked decimal.Decimal
}

type Money struct {
	Currency string
	Amount   decimal.Decimal
}

func (b *Balance) New(userID string, asset string) *Balance {
	return &Balance{
		UserID:    userID,
		Asset:     asset,
		Available: decimal.Decimal{},
		Locked:    decimal.Decimal{},
	}
}
