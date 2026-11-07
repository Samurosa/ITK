package models

import (
	"ITK_Code/m/v2/internal/core/dto"

	"github.com/shopspring/decimal"
)

type CreateOrder struct {
	UserID         string
	SpotID         string
	OrderSide      dto.OrderSide
	IdempotencyKey string

	Price         decimal.Decimal
	PriceCurrency string

	Quantity         decimal.Decimal
	QuantityCurrency string
}

type CreateOrderCommand struct {
	CommandID        string `json:"command_id"`
	OrderID          string `json:"order_id"`
	UserID           string `json:"user_id"`
	SpotID           string `json:"spot_id"`
	OrderSide        string `json:"order_side"`
	Price            string `json:"price"`
	PriceCurrency    string `json:"price_currency"`
	Quantity         string `json:"quantity"`
	QuantityCurrency string `json:"quantity_currency"`
}

type ListOrdersRequest struct {
	UserID   string
	PageSize int32
	Cursor   string
	SpotID   string
	Status   dto.OrderStatus
	Side     dto.OrderSide
}

type OrderUpdateResult struct {
	Update dto.UpdateOrder
	Err    error
}
