package dto

import (
	"time"
)

type Market struct {
	ID     string
	SpotID string

	Symbol string

	BaseAsset  string
	QuoteAsset string

	Status                SpotStatus
	LastPrice             string
	PriceChange24h        string
	PriceChangePercent24h string

	UpdatedAt time.Time
}

type DescriptionMarket struct {
	BaseAsset  string
	QuoteAsset string

	Name        string
	Description string
}
