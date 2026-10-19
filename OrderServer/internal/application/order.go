package application

import (
	"ITK_Code/m/v2/internal/application/ports"
	"ITK_Code/m/v2/internal/application/validate"
	"ITK_Code/m/v2/internal/core/dto"
	coreErorrs "ITK_Code/m/v2/internal/core/errorcore"
	"ITK_Code/m/v2/internal/core/order/models"
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var orderTracer = otel.Tracer("order-service/application/order")

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
	orderID string,
	status dto.OrderStatus,
	createdAt time.Time,
	err error,
) {
	ctx, span := orderTracer.Start(ctx, "Order.Create")
	defer func() {
		recordSpanError(span, err)
		span.End()
	}()

	log := logging.FromContext(ctx).Named("order.create").With(
		zap.String("user_id", createOrder.UserId),
		zap.String("spot_id", createOrder.SpotId),
	)

	now := time.Now()

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

	orderID, err = o.repository.Save(
		ctx,
		createOrder,
		now,
	)
	if err != nil {
		logFailure(log, "order save failed", err)

		return "", "", time.Time{}, err
	}

	log.Info("order save completed", zap.String("order_id", orderID))
	return orderID, dto.StatusNew, now, nil
}

func (o *OrderService) Get(ctx context.Context, orderID, userID string) (dto.Order, error) {
	ctx, span := orderTracer.Start(ctx, "Order.Get")
	defer span.End()

	log := logging.FromContext(ctx).Named("order.get").With(zap.String("order_id", orderID), zap.String("user_id", userID))

	order, err := o.repository.Get(
		ctx,
		orderID, userID,
	)
	if err != nil {
		recordSpanError(span, err)
		logFailure(log, "order lookup failed", err)
		return dto.Order{}, err
	}

	log.Debug("order lookup completed")
	return order, nil
}

func (o *OrderService) SubscribeOrderUpdates(ctx context.Context, orderID, userID string) (<-chan models.OrderUpdateResult, error) {
	ctx, span := orderTracer.Start(ctx, "Order.SubscribeOrderUpdates")
	log := logging.FromContext(ctx).Named("order.subscribe").With(zap.String("order_id", orderID), zap.String("user_id", userID))

	listenCtx, cancel := context.WithCancel(ctx)

	notifications, err := o.repository.ListenOrderUpdates(listenCtx, orderID)
	if err != nil {
		cancel()
		recordSpanError(span, err)
		span.End()
		logFailure(log, "order listener initialization failed", err)
		return nil, err
	}

	initial, err := o.repository.Get(listenCtx, orderID, userID)
	if err != nil {
		cancel()
		recordSpanError(span, err)
		span.End()
		logFailure(log, "initial order lookup failed", err)
		return nil, err
	}

	updates := make(chan models.OrderUpdateResult, 1)

	updates <- models.OrderUpdateResult{
		Update: orderUpdate(initial),
	}

	log.Debug("order subscription started")

	go func(last dto.Order) {
		defer span.End()
		defer cancel()
		defer close(updates)
		defer log.Debug("order subscription stopped")

		for {
			select {
			case <-listenCtx.Done():
				return

			case notification, ok := <-notifications:
				if !ok {
					if listenCtx.Err() == nil {
						err := coreErorrs.ErrOrderListenerClosed
						recordSpanError(span, err)
						logFailure(
							log,
							"order listener closed",
							err,
						)

						select {
						case updates <- models.OrderUpdateResult{
							Err: err,
						}:
						case <-listenCtx.Done():
						}
					}

					return
				}

				if notification.Err != nil {
					recordSpanError(span, notification.Err)

					logFailure(
						log,
						"order listener failed",
						notification.Err,
					)

					select {
					case updates <- models.OrderUpdateResult{
						Err: notification.Err,
					}:
					case <-listenCtx.Done():
					}

					return
				}

				current, err := o.repository.Get(listenCtx, orderID, userID)
				if err != nil {
					recordSpanError(span, err)
					logFailure(
						log,
						"order update lookup failed",
						err,
					)

					select {
					case updates <- models.OrderUpdateResult{
						Err: err,
					}:
					case <-listenCtx.Done():
					}

					return
				}

				if current.UpdatedAt.Equal(last.UpdatedAt) &&
					current.FilledQuantity == last.FilledQuantity &&
					current.OrderStatus == last.OrderStatus {
					continue
				}

				select {
				case updates <- models.OrderUpdateResult{
					Update: orderUpdate(current),
				}:
					last = current

				case <-listenCtx.Done():
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
	ctx, span := orderTracer.Start(ctx, "Order.ListOrders")
	defer span.End()

	log := logging.FromContext(ctx).Named("order.list").With(zap.String("user_id", request.UserID))

	orderList, cursor, hasMore, err := o.repository.List(ctx, request)
	if err != nil {
		recordSpanError(span, err)
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

func recordSpanError(span trace.Span, err error) {
	if err == nil ||
		errors.Is(err, context.Canceled) ||
		errors.Is(err, coreErorrs.ErrOrderNotFound) ||
		errors.Is(err, coreErorrs.ErrSpotNotFound) ||
		errors.Is(err, coreErorrs.ErrIdempotencyConflict) ||
		errors.Is(err, coreErorrs.ErrInvalidCursor) ||
		errors.Is(err, coreErorrs.ErrSpotInactive) ||
		errors.Is(err, coreErorrs.ErrRolePermissionDenied) ||
		errors.Is(err, coreErorrs.ErrInvalidOrder) {
		return
	}

	span.RecordError(err)
	span.SetStatus(codes.Error, "application operation failed")
}
