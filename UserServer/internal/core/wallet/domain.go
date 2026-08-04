package wallet

import (
	"github.com/shopspring/decimal"
)

type Balance struct {
	ID string

	UserID string
	Asset  string

	Available decimal.Decimal
	Locked    decimal.Decimal
}

type Money struct {
	Currency string
	Amount   decimal.Decimal
}
