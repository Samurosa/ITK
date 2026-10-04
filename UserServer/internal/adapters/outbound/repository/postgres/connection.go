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

func NewStorage(ctx context.Context, logger *zap.Logger, cfg config.Postgres) (*Storage, error) {
	log := logger.Named("postgres")

	connectionSettings, err := pgxpool.ParseConfig(cfg.Link)
	if err != nil {
		// Parse errors may include the DSN with database credentials.
		log.Error("invalid postgres connection configuration")
		return nil, err
	}
	connectionSettings.MaxConns = cfg.MaxConns
	connectionSettings.MinConns = cfg.MinConns

	for i := 1; i <= cfg.MaxRetries; i++ {

		pool, err := pgxpool.NewWithConfig(ctx, connectionSettings)
		if err != nil {
			log.Warn("postgres pool creation attempt failed", zap.Int("attempt", i), zap.Int("max_attempts", cfg.MaxRetries), zap.Error(err))

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

		log.Warn("postgres connection attempt failed", zap.Int("attempt", i), zap.Int("max_attempts", cfg.MaxRetries), zap.Error(err))

		pool.Close()

		time.Sleep(
			time.Duration(i) * time.Second,
		)
	}
	log.Error("postgres connection attempts exhausted", zap.Int("max_attempts", cfg.MaxRetries))
	return nil, ErrPingDB
}

func (s *Storage) GetPool() *pgxpool.Pool {
	return s.pool
}

func (s *Storage) ClosePool() {
	s.pool.Close()
}
