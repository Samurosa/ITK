package contextStatus

import (
	"context"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func New() *App {

	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (a *App) GetContext() context.Context {
	return a.ctx
}

func (a *App) Stop() {
	a.cancel()
}
