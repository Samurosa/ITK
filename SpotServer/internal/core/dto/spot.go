package dto

import "time"

type SpotStatus string

const (
	UnspecifiedStatus SpotStatus = "SPOT_STATUS_UNSPECIFIED"
	ActiveStatus                 = "SPOT_STATUS_ACTIVE"
	DisabledStatus               = "SPOT_STATUS_DISABLED"
)

type Spot struct {
	ID string

	Symbol     string //USD/USDT
	BaseAsset  string //USD
	QuoteAsset string //USDT

	PricePrecision    int32 // количество знаков после запятой у актива
	QuantityPrecision int32 // количество знаков после запятой у цены
	MinOrderSize      string
	MaxOrderSize      string

	AllowedRoles []string

	Name        string
	Description string

	Status SpotStatus

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type CreateSpot struct {
	Symbol            string
	BaseAsset         string
	QuoteAsset        string
	PricePrecision    int32
	QuantityPrecision int32
	MinOrderSize      string
	MaxOrderSize      string
	AllowedRoles      []string
	Name              string
	Description       string
}
