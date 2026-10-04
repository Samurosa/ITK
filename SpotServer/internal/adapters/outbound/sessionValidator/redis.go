package sessionValidator

import (
	"context"
	"errors"

	"ITK_Code/m/v2/internal/config"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
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
		log := logging.FromContext(ctx).Named("session.validate")
		switch {
		case errors.Is(err, context.Canceled):
			log.Debug("redis session lookup canceled", zap.Error(err))
		case errors.Is(err, context.DeadlineExceeded):
			log.Warn("redis session lookup timed out", zap.Error(err))
		default:
			log.Error("redis session lookup failed", zap.Error(err))
		}
		return err
	}
	if exists == 0 {
		logging.FromContext(ctx).Debug("session not found")
		return ErrSessionNotFound
	}
	return nil
}

func (s *Storage) Close() error {
	return s.client.Close()
}
