package redis

import (
	"ITK_Code/m/v2/internal/config"
	context2 "ITK_Code/m/v2/internal/core/context"
	"ITK_Code/m/v2/internal/core/errors"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Limiter struct {
	log    *zap.Logger
	limit  int
	timer  time.Duration
	client *redis.Client
}

func NewLimiter(log *zap.Logger, cfg config.Limiter, client *redis.Client) *Limiter {
	return &Limiter{
		log:    log,
		limit:  cfg.Limit,
		timer:  cfg.Timer,
		client: client,
	}
}

func (l *Limiter) Allow(
	ctx context.Context,
) (bool, error) {
	log := l.log.Named("limiter allow")

	ip, err := context2.GetClientIPFromContext(ctx)
	if err != nil {
		log.Error("Failed to get user ip from context", zap.Error(err))
		return false, errors.ErrInvalidContext
	}

	count, err := l.client.Incr(ctx, ip).Result()

	if err != nil {
		return false, err
	}

	if count == 1 {
		if err := l.client.Expire(
			ctx,
			ip,
			l.timer,
		).Err(); err != nil {
			return false, err
		}
	}

	return count <= int64(l.limit), nil
}
