package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

func NewPool(ctx context.Context, logger *zap.Logger, connectionString string) (*Storage, error) {
	log := logger.Named("connect to postgres")
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		log.Error("Failed to connect to postgres", zap.Error(err))
		return nil, err
	}
	return &Storage{
		pool: pool,
		log:  logger,
	}, nil
}
