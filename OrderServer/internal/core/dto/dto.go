package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderStatus string

const (
	StatusUnspecified OrderStatus = "ORDER_STATUS_UNSPECIFIED"
	StatusNew         OrderStatus = "ORDER_STATUS_NEW"
	StatusOpen        OrderStatus = "ORDER_STATUS_OPEN"
	StatusFilled      OrderStatus = "ORDER_STATUS_FILLED"
	StatusCanceled    OrderStatus = "ORDER_STATUS_CANCELED"
	StatusRejected    OrderStatus = "ORDER_STATUS_REJECTED"
)

type OrderSide string

const (
	SideUnspecified OrderSide = "ORDER_SIDE_UNSPECIFIED"
	SideBuy         OrderSide = "ORDER_SIDE_BUY"
	SideSell        OrderSide = "ORDER_SIDE_SELL"
)

type Money struct {
	Currency string
	Amount   decimal.Decimal
}

type Order struct {
	OrderID string

	UserID string
	SpotID string

	OrderSide   OrderSide
	OrderStatus OrderStatus

	Money Money

	Quantity string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateOrder struct {
	OrderID string

	OrderStatus OrderStatus

	Quantity string

	UpdatedAt time.Time
}
