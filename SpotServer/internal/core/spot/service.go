package spot

import (
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/spot/models"
	"context"
	"time"
)

type Service interface {
	CreateSpot(ctx context.Context, reqSpot models.CreateSpot) (string, time.Time, error)

	GetSpot(ctx context.Context, spotID string) (dto.Spot, error)

	EnableSpot(ctx context.Context, spotID string) error
	DisableSpot(ctx context.Context, spotID string) error

	ListSpots(ctx context.Context,
		request models.ListSpotsRequest,
	) (
		[]dto.SpotListItem,
		string,
		bool,
		error,
	)
}
