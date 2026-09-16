package service

import (
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/order/models"
	"context"
	"time"
)

type Order interface {
	Create(ctx context.Context,
		createOrder models.CreateOrder,
		userRole string,
	) (
		orderID string,
		orderStatus dto.OrderStatus,
		createdTo time.Time,
		err error,
	)

	Get(ctx context.Context, orderID string) (dto.Order, error)

	SubscribeOrderUpdates(ctx context.Context, orderId string) (<-chan dto.UpdateOrder, error)

	ListOrders(ctx context.Context,
		request models.ListOrdersRequest,
	) (
		[]dto.Order,
		string,
		bool,
		error,
	)
}
