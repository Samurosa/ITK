package redis

import (
	"ITK_Code/m/v2/internal/config"
	"context"
	"fmt"
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

func NewLimiter(log *zap.Logger, cfg config.Limiter, client *redis.Client) (*Limiter, error) {
	if cfg.Capacity <= 0 || cfg.Timer < time.Millisecond {
		return nil, fmt.Errorf("limiter capacity must be positive and timer at least 1ms")
	}
	return &Limiter{
		log:      log,
		capacity: cfg.Capacity,
		timer:    cfg.Timer,
		client:   client,
	}, nil
}

func (l *Limiter) Allow(
	ctx context.Context,
	ip string,
	deviceID string,
) (bool, error) {
	keys := []string{"rate-limiter:ip:" + ip}

	if deviceID != "" {
		keys = append(keys, "rate-limiter:device:"+deviceID)
	}
	return l.allow(ctx, keys)
}

func (l *Limiter) AllowPasswordChange(ctx context.Context, userID string) (bool, error) {
	return l.allow(ctx, []string{"rate-limiter:password:" + userID})
}

var limiterScript = redis.NewScript(`
	local allowed = 1
	for _, key in ipairs(KEYS) do
		local count = redis.call("INCR", key)
		if count == 1 then redis.call("PEXPIRE", key, ARGV[1]) end
		if count > tonumber(ARGV[2]) then allowed = 0 end
	end
	return allowed
`)

func (l *Limiter) allow(ctx context.Context, keys []string) (bool, error) {
	result, err := limiterScript.Run(ctx, l.client, keys, l.timer.Milliseconds(), l.capacity).Int64()
	if err != nil {
		l.log.Error("failed to apply rate limiter", zap.Error(err))
		return false, err
	}
	return result == 1, nil
}
