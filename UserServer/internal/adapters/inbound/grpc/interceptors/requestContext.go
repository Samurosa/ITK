package interceptors

import (
	requestContext "ITK_Code/m/v2/internal/core/context"
	"ITK_Code/m/v2/internal/core/dto"
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func RequestContextInterceptor(
	log *zap.Logger,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log := log.Named("RequestContextInterceptor")
		requestCtx := dto.RequestContext{}

		ctx = requestContext.WithRequestContext(
			ctx,
			requestCtx,
		)
		log.Debug("made request context")

		return handler(ctx, req)
	}
}
