package models

import (
	"ITK_Code/m/v2/internal/core/dto"

	"github.com/shopspring/decimal"
)

type CreateOrder struct {
	UserId string

	SpotId string

	OrderSide dto.OrderSide

	IdempotencyKey string

	Price         decimal.Decimal
	PriceCurrency string

	Quantity         string
	QuantityCurrency string
}

type ListOrdersRequest struct {
	UserID   string
	PageSize int32
	Cursor   string
	SpotId   string
	Status   dto.OrderStatus
	Side     dto.OrderSide
}
