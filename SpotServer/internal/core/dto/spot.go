package dto

import "time"

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

type CreateSpot struct {
	BaseAsset         string
	QuoteAsset        string
	PricePrecision    int32
	QuantityPrecision int32
	MinOrderSize      string
	MaxOrderSize      string
	AllowedRoles      []Role
	Name              string
	Description       string
}
