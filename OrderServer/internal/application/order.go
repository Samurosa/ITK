package application

import (
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/order/models"
	"context"
	"time"

	"go.uber.org/zap"
)

func (o *OrderService) Create(ctx context.Context,
	createOrder models.CreateOrder,
) (
	string,
	dto.OrderStatus,
	time.Time,
	error,
) {
	log := o.log.Named("Create order")

	createOrder.UserId = ctx.Value("user_id").(string)

	orderID, err := o.repository.Save(
		ctx,
		createOrder,
	)
	if err != nil {
		log.Error("failed to create order", zap.Error(err))

		return "", "", time.Time{}, err
	}

	return orderID, dto.StatusNew, time.Now(), nil
}

func (o *OrderService) Get(ctx context.Context, orderID string) (dto.Order, error) {
	panic("implement me")
}

func (o *OrderService) SubscribeOrderUpdates(ctx context.Context, orderId string) (<-chan dto.UpdateOrder, error) {
	panic("implement me")
}

func (o *OrderService) ListOrders(ctx context.Context,
	request models.ListOrdersRequest,
) (
	[]dto.Order,
	string,
	bool,
	error,
) {
	panic("implement me")
}
