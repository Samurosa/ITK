package postgres

import (
	"ITK_Code/m/v2/internal/config"
	"context"
	"errors"
	"fmt"
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
		// Parse errors may contain the connection string, including its password.
		return nil, errors.New("invalid postgres connection configuration")
	}
	if postgres.MaxRetries < 1 {
		return nil, errors.New("postgres max retries must be positive")
	}
	configPool.MaxConns = postgres.MaxConnections
	configPool.MinConns = postgres.MinConnections

	var lastErr error
	for i := 1; i <= postgres.MaxRetries; i++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		pool, err := pgxpool.NewWithConfig(ctx, configPool)
		if err == nil {
			pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err = pool.Ping(pingCtx)
			cancel()
		}

		if err == nil {
			log.Info("postgres connected", zap.Int("attempt", i))

			return &Storage{
				pool: pool,
			}, nil
		}

		lastErr = err
		if pool != nil {
			pool.Close()
		}
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		log.Warn("postgres connection attempt failed",
			zap.Int("attempt", i),
			zap.Int("max_attempts", postgres.MaxRetries),
			zap.Bool("will_retry", i < postgres.MaxRetries),
			zap.Error(err),
		)
		if i < postgres.MaxRetries {
			timer := time.NewTimer(time.Duration(i) * time.Second)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			case <-timer.C:
			}
		}
	}
	return nil, fmt.Errorf("%w after %d attempts: %w", ErrPingDB, postgres.MaxRetries, lastErr)
}

func (s *Storage) GetPool() *pgxpool.Pool {
	return s.pool
}

func (s *Storage) ClosePool() {
	s.pool.Close()
}
