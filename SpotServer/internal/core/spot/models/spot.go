package models

import "ITK_Code/m/v2/internal/core/dto"

type CreateSpot struct {
	BaseAsset         string
	QuoteAsset        string
	PricePrecision    int32
	QuantityPrecision int32
	MinOrderSize      string
	MaxOrderSize      string
	AllowedRoles      []dto.Role
	Name              string
	Description       string
}

type ListSpotsRequest struct {
	PageSize int32
	Cursor   string
	Status   dto.SpotStatus

	BaseAsset  string
	QuoteAsset string
}
