package redis

import (
	"ITK_Code/m/v2/internal/config"
	requestContext "ITK_Code/m/v2/internal/core/context"
	coreErorrs "ITK_Code/m/v2/internal/core/errors"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Limiter struct {
	log      *zap.Logger
	client   *redis.Client
	capacity float64
	timer    time.Duration
}

func NewLimiter(log *zap.Logger, cfg config.Limiter, client *redis.Client) *Limiter {
	return &Limiter{
		log:      log,
		capacity: cfg.Capacity,
		timer:    cfg.Timer,
		client:   client,
	}
}

func (l *Limiter) Allow(
	ctx context.Context,
) (bool, error) {
	log := l.log.Named("limiter allow")

	requestCtx, err := requestContext.GetRequestContext(ctx)
	if err != nil {
		log.Error("Failed to get user ip from context", zap.Error(err))
		return false, coreErorrs.ErrInvalidContext
	}
	ip := requestCtx.Metadata.ClientIP
	//TODO: сделать девайс+ip

	pipe := l.client.TxPipeline()

	count, err := pipe.Incr(ctx, ip).Result()

	if err != nil {
		return false, err
	}

	if count == 1 {
		if err = pipe.Expire(
			ctx,
			ip,
			l.timer,
		).Err(); err != nil {
			return false, err
		}
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	return count <= int64(l.capacity), nil
}
