package context

import (
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/errors"
	"context"
)

type requestContextKey struct {
}

func WithRequestContext(
	ctx context.Context,
	requestCtx dto.RequestContext,
) context.Context {
	return context.WithValue(
		ctx,
		requestContextKey{},
		requestCtx,
	)
}

func GetRequestContext(
	ctx context.Context,
) (dto.RequestContext, error) {

	requestCtx, ok := ctx.Value(
		requestContextKey{},
	).(dto.RequestContext)

	if !ok {
		return dto.RequestContext{}, errors.ErrInvalidContext
	}

	return requestCtx, nil
}

func UpdateRequestContext(
	ctx context.Context,
	update func(*dto.RequestContext),
) context.Context {

	requestCtx, ok := ctx.Value(requestContextKey{}).(dto.RequestContext)

	if !ok {
		requestCtx = dto.RequestContext{}
	}

	update(&requestCtx)

	return context.WithValue(
		ctx,
		requestContextKey{},
		requestCtx,
	)
}
