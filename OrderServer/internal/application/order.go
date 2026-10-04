package application

import (
	"ITK_Code/m/v2/internal/application/ports"
	"ITK_Code/m/v2/internal/application/validate"
	coreErorrs "ITK_Code/m/v2/internal/core/corerrors"
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/order/models"
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type OrderService struct {
	repository ports.OrderRepository

	spotClient ports.SpotClient
}

func NewOrderService(repository ports.OrderRepository, spotClient ports.SpotClient) *OrderService {
	return &OrderService{repository: repository, spotClient: spotClient}
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
	log := logging.FromContext(ctx).Named("order.create").With(
		zap.String("user_id", createOrder.UserId),
		zap.String("spot_id", createOrder.SpotId),
	)
	spot, err := o.spotClient.GetSpot(ctx, createOrder.SpotId)
	if err != nil {
		logFailure(log, "spot lookup failed", err)
		return "", "", time.Time{}, err
	}
	if spot.Status != dto.ActiveStatus {
		log.Debug("order rejected: spot is not active", zap.String("spot_status", string(spot.Status)))
		return "", "", time.Time{}, coreErorrs.ErrSpotInactive
	}

	if !slices.Contains(spot.AllowedRoles, dto.Role(userRole)) {
		log.Warn("order rejected: role is not allowed", zap.String("role", userRole))
		return "", "", time.Time{}, coreErorrs.ErrRolePermissionDenied
	}
	if createOrder.PriceCurrency != spot.QuoteAsset || createOrder.QuantityCurrency != spot.BaseAsset {
		log.Debug("order rejected: currencies do not match spot assets",
			zap.String("price_currency", createOrder.PriceCurrency),
			zap.String("quote_asset", spot.QuoteAsset),
			zap.String("quantity_currency", createOrder.QuantityCurrency),
			zap.String("base_asset", spot.BaseAsset),
		)
		return "", "", time.Time{}, coreErorrs.ErrInvalidOrder
	}

	minOrderSize, err := decimal.NewFromString(spot.MinOrderSize)
	if err != nil {
		log.Error("invalid spot minimum order size", zap.String("min_order_size", spot.MinOrderSize), zap.Error(err))
		return "", "", time.Time{}, fmt.Errorf("invalid spot minimum order size: %w", err)
	}

	maxOrderSize, err := decimal.NewFromString(spot.MaxOrderSize)
	if err != nil {
		log.Error("invalid spot maximum order size", zap.String("max_order_size", spot.MaxOrderSize), zap.Error(err))
		return "", "", time.Time{}, fmt.Errorf("invalid spot maximum order size: %w", err)
	}

	if validate.Price(createOrder.Price, spot.PricePrecision) {
		log.Debug("order rejected: price violates spot constraints", zap.String("price", createOrder.Price.String()), zap.Int32("price_precision", spot.PricePrecision))
		return "", "", time.Time{}, coreErorrs.ErrInvalidOrder
	}

	if validate.Quantity(createOrder.Quantity, spot.QuantityPrecision, minOrderSize, maxOrderSize) {
		log.Debug("order rejected: quantity violates spot constraints",
			zap.String("quantity", createOrder.Quantity.String()),
			zap.Int32("quantity_precision", spot.QuantityPrecision),
			zap.String("min_order_size", spot.MinOrderSize),
			zap.String("max_order_size", spot.MaxOrderSize),
		)
		return "", "", time.Time{}, coreErorrs.ErrInvalidOrder
	}

	orderID, err := o.repository.Save(
		ctx,
		createOrder,
	)
	if err != nil {
		logFailure(log, "order save failed", err)

		return "", "", time.Time{}, err
	}

	// Save may return an existing order for an idempotent retry.
	log.Info("order save completed", zap.String("order_id", orderID))
	return orderID, dto.StatusNew, time.Now(), nil
}

func (o *OrderService) Get(ctx context.Context, orderID, userID string) (dto.Order, error) {
	log := logging.FromContext(ctx).Named("order.get").With(zap.String("order_id", orderID), zap.String("user_id", userID))

	order, err := o.repository.Get(
		ctx,
		orderID, userID,
	)
	if err != nil {
		logFailure(log, "order lookup failed", err)
		return dto.Order{}, err
	}

	log.Debug("order lookup completed")
	return order, nil
}

// OrderUpdateResult carries a polling failure to the stream handler so that a
// failed subscription cannot be reported to the client as a successful stream.
type OrderUpdateResult struct {
	Update dto.UpdateOrder
	Err    error
}

func (o *OrderService) SubscribeOrderUpdates(ctx context.Context, orderID, userID string) (<-chan OrderUpdateResult, error) {
	log := logging.FromContext(ctx).Named("order.subscribe").With(zap.String("order_id", orderID), zap.String("user_id", userID))
	initial, err := o.repository.Get(ctx, orderID, userID)
	if err != nil {
		logFailure(log, "initial order lookup failed", err)
		return nil, err
	}
	updates := make(chan OrderUpdateResult, 1)
	updates <- OrderUpdateResult{Update: orderUpdate(initial)}
	log.Debug("order subscription started")
	go func(last dto.Order) {
		defer close(updates)
		defer log.Debug("order subscription stopped")
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				current, err := o.repository.Get(ctx, orderID, userID)
				if err != nil {
					logFailure(log, "order update polling failed", err)
					select {
					case updates <- OrderUpdateResult{Err: err}:
					case <-ctx.Done():
					}
					return
				}
				if current.UpdatedAt.Equal(last.UpdatedAt) && current.FilledQuantity == last.FilledQuantity && current.OrderStatus == last.OrderStatus {
					continue
				}
				select {
				case updates <- OrderUpdateResult{Update: orderUpdate(current)}:
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
	log := logging.FromContext(ctx).Named("order.list").With(zap.String("user_id", request.UserID))

	orderList, cursor, hasMore, err := o.repository.List(ctx, request)
	if err != nil {
		logFailure(log, "order list failed", err)
		return []dto.Order{}, "", false, err
	}
	log.Debug("order list query completed", zap.Int("count", len(orderList)), zap.Bool("has_more", hasMore))

	return orderList, cursor, hasMore, nil
}

func logFailure(log *zap.Logger, message string, err error) {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, coreErorrs.ErrOrderNotFound) || errors.Is(err, coreErorrs.ErrSpotNotFound) ||
		errors.Is(err, coreErorrs.ErrIdempotencyConflict) || errors.Is(err, coreErorrs.ErrInvalidCursor) {
		log.Debug(message, zap.Error(err))
		return
	}
	log.Error(message, zap.Error(err))
}
