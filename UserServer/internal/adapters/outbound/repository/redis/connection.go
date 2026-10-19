package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	client *redis.Client
}

func NewStorage(client *redis.Client) *Storage {
	return &Storage{client: client}
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
