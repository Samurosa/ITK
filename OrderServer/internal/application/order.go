package application

import (
	"ITK_Code/m/v2/internal/application/ports"
	coreErorrs "ITK_Code/m/v2/internal/core/corerrors"
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/order/models"
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type OrderService struct {
	log *zap.Logger

	repository ports.OrderRepository

	spotClient ports.SpotClient
}

func NewOrderService(log *zap.Logger, repository ports.OrderRepository, spotClient ports.SpotClient) *OrderService {
	return &OrderService{log: log, repository: repository, spotClient: spotClient}
}

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
	spot, err := o.spotClient.GetSpot(ctx, createOrder.SpotId)
	if err != nil {
		return "", "", time.Time{}, err
	}
	if spot.Status != dto.ActiveStatus {
		return "", "", time.Time{}, fmt.Errorf("spot %s is not active", spot.ID)
	}

	if !slices.Contains(spot.AllowedRoles, dto.Role(userRole)) {
		return "", "", time.Time{}, coreErorrs.ErrRolePermissionDenied
	}
	if createOrder.PriceCurrency != spot.QuoteAsset || createOrder.QuantityCurrency != spot.BaseAsset {
		return "", "", time.Time{}, fmt.Errorf("order currencies do not match spot %s", spot.ID)
	}
	quantity, err := decimal.NewFromString(createOrder.Quantity)
	if err != nil || !createOrder.Price.IsPositive() || !quantity.IsPositive() {
		return "", "", time.Time{}, fmt.Errorf("price and quantity must be positive decimal values")
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

func (o *OrderService) Get(ctx context.Context, orderID, userID string) (dto.Order, error) {
	log := o.log.Named("Get order")

	order, err := o.repository.Get(
		ctx,
		orderID, userID,
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
	log := o.log.Named("ListOrders")

	orderList, cursor, hasMore, err := o.repository.List(ctx, request)
	if err != nil {
		log.Error("order list failed", zap.Error(err))
		return []dto.Order{}, "", false, err
	}
	if len(orderList) == 0 {
		log.Debug("order list is empty")
		return []dto.Order{}, "", false, nil
	}
	log.Debug("order list query completed", zap.Int("count", len(orderList)), zap.Bool("hasMore", hasMore))

	return orderList, cursor, hasMore, nil
}
