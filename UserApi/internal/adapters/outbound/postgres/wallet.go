package postgres

import (
	"ITK_Code/m/v2/internal/core/wallet"
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

func (b *BalanceRepository) Create(ctx context.Context, userID string, currency string) (wallet.Balance, error) {

	query := `
		INSERT INTO balances
		(
			user_id,
			asset,
			available,
			locked
		)
		VALUES ($1,$2,$3,$4)
		RETURNING user_id, asset, available, locked
	`

	var balance wallet.Balance

	err := b.pool.QueryRow(
		ctx,
		query,
		userID,
		currency,
		0,
		0,
	).Scan(
		&balance.UserID,
		&balance.Asset,
		&balance.Available,
		&balance.Locked,
	)

	if err != nil {
		return wallet.Balance{}, err
	}

	return balance, nil
}

func (b *BalanceRepository) Get(ctx context.Context, userID string, currency string) (wallet.Balance, error) {

	query := `
		SELECT
			user_id,
			asset,
			available,
			locked
		FROM balances
		WHERE user_id=$1
		AND asset=$2
	`

	var balance wallet.Balance

	err := b.pool.QueryRow(
		ctx,
		query,
		userID,
		currency,
	).Scan(
		&balance.UserID,
		&balance.Asset,
		&balance.Available,
		&balance.Locked,
	)

	if err != nil {
		return wallet.Balance{}, err
	}

	return balance, nil
}

func (b *BalanceRepository) GetOrCreate(ctx context.Context, userID string, currency string) (wallet.Balance, error) {

	query := `
		INSERT INTO balances
		(
			user_id,
			asset
		)
		VALUES($1,$2)
		ON CONFLICT(user_id,asset)
		DO NOTHING
	`

	_, err := b.pool.Exec(
		ctx,
		query,
		userID,
		currency,
	)

	if err != nil {
		return wallet.Balance{}, err
	}

	return b.Get(
		ctx,
		userID,
		currency,
	)
}

func (b *BalanceRepository) Save(ctx context.Context, balance wallet.Balance) error {

	query := `
		UPDATE balances
		SET
			available=$1,
			locked=$2
		WHERE user_id=$3
		AND asset=$4
	`

	result, err := b.pool.Exec(
		ctx,
		query,
		balance.Available,
		balance.Locked,
		balance.UserID,
		balance.Asset,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return wallet.ErrBalanceNotFound
	}

	return nil
}

func (b *BalanceRepository) GetAll(ctx context.Context, userID string) ([]wallet.Balance, error) {

	query := `
		SELECT
			user_id,
			asset,
			available,
			locked
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

	result := make([]wallet.Balance, 0)

	for rows.Next() {

		var balance wallet.Balance

		err := rows.Scan(
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
		return result, err
	}

	return result, nil
}
