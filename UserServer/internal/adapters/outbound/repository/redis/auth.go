package redis

import (
	"ITK_Code/m/v2/internal/core/auth"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *Storage) Create(ctx context.Context, jti string, sessionModel auth.SessionModel) error {
	return s.toRedisSave(ctx, jti, &sessionModel)
}

func (s *Storage) GetByJTI(ctx context.Context, jti string) (auth.SessionModel, error) {
	return s.fromRedisByJTI(ctx, jti)
}

func (s *Storage) Update(ctx context.Context, storedJTI string, jti string, sessionModel auth.SessionModel) error {
	return s.toRedisUpdate(ctx, storedJTI, jti, &sessionModel)
}

func (s *Storage) DeleteByJTI(ctx context.Context, jti string, userID string) error {
	return s.deleteFromRedisByJTI(ctx, jti, userID)
}

func (s *Storage) DeleteByUser(ctx context.Context, userID string) error {
	return s.deleteFromRedisByUser(ctx, userID)
}

func (s *Storage) toRedisSave(ctx context.Context, jti string, model *auth.SessionModel) error {
	key := "session:" + jti

	setter := func(p redis.Pipeliner) error {
		fieldsSave(p, ctx, model, key, jti)
		return nil
	}

	_, err := s.client.TxPipelined(ctx, setter)
	return err
}

func (s *Storage) fromRedisByJTI(ctx context.Context, jti string) (auth.SessionModel, error) {
	key := "session:" + jti

	data, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		return auth.SessionModel{}, err
	}

	if len(data) == 0 {
		return auth.SessionModel{}, auth.ErrSessionNotFound
	}

	expiresAt, err := time.Parse(
		time.RFC3339Nano,
		data["expires_at"],
	)
	if err != nil {
		return auth.SessionModel{}, err
	}

	createdAt, err := time.Parse(
		time.RFC3339Nano,
		data["created_at"],
	)
	if err != nil {
		return auth.SessionModel{}, err
	}

	ttl := time.Until(expiresAt)

	if ttl <= 0 {
		return auth.SessionModel{}, auth.ErrSessionNotFound
	}

	return auth.SessionModel{
		UserID:           data["user_id"],
		DeviceID:         data["device_id"],
		RefreshTokenHash: data["refresh_token_hash"],
		TTL:              ttl,
		ExpiresAt:        expiresAt,
		CreatedAt:        createdAt,
	}, nil
}

func (s *Storage) toRedisUpdate(ctx context.Context, storedJTI string, jti string, model *auth.SessionModel) error {
	result, err := redis.NewScript(`
		if redis.call("EXISTS", KEYS[1]) == 0 then
			return 0
		end
		redis.call("HSET", KEYS[2],
			"user_id", ARGV[1],
			"device_id", ARGV[2],
			"refresh_token_hash", ARGV[3],
			"created_at", ARGV[4],
			"expires_at", ARGV[5])
		redis.call("PEXPIRE", KEYS[2], ARGV[6])
		redis.call("SADD", KEYS[3], KEYS[2])
		redis.call("DEL", KEYS[1])
		redis.call("SREM", KEYS[3], KEYS[1])
		return 1
	`).Run(ctx, s.client,
		[]string{"session:" + storedJTI, "session:" + jti, "user:" + model.UserID},
		model.UserID,
		model.DeviceID,
		model.RefreshTokenHash,
		model.CreatedAt.UTC().Format(time.RFC3339Nano),
		model.ExpiresAt.UTC().Format(time.RFC3339Nano),
		model.TTL.Milliseconds(),
	).Int64()
	if err != nil {
		return err
	}
	if result == 0 {
		return auth.ErrSessionNotFound
	}
	return nil
}

func (s *Storage) deleteFromRedisByJTI(
	ctx context.Context,
	jti string,
	userID string,
) error {
	key := "session:" + jti
	pipe := s.client.TxPipeline()

	pipe.Del(ctx, key)
	pipe.SRem(ctx, "user:"+userID, key)

	_, err := pipe.Exec(ctx)
	return err
}

func (s *Storage) deleteFromRedisByUser(
	ctx context.Context,
	userID string,
) error {
	key := "user:" + userID

	return redis.NewScript(`
		local sessions = redis.call("SMEMBERS", KEYS[1])
		for _, session in ipairs(sessions) do
			redis.call("UNLINK", session)
		end
		redis.call("DEL", KEYS[1])
		return #sessions
	`).Run(ctx, s.client, []string{key}).Err()
}

func fieldsSave(p redis.Pipeliner, ctx context.Context, model *auth.SessionModel, key string, jti string) {
	fields := map[string]interface{}{
		"user_id":            model.UserID,
		"device_id":          model.DeviceID,
		"refresh_token_hash": model.RefreshTokenHash,
		"created_at":         model.CreatedAt.UTC().Format(time.RFC3339Nano),
		"expires_at":         model.ExpiresAt.UTC().Format(time.RFC3339Nano),
	}

	p.HSet(ctx, key, fields)
	p.Expire(ctx, key, model.TTL)
	p.SAdd(ctx, "user:"+model.UserID, "session:"+jti)
}
