package postgres

import (
	"context"
	"time"

	"ITK_Code/m/v2/internal/core/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxStorage struct {
	pool *pgxpool.Pool
}

func NewOutboxStorage(pool *pgxpool.Pool) *OutboxStorage {
	return &OutboxStorage{pool: pool}
}

func (s *OutboxStorage) Claim(
	ctx context.Context,
	workerID string,
	batchSize int,
	leaseDuration time.Duration,
) ([]models.OutboxMessage, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	query := `
		WITH selected AS (
			SELECT id
			FROM order_outbox
			WHERE published_at IS NULL
				AND next_attempt_at <= NOW()
				AND (
					locked_until IS NULL
					OR locked_until < NOW()
				)
			ORDER BY created_at, id
			FOR UPDATE SKIP LOCKED
			LIMIT $1
		)
		UPDATE order_outbox AS o
		SET
			locked_by = $2,
			locked_until = NOW() + $3::interval
		FROM selected
		WHERE o.id = selected.id
		RETURNING
			o.id::text,
			o.order_id::text,
			o.message_type,
			o.payload,
			o.attempts,
			o.created_at
	`

	rows, err := tx.Query(ctx, query, batchSize, workerID, leaseDuration.String())
	if err != nil {
		return nil, err
	}

	messages := make([]models.OutboxMessage, 0, batchSize)

	for rows.Next() {
		var message models.OutboxMessage

		if err := rows.Scan(
			&message.ID,
			&message.OrderID,
			&message.MessageType,
			&message.Payload,
			&message.Attempts,
			&message.CreatedAt,
		); err != nil {
			rows.Close()
			return nil, err
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *OutboxStorage) MarkPublished(
	ctx context.Context,
	messageID string,
	workerID string,
) error {
	query := `
		UPDATE order_outbox
		SET
			published_at = NOW(),
			locked_by = NULL,
			locked_until = NULL,
			last_error = NULL
		WHERE id = $1
			AND locked_by = $2
			AND published_at IS NULL
	`

	_, err := s.pool.Exec(ctx, query, messageID, workerID)
	return err
}

func (s *OutboxStorage) MarkFailed(
	ctx context.Context,
	messageID string,
	workerID string,
	messageErr string,
	retryAfter time.Duration,
) error {
	query := `
		UPDATE order_outbox
		SET
			attempts = attempts + 1,
			last_error = $3,
			next_attempt_at = NOW() + $4::interval,
			locked_by = NULL,
			locked_until = NULL
		WHERE id = $1
			AND locked_by = $2
			AND published_at IS NULL
	`

	_, err := s.pool.Exec(
		ctx,
		query,
		messageID,
		workerID,
		messageErr,
		retryAfter.String(),
	)

	return err
}
