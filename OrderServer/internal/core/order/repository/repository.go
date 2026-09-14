package repository

import (
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/order/models"
	"context"
)

type OrderRepository interface {
	Save(ctx context.Context, order models.CreateOrder) (string, error)
	Get(ctx context.Context, orderID string) (dto.Order, error)

	List(ctx context.Context, searchReq models.ListOrdersRequest) ([]dto.Order, string, bool, error)
}
