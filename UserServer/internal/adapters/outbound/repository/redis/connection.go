package redis

import (
	"ITK_Code/m/v2/internal/config"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
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
	log = log.Named("Redis outbound adapter")
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
		log.Error("Failed to connect to Redis", zap.Error(err))
		_ = client.Close()
		return nil, ErrPingToRedis
	}
	log.Info("Redis connected")

	return client, nil
}

func (s *Storage) GetClient() *redis.Client {
	return s.client
}

func (s *Storage) Stop() error {
	return s.client.Close()
}

func (s *Storage) AcquireRefreshLock(
	ctx context.Context,
	jti string,
) (string, error) {

	key := "lock:refresh:" + jti
	owner := uuid.NewString()

	ok, err := s.client.SetNX(
		ctx,
		key,
		owner,
		30*time.Second,
	).Result()
	if err != nil {
		return "", err
	}
	if !ok {
		return "", nil
	}
	return owner, nil
}

func (s *Storage) ReleaseRefreshLock(
	ctx context.Context,
	jti string,
	owner string,
) error {
	return redis.NewScript(`
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		end
		return 0
	`).Run(ctx, s.client, []string{"lock:refresh:" + jti}, owner).Err()
}
