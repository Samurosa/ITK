package sessionValidator

import (
	"ITK_Code/m/v2/internal/core/coreErrors"
	"context"
	"log"
)

func (s *Storage) Validate(ctx context.Context, jti string) error {
	key := "session:" + jti

	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		log.Printf("REDIS EXISTS ERROR: %v", err)
		return err
	}

	if exists == 0 {
		return coreErrors.ErrSessionNotFound
	}

	return nil
}
