package spot

import (
	"ITK_Code/m/v2/internal/core/dto"
	"context"
	"time"
)

type Service interface {
	CreateSpot(ctx context.Context, reqSpot dto.CreateSpot) (string, time.Time, error)
	GetSpot(ctx context.Context, spotID string) (dto.Spot, error)
	EnableSpot(ctx context.Context, spotID string) error
	DisableSpot(ctx context.Context, spotID string) error
}
