package redis

import (
	"ITK_Code/m/v2/internal/core/corerrors"
	"context"
	"errors"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"go.uber.org/zap"
)

func (s *Storage) Validate(ctx context.Context, jti string) error {
	key := "session:" + jti

	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		log := logging.FromContext(ctx).Named("SessionValidator")
		if errors.Is(err, context.Canceled) {
			log.Debug("session lookup canceled", zap.Error(err))
		} else if errors.Is(err, context.DeadlineExceeded) {
			log.Warn("session lookup timed out", zap.Error(err))
		} else {
			log.Error("failed to look up session in redis", zap.Error(err))
		}
		return err
	}

	if exists == 0 {
		return corerrors.ErrSessionNotFound
	}

	return nil
}
