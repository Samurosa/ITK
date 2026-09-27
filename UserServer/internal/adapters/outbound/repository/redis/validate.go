package redis

import (
	"ITK_Code/m/v2/internal/core/corerrors"
	"context"
)

func (s *Storage) Validate(ctx context.Context, jti string) error {
	key := "session:" + jti

	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return err
	}

	if exists == 0 {
		return corerrors.ErrSessionNotFound
	}

	return nil
}
