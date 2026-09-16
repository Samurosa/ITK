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

	Price decimal.Decimal

	Quantity string
}

type ListOrdersRequest struct {
	PageSize int32
	Cursor   string
	SpotId   string
	Status   dto.OrderStatus
	Side     dto.OrderSide
}
