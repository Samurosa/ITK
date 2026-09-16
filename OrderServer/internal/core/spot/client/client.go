package client

import (
	"ITK_Code/m/v2/internal/core/dto"
	"context"
)

type SpotProvider interface {
	GetSpot(ctx context.Context, spotID string) (dto.Spot, error)
}
