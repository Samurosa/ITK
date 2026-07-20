package infrastructure

import (
	"context"
)

type Context struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewContext() *Context {

	ctx, cancel := context.WithCancel(context.Background())
	return &Context{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (a *Context) GetContext() context.Context {
	return a.ctx
}

func (a *Context) Stop() {
	a.cancel()
}
