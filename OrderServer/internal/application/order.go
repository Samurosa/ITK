package application

import (
	coreErorrs "ITK_Code/m/v2/internal/core/coreErrors"
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/order/models"
	"context"
	"time"

	"go.uber.org/zap"
)

func (o *OrderService) Create(ctx context.Context,
	createOrder models.CreateOrder,
	userRole string,
) (
	string,
	dto.OrderStatus,
	time.Time,
	error,
) {
	log := o.log.Named("Create order")
	var isAllowed bool

	spot, err := o.spotProvider.GetSpot(ctx, createOrder.SpotId)
	if err != nil {
		return "", "", time.Time{}, err
	}

	for _, role := range spot.AllowedRoles {
		if string(role) == userRole {
			isAllowed = true
		}
	}
	if !isAllowed {
		return "", "", time.Time{}, coreErorrs.ErrRolePermissionDenied
	}

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
	log := o.log.Named("Get order")

	order, err := o.repository.Get(
		ctx,
		orderID,
	)
	if err != nil {
		log.Error("failed to get order", zap.String("orderID", orderID), zap.Error(err))
		return dto.Order{}, err
	}

	return order, nil
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
	log := o.log.Named("Order spot")

	orderList, cursor, hasMore, err := o.repository.List(ctx, request)
	if err != nil {
		log.Error("order list failed", zap.Error(err))
		return []dto.Order{}, "", false, err
	}
	if len(orderList) == 0 {
		log.Debug("spot list is empty")
		return []dto.Order{}, "", false, nil
	}
	log.Info("Slot search completed successfully")

	return orderList, cursor, hasMore, nil
}
