package redis

import (
	"ITK_Code/m/v2/internal/core/auth"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *Storage) Create(ctx context.Context, jti string, sessionModel auth.SessionModel) error {

	if err := s.toRedisSave(ctx, jti, &sessionModel); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetByJTI(ctx context.Context, jti string) (auth.SessionModel, error) {
	session, err := s.fromRedisByJTI(ctx, jti)
	if err != nil {
		return auth.SessionModel{}, err
	}

	return session, nil
}

func (s *Storage) Update(ctx context.Context, storedJTI string, jti string, sessionModel auth.SessionModel) error {
	if err := s.toRedisUpdate(ctx, storedJTI, jti, &sessionModel); err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteByJTI(ctx context.Context, jti string, userID string) error {
	err := s.deleteFromRedisByJTI(ctx, jti, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteByUser(ctx context.Context, userID string) error {
	if err := s.deleteFromRedisByUser(ctx, userID); err != nil {
		return err
	}

	return nil
}

func (s *Storage) toRedisSave(ctx context.Context, jti string, model *auth.SessionModel) error {
	key := "session:" + jti

	setter := func(p redis.Pipeliner) error {

		fields := map[string]interface{}{
			"user_id":            model.UserID,
			"device_id":          model.DeviceID,
			"refresh_token_hash": model.RefreshTokenHash,
			"created_at":         model.CreatedAt.UTC().Format(time.RFC3339Nano),
			"expires_at":         model.ExpiresAt.UTC().Format(time.RFC3339Nano),
		}

		if err := p.HSet(ctx, key, fields).Err(); err != nil {
			return err
		}

		if err := p.Expire(ctx, key, model.TTL).Err(); err != nil {
			return err
		}
		if err := p.SAdd(ctx, "user:"+model.UserID, "session:"+jti).Err(); err != nil {
			return err
		}

		return nil
	}

	if _, err := s.client.TxPipelined(ctx, setter); err != nil {
		return err
	}
	return nil
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
	key := "session:" + jti

	setter := func(p redis.Pipeliner) error {

		fields := map[string]interface{}{
			"user_id":            model.UserID,
			"device_id":          model.DeviceID,
			"refresh_token_hash": model.RefreshTokenHash,
			"created_at":         model.CreatedAt.UTC().Format(time.RFC3339Nano),
			"expires_at":         model.ExpiresAt.UTC().Format(time.RFC3339Nano),
		}

		if err := p.HSet(ctx, key, fields).Err(); err != nil {
			return err
		}

		if err := p.Expire(ctx, key, model.TTL).Err(); err != nil {
			return err
		}
		if err := p.SAdd(ctx, "user:"+model.UserID, "session:"+jti).Err(); err != nil {
			return err
		}

		restoreSession := "session:" + storedJTI

		if err := p.Del(
			ctx,
			restoreSession,
		).Err(); err != nil {
			return err
		}

		if err := p.SRem(
			ctx,
			"user:"+model.UserID,
			restoreSession,
		).Err(); err != nil {
			return err
		}

		return nil
	}
	if _, err := s.client.TxPipelined(ctx, setter); err != nil {
		return err
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

	if err := pipe.Del(
		ctx,
		key,
	).Err(); err != nil {
		return err
	}

	if err := pipe.SRem(
		ctx,
		"user:"+userID,
		key,
	).Err(); err != nil {
		return err
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) deleteFromRedisByUser(
	ctx context.Context,
	userID string,
) error {
	key := "user:" + userID

	tokensJTI, err := s.client.SMembers(
		ctx,
		key,
	).Result()

	if err != nil {
		return err
	}

	if len(tokensJTI) == 0 {
		return nil
	}

	pipe := s.client.TxPipeline()

	if err := pipe.Unlink(ctx, tokensJTI...).Err(); err != nil {
		return err
	}

	if err = pipe.Del(ctx, key).Err(); err != nil {
		return err
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
