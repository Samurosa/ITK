package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func New(ctx context.Context, link string) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr: link,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
