package application

import (
	"ITK_Code/m/v2/internal/core/market"
	"ITK_Code/m/v2/internal/core/spot"

	"go.uber.org/zap"
)

type Spot struct {
	log *zap.Logger

	spotService spot.Service

	spotRepository spot.Repository
}

type Market struct {
	log *zap.Logger

	marketService market.Service
}

func NewSpot(log *zap.Logger, spotService spot.Service, spotRepository spot.Repository) *Spot {
	return &Spot{
		log:            log,
		spotService:    spotService,
		spotRepository: spotRepository,
	}
}

func NewMarket(log *zap.Logger, marketService market.Service) *Market {
	return &Market{
		log:           log,
		marketService: marketService,
	}
}
