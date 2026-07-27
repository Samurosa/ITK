package spot

import (
	"ITK_Code/m/v2/internal/core/dto"
	"context"
)

type Repository interface {
	Save(ctx context.Context, spot dto.CreateSpot) (string, error)
	Enable(ctx context.Context, spotID string) error
	Disable(ctx context.Context, spotID string) error
}
