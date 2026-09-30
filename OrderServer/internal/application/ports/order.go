package ports

import (
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/order/models"
	"context"
)

type OrderRepository interface {
	Save(context.Context, models.CreateOrder) (string, error)
	Get(context.Context, string, string) (dto.Order, error)
	List(context.Context, models.ListOrdersRequest) ([]dto.Order, string, bool, error)
}

type SpotClient interface {
	GetSpot(context.Context, string) (dto.Spot, error)
}
