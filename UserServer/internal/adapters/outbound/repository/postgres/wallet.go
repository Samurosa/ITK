package postgres

import (
	"ITK_Code/m/v2/internal/core/dto"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BalanceRepository struct {
	pool *pgxpool.Pool
}

func NewBalanceStorage(pool *pgxpool.Pool) *BalanceRepository {
	return &BalanceRepository{
		pool: pool,
	}
}

func (b *BalanceRepository) Deposit(ctx context.Context, userID string, asset string, amount dto.Money) (dto.Balance, error) {
	query := `
		INSERT INTO balances (user_id, asset, available, locked) VALUES ($1, $2, $3, 0)
		ON CONFLICT (user_id, asset)
    	DO UPDATE SET avalible = balances.available + EXCLUDED.available
		RETURNING id, user_id, asset, available, locked
	`

	var balance dto.Balance

	err := b.pool.QueryRow(ctx,
		query,
		userID,
		asset,
		amount,
	).Scan(
		&balance.ID,
		&balance.UserID,
		&balance.Asset,
		&balance.Available,
		&balance.Locked,
	)
	if err != nil {
		return dto.Balance{}, err
	}

	return balance, nil
}

func (b *BalanceRepository) GetAll(ctx context.Context, userID string) ([]dto.Balance, error) {

	query := `
		SELECT id, user_id, asset, available, locked
		FROM balances
		WHERE user_id=$1
	`

	rows, err := b.pool.Query(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]dto.Balance, 0)

	for rows.Next() {

		var balance dto.Balance

		err = rows.Scan(
			&balance.UserID,
			&balance.Asset,
			&balance.Available,
			&balance.Locked,
		)

		if err != nil {
			return nil, err
		}

		result = append(
			result,
			balance,
		)
	}

	if len(result) == 0 {
		return result, nil
	}

	return result, nil
}
