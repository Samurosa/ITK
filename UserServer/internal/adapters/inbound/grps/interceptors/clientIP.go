package interceptors

import (
	reqCtx "ITK_Code/m/v2/internal/core/context"
	"ITK_Code/m/v2/internal/core/dto"
	"context"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func ClientIPInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (
		interface{},
		error,
	) {
		log = log.Named("Rate limiting interceptor")
		value, ok := peer.FromContext(ctx)
		if !ok {
			log.Error("No peer info found", zap.String("value", value.String()))
			return nil, status.Error(codes.NotFound, "Client IP address is missing")
		}

		host, _, err := net.SplitHostPort(value.Addr.String())
		if err != nil {
			log.Error("failed to split host and port", zap.Error(err))
			return nil, status.Error(codes.InvalidArgument, "Client IP address is invalid")
		}

		ctx, err = reqCtx.UpdateRequestContext(ctx,
			func(baseContext *dto.RequestContext) {
				baseContext.Metadata.ClientIP = host
			})
		if err != nil {
			log.Info("failed to update request context", zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to update request context")
		}

		return handler(ctx, req)
	}
}
