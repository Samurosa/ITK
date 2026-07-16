package redis

import (
	"ITK_Code/m/v2/internal/core/auth"
	"context"
	"errors"
	"reflect"
	"strconv"
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

func (s *Storage) toRedisSave(ctx context.Context, jti string, value *auth.SessionModel) error {
	key := "session:" + jti

	val := reflect.ValueOf(value).Elem()

	if val.NumField() == 0 {
		return errors.New("value is empty")
	}

	setter := func(p redis.Pipeliner) error {

		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)

			tag := field.Tag.Get("redis")

			if err := p.HSet(ctx, key, tag, val.Field(i).Interface()).Err(); err != nil {
				return err
			}

		}
		if err := p.Expire(ctx, key, value.TTL).Err(); err != nil {
			return err
		}
		if err := p.SAdd(ctx, "user:"+value.UserID, "session:"+jti).Err(); err != nil {
			return err
		}

		return nil
	}
	if _, err := s.client.TxPipelined(ctx, setter); err != nil {
		return err
	}
	return nil
}

func (s *Storage) fromRedisByJTI(ctx context.Context, key string) (auth.SessionModel, error) {
	key = "session:" + key

	data, err := s.client.HGetAll(
		ctx,
		key,
	).Result()

	if err != nil {
		return auth.SessionModel{}, err
	}

	if len(data) == 0 {
		return auth.SessionModel{}, redis.Nil
	}

	expiresAt, err := time.Parse(
		time.RFC3339,
		data["expires_at"],
	)
	if err != nil {
		return auth.SessionModel{}, err
	}

	createsAt, err := time.Parse(
		time.RFC3339,
		data["created_at"],
	)
	if err != nil {
		return auth.SessionModel{}, err
	}

	ns, err := strconv.ParseInt(data["ttl"], 10, 64)
	if err != nil {
		return auth.SessionModel{}, err
	}

	duration := time.Duration(ns)

	return auth.SessionModel{
		UserID:           data["user_id"],
		DeviceID:         data["device_id"],
		RefreshTokenHash: data["refresh_token_hash"],
		TTL:              duration,
		ExpiresAt:        expiresAt,
		CreatedAt:        createsAt,
	}, nil
}

func (s *Storage) toRedisUpdate(ctx context.Context, storedJTI string, jti string, value *auth.SessionModel) error {
	key := "session:" + jti

	val := reflect.ValueOf(value).Elem()

	if val.NumField() == 0 {
		return errors.New("value is empty")
	}

	setter := func(p redis.Pipeliner) error {

		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)

			tag := field.Tag.Get("redis")

			if err := p.HSet(ctx, key, tag, val.Field(i).Interface()).Err(); err != nil {
				return err
			}

		}
		if err := p.Expire(ctx, key, value.TTL).Err(); err != nil {
			return err
		}
		if err := p.SAdd(ctx, "user:"+value.UserID, "session:"+jti).Err(); err != nil {
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
			"user:"+value.UserID,
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

	return err
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
		return errors.New("tokens not found")
	}

	pipe := s.client.TxPipeline()

	for _, tokenJTI := range tokensJTI {
		if err = pipe.Del(
			ctx,
			tokenJTI,
		).Err(); err != nil {
			return err
		}
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
