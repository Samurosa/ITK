package spot

import (
	"ITK_Code/m/v2/internal/core/dto"
	"context"
	"time"

	"go.uber.org/zap"
)

type Service interface {
	CreateSpot(ctx context.Context, log *zap.Logger, reqSpot dto.CreateSpot) (string, time.Time, error)
	GetSpot(ctx context.Context, log *zap.Logger, spotID string) (dto.Spot, error)
	EnableSpot(ctx context.Context, log *zap.Logger, spotID string) error
	DisableSpot(ctx context.Context, log *zap.Logger, spotID string) error
}
