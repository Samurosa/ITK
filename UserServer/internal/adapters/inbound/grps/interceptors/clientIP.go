package interceptors

import (
	context2 "ITK_Code/m/v2/internal/core/context"
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
			return nil, status.Error(codes.FailedPrecondition, "Client IP address is missing")
		}

		host, _, err := net.SplitHostPort(value.Addr.String())
		if err != nil {
			log.Error("failed to split host and port", zap.Error(err))
			return nil, status.Error(codes.FailedPrecondition, "Client IP address is invalid")
		}

		ctx = context2.WithClientIP(ctx, host)

		return handler(ctx, req)
	}
}
