package workers

import (
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/workers"
	"context"

	"go.uber.org/zap"
)

type App struct {
	log               *zap.Logger
	ctx               context.Context
	sessionRepository auth.SessionRepository
}

func NewWorker(log *zap.Logger, ctx context.Context, sessionRepository auth.SessionRepository) *App {
	return &App{
		log:               log,
		ctx:               ctx,
		sessionRepository: sessionRepository,
	}
}

func (w *App) Run() {
	cleaner := workers.NewExpiredSessionCleaner(w.log, w.sessionRepository)

	cleaner.Clean(w.ctx)
}
