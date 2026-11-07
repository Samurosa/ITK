package postgres

import (
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/errorcore"
	"ITK_Code/m/v2/internal/core/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Samurosa/exchange-common/shared/encoding/cursor"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) Save(
	ctx context.Context,
	order models.CreateOrder,
	now time.Time,
) (string, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var orderID string

	insertOrderQuery := `
		INSERT INTO orders (
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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (user_id, idempotency_key) DO NOTHING
		RETURNING id::text
	`

	err = tx.QueryRow(
		ctx,
		insertOrderQuery,
		order.IdempotencyKey,
		order.UserID,
		order.SpotID,
		order.OrderSide,
		dto.StatusNew,
		order.Price,
		order.PriceCurrency,
		order.Quantity,
		order.QuantityCurrency,
		now,
		now,
	).Scan(&orderID)

	if errors.Is(err, pgx.ErrNoRows) {

		const findExistingOrderQuery = `
			SELECT id::text
			FROM orders
			WHERE user_id = $1
				AND idempotency_key = $2
				AND spot_id = $3
				AND order_side = $4
				AND price = $5
				AND price_currency = $6
				AND quantity = $7
				AND quantity_currency = $8
		`

		err = tx.QueryRow(
			ctx,
			findExistingOrderQuery,
			order.UserID,
			order.IdempotencyKey,
			order.SpotID,
			order.OrderSide,
			order.Price,
			order.PriceCurrency,
			order.Quantity,
			order.QuantityCurrency,
		).Scan(&orderID)

		if errors.Is(err, pgx.ErrNoRows) {
			return "", errorcore.ErrIdempotencyConflict
		}
		if err != nil {
			return "", err
		}

		if err = tx.Commit(ctx); err != nil {
			return "", err
		}

		return orderID, nil
	}

	if err != nil {
		return "", err
	}

	command := models.CreateOrderCommand{
		CommandID:        uuid.NewString(),
		OrderID:          orderID,
		UserID:           order.UserID,
		SpotID:           order.SpotID,
		OrderSide:        string(order.OrderSide),
		Price:            order.Price.String(),
		PriceCurrency:    order.PriceCurrency,
		Quantity:         order.Quantity.String(),
		QuantityCurrency: order.QuantityCurrency,
	}

	payload, err := json.Marshal(command)
	if err != nil {
		return "", err
	}

	const insertOutboxQuery = `
		INSERT INTO order_outbox (
		    id,
			order_id,
			message_type,
			payload,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = tx.Exec(
		ctx,
		insertOutboxQuery,
		command.CommandID,
		orderID,
		dto.CreateOrder,
		payload,
		now,
	)
	if err != nil {
		return "", err
	}

	if err = tx.Commit(ctx); err != nil {
		return "", err
	}

	return orderID, nil
}

func (s *Storage) Get(
	ctx context.Context,
	orderID, UserID string,
) (dto.Order, error) {
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

	err := s.pool.QueryRow(
		ctx,
		query,
		orderID,
		UserID,
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

func (s *Storage) List(
	ctx context.Context,
	searchReq models.ListOrdersRequest,
) ([]dto.Order, string, bool, error) {
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
			return nil, "", false, fmt.Errorf(
				"%w: %w",
				errorcore.ErrInvalidCursor,
				err,
			)
		}

		conditions = append(
			conditions,
			fmt.Sprintf("(created_at, id) < ($%d, $%d)", argsPos, argsPos+1),
		)

		args = append(
			args,
			gotCursor.CreatedAt,
			gotCursor.ID,
		)

		argsPos += 2
	}

	if searchReq.SpotID != "" {
		conditions = append(
			conditions,
			fmt.Sprintf("spot_id = $%d", argsPos),
		)

		args = append(args, searchReq.SpotID)
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

	query := baseQuery +
		" WHERE " + strings.Join(conditions, " AND ") +
		fmt.Sprintf(" ORDER BY created_at DESC, id DESC LIMIT $%d", argsPos)

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
