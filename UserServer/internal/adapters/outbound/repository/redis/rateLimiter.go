package redis

import (
	"ITK_Code/m/v2/internal/config"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Limiter struct {
	log      *zap.Logger
	client   *redis.Client
	capacity int64
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
	ip string,
	deviceID string,
) (bool, error) {
	log := l.log.Named("limiter allow")

	key := "rate-limiter ip:" + ip + "deviceID:" + deviceID

	addLimiterByKeyScript := redis.NewScript(`
	local rate = redis.call("INCR", KEYS[1])
	
	if rate == 1 then
		redis.call("EXPIRE", KEYS[1], ARGV[1])
	end

	return rate
`)

	result, err := addLimiterByKeyScript.Run(ctx,
		l.client,
		[]string{key},
		int(l.timer.Seconds()),
	).Result()

	if err != nil {
		log.Error("Failed to add rate limiter", zap.Error(err))
		return false, err
	}
	count, ok := result.(int64)
	if !ok {
		log.Error("Failed to get rate limiter count", zap.Error(err))
		return false, err
	}

	return count <= l.capacity, nil
}
