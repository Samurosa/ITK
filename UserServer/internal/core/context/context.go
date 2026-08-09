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

func UserID(ctx context.Context) (string, error) {
	requestCtx, ok := ctx.Value(
		requestContextKey{},
	).(dto.RequestContext)

	if !ok {
		return "", errors.ErrInvalidContext
	}

	return requestCtx.Principal.UserID, nil
}

func JTI(ctx context.Context) (string, error) {
	requestCtx, ok := ctx.Value(
		requestContextKey{},
	).(dto.RequestContext)

	if !ok {
		return "", errors.ErrInvalidContext
	}

	return requestCtx.JTI, nil
}

func UpdateRequestContext(
	ctx context.Context,
	update func(*dto.RequestContext),
) (
	context.Context,
	error,
) {

	requestCtx, ok := ctx.Value(requestContextKey{}).(dto.RequestContext)

	if !ok {
		return ctx, errors.ErrCtxForUpdateNotFound
	}

	update(&requestCtx)

	return context.WithValue(
		ctx,
		requestContextKey{},
		requestCtx,
	), nil
}
