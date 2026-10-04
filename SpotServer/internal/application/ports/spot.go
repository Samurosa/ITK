package ports

import (
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/spot/models"
	"context"
	"time"
)

type SpotRepository interface {
	Save(ctx context.Context, spot models.CreateSpot, now time.Time) (string, error)
	Get(context.Context, string) (dto.Spot, error)
	Enable(context.Context, string) error
	Disable(context.Context, string) error
	List(context.Context, models.ListSpotsRequest) ([]dto.SpotListItem, string, bool, error)
}
