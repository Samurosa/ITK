package redis

import (
	"ITK_Code/m/v2/config"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	limit  int
	timer  time.Duration
	client *redis.Client
}

func NewLimiter(cfg config.Limiter, client *redis.Client) *Limiter {
	return &Limiter{
		limit:  cfg.Limit,
		timer:  cfg.Timer,
		client: client,
	}
}

func (l *Limiter) Allow(
	ctx context.Context,
	key string,
) (bool, error) {

	count, err := l.client.Incr(ctx, key).Result()

	if err != nil {
		return false, err
	}

	if count == 1 {
		l.client.Expire(
			ctx,
			key,
			l.timer,
		)
	}

	return count <= int64(l.limit), nil
}
