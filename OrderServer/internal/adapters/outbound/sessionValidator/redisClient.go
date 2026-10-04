package sessionValidator

import (
	"ITK_Code/m/v2/internal/config"
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
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

func NewRedisClient(ctx context.Context, cfg config.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:                  cfg.Addr,
		Password:              cfg.Password,
		DB:                    cfg.DB,
		MaxRetries:            cfg.MaxRetries,
		DialTimeout:           cfg.DialTimeout,
		ReadTimeout:           cfg.Timeout,
		WriteTimeout:          cfg.Timeout,
		ContextTimeoutEnabled: true,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("%w: %w", ErrPingToRedis, err)
	}

	return client, nil
}
