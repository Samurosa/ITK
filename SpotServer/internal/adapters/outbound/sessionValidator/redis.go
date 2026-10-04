package sessionValidator

import (
	"context"
	"errors"

	"ITK_Code/m/v2/internal/config"

	"github.com/redis/go-redis/v9"
)

var ErrSessionNotFound = errors.New("session not found")

type Storage struct {
	client *redis.Client
}

func NewStorage(ctx context.Context, cfg config.Redis) (*Storage, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Addr, Password: cfg.Password, DB: cfg.DB,
		MaxRetries: cfg.MaxRetries, DialTimeout: cfg.DialTimeout,
		ReadTimeout: cfg.Timeout, WriteTimeout: cfg.Timeout,
		ContextTimeoutEnabled: true,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}
	return &Storage{client: client}, nil
}

func (s *Storage) Validate(ctx context.Context, jti string) error {
	exists, err := s.client.Exists(ctx, "session:"+jti).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return ErrSessionNotFound
	}
	return nil
}

func (s *Storage) Close() error {
	return s.client.Close()
}
