package spot

import (
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/spot/models"
	"context"
)

type Repository interface {
	Save(ctx context.Context, spot models.CreateSpot) (string, error)
	Get(ctx context.Context, spotID string) (dto.Spot, error)
	Enable(ctx context.Context, spotID string) error
	Disable(ctx context.Context, spotID string) error
	List(ctx context.Context, searchReq models.ListSpotsRequest) ([]dto.SpotListItem, string, bool, error)
}
