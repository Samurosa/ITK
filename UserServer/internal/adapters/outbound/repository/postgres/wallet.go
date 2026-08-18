package postgres

import (
	"ITK_Code/m/v2/internal/core/wallet"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
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

func (b *BalanceRepository) Deposit(ctx context.Context, userID string, asset string, amount wallet.Money, idempotentKey string) (wallet.Balance, error) {
	tx, err := b.pool.Begin(ctx)
	if err != nil {
		return wallet.Balance{}, err
	}
	defer tx.Rollback(ctx)

	var balance wallet.Balance
	var idWalletOperation string

	err = tx.QueryRow(ctx,
		`
		INSERT INTO wallet_operations (idempotency_key, user_id, asset, amount) VALUES ($1, $2, $3, $4)
		ON CONFLICT (idempotency_key)
    	DO NOTHING 
		RETURNING id
	`,
		idempotentKey,
		userID,
		asset,
		amount.Amount,
	).Scan(
		&idWalletOperation,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		err = tx.QueryRow(ctx,
			`
			SELECT id, user_id, asset, available, locked FROM balances
			WHERE user_id = $1 AND asset = $2
			`,
			userID,
			asset,
		).Scan(
			&balance.ID,
			&balance.UserID,
			&balance.Asset,
			&balance.Available,
			&balance.Locked,
		)

		return balance, nil
	}

	if err != nil {
		return wallet.Balance{}, err
	}

	err = tx.QueryRow(ctx,
		`
		INSERT INTO balances (user_id, asset, available, locked) VALUES ($1, $2, $3, 0)
		ON CONFLICT (user_id, asset)
    	DO UPDATE SET available = balances.available + EXCLUDED.available
		RETURNING id, user_id, asset, available, locked
	`,
		userID,
		asset,
		amount.Amount,
	).Scan(
		&balance.ID,
		&balance.UserID,
		&balance.Asset,
		&balance.Available,
		&balance.Locked,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return wallet.Balance{}, wallet.ErrBalanceNotFound
	}

	if err != nil {
		return wallet.Balance{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return wallet.Balance{}, err
	}

	return balance, nil
}

func (b *BalanceRepository) GetAll(ctx context.Context, userID string) ([]wallet.Balance, error) {

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

	if errors.Is(err, pgx.ErrNoRows) {
		return []wallet.Balance{}, wallet.ErrBalanceNotFound
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]wallet.Balance, 0)

	for rows.Next() {

		var balance wallet.Balance

		err = rows.Scan(
			&balance.ID,
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
