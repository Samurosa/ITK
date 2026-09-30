package ports

import (
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/spot/models"
	"context"
)

type SpotRepository interface {
	Save(context.Context, models.CreateSpot) (string, error)
	Get(context.Context, string) (dto.Spot, error)
	Enable(context.Context, string) error
	Disable(context.Context, string) error
	List(context.Context, models.ListSpotsRequest) ([]dto.SpotListItem, string, bool, error)
}
