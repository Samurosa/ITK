package infrastructure

import (
	"ITK_Code/m/v2/internal/adapters/outbound/repository/postgres"
	"ITK_Code/m/v2/internal/adapters/outbound/repository/redis"
	"ITK_Code/m/v2/internal/config"
	"context"

	"go.uber.org/zap"
)

type Storages struct {
	Postgres *postgres.Storage
	Redis    *redis.Storage
}

func NewStorages(ctx context.Context,
	logger *zap.Logger,
	cfg *config.Config,
) (
	*Storages,
	error,
) {
	log := logger.Named("postgres start")

	postgresStorage, err := postgres.NewStorage(ctx,
		log,
		cfg.Postgres.Link,
		cfg.Postgres.MaxRetries,
	)
	if err != nil {
		log.Error("Error starting postgres storage", zap.Error(err))
		return nil, err
	}
	log.Info("Successfully started postgres storage")

	log.Named("Redis infrastructure")

	redisClient, err := redis.NewRedisClient(ctx, log, cfg.Redis)
	if err != nil {
		log.Error("Error creating redis client", zap.Error(err))
		return nil, err
	}

	storage := redis.NewStorage(redisClient)

	return &Storages{
		Postgres: postgresStorage,
		Redis:    storage,
	}, nil
}
