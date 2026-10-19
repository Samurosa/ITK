package postgres

import (
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/errorcore"
	"ITK_Code/m/v2/internal/core/order/models"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Samurosa/exchange-common/shared/encoding/cursor"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) Save(ctx context.Context, order models.CreateOrder, now time.Time) (string, error) {

	var orderID string

	query := `
		INSERT INTO orders
		(
			idempotency_key,
			user_id,
			spot_id,
			order_side,
			order_status,
			price,
			price_currency,
			quantity,
			quantity_currency,
		 created_at,
		 updated_at
		 )
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (user_id, idempotency_key)
		DO UPDATE SET idempotency_key = orders.idempotency_key
		WHERE orders.spot_id = EXCLUDED.spot_id
		  AND orders.order_side = EXCLUDED.order_side
		  AND orders.price = EXCLUDED.price
		  AND orders.price_currency = EXCLUDED.price_currency
		  AND orders.quantity = EXCLUDED.quantity
		  AND orders.quantity_currency = EXCLUDED.quantity_currency
		RETURNING id
	`

	err := s.pool.QueryRow(ctx,
		query,
		order.IdempotencyKey,
		order.UserId,
		order.SpotId,
		order.OrderSide,
		dto.StatusNew,
		order.Price,
		order.PriceCurrency,
		order.Quantity,
		order.QuantityCurrency,
		now,
		now,
	).Scan(
		&orderID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", errorcore.ErrIdempotencyConflict
	}
	if err != nil {
		return "", err
	}

	return orderID, err
}

func (s *Storage) Get(ctx context.Context, orderID, userID string) (dto.Order, error) {
	query := `
		SELECT
			id,
			user_id,
			spot_id,
			order_side,
			order_status,
			price,
			price_currency,
			quantity,
			quantity_currency,
			filled_quantity,
			created_at,
			updated_at
		FROM orders
		WHERE id = $1 AND user_id = $2
`

	var order dto.Order

	err := s.pool.QueryRow(ctx,
		query,
		orderID,
		userID,
	).Scan(
		&order.OrderID,
		&order.UserID,
		&order.SpotID,
		&order.OrderSide,
		&order.OrderStatus,
		&order.Price,
		&order.PriceCurrency,
		&order.Quantity,
		&order.QuantityCurrency,
		&order.FilledQuantity,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return dto.Order{}, errorcore.ErrOrderNotFound
	}
	if err != nil {
		return dto.Order{}, err
	}
	return order, nil
}

func (s *Storage) List(ctx context.Context, searchReq models.ListOrdersRequest) ([]dto.Order, string, bool, error) {

	baseQuery := `
		SELECT
			id,
			user_id,
			spot_id,
			order_side,
			order_status,
			price,
			price_currency,
			quantity,
			quantity_currency,
			filled_quantity,
			created_at,
			updated_at
		FROM orders
`
	args := []any{searchReq.UserID}
	argsPos := 2
	conditions := []string{"user_id = $1"}

	if searchReq.Cursor != "" {
		gotCursor, err := cursor.DecodeCursor(searchReq.Cursor)
		if err != nil {
			return nil, "", false, fmt.Errorf("%w: %w", errorcore.ErrInvalidCursor, err)
		}

		conditions = append(
			conditions,
			fmt.Sprintf(
				"(created_at, id) < ($%d, $%d)",
				argsPos,
				argsPos+1,
			),
		)

		args = append(
			args,
			gotCursor.CreatedAt,
			gotCursor.ID,
		)

		argsPos += 2
	}

	if searchReq.SpotId != "" {
		conditions = append(
			conditions,
			fmt.Sprintf("spot_id = $%d", argsPos),
		)

		args = append(args, searchReq.SpotId)

		argsPos++
	}

	if searchReq.Status != "" {
		conditions = append(
			conditions,
			fmt.Sprintf("order_status = $%d", argsPos),
		)

		args = append(args, searchReq.Status)

		argsPos++
	}

	if searchReq.Side != "" {
		conditions = append(
			conditions,
			fmt.Sprintf("order_side = $%d", argsPos),
		)

		args = append(args, searchReq.Side)

		argsPos++
	}

	query := baseQuery + " WHERE " + strings.Join(conditions, " AND ")

	query += " ORDER BY created_at DESC, id DESC " + fmt.Sprintf(" LIMIT $%d ", argsPos)

	pageSize := int(searchReq.PageSize)
	limit := pageSize + 1

	args = append(args, limit)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, "", false, err
	}
	defer rows.Close()

	orders := make([]dto.Order, 0, limit)

	for rows.Next() {

		var order dto.Order

		err = rows.Scan(
			&order.OrderID,
			&order.UserID,
			&order.SpotID,
			&order.OrderSide,
			&order.OrderStatus,
			&order.Price,
			&order.PriceCurrency,
			&order.Quantity,
			&order.QuantityCurrency,
			&order.FilledQuantity,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, "", false, err
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, "", false, err
	}

	hasMore := len(orders) > pageSize

	if hasMore {
		orders = orders[:pageSize]

		newCursor, err := cursor.EncodeCursor(
			cursor.Cursor{
				CreatedAt: orders[len(orders)-1].CreatedAt,
				ID:        orders[len(orders)-1].OrderID,
			},
		)
		if err != nil {
			return nil, "", false, err
		}

		return orders, newCursor, true, nil
	}

	return orders, "", false, nil
}

func (s *Storage) ListenOrderUpdates(ctx context.Context, orderID string) (<-chan models.OrderNotification, error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"acquire order listener connection: %w",
			err,
		)
	}

	_, err = conn.Exec(ctx, "LISTEN order_updates")
	if err != nil {
		conn.Release()

		return nil, fmt.Errorf(
			"listen order updates: %w",
			err,
		)
	}

	notifications := make(chan models.OrderNotification, 1)

	go func() {
		defer close(notifications)

		defer conn.Release()

		for {
			notification, err := conn.Conn().
				WaitForNotification(ctx)

			if err != nil {
				if ctx.Err() == nil {
					select {
					case notifications <- models.OrderNotification{
						Err: fmt.Errorf(
							"wait for order notification: %w",
							err,
						),
					}:
					case <-ctx.Done():
					}
				}

				return
			}

			if notification.Payload != orderID {
				continue
			}

			select {
			case notifications <- models.OrderNotification{}:
			default:
			}
		}
	}()

	return notifications, nil
}
