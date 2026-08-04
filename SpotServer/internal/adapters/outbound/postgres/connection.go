package postgres

import (
	"ITK_Code/m/v2/internal/config"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	ErrPingDB = errors.New("ping db failed")
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(ctx context.Context, logger *zap.Logger, postgres config.Postgres) (*Storage, error) {
	log := logger.Named("postgres")

	configPool, err := pgxpool.ParseConfig(postgres.Link)
	if err != nil {
		return nil, err
	}
	configPool.MaxConns = 10
	configPool.MinConns = 2

	for i := 1; i <= postgres.MaxRetries; i++ {

		pool, err := pgxpool.NewWithConfig(ctx, configPool)
		if err != nil {
			log.Error("create postgres pool failed", zap.Error(err))

			time.Sleep(time.Duration(i) * time.Second)
			continue
		}

		pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		err = pool.Ping(pingCtx)
		cancel()

		if err == nil {
			log.Info("postgres connected")

			return &Storage{
				pool: pool,
			}, nil
		}

		log.Error("postgres ping failed", zap.Error(err))

		pool.Close()

		time.Sleep(
			time.Duration(i) * time.Second,
		)
	}
	return nil, ErrPingDB
}

func (s *Storage) GetPool() *pgxpool.Pool {
	return s.pool
}

func (s *Storage) ClosePool() {
	s.pool.Close()
}
