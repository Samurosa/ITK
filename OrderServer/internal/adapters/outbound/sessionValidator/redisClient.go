package sessionValidator

import (
	"ITK_Code/m/v2/internal/config"
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	ErrPingToRedis = errors.New("ping to redis failed")
)

type Storage struct {
	client *redis.Client
}

func NewStorage(client *redis.Client) *Storage {
	return &Storage{client: client}
}

func NewRedisClient(ctx context.Context, log *zap.Logger, cfg config.Redis) (*redis.Client, error) {
	log.Named("Redis outbound adapter")
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Error("Failed to connect to Redis", zap.Error(err))
		_ = client.Close()
		return nil, ErrPingToRedis
	}
	log.Info("Redis connected")

	return client, nil
}
