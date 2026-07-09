package workers

import (
	"ITK_Code/m/v2/internal/core/auth"
	"context"
	"time"

	"go.uber.org/zap"
)

type ExpiredSessionCleaner struct {
	log               *zap.Logger
	sessionRepository auth.SessionRepository
}

func NewExpiredSessionCleaner(log *zap.Logger, sessionRepository auth.SessionRepository) *ExpiredSessionCleaner {
	return &ExpiredSessionCleaner{
		log:               log,
		sessionRepository: sessionRepository,
	}
}

func (e *ExpiredSessionCleaner) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			err := e.sessionRepository.DeleteExpiredSessions(ctx)
			if err != nil {
				e.log.Error("Error deleting expired sessions", zap.Error(err))
			}
		}
	}
}
