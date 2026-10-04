package application

import (
	"ITK_Code/m/v2/internal/application/ports"
	"ITK_Code/m/v2/internal/application/validate"
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
		log.Debug("assets new orders not contains",
			zap.String("orderPriceAssets: ", createOrder.PriceCurrency),
			zap.String("spotQuoteAssets: ", spot.QuoteAsset),
			zap.Error(err),
		)
		return "", "", time.Time{}, coreErorrs.ErrInvalidOrder
	}

	minOrderSize, err := decimal.NewFromString(spot.MinOrderSize)
	if err != nil {
		log.Error("Failed to parse minOrderSize", zap.Error(err))
		return "", "", time.Time{}, coreErorrs.ErrInvalidOrder
	}

	maxOrderSize, err := decimal.NewFromString(spot.MaxOrderSize)
	if err != nil {
		log.Error("Failed to parse maxOrderSize", zap.Error(err))
		return "", "", time.Time{}, coreErorrs.ErrInvalidOrder
	}

	if validate.Price(createOrder.Price, spot.PricePrecision) {
		log.Error("Failed to validate price", zap.String("price", createOrder.Price.String()))
		return "", "", time.Time{}, coreErorrs.ErrInvalidOrder
	}

	if validate.Quantity(createOrder.Quantity, spot.QuantityPrecision, minOrderSize, maxOrderSize) {
		log.Error("Failed to validate quantity", zap.String("quantity", createOrder.Quantity.String()))
		return "", "", time.Time{}, coreErorrs.ErrInvalidOrder
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

func (o *OrderService) SubscribeOrderUpdates(ctx context.Context, orderID, userID string) (<-chan dto.UpdateOrder, error) {
	log := o.log.Named("Subscribe order updates")
	initial, err := o.repository.Get(ctx, orderID, userID)
	if err != nil {
		log.Error("failed to get order", zap.String("orderID", orderID), zap.Error(err))
		return nil, err
	}
	updates := make(chan dto.UpdateOrder, 1)
	updates <- orderUpdate(initial)
	go func(last dto.Order) {
		defer close(updates)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				current, err := o.repository.Get(ctx, orderID, userID)
				if err != nil {
					o.log.Error("failed polling order update", zap.String("orderID", orderID), zap.Error(err))
					return
				}
				if current.UpdatedAt.Equal(last.UpdatedAt) && current.FilledQuantity == last.FilledQuantity && current.OrderStatus == last.OrderStatus {
					continue
				}
				select {
				case updates <- orderUpdate(current):
					last = current
				case <-ctx.Done():
					return
				}
			}
		}
	}(initial)
	return updates, nil
}

func orderUpdate(order dto.Order) dto.UpdateOrder {
	return dto.UpdateOrder{
		OrderID:          order.OrderID,
		OrderStatus:      order.OrderStatus,
		Quantity:         order.FilledQuantity,
		QuantityCurrency: order.QuantityCurrency,
		UpdatedAt:        order.UpdatedAt,
	}
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
	log.Debug("order list query completed", zap.Int("count", len(orderList)), zap.Bool("hasMore", hasMore))

	return orderList, cursor, hasMore, nil
}
