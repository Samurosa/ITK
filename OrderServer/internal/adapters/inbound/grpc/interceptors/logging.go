package interceptors

import (
	"context"
	"time"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RequestLog() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		started := time.Now()
		logging.FromContext(ctx).Debug("grpc request started")
		response, err := handler(ctx, req)
		logCompletion(ctx, "grpc request completed", started, err)
		return response, err
	}
}

func StreamLog() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		started := time.Now()
		logging.FromContext(stream.Context()).Debug("grpc stream started")
		err := handler(srv, stream)
		logCompletion(stream.Context(), "grpc stream completed", started, err)
		return err
	}
}

func logCompletion(ctx context.Context, message string, started time.Time, err error) {
	code := status.Code(err)
	if code == codes.Unknown {
		code = status.Code(status.FromContextError(err).Err())
	}
	logging.FromContext(ctx).Debug(message,
		zap.String("grpc_code", code.String()),
		zap.Duration("duration", time.Since(started)),
		zap.Error(err),
	)
}
