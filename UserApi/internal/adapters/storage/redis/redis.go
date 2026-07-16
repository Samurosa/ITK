package redis

import (
	"ITK_Code/m/v2/config"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	client *redis.Client
}

func NewStorage(client *redis.Client) *Storage {
	return &Storage{client: client}
}

func NewRedisClient(ctx context.Context, cfg config.Redis) (*redis.Client, error) {

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
		return nil, err
	}

	return client, nil
}

func (s *Storage) Stop() error {
	return s.client.Close()
}

func (s *Storage) AcquireRefreshLock(
	ctx context.Context,
	jti string,
) (bool, error) {

	key := "lock:refresh:" + jti

	ok, err := s.client.SetNX(
		ctx,
		key,
		"locked",
		5*time.Second,
	).Result()

	return ok, err
}

func (s *Storage) ReleaseRefreshLock(
	ctx context.Context,
	userID string,
) error {

	return s.client.Del(
		ctx,
		"lock:refresh:"+userID,
	).Err()
}
