package sessionValidator

import (
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	client *redis.Client
}

func NewStorage(client *redis.Client) *Storage {
	return &Storage{client: client}
}

func (s *Storage) Close() error {
	return s.client.Close()
}
