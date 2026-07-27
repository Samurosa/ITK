package spot

import (
	"ITK_Code/m/v2/internal/core/dto"
)

type Repository interface {
	Save(spot dto.Spot) error
	Disable(spotID string) error
}
