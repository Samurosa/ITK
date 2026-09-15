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

type Order struct {
	OrderID string

	UserID string
	SpotID string

	OrderSide   OrderSide
	OrderStatus OrderStatus

	Price decimal.Decimal

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

type SpotStatus string

const (
	UnspecifiedStatus SpotStatus = "SPOT_STATUS_UNSPECIFIED"
	ActiveStatus      SpotStatus = "SPOT_STATUS_ACTIVE"
	DisabledStatus    SpotStatus = "SPOT_STATUS_DISABLED"
)

type Role string

const (
	UnspecifiedRole Role = "ROLE_UNSPECIFIED"

	UserRole Role = "ROLE_USER"

	GuestRole Role = "ROLE_GUEST"

	PremiumRole Role = "ROLE_PREMIUM"

	AdminRole Role = "ROLE_ADMIN"
)

type Spot struct {
	ID string

	BaseAsset  string //USD
	QuoteAsset string //USDT

	PricePrecision    int32 // количество знаков после запятой у актива
	QuantityPrecision int32 // количество знаков после запятой у цены
	MinOrderSize      string
	MaxOrderSize      string

	AllowedRoles []Role

	Name        string
	Description string

	Status SpotStatus

	CreatedAt  time.Time
	UpdatedAt  time.Time
	DisabledAt *time.Time
}
