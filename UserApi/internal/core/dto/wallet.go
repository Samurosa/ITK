package dto

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
