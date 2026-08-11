package application

import (
	"ITK_Code/m/v2/internal/core/spot"

	"go.uber.org/zap"
)

type Spot struct {
	log *zap.Logger

	spotRepository spot.Repository
}

func NewSpot(log *zap.Logger, spotRepository spot.Repository) *Spot {
	return &Spot{
		log:            log,
		spotRepository: spotRepository,
	}
}
