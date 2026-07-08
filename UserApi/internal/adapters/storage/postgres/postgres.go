package postgres

import (
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
	log  *zap.Logger
}

func NewPool(ctx context.Context, logger *zap.Logger, connectionString string, maxConnections int) (*Storage, error) {
	log := logger.Named("connect to postgres")

	for i := 1; i <= maxConnections; i++ {
		pool, err := pgxpool.New(ctx, connectionString)
		if err != nil {
			time.Sleep(time.Duration(i) * time.Second)
			continue
		}

		ctxPing, cancelPing := context.WithTimeout(context.Background(), 5*time.Second)
		err = pool.Ping(ctxPing)
		cancelPing()

		if err == nil {
			return &Storage{
				pool: pool,
				log:  logger,
			}, nil
		}

		log.Error("ping to postgres failed", zap.Error(ErrPingDB))
		pool.Close()
		time.Sleep(time.Duration(i) * time.Second)

		return nil, ErrPingDB
	}
	return nil, ErrPingDB
}

func (s *Storage) GetPool() *pgxpool.Pool {
	return s.pool
}

//TODO: переделать конект к бд и логи
