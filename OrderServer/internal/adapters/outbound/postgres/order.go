package postgres

import (
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/order/models"
	"context"
)

func (s *Storage) Save(ctx context.Context, order models.CreateOrder) (string, error) {

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
			quantity
		 )
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (user_id, idempotency_key)
        DO UPDATE SET id = orders.id
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
		order.Quantity,
	).Scan(
		&orderID,
	)
	if err != nil {
		return "", err
	}

	return orderID, err
}

func (s *Storage) Get(ctx context.Context, orderID string) (dto.Order, error) {
	panic("implement me")
}

func (s *Storage) List(ctx context.Context, searchReq models.ListOrdersRequest) ([]dto.Order, string, bool, error) {

	panic("implement me")
}
